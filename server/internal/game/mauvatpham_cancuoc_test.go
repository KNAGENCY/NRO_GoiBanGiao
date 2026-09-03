// mauvatpham_cancuoc_test.go — CAN CUOC cua mau vat pham / mau option
// (KH8 task 4, thiet ke §2.1 va §2.6).
//
// CAU HOI DUY NHAT FILE NAY TRA LOI:
//
//	"Cai id ma may chu dung hom nay co con la cai id dang nam trong CSDL cua
//	 nguoi choi khong?"
//
// Ba su that phai trung khop, khong duoc ro ra:
//
//	id trong vat-pham.yaml  ==  CHI SO trong slice may chu dung
//	                        ==  dong tuong ung trong khoa-id.yaml
//
// Lech mot cai la "ao ninja cua moi nguoi bien thanh thuoc hoi mau" — khong mot
// dong log, khong mot thong bao loi. Do la ly do task 4 duoc danh dau la task
// rui ro can cuoc cao nhat cua KH8.
package game

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"nroserver/internal/noidung"
)

// thuMucKhoNoiDung la duong toi CAY NGUON kho-noi-dung tinh tu internal/game.
const thuMucKhoNoiDung = "../../kho-noi-dung"

// ---------------------------------------------------------------------------
// Do nghe dung chung: cam mot kho khac lam nguon mac dinh
// ---------------------------------------------------------------------------

// ganKhoVatPhamTam cam `du` lam nguon mau vat pham mac dinh cua goi `game` cho
// den het test, roi tra lai nguyen trang.
//
// KHONG dung mot bang Go: `du` phai den tu noidung.NapNhung (tuc van la
// vat-pham.yaml) roi moi sua — neu test tu bia mot bang thi no lai dung dung cai
// bay ma task 4 vua go bo.
func ganKhoVatPhamTam(t *testing.T, du *noidung.DuLieu) {
	t.Helper()
	kho := noidung.MoKho(du)
	goc := khoVatPhamMacDinh
	t.Cleanup(func() { khoVatPhamMacDinh = goc })
	khoVatPhamMacDinh = func() khoVatPhamAnhChup { return kho.AnhChup() }
}

// themMauVatPhamTam nap kho NHUNG, them `m` vao CUOI danh muc (dung luat "chi
// them vao cuoi"), cam lam nguon mac dinh va tra ve id cua mau moi.
func themMauVatPhamTam(t *testing.T, m noidung.MauVatPham) int {
	t.Helper()
	du, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("nap kho noi dung: %v", err)
	}
	m.Id = len(du.MauVatPham)
	du.MauVatPham = append(du.MauVatPham, m)
	ganKhoVatPhamTam(t, du)
	return m.Id
}

// ---------------------------------------------------------------------------
// 1. id == chi so, o CA HAI bang
// ---------------------------------------------------------------------------

// TestCanCuocIdBangChiSoTrongAnhChup: bo kiem cua goi noidung khang dinh
// "id == vi tri trong FILE". Test nay khang dinh nua con lai — ANH CHUP ma goi
// `game` doc ra cung giu dung thu tu do, tuc chi so slice = id.
//
// DOT BIEN NHAM TOI: mot ngay nao do co ai sap xep lai danh muc trong bo nap
// (vi du sort theo ten cho "de doc"). Bo kiem YAML van xanh (file khong doi),
// nhung chi so slice khong con la id nua.
func TestCanCuocIdBangChiSoTrongAnhChup(t *testing.T) {
	ac := khoVatPhamMacDinh()

	dsMau := ac.MauVatPham()
	if len(dsMau) == 0 {
		t.Fatal("anh chup khong co mau vat pham nao — test se xanh oan")
	}
	for i, m := range dsMau {
		if m.Id != i {
			t.Errorf("mau vat pham o chi so %d mang id %d — client dat id = chi so vong lap (Controller.cs:6498), lech la doi do cua nguoi choi", i, m.Id)
		}
	}

	dsOpt := ac.MauOption()
	if len(dsOpt) == 0 {
		t.Fatal("anh chup khong co mau option nao — test se xanh oan")
	}
	for i, o := range dsOpt {
		if o.Id != i {
			t.Errorf("mau option o chi so %d mang id %d — client dat id = chi so vong lap (Controller.cs:6479)", i, o.Id)
		}
	}
}

