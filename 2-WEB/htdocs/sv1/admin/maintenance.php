<?php
session_start();
require_once $_SERVER['DOCUMENT_ROOT'] . "/sv1/config.php";
require_once __DIR__ . "/game_bridge.php";
require_once __DIR__ . "/admin_layout.php";

if (!isset($_SESSION['account_id'])) {
    header("Location: /sv1/trangchu");
    exit;
}

$account_id = $_SESSION['account_id'];

$stmt = $pdo->prepare("SELECT is_admin FROM account WHERE id = :id LIMIT 1");
$stmt->bindParam(":id", $account_id);
$stmt->execute();
$acc = $stmt->fetch(PDO::FETCH_ASSOC);

if (!$acc || $acc['is_admin'] != 1) {
    header("Location: /sv1/trangchu");
    exit;
}

$message = '';
$csrfToken = ensure_csrf_token();

if ($_SERVER['REQUEST_METHOD'] === 'POST' && isset($_POST['start_maintenance'])) {
    require_csrf_token();
    $command = 'MAINTENANCE_NOW';
    if (append_panel_command($command)) {
        $message = 'Lệnh bảo trì đã được gửi, server sẽ thực hiện bảo trì trong game.';
    } else {
        $message = 'Không thể ghi lệnh bảo trì vào file điều khiển game.';
    }
}

admin_head('Bảo Trì');
admin_bar('Bảo Trì Server', 'tool');
?>
<div class="kadmin-body">
    <div class="kadmin-card" style="max-width:520px;margin:0 auto">
        <div class="kadmin-card__head"><i class="ti ti-alert-triangle" aria-hidden="true"></i> Bật chế độ bảo trì</div>
        <p class="kame-muted" style="margin:0 0 var(--kame-space-4)">
            Nhấn nút dưới để bắt đầu chế độ bảo trì. Tất cả người chơi sẽ bị kick và server sẽ offline để bảo trì.
        </p>
        <?php if ($message): ?>
            <div class="kame-alert kame-alert--success" style="margin-bottom:var(--kame-space-4)">
                <?= htmlspecialchars($message) ?>
            </div>
        <?php endif; ?>
        <form method="POST">
            <input type="hidden" name="csrf_token" value="<?= htmlspecialchars($csrfToken, ENT_QUOTES) ?>">
            <button type="submit" name="start_maintenance" class="kame-button kame-button--danger" style="width:100%">
                <i class="ti ti-player-power" aria-hidden="true"></i> Bắt Đầu Bảo Trì
            </button>
        </form>
        <div style="margin-top:var(--kame-space-4);text-align:center">
            <a href="/sv1/admin/panel" class="kadmin-back"><i class="ti ti-arrow-left" aria-hidden="true"></i> Quay lại Panel</a>
        </div>
    </div>
</div>
<?php
admin_foot();
