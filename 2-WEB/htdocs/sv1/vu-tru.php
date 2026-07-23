<?php
/* ==========================================================================
   NGỌC RỒNG KAME — TRANG VŨ TRỤ
   Giới thiệu 3 hành tinh, chiến binh mạnh nhất, danh sách boss và hành trình.

   MỌI SỐ LIỆU TRÊN TRANG NÀY ĐỀU LẤY TỪ DỮ LIỆU THẬT:
     - Địa danh   : bảng map_template
     - Kỹ năng    : bảng skill_template (nclass_id 0/1/2 = Trái Đất/Namếc/Xayda)
     - Mốc sức mạnh: bảng skill_template.power_require và bảng power_limit
     - Danh hiệu  : bảng caption
     - Chiến binh : kame_top_power()  — bảng player
     - Boss       : kame_bosses()     — file boss_status.json do server Java ghi
   Không có con số nào được bịa ra. Thiếu dữ liệu thì hiện "Chưa có dữ liệu".
   ========================================================================== */

declare(strict_types=1);

require_once __DIR__ . '/kame/layout.php';

/* --- Dữ liệu 3 hành tinh -------------------------------------------------
   Tên địa danh chép ĐÚNG NGUYÊN VĂN từ cột map_template.NAME.
   Tên kỹ năng chép ĐÚNG NGUYÊN VĂN từ cột skill_template.NAME, phần mô tả
   rút gọn từ cột dam_info (bỏ ký tự "#" là chỗ game thay số theo cấp).
   ------------------------------------------------------------------------ */
