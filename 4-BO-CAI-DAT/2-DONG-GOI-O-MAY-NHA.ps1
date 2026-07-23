# ===========================================================================
#  NGOC RONG KAME - BUOC 2: DONG GOI DU LIEU O MAY NHA
#  ---------------------------------------------------------------------
#  File nay chay TREN MAY HIEN TAI cua ban (may dang chay game), KHONG
#  phai tren VPS.
#
#  No lam 3 viec:
#    1. Xuat mot ban sao co so du lieu MOI NHAT (khong dung z.sql cu)
#    2. Dong goi thu muc game zz1 (da bo rac) thanh 1 file nen
#    3. Dong goi thu muc web htdocs (da bo rac) thanh 1 file nen
#
#  Ket qua: 3 file nam trong C:\Users\Admin\Music\bo-cai-vps\goi-len-vps\
#  Chep 3 file do len VPS la du.
#
#  CHAY LAI NHIEU LAN DUOC.
#
#  CACH CHAY: bam CHUOT PHAI vao file nay -> chon "Run with PowerShell".
#  Khong can quyen Administrator.
# ===========================================================================

$ErrorActionPreference = 'Stop'
$ProgressPreference    = 'SilentlyContinue'   # khong co dong nay, nen file cham gap 10 lan

# --- Cac duong dan o may nha (sua o day neu may ban khac) ----------------
$NGUON_GAME = 'C:\Users\Admin\Music\zz1'
$NGUON_WEB  = 'C:\xampp\htdocs'
$MYSQL_BIN  = 'C:\xampp\mysql\bin'
$DICH       = 'C:\Users\Admin\Music\bo-cai-vps\goi-len-vps'

function Ghi($m, $mau = 'White') { Write-Host "  $m" -ForegroundColor $mau }
function Tieu($m) {
    Write-Host ''
    Write-Host "  --- $m ---" -ForegroundColor Cyan
}

# CHU Y: trong PowerShell, ham PHAI duoc dinh nghia TRUOC cho goi no.
# Dat ham nay o cuoi file thi script se bao loi "khong nhan ra Nen-ThuMuc".
function Nen-ThuMuc {
    param($Nguon, $Dich)

    if (Test-Path $Dich) { Remove-Item $Dich -Force }

    # Uu tien 7-Zip neu may co: nhanh hon rat nhieu va khong vo o duong dan
    # dai. Compress-Archive tren PowerShell 5.1 nen 71.000 file rat cham
    # (20-40 phut) - van chay duoc, chi lau. Duong dan dai nhat trong data\
    # chi 66 ky tu nen khong dung gioi han 260 ky tu cua Windows.
    $bay = @(
        "$env:ProgramFiles\7-Zip\7z.exe",
        "${env:ProgramFiles(x86)}\7-Zip\7z.exe"
    ) | Where-Object { Test-Path $_ } | Select-Object -First 1

    if ($bay) {
        Ghi "Dung 7-Zip cho nhanh: $bay" 'Gray'
        & $bay a -tzip -mx1 "$Dich" "$Nguon\*" | Out-Null
    } else {
        Ghi 'Khong thay 7-Zip, dung cong cu nen san co cua Windows (cham hon).' 'Gray'
        Compress-Archive -Path "$Nguon\*" -DestinationPath $Dich -Force
    }

    if (Test-Path $Dich) {
        $mb = [math]::Round((Get-Item $Dich).Length / 1MB, 1)
        Ghi "[OK] Da nen xong: $Dich ($mb MB)" 'Green'
    } else {
        Ghi "LOI: khong tao duoc file nen $Dich" 'Red'
        throw "Nen that bai: $Dich"
    }
}

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  NGOC RONG KAME - BUOC 2: DONG GOI DU LIEU O MAY NHA' -ForegroundColor Cyan
Write-Host '===========================================================' -ForegroundColor Cyan
Ghi ''
Ghi 'Viec nay mat khoang 10-40 phut, phan lon la thoi gian nen'
Ghi 'thu muc anh (hon 71.000 file). Cu de may chay, dung tat.'
Ghi ''

