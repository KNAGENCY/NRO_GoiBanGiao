package noidung

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// khoDuLieuThuNghiem dung mot anh chup gia trong do MOI phan du lieu deu mang
// dau the he `g`. Nho vay mot nguoi doc co the tu kiem: neu doc thay chieu the
// he 5 nhung vat pham the he 4 thi anh chup da bi tron nua cu nua moi.
func khoDuLieuThuNghiem(g int) *DuLieu {
	dl := &DuLieu{}
	for lang := 0; lang < SoLangCoDinh; lang++ {
		chieu := make([]Chieu, 0, 3)
		for t := 0; t < 3; t++ {
			caps := make([]CapChieu, 0, 4)
			for c := 0; c < 4; c++ {
				caps = append(caps, CapChieu{
					SkillId:      int16(g),
					Diem:         byte(c + 1),
					HeSoPhanTram: int16(g),
					Gia:          int16(g),
				})
			}
			chieu = append(chieu, Chieu{
				TemplateId: byte(t),
				Ten:        fmt.Sprintf("chieu-%d-%d-g%d", lang, t, g),
				CapToiDa:   4,
				CacCap:     caps,
			})
		}
		dl.ChieuTheoLang[lang] = chieu
	}
	for i := 0; i < 20; i++ {
		dl.MauOption = append(dl.MauOption, MauOption{
			Id: i, Ten: fmt.Sprintf("opt-%d-g%d", i, g), Loai: 1,
		})
		dl.MauVatPham = append(dl.MauVatPham, MauVatPham{
			Id: i, Ten: fmt.Sprintf("vp-%d", i), Cap: 1, IconId: int16(g),
		})
	}
	dl.CuaHang = []CuaHang{{
		Id:  1,
		Ten: "cua hang thu nghiem",
		Tabs: []TabCuaHang{{
			Ten: "tab",
			Mon: []MonHang{{
				MauId:   1,
				Gia:     int32(g),
				Tien:    TienVang,
				Options: []OptionMon{{OptionId: 1, Param: uint16(g)}},
			}},
		}},
	}}
	return dl
}

// khoDuoiG lay so the he tu duoi mot ten dang "...-g12". Tra -1 neu khong co.
// Dung chuoi chu khong dung truong so hep (byte/int8) vi test chay toi the he
// 1001 — nhet vao mot byte la tran, va cai "lech" do la loi cua TEST chu khong
// phai cua kho.
func khoDuoiG(ten string) int {
	i := strings.LastIndex(ten, "-g")
	if i < 0 {
		return -1
	}
	n, err := strconv.Atoi(ten[i+2:])
	if err != nil {
		return -1
	}
	return n
}

// theHeTrongDuLieu doc dau the he ra tu MOI goc cua anh chup. Tra ve tap cac gia
// tri tim thay — dung mot phan tu nghia la anh chup NHAT QUAN.
func theHeTrongDuLieu(a *anhChup) map[int]bool {
	thay := map[int]bool{}
	for lang := byte(0); lang < SoLangCoDinh; lang++ {
		for _, c := range a.ChieuCuaLang(lang) {
			thay[khoDuoiG(c.Ten)] = true
			for _, cc := range c.CacCap {
				thay[int(cc.HeSoPhanTram)] = true
			}
		}
	}
	for _, m := range a.MauVatPham() {
		thay[int(m.IconId)] = true
	}
	for _, o := range a.MauOption() {
		thay[khoDuoiG(o.Ten)] = true
	}
	for _, ch := range a.DanhSachCuaHang() {
		for _, tab := range ch.Tabs {
			for _, mon := range tab.Mon {
				thay[int(mon.Gia)] = true
				for _, op := range mon.Options {
					thay[int(op.Param)] = true
				}
			}
		}
	}
	return thay
}

func TestKhoMoKhoTheHeMotVaThayTangDonDieu(t *testing.T) {
	k := MoKho(khoDuLieuThuNghiem(1))
	if got := k.TheHe(); got != 1 {
		t.Fatalf("the he dau tien = %d, muon 1", got)
	}
	// Dot bien "Thay khong tang the he" chet o day.
	for i := 2; i <= 10; i++ {
		moi := k.Thay(khoDuLieuThuNghiem(i))
		if moi != uint64(i) {
			t.Fatalf("Thay lan %d tra the he %d, muon %d", i, moi, i)
		}
		if got := k.TheHe(); got != uint64(i) {
			t.Fatalf("sau Thay lan %d, TheHe() = %d, muon %d", i, got, i)
		}
	}
}

