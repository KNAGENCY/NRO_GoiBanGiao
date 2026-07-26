package game

import (
	"log/slog"
	"sync"
	"time"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// KhuVuc (zone) la don vi khong gian nho nhat cua mot ban do: giu ban do dang
// choi + danh sach nguoi choi dang o trong do. Ke hoach 4 CHI can mot nguoi dung
// yen trong map trong, nen day la mo hinh TOI THIEU — khong goroutine rieng,
// khong mutex, khong nhip 100ms.
//
// VI SAO CHUA DUNG MO HINH ACTOR (goroutine moi zone) NGAY: actor chi can khi co
// TRANH CHAP that — nhieu luong cung doc/ghi trang thai zone (di chuyen, nhin thay
// nhau, quai di lai). Ke hoach 4 chua co thu nao trong so do: mot nguoi, dung yen,
// khong broadcast. Dung actor bay gio la giai bai toan chua ton tai (YAGNI). Ke
// hoach 5 (di chuyen) se nang KhuVuc thanh actor khi co tranh chap that su, kem
// nhip 100ms cho quai/vat pham (xem 2026-07-23-hoa-giai-muot §4.3). O KH4, KhuVuc
// chi lo NAP ban do va lam nen cho KH5.
//
// KE HOACH 5: KhuVuc nay tro thanh ACTOR — mot goroutine rieng so huu toan bo
// trang thai (nguoiTrongKhuVuc), nhan su kien qua hop thu (channel) va co nhip
// tim 100ms cho cac cong viec dinh ky (quai/vat pham o KH6). Moi thay doi trang
// thai PHAI di qua hop thu — khong truy cap truc tiep field tu goroutine khac.
type KhuVuc struct {
	banDo     *BanDo
	nguoiChoi []*store.NhanVat

	nguoiTrong map[int64]*nguoiTrongKhuVuc

	// quai la danh sach quai dang song trong khu vuc nay (Ke hoach 6 Task 5).
	// Tao MOT LAN trong MoKhuVuc; sau do CHI luong vongLap duoc doc/ghi (dac ta
	// §1.6: "chi luong vongLap duoc cham trang thai quai va trang thai chien
	// dau"). KHONG co getter doc truc tiep tu ben ngoai — luong khac (vd
	// luong phien khi dung goi -24) phai gui suKienXinDanhSachQuai qua hop thu
	// va nhan BAN SAO qua kenh traLoi (xem suKienXinDanhSachQuai ben duoi).
	quai []*Quai

	// dongHo la DONG HO TIEM DUOC cho toan bo logic co moc thoi gian trong
	// actor (Ke hoach 6 Task 8: hoi sinh 5s + hoi mau 30s trong nhipTim, va
	// kiem bat tu trong case suKienDanhQuai) — RANG BUOC CHUNG cua brief Task
	// 8: "khong goi time.Now() thang trong logic" de test moi moc thoi gian
	// (4.9s/5.0s/29.9s/30.0s...) MA KHONG sleep that (khong lam bo test cham
	// them hang chuc giay). MoKhuVuc gan mac dinh time.Now (hanh vi that khi
	// chay server that); test CUNG package gan thang truong nay (khong co
	// setter rieng — dung tinh than "chi test cung package moi duoc dung truc
	// tiep trang thai noi bo" da ap dung cho k.quai/k.nguoiTrong o cac Task
	// truoc) de dieu khien "bay gio" theo y muon, khong sleep.
	dongHo func() time.Time

	hopThu     chan suKienKhuVuc
	dung       chan struct{}
	xong       chan struct{}
	dungMotLan sync.Once
}

// guiTin la phan toi thieu cua mot phien (session) ma actor can de gui goi tin
// cho client. *mynet.Session thoa man interface nay; test dung phien gia.
type guiTin interface {
	Gui(*protocol.Message) error
}

// nguoiTrongKhuVuc gan mot nhan vat voi phien mang cua no trong pham vi mot
// khu vuc — de actor biet gui goi tin cho ai khi co su kien.
type nguoiTrongKhuVuc struct {
	nv    *store.NhanVat
	phien guiTin

	// chieuDangChon la CapChieu (cap 1) cua chieu nguoi choi vua chon bang goi
	// 34 (xem suKienChonChieu). nil = chua chon chieu nao. CHI luong vongLap
	// cua actor duoc doc/ghi truong nay — day chinh la trang thai dac ta §3.2
	// yeu cau dat trong actor (khong trong BoXuLy) de tranh tranh chap voi
	// luong tinh sat thuong doc no sau nay (Ke hoach 6 Task 3+).
	chieuDangChon *CapChieu

	// lanDanhCuoi la thoi diem nguoi choi nay DUNG THANH CONG mot chieu lan
	// gan nhat — Task 3 (chotkiem.go) can truong nay de tinh "da troi qua bao
	// lau" cho kiem hoi chieu. store.NhanVat KHONG co truong thoi gian nao
	// (da kiem: chi co cac truong chi so), nen buoc phai dat o day — dung vi
	// tri voi chieuDangChon (chi actor duoc doc/ghi). Zero value (chua tung
	// dung chieu nao) khien bayGio.Sub(lanDanhCuoi) ra mot khoang RAT LON, nen
	// kiem hoi chieu tu nhien CHO PHEP ngay lan dau — khong can if/else rieng.
	lanDanhCuoi time.Time

	// lucHoiSinh la thoi diem nguoi choi nay HOI SINH GAN NHAT (Ke hoach 6
	// Task 13, dac ta §1.7/§3.4/§6.4) — dat trong case suKienHoiSinh duoi
	// day. CUNG CO CHE voi Quai.LucHoiSinh (quai.go, Task 8): CHI luong
	// vongLap cua actor duoc doc/ghi (dac ta §1.6), va dung lam nguon cho
	// ViTriCoMau.LucHoiSinh (chotkiem.go) khi Task 11 (PvP) can kiem bat tu
	// 1500ms cho MUC TIEU la nguoi choi — nguoiTrongKhuVuc gio da du bo
	// truong (X/Y qua nv.ViTriX/Y, HpHienTai qua nv.HpHienTai, LucHoiSinh qua
	// truong nay) de dung lam nguon ViTriCoMau, dung y brief Task 13. Zero
	// value (nguoi choi CHUA TUNG hoi sinh lan nao tu luc dang nhap) khien
	// bayGio.Sub(zero) ra khoang thoi gian RAT LON -> kiem bat tu tu nhien
	// CHO PHEP ngay tu dau, dung hanh vi mong doi — GIONG HET ly do zero
	// value cua Quai.LucHoiSinh (quai.go) an toan.
	lucHoiSinh time.Time
}

// suKienKhuVuc la kieu su kien actor xu ly qua hop thu. Ke hoach 5 Task 5 se
// lam day logic xu ly; Task 1 chi dinh nghia khung.
type suKienKhuVuc interface{}

// suKienVao bao actor mot nhan vat vua vao khu vuc, kem phien de gui goi tin.
type suKienVao struct {
	nv    *store.NhanVat
	phien guiTin
}

// suKienDiChuyen bao actor mot nhan vat doi vi tri.
type suKienDiChuyen struct {
	idNhanVat int64
	x, y      int16
}

// suKienRa bao actor mot nhan vat vua roi khu vuc.
type suKienRa struct {
	idNhanVat int64
}

// suKienChonChieu bao actor mot nhan vat vua chon mot chieu (goi len `34`,
// dac ta §3.2). mauChieuId la template.id (KHONG phai skillId — hai truong
// khac nhau, xem bang §3.2 cua dac ta). Luong phien CHI gui su kien nay, actor
// moi la noi tra cuu ChieuCuaLang + kiem chieu co thuoc lang nhan vat khong
// va ghi chieuDangChon — tranh dung bay §6.9 (tranh chap doc/ghi nv.ViTriX/Y
// ma Ke hoach 5 da vap).
type suKienChonChieu struct {
	idNhanVat  int64
	mauChieuId byte
}

// suKienDanhQuai bao actor mot nhan vat vua gui goi len `54` (dac ta §3.1):
// xin danh con quai co Id = mobId trong khu vuc. mobId la byte THO doc truc
// tiep tu Reader (khong dau) — gia tri 255 (= -1 khi doc CO dau, dac ta §1.4)
// la truong hop "thu cung" (dac ta §0 "Ngoai KH6", chua lam): actor KHONG can
// nhanh rieng cho no, vi Quai.Id LUON nam trong [0,127] (taoQuaiMacDinh,
// quai.go) nen 255 tu nhien khong khop con quai nao — roi thang vao nhanh
// "khong tim thay muc tieu" -> im lang bo qua, dung tinh than voi cach
// suKienDiChuyen bo qua khi khong tim thay nguoi trong k.nguoiTrong.
//
// Luong phien (xuLyDanhQuai, login.go) CHI doc gia tri mobId nay tu goi (ke
// ca doc BO 4 byte masterId khi mobId==-1 de tieu thu DU than goi bien thien,
// dac ta §2) roi gui su kien nay vao hop thu — KHONG tu tra k.quai, KHONG tu
// tinh sat thuong. Toan bo kiem tra (KiemDanhHopLe) + tinh sat thuong + tru
// mau + quyet chet + dung/phat goi nam TRONG vongLap (xem case duoi day) —
// dung ky luat actor dac ta §1.6 va rang buoc cua Task 7.
type suKienDanhQuai struct {
	idNhanVat int64
	mobId     byte
}

// suKienHoiSinh bao actor mot nhan vat vua gui goi len -16 (dac ta §3.1, Ke
// hoach 6 Task 13): xin hoi sinh sau khi chet. Luong phien (xuLyHoiSinh,
// login.go) CHI gui su kien nay — actor moi la noi kiem "dang THAT SU chet"
// + hoi mau/chakra day du + dat lucHoiSinh + phat goi -16/-30-15 (dung ky
// luat actor §1.6, xem case suKienHoiSinh duoi day). Dung ten truong
// idNhanVat giong het cac su kien khac trong file nay (suKienDiChuyen,
// suKienRa, suKienDanhQuai) de nhat quan.
type suKienHoiSinh struct {
	idNhanVat int64
}

// heQuaiChuaXacDinh la gia tri he nguyen to gan cho MOI con quai khi tinh sat
// thuong (Task 7): Quai (Task 5, quai.go) CHUA co truong theo doi he nguyen
// to rieng cho tung con. Dat DUNG BANG soHeNguyenTo (satthuong.go) — tuc
// NGOAI khoang hop le [0, soHeNguyenTo) — de heSoNguyenToPhanTram tu nhien
// coi day la "he khong xac dinh" va tra ve 100% (trung tinh), khong can sua
// ham do hay them nhanh dac biet o day. Neu sau nay quai co he nguyen to that
// (vd theo mau quai), doi hang so nay thanh mot tra cuu that.
const heQuaiChuaXacDinh byte = soHeNguyenTo

// heNguyenToTheoLang anh xa Lang.ID (0..4, thu tu danhSachLang trong lang.go:
// La, Cat, Suong Mu, Da, May) sang hang so he nguyen to cua satthuong.go
// (HeHoa/HePhong/HeLoi/HeTho/HeThuy). HAI THU TU NAY LECH NHAU — KHONG duoc
// suy dien "lang i = he i": lang 2 (Suong Mu) la He Thuy (hang so 4), lang 4
// (May) la He Loi (hang so 2). Dung theo thiet ke tong §7.2 va chieuthuc.go
// (xem tenChieuTheoLang: "Lang Suong Mu — Thuy", "Lang May — Loi").
var heNguyenToTheoLang = [5]byte{HeHoa, HePhong, HeThuy, HeTho, HeLoi}

// heNguyenToNguoiChoi tra he nguyen to cua mot nhan vat theo lang, dung lam
// tham so heNguoiDanh cho TinhSatThuong (Task 7 — Task 4 chi dung ham, khong
// quyet dinh nguon du lieu nay). Lang ngoai [0,4] (du lieu hong, ve ly thuyet
// khong xay ra vi TimLang da chan luc tao nhan vat) tra heQuaiChuaXacDinh de
// roi vao nhanh trung tinh, khong panic.
func heNguyenToNguoiChoi(lang int16) byte {
	if lang < 0 || int(lang) >= len(heNguyenToTheoLang) {
		return heQuaiChuaXacDinh
	}
	return heNguyenToTheoLang[lang]
}

// suKienXinDanhSachQuai la su kien TRUY VAN (khac cac su kien "lam viec" o
// tren): mot luong BEN NGOAI vongLap (vd luong phien khi dung goi -24 MAP
// INFO o xuLyClientOk, TRUOC khi nhan vat duoc giao cho actor qua -39) can
// biet quai hien co trong khu vuc, nhung KHONG duoc doc truc tiep k.quai (dac
// ta §1.6). Luong do gui su kien nay qua Gui() roi cho tra loi tren traLoi;
// vongLap la noi DUY NHAT xu ly no va tra ve BAN SAO GIA TRI (con tro moi,
// KHONG PHAI con tro goc trong k.quai) de ben nhan khong the vo tinh sua
// trang thai actor tu ben ngoai.
type suKienXinDanhSachQuai struct {
	traLoi chan []*Quai
}

// MoKhuVuc tao mot khu vuc giu ban do da cho, chua co nguoi choi nao. Quai
// (Ke hoach 6 Task 5) duoc sinh MOT LAN o day bang taoQuaiMacDinh — day la
// diem duy nhat gan quai vao mot khu vuc, giu dung nguyen tac "quai thuoc
// actor, sinh luc MoKhuVuc" cua brief Task 5.
func MoKhuVuc(bd *BanDo) *KhuVuc {
	k := &KhuVuc{
		banDo:      bd,
		nguoiTrong: make(map[int64]*nguoiTrongKhuVuc),
		quai:       taoQuaiMacDinh(bd),
		dongHo:     time.Now,
	}
	// Ke hoach 6 Task 8: khoi tao LanHoiMauCuoi = thoi diem spawn cho MOI con
	// quai — tranh zero value (nam 1) khien lan kiem hoi mau DAU TIEN cua
	// nhipTim tuong nham "da qua 30s" ngay tick dau (xem chu thich truong
	// LanHoiMauCuoi, quai.go). Dat o day (sau khi da co k.dongHo) thay vi
	// truyen tham so vao taoQuaiMacDinh — giu chu ky taoQuaiMacDinh don gian,
	// khong ai khac goi ham do nen doi day khong pha vo gi.
	luc := k.dongHo()
	for _, q := range k.quai {
		q.LanHoiMauCuoi = luc
	}
	return k
}

// BanDo tra ban do cua khu vuc (client can no de dung MAP INFO).
func (k *KhuVuc) BanDo() *BanDo { return k.banDo }

// DanhSachQuaiSnapshot tra BAN SAO danh sach quai hien co trong khu vuc, AN
// TOAN goi tu BEN NGOAI vongLap (vd luong phien o xuLyClientOk khi dung goi
// -24 MAP INFO). Di qua hop thu actor (suKienXinDanhSachQuai) thay vi doc
// truc tiep k.quai — dung dac ta §1.6 "chi luong vongLap duoc cham trang
// thai quai". Cho toi da 100ms; qua han (actor da Dung(), hop thu day bat
// thuong, hoac k.hopThu chua duoc cap phat vi BatDau() chua chay) thi tra ve
// nil — luong phien khong bi treo, chi gui MAP INFO voi 0 quai (giong hanh
// vi KH4 cu, khong sap goi).
func (k *KhuVuc) DanhSachQuaiSnapshot() []*Quai {
	if k.hopThu == nil {
		return nil // actor chua BatDau() -> khong co hop thu de gui
	}
	traLoi := make(chan []*Quai, 1)
	k.Gui(suKienXinDanhSachQuai{traLoi: traLoi})
	select {
	case ds := <-traLoi:
		return ds
	case <-time.After(100 * time.Millisecond):
		slog.Warn("xin danh sach quai qua hop thu actor het gio cho", "mapId", k.banDo.MapID)
		return nil
	}
}

// ThemNguoiChoi dua mot nhan vat vao khu vuc. Ke hoach 5 se broadcast su xuat
// hien nay cho nguoi khac; Ke hoach 4 chi giu danh sach de lam nen.
func (k *KhuVuc) ThemNguoiChoi(nv *store.NhanVat) {
	k.nguoiChoi = append(k.nguoiChoi, nv)
}

// SoNguoiChoi tra so nguoi dang trong khu vuc (dung cho kiem thu + Ke hoach 5).
func (k *KhuVuc) SoNguoiChoi() int { return len(k.nguoiChoi) }

// BatDau khoi dong goroutine actor cua khu vuc: cap phat hop thu + kenh dieu
// khien roi chay vongLap trong goroutine rieng. Goi mot lan truoc khi dung Gui.
func (k *KhuVuc) BatDau() {
	k.hopThu = make(chan suKienKhuVuc, 256)
	k.dung = make(chan struct{})
	k.xong = make(chan struct{})
	go k.vongLap()
}

// vongLap la than actor: doc hop thu, nhip tim 100ms, hoac dung khi co lenh
// dung. Day la NOI DUY NHAT doc/ghi trang thai khu vuc (nguoiTrong, ...) —
// khong goroutine nao khac duoc dung truc tiep cac field nay.
func (k *KhuVuc) vongLap() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	defer close(k.xong)

	for {
		select {
		case sk := <-k.hopThu:
			k.xuLy(sk)
		case <-ticker.C:
			k.nhipTim()
		case <-k.dung:
			return
		}
	}
}