# --- Kiem tra cho trong: can gap doi vi vua co ban chep tam vua co file nen
$oC = Get-CimInstance Win32_LogicalDisk -Filter "DeviceID='C:'"
$trong = [math]::Round($oC.FreeSpace / 1GB, 1)
Ghi "O C: con trong $trong GB."
if ($trong -lt 4) {
    Ghi 'DUNG LAI: can it nhat 4 GB trong o C:.' 'Red'
    Ghi 'Ly do: script vua tao ban chep tam (~860 MB) vua tao file nen (~700 MB).' 'Yellow'
    Read-Host "`n  Bam Enter de thoat"
    exit 1
}

New-Item -ItemType Directory -Path $DICH -Force | Out-Null


# ===========================================================================
#  PHAN 1 - XUAT CO SO DU LIEU
# ===========================================================================
Tieu 'PHAN 1/3: XUAT CO SO DU LIEU'

# VI SAO KHONG DUNG FILE zz1\z.sql CO SAN:
#   - z.sql co 47 bang, co so du lieu dang chay co 48 bang (thieu bang "news")
#   - cot mat khau trong z.sql VAN LA BAN CHUA MA HOA (khong co chuoi pbkdf2),
#     trong khi CSDL dang chay da chuyen sang dang pbkdf2$10000$...
#   - z.sql khong he chua chu "nrokame", nghia la cac link tai game trong
#     bang adminpanel van la ban truoc khi doi thuong hieu
#   Dung z.sql = quay nguoc toan bo cong suc doi thuong hieu va ma hoa
#   mat khau, MA KHONG CO THONG BAO LOI NAO.

if (-not (Test-Path "$MYSQL_BIN\mysqldump.exe")) {
    Ghi "DUNG LAI: khong thay $MYSQL_BIN\mysqldump.exe" 'Red'
    Ghi 'Ban co cai XAMPP o cho khac khong? Sua bien $MYSQL_BIN o dau file nay.' 'Yellow'
    Read-Host "`n  Bam Enter de thoat"
    exit 1
}

# Kiem tra MySQL co dang bat khong
$mysqlSong = $false
try {
    $thu = & "$MYSQL_BIN\mysql.exe" -u root -e "SELECT 1;" 2>&1
    if ($LASTEXITCODE -eq 0) { $mysqlSong = $true }
} catch { }

if (-not $mysqlSong) {
    Ghi 'DUNG LAI: khong noi duoc vao MySQL.' 'Red'
    Ghi ''
    Ghi 'Hay mo XAMPP Control Panel va bam Start o dong MySQL,' 'Yellow'
    Ghi 'doi den khi dong do chuyen sang mau xanh, roi chay lai file nay.' 'Yellow'
    Read-Host "`n  Bam Enter de thoat"
    exit 1
}
Ghi '[OK] MySQL dang chay.' 'Green'

$fileDump = "$DICH\z-moi-nhat.sql"

# CHU Y: dung --result-file, KHONG dung dau > cua PowerShell.
# Ly do: dau > cua PowerShell ghi file theo bang ma cua cua so lenh. Tren
# Windows Server 2019 bang ma mac dinh khong phai UTF-8 -> TOAN BO CHU
# TIENG VIET trong CSDL bi hong (ten vat pham, loi thoai NPC...) ma khong
# bao loi gi. --result-file la mysqldump tu ghi thang ra file, khong qua
# cua so lenh, nen chac chan dung.
Ghi 'Dang xuat co so du lieu "z"...'
& "$MYSQL_BIN\mysqldump.exe" -u root `
    --default-character-set=utf8mb4 `
    --routines --events --single-transaction `
    --result-file="$fileDump" z

if ($LASTEXITCODE -ne 0 -or -not (Test-Path $fileDump)) {
    Ghi 'DUNG LAI: xuat co so du lieu that bai.' 'Red'
    Read-Host "`n  Bam Enter de thoat"
    exit 1
}

$mbDump = [math]::Round((Get-Item $fileDump).Length / 1MB, 2)
Ghi "[OK] Da xuat: $fileDump ($mbDump MB)" 'Green'

# --- Tu kiem tra ban dump vua xuat co "that" khong -----------------------
# Ba dau hieu nay phan biet ban moi voi ban z.sql cu.
$noiDung = Get-Content $fileDump -Raw

$soBang = ([regex]::Matches($noiDung, '(?m)^CREATE TABLE ')).Count
$coPbkdf2 = $noiDung -match 'pbkdf2'
$coNrokame = $noiDung -match 'nrokame'
$coNews = $noiDung -match '(?m)^CREATE TABLE `news`'

