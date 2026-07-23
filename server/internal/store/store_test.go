package store

import (
	"context"
	"errors"
	"os"
	"testing"
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
