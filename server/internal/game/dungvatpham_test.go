package game

// dungvatpham_test.go — KH7 Task 7: test duong BO DO / CHUYEN O / MAC THAO
// tren CSDL THAT.
//
// Moi test tra loi MOT cau "cai gi khong bao gio duoc xay ra". BA DOT BIEN bat
// buoc (da kiem tay 03/09):
//
//	DB1 dungvatpham.go: bo cau kiem bien trong dungDoKiemBien (cho `oSo` di
//	    thang toi bang / CSDL)
//	    -> TestKiemBienPhanBietVoiOTrong +
//	       TestChuyenOMonNamNgoaiBienBangMacBiTuChoi do
//
//	    ⚠️ LUU Y CHO NGUOI SAU: hai test "OSoNgoaiBienBiTuChoi" di qua handler
//	    KHONG giet duoc dot bien nay — da thu tay 03/09 va chung VAN XANH. Ly
//	    do: bang phien la mot LUOI THU HAI, no chi chua o hop le nen `oSo` = -1
//	    bi tu choi vi "o rong" ngay ca khi cau kiem bien da bien mat. Muon giet
//	    that su phai (a) kiem THANG loai loi tra ve, hoac (b) dung mot mon nam
//	    o o NGOAI BIEN nhung CO trong bang (CSDL khong rang buoc o_so < 5 cho
//	    MAC). Hai test moi lam dung hai viec do.
//	DB2 dungvatpham.go: sau khi HuyVatPham thanh cong, KHONG goi
//	    dungDoLamMoiVaGui (khong dung lai bang phien)
//	    -> TestBoDoXongBangPhienKhongConO do
//	DB3 dungvatpham.go: bo cau kiem `mau.Loai >= SoOMacChuan` trong
//	    chuyenOMacTrangBi
//	    -> TestMacMonKhongPhaiTrangBiBiTuChoi do

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	mynet "nroserver/internal/net"
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// ---------------------------------------------------------------------------
// Ha tang
// ---------------------------------------------------------------------------

// dungDoDungCanh dung MOT canh thu day du: CSDL that, mot nhan vat that, mot
// phien mang that, va mot BoXuLy da nap hanh trang (bang phien dung san).
type dungDoCanh struct {
	s      *store.Store
	pool   *pgxpool.Pool
	nv     *store.NhanVat
	b      *BoXuLy
	phien  *mynet.Session
	docGoi func() []*protocol.Message
}

func dungDoDungCanh(t *testing.T) *dungDoCanh {
	t.Helper()
	s := moKhoMuaBan(t)
	nv, pool := taoNguoiMua(t, s, 1000, 0)
	phien, docGoi := phienGia(t)

	b := NewBoXuLy(s, nil)
	b.nhanVat = nv
	c := &dungDoCanh{s: s, pool: pool, nv: nv, b: b, phien: phien, docGoi: docGoi}
	c.napLaiBang(t)
	return c
}

// datMon dat MOT mon vao (noiChua, oSo) qua duong sinh mon duy nhat cua kho.
func (c *dungDoCanh) datMon(t *testing.T, noiChua string, oSo int16, mauID int32) int64 {
	t.Helper()
	id, err := c.s.CapVatPham(context.Background(), c.nv.ID, noiChua, oSo,
		mauID, 1, store.ChiSoPhu{}, store.LyDoAdmin)
	if err != nil {
		t.Fatalf("dat mon vao %s o %d that bai: %v", noiChua, oSo, err)
	}
	return id
}

// napLaiBang dung lai bang phien tu CSDL — mo phong dung cai ma duong VAO THE
// GIOI lam, de test bat dau tu mot phien "dung nhu that".
func (c *dungDoCanh) napLaiBang(t *testing.T) {
	t.Helper()
	ds, err := c.b.hanhTrangNap(c.nv.ID)
	if err != nil {
		t.Fatalf("nap hanh trang that bai: %v", err)
	}
	bang, err := DungBangHanhTrang(ds)
	if err != nil {
		t.Fatalf("dung bang hanh trang that bai: %v", err)
	}
	c.b.hanhTrang = bang
}

