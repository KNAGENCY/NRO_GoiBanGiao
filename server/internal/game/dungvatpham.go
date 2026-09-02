package game

// dungvatpham.go — KH7 Task 7: BO DO (`-43`) va CHUYEN O (`-40`).
//
// Giai ma giao thuc (PHA A):
// docs/superpowers/research/2026-09-03-kh7-giai-ma-43-40.md
//
// PHAM VI BAN NAY — doc truoc khi trach "sao thieu":
//
//	LAM:   bo do (vut), mac trang bi, thao trang bi, chuyen o tui <-> ruong.
//	HOAN:  "dung" mon tieu hao (type 0/3 cua goi `-43`). Moi mon tieu hao cua
//	       ban 1 (Loai 6: thuoc hoi mau / chakra) deu doi HP hoac chakra — ma
//	       HP thuoc ACTOR KHU VUC (§6.9), phai phat su kien vao khuvuc.go.
//	       File do dang do agent khac giu (task 12), nen task nay KHONG cham.
//	       Goi type 0/3 den day duoc TRA LOI bang `-25` "Chua dung duoc vat
//	       pham nay" — KHONG im lang (§6.32: khong bao gio de nguoi choi bam
//	       ma khong co gi xay ra).
//
// BON DIEU CHOT — moi cai la mot dot bien bat buoc phai bi giet:
//
//  1. KIEM BIEN `oSo` O MOI NHANH, TRUOC khi tra bang (bay §7.11). `oSo` di
//     xuong day tu mot byte CO DAU phia client: `-1` la gia tri client CHU
//     DONG gui (Service.cs:418), va mot client sua doi gui duoc ca `-128`.
//     Moi nhanh kiem `0 <= oSo < SoO cua CHINH bang ma nhanh do doc` — khong
//     phai mot hang so chung, vi ba bang ba kich thuoc khac nhau.
//  2. SAU KHI CSDL thanh cong PHAI dung lai bang phien + gui lai ba goi tu
//     DUNG MOT lat cat (dungDoLamMoiVaGui). Bo buoc nay thi bang phien con
//     tro toi bo cuc CU: nguoi choi bam o 3 lan sau se dung nham mon khac —
//     khong loi, khong log (bay §7.10).
//  3. That bai thi KHONG DOI GI: mot cau `-25` va thoat. Khong gui goi hanh
//     trang nua voi.
//  4. Mac trang bi phai kiem `MauVatPham.Loai` — o mac la CHINH `Loai`
//     (0 ao, 1 quan, 2 gang, 3 giay, 4 nhan). Mon khong phai trang bi
//     (Loai >= SoOMacChuan) bi tu choi, khong duoc "nhet dai" vao o 0.
//
// BA DOT BIEN DA KIEM TAY (03/09) — moi cai lam DUNG mot test do:
//
//	1. Bo cau kiem bien trong dungDoKiemBien
//	   -> TestKiemBienPhanBietVoiOTrong +
//	      TestChuyenOMonNamNgoaiBienBangMacBiTuChoi do
//	2. Sau HuyVatPham thanh cong, khong goi dungDoLamMoiVaGui
//	   -> TestBoDoXongBangPhienKhongConO do
//	3. Bo cau kiem `mau.Loai >= SoOMacChuan` trong chuyenOMacTrangBi
//	   -> TestMacMonKhongPhaiTrangBiBiTuChoi do
//
// Chi tiet (va mot canh bao ve dot bien so 1) o dungvatpham_test.go.

import (
	"context"
	"errors"
	"fmt"

	"nroserver/internal/protocol"
	"nroserver/internal/store"

	mynet "nroserver/internal/net"
)

// ---------------------------------------------------------------------------
// 1. Hang so giao thuc (PHA A §1, §2)
// ---------------------------------------------------------------------------

