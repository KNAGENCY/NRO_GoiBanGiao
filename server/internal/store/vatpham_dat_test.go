package store

// vatpham_dat_test.go — KH7 Task 11: test duong DUOI DAT + NHAT tren CSDL THAT.
//
// Ba test duoi day la ba chot chan dot bien cua task (ghi ro o dau moi test):
//   - nhat HAI LAN cung mot id -> lan hai ErrDaCoNguoiNhat (giet dot bien "bo
//     `AND noi_chua='DUOI_DAT'` khoi cau UPDATE")
//   - lich su NHAT co DUNG MOT dong sau hai lan nhat (giet dot bien "ghi lich su
//     ngoai giao dich" / "ghi ca khi that bai")
//   - nhat vao o DA CO MON -> ErrOCoMon va mon VAN o duoi dat (giet dot bien
//     "COMMIT du UNIQUE hong")

import (
	"context"
	"errors"
	"testing"
)

// donDoDuoiDat xoa lich su roi xoa mon dang nam duoi dat (khong co chu nen
// donKhoCua khong don duoc).
func donDoDuoiDat(t *testing.T, s *Store, id int64) {
	t.Helper()
	t.Cleanup(func() {
		ctx := context.Background()
		_, _ = s.pool.Exec(ctx, `DELETE FROM vat_pham_lich_su WHERE vat_pham_id = $1`, id)
		_, _ = s.pool.Exec(ctx, `DELETE FROM vat_pham WHERE id = $1`, id)
	})
}

// noiChuaCua doc noi_chua hien tai cua mot mon.
func noiChuaCua(t *testing.T, s *Store, id int64) string {
	t.Helper()
	var noi string
	if err := s.pool.QueryRow(context.Background(),
		`SELECT noi_chua FROM vat_pham WHERE id = $1`, id).Scan(&noi); err != nil {
		t.Fatalf("doc noi_chua that bai: %v", err)
	}
	return noi
}

// ------------------------------------------------------------ SinhDoDuoiDat

// Duong vui: sinh mon duoi dat -> dong co chu NULL, noi DUOI_DAT, o NULL, va co
// DUNG MOT dong lich su SINH_RA (BB2).
func TestSinhDoDuoiDatGhiCaMonLanLichSu(t *testing.T) {
	s := moKhoVatPham(t)
	ctx := context.Background()

	id, err := s.SinhDoDuoiDat(ctx, 3, 5, khongOpt())
	if err != nil {
		t.Fatalf("sinh do duoi dat that bai: %v", err)
	}
	donDoDuoiDat(t, s, id)

	if noi := noiChuaCua(t, s, id); noi != NoiChuaDuoiDat {
		t.Fatalf("noi_chua = %q, mong doi %q", noi, NoiChuaDuoiDat)
	}
	if n := demLichSu(t, s, id); n != 1 {
		t.Fatalf("lich su phai co dung 1 dong SINH_RA, co %d", n)
	}
}

// ------------------------------------------------------------------ NhatLen

// DUONG VUI: nhat mot mon duoi dat -> mon vao tui dung o, va lich su co THEM
// dung mot dong (tong 2: SINH_RA + NHAT).
func TestNhatLenDuaMonVaoTuiVaGhiLichSu(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)
	ctx := context.Background()

	id, err := s.SinhDoDuoiDat(ctx, 3, 5, khongOpt())
	if err != nil {
		t.Fatalf("sinh do duoi dat that bai: %v", err)
	}
	donDoDuoiDat(t, s, id)

	if err := s.NhatLen(ctx, nv.ID, id, 7); err != nil {
		t.Fatalf("nhat len that bai: %v", err)
	}

	ds, err := s.LayTui(ctx, nv.ID, NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	if len(ds) != 1 || ds[0].ID != id || ds[0].OSo == nil || *ds[0].OSo != 7 {
		t.Fatalf("tui sau khi nhat sai: %+v", ds)
	}
	if n := demLichSu(t, s, id); n != 2 {
		t.Fatalf("lich su phai co 2 dong (SINH_RA + NHAT), co %d", n)
	}
}

// ⚠️ CHOT CHAN DOT BIEN CHINH CUA TASK (D8).
//
// Nhat HAI LAN cung mot id: lan hai PHAI ErrDaCoNguoiNhat, va lich su van chi
// co DUNG MOT dong NHAT. Bo `AND noi_chua = 'DUOI_DAT'` khoi cau UPDATE trong
// NhatLen -> lan hai thanh cong (mon "nhay" tu tui nguoi mot sang tui nguoi
// hai, va co HAI dong NHAT) -> test nay do.
func TestNhatLenHaiLanCungMonThiLanHaiThuaCuoc(t *testing.T) {
	s := moKhoVatPham(t)
	nvA := taoNhanVatThu(t, s, 0)
	nvB := taoNhanVatThu(t, s, 1)
	donKhoCua(t, s, nvA.ID)
	donKhoCua(t, s, nvB.ID)
	ctx := context.Background()

	id, err := s.SinhDoDuoiDat(ctx, 3, 1, khongOpt())
	if err != nil {
		t.Fatalf("sinh do duoi dat that bai: %v", err)
	}
	donDoDuoiDat(t, s, id)

	if err := s.NhatLen(ctx, nvA.ID, id, 0); err != nil {
		t.Fatalf("nguoi A nhat that bai: %v", err)
	}
	err = s.NhatLen(ctx, nvB.ID, id, 0)
	if !errors.Is(err, ErrDaCoNguoiNhat) {
		t.Fatalf("nguoi B phai nhan ErrDaCoNguoiNhat, nhan %v", err)
	}

	// Mon van thuoc A, khong nhay sang B.
	dsB, err := s.LayTui(ctx, nvB.ID, NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui B that bai: %v", err)
	}
	if len(dsB) != 0 {
		t.Fatalf("tui B phai rong, co %d mon", len(dsB))
	}
	// DUNG MOT dong NHAT — lan thua cuoc khong duoc de lai dau vet nao.
	var soNhat int
	if err := s.pool.QueryRow(ctx,
		`SELECT count(*) FROM vat_pham_lich_su WHERE vat_pham_id = $1 AND ly_do = $2`,
		id, LyDoNhat).Scan(&soNhat); err != nil {
		t.Fatalf("dem dong NHAT that bai: %v", err)
	}
	if soNhat != 1 {
		t.Fatalf("phai co dung 1 dong NHAT, co %d", soNhat)
	}
}

