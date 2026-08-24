// chotkiem.go — chot kiem truoc khi cho danh (Ke hoach 6, Task 3).
//
// KiemDanhHopLe la HAM THUAN: khong doc/ghi trang thai khu vuc, nhan bayGio lam
// tham so — de test duoc MOI moc thoi gian (bat tu, hoi chieu) ma khong can
// sleep that. Actor (khuvuc.go, Task 4+) se goi ham nay TRONG vongLap, dung
// gio That() cua may (hoac gio ghim trong test).
//
// THU TU kiem dung y het ban Java (SkillService.java:58-136, dac ta §9 phu
// luc "Chot kiem truoc khi danh"): bat tu -> chakra -> hoi chieu -> chet ->
// tam. DUNG SAI THU TU la sai vi cac kiem KHONG doc lap — vi du kiem chakra
// truoc bat tu se cho phep tieu chakra vao mot don khong bao gio trung (bat tu
// chan o buoc sau), khac hanh vi ban goc. TestThuTuKiemGiuNguyen
// (chotkiem_test.go) khoa dung thu tu nay.
//
// KE HOACH 6b (Task 5, don nen truoc phan boss): ham nay TRUOC day nhan thang
// *nguoiTrongKhuVuc, tuc chi NGUOI CHOI moi danh duoc — boss la kieu rieng,
// khong co store.NhanVat, khong co bang chakra, nen khong the truyen vao. Nay
// no nhan kieu du lieu thuan BenDanh de CA HAI BEN di chung mot duong kiem:
// mot bo luat, mot cho sua, khong co ban sao logic thu hai cho boss.
package game

import "time"

// Hai hang so da CHOT (dac ta §1 diem 7 va 8; brief Task 3 — "Hang so da
// chot"). KHONG lay tu dau khac, day la nguon su that duy nhat cho hai con so
// nay trong tang chien dau.
const (
	// thoiGianBatTuSauHoiSinh: moi don vao muc tieu vua hoi sinh trong khoang
	// nay bi tu choi (dac ta §1.7).
	thoiGianBatTuSauHoiSinh = 1500 * time.Millisecond

	// hoiChieuNoi: cho phep dung chieu khi daTroiQua >= coolDown - hoiChieuNoi
	// (dac ta §1.8) — bu do tre mang round-trip, tranh client bam dung gio
	// nhung server tinh tre 1 tick roi tu choi oan.
	hoiChieuNoi = 50 * time.Millisecond
)

// BenDanh la BEN CHU DONG cua mot don danh: chinh xac nhung gi KiemDanhHopLe
// can doc tu ke tan cong, khong hon. Kieu du lieu THUAN (khong con tro toi
// trang thai khu vuc, tru con tro chi-doc toi cap chieu) nen nguoi choi va
// boss deu dung duoc — xem benDanhTuNguoi (khuvuc.go) cho phia nguoi choi.
type BenDanh struct {
	// X, Y la vi tri ben danh — int32 DUNG kieu store.NhanVat.ViTriX/Y
	// (store.go:38-39), KHONG thu hep xuong int16.
	//
	// Vi sao khong int16 du ViTriCoMau (muc tieu) la int16: ban truoc khi tach
	// ham doc THANG int64(nguoiDanh.nv.ViTriX) tu int32, nen ep xuong int16 o
	// adapter la mot phep THU HEP mieng am tham — toa do ngoai [-32768, 32767]
	// se cuon vong va kiem tam (#5) tinh tren con so sai. Hien tai bien map
	// chan &lt; 4096 nen chua vo, nhung do la trung hop chu khong phai bao dam.
	// Giu int32 thi kiem #5 (da tinh bang int64) khong doi mot phep nao, va
	// KHONG doi mot byte nao tren day: goi -5 van tu ep int16 rieng o cho no.
	X, Y int32

	// ChakraHienTai: int32, DUNG kieu store.NhanVat.ChakraHienTai (store.go)
	// de phia nguoi choi chep sang khong phai ep kieu, khong mat gia tri.
	ChakraHienTai int32

	// LanDanhCuoi la thoi diem ben nay DUNG THANH CONG mot chieu lan gan nhat.
	// Zero value (chua tung dung chieu nao) khien bayGio.Sub(LanDanhCuoi) ra
	// mot khoang RAT LON, nen kiem hoi chieu tu nhien CHO PHEP ngay lan dau —
	// khong can if/else rieng.
	LanDanhCuoi time.Time

	// Chieu la cap chieu dang chon. nil = chua chon chieu -> tu choi ngay o
	// bao ve dau ham (truoc ca kiem bat tu).
	Chieu *CapChieu

	// VoHanChakra BO QUA hoan toan kiem chakra (#2). Danh cho BOSS: boss
	// khong co bang chakra nao ca nen "du chakra khong" la cau hoi vo nghia
	// voi no. Noi dieu do BANG KIEU DU LIEU thay vi nhoi mot con so ma thuat
	// (vd cho ChakraHienTai = 31_07_2002 roi cau mong khong chieu nao ton hon
	// the) — so ma thuat kieu do la mot qua bom hen gio: no dung cho toi khi
	// ai do them mot chieu ton nhieu chakra hon, va luc vo thi khong co mot
	// dong code nao giai thich con so ay.
	//
	// Voi NGUOI CHOI truong nay LUON false (benDanhTuNguoi, khuvuc.go).
	// No CHI mien kiem #2, KHONG mien bat tu / hoi chieu / muc tieu chet /
	// tam — TestVoHanChakraKHONGBoQuaKiemKhac khoa dieu do.
	VoHanChakra bool
}

