package store

// vatpham_kho_test.go — KH7 Task 6: test bon ham kho vat pham tren CSDL THAT.
//
// Moi ham: duong vui + tung nhanh loi. Ba test duoi day la ba "chot chan dot bien"
// bat buoc cua task (ghi ro o dau moi test): lich su SINH_RA, WHERE chu_so_huu_id
// trong ChuyenO, va "huy KHONG xoa dong".

import (
	"context"
	"errors"
	"testing"
)

// demLichSu dem so dong nhat ky cua mot mon.
func demLichSu(t *testing.T, s *Store, vatPhamID int64) int {
	t.Helper()
	var n int
	if err := s.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM vat_pham_lich_su WHERE vat_pham_id = $1`, vatPhamID).Scan(&n); err != nil {
		t.Fatalf("dem lich su that bai: %v", err)
	}
	return n
}

// donKhoCua xoa lich su TRUOC roi moi xoa mon (khoa ngoai), cho ca mon da huy —
// mon da huy khong con chu_so_huu_id nen donVatPhamCua khong don duoc.
func donKhoCua(t *testing.T, s *Store, nvID int64) {
	t.Helper()
	t.Cleanup(func() {
		ctx := context.Background()
		_, _ = s.pool.Exec(ctx,
			`DELETE FROM vat_pham_lich_su
			  WHERE tu_ai = $1 OR sang_ai = $1
			     OR vat_pham_id IN (SELECT id FROM vat_pham WHERE chu_so_huu_id = $1)`, nvID)
		_, _ = s.pool.Exec(ctx, `DELETE FROM vat_pham WHERE chu_so_huu_id = $1`, nvID)
	})
}

func khongOpt() ChiSoPhu { return ChiSoPhu{Opt: []TuyChon{}} }

// ---------------------------------------------------------------- CapVatPham

// Duong vui: cap mot mon xong thi LayTui thay dung, VA lich su co DUNG MOT dong
// SINH_RA. Day la test giet dot bien "bo INSERT lich su trong CapVatPham" (BB2/BB3).
func TestCapVatPhamGhiCaMonLanLichSu(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)
	ctx := context.Background()

	id, err := s.CapVatPham(ctx, nv.ID, NoiChuaTui, 4, 77, 3, khongOpt(), LyDoSinhRa)
	if err != nil {
		t.Fatalf("cap vat pham that bai: %v", err)
	}
	if id <= 0 {
		t.Fatalf("id moi phai duong, duoc %d", id)
	}

	ds, err := s.LayTui(ctx, nv.ID, NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	if len(ds) != 1 {
		t.Fatalf("tui phai co 1 mon, co %d", len(ds))
	}
	vp := ds[0]
	if vp.ID != id || vp.MauID != 77 || vp.SoLuong != 3 || vp.NoiChua != NoiChuaTui {
		t.Fatalf("mon doc ra sai: %+v", vp)
	}
	if vp.OSo == nil || *vp.OSo != 4 {
		t.Fatalf("o_so sai: %v", vp.OSo)
	}
	if vp.ChuSoHuuID == nil || *vp.ChuSoHuuID != nv.ID {
		t.Fatalf("chu so huu sai: %v", vp.ChuSoHuuID)
	}

	// CHOT CHAN DOT BIEN 1: bo INSERT vat_pham_lich_su -> so nay thanh 0 -> do.
	if n := demLichSu(t, s, id); n != 1 {
		t.Fatalf("phai co DUNG 1 dong lich su SINH_RA, dem duoc %d", n)
	}
	var tuAi *int64
	var sangAi *int64
	var lyDo string
	if err := s.pool.QueryRow(ctx,
		`SELECT tu_ai, sang_ai, ly_do FROM vat_pham_lich_su WHERE vat_pham_id = $1`, id,
	).Scan(&tuAi, &sangAi, &lyDo); err != nil {
		t.Fatalf("doc lich su that bai: %v", err)
	}
	if tuAi != nil {
		t.Fatalf("tu_ai phai NULL (he thong sinh ra), duoc %d", *tuAi)
	}
	if sangAi == nil || *sangAi != nv.ID {
		t.Fatalf("sang_ai phai la nguoi nhan, duoc %v", sangAi)
	}
	if lyDo != LyDoSinhRa {
		t.Fatalf("ly_do = %q, mong doi %q", lyDo, LyDoSinhRa)
	}
}

// Cap vao o da co mon -> ErrOCoMon (khong nuot loi UNIQUE), va KHONG de lai rac:
// giao dich hong thi khong co dong lich su thua nao.
func TestCapVatPhamVaoODaCoMon(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)
	ctx := context.Background()

	if _, err := s.CapVatPham(ctx, nv.ID, NoiChuaTui, 1, 10, 1, khongOpt(), LyDoSinhRa); err != nil {
		t.Fatalf("mon dau phai cap duoc: %v", err)
	}
	_, err := s.CapVatPham(ctx, nv.ID, NoiChuaTui, 1, 11, 1, khongOpt(), LyDoSinhRa)
	if !errors.Is(err, ErrOCoMon) {
		t.Fatalf("mong doi ErrOCoMon, duoc %v", err)
	}
	ds, err := s.LayTui(ctx, nv.ID, NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	if len(ds) != 1 || ds[0].MauID != 10 {
		t.Fatalf("tui phai con dung mon dau, duoc %+v", ds)
	}
}

// Vi pham CHECK (so_luong <= 0, hoac co option ma nhieu cai) -> ErrVatPhamKhongHopLe.
func TestCapVatPhamViPhamRangBuoc(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)
	ctx := context.Background()

	_, err := s.CapVatPham(ctx, nv.ID, NoiChuaTui, 1, 10, 0, khongOpt(), LyDoSinhRa)
	if !errors.Is(err, ErrVatPhamKhongHopLe) {
		t.Fatalf("so_luong 0: mong doi ErrVatPhamKhongHopLe, duoc %v", err)
	}
	coOpt := ChiSoPhu{Opt: []TuyChon{{Id: 47, Param: 120}}}
	_, err = s.CapVatPham(ctx, nv.ID, NoiChuaTui, 1, 10, 2, coOpt, LyDoSinhRa)
	if !errors.Is(err, ErrVatPhamKhongHopLe) {
		t.Fatalf("2 cai kem option: mong doi ErrVatPhamKhongHopLe, duoc %v", err)
	}
	// Noi chua khong hop le (DA_HUY doi hoi chu NULL) cung phai la loi ro.
	_, err = s.CapVatPham(ctx, nv.ID, NoiChuaDaHuy, 1, 10, 1, khongOpt(), LyDoSinhRa)
	if !errors.Is(err, ErrVatPhamKhongHopLe) {
		t.Fatalf("cap thang vao DA_HUY: mong doi ErrVatPhamKhongHopLe, duoc %v", err)
	}
	// Mot cai KEM option thi hop le (§1.3) — va option doc lai duoc.
	id, err := s.CapVatPham(ctx, nv.ID, NoiChuaMac, 2, 10, 1, coOpt, LyDoSinhRa)
	if err != nil {
		t.Fatalf("1 cai kem option phai cap duoc: %v", err)
	}
	ds, err := s.LayTui(ctx, nv.ID, NoiChuaMac)
	if err != nil || len(ds) != 1 || ds[0].ID != id {
		t.Fatalf("doc lai mon co option that bai: %v / %+v", err, ds)
	}
	if len(ds[0].ChiSoPhu.Opt) != 1 || ds[0].ChiSoPhu.Opt[0].Id != 47 || ds[0].ChiSoPhu.Opt[0].Param != 120 {
		t.Fatalf("chi_so_phu giai ma sai: %+v", ds[0].ChiSoPhu)
	}
}

// Nhan vat khong ton tai -> ErrKhongThayNhanVat (khoa ngoai), khong phai loi tho.
func TestCapVatPhamNhanVatKhongCo(t *testing.T) {
	s := moKhoVatPham(t)
	_, err := s.CapVatPham(context.Background(), -12345, NoiChuaTui, 1, 10, 1, khongOpt(), LyDoSinhRa)
	if !errors.Is(err, ErrKhongThayNhanVat) {
		t.Fatalf("mong doi ErrKhongThayNhanVat, duoc %v", err)
	}
}

// Ly do ngoai danh sach bi chan o tang Go (CSDL chua co CHECK cho cot nay).
func TestCapVatPhamLyDoLa(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)
	_, err := s.CapVatPham(context.Background(), nv.ID, NoiChuaTui, 1, 10, 1, khongOpt(), "TU_TREN_TROI")
	if !errors.Is(err, ErrLyDoKhongHopLe) {
		t.Fatalf("mong doi ErrLyDoKhongHopLe, duoc %v", err)
	}
}

// ------------------------------------------------------------------- LayTui

// LayTui sap theo o_so tang dan, va KHONG lan sang noi chua khac / nguoi khac.
func TestLayTuiSapTheoOVaKhongLanNoiKhac(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	khac := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)
	donKhoCua(t, s, khac.ID)
	ctx := context.Background()

	for _, o := range []int16{5, 1, 3} {
		if _, err := s.CapVatPham(ctx, nv.ID, NoiChuaTui, o, 10, 1, khongOpt(), LyDoSinhRa); err != nil {
			t.Fatalf("cap o %d that bai: %v", o, err)
		}
	}
	if _, err := s.CapVatPham(ctx, nv.ID, NoiChuaRuong, 0, 10, 1, khongOpt(), LyDoSinhRa); err != nil {
		t.Fatalf("cap vao ruong that bai: %v", err)
	}
	if _, err := s.CapVatPham(ctx, khac.ID, NoiChuaTui, 1, 10, 1, khongOpt(), LyDoSinhRa); err != nil {
		t.Fatalf("cap cho nguoi khac that bai: %v", err)
	}

	ds, err := s.LayTui(ctx, nv.ID, NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	if len(ds) != 3 {
		t.Fatalf("tui phai co 3 mon, co %d", len(ds))
	}
	for i, mongDoi := range []int16{1, 3, 5} {
		if ds[i].OSo == nil || *ds[i].OSo != mongDoi {
			t.Fatalf("mon thu %d o o %v, mong doi %d", i, ds[i].OSo, mongDoi)
		}
	}
}

// Tui rong khong phai loi.
func TestLayTuiRong(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)
	ds, err := s.LayTui(context.Background(), nv.ID, NoiChuaTui)
	if err != nil {
		t.Fatalf("tui rong khong duoc bao loi: %v", err)
	}
	if len(ds) != 0 {
		t.Fatalf("tui rong phai tra 0 mon, duoc %d", len(ds))
	}
}

// ------------------------------------------------------------------ ChuyenO

// Duong vui: chuyen sang o trong (ca sang noi chua khac).
func TestChuyenOSangOTrong(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)
	ctx := context.Background()

	id, err := s.CapVatPham(ctx, nv.ID, NoiChuaTui, 2, 99, 1, khongOpt(), LyDoSinhRa)
	if err != nil {
		t.Fatalf("cap that bai: %v", err)
	}
	if err := s.ChuyenO(ctx, nv.ID, NoiChuaTui, 2, NoiChuaTui, 7); err != nil {
		t.Fatalf("chuyen trong cung tui that bai: %v", err)
	}
	if err := s.ChuyenO(ctx, nv.ID, NoiChuaTui, 7, NoiChuaRuong, 0); err != nil {
		t.Fatalf("chuyen tui -> ruong that bai: %v", err)
	}
	if ds, _ := s.LayTui(ctx, nv.ID, NoiChuaTui); len(ds) != 0 {
		t.Fatalf("tui phai rong sau khi chuyen, con %d mon", len(ds))
	}
	ds, err := s.LayTui(ctx, nv.ID, NoiChuaRuong)
	if err != nil || len(ds) != 1 || ds[0].ID != id {
		t.Fatalf("mon phai nam trong ruong: %v / %+v", err, ds)
	}
	if ds[0].OSo == nil || *ds[0].OSo != 0 {
		t.Fatalf("o_so sau chuyen sai: %v", ds[0].OSo)
	}
	// Chuyen o KHONG doi chu -> khong ghi them lich su (quyet dinh ban 1, BB3).
	if n := demLichSu(t, s, id); n != 1 {
		t.Fatalf("chuyen o khong duoc ghi them lich su; dem duoc %d dong", n)
	}
}

// O dich da co mon -> ErrOCoMon, va KHONG mon nao bi xe dich.
func TestChuyenOSangODaCoMon(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)
	ctx := context.Background()

	idA, _ := s.CapVatPham(ctx, nv.ID, NoiChuaTui, 1, 10, 1, khongOpt(), LyDoSinhRa)
	idB, _ := s.CapVatPham(ctx, nv.ID, NoiChuaTui, 2, 20, 1, khongOpt(), LyDoSinhRa)

	if err := s.ChuyenO(ctx, nv.ID, NoiChuaTui, 1, NoiChuaTui, 2); !errors.Is(err, ErrOCoMon) {
		t.Fatalf("mong doi ErrOCoMon, duoc %v", err)
	}
	ds, err := s.LayTui(ctx, nv.ID, NoiChuaTui)
	if err != nil || len(ds) != 2 {
		t.Fatalf("van phai con dung 2 mon: %v / %+v", err, ds)
	}
	if ds[0].ID != idA || ds[1].ID != idB {
		t.Fatalf("hai mon phai giu nguyen cho: %+v", ds)
	}
}

// O nguon trong -> ErrOTrong.
func TestChuyenOTuOTrong(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)
	err := s.ChuyenO(context.Background(), nv.ID, NoiChuaTui, 9, NoiChuaTui, 10)
	if !errors.Is(err, ErrOTrong) {
		t.Fatalf("mong doi ErrOTrong, duoc %v", err)
	}
}

// D4 — CHOT CHAN DOT BIEN 2: nhan vat KHAC goi chuyen o cua nguoi nay phai cham
// 0 dong -> ErrOTrong, va mon cua nguoi kia KHONG duoc xe dich. Bo
// `AND chu_so_huu_id = $1` khoi UPDATE (hoac khoi cau khoa) -> test nay do.
func TestChuyenOKhongDungToiDoNguoiKhac(t *testing.T) {
	s := moKhoVatPham(t)
	chu := taoNhanVatThu(t, s, 0)
	ke := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, chu.ID)
	donKhoCua(t, s, ke.ID)
	ctx := context.Background()

	id, err := s.CapVatPham(ctx, chu.ID, NoiChuaTui, 3, 55, 1, khongOpt(), LyDoSinhRa)
	if err != nil {
		t.Fatalf("cap cho chu that bai: %v", err)
	}
	// `ke` co tui RONG. Goi chuyen o 3 -> o 4 nhan danh `ke`.
	if err := s.ChuyenO(ctx, ke.ID, NoiChuaTui, 3, NoiChuaTui, 4); !errors.Is(err, ErrOTrong) {
		t.Fatalf("nguoi khac chuyen o cua chu: mong doi ErrOTrong, duoc %v", err)
	}
	ds, err := s.LayTui(ctx, chu.ID, NoiChuaTui)
	if err != nil || len(ds) != 1 || ds[0].ID != id {
		t.Fatalf("mon cua chu phai con nguyen: %v / %+v", err, ds)
	}
	if ds[0].OSo == nil || *ds[0].OSo != 3 {
		t.Fatalf("mon cua chu bi xe dich sang o %v", ds[0].OSo)
	}
	if ds2, _ := s.LayTui(ctx, ke.ID, NoiChuaTui); len(ds2) != 0 {
		t.Fatalf("tui ke phai van rong, co %d mon", len(ds2))
	}
}

// ---------------------------------------------------------------- HuyVatPham

// CHOT CHAN DOT BIEN 3: huy KHONG xoa dong. Dong con day du (id giu nguyen),
// `noi_chua='DA_HUY'`, chu/o thanh NULL, va co dong lich su `sang_ai` NULL.
// Doi UPDATE thanh DELETE -> test nay do (dong bien mat + khoa ngoai lich su gay).
func TestHuyVatPhamGiuDongVaGhiLichSu(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)
	ctx := context.Background()

	id, err := s.CapVatPham(ctx, nv.ID, NoiChuaTui, 6, 88, 2, khongOpt(), LyDoSinhRa)
	if err != nil {
		t.Fatalf("cap that bai: %v", err)
	}
	if err := s.HuyVatPham(ctx, nv.ID, NoiChuaTui, 6, LyDoBo); err != nil {
		t.Fatalf("huy that bai: %v", err)
	}

	var noiChua string
	var chu *int64
	var oSo *int16
	if err := s.pool.QueryRow(ctx,
		`SELECT noi_chua, chu_so_huu_id, o_so FROM vat_pham WHERE id = $1`, id,
	).Scan(&noiChua, &chu, &oSo); err != nil {
		t.Fatalf("dong phai CON (huy khong duoc xoa dong): %v", err)
	}
	if noiChua != NoiChuaDaHuy || chu != nil || oSo != nil {
		t.Fatalf("trang thai sau huy sai: noi=%s chu=%v o=%v", noiChua, chu, oSo)
	}
	if ds, _ := s.LayTui(ctx, nv.ID, NoiChuaTui); len(ds) != 0 {
		t.Fatalf("tui phai rong sau khi huy, con %d mon", len(ds))
	}

	// Lich su: 1 dong SINH_RA + 1 dong tieu huy.
	if n := demLichSu(t, s, id); n != 2 {
		t.Fatalf("phai co 2 dong lich su, dem duoc %d", n)
	}
	var tuAi, sangAi *int64
	var lyDo string
	var soLuong int32
	if err := s.pool.QueryRow(ctx,
		`SELECT tu_ai, sang_ai, ly_do, so_luong FROM vat_pham_lich_su
		  WHERE vat_pham_id = $1 ORDER BY id DESC LIMIT 1`, id,
	).Scan(&tuAi, &sangAi, &lyDo, &soLuong); err != nil {
		t.Fatalf("doc lich su tieu huy that bai: %v", err)
	}
	if sangAi != nil {
		t.Fatalf("sang_ai phai NULL (bi tieu huy), duoc %d", *sangAi)
	}
	if tuAi == nil || *tuAi != nv.ID {
		t.Fatalf("tu_ai phai la chu cu, duoc %v", tuAi)
	}
	if lyDo != LyDoBo || soLuong != 2 {
		t.Fatalf("lich su tieu huy sai: ly_do=%s so_luong=%d", lyDo, soLuong)
	}

	// Id KHONG duoc tai su dung: cap mon moi vao dung o vua trong phai ra id khac.
	idMoi, err := s.CapVatPham(ctx, nv.ID, NoiChuaTui, 6, 88, 1, khongOpt(), LyDoSinhRa)
	if err != nil {
		t.Fatalf("cap lai vao o vua trong that bai: %v", err)
	}
	if idMoi == id {
		t.Fatalf("id %d bi tai su dung — pha vo D15", id)
	}
}

// Huy o trong -> ErrOTrong. Huy hai lan cung vay (mon da roi khoi o).
func TestHuyVatPhamOTrong(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)
	ctx := context.Background()

	if err := s.HuyVatPham(ctx, nv.ID, NoiChuaTui, 0, LyDoBo); !errors.Is(err, ErrOTrong) {
		t.Fatalf("huy o trong: mong doi ErrOTrong, duoc %v", err)
	}
	if _, err := s.CapVatPham(ctx, nv.ID, NoiChuaTui, 0, 10, 1, khongOpt(), LyDoSinhRa); err != nil {
		t.Fatalf("cap that bai: %v", err)
	}
	if err := s.HuyVatPham(ctx, nv.ID, NoiChuaTui, 0, LyDoBo); err != nil {
		t.Fatalf("huy lan dau that bai: %v", err)
	}
	if err := s.HuyVatPham(ctx, nv.ID, NoiChuaTui, 0, LyDoBo); !errors.Is(err, ErrOTrong) {
		t.Fatalf("huy lan hai: mong doi ErrOTrong, duoc %v", err)
	}
}

// D4 cho HuyVatPham: nguoi khac khong huy duoc do cua minh.
func TestHuyVatPhamKhongDungToiDoNguoiKhac(t *testing.T) {
	s := moKhoVatPham(t)
	chu := taoNhanVatThu(t, s, 0)
	ke := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, chu.ID)
	donKhoCua(t, s, ke.ID)
	ctx := context.Background()

	id, err := s.CapVatPham(ctx, chu.ID, NoiChuaTui, 1, 10, 1, khongOpt(), LyDoSinhRa)
	if err != nil {
		t.Fatalf("cap that bai: %v", err)
	}
	if err := s.HuyVatPham(ctx, ke.ID, NoiChuaTui, 1, LyDoBo); !errors.Is(err, ErrOTrong) {
		t.Fatalf("mong doi ErrOTrong, duoc %v", err)
	}
	ds, _ := s.LayTui(ctx, chu.ID, NoiChuaTui)
	if len(ds) != 1 || ds[0].ID != id {
		t.Fatalf("mon cua chu phai con nguyen: %+v", ds)
	}
}

// Ly do la bi chan truoc khi cham CSDL.
func TestHuyVatPhamLyDoLa(t *testing.T) {
	s := moKhoVatPham(t)
	nv := taoNhanVatThu(t, s, 0)
	donKhoCua(t, s, nv.ID)
	if err := s.HuyVatPham(context.Background(), nv.ID, NoiChuaTui, 1, "XOA_CHOI"); !errors.Is(err, ErrLyDoKhongHopLe) {
		t.Fatalf("mong doi ErrLyDoKhongHopLe, duoc %v", err)
	}
}
