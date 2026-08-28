// noidungboss_test.go — CHOT KIEM TINH DUNG DAN CUA DU LIEU boss mau
// (KH6b Task 15, nua mot: du lieu thuan Go).
//
// Day KHONG phai test logic — noidungboss.go khong co mot cau lenh if nao.
// Day la LUOI CHAN cho mot loai loi rat rieng cua du lieu noi dung: no khong
// lam chuong trinh sap, khong lam test nao khac do, no chi lam con boss HONG
// CAM trong game — dung im, vo hinh, danh khong trung, hoac sinh trong long
// dat. Kieu loi do chi lo ra khi co nguoi that vao game nhin thay, tuc la
// muon nhat co the.
//
// Vi vay moi test o day deu co dang "moi con so trong du lieu phai TRUY DUOC
// ve mot nguon that trong code": templateId phai co trong catalog chieu
// (chieuthuc.go), mapId phai co trong danh sach ban do (banDo.go), toa do sinh
// phai di duoc theo chinh luoi va cham cua map (catalog.go), cap phai qua
// CapHopLe (boss.go). Khong test nao o day tu bia ra mot con so ky vong roi
// so sanh — chung deu doi chieu voi mot bang tra CO THAT o cho khac.
package game

import (
	"strings"
	"testing"
)

// oTileVaCham la canh mot o cua luoi va cham map 1 (24x24 — khoa boi
// TestVaChamMap1DungKichThuocLuoi trong vacham_map1_test.go). Khai lai o day
// thay vi doc collisionLangLa.tileHeight LA CO Y: neu ai do doi kich thuoc o
// trong artifact, test kia do TRUOC, con test nay van kiem dung y "o ngay
// duoi chan" theo con so da chot.
const oTileVaCham int16 = 24

// mapCoLuoiVaCham: cac map co artifact va cham THAT de kiem duoc toa do sinh.
// Hien chi map 1 (Lang La) — map 2-5 chua co file collision nao, dichuyen.go
// phai dung bien du phong cho chung. Boss dat chan o map chua co luoi thi
// KHONG AI kiem duoc no co dung trong vach da khong, nen o day fail-closed:
// them mot vi tri sinh o map khac se lam test do va bat nguoi them phai mang
// theo luoi va cham cua map do.
var mapCoLuoiVaCham = map[int32]*Collision{1: collisionLangLa}

// phanDaCoAnhThat la tap id "part" (dau/than/chan) ma client CHAC CHAN co anh.
//
// NGUON (khong phai suy doan): server/internal/game/assets/asset-lock.json,
// closure.part_ids = [30, 57, 58, 2189, 2190, 2191]. Trong do:
//   - 30/57/58  la bo do hoc vien (PHASE_ONE_REPLACEMENTS trong assetpack/part.py),
//     3 mieng dau/than/chan da duoc thay bang anh Naruto that;
//   - 2189/2190/2191 la bo cua NPC Iruka (assets/catalog.json npcs[0], va
//     dongbo_tichhop_test.go:297 ghim dung ba so nay xuong day).
//
// Dien mot id NGOAI tap nay vao Outfit la con boss VO HINH: client tim anh
// khong thay thi khong ve gi ca, va khong co mot dong log nao ben may chu.
var phanDaCoAnhThat = map[int16]bool{30: true, 57: true, 58: true, 2189: true, 2190: true, 2191: true}

// duyetPha goi ham than cho MOI pha cua MOI boss mau, kem mot nhan de loi bao
// ra chi thang duoc con nao pha nao. Gom mot cho de moi test ben duoi chi lo
// dung mot dieu kien.
func duyetPha(t *testing.T, than func(t *testing.T, nhan string, mau NoiDungBoss, chiSoPha int, pha *PhaBoss)) {
	t.Helper()
	ds := DanhSachNoiDungBoss()
	if len(ds) == 0 {
		t.Fatal("danh sach boss mau RONG — task 15 khong co gi de nghiem thu")
	}
	for _, mau := range ds {
		for i := range mau.Pha {
			nhan := mau.Ma + " pha " + itoaNhanh(i+1)
			than(t, nhan, mau, i, &mau.Pha[i])
		}
	}
}

