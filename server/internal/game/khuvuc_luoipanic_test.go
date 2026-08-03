package game

import (
	"testing"
	"time"
)

// TestKhuVuc_PanicTrongSuKien_KhuVucVanChayTiep (manh A Task 2): mot su kien
// lam xuLy panic KHONG duoc giet actor — khu vuc phai xu duoc su kien KE TIEP
// (bang chung vong lap con song). Dung suKienChay (diem chen test, khuvuc.go)
// de bom panic co kiem soat.
func TestKhuVuc_PanicTrongSuKien_KhuVucVanChayTiep(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.BatDau()
	t.Cleanup(func() { _ = k.Dung() })

	daGoiKeTiep := make(chan struct{}, 1)
	// Su kien 1: panic. Su kien 2: bao hieu con song.
	k.Gui(suKienChay{f: func() { panic("bug danh nhau gia") }})
	k.Gui(suKienChay{f: func() { daGoiKeTiep <- struct{}{} }})

	select {
	case <-daGoiKeTiep:
		// PASS: vong lap song sot qua panic, xu duoc su kien ke tiep.
	case <-time.After(2 * time.Second):
		t.Fatal("khu vuc chet sau panic — vong lap khong chay tiep")
	}
}

// TestKhuVuc_PanicDaiDangQuaNguong_TuDungKhuVuc (manh A Task 3): mot bug lam
// actor panic DAI DANG (moi nhip khong dut) phai lam VAN AN TOAN dung RIENG khu
// vuc do (ghi log muc Error) de no khong quay nong CPU + ngap log mai mai. Bom
// LIEN TUC su kien panic: giua hai nhip tim (100ms) vong lap xu duoc hang tram
// panic LIEN TIEP nen chac chan cham nguong truoc khi mot nhipTim trot kip
// reset -> xac dinh, khong phu thuoc timing.
func TestKhuVuc_PanicDaiDangQuaNguong_TuDungKhuVuc(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.BatDau()
	t.Cleanup(func() { _ = k.Dung() })

	bomXong := make(chan struct{})
	go func() {
		defer close(bomXong)
		for {
			select {
			case <-k.xong:
				return
			default:
				k.Gui(suKienChay{f: func() { panic("hong dai dang") }})
			}
		}
	}()

	select {
	case <-k.xong:
		<-bomXong // don goroutine bom truoc khi ket thuc
	case <-time.After(10 * time.Second):
		t.Fatal("khu vuc panic dai dang nhung khong tu dung — van an toan khong chay")
	}
}

// TestKhuVuc_VanAnToan_PanicLienTiepDatNguongThiBaoDung kiem logic THUAN cua
// van (ghiNhanNhip) khong qua ticker: panic lien tiep dung nguong -> nhip cuoi
// bao dung. Bang chung logic ngưỡng chinh xac, xac dinh, nhanh.
func TestKhuVuc_VanAnToan_PanicLienTiepDatNguongThiBaoDung(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	// nguong-1 nhip panic dau tien KHONG duoc bao dung.
	for i := 1; i < nguongPanicDungKhuVuc; i++ {
		if k.ghiNhanNhip(true) {
			t.Fatalf("bao dung som o nhip panic thu %d, mong chi dung tai nguong %d", i, nguongPanicDungKhuVuc)
		}
	}
	// Nhip panic thu nguong -> PHAI bao dung.
	if !k.ghiNhanNhip(true) {
		t.Fatalf("panic lien tiep dat nguong %d nhung van khong bao dung", nguongPanicDungKhuVuc)
	}
}

// TestKhuVuc_VanAnToan_NhipTrotResetKhongBaoGioDung la TEST BAT BUOC cho kiem
// thu dot bien Task 3: loi THOANG QUA (moi panic deu co nhip trot xen ngay
// sau) DU tong so panic vuot XA nguong van KHONG BAO GIO duoc bao dung — vi bo
// dem reset moi lan trot. Neu bo dong reset (k.soPanicLienTiep = 0), tich luy
// se vuot nguong va test nay DO -> chung minh dong reset co tac dung that.
func TestKhuVuc_VanAnToan_NhipTrotResetKhongBaoGioDung(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	for i := 0; i < nguongPanicDungKhuVuc*3; i++ {
		if k.ghiNhanNhip(true) { // panic thoang qua
			t.Fatalf("bao dung nham o vong %d: loi thoang qua (co nhip trot xen giua) khong duoc cham nguong", i)
		}
		if k.ghiNhanNhip(false) { // nhip trot ngay sau -> reset
			t.Fatal("nhip trot khong duoc bao dung")
		}
	}
}
