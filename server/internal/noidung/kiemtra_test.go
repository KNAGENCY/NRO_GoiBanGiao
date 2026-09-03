// kiemtra_test.go — test cua BO KIEM DAY DU (KH8 task 1).
//
// Cung ky luat voi nap_test.go: lay kho noi dung THAT, va MOT DOT BIEN vao van
// ban that, doi bo nap TU CHOI. Khong dung fixture bia — mot fixture da don gian
// hoa co the lam test xanh ma khong kiem gi.
//
// Ban do dot bien -> test (theo bang §2.5 cua thiet ke):
//
//	M1  hoan vi hai mau vat pham        -> nap_test.go TestNap_HoanViHaiMauVatPhamBiTuChoi
//	                                       + TestKiem_DoiChoHaiMauRoiSuaIdChoLienTuc (day)
//	M2  doi id mot mau (2 -> 9)         -> nap_test.go TestNap_DoiIdMauVatPhamBiTuChoi
//	M3  xoa mau dang duoc CSDL tro toi  -> TestKiem_XoaMauCuoiDanhMuc,
//	                                       TestKiem_XoaMauGiuaDanhMuc (day)
//	M4  param option = 70000            -> TestKiem_ParamOptionVuotTran (day)
//	M5  param option = -1               -> TestKiem_ParamOptionAm (day)
//	M6  them mau option thu 256         -> TestKiem_MauOptionThu256 (day)
//	M9  cua hang tro toi mau khong co   -> nap_test.go TestNap_CuaHangTroToiMauKhongTonTai
//	M14 chieu khai loai la              -> nap_test.go TestNap_KieuChieuLaBiTuChoi
package noidung

import (
	"fmt"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// M3 — so khoa id (thu duy nhat bat duoc viec XOA / doi cho roi sua id)
// ---------------------------------------------------------------------------

// viTriMau tra ve chi so byte cua dong "  - id: N" TRONG KHOI mau_vat_pham.
// Bang mau option nam TRUOC trong cung file va dung y het thut le, nen tim
// thang bang strings.Index se tro nham vao khoi option — mot dot bien tro nham
// cho la mot test xanh vi ly do sai.
func viTriMau(t *testing.T, s string, id int) int {
	t.Helper()
	dau := strings.Index(s, "mau_vat_pham:")
	if dau < 0 {
		t.Fatal("khong tim thay khoi mau_vat_pham")
	}
	i := strings.Index(s[dau:], fmt.Sprintf("  - id: %d\n", id))
	if i < 0 {
		t.Fatalf("khong tim thay khoi mau id %d", id)
	}
	return dau + i
}

// TestKiem_XoaMauCuoiDanhMuc (M3): xoa mau id 4 — mau CUOI. Tap id con lai van
// la {0,1,2,3}: lien tuc, khong lo, khong trung. Moi kiem cua task 0 VAN XANH.
// Chi so khoa bat duoc.
func TestKiem_XoaMauCuoiDanhMuc(t *testing.T) {
	van := khoThat(t)
	van[TenFileVatPham] = van[TenFileVatPham][:viTriMau(t, van[TenFileVatPham], 4)]
	// Bo luon mon do khoi cua hang — nguoi xoa mot mau chac chan se lam viec nay,
	// va no lam kiem tham chieu cheo (K2) im tieng. Con lai DUY NHAT so khoa.
	iMon := strings.Index(van[TenFileCuaHang], "          - mau_id: 4")
	if iMon < 0 {
		t.Fatal("khong tim thay mon mau_id 4 trong cua hang")
	}
	van[TenFileCuaHang] = van[TenFileCuaHang][:iMon]
	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), TenFileKhoaId) {
		t.Errorf("loi bao %q — mong doi noi ro la lech voi %s", err, TenFileKhoaId)
	}
}

// TestKiem_XoaMauGiuaDanhMuc (M3): xoa mau id 2 nhung GIU nguyen so khoa. Ca hai
// bo kiem deu phai chan; chi can no do.
func TestKiem_XoaMauGiuaDanhMuc(t *testing.T) {
	van := khoThat(t)
	s := van[TenFileVatPham]
	i2, i3 := viTriMau(t, s, 2), viTriMau(t, s, 3)
	van[TenFileVatPham] = s[:i2] + s[i3:]
	mongDoiLoi(t, van)
}

