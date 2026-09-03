// boss_test.go — bo test cua boss.yaml (KH8 task 7).
//
// Cach lam GIONG HET nap_test.go: lay kho noi dung THAT (ban nhung), do ra van
// ban, VA MOT DOT BIEN vao van ban do, roi doi bo nap TU CHOI. Dung van ban that
// (thay vi fixture bia) de test khong the xanh nho mot fixture da don gian hoa.
//
// Bon dot bien co ten trong thiet ke §6 deu co test rieng o day:
//
//	M7  nhip_danh_ms khong phai boi so 100ms  -> TestBoss_M7_*
//	M8  giay_nghi cua mot pha = 0             -> TestBoss_M8_*
//	M10 boss cam chieu loai BUFF              -> TestBoss_M10_*
//	M11 (vi tri sinh doi vao o dac/lo lung)   -> KHONG o day: can LUOI VA CHAM,
//	    artifact cua goi `game`. Xem TestNoiDungBoss_ViTriSinhDungTrenMatDat
//	    (game/noidungboss_test.go). O day chi kiem duoc mien so (toa do am/0).
package noidung

import "testing"

// TestBoss_KhoThatNapDuoc: boss.yaml dang commit phai nap duoc va co du hai con.
// Day la NEO cho moi test dot bien duoi day: neu ban goc da do thi "dot bien lam
// do" khong chung minh dieu gi.
func TestBoss_KhoThatNapDuoc(t *testing.T) {
	du, err := NapNhung()
	if err != nil {
		t.Fatalf("kho noi dung dang commit khong nap duoc: %v", err)
	}
	if len(du.Boss) != 2 {
		t.Fatalf("%d con boss, mong doi 2", len(du.Boss))
	}
	// Ba duong dac trung cua he boss phai chay duoc bang chinh bo du lieu nay,
	// neu khong thi buoi nghiem thu bang client that khong chung minh duoc gi.
	var coNhieuPha, coNoiGian, coBossLon bool
	for _, b := range du.Boss {
		if len(b.Pha) >= 2 {
			coNhieuPha = true
		}
		for _, p := range b.Pha {
			if len(p.MocNoiGian) > 0 {
				coNoiGian = true
			}
			if p.LaBossLon {
				coBossLon = true
			}
		}
	}
	if !coNhieuPha {
		t.Error("khong con boss nao co tu 2 pha tro len — khong chay duoc duong bien hinh")
	}
	if !coNoiGian {
		t.Error("khong pha nao co moc_noi_gian — khong chay duoc duong noi gian theo %mau")
	}
	if !coBossLon {
		t.Error("khong pha nao co la_boss_lon — khong chay duoc duong chay chu toan server (goi 93)")
	}
}

// napVoiDotBienBoss doi mot chuoi trong boss.yaml roi nap lai, tra ve loi.
func napVoiDotBienBoss(t *testing.T, cu, moi string) error {
	t.Helper()
	van := khoThat(t)
	doiMotLan(t, van, TenFileBoss, cu, moi)
	_, err := napVan(van)
	return err
}

// doiPhaiDo goi napVoiDotBienBoss va bat buoc no DO, kem mot dong noi vi sao dot
// bien nay nguy hiem.
func doiPhaiDo(t *testing.T, cu, moi, viSao string) {
	t.Helper()
	if err := napVoiDotBienBoss(t, cu, moi); err == nil {
		t.Errorf("dot bien %q -> %q duoc nap: %s", cu, moi, viSao)
	}
}

// ── M7 ─────────────────────────────────────────────────────────────────────

// TestBoss_M7_NhipDanhKhongPhaiBoiSoNhipActor: boss cap nhat TRONG nhip 100ms
// cua khu vuc, nen 1450ms tren thuc te van la 1500ms — phan le 50ms la mot con
// so NOI DOI: no nam trong file du lieu nhung khong bao gio xay ra trong game.
func TestBoss_M7_NhipDanhKhongPhaiBoiSoNhipActor(t *testing.T) {
	doiPhaiDo(t, "nhip_danh_ms: 1500", "nhip_danh_ms: 1450",
		"1450ms khong bao gio xay ra that — actor chi tinh theo nhip 100ms")
}

