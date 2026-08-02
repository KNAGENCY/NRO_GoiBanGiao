// Package net quan ly vong doi mot ket noi TCP.
//
// Package nay khong biet gi ve luat choi. No chi lo: bat tay, doc goi tin,
// va giao goi tin cho tang tren xu ly.
package net

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
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

	// thoiGianChoGui: han gio cho MOI LAN GHI xuong client.
	//
	// Vi sao can ca hai chieu (doc VA ghi), khong chi doc: han gio doc chi
	// cuu duoc khi luong may chu dang KET O READ(). Neu client gui goi lien
	// tuc roi ngung DOC (khong phai ngung gui), bo dem gui TCP phia may chu
	// se day len, va Write() se CHAN VINH VIEN cho toi khi client doc bot —
	// luc do luong dang ket o Write(), han gio doc khong lien quan gi ca. Da
	// do duoc tren may chu that: client gui 10.910 goi roi khong doc gi, ket
	// noi song toi 190 giay truoc khi bi phat hien — han gio doc (120s)
	// khong cuu duoc vi no chi ap dung cho Read(), khong ap dung cho Write().
	//
	// Chon 30s (ngan hon nhieu so thoiGianChoGoiTin): ghi la thao tac may chu
	// CHU DONG thuc hien de tra loi, khong co ly do gi phai cho lau nhu doc
	// (noi may chu bi dong THU DONG cho client goi tin lien tuc).
	thoiGianChoGui = 30 * time.Second

	// sucChuaHangGui la so goi gameplay toi da cho mot client. Actor khong
	// duoc cho socket; vuot nguong nay thi loai client cham.
	sucChuaHangGui = 256

	// thoiGianChoWriter dung cho cleanup cua vong phuc vu. DongCham da dong
	// socket truoc, nen writer binh thuong se thoat ngay; gioi han nay chan
	// mot net.Conn loi khong ton trong production lam cleanup treo vo han.
	thoiGianChoWriter = time.Second
)

var (
	ErrHangGuiDay  = errors.New("hang gui da day")
	ErrPhienDaDong = errors.New("phien da dong")
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

	khoiTaoWriterMotLan sync.Once
	batDauWriterMotLan  sync.Once
	muHang              sync.RWMutex
	hangGui             chan *protocol.Message
	dungWriter          chan struct{}
	writerXong          chan struct{}
	writerDaBatDau      atomic.Bool
}

// Gui gui mot goi tin xuong client. An toan khi goi tu nhieu goroutine.
func (s *Session) Gui(m *protocol.Message) error {
	s.muGui.Lock()
	defer s.muGui.Unlock()
	// Dat han gio TRUOC MOI LAN GHI (xem chu thich thoiGianChoGui): thieu
	// buoc nay la ly do luong may chu co the ket vinh vien o Write() khi
	// client ngung doc du van con gui goi len.
	_ = s.conn.SetWriteDeadline(time.Now().Add(thoiGianChoGui))
	return s.framer.WriteMessage(m)
}

// XepGui xep mot goi gameplay cho writer duy nhat cua Session. Ham chi sao
// chep message va day vao queue, khong bao gio cham socket hay cho writer.
func (s *Session) XepGui(m *protocol.Message) error {
	if m == nil {
		return errors.New("khong the xep message nil")
	}
	s.khoiTaoWriter()
	s.muHang.RLock()
	defer s.muHang.RUnlock()
	select {
	case <-s.dungWriter:
		return ErrPhienDaDong
	default:
	}

	s.batDauWriterMotLan.Do(func() {
		s.writerDaBatDau.Store(true)
		go s.vongWriter()
	})

	ban := &protocol.Message{Command: m.Command}
	if len(m.Data) > 0 {
		ban.Data = append([]byte(nil), m.Data...)
	}
	select {
	case <-s.dungWriter:
		return ErrPhienDaDong
	case s.hangGui <- ban:
		return nil
	default:
		return ErrHangGuiDay
	}
}

func (s *Session) khoiTaoWriter() {
	s.khoiTaoWriterMotLan.Do(func() {
		s.hangGui = make(chan *protocol.Message, sucChuaHangGui)
		s.dungWriter = make(chan struct{})
		s.writerXong = make(chan struct{})
	})
}

func (s *Session) vongWriter() {
	defer close(s.writerXong)
	for {
		// Uu tien tin hieu dong truoc queue: sau write error/DongCham, phan
		// queue con lai phai bi bo, khong duoc ghi tiep.
		select {
		case <-s.dungWriter:
			return
		default:
		}
		select {
		case <-s.dungWriter:
			return
		case m := <-s.hangGui:
			if err := s.Gui(m); err != nil {
				slog.Debug("writer session ket thuc", "ip", s.DiaChiIP, "loi", err)
				s.DongCham()
				return
			}
		}
	}
}

// DongCham chi danh dau/dong socket, tuyet doi khong cho writer. Actor goi
// ham nay khi queue day de khong bi client cham chan vong gameplay.
func (s *Session) DongCham() {
	s.khoiTaoWriter()
	s.dongMotLan.Do(func() {
		s.muHang.Lock()
		close(s.dungWriter)
		s.muHang.Unlock()
		_ = s.conn.Close()
	})
}

func (s *Session) Dong() {
	s.DongCham()
	if !s.writerDaBatDau.Load() {
		return
	}
	select {
	case <-s.writerXong:
	case <-time.After(thoiGianChoWriter):
		slog.Warn("cho writer session thoat qua han", "ip", s.DiaChiIP)
	}
}

// PhucVu chay vong doi mot ket noi cho toi khi dong.
//
// Nhan MOT HAM TAO bo xu ly, khong phai mot bo xu ly dung chung: moi ket noi
// can trang thai rieng (dang nhap bang tai khoan nao, dang giu nhan vat nao).
// Cho moi ket noi mot bo xu ly rieng nghia la trang thai do chi la truong cua
// bo xu ly — khong can ban do phien->trang thai, va khi ket noi dong thi bo xu
// ly duoc thu gom theo, khong ro ri.
func PhucVu(conn net.Conn, taoXuLy func() XuLyGoiTin) {
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

	xuLy := taoXuLy()

	// Neu bo xu ly co diem moc "khi dong" (KH5: bao actor khu vuc go nguoi choi),
	// goi no khi phien ket thuc. Dung type-assert de KHONG bat buoc moi XuLyGoiTin
	// phai co KhiDong (mock test cu khong can doi).
	if c, ok := xuLy.(interface{ KhiDong() }); ok {
		defer c.KhiDong()
	}

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
