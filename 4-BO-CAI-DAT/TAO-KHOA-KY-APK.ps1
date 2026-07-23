# ===========================================================================
#  TAO KHOA KY CHO BAN APK (keystore)
#
#  ---------------------------------------------------------------------
#  KHOA KY LA GI VA VI SAO CAN
#  ---------------------------------------------------------------------
#  Moi file APK deu phai duoc "ky" bang mot khoa. Android dung chu ky do de
#  biet ban cap nhat co that su den tu cung mot nguoi lam ra ban cu khong.
#
#  Hien tai Unity dang tu ky bang KHOA GO LOI - mot khoa tam do chinh Unity
#  sinh ra, nam o:
#        %USERPROFILE%\.android\debug.keystore
#
#  APK ky bang khoa go loi VAN CAI DUOC binh thuong. Van de nam o luc CAP NHAT:
#
#    - Neu ban cai lai Windows, doi may, hoac lo tay xoa file do
#      -> ban APK moi se co chu ky khac
#      -> Android tu choi cap nhat de len ban cu
#      -> MOI NGUOI CHOI phai go game ra roi cai lai tu dau
#
#    - Khoa go loi cung khong dua len Google Play duoc, neu sau nay ban muon.
#
#  File nay tao ra mot khoa RIENG cua ban, dung duoc 100 nam. Giu ky file do
#  la sau nay cap nhat luc nao cung duoc.
#
#  ---------------------------------------------------------------------
#  MOT DIEU KHONG THE CUU DUOC
#  ---------------------------------------------------------------------
#  MAT FILE KHOA HOAC QUEN MAT KHAU LA MAT VINH VIEN. Khong ai khoi phuc
#  duoc - ke ca Google, ke ca toi. Hau qua: moi ban cap nhat sau do deu bat
#  nguoi choi go ra cai lai.
#
#  Vi vay sau khi tao xong, HAY CHEP FILE KHOA RA IT NHAT HAI NOI KHAC
#  (USB, Google Drive, email cho chinh minh) va GHI MAT KHAU RA GIAY.
# ===========================================================================

$ErrorActionPreference = 'Stop'
function Ghi($m, $mau = 'White') { Write-Host "  $m" -ForegroundColor $mau }

# ===========================================================================
#  FILE NAY DA CU - DUNG CHAY NUA
#
#  No chi TAO khoa, chua ky vao file APK va chua dua len trang tai ve.
#  Chay xong de tuong la da xong, trong khi ban APK van dang ky bang khoa
#  go loi - dung cai bay can tranh.
#
#  File thay the lam DU CA BA VIEC:  KY-VA-DANG-APK.ps1  (cung thu muc)
# ===========================================================================
$fileMoi = Join-Path $PSScriptRoot 'KY-VA-DANG-APK.ps1'
Write-Host ''
Write-Host '===========================================================' -ForegroundColor Yellow
Write-Host '  FILE NAY DA CU - HAY CHAY FILE KHAC' -ForegroundColor Yellow
Write-Host '===========================================================' -ForegroundColor Yellow
Write-Host ''
Ghi 'File nay chi TAO khoa. No khong ky vao ban APK va khong dua len trang' 'Yellow'
Ghi 'tai ve, nen chay xong ban APK VAN dang ky bang khoa go loi.' 'Yellow'
Ghi ''
Ghi 'Hay chay file nay thay the - no lam du ca ba viec:' 'Cyan'
Ghi "  $fileMoi" 'Green'
Ghi ''
if (Test-Path $fileMoi) {
    $tl = Read-Host '  Mo file dung ngay bay gio? (go  c  roi Enter)'
    if ($tl -eq 'c' -or $tl -eq 'C') { & $fileMoi; exit $LASTEXITCODE }
}
Ghi 'Da dung lai, chua thay doi gi.' 'Gray'
Read-Host "`n  Bam Enter de thoat"
exit 0

$THU_MUC = 'D:\NRO_Khoa'
$FILE    = Join-Path $THU_MUC 'nrokame.keystore'
$BIDANH  = 'nrokame'

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  TAO KHOA KY CHO BAN APK' -ForegroundColor Cyan
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host ''

# --- Tim keytool ---------------------------------------------------------
# keytool di kem moi ban Java. Uu tien ban OpenJDK cua Unity vi no chac chan
# tuong thich voi thu Unity dung de ky APK.
$ungVien = @(
    'D:\Unity\2022.3.62f2\Editor\Data\PlaybackEngines\AndroidPlayer\OpenJDK\bin\keytool.exe',
    (Join-Path $env:ProgramFiles 'Java\jdk-21\bin\keytool.exe'),
    (Join-Path $env:JAVA_HOME 'bin\keytool.exe')
)
$keytool = $ungVien | Where-Object { $_ -and (Test-Path $_) } | Select-Object -First 1
if (-not $keytool) {
    Ghi 'DUNG LAI: khong tim thay keytool (cong cu tao khoa cua Java).' 'Red'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}
Ghi "[OK] Dung keytool tai: $keytool" 'Green'

# --- Da co khoa chua -----------------------------------------------------
if (Test-Path $FILE) {
    Ghi ''
    Ghi "DUNG LAI: da co khoa tai $FILE" 'Yellow'
    Ghi 'Ghi de len la MAT khoa cu vinh vien, va moi nguoi choi phai go game' 'Red'
    Ghi 'ra cai lai. Neu that su muon lam moi, hay TU TAY doi ten file cu di' 'Red'
    Ghi 'roi chay lai script nay.' 'Red'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}

