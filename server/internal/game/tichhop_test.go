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
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
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

// moPoolThuChoTest mo mot pgxpool RIENG (khac pool noi bo cua *store.Store,
// khong xuat ra ngoai) — dung cho cac thao tac SQL truc tiep ma store.Store
// co y khong xuat (chen tai khoan tho, dem hang kiem chung trong kich ban dua
// tranh). Tu dang ky dong qua t.Cleanup.
func moPoolThuChoTest(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("NRO_TEST_DSN")
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("khong mo duoc pool CSDL thu: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

// taoTaiKhoanCoMatKhau chen mot tai khoan THAT vao CSDL voi mat khau da BAM
// bang auth.Hash (dung dinh dang that server doi chieu luc dang nhap — luu
// mat khau tho se khong bao gio Verify() dung). Tra ve ID tai khoan vua tao.
//
// store.Store CO Y khong xuat ham "chen tai khoan tho" nao (tao tai khoan
// that su la qua kenh dang ky, ngoai pham vi Ke hoach 2) — nen o day phai tu
// mo mot pgxpool rieng (moPoolThuChoTest) va chen thang bang SQL, dung nhu
// ghi chu trong task-7-brief.md: "Trong game phai tu chen tai khoan qua SQL
// hoac qua kho that". Tu dang ky don dep qua t.Cleanup: xoa nhan_vat TRUOC
// (khoa ngoai toi tai_khoan) roi moi xoa tai_khoan — dung thu tu Cleanup chay
// LIFO nen moPoolThuChoTest dang ky pool.Close() TRUOC, ham nay dang ky xoa
// hang SAU, de luc chay cleanup thi xoa hang thuc hien truoc, dong pool sau
// cung.
func taoTaiKhoanCoMatKhau(t *testing.T, ten, matKhau string) int64 {
	t.Helper()
	pool := moPoolThuChoTest(t)

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
// Cap ham RIENG cho kich ban dua tranh (TestTaoNhanVat_DuaTranhCungTen): an
// toan goi tu GOROUTINE KHAC voi goroutine dang chay ham test.
//
// *testing.T KHONG an toan khi FailNow() (duoc t.Fatal/t.Fatalf goi ben
// trong) chay tu mot goroutine khac goroutine thuc thi ham Test — tai lieu Go
// noi ro "FailNow must be called from the goroutine running the Test..."; goi
// tu goroutine khac se chi lam GOROUTINE DO thoat qua runtime.Goexit(), rat
// de gay treo vinh vien (vd mot goroutine trong cuoc dua thoat som ma quen
// goi wg.Done()) thay vi that bai ro rang. Vi vay hai ham duoi day KHONG dung
// t o bat ky dang nao (khong ca t.Helper()) — chi tra du lieu (loi/ket qua)
// qua gia tri tra ve, de goroutine CHINH cua test tu gom qua channel roi moi
// goi t.Fatal SAU KHI wg.Wait() (xem TestTaoNhanVat_DuaTranhCungTen).
// ---------------------------------------------------------------------------

// guiTaoNhanVatKhongFatal giong het clientGia.taoNhanVat ve dinh dang goi tin
// (khung lenhKhongTrenBanDo, than: lenhConTaoNhanVat + ten + lang + kieu toc)
// nhung tra loi qua error thay vi goi t.Fatal — an toan goi tu goroutine con.
func (cl *clientGia) guiTaoNhanVatKhongFatal(ten string, lang, toc byte) error {
	w := protocol.NewWriter()
	w.Byte(lenhConTaoNhanVat)
	if err := w.UTF(ten); err != nil {
		return fmt.Errorf("dong goi ten nhan vat that bai: %w", err)
	}
	w.Byte(lang)
	w.Byte(toc)
	return cl.f.WriteMessage(protocol.NewMessage(lenhKhongTrenBanDo, w))
}

// ketQuaDocSauKhiTao la ket qua doc duoc tren MOT phien sau khi da gui goi
// tao nhan vat, dung rieng cho kich ban dua tranh. cauThongBao (noi dung cau
// bao trong goi -26) duoc giu lai CHI de ghi vao t.Logf lam bang chung/chan
// doan de neu sau nay dot bien khac lam sai NOI DUNG cau bao (khong doi
// hanh vi tu-choi/khong-tu-choi), dong log con lai van chi ro duoc — bat bien
// chinh cua test van CHI dua vao co/khong -26, khong dua vao dung chu.
type ketQuaDocSauKhiTao struct {
	biTuChoi    bool   // true neu nhan duoc -26 (hop thoai OK, tu choi)
	cauThongBao string // noi dung cau bao trong goi -26, neu biTuChoi
	loi         error  // khac nil neu gap loi bat thuong (khac het han gio doc)
}

// docPhanHoiTaoNhanVat dat han gio NGAN (hanChoImLang) roi doc DUNG MOT lan:
//   - het han ma khong co gi toi -> Ke hoach 2 khong gui goi nao khi THANH
//     CONG (xem vaoTheGioi trong nhanvat.go), nen day la thanh cong.
//   - nhan duoc -26 -> bi tu choi.
//   - bat ky truong hop nao khac (loi khac timeout, hoac lenh la khong mong
//     doi trong kich ban nay) -> tra ve qua truong loi, de goroutine chinh
//     lam that bai ro rang thay vi phan loai nham thanh cong/tu choi.
func (cl *clientGia) docPhanHoiTaoNhanVat() ketQuaDocSauKhiTao {
	_ = cl.conn.SetReadDeadline(time.Now().Add(hanChoImLang))
	defer cl.conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	m, err := cl.f.ReadMessage()
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return ketQuaDocSauKhiTao{biTuChoi: false} // im lang -> thanh cong
		}
		return ketQuaDocSauKhiTao{loi: fmt.Errorf("loi doc KHONG PHAI het han gio (nghi ngo may chu sap/dong ket noi): %w", err)}
	}
	if m.Command != lenhThongBaoOK {
		return ketQuaDocSauKhiTao{loi: fmt.Errorf("nhan lenh %d ngoai du kien, mong doi -26 hoac im lang", m.Command)}
	}
	dr := protocol.NewReader(m.Data)
	cau, _ := dr.UTF()
	return ketQuaDocSauKhiTao{biTuChoi: true, cauThongBao: cau}
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
	// khoi dau. MOI khang dinh chi so dung GIA TRI GHIM (150...) chu KHONG
	// doi chieu voi TimLang(0): duong ghi CSDL (xuLyTaoNhanVat) da goi
	// TimLang, nen neu phep so sanh cung goi TimLang thi ca hai phia doc CHUNG
	// mot nguon (danhSachLang trong lang.go) — mot dot bien sua thang so lieu
	// trong bang se lam CA HAI phia doi theo nhau va test khong bat duoc (da
	// kiem chung tren ban sao). Ghim gia tri HIEN TAI cua lang 0 (xem
	// danhSachLang) truc tiep; neu sau nay can bang lai chi so khoi dau
	// (lang.go ghi "se can bang lai o Ke hoach 4") thi sua cac gia tri ghim
	// duoi day cho khop.
	nv, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoanID)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan sau khi tao that bai: %v", err)
	}
	if nv.Ten != tenNV {
		t.Errorf("ten nhan vat = %q, mong doi %q", nv.Ten, tenNV)
	}
	if nv.Lang != 0 {
		t.Errorf("lang = %d, mong doi 0", nv.Lang)
	}
	if nv.HpToiDa != 150 || nv.HpHienTai != 150 {
		t.Errorf("hp toi da/hien tai = %d/%d, mong doi ca hai deu = 150 (HpKhoiDau ghim cua lang 0)",
			nv.HpToiDa, nv.HpHienTai)
	}
	if nv.ChakraToiDa != 150 || nv.ChakraHienTai != 150 {
		t.Errorf("chakra toi da/hien tai = %d/%d, mong doi ca hai deu = 150 (ChakraKhoiDau ghim cua lang 0)",
			nv.ChakraToiDa, nv.ChakraHienTai)
	}
	if nv.Cap != 1 {
		t.Errorf("cap = %d, mong doi 1 (nhan vat vua tao)", nv.Cap)
	}
	// Round-trip map/vi tri/kieu toc — cung dung GIA TRI GHIM nhu HP/Chakra o
	// tren, cung mot ly do (tranh diem mu "cung nguon" voi TimLang). Bat duoc
	// neu ai do ghi sai/hoan doi truong: vd gan nham ViTriY vao cot vi_tri_x,
	// hay cong nham 1 vao kieu toc truoc khi luu.
	if nv.MapID != 1 {
		t.Errorf("map id = %d, mong doi 1 (MapKhoiDau ghim cua lang 0, xem danhSachLang trong lang.go)", nv.MapID)
	}
	if nv.ViTriX != 100 {
		t.Errorf("vi tri x = %d, mong doi 100 (ViTriXKhoiDau ghim cua lang 0)", nv.ViTriX)
	}
	if nv.ViTriY != 384 {
		t.Errorf("vi tri y = %d, mong doi 384 (ViTriYKhoiDau ghim cua lang 0)", nv.ViTriY)
	}
	// 5 la byte kieu toc that su gui trong goi tao nhan vat o tren
	// (cl.taoNhanVat(tenNV, 0, 5)) — khong lien quan TimLang, khong co van de
	// "cung nguon" nhu tren.
	if nv.KieuToc != int16(5) {
		t.Errorf("kieu toc = %d, mong doi 5 (byte kieu toc da gui trong goi tao nhan vat)", nv.KieuToc)
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
//
// Luu y pham vi: hai buoc gui goi trong kich ban nay LUON tuan tu (TK1 tao
// xong HET, roi moi toi TK2), nen nhanh tu choi o day CHAC CHAN di qua
// pre-check tang game (TenNhanVatDaDung trong xuLyTaoNhanVat) — KHONG bao gio
// cham toi nhanh anh xa loi UNIQUE that su ben trong store.TaoNhanVat (so
// sanh pgErr.ConstraintName == "nhan_vat_ten_duy_nhat" -> ErrTenNhanVatDaDung
// trong store.go). Nhanh do chi kich hoat khi hai client "lot qua" pre-check
// CUNG LUC — xem TestTaoNhanVat_DuaTranhCungTen (duoi file nay) cho kich ban
// dua tranh that su phu duong do, va internal/store/store_test.go
// (TestTenNhanVatTrungBiTuChoi) cho kiem thu truc tiep tren store.Store.
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
//
// Luu y pham vi (giong het kich ban 2 o tren): -26 o day di qua nhanh
// pre-check b.nhanVat != nil ngay dau xuLyTaoNhanVat (khong cham CSDL nua),
// KHONG phai nhanh anh xa loi UNIQUE(tai_khoan_id) (nhan_vat_mot_tai_khoan)
// cua store.TaoNhanVat. Nhanh UNIQUE(tai_khoan_id) do rieng — mot tai khoan
// gui HAI yeu cau tao nhan vat DUA NHAU (chua co lan nao thanh cong truoc do)
// — ngoai pham vi kich ban nay (kich ban nay co tinh tao THANH CONG lan dau
// roi moi thu lai lan hai, khong phai dua tranh); duoc
// internal/store/store_test.go (TestMotTaiKhoanChiMotNhanVat) phu truc tiep.
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

// ---------------------------------------------------------------------------
// Kich ban 7 (sua vong 1, theo phat hien cua nguoi soi chay dot bien): DUA
// TRANH THAT giua hai tai khoan cung tao nhan vat MOT TEN gan nhu dong thoi.
//
// Boi canh: kich ban 2 va 3 o tren CHI di qua pre-check tang game
// (TenNhanVatDaDung / b.nhanVat != nil) vi hai buoc gui goi trong hai kich
// ban do LUON TUAN TU — nen KHONG BAO GIO cham toi nhanh anh xa loi UNIQUE
// that su ben trong store.TaoNhanVat (store.go: so sanh
// pgErr.ConstraintName == "nhan_vat_ten_duy_nhat" -> ErrTenNhanVatDaDung).
// Nhanh do CHI kich hoat khi hai client "lot qua" pre-check CUNG LUC — dung
// xay ra tren may chu nhieu nguoi choi that (hai nguoi bam nut "Tao" gan nhu
// dong thoi voi cung mot ten dep). Test nay dung dung tinh huong do: hai
// goroutine gui goi tao nhan vat GAN NHU DONG THOI, khong doi nhau — dua vao
// pgxpool (an toan dung dong thoi) cua CUNG mot khoThat va rang buoc UNIQUE
// that cua Postgres de tao ra cuoc dua thuc su, khong gia lap.
func TestTaoNhanVat_DuaTranhCungTen(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	tenTK1 := "duaA" + hauToNganNgau(t)
	tenTK2 := "duaB" + hauToNganNgau(t)
	matKhau := "matkhau123"
	taoTaiKhoanCoMatKhau(t, tenTK1, matKhau)
	taoTaiKhoanCoMatKhau(t, tenTK2, matKhau)

	// "Race" (4 chu cai) + hau to: dam bao >=5 ky tu KE CA khi hau to chi la
	// 1 chu so (vd "Race3" = 5 ky tu) — mot tien to ngan hon (vd "Dua", 3 chu
	// cai) se cho ten 4 ky tu voi hau to 1 chu so, VI PHAM ngay reTenNhanVat
	// (^[A-Za-z0-9]{5,10}$) truoc ca khi toi duoc logic dua tranh can kiem.
	tenNV := "Race" + hauToNganNgau(t)

	cl1 := moPhienClientGia(t, khoThat)
	defer cl1.dong()
	cl2 := moPhienClientGia(t, khoThat)
	defer cl2.dong()

	cl1.dangNhap(tenTK1, matKhau)
	if !cl1.choDungMotLenh(t, 2) {
		t.Fatal("TK1: lan dau dang nhap phai nhan DUNG MOT lenh 2")
	}
	cl2.dangNhap(tenTK2, matKhau)
	if !cl2.choDungMotLenh(t, 2) {
		t.Fatal("TK2: lan dau dang nhap phai nhan DUNG MOT lenh 2")
	}

	// Cuoc dua: gui goi tao nhan vat CUNG TEN tren hai goroutine RIENG, khong
	// dong bo hoa gi giua chung — day chinh la dieu tao ra dua tranh thuc su
	// giua hai luong xu ly phia may chu (hai *BoXuLy rieng, dung chung
	// khoThat/pgxpool an toan dong thoi). Ket qua gom qua channel (KHONG ghi
	// chung bien tran giua cac goroutine) — xem chu thich cap ham
	// guiTaoNhanVatKhongFatal/docPhanHoiTaoNhanVat o tren ve ly do khong dung
	// t.Fatal ben trong goroutine con.
	type ketQuaDua struct {
		phien       string
		biTuChoi    bool
		cauThongBao string
		loi         error
	}
	ketQuaCh := make(chan ketQuaDua, 2)
	var wg sync.WaitGroup
	guiVaDoc := func(nhan string, cl *clientGia) {
		defer wg.Done()
		if err := cl.guiTaoNhanVatKhongFatal(tenNV, 0, 5); err != nil {
			ketQuaCh <- ketQuaDua{phien: nhan, loi: fmt.Errorf("gui goi tao nhan vat that bai: %w", err)}
			return
		}
		kq := cl.docPhanHoiTaoNhanVat()
		ketQuaCh <- ketQuaDua{phien: nhan, biTuChoi: kq.biTuChoi, cauThongBao: kq.cauThongBao, loi: kq.loi}
	}
	wg.Add(2)
	go guiVaDoc("TK1", cl1)
	go guiVaDoc("TK2", cl2)
	wg.Wait()
	close(ketQuaCh)

	// Tu day tro xuong CHAY TREN GOROUTINE CHINH cua test — an toan dung
	// t.Fatal/t.Errorf binh thuong.
	var soThanhCong, soBiTuChoi int
	for kq := range ketQuaCh {
		if kq.loi != nil {
			t.Fatalf("phien %s gap loi khi gui/doc phan hoi tao nhan vat: %v", kq.phien, kq.loi)
		}
		if kq.biTuChoi {
			soBiTuChoi++
			t.Logf("phien %s: bi tu choi -26, cau bao %q", kq.phien, kq.cauThongBao)
		} else {
			soThanhCong++
			t.Logf("phien %s: thanh cong (im lang)", kq.phien)
		}
	}

	// Bat bien phai dung BAT KE tang nao thang trong cuoc dua (pre-check kip
	// chan mot ben, HOAC ca hai lot qua pre-check roi UNIQUE cua Postgres chot
	// chan): DUNG MOT thanh cong, DUNG MOT bi tu choi. Khong duoc ca hai thanh
	// cong (vi pham UNIQUE(lower(ten))), khong duoc ca hai bi tu choi (ten hop
	// le va chua ai dung truoc do, nen it nhat mot ben PHAI thanh cong).
	if soThanhCong != 1 || soBiTuChoi != 1 {
		t.Fatalf("dua tranh cung ten %q: mong doi DUNG MOT thanh cong + DUNG MOT bi tu choi (-26), nhan duoc %d thanh cong / %d bi tu choi",
			tenNV, soThanhCong, soBiTuChoi)
	}

	// Kiem CSDL doc lap voi viec dem -26 o tren: phai co DUNG MOT hang nhan_vat
	// mang ten nay (khong phan biet hoa thuong) — khang dinh cuoi cung, chac
	// chan khong co hai hang trung ten nao lot duoc ca hai tang kiem tra.
	pool := moPoolThuChoTest(t)
	var soHangTrungTen int
	if err := pool.QueryRow(ctx,
		`SELECT count(*) FROM nhan_vat WHERE lower(ten) = lower($1)`, tenNV,
	).Scan(&soHangTrungTen); err != nil {
		t.Fatalf("dem nhan_vat trung ten that bai: %v", err)
	}
	if soHangTrungTen != 1 {
		t.Fatalf("CSDL co %d hang nhan_vat ten %q (khong phan biet hoa thuong), mong doi DUNG MOT", soHangTrungTen, tenNV)
	}
}
