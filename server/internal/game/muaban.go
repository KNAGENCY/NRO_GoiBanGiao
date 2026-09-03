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
//
// # NAP NONG KHONG DOI GIA GIUA MOT PHIEN MUA — CO Y (KH8 task 5)
//
// `ch` la mot BAN SAO lay luc mo cua hang (TimCuaHang -> timCuaHangTu tra con
// tro toi ban sao moi). Phien giu no cho toi khi mo cua hang khac. Nen neu quan
// tri nap nong cua-hang.yaml giua luc mot nguoi dang xem bang gia, NGUOI DO VAN
// MUA/BAN THEO GIA CU cho den khi dong ra mo lai.
//
// Do la lua chon, khong phai thieu sot: client da ve xong bang gia tu goi -44:
// doi gia phia may chu giua chung nghia la nguoi choi bam "mua" o gia 100 va bi
// tru 200 — khong mot dong chu nao giai thich. "Khong nap nong duoc giua phien"
// la cai gia phai tra de con so nguoi choi NHIN THAY va con so may chu TINH luon
// la mot.
//
// BA DIEM GOI cua danh muc cua hang, tat ca NGOAI nhip actor (nen deu dung ban
// cua khong tham so, doc anh chup mac dinh):
//   - TimCuaHang(id) khi mo NPC        -> goroutine PHIEN
//   - xuLyMua / xuLyBan                -> goroutine PHIEN, doc `ch` da cam san
//   - kiemMoiCuaHang / DanhSachCuaHang -> test va cong cu, ngoai gio chay
func (b *BoXuLy) DatCuaHangDangMo(ch *CuaHang) { b.muaBanCuaHangDangMo = ch }

// ===========================================================================
// KH7 Task 10 — BAN cho NPC (goi LEN `7`). Anh guong cua duong mua o tren.
//
// Giai ma giao thuc (PHA A):
// docs/superpowers/research/2026-09-03-kh7-giai-ma-ban-7.md
//
// BA DIEU KHAC MUA — moi cai la mot cho de sai:
//
//  1. O thu ba cua goi `7` la CHI SO O, khong phai MauID. Nen duong ban bat
//     buoc di qua bang phien dich o -> vat_pham.id (task 4, BangHanhTrang) —
//     `vat_pham.id` khong bao gio xuong client (dac ta 2.3 / 4.1).
//  2. Goi `7` co HAI VONG: action 0 = xin xac nhan (may chu tra hop thoai,
//     KHONG doi gi trong CSDL), action khac 0 = ban that.
//  3. Gia ban lay tu mon tuong ung trong CUA HANG DANG MO: max(1, Gia/4).
//     MauVatPham khong co truong gia (xem PHA A muc 3) — quyet dinh nay co
//     GHI RO trong tai lieu, dung sua o day ma khong sua tai lieu.
//
// Ky luat GIU NGUYEN tu duong mua, khong duoc noi long:
//   - MOT giao dich CSDL cho ca "bot do" lan "cong tien" (store.BanMotGiaoDich).
//   - Khoa idempotency `<nvID>:BAN:<soThuTu>`, cung bo dem soThuTuKinhTeKe.
//   - KIEM HET roi moi cham CSDL de ghi.
//   - Goi XUONG `6` (tien) LUON duoc gui, ke ca khi ban that bai.

// lenhBan: goi LEN `7` (client xin ban) VA goi XUONG `7` (hop thoai xac nhan) —
// cung mot ma lenh cho hai chieu, giong het cap `6`.
const lenhBan int8 = 7

// banWhereMac: gia tri o `where` nghia la "do dang MAC" (Java: itemsBody).
// Moi gia tri khac deu la TUI. Ban 1 chi ban do trong TUI (PHA A muc 4 dong 1).
const banWhereMac int8 = 0

var (
	// ErrBanDoDangMac: xin ban do dang mac — ban 1 tu choi.
	ErrBanDoDangMac = errors.New("khong ban duoc do dang mac")
	// ErrBanOTrong: o do khong co mon nao.
	ErrBanOTrong = errors.New("o do khong co mon nao")
	// ErrBanKhongMuaMon: cua hang dang mo khong ban mau nay -> cung khong mua lai.
	ErrBanKhongMuaMon = errors.New("cua hang khong mua mon nay")
	// ErrOSoKhongHopLe: chi so o ngoai bien.
	ErrOSoKhongHopLe = errors.New("o so khong hop le")
)

