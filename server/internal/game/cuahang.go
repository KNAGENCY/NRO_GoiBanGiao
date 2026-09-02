// cuahang.go — CUA HANG NPC: danh muc hang + bo dung goi XUONG -44 (loaiShop=0).
// KH7 Task 8. CHUA noi vao login.go — task nay chi dung goi va danh muc.
//
// Nguon giai ma khuon byte (PHA A GATE):
// docs/superpowers/research/2026-09-03-kh7-giai-ma-44.md
//
// BON SU THAT DE NHO — sai la HONG NGAM, client nuot loi khong bao gi:
//
//  1. Client goi ItemTemplates.get(templateId) roi truy cap .template.name NGAY
//     (Controller.cs:2627-2629). Gui mot MauID KHONG co trong DanhSachMauVatPham
//     = NullReferenceException -> catch rong (Controller.cs:4600-4602) nuot ->
//     cua hang hien DO DANG. Vi vay goiCuaHang KIEM HET roi moi ghi byte dau tien.
//  2. Hai o gia la [vang][ngoc], DEU LUON ghi (int 4 byte), o khong dung ghi 0.
//     Ten phia client goi o thu hai la buyGold nhung no la NGOC — dung chep ten.
//  3. optionId di xuong bang MOT byte khong dau (ModFunc.isShortOptionTemp=false)
//     va client KHONG kiem bien mang mau (ItemOption.cs:40) -> optionId >= so mau
//     la IndexOutOfRange, cung bi nuot.
//  4. Loai shop 1/2/3/4/8 co THAN GOI KHAC HAN (loai 2 con chen maxPage sau ten
//     tab). File nay CHI dung loai 0. Loai khac -> ham dung rieng, khong dung co.
package game

import (
	"fmt"
	"math"

	"nroserver/internal/protocol"
)

// lenhMoCuaHang: ma lenh goi XUONG mo cua hang NPC (Controller.cs case -44).
const lenhMoCuaHang int8 = -44

// loaiShopThuong: gia tri o dau than goi. Java ShopService.NORMAL_SHOP = 0.
const loaiShopThuong byte = 0

// Tran cung — lay tu PHA A, khong duoc noi.
const (
	// tranSoTab / tranSoMonMoiTab: client doc readUnsignedByte (<=255) nhung ban
	// Java ghi writeByte CO DAU. Giu trong 127 de khong bao gio phai doan xem
	// client xu ly 128..255 the nao.
	tranSoTab       = 127
	tranSoMonMoiTab = 127
	// tranSoOptionMoiMon: soOption la mot byte khong dau.
	tranSoOptionMoiMon = 255
	// tranParamOption: client doc ReadParamDQT -> readUnsignedShort (isReadInt=true).
	tranParamOption = 65535
)

// LoaiTien noi mot mon hang lay tien gi. Ban 1 CHI hai loai (bay L10 cua dac ta:
// ban Java nhet ngoc/hong ngoc/phieu chung mot o, khong chep lai).
type LoaiTien uint8

const (
	// TienVang — o gia thu nhat (client: buyCoin).
	TienVang LoaiTien = iota
	// TienNgoc — o gia thu hai (client: buyGold, nhung that ra la NGOC).
	TienNgoc
)

// OptionMon la mot dong option gan san vao mon hang bay ban.
type OptionMon struct {
	// OptionID = CHI SO trong DanhSachMauOption(). Di xuong bang 1 byte khong dau.
	OptionID uint8
	// Param di xuong bang unsigned short: 0..65535 (bay L3 — so lon hon phai doi
	// DON VI, khong duoc tran so).
	Param uint16
}

// MonHang la mot mon bay ban trong mot tab cua hang.
//
// VI SAO tach MauID khoi Gia: mau vat pham (ten, icon, chi so) la du lieu CHUNG
// cua ca game; don gia la du lieu RIENG cua tung cua hang. Cung mot mau co the
// ban gia khac nhau o hai lang. Gop lam mot la khong the.
type MonHang struct {
	// MauID = chi so trong DanhSachMauVatPham(). PHAI ton tai (su that 1).
	MauID int16
	// Gia >= 0 va <= math.MaxInt32 (o la int 4 byte co dau).
	Gia int32
	// Tien: vang hay ngoc. O con lai ghi 0.
	Tien LoaiTien
	// Options: option gan san. Rong = mon do tho.
	Options []OptionMon
	// LaMoi: client ve nhan "moi" (Controller.cs:2675-2676).
	LaMoi bool
}

