# ===========================================================================
#  BUNG GOI DU LIEU TREN VPS  (ban 2 - chia nho file lon)
#
#  Vi sao chia nho: chep nguyen file 900 MB qua o dia chia se cua Remote
#  Desktop bi hong giua chung (ERROR 995), ma robocopy van bao thanh cong.
#  Chia thanh 9 phan 100 MB thi de chep hon, va neu phan nao hong thi chi
#  phai chep lai dung phan do chu khong phai ca 900 MB.
#
#  Script kiem tra TUNG PHAN truoc khi ghep, roi kiem tra lai file sau khi
#  ghep. Chay lai nhieu lan khong sao.
# ===========================================================================

$ErrorActionPreference = 'Stop'
function Ghi($m, $mau = 'White') { Write-Host "  $m" -ForegroundColor $mau }

$GOI = 'C:\kame-goi'
$loi = @()

# --- Ma kiem tra lay tu may nha -----------------------------------------
$PHAN = [ordered]@{
    'zz1.tar.part-00' = @{ Byte = 104857600; Ma = '8A65F13EE2AEAE46' }
    'zz1.tar.part-01' = @{ Byte = 104857600; Ma = 'AF9D210DF309460A' }
    'zz1.tar.part-02' = @{ Byte = 104857600; Ma = 'C1DED9828B0F78F7' }
    'zz1.tar.part-03' = @{ Byte = 104857600; Ma = 'B2DF40509D541097' }
    'zz1.tar.part-04' = @{ Byte = 104857600; Ma = '21862A8A423503BE' }
    'zz1.tar.part-05' = @{ Byte = 104857600; Ma = '4B46E33F8DBA6B43' }
    'zz1.tar.part-06' = @{ Byte = 104857600; Ma = '72BF1D2E27AF3CE1' }
    'zz1.tar.part-07' = @{ Byte = 104857600; Ma = '5ED838FDCDA0E844' }
    'zz1.tar.part-08' = @{ Byte = 103198720; Ma = '277915B4D0125FFE' }
}
$KHAC = @{
    'htdocs.tar' = @{ Byte = 21534720; Ma = 'DA9797C4839C1C86' }
    'z.sql'      = @{ Byte = 856302;   Ma = '1E37C67A6C6382D5' }
}
$MA_TAR_DAY_DU = '72A555DF08D6A58A'   # ma cua zz1.tar sau khi ghep lai

function KiemTraFile($ten, $chuan) {
    $f = Join-Path $GOI $ten
    if (-not (Test-Path $f)) { return "THIEU" }
    $co = Get-Item $f
    if ($co.Length -ne $chuan.Byte) { return "SAI DUNG LUONG (co $($co.Length), phai $($chuan.Byte))" }
    $ma = (Get-FileHash $f -Algorithm SHA256).Hash.Substring(0, 16)
    if ($ma -ne $chuan.Ma) { return "HONG (ma $ma, phai $($chuan.Ma))" }
    return $null
}

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  BUNG GOI DU LIEU - NGOC RONG KAME' -ForegroundColor Cyan
Write-Host '===========================================================' -ForegroundColor Cyan

# --- 1. Kiem tra 9 phan + 2 file khac -----------------------------------
Write-Host ''
Ghi '--- BUOC 1/5: Kiem tra 11 file co nguyen ven khong ---' 'Cyan'
Ghi 'Dang tinh ma kiem tra, doi khoang 30 giay...' 'Gray'

foreach ($ten in $KHAC.Keys) {
    $kq = KiemTraFile $ten $KHAC[$ten]
    if ($kq) { Ghi "$ten : $kq" 'Red'; $loi += $ten } else { Ghi "[OK] $ten" 'Green' }
}
foreach ($ten in $PHAN.Keys) {
    $kq = KiemTraFile $ten $PHAN[$ten]
    if ($kq) { Ghi "$ten : $kq" 'Red'; $loi += $ten } else { Ghi "[OK] $ten" 'Green' }
}

