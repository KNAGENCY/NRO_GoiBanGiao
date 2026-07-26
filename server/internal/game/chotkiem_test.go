package game

import (
	"testing"
	"time"

	"nroserver/internal/store"
)

// nguoiDanhHopLeMau dung chung cho cac test: mot nguoi choi o (0,0), chieu
// dang chon co tam 100, ton 10 chakra, hoi chieu 1000ms, chua tung dung chieu
// nao (lanDanhCuoi la zero value). Cac gia tri o day GHIM CUNG, KHONG lay tu
// hang so thoiGianBatTuSauHoiSinh / hoiChieuNoi cua chotkiem.go — de test that
// su kiem duoc GIA TRI, khong kiem "code tu so voi chinh no".
func nguoiDanhHopLeMau() *nguoiTrongKhuVuc {
	return &nguoiTrongKhuVuc{
		nv: &store.NhanVat{
			ID:            1,
			ViTriX:        0,
			ViTriY:        0,
			ChakraHienTai: 10,
		},
		chieuDangChon: &CapChieu{
			ChakraDung: 10,
			HoiChieuMs: 1000,
			Dx:         100,
			Dy:         100,
		},
		// lanDanhCuoi: zero value = chua tung dung -> hoi chieu tu nhien qua.
	}
}

// mucTieuHopLeMau: muc tieu con song, o (0,0) (trung vi tri nguoiDanh -> luon
// trong tam), hoi sinh tu rat lau (LucHoiSinh zero -> khong bat tu).
func mucTieuHopLeMau() ViTriCoMau {
	return ViTriCoMau{X: 0, Y: 0, HpHienTai: 100}
}

// ---- Dong 1: bat tu ----------------------------------------------------

// TestKiemDanhBatTu: muc tieu hoi sinh 1000ms truoc bayGio (< 1500ms) ->
// ChoPhep phai la false. Moi kiem khac deu hop le, chi rieng bat tu hong.
func TestKiemDanhBatTu(t *testing.T) {
	bayGio := time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC)
	nguoiDanh := nguoiDanhHopLeMau()
	mucTieu := mucTieuHopLeMau()
	mucTieu.LucHoiSinh = bayGio.Add(-1000 * time.Millisecond)

	kq := KiemDanhHopLe(nguoiDanh, mucTieu, bayGio)

	if kq.ChoPhep {
		t.Fatalf("muc tieu vua hoi sinh 1000ms truoc (< 1500ms bat tu) nhung ChoPhep = true, muon false. kq=%+v", kq)
	}
}

// ---- Dong 2: chakra -----------------------------------------------------

// TestKiemDanhKhongDuChakra: nguoiDanh chi con 5 chakra, chieu ton 10 ->
// ChoPhep phai la false.
func TestKiemDanhKhongDuChakra(t *testing.T) {
	bayGio := time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC)
	nguoiDanh := nguoiDanhHopLeMau()
	nguoiDanh.nv.ChakraHienTai = 5
	nguoiDanh.chieuDangChon.ChakraDung = 10
	mucTieu := mucTieuHopLeMau()

	kq := KiemDanhHopLe(nguoiDanh, mucTieu, bayGio)

	if kq.ChoPhep {
		t.Fatalf("chi con 5/10 chakra nhung ChoPhep = true, muon false. kq=%+v", kq)
	}
}

// ---- Dong 3: hoi chieu ----------------------------------------------------

// TestKiemDanhChieuChuaHoi: lan dung truoc 500ms, hoi chieu 1000ms (noi 50ms
// -> nguong 950ms) -> 500 < 950 -> ChoPhep phai la false.
func TestKiemDanhChieuChuaHoi(t *testing.T) {
	bayGio := time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC)
	nguoiDanh := nguoiDanhHopLeMau()
	nguoiDanh.chieuDangChon.HoiChieuMs = 1000
	nguoiDanh.lanDanhCuoi = bayGio.Add(-500 * time.Millisecond)
	mucTieu := mucTieuHopLeMau()

	kq := KiemDanhHopLe(nguoiDanh, mucTieu, bayGio)

	if kq.ChoPhep {
		t.Fatalf("chieu con hoi chieu (500ms/1000ms) nhung ChoPhep = true, muon false. kq=%+v", kq)
	}
}

