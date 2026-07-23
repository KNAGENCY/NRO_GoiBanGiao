# ===========================================================================
#  CAI MYSQL VA APACHE THANH DICH VU WINDOWS
#
#  Vi sao can: XAMPP cai xong chi de hai thu nay chay dang chuong trinh
#  thuong - phai mo XAMPP Control Panel bam nut moi len. Tren may chu that
#  thi khong duoc, vi:
#    - May khoi dong lai la web va CSDL khong tu len
#    - Dang xuat khoi Remote Desktop la chung tat theo
#
#  Bien thanh dich vu Windows thi chung chay ngam, tu len cung may.
#
#  Chay lai nhieu lan khong sao.
# ===========================================================================

$ErrorActionPreference = 'Stop'
function Ghi($m, $mau = 'White') { Write-Host "  $m" -ForegroundColor $mau }

$MYSQLD = 'C:\xampp\mysql\bin\mysqld.exe'
$MYINI  = 'C:\xampp\mysql\bin\my.ini'
$HTTPD  = 'C:\xampp\apache\bin\httpd.exe'
$TEN_MYSQL  = 'mysql'
$TEN_APACHE = 'Apache2.4'
$loi = @()

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  CAI MYSQL + APACHE THANH DICH VU WINDOWS' -ForegroundColor Cyan
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host ''

# --- Kiem tra quyen ------------------------------------------------------
$quyen = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()
         ).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $quyen) {
    Ghi 'DUNG LAI: phai chay bang quyen Administrator.' 'Red'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}
Ghi '[OK] Dang chay bang quyen Administrator' 'Green'

foreach ($f in @($MYSQLD, $MYINI, $HTTPD)) {
    if (-not (Test-Path $f)) {
        Ghi "DUNG LAI: khong thay $f" 'Red'
        Ghi 'XAMPP chua cai xong? Chay lai 1-CAI-MOI-TRUONG.bat.' 'Yellow'
        Read-Host "`n  Bam Enter de thoat"; exit 1
    }
}

# --- Tat cac tien trinh dang chay kieu thuong ---------------------------
# Neu mysqld/httpd dang chay dang chuong trinh thuong thi khong cai dich vu
# duoc (file bi khoa, hoac cong 3306/80 da co nguoi giu).
Write-Host ''
Ghi '--- Tat cac tien trinh dang chay kieu thuong ---' 'Cyan'
foreach ($tt in @('mysqld', 'httpd')) {
    $ds = Get-Process -Name $tt -ErrorAction SilentlyContinue
    if ($ds) {
        $ds | Stop-Process -Force -ErrorAction SilentlyContinue
        Start-Sleep -Seconds 2
        Ghi "[OK] Da tat $($ds.Count) tien trinh $tt" 'Green'
    } else {
        Ghi "[BO QUA] $tt khong chay san" 'Gray'
    }
}

# --- 1. MySQL ------------------------------------------------------------
Write-Host ''
Ghi '--- BUOC 1/2: MySQL ---' 'Cyan'

$dv = Get-Service -Name $TEN_MYSQL -ErrorAction SilentlyContinue
if ($dv) {
    Ghi "[BO QUA] Dich vu '$TEN_MYSQL' da co san" 'Yellow'
} else {
    # --install <ten> --defaults-file=<duong dan my.ini>
    # Phai co --defaults-file, neu khong MySQL se khong biet thu muc du lieu
    # nam o dau va dich vu se bat len roi tat ngay.
    & $MYSQLD --install $TEN_MYSQL --defaults-file=$MYINI
    Start-Sleep -Seconds 2
    $dv = Get-Service -Name $TEN_MYSQL -ErrorAction SilentlyContinue
    if ($dv) { Ghi "[OK] Da cai dich vu '$TEN_MYSQL'" 'Green' }
    else { Ghi "LOI: khong cai duoc dich vu '$TEN_MYSQL'" 'Red'; $loi += 'cai-mysql' }
}

if ($dv) {
    Set-Service -Name $TEN_MYSQL -StartupType Automatic
    if ($dv.Status -ne 'Running') {
        try {
            Start-Service -Name $TEN_MYSQL
            Start-Sleep -Seconds 4
            Ghi '[OK] Da bat MySQL' 'Green'
        } catch {
            Ghi "LOI khi bat MySQL: $($_.Exception.Message)" 'Red'
            Ghi 'Xem nhat ky loi o: C:\xampp\mysql\data\*.err' 'Yellow'
            $loi += 'bat-mysql'
        }
    } else { Ghi '[OK] MySQL dang chay' 'Green' }
}

