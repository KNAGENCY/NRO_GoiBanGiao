// boss.go — KIEU DU LIEU + HAM THUAN cho he boss (KH6b, thiet ke
// docs/superpowers/specs/2026-08-03-kh6b-boss-design.md §3-§5).
//
// PHAM VI FILE NAY (co y hep): CHI mo hinh du lieu va cac ham THUAN (khong
// dong ho he thong, khong socket, khong goroutine, khong doc/ghi trang thai
// ngoai tham so). CHUA cam gi vao actor KhuVuc — viec do la task sau, va khi
// cam thi PHAI theo luat §6.9/§6.10b: boss cap nhat NGAY TRONG nhipTim 100ms
// cua khu vuc (KHONG luong 150ms rieng nhu ban Java, se tranh chap viTri), moi
// goi phat ra dung phien.XepGui(), va moi moc thoi gian lay tu k.dongHo() —
// vi vay MOI ham trong file nay nhan bayGio/thoi gian qua THAM SO neu can,
// TUYET DOI khong goi time.Now().
//
// Boss KHONG phai quai (quai.go). Voi client, boss la MOT NHAN VAT — no dung
// lai het goi nguoi choi (-5 xuat hien, -6 roi, 44 chat, -60/-45 danh nhau),
// chi them goi 93 (chay chu toan server). Vi the phan kho nam o vong doi + du
// lieu, khong nam o byte.
package game

import (
	"math/bits"
	"sync/atomic"
	"time"
)

// ---------------------------------------------------------------------------
// 1. Trang thai vong doi boss (thiet ke §5)
// ---------------------------------------------------------------------------

// TrangThaiBoss la mot buoc trong may trang thai vong doi boss (thiet ke §5):
//
//	NGHI -> HIEN -> CHAO -> HOAT_DONG -> (chet) -> TROI -> RA_DI
//	  ^                                                       |
//	  +--- het pha ---------------------------------------------+
//	       (con pha -> quay ve HIEN voi pha+1, GIU NGUYEN Id)
//
// Kieu nen la byte (khong phai int) vi day la truong trang thai nho, nam trong
// struct Boss se co hang tram ban sao — va no KHONG bao gio di thang xuong day
// dan (client khong biet gi ve may trang thai nay), nen khong bi rang buoc
// khuon giao thuc.
type TrangThaiBoss byte

const (
	// BossNghi: chua co mat trong the gioi, dang dem nguoc GiayNghi cua pha
	// (hoac cho trigger). Day la trang thai KHOI TAO (zero value) — co y:
	// mot Boss{} moi tao chua xuat hien o dau ca, dung truc giac.
	BossNghi TrangThaiBoss = iota

	// BossHien: vua chon map + vi tri, phat -5 cho ca khu vuc (va goi 93 neu
	// LaBossLon).
	BossHien

	// BossChao: phat ThoaiHien qua goi chat 44.
	BossChao

	// BossHoatDong: moi tick chon muc tieu + chieu/di chuyen, thinh thoang
	// ThoaiDanh; theo doi MocNoiGian theo %HP.
	BossHoatDong

	// BossTroi: vua chet, phat ThoaiChet (va 93 neu LaBossLon), chia thuong
	// theo dong gop.
	BossTroi

	// BossRaDi: don dep — con pha thi quay ve BossHien voi pha ke tiep (GIU
	// NGUYEN Id, bay #3 thiet ke §7), het pha thi phat -6 va ve BossNghi.
	BossRaDi
)

// String tra ten trang thai DUNG NGUYEN VAN nhu thiet ke §5 (NGHI, HIEN, CHAO,
// HOAT_DONG, TROI, RA_DI) — de log/test doc thang khong phai tra bang.
//
// Gia tri ngoai bang tra "KHONG_RO" thay vi panic: String() hay bi goi tu
// duong log/thong bao loi, panic o do se bien mot du lieu hong thanh mot su co
// nang hon han.
func (t TrangThaiBoss) String() string {
	switch t {
	case BossNghi:
		return "NGHI"
	case BossHien:
		return "HIEN"
	case BossChao:
		return "CHAO"
	case BossHoatDong:
		return "HOAT_DONG"
	case BossTroi:
		return "TROI"
	case BossRaDi:
		return "RA_DI"
	default:
		return "KHONG_RO"
	}
}

