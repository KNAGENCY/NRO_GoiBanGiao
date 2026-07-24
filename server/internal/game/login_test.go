package game

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"strings"
	"testing"
	"time"

	"nroserver/internal/auth"
	mynet "nroserver/internal/net"
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// khoGia gia lap Kho. Cac field tai khoan (tk/err) giu nguyen tu Ke hoach 1 de
// khong phai sua lai moi test cu chi dat hai field do. Ba field con lai la moi
// cua Ke hoach 2, mac dinh nil/false nen test cu khong dung toi van bien dich
// va chay binh thuong.
type khoGia struct {
	tk  *store.TaiKhoan
	err error

	nhanVat       *store.NhanVat // nil = tai khoan chua co nhan vat
	loiLayNhanVat error          // dat = store.ErrKhongThayNhanVat de gia lap "chua co"
	tenDaDung     bool
	nhanVatVuaTao *store.NhanVat
	loiTaoNhanVat error
}

func (k *khoGia) LayTaiKhoan(ctx context.Context, ten string) (*store.TaiKhoan, error) {
	if k.err != nil {
		return nil, k.err
	}
	return k.tk, nil
}

func (k *khoGia) LayNhanVatTheoTaiKhoan(ctx context.Context, taiKhoanID int64) (*store.NhanVat, error) {
	if k.loiLayNhanVat != nil {
		return nil, k.loiLayNhanVat
	}
	return k.nhanVat, nil
}

func (k *khoGia) TenNhanVatDaDung(ctx context.Context, ten string) (bool, error) {
	return k.tenDaDung, nil
}

func (k *khoGia) TaoNhanVat(ctx context.Context, nv store.NhanVatMoi) (*store.NhanVat, error) {
	if k.loiTaoNhanVat != nil {
		return nil, k.loiTaoNhanVat
	}
	return k.nhanVatVuaTao, nil
}

// taiKhoanHopLeChoTest tra mot tai khoan hop le voi mat khau bam san khop
// "matkhau", dung chung cho cac test dang nhap re nhanh nhan vat duoi day.
func taiKhoanHopLeChoTest() *store.TaiKhoan {
	bam, _ := auth.Hash("matkhau")
	return &store.TaiKhoan{ID: 1, TenDangNhap: "nguoidung", MatKhauBam: bam}
}

// dungGoiDangNhap dung goi dang nhap dung y het client gui:
// khung lenh -29, ben trong la byte lenh con 0 roi hai chuoi.
func dungGoiDangNhap(ten, matKhau string) *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(0) // lenh con: dang nhap
	w.UTF(ten)
	w.UTF(matKhau)
	return protocol.NewMessage(-29, w)
}

// Goi dung dinh dang phai duoc boc tach ra dung ten va mat khau.
func TestBocTachGoiDangNhap(t *testing.T) {
	m := dungGoiDangNhap("nguoi", "matkhau")
	r := protocol.NewReader(m.Data)

	lenhCon, err := r.Byte()
	if err != nil || lenhCon != 0 {
		t.Fatalf("lenh con = %d, err = %v, mong doi 0", lenhCon, err)
	}
	ten, _ := r.UTF()
	mk, _ := r.UTF()
	if ten != "nguoi" || mk != "matkhau" {
		t.Fatalf("boc tach ra %q / %q, mong doi \"nguoi\" / \"matkhau\"", ten, mk)
	}
}

func TestDangNhapDungMatKhau(t *testing.T) {
	bam, _ := auth.Hash("matkhau")
	b := NewBoXuLy(&khoGia{tk: &store.TaiKhoan{ID: 7, TenDangNhap: "nguoi", MatKhauBam: bam}})

	kq, err := b.KiemTraDangNhap(context.Background(), "nguoi", "matkhau")
	if err != nil {
		t.Fatalf("loi khong mong doi: %v", err)
	}
	if !kq.ThanhCong {
		t.Fatalf("phai dang nhap duoc, ly do tu choi: %q", kq.LyDo)
	}
	if kq.TaiKhoan.ID != 7 {
		t.Fatalf("ID = %d, mong doi 7", kq.TaiKhoan.ID)
	}
}

