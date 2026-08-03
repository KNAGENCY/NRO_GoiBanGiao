package game

import (
	"testing"
	"time"
)

// capNhatTreNhip la ham THUAN (khong dung dong ho): nhan khoang cach giua hai
// nhip, quyet dinh co dem la "tre" khong + cap nhat max. Test truc tiep, khong
// phu thuoc thoi gian thuc.
func TestCapNhatTreNhip(t *testing.T) {
	k := &KhuVuc{}

	// Duoi/bang nguong 200ms: KHONG tinh la tre.
	k.capNhatTreNhip(150 * time.Millisecond)
	k.capNhatTreNhip(200 * time.Millisecond)
	if got := k.soNhipTre.Load(); got != 0 {
		t.Fatalf("150ms va 200ms deu <= nguong, soNhipTre=%d, mong 0", got)
	}
	if got := k.treNhipLauNhatMs.Load(); got != 0 {
		t.Fatalf("chua co nhip tre, treNhipLauNhatMs=%d, mong 0", got)
	}

	// Vuot nguong: dem + ghi max.
	k.capNhatTreNhip(300 * time.Millisecond)
	if got := k.soNhipTre.Load(); got != 1 {
		t.Fatalf("300ms vuot nguong, soNhipTre=%d, mong 1", got)
	}
	if got := k.treNhipLauNhatMs.Load(); got != 300 {
		t.Fatalf("treNhipLauNhatMs=%d, mong 300", got)
	}

	// Tre hon: max len 500.
	k.capNhatTreNhip(500 * time.Millisecond)
	if got := k.soNhipTre.Load(); got != 2 {
		t.Fatalf("soNhipTre=%d, mong 2", got)
	}
	if got := k.treNhipLauNhatMs.Load(); got != 500 {
		t.Fatalf("treNhipLauNhatMs=%d, mong 500", got)
	}

	// Tre nhung NHO hon max cu: dem tang, max GIU 500.
	k.capNhatTreNhip(250 * time.Millisecond)
	if got := k.soNhipTre.Load(); got != 3 {
		t.Fatalf("soNhipTre=%d, mong 3", got)
	}
	if got := k.treNhipLauNhatMs.Load(); got != 500 {
		t.Fatalf("treNhipLauNhatMs=%d, mong giu 500", got)
	}
}

// TestDanhDauNhip khoa DUONG NOI (thu reviewer bat: dot bien wiring van xanh).
// danhDauNhip nhan mot moc thoi gian, BO QUA nhip dau (chua co moc truoc), roi
// tinh khoang cach toi moc truoc va goi capNhatTreNhip, cuoi cung GAN LAI moc.
// Test tat dinh (truyen time gia) -> bat duoc dot bien "quen gan mocNhipTruoc".
func TestDanhDauNhip(t *testing.T) {
	k := &KhuVuc{}
	t0 := time.Unix(1000, 0)

	// Nhip dau: chua co moc -> khong dem gi, chi ghi moc.
	k.danhDauNhip(t0)
	if got := k.soNhipTre.Load(); got != 0 {
		t.Fatalf("nhip dau khong duoc dem tre, soNhipTre=%d, mong 0", got)
	}

	// Nhip 2 cach 300ms (> nguong 200ms) -> dem 1, max 300.
	k.danhDauNhip(t0.Add(300 * time.Millisecond))
	if got := k.soNhipTre.Load(); got != 1 {
		t.Fatalf("cach 300ms phai dem tre, soNhipTre=%d, mong 1", got)
	}
	if got := k.treNhipLauNhatMs.Load(); got != 300 {
		t.Fatalf("treNhipLauNhatMs=%d, mong 300", got)
	}

	// Nhip 3 cach 100ms (< nguong) -> KHONG dem. Neu dot bien QUEN gan lai moc o
	// nhip 2, moc van la t0 -> khoang cach thanh 400ms -> bi dem sai (2). Buoc nay
	// chinh la cai khoa duong noi.
	k.danhDauNhip(t0.Add(400 * time.Millisecond))
	if got := k.soNhipTre.Load(); got != 1 {
		t.Fatalf("cach 100ms khong duoc dem (moc phai da gan lai), soNhipTre=%d, mong giu 1", got)
	}
}
