<?php
session_start();
require_once $_SERVER['DOCUMENT_ROOT'] . "/sv1/config.php";
require_once __DIR__ . '/admin_layout.php';

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

$csrfToken = ensure_csrf_token();

// ── Handle AJAX update ──
if ($_SERVER['REQUEST_METHOD'] === 'POST' && isset($_POST['action']) && $_POST['action'] === 'update') {
    require_csrf_token();
    ob_clean();
    header('Content-Type: application/json; charset=utf-8');

    $rawId = $_POST['id'] ?? '';
    $id    = intval($rawId);

    if ($id < 0) {
        echo json_encode(['status' => 'error', 'message' => 'ID không hợp lệ (nhận: ' . htmlspecialchars($rawId) . ')']);
        exit;
    }

    $type        = trim($_POST['type']          ?? '');
    $gender      = trim($_POST['gender']        ?? '');
    $name        = trim($_POST['name']          ?? '');
    $description = trim($_POST['description']   ?? '');
    $icon_id     = intval($_POST['icon_id']     ?? 0);
    $head        = intval($_POST['head']        ?? 0);
    $body        = intval($_POST['body']        ?? 0);
    $leg         = intval($_POST['leg']         ?? 0);
    $is_up_to_up = intval($_POST['is_up_to_up'] ?? 0);

    try {
        $stmt = $pdo->prepare("UPDATE item_template
            SET TYPE=:type, gender=:gender, NAME=:name, description=:desc,
                icon_id=:icon_id, head=:head, body=:body, leg=:leg, is_up_to_up=:iutu
            WHERE id=:id");
        $stmt->execute([
            ':type'   => $type,   ':gender' => $gender,
            ':name'   => $name,   ':desc'   => $description,
            ':icon_id'=> $icon_id,':head'   => $head,
            ':body'   => $body,   ':leg'    => $leg,
            ':iutu'   => $is_up_to_up, ':id' => $id,
        ]);
        echo json_encode(['status' => 'success', 'id' => $id]);
    } catch (Exception $e) {
        echo json_encode(['status' => 'error', 'message' => $e->getMessage()]);
    }
    exit;
}

// ── Pagination & filter ──
$page    = max(1, intval($_GET['page']   ?? 1));
$perPage = 500;
$offset  = ($page - 1) * $perPage;

$filterCol = $_GET['filter_col'] ?? '';
$filterVal = trim($_GET['filter_val'] ?? '');

$allowedCols = ['id','TYPE','gender','NAME','description','icon_id','head','body','leg','is_up_to_up'];

$where  = '';
$params = [];

if ($filterVal !== '') {
    if ($filterCol !== '' && in_array($filterCol, $allowedCols)) {
        $where = "WHERE `$filterCol` LIKE :fv";
        $params[':fv'] = '%' . $filterVal . '%';
    } else {
        $where = "WHERE (id LIKE :fv OR TYPE LIKE :fv OR gender LIKE :fv OR NAME LIKE :fv
                      OR description LIKE :fv OR icon_id LIKE :fv OR head LIKE :fv
                      OR body LIKE :fv OR leg LIKE :fv OR is_up_to_up LIKE :fv)";
        $params[':fv'] = '%' . $filterVal . '%';
    }
}

$totalStmt = $pdo->prepare("SELECT COUNT(*) FROM item_template $where");
$totalStmt->execute($params);
$total      = (int)$totalStmt->fetchColumn();
$totalPages = max(1, ceil($total / $perPage));

$params[':limit']  = $perPage;
$params[':offset'] = $offset;
$stmt = $pdo->prepare("SELECT * FROM item_template $where ORDER BY id ASC LIMIT :limit OFFSET :offset");
foreach ($params as $k => &$v) {
    $stmt->bindValue($k, $v, ($k === ':limit' || $k === ':offset') ? PDO::PARAM_INT : PDO::PARAM_STR);
}
unset($v);
$stmt->execute();
$items = $stmt->fetchAll();

// Quick view — all items for modal
$quickItems = $pdo->query("SELECT id, NAME, icon_id FROM item_template ORDER BY id ASC")->fetchAll();

