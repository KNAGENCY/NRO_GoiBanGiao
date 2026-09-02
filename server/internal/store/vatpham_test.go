package store

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

// vatpham_test.go — KH7 Task 1: test RANG BUOC CSDL THAT cua bang vat_pham
// (migration 0005). Day khong phai test logic Go — day la test rang CSDL tu no
// tu choi du lieu hong, ke ca khi tang Go co bug (spec §5.2). Moi test o day
// tuong ung mot rang buoc chong nhan ban trong spec.

// Ma loi PostgreSQL (SQLSTATE) — chot cung de test bat DUNG loai loi, khong bat
// "co loi la duoc" (mot loi khac, vd sai ten cot, cung se lam test xanh gia).
const (
	maLoiTrungDuyNhat = "23505" // unique_violation
	maLoiViPhamCheck  = "23514" // check_violation
)

// moKhoVatPham mo ket noi thu VA ap moi migration con lai — vi Open() khong tu ap
// migration, con DB test co the chua co bang vat_pham (0005 moi them). Idempotent
// (IF NOT EXISTS) nen goi lai bao nhieu lan cung an toan.
func moKhoVatPham(t *testing.T) *Store {
	t.Helper()
	s := moKetNoiThu(t)
	if _, err := s.ApDungMigrationConLai(context.Background()); err != nil {
		t.Fatalf("ap migration that bai: %v", err)
	}
	return s
}

// donVatPhamCua xoa moi mon cua mot nhan vat sau test. Dang ky SAU taoNhanVatThu
// de chay TRUOC cleanup xoa nhan vat (t.Cleanup la LIFO) — khong thi khoa ngoai
// chan xoa nhan vat va de rac lai DB test.
func donVatPhamCua(t *testing.T, s *Store, nvID int64) {
	t.Helper()
	t.Cleanup(func() {
		_, _ = s.pool.Exec(context.Background(),
			`DELETE FROM vat_pham WHERE chu_so_huu_id = $1`, nvID)
	})
}

// chenMon chen mot dong vat_pham bang SQL tho — co y KHONG qua ham Go nao (chua co
// o task 1), de test do thang rang buoc CSDL.
func chenMon(s *Store, chu *int64, noiChua string, oSo *int16, soLuong int, chiSoPhu string) error {
	_, err := s.pool.Exec(context.Background(),
		`INSERT INTO vat_pham (mau_id, chu_so_huu_id, noi_chua, o_so, so_luong, chi_so_phu)
		 VALUES (1, $1, $2, $3, $4, $5::jsonb)`,
		chu, noiChua, oSo, soLuong, chiSoPhu)
	return err
}

// doiLoiSQL khang dinh err la loi PostgreSQL mang dung SQLSTATE mong doi.
func doiLoiSQL(t *testing.T, err error, ma string, moTa string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: mong doi CSDL TU CHOI, nhung chen thanh cong", moTa)
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		t.Fatalf("%s: loi khong phai PgError: %v", moTa, err)
	}
	if pgErr.Code != ma {
		t.Fatalf("%s: SQLSTATE = %s (%s), mong doi %s", moTa, pgErr.Code, pgErr.Message, ma)
	}
}

func i64(v int64) *int64 { return &v }
func i16(v int16) *int16 { return &v }

// BB1 — rang buoc CHONG NHAN BAN manh nhat: mot o cua mot nguoi chi chua DUNG MOT
// mon. Chen mon thu hai vao cung (chu, noi_chua, o_so) phai bi UNIQUE tu choi.
func TestVatPhamMotOMotMon(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donVatPhamCua(t, s, nv.ID)

	if err := chenMon(s, i64(nv.ID), NoiChuaTui, i16(3), 1, `{"opt":[]}`); err != nil {
		t.Fatalf("mon dau vao o 3 phai chen duoc: %v", err)
	}
	err := chenMon(s, i64(nv.ID), NoiChuaTui, i16(3), 1, `{"opt":[]}`)
	doiLoiSQL(t, err, maLoiTrungDuyNhat, "mon thu hai cung o 3 cung tui")

	// Cung o so nhung KHAC noi chua (tui vs ruong) thi la hai o khac nhau -> hop le.
	if err := chenMon(s, i64(nv.ID), NoiChuaRuong, i16(3), 1, `{"opt":[]}`); err != nil {
		t.Fatalf("o 3 cua RUONG khac o 3 cua TUI, phai chen duoc: %v", err)
	}
}

