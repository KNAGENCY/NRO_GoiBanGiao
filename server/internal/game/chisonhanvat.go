package game

import (
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// lenhChiSoNhanVat: goi -42 mang toan bo chi so cua chinh nguoi choi. Thieu goi
// nay thi client khong co nguon nao khac de biet mau — thanh mau va thanh chakra
// tren HUD nam im o 0, va moi tinh toan hien thi lien quan den mau deu sai.
const lenhChiSoNhanVat int8 = -42

// Chi so mac dinh cho nhan vat chua duoc gan gi. Khong de 0: client coi nhan vat
// 0 mau la da chet.
const (
	hpMacDinh     int32 = 100
	chakraMacDinh int32 = 50
	tocDoMacDinh  byte  = 8
)

// goiChiSoNhanVat dung goi -42. Khuon byte doc tu client Controller.cs:1704.
//
// LUU Y kieu so: client doc bang readIntToLongDQT(), ma ham do nhin hai co
// ModFunc.isReadInt / isReadDouble de quyet dinh doc int, long hay double.
// Ban client hien tai dat isReadInt = true nen o day ghi INT 4 byte. Doi co ben
// client ma khong doi o day thi CA GOI lech, khong chi mot truong.
//
// Ban goi nay dung o cCriticalGoc. Phan sau (tlDef, tlPst, tlNeDon, …) client
// boc trong try/catch va tu dat 0 khi thieu, nen dung o day la hop le. Nhung
// khong duoc ghi THUA byte — thua thi client doc lech sang truong khac.
func goiChiSoNhanVat(nv *store.NhanVat) *protocol.Message {
	hpDay := nv.HpToiDa
	if hpDay <= 0 {
		hpDay = hpMacDinh
	}
	chakraDay := nv.ChakraToiDa
	if chakraDay <= 0 {
		chakraDay = chakraMacDinh
	}

	// Mau hien tai khong duoc vuot mau day, neu khong thanh mau ve tran khung.
	hp := nv.HpHienTai
	if hp <= 0 || hp > hpDay {
		hp = hpDay
	}
	chakra := nv.ChakraHienTai
	if chakra <= 0 || chakra > chakraDay {
		chakra = chakraDay
	}

	// Suc danh suy tu suc manh, chua co cong thuc rieng (KH8 se do lai bang so
	// lieu thuc nghiem). Giu mot moc toi thieu de danh quai co tac dung.
	satThuong := int32(nv.SucManh / 100)
	if satThuong < 10 {
		satThuong = 10
	}

	w := protocol.NewWriter()

	w.Int(hpDay)     // cHPGoc
	w.Int(chakraDay) // cMPGoc
	w.Int(satThuong) // cDamGoc

	w.Int(hpDay)     // cHPFull
	w.Int(chakraDay) // cMPFull
	w.Int(hp)        // cHP
	w.Int(chakra)    // cMP

	w.Byte(tocDoMacDinh) // cspeed — 0 thi client coi nhu khong di chuyen duoc
	w.Byte(0)            // hpFrom1000TiemNang
	w.Byte(0)            // mpFrom1000TiemNang
	w.Byte(0)            // damFrom1000TiemNang

	w.Int(satThuong) // cDamFull
	w.Int(0)         // cDefull
	w.Byte(0)        // cCriticalFull

	w.Long(nv.KinhNghiem) // cTiemNang
	w.Short(0)            // expForOneAdd
	w.Short(0)            // cDefGoc
	w.Byte(0)             // cCriticalGoc

	return protocol.NewMessage(lenhChiSoNhanVat, w)
}
