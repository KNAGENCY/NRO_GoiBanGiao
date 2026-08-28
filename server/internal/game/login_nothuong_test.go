package game

import (
	"context"
	"errors"
	"testing"

	"nroserver/internal/store"
)

// ─────────────────────────────────────────────────────────────────────────────
// Duong NHAN so no luc dang nhap.
//
// Nua kia cua tinh nang nam o khuvuc_nothuong_test.go (ghi no khi boss chet).
// O day la nua tra: dang nhap -> lay nhan vat -> tra het no dang cho -> cong
// vao ban trong bo nho -> moi vao the gioi.
//
// THU TU LA THU QUAN TRONG NHAT: phai cong vao ban NhanVat trong bo nho NGAY
// sau khi CSDL cong xong. ThoGhi ghi DE CA DONG, nen mot ban trong RAM con
// mang kinh_nghiem cu se xoa sach phan vua cong o lan luu dinh ky dau tien —
// dung cai bay da khien khoan thuong nay khong tra duoc ngay tu dau.
// ─────────────────────────────────────────────────────────────────────────────

// khoBietTraNo la khoGia CO THEM phuong thuc NhanNoThuong. Bo xu ly nhan ra no
// bang type-assert (khuon ghiDongBoKhiBiDa, login.go) — nen moi kho gia CU chi
// cai 4 phuong thuc cua Kho van chay nguyen, khong phai sua mot dong nao.
type khoBietTraNo struct {
	*khoGia

	traVe         int64 // so diem no tra ve
	kinhNghiemMoi int64 // gia tri kinh_nghiem MOI ma CSDL tra ra sau khi cong
	loi           error

	soLanGoi int
	idDaGoi  int64
}

func (k *khoBietTraNo) NhanNoThuong(ctx context.Context, nhanVatID int64) (int64, int64, error) {
	k.soLanGoi++
	k.idDaGoi = nhanVatID
	if k.loi != nil {
		return 0, 0, k.loi
	}
	return k.traVe, k.kinhNghiemMoi, nil
}

// TestDangNhapNhanNoThuongCongVaoNhanVat la bai kiem trung tam cua nua tra:
// nguoi choi co 500 diem no phai thay kinh nghiem tang NGAY khi vao the gioi,
// khong phai cho lan luu ke tiep.
func TestDangNhapNhanNoThuongCongVaoNhanVat(t *testing.T) {
	nv := &store.NhanVat{ID: 7, Ten: "Naruto", KinhNghiem: 7}
	kho := &khoBietTraNo{
		khoGia:        &khoGia{tk: taiKhoanHopLeChoTest(), nhanVat: nv},
		traVe:         500,
		kinhNghiemMoi: 507, // 7 co san + 500 no, doc ra tu chinh giao dich
	}

	chayDangNhap(t, kho, "nguoidung", "matkhau")

	if kho.soLanGoi != 1 {
		t.Fatalf("NhanNoThuong duoc goi %d lan, mong doi dung 1", kho.soLanGoi)
	}
	if kho.idDaGoi != 7 {
		t.Errorf("NhanNoThuong goi voi id %d, mong doi 7 (id nhan vat)", kho.idDaGoi)
	}
	if nv.KinhNghiem != 507 {
		t.Errorf("KinhNghiem trong bo nho = %d, mong doi 507 (7 co san + 500 no). "+
			"Khong cong vao ban trong RAM thi lan luu dinh ky dau tien se ghi de mat", nv.KinhNghiem)
	}
}

// TestDangNhapKhongCoNoThiKhongDoiKinhNghiem: duong nay chay o MOI lan dang
// nhap, nen truong hop pho bien nhat la khong co no gi ca.
func TestDangNhapKhongCoNoThiKhongDoiKinhNghiem(t *testing.T) {
	nv := &store.NhanVat{ID: 7, Ten: "Naruto", KinhNghiem: 123}
	kho := &khoBietTraNo{
		khoGia: &khoGia{tk: taiKhoanHopLeChoTest(), nhanVat: nv},
		traVe:  0,
	}

	goiGuiRa := chayDangNhap(t, kho, "nguoidung", "matkhau")

	if nv.KinhNghiem != 123 {
		t.Errorf("KinhNghiem = %d, mong doi giu nguyen 123", nv.KinhNghiem)
	}
	if !coLenh(goiGuiRa, lenhKhongTrenBanDo) {
		t.Error("khong co no thi van phai vao the gioi binh thuong")
	}
}

