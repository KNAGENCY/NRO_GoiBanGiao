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
	"bytes"
	"context"
	"fmt"
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

// chonChieu gui goi len 34 (chon chieu, dac ta §3.1/§3.2): than tin CHI co
// short templateId — la MauChieu.Id (template.id), TUYET DOI KHONG PHAI
// skillId (xem chu thich lenhChonChieu trong chieuthuc.go).
func (cl *clientGia) chonChieu(templateId int16) {
	cl.t.Helper()
	w := protocol.NewWriter()
	w.Short(templateId)
	if err := cl.f.WriteMessage(protocol.NewMessage(lenhChonChieu, w)); err != nil {
		cl.t.Fatalf("gui goi 34 (chon chieu) that bai: %v", err)
	}
}

// guiDanhQuai gui goi len 54 DANG THUONG (1 byte, mobId != -1): xin danh con
// quai co Id = mobId trong khu vuc (dac ta §3.1).
func (cl *clientGia) guiDanhQuai(mobId byte) {
	cl.t.Helper()
	w := protocol.NewWriter()
	w.Byte(mobId)
	if err := cl.f.WriteMessage(protocol.NewMessage(lenhNguoiDanhQuai, w)); err != nil {
		cl.t.Fatalf("gui goi 54 (danh quai) that bai: %v", err)
	}
}

// guiDanhQuaiThuCung gui goi len 54 DANG THU CUNG (5 byte: byte mobId=-1 roi
// int masterId) — dac ta §2/§3.1. KH6 chua lam he thu cung (dac ta §0 "Ngoai
// KH6"); test dung ham nay CHI de kiem may chu doc DU 5 byte, khong lech dong
// byte cho goi ke tiep tren cung ket noi (rui ro so mot cua KH6, dac ta §2).
func (cl *clientGia) guiDanhQuaiThuCung(masterId int32) {
	cl.t.Helper()
	w := protocol.NewWriter()
	w.Byte(255) // byte khong dau cua -1 (int8(255) == -1)
	w.Int(masterId)
	if err := cl.f.WriteMessage(protocol.NewMessage(lenhNguoiDanhQuai, w)); err != nil {
		cl.t.Fatalf("gui goi 54 (thu cung, 5 byte) that bai: %v", err)
	}
}