// itoaNhanh: strconv chi de in mot chu so trong nhan test thi nang tay hon can.
func itoaNhanh(v int) string {
	if v < 10 {
		return string(rune('0' + v))
	}
	return string(rune('0'+v/10)) + string(rune('0'+v%10))
}

// TestNoiDungBoss_DuBaDuongChayThuNghiem: bo boss mau phai du suc chay HET ba
// duong dac trung cua he boss, neu khong thi nghiem thu bang client that
// (nua hai cua task 15) khong chung minh duoc gi.
//
//	nhieu pha    -> duong bien hinh (RA_DI con pha -> ve HIEN pha+1, giu id)
//	MocNoiGian   -> duong noi gian theo %mau
//	LaBossLon    -> duong thong bao chay chu toan server (goi 93)
func TestNoiDungBoss_DuBaDuongChayThuNghiem(t *testing.T) {
	var coNhieuPha, coNoiGian, coBossLon bool
	for _, mau := range DanhSachNoiDungBoss() {
		if len(mau.Pha) >= 2 {
			coNhieuPha = true
		}
		for _, pha := range mau.Pha {
			if len(pha.MocNoiGian) > 0 {
				coNoiGian = true
			}
			if pha.LaBossLon {
				coBossLon = true
			}
		}
	}
	if !coNhieuPha {
		t.Error("khong con boss mau nao co tu 2 pha tro len — khong chay duoc duong bien hinh")
	}
	if !coNoiGian {
		t.Error("khong pha nao co MocNoiGian — khong chay duoc duong noi gian theo %mau")
	}
	if !coBossLon {
		t.Error("khong pha nao co LaBossLon — khong chay duoc duong thong bao chay chu (goi 93)")
	}
}

// TestNoiDungBoss_MoiChieuTonTaiThatTrongCatalog: tung o {templateId, cap}
// cua BoChieu phai tra ra mot CapChieu THAT qua TraCapChieu — cung cua tra
// cuu ma tang chien dau dung. Dien id khong ton tai = boss dung im khong danh
// duoc, va KHONG co thong bao loi nao.
func TestNoiDungBoss_MoiChieuTonTaiThatTrongCatalog(t *testing.T) {
	duyetPha(t, func(t *testing.T, nhan string, _ NoiDungBoss, _ int, pha *PhaBoss) {
		if len(pha.BoChieu) == 0 {
			t.Errorf("%s: BoChieu RONG — boss khong co gi de danh", nhan)
			return
		}
		if pha.Lang < 0 || int(pha.Lang) > 255 {
			t.Errorf("%s: Lang = %d, ngoai mien byte ma TraCapChieu/goi -5 dung", nhan, pha.Lang)
			return
		}
		for _, o := range pha.BoChieu {
			tid, cap := o[0], o[1]
			if tid < 1 || tid > 255 {
				t.Errorf("%s: templateId = %d — phai trong [1,255], id bat dau tu 1 (chieuthuc.go)", nhan, tid)
				continue
			}
			if cap < 1 || cap > 255 {
				t.Errorf("%s: cap chieu = %d — phai trong [1,255], cap tinh tu 1", nhan, cap)
				continue
			}
			if TraCapChieu(byte(pha.Lang), byte(tid), byte(cap)) == nil {
				t.Errorf("%s: chieu {templateId %d, cap %d} KHONG CO trong catalog cua lang %d — boss se dung im",
					nhan, tid, cap, pha.Lang)
			}
		}
	})
}

