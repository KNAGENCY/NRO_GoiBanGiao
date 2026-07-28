package game

import (
	"sync"
	"testing"
	"time"

	"nroserver/internal/store"
)

// thoGhiGia la "tho ghi gia" thoa man boXep (luudinhky.go) — CHI ghi lai MOI
// lan Xep duoc goi vao mot slice, KHONG cham CSDL/kho gia nao. Khac voi
// khoLuuGia (thogi_test.go, gia MOT TANG THAP HON — duoi ThoGhi that): o day
// gia THANG ca ThoGhi, nen test quan sat duoc NGAY LAP TUC (dong bo, trong
// cung goroutine goi kiemVaLuu/LuuTatCa — AnhChupNguoiChoi da tu no cho actor
// tra loi xong truoc khi ham do return, xem khuvuc.go), khong phai cho
// goroutine nen cua ThoGhi that flush qua polling nhu thogi_test.go phai lam.
type thoGhiGia struct {
	mu  sync.Mutex
	xep []store.NhanVat
}

func (g *thoGhiGia) Xep(nv store.NhanVat) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.xep = append(g.xep, nv)
}

func (g *thoGhiGia) banSaoXep() []store.NhanVat {
	g.mu.Lock()
	defer g.mu.Unlock()
	ban := make([]store.NhanVat, len(g.xep))
	copy(ban, g.xep)
	return ban
}

func (g *thoGhiGia) coId(id int64) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, nv := range g.xep {
		if nv.ID == id {
			return true
		}
	}
	return false
}

// soKhuVucLuuDinhKyGia dung 2 khu vuc THAT (actor that, BatDau that — CAN
// THIET vi LuuTatCa/AnhChupNguoiChoi di qua hop thu actor, khong hoat dong
// neu khong ai rut hop thu, dung ly do voi khuVucAnhChupGia o khuvuc_test.go)
// cho cac test trong file nay. Dung 2 khu vuc (khong phai 1) de LuuTatCa/
// kiemVaLuu thuc su phai DUYET NHIEU khu vuc, khong "tinh co dung" voi 1 khu
// vuc duy nhat. Vi tri hai nguoi dat THAT XA quai mac dinh cua map (dung
// tinh than khuVucAnhChupGia, khuvuc_test.go) de nhipTim that (100ms/lan,
// actor that) khong vo tinh doi HpHienTai giua luc test dang chay.
func soKhuVucLuuDinhKyGia(t *testing.T) (s *SoKhuVuc, nvA, nvB *store.NhanVat) {
	t.Helper()
	k1 := MoKhuVuc(TimBanDo(1))
	k1.BatDau()
	t.Cleanup(func() { _ = k1.Dung() })
	k2 := MoKhuVuc(TimBanDo(2))
	k2.BatDau()
	t.Cleanup(func() { _ = k2.Dung() })
	s = &SoKhuVuc{theoMap: map[int32]*KhuVuc{1: k1, 2: k2}}

	nvA = &store.NhanVat{ID: 101, Ten: "A", HpToiDa: 100, HpHienTai: 77, ViTriX: 30000, ViTriY: 30000}
	nvB = &store.NhanVat{ID: 202, Ten: "B", HpToiDa: 100, HpHienTai: 55, ViTriX: -30000, ViTriY: -30000}
	k1.Gui(suKienVao{nv: nvA, phien: &phienGiaKhuVuc{}})
	k2.Gui(suKienVao{nv: nvB, phien: &phienGiaKhuVuc{}})
	return s, nvA, nvB
}

// TestLuuTatCaDuyetMoiKhuVuc la test RIENG cho helper LuuTatCa (khong qua
// BoLuuDinhKy) — day la HELPER DUNG CHUNG Task 6 (tat may chu) se TAI DUNG
// (xem chu thich LuuTatCa, luudinhky.go), nen can mot test xac nhan hop dong
// cua rieng no: "duyet MOI khu vuc, Xep TUNG nguoi trong do", doc lap voi
// logic dong ho 60 giay cua BoLuuDinhKy.
func TestLuuTatCaDuyetMoiKhuVuc(t *testing.T) {
	s, nvA, nvB := soKhuVucLuuDinhKyGia(t)
	g := &thoGhiGia{}

	s.LuuTatCa(g)

	ban := g.banSaoXep()
	if len(ban) != 2 {
		t.Fatalf("LuuTatCa xep %d ban, mong doi 2 (1 nguoi/khu vuc x 2 khu vuc)", len(ban))
	}
	if !g.coId(nvA.ID) || !g.coId(nvB.ID) {
		t.Fatalf("LuuTatCa thieu nguoi: %+v (mong doi ca id %d va %d)", ban, nvA.ID, nvB.ID)
	}
}