// TabCuaHang la mot the trong cua so cua hang.
//
// Ten co the chua '\n' de tach thanh TOI DA 2 dong (Res.split, Controller.cs:2598).
// Nhieu hon 2 dong: client chi nhan mang do dai 1 hoac 2, do dai khac bi BO QUA
// va tab hien ten RONG — khong loi. Vi vay goiCuaHang kiem so dong.
type TabCuaHang struct {
	Ten string
	Mon []MonHang
}

// CuaHang la mot cua hang NPC hoan chinh.
type CuaHang struct {
	// ID: khoa tra cuu phia may chu (NPC nao mo cua hang nao). KHONG di xuong
	// client — gói -44 khong co o nao chua id cua hang.
	ID int32
	// Ten: chi de doc log/quan tri, cung khong di xuong client.
	Ten  string
	Tabs []TabCuaHang
}

// ---------------------------------------------------------------------------
// Danh muc cua hang mau
// ---------------------------------------------------------------------------

// cuaHangTapHoaLangLa: cua hang mau toi thieu de chay dau-cuoi. Cac mon lay tu
// danh muc mau task 2 (mauvatpham.go): 0/1 la thuoc, 2/3/4 la trang bi khoi dau.
//
// GIA DINH con no bang chung: don gia duoi day la KHOI DAU do nguoi viet dat,
// KHONG trich tu ban Java (ban Java la game Dragon Ball). Can bang kinh te
// thuoc KH8, khong phai task nay.
var cuaHangTapHoaLangLa = CuaHang{
	ID: 1, Ten: "Tap hoa Lang La",
	Tabs: []TabCuaHang{
		{
			Ten: "Thuoc",
			Mon: []MonHang{
				{MauID: 0, Gia: 100, Tien: TienVang},
				{MauID: 1, Gia: 100, Tien: TienVang},
			},
		},
		{
			Ten: "Trang bi",
			Mon: []MonHang{
				{MauID: 2, Gia: 1000, Tien: TienVang, LaMoi: true},
				{MauID: 3, Gia: 1000, Tien: TienVang, LaMoi: true},
				{MauID: 4, Gia: 5, Tien: TienNgoc},
			},
		},
	},
}

var danhSachCuaHang = []CuaHang{cuaHangTapHoaLangLa}

// DanhSachCuaHang tra toan bo cua hang NPC dang co.
func DanhSachCuaHang() []CuaHang { return danhSachCuaHang }

