package game

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"nroserver/internal/antoan"
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// KhuVuc (zone) la don vi khong gian nho nhat cua mot ban do: giu ban do dang
// choi + danh sach nguoi choi dang o trong do. Ke hoach 4 CHI can mot nguoi dung
// yen trong map trong, nen day la mo hinh TOI THIEU — khong goroutine rieng,
// khong mutex, khong nhip 100ms.
//
// VI SAO CHUA DUNG MO HINH ACTOR (goroutine moi zone) NGAY: actor chi can khi co
// TRANH CHAP that — nhieu luong cung doc/ghi trang thai zone (di chuyen, nhin thay
// nhau, quai di lai). Ke hoach 4 chua co thu nao trong so do: mot nguoi, dung yen,
// khong broadcast. Dung actor bay gio la giai bai toan chua ton tai (YAGNI). Ke
// hoach 5 (di chuyen) se nang KhuVuc thanh actor khi co tranh chap that su, kem
// nhip 100ms cho quai/vat pham (xem 2026-07-23-hoa-giai-muot §4.3). O KH4, KhuVuc
// chi lo NAP ban do va lam nen cho KH5.
//
// KE HOACH 5: KhuVuc nay tro thanh ACTOR — mot goroutine rieng so huu toan bo
// trang thai (nguoiTrongKhuVuc), nhan su kien qua hop thu (channel) va co nhip
// tim 100ms cho cac cong viec dinh ky (quai/vat pham o KH6). Moi thay doi trang
// thai PHAI di qua hop thu — khong truy cap truc tiep field tu goroutine khac.
type KhuVuc struct {
	banDo     *BanDo
	nguoiChoi []*store.NhanVat

	nguoiTrong map[int64]*nguoiTrongKhuVuc
	// anhChupDaRoi giu tam ban cuoi cua session bi actor loai vi outbound
	// error/queue day. KhiDong van xin duoc ban nay de xep luu, roi control
	// leave don no; id khong nam trong nguoiTrong nen khong tao ghost.
	anhChupDaRoi map[int64]store.NhanVat

	// quai la danh sach quai dang song trong khu vuc nay (Ke hoach 6 Task 5).
	// Tao MOT LAN trong MoKhuVuc; sau do CHI luong vongLap duoc doc/ghi (dac ta
	// §1.6: "chi luong vongLap duoc cham trang thai quai va trang thai chien
	// dau"). KHONG co getter doc truc tiep tu ben ngoai — luong khac (vd
	// luong phien khi dung goi -24) phai gui suKienXinDanhSachQuai qua hop thu
	// va nhan BAN SAO qua kenh traLoi (xem suKienXinDanhSachQuai ben duoi).
	quai []*Quai

	// boss la danh sach boss dang song trong khu vuc nay (KH6b Task 7). CUNG
	// QUY TAC SO HUU voi k.quai ngay tren: CHI luong vongLap duoc doc/ghi —
	// khong getter doc thang tu ben ngoai, khong con tro nao thoat ra ngoai
	// actor. Dat CANH k.quai co y: hai danh sach nay la toan bo "thu do may
	// dieu khien" ma actor so huu.
	//
	// KHAC k.quai o mot diem: quai duoc sinh MOT LAN trong MoKhuVuc
	// (taoQuaiMacDinh), con boss thi KHONG — MoKhuVuc de danh sach nay RONG
	// (nil slice, append van chay). Boss chi vao khu vuc qua themBoss, do bo
	// dieu phoi lich NGOAI actor goi (CHAN-3/CHAN-7 cua ke hoach: MoKhuVuc va
	// MoSoKhuVuc GIU NGUYEN hanh vi, khong tu khoi dong hay sinh gi moi — bo
	// test tich hop dung "5 khu vuc im lang" lam moc).
	boss []*Boss

	// dongHo la DONG HO TIEM DUOC cho toan bo logic co moc thoi gian trong
	// actor (Ke hoach 6 Task 8: hoi sinh 5s + hoi mau 30s trong nhipTim, va
	// kiem bat tu trong case suKienDanhQuai) — RANG BUOC CHUNG cua brief Task
	// 8: "khong goi time.Now() thang trong logic" de test moi moc thoi gian
	// (4.9s/5.0s/29.9s/30.0s...) MA KHONG sleep that (khong lam bo test cham
	// them hang chuc giay). MoKhuVuc gan mac dinh time.Now (hanh vi that khi
	// chay server that); test CUNG package gan thang truong nay (khong co
	// setter rieng — dung tinh than "chi test cung package moi duoc dung truc
	// tiep trang thai noi bo" da ap dung cho k.quai/k.nguoiTrong o cac Task
	// truoc) de dieu khien "bay gio" theo y muon, khong sleep.
	dongHo func() time.Time

	// loa la duong DUY NHAT actor duoc dung de bao mot tin ra TOAN MAY CHU
	// (KH6b Task 6, CHAN-2). nil = chua noi day -> xepThongBao im lang, KHONG
	// panic: rat nhieu test (va ca giai doan truoc khi main.go goi GanLoa)
	// dung MoKhuVuc truc tiep, khong he co loa.
	//
	// VI SAO atomic.Pointer chu khong phai *BoPhatThanh thuong: SoKhuVuc.GanLoa
	// gan loa SAU khi MoSoKhuVuc da BatDau() 5 actor (luat 4: MoSoKhuVuc giu
	// nguyen, khong tu khoi dong them gi) — nen loi ghi nay den tu goroutine
	// main TRONG LUC cac actor dang chay va co the doc truong nay tu vongLap.
	// Con tro thuong o day la data race THAT (go test -race bat duoc ngay khi
	// Task 12 them loi doc trong actor); atomic.Pointer co zero value san sang
	// dung (nil) nen MoKhuVuc khong can khoi tao gi them.
	loa atomic.Pointer[BoPhatThanh]

	hopThu          chan suKienKhuVuc
	hopThuDieuKhien chan suKienDieuKhien
	xong            chan struct{}
	dungMotLan      sync.Once
	loiDung         error

	// soLanGuiRoi dem gameplay event best-effort bi bo khi hop thu day. Snapshot,
	// leave va stop nay di qua hopThuDieuKhien reliable nen khong con nam trong
	// nhom co the roi nay.
	//
	// atomic.Int64 (KHONG PHAI int64 thuong) vi Gui() duoc goi TU NHIEU
	// goroutine cung luc: moi luong phien (mot goroutine/ket noi) + luong
	// dinh ky auto-save (§4 diem A) + luong xu ly tat may chu (§4 diem C) deu
	// co the goi Gui() DONG THOI — bien thuong se la data race (go vet -race
	// se bat), atomic.Int64 la kieu zero-value-san-sang, an toan doc/ghi tu
	// nhieu goroutine ma khong can mutex rieng.
	soLanGuiRoi atomic.Int64

	// soPanicLienTiep dem so NHIP GAMEPLAY (xuLy/nhipTim) panic LIEN TIEP GAN
	// NHAT — van an toan cua manh A §3.2a. CHI luong vongLap doc/ghi (dung ky
	// luat "mot chu so huu" nhu moi truong khac cua actor) nen la int thuong,
	// khong can atomic. Reset ve 0 ngay khi co MOT nhip trot lot -> loi THOANG
	// QUA (dung y "khu vuc cu chay tiep") khong bao gio cham nguong; chi loi
	// DAI DANG THAT (panic khong dut) moi tich luy toi nguong roi dung rieng
	// khu vuc.
	soPanicLienTiep int

	// Do tre nhip (KH6c): mocNhipTruoc chi vongLap cham (mot luong) nen la bien
	// thuong; soNhipTre/treNhipLauNhatMs la atomic vi test/giam sat doc song song.
	mocNhipTruoc     time.Time
	soNhipTre        atomic.Int64
	treNhipLauNhatMs atomic.Int64
}

// nguongPanicDungKhuVuc: so nhip panic LIEN TIEP toi da truoc khi van an toan
// dung rieng khu vuc (manh A §3.2a: "100 nhip = 10 giay panic lien khong dut" o
// nhip tim 100ms). Dat kieu hang so o mot noi de doi de + de kiem thu dot bien.
const nguongPanicDungKhuVuc = 100

// guiTin la phan toi thieu cua mot phien (session) ma actor can de gui goi tin
// cho client. *mynet.Session thoa man interface nay; test dung phien gia.
type guiTin interface {
	XepGui(*protocol.Message) error
	DongCham()
}

// nguoiTrongKhuVuc gan mot nhan vat voi phien mang cua no trong pham vi mot
// khu vuc — de actor biet gui goi tin cho ai khi co su kien.
type nguoiTrongKhuVuc struct {
	nv    *store.NhanVat
	phien guiTin

	// chieuDangChon la CapChieu (cap 1) cua chieu nguoi choi vua chon bang goi
	// 34 (xem suKienChonChieu). nil = chua chon chieu nao. CHI luong vongLap
	// cua actor duoc doc/ghi truong nay — day chinh la trang thai dac ta §3.2
	// yeu cau dat trong actor (khong trong BoXuLy) de tranh tranh chap voi
	// luong tinh sat thuong doc no sau nay (Ke hoach 6 Task 3+).
	chieuDangChon *CapChieu

	// loaiChieuDangChon la truong "Loai" (type: 1 danh, 2 buff/hoi mau, 3 tu
	// than) cua MauChieu CHA chua chieuDangChon (Fix wave cuoi Task 11, 27/07)
	// — dat CUNG LUC voi chieuDangChon trong case suKienChonChieu duoi day.
	// LY DO CAN TRUONG RIENG NAY (khong doc thang tu chieuDangChon): CapChieu
	// (chieuthuc.go) la mot CAP CU THE cua chieu, KHONG mang truong Loai —
	// Loai chi ton tai o MauChieu (SkillTemplate) cha, va chieuDangChon la con
	// tro *CapChieu, khong con duong nao doc nguoc ve MauChieu cha tu con tro
	// do. Zero value (0) KHONG trung loaiDanh/loaiBuff/loaiTuThan (1/2/3, xem
	// chieuthuc.go) nen "chua chon chieu nao" tu nhien khong bi hieu nham la
	// bat ky loai nao — an toan voi moi noi dang doc chieuDangChon==nil truoc
	// Task 11 (case suKienDanhNguoiChoi la noi DAU TIEN dung truong nay).
	loaiChieuDangChon byte

	// lanDanhCuoi la thoi diem nguoi choi nay DUNG THANH CONG mot chieu lan
	// gan nhat — Task 3 (chotkiem.go) can truong nay de tinh "da troi qua bao
	// lau" cho kiem hoi chieu. store.NhanVat KHONG co truong thoi gian nao
	// (da kiem: chi co cac truong chi so), nen buoc phai dat o day — dung vi
	// tri voi chieuDangChon (chi actor duoc doc/ghi). Zero value (chua tung
	// dung chieu nao) khien bayGio.Sub(lanDanhCuoi) ra mot khoang RAT LON, nen
	// kiem hoi chieu tu nhien CHO PHEP ngay lan dau — khong can if/else rieng.
	lanDanhCuoi time.Time

	// lucHoiSinh la thoi diem nguoi choi nay HOI SINH GAN NHAT (Ke hoach 6
	// Task 13, dac ta §1.7/§3.4/§6.4) — dat trong case suKienHoiSinh duoi
	// day. CUNG CO CHE voi Quai.LucHoiSinh (quai.go, Task 8): CHI luong
	// vongLap cua actor duoc doc/ghi (dac ta §1.6), va dung lam nguon cho
	// ViTriCoMau.LucHoiSinh (chotkiem.go) khi Task 11 (PvP) can kiem bat tu
	// 1500ms cho MUC TIEU la nguoi choi — nguoiTrongKhuVuc gio da du bo
	// truong (X/Y qua nv.ViTriX/Y, HpHienTai qua nv.HpHienTai, LucHoiSinh qua
	// truong nay) de dung lam nguon ViTriCoMau, dung y brief Task 13. Zero
	// value (nguoi choi CHUA TUNG hoi sinh lan nao tu luc dang nhap) khien
	// bayGio.Sub(zero) ra khoang thoi gian RAT LON -> kiem bat tu tu nhien
	// CHO PHEP ngay tu dau, dung hanh vi mong doi — GIONG HET ly do zero
	// value cua Quai.LucHoiSinh (quai.go) an toan.
	lucHoiSinh time.Time
}

// benDanhTuNguoi doi mot nguoi choi trong khu vuc thanh BenDanh (chotkiem.go)
// — adapter cua Ke hoach 6b Task 5. Truoc do KiemDanhHopLe nhan thang
// *nguoiTrongKhuVuc; tach ra de boss (kieu rieng, khong co store.NhanVat) di
// CHUNG mot duong kiem voi nguoi choi thay vi co ban sao logic thu hai.
//
// VoHanChakra LUON false o day: mien chakra la dac quyen rieng cua boss, nguoi
// choi khong bao gio duoc cho qua kiem #2 (TestBenDanhTuNguoiChepDungTruong
// khoa dieu nay).
//
// nguoi == nil (hoac chua co nv) tra ve BenDanh zero — Chieu == nil nen
// KiemDanhHopLe tu choi voi ly do "chua chon chieu", GIU NGUYEN hanh vi cua
// nhanh bao ve `nguoiDanh == nil` truoc khi tach ham.
func benDanhTuNguoi(nguoi *nguoiTrongKhuVuc) BenDanh {
	if nguoi == nil || nguoi.nv == nil {
		return BenDanh{}
	}
	return BenDanh{
		X:             nguoi.nv.ViTriX, // int32, KHONG ep xuong int16 (xem BenDanh.X)
		Y:             nguoi.nv.ViTriY,
		ChakraHienTai: nguoi.nv.ChakraHienTai,
		LanDanhCuoi:   nguoi.lanDanhCuoi,
		Chieu:         nguoi.chieuDangChon,
		VoHanChakra:   false, // xem chu thich tren — khong bao gio doi thanh true cho nguoi choi
	}
}

// suKienKhuVuc la kieu su kien actor xu ly qua hop thu. Ke hoach 5 Task 5 se
// lam day logic xu ly; Task 1 chi dinh nghia khung.
type suKienKhuVuc interface{}

// suKienVao bao actor mot nhan vat vua vao khu vuc, kem phien de gui goi tin.
type suKienVao struct {
	nv    *store.NhanVat
	phien guiTin
}

// suKienDiChuyen bao actor mot nhan vat doi vi tri.
type suKienDiChuyen struct {
	idNhanVat int64
	x, y      int16
}

// suKienRa bao actor mot nhan vat vua roi khu vuc.
type suKienRa struct {
	idNhanVat int64
}

// suKienChonChieu bao actor mot nhan vat vua chon mot chieu (goi len `34`,
// dac ta §3.2). mauChieuId la template.id (KHONG phai skillId — hai truong
// khac nhau, xem bang §3.2 cua dac ta). Luong phien CHI gui su kien nay, actor
// moi la noi tra cuu ChieuCuaLang + kiem chieu co thuoc lang nhan vat khong
// va ghi chieuDangChon — tranh dung bay §6.9 (tranh chap doc/ghi nv.ViTriX/Y
// ma Ke hoach 5 da vap).
type suKienChonChieu struct {
	idNhanVat  int64
	mauChieuId byte
}

// suKienDanhQuai bao actor mot nhan vat vua gui goi len `54` (dac ta §3.1):
// xin danh con quai co Id = mobId trong khu vuc. mobId la byte THO doc truc
// tiep tu Reader (khong dau) — gia tri 255 (= -1 khi doc CO dau, dac ta §1.4)
// la truong hop "thu cung" (dac ta §0 "Ngoai KH6", chua lam): actor KHONG can
// nhanh rieng cho no, vi Quai.Id LUON nam trong [0,127] (taoQuaiMacDinh,
// quai.go) nen 255 tu nhien khong khop con quai nao — roi thang vao nhanh
// "khong tim thay muc tieu" -> im lang bo qua, dung tinh than voi cach
// suKienDiChuyen bo qua khi khong tim thay nguoi trong k.nguoiTrong.
//
// Luong phien (xuLyDanhQuai, login.go) CHI doc gia tri mobId nay tu goi (ke
// ca doc BO 4 byte masterId khi mobId==-1 de tieu thu DU than goi bien thien,
// dac ta §2) roi gui su kien nay vao hop thu — KHONG tu tra k.quai, KHONG tu
// tinh sat thuong. Toan bo kiem tra (KiemDanhHopLe) + tinh sat thuong + tru
// mau + quyet chet + dung/phat goi nam TRONG vongLap (xem case duoi day) —
// dung ky luat actor dac ta §1.6 va rang buoc cua Task 7.
type suKienDanhQuai struct {
	idNhanVat int64
	mobId     byte
}

// suKienHoiSinh bao actor mot nhan vat vua gui goi len -16 (dac ta §3.1, Ke
// hoach 6 Task 13): xin hoi sinh sau khi chet. Luong phien (xuLyHoiSinh,
// login.go) CHI gui su kien nay — actor moi la noi kiem "dang THAT SU chet"
// + hoi mau/chakra day du + dat lucHoiSinh + phat goi -16/-30-15 (dung ky
// luat actor §1.6, xem case suKienHoiSinh duoi day). Dung ten truong
// idNhanVat giong het cac su kien khac trong file nay (suKienDiChuyen,
// suKienRa, suKienDanhQuai) de nhat quan.
type suKienHoiSinh struct {
	idNhanVat int64
}

// suKienDanhNguoiChoi bao actor mot nhan vat vua gui goi len -60 (dac ta
// §3.1/§2.1, Ke hoach 6 Task 11, Nhom F "nguoi danh nguoi"): xin danh nguoi
// choi co Id = idMucTieu trong khu vuc. Than goi len -60 CHI CO mot truong
// `int idNguoiBiDanh` (khac han goi 54 danh quai, KHONG co do dai bien thien
// o chieu LEN) nen luong phien (xuLyDanhNguoiChoi, login.go) chi doc DUNG mot
// gia tri roi gui thang su kien nay vao hop thu — TOAN BO kiem tra
// (KiemDanhHopLe) + tinh sat thuong + tru mau + quyet chet + dung/phat goi
// nam TRONG vongLap (xem case duoi day), dung ky luat actor §1.6 y het case
// suKienDanhQuai.
type suKienDanhNguoiChoi struct {
	idNhanVat int64 // nguoi DANH (A)
	idMucTieu int64 // nguoi BI DANH (B) — chuyen tu int32 doc duoc trong than goi len
}

// heQuaiChuaXacDinh la gia tri he nguyen to gan cho MOI con quai khi tinh sat
// thuong (Task 7): Quai (Task 5, quai.go) CHUA co truong theo doi he nguyen
// to rieng cho tung con. Dat DUNG BANG soHeNguyenTo (satthuong.go) — tuc
// NGOAI khoang hop le [0, soHeNguyenTo) — de heSoNguyenToPhanTram tu nhien
// coi day la "he khong xac dinh" va tra ve 100% (trung tinh), khong can sua
// ham do hay them nhanh dac biet o day. Neu sau nay quai co he nguyen to that
// (vd theo mau quai), doi hang so nay thanh mot tra cuu that.
const heQuaiChuaXacDinh byte = soHeNguyenTo

// heNguyenToTheoLang anh xa Lang.ID (0..4, thu tu danhSachLang trong lang.go:
// La, Cat, Suong Mu, Da, May) sang hang so he nguyen to cua satthuong.go
// (HeHoa/HePhong/HeLoi/HeTho/HeThuy). HAI THU TU NAY LECH NHAU — KHONG duoc
// suy dien "lang i = he i": lang 2 (Suong Mu) la He Thuy (hang so 4), lang 4
// (May) la He Loi (hang so 2). Dung theo thiet ke tong §7.2 va chieuthuc.go
// (xem tenChieuTheoLang: "Lang Suong Mu — Thuy", "Lang May — Loi").
var heNguyenToTheoLang = [5]byte{HeHoa, HePhong, HeThuy, HeTho, HeLoi}

