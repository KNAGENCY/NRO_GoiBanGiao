package game

// dovuong_test.go — KH7 Task 11.
//
// Hai nhom test:
//
//	A. BYTE-CHINH-XAC bon goi 68 / -20 / -19 / -21, co ca ca type 9/10/34 (loi L2).
//	B. ACTOR: hai nguoi cung xin nhat MOT mon -> chi MOT yeu cau CSDL duoc xep;
//	   that bai CSDL -> mon quay lai mat dat.
//
// Cac test dot bien duoc danh dau ⚠️ ngay tren ham.

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// thanUTF tra ve DUNG chuoi byte ma Writer.UTF sinh ra (2 byte do dai + modified
// UTF-8). Dung chinh Writer de test khong phai chep tay bang ma hoa — cai dang
// kiem o day la CAU TRUC goi, khong phai bo ma hoa UTF (da co test rieng).
func thanUTF(s string) []byte {
	w := protocol.NewWriter()
	if err := w.UTF(s); err != nil {
		panic(err)
	}
	return w.Bytes()
}

// ---------------------------------------------------------------------------
// A. Byte-chinh-xac
// ---------------------------------------------------------------------------

func kiemGoi(t *testing.T, m *protocol.Message, cmd int8, than []byte) {
	t.Helper()
	if m == nil {
		t.Fatalf("goi nil")
	}
	if m.Command != cmd {
		t.Fatalf("ma lenh = %d, mong doi %d", m.Command, cmd)
	}
	if !bytes.Equal(m.Data, than) {
		t.Fatalf("than goi = % x, mong doi % x", m.Data, than)
	}
}

func TestGoiDoXuatHienDungTungByte(t *testing.T) {
	m := goiDoXuatHien(7, 3, 100, 384)
	kiemGoi(t, m, 68, []byte{
		0, 7, // itemMapId
		0, 3, // mauId
		0, 100, // x
		1, 128, // y = 384
		255, 255, 255, 255, // chuUuTien = -1 (ai cung nhat duoc)
	})
}

// ⚠️ CHOT CHAN DOT BIEN L2 (mot nua).
// Loai THUONG: dung MOT short so luong. Ghi hai short o day la client doc du
// mot short vao phan goi sau -> lech ca luong tin.
func TestGoiNhatDuocLoaiThuongGhiMotShort(t *testing.T) {
	m, err := goiNhatDuoc(9, 6, "Bạn nhận được X", 3)
	if err != nil {
		t.Fatalf("dung goi that bai: %v", err)
	}
	// "Bạn nhận được X" = 18 byte modified UTF-8 (5 chu co dau x 2 byte + 10 x 1).
	than := []byte{0, 9}
	than = append(than, thanUTF("Bạn nhận được X")...)
	than = append(than, 0, 3) // MOT short so luong
	kiemGoi(t, m, -20, than)
}

// ⚠️ CHOT CHAN DOT BIEN L2 (nua con lai) — day la lo hong that cua ban Java.
// Type 9 (vang) / 10 (ngoc) / 34 (luong khoa): client doc HAI short va CHI cong
// tien bang cai THU HAI. Ghi mot short -> tien khong cong, im lang tuyet doi.
func TestGoiNhatDuocLoai9va10va34GhiHaiShort(t *testing.T) {
	for _, loai := range []int8{9, 10, 34} {
		m, err := goiNhatDuoc(4, loai, "", 250)
		if err != nil {
			t.Fatalf("loai %d: dung goi that bai: %v", loai, err)
		}
		than := []byte{0, 4, 0, 0, 0, 250, 0, 250} // id, UTF rong (len 0), short, short
		kiemGoi(t, m, -20, than)
	}
}

// Loai 8 va 11 (ngay canh 9/10) PHAI chi ghi mot short — chan dot bien "mo rong
// dieu kien thanh loai >= 9".
func TestGoiNhatDuocLoaiCanh9va10ChiMotShort(t *testing.T) {
	for _, loai := range []int8{8, 11, 33, 35} {
		m, err := goiNhatDuoc(4, loai, "", 250)
		if err != nil {
			t.Fatalf("loai %d: %v", loai, err)
		}
		if len(m.Data) != 6 {
			t.Fatalf("loai %d: than dai %d byte, mong doi 6 (mot short)", loai, len(m.Data))
		}
	}
}

