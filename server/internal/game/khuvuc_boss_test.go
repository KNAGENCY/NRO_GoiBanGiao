package game

import (
	"strings"
	"testing"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// ---------------------------------------------------------------------------
// KH6b Task 7 — boss cam vao actor KhuVuc: nguoi moi vao PHAI thay boss dang
// hien (goi -5), va themBoss phat -5 cho moi nguoi dang o trong khu vuc.
//
// TAT CA test o day goi k.xuLy(...) TRUC TIEP (actor CHUA BatDau) — cung khuon
// voi khuvuc_test.go: tach logic xu ly khoi vong lap channel, khong goroutine,
// khong sleep, khong dong ho that.
// ---------------------------------------------------------------------------

// goi5DaDoc la goi -5 da duoc doc nguoc ra thanh cac truong co ten. Test boss
// phai khang dinh NHIEU truong cung luc (id am, ten, cap, hp, toa do, typePk,
// outfit) nen doc mot lan roi so sanh tung o de hon han lap lai chuoi doc byte
// o moi test.
type goi5DaDoc struct {
	id       int32
	cap      byte
	typePk   byte
	lang     byte
	head     int16
	ten      string
	hp       int32
	hpToiDa  int32
	body     int16
	leg      int16
	flagBag  byte
	x, y     int16
	aura     int16
	effFront byte
}

// docGoi5 doc DU 28 truong cua goi -5 theo dung thu tu goiXuatHienDayDu ghi ra
// (xuathien.go) va khang dinh het buffer sau truong 28 — doc thieu mot truong
// la moi khang dinh sau do lech ma khong ai biet.
func docGoi5(t *testing.T, m *protocol.Message) goi5DaDoc {
	t.Helper()
	if m == nil {
		t.Fatal("goi -5 la nil")
	}
	if m.Command != lenhXuatHien {
		t.Fatalf("command = %d, mong doi %d (-5)", m.Command, lenhXuatHien)
	}
	r := protocol.NewReader(m.Data)
	var g goi5DaDoc
	kiem := func(err error, ten string) {
		t.Helper()
		if err != nil {
			t.Fatalf("doc truong %s loi: %v", ten, err)
		}
	}
	var err error
	g.id, err = r.Int()
	kiem(err, "1 playerId")
	_, err = r.Int()
	kiem(err, "2 clanId")
	g.cap, err = r.Byte()
	kiem(err, "3 level")
	_, err = r.Byte()
	kiem(err, "4 co co dinh")
	g.typePk, err = r.Byte()
	kiem(err, "5 typePk")
	g.lang, err = r.Byte()
	kiem(err, "6 gender")
	_, err = r.Byte()
	kiem(err, "7 gender lan 2")
	g.head, err = r.Short()
	kiem(err, "8 head")
	g.ten, err = r.UTF()
	kiem(err, "9 name")
	g.hp, err = r.Int()
	kiem(err, "10 hp")
	g.hpToiDa, err = r.Int()
	kiem(err, "11 hpMax")
	g.body, err = r.Short()
	kiem(err, "12 body")
	g.leg, err = r.Short()
	kiem(err, "13 leg")
	g.flagBag, err = r.Byte()
	kiem(err, "14 flagbag")
	_, err = r.Byte()
	kiem(err, "15 hang -1")
	g.x, err = r.Short()
	kiem(err, "16 x")
	g.y, err = r.Short()
	kiem(err, "17 y")
	_, err = r.Short()
	kiem(err, "18 effBuffHp")
	_, err = r.Short()
	kiem(err, "19 effBuffMp")
	_, err = r.Byte()
	kiem(err, "20 numEff")
	_, err = r.Byte()
	kiem(err, "21 idSpaceShip")
	_, err = r.Byte()
	kiem(err, "22 isMonkey")
	_, err = r.Short()
	kiem(err, "23 mount")
	_, err = r.Byte()
	kiem(err, "24 cFlag")
	_, err = r.Byte()
	kiem(err, "25 hang 0")
	g.aura, err = r.Short()
	kiem(err, "26 aura")
	g.effFront, err = r.Byte()
	kiem(err, "27 effFront")
	_, err = r.Short()
	kiem(err, "28 hat")
	if _, err := r.Byte(); err == nil {
		t.Fatal("sau truong 28 (hat) buffer phai het, nhung doc them khong loi")
	}
	return g
}

// bossMau dung mot boss du dung de test goi -5: MOI o mot con so KHAC NHAU de
// nham o nao cung lo ra (outfit 11..66, hp khac hpToiDa, toa do am/duong).
func bossMau(trangThai TrangThaiBoss) *Boss {
	return &Boss{
		Id:         -7,
		PhaHienTai: 0,
		DuLieu: []PhaBoss{{
			Ten:    "Cửu Vĩ",
			Lang:   3,
			Cap:    90,
			Outfit: [6]int16{11, 22, 33, 44, 55, 66}, // head, body, leg, flagBag, aura, effFront
			TypePk: TypePkBoss,
		}},
		Hp:        1234,
		HpToiDa:   5678,
		ViTriX:    -300,
		ViTriY:    777,
		TrangThai: trangThai,
	}
}

// nvMau dung nhan vat toi thieu de vao khu vuc. (100,384) = o spawn hop le cua
// luoi Lat 4 (52x24) — dung y het khuvuc_test.go.
func nvMau(id int64, ten string) *store.NhanVat {
	return &store.NhanVat{ID: id, Ten: ten, HpToiDa: 100, HpHienTai: 100, ViTriX: 100, ViTriY: 384}
}

// TestMoKhuVucChuaCoBoss: khu vuc vua mo KHONG co boss nao — day la dieu giu
// cho moi test cu (dem so goi khi vao khu vuc) khong doi mot con so nao.
func TestMoKhuVucChuaCoBoss(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	if len(k.boss) != 0 {
		t.Fatalf("khu vuc moi mo co %d boss, mong doi 0", len(k.boss))
	}
}

// TestNguoiMoiVaoThayBoss (brief Task 7, test bat buoc 1): khu vuc co 1 boss
// dang HIEN, mot nguoi vao -> phien cua NGUOI DO nhan DUNG goi -5 cua boss,
// doc lai byte khang dinh id AM, ten/cap/hp/toa do/outfit dung, typePk = 5.
func TestNguoiMoiVaoThayBoss(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	b := bossMau(BossHien)
	k.boss = append(k.boss, b)

	p := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvMau(1, "A"), phien: p})

	if len(p.nhan) != 1 {
		t.Fatalf("nguoi moi vao nhan %d goi, mong doi 1 (goi -5 cua boss)", len(p.nhan))
	}
	g := docGoi5(t, p.nhan[0])
	if g.id != -7 {
		t.Errorf("id trong goi -5 = %d, mong doi -7 (id boss PHAI AM)", g.id)
	}
	if g.ten != "Cửu Vĩ" {
		t.Errorf("ten = %q, mong doi %q", g.ten, "Cửu Vĩ")
	}
	if g.cap != 90 {
		t.Errorf("cap = %d, mong doi 90", g.cap)
	}
	if g.hp != 1234 || g.hpToiDa != 5678 {
		t.Errorf("hp/hpToiDa = %d/%d, mong doi 1234/5678", g.hp, g.hpToiDa)
	}
	if g.x != -300 || g.y != 777 {
		t.Errorf("toa do = (%d,%d), mong doi (-300,777)", g.x, g.y)
	}
	if g.typePk != TypePkBoss {
		t.Errorf("typePk = %d, mong doi %d — typePk 0 lam client KHONG ve thanh mau boss", g.typePk, TypePkBoss)
	}
	if g.lang != 3 {
		t.Errorf("lang (gender) = %d, mong doi 3", g.lang)
	}
	if g.head != 11 || g.body != 22 || g.leg != 33 || g.flagBag != 44 || g.aura != 55 || g.effFront != 66 {
		t.Errorf("outfit = head %d, body %d, leg %d, flagBag %d, aura %d, effFront %d; mong doi 11/22/33/44/55/66",
			g.head, g.body, g.leg, g.flagBag, g.aura, g.effFront)
	}
}

