// nap_test.go — kiem bo nap va bo kiem TOI THIEU cua task 0.
//
// Cach lam: lay kho noi dung THAT (ban nhung), do ra van ban, roi VA MOT DOT
// BIEN vao van ban do va doi bo nap TU CHOI. Dung van ban that (thay vi fixture
// bia) de test khong the xanh nho mot fixture da don gian hoa.
package noidung

import (
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"nroserver/kho-noi-dung"
)

// khoThat doc ba file YAML that ra dang van ban.
func khoThat(t *testing.T) map[string]string {
	t.Helper()
	m := map[string]string{}
	for _, ten := range TenCacFile() { // nguon duy nhat — xem chu thich TenCacFile (task 6: TenFileLang da them vao do)
		b, err := fs.ReadFile(khonoidung.FS, ten)
		if err != nil {
			t.Fatalf("khong doc duoc %s: %v", ten, err)
		}
		m[ten] = string(b)
	}
	return m
}

// napVan nap tu ba chuoi van ban.
func napVan(van map[string]string) (*DuLieu, error) {
	fsys := fstest.MapFS{}
	for ten, noiDung := range van {
		fsys[ten] = &fstest.MapFile{Data: []byte(noiDung)}
	}
	return Nap(fsys)
}

// doiMotLan doi lan xuat hien DAU TIEN cua cu thanh moi. Fatal neu khong tim
// thay — mot dot bien khong duoc va vao dau la mot dot bien VO NGHIA, va no se
// lam test "xanh" ma khong kiem gi ca.
func doiMotLan(t *testing.T, van map[string]string, ten, cu, moi string) {
	t.Helper()
	if !strings.Contains(van[ten], cu) {
		t.Fatalf("khong tim thay %q trong %s — dot bien khong va vao dau", cu, ten)
	}
	van[ten] = strings.Replace(van[ten], cu, moi, 1)
}

// TestNap_KhoThatNapDuoc: kho noi dung dang commit phai nap duoc va co du so
// muc. Day la neo cho moi test dot bien duoi day: neu ban goc da do thi "dot
// bien lam do" khong chung minh dieu gi.
func TestNap_KhoThatNapDuoc(t *testing.T) {
	du, err := NapNhung()
	if err != nil {
		t.Fatalf("kho noi dung dang commit khong nap duoc: %v", err)
	}
	soChieu := 0
	for i := range du.ChieuTheoLang {
		soChieu += len(du.ChieuTheoLang[i])
	}
	if soChieu != 20 {
		t.Errorf("%d chieu, mong doi 20", soChieu)
	}
	if len(du.MauOption) != 8 {
		t.Errorf("%d mau option, mong doi 8", len(du.MauOption))
	}
	if len(du.MauVatPham) != 5 {
		t.Errorf("%d mau vat pham, mong doi 5", len(du.MauVatPham))
	}
	if len(du.CuaHang) != 1 {
		t.Errorf("%d cua hang, mong doi 1", len(du.CuaHang))
	}
}

// TestNap_HoanViHaiMauVatPhamBiTuChoi (dot bien M1): doi cho hai mau trong file
// nhung GIU NGUYEN `id:` cua chung. Day la thao tac trong vo hai nhat trong ca
// KH8 — "sap xep lai cho gon" — va la thao tac doi mon do cua moi nguoi choi.
func TestNap_HoanViHaiMauVatPhamBiTuChoi(t *testing.T) {
	van := khoThat(t)
	i2 := strings.Index(van[TenFileVatPham], "  - id: 2\n")
	i3 := strings.Index(van[TenFileVatPham], "  - id: 3\n")
	i4 := strings.Index(van[TenFileVatPham], "  - id: 4\n")
	if i2 < 0 || i3 <= i2 || i4 <= i3 {
		t.Fatalf("khong dinh vi duoc ba khoi mau 2/3/4 (%d %d %d)", i2, i3, i4)
	}
	s := van[TenFileVatPham]
	van[TenFileVatPham] = s[:i2] + s[i3:i4] + s[i2:i3] + s[i4:]
	if _, err := napVan(van); err == nil {
		t.Fatal("hoan vi hai mau vat pham ma bo nap VAN CHAP NHAN — lo hong nguy hiem nhat cua KH8")
	}
}

// TestNap_DoiIdMauVatPhamBiTuChoi (dot bien M2): id 2 -> 9, tap id thung lo.
func TestNap_DoiIdMauVatPhamBiTuChoi(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileVatPham, "  - id: 2\n", "  - id: 9\n")
	if _, err := napVan(van); err == nil {
		t.Fatal("id mau vat pham thung lo ma bo nap van chap nhan")
	}
}