// L4 khong chep: so luong tran short -> tra LOI, KHONG gui goi mang so sai.
func TestGoiNhatDuocSoLuongTranThiLoi(t *testing.T) {
	if _, err := goiNhatDuoc(1, 6, "", 32768); err == nil {
		t.Fatalf("so luong 32768 phai loi")
	}
	if _, err := goiNhatDuoc(1, 6, "", -1); err == nil {
		t.Fatalf("so luong am phai loi")
	}
	if _, err := goiNhatDuoc(1, 6, "", 32767); err != nil {
		t.Fatalf("so luong 32767 phai hop le: %v", err)
	}
}

func TestGoiNguoiKhacNhatDungTungByte(t *testing.T) {
	kiemGoi(t, goiNguoiKhacNhat(258, 5), -19, []byte{1, 2, 0, 0, 0, 5})
}

func TestGoiDoBienMatDungTungByte(t *testing.T) {
	kiemGoi(t, goiDoBienMat(258), -21, []byte{1, 2})
}

func TestGoiNhanTinRiengDungTungByte(t *testing.T) {
	m := goiNhanTinRieng("abc")
	kiemGoi(t, m, -25, []byte{0, 3, 'a', 'b', 'c'})
}

// ---------------------------------------------------------------------------
// B. Actor
// ---------------------------------------------------------------------------

// khoNhatGia dem so lan NhatLen duoc goi va tra loi da dat truoc.
type khoNhatGia struct {
	tui     []store.VatPham
	loiNhat error
	soLan   int

	// Duong DON DO HET GIO (bo sung KH7).
	loiHuy   error
	soLanHuy int
	daHuy    []int64
}

func (k *khoNhatGia) LayTui(context.Context, int64, string) ([]store.VatPham, error) {
	return k.tui, nil
}

func (k *khoNhatGia) NhatLen(context.Context, int64, int64, int16) error {
	k.soLan++
	return k.loiNhat
}

func (k *khoNhatGia) HuyDoDuoiDatHetGio(_ context.Context, vatPhamID int64) error {
	k.soLanHuy++
	k.daHuy = append(k.daHuy, vatPhamID)
	return k.loiHuy
}

// khuVucDoGia dung mot khu vuc CHUA BatDau() voi hai nguoi dung yen tai (0,0)
// va MOT mon nam duoi dat cung cho. Tra ve itemMapId cua mon do.
func khuVucDoGia(t *testing.T) (k *KhuVuc, pA, pB *phienGiaKhuVuc, id int16) {
	t.Helper()
	k = MoKhuVuc(TimBanDo(1))
	nvA := &store.NhanVat{ID: 1, Ten: "A", HpHienTai: 100, HpToiDa: 100}
	nvB := &store.NhanVat{ID: 2, Ten: "B", HpHienTai: 100, HpToiDa: 100}
	pA, pB = &phienGiaKhuVuc{}, &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvA, phien: pA})
	k.xuLy(suKienVao{nv: nvB, phien: pB})

	tra := make(chan int16, 1)
	k.xuLy(suKienThaDo{mon: doVuong{
		vatPhamID: 555, mauID: 3, soLuong: 1, x: 0, y: 0, loai: 6, ten: "Kunai",
	}, traLoi: tra})
	id = <-tra
	pA.nhan, pB.nhan = nil, nil
	return k, pA, pB, id
}

