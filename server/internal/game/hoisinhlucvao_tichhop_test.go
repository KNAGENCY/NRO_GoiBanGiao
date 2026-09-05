package game

import (
	"context"
	"testing"
	"time"

	"nroserver/internal/store"
)

// TestTichHop_NhanVatChetTrongCSDLDangNhapLaSong: hang nhan_vat co
// hp_hien_tai = 0 (nguoi choi chet roi thoat, KH11 luu nguyen) -> dang nhap
// lai -> sau -39, ban chup tren ACTOR phai thay nguoi do song (hp = hp toi da,
// chakra day). Truoc 05/09 khong co buoc nay: -42 kep mau ve day nen HUD hien
// day mau, con actor giu HpHienTai = 0 va tu choi moi don danh (khuvuc.go
// "nguoi chet khong danh duoc") — client va may chu noi hai su that khac nhau.
//
// Dot bien de xac nhan test co rang: bo dong hoiSinhLucVao(b.nhanVat) trong
// xuLyClientOk (login.go) -> test nay do o "hp tren actor = 0".
func TestTichHop_NhanVatChetTrongCSDLDangNhapLaSong(t *testing.T) {
	kho := moKhoThatChoTest(t)
	// Tien to "hsv" (hoi sinh vao) RIENG — "chet"+hauTo da duoc doimap_tichhop_test
	// va khuvuc_tichhop_test dung, va mot hang rac "chet8469" tu 03/09 (don dep
	// hong luc PostgreSQL tat giua chung) da lam test nay do "duplicate key" o
	// lan chay full dau tien 05/09.
	hauTo := hauToNganNgau(t)
	ten := "hsv" + hauTo
	idTK := taoTaiKhoanCoMatKhau(t, ten, "1")

	ctx := context.Background()
	lang := TimLang(0)
	nv, err := kho.TaoNhanVat(ctx, store.NhanVatMoi{
		TaiKhoanID: idTK, Ten: "Hsv" + hauTo, Lang: 0, KieuToc: HeadNhanVat,
		MapID: lang.MapKhoiDau, ViTriX: lang.ViTriXKhoiDau, ViTriY: lang.ViTriYKhoiDau,
		HpToiDa: lang.HpKhoiDau, ChakraToiDa: lang.ChakraKhoiDau, Vang: lang.VangKhoiDau,
	})
	if err != nil {
		t.Fatalf("tao nhan vat: %v", err)
	}

	// Giet nhan vat ngay trong CSDL — dung duong luu that (CapNhatNhanVat) de
	// hang co dung hinh dang ma ThoGhi de lai khi nguoi choi chet roi thoat.
	chet := *nv
	chet.HpHienTai = 0
	chet.ChakraHienTai = 0
	if err := kho.CapNhatNhanVat(ctx, chet); err != nil {
		t.Fatalf("ghi nhan vat chet: %v", err)
	}
	if lai, err := kho.LayNhanVatTheoTaiKhoan(ctx, idTK); err != nil || lai.HpHienTai != 0 {
		t.Fatalf("tien de hong: nhan vat trong CSDL phai co hp 0, got %+v err=%v", lai, err)
	}

	so := moSoKhuVucKhongQuaiMacDinh(t)
	cl := moPhienClientGiaKhuVuc(t, kho, so)
	defer cl.dong()
	dangNhapLaiVaoTheGioiVaKhuVuc(t, cl, ten, "1")

	// -39 di qua hop thu actor: poll ban chup toi khi thay nguoi (co han gio).
	k := so.Tra(lang.MapKhoiDau)
	var ban []store.NhanVat
	han := time.Now().Add(3 * time.Second)
	for time.Now().Before(han) {
		ban = k.AnhChupNguoiChoi(nv.ID)
		if len(ban) == 1 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if len(ban) != 1 {
		t.Fatalf("sau -39 actor khu vuc %d khong thay nhan vat %d", lang.MapKhoiDau, nv.ID)
	}
	if ban[0].HpHienTai != lang.HpKhoiDau {
		t.Fatalf("hp tren actor = %d, mong doi %d (nhan vat chet trong CSDL phai song lai luc vao)",
			ban[0].HpHienTai, lang.HpKhoiDau)
	}
	if ban[0].ChakraHienTai != lang.ChakraKhoiDau {
		t.Fatalf("chakra tren actor = %d, mong doi %d", ban[0].ChakraHienTai, lang.ChakraKhoiDau)
	}
	// Vi tri GIU NGUYEN (khong keo ve diem sinh) — cung luat voi suKienHoiSinh.
	if ban[0].ViTriX != nv.ViTriX || ban[0].ViTriY != nv.ViTriY {
		t.Fatalf("vi tri doi thanh (%d,%d), mong doi giu (%d,%d)", ban[0].ViTriX, ban[0].ViTriY, nv.ViTriX, nv.ViTriY)
	}
}
