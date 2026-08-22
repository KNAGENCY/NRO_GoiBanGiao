// Package game lo nghiep vu choi. O ke hoach 1 moi chi co dang nhap.
package game

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"sync"

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
	lenhVaoKhuVuc       int8 = -39 // client bao da tai xong map
	lenhXinTaiNguyen    int8 = -74 // client xin resource o man ServerListScreen — Service.getResource
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

// GhiKhiRoi la phan toi thieu cua ThoGhi (thogi.go) ma BoXuLy can de luu
// nguoi choi khi phien mat ket noi (Ke hoach 11 Task 5, dac ta
// 2026-07-27-luu-trang-thai-song.md §4 diem B) — dung khuon voi Kho o tren va
// KhoLuu (thogi.go): interface HEP chi doi hoi dung phuong thuc can dung
// (Xep), de kiem thu duoc bang mot "tho ghi gia" ghi nhan NGAY vao bo nho,
// khong phai dung *ThoGhi that (di qua goroutine nen + ticker that se bat
// test phai cho/doan thoi gian, de flaky). *ThoGhi thoa man interface nay tu
// nhien, khong can glue code them.
type GhiKhiRoi interface {
	Xep(nv store.NhanVat)
}

// ghiDongBoKhiBiDa la phan MO RONG khong xuat (Ke hoach 11 Task 7, dac ta
// §7.1) ma KhiDong() can KHI PHIEN NAY VUA BI DA boi mot dang nhap moi cung
// tai khoan: ghi DONG BO, cho ket qua that su truoc khi tra ve — khac han
// Xep() (GhiKhiRoi o tren, KHONG BAO GIO CHAN, chi cham hang doi trong bo
// nho). *ThoGhi thoa man interface nay tu nhien (XepVaCho, thogi.go).
//
// TACH RIENG khoi GhiKhiRoi (khong nhet XepVaCho vao do) de MOI fake thoGhi
// hien co trong cac test KHONG lien quan kich ban "bi da" (chi implement
// Xep) khong phai sua gi ca — KhiDong() type-assert interface nay CHI khi
// can (xem daBiDa()), va rot ve nhanh Xep() binh thuong neu thoGhi cu the
// khong thoa man (an toan tuyet doi, khong panic — xem KhiDong).
type ghiDongBoKhiBiDa interface {
	XepVaCho(ctx context.Context, nv store.NhanVat) error
}