$kamePlanets = [
    'traidat' => [
        'ten'      => 'Trái Đất',
        'phu'      => 'Nơi khởi nguồn',
        'mau'      => 'kame-badge--success',
        'anh'      => 'map-traidat.jpg',
        'mota'     => 'Hành tinh cân bằng nhất: vừa đánh gần vừa đánh xa đều khá. '
                    . 'Người chơi Trái Đất mở đầu ở Làng Aru, học chiêu đấm tại ông nội rồi '
                    . 'lên Tháp Karin và Thần điện để luyện tiếp. Đây là lựa chọn dễ chơi nhất cho người mới.',
        'diadanh'  => [
            'Làng Aru',
            'Đồi hoa cúc',
            'Thung lũng tre',
            'Rừng nấm',
            'Rừng xương',
            'Đảo Kamê',
            'Tháp Karin',
            'Thần điện',
        ],
        'kynang'   => [
            ['Chiêu đấm Dragon',     'Đòn đánh gần cơ bản, tăng sức đánh theo cấp'],
            ['Chiêu Kamejoko',       'Chưởng đánh xa, tăng sức đánh theo cấp'],
            ['Thái Dương Hạ San',    'Làm choáng kẻ địch trong một khoảng thời gian'],
            ['Kaioken',              'Tăng mạnh sức đánh, đổi lại tốn rất nhiều KI'],
            ['Quả cầu kênh khi',     'Đòn cầu năng lượng gây sát thương diện rộng'],
            ['Dịch chuyển tức thời', 'Dịch chuyển tức thời và gây choáng kẻ thù'],
            ['Thôi miên',            'Ru ngủ kẻ thù trong vài giây'],
            ['Khiên năng lượng',     'Vô hiệu các đòn tấn công'],
        ],
    ],
    'namec' => [
        'ten'      => 'Namếc',
        'phu'      => 'Hành tinh của trí tuệ',
        'mau'      => 'kame-badge--info',
        'anh'      => 'map-namec.jpg',
        'mota'     => 'Hành tinh thiên về hồi phục và hỗ trợ. Namếc là phe duy nhất có kỹ năng '
                    . 'trị thương cho đồng đội, nên rất được săn đón khi đi tổ đội đánh boss. '
                    . 'Đổi lại, thể lực thấp hơn hai hành tinh còn lại.',
        'diadanh'  => [
            'Làng Mori',
            'Đồi nấm tím',
            'Thị trấn Moori',
            'Thung lũng Namếc',
            'Thung lũng Maima',
            'Vực maima',
            'Đảo Guru',
            'Nhà Moori',
        ],
        'kynang'   => [
            ['Chiêu đấm Demon',   'Đòn đánh gần cơ bản, tăng sức đánh theo cấp'],
            ['Chiêu Masenko',     'Chưởng đánh xa, tăng sức đánh theo cấp'],
            ['Trị thương',        'Phục hồi HP và KI cho đồng đội'],
            ['Makankosappo',      'Đòn xuyên thấu tầm cực xa, gây sát thương lớn'],
            ['Đẻ trứng',          'Tạo quái đi theo hỗ trợ'],
            ['Liên hoàn',         'Chuỗi đòn đánh nhanh, tăng sức đánh theo cấp'],
            ['Biến Sôcôla',       'Biến quái thành Sôcôla'],
            ['Khiên năng lượng',  'Vô hiệu các đòn tấn công'],
        ],
    ],
    'xayda' => [
        'ten'      => 'Xayda',
        'phu'      => 'Tộc chiến binh',
        'mau'      => 'kame-badge--hot',
        'anh'      => 'map-xayda.jpg',
        'mota'     => 'Hành tinh mạnh nhất về tấn công. Xayda có Biến hình để tăng vọt sức đánh, '
                    . 'HP và tốc độ, cùng Tự phát nổ đổi mạng lấy sát thương. Lối chơi thiên về '
                    . 'áp đảo nhanh, phù hợp người thích đánh nhau trực diện.',
        'diadanh'  => [
            'Làng Kakarot',
            'Đồi hoang',
            'Làng Plant',
            'Rừng nguyên sinh',
            'Rừng thông Xayda',
            'Thành phố Vegeta',
            'Vách núi đen',
            'Nhà Broly',
        ],
        'kynang'   => [
            ['Chiêu đấm Galick',    'Đòn đánh gần cơ bản, sức đánh cao nhất trong ba phe'],
            ['Chiêu Antomic',       'Chưởng đánh xa, tăng sức đánh theo cấp'],
            ['Tái tạo năng lượng',  'Tự hồi HP và KI liên tục theo từng giây'],
            ['Biến hình',           'Tăng sức đánh, HP và tốc độ'],
            ['Tự phát nổ',          'Hy sinh, gây sát thương lớn cho kẻ thù'],
            ['Huýt sáo',            'Tăng tạm thời HP cho mọi người xung quanh, thêm đòn chí mạng'],
            ['Trói',                'Trói kẻ thù tại chỗ'],
            ['Khiên năng lượng',    'Vô hiệu các đòn tấn công'],
        ],
    ],
];

/* --- Hành trình 6 bước ---------------------------------------------------
   Mỗi mốc sức mạnh dưới đây là cột power_require thật trong skill_template,
   riêng bước 6 lấy từ bảng power_limit và bảng caption.
   ------------------------------------------------------------------------ */
