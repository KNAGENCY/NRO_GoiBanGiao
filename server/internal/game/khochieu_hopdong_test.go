// khochieu_hopdong_test.go — KHOA HOP DONG "MOT LAN DAU NHIP" cua KH8 task 3.
//
// Hop dong (internal/noidung/kho.go, dau tep; thiet ke §4): actor lay anh chup
// noi dung DUNG MOT LAN o dau nhip va giu no suot nhip. Truoc file nay, hop dong
// do chi la mot doan chu thich — khong co gi ngan mot nguoi sau day nay them mot
// loi goi AnhChup() thu hai vao aiBoss.
//
// Bon test o day bien no thanh mot loi hua CO THE GIET DUOC:
//
//  1. dem: mot nhip = DUNG 1 lan lay anh chup, du trong nhip co bao nhieu boss.
//  2. su kien: moi lan xu ly su kien cung lay dung 1 lan.
//  3. the he: Thay() GIUA nhip khong doi ket qua trong chinh nhip do.
//  4. mac dinh: chua gan kho thi doc ban NHUNG (KHONG rong, KHONG bang Go cu).
package game

import (
	"testing"
	"time"

	"nroserver/internal/noidung"
)

// khoChieuDem la mot NGUON noi dung dem so lan bi hoi anh chup. Boc mot
// *noidung.Kho THAT (khong gia du lieu) — cai duy nhat no them la cai dem.
type khoChieuDem struct {
	kho   *noidung.Kho
	soLan int
}

func (d *khoChieuDem) nguon() *nguonNoiDung {
	return &nguonNoiDung{lay: func() khoChieuAnhChup {
		d.soLan++
		return d.kho.AnhChup()
	}}
}

func moKhoChieuDem(t *testing.T) *khoChieuDem {
	t.Helper()
	du, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("nap kho noi dung: %v", err)
	}
	return &khoChieuDem{kho: noidung.MoKho(du)}
}

// khuVucDemChieu dung mot khu vuc co MOT nguoi choi va MOT boss dang HOAT DONG
// nam dung cho nguoi choi, de nhipTim chac chan chay het duong
// nhipBoss -> aiBoss -> tamChieuXaNhatCuaPha -> bossThuDanh (ba cho tra cuu
// chieu). Khong co boss thi test dem se xanh vi ly do sai: khong ai hoi chieu.
func khuVucDemChieu(t *testing.T, d *khoChieuDem, luc *time.Time) *KhuVuc {
	t.Helper()
	k := khuVucAiBoss(luc)
	k.ganNguonNoiDung(d.nguon())
	themNguoiAiBoss(k, 1, xBossThuNghiem, yDiDuocBoss)
	k.boss = append(k.boss, bossChienDau(*luc, xBossThuNghiem, 0, nil))
	return k
}

// TestKhoChieuMotLanMoiNhip — DOT BIEN NHAM TOI: them mot loi goi
// k.layAnhChupNhip() (hay kho.AnhChup()) thu hai o bat ky dau trong duong nhip
// — aiBoss, bossThuDanh, tamChieuXaNhatCuaPha — se lam so lan len 2 va test nay
// do ngay.
func TestKhoChieuMotLanMoiNhip(t *testing.T) {
	luc := mocDanhBoss
	d := moKhoChieuDem(t)
	k := khuVucDemChieu(t, d, &luc)

	d.soLan = 0 // bo qua nhung lan cua buoc chuan bi
	k.nhipTim()
	if d.soLan != 1 {
		t.Fatalf("mot nhip lay anh chup %d lan, hop dong la DUNG 1 lan "+
			"(xem khoi chu thich dau internal/noidung/kho.go)", d.soLan)
	}

	// Nhip thu hai: lai dung 1 lan nua — KHONG phai 0. Nhip sau PHAI thay duoc
	// noi dung moi neu vua co nap nong; nho anh chup lai la nap nong chet.
	luc = luc.Add(time.Second)
	k.nhipTim()
	if d.soLan != 2 {
		t.Fatalf("hai nhip lay anh chup %d lan, mong doi 2", d.soLan)
	}
}

// TestKhoChieuMoiSuKienMotLan: mot su kien cung la mot don vi cong viec — chay
// tron ven tren mot the he so lieu, nen cung lay dung mot lan.
func TestKhoChieuMoiSuKienMotLan(t *testing.T) {
	luc := mocDanhBoss
	d := moKhoChieuDem(t)
	k := khuVucDemChieu(t, d, &luc)

	d.soLan = 0
	k.xuLy(suKienChonChieu{idNhanVat: 1, mauChieuId: 1})
	if d.soLan != 1 {
		t.Fatalf("mot su kien lay anh chup %d lan, mong doi 1", d.soLan)
	}
	if k.nguoiTrong[1].chieuDangChon == nil {
		t.Fatal("chuan bi test sai: chon chieu 1 cua lang 0 phai thanh cong")
	}

	// Su kien THU HAI phai lay anh chup MOI (2 lan cong don), khong dung lai
	// anh chup cua su kien truoc.
	//
	// DOT BIEN NHAM TOI: bo dong `k.acNhip = k.layAnhChupNhip()` o dau xuLy.
	// Chi dem MOT su kien thi dot bien do song (anhChupNoiDung() tu va khi
	// acNhip con nil); dem HAI su kien moi giet duoc no — mot su kien nap
	// nong xong ma su kien sau van chay so cu la nap nong hong.
	k.xuLy(suKienChonChieu{idNhanVat: 1, mauChieuId: 2})
	if d.soLan != 2 {
		t.Fatalf("hai su kien lay anh chup %d lan, mong doi 2 (moi su kien mot lan)", d.soLan)
	}
}

