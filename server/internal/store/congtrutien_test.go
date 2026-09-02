package store

import (
	"context"
	"errors"
	"testing"
)

// docTien doc thang vang/ngoc cua mot nhan vat tu CSDL (khong qua bat ky ham
// nghiep vu nao) de kiem chung ket qua that su nam trong bang.
func docTien(t *testing.T, s *Store, id int64) (int64, int64) {
	t.Helper()
	var vang, ngoc int64
	err := s.pool.QueryRow(context.Background(),
		`SELECT vang, ngoc FROM nhan_vat WHERE id = $1`, id).Scan(&vang, &ngoc)
	if err != nil {
		t.Fatalf("doc tien that bai: %v", err)
	}
	return vang, ngoc
}

// taoNhanVatThu tao mot tai khoan + nhan vat thu voi vang khoi tao cho truoc,
// tra ve nhan vat. Ngoc mac dinh 0.
func taoNhanVatThu(t *testing.T, s *Store, vang int64) *NhanVat {
	t.Helper()
	tk := taoTaiKhoanThu(t, s, "tien"+hauToNganNgau(t))
	nv, err := s.TaoNhanVat(context.Background(), NhanVatMoi{
		TaiKhoanID: tk, Ten: "Tien" + hauToNganNgau(t),
		Lang: 0, KieuToc: 5, MapID: 1, ViTriX: 100, ViTriY: 384,
		HpToiDa: 150, ChakraToiDa: 150, Vang: vang,
	})
	if err != nil {
		t.Fatalf("tao nhan vat thu loi: %v", err)
	}
	return nv
}

// TestCongTruTienCongThuong: cong duong phai tang so du dung bang delta.
func TestCongTruTienCongThuong(t *testing.T) {
	s := moKetNoiThu(t)
	ctx := context.Background()
	nv := taoNhanVatThu(t, s, 1000)

	if err := s.CongTruTien(ctx, nv.ID, 300, 50); err != nil {
		t.Fatalf("CongTruTien loi: %v", err)
	}
	vang, ngoc := docTien(t, s, nv.ID)
	if vang != 1300 {
		t.Errorf("vang = %d, mong doi 1300", vang)
	}
	if ngoc != 50 {
		t.Errorf("ngoc = %d, mong doi 50", ngoc)
	}
}

// TestCongTruTienTruThuong: tru trong pham vi so du phai giam dung bang delta.
func TestCongTruTienTruThuong(t *testing.T) {
	s := moKetNoiThu(t)
	ctx := context.Background()
	nv := taoNhanVatThu(t, s, 1000)

	if err := s.CongTruTien(ctx, nv.ID, -400, 0); err != nil {
		t.Fatalf("CongTruTien loi: %v", err)
	}
	vang, ngoc := docTien(t, s, nv.ID)
	if vang != 600 {
		t.Errorf("vang = %d, mong doi 600", vang)
	}
	if ngoc != 0 {
		t.Errorf("ngoc = %d, mong doi 0", ngoc)
	}
}

// TestCongTruTienTuChoiTruQuaSoDu la test bat bien BB4 (duong nhan ban D10):
// tru nhieu hon so du phai bi tu choi bang ErrKhongDuTien VA so du KHONG doi
// (khong am, khong ghi mot phan). Kiem ca vang lan ngoc.
func TestCongTruTienTuChoiTruQuaSoDu(t *testing.T) {
	s := moKetNoiThu(t)
	ctx := context.Background()
	nv := taoNhanVatThu(t, s, 500)

	// Vang du (500 - 100 = 400 >= 0) nhung ngoc am (0 - 10 < 0): ca cau phai
	// bi tu choi nguyen khoi, KHONG ghi vang mot minh.
	err := s.CongTruTien(ctx, nv.ID, -100, -10)
	if !errors.Is(err, ErrKhongDuTien) {
		t.Fatalf("mong doi ErrKhongDuTien, nhan %v", err)
	}
	vang, ngoc := docTien(t, s, nv.ID)
	if vang != 500 || ngoc != 0 {
		t.Errorf("so du bi doi sau khi tu choi: vang=%d ngoc=%d, mong doi 500/0", vang, ngoc)
	}

	// Tru dung bang so du (bien) phai duoc phep (500 - 500 = 0 >= 0).
	if err := s.CongTruTien(ctx, nv.ID, -500, 0); err != nil {
		t.Fatalf("tru het sach (ve 0) phai duoc phep, nhan %v", err)
	}
	vang, _ = docTien(t, s, nv.ID)
	if vang != 0 {
		t.Errorf("vang = %d, mong doi 0", vang)
	}
}

// TestCongTruTienKhongThayNhanVat: nhan vat khong ton tai phai tra
// ErrKhongThayNhanVat, phan biet ro voi ErrKhongDuTien.
func TestCongTruTienKhongThayNhanVat(t *testing.T) {
	s := moKetNoiThu(t)
	ctx := context.Background()
	nv := taoNhanVatThu(t, s, 1000)
	if _, err := s.pool.Exec(ctx, `DELETE FROM nhan_vat WHERE id = $1`, nv.ID); err != nil {
		t.Fatalf("xoa nhan vat thu loi: %v", err)
	}

	err := s.CongTruTien(ctx, nv.ID, 100, 0)
	if !errors.Is(err, ErrKhongThayNhanVat) {
		t.Fatalf("mong doi ErrKhongThayNhanVat, nhan %v", err)
	}
}

// TestThoGhiXaHangDoiKhongDeTien la test THEN CHOT cua Task 0 (duong nhan ban
// D7). Mo phong kich ban chet nguoi:
//  1. Nhan vat co 1000 vang. Actor chup mot anh (snapshot) luc nay.
//  2. Nguoi choi mua do het 400 -> CongTruTien(-400) -> CSDL con 600.
//  3. Tho ghi xa hang doi: goi CapNhatNhanVat voi ANH CHUP CU (Vang=1000).
//
// Neu CapNhatNhanVat con ghi tien tuyet doi nhu truoc, buoc 3 se ghi de 600
// tro lai 1000 -> nguoi choi co ca do lan vang (nhan ban tien). Bit dung roi
// thi tien phai van la 600.
func TestThoGhiXaHangDoiKhongDeTien(t *testing.T) {
	s := moKetNoiThu(t)
	ctx := context.Background()
	nv := taoNhanVatThu(t, s, 1000)

	// (1) Actor chup anh khi vang = 1000.
	anhChup := *nv // Vang = 1000, Ngoc = 0

	// (2) Nguoi choi mua do het 400 (con 600).
	if err := s.CongTruTien(ctx, nv.ID, -400, 0); err != nil {
		t.Fatalf("mua do (CongTruTien) loi: %v", err)
	}

	// (3) Tho ghi xa anh chup CU. Co tinh doi mot truong khac (vi tri) de
	// chac chan UPDATE that su chay, va de ANH CHUP mang Vang=1000 cu.
	anhChup.ViTriX = 999
	if err := s.CapNhatNhanVat(ctx, anhChup); err != nil {
		t.Fatalf("tho ghi (CapNhatNhanVat) loi: %v", err)
	}

	vang, ngoc := docTien(t, s, nv.ID)
	if vang != 600 {
		t.Errorf("vang = %d, mong doi 600 - tho ghi da GHI DE tien cu = nhan ban (D7)!", vang)
	}
	if ngoc != 0 {
		t.Errorf("ngoc = %d, mong doi 0", ngoc)
	}
}