func TestDangNhapSaiMatKhau(t *testing.T) {
	bam, _ := auth.Hash("matkhau")
	b := NewBoXuLy(&khoGia{tk: &store.TaiKhoan{ID: 7, MatKhauBam: bam}})

	kq, _ := b.KiemTraDangNhap(context.Background(), "nguoi", "sai")
	if kq.ThanhCong {
		t.Fatal("mat khau sai ma van cho vao")
	}
	if kq.LyDo != LyDoSaiThongTin {
		t.Fatalf("ly do = %q, mong doi %q", kq.LyDo, LyDoSaiThongTin)
	}
}

// Tai khoan khong ton tai va sai mat khau phai tra ve CUNG MOT cau bao,
// de ke xau khong do duoc tai khoan nao co that.
func TestTaiKhoanKhongTonTaiBaoGiongSaiMatKhau(t *testing.T) {
	b := NewBoXuLy(&khoGia{err: store.ErrKhongThayTaiKhoan})

	kq, _ := b.KiemTraDangNhap(context.Background(), "khong-co", "gi-do")
	if kq.ThanhCong {
		t.Fatal("tai khoan khong ton tai ma van cho vao")
	}
	if kq.LyDo != LyDoSaiThongTin {
		t.Fatalf("ly do = %q, phai giong het truong hop sai mat khau", kq.LyDo)
	}
}

func TestTaiKhoanBiKhoa(t *testing.T) {
	bam, _ := auth.Hash("matkhau")
	b := NewBoXuLy(&khoGia{tk: &store.TaiKhoan{ID: 7, MatKhauBam: bam, BiKhoa: true}})

	kq, _ := b.KiemTraDangNhap(context.Background(), "nguoi", "matkhau")
	if kq.ThanhCong {
		t.Fatal("tai khoan bi khoa ma van cho vao")
	}
	if kq.LyDo != LyDoBiKhoa {
		t.Fatalf("ly do = %q, mong doi %q", kq.LyDo, LyDoBiKhoa)
	}
}

// Ten dang nhap rong khong duoc cham toi CSDL.
func TestTenRongBiTuChoiNgay(t *testing.T) {
	b := NewBoXuLy(&khoGia{err: context.Canceled}) // kho se loi neu bi goi
	kq, err := b.KiemTraDangNhap(context.Background(), "", "gi-do")
	if err != nil {
		t.Fatalf("khong duoc goi toi kho: %v", err)
	}
	if kq.ThanhCong {
		t.Fatal("ten rong ma van cho vao")
	}
}

// ---------------------------------------------------------------------------
// Kiem thu dau-cuoi (end-to-end) qua XuLy that.
//
// Cac test o tren chi goi KiemTraDangNhap() — ham logic thuan tuy, khong he
// dung toi mang. Kiem thu dot bien (mutation testing) da chi ra dieu nay la
// mot lo hong nghiem trong: xoa han buoc gui goi -102 trong guiTuChoi(), hoac
// dao nguoc thu tu kiem khoa/mat khau trong KiemTraDangNhap(), CA HAI deu lam
// toan bo 6 test cu tren van DAT (vi khong test nao cham toi XuLy() hay
// Session.Gui() that su). Cac test duoi day dung net.Pipe() de dung mot phien
// that hoan chinh (bat tay -> ma hoa -> dang nhap) va bat dung cac goi tin
// server GUI XUONG DAY tren day, giong het cach client that nhan.
// ---------------------------------------------------------------------------

