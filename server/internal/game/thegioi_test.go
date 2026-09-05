package game

import (
	"bytes"
	"testing"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

func TestNhungTileSetDungKichThuoc(t *testing.T) {
	// File tileset = file goc Ngoc Rong 1696 byte, CONG phan khai bao tile nen
	// Lang La (60-69) them vao bo 1 = 1716 byte. Xem tileset_langla_test.go va
	// tools/map-lang-la/them-tile-vao-tileset.py.
	//
	// Van khoa kich thuoc lai vi day la file NHUNG: nhung nham file khac thi
	// tileId roi ra ngoai [1,numTileMap] va client THOAT APP (loadMap bat
	// Exception roi goi GameMidlet.instance.exit()).
	const coTileSet = 1734
	if len(duLieuTileSet) != coTileSet {
		t.Fatalf("duLieuTileSet = %d byte, mong doi %d (goc 1696 + khai bao tile Lang La: bo 60-69 va bo bo sung 83-86)", len(duLieuTileSet), coTileSet)
	}
	// So bo tile KHONG duoc doi: them tile vao bo san, khong them bo moi.
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
	m, err := goiMeInfo(nv, nil)
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
	// 05/09: byte cgender luon 0 (gioiTinhClient, hinhdang.go) — client chi biet
	// 3 hanh tinh; lang 3 ma ghi 3 thi setDefaultPart() khong dat than/chan.
	b("gender (cgender phia client)", 0)
	// 05/09: head qua dauHopLe — KieuToc 7 la toc Dragon Ball, bo Naruto chi co dau 30.
	sh("head", 30)
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
	// KH7 Task 5: ba mang vat pham mang SO O that (suc chua), roi DU so o than
	// o — nhan vat nay khong co mon nao nen moi o la mot short -1.
	bangRong := func(ten string, soOMongDoi int) {
		t.Helper()
		soO, err := r.Byte()
		if err != nil || int(soO) != soOMongDoi {
			t.Fatalf("%s.size = %d (loi %v), mong doi %d (SO O, khong phai so mon)", ten, soO, err, soOMongDoi)
		}
		for i := 0; i < int(soO); i++ {
			if id, err := r.Short(); err != nil || id != -1 {
				t.Fatalf("%s o %d = %d (loi %v), mong doi -1 (o rong)", ten, i, id, err)
			}
		}
	}
	bangRong("itemsBody", SoOMacChuan)
	bangRong("itemsBag", SoOTuiChuan)
	bangRong("itemsBox", SoORuongChuan)
	// Mot cap anh xa headId -> avatarId, de client biet ve anh nao vao o chan
	// dung tren HUD. Gui 0 cap thi Char.getAvatar tra -1 va o do trong tron.
	sh("headAvatar.size", 1)
	sh("headAvatar[0].headId", HeadNhanVat)
	sh("headAvatar[0].avatarId", AvatarNhanVat)
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
