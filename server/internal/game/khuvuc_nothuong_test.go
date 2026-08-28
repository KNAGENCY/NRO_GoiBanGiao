package game

import (
	"context"
	"testing"
	"time"
)

// ─────────────────────────────────────────────────────────────────────────────
// So no phan thuong — noi duong chia thuong boss (Task 12) vao so no ben.
//
// TRUOC BAN NAY: nguoi dong gop sat thuong cho boss nhung da roi khu vuc /
// mat ket noi dung luc boss chet thi MAT TRANG phan cua minh — traThuongBoss
// chi log Warn kem so diem bi bo. Ly do khi do la dung: actor khong so huu con
// tro nv cua ho (§6.9), va day mot ban chup cu qua ThoGhi se ghi de CA DONG.
//
// BAN NAY: phan cua ho duoc xep vao so no (bang thuong_cho_nhan) qua mot hang
// doi ngoai actor, va tra o lan dang nhap sau. Actor van khong cham CSDL, van
// khong cham con tro cua nguoi khac — chi cham mot map trong bo nho vai chuc
// nano-giay.
// ─────────────────────────────────────────────────────────────────────────────

// khuVucCoSoNo mo mot khu vuc test DA noi voi mot so no that (nhip gom dai:
// chi xa khi dong, nen test tat dinh — khong doan nhip goroutine nen) va tra
// ve ca kho gia de doc lai nhung gi da xuong "CSDL".
//
// Dung SoNoThuong THAT (khong phai mot ban gia cua chinh no) de bai kiem di
// qua CA duong that: actor -> hang doi -> lo ghi.
func khuVucCoSoNo(t *testing.T) (*KhuVuc, *khoNoGia, *SoNoThuong) {
	t.Helper()
	k := khuVucTest(t)
	kho := &khoNoGia{}
	so := MoSoNoThuong(kho, time.Hour)
	k.soNo.Store(so)
	return k, kho, so
}

// xaSoNo dong so no va cho ghi xong — goi sau khi boss chet de doc ket qua.
func xaSoNo(t *testing.T, so *SoNoThuong) {
	t.Helper()
	if err := so.DongVaChoXong(context.Background()); err != nil {
		t.Fatalf("dong so no: %v", err)
	}
}

// TestBossChetGhiNoChoNguoiDaRoiKhuVuc la bai kiem TRUNG TAM: nguoi 99 danh
// 50% mau boss roi roi khu vuc truoc khi boss chet. Phan cua ho (500/1000 diem)
// phai thanh mot khoan NO, khong bi vut di.
func TestBossChetGhiNoChoNguoiDaRoiKhuVuc(t *testing.T) {
	k, kho, so := khuVucCoSoNo(t)
	nvA := datNguoi(t, k, 11, "A")

	b := bossThuong(BossHoatDong, 1000, false)
	// 99 co ten trong so dong gop nhung KHONG co trong k.nguoiTrong.
	b.SatThuongTheoNguoi = map[int64]int64{11: 50, 99: 50}

	k.bossChet(b, 11, mocTest)
	xaSoNo(t, so)

	if got := kho.tongDaGhi(99); got != 500 {
		t.Errorf("no ghi cho nguoi da roi = %d diem, mong doi 500", got)
	}
	// Nguoi con o lai van duoc cong THANG vao nhan vat nhu cu.
	if nvA.KinhNghiem != 500 {
		t.Errorf("KinhNghiem cua nguoi con o lai = %d, mong doi 500", nvA.KinhNghiem)
	}
}

// TestBossChetKhongGhiNoChoNguoiConTrongKhuVuc: nguoi con o day da duoc cong
// THANG vao nv.KinhNghiem — ghi them mot khoan no cho ho nghia la ho se nhan
// HAI LAN o lan dang nhap sau.
func TestBossChetKhongGhiNoChoNguoiConTrongKhuVuc(t *testing.T) {
	k, kho, so := khuVucCoSoNo(t)
	datNguoi(t, k, 11, "A")

	b := bossThuong(BossHoatDong, 1000, false)
	b.SatThuongTheoNguoi = map[int64]int64{11: 100}

	k.bossChet(b, 11, mocTest)
	xaSoNo(t, so)

	if got := kho.soDong(); got != 0 {
		t.Fatalf("ghi %d dong no, mong doi 0 (nguoi con trong khu vuc da duoc cong thang)", got)
	}
}

