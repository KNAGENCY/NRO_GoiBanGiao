<?php
session_start();
require_once $_SERVER['DOCUMENT_ROOT'] . "/sv1/config.php";
require_once __DIR__ . "/admin_layout.php";

if (!isset($_SESSION['account_id'])) {
    header("Location: /sv1/trangchu");
    exit;
}
$stmt = $pdo->prepare("SELECT is_admin FROM account WHERE id = ? LIMIT 1");
$stmt->execute([(int)$_SESSION['account_id']]);
$acc = $stmt->fetch(PDO::FETCH_ASSOC);
if (!$acc || (int)$acc['is_admin'] !== 1) {
    http_response_code(403);
    exit('Không có quyền truy cập.');
}

/* --- Đọc số liệu giữ chân. Bảng có thể chưa tồn tại (chưa cắm bộ ghi) --- */
$ready = true;
$err = '';
try {
    $pdo->query("SELECT 1 FROM analytics_session LIMIT 1");
} catch (Throwable $e) {
    $ready = false;
    $err = 'Bảng analytics_session chưa có. Cần cắm bộ ghi vào server + tạo bảng trước khi có dữ liệu.';
}

$dau = [];            // người vào mỗi ngày (14 ngày gần nhất)
$totalSessions = 0;
$avgMin = 0;
$uniquePlayers = 0;
$d1 = null; $d7 = null;
$churn = [];          // rụng theo mốc sức mạnh