// TestBoLuuDinhKyChuaDu60GiayChuaLuu la brief Task 4 buoc 1 y dau: CHUA toi
// 60 giay ke tu lanLuuCuoi -> kiemVaLuu KHONG duoc lam gi ca (khong Xep bat
// ky ai). Ghim CUNG 59 giay — KHONG doc lai nhipLuuDinhKy tu code (bay oracle
// mu, xem brief: "dự án đã dính 3 lần").
func TestBoLuuDinhKyChuaDu60GiayChuaLuu(t *testing.T) {
	s, _, _ := soKhuVucLuuDinhKyGia(t)
	g := &thoGhiGia{}
	b := MoBoLuuDinhKy(s, g)

	t0 := time.Now()
	b.dongHo = func() time.Time { return t0 }
	b.lanLuuCuoi = t0

	b.kiemVaLuu(t0.Add(59 * time.Second))

	if soLan := len(g.banSaoXep()); soLan != 0 {
		t.Fatalf("kiemVaLuu tai +59s da Xep %d ban, mong doi 0 (chua du 60 giay thi chua luu gi)", soLan)
	}
}

// TestBoLuuDinhKyDuHoacQua60GiayLuuTatCa la brief Task 4 buoc 1 y hai: DUNG
// hoac VUOT 60 giay -> kiemVaLuu PHAI goi LuuTatCa, moi nguoi trong MOI khu
// vuc duoc Xep. Ghim CUNG hai moc +60s (dung nguong, "đúng") va +90s (vuot
// nguong, "qua") — dung tinh than "Đúng/qua 60 giây" cua brief, KHONG doc lai
// hang so tu code. Day cung la moc kiem thu dot bien (Step 5 brief: doi 60s
// thanh 600s) se lam do o phan +60s.
func TestBoLuuDinhKyDuHoacQua60GiayLuuTatCa(t *testing.T) {
	for _, tc := range []struct {
		ten  string
		lech time.Duration
	}{
		{"dung moc, +60s", 60 * time.Second},
		{"vuot moc, +90s", 90 * time.Second},
	} {
		t.Run(tc.ten, func(t *testing.T) {
			s, nvA, nvB := soKhuVucLuuDinhKyGia(t)
			g := &thoGhiGia{}
			b := MoBoLuuDinhKy(s, g)

			t0 := time.Now()
			b.dongHo = func() time.Time { return t0 }
			b.lanLuuCuoi = t0

			b.kiemVaLuu(t0.Add(tc.lech))

			ban := g.banSaoXep()
			if len(ban) != 2 {
				t.Fatalf("kiemVaLuu tai +%v Xep %d ban, mong doi 2 (moi nguoi trong moi khu vuc)", tc.lech, len(ban))
			}
			if !g.coId(nvA.ID) || !g.coId(nvB.ID) {
				t.Fatalf("kiemVaLuu tai +%v thieu nguoi: %+v (mong doi ca id %d va %d)", tc.lech, ban, nvA.ID, nvB.ID)
			}
		})
	}
}

// TestBoLuuDinhKySauKhiLuuCapNhatLanLuuCuoi xac nhan kiemVaLuu CAP NHAT
// lanLuuCuoi sau lan luu thanh cong — neu khong cap nhat, MOI lan kiem sau do
// (du chi vai giay sau, chua toi 60 giay tinh tu lan luu THAT) se lai luu
// THUA. Goi kiemVaLuu HAI LAN: lan dau qua nguong (+60s tinh tu t0, luu), lan
// hai chi CACH LAN DAU 5 giay (+65s tinh tu t0) — neu lanLuuCuoi da cap nhat
// dung thanh +60s, lan hai chi moi +5s ke tu do nen KHONG duoc luu them.
func TestBoLuuDinhKySauKhiLuuCapNhatLanLuuCuoi(t *testing.T) {
	s, _, _ := soKhuVucLuuDinhKyGia(t)
	g := &thoGhiGia{}
	b := MoBoLuuDinhKy(s, g)

	t0 := time.Now()
	b.dongHo = func() time.Time { return t0 }
	b.lanLuuCuoi = t0

	b.kiemVaLuu(t0.Add(60 * time.Second)) // lan 1: du nguong -> luu (2 ban)
	if soLan := len(g.banSaoXep()); soLan != 2 {
		t.Fatalf("sau lan kiem dau (+60s), da Xep %d ban, mong doi 2", soLan)
	}

	b.kiemVaLuu(t0.Add(65 * time.Second)) // lan 2: moi cach lan luu truoc 5s -> KHONG luu them
	if soLan := len(g.banSaoXep()); soLan != 2 {
		t.Fatalf("sau lan kiem hai (+65s, chi cach lan luu truoc 5s), tong so ban Xep = %d, mong doi VAN LA 2 "+
			"(lanLuuCuoi phai duoc cap nhat sau lan luu dau, khong duoc luu thua)", soLan)
	}
}

