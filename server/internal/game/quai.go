// Package game — mo hinh quai (mob) va cach ghi danh sach quai vao goi -24 MAP
// INFO (Ke hoach 6 Task 5). Quai la trang thai THUOC actor KhuVuc: tao luc
// MoKhuVuc, chi luong vongLap cua khu vuc duoc doc/ghi (xem khuvuc.go, cung quy
// tac da ap dung cho nguoiTrong o Ke hoach 5).
package game

import (
	"time"

	"nroserver/internal/protocol"
)

// Quai la mot con quai dang song trong PHAM VI mot KhuVuc. Cac Task sau cua Ke
// hoach 6 (6-9, 13) se doc/ghi them cac truong nay tu trong vongLap (AI, chien
// dau, hoi sinh) — Task 5 chi dung mo hinh + phat trong -24.
type Quai struct {
	// Id la tempId — id TAM THOI trong PHAM VI khu vuc nay (khong phai id toan
	// cuc), dung de cac goi chien dau sau nay (54/-9/-10/-11/-12/-13, Task 6+)
	// tham chieu dung con quai. Client doc byte co dau -> Id/soQuai <= 127
	// (dac ta §1.4).
	Id byte

	// Ma la CAN CUOC cua mau quai — chuoi `ma` trong kho-noi-dung/quai.yaml
	// (KH8 task 9). Moi tra cuu "con quai nay thuoc mau nao" di qua truong NAY.
	//
	// ⚠️ TACH VAI voi MauId ngay ben duoi. Truoc task 9 chi co MauId va no gong
	// HAI vai: vua la id ANH (goi xuong client), vua la khoa BANG ROI. Tu task 8,
	// MauId duoc gan bang `anh_id` cua mau — nen hai mau quai khac han nhau dung
	// chung mot anh se dung chung mot bang roi, im lang. Task 9 cat doi:
	//
	//	Ma    — can cuoc mau quai (bang roi, va moi tra cuu noi dung sau nay)
	//	MauId — CHI la id anh trong catalog, khong con lam khoa tra cuu gi nua
	//
	// Rong la hop le: quai dung trong test cu duoc dung tay khong qua
	// taoQuaiTuAnhChup. Rong -> bang roi MAC DINH (noidung.BangRoiCuaQuai).
	Ma string

	// MauId la id ANH quai (mob template trong catalog anh da ky) — KHONG xuat
	// hien trong than goi -24 (dac ta §3.6 khong co truong nay), nhung can giu
	// noi bo cho duong hien thi hinh. KHONG dung lam khoa tra cuu (xem Ma).
	MauId byte

	// Cap la truong "level" trong dac ta §3.6 (thu tu 9, dung SAU hp, TRUOC
	// hpMax) — cap dung de tinh sat thuong/song sot.
	Cap byte

	// CapTinhAnh la truong "lvMob" trong dac ta §3.6 (thu tu 14, o cuoi than
	// goi) — cap HIEN THI canh ten quai o client, TACH RIENG voi Cap vi dac
	// ta dinh nghia hai o rieng du thuong bang nhau.
	CapTinhAnh byte

	// Hp/HpToiDa la mau int32 4 byte big-endian (dac ta §1.1: readInt phia
	// client, kep tran 2_147_483_647, kep san 0).
	Hp      int32
	HpToiDa int32

	// X/Y la toa do quai trong ban do (short, dac ta §3.6 truong 11-12).
	X int16
	Y int16

	// TrangThai la truong "status" (dac ta §3.6 truong 13) — 0 = binh thuong;
	// Task sau (choang, buff...) se dung gia tri khac.
	TrangThai byte

	// --- Ba truong duoi day la Ke hoach 6 Task 8 (hoi sinh 5s + hoi mau 30s,
	// dac ta §5.2) — trang thai NOI BO cua actor, KHONG xuat hien trong than
	// goi -24 (ghiQuai o duoi khong dung toi). Chi luong vongLap cua KhuVuc
	// duoc doc/ghi (dac ta §1.6), giong moi truong khac cua Quai.

	// LucChet la thoi diem Hp xuong <= 0 GAN NHAT — nhipTim (khuvuc.go) dung
	// truong nay de tinh "da troi qua bao lau tu luc chet" cho nguong 5 giay
	// hoi sinh. Duoc gan trong case suKienDanhQuai (khuvuc.go) ngay khi phat
	// hien Hp <= 0, TRUOC khi phat goi -12. Zero value (chua tung chet) khong
	// gay hai: quai moi taoQuaiMacDinh luon Hp > 0 nen nhanh "dang chet" cua
	// nhipTim khong bao gio doc truong nay khi con zero.
	LucChet time.Time

	// LucHoiSinh la thoi diem quai HOI SINH GAN NHAT (nhipTim gan khi qua du
	// 5 giay ke tu LucChet). Dung lam ViTriCoMau.LucHoiSinh (chotkiem.go) khi
	// nguoi choi xin danh quai nay — kich hoat bat tu 1500ms sau hoi sinh (dac
	// ta §1.7), CUNG CO CHE voi nguoiTrongKhuVuc ma Task 13 se dung cho nguoi
	// choi. Zero value (quai moi spawn, CHUA TUNG hoi sinh lan nao — khac voi
	// "vua hoi sinh") khien bayGio.Sub(zero) ra khoang RAT LON -> kiem bat tu
	// tu nhien CHO PHEP ngay, dung hanh vi mong doi: quai moi mo ban do khong
	// bi coi la "vua hoi sinh" nen khong duoc mien nhiem oan.
	LucHoiSinh time.Time

	// LanHoiMauCuoi la thoi diem GAN NHAT nhipTim DA KIEM (khong nhat thiet
	// DA HOI — neu dang day mau thi kiem xong khong hoi gi) chu ky hoi mau 30
	// giay (dac ta §5.2). Khoi tao = thoi diem spawn (MoKhuVuc gan sau khi goi
	// taoQuaiMacDinh) de lan kiem DAU TIEN dung 30 giay sau khi ban do mo, chu
	// khong hoi ngay lap tuc do zero value nam 1 tao ra khoang cach "da troi
	// qua" khong lo. Sau moi lan hoi sinh (nhipTim), truong nay CUNG duoc dat
	// lai ve thoi diem hoi sinh — tranh mot con quai vua hoi (da day mau) van
	// bi tinh "no du 30s roi" tu truoc luc chet va hoi "thua" ngay tick ke tiep.
	LanHoiMauCuoi time.Time

	// --- Ba truong duoi day la Ke hoach 6 Task 9 (AI hai tang uu tien + goi
	// -11/-10, dac ta §5.2) — trang thai NOI BO cua actor, CHI luong vongLap
	// duoc doc/ghi, dung nguyen tac da ap dung cho ba truong Task 8 o tren.

	// DaTungDanh la TAP HOP id nhung nguoi choi DA TUNG danh TRUNG quai nay
	// (dac ta §5.2 tang uu tien 1: "nguoi da tung danh no"). Nap mot phan tu
	// (id -> true) trong case suKienDanhQuai (khuvuc.go), CHI trong nhanh don
	// TRUNG (khong phai nhanh Truot=true — "trung don" nghia la thuc su gay
	// sat thuong, khac "vuot tam" noi don khong thuc su cham muc tieu). XOA
	// SACH (gan nil) khi quai CHET (mucTieu.Hp <= 0, cung case do) — dung
	// nguyen van brief Task 9: "danh sach da tung danh no nap khi trung don,
	// xoa sach khi quai chet". map nil van doc an toan (tra false cho moi
	// khoa, len()==0) nen khong can khoi tao rong o taoQuaiMacDinh.
	DaTungDanh map[int64]bool

	// LanDanhCuoi la thoi diem CHINH QUAI NAY that su danh TRUNG mot nguoi
	// choi GAN NHAT — khac han nguoiTrongKhuVuc.lanDanhCuoi (do la nguoi choi
	// danh quai; day la quai danh nguoi choi, huong nguoc lai). aiQuai
	// (khuvuc.go) CHI dat lai truong nay trong nhanh "that su danh" (sau khi
	// da kiem bayGio.Sub(LanDanhCuoi) >= nhip ma chonMucTieu tra ve) — TUYET
	// DOI khong dat lai khi chonMucTieu tra ve muc tieu nil (dac ta §5.4.2:
	// bay ban Java dat lai moc thoi gian NGOAI nhanh "co muc tieu", khien
	// quai vua thay nguoi phai cho THEM mot nhip nua moi duoc danh, thay vi
	// danh NGAY — xem chu thich ham aiQuai). Zero value (chua tung danh ai)
	// khien bayGio.Sub(zero) ra khoang RAT LON -> lan danh DAU TIEN tu nhien
	// duoc phep ngay, khong can nhanh if rieng cho truong hop nay.
	LanDanhCuoi time.Time

	// SatThuongCoBan la sat thuong CO DINH quai gay ra MOI don len nguoi
	// choi. KH8 task 8: gia tri nay khong con la hang so 5 go trong code nua —
	// no den TU MAU QUAI trong kho-noi-dung/quai.yaml, gan MOT LAN luc sinh
	// (taoQuaiTuAnhChup, khoquai.go).
	//
	// Vi gan mot lan luc sinh, con quai nay GIU SO CUA THE HE LUC NO SINH RA:
	// nap nong quai.yaml khong doi sat thuong cua no. Da chot, co ly do —
	// xem khoi chu thich dau khoquai.go.
	//
	// Dac ta §5.3 con dinh nghia cong thuc DAY DU ("dame rieng ±1%" hoac
	// "hpMax×pDame%±theo cap"); day moi la o "dame rieng" cua cong thuc do.
	// Hai o con lai (pDame%, bien thien theo cap) chua co — chung la viec cua
	// buoc DO CAN BANG, khong phai cua buoc chuyen cho o.
	SatThuongCoBan int32
}

