package game

// muaban.go — KH7 Task 9: MUA HANG NPC (goi LEN `6`).
//
// Giai ma giao thuc (PHA A):
// docs/superpowers/research/2026-09-03-kh7-giai-ma-mua-6.md
//
// BON DIEU CHOT — moi cai la mot dot bien bat buoc phai bi giet:
//
//  1. MOT giao dich CSDL cho ca "tru tien" lan "cap mon"
//     (store.MuaMotGiaoDich). Tach lam hai la D2: dut giua chung thi mat tien
//     ma khong co do.
//  2. Khoa idempotency `<nhanVatId>:MUA:<soThuTu>` (D3). `soThuTu` la bo dem
//     MAY CHU giu trong phien — client khong bao gio duoc quyet dinh no, neu
//     khong ke gian chi viec gui lai so cu de... khong, gui so MOI de mua that.
//     Bo dem tang MOI LENH KINH TE, ke ca lenh bi tu choi (xem soThuTuKinhTeKe).
//  3. KIEM HET roi moi doi bat cu thu gi (§2.3 buoc 4 / bay D16 — ban Java tru
//     phi TRUOC khi kiem du lieu). Moi cau kiem duoi day nam TRUOC loi goi CSDL
//     duy nhat cua ham.
//  4. Goi XUONG `6` (tien) LUON duoc gui, ke ca khi mua that bai — vi client
//     goi GameCanvas.endDlg() trong chinh nhanh do (PHA A §2). Khong gui la
//     nguoi choi ket man hinh "vui long doi".
//
// BA DOT BIEN DA KIEM TAY (03/09) — moi cai lam DUNG mot test do:
//
//	1. store/muaban_kho.go: doi congTruTien(ctx, tx, ...) thanh s.CongTruTien(
//	   ctx, ...) — tuc tru tien thanh mot giao dich RIENG, commit ngay.
//	   -> TestMuaCapMonHongKhongMatTien do ("vang = 900, phai giu nguyen 1000").
//	2. store/muaban_kho.go: bo ON CONFLICT (khoa_idem) DO NOTHING (moi goi deu
//	   duoc coi la lan dau).
//	   -> TestMuaGuiLapChiTruTienMotLan do ("lan 2 phai duoc danh dau la goi lap").
//	3. muaban.go: tru tien TRUOC cau kiem "con o trong khong".
//	   -> TestMuaTuiDayTienKhongDoi do ("vang = 900, phai giu nguyen 1000").
//
// MOI NOI VOI TASK 5 (thegioi.go): so o tui that su cua nhan vat hien chua co
// trong BoXuLy. Ham thuan `MuaTuCuaHang` nhan `soOTui` lam THAM SO; handler tam
// dung `soOTuiChuan`. Khi task 5 dat so o tui vao BoXuLy/NhanVat, chi can doi
// MOT dong trong `xuLyMua` — khong dung toi ham thuan.

import (
	"context"
	"errors"
	"fmt"

	"nroserver/internal/protocol"
	"nroserver/internal/store"

	mynet "nroserver/internal/net"
)

// ---------------------------------------------------------------------------
// 1. Hang so giao thuc
// ---------------------------------------------------------------------------

const (
	// lenhMua: goi LEN `6` (client xin mua). CUNG con so voi goi XUONG tien —
	// giao thuc goc dung chung ma lenh cho hai chieu, khong phai nham lan.
	lenhMua int8 = 6
	// lenhMuaBanTien: goi XUONG `6` — [vang long][ngoc int][ngocKhoa int=0].
	lenhMuaBanTien int8 = 6
	// lenhMuaBanBaoNhanh: goi XUONG `-25` — dong chu chay, khong chan man hinh.
	lenhMuaBanBaoNhanh int8 = -25
)

// muaBanTypeHopLe: bang `typeBuy` cua goi LEN `6` ma ban 1 chap nhan.
// 0 = mua le (Panel.cs idAction 3000), 3 = mua nhieu cai (AutoBuyItem).
// Cac gia tri khac (1, 2) thuoc man hinh chua cai — tu choi thay vi doan.
func muaBanTypeHopLe(t int8) bool { return t == 0 || t == 3 }

