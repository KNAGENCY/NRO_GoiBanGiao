package store

import (
	"context"
	"testing"
)

// moKhoCoBangNo mo ket noi thu VA ap cac migration con thieu — KHAC moKetNoiThu
// (store_test.go, khong ap migration). Bat buoc phai co o day vi bang
// thuong_cho_nhan moi ra doi o 0004: mot CSDL nro_test da dung tu truoc se con
// dung o 0003 va moi cau SQL duoi day se hong voi loi "relation does not exist"
// — mot kieu do RAT de bi doc nham thanh "ham viet sai".
func moKhoCoBangNo(t *testing.T) *Store {
	t.Helper()
	s := moKetNoiThu(t)
	if _, err := s.ApDungMigrationConLai(context.Background()); err != nil {
		t.Fatalf("ap migration: %v", err)
	}
	return s
}

// taoNguoiThu tao mot tai khoan + nhan vat THAT (khong phai gia) roi dang ky
// don dep. Moi nguoi thu mot tai khoan rieng vi CSDL rang buoc mot tai khoan
// chi mot nhan vat (nhan_vat_mot_tai_khoan).
//
// Chi xoa tai_khoan: nhan_vat co khoa ngoai ON DELETE CASCADE toi tai_khoan
// khong? — KHONG, nen phai xoa nhan_vat truoc. t.Cleanup chay NGUOC thu tu dang
// ky, vi vay dang ky "xoa tai khoan" TRUOC roi "xoa nhan vat" SAU thi luc chay
// nhan vat bi xoa truoc, tai khoan sau — dung thu tu khoa ngoai. So no thi
// khong phai don: thuong_cho_nhan CASCADE theo nhan_vat.
func taoNguoiThu(t *testing.T, s *Store, ten string) *NhanVat {
	t.Helper()
	ctx := context.Background()

	var idTaiKhoan int64
	if err := s.pool.QueryRow(ctx,
		`INSERT INTO tai_khoan (ten_dang_nhap, mat_khau_bam) VALUES ($1, $2) RETURNING id`,
		ten, "bam-gia").Scan(&idTaiKhoan); err != nil {
		t.Fatalf("tao tai khoan %q: %v", ten, err)
	}
	t.Cleanup(func() {
		s.pool.Exec(context.Background(), `DELETE FROM tai_khoan WHERE id = $1`, idTaiKhoan)
	})

	nv, err := s.TaoNhanVat(ctx, NhanVatMoi{
		TaiKhoanID:  idTaiKhoan,
		Ten:         ten,
		Lang:        1,
		KieuToc:     1,
		MapID:       0,
		ViTriX:      100,
		ViTriY:      200,
		HpToiDa:     1000,
		ChakraToiDa: 500,
		Vang:        0,
	})
	if err != nil {
		t.Fatalf("tao nhan vat %q: %v", ten, err)
	}
	t.Cleanup(func() {
		s.pool.Exec(context.Background(), `DELETE FROM nhan_vat WHERE id = $1`, nv.ID)
	})
	return nv
}

// docKinhNghiem doc thang cot kinh_nghiem tu CSDL — KHONG qua
// LayNhanVatTheoTaiKhoan, de bai kiem khong phu thuoc vao mot ham khac cung co
// the sai (§6.12: dung gia tri ghim, dung lay tu chinh nguon dang kiem).
func docKinhNghiem(t *testing.T, s *Store, idNhanVat int64) int64 {
	t.Helper()
	var kn int64
	if err := s.pool.QueryRow(context.Background(),
		`SELECT kinh_nghiem FROM nhan_vat WHERE id = $1`, idNhanVat).Scan(&kn); err != nil {
		t.Fatalf("doc kinh_nghiem: %v", err)
	}
	return kn
}

