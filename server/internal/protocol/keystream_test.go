package protocol

import (
	"bytes"
	"testing"
)

// Khoa tho [1,2,3] sau bien doi phai thanh [1, 1^2, (1^2)^3] = [1, 3, 0].
func TestDeriveKey(t *testing.T) {
	got := DeriveKey([]byte{1, 2, 3})
	want := []byte{1, 3, 0}
	if !bytes.Equal(got, want) {
		t.Fatalf("DeriveKey = %v, mong doi %v", got, want)
	}
}

// DeriveKey khong duoc sua mang goc.
func TestDeriveKeyKhongSuaMangGoc(t *testing.T) {
	raw := []byte{1, 2, 3}
	_ = DeriveKey(raw)
	if !bytes.Equal(raw, []byte{1, 2, 3}) {
		t.Fatalf("mang goc bi sua: %v", raw)
	}
}

// XOR doi xung: ghi roi doc lai bang mot luong khac phai ra byte ban dau.
func TestGhiRoiDocRaByteBanDau(t *testing.T) {
	key := DeriveKey([]byte{9, 4, 7, 2})
	ben_ghi := NewKeyStream(key)
	ben_doc := NewKeyStream(key)

	goc := []byte{0, 1, 127, 128, 255, 42, 200, 13, 99}
	for i, b := range goc {
		maHoa := ben_ghi.Write(b)
		giaiMa := ben_doc.Read(maHoa)
		if giaiMa != b {
			t.Fatalf("byte thu %d: goc %d, giai ma ra %d", i, b, giaiMa)
		}
	}
}

// Con tro phai quay vong khi vuot qua do dai khoa.
func TestConTroQuayVong(t *testing.T) {
	key := []byte{0xAA, 0xBB}
	ks := NewKeyStream(key)
	if got := ks.Write(0x00); got != 0xAA {
		t.Fatalf("byte 1: %02X, mong doi AA", got)
	}
	if got := ks.Write(0x00); got != 0xBB {
		t.Fatalf("byte 2: %02X, mong doi BB", got)
	}
	if got := ks.Write(0x00); got != 0xAA {
		t.Fatalf("byte 3 phai quay ve dau khoa: %02X, mong doi AA", got)
	}
}

// Hai con tro doc va ghi phai doc lap nhau.
func TestHaiConTroDocLap(t *testing.T) {
	ks := NewKeyStream([]byte{0x11, 0x22})
	ks.Write(0x00) // con tro ghi -> 1
	if got := ks.Read(0x00); got != 0x11 {
		t.Fatalf("con tro doc bi anh huong boi ghi: %02X, mong doi 11", got)
	}
}

func TestGenerateRawKeyChanDoDaiSai(t *testing.T) {
	if _, err := GenerateRawKey(0); err == nil {
		t.Fatal("do dai 0 phai bao loi")
	}
	if _, err := GenerateRawKey(128); err == nil {
		t.Fatal("do dai 128 phai bao loi vi client dung sbyte lam con tro")
	}
	k, err := GenerateRawKey(8)
	if err != nil {
		t.Fatalf("do dai 8 phai hop le: %v", err)
	}
	if len(k) != 8 {
		t.Fatalf("do dai khoa = %d, mong doi 8", len(k))
	}
}
