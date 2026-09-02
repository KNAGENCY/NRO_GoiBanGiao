// cuahang_test.go 鈥?test BYTE-CHINH-XAC cho goi -44 (loaiShop=0).
//
// Byte mong doi duoc GHIM DOC LAP: dung tay bang utfCH/short/int cua chinh
// file test, KHONG goi lai protocol.Writer cua ban cai. Neu bo dung sai thu tu
// hoac sot mot truong, so byte lech ngay.
package game

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
)

// --- do nghe ghim byte doc lap ------------------------------------------------

func utfCH(s string) []byte {
	b := []byte(s) // moi chuoi trong test la ASCII -> modified UTF-8 == UTF-8
	out := []byte{byte(len(b) >> 8), byte(len(b))}
	return append(out, b...)
}

func shortCH(v int16) []byte {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], uint16(v))
	return b[:]
}

func intCH(v int32) []byte {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(v))
	return b[:]
}

func noiCH(phan ...[]byte) []byte {
	var out []byte
	for _, p := range phan {
		out = append(out, p...)
	}
	return out
}

// --- ca 1: cua hang rong ------------------------------------------------------

func TestGoiCuaHangRong(t *testing.T) {
	m, err := goiCuaHang(&CuaHang{ID: 9, Ten: "rong"})
	if err != nil {
		t.Fatalf("cua hang rong phai dung duoc: %v", err)
	}
	if m.Command != -44 {
		t.Errorf("Command = %d, mong doi -44", m.Command)
	}
	want := []byte{0, 0} // loaiShop=0, soTab=0
	if !bytes.Equal(m.Data, want) {
		t.Errorf("than goi = % x, mong doi % x", m.Data, want)
	}
}

// --- ca 2: mot tab mot mon, gia vang -----------------------------------------

func TestGoiCuaHangMotMonVang(t *testing.T) {
	ch := &CuaHang{ID: 1, Ten: "x", Tabs: []TabCuaHang{{
		Ten: "Thuoc",
		Mon: []MonHang{{MauID: 0, Gia: 100, Tien: TienVang}},
	}}}
	m, err := goiCuaHang(ch)
	if err != nil {
		t.Fatalf("khong dung duoc: %v", err)
	}
	want := noiCH(
		[]byte{0},      // loaiShop
		[]byte{1},      // soTab
		utfCH("Thuoc"), // tenTab
		[]byte{1},      // soMon
		shortCH(0),     // mauID
		intCH(100),     // giaVang
		intCH(0),       // giaNgoc
		[]byte{0},      // soOption
		[]byte{0},      // laMoi
		[]byte{0},      // coPart
	)
	if !bytes.Equal(m.Data, want) {
		t.Errorf("than goi = % x\nmong doi   = % x", m.Data, want)
	}
}

// --- ca 3: nhieu tab, gia ngoc, option, laMoi --------------------------------

func TestGoiCuaHangNhieuTabDayDu(t *testing.T) {
	ch := &CuaHang{ID: 2, Ten: "x", Tabs: []TabCuaHang{
		{Ten: "A", Mon: []MonHang{{MauID: 1, Gia: 7, Tien: TienNgoc, LaMoi: true}}},
		{Ten: "B", Mon: []MonHang{
			{MauID: 2, Gia: 1000, Tien: TienVang, Options: []OptionMon{
				{OptionID: 6, Param: 500},
				{OptionID: 7, Param: 65535},
			}},
			{MauID: 3, Gia: 0, Tien: TienVang},
		}},
	}}
	m, err := goiCuaHang(ch)
	if err != nil {
		t.Fatalf("khong dung duoc: %v", err)
	}
	want := noiCH(
		[]byte{0}, []byte{2},
		// tab A
		utfCH("A"), []byte{1},
		shortCH(1), intCH(0), intCH(7),
		[]byte{0}, // soOption
		[]byte{1}, // laMoi
		[]byte{0}, // coPart
		// tab B
		utfCH("B"), []byte{2},
		shortCH(2), intCH(1000), intCH(0),
		[]byte{2},               // soOption
		[]byte{6}, shortCH(500), // option 1
		[]byte{7}, shortCH(int16(-1)), // option 2: param 65535 -> 0xFFFF
		[]byte{0}, []byte{0},
		shortCH(3), intCH(0), intCH(0),
		[]byte{0}, []byte{0}, []byte{0},
	)
	if !bytes.Equal(m.Data, want) {
		t.Errorf("than goi = % x\nmong doi   = % x", m.Data, want)
	}
}

// --- ca 4: doc NGUOC theo dung trinh tu client -------------------------------