// heNguyenToNguoiChoi tra he nguyen to cua mot nhan vat theo lang, dung lam
// tham so heNguoiDanh cho TinhSatThuong (Task 7 — Task 4 chi dung ham, khong
// quyet dinh nguon du lieu nay). Lang ngoai [0,4] (du lieu hong, ve ly thuyet
// khong xay ra vi TimLang da chan luc tao nhan vat) tra heQuaiChuaXacDinh de
// roi vao nhanh trung tinh, khong panic.
func heNguyenToNguoiChoi(lang int16) byte {
	if lang < 0 || int(lang) >= len(heNguyenToTheoLang) {
		return heQuaiChuaXacDinh
	}
	return heNguyenToTheoLang[lang]
}

// suKienChay la DIEM CHEN (seam) danh RIENG cho kiem thu luoi panic (manh A
// Task 2/3): xuLy chi goi sk.f(). VI SAO can no thay vi loi dung mot su kien
// that co the panic: test luoi panic can bom mot panic CO KIEM SOAT, CHAC
// CHAN, vao dung goroutine actor (khong phu thuoc chi tiet cai dat cua mot su
// kien that — neu mai sau ai them nil-check vao su kien do, panic bien mat,
// test luoi hoa vo dung). Day KHONG phai nhanh "if test" trong logic that: no
// la mot loai su kien rieng, runtime that khong bao gio gui no. Giu o day (canh
// cac suKien khac) de xuLy nhan dien duoc trong type switch.
type suKienChay struct{ f func() }

// suKienPhatGoi la yeu cau "phat san goi tin nay cho ca khu vuc" den tu MOT
// GOROUTINE NGOAI actor — hien tai chi co goroutine loa (SoKhuVuc.PhatToanServer,
// sokhuvuc.go) gui no.
//
// VI SAO PHAI DI QUA HOP THU chu khong goi thang k.phatChoCaKhuVuc: danh sach
// nguoiTrong thuoc so huu cua vongLap (luat §1.6) — mot goroutine khac duyet
// no la data race. Va vi sao la hop thu GAMEPLAY (k.Gui) chu khong phai hop
// thu DIEU KHIEN: guiDieuKhien CHO toi thoiGianChoControl = 100ms; goroutine
// loa dung o do thi ca duong phat thanh nghen theo, va neu mai sau co ai goi
// PhatToanServer tu trong mot actor thi do la tu khoa chac chan. k.Gui khong
// bao gio chan — hop thu day thi bo goi (mot dong chu mat di la chap nhan
// duoc; mot actor bi treo thi khong).
type suKienPhatGoi struct{ m *protocol.Message }

// Hai su kien truy van duoi day chi con la fallback cho integration test cu
// tu cap hopThu ma khong BatDau. Runtime dung control mailbox rieng.
type suKienXinDanhSachQuai struct {
	traLoi chan []*Quai
}

// suKienXinAnhChup la su kien TRUY VAN THU HAI cung mo hinh voi
// suKienXinDanhSachQuai o tren (Ke hoach 11 Task 2, dac ta
// 2026-07-27-luu-trang-thai-song.md §5): mot luong BEN NGOAI vongLap — bo
// dinh ky luu (§4 diem A), noi xu ly mat ket noi (§4 diem B,
// BoXuLy.KhiDong), hoac nhanh tat may chu (§4 diem C) — can "chup anh" trang
// thai store.NhanVat cua nguoi choi dang trong khu vuc DE DEM DI LUU XUONG
// DIA, nhung KHONG duoc doc truc tiep k.nguoiTrong (dac ta chien dau §1.6, ap
// dung y het cho truong hop nay: chi luong vongLap duoc cham trang thai
// nguoi choi trong actor).
//
// idNhanVat = 0 nghia la xin TAT CA nguoi trong khu vuc (dung cho bo dinh ky
// luu, §4 diem A: moi nhip luu xin het mot luot roi day vao hang doi cho tho
// ghi); idNhanVat != 0 nghia la xin DUNG mot nguoi co ID do (dung cho luc mot
// phien mat ket noi, §4 diem B: chi can luu dung nguoi vua roi, khong can ca
// khu vuc).
//
// traLoi tra ve []store.NhanVat — GIA TRI, KHONG PHAI []*store.NhanVat — day
// la diem khac CO CHU DICH so voi suKienXinDanhSachQuai (tra []*Quai): con
// tro nguoiTrongKhuVuc.nv la con tro DUNG CHUNG voi luong phien (dac ta §2.1
// "Chúng là cùng một ô nhớ, không phải hai bản sao"), trong khi con tro *Quai
// goc chi mot minh actor giu — dua *store.NhanVat ra ngoai (du la con tro toi
// mot BAN SAO moi tinh) van de lai thoi quen nguy hiem; buoc kieu tra ve phai
// la GIA TRI xoa han kha nang do o ca goc do bien dich, khong con dua vao ky
// luat lap trinh vien.
type suKienXinAnhChup struct {
	idNhanVat int64
	traLoi    chan []store.NhanVat
}

// Control co mailbox rieng, khong the bi gameplay day lam roi. Cac request
// deu co kenh tra loi buffer 1 de actor khong bao gio cho nguoi goi.
type suKienDieuKhien interface{ laDieuKhien() }

type dieuKhienXinQuai struct{ traLoi chan []*Quai }

func (dieuKhienXinQuai) laDieuKhien() {}

type dieuKhienXinAnhChup struct {
	idNhanVat int64
	traLoi    chan []store.NhanVat
}

func (dieuKhienXinAnhChup) laDieuKhien() {}

type dieuKhienRa struct {
	idNhanVat int64
	xong      chan struct{}
}

func (dieuKhienRa) laDieuKhien() {}

// dieuKhienSinhBoss la LENH SINH BOSS do bo dieu phoi lich (lichboss.go, NGOAI
// actor) gui vao (KH6b Task 13, CHAN-3+CHAN-7). Con boss da duoc DUNG XONG HOAN
// TOAN o ben ngoai — actor chi viec nhan lay va cam vao khu vuc.
//
// Vi sao dung HOP THU DIEU KHIEN chu khong phai hop thu gameplay: hop thu
// gameplay la best-effort (day thi BO goi, dem soLanGuiRoi). Mot lenh sinh boss
// bi bo im lang nghia la ca chu ky lich do khong co boss ma khong ai hay —
// trong khi hop thu dieu khien bao loi RO RANG cho nguoi goi (qua han / actor
// da dung) de bo dieu phoi con biet duong hoan lai hay ngung han.
//
// CHUYEN GIAO QUYEN SO HUU: sau khi gui con tro nay vao kenh, bo dieu phoi
// TUYET DOI khong duoc cham vao *Boss nua — ke tu day no thuoc ve actor (dung
// quy tac so huu cua k.quai/k.boss). Day la ly do lenh mang mot con tro DA DUNG
// SAN thay vi mang "du lieu de actor tu dung": neu actor tu dung, no se phai
// doc lich/dong ho/chon ngau nhien — dung cai ma CHAN-7 cam.
type dieuKhienSinhBoss struct{ boss *Boss }

func (dieuKhienSinhBoss) laDieuKhien() {}

type dieuKhienDung struct{}

func (dieuKhienDung) laDieuKhien() {}

var (
	ErrActorChuaBatDau = errors.New("actor chua bat dau")
	ErrActorDaDung     = errors.New("actor da dung")
	ErrControlQuaHan   = errors.New("control actor qua han")
)

const thoiGianChoControl = 100 * time.Millisecond

// MoKhuVuc tao mot khu vuc giu ban do da cho, chua co nguoi choi nao. Quai
// (Ke hoach 6 Task 5) duoc sinh MOT LAN o day bang taoQuaiMacDinh — day la
// diem duy nhat gan quai vao mot khu vuc, giu dung nguyen tac "quai thuoc
// actor, sinh luc MoKhuVuc" cua brief Task 5.
func MoKhuVuc(bd *BanDo) *KhuVuc {
	k := &KhuVuc{
		banDo:        bd,
		nguoiTrong:   make(map[int64]*nguoiTrongKhuVuc),
		anhChupDaRoi: make(map[int64]store.NhanVat),
		quai:         taoQuaiMacDinh(bd),
		dongHo:       time.Now,
	}
	// Ke hoach 6 Task 8: khoi tao LanHoiMauCuoi = thoi diem spawn cho MOI con
	// quai — tranh zero value (nam 1) khien lan kiem hoi mau DAU TIEN cua
	// nhipTim tuong nham "da qua 30s" ngay tick dau (xem chu thich truong
	// LanHoiMauCuoi, quai.go). Dat o day (sau khi da co k.dongHo) thay vi
	// truyen tham so vao taoQuaiMacDinh — giu chu ky taoQuaiMacDinh don gian,
	// khong ai khac goi ham do nen doi day khong pha vo gi.
	luc := k.dongHo()
	for _, q := range k.quai {
		q.LanHoiMauCuoi = luc
	}
	return k
}

// BanDo tra ban do cua khu vuc (client can no de dung MAP INFO).
func (k *KhuVuc) BanDo() *BanDo { return k.banDo }

// DanhSachQuaiSnapshot tra ban sao qua control mailbox. Ban CoLoi phan biet
// actor chua bat dau/da dung; wrapper cu giu tuong thich bang cach tra nil.
func (k *KhuVuc) DanhSachQuaiSnapshot() []*Quai {
	ban, _ := k.DanhSachQuaiSnapshotCoLoi()
	return ban
}

func (k *KhuVuc) DanhSachQuaiSnapshotCoLoi() ([]*Quai, error) {
	if k.hopThu == nil {
		return nil, ErrActorChuaBatDau
	}
	traLoi := make(chan []*Quai, 1)
	// Mot so integration test cu cap phat hopThu bang tay de bom tung event,
	// khong chay actor. Giu fallback nay cho harness do; runtime BatDau luon
	// co mailbox control rieng.
	if k.hopThuDieuKhien == nil {
		k.Gui(suKienXinDanhSachQuai{traLoi: traLoi})
		select {
		case ds := <-traLoi:
			return ds, nil
		case <-time.After(100 * time.Millisecond):
			return nil, ErrActorChuaBatDau
		}
	}
	if err := k.guiDieuKhien(dieuKhienXinQuai{traLoi: traLoi}); err != nil {
		return nil, err
	}
	select {
	case ds := <-traLoi:
		return ds, nil
	case <-k.xong:
		return nil, ErrActorDaDung
	case <-time.After(thoiGianChoControl):
		return nil, ErrControlQuaHan
	}
}

// AnhChupNguoiChoi tra BAN SAO GIA TRI store.NhanVat cua nguoi choi dang
// trong khu vuc, AN TOAN goi tu BEN NGOAI vongLap (Ke hoach 11 Task 2, dac ta
// §5) — cung khuon voi DanhSachQuaiSnapshot ngay tren, chi doi doi tuong tu
// quai sang nguoi choi (xem chu thich suKienXinAnhChup ve ly do doi kieu tra
// ve tu con tro sang gia tri).
//
// idNhanVat = 0 -> TAT CA nguoi trong khu vuc; idNhanVat != 0 -> DUNG mot
// nguoi co ID do (khong tim thay -> nil, dung tinh than "khong tim thay ->
// im lang" cua cac case suKien... khac trong file nay).
//
// Runtime di qua control mailbox reliable; ban CoLoi tra loi ro khi actor chua
// bat dau/da dung. Wrapper cu giu tuong thich cho cac call site chi can slice.
func (k *KhuVuc) AnhChupNguoiChoi(idNhanVat int64) []store.NhanVat {
	ban, _ := k.AnhChupNguoiChoiCoLoi(idNhanVat)
	return ban
}

func (k *KhuVuc) AnhChupNguoiChoiCoLoi(idNhanVat int64) ([]store.NhanVat, error) {
	if k.hopThu == nil {
		return nil, ErrActorChuaBatDau
	}
	traLoi := make(chan []store.NhanVat, 1)
	if k.hopThuDieuKhien == nil {
		k.Gui(suKienXinAnhChup{idNhanVat: idNhanVat, traLoi: traLoi})
		select {
		case ban := <-traLoi:
			return ban, nil
		case <-time.After(100 * time.Millisecond):
			return nil, ErrActorChuaBatDau
		}
	}
	if err := k.guiDieuKhien(dieuKhienXinAnhChup{idNhanVat: idNhanVat, traLoi: traLoi}); err != nil {
		return nil, err
	}
	select {
	case ban := <-traLoi:
		return ban, nil
	case <-k.xong:
		return nil, ErrActorDaDung
	case <-time.After(thoiGianChoControl):
		return nil, ErrControlQuaHan
	}
}

// RoiKhuVuc la control reliable cho disconnect: khong dung gameplay mailbox
// best-effort, va chi tra ve sau khi actor da xoa nguoi khoi map.
func (k *KhuVuc) RoiKhuVuc(idNhanVat int64) error {
	xong := make(chan struct{})
	if err := k.guiDieuKhien(dieuKhienRa{idNhanVat: idNhanVat, xong: xong}); err != nil {
		return err
	}
	select {
	case <-xong:
		return nil
	case <-k.xong:
		return ErrActorDaDung
	case <-time.After(thoiGianChoControl):
		return ErrControlQuaHan
	}
}

// SinhBoss la DUONG VAO DUY NHAT tu ben ngoai actor de dua mot con boss vao
// khu vuc (KH6b Task 13). Chi bo dieu phoi lich (lichboss.go) goi — va sau nay
// la duong spawn theo trigger cua ban 2.
//
// ⚠️ TUYET DOI KHONG goi tu TRONG actor (tu vongLap/nhipTim/xuLy): guiDieuKhien
// cho toi thoiGianChoControl = 100ms tren CHINH hop thu ma vongLap dang ban ->
// tu khoa chac chan (bay §6.10b). Sinh boss tu trong actor thi goi thang
// themBoss.
//
// Tra loi thay vi nuot: bo dieu phoi PHAI phan biet duoc "khu vuc dang ban, lat
// nua thu lai" (ErrControlQuaHan) voi "khu vuc chet roi, thu lai vo ich"
// (ErrActorDaDung) — xem thuSinh trong lichboss.go.
//
// b == nil van di qua hop thu binh thuong; themBoss ben trong actor tu bo qua
// (khong panic, khong them o nil vao danh sach).
func (k *KhuVuc) SinhBoss(b *Boss) error {
	return k.guiDieuKhien(dieuKhienSinhBoss{boss: b})
}

func (k *KhuVuc) guiDieuKhien(sk suKienDieuKhien) error {
	if k.hopThuDieuKhien == nil || k.xong == nil {
		return ErrActorChuaBatDau
	}
	select {
	case <-k.xong:
		return ErrActorDaDung
	default:
	}
	select {
	case k.hopThuDieuKhien <- sk:
		return nil
	case <-k.xong:
		return ErrActorDaDung
	case <-time.After(thoiGianChoControl):
		return ErrControlQuaHan
	}
}

// ThemNguoiChoi dua mot nhan vat vao khu vuc. Ke hoach 5 se broadcast su xuat
// hien nay cho nguoi khac; Ke hoach 4 chi giu danh sach de lam nen.
func (k *KhuVuc) ThemNguoiChoi(nv *store.NhanVat) {
	k.nguoiChoi = append(k.nguoiChoi, nv)
}

// SoNguoiChoi tra so nguoi dang trong khu vuc (dung cho kiem thu + Ke hoach 5).
func (k *KhuVuc) SoNguoiChoi() int { return len(k.nguoiChoi) }

// BatDau khoi dong goroutine actor cua khu vuc: cap phat hop thu + kenh dieu
// khien roi chay vongLap trong goroutine rieng. Goi mot lan truoc khi dung Gui.
func (k *KhuVuc) BatDau() {
	k.hopThu = make(chan suKienKhuVuc, 256)
	k.hopThuDieuKhien = make(chan suKienDieuKhien, 16)
	k.xong = make(chan struct{})
	go k.vongLap()
}

// ghiNhanNhip cap nhat van an toan sau MOT nhip gameplay (xuLy/nhipTim) va tra
// ve co NEN DUNG rieng khu vuc hay khong. daPanic=true (nhip vua roi bi panic,
// da bi ChayCoLuoi chan) -> tang bo dem lien tiep; false (nhip trot) -> reset
// ve 0. Tach ham thuan RIENG (khong nhet thang vao vongLap) de test goi mot
// chuoi panic/trot MA KHONG phu thuoc ticker thoi gian thuc — dung tinh than
// da tach kiemVaLuu/chonMucTieu/hoiSinhQuai khoi vong lap o cac task truoc.
//
// CHI luong vongLap goi ham nay (qua ticker/hop thu) hoac test goi truc tiep
// khi actor chua BatDau -> soPanicLienTiep khong bao gio bi hai luong cham.
func (k *KhuVuc) ghiNhanNhip(daPanic bool) (nenDung bool) {
	if !daPanic {
		k.soPanicLienTiep = 0 // mot nhip trot -> loi thoang qua khong tich luy
		return false
	}
	k.soPanicLienTiep++
	if k.soPanicLienTiep >= nguongPanicDungKhuVuc {
		slog.Error("khu vuc hong dai dang, dung rieng no de khong quay nong CPU + ngap log",
			"mapId", k.banDo.MapID, "soPanicLienTiep", k.soPanicLienTiep)
		return true
	}
	return false
}

// nguongNhipTre: khoang cach giua hai nhip vuot con so nay thi tinh la "tre"
// (2x chu ky 100ms). Kieu hang so de kiem thu dot bien sua duoc.
const nguongNhipTre = 200 * time.Millisecond

// capNhatTreNhip nhan khoang cach tu nhip truoc toi nhip nay. Vuot nguong thi
// dem (soNhipTre) + cap nhat max (treNhipLauNhatMs). Ham THUAN theo nghia khong
// doc dong ho — vongLap truyen khoang cach vao — nen test tat dinh duoc. Chi
// vongLap ghi treNhipLauNhatMs (mot chu so huu) nen khong dua ghi; atomic chi
// de test/giam sat DOC song song an toan.
func (k *KhuVuc) capNhatTreNhip(khoangCach time.Duration) {
	if khoangCach <= nguongNhipTre {
		return
	}
	k.soNhipTre.Add(1)
	if ms := khoangCach.Milliseconds(); ms > k.treNhipLauNhatMs.Load() {
		k.treNhipLauNhatMs.Store(ms)
	}
}

// danhDauNhip la DUONG NOI (dung trong vongLap) giua dong ho that va
// capNhatTreNhip: nhan moc thoi gian NHIP NAY, bo qua nhip dau (mocNhipTruoc
// zero), tinh khoang cach toi nhip truoc roi goi capNhatTreNhip, cuoi cung GAN
// LAI moc. Tach rieng de test tat dinh KHOA duoc duong noi nay (chi test
// capNhatTreNhip roi tin wiring la chua du — xem TestDanhDauNhip). CHI vongLap
// goi -> mocNhipTruoc mot chu so huu.
func (k *KhuVuc) danhDauNhip(bayGio time.Time) {
	if !k.mocNhipTruoc.IsZero() {
		k.capNhatTreNhip(bayGio.Sub(k.mocNhipTruoc))
	}
	k.mocNhipTruoc = bayGio
}

// vongLap la than actor: doc hop thu, nhip tim 100ms, hoac dung khi co lenh
// dung. Day la NOI DUY NHAT doc/ghi trang thai khu vuc (nguoiTrong, ...) —
// khong goroutine nao khac duoc dung truc tiep cac field nay.
func (k *KhuVuc) vongLap() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	defer close(k.xong)

	for {
		// Drain control truoc gameplay. Vong default nay tranh mot mailbox
		// gameplay luon day lam snapshot/leave/stop bi doi vo han.
		//
		// LUOI PANIC (manh A §3.2a): BOC TUNG lan xu ly bang antoan.ChayCoLuoi,
		// KHONG boc ca vongLap — boc ngoai thi panic van thoat vong lap, giet
		// CA khu vuc; boc trong thi mot su kien loi chi bo dung su kien do, nhip
		// sau khu vuc chay binh thuong (quyet dinh da chot cua chu du an: khu
		// vuc gap loi cu chay tiep).
		select {
		case sk := <-k.hopThuDieuKhien:
			var dung bool
			antoan.ChayCoLuoi("khuvuc.dieuKhien", func() { dung = k.xuLyDieuKhien(sk) })
			if dung {
				return
			}
			continue
		default:
		}
		select {
		case sk := <-k.hopThuDieuKhien:
			var dung bool
			antoan.ChayCoLuoi("khuvuc.dieuKhien", func() { dung = k.xuLyDieuKhien(sk) })
			if dung {
				return
			}
		case sk := <-k.hopThu:
			daPanic := antoan.ChayCoLuoi("khuvuc.xuLy", func() { k.xuLy(sk) })
			if k.ghiNhanNhip(daPanic) {
				return
			}
		case <-ticker.C:
			// KH6c: do khoang cach THUC (wall-clock) giua hai nhip. CO Y dung
			// time.Now() thang (KHONG phai k.dongHo tiem duoc) vi day la do NHIP
			// THAT cua vong lap — dongHo chi de tiem "gio game" cho logic gameplay;
			// neu test dong bang dongHo ma o day cung dung no thi tre nhip mu tit.
			// Toan bo logic moc/tinh-khoang-cach nam trong danhDauNhip de test khoa.
			k.danhDauNhip(time.Now())
			daPanic := antoan.ChayCoLuoi("khuvuc.nhipTim", k.nhipTim)
			if k.ghiNhanNhip(daPanic) {
				return
			}
		}
	}
}

