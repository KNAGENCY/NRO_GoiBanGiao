package game

// muaban_ban_tichhop_test.go — KH7 Task 10: test duong BAN tren CSDL THAT.
//
// Moi test tra loi MOT cau hoi "cai gi khong bao gio duoc xay ra", va moi cau
// tuong ung mot dot bien bat buoc:
//
//   DB1 tach "bot do" va "cong tien" thanh HAI giao dich
//        -> TestBanTienTranKhongMatDo do
//   DB2 bo khoa idempotency (ON CONFLICT DO NOTHING)
//        -> TestBanGuiLapChiCongTienMotLan do
//   DB3 bo guard chu so huu trong menh de WHERE (D4)
//        -> TestBanDoNguoiKhacKhongDoiGi do
//   DB4 doi cau kiem "du so luong" xuong SAU khi da bot do / cong tien
//        -> TestBanQuaSoLuongKhongChamTien do

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"nroserver/internal/store"
)

// ---------------------------------------------------------------------------
// Ha tang rieng cua duong ban
// ---------------------------------------------------------------------------

// datMonVaoTui nhet mot mon vao o `oSo` cua nhan vat bang chinh duong CapVatPham
// (duong DUY NHAT sinh mon — D9), roi tra lai lat cat tui.
func datMonVaoTui(t *testing.T, s *store.Store, nvID int64, oSo int16,
	mauID int32, soLuong int32) []store.VatPham {
	t.Helper()
	ctx := context.Background()
	if _, err := s.CapVatPham(ctx, nvID, store.NoiChuaTui, oSo, mauID, soLuong,
		store.ChiSoPhu{}, store.LyDoSinhRa); err != nil {
		t.Fatalf("dat mon vao tui that bai: %v", err)
	}
	tui, err := s.LayTui(ctx, nvID, store.NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	return tui
}

// demLichSuBan dem so dong lich su BAN_NPC cua mot mon.
func demLichSuBan(t *testing.T, pool *pgxpool.Pool, vatPhamID int64) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(),
		`SELECT count(*) FROM vat_pham_lich_su
		  WHERE vat_pham_id = $1 AND ly_do = $2 AND sang_ai IS NULL`,
		vatPhamID, store.LyDoBanNPC).Scan(&n); err != nil {
		t.Fatalf("dem lich su ban that bai: %v", err)
	}
	return n
}

// ---------------------------------------------------------------------------
// Duong vui
// ---------------------------------------------------------------------------

// Ban ca chong: tien tang DUNG (gia mua 100 -> gia ban 25), mon bien khoi tui,
// lich su co dong BAN_NPC voi sang_ai NULL, va co DUNG MOT dong khoa idem.
func TestBanThanhCongCongTienVaHetMon(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 1000, 0)
	ctx := context.Background()
	tui := datMonVaoTui(t, s, nv.ID, 0, 0, 3)

	kq, err := BanChoCuaHang(ctx, s, nv.ID, cuaHangThu(), tui, 0, 0, 1)
	if err != nil {
		t.Fatalf("ban that bai: %v", err)
	}
	if kq.LaLapLai {
		t.Fatal("lan ban dau tien khong duoc danh dau la goi lap")
	}
	// soLuong = 0 -> ban ca chong 3 cai, 3 x (100/4) = 75.
	if kq.SoLuongBan != 3 || !kq.DaHuyMon || kq.SoLuongCon != 0 {
		t.Fatalf("ket qua sai: %+v", kq)
	}
	if kq.DVang != 75 || kq.DNgoc != 0 {
		t.Fatalf("tien ban: vang=%d ngoc=%d, mong doi 75/0", kq.DVang, kq.DNgoc)
	}

	vang, ngoc := docTienMuaBan(t, pool, nv.ID)
	if vang != 1075 || ngoc != 0 {
		t.Fatalf("vang=%d ngoc=%d, mong doi 1075/0", vang, ngoc)
	}

	tuiSau, err := s.LayTui(ctx, nv.ID, store.NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	if len(tuiSau) != 0 {
		t.Fatalf("tui phai rong sau khi ban ca chong, con %d mon", len(tuiSau))
	}
	if n := demLichSuBan(t, pool, kq.VatPhamID); n != 1 {
		t.Fatalf("phai co DUNG 1 dong lich su BAN_NPC (sang_ai NULL), dem duoc %d", n)
	}
	if n := demGiaoDichMuaBan(t, pool, nv.ID); n != 1 {
		t.Fatalf("phai co DUNG 1 dong vat_pham_giao_dich, dem duoc %d", n)
	}
}

