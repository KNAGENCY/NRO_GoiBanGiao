# ===========================================================================
#  KHOA CONG REMOTE DESKTOP
#
#  Lam 2 viec:
#    1. Doi cong Remote Desktop tu 3389 sang 43389
#    2. Chi cho phep dai IP 14.248.0.0/16 (VNPT khu vuc cua ban) ket noi
#
#  AN TOAN: cong 3389 cu VAN DUOC GIU LAI nhung cung bi gioi han theo dai IP.
#  Neu cong moi khong vao duoc, ban van vao duoc bang cong 3389 cu.
#  Sau khi chac chan cong moi chay, hay chay file DONG-CONG-3389-CU.ps1
#
#  MUON HUY BO TAT CA: chay file HUY-KHOA-REMOTE-DESKTOP.ps1
# ===========================================================================

$ErrorActionPreference = 'Stop'

$CONG_MOI = 43389
$DAI_IP   = '14.248.0.0/16'
$IP_NHA   = '14.248.216.163'
$RDP_REG  = 'HKLM:\SYSTEM\CurrentControlSet\Control\Terminal Server\WinStations\RDP-Tcp'

function Ghi($m, $mau = 'White') { Write-Host "  $m" -ForegroundColor $mau }

Write-Host ""
Write-Host "===========================================================" -ForegroundColor Cyan
Write-Host "  KHOA CONG REMOTE DESKTOP - NGOC RONG KAME" -ForegroundColor Cyan
Write-Host "===========================================================" -ForegroundColor Cyan
Write-Host ""

# --- Kiem tra quyen ------------------------------------------------------
$quyen = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()
         ).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $quyen) {
    Ghi "LOI: Phai chay bang quyen Administrator." 'Red'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}

# --- Canh bao ------------------------------------------------------------
Write-Host "  CANH BAO - DOC KY TRUOC KHI TIEP TUC" -ForegroundColor Yellow
Write-Host ""
Ghi "Sau khi chay xong, dia chi ket noi Remote Desktop se doi thanh:" 'Yellow'
Write-Host ""
Write-Host "        103.116.38.60:$CONG_MOI" -ForegroundColor Green
Write-Host ""
Ghi "HAY GHI CON SO NAY RA GIAY NGAY BAY GIO." 'Yellow'
Ghi "Phien lam viec hien tai se bi ngat. Ban phai ket noi lai bang dia chi tren." 'Yellow'
Write-Host ""
Ghi "Duong lui neu that bai:" 'Gray'
Ghi "  - Cong 3389 cu VAN MO (chi gioi han theo dai IP) -> van vao duoc" 'Gray'
Ghi "  - Hoac dung nut Console/VNC trong bang dieu khien DNCloud" 'Gray'
Ghi "  - Hoac chay file HUY-KHOA-REMOTE-DESKTOP.ps1 de tra lai nhu cu" 'Gray'
Write-Host ""

$traLoi = Read-Host "  Da ghi so 43389 ra giay chua? Go CO de tiep tuc"
if ($traLoi -ne 'CO' -and $traLoi -ne 'co') {
    Ghi "Da huy. Khong thay doi gi." 'Yellow'
    Read-Host "`n  Bam Enter de thoat"; exit 0
}

# --- Luu lai trang thai cu de con duong ma huy ---------------------------
$luu = "$env:USERPROFILE\Desktop\KAME-thong-tin-remote-desktop.txt"
@(
    "THONG TIN KET NOI REMOTE DESKTOP - NGOC RONG KAME",
    "Ngay doi: $(Get-Date -Format 'dd/MM/yyyy HH:mm')",
    "",
    "DIA CHI KET NOI MOI : 103.116.38.60:$CONG_MOI",
    "Cong cu (van dung duoc): 103.116.38.60:3389",
    "Chi cho phep tu dai  : $DAI_IP",
    "",
    "Neu khong vao duoc: dung nut Console/VNC trong bang dieu khien DNCloud,",
    "roi chay file HUY-KHOA-REMOTE-DESKTOP.ps1 de tra lai nhu cu."
) | Set-Content -Encoding UTF8 $luu
Ghi "[OK] Da luu thong tin ra: $luu" 'Green'

