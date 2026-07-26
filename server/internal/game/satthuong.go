// satthuong.go — cong thuc sat thuong + he nguyen to ngu hanh cho Ke hoach 6,
// Task 4. Dac ta §4.1-4.2 (docs/superpowers/specs/2026-07-26-giao-thuc-chien-dau.md).
//
// Cong thuc (§4.1):
//
//	satThuong = satThuongGoc(nhanVat) × heSoChieu(capChieu) × heSoChiMang × heSoHeNguyenTo
//
// satThuongGoc tinh tu chi so nhan vat, truyen vao qua tham so goc (KHONG
// phai viec cua file nay). heSoChieu doc tu du lieu chieu: truong
// CapChieu.SatThuong duoc dung nhu PHAN TRAM cua sat thuong goc (100 = he so
// 1.0), khong phai cong don — day la "damage trong catalog" ma dac ta noi
// toi, moi thu dac thu tung chieu nam trong DU LIEU (chuan bi cho KH8), code
// khong biet ten chieu nao ca. heSoChiMang la 2x neu chi mang, nguoc lai 1x.
// heSoHeNguyenTo xem ham HeSoNguyenTo ben duoi.
package game

// Nam hang so he nguyen to — gia tri chinh la CHI SO trong VONG KHAC CHE
// (khong phai ID lang). Dac y: thu tu khai bao TRUNG voi vong khac che §4.2
// (Hoa → Phong → Loi → Tho → Thuy → Hoa) nen phep tinh "he i khac he i+1"
// dung duoc luon iota, khong can bang tra rieng.
const (
	HeHoa byte = iota
	HePhong
	HeLoi
	HeTho
	HeThuy
	soHeNguyenTo = 5 // tong so he, dung de tinh modulo vong khac che
)

// heSoNguyenToPhanTram tra he so he nguyen te DANG PHAN TRAM (125/100/80) —
// giu nguyen so nguyen de TinhSatThuong tinh toan xuyen suot bang int64,
// khong dinh dau phay dong o buoc lo`i (tranh sai so lam tron/kho doan truoc
// gan bien tran int32).
func heSoNguyenToPhanTram(heNguoiDanh, heMucTieu byte) int64 {
	// Neu byte truyen vao ngoai [0,4] (du lieu hong), coi nhu khong xac dinh
	// duoc quan he khac che — tra ve trung tinh thay vi panic, an toan hon
	// cho mot ham nam giua vong doi chien dau.
	if heNguoiDanh >= soHeNguyenTo || heMucTieu >= soHeNguyenTo {
		return 100
	}
	// nguoiDanh khac mucTieu <=> mucTieu dung ngay sau nguoiDanh trong vong.
	if heMucTieu == (heNguoiDanh+1)%soHeNguyenTo {
		return 125
	}
	// nguoiDanh bi khac <=> nguoiDanh dung ngay sau mucTieu trong vong.
	if heNguoiDanh == (heMucTieu+1)%soHeNguyenTo {
		return 80
	}
	return 100
}

// HeSoNguyenTo tra he so nhan sat thuong theo quan he khac che ngu hanh giua
// he nguoi danh va he muc tieu (dac ta §4.2). Vong khac che:
//
//	Hoa → Phong → Loi → Tho → Thuy → Hoa
//
// (doc: Hoa bi Thuy dap, Thuy bi Tho hut, Tho bi Loi xuyen, Loi bi Phong tan,
// Phong bi Hoa dot). Khac che ×1.25 · trung tinh ×1.00 · bi khac ×0.80.
//
// Tra ve float64 de dung lam gia tri hien thi/log — TinhSatThuong ben duoi
// KHONG goi ham nay, no dung heSoNguyenToPhanTram (so nguyen) de tranh dau
// phay dong trong duong tinh sat thuong loi.
func HeSoNguyenTo(heNguoiDanh, heMucTieu byte) float64 {
	return float64(heSoNguyenToPhanTram(heNguoiDanh, heMucTieu)) / 100.0
}

// Kep tran/san sat thuong — int32, dung nhu khuon giao thuc cua truong
// SatThuong trong cac goi chien dau.
const (
	satThuongToiDa    int64 = 2_147_483_647 // math.MaxInt32, ghim so de test doc thang khong can import math
	satThuongToiThieu int64 = 0
)

// TinhSatThuong tinh sat thuong cuoi cung theo cong thuc §4.1, bang int64
// XUYEN SUOT roi moi kep xuong int32 O BUOC CUOI CUNG — tai sao: tinh thang
// tren int32 thi phep nhan (vd sat thuong goc lon x he so 125%) TRAN truoc
// khi kip kep, ra ket qua sai (co the am) ma khong bao gio thay duoc gia tri
// dung 2_147_483_647.
//
// chieu la con tro CapChieu (Task 1, chieuthuc.go) — ham nay CHI doc truong
// SatThuong (du lieu), KHONG biet chieu ten gi, khong switch theo ten chieu.
func TinhSatThuong(goc int64, chieu *CapChieu, heNguoiDanh, heMucTieu byte, chiMang bool) int32 {
	// heSoChieu: SatThuong cua cap chieu la PHAN TRAM cua sat thuong goc.
	ketQua := goc * int64(chieu.SatThuong) / 100

	// heSoHeNguyenTo: nhan roi chia 100, van la so nguyen.
	ketQua = ketQua * heSoNguyenToPhanTram(heNguoiDanh, heMucTieu) / 100

	// heSoChiMang: gap doi neu chi mang.
	if chiMang {
		ketQua *= 2
	}

	// Kep san truoc — sat thuong khong bao gio am (vi du he so < 1 lam tron
	// ve am neu goc am, du hiem khi xay ra van phai chan).
	if ketQua < satThuongToiThieu {
		ketQua = satThuongToiThieu
	}
	// Kep tran — day la buoc BAT BUOC phai lam TRUOC khi convert sang int32,
	// vi convert truc tiep mot gia tri int64 vuot pham vi int32 sang int32 la
	// hanh vi CAT BIT (wrap), co the ra so am hoac so nho vo nghia.
	if ketQua > satThuongToiDa {
		ketQua = satThuongToiDa
	}

	return int32(ketQua)
}
