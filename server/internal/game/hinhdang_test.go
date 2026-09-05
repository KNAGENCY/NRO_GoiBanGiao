package game

import (
	"testing"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// TestDauHopLe: bo hinh Naruto chi co dau 30; moi gia tri khac (toc Dragon
// Ball tu man tao nhan vat 64/31/9/29/32/6/27/28, hang cu, rac) quy ve 30.
func TestDauHopLe(t *testing.T) {
	if got := dauHopLe(HeadNhanVat); got != HeadNhanVat {
		t.Fatalf("dau 30 phai giu nguyen, got %d", got)
	}
	for _, x := range []int16{64, 31, 9, 29, 32, 6, 27, 28, 0, 5, 7, -1, 300, 2189} {
		if got := dauHopLe(x); got != HeadNhanVat {
			t.Errorf("dauHopLe(%d) = %d, mong doi %d (dau Naruto duy nhat)", x, got, HeadNhanVat)
		}
	}
}

// TestGioiTinhClientLuon0: byte cgender phia client luon 0 cho MOI lang — client
// chi co bang 3 hanh tinh, bo do Naruto nam o slot mac dinh cua gender 0.
func TestGioiTinhClientLuon0(t *testing.T) {
	for lang := int16(0); lang <= 4; lang++ {
		if g := gioiTinhClient(lang); g != 0 {
			t.Errorf("gioiTinhClient(%d) = %d, mong doi 0", lang, g)
		}
	}
	if g := gioiTinhClient(255); g != 0 {
		t.Errorf("gioiTinhClient(255) = %d, mong doi 0", g)
	}
}

// TestHinhDangNguoiChoi: than/chan la bo hoc vien 57/58, dau qua dauHopLe, ba o
// con lai 0, Lang giu nguyen de byte lop/lang cua -5 van dung.
func TestHinhDangNguoiChoi(t *testing.T) {
	hd := hinhDangNguoiChoi(&store.NhanVat{Lang: 3, KieuToc: 64})
	muon := HinhDangHien{Lang: 3, Outfit: [6]int16{30, 57, 58, 0, 0, 0}}
	if hd != muon {
		t.Fatalf("hinhDangNguoiChoi = %+v, mong doi %+v", hd, muon)
	}
	if ThanNhanVat != 57 || ChanNhanVat != 58 {
		t.Fatalf("ThanNhanVat/ChanNhanVat = %d/%d, mong doi 57/58 (mapping-id.json parts[])", ThanNhanVat, ChanNhanVat)
	}
}

// TestGoiXuatHienNguoiChoiCoThanChan: goi -5 cua nguoi choi ghi head 30, body 57,
// leg 58 va byte 7 (cgender) = 0 du lang la 4 — byte 6 (lop/lang) van la 4.
func TestGoiXuatHienNguoiChoiCoThanChan(t *testing.T) {
	nv := &store.NhanVat{ID: 9, Ten: "Gaara", Lang: 4, KieuToc: 9, Cap: 1, HpToiDa: 180, ViTriX: 100, ViTriY: 384}
	m, err := goiXuatHien(nv)
	if err != nil {
		t.Fatal(err)
	}
	r := protocol.NewReader(m.Data)
	r.Int()  // id
	r.Int()  // clan
	r.Byte() // level
	r.Byte() // co dinh
	r.Byte() // typePk
	if v, _ := r.Byte(); v != 4 {
		t.Fatalf("byte 6 (lop/lang) = %d, mong doi 4", v)
	}
	if v, _ := r.Byte(); v != 0 {
		t.Fatalf("byte 7 (cgender) = %d, mong doi 0", v)
	}
	if v, _ := r.Short(); v != 30 {
		t.Fatalf("head = %d, mong doi 30 (KieuToc 9 la toc Dragon Ball)", v)
	}
	r.UTF()
	r.Int()
	r.Int()
	if v, _ := r.Short(); v != 57 {
		t.Fatalf("body = %d, mong doi 57", v)
	}
	if v, _ := r.Short(); v != 58 {
		t.Fatalf("leg = %d, mong doi 58", v)
	}
}

// TestMeInfoGioiTinh0VaDauHopLe: ME_INFO ghi byte cgender = 0 va head qua
// dauHopLe, nhung byte lop/lang (chi so nClasss[]) VAN la nv.Lang.
func TestMeInfoGioiTinh0VaDauHopLe(t *testing.T) {
	nv := &store.NhanVat{ID: 1, Ten: "Cat", Lang: 1, KieuToc: 29}
	m, err := goiMeInfo(nv, nil)
	if err != nil {
		t.Fatal(err)
	}
	r := protocol.NewReader(m.Data)
	r.Byte() // lenh con
	r.Int()  // id
	r.Byte() // taskMain
	if v, _ := r.Byte(); v != 0 {
		t.Fatalf("cgender = %d, mong doi 0 (lang 1 ma de 1 thi client mac do Namek 59/60)", v)
	}
	if v, _ := r.Short(); v != 30 {
		t.Fatalf("head = %d, mong doi 30", v)
	}
	r.UTF()
	r.Byte()  // cPK
	r.Byte()  // typePk
	r.Long()  // power
	r.Short() // buffHp
	r.Short() // buffMp
	if v, _ := r.Byte(); v != 1 {
		t.Fatalf("byte lop/lang = %d, mong doi 1 (= nv.Lang, chi so nClasss[])", v)
	}
}

// TestHoiSinhLucVao: chet (hp <= 0) -> day mau + chakra, giu vi tri; dang song
// -> khong dong vao; nil -> khong panic.
func TestHoiSinhLucVao(t *testing.T) {
	chet := &store.NhanVat{ID: 6, Ten: "narruto", HpToiDa: 150, HpHienTai: 0, ChakraToiDa: 100, ChakraHienTai: 0, ViTriX: 873, ViTriY: 432}
	if !hoiSinhLucVao(chet) {
		t.Fatal("nhan vat chet phai duoc hoi sinh (tra true)")
	}
	if chet.HpHienTai != 150 || chet.ChakraHienTai != 100 {
		t.Fatalf("sau hoi sinh hp/chakra = %d/%d, mong doi 150/100", chet.HpHienTai, chet.ChakraHienTai)
	}
	if chet.ViTriX != 873 || chet.ViTriY != 432 {
		t.Fatalf("vi tri phai giu nguyen, got %d,%d", chet.ViTriX, chet.ViTriY)
	}

	am := &store.NhanVat{HpToiDa: 150, HpHienTai: -5, ChakraToiDa: 100}
	if !hoiSinhLucVao(am) || am.HpHienTai != 150 {
		t.Fatalf("hp am cung la chet: got %d", am.HpHienTai)
	}

	song := &store.NhanVat{HpToiDa: 150, HpHienTai: 37, ChakraToiDa: 100, ChakraHienTai: 12}
	if hoiSinhLucVao(song) {
		t.Fatal("nhan vat dang song khong duoc dong vao")
	}
	if song.HpHienTai != 37 || song.ChakraHienTai != 12 {
		t.Fatalf("nhan vat song bi doi chi so: %d/%d", song.HpHienTai, song.ChakraHienTai)
	}

	if hoiSinhLucVao(nil) {
		t.Fatal("nil phai tra false, khong panic")
	}
}