$colOpts = ['id'=>'ID','TYPE'=>'Type','gender'=>'Gen','NAME'=>'Tên',
            'description'=>'Mô Tả','icon_id'=>'Icon_ID',
            'head'=>'H','body'=>'B','leg'=>'L','is_up_to_up'=>'Cộng dồn'];

admin_head('Quản Lý Vật Phẩm');
admin_bar('Quản Lý Vật Phẩm', 'package');
?>
<meta name="csrf-token" content="<?= htmlspecialchars($csrfToken, ENT_QUOTES) ?>">
<style>
/* ── Thành phần riêng của trang (class bị JS tra cứu — giữ tên, đổi định dạng bằng --kame-*) ── */
.W { position: relative; z-index: 1; }

/* Tìm kiếm theo cột trong header bảng */
.col-search { margin-top: 5px; position: relative; font-weight: 400; text-transform: none; letter-spacing: 0; }
.col-search input {
  width: 100%; padding: 4px 7px 4px 24px;
  border: 2px solid rgb(96 0 24 / 45%); border-radius: var(--kame-radius-sm);
  background: var(--kame-white); color: var(--kame-ink); font: inherit; font-size: 13px; outline: none;
}
.col-search input:focus { border-color: var(--kame-orange-600); outline: 3px solid rgb(240 96 48 / 20%); }
.col-search .cs-ic { position: absolute; left: 6px; top: 50%; transform: translateY(-50%); font-size: .68rem; pointer-events: none; }

/* Ô đếm kết quả tìm */
.cnt-tr td { padding: 4px 12px; background: var(--kame-paper-300); text-align: center; font-size: 12px; color: var(--kame-maroon-700); }

/* Ô ID + icon */
.cell-id { color: var(--kame-violet); font-weight: 800; font-family: "Roboto Condensed", monospace; font-size: 13px; }
.icon-img { width: 32px; height: 32px; object-fit: contain; image-rendering: pixelated; display: block; }
.icon-ph { width: 32px; height: 32px; display: none; align-items: center; justify-content: center; color: var(--kame-maroon-700); font-size: .6rem; }

/* Ô nhập chỉnh sửa (khoá mặc định, mở khi hàng ở chế độ sửa) */
.editable {
  background: transparent; border: 2px solid transparent; color: var(--kame-ink);
  font: inherit; font-size: 13px; padding: 3px 6px; border-radius: 5px; width: 100%; text-align: center;
  transition: border-color .15s, background .15s, box-shadow .15s;
  pointer-events: none; user-select: none;
}
tr.editing .editable { pointer-events: auto; user-select: auto; border-color: rgb(96 0 24 / 40%); background: var(--kame-white); }
tr.editing .editable:focus { outline: none; border-color: var(--kame-orange-600); background: var(--kame-white); box-shadow: 0 0 0 3px rgb(240 96 48 / 20%); }
textarea.editable { text-align: left; resize: none; min-height: 2.8em; line-height: 1.45; overflow: hidden; white-space: pre-wrap; word-break: break-word; display: block; }
select.editable { cursor: default; }
tr.editing select.editable { cursor: pointer; }
select.editable.iutu-yes { color: var(--kame-green); font-weight: 700; }
input[type=number].editable::-webkit-inner-spin-button,
input[type=number].editable::-webkit-outer-spin-button { -webkit-appearance: none; margin: 0; }
input[type=number].editable { -moz-appearance: textfield; }

/* Nút Sửa / Lưu / Huỷ */
.action-btn {
  display: inline-flex; align-items: center; justify-content: center; gap: 5px;
  border: none; border-radius: 5px; padding: 6px 10px;
  font: inherit; font-size: 12px; font-weight: 800; text-transform: uppercase;
  color: var(--kame-white); cursor: pointer; white-space: nowrap; width: 100%;
  transition: filter .15s, transform .15s;
}
.action-btn:hover { filter: brightness(1.08); transform: translateY(-1px); }
.action-btn.edit-mode { background: var(--kame-orange-600); }
.action-btn.save-mode { background: var(--kame-green); }
.action-btn.saving-mode { background: var(--kame-gold); color: var(--kame-ink); }
.cancel-btn {
  display: none; align-items: center; justify-content: center; gap: 4px;
  border: none; border-radius: 5px; padding: 6px 8px;
  font: inherit; font-size: 11px; font-weight: 800;
  color: var(--kame-white); background: var(--kame-red); cursor: pointer; white-space: nowrap; width: 100%; margin-top: 4px;
  transition: filter .15s, transform .15s;
}
.cancel-btn:hover { filter: brightness(1.08); transform: translateY(-1px); }
tr.editing .cancel-btn { display: inline-flex; }

