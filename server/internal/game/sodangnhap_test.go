package game

// Kiem thu don vi (khong can mang/CSDL that) cho SoDangNhap — Ke hoach 11
// Task 7. Kich ban chinh (da phien cu + cho luu xong + doc lai vi tri moi
// nhat) can mang that + CSDL that (KhiDong() chi chay qua defer cua
// mynet.PhucVu), nen nam o file rieng dangnhaptrung_tichhop_test.go. File
// nay CHI kiem co che cua chinh SoDangNhap ma khong can dung toi mang/CSDL:
// dang ky khong trung khong lam gi ca, va "xoa dung phien" — canh bao rieng
// cua brief Task 7 ("cẩn thận: chỉ xoá nếu sổ đang trỏ ĐÚNG phiên đó — phiên
// cũ bị đá không được xoá mất đăng ký của phiên mới").

import "testing"

// TestSoDangNhapDangKyLanDauKhongDaAi: dang ky tai khoan CHUA co phien nao
// (coCu=false trong DangKy) phai chi don gian ghi vao so RONG roi tra ve
// NGAY — khong dong bien/kenh nao cua chinh no (biDa/daXongDon cua "moi" van
// nguyen, chua dong), vi khong co "phien cu" nao de da.
func TestSoDangNhapDangKyLanDauKhongDaAi(t *testing.T) {
	s := MoSoDangNhap()
	moi := NewBoXuLy(&khoGia{}, nil)
	var taiKhoanID int64 = 100

	s.DangKy(taiKhoanID, moi)

	s.mu.Lock()
	hienTai := s.phien[taiKhoanID]
	s.mu.Unlock()
	if hienTai != moi {
		t.Fatalf("sau DangKy lan dau, so tro toi %p, mong doi CHINH LA phien vua dang ky (%p)", hienTai, moi)
	}
	select {
	case <-moi.biDa:
		t.Fatal("dang ky lan DAU (khong trung ai) ma biDa cua chinh no da bi dong — khong co phien cu nao de da")
	default:
	}
}

// TestSoDangNhapDangKyCungPhienKhongTuDaChinhMinh: goi DangKy hai lan LIEN
// TIEP VOI CUNG MOT con tro BoXuLy (moi == cu ve dia chi) khong duoc coi la
// "trung phien can da" — nhanh "cu == moi" trong DangKy phai chan truong hop
// nay, tranh mot phien tu dong ket noi cua chinh minh mot cach vo ly.
func TestSoDangNhapDangKyCungPhienKhongTuDaChinhMinh(t *testing.T) {
	s := MoSoDangNhap()
	b := NewBoXuLy(&khoGia{}, nil)
	var taiKhoanID int64 = 200

	s.DangKy(taiKhoanID, b)
	s.DangKy(taiKhoanID, b) // goi lai VOI CHINH b — khong phai mot phien khac

	select {
	case <-b.biDa:
		t.Fatal("DangKy hai lan voi CUNG mot BoXuLy da tu da chinh no — nhanh 'cu == moi' phai chan truong hop nay")
	default:
	}
}

// TestSoDangNhapXoaNeuDungChiXoaDungPhien la kiem chung TRUC TIEP cho canh
// bao "cẩn thận" cua brief Task 7: gia lap dung tinh huong phien CU bi da
// GOI XoaNeuDung SAU KHI phien MOI da dang ky vao so (kich ban that: KhiDong()
// cua phien cu chay tren goroutine RIENG, co the tra ve SAU thoi diem
// DangKy() cua phien moi da ghi de xong s.phien[id] = moi — nhung LUON truoc
// khi DangKy tra ve, vi DangKy CHO daXongDon; o day test mo phong TRUONG HOP
// XoaNeuDung tu no goi voi b = phien KHONG con la "hien tai" nua, khong phu
// thuoc thu tu thoi gian that).
//
// Bach hoa (white-box, cung package): thao tac truc tiep s.phien thay vi di
// qua toan bo luong dang nhap that — muc dich CHI kiem logic "so sanh con
// tro roi moi xoa" cua XoaNeuDung, da co test tich hop rieng
// (dangnhaptrung_tichhop_test.go) kiem hanh vi dau-cuoi qua mang that.
func TestSoDangNhapXoaNeuDungChiXoaDungPhien(t *testing.T) {
	s := MoSoDangNhap()
	cu := NewBoXuLy(&khoGia{}, nil)
	moi := NewBoXuLy(&khoGia{}, nil)
	var taiKhoanID int64 = 42

	// Gia lap: phien MOI da dang ky xong (nhu sau khi DangKy da ghi de map).
	s.mu.Lock()
	s.phien[taiKhoanID] = moi
	s.mu.Unlock()

	// Phien CU (bi da truoc do) goi XoaNeuDung cua CHINH NO — KHONG duoc xoa
	// mat dang ky cua phien MOI, vi so dang TRO toi moi, khong phai cu.
	s.XoaNeuDung(taiKhoanID, cu)
	s.mu.Lock()
	hienTaiSauXoaCu := s.phien[taiKhoanID]
	s.mu.Unlock()
	if hienTaiSauXoaCu != moi {
		t.Fatalf("XoaNeuDung tu phien CU (bi da) da xoa/thay doi dang ky cua phien MOI (con=%p, mong doi %p) — "+
			"day chinh la lo hong brief Task 7 canh bao: phien cu bi da khong duoc xoa mat dang ky cua phien moi",
			hienTaiSauXoaCu, moi)
	}

	// Phien MOI (hien la nguoi so huu THAT SU trong so) roi di, tu goi
	// XoaNeuDung cua chinh no — PHAI xoa duoc, khong con vet trong so.
	s.XoaNeuDung(taiKhoanID, moi)
	s.mu.Lock()
	_, conSauXoaMoi := s.phien[taiKhoanID]
	s.mu.Unlock()
	if conSauXoaMoi {
		t.Fatal("XoaNeuDung tu DUNG phien dang duoc so tro toi (moi) khong xoa duoc dang ky — con sot lai 'rac' trong so")
	}
}

// TestSoDangNhapXoaNeuDungTaiKhoanKhongTonTaiKhongPanic: goi XoaNeuDung cho
// mot taiKhoanID CHUA TUNG dang ky (map rong o khoa do) phai la no-op an
// toan, khong panic — dung tinh than "khong tim thay -> im lang" cua cac
// diem tra cuu khac trong package nay (vd AnhChupNguoiChoi).
func TestSoDangNhapXoaNeuDungTaiKhoanKhongTonTaiKhongPanic(t *testing.T) {
	s := MoSoDangNhap()
	b := NewBoXuLy(&khoGia{}, nil)
	s.XoaNeuDung(999, b) // khong panic la dat
}
