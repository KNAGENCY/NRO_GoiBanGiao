<?php
session_start();
require_once $_SERVER['DOCUMENT_ROOT'] . "/sv1/config.php";
require_once __DIR__ . "/admin_layout.php";

if (empty($_SESSION['account_id'])) {
    header("Location: /sv1/trangchu");
    exit;
}

$account_id = (int)$_SESSION['account_id'];

$stmt = $pdo->prepare("SELECT id, is_admin FROM account WHERE id = ? LIMIT 1");
$stmt->execute([$account_id]);
$user = $stmt->fetch();

// ❌ Session fake hoặc user không tồn tại
if (!$user) {
    session_destroy();
    header("Location: /sv1/trangchu");
    exit;
}

// ❌ Không phải admin
if ((int)$user['is_admin'] !== 1) {
    echo "Không có quyền!";
    exit;
}

admin_head('Quản Lý Cửa Hàng');
admin_bar('Quản Lý Cửa Hàng', 'building-store');
?>
<div class="kadmin-body">
    <div class="kadmin-toolbar">
        <div class="kadmin-toolbar__grow"></div>
        <button type="button" class="kame-button kame-button--secondary" onclick="location.reload()">
            <i class="ti ti-refresh" aria-hidden="true"></i> Làm Mới
        </button>
    </div>

    <div class="kadmin-stats">
        <a href="ngoc" class="kadmin-card" style="text-decoration:none;display:block">
            <div class="kadmin-card__head"><i class="ti ti-diamond" aria-hidden="true"></i> Shop Vàng, Ngọc</div>
            <div class="kadmin-icon" style="width:auto;height:auto;padding:var(--kame-space-2);gap:8px;display:inline-flex;margin-bottom:var(--kame-space-3)">
                <img src="/img/930.png" onerror="this.src='/icon/930.png'" alt="ngọc">
                <img src="/img/932.png" onerror="this.src='/icon/932.png'" alt="vàng">
            </div>
            <p class="kame-muted" style="margin:0">Chỉnh sửa vật phẩm bán bằng vàng và ngọc</p>
        </a>

        <a href="/sv1/admin/vang" class="kadmin-card" style="text-decoration:none;display:block">
            <div class="kadmin-card__head"><i class="ti ti-ingot" aria-hidden="true"></i> Shop Thỏi Vàng</div>
            <div class="kadmin-icon" style="width:auto;height:auto;padding:var(--kame-space-2);gap:8px;display:inline-flex;margin-bottom:var(--kame-space-3)">
                <img src="/icon/4028.png" alt="item1">
                <img src="/icon/14116.png" alt="item2">
                <img src="/icon/15485.png" alt="item3">
                <img src="/icon/7223.png" alt="item4">
            </div>
            <p class="kame-muted" style="margin:0">Vật Phẩm Đặc Biệt &amp; Thỏi Vàng</p>
        </a>
    </div>
</div>
<?php
admin_foot();
