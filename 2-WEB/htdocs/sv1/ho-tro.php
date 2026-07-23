<?php
/* ==========================================================================
   NGỌC RỒNG KAME — CỘNG ĐỒNG & HỖ TRỢ
   Gồm 4 phần: thẻ cộng đồng, ô tìm kiếm lọc FAQ, FAQ theo nhóm,
   và khu bình luận cộng đồng (dùng lại bảng site_comments có sẵn).

   KHÔNG có form gửi ticket vì phía máy chủ chưa có chỗ lưu và chỗ trả lời.
   Thà không có còn hơn có một cái nút gửi vào hư không.
   ========================================================================== */

declare(strict_types=1);

require_once __DIR__ . '/kame/layout.php';   // tự kéo theo data.php + config.php + session

/* Trang này quay về chính nó sau mỗi lần đăng bài. Ghi cứng cho chắc,
   không lấy từ REQUEST_URI để khỏi ai chèn được thứ lạ vào header. */
const HT_URL = '/sv1/ho-tro';

/** Mức sức mạnh tối thiểu mới được đăng bài — chặn nick rác vừa tạo. */
const HT_MIN_POWER = 10000000000;   // 10 tỷ

/** Trần bài viết và trả lời mỗi tài khoản mỗi ngày. */
const HT_MAX_POST_PER_DAY  = 3;
const HT_MAX_REPLY_PER_DAY = 100;

/* --------------------------------------------------------------------------
   1. BẢNG DỮ LIỆU BÌNH LUẬN
   Giữ nguyên tên bảng và tên cột của bản cũ, không đổi schema.
   -------------------------------------------------------------------------- */