// ---------------------------------------------------------------------------
// 2. Mo hinh du lieu (thiet ke §4 — data-driven, KHONG 95 class con)
// ---------------------------------------------------------------------------

// MocNoiGian la mot moc "noi gian theo %mau" cua mot pha (quyet dinh 2, thiet
// ke §4): khi HP boss tut xuong duoi PhanTramHP thi doi he so dame/def va phat
// mot cau chat.
//
// HeSoDame/HeSoDef de float64 vi day la DU LIEU CAN BANG (chu du an chinh tay,
// vd 1.5 / 0.75), khong phai duong tinh sat thuong loi — duong loi van tinh
// bang so nguyen (xem TinhSatThuong, satthuong.go) de tranh sai so gan bien
// tran int32.
type MocNoiGian struct {
	PhanTramHP int     // nguong %HP kich hoat (vd 50 = con 50% mau)
	HeSoDame   float64 // nhan vao sat thuong boss gay ra
	HeSoDef    float64 // nhan vao sat thuong boss NHAN (nho hon = trau hon)
	Chat       string  // cau boss noi khi cham moc nay (goi 44)

	// HeSoNhip nhan vao Boss.NhipDanhHienTai (nho hon 1 = danh NHANH hon).
	// 0 nghia la KHONG doi nhip.
	//
	// ⚠️ TUYET DOI khong duoc "doi nhip" bang cach sua CapChieu.HoiChieuMs:
	// BenDanh.Chieu la CON TRO tro thang vao catalog chieu TOAN CUC dung chung
	// voi NGUOI CHOI — sua qua do la doi hoi chieu cua moi nguoi trong ca may
	// chu. Nhip cua boss ap RIENG, ngoai KiemDanhHopLe.
	HeSoNhip float64
}

// PhaBoss la MOT PHA cua mot boss (thiet ke §4). Mot boss cu the = []PhaBoss:
// giet pha nay -> hoi sinh pha manh hon, GIU NGUYEN Id.
//
// ⚠️ Cap la truong THEM MOI so voi thiet ke §4 (§4 thieu): goi -5 (xuathien.go,
// truong 3 "level") ghi mot BYTE cap nhan vat. Boss voi client LA nhan vat nen
// bat buoc phai co con so nay trong du lieu; thieu no thi luc cam vao actor se
// phai bia ra tai cho (bia = so ma, dung bay #7 thiet ke §7). Mien hop le
// 1..255 — 0 KHONG hop le (client hien "cap 0" vo nghia); kiem bang CapHopLe.
type PhaBoss struct {
	Ten    string   // ten hien tren dau boss (goi -5 truong 9)
	Lang   int16    // lang Naruto — thay cho "gender" Dragon Ball ban goc
	Cap    byte     // cap hien thi (goi -5 truong 3), hop le 1..255
	Outfit [6]int16 // head, body, leg, flagBag, aura, effFront
	Dame   int64    // sat thuong goc cua pha (dua vao TinhSatThuong)

	// TypePk la truong 5 cua goi -5. PHAI khac 0 — client (Char.cs paintHp)
	// THOAT SOM khong ve thanh mau khi cTypePk == 0 VA cFlag == 0, ma goi -5
	// cua ta luon ghi cFlag = 0. De 0 thi nguoi choi danh boss ma KHONG THAY
	// mau boss con bao nhieu. Dung hang so TypePkBoss, dung ghi so tran.
	TypePk byte

	// NhipDanhMs la khoang cach toi thieu giua hai don cua pha nay (mili giay),
	// nap vao Boss.NhipDanhHienTai khi vao pha. 0 = de he thong hoi chieu cua
	// chieu tu lo (xem MocNoiGian.HeSoNhip).
	NhipDanhMs int32

	HpChon  []int64    // luc sinh chon NGAU NHIEN mot gia tri trong day
	MapVao  []int32    // luc hien chon NGAU NHIEN mot khu vuc trong day
	BoChieu [][2]int32 // bo chieu: {templateId, level} — dung catalog chieu KH6

	ThoaiHien []string // chat luc xuat hien (trang thai CHAO)
	ThoaiDanh []string // chat luc dang danh (trang thai HOAT_DONG)
	ThoaiChet []string // chat luc chet (trang thai TROI)

	GiayNghi   int32        // so giay nam yen o trang thai NGHI truoc khi hien
	MocNoiGian []MocNoiGian // cac moc noi gian theo %HP, xem MocNoiGian

	// LaBossLon bat goi 93 (chay chu toan server) luc boss hien va luc bi ha
	// (quyet dinh 4, thiet ke §2). Mac dinh false — im lang, an toan hon:
	// them nham mot boss lon se lam ngap kenh toan server.
	LaBossLon bool
}

