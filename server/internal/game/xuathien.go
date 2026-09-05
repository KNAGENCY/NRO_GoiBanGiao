package game

import (
	"fmt"
	"log/slog"

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
	w.Int(int32(id))  // 1  playerId
	w.Int(-1)         // 2  clanId (chua co bang hoi)
	w.Byte(cap)       // 3  level (= cap nhan vat)
	w.Bool(false)     // 4  co co dinh (writeBoolean false)
	w.Byte(hd.TypePk) // 5  typePk
	// Client doc HAI byte "gender" khac vai (Controller.readCharInfo:6376-6377):
	//   6 -> nClass = nClasss[byte]  : LOP/LANG (ten lang, bo chieu) = id lang
	//   7 -> cgender                  : than/chan mac dinh, hieu ung theo hanh tinh
	// Byte 7 luon 0 — xem gioiTinhClient (hinhdang.go): client chi biet 3 hanh
	// tinh, bo do Naruto nam o slot mac dinh cua gender 0, lang 3/4 lam client
	// tra bang 3 o roi vo.
	w.Byte(byte(hd.Lang))              // 6  lop/lang (chi so nClasss[])
	w.Byte(gioiTinhClient(hd.Lang))    // 7  cgender phia client
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
// tham so tu *store.NhanVat. Hinh dang (dau/than/chan) lay tu hinhDangNguoiChoi
// (hinhdang.go) — truoc 05/09 than/chan ghi 0 ("chua co he trang bi") va client
// khac ve part 0 (mot part DAU 3 khung) lam than -> vo khi ve.
func goiXuatHien(nv *store.NhanVat) (*protocol.Message, error) {
	m := goiXuatHienDayDu(
		nv.ID, nv.Ten, byte(nv.Cap),
		nv.HpToiDa, nv.HpToiDa,
		int16(nv.ViTriX), int16(nv.ViTriY),
		hinhDangNguoiChoi(nv),
	)
	if m == nil {
		return nil, fmt.Errorf("dung goi -5 that bai: ten nhan vat %d dai %d byte, vuot gioi han writeUTF", nv.ID, len(nv.Ten))
	}
	return m, nil
}

// goiXuatHienBoss la vo boc cua goiXuatHienDayDu cho BOSS (KH6b Task 7): dien
// tham so tu Boss (id/hp/toa do — thu doi theo tung tick) + PhaBoss dang chay
// (ten/cap/lang/outfit/typePk — thu doi theo tung PHA).
//
// Tra nil khi KHONG dung duoc goi, va nguoi goi PHAI kiem nil truoc khi
// XepGui: xep mot con tro nil vao hang gui la panic hoac mot goi rong tren
// day. Ba truong hop tra nil:
//   - b == nil (con tro lac vao tu noi goi);
//   - b.PhaHienTaiData() == nil — PhaHienTai nam ngoai [0, len(DuLieu)), tuc
//     boss chua vao pha nao / da qua pha cuoi / du lieu rong;
//   - ten pha dai qua muc writeUTF chiu duoc (goiXuatHienDayDu tra nil).
//
// TypePk: lay tu du lieu pha, nhung neu du lieu de 0 thi THAY BANG TypePkBoss
// va ghi CANH BAO — de 0 tren day thi client (Char.cs paintHp) thoat som,
// khong ve thanh mau, nguoi choi danh boss ma khong biet con bao nhieu mau.
// Canh bao la bat buoc: thay the am tham se giau mot lo hong du lieu that.
func goiXuatHienBoss(b *Boss) *protocol.Message {
	if b == nil {
		return nil
	}
	pha := b.PhaHienTaiData()
	if pha == nil {
		slog.Error("dung goi -5 boss that bai: khong co pha dang chay",
			"idBoss", b.Id, "phaHienTai", b.PhaHienTai, "soPha", len(b.DuLieu))
		return nil
	}

	typePk := pha.TypePk
	if typePk == 0 {
		slog.Warn("du lieu pha boss thieu TypePk (=0), dung TypePkBoss thay the — client se khong ve thanh mau neu de 0",
			"idBoss", b.Id, "phaHienTai", b.PhaHienTai, "ten", pha.Ten)
		typePk = TypePkBoss
	}

	m := goiXuatHienDayDu(
		b.Id, pha.Ten, pha.Cap,
		b.Hp, b.HpToiDa,
		b.ViTriX, b.ViTriY,
		HinhDangHien{
			Lang:   pha.Lang,
			Outfit: pha.Outfit,
			TypePk: typePk,
		},
	)
	if m == nil {
		slog.Error("dung goi -5 boss that bai: ten pha vuot gioi han writeUTF",
			"idBoss", b.Id, "phaHienTai", b.PhaHienTai, "soByteTen", len(pha.Ten))
	}
	return m
}

// (Ban goc con gui kem -103 (co) va -8 (neu dang chet). KH5 BO ca hai vi nhan
// vat luon song va chua co co bang hoi — them lai o KH sau khi co he thong tuong ung.)
