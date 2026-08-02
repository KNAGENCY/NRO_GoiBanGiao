package main

import (
	"log/slog"
	"testing"
)

// mucNhatKyTuChuoi phai nhan dang dung bon gia tri hop le va TU DONG VE INFO
// (khong phai debug) cho moi gia tri khac, ke ca chuoi rong.
//
// Vi sao mac dinh phai la info: do bang thuc te tren may chu that - mot ket
// noi don gui goi 3 byte lenh la lien tuc dat 32.452 goi/giay, va truoc day
// moi goi lam login.go ghi MOT dong Debug (67 byte). Voi muc Debug bat mac
// dinh, dieu do lam o dia day 9 GB/GIO tu DUY NHAT mot ket noi. Neu quen dat
// NRO_MUC_NHAT_KY (vd trien khai voi cau hinh mac dinh) ma mac dinh la
// LevelDebug nhu truoc day, loi nay tai dien ngay lap tuc.
func TestMucNhatKyTuChuoiMacDinhLaInfoKhongPhaiDebug(t *testing.T) {
	truongHop := []struct {
		vao          string
		muc          slog.Level
		hopLeMongDoi bool
	}{
		{"", slog.LevelInfo, true},
		{"info", slog.LevelInfo, true},
		{"INFO", slog.LevelInfo, true},
		{"debug", slog.LevelDebug, true},
		{"Debug", slog.LevelDebug, true},
		{"warn", slog.LevelWarn, true},
		{"error", slog.LevelError, true},
		{"gia-tri-la", slog.LevelInfo, false},
	}

	for _, tc := range truongHop {
		muc, hopLe := mucNhatKyTuChuoi(tc.vao)
		if muc != tc.muc {
			t.Errorf("mucNhatKyTuChuoi(%q) = %v, mong doi %v", tc.vao, muc, tc.muc)
		}
		if hopLe != tc.hopLeMongDoi {
			t.Errorf("mucNhatKyTuChuoi(%q) hopLe = %v, mong doi %v", tc.vao, hopLe, tc.hopLeMongDoi)
		}
	}
}

// Main production khong dat NRO_E2E phai giu nguyen mat tan cong: khong tao
// HTTP server, khong bind listener va khong co goroutine control ngam.
func TestMainKhongMoE2EControlKhiCauHinhTat(t *testing.T) {
	mayChu, langNghe, err := moE2EControl(
		cauHinhE2EControl{},
		make(chan struct{}, 1),
	)
	if err != nil {
		t.Fatalf("control tat tra loi: %v", err)
	}
	if mayChu != nil || langNghe != nil {
		t.Fatalf(
			"E2E-CONTROL-001: production tao control ngoai y muon: server=%v listener=%v",
			mayChu,
			langNghe,
		)
	}
}
