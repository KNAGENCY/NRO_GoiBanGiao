package game

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"runtime/debug"
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
// HANG DOI HAI TANG, va day la thu phan biet no voi ThoGhi:
//
//	hangDoi   — khoan MOI, gop CONG DON theo (nhan vat, loai)
//	choThuLai — cac LO da dong goi ma ghi chua duoc, thu lai NGUYEN VEN
//
// Vi sao phai tach hai tang thay vi tron lai lam mot (ban dau lam vay, va do la
// mot lo hong): PostgreSQL co the DA COMMIT THAT roi moi tra loi (het han gio,
// dut ket noi). Lan thu lai se chen lai — chan nhan doi bang khoa lo +
// ON CONFLICT (xem store.GhiNoThuong) CHI chay duoc neu lan thu lai gui DUNG
// BANG lo cu. Tron them khoan moi vao la khoa lo cu mang so diem moi, ON CONFLICT
// bo qua, va phan moi bien mat.
//
// KHAC ThoGhi mot diem nua: hang doi khoan moi CONG DON chu khong GHI DE.
// ThoGhi ghi de vi "ban chup moi luon tot hon ban cu"; o day moi khoan la mot
// mon tien rieng — ghi de la lam bien mat phan thuong cua nguoi choi.
type SoNoThuong struct {
	kho KhoNo

	// nhipGom la khoang cach giua hai lan xa hang doi. Co y KHONG phan ung tuc
	// thi voi moi lan XepNo (khac ThoGhi, co kenh coViec): mot khoan no khong
	// gap — nguoi nhan dang offline — nen gom vai giay roi ghi mot the vua giam
	// so lan khu hoi CSDL, vua gop nhieu khoan cua cung mot nguoi thanh MOT
	// dong. Tiem duoc de test khong phai doan nhip goroutine nen.
	nhipGom time.Duration

	// maPhien la dinh danh NGAU NHIEN sinh mot lan luc mo so no, dung lam tien
	// to cho moi khoa lo. Chi dung bo dem thoi thi KHONG DU: may chu khoi dong
	// lai se dat lai bo dem ve 0, sinh ra khoa trung voi lo cu con nam trong
	// bang — va ON CONFLICT DO NOTHING se BO QUA nham mot khoan no that.
	maPhien string

	mu        sync.Mutex
	hangDoi   map[khoaNo]int64 // khoan MOI: (nhan vat, loai) -> tong diem CONG DON
	choThuLai []*loNo          // lo da dong goi, ghi chua duoc — thu lai NGUYEN VEN
	boDemLo   uint64           // sinh khoa lo, tang dan trong mot phien

	dung       chan struct{}
	xong       chan struct{}
	dungMotLan sync.Once

	// loiXaCuoi giu loi cua lan xa CUOI CUNG (luc dung). DongVaChoXong tra loi
	// nay ra ngoai: nhung khoan no khong ghi duoc o lan xa cuoi la MAT THAT
	// (goroutine chet ngay sau do), nen bao "tat sach" luc do la noi doi.
	loiXaCuoi error

	// soLoiGhi / soLoiLienTiep: cung y nghia va cung ly do voi hai bo dem cung
	// ten trong ThoGhi — loi ghi KHONG duoc nuot im lang (dac ta §6.3, bai hoc
	// PlayerDAO.updatePlayer ban Java).
	soLoiGhi      atomic.Int64
	soLoiLienTiep atomic.Int64
}

// khoaNo la khoa gop cua hang doi khoan moi: hai khoan thuong CUNG nguoi CUNG
// loai duoc cong lai lam mot dong. Khac loai thi khong gop — moi loai co duong
// tra rieng.
type khoaNo struct {
	NhanVatID int64
	Loai      string
}

// loNo la mot LO da dong goi: noi dung DONG BANG (khong nhan them khoan moi) va
// mot khoa duy nhat de CSDL nhan ra khi lo do bi ghi lai.
type loNo struct {
	khoa     string
	muc      map[khoaNo]int64
	soLanThu int
}

// danhSach doi lo thanh dang ma kho nhan.
func (l *loNo) danhSach() []store.NoThuong {
	ds := make([]store.NoThuong, 0, len(l.muc))
	for k, diem := range l.muc {
		ds = append(ds, store.NoThuong{
			NhanVatID: k.NhanVatID,
			Loai:      k.Loai,
			SoDiem:    diem,
			KhoaLo:    l.khoa,
		})
	}
	return ds
}

// KhoNo la phan CSDL toi thieu ma SoNoThuong can — tach thanh interface HEP
// dung khuon voi KhoLuu (thogi.go) va Kho (login.go), de kiem thu bang kho gia
// khong can CSDL that. *store.Store thoa man tu nhien.
type KhoNo interface {
	GhiNoThuong(ctx context.Context, ds []store.NoThuong) error
}

