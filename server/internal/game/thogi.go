package game

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"nroserver/internal/store"
)

// ThoGhi (dac ta 2026-07-27-luu-trang-thai-song.md §3/§6/§9) la CAI CHAN de
// khong lap lai loi actor bi chan cho o dia (§6.10b BAT-DAU-TU-DAY.md: actor
// khu vuc tung bi chan toi 30 giay vi tu gui goi mang cho client ngung doc).
//
// LUAT SO MOT: goroutine actor CUA MOI khu vuc (KhuVuc.vongLap, khuvuc.go)
// tuyet doi khong duoc dung lai cho ghi CSDL. Actor chi "chup anh" trang
// thai (store.NhanVat, BAN SAO GIA TRI — xem suKienXinAnhChup/AnhChupNguoiChoi
// o khuvuc.go, Ke hoach 11 Task 2) roi GOI Xep() — ham nay CHI cham hang doi
// trong bo nho (mot map duoc khoa bang mutex), KHONG BAO GIO goi toi
// KhoLuu.CapNhatNhanVat truc tiep. Viec cham dia that su xay ra trong MOT
// goroutine rieng (vongLap ben duoi) — hoan toan tach khoi luong goi Xep().
//
//	Actor khu vuc            Xep() — CHI cham map     vongLap (goroutine rieng)
//	─────────────            ────────────────────     ─────────────────────────
//	chup anh (ban sao)  ──►  hangDoi[id] = nv    ──►  lay tu hangDoi, GOI
//	di tiep NGAY              (khoa mutex ~vai ns)     KhoLuu.CapNhatNhanVat
//	KHONG BAO GIO CHO                                  (cho o day, khong sao —
//	                                                     khong actor nao dung)
type ThoGhi struct {
	kho      KhoLuu
	soThoGhi int // so goroutine ghi CHAY SONG SONG toi da trong MOI lan xa het hang doi (ghiHetHangDoi)

	mu      sync.Mutex
	hangDoi map[int64]store.NhanVat // hang doi GOP: id nhan vat -> BAN MOI NHAT (khong phai danh sach/slice — xem chu thich duoi)

	coViec     chan struct{} // tin hieu "co viec moi" — buffer 1, gui KHONG CHAN (select+default)
	dung       chan struct{}
	xong       chan struct{}
	dungMotLan sync.Once

	// soLoiGhi dem TONG so lan CapNhatNhanVat tra loi (dac ta §6.3: "dem, log
	// muc canh bao, va thu lai o luot sau" — KHONG duoc nuot im lang nhu
	// PlayerDAO.updatePlayer ban Java). atomic.Int64 vi nhieu goroutine ghi
	// (soThoGhi > 1) co the cung tang gia tri nay dong thoi.
	soLoiGhi atomic.Int64

	// soLoiLienTiep dem so lan loi LIEN TIEP GAN NHAT (reset ve 0 moi khi co
	// MOT lan ghi thanh cong, bat ke id nao) — dung de quyet dinh log muc
	// Warn (loi thoang qua) hay muc Error (dac ta §6.3: "hong lien tiep nhieu
	// lan thi log muc loi de nguoi van hanh thay"). Day la mot chi bao THO
	// (toan cuc, khong tach rieng tung id) — du de bao "CSDL co ve dang co
	// van de", khong can chinh xac tuyet doi toi tung nhan vat.
	soLoiLienTiep atomic.Int64
}

// KhoLuu la phan CSDL toi thieu ma ThoGhi can — tach thanh interface HEP (dung
// khuon voi Kho o login.go:64-69) de kiem thu tiem duoc kho gia, khong can
// CSDL that. store.Store (CapNhatNhanVat, store.go) thoa man interface nay
// tu nhien, khong can glue code them.
type KhoLuu interface {
	CapNhatNhanVat(context.Context, store.NhanVat) error
}

