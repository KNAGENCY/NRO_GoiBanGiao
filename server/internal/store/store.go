// Package store lo viec noi chuyen voi PostgreSQL.
package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrKhongThayTaiKhoan = errors.New("khong tim thay tai khoan")

var (
	ErrKhongThayNhanVat = errors.New("khong tim thay nhan vat")
	ErrTenNhanVatDaDung = errors.New("ten nhan vat da ton tai")
	ErrDaCoNhanVat      = errors.New("tai khoan da co nhan vat")
)

type TaiKhoan struct {
	ID          int64
	TenDangNhap string
	MatKhauBam  string
	BiKhoa      bool
	LaAdmin     bool
}

// NhanVat la nhan vat da luu, du cot cot loi de vao the gioi.
type NhanVat struct {
	ID            int64
	TaiKhoanID    int64
	Ten           string
	Lang          int16
	KieuToc       int16
	MapID         int32
	ViTriX        int32
	ViTriY        int32
	Cap           int32
	KinhNghiem    int64
	SucManh       int64
	HpToiDa       int32
	HpHienTai     int32
	ChakraToiDa   int32
	ChakraHienTai int32
	Vang          int64
	Ngoc          int64
}

// NhanVatMoi la du lieu can de tao mot nhan vat. Cac chi so suy ra (cap ban dau,
// mau/chakra hien tai = toi da) do TaoNhanVat tu dat, khong nhan tu ben ngoai.
type NhanVatMoi struct {
	TaiKhoanID  int64
	Ten         string
	Lang        int16
	KieuToc     int16
	MapID       int32
	ViTriX      int32
	ViTriY      int32
	HpToiDa     int32
	ChakraToiDa int32
	Vang        int64
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

// TaoNhanVat them mot nhan vat moi. Mau/chakra hien tai khoi tao bang toi da,
// cap = 1. Neu ten trung (khong phan biet hoa thuong) tra ErrTenNhanVatDaDung;
// neu tai khoan da co nhan vat tra ErrDaCoNhanVat — de tang tren doi hai loi
// nay ra hai cau bao khac nhau cho nguoi choi.
func (s *Store) TaoNhanVat(ctx context.Context, nv NhanVatMoi) (*NhanVat, error) {
	kq := &NhanVat{}
	err := s.pool.QueryRow(ctx,
		`INSERT INTO nhan_vat
		    (tai_khoan_id, ten, lang, kieu_toc, map_id, vi_tri_x, vi_tri_y,
		     hp_toi_da, hp_hien_tai, chakra_toi_da, chakra_hien_tai, vang)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$8,$9,$9,$10)
		 RETURNING id, tai_khoan_id, ten, lang, kieu_toc, map_id, vi_tri_x, vi_tri_y,
		           cap, kinh_nghiem, suc_manh, hp_toi_da, hp_hien_tai,
		           chakra_toi_da, chakra_hien_tai, vang, ngoc`,
		nv.TaiKhoanID, nv.Ten, nv.Lang, nv.KieuToc, nv.MapID, nv.ViTriX, nv.ViTriY,
		nv.HpToiDa, nv.ChakraToiDa, nv.Vang,
	).Scan(
		&kq.ID, &kq.TaiKhoanID, &kq.Ten, &kq.Lang, &kq.KieuToc, &kq.MapID,
		&kq.ViTriX, &kq.ViTriY, &kq.Cap, &kq.KinhNghiem, &kq.SucManh,
		&kq.HpToiDa, &kq.HpHienTai, &kq.ChakraToiDa, &kq.ChakraHienTai,
		&kq.Vang, &kq.Ngoc,
	)
	if err != nil {
		// Phan biet hai rang buoc duy nhat khac nhau bang TEN rang buoc, khong
		// chi bang ma loi 23505 (ca hai deu 23505). Nho vay tang tren bao dung
		// "ten da co" hay "da co nhan vat".
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "nhan_vat_ten_duy_nhat":
				return nil, ErrTenNhanVatDaDung
			case "nhan_vat_mot_tai_khoan":
				return nil, ErrDaCoNhanVat
			}
		}
		return nil, fmt.Errorf("them nhan vat that bai: %w", err)
	}
	return kq, nil
}

// LayNhanVatTheoTaiKhoan tra nhan vat cua mot tai khoan, hoac ErrKhongThayNhanVat.
func (s *Store) LayNhanVatTheoTaiKhoan(ctx context.Context, taiKhoanID int64) (*NhanVat, error) {
	kq := &NhanVat{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, tai_khoan_id, ten, lang, kieu_toc, map_id, vi_tri_x, vi_tri_y,
		        cap, kinh_nghiem, suc_manh, hp_toi_da, hp_hien_tai,
		        chakra_toi_da, chakra_hien_tai, vang, ngoc
		   FROM nhan_vat
		  WHERE tai_khoan_id = $1
		  LIMIT 1`, taiKhoanID,
	).Scan(
		&kq.ID, &kq.TaiKhoanID, &kq.Ten, &kq.Lang, &kq.KieuToc, &kq.MapID,
		&kq.ViTriX, &kq.ViTriY, &kq.Cap, &kq.KinhNghiem, &kq.SucManh,
		&kq.HpToiDa, &kq.HpHienTai, &kq.ChakraToiDa, &kq.ChakraHienTai,
		&kq.Vang, &kq.Ngoc,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrKhongThayNhanVat
	}
	if err != nil {
		return nil, fmt.Errorf("tra cuu nhan vat that bai: %w", err)
	}
	return kq, nil
}

// TenNhanVatDaDung kiem nhanh mot ten da co ai dung chua (khong phan biet hoa
// thuong). Dung de bao nguoi choi som; rang buoc duy nhat trong CSDL moi la
// chot chan that su khi co dua tranh.
func (s *Store) TenNhanVatDaDung(ctx context.Context, ten string) (bool, error) {
	var co bool
	err := s.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM nhan_vat WHERE lower(ten) = lower($1))`, ten,
	).Scan(&co)
	if err != nil {
		return false, fmt.Errorf("kiem ten nhan vat that bai: %w", err)
	}
	return co, nil
}