// soOTuiChuan: so o tui cua ban 1. Tam o day cho toi khi task 5 dua so o that
// vao nhan vat — xem chu thich "MOI NOI VOI TASK 5" dau file.
const soOTuiChuan = 30

// ---------------------------------------------------------------------------
// 2. Loi — moi cai mot cau tieng Viet cho nguoi choi
// ---------------------------------------------------------------------------

var (
	// ErrKhongMoCuaHang: chua mo cua hang nao ma da gui goi mua.
	ErrKhongMoCuaHang = errors.New("chua mo cua hang nao")
	// ErrKhongBanMonNay: mau khong nam trong cua hang dang mo.
	ErrKhongBanMonNay = errors.New("cua hang khong ban mon nay")
	// ErrSoLuongKhongHopLe: <= 0, hoac vuot tran chong cua mau.
	ErrSoLuongKhongHopLe = errors.New("so luong khong hop le")
	// ErrTuiDay: khong con o trong nao trong tui.
	ErrTuiDay = errors.New("hanh trang da day")
	// ErrTypeMuaLa: `typeBuy` ngoai bang cho phep.
	ErrTypeMuaLa = errors.New("kieu mua khong ho tro")
)

// ---------------------------------------------------------------------------
// 3. Kho — interface HEP, dung khuon voi Kho/GhiKhiRoi trong login.go
// ---------------------------------------------------------------------------

// KhoMuaBan la phan toi thieu cua *store.Store ma duong mua can. Hep de test
// duoc bang mot kho gia, va de doc mot cai la biet ham nay cham vao nhung gi.
type KhoMuaBan interface {
	LayTui(ctx context.Context, nvID int64, noiChua string) ([]store.VatPham, error)
	MuaMotGiaoDich(ctx context.Context, yc store.YeuCauMua) (store.KetQuaMua, error)
}

// ---------------------------------------------------------------------------
// 4. Ham mua THUAN — khong biet gi ve mang
// ---------------------------------------------------------------------------

// KetQuaMuaHang la thu tang tren can biet de dung goi tra ve.
type KetQuaMuaHang struct {
	// VatPhamID/OSo/SoLuong: mon vua vao tui.
	VatPhamID int64
	OSo       int16
	MauID     int16
	SoLuong   int32
	// TenMon de ghep cau "Mua thanh cong X".
	TenMon string
	// DVang/DNgoc: so tien DA tru (am hoac 0).
	DVang int64
	DNgoc int64
	// LaLapLai: goi nay la ban sao cua mot lan mua truoc (D3). Tang tren van
	// gui lai tui/tien y het lan dau — chi KHONG tru them va KHONG cap them.
	LaLapLai bool
}