// TestKiem_DoiChoHaiMauRoiSuaIdChoLienTuc (M1 bien the doc nhat): doi cho hai
// khoi VA sua luon `id:` cho lien tuc. Day la thao tac "don dep" tinh vi nhat —
// tap id van {0..4}, thu tu van tang dan, bo kiem "id = vi tri" cua task 0 XANH
// TRON VEN. No van la doi mon do cua moi nguoi choi.
func TestKiem_DoiChoHaiMauRoiSuaIdChoLienTuc(t *testing.T) {
	van := khoThat(t)
	s := van[TenFileVatPham]
	i2, i3, i4 := viTriMau(t, s, 2), viTriMau(t, s, 3), viTriMau(t, s, 4)
	khoi2, khoi3 := s[i2:i3], s[i3:i4]
	// Doi cho, roi danh so lai cho "dep": khoi cu la 3 thanh id 2, va nguoc lai.
	khoi3 = strings.Replace(khoi3, "  - id: 3\n", "  - id: 2\n", 1)
	khoi2 = strings.Replace(khoi2, "  - id: 2\n", "  - id: 3\n", 1)
	van[TenFileVatPham] = s[:i2] + khoi3 + khoi2 + s[i4:]

	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), TenFileKhoaId) {
		t.Errorf("loi bao %q — mong doi noi ro la lech voi %s (day la lo hong duy nhat "+
			"bo kiem task 0 khong thay)", err, TenFileKhoaId)
	}
}

// TestKiem_ThemMauMoiVaoCuoiKemDongKhoa: duong ĐUNG phai XANH. Neu test nay do
// thi bo kiem da bien danh muc thanh thu KHONG THE MO RONG — tuc la sai muc dich.
func TestKiem_ThemMauMoiVaoCuoiKemDongKhoa(t *testing.T) {
	van := khoThat(t)
	van[TenFileVatPham] += mauMoiId5
	van[TenFileKhoaId] += "  - id: 5\n    ma: \"shuriken-luyen-tap\"\n"

	du, err := napVan(van)
	if err != nil {
		t.Fatalf("them mot mau vao CUOI kem dong khoa ma van bi tu choi: %v", err)
	}
	if len(du.MauVatPham) != 6 {
		t.Fatalf("%d mau vat pham, mong doi 6", len(du.MauVatPham))
	}
	if du.MauVatPham[5].Ten != "Shuriken luyen tap" {
		t.Errorf("mau moi la %q", du.MauVatPham[5].Ten)
	}
}

// TestKiem_ThemMauMoiMaQuenDongKhoa: them mau ma quen ghi so -> DO. Neu bo qua,
// mau moi khong bao gio duoc doi chieu, va lan sau co nguoi xoa no thi khong ai
// bat duoc.
func TestKiem_ThemMauMoiMaQuenDongKhoa(t *testing.T) {
	van := khoThat(t)
	van[TenFileVatPham] += mauMoiId5
	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), "shuriken-luyen-tap") {
		t.Errorf("loi bao %q — mong doi noi ro ma nao con thieu dong khoa", err)
	}
}

// TestKiem_DoiTenMauDaKhoaBiTuChoi: doi ten hien thi cua mot mau da phat hanh la
// viec CO CHU DINH. Bo kiem do, va thong bao phai chi ra chinh xac phai sua gi.
func TestKiem_DoiTenMauDaKhoaBiTuChoi(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileVatPham, `ten: "Kunai luyện tập"`, `ten: "Kunai tập ném"`)
	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), "kunai-luyen-tap") {
		t.Errorf("loi bao %q — mong doi nhac ma da khoa", err)
	}
}

// TestKiem_XoaMotDongTrongSoKhoa: so khoa la file CHI-THEM. Xoa mot dong trong
// chinh no cung phai do — neu khong, ke muon xoa mot mau chi viec xoa ca hai cho.
func TestKiem_XoaMotDongTrongSoKhoa(t *testing.T) {
	van := khoThat(t)
	i := strings.Index(van[TenFileKhoaId], "  - id: 4\n    ma: \"kunai-luyen-tap\"\n")
	if i < 0 {
		t.Fatal("khong tim thay dong khoa cua mau id 4")
	}
	van[TenFileKhoaId] = van[TenFileKhoaId][:i]
	mongDoiLoi(t, van)
}

// TestKiem_SoKhoaThungLoBiTuChoi: mot so khoa hong thi moi ket luan rut ra tu no
// deu vo nghia — phai do TRUOC khi dem no ra do danh muc.
func TestKiem_SoKhoaThungLoBiTuChoi(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileKhoaId, "  - id: 3\n    ma: \"quan-ninja-lang-la\"\n",
		"  - id: 9\n    ma: \"quan-ninja-lang-la\"\n")
	mongDoiLoi(t, van)
}

// mauMoiId5 — mot mau hop le them vao CUOI danh muc.
const mauMoiId5 = `
  - id: 5
    loai: 2
    gioi_tinh: 3
    ten: "Shuriken luyen tap"
    mo_ta: "Shuriken cun dung de tap nem."
    cap: 1
    yeu_cau_suc: 0
    icon_id: 0
    phan: -1
    nang_cap_duoc: true
    tran_chong: 1
`

