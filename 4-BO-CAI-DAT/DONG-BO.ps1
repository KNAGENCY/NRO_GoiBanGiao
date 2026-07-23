# ===========================================================================
#  DONG BO MAY NHA  ->  VPS   (chay TREN VPS)
#
#  Thay cho viec go tung lenh robocopy mot. Bam chay mot lan la xong.
#
#  ---------------------------------------------------------------------
#  BA FILE TUYET DOI KHONG DUOC DONG BO - va vi sao
#  ---------------------------------------------------------------------
#  Ba file duoi day co NOI DUNG KHAC NHAU giua may nha va VPS. O may nha
#  chung dang de mat khau RONG va dia chi 127.0.0.1; tren VPS chung chua
#  mat khau CSDL that va dia chi that. Sao de len la hong ngay:
#
#    zz1\data\config\config.properties
#         may nha: database.pass=(rong), server.sv1=...127.0.0.1...
#         hau qua neu de len: may chu game khong noi duoc CSDL -> cua so
#         bao loi roi tu thoat; va neu co chay duoc thi nguoi choi nhan
#         dia chi 127.0.0.1 -> khong ai vao duoc game.
#
#    htdocs\sv1\config.php        ($DB_PASS = "";)
#         hau qua: TOAN BO web chet, moi trang hien chu do.
#
#    htdocs\sv1\boss_status.php   ('root', '')
#         hau qua: web van vao binh thuong, CHI MAT ANH MAT BOSS o trang
#         chu. Loi bi nuot im lang - day la cho de sot nhat.
#
#  Vi vay ba file nay nam trong danh sach loai tru CUNG, khong co lua chon
#  nao bat duoc chung len. Muon sua chung thi sua thang tren VPS.
#
#  ---------------------------------------------------------------------
#  SCRIPT NAY KHONG XOA GI CA
#  ---------------------------------------------------------------------
#  Khong dung /MIR. Ly do: /MIR se XOA tren VPS moi file khong co o may
#  nha - tuc la xoa sach cac ban luu .goc, file nhat ky, va thu muc data
#  nang 700 MB. Chi mot lan go nham duong dan nguon la mat het.
#  Doi lai: neu ban XOA mot file o may nha, no VAN CON tren VPS. Phai vao
#  xoa tay. Danh doi nay dang gia.
#
#  Chay lai nhieu lan khong sao.
# ===========================================================================

$ErrorActionPreference = 'Stop'
function Ghi($m, $mau = 'White') { Write-Host "  $m" -ForegroundColor $mau }

# Goc o dia C cua MAY NHA, nhin qua Remote Desktop
$NHA = '\\tsclient\C'
$loi = @()
$doi = @()

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  DONG BO MAY NHA -> VPS' -ForegroundColor Cyan
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host ''

# --- Co nhin thay o dia may nha khong ------------------------------------
if (-not (Test-Path $NHA)) {
    Ghi 'DUNG LAI: khong nhin thay o dia may nha.' 'Red'
    Ghi ''
    Ghi 'Nguyen nhan thuong gap:' 'Yellow'
    Ghi '  - Ban dang ngoi trUc tiep tren VPS, khong qua Remote Desktop'
    Ghi '  - Luc mo Remote Desktop chua bat chia se o dia'
    Ghi ''
    Ghi 'Cach bat: o cua so Remote Desktop truoc khi bam Connect ->' 'Yellow'
    Ghi '  the "Local Resources" -> nut "More..." -> tich vao o dia C'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}
Ghi '[OK] Nhin thay o dia may nha' 'Green'

# --- Nhung gi se duoc dong bo --------------------------------------------
# File loai tru: ba file cau hinh chua mat khau + cac ban luu.
$boFile = @(
    'config.properties',      # mat khau CSDL + dia chi may chu
    'config.php',             # mat khau CSDL cua web
    'boss_status.php',        # mat khau CSDL trong dong ket noi
    '*.goc', '*.bak', '*.cu', '*.orig', '*.log', '*.tmp'
)
# Thu muc loai tru: do may chu tu sinh ra, va kho anh nang.
$boThuMuc = @('update_data', 'logs', 'icon', 'img_by_name', 'big_image', '.git')