// xuLy xu ly mot su kien lay tu hop thu — day la NOI DUY NHAT doc/ghi trang
// thai khu vuc (goi tu vongLap, hoac truc tiep tu test khi actor CHUA BatDau).
//
// THU TU phat tin (quan trong, xem brief KH5 Task 5):
//   - suKienVao: gui -5 cua TUNG nguoi CU cho nguoi MOI TRUOC (nguoi moi thay
//     nguoi cu), roi moi them nguoi moi vao map (de khong tu gui -5 cho chinh
//     minh o buoc sau), roi phat -5 cua nguoi moi cho nhung nguoi con lai.
//   - suKienDiChuyen: cap nhat vi tri THAT (authoritative) truoc, roi phat -7
//     cho NGUOI KHAC (client tu ve chuyen dong cua chinh minh, khong gui lai).
//   - suKienRa: xoa khoi map TRUOC, roi phat -6 cho nguoi con lai (nguoi da
//     ra khong con trong map nen tu nhien khong nhan duoc goi cua chinh minh).
func (k *KhuVuc) xuLy(sk suKienKhuVuc) {
	switch sk := sk.(type) {
	case suKienVao:
		// 1) Gui -5 cua TUNG nguoi CU cho nguoi MOI (de nguoi moi thay nguoi cu).
		//    Lam TRUOC khi them nguoi moi vao map, de nguoi moi khong tu thay minh.
		for _, cu := range k.nguoiTrong {
			m, err := goiXuatHien(cu.nv)
			if err != nil {
				slog.Error("dung goi -5 nguoi cu that bai", "id", cu.nv.ID, "loi", err)
				continue
			}
			if err := sk.phien.Gui(m); err != nil {
				slog.Warn("gui -5 nguoi cu cho nguoi moi that bai", "loi", err)
			}
		}
		// 2) Them nguoi moi vao khu vuc.
		k.nguoiTrong[sk.nv.ID] = &nguoiTrongKhuVuc{nv: sk.nv, phien: sk.phien}
		// 3) Phat -5 cua nguoi moi cho NGUOI KHAC (khong gui lai chinh minh).
		m, err := goiXuatHien(sk.nv)
		if err != nil {
			slog.Error("dung goi -5 nguoi moi that bai", "id", sk.nv.ID, "loi", err)
			return
		}
		k.phatChoNguoiKhac(sk.nv.ID, m)

	case suKienDiChuyen:
		nguoi, co := k.nguoiTrong[sk.idNhanVat]
		if !co {
			return
		}
		// KE HOACH 5: kiem hop le PHAI o day (actor so huu nguoi.nv), KHONG o luong
		// phien — luong phien doc nv.ViTri se DUA voi actor ghi chung con tro. Kiem
		// bang vi tri AUTHORITATIVE hien tai.
		if !kiemDiChuyenHopLe(k.banDo, int16(nguoi.nv.ViTriX), int16(nguoi.nv.ViTriY), sk.x, sk.y) {
			slog.Debug("bo goi di chuyen bat hop le (ngoai bien / vuot nguong buoc)",
				"id", sk.idNhanVat, "mapId", k.banDo.MapID, "toX", sk.x, "toY", sk.y)
			return // giu vi tri cu, KHONG phat
		}
		nguoi.nv.ViTriX = int32(sk.x)
		nguoi.nv.ViTriY = int32(sk.y)
		k.phatChoNguoiKhac(sk.idNhanVat, goiDiChuyen(sk.idNhanVat, sk.x, sk.y))

	case suKienRa:
		if _, co := k.nguoiTrong[sk.idNhanVat]; !co {
			return // stub an toan
		}
		delete(k.nguoiTrong, sk.idNhanVat)
		// Phat -6 cho nhung nguoi CON LAI (nguoi ra da bi xoa nen khong nhan).
		k.phatChoNguoiKhac(sk.idNhanVat, goiRoiKhuVuc(sk.idNhanVat))

	case suKienChonChieu:
		nguoi, co := k.nguoiTrong[sk.idNhanVat]
		if !co {
			return
		}
		// Tim MauChieu co Id == mauChieuId TRONG DANH MUC CUA LANG nhan vat
		// (ChieuCuaLang(nv.Lang)) — day la buoc "kiem chieu thuoc lang" ma dac
		// ta §3.2 doi hoi actor lam (khong phai luong phien). Lang khong hop
		// le (ngoai 0..4) hoac templateId khong co trong 4 chieu cua lang do
		// -> khong tim thay -> BO QUA, GIU NGUYEN chieuDangChon cu (dung dac
		// ta: "chon chieu khong thuoc lang -> khong doi gi").
		danhMuc := ChieuCuaLang(byte(nguoi.nv.Lang))
		for i := range danhMuc {
			if danhMuc[i].Id == sk.mauChieuId && len(danhMuc[i].CacCap) > 0 {
				// Luon lay cap 1 (CacCap[0]): Ke hoach 6 Task 2 chua co mo hinh
				// cap do chieu cua tung nguoi choi (se den o task sau) — chon
				// chieu la chon LOAI, cap that do task khac quyet dinh.
				nguoi.chieuDangChon = &danhMuc[i].CacCap[0]
				return
			}
		}

	case suKienDanhQuai:
		// Dac ta §3.1 (goi len 54) + §1.6 (ky luat actor) + brief Task 7: toan
		// bo kiem tra/tinh toan/phat goi cho MOT don danh quai nam TRON trong
		// nhanh nay — luong phien (xuLyDanhQuai, login.go) chi gui su kien,
		// khong tu lam gi trong cac buoc duoi day.
		nguoi, co := k.nguoiTrong[sk.idNhanVat]
		if !co {
			return
		}

		// Ke hoach 6 Task 13 (test bat buoc cua brief): NGUOI DA CHET khong
		// duoc danh. Chu dinh dat CHOT NAY o day (truoc ca KiemDanhHopLe),
		// khong nhet vao chotkiem.go: KiemDanhHopLe (Task 3) kiem MUC TIEU
		// (mucTieu.HpHienTai — kiem #4), con day la kiem NGUOI DANH — hai
		// huong khac nhau, brief Task 13 CHI liet ke goichiendau.go/khuvuc.go
		// la file phai sua, KHONG co chotkiem.go. Ban Java (Controller.java
		// case 54) KHONG co chot nay (chi -20 nhat do co !player.isDie()) —
		// day la mot luat choi RO RANG can them, khong phai chep nguyen ban
		// goc: nguoi choi dang cho hoi sinh khong the "tranh thu" danh quai.
		if nguoi.nv.HpHienTai <= 0 {
			return // im lang bo qua, giong het cach cac kiem khac cua case nay tu choi
		}

		// Tim con quai co Id == sk.mobId. mobId=255 (thu cung, xem chu thich
		// suKienDanhQuai o tren) tu nhien khong khop quai nao -> mucTieu==nil
		// -> bo qua, dung nhu suKienDiChuyen bo qua khi khong tim thay nguoi.
		var mucTieu *Quai
		for _, q := range k.quai {
			if q.Id == sk.mobId {
				mucTieu = q
				break
			}
		}
		if mucTieu == nil {
			return
		}

		bayGio := k.dongHo() // Task 8: dong ho tiem duoc, KHONG goi time.Now() thang o day nua.
		ketQua := KiemDanhHopLe(nguoi, ViTriCoMau{
			X:         mucTieu.X,
			Y:         mucTieu.Y,
			HpHienTai: mucTieu.Hp,
			// LucHoiSinh: Ke hoach 6 Task 8 da them truong nay vao Quai
			// (quai.go) va nhipTim gan gia tri that moi lan hoi sinh -> doc
			// THANG tu mucTieu o day, KHONG con bo trong nhu truoc Task 8.
			// Zero value (quai moi spawn, chua tung hoi sinh) van tu nhien
			// cho phep danh ngay — dung dac diem da ghi o chu thich truong
			// LucHoiSinh (chotkiem.go va quai.go).
			LucHoiSinh: mucTieu.LucHoiSinh,
		}, bayGio)
		if !ketQua.ChoPhep {
			// 4 kiem dau (bat tu/chakra/hoi chieu/da chet) hong -> IM LANG,
			// KHONG phat goi gi (khac han nhanh Truot duoi day — dac ta §1.9
			// chi noi rieng ve vuot tam, KHONG ap dung cho 4 kiem nay).
			return
		}

		// Cap nhat lanDanhCuoi CHI KHI duoc phep danh (brief Task 7). Mot don
		// bi 4 kiem dau tu choi KHONG duoc tinh la "vua dung chieu", neu
		// khong nguoi choi se bi khoa hoi chieu oan boi chinh don bi tu choi.
		nguoi.lanDanhCuoi = bayGio

		// chiMang: KH6 CHUA chot ty le/cong thuc chi mang that (khong dac ta
		// nao, khong task truoc nao dinh nghia) -> luon false — lua chon an
		// toan, xem bao cao Task 7. TinhSatThuong van nhan tham so nay dung
		// chu ky da chot o Task 4; ty le chi mang that (neu co, co the la
		// KH8) chi can thay DONG GAN chiMang nay, khong doi gi khac o day.
		chiMang := false

		var satThuong int32
		if !ketQua.Truot {
			lang := TimLang(nguoi.nv.Lang)
			var satThuongGoc int64
			if lang != nil {
				// satThuongGoc: dung Lang.SatThuongKhoiDau (lang.go) — truong
				// nay co san tu Ke hoach 2, chu thich goc ghi "dung o Ke
				// hoach 4" (ten cu cua chien dau truoc khi tach thanh KH6),
				// dam bao KHAC 0 cho ca 5 lang (10-16). nv.SucManh KHONG
				// dung duoc o day: la diem xep hang hien thi rieng (bang xep
				// hang web, xem web/app/bang-xep-hang), mac dinh 0 va KHONG
				// duoc TaoNhanVat gan luc tao nhan vat — dung no se lam MOI
				// nhan vat moi tao gay 0 sat thuong vinh vien.
				satThuongGoc = int64(lang.SatThuongKhoiDau)
			}
			heNguoiDanh := heNguyenToNguoiChoi(nguoi.nv.Lang)
			satThuong = TinhSatThuong(satThuongGoc, nguoi.chieuDangChon, heNguoiDanh, heQuaiChuaXacDinh, chiMang)

			// Ke hoach 6 Task 9: don TRUNG (nhanh nay, khac nhanh Truot o
			// duoi) -> ghi nhan nguoi nay vao danh sach "da tung danh" cua
			// quai, kich hoat tang uu tien 1 cua AI (dac ta §5.2, quai.go/
			// chonMucTieu) cho NHUNG DON SAU. Dat NGAY trong nhanh "trung
			// don" (khong phai nhanh Truot) vi "trung don" (brief Task 9)
			// nghia la thuc su gay sat thuong, khac "vuot tam" noi don khong
			// thuc su cham muc tieu.
			if mucTieu.DaTungDanh == nil {
				mucTieu.DaTungDanh = make(map[int64]bool)
			}
			mucTieu.DaTungDanh[nguoi.nv.ID] = true
		}
		// Truot=true (vuot tam): satThuong giu nguyen 0 (dac ta §1.9) — KHONG
		// goi TinhSatThuong, vi don nay khong thuc su cham muc tieu.

		// EP TUONG MINH: vuot tam thi KHONG duoc chi mang — BAN GIAO Task 3
		// (chotkiem.go:75-81): "chiMang=false khi Truot=true", vi client CHI
		// hien chu "miss" khi chiMang==false (Controller.cs:4356-4363, dac ta
		// §3.3 diem 3) — Truot=true+chiMang=true se hien "-0" do thay vi
		// "miss". Dat DONG NAY o day, SAU moi cho co the gan/tinh chiMang (ke
		// ca cong thuc chi mang that trong tuong lai o nhanh !ketQua.Truot
		// phia tren): hien tai chiMang luon false nen dong nay dung "tinh co"
		// — mot khi co cong thuc chi mang that duoc nhet vao dung cho o tren,
		// dong nay la RANH GIOI CUOI CUNG bao dam nghia vu khong bao gio vo
		// lang le, du cong thuc do tinh sai o dau.
		if ketQua.Truot {
			chiMang = false
		}

		// §1.2: sat thuong khong duoc VUOT mau con lai cua muc tieu.
		if satThuong > mucTieu.Hp {
			satThuong = mucTieu.Hp
		}
		// §1.3: quai DANG DAY MAU khong chet boi MOT don — giat lai 1 mau.
		if mucTieu.Hp == mucTieu.HpToiDa && satThuong >= mucTieu.Hp {
			satThuong = mucTieu.Hp - 1
		}

		mucTieu.Hp -= satThuong
		if mucTieu.Hp < 0 { // phong thu: ve ly thuyet khong xay ra sau hai kep tren
			mucTieu.Hp = 0
		}

		// dac ta §3.2: skillId la BYTE trong goi 54 xuong (short trong catalog).
		skillId := byte(nguoi.chieuDangChon.SkillId)
		k.phatChoCaKhuVuc(goiNguoiDanhQuai(int32(nguoi.nv.ID), skillId, mucTieu.Id))

		if mucTieu.Hp <= 0 {
			// Task 8: ghi lai thoi diem chet DE nhipTim tinh duoc "da troi qua
			// bao lau" cho nguong 5 giay hoi sinh (dac ta §5.2). Dat TRUOC khi
			// phat goi -12 — thu tu phat goi khong doi, chi la ghi trang thai
			// noi bo actor truoc, giong cach nguoi.lanDanhCuoi duoc ghi truoc
			// khi phat goi 54/-9 phia tren.
			mucTieu.LucChet = bayGio
			// Task 9: quai CHET -> xoa sach danh sach "da tung danh no" (dung
			// nguyen van brief: "xoa sach khi quai chet") — lan hoi sinh sau
			// la mot "vong doi" moi, khong ai duoc coi la da tung danh no nua.
			mucTieu.DaTungDanh = nil
			// KH6 chua co he vat pham (dac ta §0) -> doRoi luon rong.
			k.phatChoCaKhuVuc(goiQuaiChet(mucTieu.Id, satThuong, chiMang, nil))
		} else {
			// hieuUngId=255 (=-1): KH6 chua co hieu ung trang thai (dac ta §0).
			k.phatChoCaKhuVuc(goiQuaiConSong(mucTieu.Id, mucTieu.Hp, satThuong, chiMang, 255))
		}

	case suKienHoiSinh:
		// Ke hoach 6 Task 13 (dac ta §3.4/§6.4): client gui -16 xin hoi sinh
		// sau khi chet.
		nguoi, co := k.nguoiTrong[sk.idNhanVat]
		if !co {
			return
		}
		if nguoi.nv.HpHienTai > 0 {
			// CHUA chet (hoac da hoi sinh roi) -> bo qua im lang, tranh hoi
			// sinh "khong" mot nguoi dang song. Day CUNG la lop chan chong
			// gui -16 hai lan lien tiep (client bam nut kep, hoac goi den
			// tre): lan dau xu ly xong HpHienTai da > 0 nen lan hai tu nhien
			// roi vao nhanh nay, khong can them bo dem thoi gian rieng.
			return
		}

		bayGio := k.dongHo() // Task 8/13: dong ho tiem duoc, KHONG goi time.Now() thang.
		nguoi.nv.HpHienTai = nguoi.nv.HpToiDa
		nguoi.nv.ChakraHienTai = nguoi.nv.ChakraToiDa
		// Chot bat tu 1500ms sau hoi sinh (dac ta §1.7) — Task 11 (PvP) se
		// doc truong nay qua ViTriCoMau.LucHoiSinh khi mucTieu la nguoi choi.
		nguoi.lucHoiSinh = bayGio

		// -16 RIENG chinh chu TRUOC, roi -30/15 CA KHU VUC — dung thu tu
		// chuoi cua dac ta §6.4.
		if err := nguoi.phien.Gui(goiHoiSinhChinhChu()); err != nil {
			slog.Warn("gui -16 (xac nhan hoi sinh) that bai", "id", sk.idNhanVat, "loi", err)
		}
		k.phatChoCaKhuVuc(goiAiDoHoiSinh(
			int32(nguoi.nv.ID), nguoi.nv.HpHienTai, nguoi.nv.ChakraHienTai,
			int16(nguoi.nv.ViTriX), int16(nguoi.nv.ViTriY),
		))

	case suKienXinDanhSachQuai:
		// SAO GIA TRI tung con quai — KHONG duoc tra thang k.quai (mang con
		// tro goc) hay tung con tro trong do ra ngoai, vi ben nhan chay tren
		// luong khac va co the giu con tro do lau hon nhip actor tiep theo
		// (dac ta §1.6). traLoi co buffer 1 nen gui khong bao gio block du
		// ben goi da het gio cho va bo di.
		ban := make([]*Quai, len(k.quai))
		for i, q := range k.quai {
			cp := *q
			ban[i] = &cp
		}
		sk.traLoi <- ban
	}
}