$kameHanhTrinh = [
    [
        'so'   => 1,
        'ten'  => 'Khởi đầu',
        'mota' => 'Sinh ra ở Làng Aru, Làng Mori hoặc Làng Kakarot tuỳ hành tinh. '
                . 'Mở rương, nhặt trang bị đầu tiên, học chiêu đấm cơ bản tại ông nội. '
                . 'Mốc mở kỹ năng đầu tiên: 1.000 sức mạnh.',
    ],
    [
        'so'   => 2,
        'ten'  => 'Bái sư',
        'mota' => 'Tìm Sư Phụ của hành tinh mình để học chiêu chưởng đánh xa — '
                . 'Kamejoko, Masenko hoặc Antomic. Cần 10.000 sức mạnh.',
    ],
    [
        'so'   => 3,
        'ten'  => 'Học bay và mở kỹ năng hỗ trợ',
        'mota' => 'Mở được kỹ năng riêng của hành tinh: Thái Dương Hạ San, Trị thương '
                . 'hoặc Tái tạo năng lượng. Cần 60.000 sức mạnh. Từ đây bay được và '
                . 'đi xa khỏi khu vực khởi đầu.',
    ],
    [
        'so'   => 4,
        'ten'  => 'Gia nhập bang hội',
        'mota' => 'Vào một bang hội để nhận nhiệm vụ bang, sinh hoạt tại Lãnh địa Bang Hội '
                . 'và đi tổ đội. Đây là lúc bắt đầu cần bạn chơi cùng.',
    ],
    [
        'so'   => 5,
        'ten'  => 'Đại chiến',
        'mota' => 'Đánh vào Trại lính Fide và Lãnh địa Fize, tham gia Võ đài Xên bọ hung, '
                . 'Đại Chiến Mabư và Đại hội võ thuật. Nhóm kỹ năng mạnh nhất mở ở mốc '
                . '150 triệu đến 1,1 tỷ sức mạnh.',
    ],
    [
        'so'   => 6,
        'ten'  => 'Về cuối game',
        'mota' => 'Nâng kỹ năng lên cấp 7 — mốc cao nhất cần 15,625 tỷ sức mạnh. '
                . 'Bảng giới hạn chỉ số của máy chủ mở dần từ 18 tỷ tới 80 tỷ sức mạnh. '
                . 'Danh hiệu cao nhất trong dữ liệu game là KaioShin cấp 2.',
    ],
];

/* --- Lấy dữ liệu thật ---------------------------------------------------- */
$kameTopPower = kame_top_power(8);
$kameBosses   = kame_bosses();
$kameBossMeta = kame_boss_meta();

/** Đổi số giây thành chuỗi dễ đọc: "3 phút 20 giây". */
function kame_vt_thoi_gian(int $giay): string
{
    if ($giay <= 0) {
        return '—';
    }
    $phut = intdiv($giay, 60);
    $du   = $giay % 60;
    if ($phut > 0) {
        return $du > 0 ? ($phut . ' phút ' . $du . ' giây') : ($phut . ' phút');
    }
    return $du . ' giây';
}

/** Định dạng số theo kiểu Việt Nam: 1.234.567 */
function kame_vt_so(int $n): string
{
    return number_format($n, 0, ',', '.');
}

kame_head(
    'Vũ trụ Kame',
    'Ba hành tinh Trái Đất, Namếc và Xayda trong Ngọc Rồng Kame: địa danh, bộ kỹ năng, '
    . 'chiến binh mạnh nhất và toàn bộ boss đang có trên máy chủ.'
);
kame_header('vutru');
?>

