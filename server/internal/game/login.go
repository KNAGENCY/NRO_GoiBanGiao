// Package game lo nghiep vu choi. O ke hoach 1 moi chi co dang nhap.
package game

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"nroserver/internal/auth"
	mynet "nroserver/internal/net"
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// Cac cau tu choi. Sai mat khau va tai khoan khong ton tai dung CHUNG mot cau:
// bao khac nhau la lo cho ke xau biet tai khoan nao co that, tu do tap trung
// do mat khau dung tai khoan do. Ban Java cu cung lam vay, giu nguyen.
const (
	LyDoSaiThongTin = "Sai tài khoản hoặc mật khẩu."
	LyDoBiKhoa      = "Tài khoản đang bị khóa. Liên hệ quản trị viên."
)

// Lenh trong giao thuc. Cac gia tri nay lay tu ban Java cu, DUNG DOI:
//   -29  khung "chua dang nhap" — Controller.java:733, ben trong con mot
//        byte lenh con nua (0 = dang nhap) — Controller.java:876
//   -26  hien hop thoai co nut OK — Service.sendThongBaoOK
//  -102  bao client thoat trang thai "dang cho" — Service.sendLoginFail
//     2  chuyen sang man tao nhan vat — Service.switchToCreateChar
const (
	lenhChuaDangNhap    int8 = -29
	lenhConDangNhap     byte = 0
	lenhThongBaoOK      int8 = -26
	lenhDangNhapThatBai int8 = -102
	lenhTaoNhanVat      int8 = 2
)

type KetQuaDangNhap struct {
	ThanhCong bool
	LyDo      string
	TaiKhoan  *store.TaiKhoan
}

// KhoTaiKhoan tach ra thanh giao dien de kiem thu khong can CSDL that.
type KhoTaiKhoan interface {
	LayTaiKhoan(ctx context.Context, ten string) (*store.TaiKhoan, error)
}

type BoXuLy struct {
	kho KhoTaiKhoan
}

func NewBoXuLy(kho KhoTaiKhoan) *BoXuLy { return &BoXuLy{kho: kho} }

// KiemTraDangNhap chua toan bo luat dang nhap, tach khoi phan mang de kiem thu.
func (b *BoXuLy) KiemTraDangNhap(ctx context.Context, ten, matKhau string) (KetQuaDangNhap, error) {
	ten = strings.TrimSpace(ten)
	if ten == "" || matKhau == "" {
		return KetQuaDangNhap{LyDo: LyDoSaiThongTin}, nil
	}

	tk, err := b.kho.LayTaiKhoan(ctx, ten)
	if errors.Is(err, store.ErrKhongThayTaiKhoan) {
		return KetQuaDangNhap{LyDo: LyDoSaiThongTin}, nil
	}
	if err != nil {
		return KetQuaDangNhap{}, fmt.Errorf("tra cuu tai khoan that bai: %w", err)
	}

	if !auth.Verify(matKhau, tk.MatKhauBam) {
		return KetQuaDangNhap{LyDo: LyDoSaiThongTin}, nil
	}
	if tk.BiKhoa {
		return KetQuaDangNhap{LyDo: LyDoBiKhoa}, nil
	}
	return KetQuaDangNhap{ThanhCong: true, TaiKhoan: tk}, nil
}

// XuLy nhan goi tin da giai ma tu tang mang.
//
// Goi dang nhap long hai tang: khung -29, ben trong la byte lenh con roi
// moi toi du lieu. Doc thieu mot tang la boc tach ra rac.
func (b *BoXuLy) XuLy(s *mynet.Session, m *protocol.Message) error {
	if m.Command != lenhChuaDangNhap {
		// Ke hoach 1 chi lam dang nhap. Cac lenh khac bo qua, khong ngat
		// ket noi — client gui vai goi phu ngay sau khi noi.
		slog.Debug("bo qua lenh chua ho tro", "lenh", m.Command)
		return nil
	}

	r := protocol.NewReader(m.Data)
	lenhCon, err := r.Byte()
	if err != nil {
		return fmt.Errorf("goi -29 thieu byte lenh con: %w", err)
	}
	if lenhCon != lenhConDangNhap {
		slog.Debug("bo qua lenh con chua ho tro", "lenh_con", lenhCon)
		return nil
	}

	ten, err := r.UTF()
	if err != nil {
		return fmt.Errorf("goi dang nhap thieu ten: %w", err)
	}
	matKhau, err := r.UTF()
	if err != nil {
		return fmt.Errorf("goi dang nhap thieu mat khau: %w", err)
	}

	kq, err := b.KiemTraDangNhap(context.Background(), ten, matKhau)
	if err != nil {
		// Loi he thong: bao chung chung, khong lo chi tiet ky thuat ra ngoai.
		slog.Error("loi he thong khi dang nhap", "ip", s.DiaChiIP, "loi", err)
		return guiTuChoi(s, "Không đăng nhập được, vui lòng thử lại sau.")
	}
	if !kq.ThanhCong {
		slog.Info("dang nhap that bai", "ip", s.DiaChiIP, "ten", ten, "ly_do", kq.LyDo)
		return guiTuChoi(s, kq.LyDo)
	}

	slog.Info("dang nhap thanh cong",
		"ip", s.DiaChiIP, "ten", kq.TaiKhoan.TenDangNhap, "id", kq.TaiKhoan.ID)
	return guiChuyenSangTaoNhanVat(s)
}

// guiTuChoi bao loi roi bao client thoat khoi trang thai "dang cho".
//
// Ban Java cu tung QUEN buoc thu hai (goi -102), khien client treo o man hinh
// "Xin cho" cho toi khi tu thoat. Day la loi da tung xay ra that va da phai
// sua — giu ca hai buoc de khong lap lai.
func guiTuChoi(s *mynet.Session, cau string) error {
	w := protocol.NewWriter()
	if err := w.UTF(cau); err != nil {
		return err
	}
	if err := s.Gui(protocol.NewMessage(lenhThongBaoOK, w)); err != nil {
		return err
	}
	// Than tin 1 byte: 0 = khong con dang trong tien trinh dang nhap.
	return s.Gui(&protocol.Message{Command: lenhDangNhapThatBai, Data: []byte{0}})
}

// guiChuyenSangTaoNhanVat bao client mo man hinh tao nhan vat (lenh 2, than rong).
// Ke hoach 2 se thay bang: neu tai khoan da co nhan vat thi vao thang the gioi.
func guiChuyenSangTaoNhanVat(s *mynet.Session) error {
	return s.Gui(&protocol.Message{Command: lenhTaoNhanVat})
}