// phatChoNguoiKhac gui goi tin `m` cho MOI nguoi trong khu vuc TRU nguoi co id `truId`.
func (k *KhuVuc) phatChoNguoiKhac(truId int64, m *protocol.Message) {
	for id, nguoi := range k.nguoiTrong {
		if id == truId {
			continue
		}
		if err := nguoi.phien.Gui(m); err != nil {
			slog.Warn("phat goi tin that bai", "id", id, "loi", err)
		}
	}
}

// phatChoCaKhuVuc gui goi tin `m` cho TOAN BO nguoi dang trong khu vuc, KE CA
// nguoi vua gay ra su kien (khac phatChoNguoiKhac o tren, loai tru nguoi gay
// su kien). Cac goi chien dau (54/-9/-12, dac ta §3.3 cot "Gui cho: Ca khu
// vuc") phai den ca nguoi danh de chinh ho cung thay hoat anh/ket qua don
// danh cua minh — khac han goi di chuyen (-7) noi client tu ve chuyen dong
// cua chinh minh nen khong can nhan lai.
func (k *KhuVuc) phatChoCaKhuVuc(m *protocol.Message) {
	for id, nguoi := range k.nguoiTrong {
		if err := nguoi.phien.Gui(m); err != nil {
			slog.Warn("phat goi tin (ca khu vuc) that bai", "id", id, "loi", err)
		}
	}
}