/* Toast */
#toast {
  position: fixed; bottom: 32px; right: 32px; z-index: 9999;
  font-size: 1.15rem; padding: 16px 26px; border-radius: var(--kame-radius-md);
  background: var(--kame-maroon-900); color: var(--kame-white); font-weight: 800;
  border: 2px solid var(--kame-orange-600); box-shadow: var(--kame-shadow-md);
  opacity: 0; transform: translateY(40px) scale(.9); transition: all .25s ease; pointer-events: none; max-width: 480px;
}
#toast.show { transform: translateY(0) scale(1.05); opacity: 1; }
#toast.ok { border-color: var(--kame-green); }
#toast.err { border-color: var(--kame-red); }

/* Modal check nhanh */
#modal-overlay { position: fixed; inset: 0; z-index: 200; background: rgb(40 21 15 / 62%); display: none; align-items: center; justify-content: center; }
#modal-overlay.open { display: flex; }
#modal-box {
  background: var(--kame-paper-100); border: 3px solid var(--kame-maroon-700); border-radius: var(--kame-radius-md);
  width: 80vw; height: 90vh; box-shadow: 0 20px 50px rgb(40 21 15 / 36%); display: flex; flex-direction: column; overflow: hidden;
}
.modal-head { display: flex; align-items: center; gap: 10px; padding: 15px 20px; border-bottom: 2px solid var(--kame-orange-600); flex-shrink: 0; }
.modal-head h2 { margin: 0; flex: 1; font-family: "Roboto Condensed", sans-serif; font-weight: 900; font-size: 1.1rem; color: var(--kame-maroon-700); }
#modal-search {
  background: var(--kame-white); border: 2px solid rgb(96 0 24 / 40%); color: var(--kame-ink);
  border-radius: var(--kame-radius-sm); padding: 8px 13px; font: inherit; font-size: .9rem; width: 230px; outline: none;
}
#modal-search:focus { border-color: var(--kame-orange-600); }
.modal-close {
  background: var(--kame-red); border: none; color: var(--kame-white); font-size: 1.15rem; font-weight: 900;
  width: 38px; height: 38px; border-radius: 50%; cursor: pointer; display: flex; align-items: center; justify-content: center; flex-shrink: 0;
  transition: filter .15s, transform .15s;
}
.modal-close:hover { filter: brightness(1.1); transform: scale(1.08); }
.modal-bar { display: flex; align-items: center; gap: 6px; flex-wrap: wrap; padding: 6px 18px; border-bottom: 1px solid rgb(96 0 24 / 22%); font-size: .8rem; color: var(--kame-maroon-700); flex-shrink: 0; }
.modal-body { overflow-y: auto; padding: 14px 16px; flex: 1; }
.item-grid { display: grid; grid-template-columns: repeat(12, 1fr); gap: 8px; }
.item-card {
  background: #fffdf7; border: 1px solid var(--kame-maroon-700); border-radius: var(--kame-radius-sm);
  padding: 9px 5px 7px; display: flex; flex-direction: column; align-items: center; gap: 4px; text-align: center;
  transition: border-color .15s, transform .15s, box-shadow .15s;
}
.item-card:hover { border-color: var(--kame-orange-600); transform: translateY(-2px); box-shadow: var(--kame-shadow-sm); }
.item-card img { width: 40px; height: 40px; object-fit: contain; image-rendering: pixelated; display: block; }
.item-card .ic-ph { width: 40px; height: 40px; background: var(--kame-paper-300); border-radius: 4px; display: none; }
.item-card .ic-id { font-size: .67rem; color: var(--kame-violet); font-weight: 700; font-family: "Roboto Condensed", monospace; line-height: 1; }
.item-card .ic-name { font-size: .63rem; color: var(--kame-maroon-700); line-height: 1.2; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 90px; width: 100%; }

