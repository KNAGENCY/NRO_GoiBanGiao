package main

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

// TestTatMayChuChungKhongPanicKhiControlLaConTroNil tai hien loi 05/09/2026:
// E2E control TAT (mac dinh) -> moE2EControl tra (*http.Server)(nil) -> gan
// vao interface mayChuCoTheTat -> `!= nil` van true -> Shutdown panic -> ca
// trinh tu tat (LUU NGUOI CHOI) bi bo qua. Dot bien: doi coControlThat ve
// `control != nil` -> test nay panic.
func TestTatMayChuChungKhongPanicKhiControlLaConTroNil(t *testing.T) {
	daDonDep := false
	donDep := func(ctx context.Context) error {
		daDonDep = true
		return nil
	}
	var controlNil *http.Server // dung kieu ma moE2EControl tra ve khi TAT
	if err := tatMayChuChung(controlNil, donDep); err != nil {
		t.Fatalf("tatMayChuChung tra loi voi control nil co kieu: %v", err)
	}
	if !daDonDep {
		t.Fatal("donDep KHONG duoc goi — buoc luu nguoi choi bi bo qua")
	}
}

// TestCoControlThat: ba truong hop — nil that, nil co kieu, may chu that.
func TestCoControlThat(t *testing.T) {
	if coControlThat(nil) {
		t.Fatal("nil that phai la false")
	}
	var nilCoKieu *http.Server
	if coControlThat(nilCoKieu) {
		t.Fatal("(*http.Server)(nil) boc trong interface phai la false")
	}
	if !coControlThat(&http.Server{}) {
		t.Fatal("may chu that phai la true")
	}
}

// TestTatMayChuChungVanGoiShutdownKhiCoControlThat: dam bao chan typed-nil
// khong lam mat duong tat control THAT — control gia ghi lai Shutdown.
type controlGia struct{ daTat bool }

func (c *controlGia) Shutdown(context.Context) error { c.daTat = true; return errors.New("loi gia") }

func TestTatMayChuChungVanGoiShutdownKhiCoControlThat(t *testing.T) {
	c := &controlGia{}
	daDonDep := false
	err := tatMayChuChung(c, func(context.Context) error { daDonDep = true; return nil })
	if !c.daTat {
		t.Fatal("control that khong duoc Shutdown")
	}
	if !daDonDep {
		t.Fatal("donDep khong chay sau khi control loi — phai LUON don dep")
	}
	if err == nil {
		t.Fatal("loi Shutdown phai duoc tra ve (errors.Join)")
	}
}
