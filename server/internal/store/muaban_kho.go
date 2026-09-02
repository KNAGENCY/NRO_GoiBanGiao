package store

// muaban_kho.go — KH7 Task 9: MOT giao dich CSDL cho ca hanh dong mua.
//
// Vi sao ham nay ton tai (thay vi tang tren tu goi CongTruTien roi CapVatPham):
// hai ham do moi ham tu mo MOT giao dich rieng. Goi lien tiep la HAI giao dich
// — dung cai dac ta cam o D2/§2.3: dut giua chung thi tien da tru xong ma mon
// chua co, va khong ai hoan lai duoc. Bang do dot bien bat buoc cua KH7 (§6)
// ghi thang: "tach thanh 2 giao dich -> test phai do".
//
// Ba viec, MOT BEGIN...COMMIT, theo dung thu tu §2.3 "ghi so truoc — giao sau":
//
//  1. Ghi khoa idempotency (D3). Khoa da co -> day la GOI LAP: doc lai ket qua
//     lan dau, tra ve nguyen van, KHONG lam lai gi ca.
//  2. Tru tien tai cho, dieu kien `>= 0` nam TRONG chinh cau UPDATE (BB4/D10).
//     Khong du -> ErrKhongDuTien -> ca giao dich ROLLBACK.
//  3. Cap mon + ghi lich su (BB2/BB3), roi ghi ket qua vao dong idempotency.
//
// Hong o bat ky buoc nao thi KHONG buoc nao xay ra: khong mat tien, khong mat
// do, khong de lai khoa idempotency chet (dong khoa cung bi rollback nen lan
// sau nguoi choi bam lai van mua duoc).

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// ErrGiaoDichHongKetQua: dong idempotency co that nhung `ket_qua` khong doc
// duoc. Khong bao gio duoc xay ra; neu xay ra thi la du lieu hong, KHONG duoc
// am tham lam lai giao dich (lam lai la nhan ban).
var ErrGiaoDichHongKetQua = errors.New("ket qua giao dich cu khong doc duoc")

// LoaiHanhDongMua la phan giua cua khoa idempotency `<nhanVatId>:<loai>:<soThuTu>`.
const LoaiHanhDongMua = "MUA"

// KhoaIdemMua dung khoa idempotency cho mot lan mua. MOT cho duy nhat dung khoa
// nay de tang tren va tang duoi khong bao gio ghep khac nhau.
func KhoaIdemMua(nvID int64, soThuTu uint64) string {
	return fmt.Sprintf("%d:%s:%d", nvID, LoaiHanhDongMua, soThuTu)
}

// YeuCauMua la mot lan mua da duoc tang tren kiem xong moi thu: mon nao, o nao,
// tru bao nhieu tien. Tang `store` KHONG biet gi ve cua hang — no chi thi hanh.
type YeuCauMua struct {
	NhanVatID int64
	// KhoaIdem: dung KhoaIdemMua de dung. Rong -> loi (khong co khoa la mat D3).
	KhoaIdem string
	// DVang/DNgoc la DELTA, thuong <= 0 khi mua. Cho phep 0.
	DVang int64
	DNgoc int64
	// Cho mon moi.
	NoiChua  string
	OSo      int16
	MauID    int32
	SoLuong  int32
	ChiSoPhu ChiSoPhu
	LyDo     string
}

// KetQuaMua la thu duoc ghi vao `vat_pham_giao_dich.ket_qua` — chinh la thu tra
// lai nguyen van khi goi lap. Co tag JSON vi no PHAI qua lai CSDL duoc.
type KetQuaMua struct {
	VatPhamID int64 `json:"vat_pham_id"`
	OSo       int16 `json:"o_so"`
	MauID     int32 `json:"mau_id"`
	SoLuong   int32 `json:"so_luong"`
	DVang     int64 `json:"d_vang"`
	DNgoc     int64 `json:"d_ngoc"`

	// LaLapLai KHONG di vao CSDL (`-`): no noi ve LAN GOI NAY, khong phai ve
	// giao dich. Lan dau false, moi lan gui lai cung khoa deu true.
	LaLapLai bool `json:"-"`
}

