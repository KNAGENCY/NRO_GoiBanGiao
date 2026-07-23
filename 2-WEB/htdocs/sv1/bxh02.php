<?php
/* ==========================================================================
   BẢNG XẾP HẠNG — TÀI SẢN & CHỈ SỐ
   Khoác giao diện Ngọc Rồng Kame. Toàn bộ cách đọc và tính giá trị giữ nguyên
   như bản cũ (xem bản sao lưu bxh02.php.cu). Chỉ đổi phần HTML/CSS xuất ra và
   bổ sung ảnh đại diện lấy từ game qua bảng head_avatar.
   ========================================================================== */

require_once __DIR__ . '/kame/layout.php';   // layout.php tự nạp data.php và config.php

/* --- Hàm tính toán: bê nguyên từ bản cũ ---------------------------------- */

/** Rút gọn số lớn: 1.234.567 -> "1.2 Tr". */
function formatK($num)
{
    if ($num >= 1000000000) return round($num / 1000000000, 1) . " Tỷ";
    if ($num >= 1000000)    return round($num / 1000000, 1) . " Tr";
    if ($num >= 1000)       return round($num / 1000, 1) . " K";
    return number_format($num);
}

/** Gộp hành trang và rương để lấy tổng số lượng của một loại vật phẩm. */
function getItemQuantityCombined($items_bag, $items_box, $id)
{
    $total = 0;
    foreach ([$items_bag, $items_box] as $itemsText) {
        $items = json_decode($itemsText, true);
        if (!$items) continue;
        foreach ($items as $it) {
            $it = json_decode($it, true);
            if (isset($it[0]) && $it[0] == $id) {
                $total += $it[1] ?? 0;
            }
        }
    }
    return $total;
}

/** Vàng chìm / ngọc xanh / hồng ngọc nằm ở [0], [1], [2] của data_inventory. */
function getInventory($data_inventory)
{
    $inv = json_decode($data_inventory, true);
    return [
        "vang_chim" => $inv[0] ?? 0,
        "ngoc_xanh" => $inv[1] ?? 0,
        "hong_ngoc" => $inv[2] ?? 0
    ];
}

/** Sức đánh ở data_point[7]; 3 sao là vật phẩm 16, 7 sao là vật phẩm 20. */
function getStat($items_bag, $items_box, $data_point)
{
    $arr = json_decode($data_point, true);
    return [
        "sd" => $arr[7] ?? 0,
        "hp" => getItemQuantityCombined($items_bag, $items_box, 16),
        "ki" => getItemQuantityCombined($items_bag, $items_box, 20)
    ];
}

/* --- Lấy dữ liệu --------------------------------------------------------- */

$vangTuoi = $vangChim = $ngocXanh = $hongNgoc = [];
$sdBXH = $hpBXH = $kiBXH = [];

try {
    /* TODO hiệu năng: mọi giá trị ở trang này đều nằm trong cột JSON
       (data_inventory, items_bag, items_box, data_point) nên bắt buộc phải đọc
       hết bảng rồi sắp xếp bằng PHP — không chuyển sang ORDER BY ... LIMIT
       trong SQL được mà không làm sai kết quả. Khi số nhân vật lên tới vài
       nghìn, cần server Java ghi sẵn các cột số riêng rồi mới tối ưu được. */
    $res = $pdo->query("SELECT name, gender, head, data_inventory, items_bag, items_box, data_point FROM player");

    while ($row = $res->fetch()) {
        $inv  = getInventory($row['data_inventory']);
        $stat = getStat($row['items_bag'], $row['items_box'], $row['data_point']);
        $vn   = getItemQuantityCombined($row['items_bag'], $row['items_box'], 457);

        $base = [
            "name"   => $row['name'],
            "gender" => (int)$row['gender'],
            "head"   => (int)$row['head']
        ];

        $vangTuoi[] = $base + ["value" => $vn];
        $vangChim[] = $base + ["value" => $inv['vang_chim']];
        $ngocXanh[] = $base + ["value" => $inv['ngoc_xanh']];
        $hongNgoc[] = $base + ["value" => $inv['hong_ngoc']];

        $sdBXH[] = $base + ["value" => $stat['sd']];
        $hpBXH[] = $base + ["value" => $stat['hp']];
        $kiBXH[] = $base + ["value" => $stat['ki']];
    }
} catch (Throwable $e) {
    // Lỗi CSDL thì hiển thị bảng rỗng chứ không để trang trắng.
    $vangTuoi = $vangChim = $ngocXanh = $hongNgoc = [];
    $sdBXH = $hpBXH = $kiBXH = [];
}

function top10(&$arr)
{
    usort($arr, fn($a, $b) => $b['value'] <=> $a['value']);
    $arr = array_slice($arr, 0, 10);
}

top10($vangTuoi);
top10($vangChim);
top10($ngocXanh);
top10($hongNgoc);
top10($sdBXH);
top10($hpBXH);
top10($kiBXH);

