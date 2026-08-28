package game

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"nroserver/internal/auth"
	"nroserver/internal/chanip"
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// File nay tai hien DUNG canh da do duoc va ghi vao so no bao mat (§7):
// "200 lan sai lien tiep khong bi chan, 76 lan/giay".
//
// Phep do di qua DUONG NOI DAY SAN XUAT (MoLuuTru -> taoXuLy), khong phai qua
// mot BoXuLy tu che — vi bai hoc cu cua du an la ha tang co the xay xong ma
// KHONG AI NOI DAY, luc do test tu che van xanh con may chu that van tho.

// khoDemDangNhap la kho gia dem so lan LayTaiKhoan thuc su bi goi.
//
// Dem cai nay chinh la thuoc do "ke xau ton duoc bao nhieu tai nguyen may
// chu": moi lan LayTaiKhoan la mot vong tra cuu CSDL CONG mot lan bam PBKDF2
// 100.000 vong (~13 ms CPU). Lan thu bi gac chan phai dung TRUOC do, khong
// duoc cham toi kho — neu van cham toi thi cai goi la "chan" chi la trang
// tri, may chu van bi vat kiet CPU nhu cu.
type khoDemDangNhap struct {
	tk     *store.TaiKhoan
	soLuot atomic.Int64
}

func (k *khoDemDangNhap) LayTaiKhoan(ctx context.Context, ten string) (*store.TaiKhoan, error) {
	k.soLuot.Add(1)
	if ten == k.tk.TenDangNhap {
		return k.tk, nil
	}
	return nil, store.ErrKhongThayTaiKhoan
}

func (k *khoDemDangNhap) LayNhanVatTheoTaiKhoan(context.Context, int64) (*store.NhanVat, error) {
	return nil, store.ErrKhongThayNhanVat
}
func (k *khoDemDangNhap) TenNhanVatDaDung(context.Context, string) (bool, error) {
	return false, nil
}
func (k *khoDemDangNhap) TaoNhanVat(context.Context, store.NhanVatMoi) (*store.NhanVat, error) {
	return nil, nil
}
func (k *khoDemDangNhap) CapNhatNhanVat(context.Context, store.NhanVat) error { return nil }

// docCauTuChoi doc cap goi tin ma guiTuChoi() luon gui (lenh -26 mang cau
// thong bao, roi lenh -102 keo client ra khoi trang thai "dang cho") va tra ve
// cau thong bao trong goi dau.
func docCauTuChoi(t *testing.T, f *protocol.Framer) string {
	t.Helper()
	m1 := docMotGoi(t, f)
	if m1.Command != lenhThongBaoOK {
		t.Fatalf("goi dau tien co lenh %d, mong doi %d (hop thoai bao loi)", m1.Command, lenhThongBaoOK)
	}
	cau, err := protocol.NewReader(m1.Data).UTF()
	if err != nil {
		t.Fatalf("khong doc duoc cau thong bao: %v", err)
	}
	m2 := docMotGoi(t, f)
	if m2.Command != lenhDangNhapThatBai {
		t.Fatalf("goi thu hai co lenh %d, mong doi %d — thieu goi nay la client treo o man 'Xin cho'",
			m2.Command, lenhDangNhapThatBai)
	}
	return cau
}