// BoXuLy giu trang thai cua MOT phien (mot ket noi). Xem Task 1: moi ket noi
// mot BoXuLy rieng, nen cac truong duoi day khong bao gio bi hai phien tranh
// — NGOAI TRU session/biDa/daXongDon (Ke hoach 11 Task 7), xem chu thich rieng
// cua tung truong do ve ly do CHUNG DUOC mot goroutine KHAC cham vao va vi
// sao van AN TOAN (chi qua thao tac dong-kenh/Session.Dong(), khong bao gio
// doc/ghi truc tiep cac truong con lai cua BoXuLy tu ben ngoai).
type BoXuLy struct {
	kho      Kho
	taiKhoan *store.TaiKhoan // nil khi chua dang nhap
	nhanVat  *store.NhanVat  // nil khi chua co / chua vao the gioi

	clientSanSang bool      // client da gui -28/13 (bat tay xong)
	daVaoKhuVuc   bool      // da gui suKienVao (chan -39 goi lap gay -5 trung)
	khuVuc        *KhuVuc   // khu vuc nhan vat dang o (dung o KH4; nen mong cho KH5)
	soKhuVuc      *SoKhuVuc // so khu vuc CHUNG (5 map), tra khu vuc theo mapId cua nhan vat

	// thoGhi (Ke hoach 11 Task 5): noi KhiDong() Xep() anh chup nguoi choi vao
	// khi phien mat ket noi (dac ta §4 diem B). nil la gia tri mac dinh HOP LE
	// (chua goi DatThoGhi) — KhiDong tu bo qua buoc luu trong truong hop do,
	// KHONG panic (xem KhiDong). Tach khoi tham so constructor (xem DatThoGhi
	// ben duoi) THAY VI them tham so thu ba vao NewBoXuLy, vi rat nhieu noi
	// trong repo dang goi NewBoXuLy(kho, soKhuVuc) voi DUNG 2 tham so
	// (nhanvat_test.go, tichhop_test.go, khuvuc_tichhop_test.go...) khong lien
	// quan gi toi luu tru — doi chu ky constructor se bat het cac noi do phai
	// sua theo, khong can thiet cho task nay.
	thoGhi GhiKhiRoi

	// session (Ke hoach 11 Task 7, dac ta §7.2): luu lai *mynet.Session cua
	// CHINH phien nay, gan MOT LAN trong xuLyDangNhap khi dang nhap thanh cong
	// — TRUOC khi dang ky vao soDangNhap (xem xuLyDangNhap), de neu MOT dang
	// nhap KHAC cung tai khoan xay ra gan nhu ngay sau do, SoDangNhap.DangKy
	// (chay tren goroutine cua phien MOI do) da co the DONG duoc ket noi vat
	// ly cua phien nay ngay. Day la truong DUY NHAT cua BoXuLy duoc DOC boi
	// mot goroutine KHAC voi goroutine so huu phien — AN TOAN vi (1) chi GAN
	// MOT LAN (khong bao gio doi lai sau do) va (2) Session.Dong() (net/
	// session.go) tu no da AN TOAN goi dong thoi tu bat ky goroutine nao
	// (sync.Once noi bo) — khong dua tranh gi voi cach chinh BoXuLy nay dung
	// no (Gui() da tu khoa muGui rieng).
	session *mynet.Session

	// biDa la kenh TIN HIEU (Ke hoach 11 Task 7): dong boi SoDangNhap.DangKy
	// (chay tren goroutine cua phien MOI) de bao PHIEN NAY vua bi DA boi mot
	// dang nhap moi cung tai khoan. Luon khoi tao khac nil (NewBoXuLy) de an
	// toan dong tu ben ngoai. Dong mot kenh la thao tac cross-goroutine AN
	// TOAN theo dung dinh nghia cua Go (khong phai doc/ghi mot truong thuong)
	// — day KHONG vi pham nguyen tac "moi truong BoXuLy chi mot phien cham
	// toi" o tren, vi kenh CHINH LA co che danh cho muc dich nay. CHI duoc
	// DOC (khong bao gio ghi) tu goroutine cua CHINH phien nay, trong
	// daBiDa() (goi ben trong KhiDong()).
	biDa chan struct{}

	// biDaMotLan (soi toan nhanh KH11, fix-cuoi-report.md, Finding 1) la
	// PHONG THU CHIEU SAU cho close(biDa): dam bao close(b.biDa) o
	// SoDangNhap.DangKy (sodangnhap.go) chi thuc su chay MOT LAN cho DU
	// DangKy vo tinh (hoac do mot loi khac trong tuong lai) goi den hai lan
	// tren CUNG mot BoXuLy — close mot channel DA dong la panic khong
	// recover duoc tren duong PhucVu (net/session.go khong co recover()),
	// panic do se lam CHET CA TIEN TRINH may chu. Chot chinh (b.taiKhoan !=
	// nil trong xuLyDangNhap) da chan duong gay ra tinh huong nay TU GOC —
	// truong nay la LUOI AN TOAN THEM, khong phai thay the cho chot chinh.
	biDaMotLan sync.Once

	// daXongDon la kenh TIN HIEU nguoc lai voi biDa: dong boi CHINH goroutine
	// cua phien nay, o CUOI KhiDong() (qua defer baoXongDon) — bao "KhiDong da
	// chay xong HOAN TOAN" (ke ca buoc ghi DONG BO neu phien bi da) cho
	// SoDangNhap.DangKy (dang cho tren goroutine cua phien MOI) biet luc nao
	// AN TOAN de tiep tuc doc lai CSDL — day CHINH LA co che thuc hien dac ta
	// §7.1 "cho no luu xong -> ROI MOI doc du lieu cho phien moi".
	daXongDon chan struct{}

	// soDangNhap (Ke hoach 11 Task 7): so phien CHUNG toan may chu — MOT nguon
	// su that duy nhat "tai khoan X dang o phien nao" (xem sodangnhap.go). nil
	// la gia tri mac dinh AN TOAN (khong dang ky/khong bao gio bi da) — cac
	// test hien co (goi NewBoXuLy(kho, soKhuVuc) truc tiep, khong goi
	// DatSoDangNhap) van bien dich va chay dung y het truoc Task 7.
	soDangNhap *SoDangNhap
}

func NewBoXuLy(kho Kho, soKhuVuc *SoKhuVuc) *BoXuLy {
	return &BoXuLy{
		kho: kho, soKhuVuc: soKhuVuc,
		// Ke hoach 11 Task 7: LUON khoi tao hai kenh nay khac nil, du phien
		// nay co bao gio bi da hay khong — an toan de SoDangNhap dong biDa tu
		// ben ngoai (dong kenh nil se panic), va daXongDon luon co the dong
		// duoc trong KhiDong() (goi mot lan duy nhat, xem baoXongDon).
		biDa:      make(chan struct{}),
		daXongDon: make(chan struct{}),
	}
}

