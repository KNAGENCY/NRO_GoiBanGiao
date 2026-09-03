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
//   DB4 bo cau kiem "du so luong" o TANG KHO
//        -> TestBanQuaSoLuongTangKhoCungTuChoi do
//        (TestBanQuaSoLuongKhongChamTien giu chot chan o tang tren)
//
// DB1 da kiem tay: doi huyVatPhamTrongTx(ctx, tx, ...) thanh s.HuyVatPham(ctx,
// ...) — tuc huy mon thanh giao dich RIENG — lam TestBanTienTranKhongMatDo treo
// (giao dich rieng doi khoa FOR UPDATE ma giao dich ngoai dang giu) roi do khi
// het gio. Treo cung la do; dieu quan trong la no KHONG con xanh.

import (
	"context"
	"errors"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"nroserver/internal/noidung"
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

// Cung cau kiem do, nhung o TANG KHO — chot chan thu hai, khong dua vao tang
// tren goi dung. Test tren (BanChoCuaHang) bi banTinhToan chan tu som nen KHONG
// cham toi cau kiem trong BanMotGiaoDich; test nay moi cham.
func TestBanQuaSoLuongTangKhoCungTuChoi(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 100, 0)
	ctx := context.Background()
	tui := datMonVaoTui(t, s, nv.ID, 0, 0, 2)

	_, err := s.BanMotGiaoDich(ctx, store.YeuCauBan{
		NhanVatID: nv.ID,
		KhoaIdem:  store.KhoaIdemBan(nv.ID, 1),
		VatPhamID: tui[0].ID,
		NoiChua:   store.NoiChuaTui,
		OSo:       0,
		SoLuong:   5, // chi co 2
		DVang:     125,
		LyDo:      store.LyDoBanNPC,
	})
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

// TestBanDonGiaDocGiaBanLaiCuaKhoNoiDung ghim gia thu lai THAT cua cua hang 1,
// doc tu kho-noi-dung/cua-hang.yaml.
//
// KH8 task 5 doi banDonGia tu cong thuc `gia/4` sang doc thang truong
// `gia_ban_lai`. Test cu chi kiem phep chia — no van xanh du du lieu YAML noi
// gi. Test nay kiem CAI DUNG: con so nguoi choi nhan khi ban lai la con so
// trong YAML, khong phai mot phep tinh trong code.
//
// Con so mong doi duoi day = gia/4 cua bang gia hien tai (100/1000/5). Doi gia
// trong YAML thi PHAI doi o day — do la muc dich.
func TestBanDonGiaDocGiaBanLaiCuaKhoNoiDung(t *testing.T) {
	ch := TimCuaHang(1)
	if ch == nil {
		t.Fatal("khong tim thay cua hang 1 trong kho noi dung")
	}
	mongDoi := map[int16]int32{0: 25, 1: 25, 2: 250, 3: 250, 4: 1}
	daThay := 0
	for _, tab := range ch.Tabs {
		for _, mon := range tab.Mon {
			can, co := mongDoi[mon.MauID]
			if !co {
				t.Errorf("cua hang 1 co mau %d ma test chua ghim gia ban lai", mon.MauID)
				continue
			}
			if duoc := banDonGia(mon); duoc != can {
				t.Errorf("banDonGia(mau %d) = %d, mong doi %d", mon.MauID, duoc, can)
			}
			daThay++
		}
	}
	if daThay != len(mongDoi) {
		t.Errorf("duyet duoc %d mon, cua hang 1 phai co %d", daThay, len(mongDoi))
	}
}

// TestBanDonGiaDoiTheoYAML la DOT BIEN "banDonGia quay ve gia/4": nap mot ban
// kho noi dung da sua gia_ban_lai cua mon dau tien thanh mot so KHAC gia/4, roi
// khang dinh duong ban tra dung so MOI do. Neu ai do doi banDonGia ve cong thuc
// cu thi test nay do (ma test tren van xanh, vi 25 tinh cờ == 100/4).
func TestBanDonGiaDoiTheoYAML(t *testing.T) {
	thuMuc := chepKhoNoiDungRaThuMucTam(t)
	doiTrongFileCuaHang(t, thuMuc, "gia_ban_lai: 25", "gia_ban_lai: 7")

	du, err := noidung.NapTuThuMuc(thuMuc)
	if err != nil {
		t.Fatalf("nap kho da sua that bai: %v", err)
	}
	ds, err := cuaHangTuYAML(du)
	if err != nil {
		t.Fatalf("doi cua hang: %v", err)
	}
	var mon *MonHang
	for i := range ds {
		if ds[i].ID != 1 {
			continue
		}
		for j := range ds[i].Tabs[0].Mon {
			if ds[i].Tabs[0].Mon[j].MauID == 0 {
				mon = &ds[i].Tabs[0].Mon[j]
			}
		}
	}
	if mon == nil {
		t.Fatal("khong tim thay mon mau 0 cua cua hang 1")
	}
	if mon.Gia != 100 {
		t.Fatalf("gia mua = %d, mong doi 100 — dot bien da lam sai cai khac", mon.Gia)
	}
	if duoc := banDonGia(*mon); duoc != 7 {
		t.Errorf("banDonGia = %d, mong doi 7 — duong ban dang KHONG doc gia_ban_lai "+
			"(mot cong thuc gia/4 se tra 25)", duoc)
	}
}

// TestKiemChanGiaBanLaiCaoHonGiaMua la DOT BIEN "bo tran tren cua gia_ban_lai":
// mua 100 ban 120 la may in tien — bo kiem PHAI tu choi nap.
func TestKiemChanGiaBanLaiCaoHonGiaMua(t *testing.T) {
	thuMuc := chepKhoNoiDungRaThuMucTam(t)
	doiTrongFileCuaHang(t, thuMuc, "gia_ban_lai: 25", "gia_ban_lai: 120")

	if _, err := noidung.NapTuThuMuc(thuMuc); err == nil {
		t.Fatal("gia_ban_lai (120) > gia (100) VAN NAP DUOC — nguoi choi mua 100 ban 120, " +
			"lap vo han, kinh te chet ma khong mot dong log nao bao")
	}
}

// TestXuLyMuaTuChoiKhiXoaCuaHangKhoiYAML la DOT BIEN "xoa cua hang id 1":
// TimCuaHang phai tra nil (khong panic), va duong mua voi cua hang nil phai TU
// CHOI chu khong no.
func TestXuLyMuaTuChoiKhiXoaCuaHangKhoiYAML(t *testing.T) {
	thuMuc := chepKhoNoiDungRaThuMucTam(t)
	ghiFileCuaHang(t, thuMuc, "cua_hang: []\n")

	du, err := noidung.NapTuThuMuc(thuMuc)
	if err != nil {
		t.Fatalf("kho khong con cua hang nao PHAI van nap duoc: %v", err)
	}
	ds, err := cuaHangTuYAML(du)
	if err != nil {
		t.Fatalf("doi cua hang: %v", err)
	}
	if len(ds) != 0 {
		t.Fatalf("mong doi 0 cua hang, duoc %d", len(ds))
	}

	// Cua tra cuu THAT (doc anh chup) phai tra nil, khong panic.
	ac := noidung.MoKho(du).AnhChup()
	if ch := timCuaHangTu(ac, 1); ch != nil {
		t.Fatalf("cua hang 1 da bi xoa khoi YAML ma timCuaHangTu van tra %+v", ch)
	}
	if n := len(cuaHangTu(ac)); n != 0 {
		t.Fatalf("cuaHangTu tra %d cua hang tren mot kho khong con cua hang nao", n)
	}

	// Cua hang dang mo = nil (dung canh "chua mo cua hang"): ca hai duong phai
	// tra loi, khong panic.
	if _, err := banTinhToan(nil, nil, 0, 0); err == nil {
		t.Error("banTinhToan voi cua hang nil khong bao loi")
	}
	// ch == nil bi chan o Kiem 1, TRUOC khi cham toi kho — nen kho nil o day la
	// an toan va test nay khong can CSDL.
	if _, err := MuaTuCuaHang(context.Background(), nil, 0, nil, 0, 0, 1, 1, soOTuiChuan); !errors.Is(err, ErrKhongMoCuaHang) {
		t.Errorf("duong mua voi cua hang nil: mong doi ErrKhongMoCuaHang, duoc %v", err)
	}
}

// doiTrongFileCuaHang doi lan xuat hien DAU TIEN cua `cu` trong cua-hang.yaml
// cua mot ban chep kho noi dung. Fatal neu khong tim thay — mot dot bien khong
// va vao dau la mot dot bien VO NGHIA.
func doiTrongFileCuaHang(t *testing.T, thuMuc, cu, moi string) {
	t.Helper()
	duongDan := filepath.Join(thuMuc, noidung.TenFileCuaHang)
	data, err := os.ReadFile(duongDan)
	if err != nil {
		t.Fatalf("doc %s: %v", duongDan, err)
	}
	s := string(data)
	if !strings.Contains(s, cu) {
		t.Fatalf("khong tim thay %q trong %s — dot bien khong va vao dau", cu, duongDan)
	}
	if err := os.WriteFile(duongDan, []byte(strings.Replace(s, cu, moi, 1)), 0o644); err != nil {
		t.Fatalf("ghi %s: %v", duongDan, err)
	}
}

// ghiFileCuaHang ghi de toan bo cua-hang.yaml cua mot ban chep kho noi dung.
func ghiFileCuaHang(t *testing.T, thuMuc, noiDung string) {
	t.Helper()
	duongDan := filepath.Join(thuMuc, noidung.TenFileCuaHang)
	if err := os.WriteFile(duongDan, []byte(noiDung), 0o644); err != nil {
		t.Fatalf("ghi %s: %v", duongDan, err)
	}
}