// Mon DUOI_DAT (khong chu, khong o) nhieu mon cung ton tai la binh thuong — khong
// rang buoc nao (UNIQUE lan CHECK chu_va_noi_khop) duoc chan chuyen tha mon thu hai
// xuong dat. LUU Y THAT: test nay KHONG bat duoc dot bien "bo WHERE khoi partial
// index", vi PostgreSQL coi cac NULL la KHAC NHAU trong UNIQUE (mac dinh, khong
// NULLS NOT DISTINCT) nen (NULL,'DUOI_DAT',NULL) x2 van chen duoc voi index thuong.
// WHERE tren index la de nho gon + noi ro y, khong phai lop chan. Dot bien BB1 that
// su phai giet la XOA HAN index — xem TestVatPhamMotOMotMon.
func TestVatPhamDuoiDatNhieuMonKhongBiChan(t *testing.T) {
	s := moKhoVatPham(t)
	t.Cleanup(func() {
		_, _ = s.pool.Exec(context.Background(),
			`DELETE FROM vat_pham WHERE chu_so_huu_id IS NULL AND mau_id = 1`)
	})
	for i := 0; i < 2; i++ {
		if err := chenMon(s, nil, NoiChuaDuoiDat, nil, 1, `{"opt":[]}`); err != nil {
			t.Fatalf("mon duoi dat thu %d phai chen duoc: %v", i+1, err)
		}
	}
}

// BB2 / §1.3 — mon CO option thi luon la mot cai: so_luong > 1 kem option phai bi
// CHECK tu choi. Day la ranh gioi cong don / don chiec, dat o CSDL de web/script
// khong lach duoc.
func TestVatPhamCoOptionThiDonChiec(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donVatPhamCua(t, s, nv.ID)

	err := chenMon(s, i64(nv.ID), NoiChuaTui, i16(1), 2, `{"opt":[{"id":47,"param":120}]}`)
	doiLoiSQL(t, err, maLoiViPhamCheck, "2 cai kem option")

	// Khong option thi cong don duoc.
	if err := chenMon(s, i64(nv.ID), NoiChuaTui, i16(1), 5, `{"opt":[]}`); err != nil {
		t.Fatalf("5 cai KHONG option phai chen duoc: %v", err)
	}
	// So luong 0 hoac am bi chan.
	err = chenMon(s, i64(nv.ID), NoiChuaTui, i16(2), 0, `{"opt":[]}`)
	doiLoiSQL(t, err, maLoiViPhamCheck, "so_luong = 0")
}

// Chu va noi chua phai khop: co chu thi phai o TUI/MAC/RUONG/KY_GUI; khong chu thi
// phai DUOI_DAT/DA_HUY. Mot mon "co chu nhung duoi dat" hay "khong chu nhung trong
// tui" la trang thai nua voi — dung thu de ra nhan ban khi dut giua chung.
func TestVatPhamChuVaNoiPhaiKhop(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donVatPhamCua(t, s, nv.ID)

	err := chenMon(s, i64(nv.ID), NoiChuaDuoiDat, nil, 1, `{"opt":[]}`)
	doiLoiSQL(t, err, maLoiViPhamCheck, "co chu nhung DUOI_DAT")

	err = chenMon(s, nil, NoiChuaTui, i16(1), 1, `{"opt":[]}`)
	doiLoiSQL(t, err, maLoiViPhamCheck, "khong chu nhung trong TUI")

	err = chenMon(s, i64(nv.ID), "TRONG_MO", i16(1), 1, `{"opt":[]}`)
	doiLoiSQL(t, err, maLoiViPhamCheck, "noi_chua ngoai danh sach")
}
