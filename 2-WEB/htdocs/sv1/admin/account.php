<?php
session_start();
require_once $_SERVER['DOCUMENT_ROOT'] . "/sv1/config.php";
require_once __DIR__ . "/admin_layout.php";

if (empty($_SESSION['account_id'])) {
    header("Location: /sv1/trangchu.php");
    exit;
}

$account_id = (int)$_SESSION['account_id'];

$stmt = $pdo->prepare("SELECT id, is_admin FROM account WHERE id = ? LIMIT 1");
$stmt->execute([$account_id]);
$user = $stmt->fetch();

// ❌ Session fake hoặc user không tồn tại
if (!$user) {
    session_destroy();
    header("Location: /sv1/trangchu.php");
    exit;
}

// ❌ Không phải admin
if ((int)$user['is_admin'] !== 1) {
    echo "Không có quyền!";
    exit;
}

$csrfToken = ensure_csrf_token();

if ($_SERVER['REQUEST_METHOD'] === 'POST' && isset($_POST['action'])) {
    header('Content-Type: application/json');
    require_csrf_token();

    if ($_POST['action'] === 'delete') {
        $id = (int)$_POST['id'];
        try {
            $pdo->prepare("DELETE FROM player WHERE account_id = ?")->execute([$id]);
            $pdo->prepare("DELETE FROM account WHERE id = ?")->execute([$id]);
            echo json_encode(['success' => true]);
        } catch (Exception $e) {
            echo json_encode(['success' => false, 'error' => $e->getMessage()]);
        }
        exit;
    }

    if ($_POST['action'] === 'save') {
        $id       = (int)$_POST['id'];
        $username = trim($_POST['username']);
        $password = trim($_POST['password']);
        $ban      = (int)$_POST['ban'];
        $is_admin = (int)$_POST['is_admin'];
        $active   = (int)$_POST['active'];
        $cash     = (int)str_replace('.', '', $_POST['cash']);
        $danap    = (int)str_replace('.', '', $_POST['danap']);

        try {
            $check = $pdo->prepare("SELECT id FROM account WHERE username = ? AND id != ? LIMIT 1");
            $check->execute([$username, $id]);

            if ($check->fetch()) {
                echo json_encode([
                    'success' => false,
                    'error' => 'Tài khoản đã tồn tại!'
                ]);
                exit;
            }

            // Mật khẩu đã được mã hoá. Chỉ đổi khi admin gõ giá trị mới khác dấu che.
            if ($password === '' || $password === '••••••') {
                $sql = "UPDATE account
                        SET username=?, ban=?, is_admin=?, active=?, cash=?, danap=?
                        WHERE id=?";
                $pdo->prepare($sql)->execute([
                    $username, $ban, $is_admin, $active, $cash, $danap, $id
                ]);
            } else {
                $sql = "UPDATE account
                        SET username=?, password=?, ban=?, is_admin=?, active=?, cash=?, danap=?
                        WHERE id=?";
                $pdo->prepare($sql)->execute([
                    $username, kame_password_hash($password), $ban, $is_admin, $active, $cash, $danap, $id
                ]);
            }

            echo json_encode(['success' => true]);

        } catch (Exception $e) {
            echo json_encode([
                'success' => false,
                'error' => $e->getMessage()
            ]);
        }
        exit;
    }

    if ($_POST['action'] === 'quick_create') {
        $tk = 'testaccount';
        try {
            $check = $pdo->prepare("SELECT id FROM account WHERE username = ?");
            $check->execute([$tk]);
            if ($check->fetch()) {
                echo json_encode(['success' => false, 'error' => 'Tài khoản đã tồn tại!']);
            } else {
                $pdo->prepare("INSERT INTO account (username, password, ban, is_admin, active, cash, danap) VALUES (?,?,0,0,0,0,0)")
                    ->execute([$tk, kame_password_hash($tk)]);
                echo json_encode(['success' => true]);
            }
        } catch (Exception $e) {
            echo json_encode(['success' => false, 'error' => $e->getMessage()]);
        }
        exit;
    }
    exit;
}

$page     = max(1, (int)($_GET['page'] ?? 1));
$perPage  = 100;
$offset   = ($page - 1) * $perPage;

$filters = [
    'id'         => trim($_GET['f_id']         ?? ''),
    'name'       => trim($_GET['f_name']       ?? ''),
    'username'   => trim($_GET['f_username']   ?? ''),
    'password'   => trim($_GET['f_password']   ?? ''),
    'ban'        => trim($_GET['f_ban']        ?? ''),
    'is_admin'   => trim($_GET['f_is_admin']   ?? ''),
    'ip_address' => trim($_GET['f_ip_address'] ?? ''),
    'active'     => trim($_GET['f_active']     ?? ''),
    'cash'       => str_replace('.', '', trim($_GET['f_cash'] ?? '')),
    'danap'      => str_replace('.', '', trim($_GET['f_danap'] ?? '')),
];

$where  = [];
$params = [];

foreach ($filters as $col => $val) {
    if ($col === 'password') { continue; } // mat khau da bam, khong loc duoc
    if ($val !== '') {
        if ($col === 'name') {
            $where[]  = "p.name LIKE ?";
            $params[] = "%$val%";
            continue;
        }
        if (in_array($col, ['id','ban','is_admin','active'])) {
            $where[]  = "a.$col = ?";
            $params[] = $val;
        } elseif (in_array($col, ['cash','danap'])) {
            $where[]  = "CAST(a.$col AS CHAR) LIKE ?";
            $params[] = "%$val%";
        } else {
            $where[]  = "a.$col LIKE ?";
            $params[] = "%$val%";
        }
    }
}
$whereSql = $where ? 'WHERE ' . implode(' AND ', $where) : '';