// TestNguoiMoiVaoKHONGThayBossDangNghi (test bat buoc 2): boss chua hien
// (NGHI) hoac da di (RA_DI) thi nguoi vao KHONG nhan goi nao cua no.
func TestNguoiMoiVaoKHONGThayBossDangNghi(t *testing.T) {
	for _, tt := range []TrangThaiBoss{BossNghi, BossRaDi} {
		t.Run(tt.String(), func(t *testing.T) {
			k := MoKhuVuc(TimBanDo(1))
			k.boss = append(k.boss, bossMau(tt))

			p := &phienGiaKhuVuc{}
			k.xuLy(suKienVao{nv: nvMau(1, "A"), phien: p})

			if len(p.nhan) != 0 {
				t.Fatalf("boss o trang thai %s: nguoi moi vao nhan %d goi, mong doi 0", tt, len(p.nhan))
			}
		})
	}
}

// TestNguoiMoiVaoThayBossMoiTrangThaiHien: ba trang thai con lai (CHAO,
// HOAT_DONG, TROI) deu la "dang tro len tren man hinh" -> van phai gui -5.
// Thieu test nay thi ai do co the thu hep dieu kien thanh "chi BossHien" ma
// khong test nao do — nguoi vao giua tran se khong thay boss.
func TestNguoiMoiVaoThayBossMoiTrangThaiHien(t *testing.T) {
	for _, tt := range []TrangThaiBoss{BossHien, BossChao, BossHoatDong, BossTroi} {
		t.Run(tt.String(), func(t *testing.T) {
			k := MoKhuVuc(TimBanDo(1))
			k.boss = append(k.boss, bossMau(tt))

			p := &phienGiaKhuVuc{}
			k.xuLy(suKienVao{nv: nvMau(1, "A"), phien: p})

			if len(p.nhan) != 1 {
				t.Fatalf("boss o trang thai %s: nguoi moi vao nhan %d goi, mong doi 1", tt, len(p.nhan))
			}
		})
	}
}

