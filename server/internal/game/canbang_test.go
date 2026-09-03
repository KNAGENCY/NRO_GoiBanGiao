// canbang_test.go - BO DO CAN BANG D1..D4 (Ke hoach 8, task 11).
//
// Thiet ke: docs/superpowers/specs/2026-09-03-kh8-kho-noi-dung-design.md muc 5.3.
//
// ============================================================================
// DAY LA BO **DO**, KHONG PHAI BO **CHOT NGUONG**.
// ============================================================================
//
// Bon phep do duoi day KHONG khang dinh mot con so nao la "dung" (8-15 giay,
// 2-3 phut...). Nguong la viec cua task 13. O day chi khang dinh ba dieu:
//
//	(1) phep do CHAY DUOC,
//	(2) phep do XAC DINH - chay hai lan ra dung mot con so,
//	(3) bang tom tat co du bon dong.
//
// Rieng TestCanBang_GoldenSoHienTai GHIM con so hien tai, de ai doi mot so
// trong kho-noi-dung/*.yaml (hoac doi cong thuc) thi thay ngay. **Task 13 DUOC
// PHEP cap nhat golden do** - do la muc dich cua no.
//
// ============================================================================
// HAI DIEU KIEN LAM PHEP DO CO NGHIA (muc 5.3)
// ============================================================================
//
//  1. DONG HO TIEM: khong mot `time.Sleep` nao, khong mot `time.Now` nao trong
//     duong do. Thoi gian tran = so nhip actor x 100ms, tinh bang dong ho ghim
//     `k.dongHo` (truong tiem duoc co san tu KH6 task 8, khuvuc.go). Nho vay
//     con so KHONG phu thuoc toc do may - dung bai hoc muc 8 BAT-DAU-TU-DAY
//     (test chop tat vi phu thuoc thoi gian that).
//  2. NGAU NHIEN TIEM: bo do nay KHONG gieo ngau nhien lan nao ca. HP boss lay
//     PHAN TU GIUA cua `hp_chon` (khong random), va vang roi tinh bang KY VONG
//     tu ti le trong bang roi (khong gieo hang nghin lan). Duong roi do
//     (`roiDoKhiChet`) khong bao gio cham bo chon vi khu vuc do khong noi
//     `BoRoiDo` - no thoat som.
//
// KY THUAT chay actor: goi `k.nhipTim()` va `k.xuLy(...)` TRUC TIEP tren mot
// KhuVuc CHUA `BatDau()` - dung khuon cua khuvuc_aiboss_test.go /
// khuvuc_danhboss_test.go: mot goroutine duy nhat, khong hop thu, khong sleep.
// Vi vay bo do KHONG can CSDL (khong DSN), dung nhu yeu cau cua DO-CAN-BANG.bat.
package game

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode/utf8"

	"nroserver/internal/noidung"
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// ---------------------------------------------------------------------------
// Hang so + kieu du lieu cua bo do
// ---------------------------------------------------------------------------

// canBangMoc la "0 giay" cua moi tran do. Ghim cung (KHONG time.Now) de moi
// moc bat tu / hoi chieu / hoi mau la tat dinh.
var canBangMoc = time.Date(2026, 9, 3, 8, 0, 0, 0, time.UTC)

const (
	// canBangNhip la mot nhip actor - dung bang chu ky nhipTim that (100ms).
	canBangNhip = 100 * time.Millisecond

	// canBangTranToiDa la tran dai nhat cho phep truoc khi bo do bo cuoc
	// (6000 nhip = 10 phut trong game). Can mot tran cham hon 10 phut nghia la
	// can bang da hong den muc con so cu the khong con y nghia.
	canBangTranToiDa = 6000

	// canBangTemplateCanChien la template_id cua chieu CAN CHIEN cap 1 - giong
	// nhau o ca 5 lang (kho-noi-dung/chieu.yaml, slot dau cua moi lang).
	canBangTemplateCanChien byte = 1

	// canBangMauVoHan la luong mau gan cho nguoi do trong che do "giet thuan"
	// (xem canBangChayD1): du lon de khong bao gio chet, du nho de khong tran
	// int32 khi tru.
	canBangMauVoHan int32 = 1 << 28

	// canBangIdCuaHang / canBangMauAo: mon do lam moc kinh te cho D4 - "Ao
	// ninja Lang La", cua hang 1 (kho-noi-dung/cua-hang.yaml).
	canBangIdCuaHang int32 = 1
	canBangMauAo     int16 = 2

	// canBangMaSoiRung la ma mau quai dung cho D1/D4 (kho-noi-dung/quai.yaml).
	canBangMaSoiRung = "soi-rung"

	// canBangSoNguoiD2 la so nguoi choi hop nhau danh boss o phep do D2.
	canBangSoNguoiD2 = 5
)

// Toa do dat boss cho phep do D2. Hang y=384 cua map Lang La di duoc suot mot
// doan dai (cung hang ma khuvuc_aiboss_test.go dung), nen boss va nam nguoi
// deu dung duoc o day va khong ai phai di chuyen trong suot tran.
const (
	xCanBangBoss int16 = 500
	yCanBangBoss int16 = 384
)

// canBangPhienIm la phien mang gia NUOT moi goi tin: bo do phat hang chuc nghin
// goi (moi don danh -> phat ca khu vuc), giu lai het se phinh bo nho ma khong
// dung vao viec gi - bo do chi quan tam TRANG THAI (mau, thoi gian), khong
// quan tam byte.
type canBangPhienIm struct{ soGoi int }

