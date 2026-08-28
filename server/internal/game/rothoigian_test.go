package game

import (
	"context"
	"sort"
	"testing"
	"time"

	"nroserver/internal/auth"
	"nroserver/internal/store"
)

// File nay do MOT thu duy nhat: hai nhanh tu choi dang nhap co ton BANG NHAU
// thoi gian khong.
//
// VI SAO quan trong: cau tra loi gui ve client la NHU NHAU cho ca hai nhanh
// ("Sai tài khoản hoặc mật khẩu." — xem LyDoSaiThongTin), nhung neu nhanh
// "khong co tai khoan" tra ve NGAY con nhanh "co tai khoan, sai mat khau"
// phai chay qua ham bam PBKDF2 100.000 vong (~13 ms), thi ke xau chi can BAM
// GIO la doc duoc dieu ma cau chu khong noi: ten nao co that. Biet ten nao co
// that roi, ho don het suc do mat khau vao dung nhung ten do — thu hep khong
// gian tim kiem di hang nghin lan so voi do bua.
//
// Cach vá: nhanh "khong tim thay tai khoan" van chay MOT lan bam GIA
// (auth.CanBangThoiGian) de hai nhanh ton cung mot khoi luong tinh toan.

// Kho dung o day la khoHaiTaiKhoan (dangnhapchong_test.go): no tra tai khoan
// theo TEN, ten khong co trong ban do thi bao khong ton tai — dung dieu phep
// do nay can. khoGia (login_test.go) khong dung duoc vi no tra CUNG mot ket
// qua cho moi ten, khong tao duoc hai nhanh khac nhau.

// trungVi tra ve gia tri GIUA cua mot day thoi gian.
//
// Dung trung vi chu KHONG dung trung binh: may dang chay test nay con chay ca
// viec khac, thinh thoang mot phep do bi he dieu hanh cuop CPU va vot len hang
// chuc mili-giay. Vai gia tri vot nhu vay keo TRUNG BINH lech han, nhung gan
// nhu khong lam xe dich TRUNG VI — nen trung vi la thuoc do trung thuc hon
// nhieu cho "mot lan goi binh thuong ton bao lau".
func trungVi(d []time.Duration) time.Duration {
	ban := append([]time.Duration(nil), d...)
	sort.Slice(ban, func(i, j int) bool { return ban[i] < ban[j] })
	return ban[len(ban)/2]
}

// doXenKe do THOI GIAN cua hai nhanh dang nhap bang cach goi XEN KE nhau
// (A, B, A, B...) thay vi chay het A roi moi chay B.
//
// VI SAO xen ke: neu do het 120 lan A truoc roi 120 lan B sau, moi thay doi
// cham cua may trong khoang giua (CPU nong nen ha xung nhip, mot tien trinh
// khac khoi dong, bo thu gom rac cua Go chay mot dot) se don HET vao mot ben
// va bi doc nham thanh "chenh lech giua hai nhanh". Xen ke lam moi thay doi
// cham do roi deu len ca hai ben nen no tu triet tieu.
func doXenKe(b *BoXuLy, tenA, tenB string, soLan int) (medA, medB time.Duration) {
	ctx := context.Background()
	a := make([]time.Duration, 0, soLan)
	bb := make([]time.Duration, 0, soLan)
	// Chay nong may truoc khi do: lan goi dau tien phai cap phat bo nho va nap
	// ma vao cache CPU, ton hon han cac lan sau — de lot vao so lieu la sai lech.
	for i := 0; i < 5; i++ {
		_, _ = b.KiemTraDangNhap(ctx, tenA, "mat-khau-sai-hoan-toan")
		_, _ = b.KiemTraDangNhap(ctx, tenB, "mat-khau-sai-hoan-toan")
	}
	for i := 0; i < soLan; i++ {
		t0 := time.Now()
		_, _ = b.KiemTraDangNhap(ctx, tenA, "mat-khau-sai-hoan-toan")
		a = append(a, time.Since(t0))

		t1 := time.Now()
		_, _ = b.KiemTraDangNhap(ctx, tenB, "mat-khau-sai-hoan-toan")
		bb = append(bb, time.Since(t1))
	}
	return trungVi(a), trungVi(bb)
}

func triTuyetDoi(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

// TestKhongDoanDuocTaiKhoanNaoCoThatQuaThoiGian la phep do CHOT cua lo hong 2.
//
// Cach do:
//  1. Do NHIEU cua chinh phep do truoc: chay hai lan do tren HAI ten DEU KHONG
//     ton tai. Ve nguyen tac hai nhanh nay giong het nhau, nen chenh lech doc
//     duoc chinh la sai so cua dung cu do — khong the doi phep do phan biet
//     duoc hai thu chenh nhau it hon the.
//  2. Do THAT: nhanh "khong co tai khoan" vs nhanh "co tai khoan, sai mat khau".
//  3. Doi hoi chenh lech that phai nam trong nguong — chon nguong la GIA TRI
//     LON HON giua (5 lan nhieu do) va 1 ms, de test khong chop tat tren may
//     ban ma van bat duoc ro ret lo hong that (~13 ms, tuc gap hon 13 lan
//     nguong 1 ms).
func TestKhongDoanDuocTaiKhoanNaoCoThatQuaThoiGian(t *testing.T) {
	bam, err := auth.Hash("matkhaudung")
	if err != nil {
		t.Fatalf("khong bam duoc mat khau mau: %v", err)
	}
	kho := &khoHaiTaiKhoan{theoTen: map[string]*store.TaiKhoan{
		"cothat": {ID: 1, TenDangNhap: "cothat", MatKhauBam: bam},
	}}
	b := NewBoXuLy(kho, nil)

	const soLan = 120

	// Buoc 1 — do nhieu: hai ten CUNG khong ton tai, ket qua phai gan nhu bang nhau.
	nhieuA, nhieuB := doXenKe(b, "khongco01", "khongco02", soLan)
	nhieu := triTuyetDoi(nhieuA - nhieuB)

	// Buoc 2 — do that.
	medKhongCo, medCoThat := doXenKe(b, "khongco01", "cothat", soLan)
	chenh := triTuyetDoi(medCoThat - medKhongCo)

	nguong := 5 * nhieu
	if nguong < time.Millisecond {
		nguong = time.Millisecond
	}

	t.Logf("nhieu do (hai nhanh giong het nhau) = %v | khong co tai khoan = %v | "+
		"co tai khoan sai mat khau = %v | chenh lech = %v | nguong = %v",
		nhieu, medKhongCo, medCoThat, chenh, nguong)

	if chenh > nguong {
		t.Fatalf(
			"RO RI THOI GIAN: nhanh 'co tai khoan' va nhanh 'khong co tai khoan' chenh nhau %v "+
				"(nguong %v, nhieu do %v). Ke xau bam gio la liet ke duoc tai khoan nao co that. "+
				"Cach vá: nhanh khong tim thay tai khoan phai chay auth.CanBangThoiGian(matKhau).",
			chenh, nguong, nhieu)
	}
}