// TestNoiDungBoss_ChieuCuaBossPhaiLaChieuDANH: catalog chieu co 3 loai —
// danh (1), buff/hoi mau (2), tu than (3). Boss chi duoc cam chieu LOAI DANH:
//   - loai 2 (hoi mau) khong sinh goi tan cong nao, va duong hoi mau PvP
//     (goiHoiMauNguoiChoi) dang co y hoan (BAT-DAU-TU-DAY.md §8);
//   - loai 3 (tu than) co tam danh = 0 (tamTuThan trong chieuthuc.go) nen moi
//     don deu vuot tam -> LUON truot.
//
// Ca hai truong hop deu cho ra con boss "vung tay vao khong khi": no co ve
// dang hoat dong nhung khong bao gio gay mot diem sat thuong nao.
func TestNoiDungBoss_ChieuCuaBossPhaiLaChieuDANH(t *testing.T) {
	duyetPha(t, func(t *testing.T, nhan string, _ NoiDungBoss, _ int, pha *PhaBoss) {
		danhMuc := ChieuCuaLang(byte(pha.Lang))
		for _, o := range pha.BoChieu {
			tid := o[0]
			var mau *MauChieu
			for i := range danhMuc {
				if int32(danhMuc[i].Id) == tid {
					mau = &danhMuc[i]
					break
				}
			}
			if mau == nil {
				continue // da bao o test tren, khong bao hai lan
			}
			if mau.Loai != loaiDanh {
				t.Errorf("%s: chieu %q (templateId %d) co Loai = %d, boss chi duoc dung chieu DANH (loai %d)",
					nhan, mau.Ten, tid, mau.Loai, loaiDanh)
			}
		}
	})
}

// TestNoiDungBoss_ChieuGayDuocSatThuongThatSu: CapChieu.SatThuong duoc dung
// nhu PHAN TRAM cua Dame (TinhSatThuong, satthuong.go), tuc ket qua chia 100.
// Dame qua nho thi phep chia lam tron XUONG 0 va boss danh trung ma khong ai
// mat mau — dung loai loi khong bao gio bao.
func TestNoiDungBoss_ChieuGayDuocSatThuongThatSu(t *testing.T) {
	duyetPha(t, func(t *testing.T, nhan string, _ NoiDungBoss, _ int, pha *PhaBoss) {
		if pha.Dame <= 0 {
			t.Errorf("%s: Dame = %d — boss khong gay duoc sat thuong nao", nhan, pha.Dame)
			return
		}
		for _, o := range pha.BoChieu {
			cap := TraCapChieu(byte(pha.Lang), byte(o[0]), byte(o[1]))
			if cap == nil {
				continue // da bao o test tren
			}
			// He bi khac (×0.80) la truong hop TE NHAT — kiem o do thi moi
			// quan he ngu hanh con lai deu qua.
			if sat := TinhSatThuong(pha.Dame, cap, HeHoa, HeThuy, false); sat <= 0 {
				t.Errorf("%s: chieu {templateId %d, cap %d} voi Dame %d ra sat thuong %d — lam tron xuong 0, don danh trung ma khong ai mat mau",
					nhan, o[0], o[1], pha.Dame, sat)
			}
		}
	})
}

// TestNoiDungBoss_MoiMapVaoTonTaiThat: MapVao phai tro toi ban do CO THAT.
// TimBanDo tra nil cho map khong ton tai, va bo dieu phoi lich se khong biet
// gui lenh sinh di dau.
func TestNoiDungBoss_MoiMapVaoTonTaiThat(t *testing.T) {
	duyetPha(t, func(t *testing.T, nhan string, _ NoiDungBoss, _ int, pha *PhaBoss) {
		if len(pha.MapVao) == 0 {
			t.Errorf("%s: MapVao RONG — boss khong biet hien o dau", nhan)
			return
		}
		for _, id := range pha.MapVao {
			if TimBanDo(id) == nil {
				t.Errorf("%s: MapVao chua map %d — khong co ban do nao mang id do (banDo.go)", nhan, id)
			}
		}
	})
}