$viec = @(
    @{ ten = 'May chu game (zz1.jar)'; nguon = "$NHA\Users\Admin\Music\zz1"; dich = 'C:\zz1';               chiFile = 'zz1.jar' },
    @{ ten = 'Ma nguon Java (src)';    nguon = "$NHA\Users\Admin\Music\zz1\src"; dich = 'C:\zz1\src';       deQuy  = $true },
    @{ ten = 'Web (htdocs)';           nguon = "$NHA\xampp\htdocs";              dich = 'C:\xampp\htdocs';  deQuy  = $true },
    @{ ten = 'Bo cai dat';             nguon = "$NHA\Users\Admin\Music\bo-cai-vps"; dich = 'C:\kame-goi\bo-cai-vps'; deQuy = $true }
)

foreach ($v in $viec) {
    Write-Host ''
    Ghi ("--- $($v.ten) ---") 'Cyan'

    if (-not (Test-Path $v.nguon)) {
        Ghi "BO QUA: khong thay $($v.nguon)" 'Yellow'
        $loi += "Khong thay nguon: $($v.ten)"
        continue
    }
    if (-not (Test-Path $v.dich)) { New-Item -ItemType Directory -Force -Path $v.dich | Out-Null }

    # Dung /XO: chi chep file NGUON MOI HON file dich. Tranh chep lai ca
    # 28 MB moi lan, va tranh de nguoc ban moi bang ban cu.
    $doiSo = @($v.nguon, $v.dich)
    if ($v.chiFile) { $doiSo += $v.chiFile }
    $doiSo += @('/XO', '/NP', '/NDL', '/NJH', '/NJS', '/R:2', '/W:2')
    if ($v.deQuy) { $doiSo += '/E' }
    $doiSo += '/XF'; $doiSo += $boFile
    $doiSo += '/XD'; $doiSo += $boThuMuc

    $cu = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    try { $kq = & robocopy @doiSo 2>&1 } finally { $ErrorActionPreference = $cu }
    $ma = $LASTEXITCODE

    # ***********************************************************************
    #  MA THOAT CUA ROBOCOPY KHONG GIONG CHUONG TRINH KHAC
    #  0-7  = THANH CONG (0 la khong co gi moi, 1 la da chep, 2 la co thua...)
    #  8+   = THAT BAI that su
    #  Neu kiem tra kieu "khac 0 la loi" thi lan nao chep duoc file cung bi
    #  bao la loi. Day la cai bay kinh dien.
    # ***********************************************************************
    if ($ma -ge 8) {
        Ghi "LOI khi chep (ma $ma):" 'Red'
        $kq | Select-Object -Last 6 | ForEach-Object { Ghi "    $_" 'Red' }
        $loi += "Chep loi: $($v.ten)"
        continue
    }

    $daChep = $kq | Where-Object { $_ -match '\S' -and $_ -notmatch '^\s*$' }
    if ($ma -eq 0) {
        Ghi '[OK] Khong co gi thay doi' 'Gray'
    } else {
        $n = ($daChep | Measure-Object).Count
        Ghi "[OK] Da cap nhat $n muc" 'Green'
        $daChep | Select-Object -First 12 | ForEach-Object { Ghi ("    " + $_.Trim()) 'DarkGray' }
        if ($n -gt 12) { Ghi "    ... va $($n - 12) muc nua" 'DarkGray' }
        $doi += $v.ten
    }
}

# --- Kiem chung file quan trong nhat -------------------------------------
Write-Host ''
Ghi '--- Kiem chung ---' 'Cyan'