// TestBoss_M7_NhipDanhNhoHonMotNhipActor: duoi mot tick thi boss danh MOI TICK —
// 10 don mot giay, khong ai song noi.
func TestBoss_M7_NhipDanhNhoHonMotNhipActor(t *testing.T) {
	doiPhaiDo(t, "nhip_danh_ms: 1500", "nhip_danh_ms: 0",
		"nhip 0 = boss danh moi tick, 10 don mot giay")
}

// TestBoss_M7_NoiGianKeoNhipXuongDuoiMotTick: he_so_nhip nhan CONG DON. Ba moc
// ×0,1 se keo 1500ms xuong 1,5ms — van "hop le" neu chi nhin tung con so rieng.
func TestBoss_M7_NoiGianKeoNhipXuongDuoiMotTick(t *testing.T) {
	doiPhaiDo(t, "he_so_nhip: 0.80", "he_so_nhip: 0.01",
		"he so nhip nhan cong don, keo nhip xuong duoi mot tick 100ms")
}

// ── M8 ─────────────────────────────────────────────────────────────────────

// TestBoss_M8_GiayNghiCuaPhaSauBangKhong: giay_nghi cua pha 2 la KHOANG LANG
// giua hai pha. Bang 0 nghia la pha 2 no ra ngay trong cung mot khoanh khac pha
// 1 chet — nguoi choi khong kip doc gi, khong kip thay mat na roi xuong.
func TestBoss_M8_GiayNghiCuaPhaSauBangKhong(t *testing.T) {
	doiPhaiDo(t, "giay_nghi: 5", "giay_nghi: 0",
		"0 = pha sau no ra ngay khoanh khac pha truoc chet, khong ai kip doc gi")
}

// TestBoss_M8_GiayNghiAm: so am con te hon 0 — bo dieu phoi lich so sanh
// `bayGio.Sub(lanSinhCuoi) < GiayNghi giay`, so am lam dieu kien LUON dung, tuc
// boss sinh lai moi nhip cho toi khi khu vuc ngop.
func TestBoss_M8_GiayNghiAm(t *testing.T) {
	doiPhaiDo(t, "giay_nghi: 1800", "giay_nghi: -1",
		"giay_nghi am = bo dieu phoi sinh boss moi nhip 100ms")
}

// ── M10 ────────────────────────────────────────────────────────────────────

// TestBoss_M10_BossCamChieuLoaiBuff: template_id 3 cua moi lang la chieu HOI MAU
// (loai BUFF). Boss cam no se "vung tay vao khong khi": trong nhu dang danh nhau
// ma khong bao gio gay mot diem sat thuong nao — va khong mot dong log.
//
// Day la dot bien duy nhat trong nhom nay can CA HAI file (boss.yaml + chieu.yaml),
// nen no chay o kiemDayDu chu khong o kiemBoss.
func TestBoss_M10_BossCamChieuLoaiBuff(t *testing.T) {
	doiPhaiDo(t, "- {template_id: 1, cap: 2}", "- {template_id: 3, cap: 2}",
		"chieu loai BUFF khong sinh goi tan cong nao — boss vung tay vao khong khi")
}

// TestBoss_M10_BossCamChieuKhongTonTai: templateId 9 khong co trong catalog cua
// lang nao — boss dung im.
func TestBoss_M10_BossCamChieuKhongTonTai(t *testing.T) {
	doiPhaiDo(t, "- {template_id: 1, cap: 2}", "- {template_id: 9, cap: 2}",
		"template_id khong co trong catalog — boss dung im, khong loi")
}

// TestBoss_M10_BossCamCapChieuVuotCatalog: catalog chi co 3 cap moi chieu.
func TestBoss_M10_BossCamCapChieuVuotCatalog(t *testing.T) {
	doiPhaiDo(t, "- {template_id: 1, cap: 2}", "- {template_id: 1, cap: 9}",
		"cap 9 khong ton tai — tra cuu tra nil, boss dung im")
}

// ── M11 (phan kiem duoc o goi nay: mien so cua toa do) ─────────────────────

// TestBoss_M11_ViTriSinhCot0: cot/hang 0 cua luoi va cham la TUONG BIEN. Phan
// con lai cua M11 (o dac / lo lung giua troi) can luoi va cham that va nam o
// game/noidungboss_test.go.
func TestBoss_M11_ViTriSinhCot0(t *testing.T) {
	doiPhaiDo(t, "{map: 1, x: 552, y: 408}", "{map: 1, x: 0, y: 408}",
		"cot 0 la tuong bien — boss nam trong da")
}

