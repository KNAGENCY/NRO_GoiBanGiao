package game

// muaban_tichhop_test.go — KH7 Task 9: test duong mua tren CSDL THAT.
//
// Moi test o day tra loi MOT cau hoi "cai gi khong bao gio duoc xay ra", va moi
// cau tuong ung mot dot bien bat buoc (ghi ro trong chu thich tung test):
//
//   DB1 tach "tru tien" va "cap mon" thanh HAI giao dich
//        -> TestMuaThatBaiKhongMatTien do
//   DB2 bo ghi/kiem khoa idempotency
//        -> TestMuaGuiLapChiTruTienMotLan do
//   DB3 tru tien TRUOC khi kiem tui day
//        -> TestMuaTuiDayTienKhongDoi do

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"nroserver/internal/store"
)

// ---------------------------------------------------------------------------
// Ha tang
// ---------------------------------------------------------------------------

// moKhoMuaBan mo CSDL that VA ap moi migration con lai (bang vat_pham /
// vat_pham_giao_dich den tu 0005, DB test cu co the chua co).
func moKhoMuaBan(t *testing.T) *store.Store {
	t.Helper()
	s := moKhoThatChoTest(t)
	if _, err := s.ApDungMigrationConLai(context.Background()); err != nil {
		t.Fatalf("ap migration that bai: %v", err)
	}
	return s
}

