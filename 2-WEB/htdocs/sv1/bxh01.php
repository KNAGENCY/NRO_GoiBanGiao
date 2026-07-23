<?php
/* ==========================================================================
   BẢNG XẾP HẠNG — SỨC MẠNH
   Khoác giao diện Ngọc Rồng Kame. Toàn bộ cách tính điểm giữ nguyên như bản cũ
   (xem bản sao lưu bxh01.php.cu). Chỉ đổi phần HTML/CSS xuất ra và bổ sung
   ảnh đại diện lấy từ game qua bảng head_avatar.
   ========================================================================== */

require_once __DIR__ . '/kame/layout.php';   // layout.php tự nạp data.php và config.php

/* --- Hàm tính toán: bê nguyên từ bản cũ ---------------------------------- */

/** Sức mạnh bản thân nằm ở phần tử [1] của cột data_point. */
function getPower($data_point)
{
    $arr = json_decode($data_point, true);
    return $arr[1] ?? 0;
}

/** Sức mạnh đệ tử: cột `pet` là JSON lồng, phần tử [1] lại là một chuỗi JSON. */
function getPetPower($pet)
{
    $pet = json_decode($pet, true);
    if (isset($pet[1])) {
        $info = json_decode($pet[1], true);
        return $info[1] ?? 0;
    }
    return 0;
}

/** Nhiệm vụ chính / nhiệm vụ phụ nằm ở [0] và [1] của cột data_task. */
function getTask($data_task)
{
    $task = json_decode($data_task, true);
    return [
        "main" => $task[0] ?? 0,
        "sub"  => $task[1] ?? 0
    ];
}

/** Rút gọn số lớn: 1.234.567 -> "1.2 Tr". */
function formatK($num)
{
    if ($num >= 1000000000) return round($num / 1000000000, 1) . " Tỷ";
    if ($num >= 1000000)    return round($num / 1000000, 1) . " Tr";
    if ($num >= 1000)       return round($num / 1000, 1) . " K";
    return number_format($num);
}

/* --- Lấy dữ liệu --------------------------------------------------------- */

$topSM   = [];
$topTask = [];
$topNap  = [];

try {
    /* Bảng sức mạnh và bảng nhiệm vụ đọc chung một lượt quét bảng `player`
       thay vì hai lượt như bản cũ — kết quả y hệt, đỡ được một lần đọc bảng.

       TODO hiệu năng: sức mạnh nằm trong cột JSON `data_point` và `pet` nên
       vẫn phải sắp xếp bằng PHP sau khi đọc hết bảng. Khi số nhân vật lên tới
       vài nghìn, cần thêm cột số riêng (ví dụ `power_cached`) do server Java
       ghi sẵn, rồi mới chuyển được sang ORDER BY ... LIMIT trong SQL. */
    $res = $pdo->query("SELECT name, gender, head, data_point, pet, data_task FROM player");
    while ($row = $res->fetch()) {
        $sm  = getPower($row['data_point']);
        $pet = getPetPower($row['pet']);
        $topSM[] = [
            "name"   => $row['name'],
            "gender" => (int)$row['gender'],
            "head"   => (int)$row['head'],
            "total"  => $sm + $pet,
            "sm"     => $sm,
            "pet"    => $pet
        ];

        $task = getTask($row['data_task']);
        $topTask[] = [
            "name"   => $row['name'],
            "gender" => (int)$row['gender'],
            "head"   => (int)$row['head'],
            "main"   => $task['main'],
            "sub"    => $task['sub']
        ];
    }

    usort($topSM, fn($a, $b) => $b['total'] <=> $a['total']);
    $topSM = array_slice($topSM, 0, 20);

    usort($topTask, function ($a, $b) {
        if ($b['main'] === $a['main']) {
            return $b['sub'] <=> $a['sub'];
        }
        return $b['main'] <=> $a['main'];
    });
    $topTask = array_slice($topTask, 0, 20);

    /* Bảng nạp xếp theo cột số `danap` nên để MySQL sắp xếp và cắt luôn 20 dòng,
       không cần kéo cả bảng về PHP như bản cũ.
       Xếp thêm theo p.id để hai người bằng điểm luôn ra cùng một thứ tự — giống
       hệt kết quả của bản cũ và không đổi vị trí mỗi lần tải trang. */
    $res = $pdo->query(
        "SELECT p.name, p.gender, p.head, a.danap
         FROM player p JOIN account a ON p.account_id = a.id
         ORDER BY a.danap DESC, p.id ASC LIMIT 20"
    );
    while ($row = $res->fetch()) {
        $topNap[] = [
            "name"   => $row['name'],
            "gender" => (int)$row['gender'],
            "head"   => (int)$row['head'],
            "danap"  => $row['danap']
        ];
    }
} catch (Throwable $e) {
    // Lỗi CSDL thì hiển thị bảng rỗng chứ không để trang trắng.
    $topSM = $topTask = $topNap = [];
}

