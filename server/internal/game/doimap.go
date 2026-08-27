// Package game — doi ban do qua cong (waypoint).
//
// TAI SAO CO FILE NAY: gui cong xuong client (MAP INFO, xem ghiDiemDi trong
// thegioi.go) MA KHONG xu ly goi -23 thi con TE HON la khong co cong. Client
// dat `Char.ischangingMap = true` + `isLockKey = true` + hien hop thoai "Xin
// cho" NGAY khi bam cong (Waypoint.perform idAction 2 -> Service.
// requestChangeMap, Waypoint.cs), va CHI go ba thu do khi may chu tra loi.
// Khong tra loi = nhan vat DUNG HINH VINH VIEN, phai thoat game.
//
// Vi vay hai duong ra cua ham xuLyDoiMap deu PHAI gui mot goi:
//
//	dung trong cong  -> -22 (xoa map cu) roi -24 (MAP INFO map moi)
//	khong o cong nao -> 46 (dat lai vi tri) de mo khoa client tai cho
package game

import (
	"log/slog"

	mynet "nroserver/internal/net"
	"nroserver/internal/protocol"
)

// Ba ma lenh cua duong doi map. Gia tri lay tu ban Java cu, DUNG DOI:
//
//	-23  client xin doi map (than RONG). Controller.java:561 dieu ve
//	     ChangeMapService.changeMapWaypoint. Client gui o Service.cs:801.
//	-22  Service.clearMap (Service.java:573) — than RONG. Client (Controller.
//	     cs:3538) xoa muc tieu dang ngam, reset nut, khoa phim cho map moi.
//	 46  Service.resetPoint (Service.java:558) — short x, short y. Client
//	     (Controller.cs:3136) MO KHOA: isLockKey = false, va setResetPoint
//	     (Char.cs:3199) an hop thoai cho + dat ischangingMap = false khi y
//	     khong doi.
const (
	lenhDoiMap      int8 = -23
	lenhXoaMap      int8 = -22
	lenhDatLaiViTri int8 = 46
)

// goiXoaMap dung goi -22: than RONG. Gui TRUOC -24 khi doi map, dung thu tu
// ban Java (ChangeMapService.changeMap: goToMap -> clearMap -> mapInfo).
func goiXoaMap() *protocol.Message {
	return protocol.NewMessage(lenhXoaMap, protocol.NewWriter())
}

// goiDatLaiViTri dung goi 46: dat nhan vat ve dung (x, y) va MO KHOA client.
//
// Dung lam duong thoat khi -23 den ma nguoi choi khong dung trong cong nao.
// Truyen dung vi tri HIEN TAI cua nguoi choi thi client di vao nhanh
// `cy - y == 0` cua Char.setResetPoint — an hop thoai cho, tha khoa phim,
// KHONG lam nhan vat truot di dau ca.
func goiDatLaiViTri(x, y int16) *protocol.Message {
	w := protocol.NewWriter()
	w.Short(x)
	w.Short(y)
	return protocol.NewMessage(lenhDatLaiViTri, w)
}

// timDiemDi tra ve cong ma (x, y) dang dung trong, hoac nil.
//
// Lay CAI DAU TIEN khop, giong ban Java (MapService.getWaypointPlayerIn dung
// findFirst). Hien moi ban do chi co mot cong nen khong co gi de tranh chap.
func timDiemDi(bd *BanDo, x, y int16) *DiemDi {
	if bd == nil {
		return nil
	}
	for i := range bd.DiemDi {
		if bd.DiemDi[i].ChuaTrong(x, y) {
			return &bd.DiemDi[i]
		}
	}
	return nil
}