// ---------------------------------------------------------------------------
// M4 / M5 / M6 — mien gia tri (K3)
// ---------------------------------------------------------------------------

// themOptionVaoAo gan MOT option vao mon "Ao ninja Lang La" trong cua hang —
// mon dau tien co `la_moi: true`. Cua hang that hom nay khong ban mon nao kem
// option, nen dot bien M4/M5 phai tu tao cho dung.
func themOptionVaoAo(t *testing.T, van map[string]string, param string) {
	t.Helper()
	doiMotLan(t, van, TenFileCuaHang,
		"options: []\n            la_moi: true",
		"options:\n              - option_id: 0\n                param: "+param+"\n            la_moi: true")
}

// TestKiem_ParamOptionVuotTran (M4): param 70000 > 65535. Client doc param bang
// readUnsignedShort — 70000 tran ra thanh mot con so khac han.
func TestKiem_ParamOptionVuotTran(t *testing.T) {
	van := khoThat(t)
	themOptionVaoAo(t, van, "70000")
	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), "70000") {
		t.Errorf("loi bao %q — mong doi neu ra chinh gia tri sai", err)
	}
}

// TestKiem_ParamOptionAm (M5): param -1. Client doc KHONG DAU nen -1 hien thanh
// 65535 — mot con so khong lo, khong mot dong log.
func TestKiem_ParamOptionAm(t *testing.T) {
	van := khoThat(t)
	themOptionVaoAo(t, van, "-1")
	mongDoiLoi(t, van)
}

// TestKiem_OptionHopLeVanNapDuoc: neo cho hai test tren. Neu ban param HOP LE
// cung bi tu choi thi hai test kia xanh vi ly do sai.
func TestKiem_OptionHopLeVanNapDuoc(t *testing.T) {
	van := khoThat(t)
	themOptionVaoAo(t, van, "500")
	du, err := napVan(van)
	if err != nil {
		t.Fatalf("mot option hop le (param 500) bi tu choi: %v", err)
	}
	mon := du.CuaHang[0].Tabs[1].Mon[0]
	if len(mon.Options) != 1 || mon.Options[0].Param != 500 {
		t.Fatalf("option gan vao khong dung: %+v", mon.Options)
	}
}

// TestKiem_MonCoOptionPhaiDonChiec: gan option vao mon "Thuoc hoi mau nho"
// (tran_chong 99). Hai lo thuoc khac option chong chung mot o thi mot lo mat
// sach option luc luu CSDL.
func TestKiem_MonCoOptionPhaiDonChiec(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileCuaHang,
		"options: []\n            la_moi: false",
		"options:\n              - option_id: 0\n                param: 200\n            la_moi: false")
	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), "tran_chong") {
		t.Errorf("loi bao %q — mong doi noi ve tran_chong", err)
	}
}

// TestKiem_MauOptionThu256 (M6): client doc so mau option bang MOT BYTE KHONG
// DAU. Mau thu 256 lam client dung mot catalog RONG, im lang.
func TestKiem_MauOptionThu256(t *testing.T) {
	van := khoThat(t)
	var them strings.Builder
	for id := 8; id <= 255; id++ {
		fmt.Fprintf(&them, "  - id: %d\n    ten: \"Tang # thu %d\"\n    loai: 0\n", id, id)
	}
	doiMotLan(t, van, TenFileVatPham,
		"  - id: 7                      # ⚠️ đích của option 23 phía client (param ×1000)\n",
		them.String()+"  - id: 7                      # ⚠️ đích của option 23 phía client (param ×1000)\n")
	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), "256") {
		t.Errorf("loi bao %q — mong doi neu so muc that su co", err)
	}
}

// TestKiem_TenOptionThieuThangBiTuChoi: ten option khong co '#' thi nguoi choi
// thay mot dong chu tinh thay vi con so.
func TestKiem_TenOptionThieuThangBiTuChoi(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileVatPham, `ten: "Tăng # giáp"`, `ten: "Tăng giáp"`)
	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), "#") {
		t.Errorf("loi bao %q — mong doi noi ve dau '#'", err)
	}
}

// TestKiem_HoiChieuLechNhipActor: 1450ms khong phai boi so cua nhip actor 100ms
// -> chieu hoi som hay muon nua nhip tuy luc bam.
func TestKiem_HoiChieuLechNhipActor(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileChieu, "hoi_chieu_ms: 1000\n", "hoi_chieu_ms: 1450\n")
	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), "1450") {
		t.Errorf("loi bao %q — mong doi neu ra gia tri sai", err)
	}
}

