# ===========================================================================
#  BUOC 1 - CAI NEN TANG CHO VPS
#  Cai JDK 21 + XAMPP, mo cong tuong lua. Chay tren VPS bang PowerShell
#  voi quyen Administrator. Chay lai nhieu lan khong sao.
# ===========================================================================

$ErrorActionPreference = 'Stop'
$ProgressPreference    = 'SilentlyContinue'   # khong co dong nay tai file se cham gap 10 lan

function Ghi($m, $mau = 'White') { Write-Host "  $m" -ForegroundColor $mau }

Write-Host ""
Write-Host "===========================================================" -ForegroundColor Cyan
Write-Host "  NGOC RONG KAME - CAI NEN TANG CHO MAY CHU" -ForegroundColor Cyan
Write-Host "===========================================================" -ForegroundColor Cyan
Write-Host ""

# --- Kiem tra quyen Administrator ---------------------------------------
$quyen = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()
         ).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $quyen) {
    Ghi "LOI: Phai chay PowerShell bang quyen Administrator." 'Red'
    Ghi "     Bam chuot phai vao PowerShell -> Run as administrator." 'Yellow'
    Read-Host "`n  Bam Enter de thoat"
    exit 1
}
Ghi "[OK] Dang chay voi quyen Administrator" 'Green'

# --- Kiem tra dung luong o C: -------------------------------------------
$oC = Get-CimInstance Win32_LogicalDisk -Filter "DeviceID='C:'"
$trong = [math]::Round($oC.FreeSpace / 1GB, 1)
if ($trong -lt 10) {
    Ghi "LOI: O C: chi con $trong GB. Can it nhat 10 GB." 'Red'
    Read-Host "`n  Bam Enter de thoat"
    exit 1
}
Ghi "[OK] O C: con trong $trong GB" 'Green'

$tam = "C:\kame-tam"
New-Item -ItemType Directory -Path $tam -Force | Out-Null

# --- 1. JDK 21 ----------------------------------------------------------
Write-Host ""
Ghi "--- BUOC 1/3: JAVA 21 ---" 'Cyan'

$javaDir = "C:\Program Files\Java\jdk-21"
if (Test-Path "$javaDir\bin\java.exe") {
    Ghi "[BO QUA] Java 21 da co san o $javaDir" 'Yellow'
} else {
    $msi = "$tam\jdk21.msi"
    if (-not (Test-Path $msi)) {
        Ghi "Dang tai JDK 21 (179 MB) - doi vai phut..."
        Invoke-WebRequest -Uri "https://api.adoptium.net/v3/installer/latest/21/ga/windows/x64/jdk/hotspot/normal/eclipse" -OutFile $msi -UseBasicParsing
    }
    Ghi "Dang cai JDK 21..."
    $p = Start-Process msiexec.exe -ArgumentList @(
        "/i", "`"$msi`"", "/qn", "/norestart",
        "ADDLOCAL=FeatureMain,FeatureEnvironment,FeatureJarFileRunWith",
        "INSTALLDIR=`"$javaDir\`""
    ) -Wait -PassThru
    if ($p.ExitCode -ne 0) { Ghi "LOI: msiexec tra ve ma $($p.ExitCode)" 'Red'; exit 1 }
    Ghi "[OK] Da cai JDK 21" 'Green'
}

# JAVA_HOME + PATH (dat o muc May, khong phai muc Nguoi dung)
[Environment]::SetEnvironmentVariable('JAVA_HOME', $javaDir, 'Machine')
$path = [Environment]::GetEnvironmentVariable('Path', 'Machine')
if ($path -notlike "*$javaDir\bin*") {
    [Environment]::SetEnvironmentVariable('Path', "$path;$javaDir\bin", 'Machine')
}
$env:JAVA_HOME = $javaDir
$env:Path      = "$env:Path;$javaDir\bin"
Ghi "[OK] Da dat JAVA_HOME = $javaDir" 'Green'
Ghi "     Phien ban: $((& "$javaDir\bin\java.exe" -version 2>&1)[0])" 'Gray'

# --- 2. XAMPP -----------------------------------------------------------
Write-Host ""
Ghi "--- BUOC 2/3: XAMPP (Apache + MariaDB + PHP) ---" 'Cyan'

