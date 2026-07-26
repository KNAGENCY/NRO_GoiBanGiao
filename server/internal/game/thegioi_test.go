package game

import (
	"bytes"
	"testing"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

func TestNhungTileSetDungKichThuoc(t *testing.T) {
	// File tileset nhung phai la file goc 1696 byte, byte dau (numTileMap) = 39.
	// Sai kich thuoc = nhung nham -> tileId ngoai [1,numTileMap] -> client thoat app.
	if len(duLieuTileSet) != 1696 {
		t.Fatalf("duLieuTileSet = %d byte, mong doi 1696 (file tile_set_info goc)", len(duLieuTileSet))
	}
	if duLieuTileSet[0] != 39 {
		t.Errorf("numTileMap (byte dau) = %d, mong doi 39", duLieuTileSet[0])
	}
}

func TestGoiTileSetInfoGuiNguyenFile(t *testing.T) {
	m := goiTileSetInfo()
	if m.Command != -82 {
		t.Fatalf("lenh = %d, mong doi -82", m.Command)
	}
	// Than = NGUYEN file, khong header, khong byte lenh con (§7).
	if !bytes.Equal(m.Data, duLieuTileSet) {
		t.Errorf("than -82 khong bang nguyen file tileset (dai %d, mong doi %d)", len(m.Data), len(duLieuTileSet))
	}
}

func TestGoiMeInfoByteDung(t *testing.T) {
	nv := &store.NhanVat{
		ID:      12345,
		Ten:     "Naruto",
		Lang:    3, // lang Da -> gender-slot 3, phai < countClass (5)
		KieuToc: 7,
		SucManh: 1_000_000,
		Vang:    250,
		Ngoc:    9,
		MapID:   4,
		ViTriX:  100,
		ViTriY:  384,
	}
	m, err := goiMeInfo(nv)
	if err != nil {
		t.Fatalf("goiMeInfo loi: %v", err)
	}
	if m.Command != -30 {
		t.Fatalf("lenh = %d, mong doi -30", m.Command)
	}

	r := protocol.NewReader(m.Data)
	b := func(ten string, muon byte) {
		g, err := r.Byte()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}
	sh := func(ten string, muon int16) {
		g, err := r.Short()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}
	in := func(ten string, muon int32) {
		g, err := r.Int()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}
	lo := func(ten string, muon int64) {
		g, err := r.Long()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}

	b("byte lenh con", 0) // §8 bo sot — messageSubCommand((byte)0)
	in("id", 12345)
	b("taskMain.id", 0)
	b("gender", 3) // = NhanVat.Lang
	sh("head", 7)  // = NhanVat.KieuToc
	if ten, err := r.UTF(); err != nil || ten != "Naruto" {
		t.Fatalf("name = %q (loi %v), mong doi \"Naruto\"", ten, err)
	}
	b("cPK", 0)
	b("typePk", 0)
	lo("power", 1_000_000) // = NhanVat.SucManh
	sh("buffHp", 0)
	sh("buffMp", 0)
	b("gender lan 2", 3) // = NhanVat.Lang, chi so nClasss[]
	b("skillCount", 0)
	lo("gold", 250) // LONG 8 byte = NhanVat.Vang
	in("ruby", 0)
	in("gem", 9) // = NhanVat.Ngoc
	b("itemsBody.size", 0)
	b("itemsBag.size", 0)
	b("itemsBox.size", 0)
	sh("headAvatar.size", 0)
	sh("conMeo0", 0)
	sh("conMeo1", 0)
	sh("conMeo2", 0)
	b("fusion", 0)
	in("deltaTime", 0)
	b("isNewMember", 0)

	// Bo 3 truong cuoi (aura/effFront/hat) -> phai HET buffer o day.
	if _, err := r.Byte(); err == nil {
		t.Error("ME_INFO con byte thua sau isNewMember — 3 truong cuoi (aura/effFront/hat) phai bo han (§8 note d)")
	}
}

func TestGoiMapInfoByteDung(t *testing.T) {
	nv := &store.NhanVat{MapID: 4, ViTriX: 100, ViTriY: 384}
	bd := TimBanDo(4)                 // Lang Da: planet 0, tile 2, bg 1, bgType 0, loai 0, ten "Làng Đá"
	m, err := goiMapInfo(nv, bd, nil) // nil = map trong quai, giong KH4
	if err != nil {
		t.Fatalf("goiMapInfo loi: %v", err)
	}
	if m.Command != -24 {
		t.Fatalf("lenh = %d, mong doi -24", m.Command)
	}

	r := protocol.NewReader(m.Data)
	b := func(ten string, muon byte) {
		g, err := r.Byte()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}
	sh := func(ten string, muon int16) {
		g, err := r.Short()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}

	b("mapId", 4)
	b("planetId", 0)
	b("tileId", 2) // ∈ [1,39]
	b("bgId", 1)
	b("type", 0) // MAP_NORMAL
	if ten, err := r.UTF(); err != nil || ten != "Làng Đá" {
		t.Fatalf("mapName = %q (loi %v), mong doi \"Làng Đá\"", ten, err)
	}
	b("zoneId", 0)
	sh("x", 100) // = NhanVat.ViTriX
	sh("y", 384) // = NhanVat.ViTriY
	b("waypointCount", 0)
	b("mobCount", 0)
	b("byte du phong", 0) // Zone.java:805 — giua mob va npc, DUNG QUEN
	b("npcCount", 0)
	b("itemCount", 0)
	sh("bgItem", 0) // short 0 (KHONG bo trang)
	sh("effMap", 0) // short 0
	b("bgType", 0)
	b("idSpaceShip", 0)
	b("isMapDouble", 0) // mapId != 148

	if _, err := r.Byte(); err == nil {
		t.Error("MAP INFO con byte thua sau isMapDouble — phai het buffer")
	}
}

func TestGoiMapInfoDungViTriNhanVat(t *testing.T) {
	// x/y trong MAP INFO lay tu VI TRI NHAN VAT (nv.ViTriX/Y), khong phai spawn
	// mac dinh cua ban do — de KH5 (di chuyen) van dung sau khi nguoi choi da dich.
	nv := &store.NhanVat{MapID: 1, ViTriX: 512, ViTriY: 240}
	m, err := goiMapInfo(nv, TimBanDo(1), nil)
	if err != nil {
		t.Fatal(err)
	}
	r := protocol.NewReader(m.Data)
	for i := 0; i < 5; i++ {
		r.Byte() // mapId,planet,tile,bg,type
	}
	r.UTF()  // mapName
	r.Byte() // zoneId
	x, _ := r.Short()
	y, _ := r.Short()
	if x != 512 || y != 240 {
		t.Errorf("x/y = %d/%d, mong doi 512/240 (vi tri nhan vat, khong phai spawn ban do)", x, y)
	}
}
