package chanip

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// dongHoGia cho phep test tua thoi gian di toi ma khong phai ngoi cho that.
// Bat buoc phai co: mot the hoi phuc sau 20 giay, test khong the sleep 20 giay.
type dongHoGia struct {
	mu  sync.Mutex
	nay time.Time
}

func (d *dongHoGia) Now() time.Time {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.nay
}

func (d *dongHoGia) Tien(khoang time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.nay = d.nay.Add(khoang)
}

func gacThu(sucChua int, nhip time.Duration) (*GacDangNhap, *dongHoGia) {
	dh := &dongHoGia{nay: time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)}
	g := MoGacDangNhap(sucChua, nhip)
	g.dongHo = dh.Now
	return g, dh
}

// Nguoi choi that go nham mat khau vai lan la chuyen binh thuong — dung so lan
// trong suc chua thi TUYET DOI khong duoc chan ho.
func TestNguoiChoiGoNhamTrongSucChuaKhongBiChan(t *testing.T) {
	g, _ := gacThu(10, 20*time.Second)
	for i := 1; i <= 10; i++ {
		if !g.ChoPhep("1.2.3.4") {
			t.Fatalf("lan thu %d da bi chan, nhung suc chua la 10 — nguoi choi that go nham "+
				"trong nguong nay khong duoc bi da ra", i)
		}
		g.GhiNhanSai("1.2.3.4")
	}
}

// Qua suc chua thi phai chan — day la chot chinh cua lo hong 1.
func TestQuaSucChuaThiBiChan(t *testing.T) {
	g, _ := gacThu(10, 20*time.Second)
	for i := 0; i < 10; i++ {
		g.ChoPhep("1.2.3.4")
		g.GhiNhanSai("1.2.3.4")
	}
	if g.ChoPhep("1.2.3.4") {
		t.Fatal("sau 10 lan sai lien tiep, IP nay VAN duoc thu tiep — khong co gi chan do mat khau")
	}
}

// Chan phai theo TUNG IP: mot ke xau do mat khau khong duoc lam nguoi choi
// khac (IP khac) dang nhap khong duoc.
func TestChanTheoTungIPKhongLanSangIPKhac(t *testing.T) {
	g, _ := gacThu(10, 20*time.Second)
	for i := 0; i < 50; i++ {
		g.ChoPhep("1.2.3.4")
		g.GhiNhanSai("1.2.3.4")
	}
	if g.ChoPhep("1.2.3.4") {
		t.Fatal("IP do mat khau dang le phai bi chan")
	}
	if !g.ChoPhep("5.6.7.8") {
		t.Fatal("IP khac bi va lay — chan phai theo tung IP, khong duoc lam ca may chu chet chum")
	}
}

// The phai tu day lai theo thoi gian, neu khong thi mot lan bi chan la bi
// chan vinh vien.
func TestTheDayLaiTheoThoiGian(t *testing.T) {
	g, dh := gacThu(10, 20*time.Second)
	for i := 0; i < 10; i++ {
		g.ChoPhep("1.2.3.4")
		g.GhiNhanSai("1.2.3.4")
	}
	if g.ChoPhep("1.2.3.4") {
		t.Fatal("dang le da het the")
	}

	// Chua du mot nhip thi van con bi chan.
	dh.Tien(19 * time.Second)
	if g.ChoPhep("1.2.3.4") {
		t.Fatal("moi 19 giay ma da hoi duoc mot the — nhip hoi phuc la 20 giay")
	}

	// Du mot nhip thi duoc them DUNG mot lan thu.
	dh.Tien(2 * time.Second)
	if !g.ChoPhep("1.2.3.4") {
		t.Fatal("qua 20 giay ma van chua hoi duoc the nao — nguoi choi that se bi khoa vinh vien")
	}
	g.GhiNhanSai("1.2.3.4")
	if g.ChoPhep("1.2.3.4") {
		t.Fatal("hoi mot the ma tieu duoc hai lan — cong thuc hoi phuc sai")
	}
}

// The khong duoc phep day qua suc chua: nghi mot ngay roi quay lai KHONG
// duoc doi thanh "duoc do thoai mai mot dot khong lo".
func TestTheKhongDayVuotSucChua(t *testing.T) {
	g, dh := gacThu(10, 20*time.Second)
	g.ChoPhep("1.2.3.4")
	g.GhiNhanSai("1.2.3.4")
	dh.Tien(24 * time.Hour)

	dem := 0
	for g.ChoPhep("1.2.3.4") {
		g.GhiNhanSai("1.2.3.4")
		dem++
		if dem > 100 {
			break
		}
	}
	if dem != 10 {
		t.Fatalf("sau khi nghi that lau, IP nay thu duoc %d lan lien tiep — chi duoc dung suc chua 10", dem)
	}
}

// Ban do IP KHONG duoc phinh vo han: moi IP mot o nho, ke xau doi IP lien tuc
// (IPv6 thi vo han) se an het bo nho may chu — lai la mot duong lam sap may chu.
func TestBanDoIPTuDonRacKhongPhinhVoHan(t *testing.T) {
	g, dh := gacThu(10, 20*time.Second)
	for i := 0; i < 5000; i++ {
		ip := fmt.Sprintf("10.0.%d.%d", i/256, i%256)
		g.ChoPhep(ip)
		g.GhiNhanSai(ip)
	}
	// Qua thoi gian du de MOI IP tren hoi day the — luc do o nho cua chung
	// khong con mang thong tin gi (giong het mot IP chua tung thay), nen phai
	// duoc don di.
	dh.Tien(10 * 20 * time.Second)
	g.ChoPhep("1.1.1.1") // mot thao tac bat ky de kich hoat don rac

	if con := g.SoIPDangNho(); con > 100 {
		t.Fatalf("con nho %d IP sau khi tat ca da hoi day the — ban do khong duoc don rac, "+
			"ke xau doi IP lien tuc se an het bo nho may chu", con)
	}
}

// Gac duoc goi tu MOI ket noi, tuc tu rat nhieu goroutine cung luc. Chay
// -race de bat dua tranh du lieu.
func TestGacAnToanKhiNhieuGoroutineCungGoi(t *testing.T) {
	g, _ := gacThu(10, 20*time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			ip := fmt.Sprintf("10.0.0.%d", n%7)
			for j := 0; j < 200; j++ {
				g.ChoPhep(ip)
				g.GhiNhanSai(ip)
			}
		}(i)
	}
	wg.Wait()
}
