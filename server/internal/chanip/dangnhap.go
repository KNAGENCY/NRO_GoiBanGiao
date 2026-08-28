// Package chanip giu hai cai gac nho, dung chung cho ca may chu, deu dem
// theo DIA CHI IP: mot cai chan do mat khau (GacDangNhap), mot cai gioi han
// so ket noi tren mot IP (GacKetNoi).
//
// VI SAO tach thanh goi rieng thay vi nhet vao game/ hay net/: CA HAI goi do
// deu can gac. net/ mo cong va dem ket noi; game/ xu ly dang nhap va dem lan
// sai. Neu dat o game/ thi net/ phai import game/ — vong import. Dat rieng o
// day, chanip CHI phu thuoc thu vien chuan (sync, time), nen ai import cung
// duoc va khong the sinh vong. Dung y het ly do package antoan (luoi chan
// panic) duoc tach ra.
//
// KY LUAT ACTOR: moi ham trong goi nay chi khoa mot mutex trong vai chuc
// nano-giay roi nha ngay — khong I/O, khong cho, khong goi nguoc len tang
// tren. Bai hoc §6.10b cua du an la goi thu gi cham tu trong actor khu vuc la
// tu khoa chac chan; gac nay duoc goi tu duong dang nhap va tu vong Accept
// nen phai re va khong bao gio chan.
package chanip

import (
	"sync"
	"time"
)

// Cac gia tri mac dinh cho gac chong do mat khau.
//
// Cach hieu: moi IP co mot "tui the". Mot lan dang nhap SAI tieu mot the.
// Het the thi bi tu choi cho toi khi the day lai. The tu day lai theo thoi
// gian, moi NhipHoiPhucMacDinh duoc mot the, toi da bang SucChuaMacDinh.
//
// VI SAO chon 10 the / 20 giay mot the:
//
//   - 10 the: nguoi choi that go nham mat khau bao nhieu lan la binh thuong?
//     Ba den nam lan la nhieu roi. Chon 10 la GAP DOI muc do, de mot nguoi
//     dang thu lai vai mat khau cu cua ho cung khong bi da ra. Nguong cua ban
//     Java cu la 15 lan sai -> khoa IP 15 phut; con so o day chat hon mot chut
//     nhung KHONG khoa cung, chi cham lai (xem duoi).
//
//   - 20 giay mot the: sau khi het the, mot IP van duoc thu tiep 3 lan/phut,
//     tuc 180 lan/gio. Do bua mat khau voi 180 lan/gio la vo vong (mot mat
//     khau 6 ky tu chu-so co hon hai ty to hop), trong khi mot nguoi choi that
//     lo bi dinh nguong van chi phai cho hai muoi giay la thu lai duoc — ho
//     KHONG bi khoa ra ngoai, chi bi cham.
//
// So do voi truoc khi va: 200 lan sai lien tiep chay het trong ~2,6 giay
// (~76 lan/giay). Sau khi va, 200 lan sai do can 10 the dau + 190 x 20 giay
// ~ 63 PHUT. Nhanh hon 1.400 lan ve phia may chu.
//
// DIEM YEU DA BIET, ghi ra day de khong ai tuong nham la da kin: nha mang
// Viet Nam dung NAT lon, hang nghin nguoi choi co the cung mot dia chi IP
// cong cong. Neu mot IP nhu vay co qua nhieu nguoi go nham cung luc, ho se
// va lay nhau. Do la ly do hai con so nay dat duoc bang bien moi truong
// (NRO_CHONGDO_SUCCHUA, NRO_CHONGDO_HOIPHUC_GIAY — xem cmd/nroserver/main.go)
// de nguoi van hanh noi long duoc NGAY trong ngay mo cua ma khong phai
// bien dich lai may chu.
const (
	SucChuaMacDinh     = 10
	NhipHoiPhucMacDinh = 20 * time.Second
)

// Cac moc don rac. Xem donRacNeuCan.
const (
	// nhipDonRacThuong: khi ban do con nho, quet moi phut mot lan la du.
	nhipDonRacThuong = time.Minute

	// nguongDonRacSom: qua so muc nay thi quet day hon (nhipDonRacSom), vi
	// day la dau hieu dang bi rai IP hang loat.
	nguongDonRacSom = 4096

	// nhipDonRacSom: KE CA khi ban do dang phinh, van khong quet nhanh hon
	// mot lan moi giay. Neu bo chan nay, mot ban do lon ma khong don duoc muc
	// nao (moi IP deu con dang hoat dong) se bi quet lai o MOI lan goi — moi
	// lan quet la O(so muc), tuc chi phi tang binh phuong va bien chinh cai
	// gac thanh mot duong lam nghen may chu.
	nhipDonRacSom = time.Second
)

