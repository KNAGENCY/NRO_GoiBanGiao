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

$item_per_page = 100;
$current_page  = max(1, (int)($_GET['page'] ?? 1));
$offset        = ($current_page - 1) * $item_per_page;
$total_row          = $pdo->query("SELECT COUNT(*) as total FROM history_transaction")->fetch();
$total_transactions = $total_row['total'];
$total_pages        = max(1, ceil($total_transactions / $item_per_page));

$itemIcons = [];
$tplStmt   = $pdo->query("SELECT name, icon_id FROM item_template");
while ($t = $tplStmt->fetch()) {
    $itemIcons[trim($t['name'])] = $t['icon_id'];
}

$ALLOW_ITEMS = [
    'Thỏi vàng',
    'Ngọc rồng 3 sao',
    'Ngọc rồng 6 sao',
    'Ngọc rồng 7 sao',
    'Đá AZ'
];
define('MAX_INVENTORY_SLOTS', 5);
function renderItems(string $s, array $map, bool $twoCol = false, bool $invGrid = false, array $newItems = []): string {
    global $ALLOW_ITEMS;

    $s = trim($s);
    $isEmpty = ($s === '' || $s === 'Gold: 0, ');

    if ($invGrid) {
        $items = [];
        if (!$isEmpty) {
            foreach (explode(',', $s) as $item) {
                $item = trim($item);
                if ($item === '') continue;
                $cleanName = preg_replace('/ \(x\d+\)/', '', $item);
                $allowed = false;
                foreach ($ALLOW_ITEMS as $ai) {
                    if (mb_stripos($cleanName, $ai) !== false) { $allowed = true; break; }
                }
                if (!$allowed) continue;
                $qty = '';
                if (preg_match('/\(x(\d+)\)/', $item, $qm))
                    $qty = number_format((int)$qm[1], 0, ',', '.');
                $icon = null;
                if (stripos($item, 'Gold') !== false) {
                    $icon = 930;
                } else {
                    if (isset($map[$cleanName])) $icon = $map[$cleanName];
                }
                $fmt      = preg_replace_callback('/\d+/', fn($m) => number_format($m[0], 0, ',', '.'), $item);
                $baseName = preg_replace('/ \(x\d+\)/', '', $fmt);
                $isNew = false;
                foreach ($newItems as $ni) {
                    $niClean = preg_replace('/ \(x\d+\)/', '', trim($ni));
                    if (mb_stripos($cleanName, $niClean) !== false || mb_stripos($niClean, $cleanName) !== false) {
                        $isNew = true; break;
                    }
                }
                $items[] = ['icon'=>$icon,'qty'=>$qty,'name'=>$baseName,'new'=>$isNew];
            }
        }
        $filled = count($items);
        $cols   = 5;
        $SLOTS  = max($filled, MAX_INVENTORY_SLOTS);
        $total  = (int)ceil($SLOTS / $cols) * $cols;
        $out = "<div class='inv-grid'>";
        for ($i = 0; $i < $total; $i++) {
            if ($i < $filled) {
                $it  = $items[$i];
                $cls = 'inv-slot filled' . ($it['new'] ? ' inv-new' : '');
                $out .= "<div class='{$cls}' title='" . htmlspecialchars($it['name']) . "'>";
                if ($it['new']) $out .= "<span class='inv-glow'></span>";
                if ($it['icon'])
                    $out .= "<img src='/icon/{$it['icon']}.png' onerror=\"this.src='/icon/1.png'\">";
                else
                    $out .= "<span class='inv-noicon'>?</span>";
                if ($it['qty'])
                    $out .= "<span class='inv-qty'>{$it['qty']}</span>";
                $out .= "</div>";
            } else {
                $out .= "<div class='inv-slot empty'></div>";
            }
        }
        $out .= "</div>";
        return $out;
    }

    if ($isEmpty) return "<span class='empty-state'>Trống</span>";
    $items = [];
    foreach (explode(',', $s) as $item) {
        $item = trim($item);
        if ($item === '') continue;
        $qty = '';
        if (preg_match('/\(x(\d+)\)/', $item, $qm))
            $qty = number_format((int)$qm[1], 0, ',', '.');
        $icon = null;
        if (stripos($item, 'Gold') !== false) {
            $icon = 930;
        } else {
            $clean = preg_replace('/ \(x\d+\)/', '', $item);
            if (isset($map[$clean])) $icon = $map[$clean];
        }
        $fmt = preg_replace_callback('/\d+/', fn($m) => number_format($m[0], 0, ',', '.'), $item);
        $baseName = preg_replace('/ \(x\d+\)/', '', $fmt);
        $html  = "<div class='ir'>";
        if ($icon) $html .= "<img src='/icon/{$icon}.png' class='ii' onerror=\"this.src='/icon/1.png'\">";
        $html .= "<span class='in'>" . htmlspecialchars($baseName) . "</span>";
        if ($qty) $html .= "<span class='iq'>×{$qty}</span>";
        $html .= "</div>";
        $items[] = $html;
    }
    if (empty($items)) return "<span class='empty-state'>Trống</span>";
    return implode('', $items);
}