func (p *canBangPhienIm) XepGui(*protocol.Message) error { p.soGoi++; return nil }
func (p *canBangPhienIm) DongCham()                      {}

// canBangDongD1 la mot dong ket qua D1/D3/D4 cua MOT lang.
type canBangDongD1 struct {
	IdLang  int16
	TenLang string

	// SatMotDon la sat thuong nguoi cap 1 cua lang nay gay len mot con quai
	// bang chieu can chien cap 1 (tinh boi chinh TinhSatThuong).
	SatMotDon int32
	HpQuai    int32

	// NhipGietThuan: so nhip de ha mot con Soi Rung trong che do "giet thuan"
	// (nguoi do duoc cho mau rat lon de KHONG chet) - day la nhip GIET that su
	// cua lang, dung lam goc tinh D3 (so sanh giua cac lang) va D4 (vang/phut).
	// Quai VAN danh tra va VAN hoi mau binh thuong trong che do nay.
	NhipGietThuan int
	GietDuocThuan bool

	// NhipThat / NguoiChet / QuaiChetThat: kich ban THAT - nguoi choi mang dung
	// HpKhoiDau cua lang, quai danh tra theo AI that.
	NhipThat     int
	NguoiChet    bool
	QuaiChetThat bool

	// D3 - ma tran PvP 5x5 (mo phong cong thuc, xem canBangDoPvP).
	PvPThang, PvPThua, PvPHoa int

	// D4 - kinh te.
	VangMoiPhut float64
	PhutMuaAo   float64
}

// canBangDongD2 la ket qua phep do D2 (5 nguoi ha boss lon pha 1).
type canBangDongD2 struct {
	TenBoss      string
	TenPha       string
	HpBoss       int32
	HpBossConLai int32
	Nhip         int
	BossChet     bool
	SoNguoiSong  int
	SoLanNoiGian int
}

// canBangDongD4 giu cac GIA DINH cua phep do D4 (de in ra bang, vi chung la
// gia dinh chu khong phai su that hien nhien).
type canBangDongD4 struct {
	VangKyVongMoiCon float64
	GiaAo            int32
	ChiTiet          []string
}

// canBangToanBo la ket qua cua ca bon phep do.
type canBangToanBo struct {
	D1 []canBangDongD1
	D2 canBangDongD2
	D4 canBangDongD4
}

// ---------------------------------------------------------------------------
// Chay bon phep do (nho MOT LAN cho ca goi test - moi phep do la tat dinh nen
// nho lai khong lam sai lech gi, chi do lap lai cong viec)
// ---------------------------------------------------------------------------

var (
	canBangMotLan sync.Once
	canBangNho    canBangToanBo
)

func canBangDo() canBangToanBo {
	canBangMotLan.Do(func() { canBangNho = canBangDoLai() })
	return canBangNho
}

// canBangDoLai chay LAI toan bo bon phep do tu dau (khong dung ban nho) - test
// tinh xac dinh goi ham nay hai lan roi so hai ket qua.
func canBangDoLai() canBangToanBo {
	var kq canBangToanBo
	for _, l := range DanhSachLang() {
		kq.D1 = append(kq.D1, canBangDoMotLang(l))
	}
	canBangDoPvP(kq.D1)
	kq.D4 = canBangDoKinhTe()
	canBangDoVangMoiPhut(kq.D1, kq.D4)
	kq.D2 = canBangDoBoss()
	return kq
}

// ---------------------------------------------------------------------------
// D1 - mot nguoi cap 1 giet mot con Soi Rung mat bao lau
// ---------------------------------------------------------------------------

// canBangCapChieuCanChien tra CapChieu cap 1 cua chieu can chien cua mot lang,
// doc tu kho noi dung NHUNG (cung ban YAML voi ban tren dia).
func canBangCapChieuCanChien(idLang int16) *CapChieu {
	return traCapChieuTu(khoChieuMacDinh(), byte(idLang), canBangTemplateCanChien, 1)
}

// canBangSatMotDon tinh sat thuong MOT don cua nguoi cap 1 lang idLang len mot
// muc tieu he nguyen to heMucTieu - dung DUNG ham TinhSatThuong ma actor dung,
// khong co ban sao cong thuc thu hai.
func canBangSatMotDon(idLang int16, heMucTieu byte) int32 {
	lang := TimLang(idLang)
	cap := canBangCapChieuCanChien(idLang)
	if lang == nil || cap == nil {
		return 0
	}
	return TinhSatThuong(int64(lang.SatThuongKhoiDau), cap,
		heNguyenToNguoiChoi(idLang), heMucTieu, false)
}

// canBangDoMotLang chay CA HAI kich ban D1 cho mot lang.
func canBangDoMotLang(l Lang) canBangDongD1 {
	d := canBangDongD1{
		IdLang:    l.ID,
		TenLang:   l.Ten,
		SatMotDon: canBangSatMotDon(l.ID, heQuaiChuaXacDinh),
	}
	d.NhipGietThuan, d.GietDuocThuan, _, d.HpQuai = canBangChayD1(l, true)
	d.NhipThat, d.QuaiChetThat, d.NguoiChet, _ = canBangChayD1(l, false)
	return d
}