// TestKhoAnhChupCuVanDocDuocSauKhiThay khoa tinh chat K4: goroutine dang cam anh
// chup cu KHONG thay du lieu doi duoi tay.
func TestKhoAnhChupCuVanDocDuocSauKhiThay(t *testing.T) {
	k := MoKho(khoDuLieuThuNghiem(1))
	cu := k.AnhChup()

	for i := 2; i <= 5; i++ {
		k.Thay(khoDuLieuThuNghiem(i))
	}

	if cu.TheHe() != 1 {
		t.Fatalf("anh chup cu doi the he thanh %d", cu.TheHe())
	}
	thay := theHeTrongDuLieu(cu)
	if len(thay) != 1 || !thay[1] {
		t.Fatalf("anh chup cu khong con nhat quan: thay cac the he %v, muon chi {1}", thay)
	}
	if k.AnhChup().TheHe() != 5 {
		t.Fatalf("anh chup moi phai la the he 5")
	}
}

// TestKhoDuaThayVaDoc — cot loi cua task. 1 goroutine Thay 1000 lan trong khi 8
// goroutine doc. Duoi -race: 0 canh bao. Moi lan doc: anh chup NHAT QUAN.
//
// DOT BIEN M12 (nap nong sua tai cho *DuLieu thay vi Thay ca anh chup) lam test
// nay do duoi -race.
func TestKhoDuaThayVaDoc(t *testing.T) {
	const soVongThay = 1000
	const soNguoiDoc = 8

	k := MoKho(khoDuLieuThuNghiem(1))
	xong := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < soNguoiDoc; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			truoc := uint64(0)
			for {
				select {
				case <-xong:
					return
				default:
				}
				// HOP DONG: lay con tro MOT LAN, dung suot "nhip".
				a := k.AnhChup()
				if a.TheHe() < truoc {
					t.Errorf("the he lui: %d sau %d", a.TheHe(), truoc)
					return
				}
				truoc = a.TheHe()

				thay := theHeTrongDuLieu(a)
				if len(thay) != 1 {
					t.Errorf("anh chup the he %d lan lon: thay %v", a.TheHe(), thay)
					return
				}
				if !thay[int(a.TheHe())] {
					t.Errorf("anh chup the he %d chua du lieu the he %v", a.TheHe(), thay)
					return
				}
				if _, co := a.TraCapChieu(0, 0, 1); !co {
					t.Errorf("TraCapChieu that bai o the he %d", a.TheHe())
					return
				}
			}
		}()
	}

	for g := 2; g <= soVongThay+1; g++ {
		if moi := k.Thay(khoDuLieuThuNghiem(g)); moi != uint64(g) {
			t.Fatalf("Thay tra the he %d, muon %d", moi, g)
		}
	}
	close(xong)
	wg.Wait()

	if got := k.TheHe(); got != uint64(soVongThay+1) {
		t.Fatalf("the he cuoi = %d, muon %d", got, soVongThay+1)
	}
}

// ---------------------------------------------------------------------------
// M13 — cua tra cuu tra ban sao sau; nguoi goi ghi vao ket qua khong lam ban kho
// ---------------------------------------------------------------------------

