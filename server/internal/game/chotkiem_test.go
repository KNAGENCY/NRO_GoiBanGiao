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

	kq := KiemDanhHopLe(benDanhTuNguoi(nguoiDanh), mucTieu, bayGio)

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

	kq := KiemDanhHopLe(benDanhTuNguoi(nguoiDanh), mucTieu, bayGio)

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

	kq := KiemDanhHopLe(benDanhTuNguoi(nguoiDanh), mucTieu, bayGio)

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

	kq := KiemDanhHopLe(benDanhTuNguoi(nguoiDanh), mucTieu, bayGio)

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

	kq := KiemDanhHopLe(benDanhTuNguoi(nguoiDanh), mucTieu, bayGio)

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

	kq := KiemDanhHopLe(benDanhTuNguoi(nguoiDanh), mucTieu, bayGio)

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

	kq := KiemDanhHopLe(benDanhTuNguoi(nguoiDanh), mucTieu, bayGio)

	if !kq.ChoPhep {
		t.Fatalf("khoang cach == dung tam (100 == 100) phai duoc CHO PHEP, nhung ChoPhep=false. kq=%+v", kq)
	}
	if kq.Truot {
		t.Fatalf("khoang cach == dung tam KHONG duoc tinh la truot, nhung Truot=true. kq=%+v", kq)
	}
}

// ---- Ben danh: thu tu kiem + vo han chakra ---------------------------------

// benDanhHopLeMau: ben danh "moi thu hop le" dung THANG kieu BenDanh (khong
// qua *nguoiTrongKhuVuc) — de cac test duoi day kiem HAM THUAN chu khong kiem
// adapter. Gia tri GHIM CUNG dung y het nguoiDanhHopLeMau() o tren: o (0,0),
// 10 chakra, chieu ton 10 chakra / hoi chieu 1000ms / tam 100, chua tung danh.
func benDanhHopLeMau() BenDanh {
	return BenDanh{
		X:             0,
		Y:             0,
		ChakraHienTai: 10,
		Chieu: &CapChieu{
			ChakraDung: 10,
			HoiChieuMs: 1000,
			Dx:         100,
			Dy:         100,
		},
		// LanDanhCuoi: zero value = chua tung dung -> hoi chieu tu nhien qua.
		// VoHanChakra: false = ben danh la nguoi choi binh thuong.
	}
}

// TestThuTuKiemGiuNguyen KHOA dung THU TU kiem cua KiemDanhHopLe:
// chua chon chieu -> bat tu -> chakra -> hoi chieu -> muc tieu chet -> tam.
//
// Moi ca duoi day co IT NHAT HAI dieu kien sai CUNG LUC; LyDo tra ve phai la
// ly do cua dieu kien duoc kiem TRUOC. Day khong phai chi tiet noi bo: cac
// kiem KHONG doc lap (chu thich dau chotkiem.go) — vi du kiem chakra truoc bat
// tu se cho phep tieu chakra vao mot don khong bao gio trung. Doi thu tu la
// doi hanh vi, nen thu tu la mot HOP DONG.
func TestThuTuKiemGiuNguyen(t *testing.T) {
	bayGio := time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC)

	cacCa := []struct {
		ten      string
		sua      func(*BenDanh, *ViTriCoMau)
		muonLyDo string
	}{
		{
			ten: "chua chon chieu THANG bat tu",
			sua: func(b *BenDanh, m *ViTriCoMau) {
				b.Chieu = nil                                      // hong: chua chon chieu
				m.LucHoiSinh = bayGio.Add(-100 * time.Millisecond) // hong: dang bat tu
			},
			muonLyDo: "chua chon chieu",
		},
		{
			ten: "bat tu THANG chakra",
			sua: func(b *BenDanh, m *ViTriCoMau) {
				m.LucHoiSinh = bayGio.Add(-100 * time.Millisecond) // hong: dang bat tu
				b.ChakraHienTai = 0                                // hong: thieu chakra
			},
			muonLyDo: "muc tieu dang bat tu",
		},
		{
			ten: "chakra THANG hoi chieu",
			sua: func(b *BenDanh, m *ViTriCoMau) {
				b.ChakraHienTai = 0    // hong: thieu chakra
				b.LanDanhCuoi = bayGio // hong: vua danh xong, chieu chua hoi
			},
			muonLyDo: "khong du chakra",
		},
		{
			ten: "hoi chieu THANG muc tieu chet",
			sua: func(b *BenDanh, m *ViTriCoMau) {
				b.LanDanhCuoi = bayGio // hong: chieu chua hoi
				m.HpHienTai = 0        // hong: muc tieu da chet
			},
			muonLyDo: "chieu chua hoi",
		},
		{
			ten: "muc tieu chet THANG vuot tam",
			sua: func(b *BenDanh, m *ViTriCoMau) {
				m.HpHienTai = 0 // hong: muc tieu da chet
				m.X = 500       // hong: vuot tam (tam chieu = 100)
			},
			muonLyDo: "muc tieu da chet",
		},
	}

	for _, ca := range cacCa {
		t.Run(ca.ten, func(t *testing.T) {
			benDanh := benDanhHopLeMau()
			mucTieu := mucTieuHopLeMau()
			ca.sua(&benDanh, &mucTieu)

			kq := KiemDanhHopLe(benDanh, mucTieu, bayGio)

			if kq.LyDo != ca.muonLyDo {
				t.Fatalf("LyDo = %q, muon %q (dieu kien kiem TRUOC phai thang) — thu tu kiem da bi doi. kq=%+v", kq.LyDo, ca.muonLyDo, kq)
			}
			if kq.ChoPhep {
				t.Fatalf("ca nay hong o 4 kiem dau nhung ChoPhep = true, muon false. kq=%+v", kq)
			}
			if kq.Truot {
				t.Fatalf("ca nay bi tu choi im lang, KHONG duoc bao Truot (Truot chi danh cho vuot tam). kq=%+v", kq)
			}
		})
	}
}

