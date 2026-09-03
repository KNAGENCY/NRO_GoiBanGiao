// noidungboss.go — CUA DUY NHAT cua goi `game` di vao KHO NOI DUNG cho phan BOSS.
//
// (KH8 task 7. Thiet ke: docs/superpowers/specs/2026-09-03-kh8-kho-noi-dung-design.md
// hang "7" bang Task, §4 "cai gi khong nap nong duoc", dot bien M7/M8/M10/M11.)
//
// # CAI GI DOI O TASK NAY
//
// Truoc task nay, hai con boss mau la mot BANG GO CUNG trong chinh tep nay
// (bien `danhSachNoiDungBoss`, ~200 dong du lieu + chu thich). Bang do DA XOA.
// Tu nay nguon duy nhat la kho-noi-dung/boss.yaml — MOI chu thich giai thich
// "vi sao lai la con so nay" da di theo sang file YAML, khong bi vut bo (do
// chinh la ly do KH8 chon YAML chu khong chon JSON, xem noidung/noidung.go).
//
// Kieu NoiDungBoss/ViTriSinh va hai ham TimNoiDungBoss/DanhSachNoiDungBoss GIU
// NGUYEN chu ky cu — lichboss.go va toan bo noidungboss_test.go goi chung khong
// doi mot dong nao. Chi noi DU LIEU o dau ra la doi.
//
// # KHUON GIONG khoquai.go (task 8) va lang.go (task 6), CO Y
//
// Cung mot sync.Once cho ban NHUNG (khoNhungMotLan, khochieu.go) — MOT anh chup
// noi dung la MOT khoi thong nhat (chieu + lang + quai + boss cung nam trong
// mot *noidung.DuLieu), hai Once cho cung mot file la hai co hoi nap lech nhau.
//
// # OUTFIT: KIEM O `game`, KHONG KIEM O `noidung`
//
// Ba o dau/than/chan cua `outfit` phai tro toi id part co that trong goi asset
// DA KY (assets/asset-lock.json, closure.part_ids). Danh sach do doc boi goi
// `game` (catalog.go); `noidung` KHONG duoc import `game`. Vi vay bo kiem chia
// doi y het duong anh_id cua quai:
//
//	noidung/boss.go kiemBoss          — so, thu tu, ma, moc noi gian, nhip...
//	noidung/boss.go KiemIdPartBoss    — ham THUAN, nhan danh sach id lam THAM SO
//	game/noidungboss.go (tep nay)     — nguoi goi ham tren, voi danh sach THAT
//
// # BOSS DANG SONG KHONG NAP NONG DUOC — DA CHOT (thiet ke §4)
//
// Bo dieu phoi lich (lichboss.go) doc kho noi dung LUC SINH: no lay so cua the
// he dang co hieu luc TAI KHOANH KHAC do, sao sau vao con boss, roi ban giao
// cho actor. Ke tu do con boss giu so cua THE HE LUC SINH cho toi khi no roi
// khoi the gioi. Nap nong boss.yaml giua tran KHONG dung toi no.
//
// Lua chon CO CHU DICH, cung ly do voi quai (khoquai.go):
//
//   - HpToiDa cua boss da di xuong client trong goi -5 luc no hien ra. Doi giua
//     chung ma khong gui lai la lam thanh mau o client sai vinh vien.
//   - Doi Dame giua tran la doi luat choi cua nhung nguoi dang danh do — ho
//     khong thay mot dong thong bao nao.
//
// Con boss SINH RA SAU lan nap nong moi lay so moi. Test khoa dieu nay:
// TestLichBoss_BossDangSongKhongDoiSoKhiNapNong (lichboss_napnong_test.go).
package game

import (
	"sync"

	"nroserver/internal/noidung"
)

// ---------------------------------------------------------------------------
// 1. Hai kieu bao ngoai — GIU NGUYEN chu ky cu cua goi `game`
// ---------------------------------------------------------------------------

// ViTriSinh la MOT cho dat chan hop le cua boss: map nao, toa do nao.
//
// VI SAO KHONG NAM TRONG PhaBoss: PhaBoss.MapVao chi noi boss hien o MAP nao,
// khong noi hien o CHO nao trong map. Ma cho dat chan lai la thu duy nhat trong
// ca bo du lieu nay co the kiem bang mot nguon KHACH QUAN (luoi va cham cua
// map), nen no dang duoc mot kieu rieng de bo test bat duoc.
//
// Toa do la "chan nhan vat dung o dau", cung quy uoc voi DiemDi.DenX/DenY
// (banDo.go) va voi luoi va cham: o chua toa do phai DI DUOC, va o ngay duoi
// phai DAC (co mat dat). Sinh giua troi thi boss roi/lo lung tuy client.
type ViTriSinh struct {
	Map  int32
	X, Y int16
}