// TestNguoiMoiVaoThayCaNguoiCuVaBoss: boss KHONG duoc lam mat goi -5 cua nguoi
// cu (va nguoc lai) — nguoi moi phai thay CA HAI, va nguoi cu chi thay nguoi
// moi (khong nhan lai goi boss lan hai).
func TestNguoiMoiVaoThayCaNguoiCuVaBoss(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.boss = append(k.boss, bossMau(BossHien))

	pA := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvMau(1, "A"), phien: pA})
	if len(pA.nhan) != 1 {
		t.Fatalf("A vao: nhan %d goi, mong doi 1 (boss)", len(pA.nhan))
	}

	pB := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvMau(2, "B"), phien: pB})
	if len(pB.nhan) != 2 {
		t.Fatalf("B vao: nhan %d goi, mong doi 2 (nguoi cu A + boss)", len(pB.nhan))
	}
	idNhan := map[int32]bool{}
	for _, m := range pB.nhan {
		idNhan[docGoi5(t, m).id] = true
	}
	if !idNhan[1] || !idNhan[-7] {
		t.Fatalf("B phai thay CA nguoi cu (id 1) va boss (id -7), thuc te nhan: %v", idNhan)
	}
	if len(pA.nhan) != 2 {
		t.Fatalf("A sau khi B vao: nhan %d goi, mong doi 2 (boss + B)", len(pA.nhan))
	}
}

// TestThemBossPhatChoMoiNguoi (test bat buoc 3): 2 nguoi dang trong khu vuc,
// themBoss -> CA HAI nhan -5 cua boss, va boss nam trong k.boss.
func TestThemBossPhatChoMoiNguoi(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	pA := &phienGiaKhuVuc{}
	pB := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvMau(1, "A"), phien: pA})
	k.xuLy(suKienVao{nv: nvMau(2, "B"), phien: pB})
	soA, soB := len(pA.nhan), len(pB.nhan)

	b := bossMau(BossHien)
	k.themBoss(b)

	if len(k.boss) != 1 || k.boss[0] != b {
		t.Fatalf("themBoss phai them dung con boss vao k.boss, thuc te len=%d", len(k.boss))
	}
	if len(pA.nhan) != soA+1 {
		t.Fatalf("A nhan %d goi sau themBoss, mong doi %d", len(pA.nhan), soA+1)
	}
	if len(pB.nhan) != soB+1 {
		t.Fatalf("B nhan %d goi sau themBoss, mong doi %d", len(pB.nhan), soB+1)
	}
	if g := docGoi5(t, pA.nhan[soA]); g.id != -7 {
		t.Errorf("A: id trong goi -5 = %d, mong doi -7", g.id)
	}
	if g := docGoi5(t, pB.nhan[soB]); g.id != -7 {
		t.Errorf("B: id trong goi -5 = %d, mong doi -7", g.id)
	}
}

// TestThemBossDangNghiKhongPhat: themBoss cho mot boss CHUA hien (NGHI) van
// them vao danh sach nhung KHONG phat goi nao — bo dieu phoi lich (Task 13) co
// the nap truoc con boss dang dem nguoc, khong duoc lam no hien som.
func TestThemBossDangNghiKhongPhat(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	p := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvMau(1, "A"), phien: p})
	so := len(p.nhan)

	k.themBoss(bossMau(BossNghi))

	if len(k.boss) != 1 {
		t.Fatalf("themBoss phai them boss vao k.boss du dang NGHI, len=%d", len(k.boss))
	}
	if len(p.nhan) != so {
		t.Fatalf("boss NGHI: A nhan %d goi, mong doi giu nguyen %d", len(p.nhan), so)
	}
}

