package store

// vatpham_kho.go — KH7 Task 6: kho vat pham trong CSDL.
//
// Bon ham: LayTui, CapVatPham, ChuyenO, HuyVatPham. TAT CA ghi thang bang giao
// dich CSDL (§5.3), KHONG qua ThoGhi — vi ThoGhi ghi de ca dong (D7).
//
// Ba dieu chot cung o day:
//   - D9/BB2: `CapVatPham` la duong DUY NHAT sinh mon. No ghi ca `vat_pham` lan
//     `vat_pham_lich_su` trong MOT giao dich — khong the co mon khong giay khai sinh.
//   - D4: moi UPDATE mon deu mang `WHERE chu_so_huu_id = $1 AND noi_chua = $2`.
//     Phien "ma" (dang nhap trung) doc du lieu cu se ghi 0 dong -> loi ro rang,
//     chu khong am tham nhan ban.
//   - D5: khoa `SELECT ... FOR UPDATE` theo `id` TANG DAN o moi noi. Khong thi hai
//     giao dich khoa nguoc thu tu -> deadlock.
//
// Quyet dinh cua ban 1 (ghi ro de nguoi sau khong doan):
//
//  1. `ChuyenO` KHONG hoan doi hai mon. O dich co mon -> tra `ErrOCoMon`, khong
//     lam gi. Hoan doi phai di qua o tam NULL trong cung giao dich (vi UNIQUE
//     `vat_pham_mot_o_mot_mon` chan buoc giua chung) — them duong do vao ban 1 la
//     them mot nhanh chua co nguoi goi. Tang tren (task 7, handler -40) cu bao
//     nguoi choi "o do da co do".
//  2. `ChuyenO` KHONG ghi `vat_pham_lich_su`. BB3 chi bat ghi khi DOI CHU SO HUU;
//     chuyen o la cung mot nguoi doi cho cat. Ghi them chi lam so nhat ky phinh
//     ma khong tra loi duoc cau hoi nao ("mon nay tu tay ai sang tay ai").

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrOCoMon: o dich da co mon khac (UNIQUE vat_pham_mot_o_mot_mon).
	ErrOCoMon = errors.New("o da co mon khac")
	// ErrOTrong: o nguon khong co mon nao cua nguoi nay (0 dong bi cham).
	ErrOTrong = errors.New("o khong co mon nao")
	// ErrVatPhamKhongHopLe: du lieu mon vi pham rang buoc CSDL (so_luong <= 0,
	// co option ma so_luong > 1, chu/noi khong khop, noi_chua ngoai danh sach).
	ErrVatPhamKhongHopLe = errors.New("du lieu vat pham khong hop le")
	// ErrLyDoKhongHopLe: ly do ghi so ngoai danh sach §5.1.
	ErrLyDoKhongHopLe = errors.New("ly do lich su khong hop le")
)

// Ly do doi chu — trung tung chu voi cot `ly_do` cua `vat_pham_lich_su` (§5.1).
const (
	LyDoSinhRa   = "SINH_RA"
	LyDoMuaNPC   = "MUA_NPC"
	LyDoBanNPC   = "BAN_NPC"
	LyDoNhat     = "NHAT"
	LyDoDung     = "DUNG"
	LyDoBo       = "BO"
	LyDoGiaoDich = "GIAO_DICH"
	LyDoAdmin    = "ADMIN"
	// LyDoHetGio: mon nam duoi dat qua lau khong ai nhat -> may chu tu don
	// (KH7, don do het gio). Tach RIENG khoi LyDoBo vi hai viec khac han nhau:
	// BO la nguoi choi tu tay vut, HET_GIO la may chu doc so. Nham hai cai nay
	// lam nhat ky doi chu noi doi (khieu nai "toi co vut dau" khong tra duoc).
	// Cot `ly_do` khong co CHECK o CSDL (migration 0005), danh sach nay la CHOT
	// DUY NHAT — them chu moi phai them CA O DAY va o lyDoHopLe ben duoi.
	LyDoHetGio = "HET_GIO"
)

func lyDoHopLe(lyDo string) bool {
	switch lyDo {
	case LyDoSinhRa, LyDoMuaNPC, LyDoBanNPC, LyDoNhat, LyDoDung, LyDoBo, LyDoGiaoDich, LyDoAdmin, LyDoHetGio:
		return true
	}
	return false
}

