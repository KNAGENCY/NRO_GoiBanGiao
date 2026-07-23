# ===========================================================================
#  NGOC RONG KAME - BUOC 3: CAI GAME VA WEB TREN VPS
#  ---------------------------------------------------------------------
#  File nay chay TREN VPS, SAU KHI da:
#    - chay xong 1-CAI-MOI-TRUONG.bat (co JDK 21 + XAMPP + dich vu)
#    - giai nen zz1-len-vps.zip vao C:\zz1
#    - giai nen web-len-vps.zip vao C:\xampp\htdocs
#    - chep file z-moi-nhat.sql len VPS
#
#  No lam nhung viec sau:
#    1. Kiem tra moi thu da san sang chua
#    2. Tao co so du lieu "z" va nhap du lieu tu z-moi-nhat.sql
#    3. HOI ban mat khau CSDL moi, roi:
#         - dat cho CA 4 tai khoan MySQL dang de trong
#         - tu dien vao CA 3 file cau hinh (day la cho hay sot nhat)
#    4. Chep config.properties, .htaccess, them cau hinh ten mien cho Apache
#    5. Don vai duong dan cu con sot trong ma nguon web
#    6. Khoi dong lai Apache + MySQL va kiem tra lai tu dau den cuoi
#
#  VE MAT KHAU:
#    Script HOI mat khau luc chay, khong luu no vao bat ky file nhat ky nao,
#    va khong dua no ra dong lenh (de nguoi khac tren may khong doc trom
#    duoc qua danh sach tien trinh). Mat khau chi duoc ghi vao dung 3 file
#    cau hinh bat buoc phai co no.
#
#  CHAY LAI NHIEU LAN DUOC. Nhung neu chay lan 2 thi CSDL da co mat khau,
#  script se hoi lai dung mat khau do.
#
#  CACH CHAY: bam CHUOT PHAI vao file nay -> "Run with PowerShell".
#  Neu Windows chan, mo PowerShell bang quyen Administrator roi go:
#      Set-ExecutionPolicy -Scope Process Bypass -Force
#      C:\duong-dan-toi\3-CAI-GAME-LEN-VPS.ps1
# ===========================================================================

$ErrorActionPreference = 'Stop'

# --- Duong dan mac dinh tren VPS ----------------------------------------
$GAME_ROOT = 'C:\zz1'
$WEB_ROOT  = 'C:\xampp\htdocs'
$MYSQL_BIN = 'C:\xampp\mysql\bin'
$APACHE    = 'C:\xampp\apache'
$THU_MUC   = Split-Path -Parent $MyInvocation.MyCommand.Path
$CAU_HINH  = Join-Path $THU_MUC 'cau-hinh'

$loi = @()   # gom loi lai, in mot the o cuoi

function Ghi($m, $mau = 'White') { Write-Host "  $m" -ForegroundColor $mau }
function Tieu($m) {
    Write-Host ''
    Write-Host "  ================================================" -ForegroundColor Cyan
    Write-Host "   $m" -ForegroundColor Cyan
    Write-Host "  ================================================" -ForegroundColor Cyan
}

# Sao luu mot file thanh <ten>.goc truoc khi sua.
# Chi sao luu LAN DAU - chay lai script khong de ban goc bi ghi de bang
# ban da sua. File .goc khong tai duoc tu ngoai vao vi .htaccess da chan.
function SaoLuu($duongDan) {
    $goc = "$duongDan.goc"
    if ((Test-Path $duongDan) -and -not (Test-Path $goc)) {
        Copy-Item $duongDan $goc -Force
        Ghi "    (da luu ban goc: $(Split-Path $goc -Leaf))" 'DarkGray'
    }
}

# Chay mot cau lenh SQL bang cach DUA QUA DUONG ONG (stdin), khong dat vao
# dong lenh. Ly do: moi thu viet tren dong lenh deu hien ra trong danh sach
# tien trinh cua Windows - ai dang nhap vao may cung doc duoc mat khau.
function ChaySQL($sql, $matKhau = $null) {
    $doiSo = @('-u', 'root', '--default-character-set=utf8mb4')
    if ($matKhau) { $doiSo += "-p$matKhau" }

    # QUAN TRONG: phai ha $ErrorActionPreference xuong 'Continue' quanh loi goi nay.
    # Trong PowerShell 5.1, khi $ErrorActionPreference='Stop' ma lai dung 2>&1 de
    # hung loi cua mot chuong trinh ben ngoai, thi CHI CAN chuong trinh do in ra
    # mot dong loi la PowerShell nem ngoai le va script chet ngay tai cho.
    # Hau qua neu khong sua: chay lai lan thu hai (luc CSDL da co mat khau),
    # mysql tra ve "Access denied" -> script chet, cua so dong ngay, va toan bo
    # cac thong bao loi than thien ben duoi khong bao gio hien ra.
    $cu = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    try {
        $kq = $sql | & "$MYSQL_BIN\mysql.exe" @doiSo 2>&1
    } finally {
        $ErrorActionPreference = $cu
    }
    return $kq
}

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  NGOC RONG KAME - BUOC 3: CAI GAME VA WEB TREN VPS' -ForegroundColor Cyan
Write-Host '===========================================================' -ForegroundColor Cyan


# ===========================================================================
Tieu 'PHAN 1/7 - KIEM TRA MOI THU DA SAN SANG CHUA'
# ===========================================================================

