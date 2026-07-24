package game

// Kiem thu tich hop (Task 7, Ke hoach 5): actor khu vuc noi vao luong dang
// nhap THAT qua mang that (net.Pipe + mynet.PhucVu) + CSDL that (nro_test),
// dung MOT SoKhuVuc CHUNG cho ca hai phien — dung y het cach main.go dung (mot
// so khu vuc song suot doi tien trinh, nhieu phien tra cung mot khu vuc theo
// mapId thay vi moi phien tu MoKhuVuc rieng nhu Ke hoach 4).
//
// Kich ban: hai nguoi choi cung vao mapId 1 (lang 0) -> ca hai gui -39 (da tai
// xong map) -> moi ben thay nguoi kia qua goi -5 -> A gui -7 (di chuyen trong
// nguong) -> B nhan goi -7 chuyen tiep, A KHONG nhan lai goi cua chinh minh.
//
// Ca "vuot nguong bi chan" va "roi di phat -6" CO CHU DICH de danh cho Task 8
// (xem ghi chu cuoi task-7-brief-kh5.md) — khong lap lai o day.

import (
	"context"
	"net"
	"testing"
	"time"

	mynet "nroserver/internal/net"
	"nroserver/internal/protocol"
)

// moPhienClientGiaKhuVuc giong het moPhienClientGia (tichhop_test.go) — cung
// buoc bat tay, cung kieu tra ve *clientGia — nhung nhan them mot *SoKhuVuc
// THAT (khac nil) de BoXuLy co the dua nhan vat vao khu vuc chung khi client
// gui -39. moPhienClientGia goc luon truyen nil (an toan cho moi test khac
// trong package nay, vi cac test do khong cham -39/-7), nen can ham rieng o
// day thay vi doi chu ky cua no.
func moPhienClientGiaKhuVuc(t *testing.T, kho Kho, soKhuVuc *SoKhuVuc) *clientGia {
	t.Helper()
	b := NewBoXuLy(kho, soKhuVuc)

	benClient, benServer := net.Pipe()
	go mynet.PhucVu(benServer, func() mynet.XuLyGoiTin { return b })

	benClient.SetReadDeadline(time.Now().Add(5 * time.Second))

	f := protocol.NewFramer(benClient, benClient)
	if err := f.WriteMessage(&protocol.Message{Command: -27}); err != nil {
		t.Fatalf("gui goi bat tay that bai: %v", err)
	}
	khoaTho := docGoiKhoaThuCong(t, benClient)
	f.SetKey(protocol.DeriveKey(khoaTho))

	return &clientGia{t: t, f: f, conn: benClient}
}

// baoVaoKhuVuc gui goi -39 than rong (client bao da tai xong map).
func (cl *clientGia) baoVaoKhuVuc() {
	cl.t.Helper()
	if err := cl.f.WriteMessage(&protocol.Message{Command: lenhVaoKhuVuc}); err != nil {
		cl.t.Fatalf("gui goi -39 (vao khu vuc) that bai: %v", err)
	}
}

// guiDiChuyen gui goi -7: byte co bay (0 = khong khinh cong, KH5 doc bo qua)
// roi short toX + short toY.
func (cl *clientGia) guiDiChuyen(toX, toY int16) {
	cl.t.Helper()
	w := protocol.NewWriter()
	w.Byte(0)
	w.Short(toX)
	w.Short(toY)
	if err := cl.f.WriteMessage(protocol.NewMessage(lenhDiChuyen, w)); err != nil {
		cl.t.Fatalf("gui goi -7 (di chuyen) that bai: %v", err)
	}
}

