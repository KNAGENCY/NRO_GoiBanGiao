package protocol

import (
	"bytes"
	"testing"
)

// Java DataOutputStream ghi so nguyen theo big-endian.
func TestGhiSoNguyenBigEndian(t *testing.T) {
	w := NewWriter()
	w.Short(0x0102)
	w.Int(0x03040506)
	got := w.Bytes()
	want := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	if !bytes.Equal(got, want) {
		t.Fatalf("got % X, mong doi % X", got, want)
	}
}

// writeUTF cua Java: 2 byte do dai big-endian roi toi noi dung.
func TestGhiChuoiKieuJava(t *testing.T) {
	w := NewWriter()
	if err := w.UTF("AB"); err != nil {
		t.Fatalf("loi khong mong doi: %v", err)
	}
	want := []byte{0x00, 0x02, 'A', 'B'}
	if !bytes.Equal(w.Bytes(), want) {
		t.Fatalf("got % X, mong doi % X", w.Bytes(), want)
	}
}

// Chuoi tieng Viet phai giu nguyen sau khi ghi roi doc lai.
func TestChuoiTiengVietVongTron(t *testing.T) {
	w := NewWriter()
	if err := w.UTF("Làng Lá"); err != nil {
		t.Fatalf("loi ghi: %v", err)
	}
	r := NewReader(w.Bytes())
	got, err := r.UTF()
	if err != nil {
		t.Fatalf("loi doc: %v", err)
	}
	if got != "Làng Lá" {
		t.Fatalf("got %q, mong doi %q", got, "Làng Lá")
	}
}

// Chuoi qua 65535 byte khong ghi duoc — phai bao loi thay vi ghi hong.
func TestChuoiQuaDaiBaoLoi(t *testing.T) {
	w := NewWriter()
	qua_dai := string(make([]byte, 70000))
	if err := w.UTF(qua_dai); err == nil {
		t.Fatal("chuoi 70000 byte phai bao loi")
	}
}

func TestDocSoNguyen(t *testing.T) {
	r := NewReader([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06})
	s, err := r.Short()
	if err != nil || s != 0x0102 {
		t.Fatalf("Short = %d, err = %v", s, err)
	}
	i, err := r.Int()
	if err != nil || i != 0x03040506 {
		t.Fatalf("Int = %d, err = %v", i, err)
	}
}

// Doc qua cuoi bo dem phai bao loi, khong duoc panic.
func TestDocQuaCuoiBaoLoi(t *testing.T) {
	r := NewReader([]byte{0x01})
	if _, err := r.Int(); err == nil {
		t.Fatal("doc Int tu 1 byte phai bao loi")
	}
}