// TimCuaHang tra cua hang theo ID, nil neu khong co. MOI cho tra cuu phai di qua
// day — id den tu client la khong tin duoc.
func TimCuaHang(id int32) *CuaHang {
	for i := range danhSachCuaHang {
		if danhSachCuaHang[i].ID == id {
			return &danhSachCuaHang[i]
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Bo dung goi -44
// ---------------------------------------------------------------------------

// goiCuaHang dung goi XUONG -44 cho mot cua hang NPC (loaiShop = 0).
//
// Khuon than goi (theo dung thu tu client doc, Controller.cs:2572-2685):
//
//	[loaiShop=0 byte][soTab byte]
//	( [tenTab UTF][soMon byte]
//	  ( [mauID short][giaVang int][giaNgoc int]
//	    [soOption byte]( [optionId byte][param ushort] )*
//	    [laMoi byte][coPart byte]( [head][body][leg][-1] short x4 neu coPart )
//	  )*
//	)*
//
// KIEM HET TRUOC, GHI SAU (§1.10 / bay L13): mot gói thieu byte lam client nem
// loi trong catch rong -> cua hang hien do dang, khong log, khong dau vet.
func goiCuaHang(ch *CuaHang) (*protocol.Message, error) {
	if ch == nil {
		return nil, fmt.Errorf("cua hang nil")
	}
	if err := kiemCuaHang(ch); err != nil {
		return nil, fmt.Errorf("cua hang %d (%s): %w", ch.ID, ch.Ten, err)
	}

	w := protocol.NewWriter()
	w.Byte(loaiShopThuong)
	w.Byte(byte(len(ch.Tabs)))
	for _, tab := range ch.Tabs {
		if err := w.UTF(tab.Ten); err != nil {
			return nil, fmt.Errorf("cua hang %d ten tab %q khong ghi duoc: %w", ch.ID, tab.Ten, err)
		}
		w.Byte(byte(len(tab.Mon)))
		for _, mon := range tab.Mon {
			w.Short(mon.MauID)
			if mon.Tien == TienVang {
				w.Int(mon.Gia)
				w.Int(0)
			} else {
				w.Int(0)
				w.Int(mon.Gia)
			}
			w.Byte(byte(len(mon.Options)))
			for _, o := range mon.Options {
				w.Byte(o.OptionID)
				// param di xuong 2 byte; client doc readUnsignedShort.
				w.Short(int16(o.Param))
			}
			if mon.LaMoi {
				w.Byte(1)
			} else {
				w.Byte(0)
			}
			// coPart: chi mau Loai == 5 moi ve manh len nguoi (ShopService:198-205).
			// TimMauVatPham da duoc kiem khac nil trong kiemCuaHang.
			mau := TimMauVatPham(mon.MauID)
			if mau.Loai == 5 {
				w.Byte(1)
				w.Short(mau.Phan)
				w.Short(mau.Phan)
				w.Short(mau.Phan)
				w.Short(-1) // o thu tu (bag) ban Java LUON ghi -1
			} else {
				w.Byte(0)
			}
		}
	}
	return protocol.NewMessage(lenhMoCuaHang, w), nil
}

// kiemCuaHang kiem TOAN BO du lieu truoc khi ghi mot byte nao. Tra loi dau tien
// gap phai, kem duong di (tab nao, mon thu may) de sua duoc ngay.
func kiemCuaHang(ch *CuaHang) error {
	if len(ch.Tabs) > tranSoTab {
		return fmt.Errorf("%d tab, toi da %d", len(ch.Tabs), tranSoTab)
	}
	soMauOption := len(DanhSachMauOption())
	for iTab, tab := range ch.Tabs {
		if err := kiemTenTab(tab.Ten); err != nil {
			return fmt.Errorf("tab %d: %w", iTab, err)
		}
		if len(tab.Mon) > tranSoMonMoiTab {
			return fmt.Errorf("tab %d co %d mon, toi da %d", iTab, len(tab.Mon), tranSoMonMoiTab)
		}
		for iMon, mon := range tab.Mon {
			if TimMauVatPham(mon.MauID) == nil {
				return fmt.Errorf(
					"tab %d mon %d: mau vat pham %d KHONG TON TAI (client se nem NullReference roi nuot)",
					iTab, iMon, mon.MauID)
			}
			if mon.Gia < 0 {
				return fmt.Errorf("tab %d mon %d: gia am (%d)", iTab, iMon, mon.Gia)
			}
			if mon.Tien != TienVang && mon.Tien != TienNgoc {
				return fmt.Errorf("tab %d mon %d: loai tien %d khong hop le", iTab, iMon, mon.Tien)
			}
			if len(mon.Options) > tranSoOptionMoiMon {
				return fmt.Errorf("tab %d mon %d: %d option, toi da %d",
					iTab, iMon, len(mon.Options), tranSoOptionMoiMon)
			}
			for iOpt, o := range mon.Options {
				if int(o.OptionID) >= soMauOption {
					return fmt.Errorf(
						"tab %d mon %d option %d: id %d >= %d mau option (client se nem IndexOutOfRange roi nuot)",
						iTab, iMon, iOpt, o.OptionID, soMauOption)
				}
				if int(o.Param) > tranParamOption {
					return fmt.Errorf("tab %d mon %d option %d: param %d vuot %d",
						iTab, iMon, iOpt, o.Param, tranParamOption)
				}
			}
		}
	}
	return nil
}

// kiemTenTab: client tach ten bang '\n' va CHI nhan mang do dai 1 hoac 2
// (Controller.cs:2598-2610). Ba dong tro len -> tab hien ten RONG, khong bao loi.
func kiemTenTab(ten string) error {
	if ten == "" {
		return fmt.Errorf("ten tab rong")
	}
	soDong := 1
	for _, r := range ten {
		if r == '\n' {
			soDong++
		}
	}
	if soDong > 2 {
		return fmt.Errorf("ten tab %q co %d dong, client chi nhan toi da 2", ten, soDong)
	}
	return nil
}

// kiemMoiCuaHang kiem toan bo danh muc cua hang. Goi trong test khoi dong de
// mot MauID lac khong bao gio ra toi nguoi choi.
func kiemMoiCuaHang() error {
	for i := range danhSachCuaHang {
		if err := kiemCuaHang(&danhSachCuaHang[i]); err != nil {
			return fmt.Errorf("cua hang %d (%s): %w",
				danhSachCuaHang[i].ID, danhSachCuaHang[i].Ten, err)
		}
	}
	return nil
}

// giaToiDa la tran o gia (int 4 byte co dau). De o day cho de tra cuu; kiemCuaHang
// khong can kiem vi Gia da la int32.
const giaToiDa = math.MaxInt32