const (
	// hanGioGhi la han gio cho MOI lan goi KhoLuu.CapNhatNhanVat (dac ta
	// §6.2: "de xuat 5 giay" — day la CHO DAU TIEN trong ca du an dung
	// context.WithTimeout, chua co tien le de bat chuoc). Khong co han gio
	// nay thi mot lan CSDL treo se giu goroutine ghi (va ca ban ghi do) vo
	// thoi han — voi soThoGhi goroutine gioi han, du vai lan treo cung du
	// lam nghen toan bo duong ghi.
	hanGioGhi = 5 * time.Second

	// khoangCachThuLai la nhip "kiem lai hang doi" NEN cua vongLap — dam bao
	// mot ban ghi loi (con nam lai trong hangDoi, xem ghiMotBan) duoc THU LAI
	// (dac ta §6.3) MA KHONG can cho ai goi Xep() them lan nao nua. Nho: khi
	// hangDoi RONG, moi nhip chi la mot lan khoa mutex + kiem len()==0 roi
	// tra ve ngay — chi phi khong dang ke (cung tinh than voi ticker 100ms
	// san co cua KhuVuc.vongLap, khuvuc.go).
	khoangCachThuLai = 50 * time.Millisecond

	// nguongLoiLienTiepCanhBaoNang la so loi LIEN TIEP (soLoiLienTiep) tro len
	// thi nang muc log tu Warn sang Error (dac ta §6.3).
	nguongLoiLienTiepCanhBaoNang = 3
)

// MoThoGhi khoi dong MOT goroutine "tho ghi" chay nen, xu ly hang doi gop
// bang toi da soThoGhi lan ghi CSDL song song trong moi dot xa hang doi.
// soThoGhi < 1 duoc ep ve 1 (phong thu goi nham, khong de ThoGhi hoa vo dung —
// vd giay phep rong (make(chan struct{}, 0)) se khoa vinh vien lan ghi dau).
func MoThoGhi(kho KhoLuu, soThoGhi int) *ThoGhi {
	if soThoGhi < 1 {
		soThoGhi = 1
	}
	t := &ThoGhi{
		kho:      kho,
		soThoGhi: soThoGhi,
		hangDoi:  make(map[int64]store.NhanVat),
		coViec:   make(chan struct{}, 1),
		dung:     make(chan struct{}),
		xong:     make(chan struct{}),
	}
	go t.vongLap()
	return t
}

// Xep dua mot ban chup store.NhanVat vao hang doi cho ghi — KHONG BAO GIO
// CHAN, du kho dang ghi cham hay hang doi dang co bao nhieu id. Day la LUAT
// SO MOT cua ca task (xem chu thich ThoGhi o tren): ham nay CHI khoa mutex
// noi bo trong vai chuc nano-giay de ghi vao map, TUYET DOI khong goi toi
// KhoLuu — nen thoi gian chay khong phu thuoc gi vao toc do dia/CSDL.
//
// hangDoi la MAP id -> ban moi nhat (khong phai slice/danh sach) — day la
// quyet dinh da chot cua chu du an (dac ta §9 muc 2): xep nhieu lan cho CUNG
// mot id chi GHI DE, khong CONG DON — vi ban moi luon tot hon ban cu nen giu
// ban cu la vo nghia, VA vi vay hang doi khong bao gio phinh to theo SO LAN
// goi Xep(), chi phinh theo SO NGUOI CHOI dang song (bi chan tren tu nhien).
func (t *ThoGhi) Xep(nv store.NhanVat) {
	t.mu.Lock()
	t.hangDoi[nv.ID] = nv
	t.mu.Unlock()

	select {
	case t.coViec <- struct{}{}:
	default: // da co tin hieu cho xu ly roi (hoac tho dang ban) — khong can gui them, KHONG CHAN
	}
}