$jarNha = "$NHA\Users\Admin\Music\zz1\zz1.jar"
$jarVps = 'C:\zz1\zz1.jar'
if ((Test-Path $jarNha) -and (Test-Path $jarVps)) {
    # So ma kiem tra chu KHONG so kich thuoc. Da tung gap truong hop robocopy
    # bao chep xong, kich thuoc khop tung byte, ma noi dung van hong.
    $a = (Get-FileHash $jarNha -Algorithm MD5).Hash
    $b = (Get-FileHash $jarVps -Algorithm MD5).Hash
    if ($a -eq $b) { Ghi "[OK] zz1.jar khop  ($($b.Substring(0,16)))" 'Green' }
    else {
        Ghi 'LOI: zz1.jar tren VPS KHONG khop voi may nha!' 'Red'
        Ghi "    may nha: $($a.Substring(0,16))" 'Red'
        Ghi "    VPS    : $($b.Substring(0,16))" 'Red'
        Ghi '    Chay lai script nay mot lan nua.' 'Yellow'
        $loi += 'zz1.jar khong khop'
    }
}

# Ba file cau hinh phai VAN CON mat khau - chung minh da loai tru dung
$kiemTra = @{
    'zz1\data\config\config.properties' = @{ f = 'C:\zz1\data\config\config.properties'; xau = 'database.pass=' }
    'htdocs\sv1\config.php'             = @{ f = 'C:\xampp\htdocs\sv1\config.php';       xau = '$DB_PASS = "";' }
    'htdocs\sv1\boss_status.php'        = @{ f = 'C:\xampp\htdocs\sv1\boss_status.php';  xau = "'root', ''" }
}
foreach ($k in $kiemTra.Keys) {
    $f = $kiemTra[$k].f
    if (Test-Path $f) {
        $nd = Get-Content $f -Raw
        # Voi config.properties: dong "database.pass=" ma KHONG co gi phia sau
        $trong = if ($k -like '*config.properties') { $nd -match '(?m)^database\.pass=\s*$' }
                 else { $nd.Contains($kiemTra[$k].xau) }
        if ($trong) {
            Ghi "LOI NANG: $k da bi de len - MAT KHAU DA MAT!" 'Red'
            $loi += "$k bi de len mat mat khau"
        } else { Ghi "[OK] $k con nguyen mat khau" 'Green' }
    }
}

# --- Neu jar doi ma may chu dang chay ------------------------------------
$doiJar = $doi -contains 'May chu game (zz1.jar)'
$dangChay = $null -ne (Get-Process java -EA SilentlyContinue)

# --- Tong ket ------------------------------------------------------------
Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
if ($loi.Count -eq 0) {
    if ($doi.Count -eq 0) {
        Write-Host '  XONG - KHONG CO GI THAY DOI' -ForegroundColor Green
    } else {
        Write-Host '  XONG - DA DONG BO' -ForegroundColor Green
    }
    Write-Host '===========================================================' -ForegroundColor Cyan
    if ($doi.Count) { Write-Host ''; Ghi 'Da cap nhat:' 'Yellow'; $doi | ForEach-Object { Ghi "  - $_" } }

    if ($doiJar) {
        Write-Host ''
        Ghi 'CHU Y: file may chu game da doi.' 'Yellow'
        if ($dangChay) {
            Ghi 'May chu dang CHAY - no van chay ban CU cho toi khi khoi dong lai.' 'Yellow'
            Ghi 'Muon dung ban moi thi tat va bat lai (nguoi choi se bi rot):' 'Yellow'
            Ghi '    Stop-Process -Name java -Force'
            Ghi '    roi mo C:\zz1 -> chuot phai run.bat -> Run as administrator'
        } else {
            Ghi 'May chu chua chay. Mo C:\zz1 -> chuot phai run.bat -> Run as administrator' 'Yellow'
        }
    }
    if ($doi -contains 'Web (htdocs)') {
        Write-Host ''
        Ghi 'Web da cap nhat - KHONG can khoi dong lai Apache.' 'Gray'
    }
} else {
    Write-Host "  CO $($loi.Count) VAN DE" -ForegroundColor Red
    Write-Host '===========================================================' -ForegroundColor Cyan
    $loi | ForEach-Object { Ghi "  - $_" 'Red' }
    Write-Host ''
    Ghi 'Chup man hinh nay gui lai cho tro ly.' 'Yellow'
}
Write-Host ''
Read-Host '  Bam Enter de dong'