// TestDangNhapNhanNoLoiThiVAN VaoDuocTheGioi: CSDL hong o buoc tra no KHONG
// duoc chan nguoi choi vao game. No van nam nguyen trong bang (giao dich bi
// huy, xem store.NhanNoThuong) nen lan dang nhap sau se tra — con chan cua vao
// vi mot khoan thuong la doi mot phien choi lay mot con so.
func TestDangNhapNhanNoLoiThiVanVaoDuocTheGioi(t *testing.T) {
	nv := &store.NhanVat{ID: 7, Ten: "Naruto", KinhNghiem: 50}
	kho := &khoBietTraNo{
		khoGia: &khoGia{tk: taiKhoanHopLeChoTest(), nhanVat: nv},
		loi:    errors.New("CSDL hong"),
	}

	goiGuiRa := chayDangNhap(t, kho, "nguoidung", "matkhau")

	if nv.KinhNghiem != 50 {
		t.Errorf("KinhNghiem = %d, mong doi giu nguyen 50 khi tra no loi", nv.KinhNghiem)
	}
	if !coLenh(goiGuiRa, lenhKhongTrenBanDo) {
		t.Error("tra no loi VAN phai vao the gioi (goi -28/4 versionGame)")
	}
	if coLenh(goiGuiRa, lenhTaoNhanVat) {
		t.Error("tra no loi khong duoc day nguoi choi ve man tao nhan vat")
	}
}

// TestDangNhapKhoKhongBietTraNoThiVanVaoDuoc: mot Kho khong cai NhanNoThuong
// (moi kho gia cu, va ca giai doan truoc khi noi day) phai chay y nhu truoc.
func TestDangNhapKhoKhongBietTraNoThiVanVaoDuoc(t *testing.T) {
	nv := &store.NhanVat{ID: 7, Ten: "Naruto", KinhNghiem: 9}
	kho := &khoGia{tk: taiKhoanHopLeChoTest(), nhanVat: nv}

	goiGuiRa := chayDangNhap(t, kho, "nguoidung", "matkhau")

	if nv.KinhNghiem != 9 {
		t.Errorf("KinhNghiem = %d, mong doi giu nguyen 9", nv.KinhNghiem)
	}
	if !coLenh(goiGuiRa, lenhKhongTrenBanDo) {
		t.Error("kho khong biet tra no thi van phai vao the gioi binh thuong")
	}
}

// TestDangNhapChuaCoNhanVatThiKhongGoiTraNo: chua co nhan vat thi khong co id
// nao de tra no — goi voi id 0 se quet nham (hoac vo nghia). Phai khong goi.
func TestDangNhapChuaCoNhanVatThiKhongGoiTraNo(t *testing.T) {
	kho := &khoBietTraNo{
		khoGia: &khoGia{tk: taiKhoanHopLeChoTest(), loiLayNhanVat: store.ErrKhongThayNhanVat},
		traVe:  500,
	}

	chayDangNhap(t, kho, "nguoidung", "matkhau")

	if kho.soLanGoi != 0 {
		t.Fatalf("NhanNoThuong duoc goi %d lan khi chua co nhan vat, mong doi 0", kho.soLanGoi)
	}
}

// TestDangNhapNhanNoBaoHoaKhongTran: mot nguoi choi gan tran kinh nghiem nhan
// them no thi phai dung o MaxInt64, khong cuon vong ve so am.
func TestDangNhapNhanNoBaoHoaKhongTran(t *testing.T) {
	const maxInt64 = int64(9_223_372_036_854_775_807)
	nv := &store.NhanVat{ID: 7, Ten: "Naruto", KinhNghiem: maxInt64 - 10}
	kho := &khoBietTraNo{
		khoGia:        &khoGia{tk: taiKhoanHopLeChoTest(), nhanVat: nv},
		traVe:         1000,
		kinhNghiemMoi: 9_223_372_036_854_775_807, // CSDL da bao hoa san
	}

	chayDangNhap(t, kho, "nguoidung", "matkhau")

	if nv.KinhNghiem != maxInt64 {
		t.Errorf("KinhNghiem = %d, mong doi bao hoa o %d", nv.KinhNghiem, maxInt64)
	}
}

// TestDangNhapGanKinhNghiemTuCSDLChuKhongTuCong khoa lua chon GAN thay vi CONG.
//
// Ban trong RAM co the da CU hon CSDL: mot ban ghi cua phien truoc con nam
// trong hang doi ThoGhi luc ta doc nhan vat. Tu cong lay tu so cu thi ra so cu
// — sai im lang. Gia tri dung la cai doc ra tu chinh giao dich vua cong.
//
// Test dung mot ban RAM CO Y lech han so voi ket qua CSDL tra ve, roi doi hoi
// ban RAM phai theo CSDL.
func TestDangNhapGanKinhNghiemTuCSDLChuKhongTuCong(t *testing.T) {
	nv := &store.NhanVat{ID: 7, Ten: "Naruto", KinhNghiem: 100} // ban RAM (cu)
	kho := &khoBietTraNo{
		khoGia:        &khoGia{tk: taiKhoanHopLeChoTest(), nhanVat: nv},
		traVe:         50,
		kinhNghiemMoi: 9999, // CSDL noi con so that la day, khong phai 100+50
	}

	chayDangNhap(t, kho, "nguoidung", "matkhau")

	if nv.KinhNghiem != 9999 {
		t.Errorf("KinhNghiem = %d, mong doi 9999 (GAN theo CSDL). Neu ra 150 tuc la "+
			"dang tu cong tu ban RAM — ban do co the da cu hon CSDL", nv.KinhNghiem)
	}
}