// Nhat vao o DA CO MON KHAC -> ErrOCoMon (UNIQUE vat_pham_mot_o_mot_mon), va
// mon VAN nam duoi dat (giao dich ROLLBACK sach). Giet dot bien "khong dich loi
// UNIQUE" / "COMMIT truoc khi INSERT lich su".
func TestNhatLenVaoODaCoMonThiMonVanNamDuoiDat(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)
	ctx := context.Background()

	if _, err := s.CapVatPham(ctx, nv.ID, NoiChuaTui, 2, 1, 1, khongOpt(), LyDoSinhRa); err != nil {
		t.Fatalf("cap mon chan o that bai: %v", err)
	}
	id, err := s.SinhDoDuoiDat(ctx, 3, 1, khongOpt())
	if err != nil {
		t.Fatalf("sinh do duoi dat that bai: %v", err)
	}
	donDoDuoiDat(t, s, id)

	if err := s.NhatLen(ctx, nv.ID, id, 2); !errors.Is(err, ErrOCoMon) {
		t.Fatalf("nhat vao o da co mon phai ErrOCoMon, nhan %v", err)
	}
	if noi := noiChuaCua(t, s, id); noi != NoiChuaDuoiDat {
		t.Fatalf("that bai xong mon phai VAN o DUOI_DAT, dang o %q", noi)
	}
	if n := demLichSu(t, s, id); n != 1 {
		t.Fatalf("that bai khong duoc ghi them lich su; co %d dong (mong doi 1)", n)
	}
}

// Nhat mot id KHONG TON TAI -> ErrDaCoNguoiNhat (0 dong bi cham), khong panic.
func TestNhatLenIdKhongTonTai(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)

	if err := s.NhatLen(context.Background(), nv.ID, 999_999_999, 0); !errors.Is(err, ErrDaCoNguoiNhat) {
		t.Fatalf("id khong ton tai phai ErrDaCoNguoiNhat, nhan %v", err)
	}
}

// --------------------------------------------------------------- ThaXuongDat

// Tha mot mon trong tui xuong dat roi nhat lai: vong tron khep kin, moi buoc
// deu ghi mot dong lich su.
func TestThaXuongDatRoiNhatLai(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)
	ctx := context.Background()

	id, err := s.CapVatPham(ctx, nv.ID, NoiChuaTui, 1, 5, 2, khongOpt(), LyDoSinhRa)
	if err != nil {
		t.Fatalf("cap vat pham that bai: %v", err)
	}
	donDoDuoiDat(t, s, id)

	idTha, err := s.ThaXuongDat(ctx, nv.ID, NoiChuaTui, 1, LyDoBo)
	if err != nil {
		t.Fatalf("tha xuong dat that bai: %v", err)
	}
	if idTha != id {
		t.Fatalf("id tha = %d, mong doi %d", idTha, id)
	}
	if noi := noiChuaCua(t, s, id); noi != NoiChuaDuoiDat {
		t.Fatalf("sau khi tha, noi_chua = %q", noi)
	}

	if err := s.NhatLen(ctx, nv.ID, id, 1); err != nil {
		t.Fatalf("nhat lai that bai: %v", err)
	}
	if n := demLichSu(t, s, id); n != 3 {
		t.Fatalf("lich su phai co 3 dong (SINH_RA + BO + NHAT), co %d", n)
	}
}

// Tha mon o mot o RONG -> ErrOTrong, khong ghi gi.
func TestThaXuongDatORongThiLoi(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)

	if _, err := s.ThaXuongDat(context.Background(), nv.ID, NoiChuaTui, 42, LyDoBo); !errors.Is(err, ErrOTrong) {
		t.Fatalf("tha o rong phai ErrOTrong, nhan %v", err)
	}
}

// Ly do ngoai danh sach -> ErrLyDoKhongHopLe, khong cham CSDL.
func TestThaXuongDatLyDoKhongHopLe(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)

	if _, err := s.ThaXuongDat(context.Background(), nv.ID, NoiChuaTui, 0, "BOC_HOI"); !errors.Is(err, ErrLyDoKhongHopLe) {
		t.Fatalf("ly do sai phai ErrLyDoKhongHopLe, nhan %v", err)
	}
}