// TestNoiDungBoss_ViTriSinhDungTrenMatDat kiem HAI dieu, thieu mot la hong:
//
//  1. o dat chan phai DI DUOC — sinh trong vach da la boss ket cung mot cho,
//     nguoi choi danh khong toi, va khong co thong bao loi nao;
//  2. o NGAY DUOI phai DAC — neu khong thi cho do khong co mat dat, boss lo
//     lung giua troi (dung khuon kiem cua bong do, bgitem_test.go).
//
// Kiem bang CHINH luoi va cham cua map (collisionLangLa) chu khong bang mot
// bang toa do chep tay: sua bo cuc map ma quen doi vi tri boss thi test nay
// bao truoc, thay vi de nguoi choi phat hien.
func TestNoiDungBoss_ViTriSinhDungTrenMatDat(t *testing.T) {
	ds := DanhSachNoiDungBoss()
	for _, mau := range ds {
		if len(mau.ViTri) == 0 {
			t.Errorf("%s: khong co vi tri sinh nao", mau.Ma)
			continue
		}
		for i, v := range mau.ViTri {
			luoi, co := mapCoLuoiVaCham[v.Map]
			if !co {
				t.Errorf("%s: vi tri #%d o map %d — map do CHUA co luoi va cham, khong ai kiem duoc toa do co dung trong da khong (fail-closed)",
					mau.Ma, i, v.Map)
				continue
			}
			if !luoi.DiDuoc(v.X, v.Y) {
				t.Errorf("%s: vi tri #%d (%d,%d) map %d nam TRONG vach da — boss sinh trong long dat",
					mau.Ma, i, v.X, v.Y, v.Map)
			}
			if luoi.DiDuoc(v.X, v.Y+oTileVaCham) {
				t.Errorf("%s: vi tri #%d (%d,%d) map %d khong co mat dat ngay duoi — boss lo lung giua troi",
					mau.Ma, i, v.X, v.Y, v.Map)
			}
		}
	}
}

// TestNoiDungBoss_ViTriSinhNamTrongMapVao: hai bang du lieu (MapVao cua pha va
// ViTri cua boss) phai noi cung mot chuyen. Lech nhau nghia la hoac boss co
// cho dat chan o mot map no khong bao gio hien toi, hoac no hien vao mot map
// khong co cho dat chan nao — ca hai deu la du lieu chet, khong loi.
func TestNoiDungBoss_ViTriSinhNamTrongMapVao(t *testing.T) {
	for _, mau := range DanhSachNoiDungBoss() {
		trongMapVao := map[int32]bool{}
		for _, pha := range mau.Pha {
			for _, id := range pha.MapVao {
				trongMapVao[id] = true
			}
		}
		coChoDat := map[int32]bool{}
		for _, v := range mau.ViTri {
			coChoDat[v.Map] = true
			if !trongMapVao[v.Map] {
				t.Errorf("%s: co vi tri sinh o map %d nhung khong pha nao hien vao map do", mau.Ma, v.Map)
			}
		}
		for id := range trongMapVao {
			if !coChoDat[id] {
				t.Errorf("%s: hien vao map %d nhung khong co vi tri dat chan nao o do", mau.Ma, id)
			}
		}
	}
}

// TestNoiDungBoss_CapVaLangHopLe: Cap la truong 3 cua goi -5 (mot byte level),
// 0 lam client hien "cap 0" vo nghia — CapHopLe la cua kiem duy nhat.
// Lang di xuong day lam "gender"; client hien tai chi xu ly duoc 0/1/2 vi
// dung chung mang toc 3x3 (ghi chu o lang.go).
func TestNoiDungBoss_CapVaLangHopLe(t *testing.T) {
	duyetPha(t, func(t *testing.T, nhan string, _ NoiDungBoss, _ int, pha *PhaBoss) {
		if !pha.CapHopLe() {
			t.Errorf("%s: Cap = %d — ngoai mien 1..255, client hien cap vo nghia", nhan, pha.Cap)
		}
		if pha.Lang < 0 || pha.Lang > 2 {
			t.Errorf("%s: Lang = %d — client hien tai chi nhan 0/1/2 (lang.go: mang toc 3x3)", nhan, pha.Lang)
		}
		if TimLang(pha.Lang) == nil {
			t.Errorf("%s: Lang = %d khong co trong danh sach lang (lang.go)", nhan, pha.Lang)
		}
	})
}

