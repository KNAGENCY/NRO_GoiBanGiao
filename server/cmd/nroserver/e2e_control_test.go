package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"nroserver/internal/game"
	"nroserver/internal/store"
)

func envE2E(giaTri map[string]string) func(string) string {
	return func(ten string) string { return giaTri[ten] }
}

func TestE2EControlChiBatKhiCoCoMoiTruongRoRang(t *testing.T) {
	t.Run("thieu co E2E thi tat", func(t *testing.T) {
		cfg, err := docCauHinhE2EControl(envE2E(map[string]string{
			"NRO_E2E_CONTROL_ADDR":  "127.0.0.1:19091",
			"NRO_E2E_CONTROL_TOKEN": "token-test-khong-phai-secret-0001",
		}))
		if err != nil {
			t.Fatalf("control dang tat khong duoc tra loi: %v", err)
		}
		if cfg.Bat {
			t.Fatal("E2E-CONTROL-001: control mo khi thieu NRO_E2E=1")
		}
	})

	t.Run("bat E2E ma thieu token thi fail closed", func(t *testing.T) {
		_, err := docCauHinhE2EControl(envE2E(map[string]string{
			"NRO_E2E":              "1",
			"NRO_E2E_CONTROL_ADDR": "127.0.0.1:19091",
		}))
		if err == nil || !strings.Contains(err.Error(), "E2E-CONTROL-002") {
			t.Fatalf("muon E2E-CONTROL-002 khi thieu token, nhan %v", err)
		}
	})

	t.Run("dia chi ngoai loopback bi chan", func(t *testing.T) {
		_, err := docCauHinhE2EControl(envE2E(map[string]string{
			"NRO_E2E":               "1",
			"NRO_E2E_CONTROL_ADDR":  "0.0.0.0:19091",
			"NRO_E2E_CONTROL_TOKEN": "token-test-khong-phai-secret-0001",
		}))
		if err == nil || !strings.Contains(err.Error(), "E2E-CONTROL-003") {
			t.Fatalf("muon E2E-CONTROL-003 cho non-loopback, nhan %v", err)
		}
	})
}