func (k *KhuVuc) xuLyDieuKhien(sk suKienDieuKhien) bool {
	switch sk := sk.(type) {
	case dieuKhienXinQuai:
		k.xuLyGameplayDaXepTruoc()
		sk.traLoi <- k.saoQuai()
	case dieuKhienXinAnhChup:
		k.xuLyGameplayDaXepTruoc()
		sk.traLoi <- k.saoNguoiChoi(sk.idNhanVat)
	case dieuKhienRa:
		k.xuLyGameplayDaXepTruoc()
		k.xuLy(suKienRa{idNhanVat: sk.idNhanVat})
		delete(k.anhChupDaRoi, sk.idNhanVat)
		close(sk.xong)
	case dieuKhienSinhBoss:
		// KHONG goi xuLyGameplayDaXepTruoc o day (khac ba case tren): ba case
		// kia la TRUY VAN — chung phai thay het gameplay da xep truoc thi cau
		// tra loi moi dung. Day la MOT LENH, khong tra loi gi ca; va du no
		// chay truoc hay sau mot suKienVao dang cho thi nguoi moi van thay
		// boss: chay truoc -> suKienVao lap qua k.boss va gui -5; chay sau ->
		// themBoss phat -5 cho ca khu vuc (da co nguoi moi trong do).
		k.themBoss(sk.boss)
	case dieuKhienDung:
		return true
	}
	return false
}

// Control duoc uu tien lay khoi mailbox rieng, nhung snapshot/leave van phai
// thay cac gameplay event da enqueue TRUOC request. Chot so luong tai thoi
// diem bat dau de khong bi mot dong gameplay moi den lien tuc lam doi control
// vo han.
func (k *KhuVuc) xuLyGameplayDaXepTruoc() {
	soLuong := len(k.hopThu)
	for i := 0; i < soLuong; i++ {
		select {
		case sk := <-k.hopThu:
			k.xuLy(sk)
		default:
			return
		}
	}
}

func (k *KhuVuc) saoQuai() []*Quai {
	ban := make([]*Quai, len(k.quai))
	for i, q := range k.quai {
		cp := *q
		ban[i] = &cp
	}
	return ban
}

func (k *KhuVuc) saoNguoiChoi(idNhanVat int64) []store.NhanVat {
	if idNhanVat == 0 {
		ban := make([]store.NhanVat, 0, len(k.nguoiTrong))
		for _, nguoi := range k.nguoiTrong {
			ban = append(ban, *nguoi.nv)
		}
		return ban
	}
	if nguoi, co := k.nguoiTrong[idNhanVat]; co {
		return []store.NhanVat{*nguoi.nv}
	}
	if ban, co := k.anhChupDaRoi[idNhanVat]; co {
		return []store.NhanVat{ban}
	}
	return nil
}

