package game

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
func MoLuuTru(kho KhoDayDu, soKhuVuc *SoKhuVuc, soThoGhi int) (thoGhi *ThoGhi, boLuu *BoLuuDinhKy, taoXuLy func() *BoXuLy) {
	thoGhi = MoThoGhi(kho, soThoGhi)
	boLuu = MoBoLuuDinhKy(soKhuVuc, thoGhi)
	boLuu.BatDau()
	taoXuLy = func() *BoXuLy {
		b := NewBoXuLy(kho, soKhuVuc)
		b.DatThoGhi(thoGhi)
		return b
	}
	return thoGhi, boLuu, taoXuLy
}
