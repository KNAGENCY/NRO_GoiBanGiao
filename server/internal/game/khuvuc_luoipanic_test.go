package game

import (
	"testing"
	"time"
)

// TestKhuVuc_PanicTrongSuKien_KhuVucVanChayTiep (manh A Task 2): mot su kien
// lam xuLy panic KHONG duoc giet actor — khu vuc phai xu duoc su kien KE TIEP
// (bang chung vong lap con song). Dung suKienChay (diem chen test, khuvuc.go)
// de bom panic co kiem soat.
func TestKhuVuc_PanicTrongSuKien_KhuVucVanChayTiep(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.BatDau()
	t.Cleanup(func() { _ = k.Dung() })

	daGoiKeTiep := make(chan struct{}, 1)
	// Su kien 1: panic. Su kien 2: bao hieu con song.
	k.Gui(suKienChay{f: func() { panic("bug danh nhau gia") }})
	k.Gui(suKienChay{f: func() { daGoiKeTiep <- struct{}{} }})

	select {
	case <-daGoiKeTiep:
		// PASS: vong lap song sot qua panic, xu duoc su kien ke tiep.
	case <-time.After(2 * time.Second):
		t.Fatal("khu vuc chet sau panic — vong lap khong chay tiep")
	}
}