// canBangChayD1 mo phong mot nguoi cap 1 cua lang `l` dung tren diem xuat phat
// cua lang do va danh con Soi Rung dau tien trong danh sach quai cua ban do.
//
// mauVoHan=true: nguoi do duoc cho mot luong mau rat lon - quai VAN danh tra
// (nhipTim chay day du, ke ca hoi mau quai 30 giay) nhung khong ha duoc nguoi,
// nen con so thu ve la NHIP GIET THUAN cua lang. Do la con so dung de so sanh
// giua cac lang (D3) va de quy ra vang/phut (D4).
//
// mauVoHan=false: kich ban THAT - mau dung bang HpKhoiDau cua lang.
func canBangChayD1(l Lang, mauVoHan bool) (nhip int, quaiChet, nguoiChet bool, hpQuai int32) {
	bd := TimBanDo(l.MapKhoiDau)
	if bd == nil {
		return 0, false, false, 0
	}
	luc := canBangMoc
	k := MoKhuVuc(bd)
	k.dongHo = func() time.Time { return luc }
	if len(k.quai) == 0 {
		return 0, false, false, 0
	}
	mucTieu := k.quai[0]
	hpQuai = mucTieu.HpToiDa

	hp := l.HpKhoiDau
	if mauVoHan {
		hp = canBangMauVoHan
	}
	nv := &store.NhanVat{
		ID: 1, Ten: "Do", Lang: l.ID,
		MapID:         l.MapKhoiDau,
		HpToiDa:       hp,
		HpHienTai:     hp,
		ChakraToiDa:   l.ChakraKhoiDau,
		ChakraHienTai: l.ChakraKhoiDau,
		ViTriX:        l.ViTriXKhoiDau,
		ViTriY:        l.ViTriYKhoiDau,
	}
	k.xuLy(suKienVao{nv: nv, phien: &canBangPhienIm{}})
	k.xuLy(suKienChonChieu{idNhanVat: nv.ID, mauChieuId: canBangTemplateCanChien})

	for nhip = 1; nhip <= canBangTranToiDa; nhip++ {
		luc = luc.Add(canBangNhip)
		k.nhipTim() // AI quai + hoi mau quai + hoi sinh quai
		if nv.HpHienTai <= 0 {
			return nhip, false, true, hpQuai
		}
		k.xuLy(suKienDanhQuai{idNhanVat: nv.ID, mobId: mucTieu.Id})
		if mucTieu.Hp <= 0 {
			return nhip, true, false, hpQuai
		}
	}
	return canBangTranToiDa, false, false, hpQuai
}

// ---------------------------------------------------------------------------
// D2 - nam nguoi chi so khoi dau danh boss lon pha 1
// ---------------------------------------------------------------------------

// canBangDoBoss dung con boss `nhan-gia-mat-na` PHA 1 voi HP = PHAN TU GIUA cua
// `hp_chon` (tiem chon, KHONG random) va cho 5 nguoi chi so khoi dau Lang La
// vay danh. Boss danh tra bang AI that (nhipBoss -> aiBoss -> bossThuDanh), ke
// ca cac moc NOI GIAN.
func canBangDoBoss() canBangDongD2 {
	var d canBangDongD2

	nd := TimNoiDungBoss(MaBossNhanGiaMatNa)
	if nd == nil || len(nd.Pha) == 0 {
		return d
	}
	pha := nd.Pha[0]
	if len(pha.HpChon) == 0 {
		return d
	}
	d.TenBoss = nd.Ma
	d.TenPha = pha.Ten
	hp := KepHp(pha.HpChon[len(pha.HpChon)/2]) // PHAN TU GIUA - khong gieo
	d.HpBoss = hp

	luc := canBangMoc
	k := MoKhuVuc(TimBanDo(1))
	k.dongHo = func() time.Time { return luc }
	k.quai = nil // do RIENG boss: dan quai cua ban do khong duoc gop vao con so

	the := &TheSongBoss{}
	the.Bat()
	boss := &Boss{
		Id:                 CapIdBoss(),
		PhaHienTai:         0,
		DuLieu:             []PhaBoss{pha},
		Hp:                 hp,
		HpToiDa:            hp,
		ViTriX:             xCanBangBoss,
		ViTriY:             yCanBangBoss,
		TrangThai:          BossHoatDong,
		MocThoiGian:        luc,
		LucHienPha:         luc,
		NhipDanhHienTai:    pha.NhipDanhMs,
		SatThuongTheoNguoi: map[int64]int64{},
		The:                the,
	}
	k.themBoss(boss)

	// Nam nguoi Lang La (lang cua chinh con boss nay) dung sat boss - trong tam
	// chieu can chien cap 1 (dx = 100) nen khong ai phai di chuyen.
	lang := TimLang(0)
	nguoi := make([]*store.NhanVat, 0, canBangSoNguoiD2)
	for i := 0; i < canBangSoNguoiD2; i++ {
		nv := &store.NhanVat{
			ID: int64(101 + i), Ten: fmt.Sprintf("N%d", i), Lang: 0,
			HpToiDa: lang.HpKhoiDau, HpHienTai: lang.HpKhoiDau,
			ChakraToiDa: lang.ChakraKhoiDau, ChakraHienTai: lang.ChakraKhoiDau,
			ViTriX: int32(xCanBangBoss) + int32(i), ViTriY: int32(yCanBangBoss),
		}
		k.xuLy(suKienVao{nv: nv, phien: &canBangPhienIm{}})
		k.xuLy(suKienChonChieu{idNhanVat: nv.ID, mauChieuId: canBangTemplateCanChien})
		nguoi = append(nguoi, nv)
	}

	for d.Nhip = 1; d.Nhip <= canBangTranToiDa; d.Nhip++ {
		luc = luc.Add(canBangNhip)
		k.nhipTim()

		song := 0
		for _, nv := range nguoi {
			if nv.HpHienTai <= 0 {
				continue
			}
			song++
			k.xuLy(suKienDanhNguoiChoi{idNhanVat: nv.ID, idMucTieu: boss.Id})
		}
		if boss.Hp <= 0 {
			d.BossChet = true
			break
		}
		if song == 0 {
			break
		}
	}

	for _, nv := range nguoi {
		if nv.HpHienTai > 0 {
			d.SoNguoiSong++
		}
	}
	d.HpBossConLai = boss.Hp
	d.SoLanNoiGian = boss.MocNoiGianDaQua
	return d
}