// xuLyDoiMap: client bam cong va gui -23 (than rong).
//
// TRINH TU BAT BUOC — doi thu tu la mat du lieu hoac tranh chap:
//
//  1. CHUP anh trang thai tu actor (AnhChupNguoiChoiCoLoi). Phai lam TRUOC
//     buoc 2 vi RoiKhuVuc xoa nguoi choi khoi khu vuc, chup sau se tra rong
//     (cung cai bay da ghi ky o KhiDong, login.go).
//     Vi tri trong anh chup nay la vi tri AUTHORITATIVE — do la thu duy nhat
//     duoc phep dung de xet "co dung trong cong khong". Doc thang
//     b.nhanVat.ViTriX o luong phien la doc mot truong ma actor DANG GHI.
//  2. RoiKhuVuc (dong bo, cho actor xu ly xong). Sau khi ham nay tra ve,
//     actor khong con giu con tro nhan vat nay nua -> luong phien so huu lai
//     b.nhanVat va duoc phep GHI vao no. Truoc do thi TUYET DOI khong.
//  3. Ghi map/vi tri moi, doi b.khuVuc, ha co daVaoKhuVuc.
//  4. Gui -22 roi -24. Client tai map moi roi tu gui -39, luc do
//     xuLyVaoKhuVuc moi dua nhan vat vao khu vuc MOI — dung y het duong vao
//     the gioi lan dau, khong co duong rieng nao ca.
//
// KHONG doi map khi nhan vat dang chet: client dang o man hinh hoi sinh,
// day no sang map khac giua chung se lam hai ben lech trang thai. Ban Java
// cung chan doi map trong hang loat tinh huong tuong tu.
func (b *BoXuLy) xuLyDoiMap(s *mynet.Session) error {
	// Chua vao khu vuc = chua choi thuc su (dang o man tao nhan vat, hoac goi
	// -23 den truoc -39). Bo qua im lang: client khong the dang cho o hop
	// thoai "Xin cho" trong trang thai do, nen khong can mo khoa gi.
	if b.nhanVat == nil || b.khuVuc == nil || !b.daVaoKhuVuc {
		return nil
	}

	bdCu := TimBanDo(b.nhanVat.MapID)
	if bdCu == nil {
		return nil
	}

	// (1) vi tri authoritative
	ban, err := b.khuVuc.AnhChupNguoiChoiCoLoi(b.nhanVat.ID)
	if err != nil || len(ban) == 0 {
		slog.Warn("khong chup duoc trang thai de doi map", "idNhanVat", b.nhanVat.ID, "loi", err)
		return nil
	}
	hienTai := ban[0]
	x, y := int16(hienTai.ViTriX), int16(hienTai.ViTriY)

	// Nhan vat dang chet -> tu choi, nhung VAN PHAI mo khoa client.
	if hienTai.HpHienTai <= 0 {
		return s.Gui(goiDatLaiViTri(x, y))
	}

	wp := timDiemDi(bdCu, x, y)
	if wp == nil {
		// Khong dung trong cong nao. Xay ra khi client va may chu lech vi tri
		// (nguoi choi bam dung luc goi di chuyen chua toi), hoac khi mot client
		// sua doi tu gui -23. Mo khoa tai cho, KHONG doi map.
		return s.Gui(goiDatLaiViTri(x, y))
	}
	bdMoi := TimBanDo(wp.DenMap)
	if bdMoi == nil || b.soKhuVuc == nil || b.soKhuVuc.Tra(wp.DenMap) == nil {
		// Cong tro toi map khong co that / chua co khu vuc. Day la loi CAU
		// HINH cua may chu (bang DiemDi sai), khong phai loi nguoi choi — ghi
		// log muc Error vi no khong bao gio duoc phep xay ra o ban chay that,
		// va van mo khoa cho nguoi choi choi tiep.
		slog.Error("cong roi map tro toi map khong ton tai",
			"tuMap", bdCu.MapID, "denMap", wp.DenMap, "ten", wp.Ten)
		return s.Gui(goiDatLaiViTri(x, y))
	}

	// (2) roi khu vuc cu — DONG BO, cho actor xu ly xong.
	if err := b.khuVuc.RoiKhuVuc(b.nhanVat.ID); err != nil {
		// Khong go duoc thi TUYET DOI khong duoc di tiep: nhan vat se nam o
		// hai khu vuc cung luc, hai actor cung ghi mot con tro.
		slog.Error("khong go duoc nguoi choi khoi khu vuc khi doi map — huy doi map",
			"idNhanVat", b.nhanVat.ID, "loi", err)
		return s.Gui(goiDatLaiViTri(x, y))
	}
	b.daVaoKhuVuc = false

	// (3) Tu day tro di luong phien so huu b.nhanVat (actor da nha ra o buoc
	// 2) va duoc phep ghi vao no.
	//
	// CHI ghi ba truong map/vi tri, KHONG chep de `*b.nhanVat = hienTai`:
	// actor sua nhan vat qua CHINH con tro nay (k.nguoiTrong giu dung
	// b.nhanVat, xem suKienVao trong khuvuc.go), nen mau/chakra/kinh nghiem
	// da la ban moi nhat san roi. Chep anh chup de len se DAY NGUOC lai moi
	// thay doi xay ra giua luc chup (buoc 1) va luc roi khu vuc (buoc 2) —
	// vi du mot don quai danh trung dung khe do se bi xoa. Kiem thu dot bien
	// da bat duoc dong chep do la vo dung, va doc ky thi no con co hai.
	b.nhanVat.MapID = wp.DenMap
	b.nhanVat.ViTriX = int32(wp.DenX)
	b.nhanVat.ViTriY = int32(wp.DenY)
	b.khuVuc = b.soKhuVuc.Tra(wp.DenMap)

	// (3b) Luu NGAY. Tu luc nay toi luc client gui -39, nguoi choi KHONG nam
	// trong khu vuc nao — ca hai duong luu dinh ky (luudinhky.go duyet theo
	// khu vuc) lan duong luu khi dong phien (KhiDong, login.go, tu thoat som
	// khi !daVaoKhuVuc) deu khong thay ho. Mat ket noi dung khe do ma khong co
	// dong nay thi nguoi choi dang nhap lai o MAP CU: khong mat du lieu,
	// nhung "vua di qua cong xong lai thay minh o lang cu" la loi nguoi choi
	// bao duoc ma khong ai lan ra duoc nguyen nhan. Xep() khong bao gio chan
	// (thogi.go) nen an toan goi tren luong phien.
	if b.thoGhi != nil {
		b.thoGhi.Xep(*b.nhanVat)
	}

	// (4) -22 roi -24.
	if err := s.Gui(goiXoaMap()); err != nil {
		return err
	}
	mapInfo, err := goiMapInfo(b.nhanVat, bdMoi, b.khuVuc.DanhSachQuaiSnapshot())
	if err != nil {
		return err
	}
	if err := s.Gui(mapInfo); err != nil {
		return err
	}

	slog.Info("doi map qua cong",
		"ten", b.nhanVat.Ten, "id", b.nhanVat.ID,
		"tuMap", bdCu.MapID, "denMap", bdMoi.MapID, "cong", wp.Ten)
	return nil
}