// dangNhapTaoNhanVatVaoTheGioi dua MOT client gia di TRON version dance +
// vao the gioi: dang nhap (tai khoan MOI, chua co nhan vat) -> nhan lenh 2 ->
// tao nhan vat -> nhan -28/4 -> xin du 4 catalog -> gui -28/13 -> nhan
// -82/-30/-24. Dung chung cho ca hai nguoi choi trong test nay de khong lap
// lai chuoi buoc da co san o dongbo_tichhop_test.go (TestVersionDanceDayDu)
// cho tung ben.
func dangNhapTaoNhanVatVaoTheGioi(t *testing.T, cl *clientGia, ten, matKhau, tenNV string) {
	t.Helper()
	cl.dangNhap(ten, matKhau)
	if !cl.choLenh(t, 2) {
		t.Fatal("tai khoan moi chua co nhan vat -> phai nhan DUNG MOT lenh 2 (man tao nhan vat)")
	}
	cl.taoNhanVat(tenNV, 0, 5) // lang 0 -> spawn mapId 1
	kiemVersionGame(t, cl.doiGoi(t, lenhKhongTrenBanDo, lenhConVersionGame))

	cl.xinData()
	cl.doiGoiLenh(t, lenhDuLieu)
	cl.xinCatalog(lenhConCatalogMap)
	cl.doiGoi(t, lenhKhongTrenBanDo, lenhConCatalogMap)
	cl.xinCatalog(lenhConCatalogSkill)
	cl.doiGoi(t, lenhKhongTrenBanDo, lenhConCatalogSkill)
	cl.xinCatalog(lenhConCatalogItem)
	cl.doiGoi(t, lenhKhongTrenBanDo, lenhConCatalogItem)

	cl.guiClientOk()
	cl.doiGoiLenh(t, -82)
	cl.doiGoiLenh(t, -30)
	cl.doiGoiLenh(t, -24)
}

// TestTichHop_HaiNguoiThayNhauVaDiChuyen la kich ban chinh cua Task 7: hai
// phien THAT dung CHUNG mot SoKhuVuc (nhu main.go), lan luot vao khu vuc
// (-39) roi thay nhau qua -5, sau do A di chuyen (-7) va B nhan duoc goi
// chuyen tiep — trong khi luong phien (xuLyDiChuyen) KHONG doc nv.ViTri, chi
// actor (kiemDiChuyenHopLe trong khuvuc.go) moi la noi kiem hop le va ghi vi
// tri authoritative (xem "SUA LOI TRANH CHAP" trong task-7-brief-kh5.md).
func TestTichHop_HaiNguoiThayNhauVaDiChuyen(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	soKhuVuc := MoSoKhuVuc()
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTKA := "khuA" + hauToNganNgau(t)
	tenTKB := "khuB" + hauToNganNgau(t)
	taiKhoanA := taoTaiKhoanCoMatKhau(t, tenTKA, matKhau)
	taiKhoanB := taoTaiKhoanCoMatKhau(t, tenTKB, matKhau)
	tenNVA := "ZoneA" + hauToNganNgau(t)
	tenNVB := "ZoneB" + hauToNganNgau(t)

	clA := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer clA.dong()
	clB := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer clB.dong()

	dangNhapTaoNhanVatVaoTheGioi(t, clA, tenTKA, matKhau, tenNVA)
	dangNhapTaoNhanVatVaoTheGioi(t, clB, tenTKB, matKhau, tenNVB)

	nvA, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoanA)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan A that bai: %v", err)
	}
	nvB, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoanB)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan B that bai: %v", err)
	}
	if nvA.MapID != 1 || nvB.MapID != 1 {
		t.Fatalf("ca hai nhan vat phai o mapId 1 (lang 0), nhan duoc A=%d B=%d", nvA.MapID, nvB.MapID)
	}

	// A vao khu vuc TRUOC, mot minh -> khong ai de thay, khong nhan goi nao.
	clA.baoVaoKhuVuc()
	if !clA.khongConGoiGi(t) {
		t.Fatal("A vao khu vuc mot minh -> khong duoc nhan goi nao (chua co ai khac)")
	}

	// B vao SAU -> actor phat -5 cua A cho B truoc (nguoi cu), roi -5 cua B
	// cho A (nguoi con lai) — dung THU TU o khuvuc.go case suKienVao. Phat tin
	// di qua goroutine actor rieng + mang that nen dung doiGoiLenh (co han
	// gio, bo qua goi khac neu lo) thay vi doc thang mot lan.
	clB.baoVaoKhuVuc()

	goiChoB := clB.doiGoiLenh(t, -5)
	khangDinhIdTrongGoi(t, goiChoB, nvA.ID)

	goiChoA := clA.doiGoiLenh(t, -5)
	khangDinhIdTrongGoi(t, goiChoA, nvB.ID)

	// A di chuyen toi vi tri gan (trong nguong 500, xem NguongBuoc trong
	// dichuyen.go) -> B phai nhan -7 [idA][toX][toY]; A KHONG duoc nhan lai
	// goi cua chinh minh (client tu ve chuyen dong).
	toX := int16(nvA.ViTriX) + 50
	toY := int16(nvA.ViTriY) + 16
	clA.guiDiChuyen(toX, toY)

	goiDiChuyenChoB := clB.doiGoiLenh(t, -7)
	r := protocol.NewReader(goiDiChuyenChoB.Data)
	if id, errID := r.Int(); errID != nil || id != int32(nvA.ID) {
		t.Fatalf("id trong goi -7 = %d (loi %v), mong doi %d", id, errID, nvA.ID)
	}
	if x, errX := r.Short(); errX != nil || x != toX {
		t.Fatalf("toX trong goi -7 = %d (loi %v), mong doi %d", x, errX, toX)
	}
	if y, errY := r.Short(); errY != nil || y != toY {
		t.Fatalf("toY trong goi -7 = %d (loi %v), mong doi %d", y, errY, toY)
	}

	if !clA.khongConGoiGi(t) {
		t.Fatal("A tu di chuyen -> KHONG duoc nhan lai goi -7 cua chinh minh")
	}
}