$countSql = "SELECT COUNT(*) FROM account a LEFT JOIN player p ON p.account_id = a.id $whereSql";
$countStmt = $pdo->prepare($countSql);
$countStmt->execute($params);
$total      = (int)$countStmt->fetchColumn();
$totalPages = max(1, ceil($total / $perPage));

$sql = "SELECT a.id, a.username, a.password, a.ban, a.is_admin, a.ip_address, a.active, a.cash, a.danap,
               p.name as char_name
        FROM account a
        LEFT JOIN player p ON p.account_id = a.id
        $whereSql
        ORDER BY a.danap DESC, a.id ASC
        LIMIT $perPage OFFSET $offset";

$stmt = $pdo->prepare($sql);
$stmt->execute($params);
$accounts = $stmt->fetchAll();

$stBan   = $pdo->query("SELECT COUNT(*) FROM account WHERE ban=1")->fetchColumn();
$stAdmin = $pdo->query("SELECT COUNT(*) FROM account WHERE is_admin=1")->fetchColumn();
$stUser  = $pdo->query("SELECT COUNT(*) FROM account WHERE is_admin=0")->fetchColumn();
$stAct   = $pdo->query("SELECT COUNT(*) FROM account WHERE active=1")->fetchColumn();
$stTotal = $pdo->query("SELECT COUNT(*) FROM account")->fetchColumn();

admin_head('Quản Lý Tài Khoản');
admin_bar('Quản Lý Tài Khoản', 'user-cog');
?>
<style>
/* ==========================================================================
   Reskin Kame cho trang Quản Lý Tài Khoản.
   Giữ NGUYÊN mọi tên class/id mà JS + form tra cứu, chỉ viết lại định dạng
   bằng biến --kame-*. Alias các token cũ (--danger, --success, --blue2...)
   sang biến Kame để inline-style trong markup và JS vẫn hiển thị đúng hệ màu.
   ========================================================================== */
#scale-root {
  --bg:      var(--kame-paper-100);
  --surf:    var(--kame-white);
  --surf2:   var(--kame-paper-100);
  --surf3:   var(--kame-paper-300);
  --bd:      var(--kame-paper-300);
  --bd2:     var(--kame-maroon-700);
  --tx:      var(--kame-ink);
  --tx2:     var(--kame-maroon-900);
  --tx3:     var(--kame-maroon-700);
  --blue:    var(--kame-blue);
  --blue2:   var(--kame-blue);
  --indigo:  var(--kame-maroon-900);
  --sky:     var(--kame-blue);
  --teal:    var(--kame-green);
  --success: var(--kame-green);
  --danger:  var(--kame-red);
  --warn:    var(--kame-orange-700);
  --blu-lt:  var(--kame-paper-300);
  --ind-lt:  var(--kame-paper-300);
  --sky-lt:  var(--kame-paper-300);
  --glow:    rgb(96 0 24 / 16%);
  --r:  var(--kame-radius-md);
  --rs: var(--kame-radius-sm);
  --sh:  var(--kame-shadow-md);
  --sh2: var(--kame-shadow-sm);
}

#scale-root {
  width: 1500px;
  transform-origin: top left;
  font-size: 14px;
  line-height: 1.6;
  color: var(--tx);
  position: relative;
}

.W {
  position: relative; z-index: 1;
  width: 1500px;
  padding: 4px 4px 24px;
}

/* TOPBAR */
.top {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  gap: 8px;
  background: #fffdf7; border: var(--kame-border);
  border-left: 5px solid var(--kame-orange-600);
  border-radius: var(--r); padding: 16px 22px 14px;
  margin-bottom: 10px; box-shadow: var(--sh);
  position: relative; overflow: hidden; text-align: center;
}
.ttit { font-size: 1.35rem; font-weight: 800; letter-spacing: -.4px; color: var(--tx); }

.tx-box {
  display: inline-flex; align-items: center; gap: 10px;
  padding: 11px 26px;
  border-radius: 999px;
  background: var(--kame-maroon-900);
  color: var(--kame-paper-100); font-weight: 800; font-size: 1.05rem;
  box-shadow: var(--sh2);
  border: 1px solid var(--kame-orange-600);
  letter-spacing: .3px; cursor: pointer;
  transition: all .15s;
  line-height: 1.3;
}
.tx-box:hover { filter: brightness(1.07); transform: translateY(-1px); }

.tr2 { display: flex; gap: 8px; flex-wrap: wrap; align-items: center; justify-content: center; }
.btn {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 11px 20px; border: none; border-radius: var(--rs);
  font-family: inherit; font-size: .9rem; font-weight: 700;
  cursor: pointer; text-decoration: none; white-space: nowrap;
  transition: all .15s; letter-spacing: .2px;
  line-height: 1;
}
.btn:hover { transform: translateY(-2px); box-shadow: var(--sh2); }
.btn:active { transform: scale(.96); }
.bblue  { background: var(--kame-blue); color: var(--kame-white); }
.bteal  { background: var(--kame-green); color: var(--kame-white); }

/* STATS */
.stats-row {
  display: flex; gap: 6px; flex-wrap: nowrap;
  justify-content: center; margin-bottom: 9px;
}
.stat-card {
  background: #fffdf7; border: var(--kame-border);
  border-left: 5px solid var(--kame-orange-600);
  border-radius: var(--r); padding: 9px 20px;
  display: flex; align-items: center; gap: 8px;
  box-shadow: var(--sh2); min-width: 100px; justify-content: center;
  text-align: center;
}
.stat-val { font-size: 1.1rem; font-weight: 900; font-family: "Roboto Condensed", "Arial Narrow", sans-serif; line-height: 1; text-align: center; }
.stat-lbl { font-size: .72rem; color: var(--tx3); font-weight: 700; line-height: 1; text-align: center; }

