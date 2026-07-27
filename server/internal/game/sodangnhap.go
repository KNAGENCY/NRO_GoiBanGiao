package game

import (
	"log/slog"
	"sync"
	"time"
)

// SoDangNhap la SO PHIEN DUY NHAT toan may chu (Ke hoach 11 Task 7, dac ta
// 2026-07-27-luu-trang-thai-song.md §7.2): MOT map tai khoan ID -> phien
// (*BoXuLy) DANG SONG, bao ve boi mutex. Day la MOT NGUON SU THAT DUY NHAT —
// TUYET DOI KHONG lap lai kieu "ba lop chap va" cua ban Java cu (khoa theo
// ten, quet don theo dinh ky, va cuoi cung khoa han tai khoan khi phat hien
// trung — dau hieu cua viec khong co mot nguon su that duy nhat). "Tai khoan
// X dang o phien nao" CHI duoc tra loi o MOT CHO DUY NHAT trong ca du an: o
// day.
type SoDangNhap struct {
	mu    sync.Mutex
	phien map[int64]*BoXuLy
}

// hanGioChoPhienCuLuu la han gio TOI DA cho buoc "cho phien cu bao da luu
// xong" trong DangKy. KhiDong() cua phien cu (goi qua XepVaCho, thogi.go) da
// tu co han gio rieng hanGioGhi (5 giay) — han o day CHI la luoi an toan
// THEM (hoi lon hon hanGioGhi mot chut) de KHONG BAO GIO treo dang nhap cua
// phien MOI vo thoi han neu phien cu vi ly do nao do (vd goroutine cua no bi
// treo o mot cho khac ngoai du kien) khong bao gio dong duoc b.daXongDon —
// mot lan dang nhap "co the doc phai du lieu chua-han-moi-nhat" van con hon
// treo ca luong dang nhap khong co diem dung.
const hanGioChoPhienCuLuu = 6 * time.Second

// MoSoDangNhap tao mot so dang nhap RONG, san sang dung ngay — khong can
// buoc khoi dong rieng nao khac (khac ThoGhi/BoLuuDinhKy, SoDangNhap khong co
// goroutine nen cua rieng no, chi la mot map co khoa).
func MoSoDangNhap() *SoDangNhap {
	return &SoDangNhap{phien: make(map[int64]*BoXuLy)}
}

// DangKy dang ky "moi" la phien DANG SONG cho taiKhoanID. Goi TU
// xuLyDangNhap NGAY sau khi dang nhap thanh cong, TRUOC LayNhanVatTheoTaiKhoan
// (xem chu thich tai diem goi) — day la thu tu BAT BUOC cua dac ta §7.1: "da
// phien cu -> cho no luu xong -> ROI MOI cho phien moi doc du lieu".
//
// Neu CHUA co phien nao cho tai khoan nay (dang nhap binh thuong, khong
// trung), chi don gian dang ky roi tra ve NGAY — khong co gi de da, khong co
// gi de cho.
//
// Neu DA co mot phien KHAC (cu) dang giu tai khoan nay: DangKy
//  1. bao phien cu no VUA bi da (dong kenh biDa cua no — KhiDong() cua no,
//     khi chay, se doc co nay de quyet dinh ghi DONG BO XAC NHAN thay vi
//     Xep() khong chan binh thuong, xem KhiDong trong login.go),
//  2. dong ket noi vat ly cua phien cu (Session.Dong(), an toan goi tu bat
//     ky goroutine nao) — buoc ReadMessage() dang chan cua no tra ve loi,
//     kich hoat defer KhiDong() CHAY TREN CHINH goroutine cua phien cu (an
//     toan tuyet doi, KHONG cham truong nao khac cua BoXuLy tu ben ngoai —
//     day la ly do DUY NHAT SoDangNhap duoc phep "dung vao" mot BoXuLy khac
//     phien minh: chi qua hai thao tac an toan cross-goroutine nay, khong
//     bao gio doc/ghi truc tiep cac truong con lai),
//  3. CHO phien cu bao "da xong don" (KhiDong() da chay xong hoan toan, ke
//     ca buoc ghi dong bo neu co) — co han gio hanGioChoPhienCuLuu de khong
//     treo vo han.
func (s *SoDangNhap) DangKy(taiKhoanID int64, moi *BoXuLy) {
	s.mu.Lock()
	cu, coCu := s.phien[taiKhoanID]
	s.phien[taiKhoanID] = moi
	s.mu.Unlock()

	if !coCu || cu == moi {
		return
	}

	slog.Info("dang nhap trung tai khoan, da phien cu va cho luu xong",
		"taiKhoanID", taiKhoanID)

	close(cu.biDa)
	if cu.session != nil {
		cu.session.Dong()
	}

	select {
	case <-cu.daXongDon:
	case <-time.After(hanGioChoPhienCuLuu):
		slog.Warn("cho phien cu bao da luu xong vuot han gio — phien moi co the doc phai du lieu CHUA phai moi nhat",
			"taiKhoanID", taiKhoanID, "hanGio", hanGioChoPhienCuLuu)
	}
}

// XoaNeuDung go dang ky cua taiKhoanID KHOI so — NHUNG CHI NEU so dang TRO
// DUNG phien b nay. Goi tu KhiDong() (login.go) khi mot phien dong, bat ke
// dong tu nhien hay bi da.
//
// Bao ve khoi tinh huong BAT BUOC PHAI CHAN (canh bao rieng cua brief Task
// 7): phien CU bi da chay KhiDong() cua no SAU KHI phien MOI da dang ky vao
// so (DangKy ghi de s.phien[id] = moi TRUOC khi da phien cu, xem DangKy o
// tren) — neu XoaNeuDung khong kiem "dung phien" ma xoa vo dieu kien, no se
// XOA MAT dang ky cua phien MOI (con dang song binh thuong), lam tai khoan
// nay tro thanh "khong co ai giu" trong so du that su co mot phien dang
// choi — mo lai chinh lo hong ma Task 7 phai chan.
func (s *SoDangNhap) XoaNeuDung(taiKhoanID int64, b *BoXuLy) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if hienTai, ok := s.phien[taiKhoanID]; ok && hienTai == b {
		delete(s.phien, taiKhoanID)
	}
}