// TestVoHanChakraBoQuaKiemChakra: co VoHanChakra (danh cho boss — boss khong
// co bang chakra) lam kiem #2 bi BO QUA hoan toan. Hai nua cua CUNG mot tinh
// huong (0 chakra, chieu ton 10) de chac chan chinh CO NAY la thu quyet dinh,
// khong phai thu gi khac vo tinh cho qua.
func TestVoHanChakraBoQuaKiemChakra(t *testing.T) {
	bayGio := time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC)
	mucTieu := mucTieuHopLeMau()

	// Nua 1: nguoi choi (VoHanChakra = false) — 0/10 chakra thi bi tu choi.
	nguoi := benDanhHopLeMau()
	nguoi.ChakraHienTai = 0
	if kq := KiemDanhHopLe(nguoi, mucTieu, bayGio); kq.ChoPhep {
		t.Fatalf("nguoi choi 0/10 chakra phai bi tu choi, nhung ChoPhep=true. kq=%+v", kq)
	}

	// Nua 2: boss (VoHanChakra = true) — CUNG 0/10 chakra nhung duoc danh.
	boss := benDanhHopLeMau()
	boss.ChakraHienTai = 0
	boss.VoHanChakra = true
	kq := KiemDanhHopLe(boss, mucTieu, bayGio)
	if !kq.ChoPhep {
		t.Fatalf("VoHanChakra=true phai BO QUA kiem chakra, nhung bi tu choi. kq=%+v", kq)
	}
	if kq.Truot {
		t.Fatalf("don hop le khong duoc bao Truot. kq=%+v", kq)
	}
}

// TestVoHanChakraKHONGBoQuaKiemKhac: co VoHanChakra CHI mien kiem chakra —
// bat tu / hoi chieu / muc tieu chet / tam van chan boss nhu thuong. Neu ai do
// noi rong co nay thanh "bo qua het" thi test nay do.
func TestVoHanChakraKHONGBoQuaKiemKhac(t *testing.T) {
	bayGio := time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC)

	boss := benDanhHopLeMau()
	boss.VoHanChakra = true
	boss.LanDanhCuoi = bayGio.Add(-500 * time.Millisecond) // hoi chieu 1000ms -> chua hoi
	mucTieu := mucTieuHopLeMau()

	kq := KiemDanhHopLe(boss, mucTieu, bayGio)

	if kq.ChoPhep {
		t.Fatalf("VoHanChakra khong duoc mien kiem hoi chieu, nhung ChoPhep=true. kq=%+v", kq)
	}
	if kq.LyDo != "chieu chua hoi" {
		t.Fatalf("LyDo = %q, muon \"chieu chua hoi\". kq=%+v", kq.LyDo, kq)
	}
}

// ---- Adapter *nguoiTrongKhuVuc -> BenDanh ---------------------------------

// TestBenDanhTuNguoiChepDungTruong: adapter (khuvuc.go) phai chep dung cac
// truong sang BenDanh va LUON dat VoHanChakra=false — nguoi choi khong bao gio
// duoc mien chakra (day la dac quyen rieng cua boss).
func TestBenDanhTuNguoiChepDungTruong(t *testing.T) {
	moc := time.Date(2026, 7, 26, 11, 59, 0, 0, time.UTC)
	nguoi := nguoiDanhHopLeMau()
	nguoi.nv.ViTriX = 123
	nguoi.nv.ViTriY = 456
	nguoi.nv.ChakraHienTai = 77
	nguoi.lanDanhCuoi = moc

	b := benDanhTuNguoi(nguoi)

	if b.X != 123 || b.Y != 456 {
		t.Fatalf("X/Y = %d/%d, muon 123/456", b.X, b.Y)
	}
	if b.ChakraHienTai != 77 {
		t.Fatalf("ChakraHienTai = %d, muon 77", b.ChakraHienTai)
	}
	if !b.LanDanhCuoi.Equal(moc) {
		t.Fatalf("LanDanhCuoi = %v, muon %v", b.LanDanhCuoi, moc)
	}
	if b.Chieu != nguoi.chieuDangChon {
		t.Fatalf("Chieu phai tro DUNG vao chieuDangChon cua nguoi choi")
	}
	if b.VoHanChakra {
		t.Fatalf("nguoi choi phai LUON VoHanChakra=false, khong bao gio duoc mien chakra")
	}
}

// TestBenDanhTuNguoiNil: nguoi == nil (hoac chua co nv) van phai ra BenDanh co
// Chieu == nil -> KiemDanhHopLe tra "chua chon chieu", GIU NGUYEN hanh vi cua
// nhanh bao ve `nguoiDanh == nil` co truoc khi tach ham.
func TestBenDanhTuNguoiNil(t *testing.T) {
	bayGio := time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC)

	kq := KiemDanhHopLe(benDanhTuNguoi(nil), mucTieuHopLeMau(), bayGio)

	if kq.ChoPhep {
		t.Fatalf("nguoi nil phai bi tu choi, nhung ChoPhep=true. kq=%+v", kq)
	}
	if kq.LyDo != "chua chon chieu" {
		t.Fatalf("LyDo = %q, muon \"chua chon chieu\". kq=%+v", kq.LyDo, kq)
	}
}
