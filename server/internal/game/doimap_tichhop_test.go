package game

import (
	"context"
	"testing"
	"time"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// guiDoiMap gui goi -23 (than RONG) — dung nhu Service.requestChangeMap ben
// client (Service.cs:801).
func (cl *clientGia) guiDoiMap() {
	cl.t.Helper()
	if err := cl.f.WriteMessage(&protocol.Message{Command: lenhDoiMap}); err != nil {
		cl.t.Fatalf("gui goi -23 (doi map) that bai: %v", err)
	}
}

// choNguoiVaoKhuVuc doi cho toi khi actor cua khu vuc THAT SU da nhan nguoi
// choi vao.
//
// Vi sao phai cho: goi -39 di qua net.Pipe (dong bo, ghi xong la may chu da
// DOC) nhung viec dua nguoi vao khu vuc di tiep mot chang nua qua hop thu cua
// actor — bat dong bo. Doc ngay sau khi ghi la doc truoc khi actor kip xu ly.
// Doi bang cach hoi lai chinh actor (khong ngu mot khoang doan mo) nen test
// khong nhanh cham theo may.
func choNguoiVaoKhuVuc(t *testing.T, k *KhuVuc, idNhanVat int64) []store.NhanVat {
	t.Helper()
	hetHan := time.Now().Add(3 * time.Second)
	for {
		ban, err := k.AnhChupNguoiChoiCoLoi(idNhanVat)
		if err == nil && len(ban) > 0 {
			return ban
		}
		if time.Now().After(hetHan) {
			t.Fatalf("sau 3s actor van chua nhan nguoi choi %d vao khu vuc (loi %v)", idNhanVat, err)
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// TestTichHop_DiQuaCongRoiLangLa la kich ban chinh cua duong doi map:
// nguoi choi di toi cong torii roi bam vao no, phai sang duoc Lang Cat.
//
// Truoc ban nay khong co duong nao roi map 1 — day la test chan viec quay
// lai tinh trang do.
func TestTichHop_DiQuaCongRoiLangLa(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	// Khong dung quai mac dinh: test nay khong lien quan gi toi quai, ma quai
	// o gan diem xuat phat co the tu danh nhan vat giua chung (xem chu thich
	// moSoKhuVucKhongQuaiMacDinh).
	soKhuVuc := moSoKhuVucKhongQuaiMacDinh(t)
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTK := "cong" + hauToNganNgau(t)
	taiKhoan := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "Cong" + hauToNganNgau(t)

	cl := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer cl.dong()
	dangNhapTaoNhanVatVaoTheGioi(t, cl, tenTK, matKhau, tenNV)
	cl.baoVaoKhuVuc()

	nv, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoan)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan that bai: %v", err)
	}
	if nv.MapID != 1 {
		t.Fatalf("nhan vat phai xuat phat o map 1, nhan duoc %d", nv.MapID)
	}

	// Di tu diem xuat phat (100,384) toi giua cong torii (441,408). Mot buoc
	// duy nhat: khoang cach ~342 < NguongBuoc 500, va ca doan duong nam tren
	// hang 16-17 nen di duoc.
	cl.guiDiChuyen(441, 408)
	cl.guiDoiMap()

	// -22 (xoa map cu) PHAI toi truoc -24: client dung no de bo muc tieu dang
	// ngam va khoa phim trong luc tai map moi.
	if m := cl.doiGoiLenh(t, lenhXoaMap); m.Command != -22 {
		t.Fatalf("mong doi goi -22 truoc khi doi map, nhan %d", m.Command)
	}
	m24 := cl.doiGoiLenh(t, lenhMapInfo)

	d := moDocGoi(t, m24)
	d.b("mapId", 2) // Lang Cat
	d.b("planetId", 0)
	d.b("tileId", 1)
	d.b("bgId", 0)
	d.b("type", 0)
	d.utf("mapName", "Làng Cát")
	d.b("zoneId", 0)
	d.sh("x", 100) // diem ha canh cua cong, ghim tay
	d.sh("y", 384)
	// Lang Cat phai co cong VE — khong thi vua sua xong mot nha tu lai xay
	// nha tu khac.
	d.b("waypointCount", 1)
	d.sh("wp.minX", 48)
	d.sh("wp.minY", 336)
	d.sh("wp.maxX", 168)
	d.sh("wp.maxY", 456)
	d.b("wp.laCong", 1)
	d.b("wp.ngoaiTuyen", 0)
	d.utf("wp.ten", "Về Làng Lá")
	d.b("mobCount", 0)
	d.b("byte du phong", 0)
	d.b("npcCount", 0)  // Iruka chi dung o map 1
	d.b("itemCount", 0) //
	d.sh("soBong", 0)   // Lang Cat chua co my thuat rieng -> chua co bong do
	d.sh("effMap", 0)
	d.b("bgType", 0)
	d.b("idSpaceShip", 0)
	d.b("isMapDouble", 0)
	d.hetGoi("MAP INFO map 2")

	// Client tai xong map moi thi gui -39; nhan vat phai duoc dua vao khu vuc
	// MOI (khong phai khu vuc cu). Mot minh trong do nen khong co goi nao them.
	cl.baoVaoKhuVuc()
	if !cl.khongConGoiGi(t) {
		t.Fatal("vao khu vuc map 2 mot minh ma van nhan goi tin — sai khu vuc?")
	}

	// MapID cua chinh ban ghi nhan vat cung phai doi. Goi -24 o tren lay mapId
	// tu BAN DO dich nen no dung ke ca khi ban ghi quen doi — chi cho nay moi
	// bat duoc. Quen doi thi lan dang nhap sau nguoi choi hien lai o Lang La.
	ban := choNguoiVaoKhuVuc(t, soKhuVuc.Tra(2), nv.ID)
	if ban[0].MapID != 2 {
		t.Errorf("MapID trong ban ghi nhan vat = %d sau khi di qua cong, mong doi 2", ban[0].MapID)
	}
	if ban[0].ViTriX != 100 || ban[0].ViTriY != 384 {
		t.Errorf("vi tri sau khi di qua cong = (%d,%d), mong doi (100,384)", ban[0].ViTriX, ban[0].ViTriY)
	}
}

// TestTichHop_DoiMapKhiKhongDungTrongCongVanMoKhoaClient la test AN TOAN
// quan trong nhat cua duong nay.
//
// Client tu khoa minh lai NGAY khi gui -23 (Char.ischangingMap = true,
// isLockKey = true, hien hop thoai "Xin cho" — Waypoint.cs / Char.cs:2322).
// Neu may chu im lang thi nhan vat DUNG HINH VINH VIEN. Vay ngay ca khi
// nguoi choi khong dung trong cong nao, may chu VAN phai tra loi — bang goi
// 46 (dat lai vi tri) de mo khoa tai cho.
func TestTichHop_DoiMapKhiKhongDungTrongCongVanMoKhoaClient(t *testing.T) {
	khoThat := moKhoThatChoTest(t)

	soKhuVuc := moSoKhuVucKhongQuaiMacDinh(t)
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTK := "ngoai" + hauToNganNgau(t)
	taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "Ngoai" + hauToNganNgau(t)

	cl := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer cl.dong()
	dangNhapTaoNhanVatVaoTheGioi(t, cl, tenTK, matKhau, tenNV)
	cl.baoVaoKhuVuc()

	// Dung yen o diem xuat phat (100,384) — cach cong torii rat xa — roi gui
	// -23 (mot client bi sua doi, hoac bam dung luc goi di chuyen chua toi).
	cl.guiDoiMap()

	m46 := cl.doiGoiLenh(t, lenhDatLaiViTri)
	d := moDocGoi(t, m46)
	d.sh("x", 100) // dung vi tri hien tai -> client an hop thoai cho, khong truot di dau
	d.sh("y", 384)
	d.hetGoi("goi 46")

	// Va TUYET DOI khong duoc doi map: khong co -24 nao theo sau.
	if !cl.khongConGoiGi(t) {
		t.Fatal("gui -23 khi khong dung trong cong nao van lam may chu gui them goi — co the da doi map nham")
	}
}

// TestTichHop_DoiMapKhiDangChetBiTuChoi: nguoi choi dang nam cho hoi sinh
// khong duoc di qua cong. Client luc do dang o luong "chet" rieng (man hoi
// sinh); day ho sang map khac giua chung la hai ben lech trang thai. Van
// PHAI tra goi 46 de client khong ket o hop thoai "Xin cho".
func TestTichHop_DoiMapKhiDangChetBiTuChoi(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	soKhuVuc := moSoKhuVucKhongQuaiMacDinh(t)
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTK := "chet" + hauToNganNgau(t)
	taiKhoan := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "Chet" + hauToNganNgau(t)

	cl := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer cl.dong()
	dangNhapTaoNhanVatVaoTheGioi(t, cl, tenTK, matKhau, tenNV)
	cl.baoVaoKhuVuc()

	nv, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoan)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan that bai: %v", err)
	}
	choNguoiVaoKhuVuc(t, soKhuVuc.Tra(1), nv.ID)

	// Di toi dung giua cong, roi "chet" ngay tai do.
	cl.guiDiChuyen(441, 408)
	k := soKhuVuc.Tra(1)
	xong := make(chan struct{})
	k.Gui(suKienChay{f: func() {
		if nguoi, co := k.nguoiTrong[nv.ID]; co {
			nguoi.nv.HpHienTai = 0
		}
		close(xong)
	}})
	<-xong

	cl.guiDoiMap()

	m46 := cl.doiGoiLenh(t, lenhDatLaiViTri)
	d := moDocGoi(t, m46)
	d.sh("x", 441)
	d.sh("y", 408)
	d.hetGoi("goi 46")
	if !cl.khongConGoiGi(t) {
		t.Fatal("nguoi choi dang chet van doi duoc map (nhan them goi sau 46)")
	}
}

