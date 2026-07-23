<?php
session_start();
require_once $_SERVER['DOCUMENT_ROOT'] . "/sv1/config.php";
require_once __DIR__ . "/admin_layout.php";

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

// ========== XỬ LÝ POST ==========
if ($_SERVER['REQUEST_METHOD'] === 'POST' && isset($_POST['action'])) {
    require_csrf_token();
    header('Content-Type: application/json');

    if ($_POST['action'] === 'save_intrinsic') {
        $id           = (int)$_POST['intrinsic_id'];
        $name         = trim($_POST['name']);
        $param_from_1 = (int)$_POST['param_from_1'];
        $param_to_1   = (int)$_POST['param_to_1'];
        $param_from_2 = (int)$_POST['param_from_2'];
        $param_to_2   = (int)$_POST['param_to_2'];
        $icon         = (int)$_POST['icon'];

        try {
            $pdo->prepare(
                "UPDATE intrinsic SET NAME=?, param_from_1=?, param_to_1=?, param_from_2=?, param_to_2=?, icon=? WHERE id=?"
            )->execute([$name, $param_from_1, $param_to_1, $param_from_2, $param_to_2, $icon, $id]);
            echo json_encode(['ok' => true]);
        } catch (Exception $e) {
            echo json_encode(['ok' => false, 'msg' => $e->getMessage()]);
        }
        exit;
    }

    echo json_encode(['ok' => false, 'msg' => 'Unknown action']);
    exit;
}

// === LOAD DATA - không phân trang, bỏ id = 0 ===
$list = $pdo->query(
    "SELECT id, NAME, param_from_1, param_to_1, param_from_2, param_to_2, icon, gender
     FROM intrinsic WHERE id > 0 ORDER BY id ASC"
)->fetchAll(PDO::FETCH_ASSOC);

$total = count($list); // tổng số sau khi lọc

admin_head('Quản Lý Nội Tại');
admin_bar('Quản Lý Nội Tại', 'bolt');
?>
<meta name="csrf-token" content="<?= htmlspecialchars($csrfToken, ENT_QUOTES) ?>">
<script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
<div class="kadmin-body">
<div id="scaler">
<div class="W">

    <div class="kadmin-stats">
        <div class="kadmin-stat">
            <div class="kadmin-stat__label">Tổng số nội tại</div>
            <div class="kadmin-stat__value"><?= number_format($total) ?></div>
        </div>
    </div>

    <div class="kadmin-card">
        <div class="kadmin-card__head"><i class="ti ti-bolt" aria-hidden="true"></i> Danh sách nội tại</div>

        <div class="kadmin-toolbar">
            <button class="kame-button kame-button--secondary" onclick="location.reload()">
                <i class="ti ti-refresh" aria-hidden="true"></i> Làm Mới
            </button>
        </div>

        <div class="kadmin-table-wrap">
        <table class="kadmin-table">
        <colgroup>
          <col class="col-id">
          <col class="col-name">
          <col class="col-p1">
          <col class="col-icon">
          <col class="col-gender">
          <col class="col-act">
        </colgroup>
        <thead>
        <tr>
          <th>ID</th>
          <th>
            <div style="display:flex;flex-direction:column;gap:6px;">
              <span>Tên Nội Tại</span>
              <input type="text" id="searchNT" class="kame-input" placeholder="🔍 Tìm tên...">
            </div>
          </th>
          <th>Chỉ Số<br><small style="font-weight:600;opacity:.7">Min / Max</small></th>
          <th>Icon</th>
          <th>Gender</th>
          <th>Hành Động</th>
        </tr>
        </thead>
        <tbody id="ntTableBody">
        <?php foreach ($list as $row):
          $g = (int)$row['gender'];
          $genderSrc = match($g) {
            0 => '/img/td.jpg',   // Trái Đất
            1 => '/img/nm.jpg',   // Namek
            2 => '/img/xd.jpg',   // Xayda
            3 => '/img/v1.png',   // Tất cả
            default => '/img/v1.png',
          };
          $genderLabel = match($g) {
            0 => 'Trái Đất',
            1 => 'Namek',
            2 => 'Xayda',
            3 => 'Tất cả',
            default => "Gender $g",
          };
          $nameEsc = htmlspecialchars($row['NAME'], ENT_QUOTES);
        ?>
        <tr class="nt-row"
            data-id="<?= $row['id'] ?>"
            data-name="<?= mb_strtolower(htmlspecialchars($row['NAME'])) ?>"
            data-editing="false"
            data-orig-name="<?= $nameEsc ?>"
            data-orig-pf1="<?= (int)$row['param_from_1'] ?>"
            data-orig-pt1="<?= (int)$row['param_to_1'] ?>"
            data-orig-pf2="<?= (int)$row['param_from_2'] ?>"
            data-orig-pt2="<?= (int)$row['param_to_2'] ?>"
            data-orig-icon="<?= (int)$row['icon'] ?>">

          <!-- ID -->
          <td class="num">
            <span style="font-weight:700;">#<?= $row['id'] ?></span>
          </td>

          <!-- NAME -->
          <td>
            <textarea class="cell-textarea" disabled><?= $nameEsc ?></textarea>
          </td>

          <!-- PARAM 1 -->
          <td>
            <div class="param-block">
              <div class="param-row">
                <label>From</label>
                <input type="number" class="param-input inp-pf1" value="<?= (int)$row['param_from_1'] ?>" min="0" disabled>
              </div>
              <div class="param-row">
                <label>To</label>
                <input type="number" class="param-input inp-pt1" value="<?= (int)$row['param_to_1'] ?>" min="0" disabled>
              </div>
            </div>
          </td>

          <!-- ICON -->
          <td>
            <div class="icon-wrap">
              <img class="icon-img icon-preview"
                   src="/icon/<?= (int)$row['icon'] ?>.png"
                   onerror="this.src='/icon/1.png'"
                   title="Icon #<?= (int)$row['icon'] ?>">
              <input type="number" class="icon-num-input inp-icon"
                     value="<?= (int)$row['icon'] ?>" min="0" disabled>
            </div>
          </td>

          <!-- GENDER -->
          <td>
            <div style="display:flex;flex-direction:column;align-items:center;gap:3px;">
              <img class="gender-img" src="<?= $genderSrc ?>" onerror="this.style.display='none'" title="<?= $genderLabel ?>">
              <span class="gender-label"><?= $genderLabel ?></span>
            </div>
          </td>

          <!-- ACTION -->
          <td class="act">
            <div class="action-cell">
              <button class="btn-sm btn-edit-p" onclick="enableEdit(this)">✏️ Sửa</button>
            </div>
          </td>
        </tr>
        <?php endforeach; ?>
        </tbody>
        </table>
        </div><!-- /kadmin-table-wrap -->
    </div><!-- /kadmin-card -->

