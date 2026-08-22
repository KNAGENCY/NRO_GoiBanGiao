package game

import (
	"sync"
	"testing"
)

// TestCapIdBoss_AmVaDuyNhatQua1000LanGoiDongThoi khoa bay #2 (thiet ke §7):
// ban goc sinh 5 boss cung id=-1, khien findCharInMap(id) phia client tra nham
// con. Yeu cau: id LUON AM va KHONG BAO GIO trung, ke ca khi 1000 luong cung
// xin mot luc.
//
// Vi sao chay dong thoi that (goroutine + WaitGroup) chu khong goi tuan tu
// 1000 lan: mot bo dem viet bang "doc -> +1 -> ghi" thuong (khong nguyen tu)
// van CHAY DUNG khi goi tuan tu — chi khi co tranh chap that no moi cap trung
// so. Goi tuan tu se KHONG bat duoc loi nay.
func TestCapIdBoss_AmVaDuyNhatQua1000LanGoiDongThoi(t *testing.T) {
	const soLan = 1000

	ids := make([]int64, soLan)
	batDau := make(chan struct{}) // vach xuat phat chung — don tranh chap vao mot diem
	var wg sync.WaitGroup
	wg.Add(soLan)
	for i := 0; i < soLan; i++ {
		go func(i int) {
			defer wg.Done()
			<-batDau
			ids[i] = CapIdBoss() // moi goroutine ghi mot o rieng -> khong tu tao dua tranh
		}(i)
	}
	close(batDau)
	wg.Wait()

	daThay := make(map[int64]int, soLan)
	for i, id := range ids {
		if id >= 0 {
			t.Fatalf("CapIdBoss lan thu %d tra id = %d — id boss PHAI AM", i, id)
		}
		if truoc, trung := daThay[id]; trung {
			t.Fatalf("id boss TRUNG: lan %d va lan %d cung tra %d — dung bay #2 thiet ke §7 (ban goc 5 boss cung id=-1)", truoc, i, id)
		}
		daThay[id] = i
	}
	if len(daThay) != soLan {
		t.Fatalf("chi co %d id khac nhau tren %d lan goi", len(daThay), soLan)
	}
}

// TestKepHp_ChanTranVaSoAm: KepHp phai kep ve dung mien [0, 2_147_483_647] —
// khong bao gio tra so am (bay #4 thiet ke §7: convert int64 vuot mien sang
// int32 la CAT BIT, ra so am vo nghia ma khong bao loi).
func TestKepHp_ChanTranVaSoAm(t *testing.T) {
	const (
		max32 int32 = 2_147_483_647
		min64 int64 = -9_223_372_036_854_775_808
		max64 int64 = 9_223_372_036_854_775_807
	)

	cacTruongHop := []struct {
		ten  string
		vao  int64
		muon int32
	}{
		{"so am nho", -1, 0},
		{"so am rat lon", min64, 0},
		{"san dung 0", 0, 0},
		{"gia tri thuong", 12_345, 12_345},
		{"dung tran", int64(max32), max32},
		{"tran hon 1", int64(max32) + 1, max32},
		{"tran xa (3 ty)", 3_000_000_000, max32},
		{"tran toi da int64", max64, max32},
	}

	for _, tc := range cacTruongHop {
		got := KepHp(tc.vao)
		if got != tc.muon {
			t.Errorf("%s: KepHp(%d) = %d, muon %d", tc.ten, tc.vao, got, tc.muon)
		}
		if got < 0 {
			t.Fatalf("%s: KepHp(%d) tra AM (%d) — dau hieu convert int64->int32 khong kep truoc", tc.ten, tc.vao, got)
		}
	}
}

// TestCongSatThuong_BoQuaGiaTriKhongDuong: don 0 hoac AM khong duoc ghi ten
// nguoi do vao so dong gop, va khong duoc TRU dong gop da co.
func TestCongSatThuong_BoQuaGiaTriKhongDuong(t *testing.T) {
	b := &Boss{}

	b.CongSatThuong(7, 0)
	b.CongSatThuong(7, -50)
	if len(b.SatThuongTheoNguoi) != 0 {
		t.Fatalf("sau hai don <= 0, so dong gop co %d muc: %v — muon rong", len(b.SatThuongTheoNguoi), b.SatThuongTheoNguoi)
	}

	b.CongSatThuong(7, 10)
	b.CongSatThuong(7, 5)
	if got := b.SatThuongTheoNguoi[7]; got != 15 {
		t.Fatalf("cong don 10 + 5 = %d, muon 15", got)
	}

	b.CongSatThuong(7, -100)
	if got := b.SatThuongTheoNguoi[7]; got != 15 {
		t.Errorf("don AM da TRU dong gop: con %d, muon giu nguyen 15", got)
	}
}