// taoNguoiMua tao mot nhan vat that voi so vang cho truoc, va dang ky don dep
// TAT CA dau vet cua no (giao dich, lich su, vat pham, nhan vat, tai khoan).
func taoNguoiMua(t *testing.T, s *store.Store, vang, ngoc int64) (*store.NhanVat, *pgxpool.Pool) {
	t.Helper()
	pool := moPoolThuChoTest(t)
	tkID := taoTaiKhoanCoMatKhau(t, "tkmua"+hauToNganNgauMuaBan(t), "matkhau123")

	nv, err := s.TaoNhanVat(context.Background(), store.NhanVatMoi{
		// Ten nhan vat phai khop `^[A-Za-z0-9]{5,10}$` (rang buoc 0002).
		TaiKhoanID: tkID, Ten: "Mua" + hauToNganNgauMuaBan(t),
		Lang: 0, KieuToc: 5, MapID: 1, ViTriX: 100, ViTriY: 384,
		HpToiDa: 150, ChakraToiDa: 150, Vang: vang,
	})
	if err != nil {
		t.Fatalf("tao nhan vat thu that bai: %v", err)
	}
	if ngoc != 0 {
		if err := s.CongTruTien(context.Background(), nv.ID, 0, ngoc); err != nil {
			t.Fatalf("nap ngoc that bai: %v", err)
		}
	}

	// Don dep: lich su TRUOC (khoa ngoai), roi vat pham, roi giao dich.
	// taoTaiKhoanCoMatKhau da dang ky xoa nhan_vat/tai_khoan; Cleanup chay LIFO
	// nen ham nay dang ky SAU se chay TRUOC — dung thu tu can.
	t.Cleanup(func() {
		ctx := context.Background()
		_, _ = pool.Exec(ctx,
			`DELETE FROM vat_pham_lich_su
			  WHERE tu_ai = $1 OR sang_ai = $1
			     OR vat_pham_id IN (SELECT id FROM vat_pham WHERE chu_so_huu_id = $1)`, nv.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM vat_pham WHERE chu_so_huu_id = $1`, nv.ID)
		_, _ = pool.Exec(ctx, `DELETE FROM vat_pham_giao_dich WHERE khoa_idem LIKE $1`,
			fmt.Sprintf("%d:%%", nv.ID))
	})
	return nv, pool
}

// hauToNganNgauMuaBan tra mot hau to CHI CHU-SO, 5 ky tu, duy nhat trong moi
// lan chay — de ghep vao ten nhan vat van khop `^[A-Za-z0-9]{5,10}$`.
func hauToNganNgauMuaBan(t *testing.T) string {
	t.Helper()
	demMuaBan++
	n := (int64(os.Getpid())*1000 + int64(demMuaBan)) % 60466176 // 36^5
	s := "00000" + strconv.FormatInt(n, 36)
	return s[len(s)-5:]
}

var demMuaBan int

// docTienMuaBan doc thang vang/ngoc tu CSDL (khong qua ham nghiep vu nao).
func docTienMuaBan(t *testing.T, pool *pgxpool.Pool, nvID int64) (int64, int64) {
	t.Helper()
	var vang, ngoc int64
	if err := pool.QueryRow(context.Background(),
		`SELECT vang, ngoc FROM nhan_vat WHERE id = $1`, nvID).Scan(&vang, &ngoc); err != nil {
		t.Fatalf("doc tien that bai: %v", err)
	}
	return vang, ngoc
}

// demGiaoDichMuaBan dem so dong khoa idempotency cua mot nhan vat.
func demGiaoDichMuaBan(t *testing.T, pool *pgxpool.Pool, nvID int64) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(),
		`SELECT count(*) FROM vat_pham_giao_dich WHERE khoa_idem LIKE $1`,
		fmt.Sprintf("%d:%%", nvID)).Scan(&n); err != nil {
		t.Fatalf("dem giao dich that bai: %v", err)
	}
	return n
}

// cuaHangThu: mot cua hang tu dung trong test — KHONG dung danh muc that, de
// test khong vo khi ai do doi gia cua hang mau.
//
//	mau 0 (cong don, TranChong 99) gia 100 vang
//	mau 4                          gia   5 ngoc
func cuaHangThu() *CuaHang {
	return &CuaHang{
		ID: 9001, Ten: "Cua hang thu",
		Tabs: []TabCuaHang{{
			Ten: "Thu",
			Mon: []MonHang{
				{MauID: 0, Gia: 100, Tien: TienVang},
				{MauID: 4, Gia: 5, Tien: TienNgoc},
			},
		}},
	}
}

// ---------------------------------------------------------------------------
// Duong vui
// ---------------------------------------------------------------------------

// Mua duoc: tien giam DUNG, tui co mon, lich su co MUA_NPC, va co DUNG MOT dong
// khoa idempotency.
func TestMuaThanhCongTruTienVaCapMon(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 1000, 0)
	ctx := context.Background()

	kq, err := MuaTuCuaHang(ctx, s, nv.ID, cuaHangThu(), 0, 0, 3, 1, 30)
	if err != nil {
		t.Fatalf("mua that bai: %v", err)
	}
	if kq.LaLapLai {
		t.Fatal("lan mua dau tien khong duoc danh dau la goi lap")
	}
	if kq.SoLuong != 3 || kq.MauID != 0 {
		t.Fatalf("ket qua sai: %+v", kq)
	}

	// Tien: 3 cai x 100 = 300.
	vang, ngoc := docTienMuaBan(t, pool, nv.ID)
	if vang != 700 || ngoc != 0 {
		t.Fatalf("vang=%d ngoc=%d, mong doi 700/0", vang, ngoc)
	}

	// Tui: dung MOT mon, o 0, so luong 3.
	tui, err := s.LayTui(ctx, nv.ID, store.NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	if len(tui) != 1 {
		t.Fatalf("tui phai co 1 mon, co %d", len(tui))
	}
	if tui[0].SoLuong != 3 || tui[0].OSo == nil || *tui[0].OSo != 0 {
		t.Fatalf("mon trong tui sai: %+v", tui[0])
	}

	// Lich su: dung mot dong, ly do MUA_NPC.
	var lyDo string
	var soDong int
	if err := pool.QueryRow(ctx,
		`SELECT count(*), coalesce(max(ly_do), '') FROM vat_pham_lich_su WHERE vat_pham_id = $1`,
		kq.VatPhamID).Scan(&soDong, &lyDo); err != nil {
		t.Fatalf("doc lich su that bai: %v", err)
	}
	if soDong != 1 || lyDo != store.LyDoMuaNPC {
		t.Fatalf("lich su: %d dong, ly do %q — mong doi 1 dong MUA_NPC", soDong, lyDo)
	}

	// Khoa idempotency: dung MOT dong.
	if n := demGiaoDichMuaBan(t, pool, nv.ID); n != 1 {
		t.Fatalf("phai co DUNG 1 dong vat_pham_giao_dich, dem duoc %d", n)
	}
}

// Mua bang ngoc: o tien con lai khong bi cham vao.
func TestMuaBangNgocKhongChamVang(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 1000, 20)
	ctx := context.Background()

	if _, err := MuaTuCuaHang(ctx, s, nv.ID, cuaHangThu(), 0, 4, 1, 1, 30); err != nil {
		t.Fatalf("mua bang ngoc that bai: %v", err)
	}
	vang, ngoc := docTienMuaBan(t, pool, nv.ID)
	if vang != 1000 || ngoc != 15 {
		t.Fatalf("vang=%d ngoc=%d, mong doi 1000/15", vang, ngoc)
	}
}

// ---------------------------------------------------------------------------
// DOT BIEN 1 — tach thanh 2 giao dich
// ---------------------------------------------------------------------------

// Khong du tien: KHONG mat tien, KHONG co mon nao vao tui, va KHONG de lai dong
// khoa idempotency chet (neu con, lan sau nguoi choi bam lai se bi coi la goi
// lap va khong bao gio mua duoc).
//
// DOT BIEN: neu tru tien duoc commit rieng truoc khi cap mon (hai giao dich),
// hoac neu dong khoa idempotency duoc commit rieng, test nay do.
func TestMuaThatBaiKhongMatTien(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 50, 0) // 50 vang, mua 1 cai gia 100
	ctx := context.Background()

	_, err := MuaTuCuaHang(ctx, s, nv.ID, cuaHangThu(), 0, 0, 1, 1, 30)
	if !errors.Is(err, store.ErrKhongDuTien) {
		t.Fatalf("mong doi ErrKhongDuTien, duoc %v", err)
	}

	vang, _ := docTienMuaBan(t, pool, nv.ID)
	if vang != 50 {
		t.Fatalf("vang = %d, phai giu nguyen 50 (giao dich ROLLBACK)", vang)
	}
	tui, err := s.LayTui(ctx, nv.ID, store.NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	if len(tui) != 0 {
		t.Fatalf("tui phai rong, co %d mon", len(tui))
	}
	if n := demGiaoDichMuaBan(t, pool, nv.ID); n != 0 {
		t.Fatalf("giao dich hong khong duoc de lai dong khoa nao, con %d", n)
	}

	// Va nguoi choi van mua duoc khi co du tien, VOI CHINH so thu tu do.
	if err := s.CongTruTien(ctx, nv.ID, 100, 0); err != nil {
		t.Fatalf("nap them vang that bai: %v", err)
	}
	if _, err := MuaTuCuaHang(ctx, s, nv.ID, cuaHangThu(), 0, 0, 1, 1, 30); err != nil {
		t.Fatalf("mua lai sau khi du tien phai duoc, loi: %v", err)
	}
}

// Cap mon HONG sau khi tien da tru duoc: tien phai quay ve nguyen ven.
//
// Day la kich ban thuc te: giua luc doc tui (thay o 0 trong) va luc ghi, mot
// duong khac chiem mat o do -> UNIQUE `vat_pham_mot_o_mot_mon` -> ErrOCoMon.
// Test goi thang store.MuaMotGiaoDich voi mot o DA CO MON de dung lai canh do
// mot cach tat dinh, khong phai cho dua tranh.
//
// DOT BIEN 1 (BAT BUOC): tach "tru tien" ra mot giao dich rieng (vd doi
// congTruTien(ctx, tx, ...) thanh s.CongTruTien(ctx, ...)) -> tien bi tru va
// COMMIT ngay, roi buoc cap mon vo -> nguoi choi mat 100 vang khong duoc gi ->
// test nay DO. Da kiem tay 03/09.
func TestMuaCapMonHongKhongMatTien(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 1000, 0)
	ctx := context.Background()

	// Chiem san o 0.
	if _, err := s.CapVatPham(ctx, nv.ID, store.NoiChuaTui, 0, 0, 1,
		store.ChiSoPhu{Opt: []store.TuyChon{}}, store.LyDoSinhRa); err != nil {
		t.Fatalf("chiem o 0 that bai: %v", err)
	}

	_, err := s.MuaMotGiaoDich(ctx, store.YeuCauMua{
		NhanVatID: nv.ID,
		KhoaIdem:  store.KhoaIdemMua(nv.ID, 1),
		DVang:     -100,
		NoiChua:   store.NoiChuaTui,
		OSo:       0, // <- o da co mon
		MauID:     0,
		SoLuong:   1,
		LyDo:      store.LyDoMuaNPC,
	})
	if !errors.Is(err, store.ErrOCoMon) {
		t.Fatalf("mong doi ErrOCoMon, duoc %v", err)
	}

	vang, _ := docTienMuaBan(t, pool, nv.ID)
	if vang != 1000 {
		t.Fatalf("vang = %d, phai giu nguyen 1000 — tru tien va cap mon PHAI cung mot giao dich", vang)
	}
	if n := demGiaoDichMuaBan(t, pool, nv.ID); n != 0 {
		t.Fatalf("giao dich hong khong duoc de lai dong khoa nao, con %d", n)
	}
}

// ---------------------------------------------------------------------------
// DOT BIEN 2 — bo khoa idempotency (D3)
// ---------------------------------------------------------------------------

// Gui LAP cung mot so thu tu: tien chi tru MOT lan, tui chi them MOT mon, va
// ket qua tra ve GIONG HET lan dau.
//
// DOT BIEN: bo INSERT/kiem `khoa_idem` -> lan hai tru them 100 vang va them mon
// thu hai -> test do o ca ba cho.
func TestMuaGuiLapChiTruTienMotLan(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 1000, 0)
	ctx := context.Background()

	kq1, err := MuaTuCuaHang(ctx, s, nv.ID, cuaHangThu(), 0, 0, 1, 7, 30)
	if err != nil {
		t.Fatalf("mua lan 1 that bai: %v", err)
	}
	kq2, err := MuaTuCuaHang(ctx, s, nv.ID, cuaHangThu(), 0, 0, 1, 7, 30)
	if err != nil {
		t.Fatalf("mua lan 2 (goi lap) phai thanh cong, loi: %v", err)
	}

	if !kq2.LaLapLai {
		t.Fatal("lan 2 phai duoc danh dau la goi lap")
	}
	if kq2.VatPhamID != kq1.VatPhamID || kq2.OSo != kq1.OSo ||
		kq2.SoLuong != kq1.SoLuong || kq2.DVang != kq1.DVang {
		t.Fatalf("goi lap phai tra DUNG ket qua lan dau:\n lan1 %+v\n lan2 %+v", kq1, kq2)
	}

	vang, _ := docTienMuaBan(t, pool, nv.ID)
	if vang != 900 {
		t.Fatalf("vang = %d, phai la 900 (chi tru MOT lan)", vang)
	}
	tui, err := s.LayTui(ctx, nv.ID, store.NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	if len(tui) != 1 {
		t.Fatalf("tui phai co DUNG 1 mon, co %d", len(tui))
	}
	if n := demGiaoDichMuaBan(t, pool, nv.ID); n != 1 {
		t.Fatalf("phai co DUNG 1 dong giao dich, dem duoc %d", n)
	}
}

// So thu tu KHAC nhau la hai lan mua that su — chot chan de dot bien "luon coi
// moi goi la goi lap" khong lot qua test tren.
func TestMuaSoThuTuKhacNhauLaHaiLanMua(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 1000, 0)
	ctx := context.Background()

	for i, stt := range []uint64{1, 2} {
		if _, err := MuaTuCuaHang(ctx, s, nv.ID, cuaHangThu(), 0, 0, 1, stt, 30); err != nil {
			t.Fatalf("mua lan %d that bai: %v", i+1, err)
		}
	}
	vang, _ := docTienMuaBan(t, pool, nv.ID)
	if vang != 800 {
		t.Fatalf("vang = %d, phai la 800 (hai lan mua that su)", vang)
	}
	tui, err := s.LayTui(ctx, nv.ID, store.NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	if len(tui) != 2 {
		t.Fatalf("tui phai co 2 mon, co %d", len(tui))
	}
}

// ---------------------------------------------------------------------------
// DOT BIEN 3 — tru tien TRUOC khi kiem tui day
// ---------------------------------------------------------------------------

// Tui day: bao loi ro, tien KHONG doi, khong co dong giao dich nao.
//
// DOT BIEN: doi cau kiem "con o trong khong" xuong SAU loi goi MuaMotGiaoDich
// (dung bay D16 cua ban Java) -> tien bi tru roi moi vo -> test do.
func TestMuaTuiDayTienKhongDoi(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 1000, 0)
	ctx := context.Background()

	// Tui chi co 2 o, va lap day ca hai.
	const soO = 2
	for o := int16(0); o < soO; o++ {
		if _, err := s.CapVatPham(ctx, nv.ID, store.NoiChuaTui, o, 0, 1,
			store.ChiSoPhu{Opt: []store.TuyChon{}}, store.LyDoSinhRa); err != nil {
			t.Fatalf("lap day o %d that bai: %v", o, err)
		}
	}

	_, err := MuaTuCuaHang(ctx, s, nv.ID, cuaHangThu(), 0, 0, 1, 1, soO)
	if !errors.Is(err, ErrTuiDay) {
		t.Fatalf("mong doi ErrTuiDay, duoc %v", err)
	}
	vang, _ := docTienMuaBan(t, pool, nv.ID)
	if vang != 1000 {
		t.Fatalf("vang = %d, phai giu nguyen 1000 (kiem tui TRUOC khi tru tien)", vang)
	}
	if n := demGiaoDichMuaBan(t, pool, nv.ID); n != 0 {
		t.Fatalf("tui day khong duoc de lai dong giao dich nao, con %d", n)
	}
}

// ---------------------------------------------------------------------------
// Gói rác: tu choi TRUOC khi cham CSDL
// ---------------------------------------------------------------------------

// So luong 0 / am / vuot tran, mau khong ban, typeBuy la, chua mo cua hang:
// tat ca phai bi tu choi va KHONG de lai dau vet nao trong CSDL.
func TestMuaGoiRacKhongChamCSDL(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 1000, 0)
	ctx := context.Background()
	ch := cuaHangThu()

	cacCa := []struct {
		ten     string
		ch      *CuaHang
		typeBuy int8
		mauID   int16
		soLuong int32
		mongDoi error
	}{
		{"chua mo cua hang", nil, 0, 0, 1, ErrKhongMoCuaHang},
		{"typeBuy la", ch, 9, 0, 1, ErrTypeMuaLa},
		{"so luong 0", ch, 0, 0, 0, ErrSoLuongKhongHopLe},
		{"so luong am", ch, 0, 0, -5, ErrSoLuongKhongHopLe},
		{"vuot tran chong", ch, 0, 0, 1000, ErrSoLuongKhongHopLe},
		{"mau khong ban", ch, 0, 1, 1, ErrKhongBanMonNay},
		{"don chiec mua 2", ch, 0, 4, 2, ErrSoLuongKhongHopLe},
	}
	for i, ca := range cacCa {
		_, err := MuaTuCuaHang(ctx, s, nv.ID, ca.ch, ca.typeBuy, ca.mauID, ca.soLuong,
			uint64(100+i), 30)
		if !errors.Is(err, ca.mongDoi) {
			t.Errorf("%s: mong doi %v, duoc %v", ca.ten, ca.mongDoi, err)
		}
	}

	vang, _ := docTienMuaBan(t, pool, nv.ID)
	if vang != 1000 {
		t.Fatalf("vang = %d, goi rac khong duoc cham vao tien", vang)
	}
	if n := demGiaoDichMuaBan(t, pool, nv.ID); n != 0 {
		t.Fatalf("goi rac khong duoc de lai dong giao dich nao, con %d", n)
	}
	tui, err := s.LayTui(ctx, nv.ID, store.NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	if len(tui) != 0 {
		t.Fatalf("tui phai rong, co %d mon", len(tui))
	}
}

// Mon co option thi LUON mot cai mot dong (§1.3), du mau khai TranChong lon.
func TestMuaMonCoOptionChiDuocMotCai(t *testing.T) {
	s := moKhoMuaBan(t)
	nv, _ := taoNguoiMua(t, s, 100000, 0)
	ctx := context.Background()

	ch := &CuaHang{ID: 9002, Ten: "Thu option", Tabs: []TabCuaHang{{
		Ten: "Thu",
		Mon: []MonHang{{
			MauID: 0, Gia: 10, Tien: TienVang, // mau 0 co TranChong 99
			Options: []OptionMon{{OptionID: 0, Param: 200}},
		}},
	}}}

	if _, err := MuaTuCuaHang(ctx, s, nv.ID, ch, 0, 0, 2, 1, 30); !errors.Is(err, ErrSoLuongKhongHopLe) {
		t.Fatalf("mua 2 cai mon co option phai bi tu choi, duoc %v", err)
	}
	kq, err := MuaTuCuaHang(ctx, s, nv.ID, ch, 0, 0, 1, 2, 30)
	if err != nil {
		t.Fatalf("mua 1 cai mon co option that bai: %v", err)
	}
	tui, err := s.LayTui(ctx, nv.ID, store.NoiChuaTui)
	if err != nil {
		t.Fatalf("lay tui that bai: %v", err)
	}
	if len(tui) != 1 || tui[0].ID != kq.VatPhamID {
		t.Fatalf("tui sai: %+v", tui)
	}
	if len(tui[0].ChiSoPhu.Opt) != 1 || tui[0].ChiSoPhu.Opt[0].Param != 200 {
		t.Fatalf("option khong duoc luu dung: %+v", tui[0].ChiSoPhu)
	}
}