@media (max-width: 900px) { .item-grid { grid-template-columns: repeat(6, 1fr); } }
</style>
<script>
(function(){
  function setPageZoom(){
    var zoom = (document.documentElement.clientWidth * 0.98) / 1500;
    document.documentElement.style.setProperty('--page-zoom', zoom);
  }
  setPageZoom();
  window.__setAdminPageZoom = setPageZoom;
})();
</script>
<div class="kadmin-body">
<div class="W" id="W">

  <!-- Thanh công cụ -->
  <div class="kadmin-card">
    <div class="kadmin-toolbar">
      <button class="kame-button" onclick="openModal()">
        <i class="ti ti-list-search" aria-hidden="true"></i>
        Check Item: <span style="font-family:'Roboto Condensed',monospace"><?= number_format($total) ?></span>
      </button>
      <button class="kame-button kame-button--secondary" onclick="location.reload()">
        <i class="ti ti-refresh" aria-hidden="true"></i> Làm Mới
      </button>
    </div>
  </div>

  <!-- PHÂN TRANG -->
  <?php if ($totalPages > 1): ?>
  <div class="kame-pagination" style="margin-bottom:var(--kame-space-4)">
    <?php
      $pS = max(1,$page-3); $pE = min($totalPages,$page+3);
      if ($pS > 1){ echo '<a class="kame-page" href="?page=1&filter_col='.urlencode($filterCol).'&filter_val='.urlencode($filterVal).'">1</a>'; if($pS>2) echo '<span class="kame-muted">…</span>'; }
      for ($i=$pS;$i<=$pE;$i++):
    ?>
      <a class="kame-page"<?= $i===$page?' aria-current="page"':'' ?> href="?page=<?= $i ?>&filter_col=<?= urlencode($filterCol) ?>&filter_val=<?= urlencode($filterVal) ?>"><?= $i ?></a>
    <?php endfor;
      if ($pE < $totalPages){ if($pE<$totalPages-1) echo '<span class="kame-muted">…</span>'; echo '<a class="kame-page" href="?page='.$totalPages.'&filter_col='.urlencode($filterCol).'&filter_val='.urlencode($filterVal).'">'.$totalPages.'</a>'; }
    ?>
  </div>
  <?php endif; ?>

  <!-- BẢNG -->
  <div class="kadmin-table-wrap">
  <table class="kadmin-table">
    <colgroup>
      <col class="c-id">
      <col class="c-icon">
      <col class="c-type">
      <col class="c-gen">
      <col class="c-name">
      <col class="c-desc">
      <col class="c-iconid">
      <col class="c-head">
      <col class="c-body">
      <col class="c-leg">
      <col class="c-don">
      <col class="c-act">
    </colgroup>
    <thead>
      <tr>
        <th class="c-id-th">ID
          <div class="col-search">
            <span class="cs-ic">🔢</span>
            <input type="text" data-col="id" placeholder="ID…" oninput="colSearch()">
          </div>
        </th>
        <th>Icon</th>
        <th>Type
          <div class="col-search">
            <span class="cs-ic">🏷</span>
            <input type="text" data-col="type" placeholder="…" oninput="colSearch()">
          </div>
        </th>
        <th>Gen
          <div class="col-search">
            <span class="cs-ic">⚥</span>
            <input type="text" data-col="gender" placeholder="…" oninput="colSearch()">
          </div>
        </th>
        <th>Tên
          <div class="col-search">
            <span class="cs-ic">🔍</span>
            <input type="text" data-col="name" placeholder="Tên…" oninput="colSearch()">
          </div>
        </th>
        <th>Mô Tả
          <div class="col-search">
            <span class="cs-ic">🔍</span>
            <input type="text" data-col="desc" placeholder="Mô tả…" oninput="colSearch()">
          </div>
        </th>
        <th>Icon_ID
          <div class="col-search">
            <span class="cs-ic">🖼</span>
            <input type="text" data-col="icon_id" placeholder="Icon ID…" oninput="colSearch()">
          </div>
        </th>
        <th>Head</th>
        <th>Body</th>
        <th>Leg</th>
        <th>Up To</th>
        <th>VIP</th>
      </tr>
      <tr class="cnt-tr" id="cnt-tr" style="display:none">
        <td colspan="12"><span id="cnt-sc"></span></td>
      </tr>
    </thead>
    <tbody id="tb">
      <?php foreach ($items as $row):
        $iutu = intval($row['is_up_to_up']);
      ?>
      <tr class="item-row" data-id="<?= intval($row['id']) ?>"
          data-id-val="<?= intval($row['id']) ?>"
          data-type-val="<?= mb_strtolower(htmlspecialchars($row['TYPE'] ?? '')) ?>"
          data-gender-val="<?= mb_strtolower(htmlspecialchars($row['gender'] ?? '')) ?>"
          data-name-val="<?= mb_strtolower(htmlspecialchars($row['NAME'] ?? '')) ?>"
          data-desc-val="<?= mb_strtolower(htmlspecialchars($row['description'] ?? '')) ?>"
          data-name="<?= htmlspecialchars($row['NAME'] ?? '') ?>"
          data-icon-val="<?= intval($row['icon_id']) ?>">

        <td><span class="cell-id"><?= intval($row['id']) ?></span></td>
        <td>
          <div class="kadmin-icon">
            <img class="icon-img" src="/icon/<?= intval($row['icon_id']) ?>.png" alt=""
                 onerror="this.style.display='none';this.nextElementSibling.style.display='flex'">
            <div class="icon-ph">?</div>
          </div>
        </td>
        <td><input class="editable" name="type"   value="<?= htmlspecialchars($row['TYPE']   ?? '') ?>" readonly></td>
        <td><input class="editable" name="gender" value="<?= htmlspecialchars($row['gender'] ?? '') ?>" readonly></td>
        <td><input class="editable" name="name"   value="<?= htmlspecialchars($row['NAME']   ?? '') ?>" readonly></td>
        <td>
          <textarea class="editable" name="description" rows="2" readonly><?= htmlspecialchars($row['description'] ?? '') ?></textarea>
        </td>
        <td><input class="editable" name="icon_id" type="number" value="<?= intval($row['icon_id']) ?>" readonly></td>
        <td><input class="editable" name="head"    type="number" value="<?= intval($row['head']) ?>" readonly></td>
        <td><input class="editable" name="body"    type="number" value="<?= intval($row['body']) ?>" readonly></td>
        <td><input class="editable" name="leg"     type="number" value="<?= intval($row['leg']) ?>" readonly></td>
        <td>
          <select class="editable<?= $iutu===1?' iutu-yes':'' ?>" name="is_up_to_up" disabled
                  onchange="this.className='editable'+(this.value==='1'?' iutu-yes':'')">
            <option value="0" <?= $iutu===0?'selected':'' ?>>0 - No</option>
            <option value="1" <?= $iutu===1?'selected':'' ?>>1 - Có</option>
          </select>
        </td>
        <td>
          <button class="action-btn edit-mode" onclick="toggleEdit(this)">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"
                    stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"
                    stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            Sửa
          </button>
          <button class="cancel-btn" onclick="cancelEdit(this)">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="none">
              <path d="M18 6 6 18M6 6l12 12" stroke="white" stroke-width="2.5"
                    stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            Hủy
          </button>
        </td>
      </tr>
      <?php endforeach; ?>
    </tbody>
  </table>
  </div>

