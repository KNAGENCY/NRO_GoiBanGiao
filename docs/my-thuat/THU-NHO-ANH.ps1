# ===========================================================================
#  THU NHỎ ẢNH VỀ ĐÚNG KÍCH THƯỚC GAME (giữ nền trong suốt)
#  ---------------------------------------------------------------------
#  Giải quyết đúng vấn đề: Flow/Nano Banana tạo ảnh TO, game cần NHỎ.
#  Công cụ tự tra kích thước đích của TỪNG ảnh theo tên file (từ
#  kich-thuoc-anh.json) rồi thu nhỏ, giữ nguyên kênh alpha (nền trong suốt).
#
#  CÁCH DÙNG:
#    Chuột phải file này -> "Run with PowerShell", rồi dán đường dẫn thư mục
#    ảnh (Flow xuất ra). Hoặc chạy:
#      powershell -File "THU-NHO-ANH.ps1" -ThuMuc "D:\...\anh-flow-xuat"
#
#  Ảnh đã thu nhỏ nằm trong thư mục con "da-thu-nho\" (KHÔNG ghi đè bản gốc).
#  Ảnh nào không có trong bảng kích thước -> báo để xử tay, không đụng.
# ===========================================================================

param(
    [string]$ThuMuc = "",
    [string]$BangKichThuoc = ""
)

Add-Type -AssemblyName System.Drawing

# --- tìm file bảng kích thước (mặc định cạnh script này) ---
if ($BangKichThuoc -eq "") {
    $BangKichThuoc = Join-Path $PSScriptRoot "kich-thuoc-anh.json"
}
if (-not (Test-Path $BangKichThuoc)) {
    Write-Host "KHONG TIM THAY BANG KICH THUOC: $BangKichThuoc" -ForegroundColor Red
    Write-Host "File nay phai nam canh THU-NHO-ANH.ps1 (cung thu muc docs\my-thuat)." -ForegroundColor Red
    Read-Host "Enter de dong"; exit 1
}

if ($ThuMuc -eq "") {
    $ThuMuc = Read-Host "Dan duong dan thu muc anh Flow xuat ra roi Enter"
}
if (-not (Test-Path $ThuMuc)) {
    Write-Host "KHONG TIM THAY THU MUC: $ThuMuc" -ForegroundColor Red
    Read-Host "Enter de dong"; exit 1
}

# --- nạp bảng tên_file -> [rộng, cao] ---
$bang = Get-Content $BangKichThuoc -Raw -Encoding UTF8 | ConvertFrom-Json
# ConvertFrom-Json cho PSCustomObject; đổi sang tra cứu theo tên thuộc tính.

$thuMucRa = Join-Path $ThuMuc "da-thu-nho"
if (-not (Test-Path $thuMucRa)) { New-Item -ItemType Directory -Path $thuMucRa | Out-Null }

$anh = Get-ChildItem $ThuMuc -Recurse -Filter *.png | Where-Object { $_.FullName -notlike "*\da-thu-nho\*" }
if ($anh.Count -eq 0) {
    Write-Host "Khong co file .png nao." -ForegroundColor Red
    Read-Host "Enter de dong"; exit 1
}

Write-Host ""
Write-Host "Dang thu nho $($anh.Count) anh..." -ForegroundColor Cyan
Write-Host ""

$xong = 0
$khongCoTrongBang = @()

foreach ($f in $anh) {
    $ten = $f.Name
    $kt = $bang.$ten     # tra theo ten file (thuoc tinh)
    if ($null -eq $kt) {
        $khongCoTrongBang += $ten
        continue
    }
    $rong = [int]$kt[0]; $cao = [int]$kt[1]

    try {
        $goc = New-Object System.Drawing.Bitmap($f.FullName)
    } catch {
        Write-Host "  Loi mo: $ten" -ForegroundColor Red
        continue
    }

    # anh dich 32bpp ARGB de GIU nen trong suot
    $moi = New-Object System.Drawing.Bitmap($rong, $cao, [System.Drawing.Imaging.PixelFormat]::Format32bppArgb)
    $g = [System.Drawing.Graphics]::FromImage($moi)
    $g.CompositingMode = [System.Drawing.Drawing2D.CompositingMode]::SourceCopy
    $g.InterpolationMode = [System.Drawing.Drawing2D.InterpolationMode]::HighQualityBicubic
    $g.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::HighQuality
    $g.PixelOffsetMode = [System.Drawing.Drawing2D.PixelOffsetMode]::HighQuality
    $g.DrawImage($goc, (New-Object System.Drawing.Rectangle(0, 0, $rong, $cao)))
    $g.Dispose()

    $duongRa = Join-Path $thuMucRa $ten
    $moi.Save($duongRa, [System.Drawing.Imaging.ImageFormat]::Png)
    $moi.Dispose(); $goc.Dispose()
    $xong++
}

# ===========================================================================
Write-Host ("=" * 68)
Write-Host "  KET QUA" -ForegroundColor White
Write-Host ("=" * 68)
Write-Host ("  Da thu nho     : {0,4} / {1}" -f $xong, $anh.Count) -ForegroundColor Green
Write-Host ("  Khong co trong bang kich thuoc: {0,4}" -f $khongCoTrongBang.Count) -ForegroundColor Yellow
if ($khongCoTrongBang.Count -gt 0) {
    Write-Host "  --- Nhung anh nay tao SAI TEN (khong khop danh sach 343) — kiem lai ten: ---" -ForegroundColor Yellow
    $khongCoTrongBang | Select-Object -First 25 | ForEach-Object { Write-Host "    $_" }
    if ($khongCoTrongBang.Count -gt 25) { Write-Host "    ... va $($khongCoTrongBang.Count - 25) file nua" }
}
Write-Host ""
Write-Host "  Anh da thu nho o: $thuMucRa" -ForegroundColor Cyan
Write-Host "  => Buoc tiep: chay KIEM-TRA-ANH.ps1 tren thu muc do de loc anh phang." -ForegroundColor Cyan
Write-Host ("=" * 68)
Read-Host "Enter de dong"