// MuaMotGiaoDich thi hanh mot lan mua. Xem chu thich dau file ve thu tu ba buoc.
//
// Tra ErrKhongDuTien / ErrOCoMon / ErrKhongThayNhanVat / ErrVatPhamKhongHopLe
// giong het hai ham goc — tang tren dich sang cau tieng Viet cho nguoi choi.
func (s *Store) MuaMotGiaoDich(ctx context.Context, yc YeuCauMua) (KetQuaMua, error) {
	if yc.KhoaIdem == "" {
		return KetQuaMua{}, fmt.Errorf("thieu khoa idempotency")
	}
	if !lyDoHopLe(yc.LyDo) {
		return KetQuaMua{}, ErrLyDoKhongHopLe
	}
	if yc.DVang > 0 || yc.DNgoc > 0 {
		// Chan tu goc: mot lan "mua" khong bao gio CONG tien cho nguoi choi.
		// Neu mot ngay nao do can duong cong tien, do la ham khac.
		return KetQuaMua{}, fmt.Errorf("mua chi duoc tru tien (dVang=%d, dNgoc=%d)", yc.DVang, yc.DNgoc)
	}
	if yc.ChiSoPhu.Opt == nil {
		yc.ChiSoPhu.Opt = []TuyChon{}
	}
	thoChiSo, err := json.Marshal(yc.ChiSoPhu)
	if err != nil {
		return KetQuaMua{}, fmt.Errorf("ma hoa chi_so_phu that bai: %w", err)
	}

	var kq KetQuaMua
	err = s.trongGiaoDich(ctx, func(tx pgx.Tx) error {
		// BUOC 1 — khoa idempotency (D3).
		//
		// ON CONFLICT DO NOTHING chu KHONG phai "SELECT xem co chua roi INSERT":
		// hai buoc la mot khe TOCTOU, hai goi lap cung luc deu thay "chua co".
		// Voi ON CONFLICT, ban thu hai CHO ban dau commit roi nhan 0 dong —
		// dung mot buoc, khong khe ho nao.
		var idGiaoDich int64
		err := tx.QueryRow(ctx,
			`INSERT INTO vat_pham_giao_dich (khoa_idem, ket_qua)
			 VALUES ($1, '{}'::jsonb)
			 ON CONFLICT (khoa_idem) DO NOTHING
			 RETURNING id`,
			yc.KhoaIdem,
		).Scan(&idGiaoDich)
		if errors.Is(err, pgx.ErrNoRows) {
			// Khoa da ton tai -> GOI LAP. Tra lai DUNG ket qua lan dau.
			var tho []byte
			if err := tx.QueryRow(ctx,
				`SELECT ket_qua FROM vat_pham_giao_dich WHERE khoa_idem = $1`,
				yc.KhoaIdem,
			).Scan(&tho); err != nil {
				return fmt.Errorf("doc ket qua giao dich cu that bai: %w", err)
			}
			var cu KetQuaMua
			if err := json.Unmarshal(tho, &cu); err != nil {
				return ErrGiaoDichHongKetQua
			}
			if cu.VatPhamID == 0 {
				// Dong co that nhung than rong: khong the xay ra vi buoc 3 ghi
				// than trong CUNG giao dich. Neu gap, KHONG lam lai — bao loi.
				return ErrGiaoDichHongKetQua
			}
			cu.LaLapLai = true
			kq = cu
			return nil
		}
		if err != nil {
			return fmt.Errorf("ghi khoa idempotency that bai: %w", err)
		}

		// BUOC 2 — tru tien TAI CHO, dieu kien >= 0 nam trong chinh cau UPDATE.
		if yc.DVang != 0 || yc.DNgoc != 0 {
			if err := congTruTien(ctx, tx, yc.NhanVatID, yc.DVang, yc.DNgoc); err != nil {
				return err
			}
		}

		// BUOC 3 — cap mon (kem lich su) roi ghi ket qua vao dong idempotency.
		idMon, err := capVatPhamTrongTx(ctx, tx, yc.NhanVatID, yc.NoiChua, yc.OSo,
			yc.MauID, yc.SoLuong, thoChiSo, yc.LyDo)
		if err != nil {
			return err
		}
		kq = KetQuaMua{
			VatPhamID: idMon, OSo: yc.OSo, MauID: yc.MauID, SoLuong: yc.SoLuong,
			DVang: yc.DVang, DNgoc: yc.DNgoc,
		}
		thoKQ, err := json.Marshal(kq)
		if err != nil {
			return fmt.Errorf("ma hoa ket qua giao dich that bai: %w", err)
		}
		if _, err := tx.Exec(ctx,
			`UPDATE vat_pham_giao_dich SET ket_qua = $2::jsonb WHERE id = $1`,
			idGiaoDich, thoKQ); err != nil {
			return fmt.Errorf("ghi ket qua giao dich that bai: %w", err)
		}
		return nil
	})
	if err != nil {
		return KetQuaMua{}, err
	}
	return kq, nil
}

