// lang.go — 5 LANG NARUTO.
//
// (KH8 task 6. Thiet ke: docs/superpowers/specs/2026-09-03-kh8-kho-noi-dung-design.md
// hang "6" bang Task §7, §2.2 bang can cuoc "Lang".)
//
// Tu task nay, hang so KHOI DAU cua tung lang nam trong kho-noi-dung/lang.yaml,
// KHONG con la mot bang Go cung nua — day la loi hua tra loi cho chinh chu
// thich cu tren dau file nay: "Ke hoach 6 se chuyen bang nay sang file YAML
// nap nong. Gio de trong code cho don gian, khong dung som he thong nap YAML."
//
// Kieu Lang va hai ham TimLang/DanhSachLang GIU NGUYEN chu ky cu — nhanvat.go,
// khuvuc.go, dongbo.go goi chung KHONG doi mot dong nao. Chi noi DU LIEU o dau
// ra la doi.
//
// # DUNG LAI CO CHE ANH CHUP CUA TASK 3, KHONG TAO SYNC.ONCE THU HAI
//
// khochieu.go (KH8 task 3) da co san mot sync.Once nap kho-noi-dung NHUNG luc
// chua ai goi GanKhoNoiDung (khoChieuMotLan/khoChieuNhung). Anh chup do la MOT
// *noidung.DuLieu duy nhat chua CA chieu LAN lang — nen o day chi can mot ham
// nho khoLangMacDinh() goi lai chinh cai Once ay qua khoNhungMotLan()
// (khochieu.go), roi doc phan Lang cua CUNG mot anh chup. Hai Once cho cung
// mot file se la hai co hoi nap lech nhau — khong lam.
package game

import "nroserver/internal/noidung"

// Lang la mot lang nhan gia (thay cho "chung toc" cua ban Dragon Ball cu).
// Chi so khoi dau la GIA TRI DAU — se can bang lai o Ke hoach 13 (vong can
// bang that, sau khi Bo do can bang Ke hoach 11/12 xong).
// Map va vi tri khoi dau la CHO DAT TRUOC — ban do that dung o Ke hoach 3/6;
// o day chi can mot toa do hop le de nhan vat ton tai.
type Lang struct {
	ID               int16
	Ten              string
	HpKhoiDau        int32
	ChakraKhoiDau    int32
	SatThuongKhoiDau int32 // dung o game/khuvuc.go TinhSatThuong
	MapKhoiDau       int32
	ViTriXKhoiDau    int32
	ViTriYKhoiDau    int32
	VangKhoiDau      int64
}

// khoLangAnhChup la phan cua mot anh chup noi dung ma phan LANG cua goi `game`
// dung — cung khuon voi khoChieuAnhChup (khochieu.go): mot interface nho vi
// *noidung.anhChup la kieu KHONG XUAT KHAU. *noidung.anhChup thoa man no.
type khoLangAnhChup interface {
	Lang(id int) (noidung.Lang, bool)
	DanhSachLang() []noidung.Lang
}

// langTuNoiDung doi noidung.Lang (kieu YAML) sang Lang (kieu goi `game`).
func langTuNoiDung(l noidung.Lang) Lang {
	return Lang{
		ID:               int16(l.Id),
		Ten:              l.Ten,
		HpKhoiDau:        l.HpKhoiDau,
		ChakraKhoiDau:    l.ChakraKhoiDau,
		SatThuongKhoiDau: l.SatThuongKhoiDau,
		MapKhoiDau:       l.MapKhoiDau,
		ViTriXKhoiDau:    l.ViTriXKhoiDau,
		ViTriYKhoiDau:    l.ViTriYKhoiDau,
		VangKhoiDau:      l.VangKhoiDau,
	}
}

// khoLangMacDinh tra ve anh chup kho NHUNG — cung mot *noidung.Kho voi
// khoChieuMacDinh (khochieu.go), qua khoNhungMotLan() dung chung sync.Once.
//
// TimLang/DanhSachLang goi ham nay o duong DANG NHAP / tao nhan vat (nhanvat.go)
// va trong actor (khuvuc.go) — CHUA co kenh nap nong cho lang.yaml (do la Ke
// hoach 10), nen chua co rui ro "hai lan goi trong cung mot nhip roi vao hai
// the he khac nhau" ma hop dong k.acNhip (khochieu.go) dang phong ngua cho
// chieu. Khi task 10 noi day len kenh nap nong, cac diem goi TimLang trong
// khuvuc.go (combat) nen doi sang doc qua k.acNhip nhu chieu — GHI RO o day de
// khong ai quen.
func khoLangMacDinh() khoLangAnhChup {
	return khoNhungMotLan().AnhChup()
}

// TimLang tra ve lang theo id, hoac nil neu khong co.
// Tra ban sao de tranh lo su sua qua con tro cua nguoi goi.
func TimLang(id int16) *Lang {
	if id < 0 {
		return nil
	}
	l, co := khoLangMacDinh().Lang(int(id))
	if !co {
		return nil
	}
	// Tra ban sao, KHONG tro thang vao anh chup goc: nguoi goi lo tay sua
	// truong cua ket qua se khong dung toi du lieu he thong. (Go tu dua
	// ban sao len heap khi lay dia chi tra ve — an toan, dung idiom.)
	cp := langTuNoiDung(l)
	return &cp
}

// DanhSachLang tra ban sao danh sach lang (chi doc), dung thu tu id 0..4.
func DanhSachLang() []Lang {
	ds := khoLangMacDinh().DanhSachLang()
	ra := make([]Lang, len(ds))
	for i, l := range ds {
		ra[i] = langTuNoiDung(l)
	}
	return ra
}