Ghi ''
Ghi 'Kiem tra ban vua xuat:'
Ghi ("  So bang         : $soBang " + $(if ($soBang -ge 48) { '(dung)' } else { '(NGHI NGO - phai la 48)' })) `
    $(if ($soBang -ge 48) { 'Green' } else { 'Red' })
Ghi ("  Co bang news    : " + $(if ($coNews) { 'CO (dung)' } else { 'KHONG (NGHI NGO)' })) `
    $(if ($coNews) { 'Green' } else { 'Red' })
Ghi ("  Mat khau ma hoa : " + $(if ($coPbkdf2) { 'CO pbkdf2 (dung)' } else { 'KHONG (NGHI NGO)' })) `
    $(if ($coPbkdf2) { 'Green' } else { 'Red' })
Ghi ("  Thuong hieu moi : " + $(if ($coNrokame) { 'CO nrokame (dung)' } else { 'KHONG (NGHI NGO)' })) `
    $(if ($coNrokame) { 'Green' } else { 'Red' })

if (-not ($soBang -ge 48 -and $coPbkdf2 -and $coNrokame -and $coNews)) {
    Ghi ''
    Ghi 'CANH BAO: ban vua xuat khong giong ban dang chay.' 'Red'
    Ghi 'Co the ban da lo nhap de file z.sql cu len CSDL dang chay.' 'Yellow'
    Ghi 'DUNG chep len VPS. Hay bao lai truoc khi lam tiep.' 'Yellow'
    $tiep = Read-Host "`n  Van muon lam tiep? Go: co"
    if ($tiep -ne 'co') { exit 1 }
}
$noiDung = $null   # tra lai bo nho, file nay ~1 MB nhung khong can giu


# ===========================================================================
#  PHAN 2 - DONG GOI THU MUC GAME
# ===========================================================================
Tieu 'PHAN 2/3: DONG GOI THU MUC GAME (zz1)'

$tamGame = "$DICH\tam-zz1"
if (Test-Path $tamGame) { Remove-Item $tamGame -Recurse -Force }
New-Item -ItemType Directory -Path $tamGame -Force | Out-Null

Ghi 'Dang chep thu muc data (822 MB, hon 71.000 file) - lau nhat o day...'
Copy-Item "$NGUON_GAME\data" $tamGame -Recurse
Ghi 'Dang chep src, lib, zz1.jar, run.bat...'
Copy-Item "$NGUON_GAME\src"  $tamGame -Recurse
Copy-Item "$NGUON_GAME\lib"  $tamGame -Recurse
Copy-Item "$NGUON_GAME\zz1.jar", "$NGUON_GAME\run.bat" $tamGame

# ***************************************************************************
#  CHU Y CHO DAU DAY - DAU * SAU DUONG DAN LA BAT BUOC
#  ------------------------------------------------------------------------
#  Viet "Get-ChildItem $NGUON_GAME -Include *.txt" (khong co \*) thi
#  PowerShell tra ve 0 FILE va KHONG BAO LOI GI. Da kiem chung:
#      Get-ChildItem "...\zz1"   -Include *.txt,... ->  0 file
#      Get-ChildItem "...\zz1\*" -Include *.txt,... -> 13 file
#  Mat 12 file nay = mat TOAN BO thiet lap ban da chinh trong trang quan
#  tri (ti le EXP, drop, boss, phuc loi, nang sao). Game van len binh
#  thuong, chi lang le quay ve so mac dinh trong code - rat kho phat hien.
# ***************************************************************************
$fileLe = Get-ChildItem "$NGUON_GAME\*" -File -Include *.txt, *.properties, *.json, *.xml, *.mf |
          Where-Object { $_.Name -ne 'sources.txt' }

if ($fileLe.Count -lt 5) {
    Ghi "DUNG LAI: chi tim thay $($fileLe.Count) file cau hinh le, dang le phai co khoang 12." 'Red'
    Ghi 'Dung chep tiep - se thieu thiet lap trang quan tri.' 'Yellow'
    Read-Host "`n  Bam Enter de thoat"
    exit 1
}
$fileLe | Copy-Item -Destination $tamGame
Ghi "[OK] Da chep $($fileLe.Count) file cau hinh le o goc zz1." 'Green'

# Bo rac trong src: file .bak nay khong duoc bien dich, nhung chua ma cu
# gay loi de cua so chong chat, va con ten mien chu cu. Trang quan tri lai
# quet ca thu muc src bang regex nen de rac trong do la rui ro thua.
Remove-Item "$tamGame\src\server\ServerManager.java.bak-update-shop" -Force -EA SilentlyContinue