</div><!-- /W -->
</div><!-- /kadmin-body -->

<div id="toast"></div>

<!-- MODAL — chỉ đóng bằng nút X -->
<div id="modal-overlay">
  <div id="modal-box">
    <div class="modal-head">
      <h2>⚡ Check Item Nhanh</h2>
      <input id="modal-search" type="text" placeholder="🔍 ID hoặc tên…" oninput="filterModal()">
      <button class="modal-close" onclick="closeModal()" title="Đóng">✕</button>
    </div>
    <div class="modal-bar" id="modal-bar" style="display:none"></div>
    <div class="modal-body">
      <div class="item-grid" id="modal-grid"></div>
    </div>
  </div>
</div>

<script>
// ══════════════════════════════
//  RESPONSIVE SCALE — 1500px always fills 98vw
//  Uses CSS zoom: browser handles reflow natively
// ══════════════════════════════
const PAGE_W = 1500;

function applyScale() {
  if (window.__setAdminPageZoom) window.__setAdminPageZoom();
}

applyScale();
window.addEventListener('resize', applyScale);

// ══════════════════════════════
//  COLUMN SEARCH (client-side)
// ══════════════════════════════
const trs   = document.querySelectorAll('#tb tr.item-row');
const cntTr = document.getElementById('cnt-tr');
const cntSc = document.getElementById('cnt-sc');

