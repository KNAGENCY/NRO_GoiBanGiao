package store

// vatpham.go — KH7 Task 1: kieu du lieu MOT MON vat pham va cac hang so noi chua.
//
// Thiet ke: docs/superpowers/specs/2026-09-03-kh7-vat-pham-design.md §1.2/§5.1.
// Task 1 CHI dinh nghia kieu + hang so + rang buoc CSDL (migration 0005). Cac ham
// doc/ghi that (LayHanhTrang, CapVatPham, ChuyenO, mua/ban...) den o cac task sau —
// va TAT CA phai di qua giao dich CSDL cong/tru tai cho, KHONG qua ThoGhi (§5.3).

// Noi chua cua mot mon. Trung TUNG CHU voi CHECK `vat_pham_noi_chua_hop_le` trong
// migration 0005 — sua mot ben phai sua ben kia, neu khong mot gia tri hop le o Go
// se bi CSDL tu choi (hoac nguoc lai). Dung hang so thay vi chuoi tran de go nham
// mot chu la loi bien dich chu khong phai loi chay ngam.
const (
	NoiChuaTui     = "TUI"      // trong tui do (co chu, co o)
	NoiChuaMac     = "MAC"      // dang mac tren nguoi (co chu, co o)
	NoiChuaRuong   = "RUONG"    // trong ruong/kho (co chu, co o)
	NoiChuaDuoiDat = "DUOI_DAT" // roi duoi dat cho nhat (khong chu, khong o)
	NoiChuaKyGui   = "KY_GUI"   // ky gui cho ban (co chu, co o) — dung o ban 3
	NoiChuaDaHuy   = "DA_HUY"   // da tieu huy (khong chu) — dong GIU LAI, khong xoa
)

// TuyChon la MOT chi so phu (option) tren mot mon: khuon JSONB `{"id","param"}`
// (§1.2). id ∈ [0,255] (client doc 1 byte), Param ∈ [0,65535] (client doc unsigned
// short — ghi so am la hong). Rang buoc mien duoc ep o tang Go khi dung/ghi (task
// sau), CSDL chi giu khuon rong bang DEFAULT `{"opt":[]}` + CHECK don-chiec.
type TuyChon struct {
	Id    int `json:"id"`
	Param int `json:"param"`
}

// ChiSoPhu la ca goi option cua mot mon — khuon JSONB chot cung `{"opt":[...]}`.
// Rong (Opt nil/[]) nghia la mon khong option, duoc phep cong don (so_luong > 1).
type ChiSoPhu struct {
	Opt []TuyChon `json:"opt"`
}

// VatPham la MOT MON (instance) — mot dong bang `vat_pham`. Con tro cho hai truong
// co the NULL trong CSDL (chu so huu va o so): mon duoi dat / da huy khong co chu
// va khong nam trong o nao, phan biet duoc voi "chu = 0" hay "o = 0" (deu la gia
// tri that). Cac truong so khong-NULL de gia tri thang.
type VatPham struct {
	ID         int64    // id ben toan cuc, KHONG BAO GIO tai su dung
	MauID      int32    // ItemTemplate.id — quyet dinh ten/icon/loai
	ChuSoHuuID *int64   // nil khi duoi dat / da huy
	NoiChua    string   // mot trong cac hang so NoiChua* o tren
	OSo        *int16   // nil khi khong nam trong o nao (duoi dat / da huy)
	SoLuong    int32    // >= 1; mon co option thi luon = 1 (§1.3)
	ChiSoPhu   ChiSoPhu // option, khuon JSONB `{"opt":[...]}`
}
