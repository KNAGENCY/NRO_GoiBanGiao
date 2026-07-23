// Package store lo viec noi chuyen voi PostgreSQL.
package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrKhongThayTaiKhoan = errors.New("khong tim thay tai khoan")

type TaiKhoan struct {
	ID          int64
	TenDangNhap string
	MatKhauBam  string
	BiKhoa      bool
	LaAdmin     bool
}

type Store struct {
	pool *pgxpool.Pool
}

func Open(ctx context.Context, dsn string) (*Store, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("tao pool that bai: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("khong ket noi duoc CSDL: %w", err)
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close() { s.pool.Close() }

// LayTaiKhoan tim theo ten dang nhap, khong phan biet hoa thuong.
func (s *Store) LayTaiKhoan(ctx context.Context, ten string) (*TaiKhoan, error) {
	var tk TaiKhoan
	err := s.pool.QueryRow(ctx,
		`SELECT id, ten_dang_nhap, mat_khau_bam, bi_khoa, la_admin
		   FROM tai_khoan
		  WHERE lower(ten_dang_nhap) = lower($1)
		  LIMIT 1`, ten,
	).Scan(&tk.ID, &tk.TenDangNhap, &tk.MatKhauBam, &tk.BiKhoa, &tk.LaAdmin)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrKhongThayTaiKhoan
	}
	if err != nil {
		return nil, fmt.Errorf("truy van tai khoan that bai: %w", err)
	}
	return &tk, nil
}