// xuLy xu ly mot su kien lay tu hop thu — day la NOI DUY NHAT doc/ghi trang
// thai khu vuc (goi tu vongLap, hoac truc tiep tu test khi actor CHUA BatDau).
//
// THU TU phat tin (quan trong, xem brief KH5 Task 5):
//   - suKienVao: gui -5 cua TUNG nguoi CU va cua TUNG BOSS dang hien cho nguoi
//     MOI TRUOC (nguoi moi thay moi thu da co san), roi moi them nguoi moi vao
//     map (de khong tu gui -5 cho chinh minh o buoc sau), roi phat -5 cua nguoi
//     moi cho nhung nguoi con lai. Boss khong nhan gi o buoc cuoi — no khong co
//     phien mang, khong can biet ai vua vao.
//   - suKienDiChuyen: cap nhat vi tri THAT (authoritative) truoc, roi phat -7
//     cho NGUOI KHAC (client tu ve chuyen dong cua chinh minh, khong gui lai).
//   - suKienRa: xoa khoi map TRUOC, roi phat -6 cho nguoi con lai (nguoi da
//     ra khong con trong map nen tu nhien khong nhan duoc goi cua chinh minh).
func (k *KhuVuc) xuLy(sk suKienKhuVuc) {
	switch sk := sk.(type) {
	case suKienChay:
		// Diem chen test luoi panic (xem suKienChay o dau file): chay ham do
		// bom vao. Runtime that khong bao gio gui su kien nay.
		sk.f()

	case suKienPhatGoi:
		// Goi da duoc DUNG SAN o ngoai (goroutine loa) — actor chi lam dung
		// mot viec re: day no toi tung nguoi dang trong khu vuc. KHONG kiem
		// nil o day: nguoi GUI chiu trach nhiem do (SoKhuVuc.PhatToanServer da
		// kiem, xem sokhuvuc.go). Mot nhanh phong thu thua o day se la nhanh
		// khong test nao giet duoc — che mat chinh cai loi no dinh chan.
		k.phatChoCaKhuVuc(sk.m)

	case suKienVao:
		delete(k.anhChupDaRoi, sk.nv.ID)
		// 1) Gui -5 cua TUNG nguoi CU cho nguoi MOI (de nguoi moi thay nguoi cu).
		//    Lam TRUOC khi them nguoi moi vao map, de nguoi moi khong tu thay minh.
		for _, cu := range k.nguoiTrong {
			m, err := goiXuatHien(cu.nv)
			if err != nil {
				slog.Error("dung goi -5 nguoi cu that bai", "id", cu.nv.ID, "loi", err)
				continue
			}
			if err := sk.phien.XepGui(m); err != nil {
				slog.Warn("gui -5 nguoi cu cho nguoi moi that bai", "loi", err)
				sk.phien.DongCham()
				return
			}
		}
		// 1b) Gui -5 cua tung BOSS dang HIEN cho nguoi MOI (KH6b Task 7). Van
		//     thuoc buoc (1) "nguoi moi thay nhung gi da co san": boss la MOT
		//     NHAN VAT voi client, dung chung goi -5 voi nguoi choi. Boss KHONG
		//     can biet ai vua vao (no khong co phien mang) nen KHONG co chieu
		//     nguoc lai o buoc (3).
		//
		//     Bo qua boss chua hien (NGHI) hoac da don di (RA_DI) — xem
		//     TrangThaiBoss.DangHien. Bo qua ca goi nil (du lieu pha hong/ten
		//     qua dai): goiXuatHienBoss da ghi log ro ly do, o day chi tuyet
		//     doi KHONG XepGui mot con tro nil.
		for _, b := range k.boss {
			if b == nil || !b.TrangThai.DangHien() {
				continue
			}
			m := goiXuatHienBoss(b)
			if m == nil {
				continue
			}
			if err := sk.phien.XepGui(m); err != nil {
				slog.Warn("gui -5 boss cho nguoi moi that bai", "loi", err)
				sk.phien.DongCham()
				return
			}
		}
		// 2) Them nguoi moi vao khu vuc.
		k.nguoiTrong[sk.nv.ID] = &nguoiTrongKhuVuc{nv: sk.nv, phien: sk.phien}
		// 3) Phat -5 cua nguoi moi cho NGUOI KHAC (khong gui lai chinh minh).
		m, err := goiXuatHien(sk.nv)
		if err != nil {
			slog.Error("dung goi -5 nguoi moi that bai", "id", sk.nv.ID, "loi", err)
			return
		}
		k.phatChoNguoiKhac(sk.nv.ID, m)

	case suKienDiChuyen:
		nguoi, co := k.nguoiTrong[sk.idNhanVat]
		if !co {
			return
		}
		// KE HOACH 5: kiem hop le PHAI o day (actor so huu nguoi.nv), KHONG o luong
		// phien — luong phien doc nv.ViTri se DUA voi actor ghi chung con tro. Kiem
		// bang vi tri AUTHORITATIVE hien tai.
		if !kiemDiChuyenHopLe(k.banDo, int16(nguoi.nv.ViTriX), int16(nguoi.nv.ViTriY), sk.x, sk.y) {
			slog.Debug("bo goi di chuyen bat hop le (ngoai bien / vuot nguong buoc)",
				"id", sk.idNhanVat, "mapId", k.banDo.MapID, "toX", sk.x, "toY", sk.y)
			return // giu vi tri cu, KHONG phat
		}
		nguoi.nv.ViTriX = int32(sk.x)
		nguoi.nv.ViTriY = int32(sk.y)
		k.phatChoNguoiKhac(sk.idNhanVat, goiDiChuyen(sk.idNhanVat, sk.x, sk.y))

	case suKienRa:
		if _, co := k.nguoiTrong[sk.idNhanVat]; !co {
			return // stub an toan
		}
		delete(k.nguoiTrong, sk.idNhanVat)
		// Phat -6 cho nhung nguoi CON LAI (nguoi ra da bi xoa nen khong nhan).
		k.phatChoNguoiKhac(sk.idNhanVat, goiRoiKhuVuc(sk.idNhanVat))

	case suKienChonChieu:
		nguoi, co := k.nguoiTrong[sk.idNhanVat]
		if !co {
			return
		}
		// Tim MauChieu co Id == mauChieuId TRONG DANH MUC CUA LANG nhan vat
		// (ChieuCuaLang(nv.Lang)) — day la buoc "kiem chieu thuoc lang" ma dac
		// ta §3.2 doi hoi actor lam (khong phai luong phien). Lang khong hop
		// le (ngoai 0..4) hoac templateId khong co trong 4 chieu cua lang do
		// -> khong tim thay -> BO QUA, GIU NGUYEN chieuDangChon cu (dung dac
		// ta: "chon chieu khong thuoc lang -> khong doi gi").
		danhMuc := ChieuCuaLang(byte(nguoi.nv.Lang))
		for i := range danhMuc {
			if danhMuc[i].Id == sk.mauChieuId && len(danhMuc[i].CacCap) > 0 {
				// Luon lay cap 1 (CacCap[0]): Ke hoach 6 Task 2 chua co mo hinh
				// cap do chieu cua tung nguoi choi (se den o task sau) — chon
				// chieu la chon LOAI, cap that do task khac quyet dinh.
				nguoi.chieuDangChon = &danhMuc[i].CacCap[0]
				// Fix wave cuoi Task 11: ghi CUNG LUC loaiChieuDangChon (xem
				// chu thich truong do o dau file) — day la DIEM DUY NHAT trong
				// ca actor gan chieuDangChon nen la noi DUNG NHAT de giu hai
				// truong nay dong bo, khong the lech nhau.
				nguoi.loaiChieuDangChon = danhMuc[i].Loai
				return
			}
		}

	case suKienDanhQuai:
		// Dac ta §3.1 (goi len 54) + §1.6 (ky luat actor) + brief Task 7: toan
		// bo kiem tra/tinh toan/phat goi cho MOT don danh quai nam TRON trong
		// nhanh nay — luong phien (xuLyDanhQuai, login.go) chi gui su kien,
		// khong tu lam gi trong cac buoc duoi day.
		nguoi, co := k.nguoiTrong[sk.idNhanVat]
		if !co {
			return
		}

		// Ke hoach 6 Task 13 (test bat buoc cua brief): NGUOI DA CHET khong
		// duoc danh. Chu dinh dat CHOT NAY o day (truoc ca KiemDanhHopLe),
		// khong nhet vao chotkiem.go: KiemDanhHopLe (Task 3) kiem MUC TIEU
		// (mucTieu.HpHienTai — kiem #4), con day la kiem NGUOI DANH — hai
		// huong khac nhau, brief Task 13 CHI liet ke goichiendau.go/khuvuc.go
		// la file phai sua, KHONG co chotkiem.go. Ban Java (Controller.java
		// case 54) KHONG co chot nay (chi -20 nhat do co !player.isDie()) —
		// day la mot luat choi RO RANG can them, khong phai chep nguyen ban
		// goc: nguoi choi dang cho hoi sinh khong the "tranh thu" danh quai.
		if nguoi.nv.HpHienTai <= 0 {
			return // im lang bo qua, giong het cach cac kiem khac cua case nay tu choi
		}

		// Tim con quai co Id == sk.mobId. mobId=255 (thu cung, xem chu thich
		// suKienDanhQuai o tren) tu nhien khong khop quai nao -> mucTieu==nil
		// -> bo qua, dung nhu suKienDiChuyen bo qua khi khong tim thay nguoi.
		var mucTieu *Quai
		for _, q := range k.quai {
			if q.Id == sk.mobId {
				mucTieu = q
				break
			}
		}
		if mucTieu == nil {
			return
		}

		bayGio := k.dongHo() // Task 8: dong ho tiem duoc, KHONG goi time.Now() thang o day nua.
		ketQua := KiemDanhHopLe(benDanhTuNguoi(nguoi), ViTriCoMau{
			X:         mucTieu.X,
			Y:         mucTieu.Y,
			HpHienTai: mucTieu.Hp,
			// LucHoiSinh: Ke hoach 6 Task 8 da them truong nay vao Quai
			// (quai.go) va nhipTim gan gia tri that moi lan hoi sinh -> doc
			// THANG tu mucTieu o day, KHONG con bo trong nhu truoc Task 8.
			// Zero value (quai moi spawn, chua tung hoi sinh) van tu nhien
			// cho phep danh ngay — dung dac diem da ghi o chu thich truong
			// LucHoiSinh (chotkiem.go va quai.go).
			LucHoiSinh: mucTieu.LucHoiSinh,
		}, bayGio)
		if !ketQua.ChoPhep {
			// 4 kiem dau (bat tu/chakra/hoi chieu/da chet) hong -> IM LANG,
			// KHONG phat goi gi (khac han nhanh Truot duoi day — dac ta §1.9
			// chi noi rieng ve vuot tam, KHONG ap dung cho 4 kiem nay).
			return
		}

		// Cap nhat lanDanhCuoi CHI KHI duoc phep danh (brief Task 7). Mot don
		// bi 4 kiem dau tu choi KHONG duoc tinh la "vua dung chieu", neu
		// khong nguoi choi se bi khoa hoi chieu oan boi chinh don bi tu choi.
		nguoi.lanDanhCuoi = bayGio

		// chiMang: KH6 CHUA chot ty le/cong thuc chi mang that (khong dac ta
		// nao, khong task truoc nao dinh nghia) -> luon false — lua chon an
		// toan, xem bao cao Task 7. TinhSatThuong van nhan tham so nay dung
		// chu ky da chot o Task 4; ty le chi mang that (neu co, co the la
		// KH8) chi can thay DONG GAN chiMang nay, khong doi gi khac o day.
		chiMang := false

		var satThuong int32
		if !ketQua.Truot {
			lang := TimLang(nguoi.nv.Lang)
			var satThuongGoc int64
			if lang != nil {
				// satThuongGoc: dung Lang.SatThuongKhoiDau (lang.go) — truong
				// nay co san tu Ke hoach 2, chu thich goc ghi "dung o Ke
				// hoach 4" (ten cu cua chien dau truoc khi tach thanh KH6),
				// dam bao KHAC 0 cho ca 5 lang (10-16). nv.SucManh KHONG
				// dung duoc o day: la diem xep hang hien thi rieng (bang xep
				// hang web, xem web/app/bang-xep-hang), mac dinh 0 va KHONG
				// duoc TaoNhanVat gan luc tao nhan vat — dung no se lam MOI
				// nhan vat moi tao gay 0 sat thuong vinh vien.
				satThuongGoc = int64(lang.SatThuongKhoiDau)
			}
			heNguoiDanh := heNguyenToNguoiChoi(nguoi.nv.Lang)
			satThuong = TinhSatThuong(satThuongGoc, nguoi.chieuDangChon, heNguoiDanh, heQuaiChuaXacDinh, chiMang)

			// Ke hoach 6 Task 9: don TRUNG (nhanh nay, khac nhanh Truot o
			// duoi) -> ghi nhan nguoi nay vao danh sach "da tung danh" cua
			// quai, kich hoat tang uu tien 1 cua AI (dac ta §5.2, quai.go/
			// chonMucTieu) cho NHUNG DON SAU. Dat NGAY trong nhanh "trung
			// don" (khong phai nhanh Truot) vi "trung don" (brief Task 9)
			// nghia la thuc su gay sat thuong, khac "vuot tam" noi don khong
			// thuc su cham muc tieu.
			if mucTieu.DaTungDanh == nil {
				mucTieu.DaTungDanh = make(map[int64]bool)
			}
			mucTieu.DaTungDanh[nguoi.nv.ID] = true
		}
		// Truot=true (vuot tam): satThuong giu nguyen 0 (dac ta §1.9) — KHONG
		// goi TinhSatThuong, vi don nay khong thuc su cham muc tieu.

		// EP TUONG MINH: vuot tam thi KHONG duoc chi mang — BAN GIAO Task 3
		// (chotkiem.go:75-81): "chiMang=false khi Truot=true", vi client CHI
		// hien chu "miss" khi chiMang==false (Controller.cs:4356-4363, dac ta
		// §3.3 diem 3) — Truot=true+chiMang=true se hien "-0" do thay vi
		// "miss". Dat DONG NAY o day, SAU moi cho co the gan/tinh chiMang (ke
		// ca cong thuc chi mang that trong tuong lai o nhanh !ketQua.Truot
		// phia tren): hien tai chiMang luon false nen dong nay dung "tinh co"
		// — mot khi co cong thuc chi mang that duoc nhet vao dung cho o tren,
		// dong nay la RANH GIOI CUOI CUNG bao dam nghia vu khong bao gio vo
		// lang le, du cong thuc do tinh sai o dau.
		if ketQua.Truot {
			chiMang = false
		}

		// §1.2: sat thuong khong duoc VUOT mau con lai cua muc tieu.
		if satThuong > mucTieu.Hp {
			satThuong = mucTieu.Hp
		}
		// §1.3: quai DANG DAY MAU khong chet boi MOT don — giat lai 1 mau.
		if mucTieu.Hp == mucTieu.HpToiDa && satThuong >= mucTieu.Hp {
			satThuong = mucTieu.Hp - 1
		}

		mucTieu.Hp -= satThuong
		if mucTieu.Hp < 0 { // phong thu: ve ly thuyet khong xay ra sau hai kep tren
			mucTieu.Hp = 0
		}

		// dac ta §3.2: skillId la BYTE trong goi 54 xuong (short trong catalog).
		skillId := byte(nguoi.chieuDangChon.SkillId)
		k.phatChoCaKhuVuc(goiNguoiDanhQuai(int32(nguoi.nv.ID), skillId, mucTieu.Id))

		if mucTieu.Hp <= 0 {
			// Task 8: ghi lai thoi diem chet DE nhipTim tinh duoc "da troi qua
			// bao lau" cho nguong 5 giay hoi sinh (dac ta §5.2). Dat TRUOC khi
			// phat goi -12 — thu tu phat goi khong doi, chi la ghi trang thai
			// noi bo actor truoc, giong cach nguoi.lanDanhCuoi duoc ghi truoc
			// khi phat goi 54/-9 phia tren.
			mucTieu.LucChet = bayGio
			// Task 9: quai CHET -> xoa sach danh sach "da tung danh no" (dung
			// nguyen van brief: "xoa sach khi quai chet") — lan hoi sinh sau
			// la mot "vong doi" moi, khong ai duoc coi la da tung danh no nua.
			mucTieu.DaTungDanh = nil
			// KH6 chua co he vat pham (dac ta §0) -> doRoi luon rong.
			k.phatChoCaKhuVuc(goiQuaiChet(mucTieu.Id, satThuong, chiMang, nil))
		} else {
			// hieuUngId=255 (=-1): KH6 chua co hieu ung trang thai (dac ta §0).
			k.phatChoCaKhuVuc(goiQuaiConSong(mucTieu.Id, mucTieu.Hp, satThuong, chiMang, 255))
		}

	case suKienDanhNguoiChoi:
		// Ke hoach 6 Task 11 (Nhom F "nguoi danh nguoi", dac ta §2.1/§3.1/§3.3)
		// + brief Task 11 Phan A: TOAN BO kiem tra + tinh toan + phat goi cho
		// MOT don PvP nam TRON trong nhanh nay — luong phien
		// (xuLyDanhNguoiChoi, login.go) chi gui su kien, dung ky luat actor
		// §1.6, cau truc gan y het case suKienDanhQuai o tren (khac chu yeu o
		// nguon sat thuong/muc tieu: nguoi that thay vi quai).
		nguoi, co := k.nguoiTrong[sk.idNhanVat]
		if !co {
			return
		}

		// Dung tinh than voi Task 13 (nhanh "nguoi da chet khong duoc danh
		// quai" o case suKienDanhQuai phia tren): nguoi DA CHET (dang cho hoi
		// sinh) khong duoc "tranh thu" danh NGUOI KHAC. Chot nay KHONG nam
		// trong KiemDanhHopLe (ham do chi kiem phia MUC TIEU, khong kiem phia
		// NGUOI DANH) nen dat truoc o day, truoc ca tra cuu muc tieu — day la
		// mot luat choi RO RANG can them (dung tinh than voi ly do da ghi o
		// case suKienDanhQuai), ap dung DOI XUNG cho ca hai huong danh.
		if nguoi.nv.HpHienTai <= 0 {
			return // im lang bo qua
		}

		// Fix wave cuoi (coordinator soi toan nhanh KH6, 27/07): chieu HOI
		// MAU (Loai == loaiBuff, chieuthuc.go) TUYET DOI KHONG duoc xu ly nhu
		// mot don danh — truoc dong nay, moi chieu deu roi thang xuong
		// TinhSatThuong ben duoi bat ke Loai, nen ai chon chieu hoi mau (vd Y
		// Lieu Nhan Thuat) roi gui -60 se lam MUC TIEU MAT MAU thay vi duoc
		// hoi — sai nguoc hoan toan y nghia chieu.
		//
		// CO CHU DICH CHUA xay co che hoi mau PvP that o day: dac ta chua du
		// de lam dung — con thieu (a) luong hoi bao nhieu (khong co cong
		// thuc nao duoc chot, khac han satThuongGoc/heSoChieu cua don danh),
		// (b) nguoi DUOC HOI biet mau moi cua minh qua goi nao (goi -60 dang
		// hoi mau, goiHoiMauNguoiChoi, dac ta §2.1 CHI co 4 truong dau roi
		// dung — KHONG mang hpConLai; va -30/9 §3.4 chi gui cho "nguoi khac",
		// tuc chinh nguoi duoc hoi KHONG tu biet mau minh qua goi nay, giong
		// he qua tuong tu vua gap voi don danh o Task 11 Phan A). Xay ngay
		// bay gio la DOAN giao thuc — dung bay "giai bai toan chua ton tai"
		// (dac ta chien dau §4.1 "cong thuc chua chot" + BAT-DAU-TU-DAY §6.11
		// "test xanh ma code sai"). PvP co luat (doi/dong minh — ai duoc hoi,
		// ai bi cam hoi ke dich) cung da hoan sang KH9 (dac ta §0).
		//
		// BO QUA AN TOAN: khong tinh sat thuong, khong tru mau ai, KHONG
		// phat goi tin nao ca (khac han nhanh "vuot tam" — day KHONG PHAI
		// mot don thuc su duoc thu, nen khong ap dung dac ta §1.9 "van phat
		// goi"). goiHoiMauNguoiChoi (goichiendau.go, dung xong tu Task 10,
		// van con nguyen ham + test byte-chinh-xac cua no) CO CHU DICH CHUA
		// duoc noi vao day — giu lai de dung NGAY khi co dac ta hoi mau PvP
		// day du, khong phai viet lai tu dau.
		//
		// Dung nguoi.loaiChieuDangChon (KHONG PHAI nguoi.chieuDangChon.Loai —
		// CapChieu, chieuthuc.go, khong co truong Loai; Loai chi ton tai o
		// MauChieu cha, xem chu thich truong loaiChieuDangChon o dau file).
		if nguoi.loaiChieuDangChon == loaiBuff {
			return
		}

		// ─── KH6b Task 8: don danh nham vao BOSS ────────────────────────────
		//
		// VI SAO NHANH NAY NAM O DAY chu khong phai mot goi len rieng: boss
		// vao the gioi bang CHINH goi -5 cua nguoi choi (goiXuatHienBoss,
		// xuathien.go) nen voi client no la mot "char" trong map — nut danh
		// gui LEN dung goi -60 nay, chi khac o con so id. Client khong he co
		// khai niem "danh boss".
		//
		// DAU CUA ID la thu duy nhat can de biet don nay nham vao ai: id boss
		// LUON AM (CapIdBoss, boss.go), id nhan vat lay tu bigserial nen LUON
		// DUONG. Chot dau AM o day (thay vi tra ca hai bang roi so ket qua) de
		// mot id bia dat KHONG BAO GIO lac duoc vao duong nguoi choi, va nguoc
		// lai — hai duong khong nuot nhau.
		//
		// CAM sua parser -60 trong login.go (CHAN-1 cua ke hoach KH6b): than
		// goi len that su la "danh sach id + cdir" (Service.cs:1256-1265) ma
		// may chu doc dung MOT int nen VO TINH DUNG voi mot muc tieu. Nhanh
		// nay CHI doc gia tri id da parse san, khong dung toi tang parser.
		if sk.idMucTieu < 0 {
			// Chi boss DANG TREN MAN HINH moi la muc tieu hop le.
			// TrangThai.DangHien() la DUNG dieu kien ma suKienVao va themBoss
			// dung de quyet dinh co gui -5 hay khong, nen o day no cho mot
			// dinh nghia duy nhat: "danh duoc" = "thay duoc". Boss dang NGHI
			// (dem nguoc, client chua he biet id do) va boss da RA_DI (da don
			// di) roi vao CUNG mot nhanh voi id bia dat: bo qua im lang, dung
			// cach nhanh nay da xu ly mot idMucTieu khong khop ai. Mot id bia
			// TUYET DOI khong duoc lam sap actor.
			//
			// Chot nay con che cho CLIENT: khi phat -60, client tim ca hai
			// dau bang findCharInMap va dung ket qua NGAY (char3.cx,
			// Controller.cs:3934) — doan do nam NGOAI khoi try/catch cua no.
			// Phat mot don nham vao con boss ma client chua he nhan -5 la ban
			// mot null-deref sang cho moi nguoi trong khu vuc.
			var boss *Boss
			for _, b := range k.boss {
				if b != nil && b.Id == sk.idMucTieu && b.TrangThai.DangHien() {
					boss = b
					break
				}
			}
			if boss == nil {
				return
			}

			bayGio := k.dongHo() // dong ho tiem duoc, KHONG time.Now() thang.
			ketQua := KiemDanhHopLe(benDanhTuNguoi(nguoi), ViTriCoMau{
				X: boss.ViTriX,
				Y: boss.ViTriY,
				// HpHienTai lam kiem #4 ("muc tieu da chet") tu nhien che not
				// CAI XAC boss: BossTroi LUON co Hp == 0 (bossChet dat) nen
				// khong can them mot nhanh trang thai rieng o day — va don vao
				// xac boss im lang DUNG NHU don vao mot nguoi choi da chet.
				HpHienTai: boss.Hp,
				// LucHienPha (KHONG PHAI MocThoiGian — xem chu thich hai
				// truong do trong boss.go): boss vua vao pha moi duoc mot
				// khoang bat tu 1500ms, giong het nguoi choi/quai vua hoi
				// sinh. MocThoiGian doi theo MOI buoc trang thai nen dung no
				// se lam boss bat tu lai o moi buoc.
				LucHoiSinh: boss.LucHienPha,
			}, bayGio)
			if !ketQua.ChoPhep {
				return // 4 kiem dau hong -> im lang, dung nhu hai nhanh danh o tren
			}
			nguoi.lanDanhCuoi = bayGio

			chiMang := false // KH6 chua chot cong thuc chi mang — xem case suKienDanhQuai
			var satThuong int32
			if !ketQua.Truot {
				lang := TimLang(nguoi.nv.Lang)
				var satThuongGoc int64
				if lang != nil {
					satThuongGoc = int64(lang.SatThuongKhoiDau)
				}
				// He nguyen to cua BOSS lay tu Lang cua PHA dang chay: boss co
				// lang that (PhaBoss.Lang, dung cho ca goi -5) nen KHONG dung
				// heQuaiChuaXacDinh nhu quai. Pha nil (du lieu noi dung hong)
				// -> roi ve trung tinh 100%, khong doan bua va khong nil-deref.
				heBoss := heQuaiChuaXacDinh
				if pha := boss.PhaHienTaiData(); pha != nil {
					heBoss = heNguyenToNguoiChoi(pha.Lang)
				}
				satThuong = TinhSatThuong(satThuongGoc, nguoi.chieuDangChon,
					heNguyenToNguoiChoi(nguoi.nv.Lang), heBoss, chiMang)
				// TODO(Task 11 — noi gian theo %mau): MocNoiGian.HeSoDef nhan
				// vao DUNG CHO NAY (sat thuong boss NHAN). Chua ap o Task 8 vi
				// may trang thai chua kich hoat moc nao.
			}
			// Vuot tam thi KHONG duoc chi mang — client CHI hien chu "miss"
			// khi chiMang == false. Dung ly do va dung vi tri voi hai nhanh
			// danh o tren (dat SAU moi cho co the gan chiMang).
			if ketQua.Truot {
				chiMang = false
			}

			// bossNhanSatThuong la DUONG DUY NHAT duoc tru mau boss (chu thich
			// cua chinh ham do): no ghi so dong gop, kep mau, va chay tron
			// nghi thuc "boss chet" (chia thuong + thong bao) khi don nay ket
			// lieu. Tu `boss.Hp -= satThuong` o day la bo qua CA hai thu do —
			// nguoi choi mat sach phan thuong ma khong mot dong log nao bao.
			hpTruoc := boss.Hp
			daChet := k.bossNhanSatThuong(boss, nguoi.nv.ID, int64(satThuong), bayGio)

			// Con so bao xuong client la SAT THUONG THAT SU DA AP (hieu mau
			// truoc/sau), khong phai con so vua tinh ra. Hai cai nay lech nhau
			// o dung cho dac ta §1.2 noi: don danh manh hon mau con lai. Lay
			// hieu thay vi tu kep lai o day co hai cai loi: (a) chi con MOT
			// noi kep mau — chinh bossNhanSatThuong — nen khong the lech nhau;
			// (b) neu ham do tu choi don (trang thai boss doi giua chung),
			// hieu bang 0 nen goi phat ra la mot cu "miss" that tha, khong bao
			// gio co chuyen goi bao mot con so ma may chu khong he tru.
			// Client tu tru con so nay khoi thanh mau boss (Char.doInjure) nen
			// bao thua la lech thanh mau ngay lap tuc.
			satThat := hpTruoc - boss.Hp

			// Goi phat ra la -60 (KHONG PHAI -9/-12 cua quai): hai goi quai
			// mang `byte mobId` — mot id boss am 4 byte khong nhet vao do
			// duoc — va client tim muc tieu cua -60 bang GameScr.findCharInMap
			// (Controller.cs:3989), dung cai kho ma boss dang nam trong nho
			// goi -5. Phat cho CA KHU VUC (ke ca nguoi danh) dung nhu hai
			// nhanh danh o tren.
			//
			// int32(boss.Id): dung phep thu hep ma goiXuatHienDayDu da dung
			// khi ghi id boss vao goi -5, nen con so o day KHOP CHINH XAC id
			// ma client dang giu.
			//
			// KHONG gui kem -30/9 (cap nhat mau) nhu nhanh PvP — hai ly do
			// khien no CAN o do deu khong ton tai o day: goi -5 cua NGUOI CHOI
			// ghi hp = HpToiDa (nguoi vao giua tran thay nguoi khac day mau)
			// nen phai co goi sua, con goi -5 cua BOSS ghi b.Hp THAT; va con
			// so trong -60 la sat thuong da ap that nen moi client tru ra dung
			// bang may chu. Them mot goi nua chi la nhan doi luu luong o dung
			// cho dong nguoi nhat cua game.
			skillId := byte(nguoi.chieuDangChon.SkillId)
			k.phatChoCaKhuVuc(goiDanhNguoiChoi(
				int32(nguoi.nv.ID), skillId, int32(boss.Id), satThat, daChet, chiMang,
			))
			return
		}

		// Tim nguoi BI DANH trong khu vuc theo idMucTieu. Khong tim thay (da
		// roi khu vuc giua chung, hoac client gui id cu/gia) -> im lang bo
		// qua, dung tinh than voi suKienDanhQuai khi mobId khong khop con
		// quai nao trong khu vuc.
		mucTieu, coMT := k.nguoiTrong[sk.idMucTieu]
		if !coMT {
			return
		}

		bayGio := k.dongHo() // Task 8/11: dong ho tiem duoc, KHONG goi time.Now() thang.
		ketQua := KiemDanhHopLe(benDanhTuNguoi(nguoi), ViTriCoMau{
			X:         int16(mucTieu.nv.ViTriX),
			Y:         int16(mucTieu.nv.ViTriY),
			HpHienTai: mucTieu.nv.HpHienTai,
			// LucHoiSinh: doc THANG tu nguoiTrongKhuVuc.lucHoiSinh cua MUC
			// TIEU (dat trong case suKienHoiSinh o duoi) — day chinh la moc
			// noi day ma chu thich truong lucHoiSinh (o dau file nay) da bao
			// truoc: "Task 11 (PvP) se doc truong nay qua ViTriCoMau.LucHoiSinh
			// khi mucTieu la nguoi choi". Kich hoat bat tu 1500ms sau hoi
			// sinh (dac ta §1.7) cho HUONG nguoi-danh-nguoi.
			LucHoiSinh: mucTieu.lucHoiSinh,
		}, bayGio)
		if !ketQua.ChoPhep {
			// 4 kiem dau (bat tu/chakra/hoi chieu/da chet) hong -> IM LANG,
			// KHONG phat goi gi ca (dac ta §1.9 chi noi rieng ve vuot tam,
			// khong ap dung cho 4 kiem nay) — dung tinh than voi
			// suKienDanhQuai.
			return
		}

		// Cap nhat lanDanhCuoi CHI KHI duoc phep danh (dung ly do voi case
		// suKienDanhQuai: mot don bi 4 kiem dau tu choi khong duoc tinh la
		// "vua dung chieu").
		nguoi.lanDanhCuoi = bayGio

		// chiMang: giong het case suKienDanhQuai — KH6 CHUA chot ty le/cong
		// thuc chi mang that (khong dac ta nao, khong task truoc nao dinh
		// nghia) -> luon false, lua chon an toan.
		chiMang := false

		var satThuong int32
		if !ketQua.Truot {
			lang := TimLang(nguoi.nv.Lang)
			var satThuongGoc int64
			if lang != nil {
				satThuongGoc = int64(lang.SatThuongKhoiDau)
			}
			// KHAC case suKienDanhQuai (muc tieu la quai, LUON dung
			// heQuaiChuaXacDinh = trung tinh tuyet doi): o day CA HAI PHIA la
			// nguoi choi that, nen dung he nguyen to THAT cua ca hai lang —
			// dung dac ta §4.2 va dung y TestHeNguyenToNguoiChoiAnhXaDungLang
			// (khuvuc_tichhop_test.go) da chuan bi san "cho cac Task sau (vd
			// Task 11 nguoi-vs-nguoi, ca hai phia deu co he nguyen to that)".
			heNguoiDanh := heNguyenToNguoiChoi(nguoi.nv.Lang)
			heMucTieu := heNguyenToNguoiChoi(mucTieu.nv.Lang)
			satThuong = TinhSatThuong(satThuongGoc, nguoi.chieuDangChon, heNguoiDanh, heMucTieu, chiMang)
		}
		// Truot=true (vuot tam): satThuong giu nguyen 0 (dac ta §1.9) — KHONG
		// goi TinhSatThuong, vi don nay khong thuc su cham muc tieu.

		// EP TUONG MINH: vuot tam thi KHONG duoc chi mang — BAN GIAO Task 3
		// (chotkiem.go:75-81, brief Task 11 nhac lai ro rang): "chiMang=false
		// khi Truot=true", vi client CHI hien chu "miss" khi chiMang==false.
		// Dat DONG NAY o day, SAU moi cho co the gan/tinh chiMang, dung y het
		// vi tri/ly do cua dong tuong tu trong case suKienDanhQuai.
		if ketQua.Truot {
			chiMang = false
		}

		// §1.2: sat thuong khong duoc VUOT mau con lai cua muc tieu. LUU Y:
		// KHONG co nhanh §1.3 o day ("quai dang day mau khong chet boi mot
		// don") — dac ta ghi RO quy tac do CHI danh cho QUAI ("Quái đang đầy
		// máu không chết bởi một đòn", §1 diem 3), KHONG danh cho nguoi choi.
		// Nguoi choi CO THE chet dung MOT don du dang day mau tuyet doi — day
		// la diem khac biet CO CHU DICH so voi case suKienDanhQuai, co test
		// rieng khang dinh (TestSuKienDanhNguoiChoiDayMauVanCoTheChetMotDon).
		if satThuong > mucTieu.nv.HpHienTai {
			satThuong = mucTieu.nv.HpHienTai
		}

		mucTieu.nv.HpHienTai -= satThuong
		if mucTieu.nv.HpHienTai < 0 { // phong thu: ve ly thuyet khong xay ra sau kep tren
			mucTieu.nv.HpHienTai = 0
		}
		daChet := mucTieu.nv.HpHienTai <= 0

		// dac ta §3.2: skillId la BYTE trong goi -60 xuong (short trong catalog).
		skillId := byte(nguoi.chieuDangChon.SkillId)
		// -60: CA KHU VUC (dac ta §3.3 cot "Gui cho": "Ca khu vuc") — ke ca A
		// (nguoi danh) VA B (nguoi bi danh), dung tinh than voi 54/-9/-12 o
		// Task 6/7: nguoi gay ra su kien cung nhan lai hoat anh/ket qua.
		k.phatChoCaKhuVuc(goiDanhNguoiChoi(
			int32(nguoi.nv.ID), skillId, int32(mucTieu.nv.ID), satThuong, daChet, chiMang,
		))

		if !ketQua.Truot {
			// -30/9: goi -60 (dac ta §3.3) KHONG mang truong hpConLai nao
			// (khac han -9/-10 da co san hpConLai) -> NHUNG NGUOI KHAC (KHONG
			// PHAI chinh B — dac ta §3.4 cot "Gui cho" cua -30/9 la "Người
			// khác", KHAC voi -30/15 la "Ca khu vuc") can goi rieng nay de ve
			// lai dung thanh mau cua B. Chi gui khi THUC SU co sat thuong
			// (khong Truot) — dac ta §6.4 "moi lan mau nguoi choi DOI vi bi
			// danh": vuot tam thi mau KHONG doi (satThuong=0, giu nguyen o
			// tren), khong co gi de cap nhat.
			k.phatChoNguoiKhac(mucTieu.nv.ID, goiCapNhatMau(
				int32(mucTieu.nv.ID), mucTieu.nv.HpHienTai, mucTieu.nv.HpToiDa,
			))
		}

		if daChet {
			// Ke hoach 6 Task 13 (dac ta §3.4/§6.4), TAI DUNG nguyen ham
			// xuLyNguoiChoiChet (dung dinh nghia o duoi, aiQuai cung goi ham
			// nay): -17 rieng chinh chu (B) roi -8 cho nguoi khac — dung THU
			// TU "dam sat thuong truoc, bao chet sau" cua dac ta §6.4.
			k.xuLyNguoiChoiChet(mucTieu)
		}

	case suKienHoiSinh:
		// Ke hoach 6 Task 13 (dac ta §3.4/§6.4): client gui -16 xin hoi sinh
		// sau khi chet.
		nguoi, co := k.nguoiTrong[sk.idNhanVat]
		if !co {
			return
		}
		if nguoi.nv.HpHienTai > 0 {
			// CHUA chet (hoac da hoi sinh roi) -> bo qua im lang, tranh hoi
			// sinh "khong" mot nguoi dang song. Day CUNG la lop chan chong
			// gui -16 hai lan lien tiep (client bam nut kep, hoac goi den
			// tre): lan dau xu ly xong HpHienTai da > 0 nen lan hai tu nhien
			// roi vao nhanh nay, khong can them bo dem thoi gian rieng.
			return
		}

		bayGio := k.dongHo() // Task 8/13: dong ho tiem duoc, KHONG goi time.Now() thang.
		nguoi.nv.HpHienTai = nguoi.nv.HpToiDa
		nguoi.nv.ChakraHienTai = nguoi.nv.ChakraToiDa
		// Chot bat tu 1500ms sau hoi sinh (dac ta §1.7) — Task 11 (PvP) se
		// doc truong nay qua ViTriCoMau.LucHoiSinh khi mucTieu la nguoi choi.
		nguoi.lucHoiSinh = bayGio

		// -16 RIENG chinh chu TRUOC, roi -30/15 CA KHU VUC — dung thu tu
		// chuoi cua dac ta §6.4.
		if err := nguoi.phien.XepGui(goiHoiSinhChinhChu()); err != nil {
			slog.Warn("gui -16 (xac nhan hoi sinh) that bai", "id", sk.idNhanVat, "loi", err)
			k.loaiNguoiLoiGui(sk.idNhanVat)
			return
		}
		k.phatChoCaKhuVuc(goiAiDoHoiSinh(
			int32(nguoi.nv.ID), nguoi.nv.HpHienTai, nguoi.nv.ChakraHienTai,
			int16(nguoi.nv.ViTriX), int16(nguoi.nv.ViTriY),
		))

	case suKienXinDanhSachQuai:
		// SAO GIA TRI tung con quai — KHONG duoc tra thang k.quai (mang con
		// tro goc) hay tung con tro trong do ra ngoai, vi ben nhan chay tren
		// luong khac va co the giu con tro do lau hon nhip actor tiep theo
		// (dac ta §1.6). traLoi co buffer 1 nen gui khong bao gio block du
		// ben goi da het gio cho va bo di.
		sk.traLoi <- k.saoQuai()

	case suKienXinAnhChup:
		// Ke hoach 11 Task 2: SAO GIA TRI tung nguoi (*nguoi.nv, KHONG PHAI
		// nguoi.nv) — dung khuon voi case suKienXinDanhSachQuai o tren, doi
		// "con tro moi toi ban sao" (*Quai) thanh "gia tri thang" (store.
		// NhanVat) vi ly do da giai thich o chu thich struct suKienXinAnhChup:
		// nguoiTrongKhuVuc.nv la con tro dung CHUNG voi luong phien (khac han
		// *Quai chi minh actor giu), tuyet doi khong duoc lo ra ngoai duoi
		// bat ky hinh thuc con tro nao — *nguoi.nv (dereference) luon tao MOT
		// GIA TRI MOI doc lap o day, gan vao slice/struct la sao chep GIA TRI
		// that su (store.NhanVat khong co truong con tro/slice/map nao ben
		// trong nen khong con duong nao de "vo tinh" con lai tham chieu).
		//
		// traLoi co buffer 1 nen gui khong bao gio block du ben goi da het
		// gio cho va bo di — dung tinh than voi suKienXinDanhSachQuai.
		sk.traLoi <- k.saoNguoiChoi(sk.idNhanVat)
	}
}