// MuaTuCuaHang thi hanh mot lan mua tu cua hang NPC.
//
// Thu tu CO Y (§2.3): moi cau kiem o day nam TRUOC loi goi CSDL duy nhat
// (kho.MuaMotGiaoDich) o cuoi ham. Doi mot cau kiem xuong duoi loi goi do la
// mo lai dung lo hong K3/D16 cua ban Java.
//
//	ctx      — ngat duoc
//	kho      — CSDL
//	nvID     — nhan vat mua
//	ch       — cua hang DANG MO (may chu giu, khong phai client gui)
//	typeBuy  — o thu nhat cua goi `6`
//	mauID    — o thu hai cua goi `6`
//	soLuong  — o thu ba (da chuan hoa: vang -> 1)
//	soThuTu  — bo dem kinh te cua phien (D3)
//	soOTui   — so o tui cua nhan vat
func MuaTuCuaHang(ctx context.Context, kho KhoMuaBan, nvID int64, ch *CuaHang,
	typeBuy int8, mauID int16, soLuong int32, soThuTu uint64, soOTui int) (KetQuaMuaHang, error) {

	// --- Kiem 1: cua hang phai dang mo (client khong gui id cua hang). ---
	if ch == nil {
		return KetQuaMuaHang{}, ErrKhongMoCuaHang
	}
	// --- Kiem 2: kieu thao tac. ---
	if !muaBanTypeHopLe(typeBuy) {
		return KetQuaMuaHang{}, fmt.Errorf("%w: typeBuy=%d", ErrTypeMuaLa, typeBuy)
	}
	// --- Kiem 3: so luong duong (tran kiem o buoc 6, sau khi biet mau). ---
	if soLuong < 1 {
		return KetQuaMuaHang{}, fmt.Errorf("%w: %d", ErrSoLuongKhongHopLe, soLuong)
	}
	// --- Kiem 4: mon co trong cua hang khong (duyet MOI tab — PHA A §1.1). ---
	mon := muaBanTimMon(ch, mauID)
	if mon == nil {
		return KetQuaMuaHang{}, fmt.Errorf("%w: mau %d", ErrKhongBanMonNay, mauID)
	}
	// --- Kiem 5: mau ton tai. `kiemCuaHang` da bao dam, nhung cua hang co the
	// duoc sua luc chay -> kiem lai chu khong tin. ---
	mau := TimMauVatPham(mon.MauID)
	if mau == nil {
		return KetQuaMuaHang{}, fmt.Errorf("mau vat pham %d khong ton tai", mon.MauID)
	}
	// --- Kiem 6: cong don hay don chiec (§1.3). Mon co option LUON 1 cai. ---
	tranChong := int32(mau.TranChong)
	if tranChong < 1 {
		tranChong = 1
	}
	if len(mon.Options) > 0 {
		tranChong = 1
	}
	if soLuong > tranChong {
		return KetQuaMuaHang{}, fmt.Errorf("%w: mua %d, mot o chua toi da %d",
			ErrSoLuongKhongHopLe, soLuong, tranChong)
	}
	// --- Kiem 7: gia. Tinh tren int64 de N * gia khong bao gio tran. ---
	tong := int64(mon.Gia) * int64(soLuong)
	if tong < 0 {
		return KetQuaMuaHang{}, fmt.Errorf("tong gia am (%d)", tong)
	}
	var dVang, dNgoc int64
	switch mon.Tien {
	case TienVang:
		dVang = -tong
	case TienNgoc:
		dNgoc = -tong
	default:
		return KetQuaMuaHang{}, fmt.Errorf("loai tien %d khong hop le", mon.Tien)
	}

	// --- Kiem 8: con o trong khong. ---
	//
	// Doc nay nam NGOAI giao dich (LayTui la mot cau doc rieng). Khong sao: neu
	// giua doc va ghi co mot duong khac chiem mat o do, buoc INSERT trong giao
	// dich dung UNIQUE `vat_pham_mot_o_mot_mon` -> ErrOCoMon -> CA giao dich
	// ROLLBACK. Nguoi choi thay loi, KHONG mat tien. Doc som chi de tu choi
	// som (§2.3 buoc 1), khong phai chot chan.
	tui, err := kho.LayTui(ctx, nvID, store.NoiChuaTui)
	if err != nil {
		return KetQuaMuaHang{}, fmt.Errorf("doc tui that bai: %w", err)
	}
	oTrong, coOTrong := muaBanTimOTrong(tui, soOTui)
	if !coOTrong {
		return KetQuaMuaHang{}, ErrTuiDay
	}

	// --- Chi TOI DAY moi cham CSDL de GHI. Mot giao dich duy nhat. ---
	kq, err := kho.MuaMotGiaoDich(ctx, store.YeuCauMua{
		NhanVatID: nvID,
		KhoaIdem:  store.KhoaIdemMua(nvID, soThuTu),
		DVang:     dVang,
		DNgoc:     dNgoc,
		NoiChua:   store.NoiChuaTui,
		OSo:       oTrong,
		MauID:     int32(mon.MauID),
		SoLuong:   soLuong,
		ChiSoPhu:  muaBanChiSoPhu(mon),
		LyDo:      store.LyDoMuaNPC,
	})
	if err != nil {
		return KetQuaMuaHang{}, err
	}
	return KetQuaMuaHang{
		VatPhamID: kq.VatPhamID,
		OSo:       kq.OSo,
		MauID:     int16(kq.MauID),
		SoLuong:   kq.SoLuong,
		TenMon:    mau.Ten,
		DVang:     kq.DVang,
		DNgoc:     kq.DNgoc,
		LaLapLai:  kq.LaLapLai,
	}, nil
}