// goiLen dung mot goi tin LEN tu cac byte tho.
func goiLenDungDo(cmd int8, than ...byte) *protocol.Message {
	return &protocol.Message{Command: cmd, Data: than}
}

// docNoiChuaCuaMon doc thang tu CSDL: mon `id` dang o dau, o nao.
func (c *dungDoCanh) docNoiChuaCuaMon(t *testing.T, id int64) (string, *int16) {
	t.Helper()
	var noi string
	var o *int16
	if err := c.pool.QueryRow(context.Background(),
		`SELECT noi_chua, o_so FROM vat_pham WHERE id = $1`, id).Scan(&noi, &o); err != nil {
		t.Fatalf("doc noi chua cua mon %d that bai: %v", id, err)
	}
	return noi, o
}

// demLichSu dem so dong lich su cua mot mon voi mot ly do cho truoc.
func (c *dungDoCanh) demLichSu(t *testing.T, id int64, lyDo string) int {
	t.Helper()
	var n int
	if err := c.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM vat_pham_lich_su WHERE vat_pham_id = $1 AND ly_do = $2`,
		id, lyDo).Scan(&n); err != nil {
		t.Fatalf("dem lich su that bai: %v", err)
	}
	return n
}

// timGoi tra goi DAU TIEN mang ma lenh cho truoc trong danh sach da gui.
func timGoi(ds []*protocol.Message, cmd int8) *protocol.Message {
	for _, m := range ds {
		if m.Command == cmd {
			return m
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// 1. BO DO
// ---------------------------------------------------------------------------

// TestBoDoXongBangPhienKhongConO — DOT BIEN DB2.
//
// Bo mot mon trong tui: dong CSDL phai chuyen sang DA_HUY, phai co dong lich su
// ly do BO, VA bang phien phai khong con tra ra o do nua. Cai cuoi la thu de
// quen nhat: neu quen dung lai bang, nguoi choi bam lai o do lan sau se cham
// vao mot id da bi huy.
func TestBoDoXongBangPhienKhongConO(t *testing.T) {
	c := dungDoDungCanh(t)
	id := c.datMon(t, store.NoiChuaTui, 3, 2)
	c.napLaiBang(t)

	if _, co := c.b.hanhTrang.TraId(store.NoiChuaTui, 3); !co {
		t.Fatal("truoc khi bo: bang phien phai co o 3")
	}

	// type 2 = xac nhan vut, where 1 = tui, o 3.
	if err := c.b.XuLy(c.phien, goiLenDungDo(lenhDungBoDo, 2, 1, 3)); err != nil {
		t.Fatalf("xu ly goi -43 that bai: %v", err)
	}

	noi, o := c.docNoiChuaCuaMon(t, id)
	if noi != "DA_HUY" || o != nil {
		t.Fatalf("sau khi bo: mon phai o DA_HUY voi o_so NULL, dang o %q o=%v", noi, o)
	}
	if n := c.demLichSu(t, id, store.LyDoBo); n != 1 {
		t.Fatalf("phai co dung 1 dong lich su ly do BO, co %d", n)
	}
	if _, co := c.b.hanhTrang.TraId(store.NoiChuaTui, 3); co {
		t.Fatal("sau khi bo: bang phien VAN tra ra o 3 — bang chua duoc dung lai (dot bien DB2)")
	}

	// Ba goi hanh trang phai duoc gui lai.
	ds := c.docGoi()
	for _, cmd := range []int8{lenhDoDangMac, lenhTui, lenhRuong} {
		if timGoi(ds, cmd) == nil {
			t.Fatalf("thieu goi %d sau khi bo do", cmd)
		}
	}
}

// TestBoDoOSoNgoaiBienBiTuChoi — DOT BIEN DB1 (bay §7.11).
//
// Ba gia tri `oSo` khong hop le: -1 (client CHU DONG gui gia tri nay),
// SoOTuiChuan (ngay ngoai bien tren), va 100 (xa hon). Ca ba phai bi tu choi
// TRUOC khi cham CSDL — kiem bang cach: khong mon nao bi dong vao.
func TestBoDoOSoNgoaiBienBiTuChoi(t *testing.T) {
	c := dungDoDungCanh(t)
	id := c.datMon(t, store.NoiChuaTui, 0, 2)
	c.napLaiBang(t)

	for _, oSo := range []int8{-1, int8(SoOTuiChuan), 100, -128} {
		if err := c.b.XuLy(c.phien, goiLenDungDo(lenhDungBoDo, 2, 1, byte(oSo))); err != nil {
			t.Fatalf("oSo=%d: xu ly goi that bai: %v", oSo, err)
		}
		noi, _ := c.docNoiChuaCuaMon(t, id)
		if noi != store.NoiChuaTui {
			t.Fatalf("oSo=%d: mon o o 0 bi dong toi (dang o %q) — cau kiem bien khong chan (dot bien DB1)",
				oSo, noi)
		}
	}
	// Bang phien phai con nguyen.
	if _, co := c.b.hanhTrang.TraId(store.NoiChuaTui, 0); !co {
		t.Fatal("bang phien mat o 0 sau cac goi bi tu choi")
	}
}

// TestBoDoTronORuongBiTuChoi: vut thang tu ruong khong duoc phep (dung y ban
// Java), nhung PHAI co cau tra loi -25 chu khong im lang.
func TestBoDoTuRuongBiTuChoi(t *testing.T) {
	c := dungDoDungCanh(t)
	id := c.datMon(t, store.NoiChuaRuong, 2, 2)
	c.napLaiBang(t)

	if err := c.b.XuLy(c.phien, goiLenDungDo(lenhDungBoDo, 2, 2, 2)); err != nil {
		t.Fatalf("xu ly that bai: %v", err)
	}
	noi, _ := c.docNoiChuaCuaMon(t, id)
	if noi != store.NoiChuaRuong {
		t.Fatalf("mon trong ruong bi dong toi: dang o %q", noi)
	}
	if timGoi(c.docGoi(), lenhDungDoBao) == nil {
		t.Fatal("phai gui -25 giai thich, khong duoc im lang")
	}
}

// TestBoDoHoiTruocKhiVut: type 1 (vut) phai tra ve goi -43 HOI LAI, va KHONG
// duoc dong vao CSDL. Thieu vong hai buoc nay thi nut "vut" cua client khong
// bao gio dan toi gap type 2 — tinh nang chet im lang.
func TestBoDoHoiTruocKhiVut(t *testing.T) {
	c := dungDoDungCanh(t)
	id := c.datMon(t, store.NoiChuaTui, 1, 2)
	c.napLaiBang(t)

	if err := c.b.XuLy(c.phien, goiLenDungDo(lenhDungBoDo, 1, 1, 1)); err != nil {
		t.Fatalf("xu ly that bai: %v", err)
	}
	noi, _ := c.docNoiChuaCuaMon(t, id)
	if noi != store.NoiChuaTui {
		t.Fatalf("buoc HOI khong duoc dong CSDL, mon dang o %q", noi)
	}
	g := timGoi(c.docGoi(), lenhDungBoDo)
	if g == nil {
		t.Fatal("phai gui goi -43 hoi lai")
	}
	r := protocol.NewReader(g.Data)
	loai, _ := r.Byte()
	where, _ := r.Byte()
	o, _ := r.Byte()
	cau, err := r.UTF()
	if err != nil {
		t.Fatalf("goi -43 hoi thieu chuoi: %v", err)
	}
	if loai != dungDoTypeVut || where != dungDoNoiTui || int8(o) != 1 {
		t.Fatalf("goi hoi sai truong: type=%d where=%d o=%d", loai, where, int8(o))
	}
	if cau == "" {
		t.Fatal("cau hoi rong")
	}
}

// TestDungThuocTraLoiTuTe: nhanh "dung" con hoan, nhung KHONG duoc im lang.
func TestDungThuocTraLoiTuTe(t *testing.T) {
	c := dungDoDungCanh(t)
	c.datMon(t, store.NoiChuaTui, 0, 0) // Loai 6 = thuoc
	c.napLaiBang(t)

	if err := c.b.XuLy(c.phien, goiLenDungDo(lenhDungBoDo, 0, 1, 0)); err != nil {
		t.Fatalf("xu ly that bai: %v", err)
	}
	if timGoi(c.docGoi(), lenhDungDoBao) == nil {
		t.Fatal("goi type 0 (dung) phai duoc tra loi bang -25, khong im lang")
	}
}

// ---------------------------------------------------------------------------
// 2. CHUYEN O
// ---------------------------------------------------------------------------

// TestChuyenOTuiSangRuongVaoOTrong: chuyen sang o trong nho nhat cua ruong.
func TestChuyenOTuiSangRuongVaoOTrong(t *testing.T) {
	c := dungDoDungCanh(t)
	id := c.datMon(t, store.NoiChuaTui, 4, 2)
	c.datMon(t, store.NoiChuaRuong, 0, 3) // chiem o 0 cua ruong
	c.napLaiBang(t)

	// loai 1 = tui -> ruong, o 4.
	if err := c.b.XuLy(c.phien, goiLenDungDo(lenhChuyenODo, 1, 4)); err != nil {
		t.Fatalf("xu ly -40 that bai: %v", err)
	}
	noi, o := c.docNoiChuaCuaMon(t, id)
	if noi != store.NoiChuaRuong || o == nil || *o != 1 {
		t.Fatalf("mon phai sang RUONG o 1 (o trong nho nhat), dang o %q o=%v", noi, o)
	}
	if _, co := c.b.hanhTrang.TraId(store.NoiChuaRuong, 1); !co {
		t.Fatal("bang phien phai co RUONG o 1 sau khi chuyen")
	}
	if _, co := c.b.hanhTrang.TraId(store.NoiChuaTui, 4); co {
		t.Fatal("bang phien van con TUI o 4 sau khi chuyen")
	}
}

// TestChuyenORuongDayBiTuChoi: het o dich -> tu choi, KHONG doi gi.
func TestChuyenORuongDayBiTuChoi(t *testing.T) {
	c := dungDoDungCanh(t)
	id := c.datMon(t, store.NoiChuaTui, 0, 2)
	for o := 0; o < SoORuongChuan; o++ {
		c.datMon(t, store.NoiChuaRuong, int16(o), 3)
	}
	c.napLaiBang(t)

	if err := c.b.XuLy(c.phien, goiLenDungDo(lenhChuyenODo, 1, 0)); err != nil {
		t.Fatalf("xu ly that bai: %v", err)
	}
	noi, o := c.docNoiChuaCuaMon(t, id)
	if noi != store.NoiChuaTui || o == nil || *o != 0 {
		t.Fatalf("ruong day: mon phai o nguyen TUI o 0, dang o %q o=%v", noi, o)
	}
	if timGoi(c.docGoi(), lenhDungDoBao) == nil {
		t.Fatal("phai gui -25 bao khong con o trong")
	}
}

// TestChuyenOOSoNgoaiBienBiTuChoi — DOT BIEN DB1, kiem O TUNG NHANH.
//
// Moi `loai` doc mot bang khac nhau voi so o khac nhau, nen moi nhanh phai co
// cau kiem bien RIENG. Test di qua CA BON nhanh tu dong tim o + nhanh mac.
func TestChuyenOOSoNgoaiBienBiTuChoi(t *testing.T) {
	c := dungDoDungCanh(t)
	idTui := c.datMon(t, store.NoiChuaTui, 0, 2)
	idMac := c.datMon(t, store.NoiChuaMac, 1, 3)
	idRuong := c.datMon(t, store.NoiChuaRuong, 0, 3)
	c.napLaiBang(t)

	truong := []struct {
		ten  string
		loai byte
		oSo  int8
	}{
		{"tui->ruong, o = -1", chuyenOTuiSangRuong, -1},
		{"tui->ruong, o = SoOTuiChuan", chuyenOTuiSangRuong, int8(SoOTuiChuan)},
		{"tui->ruong, o = 100", chuyenOTuiSangRuong, 100},
		{"ruong->tui, o = -1", chuyenORuongSangTui, -1},
		{"ruong->tui, o = SoORuongChuan", chuyenORuongSangTui, int8(SoORuongChuan)},
		{"mac->tui, o = -1", chuyenOMacSangTui, -1},
		{"mac->tui, o = SoOMacChuan", chuyenOMacSangTui, int8(SoOMacChuan)},
		{"mac->ruong, o = 100", chuyenOMacSangRuong, 100},
		{"tui->mac, o = -1", chuyenOTuiSangMac, -1},
		{"tui->mac, o = SoOTuiChuan", chuyenOTuiSangMac, int8(SoOTuiChuan)},
	}
	for _, tc := range truong {
		if err := c.b.XuLy(c.phien, goiLenDungDo(lenhChuyenODo, tc.loai, byte(tc.oSo))); err != nil {
			t.Fatalf("%s: xu ly that bai: %v", tc.ten, err)
		}
		for _, ct := range []struct {
			id  int64
			noi string
		}{{idTui, store.NoiChuaTui}, {idMac, store.NoiChuaMac}, {idRuong, store.NoiChuaRuong}} {
			noi, _ := c.docNoiChuaCuaMon(t, ct.id)
			if noi != ct.noi {
				t.Fatalf("%s: mon o %s bi dong toi (dang o %q) — thieu kiem bien (dot bien DB1)",
					tc.ten, ct.noi, noi)
			}
		}
	}
}

// TestKiemBienPhanBietVoiOTrong — DOT BIEN DB1, cach giet CHAC CHAN nhat.
//
// VI SAO CAN TEST NAY: bang phien la mot LUOI thu hai — no chi chua o hop le,
// nen mot `oSo` = -1 se "khong tim thay mon" va bi tu choi NGAY CA KHI cau kiem
// bien da bi xoa. Nghia la moi test di qua handler deu KHONG phan biet duoc hai
// truong hop. Cho nen phai kiem THANG loai loi: ngoai bien phai la
// ErrODoNgoaiBien, KHONG duoc rot xuong ErrONguonTrong.
func TestKiemBienPhanBietVoiOTrong(t *testing.T) {
	c := dungDoDungCanh(t)
	c.datMon(t, store.NoiChuaTui, 0, 2)
	c.napLaiBang(t)

	truong := []struct {
		noiChua string
		oSo     int16
	}{
		{store.NoiChuaTui, -1},
		{store.NoiChuaTui, int16(SoOTuiChuan)},
		{store.NoiChuaTui, 200},
		{store.NoiChuaMac, -1},
		{store.NoiChuaMac, int16(SoOMacChuan)},
		{store.NoiChuaRuong, -1},
		{store.NoiChuaRuong, int16(SoORuongChuan)},
	}
	for _, tc := range truong {
		_, err := c.b.dungDoKiemBien(tc.noiChua, tc.oSo)
		if !errors.Is(err, ErrODoNgoaiBien) {
			t.Fatalf("%s o %d: phai tra ErrODoNgoaiBien, tra %v (dot bien DB1)",
				tc.noiChua, tc.oSo, err)
		}
	}
	// O hop le nhung rong: phai la ErrONguonTrong, khong phai loi bien.
	if _, err := c.b.dungDoKiemBien(store.NoiChuaTui, 1); !errors.Is(err, ErrONguonTrong) {
		t.Fatalf("o hop le nhung rong phai tra ErrONguonTrong, tra %v", err)
	}
}

// TestChuyenOMonNamNgoaiBienBangMacBiTuChoi — DOT BIEN DB1, qua HANDLER.
//
// CSDL KHONG rang buoc `o_so < 5` cho MAC, va DungBangHanhTrang chi kiem theo
// tran chung 127 — nen mot dong hong (hoac mot task tuong lai ghi sai) co the
// dat mon o MAC o 10. Luc do bang phien CO o 10, va neu handler khong kiem bien
// rieng thi mon do chuyen di duoc. Day la truong hop bang phien KHONG che duoc.
func TestChuyenOMonNamNgoaiBienBangMacBiTuChoi(t *testing.T) {
	c := dungDoDungCanh(t)
	id := c.datMon(t, store.NoiChuaMac, 10, 3) // 10 >= SoOMacChuan
	c.napLaiBang(t)

	if _, co := c.b.hanhTrang.TraId(store.NoiChuaMac, 10); !co {
		t.Fatal("tien de test sai: bang phien phai chua MAC o 10")
	}
	// loai 5 = mac -> tui.
	if err := c.b.XuLy(c.phien, goiLenDungDo(lenhChuyenODo, 5, 10)); err != nil {
		t.Fatalf("xu ly that bai: %v", err)
	}
	noi, o := c.docNoiChuaCuaMon(t, id)
	if noi != store.NoiChuaMac || o == nil || *o != 10 {
		t.Fatalf("mon o o mac ngoai bien VAN chuyen duoc (dang o %q o=%v) — thieu kiem bien (dot bien DB1)",
			noi, o)
	}
}

// ---------------------------------------------------------------------------
// 3. MAC / THAO
// ---------------------------------------------------------------------------

// TestMacTrangBiVaoDungSlotTheoLoai: mau 3 co Loai = 1 (quan) -> phai vao o
// mac so 1, khong phai "o trong nho nhat" (= 0).
func TestMacTrangBiVaoDungSlotTheoLoai(t *testing.T) {
	c := dungDoDungCanh(t)
	id := c.datMon(t, store.NoiChuaTui, 2, 3) // mau 3: Loai 1
	c.napLaiBang(t)

	if mau := TimMauVatPham(3); mau == nil || mau.Loai != 1 {
		t.Fatalf("tien de test sai: mau 3 phai co Loai = 1")
	}
	// loai 4 = tui -> mac.
	if err := c.b.XuLy(c.phien, goiLenDungDo(lenhChuyenODo, 4, 2)); err != nil {
		t.Fatalf("xu ly that bai: %v", err)
	}
	noi, o := c.docNoiChuaCuaMon(t, id)
	if noi != store.NoiChuaMac || o == nil || *o != 1 {
		t.Fatalf("phai vao MAC o 1 (theo Loai), dang o %q o=%v", noi, o)
	}
}

// TestMacMonKhongPhaiTrangBiBiTuChoi — DOT BIEN DB3.
//
// Mau 0 co Loai = 6 (thuoc, ngoai [0,SoOMacChuan)). Mac no phai bi tu choi, va
// mon phai nam nguyen trong tui.
func TestMacMonKhongPhaiTrangBiBiTuChoi(t *testing.T) {
	c := dungDoDungCanh(t)
	id := c.datMon(t, store.NoiChuaTui, 0, 0) // mau 0: Loai 6
	c.napLaiBang(t)

	if mau := TimMauVatPham(0); mau == nil || int(mau.Loai) < SoOMacChuan {
		t.Fatalf("tien de test sai: mau 0 phai co Loai >= %d", SoOMacChuan)
	}
	if err := c.b.XuLy(c.phien, goiLenDungDo(lenhChuyenODo, 4, 0)); err != nil {
		t.Fatalf("xu ly that bai: %v", err)
	}
	noi, o := c.docNoiChuaCuaMon(t, id)
	if noi != store.NoiChuaTui || o == nil || *o != 0 {
		t.Fatalf("mon khong phai trang bi VAN duoc mac (dang o %q o=%v) — thieu kiem Loai (dot bien DB3)",
			noi, o)
	}
	if timGoi(c.docGoi(), lenhDungDoBao) == nil {
		t.Fatal("phai gui -25 giai thich")
	}
}

// TestMacKhiSlotDaCoDoBiTuChoi: kho KHONG hoan doi (vatpham_kho.go) -> phai bao
// nguoi choi thao do cu ra truoc, va KHONG duoc lam mat mon nao.
func TestMacKhiSlotDaCoDoBiTuChoi(t *testing.T) {
	c := dungDoDungCanh(t)
	idCu := c.datMon(t, store.NoiChuaMac, 1, 3)
	idMoi := c.datMon(t, store.NoiChuaTui, 0, 3)
	c.napLaiBang(t)

	if err := c.b.XuLy(c.phien, goiLenDungDo(lenhChuyenODo, 4, 0)); err != nil {
		t.Fatalf("xu ly that bai: %v", err)
	}
	if noi, o := c.docNoiChuaCuaMon(t, idCu); noi != store.NoiChuaMac || o == nil || *o != 1 {
		t.Fatalf("do cu phai con nguyen o MAC 1, dang o %q o=%v", noi, o)
	}
	if noi, o := c.docNoiChuaCuaMon(t, idMoi); noi != store.NoiChuaTui || o == nil || *o != 0 {
		t.Fatalf("do moi phai con nguyen o TUI 0, dang o %q o=%v", noi, o)
	}
}

// TestThaoTrangBiVeOTrongDauTienCuaTui: mac -> tui, vao o trong nho nhat.
func TestThaoTrangBiVeOTrongDauTienCuaTui(t *testing.T) {
	c := dungDoDungCanh(t)
	id := c.datMon(t, store.NoiChuaMac, 0, 2)
	c.datMon(t, store.NoiChuaTui, 0, 3) // chiem o 0 cua tui
	c.napLaiBang(t)

	// loai 5 = mac -> tui.
	if err := c.b.XuLy(c.phien, goiLenDungDo(lenhChuyenODo, 5, 0)); err != nil {
		t.Fatalf("xu ly that bai: %v", err)
	}
	noi, o := c.docNoiChuaCuaMon(t, id)
	if noi != store.NoiChuaTui || o == nil || *o != 1 {
		t.Fatalf("thao do phai ve TUI o 1, dang o %q o=%v", noi, o)
	}
}

// TestThaoTrangBiKhiTuiDayBiTuChoi: tui day -> tu choi, do van dang mac.
func TestThaoTrangBiKhiTuiDayBiTuChoi(t *testing.T) {
	c := dungDoDungCanh(t)
	id := c.datMon(t, store.NoiChuaMac, 0, 2)
	for o := 0; o < SoOTuiChuan; o++ {
		c.datMon(t, store.NoiChuaTui, int16(o), 3)
	}
	c.napLaiBang(t)

	if err := c.b.XuLy(c.phien, goiLenDungDo(lenhChuyenODo, 5, 0)); err != nil {
		t.Fatalf("xu ly that bai: %v", err)
	}
	noi, o := c.docNoiChuaCuaMon(t, id)
	if noi != store.NoiChuaMac || o == nil || *o != 0 {
		t.Fatalf("tui day: do phai con dang mac, dang o %q o=%v", noi, o)
	}
	if timGoi(c.docGoi(), lenhDungDoBao) == nil {
		t.Fatal("phai gui -25 bao tui day")
	}
}

// ---------------------------------------------------------------------------
// 4. Goi cut / goi la — khong duoc ngat ket noi
// ---------------------------------------------------------------------------

// TestGoiCutKhongNgatKetNoi: goi thieu byte o CA HAI ma lenh phai duoc bo qua
// an toan (khuon xuLyDiChuyen).
func TestGoiCutKhongNgatKetNoi(t *testing.T) {
	c := dungDoDungCanh(t)

	goiCut := []*protocol.Message{
		goiLenDungDo(lenhDungBoDo),
		goiLenDungDo(lenhDungBoDo, 2),
		goiLenDungDo(lenhDungBoDo, 2, 1),
		goiLenDungDo(lenhChuyenODo),
		goiLenDungDo(lenhChuyenODo, 1),
	}
	for i, g := range goiCut {
		if err := c.b.XuLy(c.phien, g); err != nil {
			t.Fatalf("goi cut #%d lam XuLy tra loi (se ngat ket noi): %v", i, err)
		}
	}
}

// TestLoaiChuyenOLaBiTuChoi: `loai` 6..9 (de tu / thu cung) va 2 (khong ton
// tai) phai bi tu choi co cau tra loi.
func TestLoaiChuyenOLaBiTuChoi(t *testing.T) {
	c := dungDoDungCanh(t)
	id := c.datMon(t, store.NoiChuaTui, 0, 2)
	c.napLaiBang(t)

	for _, loai := range []byte{2, 6, 7, 8, 9, 50} {
		if err := c.b.XuLy(c.phien, goiLenDungDo(lenhChuyenODo, loai, 0)); err != nil {
			t.Fatalf("loai=%d: xu ly that bai: %v", loai, err)
		}
		if noi, _ := c.docNoiChuaCuaMon(t, id); noi != store.NoiChuaTui {
			t.Fatalf("loai=%d: mon bi dong toi, dang o %q", loai, noi)
		}
	}
}
