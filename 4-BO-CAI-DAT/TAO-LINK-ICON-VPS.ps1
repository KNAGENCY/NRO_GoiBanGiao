# ===========================================================================
#  TAO LINK ICON CHO WEB TREN VPS
#
#  Panel goi anh vat pham o duong dan  /icon/<id>.png  tren web, nhung thu muc
#  do chua tung duoc tao -> anh khong hien. Kho icon that nam trong game root
#  o  data\icon\x3\  (dat ten dung theo icon_id).
#
#  Script nay tao mot LIEN KET (junction): htdocs\icon  ->  data\icon\x3
#  Khong copy gi (khoi ton 100 MB), Apache van phuc vu binh thuong.
#
#  Chay TREN VPS mot lan:
#     powershell -ExecutionPolicy Bypass -File "\\tsclient\C\Users\Admin\Music\bo-cai-vps\TAO-LINK-ICON-VPS.ps1"
# ===========================================================================

$ErrorActionPreference = 'Stop'
function Ghi($m,$c='White'){ Write-Host "  $m" -ForegroundColor $c }

$LINK = 'C:\xampp\htdocs\icon'

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  TAO LINK ICON CHO WEB' -ForegroundColor Cyan
Write-Host '===========================================================' -ForegroundColor Cyan

# Tim game root tren VPS + thu muc icon. Uu tien x3 (vua dep), lui ve x2/x4.
$goc = @('C:\zz1','C:\Users\Administrator\zz1','C:\game\zz1') | Where-Object { Test-Path "$_\data\icon" } | Select-Object -First 1
if (-not $goc) {
    Ghi 'DUNG LAI: khong tim thay thu muc data\icon trong game root.' 'Red'
    Ghi 'Sua bien $goc trong script cho dung duong dan game tren VPS.' 'Yellow'
    Read-Host "`n  Enter de thoat"; exit 1
}
$target = @("$goc\data\icon\x3","$goc\data\icon\x2","$goc\data\icon\x4") | Where-Object { (Test-Path $_) -and (Get-ChildItem $_ -File -ErrorAction SilentlyContinue | Select-Object -First 1) } | Select-Object -First 1
if (-not $target) { Ghi "DUNG LAI: $goc\data\icon khong co anh trong x2/x3/x4." 'Red'; Read-Host; exit 1 }
Ghi "[OK] Kho icon: $target" 'Green'

# Neu da co link/thu muc -> bao va dung
if (Test-Path $LINK) {
    $it = Get-Item $LINK -Force
    if ($it.LinkType) { Ghi "[OK] Da co link san: $LINK -> $($it.Target)" 'Green' }
    else { Ghi "[CHU Y] $LINK da ton tai (thu muc that, khong phai link). Giu nguyen." 'Yellow' }
} else {
    cmd /c "mklink /J `"$LINK`" `"$target`"" | ForEach-Object { Ghi $_ Gray }
    Ghi "[OK] Da tao link $LINK -> $target" 'Green'
}

# Kiem chung: tai thu 1 icon qua web
Write-Host ''
Ghi '--- Kiem chung ---' 'Cyan'
Start-Sleep -Seconds 1
$mau = Get-ChildItem $target -Filter '*.png' -File | Select-Object -First 1
if ($mau) {
    $id = [IO.Path]::GetFileNameWithoutExtension($mau.Name)
    try {
        $r = Invoke-WebRequest -Uri "http://localhost/icon/$id.png" -Method Head -UseBasicParsing -TimeoutSec 15
        if ($r.StatusCode -eq 200) { Ghi "[OK] http://localhost/icon/$id.png tra ve anh ($($r.Headers['Content-Length']) byte)" 'Green' }
        else { Ghi "[?] HTTP $($r.StatusCode)" 'Yellow' }
    } catch { Ghi ("[!!] Web chua phuc vu duoc icon: " + $_.Exception.Message) 'Red' }
}
Write-Host ''
Ghi 'Xong. Mo lai panel, anh vat pham se hien.' 'Green'
Ghi '(Neu trinh duyet van chua thay: bam Ctrl+F5 de tai lai bo cache.)' 'Gray'
Write-Host ''
Read-Host '  Bam Enter de dong'
