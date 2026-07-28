package game

import (
	"context"
	"testing"
	"time"

	"nroserver/internal/store"
)

// khoDayDuGia gop khoGia (thoa man Kho, login_test.go) + khoLuuGia (thoa
// man KhoLuu, thogi_test.go) de thoa man KhoDayDu ma MoLuuTru doi hoi — chi
// NHUNG (embed) lai hai fake da co san, khong can viet lai method nao.
type khoDayDuGia struct {
	khoGia
	khoLuuGia
}

// TestMoLuuTruNoiDuBaDiem la test WIRING BAT BUOC cua Task 6 diem B "NOI DAY
// THAT": xac nhan MoLuuTru noi DU CA BA DIEM ma brief liet ke. TRUOC Task 6,
// khong ai trong toan bo du an goi MoThoGhi / MoBoLuuDinhKy(...).BatDau() /
// DatThoGhi — nghia la may chu THAT SU KHONG LUU GI CA du Task 1-5 da xong
// va qua review. Neu mot trong ba diem duoi day bi thieu (vd sua nham xoa
// boLuu.BatDau(), hoac quen DatThoGhi trong taoXuLy), test nay PHAI do.
//
//  1. thoGhi := MoThoGhi(kho, N)      -> thoGhi != nil VA THAT SU ghi duoc
//     xuong kho (kiem bang mot lan Xep + DongVaChoXong, dung khuon
//     thogi_test.go — khong chi kiem "khac nil" cho co).
//  2. BoLuuDinhKy duoc TAO + BatDau() -> boLuu.dung/boLuu.xong (kenh dieu
//     khien cua goroutine nen, CHI duoc cap phat BEN TRONG BatDau(), xem
//     luudinhky.go) phai KHAC nil — day la dau hieu KHANG DINH ro rang hon
//     la chi dua vao "goi Dung() khong panic".
//  3. taoXuLy() tra ve MOT *BoXuLy da DatThoGhi(thoGhi) -> truong noi bo
//     b.thoGhi phai CHINH LA con tro thoGhi vua tao o diem 1 (khong phai
//     nil, cung khong phai mot ThoGhi KHAC).
func TestMoLuuTruNoiDuBaDiem(t *testing.T) {
	kho := &khoDayDuGia{}
	soKhuVuc := MoSoKhuVuc()
	t.Cleanup(func() { _ = soKhuVuc.DungTatCa() })

	thoGhi, boLuu, taoXuLy := MoLuuTru(kho, soKhuVuc, 4)
	if boLuu != nil {
		t.Cleanup(boLuu.Dung)
	}

	// --- Diem 1: MoThoGhi ---
	if thoGhi == nil {
		t.Fatal("MoLuuTru tra thoGhi = nil — diem noi day 1 (MoThoGhi) chua duoc goi")
	}
	thoGhi.Xep(store.NhanVat{ID: 321, ViTriX: 1})
	ctx, huy := context.WithTimeout(context.Background(), 2*time.Second)
	defer huy()
	if err := thoGhi.DongVaChoXong(ctx); err != nil {
		t.Fatalf("DongVaChoXong loi: %v", err)
	}
	if soLan := kho.khoLuuGia.soLanGoi(); soLan != 1 {
		t.Fatalf("kho gia nhan %d lan ghi sau Xep+DongVaChoXong, mong doi 1 — thoGhi tra ve boi MoLuuTru "+
			"khong THAT SU noi toi kho da truyen vao (diem noi day 1 chua du)", soLan)
	}

	// --- Diem 2: MoBoLuuDinhKy(...).BatDau() ---
	if boLuu == nil {
		t.Fatal("MoLuuTru tra boLuu = nil — diem noi day 2 (MoBoLuuDinhKy) chua duoc goi")
	}
	if boLuu.dung == nil || boLuu.xong == nil {
		t.Fatal("BoLuuDinhKy tra ve CHUA duoc BatDau() (kenh dung/xong van nil) — diem noi day 2 thieu goi BatDau()")
	}

	// --- Diem 3: DatThoGhi trong taoXuLy ---
	if taoXuLy == nil {
		t.Fatal("MoLuuTru tra taoXuLy = nil — diem noi day 3 khong co gi de kiem")
	}
	b := taoXuLy()
	if b == nil {
		t.Fatal("taoXuLy() tra ve nil")
	}
	got, ok := b.thoGhi.(*ThoGhi)
	if !ok || got != thoGhi {
		t.Fatalf("BoXuLy tao tu taoXuLy() co thoGhi = %#v, mong doi CHINH LA con tro thoGhi (%p) tao o diem 1 — "+
			"diem noi day 3 (DatThoGhi trong taoXuLy) chua duoc goi dung", b.thoGhi, thoGhi)
	}
}