// dichLoiCSDL doi loi PostgreSQL thanh loi Go co nghia. KHONG nuot loi: cai gi
// khong nhan ra thi boc lai nguyen van.
func dichLoiCSDL(err error, moTa string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation — chi co mot UNIQUE tren bang nay
			return ErrOCoMon
		case "23514": // check_violation
			return ErrVatPhamKhongHopLe
		case "23503": // foreign_key_violation — chu_so_huu_id tro toi nhan vat khong co
			return ErrKhongThayNhanVat
		}
	}
	return fmt.Errorf("%s: %w", moTa, err)
}

// LayTui doc moi mon cua mot nhan vat o MOT noi chua (TUI/MAC/RUONG/KY_GUI),
// sap theo `o_so` tang dan — dung thu tu tang tren dung de dung goi -36/-37/-35.
//
// Khong tim thay nhan vat va tui rong deu tra danh sach rong: doc mot tui rong
// khong phai loi. Ai can biet nhan vat co ton tai hay khong thi hoi bang nhan_vat.
func (s *Store) LayTui(ctx context.Context, nvID int64, noiChua string) ([]VatPham, error) {
	hang, err := s.pool.Query(ctx,
		`SELECT id, mau_id, chu_so_huu_id, noi_chua, o_so, so_luong, chi_so_phu
		   FROM vat_pham
		  WHERE chu_so_huu_id = $1 AND noi_chua = $2
		  ORDER BY o_so`,
		nvID, noiChua)
	if err != nil {
		return nil, fmt.Errorf("lay tui that bai: %w", err)
	}
	defer hang.Close()

	var ds []VatPham
	for hang.Next() {
		var vp VatPham
		var thoChiSo []byte
		if err := hang.Scan(&vp.ID, &vp.MauID, &vp.ChuSoHuuID, &vp.NoiChua,
			&vp.OSo, &vp.SoLuong, &thoChiSo); err != nil {
			return nil, fmt.Errorf("doc dong vat pham that bai: %w", err)
		}
		if err := json.Unmarshal(thoChiSo, &vp.ChiSoPhu); err != nil {
			return nil, fmt.Errorf("giai ma chi_so_phu cua mon %d that bai: %w", vp.ID, err)
		}
		ds = append(ds, vp)
	}
	if err := hang.Err(); err != nil {
		return nil, fmt.Errorf("duyet ket qua lay tui that bai: %w", err)
	}
	return ds, nil
}

// CapVatPham sinh MOT mon moi vao o `oSo` cua nhan vat — duong DUY NHAT dua mon
// vao game (D9). Mot giao dich: INSERT `vat_pham` + INSERT `vat_pham_lich_su`
// (`tu_ai` NULL = tu he thong sinh ra, `sang_ai` = nguoi nhan). Hong mot trong
// hai buoc thi ca hai deu khong xay ra — khong bao gio co mon khong giay khai
// sinh (BB2/BB3).
//
// Loi phan biet duoc: o da co mon -> ErrOCoMon; du lieu sai rang buoc ->
// ErrVatPhamKhongHopLe; khong co nhan vat do -> ErrKhongThayNhanVat.
func (s *Store) CapVatPham(ctx context.Context, nvID int64, noiChua string, oSo int16,
	mauID int32, soLuong int32, chiSoPhu ChiSoPhu, lyDo string) (int64, error) {

	if !lyDoHopLe(lyDo) {
		return 0, ErrLyDoKhongHopLe
	}
	if chiSoPhu.Opt == nil {
		chiSoPhu.Opt = []TuyChon{}
	}
	thoChiSo, err := json.Marshal(chiSoPhu)
	if err != nil {
		return 0, fmt.Errorf("ma hoa chi_so_phu that bai: %w", err)
	}

	var idMoi int64
	err = s.trongGiaoDich(ctx, func(tx pgx.Tx) error {
		var errTx error
		idMoi, errTx = capVatPhamTrongTx(ctx, tx, nvID, noiChua, oSo, mauID, soLuong, thoChiSo, lyDo)
		return errTx
	})
	if err != nil {
		return 0, err
	}
	return idMoi, nil
}