// phatChoNguoiKhac gui goi tin `m` cho MOI nguoi trong khu vuc TRU nguoi co id `truId`.
func (k *KhuVuc) phatChoNguoiKhac(truId int64, m *protocol.Message) {
	var loi []int64
	for id, nguoi := range k.nguoiTrong {
		if id == truId {
			continue
		}
		if err := nguoi.phien.XepGui(m); err != nil {
			slog.Warn("phat goi tin that bai", "id", id, "loi", err)
			loi = append(loi, id)
		}
	}
	k.loaiNhieuNguoiLoiGui(loi)
}

// phatChoCaKhuVuc gui goi tin `m` cho TOAN BO nguoi dang trong khu vuc, KE CA
// nguoi vua gay ra su kien (khac phatChoNguoiKhac o tren, loai tru nguoi gay
// su kien). Cac goi chien dau (54/-9/-12, dac ta §3.3 cot "Gui cho: Ca khu
// vuc") phai den ca nguoi danh de chinh ho cung thay hoat anh/ket qua don
// danh cua minh — khac han goi di chuyen (-7) noi client tu ve chuyen dong
// cua chinh minh nen khong can nhan lai.
func (k *KhuVuc) phatChoCaKhuVuc(m *protocol.Message) {
	var loi []int64
	for id, nguoi := range k.nguoiTrong {
		if err := nguoi.phien.XepGui(m); err != nil {
			slog.Warn("phat goi tin (ca khu vuc) that bai", "id", id, "loi", err)
			loi = append(loi, id)
		}
	}
	k.loaiNhieuNguoiLoiGui(loi)
}

// themBoss dua MOT con boss vao khu vuc va bao cho moi nguoi dang co mat
// (KH6b Task 7). Day la duong VAO DUY NHAT cua k.boss — bo dieu phoi lich
// (Task 13, ngoai actor) gui lenh sinh qua hop thu, may trang thai (Task 9)
// goi ham nay TU TRONG vongLap.
//
// CHAY TRONG ACTOR: chi dung phien.XepGui qua phatChoCaKhuVuc (khong Gui()),
// khong cho doi ai, khong sinh goroutine — dung ky luat cua moi ham khac trong
// file nay.
//
// b == nil -> khong lam gi (KHONG panic, KHONG them mot o nil vao danh sach
// de moi vong lap sau nay phai kiem lai).
//
// Boss CHUA hien (NGHI) van duoc THEM vao danh sach nhung KHONG phat goi nao:
// phat -5 luc do la lam boss hien som, sai vong doi. Goi nil (du lieu pha hong)
// cung khong phat — goiXuatHienBoss da ghi log ly do; boss van nam trong danh
// sach de may trang thai co co hoi don no o cac buoc sau.
//
// ⚠️ TU TASK 9: nap mot con o trang thai NGHI khong con y nghia gi nua — nhip
// tim ke tiep se GO no ra ngay (xem case BossNghi trong mayTrangThaiBoss). Ly
// do: CHAN-7 cam actor dem nguoc GiayNghi de tu danh thuc boss, nen mot con
// NGHI nam trong actor khong co duong nao tro lai the gioi. Muon boss hien lat
// nua thi bo dieu phoi giu no lai BEN NGOAI roi sinh dung luc, dung nap som.
func (k *KhuVuc) themBoss(b *Boss) {
	if b == nil {
		return
	}
	k.boss = append(k.boss, b)
	if !b.TrangThai.DangHien() {
		return
	}
	if m := goiXuatHienBoss(b); m != nil {
		k.phatChoCaKhuVuc(m)
	}
	// Thong bao toan server SAU goi -5 (KH6b Task 12, thiet ke §5 buoc HIEN).
	// Thu tu nay co y va KHONG doi duoc y nghia dau: -5 di thang toi nguoi
	// trong khu vuc ngay lap tuc, con thong bao 93 phai xep hang qua loa (nhip
	// 1,5s) nen no LUON toi sau — dat truoc hay sau o day khong doi dieu do.
	// Viet sau cho khop voi thu tu nguoi choi cam nhan: thay boss truoc, doc
	// dong chu sau.
	//
	// CO Y KHONG kiem "goi -5 dung duoc khong": mot con boss co du lieu pha
	// hong (khong ve duoc) van la mot su kien da xay ra; nhung hai viec do doc
	// lap nhau, gan chung lai se lam mot loi hien thi am tham nuot mat ca
	// thong bao.
	k.thongBaoBossHien(b)
}

// thongBaoBossHien bao toan may chu rang MOT boss lon vua xuat hien (thiet ke
// §5, buoc HIEN: "phat -5, announce 93 neu boss lon").
//
// Tach thanh ham RIENG (khong nhet thang vao themBoss) vi day la ham ma may
// trang thai vong doi (Task 9) se goi lai o moi lan BIEN HINH: pha ke tiep
// quay ve HIEN voi CUNG mot id nhung ten/hinh dang khac — do la mot lan "xuat
// hien" nua theo dung nghia nguoi choi thay.
//
// CHI bao khi trang thai dung la BossHien: DangHien() con gom CHAO/HOAT_DONG/
// TROI, ma mot con boss duoc nap vao khu vuc luc dang danh nhau (hoac dang la
// cai xac) thi khong phai "vua xuat hien".
//
// LaBossLon nam o TUNG PHA, khong phai o ca con boss — day la nut chinh do on
// duy nhat cua tinh nang nay: tac gia noi dung bat co o pha nao thi ca may chu
// nghe pha do, cac pha con lai bien hinh im lang. Vi vay o day KHONG can them
// bat ky lop "gop tin/bo trung" nao.
func (k *KhuVuc) thongBaoBossHien(b *Boss) {
	if b == nil || b.TrangThai != BossHien {
		return
	}
	pha := b.PhaHienTaiData()
	if pha == nil || !pha.LaBossLon {
		return
	}
	k.xepThongBao(cauBossXuatHien(pha.Ten))
}

// xepThongBao la THU DUY NHAT actor duoc goi de bao mot tin ra toan may chu
// (vd nhipBoss: "Boss Cuu Vi da xuat hien"). No CHI xep mot yeu cau vao hang
// doi cua loa roi tra ve NGAY — khong cham mang, khong cho ai, khong dung
// guiDieuKhien (bay §6.10b: cho 100ms tren hop thu ma chinh vongLap dang ban
// = tu khoa chac chan). Viec phat that do goroutine cua loa lam.
//
// loa == nil -> khong lam gi, KHONG panic: khu vuc chua duoc noi day (moi
// MoKhuVuc trong test, hoac may chu chua chay GanLoa) van phai chay binh
// thuong.
//
// Moc thoi gian lay tu k.dongHo() (khong phai time.Now) theo dung rang buoc
// chung cua actor — MocTao chi dung de do do tre hang doi trong log.
func (k *KhuVuc) xepThongBao(noiDung string) {
	loa := k.loa.Load()
	if loa == nil {
		return
	}
	loa.Xep(YeuCauThongBao{NoiDung: noiDung, MocTao: k.dongHo()})
}

// tranSoNguoiDongGop la SO NGUOI toi da mot con boss ghi ten vao so dong gop
// cua MOT pha (bay "map dong gop phinh vo han", KH6b Task 12).
//
// Chan tren THAT SU cua so nay la buoc xoa sach moi lan boss chet (xem
// bossChet) — mot pha chi song vai phut nen so nguoi kip danh no la huu han.
// Tran nay chi la cai chan cuoi cho truong hop benh hoan: mot con boss vi loi
// du lieu ma khong bao gio chet, nam trong khu vuc dong nguoi hang gio. Khong
// co no, so dong gop lon len mai theo SO NGUOI TUNG DI QUA, khong bao gio thu
// nho lai.
//
// 512 chon theo tinh than "cao hon moi con so that co the xay ra, thap hon
// muc gay hai": mot khu vuc khong bao gio co 512 nguoi cung danh mot con boss
// trong mot pha, nen tran nay khong bao gio cham toi trong van hanh binh
// thuong — no chi ton tai de con so co GIOI HAN.
const tranSoNguoiDongGop = 512

// bossNhanSatThuong ap MOT don sat thuong len boss: ghi so dong gop, tru mau,
// va phat hien chet (thiet ke §3 "cong satThuongTheoNguoi" + §5 "chet: hp<=0
// trong luc ap sat thuong").
//
// DAY LA DUONG DUY NHAT duoc tru mau boss. Task 8 (dinh tuyen goi -60 len
// boss) chiu trach nhiem TINH ra con so `sat` theo luat chien dau roi goi ham
// nay — no KHONG duoc tu `b.Hp -= ...`, vi lam vay la bo qua ca so dong gop
// lan chot chong-chia-hai-lan ben duoi.
//
// Tra ve true khi don nay LAM BOSS CHET (duong thuong/thong bao da chay xong
// ngay trong loi goi) — nguoi goi dung gia tri do de chon goi tin phat ra.
//
// CHAY TRONG ACTOR: khong cham mang, khong cham CSDL, khong cho ai. Cai duy
// nhat "ra ngoai" la mot cu k.xepThongBao (chi day vao hang doi loa).
func (k *KhuVuc) bossNhanSatThuong(b *Boss, idNguoi int64, sat int64, bayGio time.Time) bool {
	if b == nil || sat <= 0 {
		return false
	}
	// CHOT CHONG CHIA THUONG HAI LAN (quan trong nhat ham nay): chi con boss
	// DANG SONG va DANG TREN MAN HINH moi an don. BossTroi la cai XAC — no van
	// DangHien() (client con thay) nhung da tra thuong xong roi; cho danh tiep
	// nghia la hai nguoi bam cung mot nhip deu thay hp<=0 va deu goi bossChet.
	// BossNghi/BossRaDi thi khong co mat trong the gioi de ma danh.
	if !b.TrangThai.DangHien() || b.TrangThai == BossTroi || b.Hp <= 0 {
		return false
	}

	// Dong gop duoc KEP o so mau CON LAI, khong phai so sat thuong vung ra.
	// Neu tinh nguyen `sat`, mot nguoi den phut chot giang mot don 1 ty vao
	// con boss con 5 mau se chiem ~100% ti le va cuop trang phan cua nhung
	// nguoi da danh ca tran. Dong gop phai do "da lam boss mat bao nhieu mau".
	gop := sat
	if gop > int64(b.Hp) {
		gop = int64(b.Hp)
	}
	// Tran so nguoi: nguoi DA CO trong so van cong don binh thuong, chi chan
	// nguoi MOI. Chan ca nguoi cu thi ai danh tu dau se bi dong bang dong gop
	// giua chung tran — te hon han viec bo sot vai nguoi den muon.
	if _, daCo := b.SatThuongTheoNguoi[idNguoi]; daCo || len(b.SatThuongTheoNguoi) < tranSoNguoiDongGop {
		b.CongSatThuong(idNguoi, gop)
	} else {
		slog.Warn("so dong gop cua boss da day, KHONG ghi them nguoi moi",
			"idBoss", b.Id, "idNguoi", idNguoi, "tran", tranSoNguoiDongGop)
	}

	// KepHp (khong phai `b.Hp -= int32(sat)`): tinh bang int64 roi kep MOT LAN
	// o buoc cuoi — tru thang tren int32 se TRAN truoc khi kip kep (bay #4
	// thiet ke §7).
	b.Hp = KepHp(int64(b.Hp) - sat)
	if b.Hp > 0 {
		return false
	}
	k.bossChet(b, idNguoi, bayGio)
	return true
}

// bossChet chay TOAN BO nghi thuc "boss vua chet" theo dung thu tu thiet ke §5
// dong "Chet (hp<=0 trong luc ap sat thuong)":
//
//	announce 93 (neu boss lon) -> chia thuong theo dong gop -> XOA so dong gop
//	-> chuyen TROI
//
// idNguoiHa la nguoi giang don ket lieu — CHI dung de dien ten vao dong thong
// bao, KHONG duoc uu ai gi them trong phan thuong: thuong chia theo dong gop,
// khong theo ai bam nut cuoi.
//
// BA DIEU BAT BIEN, doc ky truoc khi sua:
//
//  1. CHIA XONG THI XOA, TRONG CUNG MOT HAM. Khong duoc tach buoc xoa ra cho
//     may trang thai lam o RA_DI: giua hai thoi diem do boss van co the an
//     them don (chot o bossNhanSatThuong chan duoc, nhung mot duong goi thu
//     hai trong tuong lai thi khong) — va mot lan chia lai la nhan doi tien
//     phat ra tu khong khi.
//
//  2. CHIA THUONG O MOI LAN PHA CHET, khong doi toi pha cuoi. Thiet ke §5 noi
//     ro TROI xay ra o MOI lan chet, va dong 62 dat "chia thuong" ngay trong
//     buoc do. Doi toi pha cuoi thi so dong gop phai song xuyen nhieu pha —
//     vua phinh to (bay §7), vua bat cong: nguoi danh gan chet pha 1 roi thoat
//     game se an ke phan cua nguoi vat lon voi pha 3.
//
//  3. CHOT "DA CHET ROI THI THOI" (TrangThai da la TROI/RA_DI) la lop bao ve
//     THU HAI, sau chot o bossNhanSatThuong. Giu ca hai co y: chot kia bao ve
//     duong danh, chot nay bao ve chinh ham nay truoc moi nguoi goi tuong lai
//     (may trang thai, lenh admin, hanh dong dac biet).
//
// CHAY TRONG ACTOR: khong cham mang/CSDL/khong cho ai — phan thuong ghi thang
// vao nv (actor so huu con tro do), thong bao di qua hang doi loa.
func (k *KhuVuc) bossChet(b *Boss, idNguoiHa int64, bayGio time.Time) {
	if b == nil {
		return
	}
	if b.TrangThai == BossTroi || b.TrangThai == BossRaDi {
		return // xem bat bien 3
	}
	b.Hp = 0

	pha := b.PhaHienTaiData() // nil khi du lieu pha hong — moi buoc duoi deu chiu duoc nil
	if pha != nil && pha.LaBossLon {
		k.xepThongBao(cauBossBiHa(k.tenNguoiTrongKhuVuc(idNguoiHa), pha.Ten))
	}

	k.traThuongBoss(b, pha)
	// XOA NGAY sau khi chia — xem bat bien 1. Gan nil (khong phai xoa tung
	// khoa) de tra luon bo nho map ve cho bo don rac: mot pha dong nguoi co
	// the de lai hang tram o trong, ma pha sau se tu tao map moi khi can.
	b.SatThuongTheoNguoi = nil

	b.TrangThai = BossTroi
	b.MocThoiGian = bayGio
}

// traThuongBoss chia DiemThuong cua pha cho nhung nguoi da danh, theo ti le
// dong gop, roi CONG vao nv.KinhNghiem cua tung nguoi.
//
// VI SAO PHAN THUONG LA nv.KinhNghiem (thiet ke §2 quyet dinh 3 chi noi
// "diem san boss", khong chot truong nao):
//   - No LA MOT BO DEM CHI TANG da co san trong store.NhanVat, dung y nghia
//     "phan thuong cho cong danh quai/boss" — khong phai bia ra ngu nghia moi.
//   - No DI XUONG CLIENT SAN ROI: chisonhanvat.go ghi nv.KinhNghiem vao truong
//     cTiemNang cua goi chi so, nen nguoi choi THAY con so tang ma khong can
//     them mot byte giao thuc nao.
//   - Ba truong con lai deu tro hon: Ngoc la tien nap (tu dong phat tien nap
//     tu boss la mot lo hong kinh te, khong phai mot tinh nang); Vang la voi
//     nuoc lam phat khi KH7 chua co gi de mua; SucManh la diem xep hang hien
//     thi cho web, khong duoc TaoNhanVat khoi tao va da bi ghi chu trong
//     chinh file nay la KHONG dung lam chi so gameplay.
//
// VI SAO KHONG PHAI GHI CSDL O DAY (yeu cau "thuong phai ben"): ghi tu actor
// la dung cai loi §6.10b. Khong can: nv la con tro actor DANG SO HUU, va
// duong luu san co tu chup lai — BoLuuDinhKy 60 giay (luudinhky.go) va
// KhiDong luc nguoi choi thoat (login.go) deu sao *nguoi.nv thanh gia tri roi
// day cho ThoGhi. Ghi vao nv la du de ben; cai gia toi da la 60 giay tien do
// neu may chu sap dot ngot — dung muc da chot cho MOI chi so khac trong game.
//
// NGUOI DA ROI KHU VUC BI BO QUA — cai gia da can nhac va chap nhan:
//   - Ho khong con trong k.nguoiTrong, nghia la actor nay KHONG con so huu con
//     tro nv cua ho (ho dang o khu vuc khac, hoac phien da dong). Ghi vao mot
//     con tro minh khong so huu la dung bay §6.9 (tranh chap du lieu that).
//   - Duong vong "day mot ban store.NhanVat cu cho ThoGhi" con te hon: ThoGhi
//     ghi de CA DONG, nen ban chup cu se KEO LUI moi thu ho lam tu luc roi di
//     (vi tri, mau, vang) — mat du lieu that de doi lay mot khoan thuong.
//   - Cach dung la mot SO NO phan thuong ben (bang CSDL rieng + duong nhan lai
//     luc dang nhap). Do la mot TINH NANG RIENG — can bang moi, migration moi,
//     duong nhan thuong luc dang nhap — khong phai task nay. CHUA CO AI LAM;
//     neu tinh cong bang nay quan trong thi phai mo mot task rieng cho no.
//
// Vi vay o day chi log muc canh bao KEM SO DIEM — de neu chuyen do xay ra
// nhieu that thi co so lieu ma quyet dinh, thay vi mat im lang.
//
// pha nhan qua THAM SO (khong tu goi PhaHienTaiData lai) de chac chan chia
// theo DUNG pha ma bossChet vua thong bao — mot cho duy nhat doc pha.
func (k *KhuVuc) traThuongBoss(b *Boss, pha *PhaBoss) {
	if b == nil || pha == nil || pha.DiemThuong <= 0 {
		return
	}
	phan := ChiaThuongTheoDongGop(b.SatThuongTheoNguoi, pha.DiemThuong)
	for id, diem := range phan {
		nguoi, co := k.nguoiTrong[id]
		if !co || nguoi == nil || nguoi.nv == nil {
			slog.Warn("nguoi dong gop khong con trong khu vuc luc boss chet — BO phan thuong",
				"idBoss", b.Id, "idNguoi", id, "diemBoQua", diem)
			continue
		}
		nguoi.nv.KinhNghiem = congBaoHoa(nguoi.nv.KinhNghiem, diem)
	}
}

