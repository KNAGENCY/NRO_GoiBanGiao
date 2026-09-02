package game

// muaban_ban_test.go — KH7 Task 10: test THUAN cua duong ban (khong can CSDL).
//
// Giai ma giao thuc: docs/superpowers/research/2026-09-03-kh7-giai-ma-ban-7.md

import (
	"encoding/binary"
	"strings"
	"testing"
)

// Goi XUONG `7` — [where byte][oSo short][cau hoi UTF].
//
// Con so 7 va thu tu ba o ghim CUNG o day (khong doc lai hang so tu code): no
// den tu Controller.cs:3669-3674 + ShopService.java:884-887. Doc lai hang so
// thi doi hang so la test van xanh.
func TestGoiBanXacNhanByteDung(t *testing.T) {
	m := goiBanXacNhan(1, 5, "abc")
	if m == nil {
		t.Fatal("goiBanXacNhan tra nil cho dau vao hop le")
	}
	if m.Command != 7 {
		t.Fatalf("ma lenh = %d, mong doi 7", m.Command)
	}
	// 1 byte where + 2 byte oSo + 2 byte do dai UTF + 3 byte "abc" = 8.
	if len(m.Data) != 8 {
		t.Fatalf("than goi %d byte, mong doi 8", len(m.Data))
	}
	if m.Data[0] != 1 {
		t.Errorf("where = %d, phai doi lai nguyen van 1", m.Data[0])
	}
	if v := int16(binary.BigEndian.Uint16(m.Data[1:3])); v != 5 {
		t.Errorf("oSo = %d, phai doi lai nguyen van 5", v)
	}
	if v := binary.BigEndian.Uint16(m.Data[3:5]); v != 3 {
		t.Errorf("do dai UTF = %d, mong doi 3", v)
	}
	if string(m.Data[5:8]) != "abc" {
		t.Errorf("noi dung = %q, mong doi \"abc\"", string(m.Data[5:8]))
	}
}

// Cau hoi xac nhan phai co du ba manh: so luong, ten mon, tong tien + loai tien.
func TestBanCauHoiCoDuBaManh(t *testing.T) {
	cau := banCauHoi(banDuTinh{
		SoLuong: 3, TenMon: "Thuoc hoi phuc", DonGia: 25, Tien: TienVang, DVang: 75,
	})
	for _, manh := range []string{"x3", "Thuoc hoi phuc", "75", "vàng"} {
		if !strings.Contains(cau, manh) {
			t.Errorf("cau hoi %q thieu manh %q", cau, manh)
		}
	}
	cauNgoc := banCauHoi(banDuTinh{
		SoLuong: 1, TenMon: "Ngoc", DonGia: 1, Tien: TienNgoc, DNgoc: 1,
	})
	if !strings.Contains(cauNgoc, "ngọc") {
		t.Errorf("cau hoi ban bang ngoc phai noi 'ngọc', duoc %q", cauNgoc)
	}
}

// Chua mo cua hang -> tu choi ngay, khong cham gi.
func TestBanTinhToanChuaMoCuaHang(t *testing.T) {
	if _, err := banTinhToan(nil, nil, 0, 0); err != ErrKhongMoCuaHang {
		t.Fatalf("mong doi ErrKhongMoCuaHang, duoc %v", err)
	}
}

// O so ngoai bien -> tu choi (bay §7.11: kiem bien oSo o MOI nhanh).
func TestBanTinhToanOSoNgoaiBien(t *testing.T) {
	for _, o := range []int16{-1, tranSoO, 30000} {
		if _, err := banTinhToan(cuaHangThu(), nil, o, 0); err == nil {
			t.Errorf("oSo = %d phai bi tu choi", o)
		}
	}
}