// TestCanCuocTimMauVatPhamTraDungMauCuaId: cua tra cuu that (TimMauVatPham, cai
// ma muaban/hanhtrang/roido dung) phai tra ve DUNG mau o chi so do.
//
// DOT BIEN NHAM TOI: doi anhChup.TimMauVatPham sang tra `a.duLieu.MauVatPham[id]`
// theo chi so mang — hom nay dung, nhung se sai ngay khi id thoi trung chi so.
func TestCanCuocTimMauVatPhamTraDungMauCuaId(t *testing.T) {
	ds := DanhSachMauVatPham()
	for i := range ds {
		m := TimMauVatPham(int16(i))
		if m == nil {
			t.Fatalf("TimMauVatPham(%d) tra nil — danh muc phai lien tuc tu 0", i)
		}
		if m.Ten != ds[i].Ten || m.Loai != ds[i].Loai || m.TranChong != ds[i].TranChong {
			t.Errorf("TimMauVatPham(%d) = %+v, danh muc o chi so %d = %+v", i, *m, i, ds[i])
		}
	}
}

// TestCanCuocTraVeBanSao: TimMauVatPham tra CON TRO. Truoc task 4 no tro thang
// vao bang Go toan cuc — mot cu ghi qua con tro do la doi mon do cua MOI nguoi
// choi. Bay gio no phai tro vao mot ban sao.
func TestCanCuocTraVeBanSao(t *testing.T) {
	m := TimMauVatPham(0)
	if m == nil {
		t.Fatal("khong tra duoc mau 0")
	}
	tenGoc := m.Ten
	m.Ten = "BI SUA"

	lai := TimMauVatPham(0)
	if lai == nil || lai.Ten != tenGoc {
		t.Fatalf("ghi vao ket qua TimMauVatPham lam ban danh muc: doc lai = %v, mong doi ten %q", lai, tenGoc)
	}
	if ds := DanhSachMauVatPham(); ds[0].Ten != tenGoc {
		t.Fatalf("DanhSachMauVatPham cung bi ban: %q", ds[0].Ten)
	}
}

// ---------------------------------------------------------------------------
// 2. Doi chieu voi SO KHOA (khoa-id.yaml) tu phia goi `game`
// ---------------------------------------------------------------------------

// docSoKhoa doc khoa-id.yaml tu CAY NGUON (khong phai ban nhung) va tra ve hai
// danh sach ma theo dung thu tu file.
//
// Doc thang file thay vi goi lai bo kiem cua goi noidung LA CO Y: neu test nay
// goi kiemDayDu thi no chi chung minh "bo kiem goi duoc bo kiem". O day ta muon
// mot phep doi chieu DOC LAP, tu phia nguoi tieu thu du lieu.
func docSoKhoa(t *testing.T) (option []string, vatPham []string) {
	t.Helper()
	duong := filepath.Join(thuMucKhoNoiDung, noidung.TenFileKhoaId)
	b, err := os.ReadFile(duong)
	if err != nil {
		t.Fatalf("doc %s: %v", duong, err)
	}
	return soKhoaTachMuc(t, string(b), "mau_option:"), soKhoaTachMuc(t, string(b), "mau_vat_pham:")
}

// soKhoaTachMuc boc cac gia tri `ma:` nam duoi mot dau muc, theo dung thu tu.
// Bo phan tich toi gian CO Y — keo ca thu vien YAML vao day chi de doc mot file
// hai truong la thua, va mot bo phan tich chuc dong thi doc duoc bang mat.
func soKhoaTachMuc(t *testing.T, noiDung, dauMuc string) []string {
	t.Helper()
	var ra []string
	trong := false
	for _, dong := range strings.Split(noiDung, "\n") {
		cat := strings.TrimSpace(strings.TrimRight(dong, "\r"))
		if cat == dauMuc {
			trong = true
			continue
		}
		// Mot dau muc khac (khong thut le, khong phai chu thich) ket thuc khoi.
		if trong && cat != "" && !strings.HasPrefix(dong, " ") &&
			!strings.HasPrefix(cat, "-") && !strings.HasPrefix(cat, "#") {
			trong = false
		}
		if !trong {
			continue
		}
		if i := strings.Index(cat, "ma:"); i >= 0 {
			ma := strings.TrimSpace(cat[i+len("ma:"):])
			ra = append(ra, strings.Trim(ma, `"`))
		}
	}
	if len(ra) == 0 {
		t.Fatalf("khong doc duoc muc nao duoi %q trong %s", dauMuc, noidung.TenFileKhoaId)
	}
	return ra
}

