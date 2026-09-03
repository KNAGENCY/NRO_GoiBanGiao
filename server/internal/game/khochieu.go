// khochieu.go — CUA DUY NHAT cua goi `game` di vao KHO NOI DUNG cho phan CHIEU.
//
// (KH8 task 3. Thiet ke: docs/superpowers/specs/2026-09-03-kh8-kho-noi-dung-design.md
// hang "3" bang Task, §4 nap nong, §6.9 ky luat actor, §7 GATE.)
//
// # CAI GI DOI O TASK NAY
//
// Truoc task nay, 5 lang x 4 chieu x 3 cap duoc SINH RA BANG CODE GO
// (chieuthuc.go: xayDungChieuTheoLang + taoCapChieu) va nam trong mot bien
// toan cuc `chieuTheoLang`. Hai ban so lieu cung ton tai: ban Go va
// kho-noi-dung/chieu.yaml.
//
// Tu task nay chi con MOT ban: chieu.yaml. Bang Go da BI XOA. Nghia la dot
// bien "actor doc bang Go cu" KHONG THE viet duoc nua — khong con bang nao de
// doc.
//
// # BA LOP, DUNG KHONG LAN NHAU
//
//	nguonNoiDung      — mot NGUON anh chup TIEM DUOC (mot ham lay ra anh chup).
//	                    KhuVuc giu no trong atomic.Pointer, dung khuon voi
//	                    loa/soNo/boNhat/boRoi (xem khuvuc.go).
//	khoChieuAnhChup   — phan cua anh chup ma goi `game` duoc phep dung. Kieu
//	                    that (*noidung.anhChup) KHONG XUAT KHAU nen o day phai
//	                    la mot interface; *noidung.anhChup thoa man no.
//	chieuCuaLangTu /  — hai cua tra cuu, doi kieu noidung.* sang MauChieu/
//	traCapChieuTu       CapChieu (khuon byte -28/7 do goi `game` so huu).
//
// # VI SAO NGUON LA MOT HAM CHU KHONG PHAI *noidung.Kho
//
// Vi AnhChup() tra ve mot kieu KHONG XUAT KHAU. Khong the viet
// `interface{ AnhChup() *noidung.anhChup }` trong goi khac — Go khong cho goi
// ten kieu do. Mot ham `func() khoChieuAnhChup` thi goi nao cung viet duoc, va
// no cho phep test cam mot nguon DEM SO LAN LAY (khochieu_hopdong_test.go) de
// khoa hop dong "mot lan dau nhip" — hop dong do khong con la loi hua trong
// chu thich nua, no co test giet duoc.
//
// # MAC DINH KHI CHUA AI GAN KHO
//
// Rat nhieu test dung MoKhuVuc truc tiep, khong qua main.go, nen KhuVuc phai
// chay duoc khi chua ai gan kho — dung khuon "nil thi im lang" cua loa/soNo.
// NHUNG o day "im lang" khong the la "khong co chieu nao" (nhu the moi test
// chien dau cu se do voi ly do sai), va TUYET DOI khong duoc la "roi ve bang
// Go cung" (bang do da xoa — do chinh la muc tieu cua task).
//
// Vay mac dinh la: NAP CHINH kho-noi-dung NHUNG (noidung.NapNhung) MOT LAN,
// nho lai bang sync.Once. Cung mot file YAML ma may chu that chay, nen test va
// may chu that khong bao gio thay hai bo so khac nhau.
//
// PANIC khi ban nhung hong: file YAML nam TRONG file thuc thi (go:embed) nen
// "nap that bai" khong phai su co moi truong ma la mot ban build hong — im
// lang cho qua se de may chu chay voi 0 chieu. Bo nap da co test rieng
// (internal/noidung/nap_test.go) nen duong nay chi no khi ai do lam hong
// chinh kho noi dung.
package game

import (
	"sync"

	"nroserver/internal/noidung"
)

// khoVatPhamAnhChup la phan MAU VAT PHAM / MAU OPTION cua mot anh chup noi dung
// (KH8 task 4). Tach rieng khoi phan chieu vi hai vung so lieu doc lap nhau:
// mot ham chi tra cuu vat pham thi nhan interface nay, doc chu ky la biet no
// KHONG dung toi chieu.
type khoVatPhamAnhChup interface {
	MauOption() []noidung.MauOption
	MauVatPham() []noidung.MauVatPham
	TimMauVatPham(id int) (noidung.MauVatPham, bool)
}