// Hai hang so hanh vi quai DA CHOT (dac ta §5.2, nguon Mob.java:347 va
// Mob.java:351-355; brief Task 8 "Hang so da chot") — nguon su that duy nhat
// cho hai con so nay, dung trong nhipTim duoi day. Dat kieu hang so (khong
// phai bien) de kiem thu dot bien (brief Step 5: "doi 5s thanh 50s") sua
// TRUC TIEP gia tri o day roi chay lai test bien, khong can sua noi khac.
const (
	// thoiGianHoiSinhQuai: quai chet roi cho DUNG 5 giay thi hoi sinh (dac ta
	// §5.2 "Hoi sinh sau khi chet — 5 giay").
	thoiGianHoiSinhQuai = 5 * time.Second

	// khoangHoiMauQuai: quai con song nhung chua day mau thi cu moi 30 giay
	// hoi them hpToiDa/10 (dac ta §5.2 "Hoi mau khi chua day — moi 30 giay,
	// hoi hpMax/10").
	khoangHoiMauQuai = 30 * time.Second
)

// nhipTim chay moi 100ms (goi tu vongLap) — Ke hoach 6 Task 8 dung mooc nay
// (de rong tu Ke hoach 5) cho ba viec cua tung con quai trong k.quai: HOI
// SINH neu da chet du 5 giay (hoiSinhQuai), HOI MAU dinh ky 30 giay neu con
// song nhung chua day mau (hoiMauQuai), va AI hai tang uu tien danh nguoi
// choi (aiQuai, Task 9). Day la NOI DUY NHAT ngoai case suKienDanhQuai cham
// vao trang thai quai/nguoiTrong — dung dac ta §1.6 (chi luong vongLap duoc
// cham trang thai quai va trang thai chien dau).
//
// bayGio LAY MOT LAN o dau ham (khong goi k.dongHo() lai trong vong lap) de
// TOAN BO quai trong cung mot nhip duoc xet tren CUNG mot moc thoi gian —
// tranh truong hop hai con quai o hai nhanh khac nhau cua cung mot lenh gio
// (vd mot con vua qua nguong o dau vong lap, con kia chua qua vi dongHo() da
// nhich len vai micro-giay giua hai lan goi — khong xay ra trong that te vi
// nhip 100ms qua tho, nhung day la thoi quen dung, tranh phu thuoc vao thu tu
// lap qua slice).
//
// Task 9: nhipTim da ~55 dong tu Task 8 (brief Task 9 de nghi tach ham khi
// them logic) nen ba viec tren duoc tach thanh ba ham rieng (hoiSinhQuai/
// hoiMauQuai/aiQuai) — chi vongLap goi chung, van chi luong nay cham trang
// thai (khong doi ky luat §1.6, chi doi cach to chuc code).
func (k *KhuVuc) nhipTim() {
	bayGio := k.dongHo()

	// Task 9: snapshot nguoi choi thanh SLICE MOT LAN cho ca nhip nay, dung
	// chung cho AI cua MOI con quai — tranh xay map->slice lai N lan (N = so
	// quai). Con tro TRO THANG vao k.nguoiTrong (KHONG sao gia tri) vi ham
	// nay chay TRONG vongLap nen doc truc tiep la an toan (dac ta §1.6); AI
	// (aiQuai) cung CHI DOC qua slice nay, khong ghi.
	var nguoi []*nguoiTrongKhuVuc
	if len(k.nguoiTrong) > 0 {
		nguoi = make([]*nguoiTrongKhuVuc, 0, len(k.nguoiTrong))
		for _, n := range k.nguoiTrong {
			nguoi = append(nguoi, n)
		}
	}

	for _, q := range k.quai {
		if q.Hp <= 0 {
			// Dang chet (cho hoi sinh, hoac vua hoi sinh xong NGAY trong
			// hoiSinhQuai() o duoi). Ca hai truong hop deu continue — quai
			// dang chet khong the hoi mau/AI, va quai vua hoi sinh chac
			// chan Hp==HpToiDa nen nhanh hoi mau se tu nhien khong lam gi;
			// AI cho quai vua hoi sinh doi sang tick SAU (khong danh ngay
			// trong chinh tick hoi sinh) — don gian hoa, khong dac ta nao
			// doi hoi phai danh NGAY trong cung tick hoi sinh.
			k.hoiSinhQuai(q, bayGio)
			continue
		}

		k.hoiMauQuai(q, bayGio)
		k.aiQuai(q, nguoi, bayGio)
	}
}