// docGoiKhoaThuCong doc goi khoa ma server tra ve trong buoc bat tay, dung
// CHINH XAC cach client that lam (xem net/session_test.go,
// TestBatTayRoiNhanGoiDaMaHoa): do dai va noi dung goi khoa duoc doc THU CONG
// theo big-endian, TUYET DOI KHONG dung Framer.ReadMessage() o day — truoc khi
// co khoa, hai chieu client/server dung hai thu tu byte KHAC NHAU (xem chu
// thich cua kieu Framer trong protocol/frame.go).
func docGoiKhoaThuCong(t *testing.T, conn net.Conn) []byte {
	t.Helper()

	var cmdByte [1]byte
	if _, err := io.ReadFull(conn, cmdByte[:]); err != nil {
		t.Fatalf("khong doc duoc lenh goi khoa: %v", err)
	}
	if cmdTraLoi := int8(cmdByte[0]); cmdTraLoi != -27 {
		t.Fatalf("lenh goi khoa = %d, mong doi -27", cmdTraLoi)
	}

	var doDaiBytes [2]byte
	if _, err := io.ReadFull(conn, doDaiBytes[:]); err != nil {
		t.Fatalf("khong doc duoc do dai goi khoa: %v", err)
	}
	soByteThan := int(doDaiBytes[0])<<8 | int(doDaiBytes[1])

	than := make([]byte, soByteThan)
	if _, err := io.ReadFull(conn, than); err != nil {
		t.Fatalf("khong doc duoc than tin goi khoa (%d byte): %v", soByteThan, err)
	}

	r := protocol.NewReader(than)
	doDai, err := r.Byte()
	if err != nil {
		t.Fatalf("doc do dai khoa that bai: %v", err)
	}
	if doDai < 1 || doDai > 127 {
		t.Fatalf("do dai khoa = %d, phai trong 1..127", doDai)
	}
	khoaTho := make([]byte, doDai)
	for i := range khoaTho {
		khoaTho[i], _ = r.Byte()
	}
	// Ba truong con lai (dia chi/cong/co server phu) khong can cho test nay,
	// nhung van doc het de dung quy uoc giong TestBatTayRoiNhanGoiDaMaHoa.
	_, _ = r.UTF()
	_, _ = r.Int()
	_, _ = r.Byte()

	return khoaTho
}

// dungPhienThat dung mot phien XuLy that qua net.Pipe(): chay mynet.PhucVu
// tren mot goroutine, thuc hien bat tay va tra ve mot Framer da co khoa ma
// phia "client gia" trong test co the dung de gui/nhan goi tin.
//
// Luu y: sau khi co khoa, hai chieu client/server DOI XUNG (xem chu thich cua
// Framer trong protocol/frame.go), nen dung LAI chinh Framer nay de vua GUI
// (gia lam client gui len server) vua DOC (nhan goi server tra ve) la hop le.
func dungPhienThat(t *testing.T, b *BoXuLy) (fClient *protocol.Framer, dong func()) {
	t.Helper()

	benClient, benServer := net.Pipe()
	go mynet.PhucVu(benServer, func() mynet.XuLyGoiTin { return b })

	benClient.SetReadDeadline(time.Now().Add(5 * time.Second))

	fClient = protocol.NewFramer(benClient, benClient)
	if err := fClient.WriteMessage(&protocol.Message{Command: -27}); err != nil {
		t.Fatalf("gui goi bat tay that bai: %v", err)
	}
	khoaTho := docGoiKhoaThuCong(t, benClient)
	fClient.SetKey(protocol.DeriveKey(khoaTho))

	return fClient, func() { benClient.Close() }
}