# Thu muc cau-hinh\ la BAT BUOC, khong phai tuy chon.
# Trong do co config.properties da doi san server.sv1 tu 127.0.0.1 sang
# 103.116.38.60. Neu thieu thu muc nay, script se roi vao nhanh du phong
# chi dien moi mat khau va DE NGUYEN 127.0.0.1 -> may chu van len, web van
# chay, nhung dien thoai nguoi choi nhan duoc dia chi 127.0.0.1 se tu goi vao
# chinh no => KHONG MOT AI VAO DUOC GAME, ma khong co bat ky bao loi nao.
if (-not (Test-Path $CAU_HINH)) {
    Ghi 'DUNG LAI: khong thay thu muc cau-hinh canh script nay.' 'Red'
    Ghi "Dang tim o: $CAU_HINH" 'Gray'
    Ghi 'Phai chep NGUYEN thu muc bo-cai-vps len VPS, khong chi rieng file .ps1.' 'Yellow'
    Read-Host "`n  Bam Enter de thoat"
    exit 1
}

$quyen = ([Security.Principal.WindowsPrincipal] `
          [Security.Principal.WindowsIdentity]::GetCurrent()
         ).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $quyen) {
    Ghi 'DUNG LAI: phai chay bang quyen Administrator.' 'Red'
    Ghi 'Bam chuot phai vao PowerShell -> Run as administrator, roi chay lai.' 'Yellow'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}
Ghi '[OK] Dang chay bang quyen Administrator' 'Green'

# --- Thu muc game ---
if (-not (Test-Path "$GAME_ROOT\zz1.jar")) {
    Ghi "DUNG LAI: khong thay $GAME_ROOT\zz1.jar" 'Red'
    Ghi ''
    Ghi 'Ban da giai nen zz1-len-vps.zip vao dung cho chua?' 'Yellow'
    Ghi "Cho dung la: $GAME_ROOT" 'Yellow'
    Ghi ''
    Ghi 'Vi sao phai dung dung cho nay: 4 file PHP cua trang quan tri co san'
    Ghi 'nhanh du phong tro toi day. Dat cho khac thi phai sua tay 5 cho.'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}
Ghi "[OK] Thu muc game: $GAME_ROOT" 'Green'

# src\ la bat buoc: trang quan tri doc thang ma nguon Java bang regex
# de hien ti le nang sao Pha Le. Thieu src\ thi trang do hien toan so 0.
# Duong dan DUNG la src\models\Combine\CombineService.java - kiem chung bang
# cach doi chieu voi game_bridge.php dong 217, noi web that su doc file nay:
#     panel_java_source_file('models/Combine/CombineService.java')
# Ban dau script tim o src\server\ nen bao thieu file trong khi file van co.
if (-not (Test-Path "$GAME_ROOT\src\models\Combine\CombineService.java")) {
    Ghi 'CANH BAO: thieu src\models\Combine\CombineService.java.' 'Red'
    Ghi 'Trang quan tri "ti le nang sao" se hien toan so 0.' 'Yellow'
    $loi += 'Thieu thu muc src (trang quan tri nang sao se hong)'
}
if (-not (Test-Path "$GAME_ROOT\src\server\icon.png")) {
    Ghi 'CANH BAO: thieu src\server\icon.png - server se vang khi khoi dong.' 'Red'
    $loi += 'Thieu src\server\icon.png'
}

# --- Thu muc web ---
if (-not (Test-Path "$WEB_ROOT\sv1\config.php")) {
    Ghi "DUNG LAI: khong thay $WEB_ROOT\sv1\config.php" 'Red'
    Ghi 'Ban da giai nen web-len-vps.zip vao C:\xampp\htdocs chua?' 'Yellow'
    Ghi 'CHU Y: ten thu muc phai dung la "sv1" - doi ten la dut het lien ket.' 'Yellow'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}
Ghi "[OK] Thu muc web: $WEB_ROOT\sv1" 'Green'

# --- Kiem tra rac XAMPP con sot khong ---
foreach ($rac in @('dashboard', 'applications.html', 'bitnami.css')) {
    if (Test-Path "$WEB_ROOT\$rac") {
        Ghi "CANH BAO: con rac cua XAMPP: $rac" 'Yellow'
        Ghi '  Trang http://nrokame.com/dashboard/ khoe ra phien ban XAMPP/PHP' 'Yellow'
        Ghi '  cho ke xau biet duong do lo hong. Nen xoa.' 'Yellow'
        $loi += "Con rac XAMPP trong htdocs: $rac"
    }
}

# --- MySQL ---
$dvMysql = Get-Service -Name 'mysql' -EA SilentlyContinue
if (-not $dvMysql) {
    Ghi 'DUNG LAI: chua co dich vu Windows ten "mysql".' 'Red'
    Ghi 'Hay chay lai 1-CAI-MOI-TRUONG.bat bang quyen Administrator.' 'Yellow'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}
if ($dvMysql.Status -ne 'Running') {
    Ghi 'MySQL chua chay, dang bat...' 'Yellow'
    Start-Service 'mysql'
    Start-Sleep -Seconds 5
}
Ghi "[OK] Dich vu MySQL: $((Get-Service 'mysql').Status)" 'Green'

# --- File ban sao CSDL ---
$fileDump = Join-Path $THU_MUC 'z-moi-nhat.sql'
if (-not (Test-Path $fileDump)) {
    Ghi ''
    Ghi "Khong thay file z-moi-nhat.sql canh script nay." 'Yellow'
    $fileDump = Read-Host '  Go duong dan day du toi file z-moi-nhat.sql'
    $fileDump = $fileDump.Trim('"', ' ')
}
if (-not (Test-Path $fileDump)) {
    Ghi 'DUNG LAI: khong tim thay file ban sao co so du lieu.' 'Red'
    Ghi 'Chay file 2-DONG-GOI-O-MAY-NHA.ps1 o may nha de tao ra no.' 'Yellow'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}
Ghi "[OK] Ban sao CSDL: $fileDump" 'Green'

# TUYET DOI khong dung file z.sql cu (47 bang, thieu bang news, mat khau
# chua ma hoa, link tai game van ten cu).
if ((Split-Path $fileDump -Leaf) -eq 'z.sql') {
    Ghi ''
    Ghi 'DUNG LAI: day la file z.sql - ban CU, khong dung duoc.' 'Red'
    Ghi 'No thieu bang "news", mat khau nguoi choi chua ma hoa, va cac link' 'Yellow'
    Ghi 'tai game van la ten thuong hieu cu. Phai dung z-moi-nhat.sql.' 'Yellow'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}


# ===========================================================================
Tieu 'PHAN 2/7 - TAO VA NHAP CO SO DU LIEU'
# ===========================================================================

# Kiem tra xem MySQL da co mat khau chua (chay lan 2 tro di)
$daCoMatKhau = $false
$thu = ChaySQL 'SELECT 1;'
if ($LASTEXITCODE -ne 0) {
    $daCoMatKhau = $true
    Ghi 'CSDL da co mat khau tu lan chay truoc.' 'Yellow'
}

$matKhauCu = $null
if ($daCoMatKhau) {
    $sec = Read-Host '  Go mat khau CSDL hien tai' -AsSecureString
    $matKhauCu = [Runtime.InteropServices.Marshal]::PtrToStringAuto(
        [Runtime.InteropServices.Marshal]::SecureStringToBSTR($sec))
    $thu = ChaySQL 'SELECT 1;' $matKhauCu
    if ($LASTEXITCODE -ne 0) {
        Ghi 'DUNG LAI: mat khau khong dung.' 'Red'
        Read-Host "`n  Bam Enter de thoat"; exit 1
    }
    Ghi '[OK] Mat khau dung.' 'Green'
}