// TestNap_HoanViHaiMauOptionBiTuChoi: nhu tren nhung voi bang option — nguy
// hiem rieng vi client BAM CUNG y nghia cua id 6 va 7.
func TestNap_HoanViHaiMauOptionBiTuChoi(t *testing.T) {
	van := khoThat(t)
	s := van[TenFileVatPham]
	i6 := strings.Index(s, "  - id: 6")
	i7 := strings.Index(s, "  - id: 7")
	iHet := i7 + strings.Index(s[i7:], "# ------") // moc ket thuc khoi option
	if i6 < 0 || i7 <= i6 || iHet <= i7 {
		t.Fatalf("khong dinh vi duoc hai khoi option 6/7 (%d %d %d)", i6, i7, iHet)
	}
	van[TenFileVatPham] = s[:i6] + s[i7:iHet] + s[i6:i7] + s[iHet:]
	if _, err := napVan(van); err == nil {
		t.Fatal("hoan vi hai mau option ma bo nap van chap nhan")
	}
}

// TestNap_KieuChieuLaBiTuChoi (dot bien M14): loai: 9 — mot kieu chieu may chu
// CHUA CAI DAT. Phai do, khong duoc im lang bo qua.
func TestNap_KieuChieuLaBiTuChoi(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileChieu, "        loai: 3\n", "        loai: 9\n")
	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), "loai") {
		t.Errorf("loi bao %q, mong doi noi ve `loai`", err)
	}
}

// TestNap_TemplateIdKhongBiaKhongDuocLaKhong: template_id = 0 kich nhanh nhoi
// hai cap gia phia client.
func TestNap_TemplateIdKhongBiTuChoi(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileChieu, "      - template_id: 1\n", "      - template_id: 0\n")
	mongDoiLoi(t, van)
}

// TestNap_ThieuMotLangBiTuChoi (dot bien M15): bo lang 4 khoi file.
func TestNap_ThieuMotLangBiTuChoi(t *testing.T) {
	van := khoThat(t)
	i := strings.Index(van[TenFileChieu], "  - id_lang: 4\n")
	if i < 0 {
		t.Fatal("khong tim thay khoi lang 4")
	}
	van[TenFileChieu] = van[TenFileChieu][:i]
	mongDoiLoi(t, van)
}

// TestNap_ThieuLang3BiTuChoi (KH8 task 6): xoa khoi lang id 3 khoi lang.yaml.
// Client cu chi gui duoc lang 0/1/2 nhung may chu van phai giu du 5, thieu mot
// lang la vo hieu "id lien tuc tu 0" ma kiemLang dua vao.
func TestNap_ThieuLang3BiTuChoi(t *testing.T) {
	van := khoThat(t)
	i := strings.Index(van[TenFileLang], "  - id: 3\n")
	j := strings.Index(van[TenFileLang], "  - id: 4\n")
	if i < 0 || j <= i {
		t.Fatalf("khong dinh vi duoc khoi lang 3 (%d %d)", i, j)
	}
	van[TenFileLang] = van[TenFileLang][:i] + van[TenFileLang][j:]
	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), TenFileLang) {
		t.Errorf("loi bao %q, mong doi noi ro %s", err, TenFileLang)
	}
}

// TestNap_SatThuongKhoiDauBangKhongBiTuChoi (KH8 task 6): sat_thuong_khoi_dau
// = 0 lam MOI nhan vat cua lang do gay 0 sat thuong vinh vien (khuvuc.go
// TinhSatThuong nhan gia tri nay lam satThuongGoc).
func TestNap_SatThuongKhoiDauBangKhongBiTuChoi(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileLang, "sat_thuong_khoi_dau: 55\n", "sat_thuong_khoi_dau: 0\n")
	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), "sat_thuong_khoi_dau") {
		t.Errorf("loi bao %q, mong doi noi ve sat_thuong_khoi_dau", err)
	}
}

// TestNap_HpKhoiDauBangKhongBiTuChoi: hp_khoi_dau <= 0 la nhan vat chet ngay
// luc tao.
func TestNap_HpKhoiDauBangKhongBiTuChoi(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileLang, "hp_khoi_dau: 200\n", "hp_khoi_dau: 0\n")
	mongDoiLoi(t, van)
}

// TestNap_ViTriLangVuotBienBiTuChoi: toa do am hoac vuot TranToaDoBanDo se
// tran khi gan vao BanDo.SpawnX/SpawnY (int16, game/banDo.go).
func TestNap_ViTriLangVuotBienBiTuChoi(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileLang, "vi_tri_x_khoi_dau: 100\n", "vi_tri_x_khoi_dau: 99999\n")
	err := mongDoiLoi(t, van)
	if !strings.Contains(err.Error(), "vi tri khoi dau") {
		t.Errorf("loi bao %q, mong doi noi ve vi tri khoi dau", err)
	}
}

// TestNap_LangHoanViBiTuChoi: hoan vi hai khoi lang (giu id) — nhu dot bien M1
// cua mau vat pham, nhung cho bang lang.
func TestNap_LangHoanViBiTuChoi(t *testing.T) {
	van := khoThat(t)
	s := van[TenFileLang]
	i0 := strings.Index(s, "  - id: 0\n")
	i1 := strings.Index(s, "  - id: 1\n")
	i2 := strings.Index(s, "  - id: 2\n")
	if i0 < 0 || i1 <= i0 || i2 <= i1 {
		t.Fatalf("khong dinh vi duoc ba khoi lang 0/1/2 (%d %d %d)", i0, i1, i2)
	}
	van[TenFileLang] = s[:i0] + s[i1:i2] + s[i0:i1] + s[i2:]
	if _, err := napVan(van); err == nil {
		t.Fatal("hoan vi hai khoi lang ma bo nap van chap nhan")
	}
}

