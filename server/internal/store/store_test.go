package store

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

// Kiem thu nay can mot CSDL PostgreSQL that. Dat bien moi truong:
//
//	NRO_TEST_DSN=postgres://postgres:nrodev123@127.0.0.1:5433/nro_test
//
// Khong dat thi bo qua, de nguoi khac chay `go test ./...` khong bi do.
func moKetNoiThu(t *testing.T) *Store {
	t.Helper()
	dsn := os.Getenv("NRO_TEST_DSN")
	if dsn == "" {
		t.Skip("bo qua: chua dat NRO_TEST_DSN")
	}
	s, err := Open(context.Background(), dsn)
	if err != nil {
		t.Fatalf("khong mo duoc CSDL: %v", err)
	}
	t.Cleanup(s.Close)
	return s
}

func TestLayTaiKhoanKhongTonTai(t *testing.T) {
	s := moKetNoiThu(t)
	_, err := s.LayTaiKhoan(context.Background(), "khong-he-ton-tai-xyz")
	if !errors.Is(err, ErrKhongThayTaiKhoan) {
		t.Fatalf("err = %v, mong doi ErrKhongThayTaiKhoan", err)
	}
}

func TestLayTaiKhoanKhongPhanBietHoaThuong(t *testing.T) {
	s := moKetNoiThu(t)
	ctx := context.Background()

	_, err := s.pool.Exec(ctx,
		`INSERT INTO tai_khoan (ten_dang_nhap, mat_khau_bam) VALUES ($1, $2)
		 ON CONFLICT DO NOTHING`, "NguoiThu", "bam-gia")
	if err != nil {
		t.Fatalf("khong chen duoc du lieu thu: %v", err)
	}
	t.Cleanup(func() {
		s.pool.Exec(ctx, `DELETE FROM tai_khoan WHERE lower(ten_dang_nhap) = 'nguoithu'`)
	})

	for _, ten := range []string{"NguoiThu", "nguoithu", "NGUOITHU"} {
		tk, err := s.LayTaiKhoan(ctx, ten)
		if err != nil {
			t.Fatalf("tim %q that bai: %v", ten, err)
		}
		if tk.MatKhauBam != "bam-gia" {
			t.Fatalf("lay nham tai khoan cho %q", ten)
		}
	}
}

// chenThuTaiKhoan chen mot dong vao tai_khoan de kiem tra rang buoc CSDL, va
// tu dang ky don dep bang t.Cleanup (khop hoac khong khop deu khong loi).
func chenThuTaiKhoan(t *testing.T, s *Store, ten string) error {
	t.Helper()
	ctx := context.Background()
	_, err := s.pool.Exec(ctx,
		`INSERT INTO tai_khoan (ten_dang_nhap, mat_khau_bam) VALUES ($1, $2)`,
		ten, "bam-gia")
	t.Cleanup(func() {
		s.pool.Exec(context.Background(),
			`DELETE FROM tai_khoan WHERE lower(ten_dang_nhap) = lower($1)`, ten)
	})
	return err
}

// TestChenTrungTenKhongPhanBietHoaThuongBiTuChoi kiem tra co che chan trung
// thuc su o tang CSDL: neu xoa chi muc duy nhat tai_khoan_ten_duy_nhat di,
// test nay se that bai (khac voi test doc du lieu o tren, chi INSERT mot
// dong roi doc lai nen khong phat hien duoc mat chi muc).
func TestChenTrungTenKhongPhanBietHoaThuongBiTuChoi(t *testing.T) {
	s := moKetNoiThu(t)

	if err := chenThuTaiKhoan(t, s, "Abc9"); err != nil {
		t.Fatalf("chen lan dau that bai: %v", err)
	}

	err := chenThuTaiKhoan(t, s, "abc9")
	if err == nil {
		t.Fatalf(`chen "abc9" trung voi "Abc9" (chi khac hoa thuong) phai bi CSDL tu choi nhung lai thanh cong`)
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		t.Fatalf("mong doi loi trung khoa (SQLSTATE 23505), nhan duoc: %v", err)
	}
}

// TestChenTenDangNhapCoDauBiTuChoi kiem tra rang buoc CHECK moi:
// ten_dang_nhap co dau tieng Viet (ngoai bang chu cai tieng Anh) phai bi
// CSDL tu choi, vi duoi locale C ham lower() khong xu ly dung cac ky tu nay,
// se pha vo quy tac "khong phan biet hoa thuong" neu duoc lot vao.
func TestChenTenDangNhapCoDauBiTuChoi(t *testing.T) {
	s := moKetNoiThu(t)

	err := chenThuTaiKhoan(t, s, "Ănguoi")
	if err == nil {
		t.Fatalf(`ten dang nhap co dau tieng Viet "Ănguoi" phai bi rang buoc CHECK tu choi nhung lai thanh cong`)
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23514" {
		t.Fatalf("mong doi loi vi pham CHECK (SQLSTATE 23514), nhan duoc: %v", err)
	}
}

// TestChenTenDangNhapDoDaiKhongHopLeBiTuChoi kiem tra rang buoc CHECK moi ve
// do dai: ten_dang_nhap phai tu 4 den 20 ky tu.
func TestChenTenDangNhapDoDaiKhongHopLeBiTuChoi(t *testing.T) {
	s := moKetNoiThu(t)

	truongHop := []struct {
		mota string
		ten  string
	}{
		{"qua ngan (3 ky tu)", "abc"},
		{"qua dai (21 ky tu)", strings.Repeat("x", 21)},
	}
	for _, tc := range truongHop {
		err := chenThuTaiKhoan(t, s, tc.ten)
		if err == nil {
			t.Fatalf("ten dang nhap %s phai bi rang buoc CHECK tu choi nhung lai thanh cong", tc.mota)
		}
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || pgErr.Code != "23514" {
			t.Fatalf("mong doi loi vi pham CHECK (SQLSTATE 23514) cho truong hop %s, nhan duoc: %v", tc.mota, err)
		}
	}
}

// TestChenTenDangNhapHopLeThanhCong la doi chung cua hai test rang buoc CHECK
// o tren: ten dang nhap dung quy tac (chu, so, gach duoi, dai 4-20) phai
// duoc chap nhan binh thuong.
func TestChenTenDangNhapHopLeThanhCong(t *testing.T) {
	s := moKetNoiThu(t)

	if err := chenThuTaiKhoan(t, s, "nguoi_choi_01"); err != nil {
		t.Fatalf("ten dang nhap hop le phai duoc chap nhan nhung lai bi tu choi: %v", err)
	}
}
