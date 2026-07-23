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

// Xử lý lưu dữ liệu từ AJAX (giống cơ chế của account.php)
if ($_SERVER['REQUEST_METHOD'] === 'POST' && isset($_POST['action']) && $_POST['action'] === 'save') {
    require_csrf_token();
    header('Content-Type: application/json');

    $id       = (int)$_POST['id'];
    $power    = (int)str_replace('.', '', $_POST['power']);
    $hp       = (int)str_replace('.', '', $_POST['hp']);
    $mp       = (int)str_replace('.', '', $_POST['mp']);
    $damage   = (int)str_replace('.', '', $_POST['damage']);
    $defense  = (int)str_replace('.', '', $_POST['defense']);
    $critical = (int)str_replace('.', '', $_POST['critical']);

    try {
        $stmt = $pdo->prepare("UPDATE power_limit SET
            power = ?, hp = ?, mp = ?, damage = ?, defense = ?, critical = ?
            WHERE id = ?");
        $stmt->execute([$power, $hp, $mp, $damage, $defense, $critical, $id]);
        echo json_encode(['success' => true]);
    } catch (Exception $e) {
        echo json_encode(['success' => false, 'error' => $e->getMessage()]);
    }
    exit;
}

// Lấy danh sách các giới hạn chỉ số
$rows = $pdo->query("SELECT id, power, hp, mp, damage, defense, critical FROM power_limit ORDER BY id ASC")->fetchAll(PDO::FETCH_ASSOC);
$totalRows = count($rows);

admin_head('Quản Lý Giới Hạn Chỉ Số');
admin_bar('Giới Hạn Chỉ Số', 'shield-lock');
?>
<meta name="csrf-token" content="<?= htmlspecialchars($csrfToken, ENT_QUOTES) ?>">
<div class="kadmin-body">

    <div id="loading" style="display:none"><div class="spinner"></div></div>
    <div id="toast" style="display:none"></div>

    <div id="scale-root">

        <div class="kadmin-stats">
            <div class="kadmin-stat">
                <p class="kadmin-stat__label">💪 Tổng số bậc giới hạn</p>
                <p class="kadmin-stat__value"><?= number_format($totalRows) ?></p>
            </div>
        </div>

        <div class="kadmin-toolbar">
            <div class="kadmin-toolbar__grow"></div>
            <button type="button" class="kame-button kame-button--secondary" onclick="location.reload()">
                <i class="ti ti-refresh" aria-hidden="true"></i> Làm Mới
            </button>
        </div>

        <div class="kadmin-card">
            <div class="kadmin-card__head"><i class="ti ti-shield-lock" aria-hidden="true"></i> Bảng giới hạn chỉ số</div>
            <div class="kadmin-table-wrap">
                <table class="kadmin-table">
                <thead>
                  <tr class="header-row">
                    <th>ID</th>
                    <th>⚡ Power</th>
                    <th>❤️ HP</th>
                    <th>💧 KI</th>
                    <th>⚔️ Damege</th>
                    <th>🛡️ Def</th>
                    <th>🎯 Crit</th>
                    <th>✏️ VIP</th>
                  </tr>
                </thead>
                <tbody>
                <?php foreach ($rows as $row):
                    $id = $row['id'];
                ?>
                <tr data-id="<?= $id ?>">
                  <td class="num"><span class="no-edit view-value"><?= $id ?></span></td>
                  <td class="num"><span class="view-value" data-field="power" data-original="<?= $row['power'] ?>"><?= number_format($row['power'], 0, ',', '.') ?></span></td>
                  <td class="num"><span class="view-value" data-field="hp" data-original="<?= $row['hp'] ?>"><?= number_format($row['hp'], 0, ',', '.') ?></span></td>
                  <td class="num"><span class="view-value" data-field="mp" data-original="<?= $row['mp'] ?>"><?= number_format($row['mp'], 0, ',', '.') ?></span></td>
                  <td class="num"><span class="view-value" data-field="damage" data-original="<?= $row['damage'] ?>"><?= number_format($row['damage'], 0, ',', '.') ?></span></td>
                  <td class="num"><span class="view-value" data-field="defense" data-original="<?= $row['defense'] ?>"><?= number_format($row['defense'], 0, ',', '.') ?></span></td>
                  <td class="num"><span class="view-value" data-field="critical" data-original="<?= $row['critical'] ?>"><?= number_format($row['critical'], 0, ',', '.') ?></span></td>
                  <td class="act">
                    <div class="action-wrap">
                      <button class="btn-sm btn-edit" onclick="editRow(this)">✏️ Sửa</button>
                    </div>
                  </td>
                </tr>
                <?php endforeach; ?>
                <?php if (empty($rows)): ?>
                <tr><td colspan="8" style="padding:32px;text-align:center" class="kame-muted">Không có dữ liệu giới hạn chỉ số.</td></tr>
                <?php endif; ?>
                </tbody>
                </table>
            </div>
        </div>

    </div><!-- /scale-root -->
</div><!-- /kadmin-body -->

<script>
/* ================= TOAST / LOADING ================= */
function showToast(msg, type='success') {
    const t = document.getElementById('toast');
    t.textContent = msg; t.className = type; t.style.display = 'block';
    setTimeout(() => t.style.display = 'none', 2200);
}
function showLoading(v) { document.getElementById('loading').style.display = v ? 'flex' : 'none'; }

/* ================= EDIT MODE LOGIC (giống account.php) ================= */
function editRow(btn) {
    const tr = btn.closest('tr');
    if (tr.classList.contains('editing')) return;

    tr.classList.add('editing');

    // Lưu lại giá trị gốc
    const originalValues = {};
    const editableSpans = tr.querySelectorAll('td .view-value[data-field]');
    editableSpans.forEach(span => {
        const field = span.getAttribute('data-field');
        originalValues[field] = span.getAttribute('data-original') || span.innerText;
    });

    // Ẩn nút Sửa gốc (sẽ thay thế tạm thời)
    const actionDiv = tr.querySelector('.action-wrap');
    const editBtn = actionDiv.querySelector('.btn-edit');
    editBtn.style.display = 'none';

    // Tạo input cho từng trường có thể sửa
    editableSpans.forEach(span => {
        const field = span.getAttribute('data-field');
        const currentValue = span.getAttribute('data-original') || span.innerText;
        const td = span.parentNode;
        let inputEl = document.createElement('input');
        inputEl.type = 'text';
        // Hiển thị giá trị thô (không dấu chấm) để dễ chỉnh sửa
        let rawValue = currentValue.toString().replace(/\./g, '');
        inputEl.value = rawValue;
        inputEl.classList.add('edit-input');
        inputEl.setAttribute('data-field', field);

        // Xử lý sự kiện nhập: chỉ cho phép số
        inputEl.addEventListener('input', function(e) {
            this.value = this.value.replace(/[^\d]/g, '');
        });

        span.style.display = 'none';
        td.appendChild(inputEl);
        inputEl.focus();
    });

    // Thay nút Sửa bằng Lưu + Hủy
    const saveBtn = document.createElement('button');
    saveBtn.className = 'btn-sm btn-save';
    saveBtn.innerHTML = '💾 Lưu';
    const cancelBtn = document.createElement('button');
    cancelBtn.className = 'btn-sm btn-cancel';
    cancelBtn.innerHTML = '✖️ Hủy';

    actionDiv.appendChild(saveBtn);
    actionDiv.appendChild(cancelBtn);

    saveBtn.onclick = () => saveEdit(tr, originalValues);
    cancelBtn.onclick = () => cancelEdit(tr, originalValues);
}

async function saveEdit(tr, originalValues) {
    const id = tr.getAttribute('data-id');
    const inputs = tr.querySelectorAll('.edit-input');
    const data = {
        action: 'save',
        id: id,
        power: '',
        hp: '',
        mp: '',
        damage: '',
        defense: '',
        critical: ''
    };
    inputs.forEach(inp => {
        const field = inp.getAttribute('data-field');
        if (field && data.hasOwnProperty(field)) {
            let value = inp.value.replace(/[^\d]/g, '');
            if (value === '') value = '0';
            data[field] = value;
        }
    });

    const fields = ['power','hp','mp','damage','defense','critical'];
    fields.forEach(f => {
        if (!data[f]) data[f] = originalValues[f] || '0';
    });

    showLoading(true);
    try {
        const res = await fetch(window.location.href, {
            method: 'POST',
            headers: { 'Content-Type': 'application/x-www-form-urlencoded', 'X-CSRF-TOKEN': document.querySelector('meta[name="csrf-token"]')?.content || '' },
            body: new URLSearchParams(data)
        });
        const json = await res.json();
        if (json.success) {
            // Xóa input trước
            inputs.forEach(inp => inp.remove());
            // Hiện lại span và cập nhật giá trị mới
            const spans = tr.querySelectorAll('.view-value[data-field]');
            spans.forEach(span => {
                const field = span.getAttribute('data-field');
                const newVal = data[field];
                if (newVal !== undefined) {
                    const formatted = Number(newVal).toLocaleString('vi-VN');
                    span.innerText = formatted;
                    span.setAttribute('data-original', newVal);
                }
                span.style.display = 'inline-block';
            });
            showToast(`✅ Đã lưu giới hạn chỉ số ID ${id}`);
        } else {
            showToast('❌ Lỗi: ' + (json.error || 'Không xác định'), 'error');
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
    // Xóa input trước, rồi mới hiện span
    const inputs = tr.querySelectorAll('.edit-input');
    inputs.forEach(inp => inp.remove());
    const spans = tr.querySelectorAll('.view-value');
    spans.forEach(span => span.style.display = 'inline-block');
    finishEdit(tr);
}

function finishEdit(tr) {
    tr.classList.remove('editing');
    const actionDiv = tr.querySelector('.action-wrap');
    // Xóa nút Lưu và Hủy nếu có
    const saveBtn = actionDiv.querySelector('.btn-save');
    const cancelBtn = actionDiv.querySelector('.btn-cancel');
    if (saveBtn) saveBtn.remove();
    if (cancelBtn) cancelBtn.remove();
    // Hiện lại nút Sửa
    const editBtn = actionDiv.querySelector('.btn-edit');
    if (editBtn) editBtn.style.display = 'inline-flex';
}
</script>
<?php require __DIR__ . '/guide_widget.php'; ?>
<?php
admin_foot();