// TestNap_CuaHangTroToiMauKhongTonTai (dot bien M9): mau_id 9 khong co that.
// Client se nem NullReference vao mot catch RONG va nuot mat.
func TestNap_CuaHangTroToiMauKhongTonTai(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileCuaHang, "- mau_id: 0 ", "- mau_id: 9 ")
	mongDoiLoi(t, van)
}

// TestNap_LoaiTienLaBiTuChoi: "vang" go nham thanh "vang " hay "gold" phai do.
func TestNap_LoaiTienLaBiTuChoi(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileCuaHang, "tien: \"vang\"", "tien: \"gold\"")
	mongDoiLoi(t, van)
}

// TestNap_TruongGoNhamTenBiTuChoi: `he_so_phan_tram` -> `he_so_phan_tam`. Neu
// bo nap bo qua truong la thi con so bi mat va thay bang 0 mot cach IM LANG.
func TestNap_TruongGoNhamTenBiTuChoi(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileChieu, "            he_so_phan_tram: 20\n", "            he_so_phan_tam: 20\n")
	mongDoiLoi(t, van)
}

// TestNap_TenTruongCuSatThuongBiTuChoi: KH8 §5.1 doi ten truong YAML tu
// `sat_thuong` (nghe nhu so tuyet doi) sang `he_so_phan_tram` (noi ro la %).
// Mot dong con dung ten cu phai bi KnownFields(true) tu choi voi thong diep
// de hieu, khong duoc am tham doc thanh 0 (dot bien bat buoc #1, task 12).
func TestNap_TenTruongCuSatThuongBiTuChoi(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileChieu, "            he_so_phan_tram: 20\n", "            sat_thuong: 20\n")
	mongDoiLoi(t, van)
}

// TestNap_CuPhapHongThiTuChoiCaBan: fail-closed — YAML hong khong duoc nap mot
// phan.
func TestNap_CuPhapHongThiTuChoiCaBan(t *testing.T) {
	van := khoThat(t)
	van[TenFileVatPham] = van[TenFileVatPham] + "\n  nay la: [ mot dong hong\n"
	mongDoiLoi(t, van)
}

// TestNap_XoaChuThichKhongAnhHuong: chu thich la thu duy nhat trong file KHONG
// duoc anh huong ket qua nap. Neu no anh huong thi cai gia da tra de chon YAML
// thay JSON bien thanh mot rui ro.
func TestNap_XoaChuThichKhongAnhHuong(t *testing.T) {
	goc, err := NapNhung()
	if err != nil {
		t.Fatalf("nap ban that loi: %v", err)
	}
	van := khoThat(t)
	for ten, noiDung := range van {
		var giu []string
		for _, dong := range strings.Split(noiDung, "\n") {
			if strings.HasPrefix(strings.TrimSpace(dong), "#") {
				continue
			}
			giu = append(giu, dong)
		}
		van[ten] = strings.Join(giu, "\n")
	}
	khongChuThich, err := napVan(van)
	if err != nil {
		t.Fatalf("bo het chu thich thi khong nap duoc nua: %v", err)
	}
	if len(khongChuThich.MauVatPham) != len(goc.MauVatPham) ||
		len(khongChuThich.MauOption) != len(goc.MauOption) ||
		len(khongChuThich.CuaHang) != len(goc.CuaHang) {
		t.Fatal("bo chu thich lam doi so muc nap duoc")
	}
	for i := range goc.ChieuTheoLang {
		if len(khongChuThich.ChieuTheoLang[i]) != len(goc.ChieuTheoLang[i]) {
			t.Fatalf("lang %d: bo chu thich lam doi so chieu", i)
		}
		for j := range goc.ChieuTheoLang[i] {
			a, b := goc.ChieuTheoLang[i][j], khongChuThich.ChieuTheoLang[i][j]
			if a.Ten != b.Ten || len(a.CacCap) != len(b.CacCap) {
				t.Fatalf("lang %d chieu %d doi khi bo chu thich", i, j)
			}
			for k := range a.CacCap {
				if a.CacCap[k] != b.CacCap[k] {
					t.Fatalf("lang %d chieu %q cap %d doi khi bo chu thich", i, a.Ten, k+1)
				}
			}
		}
	}
}

// mongDoiLoi khang dinh ban da dot bien KHONG nap duoc, va tra loi de nguoi goi
// soi them noi dung thong bao.
func mongDoiLoi(t *testing.T, van map[string]string) error {
	t.Helper()
	du, err := napVan(van)
	if err == nil {
		t.Fatalf("ban da dot bien VAN NAP DUOC (%d mau vat pham) — bo kiem co lo hong", len(du.MauVatPham))
	}
	return err
}
