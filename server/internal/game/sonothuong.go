package game

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"nroserver/internal/antoan"
	"nroserver/internal/store"
)

// SoNoThuong la duong DUY NHAT actor khu vuc duoc dung de ghi mot khoan thuong
// cho nguoi KHONG con o trong khu vuc (da roi map, hoac mat ket noi dung luc
// boss chet).
//
// VI SAO PHAI CO NO — chu khong phai ghi thang: actor chi so huu con tro
// store.NhanVat cua nhung nguoi DANG o trong no. Voi nguoi da roi, ghi vao con
// tro do la tranh chap du lieu that (§6.9 BAT-DAU-TU-DAY.md), con day mot ban
// chup cu qua ThoGhi thi con te hon — ThoGhi GHI DE CA DONG, keo lui moi thu
// nguoi do lam tu luc roi di. Nen phan thuong duoc ghi thanh MOT KHOAN NO ben
// (bang thuong_cho_nhan) va tra o lan dang nhap sau (login.go).
//
// LUAT SO MOT giong het ThoGhi (thogi.go): XepNo() CHI cham hang doi trong bo
// nho (mot map duoc khoa mutex vai chuc nano-giay), TUYET DOI khong goi CSDL.
// Viec cham dia xay ra trong goroutine vongLap rieng ben duoi.
//
// KHAC ThoGhi DUNG MOT DIEM, va day la diem quan trong nhat cua ca file:
// hang doi nay CONG DON chu khong GHI DE. ThoGhi ghi de vi "ban chup moi nhat
// luon tot hon ban cu"; o day moi khoan no la mot khoan TIEN RIENG — ghi de
// nghia la lam bien mat phan thuong cua nguoi choi.
type SoNoThuong struct {
	kho KhoNo

	// nhipGom la khoang cach giua hai lan xa hang doi. Co y KHONG phan ung tuc
	// thi voi moi lan XepNo (khac ThoGhi, co kenh coViec): mot khoan no khong
	// gap — nguoi nhan dang offline — nen gom vai giay roi ghi mot the vua giam
	// so lan khu hoi CSDL, vua gop nhieu khoan cua cung mot nguoi thanh MOT
	// dong. Tiem duoc de test khong phai doan nhip goroutine nen.
	nhipGom time.Duration

	mu      sync.Mutex
	hangDoi map[khoaNo]int64 // (nhan vat, loai) -> tong diem CONG DON

	dung       chan struct{}
	xong       chan struct{}
	dungMotLan sync.Once

	// soLoiGhi / soLoiLienTiep: cung y nghia va cung ly do voi hai bo dem cung
	// ten trong ThoGhi — loi ghi KHONG duoc nuot im lang (dac ta §6.3, bai hoc
	// PlayerDAO.updatePlayer ban Java).
	soLoiGhi      atomic.Int64
	soLoiLienTiep atomic.Int64
}

// khoaNo la khoa gop cua hang doi: hai khoan thuong CUNG nguoi CUNG loai duoc
// cong lai lam mot dong. Khac loai thi khong gop — moi loai co duong tra rieng.
type khoaNo struct {
	NhanVatID int64
	Loai      string
}

// KhoNo la phan CSDL toi thieu ma SoNoThuong can — tach thanh interface HEP
// dung khuon voi KhoLuu (thogi.go) va Kho (login.go), de kiem thu bang kho gia
// khong can CSDL that. *store.Store thoa man tu nhien.
type KhoNo interface {
	GhiNoThuong(ctx context.Context, ds []store.NoThuong) error
}

// NhipGomNoMacDinh la nhip xa hang doi dung khi chay that (main.go). Chon 2
// giay: du dai de gop nhieu khoan cua mot dot boss lam mot, du ngan de mot cu
// sap may chu chi lam mat vai giay no — dung muc rui ro da chot cho BoLuuDinhKy
// (60 giay) nhung chat hon vi day la duong DUY NHAT giu khoan thuong nay.
const NhipGomNoMacDinh = 2 * time.Second

// nguongHangDoiNoCanhBao la so dong trong hang doi tro len thi log canh bao.
// Hang doi bi chan tu nhien theo SO NGUOI (gop theo nguoi), nen vuot nguong
// nay gan nhu chac chan nghia la CSDL dang hong lien tuc va no khong xuong
// duoc — thu can biet som chu khong phai luc mat dien.
const nguongHangDoiNoCanhBao = 500

// MoSoNoThuong khoi dong goroutine ghi so no. nhipGom < 1ms bi ep ve
// NhipGomNoMacDinh (phong thu goi nham lam vong lap quay dien).
func MoSoNoThuong(kho KhoNo, nhipGom time.Duration) *SoNoThuong {
	if nhipGom < time.Millisecond {
		nhipGom = NhipGomNoMacDinh
	}
	s := &SoNoThuong{
		kho:     kho,
		nhipGom: nhipGom,
		hangDoi: make(map[khoaNo]int64),
		dung:    make(chan struct{}),
		xong:    make(chan struct{}),
	}
	go s.vongLap()
	return s
}