// TestBoss_M11_ViTriSinhOMapKhongCoPhaNaoVao: hai bang du lieu (map_vao cua pha
// va vi_tri cua boss) phai noi CUNG MOT chuyen. Lech nhau la du lieu chet.
func TestBoss_M11_ViTriSinhOMapKhongCoPhaNaoVao(t *testing.T) {
	doiPhaiDo(t, "{map: 1, x: 552, y: 408}", "{map: 4, x: 552, y: 408}",
		"co cho dat chan o map 4 nhung khong pha nao hien vao map 4")
}

// ── Cac kiem con lai (nhung gi noidungboss_test.go dang ghim) ──────────────

func TestBoss_MocNoiGianPhaiGiamDanKhongTrung(t *testing.T) {
	doiPhaiDo(t, "phan_tram_hp: 20", "phan_tram_hp: 50",
		"hai moc trung nguong: cung cham mot luc, mot cau chat bi nuot")
	doiPhaiDo(t, "phan_tram_hp: 20", "phan_tram_hp: 60",
		"moc sau CAO hon moc truoc: con tro moc chay mot chieu nen no khong bao gio toi luot")
	doiPhaiDo(t, "phan_tram_hp: 50", "phan_tram_hp: 100",
		"100 = boss noi gian ngay khi vua hien, chua ai danh")
	doiPhaiDo(t, "phan_tram_hp: 50", "phan_tram_hp: 0",
		"0 khong bao gio toi vi mau ve 0 la boss chet")
}

func TestBoss_NoiGianPhaiManhLenKhongYeuDi(t *testing.T) {
	doiPhaiDo(t, "he_so_dame: 1.60", "he_so_dame: 1.10",
		"noi gian ma danh NHE di hon muc truoc — du lieu nguoc nghia")
	doiPhaiDo(t, "he_so_def: 0.70", "he_so_def: 0.95",
		"noi gian ma MONG hon truoc (he_so_def lon hon = nhan nhieu sat thuong hon)")
	doiPhaiDo(t, "he_so_def: 0.85", "he_so_def: 0.0",
		"he_so_def = 0 lam boss BAT TU")
}

func TestBoss_HpChonPhaiDuongVaVuaInt32(t *testing.T) {
	doiPhaiDo(t, "hp_chon: [3600, 4500, 5400]", "hp_chon: [3600, 0, 5400]",
		"hp 0 = boss sinh ra la chet ngay")
	doiPhaiDo(t, "hp_chon: [3600, 4500, 5400]", "hp_chon: [3600, 3000000000, 5400]",
		"vuot int32: KepHp cat xuong tran ma khong bao gi, con boss 3 ty mau am tham thanh 2,1 ty")
	doiPhaiDo(t, "hp_chon: [3600, 4500, 5400]", "hp_chon: []",
		"hp_chon rong = luc sinh khong biet cho boss bao nhieu mau")
}

func TestBoss_LangPhaiTrongNamLang(t *testing.T) {
	doiPhaiDo(t, "lang: 2", "lang: 9",
		"lang 9 khong ton tai — bo chieu tra cuu vao mot catalog rong")
}

func TestBoss_CapVaTypePkPhaiHopLe(t *testing.T) {
	doiPhaiDo(t, "cap: 50", "cap: 0",
		"cap 0 lam client hien \"cap 0\" vo nghia (goi -5 ghi cap bang mot byte)")
	doiPhaiDo(t, "cap: 50", "cap: 300",
		"cap 300 vuot mot byte cua goi -5")
	doiPhaiDo(t, "type_pk: 5", "type_pk: 0",
		"type_pk 0: client THOAT SOM khong ve thanh mau — nguoi choi danh boss ma khong thay mau boss")
}

