package game

import (
	"context"
	"errors"
)

// SoKhuVuc giu 5 khu vuc actor co dinh (mapId 1-5) — tao san luc khoi dong may
// chu, song suot doi tien trinh. Don gian hon tao/huy dong; 5 goroutine ngu la
// chi phi khong dang ke. Tao/huy dong de danh khi co map tam (to doi/su kien).
type SoKhuVuc struct {
	theoMap map[int32]*KhuVuc

	// soNo giu chinh so no da cam vao cac khu vuc, de TatMayChu co the dong no
	// ma KHONG phai doi chu ky (main.go thuoc pham vi Task 14, dang co agent
	// khac lam — tranh cham vao la tranh dung do).
	soNo *SoNoThuong
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

// GanSoNo noi so no phan thuong (sonothuong.go) vao ca 5 khu vuc. Khac GanLoa
// o mot diem: KHONG tu khoi dong gi ca — MoSoNoThuong da chay goroutine ghi
// ngay khi mo, nen o day chi con viec cam day.
//
// CO Y KHONG goi tu MoSoKhuVuc, cung ly do voi GanLoa (luat 4: MoSoKhuVuc phai
// GIU NGUYEN hanh vi). Chi cmd/nroserver/main.go goi ham nay; khi tat may chu
// phai goi soNo.DongVaChoXong TRUOC s.DungTatCa() de khoan no cuoi cung kip
// xuong dia truoc khi khong con ai xep them.
func (s *SoKhuVuc) GanSoNo(so *SoNoThuong) {
	for _, k := range s.theoMap {
		k.soNo.Store(so)
	}
	s.soNo = so
}

// DongSoNo xa not so no roi dung goroutine ghi cua no. Tra nil khi chua noi
// day so no (kho khong biet ghi no) — khong phai loi.
//
// GOI SAU DungTatCa, khong truoc: chung nao con mot actor song thi con co the
// co boss chet va co them khoan no duoc xep. Dong truoc la vut di dung nhung
// khoan cuoi cung.
func (s *SoKhuVuc) DongSoNo(ctx context.Context) error {
	if s.soNo == nil {
		return nil
	}
	return s.soNo.DongVaChoXong(ctx)
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