const (
	// lenhDungBoDo: goi `-43`. CUNG mot ma cho ca chieu LEN (client xin dung /
	// vut) lan chieu XUONG (may chu hoi lai "chac chua?").
	lenhDungBoDo int8 = -43
	// lenhChuyenODo: goi LEN `-40` (chuyen o). KHONG co goi XUONG mang ma nay —
	// phan hoi la ba goi hanh trang -36/-37/-35.
	lenhChuyenODo int8 = -40
	// lenhDungDoBao: goi XUONG `-25` — mot dong chu chay, khong chan man hinh.
	lenhDungDoBao int8 = -25
)

// `type` cua goi `-43` (UseItem.java:52-55).
const (
	dungDoTypeDung        byte = 0 // client bam "dung"
	dungDoTypeVut         byte = 1 // client bam "vut"
	dungDoTypeXacNhanVut  byte = 2 // client bam "Co" tren hop thoai vut
	dungDoTypeXacNhanDung byte = 3 // client bam "Co" tren hop thoai dung
)

// `where` cua goi `-43` (Panel.cs:13374/13414).
const (
	dungDoNoiMac   byte = 0
	dungDoNoiTui   byte = 1
	dungDoNoiRuong byte = 2
)

// `loai` cua goi `-40` (Panel.cs:304-320 == UseItem.java:44-50).
// 6..9 la de tu / thu cung — ban 1 KHONG lam, roi vao nhanh mac dinh.
const (
	chuyenORuongSangTui byte = 0
	chuyenOTuiSangRuong byte = 1
	chuyenOMacSangRuong byte = 3
	chuyenOTuiSangMac   byte = 4
	chuyenOMacSangTui   byte = 5
)

// ---------------------------------------------------------------------------
// 2. Loi — moi cai mot cau tieng Viet cho nguoi choi
// ---------------------------------------------------------------------------

var (
	// ErrODoNgoaiBien: `oSo` client gui nam ngoai [0, so o cua bang do).
	// Day la loi cua bay §7.11 — KHONG BAO GIO duoc bo qua im lang.
	ErrODoNgoaiBien = errors.New("o so ngoai bien")
	// ErrONguonTrong: o nguon khong co mon nao (theo bang phien cua phien nay).
	ErrONguonTrong = errors.New("o nguon khong co mon")
	// ErrKhongPhaiTrangBi: mon nay khong mac len nguoi duoc (Loai ngoai 0..4).
	ErrKhongPhaiTrangBi = errors.New("mon nay khong phai trang bi")
	// ErrHetChoChua: bang dich khong con o trong nao.
	ErrHetChoChua = errors.New("khong con o trong")
	// ErrNoiChuaLa: `where` / `loai` ngoai bang da giai ma.
	ErrNoiChuaLa = errors.New("noi chua khong ho tro")
)

// ---------------------------------------------------------------------------
// 3. Kho — interface HEP, dung khuon voi KhoMuaBan (muaban.go)
// ---------------------------------------------------------------------------

// KhoDungVatPham la phan toi thieu cua *store.Store ma duong bo do / chuyen o
// can. Hep de test duoc bang kho gia va de doc mot cai la biet ham nay cham
// vao nhung gi.
type KhoDungVatPham interface {
	LayTui(ctx context.Context, nvID int64, noiChua string) ([]store.VatPham, error)
	ChuyenO(ctx context.Context, nvID int64, tuNoi string, tuO int16, denNoi string, denO int16) error
	HuyVatPham(ctx context.Context, nvID int64, noiChua string, oSo int16, lyDo string) error
}

// ---------------------------------------------------------------------------
// 4. Doi `where` / noi chua -> (ten noi chua, so o cua bang do)
// ---------------------------------------------------------------------------

// dungDoSoO tra SUC CHUA cua mot bang hanh trang. Day la ham DUY NHAT tra loi
// cau "bang nay bao nhieu o" — moi cho kiem bien deu goi no, nen khong the co
// hai cho dung hai con so khac nhau.
func dungDoSoO(noiChua string) (int, bool) {
	switch noiChua {
	case store.NoiChuaMac:
		return SoOMacChuan, true
	case store.NoiChuaTui:
		return SoOTuiChuan, true
	case store.NoiChuaRuong:
		return SoORuongChuan, true
	}
	return 0, false
}

