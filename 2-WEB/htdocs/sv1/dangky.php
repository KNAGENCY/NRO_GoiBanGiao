<?php
require_once $_SERVER['DOCUMENT_ROOT'] . "/sv1/config.php";
header('Content-Type: application/json; charset=utf-8');

if ($_SERVER["REQUEST_METHOD"] == "POST") {

    $user  = trim($_POST['username'] ?? '');
    $pass  = $_POST['password'] ?? '';
    $phone = normalize_phone($_POST['phone'] ?? ($_POST['email'] ?? ''));

    // ================= VALIDATE =================
    if (!$user || !$pass || !$phone) {
        echo json_encode([
            "status" => "error",
            "message" => "Vui lòng nhập đầy đủ thông tin!"
        ]);
        exit;
    }

    if (strlen($pass) < 6) {
        echo json_encode([
            "status" => "error",
            "message" => "Mật khẩu phải từ 6 ký tự trở lên!"
        ]);
        exit;
    }

    if (!is_valid_phone($phone)) {
        echo json_encode([
            "status" => "error",
            "message" => "Số điện thoại không hợp lệ!"
        ]);
        exit;
    }

    try {

        // ================= CHECK TRÙNG USER =================
        $stmt = $pdo->prepare("SELECT username FROM account WHERE username = ?");
        $stmt->execute([$user]);

        if ($stmt->fetch()) {
            echo json_encode([
                "status" => "error",
                "message" => "Tài khoản $user đã tồn tại!"
            ]);
            exit;
        }

        // ================= CHECK TRÙNG EMAIL =================
        $stmt = $pdo->prepare("SELECT phone FROM account WHERE phone = ?");
        $stmt->execute([$phone]);

        if ($stmt->fetch()) {
            echo json_encode([
                "status" => "error",
                "message" => "Số điện thoại $phone đã tồn tại!"
            ]);
            exit;
        }

        // ================= INSERT =================
        $stmt = $pdo->prepare("
            INSERT INTO account (username, password, phone)
            VALUES (?, ?, ?)
        ");
        $stmt->execute([$user, kame_password_hash($pass), $phone]);

        echo json_encode([
            "status" => "success",
            "message" => "Đăng Ký Thành Cong!"
        ]);

    } catch (PDOException $e) {
        echo json_encode([
            "status" => "error",
            "message" => "Lỗi hệ thống, vui lòng thử lại!"
        ]);
    }
}
?>
