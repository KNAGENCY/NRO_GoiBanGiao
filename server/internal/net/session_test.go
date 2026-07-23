package net

import (
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
	go PhucVu(ben_server, xuLy)

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

	go PhucVu(ben_server, &xuLyGhiNhan{nhanDuoc: make(chan *protocol.Message, 1)})

	fClient := protocol.NewFramer(ben_client, ben_client)
	if err := fClient.WriteMessage(&protocol.Message{Command: 99}); err != nil {
		t.Fatalf("gui that bai: %v", err)
	}

	ben_client.SetReadDeadline(time.Now().Add(2 * time.Second))
	if _, err := fClient.ReadMessage(); err == nil {
		t.Fatal("may chu phai ngat ket noi khi goi dau khong phai -27")
	}
}