// TestTichHop_DoiMapMangTheoMauVaChakra: doi map KHONG duoc lam moi nhan vat.
// Anh chup lay tu actor truoc khi roi khu vuc phai duoc chep nguyen sang map
// moi — neu ai do sau nay sua xuLyDoiMap thanh "dung b.nhanVat co san" thi
// mau/chakra vua bi mat trong tran danh o map cu se hoi day nhu chua co gi.
func TestTichHop_DoiMapMangTheoMauVaChakra(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	soKhuVuc := moSoKhuVucKhongQuaiMacDinh(t)
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTK := "mang" + hauToNganNgau(t)
	taiKhoan := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "Mang" + hauToNganNgau(t)

	cl := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer cl.dong()
	dangNhapTaoNhanVatVaoTheGioi(t, cl, tenTK, matKhau, tenNV)
	cl.baoVaoKhuVuc()

	nv, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoan)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan that bai: %v", err)
	}
	// Ha mau nhan vat NGAY TRONG actor (dung duong ma actor dung), roi doi map.
	k := soKhuVuc.Tra(1)
	xong := make(chan struct{})
	k.Gui(suKienChay{f: func() {
		if nguoi, co := k.nguoiTrong[nv.ID]; co {
			nguoi.nv.HpHienTai = 37
		}
		close(xong)
	}})
	<-xong

	cl.guiDiChuyen(441, 408)
	cl.guiDoiMap()
	cl.doiGoiLenh(t, lenhXoaMap)
	cl.doiGoiLenh(t, lenhMapInfo)
	cl.baoVaoKhuVuc()

	// Doc lai mau tu actor cua khu vuc MOI.
	ban := choNguoiVaoKhuVuc(t, soKhuVuc.Tra(2), nv.ID)
	if ban[0].HpHienTai != 37 {
		t.Errorf("mau sau khi doi map = %d, mong doi 37 — doi map da lam moi nhan vat", ban[0].HpHienTai)
	}
	if ban[0].MapID != 2 {
		t.Errorf("MapID sau khi doi map = %d, mong doi 2", ban[0].MapID)
	}
}