# Da co bang nao trong "z" chua?
$demBang = ChaySQL "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='z';" $matKhauCu
$soBangHienCo = 0
if ($LASTEXITCODE -eq 0) {
    $soBangHienCo = [int](($demBang | Select-Object -Last 1).Trim())
}

if ($soBangHienCo -ge 40) {
    Ghi "CSDL 'z' da co $soBangHienCo bang - co ve da nhap tu lan truoc." 'Yellow'
    Ghi ''
    Ghi 'Nhap de len se XOA SACH du lieu nguoi choi hien co.' 'Red'
    $traLoi = Read-Host '  Van muon nhap de len? Go: XOA VA NHAP LAI (hoac Enter de bo qua)'
    # -ceq chu KHONG phai -eq: toan tu -eq cua PowerShell khong phan biet hoa
    # thuong, nen go "xoa va nhap lai" cung xoa sach CSDL. Bat go chu HOA la
    # de nguoi dung phai dung lai suy nghi - dung -eq thi mat het y nghia do.
    if ($traLoi -ceq 'XOA VA NHAP LAI') {
        ChaySQL 'DROP DATABASE IF EXISTS z;' $matKhauCu | Out-Null
        $soBangHienCo = 0
    } else {
        Ghi 'Bo qua buoc nhap, giu nguyen du lieu dang co.' 'Green'
    }
}

if ($soBangHienCo -lt 40) {
    Ghi 'Dang tao co so du lieu "z"...'
    ChaySQL 'CREATE DATABASE IF NOT EXISTS z CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;' $matKhauCu | Out-Null

    # CHU Y: KHONG dung dau "<" o day.
    # PowerShell KHONG CO toan tu "<" - go vao la bao loi ngay:
    #     The '<' operator is reserved for future use.
    # (Lenh do chi dung trong cua so cmd doi kieu cu.) Cach dung duoc trong
    # PowerShell la nho chinh mysql doc file bang lenh "source".
    # Duong dan phai doi \ thanh / thi mysql moi hieu.
    #
    # VA: "source" PHAI DUNG DAU MOT DONG RIENG.
    # Viet "USE z; source duong-dan" tren cung mot dong thi mysql chi nhan ra
    # "USE z;", con phan "source ..." con lai bi coi la cau lenh SQL binh thuong
    # va gui thang len may chu -> ERROR 1064 syntax error. Da tai hien va kiem
    # chung ca hai cach o may nha 2026-07-21.
    # Vi sao dung "source" chu khong day noi dung file qua ong dan: file .sql
    # nang 845 KB va co chu tieng Viet co dau. Day qua ong dan cua PowerShell
    # thi chu bi doi bang ma va hong dau. De mysql tu doc file thi khong sao.
    $duongDanSql = $fileDump -replace '\\', '/'
    Ghi 'Dang nhap du lieu (mat khoang 10-60 giay)...'
    $ketQua = ChaySQL "USE z;`nsource $duongDanSql`n" $matKhauCu

    if ($LASTEXITCODE -ne 0) {
        Ghi 'LOI khi nhap du lieu:' 'Red'
        $ketQua | Select-Object -First 10 | ForEach-Object { Ghi "  $_" 'Red' }
        Read-Host "`n  Bam Enter de thoat"; exit 1
    }
}

# Kiem tra so bang - phai la 48
$demBang = ChaySQL "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='z';" $matKhauCu
$soBang = [int](($demBang | Select-Object -Last 1).Trim())
if ($soBang -ge 48) {
    Ghi "[OK] Co so du lieu 'z' co $soBang bang." 'Green'
} else {
    Ghi "CANH BAO: chi co $soBang bang, dang le phai co 48." 'Red'
    $loi += "CSDL chi co $soBang bang thay vi 48"
}