// Tha do xuong khu vuc thi CA KHU VUC nhan goi 68 va mon nam trong k.doDuoiDat.
func TestThaDoPhatGoi68ChoCaKhuVuc(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	p := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: &store.NhanVat{ID: 1, Ten: "A", HpHienTai: 1}, phien: p})
	p.nhan = nil

	tra := make(chan int16, 1)
	k.xuLy(suKienThaDo{mon: doVuong{vatPhamID: 1, mauID: 3, soLuong: 1, x: 10, y: 20}, traLoi: tra})
	id := <-tra

	if _, co := k.doDuoiDat[id]; !co {
		t.Fatalf("mon phai nam trong doDuoiDat")
	}
	tin := p.banTin()
	if len(tin) != 1 || tin[0].Command != 68 {
		t.Fatalf("phai nhan dung mot goi 68, nhan %+v", tin)
	}
}

// ⚠️ CHOT CHAN DOT BIEN CHINH CUA ACTOR.
//
// Hai nguoi xin nhat CUNG mot mon -> chi MOT yeu cau duoc xep ra bo nhat do.
// Neu actor khong go mon khoi k.doDuoiDat TRUOC khi xep yeu cau (hoac go sau),
// ca hai su kien deu thay mon con day va deu xep -> so yeu cau = 2 -> test do.
func TestHaiNguoiCungXinNhatChiMotYeuCauDuocXep(t *testing.T) {
	k, _, _, id := khuVucDoGia(t)
	kho := &khoNhatGia{}
	bo := MoBoNhatDo(kho, 8)
	k.GanBoNhat(bo)

	k.xuLy(suKienXinNhat{idNhanVat: 1, itemMapId: id})
	k.xuLy(suKienXinNhat{idNhanVat: 2, itemMapId: id})

	if n := len(bo.hang); n != 1 {
		t.Fatalf("phai co dung 1 yeu cau trong hang, co %d", n)
	}
	if _, con := k.doDuoiDat[id]; con {
		t.Fatalf("mon phai da bi go khoi doDuoiDat ngay khi nhan yeu cau dau")
	}
}

// Nhat thanh cong: nguoi nhat nhan -20, NGUOI KHAC nhan -19 (khong phai nguoc lai).
func TestNhatXongThanhCongPhat20ChoNguoiNhatVa19ChoNguoiKhac(t *testing.T) {
	k, pA, pB, id := khuVucDoGia(t)
	delete(k.doDuoiDat, id) // actor da go mon ra tu buoc suKienXinNhat

	k.xuLy(suKienNhatXong{
		idNhanVat: 1, itemMapId: id,
		mon: doVuong{vatPhamID: 555, mauID: 3, soLuong: 1, loai: 6, ten: "Kunai"},
	})

	tinA, tinB := pA.banTin(), pB.banTin()
	if len(tinA) != 1 || tinA[0].Command != -20 {
		t.Fatalf("nguoi nhat phai nhan dung mot goi -20, nhan %+v", tinA)
	}
	if len(tinB) != 1 || tinB[0].Command != -19 {
		t.Fatalf("nguoi khac phai nhan dung mot goi -19, nhan %+v", tinB)
	}
	if _, con := k.doDuoiDat[id]; con {
		t.Fatalf("nhat thanh cong roi thi mon khong duoc quay lai mat dat")
	}
}

// ⚠️ CHOT CHAN "khong mat do": CSDL bao thua cuoc -> mon QUAY LAI mat dat, va
// nguoi xin nhat nhan mot dong -25 giai thich. Bo nhanh tra mon lai trong
// xuLyNhatXong -> mon boc hoi -> test nay do.
func TestNhatXongThatBaiThiMonQuayLaiMatDat(t *testing.T) {
	k, pA, pB, id := khuVucDoGia(t)
	mon := *k.doDuoiDat[id]
	delete(k.doDuoiDat, id) // gia lap: actor da go mon ra roi

	k.xuLy(suKienNhatXong{
		idNhanVat: 1, itemMapId: id, mon: mon,
		loi: store.ErrDaCoNguoiNhat,
	})

	con, co := k.doDuoiDat[id]
	if !co {
		t.Fatalf("that bai xong mon PHAI quay lai mat dat")
	}
	if con.vatPhamID != mon.vatPhamID {
		t.Fatalf("mon quay lai sai: %+v", con)
	}
	tinA := pA.banTin()
	if len(tinA) != 1 || tinA[0].Command != -25 {
		t.Fatalf("nguoi xin nhat phai nhan dung mot goi -25, nhan %+v", tinA)
	}
	if n := len(pB.banTin()); n != 0 {
		t.Fatalf("nguoi khac khong duoc nhan gi khi nhat that bai, nhan %d goi", n)
	}
}

