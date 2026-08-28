package game

import "errors"

// SoKhuVuc giu 5 khu vuc actor co dinh (mapId 1-5) — tao san luc khoi dong may
// chu, song suot doi tien trinh. Don gian hon tao/huy dong; 5 goroutine ngu la
// chi phi khong dang ke. Tao/huy dong de danh khi co map tam (to doi/su kien).
type SoKhuVuc struct {
	theoMap map[int32]*KhuVuc
}

// MoSoKhuVuc tao 5 khu vuc mapId 1-5 va KHOI DONG actor cua tung cai.
func MoSoKhuVuc() *SoKhuVuc {
	s := &SoKhuVuc{theoMap: make(map[int32]*KhuVuc)}
	for mapId := int32(1); mapId <= 5; mapId++ {
		k := MoKhuVuc(TimBanDo(mapId))
		k.BatDau()
		s.theoMap[mapId] = k
	}
	return s
}

// Tra tra khu vuc theo mapId, nil neu khong co.
func (s *SoKhuVuc) Tra(mapId int32) *KhuVuc {
	return s.theoMap[mapId]
}

// PhatToanServer dua MOT dong thong bao (goi 93) toi MOI nguoi trong MOI khu
// vuc. Day la ham ma BoPhatThanh goi o cuoi duong day — no chay tren GOROUTINE
// LOA, KHONG PHAI tren actor nao ca.
//
// Vi chay ngoai actor, no TUYET DOI khong duoc dung guiDieuKhien (cho toi
// thoiGianChoControl = 100ms tren hop thu ma vongLap co the dang ban, xem
// khuvuc.go) — thay vao do gui suKienPhatGoi qua hop thu GAMEPLAY bang k.Gui,
// von khong bao gio chan: hop thu day thi bo goi va dem soLanGuiRoi. Nho vay
// khong ai cho ai: loa khong cho actor, actor khong cho loa.
//
// Goi 93 duoc dung MOT LAN roi dung chung con tro cho ca 5 khu vuc —
// protocol.Message chi duoc DOC khi gui (khong ai sua no) nen dung chung la an
// toan va tiet kiem. noiDung qua dai (>65535 byte sau ma hoa) khien
// goiThongBaoToanServer tra nil; khi do khong phat gi ca — mot goi thieu byte
// con te hon khong goi vi client se doc lech ca luong tin phia sau.
func (s *SoKhuVuc) PhatToanServer(noiDung string) {
	m := goiThongBaoToanServer(noiDung)
	if m == nil {
		return
	}
	for _, k := range s.theoMap {
		k.Gui(suKienPhatGoi{m: m})
	}
}

// GanLoa noi bo phat thanh vao ca 5 khu vuc VA khoi dong no voi duong phat la
// chinh s.PhatToanServer — mot loi goi lam tron ca hai dau day, de khong ai
// gan loa xong roi quen bat (thong bao se nam mai trong hang doi khong ai rut).
//
// CO Y KHONG goi tu MoSoKhuVuc (luat 4 cua du an: MoSoKhuVuc phai GIU NGUYEN
// hanh vi, khong tu khoi dong them thu gi — hang loat test tich hop dung no de
// dung 5 khu vuc roi kiem "im lang"). Chi cmd/nroserver/main.go goi ham nay
// (Task 14), va khi tat may chu phai goi loa.Dung() TRUOC s.DungTatCa() de
// khong con ai gui vao khu vuc dang dong.
func (s *SoKhuVuc) GanLoa(loa *BoPhatThanh) {
	for _, k := range s.theoMap {
		k.loa.Store(loa)
	}
	loa.BatDau(s.PhatToanServer)
}

// LoaDaGan tra ve cai loa dang cam vao TAT CA khu vuc, hoac nil neu con khu
// vuc nao chua duoc gan (hoac cac khu vuc dang cam hai cai loa khac nhau).
//
// Ham nay ton tai VI MOT LY DO DUY NHAT: cho cmd/nroserver khoa duoc bang test
// rang main.go THAT SU goi GanLoa. Truong k.loa la noi bo cua package nay, nen
// neu khong co cua so nay thi tu ben ngoai KHONG CACH NAO phan biet "da noi
// day" voi "quen noi day" — ma "quen noi day" chinh la lo hong Task 14 di va:
// xepThongBao thay loa nil thi LANG LE return, khong mot dong log.
//
// Tra ve chinh con tro (khong phai bool) de test con khang dinh duoc no la
// DUNG cai loa ma main() dang giu tay cam — gan mot cai roi tat mot cai khac
// thi goroutine phat con song mai sau khi may chu tuong da tat sach.
func (s *SoKhuVuc) LoaDaGan() *BoPhatThanh {
	var chung *BoPhatThanh
	for _, k := range s.theoMap {
		l := k.loa.Load()
		if l == nil || (chung != nil && l != chung) {
			return nil
		}
		chung = l
	}
	return chung
}

// DungTatCa dung sach moi khu vuc (goi luc tat may chu) va tra loi neu bat
// ky actor nao khong thoat. Van thu dung HET khu vuc de mot actor ket khong
// lam bo qua viec don dep cac actor con lai.
func (s *SoKhuVuc) DungTatCa() error {
	var loi error
	for _, k := range s.theoMap {
		loi = errors.Join(loi, k.Dung())
	}
	return loi
}
