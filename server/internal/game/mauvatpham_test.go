package game

import (
	"bytes"
	"testing"

	"nroserver/internal/protocol"
)

// ---------------------------------------------------------------------------
// Cac test o day la BYTE-CHINH-XAC theo dung ket qua PHA A
// (docs/superpowers/research/2026-09-03-kh7-giai-ma-28-8.md). Moi con so duoc
// GHIM DOC LAP voi ham dung goi — chep cong thuc tu ham sang test thi test chi
// chung minh ham bang chinh no.
// ---------------------------------------------------------------------------

// utfMongDoi dung mot chuoi theo dung khuon Java writeUTF: 2 byte do dai (big
// endian) + byte UTF-8. Ghim doc lap voi protocol.Writer.UTF.
func utfMongDoi(s string) []byte {
	b := []byte(s)
	if len(b) > 0xFFFF {
		panic("chuoi qua dai cho writeUTF")
	}
	return append([]byte{byte(len(b) >> 8), byte(len(b))}, b...)
}

func TestGoiCatalogItemSubtype0DungByte(t *testing.T) {
	m, err := goiCatalogItemOption()
	if err != nil {
		t.Fatalf("goiCatalogItemOption loi: %v", err)
	}
	if m.Command != -28 {
		t.Fatalf("lenh = %d, mong doi -28", m.Command)
	}
	ds := DanhSachMauOption()
	// [8][vsItem][0][soMau 1 byte KHONG DAU] ( [ten UTF][loai byte] )*
	want := []byte{lenhConCatalogItem, vsItem, 0, byte(len(ds))}
	for _, o := range ds {
		want = append(want, utfMongDoi(o.Ten)...)
		want = append(want, byte(o.Loai))
	}
	if !bytes.Equal(m.Data, want) {
		t.Errorf("than -28/8 subtype 0 = %v,\nmong doi                  %v", m.Data, want)
	}
}

func TestGoiCatalogItemSubtype1DungByte(t *testing.T) {
	m, err := goiCatalogItemMau()
	if err != nil {
		t.Fatalf("goiCatalogItemMau loi: %v", err)
	}
	if m.Command != -28 {
		t.Fatalf("lenh = %d, mong doi -28", m.Command)
	}
	ds := DanhSachMauVatPham()
	n := len(ds)
	// [8][vsItem][1][soMau short CO DAU] ( 9 truong )*
	want := []byte{lenhConCatalogItem, vsItem, 1, byte(n >> 8), byte(n)}
	for _, v := range ds {
		want = append(want, byte(v.Loai), byte(v.GioiTinh))
		want = append(want, utfMongDoi(v.Ten)...)
		want = append(want, utfMongDoi(v.MoTa)...)
		want = append(want, byte(v.Cap))
		want = append(want,
			byte(v.YeuCauSuc>>24), byte(v.YeuCauSuc>>16), byte(v.YeuCauSuc>>8), byte(v.YeuCauSuc))
		want = append(want, byte(v.IconID>>8), byte(v.IconID))
		want = append(want, byte(v.Phan>>8), byte(v.Phan))
		if v.NangCapDuoc {
			want = append(want, 1)
		} else {
			want = append(want, 0)
		}
	}
	if !bytes.Equal(m.Data, want) {
		t.Errorf("than -28/8 subtype 1 = %v,\nmong doi                  %v", m.Data, want)
	}
}

