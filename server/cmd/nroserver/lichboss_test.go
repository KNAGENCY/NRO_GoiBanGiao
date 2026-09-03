package main

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"nroserver/internal/game"
)

// ─────────────────────────────────────────────────────────────────────────────
// KH6b Task 14 (phan hai) — NOI DAY bo dieu phoi lich boss trong main.go.
//
// Truoc phan hai nay, khong noi nao trong ca du an goi game.MoBoLichBoss/
// BatDau() — bo dieu phoi da qua review tu Task 13 nam do KHONG AI GOI TOI, y
// het lo hong cua bo phat thanh truoc Task 14 phan mot (xem phatthanh_test.go).
// Cac test duoi day khoa DUNG duong noi do tu phia package main.
// ─────────────────────────────────────────────────────────────────────────────

// TestNoiLichBossBatGoroutineNen khoa CA HAI nua cua duong noi: noiLichBoss()
// phai (a) tra ve mot bo dieu phoi khac nil, VA (b) bo dieu phoi do phai DA
// BatDau() — chua bat thi GiayNghi trong kho noi dung (600s/1800s, chot san o
// noidungboss.go) qua lon de mot test doi that toi luc sinh boss trong vai
// giay, nen khong the khoa bang cach doi boss xuat hien.
//
// Cach do (b) MA KHONG DUNG DONG HO THAT: dem so goroutine truoc/sau
// noiLichBoss() giong het TestBatDauLichBossHaiLanChiChayMotGoroutine
// (internal/game/lichboss_test.go) — BatDau() luon no ra DUNG MOT goroutine
// nen (`go l.vongLap()`, xem lichboss.go), nen bo BatDau() se khien so
// goroutine KHONG doi ma khong can cho dong ho tien mot giay nao.
func TestNoiLichBossBatGoroutineNen(t *testing.T) {
	soKhuVuc := game.MoSoKhuVuc()
	defer soKhuVuc.DungTatCa()

	runtime.Gosched()
	time.Sleep(20 * time.Millisecond) // de cac goroutine con lai (neu co) on dinh truoc khi dem
	truoc := runtime.NumGoroutine()

	lichBoss := noiLichBoss(soKhuVuc)
	defer lichBoss.Dung()

	if lichBoss == nil {
		t.Fatal("noiLichBoss tra ve nil — main() khong co gi de tat luc dong may chu, " +
			"va se khong con boss nao tung duoc sinh")
	}

	time.Sleep(50 * time.Millisecond) // cho goroutine vongLap (neu co) vao han doi ticker
	sau := runtime.NumGoroutine()

	if sau <= truoc {
		t.Fatalf("so goroutine truoc=%d sau=%d KHONG tang sau noiLichBoss — bo dieu phoi "+
			"chua tung BatDau(), se khong con boss nao tung duoc sinh trong suot doi may chu",
			truoc, sau)
	}
}

