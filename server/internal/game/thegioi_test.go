package game

import (
	"bytes"
	"testing"
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