// hoiSinhQuai xu ly nhanh "quai dang cho hoi sinh" cua nhipTim — GOI KHI
// q.Hp <= 0. Neu CHUA du 5 giay tu LucChet thi khong lam gi (dac ta §5.2,
// brief Task 8 "sau 4.9s CHUA hoi sinh"); du/qua 5 giay thi hoi sinh DAY DU
// mau + phat DUNG MOT goi -13 (khong phat gi them sau do, dac ta §0.1 tuyet
// doi cam goi 87 theo sau). Tach nguyen ven logic Task 8, chi doi vi tri.
func (k *KhuVuc) hoiSinhQuai(q *Quai, bayGio time.Time) {
	if bayGio.Sub(q.LucChet) < thoiGianHoiSinhQuai {
		return
	}

	// Du/qua 5 giay -> hoi sinh DAY DU mau (dac ta §5.4 diem 5: ban Java hoi
	// day mau ngau nhien qua loi ham lvMob() khi dung goi -13 — ban Go LAM
	// VIEC NAY TUONG MINH o day, khong an sau ham doc). Ghi LucHoiSinh = bayGio
	// de kich hoat bat tu 1500ms (dac ta §1.7) qua ViTriCoMau.LucHoiSinh (xem
	// case suKienDanhQuai). Ghi lai LanHoiMauCuoi = bayGio de tranh hoi mau
	// "thua" ngay tick sau — quai vua hoi sinh da DAY MAU, khong can cong don
	// thoi gian tu TRUOC luc chet vao chu ky hoi mau moi.
	q.Hp = q.HpToiDa
	q.LucHoiSinh = bayGio
	q.LanHoiMauCuoi = bayGio

	// CHI phat DUNG MOT goi -13, KHONG phat gi them sau do — dac ta §0.1
	// tuyet doi cam goi 87 theo sau (client hieu nham thanh "quai bi gio
	// thoi", quai mau >= 16.777.216 se dung im khong di chuyen nua).
	// goiQuaiHoiSinh (goichiendau.go, Task 6) da tu no LUON ghi sys=0/
	// levelBoss=0, khong nhan tham so nao co the vo tinh mo duong cho mot
	// gia tri khac 0.
	k.phatChoCaKhuVuc(goiQuaiHoiSinh(q.Id, q.Hp))
}

