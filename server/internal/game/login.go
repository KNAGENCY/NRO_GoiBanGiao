// Package game lo nghiep vu choi. O ke hoach 1 moi chi co dang nhap.
package game

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
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
//
//	 -29  khung "chua dang nhap" — Controller.java:733, ben trong con mot
//	      byte lenh con nua (0 = dang nhap) — Controller.java:876
//	 -26  hien hop thoai co nut OK — Service.sendThongBaoOK
//	-102  bao client thoat trang thai "dang cho" — Service.sendLoginFail
//	   2  chuyen sang man tao nhan vat — Service.switchToCreateChar
const (
	lenhChuaDangNhap    int8 = -29
	lenhConDangNhap     byte = 0
	lenhThongBaoOK      int8 = -26
	lenhDangNhapThatBai int8 = -102
	lenhTaoNhanVat      int8 = 2
)

// reTenDangNhapHopLe khop dung quy tac ma CSDL da ap dung bang rang buoc CHECK
// tai_khoan_ten_hop_le (xem migrations/0001_tai_khoan.sql). Bien nay bien dich
// MOT LAN o muc goi (khong phai moi lan goi ham) — MustCompile la an toan o day
// vi bieu thuc chinh tay viet, khong doi luc chay.
//
// Vi sao phai kiem lai o tang Go dù CSDL da co CHECK: neu khong, mot ten dang
// nhap dai (vi du 60.000 ky tu) van bi ghi nguyen vao nhat ky roi day tron ven
// xuong CSDL truoc khi bi CHECK tu choi — cong toan bo qua trinh do la vo ich
// vi chac chan khong khop gi. Kiem truoc o Go tranh duoc ca hai lang phi do.
var reTenDangNhapHopLe = regexp.MustCompile(`^[A-Za-z0-9_]{4,20}$`)

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
	// Kiem dinh dang TRUOC khi cham toi kho: mot ten khong the nao khop se
	// khong bao gio tim thay tai khoan, nen goi kho la lang phi thuan tuy.
	// Tra ve CUNG mot cau bao chung (LyDoSaiThongTin) nhu moi truong hop sai
	// khac — KHONG duoc tao cau rieng cho truong hop nay, vi lam vay se lo
	// cho ke xau biet duoc luat dat ten (vd do dai, bang chu cai cho phep).
	if !reTenDangNhapHopLe.MatchString(ten) || matKhau == "" {
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
		// ket noi — client gui vai goi phu ngay sau khi noi (vi du
		// setClientType), day la chuyen BINH THUONG chu khong phai dang ngo.
		//
		// TUYET DOI KHONG ghi log o day, ke ca muc Debug: do bang thuc te tren
		// may chu that, mot ket noi don gui goi 3 byte lenh la co the dat
		// 32.452 goi/giay. Ghi MOT dong log (~67 byte) cho moi goi nhu vay
		// se khuech dai thanh 9 GB/GIO tren dia CHI TU MOT ket noi, khong gioi
		// han - day chinh la loi da tung xay ra (xem review-cuoi-report.md).
		return nil
	}

	r := protocol.NewReader(m.Data)
	lenhCon, err := r.Byte()
	if err != nil {
		return fmt.Errorf("goi -29 thieu byte lenh con: %w", err)
	}
	if lenhCon != lenhConDangNhap {
		// Cung ly do khong ghi log nhu nhanh "lenh chua ho tro" o tren: lenh
		// con la (vd setClientType) la binh thuong, ghi log se lap lai chinh
		// dieu kien gay ngap nhat ky da xay ra tren may chu that.
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