func TestBoss_MaVaTenPhaiCoVaDuyNhat(t *testing.T) {
	doiPhaiDo(t, "ma: hoc-vien-dao-tau", "ma: nhan-gia-mat-na",
		"ma trung: tra cuu tra ra con dau tien, con thu hai khong bao gio duoc sinh")
	doiPhaiDo(t, "ma: hoc-vien-dao-tau", `ma: ""`,
		"ma rong: bo dieu phoi lich khong goi duoc con boss nay")
	doiPhaiDo(t, `ten: "Mặt Nạ Vỡ"`, `ten: ""`,
		"ten rong: client ve bong ten TRONG tren dau boss")
}

func TestBoss_ThoaiPhaiDayDu(t *testing.T) {
	doiPhaiDo(t, `- "Đó mới chỉ là lớp áo... còn ta thì chưa."`, `- ""`,
		"cau thoai rong = bong bong thoai trong tren dau boss")
	doiPhaiDo(t, "thoai_chet:\n          - \"Vậy là... ta được về nhà theo cách này.\"",
		"thoai_chet: []",
		"thoai_chet rong = buoc TROI cua may trang thai im lang")
}

func TestBoss_DameVaMapVaoPhaiCo(t *testing.T) {
	doiPhaiDo(t, "dame: 35", "dame: 0",
		"dame 0 = boss khong gay duoc sat thuong nao (sat thuong = dame × %cap chieu)")
	doiPhaiDo(t, "map_vao: [1]\n        bo_chieu:", "map_vao: []\n        bo_chieu:",
		"map_vao rong = boss khong biet hien o dau")
}

// TestBoss_GoNhamTenTruongBiTuChoi: KnownFields(true) trong docYAML lam mot ten
// truong go nham thanh LOI thay vi thanh gia tri 0 IM LANG. Day la lop chan cuoi
// cung: khong bo kiem nao o tren bat duoc "truong bi bo qua vi go sai ten".
func TestBoss_GoNhamTenTruongBiTuChoi(t *testing.T) {
	doiPhaiDo(t, "nhip_danh_ms: 1200", "nhip_dahn_ms: 1200",
		"go nham ten truong: khong go nham thi nhip = 0, boss danh moi tick")
}

// ---------------------------------------------------------------------------
// Cua tra cuu tren anh chup
// ---------------------------------------------------------------------------

