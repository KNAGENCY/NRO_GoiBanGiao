package game

import (
	"context"
	"errors"
	"testing"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// khoHongHanhTrang la khoGia (login_test.go) NHUNG LayTui luon hong — dung de
// kiem ky luat "loi CSDL thi KHONG cho vao the gioi nua voi".
type khoHongHanhTrang struct {
	khoGia
}

var loiHanhTrangGia = errors.New("CSDL hong (gia lap)")

func (k *khoHongHanhTrang) LayTui(ctx context.Context, nvID int64, noiChua string) ([]store.VatPham, error) {
	return nil, loiHanhTrangGia
}

// TestVaoTheGioiLoiHanhTrangKhongGuiGiCa: LayTui hong => may chu KHONG gui MOT
// BYTE NAO (khong -31, khong -82, khong -30) va dong phien.
//
// DOT BIEN: nuot loi cua hanhTrangNap roi cho vao the gioi voi tui rong ->
// test nay do vi se nhan duoc -31.
func TestVaoTheGioiLoiHanhTrangKhongGuiGiCa(t *testing.T) {
	b := NewBoXuLy(&khoHongHanhTrang{}, nil)
	b.nhanVat = &store.NhanVat{ID: 7, Ten: "Naruto", Lang: 0, KieuToc: 5,
		SucManh: 100, Vang: 1000, Ngoc: 0, MapID: 1, ViTriX: 100, ViTriY: 384}
	fClient, dong := dungPhienThat(t, b)
	defer dong()

	w := protocol.NewWriter()
	w.Byte(lenhConClientOk)
	if err := fClient.WriteMessage(protocol.NewMessage(lenhKhongTrenBanDo, w)); err != nil {
		t.Fatalf("gui -28/13 that bai: %v", err)
	}

	if m, err := fClient.ReadMessage(); err == nil {
		t.Fatalf("nhan duoc goi lenh %d du hanh trang KHONG nap duoc — phai khong gui gi ca va dong phien",
			m.Command)
	}
	if b.hanhTrang != nil {
		t.Error("b.hanhTrang khac nil du LayTui hong — khong duoc de lai bang tra nua voi")
	}
}

// ---------------------------------------------------------------------------
// KH7 Task 5 — ME_INFO (-30/0) mang BA MANG VAT PHAM that.
//
// PHA A: docs/superpowers/research/2026-09-03-kh7-me-info-so-o.md
// Su that phai giu: con so xuong client la SO O (suc chua), khong phai so mon;
// va sau moi con so phai co DU so o than o (Controller.cs:5860-5968).
// ---------------------------------------------------------------------------

// doDenBaMangVatPham doc ME_INFO tu dau toi ngay TRUOC itemsBody.size, roi tra
// con tro doc de ba khang dinh o duoi tu doc tiep. Giu DUY NHAT o day cong
// thuc "phan dau ME_INFO dai bao nhieu" — doi khuon phan dau chi phai sua mot
// cho.
func doDenBaMangVatPham(t *testing.T, nv *store.NhanVat, m *protocol.Message) *protocol.Reader {
	t.Helper()
	r := protocol.NewReader(m.Data)
	boByte := func(n int, ten string) {
		t.Helper()
		for i := 0; i < n; i++ {
			if _, err := r.Byte(); err != nil {
				t.Fatalf("doc %s that bai: %v", ten, err)
			}
		}
	}
	boByte(1, "byte lenh con")
	boByte(4, "id")
	boByte(1, "taskMain.id")
	boByte(1, "gender")
	boByte(2, "head")
	if _, err := r.UTF(); err != nil {
		t.Fatalf("doc name that bai: %v", err)
	}
	boByte(1, "cPK")
	boByte(1, "typePk")
	boByte(8, "power")
	boByte(2, "buffHp")
	boByte(2, "buffMp")
	boByte(1, "gender lan 2")
	boByte(1, "skillCount")
	boByte(8, "gold")
	boByte(4, "ruby")
	boByte(4, "gem")
	return r
}

// docMotBang doc [soO: byte] roi DU soO o, tra ve so o va mau cua tung o
// (-1 = o rong). FAIL neu thieu byte — day chinh la loi ma dot bien "ghi so o
// nhung khong ghi than o" gay ra.
func docMotBang(t *testing.T, r *protocol.Reader, ten string) (int, []int16) {
	t.Helper()
	soO, err := r.Byte()
	if err != nil {
		t.Fatalf("doc %s.size that bai: %v", ten, err)
	}
	mau := make([]int16, 0, int(soO))
	for i := 0; i < int(soO); i++ {
		id, err := r.Short()
		if err != nil {
			t.Fatalf("%s: doc mau o %d/%d that bai: %v (goi THIEU BYTE — so o da ghi nhung than o thi khong)",
				ten, i, soO, err)
		}
		if id == -1 {
			mau = append(mau, -1)
			continue
		}
		// O co mon: bo qua phan than con lai cua o do.
		if _, err := r.Int(); err != nil {
			t.Fatalf("%s o %d: doc so luong that bai: %v", ten, i, err)
		}
		if _, err := r.UTF(); err != nil {
			t.Fatalf("%s o %d: doc info that bai: %v", ten, i, err)
		}
		if _, err := r.UTF(); err != nil {
			t.Fatalf("%s o %d: doc content that bai: %v", ten, i, err)
		}
		soOpt, err := r.Byte()
		if err != nil {
			t.Fatalf("%s o %d: doc so option that bai: %v", ten, i, err)
		}
		for j := 0; j < int(soOpt); j++ {
			if _, err := r.Byte(); err != nil {
				t.Fatalf("%s o %d: doc optionId %d that bai: %v", ten, i, j, err)
			}
			if _, err := r.Short(); err != nil {
				t.Fatalf("%s o %d: doc param %d that bai: %v", ten, i, j, err)
			}
		}
		mau = append(mau, id)
	}
	return int(soO), mau
}

func nhanVatThuMeInfo() *store.NhanVat {
	return &store.NhanVat{
		ID: 777, Ten: "Sasuke", Lang: 0, KieuToc: 3,
		SucManh: 1000, Vang: 10, Ngoc: 0, MapID: 1,
	}
}

// TestMeInfoGhiSoOKhongPhaiSoMon la khang dinh TRUNG TAM cua task 5.
//
// DOT BIEN 1 (ghi so MON thay so O): hai mon trong tui -> itemsBag.size = 2
// thay vi 20. Test do vi so o sai VA vi o 1 khong con la o rong -1.
func TestMeInfoGhiSoOKhongPhaiSoMon(t *testing.T) {
	chu := int64(777)
	o0, o2 := int16(0), int16(2)
	dsMon := []store.VatPham{
		{ID: 11, MauID: 0, ChuSoHuuID: &chu, NoiChua: store.NoiChuaTui, OSo: &o0, SoLuong: 5},
		{ID: 12, MauID: 1, ChuSoHuuID: &chu, NoiChua: store.NoiChuaTui, OSo: &o2, SoLuong: 1},
	}

	m, err := goiMeInfo(nhanVatThuMeInfo(), dsMon)
	if err != nil {
		t.Fatalf("goiMeInfo loi: %v", err)
	}
	r := doDenBaMangVatPham(t, nhanVatThuMeInfo(), m)

	soOMac, mauMac := docMotBang(t, r, "itemsBody")
	if soOMac != SoOMacChuan {
		t.Fatalf("itemsBody.size = %d, mong doi %d (SO O, khong phai so mon)", soOMac, SoOMacChuan)
	}
	for i, id := range mauMac {
		if id != -1 {
			t.Fatalf("itemsBody o %d = mau %d, mong doi -1 (khong mac gi)", i, id)
		}
	}

	soOTui, mauTui := docMotBang(t, r, "itemsBag")
	if soOTui != SoOTuiChuan {
		t.Fatalf("itemsBag.size = %d, mong doi %d (SO O, khong phai so mon: co 2 mon nhung 20 o)",
			soOTui, SoOTuiChuan)
	}
	if mauTui[0] != 0 {
		t.Errorf("itemsBag o 0 = mau %d, mong doi 0", mauTui[0])
	}
	if mauTui[1] != -1 {
		t.Errorf("itemsBag o 1 = mau %d, mong doi -1 (o RONG giua hai mon — day la cho dot bien 'so mon' lam sai)", mauTui[1])
	}
	if mauTui[2] != 1 {
		t.Errorf("itemsBag o 2 = mau %d, mong doi 1", mauTui[2])
	}
	for i := 3; i < soOTui; i++ {
		if mauTui[i] != -1 {
			t.Errorf("itemsBag o %d = mau %d, mong doi -1", i, mauTui[i])
		}
	}

	soORuong, mauRuong := docMotBang(t, r, "itemsBox")
	if soORuong != SoORuongChuan {
		t.Fatalf("itemsBox.size = %d, mong doi %d", soORuong, SoORuongChuan)
	}
	for i, id := range mauRuong {
		if id != -1 {
			t.Errorf("itemsBox o %d = mau %d, mong doi -1", i, id)
		}
	}

	// Phan duoi ME_INFO PHAI van doc duoc — neu ba mang ghi thieu/thua byte thi
	// cap anh xa headAvatar duoi day se ra so rac.
	if n, err := r.Short(); err != nil || n != 1 {
		t.Fatalf("headAvatar.size = %d (loi %v), mong doi 1 — ba mang vat pham ghi LECH byte", n, err)
	}
	if h, err := r.Short(); err != nil || h != HeadNhanVat {
		t.Fatalf("headAvatar[0].headId = %d (loi %v), mong doi %d", h, err, HeadNhanVat)
	}
}

// TestMeInfoItemsBodyDocKhongDau: client doc itemsBody bang readUnsignedByte
// (Controller.cs:5860) nen 200 o la HOP LE o rieng mang nay.
//
// DOT BIEN 2 (ghi itemsBody CO DAU / kep ve 127): con so doc ra khong con la
// 200, hoac ham tu choi — ca hai deu lam test nay do.
func TestMeInfoItemsBodyDocKhongDau(t *testing.T) {
	const soOLon = 200
	if soOLon <= 127 {
		t.Fatal("ca thu phai > 127 moi phan biet duoc co dau / khong dau")
	}
	oMac := make([]*store.VatPham, soOLon)

	m, err := goiMeInfoDaXep(nhanVatThuMeInfo(), oMac, nil, nil)
	if err != nil {
		t.Fatalf("goiMeInfoDaXep loi: %v", err)
	}
	r := doDenBaMangVatPham(t, nhanVatThuMeInfo(), m)

	soO, mau := docMotBang(t, r, "itemsBody")
	if soO != soOLon {
		t.Fatalf("itemsBody.size doc ra = %d, mong doi %d (client doc UNSIGNED byte, :5860)", soO, soOLon)
	}
	if len(mau) != soOLon {
		t.Fatalf("doc duoc %d o, mong doi %d", len(mau), soOLon)
	}
}

// TestSoOChuanNamTrongTran: tui va ruong xuong client bang byte CO DAU
// (Controller.cs:5907/:5942) va ngoai le tran o hai cho do KHONG nam trong
// try/catch -> nguoi choi ket man "Xin cho". Ba hang so phai o duoi tran.
func TestSoOChuanNamTrongTran(t *testing.T) {
	if SoOTuiChuan < 1 || SoOTuiChuan > tranSoO {
		t.Errorf("SoOTuiChuan = %d, phai trong [1,%d]", SoOTuiChuan, tranSoO)
	}
	if SoORuongChuan < 1 || SoORuongChuan > tranSoO {
		t.Errorf("SoORuongChuan = %d, phai trong [1,%d]", SoORuongChuan, tranSoO)
	}
	if SoOMacChuan < 1 || SoOMacChuan > 255 {
		t.Errorf("SoOMacChuan = %d, phai trong [1,255]", SoOMacChuan)
	}
}

// TestMeInfoLoiThiKhongGuiGiCa: mon tro toi mau khong co trong danh muc ->
// goiMeInfo tra loi va KHONG tra goi (§1.10 "khong gui goi thieu").
func TestMeInfoLoiThiKhongGuiGiCa(t *testing.T) {
	chu := int64(777)
	o0 := int16(0)
	dsMon := []store.VatPham{
		{ID: 11, MauID: 30000, ChuSoHuuID: &chu, NoiChua: store.NoiChuaTui, OSo: &o0, SoLuong: 1},
	}
	m, err := goiMeInfo(nhanVatThuMeInfo(), dsMon)
	if err == nil {
		t.Fatal("mon tro toi mau khong ton tai ma goiMeInfo van tra goi — phai tra loi")
	}
	if m != nil {
		t.Fatal("goiMeInfo tra loi NHUNG van tra mot goi — phai tra nil")
	}
}
