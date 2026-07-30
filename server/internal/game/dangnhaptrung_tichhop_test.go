package game

// Kiem thu tich hop dau-cuoi (Ke hoach 11 Task 7, dac ta
// 2026-07-27-luu-trang-thai-song.md §7.1/§7.2): "mot tai khoan — mot phien
// song". Dung CSDL PostgreSQL that (nro_test, qua moKhoThatChoTest —
// tichhop_test.go) + actor KhuVuc that + ThoGhi that + SoDangNhap that, di
// qua mang that (net.Pipe + mynet.PhucVu, dung y het main.go dung
// game.MoLuuTru) — day la CACH DUY NHAT kiem duoc dung dieu can kiem: KhiDong()
// chi chay qua defer cua mynet.PhucVu khi mot ket noi THAT dong, nen khong
// the mo phong bang cach goi thang mot ham noi bo.
//
// Kich ban chinh (brief Task 7): tai khoan A dang nhap, di chuyen toi mot vi
// tri GHIM CUNG, roi dang nhap LAN HAI CUNG tai khoan — phien CU phai bi DA,
// du lieu phien CU phai duoc LUU XUONG CSDL truoc khi phien MOI (va bat ky ai
// doc CSDL sau do) thay duoc du lieu MOI NHAT. Day la test chong dung LOI (1)
// cua ban Java cu (NDVSqlFetcher.login(): doc CSDL TRUOC roi moi da phien cu
// — phien moi chay tren du lieu CU HON, roi GHI DE len ban MOI cua phien cu
// khi chinh no dut ket noi sau nay -> mat do, "chưa bao giờ được vá" o ban
// Java).

import (
	"context"
	"net"
	"testing"
	"time"

	mynet "nroserver/internal/net"
	"nroserver/internal/protocol"
)

