<?php
/* ==========================================================================
   NGỌC RỒNG KAME — TẦNG DỮ LIỆU DÙNG CHUNG
   Mọi trang trong sv1/kame/ gọi các hàm ở đây, KHÔNG tự viết SQL.
   Nguyên tắc: mỗi hàm tự bọc try/catch, lỗi thì trả mảng rỗng / ok=false,
   tuyệt đối không để trang trắng.
   ========================================================================== */

declare(strict_types=1);

// Phiên đăng nhập phải mở trước config.php vì config.php có kiểm tra hết hạn phiên.
if (session_status() !== PHP_SESSION_ACTIVE && !headers_sent()) {
    session_start();
}

require_once $_SERVER['DOCUMENT_ROOT'] . '/sv1/config.php';

/* --- Hằng đường dẫn ------------------------------------------------------ */

/**
 * Thư mục game root. Thử lần lượt cho tới khi thấy thư mục có thật.
 * `C:/zz1` là chỗ đặt trên VPS — đường dẫn ngắn, nằm ngoài hồ sơ người dùng
 * nên không hỏng khi đăng nhập bằng tài khoản khác.
 */
const KAME_GAME_ROOTS = [
    'C:/zz1',
    'C:/Users/Admin/Music/zz1',
    'C:/Users/Administrator/Music/zz1',
];

/** File trạng thái boss do server Java ghi ra mỗi giây. */
const KAME_BOSS_FILE = 'C:/xampp/htdocs/sv1/boss_status.json';

/** Quá số giây này coi như dữ liệu boss đã cũ (server có thể đang tắt). */
const KAME_STALE_SECONDS = 60;

/* --- Tiện ích nội bộ ----------------------------------------------------- */

/** Trả về PDO dùng chung từ config.php, hoặc null nếu kết nối hỏng. */
function kame_pdo(): ?PDO
{
    global $pdo;
    return ($pdo instanceof PDO) ? $pdo : null;
}

/** Đọc và giải mã một file JSON. Trả null nếu thiếu file hoặc JSON hỏng. */
function kame_read_json(string $path): ?array
{
    try {
        if (!is_file($path) || !is_readable($path)) {
            return null;
        }
        $raw = @file_get_contents($path);
        if ($raw === false || $raw === '') {
            return null;
        }
        $data = json_decode($raw, true);
        return is_array($data) ? $data : null;
    } catch (Throwable $e) {
        return null;
    }
}

/** Tìm file panel_server_status.json trong các game root đã biết. */
function kame_find_status_file(): ?string
{
    foreach (KAME_GAME_ROOTS as $root) {
        $file = $root . '/panel_server_status.json';
        if (is_file($file)) {
            return $file;
        }
    }
    return null;
}

/**
 * Tra avatar_id cho một danh sách head_id qua bảng head_avatar.
 * Làm đúng như sv1/boss_status.php: thiếu dữ liệu thì mặc định 1.
 *
 * @param  int[] $headIds
 * @return array<int,int> head_id => avatar_id
 */
function kame_head_avatars(array $headIds): array
{
    $map = [];
    try {
        $headIds = array_values(array_unique(array_map('intval', $headIds)));
        if (!$headIds) {
            return $map;
        }
        $pdo = kame_pdo();
        if (!$pdo) {
            return $map;
        }
        $holes = implode(',', array_fill(0, count($headIds), '?'));
        $stmt = $pdo->prepare("SELECT head_id, avatar_id FROM head_avatar WHERE head_id IN ($holes)");
        $stmt->execute($headIds);
        while ($row = $stmt->fetch(PDO::FETCH_ASSOC)) {
            $map[(int)$row['head_id']] = (int)$row['avatar_id'];
        }
    } catch (Throwable $e) {
        return $map;
    }
    return $map;
}

/** Lấy một phần tử trong chuỗi JSON dạng mảng, an toàn với dữ liệu hỏng. */
function kame_json_at(?string $json, int $index, int $default = 0): int
{
    if (!is_string($json) || $json === '') {
        return $default;
    }
    $arr = json_decode($json, true);
    if (!is_array($arr) || !isset($arr[$index]) || !is_numeric($arr[$index])) {
        return $default;
    }
    return (int)$arr[$index];
}

