<?php
/* ==========================================================================
   NGỌC RỒNG KAME — KHUNG TRANG DÙNG CHUNG
   Mọi trang gọi theo thứ tự:
       kame_head($title, $desc);  kame_header($active);
       ... nội dung ...
       kame_auth_modals();  kame_footer();
   ========================================================================== */

declare(strict_types=1);

require_once __DIR__ . '/data.php';

/* --- Đăng xuất ----------------------------------------------------------
   Xử lý ngay khi nạp file, trước mọi HTML, để còn kịp chuyển hướng.
   Dùng POST kèm CSRF để không ai đăng xuất hộ người khác bằng một cái link.
   ------------------------------------------------------------------------ */
if (($_SERVER['REQUEST_METHOD'] ?? '') === 'POST' && isset($_POST['kame_logout'])) {
    $sessionToken = $_SESSION['csrf'] ?? '';
    $postedToken  = $_POST['csrf_token'] ?? '';
    if (is_string($sessionToken) && $sessionToken !== '' && is_string($postedToken) && hash_equals($sessionToken, $postedToken)) {
        clear_login_session();
    }
    if (!headers_sent()) {
        header('Location: /sv1/trangchu');
        exit;
    }
}

/** Rút gọn htmlspecialchars. */
function kame_e(?string $value): string
{
    return htmlspecialchars((string)$value, ENT_QUOTES, 'UTF-8');
}

/**
 * Gắn số phiên bản vào đường dẫn file tĩnh, lấy theo giờ sửa file lần cuối.
 *
 * Không có cái này thì trình duyệt giữ mãi bản CSS cũ trong bộ nhớ đệm: sửa
 * giao diện xong, người chơi vẫn thấy y như cũ cho tới khi họ tự xoá đệm.
 * Giờ sửa file đổi thì địa chỉ đổi theo, trình duyệt buộc phải tải lại.
 */
function kame_asset(string $duongDanWeb): string
{
    $duongDanThat = $_SERVER['DOCUMENT_ROOT'] . $duongDanWeb;
    $gio = is_file($duongDanThat) ? filemtime($duongDanThat) : false;
    return $gio === false ? $duongDanWeb : $duongDanWeb . '?v=' . $gio;
}

/**
 * Thông tin người đang đăng nhập.
 *
 * @return array{logged:bool, id:int, username:string, is_admin:bool}
 */
function kame_current_user(): array
{
    static $cache = null;
    if ($cache !== null) {
        return $cache;
    }

    $cache = ['logged' => false, 'id' => 0, 'username' => '', 'is_admin' => false];

    try {
        $accountId = (int)($_SESSION['account_id'] ?? 0);
        if ($accountId <= 0) {
            return $cache;
        }
        $pdo = kame_pdo();
        if (!$pdo) {
            return $cache;
        }
        $stmt = $pdo->prepare('SELECT id, username, is_admin FROM account WHERE id = ? LIMIT 1');
        $stmt->execute([$accountId]);
        $row = $stmt->fetch(PDO::FETCH_ASSOC);
        if ($row) {
            $cache = [
                'logged'   => true,
                'id'       => (int)$row['id'],
                'username' => (string)$row['username'],
                'is_admin' => ((int)$row['is_admin']) === 1,
            ];
        }
    } catch (Throwable $e) {
        // Lỗi DB thì coi như chưa đăng nhập, không làm vỡ trang
    }

    return $cache;
}

/** 5 mục menu chính. Khoá dùng cho tham số $active của kame_header(). */
function kame_nav_items(): array
{
    return [
        'trangchu' => ['label' => 'Trang chủ',      'href' => '/sv1/trangchu'],
        'vutru'    => ['label' => 'Vũ trụ',         'href' => '/sv1/vu-tru'],
        'tintuc'   => ['label' => 'Tin tức',        'href' => '/sv1/tin-tuc'],
        'bxh'      => ['label' => 'Bảng xếp hạng',  'href' => '/sv1/bxh01'],
        'hotro'    => ['label' => 'Hỗ trợ',         'href' => '/sv1/ho-tro'],
    ];
}

/* --- Đầu trang ----------------------------------------------------------- */