// ---------------------------------------------------------------------------
// KH7 Task 10 ? BAN cho NPC. Anh guong cua MuaMotGiaoDich o tren.
// ---------------------------------------------------------------------------

// LoaiHanhDongBan la phan giua cua khoa idempotency khi BAN.
const LoaiHanhDongBan = "BAN"

// KhoaIdemBan dung khoa idempotency cho mot lan ban. MOT cho duy nhat dung khoa
// nay ? giong het KhoaIdemMua, va CO Y khac chu o giua de mot lan mua va mot
// lan ban cung so thu tu khong bao gio dam vao nhau.
func KhoaIdemBan(nvID int64, soThuTu uint64) string {
	return fmt.Sprintf("%d:%s:%d", nvID, LoaiHanhDongBan, soThuTu)
}

// ErrKhongDuSoLuong: xin ban nhieu hon so dang co trong o do.
var ErrKhongDuSoLuong = errors.New("khong du so luong de ban")

// YeuCauBan la mot lan ban da duoc tang tren kiem xong moi thu.
type YeuCauBan struct {
	NhanVatID int64
	// KhoaIdem: dung KhoaIdemBan de dung. Rong -> loi.
	KhoaIdem string
	// Mon dem ban ? dinh vi bang CA BA: id (tra tu bang hanh trang), noi chua,
	// o so. Ca ba deu nam trong menh de WHERE (D4).
	VatPhamID int64
	NoiChua   string
	OSo       int16
	// SoLuong: so cai dem ban. Bang so dang co -> huy ca mon; it hon -> giam.
	SoLuong int32
	// DVang/DNgoc la DELTA, phai >= 0 khi ban.
	DVang int64
	DNgoc int64
	LyDo  string
}

// KetQuaBan la thu duoc ghi vao `vat_pham_giao_dich.ket_qua`.
type KetQuaBan struct {
	VatPhamID  int64 `json:"vat_pham_id"`
	SoLuongBan int32 `json:"so_luong_ban"`
	SoLuongCon int32 `json:"so_luong_con"`
	DaHuyMon   bool  `json:"da_huy_mon"`
	DVang      int64 `json:"d_vang"`
	DNgoc      int64 `json:"d_ngoc"`

	// LaLapLai KHONG di vao CSDL: no noi ve LAN GOI NAY.
	LaLapLai bool `json:"-"`
}

