package store

// vatpham_dat.go — KH7 Task 11: mon nam DUOI DAT va duong NHAT LEN.
//
// Dac ta: docs/superpowers/specs/2026-09-03-kh7-vat-pham-design.md D8, §2.3, §5.1.
// Giai ma goi: docs/superpowers/research/2026-09-03-kh7-giai-ma-do-roi-nhat.md
//
// MOT CAU CHOT CA FILE (D8): nhat = MOT cau UPDATE co dieu kien
//
//	UPDATE vat_pham SET chu_so_huu_id=$1, noi_chua='TUI', o_so=$2
//	 WHERE id=$3 AND noi_chua='DUOI_DAT'
//
// Nguoi thu hai cham 0 dong -> ErrDaCoNguoiNhat. TUYET DOI KHONG kiem-roi-ghi
// hai buoc ("SELECT xem con duoi dat khong" roi moi "UPDATE"): giua hai buoc do
// la khe hoi ma hai nguoi bam cung luc deu lot qua — chinh la duong nhan ban D8.
//
// VI SAO KHONG CAN `SELECT ... FOR UPDATE` o day (khac ChuyenO/HuyVatPham):
// dieu kien `noi_chua='DUOI_DAT'` nam TRONG chinh cau UPDATE, va Postgres khoa
// dong khi UPDATE roi kiem lai dieu kien tren ban MOI NHAT (READ COMMITTED:
// EvalPlanQual). Hai UPDATE cung luc len cung dong -> cai thu hai cho cai dau
// commit roi kiem lai -> thay noi_chua da thanh 'TUI' -> KHONG cham dong nao.
// Them FOR UPDATE chi lam cham ma khong doi ket qua.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var (
	// ErrDaCoNguoiNhat: mon khong con o DUOI_DAT nua (nguoi khac nhat truoc,
	// hoac mon da het gio va bi don). Day la ket qua BINH THUONG cua mot cuoc
	// dua, KHONG phai loi he thong — tang tren bao nguoi choi "mon da co nguoi
	// nhat" roi thoi.
	ErrDaCoNguoiNhat = errors.New("mon da co nguoi nhat")
)

