// dovuong.go — KH7 Task 11: do VUONG duoi dat trong mot khu vuc + duong NHAT (`-20`).
//
// Giai ma hai phia (PHA A):
// docs/superpowers/research/2026-09-03-kh7-giai-ma-do-roi-nhat.md
// Dac ta: docs/superpowers/specs/2026-09-03-kh7-vat-pham-design.md
// D8 (nhat mot buoc), §4.6/§4.6b (loi L2/L4), §4.6c, §2.3, §5.3.
// Ky luat actor: BAT-DAU-TU-DAY.md §6.9/§6.10b.
//
// ---------------------------------------------------------------------------
// LUONG NHAT — ai cham cai gi (day la phan de sai nhat cua ca task)
// ---------------------------------------------------------------------------
//
//	goroutine PHIEN         goroutine ACTOR khu vuc        goroutine BO NHAT DO
//	(login.go, xuLyNhat)    (khuvuc.go vongLap)            (chay o day, duoc cho o dia)
//	------------------      ------------------------       ---------------------------
//	doc short itemMapId
//	  -> k.Gui(suKienXinNhat)
//	                        mon con trong k.doDuoiDat?
//	                        nguoi con trong khu vuc?
//	                        con trong tam nhat?
//	                        -> GO MON KHOI k.doDuoiDat NGAY   <-- chot bo nho
//	                        -> boNhat.Xep(yeuCauNhat{mon})  --->
//	                                                        LayTui -> tim o trong
//	                                                        store.NhatLen (MOT UPDATE
//	                                                          co dieu kien, D8)
//	                        <--- k.Gui(suKienNhatXong{...}) <---
//	                        thanh cong: phat -20 rieng nguoi nhat
//	                                    + -19 cho nguoi khac
//	                        that bai:  TRA MON LAI k.doDuoiDat
//	                                    + bao loi rieng nguoi nhat
//
// HAI LOP CHONG NHAN BAN, khong phai mot:
//
//	Lop 1 (bo nho, nhanh): actor go mon khoi k.doDuoiDat NGAY khi nhan yeu cau
//	  hop le. Hai nguoi bam cung mot phan nghin giay -> hai su kien vao CUNG mot
//	  hop thu, actor xu ly TUAN TU, nguoi thu hai khong con thay mon trong map ->
//	  chi MOT yeu cau CSDL duoc xep.
//	Lop 2 (CSDL, dung tuyet doi): `WHERE noi_chua='DUOI_DAT'` trong chinh cau
//	  UPDATE (store.NhatLen). Ke ca khi tang Go co bug — hai tien trinh may chu,
//	  panic khoi phuc sai, ai do goi thang NhatLen — CSDL van chi cho MOT nguoi
//	  thang. Lop 1 khong thay the duoc lop 2 va nguoc lai.
//
// VI SAO ACTOR KHONG TU GOI store.NhatLen: §6.10b cam tuyet doi — vongLap cua
// khu vuc khong bao gio duoc cho o dia. Mot cau UPDATE cham 50ms se lam DUNG
// HINH ca ban do voi moi nguoi trong do. Actor chi cham map trong bo nho (vai
// chuc nano-giay) roi XEP yeu cau, giong het cach xepNoThuong lam voi thuong boss.
//
// VI SAO GOI -20/-19 DO ACTOR PHAT chu khong phai goroutine CSDL: danh sach
// nguoiTrong thuoc so huu cua vongLap (§1.6). Goroutine CSDL khong duoc duyet no.
// No chi bao "xong roi, ket qua the nay" bang mot su kien nguoc ve hop thu.
package game

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// ---------------------------------------------------------------------------
// 1. Ma lenh — doi chieu truc tiep voi `case` cua client
// ---------------------------------------------------------------------------

const (
	// lenhDoXuatHien: mot mon hien duoi dat. Client Controller.cs:3511 (case 68).
	lenhDoXuatHien int8 = 68
	// lenhNhatDuoc: "ban nhat duoc". Client Controller.cs:3395 (case -20).
	lenhNhatDuoc int8 = -20
	// lenhNguoiKhacNhat: nguoi khac vua nhat. Client Controller.cs:3475 (case -19).
	lenhNguoiKhacNhat int8 = -19
	// lenhDoBienMat: mon het gio, bien mat. Client Controller.cs:3381 (case -21).
	lenhDoBienMat int8 = -21
	// lenhNhanTinRieng: mot dong chu bao cho MOT nguoi (than goi: dung mot UTF).
	// Java Service.java:1354-1366 (`sendThongBao`) <-> client Controller.cs:2988-2990.
	// KHAC lenhThongBaoToanServer (93, thongbao.go) — 93 chay ngang cho ca may chu.
	lenhNhanTinRieng int8 = -25
)