// dungDoNoiChuaTu doi byte `where` cua goi `-43` thanh ten noi chua.
func dungDoNoiChuaTu(where byte) (string, bool) {
	switch where {
	case dungDoNoiMac:
		return store.NoiChuaMac, true
	case dungDoNoiTui:
		return store.NoiChuaTui, true
	case dungDoNoiRuong:
		return store.NoiChuaRuong, true
	}
	return "", false
}

// dungDoKiemBien la CUA DUY NHAT vao bang phien dich (bay §7.11).
//
// Ba viec, dung thu tu nay:
//  1. noi chua co that khong;
//  2. `oSo` co nam trong [0, so o cua CHINH bang do) khong — kiem TRUOC khi
//     tra bang, vi mot `oSo` = -1 hay 200 khong duoc phep cham toi CSDL;
//  3. o do co mon khong.
//
// Tra `vat_pham.id` cua mon o o do. Loi -> nguoi goi tu choi ngay, khong doan.
func (b *BoXuLy) dungDoKiemBien(noiChua string, oSo int16) (int64, error) {
	soO, co := dungDoSoO(noiChua)
	if !co {
		return 0, fmt.Errorf("%w: %q", ErrNoiChuaLa, noiChua)
	}
	if oSo < 0 || int(oSo) >= soO {
		return 0, fmt.Errorf("%w: o=%d, %s co %d o", ErrODoNgoaiBien, oSo, noiChua, soO)
	}
	id, coMon := b.hanhTrang.TraId(noiChua, oSo)
	if !coMon {
		return 0, fmt.Errorf("%w: %s o %d", ErrONguonTrong, noiChua, oSo)
	}
	return id, nil
}

// ---------------------------------------------------------------------------
// 5. Goi XUONG
// ---------------------------------------------------------------------------

// goiDungDoHoi dung goi XUONG `-43`: [byte type][byte where][byte oSo][UTF cau].
// Client mo hop thoai Co/Khong; bam "Co" no gui LEN lai `-43` voi type 2 (vut)
// hoac 3 (dung) — xem PHA A §1.2. Tra nil neu cau khong ghi duoc.
func goiDungDoHoi(loaiHoi byte, where byte, oSo int16, cau string) *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(loaiHoi)
	w.Byte(where)
	w.Byte(byte(int8(oSo)))
	if err := w.UTF(cau); err != nil {
		return nil
	}
	return protocol.NewMessage(lenhDungBoDo, w)
}

// goiDungDoBao dung goi XUONG `-25` (dong chu chay).
//
// VI SAO KHONG DUNG LAI goiMuaBanBao (muaban.go): cung noi dung byte, nhung
// file do thuoc duong MUA/BAN va co the doi theo task 10. Hai dong code trung
// nhau re hon mot phu thuoc cheo giua hai tinh nang khong lien quan.
func goiDungDoBao(noiDung string) *protocol.Message {
	w := protocol.NewWriter()
	if err := w.UTF(noiDung); err != nil {
		return nil
	}
	return protocol.NewMessage(lenhDungDoBao, w)
}

// dungDoCauLoi doi loi ky thuat thanh cau nguoi choi doc duoc. KHONG BAO GIO
// lo chi tiet ky thuat (ten bang, ten loi Go) ra man hinh.
func dungDoCauLoi(err error) string {
	switch {
	case errors.Is(err, ErrODoNgoaiBien), errors.Is(err, ErrNoiChuaLa):
		return "Không thể thực hiện"
	case errors.Is(err, ErrONguonTrong), errors.Is(err, store.ErrOTrong):
		return "Ô này không có gì"
	case errors.Is(err, ErrKhongPhaiTrangBi):
		return "Trang bị không phù hợp"
	case errors.Is(err, ErrHetChoChua), errors.Is(err, store.ErrOCoMon):
		return "Không còn ô trống"
	default:
		return "Không thể thực hiện"
	}
}

