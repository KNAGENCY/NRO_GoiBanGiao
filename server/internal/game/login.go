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
//	 -28  hop thu "khong tren ban do" — Controller.java messageNotMap, ben
//	      trong con mot byte lenh con nua (2 = tao nhan vat)
//	 -26  hien hop thoai co nut OK — Service.sendThongBaoOK
//	-102  bao client thoat trang thai "dang cho" — Service.sendLoginFail
//	   2  chuyen sang man tao nhan vat — Service.switchToCreateChar
const (
	lenhChuaDangNhap    int8 = -29
	lenhKhongTrenBanDo  int8 = -28 // hop thu "khong tren ban do" — Controller.java messageNotMap
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

// Kho gom moi thao tac CSDL ma tang game can. Tach thanh giao dien de kiem thu
// khong can CSDL that. (Ke hoach 1 chi co LayTaiKhoan; Ke hoach 2 them nhan vat.)
type Kho interface {
	LayTaiKhoan(ctx context.Context, ten string) (*store.TaiKhoan, error)
	LayNhanVatTheoTaiKhoan(ctx context.Context, taiKhoanID int64) (*store.NhanVat, error)
	TenNhanVatDaDung(ctx context.Context, ten string) (bool, error)
	TaoNhanVat(ctx context.Context, nv store.NhanVatMoi) (*store.NhanVat, error)
}

// BoXuLy giu trang thai cua MOT phien (mot ket noi). Xem Task 1: moi ket noi
// mot BoXuLy rieng, nen cac truong duoi day khong bao gio bi hai phien tranh.
type BoXuLy struct {
	kho      Kho
	taiKhoan *store.TaiKhoan // nil khi chua dang nhap
	nhanVat  *store.NhanVat  // nil khi chua co / chua vao the gioi

	clientSanSang bool    // client da gui -28/13 (bat tay xong)
	khuVuc        *KhuVuc // khu vuc nhan vat dang o (dung o KH4; nen mong cho KH5)
}

func NewBoXuLy(kho Kho) *BoXuLy { return &BoXuLy{kho: kho} }

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

// XuLy nhan goi tin da giai ma tu tang mang va dieu phoi theo lenh.
func (b *BoXuLy) XuLy(s *mynet.Session, m *protocol.Message) error {
	switch m.Command {
	case lenhChuaDangNhap: // -29
		return b.xuLyDangNhap(s, m)
	case lenhKhongTrenBanDo: // -28
		return b.xuLyKhongTrenBanDo(s, m)
	case lenhDuLieu: // -87 — client xin data
		return b.xuLyXinData(s)
	default:
		// Lenh chua ho tro: bo qua, khong ngat ket noi — client gui vai goi
		// phu ngay sau khi noi (vi du setClientType), day la chuyen BINH
		// THUONG chu khong phai dang ngo.
		//
		// TUYET DOI KHONG ghi log o day, ke ca muc Debug: do bang thuc te tren
		// may chu that, mot ket noi don gui goi 3 byte lenh la co the dat
		// 32.452 goi/giay. Ghi MOT dong log (~67 byte) cho moi goi nhu vay
		// se khuech dai thanh 9 GB/GIO tren dia CHI TU MOT ket noi, khong gioi
		// han - day chinh la loi da tung xay ra (xem review-cuoi-report.md).
		return nil
	}
}

// xuLyKhongTrenBanDo doc byte lenh con roi dieu phoi. Ke hoach 2 lam lenh con
// 2 (tao nhan vat); Ke hoach 3 them 6/7/8 (client xin catalog map/skill/item)
// va 13 (client bao san sang).
func (b *BoXuLy) xuLyKhongTrenBanDo(s *mynet.Session, m *protocol.Message) error {
	r := protocol.NewReader(m.Data)
	lenhCon, err := r.Byte()
	if err != nil {
		return fmt.Errorf("goi -28 thieu byte lenh con: %w", err)
	}
	switch lenhCon {
	case lenhConTaoNhanVat: // 2
		return b.xuLyTaoNhanVat(s, r)
	case lenhConCatalogMap: // 6
		return b.phucVuCatalog(s, goiCatalogMap())
	case lenhConCatalogSkill: // 7
		m, err := goiCatalogSkill()
		if err != nil {
			return err
		}
		return b.phucVuCatalog(s, m)
	case lenhConCatalogItem: // 8
		return b.phucVuCatalog(s, goiCatalogItem())
	case lenhConClientOk: // 13
		return b.xuLyClientOk(s)
	default:
		// Lenh con khac chua ho tro, bo qua im lang (cung ly do khong ghi
		// log nhu nhanh "lenh chua ho tro" o XuLy).
		return nil
	}
}

// xuLyXinData tra goi -87 khi client xin du lieu tinh (part...). Chi phuc vu
// khi phien da vao the gioi (vaoTheGioi da chay); truoc do bo qua im lang
// (khong log — cung ly do chong ngap nhat ky nhu nhanh "lenh chua ho tro").
func (b *BoXuLy) xuLyXinData(s *mynet.Session) error {
	if b.nhanVat == nil {
		return nil
	}
	return s.Gui(goiCatalogData())
}

// phucVuCatalog gui mot goi catalog khi client xin, chi khi da vao the gioi.
func (b *BoXuLy) phucVuCatalog(s *mynet.Session, m *protocol.Message) error {
	if b.nhanVat == nil {
		return nil
	}
	return s.Gui(m)
}

// xuLyClientOk: client bao da nhan du 4 catalog va san sang (-28/13). Ke hoach 4
// gui NGAY 3 goi de client render: -82 tileset -> -30/0 ME_INFO -> -24 MAP INFO.
// Thu tu quan trong: -82 PHAI truoc -24 (§7), khong thi client thoat han app.
func (b *BoXuLy) xuLyClientOk(s *mynet.Session) error {
	if b.nhanVat == nil {
		return nil
	}
	b.clientSanSang = true

	bd := TimBanDo(b.nhanVat.MapID)
	if bd == nil {
		// Nhan vat o mapId khong co trong danh muc ban do: khong the gui MAP INFO
		// hop le. Ghi log (SU KIEN HIEM — khong phai goi rac lap lai, nen log o day
		// KHONG vi pham luat chong ngap nhat ky) va dung. Client se ket "Xin cho"
		// nhung KHONG sap — tot hon la gui -24 rac lam client thoat app.
		slog.Error("nhan vat o mapId khong co ban do trong danh muc",
			"mapId", b.nhanVat.MapID, "ten", b.nhanVat.Ten, "id", b.nhanVat.ID)
		return nil
	}

	// Dung khu vuc toi thieu cho phien nay va dua nhan vat vao (nen mong cho KH5).
	b.khuVuc = MoKhuVuc(bd)
	b.khuVuc.ThemNguoiChoi(b.nhanVat)

	// -82 tileset — BAT BUOC truoc -24.
	if err := s.Gui(goiTileSetInfo()); err != nil {
		return err
	}

	// -30/0 ME_INFO.
	meInfo, err := goiMeInfo(b.nhanVat)
	if err != nil {
		return err
	}
	if err := s.Gui(meInfo); err != nil {
		return err
	}

	// -24 MAP INFO — gio kich switchToMe(), client roi "Xin cho".
	mapInfo, err := goiMapInfo(b.nhanVat, bd)
	if err != nil {
		return err
	}
	if err := s.Gui(mapInfo); err != nil {
		return err
	}

	slog.Info("da gui the gioi (tileset + ME_INFO + MAP INFO), client se render",
		"ten", b.nhanVat.Ten, "mapId", bd.MapID, "lang", b.nhanVat.Lang, "id", b.nhanVat.ID)
	return nil
}

// xuLyDangNhap xu ly khung -29 (dang nhap). Goi long hai tang: byte lenh con
// roi moi toi ten/mat khau. Doc thieu mot tang la boc tach ra rac.
func (b *BoXuLy) xuLyDangNhap(s *mynet.Session, m *protocol.Message) error {
	r := protocol.NewReader(m.Data)
	lenhCon, err := r.Byte()
	if err != nil {
		return fmt.Errorf("goi -29 thieu byte lenh con: %w", err)
	}
	if lenhCon != lenhConDangNhap {
		// Cung ly do khong ghi log nhu nhanh "lenh chua ho tro" o XuLy: lenh
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
	b.taiKhoan = kq.TaiKhoan // ghi nhan phien nay da dang nhap bang tai khoan nao

	nv, err := b.kho.LayNhanVatTheoTaiKhoan(context.Background(), kq.TaiKhoan.ID)
	if errors.Is(err, store.ErrKhongThayNhanVat) {
		return guiChuyenSangTaoNhanVat(s) // chua co -> mo man tao
	}
	if err != nil {
		slog.Error("loi he thong khi tai nhan vat", "ip", s.DiaChiIP, "loi", err)
		return guiTuChoi(s, "Không tải được nhân vật, vui lòng thử lại sau.")
	}
	return b.vaoTheGioi(s, nv) // da co -> vao the gioi (stub Ke hoach 2)
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
func guiChuyenSangTaoNhanVat(s *mynet.Session) error {
	return s.Gui(&protocol.Message{Command: lenhTaoNhanVat})
}

// vaoTheGioi la diem hen chung cho ca hai duong: dang nhap co san nhan vat, va
// vua tao nhan vat xong. Ke hoach 2 chi ghi nhan nhan vat vao phien va log; Ke
// hoach 3 them buoc gui goi -28/4 (versionGame) de kich client so 4 phien ban
// du lieu tinh voi cache cua no va tu xin cap nhat (xem dongbo.go).
func (b *BoXuLy) vaoTheGioi(s *mynet.Session, nv *store.NhanVat) error {
	b.nhanVat = nv
	slog.Info("nhan vat vao the gioi, bat dau bat tay phien ban",
		"ten", nv.Ten, "lang", nv.Lang, "id", nv.ID)
	// Gui versionGame de kich client so phien ban va tu xin 4 catalog. Client
	// gui -28/13 khi ca 4 phien ban khop; luc do Ke hoach 4 se gui ME_INFO + MAP
	// INFO. Ke hoach 3 dung o buoc client gui -28/13.
	return s.Gui(goiVersionGame())
}