// tenNguoiTrongKhuVuc tra ten hien thi cua mot nguoi choi dang o trong khu
// vuc, hoac tenAnDanh neu khong tra duoc.
//
// Ba nhanh cung tra ten thay the, va deu la truong hop THAT: id khong con
// trong khu vuc (vua mat ket noi ngay sau don ket lieu), ban ghi nhan vat nil
// (chi xay ra neu co loi lap trinh o cho khac — nhung mot cu nil-deref o day
// se giet nhip actor, khong dang), va ten rong trong CSDL.
func (k *KhuVuc) tenNguoiTrongKhuVuc(id int64) string {
	nguoi, co := k.nguoiTrong[id]
	if !co || nguoi == nil || nguoi.nv == nil || nguoi.nv.Ten == "" {
		return tenAnDanh
	}
	return nguoi.nv.Ten
}

func (k *KhuVuc) loaiNguoiLoiGui(id int64) {
	k.loaiNhieuNguoiLoiGui([]int64{id})
}

// loaiNhieuNguoiLoiGui chay trong actor: xoa state TRUOC, roi moi DongCham va
// bao -6 cho nhung nguoi con lai. Neu mot nguoi con lai cung day queue trong
// luc nhan -6, lap tiep cho toi khi khong con session hong; khong goroutine
// phu, khong ghi socket.
func (k *KhuVuc) loaiNhieuNguoiLoiGui(ids []int64) {
	for len(ids) > 0 {
		daLoai := make([]int64, 0, len(ids))
		for _, id := range ids {
			nguoi, co := k.nguoiTrong[id]
			if !co {
				continue
			}
			k.anhChupDaRoi[id] = *nguoi.nv
			delete(k.nguoiTrong, id)
			nguoi.phien.DongCham()
			daLoai = append(daLoai, id)
		}
		if len(daLoai) == 0 {
			return
		}

		tiep := make(map[int64]struct{})
		for _, idDaLoai := range daLoai {
			m := goiRoiKhuVuc(idDaLoai)
			for idConLai, nguoi := range k.nguoiTrong {
				if err := nguoi.phien.XepGui(m); err != nil {
					slog.Warn("gui -6 xoa client cham that bai", "id", idConLai, "loi", err)
					tiep[idConLai] = struct{}{}
				}
			}
		}
		ids = ids[:0]
		for id := range tiep {
			ids = append(ids, id)
		}
	}
}

// Hai hang so hanh vi quai DA CHOT (dac ta §5.2, nguon Mob.java:347 va
// Mob.java:351-355; brief Task 8 "Hang so da chot") — nguon su that duy nhat
// cho hai con so nay, dung trong nhipTim duoi day. Dat kieu hang so (khong
// phai bien) de kiem thu dot bien (brief Step 5: "doi 5s thanh 50s") sua
// TRUC TIEP gia tri o day roi chay lai test bien, khong can sua noi khac.
const (
	// thoiGianHoiSinhQuai: quai chet roi cho DUNG 5 giay thi hoi sinh (dac ta
	// §5.2 "Hoi sinh sau khi chet — 5 giay").
	thoiGianHoiSinhQuai = 5 * time.Second

	// khoangHoiMauQuai: quai con song nhung chua day mau thi cu moi 30 giay
	// hoi them hpToiDa/10 (dac ta §5.2 "Hoi mau khi chua day — moi 30 giay,
	// hoi hpMax/10").
	khoangHoiMauQuai = 30 * time.Second
)

// nhipTim chay moi 100ms (goi tu vongLap) — Ke hoach 6 Task 8 dung mooc nay
// (de rong tu Ke hoach 5) cho ba viec cua tung con quai trong k.quai: HOI
// SINH neu da chet du 5 giay (hoiSinhQuai), HOI MAU dinh ky 30 giay neu con
// song nhung chua day mau (hoiMauQuai), va AI hai tang uu tien danh nguoi
// choi (aiQuai, Task 9). Day la NOI DUY NHAT ngoai case suKienDanhQuai cham
// vao trang thai quai/nguoiTrong — dung dac ta §1.6 (chi luong vongLap duoc
// cham trang thai quai va trang thai chien dau).
//
// bayGio LAY MOT LAN o dau ham (khong goi k.dongHo() lai trong vong lap) de
// TOAN BO quai trong cung mot nhip duoc xet tren CUNG mot moc thoi gian —
// tranh truong hop hai con quai o hai nhanh khac nhau cua cung mot lenh gio
// (vd mot con vua qua nguong o dau vong lap, con kia chua qua vi dongHo() da
// nhich len vai micro-giay giua hai lan goi — khong xay ra trong that te vi
// nhip 100ms qua tho, nhung day la thoi quen dung, tranh phu thuoc vao thu tu
// lap qua slice).
//
// Task 9: nhipTim da ~55 dong tu Task 8 (brief Task 9 de nghi tach ham khi
// them logic) nen ba viec tren duoc tach thanh ba ham rieng (hoiSinhQuai/
// hoiMauQuai/aiQuai) — chi vongLap goi chung, van chi luong nay cham trang
// thai (khong doi ky luat §1.6, chi doi cach to chuc code).
func (k *KhuVuc) nhipTim() {
	bayGio := k.dongHo()

	// Task 9: snapshot nguoi choi thanh SLICE MOT LAN cho ca nhip nay, dung
	// chung cho AI cua MOI con quai — tranh xay map->slice lai N lan (N = so
	// quai). Con tro TRO THANG vao k.nguoiTrong (KHONG sao gia tri) vi ham
	// nay chay TRONG vongLap nen doc truc tiep la an toan (dac ta §1.6); AI
	// (aiQuai) cung CHI DOC qua slice nay, khong ghi.
	var nguoi []*nguoiTrongKhuVuc
	if len(k.nguoiTrong) > 0 {
		nguoi = make([]*nguoiTrongKhuVuc, 0, len(k.nguoiTrong))
		for _, n := range k.nguoiTrong {
			nguoi = append(nguoi, n)
		}
	}

	for _, q := range k.quai {
		if q.Hp <= 0 {
			// Dang chet (cho hoi sinh, hoac vua hoi sinh xong NGAY trong
			// hoiSinhQuai() o duoi). Ca hai truong hop deu continue — quai
			// dang chet khong the hoi mau/AI, va quai vua hoi sinh chac
			// chan Hp==HpToiDa nen nhanh hoi mau se tu nhien khong lam gi;
			// AI cho quai vua hoi sinh doi sang tick SAU (khong danh ngay
			// trong chinh tick hoi sinh) — don gian hoa, khong dac ta nao
			// doi hoi phai danh NGAY trong cung tick hoi sinh.
			k.hoiSinhQuai(q, bayGio)
			continue
		}

		k.hoiMauQuai(q, bayGio)
		k.aiQuai(q, nguoi, bayGio)
	}

	// KH6b Task 9: may trang thai vong doi boss chay TRONG CUNG mot nhip, o
	// mot VONG LAP RIENG ngay sau dan quai (ly do tach vong: xem nhipBoss).
	// Boss cap nhat trong nhip 100ms cua khu vuc, KHONG luong 150ms rieng nhu
	// ban Java — bay #5 cua thiet ke §7: hai luong cung ghi viTri la tranh chap
	// du lieu that.
	k.nhipBoss(nguoi, bayGio)
}

// hoiSinhQuai xu ly nhanh "quai dang cho hoi sinh" cua nhipTim — GOI KHI
// q.Hp <= 0. Neu CHUA du 5 giay tu LucChet thi khong lam gi (dac ta §5.2,
// brief Task 8 "sau 4.9s CHUA hoi sinh"); du/qua 5 giay thi hoi sinh DAY DU
// mau + phat DUNG MOT goi -13 (khong phat gi them sau do, dac ta §0.1 tuyet
// doi cam goi 87 theo sau). Tach nguyen ven logic Task 8, chi doi vi tri.
func (k *KhuVuc) hoiSinhQuai(q *Quai, bayGio time.Time) {
	if bayGio.Sub(q.LucChet) < thoiGianHoiSinhQuai {
		return
	}

	// Du/qua 5 giay -> hoi sinh DAY DU mau (dac ta §5.4 diem 5: ban Java hoi
	// day mau ngau nhien qua loi ham lvMob() khi dung goi -13 — ban Go LAM
	// VIEC NAY TUONG MINH o day, khong an sau ham doc). Ghi LucHoiSinh = bayGio
	// de kich hoat bat tu 1500ms (dac ta §1.7) qua ViTriCoMau.LucHoiSinh (xem
	// case suKienDanhQuai). Ghi lai LanHoiMauCuoi = bayGio de tranh hoi mau
	// "thua" ngay tick sau — quai vua hoi sinh da DAY MAU, khong can cong don
	// thoi gian tu TRUOC luc chet vao chu ky hoi mau moi.
	q.Hp = q.HpToiDa
	q.LucHoiSinh = bayGio
	q.LanHoiMauCuoi = bayGio

	// CHI phat DUNG MOT goi -13, KHONG phat gi them sau do — dac ta §0.1
	// tuyet doi cam goi 87 theo sau (client hieu nham thanh "quai bi gio
	// thoi", quai mau >= 16.777.216 se dung im khong di chuyen nua).
	// goiQuaiHoiSinh (goichiendau.go, Task 6) da tu no LUON ghi sys=0/
	// levelBoss=0, khong nhan tham so nao co the vo tinh mo duong cho mot
	// gia tri khac 0.
	k.phatChoCaKhuVuc(goiQuaiHoiSinh(q.Id, q.Hp))
}

// hoiMauQuai xu ly nhanh "hoi mau dinh ky" cua nhipTim — GOI KHI q.Hp > 0
// (nguoi goi, nhipTim, da loai truong hop Hp<=0 truoc do). Chi hoi mau khi
// CHUA day va da du 30 giay ke tu lan kiem truoc (dac ta §5.2). Quai vua day
// mau (Hp == HpToiDa) khong vao nhanh if -> LanHoiMauCuoi KHONG bi cap nhat
// khi dang day mau, dung y dinh "chi dem thoi gian tu lan KIEM GAN NHAT",
// khong phai tu luc spawn co dinh. Tach nguyen ven logic Task 8.
func (k *KhuVuc) hoiMauQuai(q *Quai, bayGio time.Time) {
	if q.Hp < q.HpToiDa && bayGio.Sub(q.LanHoiMauCuoi) >= khoangHoiMauQuai {
		q.LanHoiMauCuoi = bayGio

		q.Hp += q.HpToiDa / 10
		if q.Hp > q.HpToiDa { // khong duoc vuot hpToiDa (dac ta §5.2)
			q.Hp = q.HpToiDa
		}
		// KHONG phat goi gi o day: dac ta khong dinh nghia goi nao cho hoi
		// mau dinh ky (khac han goi -9 khi bi danh) — client chi thay mau
		// moi qua mot don danh (-9) hoac lan vao lai khu vuc (-24) sau nay,
		// dung hanh vi "hoi mau tham lang" cua ban goc.
	}
}

// aiQuai la AI hai tang uu tien cua Ke hoach 6 Task 9 (dac ta §5.2) — GOI KHI
// q.Hp > 0, MOI nhip 100ms, cho MOI con quai (goi tu nhipTim, sau hoiMauQuai).
// Dung q.chonMucTieu (ham THUAN, quai.go) de TIM muc tieu + nhip yeu cau; ham
// nay (chay trong vongLap) la noi DUY NHAT quyet dinh co "that su danh" hay
// khong va cham q.LanDanhCuoi — dung dac ta §5.4 diem 4 ("doi nhip danh nhu
// tac dung phu cua ham tim muc tieu la loi, phai tach rieng ham tra ve +
// nguoi goi tu quyet dinh") va brief Task 9 (dot bien bat buoc: "chuyen viec
// dat lai nhip ra ngoai nhanh that su danh -> test phai do").
//
// Ke hoach 6 Task 11 Phan B: bayGio (da co san tu Task 9, truyen xuong tu
// nhipTim) gio con duoc CHUYEN TIEP vao chonMucTieu — day la duong va lo hong
// "quai danh nguoi vua hoi sinh" (dac ta §1.7), xem chu thich chonMucTieu.
func (k *KhuVuc) aiQuai(q *Quai, nguoi []*nguoiTrongKhuVuc, bayGio time.Time) {
	mucTieu, nhip := q.chonMucTieu(nguoi, bayGio)
	if mucTieu == nil {
		// KHONG CO MUC TIEU: TUYET DOI khong dung toi q.LanDanhCuoi o day
		// (dac ta §5.4.2 — bay ban Java Mob.java:386-390 dat lai moc thoi
		// gian NGOAI nhanh "co nguoi choi", khien quai vua thay nguoi phai
		// cho THEM gan mot nhip nua moi duoc danh, thay vi danh NGAY).
		// Return som, KHONG cham trang thai gi ca.
		return
	}
	if bayGio.Sub(q.LanDanhCuoi) < nhip {
		return // con trong nhip hoi cua tang uu tien vua chon
	}

	// THAT SU DANH — CHI o day q.LanDanhCuoi moi duoc dat lai (dung nhanh
	// "that su danh" ma dot bien bat buoc cua brief doi hoi).
	q.LanDanhCuoi = bayGio

	// §1.2: sat thuong khong duoc vuot mau con lai cua nan nhan. SatThuongCoBan
	// la gia tri CO DINH tam thoi (xem chu thich truong Quai.SatThuongCoBan,
	// quai.go) — cong thuc day du §5.3 la viec cua KH8, ngoai pham vi Task 9.
	satThuong := q.SatThuongCoBan
	if satThuong < 0 { // phong thu: cau hinh hong (SatThuongCoBan am) khong lam mau am
		satThuong = 0
	}
	if satThuong > mucTieu.nv.HpHienTai {
		satThuong = mucTieu.nv.HpHienTai
	}

	mucTieu.nv.HpHienTai -= satThuong
	if mucTieu.nv.HpHienTai < 0 { // phong thu: ve ly thuyet khong xay ra sau kep tren
		mucTieu.nv.HpHienTai = 0
	}

	// -11: RIENG nan nhan (dac ta §3.3 cot "Gui cho": "Rieng nan nhan") —
	// KHONG dung phatChoNguoiKhac/phatChoCaKhuVuc (ca hai deu sai doi tuong:
	// mot loai tru nan nhan, mot gui thua cho nguoi khac). Gui THANG qua
	// phien cua chinh nan nhan, giong cach case suKienVao gui -5 rieng cho
	// tung nguoi cu (xem k.xuLy o tren).
	if err := mucTieu.phien.XepGui(goiQuaiDanhTrungMinh(q.Id, satThuong)); err != nil {
		slog.Warn("gui -11 (quai danh trung minh) that bai", "idNanNhan", mucTieu.nv.ID, "loi", err)
		k.loaiNguoiLoiGui(mucTieu.nv.ID)
		return
	}
	// -10: NGUOI KHAC trong khu vuc, khong phai nan nhan (dac ta §3.3 cot
	// "Gui cho": "Nguoi khac trong khu vuc"). phatChoNguoiKhac da co san (Ke
	// hoach 5), dung dung nghia: loai tru id nan nhan.
	k.phatChoNguoiKhac(mucTieu.nv.ID, goiQuaiDanhTrungNguoiKhac(q.Id, int32(mucTieu.nv.ID), mucTieu.nv.HpHienTai))

	// Ke hoach 6 Task 13 (dac ta §3.4/§6.4): don nay VUA HA CHET nan nhan
	// (HpHienTai vua xuong <= 0 o kep tren) -> phat THEM chuoi goi "chet"
	// NGAY SAU -11/-10 (dam sat thuong truoc, roi moi bao chet — dung thu tu
	// §6.4). chonMucTieu (quai.go) da LOC san nguoi HpHienTai<=0 ra khoi CA
	// HAI tang AI (dac ta §5.4.3) nen KHONG bao gio quai nay (hay bat ky
	// quai nao khac cung tick) danh trung mot nguoi DA CHET tu truoc — su
	// kien "vua chet" o day CHAC CHAN la lan dau tien, khong can chot chong
	// phat lai.
	if mucTieu.nv.HpHienTai <= 0 {
		k.xuLyNguoiChoiChet(mucTieu)
	}
}

// xuLyNguoiChoiChet phat chuoi goi bao CHET cua mot nguoi choi (Ke hoach 6
// Task 13, dac ta §3.4/§6.4): -17 RIENG chinh chu, -8 cho NHUNG NGUOI KHAC
// trong khu vuc. GOI KHI nguoi.nv.HpHienTai VUA xuong <= 0 — ham nay KHONG
// tu kiem lai dieu kien do (nguoi goi, hien tai la aiQuai, chiu trach nhiem
// dam bao chi goi dung mot lan tai thoi diem chuyen tu song sang chet).
//
// cPk LUON 0: KH6 chua co he PK (dac ta §0 "Ngoai KH6": "PvP co luat -> KH9"),
// dung tinh than voi cPk=0 cua goiMeInfo (thegioi.go). x/y la VI TRI HIEN TAI
// — Java (PlayerService.hoiSinh) KHONG doi vi tri nguoi choi luc chet hay
// luc hoi sinh, chi phat hoat anh chet/song tai cho (dac ta khong nhac toi
// bat ky dich chuyen nao), nen ham nay doc THANG nv.ViTriX/Y, khong tinh
// toan gi them.
func (k *KhuVuc) xuLyNguoiChoiChet(nguoi *nguoiTrongKhuVuc) {
	x := int16(nguoi.nv.ViTriX)
	y := int16(nguoi.nv.ViTriY)

	if err := nguoi.phien.XepGui(goiNguoiChet(0, x, y)); err != nil {
		slog.Warn("gui -17 (ban da chet) that bai", "id", nguoi.nv.ID, "loi", err)
		k.loaiNguoiLoiGui(nguoi.nv.ID)
		return
	}
	k.phatChoNguoiKhac(nguoi.nv.ID, goiAiDoChet(int32(nguoi.nv.ID), 0, x, y))
}