// ghiQuai ghi MOT con quai vao than goi -24 MAP INFO, DUNG THEO THU TU dac ta
// §3.6. KH6 chua co quai bi vo hieu hoa/khong di chuyen duoc/nguyen to lua-bang-
// gio hay la BIG BOSS nen 5 co dau va isBigBoss LUON ghi false; sys LUON = 0.
//
// ⚠️ BAY LON NHAT (§3.6, ly do co test rieng TestGhiQuaiThuTuHpCapHpToiDa): thu
// tu la hp -> level -> hpMax, KHONG PHAI hp -> hpMax -> level nhu truc giac hay
// nghi. Doi thu tu hai truong nay van cho ra tong do dai dung (client khong bao
// loi) nhung gia tri doc duoc SAI HOAN TOAN — dung dac loai "sai thu tu, tong do
// dai van khop" o bang §2.0 cua dac ta: mau am, cap ra tuong hp, khong log nao
// bao ca.
func ghiQuai(w *protocol.Writer, q *Quai) {
	w.Bool(false) // 1 isDisable
	w.Bool(false) // 2 isDontMove
	w.Bool(false) // 3 isFire
	w.Bool(false) // 4 isIce
	w.Bool(false) // 5 isWind

	w.Byte(q.Id) // 6 tempId
	w.Byte(0)    // 7 sys = 0

	w.Int(q.Hp)      // 8  hp    — TRUOC level, DUNG doi cho hpMax
	w.Byte(q.Cap)    // 9  level — GIUA hp va hpMax
	w.Int(q.HpToiDa) // 10 hpMax — SAU level

	w.Short(q.X) // 11 x
	w.Short(q.Y) // 12 y

	w.Byte(q.TrangThai)  // 13 status
	w.Byte(q.CapTinhAnh) // 14 lvMob

	w.Bool(false) // 15 isBigBoss
}

