package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"nroserver/internal/game"
	"nroserver/internal/store"
)

// khoLuuGiaMain gia lap game.KhoLuu (interface CHI co CapNhatNhanVat, xem
// thogi.go) tu package main — khong the tai dung khoLuuGia cua package game
// (unexported, khac package), nen viet mot fake NHO rieng o day.
type khoLuuGiaMain struct {
	mu     sync.Mutex
	doTre  time.Duration
	goiGhi []store.NhanVat
}

func (k *khoLuuGiaMain) CapNhatNhanVat(ctx context.Context, nv store.NhanVat) error {
	if k.doTre > 0 {
		time.Sleep(k.doTre)
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	k.goiGhi = append(k.goiGhi, nv)
	return nil
}

// khoDongGia gia lap khoDong (interface CHI co Close(), dondep.go/main.go)
// de test khong can PostgreSQL that — chi dem so lan Close() duoc goi.
type khoDongGia struct {
	mu     sync.Mutex
	soLan  int
	daDong bool
}

func (k *khoDongGia) Close() {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.soLan++
	k.daDong = true
}

func (k *khoDongGia) laDaDong() bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	return k.daDong
}

// soKhuVucGiaMain tao mot SoKhuVuc voi 1 khu vuc THAT (actor that, BatDau
// that) — dung game.MoSoKhuVuc() se tao 5 khu vuc that, hoi nang cho mot
// test wiring don gian, nhung chap nhan duoc vi khong co nguoi choi nao va
// DungTatCa()/Dung() cua actor rong deu nhanh (khong I/O).
func soKhuVucGiaMain(t *testing.T) *game.SoKhuVuc {
	t.Helper()
	return game.MoSoKhuVuc()
}

// TestDonDepKhiTatGoiDuBonBuoc la test WIRING cua Task 6 diem A tu phia
// main(): donDepKhiTat phai (1) dung boLuu, (2) chay TatMayChu (da tu no
// dam bao thu tu LuuTatCa->DongVaChoXong->DungTatCa — kiem RIENG, sau, o
// package game qua TestTatMayChuLuuTruocKhiDungKhuVuc), va (3) goi
// kho.Close() — TAT CA khong loi khi dieu kien binh thuong.
//
// Test nay KHONG lap lai kich ban "co nguoi trong khu vuc" (da co du va ky
// o package game) — muc dich o day la kiem PLUMBING rieng cua package main:
// donDepKhiTat co THAT SU goi ca bon manh (boLuu, TatMayChu qua thoGhi/
// soKhuVuc, kho) hay khong.
func TestDonDepKhiTatGoiDuBonBuoc(t *testing.T) {
	soKhuVuc := soKhuVucGiaMain(t)
	kl := &khoLuuGiaMain{}
	thoGhi := game.MoThoGhi(kl, 2)
	boLuu := game.MoBoLuuDinhKy(soKhuVuc, thoGhi)
	boLuu.BatDau()
	kd := &khoDongGia{}

	ctx, huy := context.WithTimeout(context.Background(), 5*time.Second)
	defer huy()

	if err := donDepKhiTat(ctx, kd, soKhuVuc, thoGhi, boLuu, game.MoBoPhatThanh()); err != nil {
		t.Fatalf("donDepKhiTat loi trong dieu kien binh thuong: %v", err)
	}

	if !kd.laDaDong() {
		t.Fatal("donDepKhiTat KHONG goi kho.Close() — diem 4 (kho.Close) chua duoc goi")
	}

	// boLuu phai da dung: goi Dung() lan nua PHAI tra ve ngay (sync.Once +
	// kenh xong da dong), khong duoc treo — day la bang chung boLuu.Dung()
	// da THAT SU chay ben trong donDepKhiTat (buoc 0).
	xong := make(chan struct{})
	go func() {
		boLuu.Dung()
		close(xong)
	}()
	select {
	case <-xong:
	case <-time.After(2 * time.Second):
		t.Fatal("goi lai boLuu.Dung() sau donDepKhiTat bi treo — boLuu chua thuc su duoc dung trong donDepKhiTat")
	}
}

// TestDonDepKhiTatVanDongKhoKhiTatMayChuLoi xac nhan donDepKhiTat KHONG bo
// qua buoc kho.Close() du TatMayChu (buoc 1-3, boc trong donDepKhiTat) gap
// loi (vd ctx da het gio) — dung tinh than "mot lan tat 'danh hong' con hon
// treo tien trinh vo han", ap dung o CA hai tang (game.TatMayChu VA
// donDepKhiTat deu phai giu bat bien nay doc lap nhau).
func TestDonDepKhiTatVanDongKhoKhiTatMayChuLoi(t *testing.T) {
	soKhuVuc := soKhuVucGiaMain(t)
	kl := &khoLuuGiaMain{doTre: 2 * time.Second} // kho gia CHAM hon nhieu han gio ctx duoi day
	thoGhi := game.MoThoGhi(kl, 1)
	boLuu := game.MoBoLuuDinhKy(soKhuVuc, thoGhi)
	boLuu.BatDau()
	kd := &khoDongGia{}

	ctx, huy := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer huy()
	time.Sleep(5 * time.Millisecond) // dam bao ctx CHAC CHAN da het gio truoc khi goi

	err := donDepKhiTat(ctx, kd, soKhuVuc, thoGhi, boLuu, game.MoBoPhatThanh())
	if err == nil {
		t.Fatal("donDepKhiTat khong tra loi du ctx da het gio tu truoc khi goi — mong doi ctx.Err()")
	}
	if !kd.laDaDong() {
		t.Fatal("donDepKhiTat KHONG goi kho.Close() khi TatMayChu loi — phai LUON dong kho du buoc luu co loi hay khong")
	}
}