$sql = "SELECT t.*, h1.avatar_id AS avt1, h2.avatar_id AS avt2
        FROM history_transaction t
        LEFT JOIN player p1 ON SUBSTRING_INDEX(TRIM(t.player_1),' (',1) = p1.name
        LEFT JOIN head_avatar h1 ON p1.head = h1.head_id
        LEFT JOIN player p2 ON SUBSTRING_INDEX(TRIM(t.player_2),' (',1) = p2.name
        LEFT JOIN head_avatar h2 ON p2.head = h2.head_id
        ORDER BY t.time_tran DESC
        LIMIT $item_per_page OFFSET $offset";
$rows = $pdo->query($sql)->fetchAll();

function pgHTML(int $cur, int $total): string {
    $range = 2; $start = max(1,$cur-$range); $end = min($total,$cur+$range);
    $h = '<div class="pg">';
    if ($cur>1) $h.="<a href='?page=".($cur-1)."' class='pb'>‹</a>";
    if ($start>1){ $h.="<a href='?page=1' class='pb'>1</a>"; if($start>2) $h.="<span class='pd'>…</span>"; }
    for($i=$start;$i<=$end;$i++) $h.="<a href='?page=$i' class='pb".($i==$cur?' on':'')."'>$i</a>";
    if($end<$total){ if($end<$total-1) $h.="<span class='pd'>…</span>"; $h.="<a href='?page=$total' class='pb'>$total</a>"; }
    if($cur<$total) $h.="<a href='?page=".($cur+1)."' class='pb'>›</a>";
    return $h.'</div>';
}
?>
<?php
admin_head('Lịch Sử Giao Dịch');
admin_bar('Lịch Sử Giao Dịch', 'arrows-exchange');
?>
<style>
/* ==========================================================================
   Kame restyle cho cac lop DO PHP/JS sinh ra trong trang giao dich.
   Khong doi ten lop (PHP xuat ra / JS tra cuu theo ten), chi doi dinh dang
   bang bien --kame-*. Khong hardcode mau, khong che them lop moi.
   ========================================================================== */
#app { width: 100%; }

/* --- Phan trang (pgHTML) ------------------------------------------------- */
.pg { display:flex; justify-content:center; align-items:center; flex-wrap:wrap; gap:8px; margin:0 0 var(--kame-space-4); }
.pb { display:inline-flex; align-items:center; justify-content:center; min-width:40px; height:40px; padding:0 12px;
      border:var(--kame-border); border-radius:var(--kame-radius-sm); color:var(--kame-maroon-700);
      background:var(--kame-paper-100); font-weight:800; text-decoration:none; }
.pb:hover:not(.on) { background:var(--kame-paper-300); }
.pb.on { color:var(--kame-white); background:var(--kame-orange-600); }
.pd { color:var(--kame-maroon-700); padding:0 4px; }

/* --- Bang: can giua tieu de, giu bo cuc cot co dinh --------------------- */
.kadmin-table thead th { text-align:center; vertical-align:top; }
.kadmin-table tbody td { vertical-align:top; }
.kadmin-table td.bag-td { padding:5px; }