// oNhoIP la trang thai cua MOT dia chi IP.
type oNhoIP struct {
	// soThe la so the CON LAI, de duoi dang so thuc chu khong phai so
	// nguyen: the day lai lien tuc theo thoi gian chu khong nhay tung nac,
	// nen "con 3,4 the" la trang thai co that va phai giu lai — lam tron
	// xuong o moi lan cham vao se an mat phan the vua hoi duoc.
	soThe float64

	// lanCham la lan CUOI CUNG o nho nay duoc tinh lai. Moi lan cham vao, so
	// the duoc cong them theo dung khoang thoi gian da troi qua ke tu moc
	// nay — nho vay khong can bat ky goroutine hen gio nao chay nen.
	lanCham time.Time
}

// GacDangNhap dem so lan dang nhap SAI theo tung dia chi IP va tu choi khi
// mot IP sai qua nhieu.
//
// TAI SAO DEM THEO IP MA KHONG DEM THEO TEN TAI KHOAN: neu khoa theo ten,
// ke xau chi can co tinh go sai mat khau cua NGUOI KHAC vai chuc lan la khoa
// duoc tai khoan cua ho — bien mot lop bao ve thanh mot vu khi de bat nat
// nguoi choi khac (tu choi dich vu nham muc tieu). Ban Java cu roi dung vao
// bay nay: no khoa han tai khoan khi phat hien bat thuong. Dem theo IP thi
// nguoi chiu hau qua chinh la may dang do, dung nguoi gay ra.
type GacDangNhap struct {
	mu     sync.Mutex
	theoIP map[string]*oNhoIP

	sucChua     float64
	nhipHoiPhuc time.Duration

	lanDonRac time.Time

	// dongHo tach ra thanh truong de test tua duoc thoi gian di toi ma khong
	// phai ngoi cho that hai muoi giay cho mot the hoi lai.
	dongHo func() time.Time
}

// MoGacDangNhap tao mot gac rong, dung duoc ngay.
//
// sucChua <= 0 hoac nhipHoiPhuc <= 0 se roi ve gia tri mac dinh, de mot cau
// hinh sai (vd bien moi truong go nham) KHONG lang le bien thanh "khong gioi
// han gi ca" — do la kieu hong nguy hiem nhat vi khong co trieu chung nao.
func MoGacDangNhap(sucChua int, nhipHoiPhuc time.Duration) *GacDangNhap {
	if sucChua <= 0 {
		sucChua = SucChuaMacDinh
	}
	if nhipHoiPhuc <= 0 {
		nhipHoiPhuc = NhipHoiPhucMacDinh
	}
	return &GacDangNhap{
		theoIP:      make(map[string]*oNhoIP),
		sucChua:     float64(sucChua),
		nhipHoiPhuc: nhipHoiPhuc,
		dongHo:      time.Now,
	}
}

// ChoPhep tra ve true neu IP nay con duoc thu dang nhap them mot lan nua.
//
// Ham nay CHI HOI, khong tru the — tru the la viec cua GhiNhanSai. Tach doi
// nhu vay de mot lan dang nhap DUNG khong ton gi ca: nguoi choi binh thuong
// vao ra bao nhieu lan cung khong bao gio cham toi gioi han nay.
func (g *GacDangNhap) ChoPhep(ip string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	nay := g.dongHo()
	g.donRacNeuCan(nay)
	return g.layVaTinhLai(ip, nay).soThe >= 1
}