// Ban MOT PHAN chong: mon van con, so luong giam dung, tien tang dung.
func TestBanMotPhanChongThiGiamSoLuong(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 0, 0)
	ctx := context.Background()
	tui := datMonVaoTui(t, s, nv.ID, 0, 0, 5)

	kq, err := BanChoCuaHang(ctx, s, nv.ID, cuaHangThu(), tui, 0, 2, 1)
	if err != nil {
		t.Fatalf("ban mot phan that bai: %v", err)
	}
	if kq.DaHuyMon || kq.SoLuongCon != 3 || kq.DVang != 50 {
		t.Fatalf("ket qua sai: %+v", kq)
	}
	vang, _ := docTienMuaBan(t, pool, nv.ID)
	if vang != 50 {
		t.Fatalf("vang = %d, mong doi 50", vang)
	}
	tuiSau, err := s.LayTui(ctx, nv.ID, store.NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	if len(tuiSau) != 1 || tuiSau[0].SoLuong != 3 {
		t.Fatalf("mon con lai sai: %+v", tuiSau)
	}
	if n := demLichSuBan(t, pool, kq.VatPhamID); n != 1 {
		t.Fatalf("ban mot phan van phai ghi 1 dong lich su BAN_NPC, dem duoc %d", n)
	}
}

// Ban bang ngoc: o vang khong bi cham vao.
func TestBanLayNgocKhongChamVang(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 1000, 0)
	ctx := context.Background()
	tui := datMonVaoTui(t, s, nv.ID, 0, 4, 1)

	// Mau 4 gia mua 5 ngoc -> gia ban = max(1, 5/4) = 1 ngoc.
	kq, err := BanChoCuaHang(ctx, s, nv.ID, cuaHangThu(), tui, 0, 0, 1)
	if err != nil {
		t.Fatalf("ban bang ngoc that bai: %v", err)
	}
	if kq.DNgoc != 1 || kq.DVang != 0 {
		t.Fatalf("tien ban: vang=%d ngoc=%d, mong doi 0/1", kq.DVang, kq.DNgoc)
	}
	vang, ngoc := docTienMuaBan(t, pool, nv.ID)
	if vang != 1000 || ngoc != 1 {
		t.Fatalf("vang=%d ngoc=%d, mong doi 1000/1", vang, ngoc)
	}
}

// ---------------------------------------------------------------------------
// DOT BIEN 1 — tach "bot do" va "cong tien" thanh hai giao dich
// ---------------------------------------------------------------------------

// Cong tien THAT BAI o buoc cuoi (so tien vuot tran bigint) -> do VAN CON
// nguyen trong tui, khong dong lich su nao, khong dong khoa idem nao.
//
// DOT BIEN: neu buoc huy/bot do duoc commit rieng truoc buoc cong tien, test
// nay do — nguoi choi mat do ma khong duoc dong nao.
func TestBanTienTranKhongMatDo(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 1000, 0)
	ctx := context.Background()
	tui := datMonVaoTui(t, s, nv.ID, 0, 0, 3)

	_, err := s.BanMotGiaoDich(ctx, store.YeuCauBan{
		NhanVatID: nv.ID,
		KhoaIdem:  store.KhoaIdemBan(nv.ID, 1),
		VatPhamID: tui[0].ID,
		NoiChua:   store.NoiChuaTui,
		OSo:       0,
		SoLuong:   3,
		DVang:     math.MaxInt64, // 1000 + MaxInt64 -> tran bigint -> loi SQL
		LyDo:      store.LyDoBanNPC,
	})
	if err == nil {
		t.Fatal("cong tien vuot tran bigint phai bao loi")
	}

	tuiSau, err := s.LayTui(ctx, nv.ID, store.NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	if len(tuiSau) != 1 || tuiSau[0].SoLuong != 3 {
		t.Fatalf("mon phai con nguyen 3 cai (giao dich ROLLBACK), thay: %+v", tuiSau)
	}
	if n := demLichSuBan(t, pool, tui[0].ID); n != 0 {
		t.Fatalf("giao dich hong khong duoc de lai dong lich su nao, con %d", n)
	}
	if n := demGiaoDichMuaBan(t, pool, nv.ID); n != 0 {
		t.Fatalf("giao dich hong khong duoc de lai khoa idem nao, con %d", n)
	}
	vang, _ := docTienMuaBan(t, pool, nv.ID)
	if vang != 1000 {
		t.Fatalf("vang = %d, phai giu nguyen 1000", vang)
	}
}

