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