if ($loi.Count -gt 0) {
    Write-Host ''
    Ghi "DUNG LAI: $($loi.Count) file hong hoac thieu:" 'Red'
    $loi | ForEach-Object { Ghi "   - $_" 'Red' }
    Write-Host ''
    Ghi 'CHEP LAI DUNG NHUNG FILE DO (khong phai ca goi). Chay tren VPS:' 'Yellow'
    Write-Host ''
    foreach ($t in $loi) {
        Ghi "  Remove-Item 'C:\kame-goi\$t' -Force -EA 0" 'Gray'
    }
    $ds = ($loi | ForEach-Object { "`"$_`"" }) -join ' '
    Ghi "  robocopy `"\\tsclient\C\Users\Admin\Music\goi-len-vps`" `"C:\kame-goi`" $ds /R:5 /W:5" 'Gray'
    Write-Host ''
    Ghi 'Roi chay lai script nay.' 'Yellow'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}

# --- 2. Ghep 9 phan thanh zz1.tar ---------------------------------------
Write-Host ''
Ghi '--- BUOC 2/5: Ghep 9 phan thanh zz1.tar ---' 'Cyan'

$fileTar = Join-Path $GOI 'zz1.tar'
if (Test-Path $fileTar) { Remove-Item $fileTar -Force }

$ra = [IO.File]::Create($fileTar)
try {
    foreach ($ten in $PHAN.Keys) {
        $vao = [IO.File]::OpenRead((Join-Path $GOI $ten))
        try { $vao.CopyTo($ra) } finally { $vao.Dispose() }
        Ghi "  da ghep $ten" 'Gray'
    }
} finally { $ra.Dispose() }

$maGhep = (Get-FileHash $fileTar -Algorithm SHA256).Hash.Substring(0, 16)
if ($maGhep -ne $MA_TAR_DAY_DU) {
    Ghi "LOI: file sau khi ghep khong dung (ma $maGhep, phai $MA_TAR_DAY_DU)" 'Red'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}
Ghi "[OK] Da ghep xong va kiem tra dung - $([math]::Round((Get-Item $fileTar).Length/1MB)) MB" 'Green'

# --- 3. Bung web ---------------------------------------------------------
Write-Host ''
Ghi '--- BUOC 3/5: Bung web vao C:\xampp\htdocs ---' 'Cyan'

if ((Test-Path 'C:\xampp\htdocs') -and -not (Test-Path 'C:\xampp\htdocs\sv1')) {
    Rename-Item 'C:\xampp\htdocs' 'htdocs.xampp-goc' -Force
    Ghi '[OK] Da doi ten htdocs mac dinh cua XAMPP thanh htdocs.xampp-goc (khong xoa)' 'Green'
}
tar -xf (Join-Path $GOI 'htdocs.tar') -C 'C:\xampp'
if ($LASTEXITCODE -ne 0) { Ghi "LOI: tar bao ma $LASTEXITCODE khi bung web" 'Red'; $loi += 'bung-web' }
else { Ghi '[OK] Da bung web' 'Green' }

# --- 4. Bung server game -------------------------------------------------
Write-Host ''
Ghi '--- BUOC 4/5: Bung server game vao C:\zz1 ---' 'Cyan'
Ghi 'Goi nay co 71.296 file, mat 2-5 phut. Man hinh se dung yen - dung tat.' 'Gray'

New-Item -ItemType Directory -Path 'C:\zz1' -Force | Out-Null
tar -xf $fileTar -C 'C:\zz1'
if ($LASTEXITCODE -ne 0) { Ghi "LOI: tar bao ma $LASTEXITCODE khi bung game" 'Red'; $loi += 'bung-game' }
else { Ghi '[OK] Da bung server game' 'Green' }

# --- 5. Dem lai va kiem tra file quan trong ------------------------------
Write-Host ''
Ghi '--- BUOC 5/5: Dem lai va kiem tra ---' 'Cyan'

$soGame = (Get-ChildItem 'C:\zz1' -Recurse -File -ErrorAction SilentlyContinue).Count
$soWeb  = (Get-ChildItem 'C:\xampp\htdocs' -Recurse -File -ErrorAction SilentlyContinue).Count
Ghi "So file server game : $soGame  (phai khoang 71.900)" $(if ($soGame -gt 71000) { 'Green' } else { 'Red' })
Ghi "So file web         : $soWeb  (phai khoang 115)"     $(if ($soWeb  -gt 100)   { 'Green' } else { 'Red' })
if ($soGame -lt 71000) { $loi += 'thieu-file-game' }
if ($soWeb  -lt 100)   { $loi += 'thieu-file-web' }

$phaiCo = @(
    'C:\zz1\zz1.jar',
    'C:\zz1\run.bat',
    'C:\zz1\lib\ga.jar',
    'C:\zz1\src\server\icon.png',
    'C:\zz1\data\config\config.properties',
    'C:\xampp\htdocs\sv1\trangchu.php',
    'C:\xampp\htdocs\sv1\admin\game_bridge.php',
    'C:\xampp\htdocs\icon_asset.php'
)
Write-Host ''
foreach ($p in $phaiCo) {
    if (Test-Path $p) { Ghi "[OK] $p" 'Green' } else { Ghi "THIEU: $p" 'Red'; $loi += $p }
}

# --- Tong ket ------------------------------------------------------------
Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
if ($loi.Count -eq 0) {
    Write-Host '  XONG - BUNG GOI THANH CONG, KHONG THIEU GI' -ForegroundColor Green
    Write-Host '===========================================================' -ForegroundColor Cyan
    Write-Host ''
    Ghi 'Co the xoa cac file part de lay lai 900 MB o dia:' 'Gray'
    Ghi '  Remove-Item C:\kame-goi\zz1.tar.part-* -Force' 'Gray'
    Write-Host ''
    Ghi 'Buoc tiep theo: chay 3-CAI-GAME-LEN-VPS.ps1' 'Yellow'
} else {
    Write-Host "  CO $($loi.Count) VAN DE - XEM DONG MAU DO O TREN" -ForegroundColor Red
    Write-Host '===========================================================' -ForegroundColor Cyan
    Write-Host ''
    Ghi 'DUNG chay buoc tiep theo. Bao lai kem anh chup man hinh nay.' 'Yellow'
}
Write-Host ''
Read-Host '  Bam Enter de dong'
