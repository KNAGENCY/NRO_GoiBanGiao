package game

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"strings"
	"sync"
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
	b := NewBoXuLy(&khoGia{tk: &store.TaiKhoan{ID: 7, TenDangNhap: "nguoi", MatKhauBam: bam}}, nil)

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
	b := NewBoXuLy(&khoGia{tk: &store.TaiKhoan{ID: 7, MatKhauBam: bam}}, nil)

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
	b := NewBoXuLy(&khoGia{err: store.ErrKhongThayTaiKhoan}, nil)

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
	b := NewBoXuLy(&khoGia{tk: &store.TaiKhoan{ID: 7, MatKhauBam: bam, BiKhoa: true}}, nil)

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
	b := NewBoXuLy(&khoGia{err: context.Canceled}, nil) // kho se loi neu bi goi
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

// docMotGoi doc DUNG MOT goi ma may chu gui ra qua fClient (framer client gia
// tra ve tu dungPhienThat), dung han gio 5 giay da dat san tren benClient.
// Dung cho cac test dieu phoi goi xin cap nhat (Ke hoach 3): khac chayDangNhap
// (doc TOI KHI ranh, tra ve ca danh sach) o cho ham nay chi can DUNG MOT goi
// duy nhat va bao loi ro rang (t.Fatal) neu khong doc duoc.
func docMotGoi(t *testing.T, f *protocol.Framer) *protocol.Message {
	t.Helper()
	m, err := f.ReadMessage()
	if err != nil {
		t.Fatalf("khong doc duoc goi may chu tra ve: %v", err)
	}
	return m
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

	b := NewBoXuLy(kho, nil)
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
	b := NewBoXuLy(&khoGia{tk: &store.TaiKhoan{ID: 1, TenDangNhap: "nguoi", MatKhauBam: bam}}, nil)

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
	}, nil)

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
	// Ke hoach 3: vaoTheGioi khong con la stub im lang nua — no gui DUNG MOT
	// goi -28/4 (versionGame) de kich client so phien ban. Truoc day (Ke hoach
	// 2) day la truong hop hoan toan im lang; gio phai co dung goi nay.
	if len(goiGuiRa) != 1 {
		t.Fatalf("so goi gui ra = %d, mong doi DUNG MOT (goi -28/4 versionGame): %v", len(goiGuiRa), goiGuiRa)
	}
	if goiGuiRa[0].Command != lenhKhongTrenBanDo {
		t.Fatalf("lenh = %d, mong doi -28 (versionGame)", goiGuiRa[0].Command)
	}
	r := protocol.NewReader(goiGuiRa[0].Data)
	sub, _ := r.Byte()
	if sub != lenhConVersionGame {
		t.Errorf("lenh con = %d, mong doi %d (versionGame)", sub, lenhConVersionGame)
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
	b := NewBoXuLy(&khoGia{tk: &store.TaiKhoan{ID: 7, MatKhauBam: bam, BiKhoa: true}}, nil)

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

	b := NewBoXuLy(&khoGia{err: context.Canceled}, nil) // kho khong duoc goi toi

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
	b := NewBoXuLy(&khoGia{err: context.Canceled}, nil) // kho se loi neu bi goi

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

// ---------------------------------------------------------------------------
// Ke hoach 3: dieu phoi goi xin cap nhat (-87 than rong, -28/6, -28/7, -28/8)
// sau khi da vao the gioi. Dat b.nhanVat truoc (mo phong da vao the gioi) —
// khong di qua toan bo duong dang nhap/tao nhan vat vi da co rieng cac test
// tren cho duong do; o day chi kiem TANG DIEU PHOI (XuLy/xuLyKhongTrenBanDo)
// co goi dung ham dung goi (Task 3-4) hay khong.
// ---------------------------------------------------------------------------

// khoaThuCong mo phong con tro cua protocol.KeyStream (key[cur] ^ b, con tro
// tang roi quay vong) NHUNG tu cai lai o day, vi KeyStream/Framer.ks khong
// xuat khoi package protocol nen goi game khong the muon lai con tro that. Chi
// dung cho guiGoiXinDataThuCong/docGoiLenhDacBietThuCong duoi day — xem hai
// ham do de biet tai sao can doc/ghi thu cong thay vi qua Framer.
type khoaThuCong struct {
	khoa []byte
	cur  int
}

func (k *khoaThuCong) xor(b byte) byte {
	out := k.khoa[k.cur] ^ b
	k.cur++
	if k.cur >= len(k.khoa) {
		k.cur %= len(k.khoa)
	}
	return out
}

// guiGoiXinDataThuCong ghi THU CONG goi yeu cau -87 than rong (client xin
// data) dung DUNG dinh dang HAI byte do dai ma client that gui (Session_ME.cs,
// doSendMessage — luon hai byte cho MOI goi client gui LEN, bat ke lenh gi).
//
// KHONG duoc dung fClient.WriteMessage() cho buoc nay: ham do (protocol/
// frame.go) quyet dinh hai hay BA byte do dai CHI dua vao gia tri lenh, khong
// biet gi ve chieu truyen — vi -87 nam trong danh sach lenh dac biet cua
// chieu may-chu-gui-xuong (WriteMessage se tu dong bat che do ba byte cho MOI
// goi mang lenh -87, ke ca than tin rong), dung ham do de gia lam client GUI
// LEN se tao ra mot goi SAI dinh dang thuc te. Hau qua da kiem chung: may chu
// (Framer.ReadMessage — CO Y CHI ho tro hai byte, vi client that khong bao
// gio gui lenh dac biet len may chu) doc lem byte thu ba cua truong do dai
// thanh mot phan than tin, roi cho doc them mai khong bao gio toi — treo han
// ca ket noi (i/o timeout tren net.Pipe).
func guiGoiXinDataThuCong(t *testing.T, conn net.Conn, khoa []byte) {
	t.Helper()
	k := &khoaThuCong{khoa: khoa}
	cmd := lenhDuLieu // bien (khong phai hang), de byte(cmd) la chuyen doi kieu binh thuong
	goi := []byte{
		k.xor(byte(cmd)),
		k.xor(0), k.xor(0), // do dai = 0 (than rong), dung HAI byte thuong
	}
	if _, err := conn.Write(goi); err != nil {
		t.Fatalf("gui goi xin data thu cong that bai: %v", err)
	}
}

// docGoiLenhDacBietThuCong doc THU CONG mot goi may chu gui xuong, dung DUNG
// thuat toan client that (Session_ME.cs, readMessage/readMessage2): can cu
// VAO LENH (khong can cu than tin) de quyet dinh doc hai hay BA byte do dai.
// O day chi biet phan biet lenhDuLieu (-87, package game co dung) — du de cho
// muc dich cua file nay, khong co tham vong lam lai toan bo danh sach lenh
// dac biet cua package protocol (lenhDoDaiBaByte, khong xuat).
//
// KHONG duoc dung fClient.ReadMessage()/docMotGoi cho goi -87: ham do (mo
// hinh chieu client->server) CO Y CHI ho tro hai byte do dai, trong khi may
// chu tra loi -87 bang BA byte (goiCatalogData ~87 KB, vuot gioi han hai byte
// 65535). Dung sai ham se lam chinh TEST nay doc lem, treo ket noi — day
// chinh la ly do can ham nay. Xem docGoiNhuClientThat trong
// protocol/integration_test.go cho ban dang bien the y het (khong tai su
// dung truc tiep duoc vi ham do khong xuat va khac package).
func docGoiLenhDacBietThuCong(t *testing.T, conn net.Conn, khoa []byte) *protocol.Message {
	t.Helper()
	k := &khoaThuCong{khoa: khoa}

	var cmdByte [1]byte
	if _, err := io.ReadFull(conn, cmdByte[:]); err != nil {
		t.Fatalf("khong doc duoc lenh: %v", err)
	}
	cmd := int8(k.xor(cmdByte[0]))

	var n int
	if cmd == lenhDuLieu {
		var doDai [3]byte
		for i := range doDai {
			var b [1]byte
			if _, err := io.ReadFull(conn, b[:]); err != nil {
				t.Fatalf("khong doc duoc byte do dai thu %d: %v", i, err)
			}
			doDai[i] = k.xor(b[0]) + 128
		}
		n = int(doDai[0]) | int(doDai[1])<<8 | int(doDai[2])<<16
	} else {
		var b [2]byte
		if _, err := io.ReadFull(conn, b[:]); err != nil {
			t.Fatalf("khong doc duoc do dai: %v", err)
		}
		hi := k.xor(b[0])
		lo := k.xor(b[1])
		n = int(hi)<<8 | int(lo)
	}

	data := make([]byte, n)
	if n > 0 {
		if _, err := io.ReadFull(conn, data); err != nil {
			t.Fatalf("khong doc duoc than tin %d byte: %v", n, err)
		}
		for i := range data {
			data[i] = k.xor(data[i])
		}
	}
	return &protocol.Message{Command: cmd, Data: data}
}

// TestPhucVuGoiXinData kiem tra dieu phoi -87 (xuLyXinData): client gui -87
// than rong, may chu phai tra loi bang goi -87 day du (goiCatalogData).
//
// Test nay KHONG dung dungPhienThat/docMotGoi nhu cac test dieu phoi khac
// trong file nay: -87 la lenh dac biet (do dai BA byte, vi goiCatalogData
// ~87 KB vuot qua gioi han HAI byte 65535 — xem lenhDoDaiBaByte trong
// protocol/frame.go), ma ca hai ham do deu CO Y CHI ho tro dinh dang hai byte
// (xem chu thich cua guiGoiXinDataThuCong/docGoiLenhDacBietThuCong o tren).
// Nen o day tu dung bat tay thu cong (docGoiKhoaThuCong, da dung o
// dungPhienThat) roi ghi/doc than tin bang hai ham thu cong nay.
func TestPhucVuGoiXinData(t *testing.T) {
	b := NewBoXuLy(&khoGia{}, nil)
	b.nhanVat = &store.NhanVat{ID: 1, Lang: 0}

	benClient, benServer := net.Pipe()
	go mynet.PhucVu(benServer, func() mynet.XuLyGoiTin { return b })
	defer benClient.Close()

	benClient.SetReadDeadline(time.Now().Add(5 * time.Second))

	fBatTay := protocol.NewFramer(benClient, benClient)
	if err := fBatTay.WriteMessage(&protocol.Message{Command: -27}); err != nil {
		t.Fatalf("gui goi bat tay that bai: %v", err)
	}
	khoaTho := docGoiKhoaThuCong(t, benClient)
	khoa := protocol.DeriveKey(khoaTho)

	// Client xin data: gui -87 than rong.
	guiGoiXinDataThuCong(t, benClient, khoa)

	m := docGoiLenhDacBietThuCong(t, benClient, khoa)
	if m.Command != -87 {
		t.Errorf("tra loi xin data = lenh %d, mong doi -87", m.Command)
	}
}

func TestPhucVuGoiXinSkill5Lang(t *testing.T) {
	b := NewBoXuLy(&khoGia{}, nil)
	b.nhanVat = &store.NhanVat{ID: 1, Lang: 0}
	fClient, dong := dungPhienThat(t, b)
	defer dong()
	// Client xin skill: -28 lenh con 7.
	w := protocol.NewWriter()
	w.Byte(7)
	if err := fClient.WriteMessage(protocol.NewMessage(-28, w)); err != nil {
		t.Fatal(err)
	}
	m := docMotGoi(t, fClient)
	r := protocol.NewReader(m.Data)
	sub, _ := r.Byte()
	if m.Command != -28 || sub != 7 {
		t.Fatalf("tra loi xin skill = lenh %d/%d, mong doi -28/7", m.Command, sub)
	}
	r.Byte() // vsSkill
	r.Byte() // countOption
	countClass, _ := r.Byte()
	if countClass < 5 {
		t.Errorf("catalog skill tra ve %d lop, phai >= 5 lang", countClass)
	}
}

// TestPhucVuGoiXinMap kiem tra dieu phoi lenh con 6 (client xin catalog map)
// trong xuLyKhongTrenBanDo. Cung mau voi TestPhucVuGoiXinSkill5Lang: -28
// khong nam trong lenhDoDaiBaByte nen phan hoi map dung dinh dang HAI byte do
// dai thuong, doc duoc bang dungPhienThat/docMotGoi binh thuong (khong can
// bat tay/doc thu cong nhu goi -87).
func TestPhucVuGoiXinMap(t *testing.T) {
	b := NewBoXuLy(&khoGia{}, nil)
	b.nhanVat = &store.NhanVat{ID: 1, Lang: 0}
	fClient, dong := dungPhienThat(t, b)
	defer dong()
	// Client xin map: -28 lenh con 6.
	w := protocol.NewWriter()
	w.Byte(6)
	if err := fClient.WriteMessage(protocol.NewMessage(-28, w)); err != nil {
		t.Fatal(err)
	}
	m := docMotGoi(t, fClient)
	r := protocol.NewReader(m.Data)
	sub, _ := r.Byte()
	if m.Command != -28 || sub != 6 {
		t.Fatalf("tra loi xin map = lenh %d/%d, mong doi -28/6", m.Command, sub)
	}
	vsMapNhan, _ := r.Byte()
	if vsMapNhan != vsMap {
		t.Errorf("vsMap tra ve = %d, mong doi %d (goiCatalogMap)", vsMapNhan, vsMap)
	}
}

// TestPhucVuGoiXinItem kiem tra dieu phoi lenh con 8 (client xin catalog
// item) trong xuLyKhongTrenBanDo. QUAN TRONG NHAT trong ba lenh con moi kiem
// o day: day la goi KICH clientOk phia client that (client chi gui -28/13 khi
// ca 4 phien ban da khop, VA da nhan du catalog item) — thieu dieu phoi lenh
// con nay se lam client ket vinh vien man hinh "Xin cho" (soft-lock), khong
// he co trieu chung loi ro rang phia may chu (khong log, khong crash).
func TestPhucVuGoiXinItem(t *testing.T) {
	b := NewBoXuLy(&khoGia{}, nil)
	b.nhanVat = &store.NhanVat{ID: 1, Lang: 0}
	fClient, dong := dungPhienThat(t, b)
	defer dong()
	// Client xin item: -28 lenh con 8.
	w := protocol.NewWriter()
	w.Byte(8)
	if err := fClient.WriteMessage(protocol.NewMessage(-28, w)); err != nil {
		t.Fatal(err)
	}
	m := docMotGoi(t, fClient)
	r := protocol.NewReader(m.Data)
	sub, _ := r.Byte()
	if m.Command != -28 || sub != 8 {
		t.Fatalf("tra loi xin item = lenh %d/%d, mong doi -28/8", m.Command, sub)
	}
}

// TestNhanClientOk kiem tra dieu phoi lenh con 13 (client bao san sang) trong
// xuLyKhongTrenBanDo. Khac hai test tren: xuLyClientOk KHONG goi s.Gui (chi
// ghi nhan b.clientSanSang = true va log), nen khong co goi phan hoi nao de
// doc — phai kiem TRANG THAI PHIEN (b.clientSanSang) thay vi goi tin.
//
// dungPhienThat nhan thang *BoXuLy do test tao (khong tu tao ben trong), nen
// giu duoc con tro b o day de doc lai truong sau khi gui goi.
//
// Ve khoang nghi: mynet.PhucVu chay b.XuLy() TREN GOROUTINE RIENG (khac
// goroutine cua test) — WriteMessage() cua client tra ve ngay khi may chu DOC
// xong goi (net.Pipe dong bo o muc Read/Write), nhung KHONG dam bao may chu
// da XU LY xong (goi toi xuLyClientOk, gan b.clientSanSang) tai thoi diem do.
// Cho mot khoang nghi ngan (net.Pipe khong co do tre mang that nen vai chuc ms
// la du du) truoc khi doc b.clientSanSang — dung cach phienGia trong
// nhanvat_test.go da lam cho tinh huong tuong tu (doi goroutine doc nen kip
// gom goi).
// fClientHetHanNgan khang dinh KHONG con goi nao may chu gui ra: doc mot goi va
// mong doi loi (ket noi ranh -> ReadMessage het han theo deadline da dat trong
// dungPhienThat). Tra ve true neu dung la khong co goi.
func fClientHetHanNgan(t *testing.T, f *protocol.Framer) bool {
	t.Helper()
	if _, err := f.ReadMessage(); err == nil {
		t.Error("nhan duoc goi NGOAI DU KIEN khi mong doi im lang")
		return false
	}
	return true
}

func TestClientOkGuiVaoTheGioi(t *testing.T) {
	// Client gui -28/13 khi da vao the gioi -> may chu gui DUNG 3 goi theo thu tu
	// -82 (tileset) -> -30 (ME_INFO) -> -24 (MAP INFO). -82 PHAI truoc -24 (§7).
	b := NewBoXuLy(&khoGia{}, nil)
	b.nhanVat = &store.NhanVat{ID: 7, Ten: "Naruto", Lang: 0, KieuToc: 5,
		SucManh: 100, Vang: 1000, Ngoc: 0, MapID: 1, ViTriX: 100, ViTriY: 384}
	fClient, dong := dungPhienThat(t, b)
	defer dong()

	// Gui -28/13 (client ok).
	w := protocol.NewWriter()
	w.Byte(lenhConClientOk) // 13
	if err := fClient.WriteMessage(protocol.NewMessage(lenhKhongTrenBanDo, w)); err != nil {
		t.Fatal(err)
	}

	m82 := docMotGoi(t, fClient)
	if m82.Command != -82 {
		t.Fatalf("goi 1 = lenh %d, mong doi -82 (tileset, BAT BUOC truoc -24)", m82.Command)
	}
	m30 := docMotGoi(t, fClient)
	if m30.Command != -30 {
		t.Fatalf("goi 2 = lenh %d, mong doi -30 (ME_INFO)", m30.Command)
	}
	m24 := docMotGoi(t, fClient)
	if m24.Command != -24 {
		t.Fatalf("goi 3 = lenh %d, mong doi -24 (MAP INFO)", m24.Command)
	}
	// MAP INFO byte dau (sau khong co lenh con) = mapId = 1.
	if len(m24.Data) == 0 || m24.Data[0] != 1 {
		t.Errorf("MAP INFO mapId = %v, mong doi 1 (mapId cua nhan vat)", m24.Data)
	}
}

func TestClientOkChuaVaoTheGioiKhongGui(t *testing.T) {
	// Chua co nhan vat (b.nhanVat == nil): -28/13 bi bo qua im lang, khong gui gi,
	// khong sap. Client that khong bao gio lam vay; chi ke gia lap moi gui.
	b := NewBoXuLy(&khoGia{}, nil) // b.nhanVat = nil
	fClient, dong := dungPhienThat(t, b)
	defer dong()

	w := protocol.NewWriter()
	w.Byte(lenhConClientOk)
	if err := fClient.WriteMessage(protocol.NewMessage(lenhKhongTrenBanDo, w)); err != nil {
		t.Fatal(err)
	}
	// Khong duoc gui goi nao: dat han doc ngan, mong doi timeout.
	_ = fClientHetHanNgan(t, fClient)
}

func TestNhanClientOk(t *testing.T) {
	b := NewBoXuLy(&khoGia{}, nil)
	b.nhanVat = &store.NhanVat{ID: 1, Lang: 0}
	fClient, dong := dungPhienThat(t, b)
	defer dong()
	// Client bao san sang: -28 lenh con 13.
	w := protocol.NewWriter()
	w.Byte(13)
	if err := fClient.WriteMessage(protocol.NewMessage(-28, w)); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	if !b.clientSanSang {
		t.Error("gui -28/13 (clientOk) khong lam b.clientSanSang = true")
	}
}

// ---------------------------------------------------------------------------
// Vá chốt chống gọi lại cho xuLyClientOk (-28/13) va xuLyVaoKhuVuc (-39).
//
// Boi canh: soi phat hien ca hai handler nay KHONG co chot chong tai nhap. Loi
// quan trong (tranh chap): neu client gui LAI -28/13 sau khi da gui -39 (nhan
// vat da "giao" cho actor khu vuc quan ly vi tri), xuLyClientOk chay lai se
// GOI goiMapInfo(b.nhanVat, bd) — DOC nv.ViTriX/Y tren luong phien, CUNG LUC
// actor co the dang GHI hai truong do (xuLyDiChuyen -> suKienDiChuyen) — doc
// cap toa do "xe" (nua cu nua moi). Day dung loai tranh chap Ke hoach 5 Task 7
// da xoa o noi khac; hai handler vao-the-gioi nay bi sot.
// ---------------------------------------------------------------------------

// TestClientOkKhongChayLaiLanHai la kiem chung QUAN TRONG cho chot chong goi
// lai cua xuLyClientOk. Gui -28/13 lan 1 phai nhan du 3 goi -82/-30/-24 (dung
// nhu TestClientOkGuiVaoTheGioi). Gui LAI -28/13 lan 2 (client that khong lam
// vay, nhung mot con gia hoac goi lap tren mang co the) PHAI bi bo qua HOAN
// TOAN im lang — khong duoc gui THEM bat ky goi nao, chung to xuLyClientOk
// khong chay lai than ham (tranh doc nv.ViTriX/Y tranh chap voi actor ghi).
func TestClientOkKhongChayLaiLanHai(t *testing.T) {
	b := NewBoXuLy(&khoGia{}, nil)
	b.nhanVat = &store.NhanVat{ID: 7, Ten: "Naruto", Lang: 0, KieuToc: 5,
		SucManh: 100, Vang: 1000, Ngoc: 0, MapID: 1, ViTriX: 100, ViTriY: 384}
	fClient, dong := dungPhienThat(t, b)
	defer dong()

	guiClientOk := func() {
		t.Helper()
		w := protocol.NewWriter()
		w.Byte(lenhConClientOk) // 13
		if err := fClient.WriteMessage(protocol.NewMessage(lenhKhongTrenBanDo, w)); err != nil {
			t.Fatalf("gui -28/13 that bai: %v", err)
		}
	}

	// Lan 1: phai nhan du 3 goi theo dung thu tu (nhu TestClientOkGuiVaoTheGioi).
	guiClientOk()
	m82 := docMotGoi(t, fClient)
	if m82.Command != -82 {
		t.Fatalf("lan 1, goi 1 = lenh %d, mong doi -82", m82.Command)
	}
	m30 := docMotGoi(t, fClient)
	if m30.Command != -30 {
		t.Fatalf("lan 1, goi 2 = lenh %d, mong doi -30", m30.Command)
	}
	m24 := docMotGoi(t, fClient)
	if m24.Command != -24 {
		t.Fatalf("lan 1, goi 3 = lenh %d, mong doi -24", m24.Command)
	}

	// Lan 2: gui LAI -28/13 — PHAI khong nhan them goi nao (chot chong goi lai).
	guiClientOk()
	if !fClientHetHanNgan(t, fClient) {
		t.Fatal("gui -28/13 lan hai van khien xuLyClientOk gui them goi tin — thieu chot chong goi lai")
	}
}

// TestVaoKhuVucLanHaiKhongPhatGoiTrung la kiem chung cho chot chong goi lai
// cua xuLyVaoKhuVuc. Kich ban: client gui -39 (bao da tai xong map) HAI LAN
// LIEN TIEP. Neu khong co chot: lan -39 THU HAI se khien actor xu ly LAI
// suKienVao — buoc 1 cua KhuVuc.xuLy (xem khuvuc.go) phat -5 cua TUNG "nguoi
// cu" trong khu vuc cho nguoi moi; nhung nguoi "cu" do gio CHINH LA nguoi nay
// (da duoc them vao khu vuc boi lan -39 thu nhat) — ket qua actor gui cho
// chinh nguoi choi goi -5 CUA CHINH HO, mot goi trung/nhieu thua tren day.
// Voi chot dung, lan -39 thu hai la no-op hoan toan, khong co goi nao them.
func TestVaoKhuVucLanHaiKhongPhatGoiTrung(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.BatDau()
	defer k.Dung()
	soKhuVuc := &SoKhuVuc{theoMap: map[int32]*KhuVuc{1: k}}

	b := NewBoXuLy(&khoGia{}, soKhuVuc)
	b.nhanVat = &store.NhanVat{ID: 7, Ten: "Naruto", Lang: 0, KieuToc: 5,
		SucManh: 100, Vang: 1000, Ngoc: 0, MapID: 1, ViTriX: 100, ViTriY: 384}
	fClient, dong := dungPhienThat(t, b)
	defer dong()

	// Bat tay -28/13 truoc de b.khuVuc duoc gan (xuLyClientOk). Doc het 3 goi.
	w := protocol.NewWriter()
	w.Byte(lenhConClientOk)
	if err := fClient.WriteMessage(protocol.NewMessage(lenhKhongTrenBanDo, w)); err != nil {
		t.Fatalf("gui -28/13 that bai: %v", err)
	}
	docMotGoi(t, fClient) // -82
	docMotGoi(t, fClient) // -30
	docMotGoi(t, fClient) // -24

	// Gui -39 (vao khu vuc) HAI LAN LIEN TIEP — mo phong client gui trung.
	if err := fClient.WriteMessage(&protocol.Message{Command: lenhVaoKhuVuc}); err != nil {
		t.Fatalf("gui -39 lan 1 that bai: %v", err)
	}
	if err := fClient.WriteMessage(&protocol.Message{Command: lenhVaoKhuVuc}); err != nil {
		t.Fatalf("gui -39 lan 2 that bai: %v", err)
	}

	if !fClientHetHanNgan(t, fClient) {
		t.Fatal("gui -39 lan hai lam actor phat them goi tin (vd -5 cua chinh nguoi nay cho chinh ho) — thieu chot chong goi lai")
	}
}

// TestXuLyChonChieuGuiSuKienKhongGhiThang la kiem chung QUAN TRONG cho ky
// luat actor (dac ta §3.2 + brief KH6 Task 2): xuLyChonChieu (luong phien) CHI
// duoc gui su kien vao hop thu actor, TUYET DOI KHONG tu tra ChieuCuaLang/ghi
// chieuDangChon truc tiep.
//
// Ky thuat kiem: tao KhuVuc nhung KHONG BatDau() (khong co goroutine vongLap
// nao chay song song) — chi tu gan k.hopThu bang tay. Vi khong co goroutine
// thu hai, doc lai hopThu ngay trong test la AN TOAN (khong dua tranh voi
// actor), dung nhu k.xuLy() duoc goi truc tiep o khuvuc_test.go. Sau do kiem
// k.nguoiTrong VAN RONG — chung minh xuLyChonChieu khong dung truc tiep vao
// trang thai actor so huu.
func TestXuLyChonChieuGuiSuKienKhongGhiThang(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.hopThu = make(chan suKienKhuVuc, 4) // khong BatDau: tu cap hop thu de doc lai an toan trong cung goroutine

	b := NewBoXuLy(&khoGia{}, nil)
	b.khuVuc = k
	b.nhanVat = &store.NhanVat{ID: 42, Ten: "A", Lang: 0}

	w := protocol.NewWriter()
	w.Short(2) // templateId
	if err := b.xuLyChonChieu(protocol.NewMessage(lenhChonChieu, w)); err != nil {
		t.Fatalf("xuLyChonChieu loi: %v", err)
	}

	select {
	case sk := <-k.hopThu:
		ev, ok := sk.(suKienChonChieu)
		if !ok {
			t.Fatalf("su kien gui vao hop thu la %T, mong doi suKienChonChieu", sk)
		}
		if ev.idNhanVat != 42 || ev.mauChieuId != 2 {
			t.Fatalf("suKienChonChieu = %+v, mong doi {idNhanVat:42 mauChieuId:2}", ev)
		}
	default:
		t.Fatal("xuLyChonChieu khong gui su kien nao vao hop thu actor")
	}

	if len(k.nguoiTrong) != 0 {
		t.Fatal("xuLyChonChieu khong duoc tu dong ghi vao k.nguoiTrong — do la viec cua actor, khong phai luong phien")
	}
}

// TestXuLyChonChieuChuaVaoTheGioiKhongPanic: b.khuVuc/b.nhanVat con nil (chua
// bat tay / chua vao khu vuc) -> xuLyChonChieu phai la no-op an toan.
func TestXuLyChonChieuChuaVaoTheGioiKhongPanic(t *testing.T) {
	b := NewBoXuLy(&khoGia{}, nil)
	w := protocol.NewWriter()
	w.Short(1)
	if err := b.xuLyChonChieu(protocol.NewMessage(lenhChonChieu, w)); err != nil {
		t.Fatalf("xuLyChonChieu voi nhanVat/khuVuc nil phai tra nil, nhan duoc: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Ke hoach 11 Task 5 — KhiDong(): luu nguoi choi khi phien mat ket noi, dac ta
// 2026-07-27-luu-trang-thai-song.md §4 diem B. THU TU BAT BUOC: KhiDong PHAI
// xin anh chup (AnhChupNguoiChoi) TRUOC khi bao actor go nguoi choi ra
// (khuVuc.Gui(suKienRa)) — case suKienRa (khuvuc.go) XOA nguoi choi khoi
// k.nguoiTrong NGAY khi actor xu ly, nen goi nguoc thu tu se chup phai mot
// khu vuc DA XOA nguoi do -> anh chup rong -> mat sach tien do vua choi. Xem
// chu thich day du tai KhiDong (login.go).
// ---------------------------------------------------------------------------

// ghiKhiRoiGia la "tho ghi gia" thoa man interface GhiKhiRoi (login.go) — ghi
// nhan NGAY vao slice trong bo nho, KHONG di qua goroutine nen + ticker that
// cua ThoGhi that (thogi.go) — tranh test phai cho/doan thoi gian ghi dia.
// Dung khuon voi khoGia (dau file nay) / khoLuuGia (thogi_test.go): fake toi
// thieu, chi lam dung viec BoXuLy can toi.
type ghiKhiRoiGia struct {
	mu   sync.Mutex
	nhan []store.NhanVat
}

func (g *ghiKhiRoiGia) Xep(nv store.NhanVat) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.nhan = append(g.nhan, nv)
}

func (g *ghiKhiRoiGia) banDaNhan() []store.NhanVat {
	g.mu.Lock()
	defer g.mu.Unlock()
	ban := make([]store.NhanVat, len(g.nhan))
	copy(ban, g.nhan)
	return ban
}

// TestKhiDongLuuViTriMoiNhatKhiDaVaoKhuVuc la Task 5 buoc 1 kich ban 1 (dac ta
// §4 diem B, §10 nghiem thu muc 2): "client gia" (mot store.NhanVat dua thang
// vao khu vuc qua suKienVao, dung khuon voi khuVucAnhChupGia o khuvuc_test.go)
// DI CHUYEN QUA ACTOR toi mot vi tri ghim cung, roi "ngat ket noi" bang cach
// goi KhiDong() TRUC TIEP — dung y het cach mynet.PhucVu goi no qua defer luc
// vong doc gap loi (session.go:104-106; xem ca TestTichHop_RoiKhuVucPhatGoi6
// o khuvuc_tichhop_test.go cho kich ban tuong tu nhung qua mang that).
//
// Ky vong: tho ghi gia nhan DUNG MOT ban ghi, vi tri la VI TRI MOI NHAT (sau
// di chuyen) — KHONG PHAI vi tri luc vao khu vuc.
func TestKhiDongLuuViTriMoiNhatKhiDaVaoKhuVuc(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.quai = nil // tat AI quai (dung tinh than moSoKhuVucKhongQuaiMacDinh, khuvuc_tichhop_test.go): test nay chi kiem VI TRI, khong de quai xen vao gay nhieu ngoai y muon
	k.BatDau()
	t.Cleanup(func() { _ = k.Dung() })

	viTriLucVaoX, viTriLucVaoY := int32(100), int32(100)
	viTriGhimCungX, viTriGhimCungY := int16(300), int16(300) // cach (100,100) ~283 don vi < NguongBuoc=500 (dichuyen.go) -> hop le

	nv := &store.NhanVat{ID: 55, Ten: "RoiKhuVuc", Lang: 0, MapID: 1, ViTriX: viTriLucVaoX, ViTriY: viTriLucVaoY}
	k.Gui(suKienVao{nv: nv, phien: &phienGiaKhuVuc{}})

	// Gui di chuyen QUA ACTOR (khong tu gan nv.ViTriX/Y tu ben ngoai — actor
	// moi la nguon authoritative, dac ta §2.1). k.Gui la mot hop thu (channel)
	// FIFO CHI mot nguoi gui (goroutine test nay) nen suKienVao roi
	// suKienDiChuyen CHAC CHAN duoc actor xu ly DUNG THU TU da gui — khong can
	// sleep/poll de "cho actor kip xu ly".
	k.Gui(suKienDiChuyen{idNhanVat: nv.ID, x: viTriGhimCungX, y: viTriGhimCungY})

	ghiGia := &ghiKhiRoiGia{}
	b := &BoXuLy{khuVuc: k, nhanVat: nv, daVaoKhuVuc: true, thoGhi: ghiGia}

	b.KhiDong() // = "ngat ket noi"

	nhan := ghiGia.banDaNhan()
	if len(nhan) != 1 {
		t.Fatalf("tho ghi gia nhan %d ban ghi, mong doi 1", len(nhan))
	}
	if nhan[0].ID != nv.ID {
		t.Fatalf("ban ghi sai nguoi: ID=%d, mong doi %d", nhan[0].ID, nv.ID)
	}
	if nhan[0].ViTriX != int32(viTriGhimCungX) || nhan[0].ViTriY != int32(viTriGhimCungY) {
		t.Fatalf("ban ghi vi tri = (%d,%d), mong doi VI TRI MOI NHAT sau di chuyen (%d,%d) — khong phai vi tri luc vao (%d,%d)",
			nhan[0].ViTriX, nhan[0].ViTriY, viTriGhimCungX, viTriGhimCungY, viTriLucVaoX, viTriLucVaoY)
	}

	// Bang chung THEM cho LY DO buoc thu tu bat buoc (xem chu thich KhiDong):
	// sau KhiDong, actor DA xu ly xong suKienRa va xoa nv khoi k.nguoiTrong —
	// goi AnhChupNguoiChoi LAI (mo phong ai do lo dao thu tu, chup SAU khi da
	// gui suKienRa) gio se tra RONG. Day chinh la ly do KhiDong PHAI chup
	// truoc.
	banSau := k.AnhChupNguoiChoi(nv.ID)
	if len(banSau) != 0 {
		t.Fatalf("sau KhiDong, AnhChupNguoiChoi(%d) van tra %d ban — actor le ra da xoa nguoi nay khoi k.nguoiTrong (case suKienRa, khuvuc.go)", nv.ID, len(banSau))
	}
}

// TestKhiDongSnapshotVaLeaveKhongRoiKhiGameplayDay bat hoi quy dua hai control
// event nay ve k.Gui best-effort: mailbox day se lam mat ban luu va ghost.
func TestKhiDongSnapshotVaLeaveKhongRoiKhiGameplayDay(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	nv := &store.NhanVat{ID: 56, Ten: "MailboxDay", MapID: 1, ViTriX: 123, ViTriY: 384}
	k.xuLy(suKienVao{nv: nv, phien: &phienGiaKhuVuc{}})

	k.hopThu = make(chan suKienKhuVuc, 256)
	k.hopThuDieuKhien = make(chan suKienDieuKhien, 16)
	k.xong = make(chan struct{})
	for i := 0; i < cap(k.hopThu); i++ {
		k.hopThu <- suKienRa{idNhanVat: 999}
	}
	go k.vongLap()
	t.Cleanup(func() { _ = k.Dung() })

	ghiGia := &ghiKhiRoiGia{}
	b := &BoXuLy{khuVuc: k, nhanVat: nv, daVaoKhuVuc: true, thoGhi: ghiGia}
	b.KhiDong()

	nhan := ghiGia.banDaNhan()
	if len(nhan) != 1 || nhan[0].ID != nv.ID || nhan[0].ViTriX != 123 {
		t.Fatalf("snapshot khi gameplay day = %+v, mong dung ban moi nhat cua %d", nhan, nv.ID)
	}
	banSau, err := k.AnhChupNguoiChoiCoLoi(nv.ID)
	if err != nil {
		t.Fatalf("probe actor sau KhiDong: %v", err)
	}
	if len(banSau) != 0 {
		t.Fatalf("con ghost sau KhiDong: %+v", banSau)
	}
}

// TestKhiDongChuaVaoKhuVucKhongLuuKhongPanic la Task 5 buoc 1 kich ban 2: mot
// phien CHUA thuc su vao khu vuc (b.daVaoKhuVuc == false) — du dang o buoc
// nao (chua dang nhap, moi co nhan vat nhung chua chon khu vuc, hay da chon
// khu vuc nhung chua gui -39) — bi "ngat ket noi" thi KhiDong KHONG duoc goi
// Xep (khong luu bua) va KHONG panic.
func TestKhiDongChuaVaoKhuVucKhongLuuKhongPanic(t *testing.T) {
	ghiGia := &ghiKhiRoiGia{}
	kChuaVao := MoKhuVuc(TimBanDo(1)) // mo phong da xuLyClientOk (co khuVuc) nhung chua gui -39

	truongHop := []struct {
		ten string
		b   *BoXuLy
	}{
		{"chua dang nhap (khuVuc va nhanVat deu nil)", &BoXuLy{thoGhi: ghiGia}},
		{"da co nhan vat nhung chua tung chon khu vuc", &BoXuLy{
			nhanVat: &store.NhanVat{ID: 9, Ten: "ChuaVao"}, thoGhi: ghiGia,
		}},
		{"da chon khu vuc nhung chua gui -39 (daVaoKhuVuc=false)", &BoXuLy{
			khuVuc: kChuaVao, nhanVat: &store.NhanVat{ID: 10, Ten: "ChuaVao2"}, thoGhi: ghiGia,
		}},
	}

	for _, tc := range truongHop {
		t.Run(tc.ten, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("KhiDong panic: %v", r)
				}
			}()
			tc.b.KhiDong()
		})
	}

	if nhan := ghiGia.banDaNhan(); len(nhan) != 0 {
		t.Fatalf("tho ghi gia nhan %d ban ghi, mong doi 0 (chua vao khu vuc thi khong duoc luu)", len(nhan))
	}
}