/* LEGEND */
.legend-row {
  display: flex; flex-wrap: nowrap; gap: 8px; align-items: center;
  justify-content: center; margin-bottom: 10px;
}
.legend-item {
  display: inline-flex; align-items: center; gap: 5px;
  font-size: .78rem; font-weight: 700; color: var(--tx2);
  background: #fffdf7; border: var(--kame-border);
  border-radius: 6px; padding: 4px 10px; box-shadow: var(--sh2);
}
.legend-dot { width: 12px; height: 12px; border-radius: 3px; border: 1px solid rgb(96 0 24 / 30%); flex-shrink: 0; }

/* PAGINATION */
.pag-wrap {
  display: flex; gap: 5px; flex-wrap: wrap; justify-content: center;
  margin-bottom: 10px; align-items: center;
}
.page-link {
  padding: 5px 11px; border-radius: var(--rs);
  text-decoration: none; font-weight: 700; font-size: .82rem;
  border: 1.5px solid var(--bd2);
  background: #fffdf7; color: var(--tx2); transition: all .15s;
}
.page-link:hover { background: var(--kame-paper-300); border-color: var(--kame-orange-600); }
.page-link.active { background: var(--kame-maroon-900); color: var(--kame-paper-100); border-color: var(--kame-orange-600); pointer-events: none; }
.page-info { font-size: .82rem; color: var(--tx3); font-weight: 600; }

/* TABLE */
.tc {
  background: #fffdf7; border: var(--kame-border);
  border-radius: var(--r); overflow: hidden; box-shadow: var(--sh);
}
.tsc {
  overflow: visible;
  overflow-x: visible;
  overflow-y: visible;
}

col.c-id   { width: 5%; }
col.c-name { width: 11%; }
col.c-user { width: 11%; }
col.c-pass { width: 11%; }
col.c-ban  { width: 8.2%; }
col.c-adm  { width: 8.2%; }
col.c-ip   { width: 11%; }
col.c-act  { width: 8.2%; }
col.c-cash { width: 10.2%; }
col.c-dnap { width: 10.2%; }
col.c-hd   { width: 6%; }

table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}

thead tr.header-row th {
  background: var(--kame-paper-300);
  color: var(--kame-maroon-900);
  padding: 11px 6px 9px;
  font-size: .79rem; font-weight: 800;
  text-transform: uppercase; letter-spacing: .7px;
  border-bottom: 2px solid var(--kame-maroon-700);
  white-space: nowrap; text-align: center;
  overflow: hidden; text-overflow: ellipsis;
}

thead tr.filter-row th {
  background: var(--surf3);
  padding: 5px 5px;
  border-bottom: 2px solid var(--bd2);
}
thead tr.filter-row input {
  width: 100%; padding: 5px 7px;
  border-radius: var(--rs);
  border: 1.5px solid var(--bd2);
  font-size: .78rem; font-weight: 600;
  background: var(--surf); color: var(--tx);
  font-family: inherit;
  outline: none; transition: border .15s;
  text-align: center;
}
thead tr.filter-row input:focus {
  border-color: var(--kame-orange-600);
  box-shadow: 0 0 0 3px rgb(214 106 0 / 18%);
}
thead tr.filter-row input::placeholder { color: var(--kame-maroon-700); opacity: .55; }