/* --- Thoi gian ----------------------------------------------------------- */
.ct { font-family:"Roboto Condensed","Arial Narrow",sans-serif; text-align:center; line-height:1.5; }
.ct .d { color:var(--kame-ink); font-weight:600; font-size:12px; }
.ct .t { color:var(--kame-maroon-700); font-weight:900; font-size:13px; }

/* --- Nguoi choi ---------------------------------------------------------- */
.pc { display:flex; flex-direction:column; align-items:center; gap:6px; }
.avt { width:52px; height:52px; object-fit:contain; image-rendering:pixelated;
       border-radius:var(--kame-radius-sm); border:2px solid var(--kame-maroon-700); background:var(--kame-paper-300); }
.avt:hover { border-color:var(--kame-orange-600); }
.pn { font-weight:800; font-size:12px; color:var(--kame-maroon-700); text-align:center; word-break:break-word; line-height:1.3; }

/* --- Hop vat pham -------------------------------------------------------- */
.ib { background:var(--kame-paper-100); border:1px solid var(--kame-maroon-700); border-radius:var(--kame-radius-sm); padding:6px 8px; width:100%; }
.ir { display:flex; align-items:center; gap:6px; padding:3px 0; border-bottom:1px solid rgb(96 0 24 / 22%); flex-wrap:wrap; }
.ir:last-child { border-bottom:none; }
.ii { width:24px; height:24px; flex-shrink:0; image-rendering:pixelated; border-radius:3px; }
.in { flex:1; color:var(--kame-ink); font-size:12px; line-height:1.3; word-break:break-word; }
.iq { flex-shrink:0; font-family:"Roboto Condensed","Arial Narrow",sans-serif; font-weight:800; font-size:11px;
      color:var(--kame-white); background:var(--kame-maroon-700); padding:1px 6px; border-radius:4px; white-space:nowrap; }
.empty-state { color:rgb(40 21 15 / 55%); font-style:italic; font-size:11px; }

/* --- Mui ten trao doi ---------------------------------------------------- */
.arc { text-align:center; vertical-align:middle !important; padding:0 2px !important; }
.arw { display:inline-flex; align-items:center; justify-content:center; width:30px; height:30px; border-radius:50%;
       background:var(--kame-orange-600); border:2px solid var(--kame-maroon-700); color:var(--kame-white); font-size:15px; }

/* --- Hanh trang (tui sau giao dich) ------------------------------------- */
.bp { background:var(--kame-paper-100); border:1px solid var(--kame-maroon-700); border-radius:var(--kame-radius-sm);
      padding:6px; position:relative; overflow:hidden; }
.bp::before { content:''; position:absolute; top:0; left:0; bottom:0; width:3px; }
.bs::before { background:var(--kame-orange-600); }
.br::before { background:var(--kame-green); }

/* --- Luoi o do (inventory) ---------------------------------------------- */
.inv-grid { display:grid; grid-template-columns:repeat(5,1fr); gap:3px; }
.inv-slot { aspect-ratio:1/1; border-radius:var(--kame-radius-sm); position:relative;
            display:flex; align-items:center; justify-content:center; overflow:hidden; }
.inv-slot.empty { background:rgb(40 21 15 / 6%); border:1px solid rgb(40 21 15 / 10%); }
.inv-slot.filled { background:var(--kame-paper-300); border:1px solid var(--kame-maroon-700); }
.inv-slot.filled:hover { border-color:var(--kame-orange-600); }
.inv-slot.filled img { width:60%; height:60%; object-fit:contain; image-rendering:pixelated; }
.inv-slot .inv-noicon { font-size:9px; color:var(--kame-maroon-700); font-weight:700; }
.inv-slot .inv-qty { position:absolute; bottom:2px; right:2px; font-family:"Roboto Condensed","Arial Narrow",sans-serif;
                     font-size:9px; font-weight:800; line-height:1; color:var(--kame-ink);
                     text-shadow:0 1px 2px var(--kame-paper-100); pointer-events:none; z-index:2; }
.inv-slot.inv-new { border-color:var(--kame-gold) !important; box-shadow:0 0 0 2px var(--kame-gold); }
.inv-new::after { content:'\2726'; position:absolute; top:1px; left:2px; font-size:7px; color:var(--kame-orange-600);
                  line-height:1; z-index:3; pointer-events:none; }
