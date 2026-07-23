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

if ($_SERVER['REQUEST_METHOD'] === 'POST' && isset($_POST['char_name']) && isset($_POST['amount'])) {
    require_csrf_token();
    $char_name = trim($_POST['char_name']);
    $amount = (int)$_POST['amount'];
    if ($amount <= 0) {
        $message = 'Số tiền phải > 0!';
    } else if ($char_name === '') {
        $message = 'Tên nhân vật không được để trống!';
    } else {
        $command = sprintf('BUFF_VND|%s|%d', str_replace(["\r", "\n"], '', $char_name), $amount);
        if (append_panel_command($command)) {
            $message = "Đã gửi lệnh buff cá cho nhân vật '$char_name' với số cá " . number_format($amount) . ".";
        } else {
            $message = 'Không thể ghi lệnh buff vào file điều khiển game.';
        }
    }
}
admin_head('Buff Cá');
admin_bar('Buff Cá cho Nhân Vật', 'coin');
?>
<div class="kadmin-body">
    <div class="kadmin-card" style="max-width:480px;margin:0 auto">
        <div class="kadmin-card__head"><i class="ti ti-coin" aria-hidden="true"></i> Buff cá cho nhân vật</div>
        <?php if ($message): ?>
            <div class="kame-alert <?= strpos($message, 'thành công') !== false ? 'kame-alert--success' : 'kame-alert--danger' ?>" style="margin-bottom:var(--kame-space-4)">
                <?= htmlspecialchars($message) ?>
            </div>
        <?php endif; ?>
        <form method="POST" class="kame-stack" style="gap:var(--kame-space-4)">
            <input type="hidden" name="csrf_token" value="<?= htmlspecialchars($csrfToken, ENT_QUOTES) ?>">
            <div class="kame-field">
                <label class="kame-label" for="char_name">Tên Nhân Vật:</label>
                <input class="kame-input" type="text" id="char_name" name="char_name" required>
            </div>
            <div class="kame-field">
                <label class="kame-label" for="amount">Số Cá:</label>
                <input class="kame-input" type="number" id="amount" name="amount" min="1" required>
            </div>
            <button type="submit" class="kame-button" style="width:100%"><i class="ti ti-coin" aria-hidden="true"></i> Buff Cá</button>
        </form>
        <div style="margin-top:var(--kame-space-4);text-align:center">
            <a href="panel" class="kadmin-back"><i class="ti ti-arrow-left" aria-hidden="true"></i> Quay lại Panel</a>
        </div>
    </div>
</div>
<?php
admin_foot();
