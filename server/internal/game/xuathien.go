package game

import (
	"fmt"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// Lệnh xuất hiện: người chơi vừa vào khu vực, gửi cho người khác để họ vẽ ra.
const lenhXuatHien int8 = -5

// HinhDangHien la phan "trong nhu the nao" cua goi -5, tach rieng khoi phan
// "la ai" (id/ten/cap/hp/toa do). Tach ra vi BOSS cung xuat hien bang goi -5
// nhung khong phai *store.NhanVat va co mac do that.
//
// Outfit dung DUNG thu tu ma goi -5 ghi ra day: head, body, leg, flagBag,
// aura, effFront (trung thu tu voi PhaBoss.Outfit trong boss.go).
//
// ⚠️ flagBag (o 3) va effFront (o 5) tren day CHI MOT BYTE — dat gia tri > 255
// se bi cat con byte thap. Bon o con lai la 2 byte (short).
type HinhDangHien struct {
	Lang   int16    // = "gender" ban goc; goi -5 ghi gia tri nay HAI LAN
	Outfit [6]int16 // head, body, leg, flagBag, aura, effFront
	TypePk byte
}

// goiXuatHienDayDu dung goi -5: du lieu MOT nhan vat (nguoi choi hoac boss) de
// client KHAC ve no khi no vao khu vuc. 28 truong byte-dung theo Zone.infoPlayer
// (da soi cheo). KH5 dien 0/hang cho cac he chua co (buff/pet/ngua/bang hoi).
// LUU Y: hp/hpMax la writeInt (4 byte) vi client doc readInt; gender ghi 2 LAN;
// truong 15 la hang -1 (0xFF tren day).
//
// Tra nil khi ten dai qua muc writeUTF chiu duoc (65535 byte sau ma hoa) —
// truong hop duy nhat khong dung duoc goi. Nguoi goi tu quyet xu ly.
func goiXuatHienDayDu(id int64, ten string, cap byte, hp, hpToiDa int32, x, y int16, hd HinhDangHien) *protocol.Message {
	w := protocol.NewWriter()
	w.Int(int32(id))                   // 1  playerId
	w.Int(-1)                          // 2  clanId (chua co bang hoi)
	w.Byte(cap)                        // 3  level (= cap nhan vat)
	w.Bool(false)                      // 4  co co dinh (writeBoolean false)
	w.Byte(hd.TypePk)                  // 5  typePk
	w.Byte(byte(hd.Lang))              // 6  gender (= id lang)
	w.Byte(byte(hd.Lang))              // 7  gender ghi LAN 2
	w.Short(hd.Outfit[0])              // 8  head
	if err := w.UTF(ten); err != nil { // 9 name
		return nil
	}
	w.Int(hp)                  // 10 hp (writeInt vi cn.readInt=true)
	w.Int(hpToiDa)             // 11 hpMax
	w.Short(hd.Outfit[1])      // 12 body
	w.Short(hd.Outfit[2])      // 13 leg
	w.Byte(byte(hd.Outfit[3])) // 14 flagbag (1 byte tren day)
	w.Byte(0xFF)               // 15 hang -1 (writeByte(-1) = 0xFF tren day)
	w.Short(x)                 // 16 x
	w.Short(y)                 // 17 y
	w.Short(0)                 // 18 effBuffHp
	w.Short(0)                 // 19 effBuffMp
	w.Byte(0)                  // 20 numEff
	w.Byte(0)                  // 21 idSpaceShip
	w.Byte(0)                  // 22 isMonkey
	w.Short(0)                 // 23 mount
	w.Byte(0)                  // 24 cFlag
	w.Byte(0)                  // 25 hang 0
	w.Short(hd.Outfit[4])      // 26 aura
	w.Byte(byte(hd.Outfit[5])) // 27 effFront (1 byte tren day)
	w.Short(0)                 // 28 hat
	return protocol.NewMessage(lenhXuatHien, w)
}

// goiXuatHien la VO BOC MONG cua goiXuatHienDayDu cho NGUOI CHOI: chi dien
// tham so tu *store.NhanVat, khong them logic nao. Nguoi choi chua co he trang
// bi nen outfit ngoai kieu toc va typePk deu la 0 — Y HET ban KH5.
func goiXuatHien(nv *store.NhanVat) (*protocol.Message, error) {
	m := goiXuatHienDayDu(
		nv.ID, nv.Ten, byte(nv.Cap),
		nv.HpToiDa, nv.HpToiDa,
		int16(nv.ViTriX), int16(nv.ViTriY),
		HinhDangHien{
			Lang:   nv.Lang,
			Outfit: [6]int16{nv.KieuToc, 0, 0, 0, 0, 0}, // head, body, leg, flagBag, aura, effFront
			TypePk: 0,
		},
	)
	if m == nil {
		return nil, fmt.Errorf("dung goi -5 that bai: ten nhan vat %d dai %d byte, vuot gioi han writeUTF", nv.ID, len(nv.Ten))
	}
	return m, nil
}

// (Ban goc con gui kem -103 (co) va -8 (neu dang chet). KH5 BO ca hai vi nhan
// vat luon song va chua co co bang hoi — them lai o KH sau khi co he thong tuong ung.)