.inv-glow { position:absolute; inset:0; pointer-events:none; z-index:0; }

/* --- Trang thai rong ----------------------------------------------------- */
.te { text-align:center; padding:66px 22px; color:rgb(40 21 15 / 55%); }
.te .ei { font-size:48px; margin-bottom:13px; opacity:.5; }

/* --- O tim kiem trong header cot ---------------------------------------- */
.col-search { margin-top:6px; position:relative; font-weight:400; text-transform:none; letter-spacing:0; }
.col-search input { width:100%; padding:5px 6px 5px 26px; background:rgb(255 253 245 / 88%);
                    border:2px solid rgb(96 0 24 / 45%); border-radius:var(--kame-radius-sm); color:var(--kame-ink);
                    font:inherit; font-size:12px; outline:none; }
.col-search input:focus { border-color:var(--kame-orange-600); }
.col-search input::placeholder { color:rgb(40 21 15 / 55%); font-size:11px; }
.col-search .sic { position:absolute; left:7px; top:50%; transform:translateY(-50%);
                   color:var(--kame-maroon-700); pointer-events:none; line-height:0; }
.col-search .sic svg { width:13px; height:13px; }

/* --- Dong dem ket qua tim kiem ------------------------------------------ */
.search-cnt-tr td { padding:4px 10px; background:var(--kame-paper-300); }
.search-cnt-row { text-align:center; font-size:12px; color:var(--kame-maroon-700); font-weight:700; }

/* --- An dong khi loc ----------------------------------------------------- */
tbody tr.hide { display:none; }
</style>