// Chua noi bo nhat do (boNhat == nil) -> mon PHAI o nguyen duoi dat, khong panic.
func TestXinNhatKhiChuaCoBoNhatThiMonONguyen(t *testing.T) {
	k, _, _, id := khuVucDoGia(t)
	k.xuLy(suKienXinNhat{idNhanVat: 1, itemMapId: id})
	if _, co := k.doDuoiDat[id]; !co {
		t.Fatalf("chua co bo nhat thi mon phai o nguyen duoi dat")
	}
}

// Ngoai tam -> khong xep yeu cau nao, mon o nguyen.
func TestXinNhatNgoaiTamThiBoQua(t *testing.T) {
	k, _, _, id := khuVucDoGia(t)
	bo := MoBoNhatDo(&khoNhatGia{}, 8)
	k.GanBoNhat(bo)
	k.nguoiTrong[1].nv.ViTriX = 5000 // rat xa mon o (0,0)

	k.xuLy(suKienXinNhat{idNhanVat: 1, itemMapId: id})

	if len(bo.hang) != 0 {
		t.Fatalf("ngoai tam khong duoc xep yeu cau nao")
	}
	if _, co := k.doDuoiDat[id]; !co {
		t.Fatalf("ngoai tam thi mon phai o nguyen")
	}
}

// Nguoi dang chet khong nhat duoc (giong Controller.java:724).
func TestXinNhatKhiDangChetThiBoQua(t *testing.T) {
	k, _, _, id := khuVucDoGia(t)
	bo := MoBoNhatDo(&khoNhatGia{}, 8)
	k.GanBoNhat(bo)
	k.nguoiTrong[1].nv.HpHienTai = 0

	k.xuLy(suKienXinNhat{idNhanVat: 1, itemMapId: id})
	if len(bo.hang) != 0 || len(k.doDuoiDat) != 1 {
		t.Fatalf("nguoi chet khong duoc nhat")
	}
}

// Don do het gio -> phat -21 cho ca khu vuc va go khoi map. Goi lan hai khong
// phat gi (mon da di roi).
func TestDonDoHetGioPhat21MotLan(t *testing.T) {
	k, pA, _, id := khuVucDoGia(t)

	k.xuLy(suKienDonDoHetGio{itemMapId: id})
	tin := pA.banTin()
	if len(tin) != 1 || tin[0].Command != -21 {
		t.Fatalf("phai nhan dung mot goi -21, nhan %+v", tin)
	}

	pA.nhan = nil
	k.xuLy(suKienDonDoHetGio{itemMapId: id})
	if n := len(pA.banTin()); n != 0 {
		t.Fatalf("don lan hai khong duoc phat gi, phat %d goi", n)
	}
}

// capItemMapId khong bao gio tra id da dung va khong bao gio tra so am.
func TestCapItemMapIdKhongTrungKhongAm(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.itemMapIdKe = 32766
	daThay := map[int16]bool{}
	for i := 0; i < 5; i++ {
		id := k.capItemMapId()
		if id < 0 {
			t.Fatalf("id am: %d", id)
		}
		if daThay[id] {
			t.Fatalf("id trung: %d", id)
		}
		daThay[id] = true
		k.doDuoiDat[id] = &doVuong{}
	}
}

// ---------------------------------------------------------------------------
// C. Bo nhat do (goroutine CSDL) — ham thuan + duong ghi
// ---------------------------------------------------------------------------