// TestCanCuocDanhMucKhopSoKhoa doi chieu danh muc dang chay voi so khoa: dung so
// luong, dung thu tu, dung ma (ma sinh tu ten hien thi).
//
// DOT BIEN NHAM TOI: doi `ten` cua mau 0 trong vat-pham.yaml -> ma sinh ra khac
// -> test nay do. Do la CO Y (thiet ke §2.6): doi ten mot mau da phat hanh la
// viec co chu dinh, phai sua dong khoa tuong ung trong cung mot commit.
func TestCanCuocDanhMucKhopSoKhoa(t *testing.T) {
	maOption, maVatPham := docSoKhoa(t)

	dsOpt := DanhSachMauOption()
	if len(dsOpt) != len(maOption) {
		t.Fatalf("danh muc co %d mau option, so khoa co %d dong — them mau ma quen ghi so (hay nguoc lai)",
			len(dsOpt), len(maOption))
	}
	for i, o := range dsOpt {
		if got := noidung.KhoaIdMaTu(o.Ten); got != maOption[i] {
			t.Errorf("mau option id %d ten %q -> ma %q, so khoa ghi %q", i, o.Ten, got, maOption[i])
		}
	}

	dsMau := DanhSachMauVatPham()
	if len(dsMau) != len(maVatPham) {
		t.Fatalf("danh muc co %d mau vat pham, so khoa co %d dong", len(dsMau), len(maVatPham))
	}
	for i, m := range dsMau {
		if got := noidung.KhoaIdMaTu(m.Ten); got != maVatPham[i] {
			t.Errorf("mau vat pham id %d ten %q -> ma %q, so khoa ghi %q", i, m.Ten, got, maVatPham[i])
		}
	}
}

// ---------------------------------------------------------------------------
// 3. Nap HONG thi may chu KHONG khoi dong — chu khong chay voi catalog RONG
// ---------------------------------------------------------------------------

// TestCanCuocDaoThuTuThiNapDO chung minh lai — tu phia goi `game` — rang cu dao
// hai mau roi sua id cho lien tuc (kieu "don dep cho gon") KHONG lot duoc:
// noidung.Nap tra loi, nen main.go dung lai. Khong co duong nao de may chu chay
// tiep voi mot danh muc da bi hoan vi.
//
// Dot bien nay la dot bien NGUY HIEM NHAT cua ca KH8: bo kiem "id == vi tri" mot
// minh KHONG bat duoc no (doi cho xong sua id lai thi id van lien tuc) — cai bat
// duoc no la SO KHOA (khoa-id.yaml), vi ma sinh tu ten se lech dong.
func TestCanCuocDaoThuTuThiNapDO(t *testing.T) {
	tam := chepKhoNoiDungRaThuMucTam(t)
	duong := filepath.Join(tam, "vat-pham.yaml")

	// Dao mau 0 va mau 1 ROI sua id cho lien tuc, la dieu mot nguoi "don dep" se
	// lam. Ve mat file, dao hai ban ghi roi danh so lai 0,1 CHINH LA doi cho hai
	// cai `ten` — moi truong con lai cua hai mau nay giong nhau tru mo ta, ma bo
	// kiem sinh tu `ten`. Nen day la dot bien that, khong phai ban nhai.
	b, err := os.ReadFile(duong)
	if err != nil {
		t.Fatalf("doc %s: %v", duong, err)
	}
	noiDung := string(b)
	const a, c, tam0 = "Thuốc hồi máu nhỏ", "Thuốc hồi chakra nhỏ", "\x00TAM\x00"
	if !strings.Contains(noiDung, a) || !strings.Contains(noiDung, c) {
		t.Fatalf("chuan bi dot bien sai: khong tim thay hai ten mau trong %s", duong)
	}
	noiDung = strings.ReplaceAll(noiDung, a, tam0)
	noiDung = strings.ReplaceAll(noiDung, c, a)
	noiDung = strings.ReplaceAll(noiDung, tam0, c)
	if err := os.WriteFile(duong, []byte(noiDung), 0o644); err != nil {
		t.Fatalf("ghi %s: %v", duong, err)
	}

	if _, err := noidung.NapTuThuMuc(tam); err == nil {
		t.Fatal("dao thu tu hai mau roi sua id cho lien tuc VAN NAP DUOC — so khoa khong con che duoc can cuoc, ao ninja se bien thanh thuoc hoi mau ma khong ai biet")
	}
}