/**
 * Sức mạnh đệ tử. Cột `pet` là JSON lồng: phần tử [1] lại là một chuỗi JSON,
 * trong đó [1] mới là sức mạnh. Giữ nguyên logic của sv1/bxh01.php.
 */
function kame_pet_power(?string $pet): int
{
    if (!is_string($pet) || $pet === '') {
        return 0;
    }
    $outer = json_decode($pet, true);
    if (!is_array($outer) || !isset($outer[1]) || !is_string($outer[1])) {
        return 0;
    }
    $info = json_decode($outer[1], true);
    if (!is_array($info) || !isset($info[1]) || !is_numeric($info[1])) {
        return 0;
    }
    return (int)$info[1];
}

/** Tên hành tinh suy từ cột gender — đối chiếu consts/ConstPlayer.java. */
function kame_planet_name(int $gender): string
{
    return match ($gender) {
        0 => 'Trái Đất',   // ConstPlayer.TRAI_DAT
        1 => 'Namếc',      // ConstPlayer.NAMEC
        2 => 'Xayda',      // ConstPlayer.XAYDA
        default => 'Trái Đất',
    };
}

/* --- 1. Trạng thái server ------------------------------------------------ */

/**
 * Số liệu THẬT của server. Không bịa: đọc không được thì ok = false.
 *
 * @return array{online:int, accounts:int, players:int, ok:bool}
 */
function kame_server_status(): array
{
    $out = ['online' => 0, 'accounts' => 0, 'players' => 0, 'ok' => false];

    try {
        $file = kame_find_status_file();
        if ($file !== null) {
            $data = kame_read_json($file);
            if ($data !== null && isset($data['online'])) {
                $out['online'] = max(0, (int)$data['online']);
                $out['ok'] = true;
            }
        }
    } catch (Throwable $e) {
        // Giữ ok = false
    }

    try {
        $pdo = kame_pdo();
        if ($pdo) {
            $out['accounts'] = (int)$pdo->query('SELECT COUNT(*) FROM account')->fetchColumn();
            $out['players']  = (int)$pdo->query('SELECT COUNT(*) FROM player')->fetchColumn();
        }
    } catch (Throwable $e) {
        // Đếm hỏng thì để 0, không đụng tới cờ ok của file trạng thái
    }

    return $out;
}

/* --- 2. Boss ------------------------------------------------------------- */

/**
 * Danh sách boss kèm avatar lấy từ game.
 *
 * @param  int|null $limit  null = lấy hết
 * @return array<int,array{name:string,status:string,state:string,hp:int,maxHp:int,mapName:string,respawnIn:int,avatar_id:int}>
 */
function kame_bosses(?int $limit = null): array
{
    try {
        $data = kame_read_json(KAME_BOSS_FILE);
        if ($data === null || !isset($data['bosses']) || !is_array($data['bosses'])) {
            return [];
        }

        $rows = $data['bosses'];

        // Boss đang sống (ACTIVE) xếp lên trước cho dễ nhìn
        usort($rows, static function ($a, $b) {
            $sa = (($a['status'] ?? '') === 'ACTIVE') ? 0 : 1;
            $sb = (($b['status'] ?? '') === 'ACTIVE') ? 0 : 1;
            if ($sa !== $sb) {
                return $sa <=> $sb;
            }
            return ((int)($b['maxHp'] ?? 0)) <=> ((int)($a['maxHp'] ?? 0));
        });

        if ($limit !== null && $limit > 0) {
            $rows = array_slice($rows, 0, $limit);
        }

        $heads = [];
        foreach ($rows as $boss) {
            if (isset($boss['outfit0'])) {
                $heads[] = (int)$boss['outfit0'];
            }
        }
        $avatars = kame_head_avatars($heads);

        $out = [];
        foreach ($rows as $boss) {
            $head = isset($boss['outfit0']) ? (int)$boss['outfit0'] : 0;
            $out[] = [
                'name'      => (string)($boss['name'] ?? 'Không rõ'),
                'status'    => (string)($boss['status'] ?? 'AFK'),
                'state'     => (string)($boss['state'] ?? ''),
                'hp'        => (int)($boss['hp'] ?? 0),
                'maxHp'     => max(1, (int)($boss['maxHp'] ?? 1)),
                'mapName'   => (string)($boss['mapName'] ?? ''),
                'respawnIn' => (int)($boss['respawnIn'] ?? 0),
                'avatar_id' => $avatars[$head] ?? 1,
            ];
        }
        return $out;
    } catch (Throwable $e) {
        return [];
    }
}

