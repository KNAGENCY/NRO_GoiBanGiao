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
// chan o buoc sau), khac hanh vi ban goc.
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

// KiemDanhHopLe chot 5 kiem truoc khi cho phep nguoiDanh dung chieuDangChon
// cua no len mucTieu, tai thoi diem bayGio. HAM THUAN: khong sua nguoiDanh,
// khong doc dong ho he thong — moi thu qua tham so.
func KiemDanhHopLe(nguoiDanh *nguoiTrongKhuVuc, mucTieu ViTriCoMau, bayGio time.Time) KetQuaKiem {
	// Bao ve ngoai bang 5-dong: chua chon chieu thi khong co gi de kiem tam/
	// chakra/hoi chieu ca — tranh panic khi nguoiDanh.chieuDangChon == nil.
	if nguoiDanh == nil || nguoiDanh.chieuDangChon == nil {
		return KetQuaKiem{ChoPhep: false, LyDo: "chua chon chieu"}
	}
	cap := nguoiDanh.chieuDangChon

	// #1 Bat tu: muc tieu vua hoi sinh < 1500ms truoc bayGio.
	if bayGio.Sub(mucTieu.LucHoiSinh) < thoiGianBatTuSauHoiSinh {
		return KetQuaKiem{ChoPhep: false, LyDo: "muc tieu dang bat tu"}
	}

	// #2 Chakra: nguoiDanh khong du chakra cho chieu dang chon.
	if nguoiDanh.nv.ChakraHienTai < int32(cap.ChakraDung) {
		return KetQuaKiem{ChoPhep: false, LyDo: "khong du chakra"}
	}

	// #3 Hoi chieu: daTroiQua tu lan dung chieu truoc phai >= hoiChieu - noi.
	daTroiQua := bayGio.Sub(nguoiDanh.lanDanhCuoi)
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
	dx := int64(mucTieu.X) - int64(nguoiDanh.nv.ViTriX)
	dy := int64(mucTieu.Y) - int64(nguoiDanh.nv.ViTriY)
	distSq := dx*dx + dy*dy
	tamSq := int64(cap.Dx) * int64(cap.Dx)
	if distSq > tamSq {
		// Vuot tam: VAN CHO PHEP (phat goi) nhung TRUOT — dac ta §1.9. KHONG
		// duoc im lang bo qua o day, khac voi 4 kiem tren.
		return KetQuaKiem{ChoPhep: true, Truot: true, LyDo: "vuot tam"}
	}

	return KetQuaKiem{ChoPhep: true}
}