// moPhienDangNhapTrungChoTest dung mot phien XuLy that (net.Pipe +
// mynet.PhucVu that) tren *BoXuLy b DA duoc chuan bi san (DatThoGhi +
// DatSoDangNhap goi TRUOC boi nguoi goi) — TRA VE CA clientGia LAN CHINH con
// tro b, khac moPhienClientGiaKhuVuc (khuvuc_tichhop_test.go, chi goi
// NewBoXuLy(kho, soKhuVuc) TRUC TIEP ben trong, khong cho ben ngoai giu lai
// con tro). Task 7 CAN giu lai b: kich ban chinh phai doi CHINH XAC
// b.daXongDon cua phien 2 dong (xac nhan KhiDong() cua NO da chay xong hoan
// toan, ke ca buoc Xep() neu co) truoc khi goi thoGhi.DongVaChoXong — chi
// dong ket noi phia client (cl.dong()) KHONG du: buoc may chu PHAT HIEN mat
// ket noi + chay KhiDong() la BAT DONG BO tren goroutine RIENG cua PhucVu,
// khong dam bao da xong ngay khi cl.dong() tra ve.
func moPhienDangNhapTrungChoTest(t *testing.T, b *BoXuLy) *clientGia {
	t.Helper()

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

// choKhiDongXong doi (co han gio) cho toi khi b.daXongDon dong — xac nhan
// KhiDong() cua phien b da chay XONG HOAN TOAN (ke ca buoc luu, du dong bo
// hay bat dong bo qua Xep()). Dung sau khi CHU DONG dong ket noi cua mot
// phien (test tu gia lam client roi khoi) de dam bao buoc luu cua no da THAT
// SU xay ra truoc khi kiem tra ket qua — chi dong ket noi PHIA CLIENT
// (cl.dong()) khong du, vi may chu phat hien mat ket noi + chay KhiDong() la
// bat dong bo tren mot goroutine rieng (xem chu thich moPhienDangNhapTrungChoTest).
func choKhiDongXong(t *testing.T, b *BoXuLy, han time.Duration) {
	t.Helper()
	select {
	case <-b.daXongDon:
	case <-time.After(han):
		t.Fatalf("sau %s, KhiDong() cua phien van chua bao xong (daXongDon chua dong)", han)
	}
}

// dangNhapLaiVaoTheGioiVaKhuVuc dua mot client GIA (tai khoan DA CO nhan
// vat — khac dangNhapTaoNhanVatVaoTheGioi trong khuvuc_tichhop_test.go, danh
// cho tai khoan CHUA co) di tron: dang nhap -> nhan DUNG MOT goi -28/4
// (versionGame, KHONG lenh 2) -> xin du 4 catalog -> gui -28/13 -> nhan
// -82/-30/-24 -> gui -39 (vao khu vuc).
func dangNhapLaiVaoTheGioiVaKhuVuc(t *testing.T, cl *clientGia, ten, matKhau string) {
	t.Helper()
	cl.dangNhap(ten, matKhau)
	if !cl.choVaoTheGioiRoiImLang(t) {
		t.Fatal("dang nhap lai (da co nhan vat) phai nhan DUNG MOT goi -28/4 (versionGame) roi im lang (KHONG duoc lenh 2)")
	}

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

	cl.baoVaoKhuVuc()
}

// choDiChuyenXongTrenActor doi (poll co han gio) cho toi khi vi tri
// AUTHORITATIVE cua nhan vat idNhanVat trong khu vuc k khop dung (toX, toY) —
// dung AnhChupNguoiChoi (an toan goi tu ben ngoai actor, khuvuc.go) thay vi
// sleep co dinh: dam bao buoc di chuyen (-7, xu ly bat dong bo qua hop thu
// actor) da THAT SU duoc actor ap dung TRUOC khi kich ban tiep tuc sang buoc
// dang nhap lan hai — neu khong, kich ban co the (hiem nhung co the) chup anh
// phai vi tri CU do actor chua kip xu ly xong su kien di chuyen.
func choDiChuyenXongTrenActor(t *testing.T, k *KhuVuc, idNhanVat int64, toX, toY int16) {
	t.Helper()
	hetGio := time.After(2 * time.Second)
	for {
		ban := k.AnhChupNguoiChoi(idNhanVat)
		if len(ban) == 1 && int16(ban[0].ViTriX) == toX && int16(ban[0].ViTriY) == toY {
			return
		}
		select {
		case <-hetGio:
			t.Fatalf("sau 2s actor van chua cap nhat vi tri authoritative thanh (%d,%d) cho idNhanVat=%d", toX, toY, idNhanVat)
		case <-time.After(10 * time.Millisecond):
		}
	}
}

// TestDangNhapTrung_DaPhienCuChoLuuXongRoiDocViTriMoiNhat la kiem thu TDD BAT
// BUOC cua Task 7 (xem chu thich dau file).
//
// Thiet ke khang dinh CAN THAN de THAT SU bat duoc dot bien "doc truoc, da
// sau" (khong chi tinh co dung nho thu tu thoi gian): sau khi phien 2 (B, tro
// lai bang nick B nhung THAT RA CUNG tai khoan voi A — chi la lan dang nhap
// thu hai) doc xong nhan vat va VAO KHU VUC, no TU NO dut ket noi (kich hoat
// KhiDong() cua CHINH no, duong luu TU NHIEN khong lien quan gi Task 7). Neu
// code dao NGUOC thu tu (doc CSDL truoc, da phien cu sau), phien 2 se doc
// phai BAN CU (vi tri spawn ban dau) NGAY LUC dang nhap; khi no tu dut ket
// noi sau do, no se GHI DE CSDL BANG CHINH BAN CU DO — xoa mat vi tri MOI ma
// phien 1 vua luu — day chinh la "mat do" that su ma dac ta §7.1 mo ta, va la
// dieu ma CHI kiem tra "CSDL ngay sau khi phien 2 dang nhap xong" (khong lam
// buoc nay) se BO LOT (vi luc do CSDL van con dung, chi la TAM THOI dung do
// may man ve thu tu ghi — chua bi phien 2 ghi de lai).
func TestDangNhapTrung_DaPhienCuChoLuuXongRoiDocViTriMoiNhat(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	soKhuVuc := moSoKhuVucKhongQuaiMacDinh(t)
	defer soKhuVuc.DungTatCa()

	thoGhi := MoThoGhi(khoThat, 2)
	defer func() {
		// Luoi an toan don dep (an toan goi NHIEU LAN — DongVaChoXong dung
		// sync.Once noi bo): dam bao goroutine nen cua thoGhi luon duoc dung
		// han, ke ca khi test that bai som (t.Fatalf) TRUOC khi toi duoc lan
		// goi DongVaChoXong tuong minh o cuoi ham.
		ctxD, huy := context.WithTimeout(context.Background(), 5*time.Second)
		defer huy()
		_ = thoGhi.DongVaChoXong(ctxD)
	}()
	soDangNhap := MoSoDangNhap()

	taoXuLy := func() *BoXuLy {
		b := NewBoXuLy(khoThat, soKhuVuc)
		b.DatThoGhi(thoGhi)
		b.DatSoDangNhap(soDangNhap)
		return b
	}

	matKhau := "matkhau123"
	tenTK := "dupA" + hauToNganNgau(t)
	taiKhoanID := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "DupA" + hauToNganNgau(t)

	// Phien 1 (CU): dang nhap tai khoan MOI, tao nhan vat, vao khu vuc. Giu
	// lai b1 (KHONG chi taoXuLy — xem chu thich moPhienDangNhapTrungChoTest)
	// vi buoc cuoi ham can xac nhan b1.daXongDon da dong (phien 1 THAT SU bi
	// da xong), khong chi suy tu viec ket noi khong doc duoc goi nua.
	b1 := taoXuLy()
	cl1 := moPhienDangNhapTrungChoTest(t, b1)
	defer cl1.dong()
	dangNhapTaoNhanVatVaoTheGioi(t, cl1, tenTK, matKhau, tenNV)
	cl1.baoVaoKhuVuc()
	if !cl1.khongConGoiGi(t) {
		t.Fatal("A vao khu vuc mot minh -> khong duoc nhan goi nao (chua co ai khac)")
	}

	nvBanDau, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoanID)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan (truoc khi di chuyen) that bai: %v", err)
	}

	// Di chuyen toi VI TRI GHIM CUNG (trong nguong buoc 500, xem dichuyen.go)
	// — day la "du lieu MOI" ma phien cu phai luu duoc truoc khi bi da.
	viTriMoiX := int16(nvBanDau.ViTriX) - 10
	viTriMoiY := int16(nvBanDau.ViTriY) + 10
	cl1.guiDiChuyen(viTriMoiX, viTriMoiY)

	k := soKhuVuc.Tra(1)
	if k == nil {
		t.Fatal("khong tim thay khu vuc mapId=1 (lang 0)")
	}
	choDiChuyenXongTrenActor(t, k, nvBanDau.ID, viTriMoiX, viTriMoiY)

	// Phien 2 (MOI, CUNG tai khoan): dang nhap lai -> BEN TRONG xuLyDangNhap,
	// SoDangNhap.DangKy phai DA phien 1 (dong ket noi cua no) VA CHO no bao
	// da luu xong (KhiDong dong bo qua XepVaCho) TRUOC khi ham nay tra ve —
	// nen LayNhanVatTheoTaiKhoan cho phien 2 (ben trong xuLyDangNhap, sau
	// DangKy) chac chan doc duoc CSDL da co vi tri MOI.
	b2 := taoXuLy()
	cl2 := moPhienDangNhapTrungChoTest(t, b2)
	defer cl2.dong()
	dangNhapLaiVaoTheGioiVaKhuVuc(t, cl2, tenTK, matKhau)
	if !cl2.khongConGoiGi(t) {
		t.Fatal("B vao khu vuc mot minh (A da bi da va roi khoi khu vuc truoc do) -> khong duoc nhan goi nao")
	}

	// Xac nhan phien 1 (b1) THAT SU bi da: b1.daXongDon phai da dong roi (no
	// dong TRUOC khi DangKy cua phien 2 tra ve — buoc dangNhapLaiVaoTheGioiVaKhuVuc
	// o tren da hoan tat nen chac chan dung), va biDa cua no phai dong.
	select {
	case <-b1.daXongDon:
	default:
		t.Fatal("b1.daXongDon CHUA dong sau khi phien 2 dang nhap xong — phien 1 phai da bi da VA bao xong luc nay")
	}
	select {
	case <-b1.biDa:
	default:
		t.Fatal("b1.biDa CHUA dong — phien 1 phai duoc danh dau la BI DA boi phien 2")
	}

	// Phien 1 PHAI THAT SU bi da: ket noi cu phai da bi dong (doc phai loi
	// NGAY, khong phai treo cho toi han gio).
	_ = cl1.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	if _, err := cl1.f.ReadMessage(); err == nil {
		t.Fatal("phien CU van doc duoc goi tin sau khi dang nhap lan hai — ket noi cu PHAI da bi DONG (bi da)")
	}

	// Phien 2 TU NO dut ket noi (duong luu TU NHIEN, khong bi da) — day la
	// buoc BAT BUOC de bay dot bien "doc truoc, da sau" (xem chu thich dau
	// ham): neu phien 2 da doc phai du lieu CU (spawn ban dau) luc dang nhap,
	// buoc nay se GHI DE CSDL VE LAI vi tri CU, phoi bay dung dieu test can
	// bat.
	//
	// QUAN TRONG: chi dong ket noi PHIA CLIENT (cl2.dong()) la CHUA DU de
	// khang dinh buoc luu cua phien 2 da xay ra — may chu PHAT HIEN mat ket
	// noi (ReadMessage tra loi) roi chay KhiDong() la BAT DONG BO tren
	// goroutine RIENG cua PhucVu, khong dam bao xong ngay khi cl2.dong() tra
	// ve. PHAI cho CHINH XAC b2.daXongDon (choKhiDongXong) TRUOC KHI goi
	// thoGhi.DongVaChoXong — thieu buoc cho nay se lam ca kiem thu FLAKY: neu
	// DongVaChoXong tinh co chay TRUOC khi b2.KhiDong() kip goi Xep(), no se
	// thay hang doi RONG va tra ve NGAY (khong co gi sai — nhung cung khong
	// cho duoc dieu can cho), roi Xep() (voi du lieu CU neu co dot bien) moi
	// chay SAU do — qua muon de kiem thu nay bat duoc.
	cl2.dong()
	choKhiDongXong(t, b2, 5*time.Second)

	ctxCho, huy := context.WithTimeout(context.Background(), 5*time.Second)
	defer huy()
	if err := thoGhi.DongVaChoXong(ctxCho); err != nil {
		t.Fatalf("cho tho ghi xa het hang doi (sau khi phien 2 dut ket noi) loi: %v", err)
	}

	nvSauCung, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoanID)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan (sau cung) that bai: %v", err)
	}
	if nvSauCung.ViTriX != int32(viTriMoiX) || nvSauCung.ViTriY != int32(viTriMoiY) {
		t.Fatalf("CSDL sau cung: vi tri = (%d,%d), mong doi VI TRI MOI NHAT da phien 1 luu = (%d,%d).\n"+
			"Neu day la vi tri SPAWN BAN DAU (%d,%d): phien 2 da DOC PHAI DU LIEU CU luc dang nhap "+
			"(chua kip cho phien 1 luu xong), roi TU GHI DE lai CSDL bang ban cu do khi no dut ket noi — "+
			"dung dung bay 'doc truoc, da sau' cua ban Java cu (dac ta §7.1) ma Task 7 phai chan.",
			nvSauCung.ViTriX, nvSauCung.ViTriY, viTriMoiX, viTriMoiY, nvBanDau.ViTriX, nvBanDau.ViTriY)
	}
}