/**
 * Độ tươi của dữ liệu boss.
 *
 * @return array{ok:bool, stale:bool, age:int}  age tính bằng giây
 */
function kame_boss_meta(): array
{
    try {
        $data = kame_read_json(KAME_BOSS_FILE);
        if ($data === null) {
            return ['ok' => false, 'stale' => true, 'age' => 0];
        }

        // Server Java ghi generatedAt theo mili-giây; thiếu thì lấy giờ sửa file.
        if (isset($data['generatedAt']) && is_numeric($data['generatedAt'])) {
            $generated = (int)round(((float)$data['generatedAt']) / 1000);
        } else {
            $mtime = @filemtime(KAME_BOSS_FILE);
            $generated = ($mtime === false) ? 0 : (int)$mtime;
        }

        $age = max(0, time() - $generated);
        return ['ok' => true, 'stale' => $age > KAME_STALE_SECONDS, 'age' => $age];
    } catch (Throwable $e) {
        return ['ok' => false, 'stale' => true, 'age' => 0];
    }
}

/* --- 3. Bảng xếp hạng ---------------------------------------------------- */

/**
 * Top sức mạnh. Sức mạnh nằm ở data_point[1], cộng thêm sức mạnh đệ tử —
 * đúng cách sv1/bxh01.php đang tính cột "Tổng".
 *
 * @return array<int,array{name:string,power:int,avatar_id:int,planet:string}>
 */
function kame_top_power(int $limit = 8): array
{
    try {
        $pdo = kame_pdo();
        if (!$pdo) {
            return [];
        }
        $limit = max(1, $limit);

        $rows = $pdo->query('SELECT name, gender, head, data_point, pet FROM player')->fetchAll(PDO::FETCH_ASSOC);
        if (!$rows) {
            return [];
        }

        $list = [];
        foreach ($rows as $row) {
            $power = kame_json_at($row['data_point'] ?? null, 1) + kame_pet_power($row['pet'] ?? null);
            $list[] = [
                'name'   => (string)($row['name'] ?? ''),
                'power'  => $power,
                'head'   => (int)($row['head'] ?? 0),
                'gender' => (int)($row['gender'] ?? 0),
            ];
        }

        usort($list, static fn($a, $b) => $b['power'] <=> $a['power']);
        $list = array_slice($list, 0, $limit);

        $avatars = kame_head_avatars(array_column($list, 'head'));

        $out = [];
        foreach ($list as $item) {
            $out[] = [
                'name'      => $item['name'],
                'power'     => $item['power'],
                'avatar_id' => $avatars[$item['head']] ?? 1,
                'planet'    => kame_planet_name($item['gender']),
            ];
        }
        return $out;
    } catch (Throwable $e) {
        return [];
    }
}

/**
 * Top vàng. Vàng là phần tử [0] của cột data_inventory — theo sv1/bxh02.php.
 *
 * @return array<int,array{name:string,gold:int,avatar_id:int}>
 */
function kame_top_wealth(int $limit = 8): array
{
    try {
        $pdo = kame_pdo();
        if (!$pdo) {
            return [];
        }
        $limit = max(1, $limit);

        $rows = $pdo->query('SELECT name, head, data_inventory FROM player')->fetchAll(PDO::FETCH_ASSOC);
        if (!$rows) {
            return [];
        }

        $list = [];
        foreach ($rows as $row) {
            $list[] = [
                'name' => (string)($row['name'] ?? ''),
                'gold' => kame_json_at($row['data_inventory'] ?? null, 0),
                'head' => (int)($row['head'] ?? 0),
            ];
        }

        usort($list, static fn($a, $b) => $b['gold'] <=> $a['gold']);
        $list = array_slice($list, 0, $limit);

        $avatars = kame_head_avatars(array_column($list, 'head'));

        $out = [];
        foreach ($list as $item) {
            $out[] = [
                'name'      => $item['name'],
                'gold'      => $item['gold'],
                'avatar_id' => $avatars[$item['head']] ?? 1,
            ];
        }
        return $out;
    } catch (Throwable $e) {
        return [];
    }
}