// TestBoLuuDinhKyBatDauVaDungSach kiem vong doi goroutine nen: BatDau roi
// Dung phai tra ve SACH (khong treo = khong ro goroutine), dung khuon
// TestKhuVucActorKhoiDongVaDungSach (khuvuc_test.go) va TestDungTatCa
// (sokhuvuc_test.go). KHONG kiem logic 60 giay o day (da kiem het o cac test
// tren qua kiemVaLuu truc tiep) — chi kiem duong ong dong/mo co sach hay
// khong, va Dung() goi hai lan khong duoc panic (dungMotLan).
func TestBoLuuDinhKyBatDauVaDungSach(t *testing.T) {
	s, _, _ := soKhuVucLuuDinhKyGia(t)
	g := &thoGhiGia{}
	b := MoBoLuuDinhKy(s, g)

	b.BatDau()
	b.Dung() // phai tra ve, khong treo (neu treo, go test se timeout)
	b.Dung() // goi lan hai khong duoc panic
}

// TestBoLuuDinhKyVongLapThatTuDongLuuKhiDuNguong la test WIRING: khac cac
// test tren (goi kiemVaLuu TRUC TIEP, bo qua goroutine that), test nay chay
// QUA BatDau() + ticker THAT de xac nhan day du duong day "ticker -> goi
// kiemVaLuu(dongHo())" trong vongLap that su noi dung, khong chi logic
// kiemVaLuu bi co lap dung.
//
// AN TOAN DATA RACE: dongHo duoc GAN TRUOC khi goi BatDau() (goroutine con
// thay gia tri da xuat ban truoc "go" statement, dung Go memory model — xem
// chu thich truong dongHo, luudinhky.go) va KHONG BAO GIO bi ghi lai sau do;
// "dong ho gia" tien theo THOI GIAN THAT (t0 + time.Since(batDauThat)) nen
// khong can bien dem chia se nao giua goroutine test va goroutine vongLap —
// loai tru rui ro race ma khong can atomic.
func TestBoLuuDinhKyVongLapThatTuDongLuuKhiDuNguong(t *testing.T) {
	s, nvA, nvB := soKhuVucLuuDinhKyGia(t)
	g := &thoGhiGia{}
	b := MoBoLuuDinhKy(s, g)

	t0 := time.Now()
	batDauThat := time.Now()
	// "Bay gio" gia = t0 + (60s - 300ms) + thoi gian THAT da troi qua ke tu
	// luc BatDau() — sau ~300ms that (vai chuc tick cua nhipKiemLuuDinhKy),
	// dong ho gia se VUOT nguong 60s tinh tu lanLuuCuoi=t0, kich hoat luu THAT
	// qua chinh vongLap/ticker, khong phai qua loi goi kiemVaLuu truc tiep.
	b.dongHo = func() time.Time {
		return t0.Add(60*time.Second - 300*time.Millisecond).Add(time.Since(batDauThat))
	}
	b.lanLuuCuoi = t0

	b.BatDau()
	defer b.Dung()

	hetGio := time.After(3 * time.Second)
	for {
		if len(g.banSaoXep()) >= 2 {
			break
		}
		select {
		case <-hetGio:
			t.Fatalf("sau 3s that, vongLap that (qua BatDau+ticker) van chua tu dong luu — "+
				"da Xep %d ban, mong doi 2", len(g.banSaoXep()))
		case <-time.After(10 * time.Millisecond):
		}
	}

	if !g.coId(nvA.ID) || !g.coId(nvB.ID) {
		t.Fatalf("vongLap that Xep thieu nguoi: %+v (mong doi ca id %d va %d)", g.banSaoXep(), nvA.ID, nvB.ID)
	}
}
