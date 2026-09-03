// KHO — cong bo mot ANH CHUP noi dung cho ca may chu, va trao anh chup moi khi
// nap nong, ma KHONG khoa va KHONG lam ai dang doc thay du lieu sua do.
//
// (KH8 task 2. Thiet ke: docs/superpowers/specs/2026-09-03-kh8-kho-noi-dung-design.md
// §3.2 "anh chup bat bien + atomic.Pointer", §4 K4, dot bien M12/M13.)
//
// # VI SAO KHONG DUNG RWMutex, KHONG DUNG SU KIEN
//
// RWMutex: mot lan nap nong phai doi TAT CA nguoi doc nha khoa. Actor cua may
// chu nay co the ket o socket 30 giay (§6.10b) — nghia la mot lan nap nong co
// the treo ca may chu 30 giay. Loai.
//
// Gui su kien vao hop thu tung actor: hop thu day thi tin ROI IM LANG (bai hoc
// KH6c) — nap nong "thanh cong" nhung vai khu vuc van chay so cu, khong ai biet.
// Loai.
//
// Con lai: dung mot con tro nguyen tu. Nguoi nap dung TRON mot anh chup moi o
// noi khac, roi trao con tro trong MOT lenh may. Nguoi doc hoac thay ban cu tron
// ven, hoac thay ban moi tron ven — khong bao gio thay nua no nua kia.
//
// # HOP DONG VOI ACTOR: LAY CON TRO MOT LAN O DAU NHIP
//
// atomic.Pointer chi bao dam moi LAN DOC la nguyen tu. Goi AnhChup() hai lan
// giua mot nhip thi hai lan do co the roi vao hai the he khac nhau — mot don
// danh se tinh sat thuong bang bang so nay va tru chakra bang bang so kia. Vi
// vay hop dong la:
//
//	func (kv *KhuVuc) nhipTim() {
//	    ac := kv.kho.AnhChup()      // MOT LAN, dong dau tien cua nhip
//	    for _, nv := range kv.nhanVat {
//	        chieu, co := ac.TraCapChieu(nv.IdLang, nv.MauChieu, nv.Cap)
//	        ...                     // dung ac suot nhip, KHONG goi AnhChup() lai
//	    }
//	}
//
// Nap nong xay ra giua nhip thi nhip do chay tron bang so cu; nhip sau moi thay
// so moi. Do la dieu MONG MUON, khong phai thoa hiep.
//
// # ANH CHUP BAT BIEN — CUONG CHE BANG BAN SAO SAU
//
// atomic.Pointer KHONG lam du lieu ben trong bat bien; no chi lam con tro
// nguyen tu. Neu mot cua tra cuu tra thang slice cua anh chup ra ngoai va nguoi
// goi lo `append` hay ghi vao mot phan tu, ho da sua ANH CHUP DA CONG BO — moi
// nguoi choi khac thay theo. Du an nay da dinh dung bay do hai lan
// (game/banDo.go TimBanDo, game/noidungboss.go saoNoiDungBoss).
//
// Vi vay: MOI cua tra cuu o day tra BAN SAO SAU. Chi phi da do (kho_test.go,
// BenchmarkKhoChieuCuaLang / BenchmarkKhoCuaHang) — xem ghi chu tren tung ham.
package noidung

import (
	"sync/atomic"
	"time"
)

// anhChup la mot the he noi dung DA CONG BO. Bat bien sau khi tao: khong mot ham
// nao trong tep nay ghi vao truong cua no sau khi MoKho/Thay tra ve.
//
// Kieu khong xuat khau CO Y: goi `game` chi duoc cam no bang `ac := kho.AnhChup()`
// va goi phuong thuc — khong the tu dung mot cai gia, khong the doc thang
// `ac.duLieu` de lach ban sao sau.
type anhChup struct {
	duLieu *DuLieu
	theHe  uint64
	luc    time.Time
}

// TheHe la so the he cua anh chup nay. Tang dan 1 moi lan trao. Vao log moi lan
// nap nong de con tra loi duoc "luc 21:47 may chu dang chay bo so nao".
func (a *anhChup) TheHe() uint64 { return a.theHe }

// Luc la thoi diem anh chup nay duoc cong bo.
func (a *anhChup) Luc() time.Time { return a.luc }

