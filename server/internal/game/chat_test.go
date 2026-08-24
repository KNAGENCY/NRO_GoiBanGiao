package game

import (
	"bytes"
	"testing"

	"nroserver/internal/protocol"
)

// docLaiChat44 doc lai than goi 44 bang protocol.NewReader dung THU TU/KIEU ma
// client Unity doc (Controller.cs:4569-4570: readInt roi readUTF), va KHANG
// DINH HET BUFFER — sau noiDung khong duoc con byte nao.
//
// Vi sao phai khang dinh het buffer: client doc dung 2 truong roi `break`, byte
// thua bi vut im lang o msg.cleanup() — KHONG loi, KHONG trieu chung. Neu test
// khong soi cho nay thi loi "ghi thua truong" se lot luoi hoan toan.
func docLaiChat44(t *testing.T, m *protocol.Message) (int32, string) {
	t.Helper()

	r := protocol.NewReader(m.Data)

	id, err := r.Int() // 1 id nguoi noi — client: readInt()
	if err != nil {
		t.Fatalf("doc truong 1 (id nguoi noi) loi: %v", err)
	}
	noiDung, err := r.UTF() // 2 noi dung — client: readUTF()
	if err != nil {
		t.Fatalf("doc truong 2 (noi dung) loi: %v", err)
	}

	if _, err := r.Byte(); err == nil {
		t.Fatalf("goi 44 con byte THUA sau noi dung (than dai %d byte) — "+
			"client vut im lang, dung loai loi khong the phat hien khi chay that",
			len(m.Data))
	}

	return id, noiDung
}

// TestGoiChat44ByteDung: khuon byte co ban cua goi 44 — int id + UTF noi dung,
// DUNG 2 TRUONG. Bang chung: Service.java:599-601 va Controller.cs:4567-4571.
func TestGoiChat44ByteDung(t *testing.T) {
	m := goiChat44(1001, "hi")

	if m.Command != 44 {
		t.Fatalf("Command = %d, muon 44", m.Command)
	}

	muon := []byte{
		0, 0, 3, 233, // idNguoiNoi = 1001 (int32 BE, 0x3E9) — writeInt, KHONG phai short/long
		0, 2, // do dai UTF = 2 byte
		104, 105, // "hi"
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi 44 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}

	id, noiDung := docLaiChat44(t, m)
	if id != 1001 {
		t.Errorf("id doc lai = %d, muon 1001", id)
	}
	if noiDung != "hi" {
		t.Errorf("noi dung doc lai = %q, muon %q", noiDung, "hi")
	}
}

// TestGoiChat44IdLaIntBonByteBigEndian: id PHAI la int 4 byte big-endian. Dung
// gia tri co 4 byte KHAC NHAU HET de bat duoc moi kieu sai: ghi Short (mat 2
// byte dau), ghi Long (thua 4 byte), hay ghi little-endian (dao thu tu).
func TestGoiChat44IdLaIntBonByteBigEndian(t *testing.T) {
	const id int64 = 0x01020304 // = 16909060

	m := goiChat44(id, "x")

	if m.Command != 44 {
		t.Fatalf("Command = %d, muon 44", m.Command)
	}

	muon := []byte{
		1, 2, 3, 4, // idNguoiNoi big-endian — 4 byte, thu tu byte cao truoc
		0, 1, // do dai UTF = 1
		120, // "x"
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi 44 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}

	docId, noiDung := docLaiChat44(t, m)
	if docId != 0x01020304 {
		t.Errorf("id doc lai = %d (0x%X), muon %d (0x%X)", docId, docId, id, id)
	}
	if noiDung != "x" {
		t.Errorf("noi dung doc lai = %q, muon %q", noiDung, "x")
	}
}

