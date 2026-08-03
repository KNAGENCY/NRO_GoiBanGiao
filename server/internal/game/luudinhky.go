package game

import (
	"sync"
	"time"

	"nroserver/internal/antoan"
	"nroserver/internal/store"
)

// nhipLuuDinhKy la HANG SO DUY NHAT quyet dinh chu ky luu dinh ky (Ke hoach
// 11 Task 4, dac ta 2026-07-27-luu-trang-thai-song.md §4 diem A + §9 muc 1 —
// QUYET DINH DA CHOT cua chu du an 27/07/2026: "Nhịp lưu = 60 giây. Sập thì
// mất tối đa 1 phút tiến độ ... Đặt thành hằng số một chỗ để đổi được dễ").
// Dat kieu hang so (khong phai bien) o DUY NHAT mot noi de (a) doi duoc de
// dang khi chu du an doi y, va (b) kiem thu dot bien (doi 60 giay thanh 600
// giay) sua TRUC TIEP gia tri o day, khong can sua noi khac.
const nhipLuuDinhKy = 60 * time.Second

// nhipKiemLuuDinhKy la tan suat DRIVER THAT (ticker cua vongLap ben duoi) —
// KHONG PHAI chu ky luu (dung nhipLuuDinhKy o tren cho viec do). Day chi la
// "bao lau nhin dong ho mot lan" de quyet dinh CO du nhipLuuDinhKy chua, dung
// y het cach ticker 100ms cua KhuVuc.vongLap (khuvuc.go) kiem cac nguong dai
// hon (5 giay hoi sinh, 30 giay hoi mau quai) ma khong can ticker rieng cho
// tung nguong — dung 100ms tai day de nhat quan phong cach toan du an va giu
// test wiring (TestBoLuuDinhKyVongLapThatTuDongLuuKhiDuNguong) nhanh, khong
// can cho hang giay that.
const nhipKiemLuuDinhKy = 100 * time.Millisecond

// boXep la phan toi thieu ThoGhi (thogi.go) ma LuuTatCa/BoLuuDinhKy can — CHI
// mot ham Xep, dung khuon interface HEP voi guiTin (khuvuc.go) va KhoLuu
// (thogi.go). *ThoGhi thoa man interface nay tu nhien (khong can glue code),
// nhung tach rieng interface o day (thay vi nhan thang *ThoGhi cu the) de
// test tiem duoc "tho ghi gia" ghi lai NGAY LAP TUC (dong bo) Xep nao da duoc
// goi, khong phai cho ThoGhi that flush qua goroutine nen cua no.
type boXep interface {
	Xep(nv store.NhanVat)
}

// LuuTatCa la HELPER DUNG CHUNG (Ke hoach 11 Task 4) — BAN GIAO cho Task 6
// (tat may chu, dac ta §4 diem C): ca hai diem moc A (dinh ky) va C (tat may
// chu) can y het mot thao tac "xin anh chup toan bo nguoi trong TAT CA khu
// vuc roi day vao hang doi cho tho ghi", chi khac o AI/KHI NAO goi no — Task
// 4 goi tu vong hen gio 60 giay, Task 6 se goi mot lan truoc khi DungTatCa()
// luc tat may chu. Vi vay dat thanh mot ham CONG KHAI, KHONG nhet chat vao
// vong hen gio (BoLuuDinhKy ben duoi), de goi duoc tu ben ngoai.
//
// Duyet TUNG khu vuc trong SoKhuVuc, xin AnhChupNguoiChoi(0) (0 = TAT CA
// nguoi trong khu vuc do, xem khuvuc.go) roi Xep TUNG ban vao tg.
//
// AN TOAN goi tu BEN NGOAI actor: AnhChupNguoiChoi tu no da di qua hop thu +
// tra ve BAN SAO GIA TRI (dac ta §5) — ham nay KHONG doc truc tiep bat ky
// truong noi bo nao cua KhuVuc, chi goi cac ham cong khai da co san (dung
// tinh than LUAT SO MOT, dac ta §3: khong nhet logic moi vao actor).
func (s *SoKhuVuc) LuuTatCa(tg boXep) {
	for _, k := range s.theoMap {
		for _, nv := range k.AnhChupNguoiChoi(0) {
			tg.Xep(nv)
		}
	}
}

// BoLuuDinhKy la GOROUTINE RIENG (Ke hoach 11 Task 4) chay nen, cu moi
// nhipLuuDinhKy lai goi LuuTatCa mot lan — TUYET DOI KHONG nhet vao nhipTim
// cua bat ky KhuVuc nao (dac ta §3 LUAT SO MOT + brief Task 4: "nhipTim chay
// 100ms/lan TRONG vongLap cua actor; nhet viec luu vao do la keo actor vao
// duong ghi dia"). Day la mot thanh phan HOAN TOAN DOC LAP voi actor cua
// tung khu vuc — no chi GOI cac ham cong khai (LuuTatCa -> AnhChupNguoiChoi
// -> Xep), khong bao gio dung truc tiep truong noi bo cua KhuVuc nao.
type BoLuuDinhKy struct {
	soKhuVuc *SoKhuVuc
	thoGhi   boXep

	// dongHo la DONG HO TIEM DUOC (dung khuon k.dongHo, khuvuc.go, KH6 Task 8
	// — xem chu thich truong do de biet ly do day du) de test kiem "chua toi
	// 60 giay" / "du 60 giay" MA KHONG can Sleep that 60 (hay 600, luc kiem
	// thu dot bien) giay that. MoBoLuuDinhKy gan mac dinh time.Now (hanh vi
	// that luc chay server that); test CUNG package gan thang truong nay
	// (khong co setter rieng) — CHI lam vay TRUOC khi goi BatDau(): Go
	// memory model dam bao goroutine con thay gia tri da ghi TRUOC cau lenh
	// "go" (khong can dong bo them), nhung sau BatDau() thi vongLap doc
	// truong nay moi nhip nhipKiemLuuDinhKy — ghi tu ben ngoai LUC DO la data
	// race chua dong bo. Moi test kiem logic nguong (kiemVaLuu) trong file
	// nay goi TRUC TIEP kiemVaLuu(), KHONG BatDau(), dung tinh than "goi
	// k.xuLy truc tiep, actor chua BatDau" da lap o khuvuc_test.go.
	dongHo func() time.Time

	// lanLuuCuoi la moc thoi gian cua lan luu dinh ky GAN NHAT (hoac luc khoi
	// tao BoLuuDinhKy, neu chua tung luu lan nao) — CHI vongLap (hoac test
	// goi truc tiep kiemVaLuu, xem chu thich dongHo o tren) duoc doc/ghi
	// truong nay, dung tinh than "mot luong duy nhat cham trang thai" da ap
	// dung cho KhuVuc/ThoGhi.
	lanLuuCuoi time.Time

	dung       chan struct{}
	xong       chan struct{}
	dungMotLan sync.Once
}

