package chanip

import "sync"

// ToiDaKetNoiMoiIPMacDinh la so ket noi cung luc toi da cho MOT dia chi IP.
//
// VI SAO can gioi han: mot ket noi ton mot goroutine, mot bo dem doc/ghi va
// mot hang goi tin. Khong gioi han thi mot may duy nhat mo vai chuc nghin ket
// noi la an het bo nho may chu, khong can biet mot chut nao ve luat choi —
// re, de, va lam ca may chu chet chum.
//
// VI SAO chon 32 chu khong phai 2 hay 3: nguoi choi Viet Nam thuong xuyen mo
// nhieu nick cung luc, ca nha dung chung mot duong mang, quan net thi ca chuc
// may ra ngoai bang MOT dia chi IP. Con so nay khong nham "moi nguoi mot ket
// noi" — no chi nham chan cai canh MOT may mo hang van ket noi. 32 du rong de
// mot ho gia dinh hay mot quan net nho khong bao gio cham toi, ma van chan
// dut khoat kieu tan cong ke tren.
//
// DIEM YEU DA BIET, ghi ra day cho ro: nha mang di dong Viet Nam dung NAT
// lon, ca ngan thue bao co the chung mot dia chi IP cong cong. Voi mot IP nhu
// vay thi 32 la it. Do la ly do con so nay dat duoc bang bien moi truong
// (NRO_TOIDA_KETNOI_MOI_IP — xem cmd/nroserver/main.go): neu ngay mo cua thay
// nguoi choi bao "khong vao duoc" hang loat tu cung mot nha mang, nguoi van
// hanh noi so nay len duoc ngay, khong phai bien dich lai may chu.
const ToiDaKetNoiMoiIPMacDinh = 32

// GacKetNoi dem so ket noi DANG MO cua tung dia chi IP.
//
// Khac GacDangNhap (dem lan sai, phai tu don rac theo thoi gian), gac nay
// khong can don rac gi ca: so dem ve 0 la o nho bi xoa ngay trong Tra(). Ban
// do vi vay luon chi lon bang so IP DANG THUC SU co ket noi mo.
type GacKetNoi struct {
	mu     sync.Mutex
	theoIP map[string]int
	toiDa  int
}

// MoGacKetNoi tao gac voi han muc toiDa ket noi cung luc cho mot IP.
//
// toiDa <= 0 roi ve mac dinh — mot cau hinh go nham KHONG duoc lang le bien
// thanh "khong gioi han", vi kieu hong do khong co trieu chung nao cho toi
// dung hom bi tan cong.
func MoGacKetNoi(toiDa int) *GacKetNoi {
	if toiDa <= 0 {
		toiDa = ToiDaKetNoiMoiIPMacDinh
	}
	return &GacKetNoi{theoIP: make(map[string]int), toiDa: toiDa}
}

// Nhan xin mot cho cho mot ket noi moi tu ip. Tra ve false neu IP nay da day.
//
// MOI lan Nhan tra ve true BAT BUOC phai co dung mot lan Tra tuong ung, ke ca
// khi ket noi chet vi loi hay vi panic — quen tra mot lan la cho do mat vinh
// vien, va sau vai gio cai gac se tu khoa nguoi choi that ra ngoai du chang
// con ket noi nao. Cach lam dung: defer Tra(ip) NGAY sau khi Nhan tra ve true.
func (g *GacKetNoi) Nhan(ip string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.theoIP[ip] >= g.toiDa {
		return false
	}
	g.theoIP[ip]++
	return true
}

// Tra tra lai mot cho khi ket noi dong.
func (g *GacKetNoi) Tra(ip string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	con, co := g.theoIP[ip]
	if !co {
		// Tra thua (khong co lan Nhan nao tuong ung). Bo qua thay vi de so dem
		// tut xuong am: mot so am se lam dieu kien ">= toiDa" khong bao gio
		// dung nua, tuc la han muc bi vo hieu hoa vinh vien cho IP do —
		// dung cai ma gac nay sinh ra de chan.
		return
	}
	if con <= 1 {
		delete(g.theoIP, ip)
		return
	}
	g.theoIP[ip] = con - 1
}

// SoIPDangGiu cho biet ban do dang giu bao nhieu muc (dung cho test va soi
// ro ri bo nho).
func (g *GacKetNoi) SoIPDangGiu() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return len(g.theoIP)
}

// SoKetNoi cho biet mot IP dang mo bao nhieu ket noi.
func (g *GacKetNoi) SoKetNoi(ip string) int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.theoIP[ip]
}