// TestNoiDungBoss_TypePkKhac0: client (Char.cs paintHp) THOAT SOM khong ve
// thanh mau khi typePk == 0 va cFlag == 0, ma goi -5 cua ta luon ghi cFlag=0.
// De 0 thi nguoi choi danh boss ma khong thay boss con bao nhieu mau.
func TestNoiDungBoss_TypePkKhac0(t *testing.T) {
	duyetPha(t, func(t *testing.T, nhan string, _ NoiDungBoss, _ int, pha *PhaBoss) {
		if pha.TypePk != TypePkBoss {
			t.Errorf("%s: TypePk = %d, phai la TypePkBoss = %d (0 = nguoi choi khong thay thanh mau boss)",
				nhan, pha.TypePk, TypePkBoss)
		}
	})
}

// TestNoiDungBoss_OutfitDungPhanDaCoAnhThat: dau/than/chan phai la id part ma
// client CHAC CHAN co anh (xem phanDaCoAnhThat). Ba o con lai (flagBag, aura,
// effFront) phai la 0: chua co anh Naruto nao cho chung, va hai trong ba o do
// chi MOT BYTE tren day nen so > 255 bi cat am tham.
func TestNoiDungBoss_OutfitDungPhanDaCoAnhThat(t *testing.T) {
	tenO := [6]string{"head", "body", "leg", "flagBag", "aura", "effFront"}
	duyetPha(t, func(t *testing.T, nhan string, _ NoiDungBoss, _ int, pha *PhaBoss) {
		for i := 0; i < 3; i++ {
			if !phanDaCoAnhThat[pha.Outfit[i]] {
				t.Errorf("%s: Outfit[%d] (%s) = %d — khong nam trong closure.part_ids cua asset-lock.json, boss se VO HINH",
					nhan, i, tenO[i], pha.Outfit[i])
			}
		}
		for i := 3; i < 6; i++ {
			if pha.Outfit[i] != 0 {
				t.Errorf("%s: Outfit[%d] (%s) = %d — chua co anh cho o nay, phai de 0",
					nhan, i, tenO[i], pha.Outfit[i])
			}
		}
	})
}

// TestNoiDungBoss_HpChonKhongRongVaVuaInt32: Hp/HpToiDa cua Boss la int32.
// Mot gia tri trong HpChon vuot 2_147_483_647 se bi KepHp cat xuong tran ma
// KHONG bao gi — con boss "3 ty mau" am tham thanh con boss 2,1 ty mau, va
// chenh lech do khong bao gio lo ra o dau.
func TestNoiDungBoss_HpChonKhongRongVaVuaInt32(t *testing.T) {
	duyetPha(t, func(t *testing.T, nhan string, _ NoiDungBoss, _ int, pha *PhaBoss) {
		if len(pha.HpChon) == 0 {
			t.Errorf("%s: HpChon RONG — luc sinh khong biet cho boss bao nhieu mau", nhan)
			return
		}
		for i, hp := range pha.HpChon {
			if hp <= 0 {
				t.Errorf("%s: HpChon[%d] = %d — boss sinh ra la chet ngay", nhan, i, hp)
				continue
			}
			if int64(KepHp(hp)) != hp {
				t.Errorf("%s: HpChon[%d] = %d bi KepHp cat con %d — vuot mien int32 cua Hp/HpToiDa",
					nhan, i, hp, KepHp(hp))
			}
		}
	})
}

