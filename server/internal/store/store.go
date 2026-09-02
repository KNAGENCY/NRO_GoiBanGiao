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
	ErrKhongDuTien      = errors.New("khong du tien")
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

// CapNhatNhanVat ghi lai TOAN BO trang thai nhan vat hien tai xuong CSDL -
// khong chi 4 truong dang doi hom nay (vi tri, mau, chakra o KH6) ma ca cac
// truong KH7/KH8 sap noi day (Cap, KinhNghiem, SucManh, MapID...) - TRU tien
// (vang/ngoc), xem ghi chu BB4 o duoi. Ghi het struct tu dau thi sau nay them
// truong moi vao NhanVat khong phai sua lai cau SQL nay.
//
// KHONG ghi ID, TaiKhoanID, Ten trong menh de SET: ba truong nay BAT BIEN sau
// khi TaoNhanVat tao xong (dinh danh hang, chu so huu, ten hien thi voi rang
// buoc duy nhat rieng) - cho phep ghi lai se mo duong cho loi doi ten/doi chu
// nhan vat neu tang tren goi ham nay voi mot ban NhanVat bi lech du lieu. ID
// chi dung o menh de WHERE, khong bao gio dung ben trai dau '=' cua SET.
//
// dang_nhap_cuoi duoc ghi bang now() cua CSDL (khong phai gio may Go, tranh
// lech dong ho giua cac tien trinh) trong MOI lan goi ham nay - ke ca khi goi
// tu vong lap luu dinh ky, khong chi luc dang nhap/dang xuat that su. Cot nay
// co san tu 0002_nhan_vat.sql nhung chua ai dung toi.
//
// Nhan vat khong ton tai (id sai, hoac vua bi xoa) tra ErrKhongThayNhanVat
// thay vi im lang coi nhu xong (bai hoc dat tu ban Java cu - loi ghi bi nuot,
// xem dac ta §6.3) - phai kiem RowsAffected() vi UPDATE khop 0 dong khong tu
// no la loi o tang pgx/PostgreSQL.
//
// Ham nay CHUA tu dat context.WithTimeout - dac ta §6.2 yeu cau moi lan GOI
// phai kem han gio (de xuat 5 giay), nhung do la trach nhiem cua noi GOI (vd
// vong lap luu dinh ky se lam o task sau), giong cach ba ham con lai trong
// file nay deu chi nhan ctx roi dung thang, khong tu bao boc.
//
// ⚠️ CO Y KHONG ghi vang/ngoc o day (bat bien BB4, dac ta §2.1 + duong nhan
// ban D7). Truoc kia ham nay ghi `vang = $13, ngoc = $14` bang GIA TRI TUYET
// DOI lay tu anh chup nv.Vang/nv.Ngoc. Vi duong luu dinh ky (ThoGhi ->
// ghiMotBan -> CapNhatNhanVat) chup anh roi ghi BAT DONG BO sau, neu nguoi
// choi tieu/kiem tien giua luc chup va luc xa hang doi, cu ghi tuyet doi se
// GHI DE tien ve gia tri cu = nhan ban/rollback tien. Tien tu nay CHI doi qua
// CongTruTien (cong/tru tai cho). CapNhatNhanVat khong bao gio cham vang/ngoc.
func (s *Store) CapNhatNhanVat(ctx context.Context, nv NhanVat) error {
	ct, err := s.pool.Exec(ctx,
		`UPDATE nhan_vat
		    SET lang = $1, kieu_toc = $2, map_id = $3, vi_tri_x = $4, vi_tri_y = $5,
		        cap = $6, kinh_nghiem = $7, suc_manh = $8, hp_toi_da = $9,
		        hp_hien_tai = $10, chakra_toi_da = $11, chakra_hien_tai = $12,
		        dang_nhap_cuoi = now()
		  WHERE id = $13`,
		nv.Lang, nv.KieuToc, nv.MapID, nv.ViTriX, nv.ViTriY,
		nv.Cap, nv.KinhNghiem, nv.SucManh, nv.HpToiDa,
		nv.HpHienTai, nv.ChakraToiDa, nv.ChakraHienTai,
		nv.ID,
	)
	if err != nil {
		return fmt.Errorf("cap nhat nhan vat that bai: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrKhongThayNhanVat
	}
	return nil
}

// CongTruTien cong/tru vang va ngoc cua mot nhan vat TAI CHO (dVang/dNgoc co
// the am de tru). Day la duong DUY NHAT hop le de doi tien (bat bien BB4, dac
// ta §2.1). Ket qua khong bao gio duoc am (chan duong nhan ban D10 - ban Java
// cu tru thang khong kiem so du -> tien am/in tien).
//
// Phan biet "khong du tien" voi "khong thay nhan vat": ca hai deu khien mot
// cau `UPDATE ... WHERE vang+$1>=0` khop 0 dong, khong the phan biet bang moi
// RowsAffected. Neu tach thanh hai cau (mot SELECT kiem ton tai + mot UPDATE)
// thi co khe TOCTOU: giua hai cau, phien khac co the xoa nhan vat hoac doi so
// du. Vi vay dung MOT cau duy nhat voi CTE:
//   - `ton_tai` khoa dong bang FOR UPDATE (chan phien khac chen vao giua),
//   - `cap_nhat` chi ghi khi so du sau van >= 0,
//   - cau SELECT cuoi tra ve HAI co: co nhan vat khong, va co ghi duoc khong.
//
// Tu do: khong co nhan vat -> ErrKhongThayNhanVat; co nhan vat nhung khong ghi
// -> ErrKhongDuTien; ghi duoc -> nil. Tat ca trong mot giao dich ngam, khong
// co khe ho giua cac buoc.
//
// Cong tren kieu numeric roi kep ve bigint de mot delta lon khong lam
// PostgreSQL nem loi tran (22003) huy ca cau; chan duoi van la >= 0.
func (s *Store) CongTruTien(ctx context.Context, nvID, dVang, dNgoc int64) error {
	return congTruTien(ctx, s.pool, nvID, dVang, dNgoc)
}

// CongTruTienTx la CHINH CongTruTien nhung chay trong mot giao dich DANG MO
// (KH7 Task 9). Sinh ra vi mua hang phai tru tien VA cap mon trong CUNG mot
// giao dich: neu tru tien dung ban `s.pool` thi no tu commit ngay, va buoc cap
// mon hong sau do se khong hoan tien lai duoc — dung dieu D2/§2.3 cam.
//
// Than cau lenh dung CHUNG mot ham `congTruTien` voi ban khong-tx: hai ban SQL
// khac nhau la mot ngay nao do mot ban duoc va con ban kia thi khong.
func (s *Store) CongTruTienTx(ctx context.Context, tx pgx.Tx, nvID, dVang, dNgoc int64) error {
	return congTruTien(ctx, tx, nvID, dVang, dNgoc)
}

// nguoiHoi la phan chung giua *pgxpool.Pool va pgx.Tx ma cac cau mot-dong o day
// can. Khai bao HEP dung mot phuong thuc — de khong ai nham tuong co the truyen
// mot thu gi khac vao.
type nguoiHoi interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func congTruTien(ctx context.Context, q nguoiHoi, nvID, dVang, dNgoc int64) error {
	var coNhanVat, daCapNhat bool
	err := q.QueryRow(ctx,
		`WITH ton_tai AS (
		     SELECT id, vang, ngoc FROM nhan_vat WHERE id = $1 FOR UPDATE
		 ),
		 cap_nhat AS (
		     UPDATE nhan_vat n
		        SET vang = (n.vang::numeric + $2::numeric)::bigint,
		            ngoc = (n.ngoc::numeric + $3::numeric)::bigint
		       FROM ton_tai t
		      WHERE n.id = t.id
		        AND t.vang::numeric + $2::numeric >= 0
		        AND t.ngoc::numeric + $3::numeric >= 0
		     RETURNING n.id
		 )
		 SELECT EXISTS(SELECT 1 FROM ton_tai), EXISTS(SELECT 1 FROM cap_nhat)`,
		nvID, dVang, dNgoc,
	).Scan(&coNhanVat, &daCapNhat)
	if err != nil {
		return fmt.Errorf("cong tru tien that bai: %w", err)
	}
	if !coNhanVat {
		return ErrKhongThayNhanVat
	}
	if !daCapNhat {
		return ErrKhongDuTien
	}
	return nil
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
