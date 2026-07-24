package game

import (
	"net"
	"sync"
	"testing"
	"time"

	mynet "nroserver/internal/net"
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// ---------------------------------------------------------------------------
// Helper dung chung cho ca file nay: boXuLyDaDangNhap, phienGia.
//
// login_test.go (Ke hoach 1) moi chi co dungPhienThat — dung cho kiem thu
// dau-cuoi di qua XuLy() that. Task 6 can goi THANG b.xuLyTaoNhanVat(s, r),
// bo qua tang dieu phoi -28/lenh con, nen phai co MOT *mynet.Session that su
// hoat dong duoc (khong the tu dung mynet.Session{conn:..., framer:...} tu goi
// game — cac truong do khong xuat) cong voi mot cach doc lai goi may chu gui
// ra. Trich ca hai o day de dung chung voi cac test tuong lai neu can.
// ---------------------------------------------------------------------------

// boXuLyDaDangNhap dung mot BoXuLy voi kho gia rong va b.taiKhoan da duoc gan
// san, mo phong mot phien DA dang nhap thanh cong (Ke hoach 1/Task 5), san
// sang goi thang cac ham xuLy* cua Ke hoach 2 ma khong phai chay lai
// KiemTraDangNhap.
func boXuLyDaDangNhap(t *testing.T) *BoXuLy {
	t.Helper()
	b := NewBoXuLy(&khoGia{}, nil)
	b.taiKhoan = &store.TaiKhoan{ID: 1, TenDangNhap: "nguoidung"}
	return b
}

// boBatSession la mot mynet.XuLyGoiTin GIA, dung DUY NHAT de "tom" lai chinh
// con tro *mynet.Session ma PhucVu tao ra ben trong cho MOI ket noi. Phai lam
// vay vi cac truong cua mynet.Session (conn, framer) khong xuat — goi game
// khong the tu dung len mot Session hoat dong duoc, chi co the muon lai con
// tro that qua duong duy nhat PhucVu dua no ra ngoai: goi XuLy(s, m). Khong xu
// ly gi khac, khong gui goi tin nao ca.
type boBatSession struct {
	ket chan *mynet.Session
}

func (x *boBatSession) XuLy(s *mynet.Session, m *protocol.Message) error {
	select {
	case x.ket <- s:
	default:
	}
	return nil
}

// phienGia dung mot phien mynet THAT qua net.Pipe (bat tay xong, khoa da dat
// ca hai chieu), bat lai chinh *mynet.Session phia may chu bang boBatSession,
// va tra ve ham doc TAT CA goi tin may chu da gui ra.
//
// Vi sao phai co goroutine doc NEN chay san (khong doi den luc goi docGoi()
// moi bat dau doc): net.Pipe KHONG co bo dem — Write() o mot dau CHAN cho toi
// khi dau kia Read(). Cac test goi b.xuLyTaoNhanVat(s, r) THANG tu goroutine
// cua chinh test (khong qua goroutine rieng cua PhucVu), nen neu s.Gui() ben
// trong can ghi mot goi tra ve ma chua co ai dang doc, ca hai se ket cung mot
// cho — deadlock (chi thoat sau 30s han gio ghi cua Session.Gui, lam test that
// bai rat cham va kho hieu). Cho nen goroutine doc phai chay TU LUC phienGia
// tra ve, doc lien tuc vao mot hang doi; docGoi() chi doi mot chut roi lay
// nhung gi da gom duoc — khong tu minh goi ReadMessage().
func phienGia(t *testing.T) (s *mynet.Session, docGoi func() []*protocol.Message) {
	t.Helper()

	benClient, benServer := net.Pipe()

	batDuoc := make(chan *mynet.Session, 1)
	go mynet.PhucVu(benServer, func() mynet.XuLyGoiTin { return &boBatSession{ket: batDuoc} })

	benClient.SetReadDeadline(time.Now().Add(5 * time.Second))

	fClient := protocol.NewFramer(benClient, benClient)
	if err := fClient.WriteMessage(&protocol.Message{Command: -27}); err != nil {
		t.Fatalf("gui goi bat tay that bai: %v", err)
	}
	khoaTho := docGoiKhoaThuCong(t, benClient)
	fClient.SetKey(protocol.DeriveKey(khoaTho))

	// Gui MOT goi bat ky (than rong) de PhucVu goi toi XuLy() dung mot lan,
	// qua do boBatSession bat duoc con tro Session that. Lenh dung o day khong
	// quan trong — boBatSession khong doc than tin, chi "tom" con tro Session.
	if err := fClient.WriteMessage(&protocol.Message{Command: 0}); err != nil {
		t.Fatalf("gui goi bat session that bai: %v", err)
	}

	select {
	case s = <-batDuoc:
	case <-time.After(5 * time.Second):
		t.Fatal("khong bat duoc *mynet.Session tu PhucVu")
	}

	var mu sync.Mutex
	var hangDoi []*protocol.Message
	go func() {
		for {
			m, err := fClient.ReadMessage()
			if err != nil {
				return // pipe dong hoac het han gio doc -> dung, khong con gi de doc
			}
			mu.Lock()
			hangDoi = append(hangDoi, m)
			mu.Unlock()
		}
	}()

	t.Cleanup(func() { benClient.Close() })

	docGoi = func() []*protocol.Message {
		// Doi mot chut de goroutine doc nen kip nhan het goi da gui (net.Pipe
		// dong bo, khong co do tre mang that nen vai chuc ms la du du).
		time.Sleep(100 * time.Millisecond)
		mu.Lock()
		defer mu.Unlock()
		out := make([]*protocol.Message, len(hangDoi))
		copy(out, hangDoi)
		return out
	}
	return s, docGoi
}

// thanTaoNhanVat dong goi than tin goi -28/2 (KHONG gom byte lenh -28 lan byte
// lenh con 2, vi xuLyTaoNhanVat nhan Reader da tro qua ca hai — xem
// xuLyKhongTrenBanDo trong login.go). Tra ve []byte than sau lenh con.
func thanTaoNhanVat(t *testing.T, ten string, lang, toc byte) []byte {
	t.Helper()
	w := protocol.NewWriter()
	if err := w.UTF(ten); err != nil {
		t.Fatal(err)
	}
	w.Byte(lang)
	w.Byte(toc)
	return w.Bytes() // ham xuat byte that cua Writer (xem datastream.go)
}

// kiemTraDungMotGoiTuChoi khang dinh may chu gui ve DUNG MOT goi, va goi do la
// -26 (hop thoai OK). Sieu QUAN TRONG phai kiem CA SO LUONG: guiHopThoai chi
// duoc gui DUY NHAT mot -26, KHONG duoc kem them -102 (khac guiTuChoi cua man
// dang nhap — xem chu thich guiHopThoai trong nhanvat.go). Neu chi kiem
// goi[0].Command == -26 ma khong kiem len(goi) == 1, mot dot bien them -102
// thua vao guiHopThoai van lot qua duoc — day chinh la lo hong da bi bat boi
// kiem thu dot bien.
func kiemTraDungMotGoiTuChoi(t *testing.T, goi []*protocol.Message) {
	t.Helper()
	if len(goi) != 1 {
		t.Fatalf("phai nhan DUNG MOT goi (guiHopThoai KHONG duoc kem goi nao khac), nhan %d goi: %+v", len(goi), goi)
	}
	if goi[0].Command != -26 {
		t.Fatalf("goi nhan duoc co lenh %d, mong doi -26 (hop thoai OK)", goi[0].Command)
	}
}

func TestTaoNhanVat_TenQuaNgan_TuChoi(t *testing.T) {
	b := boXuLyDaDangNhap(t)
	s, docGoi := phienGia(t)
	r := protocol.NewReader(thanTaoNhanVat(t, "abc", 0, 5)) // 3 ky tu < 5
	if err := b.xuLyTaoNhanVat(s, r); err != nil {
		t.Fatalf("xuLyTaoNhanVat loi: %v", err)
	}
	kiemTraDungMotGoiTuChoi(t, docGoi())
}

// TestTaoNhanVat_TenQuaDai_TuChoi bat dot bien noi long gioi han tren cua
// reTenNhanVat (vi du {5,10} -> {5,50}): khong test nao truoc day thu mot ten
// VUOT can tren, chi co ten qua NGAN (duoi can duoi).
func TestTaoNhanVat_TenQuaDai_TuChoi(t *testing.T) {
	b := boXuLyDaDangNhap(t)
	s, docGoi := phienGia(t)
	r := protocol.NewReader(thanTaoNhanVat(t, "Naruto12345", 0, 5)) // 11 ky tu > 10
	if err := b.xuLyTaoNhanVat(s, r); err != nil {
		t.Fatalf("xuLyTaoNhanVat loi: %v", err)
	}
	kiemTraDungMotGoiTuChoi(t, docGoi())
}

func TestTaoNhanVat_KyTuLa_TuChoi(t *testing.T) {
	b := boXuLyDaDangNhap(t)
	s, docGoi := phienGia(t)
	r := protocol.NewReader(thanTaoNhanVat(t, "na ru to", 0, 5)) // co dau cach
	_ = b.xuLyTaoNhanVat(s, r)
	kiemTraDungMotGoiTuChoi(t, docGoi())
}

// TestTaoNhanVat_TenCam_TuChoi bat dot bien bo han nhanh kiem danhSachTenCam:
// "admin" dai dung 5-10 ky tu, chi gom chu cai (khop reTenNhanVat), nen NEU
// thieu buoc kiem danh sach cam thi se lot qua toi tan buoc TenNhanVatDaDung.
func TestTaoNhanVat_TenCam_TuChoi(t *testing.T) {
	b := boXuLyDaDangNhap(t)
	s, docGoi := phienGia(t)
	r := protocol.NewReader(thanTaoNhanVat(t, "admin", 0, 5)) // trong danhSachTenCam
	if err := b.xuLyTaoNhanVat(s, r); err != nil {
		t.Fatalf("xuLyTaoNhanVat loi: %v", err)
	}
	kiemTraDungMotGoiTuChoi(t, docGoi())
}

func TestTaoNhanVat_LangKhongHopLe_TuChoi(t *testing.T) {
	b := boXuLyDaDangNhap(t)
	s, docGoi := phienGia(t)
	r := protocol.NewReader(thanTaoNhanVat(t, "Naruto", 99, 5)) // lang 99 khong ton tai
	_ = b.xuLyTaoNhanVat(s, r)
	kiemTraDungMotGoiTuChoi(t, docGoi())
}

func TestTaoNhanVat_TenDaDung_TuChoi(t *testing.T) {
	b := boXuLyDaDangNhap(t)
	b.kho.(*khoGia).tenDaDung = true
	s, docGoi := phienGia(t)
	r := protocol.NewReader(thanTaoNhanVat(t, "Naruto", 0, 5))
	_ = b.xuLyTaoNhanVat(s, r)
	kiemTraDungMotGoiTuChoi(t, docGoi())
}

func TestTaoNhanVat_ChuaDangNhap_BoQua(t *testing.T) {
	b := NewBoXuLy(&khoGia{}, nil) // b.taiKhoan == nil: chua dang nhap
	s, docGoi := phienGia(t)
	r := protocol.NewReader(thanTaoNhanVat(t, "Naruto", 0, 5))
	if err := b.xuLyTaoNhanVat(s, r); err != nil {
		t.Fatalf("khong duoc tra loi khi chua dang nhap, chi bo qua: %v", err)
	}
	if goi := docGoi(); len(goi) != 0 {
		t.Errorf("chua dang nhap ma gui create -> phai IM LANG, khong gui gi; nhan %+v", goi)
	}
}

func TestTaoNhanVat_ThanhCong_VaoTheGioi(t *testing.T) {
	b := boXuLyDaDangNhap(t)
	b.kho.(*khoGia).nhanVatVuaTao = &store.NhanVat{ID: 42, Ten: "Naruto", Lang: 0, MapID: 1}
	s, _ := phienGia(t)
	r := protocol.NewReader(thanTaoNhanVat(t, "Naruto", 0, 5))
	if err := b.xuLyTaoNhanVat(s, r); err != nil {
		t.Fatalf("xuLyTaoNhanVat loi: %v", err)
	}
	// vaoTheGioi (stub) phai gan nhan vat vao phien.
	if b.nhanVat == nil || b.nhanVat.ID != 42 {
		t.Errorf("tao thanh cong phai gan nhan vat vao phien, b.nhanVat=%+v", b.nhanVat)
	}
}

// TestXuLy_Khung28_DieuPhoiToiTaoNhanVat la kiem thu DAU-CUOI duy nhat cua
// Task 6 di qua chinh XuLy(s, m) that (khong goi thang xuLyTaoNhanVat nhu cac
// test o tren) — dung dungPhienThat (login_test.go) giong het cac test dang
// nhap dau-cuoi, nen KHONG can lo deadlock net.Pipe: PhucVu chay tren goroutine
// rieng cua no, doc/ghi doc lap voi goroutine cua test.
//
// Bat dot bien cho xuLyKhongTrenBanDo (vi du bo mat nhanh "case
// lenhConTaoNhanVat: return b.xuLyTaoNhanVat(s, r)"): neu day dieu phoi
// -28 -> lenh con 2 -> xuLyTaoNhanVat bi dut o bat ky khau nao, se KHONG co
// goi -26 nao duoc gui ve va ReadMessage() o day se timeout/loi, lam test do.
func TestXuLy_Khung28_DieuPhoiToiTaoNhanVat(t *testing.T) {
	b := NewBoXuLy(&khoGia{}, nil)
	b.taiKhoan = &store.TaiKhoan{ID: 1, TenDangNhap: "nguoidung"} // da dang nhap

	fClient, dong := dungPhienThat(t, b)
	defer dong()

	// Dong goi khung -28 dung y het client that: byte lenh con (2 = tao nhan
	// vat) roi toi ten/lang/kieu toc. Dung ten "abc" (qua ngan) — chi can
	// bat duoc RANG BUOC xuLyTaoNhanVat co duoc goi toi hay khong, khong quan
	// tam nhanh tu choi nao kich hoat.
	w := protocol.NewWriter()
	w.Byte(lenhConTaoNhanVat) // 2
	if err := w.UTF("abc"); err != nil {
		t.Fatal(err)
	}
	w.Byte(0) // lang
	w.Byte(5) // kieu toc

	if err := fClient.WriteMessage(protocol.NewMessage(-28, w)); err != nil {
		t.Fatalf("gui goi -28 that bai: %v", err)
	}

	m, err := fClient.ReadMessage()
	if err != nil {
		t.Fatalf("khong doc duoc goi tra ve tu goi -28: %v — neu bo dieu phoi -28 khong goi toi xuLyTaoNhanVat, se khong co goi nao duoc gui va dong nay se timeout", err)
	}
	if m.Command != -26 {
		t.Fatalf("lenh = %d, mong doi -26 (dieu phoi -28 -> lenh con 2 -> xuLyTaoNhanVat)", m.Command)
	}
}