# ===========================================================================
Tieu 'PHAN 3/7 - DAT MAT KHAU CO SO DU LIEU'
# ===========================================================================

if ($daCoMatKhau) {
    Ghi 'Bo qua - CSDL da co mat khau roi.' 'Green'
    $matKhau = $matKhauCu
} else {

    Ghi 'Hien gio CSDL dang KHONG CO MAT KHAU - ai vao duoc may la vao duoc' 'Yellow'
    Ghi 'toan bo du lieu nguoi choi. Bay gio dat mat khau.' 'Yellow'
    Ghi ''
    Ghi 'Mat khau nay chi ban va may chu dung, KHONG phai mat khau nguoi choi.'
    Ghi 'Ban se khong phai go lai no hang ngay - script tu dien vao 3 file.'
    Ghi ''
    Ghi 'Quy dinh (de khong lam hong file cau hinh):' 'Yellow'
    Ghi '  - it nhat 12 ky tu'
    Ghi '  - chi dung: chu cai, chu so, va cac dau  - _ . ! @ # % ^ * + ='
    Ghi '  - KHONG dung dau nhay  '' "  dau xuyet nguoc \  hay dau  $'
    Ghi '    (nhung dau nay co y nghia rieng trong file cau hinh PHP va Java,'
    Ghi '     dung vao la file hong ma khong bao loi ro rang)'
    Ghi ''
    Ghi 'Goi y: ghep 3 tu de nho + so, vi du  KameRongThieng2026'
    Ghi ''

    $matKhau = $null
    for ($lan = 1; $lan -le 5; $lan++) {
        $s1 = Read-Host '  Go mat khau moi' -AsSecureString
        $p1 = [Runtime.InteropServices.Marshal]::PtrToStringAuto(
              [Runtime.InteropServices.Marshal]::SecureStringToBSTR($s1))
        $s2 = Read-Host '  Go lai lan nua cho chac' -AsSecureString
        $p2 = [Runtime.InteropServices.Marshal]::PtrToStringAuto(
              [Runtime.InteropServices.Marshal]::SecureStringToBSTR($s2))

        if ($p1 -ne $p2)              { Ghi 'Hai lan go khong giong nhau. Thu lai.' 'Red'; continue }
        if ($p1.Length -lt 12)        { Ghi 'Ngan qua - can it nhat 12 ky tu. Thu lai.' 'Red'; continue }
        if ($p1 -notmatch '^[A-Za-z0-9\-_.!@#%^*+=]+$') {
            Ghi 'Co ky tu khong duoc phep. Chi dung chu, so va  - _ . ! @ # % ^ * + =' 'Red'; continue
        }
        $matKhau = $p1
        break
    }

    if (-not $matKhau) {
        Ghi 'DUNG LAI: khong dat duoc mat khau sau 5 lan thu.' 'Red'
        Read-Host "`n  Bam Enter de thoat"; exit 1
    }

    # ***********************************************************************
    #  PHAI DAT CHO CA 4 TAI KHOAN, KHONG PHAI 1
    #  ---------------------------------------------------------------------
    #  MySQL coi "root dang nhap tu localhost" va "root dang nhap tu
    #  127.0.0.1" la HAI TAI KHOAN KHAC NHAU. Tren may nay dang co 4 tai
    #  khoan deu de trong mat khau:
    #      root@localhost   root@127.0.0.1   root@::1   pma@localhost
    #
    #  Cach lam pho bien nhat (mysqladmin -u root password ...) chi doi
    #  DUNG MOT DONG. Ma tren Windows, ca Java lan PHP deu noi vao qua
    #  127.0.0.1 -> khop voi root@127.0.0.1, von VAN CON DE TRONG.
    #  Ket qua: web va game van chay ngon lanh, ban tin la da co mat khau,
    #  trong khi CSDL van mo toang. Day la kieu nguy hiem am tham nhat.
    # ***********************************************************************
    Ghi ''
    Ghi 'Dang dat mat khau cho tat ca tai khoan CSDL...'

    $sqlDat = @"
SET PASSWORD FOR 'root'@'localhost' = PASSWORD('$matKhau');
SET PASSWORD FOR 'root'@'127.0.0.1' = PASSWORD('$matKhau');
SET PASSWORD FOR 'root'@'::1'       = PASSWORD('$matKhau');
FLUSH PRIVILEGES;
"@
    $kq = ChaySQL $sqlDat
    if ($LASTEXITCODE -ne 0) {
        Ghi 'LOI khi dat mat khau:' 'Red'
        $kq | Select-Object -First 5 | ForEach-Object { Ghi "  $_" 'Red' }
        Read-Host "`n  Bam Enter de thoat"; exit 1
    }

    # Tai khoan pma la cua phpMyAdmin. Khong ai nhac den no, nhung no cung
    # de trong mat khau. phpMyAdmin da bi chan chi cho vao tu chinh may nay
    # (dong "Require local" trong httpd-xampp.conf) nen it nguy hiem hon,
    # nhung xoa han cho gon vi ban khong dung phpMyAdmin qua mang.
    ChaySQL "DROP USER IF EXISTS 'pma'@'localhost'; FLUSH PRIVILEGES;" $matKhau | Out-Null
    Ghi '  Da xoa tai khoan phu "pma" (cua phpMyAdmin, khong dung den).' 'Gray'

    # Kiem chung that: con tai khoan nao de trong mat khau khong?
    $conTrong = ChaySQL "SELECT CONCAT(user,'@',host) FROM mysql.user WHERE authentication_string='' AND password='';" $matKhau
    $dsTrong = $conTrong | Select-Object -Skip 1 | Where-Object { $_.Trim() -ne '' }
    if ($dsTrong) {
        Ghi 'CANH BAO: van con tai khoan CSDL khong co mat khau:' 'Red'
        $dsTrong | ForEach-Object { Ghi "    $_" 'Red' }
        $loi += 'Con tai khoan CSDL khong co mat khau'
    } else {
        Ghi '[OK] Moi tai khoan CSDL deu da co mat khau.' 'Green'
    }
}


