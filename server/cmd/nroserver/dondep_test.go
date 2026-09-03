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

// boGiaDemDung la mot boDoCoTheDung gia (KH7 Task 12), ghi lai no bi Dung()
// bao nhieu lan va — quan trong hon — ACTOR DA DUNG CHUA vao dung luc do.
//
// `so` co the nil khi test khong quan tam thu tu.
type boGiaDemDung struct {
	so *game.SoKhuVuc

	soLan int
	// actorDaDung: tai thoi diem Dung() duoc goi, cac actor khu vuc da dung
	// chua. Day la thu duoc chot o TestBoVatPhamDungSauTatMayChu.
	actorDaDung bool
}

func (b *boGiaDemDung) Dung() {
	b.soLan++
	if b.so != nil {
		k := b.so.Tra(1)
		if k != nil {
			_, err := k.AnhChupNguoiChoiCoLoi(1)
			b.actorDaDung = err != nil
		}
	}
}

// ⚠️ CHOT THU TU TAT (KH7 Task 12): hai bo ghi CSDL cua he vat pham phai dung
// SAU game.TatMayChu, khong truoc.
//
// Ly do day du o chu thich trong donDepKhiTat: chung la NOI NHAN, con NGUON la
// cac actor khu vuc. Chung nao con mot actor song thi con co the co nguoi bam
// nhat hay mot con quai chet va xep them yeu cau; dong hang doi truoc la vut di
// dung nhung yeu cau cuoi cung.
//
// Doi hai dong boRoi.Dung()/boNhat.Dung() len TRUOC game.TatMayChu -> test do.
func TestBoVatPhamDungSauTatMayChu(t *testing.T) {
	soKhuVuc := soKhuVucGiaMain(t)
	thoGhi := game.MoThoGhi(&khoLuuGiaMain{}, 2)
	boLuu := game.MoBoLuuDinhKy(soKhuVuc, thoGhi)
	boLuu.BatDau()
	kd := &khoDongGia{}

	boNhat := &boGiaDemDung{so: soKhuVuc}
	boRoi := &boGiaDemDung{so: soKhuVuc}

	if err := donDepKhiTat(context.Background(), kd, soKhuVuc, thoGhi, boLuu, nil,
		game.MoBoLichBoss(soKhuVuc, nil), game.MoBoPhatThanh(), boNhat, boRoi); err != nil {
		t.Fatalf("donDepKhiTat loi: %v", err)
	}

	if boNhat.soLan != 1 || boRoi.soLan != 1 {
		t.Fatalf("moi bo phai duoc Dung() dung mot lan: boNhat=%d boRoi=%d",
			boNhat.soLan, boRoi.soLan)
	}
	if !boNhat.actorDaDung {
		t.Fatal("boNhat.Dung() chay TRUOC khi actor dung - sai thu tu, xem chu thich trong donDepKhiTat")
	}
	if !boRoi.actorDaDung {
		t.Fatal("boRoi.Dung() chay TRUOC khi actor dung - sai thu tu, xem chu thich trong donDepKhiTat")
	}
}

// kenhGiaDemDung la mot kenhCoTheDung gia (KH8 task 10): ghi lai no bi Dung()
// bao nhieu lan va — quan trong hon — ACTOR DA DUNG CHUA vao dung luc do.
// Doi xung voi boGiaDemDung o tren, nhung khang dinh dieu NGUOC LAI: kenh nap
// nong phai dung TRUOC khi actor dung.
type kenhGiaDemDung struct {
	so *game.SoKhuVuc

	soLan       int
	actorDaDung bool
}

func (k *kenhGiaDemDung) Dung(context.Context) error {
	k.soLan++
	if k.so != nil {
		if kv := k.so.Tra(1); kv != nil {
			_, err := kv.AnhChupNguoiChoiCoLoi(1)
			k.actorDaDung = err != nil
		}
	}
	return nil
}

// ⚠️ CHOT THU TU TAT (KH8 task 10): kenh nap nong phai dung TRUOC
// game.TatMayChu, khong sau.
//
// Ly do day du o chu thich trong donDepKhiTat. Tom lai: de sau thi mot lan nap
// nong den dung luc do se tra ve "200, the he moi" cho mot may chu khong con
// actor nao doc bo so lieu do — mot cau tra loi DOI TRA, dung loai hong ma ca
// cai kenh HTTP nay sinh ra de tranh.
//
// Doi dong kenhNapLai.Dung() xuong duoi game.TatMayChu -> test nay do.
func TestKenhNapLaiDungTruocTatMayChu(t *testing.T) {
	soKhuVuc := soKhuVucGiaMain(t)
	thoGhi := game.MoThoGhi(&khoLuuGiaMain{}, 2)
	boLuu := game.MoBoLuuDinhKy(soKhuVuc, thoGhi)
	boLuu.BatDau()

	kenh := &kenhGiaDemDung{so: soKhuVuc}
	if err := donDepKhiTat(context.Background(), &khoDongGia{}, soKhuVuc, thoGhi, boLuu,
		kenh, game.MoBoLichBoss(soKhuVuc, nil), game.MoBoPhatThanh(),
		&boGiaDemDung{}, &boGiaDemDung{}); err != nil {
		t.Fatalf("donDepKhiTat loi: %v", err)
	}
	if kenh.soLan != 1 {
		t.Fatalf("kenh nap nong phai duoc Dung() dung mot lan, nhan %d", kenh.soLan)
	}
	if kenh.actorDaDung {
		t.Fatal("kenh nap nong bi dung SAU khi actor da dung — sai thu tu, xem chu thich trong donDepKhiTat")
	}
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

	if err := donDepKhiTat(ctx, kd, soKhuVuc, thoGhi, boLuu, nil, game.MoBoLichBoss(soKhuVuc, nil), game.MoBoPhatThanh(), &boGiaDemDung{}, &boGiaDemDung{}); err != nil {
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

	err := donDepKhiTat(ctx, kd, soKhuVuc, thoGhi, boLuu, nil, game.MoBoLichBoss(soKhuVuc, nil), game.MoBoPhatThanh(), &boGiaDemDung{}, &boGiaDemDung{})
	if err == nil {
		t.Fatal("donDepKhiTat khong tra loi du ctx da het gio tu truoc khi goi — mong doi ctx.Err()")
	}
	if !kd.laDaDong() {
		t.Fatal("donDepKhiTat KHONG goi kho.Close() khi TatMayChu loi — phai LUON dong kho du buoc luu co loi hay khong")
	}
}
