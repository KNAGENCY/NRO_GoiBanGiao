package game

import "testing"

// TestKiemDiChuyen kiểm tra hàm kiemDiChuyenHopLe: di chuyển hợp lệ trong biên
// và khoảng cách mỗi bước ≤ NguongBuoc (500). Khoảng cách = Euclid (so bình phương).
func TestKiemDiChuyen(t *testing.T) {
	bd := TimBanDo(1)

	// Bước ngắn, hợp lệ
	if !kiemDiChuyenHopLe(bd, 100, 100, 300, 100) {
		t.Error("buoc ngan hop le phai true")
	}

	// toX âm, ngoài biên
	if kiemDiChuyenHopLe(bd, 100, 100, -5, 100) {
		t.Error("toX am phai false")
	}

	// toX > BienRongToiDa, ngoài biên
	if kiemDiChuyenHopLe(bd, 100, 100, 5000, 100) {
		t.Error("ngoai bien phai false")
	}

	// Khoảng cách = 500 (= NguongBuoc), hợp lệ. 400²+300²=250000=500²
	if !kiemDiChuyenHopLe(bd, 0, 0, 400, 300) {
		t.Error("khoang cach 500 (=nguong) phai true")
	}

	// Khoảng cách 600 > NguongBuoc, không hợp lệ
	if kiemDiChuyenHopLe(bd, 0, 0, 600, 0) {
		t.Error("khoang cach 600 > 500 phai false")
	}

	// Đúng mép biên toX=4096, nhưng khoảng cách 4096 > 500 nên false
	if kiemDiChuyenHopLe(bd, 0, 0, 4096, 0) {
		t.Error("edge case: toX=4096 trong bien nhung khoang cach 4096 > 500 phai false")
	}

	// Đứng yên (cuX==toX && cuY==toY), khoảng cách 0 ≤ 500, hợp lệ
	if !kiemDiChuyenHopLe(bd, 100, 100, 100, 100) {
		t.Error("dung yen phai true")
	}
}
