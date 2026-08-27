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

// ---------------------------------------------------------------------------
// BON DUONG RA CUA xuLyDoiMap — hai duong duoi day truoc ban nay IM LANG,
// tuc la nguoi choi bam cong xong DUNG HINH VINH VIEN (client tu khoa minh
// roi doi may chu, xem chu thich dau doimap.go). Hai test duoi khoa lai
// chung. Chung KHONG phai gia dinh: han cho control chi 100ms
// (thoiGianChoControl, khuvuc.go) trong khi actor phai rut can hang gameplay
// truoc khi tra loi, va van an toan con DUNG HAN actor sau 100 nhip panic
// lien tiep — tu luc do MOI nguoi bam cong trong khu vuc do deu roi vao
// duong ra thu hai.
// ---------------------------------------------------------------------------

// TestTichHop_DoiMapKhiMapHienTaiKhongCoDanhMucVanMoKhoaClient khoa duong ra
// "TimBanDo(mapId hien tai) == nil".
//
// Xay ra khi ban ghi nhan vat mang mot mapId khong con trong danh muc (du
// lieu cu sau khi doi bang ban do, hoac mot ban va lam hep danh muc lai).
// Truoc ban nay may chu tra ve nil — khong goi nao — va nguoi choi ket o hop
// thoai "Xin cho".
func TestTichHop_DoiMapKhiMapHienTaiKhongCoDanhMucVanMoKhoaClient(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	soKhuVuc := moSoKhuVucKhongQuaiMacDinh(t)
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTK := "nomap" + hauToNganNgau(t)
	taiKhoan := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "Nomap" + hauToNganNgau(t)

	cl := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer cl.dong()
	dangNhapTaoNhanVatVaoTheGioi(t, cl, tenTK, matKhau, tenNV)
	cl.baoVaoKhuVuc()

	nv, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoan)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan that bai: %v", err)
	}
	k := soKhuVuc.Tra(1)
	choNguoiVaoKhuVuc(t, k, nv.ID)

	// Doi mapId cua ban ghi nhan vat sang mot so KHONG co trong danh muc.
	// Lam TU TRONG actor (suKienChay) — actor la chu so huu con tro nhan vat,
	// dung het cach cac test khac trong file nay ha mau nhan vat.
	xong := make(chan struct{})
	k.Gui(suKienChay{f: func() {
		if nguoi, co := k.nguoiTrong[nv.ID]; co {
			nguoi.nv.MapID = 99 // 99 khong nam trong 5 map cua danh muc
		}
		close(xong)
	}})
	<-xong

	cl.guiDoiMap()

	m46 := cl.doiGoiLenh(t, lenhDatLaiViTri)
	d := moDocGoi(t, m46)
	d.sh("x", 100)
	d.sh("y", 384)
	d.hetGoi("goi 46")
	if !cl.khongConGoiGi(t) {
		t.Fatal("sau goi 46 con goi nua — nhan vat o mapId la khong duoc doi map di dau ca")
	}
}

// TestTichHop_DoiMapKhiActorDaDungVanMoKhoaClient khoa duong ra "khong chup
// duoc trang thai".
//
// Kich ban THAT lam no xay ra: van an toan cua actor (nguongPanicDungKhuVuc =
// 100 nhip panic lien tiep, khuvuc.go) da dung han actor cua khu vuc. Tu luc
// do AnhChupNguoiChoiCoLoi tra ErrActorDaDung cho MOI nguoi trong khu vuc do
// — truoc ban nay tat ca deu ket cung o hop thoai "Xin cho" ngay khi bam
// cong. O day dung Dung() de dat actor vao dung trang thai do mot cach chac
// chan, khong phai bom 100 panic that.
func TestTichHop_DoiMapKhiActorDaDungVanMoKhoaClient(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	soKhuVuc := moSoKhuVucKhongQuaiMacDinh(t)
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTK := "stopa" + hauToNganNgau(t)
	taiKhoan := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "Stopa" + hauToNganNgau(t)

	cl := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer cl.dong()
	dangNhapTaoNhanVatVaoTheGioi(t, cl, tenTK, matKhau, tenNV)
	cl.baoVaoKhuVuc()

	nv, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoan)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan that bai: %v", err)
	}
	k := soKhuVuc.Tra(1)
	choNguoiVaoKhuVuc(t, k, nv.ID)

	if err := k.Dung(); err != nil {
		t.Fatalf("khong dung duoc actor khu vuc 1: %v", err)
	}

	cl.guiDoiMap()

	m46 := cl.doiGoiLenh(t, lenhDatLaiViTri)
	d := moDocGoi(t, m46)
	// Khong con actor de hoi thi vi tri duy nhat may chu con biet la ban ghi
	// trong phien — lech mot chut con hon treo (client tu chinh qua nhanh
	// `cy - y != 0` cua Char.setResetPoint).
	d.sh("x", 100)
	d.sh("y", 384)
	d.hetGoi("goi 46")
	if !cl.khongConGoiGi(t) {
		t.Fatal("sau goi 46 con goi nua — khong chup duoc trang thai thi TUYET DOI khong duoc doi map")
	}
}