// TestGoiCatalogItemSubtype1DocLaiDungThuTuTruong doc NGUOC lai goi theo dung
// trinh tu lenh doc cua client (Controller.cs:6498) va doi chieu voi danh muc.
// Test byte-chinh-xac phia tren bat duoc "byte khac"; test nay bat duoc "byte
// giong nhung Y NGHIA lech" — vi du doi cho IconID voi Phan (ca hai deu short).
func TestGoiCatalogItemSubtype1DocLaiDungThuTuTruong(t *testing.T) {
	m, err := goiCatalogItemMau()
	if err != nil {
		t.Fatalf("goiCatalogItemMau loi: %v", err)
	}
	r := protocol.NewReader(m.Data)
	if b, _ := r.Byte(); b != lenhConCatalogItem {
		t.Fatalf("lenh con = %d, mong doi 8", b)
	}
	if b, _ := r.Byte(); b != vsItem {
		t.Fatalf("vsItem = %d, mong doi %d", b, vsItem)
	}
	if b, _ := r.Byte(); b != 1 {
		t.Fatalf("subtype = %d, mong doi 1", b)
	}
	n, err := r.Short()
	if err != nil {
		t.Fatalf("doc so mau loi: %v", err)
	}
	ds := DanhSachMauVatPham()
	if int(n) != len(ds) {
		t.Fatalf("so mau trong goi = %d, danh muc co %d", n, len(ds))
	}
	for i := 0; i < int(n); i++ {
		muon := ds[i]
		loai, _ := r.Byte()
		gioi, _ := r.Byte()
		ten, _ := r.UTF()
		moTa, _ := r.UTF()
		cap, _ := r.Byte()
		suc, _ := r.Int()
		icon, _ := r.Short()
		phan, _ := r.Short()
		nangCap, errCuoi := r.Byte()
		if errCuoi != nil {
			t.Fatalf("mau %d doc thieu byte: %v", i, errCuoi)
		}
		if int8(loai) != muon.Loai || int8(gioi) != muon.GioiTinh ||
			ten != muon.Ten || moTa != muon.MoTa || int8(cap) != muon.Cap ||
			suc != muon.YeuCauSuc || icon != muon.IconID || phan != muon.Phan ||
			(nangCap == 1) != muon.NangCapDuoc {
			t.Errorf("mau %d doc lai lech: got{%d %d %q %q %d %d %d %d %d} want %+v",
				i, loai, gioi, ten, moTa, cap, suc, icon, phan, nangCap, muon)
		}
	}
	if _, err := r.Byte(); err == nil {
		t.Error("goi -28/8 subtype 1 con byte thua sau khi doc du so mau")
	}
}

// TestChuoiCatalogItemDungThuTu: subtype 2 PHAI o cuoi — no la cai duy nhat
// kich Service.clientOk() phia client (Controller.cs:6517-6534). Dao thu tu =
// nguoi choi kẹt "Xin chờ" ma khong co loi nao bat ra.
func TestChuoiCatalogItemDungThuTu(t *testing.T) {
	ds, err := goiCatalogItemDayDu()
	if err != nil {
		t.Fatalf("goiCatalogItemDayDu loi: %v", err)
	}
	if len(ds) != 3 {
		t.Fatalf("so message = %d, mong doi 3 (subtype 0, 1, 2)", len(ds))
	}
	wantSub := []byte{0, 1, 2}
	for i, m := range ds {
		if m.Command != -28 {
			t.Errorf("message %d lenh = %d, mong doi -28", i, m.Command)
		}
		if len(m.Data) < 3 {
			t.Fatalf("message %d qua ngan: %v", i, m.Data)
		}
		if m.Data[0] != lenhConCatalogItem {
			t.Errorf("message %d lenh con = %d, mong doi 8", i, m.Data[0])
		}
		if m.Data[1] != vsItem {
			t.Errorf("message %d vsItem = %d, mong doi %d", i, m.Data[1], vsItem)
		}
		if m.Data[2] != wantSub[i] {
			t.Errorf("message %d subtype = %d, mong doi %d", i, m.Data[2], wantSub[i])
		}
	}
	if ds[2].Data[2] != 2 {
		t.Error("message CUOI CUNG phai la subtype 2 — chi no moi kich clientOk")
	}
}

// ---------------------------------------------------------------------------
// Tran va bat bien — PHA A §8
// ---------------------------------------------------------------------------

func TestTranDanhMucKhongVuotSucDocCuaClient(t *testing.T) {
	if n := len(DanhSachMauOption()); n > tranMauOption {
		t.Errorf("bang mau option %d muc > %d — client doc bang MOT byte khong dau, catalog se rong im lang", n, tranMauOption)
	}
	if n := len(DanhSachMauVatPham()); n > tranMauVatPham {
		t.Errorf("danh muc %d mau > %d — client doc bang short CO DAU, so am = catalog rong im lang", n, tranMauVatPham)
	}
}

func TestThanGoiItemDuoiMocMark100000(t *testing.T) {
	ds, err := goiCatalogItemDayDu()
	if err != nil {
		t.Fatalf("goiCatalogItemDayDu loi: %v", err)
	}
	for i, m := range ds {
		if len(m.Data) >= tranThanGoiItem {
			t.Errorf("message %d = %d byte, phai < %d (client mark(100000))", i, len(m.Data), tranThanGoiItem)
		}
	}
}