// TestTichHop_DiChuyenVuotNguongBiChan la kich ban 1 cua Task 8: hai phien
// THAT (giong het TestTichHop_HaiNguoiThayNhauVaDiChuyen) da thay nhau qua
// -5, roi A gui -7 voi dich CACH vi tri hien tai > NguongBuoc (500, xem
// dichuyen.go) — actor (kiemDiChuyenHopLe trong khuvuc.go, case
// suKienDiChuyen) phai BO goi nay va KHONG phat cho ai, nen B khong duoc
// nhan goi -7 nao trong han cho im lang. Sau do A gui tiep MOT -7 HOP LE
// (dich gan, trong nguong) de chung minh ket noi VAN SONG BINH THUONG — chi
// goi bat hop le bi bo, may chu khong ngat ket noi vi no.
func TestTichHop_DiChuyenVuotNguongBiChan(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	soKhuVuc := MoSoKhuVuc()
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTKA := "teleA" + hauToNganNgau(t)
	tenTKB := "teleB" + hauToNganNgau(t)
	taiKhoanA := taoTaiKhoanCoMatKhau(t, tenTKA, matKhau)
	taiKhoanB := taoTaiKhoanCoMatKhau(t, tenTKB, matKhau)
	tenNVA := "TeleA" + hauToNganNgau(t)
	tenNVB := "TeleB" + hauToNganNgau(t)

	clA := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer clA.dong()
	clB := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer clB.dong()

	dangNhapTaoNhanVatVaoTheGioi(t, clA, tenTKA, matKhau, tenNVA)
	dangNhapTaoNhanVatVaoTheGioi(t, clB, tenTKB, matKhau, tenNVB)

	nvA, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoanA)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan A that bai: %v", err)
	}
	nvB, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoanB)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan B that bai: %v", err)
	}
	if nvA.MapID != 1 || nvB.MapID != 1 {
		t.Fatalf("ca hai nhan vat phai o mapId 1 (lang 0), nhan duoc A=%d B=%d", nvA.MapID, nvB.MapID)
	}

	clA.baoVaoKhuVuc()
	if !clA.khongConGoiGi(t) {
		t.Fatal("A vao khu vuc mot minh -> khong duoc nhan goi nao (chua co ai khac)")
	}
	clB.baoVaoKhuVuc()

	goiChoB := clB.doiGoiLenh(t, -5)
	khangDinhIdTrongGoi(t, goiChoB, nvA.ID)
	goiChoA := clA.doiGoiLenh(t, -5)
	khangDinhIdTrongGoi(t, goiChoA, nvB.ID)

	// A gui -7 toi dich CACH vi tri hien tai RAT XA (>> NguongBuoc = 500) —
	// van trong bien ban do [0, BienRongToiDa=4096] nen khong bi chan boi kiem
	// bien, CHI bi chan boi kiem khoang cach. Dung toa do co dinh xa han (2000,
	// 2000) thay vi cong tuong doi vao nvA.ViTriX/Y nhu test hop le o tren, de
	// khoang cach chac chan >> 500 bat ke spawn cua lang 0 nam o dau trong
	// bien 4096 (ViTriXKhoiDau=100 trong TestTaoRoiDangNhapLai_ThayNhanVat).
	toXQuaXa := int16(2000)
	toYQuaXa := int16(2000)
	clA.guiDiChuyen(toXQuaXa, toYQuaXa)

	// Khang dinh B KHONG nhan goi -7 nao trong han cho im lang: actor bo goi
	// bat hop le TRUOC khi phat (return som trong case suKienDiChuyen), nen
	// khong ai — ke ca B — nhan duoc gi.
	if !clB.khongConGoiGi(t) {
		t.Fatal("A gui -7 vuot nguong buoc -> B KHONG duoc nhan goi -7 nao (may chu phai bo goi bat hop le)")
	}

	// Chung minh ket noi CUA A VAN SONG BINH THUONG (khong bi ngat vi goi bat
	// hop le): gui tiep MOT -7 HOP LE (dich gan, trong nguong) -> B phai nhan
	// dung goi -7 chuyen tiep nhu kich ban duong hanh phuc.
	toXHopLe := int16(nvA.ViTriX) + 40
	toYHopLe := int16(nvA.ViTriY) + 10
	clA.guiDiChuyen(toXHopLe, toYHopLe)

	goiDiChuyenChoB := clB.doiGoiLenh(t, -7)
	r := protocol.NewReader(goiDiChuyenChoB.Data)
	if id, errID := r.Int(); errID != nil || id != int32(nvA.ID) {
		t.Fatalf("id trong goi -7 (hop le, sau khi goi vuot nguong bi bo) = %d (loi %v), mong doi %d", id, errID, nvA.ID)
	}
	if x, errX := r.Short(); errX != nil || x != toXHopLe {
		t.Fatalf("toX trong goi -7 hop le = %d (loi %v), mong doi %d", x, errX, toXHopLe)
	}
	if y, errY := r.Short(); errY != nil || y != toYHopLe {
		t.Fatalf("toY trong goi -7 hop le = %d (loi %v), mong doi %d", y, errY, toYHopLe)
	}
}