function kame_head(string $title, string $desc = ''): void
{
    $fullTitle = $title === '' ? 'Ngọc Rồng Kame' : $title . ' — Ngọc Rồng Kame';
    if ($desc === '') {
        $desc = 'Ngọc Rồng Kame — máy chủ Ngọc Rồng Online với 3 hành tinh Trái Đất, Namếc, Xayda.';
    }
    ?>
<!DOCTYPE html>
<html lang="vi">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title><?= kame_e($fullTitle) ?></title>
<meta name="description" content="<?= kame_e($desc) ?>">
<meta property="og:type" content="website">
<meta property="og:site_name" content="Ngọc Rồng Kame">
<meta property="og:title" content="<?= kame_e($fullTitle) ?>">
<meta property="og:description" content="<?= kame_e($desc) ?>">
<meta property="og:image" content="/img/brand/og.jpg">
<meta name="theme-color" content="#f06030">
<link rel="icon" type="image/png" sizes="32x32" href="/img/brand/favicon-32.png">
<link rel="icon" type="image/png" sizes="192x192" href="/img/brand/favicon-192.png">
<link rel="apple-touch-icon" href="/img/brand/favicon-180.png">
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link rel="stylesheet" href="<?= kame_e(kame_asset('/sv1/kame/ui-kit.css')) ?>">
<link rel="stylesheet" href="<?= kame_e(kame_asset('/sv1/kame/theme.css')) ?>">
</head>
<body>
    <?php
}

/* --- Header -------------------------------------------------------------- */

function kame_header(string $active): void
{
    $items = kame_nav_items();
    $user  = kame_current_user();
    ?>
<header class="kame-header">
  <div class="kame-container kame-header__inner">

    <a class="kame-header__logo" href="/sv1/trangchu">
      <img src="/img/brand/logo.png" alt="Ngọc Rồng Kame">
    </a>

    <nav class="kame-nav" aria-label="Điều hướng chính">
      <?php foreach ($items as $key => $item): ?>
        <a href="<?= kame_e($item['href']) ?>"<?= $key === $active ? ' aria-current="page"' : '' ?>><?= kame_e($item['label']) ?></a>
      <?php endforeach; ?>
    </nav>

    <div class="kame-header__actions">
      <?php if ($user['logged']): ?>
        <span class="kame-header__user kame-button--hide-sm"><?= kame_e($user['username']) ?></span>
        <?php if ($user['is_admin']): ?>
          <a class="kame-button kame-button--ghost kame-button--sm kame-button--hide-sm" href="/sv1/admin/panel">Quản trị</a>
        <?php endif; ?>
        <form method="post" style="margin:0">
          <input type="hidden" name="kame_logout" value="1">
          <input type="hidden" name="csrf_token" value="<?= kame_e(csrf_token()) ?>">
          <button type="submit" class="kame-button kame-button--ghost kame-button--sm">Đăng xuất</button>
        </form>
      <?php else: ?>
        <button type="button" class="kame-button kame-button--ghost kame-button--sm" data-kame-open="dangnhap">Đăng nhập</button>
      <?php endif; ?>

      <a class="kame-button kame-button--sm" href="/sv1/tai-game">Tải game</a>

      <button type="button" class="kame-burger" data-kame-drawer aria-label="Mở menu" aria-expanded="false">
        <span></span><span></span><span></span>
      </button>
    </div>

  </div>
</header>

<!-- Menu kéo ra, chỉ dùng ở màn hình dưới 900px -->
<div class="kame-drawer" data-kame-drawer-panel data-open="false">
  <div class="kame-drawer__panel">
    <button type="button" class="kame-button kame-button--ghost kame-button--sm kame-drawer__close" data-kame-drawer-close>Đóng</button>
    <?php foreach ($items as $key => $item): ?>
      <a href="<?= kame_e($item['href']) ?>"<?= $key === $active ? ' aria-current="page"' : '' ?>><?= kame_e($item['label']) ?></a>
    <?php endforeach; ?>
    <?php if ($user['logged']): ?>
      <a href="#" aria-disabled="true"><?= kame_e($user['username']) ?></a>
      <?php if ($user['is_admin']): ?>
        <a href="/sv1/admin/panel">Quản trị</a>
      <?php endif; ?>
    <?php endif; ?>
    <a href="/sv1/tai-game">Tải game</a>
  </div>
</div>
    <?php
}

/* --- Modal đăng nhập / đăng ký / quên mật khẩu / đổi mật khẩu ------------ */