// KhoBan la phan toi thieu cua *store.Store ma duong ban can. Hep, cung khuon
// voi KhoMuaBan.
type KhoBan interface {
	LayTui(ctx context.Context, nvID int64, noiChua string) ([]store.VatPham, error)
	BanMotGiaoDich(ctx context.Context, yc store.YeuCauBan) (store.KetQuaBan, error)
}

// banDuTinh la ket qua cua buoc TINH — moi thu da kiem xong, chua cham CSDL.
// Dung cho ca hai vong: vong xac nhan chi can no de dung cau hoi, vong ban that
// dua thang no vao giao dich.
type banDuTinh struct {
	VatPhamID int64
	MauID     int16
	TenMon    string
	SoLuong   int32
	DonGia    int32
	Tien      LoaiTien
	DVang     int64
	DNgoc     int64
}

// KetQuaBanHang la thu tang tren can biet de dung goi tra ve.
type KetQuaBanHang struct {
	VatPhamID  int64
	MauID      int16
	TenMon     string
	SoLuongBan int32
	SoLuongCon int32
	DaHuyMon   bool
	// DVang/DNgoc: so tien DA cong (duong hoac 0).
	DVang int64
	DNgoc int64
	// LaLapLai: goi nay la ban sao cua mot lan ban truoc (D3).
	LaLapLai bool
}

// banDonGia tra gia NPC mua lai mot mon: doc thang truong GiaBanLai cua mon
// hang, KHONG tinh nua.
//
// # CAI GI DOI O KH8 TASK 5
//
// Truoc task nay ham la `giaMua / 4` (san 1) — chep cong thuc ban Java
// (`cost /= 4; if (cost == 0) cost = 1;`). Nghia la muon ha gia thu lai cua
// RIENG mot mon thi phai doi cong thuc cua CA GAME. Nay gia ban lai la mot con
// so tuong minh trong kho-noi-dung/cua-hang.yaml, sua tung mon va nap nong duoc.
//
// Con so khoi dau trong YAML dung bang gia/4 cua cong thuc cu, nen ngay chuyen
// doi khong mot dong tien nao trong game doi.
//
// Bo kiem kho noi dung (noidung/nap.go kiemCuaHang) da khang dinh
// 0 <= GiaBanLai <= Gia truoc khi so lieu den duoc day.
func banDonGia(mon MonHang) int32 { return mon.GiaBanLai }

// banTinhToan lam TOAN BO phan kiem cua duong ban, KHONG cham CSDL de ghi.
//
//	ch      — cua hang DANG MO (may chu giu)
//	tui     — lat cat tui doc tu CSDL; bang phien dich o -> id duoc dung TU
//	          CHINH lat cat nay, nen khong the lech (xem hanhtrang.go)
//	oSo     — o thu ba cua goi `7`
//	soLuong — 0 nghia la BAN CA CHONG (giao thuc khong co o so luong; con so
//	          duong chi den tu test hoac tu duong goi tuong lai)
func banTinhToan(ch *CuaHang, tui []store.VatPham, oSo int16, soLuong int32) (banDuTinh, error) {
	if ch == nil {
		return banDuTinh{}, ErrKhongMoCuaHang
	}
	if oSo < 0 || int(oSo) >= tranSoO {
		return banDuTinh{}, fmt.Errorf("%w: %d", ErrOSoKhongHopLe, oSo)
	}

	// O -> vat_pham.id qua bang phien dich (task 4). Dung lat cat `tui` de dung
	// bang: dung mot lat cat cho ca id lan so luong, khong the tra o A ra mon B.
	bang, err := DungBangHanhTrang(tui)
	if err != nil {
		return banDuTinh{}, err
	}
	id, co := bang.TraId(store.NoiChuaTui, oSo)
	if !co {
		return banDuTinh{}, ErrBanOTrong
	}
	var vp *store.VatPham
	for i := range tui {
		if tui[i].ID == id {
			vp = &tui[i]
			break
		}
	}
	if vp == nil {
		// Khong the xay ra: bang duoc dung TU chinh `tui`.
		return banDuTinh{}, ErrBanOTrong
	}

	if soLuong == 0 {
		soLuong = vp.SoLuong // giao thuc: ban ca chong
	}
	if soLuong < 1 {
		return banDuTinh{}, fmt.Errorf("%w: %d", ErrSoLuongKhongHopLe, soLuong)
	}
	if soLuong > vp.SoLuong {
		return banDuTinh{}, fmt.Errorf("%w: ban %d, dang co %d",
			store.ErrKhongDuSoLuong, soLuong, vp.SoLuong)
	}

	if vp.MauID < 0 || vp.MauID > 32767 {
		return banDuTinh{}, fmt.Errorf("mau vat pham %d ngoai bien", vp.MauID)
	}
	mauID := int16(vp.MauID)
	mon := muaBanTimMon(ch, mauID)
	if mon == nil {
		return banDuTinh{}, fmt.Errorf("%w: mau %d", ErrBanKhongMuaMon, mauID)
	}
	mau := TimMauVatPham(mauID)
	if mau == nil {
		return banDuTinh{}, fmt.Errorf("mau vat pham %d khong ton tai", mauID)
	}

	donGia := banDonGia(*mon)
	tong := int64(donGia) * int64(soLuong)
	if tong < 0 {
		return banDuTinh{}, fmt.Errorf("tong gia am (%d)", tong)
	}
	dt := banDuTinh{
		VatPhamID: vp.ID, MauID: mauID, TenMon: mau.Ten,
		SoLuong: soLuong, DonGia: donGia, Tien: mon.Tien,
	}
	switch mon.Tien {
	case TienVang:
		dt.DVang = tong
	case TienNgoc:
		dt.DNgoc = tong
	default:
		return banDuTinh{}, fmt.Errorf("loai tien %d khong hop le", mon.Tien)
	}
	return dt, nil
}

