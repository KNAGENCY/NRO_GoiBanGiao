// Package noidung — KHO NOI DUNG: nhung con so cua GAME nam trong file YAML
// nguoi thuong doc duoc, thay vi nam trong ma nguon Go.
//
// (KH8 task 0 — CONG TUONG DUONG BYTE. Task nay CHUA doi cho o cua du lieu:
// cac cua tra cuu trong goi `game` van doc du lieu Go cung nhu cu. Task 0 chi
// dung bo nap + CHUNG MINH rang du lieu nap tu YAML sinh ra DUNG TUNG BYTE
// cung mot goi tin nhu du lieu Go hom nay. Chung minh xong moi duoc chinh so —
// xem docs/superpowers/specs/2026-09-03-kh8-kho-noi-dung-design.md §7.1.)
//
// # VI SAO YAML CHU KHONG PHAI JSON
//
// Thu vien chuan Go co san bo doc JSON (0 phu thuoc) — nhung JSON KHONG CO CHU
// THICH. Ma phan dat nhat cua du lieu nay khong phai con so, ma la cau tra loi
// cho "vi sao lai la con so nay". Vi du nguyen van dang nam trong noidungboss.go:
// "Dame boss 35 → don 10..21: nguoi choi 150 mau chiu ~10 don". Chuyen so sang
// JSON la vut bo phan do — va lan sau khong ai dam sua vi khong ai biet no tu dau
// ra. Vi vay chap nhan THEM MOT phu thuoc truc tiep (gopkg.in/yaml.v3 v3.0.1) —
// phu thuoc thu hai cua ca du an, sau github.com/jackc/pgx/v5.
//
// Phuong an "tu viet bo doc YAML toi gian" bi loai: tu viet trinh phan tich la
// tu tao mot nguon loi moi, doi lay viec tiet kiem mot phu thuoc da on dinh.
//
// # RANH GIOI: YAML chua cai TA quyet dinh, Go chua cai CLIENT quyet dinh
//
// Trong YAML: ten, mau, gia, sat thuong, chakra, tam danh — thu chung ta co
// quyen chon va co y muon doi hang ngay.
//
// O LAI Go: ma lenh goi (-28/7, -28/8, -44), THU TU truong trong goi, va cong
// thuc tinh sat thuong. Do la su that cua client Unity DA BIEN DICH; sua duoc
// bang mot file text nghia la mot cu go nham lam client doc lech dong byte — loi
// im lang tuyet doi, khong mot dong log.
//
// # TRAN BYTE — khai bao MOT NOI DUY NHAT o day
//
// Cac hang Tran* duoi day la SUC DOC CUA CLIENT, khong phai lua chon thiet ke.
// YAML BI KIEM theo chung, khong duoc KHAI chung. Truoc KH8 chung nam trong
// game/mauvatpham.go va game/cuahang.go; task 0 chuyen NGUYEN ve day va goi
// `game` tro nguoc lai (khong sao chep) — vi 03/09 hai task song song moi ben tu
// khai mot ban tranParamOption, gop vao main la dut build (§7.2 thiet ke).
//
// Huong import da chot: game -> noidung. TUYET DOI KHONG nguoc lai.
package noidung

const (
	// TranMauOption: client doc SO MAU OPTION bang MOT BYTE KHONG DAU
	// (myReader.readShortOptionTemp voi ModFunc.isShortOptionTemp == false).
	// Vuot tran = client dung catalog RONG, im lang.
	TranMauOption = 255

	// TranMauVatPham: client doc so mau vat pham bang SHORT CO DAU
	// (Controller.cs:6495).
	TranMauVatPham = 32767

	// TranParamOption: gia tri LON NHAT cua param mot option. Client doc param
	// bang ReadParamDQT = readUnsignedShort -> 0..65535; ghi so am hoac lon hon
	// la client hien sai ma khong bao loi.
	TranParamOption = 65535

	// TranThanGoiItem: client goi d.mark(100000) truoc khi doc goi -28/8
	// (Controller.cs:6468). Than dai hon la client doc lech ma KHONG bao loi.
	TranThanGoiItem = 100000

	// TranSoTab / TranSoMonMoiTab: client doc readUnsignedByte (<=255) nhung ban
	// Java ghi writeByte CO DAU. Giu trong 127 de khong bao gio phai doan xem
	// client xu ly 128..255 the nao.
	TranSoTab       = 127
	TranSoMonMoiTab = 127

	// TranSoOptionMoiMon: soOption trong goi -44 la mot byte khong dau.
	TranSoOptionMoiMon = 255
)