func TestTimOTrong(t *testing.T) {
	o0, o1 := int16(0), int16(1)
	ds := []store.VatPham{{OSo: &o0}, {OSo: &o1}, {OSo: nil}}
	o, co := timOTrong(ds, 30)
	if !co || o != 2 {
		t.Fatalf("o trong dau tien phai la 2, duoc %d (co=%v)", o, co)
	}
	if _, co := timOTrong(ds, 2); co {
		t.Fatalf("tui 2 o da day thi phai bao khong con o")
	}
}

// Tui day -> bo nhat tra errTuiDay, KHONG goi NhatLen (khong duoc dong vao CSDL
// khi biet chac se that bai).
func TestBoNhatTuiDayThiKhongGoiNhatLen(t *testing.T) {
	ds := make([]store.VatPham, doVuongSoOTui)
	for i := range ds {
		o := int16(i)
		ds[i].OSo = &o
	}
	kho := &khoNhatGia{tui: ds}
	bo := MoBoNhatDo(kho, 4)

	err := bo.ghi(context.Background(), yeuCauNhat{nvID: 1, mon: doVuong{vatPhamID: 9}})
	if !errors.Is(err, errTuiDay) {
		t.Fatalf("phai errTuiDay, nhan %v", err)
	}
	if kho.soLan != 0 {
		t.Fatalf("tui day thi khong duoc goi NhatLen, goi %d lan", kho.soLan)
	}
}

// Hang doi day -> Xep tra false (khong bao gio chan actor).
func TestXepKhiHangDayTraFalse(t *testing.T) {
	bo := MoBoNhatDo(&khoNhatGia{}, 1)
	if !bo.Xep(yeuCauNhat{}) {
		t.Fatalf("yeu cau dau phai vao duoc")
	}
	if bo.Xep(yeuCauNhat{}) {
		t.Fatalf("hang day thi Xep phai tra false, khong duoc chan")
	}
}

// ---------------------------------------------------------------------------
// C. HAN GIO TU DON do roi duoi dat (bo sung KH7)
// ---------------------------------------------------------------------------

// mocDonDo la moc thoi gian co dinh cho ca nhom test nay. Gia tri cu the khong
// quan trong; cai quan trong la MOI moc deu do dong ho TIEM (khuVucHetGioGia)
// sinh ra, khong mot cho nao goi time.Now().
var mocDonDo = time.Date(2026, 9, 3, 12, 0, 0, 0, time.UTC)

// khuVucHetGioGia dung khu vuc CHUA BatDau() voi dong ho tiem duoc, mot nguoi
// choi (de bat goi -21) va mot bo nhat do gia. Tra ve con tro `luc` de test tu
// BOM dong ho: `*luc = luc.Add(...)`.
func khuVucHetGioGia(t *testing.T) (k *KhuVuc, p *phienGiaKhuVuc, kho *khoNhatGia,
	bo *BoNhatDo, luc *time.Time) {
	t.Helper()
	moc := mocDonDo
	luc = &moc

	k = MoKhuVuc(TimBanDo(1))
	k.dongHo = func() time.Time { return *luc }
	p = &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: &store.NhanVat{ID: 1, Ten: "A", HpHienTai: 100, HpToiDa: 100}, phien: p})

	kho = &khoNhatGia{}
	bo = MoBoNhatDo(kho, 8)
	k.GanBoNhat(bo)
	return k, p, kho, bo, luc
}

// thaMonGia tha MOT mon xuong khu vuc va tra itemMapId (dong bo, khong BatDau).
func thaMonGia(t *testing.T, k *KhuVuc, vatPhamID int64) int16 {
	t.Helper()
	tra := make(chan int16, 1)
	k.xuLy(suKienThaDo{mon: doVuong{
		vatPhamID: vatPhamID, mauID: 3, soLuong: 1, x: 0, y: 0, loai: 6, ten: "Kunai",
	}, traLoi: tra})
	return <-tra
}

// goi21 loc ra cac goi -21 trong ban tin cua mot phien.
func goi21(p *phienGiaKhuVuc) []*protocol.Message {
	var ra []*protocol.Message
	for _, m := range p.banTin() {
		if m != nil && m.Command == lenhDoBienMat {
			ra = append(ra, m)
		}
	}
	return ra
}

