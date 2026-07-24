package protocol

import (
	"bytes"
	"io"
	"testing"
)

// Kiem thu nay bat loi nguy hiem nhat: hai chieu ma hoa lech nhau.
// No gia lap dung trinh tu that — may chu ghi, client doc — bang cach dung
// hai KeyStream rieng tren cung mot khoa, giong het hai dau ket noi that.
//
// LUU Y quan trong (khac voi ban brief goc): o day KHONG dung
// Framer.ReadMessage() truc tiep de "client" doc lai, vi ReadMessage() chi mo
// hinh chieu client->server (xem docstring Framer trong frame.go) va CO Y
// khong biet doc do dai 3 byte cua cac lenh dac biet (client that khong bao
// gio gui nhung lenh do LEN server). Client that (Session_ME.cs) doc goi tu
// SERVER GUI XUONG bang mot ham rieng (readMessage2) luon doc du 3 byte do
// dai cho lenh dac biet, bat ke than tin dai bao nhieu. Ham docGoiNhuClientThat
// duoi day mo phong dung ham do, y het cach TestGoiDacBietThanRongKhongLamLechGoiSau
// (frame_test.go) da lam thu cong cho tung truong hop rieng le — o day ghep
// lai thanh mot vong lap doc nhieu goi lien tiep de bat loi lech con tro tich
// luy dan qua ca gioi han lenh thuong LAN lenh dac biet.
func TestMayChuGhiClientDocDungThuTu(t *testing.T) {
	khoaTho := []byte{0x3A, 0x7F, 0x01, 0xC4, 0x55, 0x09, 0xE1, 0x88}
	khoa := DeriveKey(khoaTho)

	var day bytes.Buffer

	// Ben may chu: chi ghi
	mayChu := NewFramer(&bytes.Buffer{}, &day)
	mayChu.SetKey(khoa)

	// Gui nhieu goi lien tiep de bat loi lech con tro tich luy dan.
	goiGui := []*Message{
		{Command: -26, Data: []byte("Sai tài khoản hoặc mật khẩu.")},
		{Command: -30},
		{Command: -87, Data: bytes.Repeat([]byte{0xAB}, 300)}, // lenh do dai 3 byte
		{Command: 0, Data: []byte{1, 2, 3}},
	}
	for i, m := range goiGui {
		if err := mayChu.WriteMessage(m); err != nil {
			t.Fatalf("goi %d: loi ghi: %v", i, err)
		}
	}

	// Ben client: chi doc, mo phong dung thuat toan doc cua client that.
	client := NewFramer(bytes.NewReader(day.Bytes()), &bytes.Buffer{})
	client.SetKey(khoa)

	for i, mong := range goiGui {
		nhan, err := docGoiNhuClientThat(client)
		if err != nil {
			t.Fatalf("goi %d: loi doc: %v", i, err)
		}
		if nhan.Command != mong.Command {
			t.Fatalf("goi %d: lenh = %d, mong doi %d", i, nhan.Command, mong.Command)
		}
		if !bytes.Equal(nhan.Data, mong.Data) && !(len(nhan.Data) == 0 && len(mong.Data) == 0) {
			t.Fatalf("goi %d: than tin dai %d, mong doi %d", i, len(nhan.Data), len(mong.Data))
		}
	}
}

// docGoiNhuClientThat doc mot goi tin dung y het thuat toan cua client that
// (Session_ME.cs, readMessage/readMessage2): doc 1 byte lenh, roi CAN CU VAO
// LENH (khong can cu vao than tin) de quyet dinh doc 2 hay 3 byte do dai.
//
// Day KHONG phai Framer.ReadMessage() — ham do mo hinh chieu nguoc lai
// (client->server) va khong ho tro lenh dac biet, xem chu thich o
// TestMayChuGhiClientDocDungThuTu o tren.
func docGoiNhuClientThat(f *Framer) (*Message, error) {
	cmdRaw, err := f.docMotByte()
	if err != nil {
		return nil, err
	}
	cmd := int8(f.ks.Read(cmdRaw))

	var n int
	if lenhDoDaiBaByte[cmd] {
		var doDai [3]byte
		for i := range doDai {
			b, err := f.docMotByte()
			if err != nil {
				return nil, err
			}
			doDai[i] = f.ks.Read(b) + 128
		}
		n = int(doDai[0]) | int(doDai[1])<<8 | int(doDai[2])<<16
	} else {
		b1, err := f.docMotByte()
		if err != nil {
			return nil, err
		}
		b2, err := f.docMotByte()
		if err != nil {
			return nil, err
		}
		hi := f.ks.Read(b1)
		lo := f.ks.Read(b2)
		n = int(hi)<<8 | int(lo)
	}

	data := make([]byte, n)
	if n > 0 {
		if _, err := io.ReadFull(f.r, data); err != nil {
			return nil, err
		}
		for i := range data {
			data[i] = f.ks.Read(data[i])
		}
	}
	return &Message{Command: cmd, Data: data}, nil
}