// ─────────────────────────────────────────────────────────────────────────────
// KH6b Task 9 — MAY TRANG THAI VONG DOI BOSS, chay TRONG actor
//
//	NGHI -> HIEN -> CHAO -> HOAT_DONG -> (chet) -> TROI -> RA_DI -> (roi han)
//
// RANH GIOI CUA MAY NAY (muc CHAN-7 cua ke hoach, doc truoc khi sua bat cu gi):
// no CHI lo chuyen gi xay ra SAU KHI da nhan lenh sinh. No KHONG doc dong ho
// lich, KHONG chon map, KHONG tu quyet sinh lai. Ba cau hoi do thuoc ve bo dieu
// phoi lich (lichboss.go) chay NGOAI actor. GiayNghi trong du lieu pha la thu
// BO DIEU PHOI doc — may trang thai nay khong duoc nhin vao no.
//
// Vi sao ranh gioi do quan trong den muc phai viet lai o day: neu ca hai ben
// cung tra loi "khi nao sinh", ta co hai co che sinh boss chong nhau, va khong
// ai con biet con boss dang dung giua ban do la do ai sinh ra — thu ma khong
// mot dong log nao noi cho ta biet.
// ─────────────────────────────────────────────────────────────────────────────

// Ba hang so nhip cua vong doi boss. Ca ba deu la BOI SO CUA 100ms (chu ky
// nhipTim): mot con so le nhu 1450ms tren thuc te van la 1500ms, tuc la mot con
// so noi doi.
const (
	// treoMoiKyTuThoai + treoThoaiToiDa: mot cau thoai treo tren dau boss bao
	// lau truoc khi cau ke tiep duoc noi.
	//
	// BANG CHUNG (khong bia so) — ban Java goc, 1-MAY-CHU-GAME/src/boss/
	// Boss.java:507-510:
	//
	//	this.timeChatS = textChat.length() * 100;
	//	if (this.timeChatS > 2000) { this.timeChatS = 2000; }
	//
	// Tuc la: cang dai doc cang lau, nhung khong bao gio qua 2 giay.
	treoMoiKyTuThoai = 100 * time.Millisecond
	treoThoaiToiDa   = 2 * time.Second

	// thoiGianXacBoss la thoi gian TOI THIEU cai xac boss nam lai tren man
	// hinh truoc khi buoc RA_DI don no di.
	//
	// VI SAO PHAI CO CON SO NAY (ban Java KHONG co — no di ngay khi noi het
	// loi trang troi): mot pha khong khai bao ThoaiChet se bien mat NGAY trong
	// nhip boss chet. Nguoi choi vua ha duoc boss, nhin len thi no da khong
	// con — khong kip thay minh vua lam gi.
	//
	// VI SAO DUNG 3 GIAY: dong chu "X vua ha guc Boss Y" khong di thang, no
	// xep hang qua loa voi nhip toi thieu 1,5 giay (nhipToiThieuThongBao,
	// thongbao.go). Xac phai song lau hon MOT nhip loa, neu khong ca may chu
	// doc tin ve mot con boss ma nguoi trong khu vuc khong con thay dau vet.
	// 3 giay = hai nhip loa: mot nhip de tin toi noi, mot nhip lam le an toan.
	thoiGianXacBoss = 3 * time.Second
)

// treoThoai tra ve thoi gian mot cau thoai phai treo tren dau boss truoc khi
// cau ke tiep duoc noi (cong thuc + hai hang so: xem treoMoiKyTuThoai).
//
// DEM THEO KY TU (rune) chu khong theo BYTE — day la cho ban Go PHAI khac ban
// Java chu khong phai duoc phep khac: `String.length()` cua Java dem don vi
// UTF-16, ma moi chu tieng Viet co dau la DUNG MOT don vi UTF-16. Trong Go,
// len(chuoi) dem BYTE, va "Đủ rồi." la 7 ky tu nhung 10 byte. Dem byte thi moi
// cau tieng Viet treo lau hon cau khong dau cung do dai khoang 40% — mot su
// thien vi vo nghia sinh ra tu chi tiet ma hoa, khong tu y do nao ca.
func treoThoai(cau string) time.Duration {
	treo := time.Duration(utf8.RuneCountInString(cau)) * treoMoiKyTuThoai
	if treo > treoThoaiToiDa {
		return treoThoaiToiDa
	}
	return treo
}

// thoaiKeTiep tra loi cau hoi cua MOT tick: "boss co noi gi khong, va noi cau
// nao?". Ham THUAN (khong dong ho he thong, khong cham trang thai) — tach rieng
// de test go duoc tung moc thoi gian ma khong dung toi actor, dung khuon
// chonMucTieu (quai.go) da tach khoi aiQuai.
//
// Tra ve:
//   - (i >= 0, false) -> noi cau thu i o tick nay.
//   - (-1, false)     -> chua toi luot (cau vua noi con dang treo).
//   - (-1, true)      -> da noi HET va cau CUOI cung da treo du lau.
//
// THU TU HAI NHANH LA CO Y: kiem "con dang treo" TRUOC kiem "het cau". Neu dao
// lai, "xong" se dung ngay khi vua noi xong cau cuoi — cau chot cua boss vua
// hien ra da bi buoc trang thai ke tiep dam len. Ban Java giu dung thu tu nay
// (Boss.java:495-514: canDoWithTime bao ngoai, so sanh indexChatS ben trong).
//
// chiSo == 0 -> khong co cau nao dang treo -> noi NGAY, khong cho: mot con boss
// vua hien ra ma dung im mot nhip roi moi chao thi trong nhu bi lag.
func thoaiKeTiep(cau []string, chiSo int, mocCuoi, bayGio time.Time) (chiSoNoi int, xong bool) {
	if chiSo > 0 && chiSo <= len(cau) && bayGio.Sub(mocCuoi) < treoThoai(cau[chiSo-1]) {
		return -1, false
	}
	if chiSo < 0 || chiSo >= len(cau) {
		return -1, true
	}
	return chiSo, false
}

// doiTrangThaiBoss dua boss sang mot buoc moi cua vong doi VA dat lai moc thoi
// gian cua buoc do. Gop hai viec vao mot ham vi tach ra la mot cai bay: quen
// dat lai MocThoiGian thi buoc moi se do thoi gian tu buoc CU — vd cai xac se
// duoc coi la "da nam 3 giay" ngay giay dau tien vi boss da chao tu truoc do
// lau roi.
func doiTrangThaiBoss(b *Boss, tt TrangThaiBoss, bayGio time.Time) {
	b.TrangThai = tt
	b.MocThoiGian = bayGio
}

// noiThoaiBoss noi cau ke tiep trong mot mang thoai neu da toi luot, va tra ve
// true khi mang do da noi HET (buoc trang thai duoc phep di tiep).
//
// chiSo la CON TRO toi truong con tro chay cua boss (ChiSoThoaiHien hoac
// ChiSoThoaiChet) — moi mang thoai co con tro rieng cua no, xem chu thich hai
// truong do trong boss.go.
//
// Phat bang phatChoCaKhuVuc (KHONG phai phatChoNguoiKhac): boss khong co phien
// mang nen khong co ai de "tru ra", va moi nguoi trong khu vuc deu phai nghe.
func (k *KhuVuc) noiThoaiBoss(b *Boss, cau []string, chiSo *int, bayGio time.Time) bool {
	i, xong := thoaiKeTiep(cau, *chiSo, b.MocThoaiCuoi, bayGio)
	if xong {
		return true
	}
	if i < 0 {
		return false // cau truoc con dang treo
	}
	k.phatChoCaKhuVuc(goiChat44(b.Id, cau[i]))
	*chiSo = i + 1
	b.MocThoaiCuoi = bayGio
	return false
}

// nhipBoss chay may trang thai cho TUNG con boss trong khu vuc, roi GO khoi
// danh sach nhung con da roi han the gioi.
//
// VI SAO MOT VONG LAP RIENG chu khong nhet boss vao vong lap quai o tren:
//
//   - Hai loai nay khong chia se mot buoc nao — quai lo hoi sinh/hoi mau/AI,
//     boss lo vong doi pha. Nhet chung vao mot vong thi than vong lap tro thanh
//     mot chuoi if phan loai, khong phai mot y nghia.
//   - Chi phi khi khu vuc DONG: vong nay chay O(so boss), va so boss la 0 hoac
//     1 trong gan het thoi gian. `for range` tren mot lat cat rong khong ton gi
//     ca. Nguoc lai, neu nhet kiem tra boss VAO TRONG vong quai thi no se chay
//     N lan moi tick (N = so quai) chi de lam dung mot viec — dung kieu chi phi
//     am tham lon len theo do dong.
//
// bayGio nhan tu nhipTim (khong tu goi k.dongHo() lai) de MOI con boss trong
// cung mot nhip duoc xet tren CUNG mot moc thoi gian — cung ly do da ghi o
// nhipTim cho dan quai.
// nguoi la ANH CHUP nguoi choi ma nhipTim da dung MOT LAN cho ca nhip — truyen
// xuong thay vi cho tung con boss tu xay lai tu k.nguoiTrong: mot map->slice
// moi con boss moi tick la chi phi am tham lon len theo do dong, va (quan trong
// hon) hai con boss trong cung mot nhip se nhin thay hai anh chup khac nhau neu
// ai do vo tinh sua trang thai o giua.
func (k *KhuVuc) nhipBoss(nguoi []*nguoiTrongKhuVuc, bayGio time.Time) {
	if len(k.boss) == 0 {
		return
	}
	conLai := k.boss[:0] // dung lai chinh mang cu: chi don ve phia truoc
	for _, b := range k.boss {
		if k.mayTrangThaiBoss(b, nguoi, bayGio) {
			continue
		}
		conLai = append(conLai, b)
	}
	// XOA DUOI: cac o tu len(conLai) tro di van giu con tro toi nhung con boss
	// vua bi go. Khong xoa thi bo don rac khong thu hoi duoc chung — mot con
	// boss keo theo ca ban sao []PhaBoss cua no (ten, ba mang thoai), va moi
	// chu ky lich lai de lai them mot con nhu the.
	for i := len(conLai); i < len(k.boss); i++ {
		k.boss[i] = nil
	}
	k.boss = conLai
}

// mayTrangThaiBoss chay DUNG MOT BUOC vong doi cho mot con boss. Tra ve true
// khi con boss nay DA ROI HAN the gioi va phai duoc go khoi k.boss.
//
// MOI TICK NHIEU NHAT MOT BUOC (khong fallthrough sang buoc vua chuyen toi) —
// giong ban Java (Boss.java:315-349, mot `switch` moi lan update). Cai gia la
// moi buoc chuyen ton them 100ms; khong ai nhin thay 100ms, va doi lai than ham
// doc duoc tu tren xuong nhu chinh so do vong doi.
func (k *KhuVuc) mayTrangThaiBoss(b *Boss, nguoi []*nguoiTrongKhuVuc, bayGio time.Time) bool {
	if b == nil {
		return true // mot o nil lac vao danh sach: go luon, dung de moi tick sau phai kiem lai
	}

	pha := b.PhaHienTaiData()
	if pha == nil {
		// Khong co pha nao dang chay (du lieu hong, hoac chi so pha lac ra
		// ngoai). Con boss nay khong the ve, khong the noi, khong the danh.
		//
		// PHAI DON NO DI chu khong duoc de nam lai: the song cua muc lich da
		// de ra no dang BAT, nen bo dieu phoi se cho no "roi the gioi" mai mai
		// — mot dong du lieu go sai giet luon ca muc lich do, im lang.
		slog.Error("boss khong co pha dang chay — don khoi khu vuc",
			"idBoss", b.Id, "phaHienTai", b.PhaHienTai, "soPha", len(b.DuLieu),
			"trangThai", b.TrangThai.String())
		return k.bossRoiTheGioi(b, bayGio)
	}

	switch b.TrangThai {
	case BossNghi:
		// KHONG co mat trong the gioi. CHAN-7 cam actor dem nguoc GiayNghi de
		// tu danh thuc no, nen tu day KHONG CO DUONG NAO tro lai — giu lai chi
		// ton bo nho va ton mot luot lap moi 100ms. Go ra, KHONG phat goi gi
		// (client chua tung thay con nay).
		//
		// Con boss THAT SU tiep theo se do bo dieu phoi lich sinh, la mot doi
		// tuong MOI voi mot id moi.
		return true

	case BossHien:
		// Goi -5 da duoc phat truoc do roi — boi themBoss (duong lich) hoac
		// boi buoc bien hinh (Task 11). O day chi buoc tiep.
		doiTrangThaiBoss(b, BossChao, bayGio)

	case BossChao:
		if k.noiThoaiBoss(b, pha.ThoaiHien, &b.ChiSoThoaiHien, bayGio) {
			doiTrangThaiBoss(b, BossHoatDong, bayGio)
		}

	case BossHoatDong:
		// KH6b Task 10 — AI: chon muc tieu, danh bang chieu KH6, di lai gan khi
		// ngoai tam, thinh thoang noi mot cau. Chi tiet + ly do o aiBoss.
		//
		// KHONG doi trang thai o day: boss roi buoc HOAT_DONG CHI khi bi danh
		// chet (bossNhanSatThuong -> bossChet -> TROI). AI khong duoc tu quyet
		// ket thuc doi minh.
		k.aiBoss(b, pha, nguoi, bayGio)

	case BossTroi:
		// Cai xac. Hai dieu kien PHAI cung du, va thu tu kiem la co y:
		//
		//  1. Noi het loi trang troi — kiem TRUOC, vi day la loi CHET: no phai
		//     ra ngay luc chet chu khong phai sau khi cho het 3 giay.
		//  2. Xac da nam du thoiGianXacBoss — kiem SAU, nen thoi gian nam xac
		//     thuc te la max(loi trang troi, 3 giay).
		if !k.noiThoaiBoss(b, pha.ThoaiChet, &b.ChiSoThoaiChet, bayGio) {
			return false
		}
		if bayGio.Sub(b.MocThoiGian) < thoiGianXacBoss {
			return false
		}
		doiTrangThaiBoss(b, BossRaDi, bayGio)

	case BossRaDi:
		// ⚠️ CHO NOI CUA TASK 11 (bien hinh qua pha): theo thiet ke §5, con pha
		// thi o day boss phai quay ve BossHien voi pha ke tiep, GIU NGUYEN Id
		// (dat lai Hp/HpToiDa/NhipDanhHienTai/MocNoiGianDaQua va CA HAI con tro
		// thoai ve 0, ghi LucHienPha = bayGio, phat lai -5 voi ngoai hinh moi).
		//
		// TASK 9 CO Y KHONG LAM NHANH DO, va day la ly do — khong phai vi ngai
		// viec: bay #3 cua thiet ke §7 doi hoi phai KIEM CLIENT truoc, rang mot
		// id da nhan -6 co duoc quay lai bang -5 hay khong. Cuoc kiem do la
		// task 1 (GATE) va no CHUA LAM. Code truoc roi kiem sau, neu client tu
		// choi, thi boss bien hinh xong se BIEN MAT — dung loai loi chi lo ra
		// truoc mat nguoi choi that.
		//
		// HOM NAY, HAU QUA THAT: mot con boss NHIEU PHA chet o pha 1 se roi han
		// the gioi thay vi lot xac. Do la lua chon hong VE PHIA AN TOAN — thieu
		// mot pha thi nguoi choi thieu mot doan noi dung; con mot con boss ket
		// cung o RA_DI thi khong bao gio goi Tat(), va muc lich de ra no chet
		// vinh vien.
		return k.bossRoiTheGioi(b, bayGio)
	}
	return false
}

// ─────────────────────────────────────────────────────────────────────────────
// KH6b Task 10 — AI BOSS (than cua buoc HOAT_DONG)
//
// Moi thu duoi day chay TRONG nhipTim 100ms cua actor: khong cham mang, khong
// cham CSDL, khong cho ai (ky luat §6.9/§6.10b). Moi goi phat ra di qua
// phatChoCaKhuVuc/phatChoNguoiKhac nhu moi nhanh chien dau khac.
// ─────────────────────────────────────────────────────────────────────────────

const (
	// heSoTamDuoiBoss: tam DUOI (aggro) = tam chieu XA NHAT cua pha × he so
	// nay. Lay theo du lieu chu khong ghim mot con so rieng — mot con boss
	// toan chieu can chien (tam 100) khong co ly do gi de duoi nguoi tu dau
	// ben kia ban do, con boss co chieu tam xa (300) thi phai duoi xa hon.
	//
	// 2 lan: du de boss buoc ra khoi cho dung ma khong bien thanh con cho san
	// bam theo nguoi choi qua nua ban do. Day la con so CAN BANG (doan ban
	// dau) — khong co bang chung nao trong ban Java, va no se duoc chinh khi
	// co noi dung that. Ghi ro o day de lan sau ai doi thi biet minh dang doi
	// mot lua chon, khong phai mot hang so co nguon.
	heSoTamDuoiBoss int64 = 2

	// buocBossMoiNhip: boss di toi da bao nhieu don vi trong MOT nhip 100ms.
	//
	// 24 = DUNG mot o cua luoi va cham map 1 (TaiCollision bat buoc tileWidth
	// == tileHeight == 24, catalog.go). Con so nay khong phai chon cho dep:
	// buoc dai hon mot o thi mot buoc co the NHAY QUA mot buc tuong mong
	// (DiDuocDoan lay mau doc duong di nen van bat duoc, nhung se co nhung
	// khe hep ma boss "lot" qua tuy theo goc); buoc bang mot o thi moi lan
	// kiem va cham deu tren dung do phan giai ma ban do duoc ve.
	//
	// 24 don vi/100ms = 240 don vi/giay, xap xi toc do chay cua nguoi choi
	// (NguongBuoc 500 moi goi di chuyen, client gui vai goi mot giay).
	buocBossMoiNhip int16 = 24

	// nhipThoaiDanhBoss: "thinh thoang" cua ThoaiDanh (thiet ke §5) la MOT CON
	// SO — 15 giay giua hai cau.
	//
	// Vi sao khong nhanh hon: goi 44 phat cho CA khu vuc, va mot con boss noi
	// lien mom bien khung chat thanh rac, nguoi choi tat luon khung chat roi
	// khong doc duoc gi ca. Vi sao khong cham hon: mot tran boss thuong ket
	// thuc trong khoang mot phut, cham hon nua thi phan lon cau thoai khong
	// bao gio duoc ai nghe.
	//
	// Day la con so CAN BANG do em chon (ban Java khong co nhip nay — no chat
	// theo su kien). Khac han hai hang so treoThoai o tren: chung co bang
	// chung tu Boss.java:507-510, con con so nay thi khong.
	nhipThoaiDanhBoss = 15 * time.Second
)

// aiBoss chay MOT nhip AI cho mot con boss dang o buoc HOAT_DONG.
//
// THU TU BON VIEC LA CO Y:
//
//  1. Chon muc tieu. KHONG co ai -> return NGAY, khong cham mot truong nao
//     (LanDanhCuoi, MocThoaiCuoi, vi tri). Day dung la cai bay ma aiQuai da
//     ghi lai tu ban Java (Mob.java:386-390 dat lai moc thoi gian NGOAI nhanh
//     "co nguoi choi", khien con quai vua thay nguoi phai cho them gan mot
//     nhip nua moi duoc danh).
//  2. Noi thoai. Truoc khi danh, de cau noi va don danh cua cung mot nhip di
//     cung nhau; va CHI khi da co muc tieu — boss lam bam mot minh giua ban do
//     trong la loi ai cung nhin thay.
//  3. Thu danh. Neu co mot chieu voi toi thi danh va DUNG o day.
//  4. Khong chieu nao voi toi -> di lai gan.
//
// Boss dang trong nhip hoi (buoc 3 tu choi vi chua toi luot) thi KHONG di
// chuyen: no dung yen thu the, khong bo di roi lai lan vao — dung nhu ban Java
// (Boss.java chi di chuyen khi khong danh duoc vi KHOANG CACH).
func (k *KhuVuc) aiBoss(b *Boss, pha *PhaBoss, nguoi []*nguoiTrongKhuVuc, bayGio time.Time) {
	tamXaNhat := tamChieuXaNhatCuaPha(pha)
	if tamXaNhat <= 0 {
		// Pha khong khai bao BoChieu (hoac du lieu chieu hong): FAIL-CLOSED.
		// Khong co chieu thi khong co tam, khong co tam thi khong co muc tieu
		// — boss dung im thay vi bien thanh mot con danh tay khong tam vo han.
		return
	}

	mucTieu := chonMucTieuBoss(b, nguoi, tamXaNhat*heSoTamDuoiBoss, bayGio)
	if mucTieu == nil {
		return
	}

	k.thoaiDanhBoss(b, pha, bayGio)

	if k.bossThuDanh(b, pha, mucTieu, bayGio) {
		return
	}
	k.bossDiToi(b, int16(mucTieu.nv.ViTriX), int16(mucTieu.nv.ViTriY))
}