// BanChoCuaHang thi hanh mot lan ban cho cua hang NPC.
//
// Moi cau kiem nam TRONG banTinhToan, tuc TRUOC loi goi CSDL duy nhat o cuoi
// ham — doi mot cau kiem xuong duoi la mo lai dung lo hong D16 cua ban Java.
func BanChoCuaHang(ctx context.Context, kho KhoBan, nvID int64, ch *CuaHang,
	tui []store.VatPham, oSo int16, soLuong int32, soThuTu uint64) (KetQuaBanHang, error) {

	dt, err := banTinhToan(ch, tui, oSo, soLuong)
	if err != nil {
		return KetQuaBanHang{}, err
	}

	// --- Chi TOI DAY moi cham CSDL de GHI. Mot giao dich duy nhat. ---
	kq, err := kho.BanMotGiaoDich(ctx, store.YeuCauBan{
		NhanVatID: nvID,
		KhoaIdem:  store.KhoaIdemBan(nvID, soThuTu),
		VatPhamID: dt.VatPhamID,
		NoiChua:   store.NoiChuaTui,
		OSo:       oSo,
		SoLuong:   dt.SoLuong,
		DVang:     dt.DVang,
		DNgoc:     dt.DNgoc,
		LyDo:      store.LyDoBanNPC,
	})
	if err != nil {
		return KetQuaBanHang{}, err
	}
	return KetQuaBanHang{
		VatPhamID: kq.VatPhamID, MauID: dt.MauID, TenMon: dt.TenMon,
		SoLuongBan: kq.SoLuongBan, SoLuongCon: kq.SoLuongCon, DaHuyMon: kq.DaHuyMon,
		DVang: kq.DVang, DNgoc: kq.DNgoc, LaLapLai: kq.LaLapLai,
	}, nil
}

// ---------------------------------------------------------------------------
// Goi XUONG `7` — hop thoai xac nhan
// ---------------------------------------------------------------------------

// goiBanXacNhan dung goi XUONG `7`: [where byte][oSo short][cau hoi UTF].
//
// `where`/`oSo` phai DOI LAI nguyen van cai client vua gui — client giu chung de
// gui lai o vong hai (Panel.cs idAction 3003). Tra nil neu cau hoi khong ghi
// duoc (writeUTF do dai 2 byte).
func goiBanXacNhan(where int8, oSo int16, cauHoi string) *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(byte(where))
	w.Short(oSo)
	if err := w.UTF(cauHoi); err != nil {
		return nil
	}
	return protocol.NewMessage(lenhBan, w)
}

// banTenTien doi loai tien thanh chu de ghep cau cho nguoi choi.
func banTenTien(t LoaiTien) string {
	if t == TienNgoc {
		return "ngọc"
	}
	return "vàng"
}

// banCauHoi dung cau trong hop thoai xac nhan, theo khuon ban Java
// (ShopService.showConfirmSellItem).
func banCauHoi(dt banDuTinh) string {
	tong := dt.DVang + dt.DNgoc
	return fmt.Sprintf("Bạn có muốn bán\nx%d %s\nvới giá là %d %s?",
		dt.SoLuong, dt.TenMon, tong, banTenTien(dt.Tien))
}

