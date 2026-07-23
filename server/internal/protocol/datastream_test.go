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

// Chuoi ma sau khi ma hoa modified UTF-8 vuot 65535 byte phai bao loi,
// du do dai chuoi Go (so rune hoac so byte UTF-8 chuan) chua vuot.
func TestChuoiMaHoaVuotGioiHanBaoLoi(t *testing.T) {
	w := NewWriter()
	// Moi rune '߿' ma hoa modified UTF-8 thanh 2 byte, giong UTF-8 chuan,
	// nen dung emoji (4 byte UTF-8 chuan -> 6 byte modified UTF-8) de chac chan
	// do dai sau ma hoa vuot 65535 trong khi do dai UTF-8 chuan cua Go con nho hon.
	// 20000 emoji: UTF-8 chuan Go = 80000 byte (da vuot ca hai), nen thu truong hop
	// bien: so byte UTF-8 chuan CHUA vuot nhung modified UTF-8 thi vuot.
	// 15000 emoji => UTF-8 chuan = 60000 byte (< 65535), modified UTF-8 = 90000 byte (> 65535).
	runes := make([]rune, 15000)
	for i := range runes {
		runes[i] = '\U0001F389'
	}
	s := string(runes)
	if len(s) > 65535 {
		t.Fatalf("truong hop kiem tra sai gia thiet: do dai UTF-8 chuan = %d da vuot 65535", len(s))
	}
	if err := w.UTF(s); err == nil {
		t.Fatal("chuoi ma hoa modified UTF-8 vuot 65535 byte phai bao loi")
	}
}

// So byte chuan da do bang Java that: writeUTF("A\x00B") -> 00 04 41 C0 80 42.
// NUL duoc ma hoa thanh 2 byte C0 80, khong phai 1 byte 0x00.
func TestGhiChuoiCoNulKhopJava(t *testing.T) {
	w := NewWriter()
	if err := w.UTF("A\x00B"); err != nil {
		t.Fatalf("loi khong mong doi: %v", err)
	}
	want := []byte{0x00, 0x04, 0x41, 0xC0, 0x80, 0x42}
	if !bytes.Equal(w.Bytes(), want) {
		t.Fatalf("got % X, mong doi % X", w.Bytes(), want)
	}
}

// So byte chuan da do bang Java that: writeUTF("X\U0001F389Y") -> 00 08 58 ED A0 BC ED BE 89 59.
// Ky tu ngoai BMP (emoji) phai tach thanh cap surrogate, moi nua ma hoa 3 byte.
func TestGhiChuoiEmojiKhopJava(t *testing.T) {
	w := NewWriter()
	if err := w.UTF("X\U0001F389Y"); err != nil {
		t.Fatalf("loi khong mong doi: %v", err)
	}
	want := []byte{0x00, 0x08, 0x58, 0xED, 0xA0, 0xBC, 0xED, 0xBE, 0x89, 0x59}
	if !bytes.Equal(w.Bytes(), want) {
		t.Fatalf("got % X, mong doi % X", w.Bytes(), want)
	}
}

// Ghi roi doc lai (vong tron) phai ra dung chuoi ban dau, voi nhieu loai chuoi.
func TestUTFVongTronNhieuLoaiChuoi(t *testing.T) {
	cases := []string{
		"Hello World",      // ASCII thuong
		"Làng Lá xóm Nứa",  // tieng Viet co dau
		"A\x00B",           // co NUL o giua
		"X\U0001F389Y",     // co emoji ngoai BMP
		"🎉🎊🎈 chuc mung!!!", // nhieu emoji + tieng Viet
	}
	for _, s := range cases {
		w := NewWriter()
		if err := w.UTF(s); err != nil {
			t.Fatalf("loi ghi voi chuoi %q: %v", s, err)
		}
		r := NewReader(w.Bytes())
		got, err := r.UTF()
		if err != nil {
			t.Fatalf("loi doc voi chuoi %q: %v", s, err)
		}
		if got != s {
			t.Fatalf("vong tron sai: got %q, mong doi %q", got, s)
		}
	}
}

// Doc dung byte ma Java that gui (writeUTF("A\x00B")) phai ra dung chuoi goc.
func TestDocByteJavaThatCoNul(t *testing.T) {
	b := []byte{0x00, 0x04, 0x41, 0xC0, 0x80, 0x42}
	r := NewReader(b)
	got, err := r.UTF()
	if err != nil {
		t.Fatalf("loi doc: %v", err)
	}
	want := "A\x00B"
	if got != want {
		t.Fatalf("got %q, mong doi %q", got, want)
	}
}

// Doc dung byte ma Java that gui (writeUTF("X\U0001F389Y")) phai ra dung chuoi goc,
// tuc la phai ghep cap surrogate lai thanh 1 rune emoji.
func TestDocByteJavaThatCoEmoji(t *testing.T) {
	b := []byte{0x00, 0x08, 0x58, 0xED, 0xA0, 0xBC, 0xED, 0xBE, 0x89, 0x59}
	r := NewReader(b)
	got, err := r.UTF()
	if err != nil {
		t.Fatalf("loi doc: %v", err)
	}
	want := "X\U0001F389Y"
	if got != want {
		t.Fatalf("got %q, mong doi %q", got, want)
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