// ⚠️ CHOT BIEN. Chua toi han (thieu 100ms) thi KHONG don. Neu ai doi `>=` o
// nhipDonDoHetGio thanh `>` thi test NAY van xanh — cap doi cua no
// (TestDonDoDungHanThiDon) moi do. Hai test phai di CUNG NHAU.
func TestDonDoChuaToiHanThiKhongDon(t *testing.T) {
	k, p, kho, bo, luc := khuVucHetGioGia(t)
	id := thaMonGia(t, k, 555)

	*luc = luc.Add(hetGioDoRoi - 100*time.Millisecond)
	k.nhipTim()

	if _, con := k.doDuoiDat[id]; !con {
		t.Fatalf("chua toi han ma mon da bi don")
	}
	if g := goi21(p); len(g) != 0 {
		t.Fatalf("chua toi han thi khong duoc phat -21, phat %d goi", len(g))
	}
	if n := len(bo.hang); n != 0 {
		t.Fatalf("chua toi han thi khong duoc xep yeu cau huy, xep %d", n)
	}
	if kho.soLanHuy != 0 {
		t.Fatalf("actor khong duoc cham CSDL, goi %d lan", kho.soLanHuy)
	}
}

// ⚠️ CHOT BIEN + CHOT DONG HO TIEM. DUNG han thi don: mon roi khoi map, phat
// DUNG MOT goi -21 dung itemMapId, va xep DUNG MOT yeu cau huy CSDL.
//
// Test nay do neu:
//   - don ma khong phat -21 (dot bien 1);
//   - so `>` thay vi `>=` (dot bien 2 — dung han la bien);
//   - dung time.Now() thay k.dongHo() (dot bien 3 — dong ho that moi nhich vai
//     micro giay, mon khong bao gio toi han, mon van con trong map).
func TestDonDoDungHanThiDon(t *testing.T) {
	k, p, kho, bo, luc := khuVucHetGioGia(t)
	id := thaMonGia(t, k, 555)

	*luc = luc.Add(hetGioDoRoi)
	k.nhipTim()

	if _, con := k.doDuoiDat[id]; con {
		t.Fatalf("dung han ma mon van con duoi dat")
	}
	g := goi21(p)
	if len(g) != 1 {
		t.Fatalf("phai phat dung MOT goi -21, phat %d", len(g))
	}
	than := protocol.NewWriter()
	than.Short(id)
	kiemGoi(t, g[0], lenhDoBienMat, than.Bytes())

	if n := len(bo.hang); n != 1 {
		t.Fatalf("phai xep dung MOT yeu cau huy, xep %d", n)
	}
	yc := <-bo.hang
	if !yc.huyHetGio {
		t.Fatalf("yeu cau phai duoc danh dau huyHetGio")
	}
	if yc.mon.vatPhamID != 555 {
		t.Fatalf("yeu cau huy phai mang vatPhamID 555, mang %d", yc.mon.vatPhamID)
	}
	if kho.soLanHuy != 0 {
		t.Fatalf("actor KHONG duoc tu goi CSDL (§6.10b), goi %d lan", kho.soLanHuy)
	}
}

// Mon DA BI NHAT truoc do (khong con trong map) thi nhip don khong dung toi:
// khong -21 lan hai, khong yeu cau huy. Neu khong the, nguoi vua nhat duoc mon
// se bi may chu danh dau DA_HUY chinh mon dang nam trong tui ho.
func TestDonDoKhongDungMonDaNhat(t *testing.T) {
	k, p, _, bo, luc := khuVucHetGioGia(t)
	id := thaMonGia(t, k, 555)

	// Nguoi choi nhat mon (actor go khoi map ngay).
	k.xuLy(suKienXinNhat{idNhanVat: 1, itemMapId: id})
	if _, con := k.doDuoiDat[id]; con {
		t.Fatalf("tien de sai: mon phai da bi go khi xin nhat")
	}
	<-bo.hang // rut yeu cau NHAT ra de dem cho sach

	*luc = luc.Add(hetGioDoRoi * 10)
	k.nhipTim()

	if g := goi21(p); len(g) != 0 {
		t.Fatalf("mon da bi nhat, khong duoc phat -21, phat %d", len(g))
	}
	if n := len(bo.hang); n != 0 {
		t.Fatalf("mon da bi nhat, khong duoc xep yeu cau huy, xep %d", n)
	}
}

