<?php
session_start();
require_once $_SERVER['DOCUMENT_ROOT'] . "/sv1/config.php";

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    exit("Invalid request");
}

if (!isset($_SESSION['account_id'])) {
    exit("Not login");
}

require_csrf_token();

$stmt = $pdo->prepare("SELECT is_admin FROM account WHERE id = ?");
$stmt->execute([$_SESSION['account_id']]);

$me = $stmt->fetch(PDO::FETCH_ASSOC);

if (!$me || (int)$me['is_admin'] !== 1) {
    exit("Không có quyền");
}

try {

    $pdo->beginTransaction();

    $pdo->exec("
        DELETE p FROM player p
        INNER JOIN account a ON p.account_id = a.id
        WHERE a.is_admin = 0
    ");

    $pdo->exec("DELETE FROM account WHERE is_admin = 0");

    $pdo->exec("DELETE FROM history_items_diemdanh");
    $pdo->exec("DELETE FROM history_transaction");
    $pdo->exec("DELETE FROM clan");
    $pdo->exec("DELETE FROM shop_ky_gui");

    $pdo->commit();

    file_put_contents(
        __DIR__ . '/reset_log.txt',
        date('Y-m-d H:i:s') .
        " | ADMIN ID: " . $_SESSION['account_id'] .
        " | IP: " . ($_SERVER['REMOTE_ADDR'] ?? 'unknown') .
        PHP_EOL,
        FILE_APPEND
    );

    echo "✅ RESET THÀNH CÔNG";

} catch (Exception $e) {

    if ($pdo->inTransaction()) {
        $pdo->rollBack();
    }

    echo "❌ ERROR";
}