// Kho giu anh chup dang co hieu luc. Zero value KHONG dung duoc — phai qua MoKho.
type Kho struct {
	hienTai atomic.Pointer[anhChup]
}

// MoKho mo kho voi anh chup dau tien (the he 1). dl phai KHAC nil.
//
// Sau khi giao dl cho MoKho, nguoi goi KHONG duoc giu lai va sua dl nua — no da
// tro thanh du lieu bat bien cua the he 1.
func MoKho(dl *DuLieu) *Kho {
	if dl == nil {
		panic("noidung: MoKho voi du lieu nil")
	}
	k := &Kho{}
	k.hienTai.Store(&anhChup{duLieu: dl, theHe: 1, luc: time.Now()})
	return k
}

// Thay trao ca kho sang mot anh chup moi va tra ve so the he moi.
//
// Nguyen tu, khong khoa, khong bao gio doi ai. Goroutine nao dang cam anh chup
// cu van doc no binh thuong cho het viec — bo thu gom rac don khi khong con ai
// cam nua.
//
// LUU Y (dot bien M12): dung bao gio "nap nong" bang cach sua truong cua
// *DuLieu dang chay. Duong duy nhat la dung mot *DuLieu HOAN TOAN MOI roi Thay.
func (k *Kho) Thay(dl *DuLieu) (theHeMoi uint64) {
	if dl == nil {
		panic("noidung: Thay voi du lieu nil")
	}
	for {
		cu := k.hienTai.Load()
		moi := &anhChup{duLieu: dl, theHe: cu.theHe + 1, luc: time.Now()}
		if k.hienTai.CompareAndSwap(cu, moi) {
			return moi.theHe
		}
	}
}

// AnhChup tra ve anh chup dang co hieu luc.
//
// HOP DONG: actor goi ham nay DUNG MOT LAN o dau nhipTim, giu ket qua suot nhip,
// KHONG goi lai giua nhip. Xem chu thich dau tep.
func (k *Kho) AnhChup() *anhChup { return k.hienTai.Load() }

// TheHe la so the he dang co hieu luc. Tien cho log/panel; trong nhip hay dung
// AnhChup().TheHe() de chac chan cung mot anh chup.
func (k *Kho) TheHe() uint64 { return k.hienTai.Load().theHe }

// ---------------------------------------------------------------------------
// Cua tra cuu — TAT CA tra ban sao sau
// ---------------------------------------------------------------------------

// ChieuCuaLang tra ve ban sao sau danh sach chieu cua mot lang. Lang ngoai
// 0..SoLangCoDinh-1 tra nil.
//
// CHI PHI — DA DO (BenchmarkKhoChieuCuaLangThatCo, 30 chieu x 10 cap, may dev):
//
//	ChieuCuaLang (ban sao sau)   ~2 950 ns  (19,6 KB, 31 cap phat)
//	TraCapChieu  (mot cap)          ~2,1 ns (0 cap phat)
//	CuaHang      (ban sao 3 lop)   ~54   ns (0,1 KB, 3 cap phat)
//	AnhChup      (doc con tro)      ~0,25 ns
//
// Nhip la 100 ms = 100 000 000 ns. Ngay ca goi ChieuCuaLang mot lan moi nhip
// cho MOI trong 100 khu vuc van la ~0,3% cua nhip. Vi vay CHON BAN SAO SAU, ke
// ca khi con tro chi-doc re hon — vi "chung minh khong ai ghi" la loi hua ma
// tuong lai se pha (TimBanDo va saoNoiDungBoss da pha dung loi hua do mot lan
// moi), con ban sao thi khong the pha bang mot dong code vo y.
//
// Trong nhip hay dung TraCapChieu (0 cap phat); ChieuCuaLang la cua danh cho
// luc gui catalog / dung goi -28/7, khong phai cua goi moi don danh.
func (a *anhChup) ChieuCuaLang(idLang byte) []Chieu {
	if int(idLang) >= SoLangCoDinh {
		return nil
	}
	return khoSaoChieu(a.duLieu.ChieuTheoLang[idLang])
}