// SoLangCoDinh: dung 5 lang Naruto. Client cu chi gui duoc lang 0/1/2 nen
// KHONG duoc them lang thu 6 bang YAML (thiet ke §2.2); bo kiem chan viec do.
const SoLangCoDinh = 5

// Ba KIEU CHIEU ma may chu DA CAI DAT. YAML chi duoc chon trong tap nay: them
// mot chieu CUNG KIEU la them du lieu (re); them mot KIEU MOI (choang, day lui,
// doc, khien, hut mau) la VIET CODE MAY CHU MOI. Khai mot kieu la trong YAML
// phai lam bo kiem DO, khong duoc im lang bo qua (thiet ke §1.2).
const (
	LoaiChieuDanh   byte = 1
	LoaiChieuBuff   byte = 2
	LoaiChieuTuThan byte = 3
)

// ---------------------------------------------------------------------------
// Kieu du lieu YAML — chieu thuc
// ---------------------------------------------------------------------------

// CapChieu la MOT cap cua mot chieu. Moi truong tuong ung mot o trong khuon
// byte -28/7 (game/chieuthuc.go, ghiMauChieu) — nhung THU TU byte do goi `game`
// quyet dinh, khong phai file YAML.
type CapChieu struct {
	SkillId      int16  `yaml:"skill_id"`
	Diem         byte   `yaml:"diem"`
	NhanLucCan   int64  `yaml:"nhan_luc_can"`
	ChakraDung   int16  `yaml:"chakra_dung"`
	HoiChieuMs   int32  `yaml:"hoi_chieu_ms"`
	Dx           int16  `yaml:"dx"`
	Dy           int16  `yaml:"dy"`
	MaxFight     byte   `yaml:"max_fight"`
	SatThuong    int16  `yaml:"sat_thuong"`
	Gia          int16  `yaml:"gia"`
	ThemThongTin string `yaml:"them_thong_tin"`
}

// Chieu la mot MauChieu (SkillTemplate) voi toan bo cap cua no.
type Chieu struct {
	TemplateId    byte       `yaml:"template_id"`
	Ten           string     `yaml:"ten"`
	CapToiDa      byte       `yaml:"cap_toi_da"`
	KieuTonChakra byte       `yaml:"kieu_ton_chakra"`
	Loai          byte       `yaml:"loai"`
	IconId        int16      `yaml:"icon_id"`
	MoTa          string     `yaml:"mo_ta"`
	CacCap        []CapChieu `yaml:"cac_cap"`
}

// ChieuCuaLang la khoi chieu cua MOT lang. Truoc KH8 task 6 file nay con giu
// mot ban sao thu hai cua ten lang (`ten_lang`, chi de nguoi doc dinh huong,
// bo nap khong dung) — da XOA: ten lang gui xuong client (goi -28/7) gio doc
// TU MOT NGUON DUY NHAT la lang.yaml (qua game/lang.go), tra theo id_lang.
type ChieuCuaLang struct {
	IdLang byte    `yaml:"id_lang"`
	Chieu  []Chieu `yaml:"chieu"`
}

type fileChieu struct {
	Lang []ChieuCuaLang `yaml:"lang"`
}

// ---------------------------------------------------------------------------
// Kieu du lieu YAML — mau option + mau vat pham
// ---------------------------------------------------------------------------

// MauOption la MOT dong bang mau option. `id` = CHI SO, KHAI TUONG MINH.
type MauOption struct {
	Id   int    `yaml:"id"`
	Ten  string `yaml:"ten"`
	Loai int8   `yaml:"loai"`
}

// MauVatPham la MOT mau vat pham. `id` = CHI SO, KHAI TUONG MINH.
type MauVatPham struct {
	Id          int    `yaml:"id"`
	Loai        int8   `yaml:"loai"`
	GioiTinh    int8   `yaml:"gioi_tinh"`
	Ten         string `yaml:"ten"`
	MoTa        string `yaml:"mo_ta"`
	Cap         int8   `yaml:"cap"`
	YeuCauSuc   int32  `yaml:"yeu_cau_suc"`
	IconId      int16  `yaml:"icon_id"`
	Phan        int16  `yaml:"phan"`
	NangCapDuoc bool   `yaml:"nang_cap_duoc"`
	TranChong   int16  `yaml:"tran_chong"`
}

type fileVatPham struct {
	MauOption  []MauOption  `yaml:"mau_option"`
	MauVatPham []MauVatPham `yaml:"mau_vat_pham"`
}