/* Tra ảnh đại diện: player.head -> head_avatar.avatar_id (bổ sung mới). */
$bxhHeads = array_merge(
    array_column($topSM, 'head'),
    array_column($topTask, 'head'),
    array_column($topNap, 'head')
);
$bxhAvatars = kame_head_avatars($bxhHeads);

/** Ảnh đại diện của một dòng, mặc định avatar 1 nếu chưa có trong head_avatar. */
function bxhAvatarUrl(array $row): string
{
    global $bxhAvatars;
    return kame_icon_url($bxhAvatars[(int)$row['head']] ?? 1);
}

/** Lớp CSS cho chip hành tinh. */
function bxhPlanetClass(int $gender): string
{
    return match ($gender) {
        1 => 'kame-badge--success',   // Namếc
        2 => 'kame-badge--info',      // Xayda
        default => '',                // Trái Đất
    };
}

kame_head(
    'Bảng xếp hạng sức mạnh',
    'Bảng xếp hạng sức mạnh, nhiệm vụ và cá cá của máy chủ Ngọc Rồng Kame — số liệu lấy trực tiếp từ game.'
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
  color: var(--kame-maroon-700);
  font-family: "Roboto Condensed", "Arial Narrow", sans-serif;
  font-size: 20px;
  font-style: italic;
  font-weight: 900;
  white-space: nowrap;
  text-align: right;
}
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
        <span>Sức mạnh</span>
      </nav>

      <div>
        <h1 class="kame-heading kame-heading--lg">Bảng xếp hạng</h1>
        <p class="kame-muted bxh-note">Số liệu đọc thẳng từ cơ sở dữ liệu máy chủ. Tải lại trang để xem thứ hạng mới nhất.</p>
      </div>

      <!-- Chuyển qua lại giữa hai trang bảng xếp hạng -->
      <div class="kame-tabs" aria-label="Chọn bảng xếp hạng">
        <a class="kame-tab bxh-tablink" href="/sv1/bxh01" aria-current="page" aria-selected="true">Sức mạnh</a>
        <a class="kame-tab bxh-tablink" href="/sv1/bxh02" aria-selected="false">Tài sản</a>
      </div>

      <!-- Ba bảng con trong trang này -->
      <div class="kame-tabs" role="tablist" aria-label="Chọn nội dung xếp hạng">
        <button type="button" class="kame-tab" role="tab" aria-selected="true"  data-bxh-tab="sm">Sức mạnh</button>
        <button type="button" class="kame-tab" role="tab" aria-selected="false" data-bxh-tab="task">Nhiệm vụ</button>
        <button type="button" class="kame-tab" role="tab" aria-selected="false" data-bxh-tab="nap">Cá cá</button>
      </div>

      <!-- ================= TOP SỨC MẠNH ================= -->
      <div data-bxh-panel="sm">
        <h2 class="kame-section-title">Top 20 sức mạnh</h2>
        <div class="kame-table-wrap">
          <table class="kame-table bxh-table">
            <thead>
              <tr>
                <th scope="col">Hạng</th>
                <th scope="col">Ảnh</th>
                <th scope="col">Tên</th>
                <th scope="col" class="bxh-planet-head">Hành tinh</th>
                <th scope="col" class="bxh-value-head">Sức mạnh</th>
              </tr>
            </thead>
            <tbody>
            <?php if (!$topSM): ?>
              <tr><td colspan="5" class="bxh-empty">
                <strong>Chưa có dữ liệu</strong>
                hãy là người đầu tiên
              </td></tr>
            <?php else: foreach ($topSM as $i => $p): ?>
              <tr<?= $i < 3 ? ' class="bxh-row--top"' : '' ?>>
                <td class="bxh-pos"><?= $i + 1 ?></td>
                <td><span class="bxh-ava"><img src="<?= bxhAvatarUrl($p) ?>" alt="" loading="lazy" width="46" height="46"></span></td>
                <td>
                  <div class="bxh-name"><?= kame_e($p['name']) ?></div>
                  <div class="bxh-sub bxh-sub--planet"><?= kame_e(kame_planet_name($p['gender'])) ?></div>
                </td>
                <td class="bxh-planet-cell">
                  <span class="kame-badge <?= bxhPlanetClass($p['gender']) ?>"><?= kame_e(kame_planet_name($p['gender'])) ?></span>
                </td>
                <td class="bxh-value-cell">
                  <div class="bxh-value"><?= kame_e(formatK($p['total'])) ?></div>
                  <div class="bxh-sub">Sư <?= kame_e(formatK($p['sm'])) ?> · Đệ <?= kame_e(formatK($p['pet'])) ?></div>
                </td>
              </tr>
            <?php endforeach; endif; ?>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ================= TOP NHIỆM VỤ ================= -->
      <div data-bxh-panel="task" hidden>
        <h2 class="kame-section-title">Top 20 nhiệm vụ</h2>
        <div class="kame-table-wrap">
          <table class="kame-table bxh-table">
            <thead>
              <tr>
                <th scope="col">Hạng</th>
                <th scope="col">Ảnh</th>
                <th scope="col">Tên</th>
                <th scope="col" class="bxh-planet-head">Hành tinh</th>
                <th scope="col" class="bxh-value-head">Nhiệm vụ</th>
              </tr>
            </thead>
            <tbody>
            <?php if (!$topTask): ?>
              <tr><td colspan="5" class="bxh-empty">
                <strong>Chưa có dữ liệu</strong>
                hãy là người đầu tiên
              </td></tr>
            <?php else: foreach ($topTask as $i => $p): ?>
              <tr<?= $i < 3 ? ' class="bxh-row--top"' : '' ?>>
                <td class="bxh-pos"><?= $i + 1 ?></td>
                <td><span class="bxh-ava"><img src="<?= bxhAvatarUrl($p) ?>" alt="" loading="lazy" width="46" height="46"></span></td>
                <td>
                  <div class="bxh-name"><?= kame_e($p['name']) ?></div>
                  <div class="bxh-sub bxh-sub--planet"><?= kame_e(kame_planet_name($p['gender'])) ?></div>
                </td>
                <td class="bxh-planet-cell">
                  <span class="kame-badge <?= bxhPlanetClass($p['gender']) ?>"><?= kame_e(kame_planet_name($p['gender'])) ?></span>
                </td>
                <td class="bxh-value-cell">
                  <div class="bxh-value"><?= kame_e((string)$p['main']) ?>-<?= kame_e((string)$p['sub']) ?></div>
                  <div class="bxh-sub">Chính · phụ</div>
                </td>
              </tr>
            <?php endforeach; endif; ?>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ================= TOP CÁ CÁ ================= -->
      <div data-bxh-panel="nap" hidden>
        <h2 class="kame-section-title">Top 20 cá cá</h2>
        <div class="kame-table-wrap">
          <table class="kame-table bxh-table">
            <thead>
              <tr>
                <th scope="col">Hạng</th>
                <th scope="col">Ảnh</th>
                <th scope="col">Tên</th>
                <th scope="col" class="bxh-planet-head">Hành tinh</th>
                <th scope="col" class="bxh-value-head">Điểm</th>
              </tr>
            </thead>
            <tbody>
            <?php if (!$topNap): ?>
              <tr><td colspan="5" class="bxh-empty">
                <strong>Chưa có dữ liệu</strong>
                hãy là người đầu tiên
              </td></tr>
            <?php else: foreach ($topNap as $i => $p): ?>
              <tr<?= $i < 3 ? ' class="bxh-row--top"' : '' ?>>
                <td class="bxh-pos"><?= $i + 1 ?></td>
                <td><span class="bxh-ava"><img src="<?= bxhAvatarUrl($p) ?>" alt="" loading="lazy" width="46" height="46"></span></td>
                <td>
                  <div class="bxh-name"><?= kame_e($p['name']) ?></div>
                  <div class="bxh-sub bxh-sub--planet"><?= kame_e(kame_planet_name($p['gender'])) ?></div>
                </td>
                <td class="bxh-planet-cell">
                  <span class="kame-badge <?= bxhPlanetClass($p['gender']) ?>"><?= kame_e(kame_planet_name($p['gender'])) ?></span>
                </td>
                <td class="bxh-value-cell">
                  <div class="bxh-value"><?= kame_e(number_format($p['danap'] / 10)) ?></div>
                </td>
              </tr>
            <?php endforeach; endif; ?>
            </tbody>
          </table>
        </div>
      </div>

    </div>
  </section>
</main>

<script>
/* Chuyển tab giữa ba bảng con — không tải lại trang */
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