// TestTichHop_CongTroToiMapChuaCoKhuVucVanMoKhoaClient khoa chot "map dich co
// ton tai khong" (doimap.go). Bo chot do di thi b.khuVuc thanh nil roi
// b.khuVuc.DanhSachQuaiSnapshot() panic — nghia la mot dong CAU HINH sai
// trong bang DiemDi lam SAP ca goroutine phien, chu khong chi lam mot nguoi
// choi khong di duoc.
//
// Dung mot SoKhuVuc CHI CO map 1 de dung dung tinh huong "cong tro toi map
// chua duoc mo khu vuc" ma khong phai dong vao danh muc ban do he thong.
func TestTichHop_CongTroToiMapChuaCoKhuVucVanMoKhoaClient(t *testing.T) {
	khoThat := moKhoThatChoTest(t)

	// Chi mo khu vuc cho map 1 — map 2 (dich cua cong torii) khong co khu vuc.
	soKhuVuc := &SoKhuVuc{theoMap: make(map[int32]*KhuVuc)}
	k1 := MoKhuVuc(TimBanDo(1))
	k1.quai = nil // AN TOAN: goroutine vongLap chua chay (BatDau o duoi)
	k1.BatDau()
	soKhuVuc.theoMap[1] = k1
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTK := "nozon" + hauToNganNgau(t)
	taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "Nozon" + hauToNganNgau(t)

	cl := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer cl.dong()
	dangNhapTaoNhanVatVaoTheGioi(t, cl, tenTK, matKhau, tenNV)
	cl.baoVaoKhuVuc()

	cl.guiDiChuyen(441, 408)
	cl.guiDoiMap()

	m46 := cl.doiGoiLenh(t, lenhDatLaiViTri)
	d := moDocGoi(t, m46)
	d.sh("x", 441)
	d.sh("y", 408)
	d.hetGoi("goi 46")
	if !cl.khongConGoiGi(t) {
		t.Fatal("cong tro toi map chua co khu vuc ma van gui them goi — co the da doi map vao hu vo")
	}
}

