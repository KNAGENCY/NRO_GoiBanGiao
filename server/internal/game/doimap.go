// Package game — doi ban do qua cong (waypoint).
//
// TAI SAO CO FILE NAY: gui cong xuong client (MAP INFO, xem ghiDiemDi trong
// thegioi.go) MA KHONG xu ly goi -23 thi con TE HON la khong co cong. Client
// dat `Char.ischangingMap = true` + `isLockKey = true` + hien hop thoai "Xin
// cho" NGAY khi bam cong (Waypoint.perform idAction 2 -> Service.
// requestChangeMap, Waypoint.cs), va CHI go ba thu do khi may chu tra loi.
// Khong tra loi = nhan vat DUNG HINH VINH VIEN, phai thoat game.
//
// Vi vay MOI duong ra cua ham xuLyDoiMap — tru DUY NHAT mot cai — deu PHAI
// gui mot goi. Ke het TAM CAI ra day de khong ai them duong ra thu chin ma
// quen:
//
//	chua vao khu vuc            -> IM LANG (duong ra DUY NHAT im lang: client
//	                               khong the dang cho, xem chot dau ham)
//	khong chup duoc trang thai  -> 46
//	mapId hien tai ngoai danh muc -> 46
//	nhan vat dang chet          -> 46
//	khong dung trong cong nao   -> 46
//	cong tro toi map khong co   -> 46
//	go khoi khu vuc cu that bai -> 46
//	THANH CONG                  -> -22 (xoa map cu) roi -24 (MAP INFO map moi)
//
// Ban dau chu thich nay chi ke HAI duong ("khong dung trong cong" va "thanh
// cong") trong khi code co BON — hai duong "khong chup duoc trang thai" va
// "mapId ngoai danh muc" tra `return nil` im lang. Do khong phai loi ly
// thuyet: han cho actor tra loi chi 100ms (thoiGianChoControl, khuvuc.go)
// trong khi actor phai rut can hang gameplay ton dong truoc khi tra loi, va
// van an toan con DUNG HAN actor sau 100 nhip panic lien tiep — tu luc do MOI
// nguoi bam cong trong khu vuc do deu roi vao duong "khong chup duoc trang
// thai" va ket cung.
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
//	     (Controller.cs:3137) MO KHOA: isLockKey = false, va setResetPoint
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

	// (1) vi tri authoritative. Buoc nay dat TRUOC MOI CHOT KHAC (ke ca chot
	// "mapId hien tai co trong danh muc khong") CO CHU DICH: no la thu duy
	// nhat cho ra mot toa do hop le de mo khoa client, nen moi duong ra sau no
	// deu co san mot goi 46 dung dan de gui.
	ban, err := b.khuVuc.AnhChupNguoiChoiCoLoi(b.nhanVat.ID)
	if err != nil || len(ban) == 0 {
		// Duong ra 2. Xay ra that khi actor cua khu vuc da bi van an toan dung
		// han (nguongPanicDungKhuVuc, khuvuc.go) hoac dang qua ban den muc vuot
		// han 100ms cua control. Truoc day cho nay `return nil` — tuc la MOI
		// nguoi bam cong trong khu vuc do deu dung hinh vinh vien, phai thoat
		// game.
		//
		// Khong con actor de hoi thi vi tri duy nhat may chu con biet la ban
		// ghi trong phien. Doc no o day la doc mot truong ma actor CO THE dang
		// ghi (ky luat §6.9) — chap nhan CO CHU DICH: doi lai la mot toa do co
		// the lech vai chuc don vi (client tu chinh qua nhanh `cy - y != 0` cua
		// Char.setResetPoint, Char.cs:3204-3216), con khong gui gi la nguoi choi mat
		// phien choi. Tha lech con hon treo.
		slog.Warn("khong chup duoc trang thai de doi map — mo khoa client bang vi tri trong phien",
			"idNhanVat", b.nhanVat.ID, "loi", err)
		return s.XepGui(goiDatLaiViTri(int16(b.nhanVat.ViTriX), int16(b.nhanVat.ViTriY)))
	}
	hienTai := ban[0]
	x, y := int16(hienTai.ViTriX), int16(hienTai.ViTriY)

	bdCu := TimBanDo(hienTai.MapID)
	if bdCu == nil {
		// Duong ra 3: ban ghi nhan vat mang mot mapId khong con trong danh muc
		// (du lieu cu sau khi doi bang ban do, hoac mot ban va lam hep danh
		// muc). Khong tra cuu duoc cong nao ca — nhung van PHAI mo khoa, va o
		// day da co vi tri authoritative de lam viec do.
		slog.Error("nhan vat dang o mapId khong co trong danh muc ban do — khong tra cuu duoc cong",
			"idNhanVat", b.nhanVat.ID, "mapId", hienTai.MapID)
		return s.XepGui(goiDatLaiViTri(x, y))
	}

	// Nhan vat dang chet -> tu choi, nhung VAN PHAI mo khoa client.
	if hienTai.HpHienTai <= 0 {
		return s.XepGui(goiDatLaiViTri(x, y))
	}

	wp := timDiemDi(bdCu, x, y)
	if wp == nil {
		// Khong dung trong cong nao. Xay ra khi client va may chu lech vi tri
		// (nguoi choi bam dung luc goi di chuyen chua toi), hoac khi mot client
		// sua doi tu gui -23. Mo khoa tai cho, KHONG doi map.
		return s.XepGui(goiDatLaiViTri(x, y))
	}
	bdMoi := TimBanDo(wp.DenMap)
	if bdMoi == nil || b.soKhuVuc == nil || b.soKhuVuc.Tra(wp.DenMap) == nil {
		// Cong tro toi map khong co that / chua co khu vuc. Day la loi CAU
		// HINH cua may chu (bang DiemDi sai), khong phai loi nguoi choi — ghi
		// log muc Error vi no khong bao gio duoc phep xay ra o ban chay that,
		// va van mo khoa cho nguoi choi choi tiep.
		slog.Error("cong roi map tro toi map khong ton tai",
			"tuMap", bdCu.MapID, "denMap", wp.DenMap, "ten", wp.Ten)
		return s.XepGui(goiDatLaiViTri(x, y))
	}

	// (2) roi khu vuc cu — DONG BO, cho actor xu ly xong.
	if err := b.khuVuc.RoiKhuVuc(b.nhanVat.ID); err != nil {
		// Khong go duoc thi TUYET DOI khong duoc di tiep: nhan vat se nam o
		// hai khu vuc cung luc, hai actor cung ghi mot con tro.
		slog.Error("khong go duoc nguoi choi khoi khu vuc khi doi map — huy doi map",
			"idNhanVat", b.nhanVat.ID, "loi", err)
		return s.XepGui(goiDatLaiViTri(x, y))
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
	//
	// XepGui (hang doi FIFO cua phien) chu KHONG PHAI Gui (ghi thang socket).
	// Ly do KHONG phai an toan luong — Gui co mutex nen khung goi khong bao gio
	// vo — ma la THU TU: actor phat moi goi gameplay bang XepGui, nen mot goi
	// -5 (nguoi khac xuat hien o map CU) vua duoc xep vao hang van con nam do
	// cho goroutine writer. Ghi thang -22/-24 la NHAY HANG len truoc no: client
	// tai xong map moi, roi writer moi xa goi -5 cu ra — va vi -24 da goi
	// resetAllvector() lam vCharInMap rong, goi -5 den sau lai them nguoi do
	// vao. Ket qua: mot nhan vat MA cua Lang La dung giua Lang Cat.
	//
	// Cung ly do do, MOI goi 46 o cac duong ra ben tren cung dung XepGui: sau
	// khi vao the gioi thi doimap.go la cho DUY NHAT con ghi thang socket, bo
	// no di la ca duong hau-dang-nhap chay chung mot hang FIFO, khong con cho
	// nao chen ngang duoc nua.
	if err := s.XepGui(goiXoaMap()); err != nil {
		return err
	}
	mapInfo, err := goiMapInfo(b.nhanVat, bdMoi, b.khuVuc.DanhSachQuaiSnapshot())
	if err != nil {
		return err
	}
	if err := s.XepGui(mapInfo); err != nil {
		return err
	}

	slog.Info("doi map qua cong",
		"ten", b.nhanVat.Ten, "id", b.nhanVat.ID,
		"tuMap", bdCu.MapID, "denMap", bdMoi.MapID, "cong", wp.Ten)
	return nil
}
