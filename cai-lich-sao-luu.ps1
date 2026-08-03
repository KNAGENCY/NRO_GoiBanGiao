# cai-lich-sao-luu.ps1 — dang ky Windows Scheduled Task chay sao-luu-db.ps1 hang
# ngay. CHAY MOT LAN TREN VPS (can quyen admin). Khong can chay tren may dev.
#
#   powershell -ExecutionPolicy Bypass -File .\cai-lich-sao-luu.ps1 -Gio 04:00 -MatKhau '<mat khau VPS>'
#
# QUAN TRONG (reviewer F5): PHAI truyen -MatKhau dung cua VPS. Neu khong, task
# nen se chay voi mat khau mac dinh 'nrodev123' -> pg_dump fail moi ngay, backup
# rong. Cach an toan hon: dung file .pgpass roi bo -MatKhau.

param(
    [string]$Gio = '04:00',
    [string]$TenTask = 'NRO-SaoLuuDB',
    [string]$MatKhau = 'nrodev123',
    [int]$Cong = 0,
    [int]$SoBanGiu = 14
)
$ErrorActionPreference = 'Stop'

$script = Join-Path $PSScriptRoot 'sao-luu-db.ps1'
if (-not (Test-Path $script)) {
    Write-Host "[DO] khong thay $script" -ForegroundColor Red
    exit 1
}

# Truyen mat khau/cong/so ban xuong task (F5: task phai dung mat khau DUNG, khong
# de mac dinh). Luu y: mat khau nam trong dinh nghia task dang van ban — tren VPS
# rieng chap nhan duoc; muon kin hon thi dung .pgpass va bo -MatKhau.
$arg = "-NoProfile -ExecutionPolicy Bypass -File `"$script`" -MatKhau `"$MatKhau`" -SoBanGiu $SoBanGiu"
if ($Cong -gt 0) { $arg += " -Cong $Cong" }

$action = New-ScheduledTaskAction -Execute 'powershell.exe' -Argument $arg
$trigger = New-ScheduledTaskTrigger -Daily -At $Gio
Register-ScheduledTask -TaskName $TenTask -Action $action -Trigger $trigger `
    -Description 'Sao luu CSDL nro hang ngay (pg_dump, giu 14 ban, log backups/sao-luu.log)' -Force | Out-Null

Write-Host "[OK] da dang ky Scheduled Task '$TenTask' chay hang ngay luc $Gio." -ForegroundColor Green
Write-Host "Kiem log that bai o: backups\sao-luu.log. Chay thu ngay: Start-ScheduledTask -TaskName '$TenTask'"
