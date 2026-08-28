package game

import (
	"context"
	"sync"
	"testing"
	"time"

	"nroserver/internal/store"
)

// khoNoGia ghi nhan cac lo no ma SoNoThuong day xuong, va cho phep ep loi /
// ep cham de kiem hai tinh chat quan trong nhat: KHONG BAO GIO CHAN nguoi goi,
// va KHONG MAT no khi CSDL hong.
type khoNoGia struct {
	mu       sync.Mutex
	daGhi    []store.NoThuong
	soLanGoi int

	// loiConLai > 0 thi lan ghi ke tiep tra loi va tru dan — dung de dung
	// kich ban "CSDL hong vai nhip roi song lai".
	loiConLai int

	// chan != nil thi moi lan ghi se cho tren kenh nay — dung de gia lap mot
	// CSDL treo, kiem XepNo van tra ve ngay.
	chan_ chan struct{}
}

func (k *khoNoGia) GhiNoThuong(ctx context.Context, ds []store.NoThuong) error {
	k.mu.Lock()
	k.soLanGoi++
	choi := k.chan_
	if k.loiConLai > 0 {
		k.loiConLai--
		k.mu.Unlock()
		return context.DeadlineExceeded // mot loi bat ky, khong quan trong loai gi
	}
	k.mu.Unlock()

	if choi != nil {
		<-choi
	}

	k.mu.Lock()
	k.daGhi = append(k.daGhi, ds...)
	k.mu.Unlock()
	return nil
}

// tongDaGhi cong het so diem da xuong "CSDL" cho mot nhan vat.
func (k *khoNoGia) tongDaGhi(idNhanVat int64) int64 {
	k.mu.Lock()
	defer k.mu.Unlock()
	var tong int64
	for _, no := range k.daGhi {
		if no.NhanVatID == idNhanVat {
			tong += no.SoDiem
		}
	}
	return tong
}

func (k *khoNoGia) soDong() int {
	k.mu.Lock()
	defer k.mu.Unlock()
	return len(k.daGhi)
}