// chayDangNhap dung mot BoXuLy tren kho da cho, gui goi dang nhap that qua
// net.Pipe (dungPhienThat), roi doc HET cac goi may chu tra ve cho toi khi
// ket noi RANH — tuc la ReadMessage het han gio (xem han 5 giay dat trong
// dungPhienThat) — va tra ve DANH SACH cac goi da doc duoc.
//
// Vi sao doc-toi-khi-ranh thay vi doc dung mot so goi biet truoc: dung cho ca
// hai nhanh Ke hoach 2 — nhanh "da co nhan vat" (stub vaoTheGioi khong gui goi
// nao ca, danh sach ket qua se RONG) va nhanh "chua co nhan vat" (gui dung
// mot goi lenh 2). Neu doc theo so luong co dinh se khong the biet chac la
// "khong con goi nao nua" hay chi la "chua toi luot goi do".
func chayDangNhap(t *testing.T, kho Kho, ten, matKhau string) []*protocol.Message {
	t.Helper()

	b := NewBoXuLy(kho)
	fClient, dong := dungPhienThat(t, b)
	defer dong()

	if err := fClient.WriteMessage(dungGoiDangNhap(ten, matKhau)); err != nil {
		t.Fatalf("gui goi dang nhap that bai: %v", err)
	}

	var goiGuiRa []*protocol.Message
	for {
		m, err := fClient.ReadMessage()
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				break // het han gio ma khong co goi moi -> ket noi da ranh
			}
			t.Fatalf("loi khi doc goi may chu tra ve: %v", err)
		}
		goiGuiRa = append(goiGuiRa, m)
	}
	return goiGuiRa
}

// coLenh bao mot danh sach goi (tra ve tu chayDangNhap) co chua goi mang lenh
// da cho hay khong.
func coLenh(goi []*protocol.Message, lenh int8) bool {
	for _, m := range goi {
		if m.Command == lenh {
			return true
		}
	}
	return false
}

// TestTuChoiDangNhapGuiDuHaiGoiDungThuTu la kiem thu QUAN TRONG NHAT cua Task 8.
//
// Boi canh: ban Java cu TUNG co loi that trong guiTuChoi (Service.sendLoginFail):
// quen mat goi thu hai (-102, bao client thoat trang thai "dang cho"). Hau qua
// la khi dang nhap sai, client CHI nhan duoc hop thoai bao loi (-26) roi treo
// VINH VIEN o man hinh "Xin cho" — nguoi choi khong lam gi duoc ngoai tat ung
// dung. Loi nay da duoc sua trong guiTuChoi() hien tai (goi ca -26 lan -102),
// nhung KHONG CO TEST NAO bao ve dieu do: moi test truoc day chi goi
// KiemTraDangNhap() — ham logic thuan tuy, khong di qua XuLy()/guiTuChoi() nen
// khong the phat hien neu ai do sau nay lo xoa mat goi thu hai.
//
// Test nay dung mot phien that (net.Pipe + mynet.PhucVu that), gui goi dang
// nhap SAI mat khau, roi bat dung HAI goi tin server tra ve tren day va
// khang dinh: (1) co du hai goi, (2) dung lenh, (3) dung noi dung, (4) dung
// THU TU (-26 truoc, -102 sau — dung thu tu client can: HIEN HOP THOAI TRUOC
// ROI MOI THOAT TRANG THAI CHO, khong thi hop thoai hien len nhung client van
// ket qua o trang thai cho — dung loi that da tung xay ra o ban Java cu).
func TestTuChoiDangNhapGuiDuHaiGoiDungThuTu(t *testing.T) {
	bam, _ := auth.Hash("dungmatkhau")
	b := NewBoXuLy(&khoGia{tk: &store.TaiKhoan{ID: 1, TenDangNhap: "nguoi", MatKhauBam: bam}})

	fClient, dong := dungPhienThat(t, b)
	defer dong()

	// Gui goi dang nhap SAI mat khau, dung dinh dang client that gui (khung
	// -29, than: 1 byte lenh con = 0, roi ten va mat khau dang UTF).
	if err := fClient.WriteMessage(dungGoiDangNhap("nguoi", "sai-mat-khau")); err != nil {
		t.Fatalf("gui goi dang nhap that bai: %v", err)
	}

	// Goi THU NHAT phai la -26 (hop thoai OK) chua dung cau LyDoSaiThongTin.
	m1, err := fClient.ReadMessage()
	if err != nil {
		t.Fatalf("khong doc duoc goi thu nhat (-26): %v", err)
	}
	if m1.Command != -26 {
		t.Fatalf("goi thu nhat co lenh %d, mong doi -26 (hop thoai OK)", m1.Command)
	}
	r1 := protocol.NewReader(m1.Data)
	cau, err := r1.UTF()
	if err != nil {
		t.Fatalf("khong doc duoc cau bao trong goi -26: %v", err)
	}
	if cau != LyDoSaiThongTin {
		t.Fatalf("cau bao = %q, mong doi %q", cau, LyDoSaiThongTin)
	}

	// Goi THU HAI phai la -102 (thoat trang thai "dang cho"), than dung 1
	// byte gia tri 0. DAY LA PHEP KIEM CHINH: neu ai xoa mat buoc gui goi nay
	// trong guiTuChoi(), ReadMessage() o day se KHONG CO GI DE DOC va test
	// phai FAIL (nho SetReadDeadline o dungPhienThat, test se timeout thay vi
	// treo vinh vien).
	m2, err := fClient.ReadMessage()
	if err != nil {
		t.Fatalf("khong doc duoc goi thu hai (-102): %v — thieu goi nay la lap lai dung loi da tung xay ra o ban Java (client treo man hinh \"Xin cho\")", err)
	}
	if m2.Command != -102 {
		t.Fatalf("goi thu hai co lenh %d, mong doi -102 (thoat trang thai dang cho)", m2.Command)
	}
	if !bytes.Equal(m2.Data, []byte{0}) {
		t.Fatalf("than goi -102 = %v, mong doi dung 1 byte [0]", m2.Data)
	}
}

