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

// duLieuPart: du lieu sprite nhan vat (dau/than/chan). Day la PLACEHOLDER cua
// ban Dragon Ball cu — nhan vat hien ra bang sprite cu; thay bang part Naruto
// khi duong ong my thuat giao. Client ve nhan vat bang parts[head/body/leg]
// KHONG kiem bien, nen part rong = sap khi ve (Ke hoach 4). File nhung phai phu
// du chi so head/body/leg cua moi nhan vat -> forward nguyen file cu la an toan.
//
//go:embed part.bin
var duLieuPart []byte

// Phien ban du lieu tinh cua ban Naruto. PHAI khac ban Dragon Ball cu
// (vsData=10, vsMap=2, vsSkill=1, vsItem=4) de client khong dung cache cu — nhat
// la vsSkill != 1 de khong dung cache 3 lop cu. 20 la gia tri tuy chon != cu va
// != 0. Bump khi doi noi dung tuong ung o cac ke hoach sau.
const (
	vsData  byte = 20
	vsMap   byte = 20
	vsSkill byte = 20
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

// goiCatalogData dung goi -87: vsData roi 6 blob (dart/arrow/effect/image/part/
// skill) theo dung thu tu. Chi part co noi dung (du lieu sprite nhan vat); 5 blob
// kia rong (do dai 2 = mot short 0). Client chi LUU raw vao RMS, khong parse ngay
// — loi part chi lo khi VE nhan vat (Ke hoach 4). Xem tai lieu giao thuc §3.
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
	blobRong()                    // skill
	return protocol.NewMessage(lenhDuLieu, w)
}

// goiCatalogSkill dung goi -28/7: DANH SACH 5 LANG NARUTO thay cho 3 lop Dragon
// Ball cu. Client dung byte "gender" trong ME_INFO lam chi so nClasss[] nen
// countClass PHAI >= so lang (5) — id lang 3/4 se sap neu client con cache 3 lop
// cu. Moi lang gui ten + 0 template (chua co skill). Xem tai lieu §5.
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
		w.Byte(0) // countTemplate = 0 (chua co skill)
	}
	return protocol.NewMessage(lenhKhongTrenBanDo, w), nil
}

// goiCatalogMap dung goi -28/6 RONG: 0 map, 0 npc, 0 mob. Du de qua bat tay;
// noi dung map/npc/mob that se them o cac ke hoach sau. Xem tai lieu §4.
func goiCatalogMap() *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(lenhConCatalogMap) // 6
	w.Byte(vsMap)
	w.Byte(0) // countMap
	w.Byte(0) // countNpc
	w.Byte(0) // countMob
	return protocol.NewMessage(lenhKhongTrenBanDo, w)
}

// goiCatalogItem dung goi -28/8 toi thieu: dung MOT message con subtype 2 rong.
// Chi subtype 2 moi kich clientOk (khi du 4 phien ban khop). Nhan vat moi khong
// do van ve duoc voi catalog item rong. Xem tai lieu §6.
func goiCatalogItem() *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(lenhConCatalogItem) // 8
	w.Byte(vsItem)
	w.Byte(2)  // subtype 2 (item them)
	w.Short(0) // start
	w.Short(0) // end
	return protocol.NewMessage(lenhKhongTrenBanDo, w)
}