// ---- Dong 4: muc tieu da chet --------------------------------------------

// TestKiemDanhMucTieuDaChet: HpHienTai = 0 -> ChoPhep phai la false.
func TestKiemDanhMucTieuDaChet(t *testing.T) {
	bayGio := time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC)
	nguoiDanh := nguoiDanhHopLeMau()
	mucTieu := mucTieuHopLeMau()
	mucTieu.HpHienTai = 0

	kq := KiemDanhHopLe(nguoiDanh, mucTieu, bayGio)

	if kq.ChoPhep {
		t.Fatalf("muc tieu da chet (HpHienTai=0) nhung ChoPhep = true, muon false. kq=%+v", kq)
	}
}

// ---- Dong 5: vuot tam -----------------------------------------------------

// TestKiemDanhVuotTam: nguoiDanh o (0,0), muc tieu o (500,0), tam chieu = 100
// -> vuot tam. Phai ChoPhep=true, Truot=true (dac ta §1.9 + §3.3 diem 3 —
// van phat goi nhung bao trung "hut", khong im lang bo qua).
func TestKiemDanhVuotTam(t *testing.T) {
	bayGio := time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC)
	nguoiDanh := nguoiDanhHopLeMau()
	nguoiDanh.chieuDangChon.Dx = 100
	mucTieu := mucTieuHopLeMau()
	mucTieu.X = 500
	mucTieu.Y = 0

	kq := KiemDanhHopLe(nguoiDanh, mucTieu, bayGio)

	if !kq.ChoPhep {
		t.Fatalf("vuot tam phai VAN ChoPhep=true (van phat goi, chi trot), nhung duoc false. kq=%+v", kq)
	}
	if !kq.Truot {
		t.Fatalf("vuot tam phai Truot=true, nhung duoc false. kq=%+v", kq)
	}
}

// ---- Bien 1: daTroiQua == hoiChieu - 50ms DUNG BANG -----------------------

// TestKiemDanhBienHoiChieuDungBang: daTroiQua dung bang nguong (hoiChieu -
// 50ms) -> phai CHO PHEP (dac ta dung ">=" khong phai ">").
func TestKiemDanhBienHoiChieuDungBang(t *testing.T) {
	bayGio := time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC)
	nguoiDanh := nguoiDanhHopLeMau()
	nguoiDanh.chieuDangChon.HoiChieuMs = 1000
	// nguong = 1000 - 50 = 950ms. dat lanDanhCuoi dung 950ms truoc bayGio.
	nguoiDanh.lanDanhCuoi = bayGio.Add(-950 * time.Millisecond)
	mucTieu := mucTieuHopLeMau()

	kq := KiemDanhHopLe(nguoiDanh, mucTieu, bayGio)

	if !kq.ChoPhep {
		t.Fatalf("daTroiQua == hoiChieu-50ms (950ms == 950ms) phai duoc CHO PHEP, nhung ChoPhep=false. kq=%+v", kq)
	}
}

// ---- Bien 2: khoang cach == dung tam --------------------------------------

// TestKiemDanhBienTamDungBang: khoang cach dung bang tam chieu -> phai CHO
// PHEP va KHONG TRUOT (dac ta dung "<=" trong tam, khong phai "<").
func TestKiemDanhBienTamDungBang(t *testing.T) {
	bayGio := time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC)
	nguoiDanh := nguoiDanhHopLeMau()
	nguoiDanh.chieuDangChon.Dx = 100
	mucTieu := mucTieuHopLeMau()
	mucTieu.X = 100 // khoang cach dung bang tam = 100
	mucTieu.Y = 0

	kq := KiemDanhHopLe(nguoiDanh, mucTieu, bayGio)

	if !kq.ChoPhep {
		t.Fatalf("khoang cach == dung tam (100 == 100) phai duoc CHO PHEP, nhung ChoPhep=false. kq=%+v", kq)
	}
	if kq.Truot {
		t.Fatalf("khoang cach == dung tam KHONG duoc tinh la truot, nhung Truot=true. kq=%+v", kq)
	}
}