# Kiem tra icon.png con nguyen - thieu no la server vang ngay khi khoi dong
if (-not (Test-Path "$tamGame\src\server\icon.png")) {
    Ghi 'CANH BAO: thieu src\server\icon.png. Server se vang loi khi khoi dong.' 'Red'
}

$soFileData = (Get-ChildItem "$tamGame\data" -Recurse -File).Count
Ghi "So file trong data: $soFileData (o may nha la 71.296)"

Ghi 'Dang nen thu muc game - buoc nay lau, co the toi 30 phut...'
$zipGame = "$DICH\zz1-len-vps.zip"
Nen-ThuMuc -Nguon $tamGame -Dich $zipGame


# ===========================================================================
#  PHAN 3 - DONG GOI THU MUC WEB
# ===========================================================================
Tieu 'PHAN 3/3: DONG GOI THU MUC WEB (htdocs)'

$tamWeb = "$DICH\tam-web"
if (Test-Path $tamWeb) { Remove-Item $tamWeb -Recurse -Force }
New-Item -ItemType Directory -Path $tamWeb -Force | Out-Null

Copy-Item "$NGUON_WEB\sv1", "$NGUON_WEB\img" $tamWeb -Recurse
Copy-Item "$NGUON_WEB\index.php", "$NGUON_WEB\privacy-policy.php",
          "$NGUON_WEB\ddos_guard.php", "$NGUON_WEB\icon_asset.php" $tamWeb

# Don nhat ky chay thu o nha (co IP va ID admin)
Remove-Item "$tamWeb\sv1\admin\reset_log.txt" -Force -EA SilentlyContinue

# Kiem tra 2 file .htaccess con o trong sv1 - DAY LA PHAN BAO MAT.
#   sv1\down\.htaccess  : chi cho tai dung 3 ten file
#   sv1\admin\.htaccess : chan tai file .txt
# (Copy-Item CHEP DUOC file bat dau bang dau cham - da kiem chung. Kiem tra
#  o day chi de chac chan, khong phai vi Copy-Item hay bo sot.)
foreach ($h in @('down\.htaccess', 'admin\.htaccess')) {
    if (Test-Path "$tamWeb\sv1\$h") {
        Ghi "[OK] Co sv1\$h" 'Green'
    } else {
        Ghi "CANH BAO: THIEU sv1\$h - lop bao ve bi mat!" 'Red'
    }
}

$zipWeb = "$DICH\web-len-vps.zip"
Ghi 'Dang nen thu muc web...'
Nen-ThuMuc -Nguon $tamWeb -Dich $zipWeb


# ===========================================================================
#  DON DEP VA TONG KET
# ===========================================================================
Ghi ''
Ghi 'Dang xoa hai thu muc chep tam (khong con can nua)...'
Remove-Item $tamGame -Recurse -Force -EA SilentlyContinue
Remove-Item $tamWeb  -Recurse -Force -EA SilentlyContinue

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  XONG BUOC 2' -ForegroundColor Cyan
Write-Host '===========================================================' -ForegroundColor Cyan
Ghi ''
Ghi "Ba file can chep len VPS nam trong: $DICH" 'Yellow'
foreach ($f in @($fileDump, $zipGame, $zipWeb)) {
    if (Test-Path $f) {
        $mb = [math]::Round((Get-Item $f).Length / 1MB, 1)
        Ghi ("  " + (Split-Path $f -Leaf).PadRight(24) + "$mb MB") 'Green'
    } else {
        Ghi ("  " + (Split-Path $f -Leaf).PadRight(24) + 'THIEU!') 'Red'
    }
}
Ghi ''
Ghi 'Chep len VPS bang cach nao:' 'Yellow'
Ghi '  Cach 1 (de nhat): keo tha thang vao cua so Remote Desktop.'
Ghi '  Cach 2 (on dinh hon voi file lon): tai 3 file len Google Drive,'
Ghi '          roi mo trinh duyet TREN VPS de tai xuong.'
Ghi ''
Ghi 'DUNG XOA cac file nay o may nha sau khi chep xong.' 'Yellow'
Ghi 'May nha la ban goc du phong - giu it nhat 1 tuan sau khi VPS chay on.' 'Yellow'
Ghi ''
Read-Host '  Bam Enter de dong'
