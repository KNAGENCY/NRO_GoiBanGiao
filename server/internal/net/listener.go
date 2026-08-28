package net

import (
	"errors"
	"fmt"
	"log/slog"
	"net"

	"nroserver/internal/chanip"
)

// LangNghe mo cong TCP roi phuc vu tung ket noi tren mot goroutine rieng.
// taoXuLy duoc goi mot lan cho MOI ket noi (xem PhucVu).
//
// gac gioi han so ket noi cung luc tren mot dia chi IP — xem PhucVuTren.
func LangNghe(diaChi string, gac *chanip.GacKetNoi, taoXuLy func() XuLyGoiTin) error {
	ln, err := net.Listen("tcp", diaChi)
	if err != nil {
		return fmt.Errorf("khong mo duoc cong %s: %w", diaChi, err)
	}
	defer ln.Close()

	slog.Info("may chu game dang lang nghe", "dia_chi", diaChi)
	return PhucVuTren(ln, gac, taoXuLy)
}

// PhucVuTren chay vong nhan ket noi tren mot listener DA MO san.
//
// Tach khoi LangNghe de kiem thu duoc: test can mo cong tren 127.0.0.1:0 (he
// dieu hanh tu chon so cong ranh) roi doc lai so cong do de quay so vao — voi
// LangNghe (tu mo cong ben trong, khong tra so cong ra) thi khong lam duoc.
//
// gac BAT BUOC khac nil. Khong co no thi mot may duy nhat mo duoc bao nhieu
// ket noi tuy thich, moi ket noi ton mot goroutine va mot bo dem — vai chuc
// nghin ket noi la an het bo nho may chu ma khong can biet mot chut nao ve
// luat choi (no bao mat §7, lo hong 3). Thieu gac thi DUNG HAN thay vi chay
// tiep: mot lop bao ve bi tat ma khong co trieu chung gi la kieu hong te nhat
// — no chi lo ra dung hom bi tan cong. Cung tinh than "hong thi dong" nhu
// buoc ap migration luc khoi dong (cmd/nroserver/main.go).
func PhucVuTren(ln net.Listener, gac *chanip.GacKetNoi, taoXuLy func() XuLyGoiTin) error {
	if gac == nil {
		return errors.New("thieu gac gioi han ket noi moi IP — tu choi mo cong")
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			slog.Error("nhan ket noi that bai", "loi", err)
			continue
		}

		ip := ipCua(conn)
		if !gac.Nhan(ip) {
			// Dong NGAY, truoc ca buoc bat tay: mot ket noi bi tu choi khong
			// duoc phep ton them mot goroutine hay mot bo dem nao.
			//
			// KHONG ghi nhat ky o day: dung cai may dang mo ket noi hang loat
			// la cai sinh ra dong log nay, nen moi dong log lai la mot cach
			// no lam day o dia may chu — dung bai hoc da ghi o dau main.go.
			// Dau hieu cho nguoi van hanh la o cho khac: so ket noi dang mo.
			_ = conn.Close()
			continue
		}

		go func() {
			// Tra cho DU CO CHUYEN GI XAY RA: het loi, het panic, het dong
			// binh thuong. Quen tra mot lan la cho do mat vinh vien, va sau
			// vai gio cai gac se tu khoa nguoi choi that ra ngoai du chang
			// con ket noi nao dang mo. PhucVu tu no da co luoi chan panic ben
			// trong, nhung defer o day van la chot cuoi cung dung cach.
			defer gac.Tra(ip)
			PhucVu(conn, taoXuLy)
		}()
	}
}

// ipCua lay dia chi IP cua dau ben kia, bo phan so cong.
//
// PHAI bo so cong: moi ket noi tu cung mot may co mot so cong khac nhau, nen
// neu dem theo ca "IP:cong" thi moi ket noi la mot khoa rieng va han muc
// khong bao gio cham toi — gac se co ma nhu khong.
func ipCua(conn net.Conn) string {
	if a, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		return a.IP.String()
	}
	if host, _, err := net.SplitHostPort(conn.RemoteAddr().String()); err == nil {
		return host
	}
	return "khong-ro"
}