// TestNoiDungBoss_PhaLienTucKhongRong: PhaHienTai la chi so chay tu 0 len het
// mang khi boss bien hinh. Moi chi so trong [0, len) phai tra ra mot pha
// THAT (PhaHienTaiData khac nil) va pha do phai co du thit — thieu ten thi
// nguoi choi thay bong bong ten rong tren dau boss.
func TestNoiDungBoss_PhaLienTucKhongRong(t *testing.T) {
	for _, mau := range DanhSachNoiDungBoss() {
		if len(mau.Pha) == 0 {
			t.Errorf("%s: khong co pha nao", mau.Ma)
			continue
		}
		b := &Boss{DuLieu: mau.Pha}
		for i := range mau.Pha {
			b.PhaHienTai = i
			pha := b.PhaHienTaiData()
			if pha == nil {
				t.Errorf("%s: PhaHienTai = %d tra nil — day pha bi dut", mau.Ma, i)
				continue
			}
			if strings.TrimSpace(pha.Ten) == "" {
				t.Errorf("%s pha %d: Ten rong — client ve bong ten trong tren dau boss", mau.Ma, i+1)
			}
			if pha.GiayNghi <= 0 {
				t.Errorf("%s pha %d: GiayNghi = %d — phai duong, day la so giay nam yen truoc khi hien",
					mau.Ma, i+1, pha.GiayNghi)
			}
		}
		// Ngay sau pha cuoi phai la NIL — do la dau hieu "het pha" ma may
		// trang thai dung de biet khi nao phat -6 va ve NGHI.
		b.PhaHienTai = len(mau.Pha)
		if b.PhaHienTaiData() != nil {
			t.Errorf("%s: vuot qua pha cuoi ma van tra ra pha — may trang thai se bien hinh vo tan", mau.Ma)
		}
	}
}

// TestNoiDungBoss_MocNoiGianGiamDanKhongTrung: MocNoiGianDaQua la con tro chay
// MOT CHIEU (boss.go) — may trang thai duyet mang theo thu tu khai bao khi mau
// tut xuong. Vi vay mang PHAI sap theo %HP GIAM DAN va khong duoc trung
// nguong:
//   - sap sai thu tu -> moc sau khong bao gio toi luot, du lieu chet;
//   - trung nguong  -> hai moc cung cham mot luc, mot cau chat bi nuot.
//
// Nguong cung phai nam trong (0, 100): 100 la boss noi gian ngay khi vua hien
// (chua ai danh), 0 thi khong bao gio toi vi mau ve 0 la boss chet.
func TestNoiDungBoss_MocNoiGianGiamDanKhongTrung(t *testing.T) {
	duyetPha(t, func(t *testing.T, nhan string, _ NoiDungBoss, _ int, pha *PhaBoss) {
		truoc := 100
		for i, moc := range pha.MocNoiGian {
			if moc.PhanTramHP <= 0 || moc.PhanTramHP >= 100 {
				t.Errorf("%s moc #%d: PhanTramHP = %d — phai trong (0,100)", nhan, i, moc.PhanTramHP)
			}
			if moc.PhanTramHP >= truoc {
				t.Errorf("%s moc #%d: PhanTramHP = %d khong nho hon moc truoc (%d) — mang phai giam dan, khong trung",
					nhan, i, moc.PhanTramHP, truoc)
			}
			truoc = moc.PhanTramHP
			if strings.TrimSpace(moc.Chat) == "" {
				t.Errorf("%s moc #%d: Chat rong — noi gian ma khong noi gi thi nguoi choi khong biet gi doi", nhan, i)
			}
		}
	})
}

// TestNoiDungBoss_MocNoiGianChiManhLenKhongYeuDi: "noi gian" ma boss danh nhe
// di hoac mong hon la du lieu nguoc nghia. HeSoDame nhan vao sat thuong boss
// GAY RA (lon hon = manh hon); HeSoDef nhan vao sat thuong boss NHAN (nho hon
// = trau hon) — hai chieu nguoc nhau, rat de dien nham.
func TestNoiDungBoss_MocNoiGianChiManhLenKhongYeuDi(t *testing.T) {
	duyetPha(t, func(t *testing.T, nhan string, _ NoiDungBoss, _ int, pha *PhaBoss) {
		dameTruoc, defTruoc := 1.0, 1.0
		for i, moc := range pha.MocNoiGian {
			if moc.HeSoDame < dameTruoc {
				t.Errorf("%s moc #%d: HeSoDame = %g yeu hon muc truoc (%g) — noi gian ma danh nhe di",
					nhan, i, moc.HeSoDame, dameTruoc)
			}
			if moc.HeSoDef <= 0 || moc.HeSoDef > defTruoc {
				t.Errorf("%s moc #%d: HeSoDef = %g phai trong (0, %g] — noi gian ma mong hon truoc",
					nhan, i, moc.HeSoDef, defTruoc)
			}
			dameTruoc, defTruoc = moc.HeSoDame, moc.HeSoDef
		}
	})
}