// dungDoBaoLoi gui MOT cau `-25` roi thoi — that bai thi KHONG doi gi va
// KHONG gui goi hanh trang nua voi (§2.3).
func (b *BoXuLy) dungDoBaoLoi(s *mynet.Session, cau string) error {
	goi := goiDungDoBao(cau)
	if goi == nil {
		return nil
	}
	return s.XepGui(goi)
}

// ---------------------------------------------------------------------------
// 6. Lam moi hanh trang — HAM DUNG CHUNG cho moi task sau
// ---------------------------------------------------------------------------

// dungDoLamMoiVaGui doc lai CA hanh trang tu CSDL (nguon su that), dung lai
// bang phien VA ba goi -37/-36/-35 tu DUNG MOT lat cat, gan bang moi vao phien
// roi gui ba goi.
//
// DAY LA HAM MA MOI TASK SAU (ban, giao dich, ky gui...) PHAI GOI sau khi doi
// bo cuc tui — dung viet lai. Bo qua no thi bang phien con tro toi bo cuc CU
// va thao tac ke tiep cua nguoi choi se cham nham mon (bay §7.10).
//
// Thu tu goi -37 (mac) -> -36 (tui) -> -35 (ruong) giu y het duong VAO THE GIOI
// (login.go) — client cap lai ca ba mang nen thu tu phai on dinh.
func (b *BoXuLy) dungDoLamMoiVaGui(s *mynet.Session) error {
	if b.nhanVat == nil {
		return nil
	}
	dsMon, err := b.hanhTrangNap(b.nhanVat.ID)
	if err != nil {
		return fmt.Errorf("nap lai hanh trang cua %d: %w", b.nhanVat.ID, err)
	}
	ht, err := DungBangVaGoi(dsMon, b.nhanVat.KieuToc, SoOTuiChuan, SoOMacChuan, SoORuongChuan)
	if err != nil {
		return fmt.Errorf("dung lai hanh trang cua %d: %w", b.nhanVat.ID, err)
	}
	b.hanhTrang = ht.Bang
	for _, g := range []*protocol.Message{ht.GoiMac, ht.GoiTui, ht.GoiRuong} {
		if err := s.XepGui(g); err != nil {
			return err
		}
	}
	return nil
}

// dungDoTimMon tim MOT mon trong mot lat cat theo (noi chua, o so).
func dungDoTimMon(dsMon []store.VatPham, noiChua string, oSo int16) *store.VatPham {
	for i := range dsMon {
		v := &dsMon[i]
		if v.NoiChua == noiChua && v.OSo != nil && *v.OSo == oSo {
			return v
		}
	}
	return nil
}

// dungDoTenMon tra ten mon o (noiChua, oSo) de ghep cau hoi/thong bao. Khong
// tra cuu duoc thi tra "vật phẩm" — cau van doc duoc, khong bao gio de trong.
func dungDoTenMon(dsMon []store.VatPham, noiChua string, oSo int16) string {
	v := dungDoTimMon(dsMon, noiChua, oSo)
	if v == nil {
		return "vật phẩm"
	}
	mau := TimMauVatPham(int16(v.MauID))
	if mau == nil || mau.Ten == "" {
		return "vật phẩm"
	}
	return mau.Ten
}

// ---------------------------------------------------------------------------
// 7. Handler goi LEN `-43` — dung / bo do
// ---------------------------------------------------------------------------