function kame_auth_modals(): void
{
    ?>
<!-- ĐĂNG NHẬP → POST /sv1/dangnhap  (username, password) -->
<div class="kame-modal-backdrop" data-kame-modal="dangnhap" data-open="false">
  <div class="kame-modal" role="dialog" aria-modal="true" aria-label="Đăng nhập">
    <div class="kame-modal__head">
      <h2 class="kame-heading kame-heading--md">Đăng nhập</h2>
      <button type="button" class="kame-button kame-button--ghost kame-button--sm" data-kame-close>Đóng</button>
    </div>
    <div class="kame-alert kame-alert--danger" data-kame-msg><span></span><span data-kame-msg-text></span><span></span></div>
    <form data-kame-form data-endpoint="/sv1/dangnhap">
      <div class="kame-field">
        <label class="kame-label" for="dn-user">Tài khoản</label>
        <input class="kame-input" id="dn-user" name="username" autocomplete="username" required>
      </div>
      <div class="kame-field">
        <label class="kame-label" for="dn-pass">Mật khẩu</label>
        <input class="kame-input" id="dn-pass" name="password" type="password" autocomplete="current-password" required>
      </div>
      <button type="submit" class="kame-button">Đăng nhập</button>
    </form>
    <div class="kame-modal__foot">
      Chưa có tài khoản? <button type="button" data-kame-open="dangky">Đăng ký</button><br>
      <button type="button" data-kame-open="quenmatkhau">Quên mật khẩu</button> ·
      <button type="button" data-kame-open="doimatkhau">Đổi mật khẩu</button>
    </div>
  </div>
</div>

<!-- ĐĂNG KÝ → POST /sv1/dangky  (username, password, phone) -->
<div class="kame-modal-backdrop" data-kame-modal="dangky" data-open="false">
  <div class="kame-modal" role="dialog" aria-modal="true" aria-label="Đăng ký">
    <div class="kame-modal__head">
      <h2 class="kame-heading kame-heading--md">Đăng ký</h2>
      <button type="button" class="kame-button kame-button--ghost kame-button--sm" data-kame-close>Đóng</button>
    </div>
    <div class="kame-alert kame-alert--danger" data-kame-msg><span></span><span data-kame-msg-text></span><span></span></div>
    <form data-kame-form data-endpoint="/sv1/dangky">
      <div class="kame-field">
        <label class="kame-label" for="dk-user">Tài khoản</label>
        <input class="kame-input" id="dk-user" name="username" autocomplete="username" required>
      </div>
      <div class="kame-field">
        <label class="kame-label" for="dk-pass">Mật khẩu</label>
        <input class="kame-input" id="dk-pass" name="password" type="password" autocomplete="new-password" minlength="6" required>
        <span class="kame-help">Từ 6 ký tự trở lên.</span>
      </div>
      <div class="kame-field">
        <label class="kame-label" for="dk-phone">Số điện thoại</label>
        <input class="kame-input" id="dk-phone" name="phone" inputmode="numeric" autocomplete="tel" required>
        <span class="kame-help">Dùng để lấy lại mật khẩu khi quên.</span>
      </div>
      <button type="submit" class="kame-button">Tạo tài khoản</button>
    </form>
    <div class="kame-modal__foot">
      Đã có tài khoản? <button type="button" data-kame-open="dangnhap">Đăng nhập</button>
    </div>
  </div>
</div>

<!-- QUÊN MẬT KHẨU → POST /sv1/quenmatkhau  (username, phone) -->
<div class="kame-modal-backdrop" data-kame-modal="quenmatkhau" data-open="false">
  <div class="kame-modal" role="dialog" aria-modal="true" aria-label="Quên mật khẩu">
    <div class="kame-modal__head">
      <h2 class="kame-heading kame-heading--md">Quên mật khẩu</h2>
      <button type="button" class="kame-button kame-button--ghost kame-button--sm" data-kame-close>Đóng</button>
    </div>
    <div class="kame-alert kame-alert--danger" data-kame-msg><span></span><span data-kame-msg-text></span><span></span></div>
    <form data-kame-form data-endpoint="/sv1/quenmatkhau">
      <div class="kame-field">
        <label class="kame-label" for="qm-user">Tài khoản</label>
        <input class="kame-input" id="qm-user" name="username" autocomplete="username" required>
      </div>
      <div class="kame-field">
        <label class="kame-label" for="qm-phone">Số điện thoại đã đăng ký</label>
        <input class="kame-input" id="qm-phone" name="phone" inputmode="numeric" autocomplete="tel" required>
      </div>
      <button type="submit" class="kame-button">Xác minh</button>
    </form>
    <div class="kame-modal__foot">
      <button type="button" data-kame-open="dangnhap">Quay lại đăng nhập</button>
    </div>
  </div>
</div>

<!-- ĐỔI MẬT KHẨU → POST /sv1/doimatkhau  (username, oldpass, newpass, phone tuỳ chọn) -->
<div class="kame-modal-backdrop" data-kame-modal="doimatkhau" data-open="false">
  <div class="kame-modal" role="dialog" aria-modal="true" aria-label="Đổi mật khẩu">
    <div class="kame-modal__head">
      <h2 class="kame-heading kame-heading--md">Đổi mật khẩu</h2>
      <button type="button" class="kame-button kame-button--ghost kame-button--sm" data-kame-close>Đóng</button>
    </div>
    <div class="kame-alert kame-alert--danger" data-kame-msg><span></span><span data-kame-msg-text></span><span></span></div>
    <form data-kame-form data-endpoint="/sv1/doimatkhau">
      <div class="kame-field">
        <label class="kame-label" for="dm-user">Tài khoản</label>
        <input class="kame-input" id="dm-user" name="username" autocomplete="username" required>
      </div>
      <div class="kame-field">
        <label class="kame-label" for="dm-old">Mật khẩu cũ</label>
        <input class="kame-input" id="dm-old" name="oldpass" type="password" autocomplete="current-password" required>
      </div>
      <div class="kame-field">
        <label class="kame-label" for="dm-new">Mật khẩu mới</label>
        <input class="kame-input" id="dm-new" name="newpass" type="password" autocomplete="new-password" minlength="6" required>
      </div>
      <div class="kame-field">
        <label class="kame-label" for="dm-phone">Số điện thoại mới (không bắt buộc)</label>
        <input class="kame-input" id="dm-phone" name="phone" inputmode="numeric" autocomplete="tel">
        <span class="kame-help">Bỏ trống nếu giữ số cũ.</span>
      </div>
      <button type="submit" class="kame-button">Đổi mật khẩu</button>
    </form>
    <div class="kame-modal__foot">
      <button type="button" data-kame-open="dangnhap">Quay lại đăng nhập</button>
    </div>
  </div>
</div>
    <?php
}

