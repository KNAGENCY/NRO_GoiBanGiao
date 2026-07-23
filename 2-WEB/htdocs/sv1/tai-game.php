<?php
/* ==========================================================================
   NGỌC RỒNG KAME — TRANG TẢI GAME
   Ba nền tảng Windows / Android / iOS.

   Nguyên tắc quan trọng của trang này: KHÔNG bịa bản cài.
   sv1/download.php chỉ nhận đúng 3 tên file (NRO.apk, NROKAME.zip, IPA.ipa) và
   trả 404 nếu file chưa nằm trong sv1/down/. Vì vậy mỗi thẻ đều kiểm tra
   file có thật hay không rồi mới quyết định hiện nút tải hay chữ "Sắp có" —
   tuyệt đối không để nút nào dẫn tới trang 404.
   ========================================================================== */

declare(strict_types=1);

require_once __DIR__ . '/kame/layout.php';

/* --- Tiện ích riêng của trang ------------------------------------------- */

/** Thư mục chứa bản cài. Đúng thư mục mà download.php đọc. */
const TG_DOWN_DIR = __DIR__ . '/down';

/** Thư mục ảnh của bộ giao diện Kame. */
const TG_ASSET_DIR = __DIR__ . '/kame/assets';

/**
 * Đổi số byte sang chữ dễ đọc, dấu phẩy thập phân kiểu Việt Nam.
 * Chỉ gọi khi file CÓ THẬT — con số luôn là dung lượng thật.
 */
function tg_size_text(int $bytes): string
{
    if ($bytes >= 1073741824) {
        return number_format($bytes / 1073741824, 2, ',', '.') . ' GB';
    }
    if ($bytes >= 1048576) {
        return number_format($bytes / 1048576, 0, ',', '.') . ' MB';
    }
    if ($bytes >= 1024) {
        return number_format($bytes / 1024, 0, ',', '.') . ' KB';
    }
    return $bytes . ' B';
}

/**
 * Tra cứu một bản cài trong sv1/down/.
 *
 * @return array{exists:bool, size_text:string, updated:string}
 */
function tg_download_info(string $fileName): array
{
    $out = ['exists' => false, 'size_text' => '', 'updated' => ''];

    try {
        // basename() để chắc chắn không ai ghép được đường dẫn lạ vào đây.
        $path = TG_DOWN_DIR . '/' . basename($fileName);
        if (!is_file($path) || !is_readable($path)) {
            return $out;
        }

        $size = @filesize($path);
        if ($size === false || $size <= 0) {
            return $out;                      // File rỗng thì coi như chưa có
        }

        $out['exists']    = true;
        $out['size_text'] = tg_size_text((int)$size);

        $mtime = @filemtime($path);
        if ($mtime !== false) {
            $out['updated'] = date('d/m/Y', (int)$mtime);
        }
    } catch (Throwable $e) {
        return ['exists' => false, 'size_text' => '', 'updated' => ''];
    }

    return $out;
}

/** Ảnh dải đầu trang đã có chưa? Chưa có thì dùng ô giữ chỗ. */
function tg_band_ready(): bool
{
    $path = TG_ASSET_DIR . '/band-taigame.jpg';
    return is_file($path) && is_readable($path) && (int)@filesize($path) > 0;
}

/* --- Ba nền tảng --------------------------------------------------------- */

$tgPlatforms = [
    [
        'key'    => 'windows',
        'name'   => 'Windows',
        'badge'  => 'Máy tính',
        'file'   => 'NROKAME.zip',
        'note'   => 'Bản nén cho máy tính Windows. Giải nén rồi chạy file Ngoc Rong Kame.exe, không cần cài đặt.',
        // Biểu tượng màn hình máy tính — hình khối chung, không dùng logo hãng.
        'icon'   => '<rect x="5" y="9" width="38" height="25" rx="3"/><path d="M18 39h12M24 34v5"/>',
        'accent' => 'var(--kame-blue)',
    ],
    [
        'key'    => 'android',
        'name'   => 'Android',
        'badge'  => 'Điện thoại',
        'file'   => 'NRO.apk',
        'note'   => 'File APK cài thẳng vào máy Android. Khi cài cần cho phép cài ứng dụng ngoài cửa hàng.',
        // Biểu tượng điện thoại
        'icon'   => '<rect x="14" y="6" width="20" height="36" rx="4"/><path d="M21 11h6M24 37h.01"/>',
        'accent' => 'var(--kame-green)',
    ],
    [
        'key'    => 'ios',
        'name'   => 'iOS',
        'badge'  => 'iPhone / iPad',
        'file'   => 'IPA.ipa',
        'note'   => 'File IPA cho iPhone và iPad. Cách cài phụ thuộc thiết bị, sẽ có hướng dẫn kèm theo khi phát hành.',
        // Biểu tượng máy tính bảng
        'icon'   => '<rect x="11" y="5" width="26" height="38" rx="4"/><path d="M24 38h.01"/>',
        'accent' => 'var(--kame-maroon-700)',
    ],
];