// CapHopLe kiem truong Cap co nam trong mien 1..255 khong (xem chu thich
// truong Cap). byte da chan tran 255 giup, nen thuc chat day la kiem "khac 0".
// nil-safe: con tro nil coi nhu KHONG hop le — mot pha khong ton tai thi
// khong the co cap hop le, va nguoi goi khong phai kiem nil truoc.
func (p *PhaBoss) CapHopLe() bool {
	if p == nil {
		return false
	}
	return p.Cap >= 1
}

// Boss la MOT thuc the boss dang song trong pham vi mot KhuVuc (thiet ke §3).
// No la "nguoi choi do may dieu khien": co mat trong khu vuc nhu mot nhan vat
// nhung KHONG co phien mang.
//
// Cung quy tac so huu voi Quai (quai.go): khi da cam vao actor, CHI luong
// vongLap cua khu vuc duoc doc/ghi cac truong nay.
type Boss struct {
	// Id la id runtime AM va DUY NHAT (bay #2 thiet ke §7). PHAI lay tu
	// CapIdBoss() — TUYET DOI khong dung id-loai: ban goc sinh 5 con cung
	// id=-1, khien findCharInMap(id) phia client tra nham con.
	Id int64

	// PhaHienTai la CHI SO trong DuLieu. Ngoai mien [0, len(DuLieu)) nghia la
	// "khong co pha nao dang chay" — PhaHienTaiData tra nil, khong panic.
	PhaHienTai int
	DuLieu     []PhaBoss

	// Hp/HpToiDa kep int32 [0, 2_147_483_647] (§6.17 + bay #4 thiet ke §7):
	// moi phep cong/tru mau PHAI di qua KepHp, ke ca duong "hap thu mau" sau
	// nay, neu khong se tran thanh so am.
	Hp      int32
	HpToiDa int32

	ViTriX int16
	ViTriY int16

	TrangThai TrangThaiBoss

	// MocThoiGian la moc thoi gian cua BUOC trang thai hien tai (vao trang
	// thai luc nao) — may trang thai lay bayGio tu k.dongHo() roi tru truong
	// nay ra "da o trang thai nay bao lau". KHONG BAO GIO gan bang time.Now()
	// trong logic co the test.
	MocThoiGian time.Time

	// LanDanhCuoi la thoi diem boss DUNG THANH CONG mot chieu lan gan nhat —
	// KHAC HAN MocThoiGian (moc vao TRANG THAI). Day la nguon du lieu cho
	// buoc kiem hoi chieu #3 cua KiemDanhHopLe (chotkiem.go): thieu truong nay
	// thi BenDanh.LanDanhCuoi phai de zero value, bayGio.Sub(zero) ra khoang
	// cuc lon nen hoi chieu LUON QUA -> boss danh MOI TICK 100ms (10 don/giay).
	// Zero value = chua tung danh -> lan danh dau tien duoc phep ngay.
	LanDanhCuoi time.Time

	// LucHienPha la thoi diem boss VAO trang thai HIEN cua pha hien tai. Dung
	// lam ViTriCoMau.LucHoiSinh khi boss lam MUC TIEU (nguoi danh boss) — cho
	// boss khoang bat tu ngan sau moi lan bien hinh, giong nguoi choi/quai vua
	// hoi sinh. KHONG dung MocThoiGian vao viec nay: MocThoiGian doi theo MOI
	// buoc trang thai (CHAO, HOAT_DONG...) nen boss se bat tu lai moi buoc.
	LucHienPha time.Time

	// NhipDanhHienTai la khoang cach toi thieu giua hai don cua boss, tinh
	// bang mili giay. Nap tu PhaBoss.NhipDanhMs khi vao pha, roi MocNoiGian
	// nhan HeSoNhip vao day khi boss noi gian. Ap NGOAI KiemDanhHopLe (xem
	// chu thich MocNoiGian.HeSoNhip).
	NhipDanhHienTai int32

	// SatThuongTheoNguoi cong don sat thuong tung nguoi choi gay ra cho boss,
	// dung de chia thuong theo dong gop khi boss chet (quyet dinh 3, thiet ke
	// §2). Nap qua CongSatThuong; chia qua ChiaThuongTheoDongGop; xoa sach sau
	// khi da chia. map nil van DOC an toan (tra 0, len()==0) nen khong bat
	// buoc khoi tao.
	SatThuongTheoNguoi map[int64]int64

	// MocNoiGianDaQua la SO moc noi gian cua pha hien tai DA kich hoat — dung
	// lam con tro chay tien mot chieu de moi moc chi no DUNG MOT LAN. Dat lai
	// ve 0 moi khi doi pha.
	MocNoiGianDaQua int
}