// ---------------------------------------------------------------------------
// Kieu du lieu YAML — cua hang
// ---------------------------------------------------------------------------

// OptionMon la mot option gan san vao mon hang bay ban.
type OptionMon struct {
	OptionId uint8  `yaml:"option_id"`
	Param    uint16 `yaml:"param"`
}

// MonHang la mot mon bay ban. `tien` la chuoi "vang" hoac "ngoc" — CO Y dung
// chu chu khong dung so 0/1: mot con so go nham van tra ra mot loai tien hop le
// va chay em; mot chuoi go nham thi bo kiem BAT DUOC ngay.
//
// GiaBanLai (KH8 task 5, thiet ke §5.4 diem 2) la gia NGUOI CHOI BAN LAI mon nay
// cho NPC. Truoc task 5 no khong ton tai: duong ban tinh `Gia/4` trong code Go —
// nghia la muon chinh gia thu lai cua MOT mon thi phai sua cong thuc cho MOI mon.
// Nay no la mot con so TUONG MINH, sua duoc tung mon, nap nong duoc.
type MonHang struct {
	MauId int16  `yaml:"mau_id"`
	Gia   int32  `yaml:"gia"`
	Tien  string `yaml:"tien"`
	// GiaBanLai: 0 <= GiaBanLai <= Gia. Bo kiem chan ca hai dau — ban lai DAT HON
	// gia mua la MAY IN TIEN (mua 100 ban 120, lap vo han), khong phai lua chon
	// can bang.
	GiaBanLai int32       `yaml:"gia_ban_lai"`
	Options   []OptionMon `yaml:"options"`
	LaMoi     bool        `yaml:"la_moi"`
}

// TabCuaHang la mot the trong cua so cua hang.
type TabCuaHang struct {
	Ten string    `yaml:"ten"`
	Mon []MonHang `yaml:"mon"`
}

// CuaHang la mot cua hang NPC hoan chinh. ID KHONG di xuong client, nhung may
// chu nho "dang mo cua hang nao" theo id -> doi id giua phien la ban nham bang gia.
type CuaHang struct {
	Id   int32        `yaml:"id"`
	Ten  string       `yaml:"ten"`
	Tabs []TabCuaHang `yaml:"tabs"`
}

type fileCuaHang struct {
	CuaHang []CuaHang `yaml:"cua_hang"`
}

// Hai gia tri hop le cua truong `tien`.
const (
	TienVang = "vang"
	TienNgoc = "ngoc"
)

// ---------------------------------------------------------------------------
// Kieu du lieu YAML — quai (KH8 task 8)
// ---------------------------------------------------------------------------

// MauQuai la MOT mau quai. CAN CUOC la `Ma` (chuoi), KHONG phai so — khac han
// MauVatPham o tren. Ly do: id mau quai KHONG di xuong client va KHONG nam
// trong CSDL nguoi choi, no chi song trong bo nho may chu luc SINH quai. Vi vay
// them/bo/sap xep lai mau quai la an toan, mien la `sinh:` van tro dung ma.
//
// TranChong/Phan/... khong co o day vi quai khong phai do vat. Cai duy nhat di
// xuong client trong goi -24 la Hp/HpToiDa/Cap/CapTinhAnh/X/Y — khuon byte do
// goi `game` so huu (ghiQuai, game/quai.go), KHONG phai file YAML.
type MauQuai struct {
	Ma  string `yaml:"ma"`
	Ten string `yaml:"ten"`

	// AnhId la id ANH quai trong catalog da ky (assets/catalog.json, ghim SHA
	// trong asset-lock.json). Bo nap o goi nay KHONG tu kiem duoc no ton tai —
	// `noidung` khong import `game`. Nguoi TIEM du lieu goi KiemIdAnhQuai voi
	// danh sach id hop le lay tu catalog (game/khoquai.go).
	AnhId int `yaml:"anh_id"`

	Hp int32 `yaml:"hp"`

	// SatThuongCoBan: sat thuong CO DINH quai gay ra moi don len nguoi choi.
	// Truoc KH8 task 8 day la hang so 5 nam trong game/quai.go; tach ra tung
	// mau chinh la de sau nay DO CAN BANG duoc (thiet ke §5.3).
	SatThuongCoBan int32 `yaml:"sat_thuong_co_ban"`

	Cap        byte `yaml:"cap"`
	CapTinhAnh byte `yaml:"cap_tinh_anh"`
}