// ghiDanhSachQuai ghi soQuai (byte) roi tung con quai theo dac ta §3.6, dung
// trong goi -24 MAP INFO. soQuai la byte CO DAU o phia client (dac ta §1.4) nen
// PHAI <= 127; qua han se bi CAT BOT (khong ghi qua 127 con, khong ghi so am)
// de khong lam hong khung goi — day la loi cau hinh ban do (qua nhieu quai),
// khong phai loi runtime nen khong panic, chi cat cho an toan.
func ghiDanhSachQuai(w *protocol.Writer, ds []*Quai) {
	n := len(ds)
	if n > 127 {
		n = 127
	}
	w.Byte(byte(n))
	for i := 0; i < n; i++ {
		ghiQuai(w, ds[i])
	}
}

// KH8 task 8: BANG QUAI CUNG DA BI XOA KHOI TEP NAY.
//
// Truoc task 8, o day co ham taoQuaiMacDinh dung 2 con quai voi so lieu go
// thang vao code (hp tu catalog.json, SatThuongCoBan la hang so 5, diem sinh
// dung chung cho moi ban do). Ca ba thu do gio nam trong kho-noi-dung/quai.yaml
// va duong sinh quai DUY NHAT la taoQuaiTuAnhChup (khoquai.go) — doc tu ANH CHUP
// noi dung. Khong con bang nao trong Go de "roi ve" ca.