// chepKhoNoiDungRaThuMucTam chep TOAN BO kho-noi-dung sang mot thu muc tam de
// dot bien tren ban chep, khong cham vao cay nguon.
func chepKhoNoiDungRaThuMucTam(t *testing.T) string {
	t.Helper()
	dich := t.TempDir()
	muc, err := os.ReadDir(thuMucKhoNoiDung)
	if err != nil {
		t.Fatalf("doc thu muc kho noi dung: %v", err)
	}
	soYAML := 0
	for _, m := range muc {
		if m.IsDir() || !strings.HasSuffix(m.Name(), ".yaml") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(thuMucKhoNoiDung, m.Name()))
		if err != nil {
			t.Fatalf("doc %s: %v", m.Name(), err)
		}
		if err := os.WriteFile(filepath.Join(dich, m.Name()), b, 0o644); err != nil {
			t.Fatalf("ghi %s: %v", m.Name(), err)
		}
		soYAML++
	}
	if soYAML == 0 {
		t.Fatal("khong chep duoc file YAML nao — test se xanh oan")
	}
	// Ban chep CHUA dot bien phai nap duoc, neu khong thi moi ket luan sau do
	// deu vo nghia (test do vi ly do sai).
	if _, err := noidung.NapTuThuMuc(dich); err != nil {
		t.Fatalf("ban chep sach cua kho noi dung khong nap duoc: %v", err)
	}
	return dich
}

// TestCanCuocKhoHongThiPanicChuKhongRong: neu ban NHUNG hong thi cua mac dinh
// PHAI no, khong duoc tra ve danh muc rong. Catalog rong = client vao game thay
// tui do trong tron, khong mot loi nao.
//
// O day ta khong lam hong duoc ban nhung (no nam trong file thuc thi) nen test
// khang dinh dieu con lai va con quan trong hon: cua mac dinh KHONG BAO GIO tra
// ve rong.
func TestCanCuocMacDinhKhongBaoGioRong(t *testing.T) {
	if n := len(DanhSachMauVatPham()); n == 0 {
		t.Fatal("DanhSachMauVatPham() rong — nguoi choi se thay tui do trong tron ma khong co loi nao")
	}
	if n := len(DanhSachMauOption()); n == 0 {
		t.Fatal("DanhSachMauOption() rong — moi option cua moi mon do bien mat im lang")
	}
}

// ---------------------------------------------------------------------------
// 4. Nguon TIEM VAO — cua `…Tu(ac)` doc dung anh chup duoc dua
// ---------------------------------------------------------------------------

