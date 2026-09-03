// boss.go — KIEU DU LIEU + BO KIEM + CUA TRA CUU cua BOSS trong kho noi dung.
//
// (KH8 task 7. Thiet ke: docs/superpowers/specs/2026-09-03-kh8-kho-noi-dung-design.md
// hang "7" bang Task §7, dot bien M7/M8/M10/M11 §6.)
//
// # CAI GI DOI O TASK NAY
//
// Truoc task nay, hai con boss mau nam trong MOT BANG GO CUNG
// (game/noidungboss.go, bien `danhSachNoiDungBoss`). Tu task nay bang do DA XOA
// va nguon duy nhat la kho-noi-dung/boss.yaml. Nghia la dot bien "ai do sua
// bang Go cu" khong con viet duoc — khong con bang nao de sua.
//
// # TAT CA NAM TRONG MOT TEP
//
// Kieu YAML, bo kiem, VA hai cua tra cuu tren anh chup (Boss/DanhSachBoss) deu
// o day, khong rai sang noidung.go/nap.go/kho.go. Ly do rat thuc te: ngay
// 03/09/2026 co BA task KH8 chay song song cham vao dung ba tep do; moi dong
// them vao chung la mot xung dot gop. Ranh gioi "mot vung du lieu = mot tep" giu
// cho task nay chi de lai DUNG BA dong o ba tep chung (mot truong DuLieu.Boss,
// mot lan doc file trong Nap, mot lan goi trong kiemDayDu).
//
// # VI SAO BO KIEM CHIA LAM BA CHO
//
//	kiemBoss              — kiem duoc bang MOT MINH boss.yaml (so, thu tu, ma).
//	                        Goi trong Nap.
//	kiemBossThamChieuChieu — can CA chieu.yaml (bo chieu cua boss phai la chieu
//	                        CO THAT va phai la chieu loai DANH). Goi trong
//	                        kiemDayDu, sau khi ca hai file da nap.
//	KiemIdPartBoss        — can catalog ANH DA KY (asset-lock.json) ma goi nay
//	                        KHONG duoc import (huong import: game -> noidung).
//	                        Ham THUAN, nhan danh sach id lam tham so; nguoi TIEM
//	                        du lieu (game/noidungboss.go) goi voi danh sach that.
//	                        Cung khuon voi KiemIdAnhQuai (kiemtra.go, task 8).
//
// Bo kiem VA CHAM (boss co dung trong vach da khong) KHONG o day: luoi va cham
// la artifact cua goi `game`. Day chi kiem toa do co nam trong mien int16 khong
// am — phan con lai la viec cua TestNoiDungBoss_ViTriSinhDungTrenMatDat.
package noidung

import "fmt"

// ---------------------------------------------------------------------------
// Kieu du lieu YAML — boss
// ---------------------------------------------------------------------------

// ViTriSinhBoss la MOT cho dat chan hop le cua boss: map nao, toa do nao.
//
// VI SAO KHONG NAM TRONG PhaBoss: PhaBoss.MapVao chi noi boss hien o MAP nao,
// khong noi hien o CHO nao trong map. Ma cho dat chan lai la thu duy nhat trong
// ca bo du lieu nay kiem duoc bang mot nguon KHACH QUAN (luoi va cham cua map),
// nen no dang mot kieu rieng de bo test ben `game` bat duoc.
//
// Toa do la "chan nhan vat dung o dau": o chua toa do phai DI DUOC, va o ngay
// duoi phai DAC (co mat dat). Sinh giua troi thi boss roi/lo lung tuy client.
type ViTriSinhBoss struct {
	Map int32 `yaml:"map"`
	X   int16 `yaml:"x"`
	Y   int16 `yaml:"y"`
}

