# ===========================================================================
#  VA LOI: APACHE KHONG CHAY SAU KHI CHAY BUOC 3
#
#  CHUYEN GI DA XAY RA
#  ---------------------------------------------------------------------
#  File 3-CAI-GAME-LEN-VPS.ps1 (ban cu) co mot loi khi doc file mau
#  cau-hinh\nrokame-vhost.conf. No di tim hai cai moc:
#        === BAT DAU CHEP ===   va   === HET CHEP ===
#  nhung lai vo phai DONG HUONG DAN o dau file mau - dong do co nhac ten
#  ca hai cai moc trong ngoac kep. Ket qua: no cat ra dung hai chu
#        " va "
#  roi ghim hai chu do vao cuoi file cau hinh cua Apache.
#
#  Hau qua KEP:
#    1. Apache doc den dong  " va "  -> khong hieu -> KHONG chiu chay
#    2. Cau hinh ten mien nrokame.com thi CHUA HE duoc them vao
#
#  FILE NAY LAM GI
#  ---------------------------------------------------------------------
#    1. Tra file cau hinh Apache ve ban goc (da luu san khi chay buoc 3)
#    2. Them dung khoi cau hinh ten mien
#    3. Thu cu phap truoc khi bat - neu sai thi bao ro sai dong nao
#    4. Bat Apache va kiem tra that
#
#  Chay lai nhieu lan khong sao.
# ===========================================================================

$ErrorActionPreference = 'Stop'
function Ghi($m, $mau = 'White') { Write-Host "  $m" -ForegroundColor $mau }

$APACHE    = 'C:\xampp\apache'
$fileVhost = "$APACHE\conf\extra\httpd-vhosts.conf"
$fileGoc   = "$fileVhost.goc"
$httpdExe  = "$APACHE\bin\httpd.exe"
$mauVhost  = Join-Path $PSScriptRoot 'cau-hinh\nrokame-vhost.conf'
$loi = @()

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  VA LOI APACHE KHONG CHAY' -ForegroundColor Cyan
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host ''

# --- Quyen quan tri ------------------------------------------------------
$quyen = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()
         ).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $quyen) {
    Ghi 'DUNG LAI: phai chay bang quyen Administrator.' 'Red'
    Ghi 'Bam chuot phai vao file -> Run with PowerShell (voi quyen quan tri).' 'Yellow'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}
Ghi '[OK] Dang chay bang quyen Administrator' 'Green'

foreach ($f in @($fileVhost, $httpdExe, $mauVhost)) {
    if (-not (Test-Path $f)) {
        Ghi "DUNG LAI: khong thay $f" 'Red'
        Read-Host "`n  Bam Enter de thoat"; exit 1
    }
}

# --- 1. Cho xem cho hong -------------------------------------------------
Write-Host ''
Ghi '--- BUOC 1/5: Tim dong rac trong cau hinh ---' 'Cyan'

$dong = Get-Content $fileVhost
$rac  = $dong | Where-Object { $_.Trim() -eq '" va "' }
if ($rac) {
    Ghi 'Tim thay dong rac (day chinh la thu lam Apache chet):' 'Red'
    $rac | ForEach-Object { Ghi "    [$_]" 'Red' }
} else {
    Ghi '[BO QUA] Khong thay dong rac  " va "  - co the da duoc don roi.' 'Yellow'
}

# --- 2. Tra ve ban goc ---------------------------------------------------
Write-Host ''
Ghi '--- BUOC 2/5: Tra cau hinh ve ban goc ---' 'Cyan'

if (Test-Path $fileGoc) {
    # Ban .goc duoc tao MOT LAN DUY NHAT truoc khi buoc 3 dung den file,
    # nen no chac chan la ban sach.
    Copy-Item $fileGoc $fileVhost -Force
    Ghi '[OK] Da tra httpd-vhosts.conf ve ban goc' 'Green'
} else {
    # Khong co ban goc thi don tay: bo dong rac di.
    Ghi 'Khong thay ban goc - don tay thay vi tra ve.' 'Yellow'
    $sach = $dong | Where-Object { $_.Trim() -ne '" va "' }
    Set-Content -Path $fileVhost -Value $sach -Encoding UTF8
    Ghi '[OK] Da bo dong rac' 'Green'
}

# --- 3. Them dung khoi cau hinh ten mien ---------------------------------
Write-Host ''
Ghi '--- BUOC 3/5: Them cau hinh ten mien nrokame.com ---' 'Cyan'

$dangCo = Get-Content $fileVhost -Raw
if ($dangCo -match 'nrokame\.com') {
    Ghi '[BO QUA] Cau hinh ten mien da co san' 'Yellow'
} else {
    $mau = Get-Content $mauVhost -Raw
    # Bat buoc khop CA DONG: phai co dau # dau dong va day dau = hai ben.
    # Day chinh la cho ban cu lam sai.
    $moc = '(?sm)^#\s*=+\s*BAT DAU CHEP\s*=+\s*\r?$(.*?)^#\s*=+\s*HET CHEP\s*=+\s*\r?$'
    if ($mau -notmatch $moc) {
        Ghi 'DUNG LAI: khong doc duoc khoi cau hinh trong nrokame-vhost.conf' 'Red'
        Read-Host "`n  Bam Enter de thoat"; exit 1
    }
    $phan = $Matches[1].Trim()

    Ghi 'Doan sap them vao:' 'Gray'
    $phan -split "`r?`n" | ForEach-Object { Ghi "    $_" 'DarkGray' }

    Add-Content -Path $fileVhost -Value "`r`n$phan`r`n" -Encoding UTF8
    Ghi '[OK] Da them cau hinh ten mien' 'Green'
}

