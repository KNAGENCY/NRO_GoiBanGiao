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

// Lenh dac biet voi than tin RONG van phai ghi du 3 byte do dai, khong duoc
// rut xuong 2 byte. Client (Session_ME.cs, readMessage) quyet dinh doc 3 byte
// CHI dua vao byte lenh, khong xet than tin co rong hay khong — neu server rut
// xuong 2 byte thi client van doc 3, lech dong byte tu do tro di.
func TestLenhDacBietThanRongVanDungBaByte(t *testing.T) {
	var day bytes.Buffer
	f := NewFramer(&bytes.Buffer{}, &day)
	f.SetKey(DeriveKey([]byte{1, 2, 3, 4}))

	// -87 nam trong danh sach dac biet, than tin rong (nil)
	if err := f.WriteMessage(&Message{Command: -87, Data: nil}); err != nil {
		t.Fatalf("loi ghi: %v", err)
	}
	// 1 byte lenh + 3 byte do dai (khong co than tin) = 4
	if day.Len() != 4 {
		t.Fatalf("goi dai %d byte, mong doi 4 (lenh dac biet than tin rong van phai dung 3 byte do dai, khong duoc rut xuong 3)", day.Len())
	}
}

// Test manh hon TestLenhDacBietThanRongVanDungBaByte: khong chi dem do dai goi
// dau, ma ghi LIEN TIEP hai goi roi doc lai bang mot Framer khac cung khoa —
// dung cach TestVongTronSauKhiCoKhoa da lam (con tro ghi va con tro doc la hai
// Framer rieng). Goi dau la lenh dac biet than tin rong, goi sau la lenh thuong
// co than tin. Neu server ghi thieu byte cho goi dau, con tro doc cua goi sau
// se lech dung mot byte — dung hien tuong "treo im lang" ma loi nay gay ra.
//
// Vi ReadMessage() mo phong chieu client->server (luon 2 byte do dai, xem
// docstring ReadMessage), o day ta tu doc goi dac biet "bang tay" y het cach
// client that su doc (Session_ME.cs, readMessage2: luon 3 byte do dai bat ke
// than tin), roi dung ReadMessage() binh thuong cho goi thuong con lai. Neu con
// tro doc con dung cho sau goi dau thi goi sau phai ra dung Command va Data.
func TestGoiDacBietThanRongKhongLamLechGoiSau(t *testing.T) {
	raw := []byte{7, 11, 19, 23}
	var day bytes.Buffer

	ben_gui := NewFramer(&bytes.Buffer{}, &day)
	ben_gui.SetKey(DeriveKey(raw))

	goi1 := &Message{Command: -87, Data: nil}                // lenh dac biet, than tin rong
	goi2 := &Message{Command: -29, Data: []byte{0x10, 0x20, 0x30, 0x40}} // lenh thuong, co than tin

	if err := ben_gui.WriteMessage(goi1); err != nil {
		t.Fatalf("loi ghi goi 1: %v", err)
	}
	if err := ben_gui.WriteMessage(goi2); err != nil {
		t.Fatalf("loi ghi goi 2: %v", err)
	}

	// Ben nhan dung con tro DOC, Framer moi nhung cung khoa dan xuat tu cung raw.
	ben_nhan := NewFramer(bytes.NewReader(day.Bytes()), &bytes.Buffer{})
	ben_nhan.SetKey(DeriveKey(raw))

	// Doc goi 1 "bang tay" — mo phong readMessage2 cua client: cmd roi LUON 3
	// byte do dai (khong xet than tin).
	cmdByte, err := ben_nhan.docMotByte()
	if err != nil {
		t.Fatalf("loi doc lenh goi 1: %v", err)
	}
	cmd1 := int8(ben_nhan.ks.Read(cmdByte))
	if cmd1 != goi1.Command {
		t.Fatalf("Command goi 1 = %d, mong doi %d", cmd1, goi1.Command)
	}
	if !lenhDoDaiBaByte[cmd1] {
		t.Fatalf("cmd1 = %d phai nam trong danh sach lenh dac biet, sai gia thiet test", cmd1)
	}
	var doDai1 [3]byte
	for i := 0; i < 3; i++ {
		b, err := ben_nhan.docMotByte()
		if err != nil {
			t.Fatalf("loi doc byte do dai thu %d cua goi 1 (con tro lech tu day): %v", i, err)
		}
		doDai1[i] = ben_nhan.ks.Read(b) + 128
	}
	n1 := int(doDai1[0]) | int(doDai1[1])<<8 | int(doDai1[2])<<16
	if n1 != 0 {
		t.Fatalf("do dai than tin goi 1 giai ma ra %d, mong doi 0 (con tro doc da lech)", n1)
	}

	// Doc goi 2 bang ReadMessage binh thuong. Neu con tro con dung cho, phai ra
	// dung Command va Data cua goi 2; neu goi 1 ghi thieu byte, cho nay se sai.
	nhan2, err := ben_nhan.ReadMessage()
	if err != nil {
		t.Fatalf("loi doc goi 2 (con tro doc bi lech do goi 1 ghi thieu byte): %v", err)
	}
	if nhan2.Command != goi2.Command {
		t.Fatalf("Command goi 2 = %d, mong doi %d (con tro doc bi lech do goi 1 ghi thieu byte)", nhan2.Command, goi2.Command)
	}
	if !bytes.Equal(nhan2.Data, goi2.Data) {
		t.Fatalf("Data goi 2 = %v, mong doi %v (con tro doc bi lech do goi 1 ghi thieu byte)", nhan2.Data, goi2.Data)
	}
}