<main class="kame-page-wrap">

  <!-- ====== Tiêu đề trang ====== -->
  <section class="kame-section kame-section--tight">
    <div class="kame-container kame-stack">

      <nav class="kame-breadcrumb" aria-label="Đường dẫn">
        <a href="/sv1/trangchu">Trang chủ</a>
        <span aria-current="page">Vũ trụ</span>
      </nav>

      <div class="kame-center">
        <h1 class="kame-heading kame-heading--xl">Vũ trụ Kame</h1>
        <p class="kame-muted" style="max-width:720px;margin:12px auto 0">
          Ba hành tinh, ba lối chơi khác nhau. Chọn nơi mình sinh ra là quyết định
          đầu tiên và cũng là quyết định khó đổi nhất trong Ngọc Rồng Kame.
        </p>
      </div>

    </div>
  </section>

  <!-- ====== Ba hành tinh ====== -->
  <section class="kame-section" id="hanh-tinh">
    <div class="kame-container kame-stack">

      <div class="kame-section-title">
        <h2 class="kame-heading kame-heading--lg">Ba hành tinh</h2>
        <span class="kame-muted">Địa danh và kỹ năng lấy trực tiếp từ dữ liệu máy chủ</span>
      </div>

      <div class="kame-tabs" role="tablist" aria-label="Chọn hành tinh">
        <?php $dau = true; foreach ($kamePlanets as $ma => $ht): ?>
          <button type="button"
                  class="kame-tab"
                  role="tab"
                  id="tab-<?= kame_e($ma) ?>"
                  aria-controls="panel-<?= kame_e($ma) ?>"
                  aria-selected="<?= $dau ? 'true' : 'false' ?>"
                  data-kame-tab="<?= kame_e($ma) ?>"><?= kame_e($ht['ten']) ?></button>
        <?php $dau = false; endforeach; ?>
      </div>

      <?php $dau = true; foreach ($kamePlanets as $ma => $ht): ?>
        <div class="kame-planet"
             id="panel-<?= kame_e($ma) ?>"
             role="tabpanel"
             aria-labelledby="tab-<?= kame_e($ma) ?>"
             data-kame-panel="<?= kame_e($ma) ?>"
             <?= $dau ? '' : 'hidden' ?>>

          <div class="kame-grid">

            <!-- Bản đồ hành tinh -->
            <div class="kame-col-8">
              <?= kame_placeholder($ht['anh'], '1400/800') ?>
            </div>

            <!-- Giới thiệu -->
            <div class="kame-col-4">
              <div class="kame-card">
                <div class="kame-card__body kame-stack" style="gap:14px">
                  <span class="kame-badge <?= kame_e($ht['mau']) ?>"><?= kame_e($ht['phu']) ?></span>
                  <h3 class="kame-heading kame-heading--md"><?= kame_e($ht['ten']) ?></h3>
                  <p style="margin:0"><?= kame_e($ht['mota']) ?></p>
                </div>
              </div>
            </div>

            <!-- Địa danh tiêu biểu -->
            <div class="kame-col-6">
              <div class="kame-card">
                <div class="kame-card__body">
                  <h3 class="kame-heading kame-heading--md" style="margin-bottom:14px">Địa danh tiêu biểu</h3>
                  <ul class="kame-place-list">
                    <?php foreach ($ht['diadanh'] as $noi): ?>
                      <li><?= kame_e($noi) ?></li>
                    <?php endforeach; ?>
                  </ul>
                  <p class="kame-help" style="margin:14px 0 0">
                    Tên bản đồ ghi đúng như trong game.
                  </p>
                </div>
              </div>
            </div>

            <!-- Bộ kỹ năng -->
            <div class="kame-col-6">
              <div class="kame-card">
                <div class="kame-card__body">
                  <h3 class="kame-heading kame-heading--md" style="margin-bottom:14px">Bộ kỹ năng</h3>
                  <ul class="kame-skill-list">
                    <?php foreach ($ht['kynang'] as $kn): ?>
                      <li>
                        <strong><?= kame_e($kn[0]) ?></strong>
                        <span><?= kame_e($kn[1]) ?></span>
                      </li>
                    <?php endforeach; ?>
                  </ul>
                  <p class="kame-help" style="margin:14px 0 0">
                    Mỗi kỹ năng nâng được tới cấp 7.
                  </p>
                </div>
              </div>
            </div>

          </div>
        </div>
      <?php $dau = false; endforeach; ?>

    </div>
  </section>

  <!-- ====== Chiến binh mạnh nhất (dữ liệu thật từ bảng player) ====== -->
  <section class="kame-section" id="chien-binh">
    <div class="kame-container kame-stack">

      <div class="kame-section-title">
        <h2 class="kame-heading kame-heading--lg">Chiến binh mạnh nhất</h2>
        <a class="kame-link" href="/sv1/bxh01">Xem bảng xếp hạng đầy đủ</a>
      </div>

      <?php if (!$kameTopPower): ?>
        <div class="kame-alert kame-alert--info">
          <span aria-hidden="true">•</span>
          <span>Chưa có dữ liệu người chơi để xếp hạng.</span>
          <span></span>
        </div>
      <?php else: ?>
        <div class="kame-rank">
          <?php foreach ($kameTopPower as $i => $nv): ?>
            <div class="kame-rank-row<?= $i < 3 ? ' kame-rank-row--top' : '' ?>">
              <div class="kame-rank-row__pos"><?= $i + 1 ?></div>
              <div class="kame-rank-row__avatar">
                <img src="<?= kame_e(kame_icon_url((int)$nv['avatar_id'])) ?>"
                     alt="Nhân vật <?= kame_e($nv['name']) ?>" loading="lazy" width="48" height="48">
              </div>
              <div>
                <div class="kame-rank-row__name"><?= kame_e($nv['name']) ?></div>
                <div class="kame-rank-row__meta"><?= kame_e($nv['planet']) ?></div>
              </div>
              <div class="kame-rank-row__value"><?= kame_e(kame_vt_so((int)$nv['power'])) ?></div>
            </div>
          <?php endforeach; ?>
        </div>
        <p class="kame-help" style="margin:0">
          Sức mạnh tính cả sức mạnh đệ tử, giống cột "Tổng" của bảng xếp hạng trong game.
        </p>
      <?php endif; ?>

    </div>
  </section>

  <!-- ====== Boss trong thế giới Kame (đọc boss_status.json) ====== -->
  <section class="kame-section" id="boss">
    <div class="kame-container kame-stack">

      <div class="kame-section-title">
        <h2 class="kame-heading kame-heading--lg">Boss trong thế giới Kame</h2>
        <span class="kame-muted">
          <?= $kameBosses ? kame_e(count($kameBosses) . ' boss trong dữ liệu máy chủ') : 'Chưa có dữ liệu' ?>
        </span>
      </div>

      <?php if ($kameBossMeta['stale'] || !$kameBossMeta['ok']): ?>
        <div class="kame-alert kame-alert--danger">
          <span aria-hidden="true">•</span>
          <span>
            Trạng thái boss chưa được cập nhật (máy chủ game có thể đang tắt).
            Danh sách dưới đây là lần ghi nhận gần nhất, phần Trạng thái và Hồi sinh có thể không còn đúng.
          </span>
          <span></span>
        </div>
      <?php endif; ?>

      <?php if (!$kameBosses): ?>
        <div class="kame-alert kame-alert--info">
          <span aria-hidden="true">•</span>
          <span>Chưa có dữ liệu boss.</span>
          <span></span>
        </div>
      <?php else: ?>

        <div class="kame-field" style="max-width:420px">
          <label class="kame-label" for="boss-tim">Tìm boss theo tên</label>
          <input class="kame-input" id="boss-tim" type="search" autocomplete="off"
                 placeholder="Ví dụ: Xên, Fide, Mabư..." data-kame-boss-search>
          <span class="kame-help" data-kame-boss-count></span>
        </div>

        <div class="kame-table-wrap">
          <table class="kame-table" data-kame-boss-table>
            <thead>
              <tr>
                <th scope="col">Ảnh</th>
                <th scope="col">Tên</th>
                <th scope="col">Trạng thái</th>
                <th scope="col">Máu</th>
                <th scope="col">Bản đồ</th>
                <th scope="col">Hồi sinh</th>
              </tr>
            </thead>
            <tbody>
              <?php foreach ($kameBosses as $boss):
                  $song = ($boss['status'] === 'ACTIVE');
                  $ten  = $boss['name'];
                  $ban  = $boss['mapName'] !== '' ? $boss['mapName'] : 'Chưa có dữ liệu';
              ?>
                <tr data-kame-boss-name="<?= kame_e(mb_strtolower($ten . ' ' . $ban, 'UTF-8')) ?>">
                  <td>
                    <span class="kame-boss-thumb">
                      <img src="<?= kame_e(kame_icon_url((int)$boss['avatar_id'])) ?>"
                           alt="" loading="lazy" width="40" height="40">
                    </span>
                  </td>
                  <th scope="row" style="color:inherit;background:none;text-transform:none;font-weight:800">
                    <?= kame_e($ten) ?>
                  </th>
                  <td>
                    <span class="kame-badge <?= $song ? 'kame-badge--success' : '' ?>">
                      <?= $song ? 'Đang xuất hiện' : 'Chưa xuất hiện' ?>
                    </span>
                  </td>
                  <td style="white-space:nowrap"><?= kame_e(kame_vt_so((int)$boss['maxHp'])) ?></td>
                  <td><?= kame_e($ban) ?></td>
                  <td style="white-space:nowrap">
                    <?php if ($song): ?>
                      Đang có mặt
                    <?php elseif ((int)$boss['respawnIn'] > 0): ?>
                      <?= kame_e(kame_vt_thoi_gian((int)$boss['respawnIn'])) ?>
                    <?php else: ?>
                      Chưa có dữ liệu
                    <?php endif; ?>
                  </td>
                </tr>
              <?php endforeach; ?>
            </tbody>
          </table>
        </div>

        <p class="kame-alert kame-alert--info" data-kame-boss-empty hidden>
          <span aria-hidden="true">•</span>
          <span>Không tìm thấy boss nào khớp với từ khoá.</span>
          <span></span>
        </p>

      <?php endif; ?>

    </div>
  </section>

  <!-- ====== Hành trình 6 bước ====== -->
  <section class="kame-section" id="hanh-trinh">
    <div class="kame-container kame-stack">

      <div class="kame-section-title">
        <h2 class="kame-heading kame-heading--lg">Hành trình của bạn</h2>
        <span class="kame-muted">Các mốc sức mạnh lấy từ bảng kỹ năng của máy chủ</span>
      </div>

      <ol class="kame-journey">
        <?php foreach ($kameHanhTrinh as $buoc): ?>
          <li class="kame-journey__item">
            <span class="kame-journey__num" aria-hidden="true"><?= (int)$buoc['so'] ?></span>
            <div>
              <h3 class="kame-journey__title"><?= kame_e($buoc['ten']) ?></h3>
              <p class="kame-journey__desc"><?= kame_e($buoc['mota']) ?></p>
            </div>
          </li>
        <?php endforeach; ?>
      </ol>

    </div>
  </section>

  <!-- ====== Lời mời cuối trang ====== -->
  <section class="kame-section">
    <div class="kame-container">
      <div class="kame-card">
        <div class="kame-card__body kame-center kame-stack">
          <h2 class="kame-heading kame-heading--lg">Chọn hành tinh của bạn</h2>
          <p class="kame-muted" style="max-width:620px;margin:0 auto">
            Tạo tài khoản, tải game về và bắt đầu ở ngôi làng đầu tiên.
            Mọi thứ trên trang này đang chờ bạn ngoài kia.
          </p>
          <div class="kame-cluster" style="justify-content:center">
            <a class="kame-button" href="/sv1/tai-game">Tải game</a>
            <a class="kame-button kame-button--ghost" href="/sv1/bxh01">Xem bảng xếp hạng</a>
          </div>
        </div>
      </div>
    </div>
  </section>