// SinhDoDuoiDat sinh MOT mon MOI thang xuong dat (khong qua tay ai) — duong
// dung khi quai/boss chet lam roi do (task 12 se goi).
//
// Mot giao dich: INSERT `vat_pham` voi `chu_so_huu_id = NULL`, `noi_chua =
// 'DUOI_DAT'`, `o_so = NULL` + INSERT `vat_pham_lich_su` ly do SINH_RA voi CA
// HAI dau `tu_ai`/`sang_ai` deu NULL.
//
// VI SAO CA HAI NULL: mon nay chua thuoc ve ai — no nam tren mat dat. Cot
// `sang_ai` NULL trong lich su thuong doc la "bi tieu huy", nhung cap
// (ly_do='SINH_RA', sang_ai=NULL) khong the nham voi tieu huy duoc, vi tieu huy
// luon mang ly_do khac (BO/DUNG). BB2 van duoc giu: mon co giay khai sinh.
func (s *Store) SinhDoDuoiDat(ctx context.Context, mauID int32, soLuong int32,
	chiSoPhu ChiSoPhu) (int64, error) {

	if chiSoPhu.Opt == nil {
		chiSoPhu.Opt = []TuyChon{}
	}
	thoChiSo, err := json.Marshal(chiSoPhu)
	if err != nil {
		return 0, fmt.Errorf("ma hoa chi_so_phu that bai: %w", err)
	}

	var idMoi int64
	err = s.trongGiaoDich(ctx, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx,
			`INSERT INTO vat_pham (mau_id, chu_so_huu_id, noi_chua, o_so, so_luong, chi_so_phu)
			 VALUES ($1, NULL, $2, NULL, $3, $4::jsonb)
			 RETURNING id`,
			mauID, NoiChuaDuoiDat, soLuong, thoChiSo,
		).Scan(&idMoi); err != nil {
			return dichLoiCSDL(err, "sinh do duoi dat that bai")
		}
		if _, err := tx.Exec(ctx,
			`INSERT INTO vat_pham_lich_su (vat_pham_id, tu_ai, sang_ai, ly_do, so_luong)
			 VALUES ($1, NULL, NULL, $2, $3)`,
			idMoi, LyDoSinhRa, soLuong); err != nil {
			return dichLoiCSDL(err, "ghi lich su sinh do duoi dat that bai")
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return idMoi, nil
}

// ThaXuongDat dua mot mon DANG CO CHU (o (noiChua, oSo) cua nhan vat nvID)
// xuong dat: `chu_so_huu_id = NULL`, `noi_chua = 'DUOI_DAT'`, `o_so = NULL`.
// Dung cho nguoi choi vut do (task 7) va cho duong HOAN TAC khong dung toi.
//
// Cung giao dich ghi mot dong lich su `tu_ai = nvID`, `sang_ai = NULL`, ly do do
// nguoi goi dua (thuong la LyDoBo) — BB3.
//
// Tra ve id cua mon vua tha, de tang tren cam vao danh sach do duoi dat cua khu vuc.
func (s *Store) ThaXuongDat(ctx context.Context, nvID int64, noiChua string, oSo int16,
	lyDo string) (int64, error) {

	if !lyDoHopLe(lyDo) {
		return 0, ErrLyDoKhongHopLe
	}
	var id int64
	err := s.trongGiaoDich(ctx, func(tx pgx.Tx) error {
		var soLuong int32
		err := tx.QueryRow(ctx,
			`SELECT id, so_luong FROM vat_pham
			  WHERE chu_so_huu_id = $1 AND noi_chua = $2 AND o_so = $3
			  ORDER BY id
			  FOR UPDATE`,
			nvID, noiChua, oSo).Scan(&id, &soLuong)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOTrong
		}
		if err != nil {
			return fmt.Errorf("khoa mon can tha that bai: %w", err)
		}

		the, err := tx.Exec(ctx,
			`UPDATE vat_pham
			    SET chu_so_huu_id = NULL, noi_chua = $3, o_so = NULL, sua_luc = now()
			  WHERE id = $1 AND chu_so_huu_id = $2`,
			id, nvID, NoiChuaDuoiDat)
		if err != nil {
			return dichLoiCSDL(err, "tha xuong dat that bai")
		}
		if the.RowsAffected() == 0 {
			return ErrOTrong
		}

		if _, err := tx.Exec(ctx,
			`INSERT INTO vat_pham_lich_su (vat_pham_id, tu_ai, sang_ai, ly_do, so_luong)
			 VALUES ($1, $2, NULL, $3, $4)`,
			id, nvID, lyDo, soLuong); err != nil {
			return dichLoiCSDL(err, "ghi lich su tha xuong dat that bai")
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

// HuyDoDuoiDatHetGio danh dau MOT mon dang nam DUOI_DAT la DA_HUY vi het gio
// khong ai nhat. Duong nay KHONG CO CHU: mon chua tung thuoc ve ai.
//
//	UPDATE vat_pham SET noi_chua='DA_HUY' WHERE id=$1 AND noi_chua='DUOI_DAT'
//
// Cung mot khuon MOT BUOC voi NhatLen (D8): dieu kien `noi_chua='DUOI_DAT'` nam
// TRONG chinh cau UPDATE, nen neu ai do vua nhat mon nay xong (noi_chua da la
// 'TUI') thi cau nay cham 0 dong va KHONG lam gi ca — khong the "huy nham do
// trong tui nguoi khac". Vi vay 0 dong la KET QUA BINH THUONG, tra ve
// ErrDaCoNguoiNhat de tang tren chi log.
//
// Cung giao dich ghi mot dong lich su `tu_ai = NULL`, `sang_ai = NULL`, ly do
// HET_GIO (BB3 — mon nao cung co giay khai tu). Cap (sang_ai NULL, ly_do
// HET_GIO) doc duoc la "may chu don rac", khong nham voi BO (nguoi choi vut).
func (s *Store) HuyDoDuoiDatHetGio(ctx context.Context, vatPhamID int64) error {
	return s.trongGiaoDich(ctx, func(tx pgx.Tx) error {
		var soLuong int32
		err := tx.QueryRow(ctx,
			`UPDATE vat_pham
			    SET noi_chua = $1, o_so = NULL, sua_luc = now()
			  WHERE id = $2 AND noi_chua = $3
			  RETURNING so_luong`,
			NoiChuaDaHuy, vatPhamID, NoiChuaDuoiDat).Scan(&soLuong)
		if errors.Is(err, pgx.ErrNoRows) {
			// 0 dong: mon khong con duoi dat (da bi nhat, hoac da huy roi).
			return ErrDaCoNguoiNhat
		}
		if err != nil {
			return dichLoiCSDL(err, "huy do duoi dat het gio that bai")
		}

		if _, err := tx.Exec(ctx,
			`INSERT INTO vat_pham_lich_su (vat_pham_id, tu_ai, sang_ai, ly_do, so_luong)
			 VALUES ($1, NULL, NULL, $2, $3)`,
			vatPhamID, LyDoHetGio, soLuong); err != nil {
			return dichLoiCSDL(err, "ghi lich su huy do het gio that bai")
		}
		return nil
	})
}

// NhatLen chuyen mot mon tu DUOI_DAT vao o `oSo` trong TUI cua nhan vat nvID.
//
// D8 — MOT BUOC, khong kiem-roi-ghi:
//
//	UPDATE ... SET chu_so_huu_id=$1, noi_chua='TUI', o_so=$2
//	 WHERE id=$3 AND noi_chua='DUOI_DAT'
//
// `RowsAffected() == 0` -> ErrDaCoNguoiNhat (nguoi thu hai thua cuoc). O tui da
// co mon khac -> ErrOCoMon (UNIQUE `vat_pham_mot_o_mot_mon` chan, dich o
// dichLoiCSDL). Ca hai truong hop deu ROLLBACK: khong co dong lich su nao
// duoc ghi cho mot lan nhat that bai.
//
// Cung giao dich ghi mot dong lich su `tu_ai = NULL` (tu mat dat), `sang_ai =
// nvID`, ly do NHAT — BB3. Dung MOT dong, khong hai: test dem so dong NHAT sau
// hai lan nhat cung mot id phai ra dung 1.
func (s *Store) NhatLen(ctx context.Context, nvID int64, vatPhamID int64, oSo int16) error {
	return s.trongGiaoDich(ctx, func(tx pgx.Tx) error {
		var soLuong int32
		err := tx.QueryRow(ctx,
			`UPDATE vat_pham
			    SET chu_so_huu_id = $1, noi_chua = $2, o_so = $3, sua_luc = now()
			  WHERE id = $4 AND noi_chua = $5
			  RETURNING so_luong`,
			nvID, NoiChuaTui, oSo, vatPhamID, NoiChuaDuoiDat).Scan(&soLuong)
		if errors.Is(err, pgx.ErrNoRows) {
			// 0 dong bi cham: mon khong con duoi dat nua.
			return ErrDaCoNguoiNhat
		}
		if err != nil {
			return dichLoiCSDL(err, "nhat len that bai")
		}

		if _, err := tx.Exec(ctx,
			`INSERT INTO vat_pham_lich_su (vat_pham_id, tu_ai, sang_ai, ly_do, so_luong)
			 VALUES ($1, NULL, $2, $3, $4)`,
			vatPhamID, nvID, LyDoNhat, soLuong); err != nil {
			return dichLoiCSDL(err, "ghi lich su nhat that bai")
		}
		return nil
	})
}
