// mauvatpham.go — DANH MUC MAU vat pham + bo dung goi -28/8 subtype 0 va 1
// (KH7 Task 2). Du lieu THUAN Go, cung cach noidungboss.go da lam voi boss.
//
// VI SAO DU LIEU NAM TRONG Go CHU KHONG PHAI catalog.json: catalog.json bi ghim
// SHA trong asset-lock.json (catalog.go, lockHashKhop). Them mau vat pham vao do
// la lam do test khop artifact va phai ky lai mot artifact da ky. Bay §7.6 cua
// dac ta khuyen nghi duong (a) — de trong file Go rieng. Day la duong da chon.
//
// ⚠️ MOI CON SO O DAY DEU TRUY DUOC VE MOT NGUON. Nguon giai ma khuon byte:
// docs/superpowers/research/2026-09-03-kh7-giai-ma-28-8.md (PHA A GATE).
//
// BA SU THAT DE NHO — sai la hong NGAM, khong co loi nao bat ra:
//
//  1. id KHONG DI QUA MANG. Client dat id = CHI SO vong lap
//     (Controller.cs:6479 va :6498). Chen mot mau vao GIUA danh sach se doi id
//     cua moi mau sau no — va do la id da nam trong CSDL cua nguoi choi.
//     CHI DUOC THEM VAO CUOI.
//  2. So mau option la MOT BYTE KHONG DAU (myReader.readShortOptionTemp voi
//     ModFunc.isShortOptionTemp == false) -> toi da 255. So mau vat pham la
//     SHORT CO DAU -> toi da 32767. Vuot tran = client dung catalog RONG, im lang.
//  3. Client BAM CUNG mot so id option: 22 -> 6 (param x1000), 23 -> 7
//     (param x1000) trong ItemOption.cs:29-38 va Panel.cs:9889-9911; va
//     iOptionTemplates[id] KHONG kiem bien (ItemOption.cs:40) -> gui id option
//     >= so mau la client nem IndexOutOfRange. Xem §3.1 tai lieu giai ma.
package game

import (
	"fmt"

	"nroserver/internal/protocol"
)

// ---------------------------------------------------------------------------
// 1. Hai kieu du lieu
// ---------------------------------------------------------------------------

// MauOption la MOT dong cua bang mau option. id = CHI SO trong DanhSachMauOption.
//
// Ten chua ky tu '#': client thay '#' bang gia tri param khi ve
// (ItemOption.getOptionString, ItemOption.cs:43-46). Ten khong co '#' thi nguoi
// choi thay dong chu tinh, khong thay con so — khong loi, chi vo nghia.
type MauOption struct {
	Ten  string
	Loai int8 // ItemOptionTemplate.type. Y nghia CHUA giai ma duoc tu client -> 0.
}

// MauVatPham la MOT mau vat pham. id = CHI SO trong DanhSachMauVatPham.
//
// Chin truong dau di XUONG CLIENT dung thu tu nay (khong duoc dao):
// Loai, GioiTinh, Ten, MoTa, Cap, YeuCauSuc, IconID, Phan, NangCapDuoc.
// TranChong la truong CHI DUNG PHIA MAY CHU — KHONG co trong goi.
type MauVatPham struct {
	Loai        int8 // Item.type. 0..5 = trang bi mac tren nguoi (Item.isTypeBody, Item.cs:271-278)
	GioiTinh    int8 // 3 = moi gioi (quy uoc ban goc)
	Ten         string
	MoTa        string
	Cap         int8  // level yeu cau
	YeuCauSuc   int32 // strRequire
	IconID      int16 // chi so bang SmallImage
	Phan        int16 // part — manh hinh khi mac len nguoi; -1 = khong ve len nguoi
	NangCapDuoc bool  // isUpToUp

	// TranChong: toi da bao nhieu cai trong MOT o tui. 1 = khong chong duoc.
	// KHONG di xuong client — client khong co truong nay trong ItemTemplate.
	// Dung cho task 6 (kho vat pham) va task 9 (mua).
	TranChong int16
}

// ---------------------------------------------------------------------------
// 2. Tran cung — lay tu PHA A, khong duoc noi
// ---------------------------------------------------------------------------

const (
	// tranMauOption: client doc so mau bang MOT BYTE KHONG DAU.
	tranMauOption = 255
	// tranMauVatPham: client doc so mau bang SHORT CO DAU (Controller.cs:6495).
	tranMauVatPham = 32767
	// tranThanGoiItem: client goi d.mark(100000) truoc khi doc (Controller.cs:6468).
	tranThanGoiItem = 100000
)

// ---------------------------------------------------------------------------
// 3. Bang mau option — CHI THEM VAO CUOI
// ---------------------------------------------------------------------------