// khoChieuAnhChup la phan cua mot anh chup noi dung ma goi `game` dung.
// *noidung.anhChup (kieu khong xuat khau, tra ve boi Kho.AnhChup) thoa man no.
//
// KH8 task 4 NHUNG khoVatPhamAnhChup vao day thay vi tao mot nguon thu hai:
// nguonNoiDung (va khuvuc.acNhip) da mang kieu nay roi, nen mot NGUON TIEM VAO
// duy nhat phuc vu ca chieu lan vat pham — hai nguon la hai the he lech nhau.
type khoChieuAnhChup interface {
	ChieuCuaLang(idLang byte) []noidung.Chieu
	TraCapChieu(idLang, templateId, cap byte) (noidung.CapChieu, bool)
	TheHe() uint64

	// Ba cua QUAI (KH8 task 8, khoquai.go). Nam CHUNG mot interface voi chieu
	// — khong de ra interface thu hai — vi mot actor giu DUNG MOT anh chup cho
	// ca nhip, va *noidung.anhChup thoa man ca hai nhom.
	Quai(ma string) (noidung.MauQuai, bool)
	DanhSachQuai() []noidung.MauQuai
	SinhQuaiCuaMap(mapID int32) []noidung.DiemSinhQuai

	// Hai cua BANG ROI (KH8 task 9, roido.go). Cung mot anh chup voi chieu/quai
	// — quai chet la mot su kien, va no phai doc chi so quai VA bang roi tren
	// CUNG mot the he so lieu.
	BangRoiCuaQuai(ma string) []noidung.MucRoi
	BangRoiCuaBoss(ma string) []noidung.MucRoi

	// Phan MAU VAT PHAM / OPTION (KH8 task 4) — cung ly do nhu quai: MOT nguon
	// tiem, MOT anh chup moi nhip. (Hai task lam song song cung mo rong interface
	// nay -> xung dot gop 03/09, dieu phoi giu CA HAI.)
	khoVatPhamAnhChup
}

// nguonNoiDung boc mot ham lay anh chup. La mot STRUCT (khong phai ham tran)
// vi atomic.Pointer can mot kieu co the lay dia chi on dinh.
type nguonNoiDung struct {
	lay func() khoChieuAnhChup
}

// nguonTuKho dung mot nguon doc tu mot *noidung.Kho that. Day la duong ma
// main.go di (qua SoKhuVuc.GanKhoNoiDung).
func nguonTuKho(kho *noidung.Kho) *nguonNoiDung {
	return &nguonNoiDung{lay: func() khoChieuAnhChup { return kho.AnhChup() }}
}

var (
	khoChieuMotLan sync.Once
	khoChieuNhung  *noidung.Kho
)

// khoNhungMotLan tra ve *noidung.Kho da nap kho-noi-dung NHUNG, dung MOT LAN
// (sync.Once) roi dung lai mai mai — xem khoi chu thich dau tep ("MAC DINH KHI
// CHUA AI GAN KHO"). DUNG CHUNG cho moi phan cua goi `game` doc anh chup mac
// dinh (chieu: khoChieuMacDinh duoi day; lang: khoLangMacDinh trong lang.go,
// KH8 task 6) — CHI MOT sync.Once cho ca kho, vi mot anh chup noi dung la MOT
// khoi thong nhat (chieu + lang + ... cung nam trong *noidung.DuLieu).
func khoNhungMotLan() *noidung.Kho {
	khoChieuMotLan.Do(func() {
		du, err := noidung.NapNhung()
		if err != nil {
			panic("game: kho noi dung nhung hong, khong nap duoc: " + err.Error())
		}
		khoChieuNhung = noidung.MoKho(du)
	})
	return khoChieuNhung
}

// khoChieuMacDinh tra ve anh chup cua kho NHUNG — duong dung khi chua ai gan
// kho (xem khoi chu thich dau tep). Nap mot lan, dung lai mai mai.
func khoChieuMacDinh() khoChieuAnhChup {
	return khoNhungMotLan().AnhChup()
}

// ---------------------------------------------------------------------------
// Doi kieu: noidung.Chieu/CapChieu -> MauChieu/CapChieu (khuon byte -28/7)
// ---------------------------------------------------------------------------

func capChieuTuNoiDung(c noidung.CapChieu) CapChieu {
	return CapChieu{
		SkillId:      c.SkillId,
		Diem:         c.Diem,
		NhanLucCan:   c.NhanLucCan,
		ChakraDung:   c.ChakraDung,
		HoiChieuMs:   c.HoiChieuMs,
		Dx:           c.Dx,
		Dy:           c.Dy,
		MaxFight:     c.MaxFight,
		SatThuong:    c.SatThuong,
		Gia:          c.Gia,
		ThemThongTin: c.ThemThongTin,
	}
}

