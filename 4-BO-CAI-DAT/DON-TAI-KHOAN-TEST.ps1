# ===========================================================================
#  NANG MOT TAI KHOAN LEN QUAN TRI, ROI XOA 3 TAI KHOAN TEST
#
#  VI SAO PHAI LAM THEO THU TU NAY
#  ---------------------------------------------------------------------
#  Ba tai khoan test la  1 , 11 , 111 . Trong do HAI cai (1 va 111) dang
#  la quan tri vien. Neu xoa het truoc khi co tai khoan quan tri khac thi
#  ban mat sach duong vao trang quan tri: khong bat duoc bao tri, khong
#  khoa duoc nguoi pha game, khong spawn duoc boss.
#
#  Vi vay script nay lam DUNG THU TU: nang tai khoan that len truoc, kiem
#  tra chac chan da len roi, moi xoa ba cai test.
#
#  MOT COT DUY NHAT QUYET DINH QUYEN QUAN TRI
#  ---------------------------------------------------------------------
#  Da ra ca ma nguon: ca trang web (sv1\admin\giaodich.php) lan may chu
#  game (jdbc\daos\NDVSqlFetcher.java dong 93) deu doc cung mot cot
#  "is_admin". Dat is_admin = 1 la co quyen o CA HAI noi.
#  Cot "admin" nam canh do khong duoc dung o dau ca - dung nham la khong
#  co tac dung gi ma cung khong bao loi.
#
#  Tai khoan quan tri cung TU DONG BI LOAI khoi bang xep hang (cac cau
#  lenh TOP_SD, TOP_HP... deu co "WHERE account.is_admin = 0"), nen no
#  khong lam ban top.
#
#  DA THU THAT TREN BAN SAO CSDL TRUOC KHI GIAO (22/07/2026).
#  Chay lai nhieu lan khong sao.
# ===========================================================================

$ErrorActionPreference = 'Stop'
function Ghi($m, $mau = 'White') { Write-Host "  $m" -ForegroundColor $mau }

$MYSQL = 'C:\xampp\mysql\bin\mysql.exe'
$CSDL  = 'z'
$TEST  = @('1', '11', '111')

# Chay cau lenh SQL bang cach dua qua duong ong, KHONG dat tren dong lenh:
# moi thu viet tren dong lenh deu hien ra trong danh sach tien trinh cua
# Windows - ai dang nhap vao may cung doc duoc mat khau.
function ChaySQL($sql, $matKhau) {
    $doiSo = @('-u', 'root', '--default-character-set=utf8mb4')
    if ($matKhau) { $doiSo += "-p$matKhau" }
    $cu = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    try { $kq = $sql | & $MYSQL @doiSo 2>&1 } finally { $ErrorActionPreference = $cu }
    return $kq
}

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  DON TAI KHOAN TEST' -ForegroundColor Cyan
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host ''

if (-not (Test-Path $MYSQL)) {
    Ghi "DUNG LAI: khong thay $MYSQL" 'Red'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}

# --- Mat khau CSDL -------------------------------------------------------
Ghi 'Can mat khau co so du lieu (cai ban dat luc chay buoc 3).' 'Yellow'
$sec = Read-Host '  Go mat khau CSDL' -AsSecureString
$matKhau = [Runtime.InteropServices.Marshal]::PtrToStringAuto(
           [Runtime.InteropServices.Marshal]::SecureStringToBSTR($sec))

$thu = ChaySQL "USE $CSDL; SELECT 1;" $matKhau
if ($LASTEXITCODE -ne 0) {
    Ghi 'DUNG LAI: mat khau khong dung (hoac khong noi duoc CSDL).' 'Red'
    $thu | Select-Object -First 3 | ForEach-Object { Ghi "  $_" 'Red' }
    Read-Host "`n  Bam Enter de thoat"; exit 1
}
Ghi '[OK] Mat khau dung' 'Green'

# --- Cho xem tinh hinh hien tai ------------------------------------------
Write-Host ''
Ghi '--- Cac tai khoan dang co ---' 'Cyan'
$ds = ChaySQL @"
USE $CSDL;
SELECT a.id, a.username, a.is_admin,
       IF(a.password LIKE 'pbkdf2%','da ma hoa','CHU TRAN') AS mat_khau,
       (SELECT COUNT(*) FROM player p WHERE p.account_id = a.id) AS so_nhan_vat
FROM account a ORDER BY a.id;
"@ $matKhau
$ds | ForEach-Object { Ghi "  $_" }

# --- Hoi tai khoan nao se lam quan tri -----------------------------------
Write-Host ''
Ghi 'Go TEN tai khoan that ma ban vua dang ky tren web.' 'Yellow'
Ghi 'No se duoc nang len quan tri, roi ba tai khoan test se bi xoa.' 'Yellow'
$ten = (Read-Host '  Ten tai khoan quan tri moi').Trim()

# Kiem chu: chi cho chu, so va gach duoi. Ly do: ten nay duoc ghep thang
# vao cau lenh SQL. Ky tu la (dau nhay, dau xuyet) co the lam hong cau
# lenh hoac bi loi dung.
if ($ten -notmatch '^[A-Za-z0-9_]{3,20}$') {
    Ghi 'DUNG LAI: ten khong hop le (chi chu, so, gach duoi - 3 den 20 ky tu).' 'Red'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}
if ($TEST -contains $ten) {
    Ghi "DUNG LAI: '$ten' chinh la mot trong ba tai khoan test sap bi xoa." 'Red'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}