// TestDangNhapDungGuiLenhTaoNhanVat la duong thanh cong doi ung voi test tren:
// mat khau DUNG phai lam server gui xuong lenh 2 (chuyen man tao nhan vat),
// than rong. Cung dung phien that (khong chi goi KiemTraDangNhap()) de bat
// loi o ca buoc gui goi tin (guiChuyenSangTaoNhanVat), khong chi o logic kiem tra.
//
// Ke hoach 2: dang nhap xong CHI gui lenh 2 khi tai khoan CHUA co nhan vat.
// khoGia o day phai gia lap "chua co" (loiLayNhanVat: ErrKhongThayNhanVat) de
// gia thiet cu cua test nay (Ke hoach 1) van dung — xem
// TestDangNhapDaCoNhanVat_KhongGuiManTao ben duoi cho nhanh nguoc lai.
func TestDangNhapDungGuiLenhTaoNhanVat(t *testing.T) {
	bam, _ := auth.Hash("dungmatkhau")
	b := NewBoXuLy(&khoGia{
		tk:            &store.TaiKhoan{ID: 1, TenDangNhap: "nguoi", MatKhauBam: bam},
		loiLayNhanVat: store.ErrKhongThayNhanVat,
	})

	fClient, dong := dungPhienThat(t, b)
	defer dong()

	if err := fClient.WriteMessage(dungGoiDangNhap("nguoi", "dungmatkhau")); err != nil {
		t.Fatalf("gui goi dang nhap that bai: %v", err)
	}

	m, err := fClient.ReadMessage()
	if err != nil {
		t.Fatalf("khong doc duoc goi tra ve: %v", err)
	}
	if m.Command != 2 {
		t.Fatalf("lenh = %d, mong doi 2 (chuyen sang man tao nhan vat)", m.Command)
	}
	if len(m.Data) != 0 {
		t.Fatalf("than goi lenh 2 = %v (dai %d byte), mong doi rong", m.Data, len(m.Data))
	}
}

// TestDangNhapChuaCoNhanVat_GuiManTao va TestDangNhapDaCoNhanVat_KhongGuiManTao
// la hai kiem thu chinh cua Task 5 (Ke hoach 2): dang nhap thanh cong phai RE
// NHANH theo viec tai khoan da co nhan vat hay chua — thay vi luon gui lenh 2
// nhu Ke hoach 1.
func TestDangNhapChuaCoNhanVat_GuiManTao(t *testing.T) {
	kho := &khoGia{
		tk:            taiKhoanHopLeChoTest(),
		loiLayNhanVat: store.ErrKhongThayNhanVat,
	}
	goiGuiRa := chayDangNhap(t, kho, "nguoidung", "matkhau")
	if !coLenh(goiGuiRa, 2) {
		t.Error("chua co nhan vat -> phai gui lenh 2 (man tao nhan vat)")
	}
}

