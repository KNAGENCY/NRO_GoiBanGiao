package protocol

import (
	"bytes"
	"testing"
)

// Truoc khi bat tay xong, goi tin di nguyen ban khong ma hoa.
func TestDocGoiBatTayChuaMaHoa(t *testing.T) {
	// Client gui: cmd = -27, do dai 0 (hai byte 0)
	vao := bytes.NewReader([]byte{0xE5, 0x00, 0x00})
	f := NewFramer(vao, &bytes.Buffer{})
	m, err := f.ReadMessage()
	if err != nil {
		t.Fatalf("loi doc: %v", err)
	}
	if m.Command != -27 {
		t.Fatalf("Command = %d, mong doi -27", m.Command)
	}
	if len(m.Data) != 0 {
		t.Fatalf("Data dai %d, mong doi 0", len(m.Data))
	}
}

// Sau khi co khoa, ghi roi doc lai bang mot Framer khac phai ra goi ban dau.
func TestVongTronSauKhiCoKhoa(t *testing.T) {
	raw := []byte{5, 9, 13, 21}
	var day bytes.Buffer

	ben_gui := NewFramer(&bytes.Buffer{}, &day)
	ben_gui.SetKey(DeriveKey(raw))

	goc := &Message{Command: -29, Data: []byte("xin chao")}
	if err := ben_gui.WriteMessage(goc); err != nil {
		t.Fatalf("loi ghi: %v", err)
	}

	// Ben nhan dung KeyStream rieng nhung cung khoa. Vi ben gui dung con tro
	// ghi con ben nhan dung con tro doc, ta phai gia lap dung chieu: tao Framer
	// moi va doc bang con tro DOC.
	ben_nhan := NewFramer(bytes.NewReader(day.Bytes()), &bytes.Buffer{})
	ben_nhan.SetKey(DeriveKey(raw))

	nhan, err := ben_nhan.ReadMessage()
	if err != nil {
		t.Fatalf("loi doc: %v", err)
	}
	if nhan.Command != goc.Command {
		t.Fatalf("Command = %d, mong doi %d", nhan.Command, goc.Command)
	}
	if !bytes.Equal(nhan.Data, goc.Data) {
		t.Fatalf("Data = %q, mong doi %q", nhan.Data, goc.Data)
	}
}

// Cac lenh trong danh sach dac biet dung do dai 3 byte thay vi 2.
func TestLenhDacBietDungDoDaiBaByte(t *testing.T) {
	var day bytes.Buffer
	f := NewFramer(&bytes.Buffer{}, &day)
	f.SetKey(DeriveKey([]byte{1, 2, 3, 4}))

	// -87 nam trong danh sach dac biet
	if err := f.WriteMessage(&Message{Command: -87, Data: []byte{0xAA}}); err != nil {
		t.Fatalf("loi ghi: %v", err)
	}
	// 1 byte lenh + 3 byte do dai + 1 byte than tin = 5
	if day.Len() != 5 {
		t.Fatalf("goi dai %d byte, mong doi 5 (lenh dac biet dung 3 byte do dai)", day.Len())
	}
}

// Lenh thuong dung do dai 2 byte.
func TestLenhThuongDungDoDaiHaiByte(t *testing.T) {
	var day bytes.Buffer
	f := NewFramer(&bytes.Buffer{}, &day)
	f.SetKey(DeriveKey([]byte{1, 2, 3, 4}))

	if err := f.WriteMessage(&Message{Command: -29, Data: []byte{0xAA}}); err != nil {
		t.Fatalf("loi ghi: %v", err)
	}
	if day.Len() != 4 {
		t.Fatalf("goi dai %d byte, mong doi 4", day.Len())
	}
}

// Goi khong co than tin van phai ghi 2 byte do dai bang 0.
func TestGoiRongVanGhiDoDai(t *testing.T) {
	var day bytes.Buffer
	f := NewFramer(&bytes.Buffer{}, &day)
	if err := f.WriteMessage(&Message{Command: -27}); err != nil {
		t.Fatalf("loi ghi: %v", err)
	}
	if day.Len() != 3 {
		t.Fatalf("goi dai %d byte, mong doi 3", day.Len())
	}
}

// Goi qua lon phai bi tu choi de tranh ke xau bao do dai khong lo lam can RAM.
func TestTuChoiGoiQuaLon(t *testing.T) {
	// cmd -29, do dai 0xFFFF, nhung khong co du lieu theo sau
	vao := bytes.NewReader([]byte{0xE3, 0xFF, 0xFF})
	f := NewFramer(vao, &bytes.Buffer{})
	if _, err := f.ReadMessage(); err == nil {
		t.Fatal("goi bao do dai 65535 ma khong co du lieu phai bao loi")
	}
}