// TestOptionBamCungIdSauCuaClient ghim hai chi so client BAM CUNG:
// ItemOption.cs:29-38 doi option 22 -> 6 va 23 -> 7 (param x1000). Neu ai do
// sap xep lai bang mau option thi hai dong nay bao dong truoc khi nguoi choi
// nhin thay chi so sai gap nghin lan.
func TestOptionBamCungIdSauCuaClient(t *testing.T) {
	ds := DanhSachMauOption()
	if len(ds) <= 7 {
		t.Fatalf("bang mau option chi %d muc — client BAM CUNG iOptionTemplates[6] va [7] (ItemOption.cs:31,36) va KHONG kiem bien -> IndexOutOfRange", len(ds))
	}
	if ds[6].Ten != "Tăng # HP" {
		t.Errorf("mau option[6] = %q, phai la chi so HP — client doi option 22 sang day voi param x1000", ds[6].Ten)
	}
	if ds[7].Ten != "Tăng # chakra" {
		t.Errorf("mau option[7] = %q, phai la chi so chakra — client doi option 23 sang day voi param x1000", ds[7].Ten)
	}
}

// TestTenMauOptionCoDauThang: client thay '#' bang param khi ve
// (ItemOption.getOptionString, ItemOption.cs:45). Ten khong co '#' = nguoi choi
// khong thay con so nao, khong loi.
func TestTenMauOptionCoDauThang(t *testing.T) {
	for i, o := range DanhSachMauOption() {
		if !bytes.ContainsRune([]byte(o.Ten), '#') {
			t.Errorf("mau option %d (%q) khong co '#' — client se ve chu tinh, mat con so", i, o.Ten)
		}
	}
}

// TestDanhMucVatPhamHopLe: cac bat bien cua chinh danh muc.
func TestDanhMucVatPhamHopLe(t *testing.T) {
	for i, m := range DanhSachMauVatPham() {
		if m.Ten == "" {
			t.Errorf("mau %d khong co ten", i)
		}
		if m.TranChong < 1 {
			t.Errorf("mau %d TranChong = %d, phai >= 1", i, m.TranChong)
		}
		// Loai 0..5 = trang bi mac tren nguoi (Item.isTypeBody, Item.cs:271-278).
		// Trang bi thi khong duoc chong dong: mot mon co option la mot ca the rieng
		// (rang buoc CSDL vat_pham_co_option_thi_don_chiec, task 1).
		if m.Loai >= 0 && m.Loai < 6 && m.TranChong != 1 {
			t.Errorf("mau %d la trang bi (loai %d) nhung TranChong = %d, phai la 1", i, m.Loai, m.TranChong)
		}
	}
}

// TestTimMauVatPhamKiemBien: bay §7.11 — chi so tu client KHONG BAO GIO duoc
// dung thang. TimMauVatPham phai tu choi thay vi panic.
func TestTimMauVatPhamKiemBien(t *testing.T) {
	n := int16(len(DanhSachMauVatPham()))
	for _, id := range []int16{-1, -32768, n, n + 1, 32767} {
		if m := TimMauVatPham(id); m != nil {
			t.Errorf("TimMauVatPham(%d) tra %+v, mong doi nil", id, m)
		}
	}
	for id := int16(0); id < n; id++ {
		if TimMauVatPham(id) == nil {
			t.Errorf("TimMauVatPham(%d) tra nil — danh muc phai lien tuc tu 0", id)
		}
	}
}

// ---------------------------------------------------------------------------
// CA BIEN — dung bo dung goi voi du lieu bien, khong dung danh muc that
// ---------------------------------------------------------------------------

// ghiSubtype1 la ban sao TOI THIEU cua vong ghi trong goiCatalogItemMau, chi de
// dung goi tu du lieu gia cho cac ca bien. Neu ai sua khuon o goiCatalogItemMau
// ma quen sua o day thi TestGoiCatalogItemSubtype1DungByte se do truoc.
func ghiSubtype1(ds []MauVatPham) []byte {
	w := protocol.NewWriter()
	w.Byte(lenhConCatalogItem)
	w.Byte(vsItem)
	w.Byte(1)
	w.Short(int16(len(ds)))
	for _, m := range ds {
		w.Byte(byte(m.Loai))
		w.Byte(byte(m.GioiTinh))
		_ = w.UTF(m.Ten)
		_ = w.UTF(m.MoTa)
		w.Byte(byte(m.Cap))
		w.Int(m.YeuCauSuc)
		w.Short(m.IconID)
		w.Short(m.Phan)
		w.Bool(m.NangCapDuoc)
	}
	return w.Bytes()
}

// TestCaBienKhongMauNao: 0 mau -> chi con 5 byte dau, khong byte thua.
func TestCaBienKhongMauNao(t *testing.T) {
	got := ghiSubtype1(nil)
	want := []byte{lenhConCatalogItem, vsItem, 1, 0, 0}
	if !bytes.Equal(got, want) {
		t.Errorf("subtype 1 voi 0 mau = %v, mong doi %v", got, want)
	}
}