// CapVatPhamTx la CHINH CapVatPham nhung chay trong mot giao dich DANG MO
// (KH7 Task 9) — xem ly do o CongTruTienTx (store.go). Cung than SQL, cung ham
// `capVatPhamTrongTx`, nen khong the lech nhau.
func (s *Store) CapVatPhamTx(ctx context.Context, tx pgx.Tx, nvID int64, noiChua string, oSo int16,
	mauID int32, soLuong int32, chiSoPhu ChiSoPhu, lyDo string) (int64, error) {

	if !lyDoHopLe(lyDo) {
		return 0, ErrLyDoKhongHopLe
	}
	if chiSoPhu.Opt == nil {
		chiSoPhu.Opt = []TuyChon{}
	}
	thoChiSo, err := json.Marshal(chiSoPhu)
	if err != nil {
		return 0, fmt.Errorf("ma hoa chi_so_phu that bai: %w", err)
	}
	return capVatPhamTrongTx(ctx, tx, nvID, noiChua, oSo, mauID, soLuong, thoChiSo, lyDo)
}

// capVatPhamTrongTx la than that su: MOT INSERT mon + MOT INSERT lich su, khong
// bao gio tach roi (BB2/BB3). Nhan `thoChiSo` da ma hoa san de hai duong goi
// (CapVatPham / CapVatPhamTx) khong ma hoa hai kieu khac nhau.
func capVatPhamTrongTx(ctx context.Context, tx pgx.Tx, nvID int64, noiChua string, oSo int16,
	mauID int32, soLuong int32, thoChiSo []byte, lyDo string) (int64, error) {

	var idMoi int64
	if err := tx.QueryRow(ctx,
		`INSERT INTO vat_pham (mau_id, chu_so_huu_id, noi_chua, o_so, so_luong, chi_so_phu)
		 VALUES ($1, $2, $3, $4, $5, $6::jsonb)
		 RETURNING id`,
		mauID, nvID, noiChua, oSo, soLuong, thoChiSo,
	).Scan(&idMoi); err != nil {
		return 0, dichLoiCSDL(err, "chen vat pham that bai")
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO vat_pham_lich_su (vat_pham_id, tu_ai, sang_ai, ly_do, so_luong)
		 VALUES ($1, NULL, $2, $3, $4)`,
		idMoi, nvID, lyDo, soLuong); err != nil {
		return 0, dichLoiCSDL(err, "ghi lich su sinh ra that bai")
	}
	return idMoi, nil
}

// ChuyenO doi cho mot mon cua CHINH nhan vat do: tu (tuNoi, tuO) sang (denNoi, denO).
//
// Mot giao dich. Khoa ca dong nguon lan dong dich (neu o dich co mon) trong MOT cau
// SELECT ... ORDER BY id FOR UPDATE — thu tu id TANG DAN o moi noi, chong deadlock (D5).
// Cau UPDATE mang du `chu_so_huu_id`, `noi_chua`, `o_so` (D4): nguoi khac goi chuyen o
// cua nguoi nay se cham 0 dong -> ErrOTrong, khong bao gio ghi nham.
//
// O dich da co mon -> ErrOCoMon, khong hoan doi (xem ghi chu dau file).
func (s *Store) ChuyenO(ctx context.Context, nvID int64, tuNoi string, tuO int16,
	denNoi string, denO int16) error {

	if tuNoi == denNoi && tuO == denO {
		return nil // chuyen ve chinh no: khong lam gi, khong loi
	}
	return s.trongGiaoDich(ctx, func(tx pgx.Tx) error {
		// Khoa ca hai o cung luc, theo id tang dan.
		hang, err := tx.Query(ctx,
			`SELECT id, noi_chua, o_so
			   FROM vat_pham
			  WHERE chu_so_huu_id = $1
			    AND ((noi_chua = $2 AND o_so = $3) OR (noi_chua = $4 AND o_so = $5))
			  ORDER BY id
			  FOR UPDATE`,
			nvID, tuNoi, tuO, denNoi, denO)
		if err != nil {
			return fmt.Errorf("khoa hai o that bai: %w", err)
		}
		coNguon, coDich := false, false
		for hang.Next() {
			var id int64
			var noi string
			var o int16
			if err := hang.Scan(&id, &noi, &o); err != nil {
				hang.Close()
				return fmt.Errorf("doc dong khoa that bai: %w", err)
			}
			if noi == tuNoi && o == tuO {
				coNguon = true
			}
			if noi == denNoi && o == denO {
				coDich = true
			}
		}
		hang.Close()
		if err := hang.Err(); err != nil {
			return fmt.Errorf("duyet dong khoa that bai: %w", err)
		}
		if !coNguon {
			return ErrOTrong
		}
		if coDich {
			return ErrOCoMon
		}

		the, err := tx.Exec(ctx,
			`UPDATE vat_pham
			    SET noi_chua = $4, o_so = $5, sua_luc = now()
			  WHERE chu_so_huu_id = $1 AND noi_chua = $2 AND o_so = $3`,
			nvID, tuNoi, tuO, denNoi, denO)
		if err != nil {
			return dichLoiCSDL(err, "chuyen o that bai")
		}
		if the.RowsAffected() == 0 {
			return ErrOTrong
		}
		return nil
	})
}

// HuyVatPham tieu huy mon dang o (noiChua, oSo) cua nhan vat: dat
// `noi_chua='DA_HUY'`, `chu_so_huu_id=NULL`, `o_so=NULL` — KHONG xoa dong, vi id
// khong bao gio duoc tai su dung (D15) va nhat ky doi chu phai con tro toi duoc.
//
// Cung giao dich ghi mot dong `vat_pham_lich_su` voi `sang_ai` NULL (= bi tieu huy),
// `tu_ai` = chu cu (BB3).
func (s *Store) HuyVatPham(ctx context.Context, nvID int64, noiChua string, oSo int16, lyDo string) error {
	if !lyDoHopLe(lyDo) {
		return ErrLyDoKhongHopLe
	}
	return s.trongGiaoDich(ctx, func(tx pgx.Tx) error {
		return huyVatPhamTrongTx(ctx, tx, nvID, noiChua, oSo, lyDo)
	})
}

// HuyVatPhamTx la CHINH HuyVatPham nhung chay trong mot giao dich DANG MO
// (KH7 Task 10) — xem ly do o CongTruTienTx (store.go): duong BAN phai huy mon
// VA cong tien trong CUNG mot giao dich, neu khong thi dut giua chung la mat do
// ma khong duoc tien (D2).
//
// Cung than `huyVatPhamTrongTx` voi ban khong-tx, nen hai duong khong the lech.
func (s *Store) HuyVatPhamTx(ctx context.Context, tx pgx.Tx, nvID int64, noiChua string,
	oSo int16, lyDo string) error {

	if !lyDoHopLe(lyDo) {
		return ErrLyDoKhongHopLe
	}
	return huyVatPhamTrongTx(ctx, tx, nvID, noiChua, oSo, lyDo)
}

// huyVatPhamTrongTx la than that su cua HuyVatPham/HuyVatPhamTx.
func huyVatPhamTrongTx(ctx context.Context, tx pgx.Tx, nvID int64, noiChua string,
	oSo int16, lyDo string) error {
	{
		var id int64
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
			return fmt.Errorf("khoa mon can huy that bai: %w", err)
		}

		the, err := tx.Exec(ctx,
			`UPDATE vat_pham
			    SET noi_chua = 'DA_HUY', chu_so_huu_id = NULL, o_so = NULL, sua_luc = now()
			  WHERE id = $1 AND chu_so_huu_id = $2 AND noi_chua = $3`,
			id, nvID, noiChua)
		if err != nil {
			return dichLoiCSDL(err, "huy vat pham that bai")
		}
		if the.RowsAffected() == 0 {
			return ErrOTrong
		}

		if _, err := tx.Exec(ctx,
			`INSERT INTO vat_pham_lich_su (vat_pham_id, tu_ai, sang_ai, ly_do, so_luong)
			 VALUES ($1, $2, NULL, $3, $4)`,
			id, nvID, lyDo, soLuong); err != nil {
			return dichLoiCSDL(err, "ghi lich su tieu huy that bai")
		}
		return nil
	}
}

// trongGiaoDich chay `viec` trong mot giao dich CSDL: hong o bat cu buoc nao thi
// ROLLBACK toan bo. Rollback goi voi context.Background() de con don duoc ca khi
// ctx cua nguoi goi da bi huy (mat mang giua chung).
func (s *Store) trongGiaoDich(ctx context.Context, viec func(tx pgx.Tx) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("mo giao dich that bai: %w", err)
	}
	if err := viec(tx); err != nil {
		_ = tx.Rollback(context.Background())
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		_ = tx.Rollback(context.Background())
		return fmt.Errorf("chot giao dich that bai: %w", err)
	}
	return nil
}
