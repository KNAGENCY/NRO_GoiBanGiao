package main

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"nroserver/internal/game"
	"nroserver/internal/store"
)

// ─────────────────────────────────────────────────────────────────────────────
// Phan 1: NOI DAY luc khoi dong — lo hong tang tai KH6b Task 14
//
// Truoc Task 14, main.go KHONG BAO GIO goi GanLoa. He qua khong he co mot dong
// log nao: KhuVuc.xepThongBao (khuvuc.go) doc k.loa.Load(), thay nil thi LANG
// LE `return`. Nghia la moi thong bao boss (goi 93 — chay chu toan may chu)
// roi vao hu khong, va khong the phat hien bang cach doc log. Hai test duoi
// day khoa DUNG cai lo hong do tu phia package main.
// ─────────────────────────────────────────────────────────────────────────────

// TestNoiBoPhatThanhGanLoaVaoMoiKhuVuc khoa NUA DAU cua duong noi: cai loa ma
// main() giu trong tay PHAI la dung cai loa da duoc cam vao CA 5 khu vuc.
//
// Bo loi goi soKhuVuc.GanLoa(loa) trong noiBoPhatThanh -> LoaDaGan() tra nil
// -> test nay DO. Do chinh la lo hong dang ton tai truoc Task 14.
func TestNoiBoPhatThanhGanLoaVaoMoiKhuVuc(t *testing.T) {
	soKhuVuc := game.MoSoKhuVuc()
	loa := noiBoPhatThanh(soKhuVuc)
	defer func() {
		// Dung thu tu tat that: loa TRUOC, khu vuc SAU (xem donDepKhiTat).
		loa.Dung()
		soKhuVuc.DungTatCa()
	}()

	if loa == nil {
		t.Fatal("noiBoPhatThanh tra ve nil — main() khong co gi de tat luc dong may chu")
	}
	daGan := soKhuVuc.LoaDaGan()
	if daGan == nil {
		t.Fatal("sau noiBoPhatThanh van co khu vuc CHUA duoc gan loa — moi thong bao " +
			"boss se roi vao hu khong (xepThongBao thay loa nil thi lang le return)")
	}
	if daGan != loa {
		t.Fatal("loa gan vao khu vuc KHAC cai loa tra ve cho main() — main() se tat " +
			"nham cai loa, cai dang chay thanh goroutine mo coi")
	}
}

