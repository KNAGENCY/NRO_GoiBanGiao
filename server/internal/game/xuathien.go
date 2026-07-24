package game

import (
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// Lệnh xuất hiện: người chơi vừa vào khu vực, gửi cho người khác để họ vẽ ra.
const lenhXuatHien int8 = -5

// goiXuatHien dung goi -5: du lieu mot nhan vat de client KHAC ve no khi no vao
// khu vuc. 28 truong byte-dung theo Zone.infoPlayer (da soi cheo). KH5 dien 0/hang
// cho cac he chua co (buff/pet/ngua/aura/bang hoi). LUU Y: hp/hpMax la writeInt
// (4 byte) vi client doc readInt; gender ghi 2 LAN; truong 15 la hang -1 (0xFF tren day).
func goiXuatHien(nv *store.NhanVat) (*protocol.Message, error) {
	w := protocol.NewWriter()
	w.Int(int32(nv.ID))                   // 1  playerId
	w.Int(-1)                             // 2  clanId (chua co bang hoi)
	w.Byte(byte(nv.Cap))                  // 3  level (= cap nhan vat)
	w.Bool(false)                         // 4  co co dinh (writeBoolean false)
	w.Byte(0)                             // 5  typePk
	w.Byte(byte(nv.Lang))                 // 6  gender (= id lang)
	w.Byte(byte(nv.Lang))                 // 7  gender ghi LAN 2
	w.Short(nv.KieuToc)                   // 8  head
	if err := w.UTF(nv.Ten); err != nil { // 9 name
		return nil, err
	}
	w.Int(nv.HpToiDa)         // 10 hp (writeInt vi cn.readInt=true)
	w.Int(nv.HpToiDa)         // 11 hpMax
	w.Short(0)                // 12 body (chua co trang bi)
	w.Short(0)                // 13 leg
	w.Byte(0)                 // 14 flagbag
	w.Byte(0xFF)              // 15 hang -1 (writeByte(-1) = 0xFF tren day)
	w.Short(int16(nv.ViTriX)) // 16 x
	w.Short(int16(nv.ViTriY)) // 17 y
	w.Short(0)                // 18 effBuffHp
	w.Short(0)                // 19 effBuffMp
	w.Byte(0)                 // 20 numEff
	w.Byte(0)                 // 21 idSpaceShip
	w.Byte(0)                 // 22 isMonkey
	w.Short(0)                // 23 mount
	w.Byte(0)                 // 24 cFlag
	w.Byte(0)                 // 25 hang 0
	w.Short(0)                // 26 aura
	w.Byte(0)                 // 27 effFront
	w.Short(0)                // 28 hat
	return protocol.NewMessage(lenhXuatHien, w), nil
}

// (Ban goc con gui kem -103 (co) va -8 (neu dang chet). KH5 BO ca hai vi nhan
// vat luon song va chua co co bang hoi — them lai o KH sau khi co he thong tuong ung.)