// ---------------------------------------------------------------------------
// Ke hoach 6 Task 9 — AI quai hai tang uu tien (dac ta §5.2), goi -11/-10.
// ---------------------------------------------------------------------------

// Bon hang so hanh vi AI DA CHOT (dac ta §5.2, nguon Mob.java:423/438 va
// Mob.java:398/412) — nguon su that duy nhat, dung trong chonMucTieu ben
// duoi. Dat kieu hang so (khong phai bien) de kiem thu dot bien sua truc
// tiep gia tri o day roi chay lai test bien, khong can sua noi khac.
const (
	// tamTraThu/nhipTraThu: tang uu tien 1 — nguoi DA TUNG danh quai nay.
	// Tam XA HON va nhip NHANH HON tang 2 — "da danh thi phai chiu tran"
	// (dac ta §5.2, dong "Luat 'tra thu nhanh gap doi, tam gap ba' la luat
	// choi tot ... da danh thi phai chiu tran").
	tamTraThu  int64         = 300
	nhipTraThu time.Duration = 1000 * time.Millisecond

	// tamTuGaySu/nhipTuGaySu: tang uu tien 2 — nguoi BAT KY di ngang qua.
	// Tam GAN HON va nhip CHAM HON tang 1 — "nguoi choi di ngang qua thi
	// tuong doi an toan" (dac ta §5.2).
	tamTuGaySu  int64         = 100
	nhipTuGaySu time.Duration = 2000 * time.Millisecond
)

