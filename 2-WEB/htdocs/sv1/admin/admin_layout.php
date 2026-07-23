<?php
/* ==========================================================================
   KAME ADMIN — khung chung cho MOI trang quan tri.

   Cach dung trong tung trang admin (sau khi da kiem dang nhap + is_admin):

       require_once __DIR__ . '/admin_layout.php';
       admin_head('Bao Tri');                       // mo <html><head><body> + shell
       admin_bar('Bao Tri', 'ti-tool');             // thanh header nau + nut quay lai
       echo '<div class="kadmin-body"> ... noi dung ... </div>';
       admin_foot();                                // dong shell + </body></html>

   KHONG trang nao duoc tu viet <html><head><style> nua. Sua khung o day la
   ca 20 trang doi theo. Chi dung class trong ui-kit.css + kame-admin.css.
   ========================================================================== */

if (!function_exists('kadmin_e')) {
    function kadmin_e(string $s): string { return htmlspecialchars($s, ENT_QUOTES, 'UTF-8'); }
}

/** Mo <html><head> nap dung 3 file CSS Kame + mo khung .kadmin. */
function admin_head(string $title): void
{
    $full = ($title === '' ? 'Quan tri' : $title) . ' — Ngọc Rồng Kame';
    ?>
<!DOCTYPE html>
<html lang="vi">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title><?= kadmin_e($full) ?></title>
<meta name="robots" content="noindex,nofollow">
<link rel="icon" type="image/png" sizes="32x32" href="/img/brand/favicon-32.png">
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link rel="stylesheet" href="/sv1/kame/ui-kit.css">
<link rel="stylesheet" href="/sv1/kame/theme.css">
<link rel="stylesheet" href="/sv1/admin/kame-admin.css">
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@tabler/icons-webfont@3.1.0/dist/tabler-icons.min.css">
</head>
<body>
<div class="kadmin">
    <?php
}

/**
 * Thanh header nau cua trang admin.
 * @param string $title  Ten trang (vd "Bao Tri")
 * @param string $icon   Ten icon Tabler khong co tien to (vd "tool", "users")
 * @param string $back   Duong dan nut quay lai (mac dinh ve panel)
 */
function admin_bar(string $title, string $icon = 'layout-dashboard', string $back = '/sv1/admin/panel'): void
{
    ?>
    <div class="kadmin-bar">
        <h1 class="kadmin-bar__title"><i class="ti ti-<?= kadmin_e($icon) ?>" aria-hidden="true"></i><?= kadmin_e($title) ?></h1>
        <div class="kadmin-bar__actions">
            <a class="kame-button kame-button--secondary" href="<?= kadmin_e($back) ?>"><i class="ti ti-arrow-left" aria-hidden="true"></i> Quay lại Panel</a>
        </div>
    </div>
    <?php
}

/** Dong khung + </body></html>. */
function admin_foot(): void
{
    ?>
</div>
</body>
</html>
    <?php
}
