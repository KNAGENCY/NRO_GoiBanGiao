package game

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"nroserver/internal/store"
)

// ─────────────────────────────────────────────────────────────────────────────
// SO NO — DUONG HONG.
//
// File sonothuong_test.go khoa cac duong THANH CONG. File nay khoa duong HONG,
// va no ra doi vi mot vong soi doi khang chi ra rang toan bo nhom nay dang bo
// trong — trong do co MOT LOI THAT: mot cu panic cua pgx lam mat trang ca lo no,
// dung nguoc voi dieu chu thich cua ham do cam ket.
//
// Bai hoc chung: hang doi nay giu TIEN cua nguoi choi. Duong thanh cong chay
// dung la muc toi thieu; thu quyet dinh no co dang tin hay khong la nhung gi
// xay ra khi CSDL do chung.
// ─────────────────────────────────────────────────────────────────────────────

// khoNoHong la kho gia dieu khien duoc tung lan ghi: panic, tra loi, hoac
// thanh cong — de dung tung kich ban hong mot cach tat dinh.
type khoNoHong struct {
	mu    sync.Mutex
	daGhi []store.NoThuong

	// hanhVi la danh sach viec se lam cho tung lan ghi, theo thu tu. Het danh
	// sach thi ghi binh thuong.
	hanhVi []string // "panic" | "loi" | "ok"
	soLan  int

	// idHongMai: lo nao CO CHUA nhan vat nay thi LUON bi tu choi, bao nhieu lan
	// cung vay. Khac hanhVi (dem theo LAN GOI) — de dung mot lo hong VINH VIEN,
	// thu duy nhat phan biet duoc "co van xa" voi "khong co van xa".
	idHongMai int64
	// soLanGapLoHong dem so lan kho duoc goi voi lo chua idHongMai.
	soLanGapLoHong int
}

func (k *khoNoHong) GhiNoThuong(ctx context.Context, ds []store.NoThuong) error {
	k.mu.Lock()
	n := k.soLan
	k.soLan++
	var lam string
	if n < len(k.hanhVi) {
		lam = k.hanhVi[n]
	}
	coLoHong := false
	if k.idHongMai != 0 {
		for _, no := range ds {
			if no.NhanVatID == k.idHongMai {
				coLoHong = true
				k.soLanGapLoHong++
				break
			}
		}
	}
	k.mu.Unlock()

	if coLoHong {
		return errors.New("lo nay hong vinh vien (vd nhan vat da bi xoa)")
	}

	switch lam {
	case "panic":
		panic("pgx no giua chung")
	case "loi":
		return errors.New("CSDL tu choi")
	}

	k.mu.Lock()
	k.daGhi = append(k.daGhi, ds...)
	k.mu.Unlock()
	return nil
}

