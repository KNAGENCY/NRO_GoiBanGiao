package game

import (
	"encoding/hex"
	"math"
	"strings"
	"testing"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// TestGoiXuatHien kiem tra goi -5 xuat hien: doc lai DU 28 truong dung
// thu tu/kieu, khang dinh het buffer sau truong 28 (hat).
func TestGoiXuatHien(t *testing.T) {
	nv := &store.NhanVat{
		ID:      42,
		Ten:     "Naruto",
		Lang:    0,
		KieuToc: 7,
		Cap:     5,
		HpToiDa: 1000,
		ViTriX:  100,
		ViTriY:  384,
	}

	m, err := goiXuatHien(nv)
	if err != nil {
		t.Fatalf("goiXuatHien loi khong mong doi: %v", err)
	}

	if m.Command != -5 {
		t.Errorf("expect command -5, got %d", m.Command)
	}

	r := protocol.NewReader(m.Data)

	// 1 playerId
	if v, err := r.Int(); err != nil || v != 42 {
		t.Fatalf("1 playerId: got %d, err=%v", v, err)
	}
	// 2 clanId
	if v, err := r.Int(); err != nil || v != -1 {
		t.Fatalf("2 clanId: got %d, err=%v", v, err)
	}
	// 3 level
	if v, err := r.Byte(); err != nil || v != 5 {
		t.Fatalf("3 level: got %d, err=%v", v, err)
	}
	// 4 co co dinh (writeBoolean false -> doc lai la byte 0)
	if v, err := r.Byte(); err != nil || v != 0 {
		t.Fatalf("4 co co dinh: got %d, err=%v", v, err)
	}
	// 5 typePk
	if v, err := r.Byte(); err != nil || v != 0 {
		t.Fatalf("5 typePk: got %d, err=%v", v, err)
	}
	// 6 gender
	if v, err := r.Byte(); err != nil || v != 0 {
		t.Fatalf("6 gender: got %d, err=%v", v, err)
	}
	// 7 gender lan 2
	if v, err := r.Byte(); err != nil || v != 0 {
		t.Fatalf("7 gender lan 2: got %d, err=%v", v, err)
	}
	// 8 head
	if v, err := r.Short(); err != nil || v != 7 {
		t.Fatalf("8 head: got %d, err=%v", v, err)
	}
	// 9 name
	if v, err := r.UTF(); err != nil || v != "Naruto" {
		t.Fatalf("9 name: got %q, err=%v", v, err)
	}
	// 10 hp
	if v, err := r.Int(); err != nil || v != 1000 {
		t.Fatalf("10 hp: got %d, err=%v", v, err)
	}
	// 11 hpMax
	if v, err := r.Int(); err != nil || v != 1000 {
		t.Fatalf("11 hpMax: got %d, err=%v", v, err)
	}
	// 12 body
	if v, err := r.Short(); err != nil || v != 0 {
		t.Fatalf("12 body: got %d, err=%v", v, err)
	}
	// 13 leg
	if v, err := r.Short(); err != nil || v != 0 {
		t.Fatalf("13 leg: got %d, err=%v", v, err)
	}
	// 14 flagbag
	if v, err := r.Byte(); err != nil || v != 0 {
		t.Fatalf("14 flagbag: got %d, err=%v", v, err)
	}
	// 15 hang -1 (0xFF)
	if v, err := r.Byte(); err != nil || v != 0xFF {
		t.Fatalf("15 hang -1: got %d (0x%02X), err=%v", v, v, err)
	}
	// 16 x
	if v, err := r.Short(); err != nil || v != 100 {
		t.Fatalf("16 x: got %d, err=%v", v, err)
	}
	// 17 y
	if v, err := r.Short(); err != nil || v != 384 {
		t.Fatalf("17 y: got %d, err=%v", v, err)
	}
	// 18 effBuffHp
	if v, err := r.Short(); err != nil || v != 0 {
		t.Fatalf("18 effBuffHp: got %d, err=%v", v, err)
	}
	// 19 effBuffMp
	if v, err := r.Short(); err != nil || v != 0 {
		t.Fatalf("19 effBuffMp: got %d, err=%v", v, err)
	}
	// 20 numEff
	if v, err := r.Byte(); err != nil || v != 0 {
		t.Fatalf("20 numEff: got %d, err=%v", v, err)
	}
	// 21 idSpaceShip
	if v, err := r.Byte(); err != nil || v != 0 {
		t.Fatalf("21 idSpaceShip: got %d, err=%v", v, err)
	}
	// 22 isMonkey
	if v, err := r.Byte(); err != nil || v != 0 {
		t.Fatalf("22 isMonkey: got %d, err=%v", v, err)
	}
	// 23 mount
	if v, err := r.Short(); err != nil || v != 0 {
		t.Fatalf("23 mount: got %d, err=%v", v, err)
	}
	// 24 cFlag
	if v, err := r.Byte(); err != nil || v != 0 {
		t.Fatalf("24 cFlag: got %d, err=%v", v, err)
	}
	// 25 hang 0
	if v, err := r.Byte(); err != nil || v != 0 {
		t.Fatalf("25 hang 0: got %d, err=%v", v, err)
	}
	// 26 aura
	if v, err := r.Short(); err != nil || v != 0 {
		t.Fatalf("26 aura: got %d, err=%v", v, err)
	}
	// 27 effFront
	if v, err := r.Byte(); err != nil || v != 0 {
		t.Fatalf("27 effFront: got %d, err=%v", v, err)
	}
	// 28 hat
	if v, err := r.Short(); err != nil || v != 0 {
		t.Fatalf("28 hat: got %d, err=%v", v, err)
	}

	// Het buffer: doc them phai loi.
	if _, err := r.Byte(); err == nil {
		t.Fatal("sau truong 28 (hat) buffer phai het, nhung doc them khong loi")
	}
}

// nvVang la bo mau nhan vat kem NGUYEN VAN byte than goi -5 ma ban code TRUOC
// khi tach ham (Task 4 KH6b) sinh ra. Day la "vang" (golden) chup tu chinh code
// cu chu KHONG go tay: lech mot bit la goi -5 doi tren day, ma goi -5 la thu
// client Unity doc de VE nguoi choi — khong co cho thuong luong.
//
// Bo mau co y chon gia tri BIEN: ten tieng Viet co dau (modified UTF-8 nhieu
// byte), hp = 0, hp = MaxInt32, x/y am, cap 1 va cap 255, id vuot int32.
var nvVang = []struct {
	ten  string
	nv   *store.NhanVat
	vang string
}{
	{
		"co ban",
		&store.NhanVat{ID: 42, Ten: "Naruto", Lang: 0, KieuToc: 7, Cap: 5, HpToiDa: 1000, HpHienTai: 300, ViTriX: 100, ViTriY: 384},
		"0000002affffffff0500000000000700064e617275746f000003e8000003e80000000000ff0064018000000000000000000000000000000000",
	},
	{
		"ten tieng Viet co dau, hp 0, x/y am, cap 1",
		&store.NhanVat{ID: 1, Ten: "Hắc Ảnh Đệ Nhất", Lang: 2, KieuToc: 300, Cap: 1, HpToiDa: 0, HpHienTai: 0, ViTriX: -1, ViTriY: -32768},
		"00000001ffffffff0100000202012c001848e1baaf6320e1baa26e6820c490e1bb87204e68e1baa57400000000000000000000000000ffffff800000000000000000000000000000000000",
	},
	{
		"cap 255, hp MaxInt32",
		&store.NhanVat{ID: 2147483647, Ten: "Z", Lang: 4, KieuToc: -1, Cap: 255, HpToiDa: math.MaxInt32, HpHienTai: 7, ViTriX: 32767, ViTriY: -1},
		"7fffffffffffffffff00000404ffff00015a7fffffff7fffffff0000000000ff7fffffff00000000000000000000000000000000",
	},
	{
		"id vuot int32, ten rong, cap 256 (tran byte)",
		&store.NhanVat{ID: 9000000000, Ten: "", Lang: 255, KieuToc: 0, Cap: 256, HpToiDa: -1, HpHienTai: -1, ViTriX: 70000, ViTriY: 70000},
		"18711a00ffffffff000000ffff00000000ffffffffffffffff0000000000ff1170117000000000000000000000000000000000",
	},
}

// TestGoiXuatHienByteVang chot rang goi -5 gui cho NGUOI CHOI khong doi mot bit
// nao sau khi tach ham: so sanh BYTE-BY-BYTE voi vang chup truoc refactor.
func TestGoiXuatHienByteVang(t *testing.T) {
	for _, c := range nvVang {
		t.Run(c.ten, func(t *testing.T) {
			m, err := goiXuatHien(c.nv)
			if err != nil {
				t.Fatalf("goiXuatHien loi khong mong doi: %v", err)
			}
			if m.Command != lenhXuatHien {
				t.Errorf("command: got %d, want %d", m.Command, lenhXuatHien)
			}
			if got := hex.EncodeToString(m.Data); got != c.vang {
				t.Errorf("than goi -5 LECH so voi vang:\ngot  %s\nwant %s", got, c.vang)
			}
		})
	}
}

// TestGoiXuatHienDayDuKhopVoBoc chot rang vo boc goiXuatHien khong lam gi khac
// ngoai viec dien tham so tu *store.NhanVat: goi THANG ham day du voi dung cac
// gia tri do phai ra Y HET byte vang.
func TestGoiXuatHienDayDuKhopVoBoc(t *testing.T) {
	for _, c := range nvVang {
		t.Run(c.ten, func(t *testing.T) {
			m := goiXuatHienDayDu(
				c.nv.ID, c.nv.Ten, byte(c.nv.Cap),
				c.nv.HpToiDa, c.nv.HpToiDa,
				int16(c.nv.ViTriX), int16(c.nv.ViTriY),
				HinhDangHien{Lang: c.nv.Lang, Outfit: [6]int16{c.nv.KieuToc, 0, 0, 0, 0, 0}},
			)
			if m == nil {
				t.Fatal("goiXuatHienDayDu tra nil")
			}
			if m.Command != lenhXuatHien {
				t.Errorf("command: got %d, want %d", m.Command, lenhXuatHien)
			}
			if got := hex.EncodeToString(m.Data); got != c.vang {
				t.Errorf("than goi -5 LECH so voi vang:\ngot  %s\nwant %s", got, c.vang)
			}
		})
	}
}

// TestGoiXuatHienDayDuOutfitThat kiem tung o cua Outfit roi DUNG cho tren day:
// head, body, leg, flagBag, aura, effFront. Sai thu tu la boss mac nham do —
// mac duoc do that chinh la ly do ham day du ton tai. Moi o mot con so KHAC
// NHAU de nham o nao cung lo, va hp KHAC hpToiDa de hai truong khong lan nhau.
func TestGoiXuatHienDayDuOutfitThat(t *testing.T) {
	hd := HinhDangHien{
		Lang:   4,
		Outfit: [6]int16{11, 22, 33, 44, 55, 66}, // head, body, leg, flagBag, aura, effFront
		TypePk: 3,
	}
	m := goiXuatHienDayDu(-9, "Cửu Vĩ", 90, 1234, 5678, -300, 777, hd)
	if m == nil {
		t.Fatal("goiXuatHienDayDu tra nil")
	}
	if m.Command != lenhXuatHien {
		t.Errorf("command: got %d, want %d", m.Command, lenhXuatHien)
	}

	r := protocol.NewReader(m.Data)
	docInt := func(so int, ten string, muon int32) {
		t.Helper()
		v, err := r.Int()
		if err != nil || v != muon {
			t.Fatalf("%d %s: got %d, want %d, err=%v", so, ten, v, muon, err)
		}
	}
	docByte := func(so int, ten string, muon byte) {
		t.Helper()
		v, err := r.Byte()
		if err != nil || v != muon {
			t.Fatalf("%d %s: got %d, want %d, err=%v", so, ten, v, muon, err)
		}
	}
	docShort := func(so int, ten string, muon int16) {
		t.Helper()
		v, err := r.Short()
		if err != nil || v != muon {
			t.Fatalf("%d %s: got %d, want %d, err=%v", so, ten, v, muon, err)
		}
	}

	docInt(1, "playerId", -9)
	docInt(2, "clanId", -1)
	docByte(3, "level", 90)
	docByte(4, "co co dinh", 0)
	docByte(5, "typePk", 3)
	docByte(6, "gender", 4)
	docByte(7, "gender lan 2", 4)
	docShort(8, "head", 11)
	if v, err := r.UTF(); err != nil || v != "Cửu Vĩ" {
		t.Fatalf("9 name: got %q, err=%v", v, err)
	}
	docInt(10, "hp", 1234)
	docInt(11, "hpMax", 5678)
	docShort(12, "body", 22)
	docShort(13, "leg", 33)
	docByte(14, "flagbag", 44)
	docByte(15, "hang -1", 0xFF)
	docShort(16, "x", -300)
	docShort(17, "y", 777)
	docShort(18, "effBuffHp", 0)
	docShort(19, "effBuffMp", 0)
	docByte(20, "numEff", 0)
	docByte(21, "idSpaceShip", 0)
	docByte(22, "isMonkey", 0)
	docShort(23, "mount", 0)
	docByte(24, "cFlag", 0)
	docByte(25, "hang 0", 0)
	docShort(26, "aura", 55)
	docByte(27, "effFront", 66)
	docShort(28, "hat", 0)

	if _, err := r.Byte(); err == nil {
		t.Fatal("sau truong 28 (hat) buffer phai het, nhung doc them khong loi")
	}
}

// TestGoiXuatHienDayDuFlagBagEffFrontChiMotByte ghi thanh test cai bay: hai o
// flagBag va effFront la int16 trong HinhDangHien nhung tren day CHI MOT BYTE.
// Gia tri > 255 bi cat con byte thap — ai dat so lon phai biet truoc, dung de
// den luc boss mac nham do moi phat hien.
func TestGoiXuatHienDayDuFlagBagEffFrontChiMotByte(t *testing.T) {
	hd := HinhDangHien{Outfit: [6]int16{0, 0, 0, 0x1234, 0, 0x0102}} // flagBag, effFront
	m := goiXuatHienDayDu(1, "x", 1, 1, 1, 0, 0, hd)
	if m == nil {
		t.Fatal("goiXuatHienDayDu tra nil")
	}
	r := protocol.NewReader(m.Data)
	for i := 0; i < 13; i++ { // bo qua truong 1..13 bang cach doc lai dung kieu
		switch i {
		case 0, 1: // 1 playerId, 2 clanId
			if _, err := r.Int(); err != nil {
				t.Fatalf("bo qua truong %d: %v", i+1, err)
			}
		case 2, 3, 4, 5, 6: // 3 level, 4 co, 5 typePk, 6+7 gender
			if _, err := r.Byte(); err != nil {
				t.Fatalf("bo qua truong %d: %v", i+1, err)
			}
		case 7: // 8 head
			if _, err := r.Short(); err != nil {
				t.Fatalf("bo qua truong %d: %v", i+1, err)
			}
		case 8: // 9 name
			if _, err := r.UTF(); err != nil {
				t.Fatalf("bo qua truong %d: %v", i+1, err)
			}
		case 9, 10: // 10 hp, 11 hpMax
			if _, err := r.Int(); err != nil {
				t.Fatalf("bo qua truong %d: %v", i+1, err)
			}
		case 11, 12: // 12 body, 13 leg
			if _, err := r.Short(); err != nil {
				t.Fatalf("bo qua truong %d: %v", i+1, err)
			}
		}
	}
	if v, err := r.Byte(); err != nil || v != 0x34 {
		t.Fatalf("14 flagbag: got 0x%02X, want 0x34 (byte thap cua 0x1234), err=%v", v, err)
	}
}

// TestGoiXuatHienDayDuTenQuaDai chot nhanh loi: ten dai qua muc writeUTF chiu
// duoc (65535 byte) thi ham day du tra nil, va vo boc bien nil do thanh loi —
// dung nhu ban cu (khuvuc.go log roi bo qua, KHONG gui goi hong).
func TestGoiXuatHienDayDuTenQuaDai(t *testing.T) {
	tenDai := strings.Repeat("a", 70000)
	if m := goiXuatHienDayDu(1, tenDai, 1, 1, 1, 0, 0, HinhDangHien{}); m != nil {
		t.Fatalf("ten %d byte: goiXuatHienDayDu phai tra nil, got %d byte than", len(tenDai), len(m.Data))
	}
	m, err := goiXuatHien(&store.NhanVat{ID: 1, Ten: tenDai})
	if err == nil {
		t.Fatal("ten qua dai: goiXuatHien phai bao loi")
	}
	if m != nil {
		t.Fatal("ten qua dai: goiXuatHien phai tra goi nil")
	}
}
