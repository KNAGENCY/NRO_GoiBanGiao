// kh8gate_test.go — ⛔ CONG TUONG DUONG BYTE cua KH8 (task 0).
//
// CAU HOI DUY NHAT FILE NAY TRA LOI:
//
//	"Du lieu nap tu YAML co sinh ra DUNG TUNG BYTE cung mot goi tin nhu du lieu
//	 Go hom nay khong?"
//
// Chua tra loi duoc cau do thi KHONG duoc chinh mot con so nao trong kho noi
// dung — vi neu chinh truoc, khong ai con phan biet duoc "so moi" voi "bo nap
// doc sai". Do la ly do KH8 dat task nay lam CONG (thiet ke §7.1).
//
// HAI LOP KIEM, CO Y KHAC NHAU:
//
//  1. Go == YAML. Bat duoc "bo nap lam lech so".
//  2. Go == BAN VANG ghim tay va YAML == BAN VANG (kh8gate_vang_test.go). Bat
//     duoc truong hop HAI BEN CUNG SAI — dung diem mu oracle ma
//     mauvatpham_test.go:47 dang mac phai (no dung `want` tu chinh nguon du lieu
//     no phai kiem).
//
// Chi lop 1 thi mot loi trong CHINH bo dung goi se lot; chi lop 2 thi khong noi
// duoc gi ve bo nap. Phai co ca hai.
//
// LUU Y VE CACH DUNG: cac phep so duoi day goi CHINH cac ham dung goi that
// (goiCatalogSkillTu, goiCatalogItemOptionTu, goiCatalogItemMauTu, goiCuaHang),
// KHONG viet lai mot ban ghi byte thu hai trong test. Viet lai la chi chung minh
// ban sao dung, khong chung minh du lieu dung.
package game

import (
	"bytes"
	"encoding/hex"
	"testing"

	"nroserver/internal/noidung"
)

// napKhoYAML nap kho noi dung NHUNG (ban di cung file thuc thi). Fatal chu khong
// Skip: kho noi dung thieu la mot loi that, khong phai mot dieu kien moi truong.
func napKhoYAML(t *testing.T) *noidung.DuLieu {
	t.Helper()
	du, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("nap kho noi dung YAML that bai: %v", err)
	}
	return du
}

// soByte so hai than goi va bao cho biet lech o BYTE THU MAY — mot thong bao
// "khong bang nhau" tren 3483 byte thi vo dung khi phai di tim.
func soByte(t *testing.T, ten string, duoc, mongDoi []byte) {
	t.Helper()
	if bytes.Equal(duoc, mongDoi) {
		return
	}
	if len(duoc) != len(mongDoi) {
		t.Errorf("%s: dai %d byte, mong doi %d byte", ten, len(duoc), len(mongDoi))
	}
	n := len(duoc)
	if len(mongDoi) < n {
		n = len(mongDoi)
	}
	for i := 0; i < n; i++ {
		if duoc[i] != mongDoi[i] {
			tu, den := i-8, i+8
			if tu < 0 {
				tu = 0
			}
			if den > n {
				den = n
			}
			t.Fatalf("%s: lech tu byte %d — duoc 0x%02x, mong doi 0x%02x\n  duoc     ...%s...\n  mong doi ...%s...",
				ten, i, duoc[i], mongDoi[i],
				hex.EncodeToString(duoc[tu:den]), hex.EncodeToString(mongDoi[tu:den]))
		}
	}
	t.Fatalf("%s: mot ben la doan dau cua ben kia (dai %d vs %d)", ten, len(duoc), len(mongDoi))
}

func giaiVang(t *testing.T, ten, hexVang string) []byte {
	t.Helper()
	b, err := hex.DecodeString(hexVang)
	if err != nil {
		t.Fatalf("ban vang %s hong: %v", ten, err)
	}
	return b
}

// ---------------------------------------------------------------------------
// GATE 1 — goi -28/7 (chieu thuc): 5 lang x 4 chieu x 3 cap
// ---------------------------------------------------------------------------