<div class="kadmin-body">
<div id="app">

  <div class="kadmin-stats">
    <div class="kadmin-stat">
      <div class="kadmin-stat__label">Tổng giao dịch</div>
      <div class="kadmin-stat__value"><?= number_format($total_transactions) ?></div>
    </div>
  </div>

  <div class="kadmin-toolbar">
    <a href="/sv1/admin/panel" class="kame-button kame-button--secondary"><i class="ti ti-home" aria-hidden="true"></i> Panel</a>
    <button type="button" class="kame-button" onclick="location.reload()"><i class="ti ti-refresh" aria-hidden="true"></i> Làm mới</button>
  </div>

  <!-- PAGINATION -->
  <?= pgHTML($current_page, $total_pages) ?>

  <!-- TABLE -->
  <div class="kadmin-table-wrap">
  <table class="kadmin-table" style="table-layout:fixed;min-width:1465px">
    <colgroup>
      <col style="width:100px">   <!-- TIME          -->
      <col style="width:200px">   <!-- Hành Trang 1  -->
      <col style="width:150px">   <!-- Người Gửi     -->
      <col style="width:300px">   <!-- ITEM Gửi      -->
      <col style="width:65px">    <!-- Arrow         -->
      <col style="width:300px">   <!-- ITEM Nhận     -->
      <col style="width:150px">   <!-- Người Nhận    -->
      <col style="width:200px">   <!-- Hành Trang 2  -->
    </colgroup>
    <thead>
      <tr>
        <th>TIME</th>
        <th>Hành Trang</th>
        <th>
          Người Gửi
          <div class="col-search">
           <span class="sic">
    <svg viewBox="0 0 24 24" width="14" height="14" fill="none">
      <path d="M12 12a5 5 0 1 0 0-10 5 5 0 0 0 0 10Z"
            stroke="currentColor" stroke-width="2"/>
      <path d="M3 21c0-4.5 4-7 9-7s9 2.5 9 7"
            stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
    </svg>
  </span>
            <input type="text" id="s1" placeholder="Tìm người gửi…">
          </div>
        </th>
        <th>ITEM Gửi</th>
        <th></th>
        <th>ITEM Nhận</th>
        <th>
          Người Nhận
          <div class="col-search">
        <span class="sic">
    <svg viewBox="0 0 24 24" width="14" height="14" fill="none">
      <path d="M12 12a5 5 0 1 0 0-10 5 5 0 0 0 0 10Z"
            stroke="currentColor" stroke-width="2"/>
      <path d="M3 21c0-4.5 4-7 9-7s9 2.5 9 7"
            stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
    </svg>
  </span>
            <input type="text" id="s2" placeholder="Tìm người nhận…">
          </div>
        </th>
        <th>Hành Trang</th>
      </tr>
      <tr class="search-cnt-tr" id="search-cnt-tr" style="display:none">
        <td colspan="8"><div class="search-cnt-row"><span id="sc"></span></div></td>
      </tr>
    </thead>
    <tbody id="tb">
    <?php if(empty($rows)): ?>
      <tr><td colspan="8"><div class="te"><div class="ei">📭</div><p>Chưa có giao dịch nào</p></div></td></tr>
    <?php else: foreach($rows as $i => $row):
      $a1 = $row['avt1'] ?? 1; $a2 = $row['avt2'] ?? 1;
      $n1 = explode(" (", $row['player_1'])[0];
      $n2 = explode(" (", $row['player_2'])[0];
      [$date,$time] = array_pad(explode(" ", $row['time_tran'], 2), 2, '');
    ?>
    <tr style="animation-delay:<?= min($i*6,300) ?>ms"
        data-p1="<?= mb_strtolower(htmlspecialchars($n1)) ?>"
        data-p2="<?= mb_strtolower(htmlspecialchars($n2)) ?>">

      <!-- TIME -->
      <td><div class="ct"><div class="d"><?= htmlspecialchars($date) ?></div><div class="t"><?= htmlspecialchars($time) ?></div></div></td>

      <!-- BAG SENDER -->
      <td class="bag-td">
        <div class="bp bs">
          <?php
            $newFor1 = array_filter(array_map('trim', explode(',', $row['item_player_2'])));
            echo renderItems($row['bag_1_after_tran'], $itemIcons, false, true, $newFor1);
          ?>
        </div>
      </td>

      <!-- SENDER -->
      <td>
        <div class="pc">
          <img src="/icon/<?= $a1 ?>.png" class="avt" onerror="this.src='/icon/1.png'">
          <span class="pn"><?= htmlspecialchars($n1) ?></span>
        </div>
      </td>

      <!-- ITEMS SENT -->
      <td><div class="ib"><?= renderItems($row['item_player_1'], $itemIcons) ?></div></td>

      <!-- ARROW -->
      <td class="arc"><div class="arw">⇄</div></td>

      <!-- ITEMS RECEIVED -->
      <td><div class="ib"><?= renderItems($row['item_player_2'], $itemIcons) ?></div></td>

      <!-- RECEIVER -->
      <td>
        <div class="pc">
          <img src="/icon/<?= $a2 ?>.png" class="avt" onerror="this.src='/icon/1.png'">
          <span class="pn"><?= htmlspecialchars($n2) ?></span>
        </div>
      </td>

      <!-- BAG RECEIVER -->
      <td class="bag-td">
        <div class="bp br">
          <?php
            $newFor2 = array_filter(array_map('trim', explode(',', $row['item_player_1'])));
            echo renderItems($row['bag_2_after_tran'], $itemIcons, false, true, $newFor2);
          ?>
        </div>
      </td>

     </tr>
    <?php endforeach; endif; ?>
    </tbody>
  </table>
  </div>

</div><!-- /#app -->
</div><!-- /.kadmin-body -->

<script>
/* ── TÌM KIẾM ────────────────────────────────── */
const s1 = document.getElementById('s1'),
      s2 = document.getElementById('s2'),
      sc = document.getElementById('sc'),
      cntRow = document.getElementById('search-cnt-tr');
const trs = document.querySelectorAll('#tb tr[data-p1]');

function flt(){
  const v1 = s1.value.toLowerCase().trim(),
        v2 = s2.value.toLowerCase().trim();
  let n = 0;
  trs.forEach(r => {
    const ok = r.dataset.p1.includes(v1) && r.dataset.p2.includes(v2);
    r.classList.toggle('hide', !ok);
    if(ok) n++;
  });
  if(v1 || v2){
    sc.textContent = `${n} / ${trs.length} kết quả`;
    cntRow.style.display = '';
  } else {
    cntRow.style.display = 'none';
  }
}
s1.addEventListener('input', flt);
s2.addEventListener('input', flt);
</script>
<?php require __DIR__ . '/guide_widget.php'; ?>
<?php admin_foot(); ?>
