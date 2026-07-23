# ===========================================================================
#  KY BAN APK BANG KHOA RIENG ROI DUA LEN TRANG TAI VE
#
#  Chay MOT LAN duy nhat truoc khi cho nguoi choi tai. Sau nay moi lan build
#  APK moi thi chay lai file nay - no dung lai dung khoa cu.
#
#  ---------------------------------------------------------------------
#  VI SAO PHAI LAM, NOI BANG TIENG VIET
#  ---------------------------------------------------------------------
#  Moi file APK deu duoc "ky" bang mot khoa. Android dung chu ky do de biet
#  ban cap nhat co that su den tu cung mot nguoi lam ra ban cu hay khong.
#
#  Ban APK Unity vua tao ra dang ky bang KHOA GO LOI - mot khoa tam do chinh
#  Unity sinh ra o  %USERPROFILE%\.android\debug.keystore
#
#  APK ky bang khoa go loi VAN CAI DUOC binh thuong. Van de nam o luc CAP NHAT:
#
#    - Cai lai Windows, doi may, hoac lo tay xoa file khoa do
#      -> ban APK moi se co chu ky khac
#      -> Android TU CHOI cap nhat de len ban cu
#      -> MOI NGUOI CHOI phai go game ra roi cai lai tu dau
#
#  Vi vay: nguoi choi dau tien tai ve PHAI la ban da ky bang khoa rieng.
#  Sua sau khi da co nguoi cai roi thi khong sua duoc nua.
#
#  ---------------------------------------------------------------------
#  MOT DIEU KHONG THE CUU DUOC
#  ---------------------------------------------------------------------
#  MAT FILE KHOA HOAC QUEN MAT KHAU LA MAT VINH VIEN. Khong ai khoi phuc
#  duoc - ke ca Google, ke ca toi.
#
#  Sau khi chay xong: CHEP FILE KHOA RA IT NHAT HAI NOI KHAC va GHI MAT KHAU
#  RA GIAY.
# ===========================================================================

param(
    # Chi dung khi KIEM THU script nay. Nguoi dung binh thuong KHONG truyen
    # tham so nay - de trong thi script se tu hoi mat khau.
    [System.Security.SecureString]$MatKhauKiemThu,

    [string]$FileKhoa   = 'D:\NRO_Khoa\nrokame.keystore',
    [string]$FileApk    = 'D:\NRO_Build\Android\NgocRongKame.apk',
    [string]$ThuMucWeb  = 'C:\xampp\htdocs\sv1\down',
    [string]$ThuMucGoi  = 'C:\Users\Admin\Music\goi-len-vps',

    # Bo qua buoc chep ra web / goi gui VPS (dung khi kiem thu).
    [switch]$KhongChep
)

$ErrorActionPreference = 'Stop'
$BIDANH = 'nrokame'

function Ghi($m, $mau = 'White') { Write-Host "  $m" -ForegroundColor $mau }
function Dung($m) { Ghi ''; Ghi "DUNG LAI: $m" 'Red'; if (-not $KhongChep) { Read-Host "`n  Bam Enter de thoat" }; exit 1 }

# Chay mot chuong trinh ngoai va TRA VE ca dau ra lan ma thoat.
# Vi sao can ham nay: PowerShell 5.1 coi moi dong stderr cua chuong trinh
# ngoai la "loi" khi $ErrorActionPreference = 'Stop', ke ca khi chuong trinh
# chay thanh cong. Phai ha xuong 'Continue' trong luc goi roi dat lai.
function ChayNgoai {
    param([string]$File, [string[]]$ThamSo)
    $cu = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    try {
        $ra = & $File @ThamSo 2>&1
        $ma = $LASTEXITCODE
    } finally { $ErrorActionPreference = $cu }
    return [pscustomobject]@{ Ra = ($ra | Out-String); Ma = $ma }
}

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  KY APK BANG KHOA RIENG VA DUA LEN TRANG TAI VE' -ForegroundColor Cyan
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host ''