// TestTichHop_RoiKhuVucPhatGoi6 la kich ban 2 cua Task 8: hai phien THAT da
// thay nhau qua -5, roi DONG ket noi cua A — PhucVu (session.go) chay defer
// KhiDong() khi vong doc goi tin loi (ket noi dong), KhiDong gui suKienRa cho
// actor khu vuc, actor xoa A khoi map va phat -6 [idA] cho nguoi con lai (B).
func TestTichHop_RoiKhuVucPhatGoi6(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	soKhuVuc := MoSoKhuVuc()
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTKA := "roiA" + hauToNganNgau(t)
	tenTKB := "roiB" + hauToNganNgau(t)
	taiKhoanA := taoTaiKhoanCoMatKhau(t, tenTKA, matKhau)
	taiKhoanB := taoTaiKhoanCoMatKhau(t, tenTKB, matKhau)
	tenNVA := "RoiA" + hauToNganNgau(t)
	tenNVB := "RoiB" + hauToNganNgau(t)

	clA := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	clB := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer clB.dong()

	dangNhapTaoNhanVatVaoTheGioi(t, clA, tenTKA, matKhau, tenNVA)
	dangNhapTaoNhanVatVaoTheGioi(t, clB, tenTKB, matKhau, tenNVB)

	nvA, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoanA)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan A that bai: %v", err)
	}
	nvB, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoanB)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan B that bai: %v", err)
	}
	if nvA.MapID != 1 || nvB.MapID != 1 {
		t.Fatalf("ca hai nhan vat phai o mapId 1 (lang 0), nhan duoc A=%d B=%d", nvA.MapID, nvB.MapID)
	}

	clA.baoVaoKhuVuc()
	if !clA.khongConGoiGi(t) {
		t.Fatal("A vao khu vuc mot minh -> khong duoc nhan goi nao (chua co ai khac)")
	}
	clB.baoVaoKhuVuc()

	// Xa het hang doi -5 truoc khi dong A, de goi -6 sau nay khong bi lan voi
	// -5 con sot lai trong bo dem doc.
	goiChoB := clB.doiGoiLenh(t, -5)
	khangDinhIdTrongGoi(t, goiChoB, nvA.ID)
	goiChoA := clA.doiGoiLenh(t, -5)
	khangDinhIdTrongGoi(t, goiChoA, nvB.ID)

	// DONG ket noi cua A som (thay vi doi defer cuoi ham) de kich hoat KhiDong
	// ngay tai day: vong doc cua PhucVu phia may chu se gap loi doc (ket noi
	// dong) -> defer KhiDong() chay -> actor nhan suKienRa -> phat -6 cho B.
	clA.dong()

	goiRaChoB := clB.doiGoiLenh(t, -6)
	khangDinhIdTrongGoi(t, goiRaChoB, nvA.ID)
}
