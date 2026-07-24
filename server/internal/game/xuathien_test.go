package game

import (
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
