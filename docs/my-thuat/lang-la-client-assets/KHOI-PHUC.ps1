# Khoi phuc art client Lang La (logo + nen + tile map 1) vao D:\NRO_Project.
# Chay: powershell -File KHOI-PHUC.ps1   (hoac chuot phai > Run with PowerShell)
$ErrorActionPreference = "Stop"
$src = $PSScriptRoot
$dst = "D:\NRO_Project\Assets"

if (-not (Test-Path $dst)) { Write-Error "Khong thay $dst — client Unity chua co o may nay."; exit 1 }

# Tao thu muc dich neu thieu
foreach ($d in @("Logo","Resources\res\x2\bg","Resources\res\x4\bg","Resources\res\x2\t","Resources\res\x4\t")) {
    $p = Join-Path $dst $d
    if (-not (Test-Path $p)) { New-Item -ItemType Directory -Force -Path $p | Out-Null }
}

Copy-Item (Join-Path $src "Logo\1.png") (Join-Path $dst "Logo\1.png") -Force
Copy-Item (Join-Path $src "res\x2\bg\*") (Join-Path $dst "Resources\res\x2\bg\") -Force
Copy-Item (Join-Path $src "res\x4\bg\*") (Join-Path $dst "Resources\res\x4\bg\") -Force
Copy-Item (Join-Path $src "res\x2\t\*")  (Join-Path $dst "Resources\res\x2\t\")  -Force
Copy-Item (Join-Path $src "res\x4\t\*")  (Join-Path $dst "Resources\res\x4\t\")  -Force

Write-Output "OK: da copy logo + nen + tile vao $dst"
Write-Output ""
Write-Output "!! BAT BUOC lam trong Unity sau khi copy:"
Write-Output "   1. Chon 8 file Resources/res/x2|x4/bg/b00..b03.png"
Write-Output "   2. Inspector -> tick 'Read/Write Enabled' -> Apply"
Write-Output "      (khong lam thi TROI se DEN — code lay mau troi tu pixel anh)."
Write-Output "   3. Vao map 1 (Lang La) kiem tra nen + mat dat."
