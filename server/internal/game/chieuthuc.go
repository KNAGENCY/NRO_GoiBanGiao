// chieuthuc.go — danh muc chieu Naruto cho goi -28/7 (Ke hoach 6, Task 1).
//
// Khuon byte lay NGUYEN VAN dac ta §3.7 (data/DataGame.java:150-209 ban Java):
// mot MauChieu (SkillTemplate) la mot chieu voi nhieu CapChieu (Skill, moi cap
// mot dong so lieu). Xem ghiMauChieu() de biet thu tu byte chinh xac.
package game

import "nroserver/internal/protocol"

// ghiChuChieu la truong "ghiChu" ghi cho MOI chieu. Ban Java ghi cung
// "Ngoc Rong Kame" (dau vet thuong hieu cu) — spec §3.7 diem 3 yeu cau ghi
// chuoi rong de xoa dau vet do, khong phat lai ten thuong hieu cu.
const ghiChuChieu = ""

// MauChieu la mot SkillTemplate: mot chieu voi CapToiDa cap do (CacCap).
// Khop dac ta §3.7 — moi truong tuong ung dung 1 byte/UTF/short trong khuon,
// TRU NhanLucCan trong CapChieu la long 8 byte (xem CapChieu).
type MauChieu struct {
	Id            byte // templateId — BAT DAU TU 1, xem ghi chu o ChieuCuaLang
	Ten           string
	CapToiDa      byte // maxPoint: so cap toi da cua chieu (KH6 luon la 3)
	KieuTonChakra byte // manaUseType: 0 = so chakra co dinh, 1 = % chakra toi da
	Loai          byte // type: 1 = danh (can chien/tam xa), 2 = buff nguoi khac (hoi mau), 3 = tu than
	IconId        int16
	MoTa          string // damInfo: mo ta sat thuong/hieu qua
	CacCap        []CapChieu
}

// CapChieu la mot Skill: so lieu cua MOT cap trong mot chieu.
// Thu tu truong PHAI khop dac ta §3.7 tung byte mot.
type CapChieu struct {
	SkillId    int16
	Diem       byte  // point: so thu tu cap (1, 2, 3, ...)
	NhanLucCan int64 // powRequire — LONG 8 BYTE, truong duy nhat trong ca khuon
	// dung 8 byte; nham thanh int (4 byte) se lech byte moi cap ve sau (bay §3.7.1).
	ChakraDung   int16 // manaUse
	HoiChieuMs   int32 // coolDown, mili-giay
	Dx, Dy       int16 // tam danh (client doc ca hai, dung chung 1 gia tri tam)
	MaxFight     byte
	SatThuong    int16  // damage — cung truong nay ban Java tai su dung lam "luong hoi mau" cho chieu type=2
	Gia          int16  // price
	ThemThongTin string // moreInfo
}

// ghiMauChieu ghi MOT MauChieu (template + toan bo cap) vao w, dung khuon
// dac ta §3.7. Tach rieng ham nay (thay vi nhoi thang vao goiCatalogSkill)
// de test byte-chinh-xac (chieuthuc_test.go) ghim duoc dung MOT chieu, MOT
// cap ma khong can dung toan bo du lieu 20-chieu that.
func ghiMauChieu(w *protocol.Writer, m MauChieu) error {
	w.Byte(m.Id)
	if err := w.UTF(m.Ten); err != nil {
		return err
	}
	w.Byte(m.CapToiDa)
	w.Byte(m.KieuTonChakra)
	w.Byte(m.Loai)
	w.Short(m.IconId)
	if err := w.UTF(m.MoTa); err != nil {
		return err
	}
	if err := w.UTF(ghiChuChieu); err != nil { // luon chuoi rong, xem ghi chu o const
		return err
	}
	w.Byte(byte(len(m.CacCap)))
	for _, c := range m.CacCap {
		w.Short(c.SkillId)
		w.Byte(c.Diem)
		w.Long(c.NhanLucCan) // 8 byte — xem chu thich truong NhanLucCan
		w.Short(c.ChakraDung)
		w.Int(c.HoiChieuMs)
		w.Short(c.Dx)
		w.Short(c.Dy)
		w.Byte(c.MaxFight)
		w.Short(c.SatThuong)
		w.Short(c.Gia)
		if err := w.UTF(c.ThemThongTin); err != nil {
			return err
		}
	}
	return nil
}

// taoCapChieu dung 3 cap tang dan cho mot chieu tu vai con so goc. Gia tri chi
// la KHOI DAU can chinh o KH8 (kho YAML) — nhiem vu Task 1 la khuon byte dung,
// khong phai can bang so liệu cuoi cung.
func taoCapChieu(skillIdGoc, dx, dy, manaGoc, dameGoc, giaGoc int16) []CapChieu {
	cac := make([]CapChieu, 0, 3)
	for cap := int16(1); cap <= 3; cap++ {
		cac = append(cac, CapChieu{
			SkillId:      skillIdGoc + cap - 1,
			Diem:         byte(cap),
			NhanLucCan:   int64(cap) * 1000,
			ChakraDung:   manaGoc + (cap-1)*10,
			HoiChieuMs:   1000 + int32(cap-1)*300,
			Dx:           dx,
			Dy:           dy,
			MaxFight:     1,
			SatThuong:    dameGoc + (cap-1)*20,
			Gia:          giaGoc + (cap-1)*50,
			ThemThongTin: "",
		})
	}
	return cac
}