// ---------------------------------------------------------------------------
// DOT BIEN 2 — bo khoa idempotency
// ---------------------------------------------------------------------------

// Gui LAP cung so thu tu: tien chi duoc cong MOT lan, do chi mat MOT lan.
func TestBanGuiLapChiCongTienMotLan(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 0, 0)
	ctx := context.Background()
	tui := datMonVaoTui(t, s, nv.ID, 0, 0, 4)

	kq1, err := BanChoCuaHang(ctx, s, nv.ID, cuaHangThu(), tui, 0, 2, 7)
	if err != nil {
		t.Fatalf("ban lan 1 that bai: %v", err)
	}
	if kq1.LaLapLai {
		t.Fatal("lan 1 khong duoc danh dau la goi lap")
	}

	// Gui LAI y het, CUNG so thu tu 7 (lat cat tui cu — dung nhu goi lap that).
	kq2, err := BanChoCuaHang(ctx, s, nv.ID, cuaHangThu(), tui, 0, 2, 7)
	if err != nil {
		t.Fatalf("ban lan 2 (goi lap) khong duoc bao loi: %v", err)
	}
	if !kq2.LaLapLai {
		t.Fatal("lan 2 phai duoc danh dau la goi lap")
	}

	vang, _ := docTienMuaBan(t, pool, nv.ID)
	if vang != 50 {
		t.Fatalf("vang = %d, phai la 50 (chi cong MOT lan)", vang)
	}
	tuiSau, err := s.LayTui(ctx, nv.ID, store.NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	if len(tuiSau) != 1 || tuiSau[0].SoLuong != 2 {
		t.Fatalf("mon phai con 2 cai (chi bot MOT lan), thay: %+v", tuiSau)
	}
	if n := demGiaoDichMuaBan(t, pool, nv.ID); n != 1 {
		t.Fatalf("phai co DUNG 1 dong khoa idem, dem duoc %d", n)
	}
}

// ---------------------------------------------------------------------------
// DOT BIEN 3 — bo guard chu so huu (D4)
// ---------------------------------------------------------------------------

// Ban do CUA NGUOI KHAC: bi tu choi, va KHONG doi gi ca — do nguoi kia con
// nguyen, tien ca hai ben khong doi.
//
// DOT BIEN: bo `chu_so_huu_id = $2` khoi menh de WHERE cua cau khoa/UPDATE ->
// test nay do.
func TestBanDoNguoiKhacKhongDoiGi(t *testing.T) {
	s := moKhoMuaBan(t)
	keGian, pool := taoNguoiMua(t, s, 1000, 0)
	nanNhan, _ := taoNguoiMua(t, s, 500, 0)
	ctx := context.Background()

	tuiNanNhan := datMonVaoTui(t, s, nanNhan.ID, 0, 0, 3)

	_, err := s.BanMotGiaoDich(ctx, store.YeuCauBan{
		NhanVatID: keGian.ID, // ke gian
		KhoaIdem:  store.KhoaIdemBan(keGian.ID, 1),
		VatPhamID: tuiNanNhan[0].ID, // do cua nan nhan
		NoiChua:   store.NoiChuaTui,
		OSo:       0,
		SoLuong:   3,
		DVang:     75,
		LyDo:      store.LyDoBanNPC,
	})
	if !errors.Is(err, store.ErrOTrong) {
		t.Fatalf("mong doi ErrOTrong (guard D4), duoc %v", err)
	}

	tuiSau, err := s.LayTui(ctx, nanNhan.ID, store.NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui nan nhan that bai: %v", err)
	}
	if len(tuiSau) != 1 || tuiSau[0].SoLuong != 3 {
		t.Fatalf("do cua nan nhan phai con nguyen, thay: %+v", tuiSau)
	}
	vangKeGian, _ := docTienMuaBan(t, pool, keGian.ID)
	if vangKeGian != 1000 {
		t.Fatalf("ke gian khong duoc cong dong nao, vang = %d", vangKeGian)
	}
	vangNanNhan, _ := docTienMuaBan(t, pool, nanNhan.ID)
	if vangNanNhan != 500 {
		t.Fatalf("tien nan nhan khong duoc dong toi, vang = %d", vangNanNhan)
	}
}

