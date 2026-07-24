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

// MoKhuVuc tao mot khu vuc giu ban do da cho, chua co nguoi choi nao.
func MoKhuVuc(bd *BanDo) *KhuVuc {
	return &KhuVuc{
		banDo:      bd,
		nguoiTrong: make(map[int64]*nguoiTrongKhuVuc),
	}
}

// BanDo tra ban do cua khu vuc (client can no de dung MAP INFO).
func (k *KhuVuc) BanDo() *BanDo { return k.banDo }

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

// nhipTim chay moi 100ms. Ke hoach 5 chua can dung den — day la CHO SAN cho
// Ke hoach 6 (AI quai, hoi phuc vat pham, ...).
func (k *KhuVuc) nhipTim() {
	// KH6: cap nhat quai/vat pham dinh ky se dat o day.
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