function colSearch() {
  const filters = {};
  document.querySelectorAll('.col-search input').forEach(inp => {
    const v = inp.value.toLowerCase().trim();
    if (v) filters[inp.dataset.col] = v;
  });

  let n = 0;
  trs.forEach(r => {
    let ok = true;
    for (const [col, val] of Object.entries(filters)) {
      const attr = r.dataset[col + 'Val'] ?? '';
      if (!attr.includes(val)) { ok = false; break; }
    }
    r.style.display = ok ? '' : 'none';
    if (ok) n++;
  });

  if (Object.keys(filters).length) {
    cntSc.textContent = `${n} / ${trs.length} kết quả`;
    cntTr.style.display = '';
  } else {
    cntTr.style.display = 'none';
  }
}

// ══════════════════════════════
//  TOGGLE EDIT / SAVE
// ══════════════════════════════
function toggleEdit(btn) {
  const tr = btn.closest('tr');

  if (btn.classList.contains('edit-mode')) {
    // Snapshot current values for cancel
    tr.querySelectorAll('.editable').forEach(el => {
      el.dataset.orig = el.value;
    });

    tr.classList.add('editing');
    tr.querySelectorAll('.editable').forEach(el => {
      if (el.tagName === 'SELECT') el.disabled = false;
      else el.removeAttribute('readonly');
    });

    btn.classList.remove('edit-mode');
    btn.classList.add('save-mode');
    btn.innerHTML = `
      <svg width="12" height="12" viewBox="0 0 24 24" fill="none">
        <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"
              stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        <polyline points="17 21 17 13 7 13 7 21" stroke="white" stroke-width="2"
              stroke-linecap="round" stroke-linejoin="round"/>
        <polyline points="7 3 7 8 15 8" stroke="white" stroke-width="2"
              stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      Lưu`;

  } else if (btn.classList.contains('save-mode')) {
    saveRow(btn);
  }
}

// ══════════════════════════════
//  CANCEL EDIT — restore snapshot
// ══════════════════════════════
function cancelEdit(cancelBtn) {
  const tr = cancelBtn.closest('tr');
  tr.querySelectorAll('.editable').forEach(el => {
    if (el.dataset.orig !== undefined) el.value = el.dataset.orig;
  });
  lockRow(tr);
}

function lockRow(tr) {
  tr.classList.remove('editing');
  tr.querySelectorAll('.editable').forEach(el => {
    if (el.tagName === 'SELECT') el.disabled = true;
    else el.setAttribute('readonly', '');
  });
  const saveBtn = tr.querySelector('.action-btn');
  if (saveBtn) {
    saveBtn.classList.remove('save-mode','saving-mode');
    saveBtn.classList.add('edit-mode');
    saveBtn.innerHTML = `
      <svg width="12" height="12" viewBox="0 0 24 24" fill="none">
        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"
              stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"
              stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      Sửa`;
  }
}