// ---------------------------------------------------------------------------
// 3. Bo dem id boss (bay #2 thiet ke §7)
// ---------------------------------------------------------------------------

// boDemIdBoss la bo dem TOAN CUC cho id boss. Dung sync/atomic thay vi mutex vi
// day la mot phep tang don thuan; day cung la thu DUY NHAT trong file nay cham
// toi trang thai chia se giua cac luong.
var boDemIdBoss atomic.Int64

// CapIdBoss cap mot id boss AM, DUY NHAT: -1, -2, -3, ... An toan khi goi tu
// nhieu luong cung luc (Add la mot phep doc-sua-ghi nguyen tu, hai luong khong
// bao gio nhan cung mot so).
//
// Vi sao AM: client phan biet "nguoi choi that" (id duong) voi thuc the do may
// dieu khien bang DAU cua id; dinh tuyen goi -60 (nguoi danh boss) cung dua vao
// day. Vi sao DUY NHAT: bay #2 thiet ke §7 — ban goc sinh 5 boss cung id=-1,
// findCharInMap(id) phia client tra nham con dau tien tim thay.
func CapIdBoss() int64 {
	return -boDemIdBoss.Add(1)
}

// ---------------------------------------------------------------------------
// 4. Ham thuan — mau, dong gop, chia thuong
// ---------------------------------------------------------------------------

// Kep tran/san mau — ghim thang con so thay vi import math (cung le voi
// satThuongToiDa o satthuong.go: test doc thang khong phai doi chieu math).
const (
	hpToiThieu int64 = 0
	hpToiDa    int64 = 2_147_483_647 // = math.MaxInt32
)

// Tran int64, dung lam moc bao hoa khi cong don sat thuong (xem CongSatThuong).
const tongSatThuongToiDa int64 = 9_223_372_036_854_775_807 // = math.MaxInt64