// TestTichHop_DoiMapGoNguoiChoiKhoiKhuVucCu: nguoi choi KHONG duoc con nam o
// khu vuc cu sau khi doi map. Con lai la hai actor cung ghi mot con tro nhan
// vat — dung loai tranh chap ma luat "chi vongLap duoc cham trang thai khu
// vuc" (§6.9) sinh ra de chan.
func TestTichHop_DoiMapGoNguoiChoiKhoiKhuVucCu(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	soKhuVuc := moSoKhuVucKhongQuaiMacDinh(t)
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTK := "go" + hauToNganNgau(t)
	taiKhoan := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "Go" + hauToNganNgau(t)

	cl := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer cl.dong()
	dangNhapTaoNhanVatVaoTheGioi(t, cl, tenTK, matKhau, tenNV)
	cl.baoVaoKhuVuc()

	nv, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoan)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan that bai: %v", err)
	}

	cl.guiDiChuyen(441, 408)
	cl.guiDoiMap()
	cl.doiGoiLenh(t, lenhXoaMap)
	cl.doiGoiLenh(t, lenhMapInfo)

	k1 := soKhuVuc.Tra(1)
	conLai, err := k1.AnhChupNguoiChoiCoLoi(nv.ID)
	if err != nil {
		t.Fatalf("khong hoi duoc khu vuc map 1: %v", err)
	}
	if len(conLai) != 0 {
		t.Errorf("khu vuc map 1 van con nguoi choi %d sau khi ho da di qua cong", nv.ID)
	}
}