// NoiDungBoss la MOT con boss trong kho noi dung: cac pha cua no + cac cho no
// duoc phep dat chan. Ma la khoa tra cuu on dinh (bo dieu phoi lich goi boss
// theo ma, khong theo chi so mang — chen mot con vao giua danh sach khong duoc
// lam doi y nghia cua lich).
type NoiDungBoss struct {
	Ma    string
	Pha   []PhaBoss
	ViTri []ViTriSinh
}

// MaBossNhanGiaMatNa / MaBossHocVienDaoTau: khoa tra cuu cua hai con boss mau.
// Hang so (khong phai chuoi tran) de bo dieu phoi lich va bo test cung tro toi
// MOT cho — go nham mot chu trong chuoi tran se im lang tra nil.
//
// Hai hang nay la CAU NOI voi boss.yaml: doi `ma` trong YAML ma quen sua o day
// (hoac nguoc lai) se lam TimNoiDungBoss tra nil — bo test bat duoc ngay
// (TestNoiDungBoss_MaDuyNhatVaTraCuuDung).
const (
	MaBossNhanGiaMatNa  = "nhan-gia-mat-na"
	MaBossHocVienDaoTau = "hoc-vien-dao-tau"
)

// ---------------------------------------------------------------------------
// 2. Doi kieu: noidung.Boss -> NoiDungBoss (kieu cua goi `game`)
// ---------------------------------------------------------------------------

// khoBossAnhChup la phan BOSS cua mot anh chup noi dung — cung khuon voi
// khoLangAnhChup (lang.go): mot interface nho vi *noidung.anhChup la kieu KHONG
// XUAT KHAU. *noidung.anhChup thoa man no (hai cua nay cung nam trong
// khoChieuAnhChup de actor va lich dung chung MOT anh chup).
type khoBossAnhChup interface {
	Boss(ma string) (noidung.Boss, bool)
	DanhSachBoss() []noidung.Boss
}

// phaBossTuNoiDung doi noidung.PhaBoss (kieu YAML) sang PhaBoss (kieu goi
// `game`, khuon byte -5). Cac lat cat deu duoc DUNG MOI o day — anh chup da tra
// ban sao sau, va ham nay khong bao gio tra ve mot lat cat dung chung voi no.
//
// Ep kieu byte cho Cap/TypePk la AN TOAN o day vi bo kiem cua noidung da chan
// mien 1..255 tu luc nap (kiemMotPhaBoss); mot gia tri vuot mien khong bao gio
// di den duoc day.
func phaBossTuNoiDung(p noidung.PhaBoss) PhaBoss {
	ra := PhaBoss{
		Ten:        p.Ten,
		Lang:       int16(p.Lang),
		Cap:        byte(p.Cap),
		Outfit:     p.Outfit,
		Dame:       p.Dame,
		TypePk:     byte(p.TypePk),
		NhipDanhMs: p.NhipDanhMs,
		GiayNghi:   p.GiayNghi,
		LaBossLon:  p.LaBossLon,
		DiemThuong: p.DiemThuong,
	}
	if p.HpChon != nil {
		ra.HpChon = append([]int64(nil), p.HpChon...)
	}
	if p.MapVao != nil {
		ra.MapVao = append([]int32(nil), p.MapVao...)
	}
	if p.BoChieu != nil {
		ra.BoChieu = make([][2]int32, len(p.BoChieu))
		for i, o := range p.BoChieu {
			ra.BoChieu[i] = [2]int32{o.TemplateId, o.Cap}
		}
	}
	if p.MocNoiGian != nil {
		ra.MocNoiGian = make([]MocNoiGian, len(p.MocNoiGian))
		for i, m := range p.MocNoiGian {
			ra.MocNoiGian[i] = MocNoiGian{
				PhanTramHP: m.PhanTramHP,
				HeSoDame:   m.HeSoDame,
				HeSoDef:    m.HeSoDef,
				HeSoNhip:   m.HeSoNhip,
				Chat:       m.Chat,
			}
		}
	}
	if p.ThoaiHien != nil {
		ra.ThoaiHien = append([]string(nil), p.ThoaiHien...)
	}
	if p.ThoaiDanh != nil {
		ra.ThoaiDanh = append([]string(nil), p.ThoaiDanh...)
	}
	if p.ThoaiChet != nil {
		ra.ThoaiChet = append([]string(nil), p.ThoaiChet...)
	}
	return ra
}

func noiDungBossTuNoiDung(b noidung.Boss) NoiDungBoss {
	ra := NoiDungBoss{Ma: b.Ma}
	if b.ViTri != nil {
		ra.ViTri = make([]ViTriSinh, len(b.ViTri))
		for i, v := range b.ViTri {
			ra.ViTri[i] = ViTriSinh{Map: v.Map, X: v.X, Y: v.Y}
		}
	}
	if b.Pha != nil {
		ra.Pha = make([]PhaBoss, len(b.Pha))
		for i, p := range b.Pha {
			ra.Pha[i] = phaBossTuNoiDung(p)
		}
	}
	return ra
}