// TestCongSatThuong_KhongTranInt64: cong don gan tran PHAI bao hoa o MaxInt64,
// KHONG duoc cuon vong ve so am (tong AM se lam ChiaThuongTheoDongGop tinh sai
// toan bo ti le).
func TestCongSatThuong_KhongTranInt64(t *testing.T) {
	const max64 int64 = 9_223_372_036_854_775_807

	b := &Boss{SatThuongTheoNguoi: map[int64]int64{7: max64 - 5}}
	b.CongSatThuong(7, 100)

	got := b.SatThuongTheoNguoi[7]
	if got < 0 {
		t.Fatalf("cong don TRAN thanh so am: %d", got)
	}
	if got != max64 {
		t.Errorf("cong don gan tran = %d, muon bao hoa o %d", got, max64)
	}
}

// TestCongSatThuong_BossNilKhongPanic: ham nam tren duong chien dau, goi tren
// con tro nil khong duoc lam sap actor khu vuc.
func TestCongSatThuong_BossNilKhongPanic(t *testing.T) {
	var b *Boss
	b.CongSatThuong(1, 100) // khong panic la dat
}

// TestChiaThuongTheoDongGop_DungTiLe: chia dung ti le dong gop, lam tron XUONG.
func TestChiaThuongTheoDongGop_DungTiLe(t *testing.T) {
	// 50/30/20 tren tong 100, tongDiem 1000 -> 500/300/200 (chia het, khong
	// dinh lam tron).
	got := ChiaThuongTheoDongGop(map[int64]int64{1: 50, 2: 30, 3: 20}, 1000)
	muon := map[int64]int64{1: 500, 2: 300, 3: 200}
	soSanhMap(t, "chia het", got, muon)

	// 1/2 tren tong 3, tongDiem 10 -> floor(10*1/3)=3, floor(10*2/3)=6. Tong
	// phan = 9 < 10: lam tron XUONG, khong bao gio phat qua tongDiem.
	got = ChiaThuongTheoDongGop(map[int64]int64{1: 1, 2: 2}, 10)
	muon = map[int64]int64{1: 3, 2: 6}
	soSanhMap(t, "lam tron xuong", got, muon)

	// Mot minh danh boss -> lay tron tongDiem.
	got = ChiaThuongTheoDongGop(map[int64]int64{9: 12_345}, 777)
	muon = map[int64]int64{9: 777}
	soSanhMap(t, "mot nguoi duy nhat", got, muon)
}

// TestChiaThuongTheoDongGop_DongGopKhongLoKhongTran: hai nguoi cung dong gop 4
// ty ty (4e18) — neu ham tinh "tongDiem × gop" bang int64 thuong thi tich TRAN
// ngay, cho ra phan thuong am hoac be xiu. Ket qua dung van la 50/50.
func TestChiaThuongTheoDongGop_DongGopKhongLoKhongTran(t *testing.T) {
	const gopKhongLo int64 = 4_000_000_000_000_000_000

	got := ChiaThuongTheoDongGop(map[int64]int64{1: gopKhongLo, 2: gopKhongLo}, 1000)
	muon := map[int64]int64{1: 500, 2: 500}
	soSanhMap(t, "dong gop khong lo", got, muon)
}

// TestChiaThuongTheoDongGop_TongBangKhongTraMapRong: khong ai dong gop thi
// khong chia gi ca — tra map RONG (khac nil, de nguoi goi range thang), va
// TUYET DOI khong duoc panic vi chia cho 0.
func TestChiaThuongTheoDongGop_TongBangKhongTraMapRong(t *testing.T) {
	cacTruongHop := []struct {
		ten      string
		tong     map[int64]int64
		tongDiem int64
	}{
		{"map nil", nil, 1000},
		{"map rong", map[int64]int64{}, 1000},
		{"tat ca dong gop = 0", map[int64]int64{1: 0, 2: 0}, 1000},
		{"dong gop am", map[int64]int64{1: -50, 2: 0}, 1000},
		{"tongDiem = 0", map[int64]int64{1: 50}, 0},
		{"tongDiem am", map[int64]int64{1: 50}, -100},
	}

	for _, tc := range cacTruongHop {
		got := ChiaThuongTheoDongGop(tc.tong, tc.tongDiem)
		if got == nil {
			t.Errorf("%s: tra nil, muon map RONG", tc.ten)
			continue
		}
		if len(got) != 0 {
			t.Errorf("%s: tra %v, muon map rong", tc.ten, got)
		}
	}
}

