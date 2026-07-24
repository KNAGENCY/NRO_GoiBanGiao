package game

// Kiem thu tich hop dau-cuoi (Task 7, Ke hoach 2): dung CSDL PostgreSQL that
// (nro_test) va store.Store that (KHONG phai khoGia), di qua net.Pipe +
// mynet.PhucVu that giong het cach client that noi — tao nhan vat roi dang
// nhap lai phai thay dung nhan vat do, va moi duong tu choi phai tra dung
// -26.
//
// Vi sao mot file rieng thay vi ghep vao login_test.go/nhanvat_test.go: hai
// file do dung khoGia (kho gia lap trong bo nho) cho hau het test — nhanh va
// khong can CSDL. File nay la lop kiem thu KHAC HAN ve ban chat: no la phep
// thu "tin tuong nhung xac minh" cho toan bo chuoi that (Go + mang + CSDL that
// su), nen tach rieng de ai doc cung thay ngay day la nhom test can
// NRO_TEST_DSN va se cham toi CSDL that.
//
// Bo qua toan bo file (tung test tu goi t.Skip qua moKhoThatChoTest) neu chua
// dat NRO_TEST_DSN, dung y het loi package store dang lam (xem moKetNoiThu
// trong internal/store/store_test.go) — de ai chay `go test ./...` ma khong
// co CSDL van khong bi do.