if ($ready) {
    // DAU 14 ngày
    $rows = $pdo->query("SELECT DATE(login_at) d, COUNT(DISTINCT account_id) n
        FROM analytics_session WHERE login_at >= DATE_SUB(CURDATE(), INTERVAL 13 DAY)
        GROUP BY DATE(login_at) ORDER BY d")->fetchAll(PDO::FETCH_ASSOC);
    foreach ($rows as $r) { $dau[$r['d']] = (int)$r['n']; }

    $totalSessions = (int)$pdo->query("SELECT COUNT(*) FROM analytics_session")->fetchColumn();
    $uniquePlayers = (int)$pdo->query("SELECT COUNT(DISTINCT account_id) FROM analytics_session")->fetchColumn();
    $avgSec = $pdo->query("SELECT AVG(duration_sec) FROM analytics_session WHERE duration_sec IS NOT NULL")->fetchColumn();
    $avgMin = $avgSec ? round($avgSec / 60, 1) : 0;

    // Giữ chân D1 / D7: trong số tài khoản tạo cách đây >=1 (>=7) ngày,
    // bao nhiêu % có đăng nhập vào ngày kế (trong 7 ngày kế).
    try {
        $d1row = $pdo->query("
          SELECT
            SUM(CASE WHEN s.account_id IS NOT NULL THEN 1 ELSE 0 END) giu,
            COUNT(*) tong
          FROM (SELECT id, DATE(create_time) cd FROM account
                WHERE create_time < CURDATE()) a
          LEFT JOIN (SELECT DISTINCT account_id, DATE(login_at) ld FROM analytics_session) s
            ON s.account_id = a.id AND s.ld = DATE_ADD(a.cd, INTERVAL 1 DAY)
        ")->fetch(PDO::FETCH_ASSOC);
        if ($d1row && (int)$d1row['tong'] > 0) $d1 = round(100 * $d1row['giu'] / $d1row['tong']);

        $d7row = $pdo->query("
          SELECT
            SUM(CASE WHEN s.account_id IS NOT NULL THEN 1 ELSE 0 END) giu,
            COUNT(*) tong
          FROM (SELECT id, DATE(create_time) cd FROM account
                WHERE create_time < DATE_SUB(CURDATE(), INTERVAL 6 DAY)) a
          LEFT JOIN (SELECT DISTINCT account_id, DATE(login_at) ld FROM analytics_session) s
            ON s.account_id = a.id AND s.ld BETWEEN DATE_ADD(a.cd, INTERVAL 1 DAY) AND DATE_ADD(a.cd, INTERVAL 7 DAY)
        ")->fetch(PDO::FETCH_ASSOC);
        if ($d7row && (int)$d7row['tong'] > 0) $d7 = round(100 * $d7row['giu'] / $d7row['tong']);
    } catch (Throwable $e) { /* giữ null */ }

    // Rụng theo sức mạnh: người >7 ngày không vào, xếp theo mốc power lần cuối
    try {
        $churn = $pdo->query("
          SELECT
            CASE
              WHEN mp < 1000 THEN 'Dưới 1N'
              WHEN mp < 100000 THEN '1N - 100N'
              WHEN mp < 10000000 THEN '100N - 10Tr'
              WHEN mp < 1000000000 THEN '10Tr - 1Tỷ'
              ELSE 'Trên 1 Tỷ'
            END mocSM, COUNT(*) n
          FROM (
            SELECT account_id, MAX(power) mp, MAX(login_at) ll
            FROM analytics_session GROUP BY account_id
            HAVING ll < DATE_SUB(NOW(), INTERVAL 7 DAY)
          ) t GROUP BY mocSM ORDER BY MIN(mp)
        ")->fetchAll(PDO::FETCH_ASSOC);
    } catch (Throwable $e) { $churn = []; }
}

$maxDau = $dau ? max($dau) : 1;
$maxChurn = $churn ? max(array_column($churn, 'n')) : 1;

admin_head('Thống Kê Giữ Chân');
admin_bar('Thống Kê Giữ Chân Người Chơi', 'chart-line');
?>
<div class="kadmin-body">
<?php if (!$ready): ?>
  <div class="kame-alert" style="margin-bottom:var(--kame-space-4)"><?= htmlspecialchars($err) ?></div>
<?php endif; ?>

  <div class="kadmin-stats">
    <div class="kadmin-stat"><p class="kadmin-stat__label">Giữ chân D1 (mai quay lại)</p><div class="kadmin-stat__value"><?= $d1 === null ? '—' : $d1 . '%' ?></div></div>
    <div class="kadmin-stat"><p class="kadmin-stat__label">Giữ chân D7 (trong 7 ngày)</p><div class="kadmin-stat__value"><?= $d7 === null ? '—' : $d7 . '%' ?></div></div>
    <div class="kadmin-stat"><p class="kadmin-stat__label">Thời lượng TB / phiên</p><div class="kadmin-stat__value"><?= $avgMin ?: '—' ?><?= $avgMin ? ' phút' : '' ?></div></div>
    <div class="kadmin-stat"><p class="kadmin-stat__label">Tổng phiên chơi</p><div class="kadmin-stat__value"><?= number_format($totalSessions) ?></div></div>
  </div>

  <div class="kadmin-card">
    <div class="kadmin-card__head"><i class="ti ti-calendar-stats" aria-hidden="true"></i> Người vào mỗi ngày (14 ngày gần nhất)</div>
    <?php if (!$dau): ?>
      <p class="kame-muted" style="margin:0">Chưa có dữ liệu. Sẽ có sau khi người chơi bắt đầu đăng nhập.</p>
    <?php else: ?>
      <div style="display:flex;align-items:flex-end;gap:6px;height:180px;padding-top:var(--kame-space-3)">
        <?php foreach ($dau as $day => $n): $h = max(6, round(150 * $n / $maxDau)); ?>
          <div style="flex:1;display:flex;flex-direction:column;align-items:center;gap:4px" title="<?= $day ?>: <?= $n ?> người">
            <span style="font-size:12px;font-weight:800;color:var(--kame-maroon-700)"><?= $n ?></span>
            <div style="width:100%;height:<?= $h ?>px;background:var(--kame-orange-600);border:1px solid var(--kame-maroon-700);border-radius:var(--kame-radius-sm) var(--kame-radius-sm) 0 0"></div>
            <span style="font-size:11px;color:var(--kame-ink)"><?= date('d/m', strtotime($day)) ?></span>
          </div>
        <?php endforeach; ?>
      </div>
    <?php endif; ?>
  </div>

  <div class="kadmin-card">
    <div class="kadmin-card__head"><i class="ti ti-user-off" aria-hidden="true"></i> Người bỏ game (7+ ngày không vào) theo sức mạnh</div>
    <?php if (!$churn): ?>
      <p class="kame-muted" style="margin:0">Chưa có ai rụng, hoặc chưa đủ dữ liệu.</p>
    <?php else: ?>
      <div class="kadmin-table-wrap">
        <table class="kadmin-table">
          <thead><tr><th>Mốc sức mạnh</th><th>Số người rụng</th><th>Biểu đồ</th></tr></thead>
          <tbody>
            <?php foreach ($churn as $c): $w = round(100 * $c['n'] / $maxChurn); ?>
              <tr>
                <td><?= htmlspecialchars($c['mocSM']) ?></td>
                <td class="num"><?= number_format($c['n']) ?></td>
                <td><div style="height:14px;width:<?= $w ?>%;min-width:8px;background:var(--kame-maroon-700);border-radius:var(--kame-radius-sm)"></div></td>
              </tr>
            <?php endforeach; ?>
          </tbody>
        </table>
      </div>
      <p class="kame-muted" style="margin:var(--kame-space-3) 0 0;font-size:13px">Mốc nào đông người rụng = chỗ người chơi chán, nên xem lại nội dung ở đó.</p>
    <?php endif; ?>
  </div>

  <p class="kame-muted" style="font-size:13px">Số người chơi từng vào: <?= number_format($uniquePlayers) ?>.</p>
</div>
<?php
admin_foot();