// hoiMauQuai xu ly nhanh "hoi mau dinh ky" cua nhipTim — GOI KHI q.Hp > 0
// (nguoi goi, nhipTim, da loai truong hop Hp<=0 truoc do). Chi hoi mau khi
// CHUA day va da du 30 giay ke tu lan kiem truoc (dac ta §5.2). Quai vua day
// mau (Hp == HpToiDa) khong vao nhanh if -> LanHoiMauCuoi KHONG bi cap nhat
// khi dang day mau, dung y dinh "chi dem thoi gian tu lan KIEM GAN NHAT",
// khong phai tu luc spawn co dinh. Tach nguyen ven logic Task 8.
func (k *KhuVuc) hoiMauQuai(q *Quai, bayGio time.Time) {
	if q.Hp < q.HpToiDa && bayGio.Sub(q.LanHoiMauCuoi) >= khoangHoiMauQuai {
		q.LanHoiMauCuoi = bayGio

		q.Hp += q.HpToiDa / 10
		if q.Hp > q.HpToiDa { // khong duoc vuot hpToiDa (dac ta §5.2)
			q.Hp = q.HpToiDa
		}
		// KHONG phat goi gi o day: dac ta khong dinh nghia goi nao cho hoi
		// mau dinh ky (khac han goi -9 khi bi danh) — client chi thay mau
		// moi qua mot don danh (-9) hoac lan vao lai khu vuc (-24) sau nay,
		// dung hanh vi "hoi mau tham lang" cua ban goc.
	}
}