# ===========================================================================
Tieu 'PHAN 4/7 - DIEN MAT KHAU VAO 3 FILE CAU HINH'
# ===========================================================================

# Thieu MOT trong ba cho la hong, va TRIEU CHUNG KHAC HAN NHAU nen rat de
# chan doan nham:
#   Sot (1) config.properties : may chu game khong noi duoc CSDL -> cua so
#       bao loi roi tu thoat -> run.bat chay lai vo tan, cua so nhay qua
#       nhanh de doc kip loi.
#   Sot (2) config.php        : toan bo web chet, moi trang hien chu do.
#   Sot (3) boss_status.php   : web van vao binh thuong, CHI MAT ANH MAT
#       BOSS o trang chu. Loi bi nuot im lang trong khoi catch. Day la cho
#       de sot nhat vi mat khau nam lot thom giua dong lenh.

$soChoDaSua = 0

# --- (1) config.properties cua may chu game ---
$fileProp = "$GAME_ROOT\data\config\config.properties"
$mauProp  = Join-Path $CAU_HINH 'config.properties'

if (Test-Path $mauProp) {
    SaoLuu $fileProp
    $noiDung = Get-Content $mauProp -Raw
    $noiDung = $noiDung -replace '__MAT_KHAU_MYSQL__', $matKhau
    # Ghi khong dau BOM: Java doc file .properties bang bang ma ASCII/Latin,
    # co BOM thi khoa dau tien bi hong ten.
    [IO.File]::WriteAllText($fileProp, $noiDung, (New-Object Text.UTF8Encoding $false))
    Ghi "[OK] (1/3) Da chep va dien mat khau: $fileProp" 'Green'
    $soChoDaSua++
} elseif (Test-Path $fileProp) {
    SaoLuu $fileProp
    $noiDung = Get-Content $fileProp -Raw
    $noiDung = $noiDung -replace '(?m)^database\.pass=.*$', "database.pass=$matKhau"
    [IO.File]::WriteAllText($fileProp, $noiDung, (New-Object Text.UTF8Encoding $false))
    Ghi "[OK] (1/3) Da dien mat khau vao: $fileProp" 'Green'
    $soChoDaSua++
} else {
    Ghi "LOI: khong thay $fileProp" 'Red'
    $loi += 'Khong dien duoc mat khau vao config.properties'
}

# --- (2) sv1\config.php ---
$fileCfg = "$WEB_ROOT\sv1\config.php"
if (Test-Path $fileCfg) {
    SaoLuu $fileCfg
    $noiDung = Get-Content $fileCfg -Raw
    $moi = $noiDung -replace '(?m)^\$DB_PASS\s*=\s*".*?";', "`$DB_PASS = `"$matKhau`";"
    if ($moi -eq $noiDung) {
        Ghi 'CANH BAO: khong tim thay dong $DB_PASS trong config.php.' 'Red'
        $loi += 'Khong dien duoc mat khau vao sv1\config.php'
    } else {
        [IO.File]::WriteAllText($fileCfg, $moi, (New-Object Text.UTF8Encoding $false))
        Ghi "[OK] (2/3) Da dien mat khau vao: $fileCfg" 'Green'
        $soChoDaSua++
    }
}

# --- (3) sv1\boss_status.php ---
$fileBoss = "$WEB_ROOT\sv1\boss_status.php"
if (Test-Path $fileBoss) {
    SaoLuu $fileBoss
    $noiDung = Get-Content $fileBoss -Raw
    # Dong that trong file:
    #   $pdo = new PDO('mysql:host=localhost;dbname=z;charset=utf8mb4', 'root', '');
    $moi = $noiDung -replace "('mysql:host=localhost;dbname=z;charset=utf8mb4'\s*,\s*'root'\s*,\s*)'[^']*'", "`${1}'$matKhau'"
    if ($moi -eq $noiDung) {
        Ghi 'CANH BAO: khong tim thay dong ket noi CSDL trong boss_status.php.' 'Red'
        Ghi '  Hau qua: web van chay, chi MAT ANH MAT BOSS o trang chu.' 'Yellow'
        $loi += 'Khong dien duoc mat khau vao sv1\boss_status.php'
    } else {
        [IO.File]::WriteAllText($fileBoss, $moi, (New-Object Text.UTF8Encoding $false))
        Ghi "[OK] (3/3) Da dien mat khau vao: $fileBoss" 'Green'
        $soChoDaSua++
    }
}

if ($soChoDaSua -eq 3) {
    Ghi ''
    Ghi '[OK] Da dien du CA 3 noi. Day la cho hay sot nhat khi lam tay.' 'Green'
} else {
    Ghi ''
    Ghi "CANH BAO: chi dien duoc $soChoDaSua/3 noi. Xem canh bao o tren." 'Red'
}


# ===========================================================================
Tieu 'PHAN 5/7 - CAU HINH APACHE VA WEB'
# ===========================================================================

# --- .htaccess o goc htdocs ---
$mauHt = Join-Path $CAU_HINH 'htaccess-vps.txt'
$dichHt = "$WEB_ROOT\.htaccess"
if (Test-Path $mauHt) {
    SaoLuu $dichHt
    Copy-Item $mauHt $dichHt -Force
    Ghi "[OK] Da dat $dichHt" 'Green'
} else {
    Ghi "CANH BAO: khong thay $mauHt" 'Yellow'
    $loi += 'Thieu file htaccess-vps.txt'
}