# --- 1. Tim cong cu -------------------------------------------------------
$ap = 'D:\Unity\2022.3.62f2\Editor\Data\PlaybackEngines\AndroidPlayer'

$keytool = @(
    "$ap\OpenJDK\bin\keytool.exe",
    (Join-Path $env:ProgramFiles 'Java\jdk-21\bin\keytool.exe')
) | Where-Object { $_ -and (Test-Path $_) } | Select-Object -First 1
if (-not $keytool) { Dung 'khong tim thay keytool (cong cu tao khoa cua Java).' }

$bt = Get-ChildItem "$ap\SDK\build-tools" -Directory -ErrorAction SilentlyContinue |
      Sort-Object Name -Descending | Select-Object -First 1
if (-not $bt) { Dung "khong tim thay build-tools trong $ap\SDK" }

$apksigner = Join-Path $bt.FullName 'apksigner.bat'
$zipalign  = Join-Path $bt.FullName 'zipalign.exe'
foreach ($f in @($apksigner, $zipalign)) { if (-not (Test-Path $f)) { Dung "thieu $f" } }

# apksigner.bat can JAVA_HOME de tim java.
$env:JAVA_HOME = Split-Path (Split-Path $keytool -Parent) -Parent

Ghi "[OK] keytool   : $keytool" 'Green'
Ghi "[OK] apksigner : $apksigner" 'Green'

if (-not (Test-Path $FileApk)) {
    Dung "khong thay file APK tai $FileApk`n         Hay build APK truoc (nho tro ly build lai)."
}
$kt = Get-Item $FileApk
Ghi ("[OK] APK       : {0}  ({1:N1} MB)" -f $FileApk, ($kt.Length / 1MB)) 'Green'

# --- 2. Mat khau ----------------------------------------------------------
function DocChuoi([System.Security.SecureString]$s) {
    [Runtime.InteropServices.Marshal]::PtrToStringAuto(
        [Runtime.InteropServices.Marshal]::SecureStringToBSTR($s))
}

$daCoKhoa = Test-Path $FileKhoa
$matKhau  = $null

if ($MatKhauKiemThu) {
    $matKhau = DocChuoi $MatKhauKiemThu
} else {
    Write-Host ''
    if ($daCoKhoa) {
        Ghi "Da co san khoa tai: $FileKhoa" 'Yellow'
        Ghi 'Go dung mat khau cua khoa do (mat khau ban da ghi ra giay lan truoc).'
        Write-Host ''
        for ($lan = 1; $lan -le 5; $lan++) {
            $s = Read-Host '  Mat khau khoa' -AsSecureString
            $matKhau = DocChuoi $s
            if ($matKhau.Length -ge 8) { break }
            Ghi 'Ngan qua. Thu lai.' 'Red'; $matKhau = $null
        }
    } else {
        Ghi 'Chua co khoa - se tao moi. Dat mat khau cho khoa:' 'Yellow'
        Ghi '  - it nhat 8 ky tu'
        Ghi '  - nen dai va de nho, vi dan ra giay chu khong go hang ngay'
        Ghi '  - KHONG dung dau nhay  ''  "  hay dau xuyet nguoc  \'
        Write-Host ''
        Ghi 'GHI RA GIAY TRUOC KHI GO. Quen la mat vinh vien.' 'Red'
        Write-Host ''
        for ($lan = 1; $lan -le 5; $lan++) {
            $p1 = DocChuoi (Read-Host '  Go mat khau'     -AsSecureString)
            $p2 = DocChuoi (Read-Host '  Go lai lan nua'  -AsSecureString)
            if ($p1 -ne $p2)          { Ghi 'Hai lan khong giong nhau. Thu lai.' 'Red'; continue }
            if ($p1.Length -lt 8)     { Ghi 'Ngan qua - can it nhat 8 ky tu. Thu lai.' 'Red'; continue }
            if ($p1 -match '["''\\]') { Ghi 'Co ky tu khong duoc phep. Thu lai.' 'Red'; continue }
            $matKhau = $p1; break
        }
    }
}
if (-not $matKhau) { Dung 'khong dat duoc mat khau.' }