// ---------------------------------------------------------------------------
// 3. Kiem outfit doi chieu catalog anh DA KY
// ---------------------------------------------------------------------------

var (
	khoBossKiemMotLan sync.Once
	khoBossLoiKiem    error
)

// idPartDaKy tra danh sach id part hop le tu asset-lock (closure.part_ids) — id
// ma client CHAC CHAN co anh. catalogLangLa.PartIDs la []int; PhaBoss.Outfit la
// [6]int16, nen doi kieu o day chu khong o bo kiem (bo kiem la ham THUAN, khong
// duoc biet gi ve catalog).
func idPartDaKy() []int16 {
	ra := make([]int16, 0, len(catalogLangLa.PartIDs))
	for _, id := range catalogLangLa.PartIDs {
		ra = append(ra, int16(id))
	}
	return ra
}

// kiemOutfitBossMotLan doi chieu outfit cua MOI pha voi danh sach id part da ky.
// Chay MOT LAN cho ca doi may chu — cung ky luat voi kiemIdAnhQuaiMotLan
// (khoquai.go).
//
// PANIC khi lech: file YAML nam TRONG file thuc thi (go:embed) nen "outfit tro
// vao khong khi" khong phai su co moi truong ma la mot ban build hong — im lang
// cho qua la de may chu chay voi mot con boss VO HINH, loai loi khong co mot
// dong log nao va chi lo ra khi co nguoi that vao game nhin.
func kiemOutfitBossMotLan(ac khoBossAnhChup) {
	khoBossKiemMotLan.Do(func() {
		khoBossLoiKiem = noidung.KiemIdPartBoss(ac.DanhSachBoss(), idPartDaKy())
	})
	if khoBossLoiKiem != nil {
		panic("game: kho noi dung boss lech catalog anh da ky: " + khoBossLoiKiem.Error())
	}
}

// ---------------------------------------------------------------------------
// 4. Hai cua tra cuu — GIU NGUYEN chu ky cu
// ---------------------------------------------------------------------------

// khoBossMacDinh tra ve anh chup kho NHUNG — cung mot *noidung.Kho voi
// khoChieuMacDinh (khochieu.go) va khoLangMacDinh (lang.go), qua
// khoNhungMotLan() dung chung sync.Once.
func khoBossMacDinh() khoBossAnhChup {
	return khoNhungMotLan().AnhChup()
}

// timNoiDungBossTu / danhSachNoiDungBossTu — hai cua NHAN ANH CHUP tu ngoai.
// Bo dieu phoi lich dung chung de chot "anh chup nao" o dung khoanh khac sinh
// boss (xem chu thich dau tep ve "boss dang song khong nap nong duoc").
func timNoiDungBossTu(ac khoBossAnhChup, ma string) *NoiDungBoss {
	kiemOutfitBossMotLan(ac)
	b, co := ac.Boss(ma)
	if !co {
		return nil
	}
	cp := noiDungBossTuNoiDung(b)
	return &cp
}

func danhSachNoiDungBossTu(ac khoBossAnhChup) []NoiDungBoss {
	kiemOutfitBossMotLan(ac)
	ds := ac.DanhSachBoss()
	ra := make([]NoiDungBoss, len(ds))
	for i, b := range ds {
		ra[i] = noiDungBossTuNoiDung(b)
	}
	return ra
}

// TimNoiDungBoss tra ban sao cua con boss mang ma nay, hoac nil neu khong co.
//
// Fail-closed (tra nil thay vi mot NoiDungBoss rong): mot NoiDungBoss rong se di
// tiep vao may trang thai va sinh ra con boss khong pha, khong mau, khong ten —
// hong CAM. nil thi noi goi buoc phai kiem.
//
// BAN SAO SAU hai lop: cua Boss() cua anh chup da sao sau, roi
// noiDungBossTuNoiDung dung lat cat MOI. Nguoi goi sua ket qua khong bao gio
// dung toi anh chup dang cong bo — day la loi sao chep nong kinh dien ma du an
// nay da dinh hai lan (banDo.go TimBanDo, va chinh ham nay o ban truoc KH8).
func TimNoiDungBoss(ma string) *NoiDungBoss {
	return timNoiDungBossTu(khoBossMacDinh(), ma)
}

// DanhSachNoiDungBoss tra ban sao SAU toan bo kho noi dung boss (chi doc).
func DanhSachNoiDungBoss() []NoiDungBoss {
	return danhSachNoiDungBossTu(khoBossMacDinh())
}