// MocNoiGianBoss la mot moc "noi gian theo %mau" cua mot pha: khi HP boss tut
// xuong duoi PhanTramHP thi doi he so dame/def/nhip va phat mot cau chat.
//
// HeSo* de float64 vi day la DU LIEU CAN BANG (chu du an chinh tay, vd 1.5 /
// 0.75), khong phai duong tinh sat thuong loi — duong loi van tinh bang so
// nguyen (game/satthuong.go) de tranh sai so gan bien tran int32.
type MocNoiGianBoss struct {
	PhanTramHP int     `yaml:"phan_tram_hp"` // nguong %HP kich hoat (vd 50 = con 50% mau)
	HeSoDame   float64 `yaml:"he_so_dame"`   // nhan vao sat thuong boss GAY RA
	HeSoDef    float64 `yaml:"he_so_def"`    // nhan vao sat thuong boss NHAN (nho hon = trau hon)
	HeSoNhip   float64 `yaml:"he_so_nhip"`   // nhan vao nhip danh (nho hon 1 = danh nhanh hon); 0 = khong doi
	Chat       string  `yaml:"chat"`         // cau boss noi khi cham moc nay (goi 44)
}

// OChieuBoss la MOT o trong bo chieu cua mot pha: {template_id, cap}.
//
// VI SAO LA STRUCT CO TEN TRUONG chu khong phai mot cap so `[1, 2]`: trong Go
// ban cu no la [2]int32 va khong ai doc duoc "so nao la template, so nao la
// cap" neu khong mo code ra. Trong YAML thi mot cap so tran se cho phep go
// nguoc thu tu ma khong bo kiem nao thay — ca hai so deu la so nguyen nho.
type OChieuBoss struct {
	TemplateId int32 `yaml:"template_id"`
	Cap        int32 `yaml:"cap"`
}

// PhaBoss la MOT PHA cua mot boss. Mot boss cu the = nhieu pha: giet pha nay ->
// hoi sinh pha manh hon, GIU NGUYEN id.
//
// Kieu o day co y KHONG trung khop tung kieu nen voi game.PhaBoss (Cap la int
// chu khong phai byte, TypePk la int chu khong phai byte): mot so vuot tran
// trong YAML phai la RA O BO KIEM voi mot cau tieng nguoi, chu khong duoc bi
// trinh doc YAML cat am tham roi thanh mot con so hop le khac.
type PhaBoss struct {
	Ten    string   `yaml:"ten"`     // ten hien tren dau boss (goi -5 truong 9)
	Lang   int      `yaml:"lang"`    // lang Naruto — thay cho "gender" Dragon Ball ban goc
	Cap    int      `yaml:"cap"`     // cap hien thi (goi -5 truong 3), hop le 1..255
	Outfit [6]int16 `yaml:"outfit"`  // head, body, leg, flagBag, aura, effFront
	Dame   int64    `yaml:"dame"`    // sat thuong goc cua pha (dua vao TinhSatThuong)
	TypePk int      `yaml:"type_pk"` // truong 5 goi -5; PHAI khac 0, xem kiemBoss

	// NhipDanhMs la khoang cach toi thieu giua hai don cua pha nay (mili giay).
	NhipDanhMs int32 `yaml:"nhip_danh_ms"`

	HpChon  []int64      `yaml:"hp_chon"`  // luc sinh chon NGAU NHIEN mot gia tri trong day
	MapVao  []int32      `yaml:"map_vao"`  // luc hien chon NGAU NHIEN mot khu vuc trong day
	BoChieu []OChieuBoss `yaml:"bo_chieu"` // bo chieu, tra cuu trong catalog cua Lang

	GiayNghi   int32            `yaml:"giay_nghi"`    // so giay nam yen truoc khi hien
	MocNoiGian []MocNoiGianBoss `yaml:"moc_noi_gian"` // cac moc noi gian theo %HP
	LaBossLon  bool             `yaml:"la_boss_lon"`  // bat goi 93 (chay chu toan server)
	DiemThuong int64            `yaml:"diem_thuong"`  // tong diem chia khi GIET PHA NAY

	ThoaiHien []string `yaml:"thoai_hien"` // chat luc xuat hien (trang thai CHAO)
	ThoaiDanh []string `yaml:"thoai_danh"` // chat luc dang danh (trang thai HOAT_DONG)
	ThoaiChet []string `yaml:"thoai_chet"` // chat luc chet (trang thai TROI)
}