// xuLyDungBoDo xu ly goi LEN `-43`.
//
// Than goi: [byte type][byte where][byte oSo] (+ [short mauId] CHI KHI oSo ==
// -1). Ban nay KHONG doc truong short do — no chi xuat hien o nhanh "dung theo
// mau" ma ban nay hoan lai; doc thua mot truong khong dung toi la them mot cho
// de sai.
//
// Goi cut (thieu byte) -> KHONG tra loi gi va KHONG ngat ket noi, dung khuon
// xuLyDiChuyen / XuLyNhat.
func (b *BoXuLy) xuLyDungBoDo(s *mynet.Session, m *protocol.Message) error {
	if b.nhanVat == nil {
		return nil
	}
	r := protocol.NewReader(m.Data)
	loaiThaoTac, err := r.Byte()
	if err != nil {
		return nil
	}
	where, err := r.Byte()
	if err != nil {
		return nil
	}
	thoO, err := r.Byte()
	if err != nil {
		return nil
	}
	// Byte tren day CO DAU phia client (`sbyte`): -1 la gia tri hop le ma client
	// chu dong gui. Doi qua int8 TRUOC roi moi noi ra int16 de giu dau — doi
	// thang byte -> int16 se bien -1 thanh 255 va lot qua moi cau kiem bien.
	oSo := int16(int8(thoO))

	switch loaiThaoTac {
	case dungDoTypeDung, dungDoTypeXacNhanDung:
		// HOAN sang task sau (xem dau file). Tra loi tu te thay vi im lang.
		return b.dungDoBaoLoi(s, "Chưa dùng được vật phẩm này")

	case dungDoTypeVut:
		return b.dungDoHoiTruocKhiVut(s, where, oSo)

	case dungDoTypeXacNhanVut:
		return b.dungDoVutThat(s, where, oSo)

	default:
		return b.dungDoBaoLoi(s, "Không thể thực hiện")
	}
}

// dungDoHoiTruocKhiVut: buoc MOT cua vong hai buoc (PHA A §1.2). Kiem het roi
// moi hoi — hoi mot cau ve mot o rong la lam nguoi choi bam "Co" vao hu khong.
//
// KHONG cham CSDL de GHI o buoc nay.
func (b *BoXuLy) dungDoHoiTruocKhiVut(s *mynet.Session, where byte, oSo int16) error {
	noiChua, co := dungDoNoiChuaTu(where)
	if !co {
		return b.dungDoBaoLoi(s, "Không thể thực hiện")
	}
	// Ban 1 KHONG cho vut thang tu ruong — dung y ban Java (throwItem chi xu ly
	// where 0/1, InventoryService.java:116-132), nhung co cau tra loi chu khong
	// im lang.
	if noiChua == store.NoiChuaRuong {
		return b.dungDoBaoLoi(s, "Hãy lấy ra khỏi rương trước")
	}
	if _, err := b.dungDoKiemBien(noiChua, oSo); err != nil {
		return b.dungDoBaoLoi(s, dungDoCauLoi(err))
	}

	ten := "vật phẩm"
	if kho, ok := b.kho.(KhoDungVatPham); ok {
		ctx, huy := context.WithTimeout(context.Background(), hanhTrangHanCho)
		defer huy()
		if ds, err := kho.LayTui(ctx, b.nhanVat.ID, noiChua); err == nil {
			ten = dungDoTenMon(ds, noiChua, oSo)
		}
	}
	goi := goiDungDoHoi(dungDoTypeVut, where, oSo, "Bạn chắc chắn muốn bỏ "+ten+"?")
	if goi == nil {
		return nil
	}
	return s.XepGui(goi)
}

// dungDoVutThat: buoc HAI — thi hanh that.
//
// Trinh tu §2.3: kiem bien -> tra bang -> MOT loi goi CSDL -> chi khi no tra
// nil moi dung lai bang + gui goi.
func (b *BoXuLy) dungDoVutThat(s *mynet.Session, where byte, oSo int16) error {
	noiChua, co := dungDoNoiChuaTu(where)
	if !co {
		return b.dungDoBaoLoi(s, "Không thể thực hiện")
	}
	if noiChua == store.NoiChuaRuong {
		return b.dungDoBaoLoi(s, "Hãy lấy ra khỏi rương trước")
	}
	// KIEM BIEN TRUOC KHI TRA BANG — bay §7.11.
	if _, err := b.dungDoKiemBien(noiChua, oSo); err != nil {
		return b.dungDoBaoLoi(s, dungDoCauLoi(err))
	}
	kho, ok := b.kho.(KhoDungVatPham)
	if !ok {
		return b.dungDoBaoLoi(s, "Không thể thực hiện")
	}

	ctx, huy := context.WithTimeout(context.Background(), hanhTrangHanCho)
	defer huy()

	// Lay ten TRUOC khi huy — sau khi huy thi khong con doc duoc nua.
	ten := "vật phẩm"
	if ds, err := kho.LayTui(ctx, b.nhanVat.ID, noiChua); err == nil {
		ten = dungDoTenMon(ds, noiChua, oSo)
	}

	if err := kho.HuyVatPham(ctx, b.nhanVat.ID, noiChua, oSo, store.LyDoBo); err != nil {
		return b.dungDoBaoLoi(s, dungDoCauLoi(err))
	}
	// CHI SAU KHI commit tra nil moi dung lai bang + gui ba goi.
	if err := b.dungDoLamMoiVaGui(s); err != nil {
		return err
	}
	return b.dungDoBaoLoi(s, "Đã bỏ "+ten)
}