// TypePkBoss la gia tri truong 5 (typePk) cua goi -5 danh cho boss.
//
// BANG CHUNG (khong bia so):
//   - 1-MAY-CHU-GAME/src/consts/ConstPlayer.java:21  ->  PK_ALL = 5
//   - 1-MAY-CHU-GAME/src/boss/Boss.java:280          ->  boss doi typePk sang PK_ALL
//   - 1-MAY-CHU-GAME/src/boss/Boss.java:553          ->  boss CHI danh khi typePk == PK_ALL
//
// Vi sao BAT BUOC khac 0: client chi ve thanh mau cho nhan vat khi typePk/flag
// khac 0. De 0 thi nguoi choi danh boss ma khong thay mau boss con bao nhieu.
const TypePkBoss byte = 5

// KepHp kep mot gia tri mau ve mien [0, 2_147_483_647] roi tra ve int32.
//
// Nhan int64 (khong phai int32) LA CO Y: nguoi goi tinh mau bang int64 —
// hp + hoi, hp - satThuong, hp * heSo... — roi moi kep MOT LAN o buoc cuoi.
// Neu tinh thang tren int32 thi phep tinh TRAN TRUOC khi kip kep, cho ra so am
// vo nghia (bay #4 thiet ke §7) — cung ly do voi TinhSatThuong (satthuong.go).
func KepHp(v int64) int32 {
	if v < hpToiThieu {
		v = hpToiThieu
	}
	if v > hpToiDa {
		v = hpToiDa
	}
	return int32(v)
}

// PhaHienTaiData tra con tro toi pha dang chay, hoac nil neu PhaHienTai nam
// NGOAI mien [0, len(DuLieu)) — ke ca khi b la nil.
//
// Vi sao tra CON TRO vao thang phan tu cua lat cat (khong phai ban sao): may
// trang thai can DOC du lieu pha thuong xuyen moi tick, sao chep ca PhaBoss
// (chua nam lat cat + ba mang chuoi) moi lan la lang phi. Nguoi goi CHI DOC.
//
// Vi sao nil-safe thay vi panic: PhaHienTai di chuyen theo may trang thai (RA_DI
// tang pha), nen "vuot qua pha cuoi" la trang thai BINH THUONG cua vong doi,
// khong phai loi lap trinh.
func (b *Boss) PhaHienTaiData() *PhaBoss {
	if b == nil {
		return nil
	}
	if b.PhaHienTai < 0 || b.PhaHienTai >= len(b.DuLieu) {
		return nil
	}
	return &b.DuLieu[b.PhaHienTai]
}

// CongSatThuong cong don sat thuong mot nguoi choi gay ra cho boss (thiet ke
// §3: "cong satThuongTheoNguoi").
//
// Bo qua sat <= 0: don truot/do bi chan/loi tinh toan KHONG duoc ghi ten nguoi
// do vao so dong gop — neu ghi, mot nguoi chi dung xa hut vao boss van co ten
// trong danh sach chia thuong voi 0 diem, va (te hon) sat AM co the TRU dong
// gop that cua nguoi khac neu ai do goi nham.
//
// Chan tran int64: bao hoa o MaxInt64 thay vi cuon vong ve so am. Mot tong dong
// gop AM se lam ChiaThuongTheoDongGop tinh sai toan bo ti le.
func (b *Boss) CongSatThuong(idNguoi int64, sat int64) {
	if b == nil || sat <= 0 {
		return
	}
	if b.SatThuongTheoNguoi == nil {
		b.SatThuongTheoNguoi = make(map[int64]int64)
	}
	cu := b.SatThuongTheoNguoi[idNguoi]
	// cu + sat > MaxInt64 <=> cu > MaxInt64 - sat. Viet theo kieu TRU (khong
	// phai cong roi so sanh) vi phep cong da tran roi thi so sanh vo nghia.
	if cu > tongSatThuongToiDa-sat {
		b.SatThuongTheoNguoi[idNguoi] = tongSatThuongToiDa
		return
	}
	b.SatThuongTheoNguoi[idNguoi] = cu + sat
}

