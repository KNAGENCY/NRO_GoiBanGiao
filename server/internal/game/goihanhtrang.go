// goihanhtrang.go — KH7 Task 3: bo dung BA goi hanh trang gui XUONG client:
//
//	-36  tui do
//	-37  do dang mac tren nguoi
//	-35  ruong (kho)
//
// Nguon giai ma hai phia (PHA A):
// docs/superpowers/research/2026-09-03-kh7-giai-ma-36-37-35.md
// Dac ta: docs/superpowers/specs/2026-09-03-kh7-vat-pham-design.md §4.2-§4.4, §4.6b.
//
// BA SU THAT PHAI NHO — sai la hong NGAM, khong loi, khong log:
//
//  1. O RONG PHAI GHI short -1. Client luon lap DU soO lan va chi bo qua khi doc
//     duoc -1 (Controller.cs:2214-2218 / :2019 / :2277). Ban Java bo qua o rong
//     ma KHONG ghi -1 o goi -36 (InventoryService.java:570-572) — do la loi L1,
//     goi thieu byte, client doc lan sang phan khac roi nem loi va bi
//     `catch (Exception) { }` rong nuot sach. Nguoi choi thay hanh trang trong.
//
//  2. KHOI OPTION CHOT CUNG `byte optionId` + `unsigned short param`, vi
//     ModFunc.isShortOptionTemp == false (ModFunc.cs:300) va ModFunc.isReadInt
//     == true (ModFunc.cs:45). param AM se hien thanh so khong lo tren man hinh
//     nguoi choi ma khong bao gi (§7.1 dac ta).
//
//  3. DUNG XONG TOAN BO DU LIEU ROI MOI GHI BYTE. Loi giua chung (mau khong ton
//     tai, o_so ngoai bien, option khong co mau, param tran) => tra error va
//     KHONG gui gi ca. Gui mot goi thieu byte con te hon khong gui — no de lai
//     tui do hien thi do dang (§1.10 / L13).
package game