# Mat khau di qua BIEN MOI TRUONG chu KHONG viet thang tren dong lenh.
# Ly do: moi thu viet tren dong lenh deu hien ra trong danh sach tien trinh
# cua Windows - ai dang nhap vao may cung doc duoc.
$env:MKKHOA = $matKhau

try {
    # --- 3. Tao khoa neu chua co -----------------------------------------
    if (-not $daCoKhoa) {
        Write-Host ''
        Ghi 'Dang tao khoa moi...' 'Cyan'
        New-Item -ItemType Directory -Force -Path (Split-Path $FileKhoa -Parent) | Out-Null

        # -dname: thong tin chu so huu nam trong chu ky. Nguoi choi khong nhin
        # thay; dat san cho khoi phai go tay.
        $dname = 'CN=Ngoc Rong Kame, OU=Game, O=Kame, L=Ha Noi, S=Ha Noi, C=VN'

        $r = ChayNgoai $keytool @(
            '-genkeypair', '-keystore', $FileKhoa, '-alias', $BIDANH,
            '-keyalg', 'RSA', '-keysize', '2048', '-validity', '36500',
            '-dname', $dname,
            '-storepass:env', 'MKKHOA', '-keypass:env', 'MKKHOA')

        if (-not (Test-Path $FileKhoa)) {
            Ghi 'LOI khi tao khoa:' 'Red'
            ($r.Ra -split "`n" | Select-Object -First 6) | ForEach-Object { Ghi "    $_" 'Red' }
            Dung 'khong tao duoc khoa.'
        }
        Ghi ("[OK] Da tao khoa: {0}  ({1:N0} byte)" -f $FileKhoa, (Get-Item $FileKhoa).Length) 'Green'
    } else {
        # Khoa da co - thu mo bang mat khau vua go de bao loi som, ro rang,
        # thay vi de apksigner bao loi kho hieu o buoc sau.
        $r = ChayNgoai $keytool @('-list', '-keystore', $FileKhoa, '-storepass:env', 'MKKHOA')
        if ($r.Ma -ne 0) { Dung 'sai mat khau khoa. Kiem tra lai to giay ghi mat khau.' }
        Ghi '[OK] Mo duoc khoa cu bang mat khau vua go' 'Green'
    }

    # --- 4. Can le (zipalign) --------------------------------------------
    # zipalign xep lai cac file ben trong APK cho thang hang 4 byte. Android
    # doc file da can le nhanh hon va ton it RAM hon. PHAI lam TRUOC khi ky:
    # can le sau khi ky se pha hong chu ky.
    Write-Host ''
    Ghi 'Dang can le APK...' 'Cyan'
    $tam = "$FileApk.tam"
    if (Test-Path $tam) { Remove-Item $tam -Force }

    $r = ChayNgoai $zipalign @('-p', '-f', '4', $FileApk, $tam)
    if ($r.Ma -ne 0 -or -not (Test-Path $tam)) {
        Ghi ($r.Ra.Trim()) 'Red'; Dung 'zipalign that bai.'
    }
    Ghi '[OK] Da can le' 'Green'

    # --- 5. Ky -----------------------------------------------------------
    Ghi 'Dang ky...' 'Cyan'
    $r = ChayNgoai $apksigner @(
        'sign', '--ks', $FileKhoa, '--ks-key-alias', $BIDANH,
        '--ks-pass', 'env:MKKHOA', '--key-pass', 'env:MKKHOA',
        '--v1-signing-enabled', 'true',
        '--v2-signing-enabled', 'true',
        '--v3-signing-enabled', 'true',
        '--out', $FileApk, $tam)

    Remove-Item $tam -Force -ErrorAction SilentlyContinue
    Remove-Item "$FileApk.idsig" -Force -ErrorAction SilentlyContinue

    if ($r.Ma -ne 0) { Ghi ($r.Ra.Trim()) 'Red'; Dung 'ky that bai.' }
    Ghi '[OK] Da ky' 'Green'

    # --- 6. Kiem chung that ----------------------------------------------
    # Khong tin buoc tren "chay khong bao loi la xong". Doc lai chu ky trong
    # file APK va doi chieu: phai la khoa Kame, KHONG duoc la khoa go loi.
    Write-Host ''
    Ghi '--- Kiem chung ---' 'Cyan'
    $r = ChayNgoai $apksigner @('verify', '--print-certs', '-v', $FileApk)
    $ra = $r.Ra

    $loi = @()
    if ($ra -match 'Android Debug')                    { $loi += 'VAN dang la khoa go loi' }
    if ($ra -notmatch 'CN=Ngoc Rong Kame')             { $loi += 'khong thay ten Ngoc Rong Kame trong chu ky' }
    if ($ra -notmatch 'v1 scheme \(JAR signing\): true'){ $loi += 'thieu chu ky v1 (may Android cu se tu choi)' }
    if ($ra -notmatch 'v2 scheme[^\r\n]*: true')       { $loi += 'thieu chu ky v2' }
    if ($ra -notmatch 'v3 scheme[^\r\n]*: true')       { $loi += 'thieu chu ky v3' }

    ($ra -split "`r?`n" | Where-Object { $_ -match 'Verified using|certificate DN|SHA-256 digest' }) |
        ForEach-Object { Ghi ("    " + $_.Trim()) 'Gray' }

    if ($loi.Count -gt 0) {
        Write-Host ''
        $loi | ForEach-Object { Ghi "  - $_" 'Red' }
        Dung 'chu ky khong dat yeu cau. KHONG dua file nay cho nguoi choi.'
    }
    Ghi '[OK] Chu ky dung khoa rieng cua ban' 'Green'

} finally {
    # Xoa mat khau khoi bien moi truong du thanh cong hay that bai.
    Remove-Item Env:\MKKHOA -ErrorAction SilentlyContinue
}