func TestE2EControlAuthMethodVaEndpointFailClosed(t *testing.T) {
	const token = "token-test-khong-phai-secret"
	yeuCauTat := make(chan struct{}, 1)
	handler := taoE2EControlHandler(token, yeuCauTat)

	thu := func(method, path, bearer string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(method, path, nil)
		if bearer != "" {
			req.Header.Set("Authorization", "Bearer "+bearer)
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		return rec
	}
	khongCoTinHieu := func() {
		t.Helper()
		select {
		case <-yeuCauTat:
			t.Fatal("E2E-CONTROL-004: auth/method sai van phat tin hieu shutdown")
		default:
		}
	}

	if got := thu(http.MethodPost, "/shutdown", "").Code; got != http.StatusUnauthorized {
		t.Fatalf("thieu token: status=%d, muon 401", got)
	}
	khongCoTinHieu()
	if got := thu(http.MethodPost, "/shutdown", "token-sai").Code; got != http.StatusUnauthorized {
		t.Fatalf("token sai: status=%d, muon 401", got)
	}
	khongCoTinHieu()
	if got := thu(http.MethodGet, "/shutdown", token).Code; got != http.StatusMethodNotAllowed {
		t.Fatalf("GET /shutdown: status=%d, muon 405", got)
	}
	khongCoTinHieu()
	if got := thu(http.MethodPost, "/exit", token).Code; got != http.StatusNotFound {
		t.Fatalf("endpoint ngoai /shutdown: status=%d, muon 404", got)
	}
	khongCoTinHieu()
}

func TestE2EControlShutdownPhatMotTinHieuVaTraNgay(t *testing.T) {
	yeuCauTat := make(chan struct{}, 1)
	handler := taoE2EControlHandler("token-test-khong-phai-secret", yeuCauTat)
	for lan := 0; lan < 2; lan++ {
		req := httptest.NewRequest(http.MethodPost, "/shutdown", nil)
		req.Header.Set("Authorization", "Bearer token-test-khong-phai-secret")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusAccepted {
			t.Fatalf("lan %d status=%d, muon 202", lan+1, rec.Code)
		}
	}
	select {
	case <-yeuCauTat:
	case <-time.After(time.Second):
		t.Fatal("E2E-CONTROL-005: handler khong phat tin hieu shutdown")
	}
	select {
	case <-yeuCauTat:
		t.Fatal("E2E-CONTROL-005: hai request lam day hai tin hieu")
	default:
	}
}

func TestE2EControlListenerThucTeChiBindLoopback(t *testing.T) {
	cfg := cauHinhE2EControl{
		Bat:    true,
		DiaChi: "127.0.0.1:0",
		Token:  "token-test-khong-phai-secret-0001",
	}
	mayChu, langNghe, err := moE2EControl(cfg, make(chan struct{}, 1))
	if err != nil {
		t.Fatalf("mo listener loopback: %v", err)
	}
	if mayChu == nil || langNghe == nil {
		t.Fatal("E2E-CONTROL-003: cau hinh bat nhung khong co listener")
	}
	defer mayChu.Close()
	defer langNghe.Close()
	host, _, err := net.SplitHostPort(langNghe.Addr().String())
	if err != nil {
		t.Fatalf("listener addr sai: %v", err)
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		t.Fatalf("E2E-CONTROL-003: listener khong loopback: %s", langNghe.Addr())
	}
}

func TestShutdownQuaControlDiQuaDonDepVaFlushThoGhiThat(t *testing.T) {
	yeuCauTat := make(chan struct{}, 1)
	handler := taoE2EControlHandler("token-test-khong-phai-secret", yeuCauTat)
	req := httptest.NewRequest(http.MethodPost, "/shutdown", nil)
	req.Header.Set("Authorization", "Bearer token-test-khong-phai-secret")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("E2E-PERSIST-004: control status=%d, muon 202", rec.Code)
	}

	khoLuu := &khoLuuGiaMain{}
	thoGhi := game.MoThoGhi(khoLuu, 1)
	muonLuu := store.NhanVat{ID: 991, HpHienTai: 137, HpToiDa: 150, ViTriX: 121, ViTriY: 384}
	thoGhi.Xep(muonLuu)
	soKhuVuc := game.MoSoKhuVuc()
	boLuu := game.MoBoLuuDinhKy(soKhuVuc, thoGhi)
	boLuu.BatDau()
	khoDong := &khoDongGia{}
	maThoat, loiNguon, loiDonDep := choVaTatMayChu(
		nil,
		yeuCauTat,
		nil,
		nil,
		nil,
		func(ctx context.Context) error {
			return donDepKhiTat(ctx, khoDong, soKhuVuc, thoGhi, boLuu, game.MoBoLichBoss(soKhuVuc, nil), game.MoBoPhatThanh(), &boGiaDemDung{}, &boGiaDemDung{})
		},
	)
	if maThoat != 0 || loiNguon != nil || loiDonDep != nil {
		t.Fatalf(
			"E2E-PERSIST-004: lifecycle control tra ma=%d nguon=%v don_dep=%v",
			maThoat,
			loiNguon,
			loiDonDep,
		)
	}

	khoLuu.mu.Lock()
	defer khoLuu.mu.Unlock()
	if len(khoLuu.goiGhi) != 1 || khoLuu.goiGhi[0].ID != muonLuu.ID ||
		khoLuu.goiGhi[0].HpHienTai != 137 {
		t.Fatalf("E2E-PERSIST-004: persistence chua flush ban 137/150: %#v", khoLuu.goiGhi)
	}
	if !khoDong.laDaDong() {
		t.Fatal("E2E-PERSIST-004: graceful path khong dong kho sau flush")
	}
}