tbody tr { transition: filter .12s; }
tbody tr:hover td { filter: brightness(.96); }
tbody tr:nth-child(odd)  td { background: #fffdf7; }
tbody tr:nth-child(even) td { background: #fdf7e8; }

tbody tr.row-ban    td { background: #f6d6d0 !important; }
tbody tr.row-admin  td { background: #fbe6bf !important; }
tbody tr.row-active1 td { background: #d8ecd2 !important; }
tbody tr.row-active0 td { background: #ece5d6 !important; }

td {
  padding: 7px 6px;
  text-align: center;
  border-bottom: 1px solid var(--bd);
  vertical-align: middle;
  overflow: hidden;
}

.view-value {
  display: inline-block;
  width: 100%;
  font-size: .82rem;
  font-weight: 600;
  font-family: inherit;
  background: transparent;
  padding: 4px 6px;
  border-radius: var(--rs);
  color: var(--tx2);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.edit-input {
  width: 100%;
  padding: 4px 6px;
  border: 1.5px solid var(--bd2);
  border-radius: var(--rs);
  font-size: .8rem; font-weight: 600;
  background: rgb(255 253 245 / 85%);
  color: var(--tx); text-align: center;
  font-family: inherit;
  outline: none;
  transition: border .15s, background .15s, box-shadow .15s;
}
.edit-input:focus {
  border-color: var(--kame-orange-600);
  background: var(--surf);
  box-shadow: 0 0 0 3px rgb(214 106 0 / 18%);
}
select.edit-input { cursor: pointer; }

.action-wrap {
  display: flex; gap: 4px; justify-content: center;
  flex-direction: column; align-items: stretch;
}
.btn-sm {
  display: inline-flex; align-items: center; justify-content: center; gap: 4px;
  padding: 5px 8px; font-size: .74rem; font-weight: 700;
  border-radius: var(--rs); border: 1.5px solid; cursor: pointer;
  transition: all .14s; font-family: inherit; white-space: nowrap;
}
.btn-edit    { background: var(--kame-paper-100); color: var(--kame-blue); border-color: var(--kame-maroon-700); }
.btn-edit:hover  { background: var(--kame-paper-300); }
.btn-del-p   { background: var(--kame-paper-100); color: var(--kame-red); border-color: var(--kame-maroon-700); }
.btn-del-p:hover   { background: var(--kame-paper-300); }
.btn-save-p  { background: var(--kame-paper-100); color: var(--kame-green); border-color: var(--kame-maroon-700); }
.btn-save-p:hover  { background: var(--kame-paper-300); }
.btn-cancel-p {
  background: var(--kame-paper-100); color: var(--kame-red); border-color: var(--kame-maroon-700);
}
.btn-cancel-p:hover { background: var(--kame-paper-300); }

/* TOAST */
#toast {
  position: fixed;
  bottom: 32px;
  right: 32px;
  z-index: 9999;
  font-size: 1.5rem;
  padding: 21px 33px;
  border-radius: var(--kame-radius-md);
  background: var(--kame-maroon-900);
  color: var(--kame-white);
  font-weight: 800;
  box-shadow: var(--sh);
  border: 2px solid var(--kame-orange-600);
  display: none;
  animation: pop .22s ease;
  max-width: 480px;
  font-family: inherit;
}
#toast.error {
  background: var(--kame-red);
}
@keyframes pop {
  from { transform: scale(.88); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}

#loading {
  position: fixed; inset: 0;
  background: rgb(40 21 15 / 22%);
  display: none; align-items: center; justify-content: center; z-index: 9998;
}
#loading .spinner {
  width: 44px; height: 44px;
  border: 4px solid var(--kame-paper-300);
  border-top-color: var(--kame-orange-600);
  border-radius: 50%;
  animation: spin .7s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* ========== CONFIRM MODAL ========== */
.confirm-modal {
  display: none; position: fixed; z-index: 10000;
  left: 0; top: 0; width: 100%; height: 100%;
  background: rgb(40 21 15 / 65%);
  align-items: center; justify-content: center;
}
.confirm-box {
  position: relative;
  background: var(--surf); border: var(--kame-border);
  border-top: 4px solid var(--kame-orange-600);
  border-radius: var(--r); padding: 2rem;
  max-width: 550px; width: 90%; text-align: center;
  animation: pop .2s ease; box-shadow: var(--sh);
  font-size: 1.2rem;
}
/* Dấu X có khung tròn nổi bật */
.close-x {
  position: absolute;
  top: 12px;
  right: 18px;
  width: 32px;
  height: 32px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: bold;
  cursor: pointer;
  color: var(--tx3);
  background: rgb(96 0 24 / 8%);
  border-radius: 50%;
  transition: all 0.2s;
  line-height: 1;
}
.close-x:hover {
  background: var(--kame-red);
  color: var(--kame-white);
  transform: scale(1.05);
}
.confirm-box h3 {
  color: var(--kame-orange-700);
  font-size: 1.4rem; margin-bottom: 16px;
}
.confirm-msg-main {
  color: var(--kame-green);
  font-weight: 700;
  font-size: 1.5rem;
}
.confirm-msg-warning {
  color: var(--kame-red);
  font-weight: 600;
  display: block;
  margin-top: 8px;
  font-size: 1.5rem;
}
.confirm-box p {
  color: var(--tx3); margin-bottom: 24px;
  font-size: 1em;
  word-break: break-word;
  white-space: normal;
  line-height: 1.4;
}
.confirm-btns {
  display: flex; gap: 20px; justify-content: center;
}
.confirm-btns .btn-sm {
  font-size: 1rem; padding: 8px 20px; width: auto;
}
.btn-cancel-con {
  background: var(--kame-paper-100); color: var(--kame-red); border: 1.5px solid var(--kame-maroon-700);
}
.btn-cancel-con:hover { background: var(--kame-paper-300); }
.btn-ok-con {
  background: var(--kame-red); color: var(--kame-white); border: none;
}
.btn-ok-con:hover { filter: brightness(1.08); transform: translateY(-2px); }

.btn.account-guide-btn {
  color: var(--kame-paper-100);
  background: var(--kame-maroon-900);
  border: 2px solid var(--kame-orange-600);
  box-shadow: var(--sh2);
}
.btn.account-guide-btn:hover {
  filter: brightness(1.08);
  border-color: var(--kame-orange-400);
}

#scale-root ::-webkit-scrollbar { width: 5px; height: 5px; }
#scale-root ::-webkit-scrollbar-track { background: var(--surf2); }
#scale-root ::-webkit-scrollbar-thumb { background: var(--bd2); border-radius: 3px; }
#scale-root ::-webkit-scrollbar-thumb:hover { background: var(--tx3); }
</style>

<div class="kadmin-body">

<div id="loading"><div class="spinner"></div></div>
<div id="toast"></div>

<div id="scale-root">
<div class="W">

<!-- TOPBAR -->
<div class="top">
  <div class="ttit">
    <button class="tx-box" onclick="quickCreate()" style="border:none;">
      <svg viewBox="0 0 24 24" width="16" height="16" fill="none">
        <circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="1.8"/>
        <circle cx="17" cy="9" r="3" stroke="currentColor" stroke-width="1.6"/>
        <path d="M2 21c0-4 3.1-7 7-7s7 3 7 7" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
        <path d="M22 21c0-2.8-2-5-4.5-5.5" stroke="currentColor" stroke-width="1.6" stroke-linecap="round"/>
      </svg>
      Tạo Tài Khoản Nhanh:
      <span style="font-family:'Roboto Condensed', monospace;font-size:.95em"><?= number_format($stTotal) ?></span>
    </button>
  </div>

  <div class="tr2">
    <a href="/sv1/admin/panel.php" class="btn bblue">
      <svg width="15" height="15" viewBox="0 0 24 24" fill="none">
        <path d="M3 10.5L12 3L21 10.5V20A1 1 0 0 1 20 21H15V14H9V21H4A1 1 0 0 1 3 20V10.5Z"
              stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
      </svg>
      <span>Panel</span>
    </a>
    <button class="btn bteal" onclick="location.reload()">
      <svg width="15" height="15" viewBox="0 0 24 24" fill="none">
        <path d="M20 12A8 8 0 1 1 17.3 6M20 4V12H12"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      <span>Làm Mới</span>
    </button>
  </div>
</div>

<!-- STATS -->
<div class="stats-row">
  <div class="stat-card">
    <div>
      <div class="stat-val" style="color:var(--success)"><?= number_format($stAct) ?></div>
      <div class="stat-lbl">MTV</div>
    </div>
  </div>
  <div class="stat-card">
    <div>
      <div class="stat-val" style="color:var(--danger)"><?= number_format($stBan) ?></div>
      <div class="stat-lbl">Khóa</div>
    </div>
  </div>
  <div class="stat-card">
    <div>
      <div class="stat-val" style="color:var(--warn)"><?= number_format($stAdmin) ?></div>
      <div class="stat-lbl">Admin</div>
    </div>
  </div>
  <div class="stat-card">
    <div>
      <div class="stat-val" style="color:var(--tx3)"><?= number_format($stUser) ?></div>
      <div class="stat-lbl">User</div>
    </div>
  </div>
</div>

<!-- LEGEND -->
<div class="legend-row">
  <div class="legend-item"><div class="legend-dot" style="background:#d8ecd2"></div>Đã MTV</div>
  <div class="legend-item"><div class="legend-dot" style="background:#f6d6d0"></div>Bị khóa</div>
  <div class="legend-item"><div class="legend-dot" style="background:#fbe6bf"></div>Admin</div>
  <div class="legend-item"><div class="legend-dot" style="background:#ece5d6"></div>User</div>
</div>

<!-- PAGINATION -->
<?php
$qs = $_GET;
$start_loop = max(1, $page - 4);
$end_loop   = min($totalPages, $page + 5);
?>
<div class="pag-wrap">
  <?php for ($i = $start_loop; $i <= $end_loop; $i++): $qs['page'] = $i; ?>
    <a href="?<?= http_build_query($qs) ?>" class="page-link <?= ($i == $page) ? 'active' : '' ?>"><?= $i ?></a>
  <?php endfor; ?>
</div>

<!-- FILTER FORM -->
<form id="filterForm" method="GET">
  <input type="hidden" name="page" value="1">
</form>

<!-- TABLE -->
<div class="tc"><div class="tsc">
<table>
<colgroup>
  <col class="c-id">
  <col class="c-name">
  <col class="c-user">
  <col class="c-pass">
  <col class="c-ban">
  <col class="c-adm">
  <col class="c-ip">
  <col class="c-act">
  <col class="c-cash">
  <col class="c-dnap">
  <col class="c-hd">
</colgroup>
<thead>
  <tr class="header-row">
    <th>ID</th>
    <th>Name</th>
    <th>Tài Khoản</th>
    <th>Pass</th>
    <th>Khóa</th>
    <th>Admin</th>
    <th>IP</th>
    <th>Thành Viên</th>
    <th>Cá Tươi</th>
    <th>Tổng Cá</th>
    <th>VIP</th>
  </tr>
  <tr class="filter-row">
    <th><input type="text" name="f_id"         form="filterForm" value="<?= htmlspecialchars($filters['id']) ?>"         placeholder="ID"        oninput="submitFilter()"></th>
    <th><input type="text" name="f_name"       form="filterForm" value="<?= htmlspecialchars($filters['name']) ?>"       placeholder="Tên NV"    oninput="submitFilter()"></th>
    <th><input type="text" name="f_username"   form="filterForm" value="<?= htmlspecialchars($filters['username']) ?>"   placeholder="Tài khoản" oninput="submitFilter()"></th>
    <th><input type="text" name="f_password"   form="filterForm" value="<?= htmlspecialchars($filters['password']) ?>"   placeholder="(đã mã hoá)" disabled title="Mật khẩu đã được mã hoá, không thể tìm kiếm"></th>
    <th><input type="text" name="f_ban"        form="filterForm" value="<?= htmlspecialchars($filters['ban']) ?>"        placeholder="0/1"       oninput="submitFilter()"></th>
    <th><input type="text" name="f_is_admin"   form="filterForm" value="<?= htmlspecialchars($filters['is_admin']) ?>"   placeholder="0/1"       oninput="submitFilter()"></th>
    <th><input type="text" name="f_ip_address" form="filterForm" value="<?= htmlspecialchars($filters['ip_address']) ?>" placeholder="IP"        oninput="submitFilter()"></th>
    <th><input type="text" name="f_active"     form="filterForm" value="<?= htmlspecialchars($filters['active']) ?>"     placeholder="0/1"       oninput="submitFilter()"></th>
    <th><input type="text" name="f_cash"       form="filterForm" value="<?= htmlspecialchars($filters['cash']) ?>"       placeholder="Số dư"     oninput="submitFilter()"></th>
    <th><input type="text" name="f_danap"      form="filterForm" value="<?= htmlspecialchars($filters['danap']) ?>"      placeholder="Tổng nạp"  oninput="submitFilter()"></th>
    <th></th>
  </tr>
</thead>
<tbody>
<?php foreach ($accounts as $row):
  $isBan   = $row['ban']      == 1;
  $isActive= $row['active']   == 1;
  $isAdmin = $row['is_admin'] == 1;
  $rowClass = $isBan ? 'row-ban' : ($isAdmin ? 'row-admin' : ($isActive ? 'row-active1' : 'row-active0'));
?>
<tr class="<?= $rowClass ?>" data-id="<?= $row['id'] ?>">
  <td><span class="no-edit view-value"><?= $row['id'] ?></span></td>
  <td><strong style="color:var(--blue2);font-size:.83rem;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;display:block"><?= htmlspecialchars($row['char_name'] ?? '—') ?></strong></td>
  <td><span class="view-value" data-field="username" data-original="<?= htmlspecialchars($row['username']) ?>"><?= htmlspecialchars($row['username']) ?></span></td>
  <td><span class="view-value" data-field="password" data-original="••••••" title="Mật khẩu đã mã hoá. Gõ mật khẩu mới vào đây để đặt lại.">••••••</span></td>
  <td><span class="view-value" data-field="ban" data-original="<?= $row['ban'] ?>"><?= $row['ban'] == 1 ? '1 - Khóa' : '0 - Mở' ?></span></td>
  <td><span class="view-value" data-field="is_admin" data-original="<?= $row['is_admin'] ?>"><?= $row['is_admin'] == 1 ? '1 - Admin' : '0 - User' ?></span></td>
  <td><span class="no-edit view-value"><?= htmlspecialchars($row['ip_address'] ?? '—') ?></span></td>
  <td><span class="view-value" data-field="active" data-original="<?= $row['active'] ?>"><?= $row['active'] == 1 ? '1 - Mở' : '0 - Đóng' ?></span></td>
  <td><span class="view-value" data-field="cash" data-original="<?= $row['cash'] ?>"><?= number_format($row['cash'], 0, ',', '.') ?></span></td>
  <td><span class="view-value" data-field="danap" data-original="<?= $row['danap'] ?>"><?= number_format($row['danap'], 0, ',', '.') ?></span></td>
  <td>
    <div class="action-wrap">
      <button class="btn-sm btn-edit" onclick="editRow(this)">✏️ Sửa</button>
      <button class="btn-sm btn-del-p" onclick="deleteRow(this, <?= $row['id'] ?>, '<?= htmlspecialchars($row['username'], ENT_QUOTES) ?>')">🗑️ Xóa</button>
    </div>
  </td>
</tr>
<?php endforeach; ?>
<?php if (empty($accounts)): ?>
<tr><td colspan="11" style="padding:32px;color:var(--tx3);font-style:italic">Không tìm thấy tài khoản nào.</td></tr>
<?php endif; ?>
</tbody>
</table>
</div></div>

</div><!-- /W -->
</div><!-- /scale-root -->

<!-- CONFIRM MODAL -->
<div id="confirmModal" class="confirm-modal">
  <div class="confirm-box">
    <span class="close-x">&times;</span>
    <h3 id="confirmTitle">⚠️ Xác Nhận</h3>
    <p id="confirmMsg">Bạn có chắc chắn?</p>
    <div class="confirm-btns">
      <button class="btn-sm btn-cancel-con" id="confirmCancelBtn">✕ Hủy</button>
      <button class="btn-sm btn-ok-con" id="confirmOkBtn">OK</button>
    </div>
  </div>
</div>

<script>
const CSRF_TOKEN = <?= json_encode($csrfToken) ?>;
function csrfHeaders(extra = {}) {
    return Object.assign({}, extra, { 'X-CSRF-TOKEN': CSRF_TOKEN });
}

/* ================= SCALE ================= */
const DESIGN_WIDTH = 1500;
const root = document.getElementById('scale-root');
function applyScale() {
    const vw = window.innerWidth;
    const ratio = (vw * 0.98) / DESIGN_WIDTH;
    root.style.transform = `scale(${ratio})`;
    root.style.transformOrigin = 'top left';
    root.style.marginLeft = '1vw';
    document.body.style.height = (root.offsetHeight * ratio) + 'px';
    document.body.style.width  = '100%';
}
applyScale();
window.addEventListener('resize', applyScale);

/* ================= TOAST / LOADING ================= */
function showToast(msg, type='success') {
    const t = document.getElementById('toast');
    t.textContent = msg; t.className = type; t.style.display = 'block';
    setTimeout(() => t.style.display = 'none', 2000);
}
function showLoading(v) { document.getElementById('loading').style.display = v ? 'flex' : 'none'; }

/* ================= CONFIRM MODAL ================= */
let confirmCallback = null;
function showConfirm({title, msgMain, msgWarning, okText='OK'}, callback) {
    const modal = document.getElementById('confirmModal');
    const titleEl = document.getElementById('confirmTitle');
    const msgEl = document.getElementById('confirmMsg');
    titleEl.textContent = title || '⚠️ Xác Nhận';
    msgEl.innerHTML = `
        <span class="confirm-msg-main">${escapeHtml(msgMain)}</span>
        ${msgWarning ? `<span class="confirm-msg-warning">${escapeHtml(msgWarning)}</span>` : ''}
    `;
    const okBtn = document.getElementById('confirmOkBtn');
    okBtn.textContent = okText || 'OK';
    confirmCallback = callback;
    modal.style.display = 'flex';
}
function closeConfirmModal() {
    document.getElementById('confirmModal').style.display = 'none';
    confirmCallback = null;
}
document.getElementById('confirmOkBtn').onclick = () => {
    if (confirmCallback) confirmCallback();
    closeConfirmModal();
};
document.getElementById('confirmCancelBtn').onclick = closeConfirmModal;
document.querySelector('.close-x').onclick = closeConfirmModal;
document.addEventListener('keydown', e => { if (e.key === 'Escape') closeConfirmModal(); });

function escapeHtml(str) {
    if(!str) return '';
    return str.replace(/[&<>]/g, function(m) {
        if(m === '&') return '&amp;';
        if(m === '<') return '&lt;';
        if(m === '>') return '&gt;';
        return m;
    });
}

/* ================= DELETE ROW ================= */
async function deleteRow(btn, id, username) {
    showConfirm({
        title: '🗑️ Xóa tài khoản',
        msgMain: `Xóa TK "${escapeHtml(username)}" — ID ${id}.`,
        msgWarning: 'Không thể hoàn tác!',
        okText: 'Xóa'
    }, async () => {
        showLoading(true);
        try {
            const res  = await fetch('account.php', { method:'POST', headers:csrfHeaders({'Content-Type':'application/x-www-form-urlencoded'}), body: new URLSearchParams({action:'delete', id}) });
            const json = await res.json();
            if (json.success) {
                btn.closest('tr').remove();
                showToast('🗑️ Đã xóa Tài Khoản: ' + username);
            }
            else showToast('❌ ' + json.error, 'error');
        } catch(e) { showToast('❌ Lỗi kết nối!', 'error'); }
        showLoading(false);
    });
}

/* ================= QUICK CREATE ================= */
function quickCreate() {
    showConfirm({
        title: '⚡ Tạo Tài Khoản Nhanh',
        msgMain: 'Ấn nút OK là được',
        msgWarning: '',
        okText: 'OK'
    }, async () => {
        showLoading(true);
        try {
            const res  = await fetch('account.php', { method:'POST', headers:csrfHeaders({'Content-Type':'application/x-www-form-urlencoded'}), body: new URLSearchParams({action:'quick_create'}) });
            const json = await res.json();
            if (json.success) { showToast('✅ Tạo thành công!'); setTimeout(() => location.reload(), 1200); }
            else showToast('❌ ' + json.error, 'error');
        } catch(e) { showToast('❌ Lỗi kết nối!', 'error'); }
        showLoading(false);
    });
}

/* ================= EDIT MODE LOGIC ================= */
function editRow(btn) {
    const tr = btn.closest('tr');
    if (tr.classList.contains('editing')) return;

    tr.classList.add('editing');

    const originalValues = {};
    const editableSpans = tr.querySelectorAll('td .view-value[data-field]');
    editableSpans.forEach(span => {
        const field = span.getAttribute('data-field');
        originalValues[field] = span.getAttribute('data-original') || span.innerText;
    });

    const delBtn = tr.querySelector('.btn-del-p');
    if (delBtn) delBtn.style.display = 'none';

    editableSpans.forEach(span => {
        const field = span.getAttribute('data-field');
        const currentValue = span.getAttribute('data-original') || span.innerText;
        const td = span.parentNode;
        let inputEl = null;

        if (field === 'username' || field === 'password') {
            inputEl = document.createElement('input');
            inputEl.type = 'text';
            inputEl.value = currentValue;
            inputEl.classList.add('edit-input');
        } else if (field === 'ban' || field === 'is_admin' || field === 'active') {
            inputEl = document.createElement('select');
            inputEl.classList.add('edit-input');
            if (field === 'ban') {
                inputEl.innerHTML = '<option value="0">0 - Mở</option><option value="1">1 - Khóa</option>';
                inputEl.value = currentValue;
            } else if (field === 'is_admin') {
                inputEl.innerHTML = '<option value="0">0 - User</option><option value="1">1 - Admin</option>';
                inputEl.value = currentValue;
            } else if (field === 'active') {
                inputEl.innerHTML = '<option value="0">0 - Đóng</option><option value="1">1 - Mở</option>';
                inputEl.value = currentValue;
            }
            inputEl.style.color = inputEl.value == '1' ? 'var(--danger)' : 'var(--success)';
            inputEl.addEventListener('change', function() {
                this.style.color = this.value == '1' ? 'var(--danger)' : 'var(--success)';
            });
        } else if (field === 'cash' || field === 'danap') {
            inputEl = document.createElement('input');
            inputEl.type = 'text';
            let raw = currentValue;
            if (raw.indexOf('.') !== -1) raw = raw.replace(/\./g, '');
            inputEl.value = Number(raw).toLocaleString('vi-VN');
            inputEl.classList.add('edit-input');
            inputEl.addEventListener('input', function() {
                let v = this.value.replace(/\D/g, '');
                this.value = v ? Number(v).toLocaleString('vi-VN') : '';
            });
        }

        if (inputEl) {
            inputEl.setAttribute('data-field', field);
            span.style.display = 'none';
            td.appendChild(inputEl);
            inputEl.focus();
        }
    });

    const actionDiv = tr.querySelector('.action-wrap');
    const editBtn = actionDiv.querySelector('.btn-edit');
    const saveBtn = document.createElement('button');
    saveBtn.className = 'btn-sm btn-save-p';
    saveBtn.innerHTML = '💾 Lưu';
    const cancelBtn = document.createElement('button');
    cancelBtn.className = 'btn-sm btn-cancel-p';
    cancelBtn.innerHTML = '✖️ Hủy';

    editBtn.remove();
    actionDiv.prepend(saveBtn);
    actionDiv.prepend(cancelBtn);

    saveBtn.onclick = () => saveEdit(tr, originalValues);
    cancelBtn.onclick = () => cancelEdit(tr, originalValues);
}

async function saveEdit(tr, originalValues) {
    const id = tr.getAttribute('data-id');
    const inputs = tr.querySelectorAll('.edit-input');
    const data = {
        action: 'save',
        id: id,
        username: '',
        password: '',
        ban: '',
        is_admin: '',
        active: '',
        cash: '',
        danap: ''
    };
    inputs.forEach(inp => {
        const field = inp.getAttribute('data-field');
        if (field) {
            let value = inp.value;
            if (field === 'cash' || field === 'danap') {
                value = value.replace(/\./g, '');
            }
            data[field] = value;
        }
    });

    if (!data.username) data.username = originalValues.username || '';
    if (!data.password) data.password = originalValues.password || '';
    if (!data.ban) data.ban = originalValues.ban || '0';
    if (!data.is_admin) data.is_admin = originalValues.is_admin || '0';
    if (!data.active) data.active = originalValues.active || '0';
    if (!data.cash) data.cash = originalValues.cash || '0';
    if (!data.danap) data.danap = originalValues.danap || '0';

    showLoading(true);
    try {
        const res = await fetch('account.php', {
            method: 'POST',
            headers: csrfHeaders({ 'Content-Type': 'application/x-www-form-urlencoded' }),
            body: new URLSearchParams(data)
        });
        const json = await res.json();
        if (json.success) {
            const updatedFields = {};
            inputs.forEach(inp => {
                const field = inp.getAttribute('data-field');
                let newVal = inp.value;
                if (field === 'cash' || field === 'danap') {
                    newVal = newVal.replace(/\./g, '');
                    updatedFields[field] = newVal;
                } else if (field === 'ban' || field === 'is_admin' || field === 'active') {
                    updatedFields[field] = newVal;
                } else {
                    updatedFields[field] = newVal;
                }
            });
            const spans = tr.querySelectorAll('.view-value[data-field]');
            spans.forEach(span => {
                const field = span.getAttribute('data-field');
                const newVal = updatedFields[field];
                if (newVal !== undefined) {
                    if (field === 'ban') {
                        span.innerText = newVal == '1' ? '1 - Khóa' : '0 - Mở';
                        span.setAttribute('data-original', newVal);
                    } else if (field === 'is_admin') {
                        span.innerText = newVal == '1' ? '1 - Admin' : '0 - User';
                        span.setAttribute('data-original', newVal);
                    } else if (field === 'active') {
                        span.innerText = newVal == '1' ? '1 - Mở' : '0 - Đóng';
                        span.setAttribute('data-original', newVal);
                    } else if (field === 'cash' || field === 'danap') {
                        const formatted = Number(newVal).toLocaleString('vi-VN');
                        span.innerText = formatted;
                        span.setAttribute('data-original', newVal);
                    } else {
                        span.innerText = newVal;
                        span.setAttribute('data-original', newVal);
                    }
                }
            });
            tr.querySelectorAll('.view-value').forEach(span => span.style.display = 'inline-block');

            const newBan = updatedFields.ban !== undefined ? parseInt(updatedFields.ban) : (originalValues.ban || 0);
            const newActive = updatedFields.active !== undefined ? parseInt(updatedFields.active) : (originalValues.active || 0);
            const newAdmin = updatedFields.is_admin !== undefined ? parseInt(updatedFields.is_admin) : (originalValues.is_admin || 0);
            tr.className = newBan ? 'row-ban' : (newAdmin ? 'row-admin' : (newActive ? 'row-active1' : 'row-active0'));
            tr.setAttribute('data-id', id);
          const username = data.username || originalValues.username || '';
showToast('✅ Đã lưu Tài khoản: ' + username);

        } else {
            showToast('❌ Lỗi: ' + (json.error||'Không rõ'), 'error');
            showLoading(false);
            return;
        }
    } catch(e) {
        showToast('❌ Lỗi kết nối!', 'error');
        showLoading(false);
        return;
    }
    showLoading(false);
    finishEdit(tr);
}

function cancelEdit(tr, originalValues) {
    const inputs = tr.querySelectorAll('.edit-input');
    inputs.forEach(inp => inp.remove());
    const spans = tr.querySelectorAll('.view-value');
    spans.forEach(span => span.style.display = 'inline-block');
    finishEdit(tr);
}

function finishEdit(tr) {
    const actionDiv = tr.querySelector('.action-wrap');
    const saveBtn = actionDiv.querySelector('.btn-save-p');
    const cancelBtn = actionDiv.querySelector('.btn-cancel-p');
    if (saveBtn) saveBtn.remove();
    if (cancelBtn) cancelBtn.remove();
    if (!actionDiv.querySelector('.btn-edit')) {
        const editBtn = document.createElement('button');
        editBtn.className = 'btn-sm btn-edit';
        editBtn.innerHTML = '✏️ Sửa';
        editBtn.onclick = function() { editRow(this); };
        actionDiv.prepend(editBtn);
    }
    const delBtn = tr.querySelector('.btn-del-p');
    if (delBtn) delBtn.style.display = 'inline-flex';
    tr.classList.remove('editing');
    tr.querySelectorAll('.edit-input').forEach(inp => inp.remove());
}

/* ================= FILTER ================= */
let debounce;
function submitFilter() {
    clearTimeout(debounce);
    debounce = setTimeout(async () => {
        const form = document.getElementById('filterForm');
        const params = new URLSearchParams(new FormData(form));
        document.querySelectorAll('thead tr.filter-row input').forEach(inp => {
            if (inp.name) params.set(inp.name, inp.value);
        });
        try {
            const res = await fetch('account.php?' + params.toString());
            const html = await res.text();
            const doc = new DOMParser().parseFromString(html, 'text/html');
            const newTbody = doc.querySelector('tbody');
            if (newTbody) document.querySelector('tbody').innerHTML = newTbody.innerHTML;
            setTimeout(applyScale, 50);
        } catch(e) { console.log('Lỗi filter'); }
    }, 300);
}
</script>
<?php require __DIR__ . '/guide_widget.php'; ?>
</div><!-- /kadmin-body -->
<?php
admin_foot();
