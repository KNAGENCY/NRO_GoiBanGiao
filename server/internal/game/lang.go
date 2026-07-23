package game

// Lang la mot lang nhan gia (thay cho "chung toc" cua ban Dragon Ball cu).
// Chi so khoi dau la GIA TRI DAU — se can bang lai o Ke hoach 4 (chien dau).
// Map va vi tri khoi dau la CHO DAT TRUOC — ban do that dung o Ke hoach 3/6;
// o day chi can mot toa do hop le de nhan vat ton tai.
type Lang struct {
	ID               int16
	Ten              string
	HpKhoiDau        int32
	ChakraKhoiDau    int32
	SatThuongKhoiDau int32 // dung o Ke hoach 4
	MapKhoiDau       int32
	ViTriXKhoiDau    int32
	ViTriYKhoiDau    int32
	VangKhoiDau      int64
}

// danhSachLang: 5 lang theo thiet ke §7.2. Chi so phan anh tinh cach tung lang
// (Cat chakra cao danh xa, Suong Mu can chien sat thuong, Da mau trau phong thu,
// May toc do chi mang, La can bang). Client cu CHI gui duoc lang 0/1/2 vi cung
// dung mang toc 3x3 — lang 3/4 dung duoc khi client duoc dung lai.
//
// Ke hoach 6 se chuyen bang nay sang file YAML nap nong. Gio de trong code cho
// don gian, khong dung som he thong nap YAML.
var danhSachLang = []Lang{
	{ID: 0, Ten: "Làng Lá", HpKhoiDau: 150, ChakraKhoiDau: 150, SatThuongKhoiDau: 12, MapKhoiDau: 1, ViTriXKhoiDau: 100, ViTriYKhoiDau: 384, VangKhoiDau: 1000},
	{ID: 1, Ten: "Làng Cát", HpKhoiDau: 100, ChakraKhoiDau: 250, SatThuongKhoiDau: 10, MapKhoiDau: 2, ViTriXKhoiDau: 100, ViTriYKhoiDau: 384, VangKhoiDau: 1000},
	{ID: 2, Ten: "Làng Sương Mù", HpKhoiDau: 130, ChakraKhoiDau: 100, SatThuongKhoiDau: 16, MapKhoiDau: 3, ViTriXKhoiDau: 100, ViTriYKhoiDau: 384, VangKhoiDau: 1000},
	{ID: 3, Ten: "Làng Đá", HpKhoiDau: 250, ChakraKhoiDau: 100, SatThuongKhoiDau: 10, MapKhoiDau: 4, ViTriXKhoiDau: 100, ViTriYKhoiDau: 384, VangKhoiDau: 1000},
	{ID: 4, Ten: "Làng Mây", HpKhoiDau: 120, ChakraKhoiDau: 130, SatThuongKhoiDau: 13, MapKhoiDau: 5, ViTriXKhoiDau: 100, ViTriYKhoiDau: 384, VangKhoiDau: 1000},
}

// TimLang tra ve lang theo id, hoac nil neu khong co.
func TimLang(id int16) *Lang {
	for i := range danhSachLang {
		if danhSachLang[i].ID == id {
			return &danhSachLang[i]
		}
	}
	return nil
}

// DanhSachLang tra ban sao danh sach lang (chi doc).
func DanhSachLang() []Lang {
	ds := make([]Lang, len(danhSachLang))
	copy(ds, danhSachLang)
	return ds
}
