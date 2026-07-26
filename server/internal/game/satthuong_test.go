package game

import "testing"

// TestHeSoNguyenTo_Bang5x5 ghim cung toan bo 25 cap he theo vong khac che
// §4.2: Hoa → Phong → Loi → Tho → Thuy → Hoa. Hang = he nguoi danh,
// cot = he muc tieu.
func TestHeSoNguyenTo_Bang5x5(t *testing.T) {
	const (
		khacChe   = 1.25
		trungTinh = 1.00
		biKhac    = 0.80
	)

	// bang[nguoiDanh][mucTieu] — ghim tay tung gia tri, KHONG suy ra tu vong
	// bang cong thuc (moi diem la mot khang dinh doc lap voi cai dung trong
	// satthuong.go).
	bang := [5][5]float64{
		// muc tieu:    Hoa        Phong      Loi        Tho        Thuy
		/* Hoa   */ {trungTinh, khacChe, trungTinh, trungTinh, biKhac},
		/* Phong */ {biKhac, trungTinh, khacChe, trungTinh, trungTinh},
		/* Loi   */ {trungTinh, biKhac, trungTinh, khacChe, trungTinh},
		/* Tho   */ {trungTinh, trungTinh, biKhac, trungTinh, khacChe},
		/* Thuy  */ {khacChe, trungTinh, trungTinh, biKhac, trungTinh},
	}
	tenHe := [5]string{"Hoa", "Phong", "Loi", "Tho", "Thuy"}
	cacHe := [5]byte{HeHoa, HePhong, HeLoi, HeTho, HeThuy}

	for nd := 0; nd < 5; nd++ {
		for mt := 0; mt < 5; mt++ {
			nguoiDanh, mucTieu := cacHe[nd], cacHe[mt]
			muon := bang[nd][mt]
			got := HeSoNguyenTo(nguoiDanh, mucTieu)
			if got != muon {
				t.Errorf("HeSoNguyenTo(%s, %s) = %v, muon %v", tenHe[nd], tenHe[mt], got, muon)
			}
		}
	}
}

// TestTinhSatThuong_ApDungHeSoHeNguyenTo kiem tra cong thuc §4.1 co nhan
// dung he so he nguyen to (khong chi mang) — 1000 goc, chieu 100% sat thuong.
func TestTinhSatThuong_ApDungHeSoHeNguyenTo(t *testing.T) {
	chieu := &CapChieu{SatThuong: 100}

	cacTruongHop := []struct {
		ten                    string
		heNguoiDanh, heMucTieu byte
		muon                   int32
	}{
		{"khac che", HeHoa, HePhong, 1250},
		{"trung tinh", HeHoa, HeLoi, 1000},
		{"bi khac", HeHoa, HeThuy, 800},
	}

	for _, tc := range cacTruongHop {
		got := TinhSatThuong(1000, chieu, tc.heNguoiDanh, tc.heMucTieu, false)
		if got != tc.muon {
			t.Errorf("%s: TinhSatThuong = %d, muon %d", tc.ten, got, tc.muon)
		}
	}
}

// TestTinhSatThuong_ChiMangNhanDoi kiem tra he so chi mang la boi so 2, doc
// lap voi he so he nguyen to va he so chieu.
func TestTinhSatThuong_ChiMangNhanDoi(t *testing.T) {
	chieu := &CapChieu{SatThuong: 100}

	thuong := TinhSatThuong(1000, chieu, HeHoa, HeLoi, false) // trung tinh
	chiMang := TinhSatThuong(1000, chieu, HeHoa, HeLoi, true)

	if thuong != 1000 {
		t.Fatalf("dieu kien tien de sai: thuong = %d, muon 1000", thuong)
	}
	if chiMang != thuong*2 {
		t.Errorf("TinhSatThuong chi mang = %d, muon gap doi %d = %d", chiMang, thuong, thuong*2)
	}
}

// TestTinhSatThuong_KepTran: dau vao du lon de vuot int32 max PHAI kep dung
// ve 2_147_483_647, KHONG duoc trAn thanh so am (bay convert int64 -> int32
// khong kep truoc).
func TestTinhSatThuong_KepTran(t *testing.T) {
	chieu := &CapChieu{SatThuong: 100}

	// goc = 3 ty, he so he nguyen to khac che (x1.25) -> 3.75 ty, vuot xa
	// int32 max (~2.147 ty).
	got := TinhSatThuong(3_000_000_000, chieu, HeThuy, HeHoa, false)

	const kyVong int32 = 2_147_483_647
	if got != kyVong {
		t.Errorf("TinhSatThuong kep tran = %d, muon %d", got, kyVong)
	}
	if got < 0 {
		t.Fatalf("TinhSatThuong tran AM: %d — dau hieu tran int32 truoc khi kep", got)
	}
}

// TestTinhSatThuong_KepSan: khong bao gio ra sat thuong am, du dau vao goc am
// (vi du mot cong thuc phia tren tinh sai va truyen goc am xuong).
func TestTinhSatThuong_KepSan(t *testing.T) {
	chieu := &CapChieu{SatThuong: 100}

	got := TinhSatThuong(-100, chieu, HeHoa, HeLoi, false) // trung tinh, khong lien quan kep san
	if got != 0 {
		t.Errorf("TinhSatThuong kep san = %d, muon 0 (khong duoc am)", got)
	}
}