// muaBanTimMon duyet MOI tab tim mon theo MauID — dung y ban Java
// (Shop.getItemShop, PHA A §1.1). Tra nil neu cua hang khong ban mau do.
func muaBanTimMon(ch *CuaHang, mauID int16) *MonHang {
	for iTab := range ch.Tabs {
		for iMon := range ch.Tabs[iTab].Mon {
			if ch.Tabs[iTab].Mon[iMon].MauID == mauID {
				return &ch.Tabs[iTab].Mon[iMon]
			}
		}
	}
	return nil
}

// muaBanTimOTrong tra o TRONG NHO NHAT trong [0, soOTui). Chon o nho nhat (chu
// khong phai o bat ky) de ket qua tat dinh — test khong phai doan.
func muaBanTimOTrong(tui []store.VatPham, soOTui int) (int16, bool) {
	daDung := make(map[int16]bool, len(tui))
	for _, vp := range tui {
		if vp.NoiChua == store.NoiChuaTui && vp.OSo != nil {
			daDung[*vp.OSo] = true
		}
	}
	for o := 0; o < soOTui; o++ {
		if !daDung[int16(o)] {
			return int16(o), true
		}
	}
	return 0, false
}

// muaBanChiSoPhu doi option gan san cua mon hang thanh khuon JSONB cua kho.
func muaBanChiSoPhu(mon *MonHang) store.ChiSoPhu {
	opt := make([]store.TuyChon, 0, len(mon.Options))
	for _, o := range mon.Options {
		opt = append(opt, store.TuyChon{Id: int(o.OptionID), Param: int(o.Param)})
	}
	return store.ChiSoPhu{Opt: opt}
}

// ---------------------------------------------------------------------------
// 5. Goi XUONG `6` — tien
// ---------------------------------------------------------------------------

// goiMuaBanTien dung goi XUONG `6`: [vang long][ngoc int][ngocKhoa int].
//
// O thu ba LUON 0: ban 1 chi co vang va ngoc (bay L10). O thu nhat LUON long —
// ban Java ghi int cho client cu, client LUON doc long (bay L9).
func goiMuaBanTien(vang, ngoc int64) *protocol.Message {
	w := protocol.NewWriter()
	w.Long(vang)
	w.Int(int32(ngoc))
	w.Int(0)
	return protocol.NewMessage(lenhMuaBanTien, w)
}

// goiMuaBanBao dung goi XUONG `-25`: mot dong chu chay. Tra nil neu noi dung
// khong ghi duoc (writeUTF do dai 2 byte) — nguoi goi bo qua, khong gui goi loi.
func goiMuaBanBao(noiDung string) *protocol.Message {
	w := protocol.NewWriter()
	if err := w.UTF(noiDung); err != nil {
		return nil
	}
	return protocol.NewMessage(lenhMuaBanBaoNhanh, w)
}

// muaBanCauLoi doi loi ky thuat thanh cau nguoi choi doc duoc. Khong bao gio lo
// chi tiet ky thuat ra man hinh.
func muaBanCauLoi(err error) string {
	switch {
	case errors.Is(err, store.ErrKhongDuTien):
		return "Bạn không đủ tiền"
	case errors.Is(err, ErrTuiDay), errors.Is(err, store.ErrOCoMon):
		return "Hành trang đã đầy"
	case errors.Is(err, ErrKhongBanMonNay), errors.Is(err, ErrKhongMoCuaHang):
		return "Không thể thực hiện"
	case errors.Is(err, ErrSoLuongKhongHopLe):
		return "Số lượng không hợp lệ"
	default:
		return "Không thể thực hiện"
	}
}

// ---------------------------------------------------------------------------
// 6. Handler goi LEN `6`
// ---------------------------------------------------------------------------

// soThuTuKinhTeKe tra so thu tu ke tiep cho mot hanh dong kinh te cua phien nay
// (D3). Tang MOI LAN goi, ke ca khi lenh sau do bi tu choi: neu chi tang khi
// thanh cong thi mot lan mua that bai roi mot lan mua thanh cong se dung CUNG
// mot so — lan thu hai bi cham vao khoa cu va bi coi la "goi lap".
func (b *BoXuLy) soThuTuKinhTeKe() uint64 {
	b.muaBanSoThuTu++
	return b.muaBanSoThuTu
}