// ---------------------------------------------------------------------------
// D3 - ma tran PvP 5x5 (MO PHONG CONG THUC, khong qua actor)
// ---------------------------------------------------------------------------

// canBangDoPvP dien ba o PvPThang/PvPThua/PvPHoa cho tung lang.
//
// CANH BAO: DAY LA MO PHONG CONG THUC, KHONG PHAI TRAN THAT QUA ACTOR (muc 5.3
// cho phep). No dung DUNG ham TinhSatThuong ma actor dung va DUNG phep kep mau
// cua actor (sat thuong khong vuot mau con lai), roi dem so DON can de ha doi
// phuong. Hai ben cap 1, cung chieu can chien, cung hoi chieu 1000ms -> dem so
// don la phep so sanh cong bang. Ai can IT don hon thi thang; bang nhau la hoa.
func canBangDoPvP(ds []canBangDongD1) {
	for i := range ds {
		for j := range ds {
			if i == j {
				continue
			}
			vongA := canBangSoVongHa(ds[i].IdLang, ds[j].IdLang)
			vongB := canBangSoVongHa(ds[j].IdLang, ds[i].IdLang)
			switch {
			case vongA == vongB:
				ds[i].PvPHoa++
			case vongA < vongB:
				ds[i].PvPThang++
			default:
				ds[i].PvPThua++
			}
		}
	}
}

// canBangSoVongHa tra so don ma nguoi cap 1 lang `danh` can de ha nguoi cap 1
// lang `bi`. Tra canBangTranToiDa khi khong bao gio ha duoc (sat thuong 0) -
// mot con so lon hon moi so don that, nen no thua moi doi thu, dung y nghia.
func canBangSoVongHa(danh, bi int16) int {
	sat := canBangSatMotDon(danh, heNguyenToNguoiChoi(bi))
	if sat <= 0 {
		return canBangTranToiDa
	}
	l := TimLang(bi)
	if l == nil {
		return canBangTranToiDa
	}
	return int((l.HpKhoiDau + sat - 1) / sat)
}

// ---------------------------------------------------------------------------
// D4 - bao lau moi mua noi mot mon trong cua hang
// ---------------------------------------------------------------------------

// canBangDoKinhTe tinh KY VONG vang moi con Soi Rung.
//
// GIA DINH LON, PHAI DOC: bang roi hien nay KHONG roi vang truc tiep - no roi
// VAT PHAM (thuoc hoi mau / hoi chakra). Nen "vang kiem duoc" o day la tien
// BAN LAI cho NPC (`gia_ban_lai` trong cua-hang.yaml) cua nhung mon roi ra,
// gia su nguoi choi ban lai het. Mon nao khong co trong cua hang thi KHONG
// tinh vang (khong co cho ban).
//
// Dung KY VONG (ti le x gia) thay vi gieo hang nghin lan: gia tri ky vong tinh
// duoc chinh xac tu ti le, nen gieo chi them nhieu va them ngau nhien.
func canBangDoKinhTe() canBangDongD4 {
	var d canBangDongD4
	ac := khoChieuMacDinh()

	gia := map[int16]int32{}
	if ch, co := ac.CuaHang(canBangIdCuaHang); co {
		for _, tab := range ch.Tabs {
			for _, mon := range tab.Mon {
				if mon.Tien == noidung.TienVang {
					gia[mon.MauId] = mon.GiaBanLai
				}
				if mon.MauId == canBangMauAo {
					d.GiaAo = mon.Gia
				}
			}
		}
	}

	for _, m := range mucRoiTuNoiDung(ac.BangRoiCuaQuai(canBangMaSoiRung), false) {
		banLai, co := gia[int16(m.MauID)]
		ten := "mau " + fmt.Sprint(m.MauID)
		if mau, okMau := ac.TimMauVatPham(int(m.MauID)); okMau {
			ten = mau.Ten
		}
		// TiLe la PHAN VAN (mau so roiDoMauSo = 10000, xem roido.go) - chia dung
		// mau so do, KHONG chia 100. Chia nham la con so kinh te lech 100 lan.
		phanTram := float64(m.TiLe) / float64(roiDoMauSo) * 100
		if !co {
			d.ChiTiet = append(d.ChiTiet,
				fmt.Sprintf("%s: roi %s%%, KHONG co cho ban lai -> tinh 0 vang", ten, canBangSo(phanTram, 1)))
			continue
		}
		gopVang := float64(m.TiLe) / float64(roiDoMauSo) * float64(m.SoLuong) * float64(banLai)
		d.VangKyVongMoiCon += gopVang
		d.ChiTiet = append(d.ChiTiet,
			fmt.Sprintf("%s: roi %s%% x %d cai x ban lai %d vang = %s vang/con",
				ten, canBangSo(phanTram, 1), m.SoLuong, banLai, canBangSo(gopVang, 2)))
	}
	return d
}