// TestTichHop_DoiMapLuuMapMoiXuongCSDLNgay khoa dong `b.thoGhi.Xep(*b.nhanVat)`
// trong xuLyDoiMap.
//
// VI SAO PHAI CO: tu luc may chu ghi mapId moi vao ban ghi cho toi luc client
// gui -39, nguoi choi KHONG nam trong khu vuc nao — ca duong luu dinh ky
// (luudinhky.go duyet theo khu vuc) lan duong luu khi dong phien deu khong
// thay ho. Mat ket noi dung khe do ma thieu dong Xep() thi lan dang nhap sau
// nguoi choi hien lai o LANG CU. Dong do co 8 dong chu thich bien ho nhung
// truoc ban nay khong test nao chung minh no con day.
//
// Test CO Y khong gui -39: no kiem tra dung khe cua so do.
func TestTichHop_DoiMapLuuMapMoiXuongCSDLNgay(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	soKhuVuc := moSoKhuVucKhongQuaiMacDinh(t)
	defer soKhuVuc.DungTatCa()

	tho := MoThoGhi(khoThat, 2)
	defer func() {
		if err := tho.DongVaChoXong(context.Background()); err != nil {
			t.Errorf("dong tho ghi that bai: %v", err)
		}
	}()

	matKhau := "matkhau123"
	tenTK := "luuma" + hauToNganNgau(t)
	taiKhoan := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "Luuma" + hauToNganNgau(t)

	// Dung duong moPhienDangNhapTrungChoTest vi no nhan mot *BoXuLy da chuan
	// bi san — moPhienClientGiaKhuVuc tu tao BoXuLy ben trong nen khong cach
	// nao gan tho ghi vao.
	b := NewBoXuLy(khoThat, soKhuVuc)
	b.DatThoGhi(tho)
	cl := moPhienDangNhapTrungChoTest(t, b)
	defer cl.dong()
	dangNhapTaoNhanVatVaoTheGioi(t, cl, tenTK, matKhau, tenNV)
	cl.baoVaoKhuVuc()

	cl.guiDiChuyen(441, 408)
	cl.guiDoiMap()
	cl.doiGoiLenh(t, lenhXoaMap)
	cl.doiGoiLenh(t, lenhMapInfo)

	// CHUA gui -39: nguoi choi dang o "khe" khong thuoc khu vuc nao. Ban ghi
	// trong CSDL VAN phai da sang map 2. Cho co han vi tho ghi cham dia tren
	// mot goroutine rieng (thogi.go) — day la doi do TRE, khong phai doan mo.
	hetHan := time.Now().Add(5 * time.Second)
	for {
		nv, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoan)
		if err != nil {
			t.Fatalf("LayNhanVatTheoTaiKhoan that bai: %v", err)
		}
		if nv.MapID == 2 {
			if nv.ViTriX != 100 || nv.ViTriY != 384 {
				t.Errorf("vi tri da luu = (%d,%d), mong doi diem ha canh (100,384)", nv.ViTriX, nv.ViTriY)
			}
			return
		}
		if time.Now().After(hetHan) {
			t.Fatalf("sau 5s ban ghi trong CSDL van o mapId %d — doi map khong luu ngay, mat ket noi luc nay la nguoi choi hien lai o lang cu", nv.MapID)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// TestTichHop_DoiMapXetCongBangViTriLayTuActor khoa buoc "chup anh trang thai
// tu actor" — cu the la chuyen actor RUT CAN hang gameplay truoc khi tra loi
// (xuLyGameplayDaXepTruoc, khuvuc.go).
//
// Kich ban: actor dang ban (o day chan co kiem soat bang suKienChay; ngoai
// doi that la mot nhip nang hoac mot dot su kien don day), buoc di toi cong
// DA nam trong hop thu nhung CHUA duoc xu ly, roi goi -23 toi. Doc thang
// b.nhanVat.ViTriX luc nay se thay vi tri CU (ngoai cong) va tra ve goi 46 —
// nguoi choi bam cong ma khong di duoc, khong hieu tai sao. Chup anh qua hop
// thu dieu khien thi actor buoc phai xu ly not buoc di truoc khi tra loi, nen
// vi tri nhan duoc la vi tri THAT.
func TestTichHop_DoiMapXetCongBangViTriLayTuActor(t *testing.T) {
	khoThat := moKhoThatChoTest(t)
	ctx := context.Background()

	soKhuVuc := moSoKhuVucKhongQuaiMacDinh(t)
	defer soKhuVuc.DungTatCa()

	matKhau := "matkhau123"
	tenTK := "actor" + hauToNganNgau(t)
	taiKhoan := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
	tenNV := "Actor" + hauToNganNgau(t)

	cl := moPhienClientGiaKhuVuc(t, khoThat, soKhuVuc)
	defer cl.dong()
	dangNhapTaoNhanVatVaoTheGioi(t, cl, tenTK, matKhau, tenNV)
	cl.baoVaoKhuVuc()

	nv, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoan)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan that bai: %v", err)
	}
	k := soKhuVuc.Tra(1)
	choNguoiVaoKhuVuc(t, k, nv.ID)

	// Chan actor lai — bat tay hai chieu de CHAC CHAN no dang dung trong ham
	// nay chu khong phai "sap toi".
	daVao := make(chan struct{})
	choTha := make(chan struct{})
	k.Gui(suKienChay{f: func() {
		close(daVao)
		<-choTha
	}})
	<-daVao

	// Buoc di toi giua cong: su kien nam trong hop thu, actor dang bi chan nen
	// CHUA ai xu ly no -> ban ghi nhan vat van dang o (100,384).
	cl.guiDiChuyen(441, 408)
	// -23 den ngay sau. Luong phien se dung lai o buoc chup anh (hop thu dieu
	// khien) cho toi khi actor duoc tha.
	cl.guiDoiMap()

	// Doi cho yeu cau chup anh THAT SU vao hop thu dieu khien roi moi tha
	// actor — tha som thi actor kip xu ly buoc di TRUOC khi -23 toi, luc do
	// test se "xanh gia". Nhip cho rat ngan vi han cho control chi 100ms
	// (thoiGianChoControl).
	hetHan := time.Now().Add(2 * time.Second)
	for len(k.hopThuDieuKhien) == 0 && time.Now().Before(hetHan) {
		time.Sleep(200 * time.Microsecond)
	}
	close(choTha)

	// Actor tha ra: rut can hang gameplay (ap dung buoc di) ROI moi tra loi
	// chup anh -> vi tri la (441,408) -> dung trong cong -> doi map that.
	if m := cl.doiGoiLenh(t, lenhXoaMap); m.Command != lenhXoaMap {
		t.Fatalf("mong doi goi -22, nhan %d", m.Command)
	}
	m24 := cl.doiGoiLenh(t, lenhMapInfo)
	d := moDocGoi(t, m24)
	d.b("mapId", 2)
}