// DatSoDangNhap gan so dang nhap CHUNG cho phien nay (Ke hoach 11 Task 7) —
// goi TU main.go (qua MoLuuTru) SAU NewBoXuLy, dung khuon "thiet lap sau khoi
// tao" giong het DatThoGhi o tren (xem chu thich truong thoGhi ve ly do tach
// khoi constructor). AN TOAN KHONG goi ham nay (soDangNhap giu nguyen nil) —
// phien se khong dang ky vao so va khong bao gio bi da.
func (b *BoXuLy) DatSoDangNhap(s *SoDangNhap) {
	b.soDangNhap = s
}

// DatThoGhi gan tho ghi cho phien nay (Ke hoach 11 Task 5) — goi TU main.go
// SAU NewBoXuLy, dung khuon "thiet lap sau khoi tao" thay vi tham so
// constructor thu ba (xem ly do o chu thich truong thoGhi phia tren). An toan
// KHONG goi ham nay (thoGhi giu nguyen nil) — KhiDong() tu bo qua buoc luu.
func (b *BoXuLy) DatThoGhi(t GhiKhiRoi) {
	b.thoGhi = t
}

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
	case lenhXinTaiNguyen: // -74 — client xin resource tu server (§ScrList)
		return b.xuLyXinTaiNguyen(s, m)
	case lenhVaoKhuVuc: // -39 — client bao da tai xong map
		return b.xuLyVaoKhuVuc(s)
	case lenhDiChuyen: // -7 — client xin di chuyen
		return b.xuLyDiChuyen(m)
	case lenhChonChieu: // 34 — client chon chieu
		return b.xuLyChonChieu(m)
	case lenhNguoiDanhQuai: // 54 — client gui danh quai (dac ta §3.1)
		return b.xuLyDanhQuai(m)
	case lenhNguoiDanhNguoi: // -60 — client gui danh nguoi choi (dac ta §3.1, Ke hoach 6 Task 11)
		return b.xuLyDanhNguoiChoi(m)
	case lenhHoiSinh: // -16 — client xin hoi sinh (dac ta §3.1, Ke hoach 6 Task 13)
		return b.xuLyHoiSinh()
	case lenhMoNPC: // 33 — client mo hoi thoai NPC
		return b.xuLyMoNPC(s, m)
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