# --- Them cau hinh ten mien vao httpd-vhosts.conf ---
$fileVhost = "$APACHE\conf\extra\httpd-vhosts.conf"
$mauVhost  = Join-Path $CAU_HINH 'nrokame-vhost.conf'
if ((Test-Path $fileVhost) -and (Test-Path $mauVhost)) {
    $dangCo = Get-Content $fileVhost -Raw
    if ($dangCo -match 'nrokame\.com') {
        Ghi '[BO QUA] Cau hinh ten mien da co trong httpd-vhosts.conf' 'Yellow'
    } else {
        # Chi lay phan nam giua hai DONG danh dau trong file mau.
        #
        # BAY DA DINH THAT (2026-07-21). Bieu thuc tim cu la:
        #     '(?s)=== BAT DAU CHEP ===(.*?)=== HET CHEP ==='
        # No vo phai DONG HUONG DAN o dau file mau - dong do co nhac ten hai
        # cai moc nay trong ngoac kep:
        #     nam giua hai dong "=== BAT DAU CHEP ===" va "=== HET CHEP ==="
        # nen thu cat ra duoc dung hai chu   " va "   roi ghim vao Apache.
        # Hau qua KEP va rat kho doan:
        #   - Apache gap dong rac  " va "  -> bao sai cu phap -> KHONG chay
        #   - cau hinh ten mien that thi CHUA HE duoc them vao
        # Cach sua: bat buoc khop CA DONG (phai co dau # dau dong va day dau =),
        # nho hai neo ^ $ cung co (?m).
        $mau = Get-Content $mauVhost -Raw
        $moc = '(?sm)^#\s*=+\s*BAT DAU CHEP\s*=+\s*\r?$(.*?)^#\s*=+\s*HET CHEP\s*=+\s*\r?$'
        if ($mau -notmatch $moc) {
            Ghi 'LOI: khong doc duoc khoi cau hinh trong nrokame-vhost.conf' 'Red'
            $loi += 'Khong them duoc cau hinh ten mien cho Apache'
        } else {
            $phan = $Matches[1].Trim()
            SaoLuu $fileVhost
            Add-Content -Path $fileVhost -Value "`r`n$phan`r`n" -Encoding UTF8
            Ghi '[OK] Da them cau hinh ten mien nrokame.com vao Apache' 'Green'
        }
    }
}

# --- Sua ServerName trong httpd.conf ---
$fileHttpd = "$APACHE\conf\httpd.conf"
if (Test-Path $fileHttpd) {
    $noiDung = Get-Content $fileHttpd -Raw
    if ($noiDung -match '(?m)^ServerName\s+nrokame\.com') {
        Ghi '[BO QUA] ServerName da la nrokame.com' 'Yellow'
    } else {
        SaoLuu $fileHttpd
        $moi = $noiDung -replace '(?m)^ServerName\s+localhost:80', 'ServerName nrokame.com:80'
        if ($moi -ne $noiDung) {
            [IO.File]::WriteAllText($fileHttpd, $moi, (New-Object Text.UTF8Encoding $false))
            Ghi '[OK] Da doi ServerName thanh nrokame.com:80' 'Green'
        } else {
            Ghi 'CANH BAO: khong tim thay dong "ServerName localhost:80" de sua.' 'Yellow'
        }
    }
}

# --- Don duong dan cu con sot trong game_bridge.php ---
# Dong 131 co hai duong dan cua may nha: ...\Admin\Music\zz1 va ...\Admin\Music\j
# Tren VPS hai thu muc nay khong ton tai nen vong lap khong them gi -> KHONG
# gay hong. Nhung don di cho sach, de nguoi sau (ke ca AI) khong tuong day
# la cau hinh that.
$fileBridge = "$WEB_ROOT\sv1\admin\game_bridge.php"
if (Test-Path $fileBridge) {
    $noiDung = Get-Content $fileBridge -Raw
    # Ghi chu ve so dau xuyet nguoc: trong file PHP, duong dan viet la
    #   'C:\\Users\\Admin\\Music\\zz1'   (hai dau xuyet, vi PHP cung phai
    #                                     nhan doi chung len)
    # De tim hai dau xuyet do bang bieu thuc tim kiem cua PowerShell thi
    # phai go BON dau. Go tam dau hay nhieu hon la khong khop gi ca va
    # script se lang le bo qua. (Da kiem chung bang cach chay that.)
    if ($noiDung -match "C:\\\\Users\\\\Admin\\\\Music") {
        SaoLuu $fileBridge
        $moi = $noiDung -replace "\[\s*'C:\\\\Users\\\\Admin\\\\Music\\\\zz1'\s*,\s*'C:\\\\Users\\\\Admin\\\\Music\\\\j'\s*\]", '[]'
        if ($moi -ne $noiDung) {
            [IO.File]::WriteAllText($fileBridge, $moi, (New-Object Text.UTF8Encoding $false))
            Ghi '[OK] Da don duong dan may nha con sot trong game_bridge.php' 'Green'
        } else {
            Ghi 'Ghi chu: game_bridge.php con duong dan may nha, nhung khong gay hong.' 'Gray'
        }
    }
}

