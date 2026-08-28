package chanip

import (
	"fmt"
	"sync"
	"testing"
)

// Trong han muc thi nhan het — mot nha nhieu nguoi choi, mot quan net, mot
// nguoi mo may nick cung luc deu la chuyen binh thuong.
func TestNhanDuKetNoiTrongHanMuc(t *testing.T) {
	g := MoGacKetNoi(8)
	for i := 1; i <= 8; i++ {
		if !g.Nhan("1.2.3.4") {
			t.Fatalf("ket noi thu %d bi tu choi trong khi han muc la 8", i)
		}
	}
}

// Vuot han muc thi tu choi — day la chot cua lo hong 3.
func TestVuotHanMucThiTuChoi(t *testing.T) {
	g := MoGacKetNoi(8)
	for i := 0; i < 8; i++ {
		g.Nhan("1.2.3.4")
	}
	if g.Nhan("1.2.3.4") {
		t.Fatal("mot IP mo duoc ket noi thu 9 trong khi han muc la 8 — khong co gi chan " +
			"mot may mo ket noi vo han lam can tai nguyen may chu")
	}
}

// Han muc theo TUNG IP: mot ke xau mo kin cho minh khong duoc lam nguoi choi
// khac khong vao duoc.
func TestHanMucTheoTungIP(t *testing.T) {
	g := MoGacKetNoi(2)
	g.Nhan("1.2.3.4")
	g.Nhan("1.2.3.4")
	if g.Nhan("1.2.3.4") {
		t.Fatal("IP nay dang le da day")
	}
	if !g.Nhan("5.6.7.8") {
		t.Fatal("IP khac bi va lay — han muc phai tinh rieng tung IP")
	}
}

// Tra lai cho khi ket noi dong: neu quen tra, choi vai tieng la ca dai IP cua
// nha mang bi khoa sach du chang ai con ket noi nao.
func TestTraChoKhiDongThiNhanLaiDuoc(t *testing.T) {
	g := MoGacKetNoi(2)
	g.Nhan("1.2.3.4")
	g.Nhan("1.2.3.4")
	g.Tra("1.2.3.4")
	if !g.Nhan("1.2.3.4") {
		t.Fatal("da tra lai mot cho ma van khong nhan duoc ket noi moi")
	}
}

// Khi mot IP khong con ket noi nao, o nho cua no phai bien mat khoi ban do —
// neu khong, ban do lon dan theo tong so IP tung ghe qua, khong bao gio nho lai.
func TestVeKhongThiXoaKhoiBanDo(t *testing.T) {
	g := MoGacKetNoi(4)
	for i := 0; i < 3000; i++ {
		ip := fmt.Sprintf("10.0.%d.%d", i/256, i%256)
		g.Nhan(ip)
		g.Tra(ip)
	}
	if con := g.SoIPDangGiu(); con != 0 {
		t.Fatalf("con nho %d IP du moi ket noi deu da dong — ban do khong tu co lai", con)
	}
}

// Tra nhieu hon nhan KHONG duoc lam so dem am: so am se bien han muc thanh vo
// han (am + 1 van < han muc mai mai).
func TestTraThuaKhongLamSoDemAm(t *testing.T) {
	g := MoGacKetNoi(2)
	g.Tra("1.2.3.4")
	g.Tra("1.2.3.4")
	g.Nhan("1.2.3.4")
	g.Nhan("1.2.3.4")
	if g.Nhan("1.2.3.4") {
		t.Fatal("tra thua da lam hong so dem — han muc khong con tac dung")
	}
}

// Nhan/Tra duoc goi tu vong Accept va tu moi goroutine ket noi, nen phai an
// toan khi chay dong thoi. Chay voi -race.
func TestGacKetNoiAnToanDongThoi(t *testing.T) {
	g := MoGacKetNoi(64)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			ip := fmt.Sprintf("10.0.0.%d", n%5)
			for j := 0; j < 500; j++ {
				if g.Nhan(ip) {
					g.Tra(ip)
				}
			}
		}(i)
	}
	wg.Wait()
	if con := g.SoIPDangGiu(); con != 0 {
		t.Fatalf("con %d IP trong ban do sau khi moi ket noi deu da tra — ro ri o nho", con)
	}
}
