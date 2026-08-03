package net

import (
	"io"
	"net"
	"testing"
	"time"

	"nroserver/internal/protocol"
)

// xuLyPanicGia panic ngay khi xu goi tin dau tien — diem chen test luoi panic
// (manh A Task 4): mo phong mot bug trong tang xu ly goi tin cho mot ket noi.
type xuLyPanicGia struct{}

func (xuLyPanicGia) XuLy(*Session, *protocol.Message) error {
	panic("xu ly goi tin panic gia")
}

// TestPhucVu_XuLyPanic_DongKetNoiKhongLamSap (manh A §3.2b): mot ket noi lam
// XuLy panic phai bi DONG SACH (PhucVu thoat, socket dong) MA KHONG lam sap ca
// tien trinh — ket noi khac van phuc vu binh thuong.
func TestPhucVu_XuLyPanic_DongKetNoiKhongLamSap(t *testing.T) {
	client, serverGoc := net.Pipe()
	defer client.Close()

	phucVuXong := make(chan struct{})
	go func() {
		PhucVu(serverGoc, func() XuLyGoiTin { return xuLyPanicGia{} })
		close(phucVuXong)
	}()

	f := batTayClientTest(t, client)
	// Gui mot goi bat ky -> server doc -> XuLy panic.
	if err := f.WriteMessage(&protocol.Message{Command: 1}); err != nil {
		t.Fatalf("gui goi kich panic: %v", err)
	}

	// PhucVu PHAI thoat sach (khong lam test chet vi panic lan ra).
	select {
	case <-phucVuXong:
	case <-time.After(2 * time.Second):
		t.Fatal("PhucVu khong thoat sau khi XuLy panic — luoi khong chan")
	}

	// Ket noi bi dong: client doc phai tra loi (EOF/dong).
	_ = client.SetReadDeadline(time.Now().Add(time.Second))
	if _, err := io.ReadFull(client, make([]byte, 1)); err == nil {
		t.Fatal("ket noi gap panic phai duoc dong (client doc phai loi)")
	}

	// Ket noi THU HAI voi handler binh thuong phai bat tay va trao goi duoc —
	// bang chung tien trinh van song sau panic o ket noi truoc.
	client2, serverGoc2 := net.Pipe()
	defer client2.Close()
	xuLy2 := &xuLyGhiNhan{nhanDuoc: make(chan *protocol.Message, 1)}
	go PhucVu(serverGoc2, func() XuLyGoiTin { return xuLy2 })
	f2 := batTayClientTest(t, client2)
	if err := f2.WriteMessage(&protocol.Message{Command: -29, Data: []byte("sau-panic")}); err != nil {
		t.Fatalf("ket noi 2 gui goi: %v", err)
	}
	select {
	case m := <-xuLy2.nhanDuoc:
		if string(m.Data) != "sau-panic" {
			t.Fatalf("ket noi 2 nhan than %q, mong 'sau-panic'", m.Data)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("ket noi 2 khong phuc vu duoc — tien trinh khong con song sau panic")
	}
}

// connGhiPanic la net.Conn gia panic NGAY khi Write — diem chen test luoi
// vongWriter (manh A Task 4): mo phong loi la luc ghi socket cho mot ket noi.
type connGhiPanic struct{}

func (connGhiPanic) Read([]byte) (int, error)         { return 0, io.EOF }
func (connGhiPanic) Write([]byte) (int, error)        { panic("ghi socket panic gia") }
func (connGhiPanic) Close() error                     { return nil }
func (connGhiPanic) LocalAddr() net.Addr              { return diaChiGia("local") }
func (connGhiPanic) RemoteAddr() net.Addr             { return diaChiGia("remote") }
func (connGhiPanic) SetDeadline(time.Time) error      { return nil }
func (connGhiPanic) SetReadDeadline(time.Time) error  { return nil }
func (connGhiPanic) SetWriteDeadline(time.Time) error { return nil }

// TestVongWriter_GhiPanic_ThoatSachKhongLamSap (manh A §3.2b): panic luc GHI
// cho mot ket noi phai bi chan, writer thoat sach (dong writerXong) — khong lam
// sap ca tien trinh.
func TestVongWriter_GhiPanic_ThoatSachKhongLamSap(t *testing.T) {
	conn := connGhiPanic{}
	s := &Session{conn: conn, framer: protocol.NewFramer(conn, conn)}

	if err := s.XepGui(&protocol.Message{Command: 1}); err != nil {
		t.Fatalf("xep goi kich writer: %v", err)
	}

	select {
	case <-s.writerXong:
		// PASS: vongWriter chan panic khi ghi, dong writerXong sach.
	case <-time.After(2 * time.Second):
		t.Fatal("vongWriter khong thoat sau panic khi ghi — luoi khong chan")
	}
}