</main>

<style>
/* --- Riêng trang Vũ trụ -------------------------------------------------- */
.kame-planet[hidden] { display: none; }

.kame-place-list,
.kame-skill-list { margin: 0; padding: 0; list-style: none; display: grid; gap: 8px; }

.kame-place-list li {
  padding: 9px 12px;
  border: 2px solid rgb(96 0 24 / 32%);
  border-radius: var(--kame-radius-sm);
  background: rgb(255 253 245 / 62%);
  font-weight: 600;
}

.kame-skill-list li {
  display: grid;
  gap: 2px;
  padding: 9px 12px;
  border-left: 4px solid var(--kame-orange-600);
  background: rgb(255 253 245 / 62%);
}
.kame-skill-list strong { color: var(--kame-maroon-700); }
.kame-skill-list span { font-size: 13px; color: rgb(40 21 15 / 68%); }

.kame-boss-thumb {
  display: grid;
  width: 40px;
  aspect-ratio: 1;
  place-items: center;
  overflow: hidden;
  border: 2px solid var(--kame-maroon-700);
  border-radius: 50%;
  background: var(--kame-paper-300);
}
.kame-boss-thumb img { width: 80%; height: 80%; object-fit: contain; image-rendering: pixelated; }

/* Hành trình 6 bước */
.kame-journey {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: var(--kame-space-4);
  margin: 0;
  padding: 0;
  list-style: none;
  counter-reset: none;
}