// Chi so 6 va 7 KHONG duoc doi y nghia: client bam cung 22->6 va 23->7 voi
// param x1000 (ItemOption.cs:29-38). Neu 6/7 khong phai "tang HP"/"tang chakra"
// thi mot mon do co option 22/23 se hien sai hoan toan ma khong bao gi.
var danhSachMauOption = []MauOption{
	0: {Ten: "Hồi phục # HP", Loai: 0},
	1: {Ten: "Hồi phục # chakra", Loai: 0},
	2: {Ten: "Tăng # sức đánh", Loai: 0},
	3: {Ten: "Tăng # giáp", Loai: 0},
	4: {Ten: "Tăng # chính xác", Loai: 0},
	5: {Ten: "Tăng # né đòn", Loai: 0},
	6: {Ten: "Tăng # HP", Loai: 0},     // ⚠️ dich cua option id 22 (param x1000)
	7: {Ten: "Tăng # chakra", Loai: 0}, // ⚠️ dich cua option id 23 (param x1000)
}

// DanhSachMauOption tra bang mau option. Chi so = id option gui trong goi vat pham.
func DanhSachMauOption() []MauOption { return danhSachMauOption }

// ---------------------------------------------------------------------------
// 4. Danh muc mau vat pham — KHOI DAU, du de chay dau-cuoi
// ---------------------------------------------------------------------------
//
// Day la "vai mau toi thieu de chay duoc dau-cuoi", dung tinh than §9 dac ta:
// don gia / danh muc hang NPC / ti le roi / chi so option thuoc KH8. Cac con so
// duoi day la KHOI DAU do nguoi viet dat, KHONG phai trich tu ban Java —
// ban Java la game Dragon Ball, day la game Naruto.
//
// GIA DINH con no bang chung (task 13 nghiem thu tay se tra loi):
//   - Loai 0..5 = trang bi mac tren nguoi: SUY RA tu Item.isTypeBody
//     (Item.cs:271-278). Loai 6 chon lam "do tieu hao" vi no NAM NGOAI dai do —
//     chua doi chieu duoc voi mot cho nao trong client noi ro "6 la thuoc".
//   - GioiTinh = 3 nghia la "moi gioi": quy uoc ban goc, chua xac minh o client.
//   - IconID = 0 cho TAT CA: bang SmallImage phia client hien dang RONG
//     (goiCatalogData gui blob "image" count=0, dongbo.go:79-82). Bay §7.7.
//     SmallImage.drawSmallImage KHONG sap khi thieu anh — no goi createImage(id)
//     roi tra ve (SmallImage.cs:151-161) — nen o icon chi TRONG, khong vo game.
//     Sua khi co bo icon di qua quy trinh my thuat + asset lock.
//   - Phan = -1 nghia la "khong ve manh nao len nguoi". Ba mon trang bi duoi day
//     de -1 vi bo manh do vat pham CHUA co trong part.bin da khoa SHA; dat mot
//     so bua se ve nham manh nguoi choi len nguoi.
var danhSachMauVatPham = []MauVatPham{
	0: {
		Loai: 6, GioiTinh: 3,
		Ten: "Thuốc hồi máu nhỏ", MoTa: "Hồi 200 HP. Dùng được khi đang di chuyển.",
		Cap: 1, YeuCauSuc: 0, IconID: 0, Phan: -1, NangCapDuoc: false,
		TranChong: 99,
	},
	1: {
		Loai: 6, GioiTinh: 3,
		Ten: "Thuốc hồi chakra nhỏ", MoTa: "Hồi 200 chakra.",
		Cap: 1, YeuCauSuc: 0, IconID: 0, Phan: -1, NangCapDuoc: false,
		TranChong: 99,
	},
	2: {
		Loai: 0, GioiTinh: 3,
		Ten: "Áo ninja Làng Lá", MoTa: "Áo tập của học viện Làng Lá.",
		Cap: 1, YeuCauSuc: 0, IconID: 0, Phan: -1, NangCapDuoc: true,
		TranChong: 1, // co option -> BAT BUOC don chiec (rang buoc CSDL task 1)
	},
	3: {
		Loai: 1, GioiTinh: 3,
		Ten: "Quần ninja Làng Lá", MoTa: "Quần tập của học viện Làng Lá.",
		Cap: 1, YeuCauSuc: 0, IconID: 0, Phan: -1, NangCapDuoc: true,
		TranChong: 1,
	},
	4: {
		Loai: 2, GioiTinh: 3,
		Ten: "Kunai luyện tập", MoTa: "Kunai cùn dùng để tập ném.",
		Cap: 1, YeuCauSuc: 0, IconID: 0, Phan: -1, NangCapDuoc: true,
		TranChong: 1,
	},
}

// DanhSachMauVatPham tra ca danh muc. Chi so = templateId gui xuong client.
func DanhSachMauVatPham() []MauVatPham { return danhSachMauVatPham }

