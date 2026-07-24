// Package game — cac goi "vao the gioi" (Ke hoach 4): -82 tileset, -30/0 ME_INFO,
// -24 MAP INFO. Gui sau khi client bao san sang (-28/13, xem xuLyClientOk trong
// login.go) de client roi man "Xin cho" va hien nhan vat dung trong lang.
package game

import (
	_ "embed"

	"nroserver/internal/protocol"
)

// duLieuTileSet: NGUYEN file 1-MAY-CHU-GAME/data/map/tile_set_info (1696 byte,
// 39 tileset). Gui nguyen trong goi -82, khong header, khong bien doi (§7). Thieu
// -82, gui SAU -24, hoac map.tileId ngoai [1,numTileMap] -> client GameMidlet.exit()
// (THOAT HAN app). Nhung bang go:embed nhu part.bin o Ke hoach 3.
//
//go:embed tileset.bin
var duLieuTileSet []byte

// Lenh cac goi vao the gioi. Gia tri tu ban Java cu, DUNG DOI:
//
//	-82  DataGame.sendTileSetInfo (DataGame.java:424) — Message(-82) + nguyen file,
//	     KHONG byte lenh con.
//	-30  Service.player (Service.java:782) — messageSubCommand((byte)0) =
//	     Message(-30) + writeByte(0). Byte lenh con 0 la byte DAU than (bang §8 bo sot).
//	-24  Zone.mapInfo (Zone.java:741) — Message(-24) truc tiep, KHONG byte lenh con.
const (
	lenhTileSetInfo int8 = -82
	lenhMeInfo      int8 = -30
	lenhMapInfo     int8 = -24
)

// goiTileSetInfo dung goi -82: gui NGUYEN file tile_set_info. -82 KHONG nam trong
// lenhDoDaiBaByte (protocol/frame.go) nen dung do dai 2 byte — 1696 < 65535, vua.
// Xem tai lieu giao thuc §7.
func goiTileSetInfo() *protocol.Message {
	w := protocol.NewWriter()
	w.ThemMang(duLieuTileSet)
	return protocol.NewMessage(lenhTileSetInfo, w)
}