# Co that khong?
$co = ChaySQL "USE $CSDL; SELECT COUNT(*) FROM account WHERE username = '$ten';" $matKhau
$soCo = 0
if ($LASTEXITCODE -eq 0) { $soCo = [int](($co | Select-Object -Last 1).Trim()) }
if ($soCo -ne 1) {
    Ghi "DUNG LAI: khong tim thay tai khoan '$ten' trong CSDL." 'Red'
    Ghi 'Kiem tra lai chinh ta, hoac vao nrokame.com/sv1/dangky dang ky truoc.' 'Yellow'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}
Ghi "[OK] Tim thay tai khoan '$ten'" 'Green'

# --- Xac nhan ------------------------------------------------------------
Write-Host ''
Ghi 'Sap lam hai viec:' 'Yellow'
Ghi "  1. Nang '$ten' len quan tri (is_admin = 1)"
Ghi "  2. XOA HAN ba tai khoan test: $($TEST -join ', ') - ke ca nhan vat cua chung"
Write-Host ''
$tl = Read-Host '  Dong y? Go: DONG Y (chu HOA)'
# -ceq chu khong phai -eq: -eq cua PowerShell khong phan biet hoa thuong,
# go "dong y" thuong cung xoa. Bat go chu HOA de nguoi dung dung lai mot nhip.
if ($tl -cne 'DONG Y') {
    Ghi 'Da huy, khong thay doi gi.' 'Yellow'
    Read-Host "`n  Bam Enter de thoat"; exit 0
}

# --- Buoc 1: nang quyen --------------------------------------------------
Write-Host ''
Ghi '--- Buoc 1: nang quyen quan tri ---' 'Cyan'
ChaySQL "USE $CSDL; UPDATE account SET is_admin = 1 WHERE username = '$ten';" $matKhau | Out-Null
if ($LASTEXITCODE -ne 0) {
    Ghi 'DUNG LAI: khong nang quyen duoc. Chua xoa gi ca.' 'Red'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}

# KIEM CHUNG THAT truoc khi xoa - khong tin lenh chay xong la da co tac dung
$kt = ChaySQL "USE $CSDL; SELECT is_admin FROM account WHERE username = '$ten';" $matKhau
$giaTri = ($kt | Select-Object -Last 1).Trim()
if ($giaTri -ne '1') {
    Ghi "DUNG LAI: '$ten' van chua co quyen quan tri (is_admin = $giaTri)." 'Red'
    Ghi 'KHONG xoa tai khoan test - de ban con duong vao trang quan tri.' 'Yellow'
    Read-Host "`n  Bam Enter de thoat"; exit 1
}
Ghi "[OK] '$ten' da la quan tri vien" 'Green'

# --- Buoc 2: xoa tai khoan test ------------------------------------------
Write-Host ''
Ghi '--- Buoc 2: xoa ba tai khoan test ---' 'Cyan'
$dsTest = "'" + ($TEST -join "','") + "'"
$xoa = ChaySQL @"
USE $CSDL;
CREATE TEMPORARY TABLE tk_test AS SELECT id FROM account WHERE username IN ($dsTest);
CREATE TEMPORARY TABLE nv_test AS SELECT id FROM player WHERE account_id IN (SELECT id FROM tk_test);
DELETE FROM shop_ky_gui            WHERE player_id  IN (SELECT id FROM nv_test);
DELETE FROM site_comment_likes     WHERE account_id IN (SELECT id FROM tk_test)
                                      OR comment_id IN (SELECT id FROM site_comments
                                                        WHERE account_id IN (SELECT id FROM tk_test));
DELETE FROM site_comments          WHERE account_id IN (SELECT id FROM tk_test);
DELETE FROM history_items_diemdanh WHERE account_id IN (SELECT id FROM tk_test);
DELETE FROM phuc_loi_progress      WHERE account_id IN (SELECT id FROM tk_test);
DELETE FROM player                 WHERE account_id IN (SELECT id FROM tk_test);
DELETE FROM account                WHERE id IN (SELECT id FROM tk_test);
"@ $matKhau
if ($LASTEXITCODE -ne 0) {
    Ghi 'LOI khi xoa:' 'Red'
    $xoa | Select-Object -First 8 | ForEach-Object { Ghi "  $_" 'Red' }
    Read-Host "`n  Bam Enter de thoat"; exit 1
}
Ghi '[OK] Da xoa' 'Green'

# --- Kiem chung lai tu dau -----------------------------------------------
Write-Host ''
Ghi '--- Kiem chung ---' 'Cyan'
$sau = ChaySQL @"
USE $CSDL;
SELECT a.id, a.username, a.is_admin,
       (SELECT COUNT(*) FROM player p WHERE p.account_id = a.id) AS so_nhan_vat
FROM account a ORDER BY a.id;
"@ $matKhau
$sau | ForEach-Object { Ghi "  $_" }

Write-Host ''
$mo = ChaySQL @"
USE $CSDL;
SELECT 'nhan vat mo coi'   AS cho, COUNT(*) AS so FROM player p            LEFT JOIN account a ON a.id=p.account_id WHERE a.id IS NULL
UNION ALL SELECT 'diem danh mo coi', COUNT(*) FROM history_items_diemdanh h LEFT JOIN account a ON a.id=h.account_id WHERE a.id IS NULL
UNION ALL SELECT 'phuc loi mo coi',  COUNT(*) FROM phuc_loi_progress q      LEFT JOIN account a ON a.id=q.account_id WHERE a.id IS NULL;
"@ $matKhau
$mo | ForEach-Object { Ghi "  $_" }

Write-Host ''
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host '  XONG' -ForegroundColor Green
Write-Host '===========================================================' -ForegroundColor Cyan
Write-Host ''
Ghi "Tu gio dang nhap trang quan tri bang tai khoan: $ten" 'Yellow'
Ghi '  nrokame.com/sv1/admin/panel'
Write-Host ''
Read-Host '  Bam Enter de dong'