# --- Bat bind-address de MariaDB chi nghe tu chinh may nay ---
$fileMyIni = "$MYSQL_BIN\my.ini"
if (Test-Path $fileMyIni) {
    $noiDung = Get-Content $fileMyIni -Raw
    if ($noiDung -match '(?m)^\s*bind-address') {
        Ghi '[BO QUA] bind-address da bat san' 'Yellow'
    } elseif ($noiDung -match '(?m)^#\s*bind-address\s*=\s*"127\.0\.0\.1"') {
        SaoLuu $fileMyIni
        $moi = $noiDung -replace '(?m)^#(\s*bind-address\s*=\s*"127\.0\.0\.1")', '$1'
        [IO.File]::WriteAllText($fileMyIni, $moi, (New-Object Text.UTF8Encoding $false))
        Ghi '[OK] Da bat bind-address - MariaDB chi nhan ket noi tu chinh may nay.' 'Green'
        Ghi '     (Lop khoa thu hai, phong khi tuong lua bi tat.)' 'Gray'

        # Da bat bind-address thi PHAI doi localhost -> 127.0.0.1 trong
        # config.properties. Ly do: tren Windows "localhost" doi khi bi dich
        # thanh dia chi IPv6 ::1, ma MariaDB luc nay chi nghe 127.0.0.1
        # -> may chu game bao loi khong noi duoc CSDL.
        if (Test-Path $fileProp) {
            $p = Get-Content $fileProp -Raw
            $p2 = $p -replace '(?m)^database\.host=localhost\s*$', 'database.host=127.0.0.1'
            if ($p2 -ne $p) {
                [IO.File]::WriteAllText($fileProp, $p2, (New-Object Text.UTF8Encoding $false))
                Ghi '     Da doi database.host thanh 127.0.0.1 cho khop.' 'Gray'
            }
        }
    }
}

# --- Hoi ve che do go loi ---
if (Test-Path $fileProp) {
    $p = Get-Content $fileProp -Raw
    if ($p -match '(?m)^server\.debug=true') {
        Ghi ''
        Ghi 'Che do go loi (server.debug) dang BAT.' 'Yellow'
        Ghi 'Bat thi server in rat nhieu thong tin noi bo ra file nhat ky.'
        Ghi 'De lau ngay, file nhat ky co the phinh to an het o C: 70 GB.'
        Ghi 'Ngay mo cua nen TAT.'
        $tl = Read-Host '  Tat che do go loi bay gio? (co/khong)'
        if ($tl -eq 'co') {
            $p2 = $p -replace '(?m)^server\.debug=true', 'server.debug=false'
            [IO.File]::WriteAllText($fileProp, $p2, (New-Object Text.UTF8Encoding $false))
            Ghi '[OK] Da tat che do go loi.' 'Green'
        } else {
            Ghi 'Giu nguyen. Nho tat truoc ngay mo cua.' 'Yellow'
            $loi += 'server.debug van dang bat - nho tat truoc ngay mo cua'
        }
    }
}


# ===========================================================================
Tieu 'PHAN 6/7 - KHOI DONG LAI APACHE VA MYSQL'
# ===========================================================================

# --- Thu cu phap cau hinh Apache TRUOC khi khoi dong lai ---
# Vi sao can buoc nay: neu cau hinh sai mot dong, Restart-Service chi bao
# cut lut "Stopped" ma KHONG he noi sai o dau. Nguoi dung nhin man hinh
# khong hieu gi, phai di doc file nhat ky moi ra. Lenh "httpd -t" in thang
# ra so dong va noi dung dong sai.
# Da mat mot vong di lai vi thieu buoc nay (2026-07-21).
# Luu y: httpd -t in ket qua ra luong LOI chu khong phai luong thuong, nen
# phai gop 2>&1 va tam ha $ErrorActionPreference thi moi doc duoc.
$httpdExe = "$APACHE\bin\httpd.exe"
if (Test-Path $httpdExe) {
    $cu = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    try { $kqCauHinh = & $httpdExe -t 2>&1 } finally { $ErrorActionPreference = $cu }
    if ($kqCauHinh -match 'Syntax OK') {
        Ghi '[OK] Cau hinh Apache dung cu phap' 'Green'
    } else {
        Ghi 'LOI CU PHAP trong cau hinh Apache - Apache se KHONG chay duoc:' 'Red'
        $kqCauHinh | ForEach-Object { Ghi "  $_" 'Red' }
        $loi += 'Cau hinh Apache sai cu phap'
    }
}

foreach ($ten in @('Apache2.4', 'mysql')) {
    $dv = Get-Service -Name $ten -EA SilentlyContinue
    if ($dv) {
        Ghi "Dang khoi dong lai $ten..."
        Restart-Service $ten -Force -EA SilentlyContinue
        Start-Sleep -Seconds 3
        $tt = (Get-Service $ten).Status
        if ($tt -eq 'Running') { Ghi "[OK] $ten : $tt" 'Green' }
        else {
            Ghi "CANH BAO: $ten : $tt" 'Red'
            $loi += "$ten khong chay"
        }
    } else {
        Ghi "CANH BAO: khong co dich vu $ten" 'Red'
        $loi += "Thieu dich vu $ten"
    }
}


# ===========================================================================
Tieu 'PHAN 7/7 - KIEM TRA LAI TU DAU DEN CUOI'
# ===========================================================================

# 1. Web co tra loi khong
try {
    $r = Invoke-WebRequest -Uri 'http://127.0.0.1/' -UseBasicParsing -TimeoutSec 15
    Ghi "[OK] Web tra loi ma $($r.StatusCode)" 'Green'
} catch {
    Ghi "CANH BAO: web khong tra loi ($($_.Exception.Message))" 'Red'
    $loi += 'Web khong tra loi tren cong 80'
}