/* Tra ảnh đại diện: player.head -> head_avatar.avatar_id (bổ sung mới). */
$bxhHeads = array_merge(
    array_column($vangTuoi, 'head'),
    array_column($vangChim, 'head'),
    array_column($ngocXanh, 'head'),
    array_column($hongNgoc, 'head'),
    array_column($sdBXH, 'head'),
    array_column($hpBXH, 'head'),
    array_column($kiBXH, 'head')
);
$bxhAvatars = kame_head_avatars($bxhHeads);

/** Lớp CSS cho chip hành tinh. */
function bxhPlanetClass(int $gender): string
{
    return match ($gender) {
        1 => 'kame-badge--success',   // Namếc
        2 => 'kame-badge--info',      // Xayda
        default => '',                // Trái Đất
    };
}

/** In một bảng xếp hạng theo mẫu Hạng · Ảnh · Tên · Hành tinh · Chỉ số. */
function renderTable(string $title, array $data, string $icon, string $unit): void
{
    global $bxhAvatars;
    ?>
    <h2 class="kame-section-title"><?= kame_e($title) ?></h2>
    <div class="kame-table-wrap">
      <table class="kame-table bxh-table">
        <thead>
          <tr>
            <th scope="col">Hạng</th>
            <th scope="col">Ảnh</th>
            <th scope="col">Tên</th>
            <th scope="col" class="bxh-planet-head">Hành tinh</th>
            <th scope="col" class="bxh-value-head"><?= kame_e($unit) ?></th>
          </tr>
        </thead>
        <tbody>
        <?php if (!$data): ?>
          <tr><td colspan="5" class="bxh-empty">
            <strong>Chưa có dữ liệu</strong>
            hãy là người đầu tiên
          </td></tr>
        <?php else: foreach ($data as $i => $p): ?>
          <?php $avatar = kame_icon_url($bxhAvatars[(int)$p['head']] ?? 1); ?>
          <tr<?= $i < 3 ? ' class="bxh-row--top"' : '' ?>>
            <td class="bxh-pos"><?= $i + 1 ?></td>
            <td><span class="bxh-ava"><img src="<?= $avatar ?>" alt="" loading="lazy" width="46" height="46"></span></td>
            <td>
              <div class="bxh-name"><?= kame_e($p['name']) ?></div>
              <div class="bxh-sub bxh-sub--planet"><?= kame_e(kame_planet_name($p['gender'])) ?></div>
            </td>
            <td class="bxh-planet-cell">
              <span class="kame-badge <?= bxhPlanetClass($p['gender']) ?>"><?= kame_e(kame_planet_name($p['gender'])) ?></span>
            </td>
            <td class="bxh-value-cell">
              <span class="bxh-value">
                <img class="bxh-unit-icon" src="<?= kame_e($icon) ?>" alt="" width="22" height="22">
                <?= kame_e(formatK($p['value'])) ?>
              </span>
            </td>
          </tr>
        <?php endforeach; endif; ?>
        </tbody>
      </table>
    </div>
    <?php
}

kame_head(
    'Bảng xếp hạng tài sản',
    'Bảng xếp hạng vàng, ngọc và chỉ số của máy chủ Ngọc Rồng Kame — số liệu lấy trực tiếp từ game.'
);
kame_header('bxh');
?>

<style>
/* CSS riêng cho hai trang bảng xếp hạng — chỉ dùng biến màu của bộ UI kit */
[hidden] { display: none !important; }

.bxh-tablink {
  display: inline-flex;
  align-items: center;
  text-decoration: none;
  white-space: nowrap;
}

.bxh-table th:first-child, .bxh-table td:first-child { text-align: center; }
.bxh-table td { vertical-align: middle; }
.bxh-table tbody tr:last-child td { border-bottom: 0; }

.bxh-pos {
  width: 72px;
  color: var(--kame-maroon-700);
  font-family: "Roboto Condensed", "Arial Narrow", sans-serif;
  font-size: 20px;
  font-style: italic;
  font-weight: 900;
}

.bxh-ava {
  display: grid;
  width: 46px;
  height: 46px;
  overflow: hidden;
  place-items: center;
  border: 2px solid var(--kame-maroon-700);
  border-radius: 50%;
  background: var(--kame-paper-300);
}
.bxh-ava img { width: 82%; height: 82%; object-fit: contain; image-rendering: pixelated; }

.bxh-name { font-weight: 800; overflow-wrap: anywhere; }
.bxh-sub { font-size: 12px; color: rgb(40 21 15 / 64%); }
.bxh-sub--planet { display: none; }

.bxh-value {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--kame-maroon-700);
  font-family: "Roboto Condensed", "Arial Narrow", sans-serif;
  font-size: 20px;
  font-style: italic;
  font-weight: 900;
  white-space: nowrap;
}
.bxh-unit-icon { width: 22px; height: 22px; object-fit: contain; image-rendering: pixelated; }
.bxh-table td.bxh-value-cell, .bxh-table th.bxh-value-head { text-align: right; }