// ---------------------------------------------------------------------------
// DOT BIEN 4 — kiem so luong sau khi da cham tien
// ---------------------------------------------------------------------------

// Ban nhieu hon so dang co: bi tu choi TRUOC khi cham tien hay cham do.
func TestBanQuaSoLuongKhongChamTien(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 100, 0)
	ctx := context.Background()
	tui := datMonVaoTui(t, s, nv.ID, 0, 0, 2)

	_, err := BanChoCuaHang(ctx, s, nv.ID, cuaHangThu(), tui, 0, 5, 1)
	if !errors.Is(err, store.ErrKhongDuSoLuong) {
		t.Fatalf("mong doi ErrKhongDuSoLuong, duoc %v", err)
	}
	vang, _ := docTienMuaBan(t, pool, nv.ID)
	if vang != 100 {
		t.Fatalf("vang = %d, phai giu nguyen 100", vang)
	}
	tuiSau, err := s.LayTui(ctx, nv.ID, store.NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	if len(tuiSau) != 1 || tuiSau[0].SoLuong != 2 {
		t.Fatalf("mon phai con nguyen 2 cai, thay: %+v", tuiSau)
	}
	if n := demGiaoDichMuaBan(t, pool, nv.ID); n != 0 {
		t.Fatalf("tu choi som khong duoc de lai khoa idem nao, con %d", n)
	}
}

// ---------------------------------------------------------------------------
// Tu choi som — khong can CSDL ghi
// ---------------------------------------------------------------------------

// Chua mo cua hang nao ma da gui goi `7`: tu choi.
func TestBanChuaMoCuaHangBiTuChoi(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 0, 0)
	ctx := context.Background()
	tui := datMonVaoTui(t, s, nv.ID, 0, 0, 1)

	_, err := BanChoCuaHang(ctx, s, nv.ID, nil, tui, 0, 0, 1)
	if !errors.Is(err, ErrKhongMoCuaHang) {
		t.Fatalf("mong doi ErrKhongMoCuaHang, duoc %v", err)
	}
	tuiSau, err := s.LayTui(ctx, nv.ID, store.NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	if len(tuiSau) != 1 {
		t.Fatalf("mon phai con nguyen, thay: %+v", tuiSau)
	}
	if n := demGiaoDichMuaBan(t, pool, nv.ID); n != 0 {
		t.Fatalf("tu choi som khong duoc de lai khoa idem nao, con %d", n)
	}
}

// Ban o TRONG: tu choi, khong cham CSDL.
func TestBanORongBiTuChoi(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, _ := taoNguoiMua(t, s, 0, 0)
	ctx := context.Background()
	tui := datMonVaoTui(t, s, nv.ID, 0, 0, 1)

	if _, err := BanChoCuaHang(ctx, s, nv.ID, cuaHangThu(), tui, 5, 0, 1); !errors.Is(err, ErrBanOTrong) {
		t.Fatalf("mong doi ErrBanOTrong, duoc %v", err)
	}
}

// Cua hang dang mo khong ban mau nay -> cung khong mua lai (quyet dinh ban 1,
// xem PHA A muc 3).
func TestBanMonCuaHangKhongBanBiTuChoi(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, _ := taoNguoiMua(t, s, 0, 0)
	ctx := context.Background()
	// Mau 1 khong nam trong cuaHangThu (chi co mau 0 va 4).
	tui := datMonVaoTui(t, s, nv.ID, 0, 1, 1)

	if _, err := BanChoCuaHang(ctx, s, nv.ID, cuaHangThu(), tui, 0, 0, 1); !errors.Is(err, ErrBanKhongMuaMon) {
		t.Fatalf("mong doi ErrBanKhongMuaMon, duoc %v", err)
	}
}

// ---------------------------------------------------------------------------
// Gia ban — ham thuan, khong can CSDL
// ---------------------------------------------------------------------------

func TestBanDonGiaChiaBonVaSan1(t *testing.T) {
	cac := []struct {
		giaMua, mongDoi int32
	}{
		{100, 25}, {4, 1}, {3, 1}, {1, 1}, {0, 1}, {7, 1}, {8, 2},
	}
	for _, c := range cac {
		if duoc := banDonGia(c.giaMua); duoc != c.mongDoi {
			t.Errorf("banDonGia(%d) = %d, mong doi %d", c.giaMua, duoc, c.mongDoi)
		}
	}
}