// DongVaChoXong bao tho ghi dung nhan viec moi, GHI HET nhung gi con trong
// hang doi, roi moi tra ve — dung cho luc tat may chu (dac ta §4 diem C):
// "xin anh chup toan bo -> cho tho ghi lam xong -> ROI MOI DungTatCa()".
//
// ctx dat CAN nay cho ca qua trinh cho — neu tho ghi khong xong kip han gio
// (vd CSDL treo that su, vuot han hanGioGhi nhieu lan lien tiep), ham tra ve
// ctx.Err() thay vi treo may chu vo han luc tat: mot lan tat "danh hong" con
// hon treo ca tien trinh khong bao gio thoat duoc.
func (t *ThoGhi) DongVaChoXong(ctx context.Context) error {
	t.dungMotLan.Do(func() {
		close(t.dung)
	})
	select {
	case <-t.xong:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// vongLap la than tho ghi: NHAN TIN HIEU "co viec" (Xep vua goi), HOAC nhip
// thu lai dinh ky (khoangCachThuLai — bat lai nhung ban ghi loi lan truoc),
// HOAC lenh dung. Day la NOI DUY NHAT goi ghiHetHangDoi() — dung cau truc
// voi KhuVuc.vongLap (khuvuc.go) de nhat quan phong cach actor toan du an.
func (t *ThoGhi) vongLap() {
	ticker := time.NewTicker(khoangCachThuLai)
	defer ticker.Stop()
	defer close(t.xong)

	for {
		select {
		case <-t.coViec:
			t.ghiHetHangDoi()
		case <-ticker.C:
			t.ghiHetHangDoi()
		case <-t.dung:
			// Lenh dung: XA HET hang doi con lai MOT LAN CUOI (dac ta §4
			// diem C — "ghi het nhung gi con trong hang doi") roi thoat.
			// Day la BEST-EFFORT mot dot (khong lap vo han neu cu loi mai) —
			// ctx cua nguoi goi DongVaChoXong la cai chan thoi gian cho o
			// tang tren, khong phai o day.
			t.ghiHetHangDoi()
			return
		}
	}
}

// ghiHetHangDoi lay TOAN BO noi dung hangDoi hien tai (khoa mutex mot lan de
// "doat" con tro map roi thay bang map RONG moi — nen Xep() goi trong luc
// dang ghi se roi vao map MOI, khong dung cham voi ban dang xu ly), roi ghi
// tung ban bang toi da soThoGhi goroutine song song.
func (t *ThoGhi) ghiHetHangDoi() {
	t.mu.Lock()
	if len(t.hangDoi) == 0 {
		t.mu.Unlock()
		return
	}
	banSao := t.hangDoi
	t.hangDoi = make(map[int64]store.NhanVat)
	t.mu.Unlock()

	var cho sync.WaitGroup
	giayPhep := make(chan struct{}, t.soThoGhi)
	for _, nv := range banSao {
		cho.Add(1)
		giayPhep <- struct{}{}
		go func(nv store.NhanVat) {
			defer cho.Done()
			defer func() { <-giayPhep }()
			t.ghiMotBan(nv)
		}(nv)
	}
	cho.Wait()
}

// ghiMotBan ghi MOT store.NhanVat xuong CSDL, boc context.WithTimeout(hanGioGhi)
// (dac ta §6.2 — CHO DAU TIEN dung timeout trong du an). Loi KHONG bi nuot
// (dac ta §6.3, bai hoc PlayerDAO.updatePlayer ban Java): dem soLoiGhi, log
// canh bao (hoac loi neu hong lien tiep nhieu lan), roi XEP LAI ban do vao
// hangDoi de THU LAI o nhip sau — DU LIEU VAN CON TRONG BO NHO nen khong mat.
//
// Chi xep lai NEU hangDoi CHUA co ban moi hon cho id nay (!coBanMoiHon): neu
// trong luc dang ghi (co the mat toi hanGioGhi=5s) da co Xep() moi cho CUNG
// id nay, ban moi do dang nam trong hangDoi roi va CHAC CHAN moi hon ban vua
// that bai — ghi de len no bang ban loi (cu hon) se lam LUI du lieu, dung
// nguyen bay §7.1 cua dac ta ("ghi de ban moi bang ban cu").
func (t *ThoGhi) ghiMotBan(nv store.NhanVat) {
	ctx, huy := context.WithTimeout(context.Background(), hanGioGhi)
	defer huy()

	if err := t.kho.CapNhatNhanVat(ctx, nv); err != nil {
		t.soLoiGhi.Add(1)
		lienTiep := t.soLoiLienTiep.Add(1)
		if lienTiep >= nguongLoiLienTiepCanhBaoNang {
			slog.Error("ghi nhan vat xuong CSDL that bai NHIEU LAN LIEN TIEP — co the CSDL dang co van de",
				"idNhanVat", nv.ID, "loi", err, "soLanLienTiep", lienTiep)
		} else {
			slog.Warn("ghi nhan vat xuong CSDL that bai, se tu dong thu lai o nhip sau",
				"idNhanVat", nv.ID, "loi", err)
		}

		t.mu.Lock()
		if _, coBanMoiHon := t.hangDoi[nv.ID]; !coBanMoiHon {
			t.hangDoi[nv.ID] = nv
		}
		t.mu.Unlock()
		return
	}

	t.soLoiLienTiep.Store(0)
}
