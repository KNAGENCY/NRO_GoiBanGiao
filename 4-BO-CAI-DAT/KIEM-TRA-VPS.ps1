# ===========================================================================
#  KIEM TRA CAU HINH VPS + TIM CHO LAM CHAM
#
#  Chay TREN VPS (khong phai may nha). Chi ĐỌC thong tin, KHONG sua gi.
#  Chay xong, chep TOAN BO ket qua gui lai cho tro ly de tu van.
#
#  Cach chay: mo PowerShell tren VPS, go:
#     powershell -ExecutionPolicy Bypass -File "\\tsclient\C\Users\Admin\Music\bo-cai-vps\KIEM-TRA-VPS.ps1"
# ===========================================================================

$ErrorActionPreference = 'SilentlyContinue'
function Muc($t) { Write-Host ''; Write-Host "=== $t ===" -ForegroundColor Cyan }
function D($k, $v) { Write-Host ("  {0,-26}: {1}" -f $k, $v) }

Write-Host ''
Write-Host '###########################################################'
Write-Host '#  BAO CAO CAU HINH VPS - NGOC RONG KAME'
Write-Host '###########################################################'

# --- CPU ------------------------------------------------------------------
Muc 'CPU (bo xu ly)'
$cpu = Get-CimInstance Win32_Processor
D 'Ten'          $cpu.Name
D 'So nhan'      $cpu.NumberOfCores
D 'So luong'     $cpu.NumberOfLogicalProcessors
D 'Toc do (MHz)' $cpu.MaxClockSpeed
$load = (Get-CimInstance Win32_Processor | Measure-Object -Property LoadPercentage -Average).Average
D 'Dang dung (%)' $load

# --- RAM ------------------------------------------------------------------
Muc 'RAM (bo nho)'
$os = Get-CimInstance Win32_OperatingSystem
$tong = [math]::Round($os.TotalVisibleMemorySize/1MB, 2)
$tr   = [math]::Round($os.FreePhysicalMemory/1MB, 2)
D 'Tong cong (GB)' $tong
D 'Con trong (GB)' $tr
D 'Da dung (%)'   ([math]::Round(($tong-$tr)/$tong*100))

# --- DIA ------------------------------------------------------------------
Muc 'O DIA'
Get-CimInstance Win32_LogicalDisk -Filter "DriveType=3" | ForEach-Object {
    D $_.DeviceID ("{0:N1} GB trong / {1:N1} GB tong" -f ($_.FreeSpace/1GB), ($_.Size/1GB))
}
# Loai o dia: SSD hay HDD (SSD nhanh hon nhieu)
Get-PhysicalDisk | ForEach-Object { D ('Loai dia: ' + $_.FriendlyName) $_.MediaType }

# --- WINDOWS --------------------------------------------------------------
Muc 'HE DIEU HANH'
D 'Ban Windows' $os.Caption
D 'Ngay cai'    $os.InstallDate
$uptime = (Get-Date) - $os.LastBootUpTime
D 'Da chay lien tuc' ("{0} ngay {1} gio" -f $uptime.Days, $uptime.Hours)

# --- TIEN TRINH AN NHIEU CPU/RAM NHAT --------------------------------------
Muc '5 TIEN TRINH AN RAM NHAT'
Get-Process | Sort-Object WorkingSet64 -Descending | Select-Object -First 5 |
    ForEach-Object { D $_.ProcessName ("{0:N0} MB" -f ($_.WorkingSet64/1MB)) }

# --- PHP: CO BAT BO NHO DEM (OPcache) CHUA ---------------------------------
Muc 'PHP - BO NHO DEM (yeu to lam cham web nhieu nhat)'
$php = 'C:\xampp\php\php.exe'
if (Test-Path $php) {
    $opcache = & $php -r "echo function_exists('opcache_get_status') && ini_get('opcache.enable') ? 'DA BAT' : 'CHUA BAT';" 2>$null
    D 'OPcache' $opcache
    D 'Phien ban PHP' (& $php -r "echo PHP_VERSION;" 2>$null)
    D 'memory_limit'  (& $php -r "echo ini_get('memory_limit');" 2>$null)
} else { D 'PHP' 'khong thay C:\xampp\php\php.exe' }

# --- WINDOWS DEFENDER: CO QUET THU MUC WEB KHONG ---------------------------
Muc 'WINDOWS DEFENDER (quet file lam cham moi lan mo trang)'
$def = Get-MpPreference
D 'Quet thoi gian thuc' ($(if ($def.DisableRealtimeMonitoring) {'DA TAT'} else {'DANG BAT'}))
$loai = $def.ExclusionPath
if ($loai) { $loai | ForEach-Object { D 'Da loai tru' $_ } } else { D 'Da loai tru' '(chua loai tru thu muc nao)' }

# --- APACHE / MYSQL DANG CHAY? --------------------------------------------
Muc 'DICH VU'
Get-Process httpd, mysqld, java -ErrorAction SilentlyContinue |
    Group-Object ProcessName | ForEach-Object {
        $m = ($_.Group | Measure-Object WorkingSet64 -Sum).Sum/1MB
        D $_.Name ("{0} tien trinh, {1:N0} MB" -f $_.Count, $m)
    }

# --- MANG: nha cung cap / vi tri may chu -----------------------------------
Muc 'MANG'
$ip = (Invoke-RestMethod -Uri 'http://ip-api.com/json/?fields=query,country,city,isp' -TimeoutSec 8)
if ($ip) {
    D 'IP cong khai' $ip.query
    D 'Vi tri'       ("{0}, {1}" -f $ip.city, $ip.country)
    D 'Nha cung cap' $ip.isp
}

Write-Host ''
Write-Host '###########################################################'
Write-Host '#  XONG - chep TOAN BO tu dong dau den day gui cho tro ly'
Write-Host '###########################################################'
Write-Host ''
Read-Host '  Bam Enter de dong'