# --- 2. Apache -----------------------------------------------------------
Write-Host ''
Ghi '--- BUOC 2/2: Apache ---' 'Cyan'

$dvA = Get-Service -Name $TEN_APACHE -ErrorAction SilentlyContinue
if ($dvA) {
    Ghi "[BO QUA] Dich vu '$TEN_APACHE' da co san" 'Yellow'
} else {
    & $HTTPD -k install -n $TEN_APACHE
    Start-Sleep -Seconds 2
    $dvA = Get-Service -Name $TEN_APACHE -ErrorAction SilentlyContinue
    if ($dvA) { Ghi "[OK] Da cai dich vu '$TEN_APACHE'" 'Green' }
    else { Ghi "LOI: khong cai duoc dich vu '$TEN_APACHE'" 'Red'; $loi += 'cai-apache' }
}

if ($dvA) {
    Set-Service -Name $TEN_APACHE -StartupType Automatic
    if ($dvA.Status -ne 'Running') {
        try {
            Start-Service -Name $TEN_APACHE
            Start-Sleep -Seconds 3
            Ghi '[OK] Da bat Apache' 'Green'
        } catch {
            Ghi "LOI khi bat Apache: $($_.Exception.Message)" 'Red'
            Ghi 'Kiem tra cau hinh: C:\xampp\apache\bin\httpd.exe -t' 'Yellow'
            $loi += 'bat-apache'
        }
    } else { Ghi '[OK] Apache dang chay' 'Green' }
}

# --- Kiem chung that -----------------------------------------------------
Write-Host ''
Ghi '--- Kiem chung ---' 'Cyan'

foreach ($t in @($TEN_MYSQL, $TEN_APACHE)) {
    $s = Get-Service -Name $t -ErrorAction SilentlyContinue
    if ($s) {
        $mau = if ($s.Status -eq 'Running') { 'Green' } else { 'Red' }
        Ghi "$t : $($s.Status) | tu khoi dong: $($s.StartType)" $mau
        if ($s.Status -ne 'Running') { $loi += "$t-khong-chay" }
    } else { Ghi "$t : KHONG CO" 'Red'; $loi += "$t-khong-co" }
}

# Thu ket noi that, khong tin trang thai dich vu
$c = Test-NetConnection -ComputerName 127.0.0.1 -Port 3306 -WarningAction SilentlyContinue
Ghi "Cong 3306 (CSDL)  : $(if ($c.TcpTestSucceeded) {'nghe duoc'} else {'KHONG nghe'})" $(if ($c.TcpTestSucceeded) {'Green'} else {'Red'})
if (-not $c.TcpTestSucceeded) { $loi += 'cong-3306' }

$c = Test-NetConnection -ComputerName 127.0.0.1 -Port 80 -WarningAction SilentlyContinue
Ghi "Cong 80 (web)     : $(if ($c.TcpTestSucceeded) {'nghe duoc'} else {'KHONG nghe'})" $(if ($c.TcpTestSucceeded) {'Green'} else {'Red'})
if (-not $c.TcpTestSucceeded) { $loi += 'cong-80' }

# --- Tong ket ------------------------------------------------------------
Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
if ($loi.Count -eq 0) {
    Write-Host '  XONG - CA HAI DICH VU DA CHAY VA TU LEN CUNG MAY' -ForegroundColor Green
    Write-Host '===========================================================' -ForegroundColor Cyan
    Write-Host ''
    Ghi 'Buoc tiep theo: chay 3-CAI-GAME-LEN-VPS.ps1' 'Yellow'
} else {
    Write-Host "  CO $($loi.Count) VAN DE" -ForegroundColor Red
    Write-Host '===========================================================' -ForegroundColor Cyan
    $loi | ForEach-Object { Ghi "  - $_" 'Red' }
    Write-Host ''
    Ghi 'Chup man hinh nay gui lai cho tro ly.' 'Yellow'
}
Write-Host ''
Read-Host '  Bam Enter de dong'