// ⚠️ KH8 TASK 3 DA DOI CAU HOI CUA GATE NAY.
//
// Truoc task 3, gate chieu so BA ban: Go vs YAML vs ban vang. Task 3 XOA bang
// Go (chieuthuc.go khong con taoCapChieu/chieuTheoLang) — chi con MOT nguon.
// "Go vs YAML" bay gio la YAML so voi CHINH NO: mot phep so luon dung, tuc mot
// test khong bao gio do, tuc rac. Da bo.
//
// Con lai MOT phep so, va no van la phep so co gia tri nhat cua ca gate: goi
// -28/7 dung tu kho noi dung phai KHOP TUNG BYTE voi BAN VANG ghim tay
// (kh8gate_vang_test.go) — ban vang duoc chot tu du lieu Go hom truoc task 3.
// Doi mot con so trong chieu.yaml -> test nay do ngay.
//
// Kiem CA HAI duong vao du lieu, vi tu task 3 chung la hai duong khac nhau toi
// cung mot kho:
//
//	goiCatalogSkill()          — duong THAT cua may chu (qua anh chup kho nhung)
//	goiCatalogSkillTu(chieuTuYAML(...)) — duong doc thang file YAML
//
// Hai duong nay lech nhau = bo doi kieu (khochieu.go hoac noidungyaml.go) sai.
func TestKH8Gate_ChieuThuc_DungByteVoiBanVang(t *testing.T) {
	vang := giaiVang(t, "-28/7", vangCatalogSkill)

	tuAnhChup, err := goiCatalogSkill()
	if err != nil {
		t.Fatalf("dung goi -28/7 tu anh chup kho noi dung loi: %v", err)
	}
	soByte(t, "-28/7 (anh chup kho) vs BAN VANG", tuAnhChup.Data, vang)

	tuYAML, err := goiCatalogSkillTu(chieuTuYAML(napKhoYAML(t)))
	if err != nil {
		t.Fatalf("dung goi -28/7 tu du lieu YAML loi: %v", err)
	}
	if tuAnhChup.Command != tuYAML.Command {
		t.Fatalf("ma lenh khac nhau: anh chup %d, YAML %d", tuAnhChup.Command, tuYAML.Command)
	}
	soByte(t, "-28/7 (doc thang YAML) vs BAN VANG", tuYAML.Data, vang)
}

// ---------------------------------------------------------------------------
// GATE 2 — goi -28/8 subtype 0, 1 va 2 (mau option + mau vat pham)
// ---------------------------------------------------------------------------

// ⚠️ KH8 TASK 4 DA DOI CAU HOI CUA HAI GATE NAY — cung ly do task 3 doi gate
// chieu. Bang Go (danhSachMauOption / danhSachMauVatPham) DA BI XOA, chi con
// vat-pham.yaml. "Go vs YAML" bay gio la YAML so voi CHINH NO: mot phep so luon
// dung, tuc rac. Da bo.
//
// Con lai MOT phep so, va no van la phep so co gia tri nhat: goi -28/8 phai KHOP
// TUNG BYTE voi BAN VANG ghim tay (kh8gate_vang_test.go) — ban vang chot tu du
// lieu Go hom truoc task 4, KHONG doi mot byte nao o task nay.
//
// Kiem CA HAI duong vao du lieu, vi tu task 4 chung la hai duong khac nhau toi
// cung mot kho:
//
//	goiCatalogItemOption()/goiCatalogItemMau()  — duong THAT (qua anh chup nhung)
//	goiCatalogItem…Tu(…TuYAML(napKhoYAML))      — duong doc thang file YAML
//
// Hai duong lech nhau = bo doi kieu (mauvatpham.go hoac noidungyaml.go) sai.
func TestKH8Gate_MauOption_DungByteVoiBanVang(t *testing.T) {
	vang := giaiVang(t, "-28/8 subtype 0", vangCatalogItemOption)

	tuAnhChup, err := goiCatalogItemOption()
	if err != nil {
		t.Fatalf("dung goi -28/8 subtype 0 tu anh chup kho noi dung loi: %v", err)
	}
	soByte(t, "-28/8 subtype 0 (anh chup kho) vs BAN VANG", tuAnhChup.Data, vang)

	tuYAML, err := goiCatalogItemOptionTu(mauOptionTuYAML(napKhoYAML(t)))
	if err != nil {
		t.Fatalf("dung goi -28/8 subtype 0 tu du lieu YAML loi: %v", err)
	}
	if tuAnhChup.Command != tuYAML.Command {
		t.Fatalf("ma lenh khac nhau: anh chup %d, YAML %d", tuAnhChup.Command, tuYAML.Command)
	}
	soByte(t, "-28/8 subtype 0 (doc thang YAML) vs BAN VANG", tuYAML.Data, vang)
}

func TestKH8Gate_MauVatPham_DungByteVoiBanVang(t *testing.T) {
	vang := giaiVang(t, "-28/8 subtype 1", vangCatalogItemMau)

	tuAnhChup, err := goiCatalogItemMau()
	if err != nil {
		t.Fatalf("dung goi -28/8 subtype 1 tu anh chup kho noi dung loi: %v", err)
	}
	soByte(t, "-28/8 subtype 1 (anh chup kho) vs BAN VANG", tuAnhChup.Data, vang)

	tuYAML, err := goiCatalogItemMauTu(mauVatPhamTuYAML(napKhoYAML(t)))
	if err != nil {
		t.Fatalf("dung goi -28/8 subtype 1 tu du lieu YAML loi: %v", err)
	}
	if tuAnhChup.Command != tuYAML.Command {
		t.Fatalf("ma lenh khac nhau: anh chup %d, YAML %d", tuAnhChup.Command, tuYAML.Command)
	}
	soByte(t, "-28/8 subtype 1 (doc thang YAML) vs BAN VANG", tuYAML.Data, vang)
}