// ViTriCoMau la vi tri + luong mau cua MOT MUC TIEU (nguoi choi hoac quai) can
// de chot kiem truoc khi danh. Kieu MOI — brief Task 3 cho phep dinh nghia tai
// day vi chua ton tai o dau khac trong package.
//
// LucHoiSinh de gia tri ZERO (time.Time{}) khi muc tieu CHUA TUNG hoi sinh lan
// nao (vi du quai/nguoi choi tu luc vao map): bayGio.Sub(zero) luon ra mot
// khoang thoi gian RAT LON (zero la nam 1), nen kiem bat tu tu nhien CHO PHEP
// ma khong can nhanh if rieng cho truong hop "chua tung chet".
type ViTriCoMau struct {
	X, Y int16

	// HpHienTai <= 0 nghia la muc tieu da chet (kiem #4).
	HpHienTai int32

	// LucHoiSinh la thoi diem muc tieu hoi sinh GAN NHAT.
	LucHoiSinh time.Time
}

// KetQuaKiem la ket qua cua KiemDanhHopLe.
type KetQuaKiem struct {
	// ChoPhep: false o 4 kiem dau (bat tu/chakra/hoi chieu/chet) nghia la
	// KHONG duoc phat goi gi ca (im lang bo qua goi len). true o kiem #5
	// (vuot tam) VA khi ca 5 kiem deu qua — hai truong hop nay phan biet
	// bang Truot.
	ChoPhep bool

	// Truot = true CHI xay ra khi vuot tam (kiem #5 hong nhung ChoPhep van
	// true — dac ta §1.9: "vuot tam danh -> truot = true, sat thuong 0,
	// NHUNG VAN PHAT GOI"). Noi goi ham nay dung goi (Task 10:
	// goiDanhNguoiChoi / goi -9) PHAI doc Truot va EP chiMang=false khi
	// Truot=true — dac ta §3.3 diem 3: "client chi hien chu miss khi
	// chiMang == false". KetQuaKiem KHONG tu no co truong chiMang (brief
	// chot dung 3 truong ChoPhep/Truot/LyDo); day la hop dong cho ben goi.
	Truot bool

	// LyDo la chuoi tieng Viet ngan giai thich vi sao hong, dung de log/debug.
	// Khong phat ra client.
	LyDo string
}

// KiemDanhHopLe chot 5 kiem truoc khi cho phep benDanh dung chieu cua no len
// mucTieu, tai thoi diem bayGio. HAM THUAN: khong sua gi, khong doc dong ho
// he thong — moi thu qua tham so.
func KiemDanhHopLe(benDanh BenDanh, mucTieu ViTriCoMau, bayGio time.Time) KetQuaKiem {
	// Bao ve ngoai bang 5-dong: chua chon chieu thi khong co gi de kiem tam/
	// chakra/hoi chieu ca — tranh panic khi benDanh.Chieu == nil.
	if benDanh.Chieu == nil {
		return KetQuaKiem{ChoPhep: false, LyDo: "chua chon chieu"}
	}
	cap := benDanh.Chieu

	// #1 Bat tu: muc tieu vua hoi sinh < 1500ms truoc bayGio.
	if bayGio.Sub(mucTieu.LucHoiSinh) < thoiGianBatTuSauHoiSinh {
		return KetQuaKiem{ChoPhep: false, LyDo: "muc tieu dang bat tu"}
	}

	// #2 Chakra: benDanh khong du chakra cho chieu dang chon. Boss
	// (VoHanChakra) bo qua HAN buoc nay — xem chu thich truong do o tren.
	if !benDanh.VoHanChakra && benDanh.ChakraHienTai < int32(cap.ChakraDung) {
		return KetQuaKiem{ChoPhep: false, LyDo: "khong du chakra"}
	}

	// #3 Hoi chieu: daTroiQua tu lan dung chieu truoc phai >= hoiChieu - noi.
	daTroiQua := bayGio.Sub(benDanh.LanDanhCuoi)
	hoiChieu := time.Duration(cap.HoiChieuMs) * time.Millisecond
	if daTroiQua < hoiChieu-hoiChieuNoi {
		return KetQuaKiem{ChoPhep: false, LyDo: "chieu chua hoi"}
	}

	// #4 Muc tieu da chet.
	if mucTieu.HpHienTai <= 0 {
		return KetQuaKiem{ChoPhep: false, LyDo: "muc tieu da chet"}
	}

	// #5 Tam: so binh phuong khoang cach de tranh sqrt/float (cung cach lam
	// voi kiemDiChuyenHopLe o dichuyen.go). cap.Dx la tam danh cua chieu (xem
	// chu thich truong Dx trong chieuthuc.go: "client doc ca hai, dung chung
	// 1 gia tri tam").
	dx := int64(mucTieu.X) - int64(benDanh.X)
	dy := int64(mucTieu.Y) - int64(benDanh.Y)
	distSq := dx*dx + dy*dy
	tamSq := int64(cap.Dx) * int64(cap.Dx)
	if distSq > tamSq {
		// Vuot tam: VAN CHO PHEP (phat goi) nhung TRUOT — dac ta §1.9. KHONG
		// duoc im lang bo qua o day, khac voi 4 kiem tren.
		return KetQuaKiem{ChoPhep: true, Truot: true, LyDo: "vuot tam"}
	}

	return KetQuaKiem{ChoPhep: true}
}