import (
	"context"
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"nroserver/internal/auth"
	mynet "nroserver/internal/net"
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// ---------------------------------------------------------------------------
// Ha tang dung chung cho file nay: mo CSDL that, chen tai khoan that, va mot
// "client gia" (clientGia) giu duoc net.Conn de tu dieu khien han gio doc.
// ---------------------------------------------------------------------------

// moKhoThatChoTest mo mot *store.Store THAT toi CSDL nro_test (qua bien moi
// truong NRO_TEST_DSN). Bo qua test (t.Skip) neu chua dat bien nay.
func moKhoThatChoTest(t *testing.T) *store.Store {
	t.Helper()
	dsn := os.Getenv("NRO_TEST_DSN")
	if dsn == "" {
		t.Skip("bo qua: chua dat NRO_TEST_DSN")
	}
	s, err := store.Open(context.Background(), dsn)
	if err != nil {
		t.Fatalf("khong mo duoc CSDL: %v", err)
	}
	t.Cleanup(s.Close)
	return s
}

// taoTaiKhoanCoMatKhau chen mot tai khoan THAT vao CSDL voi mat khau da BAM
// bang auth.Hash (dung dinh dang that server doi chieu luc dang nhap — luu
// mat khau tho se khong bao gio Verify() dung). Tra ve ID tai khoan vua tao.
//
// store.Store CO Y khong xuat ham "chen tai khoan tho" nao (tao tai khoan
// that su la qua kenh dang ky, ngoai pham vi Ke hoach 2) — nen o day phai tu
// mo mot pgxpool rieng va chen thang bang SQL, dung nhu ghi chu trong
// task-7-brief.md: "Trong game phai tu chen tai khoan qua SQL hoac qua kho
// that". Tu dang ky don dep qua t.Cleanup: xoa nhan_vat TRUOC (khoa ngoai toi
// tai_khoan) roi moi xoa tai_khoan — dung thu tu Cleanup chay LIFO nen dang
// ky pool.Close() TRUOC, dang ky xoa hang SAU, de luc chay cleanup thi xoa
// hang thuc hien truoc, dong pool sau cung.
func taoTaiKhoanCoMatKhau(t *testing.T, ten, matKhau string) int64 {
	t.Helper()
	dsn := os.Getenv("NRO_TEST_DSN")

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("khong mo duoc pool CSDL de chen tai khoan thu: %v", err)
	}
	t.Cleanup(pool.Close)

	bam, err := auth.Hash(matKhau)
	if err != nil {
		t.Fatalf("bam mat khau that bai: %v", err)
	}

	var id int64
	err = pool.QueryRow(context.Background(),
		`INSERT INTO tai_khoan (ten_dang_nhap, mat_khau_bam) VALUES ($1, $2) RETURNING id`,
		ten, bam,
	).Scan(&id)
	if err != nil {
		t.Fatalf("chen tai khoan thu %q that bai: %v", ten, err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM nhan_vat WHERE tai_khoan_id = $1`, id)
		_, _ = pool.Exec(context.Background(), `DELETE FROM tai_khoan WHERE id = $1`, id)
	})

	return id
}

// demHauToTichHop la bo dem nguyen tu RIENG cho file nay (khac demHauTo cua
// package store — hai package tach biet, khong dung chung bien duoc, cung
// khong can thiet vi moi package tu quan ly khong gian ten cua minh).
var demHauToTichHop int64

// hauToNganNgau tra mot chuoi so duy nhat cho moi lan goi, de ten tai
// khoan/nhan vat giua cac test (va giua cac lan chay lai) khong dung nhau.
// Dung bo dem nguyen tu (sync/atomic), KHONG dung dia chi con tro (%p co the
// trung khi stack duoc tai su dung) — giong het ly do da giai thich o
// hauToNganNgau trong internal/store/store_test.go. Hau to luon ngan (toi da
// vai chu so voi so luong test o day), nen "Naru"+hau_to luon <= 10 ky tu,
// khop CHECK nhan_vat_ten_hop_le (5-10 ky tu).
func hauToNganNgau(t *testing.T) string {
	t.Helper()
	return strconv.FormatInt(atomic.AddInt64(&demHauToTichHop, 1), 10)
}

// hanChoImLang la khoang thoi gian toi da clientGia doi de khang dinh "khong
// con goi nao them" (cac kich ban tu choi im lang, hoac thanh cong khong tra
// loi). Ngan hon NHIEU so voi han 5 giay co dinh ma dungPhienThat/chayDangNhap
// (login_test.go) dat MOT LAN cho ca doi phien:
//
//   - net.Pipe() khong co do tre mang that: Write() mot dau CHAN cho toi khi
//     dau kia Read(), nen phan mang cua duong di gan nhu tuc thi.
//   - Phan CO the co do tre that la cac lenh goi CSDL that (TenNhanVatDaDung,
//     TaoNhanVat) ben trong xuLyTaoNhanVat — nhung ket noi toi Postgres da
//     duoc "lam nong" tu truoc (store.Open goi Ping() khi mo, va cac tai
//     khoan thu cung da chen qua mot pool khac truoc do), nen moi truy van
//     tiep theo tren localhost thuong chi mat vai mili-giay.
//
// 1 giay la du du (hang tram lan lon hon do tre thuc te ky vong) de khong
// bao gio nham "chua toi luot goi" thanh "se khong bao gio co goi", ma van
// nhanh hon nhieu so voi doi nguyen 5 giay o MOI buoc khang dinh — file nay
// co hang chuc buoc nhu vay, cong don se rat cham neu giu nguyen 5s.
const hanChoImLang = 1 * time.Second

// clientGia la mot lop mong quanh MOT phien XuLy that (net.Pipe +
// mynet.PhucVu that), dung dung buoc bat tay nhu dungPhienThat trong
// login_test.go — nhung KHAC dungPhienThat o cho no GIU LAI net.Conn phia
// client, de co the tu dat lai han gio doc NGAN truoc moi lan khang dinh
// "im lang" (xem hanChoImLang) thay vi dung han 5 giay co dinh tu luc tao
// phien. Day chinh la "cach gon hon" ma task-7-brief cho phep dung khi van
// giu duoc do chac chan (dat han gio ngan tren KET NOI, khong bo qua buoc
// cho).
type clientGia struct {
	t    *testing.T
	f    *protocol.Framer
	conn net.Conn
}

// moPhienClientGia dung mot phien XuLy that tren kho da cho (o file nay luon
// la khoThat, *store.Store that) va thuc hien bat tay y het dungPhienThat.
func moPhienClientGia(t *testing.T, kho Kho) *clientGia {
	t.Helper()
	b := NewBoXuLy(kho)

	benClient, benServer := net.Pipe()
	go mynet.PhucVu(benServer, func() mynet.XuLyGoiTin { return b })

	benClient.SetReadDeadline(time.Now().Add(5 * time.Second))

	f := protocol.NewFramer(benClient, benClient)
	if err := f.WriteMessage(&protocol.Message{Command: -27}); err != nil {
		t.Fatalf("gui goi bat tay that bai: %v", err)
	}
	khoaTho := docGoiKhoaThuCong(t, benClient)
	f.SetKey(protocol.DeriveKey(khoaTho))

	return &clientGia{t: t, f: f, conn: benClient}
}

func (cl *clientGia) dong() { cl.conn.Close() }

// dangNhap gui goi dang nhap dung dinh dang that (tai su dung dungGoiDangNhap
// cua login_test.go — cung mot ham dong goi cho ca hai file, khong lap lai).
func (cl *clientGia) dangNhap(ten, matKhau string) {
	cl.t.Helper()
	if err := cl.f.WriteMessage(dungGoiDangNhap(ten, matKhau)); err != nil {
		cl.t.Fatalf("gui goi dang nhap that bai: %v", err)
	}
}

// taoNhanVat gui goi -28/2 dung dinh dang that (khung lenhKhongTrenBanDo,
// than: byte lenh con lenhConTaoNhanVat, roi ten/lang/kieu toc).
func (cl *clientGia) taoNhanVat(ten string, lang, toc byte) {
	cl.t.Helper()
	w := protocol.NewWriter()
	w.Byte(lenhConTaoNhanVat)
	if err := w.UTF(ten); err != nil {
		cl.t.Fatalf("dong goi ten nhan vat that bai: %v", err)
	}
	w.Byte(lang)
	w.Byte(toc)
	if err := cl.f.WriteMessage(protocol.NewMessage(lenhKhongTrenBanDo, w)); err != nil {
		cl.t.Fatalf("gui goi tao nhan vat that bai: %v", err)
	}
}

// choLenh doc MOT goi tra ve va bao co dung lenh mong doi hay khong. Dung han
// gio HIEN CO tren ket noi (mac dinh 5s tu luc tao phien, hoac vua duoc dat
// lai ve 5s boi khongConGoiGi/nhanLenhTrongKhoang) — day la khang dinh DUONG
// hanh phuc (may chu tra loi gan nhu tuc thi qua mot pipe dong bo), nen trong
// thuc te doc rat nhanh; han dai chi la luoi an toan cho truong hop that bai.
func (cl *clientGia) choLenh(t *testing.T, lenh int8) bool {
	t.Helper()
	m, err := cl.f.ReadMessage()
	if err != nil {
		t.Logf("khong doc duoc goi tra ve (mong doi lenh %d): %v", lenh, err)
		return false
	}
	if m.Command != lenh {
		t.Logf("nhan lenh %d, mong doi %d", m.Command, lenh)
		return false
	}
	return true
}

// khongConGoiGi dat han gio NGAN (hanChoImLang) roi thu doc mot goi: bao true
// NEU VA CHI NEU het han ma khong co goi nao toi — tuc la thuc su im lang.
// Neu doc duoc mot goi (bat ky lenh nao) thi that bai ro rang (gui thua, dang
// nghi ngo dung ham nay sai cho); neu gap loi KHAC timeout (vd EOF vi ket noi
// bi dong/may chu sap) thi FAIL LUON bang t.Fatalf — khong duoc coi im lang vi
// "sap" khong phai "im lang", va day chinh la dieu kich ban 6 (chua dang nhap)
// phai loai tru ro rang.
func (cl *clientGia) khongConGoiGi(t *testing.T) bool {
	t.Helper()
	_ = cl.conn.SetReadDeadline(time.Now().Add(hanChoImLang))
	defer cl.conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	m, err := cl.f.ReadMessage()
	if err == nil {
		t.Logf("nhan duoc goi NGOAI DU KIEN trong luc cho im lang: lenh=%d du_lieu=%v", m.Command, m.Data)
		return false
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true // het han ma khong co goi nao them -> thuc su im lang
	}
	t.Fatalf("ket noi gap loi KHONG PHAI het han gio khi dang cho xac nhan im lang (nghi ngo may chu sap hoac tu dong ket noi): %v", err)
	return false
}

// nhanLenhTrongKhoang dat han gio NGAN (hanChoImLang) roi doc TOI KHI het han,
// bao true neu THAY duoc mot goi mang dung lenh da cho trong khoang do. Dung
// cho khang dinh am: "lenh X KHONG duoc xuat hien". Cung loai tru "sap" khoi
// "im lang" giong het khongConGoiGi.
func (cl *clientGia) nhanLenhTrongKhoang(t *testing.T, lenh int8) bool {
	t.Helper()
	_ = cl.conn.SetReadDeadline(time.Now().Add(hanChoImLang))
	defer cl.conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	for {
		m, err := cl.f.ReadMessage()
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				return false // het han, khong thay lenh can tim
			}
			t.Fatalf("ket noi gap loi KHONG PHAI het han gio khi dang cho xac nhan lenh %d (nghi ngo may chu sap hoac tu dong ket noi): %v", lenh, err)
		}
		if m.Command == lenh {
			return true
		}
	}
}

// choDungMotLenh doc DUNG MOT goi tra ve, khang dinh dung lenh mong doi, ROI
// khang dinh KHONG con goi nao them (dung ham nay cho moi khang dinh "mot loi
// tu choi -26 duy nhat" — cung tinh than voi kiemTraDungMotGoiTuChoi trong
// nhanvat_test.go, o day o muc dau-cuoi qua XuLy that + CSDL that thay vi goi
// thang xuLyTaoNhanVat).
func (cl *clientGia) choDungMotLenh(t *testing.T, lenh int8) bool {
	t.Helper()
	if !cl.choLenh(t, lenh) {
		return false
	}
	return cl.khongConGoiGi(t)
}

// ---------------------------------------------------------------------------
// Kich ban 1 (duong hanh phuc): dang ky -> dang nhap -> nhan lenh 2 -> tao
// nhan vat hop le -> khong -26 -> kiem CSDL -> ngat -> phien MOI dang nhap
// lai -> KHONG con nhan lenh 2 nua.
// ---------------------------------------------------------------------------

func TestTaoRoiDangNhapLai_ThayNhanVat(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	tenTK := "nguoi" + hauToNganNgau(t)
	matKhau := "matkhau123"
	taiKhoanID := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)

	tenNV := "Naru" + hauToNganNgau(t)

	// Lan 1: dang nhap -> chua co nhan vat -> DUNG MOT lenh 2 -> tao nhan vat
	// hop le -> hoan toan im lang (thanh cong khong tra loi goi nao).
	{
		cl := moPhienClientGia(t, khoThat)
		defer cl.dong()

		cl.dangNhap(tenTK, matKhau)
		if !cl.choDungMotLenh(t, 2) {
			t.Fatal("lan dau dang nhap phai nhan DUNG MOT lenh 2 (man tao nhan vat)")
		}

		cl.taoNhanVat(tenNV, 0, 5)
		if !cl.khongConGoiGi(t) {
			t.Fatal("tao nhan vat hop le khong duoc tra ve goi nao (dac biet la -26)")
		}
	}

	// Kiem CSDL: dung hang nhan vat, dung ten, dung lang, dung mau/chakra
	// khoi dau THEO LANG (lay TimLang(0) lam gia tri doi chieu — dam bao ca
	// duong ong xuLyTaoNhanVat -> store.TaoNhanVat khong lam sai lech du lieu
	// giua buoc tra cuu lang va buoc ghi CSDL).
	nv, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoanID)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan sau khi tao that bai: %v", err)
	}
	langMongDoi := TimLang(0)
	if langMongDoi == nil {
		t.Fatal("TimLang(0) tra ve nil - lang 0 phai luon ton tai")
	}
	if nv.Ten != tenNV {
		t.Errorf("ten nhan vat = %q, mong doi %q", nv.Ten, tenNV)
	}
	if nv.Lang != 0 {
		t.Errorf("lang = %d, mong doi 0", nv.Lang)
	}
	if nv.HpToiDa != langMongDoi.HpKhoiDau || nv.HpHienTai != langMongDoi.HpKhoiDau {
		t.Errorf("hp toi da/hien tai = %d/%d, mong doi ca hai deu = %d (HpKhoiDau cua lang 0)",
			nv.HpToiDa, nv.HpHienTai, langMongDoi.HpKhoiDau)
	}
	if nv.ChakraToiDa != langMongDoi.ChakraKhoiDau || nv.ChakraHienTai != langMongDoi.ChakraKhoiDau {
		t.Errorf("chakra toi da/hien tai = %d/%d, mong doi ca hai deu = %d (ChakraKhoiDau cua lang 0)",
			nv.ChakraToiDa, nv.ChakraHienTai, langMongDoi.ChakraKhoiDau)
	}
	if nv.Cap != 1 {
		t.Errorf("cap = %d, mong doi 1 (nhan vat vua tao)", nv.Cap)
	}

	// Lan 2: NGAT phien cu (defer cl.dong() o tren da chay), mo PHIEN MOI,
	// dang nhap lai CUNG tai khoan -> da co nhan vat -> KHONG duoc nhan lenh
	// 2 nua (hoan toan im lang, vaoTheGioi stub Ke hoach 2 khong gui gi).
	{
		cl := moPhienClientGia(t, khoThat)
		defer cl.dong()
		cl.dangNhap(tenTK, matKhau)
		if !cl.khongConGoiGi(t) {
			t.Fatal("da co nhan vat -> dang nhap lai KHONG duoc nhan goi nao (dac biet la lenh 2)")
		}
	}
}

// ---------------------------------------------------------------------------
// Kich ban 2: hai tai khoan khac nhau, tai khoan thu hai tao nhan vat trung
// ten (chi khac hoa thuong) voi tai khoan thu nhat -> phai nhan -26.
// ---------------------------------------------------------------------------

func TestTichHop_TenTrungGiuaHaiTaiKhoan_TuChoi(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	matKhau := "matkhau123"

	tenTK1 := "trungA" + hauToNganNgau(t)
	tenTK2 := "trungB" + hauToNganNgau(t)
	taoTaiKhoanCoMatKhau(t, tenTK1, matKhau)
	taoTaiKhoanCoMatKhau(t, tenTK2, matKhau)

	tenNV := "Trung" + hauToNganNgau(t)

	// TK1 tao nhan vat truoc, thanh cong.
	{
		cl := moPhienClientGia(t, khoThat)
		defer cl.dong()
		cl.dangNhap(tenTK1, matKhau)
		if !cl.choDungMotLenh(t, 2) {
			t.Fatal("TK1: lan dau dang nhap phai nhan DUNG MOT lenh 2")
		}
		cl.taoNhanVat(tenNV, 0, 5)
		if !cl.khongConGoiGi(t) {
			t.Fatal("TK1: tao nhan vat hop le khong duoc tra ve goi nao")
		}
	}

	// TK2 tao nhan vat TRUNG TEN (chi khac hoa thuong voi TK1) -> phai bi tu
	// choi bang DUNG MOT goi -26.
	{
		cl := moPhienClientGia(t, khoThat)
		defer cl.dong()
		cl.dangNhap(tenTK2, matKhau)
		if !cl.choDungMotLenh(t, 2) {
			t.Fatal("TK2: lan dau dang nhap phai nhan DUNG MOT lenh 2")
		}
		tenTrung := strings.ToUpper(tenNV)
		cl.taoNhanVat(tenTrung, 0, 5)
		if !cl.choDungMotLenh(t, -26) {
			t.Fatalf("TK2: ten %q trung voi %q (chi khac hoa thuong) phai bi tu choi bang DUNG MOT goi -26", tenTrung, tenNV)
		}
	}
}

// ---------------------------------------------------------------------------
// Kich ban 3: mot tai khoan da co nhan vat (qua duong dang nhap lai, tuc
// vaoTheGioi da gan b.nhanVat tu CSDL), gui tiep -28/2 -> phai nhan -26.
// ---------------------------------------------------------------------------

func TestTichHop_TaiKhoanDaCoNhanVat_TuChoiTaoTiep(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	tenTK := "solo" + hauToNganNgau(t)
	matKhau := "matkhau123"
	taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "Naru" + hauToNganNgau(t)

	// Tao nhan vat lan dau, thanh cong.
	{
		cl := moPhienClientGia(t, khoThat)
		defer cl.dong()
		cl.dangNhap(tenTK, matKhau)
		if !cl.choDungMotLenh(t, 2) {
			t.Fatal("lan dau dang nhap phai nhan DUNG MOT lenh 2")
		}
		cl.taoNhanVat(tenNV, 0, 5)
		if !cl.khongConGoiGi(t) {
			t.Fatal("tao nhan vat lan dau hop le khong duoc tra ve goi nao")
		}
	}

	// Phien MOI, dang nhap lai: da co nhan vat -> hoan toan im lang (khong co
	// lenh 2) -> b.nhanVat da duoc gan tu CSDL qua vaoTheGioi. Van gui tiep
	// -28/2 (mo phong client hong/ke gia lap co tinh gui lai) -> phai bi tu
	// choi bang DUNG MOT goi -26 ("Ban da co nhan vat.").
	cl := moPhienClientGia(t, khoThat)
	defer cl.dong()

	cl.dangNhap(tenTK, matKhau)
	if !cl.khongConGoiGi(t) {
		t.Fatal("dang nhap lai voi tai khoan da co nhan vat -> khong duoc nhan goi nao (dac biet la lenh 2)")
	}

	cl.taoNhanVat("Khac"+hauToNganNgau(t), 0, 5)
	if !cl.choDungMotLenh(t, -26) {
		t.Fatal("tai khoan da co nhan vat, gui tiep -28/2 phai bi tu choi bang DUNG MOT goi -26")
	}
}

// ---------------------------------------------------------------------------
// Kich ban 4: ten khong hop le qua duong THAT (XuLy that + CSDL that, khac
// voi nhanvat_test.go von goi thang xuLyTaoNhanVat) — ten qua ngan va ten co
// ky tu la, moi cai phai nhan DUNG MOT -26.
// ---------------------------------------------------------------------------

func TestTichHop_TenKhongHopLeQuaDuongThat_TuChoi(t *testing.T) {
	khoThat := moKhoThatChoTest(t)

	truongHop := []struct {
		mota string
		ten  string
	}{
		{"qua ngan (3 ky tu)", "abc"},
		{"ky tu la (khong phai chu/so)", "Naru#123"},
	}

	for _, tc := range truongHop {
		t.Run(tc.mota, func(t *testing.T) {
			tenTK := "invtk" + hauToNganNgau(t)
			matKhau := "matkhau123"
			taoTaiKhoanCoMatKhau(t, tenTK, matKhau)

			cl := moPhienClientGia(t, khoThat)
			defer cl.dong()

			cl.dangNhap(tenTK, matKhau)
			if !cl.choDungMotLenh(t, 2) {
				t.Fatal("lan dau dang nhap phai nhan DUNG MOT lenh 2 (man tao nhan vat)")
			}

			cl.taoNhanVat(tc.ten, 0, 5)
			if !cl.choDungMotLenh(t, -26) {
				t.Fatalf("ten %q khong hop le qua duong that (XuLy that + CSDL that) phai bi tu choi bang DUNG MOT goi -26", tc.ten)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Kich ban 5: lang khong hop le (99, ngoai 0..4) -> phai nhan -26.
// ---------------------------------------------------------------------------

func TestTichHop_LangKhongHopLe_TuChoi(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	tenTK := "lang99" + hauToNganNgau(t)
	matKhau := "matkhau123"
	taoTaiKhoanCoMatKhau(t, tenTK, matKhau)

	cl := moPhienClientGia(t, khoThat)
	defer cl.dong()

	cl.dangNhap(tenTK, matKhau)
	if !cl.choDungMotLenh(t, 2) {
		t.Fatal("lan dau dang nhap phai nhan DUNG MOT lenh 2 (man tao nhan vat)")
	}

	cl.taoNhanVat("Naru"+hauToNganNgau(t), 99, 5) // lang 99 khong ton tai (chi co 0..4)
	if !cl.choDungMotLenh(t, -26) {
		t.Fatal("lang 99 khong hop le phai bi tu choi bang DUNG MOT goi -26")
	}
}

// ---------------------------------------------------------------------------
// Kich ban 6: gui -28/2 NGAY SAU bat tay, CHUA dang nhap -> phai IM LANG
// hoan toan (khong goi nao, khong sap ket noi).
// ---------------------------------------------------------------------------

func TestTichHop_ChuaDangNhapMaTaoNhanVat_ImLangKhongSap(t *testing.T) {
	khoThat := moKhoThatChoTest(t)

	cl := moPhienClientGia(t, khoThat)
	defer cl.dong()

	// Gui thang -28/2, CHUA gui goi dang nhap.
	cl.taoNhanVat("Naru"+hauToNganNgau(t), 0, 5)
	if !cl.khongConGoiGi(t) {
		t.Fatal("chua dang nhap ma gui tao nhan vat -> phai IM LANG hoan toan, khong goi tin nao")
	}

	// Khang dinh them "khong sap": ket noi VAN CON SONG va xu ly binh thuong
	// sau do — gui mot goi dang nhap HOP LE tren CUNG phien va phai nhan
	// duoc lenh 2 nhu binh thuong (khong bi may chu tu dong ket noi ngam).
	tenTK := "songsau" + hauToNganNgau(t)
	matKhau := "matkhau123"
	taoTaiKhoanCoMatKhau(t, tenTK, matKhau)

	cl.dangNhap(tenTK, matKhau)
	if !cl.choDungMotLenh(t, 2) {
		t.Fatal("sau khi gui goi tao nhan vat luc chua dang nhap, ket noi phai VAN hoat dong binh thuong (dang nhap tiep theo phai thanh cong)")
	}
}