// Boss la MOT con boss trong kho noi dung: cac pha cua no + cac cho no duoc phep
// dat chan.
//
// CAN CUOC la `Ma` (chuoi), KHONG phai so — cung ly do voi MauQuai (noidung.go):
// ma boss KHONG di xuong client va KHONG nam trong CSDL nguoi choi, no chi la
// khoa tra cuu cua bo dieu phoi lich. Them/bo/sap xep lai boss la an toan.
type Boss struct {
	Ma    string          `yaml:"ma"`
	ViTri []ViTriSinhBoss `yaml:"vi_tri"`
	Pha   []PhaBoss       `yaml:"pha"`
}

type fileBoss struct {
	Boss []Boss `yaml:"boss"`
}

// TranCapBoss / TranTypePkBoss: Cap va TypePk la MOT BYTE trong goi -5
// (game/xuathien.go). Day la SUC DOC CUA CLIENT, khong phai lua chon thiet ke.
const (
	TranCapBoss    = 255
	TranTypePkBoss = 255
)

// Bo kiem nhip danh dung NhipActorMs (khai o kiemtra.go, hang KHONG KIEU nen so
// sanh thang duoc voi int32). KHONG khai lai con so 100 o day: boss cap nhat NGAY
// TRONG nhip actor, nen mot nhip danh 1450ms tren thuc te van la 1500ms — va con
// so do phai la CUNG MOT con so voi cai ma chieu thuc bi kiem (dot bien M7).

// TranHpBoss la tran cua Boss.Hp/HpToiDa (int32, game/boss.go KepHp). Mot gia
// tri hp_chon vuot no se bi cat xuong tran ma KHONG bao gi — con boss "3 ty mau"
// am tham thanh con boss 2,1 ty mau.
const TranHpBoss int64 = 2147483647

// ---------------------------------------------------------------------------
// Bo kiem — chi can boss.yaml
// ---------------------------------------------------------------------------

// kiemBoss la bo kiem doc lap cua boss.yaml. Goi trong Nap.
//
// Moi dieu kien duoi day deu la "sai thi con boss HONG CAM": no khong lam
// chuong trinh sap, khong lam test nao khac do, no chi lam boss dung im, vo
// hinh, danh khong trung, hoac sinh trong long dat — loai loi chi lo ra khi co
// nguoi that vao game nhin thay, tuc muon nhat co the.
func kiemBoss(ds []Boss) error {
	if len(ds) == 0 {
		return fmt.Errorf("khong co con boss nao — he boss se chay voi lich RONG va khong ai biet")
	}
	daThayMa := map[string]bool{}
	for i, b := range ds {
		if b.Ma == "" {
			return fmt.Errorf("boss thu %d: `ma` rong — `ma` la can cuoc, bo dieu phoi lich goi boss theo chuoi nay", i)
		}
		if daThayMa[b.Ma] {
			return fmt.Errorf("ma boss %q khai hai lan — tra cuu tra ra con dau tien, con thu hai KHONG BAO GIO "+
				"duoc sinh: mot con boss bien mat khoi game ma khong mot dong loi", b.Ma)
		}
		daThayMa[b.Ma] = true

		if len(b.Pha) == 0 {
			return fmt.Errorf("boss %q: khong co pha nao — khong biet sinh con gi", b.Ma)
		}
		if len(b.ViTri) == 0 {
			return fmt.Errorf("boss %q: khong co `vi_tri` nao — khong biet dat chan o dau", b.Ma)
		}
		for j, v := range b.ViTri {
			// CHI kiem mien so o day. "O do co dac khong, duoi do co mat dat
			// khong" phai kiem bang LUOI VA CHAM that — artifact cua goi `game`
			// (TestNoiDungBoss_ViTriSinhDungTrenMatDat, dot bien M11).
			if v.Map < 1 {
				return fmt.Errorf("boss %q vi tri #%d: map = %d, id ban do bat dau tu 1", b.Ma, j, v.Map)
			}
			if v.X < 1 || v.Y < 1 {
				return fmt.Errorf("boss %q vi tri #%d: toa do (%d,%d) — cot/hang 0 cua luoi la tuong bien, "+
					"boss se nam trong da", b.Ma, j, v.X, v.Y)
			}
		}

		// Hai bang du lieu (MapVao cua pha va ViTri cua boss) phai noi cung mot
		// chuyen. Lech nhau nghia la hoac boss co cho dat chan o map no khong
		// bao gio hien toi, hoac no hien vao map khong co cho dat chan nao — ca
		// hai deu la du lieu chet, khong loi.
		coChoDat := map[int32]bool{}
		for _, v := range b.ViTri {
			coChoDat[v.Map] = true
		}
		trongMapVao := map[int32]bool{}

		for j := range b.Pha {
			if err := kiemMotPhaBoss(b.Ma, j, &b.Pha[j]); err != nil {
				return err
			}
			for _, id := range b.Pha[j].MapVao {
				trongMapVao[id] = true
			}
		}
		for id := range trongMapVao {
			if !coChoDat[id] {
				return fmt.Errorf("boss %q: hien vao map %d nhung khong co `vi_tri` nao o do", b.Ma, id)
			}
		}
		for id := range coChoDat {
			if !trongMapVao[id] {
				return fmt.Errorf("boss %q: co `vi_tri` o map %d nhung khong pha nao hien vao map do", b.Ma, id)
			}
		}
	}
	return nil
}