// NhipGomNoMacDinh la nhip xa hang doi dung khi chay that (MoLuuTru). Chon 2
// giay: du dai de gop nhieu khoan cua mot dot boss lam mot, du ngan de mot cu
// sap may chu chi lam mat vai giay no.
const NhipGomNoMacDinh = 2 * time.Second

// soLanThuLoToiDa la so lan mot LO duoc thu ghi truoc khi bi BO va log muc
// Error — VAN XA, va no bat buoc phai co.
//
// Khong co van xa thi mot dong rac DUY NHAT (vd nhan vat da bi xoa nen vi pham
// khoa ngoai) lam ca lo bi tu choi, quay vong mai mai, va tu do KHONG khoan no
// nao cua BAT KY AI xuong duoc dia nua — hang doi phinh tren RAM cho toi luc
// tat may chu roi mat sach. Mot lo hong lam ket ca may chu la cai gia dat hon
// nhieu so voi mat mot lo.
//
// Chon 5: voi nhip 2 giay that la ~10 giay co gang. Du de di qua mot cu chap
// choang cua CSDL (khoi dong lai, chuyen master), khong du de om mai mot dong
// rac vinh vien hong.
const soLanThuLoToiDa = 5

// nguongHangDoiNoCanhBao la so dong trong hang doi tro len thi log canh bao.
// Hang doi bi chan tu nhien theo SO NGUOI (gop theo nguoi), nen vuot nguong
// nay gan nhu chac chan nghia la CSDL dang hong lien tuc va no khong xuong
// duoc — thu can biet som chu khong phai luc mat dien.
const nguongHangDoiNoCanhBao = 500

// MoSoNoThuong khoi dong goroutine ghi so no. nhipGom < 1ms bi ep ve
// NhipGomNoMacDinh — khong phai phong thu thua: time.NewTicker(0) PANIC, va no
// panic trong goroutine nen, tuc la sap ca tien trinh may chu ngay luc khoi dong.
func MoSoNoThuong(kho KhoNo, nhipGom time.Duration) *SoNoThuong {
	if nhipGom < time.Millisecond {
		nhipGom = NhipGomNoMacDinh
	}
	s := &SoNoThuong{
		kho:     kho,
		nhipGom: nhipGom,
		maPhien: maPhienNgauNhien(),
		hangDoi: make(map[khoaNo]int64),
		dung:    make(chan struct{}),
		xong:    make(chan struct{}),
	}
	go s.vongLap()
	return s
}

// maPhienNgauNhien sinh 8 byte ngau nhien dang hex. Rot ve mot chuoi co dinh
// neu bo sinh ngau nhien cua he dieu hanh hong — luc do khoa lo khong con duy
// nhat GIUA CAC LAN KHOI DONG (chi con duy nhat trong mot phien), nhung mot so
// no chay duoc voi mot rui ro hep van hon la may chu khong khoi dong duoc.
func maPhienNgauNhien() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		slog.Error("khong sinh duoc ma phien ngau nhien cho so no — khoa lo chi con duy nhat trong mot phien",
			"loi", err)
		return "sono"
	}
	return hex.EncodeToString(b)
}

// XepNo dua mot khoan no vao hang doi — KHONG BAO GIO CHAN. Goi duoc tu trong
// actor khu vuc.
//
// Bo qua so diem <= 0 VA loai rong: CSDL co rang buoc so_diem > 0 va
// loai_thuong IN (...), nen mot dong rac se lam CA LO bi tu choi va bi thu lai,
// keo theo ca nhung khoan no THAT cua nguoi khac trong cung lo. Chan ngay tu
// day re hon nhieu.
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
//
// Tra LOI khi lan xa cuoi cung that bai — khac ThoGhi. Ly do: sau lan xa cuoi,
// goroutine ghi chet han, nen nhung gi khong ghi duoc luc do la MAT THAT. Bao
// nil o day la noi doi voi TatMayChu, va nguoi van hanh se khong bao gio biet.
func (s *SoNoThuong) DongVaChoXong(ctx context.Context) error {
	s.dungMotLan.Do(func() {
		close(s.dung)
	})
	select {
	case <-s.xong:
		s.mu.Lock()
		defer s.mu.Unlock()
		return s.loiXaCuoi
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
			antoan.ChayCoLuoi("sono.xa", func() { s.xaMotDot(false) })
		case <-s.dung:
			antoan.ChayCoLuoi("sono.xa", func() { s.xaMotDot(true) })
			return
		}
	}
}

// xaMotDot ghi het cac lo dang cho: cac lo thu lai TRUOC (chung da doi lau
// hon), roi toi lo moi dong goi tu hang doi khoan moi.
func (s *SoNoThuong) xaMotDot(cuoiCung bool) {
	for _, lo := range s.layCacLoCanGhi() {
		s.ghiMotLo(lo, cuoiCung)
	}
}