# --- Mat khau ------------------------------------------------------------
Ghi ''
Ghi 'Dat mat khau cho khoa. Quy dinh:' 'Yellow'
Ghi '  - it nhat 8 ky tu (Java bat buoc)'
Ghi '  - nen dai va de nho, vi dan ra giay chu khong go hang ngay'
Ghi '  - KHONG dung dau nhay  ''  "  hay dau xuyet nguoc  \'
Ghi ''
Ghi 'GHI RA GIAY TRUOC KHI GO. Quen la mat vinh vien.' 'Red'
Ghi ''

$matKhau = $null
for ($lan = 1; $lan -le 5; $lan++) {
    $s1 = Read-Host '  Go mat khau' -AsSecureString
    $p1 = [Runtime.InteropServices.Marshal]::PtrToStringAuto(
          [Runtime.InteropServices.Marshal]::SecureStringToBSTR($s1))
    $s2 = Read-Host '  Go lai lan nua' -AsSecureString
    $p2 = [Runtime.InteropServices.Marshal]::PtrToStringAuto(
          [Runtime.InteropServices.Marshal]::SecureStringToBSTR($s2))

    if ($p1 -ne $p2)       { Ghi 'Hai lan khong giong nhau. Thu lai.' 'Red'; continue }
    if ($p1.Length -lt 8)  { Ghi 'Ngan qua - can it nhat 8 ky tu. Thu lai.' 'Red'; continue }
    if ($p1 -match '["''\\]') { Ghi 'Co ky tu khong duoc phep. Thu lai.' 'Red'; continue }
    $matKhau = $p1
    break
}
if (-not $matKhau) {
    Ghi 'DUNG LAI: khong dat duoc mat khau sau 5 lan.' 'Red'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}

# --- Tao khoa ------------------------------------------------------------
New-Item -ItemType Directory -Force -Path $THU_MUC | Out-Null

# -dname: thong tin chu so huu, nam trong chu ky. Nguoi choi khong nhin thay,
# nhung Google Play co hien khi dua len cho. Dat san cho khoi phai go.
$dname = 'CN=Ngoc Rong Kame, OU=Game, O=Kame, L=Ha Noi, S=Ha Noi, C=VN'

Write-Host ''
Ghi 'Dang tao khoa...' 'Cyan'

# Mat khau dua qua BIEN MOI TRUONG (-storepass:env) chu KHONG viet thang tren
# dong lenh. Ly do: moi thu viet tren dong lenh deu hien ra trong danh sach
# tien trinh cua Windows - ai dang nhap vao may cung doc duoc.
# Bien moi truong chi ton tai trong tien trinh nay va bi xoa ngay sau do.
$cu = $ErrorActionPreference
$env:MKKHOA = $matKhau
$ErrorActionPreference = 'Continue'
try {
    $kq = & $keytool -genkeypair -v `
        -keystore $FILE `
        -alias $BIDANH `
        -keyalg RSA -keysize 2048 -validity 36500 `
        -dname $dname `
        -storepass:env MKKHOA `
        -keypass:env MKKHOA 2>&1
} finally {
    $ErrorActionPreference = $cu
    Remove-Item Env:\MKKHOA -ErrorAction SilentlyContinue
}

if (-not (Test-Path $FILE)) {
    Ghi 'LOI: khong tao duoc khoa.' 'Red'
    $kq | Select-Object -First 8 | ForEach-Object { Ghi "    $_" 'Red' }
    Read-Host "`n  Bam Enter de thoat"; exit 1
}

Ghi "[OK] Da tao khoa: $FILE" 'Green'
Ghi ("     dung luong: {0:N0} byte" -f (Get-Item $FILE).Length) 'Gray'

# --- Kiem chung that -----------------------------------------------------
Write-Host ''
Ghi '--- Kiem chung ---' 'Cyan'
$env:MKKHOA = $matKhau
$ErrorActionPreference = 'Continue'
try { $xem = & $keytool -list -v -keystore $FILE -storepass:env MKKHOA 2>&1 }
finally { $ErrorActionPreference = $cu; Remove-Item Env:\MKKHOA -ErrorAction SilentlyContinue }

$xem | Where-Object { $_ -match 'Alias name|Valid from|Signature algorithm|Entry type' } |
       ForEach-Object { Ghi ("    " + $_.Trim()) 'Gray' }

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  XONG' -ForegroundColor Green
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host ''
Ghi 'LAM NGAY BAY GIO, DUNG DE MAI:' 'Red'
Ghi ''
Ghi ("  1. Chep file  $FILE") 'Yellow'
Ghi '     ra it nhat HAI noi khac (USB, Google Drive, gui email cho chinh minh)'
Ghi ''
Ghi '  2. Ghi mat khau ra GIAY, cat cung cho voi giay ghi mat khau CSDL'
Ghi ''
Ghi '  3. Nhan tro ly: "da tao khoa ky roi" de gan vao Unity'
Ghi ''
Ghi 'Mat file nay hoac quen mat khau = moi nguoi choi phai go game cai lai.' 'Red'
Write-Host ''
Read-Host '  Bam Enter de dong'