// canBangDoVangMoiPhut quy nhip giet thuan cua tung lang thanh vang/phut va so
// phut can de mua noi cai ao.
func canBangDoVangMoiPhut(ds []canBangDongD1, kt canBangDongD4) {
	for i := range ds {
		if !ds[i].GietDuocThuan || ds[i].NhipGietThuan <= 0 {
			continue
		}
		giay := float64(ds[i].NhipGietThuan) * canBangNhip.Seconds()
		conMoiPhut := 60.0 / giay
		ds[i].VangMoiPhut = conMoiPhut * kt.VangKyVongMoiCon
		if ds[i].VangMoiPhut > 0 && kt.GiaAo > 0 {
			ds[i].PhutMuaAo = float64(kt.GiaAo) / ds[i].VangMoiPhut
		}
	}
}

// ---------------------------------------------------------------------------
// Bang tom tat - mot file text, doc duoc boi nguoi khong lap trinh
// ---------------------------------------------------------------------------

// canBangSo dinh dang mot so thap phan theo kieu Viet (dau phay thap phan).
func canBangSo(v float64, leSo int) string {
	return strings.Replace(fmt.Sprintf("%.*f", leSo, v), ".", ",", 1)
}

// canBangDoRongCot la be rong (tinh bang KY TU, khong phai byte) cua sau cot
// trong bang tom tat. Cot am = can trai, duong = can phai.
var canBangDoRongCot = [6]int{-16, 16, 9, -30, 12, 12}

// canBangKhung dat sau o cua mot dong vao dung cot.
//
// VI SAO KHONG DUNG `%-16s` cua fmt: fmt dem BYTE. "L脿ng S瓢啤ng M霉" la 13 ky tu
// nhung 18 byte, nen moi ten lang co dau bi keo lech mot doan khac nhau va ca
// bang nhin nhu vo. O day dem bang utf8.RuneCountInString.
func canBangKhung(o []string) string {
	var b strings.Builder
	for i, s := range o {
		if i > 0 {
			b.WriteByte(' ')
		}
		rong := 0
		if i < len(canBangDoRongCot) {
			rong = canBangDoRongCot[i]
		}
		dem := utf8.RuneCountInString(s)
		switch {
		case rong < 0 && dem < -rong:
			b.WriteString(s + strings.Repeat(" ", -rong-dem))
		case rong > 0 && dem < rong:
			b.WriteString(strings.Repeat(" ", rong-dem) + s)
		default:
			b.WriteString(s)
		}
	}
	return strings.TrimRight(b.String(), " ")
}

// canBangGiay doi so nhip actor thanh chuoi "x,y giay".
func canBangGiay(nhip int) string {
	return canBangSo(float64(nhip)*canBangNhip.Seconds(), 1) + " giay"
}

// canBangPhutGiay doi so nhip actor thanh "x phut y giay" (cho tran dai).
func canBangPhutGiay(nhip int) string {
	tong := float64(nhip) * canBangNhip.Seconds()
	phut := int(tong) / 60
	giay := tong - float64(phut*60)
	if phut == 0 {
		return canBangSo(giay, 1) + " giay"
	}
	return fmt.Sprintf("%d phut %s giay", phut, canBangSo(giay, 1))
}