// ══════════════════════════════
//  SAVE ROW
// ══════════════════════════════
function saveRow(btn) {
  const tr = btn.closest('tr');
  const id = parseInt(tr.dataset.id, 10);
  if (isNaN(id) || id < 0) { toast('❌ Không đọc được ID hàng!', 'err'); return; }

  const g = n => tr.querySelector(`[name="${n}"]`).value;

  const body = new URLSearchParams({
    action: 'update', id,
    type: g('type'), gender: g('gender'), name: g('name'),
    description: g('description'), icon_id: g('icon_id'),
    head: g('head'), body: g('body'), leg: g('leg'),
    is_up_to_up: g('is_up_to_up'),
  });

  btn.classList.remove('save-mode');
  btn.classList.add('saving-mode');
  btn.innerHTML = `
    <svg width="12" height="12" viewBox="0 0 24 24" fill="none">
      <circle cx="12" cy="12" r="9" stroke="white" stroke-width="2" stroke-dasharray="28" stroke-dashoffset="8">
        <animateTransform attributeName="transform" type="rotate" from="0 12 12" to="360 12 12" dur=".7s" repeatCount="indefinite"/>
      </circle>
    </svg>
    …`;

  fetch('item', {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded', 'X-CSRF-TOKEN': document.querySelector('meta[name="csrf-token"]')?.content || '' },
    body: body.toString(),
  })
  .then(r => r.text())
  .then(raw => {
    let res;
    try { res = JSON.parse(raw); }
    catch { throw new Error('Phản hồi lạ: ' + raw.substring(0, 80)); }

    if (res.status === 'success') {
      lockRow(tr);
      const name = tr.dataset.name || `ID ${id}`;
      toast(`✅ Đã lưu: ${name}`, 'ok');
    } else {
      btn.classList.remove('saving-mode');
      btn.classList.add('save-mode');
      btn.innerHTML = `<svg width="12" height="12" viewBox="0 0 24 24" fill="none">
        <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"
              stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        <polyline points="17 21 17 13 7 13 7 21" stroke="white" stroke-width="2"/>
        <polyline points="7 3 7 8 15 8" stroke="white" stroke-width="2"/>
      </svg>Lưu`;
      toast('❌ ' + (res.message || 'Lỗi'), 'err');
    }
  })
  .catch(err => {
    btn.classList.remove('saving-mode');
    btn.classList.add('save-mode');
    btn.innerHTML = `<svg width="12" height="12" viewBox="0 0 24 24" fill="none">
      <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"
            stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      <polyline points="17 21 17 13 7 13 7 21" stroke="white" stroke-width="2"/>
      <polyline points="7 3 7 8 15 8" stroke="white" stroke-width="2"/>
    </svg>Lưu`;
    toast('❌ ' + err.message, 'err');
  });
}

// ══════════════════════════════
//  TOAST
// ══════════════════════════════
let toastTimer;
function toast(msg, type = 'ok') {
  const el = document.getElementById('toast');
  el.textContent = msg;
  el.className = 'show ' + type;
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => { el.className = ''; }, 2000);
}

// ══════════════════════════════
//  MODAL — chỉ đóng bằng nút X
//  KHÔNG click ngoài, KHÔNG Escape
// ══════════════════════════════
const ALL_ITEMS = <?= json_encode($quickItems, JSON_UNESCAPED_UNICODE) ?>;
let mItems = [...ALL_ITEMS];

function openModal() {
  mItems = [...ALL_ITEMS];
  document.getElementById('modal-search').value = '';
  renderModal();
  document.getElementById('modal-overlay').classList.add('open');
  document.body.style.overflow = 'hidden';
  // focus vào ô search để gõ ngay
  setTimeout(() => document.getElementById('modal-search').focus(), 80);
}

function closeModal() {
  document.getElementById('modal-overlay').classList.remove('open');
  document.body.style.overflow = '';
}

function filterModal() {
  const q = document.getElementById('modal-search').value.toLowerCase().trim();
  mItems = q
    ? ALL_ITEMS.filter(it => String(it.id).includes(q) || (it.NAME || '').toLowerCase().includes(q))
    : [...ALL_ITEMS];
  renderModal();

  // cập nhật bar đếm
  const bar = document.getElementById('modal-bar');
  if (q) {
    bar.style.display = 'flex';
    bar.textContent = `${mItems.length} / ${ALL_ITEMS.length} kết quả`;
  } else {
    bar.style.display = 'none';
  }
}

function renderModal() {
  document.getElementById('modal-grid').innerHTML = mItems.map(it => {
    const name = it.NAME ? it.NAME.substring(0, 18) : '—';
    return `<div class="item-card" title="#${it.id} — ${it.NAME || ''}">
      <img src="/icon/${it.icon_id}.png" alt=""
           onerror="this.style.display='none';this.nextElementSibling.style.display='block'">
      <div class="ic-ph"></div>
      <span class="ic-id">#${it.id}</span>
      <span class="ic-name">${name}</span>
    </div>`;
  }).join('');
}
</script>
<?php require __DIR__ . '/guide_widget.php'; ?>
<?php
admin_foot();
