// chieuthuc.go — danh muc chieu Naruto cho goi -28/7 (Ke hoach 6, Task 1).
//
// Khuon byte lay NGUYEN VAN dac ta §3.7 (data/DataGame.java:150-209 ban Java):
// mot MauChieu (SkillTemplate) la mot chieu voi nhieu CapChieu (Skill, moi cap
// mot dong so lieu). Xem ghiMauChieu() de biet thu tu byte chinh xac.
package game

import "nroserver/internal/protocol"

// lenhChonChieu: goi LEN client gui de chon mot chieu de dung (dac ta §3.2,
// bang §3.1). Than goi: `short templateId` — la MauChieu.Id (template.id),
// TUYET DOI KHONG PHAI skillId (hai truong khac kieu O KHAC O CHO, xem bang
// §3.2 cua dac ta — nham la loi im lang, khong triệu chung).
const lenhChonChieu int8 = 34

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

// ⚠️ KH8 TASK 3 — BANG SO LIEU GO DA BI XOA KHOI TEP NAY.
//
// Truoc task 3, tep nay con `taoCapChieu` + `xayDungChieuTheoLang` +
// `chieuTheoLang`: 5 lang x 4 chieu x 3 cap SINH RA BANG CODE. Nay nguon duy
// nhat la kho-noi-dung/chieu.yaml (60 cap tuong minh), doc qua anh chup —
// xem khochieu.go. Tep nay chi con KHUON BYTE (ghiMauChieu) va cac kieu.
//
// Dung "khoi phuc lai cho tien": hai ban so lieu la hai ban se lech nhau, va
// lech im lang.

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

// Ghi chu ve TEN va MO TA chieu: chung tung nam o day (tenChieuTheoLang,
// moTaTheoLoai). Nay nam trong kho-noi-dung/chieu.yaml, truong `ten` / `mo_ta`
// cua tung chieu — sua ten chieu la sua file YAML, khong sua code Go.
//
// Hai cua tra cuu ChieuCuaLang / TraCapChieu da CHUYEN sang khochieu.go (chung
// doc anh chup kho noi dung, khong con doc bang Go nao).