// docNguocCuaHang mo phong dung trinh tu Controller.cs case -44 (loaiShop 0).
// Doc het ma con byte thua, hoac thieu byte, deu bao loi.
func docNguocCuaHang(t *testing.T, than []byte) []TabCuaHang {
	t.Helper()
	p := 0
	lay := func(n int) []byte {
		if p+n > len(than) {
			t.Fatalf("thieu byte tai vi tri %d (can %d, con %d)", p, n, len(than)-p)
		}
		b := than[p : p+n]
		p += n
		return b
	}
	docByte := func() byte { return lay(1)[0] }
	docShort := func() int16 { return int16(binary.BigEndian.Uint16(lay(2))) }
	docInt := func() int32 { return int32(binary.BigEndian.Uint32(lay(4))) }
	docUTF := func() string {
		n := int(binary.BigEndian.Uint16(lay(2)))
		return string(lay(n))
	}

	if lt := docByte(); lt != 0 {
		t.Fatalf("loaiShop = %d, mong doi 0", lt)
	}
	soTab := int(docByte())
	tabs := make([]TabCuaHang, 0, soTab)
	for i := 0; i < soTab; i++ {
		tab := TabCuaHang{Ten: docUTF()}
		soMon := int(docByte())
		for j := 0; j < soMon; j++ {
			mon := MonHang{MauID: docShort()}
			vang := docInt()
			ngoc := docInt()
			if vang != 0 && ngoc != 0 {
				t.Fatalf("tab %d mon %d: ca hai o gia deu khac 0", i, j)
			}
			if ngoc != 0 {
				mon.Tien, mon.Gia = TienNgoc, ngoc
			} else {
				mon.Tien, mon.Gia = TienVang, vang
			}
			soOpt := int(docByte())
			for k := 0; k < soOpt; k++ {
				id := docByte()
				param := uint16(docShort())
				mon.Options = append(mon.Options, OptionMon{OptionID: id, Param: param})
			}
			mon.LaMoi = docByte() != 0
			if docByte() == 1 {
				_ = docShort()
				_ = docShort()
				_ = docShort()
				_ = docShort()
			}
			tab.Mon = append(tab.Mon, mon)
		}
		tabs = append(tabs, tab)
	}
	if p != len(than) {
		t.Fatalf("con thua %d byte sau khi doc xong", len(than)-p)
	}
	return tabs
}

func TestDocNguocCuaHangMau(t *testing.T) {
	ch := TimCuaHang(1)
	if ch == nil {
		t.Fatal("khong tim thay cua hang mau id=1")
	}
	m, err := goiCuaHang(ch)
	if err != nil {
		t.Fatalf("khong dung duoc: %v", err)
	}
	tabs := docNguocCuaHang(t, m.Data)
	if len(tabs) != len(ch.Tabs) {
		t.Fatalf("doc nguoc ra %d tab, mong doi %d", len(tabs), len(ch.Tabs))
	}
	for i, tab := range tabs {
		if tab.Ten != ch.Tabs[i].Ten {
			t.Errorf("tab %d ten = %q, mong doi %q", i, tab.Ten, ch.Tabs[i].Ten)
		}
		if len(tab.Mon) != len(ch.Tabs[i].Mon) {
			t.Fatalf("tab %d co %d mon, mong doi %d", i, len(tab.Mon), len(ch.Tabs[i].Mon))
		}
		for j, mon := range tab.Mon {
			goc := ch.Tabs[i].Mon[j]
			if mon.MauID != goc.MauID || mon.Gia != goc.Gia || mon.Tien != goc.Tien ||
				mon.LaMoi != goc.LaMoi {
				t.Errorf("tab %d mon %d = %+v, mong doi %+v", i, j, mon, goc)
			}
		}
	}
}

// --- ca 4b: nhanh coPart (mau Loai == 5) -------------------------------------

// Danh muc Naruto hien CHUA co mau Loai == 5 nen nhanh nay khong co du lieu that
// de thu. Test tu them mot mau tam vao CUOI danh muc (dung luat "chi them vao
// cuoi") roi tra lai nguyen trang — de nhanh coPart khong bao gio thanh code chet.
func TestGoiCuaHangCoPart(t *testing.T) {
	goc := danhSachMauVatPham
	t.Cleanup(func() { danhSachMauVatPham = goc })
	danhSachMauVatPham = append(append([]MauVatPham(nil), goc...), MauVatPham{
		Loai: 5, GioiTinh: 3, Ten: "bo do tam", MoTa: "chi dung trong test",
		Cap: 1, IconID: 0, Phan: 42, TranChong: 1,
	})
	idMoi := int16(len(danhSachMauVatPham) - 1)

	ch := &CuaHang{ID: 3, Ten: "x", Tabs: []TabCuaHang{{
		Ten: "A", Mon: []MonHang{{MauID: idMoi, Gia: 1, Tien: TienVang}},
	}}}
	m, err := goiCuaHang(ch)
	if err != nil {
		t.Fatalf("khong dung duoc: %v", err)
	}
	want := noiCH(
		[]byte{0}, []byte{1},
		utfCH("A"), []byte{1},
		shortCH(idMoi), intCH(1), intCH(0),
		[]byte{0}, // soOption
		[]byte{0}, // laMoi
		[]byte{1}, // coPart = 1 vi Loai == 5
		shortCH(42), shortCH(42), shortCH(42), shortCH(-1),
	)
	if !bytes.Equal(m.Data, want) {
		t.Errorf("than goi = % x\nmong doi   = % x", m.Data, want)
	}
	// doc nguoc phai tieu thu het byte (khong thua, khong thieu)
	docNguocCuaHang(t, m.Data)
}