</div><!-- /W -->
</div><!-- /scaler -->
</div><!-- /kadmin-body -->

<div id="toast"></div>

<style>
/* Chi con dinh dang cho phan chinh sua trong bang (JS phu thuoc cac class nay) */
.cell-textarea {
  width: 100%; min-height: 52px;
  background: transparent; border: 1.5px solid transparent;
  border-radius: var(--kame-radius-sm); color: var(--kame-green);
  font: inherit; font-weight: 800;
  padding: 2px; outline: none; resize: none;
  white-space: pre-wrap; word-break: break-word; line-height: 1.45;
}
.cell-textarea:not(:disabled) {
  background: rgb(255 253 245 / 78%); border-color: rgb(96 0 24 / 40%);
  padding: 5px 7px; resize: vertical; cursor: text;
}
.cell-textarea:not(:disabled):focus { border-color: var(--kame-orange-600); background: var(--kame-white); }
.param-block { display: flex; flex-direction: column; gap: 4px; width: 100%; }
.param-row   { display: flex; align-items: center; gap: 4px; }
.param-row label { font-size: 12px; font-weight: 700; color: var(--kame-maroon-700); width: 30px; text-align: right; flex-shrink: 0; }
.param-input {
  flex: 1; background: transparent; border: 1.5px solid transparent;
  border-radius: var(--kame-radius-sm); color: var(--kame-ink);
  font: inherit; font-weight: 800;
  padding: 2px 4px; outline: none; text-align: center;
}
.param-input:not(:disabled) { background: rgb(255 253 245 / 78%); border-color: rgb(96 0 24 / 40%); }
.param-input:not(:disabled):focus { border-color: var(--kame-orange-600); background: var(--kame-white); }
.icon-wrap { display: flex; flex-direction: column; align-items: center; gap: 4px; }
.icon-img  { width: 48px; height: 48px; object-fit: contain; border-radius: 8px; border: 2px solid var(--kame-maroon-700); image-rendering: pixelated; }
.icon-num-input {
  width: 72px; text-align: center;
  background: transparent; border: 1.5px solid transparent;
  border-radius: var(--kame-radius-sm); color: var(--kame-ink);
  font: inherit; font-weight: 800;
  padding: 2px 4px; outline: none;
}
.icon-num-input:not(:disabled) { background: rgb(255 253 245 / 78%); border-color: rgb(96 0 24 / 40%); }
.icon-num-input:not(:disabled):focus { border-color: var(--kame-orange-600); background: var(--kame-white); }
.gender-img   { width: 46px; height: 46px; object-fit: cover; border-radius: 8px; border: 2px solid var(--kame-maroon-700); }
.gender-label { font-size: 12px; color: var(--kame-maroon-700); font-weight: 600; margin-top: 2px; }
.action-cell { display: flex; flex-direction: column; gap: 5px; align-items: stretch; }
.btn-sm {
  display: inline-flex; align-items: center; justify-content: center;
  width: 100%; padding: 6px 0; font-size: 13px; font-weight: 700;
  border-radius: var(--kame-radius-sm); border: 2px solid var(--kame-maroon-700);
  cursor: pointer; font: inherit; font-weight: 700;
}
.btn-edit-p   { background: var(--kame-paper-100); color: var(--kame-blue); }
.btn-save-p   { background: var(--kame-paper-100); color: var(--kame-green); }
.btn-cancel-p { background: var(--kame-paper-100); color: var(--kame-red); }
#toast {
  position: fixed; bottom: 32px; right: 32px; z-index: 9999;
  font-size: 1.2rem; padding: 18px 28px; border-radius: var(--kame-radius-md);
  background: var(--kame-maroon-900); color: var(--kame-white); font-weight: 800;
  border: 2px solid var(--kame-orange-600);
  display: none; max-width: 480px; line-height: 1.2;
}
#toast.error { background: var(--kame-red); }
</style>