// TestGoiChat44TiengVietCoDau: chuoi tieng Viet co dau — cac ky tu nhu "à"/"á"
// (U+00E0/U+00E1) ton 2 byte modified UTF-8, nen do dai UTF (16) KHAC han so
// ky tu (13). Byte mong doi tinh doc lap bang bo ma UTF-8 cua Python, khong
// lay tu chinh ban Go.
func TestGoiChat44TiengVietCoDau(t *testing.T) {
	const cau = "Chào Làng Lá!"

	m := goiChat44(7, cau)

	if m.Command != 44 {
		t.Fatalf("Command = %d, muon 44", m.Command)
	}

	muon := []byte{
		0, 0, 0, 7, // idNguoiNoi = 7
		0, 16, // do dai UTF = 16 byte (13 ky tu, 3 ky tu co dau moi cai 2 byte)
		67, 104, 195, 160, 111, // "Chào"  (à = 0xC3 0xA0)
		32,                     // dau cach
		76, 195, 160, 110, 103, // "Làng"  (à = 0xC3 0xA0)
		32,           // dau cach
		76, 195, 161, // "Lá"    (á = 0xC3 0xA1)
		33, // "!"
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi 44 tieng Viet byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}

	id, noiDung := docLaiChat44(t, m)
	if id != 7 {
		t.Errorf("id doc lai = %d, muon 7", id)
	}
	if noiDung != cau {
		t.Errorf("noi dung doc lai = %q, muon %q", noiDung, cau)
	}
}

// TestGoiChat44TiengVietKyTuBaByte: chuoi tieng Viet co ky tu ngoai Latin-1 —
// "ệ" (U+1EC7 >= 0x800) ton 3 byte modified UTF-8, "Đ" (U+0110) ton 2 byte.
// Bao phu ca hai nhanh 2-byte va 3-byte cua bo ma, dieu ma "Chào Làng Lá!"
// (chi co nhanh 2-byte) khong lam duoc.
func TestGoiChat44TiengVietKyTuBaByte(t *testing.T) {
	const cau = "Chào Đệ!"

	m := goiChat44(7, cau)

	if m.Command != 44 {
		t.Fatalf("Command = %d, muon 44", m.Command)
	}

	muon := []byte{
		0, 0, 0, 7, // idNguoiNoi = 7
		0, 12, // do dai UTF = 12 byte (8 ky tu)
		67, 104, 195, 160, 111, // "Chào"
		32,       // dau cach
		196, 144, // "Đ" = U+0110 -> 2 byte 0xC4 0x90
		225, 187, 135, // "ệ" = U+1EC7 -> 3 byte 0xE1 0xBB 0x87
		33, // "!"
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi 44 ky tu 3 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}

	id, noiDung := docLaiChat44(t, m)
	if id != 7 {
		t.Errorf("id doc lai = %d, muon 7", id)
	}
	if noiDung != cau {
		t.Errorf("noi dung doc lai = %q, muon %q", noiDung, cau)
	}
}

// TestGoiChat44ChuoiRong: chuoi rong VAN phai ghi 2 byte do dai = 0 (writeUTF
// luon co tien to do dai). Bo hai byte nay se lam client readUTF an mat 2 byte
// cua goi ke tiep tren cung luong — day la ly do phai co test rieng cho chuoi
// rong thay vi coi la truong hop hien nhien.
func TestGoiChat44ChuoiRong(t *testing.T) {
	m := goiChat44(1001, "")

	if m.Command != 44 {
		t.Fatalf("Command = %d, muon 44", m.Command)
	}

	muon := []byte{
		0, 0, 3, 233, // idNguoiNoi = 1001
		0, 0, // do dai UTF = 0 — VAN PHAI CO 2 byte nay
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi 44 chuoi rong byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
	if len(m.Data) != 6 {
		t.Fatalf("goi 44 chuoi rong dai %d byte, phai DUNG 6 (int 4 + do dai 2)", len(m.Data))
	}

	id, noiDung := docLaiChat44(t, m)
	if id != 1001 {
		t.Errorf("id doc lai = %d, muon 1001", id)
	}
	if noiDung != "" {
		t.Errorf("noi dung doc lai = %q, muon chuoi rong", noiDung)
	}
}