func mauChieuTuNoiDung(c noidung.Chieu) MauChieu {
	cac := make([]CapChieu, 0, len(c.CacCap))
	for _, cap := range c.CacCap {
		cac = append(cac, capChieuTuNoiDung(cap))
	}
	return MauChieu{
		Id:            c.TemplateId,
		Ten:           c.Ten,
		CapToiDa:      c.CapToiDa,
		KieuTonChakra: c.KieuTonChakra,
		Loai:          c.Loai,
		IconId:        c.IconId,
		MoTa:          c.MoTa,
		CacCap:        cac,
	}
}

// ---------------------------------------------------------------------------
// Hai cua tra cuu — nhan ANH CHUP tu ngoai
// ---------------------------------------------------------------------------

// chieuCuaLangTu tra ve 4 MauChieu cua lang idLang, doc TU MOT ANH CHUP da cam
// san. Lang ngoai bang -> nil (fail-closed, giu nguyen hop dong cu).
//
// Ket qua da la BAN SAO RIENG cua nguoi goi hai lop (anh chup tra ban sao sau,
// roi ham nay dung slice moi) — nguoi goi duoc phep giu con tro vao mot phan tu
// (case suKienChonChieu lam dung the) ma khong cham vao anh chup dang cong bo.
func chieuCuaLangTu(ac khoChieuAnhChup, idLang byte) []MauChieu {
	ds := ac.ChieuCuaLang(idLang)
	if ds == nil {
		return nil
	}
	ra := make([]MauChieu, 0, len(ds))
	for _, c := range ds {
		ra = append(ra, mauChieuTuNoiDung(c))
	}
	return ra
}

// traCapChieuTu tra ve MOT CAP cu the cua mot chieu, doc TU MOT ANH CHUP da cam
// san — cua tra cuu DUY NHAT dung chung cho nguoi choi va boss.
//
// Fail-closed y het ban truoc task 3: nil khi lang ngoai bang, khong tim thay
// templateId, hoac cap ngoai [1, len(CacCap)]. cap = 0 la KHONG hop le.
//
// Tra ve CON TRO toi mot ban sao MOI (khong phai vao anh chup): CapChieu chi co
// truong vo huong + mot chuoi nen mot phep gan la sao xong; con tro chi de giu
// nguyen chu ky ma KiemDanhHopLe/TinhSatThuong dang doi (*CapChieu).
func traCapChieuTu(ac khoChieuAnhChup, idLang byte, templateId byte, cap byte) *CapChieu {
	c, co := ac.TraCapChieu(idLang, templateId, cap)
	if !co {
		return nil
	}
	ra := capChieuTuNoiDung(c)
	return &ra
}

// bangChieuMacDinh dung [5][]MauChieu tu anh chup mac dinh — dung cho goi
// -28/7 (dongbo.go, goiCatalogSkill) va cho cac test cu.
func bangChieuMacDinh() [5][]MauChieu {
	return bangChieuTu(khoChieuMacDinh())
}

func bangChieuTu(ac khoChieuAnhChup) [5][]MauChieu {
	var ra [5][]MauChieu
	for lang := 0; lang < noidung.SoLangCoDinh; lang++ {
		ra[lang] = chieuCuaLangTu(ac, byte(lang))
	}
	return ra
}

// ---------------------------------------------------------------------------
// Cua CU giu nguyen chu ky — doc tu anh chup MAC DINH
// ---------------------------------------------------------------------------

// ChieuCuaLang tra ve 4 MauChieu cua lang idLang tu kho NHUNG.
//
// KHONG dung ham nay trong actor (dung k.acNhip qua chieuCuaLangTu): moi lan
// goi lay mot anh chup MOI, tuc hai loi goi trong cung mot nhip co the roi vao
// hai the he khac nhau — dung cai bay ma hop dong "mot lan dau nhip" sinh ra de
// tranh. Day la cua cho duong DANG NHAP / dung catalog, ngoai nhip actor.
func ChieuCuaLang(idLang byte) []MauChieu {
	return chieuCuaLangTu(khoChieuMacDinh(), idLang)
}

// TraCapChieu tra ve MOT CAP cua mot chieu tu kho NHUNG. Cung canh bao voi
// ChieuCuaLang o tren: trong actor hay dung traCapChieuTu(k.acNhip, ...).
func TraCapChieu(idLang byte, templateId byte, cap byte) *CapChieu {
	return traCapChieuTu(khoChieuMacDinh(), idLang, templateId, cap)
}