// Đếm xem đã có bản cài nào chưa, để chọn câu dẫn ở đầu trang cho đúng sự thật.
$tgInfo  = [];
$tgReady = 0;
foreach ($tgPlatforms as $platform) {
    $info = tg_download_info($platform['file']);
    $tgInfo[$platform['key']] = $info;
    if ($info['exists']) {
        $tgReady++;
    }
}

kame_head(
    'Tải game',
    'Tải Ngọc Rồng Kame cho Windows, Android và iOS. Hướng dẫn cài đặt ba bước và giải đáp lỗi thường gặp.'
);
// 'taigame' khong co trong thanh menu chinh nen khong muc nao duoc to sang.
// Truoc day truyen 'trangchu' -> menu to sang nham muc Trang chu.
kame_header('taigame');
?>

<style>
/* --- CSS riêng của trang tải game. Chỉ dùng biến màu của ui-kit. -------- */

/* Dải đầu trang */
.tg-band {
  position: relative;
  display: flex;
  align-items: flex-end;
  overflow: hidden;
  min-height: clamp(220px, 28vw, 340px);
  border-bottom: var(--kame-border);
  background: var(--kame-maroon-900);
}
.tg-band__bg {
  position: absolute;
  inset: 0;
  background: center / cover no-repeat;
}
/* Phủ tối để chữ trắng luôn đọc được trên mọi ảnh */
.tg-band__bg::after {
  position: absolute;
  inset: 0;
  content: "";
  background: linear-gradient(to top, rgb(54 0 14 / 86%) 0%, rgb(54 0 14 / 44%) 55%, rgb(54 0 14 / 18%) 100%);
}
.tg-band__ph { position: absolute; inset: 0; }
/* Ô giữ chỗ phải phủ kín dải, nên bỏ tỉ lệ khung mà kame_placeholder đặt sẵn */
.tg-band__ph .kame-ph { height: 100%; aspect-ratio: auto !important; border: 0; border-radius: 0; }
.tg-band__inner {
  position: relative;
  z-index: 1;
  width: min(var(--kame-container), calc(100% - 32px));
  margin-inline: auto;
  padding-block: var(--kame-space-6) var(--kame-space-5);
}
.tg-band__title { color: var(--kame-white); text-shadow: 0 4px 0 rgb(54 0 14 / 55%); }
.tg-band .kame-breadcrumb { margin-top: var(--kame-space-3); color: rgb(255 253 245 / 82%); }
.tg-band .kame-breadcrumb a { color: var(--kame-gold); }

/* Thẻ nền tảng */
.tg-plat { display: flex; height: 100%; flex-direction: column; }
.tg-plat .kame-card__body { display: flex; flex: 1; flex-direction: column; gap: var(--kame-space-4); }
.tg-plat__head { display: flex; align-items: center; gap: var(--kame-space-4); }
.tg-plat__icon {
  display: grid;
  width: 62px;
  height: 62px;
  flex: none;
  place-items: center;
  border: var(--kame-border);
  border-radius: var(--kame-radius-md);
  background: var(--kame-paper-300);
}
.tg-plat__icon svg { width: 34px; height: 34px; }
.tg-plat__name { margin: 0 0 4px; font-size: 24px; }
.tg-plat__note { margin: 0; font-size: 14px; color: rgb(40 21 15 / 70%); }

/* Bảng thông số nhỏ trong thẻ */
.tg-spec { display: grid; gap: 6px; margin: 0; font-size: 14px; }
.tg-spec div {
  display: flex;
  justify-content: space-between;
  gap: var(--kame-space-3);
  padding-block: 6px;
  border-bottom: 1px dashed rgb(96 0 24 / 26%);
}
.tg-spec dt { color: rgb(40 21 15 / 66%); }
.tg-spec dd { margin: 0; font-weight: 800; text-align: right; overflow-wrap: anywhere; }

