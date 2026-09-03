// mauvatpham.go — CUA cua goi `game` di vao MAU VAT PHAM + MAU OPTION cua kho
// noi dung, va bo dung goi -28/8 subtype 0 va 1.
//
// # CAI GI DOI O KH8 TASK 4
//
// Truoc task nay du lieu la THUAN Go: hai bien toan cuc danhSachMauOption /
// danhSachMauVatPham, song song voi kho-noi-dung/vat-pham.yaml. Hai ban so lieu
// cho cung mot su that.
//
// Tu task nay chi con MOT ban: vat-pham.yaml. Hai bang Go DA BI XOA. Ba cua tra
// cuu (DanhSachMauOption, DanhSachMauVatPham, TimMauVatPham) doc tu ANH CHUP
// kho noi dung — mac dinh la ban NHUNG (cung sync.Once cua task 3, khochieu.go),
// hoac mot anh chup tiem vao khi nguoi goi dua san.
//
// VI SAO KHONG PHAI catalog.json: catalog.json bi ghim SHA trong asset-lock.json
// (catalog.go, lockHashKhop). Do la ly do KH7 khong dat du lieu vao do; KH8 dat
// vao kho noi dung YAML, mot cho khac han, khong dung toi asset lock.
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

	"nroserver/internal/noidung"
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

// ⚠️ KH8 task 0 da CHUYEN cac tran nay ve goi noidung (noidung.Tran*) — mot noi
// duy nhat cho ca kho noi dung, vi bo kiem YAML cung phai dung dung nhung so
// nay. O day chi con BI DANH cho ngan gon; KHONG duoc go so vao thay.
// (03/09 hai task song song moi ben tu khai mot ban tranParamOption, gop vao
// main la dut build — bai hoc §8. Mot noi khai, moi noi tro toi.)
const (
	// tranMauOption: client doc so mau bang MOT BYTE KHONG DAU.
	tranMauOption = noidung.TranMauOption

	// tranParamOption: gia tri LON NHAT cua param mot option. Client doc param
	// bang ReadParamDQT = readUnsignedShort (ModFunc.isReadInt == true) -> 0..65535;
	// ghi so am hoac lon hon la client hien sai ma khong bao loi (bay L3 spec KH7).
	tranParamOption = noidung.TranParamOption
	// tranMauVatPham: client doc so mau bang SHORT CO DAU (Controller.cs:6495).
	tranMauVatPham = noidung.TranMauVatPham
	// tranThanGoiItem: client goi d.mark(100000) truoc khi doc (Controller.cs:6468).
	tranThanGoiItem = noidung.TranThanGoiItem
)

// ---------------------------------------------------------------------------
// 3. Ba cua tra cuu — DOC TU ANH CHUP KHO NOI DUNG
// ---------------------------------------------------------------------------
//
// ⚠️ KH8 TASK 4 DA XOA HAI BANG GO (`danhSachMauOption`, `danhSachMauVatPham`).
// Nguon duy nhat bay gio la kho-noi-dung/vat-pham.yaml. Nghia la dot bien "sua
// bang Go cho lech YAML" KHONG THE VIET DUOC NUA — khong con bang nao de sua.
//
// Ba cua duoi day co doi: mot ban KHONG THAM SO doc anh chup MAC DINH (ban
// nhung, dung sync.Once cua task 3 — mot bo nap duy nhat cho ca goi), va mot
// ban `…Tu(ac)` nhan MOT anh chup da cam san. Trong actor phai dung ban `…Tu`
// voi `k.anhChupNoiDung()` de ca nhip chay tron mot the he (hop dong o
// internal/noidung/kho.go). Hom nay KHONG co cho goi nao nam trong actor — moi
// loi goi deu tren goroutine phien hoac goroutine tho — nen ban khong tham so
// van dung; xem bang liet ke o dau §5.

// khoVatPhamMacDinh la CUA lay anh chup mac dinh. La mot BIEN chu khong phai ham
// de test cam duoc mot kho khac (vi du them mot mau Loai == 5 de thu nhanh
// coPart) MA KHONG dung lai mot bang Go thu hai — test van phai dung
// noidung.NapNhung + noidung.MoKho, tuc van la du lieu YAML.
var khoVatPhamMacDinh = func() khoVatPhamAnhChup { return khoChieuMacDinh() }

// DanhSachMauOption tra bang mau option. Chi so = id option gui trong goi vat pham.
func DanhSachMauOption() []MauOption { return mauOptionTu(khoVatPhamMacDinh()) }

// mauOptionTu doi bang mau option cua MOT anh chup sang kieu cua goi `game`.
func mauOptionTu(ac khoVatPhamAnhChup) []MauOption {
	ds := ac.MauOption()
	ra := make([]MauOption, 0, len(ds))
	for _, o := range ds {
		ra = append(ra, mauOptionTuNoiDung(o))
	}
	return ra
}

// DanhSachMauVatPham tra ca danh muc. Chi so = templateId gui xuong client.
func DanhSachMauVatPham() []MauVatPham { return mauVatPhamTu(khoVatPhamMacDinh()) }

// mauVatPhamTu doi danh muc mau vat pham cua MOT anh chup sang kieu goi `game`.
func mauVatPhamTu(ac khoVatPhamAnhChup) []MauVatPham {
	ds := ac.MauVatPham()
	ra := make([]MauVatPham, 0, len(ds))
	for _, m := range ds {
		ra = append(ra, mauVatPhamTuNoiDung(m))
	}
	return ra
}

// TimMauVatPham tra mau theo id, nil neu ngoai mien. MOI cho tra cuu mau phai
// di qua ham nay — id den tu client hoac tu CSDL deu co the ngoai mien (bay §7.11).
func TimMauVatPham(id int16) *MauVatPham { return timMauVatPhamTu(khoVatPhamMacDinh(), id) }

// timMauVatPhamTu tra ve CON TRO toi mot BAN SAO MOI (khong phai vao anh chup):
// nguoi goi lo ghi vao ket qua cung khong lam ban so lieu dang cong bo. Truoc
// task 4 ham nay tra con tro vao bang Go toan cuc — mot cu ghi la doi mon do cua
// moi nguoi choi.
func timMauVatPhamTu(ac khoVatPhamAnhChup, id int16) *MauVatPham {
	if id < 0 {
		return nil
	}
	m, co := ac.TimMauVatPham(int(id))
	if !co {
		return nil
	}
	ra := mauVatPhamTuNoiDung(m)
	return &ra
}

// ---------------------------------------------------------------------------
// 5. Bo dung goi
// ---------------------------------------------------------------------------

// goiCatalogItemOption dung -28/8 subtype 0: bang mau option.
//
//	[8][vsItem][0][soMau: 1 byte KHONG DAU]( [ten UTF][loai byte] )*
//
// (KH8 task 0) Tach lam hai: goiCatalogItemOption() giu nguyen chu ky cu va van
// doc du lieu Go; goiCatalogItemOptionTu(ds) nhan danh muc tu ngoai vao. Tach de
// CONG TUONG DUONG BYTE cho DUNG bo dung goi nay chay tren du lieu nap tu YAML
// — neu test tu viet lai mot ban ghi byte thu hai thi no chi chung minh ban sao
// dung, khong chung minh du lieu dung.
func goiCatalogItemOption() (*protocol.Message, error) {
	return goiCatalogItemOptionTu(DanhSachMauOption())
}

func goiCatalogItemOptionTu(ds []MauOption) (*protocol.Message, error) {
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
	return goiCatalogItemMauTu(DanhSachMauVatPham())
}

func goiCatalogItemMauTu(ds []MauVatPham) (*protocol.Message, error) {
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