func kiemMotPhaBoss(ma string, chiSo int, p *PhaBoss) error {
	nhan := fmt.Sprintf("boss %q pha %d", ma, chiSo+1)

	if p.Ten == "" {
		return fmt.Errorf("%s: `ten` rong — client ve bong ten trong tren dau boss", nhan)
	}
	if p.Lang < 0 || p.Lang >= SoLangCoDinh {
		return fmt.Errorf("%s: lang = %d, hop le 0..%d — bo chieu tra cuu trong catalog cua lang nay",
			nhan, p.Lang, SoLangCoDinh-1)
	}
	if p.Cap < 1 || p.Cap > TranCapBoss {
		return fmt.Errorf("%s: cap = %d, hop le 1..%d — goi -5 ghi cap bang MOT BYTE, 0 lam client hien "+
			"\"cap 0\" vo nghia", nhan, p.Cap, TranCapBoss)
	}
	if p.TypePk <= 0 || p.TypePk > TranTypePkBoss {
		return fmt.Errorf("%s: type_pk = %d, hop le 1..%d — client (Char.cs paintHp) THOAT SOM khong ve "+
			"thanh mau khi typePk == 0, tuc nguoi choi danh boss ma khong thay boss con bao nhieu mau",
			nhan, p.TypePk, TranTypePkBoss)
	}
	if p.Dame <= 0 {
		return fmt.Errorf("%s: dame = %d — boss khong gay duoc sat thuong nao (sat thuong = dame × %%cap chieu)",
			nhan, p.Dame)
	}

	// M7 — nhip danh phai la boi so cua mot nhip actor VA khong nho hon mot nhip.
	if p.NhipDanhMs < NhipActorMs {
		return fmt.Errorf("%s: nhip_danh_ms = %d — nho hon mot nhip actor (%dms), boss danh MOI TICK "+
			"(10 don mot giay)", nhan, p.NhipDanhMs, NhipActorMs)
	}
	if p.NhipDanhMs%NhipActorMs != 0 {
		return fmt.Errorf("%s: nhip_danh_ms = %d khong phai boi so cua %dms — boss cap nhat trong nhip "+
			"%dms cua khu vuc, phan le khong bao gio xay ra that", nhan, p.NhipDanhMs, NhipActorMs, NhipActorMs)
	}

	// M8 — khoang lang giua hai pha phai > 0. GiayNghi = 0 nghia la pha sau no ra
	// NGAY trong cung mot khoanh khac pha truoc chet, nguoi choi khong kip doc gi;
	// voi pha dau thi 0 con nghia la boss sinh lai lien tuc khong nghi.
	if p.GiayNghi <= 0 {
		return fmt.Errorf("%s: giay_nghi = %d, phai > 0 — day la so giay nam yen truoc khi hien; 0 nghia la "+
			"pha nay no ra ngay trong cung khoanh khac pha truoc chet, nguoi choi khong kip doc gi",
			nhan, p.GiayNghi)
	}

	if len(p.MapVao) == 0 {
		return fmt.Errorf("%s: `map_vao` RONG — boss khong biet hien o dau", nhan)
	}
	for _, id := range p.MapVao {
		if id < 1 {
			return fmt.Errorf("%s: map_vao chua map %d — id ban do bat dau tu 1", nhan, id)
		}
	}

	// hp_chon: moi gia tri > 0 va VUA int32. Vuot tran thi KepHp cat xuong ma
	// khong bao gi, chenh lech khong bao gio lo ra o dau.
	if len(p.HpChon) == 0 {
		return fmt.Errorf("%s: `hp_chon` RONG — luc sinh khong biet cho boss bao nhieu mau", nhan)
	}
	for i, hp := range p.HpChon {
		if hp <= 0 {
			return fmt.Errorf("%s: hp_chon[%d] = %d — boss sinh ra la chet ngay", nhan, i, hp)
		}
		if hp > TranHpBoss {
			return fmt.Errorf("%s: hp_chon[%d] = %d vuot tran int32 (%d) cua Boss.Hp — KepHp se cat xuong "+
				"tran ma KHONG bao gi", nhan, i, hp, TranHpBoss)
		}
	}

	if len(p.BoChieu) == 0 {
		return fmt.Errorf("%s: `bo_chieu` RONG — boss khong co gi de danh", nhan)
	}
	for i, o := range p.BoChieu {
		if o.TemplateId < 1 || o.TemplateId > 255 {
			return fmt.Errorf("%s: bo_chieu[%d].template_id = %d — hop le 1..255, id bat dau tu 1",
				nhan, i, o.TemplateId)
		}
		if o.Cap < 1 || o.Cap > 255 {
			return fmt.Errorf("%s: bo_chieu[%d].cap = %d — hop le 1..255, cap tinh tu 1", nhan, i, o.Cap)
		}
	}

	if err := kiemMocNoiGianBoss(nhan, p); err != nil {
		return err
	}
	return kiemThoaiBoss(nhan, p)
}