// Hai mon tha o hai moc khac nhau: nhip chi don DUNG cai da toi han.
func TestDonDoChiDonMonToiHan(t *testing.T) {
	k, p, _, bo, luc := khuVucHetGioGia(t)
	idCu := thaMonGia(t, k, 111)

	*luc = luc.Add(30 * time.Second)
	idMoi := thaMonGia(t, k, 222)

	// Toi han cua mon CU (30s + 20s = 50s), mon MOI moi duoc 20s.
	*luc = luc.Add(hetGioDoRoi - 30*time.Second)
	k.nhipTim()

	if _, con := k.doDuoiDat[idCu]; con {
		t.Fatalf("mon cu da toi han, phai bi don")
	}
	if _, con := k.doDuoiDat[idMoi]; !con {
		t.Fatalf("mon moi chua toi han, khong duoc don")
	}
	g := goi21(p)
	if len(g) != 1 {
		t.Fatalf("phai phat dung MOT goi -21, phat %d", len(g))
	}
	than := protocol.NewWriter()
	than.Short(idCu)
	kiemGoi(t, g[0], lenhDoBienMat, than.Bytes())

	if n := len(bo.hang); n != 1 {
		t.Fatalf("phai xep dung MOT yeu cau huy, xep %d", n)
	}
	if yc := <-bo.hang; yc.mon.vatPhamID != 111 {
		t.Fatalf("phai huy mon cu (111), huy %d", yc.mon.vatPhamID)
	}
}

// Bo nhat do xu ly yeu cau huyHetGio: goi DUNG HuyDoDuoiDatHetGio voi dung id,
// KHONG goi NhatLen, va KHONG bao nguoc ve actor (mon da bien mat khoi game).
func TestBoNhatDoXuLyYeuCauHuyHetGio(t *testing.T) {
	k, _, kho, bo, _ := khuVucHetGioGia(t)

	bo.lam(yeuCauNhat{khu: k, itemMapId: 7, huyHetGio: true,
		mon: doVuong{vatPhamID: 999}})

	if kho.soLanHuy != 1 || len(kho.daHuy) != 1 || kho.daHuy[0] != 999 {
		t.Fatalf("phai goi HuyDoDuoiDatHetGio dung mot lan voi id 999, nhan %d lan %v",
			kho.soLanHuy, kho.daHuy)
	}
	if kho.soLan != 0 {
		t.Fatalf("duong huy khong duoc goi NhatLen, goi %d lan", kho.soLan)
	}
	if n := len(k.hopThu); n != 0 {
		t.Fatalf("duong huy KHONG duoc bao nguoc ve actor, hop thu co %d su kien", n)
	}
}

// Huy CSDL that bai thi chi log — mon KHONG quay lai mat dat (khac han duong
// nhat). Neu ai do "sua cho dung" bang cach tra mon lai, test nay do.
func TestHuyHetGioThatBaiThiMonKhongQuayLai(t *testing.T) {
	k, _, kho, bo, luc := khuVucHetGioGia(t)
	kho.loiHuy = errors.New("CSDL sap")
	id := thaMonGia(t, k, 555)

	*luc = luc.Add(hetGioDoRoi)
	k.nhipTim()
	bo.lam(<-bo.hang)

	if _, con := k.doDuoiDat[id]; con {
		t.Fatalf("huy CSDL loi thi mon van phai bien mat khoi game")
	}
	if n := len(k.hopThu); n != 0 {
		t.Fatalf("loi huy khong duoc sinh su kien nao cho actor, co %d", n)
	}
}
