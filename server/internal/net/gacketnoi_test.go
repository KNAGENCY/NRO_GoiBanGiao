package net

import (
	"errors"
	"net"
	"testing"
	"time"

	"nroserver/internal/chanip"
	"nroserver/internal/protocol"
)

// Lo hong 3: khong co gi gioi han so ket noi mot may mo duoc. Mot may mo vai
// chuc nghin ket noi la an het bo nho may chu ma khong can biet mot chut gi ve
// luat choi. Cac test o day dung CONG TCP THAT tren 127.0.0.1 (khong phai
// net.Pipe) vi cho can kiem chinh la vong Accept — net.Pipe khong co vong do.

type xuLyRong struct{}

func (xuLyRong) XuLy(*Session, *protocol.Message) error { return nil }

// moMayChuThu mo mot cong tam thoi tren 127.0.0.1 va chay vong phuc vu voi
// gac da cho. Tra ve dia chi de quay so va ham dong.
func moMayChuThu(t *testing.T, gac *chanip.GacKetNoi) (diaChi string, dong func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("khong mo duoc cong thu: %v", err)
	}
	go func() { _ = PhucVuTren(ln, gac, func() XuLyGoiTin { return xuLyRong{} }) }()
	return ln.Addr().String(), func() { _ = ln.Close() }
}

// conSong kiem mot ket noi con mo hay da bi may chu dong ngay lap tuc.
//
// Cach kiem: thu DOC voi han gio ngan. May chu chua gui gi ca, nen mot ket
// noi CON SONG se het han gio (loi timeout) — con mot ket noi DA BI DONG se
// tra ve EOF hoac loi reset ngay.
func conSong(t *testing.T, c net.Conn) bool {
	t.Helper()
	_ = c.SetReadDeadline(time.Now().Add(700 * time.Millisecond))
	var b [1]byte
	_, err := c.Read(b[:])
	if err == nil {
		return true // doc duoc du lieu -> chac chan con song
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true // het han gio ma khong co gi toi -> ket noi van dang mo
	}
	return false // EOF hoac reset -> may chu da dong ket noi nay
}

// choDenKhi doi cho toi khi dieu kien dung, toi da 3 giay. Can vi vong Accept
// chay tren goroutine khac: sau khi quay so xong, may chu co the chua kip ghi
// nhan ket noi do vao gac.
func choDenKhi(t *testing.T, mo string, dk func() bool) {
	t.Helper()
	han := time.Now().Add(3 * time.Second)
	for time.Now().Before(han) {
		if dk() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("cho qua han: %s", mo)
}

// TestMotIPKhongMoDuocQuaHanMucKetNoi la chot cua lo hong 3.
func TestMotIPKhongMoDuocQuaHanMucKetNoi(t *testing.T) {
	const hanMuc = 3
	gac := chanip.MoGacKetNoi(hanMuc)
	diaChi, dongMayChu := moMayChuThu(t, gac)
	defer dongMayChu()

	var daMo []net.Conn
	defer func() {
		for _, c := range daMo {
			_ = c.Close()
		}
	}()

	for i := 1; i <= hanMuc; i++ {
		c, err := net.Dial("tcp", diaChi)
		if err != nil {
			t.Fatalf("ket noi thu %d khong quay so duoc: %v", i, err)
		}
		daMo = append(daMo, c)
		so := i
		choDenKhi(t, "gac ghi nhan du ket noi", func() bool { return gac.SoKetNoi("127.0.0.1") >= so })
	}

	// Ket noi vuot han muc: may chu phai dong no NGAY.
	thua, err := net.Dial("tcp", diaChi)
	if err != nil {
		// Bi tu choi ngay tu buoc quay so cung la ket qua dung.
		return
	}
	defer thua.Close()
	if conSong(t, thua) {
		t.Fatalf("mo duoc ket noi thu %d trong khi han muc la %d — khong co gi chan mot may "+
			"mo ket noi vo han lam can tai nguyen may chu", hanMuc+1, hanMuc)
	}
}

// Khi mot ket noi dong, cho cua no phai duoc tra lai — neu khong, choi vai
// tieng la ca dai IP cua nha mang bi khoa sach du chang con ket noi nao.
func TestDongKetNoiThiTraLaiChoChoIPDo(t *testing.T) {
	const hanMuc = 2
	gac := chanip.MoGacKetNoi(hanMuc)
	diaChi, dongMayChu := moMayChuThu(t, gac)
	defer dongMayChu()

	c1, err := net.Dial("tcp", diaChi)
	if err != nil {
		t.Fatalf("quay so 1 that bai: %v", err)
	}
	c2, err := net.Dial("tcp", diaChi)
	if err != nil {
		t.Fatalf("quay so 2 that bai: %v", err)
	}
	defer c2.Close()
	choDenKhi(t, "gac ghi nhan 2 ket noi", func() bool { return gac.SoKetNoi("127.0.0.1") == 2 })

	_ = c1.Close()
	choDenKhi(t, "gac tra lai cho khi ket noi dong", func() bool { return gac.SoKetNoi("127.0.0.1") == 1 })

	c3, err := net.Dial("tcp", diaChi)
	if err != nil {
		t.Fatalf("quay so lai sau khi da tra cho that bai: %v", err)
	}
	defer c3.Close()
	if !conSong(t, c3) {
		t.Fatal("da dong mot ket noi cu ma van khong mo duoc ket noi moi — gac quen tra cho")
	}
}

// Thieu gac la loi cau hinh NGHIEM TRONG (may chu chay khong co gioi han ket
// noi nao). Phai bao loi ngay luc khoi dong chu khong duoc lang le chay tiep:
// mot lop bao ve bi tat ma khong co trieu chung la kieu hong te nhat.
func TestThieuGacThiKhongChiuChay(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("khong mo duoc cong thu: %v", err)
	}
	defer ln.Close()
	if err := PhucVuTren(ln, nil, func() XuLyGoiTin { return xuLyRong{} }); err == nil {
		t.Fatal("PhucVuTren chay binh thuong voi gac = nil — may chu se mo cong ma khong co " +
			"gioi han ket noi nao, dung lo hong dang phai va")
	}
}