// demNoChuaNhan dem so dong no CHUA nhan cua mot nguoi.
func demNoChuaNhan(t *testing.T, s *Store, idNhanVat int64) int {
	t.Helper()
	var n int
	if err := s.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM thuong_cho_nhan WHERE nhan_vat_id = $1 AND da_nhan = false`,
		idNhanVat).Scan(&n); err != nil {
		t.Fatalf("dem no chua nhan: %v", err)
	}
	return n
}

// TestNhanNoThuong_CongVaoKinhNghiemRoiDanhDauDaNhan la bai kiem TRUNG TAM cua
// ca tinh nang: hai khoan no rieng le phai cong DUNG TONG vao cot kinh_nghiem,
// va phai bi danh dau da nhan trong CUNG mot lan.
func TestNhanNoThuong_CongVaoKinhNghiemRoiDanhDauDaNhan(t *testing.T) {
	s := moKhoCoBangNo(t)
	ctx := context.Background()
	nv := taoNguoiThu(t, s, "notestaa")

	if err := s.GhiNoThuong(ctx, []NoThuong{
		{NhanVatID: nv.ID, Loai: LoaiThuongKinhNghiem, SoDiem: 300},
		{NhanVatID: nv.ID, Loai: LoaiThuongKinhNghiem, SoDiem: 200},
	}); err != nil {
		t.Fatalf("ghi no: %v", err)
	}
	if got := demNoChuaNhan(t, s, nv.ID); got != 2 {
		t.Fatalf("sau khi ghi: %d dong no chua nhan, mong 2", got)
	}

	tong, err := s.NhanNoThuong(ctx, nv.ID)
	if err != nil {
		t.Fatalf("nhan no: %v", err)
	}
	if tong != 500 {
		t.Fatalf("NhanNoThuong tra %d, mong 500", tong)
	}
	if got := docKinhNghiem(t, s, nv.ID); got != 500 {
		t.Fatalf("kinh_nghiem = %d, mong 500", got)
	}
	if got := demNoChuaNhan(t, s, nv.ID); got != 0 {
		t.Fatalf("sau khi nhan: con %d dong no chua nhan, mong 0", got)
	}
}

// TestNhanNoThuong_GoiLanHaiKhongNhanDoi khoa lo hong dat nhat cua tinh nang:
// mot khoan thuong duoc phat HAI LAN la in tien tu khong khi. Day cung la bai
// kiem se DO neu ai do bo menh de "da_nhan = false" hoac quen buoc danh dau.
func TestNhanNoThuong_GoiLanHaiKhongNhanDoi(t *testing.T) {
	s := moKhoCoBangNo(t)
	ctx := context.Background()
	nv := taoNguoiThu(t, s, "notestbb")

	if err := s.GhiNoThuong(ctx, []NoThuong{
		{NhanVatID: nv.ID, Loai: LoaiThuongKinhNghiem, SoDiem: 777},
	}); err != nil {
		t.Fatalf("ghi no: %v", err)
	}
	if _, err := s.NhanNoThuong(ctx, nv.ID); err != nil {
		t.Fatalf("nhan lan 1: %v", err)
	}

	tong2, err := s.NhanNoThuong(ctx, nv.ID)
	if err != nil {
		t.Fatalf("nhan lan 2: %v", err)
	}
	if tong2 != 0 {
		t.Fatalf("nhan lan 2 tra %d, mong 0 (khong nhan doi)", tong2)
	}
	if got := docKinhNghiem(t, s, nv.ID); got != 777 {
		t.Fatalf("kinh_nghiem = %d sau hai lan nhan, mong 777", got)
	}
}

// TestNhanNoThuong_KhongCoNoThiKhongDoiGi: duong nhan chay o MOI lan dang nhap,
// nen truong hop pho bien nhat la KHONG CO NO — no phai khong lam gi ca, khong
// loi, va tuyet doi khong cham vao kinh_nghiem.
func TestNhanNoThuong_KhongCoNoThiKhongDoiGi(t *testing.T) {
	s := moKhoCoBangNo(t)
	ctx := context.Background()
	nv := taoNguoiThu(t, s, "notestcc")

	if _, err := s.pool.Exec(ctx,
		`UPDATE nhan_vat SET kinh_nghiem = 12345 WHERE id = $1`, nv.ID); err != nil {
		t.Fatalf("dat kinh nghiem san: %v", err)
	}

	tong, err := s.NhanNoThuong(ctx, nv.ID)
	if err != nil {
		t.Fatalf("nhan no khi khong co no: %v", err)
	}
	if tong != 0 {
		t.Fatalf("tra %d, mong 0", tong)
	}
	if got := docKinhNghiem(t, s, nv.ID); got != 12345 {
		t.Fatalf("kinh_nghiem = %d, mong giu nguyen 12345", got)
	}
}

// TestNhanNoThuong_KhongDungToiNoNguoiKhac: menh de WHERE phai loc theo dung
// nguoi. Khong co bai kiem nay thi mot lan quen dieu kien nhan_vat_id se lam
// nguoi dang nhap dau tien om het thuong cua ca may chu — va bo test van xanh.
func TestNhanNoThuong_KhongDungToiNoNguoiKhac(t *testing.T) {
	s := moKhoCoBangNo(t)
	ctx := context.Background()
	a := taoNguoiThu(t, s, "notestdd")
	b := taoNguoiThu(t, s, "notestee")

	if err := s.GhiNoThuong(ctx, []NoThuong{
		{NhanVatID: a.ID, Loai: LoaiThuongKinhNghiem, SoDiem: 100},
		{NhanVatID: b.ID, Loai: LoaiThuongKinhNghiem, SoDiem: 900},
	}); err != nil {
		t.Fatalf("ghi no: %v", err)
	}

	tong, err := s.NhanNoThuong(ctx, a.ID)
	if err != nil {
		t.Fatalf("nhan no cua a: %v", err)
	}
	if tong != 100 {
		t.Fatalf("a nhan %d, mong 100", tong)
	}
	if got := docKinhNghiem(t, s, b.ID); got != 0 {
		t.Fatalf("kinh_nghiem cua b = %d, mong 0 (khong bi ai lay)", got)
	}
	if got := demNoChuaNhan(t, s, b.ID); got != 1 {
		t.Fatalf("b con %d dong no chua nhan, mong 1", got)
	}
}

// TestNhanNoThuong_BaoHoaKhongTran: kinh_nghiem la BIGINT, cong them mot khoan
// no vao mot nguoi da gan tran se lam PostgreSQL nem loi 22003 va HUY ca giao
// dich — nguoi do se khong bao gio dang nhap duoc vao the gioi voi thuong da
// tra. Phai bao hoa o MaxInt64 thay vi tran.
func TestNhanNoThuong_BaoHoaKhongTran(t *testing.T) {
	s := moKhoCoBangNo(t)
	ctx := context.Background()
	nv := taoNguoiThu(t, s, "notestff")

	const ganTran = int64(9_223_372_036_854_775_800) // MaxInt64 - 7
	if _, err := s.pool.Exec(ctx,
		`UPDATE nhan_vat SET kinh_nghiem = $1 WHERE id = $2`, ganTran, nv.ID); err != nil {
		t.Fatalf("dat kinh nghiem gan tran: %v", err)
	}
	if err := s.GhiNoThuong(ctx, []NoThuong{
		{NhanVatID: nv.ID, Loai: LoaiThuongKinhNghiem, SoDiem: 1_000_000},
	}); err != nil {
		t.Fatalf("ghi no: %v", err)
	}

	if _, err := s.NhanNoThuong(ctx, nv.ID); err != nil {
		t.Fatalf("nhan no gan tran: %v", err)
	}
	const maxInt64 = int64(9_223_372_036_854_775_807)
	if got := docKinhNghiem(t, s, nv.ID); got != maxInt64 {
		t.Fatalf("kinh_nghiem = %d, mong bao hoa o %d", got, maxInt64)
	}
	if got := demNoChuaNhan(t, s, nv.ID); got != 0 {
		t.Fatalf("con %d dong no chua nhan sau khi bao hoa, mong 0", got)
	}
}

// TestGhiNoThuong_LoRongKhongLamGi: actor co the goi voi lo rong (khong ai vang
// mat luc boss chet) — duong pho bien nhat. Phai khong loi va khong sinh dong
// rac nao.
func TestGhiNoThuong_LoRongKhongLamGi(t *testing.T) {
	s := moKhoCoBangNo(t)
	if err := s.GhiNoThuong(context.Background(), nil); err != nil {
		t.Fatalf("ghi lo rong: %v", err)
	}
}

// TestGhiNoThuong_TuChoiSoDiemKhongDuong: rang buoc so_diem > 0 la chot chan
// cuoi cung chong mot khoan no AM (se TRU kinh nghiem cua nguoi choi khi nhan).
func TestGhiNoThuong_TuChoiSoDiemKhongDuong(t *testing.T) {
	s := moKhoCoBangNo(t)
	ctx := context.Background()
	nv := taoNguoiThu(t, s, "notestgg")

	for _, diem := range []int64{0, -50} {
		if err := s.GhiNoThuong(ctx, []NoThuong{
			{NhanVatID: nv.ID, Loai: LoaiThuongKinhNghiem, SoDiem: diem},
		}); err == nil {
			t.Fatalf("ghi no voi so diem %d KHONG loi, mong bi tu choi", diem)
		}
	}
	if got := demNoChuaNhan(t, s, nv.ID); got != 0 {
		t.Fatalf("co %d dong no duoc ghi du bi tu choi, mong 0", got)
	}
}
