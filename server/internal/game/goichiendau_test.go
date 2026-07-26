package game

import (
	"bytes"
	"testing"
)

// TestGoiNguoiDanhQuaiByteDung: goi 54, dac ta §3.3 — CHI 3 truong, KHONG co
// truong dem sau mobId (dac ta §2.0: client doc mobId cho toi het than goi).
func TestGoiNguoiDanhQuaiByteDung(t *testing.T) {
	m := goiNguoiDanhQuai(1001, 5, 3)

	if m.Command != 54 {
		t.Fatalf("Command = %d, muon 54", m.Command)
	}

	muon := []byte{
		0, 0, 3, 233, // idNguoiDanh = 1001 (int32 BE, 0x3E9)
		5, // skillId
		3, // mobId — PHAI la byte CUOI CUNG cua than goi
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi 54 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// TestGoiQuaiConSongByteDung: goi -9, dac ta §3.3. Hieu ung -1 (khong hieu
// ung) PHAI la 1 byte 0xFF, KHONG phai 4 byte writeInt(-1) nhu ban Java.
func TestGoiQuaiConSongByteDung(t *testing.T) {
	m := goiQuaiConSong(3, 80, 20, false, 255) // 255 = byte cua -1 (khong hieu ung)

	if m.Command != -9 {
		t.Fatalf("Command = %d, muon -9", m.Command)
	}

	muon := []byte{
		3,           // mobId
		0, 0, 0, 80, // hpConLai = 80 (int32)
		0, 0, 0, 20, // satThuong = 20 (int32)
		0,   // chiMang = false
		255, // hieuUngId = -1 (1 BYTE, khong phai 4 byte)
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -9 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// TestGoiQuaiConSongChiMangTrueByteDung: kiem gia tri chiMang=true va
// hieuUngId khac -1 de khong trung tinh co voi test truoc.
func TestGoiQuaiConSongChiMangTrueByteDung(t *testing.T) {
	m := goiQuaiConSong(7, 500, 999, true, 4)

	muon := []byte{
		7,            // mobId
		0, 0, 1, 244, // hpConLai = 500
		0, 0, 3, 231, // satThuong = 999
		1, // chiMang = true
		4, // hieuUngId = 4
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -9 (chiMang=true) byte sai:\n got  = %v\n muon = %v", m.Data, muon)
	}
}

// TestGoiQuaiChetByteDung_SoDoRoiRong: goi -12 voi soDoRoi = 0 (KH6 chua lam
// he vat pham, dac ta §0 "Ngoai KH6"), theo dung yeu cau cua brief Task 6.
func TestGoiQuaiChetByteDung_SoDoRoiRong(t *testing.T) {
	m := goiQuaiChet(2, 150, true, nil)

	if m.Command != -12 {
		t.Fatalf("Command = %d, muon -12", m.Command)
	}

	muon := []byte{
		2,            // mobId
		0, 0, 0, 150, // satThuong = 150
		1, // chiMang = true
		0, // soDoRoi = 0
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -12 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// TestGoiQuaiChetByteDung_CoDoRoi: goi -12 voi 2 do roi, kiem TOAN BO 5 truong
// cua moi phan tu lap (itemMapId, templateId, x, y, idNguoiNhat) dung thu tu
// dac ta §3.3.
func TestGoiQuaiChetByteDung_CoDoRoi(t *testing.T) {
	doRoi := []DoRoi{
		{ItemMapId: 10, TemplateId: 20, X: 100, Y: 200, IdNguoiNhat: 1001},
		{ItemMapId: 11, TemplateId: 21, X: 101, Y: 201, IdNguoiNhat: 1002},
	}
	m := goiQuaiChet(5, 999, false, doRoi)

	muon := []byte{
		5,            // mobId
		0, 0, 3, 231, // satThuong = 999
		0, // chiMang = false
		2, // soDoRoi = 2

		0, 10, // itemMapId = 10
		0, 20, // templateId = 20
		0, 100, // x = 100
		0, 200, // y = 200
		0, 0, 3, 233, // idNguoiNhat = 1001

		0, 11, // itemMapId = 11
		0, 21, // templateId = 21
		0, 101, // x = 101
		0, 201, // y = 201
		0, 0, 3, 234, // idNguoiNhat = 1002
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -12 (co do roi) byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// TestGoiQuaiHoiSinhByteDung: goi -13, dac ta §3.3. sys va levelBoss LUON 0
// (dac ta: client khong doi loai quai khi hoi sinh, levelBoss != 0 bat hieu
// ung sieu cap ngau nhien khong mong muon).
func TestGoiQuaiHoiSinhByteDung(t *testing.T) {
	m := goiQuaiHoiSinh(4, 250)

	if m.Command != -13 {
		t.Fatalf("Command = %d, muon -13", m.Command)
	}

	muon := []byte{
		4,            // mobId
		0,            // sys = 0
		0,            // levelBoss = 0
		0, 0, 0, 250, // hp = 250
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -13 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}