// banCauLoi doi loi ky thuat thanh cau nguoi choi doc duoc.
func banCauLoi(err error) string {
	switch {
	case errors.Is(err, ErrBanDoDangMac):
		return "Không thể bán đồ đang mặc"
	case errors.Is(err, ErrBanOTrong):
		return "Ô đó không có món nào"
	case errors.Is(err, ErrBanKhongMuaMon):
		return "Cửa hàng không mua món này"
	case errors.Is(err, store.ErrKhongDuSoLuong), errors.Is(err, ErrSoLuongKhongHopLe):
		return "Số lượng không hợp lệ"
	default:
		return "Không thể thực hiện"
	}
}

// ---------------------------------------------------------------------------
// Handler goi LEN `7`
// ---------------------------------------------------------------------------

// xuLyBan xu ly goi LEN `7` — CA HAI VONG.
//
//	action == 0 : xin xac nhan. KHONG doi gi trong CSDL, tra goi `7` xuong.
//	action != 0 : ban that. Tra -36 (tui moi) -> -25 (dong chu) -> 6 (tien).
//
// That bai (ca hai vong): -25 (ly do) -> 6 (tien KHONG doi). Goi `6` LUON la
// goi cuoi — no dong hop thoai cho phia client.
//
// Goi rac (thieu byte) thi KHONG tra loi gi va KHONG ngat ket noi.
func (b *BoXuLy) xuLyBan(s *mynet.Session, m *protocol.Message) error {
	if b.nhanVat == nil {
		return nil
	}
	r := protocol.NewReader(m.Data)
	action, err := r.Byte()
	if err != nil {
		return nil
	}
	where, err := r.Byte()
	if err != nil {
		return nil
	}
	oSo, err := r.Short()
	if err != nil {
		return nil
	}

	// Ban 1 chi ban do trong TUI (PHA A muc 4 dong 1).
	if int8(where) == banWhereMac {
		return b.muaBanGuiLoi(s, banCauLoi(ErrBanDoDangMac))
	}

	ctx := context.Background()
	kho, ok := b.kho.(KhoBan)
	if !ok {
		return b.muaBanGuiLoi(s, "Không thể thực hiện")
	}
	tui, err := kho.LayTui(ctx, b.nhanVat.ID, store.NoiChuaTui)
	if err != nil {
		return b.muaBanGuiLoi(s, "Không thể thực hiện")
	}

	// --- VONG MOT: xin xac nhan. Khong dong bo dem kinh te (chua ban gi ca). ---
	if action == 0 {
		dt, errTinh := banTinhToan(b.muaBanCuaHangDangMo, tui, oSo, 0)
		if errTinh != nil {
			return b.muaBanGuiLoi(s, banCauLoi(errTinh))
		}
		goi := goiBanXacNhan(int8(where), oSo, banCauHoi(dt))
		if goi == nil {
			return b.muaBanGuiLoi(s, "Không thể thực hiện")
		}
		return s.XepGui(goi)
	}

	// --- VONG HAI: ban that. soLuong = 0 -> ban CA CHONG (giao thuc khong co
	// o so luong). ---
	kq, errBan := BanChoCuaHang(ctx, kho, b.nhanVat.ID, b.muaBanCuaHangDangMo,
		tui, oSo, 0, b.soThuTuKinhTeKe())
	if errBan != nil {
		return b.muaBanGuiLoi(s, banCauLoi(errBan))
	}

	// Chi SAU KHI commit tra ve nil moi cham ban trong bo nho. Goi lap KHONG
	// duoc cong them lan nua — so tien trong ket qua la so cua lan dau, va lan
	// dau da cong roi.
	if !kq.LaLapLai {
		b.nhanVat.Vang += kq.DVang
		b.nhanVat.Ngoc += kq.DNgoc
	}

	// Doc lai tui tu CSDL (nguon su that) roi dung goi -36.
	tuiMoi, err := kho.LayTui(ctx, b.nhanVat.ID, store.NoiChuaTui)
	if err == nil {
		if goi, err := GoiTui(tuiMoi, soOTuiChuan); err == nil {
			if err := s.XepGui(goi); err != nil {
				return err
			}
		}
	}
	thu := kq.DVang + kq.DNgoc
	cau := fmt.Sprintf("Đã bán %s thu được %d %s", kq.TenMon, thu,
		banTenTien(banLoaiTienCua(kq)))
	if goi := goiMuaBanBao(cau); goi != nil {
		if err := s.XepGui(goi); err != nil {
			return err
		}
	}
	return s.XepGui(goiMuaBanTien(b.nhanVat.Vang, b.nhanVat.Ngoc))
}

// banLoaiTienCua doan loai tien tu ket qua ban: chi mot trong hai o khac 0.
func banLoaiTienCua(kq KetQuaBanHang) LoaiTien {
	if kq.DNgoc != 0 {
		return TienNgoc
	}
	return TienVang
}