// canBangDungBang dung BANG TOM TAT (mau o muc 5.3) tu ket qua bon phep do.
func canBangDungBang() string {
	kq := canBangDo()
	var b strings.Builder

	fmt.Fprintf(&b, "BANG CAN BANG - the he kho noi dung %d - %s\n",
		khoChieuMacDinh().TheHe(), time.Now().Format("02/01/2006"))
	b.WriteString(strings.Repeat("=", 78) + "\n\n")
	b.WriteString("Bang nay do ra tu chinh so lieu trong server/kho-noi-dung/*.yaml.\n")
	b.WriteString("Sua mot so trong do -> bam dup server/DO-CAN-BANG.bat -> doc lai bang nay.\n")
	b.WriteString("Bo do KHONG phan xet so nao dung/sai - no chi noi hien nay dang la bao nhieu.\n\n")

	// D1 + D3 + D4 theo tung lang.
	b.WriteString("D1 - MOT NGUOI MOI (cap 1) GIET MOT CON SOI RUNG MAT BAO LAU?\n")
	b.WriteString("D3 - LANG NAO THANG? (so sanh D1 giua cac lang + doi khang PvP)\n")
	b.WriteString("D4 - DI DANH QUAI BAO LAU MOI MUA NOI MOT MON?\n\n")
	dong := func(o ...string) {
		b.WriteString(canBangKhung(o) + "\n")
	}
	dong("Lang", "D1 giet thuan", "moi don", "D1 kich ban that", "D3 PvP", "D4 mua ao")
	dong(strings.Repeat("-", 16), strings.Repeat("-", 16), strings.Repeat("-", 9),
		strings.Repeat("-", 30), strings.Repeat("-", 12), strings.Repeat("-", 12))

	for _, d := range kq.D1 {
		giet := "khong giet noi"
		if d.GietDuocThuan {
			giet = canBangPhutGiay(d.NhipGietThuan)
		}
		that := "hoa (het gio do)"
		switch {
		case d.NguoiChet:
			that = "NGUOI CHOI CHET sau " + canBangGiay(d.NhipThat)
		case d.QuaiChetThat:
			that = "giet duoc sau " + canBangGiay(d.NhipThat)
		}
		muaAo := "khong tinh duoc"
		if d.PhutMuaAo > 0 {
			muaAo = canBangSo(d.PhutMuaAo, 0) + " phut"
		}
		dong(d.TenLang, giet, fmt.Sprintf("%d mau", d.SatMotDon), that,
			fmt.Sprintf("%dT/%dB/%dH", d.PvPThang, d.PvPThua, d.PvPHoa), muaAo)
	}

	b.WriteString("\n  Cot \"D1 giet thuan\": nguoi do duoc cho mau rat lon de KHONG chet, quai van\n")
	b.WriteString("  danh tra va van hoi mau - do la nhip giet THAT SU cua lang, dung lam goc cho\n")
	b.WriteString("  cot D3 va cot D4.\n")
	b.WriteString("  Cot \"D1 kich ban that\": nguoi choi mang dung so mau khoi dau cua lang minh,\n")
	b.WriteString("  dung tren diem xuat phat, HAI con Soi Rung cua ban do danh tra theo AI that.\n")
	b.WriteString("  Cot \"D3 PvP\": T = thang, B = thua, H = hoa, tren 4 doi thu. CANH BAO: day la\n")
	b.WriteString("  MO PHONG CONG THUC (dung ham tinh sat thuong that + phep kep mau that, dem so\n")
	b.WriteString("  don can de ha doi phuong), KHONG phai mot tran that chay qua actor.\n\n")

	// D2.
	b.WriteString("D2 - NAM NGUOI CHI SO KHOI DAU HA BOSS LON PHA 1 MAT BAO LAU?\n\n")
	if kq.D2.HpBoss == 0 {
		b.WriteString("  Khong do duoc: khong doc duoc du lieu boss.\n\n")
	} else {
		fmt.Fprintf(&b, "  Boss:            %s - pha 1 \"%s\"\n", kq.D2.TenBoss, kq.D2.TenPha)
		fmt.Fprintf(&b, "  Mau boss:        %d (phan tu giua cua hp_chon, khong gieo ngau nhien)\n", kq.D2.HpBoss)
		fmt.Fprintf(&b, "  Nguoi danh:      %d nguoi Lang La, chi so khoi dau, chieu can chien cap 1\n", canBangSoNguoiD2)
		if kq.D2.BossChet {
			fmt.Fprintf(&b, "  KET QUA:         HA DUOC BOSS sau %s\n", canBangPhutGiay(kq.D2.Nhip))
		} else {
			fmt.Fprintf(&b, "  KET QUA:         KHONG HA NOI - ca %d nguoi guc sau %s\n",
				canBangSoNguoiD2, canBangPhutGiay(kq.D2.Nhip))
			fmt.Fprintf(&b, "                   Boss con %d / %d mau khi nguoi cuoi cung guc.\n",
				kq.D2.HpBossConLai, kq.D2.HpBoss)
		}
		fmt.Fprintf(&b, "  Con song cuoi:   %d / %d nguoi\n", kq.D2.SoNguoiSong, canBangSoNguoiD2)
		fmt.Fprintf(&b, "  Boss noi gian:   %d lan\n\n", kq.D2.SoLanNoiGian)
	}

	// D4 - gia dinh.
	b.WriteString("D4 - GIA DINH CUA PHEP DO KINH TE (doc ky truoc khi tin cot \"D4 mua ao\")\n\n")
	b.WriteString("  CANH BAO: bang roi hien nay KHONG roi vang truc tiep - no roi VAT PHAM. Nen\n")
	b.WriteString("  \"vang kiem duoc\" o day la tien BAN LAI cho NPC (gia_ban_lai) cua nhung mon\n")
	b.WriteString("  roi ra, gia su nguoi choi ban lai het sach.\n\n")
	for _, c := range kq.D4.ChiTiet {
		fmt.Fprintf(&b, "  - %s\n", c)
	}
	fmt.Fprintf(&b, "  => Ky vong: %s vang moi con Soi Rung.\n", canBangSo(kq.D4.VangKyVongMoiCon, 2))
	fmt.Fprintf(&b, "  => Mon do lam moc: \"Ao ninja Lang La\", gia %d vang.\n\n", kq.D4.GiaAo)

	b.WriteString(strings.Repeat("=", 78) + "\n")
	b.WriteString("Bang do bo do sinh ra tu internal/game/canbang_test.go (KH8 task 11).\n")
	return b.String()
}

// ---------------------------------------------------------------------------
// BON PHEP DO - moi phep do la MOT test co ten
// ---------------------------------------------------------------------------

