package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type cauHinhE2EControl struct {
	Bat    bool
	DiaChi string
	Token  string
}

func docCauHinhE2EControl(getenv func(string) string) (cauHinhE2EControl, error) {
	if strings.TrimSpace(getenv("NRO_E2E")) != "1" {
		return cauHinhE2EControl{}, nil
	}

	diaChi := strings.TrimSpace(getenv("NRO_E2E_CONTROL_ADDR"))
	token := getenv("NRO_E2E_CONTROL_TOKEN")
	if diaChi == "" {
		return cauHinhE2EControl{}, fmt.Errorf(
			"E2E-CONTROL-002: NRO_E2E=1 can NRO_E2E_CONTROL_ADDR")
	}
	// Token E2E la 32 byte ngau nhien. Chan token ngan de mot cau hinh mau/
	// bo trong khong vo tinh thanh cong dieu khien may chu.
	if len(token) < 32 {
		return cauHinhE2EControl{}, fmt.Errorf(
			"E2E-CONTROL-002: NRO_E2E=1 can token toi thieu 32 byte")
	}
	host, _, err := net.SplitHostPort(diaChi)
	if err != nil {
		return cauHinhE2EControl{}, fmt.Errorf(
			"E2E-CONTROL-003: dia chi control khong hop le: %w", err)
	}
	ip := net.ParseIP(strings.Trim(host, "[]"))
	if ip == nil || !ip.IsLoopback() {
		return cauHinhE2EControl{}, fmt.Errorf(
			"E2E-CONTROL-003: control chi duoc bind dia chi IP loopback")
	}
	return cauHinhE2EControl{Bat: true, DiaChi: diaChi, Token: token}, nil
}

func taoE2EControlHandler(token string, yeuCauTat chan<- struct{}) http.Handler {
	var chiPhatMotLan sync.Once
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/shutdown" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		const tienTo = "Bearer "
		authorization := r.Header.Get("Authorization")
		nhan := ""
		if strings.HasPrefix(authorization, tienTo) {
			nhan = strings.TrimPrefix(authorization, tienTo)
		}
		// Bam hai phia ve cung 32 byte truoc khi so sanh: ke ca token thieu,
		// sai do dai hay sai noi dung deu di qua cung mot phep constant-time.
		muonBam := sha256.Sum256([]byte(token))
		nhanBam := sha256.Sum256([]byte(nhan))
		if nhan == "" || subtle.ConstantTimeCompare(muonBam[:], nhanBam[:]) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		chiPhatMotLan.Do(func() {
			// Channel phai buffered; select co default giu handler khong bao
			// gio bi treo neu main da nhan yeu cau dau tien.
			select {
			case yeuCauTat <- struct{}{}:
			default:
			}
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "shutdown_requested",
		})
	})
}

func moE2EControl(cauHinh cauHinhE2EControl, yeuCauTat chan<- struct{}) (*http.Server, net.Listener, error) {
	if !cauHinh.Bat {
		return nil, nil, nil
	}
	langNghe, err := net.Listen("tcp", cauHinh.DiaChi)
	if err != nil {
		return nil, nil, fmt.Errorf("E2E-CONTROL-006: khong bind duoc control: %w", err)
	}
	host, _, err := net.SplitHostPort(langNghe.Addr().String())
	if err != nil {
		_ = langNghe.Close()
		return nil, nil, fmt.Errorf("E2E-CONTROL-003: listener khong co dia chi hop le: %w", err)
	}
	ip := net.ParseIP(strings.Trim(host, "[]"))
	if ip == nil || !ip.IsLoopback() {
		_ = langNghe.Close()
		return nil, nil, fmt.Errorf("E2E-CONTROL-003: listener thuc te khong phai loopback")
	}
	mayChu := &http.Server{
		Handler:           taoE2EControlHandler(cauHinh.Token, yeuCauTat),
		ReadHeaderTimeout: 2 * time.Second,
	}
	return mayChu, langNghe, nil
}
