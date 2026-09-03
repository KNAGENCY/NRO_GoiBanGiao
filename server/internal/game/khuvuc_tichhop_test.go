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
// Ca "ngoai collision bi chan" va "roi di phat -6" CO CHU DICH de danh cho Task 8
// (xem ghi chu cuoi task-7-brief-kh5.md) — khong lap lai o day.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	mynet "nroserver/internal/net"
	"nroserver/internal/protocol"
	"nroserver/internal/store"
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

	return &clientGia{t: t, f: f, conn: benClient, b: b}
}

// moSoKhuVucKhongQuaiMacDinh xay 5 khu vuc GIONG HET MoSoKhuVuc (san xuat
// that) NHUNG XOA SACH quai mac dinh (taoQuaiMacDinh, quai.go) TRUOC KHI
// BatDau() — Ke hoach 6 Task 9 them ham nay VI: AI quai (aiQuai, khuvuc.go)
// nay chay tren TICKER THAT ma MoSoKhuVuc() khoi dong, va nhan vat MOI TAO
// luon spawn dung diem SpawnX/SpawnY cua ban do (lang.go: ViTriXKhoiDau/
// ViTriYKhoiDau) — CACH DUNG 60 don vi voi CA HAI quai mac dinh
// (taoQuaiMacDinh), nam TRONG tam tang 2 cua AI (100, dac ta §5.2). Cac test
// KHONG lien quan gi toi quai (di chuyen, vao/roi khu vuc) se bi quai TU DONG
// DANH ngay khi nhan vat vua vao khu vuc mot minh, pha vo gia dinh "khong ai
// gui goi gi neu khong co hanh dong nao cua nguoi choi that".
//
// AN TOAN tuyet doi de ghi truc tiep k.quai = nil o day (khong dua tranh du
// lieu voi goroutine vongLap): thoi diem nay goroutine ay CHUA duoc tao —
// BatDau() (dong goi tao goroutine) duoc goi SAU dong gan k.quai = nil, dung
// ky thuat "AN TOAN sua truc tiep k.quai" da dung o TestTichHop_DanhQuaiChet,
// chi khac o CHO trong vong doi (truoc BatDau(), khong phai sau khi bom su
// kien tay).
//
// Cac test "danh quai" (TestTichHop_DanhQuai...) VAN can quai that de danh ->
// KHONG dung ham nay, GIU MoSoKhuVuc() nhu cu va tu xu ly nhieu AI rieng bang
// khongConGoiGiTruQuaiDanh (xem ham do o duoi).
func moSoKhuVucKhongQuaiMacDinh(t *testing.T) *SoKhuVuc {
	t.Helper()
	s := &SoKhuVuc{theoMap: make(map[int32]*KhuVuc)}
	for mapId := int32(1); mapId <= 5; mapId++ {
		k := MoKhuVuc(TimBanDo(mapId))
		k.quai = nil // AN TOAN: goroutine vongLap CHUA chay (BatDau() o duoi)
		k.BatDau()
		s.theoMap[mapId] = k
	}
	return s
}

// khongConGoiGiTruQuaiDanh giong het khongConGoiGi (tichhop_test.go) NHUNG
// BO QUA (ghi log, KHONG that bai) cac goi -11/-10 roi doc tiep — Ke hoach 6
// Task 9: cac test "danh quai" trong file nay CO CHU DICH GIU quai mac dinh
// (taoQuaiMacDinh) that de lam doi tuong danh, ma quai do dung CACH 60 don vi
// diem spawn — nam TRONG tam tang 2 cua AI (100, dac ta §5.2) nen co the TU
// DONG danh nhan vat NGAY khi vua vao khu vuc, truoc ca khi test kip lam gi.
// Day la NHIEU NEN HOP LE tu tinh nang AI that (khong phai loi), khac han cac
// lenh khac (van phai that su im lang tuyet doi).
//
// CHI dung ham nay o nhung diem ket noi VAN DANG DUOC DOC (client con song
// binh thuong) — Gui() phia actor se hoan tat NHANH (khong co rui ro chan lau
// toi thoiGianChoGui=30s nhu khi ket noi sap dong, xem session.go). KHONG dung
// cho kich ban dong ket noi giua chung (TestTichHop_RoiKhuVucPhatGoi6) — noi
// do dung moSoKhuVucKhongQuaiMacDinh de loai bo quai HOAN TOAN thay vi khoan
// dung nhieu, vi mot ket noi dang dong co the khien actor CHAN THAT SU toi
// 30s moi lan Gui() that bai (da do duoc: 120s cho 3 lan thu tren mot khu
// vuc that khi chua co ham nay).
func (cl *clientGia) khongConGoiGiTruQuaiDanh(t *testing.T) bool {
	t.Helper()
	for {
		_ = cl.conn.SetReadDeadline(time.Now().Add(hanChoImLang))
		m, err := cl.f.ReadMessage()
		if err != nil {
			cl.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				return true // het han, khong con goi nao (ke ca nhieu AI) -> thuc su im lang
			}
			t.Fatalf("ket noi gap loi KHONG PHAI het han gio khi dang cho xac nhan im lang: %v", err)
			return false
		}
		if m.Command == -11 || m.Command == -10 {
			t.Logf("bo qua goi quai AI (nhieu hop le, Ke hoach 6 Task 9): lenh=%d du_lieu=%v", m.Command, m.Data)
			continue // doc tiep voi han gio MOI — quai co the danh nhieu lan trong khoang cho
		}
		cl.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		t.Logf("nhan duoc goi NGOAI DU KIEN trong luc cho im lang: lenh=%d du_lieu=%v", m.Command, m.Data)
		return false
	}
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