// kiemMocNoiGianBoss — MocNoiGianDaQua (game/boss.go) la con tro chay MOT CHIEU:
// may trang thai duyet mang theo THU TU KHAI BAO khi mau tut xuong. Vi vay mang
// PHAI sap theo %HP giam dan va khong duoc trung nguong:
//
//   - sap sai thu tu -> moc sau khong bao gio toi luot, du lieu chet;
//   - trung nguong   -> hai moc cung cham mot luc, mot cau chat bi nuot.
//
// "Noi gian" ma boss danh nhe di hoac mong hon la du lieu NGUOC NGHIA: HeSoDame
// nhan vao sat thuong boss GAY RA (lon hon = manh hon), HeSoDef nhan vao sat
// thuong boss NHAN (nho hon = trau hon) — hai chieu nguoc nhau, rat de dien nham.
//
// Va nhip sau khi da nhan HET he so (chung nhan CONG DON) van phai >= mot nhip
// actor: noi gian chong nhau co the keo nhip xuong duoi mot tick, luc do boss
// danh moi tick — 10 don mot giay (dot bien M7 o dang cong don).
func kiemMocNoiGianBoss(nhan string, p *PhaBoss) error {
	truoc := 100
	dameTruoc, defTruoc := 1.0, 1.0
	nhip := float64(p.NhipDanhMs)

	for i, m := range p.MocNoiGian {
		if m.PhanTramHP <= 0 || m.PhanTramHP >= 100 {
			return fmt.Errorf("%s moc_noi_gian[%d]: phan_tram_hp = %d — phai trong (0,100). 100 la boss noi "+
				"gian ngay khi vua hien (chua ai danh), 0 thi khong bao gio toi vi mau ve 0 la boss chet",
				nhan, i, m.PhanTramHP)
		}
		if m.PhanTramHP >= truoc {
			return fmt.Errorf("%s moc_noi_gian[%d]: phan_tram_hp = %d khong nho hon moc truoc (%d) — mang phai "+
				"GIAM DAN va khong trung, vi con tro moc chay MOT CHIEU", nhan, i, m.PhanTramHP, truoc)
		}
		truoc = m.PhanTramHP

		if m.HeSoDame < dameTruoc {
			return fmt.Errorf("%s moc_noi_gian[%d]: he_so_dame = %g yeu hon muc truoc (%g) — noi gian ma danh nhe di",
				nhan, i, m.HeSoDame, dameTruoc)
		}
		if m.HeSoDef <= 0 || m.HeSoDef > defTruoc {
			return fmt.Errorf("%s moc_noi_gian[%d]: he_so_def = %g phai trong (0, %g] — noi gian ma mong hon truoc "+
				"(he_so_def nhan vao sat thuong boss NHAN: nho hon = trau hon)", nhan, i, m.HeSoDef, defTruoc)
		}
		dameTruoc, defTruoc = m.HeSoDame, m.HeSoDef

		if m.Chat == "" {
			return fmt.Errorf("%s moc_noi_gian[%d]: `chat` rong — noi gian ma khong noi gi thi nguoi choi "+
				"khong biet co gi doi", nhan, i)
		}
		if soByteUTF(m.Chat) > TranThanGoiChat {
			return fmt.Errorf("%s moc_noi_gian[%d]: cau chat dai %d byte UTF, vuot %d ma goi 44 mang duoc",
				nhan, i, soByteUTF(m.Chat), TranThanGoiChat)
		}

		if m.HeSoNhip < 0 {
			return fmt.Errorf("%s moc_noi_gian[%d]: he_so_nhip = %g am", nhan, i, m.HeSoNhip)
		}
		if m.HeSoNhip == 0 {
			continue // 0 = KHONG doi nhip, dung nghia trong game/boss.go
		}
		nhip *= m.HeSoNhip
		if nhip < float64(NhipActorMs) {
			return fmt.Errorf("%s: sau %d moc noi gian, nhip danh con %.0fms — duoi mot tick %dms, boss danh "+
				"10 don mot giay", nhan, i+1, nhip, NhipActorMs)
		}
	}
	return nil
}