func (k *khoNoHong) tongDaGhi(idNhanVat int64) int64 {
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

func (k *khoNoHong) soLanDaGoi() int {
	k.mu.Lock()
	defer k.mu.Unlock()
	return k.soLan
}

// khoaLoDaGhi tra danh sach khoa lo cua cac dong da xuong kho, theo thu tu.
func (k *khoNoHong) lanGapLoHong() int {
	k.mu.Lock()
	defer k.mu.Unlock()
	return k.soLanGapLoHong
}

func (k *khoNoHong) khoaLoDaGhi() []string {
	k.mu.Lock()
	defer k.mu.Unlock()
	var ra []string
	for _, no := range k.daGhi {
		ra = append(ra, no.KhoaLo)
	}
	return ra
}

// TestPanicKhiGhiThiLoNoVanDuocThuLai — LOI THAT da tim ra bang kiem thu dot
// bien. ghiHetHangDoi LAY HAN lo ra khoi hang doi roi moi goi kho; xepLai chi
// chay tren duong TRA LOI. Neu kho PANIC (pgx gap pool da dong, con tro nil...)
// thi luoi o vongLap bat duoc panic — nhung luc do lo da roi hang doi va KHONG
// AI xep lai. Ca lo no bien mat vinh vien.
//
// ThoGhi (thogi.go) khong dinh loi nay vi no co luoi HAI TANG: "thogi.xa" boc
// vong lap VA "thogi.ghi" boc TUNG ban ghi ben trong. So no chi chep tang ngoai.
func TestPanicKhiGhiThiLoNoVanDuocThuLai(t *testing.T) {
	kho := &khoNoHong{hanhVi: []string{"panic"}}
	so := MoSoNoThuong(kho, 10*time.Millisecond)
	defer so.DongVaChoXong(context.Background())

	so.XepNo(store.NoThuong{NhanVatID: 41, Loai: store.LoaiThuongKinhNghiem, SoDiem: 5})

	choToi(t, 3*time.Second, "khoan no bi panic phai duoc thu lai va ghi duoc", func() bool {
		return kho.tongDaGhi(41) == 5
	})
}

// TestPanicKhiGhiKhongGietVongLap: sau cu panic, so no phai con song de ghi
// nhung khoan XEP SAU do. Panic giet vong lap = mat toan bo duong ghi no cho
// toi luc khoi dong lai may chu.
func TestPanicKhiGhiKhongGietVongLap(t *testing.T) {
	kho := &khoNoHong{hanhVi: []string{"panic"}}
	so := MoSoNoThuong(kho, 10*time.Millisecond)
	defer so.DongVaChoXong(context.Background())

	so.XepNo(store.NoThuong{NhanVatID: 41, Loai: store.LoaiThuongKinhNghiem, SoDiem: 5})
	choToi(t, 3*time.Second, "kho bi goi lan dau (va panic)", func() bool {
		return kho.soLanDaGoi() >= 1
	})

	so.XepNo(store.NoThuong{NhanVatID: 42, Loai: store.LoaiThuongKinhNghiem, SoDiem: 9})
	choToi(t, 3*time.Second, "khoan xep SAU cu panic van xuong duoc kho", func() bool {
		return kho.tongDaGhi(42) == 9
	})
}

// TestLoLoiThuLaiNguyenVenKhongTronKhoanMoi khoa cai gia phai tra de chong
// NHAN DOI: mot lo ghi hong phai duoc thu lai NGUYEN VEN — cung so diem, cung
// khoa lo — chu khong duoc tron voi nhung khoan xep them trong luc do.
//
// VI SAO: PostgreSQL co the da commit that roi moi tra loi (het han gio, dut
// ket noi). Lan thu lai se chen lai. Chan nhan doi bang khoa lo + ON CONFLICT
// chi chay duoc neu lan thu lai gui DUNG BANG lo cu. Tron them khoan moi vao
// la khoa lo cu mang so diem moi -> ON CONFLICT bo qua -> mat phan moi.
func TestLoLoiThuLaiNguyenVenKhongTronKhoanMoi(t *testing.T) {
	kho := &khoNoHong{hanhVi: []string{"loi"}}
	so := MoSoNoThuong(kho, 10*time.Millisecond)
	defer so.DongVaChoXong(context.Background())

	so.XepNo(store.NoThuong{NhanVatID: 50, Loai: store.LoaiThuongKinhNghiem, SoDiem: 100})
	choToi(t, 3*time.Second, "lan ghi dau (that bai) da xay ra", func() bool {
		return kho.soLanDaGoi() >= 1
	})

	// Khoan MOI xep sau khi lo cu da hong.
	so.XepNo(store.NoThuong{NhanVatID: 50, Loai: store.LoaiThuongKinhNghiem, SoDiem: 30})

	choToi(t, 3*time.Second, "ca hai khoan deu xuong kho", func() bool {
		return kho.tongDaGhi(50) == 130
	})

	// Hai lo RIENG BIET, moi lo mot khoa — khong duoc gop lam mot dong 130.
	khoa := kho.khoaLoDaGhi()
	if len(khoa) != 2 {
		t.Fatalf("kho nhan %d dong, mong doi 2 (lo cu thu lai nguyen ven + lo moi rieng): %v",
			len(khoa), khoa)
	}
	if khoa[0] == khoa[1] {
		t.Errorf("hai lo mang CUNG khoa %q — khoa lo phai duy nhat moi lo", khoa[0])
	}
	for i, k := range khoa {
		if k == "" {
			t.Errorf("dong %d khong co khoa lo — khong co gi chan nhan doi khi thu lai", i)
		}
	}
}

// TestLoHongMaiThiBiBoSauSoLanThuToiDa — VAN XA.
//
// Mot dong rac (vd nhan vat da bi xoa nen vi pham khoa ngoai) lam ca lo bi tu
// choi MAI MAI. Khong co van xa thi lo do quay vong vo tan, choThuLai khong bao
// gio voi di, va bo nho phinh dan suot doi may chu.
//
// Cach do: dem so lan kho bi goi VOI lo hong. Con van xa thi so do dung han o
// soLanThuLoToiDa; khong con thi no tang mai.
//
// BAI HOC KHI VIET TEST NAY: ban dau toi dung "kho gia tra loi N lan dau roi
// thoi", va dot bien "bo van xa" SONG SOT — vi het N lan thi lan thu N+1 thanh
// cong o CA HAI ban. Phai co mot lo hong VINH VIEN moi phan biet duoc.
func TestLoHongMaiThiBiBoSauSoLanThuToiDa(t *testing.T) {
	kho := &khoNoHong{idHongMai: 60}
	so := MoSoNoThuong(kho, 5*time.Millisecond)
	defer so.DongVaChoXong(context.Background())

	so.XepNo(store.NoThuong{NhanVatID: 60, Loai: store.LoaiThuongKinhNghiem, SoDiem: 100})

	choToi(t, 3*time.Second, "lo hong bi thu du so lan toi da", func() bool {
		return kho.lanGapLoHong() >= soLanThuLoToiDa
	})

	// Sau khi bi bo, con so KHONG duoc tang them nua.
	dung := kho.lanGapLoHong()
	choDen := time.Now().Add(200 * time.Millisecond) // ~40 nhip nua
	for time.Now().Before(choDen) {
		time.Sleep(10 * time.Millisecond)
	}
	if sau := kho.lanGapLoHong(); sau > dung {
		t.Fatalf("lo hong van bi thu tiep (%d -> %d sau ~40 nhip) — van xa khong hoat dong, "+
			"choThuLai se phinh mai suot doi may chu", dung, sau)
	}
	if dung > soLanThuLoToiDa {
		t.Errorf("lo hong bi thu %d lan, mong doi dung %d", dung, soLanThuLoToiDa)
	}
}

// TestLoHongKhongChanDuongCuaLoKhac: mot lo hong vinh vien KHONG duoc lam nhung
// khoan cua nguoi khac ket theo. Day chinh la thu ma viec tach hang doi hai tang
// mua duoc — truoc do lo hong tron chung voi khoan moi nen keo ca hai xuong ho.
func TestLoHongKhongChanDuongCuaLoKhac(t *testing.T) {
	kho := &khoNoHong{idHongMai: 60}
	so := MoSoNoThuong(kho, 5*time.Millisecond)
	defer so.DongVaChoXong(context.Background())

	so.XepNo(store.NoThuong{NhanVatID: 60, Loai: store.LoaiThuongKinhNghiem, SoDiem: 100})
	choToi(t, 3*time.Second, "lo hong da bi thu it nhat mot lan", func() bool {
		return kho.lanGapLoHong() >= 1
	})

	so.XepNo(store.NoThuong{NhanVatID: 61, Loai: store.LoaiThuongKinhNghiem, SoDiem: 7})
	choToi(t, 3*time.Second, "khoan cua nguoi khac van xuong duoc", func() bool {
		return kho.tongDaGhi(61) == 7
	})
}

// TestDongVaChoXongBaoLoiKhiXaCuoiThatBai: luc tat may chu, neu lan xa cuoi
// cung khong ghi duoc thi nhung khoan no do MAT (goroutine chet ngay sau do).
// TatMayChu phai BIET dieu do de bao ra ngoai, thay vi bao "tat sach" trong khi
// vua danh roi tien cua nguoi choi.
func TestDongVaChoXongBaoLoiKhiXaCuoiThatBai(t *testing.T) {
	kho := &khoNoHong{hanhVi: []string{"loi", "loi", "loi", "loi", "loi", "loi", "loi", "loi"}}
	so := MoSoNoThuong(kho, time.Hour) // chi xa khi dong

	so.XepNo(store.NoThuong{NhanVatID: 70, Loai: store.LoaiThuongKinhNghiem, SoDiem: 1})

	if err := so.DongVaChoXong(context.Background()); err == nil {
		t.Fatal("DongVaChoXong tra nil du lan xa cuoi that bai — nguoi goi khong co cach nao " +
			"biet la vua mat no cua nguoi choi")
	}
}

// TestDongVaChoXongGoiHaiLanKhongPanic: tat may chu hai lan (duong tat co the
// chay tu nhieu nguon — tin hieu, E2E control, loi listener) khong duoc lam
// "close of closed channel". ThoGhi da co dungMotLan vi dung ly do nay.
func TestDongVaChoXongGoiHaiLanKhongPanic(t *testing.T) {
	so := MoSoNoThuong(&khoNoGia{}, time.Hour)
	if err := so.DongVaChoXong(context.Background()); err != nil {
		t.Fatalf("lan dong thu nhat: %v", err)
	}
	if err := so.DongVaChoXong(context.Background()); err != nil {
		t.Fatalf("lan dong thu hai: %v", err)
	}
}

// TestMoSoNoThuongNhipKhongHopLeKhongPanic: nhip 0 lam time.NewTicker PANIC —
// va no panic trong goroutine nen, tuc la SAP CA TIEN TRINH may chu ngay luc
// khoi dong. Chot ep ve mac dinh phai giu duoc dieu do.
func TestMoSoNoThuongNhipKhongHopLeKhongPanic(t *testing.T) {
	for _, nhip := range []time.Duration{0, -time.Second} {
		so := MoSoNoThuong(&khoNoGia{}, nhip)
		if err := so.DongVaChoXong(context.Background()); err != nil {
			t.Fatalf("nhip %v: %v", nhip, err)
		}
	}
}

// TestXepNoBoQuaLoaiRong: nua guard con lai cua XepNo. Nua "SoDiem <= 0" da co
// test; nua nay thi chua — mot bat doi xung ma kiem thu dot bien phoi ra. Mot
// dong loai rong lam CA LO bi CSDL tu choi (rang buoc loai_thuong), keo theo
// thuong that cua nguoi khac trong cung lo.
func TestXepNoBoQuaLoaiRong(t *testing.T) {
	kho := &khoNoGia{}
	so := MoSoNoThuong(kho, time.Hour)

	so.XepNo(store.NoThuong{NhanVatID: 21, Loai: "", SoDiem: 100})
	so.XepNo(store.NoThuong{NhanVatID: 22, Loai: store.LoaiThuongKinhNghiem, SoDiem: 7})
	if err := so.DongVaChoXong(context.Background()); err != nil {
		t.Fatalf("dong so no: %v", err)
	}

	if got := kho.soDong(); got != 1 {
		t.Fatalf("kho nhan %d dong, mong doi 1 (dong loai rong phai bi bo)", got)
	}
	if got := kho.tongDaGhi(21); got != 0 {
		t.Errorf("dong loai rong duoc ghi %d diem, mong doi 0", got)
	}
}

// TestKhoaLoKhacNhauGiuaHaiLanMoSoNo — sinh tu dot bien "bo ma phien ngau
// nhien, chi con bo dem", da SONG SOT.
//
// Hai *SoNoThuong* khac nhau la mo phong HAI LAN KHOI DONG may chu. Chi dung bo
// dem thi ca hai deu bat dau tu 1, sinh ra khoa TRUNG voi lo cu con nam trong
// bang — va ON CONFLICT DO NOTHING se BO QUA nham mot khoan no THAT cua nguoi
// choi. Mot lo hong im lang tuyet doi: khong loi, khong log, chi la nguoi choi
// khong bao gio nhan duoc.
func TestKhoaLoKhacNhauGiuaHaiLanMoSoNo(t *testing.T) {
	lay := func() string {
		kho := &khoNoGia{}
		so := MoSoNoThuong(kho, time.Hour)
		so.XepNo(store.NoThuong{NhanVatID: 1, Loai: store.LoaiThuongKinhNghiem, SoDiem: 1})
		if err := so.DongVaChoXong(context.Background()); err != nil {
			t.Fatalf("dong so no: %v", err)
		}
		kho.mu.Lock()
		defer kho.mu.Unlock()
		if len(kho.daGhi) != 1 {
			t.Fatalf("ghi %d dong, mong doi 1", len(kho.daGhi))
		}
		return kho.daGhi[0].KhoaLo
	}

	if a, b := lay(), lay(); a == b {
		t.Fatalf("hai lan mo so no sinh CUNG khoa lo %q — may chu khoi dong lai se "+
			"trung khoa voi lo cu trong bang, va ON CONFLICT se nuot mat mot khoan no that", a)
	}
}

// TestTatMayChuVanXaSoNoDuNganSachChungDaCanKiet — sinh tu dot bien "TatMayChu
// dung chung ctx", da SONG SOT.
//
// ThoGhi chay TRUOC va co the tieu het ngan sach thoi gian (mot CSDL cham la
// du). Dung chung ctx thi DongSoNo nhan mot ctx da het han va tra ngay
// ctx.Err() ma khong xa lan nao — tien cua nguoi choi bi vut di vi mot buoc
// TRUOC no cham. Test truyen thang mot ctx DA HUY de dung dung tinh huong do.
func TestTatMayChuVanXaSoNoDuNganSachChungDaCanKiet(t *testing.T) {
	// Kho phai thoa CA KhoLuu (cho ThoGhi) va KhoNo (cho so no) — dung khuon
	// khoDayDuBietGhiNo o moluutru_nothuong_test.go.
	kho := &khoDayDuBietGhiNo{khoDayDuGia: &khoDayDuGia{}, no: &khoNoGia{}}
	soKhuVuc := MoSoKhuVuc()
	soKhuVuc.GanSoNo(MoSoNoThuong(kho, time.Hour)) // chi xa khi dong

	var k *KhuVuc
	for _, kv := range soKhuVuc.theoMap {
		k = kv
		break
	}
	k.xepNoThuong(-1, 777, 33)

	// ctx DA HET HAN — dung nhu khi thoGhi vua tieu het ngan sach chung.
	ctx, huy := context.WithCancel(context.Background())
	huy()
	_ = TatMayChu(ctx, soKhuVuc, MoThoGhi(kho, 1)) // loi la duong nhien, khong phai thu dang kiem

	if got := kho.no.tongDaGhi(777); got != 33 {
		t.Fatalf("so no ghi duoc %d diem, mong doi 33 — lan xa cuoi phai co han gio RIENG, "+
			"khong duoc chet theo ngan sach ma buoc truoc da tieu het", got)
	}
}
