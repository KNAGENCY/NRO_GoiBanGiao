package game

import (
	"context"
	"testing"
	"time"

	"nroserver/internal/store"
)

// ─────────────────────────────────────────────────────────────────────────────
// NOI DAY so no — bai hoc Task 6 diem B, nguyen van: "ha tang khong duoc noi
// day o MoLuuTru la ha tang nam do khong ai goi toi, test tu che van xanh ma
// may chu that van tho".
//
// Chinh commit Task 12 vua roi cung dam vao viec do: CongSatThuong va
// ChiaThuongTheoDongGop nam trong boss.go tu Task 3, khong noi nao goi toi
// ngoai test, nen "chia thuong theo dong gop" tren giay chua bao gio chay.
// So no se lap lai y het neu khong co hai bai kiem duoi day.
// ─────────────────────────────────────────────────────────────────────────────

// khoDayDuBietGhiNo la khoDayDuGia CO THEM GhiNoThuong — tuc la mot kho THAT
// SU day du nhu *store.Store. MoLuuTru nhan ra bang type-assert.
type khoDayDuBietGhiNo struct {
	*khoDayDuGia
	no *khoNoGia
}

func (k *khoDayDuBietGhiNo) GhiNoThuong(ctx context.Context, ds []store.NoThuong) error {
	return k.no.GhiNoThuong(ctx, ds)
}

// TestMoLuuTruNoiSoNoVaoMoiKhuVuc: kho biet ghi no -> MoLuuTru phai mo mot so
// no VA cam no vao CA 5 khu vuc. Thieu buoc cam thi actor van log "BO phan
// thuong" y nhu truoc, khong ai biet.
func TestMoLuuTruNoiSoNoVaoMoiKhuVuc(t *testing.T) {
	kho := &khoDayDuBietGhiNo{khoDayDuGia: &khoDayDuGia{}, no: &khoNoGia{}}
	soKhuVuc := MoSoKhuVuc()
	t.Cleanup(func() { _ = soKhuVuc.DungTatCa() })

	_, boLuu, _ := MoLuuTru(kho, soKhuVuc, 1, 0, 0)
	if boLuu != nil {
		t.Cleanup(boLuu.Dung)
	}
	t.Cleanup(func() { _ = soKhuVuc.DongSoNo(context.Background()) })

	for mapId, k := range soKhuVuc.theoMap {
		if k.soNo.Load() == nil {
			t.Fatalf("khu vuc map %d chua duoc cam so no — hang doi ghi no nam do khong ai goi toi", mapId)
		}
	}
}

// TestMoLuuTruSoNoChayThatXuongKho: khong chi "khac nil" cho co — khoan no xep
// vao phai THAT SU xuong toi kho da truyen vao.
func TestMoLuuTruSoNoChayThatXuongKho(t *testing.T) {
	kho := &khoDayDuBietGhiNo{khoDayDuGia: &khoDayDuGia{}, no: &khoNoGia{}}
	soKhuVuc := MoSoKhuVuc()
	_, boLuu, _ := MoLuuTru(kho, soKhuVuc, 1, 0, 0)
	if boLuu != nil {
		t.Cleanup(boLuu.Dung)
	}

	var k *KhuVuc
	for _, kv := range soKhuVuc.theoMap {
		k = kv
		break
	}
	k.xepNoThuong(-1, 555, 42)

	// Tat may chu la duong DUY NHAT xa het so no khi dong — kiem luon o day.
	ctx, huy := context.WithTimeout(context.Background(), 2*time.Second)
	defer huy()
	if err := TatMayChu(ctx, soKhuVuc, MoThoGhi(kho, 1)); err != nil {
		t.Fatalf("TatMayChu: %v", err)
	}

	if got := kho.no.tongDaGhi(555); got != 42 {
		t.Fatalf("kho nhan %d diem no, mong doi 42 — so no khong noi that toi kho, "+
			"hoac TatMayChu khong dong so no", got)
	}
}

// TestMoLuuTruKhoKhongBietGhiNoThiVanChay: mot kho khong cai GhiNoThuong (moi
// kho gia trong test cu) phai chay y nhu truoc — khong panic, chi la khong co
// so no. Day la ly do MoLuuTru dung type-assert thay vi nhet KhoNo vao
// KhoDayDu: nhet vao thi hang chuc kho gia khong lien quan phai sua theo.
func TestMoLuuTruKhoKhongBietGhiNoThiVanChay(t *testing.T) {
	kho := &khoDayDuGia{}
	soKhuVuc := MoSoKhuVuc()
	t.Cleanup(func() { _ = soKhuVuc.DungTatCa() })

	_, boLuu, _ := MoLuuTru(kho, soKhuVuc, 1, 0, 0)
	if boLuu != nil {
		t.Cleanup(boLuu.Dung)
	}

	t.Cleanup(func() { _ = soKhuVuc.DongSoNo(context.Background()) })

	for mapId, k := range soKhuVuc.theoMap {
		if k.soNo.Load() != nil {
			t.Fatalf("khu vuc map %d duoc cam so no du kho khong biet ghi no", mapId)
		}
		// Va duong ghi no phai di qua duoc ma khong panic.
		k.xepNoThuong(-1, 1, 10)
	}
}