// layCacLoCanGhi khoa mutex MOT LAN: doat danh sach lo cho thu lai, va dong goi
// hang doi khoan moi thanh mot lo moi co khoa rieng. Sau khi nha khoa, XepNo goi
// trong luc dang ghi se roi vao map MOI, khong dung cham voi lo dang xu ly.
func (s *SoNoThuong) layCacLoCanGhi() []*loNo {
	s.mu.Lock()
	defer s.mu.Unlock()

	ds := s.choThuLai
	s.choThuLai = nil

	if len(s.hangDoi) > 0 {
		s.boDemLo++
		ds = append(ds, &loNo{
			khoa: fmt.Sprintf("%s-%d", s.maPhien, s.boDemLo),
			muc:  s.hangDoi,
		})
		s.hangDoi = make(map[khoaNo]int64)
	}
	return ds
}

// ghiMotLo thu ghi MOT lo. Lo hong duoc tra ve hang cho thu lai NGUYEN VEN
// (cung khoa, cung so diem) cho toi khi qua soLanThuLoToiDa thi bi BO.
func (s *SoNoThuong) ghiMotLo(lo *loNo, cuoiCung bool) {
	loi := s.thuGhiMotLo(lo)
	if loi == nil {
		s.soLoiLienTiep.Store(0)
		return
	}

	soLoi := s.soLoiGhi.Add(1)
	lienTiep := s.soLoiLienTiep.Add(1)
	lo.soLanThu++

	if cuoiCung {
		// Lan xa cuoi: khong con nhip nao de thu lai nua, nen day la mat that.
		s.mu.Lock()
		s.loiXaCuoi = loi
		s.mu.Unlock()
		slog.Error("lan xa so no CUOI CUNG that bai — nhung khoan no nay MAT",
			"loi", loi, "soDong", len(lo.muc), "khoaLo", lo.khoa)
		return
	}

	if lo.soLanThu >= soLanThuLoToiDa {
		// VAN XA: bo lo nay di. Giu lai la ket ca duong ghi no cua toan may chu
		// (xem chu thich soLanThuLoToiDa).
		slog.Error("BO mot lo so no sau qua nhieu lan ghi that bai — nhung khoan nay MAT",
			"loi", loi, "soLanThu", lo.soLanThu, "soDong", len(lo.muc),
			"khoaLo", lo.khoa, "muc", lo.muc)
		return
	}

	s.mu.Lock()
	s.choThuLai = append(s.choThuLai, lo)
	s.mu.Unlock()

	if lienTiep >= nguongLoiLienTiepCanhBaoNang {
		slog.Error("ghi so no that bai nhieu lan lien tiep — thuong dang un lai trong bo nho",
			"loi", loi, "soDong", len(lo.muc), "soLoiGhi", soLoi, "lienTiep", lienTiep)
	} else {
		slog.Warn("ghi so no that bai, se thu lai",
			"loi", loi, "soDong", len(lo.muc), "soLoiGhi", soLoi)
	}
}

// thuGhiMotLo goi kho MOT lan, va BOC PANIC NGAY TAI DAY — day la diem quan
// trong nhat cua ca file, tim ra bang kiem thu dot bien.
//
// Luoi o vongLap ("sono.xa") KHONG DU: luc panic bay toi do thi lo da bi
// layCacLoCanGhi doat khoi hang doi roi, va khong con ai tra no ve — ca lo no
// bien mat vinh vien, im lang, dung nguoc voi dieu chu thich cua ham nay tung
// cam ket ("du lieu van nam trong bo nho nen khong mat").
//
// ThoGhi khong dinh loi nay vi no co luoi HAI TANG: "thogi.xa" boc vong lap VA
// "thogi.ghi" boc tung ban ghi ben trong. Doi panic thanh mot loi binh thuong o
// day cho lo di dung duong thu lai nhu moi loi khac.
func (s *SoNoThuong) thuGhiMotLo(lo *loNo) (loi error) {
	defer func() {
		if r := recover(); r != nil {
			loi = fmt.Errorf("ghi so no panic: %v", r)
			slog.Error("ghi so no PANIC — lo duoc giu lai de thu lai",
				"hoang_so", r, "khoaLo", lo.khoa, "soDong", len(lo.muc),
				"stack", string(debug.Stack()))
		}
	}()

	ctx, huy := context.WithTimeout(context.Background(), hanGioGhi)
	defer huy()
	return s.kho.GhiNoThuong(ctx, lo.danhSach())
}

// congBaoHoaNo cong hai so khong am, dung o MaxInt64 thay vi cuon vong ve am.
// Mot khoan no AM se bi CSDL tu choi (rang buoc so_diem > 0) va keo ca lo —
// gom ca no that cua nguoi khac — vao vong thu lai roi bi bo.
func congBaoHoaNo(a, b int64) int64 {
	const tran int64 = 9_223_372_036_854_775_807
	if b > 0 && a > tran-b {
		return tran
	}
	return a + b
}
