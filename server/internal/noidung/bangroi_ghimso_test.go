package noidung

import (
	"reflect"
	"testing"
)

// Ghim TỪNG CON SỐ của bang-roi.yaml nhúng theo bảng Go cũ (roido.go trước
// KH8 task 9). Các test khác của bảng rơi đọc YAML rồi so với chính nó, nên
// đổi 15 → 99 trong YAML không test nào đỏ (đã thử tay lúc duyệt task 9).
// Test này là lưới cho lời "chép nguyên văn": ai đổi số ở YAML (đúng việc của
// task 13 — vòng chỉnh số) thì phải đổi cả đây, có chủ đích.
func TestBangRoiNhung_GhimSoKhoiDauTheoBangGoCu(t *testing.T) {
	du, err := NapNhung()
	if err != nil {
		t.Fatal(err)
	}
	muon := BangRoi{
		QuaiMacDinh: []MucRoi{
			{MauId: 0, SoLuong: 1, TiLePhanTram: 15},
			{MauId: 1, SoLuong: 1, TiLePhanTram: 10},
			{MauId: 4, SoLuong: 1, TiLePhanTram: 2},
		},
		Quai: map[string][]MucRoi{
			"soi-rung": {
				{MauId: 0, SoLuong: 1, TiLePhanTram: 10},
				{MauId: 1, SoLuong: 1, TiLePhanTram: 5},
			},
		},
		Boss: map[string][]MucRoi{
			MaBangRoiBossChung: {
				{MauId: 0, SoLuong: 1, ChacChan: true},
				{MauId: 1, SoLuong: 1, ChacChan: true},
				{MauId: 4, SoLuong: 1, TiLePhanTram: 30},
				{MauId: 2, SoLuong: 1, TiLePhanTram: 50, ChiBossLon: true},
				{MauId: 3, SoLuong: 1, TiLePhanTram: 50, ChiBossLon: true},
			},
		},
	}
	if !reflect.DeepEqual(du.BangRoi, muon) {
		t.Fatalf("bang-roi.yaml lech bang Go cu:\n co: %+v\nmuon: %+v", du.BangRoi, muon)
	}
}
