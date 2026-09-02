package game

import (
	"bytes"
	"testing"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// ---------------------------------------------------------------------------
// Test BYTE-CHINH-XAC cho ba goi -36 / -37 / -35.
//
// Moi day byte mong doi duoc GO TAY tu bang §3 cua
// docs/superpowers/research/2026-09-03-kh7-giai-ma-36-37-35.md — KHONG sinh ra
// tu chinh ham dang kiem. Neu ai do sua ham dung goi thi cac day nay phai duoc
// sua bang tay, va do la luc phai doc lai client.
// ---------------------------------------------------------------------------

// mon dung mot dong vat pham toi thieu de gui xuong.
func mon(id int64, mauID int32, oSo int16, noiChua string, soLuong int32, opt ...store.TuyChon) store.VatPham {
	chu := int64(7)
	o := oSo
	return store.VatPham{
		ID:         id,
		MauID:      mauID,
		ChuSoHuuID: &chu,
		NoiChua:    noiChua,
		OSo:        &o,
		SoLuong:    soLuong,
		ChiSoPhu:   store.ChiSoPhu{Opt: opt},
	}
}

// utfRong la day byte cua writeUTF("") — 2 byte do dai 0. Ghim doc lap.
var utfRong = []byte{0x00, 0x00}

// oRong la day byte cua MOT o rong: dung mot short -1.
var oRong = []byte{0xFF, 0xFF}

// ---------------------------------------------------------------------------
// 1. Tui rong hoan toan
// ---------------------------------------------------------------------------

func TestGoiTuiRongHoanToan(t *testing.T) {
	m, err := GoiTui(nil, 5)
	if err != nil {
		t.Fatalf("GoiTui loi: %v", err)
	}
	if m.Command != -36 {
		t.Fatalf("lenh = %d, mong doi -36 (Controller.cs:2202)", m.Command)
	}
	want := []byte{
		0x00,       // cAction = 0
		0x05,       // soO = 5
		0xFF, 0xFF, // o 0 rong
		0xFF, 0xFF, // o 1 rong
		0xFF, 0xFF, // o 2 rong
		0xFF, 0xFF, // o 3 rong
		0xFF, 0xFF, // o 4 rong
	}
	if !bytes.Equal(m.Data, want) {
		t.Errorf("than -36 = % X,\nmong doi        % X", m.Data, want)
	}
}

// ---------------------------------------------------------------------------
// 2. O RONG O GIUA — day la test bat loi L1
// ---------------------------------------------------------------------------

// TestGoiTuiORongOGiua: o 0 co mon, o 1 TRONG, o 2 co mon.
//
// Ban Java bo qua o 1 mà khong ghi gi (InventoryService.java:570-572) nen goi
// thieu dung 2 byte; client van lap du 3 lan va doc lan sang phan khac. Bo dong
// `w.Short(-1)` trong ghiCacO se lam test nay do.
func TestGoiTuiORongOGiua(t *testing.T) {
	ds := []store.VatPham{
		mon(1, 0, 0, store.NoiChuaTui, 1),
		mon(2, 1, 2, store.NoiChuaTui, 1),
	}
	m, err := GoiTui(ds, 3)
	if err != nil {
		t.Fatalf("GoiTui loi: %v", err)
	}

	var want []byte
	want = append(want, 0x00, 0x03) // cAction=0, soO=3
	// o 0: mau 0, so luong 1, info "", content "", 0 option
	want = append(want, 0x00, 0x00) // templateId = 0
	want = append(want, 0x00, 0x00, 0x00, 0x01)
	want = append(want, utfRong...)
	want = append(want, utfRong...)
	want = append(want, 0x00) // soOption
	// o 1: RONG
	want = append(want, oRong...)
	// o 2: mau 1, so luong 1
	want = append(want, 0x00, 0x01)
	want = append(want, 0x00, 0x00, 0x00, 0x01)
	want = append(want, utfRong...)
	want = append(want, utfRong...)
	want = append(want, 0x00)

	if !bytes.Equal(m.Data, want) {
		t.Errorf("than -36 (o rong o giua) = % X,\nmong doi                     % X", m.Data, want)
	}
	// Khang dinh do dai rieng: neu ai do bo dau -1 thi goi ngan di dung 2 byte.
	if len(m.Data) != len(want) {
		t.Errorf("do dai than = %d byte, mong doi %d — thieu byte la client doc lan, im lang",
			len(m.Data), len(want))
	}
}

// ---------------------------------------------------------------------------
// 3. Mon cong don so_luong > 1
// ---------------------------------------------------------------------------

func TestGoiTuiMonCongDon(t *testing.T) {
	ds := []store.VatPham{mon(1, 0, 0, store.NoiChuaTui, 99)}
	m, err := GoiTui(ds, 1)
	if err != nil {
		t.Fatalf("GoiTui loi: %v", err)
	}
	want := []byte{
		0x00, 0x01, // cAction, soO
		0x00, 0x00, // templateId = 0
		0x00, 0x00, 0x00, 0x63, // so luong = 99, INT 4 byte
		0x00, 0x00, // info ""
		0x00, 0x00, // content ""
		0x00, // 0 option
	}
	if !bytes.Equal(m.Data, want) {
		t.Errorf("than = % X,\nmong doi % X", m.Data, want)
	}
}

// ---------------------------------------------------------------------------
// 4. Mon co HAI option
// ---------------------------------------------------------------------------

func TestGoiTuiMonHaiOption(t *testing.T) {
	ds := []store.VatPham{mon(1, 2, 0, store.NoiChuaTui, 1,
		store.TuyChon{Id: 2, Param: 5},
		store.TuyChon{Id: 3, Param: 65535},
	)}
	m, err := GoiTui(ds, 1)
	if err != nil {
		t.Fatalf("GoiTui loi: %v", err)
	}
	want := []byte{
		0x00, 0x01,
		0x00, 0x02, // templateId = 2
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00,
		0x00, 0x00,
		0x02,       // 2 option
		0x02,       // option id 2 — MOT byte
		0x00, 0x05, // param 5 — HAI byte khong dau
		0x03,       // option id 3
		0xFF, 0xFF, // param 65535
	}
	if !bytes.Equal(m.Data, want) {
		t.Errorf("than = % X,\nmong doi % X", m.Data, want)
	}
}

// ---------------------------------------------------------------------------
// 5. Doc NGUOC theo dung trinh tu lenh doc cua client
// ---------------------------------------------------------------------------

// Test byte-chinh-xac phia tren bat duoc "byte khac"; test nay bat duoc "byte
// giong nhung Y NGHIA lech" — vi du dao cho info voi content (deu la UTF), hay
// ghi so luong bang short thay vi int.
func TestGoiTuiDocLaiDungTrinhTuClient(t *testing.T) {
	ds := []store.VatPham{
		mon(1, 0, 0, store.NoiChuaTui, 12, store.TuyChon{Id: 1, Param: 300}),
		mon(2, 4, 2, store.NoiChuaTui, 1),
	}
	m, err := GoiTui(ds, 3)
	if err != nil {
		t.Fatalf("GoiTui loi: %v", err)
	}
	r := protocol.NewReader(m.Data)

	if b, _ := r.Byte(); b != 0 {
		t.Fatalf("cAction = %d, mong doi 0", b)
	}
	soO, _ := r.Byte() // client: readUnsignedByte
	if soO != 3 {
		t.Fatalf("soO = %d, mong doi 3", soO)
	}
	kiem := []struct {
		mau     int16
		soLuong int32
		opt     []store.TuyChon
	}{
		{0, 12, []store.TuyChon{{Id: 1, Param: 300}}},
		{-1, 0, nil}, // o rong
		{4, 1, nil},
	}
	for i, k := range kiem {
		tid, err := r.Short()
		if err != nil {
			t.Fatalf("o %d: het byte khi doc templateId: %v", i, err)
		}
		if tid != k.mau {
			t.Fatalf("o %d: templateId = %d, mong doi %d", i, tid, k.mau)
		}
		if tid == -1 {
			continue // dung nhu client: khong doc them byte nao cho o rong
		}
		sl, _ := r.Int()
		if sl != k.soLuong {
			t.Errorf("o %d: so luong = %d, mong doi %d", i, sl, k.soLuong)
		}
		if s, _ := r.UTF(); s != "" {
			t.Errorf("o %d: info = %q, mong doi rong", i, s)
		}
		if s, _ := r.UTF(); s != "" {
			t.Errorf("o %d: content = %q, mong doi rong", i, s)
		}
		n, _ := r.Byte()
		if int(n) != len(k.opt) {
			t.Fatalf("o %d: soOption = %d, mong doi %d", i, n, len(k.opt))
		}
		for j := 0; j < int(n); j++ {
			id, _ := r.Byte() // client: readShortOptionTemp = 1 byte
			p, _ := r.Short() // client: ReadParamDQT = unsigned short
			if int(id) != k.opt[j].Id || int(uint16(p)) != k.opt[j].Param {
				t.Errorf("o %d option %d = (%d,%d), mong doi (%d,%d)",
					i, j, id, uint16(p), k.opt[j].Id, k.opt[j].Param)
			}
		}
	}
	if _, err := r.Byte(); err == nil {
		t.Error("con byte thua sau khi doc het soO o — goi dai hon client mong doi")
	}
}

// ---------------------------------------------------------------------------
// 6. Goi -37: co truong head chen giua
// ---------------------------------------------------------------------------

func TestGoiDoDangMacCoHeadVaODung(t *testing.T) {
	ds := []store.VatPham{mon(1, 2, 0, store.NoiChuaMac, 1)}
	m, err := GoiDoDangMac(ds, 300, SoOMacChuan)
	if err != nil {
		t.Fatalf("GoiDoDangMac loi: %v", err)
	}
	if m.Command != -37 {
		t.Fatalf("lenh = %d, mong doi -37 (Controller.cs:2005)", m.Command)
	}
	want := []byte{
		0x00,       // cAction
		0x01, 0x2C, // head = 300 — SHORT, nam TRUOC soO
		0x05,       // soO = 5
		0x00, 0x02, // o 0: mau 2
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00,
		0x00, 0x00,
		0x00,
		0xFF, 0xFF, // o 1 rong
		0xFF, 0xFF, // o 2 rong
		0xFF, 0xFF, // o 3 rong
		0xFF, 0xFF, // o 4 rong
	}
	if !bytes.Equal(m.Data, want) {
		t.Errorf("than -37 = % X,\nmong doi        % X", m.Data, want)
	}
}

// ---------------------------------------------------------------------------
// 7. Goi -35: than giong het -36
// ---------------------------------------------------------------------------

func TestGoiRuongCungKhuonVoiTui(t *testing.T) {
	dsRuong := []store.VatPham{mon(1, 1, 1, store.NoiChuaRuong, 5)}
	dsTui := []store.VatPham{mon(1, 1, 1, store.NoiChuaTui, 5)}
	mR, err := GoiRuong(dsRuong, 3)
	if err != nil {
		t.Fatalf("GoiRuong loi: %v", err)
	}
	if mR.Command != -35 {
		t.Fatalf("lenh = %d, mong doi -35 (Controller.cs:2265)", mR.Command)
	}
	mT, err := GoiTui(dsTui, 3)
	if err != nil {
		t.Fatalf("GoiTui loi: %v", err)
	}
	if !bytes.Equal(mR.Data, mT.Data) {
		t.Errorf("than -35 = % X,\nthan -36 = % X — hai goi phai cung khuon", mR.Data, mT.Data)
	}
	want := []byte{
		0x00, 0x03,
		0xFF, 0xFF, // o 0 rong
		0x00, 0x01, // o 1: mau 1
		0x00, 0x00, 0x00, 0x05,
		0x00, 0x00,
		0x00, 0x00,
		0x00,
		0xFF, 0xFF, // o 2 rong
	}
	if !bytes.Equal(mR.Data, want) {
		t.Errorf("than -35 = % X,\nmong doi        % X", mR.Data, want)
	}
}

// ---------------------------------------------------------------------------
// 8. Tran 127 o
// ---------------------------------------------------------------------------

func TestGoiTuiTran127(t *testing.T) {
	m, err := GoiTui(nil, 127)
	if err != nil {
		t.Fatalf("127 o phai chap nhan duoc, nhung loi: %v", err)
	}
	if m.Data[1] != 127 {
		t.Errorf("byte soO = %d, mong doi 127", m.Data[1])
	}
	if len(m.Data) != 2+127*2 {
		t.Errorf("do dai = %d, mong doi %d (2 byte dau + 127 o rong x 2)", len(m.Data), 2+127*2)
	}
	if _, err := GoiTui(nil, 128); err == nil {
		t.Error("128 o phai bi tu choi: ME_INFO ghi so o bang byte CO DAU")
	}
	if _, err := GoiTui(nil, -1); err == nil {
		t.Error("so o am phai bi tu choi")
	}
}

// ---------------------------------------------------------------------------
// 9. Cac ca LOI — phai tra error va KHONG co goi (§1.10 / L13)
// ---------------------------------------------------------------------------

func TestGoiTuiTuChoiDuLieuHong(t *testing.T) {
	soMauOption := len(DanhSachMauOption())
	cac := []struct {
		ten string
		ds  []store.VatPham
		soO int
	}{
		{"mau khong ton tai", []store.VatPham{mon(1, 9999, 0, store.NoiChuaTui, 1)}, 3},
		{"mau am", []store.VatPham{mon(1, -1, 0, store.NoiChuaTui, 1)}, 3},
		{"o_so ngoai bien", []store.VatPham{mon(1, 0, 9, store.NoiChuaTui, 1)}, 3},
		{"so luong 0", []store.VatPham{mon(1, 0, 0, store.NoiChuaTui, 0)}, 3},
		{"option id khong co mau", []store.VatPham{
			mon(1, 0, 0, store.NoiChuaTui, 1, store.TuyChon{Id: soMauOption, Param: 1})}, 3},
		{"param am", []store.VatPham{
			mon(1, 0, 0, store.NoiChuaTui, 1, store.TuyChon{Id: 0, Param: -1})}, 3},
		{"param tran", []store.VatPham{
			mon(1, 0, 0, store.NoiChuaTui, 1, store.TuyChon{Id: 0, Param: 65536})}, 3},
		{"hai mon cung o", []store.VatPham{
			mon(1, 0, 1, store.NoiChuaTui, 1),
			mon(2, 1, 1, store.NoiChuaTui, 1)}, 3},
	}
	for _, c := range cac {
		t.Run(c.ten, func(t *testing.T) {
			m, err := GoiTui(c.ds, c.soO)
			if err == nil {
				t.Fatalf("mong doi loi, nhung dung duoc goi: % X", m.Data)
			}
			if m != nil {
				t.Error("co loi ma van tra ve goi — tuyet doi khong duoc gui goi thieu byte")
			}
		})
	}
}

// TestGoiTuiMonThieuOSo: mon nam trong TUI ma o_so = NULL la du lieu hong,
// khong duoc lang le bo qua.
func TestGoiTuiMonThieuOSo(t *testing.T) {
	chu := int64(7)
	ds := []store.VatPham{{
		ID: 1, MauID: 0, ChuSoHuuID: &chu,
		NoiChua: store.NoiChuaTui, OSo: nil, SoLuong: 1,
	}}
	if _, err := GoiTui(ds, 3); err == nil {
		t.Error("mon trong tui khong co o_so phai bi tu choi")
	}
}

// ---------------------------------------------------------------------------
// 10. Loc theo noi chua: mon o noi khac KHONG duoc lot vao goi
// ---------------------------------------------------------------------------

func TestBaGoiLocDungNoiChua(t *testing.T) {
	ds := []store.VatPham{
		mon(1, 0, 0, store.NoiChuaTui, 1),
		mon(2, 1, 0, store.NoiChuaMac, 1),
		mon(3, 2, 0, store.NoiChuaRuong, 1),
		mon(4, 0, 0, store.NoiChuaDuoiDat, 1),
	}
	// Moi ham chi duoc thay DUNG mon cua minh o o 0, con lai rong.
	mTui, err := GoiTui(ds, 2)
	if err != nil {
		t.Fatalf("GoiTui loi: %v", err)
	}
	if mTui.Data[2] != 0x00 || mTui.Data[3] != 0x00 {
		t.Errorf("-36 o 0 phai la mau 0 (mon cua TUI), duoc % X", mTui.Data[2:4])
	}
	mRuong, err := GoiRuong(ds, 2)
	if err != nil {
		t.Fatalf("GoiRuong loi: %v", err)
	}
	if mRuong.Data[2] != 0x00 || mRuong.Data[3] != 0x02 {
		t.Errorf("-35 o 0 phai la mau 2 (mon cua RUONG), duoc % X", mRuong.Data[2:4])
	}
	mMac, err := GoiDoDangMac(ds, 0, 2)
	if err != nil {
		t.Fatalf("GoiDoDangMac loi: %v", err)
	}
	// -37 co them 2 byte head nen o 0 bat dau tu chi so 4.
	if mMac.Data[4] != 0x00 || mMac.Data[5] != 0x01 {
		t.Errorf("-37 o 0 phai la mau 1 (mon dang MAC), duoc % X", mMac.Data[4:6])
	}
}
