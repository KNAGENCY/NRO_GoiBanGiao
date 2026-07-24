package game

import (
	"context"
	"errors"
	"log/slog"
	"regexp"

	mynet "nroserver/internal/net"
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// Lenh con cua khung -28 ma Ke hoach 2 xu ly. Cac lenh con khac (versionGame,
// clientOk...) la Ke hoach 3.
const lenhConTaoNhanVat byte = 2

// reTenNhanVat: 5-10 ky tu, chi chu cai tieng Anh va chu so. Khop dung rang
// buoc cua may chu Java cu (Controller.java:1043 + Util.haveSpecialCharacter).
// Bien dich mot lan o muc goi.
var reTenNhanVat = regexp.MustCompile(`^[A-Za-z0-9]{5,10}$`)

// danhSachTenCam: ten khong duoc dat de tranh gia mao vai tro quan tri. So
// sanh khong phan biet hoa thuong. Danh sach ngan, mo rong sau khi can.
var danhSachTenCam = map[string]bool{
	"admin": true, "administrator": true, "gm": true, "system": true, "quantri": true,
}

// xuLyTaoNhanVat xu ly than tin goi -28/2 (Reader da tro qua byte lenh con).
//
// KHONG tin client: phai kiem ten o day truoc khi cham CSDL, va phai kiem
// phien da dang nhap chua — client co the gui create truoc khi dang nhap.
func (b *BoXuLy) xuLyTaoNhanVat(s *mynet.Session, r *protocol.Reader) error {
	if b.taiKhoan == nil {
		// Chua dang nhap ma doi tao nhan vat: bo qua im lang, KHONG log (cung
		// ly do ngap nhat ky nhu cac lenh khong ho tro o login.go). Client
		// that khong bao gio lam vay; chi ke gia lap moi gui.
		return nil
	}
	if b.nhanVat != nil {
		return guiHopThoai(s, "Bạn đã có nhân vật.")
	}

	ten, err := r.UTF()
	if err != nil {
		// Than goi hong: bo ket noi (tra loi de tang net dong phien).
		return err
	}
	langByte, err := r.Byte()
	if err != nil {
		return err
	}
	kieuTocByte, err := r.Byte()
	if err != nil {
		return err
	}

	if !reTenNhanVat.MatchString(ten) {
		return guiHopThoai(s, "Tên nhân vật dài 5–10 ký tự, chỉ gồm chữ cái và chữ số.")
	}
	if danhSachTenCam[toLowerASCII(ten)] {
		return guiHopThoai(s, "Tên này không được phép dùng.")
	}

	lang := TimLang(int16(langByte))
	if lang == nil {
		return guiHopThoai(s, "Làng không hợp lệ.")
	}

	ctx := context.Background()
	daDung, err := b.kho.TenNhanVatDaDung(ctx, ten)
	if err != nil {
		slog.Error("loi kiem ten nhan vat", "ip", s.DiaChiIP, "loi", err)
		return guiHopThoai(s, "Lỗi hệ thống, vui lòng thử lại.")
	}
	if daDung {
		return guiHopThoai(s, "Tên nhân vật đã tồn tại.")
	}

	nv, err := b.kho.TaoNhanVat(ctx, store.NhanVatMoi{
		TaiKhoanID:  b.taiKhoan.ID,
		Ten:         ten,
		Lang:        lang.ID,
		KieuToc:     int16(kieuTocByte),
		MapID:       lang.MapKhoiDau,
		ViTriX:      lang.ViTriXKhoiDau,
		ViTriY:      lang.ViTriYKhoiDau,
		HpToiDa:     lang.HpKhoiDau,
		ChakraToiDa: lang.ChakraKhoiDau,
		Vang:        lang.VangKhoiDau,
	})
	if errors.Is(err, store.ErrTenNhanVatDaDung) {
		// Dua tranh: co nguoi khac vua lay ten nay giua luc kiem va luc chen.
		return guiHopThoai(s, "Tên nhân vật đã tồn tại.")
	}
	if errors.Is(err, store.ErrDaCoNhanVat) {
		return guiHopThoai(s, "Bạn đã có nhân vật.")
	}
	if err != nil {
		slog.Error("loi tao nhan vat", "ip", s.DiaChiIP, "loi", err)
		return guiHopThoai(s, "Lỗi hệ thống, vui lòng thử lại.")
	}

	slog.Info("tao nhan vat thanh cong",
		"ip", s.DiaChiIP, "ten", nv.Ten, "lang", nv.Lang, "tai_khoan", b.taiKhoan.ID)
	return b.vaoTheGioi(s, nv)
}

// guiHopThoai gui MOT hop thoai co nut OK (lenh -26). Khac guiTuChoi o
// login.go: guiTuChoi con gui kem -102 de thoat trang thai "dang cho" cua man
// DANG NHAP; o man TAO NHAN VAT, client tu tat vong xoay khi nhan -26
// (InfoDlg.hide), nen KHONG duoc gui -102 — gui thua se lam sai trang thai
// client.
func guiHopThoai(s *mynet.Session, cau string) error {
	w := protocol.NewWriter()
	if err := w.UTF(cau); err != nil {
		return err
	}
	return s.Gui(protocol.NewMessage(lenhThongBaoOK, w))
}

// toLowerASCII ha chu thuong CHI cho [A-Z]. Du dung vi ten nhan vat da qua
// reTenNhanVat (chi con [A-Za-z0-9]) — khong can xu ly unicode.
func toLowerASCII(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 'a' - 'A'
		}
	}
	return string(b)
}