// TestCanCuocDocTuAnhChupDuocDua: ba cua `…Tu(ac)` phai doc DUNG anh chup nguoi
// goi cam vao, khong len lut doc ban nhung. Day la cai bao dam mot actor dung
// k.anhChupNoiDung() se chay tron ven mot the he.
//
// DOT BIEN NHAM TOI: sua mauVatPhamTu/timMauVatPhamTu goi khoVatPhamMacDinh()
// thay vi dung `ac` — test nay do vi so lieu la ban da sua.
func TestCanCuocDocTuAnhChupDuocDua(t *testing.T) {
	du, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("nap kho noi dung: %v", err)
	}
	du.MauVatPham[0].Ten = "TEN CUA ANH CHUP TIEM VAO"
	du.MauOption[0].Ten = "OPTION CUA ANH CHUP TIEM VAO #"
	ac := noidung.MoKho(du).AnhChup()

	if ds := mauVatPhamTu(ac); len(ds) == 0 || ds[0].Ten != "TEN CUA ANH CHUP TIEM VAO" {
		t.Errorf("mauVatPhamTu doc nham nguon: %v", ds[0].Ten)
	}
	if m := timMauVatPhamTu(ac, 0); m == nil || m.Ten != "TEN CUA ANH CHUP TIEM VAO" {
		t.Errorf("timMauVatPhamTu doc nham nguon: %v", m)
	}
	if ds := mauOptionTu(ac); len(ds) == 0 || ds[0].Ten != "OPTION CUA ANH CHUP TIEM VAO #" {
		t.Errorf("mauOptionTu doc nham nguon: %v", ds[0].Ten)
	}

	// Ban MAC DINH van phai nguyen ven — anh chup tiem vao khong duoc ro ri.
	if TimMauVatPham(0).Ten == "TEN CUA ANH CHUP TIEM VAO" {
		t.Error("anh chup tiem vao lam ban ban mac dinh")
	}
}

// TestCanCuocKhuVucCoDuMauVatPham: anh chup ma ACTOR cam duoc (k.anhChupNoiDung)
// phai phuc vu duoc ca mau vat pham, khong chi chieu. Neu mai co mot cho tra cuu
// mau CHAY TRONG actor thi no da co duong dung san — khong phai them nguon thu
// hai (hai nguon = hai the he lech nhau).
func TestCanCuocKhuVucCoDuMauVatPham(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	ac := k.anhChupNoiDung()
	if n := len(mauVatPhamTu(ac)); n == 0 {
		t.Fatal("anh chup cua khu vuc khong co mau vat pham nao")
	}
	if m := timMauVatPhamTu(ac, 0); m == nil {
		t.Fatal("anh chup cua khu vuc khong tra duoc mau 0")
	}
}

// TestCanCuocKhuVucDocKhoDaGan khoa nua sau cua duong nap nong CHO MAU VAT PHAM:
// kho ma main.go cam vao (GanKhoNoiDung) phai la kho ma anh chup cua khu vuc doc
// ra, KHONG phai ban nhung.
//
// ⚠️ DAY LA HANG RAO cho bat ky cho tra cuu mau nao SAU NAY chay TRONG actor:
// no phai di qua `…Tu(k.anhChupNoiDung())`. Neu ai do dung cua khong tham so
// (DanhSachMauVatPham / TimMauVatPham) trong actor, kho da gan se KHONG toi
// duoc no — nap nong "thanh cong" nhung actor van chay so cu, im lang. Hom nay
// (task 4) chua co cho goi nao trong actor nen chua co test nao giet duoc dot
// bien do; test nay it nhat lam cho DUONG DUNG ton tai va co bang chung.
func TestCanCuocKhuVucDocKhoDaGan(t *testing.T) {
	du, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("nap kho noi dung: %v", err)
	}
	const tenMoi = "TEN CUA KHO DA GAN"
	du.MauVatPham[0].Ten = tenMoi

	s := MoSoKhuVuc()
	defer s.DungTatCa()
	s.GanKhoNoiDung(noidung.MoKho(du))

	k := s.Tra(1)
	if k == nil {
		t.Fatal("khong lay duoc khu vuc mapId 1")
	}
	m := timMauVatPhamTu(k.layAnhChupNhip(), 0)
	if m == nil || m.Ten != tenMoi {
		t.Fatalf("khu vuc khong doc kho da gan: %v (mong doi ten %q)", m, tenMoi)
	}
	// Cua khong tham so van doc ban NHUNG — day chinh la ly do khong duoc dung
	// no trong actor.
	if TimMauVatPham(0).Ten == tenMoi {
		t.Error("cua mac dinh bong nhien doc kho da gan — gia dinh cua test nay sai, doc lai chu thich")
	}
}
