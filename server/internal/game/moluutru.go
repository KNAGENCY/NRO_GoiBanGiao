package game

import (
	"time"

	"nroserver/internal/chanip"
)

// KhoDayDu la giao dien CSDL DAY DU ma may chu can luc khoi dong: bon ham tra
// cuu/tao nhan vat (Kho, login.go) CONG mot ham ghi lai toan bo trang thai
// (KhoLuu, thogi.go). *store.Store (kho that) thoa man ca hai tu nhien,
// khong can glue code them — tach rieng interface nay CHI de MoLuuTru (duoi
// day) nhan MOT tham so kho DUY NHAT thay vi truyen hai tham so kho trung
// nhau (mot cho Kho, mot cho KhoLuu) vao cung mot loi goi.
type KhoDayDu interface {
	Kho
	KhoLuu
}

// MoLuuTru lap rap TOAN BO duong day luu tru cua may chu (Ke hoach 11 Task 6
// diem B "NOI DAY THAT", dac ta 2026-07-27-luu-trang-thai-song.md §4): mot
// ThoGhi CHUNG cho ca may chu, mot BoLuuDinhKy DA BATDAU (tu dong luu moi
// nhipLuuDinhKy — diem moc A), va mot ham taoXuLy tra ve MOT BoXuLy MOI da
// duoc DatThoGhi (diem moc B) cho MOI ket noi.
//
// Ke hoach 11 Task 7 (dac ta §7.2 "mot tai khoan, mot phien"): taoXuLy CUNG
// DatSoDangNhap mot *SoDangNhap CHUNG (tao MOI LAN goi MoLuuTru — dung MOT
// so cho ca may chu, dung khuon voi thoGhi/boLuu o tren) cho MOI BoXuLy —
// day la "MOT cho duy nhat biet tai khoan X dang o phien nao" ma dac ta doi
// hoi. Dat trong ham nay (khong phai rieng trong main.go) vi TRUOC Task 6,
// bai hoc da ro: bat cu manh ha tang nao khong duoc NOI DAY o day se nam do
// khong ai goi toi — MoLuuTru la diem duy nhat lap rap moi BoXuLy san xuat.
//
// TRUOC Task 6, KHONG NOI NAO trong ca du an goi MoThoGhi / MoBoLuuDinhKy(...)
// .BatDau() / DatThoGhi — nghia la toan bo ha tang Task 1-5 (ThoGhi,
// LuuTatCa, BoLuuDinhKy, DatThoGhi) NAM DO KHONG AI GOI TOI, may chu THAT SU
// KHONG LUU GI CA du ha tang da xay xong het. Ham nay la noi NOI DAY that:
// goi tu main.go (cmd/nroserver) DUNG MOT LAN luc khoi dong may chu.
//
// Tach thanh ham RIENG (khong nhet thang vao main()) de TEST kiem duoc "da
// noi du ca 3 diem" TRUC TIEP tu package nay — package game truy cap duoc
// truong noi bo cua ThoGhi/BoLuuDinhKy/BoXuLy (xem TestMoLuuTruNoiDuBaDiem)
// — ma khong can gia lap toan bo giao thuc mang qua mot ket noi TCP that.
//
// soThoGhi: xem MoThoGhi (thogi.go) — so goroutine ghi CSDL song song toi da
// trong MOI dot xa hang doi. Quyet dinh chon N bao nhieu la cua tang goi
// (main.go), ham nay chi chuyen tiep thang cho MoThoGhi.
// No bao mat §7 (lo hong 1 — chan do mat khau): MoLuuTru cung tao MOT
// GacDangNhap CHUNG cho ca may chu va gan vao MOI BoXuLy. Phai la MOT cai
// chung, khong phai moi phien mot cai: ke do mat khau mo ket noi MOI cho moi
// lan thu, nen mot bo dem nam trong tung phien se luon bat dau lai tu 0 va
// khong bao gio cham nguong — dung nghia "co ma nhu khong".
//
// Dat o day, giong het thoGhi/boLuu/soDangNhap, vi bai hoc Task 6 van con
// nguyen gia tri: ha tang khong duoc noi day o MoLuuTru la ha tang nam do
// khong ai goi toi, test tu che van xanh ma may chu that van tho.
//
// sucChua / nhipHoiPhuc: xem chanip.SucChuaMacDinh va chanip.NhipHoiPhucMacDinh
// de biet vi sao chon nhung con so do. Truyen 0 cho ca hai la lay mac dinh.
func MoLuuTru(kho KhoDayDu, soKhuVuc *SoKhuVuc, soThoGhi int, sucChua int, nhipHoiPhuc time.Duration) (thoGhi *ThoGhi, boLuu *BoLuuDinhKy, taoXuLy func() *BoXuLy) {
	thoGhi = MoThoGhi(kho, soThoGhi)
	boLuu = MoBoLuuDinhKy(soKhuVuc, thoGhi)
	boLuu.BatDau()
	soDangNhap := MoSoDangNhap()
	gacDangNhap := chanip.MoGacDangNhap(sucChua, nhipHoiPhuc)
	taoXuLy = func() *BoXuLy {
		b := NewBoXuLy(kho, soKhuVuc)
		b.DatThoGhi(thoGhi)
		b.DatSoDangNhap(soDangNhap)
		b.DatGacDangNhap(gacDangNhap)
		return b
	}
	return thoGhi, boLuu, taoXuLy
}