// LOI 1 (nghiem trong): lenh THUONG (khong nam trong lenhDoDaiBaByte) voi than
// tin RONG, SAU KHI DA CO KHOA, van phai ma hoa 2 byte do dai bang XOR — dieu
// kien "n > 0" o nhanh co khoa la SAI. Bang chung:
//
//	a) Client doc (Session_ME.cs, readMessage): chi xet getKeyComplete, khong
//	   he xet do dai than tin — "if (getKeyComplete) { ... doc 2 byte da XOR }".
//	b) Chinh client GUI goi RONG (Session_ME.cs, doSendMessage, nhanh
//	   data == null) van tu ma hoa 2 byte khong qua writeKey khi da co khoa.
//	c) Java cu (ServerManager.java): data != null, size == 0 van roi vao nhanh
//	   "else if (session.sentKey())" -> ghi 2 byte CO ma hoa.
//
// Neu server ghi 2 byte RAW (khong XOR) cho goi rong luc da co khoa, client se
// XOR chung ra do dai rac, VA quan trong hon: con tro khoa cua KeyStream phia
// ghi khong duoc tien (vi ks.Write khong duoc goi cho 2 byte do dai), trong
// khi con tro khoa phia doc lai tien (vi ReadMessage luon goi ks.Read cho 2
// byte do dai khi coKhoa()). Hai con tro lech nhau vinh vien tu do, lam MOI
// goi tin sau do giai ma sai — khong chi goi rong nay.
//
// Test nay bat dung hien tuong do: ghi hai goi LIEN TIEP (goi 1 rong, goi 2 co
// than tin) roi doc lai bang Framer khac cung khoa. Neu con tro lech, goi 2 se
// giai ma sai hoac loi.
func TestLenhThuongThanRongVanMaHoaDoDai(t *testing.T) {
	for _, ten := range []string{"Data=nil", "Data=[]byte{}"} {
		ten := ten
		t.Run(ten, func(t *testing.T) {
			var rongData []byte
			if ten == "Data=[]byte{}" {
				rongData = []byte{}
			}

			raw := []byte{2, 4, 6, 8}
			var day bytes.Buffer

			ben_gui := NewFramer(&bytes.Buffer{}, &day)
			ben_gui.SetKey(DeriveKey(raw))

			goi1 := &Message{Command: -29, Data: rongData} // lenh thuong, than tin rong
			goi2 := &Message{Command: -29, Data: []byte{0x11, 0x22, 0x33}} // lenh thuong, co than tin

			if err := ben_gui.WriteMessage(goi1); err != nil {
				t.Fatalf("loi ghi goi 1: %v", err)
			}
			if err := ben_gui.WriteMessage(goi2); err != nil {
				t.Fatalf("loi ghi goi 2: %v", err)
			}

			ben_nhan := NewFramer(bytes.NewReader(day.Bytes()), &bytes.Buffer{})
			ben_nhan.SetKey(DeriveKey(raw))

			nhan1, err := ben_nhan.ReadMessage()
			if err != nil {
				t.Fatalf("loi doc goi 1: %v", err)
			}
			if nhan1.Command != goi1.Command {
				t.Fatalf("Command goi 1 = %d, mong doi %d", nhan1.Command, goi1.Command)
			}
			if len(nhan1.Data) != 0 {
				t.Fatalf("Data goi 1 dai %d, mong doi 0", len(nhan1.Data))
			}

			nhan2, err := ben_nhan.ReadMessage()
			if err != nil {
				t.Fatalf("loi doc goi 2 (con tro khoa lech vi goi 1 ghi do dai khong ma hoa): %v", err)
			}
			if nhan2.Command != goi2.Command {
				t.Fatalf("Command goi 2 = %d, mong doi %d (con tro khoa lech)", nhan2.Command, goi2.Command)
			}
			if !bytes.Equal(nhan2.Data, goi2.Data) {
				t.Fatalf("Data goi 2 = %v, mong doi %v (con tro khoa lech)", nhan2.Data, goi2.Data)
			}
		})
	}
}

// LOI 2 (quan trong): goi WriteMessage voi lenh dac biet TRUOC KHI SetKey phai
// tra ve error, KHONG duoc panic. Nhanh lenh dac biet goi f.ks.Write(...) ma
// khong kiem f.coKhoa() truoc — f.ks con nil luc chua SetKey, cham vao se
// panic. Panic khong duoc recover() se lam sap ca tien trinh may chu, da vang
// TOAN BO nguoi choi dang online chu khong rieng mot ket noi.
func TestLenhDacBietChuaCoKhoaTraLoiKhongPanic(t *testing.T) {
	f := NewFramer(&bytes.Buffer{}, &bytes.Buffer{})
	// -87 nam trong danh sach dac biet, f chua goi SetKey nen f.ks == nil.
	err := f.WriteMessage(&Message{Command: -87, Data: []byte{0x01}})
	if err == nil {
		t.Fatal("mong doi error khi ghi lenh dac biet truoc khi SetKey, nhung khong co loi (va khong panic)")
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