// TestBossPhaNilKhongPanic (test bat buoc 4): du lieu boss hong (khong co pha
// nao, chi so pha ngoai bien, ten dai qua muc writeUTF chiu duoc) KHONG duoc
// lam panic actor va KHONG duoc gui goi hong — ca duong suKienVao lan duong
// themBoss.
func TestBossPhaNilKhongPanic(t *testing.T) {
	bossHong := map[string]*Boss{
		"DuLieu rong": {
			Id:        -1,
			TrangThai: BossHien,
		},
		"PhaHienTai ngoai bien tren": {
			Id:         -2,
			PhaHienTai: 5,
			DuLieu:     []PhaBoss{{Ten: "x", Cap: 1, TypePk: TypePkBoss}},
			TrangThai:  BossHien,
		},
		"PhaHienTai am": {
			Id:         -3,
			PhaHienTai: -1,
			DuLieu:     []PhaBoss{{Ten: "x", Cap: 1, TypePk: TypePkBoss}},
			TrangThai:  BossHien,
		},
		"ten dai qua muc writeUTF": {
			Id:        -4,
			DuLieu:    []PhaBoss{{Ten: strings.Repeat("a", 70000), Cap: 1, TypePk: TypePkBoss}},
			TrangThai: BossHien,
		},
	}

	for ten, b := range bossHong {
		t.Run("suKienVao/"+ten, func(t *testing.T) {
			k := MoKhuVuc(TimBanDo(1))
			k.boss = append(k.boss, b)
			p := &phienGiaKhuVuc{}
			k.xuLy(suKienVao{nv: nvMau(1, "A"), phien: p}) // khong duoc panic
			if len(p.nhan) != 0 {
				t.Fatalf("boss hong (%s): nguoi moi vao nhan %d goi, mong doi 0", ten, len(p.nhan))
			}
			if _, co := k.nguoiTrong[1]; !co {
				t.Fatal("boss hong khong duoc chan nguoi choi vao khu vuc")
			}
		})
		t.Run("themBoss/"+ten, func(t *testing.T) {
			k := MoKhuVuc(TimBanDo(1))
			p := &phienGiaKhuVuc{}
			k.xuLy(suKienVao{nv: nvMau(1, "A"), phien: p})
			so := len(p.nhan)
			k.themBoss(b) // khong duoc panic
			if len(p.nhan) != so {
				t.Fatalf("boss hong (%s): A nhan %d goi, mong doi giu nguyen %d", ten, len(p.nhan), so)
			}
		})
	}
}

// TestThemBossNilKhongPanic: themBoss(nil) la loi lap trinh o noi goi, nhung
// mot con tro nil lot vao KHONG duoc lam sap ca khu vuc.
func TestThemBossNilKhongPanic(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	p := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvMau(1, "A"), phien: p})
	so := len(p.nhan)

	k.themBoss(nil)

	if len(k.boss) != 0 {
		t.Fatalf("themBoss(nil) khong duoc them gi vao k.boss, len=%d", len(k.boss))
	}
	if len(p.nhan) != so {
		t.Fatalf("themBoss(nil): A nhan %d goi, mong doi giu nguyen %d", len(p.nhan), so)
	}
}

// TestBossThieuTypePkVanDung5: du lieu pha quen dien TypePk (=0) thi goi -5 van
// phai mang 5 — de 0 thi nguoi choi danh boss ma khong thay thanh mau (chu
// thich PhaBoss.TypePk, boss.go).
func TestBossThieuTypePkVanDung5(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	b := bossMau(BossHien)
	b.DuLieu[0].TypePk = 0
	k.boss = append(k.boss, b)

	p := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvMau(1, "A"), phien: p})

	if len(p.nhan) != 1 {
		t.Fatalf("nguoi moi vao nhan %d goi, mong doi 1", len(p.nhan))
	}
	if g := docGoi5(t, p.nhan[0]); g.typePk != TypePkBoss {
		t.Errorf("typePk = %d, mong doi %d (thay the khi du lieu de 0)", g.typePk, TypePkBoss)
	}
}

// TestGoiXuatHienBossNilTraNil: ham dung goi khong duoc panic voi con tro nil.
func TestGoiXuatHienBossNilTraNil(t *testing.T) {
	if m := goiXuatHienBoss(nil); m != nil {
		t.Fatal("goiXuatHienBoss(nil) phai tra nil")
	}
}

// TestNguoiMoiVaoPhienHongKhongPanic: phien nguoi moi hong ngay tu goi boss
// dau tien -> dong phien, KHONG them vao khu vuc (dung xu ly nhu goi cua nguoi
// cu that bai, xem case suKienVao).
func TestNguoiMoiVaoPhienHongKhongPanic(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.boss = append(k.boss, bossMau(BossHien))

	p := &phienGiaKhuVuc{}
	p.batDauLoiXep()
	k.xuLy(suKienVao{nv: nvMau(1, "A"), phien: p})

	if _, co := k.nguoiTrong[1]; co {
		t.Fatal("phien hong ngay goi dau: khong duoc them nguoi vao khu vuc")
	}
	if p.soLanDong.Load() == 0 {
		t.Fatal("phien hong phai bi DongCham")
	}
}