// TestDoMatKhau200LanLienTiepBiChan la phep do CHOT cua lo hong 1.
//
// Kich ban dung y canh da do duoc truoc day: MOT ket noi, MOT tai khoan CO
// THAT, 200 mat khau sai lien tiep. Truoc ban va: ca 200 lan deu duoc phuc vu
// (~76 lan/giay). Sau ban va: chi mot so it dau tien di toi duoc kho, phan con
// lai bi tu choi TRUOC KHI cham CSDL/ham bam.
func TestDoMatKhau200LanLienTiepBiChan(t *testing.T) {
	bam, err := auth.Hash("matkhaudung")
	if err != nil {
		t.Fatalf("khong bam duoc mat khau mau: %v", err)
	}
	kho := &khoDemDangNhap{tk: &store.TaiKhoan{ID: 1, TenDangNhap: "nguoidung", MatKhauBam: bam}}

	soKhuVuc := MoSoKhuVuc()
	t.Cleanup(func() { _ = soKhuVuc.DungTatCa() })

	// BAT BUOC di qua MoLuuTru: day la duong lap rap SAN XUAT. Neu gac chong
	// do duoc viet xong ma quen noi vao day, test nay phai do.
	thoGhi, boLuu, taoXuLy := MoLuuTru(kho, soKhuVuc, 1, 0, 0)
	t.Cleanup(boLuu.Dung)
	t.Cleanup(func() { _ = thoGhi.DongVaChoXong(context.Background()) })

	b := taoXuLy()
	fClient, dong := dungPhienThat(t, b)
	defer dong()

	const soLanThu = 200
	batDau := time.Now()
	soBiChan := 0
	for i := 0; i < soLanThu; i++ {
		if err := fClient.WriteMessage(dungGoiDangNhap("nguoidung", fmt.Sprintf("saibet%d", i))); err != nil {
			t.Fatalf("gui lan thu %d that bai: %v", i+1, err)
		}
		// Bat ky cau tra loi nao KHAC cau "sai tai khoan hoac mat khau" deu la
		// dau hieu gac da chan. Khong so khop cung mot hang so o day, de chinh
		// phep do nay chay duoc CA TRUOC lan SAU khi vá — co nhu vay con so
		// "truoc" va "sau" moi so sanh duoc voi nhau.
		cau := docCauTuChoi(t, fClient)
		if cau == LyDoSaiThongTin {
			continue
		}
		soBiChan++
		if cau != LyDoQuaNhieuLanSai {
			t.Fatalf("lan thu %d bi chan nhung tra ve cau la: %q", i+1, cau)
		}
	}
	tong := time.Since(batDau)
	soDenKho := kho.soLuot.Load()

	// Hai con so phai doc CUNG NHAU thi moi khong hieu nham:
	//
	//   - "dot dau": bao nhieu lan thu di toi duoc CSDL + ham bam TRUOC khi
	//     het the. Day la toan bo so lan doan ma ke xau duoc huong ngay lap
	//     tuc. Truoc khi va no la 200/200 (~2,1 giay cho ca 200 lan).
	//
	//   - "duy tri": sau dot dau, moi nhip hoi phuc moi duoc them DUNG mot
	//     lan thu. Day moi la con so quyet dinh do bua co kha thi hay khong —
	//     va la con so phai doi chieu voi "76 lan/giay" da ghi trong so.
	nhipHoiPhuc := chanip.NhipHoiPhucMacDinh
	duyTriMoiGio := time.Hour.Seconds() / nhipHoiPhuc.Seconds()
	canBaoLau := time.Duration(soLanThu-chanip.SucChuaMacDinh) * nhipHoiPhuc
	t.Logf("200 lan sai lien tiep: dot dau %d lan di toi duoc CSDL+ham bam (het %v), "+
		"%d lan con lai bi gac chan truoc khi cham CSDL. "+
		"Toc do duy tri sau do: 1 lan moi %v = %.0f lan/gio. "+
		"Muon do du 200 lan phai mat %v.",
		soDenKho, tong, soBiChan, nhipHoiPhuc, duyTriMoiGio, canBaoLau)

	// Nguong: chi cho phep so lan "ton tai nguyen may chu" bang dung suc chua
	// cua gac. Vuot qua la gac khong chan, hoac chan SAU khi da tra gia CPU roi.
	if soDenKho > int64(chanip.SucChuaMacDinh) {
		t.Fatalf("CHUA CHONG DO: %d/200 lan sai van di toi CSDL va ham bam PBKDF2 "+
			"(chi duoc phep %d). Ke xau do mat khau khong bi chan, dong thoi vat kiet CPU may chu.",
			soDenKho, chanip.SucChuaMacDinh)
	}
	if soBiChan < soLanThu-chanip.SucChuaMacDinh {
		t.Fatalf("chi %d/200 lan bi gac chan — phan con lai van duoc phuc vu binh thuong", soBiChan)
	}

	// DIEU QUAN TRONG NHAT cua ca ban va nay, quan trong hon ca viec chan
	// duoc: 200 lan do vao tai khoan "nguoidung" KHONG DUOC lam chinh chu
	// tai khoan do mat quyen vao game.
	//
	// Neu bo dem duoc dat theo TEN TAI KHOAN, mot ke xau chi can co tinh go
	// sai mat khau cua nguoi khac vai chuc lan la khoa duoc ho ra ngoai —
	// lop bao ve tro thanh vu khi de bat nat nguoi choi. Bo dem o day dat
	// theo DIA CHI MANG, nen hai khang dinh duoi day phai dung cung luc:
	// dia chi vua do thi het luot, moi dia chi khac van con nguyen luot.
	if b.gacDangNhap.ChoPhep("khong-ro") {
		t.Fatal("dia chi vua do 200 lan van con luot thu — gac khong chan")
	}
	if !b.gacDangNhap.ChoPhep("113.161.0.7") {
		t.Fatal("mot dia chi mang KHAC bi va lay sau khi ke xau do mat khau tai khoan nay — " +
			"nghia la bo dem dang dat theo TEN TAI KHOAN chu khong theo dia chi mang. " +
			"Nhu vay ke xau co the co tinh go sai de KHOA TAI KHOAN cua nguoi khac.")
	}
}

