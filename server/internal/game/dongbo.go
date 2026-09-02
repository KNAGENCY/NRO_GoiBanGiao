// Package game — phan dong bo du lieu tinh (Ke hoach 3).
//
// Khi vao the gioi, client so 4 phien ban (data/map/skill/item) may chu gui
// trong -28/4 voi cache cua no; lech thi tu xin cap nhat. Client MOI (RMS trong)
// luon xin lai ca 4 du phien ban co khop hay khong, nen may chu PHAI phuc vu du
// bon o lan dau. Day la cho ta thay 3 lop Dragon Ball cu bang 5 lang Naruto.
package game

import (
	_ "embed"

	"nroserver/internal/protocol"
)

// Năm runtime dưới đây được chép nguyên từ một build đã khóa SHA. Không dựng lại
// trong Go vì chỉ cần lệch một byte là server/client dùng hai thế giới khác nhau.
//
//go:embed assets/part.bin
var duLieuPart []byte

//go:embed assets/skill.bin
var duLieuSkill []byte

//go:embed assets/map-1-collision.bin
var duLieuCollision []byte

//go:embed assets/catalog.json
var duLieuCatalog []byte

//go:embed assets/asset-lock.json
var duLieuAssetLock []byte

// Ba version runtime phải bằng đúng asset-lock, không còn là số chọn tay. Item
// chưa thuộc lát cắt artifact nên giữ mốc Naruto hiện hành của server.
const (
	vsData  byte = 26 // moi lan part.bin doi phai tang, khong client dung ban cu trong RMS
	vsMap   byte = 22
	vsSkill byte = 21
	vsItem  byte = 20
)

// lenhDuLieu: goi data. Client gui -87 than rong de xin, may chu tra -87 day du.
const lenhDuLieu int8 = -87

// Cac lenh con cua khung -28 lien quan dong bo tinh.
const (
	lenhConVersionGame  byte = 4  // may chu -> client
	lenhConCatalogMap   byte = 6  // client xin map / may chu tra
	lenhConCatalogSkill byte = 7  // client xin skill / may chu tra
	lenhConCatalogItem  byte = 8  // client xin item / may chu tra
	lenhConClientOk     byte = 13 // client bao san sang
)

// goiVersionGame dung goi -28/4: bao client 4 phien ban du lieu tinh + bang moc
// kinh nghiem rong. Client so 4 phien ban nay voi cache; lech thi tu xin cap nhat.
// Xem tai lieu giao thuc §2.
func goiVersionGame() *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(lenhConVersionGame) // 4
	w.Byte(vsData)
	w.Byte(vsMap)
	w.Byte(vsSkill)
	w.Byte(vsItem)
	w.Byte(0) // byte dem, client doc roi bo
	w.Byte(0) // countCaption = 0 (bang moc kinh nghiem rong, an toan cho vao map)
	return protocol.NewMessage(lenhKhongTrenBanDo, w)
}

// goiCatalogData dùng gói -87: vsData rồi 6 blob đúng thứ tự. Part và skill là
// byte artifact thật; dart/arrow/effect/image chưa dùng nên là count=0 hợp lệ.
func goiCatalogData() *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(vsData)
	// blob rong = do dai 2 byte, noi dung la mot short 0 (count = 0 phan tu).
	blobRong := func() {
		w.Int(2)
		w.Short(0)
	}
	blobRong()                    // dart
	blobRong()                    // arrow
	blobRong()                    // effect
	blobRong()                    // image
	w.Int(int32(len(duLieuPart))) // part: do dai that
	w.ThemMang(duLieuPart)        // part: noi dung
	w.Int(int32(len(duLieuSkill)))
	w.ThemMang(duLieuSkill)
	return protocol.NewMessage(lenhDuLieu, w)
}

// goiCatalogSkill dung goi -28/7: DANH SACH 5 LANG NARUTO thay cho 3 lop Dragon
// Ball cu. Client dung byte "gender" trong ME_INFO lam chi so nClasss[] nen
// countClass PHAI >= so lang (5) — id lang 3/4 se sap neu client con cache 3 lop
// cu. Moi lang gui ten + 4 MauChieu that (Ke hoach 6, Task 1) — xem chieuthuc.go
// va tai lieu giao thuc chien dau §3.7/§4.3.
func goiCatalogSkill() (*protocol.Message, error) {
	w := protocol.NewWriter()
	w.Byte(lenhConCatalogSkill) // 7
	w.Byte(vsSkill)
	w.Byte(0) // countOption = 0
	dsLang := DanhSachLang()
	w.Byte(byte(len(dsLang))) // countClass = 5
	for _, l := range dsLang {
		if err := w.UTF(l.Ten); err != nil {
			return nil, err
		}
		chieu := ChieuCuaLang(byte(l.ID))
		w.Byte(byte(len(chieu))) // countTemplate = 4
		for _, m := range chieu {
			if err := ghiMauChieu(w, m); err != nil {
				return nil, err
			}
		}
	}
	return protocol.NewMessage(lenhKhongTrenBanDo, w), nil
}

func goiCatalogMap() *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(lenhConCatalogMap) // 6
	w.Byte(vsMap)
	c := catalogLangLa
	w.Byte(byte(len(c.Maps)))
	for _, banDo := range c.Maps {
		_ = w.UTF(banDo.Name)
	}
	w.Byte(1)
	_ = w.UTF(c.NPC.Name)
	w.Short(c.NPC.Head)
	w.Short(c.NPC.Body)
	w.Short(c.NPC.Leg)
	w.Byte(0)
	w.Byte(1)
	w.Byte(0)
	_ = w.UTF(c.Mob.Name)
	w.Int(c.Mob.HP)
	w.Byte(c.Mob.Range)
	w.Byte(c.Mob.Speed)
	w.Byte(c.Mob.Dart)
	return protocol.NewMessage(lenhKhongTrenBanDo, w)
}

// goiCatalogItem dung message -28/8 subtype 2 khoang RONG — goi DONG DAI cua
// chuoi catalog vat pham. Chi subtype 2 moi kich Service.clientOk() phia client
// (Controller.cs:6517-6534), nen no PHAI la message -28/8 gui cuoi cung.
//
// Noi dung mau vat pham that nam o subtype 0 va 1 — xem mauvatpham.go,
// goiCatalogItemDayDu(). Ham nay giu nguyen 7 byte da duoc ghim tu Ke hoach 3.
func goiCatalogItem() *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(lenhConCatalogItem) // 8
	w.Byte(vsItem)
	w.Byte(2)  // subtype 2 (item them)
	w.Short(0) // start
	w.Short(0) // end
	return protocol.NewMessage(lenhKhongTrenBanDo, w)
}