// aiQuai la AI hai tang uu tien cua Ke hoach 6 Task 9 (dac ta §5.2) — GOI KHI
// q.Hp > 0, MOI nhip 100ms, cho MOI con quai (goi tu nhipTim, sau hoiMauQuai).
// Dung q.chonMucTieu (ham THUAN, quai.go) de TIM muc tieu + nhip yeu cau; ham
// nay (chay trong vongLap) la noi DUY NHAT quyet dinh co "that su danh" hay
// khong va cham q.LanDanhCuoi — dung dac ta §5.4 diem 4 ("doi nhip danh nhu
// tac dung phu cua ham tim muc tieu la loi, phai tach rieng ham tra ve +
// nguoi goi tu quyet dinh") va brief Task 9 (dot bien bat buoc: "chuyen viec
// dat lai nhip ra ngoai nhanh that su danh -> test phai do").
func (k *KhuVuc) aiQuai(q *Quai, nguoi []*nguoiTrongKhuVuc, bayGio time.Time) {
	mucTieu, nhip := q.chonMucTieu(nguoi)
	if mucTieu == nil {
		// KHONG CO MUC TIEU: TUYET DOI khong dung toi q.LanDanhCuoi o day
		// (dac ta §5.4.2 — bay ban Java Mob.java:386-390 dat lai moc thoi
		// gian NGOAI nhanh "co nguoi choi", khien quai vua thay nguoi phai
		// cho THEM gan mot nhip nua moi duoc danh, thay vi danh NGAY).
		// Return som, KHONG cham trang thai gi ca.
		return
	}
	if bayGio.Sub(q.LanDanhCuoi) < nhip {
		return // con trong nhip hoi cua tang uu tien vua chon
	}

	// THAT SU DANH — CHI o day q.LanDanhCuoi moi duoc dat lai (dung nhanh
	// "that su danh" ma dot bien bat buoc cua brief doi hoi).
	q.LanDanhCuoi = bayGio

	// §1.2: sat thuong khong duoc vuot mau con lai cua nan nhan. SatThuongCoBan
	// la gia tri CO DINH tam thoi (xem chu thich truong Quai.SatThuongCoBan,
	// quai.go) — cong thuc day du §5.3 la viec cua KH8, ngoai pham vi Task 9.
	satThuong := q.SatThuongCoBan
	if satThuong < 0 { // phong thu: cau hinh hong (SatThuongCoBan am) khong lam mau am
		satThuong = 0
	}
	if satThuong > mucTieu.nv.HpHienTai {
		satThuong = mucTieu.nv.HpHienTai
	}

	mucTieu.nv.HpHienTai -= satThuong
	if mucTieu.nv.HpHienTai < 0 { // phong thu: ve ly thuyet khong xay ra sau kep tren
		mucTieu.nv.HpHienTai = 0
	}

	// -11: RIENG nan nhan (dac ta §3.3 cot "Gui cho": "Rieng nan nhan") —
	// KHONG dung phatChoNguoiKhac/phatChoCaKhuVuc (ca hai deu sai doi tuong:
	// mot loai tru nan nhan, mot gui thua cho nguoi khac). Gui THANG qua
	// phien cua chinh nan nhan, giong cach case suKienVao gui -5 rieng cho
	// tung nguoi cu (xem k.xuLy o tren).
	if err := mucTieu.phien.Gui(goiQuaiDanhTrungMinh(q.Id, satThuong)); err != nil {
		slog.Warn("gui -11 (quai danh trung minh) that bai", "idNanNhan", mucTieu.nv.ID, "loi", err)
	}
	// -10: NGUOI KHAC trong khu vuc, khong phai nan nhan (dac ta §3.3 cot
	// "Gui cho": "Nguoi khac trong khu vuc"). phatChoNguoiKhac da co san (Ke
	// hoach 5), dung dung nghia: loai tru id nan nhan.
	k.phatChoNguoiKhac(mucTieu.nv.ID, goiQuaiDanhTrungNguoiKhac(q.Id, int32(mucTieu.nv.ID), mucTieu.nv.HpHienTai))

	// Ke hoach 6 Task 13 (dac ta §3.4/§6.4): don nay VUA HA CHET nan nhan
	// (HpHienTai vua xuong <= 0 o kep tren) -> phat THEM chuoi goi "chet"
	// NGAY SAU -11/-10 (dam sat thuong truoc, roi moi bao chet — dung thu tu
	// §6.4). chonMucTieu (quai.go) da LOC san nguoi HpHienTai<=0 ra khoi CA
	// HAI tang AI (dac ta §5.4.3) nen KHONG bao gio quai nay (hay bat ky
	// quai nao khac cung tick) danh trung mot nguoi DA CHET tu truoc — su
	// kien "vua chet" o day CHAC CHAN la lan dau tien, khong can chot chong
	// phat lai.
	if mucTieu.nv.HpHienTai <= 0 {
		k.xuLyNguoiChoiChet(mucTieu)
	}
}

