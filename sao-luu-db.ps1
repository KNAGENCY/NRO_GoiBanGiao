# sao-luu-db.ps1 — sao luu CSDL 'nro' bang pg_dump (dinh dang custom -Fc, nen +
# pg_restore doc duoc). Giu N ban moi nhat, tu xoa ban cu hon.
#
# Chay tay:   powershell -ExecutionPolicy Bypass -File .\sao-luu-db.ps1
# Cai lich chay hang ngay: xem cai-lich-sao-luu.ps1
#
# Khoi phuc tu mot ban:
#   D:\PostgreSQL\bin\pg_restore.exe -h 127.0.0.1 -p <cong> -U postgres -d nro --clean backups\nro-XXX.dump
#   (--clean xoa doi tuong cu truoc khi restore — GHI DE du lieu hien tai!)

param(
    [int]$SoBanGiu = 14,
    [int]$Cong = 0,                 # 0 = tu do cong 5432/5433
    [string]$MatKhau = 'nrodev123', # doi tren VPS that
    [string]$TenDB = 'nro'
)
$ErrorActionPreference = 'Stop'

$pgDump = 'D:\PostgreSQL\bin\pg_dump.exe'
if (-not (Test-Path $pgDump)) {
    Write-Host "[DO] khong thay $pgDump — sua duong dan PostgreSQL trong script" -ForegroundColor Red
    exit 1
}

if ($Cong -le 0) {
    $Cong = Get-NetTCPConnection -State Listen -ErrorAction SilentlyContinue |
        Where-Object { $_.LocalPort -in 5432, 5433 } |
        Select-Object -ExpandProperty LocalPort -Unique | Select-Object -First 1
    if (-not $Cong) {
        Write-Host "[DO] PostgreSQL khong chay (khong thay cong 5432/5433)" -ForegroundColor Red
        exit 1
    }
}

$thuMuc = Join-Path $PSScriptRoot 'backups'
New-Item -ItemType Directory -Force -Path $thuMuc | Out-Null
$ts = Get-Date -Format 'yyyyMMdd-HHmmss'
$file = Join-Path $thuMuc "nro-$ts.dump"

Write-Host "Sao luu DB '$TenDB' (cong $Cong) -> $file ..."
$env:PGPASSWORD = $MatKhau
& $pgDump -h 127.0.0.1 -p $Cong -U postgres -Fc -d $TenDB -f $file
if ($LASTEXITCODE -ne 0) {
    Write-Host "[DO] pg_dump that bai (exit $LASTEXITCODE)" -ForegroundColor Red
    exit 1
}
$kb = [math]::Round((Get-Item $file).Length / 1KB, 1)
Write-Host "[OK] da sao luu $kb KB -> $file" -ForegroundColor Green

# Giu N ban moi nhat, xoa cu hon.
$cu = Get-ChildItem $thuMuc -Filter 'nro-*.dump' |
    Sort-Object LastWriteTime -Descending | Select-Object -Skip $SoBanGiu
foreach ($f in $cu) {
    Remove-Item $f.FullName -Force
    Write-Host "  da xoa ban cu: $($f.Name)"
}
$conLai = (Get-ChildItem $thuMuc -Filter 'nro-*.dump').Count
Write-Host "Xong. Giu $conLai ban (toi da $SoBanGiu)."