# --- 7. Chep ra trang tai ve ---------------------------------------------
if (-not $KhongChep) {
    Write-Host ''
    Ghi '--- Dua len trang tai ve ---' 'Cyan'
    foreach ($d in @(
        @{ TM = $ThuMucWeb; Ten = 'NRO.apk'; Nhan = 'web dang chay tren may nay' },
        @{ TM = $ThuMucGoi; Ten = 'NRO.apk'; Nhan = 'goi de gui len VPS' }
    )) {
        if (-not (Test-Path $d.TM)) { Ghi ("[bo qua] khong co thu muc " + $d.TM) 'Yellow'; continue }
        $dich = Join-Path $d.TM $d.Ten
        Copy-Item $FileApk $dich -Force
        $c = Get-Item $dich
        Ghi ("[OK] {0,-46} {1,7:N1} MB   ({2})" -f $dich, ($c.Length / 1MB), $d.Nhan) 'Green'
    }
}

# --- 8. Nhac viec ---------------------------------------------------------
Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  XONG' -ForegroundColor Green
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host ''
Ghi 'LAM NGAY BAY GIO, DUNG DE MAI:' 'Red'
Ghi ''
Ghi ("  1. Chep file  $FileKhoa") 'Yellow'
Ghi '     ra it nhat HAI noi khac (USB, Google Drive, gui email cho chinh minh)'
Ghi ''
Ghi '  2. Ghi mat khau ra GIAY, cat cung cho voi giay ghi mat khau CSDL'
Ghi ''
Ghi '  3. Cai thu file APK vao dien thoai cua ban, mo game, dang nhap thu'
Ghi '     truoc khi bao cho nguoi choi biet link tai.'
Ghi ''
Ghi '  4. Nhan tro ly: "da ky APK xong" de gui file len VPS'
Ghi ''
Ghi 'Mat file khoa hoac quen mat khau = moi nguoi choi phai go game cai lai.' 'Red'
Write-Host ''
if (-not $KhongChep) { Read-Host '  Bam Enter de dong' }