// TestCaBienTenDai: ten 300 byte van ghi dung (writeUTF dung 2 byte do dai,
// khong phai 1 — mot byte se cat vong tai 256 va lam lech CA GOI tu do tro di).
func TestCaBienTenDai(t *testing.T) {
	tenDai := string(bytes.Repeat([]byte("A"), 300))
	got := ghiSubtype1([]MauVatPham{{Ten: tenDai, MoTa: "", Phan: -1}})
	// [8][vs][1][0 1] roi [loai 0][gioi 0][UTF 300][UTF 0][cap 0][int 0][short 0][short -1][bool 0]
	want := []byte{lenhConCatalogItem, vsItem, 1, 0, 1, 0, 0}
	want = append(want, 0x01, 0x2C) // do dai 300 = 0x012C, HAI byte
	want = append(want, bytes.Repeat([]byte("A"), 300)...)
	want = append(want, 0, 0) // MoTa rong
	want = append(want, 0)    // cap
	want = append(want, 0, 0, 0, 0)
	want = append(want, 0, 0)       // icon 0
	want = append(want, 0xFF, 0xFF) // phan -1
	want = append(want, 0)          // nangCap false
	if !bytes.Equal(got, want) {
		t.Errorf("subtype 1 ten dai lech:\ngot  %v\nwant %v", got, want)
	}
}

// TestCaBienOptionId255: 255 mau option la GIOI HAN TREN vua khit cua mot byte
// khong dau. 255 phai ghi ra byte 0xFF (khong phai 0), va toan bo van doc lai duoc.
func TestCaBienOptionId255(t *testing.T) {
	ds := make([]MauOption, 255)
	for i := range ds {
		ds[i] = MauOption{Ten: "o#", Loai: 0}
	}
	w := protocol.NewWriter()
	w.Byte(lenhConCatalogItem)
	w.Byte(vsItem)
	w.Byte(0)
	w.Byte(byte(len(ds)))
	for _, o := range ds {
		_ = w.UTF(o.Ten)
		w.Byte(byte(o.Loai))
	}
	got := w.Bytes()
	if got[3] != 0xFF {
		t.Fatalf("so mau 255 ghi ra byte %d, mong doi 255 (0xFF)", got[3])
	}
	// 256 thi tran ve 0 — day chinh la lech L-A cua PHA A. Ghim de ai do noi tran
	// se thay ngay vi sao khong duoc noi.
	soVuot := 256 // bien, khong phai hang: Go chan byte(256) ngay luc bien dich
	if byte(soVuot) != 0 {
		t.Fatal("gia dinh sai: byte(256) phai tran ve 0")
	}
	if len(DanhSachMauOption()) > 255 {
		t.Error("bang mau option that da vuot 255 — se tran im lang")
	}
}

// TestCaBienParamUnsignedShort ghim su that PHA A §7.1/§8.4: param cua option
// la UNSIGNED short phia client (myReader.ReadParamDQT -> readUnsignedShort vi
// ModFunc.isReadInt == true), mien 0..65535. param KHONG nam trong goi -28/8 —
// no thuoc goi vat pham (task 3). Test nay ghim CACH GHI dung de task 3 khong
// phai giai ma lai: ghi 65535 bang Short co dau phai ra 0xFF 0xFF, va doc lai
// KHONG DAU phai ra 65535 chu khong phai -1.
func TestCaBienParamUnsignedShort(t *testing.T) {
	w := protocol.NewWriter()
	paramLon := uint16(65535)
	w.Short(int16(paramLon)) // uint16 -> int16: bien, khong phai hang (Go chan hang tran)
	got := w.Bytes()
	if !bytes.Equal(got, []byte{0xFF, 0xFF}) {
		t.Fatalf("param 65535 ghi ra %v, mong doi [255 255]", got)
	}
	if v := uint16(got[0])<<8 | uint16(got[1]); v != 65535 {
		t.Fatalf("doc lai KHONG DAU = %d, mong doi 65535", v)
	}
	// Neu ai do ghi param bang so CO DAU: -1 ra dung hai byte 0xFF 0xFF nhung
	// moi so am khac deu ra mot con so khong lo phia client, KHONG loi, KHONG log.
	w2 := protocol.NewWriter()
	w2.Short(-1)
	if v := uint16(w2.Bytes()[0])<<8 | uint16(w2.Bytes()[1]); v != 65535 {
		t.Fatalf("short -1 doc khong dau = %d — client se hien %d", v, v)
	}
}
