package game

import (
	"testing"

	"nroserver/internal/protocol"
)

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

// TestGoiDiChuyen kiểm tra gói -7 phát di chuyển: writeInt(id), writeShort(x), writeShort(y)
// Đọc lại bằng Reader, khẳng định hết buffer.
func TestGoiDiChuyen(t *testing.T) {
	m := goiDiChuyen(12345, 100, 384)

	// Kiểm Command = -7
	if m.Command != -7 {
		t.Errorf("expect command -7, got %d", m.Command)
	}

	// Đọc lại bằng Reader
	r := protocol.NewReader(m.Data)

	id, err := r.Int()
	if err != nil {
		t.Fatalf("read id: %v", err)
	}
	if id != 12345 {
		t.Errorf("expect id 12345, got %d", id)
	}

	x, err := r.Short()
	if err != nil {
		t.Fatalf("read x: %v", err)
	}
	if x != 100 {
		t.Errorf("expect x 100, got %d", x)
	}

	y, err := r.Short()
	if err != nil {
		t.Fatalf("read y: %v", err)
	}
	if y != 384 {
		t.Errorf("expect y 384, got %d", y)
	}

	// Kiểm hết buffer: r.Byte() phải lỗi
	_, err = r.Byte()
	if err == nil {
		t.Error("expect error when reading past end of buffer (hết data)")
	}
}

// TestGoiRoiKhuVuc kiểm tra gói -6 rời khu vực: writeInt(id)
// Đọc lại bằng Reader, khẳng định hết buffer.
func TestGoiRoiKhuVuc(t *testing.T) {
	m := goiRoiKhuVuc(777)

	// Kiểm Command = -6
	if m.Command != -6 {
		t.Errorf("expect command -6, got %d", m.Command)
	}

	// Đọc lại bằng Reader
	r := protocol.NewReader(m.Data)

	id, err := r.Int()
	if err != nil {
		t.Fatalf("read id: %v", err)
	}
	if id != 777 {
		t.Errorf("expect id 777, got %d", id)
	}

	// Kiểm hết buffer: r.Byte() phải lỗi
	_, err = r.Byte()
	if err == nil {
		t.Error("expect error when reading past end of buffer (hết data)")
	}
}