// TestAnhChupBoss_TraBanSaoSau: nguoi goi sua ket qua KHONG duoc dung toi anh
// chup DA CONG BO — sua no la sua cho MOI nguoi choi khac, va loi chi lo ra o
// LAN SINH BOSS SAU, o mot khu vuc khac, rat lau sau khi cho sua tay xay ra.
//
// Sao BA LOP: boss -> pha -> cac lat cat trong pha. Bo sot mot lop la thung.
func TestAnhChupBoss_TraBanSaoSau(t *testing.T) {
	du, err := NapNhung()
	if err != nil {
		t.Fatalf("nap: %v", err)
	}
	ac := MoKho(du).AnhChup()

	ds := ac.DanhSachBoss()
	if len(ds) == 0 || len(ds[0].Pha) == 0 || len(ds[0].ViTri) == 0 {
		t.Fatal("khong co du lieu de kiem")
	}
	maGoc := ds[0].Ma
	tenGoc := ds[0].Pha[0].Ten
	hpGoc := ds[0].Pha[0].HpChon[0]
	chatGoc := ds[0].Pha[0].MocNoiGian[0].Chat
	thoaiGoc := ds[0].Pha[0].ThoaiHien[0]
	chieuGoc := ds[0].Pha[0].BoChieu[0]
	mapGoc := ds[0].Pha[0].MapVao[0]

	// Lop 1 + 2 + 3.
	ds[0].Ma = "BI-GHI-DE"
	ds[0].ViTri[0].X = -12345
	ds[0].Pha[0].Ten = "BI-GHI-DE"
	ds[0].Pha[0].HpChon[0] = 1
	ds[0].Pha[0].MapVao[0] = 999
	ds[0].Pha[0].BoChieu[0] = OChieuBoss{TemplateId: 99, Cap: 99}
	ds[0].Pha[0].MocNoiGian[0].Chat = "BI-GHI-DE"
	ds[0].Pha[0].ThoaiHien[0] = "BI-GHI-DE"

	lai := ac.DanhSachBoss()
	if lai[0].Ma != maGoc {
		t.Errorf("Ma bi ghi de: %q", lai[0].Ma)
	}
	if lai[0].ViTri[0].X == -12345 {
		t.Error("ViTri bi ghi de qua ban sao")
	}
	if lai[0].Pha[0].Ten != tenGoc {
		t.Errorf("Ten pha bi ghi de: %q", lai[0].Pha[0].Ten)
	}
	if lai[0].Pha[0].HpChon[0] != hpGoc {
		t.Errorf("HpChon bi ghi de: %d", lai[0].Pha[0].HpChon[0])
	}
	if lai[0].Pha[0].MapVao[0] != mapGoc {
		t.Errorf("MapVao bi ghi de: %d", lai[0].Pha[0].MapVao[0])
	}
	if lai[0].Pha[0].BoChieu[0] != chieuGoc {
		t.Errorf("BoChieu bi ghi de: %+v", lai[0].Pha[0].BoChieu[0])
	}
	if lai[0].Pha[0].MocNoiGian[0].Chat != chatGoc {
		t.Errorf("MocNoiGian bi ghi de: %q", lai[0].Pha[0].MocNoiGian[0].Chat)
	}
	if lai[0].Pha[0].ThoaiHien[0] != thoaiGoc {
		t.Errorf("ThoaiHien bi ghi de: %q", lai[0].Pha[0].ThoaiHien[0])
	}

	// Cua Boss(ma) cung phai sao sau — no la cua ma bo dieu phoi lich dung.
	b, co := ac.Boss(maGoc)
	if !co {
		t.Fatalf("Boss(%q) tra false", maGoc)
	}
	b.Pha[0].Ten = "BI-GHI-DE-LAN-HAI"
	b.Pha[0].ThoaiDanh[0] = "BI-GHI-DE-LAN-HAI"
	lai2, _ := ac.Boss(maGoc)
	if lai2.Pha[0].Ten != tenGoc {
		t.Error("sua ket qua Boss(ma) da dung toi anh chup goc")
	}
	if lai2.Pha[0].ThoaiDanh[0] == "BI-GHI-DE-LAN-HAI" {
		t.Error("sua ThoaiDanh cua ban sao da dung toi anh chup goc")
	}

	// Fail-closed: ma khong ton tai -> false, KHONG phai mot Boss rong "hop le".
	if _, co := ac.Boss("khong-co-con-nao-ten-nay"); co {
		t.Error("Boss(ma khong ton tai) tra true — phai fail-closed")
	}
}

// ---------------------------------------------------------------------------
// KiemIdPartBoss — ham THUAN, test bang danh sach id gia
// ---------------------------------------------------------------------------

// TestKiemIdPartBoss: outfit tro vao mot id part KHONG co trong catalog anh da
// ky = con boss VO HINH (client khong tim thay anh thi khong ve gi, may chu
// khong co mot dong log nao). Ba o con lai phai la 0 — chua co anh cho chung.
func TestKiemIdPartBoss(t *testing.T) {
	hopLe := []int16{30, 57, 58}
	tot := []Boss{{Ma: "x", Pha: []PhaBoss{{Outfit: [6]int16{30, 57, 58, 0, 0, 0}}}}}
	if err := KiemIdPartBoss(tot, hopLe); err != nil {
		t.Fatalf("outfit hop le ma bi tu choi: %v", err)
	}

	voHinh := []Boss{{Ma: "x", Pha: []PhaBoss{{Outfit: [6]int16{9999, 57, 58, 0, 0, 0}}}}}
	if err := KiemIdPartBoss(voHinh, hopLe); err == nil {
		t.Error("id part khong co trong catalog ma qua duoc — boss se VO HINH")
	}

	thuaO := []Boss{{Ma: "x", Pha: []PhaBoss{{Outfit: [6]int16{30, 57, 58, 0, 7, 0}}}}}
	if err := KiemIdPartBoss(thuaO, hopLe); err == nil {
		t.Error("o aura khac 0 ma qua duoc — chua co anh cho o do")
	}

	// Kho THAT phai qua duoc bo kiem nay voi dung tap id da ky.
	du, err := NapNhung()
	if err != nil {
		t.Fatalf("nap: %v", err)
	}
	if err := KiemIdPartBoss(du.Boss, []int16{30, 57, 58, 2189, 2190, 2191}); err != nil {
		t.Errorf("boss.yaml dang commit lech catalog anh da ky: %v", err)
	}
}
