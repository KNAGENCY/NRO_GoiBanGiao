package game

import (
	"encoding/binary"
	"testing"
)

// Bo tile cua Lang La danh so 60-69, nam NGOAI bang tile_set_info goc cua Ngoc
// Rong (chi khai bao 1-39). Thieu khai bao thi client dat types[] = 0 cho ca
// map -> nhan vat roi xuyen dat xuong day. Cac test duoi khoa dieu do lai.
//
// Y nghia bit (doc tu client Char.cs):
//
//	0x0002  san, dung duoc len tren      -> Char.cs:2735 tileTypeAtPixel(cx, cy+1) & 2
//	0x0004  chan ben phai                -> Char.cs:2420 tileTypeAt(cx+chw, cy-chh, 4)
//	0x0008  chan ben trai                -> Char.cs:2427 tileTypeAt(cx-chw-1, cy-chh, 8)
//	0x2000  tran, chan tu duoi nhay len  -> Char.cs:2097 tileTypeAt(cx, cy-ch, 8192)
const (
	bitSan    = 0x0002
	bitChanP  = 0x0004
	bitChanT  = 0x0008
	bitTran   = 0x2000
	boTileMap = 1 // map 1 dung tileID = 1 -> chi so bo 0
)

// phanTichTileSet doc file tile_set_info thanh: boTile -> loai -> tap tile.
//
// Dinh dang (doc tu client Controller.cs case -82):
//
//	byte  soBo
//	  byte  soNhom
//	    int32 loai (big-endian)
//	    byte  soTile
//	      byte tile ...
func phanTichTileSet(t *testing.T, data []byte) []map[int32]map[byte]bool {
	t.Helper()
	if len(data) == 0 {
		t.Fatal("tile_set_info rong")
	}
	i := 0
	doc1 := func() byte {
		if i >= len(data) {
			t.Fatalf("tile_set_info het du lieu som o byte %d", i)
		}
		b := data[i]
		i++
		return b
	}
	soBo := int(doc1())
	ket := make([]map[int32]map[byte]bool, soBo)
	for bo := 0; bo < soBo; bo++ {
		ket[bo] = map[int32]map[byte]bool{}
		soNhom := int(doc1())
		for n := 0; n < soNhom; n++ {
			if i+4 > len(data) {
				t.Fatalf("tile_set_info thieu 4 byte loai o bo %d nhom %d", bo, n)
			}
			loai := int32(binary.BigEndian.Uint32(data[i : i+4]))
			i += 4
			soTile := int(doc1())
			if ket[bo][loai] == nil {
				ket[bo][loai] = map[byte]bool{}
			}
			for k := 0; k < soTile; k++ {
				ket[bo][loai][doc1()] = true
			}
		}
	}
	if i != len(data) {
		t.Fatalf("tile_set_info thua %d byte cuoi (doc het %d/%d)", len(data)-i, i, len(data))
	}
	return ket
}

func TestTileSetKhaiBaoTileLangLaLaSan(t *testing.T) {
	bo := phanTichTileSet(t, duLieuTileSet)
	if len(bo) <= boTileMap-1 {
		t.Fatalf("tile_set_info chi co %d bo, map 1 can bo thu %d", len(bo), boTileMap)
	}
	san := bo[boTileMap-1][bitSan]
	if san == nil {
		t.Fatalf("bo tile %d khong co nhom san (0x%04X)", boTileMap, bitSan)
	}
	// 60..69 la toan bo tile nen Lang La: dat, dat_2, mat co, hai goc doc,
	// hai vach, hai tile du phong, duong lat da. Thieu bat ky cai nao thi cho
	// do nhan vat lot xuyen.
	for tile := byte(60); tile <= 69; tile++ {
		if !san[tile] {
			t.Errorf("tile %d KHONG duoc khai bao la san (bit 0x%04X) -> nhan vat roi xuyen cho nay", tile, bitSan)
		}
	}
}

func TestTileSetKhoiDatDacChanCaBonPhia(t *testing.T) {
	bo := phanTichTileSet(t, duLieuTileSet)[boTileMap-1]
	// Dat than (60, 61) va hai vach go (65, 66) la khoi DAC: khong duoc di
	// xuyen ngang vao long dat, cung khong nhay xuyen tu duoi len.
	for _, tr := range []struct {
		bit int32
		ten string
	}{
		{bitChanP, "chan ben phai"},
		{bitChanT, "chan ben trai"},
	} {
		tap := bo[tr.bit]
		if tap == nil {
			t.Fatalf("bo tile khong co nhom %s (0x%04X)", tr.ten, tr.bit)
		}
		for _, tile := range []byte{60, 61, 65, 66} {
			if !tap[tile] {
				t.Errorf("tile %d thieu bit %s (0x%04X) -> di xuyen duoc vao long dat", tile, tr.ten, tr.bit)
			}
		}
	}
	tran := bo[bitTran]
	if tran == nil {
		t.Fatalf("bo tile khong co nhom tran (0x%04X)", bitTran)
	}
	for _, tile := range []byte{60, 61} {
		if !tran[tile] {
			t.Errorf("tile %d thieu bit tran (0x%04X) -> nhay tu duoi xuyen len duoc", tile, bitTran)
		}
	}
}

func TestTileSetGiuNguyenKhaiBaoTileGoc(t *testing.T) {
	// Them tile Lang La KHONG duoc lam mat khai bao cu: cac map khac van dung
	// bo tile goc, xoa nham la vo map cu ma khong bao loi gi.
	bo := phanTichTileSet(t, duLieuTileSet)[boTileMap-1]
	goc := map[int32][]byte{
		bitSan:   {1, 2, 3, 4, 5, 6, 18, 31, 32, 33, 34, 39},
		bitChanP: {1, 23, 28, 39, 29, 30},
		bitChanT: {18, 24, 29, 39, 28, 30},
		bitTran:  {28, 29, 30, 39},
		0x200:    {26},
	}
	for bit, dsTile := range goc {
		tap := bo[bit]
		if tap == nil {
			t.Fatalf("mat han nhom 0x%04X cua bang tile goc", bit)
		}
		for _, tile := range dsTile {
			if !tap[tile] {
				t.Errorf("nhom 0x%04X mat tile goc %d", bit, tile)
			}
		}
	}
}