// TestKH8Gate_MauVatPhamSubtype2 ghim subtype 2. No khong mang du lieu noi dung
// (khoang rong [0,0)) nen KHONG co ban YAML de so — nhung no van phai dung y
// nguyen, vi no la cai duy nhat kich Service.clientOk() phia client: doi no la
// client ket man hinh "Xin cho" mai mai.
func TestKH8Gate_MauVatPhamSubtype2_KhongDoi(t *testing.T) {
	m := goiCatalogItem()
	soByte(t, "-28/8 subtype 2 vs BAN VANG", m.Data, giaiVang(t, "-28/8 subtype 2", vangCatalogItemSub2))
}

// ---------------------------------------------------------------------------
// GATE 3 — goi -44 (moi cua hang mau)
// ---------------------------------------------------------------------------

func TestKH8Gate_CuaHang_YAMLDungByteVoiGo(t *testing.T) {
	tuYAMLDS, err := cuaHangTuYAML(napKhoYAML(t))
	if err != nil {
		t.Fatalf("doi cua hang tu YAML loi: %v", err)
	}
	dsGo := DanhSachCuaHang()
	if len(tuYAMLDS) != len(dsGo) {
		t.Fatalf("YAML co %d cua hang, du lieu Go co %d", len(tuYAMLDS), len(dsGo))
	}
	if len(dsGo) == 0 {
		t.Fatal("khong co cua hang nao de so — cong nay se xanh oan")
	}

	for i := range dsGo {
		if tuYAMLDS[i].ID != dsGo[i].ID {
			t.Fatalf("cua hang thu %d: id YAML %d, id Go %d — id la khoa tra cuu, lech la ban nham bang gia",
				i, tuYAMLDS[i].ID, dsGo[i].ID)
		}
		tuGo, err := goiCuaHang(&dsGo[i])
		if err != nil {
			t.Fatalf("dung goi -44 tu du lieu Go (cua hang %d) loi: %v", dsGo[i].ID, err)
		}
		tuYAML, err := goiCuaHang(&tuYAMLDS[i])
		if err != nil {
			t.Fatalf("dung goi -44 tu du lieu YAML (cua hang %d) loi: %v", tuYAMLDS[i].ID, err)
		}
		soByte(t, "-44 Go vs YAML", tuYAML.Data, tuGo.Data)
	}

	// Ban vang ghim cua hang mau id 1 (Tap hoa Lang La).
	chGo := TimCuaHang(1)
	if chGo == nil {
		t.Fatal("khong con cua hang id 1 — ban vang khong con doi tuong de so")
	}
	mGo, err := goiCuaHang(chGo)
	if err != nil {
		t.Fatalf("dung goi -44 cua hang 1 loi: %v", err)
	}
	vang := giaiVang(t, "-44", vangCuaHangTapHoaLangLa)
	soByte(t, "-44 Go vs BAN VANG", mGo.Data, vang)

	var chYAML *CuaHang
	for i := range tuYAMLDS {
		if tuYAMLDS[i].ID == 1 {
			chYAML = &tuYAMLDS[i]
		}
	}
	if chYAML == nil {
		t.Fatal("YAML khong con cua hang id 1")
	}
	mYAML, err := goiCuaHang(chYAML)
	if err != nil {
		t.Fatalf("dung goi -44 cua hang 1 tu YAML loi: %v", err)
	}
	soByte(t, "-44 YAML vs BAN VANG", mYAML.Data, vang)
}

// ---------------------------------------------------------------------------
// Cong khong duoc xanh oan
// ---------------------------------------------------------------------------

// TestKH8Gate_KhoYAMLKhongRong khang dinh kho noi dung thuc su co du lieu.
// Neu mot file YAML rong thi MOI phep so o tren van co the xanh (goi rong ==
// goi rong) — day la kieu "test xanh vi khong kiem gi ca" ma du an da vap.
func TestKH8Gate_KhoYAMLKhongRong(t *testing.T) {
	du := napKhoYAML(t)
	soChieu := 0
	for lang := range du.ChieuTheoLang {
		if len(du.ChieuTheoLang[lang]) == 0 {
			t.Errorf("lang %d khong co chieu nao trong YAML", lang)
		}
		soChieu += len(du.ChieuTheoLang[lang])
	}
	if soChieu != 20 {
		t.Errorf("YAML co %d chieu, mong doi 20 (5 lang x 4 chieu — chu du an da CHOT giu 4 chieu/lang o KH8)", soChieu)
	}
	// KH8 task 4: bang Go da bi xoa, nen "YAML == Go" khong con y nghia (YAML so
	// voi chinh no). Ghim SO LUONG bang tay — mot con so o day thay doi la mot
	// viec CO CHU DINH phai ghi vao khoa-id.yaml cung luc.
	if len(du.MauOption) != 8 {
		t.Errorf("YAML co %d mau option, mong doi 8 (them mot mau = them mot dong khoa-id.yaml, va sua con so nay)", len(du.MauOption))
	}
	if len(du.MauVatPham) != 5 {
		t.Errorf("YAML co %d mau vat pham, mong doi 5 (them mot mau = them mot dong khoa-id.yaml, va sua con so nay)", len(du.MauVatPham))
	}
	if len(du.CuaHang) != len(DanhSachCuaHang()) {
		t.Errorf("YAML co %d cua hang, du lieu Go co %d", len(du.CuaHang), len(DanhSachCuaHang()))
	}
}
