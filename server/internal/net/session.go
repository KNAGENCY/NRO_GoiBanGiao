// Package net quan ly vong doi mot ket noi TCP.
//
// Package nay khong biet gi ve luat choi. No chi lo: bat tay, doc goi tin,
// va giao goi tin cho tang tren xu ly.
package net

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"nroserver/internal/protocol"
)

const (
	// doDaiKhoa 8 byte la du. Client gioi han 127, nhung khoa dai khong lam
	// an toan hon dang ke — day chi la lop xao tron, khong phai ma hoa that.
	doDaiKhoa = 8

	// thoiGianChoBatTay: ke tan cong mo hang loat ket noi roi khong gui gi
	// se an het bo nho. Khong bat tay trong 5 giay thi ngat.
	thoiGianChoBatTay = 5 * time.Second

	// thoiGianChoGoiTin: ket noi im lang qua lau coi nhu da chet.
	thoiGianChoGoiTin = 120 * time.Second
)

// XuLyGoiTin la thu tang tren cai dat de nhan goi tin da giai ma.
type XuLyGoiTin interface {
	XuLy(s *Session, m *protocol.Message) error
}

// Session la mot ket noi dang song.
type Session struct {
	conn     net.Conn
	framer   *protocol.Framer
	DiaChiIP string

	// muGui bao dam hai goroutine khong ghi chen nhau vao cung mot ket noi.
	// Luong khoa ma co con tro chay tuan tu — ghi chen nhau la lech khoa
	// vinh vien, va trieu chung la client nhan rac chu khong bao loi.
	muGui sync.Mutex

	dongMotLan sync.Once
}

// Gui gui mot goi tin xuong client. An toan khi goi tu nhieu goroutine.
func (s *Session) Gui(m *protocol.Message) error {
	s.muGui.Lock()
	defer s.muGui.Unlock()
	return s.framer.WriteMessage(m)
}

func (s *Session) Dong() {
	s.dongMotLan.Do(func() { s.conn.Close() })
}

// PhucVu chay vong doi mot ket noi cho toi khi dong.
func PhucVu(conn net.Conn, xuLy XuLyGoiTin) {
	diaChi := "khong-ro"
	if a, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		diaChi = a.IP.String()
	}

	s := &Session{
		conn:     conn,
		framer:   protocol.NewFramer(conn, conn),
		DiaChiIP: diaChi,
	}
	defer s.Dong()

	if err := batTay(s); err != nil {
		slog.Warn("bat tay that bai", "ip", diaChi, "loi", err)
		return
	}

	for {
		_ = conn.SetReadDeadline(time.Now().Add(thoiGianChoGoiTin))
		m, err := s.framer.ReadMessage()
		if err != nil {
			slog.Debug("ket noi ket thuc", "ip", diaChi, "loi", err)
			return
		}
		if err := xuLy.XuLy(s, m); err != nil {
			slog.Warn("xu ly goi tin that bai", "ip", diaChi, "lenh", m.Command, "loi", err)
			return
		}
	}
}

// batTay thuc hien buoc trao khoa. Client mo ket noi la gui ngay goi -27
// rong; may chu tra ve khoa tho roi ca hai bat ma hoa.
func batTay(s *Session) error {
	_ = s.conn.SetReadDeadline(time.Now().Add(thoiGianChoBatTay))

	dau, err := s.framer.ReadMessage()
	if err != nil {
		return fmt.Errorf("khong doc duoc goi dau tien: %w", err)
	}
	if dau.Command != -27 {
		return fmt.Errorf("goi dau tien co lenh %d, mong doi -27", dau.Command)
	}

	khoaTho, err := protocol.GenerateRawKey(doDaiKhoa)
	if err != nil {
		return fmt.Errorf("sinh khoa that bai: %w", err)
	}

	w := protocol.NewWriter()
	w.Byte(byte(len(khoaTho)))
	for _, b := range khoaTho {
		w.Byte(b)
	}
	// Ba truong sau la dia chi may chu phu. Du an nay chi chay mot may chu,
	// nhung client van doc du ba truong — thieu la no ket o vong doc.
	if err := w.UTF(""); err != nil {
		return err
	}
	w.Int(0)
	w.Byte(0)

	if err := s.Gui(protocol.NewMessage(-27, w)); err != nil {
		return fmt.Errorf("gui khoa that bai: %w", err)
	}

	// Bat ma hoa SAU khi da gui goi khoa: chinh goi do phai di nguyen ban.
	s.framer.SetKey(protocol.DeriveKey(khoaTho))
	return nil
}