// TestNoiLichBossDungHaiLanKhongPanic la luoi an toan cho bay so 3 cua brief:
// Dung() goi hai lan (hoac goi khi chua BatDau) khong duoc panic — dung ly do
// voi TestGanLoaHaiLanKhongPanicLucTatMayChu (phatthanh_test.go).
func TestNoiLichBossDungHaiLanKhongPanic(t *testing.T) {
	soKhuVuc := game.MoSoKhuVuc()
	defer soKhuVuc.DungTatCa()

	lichBoss := noiLichBoss(soKhuVuc)

	xong := make(chan struct{})
	go func() {
		defer close(xong)
		lichBoss.Dung()
		lichBoss.Dung() // lan thu hai PHAI an toan (sync.Once trong BoLichBoss.Dung)
	}()
	select {
	case <-xong:
	case <-time.After(3 * time.Second):
		t.Fatal("goi Dung() hai lan bi treo hoac panic")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Phan 2: THU TU TAT — lichBoss.Dung() nam o dau trong donDepKhiTat.
//
// Thu tu da chot theo khung chu thich trong main.go: boLuu.Dung() ->
// lichBoss.Dung() -> loa.Dung() -> TatMayChu -> kho.Close(). Dung CUNG ky
// thuat "loa gia ghi lai trang thai the gioi" nhu TestDonDepKhiTatDungLoaTruoc
// KhiLuuVaDungKhuVuc (phatthanh_test.go), nhung lan nay ca HAI fake (lich boss
// gia + loa gia) cung ghi lai THU TU chung duoc goi — day moi la thu can khoa:
// khong phai "lichBoss.Dung() co chay" (test phan 1 da khoa roi) ma la "no
// chay TRUOC loa.Dung()".
// ─────────────────────────────────────────────────────────────────────────────

// dienBienDungGia dem thu tu cac lan Dung() duoc goi trong MOT lan donDepKhiTat.
type dienBienDungGia struct {
	mu       sync.Mutex
	thuTuGoi []string
}

func (d *dienBienDungGia) ghi(ten string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.thuTuGoi = append(d.thuTuGoi, ten)
}

func (d *dienBienDungGia) thuTu() []string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return append([]string(nil), d.thuTuGoi...)
}

// lichBossGiamSat dong vai *game.BoLichBoss trong day don dep — chi ghi lai
// NO duoc Dung() vao luc nao so voi loa, khong lam gi khac.
type lichBossGiamSat struct {
	dienBien *dienBienDungGia
	soLan    atomic.Int64
}

func (l *lichBossGiamSat) Dung() {
	l.soLan.Add(1)
	l.dienBien.ghi("lichBoss")
}

// loaThuTuGia dong vai *game.BoPhatThanh, chi ghi thu tu (khac loaGiamSat cua
// phatthanh_test.go — file do kiem trang thai the gioi, file nay kiem THU TU
// giua hai manh moi cua Task 14).
type loaThuTuGia struct {
	dienBien *dienBienDungGia
}

func (l *loaThuTuGia) Dung() {
	l.dienBien.ghi("loa")
}

func TestDonDepKhiTatDungLichBossTruocLoa(t *testing.T) {
	soKhuVuc := game.MoSoKhuVuc()
	defer soKhuVuc.DungTatCa()
	kl := &khoLuuGiaMain{}
	thoGhi := game.MoThoGhi(kl, 2)
	boLuu := game.MoBoLuuDinhKy(soKhuVuc, thoGhi)
	boLuu.BatDau()
	kd := &khoDongGia{}

	dienBien := &dienBienDungGia{}
	lichBoss := &lichBossGiamSat{dienBien: dienBien}
	loa := &loaThuTuGia{dienBien: dienBien}

	ctx, huy := context.WithTimeout(context.Background(), 10*time.Second)
	defer huy()

	if err := donDepKhiTat(ctx, kd, soKhuVuc, thoGhi, boLuu, nil, lichBoss, loa, &boGiaDemDung{}, &boGiaDemDung{}); err != nil {
		t.Fatalf("donDepKhiTat loi trong dieu kien binh thuong: %v", err)
	}

	if n := lichBoss.soLan.Load(); n != 1 {
		t.Fatalf("lichBoss.Dung() duoc goi %d lan, muon 1 — 0 nghia la goroutine lich boss "+
			"RO RI qua ca doi tien trinh", n)
	}

	thuTu := dienBien.thuTu()
	if len(thuTu) != 2 || thuTu[0] != "lichBoss" || thuTu[1] != "loa" {
		t.Fatalf("thu tu Dung() = %v, muon [lichBoss loa] — lichBoss la NGUON sinh thong bao, "+
			"loa la duong ONG cho no chay qua: tat nguon truoc, ong sau (xem khung chu thich "+
			"CHO NOI BO DIEU PHOI LICH BOSS trong main.go)", thuTu)
	}
}