// ChiaThuongTheoDongGop chia tongDiem cho tung nguoi THEO TI LE dong gop
// (quyet dinh 3, thiet ke §2):
//
//	phan[i] = tongDiem × gop[i] / tongGop   (lam tron XUONG)
//
// Lam tron xuong nghia la tong cac phan co the NHO HON tongDiem vai diem — co
// y: khong bao gio phat ra NHIEU HON tongDiem.
//
// Tra ve map RONG (khong phai nil) khi khong co gi de chia — tong dong gop = 0,
// map vao nil/rong, hoac tongDiem <= 0. Nguoi goi luon range duoc tren ket qua
// ma khong phai kiem nil. "KHONG chia khi tong = 0" la bat buoc: chia cho 0 se
// panic, va mot boss chet ma khong ai danh (vd het gio) thi khong ai duoc gi.
//
// Bo qua cac dong gop <= 0 khi tinh tongGop (phong khi map duoc nap tu duong
// khac ngoai CongSatThuong) va KHONG ghi phan = 0 vao ket qua — nguoi gop qua
// it de phan lam tron xuong con 0 thi khong nhan gi, ghi vao chi lam rac map.
func ChiaThuongTheoDongGop(tong map[int64]int64, tongDiem int64) map[int64]int64 {
	ketQua := make(map[int64]int64, len(tong))
	if tongDiem <= 0 {
		return ketQua
	}

	var tongGop int64
	for _, gop := range tong {
		if gop <= 0 {
			continue
		}
		// Chan tran khi cong tong (moi phan tu co the toi MaxInt64 do bao hoa
		// trong CongSatThuong) — bao hoa thay vi cuon vong ve am.
		if tongGop > tongSatThuongToiDa-gop {
			tongGop = tongSatThuongToiDa
			break
		}
		tongGop += gop
	}
	// Chan chia cho 0 — DAY LA LAN CHAN THU HAI (vong duoi da bo qua gop <= 0
	// nen tongGop == 0 keo theo khong con ai de chia). Giu lai co y: nhanRoiChia
	// dung bits.Div64, ham nay PANIC khi mau = 0 — neu mai sau ai do bo bo loc
	// "gop <= 0" o vong duoi thi day la thu duy nhat con dung giua mot boss chet
	// va mot cu panic lam sap ca actor khu vuc. (Kiem thu dot bien: xoa ba dong
	// nay test VAN XANH — dung nhu phan tich tren, day la dot bien TUONG DUONG,
	// khong phai lo hong test.)
	if tongGop <= 0 {
		return ketQua
	}

	for id, gop := range tong {
		if gop <= 0 {
			continue
		}
		phan := nhanRoiChia(tongDiem, gop, tongGop)
		if phan <= 0 {
			continue
		}
		ketQua[id] = phan
	}
	return ketQua
}

// nhanRoiChia tinh floor(a × b / c) CHINH XAC cho a, b, c > 0, khong bi tran o
// buoc nhan: bits.Mul64 giu day du 128 bit cua tich, bits.Div64 chia lai xuong
// 64 bit. Neu tinh thang "a*b/c" bang int64 thi voi dong gop co hang ty (rat
// de xay ra: sat thuong cong don ca tran boss) tich a*b tran ngay, cho ra phan
// thuong AM hoac be xiu ma khong bao loi gi.
//
// bits.Div64 PANIC khi thuong khong vua 64 bit (hi >= c), nen phai chan truoc —
// truong hop do tra ve tran int64 (bao hoa), dung tinh than voi CongSatThuong.
func nhanRoiChia(a, b, c int64) int64 {
	hi, lo := bits.Mul64(uint64(a), uint64(b))
	if hi >= uint64(c) {
		return tongSatThuongToiDa
	}
	q, _ := bits.Div64(hi, lo, uint64(c))
	if q > uint64(tongSatThuongToiDa) {
		return tongSatThuongToiDa
	}
	return int64(q)
}