// Bon vai tro bat buoc theo dac ta §4.3, dung cho ca 5 lang:
//
//	slot 0: can chien (type=1, tam gan 100)
//	slot 1: tam xa theo he (type=1, tam xa 300)
//	slot 2: hoi mau dong doi (type=2, tam 300) — BAT BUOC, goi -60 dang hoi mau
//	        (Task 10) chi sinh tu chieu type=2, thieu no goiHoiMauNguoiChoi
//	        khong co duong chay that (xem dac ta §4.3 va brief Task 1).
//	slot 3: tu than (type=3, khong tam — khong phat goi chien dau)
const (
	loaiDanh    byte  = 1
	loaiBuff    byte  = 2
	loaiTuThan  byte  = 3
	tamCanChien int16 = 100
	tamXa       int16 = 300
	tamTuThan   int16 = 0
)

// tenChieuTheoLang: 4 ten chieu (can chien, tam xa, hoi mau, tu than) cho tung
// lang trong 5 lang, dat theo he nguyen to cua lang (§4.2 dac ta).
var tenChieuTheoLang = [5][4]string{
	{"Cước Liên Hoàn", "Hoả Độn: Hào Hoả Cầu", "Y Liệu Nhẫn Thuật", "Tụ Chakra"},                    // Làng Lá — Hoả
	{"Trảm Phong Cước", "Phong Độn: Cuồng Phong Trảm", "Sa Thổ Y Liệu Thuật", "Tụ Chakra Sa Mạc"},   // Làng Cát — Phong
	{"Thuỷ Ảnh Trảm", "Thuỷ Độn: Đại Thuỷ Long Đạn", "Sương Mù Y Liệu Thuật", "Tụ Chakra Sương Mù"}, // Làng Sương Mù — Thuỷ
	{"Nham Thạch Quyền", "Thổ Độn: Nham Thạch Đạn", "Nham Y Liệu Thuật", "Tụ Chakra Nham Thạch"},    // Làng Đá — Thổ
	{"Lôi Ảnh Cước", "Lôi Độn: Cầu Sét", "Vân Y Liệu Thuật", "Tụ Chakra Vân"},                       // Làng Mây — Lôi
}

var moTaTheoLoai = [4]string{
	"Đòn cận chiến gây sát thương vật lý",
	"Đòn tấn công tầm xa theo hệ nguyên tố của làng",
	"Hồi máu cho đồng đội trong tầm",
	"Tự thân hồi chakra, không nhắm mục tiêu",
}

// chieuTheoLang: du lieu that cua 5 lang x 4 chieu x 3 cap = 20 MauChieu, 60
// CapChieu. Dung ham (khong ghi tay 20 struct-literal) de giam sai sot lap
// lai, nhung gia tri sinh ra van xac dinh (deterministic) — cung tham so vao
// luon ra cung ket qua, khong ngau nhien.
var chieuTheoLang = xayDungChieuTheoLang()

func xayDungChieuTheoLang() [5][]MauChieu {
	var ketQua [5][]MauChieu
	for lang := 0; lang < 5; lang++ {
		loaiCuaSlot := [4]byte{loaiDanh, loaiDanh, loaiBuff, loaiTuThan}
		tamCuaSlot := [4]int16{tamCanChien, tamXa, tamXa, tamTuThan}
		danhSach := make([]MauChieu, 0, 4)
		for slot := 0; slot < 4; slot++ {
			// templateId BAT DAU TU 1 (slot+1, khong bao gio la 0) — nhu dac ta
			// §3.7 diem 2 yeu cau, de tranh nhanh nhoi 2 cap gia 105/106 cua ban
			// Java khi templateId == 0.
			templateId := byte(slot + 1)
			// skillId goc: moi mot trong 20 chieu chiem 3 skillId lien tiep
			// (3 cap), khong trung giua cac chieu/lang.
			skillIdGoc := int16(1 + (lang*4+slot)*3)
			iconId := int16(1 + lang*4 + slot)
			manaGoc := int16(10 + slot*5)
			dameGoc := int16(20 + slot*10)
			giaGoc := int16(50 + slot*20)

			danhSach = append(danhSach, MauChieu{
				Id:            templateId,
				Ten:           tenChieuTheoLang[lang][slot],
				CapToiDa:      3,
				KieuTonChakra: 0,
				Loai:          loaiCuaSlot[slot],
				IconId:        iconId,
				MoTa:          moTaTheoLoai[slot],
				CacCap:        taoCapChieu(skillIdGoc, tamCuaSlot[slot], tamCuaSlot[slot], manaGoc, dameGoc, giaGoc),
			})
		}
		ketQua[lang] = danhSach
	}
	return ketQua
}

// ChieuCuaLang tra ve 4 MauChieu cua lang idLang (0..4, khop Lang.ID trong
// lang.go). idLang ngoai khoang tra ve nil (khong co lang do).
func ChieuCuaLang(idLang byte) []MauChieu {
	if int(idLang) >= len(chieuTheoLang) {
		return nil
	}
	return chieuTheoLang[idLang]
}
