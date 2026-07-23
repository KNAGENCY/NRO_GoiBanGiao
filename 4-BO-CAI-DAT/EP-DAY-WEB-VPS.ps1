# ===========================================================================
#  EP CHEP CAC FILE WEB VUA SUA LEN VPS  (khong dua vao gio)
#
#  DONG-BO dung robocopy /XO (bo qua file cu hon theo GIO). Neu dong ho VPS va
#  may nha lech nhau, no co the tuong file dich moi hon nen KHONG chep file da
#  sua. Script nay chep THANG cac file da sua, ghi de bang noi dung, khong dua
#  vao gio, nen chac chan len duoc.
#
#  KHONG dung toi file mat khau (config.php, boss_status.php).
#
#  Chay TREN VPS:
#     powershell -ExecutionPolicy Bypass -File "\\tsclient\C\Users\Admin\Music\bo-cai-vps\EP-DAY-WEB-VPS.ps1"
# ===========================================================================

$ErrorActionPreference = 'Stop'
function Ghi($m,$c='White'){ Write-Host "  $m" -ForegroundColor $c }

$NHA = '\\tsclient\C\xampp\htdocs'
$VPS = 'C:\xampp\htdocs'

# Danh sach file web da sua can ep len (duong dan tuong doi trong htdocs)
$files = @(
    'sv1\admin\buff_item.php',   # sua loi "khong tai duoc danh sach vat pham"
    'sv1\admin\panel.php',       # them dong ho do thoi gian
    'img\icon.jpg',              # logo Kame moi cho khu admin
    'icon_asset.php'             # SUA duong dan icon cho VPS (C:\zz1) - fix anh 503
)

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  EP CHEP FILE WEB DA SUA LEN VPS' -ForegroundColor Cyan
Write-Host '===========================================================' -ForegroundColor Cyan

if (-not (Test-Path $NHA)) {
    Ghi 'DUNG LAI: khong nhin thay o dia may nha (\\tsclient\C).' 'Red'
    Ghi 'Vao Remote Desktop nho tich chia se o C.' 'Yellow'
    Read-Host "`n  Enter de thoat"; exit 1
}

Write-Host ''
foreach ($f in $files) {
    $src = Join-Path $NHA $f
    $dst = Join-Path $VPS $f
    if (-not (Test-Path $src)) { Ghi ("[bo qua] khong thay nguon: " + $f) 'Yellow'; continue }

    $sizeNha = (Get-Item $src).Length
    $sizeCu  = if (Test-Path $dst) { (Get-Item $dst).Length } else { 0 }

    # Sao luu ban VPS truoc khi de (chi lan dau)
    if ((Test-Path $dst) -and -not (Test-Path "$dst.truoc-epday")) { Copy-Item $dst "$dst.truoc-epday" }

    Copy-Item $src $dst -Force
    $sizeMoi = (Get-Item $dst).Length
    Ghi ("[OK] {0}" -f $f) 'Green'
    Ghi ("     VPS cu: {0:N0} byte  ->  moi: {1:N0} byte  (khop may nha: {2})" -f $sizeCu, $sizeMoi, ($sizeMoi -eq $sizeNha)) 'Gray'
}

Write-Host ''
Ghi '--- Kiem chung nhanh (grep dau hieu ban sua trong buff_item.php) ---' 'Cyan'
$bi = Join-Path $VPS 'sv1\admin\buff_item.php'
if (Test-Path $bi) {
    $co = Select-String -Path $bi -Pattern 'JSON_INVALID_UTF8_SUBSTITUTE' -SimpleMatch
    if ($co) { Ghi '[OK] buff_item.php DA co ban va (JSON_INVALID_UTF8_SUBSTITUTE)' 'Green' }
    else { Ghi '[!!] buff_item.php CHUA co ban va - kiem lai' 'Red' }
}

Write-Host ''
Ghi 'Xong. Mo lai panel + trang tang vat pham, bam Ctrl+F5.' 'Green'
Write-Host ''
Read-Host '  Bam Enter de dong'