// ---------------------------------------------------------------------------
// 8. Handler goi LEN `-40` — chuyen o
// ---------------------------------------------------------------------------

// xuLyChuyenO xu ly goi LEN `-40`: [byte loai][byte oSo].
//
// Client KHONG gui o dich (PHA A §2.1) — may chu tu tim o trong nho nhat cua
// bang dich, tru nhanh "mac do": o mac la CHINH `MauVatPham.Loai` cua mon.
//
// Goi cut -> bo qua an toan, khong ngat ket noi.
func (b *BoXuLy) xuLyChuyenO(s *mynet.Session, m *protocol.Message) error {
	if b.nhanVat == nil {
		return nil
	}
	r := protocol.NewReader(m.Data)
	loai, err := r.Byte()
	if err != nil {
		return nil
	}
	thoO, err := r.Byte()
	if err != nil {
		return nil
	}
	oSo := int16(int8(thoO)) // giu dau, xem chu thich trong xuLyDungBoDo

	switch loai {
	case chuyenOTuiSangRuong:
		return b.chuyenOTuDongTimO(s, store.NoiChuaTui, oSo, store.NoiChuaRuong)
	case chuyenORuongSangTui:
		return b.chuyenOTuDongTimO(s, store.NoiChuaRuong, oSo, store.NoiChuaTui)
	case chuyenOMacSangRuong:
		return b.chuyenOTuDongTimO(s, store.NoiChuaMac, oSo, store.NoiChuaRuong)
	case chuyenOMacSangTui:
		return b.chuyenOTuDongTimO(s, store.NoiChuaMac, oSo, store.NoiChuaTui)
	case chuyenOTuiSangMac:
		return b.chuyenOMacTrangBi(s, oSo)
	default:
		// 6..9 la de tu / thu cung (ban 1 khong lam), 2 khong ton tai o ca hai
		// phia. Tra loi tu te thay vi im lang.
		return b.dungDoBaoLoi(s, "Không thể thực hiện")
	}
}

// chuyenOTuDongTimO chuyen mon tu (tuNoi, tuO) sang o TRONG NHO NHAT cua denNoi.
//
// KIEM BIEN `tuO` theo so o cua CHINH tuNoi — day la mot trong bon nhanh ma bay
// §7.11 noi toi, va la ly do ba bang khong duoc dung chung mot hang so.
func (b *BoXuLy) chuyenOTuDongTimO(s *mynet.Session, tuNoi string, tuO int16, denNoi string) error {
	if _, err := b.dungDoKiemBien(tuNoi, tuO); err != nil {
		return b.dungDoBaoLoi(s, dungDoCauLoi(err))
	}
	soODen, co := dungDoSoO(denNoi)
	if !co {
		return b.dungDoBaoLoi(s, "Không thể thực hiện")
	}
	kho, ok := b.kho.(KhoDungVatPham)
	if !ok {
		return b.dungDoBaoLoi(s, "Không thể thực hiện")
	}

	ctx, huy := context.WithTimeout(context.Background(), hanhTrangHanCho)
	defer huy()

	dsDich, err := kho.LayTui(ctx, b.nhanVat.ID, denNoi)
	if err != nil {
		return b.dungDoBaoLoi(s, "Không thể thực hiện")
	}
	// Doc nay nam NGOAI giao dich. Khong sao: neu giua doc va ghi co duong khac
	// chiem mat o do thi UNIQUE `vat_pham_mot_o_mot_mon` bat, ChuyenO tra
	// ErrOCoMon va CA giao dich rollback — nguoi choi thay loi, khong mat mon.
	oDich, conCho := timOTrong(dsDich, soODen)
	if !conCho {
		return b.dungDoBaoLoi(s, dungDoCauLoi(ErrHetChoChua))
	}

	if err := kho.ChuyenO(ctx, b.nhanVat.ID, tuNoi, tuO, denNoi, oDich); err != nil {
		return b.dungDoBaoLoi(s, dungDoCauLoi(err))
	}
	return b.dungDoLamMoiVaGui(s)
}