// xuLyMua xu ly goi LEN `6`.
//
// Doc goi -> goi ham thuan -> gui phan hoi. Ba dieu ve phan hoi (PHA A §2):
//   - Thanh cong: -36 (tui moi) -> -25 ("Mua thanh cong X") -> 6 (tien moi).
//   - That bai:   -25 (ly do) -> 6 (tien KHONG doi).
//   - Goi `6` LUON la goi cuoi — no dong hop thoai "vui long doi" phia client.
//
// Goi rac (thieu byte) thi KHONG tra loi gi va KHONG ngat ket noi: dung khuon
// cua cac handler khac trong login.go.
func (b *BoXuLy) xuLyMua(s *mynet.Session, m *protocol.Message) error {
	if b.nhanVat == nil {
		return nil
	}
	r := protocol.NewReader(m.Data)
	typeBuy, err := r.Byte()
	if err != nil {
		return nil
	}
	mauID, err := r.Short()
	if err != nil {
		return nil
	}
	// O thu ba CO THE VANG (client chi ghi khi quantity > 1) — thieu la 1.
	soLuong := int32(1)
	if n, err := r.Short(); err == nil && n > 1 {
		soLuong = int32(n)
	}

	ctx := context.Background()
	kho, ok := b.kho.(KhoMuaBan)
	if !ok {
		return b.muaBanGuiLoi(s, "Không thể thực hiện")
	}

	kq, errMua := MuaTuCuaHang(ctx, kho, b.nhanVat.ID, b.muaBanCuaHangDangMo,
		int8(typeBuy), mauID, soLuong, b.soThuTuKinhTeKe(), soOTuiChuan)
	if errMua != nil {
		return b.muaBanGuiLoi(s, muaBanCauLoi(errMua))
	}

	// Chi SAU KHI commit tra ve nil moi cham ban trong bo nho va gui goi
	// (§2.3 buoc 7). Goi lap cung cap nhat lai bo nho — vi so tien trong
	// KetQuaMuaHang la so DA tru cua lan dau, va lan dau da tru roi nen o day
	// KHONG duoc tru them.
	if !kq.LaLapLai {
		b.nhanVat.Vang += kq.DVang
		b.nhanVat.Ngoc += kq.DNgoc
	}

	// Doc lai tui tu CSDL (nguon su that) roi dung goi -36.
	tui, err := kho.LayTui(ctx, b.nhanVat.ID, store.NoiChuaTui)
	if err == nil {
		if goi, err := GoiTui(tui, soOTuiChuan); err == nil {
			if err := s.XepGui(goi); err != nil {
				return err
			}
		}
	}
	if goi := goiMuaBanBao("Mua thành công " + kq.TenMon); goi != nil {
		if err := s.XepGui(goi); err != nil {
			return err
		}
	}
	return s.XepGui(goiMuaBanTien(b.nhanVat.Vang, b.nhanVat.Ngoc))
}

// muaBanGuiLoi gui cau tu choi roi gui goi tien KHONG DOI — goi tien la thu
// dong hop thoai cho phia client, thieu no la nguoi choi ket man hinh.
func (b *BoXuLy) muaBanGuiLoi(s *mynet.Session, cau string) error {
	if goi := goiMuaBanBao(cau); goi != nil {
		if err := s.XepGui(goi); err != nil {
			return err
		}
	}
	var vang, ngoc int64
	if b.nhanVat != nil {
		vang, ngoc = b.nhanVat.Vang, b.nhanVat.Ngoc
	}
	return s.XepGui(goiMuaBanTien(vang, ngoc))
}

// DatCuaHangDangMo ghi nho cua hang phien nay vua mo (goi -44). Client KHONG
// gui id cua hang trong goi `6`, nen day la nguon su that duy nhat.
//
// Task mo NPC se goi ham nay. De o day (khong phai trong cuahang.go) vi no
// cham truong cua BoXuLy.
func (b *BoXuLy) DatCuaHangDangMo(ch *CuaHang) { b.muaBanCuaHangDangMo = ch }