func TestDangNhapDaCoNhanVat_KhongGuiManTao(t *testing.T) {
	kho := &khoGia{
		tk:      taiKhoanHopLeChoTest(),
		nhanVat: &store.NhanVat{ID: 7, Ten: "Naruto", Lang: 0},
	}
	goiGuiRa := chayDangNhap(t, kho, "nguoidung", "matkhau")
	if coLenh(goiGuiRa, 2) {
		t.Error("da co nhan vat -> KHONG duoc gui lenh 2 (khong quay lai man tao)")
	}
}

// TestSaiMatKhauTrenTaiKhoanBiKhoaVanBaoSaiThongTin khoa cung THU TU kiem tra
// ben trong KiemTraDangNhap: MAT KHAU phai duoc kiem TRUOC, BI KHOA kiem SAU.
//
// Tinh huong: mot tai khoan bi khoa (BiKhoa: true) nhung nguoi go SAI mat
// khau. Neu code dao nguoc thu tu (kiem bi khoa TRUOC mat khau), ke go bay se
// nhan duoc LyDoBiKhoa — vo tinh XAC NHAN voi ho rang tai khoan do LA CO THAT
// (chi dang bi khoa), pha vo dung nguyen tac "mot cau bao chung" ma
// TestTaiKhoanKhongTonTaiBaoGiongSaiMatKhau o tren da dat ra: khong duoc de
// cau bao lo tai khoan nao co that. Kiem thu dot bien da chi ra day la mot
// nhanh khong test nao trong file nay tung cham toi (cac test bi khoa/sai mat
// khau truoc gio deu dung rieng, khong ket hop ca hai dieu kien cung luc).
func TestSaiMatKhauTrenTaiKhoanBiKhoaVanBaoSaiThongTin(t *testing.T) {
	bam, _ := auth.Hash("matkhau-dung")
	b := NewBoXuLy(&khoGia{tk: &store.TaiKhoan{ID: 7, MatKhauBam: bam, BiKhoa: true}})

	kq, _ := b.KiemTraDangNhap(context.Background(), "nguoi", "mat-khau-sai")
	if kq.ThanhCong {
		t.Fatal("mat khau sai tren tai khoan bi khoa ma van cho vao")
	}
	if kq.LyDo != LyDoSaiThongTin {
		t.Fatalf("ly do = %q, mong doi %q (KHONG duoc la %q — se lo cho ke go bay biet tai khoan nay co that)",
			kq.LyDo, LyDoSaiThongTin, LyDoBiKhoa)
	}
}

// ---------------------------------------------------------------------------
// PHAN A ①: chan ngap nhat ky. Xem chu thich chi tiet o login.go va o
// main.go (mucNhatKyTuChuoi) de biet boi canh day du (9 GB/gio tu 1 ket noi).
// ---------------------------------------------------------------------------

// demNhatKy la mot slog.Handler chi dem so dong da ghi, dung de kiem tra
// "khong ghi nhat ky nao" ma khong can doc noi dung tung dong.
type demNhatKy struct {
	soDong int
}

func (d *demNhatKy) Enabled(context.Context, slog.Level) bool { return true }
func (d *demNhatKy) Handle(context.Context, slog.Record) error {
	d.soDong++
	return nil
}
func (d *demNhatKy) WithAttrs([]slog.Attr) slog.Handler { return d }
func (d *demNhatKy) WithGroup(string) slog.Handler      { return d }