// chonMucTieu chon MUC TIEU + NHIP DANH cho quai q trong so danh sach nguoi
// choi hien co trong khu vuc (dac ta §5.2, brief Task 9). HAM THUAN: chi doc
// q.X/Y/DaTungDanh va tung phan tu cua nguoi, KHONG sua bat ky trang thai nao
// (khong ghi q.LanDanhCuoi, khong ghi gi vao tung nguoiTrongKhuVuc) — noi GOI
// (aiQuai, khuvuc.go) moi la noi quyet dinh "co that su danh hay khong" va
// cham q.LanDanhCuoi. Dung dung dac diem ban giao cua brief: "tra ve CA muc
// tieu LAN nhip danh, KHONG sua trang thai ngam" — sua loi #4 cua dac ta §5.4
// ("doi nhip danh nhu tac dung phu cua ham tim muc tieu", Mob.java:412,438).
//
// Tra ve (nil, 0) khi khong tim thay muc tieu o CA HAI tang — ben goi (aiQuai)
// PHAI hieu day la "khong lam gi ca", TUYET DOI khong duoc dung ket qua nay
// de dat lai q.LanDanhCuoi (dac ta §5.4.2: quai vua thay nguoi phai duoc
// danh NGAY, khong phai cho them mot nhip — xem chu thich ham aiQuai).
//
// NGOAC DIEU KIEN LOC TUONG MINH (dac ta §5.4.3): moi dieu kien loc (con
// song, da tung danh, trong tam) la MOT cau if/return DOC LAP trong closure
// hopLe truyen cho timGanNhatTrongTam — KHONG gop chung thanh MOT bieu thuc
// bool dai roi dua vao thu tu &&/||. Day chinh la bay ban Java mac phai
// (Mob.java:402-404): voi big boss, sai thu tu toan tu lam MOI bo loc phia
// truoc — KE CA bo loc "chua chet" — bi bo qua, boss nham ca xac chet.
//
// bayGio (Ke hoach 6 Task 11 Phan B, BAN GIAO tu Task 13): THAM SO MOI —
// truoc Task 11 ham nay KHONG nhan bayGio nen KHONG THE kiem bat tu, day
// chinh la LO HONG dac ta §1.7 ("bat tu 1500ms sau hoi sinh — MOI don vao
// muc tieu vua hoi sinh bi tu choi, KHONG gioi han huong"): suKienDanhQuai
// (khuvuc.go) co KiemDanhHopLe kiem bat tu cho HUONG nguoi-danh-quai, nhung
// KHONG co gi kiem HUONG quai-danh-nguoi ca — mot nguoi vua hoi sinh se bi
// quai AI danh trung NGAY o tick ke tiep (100ms), vi pham "KHONG gioi han
// huong" cua dac ta. Them bayGio + dieu kien hopLeBatTu duoi day o CA HAI
// tang de va lo hong nay, TAI DUNG thoiGianBatTuSauHoiSinh (chotkiem.go, hang
// so da chot chung cho toan bo he chien dau — KHONG dinh nghia hang so rieng
// o day). Ham VAN LA THUAN (bayGio la tham so, khong doc dong ho he thong).
func (q *Quai) chonMucTieu(nguoi []*nguoiTrongKhuVuc, bayGio time.Time) (*nguoiTrongKhuVuc, time.Duration) {
	// hopLeBatTu: n VUA hoi sinh (< 1500ms truoc bayGio) thi KHONG duoc chon
	// lam muc tieu — dac ta §1.7 "KHONG gioi han huong", ap dung cho CA HAI
	// tang uu tien ben duoi. n.lucHoiSinh la zero value (time.Time{}) cho
	// nguoi CHUA TUNG hoi sinh lan nao -> bayGio.Sub(zero) ra khoang RAT LON
	// -> tu nhien hopLe (KHONG bi mien nhiem oan), dung tinh than voi moi noi
	// khac dung zero value cua truong lucHoiSinh trong file nay.
	hopLeBatTu := func(n *nguoiTrongKhuVuc) bool {
		return bayGio.Sub(n.lucHoiSinh) >= thoiGianBatTuSauHoiSinh
	}

	// Tang 1 (uu tien cao hon): nguoi DA TUNG danh quai nay, con song, tam 300,
	// KHONG dang trong 1500ms bat tu sau hoi sinh.
	if mt := timGanNhatTrongTam(q, nguoi, tamTraThu, func(n *nguoiTrongKhuVuc) bool {
		if n.nv.HpHienTai <= 0 { // (a) da chet -> loai, DOC LAP voi (b)/(c)
			return false
		}
		if !q.DaTungDanh[n.nv.ID] { // (b) chua tung danh quai nay -> loai
			return false
		}
		if !hopLeBatTu(n) { // (c) vua hoi sinh <1500ms -> bat tu, Ke hoach 6 Task 11 Phan B
			return false
		}
		return true
	}); mt != nil {
		return mt, nhipTraThu
	}

	// Tang 2: bat ky ai con song, tam 100, KHONG dang trong 1500ms bat tu.
	if mt := timGanNhatTrongTam(q, nguoi, tamTuGaySu, func(n *nguoiTrongKhuVuc) bool {
		if n.nv.HpHienTai <= 0 { // (a) da chet -> loai, DOC LAP voi (b)
			return false
		}
		if !hopLeBatTu(n) { // (b) vua hoi sinh <1500ms -> bat tu, Ke hoach 6 Task 11 Phan B
			return false
		}
		return true
	}); mt != nil {
		return mt, nhipTuGaySu
	}

	return nil, 0
}

// timGanNhatTrongTam duyet nguoi, loc bang hopLe(n), giu lai nguoi GAN q NHAT
// trong pham vi tam (so sanh khoang cach BINH PHUONG, tranh sqrt/float — cung
// ky thuat voi kiemDiChuyenHopLe va KiemDanhHopLe da dung o cac Task truoc).
// nil-safe: phan tu nil hoac n.nv nil bi bo qua (phong thu, khong panic) —
// ve ly thuyet khong xay ra (k.nguoiTrong luon gan nv khac nil khi them vao,
// xem case suKienVao), nhung giu cho an toan khong doi hoi nguoi goi dam bao.
func timGanNhatTrongTam(q *Quai, nguoi []*nguoiTrongKhuVuc, tam int64, hopLe func(*nguoiTrongKhuVuc) bool) *nguoiTrongKhuVuc {
	tamSq := tam * tam
	var ganNhat *nguoiTrongKhuVuc
	var ganNhatSq int64
	for _, n := range nguoi {
		if n == nil || n.nv == nil {
			continue
		}
		if !hopLe(n) {
			continue
		}
		dx := int64(n.nv.ViTriX) - int64(q.X)
		dy := int64(n.nv.ViTriY) - int64(q.Y)
		distSq := dx*dx + dy*dy
		if distSq > tamSq {
			continue
		}
		if ganNhat == nil || distSq < ganNhatSq {
			ganNhat = n
			ganNhatSq = distSq
		}
	}
	return ganNhat
}