# 2. Trang game co vao duoc khong
try {
    $r = Invoke-WebRequest -Uri 'http://127.0.0.1/sv1/trangchu' -UseBasicParsing -TimeoutSec 15
    Ghi "[OK] Trang game /sv1/trangchu tra loi ma $($r.StatusCode)" 'Green'
} catch {
    Ghi 'CANH BAO: /sv1/trangchu khong vao duoc.' 'Red'
    Ghi '  Thuong la do Apache chua doc file .htaccess (duong dan sach hong).' 'Yellow'
    $loi += 'Trang /sv1/trangchu khong vao duoc'
}

# 2b. Dia chi may chu gui xuong cho nguoi choi da dung chua.
# Day la loi IM LANG nguy hiem nhat: moi thu khac deu bao OK, cong 14445 mo,
# server chay ngon - nhung neu server.sv1 van la 127.0.0.1 thi dien thoai
# nguoi choi nhan dia chi do se tu goi vao chinh no => khong ai vao duoc game.
if (Test-Path $fileProp) {
    $noiDungProp = Get-Content $fileProp -Raw
    if ($noiDungProp -match '(?m)^\s*server\.sv1\s*=.*127\.0\.0\.1') {
        Ghi 'LOI NANG: server.sv1 van dang la 127.0.0.1 - KHONG AI VAO DUOC GAME!' 'Red'
        Ghi "  Sua file: $fileProp" 'Yellow'
        Ghi '  Dong server.sv1 phai la: Ngoc Rong Kame:103.116.38.60:14445:0,0,0' 'Yellow'
        $loi += 'server.sv1 van la 127.0.0.1 - nguoi choi khong ket noi duoc'
    } elseif ($noiDungProp -match '(?m)^\s*server\.sv1\s*=.*103\.116\.38\.60') {
        Ghi '[OK] server.sv1 da tro dung IP VPS 103.116.38.60' 'Green'
    } else {
        Ghi 'CANH BAO: khong doc duoc dong server.sv1 - hay mo file kiem tra bang tay.' 'Yellow'
        $loi += 'Khong xac minh duoc server.sv1'
    }
}

# 3. Mat khau CSDL da vao 3 file chua - kiem tra that, khong tin loi script
$kiemTra = @{
    'config.properties'  = @{ F = $fileProp;  M = '__MAT_KHAU_MYSQL__' }
    'sv1\config.php'     = @{ F = $fileCfg;   M = '$DB_PASS = "";'     }
    'sv1\boss_status.php'= @{ F = $fileBoss;  M = "'root', ''"         }
}
foreach ($k in $kiemTra.Keys) {
    $f = $kiemTra[$k].F
    if (Test-Path $f) {
        $nd = Get-Content $f -Raw
        if ($nd.Contains($kiemTra[$k].M)) {
            Ghi "CANH BAO: $k VAN CHUA co mat khau!" 'Red'
            $loi += "$k chua co mat khau"
        } else {
            Ghi "[OK] $k da co mat khau" 'Green'
        }
    }
}

# 4. Cong 3306 phai KHONG mo ra ngoai
$luat3306 = netsh advfirewall firewall show rule name="Kame - Chan MySQL 3306" 2>&1
if ($LASTEXITCODE -eq 0) {
    Ghi '[OK] Cong 3306 dang bi chan khong cho vao tu internet' 'Green'
} else {
    Ghi 'CANH BAO: khong thay luat chan cong 3306!' 'Red'
    $loi += 'Thieu luat chan cong 3306'
}


# ===========================================================================
#  TONG KET
# ===========================================================================
Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
if ($loi.Count -eq 0) {
    Write-Host '  XONG BUOC 3 - MOI THU DEU OK' -ForegroundColor Green
} else {
    Write-Host '  XONG BUOC 3 - NHUNG CO VIEC CHUA TRON' -ForegroundColor Yellow
}
Write-Host '===========================================================' -ForegroundColor Cyan
Ghi ''

if ($loi.Count -gt 0) {
    Ghi 'Nhung viec chua tron:' 'Yellow'
    $loi | ForEach-Object { Ghi "  - $_" 'Yellow' }
    Ghi ''
    Ghi 'Hay chup man hinh nay gui lai cho tro ly de duoc huong dan tiep.' 'Yellow'
    Ghi 'Ghi chu ky thuat co trong: cau-hinh\GHI-CHU.md' 'Gray'
    Ghi ''
}

Ghi 'VIEC TIEP THEO:' 'Cyan'
Ghi ''
Ghi "  1. Bat may chu game: mo $GAME_ROOT"
Ghi '     bam CHUOT PHAI vao run.bat -> "Run as administrator"'
Ghi '     Doi cua so may chu hien len (khoang 30-60 giay).'
Ghi ''
Ghi '  2. Kiem tra may chu da nghe cong chua, mo PowerShell go:'
Ghi '       netstat -ano | findstr ":14445"'
Ghi '     Phai thay chu LISTENING.'
Ghi ''
Ghi '  3. Mo http://nrokame.com tren DIEN THOAI da TAT WIFI (dung 4G).'
Ghi '     Dung dien thoai vi may tinh VPS luon tu vao duoc chinh no,'
Ghi '     nen thu tren VPS khong chung minh duoc nguoi ngoai vao duoc.'
Ghi ''
Ghi '  4. Khi game da chay ngon lanh roi MOI chay 4-TU-DONG-BAT.bat'
Ghi '     de may chu tu len sau khi VPS khoi dong lai.'
Ghi ''
Ghi 'DUNG QUEN: tu bay gio, khi roi khoi Remote Desktop phai bam' 'Red'
Ghi 'dau X tren thanh xanh o dinh man hinh, TUYET DOI KHONG bam' 'Red'
Ghi '"Sign out" - bam Sign out la may chu game TAT, nguoi choi rot het.' 'Red'
Ghi ''
Read-Host '  Bam Enter de dong'