// xuLyNguoiChoiChet phat chuoi goi bao CHET cua mot nguoi choi (Ke hoach 6
// Task 13, dac ta §3.4/§6.4): -17 RIENG chinh chu, -8 cho NHUNG NGUOI KHAC
// trong khu vuc. GOI KHI nguoi.nv.HpHienTai VUA xuong <= 0 — ham nay KHONG
// tu kiem lai dieu kien do (nguoi goi, hien tai la aiQuai, chiu trach nhiem
// dam bao chi goi dung mot lan tai thoi diem chuyen tu song sang chet).
//
// cPk LUON 0: KH6 chua co he PK (dac ta §0 "Ngoai KH6": "PvP co luat -> KH9"),
// dung tinh than voi cPk=0 cua goiMeInfo (thegioi.go). x/y la VI TRI HIEN TAI
// — Java (PlayerService.hoiSinh) KHONG doi vi tri nguoi choi luc chet hay
// luc hoi sinh, chi phat hoat anh chet/song tai cho (dac ta khong nhac toi
// bat ky dich chuyen nao), nen ham nay doc THANG nv.ViTriX/Y, khong tinh
// toan gi them.
func (k *KhuVuc) xuLyNguoiChoiChet(nguoi *nguoiTrongKhuVuc) {
	x := int16(nguoi.nv.ViTriX)
	y := int16(nguoi.nv.ViTriY)

	if err := nguoi.phien.Gui(goiNguoiChet(0, x, y)); err != nil {
		slog.Warn("gui -17 (ban da chet) that bai", "id", nguoi.nv.ID, "loi", err)
	}
	k.phatChoNguoiKhac(nguoi.nv.ID, goiAiDoChet(int32(nguoi.nv.ID), 0, x, y))
}

// Gui day mot su kien vao hop thu cua actor. KHONG CHAN (non-blocking): neu
// hop thu day thi bo goi va ghi log debug; neu actor da Dung() thi cung bo
// goi an toan (khong panic, vi hopThu khong bao gio bi close — chi kenh dung
// bi close).
func (k *KhuVuc) Gui(sk suKienKhuVuc) {
	select {
	case k.hopThu <- sk:
	case <-k.dung:
		// actor da dung, bo su kien.
	default:
		slog.Debug("hop thu day, bo su kien", "mapId", k.banDo.MapID)
	}
}

// Dung dong actor mot cach sach: dong kenh dung (mot lan, an toan khi goi
// nhieu lan) roi cho vongLap thoat han truoc khi tra ve — dam bao khong ro
// goroutine.
func (k *KhuVuc) Dung() {
	k.dungMotLan.Do(func() {
		close(k.dung)
	})
	<-k.xong
}
