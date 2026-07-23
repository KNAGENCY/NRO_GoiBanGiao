package net

import (
	"errors"
	"io"
	"net"
	"testing"
	"time"

	"nroserver/internal/protocol"
)

type xuLyGhiNhan struct {
	nhanDuoc chan *protocol.Message
}

func (x *xuLyGhiNhan) XuLy(s *Session, m *protocol.Message) error {
	x.nhanDuoc <- m
	return nil
}

// Kich ban day du: client gui -27, may chu tra khoa, roi ca hai noi chuyen
// bang goi da ma hoa.
func TestBatTayRoiNhanGoiDaMaHoa(t *testing.T) {
	ben_client, ben_server := net.Pipe()
	defer ben_client.Close()

	xuLy := &xuLyGhiNhan{nhanDuoc: make(chan *protocol.Message, 4)}
	go PhucVu(ben_server, func() XuLyGoiTin { return xuLy })

	fClient := protocol.NewFramer(ben_client, ben_client)

	// 1. Client gui goi bat tay
	if err := fClient.WriteMessage(&protocol.Message{Command: -27}); err != nil {
		t.Fatalf("gui goi bat tay that bai: %v", err)
	}

	// 2. Client nhan khoa. Doc THU CONG dung nhu client that (Session_ME.cs,
	// readMessage), KHONG dung fClient.ReadMessage(): ReadMessage() cua Framer
	// mo phong chieu client->server (do dai doc theo little-endian), con day
	// la chieu server->client, ma client that doc do dai theo BIG-ENDIAN
	// (b1<<8 | b2). Hai chieu von dung hai thu tu byte khac nhau, dung
	// ReadMessage() cho chieu nay se hieu sai do dai (xem chu thich o
	// protocol/frame.go, kieu Framer).
	ben_client.SetReadDeadline(time.Now().Add(3 * time.Second))

	var cmdByte [1]byte
	if _, err := io.ReadFull(ben_client, cmdByte[:]); err != nil {
		t.Fatalf("khong doc duoc lenh goi khoa: %v", err)
	}
	cmdTraLoi := int8(cmdByte[0])
	if cmdTraLoi != -27 {
		t.Fatalf("lenh tra loi = %d, mong doi -27", cmdTraLoi)
	}

	var doDaiBytes [2]byte
	if _, err := io.ReadFull(ben_client, doDaiBytes[:]); err != nil {
		t.Fatalf("khong doc duoc do dai goi khoa: %v", err)
	}
	soByteThan := int(doDaiBytes[0])<<8 | int(doDaiBytes[1])

	than := make([]byte, soByteThan)
	if _, err := io.ReadFull(ben_client, than); err != nil {
		t.Fatalf("khong doc duoc than tin goi khoa (%d byte): %v", soByteThan, err)
	}

	r := protocol.NewReader(than)
	doDai, err := r.Byte()
	if err != nil {
		t.Fatalf("doc do dai khoa that bai: %v", err)
	}
	if doDai < 1 || doDai > 127 {
		t.Fatalf("do dai khoa = %d, phai trong 1..127", doDai)
	}
	khoaTho := make([]byte, doDai)
	for i := range khoaTho {
		khoaTho[i], _ = r.Byte()
	}
	if _, err := r.UTF(); err != nil {
		t.Fatalf("thieu truong dia chi server phu: %v", err)
	}
	if _, err := r.Int(); err != nil {
		t.Fatalf("thieu truong cong server phu: %v", err)
	}
	if _, err := r.Byte(); err != nil {
		t.Fatalf("thieu truong co server phu: %v", err)
	}

	// 3. Ca hai bat ma hoa
	fClient.SetKey(protocol.DeriveKey(khoaTho))

	// 4. Client gui goi da ma hoa, may chu phai giai ma dung
	if err := fClient.WriteMessage(&protocol.Message{Command: -29, Data: []byte("test")}); err != nil {
		t.Fatalf("gui goi da ma hoa that bai: %v", err)
	}

	select {
	case m := <-xuLy.nhanDuoc:
		if m.Command != -29 {
			t.Fatalf("lenh = %d, mong doi -29", m.Command)
		}
		if string(m.Data) != "test" {
			t.Fatalf("than tin = %q, mong doi \"test\"", m.Data)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("may chu khong giai ma duoc goi sau bat tay")
	}
}

// Goi dau tien khong phai -27 thi phai bi ngat ket noi.
func TestGoiDauKhongPhaiBatTayThiNgat(t *testing.T) {
	ben_client, ben_server := net.Pipe()
	defer ben_client.Close()

	xuLy2 := &xuLyGhiNhan{nhanDuoc: make(chan *protocol.Message, 1)}
	go PhucVu(ben_server, func() XuLyGoiTin { return xuLy2 })

	fClient := protocol.NewFramer(ben_client, ben_client)
	if err := fClient.WriteMessage(&protocol.Message{Command: 99}); err != nil {
		t.Fatalf("gui that bai: %v", err)
	}

	ben_client.SetReadDeadline(time.Now().Add(2 * time.Second))
	if _, err := fClient.ReadMessage(); err == nil {
		t.Fatal("may chu phai ngat ket noi khi goi dau khong phai -27")
	}
}

// ---------------------------------------------------------------------------
// PHAN A ③: Gui phai co han gio ghi, khong chi han gio doc.
// ---------------------------------------------------------------------------

// connHanGioNgan bien moi loi goi SetWriteDeadline thanh mot han RAT NGAN
// (50ms), bat ke Session.Gui yeu cau han gio bao lau. Nho vay test co the
// chung minh co che han gio ghi thuc su phat huy tac dung — Write() tren
// net.Pipe tra ve loi timeout — ma khong phai cho that su het thoiGianChoGui
// (vd 30 giay) trong luc chay test.
type connHanGioNgan struct {
	net.Conn
}

func (c *connHanGioNgan) SetWriteDeadline(t time.Time) error {
	return c.Conn.SetWriteDeadline(time.Now().Add(50 * time.Millisecond))
}

// TestGuiCoHanGioGhiKhongTreoVinhVien la bang chung cho viec ③ trong PHAN A.
//
// Boi canh (do bang thuc te): client gui 10.910 goi tin roi khong doc gi ca.
// Session.Gui() truoc day KHONG dat han gio ghi, nen khi bo dem gui TCP day
// (client khong doc), luong may chu goi Write() se KET VINH VIEN o do. Han
// gio DOC (thoiGianChoGoiTin, 120s) khong cuu duoc: luc do luong dang ket o
// Write(), khong phai o Read(). Ket qua do duoc: ket noi song toi 190 giay.
//
// net.Pipe() khong co bo dem: Write() o mot dau CHAN cho toi khi dau kia Read,
// nen neu khong dat han gio ghi, khong ai doc tu ben_client -> Gui() se treo
// vinh vien. Test dung connHanGioNgan de ep han gio thuc te xuong 50ms (xem
// chu thich kieu do) thay vi phai cho that thoiGianChoGui that (vd 30s).
func TestGuiCoHanGioGhiKhongTreoVinhVien(t *testing.T) {
	ben_client, ben_server := net.Pipe()
	defer ben_client.Close()
	defer ben_server.Close()
	// CO Y khong doc gi tu ben_client: mo phong client gui roi bo doc.

	s := &Session{
		conn:   &connHanGioNgan{Conn: ben_server},
		framer: protocol.NewFramer(ben_server, ben_server),
	}

	xong := make(chan error, 1)
	go func() { xong <- s.Gui(&protocol.Message{Command: -27}) }()

	select {
	case err := <-xong:
		if err == nil {
			t.Fatal("Gui phai tra ve loi timeout khi client khong doc gi (han gio ghi), nhung lai thanh cong")
		}
		var netErr net.Error
		if !errors.As(err, &netErr) || !netErr.Timeout() {
			t.Fatalf("loi = %v, mong doi mot net.Error co Timeout() == true (tuc la Gui co dat han gio ghi)", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Gui() TREO VINH VIEN khi client khong doc gi — thieu han gio ghi truoc moi lan ghi")
	}
}