func TestKhoNguoiGoiSuaKetQuaKhongLamBanAnhChup(t *testing.T) {
	k := MoKho(khoDuLieuThuNghiem(7))
	a := k.AnhChup()

	// chieu: sua ca phan tu ngoai lan cap ben trong
	ds := a.ChieuCuaLang(0)
	if len(ds) == 0 || len(ds[0].CacCap) == 0 {
		t.Fatal("du lieu thu nghiem rong")
	}
	ds[0].Ten = "BI SUA"
	ds[0].CapToiDa = 99
	ds[0].CacCap[0].HeSoPhanTram = 12345
	ds = append(ds, Chieu{TemplateId: 200})

	lai := a.ChieuCuaLang(0)
	if lai[0].Ten == "BI SUA" || lai[0].CapToiDa == 99 {
		t.Error("ChieuCuaLang: sua lop ngoai lam ban anh chup")
	}
	if lai[0].CacCap[0].HeSoPhanTram == 12345 {
		t.Error("ChieuCuaLang: sua CacCap lam ban anh chup (sao chep NONG)")
	}
	if len(lai) != len(ds)-1 {
		t.Error("ChieuCuaLang: append cua nguoi goi lam doi anh chup")
	}

	// vat pham
	vp := a.MauVatPham()
	vp[0].Ten = "BI SUA"
	if a.MauVatPham()[0].Ten == "BI SUA" {
		t.Error("MauVatPham: sua ket qua lam ban anh chup")
	}
	mot, co := a.TimMauVatPham(0)
	if !co {
		t.Fatal("TimMauVatPham(0) khong thay")
	}
	mot.Ten = "BI SUA"
	if lai, _ := a.TimMauVatPham(0); lai.Ten == "BI SUA" {
		t.Error("TimMauVatPham: sua ket qua lam ban anh chup")
	}

	// mau option
	op := a.MauOption()
	op[0].Ten = "BI SUA"
	if a.MauOption()[0].Ten == "BI SUA" {
		t.Error("MauOption: sua ket qua lam ban anh chup")
	}

	// cua hang — ba lop
	ch, co := a.CuaHang(1)
	if !co {
		t.Fatal("CuaHang(1) khong thay")
	}
	ch.Ten = "BI SUA"
	ch.Tabs[0].Ten = "TAB SUA"
	ch.Tabs[0].Mon[0].Gia = 999999
	ch.Tabs[0].Mon[0].Options[0].Param = 4242

	lai2, _ := a.CuaHang(1)
	if lai2.Ten == "BI SUA" {
		t.Error("CuaHang: sua lop 1 lam ban anh chup")
	}
	if lai2.Tabs[0].Ten == "TAB SUA" {
		t.Error("CuaHang: sua lop tabs lam ban anh chup")
	}
	if lai2.Tabs[0].Mon[0].Gia == 999999 {
		t.Error("CuaHang: sua lop mon lam ban anh chup (gia hang doi cho MOI nguoi choi)")
	}
	if lai2.Tabs[0].Mon[0].Options[0].Param == 4242 {
		t.Error("CuaHang: sua lop options lam ban anh chup")
	}

	// va DanhSachCuaHang cung phai sau
	dsch := a.DanhSachCuaHang()
	dsch[0].Tabs[0].Mon[0].Options[0].Param = 777
	if a.DanhSachCuaHang()[0].Tabs[0].Mon[0].Options[0].Param == 777 {
		t.Error("DanhSachCuaHang: sao chep khong du sau")
	}
}

func TestKhoTraCapChieuBienNgoai(t *testing.T) {
	a := MoKho(khoDuLieuThuNghiem(3)).AnhChup()

	if _, co := a.TraCapChieu(SoLangCoDinh, 0, 1); co {
		t.Error("lang ngoai bien phai tra false")
	}
	if _, co := a.TraCapChieu(0, 0, 0); co {
		t.Error("cap 0 phai tra false (cap dem tu 1)")
	}
	if _, co := a.TraCapChieu(0, 0, 5); co {
		t.Error("cap vuot so cap phai tra false")
	}
	if _, co := a.TraCapChieu(0, 99, 1); co {
		t.Error("template khong ton tai phai tra false")
	}
	c, co := a.TraCapChieu(0, 1, 2)
	if !co || c.Diem != 2 {
		t.Errorf("TraCapChieu(0,1,2) = %+v,%v", c, co)
	}
	if a.ChieuCuaLang(SoLangCoDinh) != nil {
		t.Error("ChieuCuaLang lang ngoai bien phai tra nil")
	}
}

// ---------------------------------------------------------------------------
// Do chi phi ban sao sau (co dung trong chu thich kho.go)
// ---------------------------------------------------------------------------

func BenchmarkKhoChieuCuaLang(b *testing.B) {
	a := MoKho(khoDuLieuThuNghiem(1)).AnhChup()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = a.ChieuCuaLang(0)
	}
}

// BenchmarkKhoChieuCuaLangThatCo — co gan voi bo du lieu that hon: 30 chieu moi
// lang, 10 cap moi chieu. Day la con so dung de bien luan trong chu thich kho.go.
func BenchmarkKhoChieuCuaLangThatCo(b *testing.B) {
	dl := khoDuLieuThuNghiem(1)
	var chieu []Chieu
	for t := 0; t < 30; t++ {
		caps := make([]CapChieu, 10)
		for c := range caps {
			caps[c] = CapChieu{Diem: byte(c + 1), ThemThongTin: "ghi chu dai vua phai"}
		}
		chieu = append(chieu, Chieu{TemplateId: byte(t), Ten: "chieu", CacCap: caps})
	}
	dl.ChieuTheoLang[0] = chieu
	a := MoKho(dl).AnhChup()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = a.ChieuCuaLang(0)
	}
}

func BenchmarkKhoTraCapChieu(b *testing.B) {
	a := MoKho(khoDuLieuThuNghiem(1)).AnhChup()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = a.TraCapChieu(0, 2, 3)
	}
}

func BenchmarkKhoCuaHang(b *testing.B) {
	a := MoKho(khoDuLieuThuNghiem(1)).AnhChup()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = a.CuaHang(1)
	}
}

func BenchmarkKhoAnhChup(b *testing.B) {
	k := MoKho(khoDuLieuThuNghiem(1))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = k.AnhChup()
	}
}
