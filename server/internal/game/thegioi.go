// Package game — cac goi "vao the gioi" (Ke hoach 4): -82 tileset, -30/0 ME_INFO,
// -24 MAP INFO. Gui sau khi client bao san sang (-28/13, xem xuLyClientOk trong
// login.go) de client roi man "Xin cho" va hien nhan vat dung trong lang.
package game

import (
	_ "embed"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
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

// goiMeInfo dung goi -30/0 ME_INFO: du lieu nhan vat cua chinh nguoi choi de client
// dung "nhan vat cua toi" va set statusMe. Byte-exact theo tai lieu §8 (da doi chieu
// Service.player, Service.java:776-897). LUU Y bay cam:
//
//	(a) THAN BAT DAU BANG byte lenh con 0 (messageSubCommand((byte)0)) — bang §8 bo
//	    sot byte nay; giong quy uoc -28/X.
//	(b) gold la LONG 8 BYTE (client luon readLong) — ghi 4 byte lech toan bo phan sau.
//	(c) gender ghi 2 LAN, deu = NhanVat.Lang (0-4); lan 2 la chi so nClasss[], phai
//	    < countClass (=5, xem -28/7). Sai -> IndexOutOfRange bi nuot -> ket "Xin cho".
//	(d) 3 mang item (Body/Bag/Box) size = 0 cho nhan vat moi — item id la khong co
//	    trong catalog -28/8 se NPE giua goi.
//	(e) 3 truong cuoi (aura short, effFront byte, hat short) BO HAN — client bat loi em.
func goiMeInfo(nv *store.NhanVat) (*protocol.Message, error) {
	w := protocol.NewWriter()
	w.Byte(0) // (a) byte lenh con 0

	w.Int(int32(nv.ID))
	w.Byte(0)             // taskMain.id — chua co nhiem vu
	w.Byte(byte(nv.Lang)) // (c) gender = id lang
	w.Short(nv.KieuToc)   // head
	if err := w.UTF(nv.Ten); err != nil {
		return nil, err
	}
	w.Byte(0)             // cPK
	w.Byte(0)             // typePk
	w.Long(nv.SucManh)    // power (long)
	w.Short(0)            // buffHp
	w.Short(0)            // buffMp
	w.Byte(byte(nv.Lang)) // (c) gender lan 2 — chi so nClasss[], phai < countClass
	w.Byte(0)             // skillCount = 0 (khong ghi short skillId nao)

	w.Long(nv.Vang)       // (b) gold LONG 8 byte
	w.Int(0)              // ruby
	w.Int(int32(nv.Ngoc)) // gem

	w.Byte(0) // (d) itemsBody.size
	w.Byte(0) // (d) itemsBag.size
	w.Byte(0) // (d) itemsBox.size

	w.Short(0) // headAvatar.size (DataGame.sendHeadAvatar rong)
	w.Short(0) // conMeo[gender][0]
	w.Short(0) // conMeo[gender][1]
	w.Short(0) // conMeo[gender][2]

	w.Byte(0) // fusion (khong nhap the)
	w.Int(0)  // deltaTime
	w.Byte(0) // isNewMember — store chua co truong nay, ghi 0 an toan (§8 field 23)

	// (e) BO 3 truong cuoi (aura/effFront/hat).
	return protocol.NewMessage(lenhMeInfo, w), nil
}