// BanMotGiaoDich thi hanh mot lan ban cho NPC: MOT BEGIN...COMMIT cho ca "bot
// do" lan "cong tien".
//
// Thu tu CO Y ? nguoc voi mua, va day la diem quan trong nhat cua ham:
//
//  1. Khoa idempotency (D3). Da co -> GOI LAP, tra lai ket qua lan dau, khong
//     lam gi them.
//  2. Khoa dong mon `FOR UPDATE` roi kiem: dung chu, dung noi, dung o, du so
//     luong. Guard chu so huu (D4) nam trong CHINH menh de WHERE ? phien "ma"
//     doc du lieu cu se cham 0 dong va nhan ErrOTrong, khong bao gio ban ho do
//     nguoi khac.
//  3. BOT DO TRUOC, cong tien SAU. Neu buoc cong tien hong thi ca giao dich
//     ROLLBACK -> do van con. Nguoc lai (cong tien truoc) cung an toan nho
//     rollback, nhung thu tu nay khop voi ky luat chung "giao thu cua minh
//     truoc, nhan ve sau".
//
// Hong o bat ky buoc nao: khong mat do, khong duoc tien, khong de lai khoa chet.
func (s *Store) BanMotGiaoDich(ctx context.Context, yc YeuCauBan) (KetQuaBan, error) {
	if yc.KhoaIdem == "" {
		return KetQuaBan{}, fmt.Errorf("thieu khoa idempotency")
	}
	if !lyDoHopLe(yc.LyDo) {
		return KetQuaBan{}, ErrLyDoKhongHopLe
	}
	if yc.SoLuong < 1 {
		return KetQuaBan{}, ErrVatPhamKhongHopLe
	}
	if yc.DVang < 0 || yc.DNgoc < 0 {
		// Chan tu goc: mot lan "ban" khong bao gio TRU tien cua nguoi choi.
		return KetQuaBan{}, fmt.Errorf("ban chi duoc cong tien (dVang=%d, dNgoc=%d)", yc.DVang, yc.DNgoc)
	}

	var kq KetQuaBan
	err := s.trongGiaoDich(ctx, func(tx pgx.Tx) error {
		// BUOC 1 ? khoa idempotency (D3). Xem giai thich ON CONFLICT o
		// MuaMotGiaoDich: mot buoc, khong khe TOCTOU nao.
		var idGiaoDich int64
		err := tx.QueryRow(ctx,
			`INSERT INTO vat_pham_giao_dich (khoa_idem, ket_qua)
			 VALUES ($1, '{}'::jsonb)
			 ON CONFLICT (khoa_idem) DO NOTHING
			 RETURNING id`,
			yc.KhoaIdem,
		).Scan(&idGiaoDich)
		if errors.Is(err, pgx.ErrNoRows) {
			var tho []byte
			if err := tx.QueryRow(ctx,
				`SELECT ket_qua FROM vat_pham_giao_dich WHERE khoa_idem = $1`,
				yc.KhoaIdem,
			).Scan(&tho); err != nil {
				return fmt.Errorf("doc ket qua giao dich cu that bai: %w", err)
			}
			var cu KetQuaBan
			if err := json.Unmarshal(tho, &cu); err != nil {
				return ErrGiaoDichHongKetQua
			}
			if cu.VatPhamID == 0 {
				return ErrGiaoDichHongKetQua
			}
			cu.LaLapLai = true
			kq = cu
			return nil
		}
		if err != nil {
			return fmt.Errorf("ghi khoa idempotency that bai: %w", err)
		}

		// BUOC 2 ? khoa dong mon va kiem. CA BA dieu kien (chu, noi, o) VA id
		// deu o trong WHERE: guard D4.
		var soDangCo int32
		err = tx.QueryRow(ctx,
			`SELECT so_luong FROM vat_pham
			  WHERE id = $1 AND chu_so_huu_id = $2 AND noi_chua = $3 AND o_so = $4
			  FOR UPDATE`,
			yc.VatPhamID, yc.NhanVatID, yc.NoiChua, yc.OSo).Scan(&soDangCo)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOTrong
		}
		if err != nil {
			return fmt.Errorf("khoa mon dem ban that bai: %w", err)
		}
		if yc.SoLuong > soDangCo {
			return ErrKhongDuSoLuong
		}

		// BUOC 3 ? bot do.
		daHuy := yc.SoLuong == soDangCo
		if daHuy {
			if err := huyVatPhamTrongTx(ctx, tx, yc.NhanVatID, yc.NoiChua, yc.OSo, yc.LyDo); err != nil {
				return err
			}
		} else {
			the, err := tx.Exec(ctx,
				`UPDATE vat_pham
				    SET so_luong = so_luong - $5, sua_luc = now()
				  WHERE id = $1 AND chu_so_huu_id = $2 AND noi_chua = $3 AND o_so = $4`,
				yc.VatPhamID, yc.NhanVatID, yc.NoiChua, yc.OSo, yc.SoLuong)
			if err != nil {
				return dichLoiCSDL(err, "bot so luong mon that bai")
			}
			if the.RowsAffected() == 0 {
				return ErrOTrong
			}
			// BB3: doi chu (tu nguoi choi sang he thong) thi PHAI co dong lich su.
			// `sang_ai` NULL = di ra khoi tay nguoi choi.
			if _, err := tx.Exec(ctx,
				`INSERT INTO vat_pham_lich_su (vat_pham_id, tu_ai, sang_ai, ly_do, so_luong)
				 VALUES ($1, $2, NULL, $3, $4)`,
				yc.VatPhamID, yc.NhanVatID, yc.LyDo, yc.SoLuong); err != nil {
				return dichLoiCSDL(err, "ghi lich su ban that bai")
			}
		}

		// BUOC 4 ? cong tien, cung giao dich.
		if yc.DVang != 0 || yc.DNgoc != 0 {
			if err := congTruTien(ctx, tx, yc.NhanVatID, yc.DVang, yc.DNgoc); err != nil {
				return err
			}
		}

		kq = KetQuaBan{
			VatPhamID: yc.VatPhamID, SoLuongBan: yc.SoLuong,
			SoLuongCon: soDangCo - yc.SoLuong, DaHuyMon: daHuy,
			DVang: yc.DVang, DNgoc: yc.DNgoc,
		}
		thoKQ, err := json.Marshal(kq)
		if err != nil {
			return fmt.Errorf("ma hoa ket qua giao dich that bai: %w", err)
		}
		if _, err := tx.Exec(ctx,
			`UPDATE vat_pham_giao_dich SET ket_qua = $2::jsonb WHERE id = $1`,
			idGiaoDich, thoKQ); err != nil {
			return fmt.Errorf("ghi ket qua giao dich that bai: %w", err)
		}
		return nil
	})
	if err != nil {
		return KetQuaBan{}, err
	}
	return kq, nil
}
