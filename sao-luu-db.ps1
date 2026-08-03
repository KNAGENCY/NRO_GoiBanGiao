# sao-luu-db.ps1 — sao luu CSDL 'nro' bang pg_dump (dinh dang custom -Fc, nen +
# pg_restore doc duoc). Giu N ban moi nhat, tu xoa ban cu hon.
#
# CHONG TU HUY (reviewer F5): dump ra file .partial roi MOI doi thanh .dump khi
# pg_dump thanh cong -> KHONG bao gio co ban .dump 0-byte; retention chi giu ban
# > 0 byte; moi lan ghi vao backups/sao-luu.log de task nen that bai co dau vet.
#
# Chay tay:   powershell -ExecutionPolicy Bypass -File .\sao-luu-db.ps1
# Cai lich chay hang ngay: xem cai-lich-sao-luu.ps1
#
# Khoi phuc tu mot ban (--clean GHI DE du lieu hien tai!):
#   D:\PostgreSQL\bin\pg_restore.exe -h 127.0.0.1 -p <cong> -U postgres -d nro --clean backups\nro-XXX.dump

param(
    [int]$SoBanGiu = 14,
    [int]$Cong = 0,                 # 0 = tu do cong 5432/5433
    [string]$MatKhau = 'nrodev123', # doi tren VPS that (hoac dung .pgpass)
    [string]$TenDB = 'nro'
)
$ErrorActionPreference = 'Stop'

$thuMuc = Join-Path $PSScriptRoot 'backups'
New-Item -ItemType Directory -Force -Path $thuMuc | Out-Null
$logFile = Join-Path $thuMuc 'sao-luu.log'
function Ghi($muc, $msg) {
    $dong = "$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss') [$muc] $msg"
    Add-Content -Path $logFile -Value $dong -Encoding utf8
    if ($muc -eq 'DO') { Write-Host $dong -ForegroundColor Red }
    elseif ($muc -eq 'OK') { Write-Host $dong -ForegroundColor Green }
    else { Write-Host $dong }
}

$pgDump = 'D:\PostgreSQL\bin\pg_dump.exe'
if (-not (Test-Path $pgDump)) { Ghi 'DO' "khong thay $pgDump — sua duong dan PostgreSQL trong script"; exit 1 }

if ($Cong -le 0) {
    $Cong = Get-NetTCPConnection -State Listen -ErrorAction SilentlyContinue |
        Where-Object { $_.LocalPort -in 5432, 5433 } |
        Select-Object -ExpandProperty LocalPort -Unique | Select-Object -First 1
    if (-not $Cong) { Ghi 'DO' "PostgreSQL khong chay (khong thay cong 5432/5433)"; exit 1 }
}

$ts = Get-Date -Format 'yyyyMMdd-HHmmss'
$file = Join-Path $thuMuc "nro-$ts.dump"
$partial = "$file.partial"

Ghi 'INFO' "Sao luu DB '$TenDB' (cong $Cong) -> $file ..."
$env:PGPASSWORD = $MatKhau
# pg_dump la native command: TAM tat Stop quanh no de stderr (vd loi ket noi)
# KHONG thanh terminating error lam nhay qua doan don dep ben duoi. Kiem
# $LASTEXITCODE thu cong.
$eapCu = $ErrorActionPreference
$ErrorActionPreference = 'Continue'
& $pgDump -h 127.0.0.1 -p $Cong -U postgres -Fc -d $TenDB -f $partial
$maDump = $LASTEXITCODE
$ErrorActionPreference = $eapCu
if ($maDump -ne 0) {
    Remove-Item $partial -Force -ErrorAction SilentlyContinue  # KHONG de lai file do dang
    Ghi 'DO' "pg_dump that bai (exit $maDump) - da xoa file do dang, GIU NGUYEN cac ban cu (khong retention)"
    exit 1
}
# Chi doi .partial -> .dump khi pg_dump THANH CONG. Nho vay khong bao gio co ban
# .dump 0-byte de retention giu nham roi day ban that ra ngoai.
Move-Item $partial $file -Force
$kb = [math]::Round((Get-Item $file).Length / 1KB, 1)
Ghi 'OK' "da sao luu $kb KB -> $file"

# Retention: chi giu ban .dump > 0 byte (phong xa), N ban moi nhat, xoa cu hon.
$hopLe = Get-ChildItem $thuMuc -Filter 'nro-*.dump' | Where-Object { $_.Length -gt 0 }
$cu = $hopLe | Sort-Object LastWriteTime -Descending | Select-Object -Skip $SoBanGiu
foreach ($f in $cu) {
    Remove-Item $f.FullName -Force
    Ghi 'INFO' "da xoa ban cu: $($f.Name)"
}
$conLai = (Get-ChildItem $thuMuc -Filter 'nro-*.dump' | Where-Object { $_.Length -gt 0 }).Count
Ghi 'OK' "Xong. Giu $conLai ban hop le (toi da $SoBanGiu)."