// XepNo dua mot khoan no vao hang doi — KHONG BAO GIO CHAN. Goi duoc tu trong
// actor khu vuc.
//
// Bo qua so diem <= 0: CSDL co rang buoc so_diem > 0, nen mot dong rac se lam
// CA LO bi tu choi (xem GhiNoThuong) va bi thu lai vo han, keo theo ca nhung
// khoan no THAT cua nguoi khac trong cung lo. Chan ngay tu day re hon nhieu.
//
// Bo qua loai rong: cung ly do — rang buoc thuong_cho_nhan_loai_hop_le se tu
// choi ca lo.
func (s *SoNoThuong) XepNo(no store.NoThuong) {
	if no.SoDiem <= 0 || no.Loai == "" {
		return
	}
	k := khoaNo{NhanVatID: no.NhanVatID, Loai: no.Loai}

	s.mu.Lock()
	s.hangDoi[k] = congBaoHoaNo(s.hangDoi[k], no.SoDiem)
	soDong := len(s.hangDoi)
	s.mu.Unlock()

	if soDong >= nguongHangDoiNoCanhBao {
		slog.Warn("hang doi so no phinh to — CSDL co the dang hong",
			"soDong", soDong, "soLoiGhi", s.soLoiGhi.Load())
	}
}

// DongVaChoXong bao goroutine ghi dung nhan viec, GHI NOT nhung gi con trong
// hang doi, roi moi tra ve — dung khuon (va dung ly do) voi ThoGhi.DongVaChoXong:
// luc tat may chu, khoan no vua xep phai kip xuong dia. ctx la can thoi gian
// cho, de mot CSDL treo that su khong giu tien trinh lai vo han.
func (s *SoNoThuong) DongVaChoXong(ctx context.Context) error {
	s.dungMotLan.Do(func() {
		close(s.dung)
	})
	select {
	case <-s.xong:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// vongLap la than cua so no: xa theo nhip, hoac xa mot lan cuoi khi co lenh
// dung. Luoi panic boc THAN moi vong (khong boc ngoai) — cung ly do voi
// ThoGhi.vongLap: boc ngoai thi mot cu panic giet ca duong ghi no.
func (s *SoNoThuong) vongLap() {
	ticker := time.NewTicker(s.nhipGom)
	defer ticker.Stop()
	defer close(s.xong)

	for {
		select {
		case <-ticker.C:
			antoan.ChayCoLuoi("sono.xa", s.ghiHetHangDoi)
		case <-s.dung:
			antoan.ChayCoLuoi("sono.xa", s.ghiHetHangDoi)
			return
		}
	}
}

// ghiHetHangDoi doat toan bo hang doi (khoa mutex mot lan, thay bang map rong
// moi — nen XepNo goi trong luc dang ghi roi vao map MOI, khong dung cham),
// roi ghi CA LO bang mot lan goi CSDL.
//
// Ghi lo hong thi XEP LAI TOAN BO lo vao hang doi de thu o nhip sau — cong don
// vao nhung gi vua xep them, khong ghi de. Du lieu van nam trong bo nho nen
// khong mat; cai gia la mot cu sap may chu giua chung se mat dung so no chua
// kip xuong (toi da mot nhip gom).
func (s *SoNoThuong) ghiHetHangDoi() {
	s.mu.Lock()
	if len(s.hangDoi) == 0 {
		s.mu.Unlock()
		return
	}
	lo := s.hangDoi
	s.hangDoi = make(map[khoaNo]int64)
	s.mu.Unlock()

	ds := make([]store.NoThuong, 0, len(lo))
	for k, diem := range lo {
		ds = append(ds, store.NoThuong{NhanVatID: k.NhanVatID, Loai: k.Loai, SoDiem: diem})
	}

	ctx, huy := context.WithTimeout(context.Background(), hanGioGhi)
	defer huy()
	if err := s.kho.GhiNoThuong(ctx, ds); err != nil {
		s.xepLai(lo)
		soLoi := s.soLoiGhi.Add(1)
		lienTiep := s.soLoiLienTiep.Add(1)
		if lienTiep >= nguongLoiLienTiepCanhBaoNang {
			slog.Error("ghi so no that bai nhieu lan lien tiep — thuong dang un lai trong bo nho",
				"loi", err, "soDong", len(ds), "soLoiGhi", soLoi, "lienTiep", lienTiep)
		} else {
			slog.Warn("ghi so no that bai, se thu lai",
				"loi", err, "soDong", len(ds), "soLoiGhi", soLoi)
		}
		return
	}
	s.soLoiLienTiep.Store(0)
}

// xepLai tra mot lo ghi hong ve hang doi, CONG DON vao nhung gi da xep them
// trong luc ghi (khong ghi de — day chinh la ly do hang doi phai cong don).
func (s *SoNoThuong) xepLai(lo map[khoaNo]int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, diem := range lo {
		s.hangDoi[k] = congBaoHoaNo(s.hangDoi[k], diem)
	}
}

// congBaoHoaNo cong hai so khong am, dung o MaxInt64 thay vi cuon vong ve am.
// Mot khoan no AM se bi CSDL tu choi (rang buoc so_diem > 0) va keo ca lo —
// gom ca no that cua nguoi khac — vao vong thu lai vo tan.
func congBaoHoaNo(a, b int64) int64 {
	const tran int64 = 9_223_372_036_854_775_807
	if b > 0 && a > tran-b {
		return tran
	}
	return a + b
}