/* --- Chân trang ---------------------------------------------------------- */

function kame_footer(): void
{
    ?>
<footer class="kame-footer">
  <div class="kame-container">
    <div class="kame-footer__grid">

      <div class="kame-footer__brand">
        <img src="/img/brand/logo.png" alt="Ngọc Rồng Kame">
        <p>Máy chủ Ngọc Rồng Online với ba hành tinh Trái Đất, Namếc và Xayda. Chơi trên điện thoại và máy tính.</p>
      </div>

      <div>
        <h3 class="kame-footer__title">Khám phá</h3>
        <ul class="kame-footer__list">
          <li><a href="/sv1/trangchu">Trang chủ</a></li>
          <li><a href="/sv1/vu-tru">Vũ trụ</a></li>
          <li><a href="/sv1/tin-tuc">Tin tức</a></li>
          <li><a href="/sv1/bxh01">Bảng xếp hạng</a></li>
        </ul>
      </div>

      <div>
        <h3 class="kame-footer__title">Hỗ trợ</h3>
        <ul class="kame-footer__list">
          <li><a href="/sv1/tai-game">Tải game</a></li>
          <li><a href="/sv1/ho-tro">Trung tâm hỗ trợ</a></li>
          <li><a href="/privacy-policy">Chính sách bảo mật</a></li>
        </ul>
      </div>

      <div>
        <h3 class="kame-footer__title">Cộng đồng</h3>
        <ul class="kame-footer__list">
          <li><a href="/sv1/ho-tro">Zalo</a></li>
          <li><a href="/sv1/ho-tro">Facebook</a></li>
          <li><a href="/sv1/ho-tro">YouTube</a></li>
          <li><a href="/sv1/ho-tro">TikTok</a></li>
        </ul>
        <h3 class="kame-footer__title" style="margin-top:24px">Về chúng tôi</h3>
        <ul class="kame-footer__list">
          <li><a href="/sv1/ho-tro">Liên hệ</a></li>
          <li><a href="/sv1/dieu-khoan">Điều khoản sử dụng</a></li>
        </ul>
      </div>

    </div>
  </div>

  <div class="kame-agewarn">
    <div class="kame-container">
      <strong>Dành cho người chơi trên 12 tuổi. Chơi quá 180 phút mỗi ngày sẽ hại sức khỏe.</strong>
      <div style="margin-top:6px">© 2026 Ngọc Rồng Kame</div>
    </div>
  </div>
</footer>

<script>
(function () {
  "use strict";

  /* --- Menu kéo ra --------------------------------------------------- */
  var drawer = document.querySelector("[data-kame-drawer-panel]");
  var burger = document.querySelector("[data-kame-drawer]");

  function setDrawer(open) {
    if (!drawer) return;
    drawer.setAttribute("data-open", open ? "true" : "false");
    if (burger) burger.setAttribute("aria-expanded", open ? "true" : "false");
  }

  if (burger) burger.addEventListener("click", function () { setDrawer(true); });
  if (drawer) {
    drawer.addEventListener("click", function (e) {
      // Bấm ra ngoài panel hoặc bấm nút Đóng thì đóng lại
      if (e.target === drawer || e.target.closest("[data-kame-drawer-close]")) setDrawer(false);
    });
  }

  /* --- Modal ---------------------------------------------------------- */
  function closeAllModals() {
    document.querySelectorAll("[data-kame-modal]").forEach(function (m) {
      m.setAttribute("data-open", "false");
    });
  }

  function openModal(name) {
    closeAllModals();
    setDrawer(false);
    var modal = document.querySelector('[data-kame-modal="' + name + '"]');
    if (!modal) return;
    modal.setAttribute("data-open", "true");
    var box = modal.querySelector("[data-kame-msg]");
    if (box) box.setAttribute("data-show", "false");
    var first = modal.querySelector("input");
    if (first) first.focus();
  }

  document.addEventListener("click", function (e) {
    var opener = e.target.closest("[data-kame-open]");
    if (opener) {
      e.preventDefault();
      openModal(opener.getAttribute("data-kame-open"));
      return;
    }
    if (e.target.closest("[data-kame-close]")) { closeAllModals(); return; }
    // Bấm vào nền tối phía ngoài hộp thoại
    if (e.target.matches("[data-kame-modal]")) closeAllModals();
  });

  document.addEventListener("keydown", function (e) {
    if (e.key === "Escape") { closeAllModals(); setDrawer(false); }
  });

  /* --- Gửi form bằng fetch tới các endpoint sẵn có --------------------- */
  function showMessage(form, ok, text) {
    var modal = form.closest("[data-kame-modal]");
    var box = modal ? modal.querySelector("[data-kame-msg]") : null;
    if (!box) return;
    box.className = "kame-alert " + (ok ? "kame-alert--success" : "kame-alert--danger");
    box.setAttribute("data-show", "true");
    var slot = box.querySelector("[data-kame-msg-text]");
    if (slot) slot.textContent = text;
  }

  document.querySelectorAll("[data-kame-form]").forEach(function (form) {
    form.addEventListener("submit", function (e) {
      e.preventDefault();

      var button = form.querySelector('button[type="submit"]');
      var label = button ? button.textContent : "";
      if (button) { button.disabled = true; button.textContent = "Đang gửi..."; }

      fetch(form.getAttribute("data-endpoint"), {
        method: "POST",
        body: new FormData(form),
        headers: { "X-Requested-With": "XMLHttpRequest", "Accept": "application/json" },
        credentials: "same-origin"
      })
        .then(function (r) { return r.text(); })
        .then(function (text) {
          var data;
          try { data = JSON.parse(text); }
          catch (err) {
            // Lớp bảo vệ chống DDoS trả về trang xác minh khi cookie hết hạn.
            // Trường hợp này chỉ cần tải lại trang là dùng tiếp được.
            if (text.indexOf("xác minh trình duyệt") !== -1 || text.indexOf("bảo vệ máy chủ") !== -1) {
              throw new Error("Phiên bảo vệ đã hết hạn. Vui lòng tải lại trang rồi thử lại.");
            }
            throw new Error("Máy chủ trả về dữ liệu không hợp lệ.");
          }

          var ok = data && data.status === "success";
          // Đăng nhập trả về success không kèm message
          showMessage(form, ok, (data && data.message) || (ok ? "Thành công!" : "Có lỗi xảy ra."));

          if (ok && form.getAttribute("data-endpoint") === "/sv1/dangnhap") {
            // Máy chủ gửi kèm "chuyen_toi" thì đi tới đó (hiện chỉ dùng cho tài
            // khoản quản trị → trang quản trị). Không có thì tải lại trang đang
            // đứng, đúng như cách cũ dành cho người chơi thường.
            //
            // Chỉ nhận đường dẫn bắt đầu bằng đúng một dấu "/". Chặn dạng
            // "//trang-la.com" và "http://..." — nếu không, chỉ cần máy chủ bị
            // chèn dữ liệu là người chơi bị ném sang trang lừa đảo ngay sau khi
            // vừa gõ mật khẩu.
            var diToi = data && typeof data.chuyen_toi === "string" ? data.chuyen_toi : "";
            if (/^\/(?!\/)/.test(diToi)) {
              window.location.href = diToi;
            } else {
              window.location.reload();
            }
          } else if (ok) {
            form.reset();
          }
        })
        .catch(function (err) {
          showMessage(form, false, err.message || "Không kết nối được máy chủ.");
        })
        .finally(function () {
          if (button) { button.disabled = false; button.textContent = label; }
        });
    });
  });
})();
</script>

</body>
</html>
    <?php
}