// TestNoiDungBoss_NhipDanhHopVoiNhipActor: boss cap nhat TRONG nhip 100ms cua
// khu vuc (boss.go dau file), nen mot nhip danh 950ms tren thuc te van la
// 1000ms — phan le 50ms la con so noi doi. Yeu cau boi so cua 100 de con so
// trong du lieu dung bang con so xay ra trong game.
//
// Kiem ca nhip SAU KHI da nhan het he so noi gian (chung nhan don len
// NhipDanhHienTai, tuc CONG DON qua tung moc): noi gian chong nhau co the keo
// nhip xuong duoi mot tick, luc do boss danh moi tick — 10 don mot giay.
func TestNoiDungBoss_NhipDanhHopVoiNhipActor(t *testing.T) {
	const motTickMs = 100
	duyetPha(t, func(t *testing.T, nhan string, _ NoiDungBoss, _ int, pha *PhaBoss) {
		if pha.NhipDanhMs < motTickMs {
			t.Errorf("%s: NhipDanhMs = %d — nho hon mot nhip actor (%dms), boss danh moi tick",
				nhan, pha.NhipDanhMs, motTickMs)
			return
		}
		if pha.NhipDanhMs%motTickMs != 0 {
			t.Errorf("%s: NhipDanhMs = %d khong phai boi so cua %dms — phan le khong bao gio xay ra that",
				nhan, pha.NhipDanhMs, motTickMs)
		}
		nhip := float64(pha.NhipDanhMs)
		for i, moc := range pha.MocNoiGian {
			if moc.HeSoNhip < 0 {
				t.Errorf("%s moc #%d: HeSoNhip = %g am", nhan, i, moc.HeSoNhip)
				continue
			}
			if moc.HeSoNhip == 0 {
				continue // 0 = khong doi nhip, dung nghia trong boss.go
			}
			nhip *= moc.HeSoNhip
			if nhip < motTickMs {
				t.Errorf("%s: sau %d moc noi gian, nhip danh con %.0fms — duoi mot tick %dms, boss danh 10 don/giay",
					nhan, i+1, nhip, motTickMs)
			}
		}
	})
}

// TestNoiDungBoss_ThoaiDayDuVaVuaGoi44: ba mang thoai deu duoc may trang thai
// dung o mot buoc khac nhau (CHAO / HOAT_DONG / TROI). Mang rong = buoc do im
// lang; cau rong = mot bong bong thoai trong tren dau boss. Do dai phai vua
// khuon writeUTF cua goi 44, neu khong cau bi cat giua chung.
func TestNoiDungBoss_ThoaiDayDuVaVuaGoi44(t *testing.T) {
	duyetPha(t, func(t *testing.T, nhan string, _ NoiDungBoss, _ int, pha *PhaBoss) {
		for _, bo := range []struct {
			ten  string
			cauc []string
		}{
			{"ThoaiHien", pha.ThoaiHien},
			{"ThoaiDanh", pha.ThoaiDanh},
			{"ThoaiChet", pha.ThoaiChet},
		} {
			if len(bo.cauc) == 0 {
				t.Errorf("%s: %s RONG — buoc tuong ung cua may trang thai se im lang", nhan, bo.ten)
				continue
			}
			for i, cau := range bo.cauc {
				if strings.TrimSpace(cau) == "" {
					t.Errorf("%s: %s[%d] rong — bong bong thoai trong tren dau boss", nhan, bo.ten, i)
				}
				if len([]rune(cau)) > soKyTuToiDaChat {
					t.Errorf("%s: %s[%d] dai %d ky tu, vuot %d ma goi 44 mang duoc — se bi cat",
						nhan, bo.ten, i, len([]rune(cau)), soKyTuToiDaChat)
				}
			}
		}
	})
}