// DiemSinhQuai la MOT con quai duoc dat san tren mot ban do: mau nao, tempId
// bao nhieu, dung o dau.
type DiemSinhQuai struct {
	Ma string `yaml:"ma"`
	Id byte   `yaml:"id"`
	X  int16  `yaml:"x"`
	Y  int16  `yaml:"y"`
}

// SinhQuaiTrenMap la toan bo quai duoc dat san tren MOT ban do.
type SinhQuaiTrenMap struct {
	Map int32          `yaml:"map"`
	Dat []DiemSinhQuai `yaml:"dat"`
}

type fileQuai struct {
	MauQuai []MauQuai         `yaml:"mau_quai"`
	Sinh    []SinhQuaiTrenMap `yaml:"sinh"`
}

// TranTempIdQuai: client doc tempId (va soQuai) bang BYTE CO DAU (dac ta §1.4),
// nen moi id quai phai <= 127. Day la SUC DOC CUA CLIENT, khong phai lua chon.
const TranTempIdQuai = 127

// ---------------------------------------------------------------------------
// Kieu du lieu YAML — lang (KH8 task 6)
// ---------------------------------------------------------------------------

// Lang la mot lang nhan gia (thay cho "chung toc" cua ban Dragon Ball cu). Chi
// so khoi dau la GIA TRI DAU — se can bang lai o Ke hoach 13 (vong can bang
// that). `id` = CHI SO, KHAI TUONG MINH (nhu MauVatPham/MauOption) vi id nay
// DANG NAM TRONG CSDL nguoi choi (cot nhan_vat.lang); client cu CHI GUI DUOC
// lang 0/1/2 (con dung mang toc 3x3) nen KHONG duoc them lang thu 6 bang YAML
// (thiet ke §2.2, bang "Lang").
//
// MapKhoiDau/ViTriXKhoiDau/ViTriYKhoiDau la CHO DAT TRUOC — ban do that dung o
// Ke hoach 3/6 cua lo trinh game; o day chi can mot toa do hop le de nhan vat
// ton tai luc tao moi.
type Lang struct {
	Id               int    `yaml:"id"`
	Ten              string `yaml:"ten"`
	HpKhoiDau        int32  `yaml:"hp_khoi_dau"`
	ChakraKhoiDau    int32  `yaml:"chakra_khoi_dau"`
	SatThuongKhoiDau int32  `yaml:"sat_thuong_khoi_dau"` // dung o game/khuvuc.go TinhSatThuong
	MapKhoiDau       int32  `yaml:"map_khoi_dau"`
	ViTriXKhoiDau    int32  `yaml:"vi_tri_x_khoi_dau"`
	ViTriYKhoiDau    int32  `yaml:"vi_tri_y_khoi_dau"`
	VangKhoiDau      int64  `yaml:"vang_khoi_dau"`
}

type fileLang struct {
	Lang []Lang `yaml:"lang"`
}

// ---------------------------------------------------------------------------
// Anh chup
// ---------------------------------------------------------------------------

// DuLieu la MOT ANH CHUP hoan chinh cua kho noi dung. Task 0 chi dung va tra ve
// no; viec cong bo bang atomic.Pointer + so the he la task 2.
type DuLieu struct {
	// ChieuTheoLang[i] = danh sach chieu cua lang i (0..4). Mang co dinh 5 phan
	// tu de khong mot lan tra cuu nao phai lo "lang nay co trong map khong".
	ChieuTheoLang [SoLangCoDinh][]Chieu
	// Lang[i] = hang so khoi dau cua lang i (0..4). Cung khuon mang co dinh voi
	// ChieuTheoLang (KH8 task 6).
	Lang       [SoLangCoDinh]Lang
	MauOption  []MauOption
	MauVatPham []MauVatPham
	CuaHang    []CuaHang

	// MauQuai / SinhQuai — KH8 task 8. Tra cuu qua cua Quai()/DanhSachQuai()/
	// SinhQuaiCuaMap() cua anh chup (kho.go), khong doc thang.
	MauQuai  []MauQuai
	SinhQuai []SinhQuaiTrenMap

	// BangRoi — KH8 task 9 (bangroi.go). Tra cuu qua BangRoiCuaQuai(ma) /
	// BangRoiCuaBoss(ma) cua anh chup, khong doc thang.
	BangRoi BangRoi
	// Boss — KH8 task 7. Kieu, bo kiem VA hai cua tra cuu Boss()/DanhSachBoss()
	// deu nam trong boss.go (mot vung du lieu = mot tep, xem chu thich dau tep do).
	Boss []Boss
}