// choToi cho toi khi dieuKien() dung, hoac het han — thay cho sleep co dinh
// (test khong duoc phu thuoc vao viec doan dung nhip cua goroutine nen).
func choToi(t *testing.T, han time.Duration, mo string, dieuKien func() bool) {
	t.Helper()
	hetHan := time.Now().Add(han)
	for time.Now().Before(hetHan) {
		if dieuKien() {
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatalf("qua %v van chua: %s", han, mo)
}

// TestXepNoKhongChanDuKhoDangTreo la BAI KIEM QUAN TRONG NHAT cua file nay:
// XepNo duoc goi TU TRONG actor khu vuc, nen no phai tra ve ngay lap tuc du
// CSDL dang treo. Neu no cho, mot lan CSDL cham se lam dung ca ban do — dung
// cai loi §6.10b da phai va mot lan roi.
func TestXepNoKhongChanDuKhoDangTreo(t *testing.T) {
	kho := &khoNoGia{chan_: make(chan struct{})}
	so := MoSoNoThuong(kho, 10*time.Millisecond)
	t.Cleanup(func() {
		close(kho.chan_)
		so.DongVaChoXong(context.Background())
	})

	// Xep lan dau de goroutine nen bat dau mot lan ghi va MAC KET o do.
	so.XepNo(store.NoThuong{NhanVatID: 1, Loai: store.LoaiThuongKinhNghiem, SoDiem: 10})
	choToi(t, time.Second, "kho bat dau lan ghi dau", func() bool {
		kho.mu.Lock()
		defer kho.mu.Unlock()
		return kho.soLanGoi > 0
	})

	// Trong luc kho dang treo, 100 lan xep nua phai xong gan nhu tuc thi.
	batDau := time.Now()
	for i := 0; i < 100; i++ {
		so.XepNo(store.NoThuong{NhanVatID: 2, Loai: store.LoaiThuongKinhNghiem, SoDiem: 1})
	}
	if mat := time.Since(batDau); mat > 100*time.Millisecond {
		t.Fatalf("100 lan XepNo mat %v khi kho dang treo — phai gan nhu tuc thi", mat)
	}
}

// TestXepNoCongDonCungNguoiCungLoai: hang doi phai CONG DON chu khong ghi de.
// Day la khac biet co y so voi ThoGhi (ghi de ban moi nhat) — ghi de o day la
// lam bien mat phan thuong cua nguoi choi.
func TestXepNoCongDonCungNguoiCungLoai(t *testing.T) {
	kho := &khoNoGia{}
	so := MoSoNoThuong(kho, 10*time.Millisecond)
	defer so.DongVaChoXong(context.Background())

	so.XepNo(store.NoThuong{NhanVatID: 7, Loai: store.LoaiThuongKinhNghiem, SoDiem: 100})
	so.XepNo(store.NoThuong{NhanVatID: 7, Loai: store.LoaiThuongKinhNghiem, SoDiem: 200})
	so.XepNo(store.NoThuong{NhanVatID: 7, Loai: store.LoaiThuongKinhNghiem, SoDiem: 50})

	choToi(t, 2*time.Second, "kho nhan du 350 diem cho nguoi 7", func() bool {
		return kho.tongDaGhi(7) == 350
	})
}

// TestXepNoBoQuaDiemKhongDuong: mot khoan no <= 0 khong co y nghia gi va se bi
// CSDL tu choi ca lo (rang buoc so_diem > 0) — chan ngay tu day de mot dong rac
// khong keo theo ca lo that bai roi thu lai vo han.
func TestXepNoBoQuaDiemKhongDuong(t *testing.T) {
	kho := &khoNoGia{}
	so := MoSoNoThuong(kho, 10*time.Millisecond)
	defer so.DongVaChoXong(context.Background())

	so.XepNo(store.NoThuong{NhanVatID: 3, Loai: store.LoaiThuongKinhNghiem, SoDiem: 0})
	so.XepNo(store.NoThuong{NhanVatID: 3, Loai: store.LoaiThuongKinhNghiem, SoDiem: -5})
	so.XepNo(store.NoThuong{NhanVatID: 4, Loai: store.LoaiThuongKinhNghiem, SoDiem: 9})

	choToi(t, 2*time.Second, "kho nhan khoan hop le cua nguoi 4", func() bool {
		return kho.tongDaGhi(4) == 9
	})
	if got := kho.tongDaGhi(3); got != 0 {
		t.Fatalf("nguoi 3 duoc ghi %d diem, mong 0 (diem khong duong phai bi bo)", got)
	}
}

// TestSoNoThuongThuLaiKhiGhiLoi: CSDL hong vai nhip thi no phai NAM LAI trong
// hang doi va duoc thu lai — khong duoc bien mat im lang (bai hoc
// PlayerDAO.updatePlayer ban Java, dac ta §6.3).
func TestSoNoThuongThuLaiKhiGhiLoi(t *testing.T) {
	kho := &khoNoGia{loiConLai: 3}
	so := MoSoNoThuong(kho, 10*time.Millisecond)
	defer so.DongVaChoXong(context.Background())

	so.XepNo(store.NoThuong{NhanVatID: 5, Loai: store.LoaiThuongKinhNghiem, SoDiem: 42})

	choToi(t, 3*time.Second, "no duoc thu lai va ghi thanh cong sau 3 lan loi", func() bool {
		return kho.tongDaGhi(5) == 42
	})
}

// TestDongVaChoXongXaHetHangDoi: luc tat may chu, nhung khoan no vua xep phai
// duoc ghi not truoc khi tien trinh thoat.
func TestDongVaChoXongXaHetHangDoi(t *testing.T) {
	kho := &khoNoGia{}
	so := MoSoNoThuong(kho, time.Hour)

	so.XepNo(store.NoThuong{NhanVatID: 8, Loai: store.LoaiThuongKinhNghiem, SoDiem: 11})
	so.XepNo(store.NoThuong{NhanVatID: 9, Loai: store.LoaiThuongKinhNghiem, SoDiem: 22})

	if err := so.DongVaChoXong(context.Background()); err != nil {
		t.Fatalf("dong so no: %v", err)
	}
	if got := kho.tongDaGhi(8); got != 11 {
		t.Fatalf("nguoi 8 duoc ghi %d, mong 11 sau khi dong", got)
	}
	if got := kho.tongDaGhi(9); got != 22 {
		t.Fatalf("nguoi 9 duoc ghi %d, mong 22 sau khi dong", got)
	}
}

// TestXepNoGopThanhMotDongMoiNguoi: hai khoan cua cung mot nguoi phai xuong
// CSDL thanh MOT dong (da gop trong bo nho), khong phai hai — de bang khong
// phinh theo so lan boss chet.
func TestXepNoGopThanhMotDongMoiNguoi(t *testing.T) {
	kho := &khoNoGia{}
	so := MoSoNoThuong(kho, time.Hour)

	so.XepNo(store.NoThuong{NhanVatID: 6, Loai: store.LoaiThuongKinhNghiem, SoDiem: 1})
	so.XepNo(store.NoThuong{NhanVatID: 6, Loai: store.LoaiThuongKinhNghiem, SoDiem: 2})
	if err := so.DongVaChoXong(context.Background()); err != nil {
		t.Fatalf("dong so no: %v", err)
	}

	if got := kho.soDong(); got != 1 {
		t.Fatalf("kho nhan %d dong, mong 1 (hai khoan cung nguoi phai gop)", got)
	}
	if got := kho.tongDaGhi(6); got != 3 {
		t.Fatalf("tong ghi %d, mong 3", got)
	}
}

// TestXepNoBaoHoaKhongTran: cong don trong bo nho cung phai bao hoa. Cuon vong
// ve so am se sinh mot khoan no AM, thu ma CSDL se tu choi — lam ca lo (gom ca
// no that cua nguoi khac) bi vut di roi thu lai vo han.
func TestXepNoBaoHoaKhongTran(t *testing.T) {
	kho := &khoNoGia{}
	so := MoSoNoThuong(kho, time.Hour)

	const gan = int64(9_223_372_036_854_775_800)
	so.XepNo(store.NoThuong{NhanVatID: 10, Loai: store.LoaiThuongKinhNghiem, SoDiem: gan})
	so.XepNo(store.NoThuong{NhanVatID: 10, Loai: store.LoaiThuongKinhNghiem, SoDiem: 1000})
	if err := so.DongVaChoXong(context.Background()); err != nil {
		t.Fatalf("dong so no: %v", err)
	}

	const maxInt64 = int64(9_223_372_036_854_775_807)
	if got := kho.tongDaGhi(10); got != maxInt64 {
		t.Fatalf("tong ghi %d, mong bao hoa o %d", got, maxInt64)
	}
}