// ---------------------------------------------------------------------------
// 2. Tran cung + hang so lay tu PHA A
// ---------------------------------------------------------------------------

const (
	// doVuongKhongChu la gia tri truong thu 5 cua goi 68 — "mon nay ai cung
	// nhat duoc". Ban Java ghi cung so 3 (Service.java:1622) ma khong tra duoc
	// y nghia; ban Go dung -1 vi do chinh la quy uoc "khong co chu" ma may chu
	// Java dung o cho kiem quyen nhat (Zone.java:450). Client KHONG phan nhanh
	// gia tri nay — gia tri dac biet DUY NHAT no doc them byte la -2 (hao quang,
	// Controller.cs:3519-3522), nen -1 an toan tuyet doi.
	doVuongKhongChu int32 = -1

	// doVuongSoLuongToiDa: client doc so luong bang readShort CO DAU. Ban Java
	// ghi cung 9999 khi tran (Zone.java:498) -> nguoi choi thay SAI so ma khong
	// bao gi (loi L4). Ban Go KHONG chep meo do: tran thi tra loi, khong gui goi.
	doVuongSoLuongToiDa = 32767

	// doVuongSoOTui: so o tui cua ban 1. Phai < tranSoO (127) vi con so nay con
	// di xuong ME_INFO doc bang byte CO DAU (goihanhtrang.go). Dat 30 lam khoi
	// dau — con so can bang thuoc KH8.
	doVuongSoOTui = 30

	// tamNhatBinhPhuong: binh phuong tam nhat toi da (pixel). So binh phuong de
	// tranh sqrt/float, dung cach lam voi KiemDanhHopLe (chotkiem.go) va
	// kiemDiChuyenHopLe (dichuyen.go). 120px la khoang rong rai — client chi cho
	// bam nhat mon dang o gan, kiem nay la LUOI CHONG GIAN LAN (goi -20 tu tay),
	// khong phai co che gameplay, nen khong siet chat toi muc nguoi choi that bi
	// tu choi oan vi tre mang.
	tamNhatBinhPhuong int64 = 120 * 120
)

// Loai vat pham (Item.type) can GHI HAI SHORT o goi -20 — day chinh la loi L2
// cua ban Java (client doc hai, Java ghi mot -> tien khong cong, im lang).
const (
	doVuongLoaiVang     int8 = 9
	doVuongLoaiNgoc     int8 = 10
	doVuongLoaiLuongKho int8 = 34
)

// doVuongGhiHaiShort tra true voi cac loai ma client doc HAI short so luong
// (Controller.cs:3413-3430). Tach thanh mot ham rieng de test dot bien "ghi mot
// short cho type 9/10" co dung MOT cho de doi.
func doVuongGhiHaiShort(loai int8) bool {
	return loai == doVuongLoaiVang || loai == doVuongLoaiNgoc || loai == doVuongLoaiLuongKho
}

// ---------------------------------------------------------------------------
// 3. Bo dung goi
// ---------------------------------------------------------------------------

// goiDoXuatHien dung goi 68 — mot mon hien duoi dat, ai cung nhat duoc.
//
//	[short itemMapId][short mauId][short x][short y][int -1]
//
// KHONG co truong hao quang: truong do chi ton tai khi truong thu 5 == -2
// (do dai bien thien, bay §6.6). Mon co hao quang can mot ham dung RIENG, khong
// them co vao ham nay.
func goiDoXuatHien(itemMapId, mauId, x, y int16) *protocol.Message {
	w := protocol.NewWriter()
	w.Short(itemMapId)
	w.Short(mauId)
	w.Short(x)
	w.Short(y)
	w.Int(doVuongKhongChu)
	return protocol.NewMessage(lenhDoXuatHien, w)
}