// TestKhoChieuThayGiuaNhipKhongDoiKetQuaTrongNhip la test QUAN TRONG NHAT cua
// task: no chung minh nap nong KHONG the xen ngang mot nhip.
//
// Cach lam: mot nguon dem boc mot kho THAT. Sau khi nhip da lay anh chup, ta
// goi Kho.Thay() voi mot bo so KHAC HAN. Neu o dau do trong nhip con mot loi
// goi AnhChup() rai rac, no se thay so moi — va test nay do.
func TestKhoChieuThayGiuaNhipKhongDoiKetQuaTrongNhip(t *testing.T) {
	d := moKhoChieuDem(t)
	k := MoKhuVuc(TimBanDo(1))
	k.ganNguonNoiDung(d.nguon())

	// Bat dau mot "nhip": lay anh chup dung mot lan, y nhu nhipTim lam.
	k.acNhip = k.layAnhChupNhip()
	theHeDau := k.acNhip.TheHe()

	truoc := traCapChieuTu(k.anhChupNoiDung(), 0, 1, 1)
	if truoc == nil {
		t.Fatal("chuan bi test sai: khong tra duoc cap chieu (lang 0, chieu 1, cap 1)")
	}

	// NAP NONG GIUA NHIP: mot bo so khac han.
	du2, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("nap lan hai: %v", err)
	}
	du2.ChieuTheoLang[0][0].CacCap[0].HeSoPhanTram = truoc.SatThuong + 999
	theHeMoi := d.kho.Thay(du2)
	if theHeMoi <= theHeDau {
		t.Fatalf("Thay khong tang the he: %d -> %d", theHeDau, theHeMoi)
	}

	// TRONG CUNG NHIP: phai van thay so CU.
	sau := traCapChieuTu(k.anhChupNoiDung(), 0, 1, 1)
	if sau == nil || sau.SatThuong != truoc.SatThuong {
		t.Fatalf("nap nong xen ngang mot nhip: sat thuong %d -> %v (phai giu nguyen trong nhip)",
			truoc.SatThuong, sau)
	}
	if got := k.anhChupNoiDung().TheHe(); got != theHeDau {
		t.Fatalf("the he doi GIUA nhip: %d -> %d", theHeDau, got)
	}

	// NHIP SAU: phai thay so MOI (neu khong thi nap nong vo tac dung).
	k.acNhip = k.layAnhChupNhip()
	moi := traCapChieuTu(k.anhChupNoiDung(), 0, 1, 1)
	if moi == nil || moi.SatThuong != truoc.SatThuong+999 {
		t.Fatalf("nhip sau khong thay so moi: %v", moi)
	}
}

// TestKhoChieuChuaGanThiDungBanNhung: khu vuc chua ai gan kho van phai co day
// du chieu (doc ban NHUNG) — KHONG roi ve bang Go (bang do da bi xoa o task 3)
// va KHONG tra ve rong.
func TestKhoChieuChuaGanThiDungBanNhung(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	ds := chieuCuaLangTu(k.anhChupNoiDung(), 0)
	if len(ds) != 4 {
		t.Fatalf("khu vuc chua gan kho: lang 0 co %d chieu, mong doi 4 (ban nhung)", len(ds))
	}
	if traCapChieuTu(k.anhChupNoiDung(), 0, ds[0].Id, 1) == nil {
		t.Fatal("khu vuc chua gan kho: khong tra duoc cap 1 cua chieu dau")
	}
}

// TestGanKhoNoiDungDoiDuocSoLieu khoa NUA SAU cua duong noi day: cai kho ma
// main.go cam vao (qua SoKhuVuc.GanKhoNoiDung) PHAI la cai kho ma actor doc.
//
// DOT BIEN NHAM TOI: bo dong soKhuVuc.GanKhoNoiDung(...) trong main.go, hoac
// lam GanKhoNoiDung khong luu gi — khu vuc se roi ve ban nhung va khong bao gio
// thay so da nap nong; test nay do.
func TestGanKhoNoiDungDoiDuocSoLieu(t *testing.T) {
	du, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("nap kho noi dung: %v", err)
	}
	goc := du.ChieuTheoLang[0][0].CacCap[0].HeSoPhanTram
	du.ChieuTheoLang[0][0].CacCap[0].HeSoPhanTram = goc + 777
	kho := noidung.MoKho(du)

	s := MoSoKhuVuc()
	defer s.DungTatCa()
	s.GanKhoNoiDung(kho)

	k := s.Tra(1)
	if k == nil {
		t.Fatal("khong lay duoc khu vuc mapId 1")
	}
	cap := traCapChieuTu(k.layAnhChupNhip(), 0, 1, 1)
	if cap == nil || cap.SatThuong != goc+777 {
		t.Fatalf("khu vuc khong doc kho da gan: %v (mong doi SatThuong %d)", cap, goc+777)
	}
}