// xuLyXinTaiNguyen tra goi -74 stub khi client o ServerListScreen goi
// Service.getResource(1, null) — client kiem "phien ban resource" cua may
// chu roi so voi RMS "ResVersion". Neu khop -> di tiep sang LoginScr; neu
// khac -> hien dialog xin download. Vi may chu Go MOI khong host bo
// resource cu (do la trach nhiem cua bo cai client, khong cua game server),
// tra ResVersion = 0 CO DINH, VA yeu cau ban Unity o may nay cung ghi RMS
// "ResVersion" = "0" -> khop, client di tiep binh thuong.
//
// Format client gui (Service.cs:1936): 40 byte token (HMAC-SHA256 anti-bot)
// + 1 byte action (1/2/3, xem Controller.cs:1133 case -74). Ta CHI can
// dieu action, khong xet token (khong phai muc dich chong bot cua ke hoach
// nay). action=1: xin danh sach resource (tra -74/0/int0). action=2/3: bo
// qua im lang — client khong cho tra loi cho hai action nay.
//
// Ban KHONG CAN doi trang thai phien (b.nhanVat = nil hay khac): -74 chay
// TRUOC login, khi phien chua co nhan vat. Kiem tra o day chi phi vi
// duong voi (client mo ket noi la vao ServerListScreen -> gui -74 luon).
func (b *BoXuLy) xuLyXinTaiNguyen(s *mynet.Session, m *protocol.Message) error {
	// Body toi thieu = 40 byte token + 1 byte action = 41 byte. Ngan hon
	// la client hong hoac lenh bi trang phong, bo qua im lang (khong log,
	// cung ly do chong ngap nhu nhanh "lenh chua ho tro" o XuLy).
	if len(m.Data) < 41 {
		return nil
	}
	action := m.Data[40] // byte action ngay sau 40 byte token
	if action != 1 {
		// action=2 (xin resource cu the) va action=3 (bao xong): may chu
		// khong host resource -> khong tra gi, client tu xu ly khong cho
		// (khong treo, xem Controller.cs case -74 khong co else nhanh
		// hong cho action nay).
		return nil
	}
	w := protocol.NewWriter()
	w.Byte(0) // b42 = 0 (danh sach phien ban resource)
	w.Int(1)  // ResVersion = 1 — PHAI khop RMS "ResVersion" da tao thu cong tren
	// may client (writeUTF "1"). Neu tra ve gia tri khac, client so lech ->
	// vao nhanh show2() (hop thoai xin download) thay vi di tiep binh thuong.
	// Vi sao 1 chu khong phai 0: switchToMe cua ServerListScreen (line 660)
	// chi set loadScreen=true khi int.Parse(RMS "ResVersion") > 0. Voi "0"
	// loadScreen mai false, client ket o gif logo, khong bao gio hien 4 nut
	// (Dang nhap/Doi tk/Server/Option) de nguoi choi bam.
	return s.Gui(protocol.NewMessage(lenhXinTaiNguyen, w))
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
	// Chi chay MOT LAN moi phien: neu chua co nhan vat, hoac da san sang roi thi bo
	// qua. Chan tai nhap quan trong vi sau khi -39 giao nv cho actor, chay lai
	// xuLyClientOk se khien goiMapInfo DOC nv.ViTriX/Y tren luong phien trong khi
	// actor dang GHI (tranh chap). Client that chi gui -28/13 mot lan.
	if b.nhanVat == nil || b.clientSanSang {
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

	// Ke hoach 5: dung khu vuc CHUNG tu so (nhieu nguoi cung mapId thay nhau).
	// CHUA gui suKienVao o day — doi -39 (client tai xong map) moi dua vao khu vuc.
	if b.soKhuVuc != nil {
		b.khuVuc = b.soKhuVuc.Tra(b.nhanVat.MapID)
	}

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

	// -24 MAP INFO — gio kich switchToMe(), client roi "Xin cho". Danh sach quai
	// (Ke hoach 6 Task 5): lay tu KhuVuc NEU da co (b.khuVuc vua gan o tren),
	// qua DanhSachQuaiSnapshot() — ham nay di qua hop thu actor (dac ta §1.6:
	// chi vongLap duoc cham trang thai quai) roi tra ve BAN SAO, khong doc
	// truc tiep truong quai cua KhuVuc tu luong phien nay. Buoc nay chay
	// TRUOC khi nhan vat vao khu vuc (-39, xuLyVaoKhuVuc) nhung actor CUA KHU
	// VUC da chay san tu luc may chu khoi dong (MoSoKhuVuc goi BatDau() ngay),
	// nen van co the co nguoi khac dang chien dau lam quai doi HP dong thoi —
	// buoc nay PHAI an toan voi truong hop do.
	var quai []*Quai
	if b.khuVuc != nil {
		quai = b.khuVuc.DanhSachQuaiSnapshot()
	}
	mapInfo, err := goiMapInfo(b.nhanVat, bd, quai)
	if err != nil {
		return err
	}
	if err := s.Gui(mapInfo); err != nil {
		return err
	}

	// -42 CHI SO: mau, chakra, suc danh. Gui SAU MAP INFO vi client chi bat dau
	// ve HUD khi da vao man choi. Thieu goi nay thi thanh mau va thanh chakra
	// nam im o 0 — khong phai loi hien thi, ma la client that su khong biet.
	if err := s.Gui(goiChiSoNhanVat(b.nhanVat)); err != nil {
		return err
	}

	slog.Info("da gui the gioi (tileset + ME_INFO + MAP INFO + chi so), client se render",
		"ten", b.nhanVat.Ten, "mapId", bd.MapID, "lang", b.nhanVat.Lang, "id", b.nhanVat.ID)
	return nil
}

// xuLyVaoKhuVuc: client gui -39 (da tai xong map) -> dua nhan vat vao khu vuc chung
// (actor phat -5 cua nguoi cu cho minh, va phat -5 cua minh cho nguoi cu).
//
// Chi chay MOT LAN moi phien (chot b.daVaoKhuVuc): goi lap -39 se khien actor xu ly
// LAI suKienVao, va nguoi "cu" duy nhat trong buoc phat -5 luc do chinh la nguoi nay
// (da duoc them vao khu vuc boi lan -39 truoc) — actor se gui cho ho goi -5 CUA
// CHINH HO, mot goi trung thua tren day. Ke hoach 5 chua co doi map; khi them doi
// map se can reset daVaoKhuVuc luc roi khu vuc — hien mot phien = mot lan vao.
func (b *BoXuLy) xuLyVaoKhuVuc(s *mynet.Session) error {
	if b.khuVuc == nil || b.nhanVat == nil || b.daVaoKhuVuc {
		return nil
	}
	b.daVaoKhuVuc = true
	b.khuVuc.Gui(suKienVao{nv: b.nhanVat, phien: s})
	return nil
}

// xuLyMoNPC chi doc ID va tra dialog bat bien. Khong dong vao state KhuVuc;
// actor van la chu duy nhat cua di chuyen va chien dau.
func (b *BoXuLy) xuLyMoNPC(s *mynet.Session, m *protocol.Message) error {
	if b.nhanVat == nil || b.nhanVat.MapID != 1 {
		return nil
	}
	r := protocol.NewReader(m.Data)
	id, err := r.Short()
	if err != nil || id < 0 || id > 255 {
		return nil
	}
	goi, err := goiMoNPC(byte(id))
	if err != nil {
		return nil
	}
	// Hoi thoai NPC la gameplay sau khi vao map: vao hang gui FIFO nhu cac
	// broadcast cua actor, khong ghi socket dong bo tren luong xu ly command.
	return s.XepGui(goi)
}

// xuLyDiChuyen: client gui -7 (byte co bay + short toX + short toY). Ke hoach 5:
// KHONG kiem hop le o day (doc nv.ViTri se dua voi actor ghi) — chuyen tiep toX/toY
// tho cho actor, actor kiem bang vi tri authoritative. Thieu toX/toY -> bo goi.
func (b *BoXuLy) xuLyDiChuyen(m *protocol.Message) error {
	if b.khuVuc == nil || b.nhanVat == nil {
		return nil
	}
	r := protocol.NewReader(m.Data)
	if _, err := r.Byte(); err != nil { // co bay (khinh cong) — KH5 doc bo qua
		return nil
	}
	toX, err := r.Short()
	if err != nil {
		return nil
	}
	toY, err := r.Short()
	if err != nil {
		return nil // KH5: thieu toY -> bo goi (client that luon gui du ca hai)
	}
	b.khuVuc.Gui(suKienDiChuyen{idNhanVat: b.nhanVat.ID, x: toX, y: toY})
	return nil
}

// xuLyChonChieu: client gui 34 (short templateId — day la MauChieu.Id, KHONG
// PHAI skillId, xem dac ta §3.2). Luong phien CHI doc gói va CHUYEN TIEP mot
// su kien cho actor khu vuc; TUYET DOI KHONG tu tra ChieuCuaLang/ghi
// chieuDangChon o day — do la trach nhiem cua actor (khuvuc.go: xuLy case
// suKienChonChieu), vi luong phien ghi + actor doc/ghi CUNG mot truong se
// tranh chap du lieu (bay §6.9, KH5 da vap voi nv.ViTriX/Y).
func (b *BoXuLy) xuLyChonChieu(m *protocol.Message) error {
	if b.khuVuc == nil || b.nhanVat == nil {
		return nil
	}
	r := protocol.NewReader(m.Data)
	templateId, err := r.Short()
	if err != nil {
		return nil // goi ngan hon khuon -> bo qua an toan, khong ngat ket noi
	}
	// templateId chi co gia tri 1..4 tren thuc te (4 slot chieu/lang, xem
	// chieuthuc.go) nen luon vua byte; ep kieu ve byte de khop kieu
	// suKienChonChieu.mauChieuId (MauChieu.Id la byte trong catalog -28/7).
	// Gia tri gia mao ngoai khoang byte se wrap nhung khong khop Id nao trong
	// ChieuCuaLang nen actor tu bo qua (xem xuLy case suKienChonChieu).
	b.khuVuc.Gui(suKienChonChieu{idNhanVat: b.nhanVat.ID, mauChieuId: byte(templateId)})
	return nil
}

// xuLyDanhQuai xu ly goi LEN `54` (dac ta §3.1): client bao muon danh mot con
// quai. THAN GOI DO DAI BIEN THIEN — day la goi RUI RO SO MOT cua Ke hoach 6
// (dac ta §2): `byte mobId`; NEU mobId == -1 (thu cung, ngoai pham vi KH6 —
// dac ta §0 "Ngoai KH6") thi con THEM `int masterId` (4 byte), tong 5 byte.
// PHAI doc DU ca hai truong hop de goi tin nay duoc tieu thu dung so byte no
// khai bao — khong doc du la nguyen nhan hang dau gay lech dong byte theo dac
// ta, va la ly do Test 3 cua Task 7 kiem rieng kich ban nay.
//
// Ky luat actor (dac ta §1.6 + brief Task 7): ham nay CHI doc gia tri mobId
// tu goi roi gui suKienDanhQuai vao hop thu actor khu vuc — TUYET DOI KHONG
// tu tra k.quai, KHONG tu tinh sat thuong o day. Toan bo kiem tra + tinh toan
// + phat goi nam TRONG vongLap (xem case suKienDanhQuai, khuvuc.go).
func (b *BoXuLy) xuLyDanhQuai(m *protocol.Message) error {
	if b.khuVuc == nil || b.nhanVat == nil {
		return nil
	}
	r := protocol.NewReader(m.Data)
	mobIdRaw, err := r.Byte()
	if err != nil {
		return nil // goi rong bat thuong -> bo qua an toan, khong ngat ket noi
	}
	if int8(mobIdRaw) == -1 {
		// Thu cung (dac ta §0/§3.1): doc THEM 4 byte de tieu thu DU than goi
		// 5-byte. Gia tri masterId khong dung toi o KH6 (chua co he thu
		// cung) nen co tinh KHONG gan vao bien nao — chi doc de tranh lech.
		if _, err := r.Int(); err != nil {
			return nil // bao mobId=-1 nhung thieu du 4 byte con lai -> bo qua an toan
		}
	}
	b.khuVuc.Gui(suKienDanhQuai{idNhanVat: b.nhanVat.ID, mobId: mobIdRaw})
	return nil
}

// xuLyDanhNguoiChoi xu ly goi LEN `-60` (dac ta §3.1, Ke hoach 6 Task 11, Nhom
// F "nguoi danh nguoi"): client bao muon danh nguoi choi co id = doc duoc.
// Than goi CHI CO MOT truong `int idNguoiBiDanh` (dac ta §3.1) — KHONG co do
// dai bien thien nhu goi 54 (danh quai), nen ham nay don gian hon xuLyDanhQuai.
//
// Ky luat actor (dac ta §1.6 + brief Task 11): ham nay CHI doc idNguoiBiDanh
// tu goi roi gui suKienDanhNguoiChoi vao hop thu actor khu vuc — TUYET DOI
// KHONG tu tra k.nguoiTrong, KHONG tu tinh sat thuong o day. Toan bo kiem tra
// + tinh toan + phat goi nam TRONG vongLap (xem case suKienDanhNguoiChoi,
// khuvuc.go).
func (b *BoXuLy) xuLyDanhNguoiChoi(m *protocol.Message) error {
	if b.khuVuc == nil || b.nhanVat == nil {
		return nil
	}
	r := protocol.NewReader(m.Data)
	idMucTieu, err := r.Int()
	if err != nil {
		return nil // goi ngan hon khuon -> bo qua an toan, khong ngat ket noi
	}
	b.khuVuc.Gui(suKienDanhNguoiChoi{idNhanVat: b.nhanVat.ID, idMucTieu: int64(idMucTieu)})
	return nil
}

// xuLyHoiSinh xu ly goi LEN -16 (dac ta §3.1, Ke hoach 6 Task 13): client xin
// hoi sinh sau khi chet. Than RONG — khong co gi de doc, khac han xuLyDanhQuai
// (goi 54 co than). Luong phien CHI gui su kien vao hop thu actor (ky luat
// actor §1.6 + brief Task 13) — actor moi la noi quyet dinh co dang THAT SU
// chet hay khong + hoi mau/chakra day du + dat lucHoiSinh + phat goi
// -16/-30-15 (xem case suKienHoiSinh, khuvuc.go).
func (b *BoXuLy) xuLyHoiSinh() error {
	if b.khuVuc == nil || b.nhanVat == nil {
		return nil
	}
	b.khuVuc.Gui(suKienHoiSinh{idNhanVat: b.nhanVat.ID})
	return nil
}

// KhiDong duoc PhucVu goi khi phien dong (defer). Neu nhan vat da vao khu vuc,
// bao actor go no ra (actor phat -6 cho nguoi con lai). An toan goi khi khuVuc nil.
// Chi gui suKienRa neu da thuc su vao (daVaoKhuVuc) — chua vao thi khong can go,
// du actor xu ly suKienRa cua id khong co van la no-op an toan.
//
// Ke hoach 11 Task 5 (dac ta §4 diem B): TRUOC khi bao actor go nguoi choi ra
// (Gui(suKienRa)), xin "anh chup" trang thai nguoi choi (AnhChupNguoiChoi,
// khuvuc.go) roi Xep() vao tho ghi de luu xuong dia. THU TU NAY BAT BUOC,
// TUYET DOI KHONG DUOC DAO NGUOC:
//
//	DUNG:  AnhChupNguoiChoiCoLoi(id) ──► Xep(ban) ──► RoiKhuVuc(id)
//	SAI :  RoiKhuVuc(id) ──► AnhChupNguoiChoiCoLoi(id) tra VE RONG
//
// Ly do: case suKienRa (khuvuc.go) XOA nguoi choi khoi k.nguoiTrong NGAY khi
// actor xu ly su kien do. Neu goi Gui(suKienRa) truoc, den luc
// AnhChupNguoiChoi duoc goi thi actor rat co the (thuong la se) DA xu ly xong
// suKienRa va xoa nguoi nay khoi map — AnhChupNguoiChoi(id) khong con tim
// thay ai, tra ve slice RONG, nghia la MAT SACH tien do vua choi ngay truoc
// luc ngat ket noi. Nguoi sua sau TUYET DOI giu dung thu tu nay.
//
// AnhChupNguoiChoi da tu co han gio rieng (100ms, khuvuc.go) va Xep khong bao
// gio chan (thogi.go, Ke hoach 11 Task 3) nen ca hai buoc o day an toan de
// goi tren luong phien luc dong ket noi — KHONG duoc them buoc cho nao khac
// vao ham nay (dac ta §3: actor/luong phien khong bao gio duoc cho o dia).
//
// Ke hoach 11 Task 7 (dac ta §7.1 — chan dung bay "doc truoc, da sau" cua ban
// Java cu): NEU phien nay VUA bi DA boi mot dang nhap moi cung tai khoan
// (b.biDa da dong — xem daBiDa()), buoc luu o day PHAI la GHI DONG BO XAC
// NHAN (XepVaCho, cho ket qua that su tu CSDL) THAY VI Xep() khong chan binh
// thuong — vi SoDangNhap.DangKy (goroutine cua phien MOI) dang CHO tin hieu
// b.daXongDon o day de biet luc nao AN TOAN doc lai du lieu. Neu chi Xep()
// (dua vao hang doi trong bo nho roi bao "xong" ngay), phien MOI co the doc
// CSDL TRUOC KHI ban ghi that su xuong dia — dung dung loi (1) cua ban Java
// cu (NDVSqlFetcher.login(): doc CSDL TRUOC roi moi da phien cu) ma Task 7
// phai chan.
//
// Voi phien dut ket noi TU NHIEN (KHONG bi da): giu NGUYEN hanh vi CU (Task
// 5) — Xep() khong chan. Khong co ai dang cho tin hieu b.daXongDon de doc lai
// CSDL ngay lap tuc trong truong hop nay, nen khong can tra gia dong bo; giu
// dung tinh than "khong cho o dia lau hon muc can thiet" cho duong pho bien
// nhat (nguoi choi thoat game binh thuong).
//
// baoXongDon() LUON chay (defer, ke ca khi som return vi chua vao khu vuc) —
// SoDangNhap.DangKy phai luon nhan duoc tin hieu "xong", bat ke phien bi da
// co gi de luu hay khong, neu khong no se cho toi het han gio moi thoi.
//
// XoaNeuDung o dau ham (truoc ca kiem tra khuVuc/nhanVat/daVaoKhuVuc): mot
// phien dang ky vao so NGAY khi dang nhap thanh cong (xuLyDangNhap, TRUOC ca
// khi co nhan vat) — mot phien dut ket noi TRUOC khi vao the gioi (vd dang o
// man chon nhan vat) van phai duoc go khoi so, khong thi so giu "rac" mai.
func (b *BoXuLy) KhiDong() {
	defer b.baoXongDon()

	if b.taiKhoan != nil && b.soDangNhap != nil {
		b.soDangNhap.XoaNeuDung(b.taiKhoan.ID, b)
	}

	if b.khuVuc == nil || b.nhanVat == nil || !b.daVaoKhuVuc {
		return // chua vao khu vuc -> khong co gi de luu, khong can bao actor
	}

	if b.thoGhi != nil {
		ban, err := b.khuVuc.AnhChupNguoiChoiCoLoi(b.nhanVat.ID)
		if err != nil {
			slog.Warn("khong chup duoc trang thai khi dong phien",
				"idNhanVat", b.nhanVat.ID, "loi", err)
		}
		if len(ban) > 0 {
			if b.daBiDa() {
				if gdb, ok := b.thoGhi.(ghiDongBoKhiBiDa); ok {
					ctx, huy := context.WithTimeout(context.Background(), hanGioGhi)
					if err := gdb.XepVaCho(ctx, ban[0]); err != nil {
						slog.Error("ghi dong bo phien BI DA that bai — co the mat vai thay doi GAN NHAT cua phien bi da",
							"idNhanVat", ban[0].ID, "loi", err)
					}
					huy()
				} else {
					// Khong nen xay ra o may chu that (thoGhi that luon la
					// *ThoGhi, thoa man ghiDongBoKhiBiDa) — chi cham toi khi
					// mot test tu gan mot fake CHI co Xep(). Rot ve Xep() binh
					// thuong thay vi bo qua hoan toan buoc luu.
					b.thoGhi.Xep(ban[0])
				}
			} else {
				b.thoGhi.Xep(ban[0])
			}
		}
	}

	if err := b.khuVuc.RoiKhuVuc(b.nhanVat.ID); err != nil {
		slog.Warn("khong the go nguoi choi khoi actor khi dong phien",
			"idNhanVat", b.nhanVat.ID, "loi", err)
	}
}

// daBiDa bao phien nay da bi da (b.biDa da dong) hay chua — AN TOAN goi tren
// kenh nil (select+default khong bao gio chan tren kenh nil, chi khong bao
// gio "san sang", nen luon roi vao default = false).
func (b *BoXuLy) daBiDa() bool {
	select {
	case <-b.biDa:
		return true
	default:
		return false
	}
}

// baoXongDon dong b.daXongDon — tin hieu "KhiDong da chay xong hoan toan" cho
// SoDangNhap.DangKy dang cho tren goroutine cua phien khac. An toan goi dung
// MOT LAN (KhiDong chi chay dung mot lan qua defer trong mynet.PhucVu).
func (b *BoXuLy) baoXongDon() {
	if b.daXongDon != nil {
		close(b.daXongDon)
	}
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

	// Soi toan nhanh KH11 (fix-cuoi-report.md, Finding 1 — CRITICAL): chot
	// chan dang nhap CHONG tren CUNG mot ket noi. Bat bien "mot BoXuLy <->
	// mot tai khoan" ma SoDangNhap.DangKy dua vao (sodangnhap.go) chi con
	// dung neu MOI ket noi dang ky vao so DUNG MOT LAN. Neu khong chan:
	// mot -29 THANH CONG lan hai tren CUNG ket noi se gan lai b.taiKhoan
	// roi DangKy(tai_khoan_MOI, b) — b GIO NAM O HAI KHOA trong so cung
	// luc (khoa CU van con tro ve chinh no). Khi sau nay co ai da tai
	// khoan MOI, KhiDong() cua b chi XoaNeuDung duoc khoa MOI (so dang tro
	// dung no) — khoa CU tro thanh "rac" vinh vien, va lan da TIEP THEO
	// nham dung tai khoan CU se goi close() len CHINH b.biDa DA DONG TU
	// TRUOC (o lan da cua tai khoan moi) -> double-close mot channel ->
	// PANIC khong recover duoc (net/session.go khong co recover() tren
	// duong PhucVu) -> CHET CA TIEN TRINH may chu (DoS tu xa, tai khoan tu
	// dang ky duoc). Dang nhap SAI (KiemTraDangNhap tra ThanhCong=false)
	// KHONG gan b.taiKhoan, nen nguoi choi go sai mat khau van thu lai
	// binh thuong tren cung ket noi — nhanh nay CHI chan dang nhap THANH
	// CONG lan thu hai tro di, khong chan thu lai sau khi go sai.
	if b.taiKhoan != nil {
		return nil
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

	// Ke hoach 11 Task 7 (dac ta §7.1/§7.2): luu Session cua CHINH phien nay
	// TRUOC khi dang ky vao so (de mot dang nhap thu ba, gan nhu ngay sau,
	// van co the da duoc phien nay), ROI GOI DangKy — thu tu BAT BUOC:
	//
	//	DUNG: DangKy (da phien cu + CHO no luu xong) ──► LayNhanVatTheoTaiKhoan
	//	SAI : LayNhanVatTheoTaiKhoan ──► DangKy (da phien cu)
	//
	// DangKy (neu co phien cu cung tai khoan) DA ket noi cu VA CHO no bao "da
	// luu xong" (thong qua b.daXongDon cua phien cu) TRUOC KHI ham nay tra ve
	// — dam bao dong nay CHAC CHAN doc duoc du lieu MOI NHAT tu CSDL, khong
	// phai ban CU HON con nam lai tu truoc khi phien cu kip ghi. Day CHINH LA
	// thu tu dao nguoc voi loi (1) cua ban Java cu (NDVSqlFetcher.login():
	// doc CSDL TRUOC roi moi da phien cu o dong sau) — dao nguoc lai dong nay
	// se lam TestDangNhapTrung_DaPhienCuChoLuuXongRoiDocViTriMoiNhat do (xem
	// bao cao Task 7, muc dot bien).
	b.session = s
	if b.soDangNhap != nil {
		b.soDangNhap.DangKy(kq.TaiKhoan.ID, b)
	}

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