.tg-plat__foot { margin-top: auto; display: grid; gap: var(--kame-space-3); }
.tg-plat__foot .kame-button { width: 100%; }
.tg-plat__foot .kame-alert { padding: 12px 14px; font-size: 13px; }

/* Ba bước cài đặt */
.tg-step .kame-card__body { display: grid; gap: var(--kame-space-3); }
.tg-step__num {
  display: grid;
  width: 52px;
  height: 52px;
  place-items: center;
  border: var(--kame-border);
  border-radius: var(--kame-radius-sm);
  color: var(--kame-white);
  background: var(--kame-orange-600);
  font-family: "Roboto Condensed", "Arial Narrow", sans-serif;
  font-size: 26px;
  font-style: italic;
  font-weight: 900;
}
.tg-step__title { margin: 0; font-size: 19px; }
.tg-step p { margin: 0; font-size: 14px; color: rgb(40 21 15 / 74%); }

/* Câu hỏi thường gặp */
.tg-faq summary { position: relative; padding-right: 46px; }
.tg-faq summary::-webkit-details-marker { display: none; }
.tg-faq summary::after {
  position: absolute;
  top: 50%;
  right: 18px;
  content: "+";
  color: var(--kame-orange-600);
  font-family: "Roboto Condensed", "Arial Narrow", sans-serif;
  font-size: 24px;
  font-weight: 900;
  line-height: 1;
  transform: translateY(-50%);
}
.tg-faq details[open] summary::after { content: "\2013"; }
.tg-faq details[open] summary { color: var(--kame-orange-700); }
.tg-faq details > div { font-size: 15px; color: rgb(40 21 15 / 80%); }
.tg-faq details > div p { margin: 0 0 10px; }
.tg-faq details > div p:last-child { margin-bottom: 0; }

/* Khối kêu gọi tạo tài khoản */
.tg-cta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: var(--kame-space-5);
  padding: var(--kame-space-6);
  border: var(--kame-border);
  border-radius: var(--kame-radius-md);
  background: rgb(255 243 213 / 94%);
  box-shadow: var(--kame-shadow-sm);
}
.tg-cta p { margin: var(--kame-space-2) 0 0; max-width: 62ch; font-size: 15px; }

@media (max-width: 900px) {
  .tg-cta { padding: var(--kame-space-5); }
}
@media (max-width: 600px) {
  .tg-plat__icon { width: 52px; height: 52px; }
  .tg-plat__icon svg { width: 28px; height: 28px; }
  .tg-cta .kame-button { width: 100%; }
}
</style>