// TimMauVatPham tra mau theo id, nil neu ngoai mien. MOI cho tra cuu mau phai
// di qua ham nay — id den tu client hoac tu CSDL deu co the ngoai mien (bay §7.11).
func TimMauVatPham(id int16) *MauVatPham {
	if id < 0 || int(id) >= len(danhSachMauVatPham) {
		return nil
	}
	return &danhSachMauVatPham[id]
}

// ---------------------------------------------------------------------------
// 5. Bo dung goi
// ---------------------------------------------------------------------------

// goiCatalogItemOption dung -28/8 subtype 0: bang mau option.
//
//	[8][vsItem][0][soMau: 1 byte KHONG DAU]( [ten UTF][loai byte] )*
func goiCatalogItemOption() (*protocol.Message, error) {
	ds := DanhSachMauOption()
	if len(ds) > tranMauOption {
		return nil, fmt.Errorf(
			"bang mau option %d muc, client chi doc duoc %d (mot byte khong dau)",
			len(ds), tranMauOption)
	}
	w := protocol.NewWriter()
	w.Byte(lenhConCatalogItem) // 8
	w.Byte(vsItem)
	w.Byte(0)             // subtype 0 — bang mau option
	w.Byte(byte(len(ds))) // KHONG DAU: client doc readUnsignedByte
	for i, o := range ds {
		if err := w.UTF(o.Ten); err != nil {
			return nil, fmt.Errorf("mau option %d ten khong ghi duoc: %w", i, err)
		}
		w.Byte(byte(o.Loai))
	}
	return dongGoiItem(w)
}

// goiCatalogItemMau dung -28/8 subtype 1: nap lai TOAN BO mau vat pham.
//
//	[8][vsItem][1][soMau: short CO DAU]
//	( [loai b][gioiTinh b][ten UTF][moTa UTF][cap b][yeuCauSuc int][icon short][phan short][nangCap bool] )*
//
// Client XOA SACH bang mau truoc khi doc (ItemTemplates.clear, Controller.cs:6494)
// nen goi nay la anh chup DAY DU, khong phai ban vá.
func goiCatalogItemMau() (*protocol.Message, error) {
	ds := DanhSachMauVatPham()
	if len(ds) > tranMauVatPham {
		return nil, fmt.Errorf(
			"danh muc %d mau, client chi doc duoc %d (short co dau)",
			len(ds), tranMauVatPham)
	}
	w := protocol.NewWriter()
	w.Byte(lenhConCatalogItem) // 8
	w.Byte(vsItem)
	w.Byte(1) // subtype 1 — nap lai toan bo
	w.Short(int16(len(ds)))
	for i, m := range ds {
		w.Byte(byte(m.Loai))
		w.Byte(byte(m.GioiTinh))
		if err := w.UTF(m.Ten); err != nil {
			return nil, fmt.Errorf("mau vat pham %d ten khong ghi duoc: %w", i, err)
		}
		if err := w.UTF(m.MoTa); err != nil {
			return nil, fmt.Errorf("mau vat pham %d mo ta khong ghi duoc: %w", i, err)
		}
		w.Byte(byte(m.Cap))
		w.Int(m.YeuCauSuc)
		w.Short(m.IconID)
		w.Short(m.Phan)
		w.Bool(m.NangCapDuoc)
	}
	return dongGoiItem(w)
}

// dongGoiItem kiem tran than goi roi dong lai. Client goi d.mark(100000) truoc
// khi doc; than dai hon la client doc lech ma KHONG bao loi.
func dongGoiItem(w *protocol.Writer) (*protocol.Message, error) {
	if len(w.Bytes()) >= tranThanGoiItem {
		return nil, fmt.Errorf(
			"than goi -28/8 = %d byte, phai < %d (client mark(100000))",
			len(w.Bytes()), tranThanGoiItem)
	}
	return protocol.NewMessage(lenhKhongTrenBanDo, w), nil
}

// goiCatalogItemDayDu tra DU BA message -28/8 theo DUNG THU TU PHAI GUI:
//
//	subtype 0 (bang option) -> subtype 1 (mau vat pham) -> subtype 2 (dong dai)
//
// ⚠️ THU TU LA BAT BUOC, KHONG PHAI SO THICH: Service.gI().clientOk() chi duoc
// goi trong nhanh subtype 2 cua client (Controller.cs:6517-6534). Dat subtype 2
// khong o cuoi = client kẹt man hinh "Xin chờ" mai mai.
//
// Subtype 2 gui khoang RONG [0, 0) vi subtype 1 da nap du 0..n-1 roi; no o day
// thuan tuy de KICH clientOk. Ban Java thay vao do bam cung so 750
// (ItemData.java:21-22) — con so cua game Dragon Ball cu, khong chep lai.
func goiCatalogItemDayDu() ([]*protocol.Message, error) {
	mOption, err := goiCatalogItemOption()
	if err != nil {
		return nil, err
	}
	mMau, err := goiCatalogItemMau()
	if err != nil {
		return nil, err
	}
	return []*protocol.Message{mOption, mMau, goiCatalogItem()}, nil
}
