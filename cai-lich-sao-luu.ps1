# cai-lich-sao-luu.ps1 — dang ky Windows Scheduled Task chay sao-luu-db.ps1 hang
# ngay. CHAY MOT LAN TREN VPS (can quyen admin). Khong can chay tren may dev.
#
#   powershell -ExecutionPolicy Bypass -File .\cai-lich-sao-luu.ps1 -Gio 04:00

param(
    [string]$Gio = '04:00',
    [string]$TenTask = 'NRO-SaoLuuDB'
)
$ErrorActionPreference = 'Stop'

$script = Join-Path $PSScriptRoot 'sao-luu-db.ps1'
if (-not (Test-Path $script)) {
    Write-Host "[DO] khong thay $script" -ForegroundColor Red
    exit 1
}

$action = New-ScheduledTaskAction -Execute 'powershell.exe' `
    -Argument "-NoProfile -ExecutionPolicy Bypass -File `"$script`""
$trigger = New-ScheduledTaskTrigger -Daily -At $Gio
Register-ScheduledTask -TaskName $TenTask -Action $action -Trigger $trigger `
    -Description 'Sao luu CSDL nro hang ngay (pg_dump, giu 14 ban)' -Force | Out-Null

Write-Host "[OK] da dang ky Scheduled Task '$TenTask' chay hang ngay luc $Gio." -ForegroundColor Green
Write-Host "Xem/sua: Task Scheduler (taskschd.msc). Chay thu ngay: Start-ScheduledTask -TaskName '$TenTask'"
