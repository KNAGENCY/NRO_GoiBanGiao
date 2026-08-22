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
//go:embed assets/tile_set_info.bin
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

	// Bang anh xa headId -> avatarId. Client tra bang nay de biet ve anh nao vao
	// o chan dung tren HUD (Char.getAvatar, dung o Panel.cs). Gui 0 cap thi
	// getAvatar tra -1 va client bo qua — o chan dung TRONG TRON, dung loi da thay
	// trong game.
	w.Short(1) // headAvatar.size = 1 cap
	w.Short(HeadNhanVat)
	w.Short(AvatarNhanVat)
	w.Short(0) // conMeo[gender][0]
	w.Short(0) // conMeo[gender][1]
	w.Short(0) // conMeo[gender][2]

	w.Byte(0) // fusion (khong nhap the)
	w.Int(0)  // deltaTime
	w.Byte(0) // isNewMember — store chua co truong nay, ghi 0 an toan (§8 field 23)

	// (e) BO 3 truong cuoi (aura/effFront/hat).
	return protocol.NewMessage(lenhMeInfo, w), nil
}

// HeadNhanVat / AvatarNhanVat: cap anh xa cho o chan dung tren HUD.
//
//	HeadNhanVat   = part dau cua nhan vat nguoi choi (mapping-id.json, parts[].id)
//	AvatarNhanVat = sprite chan dung rieng trong bo ArtRoot (27041-27082)
//
// Hai so nay phai KHAC nhau: tro chan dung ve chinh anh part dau thi o HUD hien
// mot cai dau bi cat cup chu khong phai chan dung.
const (
	HeadNhanVat   int16 = 30
	AvatarNhanVat int16 = 27078
)

// goiMapInfo dung goi -24 MAP INFO: gio kich switchToMe() ben client (roi man
// "Xin cho"). Byte-exact theo tai lieu §9 (da doi chieu Zone.mapInfo, Zone.java:
// 738-869). LUU Y bay cam:
//
//	(a) BYTE DU PHONG 0 giua mob va npc (Zone.java:805) — QUEN la lech toan bo phan sau.
//	(b) bg-item RONG van phai short 0 (2 byte), KHONG bo trang.
//	(c) x/y lay tu VI TRI NHAN VAT (nv.ViTriX/Y) — Zone dung pl.location, khong phai
//	    spawn ban do; giu dung de KH5 van dung sau khi nguoi choi da di chuyen.
//	(d) -82 PHAI da duoc gui truoc goi nay (xuLyClientOk dam bao thu tu), va tileId
//	    phai ∈ [1,numTileMap] — neu khong client THOAT HAN app.
//
// Ke hoach 6 Task 5: tham so `quai` la danh sach quai dang song trong khu vuc
// nay tai THOI DIEM goi nay dung (§3.6). Ham nay CHI GHI byte — khong tu doc
// trang thai KhuVuc — de giu logic ma hoa THUAN (dev thuan tien test byte-
// chinh-xac ma khong can dung actor). Noi GOI ham nay (xuLyClientOk trong
// login.go) chiu trach nhiem lay dung danh sach quai truoc khi truyen vao.
func goiMapInfo(nv *store.NhanVat, bd *BanDo, quai []*Quai) (*protocol.Message, error) {
	w := protocol.NewWriter()
	w.Byte(byte(bd.MapID)) // mapId (client doc unsigned byte, <=255)
	w.Byte(bd.PlanetID)    // planetId
	w.Byte(bd.TileID)      // tileId ∈ [1,39]
	w.Byte(bd.BgID)        // bgId
	w.Byte(bd.Loai)        // type = 0 MAP_NORMAL
	if err := w.UTF(bd.Ten); err != nil {
		return nil, err
	}
	w.Byte(0) // zoneId = 0 (map mot zone)

	w.Short(int16(nv.ViTriX)) // (c) x = vi tri nhan vat
	w.Short(int16(nv.ViTriY)) // (c) y

	w.Byte(0)                // waypointCount = 0
	ghiDanhSachQuai(w, quai) // mobCount + tung con quai (§3.6, Ke hoach 6 Task 5)
	w.Byte(0)                // (a) BYTE DU PHONG — giua mob va npc, DUNG QUEN
	ghiNPCMap(w, bd)
	w.Byte(0) // itemCount = 0

	w.Short(0) // (b) bg-item: short 0 (map trong, khong bo trang)
	w.Short(0) // effMapCount = 0

	w.Byte(bd.BgType) // bgType
	w.Byte(0)         // idSpaceShip = 0
	// isMapDouble: rieng mapId 148 la map doi (Zone.java:861). Lang 1-5 khong phai
	// 148 nen luon 0, nhung tinh theo mapId de dung ca khi them map 148 sau nay.
	if bd.MapID == 148 {
		w.Byte(1)
	} else {
		w.Byte(0)
	}
	return protocol.NewMessage(lenhMapInfo, w), nil
}