# --- 1. Mo cong moi, chi cho dai IP cua ban -----------------------------
Write-Host ""
Ghi "--- BUOC 1: Mo cong $CONG_MOI cho dai $DAI_IP ---" 'Cyan'

Get-NetFirewallRule -DisplayName 'Kame - RDP cong moi' -ErrorAction SilentlyContinue | Remove-NetFirewallRule
New-NetFirewallRule -DisplayName 'Kame - RDP cong moi' -Direction Inbound -Protocol TCP `
    -LocalPort $CONG_MOI -RemoteAddress @($DAI_IP, $IP_NHA) -Action Allow -Profile Any | Out-Null
Ghi "[OK] Da mo cong $CONG_MOI, chi cho $DAI_IP" 'Green'

# --- 2. Gioi han cac luat RDP san co theo dai IP -------------------------
Write-Host ""
Ghi "--- BUOC 2: Gioi han cong 3389 cu theo dai IP ---" 'Cyan'

$soLuat = 0
Get-NetFirewallRule -ErrorAction SilentlyContinue |
    Where-Object { $_.DisplayName -match 'Remote Desktop' -and $_.Direction -eq 'Inbound' -and $_.Enabled -eq 'True' } |
    ForEach-Object {
        Set-NetFirewallRule -Name $_.Name -RemoteAddress @($DAI_IP, $IP_NHA) -ErrorAction SilentlyContinue
        $soLuat++
    }
Ghi "[OK] Da gioi han $soLuat luat Remote Desktop san co" 'Green'
Ghi "     Cong 3389 van mo NHUNG chi tu dai IP cua ban - day la duong lui" 'Gray'

# --- 3. Doi cong Remote Desktop -----------------------------------------
Write-Host ""
Ghi "--- BUOC 3: Doi cong Remote Desktop sang $CONG_MOI ---" 'Cyan'

$congCu = (Get-ItemProperty -Path $RDP_REG -Name PortNumber).PortNumber
Ghi "Cong hien tai: $congCu" 'Gray'
Set-ItemProperty -Path $RDP_REG -Name PortNumber -Value $CONG_MOI -Type DWord
$congSau = (Get-ItemProperty -Path $RDP_REG -Name PortNumber).PortNumber
if ($congSau -ne $CONG_MOI) {
    Ghi "LOI: Doi cong that bai. Khong thay doi gi them." 'Red'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}
Ghi "[OK] Da doi sang cong $congSau" 'Green'

# --- Tong ket ------------------------------------------------------------
Write-Host ""
Write-Host "===========================================================" -ForegroundColor Cyan
Write-Host "  DA XONG. CAN KHOI DONG LAI DE CO HIEU LUC." -ForegroundColor Cyan
Write-Host "===========================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Ket noi lai bang:  103.116.38.60:$CONG_MOI" -ForegroundColor Green
Write-Host "  Neu khong duoc:    103.116.38.60:3389  (van con hoat dong)" -ForegroundColor Yellow
Write-Host ""

$kd = Read-Host "  Khoi dong lai may bay gio? (go CO de khoi dong lai, Enter de tu lam sau)"
if ($kd -eq 'CO' -or $kd -eq 'co') {
    Ghi "May se khoi dong lai sau 10 giay. Doi khoang 1-2 phut roi ket noi lai." 'Yellow'
    Start-Sleep -Seconds 3
    shutdown /r /t 10 /c "Doi cong Remote Desktop - Ngoc Rong Kame"
} else {
    Ghi "Nho khoi dong lai may sau. Truoc khi khoi dong lai, cong cu 3389 van dung duoc." 'Yellow'
    Read-Host "`n  Bam Enter de dong"
}
