<?php
/* ==========================================================================
   NGỌC RỒNG KAME — TRANG BÁO LỖI 404 / 403

   Apache gọi file này qua ErrorDocument trong .htaccess, nên nó phải tự đặt
   đúng mã trạng thái: nếu trả về 200 thì Google sẽ lập chỉ mục hàng loạt
   đường dẫn rác, còn trình duyệt tưởng trang tồn tại thật.
   ========================================================================== */

declare(strict_types=1);

require_once __DIR__ . '/kame/layout.php';

$maLoi = (int)($_SERVER['REDIRECT_STATUS'] ?? 404);
if ($maLoi !== 403) {
    $maLoi = 404;
}
http_response_code($maLoi);

/* Đường dẫn người dùng vừa gõ, để hiện lại cho họ thấy mình gõ nhầm chỗ nào.
   Cắt ngắn và lọc thẻ vì đây là dữ liệu do người ngoài nhập. */
$duongDan = (string)($_SERVER['REDIRECT_URL'] ?? $_SERVER['REQUEST_URI'] ?? '');
$duongDan = mb_substr(strtok($duongDan, '?') ?: '', 0, 120, 'UTF-8');

$tieuDe  = $maLoi === 403 ? 'Không có quyền xem' : 'Không tìm thấy trang';
$loiNhan = $maLoi === 403
    ? 'Trang này có tồn tại nhưng không mở cho người ngoài. Nếu bạn nghĩ đây là nhầm lẫn, hãy báo cho Ban quản trị.'
    : 'Đường dẫn bạn vừa mở không có trên máy chủ. Có thể bạn gõ nhầm, hoặc trang đã được đổi tên.';

kame_head($tieuDe);
kame_header('');
?>

<style>
.kt404 { padding: 96px 0 120px; text-align: center; }
.kt404__ma {
  font-family: var(--kame-font-title, 'Roboto Condensed', sans-serif);
  font-size: clamp(88px, 18vw, 190px);
  line-height: .9;
  color: var(--kame-cam, #f06030);
  opacity: .28;
  letter-spacing: 4px;
}
.kt404__tieude { margin: -18px 0 14px; }
.kt404__duongdan {
  display: inline-block;
  max-width: 100%;
  margin: 18px 0 26px;
  padding: 8px 14px;
  border-radius: 6px;
  background: rgba(0,0,0,.05);
  font-family: ui-monospace, Consolas, monospace;
  font-size: 14px;
  word-break: break-all;
}
.kt404__nut { display: flex; gap: 12px; justify-content: center; flex-wrap: wrap; }
@media (max-width: 600px) { .kt404 { padding: 64px 0 80px; } }
</style>

<main>
  <section class="kame-section kt404">
    <div class="kame-container">

      <div class="kt404__ma"><?= $maLoi ?></div>
      <h1 class="kame-heading kame-heading--lg kt404__tieude"><?= kame_e($tieuDe) ?></h1>
      <p class="kame-muted" style="max-width:560px;margin:0 auto"><?= kame_e($loiNhan) ?></p>

      <?php if ($duongDan !== ''): ?>
        <div class="kt404__duongdan"><?= kame_e($duongDan) ?></div>
      <?php endif; ?>

      <div class="kt404__nut">
        <a class="kame-button" href="/sv1/trangchu">Về trang chủ</a>
        <a class="kame-button kame-button--secondary" href="/sv1/tin-tuc">Xem tin tức</a>
        <a class="kame-button kame-button--secondary" href="/sv1/ho-tro">Cần hỗ trợ</a>
      </div>

    </div>
  </section>
</main>

<?php
kame_auth_modals();
kame_footer();