// TestKiem_HeSoPhanTramBang0: mot cap chieu he_so_phan_tram = 0 -> chieu danh
// trung ma khong bao gio gay sat thuong (dot bien bat buoc #2, task 12).
func TestKiem_HeSoPhanTramBang0(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileChieu, "            he_so_phan_tram: 20\n", "            he_so_phan_tram: 0\n")
	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), "he_so_phan_tram") {
		t.Errorf("loi bao %q — mong doi noi ve he_so_phan_tram", err)
	}
}

// TestKiem_HeSoPhanTramVuotTran: go nham them mot chu so 0 (90 -> 900) phai bi
// tu choi thay vi lang le nap mot chieu manh gap 10 lan du dinh.
func TestKiem_HeSoPhanTramVuotTran(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileChieu, "            he_so_phan_tram: 90\n", "            he_so_phan_tram: 9000\n")
	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), "he_so_phan_tram") {
		t.Errorf("loi bao %q — mong doi noi ve he_so_phan_tram", err)
	}
}

// TestKiem_SkillIdTrungGiuaHaiLang: nap.go chi kiem trung TRONG mot lang. Client
// tra cuu chieu bang skillId khong kem theo lang, nen trung giua hai lang cung
// la mot chieu "an mat" chieu kia.
func TestKiem_SkillIdTrungGiuaHaiLang(t *testing.T) {
	van := khoThat(t)
	// skill_id 13 la cap dau cua lang 1; doi thanh 1 (da thuoc lang 0).
	doiMotLan(t, van, TenFileChieu, "          - skill_id: 13\n", "          - skill_id: 1\n")
	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), "skill_id") {
		t.Errorf("loi bao %q — mong doi noi ve skill_id", err)
	}
}

// TestKiem_TenTabBaDong: client tach ten tab bang '\n' va chi ve duoc 1-2 dong.
func TestKiem_TenTabBaDong(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileCuaHang, `- ten: "Thuoc"`, `- ten: "Thuoc\nhoi\nphuc"`)
	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), "dong") {
		t.Errorf("loi bao %q — mong doi noi ve so dong cua ten tab", err)
	}
}

// ---------------------------------------------------------------------------
// Sinh ma — kiem chinh cai thuoc do
// ---------------------------------------------------------------------------

// TestKhoaIdMaTu: neu ham sinh ma sai thi CA bo doi chieu khoa vo nghia — no la
// thuoc do, phai duoc do rieng.
func TestKhoaIdMaTu(t *testing.T) {
	ca := []struct{ vao, ra string }{
		{"Áo ninja Làng Lá", "ao-ninja-lang-la"},
		{"Thuốc hồi máu nhỏ", "thuoc-hoi-mau-nho"},
		{"Hồi phục # HP", "hoi-phuc-hp"},
		{"Tăng # né đòn", "tang-ne-don"},
		{"Kunai luyện tập", "kunai-luyen-tap"},
		{"ĐƯỜNG Ống 2", "duong-ong-2"},
		{"  nhiều   khoảng  trắng  ", "nhieu-khoang-trang"},
	}
	for _, c := range ca {
		if ra := khoaIdMaTu(c.vao); ra != c.ra {
			t.Errorf("khoaIdMaTu(%q) = %q, mong doi %q", c.vao, ra, c.ra)
		}
	}
}

// TestKhoaIdMaKhongTrungTrongKhoThat: hai mau trung ten thi trung ma, va so khoa
// mat kha nang phan biet chung. Phai bat duoc ngay tren ban that.
func TestKhoaIdMaKhongTrungTrongKhoThat(t *testing.T) {
	du, err := NapNhung()
	if err != nil {
		t.Fatalf("nap ban that loi: %v", err)
	}
	for _, bo := range []struct {
		ten string
		ma  []string
	}{
		{"mau option", maCuaOption(du.MauOption)},
		{"mau vat pham", maCuaVatPham(du.MauVatPham)},
	} {
		daThay := map[string]int{}
		for i, ma := range bo.ma {
			if ma == "" {
				t.Errorf("%s %d: ten sinh ra ma RONG", bo.ten, i)
			}
			if cu, co := daThay[ma]; co {
				t.Errorf("%s: id %d va id %d cung sinh ra ma %q", bo.ten, cu, i, ma)
			}
			daThay[ma] = i
		}
	}
}

// TestKiem_ThieuFileKhoaIdThiTuChoi: thieu so khoa KHONG duoc coi la "khong co gi
// de kiem" — do la duong im lang de vo hieu ca bo kiem.
func TestKiem_ThieuFileKhoaIdThiTuChoi(t *testing.T) {
	van := khoThat(t)
	delete(van, TenFileKhoaId)
	mongDoiLoi(t, van)
}