// tamChieuXaNhatCuaPha tra ve tam danh LON NHAT trong bo chieu cua pha (0 khi
// pha khong co chieu nao tra cuu duoc). Dung lam goc tinh tam duoi.
//
// Tra cuu qua TraCapChieu (chieuthuc.go) — cua DUY NHAT doi {templateId, cap}
// thanh *CapChieu, dung chung voi nguoi choi. Tu doc catalog o day se tao ban
// sao logic thu hai, va khac biet giua hai ban se hong IM LANG.
func tamChieuXaNhatCuaPha(pha *PhaBoss) int64 {
	if pha == nil {
		return 0
	}
	var xaNhat int64
	for _, bo := range pha.BoChieu {
		cap := TraCapChieu(byte(pha.Lang), byte(bo[0]), byte(bo[1]))
		if cap == nil {
			continue
		}
		if int64(cap.Dx) > xaNhat {
			xaNhat = int64(cap.Dx)
		}
	}
	return xaNhat
}

// chonMucTieuBoss chon nguoi boss se danh trong nhip nay, theo HAI TANG UU TIEN
// — dung khuon AI quai (Quai.chonMucTieu, quai.go) vi cung mot cau hoi:
//
//	Tang 1: nguoi DA DANH boss nay, ai gay NHIEU sat thuong nhat thi bi nham
//	        truoc. Boss quay lai voi ke danh no dau nhat — de nguoi tank giu
//	        duoc boss, va de mot nguoi di ngang qua khong bi cuop mat su chu y
//	        giua tran.
//	Tang 2: khong ai tung danh (boss vua hien) -> ai GAN NHAT.
//
// Ba dieu kien loai, deu bat buoc va deu co dot bien rieng trong test:
//   - da chet (HpHienTai <= 0): mot cai xac khong phai muc tieu.
//   - dang bat tu sau hoi sinh (< 1500ms, dac ta §1.7): boss phai ton trong y
//     het nguoi choi va quai, neu khong no don chet nguoi ta ngay khi ho vua
//     bam nut hoi sinh.
//   - ngoai tam duoi: boss khong duoc duoi nguoi qua nua ban do.
//
// "Da roi khu vuc" KHONG can dieu kien rieng: `nguoi` la anh chup cua
// k.nguoiTrong ngay trong nhip nay, ai da ra thi khong con trong do.
//
// HOA thi lay ID NHO HON. Khong phai vi id nho quan trong hon — ma vi `nguoi`
// duoc xay tu mot map, va thu tu duyet map trong Go la NGAU NHIEN moi lan chay.
// Khong co chot pha hoa nay thi hai nguoi dung cach deu nhau se lam boss quay
// qua quay lai giua ho moi tick, va khong test nao chay on dinh duoc.
func chonMucTieuBoss(b *Boss, nguoi []*nguoiTrongKhuVuc, tamDuoi int64, bayGio time.Time) *nguoiTrongKhuVuc {
	tamSq := tamDuoi * tamDuoi

	var (
		gopNhat    *nguoiTrongKhuVuc // tang 1
		gopCaoNhat int64
		ganNhat    *nguoiTrongKhuVuc // tang 2
		ganNhatSq  int64
	)

	for _, n := range nguoi {
		if n == nil || n.nv == nil {
			continue
		}
		if n.nv.HpHienTai <= 0 {
			continue
		}
		if bayGio.Sub(n.lucHoiSinh) < thoiGianBatTuSauHoiSinh {
			continue
		}
		dx := int64(n.nv.ViTriX) - int64(b.ViTriX)
		dy := int64(n.nv.ViTriY) - int64(b.ViTriY)
		distSq := dx*dx + dy*dy
		if distSq > tamSq {
			continue
		}

		if gop := b.SatThuongTheoNguoi[n.nv.ID]; gop > 0 {
			if gopNhat == nil || gop > gopCaoNhat ||
				(gop == gopCaoNhat && n.nv.ID < gopNhat.nv.ID) {
				gopNhat, gopCaoNhat = n, gop
			}
		}
		if ganNhat == nil || distSq < ganNhatSq ||
			(distSq == ganNhatSq && n.nv.ID < ganNhat.nv.ID) {
			ganNhat, ganNhatSq = n, distSq
		}
	}

	if gopNhat != nil {
		return gopNhat
	}
	return ganNhat
}

// nhipDanhBoss tra khoang cach TOI THIEU giua hai don cua boss.
//
// 0 (pha khong khai NhipDanhMs) nghia la "de hoi chieu cua chinh chieu do lo" —
// KiemDanhHopLe van kiem cap.HoiChieuMs nen boss khong bao gio danh moi tick.
//
// TODO(Task 11 — noi gian theo %mau): Task 11 nhan MocNoiGian.HeSoNhip vao
// b.NhipDanhHienTai NGAY LUC cham moc (mot lan), khong phai o day moi tick —
// nhan moi tick thi he so se cong don vo tan. Ham nay chi DOC con so da co.
func nhipDanhBoss(b *Boss) time.Duration {
	if b.NhipDanhHienTai <= 0 {
		return 0
	}
	return time.Duration(b.NhipDanhHienTai) * time.Millisecond
}

// bossThuDanh thu tung chieu trong BoChieu cua pha va danh bang chieu DAU TIEN
// duoc phep — thu tu khai bao trong du lieu CHINH LA thu tu uu tien, do nguoi
// lam noi dung quyet dinh (khong co luat an nao o code).
//
// Tra true khi da THAT SU danh (co phat goi). false = khong chieu nao voi toi,
// nguoi goi se cho boss di lai gan.
//
// KiemDanhHopLe la duong kiem DUY NHAT — cung mot bo luat voi nguoi choi
// (chotkiem.go, Task 5 don nen). BenDanh.VoHanChakra = true vi boss khong co
// bang chakra nao ca; no CHI mien kiem chakra, KHONG mien bat tu/hoi chieu/muc
// tieu chet/tam.
//
// ketQua.Truot (vuot tam) o day KHONG phat goi "miss" nhu huong nguoi->nguoi:
// boss chua he quyet dinh danh, no dang HOI xem chieu nay co toi khong. Thu
// chieu ke tiep, roi di lai gan — dung nhu mot con boss that lam.
func (k *KhuVuc) bossThuDanh(b *Boss, pha *PhaBoss, mucTieu *nguoiTrongKhuVuc, bayGio time.Time) bool {
	if nhip := nhipDanhBoss(b); nhip > 0 && bayGio.Sub(b.LanDanhCuoi) < nhip {
		return false // con trong nhip cua pha -> dung yen thu the, khong di
	}

	for _, bo := range pha.BoChieu {
		cap := TraCapChieu(byte(pha.Lang), byte(bo[0]), byte(bo[1]))
		if cap == nil {
			continue
		}
		ketQua := KiemDanhHopLe(BenDanh{
			X:           int32(b.ViTriX),
			Y:           int32(b.ViTriY),
			LanDanhCuoi: b.LanDanhCuoi,
			Chieu:       cap,
			VoHanChakra: true, // boss khong co bang chakra — xem chu thich truong do
		}, ViTriCoMau{
			X:          int16(mucTieu.nv.ViTriX),
			Y:          int16(mucTieu.nv.ViTriY),
			HpHienTai:  mucTieu.nv.HpHienTai,
			LucHoiSinh: mucTieu.lucHoiSinh,
		}, bayGio)
		if !ketQua.ChoPhep || ketQua.Truot {
			continue
		}
		k.bossDanhNguoi(b, pha, cap, mucTieu, bayGio)
		return true
	}
	return false
}

// bossDanhNguoi ap MOT don cua boss len mot nguoi choi va phat dung bo goi ma
// huong NGUOI-DANH-NGUOI (case suKienDanhNguoiChoi) da dung — KHONG PHAI -9/-12
// cua quai.
//
// VI SAO -60 CHU KHONG PHAI -11/-10 (goi quai danh nguoi): hai goi do mang
// `byte mobId`, mot id boss AM 4 byte khong nhet vao duoc. Voi client, boss la
// mot NHAN VAT (no vao the gioi bang chinh goi -5 cua nguoi choi), nen moi thu
// nham vao no hay phat ra tu no deu di duong nhan vat.
//
// BA GOI, dung thu tu va dung nguoi nhan:
//   - -60 CA KHU VUC: hoat anh + con so sat thuong. Client cua nan nhan tu tru
//     con so nay khoi thanh mau cua chinh minh, nen no KHONG can goi rieng.
//   - -30/9 NGUOI KHAC (khong phai nan nhan): goi -60 khong mang truong "hp con
//     lai" nen nhung nguoi khac can goi nay de ve lai dung thanh mau nan nhan.
//     CHI gui khi mau THUC SU doi.
//   - -17/-8 khi nan nhan chet: di DUNG duong xuLyNguoiChoiChet da co (Ke hoach
//     6 Task 13) — khong che duong bao chet thu hai.
func (k *KhuVuc) bossDanhNguoi(b *Boss, pha *PhaBoss, cap *CapChieu, mucTieu *nguoiTrongKhuVuc, bayGio time.Time) {
	// CHI o day LanDanhCuoi moi duoc dat lai — dung nhanh "that su danh", cung
	// ky luat voi aiQuai va voi hai nhanh danh cua nguoi choi.
	b.LanDanhCuoi = bayGio

	// chiMang: KH6 CHUA chot cong thuc/ti le chi mang (khong dac ta nao) ->
	// luon false, dung lua chon an toan ma ca ba nhanh danh hien co dang dung.
	chiMang := false

	// He nguyen to THAT cua ca hai phia: boss co lang that (PhaBoss.Lang, dung
	// cho ca goi -5) nen KHONG dung heQuaiChuaXacDinh nhu quai.
	//
	// TODO(Task 11 — noi gian theo %mau): MocNoiGian.HeSoDame nhan vao DUNG CHO
	// NAY (sat thuong boss GAY RA). HeSoDef thi nam o huong nguoc lai, trong
	// nhanh nguoi-danh-boss (case suKienDanhNguoiChoi) — hai he so hai cho, dung
	// nham cho la lam boss vua trau vua manh gap doi.
	satThuong := TinhSatThuong(pha.Dame, cap,
		heNguyenToNguoiChoi(pha.Lang), heNguyenToNguoiChoi(mucTieu.nv.Lang), chiMang)

	// §1.2: sat thuong khong duoc VUOT mau con lai. Khong co nhanh §1.3 ("day
	// mau khong chet boi mot don") — dac ta ghi RO quy tac do CHI danh cho
	// QUAI, khong cho nguoi choi.
	if satThuong > mucTieu.nv.HpHienTai {
		satThuong = mucTieu.nv.HpHienTai
	}
	if satThuong < 0 { // phong thu: du lieu Dame am khong duoc HOI MAU nguoi choi
		satThuong = 0
	}
	mucTieu.nv.HpHienTai -= satThuong
	daChet := mucTieu.nv.HpHienTai <= 0

	// int32(b.Id): dung phep thu hep ma goiXuatHienDayDu da dung khi ghi id
	// boss vao goi -5, nen con so o day KHOP CHINH XAC id client dang giu.
	k.phatChoCaKhuVuc(goiDanhNguoiChoi(
		int32(b.Id), byte(cap.SkillId), int32(mucTieu.nv.ID), satThuong, daChet, chiMang,
	))

	if satThuong > 0 {
		k.phatChoNguoiKhac(mucTieu.nv.ID, goiCapNhatMau(
			int32(mucTieu.nv.ID), mucTieu.nv.HpHienTai, mucTieu.nv.HpToiDa,
		))
	}

	if daChet {
		k.xuLyNguoiChoiChet(mucTieu)
	}
}

// bossDiToi buoc boss MOT buoc (toi da buocBossMoiNhip) ve phia (dichX, dichY)
// va bao ca khu vuc bang goi -7. Tra true khi that su di duoc.
//
// VA CHAM: di qua kiemDiChuyenHopLe (dichuyen.go) — CHINH ham ma goi di chuyen
// cua nguoi choi phai qua. Boss khong duoc co mot luat rieng: neu no di xuyen
// nui ma nguoi choi thi khong, no se dung giua da roi ban ra, khong ai danh lai
// duoc. Ham do cung lo luon bien ban do cho cac map chua co luoi va cham.
//
// Vuong tuong -> DUNG YEN (khong truot doc tuong, khong tim duong vong): tim
// duong that su la mot bai toan rieng, va lam nua voi thi sinh ra kieu boss
// rung lac tai cho. Boss dung im la de nhin ra "cho nay chan duong", con boss
// rung lac thi trong nhu game hong.
//
// Phat cho CA KHU VUC (khong phai phatChoNguoiKhac): boss khong co phien mang
// nen khong co ai de "tru ra" — moi nguoi trong khu vuc deu phai thay no di.
func (k *KhuVuc) bossDiToi(b *Boss, dichX, dichY int16) bool {
	dx := float64(dichX) - float64(b.ViTriX)
	dy := float64(dichY) - float64(b.ViTriY)
	khoang := math.Hypot(dx, dy)
	if khoang < 1 {
		return false // da dung ngay tren muc tieu
	}

	buoc := float64(buocBossMoiNhip)
	if buoc > khoang {
		buoc = khoang // khong buoc qua muc tieu roi lai quay dau
	}
	// math.Round CHI dung cho HUONG DI (toa do), TUYET DOI khong dung o duong
	// tinh sat thuong — sat thuong van la so nguyen tu dau den cuoi
	// (TinhSatThuong, satthuong.go), day la hai the gioi khac nhau.
	toiX := b.ViTriX + int16(math.Round(dx/khoang*buoc))
	toiY := b.ViTriY + int16(math.Round(dy/khoang*buoc))
	if toiX == b.ViTriX && toiY == b.ViTriY {
		return false
	}
	if !kiemDiChuyenHopLe(k.banDo, b.ViTriX, b.ViTriY, toiX, toiY) {
		return false
	}

	b.ViTriX, b.ViTriY = toiX, toiY
	k.phatChoCaKhuVuc(goiDiChuyen(b.Id, toiX, toiY))
	return true
}

// thoaiDanhBoss noi mot cau ThoaiDanh khi da du nhipThoaiDanhBoss ke tu cau
// truoc. GOI KHI da co muc tieu (xem aiBoss) — boss lam bam mot minh giua ban
// do trong la loi ai cung nhin thay.
//
// QUAY VONG tren mang thoai (chi so % len) — KHAC han noiThoaiBoss cua buoc
// CHAO/TROI, cai do doc HET mang mot lan roi thoi. Ly do khac nhau: loi chao va
// loi trang troi la mot MAN, co dau co cuoi; con thoai danh la ban nhac nen
// chay suot tran — mot tran dai hon so cau khai bao la chuyen binh thuong, va
// im bat sau cau cuoi thi trong nhu boss het pin.
//
// Dung chung b.MocThoaiCuoi voi hai buoc kia LA CO Y: mot con boss chi co mot
// cai mieng. Cau chao cuoi cung vua ra xong thi cau thoai danh dau tien phai
// doi, khong dam len nhau.
func (k *KhuVuc) thoaiDanhBoss(b *Boss, pha *PhaBoss, bayGio time.Time) {
	if len(pha.ThoaiDanh) == 0 {
		return
	}
	// MocThoaiCuoi zero (pha khong co ThoaiHien nen buoc CHAO khong noi gi) ->
	// bayGio.Sub(zero) rat lon -> cau dau ra NGAY. Dung y do: boss vua vao tran
	// thi noi mot cau, khong phai dung im 15 giay roi moi len tieng.
	if bayGio.Sub(b.MocThoaiCuoi) < nhipThoaiDanhBoss {
		return
	}
	cau := pha.ThoaiDanh[b.ChiSoThoaiDanh%len(pha.ThoaiDanh)]
	k.phatChoCaKhuVuc(goiChat44(b.Id, cau))
	b.ChiSoThoaiDanh = (b.ChiSoThoaiDanh + 1) % len(pha.ThoaiDanh)
	b.MocThoaiCuoi = bayGio
}

// bossRoiTheGioi don MOT con boss ra khoi the gioi: bao client xoa hinh, dua
// trang thai ve NGHI, roi mo lai cua so sinh cua muc lich. Luon tra true —
// nguoi goi dung gia tri do de go no khoi k.boss.
//
// THU TU BA VIEC LA CO Y:
//
//  1. Phat -6 TRUOC. Day la thu duy nhat nguoi choi nhin thay; lam truoc thi du
//     hai buoc sau co hong the nao, man hinh nguoi choi van sach.
//  2. Dat trang thai NGHI. Thiet ke §5 goi day la buoc cuoi vong doi. No con
//     la lop chan cuoi cho moi duong goi tuong lai: DangHien() = false nen mot
//     con tro con sot lai o dau do khong the lam ai gui -5 cua no nua.
//  3. Tat the song SAU CUNG. Day chinh la giay phut bo dieu phoi duoc phep
//     sinh con moi — tat no som hon (vd truoc khi phat -6) thi trong mot cua so
//     rat hep, mot con boss thu hai co the duoc sinh trong khi con thu nhat
//     chua kip xoa khoi man hinh.
//
// Tat() la mot phep ghi NGUYEN TU (atomic.Bool) nen goi tu trong actor la an
// toan tuyet doi: khong cap phat, khong I/O, khong bao gio chan — dung yeu cau
// cua bay §6.10b. nil-safe: boss khong thuoc muc lich nao (test, hoac sau nay
// la duong sinh theo trigger) thi day la mot lenh khong lam gi.
func (k *KhuVuc) bossRoiTheGioi(b *Boss, bayGio time.Time) bool {
	k.phatChoCaKhuVuc(goiRoiKhuVuc(b.Id))
	doiTrangThaiBoss(b, BossNghi, bayGio)
	b.The.Tat()
	return true
}

// Gui day mot su kien vao hop thu cua actor. KHONG CHAN (non-blocking): neu
// hop thu day thi bo goi — DEM + LOG muc Warn (Ke hoach 11 Task 2, xem chu
// thich truong soLanGuiRoi o dinh nghia KhuVuc: TRUOC Task 2 nhanh nay chi
// ghi slog.Debug, muc thuong TAT trong production nen la roi IM LANG THAT
// SU; mot su kien roi co the la mot yeu cau chup anh de luu, khien ca khu
// vuc bi bo qua mot luot luu ma khong ai biet). fmt.Sprintf("%T", sk) ghi
// kieu su kien bi roi — phan biet "roi mot goi di chuyen" (vo hai, se co goi
// moi hon ngay sau) voi "roi mot yeu cau chup anh" (dang lo ngai hon nhieu)
// khi doc log. Neu actor da Dung() thi cung bo goi an toan (khong panic, vi
// hopThu khong bao gio bi close — chi kenh dung bi close).
func (k *KhuVuc) Gui(sk suKienKhuVuc) {
	if k.xong != nil {
		select {
		case <-k.xong:
			return
		default:
		}
	}
	select {
	case k.hopThu <- sk:
	case <-k.xong:
		// actor da dung, bo su kien.
	default:
		tong := k.soLanGuiRoi.Add(1)
		slog.Warn("hop thu actor day, bo su kien — mot luot xu ly (vd chup anh de luu) co the bi bo qua",
			"mapId", k.banDo.MapID, "kieuSuKien", fmt.Sprintf("%T", sk), "tongSoLanRoi", tong)
	}
}

// Dung dong actor mot cach sach va tra loi ro neu control khong the den actor
// trong han. KhiDong va don dep luc tat may chu khong duoc ket vinh vien chi
// vi actor da ket logic hay mailbox control bi day.
func (k *KhuVuc) Dung() error {
	k.dungMotLan.Do(func() {
		if err := k.guiDieuKhien(dieuKhienDung{}); err != nil {
			k.loiDung = err
			return
		}
		select {
		case <-k.xong:
		case <-time.After(thoiGianChoControl):
			k.loiDung = ErrControlQuaHan
		}
	})
	return k.loiDung
}