import (
	"fmt"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// ---------------------------------------------------------------------------
// 1. Ma lenh — doi chieu truc tiep voi `case` cua client
// ---------------------------------------------------------------------------

const (
	// lenhTui: client doc o Controller.cs:2202 (case -36).
	lenhTui int8 = -36
	// lenhDoDangMac: client doc o Controller.cs:2005 (case -37).
	lenhDoDangMac int8 = -37
	// lenhRuong: client doc o Controller.cs:2265 (case -35).
	lenhRuong int8 = -35
)

// Lenh con (cAction) — byte dau tien cua ca ba goi.
const (
	// hanhTrangNapCaBang: nap lai TOAN BO danh sach o. La lenh con duy nhat
	// task 3 dung; cac lenh con khac (doi so luong mot o, mo giao dien ruong)
	// thuoc task sau.
	hanhTrangNapCaBang byte = 0
)

// ---------------------------------------------------------------------------
// 2. Tran cung
// ---------------------------------------------------------------------------

const (
	// tranSoO: so o toi da cua MOT bang (tui / mac / ruong).
	//
	// VI SAO 127 CHU KHONG PHAI 255: client doc soO bang readUnsignedByte
	// (toi da 255, Controller.cs:2208/:2013/:2271) — tran that cua rieng ba goi
	// nay la 255. NHUNG cung con so do con di xuong lan nua trong ME_INFO
	// (-30/0), noi client doc bang byte CO DAU (spec vao the gioi §11.2, dac ta
	// KH7 §4.9). Hai cho phai cung mot con so, nen tran chung la 127. Muon noi
	// len phai sua ME_INFO TRUOC.
	tranSoO = 127

	// tranOptionMoiMon: client doc soOption bang readUnsignedByte.
	tranOptionMoiMon = 255

	// tranParamOption: client doc param bang readUnsignedShort (ReadParamDQT
	// voi isReadInt == true). Ghi ngoai mien nay la hien sai, khong bao loi.
	tranParamOption = 65535
)

// soOMacChuan: so o do dang mac cua ban 1 — 0 ao, 1 quan, 2 gang, 3 giay,
// 4 nhan (bo cuc goc, Player.java:878-882; client ve theo dung bo cuc nay).
const soOMacChuan = 5

// ---------------------------------------------------------------------------
// 3. Bo dung — ba ham, mot khuon chung
// ---------------------------------------------------------------------------

// GoiTui dung goi -36 (nap ca tui).
//
//	[0][soO: byte]( [short -1]  |  MOT MON )* soO
//
// dsMon la CA danh sach vat pham cua mot nguoi; ham tu loc theo NoiChua = TUI.
// Loc ben trong (thay vi tin nguoi goi loc san) de mot mon nam nham noi chua
// khong bao gio lot xuong client o sai bang.
func GoiTui(dsMon []store.VatPham, soO int) (*protocol.Message, error) {
	o, err := xepVaoO(dsMon, store.NoiChuaTui, soO)
	if err != nil {
		return nil, fmt.Errorf("goi -36 (tui): %w", err)
	}
	w := protocol.NewWriter()
	w.Byte(hanhTrangNapCaBang)
	w.Byte(byte(soO))
	if err := ghiCacO(w, o); err != nil {
		return nil, fmt.Errorf("goi -36 (tui): %w", err)
	}
	return protocol.NewMessage(lenhTui, w), nil
}

// GoiDoDangMac dung goi -37 (do dang mac tren nguoi).
//
//	[0][head: short][soO: byte]( [short -1] | MOT MON )* soO
//
// Khac hai goi kia DUNG MOT truong: `head` chen giua cAction va soO
// (Controller.cs:2011). Quen no la lech 2 byte va ca goi vo im lang.
func GoiDoDangMac(dsMon []store.VatPham, head int16, soO int) (*protocol.Message, error) {
	o, err := xepVaoO(dsMon, store.NoiChuaMac, soO)
	if err != nil {
		return nil, fmt.Errorf("goi -37 (do mac): %w", err)
	}
	w := protocol.NewWriter()
	w.Byte(hanhTrangNapCaBang)
	w.Short(head)
	w.Byte(byte(soO))
	if err := ghiCacO(w, o); err != nil {
		return nil, fmt.Errorf("goi -37 (do mac): %w", err)
	}
	return protocol.NewMessage(lenhDoDangMac, w), nil
}

// GoiRuong dung goi -35 (ruong / kho). Than goi giong het -36.
func GoiRuong(dsMon []store.VatPham, soO int) (*protocol.Message, error) {
	o, err := xepVaoO(dsMon, store.NoiChuaRuong, soO)
	if err != nil {
		return nil, fmt.Errorf("goi -35 (ruong): %w", err)
	}
	w := protocol.NewWriter()
	w.Byte(hanhTrangNapCaBang)
	w.Byte(byte(soO))
	if err := ghiCacO(w, o); err != nil {
		return nil, fmt.Errorf("goi -35 (ruong): %w", err)
	}
	return protocol.NewMessage(lenhRuong, w), nil
}

// ---------------------------------------------------------------------------
// 4. Hai buoc: xep vao o (KIEM HET) roi moi ghi byte
// ---------------------------------------------------------------------------

// xepVaoO loc dsMon theo noiChua roi xep vao mang soO phan tu, nil = o rong.
// TAT CA viec kiem tra nam o day — sau khi ham nay tra ve khong loi thi buoc
// ghi byte khong con cach nao that bai vi du lieu.
func xepVaoO(dsMon []store.VatPham, noiChua string, soO int) ([]*store.VatPham, error) {
	if soO < 0 || soO > tranSoO {
		return nil, fmt.Errorf("so o = %d, phai trong [0,%d]", soO, tranSoO)
	}
	o := make([]*store.VatPham, soO)
	for i := range dsMon {
		m := &dsMon[i]
		if m.NoiChua != noiChua {
			continue
		}
		if m.OSo == nil {
			return nil, fmt.Errorf("mon id=%d o %s khong co o_so", m.ID, noiChua)
		}
		vt := int(*m.OSo)
		if vt < 0 || vt >= soO {
			return nil, fmt.Errorf("mon id=%d co o_so=%d, ngoai bien [0,%d)", m.ID, vt, soO)
		}
		if o[vt] != nil {
			return nil, fmt.Errorf("hai mon cung o %d (id=%d va id=%d)", vt, o[vt].ID, m.ID)
		}
		if err := kiemMon(m); err != nil {
			return nil, err
		}
		o[vt] = m
	}
	return o, nil
}

// kiemMon kiem MOT mon co gui xuong client duoc khong.
func kiemMon(m *store.VatPham) error {
	if m.MauID < 0 || m.MauID > 32767 {
		return fmt.Errorf("mon id=%d co mau_id=%d, ngoai mien short", m.ID, m.MauID)
	}
	if TimMauVatPham(int16(m.MauID)) == nil {
		return fmt.Errorf("mon id=%d tro toi mau %d khong ton tai trong danh muc", m.ID, m.MauID)
	}
	if m.SoLuong < 1 {
		return fmt.Errorf("mon id=%d co so_luong=%d, phai >= 1", m.ID, m.SoLuong)
	}
	opt := m.ChiSoPhu.Opt
	if len(opt) > tranOptionMoiMon {
		return fmt.Errorf("mon id=%d co %d option, client chi doc duoc %d (mot byte khong dau)",
			m.ID, len(opt), tranOptionMoiMon)
	}
	soMauOption := len(DanhSachMauOption())
	for _, t := range opt {
		// L12: client cap phat mang mau option bang mot byte va tra
		// iOptionTemplates[id] KHONG kiem bien (ItemOption.cs:40) — gui id
		// khong co mau la client nem IndexOutOfRange roi bi nuot.
		if t.Id < 0 || t.Id >= soMauOption {
			return fmt.Errorf("mon id=%d dung option id=%d, danh muc chi co %d mau",
				m.ID, t.Id, soMauOption)
		}
		if t.Param < 0 || t.Param > tranParamOption {
			return fmt.Errorf("mon id=%d co option id=%d param=%d, phai trong [0,%d] (client doc unsigned short)",
				m.ID, t.Id, t.Param, tranParamOption)
		}
	}
	return nil
}

// ghiCacO ghi than goi: mot o rong la DUNG MOT short -1, mot o co mon la khuon
// bay truong. Chi duoc goi sau xepVaoO — loi duy nhat con lai la chuoi qua dai.
func ghiCacO(w *protocol.Writer, o []*store.VatPham) error {
	for vt, m := range o {
		if m == nil {
			// DAY LA DONG SUA LOI L1. Bo dong nay = goi thieu byte moi khi
			// co o trong, va client im lang khong hien hanh trang.
			w.Short(-1)
			continue
		}
		w.Short(int16(m.MauID))
		w.Int(m.SoLuong)
		if err := w.UTF(moTaNgan(m)); err != nil {
			return fmt.Errorf("o %d: khong ghi duoc info: %w", vt, err)
		}
		if err := w.UTF(moTaDai(m)); err != nil {
			return fmt.Errorf("o %d: khong ghi duoc content: %w", vt, err)
		}
		w.Byte(byte(len(m.ChiSoPhu.Opt)))
		for _, t := range m.ChiSoPhu.Opt {
			w.Byte(byte(t.Id))
			// param KHONG DAU: da kep 0..65535 o kiemMon.
			w.Short(int16(uint16(t.Param)))
		}
	}
	return nil
}

// moTaNgan / moTaDai tra hai chuoi hien thi cua mot mon (client goi la `info`
// va `content`). BAN 1 TRA CHUOI RONG co chu dich: chua lan duoc trong client
// cho nao VE hai chuoi nay, nen chua biet khuon dung. Chuoi rong = 2 byte 00 00,
// dung so byte, chi lam o mo ta trong. Xem §7 tai lieu giai ma PHA A.
func moTaNgan(*store.VatPham) string { return "" }
func moTaDai(*store.VatPham) string  { return "" }