// TestNoiDungBoss_MaDuyNhatVaTraCuuDung: bo dieu phoi lich (task 13) goi boss
// theo ma. Ma trung nghia la TimNoiDungBoss tra ra con dau tien, con thu hai
// khong bao gio duoc sinh — mot con boss bien mat khoi game ma khong loi.
func TestNoiDungBoss_MaDuyNhatVaTraCuuDung(t *testing.T) {
	daThay := map[string]bool{}
	for _, mau := range DanhSachNoiDungBoss() {
		if strings.TrimSpace(mau.Ma) == "" {
			t.Error("co boss mau khong co Ma — khong tra cuu duoc")
			continue
		}
		if daThay[mau.Ma] {
			t.Errorf("Ma %q trung — con thu hai khong bao gio duoc sinh", mau.Ma)
		}
		daThay[mau.Ma] = true

		got := TimNoiDungBoss(mau.Ma)
		if got == nil {
			t.Errorf("TimNoiDungBoss(%q) tra nil du ma nay co trong danh sach", mau.Ma)
			continue
		}
		if got.Ma != mau.Ma || len(got.Pha) != len(mau.Pha) {
			t.Errorf("TimNoiDungBoss(%q) tra nham con: ma %q, %d pha", mau.Ma, got.Ma, len(got.Pha))
		}
	}
	if TimNoiDungBoss("khong-co-con-nao-ten-nay") != nil {
		t.Error("TimNoiDungBoss tra khac nil cho ma khong ton tai — phai fail-closed")
	}
}

// TestDanhSachNoiDungBoss_TraBanSaoSau: nguoi goi sua ket qua KHONG duoc dung
// toi danh muc he thong — cung ki luat voi DanhSachBanDo/TimBanDo (banDo.go).
// Day la loi sao chep nong kinh dien: sao struct chi sao CON TRO cua lat cat,
// nen nguoi goi van ghi de duoc du lieu goc qua ban "sao chep".
func TestDanhSachNoiDungBoss_TraBanSaoSau(t *testing.T) {
	ds := DanhSachNoiDungBoss()
	if len(ds) == 0 || len(ds[0].Pha) == 0 {
		t.Fatal("khong co du lieu de kiem")
	}
	tenGoc := ds[0].Pha[0].Ten
	maGoc := ds[0].Ma

	ds[0].Ma = "BI-GHI-DE"
	ds[0].Pha[0].Ten = "BI-GHI-DE"
	if len(ds[0].ViTri) > 0 {
		ds[0].ViTri[0].X = -12345
	}

	lai := DanhSachNoiDungBoss()
	if lai[0].Ma != maGoc || lai[0].Pha[0].Ten != tenGoc {
		t.Errorf("sua ban sao da dung toi danh muc goc: ma %q, ten pha %q", lai[0].Ma, lai[0].Pha[0].Ten)
	}
	if len(lai[0].ViTri) > 0 && lai[0].ViTri[0].X == -12345 {
		t.Error("sua ViTri cua ban sao da dung toi danh muc goc")
	}

	tra := TimNoiDungBoss(maGoc)
	if tra == nil {
		t.Fatalf("TimNoiDungBoss(%q) tra nil sau khi ban sao bi sua", maGoc)
	}
	tra.Pha[0].Ten = "BI-GHI-DE-LAN-HAI"
	if TimNoiDungBoss(maGoc).Pha[0].Ten != tenGoc {
		t.Error("sua ket qua TimNoiDungBoss da dung toi danh muc goc")
	}
}
