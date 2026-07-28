package net

import (
	"bytes"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
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

// connChanLanGhiDau mo phong socket cua mot client khong doc: lan Write dau
// dung lai cho toi khi test mo khoa. Cac lan sau ghi vao bo nho de test doc
// dung thu tu khung tin that, khong khang dinh tren mock call.
type connChanLanGhiDau struct {
	batDau       chan struct{}
	mo           chan struct{}
	daBaoBatDau  sync.Once
	soLanGhi     atomic.Int32
	dangGhi      atomic.Int32
	soGhiDongLuc atomic.Int32
	mu           sync.Mutex
	duLieu       bytes.Buffer
	dong         atomic.Int32
}

func newConnChanLanGhiDau() *connChanLanGhiDau {
	return &connChanLanGhiDau{
		batDau: make(chan struct{}),
		mo:     make(chan struct{}),
	}
}

func (c *connChanLanGhiDau) Read([]byte) (int, error) { return 0, io.EOF }
func (c *connChanLanGhiDau) Write(p []byte) (int, error) {
	if c.dangGhi.Add(1) != 1 {
		c.soGhiDongLuc.Add(1)
	}
	defer c.dangGhi.Add(-1)

	if c.soLanGhi.Add(1) == 1 {
		c.daBaoBatDau.Do(func() { close(c.batDau) })
		<-c.mo
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.duLieu.Write(p)
}
func (c *connChanLanGhiDau) Close() error {
	c.dong.Add(1)
	return nil
}
func (c *connChanLanGhiDau) LocalAddr() net.Addr              { return diaChiGia("local") }
func (c *connChanLanGhiDau) RemoteAddr() net.Addr             { return diaChiGia("remote") }
func (c *connChanLanGhiDau) SetDeadline(time.Time) error      { return nil }
func (c *connChanLanGhiDau) SetReadDeadline(time.Time) error  { return nil }
func (c *connChanLanGhiDau) SetWriteDeadline(time.Time) error { return nil }

type diaChiGia string

func (d diaChiGia) Network() string { return "test" }
func (d diaChiGia) String() string  { return string(d) }

// TestXepGuiKhongChanFIFO256ImmutableMotWriter bat bon loi san xuat:
// XepGui ghi socket truc tiep (SERVER-ACTOR-005), hang sai nguong 256, giu
// con tro Data de ben goi sua duoc goi da xep, hoac tao nhieu writer ghi chen.
func TestXepGuiKhongChanFIFO256ImmutableMotWriter(t *testing.T) {
	conn := newConnChanLanGhiDau()
	s := &Session{conn: conn, framer: protocol.NewFramer(conn, conn)}
	t.Cleanup(s.Dong)

	dau := &protocol.Message{Command: 7, Data: []byte{1}}
	if err := s.XepGui(dau); err != nil {
		t.Fatalf("xep goi dau that bai: %v", err)
	}
	select {
	case <-conn.batDau:
	case <-time.After(time.Second):
		t.Fatal("writer khong nhan goi dau")
	}
	dau.Data[0] = 99 // message da enqueue phai bat bien.

	// Writer dang ket o goi dau, vi vay 256 goi sau nam tron trong queue.
	for i := 0; i < 256; i++ {
		if err := s.XepGui(&protocol.Message{Command: 7, Data: []byte{byte(i + 2)}}); err != nil {
			t.Fatalf("queue chua du 256 nhung goi %d bi tu choi: %v", i, err)
		}
	}
	batDau := time.Now()
	err := s.XepGui(&protocol.Message{Command: 7, Data: []byte{0}})
	if !errors.Is(err, ErrHangGuiDay) {
		t.Fatalf("goi thu 257 trong queue = %v, mong ErrHangGuiDay", err)
	}
	if mat := time.Since(batDau); mat >= 10*time.Millisecond {
		t.Fatalf("XepGui khi queue day mat %s, phai tra duoi 10ms", mat)
	}

	close(conn.mo)
	han := time.After(2 * time.Second)
	for conn.soLanGhi.Load() < 257 {
		select {
		case <-han:
			t.Fatalf("writer moi ghi %d/257 goi", conn.soLanGhi.Load())
		default:
			time.Sleep(time.Millisecond)
		}
	}
	if got := conn.soGhiDongLuc.Load(); got != 0 {
		t.Fatalf("co %d lan Write chay dong thoi, mong dung mot writer", got)
	}

	conn.mu.Lock()
	raw := append([]byte(nil), conn.duLieu.Bytes()...)
	conn.mu.Unlock()
	for i := 0; i < 257; i++ {
		if len(raw) < 4 {
			t.Fatalf("goi %d bi thieu khung, con %d byte", i, len(raw))
		}
		if int8(raw[0]) != 7 || raw[1] != 0 || raw[2] != 1 {
			t.Fatalf("goi %d co header %v, mong [7 0 1]", i, raw[:3])
		}
		want := byte(i + 1)
		if raw[3] != want {
			t.Fatalf("goi %d co data=%d, mong %d (FIFO/immutable bi vo)", i, raw[3], want)
		}
		raw = raw[4:]
	}
	if len(raw) != 0 {
		t.Fatalf("con %d byte rac sau 257 goi", len(raw))
	}
}

type connLoiGhi struct {
	*connChanLanGhiDau
	loi error
}

func (c *connLoiGhi) Write([]byte) (int, error) {
	if c.soLanGhi.Add(1) == 1 {
		c.daBaoBatDau.Do(func() { close(c.batDau) })
		<-c.mo
	}
	return 0, c.loi
}

// TestWriterLoiDongChamBoQueueVaThoat bat dot bien bo DongCham sau Write error:
// socket chi dong mot lan, queue khong duoc ghi tiep va writer phai thoat.
func TestWriterLoiDongChamBoQueueVaThoat(t *testing.T) {
	loiGhi := errors.New("mat ket noi")
	goc := newConnChanLanGhiDau()
	conn := &connLoiGhi{connChanLanGhiDau: goc, loi: loiGhi}
	s := &Session{conn: conn, framer: protocol.NewFramer(conn, conn)}

	if err := s.XepGui(&protocol.Message{Command: 1}); err != nil {
		t.Fatalf("xep goi dau: %v", err)
	}
	<-conn.batDau
	for i := 0; i < 10; i++ {
		if err := s.XepGui(&protocol.Message{Command: int8(i + 2)}); err != nil {
			t.Fatalf("xep goi cho %d: %v", i, err)
		}
	}
	close(conn.mo)

	select {
	case <-s.writerXong:
	case <-time.After(time.Second):
		t.Fatal("writer khong thoat sau Write error")
	}
	if got := conn.soLanGhi.Load(); got != 1 {
		t.Fatalf("writer van ghi %d lan sau loi dau, queue chua bi bo", got)
	}
	if got := conn.dong.Load(); got != 1 {
		t.Fatalf("Close duoc goi %d lan, mong dung 1", got)
	}
	s.DongCham()
	s.DongCham()
	if got := conn.dong.Load(); got != 1 {
		t.Fatalf("DongCham khong idempotent: Close duoc goi %d lan", got)
	}
	if err := s.XepGui(&protocol.Message{Command: 99}); !errors.Is(err, ErrPhienDaDong) {
		t.Fatalf("XepGui sau write error = %v, mong ErrPhienDaDong", err)
	}
}

type connLoiTuLanGhiThuHai struct {
	net.Conn
	soLanGhi  atomic.Int32
	soLanDong atomic.Int32
	loi       error
}

func (c *connLoiTuLanGhiThuHai) Write(p []byte) (int, error) {
	if c.soLanGhi.Add(1) == 2 {
		return 0, c.loi
	}
	return c.Conn.Write(p)
}

func (c *connLoiTuLanGhiThuHai) Close() error {
	c.soLanDong.Add(1)
	return c.Conn.Close()
}

type xuLyXepRoiDemKhiDong struct {
	xepXong   chan error
	soKhiDong atomic.Int32
}

func (x *xuLyXepRoiDemKhiDong) XuLy(s *Session, _ *protocol.Message) error {
	x.xepXong <- s.XepGui(&protocol.Message{Command: 8, Data: []byte{1}})
	return nil
}

func (x *xuLyXepRoiDemKhiDong) KhiDong() { x.soKhiDong.Add(1) }

// batTayClientTest tra Framer da co khoa, dung dung byte-order cua client that.
func batTayClientTest(t *testing.T, conn net.Conn) *protocol.Framer {
	t.Helper()
	f := protocol.NewFramer(conn, conn)
	if err := f.WriteMessage(&protocol.Message{Command: -27}); err != nil {
		t.Fatalf("gui bat tay: %v", err)
	}
	var header [3]byte
	if _, err := io.ReadFull(conn, header[:]); err != nil {
		t.Fatalf("doc header bat tay: %v", err)
	}
	if int8(header[0]) != -27 {
		t.Fatalf("lenh bat tay = %d", int8(header[0]))
	}
	n := int(header[1])<<8 | int(header[2])
	than := make([]byte, n)
	if _, err := io.ReadFull(conn, than); err != nil {
		t.Fatalf("doc than bat tay: %v", err)
	}
	r := protocol.NewReader(than)
	doDai, err := r.Byte()
	if err != nil {
		t.Fatalf("doc do dai khoa: %v", err)
	}
	khoa := make([]byte, doDai)
	for i := range khoa {
		khoa[i], err = r.Byte()
		if err != nil {
			t.Fatalf("doc khoa byte %d: %v", i, err)
		}
	}
	f.SetKey(protocol.DeriveKey(khoa))
	return f
}

// TestWriteErrorLamPhucVuCleanupDungMotLan chung minh write error cua writer
// dong socket de vong Read thoat va moc KhiDong chay dung mot lan.
func TestWriteErrorLamPhucVuCleanupDungMotLan(t *testing.T) {
	client, serverGoc := net.Pipe()
	defer client.Close()
	conn := &connLoiTuLanGhiThuHai{Conn: serverGoc, loi: errors.New("reset")}
	xuLy := &xuLyXepRoiDemKhiDong{xepXong: make(chan error, 1)}
	phucVuXong := make(chan struct{})
	go func() {
		PhucVu(conn, func() XuLyGoiTin { return xuLy })
		close(phucVuXong)
	}()

	f := batTayClientTest(t, client)
	if err := f.WriteMessage(&protocol.Message{Command: 1}); err != nil {
		t.Fatalf("gui gameplay kich writer: %v", err)
	}
	select {
	case err := <-xuLy.xepXong:
		if err != nil {
			t.Fatalf("XepGui trong handler: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("handler khong xep goi gameplay")
	}
	select {
	case <-phucVuXong:
	case <-time.After(time.Second):
		t.Fatal("PhucVu khong thoat sau write error")
	}
	if got := xuLy.soKhiDong.Load(); got != 1 {
		t.Fatalf("KhiDong chay %d lan, mong dung 1", got)
	}
	if got := conn.soLanDong.Load(); got != 1 {
		t.Fatalf("server conn Close %d lan, mong dung 1", got)
	}
}