<main class="kame-page-wrap">

  <!-- ===== Dải đầu trang ================================================ -->
  <section class="tg-band">
    <?php if (tg_band_ready()): ?>
      <div class="tg-band__bg" style="background-image:url('/sv1/kame/assets/band-taigame.jpg')"></div>
    <?php else: ?>
      <div class="tg-band__ph"><?= kame_placeholder('band-taigame.jpg', '1920/600') ?></div>
    <?php endif; ?>

    <div class="tg-band__inner">
      <h1 class="kame-heading kame-heading--lg tg-band__title">Tải game</h1>
      <nav class="kame-breadcrumb" aria-label="Đường dẫn">
        <a href="/sv1/trangchu">Trang chủ</a>
        <span>Tải game</span>
      </nav>
    </div>
  </section>

  <!-- ===== Chọn nền tảng ================================================ -->
  <section class="kame-section">
    <div class="kame-container">

      <div class="kame-section-title">
        <h2 class="kame-heading kame-heading--md">Chọn nền tảng phù hợp</h2>
        <span class="kame-muted">Ngọc Rồng Kame chơi được trên máy tính và điện thoại</span>
      </div>

      <?php if ($tgReady === 0): ?>
        <div class="kame-alert kame-alert--info" style="margin-bottom:var(--kame-space-5)">
          <span aria-hidden="true">!</span>
          <span>
            <strong>Bản cài đang được hoàn thiện.</strong><br>
            Cả ba bản cài chưa được đăng tải nên chưa có nút tải nào hoạt động. Bạn có thể
            đăng ký tài khoản trước ngay bây giờ — tài khoản dùng được luôn khi bản cài lên trang.
          </span>
          <span></span>
        </div>
      <?php endif; ?>

      <div class="kame-grid">
        <?php foreach ($tgPlatforms as $platform):
            $info = $tgInfo[$platform['key']];
        ?>
          <div class="kame-col-4">
            <article class="kame-card tg-plat">
              <div class="kame-card__body">

                <div class="tg-plat__head">
                  <span class="tg-plat__icon" aria-hidden="true">
                    <svg viewBox="0 0 48 48" fill="none"
                         stroke="<?= kame_e($platform['accent']) ?>"
                         stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                      <?= $platform['icon'] /* Chuỗi SVG cố định do trang tự khai báo, không phải dữ liệu người dùng */ ?>
                    </svg>
                  </span>
                  <div>
                    <h3 class="kame-heading tg-plat__name"><?= kame_e($platform['name']) ?></h3>
                    <span class="kame-badge"><?= kame_e($platform['badge']) ?></span>
                  </div>
                </div>

                <p class="tg-plat__note"><?= kame_e($platform['note']) ?></p>

                <dl class="tg-spec">
                  <div>
                    <dt>Tên tệp</dt>
                    <dd><?= kame_e($platform['file']) ?></dd>
                  </div>
                  <div>
                    <dt>Dung lượng</dt>
                    <dd><?= $info['exists'] ? kame_e($info['size_text']) : 'Chưa có dữ liệu' ?></dd>
                  </div>
                  <div>
                    <dt>Cập nhật</dt>
                    <dd><?= ($info['exists'] && $info['updated'] !== '') ? kame_e($info['updated']) : 'Chưa có dữ liệu' ?></dd>
                  </div>
                  <div>
                    <dt>Trạng thái</dt>
                    <dd>
                      <?php if ($info['exists']): ?>
                        <span class="kame-badge kame-badge--success">Sẵn sàng</span>
                      <?php else: ?>
                        <span class="kame-badge">Sắp có</span>
                      <?php endif; ?>
                    </dd>
                  </div>
                </dl>

                <div class="tg-plat__foot">
                  <?php if ($info['exists']): ?>
                    <a class="kame-button" href="/sv1/download?file=<?= rawurlencode($platform['file']) ?>">
                      Tải ngay
                    </a>
                  <?php else: ?>
                    <!-- Nút mờ, không có href: không bao giờ dẫn tới trang 404 -->
                    <span class="kame-button" aria-disabled="true" role="link">Sắp có</span>
                    <div class="kame-alert kame-alert--info">
                      <span aria-hidden="true">!</span>
                      <span>Bản cài đang được hoàn thiện</span>
                      <span></span>
                    </div>
                  <?php endif; ?>
                </div>

              </div>
            </article>
          </div>
        <?php endforeach; ?>
      </div>

    </div>
  </section>

  <!-- ===== Cấu hình đề nghị ============================================= -->
  <section class="kame-section kame-section--tight">
    <div class="kame-container">

      <div class="kame-section-title">
        <h2 class="kame-heading kame-heading--md">Cấu hình đề nghị</h2>
        <span class="kame-muted">Áp dụng cho bản Windows</span>
      </div>

      <div class="kame-grid">
        <div class="kame-col-8">
          <div class="kame-table-wrap">
            <table class="kame-table">
              <caption class="kame-help" style="padding:12px 16px;text-align:left">
                Bảng dưới là mức đề nghị để chơi mượt trên máy tính.
              </caption>
              <thead>
                <tr>
                  <th scope="col">Mục</th>
                  <th scope="col">Mức đề nghị</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <th scope="row">Hệ điều hành</th>
                  <td>Windows 10 bản 64-bit trở lên</td>
                </tr>
                <tr>
                  <th scope="row">Bộ nhớ RAM</th>
                  <td>4 GB</td>
                </tr>
                <tr>
                  <th scope="row">Dung lượng trống</th>
                  <td>1 GB</td>
                </tr>
                <tr>
                  <th scope="row">Kết nối mạng</th>
                  <td>Bắt buộc — game chơi trực tuyến</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="kame-col-4">
          <div class="kame-alert kame-alert--info">
            <span aria-hidden="true">!</span>
            <span>
              <strong>Đây là con số ước tính.</strong><br>
              Chúng tôi chưa đo trên nhiều dòng máy khác nhau, nên hãy xem đây là mức tham khảo
              chứ không phải yêu cầu cứng. Máy nhẹ hơn vẫn có thể chơi được, chỉ là chưa được kiểm chứng.
            </span>
            <span></span>
          </div>
          <p class="kame-help" style="margin-top:var(--kame-space-4)">
            Cấu hình cho bản Android và iOS chưa được đo, sẽ bổ sung khi bản cài lên trang.
          </p>
        </div>
      </div>

    </div>
  </section>

  <!-- ===== Ba bước cài đặt ============================================== -->
  <section class="kame-section kame-section--tight">
    <div class="kame-container">

      <div class="kame-section-title">
        <h2 class="kame-heading kame-heading--md">Cài đặt trong ba bước</h2>
        <span class="kame-muted">Làm lần lượt từ trên xuống</span>
      </div>

      <div class="kame-grid">

        <div class="kame-col-4">
          <article class="kame-card tg-step">
            <div class="kame-card__body">
              <span class="tg-step__num" aria-hidden="true">1</span>
              <h3 class="kame-heading tg-step__title">Tải tệp cài đặt</h3>
              <p>
                Chọn đúng nền tảng bạn đang dùng ở phần trên rồi bấm nút tải. Chỉ tải từ trang này —
                tệp lấy ở nơi khác có thể đã bị chỉnh sửa.
              </p>
            </div>
          </article>
        </div>

        <div class="kame-col-4">
          <article class="kame-card tg-step">
            <div class="kame-card__body">
              <span class="tg-step__num" aria-hidden="true">2</span>
              <h3 class="kame-heading tg-step__title">Giải nén hoặc cài đặt</h3>
              <p>
                Bản máy tính là tệp nén, cần bấm chuột phải rồi chọn giải nén ra một thư mục trống.
                Bản điện thoại thì mở tệp vừa tải và làm theo hướng dẫn trên máy.
              </p>
            </div>
          </article>
        </div>

        <div class="kame-col-4">
          <article class="kame-card tg-step">
            <div class="kame-card__body">
              <span class="tg-step__num" aria-hidden="true">3</span>
              <h3 class="kame-heading tg-step__title">Đăng nhập và chơi</h3>
              <p>
                Mở game lên, đăng nhập bằng tài khoản đã đăng ký trên trang này, chọn hành tinh
                rồi bắt đầu. Chưa có tài khoản thì đăng ký ngay bên dưới.
              </p>
            </div>
          </article>
        </div>

      </div>
    </div>
  </section>

  <!-- ===== Câu hỏi thường gặp =========================================== -->
  <section class="kame-section kame-section--tight">
    <div class="kame-container">

      <div class="kame-section-title">
        <h2 class="kame-heading kame-heading--md">Câu hỏi thường gặp</h2>
        <span class="kame-muted">Bấm vào từng câu để mở câu trả lời</span>
      </div>

      <div class="kame-grid">
        <div class="kame-col-8">
          <div class="kame-accordion tg-faq">

            <details>
              <summary>Bấm nút tải mà không có gì xảy ra?</summary>
              <div>
                <p>
                  Nếu nút đang là chữ <strong>“Sắp có”</strong> và bị làm mờ thì đúng như vậy: bản cài
                  chưa được đăng tải, nút cố tình không bấm được để bạn khỏi mất công vào trang lỗi.
                </p>
                <p>
                  Khi nút đã sáng thành <strong>“Tải ngay”</strong> mà bấm vẫn không tải được, hãy thử
                  lần lượt: tắt phần mềm chặn quảng cáo, thử một trình duyệt khác, kiểm tra ổ đĩa còn
                  đủ chỗ trống, rồi tải lại trang và bấm lại.
                </p>
              </div>
            </details>

            <details>
              <summary>Máy báo virus hoặc chặn không cho chạy tệp?</summary>
              <div>
                <p>
                  Đây là chuyện thường gặp với các bản cài chưa mua chứng chỉ ký số. Windows sẽ hiện
                  bảng xanh <em>“Windows đã bảo vệ PC của bạn”</em> — bấm dòng
                  <strong>“Thông tin thêm”</strong> rồi chọn <strong>“Vẫn chạy”</strong>. Một số phần
                  mềm diệt virus cũng cảnh báo tương tự và cần bạn cho phép thủ công.
                </p>
                <p>
                  Chỉ làm vậy với tệp tải trực tiếp từ trang này. Nếu bạn tải từ nơi khác thì cảnh báo
                  đó rất đáng để nghe theo.
                </p>
              </div>
            </details>

            <details>
              <summary>Quên mật khẩu thì lấy lại thế nào?</summary>
              <div>
                <p>
                  Bấm <strong>Đăng nhập</strong> ở góc trên rồi chọn <strong>Quên mật khẩu</strong>.
                  Bạn cần nhập đúng tên tài khoản và <strong>số điện thoại đã khai lúc đăng ký</strong>.
                </p>
                <p>
                  Nếu không nhớ số điện thoại đã dùng thì hệ thống không tự khôi phục được —
                  hãy liên hệ qua trang <a class="kame-link" href="/sv1/ho-tro">Hỗ trợ</a>.
                </p>
              </div>
            </details>

            <details>
              <summary>Chơi trên điện thoại được không?</summary>
              <div>
                <p>
                  Được. Máy Android dùng tệp <strong>NRO.apk</strong>, iPhone và iPad dùng tệp
                  <strong>IPA.ipa</strong>.
                </p>
                <p>
                  Riêng Android, lần đầu mở tệp APK máy sẽ hỏi có cho phép cài ứng dụng từ nguồn
                  ngoài cửa hàng hay không — bạn cần bật cho phép thì mới cài tiếp được. Đây là thiết
                  lập sẵn có của Android, không phải lỗi.
                </p>
              </div>
            </details>

            <details>
              <summary>Máy cấu hình yếu có chơi được không?</summary>
              <div>
                <p>
                  Nhiều khả năng là được. Bảng cấu hình phía trên chỉ là mức ước tính để chơi mượt,
                  không phải mức tối thiểu bắt buộc — chúng tôi chưa đo trên nhiều dòng máy nên chưa
                  dám khẳng định con số thấp nhất.
                </p>
                <p>
                  Nếu thấy giật, hãy thử đóng bớt các chương trình đang chạy nền và giảm số cửa sổ
                  trình duyệt đang mở trước khi kết luận là máy không đủ sức.
                </p>
              </div>
            </details>

            <details>
              <summary>Cài xong nhưng không kết nối được vào game?</summary>
              <div>
                <p>
                  Trước hết kiểm tra máy có mạng không — game bắt buộc phải trực tuyến mới chơi được.
                </p>
                <p>
                  Nếu mạng vẫn bình thường thì có thể máy chủ đang bảo trì hoặc tạm dừng. Bạn xem
                  thông báo mới nhất ở trang <a class="kame-link" href="/sv1/tin-tuc">Tin tức</a>,
                  hoặc nhắn cho đội hỗ trợ qua trang
                  <a class="kame-link" href="/sv1/ho-tro">Hỗ trợ</a> để được xem giúp.
                </p>
              </div>
            </details>

          </div>
        </div>

        <div class="kame-col-4">
          <article class="kame-card">
            <div class="kame-card__body">
              <h3 class="kame-heading kame-heading--md" style="font-size:22px">Giữ tài khoản an toàn</h3>
              <ul style="margin:var(--kame-space-4) 0 0;padding-left:20px;font-size:14px;line-height:1.8">
                <li>Không đưa mật khẩu cho bất kỳ ai, kể cả người tự xưng là quản trị viên.</li>
                <li>Không dùng phần mềm hack hoặc tự động — tài khoản có thể bị khoá.</li>
                <li>Chỉ tải bản cài từ trang chính thức này.</li>
                <li>Đặt mật khẩu khác với mật khẩu bạn dùng ở nơi khác.</li>
              </ul>
            </div>
          </article>
        </div>

      </div>
    </div>
  </section>

  <!-- ===== Kêu gọi đăng ký ============================================== -->
  <section class="kame-section kame-section--tight">
    <div class="kame-container">
      <div class="tg-cta">
        <div>
          <h2 class="kame-heading kame-heading--md">Chưa có tài khoản?</h2>
          <p>
            Tạo tài khoản trước rồi tải game sau cũng được. Đăng ký chỉ cần tên tài khoản, mật khẩu
            và một số điện thoại để lấy lại mật khẩu khi quên.
          </p>
        </div>
        <div class="kame-cluster">
          <button type="button" class="kame-button" data-kame-open="dangky">Đăng ký tài khoản</button>
          <a class="kame-button kame-button--secondary" href="/sv1/ho-tro">Cần hỗ trợ</a>
        </div>
      </div>
    </div>
  </section>

</main>

<?php
kame_auth_modals();
kame_footer();