// TestCanBang_D1_GietSoiRung do D1 cho ca 5 lang.
//
// KHONG khang dinh nguong 8-15 giay - nguong la viec task 13. Chi khang dinh
// phep do CHAY va doc dung nguon (HP quai lay tu quai.yaml, khong go cung 100).
func TestCanBang_D1_GietSoiRung(t *testing.T) {
	kq := canBangDo()
	if len(kq.D1) != 5 {
		t.Fatalf("D1 do duoc %d lang, mong 5", len(kq.D1))
	}

	mau, co := khoChieuMacDinh().Quai(canBangMaSoiRung)
	if !co {
		t.Fatalf("khong tim thay mau quai %q trong kho noi dung", canBangMaSoiRung)
	}
	for _, d := range kq.D1 {
		if d.HpQuai != int32(mau.Hp) {
			t.Errorf("%s: HP quai do duoc = %d, kho noi dung khai %d - phep do dang doc nham nguon",
				d.TenLang, d.HpQuai, mau.Hp)
		}
		if d.NhipGietThuan <= 0 {
			t.Errorf("%s: khong do duoc nhip giet thuan", d.TenLang)
		}
		t.Logf("D1 %-16s sat mot don=%d  giet thuan=%s  kich ban that=%s (nguoiChet=%v quaiChet=%v)",
			d.TenLang, d.SatMotDon, canBangPhutGiay(d.NhipGietThuan),
			canBangGiay(d.NhipThat), d.NguoiChet, d.QuaiChetThat)
	}
}

// TestCanBang_D2_NamNguoiDanhBossLon do D2.
func TestCanBang_D2_NamNguoiDanhBossLon(t *testing.T) {
	d := canBangDo().D2
	if d.HpBoss == 0 {
		t.Fatal("D2 khong do duoc - khong dung duoc boss tu kho noi dung")
	}
	nd := TimNoiDungBoss(MaBossNhanGiaMatNa)
	giua := KepHp(nd.Pha[0].HpChon[len(nd.Pha[0].HpChon)/2])
	if d.HpBoss != giua {
		t.Errorf("HP boss dung trong phep do = %d, phan tu giua cua hp_chon = %d", d.HpBoss, giua)
	}
	if d.Nhip <= 0 {
		t.Error("D2 khong chay duoc nhip nao")
	}
	t.Logf("D2 boss=%s pha=%q hp=%d ket thuc sau %s (bossChet=%v, boss con %d mau, con song %d/%d, noi gian %d lan)",
		d.TenBoss, d.TenPha, d.HpBoss, canBangPhutGiay(d.Nhip), d.BossChet,
		d.HpBossConLai, d.SoNguoiSong, canBangSoNguoiD2, d.SoLanNoiGian)
}

// TestCanBang_D3_LangNaoThang do D3: chenh lech D1 giua lang nhanh nhat va
// cham nhat, cong ma tran PvP 5x5.
func TestCanBang_D3_LangNaoThang(t *testing.T) {
	kq := canBangDo()
	nhanh, cham := kq.D1[0], kq.D1[0]
	for _, d := range kq.D1 {
		if d.NhipGietThuan < nhanh.NhipGietThuan {
			nhanh = d
		}
		if d.NhipGietThuan > cham.NhipGietThuan {
			cham = d
		}
	}
	if nhanh.NhipGietThuan <= 0 {
		t.Fatal("khong co lang nao do duoc nhip giet")
	}
	chenh := float64(cham.NhipGietThuan-nhanh.NhipGietThuan) / float64(nhanh.NhipGietThuan) * 100
	t.Logf("D3 nhanh nhat=%s (%s), cham nhat=%s (%s), chenh lech %s%%",
		nhanh.TenLang, canBangPhutGiay(nhanh.NhipGietThuan),
		cham.TenLang, canBangPhutGiay(cham.NhipGietThuan), canBangSo(chenh, 1))

	// Moi lang phai danh du 4 tran PvP (khong tu danh minh).
	for _, d := range kq.D1 {
		if n := d.PvPThang + d.PvPThua + d.PvPHoa; n != 4 {
			t.Errorf("%s: ma tran PvP co %d tran, mong 4", d.TenLang, n)
		}
		t.Logf("D3 PvP %-16s %d thang / %d thua / %d hoa", d.TenLang, d.PvPThang, d.PvPThua, d.PvPHoa)
	}
}

// TestCanBang_D4_MuaNoiMotMon do D4 bang KY VONG (khong gieo ngau nhien).
func TestCanBang_D4_MuaNoiMotMon(t *testing.T) {
	kq := canBangDo()
	if kq.D4.GiaAo <= 0 {
		t.Fatal("khong doc duoc gia \"Ao ninja Lang La\" tu cua-hang.yaml")
	}
	if len(kq.D4.ChiTiet) == 0 {
		t.Fatal("bang roi cua soi-rung rong - D4 khong co gi de tinh")
	}

	// Ky vong phai khop voi phep tinh tay tu chinh bang roi + gia ban lai:
	// day la chot chong dot bien "quen nhan gia_ban_lai".
	ac := khoChieuMacDinh()
	ch, co := ac.CuaHang(canBangIdCuaHang)
	if !co {
		t.Fatalf("khong co cua hang %d", canBangIdCuaHang)
	}
	banLai := map[int16]int32{}
	for _, tab := range ch.Tabs {
		for _, mon := range tab.Mon {
			if mon.Tien == noidung.TienVang {
				banLai[mon.MauId] = mon.GiaBanLai
			}
		}
	}
	var tay float64
	for _, m := range mucRoiTuNoiDung(ac.BangRoiCuaQuai(canBangMaSoiRung), false) {
		tay += float64(m.TiLe) / float64(roiDoMauSo) * float64(m.SoLuong) * float64(banLai[int16(m.MauID)])
	}
	if tay != kq.D4.VangKyVongMoiCon {
		t.Errorf("ky vong vang/con = %v, tinh tay tu bang roi x gia_ban_lai = %v",
			kq.D4.VangKyVongMoiCon, tay)
	}

	for _, d := range kq.D1 {
		t.Logf("D4 %-16s %s vang/phut -> %s phut moi mua noi ao %d vang",
			d.TenLang, canBangSo(d.VangMoiPhut, 2), canBangSo(d.PhutMuaAo, 0), kq.D4.GiaAo)
	}
}

