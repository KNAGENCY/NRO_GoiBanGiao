package game

// muaban_test.go — KH7 Task 9: khuon byte hai goi XUONG cua duong mua, va cac
// ham thuan khong can CSDL. Khong test nao o day cham CSDL.

import (
	"encoding/binary"
	"testing"

	"nroserver/internal/store"
)

// Goi XUONG `6` — 16 byte: [vang long][ngoc int][ngocKhoa int=0].
// Con so 6 ghim CUNG o day (khong doc lai hang so tu code): no den tu
// Controller.cs:3677 + Service.java:1390, doc lai hang so thi doi hang so la
// test van xanh.
func TestGoiMuaBanTienByteDung(t *testing.T) {
	m := goiMuaBanTien(1234567890123, 42)
	if m.Command != 6 {
		t.Fatalf("ma lenh = %d, mong doi 6", m.Command)
	}
	if len(m.Data) != 16 {
		t.Fatalf("than goi %d byte, mong doi 16 (long + int + int)", len(m.Data))
	}
	if v := int64(binary.BigEndian.Uint64(m.Data[0:8])); v != 1234567890123 {
		t.Errorf("vang = %d, mong doi 1234567890123", v)
	}
	if v := int32(binary.BigEndian.Uint32(m.Data[8:12])); v != 42 {
		t.Errorf("ngoc = %d, mong doi 42", v)
	}
	// O thu ba la "ngoc khoa" cua client — ban 1 khong co loai tien nay (bay
	// L10). Phai LUON 0; ghi thu khac vao la nguoi choi thay so tien la.
	if v := int32(binary.BigEndian.Uint32(m.Data[12:16])); v != 0 {
		t.Errorf("o thu ba = %d, phai luon 0", v)
	}
}

// Goi XUONG `-25` — writeUTF: [do dai ushort][byte UTF-8].
func TestGoiMuaBanBaoByteDung(t *testing.T) {
	m := goiMuaBanBao("Mua thành công")
	if m == nil {
		t.Fatal("goiMuaBanBao tra nil cho chuoi hop le")
	}
	if m.Command != -25 {
		t.Fatalf("ma lenh = %d, mong doi -25", m.Command)
	}
	than := []byte("Mua thành công")
	if int(binary.BigEndian.Uint16(m.Data[0:2])) != len(than) {
		t.Fatalf("do dai ghi %d, mong doi %d", binary.BigEndian.Uint16(m.Data[0:2]), len(than))
	}
	if string(m.Data[2:]) != string(than) {
		t.Fatalf("noi dung = %q", string(m.Data[2:]))
	}
}

// muaBanTimMon phai duyet MOI tab (PHA A §1.1 — client gui MauID, khong gui tab).
func TestMuaBanTimMonDuyetMoiTab(t *testing.T) {
	ch := &CuaHang{Tabs: []TabCuaHang{
		{Ten: "A", Mon: []MonHang{{MauID: 0, Gia: 1}}},
		{Ten: "B", Mon: []MonHang{{MauID: 7, Gia: 2}}},
	}}
	if m := muaBanTimMon(ch, 7); m == nil || m.Gia != 2 {
		t.Fatalf("khong tim thay mon o tab thu hai: %+v", m)
	}
	if m := muaBanTimMon(ch, 99); m != nil {
		t.Fatalf("mau khong ban ma van tim thay: %+v", m)
	}
}

// muaBanTimOTrong phai tra o NHO NHAT con trong — ket qua tat dinh.
func TestMuaBanTimOTrong(t *testing.T) {
	o := func(v int16) *int16 { return &v }
	tui := []store.VatPham{
		{NoiChua: store.NoiChuaTui, OSo: o(0)},
		{NoiChua: store.NoiChuaTui, OSo: o(2)},
		// Mon o noi chua KHAC khong duoc tinh la chiem o tui.
		{NoiChua: store.NoiChuaRuong, OSo: o(1)},
	}
	got, ok := muaBanTimOTrong(tui, 5)
	if !ok || got != 1 {
		t.Fatalf("o trong = %d (%v), mong doi 1", got, ok)
	}
	// Tui day.
	day := []store.VatPham{
		{NoiChua: store.NoiChuaTui, OSo: o(0)},
		{NoiChua: store.NoiChuaTui, OSo: o(1)},
	}
	if _, ok := muaBanTimOTrong(day, 2); ok {
		t.Fatal("tui day ma van bao con o trong")
	}
}

// typeBuy: chi 0 va 3 duoc nhan (PHA A §1.1).
func TestMuaBanTypeHopLe(t *testing.T) {
	for _, v := range []int8{0, 3} {
		if !muaBanTypeHopLe(v) {
			t.Errorf("typeBuy %d phai duoc nhan", v)
		}
	}
	for _, v := range []int8{-1, 1, 2, 4, 9} {
		if muaBanTypeHopLe(v) {
			t.Errorf("typeBuy %d phai bi tu choi", v)
		}
	}
}

// soThuTuKinhTeKe phai tang MOI LAN goi — ke ca khi lenh sau do that bai.
// Neu no chi tang khi mua thanh cong thi mot lan mua hong roi mot lan mua that
// se dung CUNG so, va lan thu hai bi coi nham la goi lap (mat do that su).
func TestSoThuTuKinhTeLuonTang(t *testing.T) {
	b := &BoXuLy{}
	var truoc uint64
	for i := 0; i < 5; i++ {
		n := b.soThuTuKinhTeKe()
		if n <= truoc {
			t.Fatalf("so thu tu khong tang: %d sau %d", n, truoc)
		}
		truoc = n
	}
}

// Khoa idempotency phai ghep dung khuon `<nhanVatId>:MUA:<soThuTu>` (D3).
func TestKhoaIdemMuaKhuon(t *testing.T) {
	if got := store.KhoaIdemMua(12, 7); got != "12:MUA:7" {
		t.Fatalf("khoa = %q, mong doi \"12:MUA:7\"", got)
	}
}