// TraCapChieu tra ve ban sao MOT cap cua mot chieu. CapChieu chi co truong vo
// huong va mot chuoi (chuoi Go bat bien) — gan gia tri la da sao xong.
//
// `cap` dem tu 1 nhu trong game.
func (a *anhChup) TraCapChieu(idLang, templateId, cap byte) (CapChieu, bool) {
	if int(idLang) >= SoLangCoDinh || cap == 0 {
		return CapChieu{}, false
	}
	for _, c := range a.duLieu.ChieuTheoLang[idLang] {
		if c.TemplateId != templateId {
			continue
		}
		if int(cap) > len(c.CacCap) {
			return CapChieu{}, false
		}
		return c.CacCap[cap-1], true
	}
	return CapChieu{}, false
}

// MauOption tra ve ban sao bang mau option. MauOption toan truong vo huong +
// chuoi nen mot lan copy slice la du sau.
func (a *anhChup) MauOption() []MauOption {
	if a.duLieu.MauOption == nil {
		return nil
	}
	ds := make([]MauOption, len(a.duLieu.MauOption))
	copy(ds, a.duLieu.MauOption)
	return ds
}

// MauVatPham tra ve ban sao toan bo bang mau vat pham. Bang nay co the vai nghin
// dong — dung goi trong vong lap moi nhip; de gui catalog luc dang nhap thoi.
func (a *anhChup) MauVatPham() []MauVatPham {
	if a.duLieu.MauVatPham == nil {
		return nil
	}
	ds := make([]MauVatPham, len(a.duLieu.MauVatPham))
	copy(ds, a.duLieu.MauVatPham)
	return ds
}

// TimMauVatPham tra ve ban sao MOT mau vat pham theo Id (khong phai chi so
// mang). Day la cua dung trong nhip: gan gia tri, khong cap phat.
func (a *anhChup) TimMauVatPham(id int) (MauVatPham, bool) {
	for _, m := range a.duLieu.MauVatPham {
		if m.Id == id {
			return m, true
		}
	}
	return MauVatPham{}, false
}

// SoMauVatPham / SoMauOption — dem ma khong sao chep gi.
func (a *anhChup) SoMauVatPham() int { return len(a.duLieu.MauVatPham) }
func (a *anhChup) SoMauOption() int  { return len(a.duLieu.MauOption) }

// CuaHang tra ve ban sao sau mot cua hang theo Id. Sao ba lop: tabs -> mon ->
// options. Bo sot mot lop la nguoi goi sua duoc gia hang cho MOI nguoi choi.
func (a *anhChup) CuaHang(id int32) (CuaHang, bool) {
	for _, ch := range a.duLieu.CuaHang {
		if ch.Id == id {
			return khoSaoMotCuaHang(ch), true
		}
	}
	return CuaHang{}, false
}

// DanhSachCuaHang tra ve ban sao sau toan bo cua hang.
func (a *anhChup) DanhSachCuaHang() []CuaHang {
	if a.duLieu.CuaHang == nil {
		return nil
	}
	ds := make([]CuaHang, len(a.duLieu.CuaHang))
	for i, ch := range a.duLieu.CuaHang {
		ds[i] = khoSaoMotCuaHang(ch)
	}
	return ds
}

// ---------------------------------------------------------------------------
// Sao chep sau — dung chung
// ---------------------------------------------------------------------------

func khoSaoChieu(ds []Chieu) []Chieu {
	if ds == nil {
		return nil
	}
	ra := make([]Chieu, len(ds))
	for i, c := range ds {
		ra[i] = c
		if c.CacCap != nil {
			ra[i].CacCap = make([]CapChieu, len(c.CacCap))
			copy(ra[i].CacCap, c.CacCap)
		}
	}
	return ra
}

func khoSaoMotCuaHang(ch CuaHang) CuaHang {
	ra := ch
	if ch.Tabs == nil {
		return ra
	}
	ra.Tabs = make([]TabCuaHang, len(ch.Tabs))
	for i, tab := range ch.Tabs {
		ra.Tabs[i] = tab
		if tab.Mon == nil {
			continue
		}
		ra.Tabs[i].Mon = make([]MonHang, len(tab.Mon))
		for j, mon := range tab.Mon {
			ra.Tabs[i].Mon[j] = mon
			if mon.Options != nil {
				ra.Tabs[i].Mon[j].Options = make([]OptionMon, len(mon.Options))
				copy(ra.Tabs[i].Mon[j].Options, mon.Options)
			}
		}
	}
	return ra
}