// TestNoiBoPhatThanhBatGoroutinePhat khoa NUA SAU cua duong noi: gan loa vao
// khu vuc thoi CHUA DU — phai co goroutine phat DANG CHAY de rut hang doi.
// Gan ma quen bat thi thong bao van "mat", chi khac la mat trong hang doi thay
// vi mat o cho k.loa == nil.
//
// Cach do KHONG dung dong ho: do day hang doi den TRAN roi xem co bao gio he
// ra mot cho trong khong. Loa dang chay se rut cai dau tien gan nhu tuc thi
// (lanPhatCuoi khoi tao la zero value — thongbao.go co y de thong bao dau tien
// len song NGAY). Loa khong chay thi hang doi day VINH VIEN, moi Xep deu bi BO.
func TestNoiBoPhatThanhBatGoroutinePhat(t *testing.T) {
	soKhuVuc := game.MoSoKhuVuc()
	loa := noiBoPhatThanh(soKhuVuc)
	defer func() {
		loa.Dung()
		soKhuVuc.DungTatCa()
	}()

	// Ghim CUNG 64 thay vi doc lai hang so tranHangDoiThongBao (no unexported
	// o package game, va cung cung ly do da ghi trong thongbao_test.go: doc lai
	// hang so lam test dong tinh voi code neu ai do noi tran ra).
	const sucChua = 64
	for i := 0; i < sucChua; i++ {
		loa.Xep(game.YeuCauThongBao{NoiDung: "do cho day hang doi", MocTao: time.Now()})
	}

	lot := false
	for han := time.Now().Add(3 * time.Second); time.Now().Before(han); {
		truoc := loa.SoLanBo()
		loa.Xep(game.YeuCauThongBao{NoiDung: "kiem xem co cho trong chua", MocTao: time.Now()})
		if loa.SoLanBo() == truoc {
			lot = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if !lot {
		t.Fatalf("sau 3 giay hang doi loa van DAY nguyen (SoLanBo=%d) — khong ai rut no ra, "+
			"tuc goroutine phat CHUA duoc bat: thong bao se nam trong hang doi den chet",
			loa.SoLanBo())
	}
}

// TestGanLoaHaiLanKhongPanicLucTatMayChu la LUOI AN TOAN cho bay so 1 cua
// Task 14: BatDau() goi hai lan sinh HAI goroutine, ca hai cung `defer
// close(b.xong)` -> lan thu hai panic "close of closed channel" tren mot
// goroutine nen, khong ai bat duoc, CA TIEN TRINH MAY CHU CHET dung luc dang
// tat. Lop chan (CompareAndSwap) da vao thongbao.go ngay 28/08; test nay dam
// bao duong noi cua main.go khong lach qua duoc no.
func TestGanLoaHaiLanKhongPanicLucTatMayChu(t *testing.T) {
	soKhuVuc := game.MoSoKhuVuc()
	defer soKhuVuc.DungTatCa()

	loa := game.MoBoPhatThanh()
	soKhuVuc.GanLoa(loa)
	soKhuVuc.GanLoa(loa) // lan thu hai PHAI bi bo qua, khong sinh goroutine thu hai

	xong := make(chan struct{})
	go func() {
		defer close(xong)
		loa.Dung()
	}()
	select {
	case <-xong:
	case <-time.After(3 * time.Second):
		t.Fatal("loa.Dung() treo sau khi GanLoa duoc goi hai lan — co goroutine phat thu hai " +
			"con song (lop chan CompareAndSwap trong BatDau da bi lach qua)")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 2: THU TU TAT — loa.Dung() nam o dau trong day don dep
//
// Thu tu da chon: boLuu.Dung() -> loa.Dung() -> TatMayChu(LuuTatCa ->
// DongVaChoXong -> DungTatCa) -> kho.Close().
//
// Cach do KHONG dung dong ho: cai loa gia duoi day, ngay TRONG Dung() cua no,
// chup lai TRANG THAI THE GIOI tai dung khoanh khac ay — khu vuc con song
// khong, tho ghi con nhan viec khong. Hai cau hoi do la CHINH XAC hai rang
// buoc ma thu tu nay ton tai de bao ve.
// ─────────────────────────────────────────────────────────────────────────────

// idBanGhiMoc la mot id nhan vat khong trung voi test nao khac trong package
// nay — no la "moc do" ma loa gia tha vao tho ghi ngay luc bi dung.
const idBanGhiMoc int64 = 424242

// loaGiamSat dong vai *game.BoPhatThanh trong day don dep. No KHONG phat gi
// ca; viec duy nhat cua no la ghi lai the gioi trong ra sao tai dung giay phut
// donDepKhiTat goi Dung().
type loaGiamSat struct {
	soKhuVuc *game.SoKhuVuc
	thoGhi   *game.ThoGhi

	soLanDung     atomic.Int64
	khuVucConSong atomic.Bool
	loiHoiKhuVuc  atomic.Value // string — loi that khi hoi khu vuc, de bao cho ro
}

func (l *loaGiamSat) Dung() {
	l.soLanDung.Add(1)

	// Cau hoi 1: khu vuc con song khong? Neu loa bi dung SAU DungTatCa thi
	// AnhChupNguoiChoiCoLoi tra ErrActorDaDung. Do la kich ban phai chan:
	// PhatToanServer bom goi vao cac khu vuc dang bi dong — phat cho ai doc?
	_, err := l.soKhuVuc.Tra(1).AnhChupNguoiChoiCoLoi(0)
	l.khuVucConSong.Store(err == nil)
	if err != nil {
		l.loiHoiKhuVuc.Store(err.Error())
	}

	// Cau hoi 2: tho ghi con nhan viec khong? Tha mot ban ghi MOC vao day. Neu
	// loa bi dung TRUOC DongVaChoXong (dung thu tu), ban moc se duoc ghi xuong
	// "dia" truoc khi donDepKhiTat tra ve. Neu bi dung SAU, tho ghi da dong,
	// Xep() khong panic nhung ban ghi nam lai trong hang doi VINH VIEN — day
	// dung la kieu mat du lieu am tham ma ca trinh tu tat may chu duoc dung
	// len de chan.
	l.thoGhi.Xep(store.NhanVat{ID: idBanGhiMoc, HpHienTai: 42, HpToiDa: 99})
}

// khoLuuDemBan la KhoLuu gia chi lam mot viec: nho lai nhung id da THUC SU
// duoc ghi xuong.
type khoLuuDemBan struct {
	mu    sync.Mutex
	daGhi map[int64]bool
}

func (k *khoLuuDemBan) CapNhatNhanVat(_ context.Context, nv store.NhanVat) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.daGhi == nil {
		k.daGhi = make(map[int64]bool)
	}
	k.daGhi[nv.ID] = true
	return nil
}

func (k *khoLuuDemBan) coBan(id int64) bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	return k.daGhi[id]
}

// TestDonDepKhiTatDungLoaTruocKhiLuuVaDungKhuVuc la test THU TU cua Task 14.
//
// Doi loa.Dung() xuong SAU game.TatMayChu(...) -> ca hai khang dinh duoi day
// cung do: khu vuc da dung, va ban ghi moc khong bao gio xuong dia.
func TestDonDepKhiTatDungLoaTruocKhiLuuVaDungKhuVuc(t *testing.T) {
	soKhuVuc := game.MoSoKhuVuc()
	kl := &khoLuuDemBan{}
	thoGhi := game.MoThoGhi(kl, 2)
	boLuu := game.MoBoLuuDinhKy(soKhuVuc, thoGhi)
	boLuu.BatDau()
	kd := &khoDongGia{}
	loa := &loaGiamSat{soKhuVuc: soKhuVuc, thoGhi: thoGhi}

	ctx, huy := context.WithTimeout(context.Background(), 10*time.Second)
	defer huy()

	if err := donDepKhiTat(ctx, kd, soKhuVuc, thoGhi, boLuu, game.MoBoLichBoss(soKhuVuc, nil), loa, &boGiaDemDung{}, &boGiaDemDung{}); err != nil {
		t.Fatalf("donDepKhiTat loi trong dieu kien binh thuong: %v", err)
	}

	if n := loa.soLanDung.Load(); n != 1 {
		t.Fatalf("loa.Dung() duoc goi %d lan trong donDepKhiTat, muon dung 1 — "+
			"0 nghia la goroutine phat thanh RO RI qua ca doi tien trinh", n)
	}
	if !loa.khuVucConSong.Load() {
		loi, _ := loa.loiHoiKhuVuc.Load().(string)
		t.Fatalf("luc loa.Dung() chay thi khu vuc DA DUNG roi (%s) — loa.Dung() dat QUA MUON: "+
			"tu day toi luc do, PhatToanServer con bom goi vao cac khu vuc dang bi dong", loi)
	}
	if !kl.coBan(idBanGhiMoc) {
		t.Fatal("ban ghi tha vao tho ghi ngay luc loa.Dung() chay KHONG bao gio xuong dia — " +
			"loa.Dung() dat SAU DongVaChoXong: bat cu gi sinh ra tu day tro di deu mat am tham")
	}
	if !kd.laDaDong() {
		t.Fatal("donDepKhiTat khong con goi kho.Close() sau khi them buoc tat loa")
	}
}

// TestDonDepKhiTatVoiLoaThatKhongTreoKhongPanic phu hai kich ban song sinh cua
// bay so 2: (a) tat may chu khi loa CHUA duoc BatDau bao gio, (b) tat may chu
// HAI LAN. Ca hai deu phai tra ve, khong treo, khong panic.
//
// Dung loa THAT (khong phai loaGiamSat) vi day dung la cho hai duong do gap
// nhau: BoPhatThanh.Dung() cho <-b.xong, ma khong ai dong kenh do neu chua ai
// BatDau — co daBatDau chinh la de chan cai treo do.
func TestDonDepKhiTatVoiLoaThatKhongTreoKhongPanic(t *testing.T) {
	t.Run("loa chua bao gio BatDau", func(t *testing.T) {
		soKhuVuc := game.MoSoKhuVuc()
		thoGhi := game.MoThoGhi(&khoLuuDemBan{}, 1)
		boLuu := game.MoBoLuuDinhKy(soKhuVuc, thoGhi)
		boLuu.BatDau()
		loa := game.MoBoPhatThanh() // CO Y khong GanLoa/BatDau

		xong := make(chan error, 1)
		go func() {
			ctx, huy := context.WithTimeout(context.Background(), 10*time.Second)
			defer huy()
			xong <- donDepKhiTat(ctx, &khoDongGia{}, soKhuVuc, thoGhi, boLuu, game.MoBoLichBoss(soKhuVuc, nil), loa, &boGiaDemDung{}, &boGiaDemDung{})
		}()
		select {
		case err := <-xong:
			if err != nil {
				t.Fatalf("donDepKhiTat loi khi loa chua tung BatDau: %v", err)
			}
		case <-time.After(15 * time.Second):
			t.Fatal("donDepKhiTat TREO khi loa chua tung BatDau — Dung() dang cho mot " +
				"goroutine khong bao gio ton tai")
		}
	})

	t.Run("tat may chu hai lan", func(t *testing.T) {
		soKhuVuc := game.MoSoKhuVuc()
		thoGhi := game.MoThoGhi(&khoLuuDemBan{}, 1)
		boLuu := game.MoBoLuuDinhKy(soKhuVuc, thoGhi)
		boLuu.BatDau()
		loa := noiBoPhatThanh(soKhuVuc)
		kd := &khoDongGia{}

		ctx, huy := context.WithTimeout(context.Background(), 10*time.Second)
		defer huy()
		if err := donDepKhiTat(ctx, kd, soKhuVuc, thoGhi, boLuu, game.MoBoLichBoss(soKhuVuc, nil), loa, &boGiaDemDung{}, &boGiaDemDung{}); err != nil {
			t.Fatalf("lan tat thu nhat loi: %v", err)
		}

		xong := make(chan error, 1)
		go func() {
			ctx2, huy2 := context.WithTimeout(context.Background(), 10*time.Second)
			defer huy2()
			xong <- donDepKhiTat(ctx2, kd, soKhuVuc, thoGhi, boLuu, game.MoBoLichBoss(soKhuVuc, nil), loa, &boGiaDemDung{}, &boGiaDemDung{})
		}()
		select {
		case <-xong: // loi hay khong deu duoc — mien la KHONG panic va KHONG treo
		case <-time.After(15 * time.Second):
			t.Fatal("lan tat thu hai TREO — trinh tu tat may chu khong chiu duoc goi lai")
		}
	})
}