// Gac chong do phai duoc noi vao duong lap rap SAN XUAT. Ha tang xay xong ma
// khong ai goi toi la bai hoc cu cua du an (Task 6): test tu che van xanh
// trong khi may chu that van tho nguyen.
func TestMoLuuTruNoiGacChongDo(t *testing.T) {
	kho := &khoDayDuGia{}
	soKhuVuc := MoSoKhuVuc()
	t.Cleanup(func() { _ = soKhuVuc.DungTatCa() })

	thoGhi, boLuu, taoXuLy := MoLuuTru(kho, soKhuVuc, 1, 0, 0)
	t.Cleanup(boLuu.Dung)
	t.Cleanup(func() { _ = thoGhi.DongVaChoXong(context.Background()) })

	b1 := taoXuLy()
	b2 := taoXuLy()
	if b1.gacDangNhap == nil {
		t.Fatal("MoLuuTru khong gan gac chong do vao BoXuLy — may chu chay khong co gi chan do mat khau")
	}
	// PHAI la CUNG MOT gac cho moi phien. Neu moi ket noi mot gac rieng thi ke
	// do mat khau chi viec mo ket noi moi cho moi lan thu, bo dem lai bat dau
	// tu 0, va gac co cung nhu khong.
	if b1.gacDangNhap != b2.gacDangNhap {
		t.Fatal("moi ket noi duoc mot gac RIENG — ke do mat khau chi can mo ket noi moi moi lan " +
			"thu la bo dem lai ve 0, gac tro thanh do trang tri")
	}
}

// Chan do mat khau KHONG duoc lam nguoi choi that mat quyen vao game: go nham
// vai lan roi go dung thi VAN PHAI vao duoc.
func TestGoNhamVaiLanRoiGoDungVanVaoDuoc(t *testing.T) {
	bam, err := auth.Hash("matkhaudung")
	if err != nil {
		t.Fatalf("khong bam duoc mat khau mau: %v", err)
	}
	kho := &khoDemDangNhap{tk: &store.TaiKhoan{ID: 1, TenDangNhap: "nguoidung", MatKhauBam: bam}}

	soKhuVuc := MoSoKhuVuc()
	t.Cleanup(func() { _ = soKhuVuc.DungTatCa() })
	thoGhi, boLuu, taoXuLy := MoLuuTru(kho, soKhuVuc, 1, 0, 0)
	t.Cleanup(boLuu.Dung)
	t.Cleanup(func() { _ = thoGhi.DongVaChoXong(context.Background()) })

	fClient, dong := dungPhienThat(t, taoXuLy())
	defer dong()

	// Nam lan go nham — con so nay phai nam THOAI MAI trong suc chua, day la
	// muc ma mot nguoi choi that hoan toan co the cham toi.
	for i := 0; i < 5; i++ {
		if err := fClient.WriteMessage(dungGoiDangNhap("nguoidung", "gonham")); err != nil {
			t.Fatalf("gui lan nham thu %d that bai: %v", i+1, err)
		}
		if cau := docCauTuChoi(t, fClient); cau != LyDoSaiThongTin {
			t.Fatalf("lan go nham thu %d bi tra loi %q — nguoi choi that dang bi da ra qua som", i+1, cau)
		}
	}

	// Lan thu sau go DUNG: phai vao duoc (tai khoan chua co nhan vat -> lenh 2).
	if err := fClient.WriteMessage(dungGoiDangNhap("nguoidung", "matkhaudung")); err != nil {
		t.Fatalf("gui goi dang nhap dung that bai: %v", err)
	}
	m := docMotGoi(t, fClient)
	if m.Command != lenhTaoNhanVat {
		t.Fatalf("go dung mat khau sau 5 lan nham nhan lenh %d, mong doi %d (vao duoc)", m.Command, lenhTaoNhanVat)
	}
}
