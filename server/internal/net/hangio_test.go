package net

import (
	"testing"
	"time"
)

// TestHanDocDuChoNguoiChoiDungYen chot rang han doc mot goi tin khong ngan hon
// 5 phut: client Unity khong gui gi khi nguoi choi khong thao tac, han 120s cu
// lam nguoi dung yen 2 phut bi ngat (do that 05/09/2026). Ket noi chet that do
// TCP keep-alive lo, khong phai han nay.
func TestHanDocDuChoNguoiChoiDungYen(t *testing.T) {
	if thoiGianChoGoiTin < 5*time.Minute {
		t.Fatalf("thoiGianChoGoiTin = %s — nguoi choi dung yen se bi ngat; toi thieu 5 phut", thoiGianChoGoiTin)
	}
	if thoiGianChoBatTay > 10*time.Second {
		t.Fatalf("thoiGianChoBatTay = %s — bat tay phai ngan de chan ket noi rac", thoiGianChoBatTay)
	}
}
