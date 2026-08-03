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