// TestPhaHienTaiData_NilSafe: PhaHienTai chay theo may trang thai nen "ngoai
// bien" la trang thai BINH THUONG cua vong doi — phai tra nil, khong panic.
func TestPhaHienTaiData_NilSafe(t *testing.T) {
	var bossNil *Boss
	if got := bossNil.PhaHienTaiData(); got != nil {
		t.Errorf("boss nil: PhaHienTaiData = %v, muon nil", got)
	}

	cacTruongHop := []struct {
		ten string
		b   *Boss
	}{
		{"chua co du lieu pha", &Boss{PhaHienTai: 0}},
		{"chi so am", &Boss{PhaHienTai: -1, DuLieu: []PhaBoss{{Ten: "Pha 1"}}}},
		{"chi so vuot cuoi", &Boss{PhaHienTai: 5, DuLieu: []PhaBoss{{Ten: "Pha 1"}, {Ten: "Pha 2"}}}},
		{"chi so dung bang len", &Boss{PhaHienTai: 2, DuLieu: []PhaBoss{{Ten: "Pha 1"}, {Ten: "Pha 2"}}}},
	}
	for _, tc := range cacTruongHop {
		if got := tc.b.PhaHienTaiData(); got != nil {
			t.Errorf("%s: PhaHienTaiData = %+v, muon nil", tc.ten, got)
		}
	}

	// Pha hop le: tra dung phan tu, va la CON TRO VAO lat cat (sua qua con tro
	// thi du lieu goc doi theo) — may trang thai dua vao dieu nay de cap nhat
	// tai cho ma khong phai ghi lai ca phan tu.
	b := &Boss{PhaHienTai: 1, DuLieu: []PhaBoss{{Ten: "Pha 1", Cap: 10}, {Ten: "Pha 2", Cap: 20}}}
	pha := b.PhaHienTaiData()
	if pha == nil {
		t.Fatalf("pha hop le: PhaHienTaiData = nil, muon pha thu 1")
	}
	if pha.Ten != "Pha 2" || pha.Cap != 20 {
		t.Fatalf("PhaHienTaiData tra %+v, muon pha thu 1 (Ten=Pha 2, Cap=20)", pha)
	}
	pha.Cap = 99
	if b.DuLieu[1].Cap != 99 {
		t.Errorf("PhaHienTaiData tra BAN SAO (sua qua con tro khong doi du lieu goc) — muon con tro vao lat cat")
	}
}

// TestPhaBoss_CapHopLe: goi -5 truong 3 ghi mot byte cap; cap 0 la vo nghia voi
// client nen KHONG hop le. nil coi nhu khong hop le.
func TestPhaBoss_CapHopLe(t *testing.T) {
	var phaNil *PhaBoss
	if phaNil.CapHopLe() {
		t.Errorf("pha nil: CapHopLe = true, muon false")
	}

	cacTruongHop := []struct {
		cap  byte
		muon bool
	}{
		{0, false},
		{1, true},
		{50, true},
		{255, true},
	}
	for _, tc := range cacTruongHop {
		p := &PhaBoss{Cap: tc.cap}
		if got := p.CapHopLe(); got != tc.muon {
			t.Errorf("PhaBoss{Cap: %d}.CapHopLe() = %v, muon %v", tc.cap, got, tc.muon)
		}
	}
}

// TestTrangThaiBoss_String ghim ten sau trang thai DUNG NGUYEN VAN thiet ke §5
// — log va test doc thang ten nay, doi ten la vo cho khac.
func TestTrangThaiBoss_String(t *testing.T) {
	cacTruongHop := []struct {
		tt   TrangThaiBoss
		muon string
	}{
		{BossNghi, "NGHI"},
		{BossHien, "HIEN"},
		{BossChao, "CHAO"},
		{BossHoatDong, "HOAT_DONG"},
		{BossTroi, "TROI"},
		{BossRaDi, "RA_DI"},
		{TrangThaiBoss(200), "KHONG_RO"},
	}
	for _, tc := range cacTruongHop {
		if got := tc.tt.String(); got != tc.muon {
			t.Errorf("TrangThaiBoss(%d).String() = %q, muon %q", byte(tc.tt), got, tc.muon)
		}
	}

	// Zero value PHAI la NGHI: mot Boss{} moi tao chua xuat hien o dau ca.
	var b Boss
	if b.TrangThai != BossNghi {
		t.Errorf("zero value TrangThai = %v, muon NGHI", b.TrangThai)
	}
}

// soSanhMap so sanh hai map ket qua chia thuong, bao ro khoa nao lech.
func soSanhMap(t *testing.T, ten string, got, muon map[int64]int64) {
	t.Helper()
	if len(got) != len(muon) {
		t.Errorf("%s: ket qua co %d muc (%v), muon %d muc (%v)", ten, len(got), got, len(muon), muon)
		return
	}
	for id, phanMuon := range muon {
		if phan, co := got[id]; !co || phan != phanMuon {
			t.Errorf("%s: phan cua id=%d la %d (co=%v), muon %d", ten, id, phan, co, phanMuon)
		}
	}
}