// guiHoiSinh gui goi len -16 than RONG (dac ta §3.1, Ke hoach 6 Task 13):
// client xin hoi sinh sau khi chet.
func (cl *clientGia) guiHoiSinh() {
	cl.t.Helper()
	if err := cl.f.WriteMessage(&protocol.Message{Command: lenhHoiSinh}); err != nil {
		cl.t.Fatalf("gui goi -16 (xin hoi sinh) that bai: %v", err)
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
	doiChuoiCatalogItem(t, cl) // BA message: subtype 0, 1, 2 (KH7 task 2)

	cl.guiClientOk()
	cl.doiGoiLenh(t, -82)
	cl.doiGoiLenh(t, -30)
	cl.doiGoiLenh(t, -24)
	// -42 CHI SO gui NGAY SAU MAP INFO (login.go xuLyClientOk). Khong doc thi ong
	// pipe day, may chu ket thuc bang loi GHI het gio 30 giay — trong nhu treo,
	// khong phai loi logic. Da lam do 13 test tich hop vi thieu dung dong nay.
	cl.doiGoiLenh(t, lenhChiSoNhanVat)
	// KH7 Task 5: ba goi hanh trang -37/-36/-35 khep lai chuoi vao the gioi
	// (login.go xuLyClientOk, gui SAU -42). Khong doc thi moi khang dinh "im
	// lang" phia sau se thay ba goi nay va bao that bai — cung ly do voi -42 o
	// tren.
	cl.doiGoiLenh(t, lenhDoDangMac)
	cl.doiGoiLenh(t, lenhTui)
	cl.doiGoiLenh(t, lenhRuong)
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

	// Task 9: khong dung MoSoKhuVuc() — test nay khong lien quan gi toi quai,
	// dung ban khong quai de tranh AI tu danh nhan vat dang doc (xem chu
	// thich moSoKhuVucKhongQuaiMacDinh).
	soKhuVuc := moSoKhuVucKhongQuaiMacDinh(t)
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
	toX := int16(nvA.ViTriX) - 10
	toY := int16(nvA.ViTriY) + 10
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

// TestTichHop_DiChuyenNgoaiCollisionBiChan la kich ban 1 cua Task 8: hai phien
// THAT (giong het TestTichHop_HaiNguoiThayNhauVaDiChuyen) da thay nhau qua
// -5, roi A gui -7 voi dich nam ngoai collision artifact cua map 1. Actor
// (kiemDiChuyenHopLe trong khuvuc.go, case
// suKienDiChuyen) phai BO goi nay va KHONG phat cho ai, nen B khong duoc
// nhan goi -7 nao trong han cho im lang. Sau do A gui tiep MOT -7 HOP LE
// (dich gan, trong nguong) de chung minh ket noi VAN SONG BINH THUONG — chi
// goi bat hop le bi bo, may chu khong ngat ket noi vi no.
func TestTichHop_DiChuyenNgoaiCollisionBiChan(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	// Task 9: xem chu thich moSoKhuVucKhongQuaiMacDinh — test nay khong lien
	// quan gi toi quai.
	soKhuVuc := moSoKhuVucKhongQuaiMacDinh(t)
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

	// A gui -7 toi (2000,2000), nam ngoai footprint collision 120x408 cua map
	// 1. Collision artifact la nguon bien authoritative, nen actor phai bo goi.
	toXQuaXa := int16(2000)
	toYQuaXa := int16(2000)
	clA.guiDiChuyen(toXQuaXa, toYQuaXa)

	// Khang dinh B KHONG nhan goi -7 nao trong han cho im lang: actor bo goi
	// bat hop le TRUOC khi phat (return som trong case suKienDiChuyen), nen
	// khong ai — ke ca B — nhan duoc gi.
	if !clB.khongConGoiGi(t) {
		t.Fatal("A gui -7 ngoai collision -> B KHONG duoc nhan goi -7 nao (may chu phai bo goi bat hop le)")
	}

	// Chung minh ket noi CUA A VAN SONG BINH THUONG (khong bi ngat vi goi bat
	// hop le): gui tiep MOT -7 HOP LE (dich gan, trong nguong) -> B phai nhan
	// dung goi -7 chuyen tiep nhu kich ban duong hanh phuc.
	toXHopLe := int16(nvA.ViTriX) - 10
	toYHopLe := int16(nvA.ViTriY) + 10
	clA.guiDiChuyen(toXHopLe, toYHopLe)

	goiDiChuyenChoB := clB.doiGoiLenh(t, -7)
	r := protocol.NewReader(goiDiChuyenChoB.Data)
	if id, errID := r.Int(); errID != nil || id != int32(nvA.ID) {
		t.Fatalf("id trong goi -7 (hop le, sau khi goi ngoai collision bi bo) = %d (loi %v), mong doi %d", id, errID, nvA.ID)
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

	// Task 9: BAT BUOC dung ban khong quai — test nay DONG ket noi cua A giua
	// chung (xem clA.dong() ben duoi). Neu quai that con trong khu vuc, AI co
	// the dang thu Gui() toi A dung luc do -> chan actor toi 30s/lan (da do
	// duoc that: 120s cho 3 lan thu truoc khi co ham nay). Xem chu thich
	// moSoKhuVucKhongQuaiMacDinh.
	soKhuVuc := moSoKhuVucKhongQuaiMacDinh(t)
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
	// Task 9: khongConGoiGiTruQuaiDanh (khong phai khongConGoiGi) — quai mac
	// dinh dung canh diem spawn (tam tang 2 AI, dac ta §5.2) co the tu danh A
	// ngay khi vua vao, day la nhieu hop le (xem chu thich ham).
	if !clA.khongConGoiGiTruQuaiDanh(t) {
		t.Fatal("A vao khu vuc mot minh -> khong duoc nhan goi nao NGOAI quai AI (chua co ai khac)")
	}
	clB.baoVaoKhuVuc()
	khangDinhIdTrongGoi(t, clB.doiGoiLenh(t, -5), nvA.ID)
	khangDinhIdTrongGoi(t, clA.doiGoiLenh(t, -5), nvB.ID)

	// A chon chieu can chien (templateId=1, tam 100) — khong phat goi gi ca
	// (dac ta §3.2: chon chieu la trang thai noi bo, khong co goi xuong).
	clA.chonChieu(1)
	if !clA.khongConGoiGiTruQuaiDanh(t) {
		t.Fatal("chon chieu (34) khong duoc phat goi tin nao NGOAI quai AI")
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
			0, 0, 0, 89, // hpConLai = 100 - 11 = 89
			0, 0, 0, 11, // satThuong = 11 (55 * 20 / 100, xem chu thich dau nhom test)
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
	// Task 9: khongConGoiGiTruQuaiDanh — quai mac dinh (Id=1) trong tam tang 2
	// AI co the tu danh A ngay, xem chu thich ham.
	if !clA.khongConGoiGiTruQuaiDanh(t) {
		t.Fatal("A vao khu vuc mot minh -> khong duoc nhan goi nao NGOAI quai AI (chua co ai khac)")
	}

	// Chon chieu can chien (tam 100) — quai Id=1 mac dinh o bd.SpawnX-60, cach
	// A dung 60 don vi, trong tam.
	clA.chonChieu(1)
	if !clA.khongConGoiGiTruQuaiDanh(t) {
		t.Fatal("chon chieu (34) khong duoc phat goi tin nao NGOAI quai AI")
	}

	// Goi 54 DANG THU CUNG: 5 byte. KHONG co quai nao Id=255 trong khu vuc ->
	// may chu PHAI im lang, KHONG duoc sap ket noi hay tra loi gi ca.
	clA.guiDanhQuaiThuCung(999999)
	if !clA.khongConGoiGiTruQuaiDanh(t) {
		t.Fatal("goi 54 dang thu cung (mobId=-1) phai bi lang le bo qua (KH6 chua co he thu cung) NGOAI quai AI, khong duoc phat goi nao khac")
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
		0, 0, 0, 89, // hpConLai = 100 - 11 = 89
		0, 0, 0, 11, // satThuong = 11
		0,   // chiMang = false
		255, // hieuUngId = -1
	}
	if !bytes.Equal(goi9.Data, muon) {
		t.Fatalf("goi -9 (sau goi thu cung) byte sai — dau hieu LECH DONG BYTE:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			goi9.Data, len(goi9.Data), muon, len(muon))
	}
}

// ---------------------------------------------------------------------------
// Vong sua 1 (review coordinator, 27/07): hai luoi an toan bi bo sot o vong
// dau — logic da dung nhung khong test nao thuc su kich hoat nhanh do.
// ---------------------------------------------------------------------------

// TestTichHop_DanhQuaiDayMauKhongChetBoiMotDon la Finding 1: quy tac §1.3
// "quai dang DAY MAU khong chet boi MOT don" (khuvuc.go, nhanh
// `if mucTieu.Hp == mucTieu.HpToiDa && satThuong >= mucTieu.Hp { satThuong =
// mucTieu.Hp - 1 }`) TRUOC DO khong test nao thuc su kich hoat: Test 1
// (TestTichHop_DanhQuaiConSong) danh voi sat thuong nho hon nhieu HP; Test 2
// (TestTichHop_DanhQuaiChet) CO CHU DICH ha Hp=1 nhung GIU HpToiDa=100 (!=Hp)
// de NE dung nhanh nay. Neu ai doi `>=` thanh `>` hay xoa han nhanh, khong
// test nao truoc do phat hien.
//
// Ky thuat giong het TestTichHop_DanhQuaiChet (khu vuc RIENG, KHONG BatDau(),
// bom tay su kien) — nhung dat CA Hp VA HpToiDa cua quai Id=0 DEU = 2 (thay
// vi chi Hp=1) de dieu kien "Hp == HpToiDa" cua §1.3 THUC SU dung. Sat thuong
// tho tinh duoc (lang0 slot0 cap1) la 2 — dung bang Hp/HpToiDa nen kich hoat
// DUNG nhanh §1.3: quai phai CON SONG voi 1 mau, nhan `-9` (KHONG PHAI
// `-12`), satThuong bao xuong = 1 (khong phai 2 tho).
func TestTichHop_DanhQuaiDayMauKhongChetBoiMotDon(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	k := MoKhuVuc(TimBanDo(1))
	k.hopThu = make(chan suKienKhuVuc, 8) // KHONG goi k.BatDau() — xem TestTichHop_DanhQuaiChet

	soKhuVuc := &SoKhuVuc{theoMap: map[int32]*KhuVuc{1: k}}

	matKhau := "matkhau123"
	tenTK := "fullA" + hauToNganNgau(t)
	taiKhoanA := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "FullA" + hauToNganNgau(t)

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

	// Vet su kien suKienXinDanhSachQuai con sot lai tu xuLyClientOk (xem chu
	// thich chi tiet o TestTichHop_DanhQuaiChet).
	<-k.hopThu

	clA.baoVaoKhuVuc()
	bomSuKienKhongPhatGoi(k, <-k.hopThu)
	if _, co := k.nguoiTrong[nvA.ID]; !co {
		t.Fatal("A phai co trong k.nguoiTrong sau khi xu ly suKienVao")
	}

	clA.chonChieu(1)
	bomSuKienKhongPhatGoi(k, <-k.hopThu)
	if k.nguoiTrong[nvA.ID].chieuDangChon == nil {
		t.Fatal("A phai co chieuDangChon sau khi xu ly suKienChonChieu")
	}

	// AN TOAN sua truc tiep k.quai o day (khong co goroutine thu hai nao
	// dung toi, xem TestTichHop_DanhQuaiChet). DAT CA HAI Hp va HpToiDa = 2
	// — khac ky thuat cua TestTichHop_DanhQuaiChet (chi Hp=1, HpToiDa=100
	// giu nguyen) de dung THUC SU kich hoat nhanh "day mau" cua §1.3.
	k.quai[0].Hp = 2
	k.quai[0].HpToiDa = 2

	clA.guiDanhQuai(0)
	go k.xuLy(<-k.hopThu) // co the phat Gui() qua mang that -> fire-and-forget, xem TestTichHop_DanhQuaiChet

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

	// PHAI la -9 (quai con song), KHONG duoc la -12 (quai chet) — day la
	// khang dinh cot loi cua Finding 1: neu nhanh "day mau khong chet boi
	// mot don" bi xoa/sai, quai se chet that (Hp 2-2=0) va gui -12 thay vi
	// -9, lam doiGoiLenh(t, -9) duoi day cho het han 5s roi t.Fatal (do -12
	// se bi bo qua vi khong khop lenh dang cho).
	goi9 := clA.doiGoiLenh(t, -9)
	muon := []byte{
		0,          // mobId
		0, 0, 0, 1, // hpConLai = 1 (2 mau - 1 mau bi §1.3 kep, KHONG PHAI 0)
		0, 0, 0, 1, // satThuong = 1 (dame tho 2, nhung §1.3 ep xuong hp-1 = 1)
		0,   // chiMang = false
		255, // hieuUngId = -1
	}
	if !bytes.Equal(goi9.Data, muon) {
		t.Fatalf("goi -9 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			goi9.Data, len(goi9.Data), muon, len(muon))
	}
}

// TestTichHop_DanhQuaiVuotTamVanPhatGoiChiMangFalse la Finding 2(b): dac ta
// §1.9 "vuot tam -> sat thuong 0 NHUNG VAN PHAT GOI" — truoc do khong test
// nao dung mot muc tieu THAT SU ngoai tam (ca ba test truoc deu dung quai
// mac dinh, luon trong tam 100 tinh tu vi tri spawn). Kich ban: A di chuyen
// ra xa (qua goi THAT -7) roi danh quai Id=0 tu khoang cach > 100 (tam can
// chien) -> PHAI van nhan `-9` (KHONG im lang), voi satThuong=0 va
// chiMang=false — khang dinh dong "EP TUONG MINH" moi them trong case
// suKienDanhQuai (khuvuc.go, Finding 2a: `if ketQua.Truot { chiMang = false
// }`), dam bao nghia vu bang giao Task 3 (chotkiem.go:75-81) khong con dung
// tinh co (truoc do chiMang luon gan cung false o MOI nhanh nen "dung" ma
// khong co gi thuc su ep rieng cho truong hop Truot).
func TestTichHop_DanhQuaiVuotTamVanPhatGoiChiMangFalse(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	soKhuVuc := MoSoKhuVuc()
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTK := "farA" + hauToNganNgau(t)
	taiKhoanA := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "FarA" + hauToNganNgau(t)

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
	// Task 9: khongConGoiGiTruQuaiDanh — quai mac dinh trong tam tang 2 AI co
	// the tu danh A ngay, xem chu thich ham.
	if !clA.khongConGoiGiTruQuaiDanh(t) {
		t.Fatal("A vao khu vuc mot minh -> khong duoc nhan goi nao NGOAI quai AI (chua co ai khac)")
	}

	clA.chonChieu(1) // can chien, tam 100
	if !clA.khongConGoiGiTruQuaiDanh(t) {
		t.Fatal("chon chieu (34) khong duoc phat goi tin nao NGOAI quai AI")
	}

	// A di chuyen tu (100,384) ra o di duoc XA NHAT ve ben trai, roi danh quai
	// Id=0 o (160,384) tu ngoai tam can chien 100.
	//
	// PHAI la 48, KHONG duoc la 0: tu khi map 1 Lang La co TUONG VO HINH o hai
	// cot bien (de nguoi choi khong roi xuong ho o mep map), cac o x=0..36 deu
	// bi chan. Gui -7 toi x=0 bi actor tu choi, A DUNG YEN o x=100, khoang cach
	// toi quai chi con |100-160| = 60 < 100 -> quai TRUNG don 2 sat thuong va
	// test do voi "goi -9 byte sai". Doi bo cuc map ma quen dong nay thi trieu
	// chung se lai la "sai byte", khong he giong nguyen nhan that.
	//
	// |48-160| = 112 > 100 -> vuot tam, dung y do cua test.
	toXXa := int16(48)
	toYXa := int16(nvA.ViTriY)
	clA.guiDiChuyen(toXXa, toYXa)
	// -7 chi phat cho NGUOI KHAC (client tu ve chuyen dong cua chinh minh) —
	// A mot minh trong khu vuc nen KHONG nhan goi nao o buoc nay (ngoai co
	// the co quai AI con sot lai TU TRUOC luc di chuyen).
	if !clA.khongConGoiGiTruQuaiDanh(t) {
		t.Fatal("A tu di chuyen -> khong duoc nhan lai goi -7 cua chinh minh (ngoai quai AI)")
	}

	clA.guiDanhQuai(0)

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

	// VAN phai nhan -9 (khong im lang) — dac ta §1.9. satThuong=0, chiMang=
	// false, hpConLai KHONG doi (100, vi khong co sat thuong that su cham
	// muc tieu).
	goi9 := clA.doiGoiLenh(t, -9)
	muon := []byte{
		0,            // mobId
		0, 0, 0, 100, // hpConLai = 100 (khong doi, vuot tam nen khong gay sat thuong)
		0, 0, 0, 0, // satThuong = 0
		0,   // chiMang = false — khang dinh dong "EP TUONG MINH" (Finding 2a)
		255, // hieuUngId = -1
	}
	if !bytes.Equal(goi9.Data, muon) {
		t.Fatalf("goi -9 (vuot tam) byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			goi9.Data, len(goi9.Data), muon, len(muon))
	}
}

// ---------------------------------------------------------------------------
// Ke hoach 6 Task 13 — chet va hoi sinh nguoi choi qua mang THAT (dac ta
// §3.4, §6.4).
// ---------------------------------------------------------------------------

// moKhuVucMotQuaiChetNguoiNgay tao MOT khu vuc (chi mapId 1 — du cho lang 0,
// lang duy nhat cac test o day dung) voi DUNG MOT con quai dung tai diem
// spawn (khoang cach 0 toi nguoi choi vua vao — chac chan trong CA HAI tam
// AI, dac ta §5.2) va SatThuongCoBan RAT LON. Ghi truc tiep k.quai TRUOC khi
// goi k.BatDau() la AN TOAN (goroutine vongLap CHUA duoc tao — dung ky thuat
// da dung o moSoKhuVucKhongQuaiMacDinh, file nay, dong tren).
//
// Ly do can ham nay thay vi MoSoKhuVuc(): quai mac dinh (taoQuaiMacDinh,
// SatThuongCoBan=5) can ~30 don (~40 giay o nhip tang 2 2000ms) moi ha het
// 150 mau mac dinh cua lang 0 — qua cham cho muc dich test nay (kiem CHUOI
// GOI khi nguoi choi chet, khong phai kiem tich luy sat thuong qua nhieu
// don). SatThuong van bi kep boi §1.2 (khong vuot HpHienTai con lai) nen
// KHONG lam mau am, chi vua du giet dung MOT don — kich hoat gan nhu ngay
// (nhip 100ms) sau khi nguoi choi vao khu vuc.
func moKhuVucMotQuaiChetNguoiNgay(t *testing.T) *SoKhuVuc {
	t.Helper()
	bd := TimBanDo(1)
	k := MoKhuVuc(bd)
	k.quai = []*Quai{{Id: 0, Hp: 100, HpToiDa: 100, X: bd.SpawnX, Y: bd.SpawnY, SatThuongCoBan: 1_000_000}}
	k.BatDau()
	return &SoKhuVuc{theoMap: map[int32]*KhuVuc{1: k}}
}

// TestTichHop_ChetVaHoiSinh la kich ban chinh cua Task 13 (dac ta §6.4): mot
// client THAT bi quai AI ha CHET (nhan -17, than byte-chinh-xac 5 byte KHONG
// co id — test bat buoc (a) cua brief) -> gui -16 xin hoi sinh -> nhan lai
// -16 (rong) roi -30/15 (day mau/chakra, vi tri hien tai) — dung DUNG chuoi
// cua dac ta §6.4. Chi mot nguoi choi (khong can nguoi thu hai): -8 (nguoi
// khac thay chet) da duoc kiem BYTE-CHINH-XAC rieng o goichiendau_test.go
// (test bat buoc (b)) va kiem SU KIEN "phat cho nguoi khac" rieng o
// khuvuc_test.go (TestNhipTimAiQuaiHaChetPhatChuoiGoiChetChoCaHaiPhia) — dung
// ca hai vi mot kich ban 2-nguoi qua mang that se can dieu khien AI nham
// DUNG nguoi bi giet (khong the dam bao tat dinh, vi ca hai nguoi choi lang0
// deu spawn CUNG mot diem), phuc tap khong can thiet cho phan con lai da
// duoc hai test kia bao phu.
func TestTichHop_ChetVaHoiSinh(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	soKhuVuc := moKhuVucMotQuaiChetNguoiNgay(t)
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTK := "chet" + hauToNganNgau(t)
	taiKhoan := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "Chet" + hauToNganNgau(t)

	cl := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer cl.dong()
	dangNhapTaoNhanVatVaoTheGioi(t, cl, tenTK, matKhau, tenNV)

	nv, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoan)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan that bai: %v", err)
	}
	if nv.MapID != 1 {
		t.Fatalf("nhan vat phai o mapId 1 (lang 0), nhan duoc %d", nv.MapID)
	}

	cl.baoVaoKhuVuc()

	// --- Chet: -17 RIENG chinh chu, than DUNG 5 BYTE, KHONG co id (dac ta
	// §7.2/§3.4 — test bat buoc (a) cua brief Task 13) ---
	goi17 := cl.doiGoiLenh(t, -17)
	if len(goi17.Data) != 5 {
		t.Fatalf("goi -17 dai %d byte, dac ta doi hoi DUNG 5 (byte cPk + short x + short y, KHONG co id)", len(goi17.Data))
	}
	r := protocol.NewReader(goi17.Data)
	cPk, errCPk := r.Byte()
	if errCPk != nil || cPk != 0 {
		t.Fatalf("cPk trong -17 = %d (loi %v), mong doi 0 (chua co he PK)", cPk, errCPk)
	}
	x, errX := r.Short()
	y, errY := r.Short()
	if errX != nil || errY != nil || x != int16(nv.ViTriX) || y != int16(nv.ViTriY) {
		t.Fatalf("vi tri trong -17 = (%d,%d) (loi %v/%v), mong doi (%d,%d) (vi tri hien tai — chet KHONG doi vi tri)",
			x, y, errX, errY, nv.ViTriX, nv.ViTriY)
	}

	// --- Hoi sinh: client gui -16 (rong) -> nhan lai -16 (rong) roi -30/15
	// (day mau/chakra) — dac ta §6.4 ---
	cl.guiHoiSinh()

	goi16 := cl.doiGoiLenh(t, lenhHoiSinh)
	if len(goi16.Data) != 0 {
		t.Fatalf("goi -16 (xac nhan hoi sinh) dai %d byte, dac ta §3.4 doi hoi RONG", len(goi16.Data))
	}

	// -30 da bi tieu thu MOT LAN cho ME_INFO (-30/0) trong
	// dangNhapTaoNhanVatVaoTheGioi -> lan toi goi -30 tren ket noi nay CHAC
	// CHAN la -30/15 (ai do hoi sinh).
	goi30 := cl.doiGoiLenh(t, -30)
	r2 := protocol.NewReader(goi30.Data)
	sub, errSub := r2.Byte()
	if errSub != nil || sub != 15 {
		t.Fatalf("byte lenh con cua -30 = %d (loi %v), mong doi 15 (ai do hoi sinh)", sub, errSub)
	}
	id, errID := r2.Int()
	if errID != nil || id != int32(nv.ID) {
		t.Fatalf("idNguoiChoi trong -30/15 = %d (loi %v), mong doi %d", id, errID, nv.ID)
	}
	hp, errHp := r2.Int()
	if errHp != nil || hp != nv.HpToiDa {
		t.Fatalf("hp trong -30/15 = %d (loi %v), mong doi %d (day mau)", hp, errHp, nv.HpToiDa)
	}
	mp, errMp := r2.Int()
	if errMp != nil || mp != nv.ChakraToiDa {
		t.Fatalf("mp trong -30/15 = %d (loi %v), mong doi %d (day chakra)", mp, errMp, nv.ChakraToiDa)
	}
}

// ---------------------------------------------------------------------------
// Ke hoach 6 Task 11 Phan A — nguoi choi danh nguoi choi (PvP) qua mang THAT,
// goi len -60 (dac ta §2.1/§3.1/§3.3), case suKienDanhNguoiChoi (khuvuc.go).
// ---------------------------------------------------------------------------

// guiDanhNguoiChoi gui goi len -60: than CHI CO int idNguoiBiDanh (dac ta §3.1).
func (cl *clientGia) guiDanhNguoiChoi(idMucTieu int32) {
	cl.t.Helper()
	w := protocol.NewWriter()
	w.Int(idMucTieu)
	if err := cl.f.WriteMessage(protocol.NewMessage(lenhNguoiDanhNguoi, w)); err != nil {
		cl.t.Fatalf("gui goi -60 (danh nguoi choi) that bai: %v", err)
	}
}

// goi60DaDoc gom lai toan bo 9 truong cua goi XUONG -60 dang "don danh" (dac
// ta §2.1/§3.3) sau khi doc, dung chung cho ca kiemChoA/kiemChoB duoi day —
// tranh doc lap lai code doc tung truong.
type goi60DaDoc struct {
	idDanh, idTrung, satThuong             int32
	skillId, soMucTieu, docTiep, loaiChieu byte
	daChet, chiMang                        byte
}

func docGoi60(data []byte) (goi60DaDoc, error) {
	r := protocol.NewReader(data)
	var g goi60DaDoc
	var err error
	if g.idDanh, err = r.Int(); err != nil {
		return g, fmt.Errorf("idNguoiDanh: %w", err)
	}
	if g.skillId, err = r.Byte(); err != nil {
		return g, fmt.Errorf("skillId: %w", err)
	}
	if g.soMucTieu, err = r.Byte(); err != nil {
		return g, fmt.Errorf("soMucTieu: %w", err)
	}
	if g.idTrung, err = r.Int(); err != nil {
		return g, fmt.Errorf("idNguoiTrung: %w", err)
	}
	if g.docTiep, err = r.Byte(); err != nil {
		return g, fmt.Errorf("docTiep: %w", err)
	}
	if g.loaiChieu, err = r.Byte(); err != nil {
		return g, fmt.Errorf("loaiChieu: %w", err)
	}
	if g.satThuong, err = r.Int(); err != nil {
		return g, fmt.Errorf("satThuong: %w", err)
	}
	if g.daChet, err = r.Byte(); err != nil {
		return g, fmt.Errorf("daChet: %w", err)
	}
	if g.chiMang, err = r.Byte(); err != nil {
		return g, fmt.Errorf("chiMang: %w", err)
	}
	return g, nil
}

// TestTichHop_DanhNguoiChoiConSong la kich ban chinh cua brief Task 11 buoc
// 1: hai client THAT trong CUNG khu vuc — A chon chieu can chien roi gui -60
// nham B. CA HAI phai nhan goi -60 (dac ta §3.3 cot "Gui cho": "Ca khu
// vuc"); A (nguoi KHAC cua B, KHONG PHAI B) nhan THEM mot goi -30/9 cap nhat
// mau cua B (dac ta §3.4 cot "Gui cho" cua -30/9: "Người khác"). Dung ban
// KHONG quai mac dinh (moSoKhuVucKhongQuaiMacDinh) — PvP khong lien quan gi
// toi quai, tranh AI quai xen vao lam phuc tap kich ban.
//
// Day CUNG la kich ban DOT BIEN CHINH cua Phan A (brief buoc 5: "xoa nhanh
// dieu phoi -60 -> test phai do"): neu xoa `case lenhNguoiDanhNguoi` khoi
// BoXuLy.XuLy (login.go), goi -60 roi vao nhanh `default: return nil` (im
// lang) -> actor khong bao gio nhan duoc suKienDanhNguoiChoi -> A/B KHONG
// nhan duoc goi -60 nao -> doiGoiLenhKhongFatal(lenhNguoiDanhNguoi) duoi day
// het han hanDoiGoi roi tra loi -> t.Fatalf. CHI test O MUC NAY (di qua dung
// XuLy that qua mang that) moi bat duoc dot bien nay — cac test goi thang
// k.xuLy(...) o khuvuc_test.go bo qua han toan lop dieu phoi cua login.go
// nen KHONG the bat duoc dot bien nay (da xac minh bang tay, xem bao cao).
func TestTichHop_DanhNguoiChoiConSong(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	// Task 11: khong dung MoSoKhuVuc() — kich ban PvP nay khong lien quan gi
	// toi quai, dung ban khong quai de tranh AI xen vao (xem chu thich
	// moSoKhuVucKhongQuaiMacDinh).
	soKhuVuc := moSoKhuVucKhongQuaiMacDinh(t)
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTKA := "pvpA" + hauToNganNgau(t)
	tenTKB := "pvpB" + hauToNganNgau(t)
	taiKhoanA := taoTaiKhoanCoMatKhau(t, tenTKA, matKhau)
	taiKhoanB := taoTaiKhoanCoMatKhau(t, tenTKB, matKhau)
	tenNVA := "PvpA" + hauToNganNgau(t)
	tenNVB := "PvpB" + hauToNganNgau(t)

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
	// (dac ta §3.2: chon chieu la trang thai noi bo).
	clA.chonChieu(1)
	if !clA.khongConGoiGi(t) {
		t.Fatal("chon chieu (34) khong duoc phat goi tin nao")
	}

	// A danh B — ca hai lang 0, cung diem spawn (khoang cach 0) -> trong tam
	// can chien 100; chua tung dung chieu -> hoi chieu tu nhien qua
	// (lanDanhCuoi zero value). satThuong tho tinh duoc = 12*20/100 = 2, dung
	// cong thuc GIONG HET cac test danh quai lang 0 (xem chu thich dau nhom
	// TestTichHop_DanhQuaiConSong).
	clA.guiDanhNguoiChoi(int32(nvB.ID))

	kiemChoA := func(cl *clientGia) error {
		goi60, err := cl.doiGoiLenhKhongFatal(lenhNguoiDanhNguoi)
		if err != nil {
			return err
		}
		g, err := docGoi60(goi60.Data)
		if err != nil {
			return fmt.Errorf("A: doc goi -60 loi: %w", err)
		}
		if g.idDanh != int32(nvA.ID) || g.skillId != 1 || g.soMucTieu != 1 ||
			g.idTrung != int32(nvB.ID) || g.docTiep != 1 || g.loaiChieu != 0 ||
			g.satThuong != 11 || g.daChet != 0 || g.chiMang != 0 {
			return fmt.Errorf("A: goi -60 noi dung sai: %+v (mong idDanh=%d idTrung=%d skillId=1 soMucTieu=1 docTiep=1 loaiChieu=0 satThuong=11 daChet=0 chiMang=0)",
				g, nvA.ID, nvB.ID)
		}

		goi309, err := cl.doiGoiLenhKhongFatal(-30)
		if err != nil {
			return err
		}
		r := protocol.NewReader(goi309.Data)
		sub, e1 := r.Byte()
		id, e2 := r.Int()
		hp, e3 := r.Int()
		hpMax, e4 := r.Int()
		if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
			return fmt.Errorf("A: doc goi -30/9 loi: sub=%v id=%v hp=%v hpMax=%v", e1, e2, e3, e4)
		}
		if sub != 9 {
			return fmt.Errorf("A: byte lenh con cua -30 = %d, mong doi 9 (cap nhat mau)", sub)
		}
		if id != int32(nvB.ID) {
			return fmt.Errorf("A: idNguoiChoi trong -30/9 = %d, mong doi %d (B)", id, nvB.ID)
		}
		if hp != nvB.HpToiDa-11 {
			return fmt.Errorf("A: hp trong -30/9 = %d, mong doi %d (200-11)", hp, nvB.HpToiDa-11)
		}
		if hpMax != nvB.HpToiDa {
			return fmt.Errorf("A: hpMax trong -30/9 = %d, mong doi %d", hpMax, nvB.HpToiDa)
		}
		return nil
	}
	kiemChoB := func(cl *clientGia) error {
		goi60, err := cl.doiGoiLenhKhongFatal(lenhNguoiDanhNguoi)
		if err != nil {
			return err
		}
		g, err := docGoi60(goi60.Data)
		if err != nil {
			return fmt.Errorf("B: doc goi -60 loi: %w", err)
		}
		if g.idDanh != int32(nvA.ID) || g.idTrung != int32(nvB.ID) || g.satThuong != 11 || g.daChet != 0 || g.chiMang != 0 {
			return fmt.Errorf("B: goi -60 noi dung sai: %+v", g)
		}
		return nil
	}

	type ketQua struct {
		ten string
		err error
	}
	ketQuaCh := make(chan ketQua, 2)
	go func() { ketQuaCh <- ketQua{ten: "A (nguoi danh, cung la 'nguoi khac' cua B)", err: kiemChoA(clA)} }()
	go func() { ketQuaCh <- ketQua{ten: "B (nguoi bi danh)", err: kiemChoB(clB)} }()
	for i := 0; i < 2; i++ {
		kq := <-ketQuaCh
		if kq.err != nil {
			t.Fatalf("%s: %v", kq.ten, kq.err)
		}
	}

	// B KHONG duoc nhan them goi nao sau -60 — DAC BIET KHONG duoc nhan
	// -30/9 cua chinh minh (dac ta §3.4 cot "Gui cho" cua -30/9: "Người
	// khác" — B khong phai "nguoi khac" cua chinh no).
	if !clB.khongConGoiGi(t) {
		t.Fatal("B khong duoc nhan them goi nao sau -60 (dac biet KHONG duoc nhan -30/9 cua chinh minh)")
	}
	// A cung khong duoc nhan them gi ngoai hai goi da kiem o tren.
	if !clA.khongConGoiGi(t) {
		t.Fatal("A khong duoc nhan them goi nao ngoai -60 va -30/9 da kiem")
	}
}

// ---------------------------------------------------------------------------
// Ke hoach 6 Task 14 — test dau-cuoi (E2E) TOAN BO chuoi chien dau qua HAI
// client THAT (mang that qua net.Pipe + mynet.PhucVu that + CSDL that
// nro_test): A danh quai con song -> A ha quai chet -> quai hoi sinh sau 5s
// -> quai danh tra A -> A danh B toi chet -> B hoi sinh. Day la cot moc
// nghiem thu tu dong cua toan bo Ke hoach 6 (brief Task 14): tung manh
// (Task 1-13) da co test rieng, nhung chua co test nao GHEP CA CHUOI lai
// thanh MOT kich ban lien tuc nhu nguoi choi that se trai qua.
// ---------------------------------------------------------------------------

// dongHoDauCuoi la DONG HO TIEM DUOC (k.dongHo) rieng cho test nay: mot moc
// thoi gian TEST tu tay "tien toi" de vuot qua cac nguong cho (hoi chieu 1s,
// hoi sinh quai 5s) MA KHONG sleep that — dung y het rang buoc cua brief:
// "Dong ho tiem duoc (k.dongHo()) de khong cho 5s that". Boc mutex vi CA test
// (ghi qua tienToi) LAN cac lan goi k.xuLy/k.nhipTim (doc qua k.dongHo(), moi
// lan chay tren MOT GOROUTINE RIENG qua ky thuat "fire-and-forget" cua
// chayVaDocDongThoi duoi day) co the cham vao gia tri nay quanh cung thoi
// diem — dung mot khoa don gian thay vi phai chung minh happens-before tinh
// vi cho tung buoc.
type dongHoDauCuoi struct {
	mu  sync.Mutex
	luc time.Time
}

func (d *dongHoDauCuoi) bayGio() time.Time {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.luc
}

func (d *dongHoDauCuoi) datLai(t time.Time) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.luc = t
}

func (d *dongHoDauCuoi) tienToi(khoang time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.luc = d.luc.Add(khoang)
}

// moKhuVucDauCuoiMotQuai xay MOT KhuVuc CHUA goi BatDau() — cung ky thuat voi
// TestTichHop_DanhQuaiChet/TestTichHop_ChetVaHoiSinh: hopThu duoc CAP PHAT TAY
// nhung KHONG co goroutine vongLap nao tieu thu no. Test tu bom su kien/goi
// k.nhipTim() TRUC TIEP de LAM CHU HOAN TOAN thu tu xay ra, thay vi phu thuoc
// vao mot ticker 100ms that chay song song voi HAI ket noi mang — can thiet o
// day vi kich ban co toi 5 buoc noi tiep brief doi hoi DUNG THU TU; mot AI
// quai "tu y" xen vao som (vd tan cong truoc khi A kip danh no) qua mot ticker
// that se pha vo dung thu tu ma khong cach nao dam bao duoc tu phia test.
//
// Quai DUY NHAT Hp=HpToiDa=3 (thay vi 100 mac dinh cua taoQuaiMacDinh): sat
// thuong lang0/cap1 chi ~2 moi don (xem chu thich dau nhom
// TestTichHop_DanhQuaiConSong), can hang chuc don that moi ha du 100 mau — qua
// cham cho muc dich Task 14 (kiem CHUOI GOI dau-cuoi, khong phai tich luy sat
// thuong qua nhieu don, cai do da co test rieng). Voi Hp=3: don 1 (2 sat
// thuong) con lai 1 mau -> con song (-9, dung brief buoc 1); don 2 (dame goc 2
// nhung quai chi con 1 mau -> kep con 1, §1.2) -> het mau (-12, dung brief
// buoc 2) — DUNG 2 don khop brief, khong can quai gia (Hp=1) nhu
// TestTichHop_DanhQuaiChet vi o day CAN mot don "con song" THAT truoc khi
// don chet, de co goi -9 cho ca hai client kiem.
func moKhuVucDauCuoiMotQuai(t *testing.T, dh *dongHoDauCuoi) (*KhuVuc, *SoKhuVuc) {
	t.Helper()
	bd := TimBanDo(1)
	k := MoKhuVuc(bd)
	k.dongHo = dh.bayGio
	k.quai = []*Quai{{
		Id: 0, MauId: 1, Cap: 1, CapTinhAnh: 1,
		Hp: 12, HpToiDa: 12, // 12 - 11 (mot don cap 1 lang 0 sau vong can bang KH8 task 13) = 1
		X: bd.SpawnX + 30, Y: bd.SpawnY,
		SatThuongCoBan: 7,
	}}
	k.hopThu = make(chan suKienKhuVuc, 16) // KHONG goi BatDau() — xem chu thich ham

	return k, &SoKhuVuc{theoMap: map[int32]*KhuVuc{1: k}}
}

// chayVaDocDongThoi chay `viec` (thuong la MOT lan k.xuLy(<-k.hopThu) hoac
// k.nhipTim()) TRONG MOT GOROUTINE RIENG — `viec` co the goi Gui() qua mang
// BEN TRONG, se BLOCK cho toi khi dau kia doc (xem canh bao bomSuKienKhongPhatGoi
// dau file nay) — CHAY DONG THOI moi ham trong `docs` (moi ham MOT goroutine
// RIENG cua no: actor co the ghi xen ke vao NHIEU ket noi cung luc, doc TUAN
// TU tung ket noi se be tac, cung ly do voi doiGoiLenhKhongFatal), ROI CHO
// `viec` HOAN TOAN tra ve (kenh xong) truoc khi ham nay tra ve. Diem dong bo
// cuoi cung nay la BAT BUOC de test duoc phep CHAM TRUC TIEP trang thai actor
// (vd k.nguoiTrong[...].nv.ViTriX/HpHienTai) NGAY SAU khi goi ham nay — dung
// tinh than voi bomSuKienKhongPhatGoi, mo rong cho truong hop CO Gui() ben
// trong (khac han bomSuKienKhongPhatGoi chi an toan khi KHONG co Gui() nao).
func chayVaDocDongThoi(t *testing.T, viec func(), docs ...func() error) {
	t.Helper()
	xong := make(chan struct{})
	go func() {
		viec()
		close(xong)
	}()

	ketQuaCh := make(chan error, len(docs))
	for _, doc := range docs {
		doc := doc
		go func() { ketQuaCh <- doc() }()
	}
	for i := 0; i < len(docs); i++ {
		if err := <-ketQuaCh; err != nil {
			t.Fatal(err)
		}
	}
	<-xong
}

// kiemGoi3015 kiem than goi -30/15 (goiAiDoHoiSinh, dac ta §3.4): sub=15,
// id=nvNguoiHoiSinh.ID, hp/mp = DAY DU (HpToiDa/ChakraToiDa — hoi sinh luon
// day mau/chakra, dac ta §6.4), x/y = vi tri HIEN TAI cua nguoi hoi sinh
// (KHONG doi luc chet/hoi sinh — xuLyNguoiChoiChet/case suKienHoiSinh deu
// chi doc thang nv.ViTriX/Y, khong tinh toan lai, xem khuvuc.go).
//
// xMongDoi/yMongDoi GHIM CUNG (Vong sua 1, theo yeu cau nguoi soi): truyen tu
// noi goi thay vi doc lai nv.ViTriX/Y cua actor — o day cu the la vi tri B
// ĐA BI TEST DOI (k.nguoiTrong[nvB.ID].nv.ViTriX = bd.SpawnX-90, xem chu
// thich buoc "dua B ra xa quai" trong TestTichHop_DauCuoiHaiNguoiChoiChienDauVaQuai),
// KHONG con la vi tri spawn goc trong nvNguoiHoiSinh (snapshot CSDL lay TRUOC
// khi doi vi tri) — vi vay ham nay KHONG the tu suy ra x/y mong doi tu chinh
// nvNguoiHoiSinh, phai nhan tu tham so. Day chinh la khang dinh dong vong: neu
// goiAiDoHoiSinh vo tinh doc sai truong (vd tra ve vi tri spawn goc thay vi vi
// tri hien tai) test se bat duoc, khac voi truoc Vong sua 1 (chi doc de tieu
// thu du than goi, khong kiem gia tri).
func kiemGoi3015(m *protocol.Message, nvNguoiHoiSinh *store.NhanVat, xMongDoi, yMongDoi int16) error {
	r := protocol.NewReader(m.Data)
	sub, e1 := r.Byte()
	id, e2 := r.Int()
	hp, e3 := r.Int()
	mp, e4 := r.Int()
	x, e5 := r.Short()
	y, e6 := r.Short()
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil || e5 != nil || e6 != nil {
		return fmt.Errorf("doc -30/15 loi: sub=%v id=%v hp=%v mp=%v x=%v y=%v", e1, e2, e3, e4, e5, e6)
	}
	if sub != subAiDoHoiSinh || id != int32(nvNguoiHoiSinh.ID) || hp != nvNguoiHoiSinh.HpToiDa || mp != nvNguoiHoiSinh.ChakraToiDa || x != xMongDoi || y != yMongDoi {
		return fmt.Errorf("-30/15 sai: sub=%d id=%d hp=%d mp=%d x=%d y=%d, mong sub=%d id=%d hp=%d mp=%d x=%d y=%d",
			sub, id, hp, mp, x, y, subAiDoHoiSinh, nvNguoiHoiSinh.ID, nvNguoiHoiSinh.HpToiDa, nvNguoiHoiSinh.ChakraToiDa, xMongDoi, yMongDoi)
	}
	return nil
}

// TestTichHop_DauCuoiHaiNguoiChoiChienDauVaQuai la kich ban chinh cua Task 14
// (brief Step 1): hai client THAT dang nhap that, tao nhan vat that, vao CUNG
// khu vuc, roi dien TRON chuoi chien dau brief doi hoi:
//
//  1. A chon chieu, danh quai -> B thay 54 va -9.
//  2. A danh quai toi chet -> B thay -12; sau 5s (dong ho tiem duoc, KHONG
//     sleep that) ca hai thay -13.
//  3. Quai danh A -> A nhan -11, B nhan -10.
//  4. A danh B toi chet -> ca hai thay -60, A thay -8 cua B.
//  5. B hoi sinh -> ca hai thay -30/15.
func TestTichHop_DauCuoiHaiNguoiChoiChienDauVaQuai(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	dh := &dongHoDauCuoi{}
	dh.datLai(time.Now())
	k, soKhuVuc := moKhuVucDauCuoiMotQuai(t, dh)
	bd := k.BanDo()
	// KHONG defer soKhuVuc.DungTatCa(): khu vuc nay KHONG goi BatDau() (xem chu
	// thich moKhuVucDauCuoiMotQuai) nen k.dung/k.xong con nil — Dung() se PANIC
	// (close(nil)) roi treo vinh vien (<-nil), dung y het ly do
	// TestTichHop_DanhQuaiChet/TestTichHop_ChetVaHoiSinh o tren KHONG defer no.

	matKhau := "matkhau123"
	tenTKA := "e2eA" + hauToNganNgau(t)
	tenTKB := "e2eB" + hauToNganNgau(t)
	taiKhoanA := taoTaiKhoanCoMatKhau(t, tenTKA, matKhau)
	taiKhoanB := taoTaiKhoanCoMatKhau(t, tenTKB, matKhau)
	tenNVA := "E2eA" + hauToNganNgau(t)
	tenNVB := "E2eB" + hauToNganNgau(t)

	clA := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer clA.dong()
	clB := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer clB.dong()

	dangNhapTaoNhanVatVaoTheGioi(t, clA, tenTKA, matKhau, tenNVA)
	<-k.hopThu // xa suKienXinDanhSachQuai cua A, xem chu thich TestTichHop_DanhQuaiChet
	dangNhapTaoNhanVatVaoTheGioi(t, clB, tenTKB, matKhau, tenNVB)
	<-k.hopThu // xa suKienXinDanhSachQuai cua B

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

	// A vao khu vuc mot minh -> suKienVao KHONG Gui() gi ca (map con rong
	// truoc do) -> an toan bomSuKienKhongPhatGoi, dung ky thuat
	// TestTichHop_DanhQuaiChet.
	clA.baoVaoKhuVuc()
	bomSuKienKhongPhatGoi(k, <-k.hopThu)
	if _, co := k.nguoiTrong[nvA.ID]; !co {
		t.Fatal("A phai co trong k.nguoiTrong sau khi xu ly suKienVao")
	}

	// B vao sau -> suKienVao PHAT -5 cho CA HAI chieu (B thay A, A thay B) ->
	// can doc dong thoi tren ca hai ket noi (xem chu thich chayVaDocDongThoi).
	clB.baoVaoKhuVuc()
	chayVaDocDongThoi(t, func() { k.xuLy(<-k.hopThu) },
		func() error {
			m, err := clB.doiGoiLenhKhongFatal(-5)
			if err != nil {
				return fmt.Errorf("B doc -5 (thay A) loi: %w", err)
			}
			r := protocol.NewReader(m.Data)
			id, errID := r.Int()
			if errID != nil || id != int32(nvA.ID) {
				return fmt.Errorf("B: id trong -5 = %d (loi %v), mong doi %d (A)", id, errID, nvA.ID)
			}
			return nil
		},
		func() error {
			m, err := clA.doiGoiLenhKhongFatal(-5)
			if err != nil {
				return fmt.Errorf("A doc -5 (thay B) loi: %w", err)
			}
			r := protocol.NewReader(m.Data)
			id, errID := r.Int()
			if errID != nil || id != int32(nvB.ID) {
				return fmt.Errorf("A: id trong -5 = %d (loi %v), mong doi %d (B)", id, errID, nvB.ID)
			}
			return nil
		},
	)

	// Dua B ra XA quai (nhung VAN trong tam can chien 100 cua A) NGAY SAU khi
	// B vao khu vuc — CHI MOT LAN, dung suot phan con lai cua kich ban. Vi sao
	// can: brief buoc (3) "quai danh A" dua vao AI tang 2 (chonMucTieu,
	// quai.go) — "bat ky ai TRONG TAM 100, GAN NHAT" — quai vua hoi sinh o
	// buoc (2) MAT SACH danh sach "da tung danh" (tang 1), nen KHONG con ly do
	// nao de AI uu tien A hon B. Neu A/B dung CUNG mot diem (ca hai spawn tai
	// (100,384) cua lang 0), khoang cach toi quai BANG NHAU -> AI chon ai la
	// NGAU NHIEN (thu tu duyet map[int64]*nguoiTrongKhuVuc cua Go KHONG dinh)
	// -> vi pham yeu cau "khong chap chon" cua task. Chon do lech de CA HAI
	// dieu kien cung thoa MOT LUC, deu co bien an toan, khong ke sat bien:
	//   - quai o (SpawnX+30, SpawnY); A dung nguyen (SpawnX, SpawnY) -> cach
	//     quai 30 (trong tam AI 100, con nhieu bien).
	//   - B doi sang (SpawnX-90, SpawnY) -> cach A DUNG 90 (<=100, VAN trong
	//     tam can chien cua A cho brief buoc (4) sau nay, bien 10) NHUNG cach
	//     quai 30+90=120 (>100, NGOAI tam AI, bien 20) -> AI tang 2 CHI CON
	//     MOT ung vien hop le (A), khong con gi de "boc tham".
	// Ghi TRUC TIEP nv.ViTriX (KHONG qua goi -7/kiemDiChuyenHopLe): AN TOAN vi
	// dung SAU khi chayVaDocDongThoi() o tren da tra ve (<-xong da xay ra —
	// khong goroutine nao khac dang cham k.nguoiTrong luc nay), dung ky thuat
	// "cham truc tiep trang thai actor giua hai diem dong bo" da dung cho
	// k.quai[0].Hp trong TestTichHop_DanhQuaiChet.
	k.nguoiTrong[nvB.ID].nv.ViTriX = int32(bd.SpawnX) - 90

	// A chon chieu can chien (templateId=1, tam 100) — suKienChonChieu KHONG
	// Gui() gi (xem khuvuc.go) -> an toan bomSuKienKhongPhatGoi.
	clA.chonChieu(1)
	bomSuKienKhongPhatGoi(k, <-k.hopThu)
	if k.nguoiTrong[nvA.ID].chieuDangChon == nil {
		t.Fatal("A phai co chieuDangChon sau khi xu ly suKienChonChieu")
	}

	// docChuoi54RoiTiep doc goi 54 (kiem idDanh/skillId/mobId co dinh cua kich
	// ban nay) ROI goi `lenhSau` (vd -9 hoac -12), giao than goi do cho
	// `kiemThan` kiem tiep — dung chung cho ca hai buoc "A danh quai" (con
	// song va toi chet) de khong lap code doc 54.
	docChuoi54RoiTiep := func(cl *clientGia, lenhSau int8, kiemThan func(*protocol.Message) error, taiSao string) func() error {
		return func() error {
			goi54, err := cl.doiGoiLenhKhongFatal(lenhNguoiDanhQuai)
			if err != nil {
				return fmt.Errorf("%s: doc 54 loi: %w", taiSao, err)
			}
			r := protocol.NewReader(goi54.Data)
			idDanh, e1 := r.Int()
			skillId, e2 := r.Byte()
			mobId, e3 := r.Byte()
			if e1 != nil || e2 != nil || e3 != nil || idDanh != int32(nvA.ID) || skillId != 1 || mobId != 0 {
				return fmt.Errorf("%s: goi 54 sai: idDanh=%d skillId=%d mobId=%d (loi %v/%v/%v), mong idDanh=%d skillId=1 mobId=0",
					taiSao, idDanh, skillId, mobId, e1, e2, e3, nvA.ID)
			}
			goiSau, err := cl.doiGoiLenhKhongFatal(lenhSau)
			if err != nil {
				return fmt.Errorf("%s: doc %d loi: %w", taiSao, lenhSau, err)
			}
			if err := kiemThan(goiSau); err != nil {
				return fmt.Errorf("%s: %w", taiSao, err)
			}
			return nil
		}
	}

	// --- Brief buoc (1): A danh quai (con song) -> CA HAI thay 54 va -9. ---
	kiemGoi9 := func(m *protocol.Message) error {
		r := protocol.NewReader(m.Data)
		mobId, e1 := r.Byte()
		hpConLai, e2 := r.Int()
		satThuong, e3 := r.Int()
		chiMang, e4 := r.Byte() // bool tren day la 1 byte (0/1), KHONG co r.Bool()
		hieuUngId, e5 := r.Byte()
		if e1 != nil || e2 != nil || e3 != nil || e4 != nil || e5 != nil {
			return fmt.Errorf("doc -9 loi: %v/%v/%v/%v/%v", e1, e2, e3, e4, e5)
		}
		if mobId != 0 || hpConLai != 1 || satThuong != 11 || chiMang != 0 || hieuUngId != 255 {
			return fmt.Errorf("-9 sai: mobId=%d hpConLai=%d satThuong=%d chiMang=%d hieuUngId=%d, mong mobId=0 hpConLai=1 satThuong=11 chiMang=0 hieuUngId=255",
				mobId, hpConLai, satThuong, chiMang, hieuUngId)
		}
		return nil
	}
	clA.guiDanhQuai(0)
	chayVaDocDongThoi(t, func() { k.xuLy(<-k.hopThu) },
		docChuoi54RoiTiep(clA, lenhQuaiConSong, kiemGoi9, "A (nguoi danh, buoc 1)"),
		docChuoi54RoiTiep(clB, lenhQuaiConSong, kiemGoi9, "B (nguoi ngoai, buoc 1 — brief doi hoi B thay 54 va -9)"),
	)

	// Vuot qua hoi chieu 1000ms (cap 1 = 1000ms + noi cho phep 50ms tre mang,
	// chotkiem.go) MA KHONG sleep that — dung dong ho tiem duoc.
	dh.tienToi(1100 * time.Millisecond)

	// --- Brief buoc (2) phan 1: A danh quai TOI CHET -> ca hai thay -12. ---
	kiemGoi12 := func(m *protocol.Message) error {
		r := protocol.NewReader(m.Data)
		mobId, e1 := r.Byte()
		satThuong, e2 := r.Int()
		chiMang, e3 := r.Byte()
		soDoRoi, e4 := r.Byte()
		if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
			return fmt.Errorf("doc -12 loi: %v/%v/%v/%v", e1, e2, e3, e4)
		}
		if mobId != 0 || satThuong != 1 || chiMang != 0 || soDoRoi != 0 {
			return fmt.Errorf("-12 sai: mobId=%d satThuong=%d chiMang=%d soDoRoi=%d, mong mobId=0 satThuong=1 chiMang=0 soDoRoi=0",
				mobId, satThuong, chiMang, soDoRoi)
		}
		return nil
	}
	clA.guiDanhQuai(0)
	chayVaDocDongThoi(t, func() { k.xuLy(<-k.hopThu) },
		docChuoi54RoiTiep(clA, lenhQuaiChet, kiemGoi12, "A (nguoi danh, buoc 2 — quai chet)"),
		docChuoi54RoiTiep(clB, lenhQuaiChet, kiemGoi12, "B (nguoi ngoai, buoc 2 — brief doi hoi B thay -12)"),
	)

	// --- Brief buoc (2) phan 2: "sau 5s ca hai thay -13" — tien dong ho qua
	// nguong hoi sinh 5s MOT LAN, goi TRUC TIEP k.nhipTim() (KHONG phai su
	// kien tu hop thu — day la ham vong lap goi dinh ky, o day test goi thay
	// cho ticker that de "khong cho 5s that"). GHIM SO 5*time.Second THANG
	// (KHONG doc lai hang so thoiGianHoiSinhQuai cua khuvuc.go) — co tinh,
	// dung y het TestNhipTimQuaiHoiSinhDung5Giay (quai_test.go): neu test doc
	// lai chinh hang so san xuat de tinh khoang tien, mot dot bien sua hang so
	// do (vd 5s -> 50s) se KHONG BAO GIO bi test nay bat duoc (ca hai phia
	// "tu dong khop nhau"). Ghim boc lap thanh 100ms bien an toan cho tre xu
	// ly noi bo (khong phai tre mang — actor duoc goi TRUC TIEP, khong qua
	// ticker that). ---
	dh.tienToi(5*time.Second + 100*time.Millisecond)
	kiemGoi13 := func(m *protocol.Message) error {
		r := protocol.NewReader(m.Data)
		mobId, e1 := r.Byte()
		sys, e2 := r.Byte()
		levelBoss, e3 := r.Byte()
		hp, e4 := r.Int()
		if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
			return fmt.Errorf("doc -13 loi: %v/%v/%v/%v", e1, e2, e3, e4)
		}
		if mobId != 0 || sys != 0 || levelBoss != 0 || hp != 12 {
			return fmt.Errorf("-13 sai: mobId=%d sys=%d levelBoss=%d hp=%d, mong mobId=0 sys=0 levelBoss=0 hp=12 (HpToiDa)",
				mobId, sys, levelBoss, hp)
		}
		return nil
	}
	chayVaDocDongThoi(t, func() { k.nhipTim() },
		func() error {
			m, err := clA.doiGoiLenhKhongFatal(lenhQuaiHoiSinh)
			if err != nil {
				return fmt.Errorf("A doc -13 loi: %w", err)
			}
			return kiemGoi13(m)
		},
		func() error {
			m, err := clB.doiGoiLenhKhongFatal(lenhQuaiHoiSinh)
			if err != nil {
				return fmt.Errorf("B doc -13 loi: %w", err)
			}
			return kiemGoi13(m)
		},
	)

	// --- Brief buoc (3): quai (vua hoi sinh, con song) danh TRA A qua AI (goi
	// TRUC TIEP k.nhipTim() LAN HAI — lan dau chi lam hoi sinh roi "continue",
	// KHONG chay AI trong CUNG mot lan goi, xem chu thich nhipTim/hoiSinhQuai
	// trong khuvuc.go) -> A nhan -11 (rieng), B nhan -10 (nguoi khac). B da bi
	// day ra ngoai tam AI (xem buoc dua B o tren) nen day la ket qua TAT DINH,
	// khong phai may run. ---
	chayVaDocDongThoi(t, func() { k.nhipTim() },
		func() error {
			m, err := clA.doiGoiLenhKhongFatal(lenhQuaiDanhTrungMinh)
			if err != nil {
				return fmt.Errorf("A doc -11 loi: %w", err)
			}
			r := protocol.NewReader(m.Data)
			mobId, e1 := r.Byte()
			satThuong, e2 := r.Int()
			if e1 != nil || e2 != nil || mobId != 0 || satThuong != 7 {
				return fmt.Errorf("A: -11 sai: mobId=%d(loi %v) satThuong=%d(loi %v), mong mobId=0 satThuong=7 (SatThuongCoBan)",
					mobId, e1, satThuong, e2)
			}
			return nil
		},
		func() error {
			m, err := clB.doiGoiLenhKhongFatal(lenhQuaiDanhTrungNguoiKhac)
			if err != nil {
				return fmt.Errorf("B doc -10 loi: %w", err)
			}
			r := protocol.NewReader(m.Data)
			mobId, e1 := r.Byte()
			idNanNhan, e2 := r.Int()
			hpConLai, e3 := r.Int()
			if e1 != nil || e2 != nil || e3 != nil || mobId != 0 || idNanNhan != int32(nvA.ID) || hpConLai != nvA.HpToiDa-7 {
				return fmt.Errorf("B: -10 sai: mobId=%d idNanNhan=%d hpConLai=%d, mong mobId=0 idNanNhan=%d hpConLai=%d",
					mobId, idNanNhan, hpConLai, nvA.ID, nvA.HpToiDa-7)
			}
			return nil
		},
	)

	// Vuot qua hoi chieu cua A (lan danh gan nhat la don ha quai o buoc (2);
	// don quai danh TRA A o tren KHONG cham nguoiTrongKhuVuc.lanDanhCuoi cua
	// A — do la truong rieng cua PHIA A tan cong, khac Quai.LanDanhCuoi cua
	// phia quai tan cong) truoc khi A ra don PvP o brief buoc (4).
	dh.tienToi(1100 * time.Millisecond)

	// Ha thang HP cua B xuong du thap de MOT don PvP (~2 sat thuong, cung cong
	// thuc lang0/cap1 nhu danh quai vi B cung lang0) la CHET NGAY — tranh phai
	// grind hang chuc don that (B mac dinh 150 mau). Ghi TRUC TIEP
	// nv.HpHienTai: AN TOAN vi dang o giua hai diem dong bo (chayVaDocDongThoi
	// truoc do da tra ve), dung ky thuat da dung cho k.quai[0].Hp
	// (TestTichHop_DanhQuaiChet) va k.nguoiTrong[nvB.ID].nv.ViTriX (o tren) —
	// KHONG vi pham ky luat actor vi KHONG co goroutine nao khac dang chay tai
	// thoi diem nay.
	k.nguoiTrong[nvB.ID].nv.HpHienTai = 11 // = dung mot don cap 1 (KH8 task 13): don ket lieu bao dung 11, khong bi kep

	// --- Brief buoc (4): A danh B TOI CHET -> ca hai thay -60, A thay -8 (bao
	// "nguoi khac" B da chet) — B KHONG nhan -8 cua chinh minh, thay vao do
	// nhan -17 (rieng chinh chu, xem xuLyNguoiChoiChet). ---
	clA.guiDanhNguoiChoi(int32(nvB.ID))
	chayVaDocDongThoi(t, func() { k.xuLy(<-k.hopThu) },
		func() error {
			m, err := clA.doiGoiLenhKhongFatal(lenhNguoiDanhNguoi)
			if err != nil {
				return fmt.Errorf("A doc -60 loi: %w", err)
			}
			g, err := docGoi60(m.Data)
			if err != nil {
				return fmt.Errorf("A: doc than -60 loi: %w", err)
			}
			if g.idDanh != int32(nvA.ID) || g.idTrung != int32(nvB.ID) || g.satThuong != 11 || g.daChet != 1 {
				return fmt.Errorf("A: -60 sai: %+v, mong idDanh=%d idTrung=%d satThuong=11 daChet=1", g, nvA.ID, nvB.ID)
			}
			// A la "nguoi khac" cua B (khong phai chinh B) -> nhan THEM -8 bao
			// B da chet (dac ta §3.4, brief "A thay -8 cua B"). doiGoiLenhKhongFatal
			// bo qua duoc goi -30/9 (cap nhat mau) xen giua ma khong lam sai
			// ket qua — day chinh la ly do dung ham doc "bo qua goi khac" thay
			// vi doc dung MOT goi ke tiep.
			m8, err := clA.doiGoiLenhKhongFatal(lenhAiDoChet)
			if err != nil {
				return fmt.Errorf("A doc -8 loi: %w", err)
			}
			r := protocol.NewReader(m8.Data)
			id, e1 := r.Int()
			if e1 != nil || id != int32(nvB.ID) {
				return fmt.Errorf("A: id trong -8 = %d (loi %v), mong doi %d (B)", id, e1, nvB.ID)
			}
			return nil
		},
		func() error {
			m, err := clB.doiGoiLenhKhongFatal(lenhNguoiDanhNguoi)
			if err != nil {
				return fmt.Errorf("B doc -60 loi: %w", err)
			}
			g, err := docGoi60(m.Data)
			if err != nil {
				return fmt.Errorf("B: doc than -60 loi: %w", err)
			}
			if g.idDanh != int32(nvA.ID) || g.idTrung != int32(nvB.ID) || g.satThuong != 11 || g.daChet != 1 {
				return fmt.Errorf("B: -60 sai: %+v, mong idDanh=%d idTrung=%d satThuong=11 daChet=1", g, nvA.ID, nvB.ID)
			}
			// B la CHINH CHU vua chet -> nhan -17 (RIENG, KHONG phai -8, xem
			// dac ta §3.4/§7.2) — kiem them de chac chan hai chieu khong lan.
			m17, err := clB.doiGoiLenhKhongFatal(lenhNguoiChet)
			if err != nil {
				return fmt.Errorf("B doc -17 loi: %w", err)
			}
			if len(m17.Data) != 5 {
				return fmt.Errorf("B: -17 dai %d byte, mong 5 (dac ta §7.2)", len(m17.Data))
			}
			return nil
		},
	)

	// --- Brief buoc (5): B hoi sinh -> ca hai thay -30/15. ---
	// xB/yB mong doi = vi tri B DA BI TEST DOI toi (buoc "dua B ra xa quai" o
	// tren: k.nguoiTrong[nvB.ID].nv.ViTriX = bd.SpawnX-90) — GHIM CUNG tai
	// day (khong doc lai qua k.nguoiTrong) dung yeu cau "ghim cung" cua nguoi
	// soi, xem chu thich day du o kiemGoi3015.
	xB, yB := bd.SpawnX-90, bd.SpawnY
	clB.guiHoiSinh()
	chayVaDocDongThoi(t, func() { k.xuLy(<-k.hopThu) },
		func() error {
			m16, err := clB.doiGoiLenhKhongFatal(lenhHoiSinh)
			if err != nil {
				return fmt.Errorf("B doc -16 loi: %w", err)
			}
			if len(m16.Data) != 0 {
				return fmt.Errorf("B: -16 dai %d byte, mong RONG (dac ta §3.4)", len(m16.Data))
			}
			m30, err := clB.doiGoiLenhKhongFatal(lenhMeInfo)
			if err != nil {
				return fmt.Errorf("B doc -30/15 loi: %w", err)
			}
			return kiemGoi3015(m30, nvB, xB, yB)
		},
		func() error {
			m30, err := clA.doiGoiLenhKhongFatal(lenhMeInfo)
			if err != nil {
				return fmt.Errorf("A doc -30/15 loi: %w", err)
			}
			return kiemGoi3015(m30, nvB, xB, yB)
		},
	)

	// --- Vong sua 1 (theo yeu cau nguoi soi): chot CUOI kiem KHONG CON goi
	// thua nao chua doc tren CA HAI ket noi, sau khi da doc HET chuoi 5 buoc —
	// bat lop loi "gui thua goi cho SAI nguoi" o PHAN DUOI cua kich ban (vd A
	// vo tinh nhan them -10 cua buoc (3) danh cho B, B vo tinh nhan them
	// -11/-8 danh cho A) MA cac buoc giua khong the tu bat, vi
	// doiGoiLenhKhongFatal CHU DICH bo qua goi khong khop de tim goi dang cho
	// (vd doc -8 cua A o buoc (4) chu dong bo qua -30/9 hop le xen giua — xem
	// chu thich buoc do).
	//
	// AN TOAN dat CUOI cung (KHONG lam test chap chon): kich ban nay KHONG
	// BAO GIO goi BatDau() — khong co goroutine vongLap/ticker nao chay song
	// song nen khong co "nhieu AI" tu phat sinh ngoai y muon nhu cac test
	// dung MoSoKhuVuc() (xem khongConGoiGiTruQuaiDanh). Moi goi tren day ket
	// noi la KET QUA TAT DINH cua dung cac lan test tu tay goi k.xuLy/k.nhipTim,
	// va tat ca da duoc doc HET qua tung buoc o tren — goi con SOT lai o day
	// (neu co) chi co the la dau hieu that su cua loi "gui sai nguoi", khong
	// phai nhieu hop le.
	//
	// Han che CON LAI (chot nay KHONG khac phuc duoc): mot goi sai nguoi bi
	// "nuot" GIUA CHUNG boi mot lan doiGoiLenhKhongFatal dang bo qua de tim
	// goi khac (truong hop -8 cua A o buoc (4) neu vi du co them -10 sai gui
	// toi A xen vao) se KHONG bi bat o day, vi da bi doc mat truoc do. Lop loi
	// do (mid-stream) da duoc phu boi cac test byte-chinh-xac/kiem-tung-buoc
	// cua tung Task rieng (Task 9: TestTichHop_DanhQuaiConSong; Task 11:
	// TestTichHop_DanhNguoiChoiConSong dung khongConGoiGi ngay sau MOI buoc;
	// Task 13: cac test muc actor trong khuvuc_test.go).
	if !clA.khongConGoiGi(t) {
		t.Fatal("A khong duoc nhan them goi nao sau khi ca chuoi 5 buoc ket thuc (nghi ngo gui thua goi cho sai nguoi)")
	}
	if !clB.khongConGoiGi(t) {
		t.Fatal("B khong duoc nhan them goi nao sau khi ca chuoi 5 buoc ket thuc (nghi ngo gui thua goi cho sai nguoi)")
	}
}