// ---------------------------------------------------------------------------
// Ba chot bao ve chinh bo do
// ---------------------------------------------------------------------------

// TestCanBang_PhepDoXacDinh chay LAI ca bon phep do tu dau va doi hoi ket qua y
// HET lan truoc. Day la chot chong "dung time.Now / dung ngau nhien khong
// tiem" - mot bo do phu thuoc dong ho that hay may man se lech ngay.
func TestCanBang_PhepDoXacDinh(t *testing.T) {
	a := canBangDoLai()
	b := canBangDoLai()
	if fmt.Sprint(a) != fmt.Sprint(b) {
		t.Fatalf("hai lan do ra hai ket qua khac nhau - phep do KHONG xac dinh\nlan 1: %v\nlan 2: %v", a, b)
	}
}

// TestCanBang_BangCoDuBonPhepDo giu dieu kien toi thieu: mot bang, du bon phep
// do, mot dong cho moi lang.
func TestCanBang_BangCoDuBonPhepDo(t *testing.T) {
	b := canBangDungBang()
	for _, phai := range []string{"D1 -", "D2 -", "D3 -", "D4 -"} {
		if !strings.Contains(b, phai) {
			t.Errorf("bang tom tat thieu phan %q", phai)
		}
	}
	for _, l := range DanhSachLang() {
		if !strings.Contains(b, l.Ten) {
			t.Errorf("bang tom tat thieu dong cua lang %q", l.Ten)
		}
	}
}

// TestCanBang_GhiBangTomTat GHI bang ra server/bang-can-bang.txt - day la thu
// ma chu du an doc sau khi bam dup DO-CAN-BANG.bat.
//
// Duong dan tinh tu thu muc goi test (internal/game) len goc module server.
func TestCanBang_GhiBangTomTat(t *testing.T) {
	duong := filepath.Join("..", "..", "bang-can-bang.txt")
	if err := os.WriteFile(duong, []byte(canBangDungBang()), 0o644); err != nil {
		t.Fatalf("khong ghi duoc bang tom tat: %v", err)
	}
	tuyetDoi, _ := filepath.Abs(duong)
	t.Logf("da ghi bang tom tat: %s", tuyetDoi)
}

// TestCanBang_GoldenSoHienTai GHIM con so HIEN NAY.
//
// TASK 13 (vong can bang) DUOC PHEP - va co nhiem vu phai - CAP NHAT cac con so
// trong test nay khi no doi so trong kho-noi-dung/*.yaml. Muc dich cua test
// khong phai "cam doi", ma la "doi thi phai THAY minh vua doi cai gi".
//
// Test nay chinh la thu bat duoc ba loai dot bien am tham: doc HP quai go cung
// thay vi doc quai.yaml, bo dong ho tiem dung time.Now, va quen nhan
// gia_ban_lai trong D4.
func TestCanBang_GoldenSoHienTai(t *testing.T) {
	kq := canBangDo()

	mongD1 := []struct {
		ten           string
		satMotDon     int32
		nhipGietThuan int
		nguoiChet     bool
		nhipThat      int
	}{
		{"Làng Lá", 2, 491, true, 191},
		{"Làng Cát", 2, 491, true, 121},
		{"Làng Sương Mù", 3, 331, true, 161},
		{"Làng Đá", 2, 491, true, 321},
		{"Làng Mây", 2, 491, true, 151},
	}
	if len(kq.D1) != len(mongD1) {
		t.Fatalf("co %d lang, golden ghim %d", len(kq.D1), len(mongD1))
	}
	for i, m := range mongD1 {
		d := kq.D1[i]
		if d.TenLang != m.ten {
			t.Errorf("lang #%d ten = %q, golden %q", i, d.TenLang, m.ten)
			continue
		}
		if d.SatMotDon != m.satMotDon {
			t.Errorf("%s: sat mot don = %d, golden %d", d.TenLang, d.SatMotDon, m.satMotDon)
		}
		if d.NhipGietThuan != m.nhipGietThuan {
			t.Errorf("%s: nhip giet thuan = %d, golden %d", d.TenLang, d.NhipGietThuan, m.nhipGietThuan)
		}
		if d.NguoiChet != m.nguoiChet {
			t.Errorf("%s: nguoi chet o kich ban that = %v, golden %v", d.TenLang, d.NguoiChet, m.nguoiChet)
		}
		if d.NhipThat != m.nhipThat {
			t.Errorf("%s: nhip kich ban that = %d, golden %d", d.TenLang, d.NhipThat, m.nhipThat)
		}
	}

	if kq.D4.VangKyVongMoiCon != 3.75 {
		t.Errorf("D4 ky vong vang moi con = %v, golden 3,75", kq.D4.VangKyVongMoiCon)
	}
	if kq.D4.GiaAo != 1000 {
		t.Errorf("D4 gia ao = %d, golden 1000", kq.D4.GiaAo)
	}
	if kq.D2.HpBoss != 3000 {
		t.Errorf("D2 hp boss = %d, golden 3000", kq.D2.HpBoss)
	}
	if kq.D2.BossChet {
		t.Error("D2 golden hien tai: KHONG ha noi boss - nay lai ha duoc, hay cap nhat golden")
	}
}