// MoBoLuuDinhKy tao MOT BoLuuDinhKy CHUA CHAY (dung khuon MoKhuVuc, khuvuc.go:
// KHONG tu khoi dong goroutine — goi BatDau() rieng de bat dau). Tach lam hai
// buoc, KHAC voi MoThoGhi (thogi.go, tu khoi dong luon), de test co the gan
// dongHo/lanLuuCuoi TRUOC khi bat ky goroutine nao co the doc chung (xem chu
// thich truong dongHo o tren).
//
// lanLuuCuoi khoi tao = dongHo() TAI LUC GOI (dung tinh than MoKhuVuc gan
// LanHoiMauCuoi = luc spawn cho tung con quai, khuvuc.go): tranh zero value
// (nam 1) khien lan kiem DAU TIEN sau BatDau() hieu nham "da qua 60 giay tu
// lau", gay mot lan luu THUA ngay luc may chu vua khoi dong, chua ai kip choi
// gi.
func MoBoLuuDinhKy(soKhuVuc *SoKhuVuc, thoGhi boXep) *BoLuuDinhKy {
	b := &BoLuuDinhKy{
		soKhuVuc: soKhuVuc,
		thoGhi:   thoGhi,
		dongHo:   time.Now,
	}
	b.lanLuuCuoi = b.dongHo()
	return b
}

// BatDau khoi dong goroutine nen cua bo luu dinh ky. Goi mot lan sau khi da
// thiet lap xong (vd sau khi tiem dongHo rieng cho test, neu can).
func (b *BoLuuDinhKy) BatDau() {
	b.dung = make(chan struct{})
	b.xong = make(chan struct{})
	go b.vongLap()
}

// vongLap la than cua bo luu dinh ky: ticker nhipKiemLuuDinhKy (driver THAT,
// KHAC nhipLuuDinhKy — xem chu thich hang so o tren) goi kiemVaLuu() moi lan,
// hoac thoat sach khi co lenh dung. Cau truc select giong het KhuVuc.vongLap
// (khuvuc.go) va ThoGhi.vongLap (thogi.go) — nhat quan phong cach
// actor/goroutine-nen toan du an.
func (b *BoLuuDinhKy) vongLap() {
	ticker := time.NewTicker(nhipKiemLuuDinhKy)
	defer ticker.Stop()
	defer close(b.xong)

	// LUOI PANIC (manh A §3.2c #4): boc THAN moi vong, KHONG boc ngoai — boc
	// ngoai thi panic giet ca vong lap, mat han tinh nang luu dinh ky. Mot nhip
	// luu loi -> bo nhip do, vong van quay nhip sau.
	for {
		select {
		case <-ticker.C:
			antoan.ChayCoLuoi("luudinhky.nhip", func() { b.kiemVaLuu(b.dongHo()) })
		case <-b.dung:
			return
		}
	}
}

// kiemVaLuu la logic THUAN "da du nhipLuuDinhKy chua" — tach rieng khoi
// vongLap (dung ly do KH6 Task 9 tach hoiSinhQuai/hoiMauQuai/aiQuai khoi
// nhipTim: de test goi TRUC TIEP voi bayGio tiem san, khong phu thuoc ticker
// that). CHUA du nhipLuuDinhKy ke tu lanLuuCuoi -> KHONG lam gi (chua luu gi
// ca). DU/QUA -> goi LuuTatCa MOT LAN roi cap nhat lanLuuCuoi = bayGio (dat
// TRUOC khi goi LuuTatCa, dung tinh than q.LanDanhCuoi duoc dat TRUOC khi
// phat goi trong case suKienDanhQuai, khuvuc.go: "gio luu" phan anh dung luc
// BAT DAU luu, va dam bao lanLuuCuoi luon duoc cap nhat du LuuTatCa co mat
// bao lau di nua).
func (b *BoLuuDinhKy) kiemVaLuu(bayGio time.Time) {
	if bayGio.Sub(b.lanLuuCuoi) < nhipLuuDinhKy {
		return
	}
	b.lanLuuCuoi = bayGio
	b.soKhuVuc.LuuTatCa(b.thoGhi)
}

// Dung dung goroutine nen mot cach sach: dong kenh dung (mot lan, an toan khi
// goi nhieu lan) roi cho vongLap thoat han truoc khi tra ve — dam bao khong
// ro goroutine, dung khuon KhuVuc.Dung (khuvuc.go).
func (b *BoLuuDinhKy) Dung() {
	b.dungMotLan.Do(func() {
		close(b.dung)
	})
	<-b.xong
}
