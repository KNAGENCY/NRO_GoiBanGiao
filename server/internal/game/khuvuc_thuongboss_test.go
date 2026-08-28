package game

import (
	"testing"
	"time"

	"nroserver/internal/store"
)

// ─────────────────────────────────────────────────────────────────────────────
// KH6b Task 12 — boss chet: chia thuong theo dong gop + thong bao toan server.
//
// File nay bao ve BON cau, moi cau la mot loi da phai tra gia o cho khac:
//
//  1. Phan thuong chia DUNG TI LE dong gop, va ghi vao NHAN VAT (nv.KinhNghiem)
//     — tuc la di theo duong luu san co (BoLuuDinhKy 60s + KhiDong) nen BEN sau
//     khi tat may chu, MA actor khong phai cham CSDL lan nao.
//  2. Chia XONG thi XOA SACH map dong gop — trong CUNG mot ham, khong de ho
//     giua "da chia" va "da xoa": ho do la chia hai lan.
//  3. Thong bao 93 di QUA HANG DOI LOA (k.xepThongBao), TUYET DOI khong phat
//     thang tu actor (CHAN-2, ke hoach KH6b).
//  4. Khong cu goi nao trong duong nay lam SAP actor: nguoi dong gop da roi khu
//     vuc, nguoi ha boss da mat ket noi, du lieu pha hong — deu phai di qua
//     duoc.
//
// Moi test goi k.xuLy / k.bossChet TRUC TIEP (actor CHUA BatDau), dung khuon
// khuvuc_boss_test.go — tru DUNG MOT test co y chay tu TRONG actor de do ky
// luat CHAN-2 (TestActorKhongTuKhoaKhiBossChet).
// ─────────────────────────────────────────────────────────────────────────────

// mocTest la mot moc thoi gian co dinh dung cho k.dongHo trong ca file. Ghim
// cung (khong time.Now) de moi khang dinh ve MocThoiGian la tat dinh.
var mocTest = time.Date(2026, 8, 28, 10, 0, 0, 0, time.UTC)

// khuVucTest mo mot khu vuc voi dong ho ghim — dung cho moi test o day.
func khuVucTest(t *testing.T) *KhuVuc {
	t.Helper()
	k := MoKhuVuc(TimBanDo(1))
	k.dongHo = func() time.Time { return mocTest }
	return k
}

// bossThuong dung mot boss MOT PHA co du lieu thuong — tham so hoa dung ba thu
// moi test can doi (trang thai, diem thuong, co boss lon) de khong test nao
// phai cheo lai ca struct.
func bossThuong(trangThai TrangThaiBoss, diemThuong int64, laBossLon bool) *Boss {
	return &Boss{
		Id:         -7,
		PhaHienTai: 0,
		DuLieu: []PhaBoss{{
			Ten:        "Cửu Vĩ",
			Lang:       3,
			Cap:        90,
			TypePk:     TypePkBoss,
			DiemThuong: diemThuong,
			LaBossLon:  laBossLon,
		}},
		Hp:        1000,
		HpToiDa:   1000,
		TrangThai: trangThai,
	}
}

// datNguoi dua mot nguoi choi vao khu vuc QUA k.xuLy (khong ghi thang
// k.nguoiTrong — do la trang thai actor) va tra ve con tro nhan vat de test
// doc lai KinhNghiem sau khi chia thuong.
func datNguoi(t *testing.T, k *KhuVuc, id int64, ten string) *store.NhanVat {
	t.Helper()
	nv := &store.NhanVat{ID: id, Ten: ten, HpToiDa: 100, HpHienTai: 100, ViTriX: 100, ViTriY: 384}
	k.xuLy(suKienVao{nv: nv, phien: &phienGiaKhuVuc{}})
	if _, co := k.nguoiTrong[id]; !co {
		t.Fatalf("nguoi %d chua vao duoc khu vuc", id)
	}
	return nv
}