// goiNhatDuoc dung goi -20 gui RIENG nguoi vua nhat duoc.
//
//	[short itemMapId][UTF loiNhan][short soLuong]        loai thuong
//	[short itemMapId][UTF loiNhan][short soLuong][short soLuong]  loai 9/10/34
//
// ⚠️ SHORT THU HAI la ca cai lam nen loi L2. Client voi type 9/10/34 doc HAI
// short lien tiep va CHI cong tien bang cai THU HAI (Controller.cs:3413-3430:
// `num = readShort()` lan hai roi `xu += num`). Thieu no thi client doc hut,
// nem loi, bi `catch (Exception) { }` RONG nuot -> tien khong cong, khong mot
// dau hieu nao. Dung xoa dong do.
//
// loiNhan RONG co y nghia: client hien CHU BAY `+N` mau vang/xanh/do thay vi
// mot dong chu (Controller.cs:3435-3451). Tien nen de rong.
//
// Tra loi va KHONG dung goi khi soLuong ngoai [0, 32767] (L4) hoac loiNhan qua
// dai — dung xong toan bo du lieu roi moi ghi byte (§1.10/L13).
func goiNhatDuoc(itemMapId int16, loai int8, loiNhan string, soLuong int32) (*protocol.Message, error) {
	if soLuong < 0 || soLuong > doVuongSoLuongToiDa {
		return nil, &loiDoVuong{"goi -20: so luong " + itoa(soLuong) + " ngoai [0,32767]"}
	}
	w := protocol.NewWriter()
	w.Short(itemMapId)
	if err := w.UTF(loiNhan); err != nil {
		return nil, err
	}
	w.Short(int16(soLuong))
	if doVuongGhiHaiShort(loai) {
		w.Short(int16(soLuong))
	}
	return protocol.NewMessage(lenhNhatDuoc, w), nil
}

// goiNguoiKhacNhat dung goi -19 — bao NHUNG NGUOI KHAC rang mon do vua bi ai
// nhat, de client ho cho mon bay ve phia nguoi do roi tu xoa.
//
//	[short itemMapId][int idNguoiNhat]
//
// ⚠️ Client `return` ngay neu khong tim thay nhan vat do trong map
// (Controller.cs:3486-3489) — nghia la mon se KET LAI tren man hinh nguoi khac
// vinh vien. Vi vay goi nay phai duoc phat NGAY trong nhip actor xu ly ket qua
// nhat, khong duoc hoan lai.
func goiNguoiKhacNhat(itemMapId int16, idNguoiNhat int32) *protocol.Message {
	w := protocol.NewWriter()
	w.Short(itemMapId)
	w.Int(idNguoiNhat)
	return protocol.NewMessage(lenhNguoiKhacNhat, w)
}

// goiDoBienMat dung goi -21 — mon HET GIO, bien mat khong hoat anh.
//
//	[short itemMapId]
//
// KHONG dung sau khi nhat: client tu xoa mon sau -20/-19 (setPoint -> status=2
// -> update() removeElement, ItemMap.cs:96-115). Xem PHA A §0.
func goiDoBienMat(itemMapId int16) *protocol.Message {
	w := protocol.NewWriter()
	w.Short(itemMapId)
	return protocol.NewMessage(lenhDoBienMat, w)
}

// goiNhanTinRieng dung goi -25 — mot dong chu cho RIENG mot nguoi.
//
//	[UTF noiDung]
//
// Tra nil khi noiDung qua dai (Writer.UTF ghi do dai 2 byte khong dau) — nguoi
// goi PHAI kiem nil, dung tinh than voi goiThongBaoToanServer (thongbao.go).
func goiNhanTinRieng(noiDung string) *protocol.Message {
	w := protocol.NewWriter()
	if err := w.UTF(noiDung); err != nil {
		slog.Error("tin rieng qua dai, bo qua goi -25", "loi", err)
		return nil
	}
	return protocol.NewMessage(lenhNhanTinRieng, w)
}

// loiDoVuong la kieu loi rieng cua file nay — du de phan biet khi test, khong
// keo them phu thuoc fmt vao duong nong.
type loiDoVuong struct{ s string }

func (e *loiDoVuong) Error() string { return e.s }