<script>
(function () {
  var PAGE_W = 1500, FILL = 0.98;
  var scaler = document.getElementById('scaler');
  function applyScale() {
    var vw    = window.innerWidth;
    var avail = vw * FILL;
    var ratio = avail / PAGE_W;
    scaler.style.width           = PAGE_W + 'px';
    scaler.style.transform       = 'scale(' + ratio + ')';
    scaler.style.transformOrigin = 'top left';
    scaler.style.marginLeft      = ((vw - avail) / 2) + 'px';
    document.body.style.minHeight = (scaler.offsetHeight * ratio) + 'px';
  }
  applyScale();
  window.addEventListener('resize', applyScale);
})();

function showToast(msg, isError) {
  const t = $('#toast');
  t.text(msg).removeClass('error');
  if (isError) t.addClass('error');
  t.stop(true).fadeIn(200).delay(2400).fadeOut(400);
}

$('#searchNT').on('keyup', function () {
  const val = $(this).val().toLowerCase().trim();
  $('#ntTableBody tr').each(function () {
    $(this).toggle(($(this).attr('data-name') || '').includes(val));
  });
});

function enableEdit(btn) {
  const $row = $(btn).closest('tr');
  if ($row.attr('data-editing') === 'true') return;
  $row.attr('data-editing', 'true');

  $row.find('.cell-textarea').prop('disabled', false);
  $row.find('.inp-pf1, .inp-pt1, .inp-icon').prop('disabled', false);

  const $ac = $row.find('.action-cell');
  $ac.find('.btn-edit-p').hide();
  $ac.append('<button class="btn-sm btn-save-p"   onclick="saveEdit(this)">💾 Lưu</button>');
  $ac.append('<button class="btn-sm btn-cancel-p" onclick="cancelEdit(this)">✖️ Hủy</button>');

  $row.find('.inp-icon').off('input').on('input', function () {
    const v = parseInt($(this).val()) || 0;
    $row.find('.icon-preview').attr('src', '/icon/' + v + '.png');
  });
}

function cancelEdit(btn) {
  const $row = $(btn).closest('tr');
  $row.find('.cell-textarea').val($row.attr('data-orig-name')).prop('disabled', true);
  $row.find('.inp-pf1').val($row.attr('data-orig-pf1')).prop('disabled', true);
  $row.find('.inp-pt1').val($row.attr('data-orig-pt1')).prop('disabled', true);
  $row.find('.inp-icon').val($row.attr('data-orig-icon')).prop('disabled', true);
  $row.find('.icon-preview').attr('src', '/icon/' + $row.attr('data-orig-icon') + '.png');
  $row.find('.btn-save-p, .btn-cancel-p').remove();
  $row.find('.btn-edit-p').show();
  $row.attr('data-editing', 'false');
}

async function saveEdit(btn) {
  const $row = $(btn).closest('tr');
  const id          = $row.data('id');
  const name        = $row.find('.cell-textarea').val().trim();
  const param_from_1 = parseInt($row.find('.inp-pf1').val()) || 0;
  const param_to_1   = parseInt($row.find('.inp-pt1').val()) || 0;
  // Giữ nguyên param2 từ data cũ
  const param_from_2 = parseInt($row.attr('data-orig-pf2')) || 0;
  const param_to_2   = parseInt($row.attr('data-orig-pt2')) || 0;
  const icon         = parseInt($row.find('.inp-icon').val()) || 0;

  try {
    const res = await $.post('', {
      csrf_token: document.querySelector('meta[name="csrf-token"]')?.content || '',
      action: 'save_intrinsic',
      intrinsic_id: id,
      name, param_from_1, param_to_1, param_from_2, param_to_2, icon
    });
    if (res.ok) {
      showToast('✅ Đã lưu nội tại #' + id + '!');
      $row.attr('data-orig-name', name);
      $row.attr('data-orig-pf1',  param_from_1);
      $row.attr('data-orig-pt1',  param_to_1);
      $row.attr('data-orig-icon', icon);
      $row.attr('data-name', name.toLowerCase());
      $row.find('.cell-textarea').prop('disabled', true);
      $row.find('.inp-pf1, .inp-pt1, .inp-icon').prop('disabled', true);
      $row.find('.btn-save-p, .btn-cancel-p').remove();
      $row.find('.btn-edit-p').show();
      $row.attr('data-editing', 'false');
    } else {
      showToast('❌ Lỗi: ' + (res.msg || 'Không rõ'), true);
    }
  } catch (e) {
    showToast('❌ Lỗi kết nối', true);
  }
}
</script>
<?php
require __DIR__ . '/guide_widget.php';
admin_foot();