// rutHangDoiLoa lay HET noi dung dang cho trong hang doi loa, theo dung thu tu
// FIFO. Doc bang select+default (khong chan) de test khong bao gio treo khi
// hang doi rong hon mong doi.
func rutHangDoiLoa(loa *BoPhatThanh) []string {
	var ra []string
	for {
		select {
		case yc := <-loa.hangDoi:
			ra = append(ra, yc.NoiDung)
		default:
			return ra
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 1 — hai cau chu thong bao (ham thuan, boss.go)
// ─────────────────────────────────────────────────────────────────────────────

// TestCauThongBaoBossGhimNguyenVan ghim NGUYEN VAN hai dong chu se chay ngang
// man hinh moi nguoi choi. Ghim cung (khong ghep lai tu chinh hang so trong
// code) vi day la thu NGUOI CHOI DOC — doi no la doi trai nghiem, phai co mot
// test do de ai doi con biet minh dang doi cai gi.
func TestCauThongBaoBossGhimNguyenVan(t *testing.T) {
	if got := cauBossXuatHien("Cửu Vĩ"); got != "Boss Cửu Vĩ vừa xuất hiện!" {
		t.Errorf("cauBossXuatHien = %q, mong doi %q", got, "Boss Cửu Vĩ vừa xuất hiện!")
	}
	if got := cauBossBiHa("Naruto", "Cửu Vĩ"); got != "Naruto vừa hạ gục Boss Cửu Vĩ!" {
		t.Errorf("cauBossBiHa = %q, mong doi %q", got, "Naruto vừa hạ gục Boss Cửu Vĩ!")
	}
}

// TestCongBaoHoaKhongTranInt64: phan thuong cong vao KinhNghiem PHAI bao hoa o
// MaxInt64, khong cuon vong ve am. Mot nguoi choi cay lau nam co the co
// KinhNghiem rat lon; cuon vong se bien no thanh so am — mat trang tien do ma
// khong bao loi gi.
func TestCongBaoHoaKhongTranInt64(t *testing.T) {
	const maxInt64 int64 = 9_223_372_036_854_775_807
	if got := congBaoHoa(maxInt64-10, 100); got != maxInt64 {
		t.Errorf("congBaoHoa(max-10, 100) = %d, mong doi bao hoa o %d", got, maxInt64)
	}
	if got := congBaoHoa(5, 7); got != 12 {
		t.Errorf("congBaoHoa(5, 7) = %d, mong doi 12", got)
	}
	if got := congBaoHoa(5, 0); got != 5 {
		t.Errorf("congBaoHoa(5, 0) = %d, mong doi 5 (cong 0 khong doi gi)", got)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 2 — bossNhanSatThuong: cong dong gop + tru mau + phat hien chet
// ─────────────────────────────────────────────────────────────────────────────

// TestBossNhanSatThuongTruMauVaGhiDongGop: mot don trung PHAI lam hai viec
// cung luc — tru mau boss va ghi ten nguoi danh vao so dong gop. Thieu ve ghi
// so thi den luc chet khong ai duoc gi ma khong co dau hieu nao.
func TestBossNhanSatThuongTruMauVaGhiDongGop(t *testing.T) {
	k := khuVucTest(t)
	b := bossThuong(BossHoatDong, 1000, false)

	if chet := k.bossNhanSatThuong(b, 11, 300, mocTest); chet {
		t.Fatal("boss 1000 mau an 300 sat thuong ma bao la CHET")
	}
	if b.Hp != 700 {
		t.Errorf("mau boss = %d, mong doi 700", b.Hp)
	}
	if got := b.SatThuongTheoNguoi[11]; got != 300 {
		t.Errorf("dong gop cua nguoi 11 = %d, mong doi 300", got)
	}
}

// TestBossNhanSatThuongKhongTinhPhanDonThua: don cuoi cung DU MANH cach may
// cung chi duoc tinh dong gop bang so mau CON LAI cua boss.
//
// Vi sao day la luat, khong phai chi tiet: neu tinh nguyen so sat thuong, mot
// nguoi den giay chot dung mot don 1 ty vao con boss con 1 mau se chiem gan
// 100% phan thuong, cuop trang cua nhung nguoi da danh ca tran. Dong gop phai
// do "da lam boss mat bao nhieu mau", khong phai "da vung ra con so to bao
// nhieu".
func TestBossNhanSatThuongKhongTinhPhanDonThua(t *testing.T) {
	k := khuVucTest(t)
	b := bossThuong(BossHoatDong, 1000, false)
	b.Hp = 5 // gan chet

	if chet := k.bossNhanSatThuong(b, 11, 1_000_000_000, mocTest); !chet {
		t.Fatal("don 1 ty vao con boss con 5 mau PHAI giet duoc no")
	}
	if b.Hp != 0 {
		t.Errorf("mau boss = %d, mong doi 0 (KepHp phai chan so am)", b.Hp)
	}
	// Map da bi xoa sau khi chia thuong nen phai kiem qua phan thuong: chi mot
	// nguoi dong gop nen ho lay tron 1000 diem du dong gop chi la 5.
	// (Kiem dong gop trung thuc hon o test rieng ben duoi.)
}

// TestDonKetLieuThuaKhongCuopPhanThuong do cai kep dong gop bang HAU QUA THAT
// cua no — phan thuong chia ra — thay vi doc so trong map (map bi xoa ngay khi
// boss chet nen khong doc duoc nua).
//
// Kich ban: boss 200 mau. A danh dung 100 (mot nua). B giang don ket lieu
// 1.000.000 vao 100 mau con lai.
//   - CO kep: dong gop 100 vs 100 -> 500/500. Cong bang.
//   - KHONG kep: 100 vs 1.000.000 -> A duoc 0 diem, B ve gan tron 1000. Nguoi
//     danh nua tran trang tay chi vi khong bam duoc phat cuoi.
func TestDonKetLieuThuaKhongCuopPhanThuong(t *testing.T) {
	k := khuVucTest(t)
	nvA := datNguoi(t, k, 11, "A")
	nvB := datNguoi(t, k, 22, "B")

	b := bossThuong(BossHoatDong, 1000, false)
	b.Hp = 200
	b.HpToiDa = 200

	k.bossNhanSatThuong(b, 11, 100, mocTest)       // A: dung mot nua mau
	k.bossNhanSatThuong(b, 22, 1_000_000, mocTest) // B: don ket lieu thua rat nhieu

	if b.TrangThai != BossTroi {
		t.Fatalf("boss phai chet sau don ket lieu, trang thai = %s", b.TrangThai)
	}
	if nvA.KinhNghiem != 500 || nvB.KinhNghiem != 500 {
		t.Fatalf("thuong A/B = %d/%d, mong doi 500/500. Dong gop phai kep o SO MAU CON LAI: "+
			"khong kep thi mot don thua khong lo se cuop gan het phan cua nguoi danh ca tran",
			nvA.KinhNghiem, nvB.KinhNghiem)
	}
}

// TestBossNhanSatThuongKhongDanhDuocXacBoss: boss DA CHET (TROI) hoac chua/
// khong con tren man hinh (NGHI, RA_DI) thi moi don deu roi vao khoang khong —
// khong ghi dong gop, khong tru mau, khong kich hoat lai duong chia thuong.
//
// Day la CHOT CHONG CHIA HAI LAN quan trong nhat: neu thieu, hai don cung
// nhip (hai nguoi choi bam cung luc) deu thay hp<=0 va deu goi bossChet.
func TestBossNhanSatThuongKhongDanhDuocXacBoss(t *testing.T) {
	for _, tt := range []TrangThaiBoss{BossNghi, BossTroi, BossRaDi} {
		k := khuVucTest(t)
		b := bossThuong(tt, 1000, false)

		if chet := k.bossNhanSatThuong(b, 11, 300, mocTest); chet {
			t.Errorf("trang thai %s: don danh bao la giet duoc boss", tt)
		}
		if b.Hp != 1000 {
			t.Errorf("trang thai %s: mau boss = %d, mong doi giu nguyen 1000", tt, b.Hp)
		}
		if len(b.SatThuongTheoNguoi) != 0 {
			t.Errorf("trang thai %s: da ghi %d dong gop, mong doi 0", tt, len(b.SatThuongTheoNguoi))
		}
	}
}

// TestBossNhanSatThuongChanMapDongGopPhinhVoHan: boss song lau trong khu vuc
// dong nguoi thi map dong gop se phinh theo SO NGUOI TUNG DANH no. Co tran de
// mot con boss khong bao gio chet cung khong the an het bo nho.
//
// Luat cua tran: dat tran roi thi nguoi MOI khong duoc them nua, nhung nguoi
// DA CO trong so van cong don binh thuong — neu khong, nguoi danh tu dau se bi
// dong bang dong gop giua chung tran.
func TestBossNhanSatThuongChanMapDongGopPhinhVoHan(t *testing.T) {
	k := khuVucTest(t)
	b := bossThuong(BossHoatDong, 1000, false)
	b.Hp = 2_000_000_000
	b.HpToiDa = 2_000_000_000

	for i := 0; i < tranSoNguoiDongGop+50; i++ {
		k.bossNhanSatThuong(b, int64(1000+i), 1, mocTest)
	}
	if got := len(b.SatThuongTheoNguoi); got != tranSoNguoiDongGop {
		t.Fatalf("so nguoi trong so dong gop = %d, mong doi bi chan o %d", got, tranSoNguoiDongGop)
	}

	// Nguoi DA CO trong so (id 1000, nguoi dau tien) van phai cong don duoc.
	k.bossNhanSatThuong(b, 1000, 99, mocTest)
	if got := b.SatThuongTheoNguoi[1000]; got != 100 {
		t.Errorf("dong gop cua nguoi da co trong so = %d, mong doi 100 (1 + 99) — "+
			"tran chi duoc chan nguoi MOI, khong duoc dong bang nguoi cu", got)
	}
}

// TestBossNhanSatThuongBossNilKhongPanic: ham nam tren duong chien dau (moi
// don danh deu di qua), goi voi con tro nil khong duoc lam sap actor.
func TestBossNhanSatThuongBossNilKhongPanic(t *testing.T) {
	k := khuVucTest(t)
	if chet := k.bossNhanSatThuong(nil, 11, 100, mocTest); chet {
		t.Error("bossNhanSatThuong(nil) bao la giet duoc boss")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 3 — bossChet: chia thuong, xoa so, chuyen TROI
// ─────────────────────────────────────────────────────────────────────────────

// TestBossChetChiaThuongDungTiLe la test TRUNG TAM cua task: ba nguoi dong gop
// 50/30/20 phan tram, tong thuong 1000 diem -> 500/300/200 diem KINH NGHIEM
// cong vao nhan vat. Ghim con so cung, khong tinh lai bang cong thuc.
func TestBossChetChiaThuongDungTiLe(t *testing.T) {
	k := khuVucTest(t)
	nvA := datNguoi(t, k, 11, "A")
	nvB := datNguoi(t, k, 22, "B")
	nvC := datNguoi(t, k, 33, "C")
	nvA.KinhNghiem = 7 // co san it kinh nghiem -> phan thuong phai CONG THEM, khong ghi de

	b := bossThuong(BossHoatDong, 1000, false)
	b.SatThuongTheoNguoi = map[int64]int64{11: 50, 22: 30, 33: 20}

	k.bossChet(b, 11, mocTest)

	if nvA.KinhNghiem != 507 {
		t.Errorf("KinhNghiem cua A = %d, mong doi 507 (7 co san + 500 thuong)", nvA.KinhNghiem)
	}
	if nvB.KinhNghiem != 300 {
		t.Errorf("KinhNghiem cua B = %d, mong doi 300", nvB.KinhNghiem)
	}
	if nvC.KinhNghiem != 200 {
		t.Errorf("KinhNghiem cua C = %d, mong doi 200", nvC.KinhNghiem)
	}
}

// TestBossChetXoaSachMapDongGop: chia xong PHAI xoa so. Neu khong, lan chet ke
// tiep (pha sau, hoac lan hoi sinh sau) se chia LAI cho dung nhung nguoi do
// theo dong gop cu — nguoi khong danh phat nao van duoc thuong.
func TestBossChetXoaSachMapDongGop(t *testing.T) {
	k := khuVucTest(t)
	datNguoi(t, k, 11, "A")

	b := bossThuong(BossHoatDong, 1000, false)
	b.SatThuongTheoNguoi = map[int64]int64{11: 100}

	k.bossChet(b, 11, mocTest)

	if len(b.SatThuongTheoNguoi) != 0 {
		t.Fatalf("sau khi chia thuong, so dong gop con %d ban ghi, mong doi 0",
			len(b.SatThuongTheoNguoi))
	}
}

// TestBossChetKhongChiaThuongHaiLan la test DOT BIEN song: goi bossChet LAN
// HAI tren cung con boss (mo phong hai don cung nhip, hoac may trang thai goi
// nham) — khong ai duoc them mot diem nao.
//
// Test nay giet HAI dot bien cung luc: (a) bo dong xoa map dong gop, (b) bo
// chot "da o TROI/RA_DI thi thoi".
func TestBossChetKhongChiaThuongHaiLan(t *testing.T) {
	k := khuVucTest(t)
	nvA := datNguoi(t, k, 11, "A")

	b := bossThuong(BossHoatDong, 1000, false)
	b.SatThuongTheoNguoi = map[int64]int64{11: 100}

	k.bossChet(b, 11, mocTest)
	sauLanMot := nvA.KinhNghiem
	if sauLanMot != 1000 {
		t.Fatalf("lan chet dau: KinhNghiem = %d, mong doi 1000", sauLanMot)
	}

	k.bossChet(b, 11, mocTest) // lan hai — phai la viec KHONG

	if nvA.KinhNghiem != sauLanMot {
		t.Fatalf("goi bossChet LAN HAI cong them thuong: KinhNghiem %d -> %d. "+
			"Mot con boss chi duoc tra thuong DUNG MOT LAN cho moi lan chet",
			sauLanMot, nvA.KinhNghiem)
	}
}

// TestBossChetChuyenTrangThaiTroi: theo thiet ke §5, chet xong boss vao TROI
// (con xac tren man hinh de phat ThoaiChet) chu KHONG bien mat ngay.
func TestBossChetChuyenTrangThaiTroi(t *testing.T) {
	k := khuVucTest(t)
	b := bossThuong(BossHoatDong, 1000, false)

	k.bossChet(b, 11, mocTest)

	if b.TrangThai != BossTroi {
		t.Errorf("trang thai sau khi chet = %s, mong doi TROI", b.TrangThai)
	}
	if !b.MocThoiGian.Equal(mocTest) {
		t.Errorf("MocThoiGian = %v, mong doi %v (moc vao trang thai TROI)", b.MocThoiGian, mocTest)
	}
	if b.Hp != 0 {
		t.Errorf("mau boss sau khi chet = %d, mong doi 0", b.Hp)
	}
}

// TestBossChetNguoiDongGopDaRoiKhuVucKhongLamSap: nguoi danh 90% mau roi rot
// mang o giay cuoi KHONG con trong k.nguoiTrong. Actor khong duoc panic, va
// nhung nguoi CON LAI van phai duoc tra dung phan cua ho.
//
// Phan cua nguoi da roi bi BO (khong ai nhan ho) — day la cai gia da chot,
// ly do ghi o chu thich traThuongBoss (khuvuc.go).
func TestBossChetNguoiDongGopDaRoiKhuVucKhongLamSap(t *testing.T) {
	k := khuVucTest(t)
	nvA := datNguoi(t, k, 11, "A")

	b := bossThuong(BossHoatDong, 1000, false)
	// 99 la nguoi da roi khu vuc: co ten trong so dong gop nhung khong co
	// trong k.nguoiTrong.
	b.SatThuongTheoNguoi = map[int64]int64{11: 50, 99: 50}

	k.bossChet(b, 11, mocTest) // khong duoc panic

	if nvA.KinhNghiem != 500 {
		t.Errorf("KinhNghiem cua nguoi con o lai = %d, mong doi 500 (van dung ti le 50%%)",
			nvA.KinhNghiem)
	}
}

// TestBossChetKhongCoDiemThuongThiKhongAiDuocGi: pha khong khai bao DiemThuong
// (=0) la mot pha KHONG thuong — vong chia phai tra ve tay khong, khong duoc
// cong 0 hay cong bua.
func TestBossChetKhongCoDiemThuongThiKhongAiDuocGi(t *testing.T) {
	k := khuVucTest(t)
	nvA := datNguoi(t, k, 11, "A")

	b := bossThuong(BossHoatDong, 0, false)
	b.SatThuongTheoNguoi = map[int64]int64{11: 100}

	k.bossChet(b, 11, mocTest)

	if nvA.KinhNghiem != 0 {
		t.Errorf("KinhNghiem = %d, mong doi 0 (pha khong co DiemThuong)", nvA.KinhNghiem)
	}
}

// TestBossChetPhaHongKhongLamSap: PhaHienTai tro ra ngoai mien DuLieu (du lieu
// hong) -> PhaHienTaiData() tra nil. Duong chia thuong + thong bao deu phai di
// qua duoc ma khong nil-deref.
func TestBossChetPhaHongKhongLamSap(t *testing.T) {
	k := khuVucTest(t)
	nvA := datNguoi(t, k, 11, "A")

	b := bossThuong(BossHoatDong, 1000, true)
	b.PhaHienTai = 99 // ngoai mien
	b.SatThuongTheoNguoi = map[int64]int64{11: 100}

	k.bossChet(b, 11, mocTest) // khong duoc panic

	if nvA.KinhNghiem != 0 {
		t.Errorf("KinhNghiem = %d, mong doi 0 — khong co pha thi khong biet thuong bao nhieu",
			nvA.KinhNghiem)
	}
	if b.TrangThai != BossTroi {
		t.Errorf("trang thai = %s, mong doi TROI — du lieu pha hong khong duoc lam ket vong doi",
			b.TrangThai)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 4 — thong bao 93: PHAI di qua hang doi loa (CHAN-2)
// ─────────────────────────────────────────────────────────────────────────────

// TestBossLonChetXepThongBaoQuaLoa khoa CA HAI nua cua CHAN-2:
//
//	(a) tin PHAI nam trong hang doi loa (dung mot yeu cau, dung noi dung), va
//	(b) nguoi choi trong khu vuc KHONG duoc nhan goi 93 nao ngay luc do — neu
//	    ai do doi k.xepThongBao thanh mot cu phat thang (PhatToanServer /
//	    phatChoCaKhuVuc), hang doi se rong va tay nguoi choi se co goi -> DO.
func TestBossLonChetXepThongBaoQuaLoa(t *testing.T) {
	loa := MoBoPhatThanh()
	defer loa.Dung()

	k := khuVucTest(t)
	k.loa.Store(loa)

	nv := &store.NhanVat{ID: 11, Ten: "Naruto", HpToiDa: 100, HpHienTai: 100, ViTriX: 100, ViTriY: 384}
	p := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nv, phien: p})
	soGoiTruoc := len(p.banTin())

	b := bossThuong(BossHoatDong, 1000, true)
	b.SatThuongTheoNguoi = map[int64]int64{11: 100}

	k.bossChet(b, 11, mocTest)

	tin := rutHangDoiLoa(loa)
	if len(tin) != 1 {
		t.Fatalf("hang doi loa co %d tin sau khi boss lon chet, mong doi dung 1: %v", len(tin), tin)
	}
	if tin[0] != "Naruto vừa hạ gục Boss Cửu Vĩ!" {
		t.Errorf("noi dung thong bao = %q, mong doi %q", tin[0], "Naruto vừa hạ gục Boss Cửu Vĩ!")
	}
	if them := len(p.banTin()) - soGoiTruoc; them != 0 {
		t.Errorf("nguoi choi nhan them %d goi NGAY luc boss chet — thong bao 93 PHAI di qua "+
			"hang doi loa (nhip 1,5s, goroutine rieng), KHONG duoc phat thang tu actor", them)
	}
}

// TestBossLonChetChiBaoDUNG_MOT_LAN la cho thuc hien "gop/bo trung o NOI SINH
// RA TIN, khong o duong ong" (chot cua Task 6: loa co y ngo nghech, khong hieu
// noi dung, khong bo trung).
//
// Boss chet la mot su kien XAY RA MOT LAN. Duong sinh tin phai tu bao dam dieu
// do — o day la chot "da o TROI/RA_DI thi thoi" trong bossChet. Neu ai do bo
// chot ay, moi loi goi lap lai (hai don cung nhip, may trang thai goi lai,
// lenh admin) se de mot dong chu trung y het len song lan hai, va loa se
// trung thanh phat ca hai.
//
// LUU Y: test "khong chia thuong hai lan" KHONG bat duoc dot bien nay (da
// kiem: bo rieng chot do, test do van xanh) — vi so dong gop da bi xoa nen lan
// hai khong chia duoc gi. Rieng THONG BAO thi khong phu thuoc so dong gop,
// nen day la test DUY NHAT khoa cai chot ay.
func TestBossLonChetChiBaoDUNG_MOT_LAN(t *testing.T) {
	loa := MoBoPhatThanh()
	defer loa.Dung()

	k := khuVucTest(t)
	k.loa.Store(loa)
	datNguoi(t, k, 11, "Naruto")

	b := bossThuong(BossHoatDong, 1000, true)
	b.SatThuongTheoNguoi = map[int64]int64{11: 100}

	k.bossChet(b, 11, mocTest)
	k.bossChet(b, 11, mocTest) // lan hai — phai la viec KHONG

	tin := rutHangDoiLoa(loa)
	if len(tin) != 1 {
		t.Fatalf("goi bossChet hai lan sinh ra %d thong bao, mong doi DUNG 1: %v. "+
			"Bo trung phai lam o day (noi biet do la boss nao), KHONG o duong ong loa", len(tin), tin)
	}
}

// TestBossThuongKhongPhaiBossLonThiImLang: LaBossLon=false -> khong mot tin
// nao. Day la mac dinh an toan cua thiet ke: them nham mot boss lon se lam
// ngap kenh toan server.
func TestBossThuongKhongPhaiBossLonThiImLang(t *testing.T) {
	loa := MoBoPhatThanh()
	defer loa.Dung()

	k := khuVucTest(t)
	k.loa.Store(loa)
	datNguoi(t, k, 11, "Naruto")

	b := bossThuong(BossHoatDong, 1000, false)
	b.SatThuongTheoNguoi = map[int64]int64{11: 100}

	k.bossChet(b, 11, mocTest)

	if tin := rutHangDoiLoa(loa); len(tin) != 0 {
		t.Fatalf("boss THUONG chet ma van phat %d thong bao toan server: %v", len(tin), tin)
	}
}

// TestBossLonChetKhiNguoiHaDaRoiKhuVuc: nguoi giang don ket lieu roi mat ket
// noi ngay sau do -> khong con trong k.nguoiTrong. Thong bao van phai len song
// voi mot cai ten thay the, KHONG duoc nil-deref lam sap actor.
func TestBossLonChetKhiNguoiHaDaRoiKhuVuc(t *testing.T) {
	loa := MoBoPhatThanh()
	defer loa.Dung()

	k := khuVucTest(t)
	k.loa.Store(loa)

	b := bossThuong(BossHoatDong, 1000, true)
	k.bossChet(b, 12345, mocTest) // 12345 chua bao gio o trong khu vuc nay

	tin := rutHangDoiLoa(loa)
	if len(tin) != 1 {
		t.Fatalf("hang doi loa co %d tin, mong doi 1: %v", len(tin), tin)
	}
	if tin[0] != "Một ninja bí ẩn vừa hạ gục Boss Cửu Vĩ!" {
		t.Errorf("noi dung = %q, mong doi %q", tin[0], "Một ninja bí ẩn vừa hạ gục Boss Cửu Vĩ!")
	}
}

// TestThemBossLonHienXepThongBao (thiet ke §5, buoc HIEN): boss lon vao khu vuc
// o trang thai HIEN -> mot tin "vua xuat hien" di vao hang doi loa.
func TestThemBossLonHienXepThongBao(t *testing.T) {
	loa := MoBoPhatThanh()
	defer loa.Dung()

	k := khuVucTest(t)
	k.loa.Store(loa)

	k.themBoss(bossThuong(BossHien, 1000, true))

	tin := rutHangDoiLoa(loa)
	if len(tin) != 1 {
		t.Fatalf("hang doi loa co %d tin sau khi boss lon hien, mong doi 1: %v", len(tin), tin)
	}
	if tin[0] != "Boss Cửu Vĩ vừa xuất hiện!" {
		t.Errorf("noi dung = %q, mong doi %q", tin[0], "Boss Cửu Vĩ vừa xuất hiện!")
	}
}

// TestThemBossKhongOBuocHIENThiKhongBaoXuatHien: chi DUNG buoc HIEN moi la
// "vua xuat hien".
//
// Bon trang thai kia deu la truong hop that va deu phai im lang:
//   - NGHI: bo dieu phoi lich (Task 13) nap truoc mot con dang dem nguoc.
//   - CHAO/HOAT_DONG: mot con boss dang giua tran duoc chuyen/nap lai — no da
//     xuat hien tu truoc roi, bao lan nua la noi doi.
//   - TROI: do la CAI XAC. Bao "vua xuat hien" cho mot cai xac thi nguoi choi
//     se chay den mot con boss khong con danh duoc.
//
// Bay kem: dieu kien dung phai la `TrangThai == BossHien`, KHONG duoc noi long
// thanh `DangHien()` — DangHien() gom ca CHAO/HOAT_DONG/TROI (no tra loi cau
// hoi KHAC: "client co dang thay con nay khong").
func TestThemBossKhongOBuocHIENThiKhongBaoXuatHien(t *testing.T) {
	for _, tt := range []TrangThaiBoss{BossNghi, BossChao, BossHoatDong, BossTroi, BossRaDi} {
		loa := MoBoPhatThanh()
		k := khuVucTest(t)
		k.loa.Store(loa)

		k.themBoss(bossThuong(tt, 1000, true))

		if tin := rutHangDoiLoa(loa); len(tin) != 0 {
			t.Errorf("boss o trang thai %s ma da bao %d tin \"vua xuat hien\": %v", tt, len(tin), tin)
		}
		loa.Dung()
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 5 — ky luat actor (§6.9/§6.10b): duong chet boss chay TU TRONG actor
// ─────────────────────────────────────────────────────────────────────────────

// TestActorKhongTuKhoaKhiBossChet chay CA duong chet boss tu TRONG vongLap
// (bang suKienChay — diem chen chay tren dung goroutine actor) va do hai thu,
// y het TestActorKhongGoiGuiDieuKhienKhiXepThongBao (phattoanserver_test.go):
//
//  1. KHONG TREO trong 2 giay. Neu ai do nhet mot loi goi "gui roi cho tra
//     loi" (guiDieuKhien / RoiKhuVuc / AnhChupNguoiChoi / DanhSachQuaiSnapshot)
//     vao duong nay thi actor tu cho chinh minh -> het gio -> test DO ro rang.
//  2. KHONG CHAM HOP THU DIEU KHIEN: dem TU TRONG actor, ngay truoc va ngay
//     sau. Luc nay actor la nguoi duy nhat rut hop thu do ma no dang ban chay
//     ham cua ta, nen bat ky su kien dieu khien nao lo gui vao deu con nguyen.
func TestActorKhongTuKhoaKhiBossChet(t *testing.T) {
	loa := MoBoPhatThanh()
	defer loa.Dung()

	k := MoKhuVuc(TimBanDo(1))
	k.loa.Store(loa)
	k.BatDau()
	defer k.Dung()

	type doDuoc struct{ truoc, sau int }
	ketQua := make(chan doDuoc, 1)

	b := bossThuong(BossHoatDong, 1000, true)
	b.SatThuongTheoNguoi = map[int64]int64{11: 100}

	k.Gui(suKienChay{f: func() {
		truoc := len(k.hopThuDieuKhien)
		k.bossChet(b, 11, k.dongHo())
		ketQua <- doDuoc{truoc: truoc, sau: len(k.hopThuDieuKhien)}
	}})

	select {
	case d := <-ketQua:
		if d.sau != d.truoc {
			t.Fatalf("bossChet lam hop thu DIEU KHIEN tang tu %d len %d — actor KHONG duoc "+
				"goi guiDieuKhien tu trong chinh no (cho toi 100ms tren hop thu ma no dang ban "+
				"= tu khoa chac chan)", d.truoc, d.sau)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("bossChet goi TU TRONG actor bi TREO qua 2 giay — dau hieu tu khoa")
	}

	if tin := rutHangDoiLoa(loa); len(tin) != 1 {
		t.Fatalf("hang doi loa co %d tin, mong doi 1 (thong bao boss lon bi ha)", len(tin))
	}
}