// --- ca 5: cac loi phai chan --------------------------------------------------

func TestGoiCuaHangLoi(t *testing.T) {
	cases := []struct {
		ten   string
		ch    CuaHang
		chuoi string
	}{
		{
			ten:   "mau vat pham khong ton tai",
			ch:    CuaHang{Tabs: []TabCuaHang{{Ten: "A", Mon: []MonHang{{MauID: 9999}}}}},
			chuoi: "KHONG TON TAI",
		},
		{
			ten:   "mau vat pham am",
			ch:    CuaHang{Tabs: []TabCuaHang{{Ten: "A", Mon: []MonHang{{MauID: -1}}}}},
			chuoi: "KHONG TON TAI",
		},
		{
			ten:   "gia am",
			ch:    CuaHang{Tabs: []TabCuaHang{{Ten: "A", Mon: []MonHang{{MauID: 0, Gia: -1}}}}},
			chuoi: "gia am",
		},
		{
			ten: "option id vuot bang mau",
			ch: CuaHang{Tabs: []TabCuaHang{{Ten: "A", Mon: []MonHang{
				{MauID: 0, Options: []OptionMon{{OptionID: 200}}}}}}},
			chuoi: "mau option",
		},
		{
			ten:   "ten tab rong",
			ch:    CuaHang{Tabs: []TabCuaHang{{Ten: "", Mon: nil}}},
			chuoi: "ten tab rong",
		},
		{
			ten:   "ten tab 3 dong",
			ch:    CuaHang{Tabs: []TabCuaHang{{Ten: "a\nb\nc"}}},
			chuoi: "toi da 2",
		},
		{
			ten:   "loai tien la",
			ch:    CuaHang{Tabs: []TabCuaHang{{Ten: "A", Mon: []MonHang{{MauID: 0, Tien: LoaiTien(9)}}}}},
			chuoi: "loai tien",
		},
	}
	for _, c := range cases {
		t.Run(c.ten, func(t *testing.T) {
			m, err := goiCuaHang(&c.ch)
			if err == nil {
				t.Fatalf("phai bao loi, nhung dung duoc goi % x", m.Data)
			}
			if !strings.Contains(err.Error(), c.chuoi) {
				t.Errorf("loi = %q, mong doi chua %q", err, c.chuoi)
			}
		})
	}
}

// --- ca 6: tran so tab / so mon ----------------------------------------------

func TestGoiCuaHangTranSoLuong(t *testing.T) {
	// tran so tab
	ch := CuaHang{}
	for i := 0; i <= tranSoTab; i++ {
		ch.Tabs = append(ch.Tabs, TabCuaHang{Ten: "A"})
	}
	if _, err := goiCuaHang(&ch); err == nil {
		t.Errorf("%d tab phai bao loi (tran %d)", len(ch.Tabs), tranSoTab)
	}
	// dung tran thi phai dung duoc
	ch.Tabs = ch.Tabs[:tranSoTab]
	if _, err := goiCuaHang(&ch); err != nil {
		t.Errorf("%d tab (dung tran) phai dung duoc: %v", tranSoTab, err)
	}

	// tran so mon mot tab
	tab := TabCuaHang{Ten: "A"}
	for i := 0; i <= tranSoMonMoiTab; i++ {
		tab.Mon = append(tab.Mon, MonHang{MauID: 0, Gia: 1})
	}
	ch2 := CuaHang{Tabs: []TabCuaHang{tab}}
	if _, err := goiCuaHang(&ch2); err == nil {
		t.Errorf("%d mon phai bao loi (tran %d)", len(tab.Mon), tranSoMonMoiTab)
	}
}

// --- ca 7: danh muc cua hang thuc te khong co MauID lac ----------------------

func TestDanhSachCuaHangHopLe(t *testing.T) {
	if err := kiemMoiCuaHang(); err != nil {
		t.Fatalf("danh muc cua hang co loi: %v", err)
	}
	if len(DanhSachCuaHang()) == 0 {
		t.Fatal("phai co it nhat mot cua hang mau de nghiem thu dau-cuoi")
	}
	for _, ch := range DanhSachCuaHang() {
		if _, err := goiCuaHang(&ch); err != nil {
			t.Errorf("cua hang %d khong dung duoc goi: %v", ch.ID, err)
		}
	}
}