if (Test-Path "C:\xampp\xampp-control.exe") {
    Ghi "[BO QUA] XAMPP da co san o C:\xampp" 'Yellow'
} else {
    $exe = "$tam\xampp-installer.exe"
    if (-not (Test-Path $exe)) {
        Ghi "Dang tai XAMPP 8.2.12 (150 MB) - doi vai phut..."
        Invoke-WebRequest -UseBasicParsing -OutFile $exe `
            -Uri "https://sourceforge.net/projects/xampp/files/XAMPP%20Windows/8.2.12/xampp-windows-x64-8.2.12-0-VS16-installer.exe/download"
    }
    Ghi "Dang cai XAMPP vao C:\xampp - buoc nay lau nhat, co the toi 10 phut..."
    $p = Start-Process $exe -ArgumentList @(
        "--mode", "unattended", "--unattendedmodeui", "minimal", "--prefix", "C:\xampp"
    ) -Wait -PassThru
    if (-not (Test-Path "C:\xampp\xampp-control.exe")) {
        Ghi "LOI: Cai XAMPP that bai (ma thoat $($p.ExitCode))" 'Red'
        exit 1
    }
    Ghi "[OK] Da cai XAMPP" 'Green'
}

# --- 3. Tuong lua -------------------------------------------------------
Write-Host ""
Ghi "--- BUOC 3/3: MO CONG TUONG LUA ---" 'Cyan'

$congCanMo = @(
    @{ Ten = 'Kame - Web (HTTP)';   Cong = 80    },
    @{ Ten = 'Kame - Web (HTTPS)';  Cong = 443   },
    @{ Ten = 'Kame - Server game';  Cong = 14445 }
)
foreach ($c in $congCanMo) {
    if (Get-NetFirewallRule -DisplayName $c.Ten -ErrorAction SilentlyContinue) {
        Ghi "[BO QUA] Cong $($c.Cong) da mo tu truoc" 'Yellow'
    } else {
        New-NetFirewallRule -DisplayName $c.Ten -Direction Inbound -Protocol TCP `
            -LocalPort $c.Cong -Action Allow -Profile Any | Out-Null
        Ghi "[OK] Da mo cong $($c.Cong) - $($c.Ten)" 'Green'
    }
}

# Cong 3306 (MySQL) TUYET DOI khong duoc mo ra internet.
if (Get-NetFirewallRule -DisplayName 'Kame - CHAN MySQL tu ngoai' -ErrorAction SilentlyContinue) {
    Ghi "[BO QUA] Luat chan MySQL da co" 'Yellow'
} else {
    New-NetFirewallRule -DisplayName 'Kame - CHAN MySQL tu ngoai' -Direction Inbound `
        -Protocol TCP -LocalPort 3306 -Action Block -Profile Any | Out-Null
    Ghi "[OK] Da CHAN cong 3306 - CSDL khong lo ra internet" 'Green'
}

# --- Mui gio ------------------------------------------------------------
Set-TimeZone -Id 'SE Asia Standard Time' -ErrorAction SilentlyContinue
Ghi "[OK] Mui gio: $((Get-TimeZone).Id)" 'Green'

# --- Tong ket -----------------------------------------------------------
Write-Host ""
Write-Host "===========================================================" -ForegroundColor Cyan
Write-Host "  XONG BUOC 1. TINH TRANG HIEN TAI:" -ForegroundColor Cyan
Write-Host "===========================================================" -ForegroundColor Cyan
Ghi "Java 21      : $(if (Test-Path "$javaDir\bin\java.exe") {'DA CAI'} else {'THAT BAI'})"
Ghi "XAMPP        : $(if (Test-Path 'C:\xampp\xampp-control.exe') {'DA CAI'} else {'THAT BAI'})"
Ghi "Cong 80/443  : da mo"
Ghi "Cong 14445   : da mo"
Ghi "Cong 3306    : DA CHAN (dung nhu vay)"
Write-Host ""
Ghi "Buoc tiep theo: bao lai cho tro ly de nhan bo du lieu game." 'Yellow'
Write-Host ""
Read-Host "  Bam Enter de dong"