.kame-journey__item {
  display: grid;
  grid-template-columns: 52px 1fr;
  gap: var(--kame-space-4);
  align-items: start;
  padding: var(--kame-space-4);
  border: var(--kame-border);
  border-radius: var(--kame-radius-md);
  background: rgb(255 243 213 / 92%);
  box-shadow: var(--kame-shadow-sm);
}

.kame-journey__num {
  display: grid;
  width: 52px;
  aspect-ratio: 1;
  place-items: center;
  border: 3px solid var(--kame-maroon-700);
  border-radius: 50%;
  color: var(--kame-white);
  background: var(--kame-orange-600);
  font-family: "Roboto Condensed", "Arial Narrow", sans-serif;
  font-size: 24px;
  font-style: italic;
  font-weight: 900;
}

.kame-journey__title {
  margin: 0 0 6px;
  color: var(--kame-maroon-700);
  font-size: 18px;
  font-weight: 800;
  line-height: 1.25;
}
.kame-journey__desc { margin: 0; font-size: 14px; color: rgb(40 21 15 / 76%); }

@media (max-width: 1180px) {
  .kame-journey { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}
@media (max-width: 900px) {
  .kame-journey { grid-template-columns: 1fr; }
}
@media (max-width: 600px) {
  .kame-journey__item { grid-template-columns: 44px 1fr; }
  .kame-journey__num { width: 44px; font-size: 20px; }
}
</style>

<script>
(function () {
  "use strict";

  /* --- Tab ba hành tinh ------------------------------------------------ */
  var tabs = Array.prototype.slice.call(document.querySelectorAll("[data-kame-tab]"));
  var panels = Array.prototype.slice.call(document.querySelectorAll("[data-kame-panel]"));

  function chonTab(ma) {
    tabs.forEach(function (t) {
      t.setAttribute("aria-selected", t.getAttribute("data-kame-tab") === ma ? "true" : "false");
    });
    panels.forEach(function (p) {
      if (p.getAttribute("data-kame-panel") === ma) {
        p.removeAttribute("hidden");
      } else {
        p.setAttribute("hidden", "");
      }
    });
  }

  tabs.forEach(function (t) {
    t.addEventListener("click", function () {
      chonTab(t.getAttribute("data-kame-tab"));
    });
  });

  /* --- Lọc boss theo tên ----------------------------------------------- */
  var o = document.querySelector("[data-kame-boss-search]");
  var bang = document.querySelector("[data-kame-boss-table]");
  var demSlot = document.querySelector("[data-kame-boss-count]");
  var khongCo = document.querySelector("[data-kame-boss-empty]");

  if (o && bang) {
    var hang = Array.prototype.slice.call(bang.querySelectorAll("tbody tr"));
    var tong = hang.length;

    function capNhatDem(hien) {
      if (demSlot) {
        demSlot.textContent = "Đang hiện " + hien + " / " + tong + " boss.";
      }
      if (khongCo) {
        if (hien === 0) {
          khongCo.removeAttribute("hidden");
        } else {
          khongCo.setAttribute("hidden", "");
        }
      }
    }

    o.addEventListener("input", function () {
      var tu = o.value.trim().toLowerCase();
      var hien = 0;
      hang.forEach(function (r) {
        var khop = tu === "" || (r.getAttribute("data-kame-boss-name") || "").indexOf(tu) !== -1;
        r.style.display = khop ? "" : "none";
        if (khop) hien++;
      });
      capNhatDem(hien);
    });

    capNhatDem(tong);
  }
})();
</script>

<?php
kame_auth_modals();
kame_footer();