// doiGoiLenhKhongFatal giong het (cl *clientGia).doiGoiLenh (dongbo_tichhop_test.go)
// ve tinh than — doc toi khi gap dung lenh mong doi, bo qua goi khac, co han
// gio hanDoiGoi — nhung TRA LOI QUA ERROR thay vi goi t.Fatal, nen AN TOAN
// goi tu MOT GOROUTINE KHAC voi goroutine dang chay ham Test (xem canh bao ve
// t.Fatal/FailNow tu goroutine con o dau tichhop_test.go).
//
// Vi sao Task 7 can ban nay ma cac Task truoc (KH5) khong can: k.phatChoCaKhuVuc
// (khuvuc.go, moi tu Task 7) la CHO DAU TIEN mot lan phat ghi vao HAI HAY
// NHIEU ket noi cung luc (goi chien dau 54/-9/-12 gui ca nguoi danh LAN nguoi
// ngoai). net.Pipe() dong bo TUYET DOI (Write() cho toi khi dau kia Read()
// khop) — neu test doc TUAN TU tung ket noi mot (het A moi toi B), actor co
// the dang ket o Write() cho B doc goi 54 trong khi test dang ket o doc "-9"
// cua A (goi ma actor CHUA GUI, vi con dang ket o vong phat 54 cho B) -> BE
// TAC cho toi khi hanDoiGoi het han. Phai doc DONG THOI (goroutine rieng cho
// moi ket noi) de khop dung cach actor phat: ghi xen ke vao nhieu ket noi.
func (cl *clientGia) doiGoiLenhKhongFatal(lenh int8) (*protocol.Message, error) {
	k := cl.f.KeyStream()
	_ = cl.conn.SetReadDeadline(time.Now().Add(hanDoiGoi))
	defer cl.conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	for {
		m, err := docMotGoiThuCong(cl.conn, k)
		if err != nil {
			return nil, fmt.Errorf("khong tim thay goi lenh %d truoc khi het han %s: %w", lenh, hanDoiGoi, err)
		}
		if m.Command == lenh {
			return m, nil
		}
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

// ---------------------------------------------------------------------------
// Task 7 (Ke hoach 6): nguoi choi danh quai — noi dieu phoi goi len 54 (dac ta
// §3.1) toi su kien suKienDanhQuai (khuvuc.go), xu ly TRONG vongLap actor.
//
// Ca ba kich ban duoi day dung quai MAC DINH cua taoQuaiMacDinh (quai.go):
// Id=0 tai (bd.SpawnX+60, bd.SpawnY), Id=1 tai (bd.SpawnX-60, bd.SpawnY), ca
// hai Hp=HpToiDa=100. Nguoi choi lang 0 (La) spawn dung (bd.SpawnX, bd.SpawnY)
// = (100,384) (xem banDo.go/lang.go) nen ca hai quai deu cach dung 60 don vi —
// trong tam can chien 100 (chieuthuc.go, tamCanChien). Voi chieu can chien cap
// 1 cua lang 0 (templateId=1, slot 0): satThuongGoc = lang0.SatThuongKhoiDau =
// 12 (lang.go), chieu.SatThuong = dameGoc(slot0,cap1) = 20 (taoCapChieu),
// he nguyen to nguoi-vs-quai LUON trung tinh (heQuaiChuaXacDinh, khuvuc.go) ->
// satThuong = 12*20/100 = 2 (nguyen, khong lam tron). skillId (lang0 slot0
// cap1) = 1 — cong thuc giong het TestKhuVucChonChieuCapNhatChieuDangChon
// (khuvuc_test.go): skillIdGoc = 1 + (lang*4+slot)*3 = 1 + (0*4+0)*3 = 1, cap
// 1 -> SkillId = skillIdGoc + cap - 1 = 1.
// ---------------------------------------------------------------------------

// TestHeNguyenToNguoiChoiAnhXaDungLang kiem tra bang anh xa Lang.ID -> he
// nguyen to (heNguyenToTheoLang, khuvuc.go) KHONG bi suy dien nham "lang i =
// he i" — hai lang (2: Suong Mu, 4: May) co thu tu LECH so voi hang so he
// nguyen to (satthuong.go). Gia tri ghim cung, doi chieu truc tiep voi thiet
// ke tong §7.2 va chinh comment trong chieuthuc.go (tenChieuTheoLang: "Lang
// Suong Mu — Thuy", "Lang May — Loi").
//
// Test rieng nay CAN THIET du khong nam trong 3 kich ban brief Task 7 doi
// hoi: ca ba test tich hop o duoi day CHI dung nhan vat lang 0 (La=HeHoa,
// PHAN TU DAU cua mang) — mot ban do lech kieu "identity" sai (vd
// {HeHoa,HePhong,HeLoi,HeTho,HeThuy}, khong hoan doi Suong Mu/May) van cho
// lang 0 dung, nen se KHONG bi ba test tich hop bat duoc.
//
// Luu y: mapping nay HIEN chua anh huong byte nao cua Task 7 (heQuaiChuaXacDinh
// lam quai luon "trung tinh" bat ke he nguoi danh la gi — xem chu thich
// heQuaiChuaXacDinh) — test nay chot dung bang tra CHO CAC TASK SAU (vd Task
// 11 nguoi-vs-nguoi, ca hai phia deu co he nguyen to that) dung lai ham nay
// ma khong can viet lai/soat lai.
func TestHeNguyenToNguoiChoiAnhXaDungLang(t *testing.T) {
	truongHop := []struct {
		lang int16
		muon byte
	}{
		{0, HeHoa},   // Lang La
		{1, HePhong}, // Lang Cat
		{2, HeThuy},  // Lang Suong Mu — LECH: he thu 5 trong ngu hanh (iota 4), KHONG phai lang thu 3
		{3, HeTho},   // Lang Da
		{4, HeLoi},   // Lang May — LECH: he thu 3 trong ngu hanh (iota 2), KHONG phai lang thu 5
	}
	for _, tc := range truongHop {
		if got := heNguyenToNguoiChoi(tc.lang); got != tc.muon {
			t.Errorf("heNguyenToNguoiChoi(lang=%d) = %d, mong doi %d", tc.lang, got, tc.muon)
		}
	}
}

// TestHeNguyenToNguoiChoiLangHongTraTrungTinh: lang ngoai [0,4] (du lieu
// hong, ve ly thuyet khong xay ra vi TimLang da chan luc tao nhan vat) phai
// tra heQuaiChuaXacDinh (>= soHeNguyenTo) de heSoNguyenToPhanTram tu nhien coi
// la trung tinh, khong panic index-out-of-range.
func TestHeNguyenToNguoiChoiLangHongTraTrungTinh(t *testing.T) {
	if got := heNguyenToNguoiChoi(-1); got != heQuaiChuaXacDinh {
		t.Errorf("heNguyenToNguoiChoi(-1) = %d, mong doi %d (heQuaiChuaXacDinh)", got, heQuaiChuaXacDinh)
	}
	if got := heNguyenToNguoiChoi(5); got != heQuaiChuaXacDinh {
		t.Errorf("heNguyenToNguoiChoi(5) = %d, mong doi %d (heQuaiChuaXacDinh)", got, heQuaiChuaXacDinh)
	}
}

// TestTichHop_DanhQuaiConSong la kich ban 1 cua Task 7: hai client THAT trong
// CUNG khu vuc (actor THAT, BatDau() qua MoSoKhuVuc — dung y het cac test
// khac trong file nay); A chon chieu roi danh quai Id=0. CA HAI A (nguoi
// danh) VA B (nguoi ngoai) deu phai nhan dung chuoi 54 roi -9 — dac ta §3.3
// cot "Gui cho: Ca khu vuc", brief Task 7 nhan manh ro "ke ca nguoi danh".
func TestTichHop_DanhQuaiConSong(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	soKhuVuc := MoSoKhuVuc()
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTKA := "danhA" + hauToNganNgau(t)
	tenTKB := "danhB" + hauToNganNgau(t)
	taiKhoanA := taoTaiKhoanCoMatKhau(t, tenTKA, matKhau)
	taiKhoanB := taoTaiKhoanCoMatKhau(t, tenTKB, matKhau)
	tenNVA := "DanhA" + hauToNganNgau(t)
	tenNVB := "DanhB" + hauToNganNgau(t)

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
	khangDinhIdTrongGoi(t, clB.doiGoiLenh(t, -5), nvA.ID)
	khangDinhIdTrongGoi(t, clA.doiGoiLenh(t, -5), nvB.ID)

	// A chon chieu can chien (templateId=1, tam 100) — khong phat goi gi ca
	// (dac ta §3.2: chon chieu la trang thai noi bo, khong co goi xuong).
	clA.chonChieu(1)
	if !clA.khongConGoiGi(t) {
		t.Fatal("chon chieu (34) khong duoc phat goi tin nao")
	}

	// A danh quai Id=0 — trong tam, chieu da chon, chua tung danh -> hoi chieu
	// tu nhien qua (lanDanhCuoi la zero value).
	clA.guiDanhQuai(0)

	// kiemChuoi54Roi9 doc tu MOT ket noi va tra loi qua error (KHONG t.Fatal)
	// — bat buoc chay tren GOROUTINE RIENG cho tung client (xem chu thich
	// doiGoiLenhKhongFatal ve ly do: k.phatChoCaKhuVuc ghi xen ke vao ca hai
	// ket noi, doc tuan tu se be tac).
	kiemChuoi54Roi9 := func(cl *clientGia) error {
		goi54, err := cl.doiGoiLenhKhongFatal(lenhNguoiDanhQuai)
		if err != nil {
			return err
		}
		r := protocol.NewReader(goi54.Data)
		idDanh, errID := r.Int()
		if errID != nil || idDanh != int32(nvA.ID) {
			return fmt.Errorf("goi 54 idNguoiDanh = %d (loi %v), mong doi %d", idDanh, errID, nvA.ID)
		}
		skillId, errSk := r.Byte()
		if errSk != nil || skillId != 1 {
			return fmt.Errorf("goi 54 skillId = %d (loi %v), mong doi 1 (lang0 slot0 cap1)", skillId, errSk)
		}
		mobId, errMob := r.Byte()
		if errMob != nil || mobId != 0 {
			return fmt.Errorf("goi 54 mobId = %d (loi %v), mong doi 0", mobId, errMob)
		}

		goi9, err := cl.doiGoiLenhKhongFatal(-9)
		if err != nil {
			return err
		}
		muon := []byte{
			0,           // mobId
			0, 0, 0, 98, // hpConLai = 100 - 2 = 98
			0, 0, 0, 2, // satThuong = 2 (12 * 20 / 100, xem chu thich dau nhom test)
			0,   // chiMang = false
			255, // hieuUngId = -1 (khong hieu ung, dac ta §3.3)
		}
		if !bytes.Equal(goi9.Data, muon) {
			return fmt.Errorf("goi -9 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
				goi9.Data, len(goi9.Data), muon, len(muon))
		}
		return nil
	}

	type ketQua struct {
		ten string
		err error
	}
	ketQuaCh := make(chan ketQua, 2)
	go func() { ketQuaCh <- ketQua{ten: "A (nguoi danh)", err: kiemChuoi54Roi9(clA)} }()
	go func() { ketQuaCh <- ketQua{ten: "B (nguoi ngoai)", err: kiemChuoi54Roi9(clB)} }()
	for i := 0; i < 2; i++ {
		kq := <-ketQuaCh
		if kq.err != nil {
			t.Fatalf("%s: %v", kq.ten, kq.err)
		}
	}
}

// bomSuKienKhongPhatGoi bom TAY dung mot su kien da lay tu hop thu qua
// k.xuLy, CHAY TRONG MOT GOROUTINE RIENG roi CHO no chay xong HOAN TOAN truoc
// khi tra ve. Dung "close(xong)" lam diem dong bo — theo mo hinh bo nho Go,
// day la quan he "happens-before" hop le nen doc k.quai/k.nguoiTrong NGAY SAU
// khi ham nay tra ve la AN TOAN (khong dua tranh voi goroutine con, no da ket
// thuc han truoc do).
//
// ⚠️ CHI AN TOAN DUNG cho su kien nao KHONG lam k.xuLy goi Gui() (goi tin qua
// mang) ben trong — o Task 7 la suKienVao (khi khu vuc dang CHI CO mot minh
// nguoi vua vao, phatChoNguoiKhac() la mot vong lap RONG) va suKienChonChieu
// (khong co dong Gui() nao trong nhanh do, xem khuvuc.go). TUYET DOI KHONG
// dung cho suKienDanhQuai (co the phat 54/-9/-12 qua ket noi that): Gui() se
// BLOCK cho toi khi dau kia Read() — va o day khong ai dang doc trong luc
// goroutine con con dang cho ben trong ham nay -> BE TAC VINH VIEN. Voi
// suKienDanhQuai, goi k.xuLy trong "go func(){...}" (fire-and-forget) ROI DOC
// NGAY goi tra loi tu ket noi — chinh thao tac doc do la diem dong bo, giong
// het cach actor THAT (goroutine vongLap rieng) hoat dong o cac test khac.
func bomSuKienKhongPhatGoi(k *KhuVuc, sk suKienKhuVuc) {
	xong := make(chan struct{})
	go func() {
		k.xuLy(sk)
		close(xong)
	}()
	<-xong
}

// TestTichHop_DanhQuaiChet la kich ban 2 cua Task 7: mot client THAT danh
// CHET mot con quai -> phai nhan dung chuoi 54 roi -12 (dac ta §6.2).
//
// Ky thuat: khu vuc RIENG cho test nay (khong dung MoSoKhuVuc) va KHONG goi
// k.BatDau() — khong co goroutine vongLap nao chay song song. Nho vay, sau
// khi bom xong hai su kien khong phat goi (suKienVao/suKienChonChieu) qua
// bomSuKienKhongPhatGoi (an toan, xem chu thich ham do), test goroutine la
// NGUOI DUY NHAT dang dung k.quai — ha HP quai xuong 1 TRUC TIEP o day KHONG
// vi pham ky luat actor (dac ta §1.6 cam DUA TRANH giua NHIEU goroutine; o
// day tai moi thoi diem chi co DUNG MOT). Day la cach duy nhat lam "quai
// chet" NHANH trong mot test don vi: sat thuong cap 1 hien chi ~2/don (xem
// nhom chu thich dau file), can hang chuc don that (hoac hang chuc nguoi choi
// that) moi ha du 100 mau mac dinh — qua cham/qua cong kenh cho muc dich cua
// test nay (kiem CHUOI GOI khi quai chet, khong phai kiem tich luy sat
// thuong qua nhieu don — cai do da co test rieng o TestTichHop_DanhQuaiConSong).
func TestTichHop_DanhQuaiChet(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	k := MoKhuVuc(TimBanDo(1))
	k.hopThu = make(chan suKienKhuVuc, 8) // KHONG goi k.BatDau() — xem chu thich ham

	soKhuVuc := &SoKhuVuc{theoMap: map[int32]*KhuVuc{1: k}}

	matKhau := "matkhau123"
	tenTK := "killA" + hauToNganNgau(t)
	taiKhoanA := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "KillA" + hauToNganNgau(t)

	clA := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer clA.dong()
	dangNhapTaoNhanVatVaoTheGioi(t, clA, tenTK, matKhau, tenNV)

	nvA, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoanA)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan A that bai: %v", err)
	}
	if nvA.MapID != 1 {
		t.Fatalf("nhan vat phai o mapId 1 (lang 0), nhan duoc %d", nvA.MapID)
	}

	// dangNhapTaoNhanVatVaoTheGioi da kich hoat xuLyClientOk (buoc -28/13),
	// va ham do GOI DanhSachQuaiSnapshot() khi dung -24 MAP INFO -> k.Gui(
	// suKienXinDanhSachQuai{...}). Vi khu vuc nay KHONG BatDau() nen KHONG ai
	// tieu thu su kien do — no con NAM YEN o DAU k.hopThu. Da doc thanh cong
	// goi -24 (ben trong dangNhapTaoNhanVatVaoTheGioi) tuc DanhSachQuaiSnapshot
	// da goi k.Gui(...) XONG (cung mot goroutine PhucVu, tuan tu: goi -82 ->
	// -30 -> DanhSachQuaiSnapshot() -> goi -24), nen xa NGAY su kien nay o day
	// la AN TOAN VA CHAC CHAN (khong phai doan mo, la quan he "happens-before"
	// qua viec doc goi -24 thanh cong) — neu khong xa, buoc doc suKienVao duoi
	// day se boc NHAM phai su kien nay truoc.
	<-k.hopThu

	// A vao khu vuc RIENG nay MOT MINH -> phatChoNguoiKhac ben trong xuLy la
	// vong lap RONG (khong co "nguoi khac" nao) -> KHONG Gui() nao chay -> an
	// toan tuyet doi cho bomSuKienKhongPhatGoi.
	clA.baoVaoKhuVuc()
	bomSuKienKhongPhatGoi(k, <-k.hopThu)
	if _, co := k.nguoiTrong[nvA.ID]; !co {
		t.Fatal("A phai co trong k.nguoiTrong sau khi xu ly suKienVao")
	}

	// A chon chieu can chien (templateId=1) — nhanh suKienChonChieu KHONG co
	// Gui() nao (xem khuvuc.go) -> an toan cho bomSuKienKhongPhatGoi.
	clA.chonChieu(1)
	bomSuKienKhongPhatGoi(k, <-k.hopThu)
	if k.nguoiTrong[nvA.ID].chieuDangChon == nil {
		t.Fatal("A phai co chieuDangChon sau khi xu ly suKienChonChieu")
	}

	// AN TOAN sua truc tiep k.quai o day: xem chu thich ham TestTichHop_DanhQuaiChet
	// va bomSuKienKhongPhatGoi o tren — khong co goroutine thu hai nao dung
	// toi k.quai tai thoi diem nay. Ha HP quai Id=0 xuong 1 (GIU NGUYEN
	// HpToiDa=100): vi Hp(1) != HpToiDa(100), quy tac "quai day mau khong
	// chet boi mot don" (dac ta §1.3) KHONG can thiep — bat ky don hop le nao
	// (>= 1 sat thuong) cung du giet no ngay, dung trong tam cua test nay.
	k.quai[0].Hp = 1

	clA.guiDanhQuai(0)
	// suKienDanhQuai CO THE phat Gui() qua ket noi that (54 roi -12) -> KHONG
	// dung bomSuKienKhongPhatGoi (se be tac, xem chu thich ham do). Goi
	// k.xuLy trong goroutine rieng (fire-and-forget) roi DOC NGAY goi tra loi
	// — chinh thao tac doc la diem dong bo.
	go k.xuLy(<-k.hopThu)

	goi54 := clA.doiGoiLenh(t, lenhNguoiDanhQuai)
	r := protocol.NewReader(goi54.Data)
	idDanh, errID := r.Int()
	if errID != nil || idDanh != int32(nvA.ID) {
		t.Fatalf("goi 54: idNguoiDanh = %d (loi %v), mong doi %d", idDanh, errID, nvA.ID)
	}
	skillId, errSk := r.Byte()
	if errSk != nil || skillId != 1 {
		t.Fatalf("goi 54: skillId = %d (loi %v), mong doi 1 (lang0 slot0 cap1)", skillId, errSk)
	}
	mobId, errMob := r.Byte()
	if errMob != nil || mobId != 0 {
		t.Fatalf("goi 54: mobId = %d (loi %v), mong doi 0", mobId, errMob)
	}

	goi12 := clA.doiGoiLenh(t, -12)
	muon := []byte{
		0,          // mobId
		0, 0, 0, 1, // satThuong = 1 (dame goc 2 nhung quai chi con 1 mau -> kep theo §1.2)
		0, // chiMang = false
		0, // soDoRoi = 0 (KH6 chua co he vat pham)
	}
	if !bytes.Equal(goi12.Data, muon) {
		t.Fatalf("goi -12 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			goi12.Data, len(goi12.Data), muon, len(muon))
	}
}

// TestTichHop_DanhQuaiThuCungDocDu5ByteKhongLechDong la kich ban 3 (BAT BUOC
// theo brief Task 7) — goi len 54 la goi DO DAI BIEN THIEN (dac ta §2), rui
// ro so mot cua Ke hoach 6: mobId=-1 (thu cung, ngoai pham vi KH6, dac ta §0)
// mang THEM 4 byte masterId, tong 5 byte. May chu PHAI doc DU ca 5 byte, neu
// khong goi KE TIEP tren CUNG ket noi se bi hieu sai.
//
// Kich ban: A gui 54 dang thu cung (5 byte) -> mobId=255 khong khop quai nao
// trong khu vuc (Quai.Id luon 0..127) -> may chu phai IM LANG, KHONG phat goi
// nao (dung nhu KH6 chua co he thu cung). NGAY SAU DO, tren CUNG ket noi, A
// gui MOT goi 54 THUONG (1 byte, mobId=1 — con quai mac dinh con lai) -> PHAI
// nhan dung chuoi 54 roi -9 nhu binh thuong. Neu server doc THIEU/THUA byte o
// buoc truoc, goi nay se hoac khong toi, hoac toi voi du lieu sai lech —
// chinh la dieu test nay bat duoc.
func TestTichHop_DanhQuaiThuCungDocDu5ByteKhongLechDong(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	soKhuVuc := MoSoKhuVuc()
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTK := "petA" + hauToNganNgau(t)
	taiKhoanA := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "PetA" + hauToNganNgau(t)

	clA := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer clA.dong()
	dangNhapTaoNhanVatVaoTheGioi(t, clA, tenTK, matKhau, tenNV)

	nvA, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoanA)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan A that bai: %v", err)
	}
	if nvA.MapID != 1 {
		t.Fatalf("nhan vat phai o mapId 1 (lang 0), nhan duoc %d", nvA.MapID)
	}

	clA.baoVaoKhuVuc()
	if !clA.khongConGoiGi(t) {
		t.Fatal("A vao khu vuc mot minh -> khong duoc nhan goi nao (chua co ai khac)")
	}

	// Chon chieu can chien (tam 100) — quai Id=1 mac dinh o bd.SpawnX-60, cach
	// A dung 60 don vi, trong tam.
	clA.chonChieu(1)
	if !clA.khongConGoiGi(t) {
		t.Fatal("chon chieu (34) khong duoc phat goi tin nao")
	}

	// Goi 54 DANG THU CUNG: 5 byte. KHONG co quai nao Id=255 trong khu vuc ->
	// may chu PHAI im lang, KHONG duoc sap ket noi hay tra loi gi ca.
	clA.guiDanhQuaiThuCung(999999)
	if !clA.khongConGoiGi(t) {
		t.Fatal("goi 54 dang thu cung (mobId=-1) phai bi lang le bo qua (KH6 chua co he thu cung), khong duoc phat goi nao")
	}

	// NGAY SAU DAY, tren CUNG ket noi: goi 54 THUONG, mobId=1 -> phai nhan
	// dung chuoi 54 roi -9 — day chinh la bang chung ket noi KHONG bi lech
	// dong byte boi goi thu cung truoc do.
	clA.guiDanhQuai(1)

	goi54 := clA.doiGoiLenh(t, lenhNguoiDanhQuai)
	r := protocol.NewReader(goi54.Data)
	idDanh, errID := r.Int()
	if errID != nil || idDanh != int32(nvA.ID) {
		t.Fatalf("goi 54 (sau goi thu cung): idNguoiDanh = %d (loi %v), mong doi %d — dau hieu LECH DONG BYTE", idDanh, errID, nvA.ID)
	}
	skillId, errSk := r.Byte()
	if errSk != nil || skillId != 1 {
		t.Fatalf("goi 54 (sau goi thu cung): skillId = %d (loi %v), mong doi 1 — dau hieu LECH DONG BYTE", skillId, errSk)
	}
	mobId, errMob := r.Byte()
	if errMob != nil || mobId != 1 {
		t.Fatalf("goi 54 (sau goi thu cung): mobId = %d (loi %v), mong doi 1 — dau hieu LECH DONG BYTE", mobId, errMob)
	}

	goi9 := clA.doiGoiLenh(t, -9)
	muon := []byte{
		1,           // mobId
		0, 0, 0, 98, // hpConLai = 100 - 2 = 98
		0, 0, 0, 2, // satThuong = 2
		0,   // chiMang = false
		255, // hieuUngId = -1
	}
	if !bytes.Equal(goi9.Data, muon) {
		t.Fatalf("goi -9 (sau goi thu cung) byte sai — dau hieu LECH DONG BYTE:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			goi9.Data, len(goi9.Data), muon, len(muon))
	}
}