// chuyenOMacTrangBi: tui -> mac. Khac ba nhanh kia o CHO DUY NHAT: o dich
// KHONG phai "o trong nho nhat" ma la CHINH `MauVatPham.Loai` cua mon
// (0 ao, 1 quan, 2 gang, 3 giay, 4 nhan) — dung y ban Java
// (InventoryService.java:355-363, o mac = template.type).
//
// Mon khong phai trang bi (Loai ngoai [0, SoOMacChuan)) bi TU CHOI. Bo cau kiem
// nay thi mot chai thuoc (Loai 6) se bi nhet vao mot o mac ngoai bien hoac o 0
// — dot bien so 3, phai lam test do.
func (b *BoXuLy) chuyenOMacTrangBi(s *mynet.Session, tuO int16) error {
	if _, err := b.dungDoKiemBien(store.NoiChuaTui, tuO); err != nil {
		return b.dungDoBaoLoi(s, dungDoCauLoi(err))
	}
	kho, ok := b.kho.(KhoDungVatPham)
	if !ok {
		return b.dungDoBaoLoi(s, "Không thể thực hiện")
	}

	ctx, huy := context.WithTimeout(context.Background(), hanhTrangHanCho)
	defer huy()

	tui, err := kho.LayTui(ctx, b.nhanVat.ID, store.NoiChuaTui)
	if err != nil {
		return b.dungDoBaoLoi(s, "Không thể thực hiện")
	}
	mon := dungDoTimMon(tui, store.NoiChuaTui, tuO)
	if mon == nil {
		return b.dungDoBaoLoi(s, dungDoCauLoi(ErrONguonTrong))
	}
	mau := TimMauVatPham(int16(mon.MauID))
	if mau == nil {
		return b.dungDoBaoLoi(s, "Không thể thực hiện")
	}
	// KIEM `Loai`: chi 0..4 moi mac duoc. `SoOMacChuan` la con so DUY NHAT
	// quyet dinh — khong ghim 5 o day, de doi so o mac khong phai sua hai cho.
	if mau.Loai < 0 || int(mau.Loai) >= SoOMacChuan {
		return b.dungDoBaoLoi(s, dungDoCauLoi(ErrKhongPhaiTrangBi))
	}
	oMac := int16(mau.Loai)

	// O mac do da co mon -> ChuyenO tra ErrOCoMon (kho KHONG hoan doi, xem
	// vatpham_kho.go). Nguoi choi phai thao mon cu ra truoc. Ban 1 chap nhan
	// nhu vay va noi ro thay vi im lang.
	if err := kho.ChuyenO(ctx, b.nhanVat.ID, store.NoiChuaTui, tuO, store.NoiChuaMac, oMac); err != nil {
		if errors.Is(err, store.ErrOCoMon) {
			return b.dungDoBaoLoi(s, "Hãy tháo trang bị cũ ra trước")
		}
		return b.dungDoBaoLoi(s, dungDoCauLoi(err))
	}
	return b.dungDoLamMoiVaGui(s)
}