/* Ba hạng đầu: số hạng to hơn, màu vàng, nền nhạt hơn */
.bxh-row--top td { background: linear-gradient(90deg, rgb(255 200 61 / 34%), rgb(255 243 213 / 55%)); }
.bxh-row--top .bxh-pos {
  color: var(--kame-gold);
  font-size: 34px;
  text-shadow: 0 2px 0 var(--kame-maroon-700), 0 -1px 0 var(--kame-maroon-700),
               1px 0 0 var(--kame-maroon-700), -1px 0 0 var(--kame-maroon-700);
}
.bxh-row--top .bxh-ava { width: 56px; height: 56px; border-width: 3px; border-color: var(--kame-gold); }
.bxh-row--top .bxh-name { font-size: 18px; }
.bxh-row--top .bxh-value { font-size: 25px; }

.bxh-empty {
  padding: var(--kame-space-7) var(--kame-space-4);
  color: rgb(40 21 15 / 62%);
  text-align: center;
}
.bxh-empty strong { display: block; color: var(--kame-maroon-700); font-size: 18px; }

.bxh-note { font-size: 14px; }
.bxh-board + .bxh-board { margin-top: var(--kame-space-6); }

@media (max-width: 1180px) {
  .bxh-table th, .bxh-table td { padding: 12px 12px; }
}

@media (max-width: 900px) {
  .bxh-pos { width: 56px; font-size: 18px; }
  .bxh-row--top .bxh-pos { font-size: 28px; }
}

@media (max-width: 600px) {
  .bxh-table th, .bxh-table td { padding: 10px 8px; font-size: 14px; }
  .bxh-pos { width: 40px; }
  .bxh-ava { width: 38px; height: 38px; }
  .bxh-row--top .bxh-ava { width: 44px; height: 44px; }
  .bxh-planet-cell, .bxh-planet-head { display: none; }
  .bxh-sub--planet { display: block; }
  .bxh-value { font-size: 17px; }
  .bxh-row--top .bxh-value { font-size: 20px; }
}
</style>

<main class="kame-page-wrap">
  <section class="kame-section">
    <div class="kame-container kame-stack">

      <nav class="kame-breadcrumb" aria-label="Đường dẫn">
        <a href="/sv1/trangchu">Trang chủ</a>
        <span>Bảng xếp hạng</span>
        <span>Tài sản</span>
      </nav>

      <div>
        <h1 class="kame-heading kame-heading--lg">Bảng xếp hạng</h1>
        <p class="kame-muted bxh-note">Số liệu đọc thẳng từ cơ sở dữ liệu máy chủ. Tải lại trang để xem thứ hạng mới nhất.</p>
      </div>

      <!-- Chuyển qua lại giữa hai trang bảng xếp hạng -->
      <div class="kame-tabs" aria-label="Chọn bảng xếp hạng">
        <a class="kame-tab bxh-tablink" href="/sv1/bxh01" aria-selected="false">Sức mạnh</a>
        <a class="kame-tab bxh-tablink" href="/sv1/bxh02" aria-current="page" aria-selected="true">Tài sản</a>
      </div>

      <!-- Hai nhóm bảng con trong trang này -->
      <div class="kame-tabs" role="tablist" aria-label="Chọn nội dung xếp hạng">
        <button type="button" class="kame-tab" role="tab" aria-selected="true"  data-bxh-tab="tn">Vàng &amp; ngọc</button>
        <button type="button" class="kame-tab" role="tab" aria-selected="false" data-bxh-tab="cs">Chỉ số</button>
      </div>

      <div data-bxh-panel="tn">
        <div class="bxh-board"><?php renderTable("Vàng nổi",  $vangTuoi, "/img/4028.png", "Số lượng"); ?></div>
        <div class="bxh-board"><?php renderTable("Vàng chìm", $vangChim, "/img/930.png",  "Số lượng"); ?></div>
        <div class="bxh-board"><?php renderTable("Ngọc xanh", $ngocXanh, "/img/932.png",  "Số lượng"); ?></div>
        <div class="bxh-board"><?php renderTable("Hồng ngọc", $hongNgoc, "/img/7743.png", "Số lượng"); ?></div>
      </div>

      <div data-bxh-panel="cs" hidden>
        <div class="bxh-board"><?php renderTable("Sức đánh", $sdBXH, "/img/2754.png", "Sức đánh"); ?></div>
        <div class="bxh-board"><?php renderTable("3 sao",    $hpBXH, "/img/114.png",  "Số lượng"); ?></div>
        <div class="bxh-board"><?php renderTable("7 sao",    $kiBXH, "/img/118.png",  "Số lượng"); ?></div>
      </div>

    </div>
  </section>
</main>

<script>
/* Chuyển tab giữa hai nhóm bảng con — không tải lại trang */
(function () {
  "use strict";
  var tabs = document.querySelectorAll("[data-bxh-tab]");
  tabs.forEach(function (tab) {
    tab.addEventListener("click", function () {
      var key = tab.getAttribute("data-bxh-tab");
      tabs.forEach(function (t) {
        t.setAttribute("aria-selected", t === tab ? "true" : "false");
      });
      document.querySelectorAll("[data-bxh-panel]").forEach(function (panel) {
        panel.hidden = panel.getAttribute("data-bxh-panel") !== key;
      });
    });
  });
})();
</script>

<?php
kame_auth_modals();
kame_footer();