// itoa doi int32 sang chuoi (chi dung trong thong bao loi).
func itoa(v int32) string {
	if v == 0 {
		return "0"
	}
	am := v < 0
	if am {
		v = -v
	}
	var b [12]byte
	i := len(b)
	for v > 0 {
		i--
		b[i] = byte('0' + v%10)
		v /= 10
	}
	if am {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

// ---------------------------------------------------------------------------
// 4. Trang thai "do duoi dat" — THUOC ACTOR khu vuc
// ---------------------------------------------------------------------------

// doVuong la MOT mon dang nam duoi dat trong khu vuc. Ai cung nhin thay no, nen
// no thuoc actor khu vuc (§6.9) — dung noi voi hanh trang (chi mot nguoi thay,
// thuoc goroutine phien).
//
// vatPhamID la id THAT trong CSDL; itemMapId (khoa cua map k.doDuoiDat) la mot
// so ngan CAP TRONG KHU VUC chi de noi chuyen voi client — id CSDL KHONG BAO
// GIO xuong client (§4.1).
type doVuong struct {
	vatPhamID int64
	mauID     int32
	soLuong   int32
	x, y      int16
	// loai la MauVatPham.Loai, chup lai luc tha xuong. Giu o day thay vi tra
	// cuu lai luc nhat de goroutine CSDL khong phai dung toi danh muc mau.
	loai int8
	// ten dung cho loi nhan "Ban nhan duoc <ten>". Chup cung luc voi loai.
	ten string
}

// suKienThaDo bao actor: dat MOT mon (da co dong trong CSDL o DUOI_DAT) xuong
// khu vuc. Nguoi gui la goroutine NGOAI actor — task 12 (quai/boss chet) se goi
// no sau khi da ghi dong CSDL xong.
//
// traLoi (co the nil) nhan itemMapId ma actor vua cap — de nguoi goi biet ma
// theo doi mon. Buffer 1 o phia nguoi goi; actor khong bao gio cho.
type suKienThaDo struct {
	mon    doVuong
	traLoi chan int16
}

// suKienXinNhat bao actor: nguoi choi idNhanVat xin nhat mon itemMapId (goi len
// `-20`). Goroutine phien CHI gui su kien nay — khong tra k.doDuoiDat, khong
// cham CSDL. Toan bo viec kiem nam trong vongLap.
type suKienXinNhat struct {
	idNhanVat int64
	itemMapId int16
}

// suKienNhatXong la duong VE tu goroutine CSDL. `loi == nil` nghia la mon da
// doi chu that su trong CSDL (da COMMIT) — chi luc do actor moi phat goi cho
// client ("ghi so truoc — giao sau", §2.3 buoc 7).
//
// `loi != nil` nghia la KHONG co gi doi trong CSDL: actor TRA MON LAI
// k.doDuoiDat (khong ai mat do) va bao loi rieng nguoi xin nhat.
type suKienNhatXong struct {
	idNhanVat int64
	itemMapId int16
	mon       doVuong
	loi       error
}

// suKienDonDoHetGio bao actor: don mon itemMapId khoi mat dat (het gio) va phat
// -21 cho ca khu vuc. Chua co ai goi o task 11 (bo dem gio thuoc task 12) —
// nhung goi -21 va duong don da co san de task 12 khong phai sua file nay.
type suKienDonDoHetGio struct {
	itemMapId int16
}

// capItemMapId cap mot itemMapId moi trong khu vuc. CHI vongLap goi.
//
// Cuon vong khi cham tran short duong: id khong bao gio am (client doc short co
// dau, so am se khong khop mon nao). Bo qua id dang duoc dung — mot khu vuc
// khong the co 32767 mon duoi dat cung luc, nhung vong for nay lam cho dieu do
// KHONG THE gay trung id ke ca khi co.
func (k *KhuVuc) capItemMapId() int16 {
	for i := 0; i < doVuongSoLuongToiDa; i++ {
		k.itemMapIdKe++
		if k.itemMapIdKe < 0 {
			k.itemMapIdKe = 0
		}
		if _, dungRoi := k.doDuoiDat[k.itemMapIdKe]; !dungRoi {
			return k.itemMapIdKe
		}
	}
	return -1 // het id: khu vuc day nghet do — nguoi goi bo qua
}

// thaDoXuongKhuVuc cam mot mon vao mat dat va phat goi 68 cho ca khu vuc.
// CHAY TRONG ACTOR.
func (k *KhuVuc) thaDoXuongKhuVuc(mon doVuong) int16 {
	if k.doDuoiDat == nil {
		k.doDuoiDat = make(map[int16]*doVuong)
	}
	id := k.capItemMapId()
	if id < 0 {
		slog.Warn("khu vuc het itemMapId, bo qua mon roi", "vatPhamID", mon.vatPhamID)
		return -1
	}
	ban := mon
	k.doDuoiDat[id] = &ban
	if mon.mauID < 0 || mon.mauID > 32767 {
		slog.Error("mau_id ngoai mien short, khong phat goi 68", "mauID", mon.mauID)
		return id
	}
	k.phatChoCaKhuVuc(goiDoXuatHien(id, int16(mon.mauID), mon.x, mon.y))
	return id
}

// donDoHetGio go mon khoi mat dat va bao ca khu vuc bang -21. CHAY TRONG ACTOR.
func (k *KhuVuc) donDoHetGio(itemMapId int16) {
	if _, co := k.doDuoiDat[itemMapId]; !co {
		return // da bi nhat truoc do — khong phat gi, khong loi
	}
	delete(k.doDuoiDat, itemMapId)
	k.phatChoCaKhuVuc(goiDoBienMat(itemMapId))
}

// xuLyXinNhat la than cua case suKienXinNhat. CHAY TRONG ACTOR.
//
// Ba kiem, roi GO MON NGAY. Thu tu quan trong: go mon TRUOC khi xep yeu cau —
// neu xep truoc roi moi go, hai su kien lien tiep trong cung mot nhip van co
// the deu thay mon con day va deu xep (dot bien 2).
func (k *KhuVuc) xuLyXinNhat(sk suKienXinNhat) {
	nguoi, co := k.nguoiTrong[sk.idNhanVat]
	if !co || nguoi.nv == nil {
		return // khong con trong khu vuc -> im lang bo qua
	}
	mon, coMon := k.doDuoiDat[sk.itemMapId]
	if !coMon {
		// Nguoi thua cuoc (lop 1). Khong phat -19/-20 gi ca: client cua ho van
		// dang thay mon, nhung goi -19 cua nguoi THANG (phat ngay sau khi CSDL
		// xong) se don no di. Neu nguoi thang loi CSDL thi mon quay lai — dung
		// nhu tren man hinh ho dang thay.
		return
	}
	if nguoi.nv.HpHienTai <= 0 {
		return // dang chet thi khong nhat duoc (giong Controller.java:724)
	}
	dx := int64(nguoi.nv.ViTriX) - int64(mon.x)
	dy := int64(nguoi.nv.ViTriY) - int64(mon.y)
	if dx*dx+dy*dy > tamNhatBinhPhuong {
		return // ngoai tam -> im lang bo qua (goi tu tay)
	}

	// CHOT LOP 1: go khoi mat dat NGAY. Tu giay phut nay khong ai khac xin
	// nhat duoc mon nay nua.
	ban := *mon
	delete(k.doDuoiDat, sk.itemMapId)

	bo := k.boNhat.Load()
	if bo == nil || !bo.Xep(yeuCauNhat{khu: k, nvID: sk.idNhanVat, itemMapId: sk.itemMapId, mon: ban}) {
		// Chua noi bo nhat, hoac hang doi day. TRA MON LAI ngay — khong duoc
		// de mon boc hoi chi vi may chu ban.
		k.doDuoiDat[sk.itemMapId] = &ban
		return
	}
}

// xuLyNhatXong la than cua case suKienNhatXong. CHAY TRONG ACTOR.
func (k *KhuVuc) xuLyNhatXong(sk suKienNhatXong) {
	if sk.loi != nil {
		// Khong co gi doi trong CSDL -> tra mon lai mat dat. KHONG phat 68 lan
		// nua: client cua moi nguoi VAN dang thay mon do (chua he nhan -19/-21).
		if _, dungRoi := k.doDuoiDat[sk.itemMapId]; !dungRoi {
			ban := sk.mon
			k.doDuoiDat[sk.itemMapId] = &ban
		}
		if nguoi, co := k.nguoiTrong[sk.idNhanVat]; co {
			loiNhan := "Không nhặt được"
			if errors.Is(sk.loi, store.ErrDaCoNguoiNhat) {
				loiNhan = "Món đã có người nhặt"
			} else if errors.Is(sk.loi, store.ErrOCoMon) || errors.Is(sk.loi, errTuiDay) {
				loiNhan = "Túi đã đầy"
			}
			if m := goiNhanTinRieng(loiNhan); m != nil {
				if err := nguoi.phien.XepGui(m); err != nil {
					k.loaiNguoiLoiGui(sk.idNhanVat)
				}
			}
		}
		return
	}

	// Thanh cong: mon da doi chu trong CSDL roi moi den luot goi tin (§2.3 b7).
	loiNhan := "Bạn nhận được " + sk.mon.ten
	if doVuongGhiHaiShort(sk.mon.loai) {
		loiNhan = "" // tien: de rong de client hien chu bay +N
	}
	m, err := goiNhatDuoc(sk.itemMapId, sk.mon.loai, loiNhan, sk.mon.soLuong)
	if err != nil {
		slog.Error("dung goi -20 that bai", "itemMapId", sk.itemMapId, "loi", err)
	} else if nguoi, co := k.nguoiTrong[sk.idNhanVat]; co {
		if err := nguoi.phien.XepGui(m); err != nil {
			k.loaiNguoiLoiGui(sk.idNhanVat)
		}
	}
	// -19 phat NGAY trong cung nhip (xem chu thich goiNguoiKhacNhat).
	k.phatChoNguoiKhac(sk.idNhanVat, goiNguoiKhacNhat(sk.itemMapId, int32(sk.idNhanVat)))
}

// ---------------------------------------------------------------------------
// 5. Bo nhat do — goroutine DUY NHAT duoc cham CSDL trong ca luong nay
// ---------------------------------------------------------------------------

// errTuiDay: khong con o trong nao trong tui.
var errTuiDay = errors.New("tui da day")

// KhoNhat la phan toi thieu cua *store.Store ma bo nhat do can — interface HEP
// dung khuon voi Kho/GhiKhiRoi (login.go), de test duoc bang kho gia dong bo,
// khong can CSDL that.
type KhoNhat interface {
	LayTui(ctx context.Context, nvID int64, noiChua string) ([]store.VatPham, error)
	NhatLen(ctx context.Context, nvID int64, vatPhamID int64, oSo int16) error
}

// yeuCauNhat la MOT yeu cau ghi CSDL do actor xep ra.
type yeuCauNhat struct {
	khu       *KhuVuc
	nvID      int64
	itemMapId int16
	mon       doVuong
}

// BoNhatDo la goroutine ghi CSDL cho duong nhat do. Mot bo cho ca may chu la
// du (moi yeu cau la mot giao dich ngan); muon nhieu hon thi chay nhieu Chay().
type BoNhatDo struct {
	kho  KhoNhat
	hang chan yeuCauNhat
	dung chan struct{}
	// hetGio la han cho MOT yeu cau CSDL. Vuot han -> tra mon lai mat dat con
	// hon de yeu cau treo mai.
	hetGio time.Duration
}

// MoBoNhatDo tao bo nhat do voi hang doi suc chua `sucChua`. Chua chay gi —
// goi Chay() (thuong la `go bo.Chay()`) de bat dau.
func MoBoNhatDo(kho KhoNhat, sucChua int) *BoNhatDo {
	if sucChua <= 0 {
		sucChua = 256
	}
	return &BoNhatDo{
		kho:    kho,
		hang:   make(chan yeuCauNhat, sucChua),
		dung:   make(chan struct{}),
		hetGio: 5 * time.Second,
	}
}

// Xep dua mot yeu cau vao hang. KHONG BAO GIO CHAN — hang day thi tra false va
// actor tra mon lai mat dat. Actor cho o day mot phan nghin giay cung la treo
// ca ban do (§6.10b).
func (b *BoNhatDo) Xep(yc yeuCauNhat) bool {
	if b == nil {
		return false
	}
	select {
	case b.hang <- yc:
		return true
	default:
		return false
	}
}

// Chay la vong lap cua goroutine ghi CSDL. Tra ve khi Dung() duoc goi.
func (b *BoNhatDo) Chay() {
	for {
		select {
		case <-b.dung:
			return
		case yc := <-b.hang:
			b.lam(yc)
		}
	}
}

// Dung dung goroutine. An toan goi nhieu lan? KHONG — goi dung mot lan (chu
// so huu la main). Giu don gian: khong ai trong repo goi Dung hai lan.
func (b *BoNhatDo) Dung() { close(b.dung) }

// lam chay MOT yeu cau: tim o trong -> NhatLen -> bao nguoc ve actor.
//
// Loi o BAT CU buoc nao deu di nguoc ve actor duoi dang suKienNhatXong{loi:...}
// — actor tra mon lai mat dat. Khong co duong nao mon bien mat lang le.
func (b *BoNhatDo) lam(yc yeuCauNhat) {
	ctx, huy := context.WithTimeout(context.Background(), b.hetGio)
	defer huy()

	loi := b.ghi(ctx, yc)
	yc.khu.Gui(suKienNhatXong{
		idNhanVat: yc.nvID,
		itemMapId: yc.itemMapId,
		mon:       yc.mon,
		loi:       loi,
	})
}

func (b *BoNhatDo) ghi(ctx context.Context, yc yeuCauNhat) error {
	dsMon, err := b.kho.LayTui(ctx, yc.nvID, store.NoiChuaTui)
	if err != nil {
		return err
	}
	oSo, co := timOTrong(dsMon, doVuongSoOTui)
	if !co {
		return errTuiDay
	}
	return b.kho.NhatLen(ctx, yc.nvID, yc.mon.vatPhamID, oSo)
}

// timOTrong tra o trong DAU TIEN trong [0, soO). Ham thuan, de test rieng.
//
// Mon co o_so nil (dang duoi dat / da huy lot vao lat cat) duoc bo qua — chung
// khong chiem o nao.
func timOTrong(dsMon []store.VatPham, soO int) (int16, bool) {
	daDung := make(map[int16]bool, len(dsMon))
	for i := range dsMon {
		if dsMon[i].OSo != nil {
			daDung[*dsMon[i].OSo] = true
		}
	}
	for o := 0; o < soO; o++ {
		if !daDung[int16(o)] {
			return int16(o), true
		}
	}
	return 0, false
}

// ---------------------------------------------------------------------------
// 6. Handler goi LEN `-20` — chay tren goroutine PHIEN
// ---------------------------------------------------------------------------

// XuLyNhat xu ly goi LEN `-20` (than goi: `short itemMapId`, Service.cs:1276-1291).
//
// Dieu phoi noi MOT DONG vao login.go sau khi gop:
//
//	case lenhNhatDo: // -20
//		return b.XuLyNhat(m)
//
// Ham nay CHI doc 2 byte roi gui su kien vao actor — dung khuon voi
// xuLyDanhNguoiChoi (login.go): khong tra doDuoiDat, khong cham CSDL, khong
// tinh gi ca. Goi ngan hon khuon -> bo qua an toan, KHONG ngat ket noi (client
// gui goi rac la chuyen binh thuong).
func (b *BoXuLy) XuLyNhat(m *protocol.Message) error {
	if b.khuVuc == nil || b.nhanVat == nil {
		return nil
	}
	r := protocol.NewReader(m.Data)
	itemMapId, err := r.Short()
	if err != nil {
		return nil
	}
	b.khuVuc.Gui(suKienXinNhat{idNhanVat: b.nhanVat.ID, itemMapId: itemMapId})
	return nil
}

// ThaDoVaoKhuVuc la duong VAO tu ngoai actor: dat mot mon (da co dong CSDL o
// DUOI_DAT) xuong khu vuc. Task 12 goi ham nay sau khi store.SinhDoDuoiDat tra
// ve id. KHONG BAO GIO chan: hop thu day thi mat mon roi (mot mon roi mat di
// chap nhan duoc; mot actor bi treo thi khong).
func (k *KhuVuc) ThaDoVaoKhuVuc(vatPhamID int64, mauID, soLuong int32, x, y int16,
	loai int8, ten string) {
	k.Gui(suKienThaDo{mon: doVuong{
		vatPhamID: vatPhamID, mauID: mauID, soLuong: soLuong,
		x: x, y: y, loai: loai, ten: ten,
	}})
}