// GhiNhanSai ghi nhan MOT lan dang nhap sai tu ip: tru mot the.
//
// Goi ham nay cho MOI lan dang nhap khong thanh cong, ke ca khi ten dang nhap
// sai dinh dang hay tai khoan khong ton tai — vi dung "kieu that bai" de
// quyet dinh co tru the hay khong chinh la mo lai canh cua ma lo hong 2 vua
// bit: no bien so the con lai thanh mot cai chuong bao "ten nay co that".
//
// Tra ve true DUNG MOT LAN cho moi dot: lan sai vua lam IP nay het the. Noi
// goi dung no de ghi nhat ky MOT dong cho moi dot bi chan, thay vi mot dong
// cho moi lan bi chan — ke xau do mat khau co the go cua hang nghin lan mot
// giay, moi lan mot dong log la tu tay bien lop bao ve thanh mot duong lam
// day o dia (dung bai hoc da ghi o dau main.go: mot ket noi ghi log moi goi
// tin tung khuech dai thanh 9 GB/gio).
func (g *GacDangNhap) GhiNhanSai(ip string) (vuaHetThe bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	nay := g.dongHo()
	g.donRacNeuCan(nay)
	o := g.layVaTinhLai(ip, nay)
	conThe := o.soThe >= 1
	o.soThe--
	// Khong cho am: neu de am, mot dot do bua 10.000 lan se dao mot cai ho sau
	// toi muc IP do bi khoa hang ngay troi — vuot xa y do "cham lai", va thanh
	// mot cai bay cho nguoi choi that dung chung IP.
	if o.soThe < 0 {
		o.soThe = 0
	}
	return conThe && o.soThe < 1
}

// SoIPDangNho cho biet ban do dang giu bao nhieu muc. Dung cho test va cho
// nguoi van hanh soi khi nghi ngo ro ri bo nho.
func (g *GacDangNhap) SoIPDangNho() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return len(g.theoIP)
}

// layVaTinhLai lay o nho cua ip (tao moi neu chua co) va cong bu so the da
// hoi duoc ke tu lan cham cuoi. PHAI goi khi dang giu g.mu.
func (g *GacDangNhap) layVaTinhLai(ip string, nay time.Time) *oNhoIP {
	o, co := g.theoIP[ip]
	if !co {
		// IP moi bat dau voi TUI THE DAY: mot nguoi choi chua tung sai lan nao
		// phai duoc huong tron suc chua.
		o = &oNhoIP{soThe: g.sucChua, lanCham: nay}
		g.theoIP[ip] = o
		return o
	}
	if troi := nay.Sub(o.lanCham); troi > 0 {
		o.soThe += float64(troi) / float64(g.nhipHoiPhuc)
		if o.soThe > g.sucChua {
			o.soThe = g.sucChua
		}
		o.lanCham = nay
	}
	return o
}

// donRacNeuCan xoa cac o nho khong con mang thong tin gi.
//
// VI SAO BAT BUOC PHAI CO: khong don thi ban do lon dan theo tong so dia chi
// IP tung ghe qua may chu va khong bao gio nho lai. Voi IPv6 thi so dia chi
// la vo han tren thuc te — ke xau chi can doi IP moi lan thu la bien chinh
// cai gac chong do thanh mot duong lam can bo nho may chu, tuc va mot lo
// hong roi tu tay dao mot lo khac.
//
// Chi xoa nhung o DA HOI DAY THE: mot o day the khong khac gi mot IP chua
// tung thay (layVaTinhLai se tao lai y het), nen xoa no KHONG lam mat mot
// chut thong tin nao. Nho vay cach don rac nay khong the vo tinh tha cho ke
// dang do do.
//
// He qua ve bo nho, noi bang loi thuong: ban do chi giu nhung IP vua dang
// nhap sai trong khoang gan day (day the lai tu 0 mat sucChua x nhipHoiPhuc,
// mac dinh la 10 x 20 giay = khoang 3 phut ruoi). Kich thuoc on dinh cua no
// bang so IP co nguoi go sai mat khau trong ba phut ruoi vua roi — mot con
// so gan voi luong nguoi choi that, khong phai voi tong luong tan cong.
//
// PHAI goi khi dang giu g.mu.
func (g *GacDangNhap) donRacNeuCan(nay time.Time) {
	keTuLanTruoc := nay.Sub(g.lanDonRac)
	quaDong := len(g.theoIP) >= nguongDonRacSom
	if !((quaDong && keTuLanTruoc >= nhipDonRacSom) || keTuLanTruoc >= nhipDonRacThuong) {
		return
	}
	g.lanDonRac = nay
	for ip, o := range g.theoIP {
		troi := nay.Sub(o.lanCham)
		if troi <= 0 {
			continue
		}
		if o.soThe+float64(troi)/float64(g.nhipHoiPhuc) >= g.sucChua {
			delete(g.theoIP, ip)
		}
	}
}