/* --- 4. Tin tức ---------------------------------------------------------- */

/**
 * Danh sách bài viết, mới nhất trước.
 *
 * @return array<int,array{id:int,slug:string,title:string,excerpt:string,body:string,image:?string,tag:string,created_at:string}>
 */
function kame_news(int $limit, int $offset = 0): array
{
    try {
        $pdo = kame_pdo();
        if (!$pdo) {
            return [];
        }
        // LIMIT/OFFSET không nhận tham số bind ở chế độ emulate — ép kiểu số nguyên.
        $limit  = max(1, min(100, $limit));
        $offset = max(0, $offset);

        $stmt = $pdo->query(
            "SELECT id, slug, title, excerpt, body, image, tag, created_at
             FROM news ORDER BY created_at DESC, id DESC
             LIMIT $limit OFFSET $offset"
        );
        return $stmt->fetchAll(PDO::FETCH_ASSOC) ?: [];
    } catch (Throwable $e) {
        return [];
    }
}

/** Tổng số bài viết, dùng cho phân trang. */
function kame_news_count(): int
{
    try {
        $pdo = kame_pdo();
        if (!$pdo) {
            return 0;
        }
        return (int)$pdo->query('SELECT COUNT(*) FROM news')->fetchColumn();
    } catch (Throwable $e) {
        return 0;
    }
}

/** Một bài viết theo slug, null nếu không có. */
function kame_news_one(string $slug): ?array
{
    try {
        $pdo = kame_pdo();
        if (!$pdo || $slug === '') {
            return null;
        }
        $stmt = $pdo->prepare(
            'SELECT id, slug, title, excerpt, body, image, tag, created_at
             FROM news WHERE slug = ? LIMIT 1'
        );
        $stmt->execute([$slug]);
        $row = $stmt->fetch(PDO::FETCH_ASSOC);
        return $row ?: null;
    } catch (Throwable $e) {
        return null;
    }
}

/* --- 5. Ảnh -------------------------------------------------------------- */

/** Đường dẫn icon lấy thẳng từ dữ liệu game (data/icon/x4). */
function kame_icon_url(int $id): string
{
    // Bỏ đuôi .php: .htaccess ép mọi URL .php mở bằng GET phải chuyển hướng 302
    // sang URL sạch, khiến MỖI ảnh tốn 2 lượt gọi thay vì 1.
    return '/icon_asset?file=' . rawurlencode((string)max(0, $id)) . '.png';
}

/**
 * Ô giữ chỗ cho ảnh chưa có trong sv1/kame/assets/.
 * Trả về HTML sẵn sàng echo.
 */
function kame_placeholder(string $file, string $ratio = '16/9'): string
{
    $safeFile  = htmlspecialchars($file, ENT_QUOTES, 'UTF-8');
    // Chỉ cho phép dạng "số/số" để không chèn được CSS lạ vào thuộc tính style.
    $safeRatio = preg_match('#^\d+(?:\.\d+)?\s*/\s*\d+(?:\.\d+)?$#', trim($ratio)) ? trim($ratio) : '16/9';

    // Ảnh đã có thật trong sv1/kame/assets/ thì dùng luôn, khỏi vẽ ô xám.
    // Nhờ vậy chỉ cần chép ảnh vào thư mục là mọi trang tự hiện, không phải sửa code.
    $duongDan = __DIR__ . DIRECTORY_SEPARATOR . 'assets' . DIRECTORY_SEPARATOR . basename($file);
    if (basename($file) === $file && is_file($duongDan)) {
        return '<img class="kame-ph-img" src="/sv1/kame/assets/' . rawurlencode($file) . '"'
             . ' style="aspect-ratio:' . $safeRatio . ';width:100%;object-fit:cover;display:block"'
             . ' alt="" loading="lazy" decoding="async">';
    }

    return '<div class="kame-ph" style="aspect-ratio:' . $safeRatio . '" role="img" aria-label="Ảnh chưa có: ' . $safeFile . '">'
         . '<div>'
         . '<div class="kame-ph__name">' . $safeFile . '</div>'
         . '<div class="kame-ph__note">Ảnh đang chuẩn bị</div>'
         . '</div>'
         . '</div>';
}
