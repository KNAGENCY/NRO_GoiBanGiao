// Package game — cac goi "vao the gioi" (Ke hoach 4): -82 tileset, -30/0 ME_INFO,
// -24 MAP INFO. Gui sau khi client bao san sang (-28/13, xem xuLyClientOk trong
// login.go) de client roi man "Xin cho" va hien nhan vat dung trong lang.
package game

import (
	_ "embed"
	"fmt"

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
//
// (d2) BA MANG VAT PHAM — PHA A (docs/superpowers/research/2026-09-03-kh7-me-info-so-o.md):
// SAU MOI con so kich co, client doc LUON noi dung tung o NGAY TRONG GOI NAY
// (Controller.cs:5860-5968). Khuon than mot o GIONG HET goi -36/-37/-35
// (ghiCacO, goihanhtrang.go). Ghi so o that ma KHONG ghi than o = lech toan bo
// phan duoi goi.
//
// Con so la SO O (suc chua), KHONG phai so mon: client lap DU so o lan, o rong
// la short -1, va indexUI = chi so o (Controller.cs:5921) — chinh la con so
// client gui nguoc len khi nguoi choi bam vao mot o.
//
// DAU: itemsBody doc unsigned (<=255, :5860); itemsBag va itemsBox doc CO DAU
// (<=127, :5907/:5942), va ngoai le tran o hai cho do KHONG nam trong try/catch
// -> statusMe = 4 (:5969) khong bao gio chay -> nguoi choi ket man "Xin cho".
// Hang tranSoO = 127 (goihanhtrang.go) kep ca ba.
//
// dsMon la CA lat cat vat pham cua nguoi choi (moi noi chua); ham tu loc theo
// NoiChua giong ba goi hanh trang. Truyen LAT CAT chu khong truyen ba goi da
// dung san — de ME_INFO va -36/-37/-35 luon noi cung mot su that.
func goiMeInfo(nv *store.NhanVat, dsMon []store.VatPham) (*protocol.Message, error) {
	// DUNG XONG TOAN BO DU LIEU ROI MOI GHI BYTE: loi o buoc xep o => khong gui
	// gi ca, thay vi mot ME_INFO thieu byte lam client ket "Xin cho".
	oMac, err := xepVaoO(dsMon, store.NoiChuaMac, SoOMacChuan)
	if err != nil {
		return nil, fmt.Errorf("ME_INFO (do mac): %w", err)
	}
	oTui, err := xepVaoO(dsMon, store.NoiChuaTui, SoOTuiChuan)
	if err != nil {
		return nil, fmt.Errorf("ME_INFO (tui): %w", err)
	}
	oRuong, err := xepVaoO(dsMon, store.NoiChuaRuong, SoORuongChuan)
	if err != nil {
		return nil, fmt.Errorf("ME_INFO (ruong): %w", err)
	}
	return goiMeInfoDaXep(nv, oMac, oTui, oRuong)
}

// goiMeInfoDaXep ghi byte sau khi ba mang o da duoc kiem het.
func goiMeInfoDaXep(nv *store.NhanVat, oMac, oTui, oRuong []*store.VatPham) (*protocol.Message, error) {
	w := protocol.NewWriter()
	w.Byte(0) // (a) byte lenh con 0

	w.Int(int32(nv.ID))
	w.Byte(0) // taskMain.id — chua co nhiem vu
	// (c) byte "gender" thu nhat -> Char.cgender (Controller.cs:5832): client dung
	// de chon than/chan MAC DINH (setDefaultPart, :5863) va hieu ung theo hanh
	// tinh — luon 0, xem gioiTinhClient (hinhdang.go). Byte "gender lan 2" o duoi
	// moi la chi so lop/lang (nClasss[]) va VAN = nv.Lang.
	w.Byte(gioiTinhClient(nv.Lang))
	w.Short(dauHopLe(nv.KieuToc)) // head — quy ve dau Naruto co that (hinhdang.go)
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

	// (d2) itemsBody: so o (unsigned byte) roi DU so o than o.
	w.Byte(byte(len(oMac)))
	if err := ghiCacO(w, oMac); err != nil {
		return nil, fmt.Errorf("ME_INFO (do mac): %w", err)
	}
	// (d2) itemsBag: so o (byte CO DAU ben client) roi DU so o than o.
	w.Byte(byte(len(oTui)))
	if err := ghiCacO(w, oTui); err != nil {
		return nil, fmt.Errorf("ME_INFO (tui): %w", err)
	}
	// (d2) itemsBox: so o (byte CO DAU ben client) roi DU so o than o.
	w.Byte(byte(len(oRuong)))
	if err := ghiCacO(w, oRuong); err != nil {
		return nil, fmt.Errorf("ME_INFO (ruong): %w", err)
	}

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
//	AvatarNhanVat = sprite chan dung RIENG cua nguoi choi
//
// KHONG duoc dung 27078: do la 'iruka-portrait', chan dung NPC thay Iruka da
// dang ky trong catalog.json. Dung chung thi nguoi choi va NPC co cung mot khuon
// mat. 27084 la anh cat tu chinh anh dau nhan vat (27041), thu nho vua o xanh
// tren HUD — xem tools/map-lang-la/tao-chan-dung-nguoi-choi.py.
const (
	HeadNhanVat   int16 = 30
	AvatarNhanVat int16 = 27084
)

// goiMapInfo dung goi -24 MAP INFO: gio kich switchToMe() ben client (roi man
// "Xin cho"). Byte-exact theo tai lieu §9 (da doi chieu Zone.mapInfo, Zone.java:
// 738-869). LUU Y bay cam:
//
//	(a) BYTE DU PHONG 0 giua mob va npc (Zone.java:805) — QUEN la lech toan bo phan sau.
//	(b) bg-item (bong do) RONG van phai short 0 (2 byte), KHONG bo trang —
//	    xem ghiBgItem trong bgitem.go. Bong do CHI hien neu goi -31 (bang mau)
//	    da duoc gui truoc do trong phien, khong thi client bo qua im lang.
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

	if err := ghiDiemDi(w, bd); err != nil { // waypointCount + tung cong roi map
		return nil, err
	}
	ghiDanhSachQuai(w, quai) // mobCount + tung con quai (§3.6, Ke hoach 6 Task 5)
	w.Byte(0)                // (a) BYTE DU PHONG — giua mob va npc, DUNG QUEN
	ghiNPCMap(w, bd)
	w.Byte(0) // itemCount = 0

	ghiBgItem(w, bd) // (b) bg-item — map trong van ghi short 0, KHONG bo trang
	w.Short(0)       // effMapCount = 0

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

// ghiDiemDi ghi o waypoint cua MAP INFO.
//
// KHUON BYTE (Zone.java:753-767 == Controller.cs:5055 — DA DOI CHIEU CA HAI
// PHIA theo luat §6.16 cua BAT-DAU-TU-DAY):
//
//	byte soCong
//	moi cong: short MinX | short MinY | short MaxX | short MaxY |
//	          bool LaCong | bool NgoaiTuyen | UTF Ten
//
// soCong la MOT byte (khong phai short) — Zone.java writeByte, client
// readByte. Mot ban do co hon 255 cong la khong the, nhung neu co thi phai
// cat chu khong duoc de tran am tham; hien moi ban do chi co 1 cong nen chi
// can chan bang mot loi ro rang.
//
// GIU DUNG MOT CONG MOI BAN DO cho toi khi co ly do that su phai them:
// Char.isInEnterOnlinePoint (Char.cs:1656-1668) tra cuu song song hai danh
// sach TileMap.vGo va PopUp.vPopups BANG CUNG MOT CHI SO, tuc no ngam gia
// dinh hai danh sach xep cung thu tu. Voi mot cong thi gia dinh do luon
// dung; voi nhieu cong no phu thuoc vao thu tu cac popup khac duoc them —
// mot cho de vo am tham.
func ghiDiemDi(w *protocol.Writer, bd *BanDo) error {
	if bd == nil {
		w.Byte(0)
		return nil
	}
	if len(bd.DiemDi) > 255 {
		return fmt.Errorf("map %d co %d cong roi map, vuot 255 (soCong la 1 byte)", bd.MapID, len(bd.DiemDi))
	}
	w.Byte(byte(len(bd.DiemDi)))
	for _, d := range bd.DiemDi {
		w.Short(d.MinX)
		w.Short(d.MinY)
		w.Short(d.MaxX)
		w.Short(d.MaxY)
		w.Bool(d.LaCong)
		w.Bool(d.NgoaiTuyen)
		if err := w.UTF(d.Ten); err != nil {
			return fmt.Errorf("ten cong %q cua map %d khong ghi duoc: %w", d.Ten, bd.MapID, err)
		}
	}
	return nil
}