function ht_ensure_comment_tables(PDO $pdo): void
{
    try {
        $pdo->exec("CREATE TABLE IF NOT EXISTS site_comments (
            id INT AUTO_INCREMENT PRIMARY KEY,
            parent_id INT NULL,
            section ENUM('admin','player') NOT NULL DEFAULT 'player',
            account_id INT NULL,
            author_name VARCHAR(80) NOT NULL DEFAULT 'conga',
            author_icon INT NOT NULL DEFAULT 1474,
            is_anonymous TINYINT(1) NOT NULL DEFAULT 0,
            is_admin TINYINT(1) NOT NULL DEFAULT 0,
            content TEXT NOT NULL,
            status ENUM('visible','hidden') NOT NULL DEFAULT 'visible',
            created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
            INDEX idx_section_parent (section, parent_id, created_at),
            INDEX idx_parent (parent_id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4");

        $pdo->exec("CREATE TABLE IF NOT EXISTS site_comment_likes (
            id INT AUTO_INCREMENT PRIMARY KEY,
            comment_id INT NOT NULL,
            account_id INT NULL,
            client_hash CHAR(64) NOT NULL,
            created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
            UNIQUE KEY uniq_comment_client (comment_id, client_hash),
            INDEX idx_comment_id (comment_id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4");
    } catch (Throwable $e) {
        // Bảng đã có sẵn từ trước; lỗi ở đây không được làm trắng trang.
    }
}

/**
 * Người đang đăng nhập, kèm nhân vật mạnh nhất của tài khoản đó.
 * Sức mạnh lấy ở data_point[1] — đúng như bản cũ, không cộng đệ tử.
 *
 * @return array{account_id:int,is_logged:bool,is_admin:bool,player_name:string,power:int,has_player:bool,icon:int}
 */
function ht_forum_user(PDO $pdo): array
{
    $out = [
        'account_id' => 0, 'is_logged' => false, 'is_admin' => false,
        'player_name' => 'conga', 'power' => 0, 'has_player' => false, 'icon' => 1474,
    ];

    try {
        $accountId = (int)($_SESSION['account_id'] ?? 0);
        if ($accountId <= 0) {
            return $out;
        }

        $stmt = $pdo->prepare('SELECT id, username, is_admin FROM account WHERE id = ? LIMIT 1');
        $stmt->execute([$accountId]);
        $account = $stmt->fetch(PDO::FETCH_ASSOC);
        if (!$account) {
            return $out;
        }

        $out['account_id'] = $accountId;
        $out['is_logged']  = true;
        $out['is_admin']   = !empty($account['is_admin']);

        $stmt = $pdo->prepare('SELECT name, gender, data_point FROM player WHERE account_id = ? ORDER BY id ASC LIMIT 1');
        $stmt->execute([$accountId]);
        $player = $stmt->fetch(PDO::FETCH_ASSOC);

        $gender = 0;
        if ($player) {
            $out['has_player'] = true;
            $name = trim((string)($player['name'] ?? ''));
            if ($name !== '') {
                $out['player_name'] = $name;
            }
            $gender = (int)($player['gender'] ?? 0);
            $out['power'] = kame_json_at($player['data_point'] ?? null, 1);
        }

        // Icon đầu nhân vật theo hành tinh; admin có icon riêng.
        $out['icon'] = $out['is_admin'] ? 16103 : match ($gender) {
            2 => 1475,   // Xayda
            1 => 1476,   // Namếc
            default => 1474,   // Trái Đất
        };
    } catch (Throwable $e) {
        return $out;
    }

    return $out;
}

/** Vân tay người bấm Thích — để một người không bấm được hai lần. */
function ht_client_hash(): string
{
    return hash('sha256',
        ($_SESSION['account_id'] ?? 'khach') . '|' .
        ($_SERVER['REMOTE_ADDR'] ?? '') . '|' .
        ($_SERVER['HTTP_USER_AGENT'] ?? '')
    );
}

/** Đếm số bài (hoặc số trả lời) tài khoản đã đăng từ 0h hôm nay. */
function ht_daily_count(PDO $pdo, int $accountId, bool $isReply): int
{
    try {
        if ($accountId <= 0) {
            return 0;
        }
        $where = $isReply ? 'parent_id IS NOT NULL' : 'parent_id IS NULL';
        $stmt = $pdo->prepare(
            "SELECT COUNT(*) FROM site_comments
             WHERE account_id = ? AND section = 'player' AND $where AND created_at >= CURDATE()"
        );
        $stmt->execute([$accountId]);
        return (int)$stmt->fetchColumn();
    } catch (Throwable $e) {
        return 0;
    }
}

/** "cách đây 5 phút" — dễ đọc hơn ngày giờ đầy đủ. */
function ht_relative_time($datetime): string
{
    $ts = strtotime((string)$datetime);
    if (!$ts) {
        return '';
    }
    $diff = max(0, time() - $ts);
    if ($diff < 60) {
        return 'vừa xong';
    }
    $minutes = intdiv($diff, 60);
    if ($minutes < 60) {
        return 'cách đây ' . $minutes . ' phút';
    }
    $hours = intdiv($minutes, 60);
    if ($hours < 24) {
        return 'cách đây ' . $hours . ' giờ';
    }
    $days = intdiv($hours, 24);
    if ($days < 30) {
        return 'cách đây ' . $days . ' ngày';
    }
    $months = intdiv($days, 30);
    if ($months < 12) {
        return 'cách đây ' . $months . ' tháng';
    }
    return 'cách đây ' . intdiv($months, 12) . ' năm';
}

/** Thời điểm hoạt động gần nhất của một bài, tính cả trả lời con. */
function ht_latest_activity(array $comment, array $children): int
{
    $latest = strtotime((string)($comment['created_at'] ?? '')) ?: 0;
    foreach ($children[(int)($comment['id'] ?? 0)] ?? [] as $child) {
        $latest = max($latest, ht_latest_activity($child, $children));
    }
    return $latest;
}

/* --------------------------------------------------------------------------
   2. XỬ LÝ FORM — phải chạy TRƯỚC mọi HTML vì có chuyển hướng
   -------------------------------------------------------------------------- */

$pdo         = kame_pdo();
$htNotice    = '';
$htNoticeOk  = false;
$htUser      = ['account_id' => 0, 'is_logged' => false, 'is_admin' => false,
                'player_name' => 'conga', 'power' => 0, 'has_player' => false, 'icon' => 1474];

if ($pdo) {
    ht_ensure_comment_tables($pdo);
    $htUser = ht_forum_user($pdo);
}

if ($pdo && ($_SERVER['REQUEST_METHOD'] ?? '') === 'POST' && isset($_POST['comment_action'])) {
    require_csrf_token();
    $action = (string)$_POST['comment_action'];

    /* --- Đăng bài mới hoặc trả lời --- */
    if ($action === 'create_comment') {
        $content   = mb_substr(trim((string)($_POST['content'] ?? '')), 0, 1000, 'UTF-8');
        $parentId  = max(0, (int)($_POST['parent_id'] ?? 0));
        // Chỉ quản trị viên mới được đăng dưới danh nghĩa Ban quản trị. Ô tick của
        // người thường có gửi lên cũng bị bỏ qua vì luôn nhân với quyền thật.
        $asAdmin   = !empty($_POST['as_admin']) && $htUser['is_admin'];
        // Thông báo của Ban quản trị luôn đứng tên, không cho ẩn danh.
        $anonymous = !empty($_POST['anonymous']) && !$asAdmin;

        if ($content === '') {
            $htNotice = 'Bạn chưa nhập nội dung.';
        } else {
            try {
                // Kiểm tra bài cha có thật không. Không trả lời vào bài của admin.
                if ($parentId > 0) {
                    $stmt = $pdo->prepare("SELECT id, is_admin FROM site_comments WHERE id = ? AND status = 'visible' LIMIT 1");
                    $stmt->execute([$parentId]);
                    $parent = $stmt->fetch(PDO::FETCH_ASSOC);
                    if (!$parent || !empty($parent['is_admin'])) {
                        $parentId = 0;
                    }
                }
                $isReply = $parentId > 0;

                if (!$htUser['is_logged']) {
                    $htNotice = 'Bạn cần đăng nhập trước khi đăng bài.';
                } elseif (!$htUser['is_admin'] && (!$htUser['has_player'] || $htUser['power'] < HT_MIN_POWER)) {
                    $htNotice = 'Nhân vật cần trên 10 tỷ sức mạnh mới được đăng bài hoặc trả lời.';
                } elseif (!$htUser['is_admin'] && !$isReply && ht_daily_count($pdo, $htUser['account_id'], false) >= HT_MAX_POST_PER_DAY) {
                    $htNotice = 'Mỗi tài khoản chỉ được đăng ' . HT_MAX_POST_PER_DAY . ' bài mỗi ngày.';
                } elseif (!$htUser['is_admin'] && $isReply && ht_daily_count($pdo, $htUser['account_id'], true) >= HT_MAX_REPLY_PER_DAY) {
                    $htNotice = 'Mỗi tài khoản chỉ được trả lời ' . HT_MAX_REPLY_PER_DAY . ' lần mỗi ngày.';
                } else {
                    $author = $anonymous ? 'Ẩn danh' : (string)$htUser['player_name'];
                    $icon   = $anonymous ? 1474 : (int)$htUser['icon'];
                    $stmt = $pdo->prepare(
                        'INSERT INTO site_comments
                         (parent_id, section, account_id, author_name, author_icon, is_anonymous, is_admin, content)
                         VALUES (?, ?, ?, ?, ?, ?, ?, ?)'
                    );
                    $stmt->execute([
                        $parentId ?: null, $asAdmin ? 'admin' : 'player', $htUser['account_id'],
                        $author, $icon, $anonymous ? 1 : 0, $asAdmin ? 1 : 0, $content,
                    ]);
                    header('Location: ' . HT_URL . '#binh-luan');
                    exit;
                }
            } catch (Throwable $e) {
                $htNotice = 'Không lưu được bài viết. Vui lòng thử lại.';
            }
        }

    /* --- Sửa bài: chỉ admin --- */
    } elseif ($action === 'update_comment') {
        $commentId = max(0, (int)($_POST['comment_id'] ?? 0));
        $content   = mb_substr(trim((string)($_POST['content'] ?? '')), 0, 1000, 'UTF-8');
        if (!$htUser['is_admin']) {
            $htNotice = 'Chỉ quản trị viên được chỉnh sửa bài viết.';
        } elseif ($commentId > 0 && $content !== '') {
            try {
                $stmt = $pdo->prepare("UPDATE site_comments SET content = ? WHERE id = ? AND status = 'visible'");
                $stmt->execute([$content, $commentId]);
                header('Location: ' . HT_URL . '#binh-luan-' . $commentId);
                exit;
            } catch (Throwable $e) {
                $htNotice = 'Không sửa được bài viết.';
            }
        }

    /* --- Xoá bài: admin xoá được tất cả, người chơi chỉ xoá bài mình --- */
    } elseif ($action === 'delete_comment') {
        $commentId = max(0, (int)($_POST['comment_id'] ?? 0));
        if ($commentId > 0) {
            try {
                $stmt = $pdo->prepare("SELECT id, account_id, is_admin FROM site_comments WHERE id = ? AND status = 'visible' LIMIT 1");
                $stmt->execute([$commentId]);
                $target = $stmt->fetch(PDO::FETCH_ASSOC);

                $canDelete = $target && (
                    $htUser['is_admin'] ||
                    (empty($target['is_admin']) && $htUser['is_logged'] && (int)$target['account_id'] === $htUser['account_id'])
                );

                if (!$canDelete) {
                    $htNotice = 'Bạn không có quyền xoá bài viết này.';
                } else {
                    // Xoá cả cây trả lời bên dưới, không để lại bài mồ côi.
                    $ids = [$commentId];
                    $queue = [$commentId];
                    while ($queue) {
                        $current = array_pop($queue);
                        $stmt = $pdo->prepare('SELECT id FROM site_comments WHERE parent_id = ?');
                        $stmt->execute([$current]);
                        foreach ($stmt->fetchAll(PDO::FETCH_COLUMN) as $childId) {
                            $ids[] = (int)$childId;
                            $queue[] = (int)$childId;
                        }
                    }
                    $holes = implode(',', array_fill(0, count($ids), '?'));
                    $pdo->prepare("DELETE FROM site_comment_likes WHERE comment_id IN ($holes)")->execute($ids);
                    $pdo->prepare("DELETE FROM site_comments WHERE id IN ($holes)")->execute($ids);
                    header('Location: ' . HT_URL . '#binh-luan');
                    exit;
                }
            } catch (Throwable $e) {
                $htNotice = 'Không xoá được bài viết.';
            }
        }

    /* --- Thích / bỏ thích --- */
    } elseif ($action === 'like') {
        $commentId = max(0, (int)($_POST['comment_id'] ?? 0));
        if ($commentId > 0) {
            try {
                $hash = ht_client_hash();
                $stmt = $pdo->prepare('SELECT id FROM site_comment_likes WHERE comment_id = ? AND client_hash = ? LIMIT 1');
                $stmt->execute([$commentId, $hash]);
                $existing = $stmt->fetchColumn();
                if ($existing) {
                    $pdo->prepare('DELETE FROM site_comment_likes WHERE id = ?')->execute([$existing]);
                } else {
                    $pdo->prepare('INSERT IGNORE INTO site_comment_likes (comment_id, account_id, client_hash) VALUES (?, ?, ?)')
                        ->execute([$commentId, $htUser['account_id'] ?: null, $hash]);
                }
                header('Location: ' . HT_URL . '#binh-luan-' . $commentId);
                exit;
            } catch (Throwable $e) {
                $htNotice = 'Không ghi nhận được lượt thích.';
            }
        }
    }
}

/* --------------------------------------------------------------------------
   3. ĐỌC BÌNH LUẬN
   -------------------------------------------------------------------------- */

$htRoots    = [];   // bài gốc, mới hoạt động nhất lên đầu
$htChildren = [];   // parent_id => danh sách trả lời
$htLikes    = [];   // comment_id => số lượt thích

if ($pdo) {
    try {
        $rows = $pdo->query("SELECT * FROM site_comments WHERE status = 'visible' ORDER BY created_at ASC, id ASC")
                    ->fetchAll(PDO::FETCH_ASSOC) ?: [];

        foreach ($rows as $row) {
            $parentId = (int)($row['parent_id'] ?? 0);
            if ($parentId > 0) {
                $htChildren[$parentId][] = $row;
            } else {
                $htRoots[] = $row;
            }
        }

        // Bài admin luôn nằm trên, còn lại xếp theo hoạt động gần nhất.
        usort($htRoots, static function ($a, $b) use ($htChildren) {
            $adminA = empty($a['is_admin']) ? 1 : 0;
            $adminB = empty($b['is_admin']) ? 1 : 0;
            if ($adminA !== $adminB) {
                return $adminA <=> $adminB;
            }
            return ht_latest_activity($b, $htChildren) <=> ht_latest_activity($a, $htChildren);
        });

        foreach ($pdo->query('SELECT comment_id, COUNT(*) total FROM site_comment_likes GROUP BY comment_id') as $row) {
            $htLikes[(int)$row['comment_id']] = (int)$row['total'];
        }
    } catch (Throwable $e) {
        $htRoots = [];
        $htChildren = [];
        $htLikes = [];
    }
}

$htTotalComments = count($htRoots);
foreach ($htChildren as $group) {
    $htTotalComments += count($group);
}

/** Vẽ một bài viết cùng toàn bộ trả lời bên dưới. */
function ht_render_comment(array $comment, array $children, array $likes, array $user, int $depth = 0): string
{
    $id       = (int)$comment['id'];
    $isAdmin  = !empty($comment['is_admin']);
    $isOwner  = !empty($user['is_logged']) && (int)($comment['account_id'] ?? 0) === (int)$user['account_id'];
    $canEdit  = !empty($user['is_admin']);
    $canDel   = !empty($user['is_admin']) || (!$isAdmin && $isOwner);
    $replies  = $children[$id] ?? [];
    $likeCnt  = (int)($likes[$id] ?? 0);
    $csrf     = kame_e(csrf_token());

    $icon = $isAdmin ? 16103 : (int)($comment['author_icon'] ?? 1474);
    if (!empty($comment['is_anonymous'])) {
        $icon = 1474;
    }

    if ($isAdmin) {
        $badge = '<span class="kame-badge kame-badge--hot">Quản trị</span>';
    } elseif (!empty($comment['is_anonymous'])) {
        $badge = '<span class="kame-badge">Ẩn danh</span>';
    } else {
        $badge = '<span class="kame-badge kame-badge--info">Người chơi</span>';
    }

    $created  = (string)($comment['created_at'] ?? '');
    $fullTime = $created !== '' ? date('d/m/Y H:i', (int)(strtotime($created) ?: time())) : '';

    $html  = '<article class="ht-cmt' . ($isAdmin ? ' ht-cmt--admin' : '') . '" id="binh-luan-' . $id . '">';
    $html .= '<img class="ht-cmt__avatar" src="' . kame_e(kame_icon_url($icon)) . '" alt="" loading="lazy" width="40" height="40">';
    $html .= '<div class="ht-cmt__body">';

    $html .= '<div class="ht-cmt__head">'
           . '<strong>' . kame_e((string)($comment['author_name'] ?? '')) . '</strong>'
           . $badge
           . '<time datetime="' . kame_e($created) . '" title="' . kame_e($fullTime) . '">'
           . kame_e(ht_relative_time($created)) . '</time>'
           . '</div>';

    $html .= '<div class="ht-cmt__text">' . nl2br(kame_e((string)$comment['content'])) . '</div>';

    /* Hàng nút thao tác */
    $html .= '<div class="ht-cmt__actions">';
    $html .= '<form method="post">'
           . '<input type="hidden" name="csrf_token" value="' . $csrf . '">'
           . '<input type="hidden" name="comment_action" value="like">'
           . '<input type="hidden" name="comment_id" value="' . $id . '">'
           . '<button type="submit">Thích (' . $likeCnt . ')</button>'
           . '</form>';

    if (!$isAdmin && $depth === 0) {
        $html .= '<button type="button" data-ht-toggle="tra-loi-' . $id . '">Trả lời</button>';
    }
    if ($canEdit) {
        $html .= '<button type="button" data-ht-toggle="sua-' . $id . '">Sửa</button>';
    }
    if ($canDel) {
        $html .= '<form method="post" onsubmit="return confirm(\'Xoá bài viết này và toàn bộ trả lời bên dưới?\');">'
               . '<input type="hidden" name="csrf_token" value="' . $csrf . '">'
               . '<input type="hidden" name="comment_action" value="delete_comment">'
               . '<input type="hidden" name="comment_id" value="' . $id . '">'
               . '<button type="submit" class="ht-cmt__danger">Xoá</button>'
               . '</form>';
    }
    $html .= '</div>';

    /* Ô sửa bài, ẩn sẵn */
    if ($canEdit) {
        $html .= '<form method="post" class="ht-cmt__form" id="sua-' . $id . '" hidden>'
               . '<input type="hidden" name="csrf_token" value="' . $csrf . '">'
               . '<input type="hidden" name="comment_action" value="update_comment">'
               . '<input type="hidden" name="comment_id" value="' . $id . '">'
               . '<textarea class="kame-textarea" name="content" rows="3" maxlength="1000" required>'
               . kame_e((string)$comment['content']) . '</textarea>'
               . '<button type="submit" class="kame-button kame-button--sm">Lưu</button>'
               . '</form>';
    }

    /* Ô trả lời, ẩn sẵn */
    if (!$isAdmin && $depth === 0) {
        $html .= '<form method="post" class="ht-cmt__form" id="tra-loi-' . $id . '" hidden>'
               . '<input type="hidden" name="csrf_token" value="' . $csrf . '">'
               . '<input type="hidden" name="comment_action" value="create_comment">'
               . '<input type="hidden" name="parent_id" value="' . $id . '">'
               . '<textarea class="kame-textarea" name="content" rows="3" maxlength="1000" placeholder="Viết trả lời..." required></textarea>'
               . '<label class="ht-cmt__anon"><input type="checkbox" name="anonymous" value="1"> Gửi ẩn danh</label>'
               . '<button type="submit" class="kame-button kame-button--sm">Gửi trả lời</button>'
               . '</form>';
    }

    if ($replies) {
        $html .= '<div class="ht-cmt__replies">';
        foreach ($replies as $reply) {
            $html .= ht_render_comment($reply, $children, $likes, $user, $depth + 1);
        }
        $html .= '</div>';
    }

    return $html . '</div></article>';
}

/* --------------------------------------------------------------------------
   4. NỘI DUNG TĨNH: kênh cộng đồng + câu hỏi thường gặp
   -------------------------------------------------------------------------- */

/* Link tạm, chờ user đưa địa chỉ thật của từng kênh. */
const HT_LINK_TAM = 'https://nrokame.com/lien-he';

$htChannels = [
    ['ten' => 'Zalo',     'mo_ta' => 'Nhóm trò chuyện nhanh, hỏi gì đáp nấy trong ngày.',        'nut' => 'Tham gia'],
    ['ten' => 'Facebook', 'mo_ta' => 'Trang thông báo sự kiện, bảo trì và cập nhật máy chủ.',    'nut' => 'Tham gia'],
    ['ten' => 'YouTube',  'mo_ta' => 'Video hướng dẫn người mới và ghi hình các trận đánh boss.', 'nut' => 'Tham gia'],
    ['ten' => 'TikTok',   'mo_ta' => 'Clip ngắn mẹo chơi, sự kiện và khoảnh khắc trong game.',   'nut' => 'Tham gia'],
];

/* FAQ: mỗi nhóm là một kame-accordion. Trường "tim" là từ khoá không dấu
   để ô tìm kiếm bên trên lọc được kể cả khi người dùng gõ không dấu. */
$htFaq = [
    [
        'nhom' => 'Tài khoản',
        'cau'  => [
            [
                'hoi' => 'Đăng ký tài khoản như thế nào?',
                'dap' => 'Bấm nút <strong>Đăng nhập</strong> ở góc trên bên phải rồi chọn <strong>Đăng ký</strong>. Bạn cần tên tài khoản, mật khẩu từ 6 ký tự trở lên và một số điện thoại. Số điện thoại chỉ dùng để lấy lại mật khẩu khi quên, nên hãy điền số bạn đang dùng thật.',
                'tim' => 'dang ky tao tai khoan mat khau so dien thoai',
            ],
            [
                'hoi' => 'Quên mật khẩu thì lấy lại bằng cách nào?',
                'dap' => 'Mở ô <strong>Đăng nhập</strong> rồi chọn <strong>Quên mật khẩu</strong>. Nhập đúng tên tài khoản và số điện thoại đã đăng ký để hệ thống xác minh. Nếu số điện thoại không khớp, bạn cần liên hệ qua kênh cộng đồng chính thức và cung cấp thông tin nhân vật để được hỗ trợ.',
                'tim' => 'quen mat khau lay lai khoi phuc xac minh',
            ],
            [
                'hoi' => 'Đổi mật khẩu ở đâu?',
                'dap' => 'Trong ô <strong>Đăng nhập</strong> có mục <strong>Đổi mật khẩu</strong>. Bạn phải nhập được mật khẩu cũ thì mới đổi được. Nhân tiện có thể cập nhật luôn số điện thoại mới, bỏ trống nếu vẫn giữ số cũ.',
                'tim' => 'doi mat khau thay doi password',
            ],
            [
                'hoi' => 'Tài khoản bị khoá thì phải làm sao?',
                'dap' => 'Hãy liên hệ qua một trong các kênh cộng đồng chính thức phía trên, kèm theo tên tài khoản và tên nhân vật. Ban quản trị sẽ kiểm tra lý do khoá và phản hồi. Đừng nhờ người lạ "mở khoá hộ" — đó gần như luôn là lừa đảo.',
                'tim' => 'bi khoa tai khoan ban cam choi mo khoa',
            ],
            [
                'hoi' => 'Có nên cho người khác mượn tài khoản không?',
                'dap' => 'Không. Ban quản trị <strong>không bao giờ</strong> hỏi mật khẩu của bạn, dù với bất kỳ lý do gì. Mọi mất mát do bạn tự đưa tài khoản cho người khác sẽ không được khôi phục.',
                'tim' => 'cho muon tai khoan lua dao bao mat chia se mat khau',
            ],
        ],
    ],
    [
        'nhom' => 'Cài đặt &amp; Kết nối',
        'cau'  => [
            [
                'hoi' => 'Tải game ở đâu?',
                'dap' => 'Vào mục <a class="kame-link" href="/sv1/tai-game">Tải game</a> trên thanh menu. Ở đó có đủ các bản cho từng loại thiết bị. Chỉ tải từ trang chính thức này, các bản tải trôi nổi bên ngoài có thể bị cài thêm mã độc.',
                'tim' => 'tai game download cai dat ban tai apk pc',
            ],
            [
                'hoi' => 'Mở game nhưng đứng mãi ở màn hình đăng nhập?',
                'dap' => 'Thường là do mạng chập chờn hoặc máy chủ đang bảo trì. Hãy thử: đổi giữa Wi-Fi và 4G, tắt hẳn game rồi mở lại, kiểm tra mục <a class="kame-link" href="/sv1/tin-tuc">Tin tức</a> xem có thông báo bảo trì không. Nếu vẫn không được, báo lại ở kênh cộng đồng kèm ảnh chụp màn hình.',
                'tim' => 'khong vao duoc game dung man hinh dang nhap loi ket noi bao tri',
            ],
            [
                'hoi' => 'Đang chơi thì bị mất kết nối, nhân vật có sao không?',
                'dap' => 'Không sao. Toàn bộ dữ liệu nhân vật nằm trên máy chủ chứ không nằm ở máy bạn, nên khi đăng nhập lại mọi thứ vẫn nguyên. Chỉ cần lưu ý là nếu đang đánh boss mà rớt mạng thì lượt đánh đó có thể mất.',
                'tim' => 'mat ket noi rot mang disconnect mat du lieu nhan vat',
            ],
            [
                'hoi' => 'Máy cấu hình yếu có chơi được không?',
                'dap' => 'Game chạy được trên cả điện thoại và máy tính, không đòi hỏi máy mạnh. Nếu máy cũ và thấy giật, hãy giảm bớt hiệu ứng trong phần cài đặt của game và tránh đứng ở những khu vực quá đông người chơi.',
                'tim' => 'cau hinh may yeu giat lag toi thieu dien thoai may tinh',
            ],
            [
                'hoi' => 'Máy chủ bảo trì thì có báo trước không?',
                'dap' => 'Có. Mọi đợt bảo trì có kế hoạch đều được đăng ở mục <a class="kame-link" href="/sv1/tin-tuc">Tin tức</a> và các kênh cộng đồng trước giờ tắt máy chủ. Với sự cố đột xuất, chúng tôi thông báo ngay khi xử lý xong.',
                'tim' => 'bao tri may chu thong bao lich tat server',
            ],
        ],
    ],
    [
        'nhom' => 'Trong game',
        'cau'  => [
            [
                'hoi' => 'Làm sao lên sức mạnh nhanh?',
                'dap' => 'Máy chủ đang để <strong>tỉ lệ kinh nghiệm 10x</strong>, nên cách nhanh nhất vẫn là đánh quái đúng khu vực với cấp độ của mình, làm nhiệm vụ theo tuyến và tận dụng thời gian tu luyện khi không online. Đừng bỏ qua nhiệm vụ — phần thưởng của nó thường hơn hẳn việc cắm đánh quái.',
                'tim' => 'len suc manh nhanh kinh nghiem exp 10x cay quai tu luyen',
            ],
            [
                'hoi' => 'Nhận nhiệm vụ ở đâu?',
                'dap' => 'Nói chuyện với các NPC ở làng khởi đầu của hành tinh bạn chọn. Nhiệm vụ được mở dần theo cấp độ, làm xong nhiệm vụ này thì nhiệm vụ tiếp theo mới hiện ra. Bảng nhiệm vụ trong game luôn ghi rõ bước kế tiếp cần làm gì.',
                'tim' => 'nhiem vu quest npc lam nhiem vu tuyen',
            ],
            [
                'hoi' => 'Bang hội hoạt động thế nào?',
                'dap' => 'Bạn có thể xin vào một bang có sẵn hoặc tự lập bang khi đủ điều kiện. Vào bang để cùng đánh boss, tham gia các sự kiện theo nhóm và nhận thưởng chung. Trưởng bang là người duyệt thành viên và phân quyền trong bang.',
                'tim' => 'bang hoi clan guild lap bang tham gia',
            ],
            [
                'hoi' => 'Có bao nhiêu boss và bản đồ?',
                'dap' => 'Máy chủ hiện có <strong>98 boss</strong> rải trên <strong>183 bản đồ</strong>. Bạn xem được boss nào đang xuất hiện, ở đâu, còn bao nhiêu máu ngay tại <a class="kame-link" href="/sv1/trangchu">Trang chủ</a> — danh sách này cập nhật liên tục theo máy chủ game.',
                'tim' => 'boss ban do map bao nhieu 98 183 san boss o dau',
            ],
            [
                'hoi' => 'Các sự kiện diễn ra vào giờ nào?',
                'dap' => 'Lịch cố định hằng ngày:<ul>'
                       . '<li><strong>Mabư</strong> — 12h00 và 22h00</li>'
                       . '<li><strong>Black Ball War</strong> — 20h00 đến 21h00</li>'
                       . '<li><strong>Map War</strong> — 21h29</li>'
                       . '<li><strong>Đại hội võ thuật</strong> — mỗi giờ một lần, phải đăng ký trước phút thứ 25</li>'
                       . '</ul>Nên vào sớm vài phút, đặc biệt với Đại hội võ thuật vì quá phút 25 là không đăng ký được nữa.',
                'tim' => 'su kien theo gio lich mabu maburu black ball war map war dai hoi vo thuat gio nao 12h 22h 20h 21h29 phut 25',
            ],
            [
                'hoi' => 'Tỉ lệ rơi đồ ra sao?',
                'dap' => 'Máy chủ để <strong>tỉ lệ rơi đồ 5x</strong> so với bản gốc. Đồ tốt vẫn phụ thuộc may mắn, nhưng đánh boss theo nhóm và tham gia sự kiện theo giờ sẽ cho cơ hội nhận đồ hiếm cao hơn hẳn đánh quái thường.',
                'tim' => 'ti le roi do drop rate 5x do hiem vat pham',
            ],
        ],
    ],
    [
        'nhom' => 'Nạp &amp; Vật phẩm',
        'cau'  => [
            [
                'hoi' => 'Nạp thẻ bằng cách nào?',
                'dap' => 'Chỉ nạp qua trang chính thức <strong>nrokame.com</strong> hoặc mục nạp thẻ ngay trong game. Nếu có ai đó nhắn tin riêng nhận "nạp hộ giá rẻ", đó là lừa đảo — chúng tôi không có đại lý hay cộng tác viên nào cả.',
                'tim' => 'nap the cach nap thanh toan nap tien dai ly',
            ],
            [
                'hoi' => 'Nạp rồi nhưng chưa nhận được vật phẩm?',
                'dap' => 'Trước hết hãy thoát game rồi vào lại — phần lớn trường hợp vật phẩm về ngay sau đó vài phút. Nếu quá 30 phút vẫn chưa có, hãy liên hệ kênh cộng đồng kèm <strong>tên tài khoản, tên nhân vật, thời điểm nạp và mã giao dịch</strong>. Thiếu mã giao dịch sẽ rất khó tra.',
                'tim' => 'nap roi khong nhan duoc do vat pham mat tien ma giao dich',
            ],
            [
                'hoi' => 'Giftcode nhập ở đâu?',
                'dap' => 'Nhập trực tiếp trong game ở mục Giftcode. Mỗi mã chỉ dùng được một lần cho mỗi tài khoản và đều có hạn sử dụng. Nếu báo lỗi, hãy kiểm tra xem mã đã hết hạn, đã dùng rồi, hay bạn gõ nhầm chữ hoa chữ thường.',
                'tim' => 'giftcode gift code nhap ma qua tang khuyen mai hết hạn',
            ],
            [
                'hoi' => 'Mua nhầm vật phẩm có đổi lại được không?',
                'dap' => 'Không. Vật phẩm đã mua trong game không hoàn lại được, kể cả khi bạn bấm nhầm. Hãy đọc kỹ mô tả và giá trước khi xác nhận.',
                'tim' => 'mua nham hoan tra doi lai vat pham refund',
            ],
            [
                'hoi' => 'Giao dịch với người chơi khác có an toàn không?',
                'dap' => 'Chỉ giao dịch bằng chức năng giao dịch có sẵn trong game, nơi cả hai bên cùng xác nhận. Tuyệt đối không "đưa đồ trước, trả sau" qua lời hứa ngoài game — mất đồ theo cách đó thì không khôi phục được.',
                'tim' => 'giao dich trao doi nguoi choi lua dao an toan mua ban',
            ],
        ],
    ],
];

/* --------------------------------------------------------------------------
   5. HIỂN THỊ
   -------------------------------------------------------------------------- */

kame_head(
    'Cộng đồng & Hỗ trợ',
    'Kênh cộng đồng chính thức, câu hỏi thường gặp và nơi trao đổi của người chơi Ngọc Rồng Kame.'
);
kame_header('hotro');
?>

<style>
/* --- Băng đầu trang: dùng lại khung của kame-hero, chỉ đổi ảnh nền ------- */
.ht-band { align-items: center; min-height: clamp(240px, 34vh, 380px); }
.ht-band .kame-hero__content { max-width: 720px; }

/* --- Thẻ kênh cộng đồng -------------------------------------------------- */
.ht-kenh { display: flex; flex-direction: column; height: 100%; padding: var(--kame-space-5); }
.ht-kenh__ten {
  margin-bottom: 8px;
  color: var(--kame-maroon-700);
  font-family: "Roboto Condensed", "Arial Narrow", sans-serif;
  font-size: 26px; font-style: italic; font-weight: 900;
  text-transform: uppercase;
}
.ht-kenh__mo-ta { flex: 1; margin-bottom: var(--kame-space-4); color: rgb(40 21 15 / 80%); font-size: 15px; }

/* --- Ô tìm kiếm FAQ ------------------------------------------------------ */
.ht-tim { display: flex; flex-wrap: wrap; gap: var(--kame-space-3); align-items: center; }
.ht-tim .kame-input { flex: 1 1 260px; }
.ht-goi-y { display: flex; flex-wrap: wrap; gap: 8px; margin-top: var(--kame-space-4); }
.ht-goi-y button {
  padding: 6px 14px;
  border: 1px solid rgb(96 0 24 / 35%); border-radius: 999px;
  background: rgb(255 243 213 / 80%);
  color: var(--kame-maroon-700);
  font: inherit; font-size: 14px; font-weight: 600;
  cursor: pointer;
}
.ht-goi-y button:hover { border-color: var(--kame-orange-600); color: var(--kame-orange-700); }

.ht-faq-nhom { margin-bottom: var(--kame-space-6); }
.ht-faq-nhom__ten {
  margin-bottom: var(--kame-space-3);
  color: var(--kame-maroon-700);
  font-family: "Roboto Condensed", "Arial Narrow", sans-serif;
  font-size: 22px; font-style: italic; font-weight: 900;
  text-transform: uppercase;
}
.ht-faq-nhom[hidden] { display: none; }
.kame-accordion details > div ul { margin: 8px 0 0; padding-left: 20px; }
.kame-accordion details > div li { margin-bottom: 4px; }
.ht-khong-thay { display: none; }
.ht-khong-thay[data-hien="true"] { display: block; }

/* --- Bình luận ----------------------------------------------------------- */
.ht-cmt {
  display: grid; grid-template-columns: 40px 1fr; gap: var(--kame-space-4);
  padding: var(--kame-space-4);
  border: var(--kame-border); border-radius: var(--kame-radius-md);
  background: rgb(255 243 213 / 88%);
}
.ht-cmt + .ht-cmt { margin-top: var(--kame-space-4); }
.ht-cmt--admin { border-color: var(--kame-orange-600); border-width: 2px; }
.ht-cmt__avatar {
  width: 40px; height: 40px;
  border: 1px solid rgb(96 0 24 / 30%); border-radius: var(--kame-radius-sm);
  background: var(--kame-paper-300);
  object-fit: contain;
  image-rendering: pixelated;
}
.ht-cmt__head { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; margin-bottom: 6px; }
.ht-cmt__head strong { color: var(--kame-maroon-700); }
.ht-cmt__head time { color: rgb(40 21 15 / 55%); font-size: 13px; }
.ht-cmt__text { color: var(--kame-ink); overflow-wrap: anywhere; }
.ht-cmt__actions { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; margin-top: var(--kame-space-3); }
.ht-cmt__actions form { margin: 0; }
.ht-cmt__actions button {
  padding: 4px 12px;
  border: 1px solid rgb(96 0 24 / 30%); border-radius: var(--kame-radius-sm);
  background: transparent;
  color: var(--kame-maroon-700);
  font: inherit; font-size: 13px; font-weight: 600;
  cursor: pointer;
}
.ht-cmt__actions button:hover { border-color: var(--kame-orange-600); color: var(--kame-orange-700); }
.ht-cmt__danger:hover { border-color: var(--kame-red) !important; color: var(--kame-red) !important; }
.ht-cmt__form { display: grid; gap: 8px; margin-top: var(--kame-space-3); justify-items: start; }
.ht-cmt__form .kame-textarea { width: 100%; }
.ht-cmt__anon { display: flex; gap: 6px; align-items: center; font-size: 14px; }
.ht-cmt__replies { margin-top: var(--kame-space-4); padding-left: var(--kame-space-4); border-left: 3px solid rgb(240 96 48 / 40%); }
.ht-dang-bai { display: grid; gap: var(--kame-space-3); justify-items: start; }
.ht-dang-bai .kame-textarea { width: 100%; }

@media (max-width: 600px) {
  .ht-cmt { grid-template-columns: 32px 1fr; gap: var(--kame-space-3); padding: var(--kame-space-3); }
  .ht-cmt__avatar { width: 32px; height: 32px; }
  .ht-cmt__replies { padding-left: var(--kame-space-3); }
}
</style>

<!-- ===== BĂNG ĐẦU TRANG ================================================== -->
<section class="kame-hero ht-band" style="background-image:url('/img/brand/band-congdong.jpg')">
  <div class="kame-hero__inner">
    <div class="kame-container">
      <div class="kame-hero__content">
        <nav class="kame-breadcrumb" aria-label="Đường dẫn" style="color:rgb(255 253 245 / 82%);margin-bottom:12px">
          <a href="/sv1/trangchu" style="color:var(--kame-gold)">Trang chủ</a>
          <span>Cộng đồng &amp; Hỗ trợ</span>
        </nav>
        <h1 class="kame-heading kame-heading--xl kame-hero__title">Cộng đồng &amp; <em>Hỗ trợ</em></h1>
        <p class="kame-hero__desc">
          Gặp khó ở đâu thì tìm ở đây. Bên dưới là các kênh liên lạc chính thức,
          những câu hỏi người chơi hay gặp nhất, và chỗ để bạn trao đổi với nhau.
        </p>
      </div>
    </div>
  </div>
</section>

<!-- ===== KÊNH CỘNG ĐỒNG ================================================== -->
<section class="kame-section">
  <div class="kame-container">
    <div class="kame-section-title">
      <h2 class="kame-heading kame-heading--lg">Kênh cộng đồng chính thức</h2>
      <span class="kame-muted">Ngoài các kênh này, mọi nơi khác đều không phải của chúng tôi.</span>
    </div>

    <div class="kame-grid">
      <?php foreach ($htChannels as $kenh): ?>
        <div class="kame-col-3">
          <div class="kame-card ht-kenh">
            <h3 class="ht-kenh__ten"><?= kame_e($kenh['ten']) ?></h3>
            <p class="ht-kenh__mo-ta"><?= kame_e($kenh['mo_ta']) ?></p>
            <a class="kame-button kame-button--sm" href="<?= kame_e(HT_LINK_TAM) ?>" rel="noopener">
              <?= kame_e($kenh['nut']) ?>
            </a>
          </div>
        </div>
      <?php endforeach; ?>
    </div>

    <div class="kame-alert kame-alert--danger" data-show="true" style="margin-top:var(--kame-space-5)">
      <span></span>
      <span>
        <strong>Cảnh giác lừa đảo.</strong>
        Ban quản trị không bao giờ hỏi mật khẩu của bạn và không nhờ ai nạp hộ.
        Mọi lời mời "mua vàng giá rẻ", "mở khoá tài khoản có phí" đều là giả mạo.
      </span>
      <span></span>
    </div>
  </div>
</section>

<!-- ===== TÌM KIẾM + FAQ ================================================== -->
<section class="kame-section" id="cau-hoi">
  <div class="kame-container">
    <div class="kame-section-title">
      <h2 class="kame-heading kame-heading--lg">Câu hỏi thường gặp</h2>
      <span class="kame-muted">Gõ từ khoá để lọc nhanh, có dấu hay không dấu đều được.</span>
    </div>

    <div class="kame-card" style="padding:var(--kame-space-5);margin-bottom:var(--kame-space-6)">
      <div class="ht-tim">
        <label class="kame-label" for="ht-o-tim" style="margin:0">Tìm trợ giúp</label>
        <input class="kame-input" id="ht-o-tim" type="search" autocomplete="off"
               placeholder="Ví dụ: quên mật khẩu, nạp thẻ, boss...">
      </div>
      <div class="ht-goi-y" aria-label="Từ khoá gợi ý">
        <button type="button" data-ht-goi-y="quên mật khẩu">Quên mật khẩu</button>
        <button type="button" data-ht-goi-y="nạp thẻ">Nạp thẻ</button>
        <button type="button" data-ht-goi-y="giftcode">Giftcode</button>
        <button type="button" data-ht-goi-y="boss">Boss</button>
        <button type="button" data-ht-goi-y="sự kiện">Sự kiện theo giờ</button>
        <button type="button" data-ht-goi-y="tải game">Tải game</button>
      </div>
      <p class="kame-help" id="ht-ket-qua" style="margin-top:var(--kame-space-3)" aria-live="polite"></p>
    </div>

    <?php foreach ($htFaq as $nhom): ?>
      <div class="ht-faq-nhom" data-ht-nhom>
        <h3 class="ht-faq-nhom__ten"><?= $nhom['nhom'] /* tên nhóm là HTML tĩnh do chính trang này viết */ ?></h3>
        <div class="kame-accordion">
          <?php foreach ($nhom['cau'] as $cau): ?>
            <details data-ht-cau
                     data-ht-tim="<?= kame_e(mb_strtolower($cau['hoi'] . ' ' . strip_tags($cau['dap']) . ' ' . $cau['tim'], 'UTF-8')) ?>">
              <summary><?= kame_e($cau['hoi']) ?></summary>
              <div><?= $cau['dap'] /* nội dung HTML tĩnh do chính trang này viết */ ?></div>
            </details>
          <?php endforeach; ?>
        </div>
      </div>
    <?php endforeach; ?>

    <div class="kame-alert kame-alert--info ht-khong-thay" id="ht-khong-thay" data-show="true">
      <span></span>
      <span>Không có câu hỏi nào khớp với từ khoá này. Bạn thử gõ ngắn gọn hơn, hoặc hỏi trực tiếp ở phần bình luận bên dưới.</span>
      <span></span>
    </div>
  </div>
</section>

<!-- ===== BÌNH LUẬN CỘNG ĐỒNG ============================================= -->
<section class="kame-section" id="binh-luan">
  <div class="kame-container">
    <div class="kame-section-title">
      <h2 class="kame-heading kame-heading--lg">Người chơi trao đổi</h2>
      <span class="kame-muted">
        <?= $htTotalComments > 0 ? kame_e((string)$htTotalComments) . ' bài viết và trả lời' : 'Chưa có bài viết nào' ?>
      </span>
    </div>

    <?php if ($htNotice !== ''): ?>
      <div class="kame-alert <?= $htNoticeOk ? 'kame-alert--success' : 'kame-alert--danger' ?>" data-show="true" style="margin-bottom:var(--kame-space-5)">
        <span></span><span><?= kame_e($htNotice) ?></span><span></span>
      </div>
    <?php endif; ?>

    <?php if (!$pdo): ?>
      <div class="kame-alert kame-alert--danger" data-show="true">
        <span></span><span>Chưa có dữ liệu — không kết nối được cơ sở dữ liệu.</span><span></span>
      </div>

    <?php else: ?>

      <!-- Ô đăng bài mới -->
      <div class="kame-card" style="padding:var(--kame-space-5);margin-bottom:var(--kame-space-6)">
        <?php if (!$htUser['is_logged']): ?>
          <p class="kame-muted" style="margin-bottom:var(--kame-space-4)">
            Bạn cần đăng nhập mới viết được bài. Để chống spam, nhân vật phải có trên 10 tỷ sức mạnh.
          </p>
          <button type="button" class="kame-button" data-kame-open="dangnhap">Đăng nhập để viết bài</button>

        <?php elseif (!$htUser['is_admin'] && (!$htUser['has_player'] || $htUser['power'] < HT_MIN_POWER)): ?>
          <p class="kame-muted">
            Nhân vật của bạn cần trên <strong>10 tỷ sức mạnh</strong> mới được đăng bài hoặc trả lời.
            Đây là cách đơn giản nhất để chặn nick rác vừa tạo vào spam.
            Bạn vẫn xem và bấm thích được mọi bài viết.
          </p>

        <?php else: ?>
          <form method="post" class="ht-dang-bai">
            <input type="hidden" name="csrf_token" value="<?= kame_e(csrf_token()) ?>">
            <input type="hidden" name="comment_action" value="create_comment">
            <input type="hidden" name="parent_id" value="0">
            <div class="kame-field" style="width:100%">
              <label class="kame-label" for="ht-noi-dung">Viết bài mới</label>
              <textarea class="kame-textarea" id="ht-noi-dung" name="content" rows="4" maxlength="1000"
                        placeholder="Hỏi han, góp ý hoặc chia sẻ kinh nghiệm..." required></textarea>
              <span class="kame-help">
                Tối đa 1000 ký tự. Mỗi ngày <?= HT_MAX_POST_PER_DAY ?> bài và <?= HT_MAX_REPLY_PER_DAY ?> lượt trả lời.
              </span>
            </div>
            <label class="ht-cmt__anon"><input type="checkbox" name="anonymous" value="1"> Gửi ẩn danh</label>
            <?php if ($htUser['is_admin']): ?>
              <label class="ht-cmt__anon ht-cmt__anon--admin">
                <input type="checkbox" name="as_admin" value="1">
                Đăng với tư cách <strong>Ban quản trị</strong> — bài sẽ ghim lên đầu và có nhãn Quản trị
              </label>
            <?php endif; ?>
            <button type="submit" class="kame-button">Đăng bài</button>
          </form>
        <?php endif; ?>
      </div>

      <!-- Danh sách bài viết -->
      <?php if (!$htRoots): ?>
        <div class="kame-alert kame-alert--info" data-show="true">
          <span></span><span>Chưa có dữ liệu. Hãy là người mở đầu câu chuyện ở đây.</span><span></span>
        </div>
      <?php else: ?>
        <?php foreach ($htRoots as $root): ?>
          <?= ht_render_comment($root, $htChildren, $htLikes, $htUser) ?>
        <?php endforeach; ?>
      <?php endif; ?>

    <?php endif; ?>
  </div>
</section>

<script>
(function () {
  "use strict";

  /* --- Lọc FAQ theo từ khoá ------------------------------------------------
     Bỏ dấu tiếng Việt cả hai phía để gõ "quen mat khau" vẫn ra kết quả.     */
  var oTim     = document.getElementById("ht-o-tim");
  var oKetQua  = document.getElementById("ht-ket-qua");
  var oTrong   = document.getElementById("ht-khong-thay");
  var cacCau   = Array.prototype.slice.call(document.querySelectorAll("[data-ht-cau]"));
  var cacNhom  = Array.prototype.slice.call(document.querySelectorAll("[data-ht-nhom]"));

  function boDau(chuoi) {
    return (chuoi || "")
      .toLowerCase()
      .normalize("NFD")
      .replace(/[\u0300-\u036f]/g, "")
      .replace(/\u0111/g, "d");
  }

  // Tính sẵn chuỗi không dấu một lần, khỏi làm lại sau mỗi phím gõ.
  cacCau.forEach(function (cau) {
    cau.dataset.htKhongDau = boDau(cau.getAttribute("data-ht-tim"));
  });

  function loc() {
    if (!oTim) return;
    var tuKhoa = boDau(oTim.value.trim());
    var soKhop = 0;

    cacCau.forEach(function (cau) {
      var khop = tuKhoa === "" || cau.dataset.htKhongDau.indexOf(tuKhoa) !== -1;
      cau.hidden = !khop;
      if (khop) {
        soKhop++;
        // Đang tìm thì mở sẵn để đọc ngay, xoá ô tìm thì thu lại.
        cau.open = tuKhoa !== "";
      } else {
        cau.open = false;
      }
    });

    // Nhóm nào không còn câu nào khớp thì giấu cả tiêu đề nhóm đi.
    cacNhom.forEach(function (nhom) {
      var con = nhom.querySelectorAll("[data-ht-cau]:not([hidden])");
      nhom.hidden = con.length === 0;
    });

    if (oTrong) oTrong.setAttribute("data-hien", soKhop === 0 ? "true" : "false");

    if (oKetQua) {
      if (tuKhoa === "") oKetQua.textContent = "";
      else if (soKhop === 0) oKetQua.textContent = "Không tìm thấy câu hỏi nào.";
      else oKetQua.textContent = "Tìm thấy " + soKhop + " câu hỏi.";
    }
  }

  if (oTim) oTim.addEventListener("input", loc);

  document.querySelectorAll("[data-ht-goi-y]").forEach(function (nut) {
    nut.addEventListener("click", function () {
      if (!oTim) return;
      oTim.value = nut.getAttribute("data-ht-goi-y");
      loc();
      oTim.focus();
    });
  });

  /* --- Mở/đóng ô trả lời và ô sửa bài ------------------------------------ */
  document.addEventListener("click", function (e) {
    var nut = e.target.closest("[data-ht-toggle]");
    if (!nut) return;
    var o = document.getElementById(nut.getAttribute("data-ht-toggle"));
    if (!o) return;
    o.hidden = !o.hidden;
    if (!o.hidden) {
      var textarea = o.querySelector("textarea");
      if (textarea) textarea.focus();
    }
  });
})();
</script>

<?php
kame_auth_modals();
kame_footer();