// TestGoiLenhLaKhongGhiNhatKy la kiem chung cho viec ① trong PHAN A: nhan mot
// goi mang lenh la (khac -29), hoac mot goi -29 voi lenh con la (khac 0),
// TUYET DOI KHONG duoc ghi bat ky dong nhat ky nao — du o muc Debug.
//
// Boi canh: do bang thuc te tren may chu that, mot ket noi don gui goi 3 byte
// lenh la co the dat 32.452 goi/giay. Truoc day moi goi nhu vay lam ghi MOT
// dong Debug (~67 byte) trong XuLy() — khuech dai thanh 9 GB/GIO tren dia CHI
// TU MOT ket noi. Ha muc log mac dinh xuong info (viec o main.go) la chua du
// neu ai do dieu tra su co roi bat NRO_MUC_NHAT_KY=debug: dong log van con do
// va van khuech dai y het cu. Cach sua dung la BO HAN dong log, khong phai ha
// muc — client gui cac goi nay mot cach binh thuong (vi du lenh con
// setClientType), khong phai chuyen dang ghi lai.
func TestGoiLenhLaKhongGhiNhatKy(t *testing.T) {
	dem := &demNhatKy{}
	cu := slog.Default()
	slog.SetDefault(slog.New(dem))
	defer slog.SetDefault(cu)

	b := NewBoXuLy(&khoGia{err: context.Canceled}) // kho khong duoc goi toi

	// 1. Lenh top-level la (khac -29 = lenhChuaDangNhap).
	if err := b.XuLy(&mynet.Session{}, &protocol.Message{Command: 42}); err != nil {
		t.Fatalf("XuLy loi khong mong doi voi lenh la: %v", err)
	}

	// 2. Lenh con la, sau khi da vao dung khung -29 (vd client goi
	// setClientType, mot lenh con khac lenhConDangNhap = 0).
	w := protocol.NewWriter()
	w.Byte(9) // lenh con la
	if err := b.XuLy(&mynet.Session{}, protocol.NewMessage(-29, w)); err != nil {
		t.Fatalf("XuLy loi khong mong doi voi lenh con la: %v", err)
	}

	if dem.soDong != 0 {
		t.Fatalf("XuLy ghi %d dong nhat ky cho lenh la/lenh con la, mong doi 0 (phai bo han, khong chi ha muc)", dem.soDong)
	}
}

// ---------------------------------------------------------------------------
// PHAN A ②: kiem dinh dang ten dang nhap o tang Go, TRUOC khi cham toi kho.
// ---------------------------------------------------------------------------

// TestTenSaiDinhDangBiChanKhongGoiKho la kiem chung cho viec ② trong PHAN A:
// ten dang nhap khong khop ^[A-Za-z0-9_]{4,20}$ phai bi KiemTraDangNhap tu
// choi NGAY, khong duoc cham toi kho du lieu (truoc day chi CSDL kiem qua
// rang buoc CHECK, nen mot ten 60.000 ky tu van bi ghi nguyen vao nhat ky roi
// day xuong CSDL — cong voi ich vi chac chan khong khop gi).
//
// Dung khoGia bao loi neu bi goi de bat dung viec "khong goi toi kho".
func TestTenSaiDinhDangBiChanKhongGoiKho(t *testing.T) {
	b := NewBoXuLy(&khoGia{err: context.Canceled}) // kho se loi neu bi goi

	truongHop := []struct {
		mota string
		ten  string
	}{
		{"qua ngan (3 ky tu)", "abc"},
		{"qua dai (21 ky tu)", strings.Repeat("x", 21)},
		{"co dau tieng Viet", "nguơi_chơi"},
		{"ky tu dac biet", "nguoi!choi"},
	}

	for _, tc := range truongHop {
		t.Run(tc.mota, func(t *testing.T) {
			kq, err := b.KiemTraDangNhap(context.Background(), tc.ten, "mat-khau-gi-do")
			if err != nil {
				t.Fatalf("khong duoc goi toi kho: %v", err)
			}
			if kq.ThanhCong {
				t.Fatalf("ten %q sai dinh dang ma van cho dang nhap", tc.ten)
			}
			if kq.LyDo != LyDoSaiThongTin {
				t.Fatalf("ly do = %q, mong doi %q (khong duoc tao cau bao rieng — se lo luat dat ten cho ke xau)",
					kq.LyDo, LyDoSaiThongTin)
			}
		})
	}
}