// TestBossChetGhiNoChoNhieuNguoiVangMat: hai nguoi vang mat phai ra HAI khoan
// no rieng, moi khoan dung phan cua nguoi do — khong gop nham vao mot ai.
func TestBossChetGhiNoChoNhieuNguoiVangMat(t *testing.T) {
	k, kho, so := khuVucCoSoNo(t)
	datNguoi(t, k, 11, "A")

	b := bossThuong(BossHoatDong, 1000, false)
	b.SatThuongTheoNguoi = map[int64]int64{11: 20, 88: 50, 99: 30}

	k.bossChet(b, 11, mocTest)
	xaSoNo(t, so)

	if got := kho.tongDaGhi(88); got != 500 {
		t.Errorf("no cua nguoi 88 = %d, mong doi 500", got)
	}
	if got := kho.tongDaGhi(99); got != 300 {
		t.Errorf("no cua nguoi 99 = %d, mong doi 300", got)
	}
	if got := kho.soDong(); got != 2 {
		t.Errorf("ghi %d dong no, mong doi dung 2", got)
	}
}

// TestBossChetGhiNoKhiBanGhiNhanVatNil: nguoi con ten trong k.nguoiTrong nhung
// ban ghi nhan vat da nil (mot loi o cho khac, hoac phien dang bi thao do). Cu
// cu la BO phan thuong — nhung ta VAN biet id cua ho, nen van ghi no duoc.
func TestBossChetGhiNoKhiBanGhiNhanVatNil(t *testing.T) {
	k, kho, so := khuVucCoSoNo(t)
	datNguoi(t, k, 11, "A")
	// Thao ban ghi nhan vat cua nguoi 11 — gan thang trang thai noi bo, dung
	// tinh than "chi test cung package moi duoc cham" da ap dung o cac task truoc.
	k.nguoiTrong[11].nv = nil

	b := bossThuong(BossHoatDong, 1000, false)
	b.SatThuongTheoNguoi = map[int64]int64{11: 100}

	k.bossChet(b, 11, mocTest) // khong duoc panic
	xaSoNo(t, so)

	if got := kho.tongDaGhi(11); got != 1000 {
		t.Errorf("no ghi cho nguoi co ban ghi nil = %d, mong doi 1000", got)
	}
}

// TestBossChetKhongCoSoNoThiVanChayBinhThuong: khu vuc chua noi day so no
// (moi test cu, va ca giai doan truoc khi main.go goi GanSoNo) phai chay y
// nhu truoc — khong panic, nguoi con o lai van nhan dung phan.
func TestBossChetKhongCoSoNoThiVanChayBinhThuong(t *testing.T) {
	k := khuVucTest(t) // KHONG gan so no
	nvA := datNguoi(t, k, 11, "A")

	b := bossThuong(BossHoatDong, 1000, false)
	b.SatThuongTheoNguoi = map[int64]int64{11: 50, 99: 50}

	k.bossChet(b, 11, mocTest) // khong duoc panic

	if nvA.KinhNghiem != 500 {
		t.Errorf("KinhNghiem = %d, mong doi 500", nvA.KinhNghiem)
	}
}

// TestBossChetKhongGhiNoKhiPhaKhongCoThuong: pha khong co DiemThuong thi khong
// ai duoc gi — ke ca duoi dang no.
func TestBossChetKhongGhiNoKhiPhaKhongCoThuong(t *testing.T) {
	k, kho, so := khuVucCoSoNo(t)
	datNguoi(t, k, 11, "A")

	b := bossThuong(BossHoatDong, 0, false)
	b.SatThuongTheoNguoi = map[int64]int64{11: 50, 99: 50}

	k.bossChet(b, 11, mocTest)
	xaSoNo(t, so)

	if got := kho.soDong(); got != 0 {
		t.Fatalf("ghi %d dong no cho pha khong co thuong, mong doi 0", got)
	}
}

// TestXepNoThuongDungLoaiKinhNghiem: khoan no phai mang dung loai ma duong
// nhan (store.NhanNoThuong) biet tra. Ghi sai loai thi no nam lai vinh vien —
// mot loi im lang tuyet doi, khong ai thay cho toi khi nguoi choi keu.
func TestXepNoThuongDungLoaiKinhNghiem(t *testing.T) {
	k, kho, so := khuVucCoSoNo(t)
	datNguoi(t, k, 11, "A")

	b := bossThuong(BossHoatDong, 1000, false)
	b.SatThuongTheoNguoi = map[int64]int64{11: 50, 99: 50}

	k.bossChet(b, 11, mocTest)
	xaSoNo(t, so)

	kho.mu.Lock()
	defer kho.mu.Unlock()
	if len(kho.daGhi) != 1 {
		t.Fatalf("ghi %d dong no, mong doi 1", len(kho.daGhi))
	}
	// GHIM CHUOI CUNG, khong dung store.LoaiThuongKinhNghiem — day la diem mu
	// oracle ma kiem thu dot bien phoi ra (§6.12): lay ky vong tu chinh hang so
	// ma code dang dung nghia la doi ten hang so thi test VAN XANH, trong khi
	// khoan no ghi ra mang loai la va nam lai vinh vien khong ai tra.
	if kho.daGhi[0].Loai != "kinh_nghiem" {
		t.Errorf("loai no = %q, mong doi %q", kho.daGhi[0].Loai, "kinh_nghiem")
	}
}