# --- 4. Thu cu phap TRUOC khi bat ----------------------------------------
Write-Host ''
Ghi '--- BUOC 4/5: Thu cu phap cau hinh ---' 'Cyan'

# httpd -t in ket qua ra luong LOI chu khong phai luong thuong, nen phai
# gop 2>&1. Va phai tam ha $ErrorActionPreference, neu khong PowerShell coi
# moi dong o luong loi la loi that va nem ngoai le ngay.
$cu = $ErrorActionPreference
$ErrorActionPreference = 'Continue'
try { $kq = & $httpdExe -t 2>&1 } finally { $ErrorActionPreference = $cu }

if ($kq -match 'Syntax OK') {
    Ghi '[OK] Cau hinh dung cu phap' 'Green'
} else {
    Ghi 'VAN CON LOI CU PHAP - khong bat Apache:' 'Red'
    $kq | ForEach-Object { Ghi "    $_" 'Red' }
    Ghi ''
    Ghi 'Chup man hinh nay gui lai cho tro ly.' 'Yellow'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}

# --- 5. Bat Apache va kiem tra that --------------------------------------
Write-Host ''
Ghi '--- BUOC 5/5: Bat Apache va kiem tra ---' 'Cyan'

$dv = Get-Service -Name 'Apache2.4' -EA SilentlyContinue
if (-not $dv) {
    Ghi 'DUNG LAI: khong co dich vu Apache2.4. Chay CAI-DICH-VU.ps1 truoc.' 'Red'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}

if ($dv.Status -ne 'Running') {
    try {
        Start-Service 'Apache2.4'
        Start-Sleep -Seconds 4
    } catch {
        Ghi "LOI khi bat Apache: $($_.Exception.Message)" 'Red'
        $loi += 'Khong bat duoc Apache'
    }
} else {
    Restart-Service 'Apache2.4' -Force -EA SilentlyContinue
    Start-Sleep -Seconds 4
}

$tt = (Get-Service 'Apache2.4').Status
if ($tt -eq 'Running') { Ghi "[OK] Apache2.4 : $tt" 'Green' }
else { Ghi "Apache2.4 : $tt" 'Red'; $loi += 'Apache khong chay' }

# Khong tin trang thai dich vu - thu ket noi that.
$c = Test-NetConnection -ComputerName 127.0.0.1 -Port 80 -WarningAction SilentlyContinue
if ($c.TcpTestSucceeded) { Ghi '[OK] Cong 80 nghe duoc' 'Green' }
else { Ghi 'Cong 80 KHONG nghe' 'Red'; $loi += 'Cong 80 khong nghe' }

# Thu lay that mot trang, va thu ca duong dan sach (/sv1/trangchu) vi do la
# thu chung minh Apache co doc file .htaccess hay khong.
foreach ($duong in @('http://127.0.0.1/', 'http://127.0.0.1/sv1/trangchu')) {
    try {
        $r = Invoke-WebRequest -Uri $duong -UseBasicParsing -TimeoutSec 10
        Ghi "[OK] $duong -> ma $($r.StatusCode)" 'Green'
    } catch {
        Ghi "$duong -> khong vao duoc" 'Red'
        $loi += "Khong vao duoc $duong"
    }
}

# Xem Apache da nhan ten mien chua
Write-Host ''
Ghi 'Danh sach ten mien Apache dang phuc vu:' 'Cyan'
$cu = $ErrorActionPreference
$ErrorActionPreference = 'Continue'
try { $ds = & $httpdExe -S 2>&1 } finally { $ErrorActionPreference = $cu }
$ds | Where-Object { $_ -match 'nrokame|port 80|namevhost|default server' } |
      ForEach-Object { Ghi "    $_" 'Gray' }

# --- Tong ket ------------------------------------------------------------
Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
if ($loi.Count -eq 0) {
    Write-Host '  XONG - APACHE DA CHAY VA WEB TRA LOI' -ForegroundColor Green
    Write-Host '===========================================================' -ForegroundColor Cyan
    Write-Host ''
    Ghi 'Viec tiep theo:' 'Yellow'
    Ghi '  1. Mo C:\zz1 -> chuot phai run.bat -> Run as administrator'
    Ghi '  2. Doi cua so may chu game hien len (30-60 giay)'
    Ghi '  3. Lay DIEN THOAI, TAT WIFI (dung 4G), vao http://nrokame.com'
} else {
    Write-Host "  CON $($loi.Count) VAN DE" -ForegroundColor Red
    Write-Host '===========================================================' -ForegroundColor Cyan
    $loi | ForEach-Object { Ghi "  - $_" 'Red' }
    Write-Host ''
    Ghi 'Chup man hinh nay gui lai cho tro ly.' 'Yellow'
}
Write-Host ''
Read-Host '  Bam Enter de dong'