// TranThanGoiChat: goi 44 mang cau chat bang writeUTF (than toi da 65535 byte).
// Dai hon la cau bi cat giua chung.
const TranThanGoiChat = 65535

// kiemThoaiBoss — ba mang thoai deu duoc may trang thai dung o mot buoc khac
// nhau (CHAO / HOAT_DONG / TROI). Mang rong = buoc do IM LANG; cau rong = mot
// bong bong thoai TRONG tren dau boss.
func kiemThoaiBoss(nhan string, p *PhaBoss) error {
	for _, bo := range []struct {
		ten  string
		cauc []string
	}{
		{"thoai_hien", p.ThoaiHien},
		{"thoai_danh", p.ThoaiDanh},
		{"thoai_chet", p.ThoaiChet},
	} {
		if len(bo.cauc) == 0 {
			return fmt.Errorf("%s: `%s` RONG — buoc tuong ung cua may trang thai se im lang", nhan, bo.ten)
		}
		for i, cau := range bo.cauc {
			if cau == "" {
				return fmt.Errorf("%s: %s[%d] rong — bong bong thoai trong tren dau boss", nhan, bo.ten, i)
			}
			if soByteUTF(cau) > TranThanGoiChat {
				return fmt.Errorf("%s: %s[%d] dai %d byte UTF, vuot %d ma goi 44 mang duoc — se bi cat",
					nhan, bo.ten, i, soByteUTF(cau), TranThanGoiChat)
			}
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Bo kiem — can CA chieu.yaml (goi tu kiemDayDu)
// ---------------------------------------------------------------------------

// kiemBossThamChieuChieu doi chieu MOI o trong `bo_chieu` voi catalog chieu THAT
// cua chinh lang cua pha (dot bien M10).
//
// HAI dieu, thieu mot la hong:
//
//  1. Chieu phai CO THAT o cap do — dien mot templateId hay mot cap khong ton
//     tai la boss dung im khong danh duoc, va KHONG co thong bao loi nao.
//  2. Chieu phai la LOAI DANH. Catalog co 3 loai: danh (1), buff/hoi mau (2),
//     tu than (3). Loai 2 khong sinh goi tan cong nao; loai 3 co tam danh = 0
//     nen moi don deu vuot tam -> LUON truot. Ca hai cho ra con boss "vung tay
//     vao khong khi": trong nhu dang danh nhau ma khong bao gio gay mot diem sat
//     thuong nao.
func kiemBossThamChieuChieu(du *DuLieu) error {
	for _, b := range du.Boss {
		for j := range b.Pha {
			p := &b.Pha[j]
			if p.Lang < 0 || p.Lang >= SoLangCoDinh {
				continue // da bao o kiemBoss
			}
			danhMuc := du.ChieuTheoLang[p.Lang]
			for i, o := range p.BoChieu {
				var c *Chieu
				for k := range danhMuc {
					if int32(danhMuc[k].TemplateId) == o.TemplateId {
						c = &danhMuc[k]
						break
					}
				}
				if c == nil {
					return fmt.Errorf("boss %q pha %d bo_chieu[%d]: template_id %d KHONG CO trong catalog "+
						"chieu cua lang %d (%s) — boss se dung im, khong mot dong loi",
						b.Ma, j+1, i, o.TemplateId, p.Lang, TenFileChieu)
				}
				if int(o.Cap) > len(c.CacCap) {
					return fmt.Errorf("boss %q pha %d bo_chieu[%d]: chieu %q chi co %d cap nhung boss cam cap %d",
						b.Ma, j+1, i, c.Ten, len(c.CacCap), o.Cap)
				}
				if c.Loai != LoaiChieuDanh {
					return fmt.Errorf("boss %q pha %d bo_chieu[%d]: chieu %q (template_id %d) co loai = %d, "+
						"boss CHI duoc dung chieu DANH (loai %d) — chieu buff khong sinh goi tan cong nao va "+
						"chieu tu than co tam 0 nen LUON truot: boss se vung tay vao khong khi",
						b.Ma, j+1, i, c.Ten, o.TemplateId, c.Loai, LoaiChieuDanh)
				}
			}
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Bo kiem — can catalog ANH DA KY (goi tu game/noidungboss.go)
// ---------------------------------------------------------------------------

// KiemIdPartBoss doi chieu ba o dau/than/chan cua `outfit` voi danh sach id part
// HOP LE, va bat ba o con lai (flagBag, aura, effFront) phai la 0.
//
// VI SAO HAM NAY XUAT KHAU va nhan danh sach lam THAM SO — y het KiemIdAnhQuai
// (kiemtra.go): nguon su that cua id part la `internal/game/assets/catalog.json`,
// bi GHIM SHA trong `asset-lock.json` (`closure.part_ids`), va no do goi `game`
// doc. Huong import da chot la game -> noidung, TUYET DOI khong nguoc lai, nen
// goi nay khong the tu di lay danh sach do.
//
// Ham THUAN — test duoc bang mot danh sach id gia, khong can dung toi catalog that.
func KiemIdPartBoss(ds []Boss, idHopLe []int16) error {
	co := make(map[int16]bool, len(idHopLe))
	for _, id := range idHopLe {
		co[id] = true
	}
	tenO := [6]string{"head", "body", "leg", "flagBag", "aura", "effFront"}
	for _, b := range ds {
		for j := range b.Pha {
			p := &b.Pha[j]
			for i := 0; i < 3; i++ {
				if !co[p.Outfit[i]] {
					return fmt.Errorf("boss %q pha %d: outfit[%d] (%s) = %d KHONG co trong catalog anh da ky "+
						"(asset-lock.json, closure.part_ids) — client khong tim thay anh thi KHONG VE GI CA, "+
						"con boss VO HINH, va may chu khong co mot dong log nao",
						b.Ma, j+1, i, tenO[i], p.Outfit[i])
				}
			}
			for i := 3; i < 6; i++ {
				if p.Outfit[i] != 0 {
					return fmt.Errorf("boss %q pha %d: outfit[%d] (%s) = %d — chua co anh Naruto nao cho o nay, "+
						"phai de 0 (va hai trong ba o do chi MOT BYTE tren day nen so > 255 bi cat am tham)",
						b.Ma, j+1, i, tenO[i], p.Outfit[i])
				}
			}
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Cua tra cuu tren anh chup — TAT CA tra BAN SAO SAU
// ---------------------------------------------------------------------------

// Boss tra ve ban sao SAU cua mot con boss theo ma, hoac false neu khong co.
//
// Fail-closed (khong tra mot Boss rong): mot Boss rong se di tiep vao may trang
// thai va sinh ra con boss khong pha, khong mau, khong ten — hong CAM.
func (a *anhChup) Boss(ma string) (Boss, bool) {
	for _, b := range a.duLieu.Boss {
		if b.Ma == ma {
			return khoSaoMotBoss(b), true
		}
	}
	return Boss{}, false
}

// DanhSachBoss tra ve ban sao SAU toan bo kho boss. Danh cho luc dung LICH boss
// (mot lan luc khoi dong / luc nap nong), KHONG danh cho vong lap moi nhip.
func (a *anhChup) DanhSachBoss() []Boss {
	if a.duLieu.Boss == nil {
		return nil
	}
	ds := make([]Boss, len(a.duLieu.Boss))
	for i, b := range a.duLieu.Boss {
		ds[i] = khoSaoMotBoss(b)
	}
	return ds
}

// khoSaoMotBoss sao BA LOP: boss -> pha -> cac lat cat trong pha.
//
// VI SAO PHAI SAO SAU: sao mot struct trong Go chi sao CON TRO cua lat cat, nen
// nguoi goi van ghi de duoc du lieu goc qua ban "sao chep" — con boss ma dong
// tren dinh cam. Du an nay da dinh dung bay do hai lan (game/banDo.go TimBanDo,
// game/noidungboss.go saoNoiDungBoss). Voi boss thi hau qua nang hon: du lieu bi
// sua chi lo ra o LAN SINH SAU, o mot khu vuc khac, rat lau sau khi cho sua tay
// xay ra.
//
// KHONG can sao [6]int16 va cac truong so: mang co do dai co dinh trong Go la
// GIA TRI, copy struct da chep tron.
func khoSaoMotBoss(goc Boss) Boss {
	ra := goc
	if goc.ViTri != nil {
		ra.ViTri = make([]ViTriSinhBoss, len(goc.ViTri))
		copy(ra.ViTri, goc.ViTri)
	}
	if goc.Pha == nil {
		return ra
	}
	ra.Pha = make([]PhaBoss, len(goc.Pha))
	for i, pha := range goc.Pha {
		p := pha
		if pha.HpChon != nil {
			p.HpChon = make([]int64, len(pha.HpChon))
			copy(p.HpChon, pha.HpChon)
		}
		if pha.MapVao != nil {
			p.MapVao = make([]int32, len(pha.MapVao))
			copy(p.MapVao, pha.MapVao)
		}
		if pha.BoChieu != nil {
			p.BoChieu = make([]OChieuBoss, len(pha.BoChieu))
			copy(p.BoChieu, pha.BoChieu)
		}
		if pha.MocNoiGian != nil {
			p.MocNoiGian = make([]MocNoiGianBoss, len(pha.MocNoiGian))
			copy(p.MocNoiGian, pha.MocNoiGian)
		}
		p.ThoaiHien = khoSaoChuoi(pha.ThoaiHien)
		p.ThoaiDanh = khoSaoChuoi(pha.ThoaiDanh)
		p.ThoaiChet = khoSaoChuoi(pha.ThoaiChet)
		ra.Pha[i] = p
	}
	return ra
}

func khoSaoChuoi(ds []string) []string {
	if ds == nil {
		return nil
	}
	ra := make([]string, len(ds))
	copy(ra, ds)
	return ra
}
