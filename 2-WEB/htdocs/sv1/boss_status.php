<?php
header('Content-Type: application/json; charset=utf-8');

$bossStatusPath = __DIR__ . '/boss_status.json';
if (!file_exists($bossStatusPath)) {
    http_response_code(500);
    echo json_encode(["error" => "Không tìm thấy file dữ liệu boss."], JSON_UNESCAPED_UNICODE);
    exit;
}

$raw = @file_get_contents($bossStatusPath);
if ($raw === false) {
    http_response_code(500);
    echo json_encode(["error" => "Không thể đọc file dữ liệu boss."], JSON_UNESCAPED_UNICODE);
    exit;
}

$data = json_decode($raw, true);
if ($data === null) {
    http_response_code(500);
    echo json_encode(["error" => "Dữ liệu boss không hợp lệ."], JSON_UNESCAPED_UNICODE);
    exit;
}

$headIds = [];
foreach ($data['bosses'] as $boss) {
    if (isset($boss['outfit0'])) {
        $headIds[] = (int)$boss['outfit0'];
    }
}
$headIds = array_values(array_unique($headIds));
$avatars = [];
if (!empty($headIds)) {
    try {
        $pdo = new PDO('mysql:host=localhost;dbname=z;charset=utf8mb4', 'root', '');
        $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
        $pdo->setAttribute(PDO::ATTR_DEFAULT_FETCH_MODE, PDO::FETCH_ASSOC);
        $placeholders = implode(',', array_fill(0, count($headIds), '?'));
        $stmt = $pdo->prepare("SELECT head_id, avatar_id FROM head_avatar WHERE head_id IN ($placeholders)");
        $stmt->execute($headIds);
        while ($row = $stmt->fetch()) {
            $avatars[(int)$row['head_id']] = (int)$row['avatar_id'];
        }
    } catch (Exception $e) {
        // Không cần dừng; vẫn trả dữ liệu boss mà không có avatar.
    }
}

foreach ($data['bosses'] as &$boss) {
    $head = isset($boss['outfit0']) ? (int)$boss['outfit0'] : 0;
    $boss['avatar_id'] = $avatars[$head] ?? 1;
}
unset($boss);

echo json_encode($data, JSON_UNESCAPED_UNICODE);
