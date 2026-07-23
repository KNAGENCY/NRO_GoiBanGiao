# ===========================================================================
#  HUY KHOA REMOTE DESKTOP - TRA LAI NHU CU
#
#  Chay file nay khi ban bi khoa ngoai, hoac muon quay lai trang thai ban dau.
#  Vao may bang nut Console/VNC trong bang dieu khien DNCloud roi chay file nay.
#
#  Sau khi chay: cong tro lai 3389, cho phep moi dia chi ket noi.
# ===========================================================================

$ErrorActionPreference = 'Stop'
$RDP_REG = 'HKLM:\SYSTEM\CurrentControlSet\Control\Terminal Server\WinStations\RDP-Tcp'

function Ghi($m, $mau = 'White') { Write-Host "  $m" -ForegroundColor $mau }

Write-Host ""
Write-Host "===========================================================" -ForegroundColor Yellow
Write-Host "  HUY KHOA REMOTE DESKTOP - TRA VE MAC DINH" -ForegroundColor Yellow
Write-Host "===========================================================" -ForegroundColor Yellow
Write-Host ""

$quyen = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()
         ).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $quyen) {
    Ghi "LOI: Phai chay bang quyen Administrator." 'Red'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}

# 1. Tra cong ve 3389
Set-ItemProperty -Path $RDP_REG -Name PortNumber -Value 3389 -Type DWord
Ghi "[OK] Da tra cong ve 3389" 'Green'

# 2. Bo gioi han dia chi tren cac luat Remote Desktop
$n = 0
Get-NetFirewallRule -ErrorAction SilentlyContinue |
    Where-Object { $_.DisplayName -match 'Remote Desktop' -and $_.Direction -eq 'Inbound' } |
    ForEach-Object { Set-NetFirewallRule -Name $_.Name -RemoteAddress Any -ErrorAction SilentlyContinue; $n++ }
Ghi "[OK] Da bo gioi han dia chi tren $n luat" 'Green'

# 3. Xoa luat cong moi
Get-NetFirewallRule -DisplayName 'Kame - RDP cong moi' -ErrorAction SilentlyContinue | Remove-NetFirewallRule
Ghi "[OK] Da xoa luat cong 43389" 'Green'

# 4. Bat lai luat Remote Desktop mac dinh phong khi bi tat
Get-NetFirewallRule -ErrorAction SilentlyContinue |
    Where-Object { $_.DisplayName -match 'Remote Desktop' -and $_.Direction -eq 'Inbound' } |
    Enable-NetFirewallRule -ErrorAction SilentlyContinue
Ghi "[OK] Da bat lai cac luat Remote Desktop mac dinh" 'Green'

Write-Host ""
Ghi "XONG. Ket noi lai bang: 103.116.38.60:3389" 'Cyan'
Ghi "CANH BAO: cong 3389 gio mo cho CA THE GIOI. Nho khoa lai sau." 'Yellow'
Write-Host ""

$kd = Read-Host "  Khoi dong lai may de co hieu luc? (go CO)"
if ($kd -eq 'CO' -or $kd -eq 'co') { shutdown /r /t 10 /c "Huy khoa Remote Desktop" }
else { Read-Host "`n  Bam Enter de dong" }
