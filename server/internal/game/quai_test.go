package game

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// haiQuaiMau la 2 con quai mau dung chung cho cac test byte-chinh-xac trong
// file nay — gia tri chon co chu y KHAC NHAU o moi truong (level != hpMax !=
// hp != x != y) de mot phep hoan doi truong nao cung lam lech mang byte, thay
// vi tinh co bang gia tri trung nhau che mat loi.
func haiQuaiMau() []*Quai {
	return []*Quai{
		{Id: 0, MauId: 5, Cap: 3, CapTinhAnh: 3, Hp: 80, HpToiDa: 100, X: 120, Y: 200, TrangThai: 0},
		{Id: 1, MauId: 7, Cap: 10, CapTinhAnh: 10, Hp: 500, HpToiDa: 500, X: 300, Y: 150, TrangThai: 0},
	}
}

// TestTaoQuaiMacDinhDungCatalogKhongSuyDienSpawn giu NGUYEN VEN cac con so cua
// ban truoc KH8 task 8 (khi chung con nam trong mot bang cung o quai.go). Task 8
// chi CHUYEN CHO O cua so lieu sang kho-noi-dung/quai.yaml — mot con so nao lech
// la task da lam sai viec, khong phai "gia tri moi".
//
// Them SatThuongCoBan vao phep so (ban cu khong so): do la chinh con so ma task
// 8 lay ra khoi hang so 5 go cung trong quai.go.
func TestTaoQuaiMacDinhDungCatalogKhongSuyDienSpawn(t *testing.T) {
	quai := taoQuaiTuAnhChup(khoChieuMacDinh(), TimBanDo(1))
	// 05/09: y 384 -> 432 (mat co hang 18; 384 la giua troi, client khong cho quai
	// roi nen soi lo lung 2 o) — doi co chu dich, xem quai.yaml.
	if len(quai) != 2 || quai[0].MauId != 0 || quai[0].Hp != 100 || quai[0].HpToiDa != 100 || quai[0].X != 160 || quai[0].Y != 432 || quai[1].X != 40 || quai[1].Y != 432 {
		t.Fatalf("quai mac dinh phai theo kho noi dung, got=%+v", quai)
	}
	for i, q := range quai {
		if q.Id != byte(i) || q.Cap != 1 || q.CapTinhAnh != 1 || q.SatThuongCoBan != 5 {
			t.Fatalf("quai[%d] lech so cu: %+v", i, q)
		}
	}
}

// TestGhiDanhSachQuaiByteDung: so TOAN BO mang byte ma ghiDanhSachQuai ghi ra
// cho 2 con quai voi gia tri ghim cung, dung dac ta §3.6.
func TestGhiDanhSachQuaiByteDung(t *testing.T) {
	w := protocol.NewWriter()
	ghiDanhSachQuai(w, haiQuaiMau())

	muon := []byte{
		2, // soQuai

		// --- quai 1: Id=0, MauId=5, Cap=3, Hp=80, HpToiDa=100, X=120, Y=200 ---
		0, 0, 0, 0, 0, // isDisable/isDontMove/isFire/isIce/isWind = false
		5,           // o 6 = MauId (templateId phia client) — KHONG phai Id (05/09, xem ghiQuai)
		0,           // sys = 0
		0, 0, 0, 80, // hp = 80 (int32, TRUOC level)
		3,            // level = 3 (GIUA hp va hpMax)
		0, 0, 0, 100, // hpMax = 100 (int32, SAU level)
		0, 120, // x = 120
		0, 200, // y = 200 (200 < 256, 1 byte du dung 2 byte short)
		0, // status = 0
		3, // lvMob = 3
		0, // isBigBoss = false

		// --- quai 2: Id=1, MauId=7, Cap=10, Hp=500, HpToiDa=500, X=300, Y=150 ---
		0, 0, 0, 0, 0,
		7,            // o 6 = MauId = 7 (client: Mob.arrMobTemplate[7]); Id=1 chi la vi tri trong danh sach
		0,            // sys = 0
		0, 0, 1, 244, // hp = 500
		10,           // level = 10
		0, 0, 1, 244, // hpMax = 500
		1, 44, // x = 300 (0x012C)
		0, 150, // y = 150
		0,  // status = 0
		10, // lvMob = 10
		0,  // isBigBoss = false
	}

	got := w.Bytes()
	if !bytes.Equal(got, muon) {
		t.Fatalf("ghiDanhSachQuai byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			got, len(got), muon, len(muon))
	}
}

// TestGhiQuaiThuTuHpCapHpToiDa khang dinh RIENG thu tu doc duoc cua 3 truong
// hp -> level -> hpMax (dac ta §3.6) — day la cho DE CHEP SAI NHAT (truc giac
// hay nghi hp -> hpMax -> level). Test doc TUNG truong theo dung kieu cua no
// (int, byte, int) ngay sau tempId/sys, dam bao thu tu that su nhu the.
func TestGhiQuaiThuTuHpCapHpToiDa(t *testing.T) {
	q := &Quai{Id: 9, Cap: 42, Hp: 111, HpToiDa: 222, X: 1, Y: 2, TrangThai: 0, CapTinhAnh: 42}

	w := protocol.NewWriter()
	ghiQuai(w, q)
	r := protocol.NewReader(w.Bytes())

	// 5 co bool + templateId (MauId) + sys.
	for i := 0; i < 5; i++ {
		if _, err := r.Byte(); err != nil {
			t.Fatalf("doc co bool thu %d loi: %v", i+1, err)
		}
	}
	if _, err := r.Byte(); err != nil { // templateId = MauId
		t.Fatalf("doc templateId loi: %v", err)
	}
	if _, err := r.Byte(); err != nil { // sys
		t.Fatalf("doc sys loi: %v", err)
	}

	hp, err := r.Int()
	if err != nil || hp != 111 {
		t.Fatalf("truong thu 8 (hp) = %d (loi %v), mong doi 111 — phai la hp, DUNG truoc level", hp, err)
	}
	level, err := r.Byte()
	if err != nil || level != 42 {
		t.Fatalf("truong thu 9 (level) = %d (loi %v), mong doi 42 — phai la level, GIUA hp va hpMax", level, err)
	}
	hpMax, err := r.Int()
	if err != nil || hpMax != 222 {
		t.Fatalf("truong thu 10 (hpMax) = %d (loi %v), mong doi 222 — phai la hpMax, SAU level", hpMax, err)
	}
}

// TestGoiMapInfoVoiQuaiByteDung dung goi -24 that (qua goiMapInfo) voi 2 con
// quai, so TOAN BO than goi voi gia tri ghim cung — dam bao mobCount, danh
// sach quai, va byte du phong (a) sau quai deu dung vi tri, KHONG bi nuot mat.
func TestGoiMapInfoVoiQuaiByteDung(t *testing.T) {
	nv := &store.NhanVat{MapID: 1, ViTriX: 100, ViTriY: 384}
	bd := TimBanDo(1)

	m, err := goiMapInfo(nv, bd, haiQuaiMau())
	if err != nil {
		t.Fatalf("goiMapInfo loi: %v", err)
	}
	r := protocol.NewReader(m.Data)

	b := func(ten string, muon byte) {
		t.Helper()
		g, err := r.Byte()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}
	sh := func(ten string, muon int16) {
		t.Helper()
		g, err := r.Short()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}
	in := func(ten string, muon int32) {
		t.Helper()
		g, err := r.Int()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}

	b("mapId", 1)
	b("planetId", 0)
	b("tileId", 1)
	b("bgId", 0)
	b("type", 0)
	if ten, err := r.UTF(); err != nil || ten != "Làng Lá" {
		t.Fatalf("mapName = %q (loi %v)", ten, err)
	}
	b("zoneId", 0)
	sh("x", 100)
	sh("y", 384)
	// Map 1 co DUNG MOT cong roi map (Cong Lang La). Byte-chinh-xac cua o nay
	// thuoc ve TestGoiMapInfoMap1CoCongRoiMap (diemdi_test.go); o day chi doc
	// qua cho dung so byte, vi test nay soi phan QUAI phia sau.
	b("waypointCount", 1)
	sh("wp.minX", 408)
	sh("wp.minY", 384)
	sh("wp.maxX", 472)
	sh("wp.maxY", 432)
	b("wp.laCong", 1)
	b("wp.ngoaiTuyen", 0)
	if ten, err := r.UTF(); err != nil || ten != "Cổng Làng Lá" {
		t.Fatalf("wp.ten = %q (loi %v)", ten, err)
	}

	b("mobCount", 2)

	// --- quai 1 ---
	b("q1.isDisable", 0)
	b("q1.isDontMove", 0)
	b("q1.isFire", 0)
	b("q1.isIce", 0)
	b("q1.isWind", 0)
	b("q1.templateId(MauId)", 5) // 05/09: o 6 la MAU quai, khong phai Id — xem ghiQuai
	b("q1.sys", 0)
	in("q1.hp", 80)
	b("q1.level", 3)
	in("q1.hpMax", 100)
	sh("q1.x", 120)
	sh("q1.y", 200)
	b("q1.status", 0)
	b("q1.lvMob", 3)
	b("q1.isBigBoss", 0)

	// --- quai 2 ---
	b("q2.isDisable", 0)
	b("q2.isDontMove", 0)
	b("q2.isFire", 0)
	b("q2.isIce", 0)
	b("q2.isWind", 0)
	b("q2.templateId(MauId)", 7)
	b("q2.sys", 0)
	in("q2.hp", 500)
	b("q2.level", 10)
	in("q2.hpMax", 500)
	sh("q2.x", 300)
	sh("q2.y", 150)
	b("q2.status", 0)
	b("q2.lvMob", 10)
	b("q2.isBigBoss", 0)

	b("byte du phong", 0) // (a) — giua quai va npc, DUNG QUEN
	b("npcCount", 1)
	b("npc.status", 0)
	sh("npc.x", 172)
	sh("npc.y", 429) // 05/09: +3 phia client = 432 mat co (truoc 384 lo lung)
	b("npc.template", 0)
	sh("npc.portrait", 27078)
	b("itemCount", 0)
	// Map 1 co 11 bong do; byte-chinh-xac thuoc ve TestGoiMapInfoMap1CoBongDo
	// (bgitem_test.go). O day chi doc qua cho dung so byte.
	sh("soBong", 11)
	for i := 0; i < 11; i++ {
		if _, err := r.Short(); err != nil {
			t.Fatalf("bong do #%d thieu truong mau: %v", i, err)
		}
		if _, err := r.Short(); err != nil {
			t.Fatalf("bong do #%d thieu truong tileX: %v", i, err)
		}
		if _, err := r.Short(); err != nil {
			t.Fatalf("bong do #%d thieu truong tileY: %v", i, err)
		}
	}
	sh("effMap", 0)
	b("bgType", 0)
	b("idSpaceShip", 0)
	b("isMapDouble", 0)

	if _, err := r.Byte(); err == nil {
		t.Error("MAP INFO con byte thua sau isMapDouble — phai het buffer")
	}
}

// TestQuaiSoLuongVuot127BiCat: soQuai la byte co dau phia client (dac ta §1.4),
// nen danh sach qua 127 con PHAI bi cat con 127, khong duoc ghi so am.
func TestQuaiSoLuongVuot127BiCat(t *testing.T) {
	ds := make([]*Quai, 130)
	for i := range ds {
		ds[i] = &Quai{Id: byte(i % 128)}
	}
	w := protocol.NewWriter()
	ghiDanhSachQuai(w, ds)
	r := protocol.NewReader(w.Bytes())
	n, err := r.Byte()
	if err != nil || n != 127 {
		t.Fatalf("soQuai = %d (loi %v), mong doi cat con 127", n, err)
	}
}

// TestDanhSachQuaiSnapshotQuaHopThuActor kham nghiem VONG SUA 1 (theo review
// coordinator): duong lay quai cho goi -24 PHAI di qua hop thu actor
// (suKienXinDanhSachQuai), KHONG doc truc tiep truong k.quai tu luong khac
// (dac ta §1.6). Vi test nam CUNG package (game), co the so sanh dia chi con
// tro giua ban sao tra ve va phan tu goc trong k.quai — chung minh
// DanhSachQuaiSnapshot tra CON TRO MOI (sao gia tri), khong phai con tro goc.
func TestDanhSachQuaiSnapshotQuaHopThuActor(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.BatDau()
	defer k.Dung()

	snap := k.DanhSachQuaiSnapshot()

	if len(snap) != len(k.quai) {
		t.Fatalf("snapshot co %d quai, mong doi %d (bang so quai noi bo cua khu vuc)",
			len(snap), len(k.quai))
	}
	for i, q := range snap {
		goc := k.quai[i]
		if q == goc {
			t.Fatalf("quai thu %d: snapshot tra VE CUNG con tro voi k.quai — phai la BAN SAO GIA TRI (dac ta §1.6), khong duoc lo con tro goc ra ngoai vongLap", i)
		}
		// reflect.DeepEqual (khong phai *q != *goc): tu Task 9, Quai co them
		// truong DaTungDanh kieu map[int64]bool — struct chua map KHONG so
		// sanh duoc bang != o Go (bao loi bien dich), nen phai doi sang
		// DeepEqual. Van giu dung y dinh test goc: gia tri snapshot phai
		// KHOP HOAN TOAN voi gia tri goc.
		if !reflect.DeepEqual(*q, *goc) {
			t.Fatalf("quai thu %d: gia tri snapshot khac gia tri goc: %+v != %+v", i, *q, *goc)
		}
	}
}

// TestDanhSachQuaiSnapshotActorChuaBatDauTraNilNgay: khu vuc chua BatDau() thi
// chua co hop thu (k.hopThu == nil) — DanhSachQuaiSnapshot PHAI tra ve nil
// NGAY, khong duoc cho het 100ms timeout (luong phien khong bi treo vo ich).
func TestDanhSachQuaiSnapshotActorChuaBatDauTraNilNgay(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1)) // KHONG goi BatDau()

	batDau := time.Now()
	snap := k.DanhSachQuaiSnapshot()
	trongVongMs := time.Since(batDau)

	if snap != nil {
		t.Fatalf("actor chua BatDau(), snapshot phai nil, duoc %v", snap)
	}
	if trongVongMs > 20*time.Millisecond {
		t.Fatalf("DanhSachQuaiSnapshot mat %v de tra ve khi actor chua BatDau — phai tra NGAY (khong hop thu de gui), khong duoc cho het 100ms timeout", trongVongMs)
	}
}

// ---------------------------------------------------------------------------
// Ke hoach 6 Task 8 — quai hoi sinh 5 giay + hoi mau 30 giay, trong nhipTim.
//
// DONG HO TIEM DUOC (rang buoc bat buoc cua brief): k.dongHo la mot closure
// doc gia tri HIEN TAI cua con tro *time.Time — test "di chuyen" dong ho bang
// cach doi gia tri con tro do, KHONG sleep that. Tat ca cac test duoi day
// dung MOT khu vuc CHUA goi BatDau() (khong co goroutine vongLap song song),
// giong ky thuat da dung o Task 7 (xem quai_test.go/khuvuc_tichhop_test.go):
// chi mot goroutine (goroutine cua chinh ham Test) cham vao k.quai/k.nguoiTrong
// tai moi thoi diem nen an toan, khong can dong bo them.
// ---------------------------------------------------------------------------

// khuVucNhipTimGia dung mot KhuVuc CHUA BatDau(), gan dongHo tiem duoc tro
// toi bien luc (con tro dung chung voi ham goi, de test doi "bay gio" giua
// cac lan goi nhipTim), va them SAN mot nguoi choi voi phien gia
// (phienGiaKhuVuc, dinh nghia o khuvuc_test.go, cung package) de co the
// khang dinh goi tin ma phatChoCaKhuVuc phat ra tu nhipTim.
func khuVucNhipTimGia(luc *time.Time) (*KhuVuc, *phienGiaKhuVuc) {
	k := MoKhuVuc(TimBanDo(1))
	k.dongHo = func() time.Time { return *luc }
	p := &phienGiaKhuVuc{}
	k.nguoiTrong[1] = &nguoiTrongKhuVuc{nv: &store.NhanVat{ID: 1, Ten: "NguoiXem"}, phien: p}
	return k, p
}

// TestNhipTimQuaiChuaHoiSinhTruoc5Giay: quai chet luc t0, kiem tai t0+4.9s —
// brief Task 8 buoc 1 doi hoi CHUA hoi sinh, KHONG phat goi gi (day la vach
// bien duoi cua nguong 5 giay, dac ta §5.2).
func TestNhipTimQuaiChuaHoiSinhTruoc5Giay(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucNhipTimGia(&luc)
	k.quai = []*Quai{{Id: 0, Hp: 0, HpToiDa: 100, LucChet: t0}}

	luc = t0.Add(4900 * time.Millisecond) // 4.9s — CHUA du 5s
	k.nhipTim()

	if k.quai[0].Hp != 0 {
		t.Fatalf("Hp luc 4.9s sau khi chet = %d, mong doi VAN LA 0 (chua du 5s de hoi sinh)", k.quai[0].Hp)
	}
	if len(p.nhan) != 0 {
		t.Fatalf("nhan %d goi luc 4.9s (chua hoi sinh), mong doi 0 — khong duoc phat goi gi truoc khi du nguong", len(p.nhan))
	}
}

// TestNhipTimQuaiHoiSinhDung5Giay: quai chet luc t0, kiem tai DUNG t0+5.0s —
// dac ta §5.2 + brief Task 8: "dung/qua 5.0s -> hoi sinh, phat -13 cho ca khu
// vuc". Kiem byte-chinh-xac goi -13 (mobId, sys=0, levelBoss=0, hp=HpToiDa)
// va hai truong trang thai noi bo (Hp day lai, LucHoiSinh duoc ghi).
func TestNhipTimQuaiHoiSinhDung5Giay(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucNhipTimGia(&luc)
	k.quai = []*Quai{{Id: 3, Hp: 0, HpToiDa: 250, LucChet: t0}}

	luc = t0.Add(5 * time.Second) // DUNG 5.0s
	k.nhipTim()

	q := k.quai[0]
	if q.Hp != 250 {
		t.Fatalf("Hp sau hoi sinh = %d, mong doi day lai HpToiDa = 250", q.Hp)
	}
	if !q.LucHoiSinh.Equal(luc) {
		t.Fatalf("LucHoiSinh = %v, mong doi = thoi diem hoi sinh %v", q.LucHoiSinh, luc)
	}

	if len(p.nhan) != 1 {
		t.Fatalf("nhan %d goi luc hoi sinh, mong doi DUNG 1 goi (-13)", len(p.nhan))
	}
	m := p.nhan[0]
	if m.Command != -13 {
		t.Fatalf("Command = %d, mong doi -13", m.Command)
	}
	muon := []byte{
		3,            // mobId
		0,            // sys = 0 (dac ta §3.3: khong doi loai quai khi hoi sinh)
		0,            // levelBoss = 0 (dac ta §3.3: khac 0 se bat hieu ung sieu cap ngoai y)
		0, 0, 0, 250, // hp = 250 = HpToiDa
	}
	if !bytes.Equal(m.Data, muon) {
		t.Fatalf("goi -13 byte sai:\n got  = %v (%d byte)\n muon = %v (%d byte)",
			m.Data, len(m.Data), muon, len(muon))
	}
}

// TestNhipTimQuaiHoiSinhKhongPhatGoiKhacSauDo la TEST BAT BUOC cua brief Task
// 8 va dac ta §0.1: sau goi -13, TUYET DOI khong duoc phat them goi nao khac,
// DAC BIET KHONG goi 87 (ban Java phat 87 ngay sau -13 tai Mob.java:525, rat
// de chep theo mot cach vo tinh — nhung ma 87 phia client la "quai bi gio
// thoi" (byte mobId + bool), khong phai maxHp; voi quai mau >= 16.777.216,
// byte cao nhat cua maxHp se bi client hieu nham thanh isWind=true, lam quai
// DUNG IM khong di chuyen nua). HpToiDa duoc chon RAT LON (co dinh dac diem
// do) de neu ai sau nay vo tinh them lai gioi han goi 87, gia tri nay chinh
// la truong hop se lo hau qua that.
func TestNhipTimQuaiHoiSinhKhongPhatGoiKhacSauDo(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucNhipTimGia(&luc)
	k.quai = []*Quai{{Id: 0, Hp: 0, HpToiDa: 20_000_000, LucChet: t0}}

	luc = t0.Add(5 * time.Second)
	k.nhipTim()

	if len(p.nhan) != 1 {
		t.Fatalf("sau khi hoi sinh, nhan %d goi — dac ta §0.1 doi hoi DUNG 1 goi (-13), KHONG mot goi nao khac", len(p.nhan))
	}
	if p.nhan[0].Command == 87 {
		t.Fatal("goi 87 BI PHAT sau -13 — TUYET DOI cam theo dac ta §0.1: day la lenh 'quai bi gio thoi' phia client, khong phai maxHp")
	}
	if p.nhan[0].Command != -13 {
		t.Fatalf("Command = %d, mong doi -13", p.nhan[0].Command)
	}

	// Tick ke tiep (quai da day mau, chua co su kien nao khac xay ra) khong
	// duoc phat THEM goi nao nua.
	luc = luc.Add(100 * time.Millisecond)
	k.nhipTim()
	if len(p.nhan) != 1 {
		t.Fatalf("tick ke tiep sau khi hoi sinh: tong nhan %d goi, mong doi VAN LA 1 (khong phat gi them)", len(p.nhan))
	}
}

// TestNhipTimQuaiHoiMauChuaDu30GiayKhongHoi: quai mat mau, kiem tai t0+29.9s
// — CHUA du 30 giay nen KHONG duoc hoi (dac ta §5.2, vach bien duoi).
func TestNhipTimQuaiHoiMauChuaDu30GiayKhongHoi(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	luc := t0
	k, _ := khuVucNhipTimGia(&luc)
	k.quai = []*Quai{{Id: 0, Hp: 50, HpToiDa: 100, LanHoiMauCuoi: t0}}

	luc = t0.Add(29900 * time.Millisecond) // 29.9s
	k.nhipTim()

	if k.quai[0].Hp != 50 {
		t.Fatalf("Hp luc 29.9s = %d, mong doi VAN LA 50 (chua du 30s de hoi mau)", k.quai[0].Hp)
	}
}

// TestNhipTimQuaiHoiMauDung30GiayHoiDungHpToiDaChia10: kiem tai DUNG t0+30.0s
// — dac ta §5.2: hoi dung hpToiDa/10, va moc LanHoiMauCuoi phai duoc dat lai
// (de chu ky tiep theo tinh tu day, khong phai tu diem cu).
func TestNhipTimQuaiHoiMauDung30GiayHoiDungHpToiDaChia10(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	luc := t0
	k, _ := khuVucNhipTimGia(&luc)
	k.quai = []*Quai{{Id: 0, Hp: 50, HpToiDa: 100, LanHoiMauCuoi: t0}}

	luc = t0.Add(30 * time.Second) // DUNG 30.0s
	k.nhipTim()

	q := k.quai[0]
	if q.Hp != 60 { // 50 + 100/10 = 60
		t.Fatalf("Hp sau hoi mau = %d, mong doi 60 (50 + hpToiDa/10)", q.Hp)
	}
	if !q.LanHoiMauCuoi.Equal(luc) {
		t.Fatalf("LanHoiMauCuoi = %v, mong doi duoc dat lai ve %v", q.LanHoiMauCuoi, luc)
	}
}

// TestNhipTimQuaiHoiMauKhongVuotHpToiDa: Hp=95, HpToiDa=100 — hoi 10 se ra
// 105, PHAI bi kep con 100 (dac ta §5.2 "khong vuot hpToiDa").
func TestNhipTimQuaiHoiMauKhongVuotHpToiDa(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	luc := t0
	k, _ := khuVucNhipTimGia(&luc)
	k.quai = []*Quai{{Id: 0, Hp: 95, HpToiDa: 100, LanHoiMauCuoi: t0}}

	luc = t0.Add(30 * time.Second)
	k.nhipTim()

	if k.quai[0].Hp != 100 {
		t.Fatalf("Hp sau hoi mau = %d, mong doi kep con 100 (khong vuot hpToiDa)", k.quai[0].Hp)
	}
}

// TestNhipTimQuaiDayMauKhongHoiMauThuaKhongPhatGoi: quai da day mau — nhanh
// hoi mau khong duoc chay du da qua 30s (khong hoi thua ra ngoai hpToiDa), va
// khong phat goi tin nao (dac ta khong dinh nghia goi nao rieng cho hoi mau
// dinh ky, khac han goi -9 khi bi danh).
func TestNhipTimQuaiDayMauKhongHoiMauThuaKhongPhatGoi(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucNhipTimGia(&luc)
	k.quai = []*Quai{{Id: 0, Hp: 100, HpToiDa: 100, LanHoiMauCuoi: t0}}

	luc = t0.Add(60 * time.Second)
	k.nhipTim()

	if k.quai[0].Hp != 100 {
		t.Fatalf("Hp = %d, mong doi van la 100", k.quai[0].Hp)
	}
	if len(p.nhan) != 0 {
		t.Fatalf("nhan %d goi tu nhanh hoi mau, mong doi 0 (hoi mau dinh ky KHONG phat goi)", len(p.nhan))
	}
}

// TestNhipTimHaiQuaiDocLapVoiNhau: mot khu vuc co 2 con quai o hai trang thai
// khac nhau (mot dang cho hoi sinh CHUA du 5s, mot da day mau) trong CUNG mot
// tick — moi con phai duoc xu ly DOC LAP, khong anh huong lan nhau (dac ta
// §1.6 ngu y actor duyet toan bo danh sach quai moi nhip, khong dung som).
func TestNhipTimHaiQuaiDocLapVoiNhau(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucNhipTimGia(&luc)
	k.quai = []*Quai{
		{Id: 0, Hp: 0, HpToiDa: 100, LucChet: t0},        // moi chet, se hoi sinh dung luc
		{Id: 1, Hp: 80, HpToiDa: 100, LanHoiMauCuoi: t0}, // con song, chua du 30s hoi mau
	}

	luc = t0.Add(5 * time.Second) // du 5s cho quai 0, nhung moi 5s cho quai 1 (chua du 30s)
	k.nhipTim()

	if k.quai[0].Hp != 100 {
		t.Fatalf("quai 0: Hp = %d, mong doi hoi sinh day mau 100", k.quai[0].Hp)
	}
	if k.quai[1].Hp != 80 {
		t.Fatalf("quai 1: Hp = %d, mong doi VAN LA 80 (chua du 30s hoi mau)", k.quai[1].Hp)
	}
	if len(p.nhan) != 1 {
		t.Fatalf("nhan %d goi, mong doi DUNG 1 (chi -13 cua quai 0)", len(p.nhan))
	}
	if p.nhan[0].Command != -13 {
		t.Fatalf("Command = %d, mong doi -13", p.nhan[0].Command)
	}
}

// TestSuKienDanhQuaiBatTuTrongKhoangSauHoiSinh: kiem tra day du day noi day —
// Task 8 them mucTieu.LucHoiSinh vao ViTriCoMau (khuvuc.go, case
// suKienDanhQuai) thay vi de trong nhu truoc day (xem chu thich cu trich
// nguyen van trong bao cao Task 7: "se them o Task 8"). Quai vua hoi sinh
// 500ms truoc (< 1500ms bat tu, dac ta §1.7) PHAI bi tu choi danh — im lang,
// khong phat goi gi — CHUNG MINH bang mot chieu HOP LE (tam, chakra, hoi
// chieu deu qua) de loai tru kha nang bi chan boi ly do khac ngoai bat tu.
func TestSuKienDanhQuaiBatTuTrongKhoangSauHoiSinh(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	luc := t0
	k := MoKhuVuc(TimBanDo(1))
	k.dongHo = func() time.Time { return luc }

	// HpHienTai: 100 (con song) — Ke hoach 6 Task 13 them chot "nguoi da chet
	// (HpHienTai<=0) khong danh duoc" vao case suKienDanhQuai (khuvuc.go); bo
	// truong nay se mac dinh HpHienTai=0 (zero value) va bi CHOT DO chan
	// TRUOC ca kiem dang test o day, lam sai lech y nghia test.
	nvA := &store.NhanVat{ID: 1, Ten: "A", Lang: 0, ChakraHienTai: 100, HpHienTai: 100}
	pA := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvA, phien: pA})
	// templateId=1 = slot can chien cua lang 0 (chieuthuc.go): SkillId=1,
	// ChakraDung=10, Dx=Dy=tamCanChien=100 — du de nvA (0,0) danh trung quai
	// cung o (0,0) neu KHONG bi chan boi bat tu.
	k.xuLy(suKienChonChieu{idNhanVat: nvA.ID, mauChieuId: 1})
	pA.nhan = nil // don goi (neu co) tu hai buoc chuan bi tren truoc khi kiem

	k.quai = []*Quai{{Id: 0, Hp: 100, HpToiDa: 100, LucHoiSinh: t0}}
	luc = t0.Add(500 * time.Millisecond) // 500ms < 1500ms bat tu

	k.xuLy(suKienDanhQuai{idNhanVat: nvA.ID, mobId: 0})

	if len(pA.nhan) != 0 {
		t.Fatalf("con trong 1500ms bat tu sau hoi sinh nhung van phat %d goi — LucHoiSinh chua duoc noi day dung vao ViTriCoMau (case suKienDanhQuai)", len(pA.nhan))
	}
	if k.quai[0].Hp != 100 {
		t.Fatalf("Hp quai = %d, mong doi VAN LA 100 (don danh phai bi tu choi hoan toan boi bat tu)", k.quai[0].Hp)
	}
}

// ---------------------------------------------------------------------------
// Ke hoach 6 Task 9 — AI quai hai tang uu tien (dac ta §5.2) + goi -11/-10.
//
// nguoiGia la mot helper NHO dung rieng cho nhom test chonMucTieu (thuan,
// khong dung actor/hop thu): dung tao nhanh *nguoiTrongKhuVuc voi id, vi tri,
// mau con lai — CHI can 3 truong nay cho chonMucTieu (X/Y de tinh khoang
// cach, HpHienTai de loc "da chet", ID de tra cuu DaTungDanh).
// ---------------------------------------------------------------------------

func nguoiGia(id int64, x, y int32, hpHienTai int32) *nguoiTrongKhuVuc {
	return &nguoiTrongKhuVuc{nv: &store.NhanVat{ID: id, ViTriX: x, ViTriY: y, HpHienTai: hpHienTai}}
}

// bayGioChonMucTieuGia la "bay gio" gia dung CHUNG cho cac test chonMucTieu
// KHONG lien quan toi bat tu (Ke hoach 6 Task 11 Phan B them tham so bayGio
// vao chu ky ham) — nguoiGia() KHONG dat lucHoiSinh (zero value, nam 1) nen
// bat ky moc "bay gio" nao du xa trong tuong lai cung an toan (khong bi loc
// oan boi kiem bat tu moi them). Cac test RIENG kiem CHINH tinh nang bat tu
// (TestChonMucTieu...VuaHoiSinh... o duoi) tu dinh nghia moc gio cua minh.
var bayGioChonMucTieuGia = time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)

// TestChonMucTieuTangMotUuTienNguoiDaTungDanh: quai co MOT nguoi DA TUNG danh
// no (tam 300) VA mot nguoi KHAC gan hon nhieu nhung CHUA TUNG danh (chi hop
// le o tam 2 100) — tang 1 PHAI thang tuyet doi, khong so khoang cach xuyen
// tang (dac ta §5.2: hai tang la HAI CHE DO khac nhau, khong phai mot danh
// sach xep hang chung theo khoang cach).
func TestChonMucTieuTangMotUuTienNguoiDaTungDanh(t *testing.T) {
	q := &Quai{X: 0, Y: 0, DaTungDanh: map[int64]bool{2: true}}

	nguoiA := nguoiGia(1, 50, 0, 100)  // gan (50), NHUNG chua tung danh quai -> chi hop le o tang 2
	nguoiB := nguoiGia(2, 250, 0, 100) // xa hon nhieu (250), NHUNG da tung danh -> hop le tang 1

	mt, nhip := q.chonMucTieu([]*nguoiTrongKhuVuc{nguoiA, nguoiB}, bayGioChonMucTieuGia)

	if mt == nil {
		t.Fatal("chonMucTieu tra nil, mong doi chon nguoiB (tang 1)")
	}
	if mt.nv.ID != 2 {
		t.Fatalf("chon nguoi id=%d, mong doi id=2 (da tung danh, tang 1 uu tien hon tang 2 du xa hon)", mt.nv.ID)
	}
	if nhip != nhipTraThu {
		t.Fatalf("nhip = %v, mong doi %v (nhip tang 1)", nhip, nhipTraThu)
	}
}

// TestChonMucTieuTangHaiBatKyGanNhat: khong ai trong DaTungDanh (tang 1 rong)
// -> roi xuong tang 2, chon nguoi GAN NHAT trong tam 100 giua nhieu ung vien.
func TestChonMucTieuTangHaiBatKyGanNhat(t *testing.T) {
	q := &Quai{X: 0, Y: 0}

	xa := nguoiGia(10, 80, 0, 50)  // trong tam 100, nhung xa hon
	gan := nguoiGia(11, 30, 0, 50) // trong tam 100, gan nhat

	mt, nhip := q.chonMucTieu([]*nguoiTrongKhuVuc{xa, gan}, bayGioChonMucTieuGia)

	if mt == nil {
		t.Fatal("chonMucTieu tra nil, mong doi chon nguoi gan nhat (tang 2)")
	}
	if mt.nv.ID != 11 {
		t.Fatalf("chon nguoi id=%d, mong doi id=11 (gan nhat trong tam 100)", mt.nv.ID)
	}
	if nhip != nhipTuGaySu {
		t.Fatalf("nhip = %v, mong doi %v (nhip tang 2)", nhip, nhipTuGaySu)
	}
}

// TestChonMucTieuMucTieuDaChetKhongDuocChon la TEST BAT BUOC cua brief Task 9
// (dac ta §5.4.3 — bay ban Java sai thu tu &&/|| khien big boss nham ca xac
// chet). Nguoi duy nhat trong tam THOA het cac dieu kien khac cua CA HAI tang
// (da tung danh + trong tam 300 cho tang 1; trong tam 100 cho tang 2) nhung
// HpHienTai <= 0 (da chet) — PHAI bi loai o CA HAI tang, chonMucTieu tra nil.
func TestChonMucTieuMucTieuDaChetKhongDuocChon(t *testing.T) {
	q := &Quai{X: 0, Y: 0, DaTungDanh: map[int64]bool{1: true}}

	nguoiChet := nguoiGia(1, 10, 0, 0) // da tung danh + rat gan (ca hai tam) + HP=0

	mt, nhip := q.chonMucTieu([]*nguoiTrongKhuVuc{nguoiChet}, bayGioChonMucTieuGia)

	if mt != nil {
		t.Fatalf("chonMucTieu chon nguoi id=%d dang CHET (HpHienTai<=0) — vi pham dac ta §5.4.3 (bay big-boss nham xac chet cua ban Java)", mt.nv.ID)
	}
	if nhip != 0 {
		t.Fatalf("nhip = %v, mong doi 0 khi khong co muc tieu", nhip)
	}
}

// TestChonMucTieuKhongAiTraNilVaKhongTieuNhip: danh sach nguoi RONG -> tra
// (nil, 0) — chonMucTieu la HAM THUAN nen KHONG the "tieu nhip" cua bat ky ai
// (khong co truong nao de sua); phan "khong tieu nhip that su" (dac ta §5.4.2)
// duoc chung minh o muc actor bang TestNhipTimAiQuai... ben duoi (chonMucTieu
// chi la mot nua cau chuyen — nua con lai la aiQuai KHONG duoc dung ket qua
// nil nay de dat lai q.LanDanhCuoi).
func TestChonMucTieuKhongAiTraNilVaKhongTieuNhip(t *testing.T) {
	q := &Quai{X: 0, Y: 0}

	mt, nhip := q.chonMucTieu(nil, bayGioChonMucTieuGia)

	if mt != nil {
		t.Fatalf("chonMucTieu tra muc tieu %+v, mong doi nil khi danh sach nguoi rong", mt)
	}
	if nhip != 0 {
		t.Fatalf("nhip = %v, mong doi 0 khi khong co muc tieu", nhip)
	}
}

// ---------------------------------------------------------------------------
// Ke hoach 6 Task 11 Phan B — VA GAP bat tu cho quai (ban giao tu Task 13):
// chonMucTieu (va do do aiQuai/nhipTim, khuvuc.go) truoc day CHI loc
// HpHienTai<=0, KHONG kiem lucHoiSinh — vi pham dac ta §1.7 ("bat tu 1500ms
// sau hoi sinh — MOI don vao muc tieu vua hoi sinh bi tu choi, KHONG gioi han
// huong"): huong nguoi-danh-quai da duoc KiemDanhHopLe (chotkiem.go, Task 3)
// bao ve tu truoc, nhung huong quai-danh-nguoi thi KHONG co gi chan ca. Bon
// test duoi day la TEST BAT BUOC cua brief Task 11 Phan B: "quai KHONG danh
// nguoi vua hoi sinh <1500ms; sau 1500ms thi danh lai duoc" — kiem ca hai
// tang uu tien (dac ta §1.7 khong gioi han huong/tang).
// ---------------------------------------------------------------------------

// TestChonMucTieuTangHaiBoQuaNguoiVuaHoiSinh: nguoi duy nhat trong tam (tang
// 2, "tu gay su") vua hoi sinh 500ms truoc bayGio (< 1500ms bat tu) — PHAI bi
// loai, chonMucTieu tra nil, DU la ung vien HOP LE duy nhat (con song, trong
// tam). Day la TEST DOT BIEN CHINH cua Phan B: bo dieu kien hopLeBatTu o tang
// 2 (chonMucTieu, quai.go) se lam test nay DO.
func TestChonMucTieuTangHaiBoQuaNguoiVuaHoiSinh(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	q := &Quai{X: 0, Y: 0}

	nguoiVuaHoiSinh := nguoiGia(1, 10, 0, 100) // rat gan (tam 2 = 100), con song
	nguoiVuaHoiSinh.lucHoiSinh = t0

	bayGio := t0.Add(500 * time.Millisecond) // 500ms < 1500ms bat tu

	mt, nhip := q.chonMucTieu([]*nguoiTrongKhuVuc{nguoiVuaHoiSinh}, bayGio)

	if mt != nil {
		t.Fatalf("chonMucTieu chon nguoi id=%d dang BAT TU (vua hoi sinh 500ms truoc) — vi pham dac ta §1.7", mt.nv.ID)
	}
	if nhip != 0 {
		t.Fatalf("nhip = %v, mong doi 0 khi khong co muc tieu", nhip)
	}
}

// TestChonMucTieuTangHaiDanhDuocSauKhiHetBatTu: nguoi CUNG nguoi cua test
// tren, nhung bayGio gio la DUNG 1500ms sau lucHoiSinh (bien tren cua nguong,
// giong quy uoc ">= " cua KiemDanhHopLe) — PHAI duoc chon lai binh thuong.
func TestChonMucTieuTangHaiDanhDuocSauKhiHetBatTu(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	q := &Quai{X: 0, Y: 0}

	nguoi := nguoiGia(1, 10, 0, 100)
	nguoi.lucHoiSinh = t0

	bayGio := t0.Add(1500 * time.Millisecond) // DUNG 1500ms — het bat tu

	mt, nhip := q.chonMucTieu([]*nguoiTrongKhuVuc{nguoi}, bayGio)

	if mt == nil {
		t.Fatal("chonMucTieu tra nil, mong doi chon lai nguoi sau khi het 1500ms bat tu")
	}
	if mt.nv.ID != 1 {
		t.Fatalf("chon nguoi id=%d, mong doi id=1", mt.nv.ID)
	}
	if nhip != nhipTuGaySu {
		t.Fatalf("nhip = %v, mong doi %v (nhip tang 2)", nhip, nhipTuGaySu)
	}
}

// TestChonMucTieuTangMotCungBoQuaNguoiVuaHoiSinh: dac ta §1.7 ghi ro "KHONG
// gioi han huong" — ap dung cho CA tang 1 (tra thu), khong chi tang 2. Nguoi
// duy nhat THOA tat ca dieu kien khac cua tang 1 (da tung danh quai, trong
// tam 300) nhung vua hoi sinh 100ms truoc -> van phai bi loai, KHONG duoc
// "tut" xuong tang 2 (vi tang 2 cung loai nguoi nay boi cung mot ly do).
func TestChonMucTieuTangMotCungBoQuaNguoiVuaHoiSinh(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	q := &Quai{X: 0, Y: 0, DaTungDanh: map[int64]bool{1: true}}

	nguoi := nguoiGia(1, 10, 0, 100) // da tung danh + rat gan (ca hai tam)
	nguoi.lucHoiSinh = t0

	bayGio := t0.Add(100 * time.Millisecond) // 100ms < 1500ms bat tu

	mt, nhip := q.chonMucTieu([]*nguoiTrongKhuVuc{nguoi}, bayGio)

	if mt != nil {
		t.Fatalf("chonMucTieu chon nguoi id=%d dang BAT TU o tang 1 (tra thu) — dac ta §1.7 khong gioi han tang", mt.nv.ID)
	}
	if nhip != 0 {
		t.Fatalf("nhip = %v, mong doi 0 khi khong co muc tieu", nhip)
	}
}

// TestNhipTimAiQuaiBoQuaNguoiVuaHoiSinhRoiDanhLaiSauBatTu la kich ban O MUC
// ACTOR (nhipTim/aiQuai, khuvuc.go) — bo sung cho ba test chonMucTieu thuan o
// tren bang mot kich ban CO THAT qua actor that su: quai AI khong duoc phat
// goi -11 nao trong khoang bat tu, roi PHAI danh duoc ngay khi het bat tu (tai
// khoang khac voi nhip 2000ms cua tang 2, vi day la lan danh DAU TIEN — dac ta
// §5.4.2 "quai vua thay nguoi phai duoc danh NGAY" van dung song song voi bat
// tu, khong tron lan hai co che).
func TestNhipTimAiQuaiBoQuaNguoiVuaHoiSinhRoiDanhLaiSauBatTu(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucNhipTimGia(&luc)
	k.nguoiTrong[1].nv.HpHienTai = 100
	k.nguoiTrong[1].nv.ViTriX = 0
	k.nguoiTrong[1].nv.ViTriY = 0
	k.nguoiTrong[1].lucHoiSinh = t0 // vua hoi sinh DUNG luc t0
	k.quai = []*Quai{{Id: 0, Hp: 30, HpToiDa: 30, X: 0, Y: 0, SatThuongCoBan: 5, LanHoiMauCuoi: t0}}

	// Tick tai t0+500ms: van trong 1500ms bat tu -> quai KHONG duoc danh, du
	// nguoi choi dang dung ngay trong tam (0,0) va khong ai khac de chon.
	luc = t0.Add(500 * time.Millisecond)
	k.nhipTim()
	if len(p.nhan) != 0 {
		t.Fatalf("nhan %d goi trong khoang bat tu (500ms sau hoi sinh), mong doi 0 — dac ta §1.7", len(p.nhan))
	}
	if !k.quai[0].LanDanhCuoi.IsZero() {
		t.Fatalf("LanDanhCuoi = %v, mong doi VAN LA zero value — khong co muc tieu (bi bat tu loai) thi TUYET DOI khong duoc cham truong nay (dac ta §5.4.2)", k.quai[0].LanDanhCuoi)
	}

	// Tick tai t0+1500ms (DUNG nguong bat tu): quai duoc phep danh NGAY (lan
	// danh dau tien, LanDanhCuoi con la zero value nen khong bi chan boi nhip).
	luc = t0.Add(1500 * time.Millisecond)
	k.nhipTim()
	if len(p.nhan) != 1 {
		t.Fatalf("nhan %d goi tai t0+1500ms (het bat tu), mong doi DUNG 1 (-11)", len(p.nhan))
	}
	if p.nhan[0].Command != -11 {
		t.Fatalf("Command = %d, mong doi -11", p.nhan[0].Command)
	}
	muon := []byte{
		0,          // mobId
		0, 0, 0, 5, // satThuong = SatThuongCoBan = 5
	}
	if !bytes.Equal(p.nhan[0].Data, muon) {
		t.Fatalf("goi -11 byte sai:\n got  = %v\n muon = %v", p.nhan[0].Data, muon)
	}
}

// ---------------------------------------------------------------------------
// Task 9 — noi day case suKienDanhQuai: nap DaTungDanh khi trung don, xoa
// sach khi quai chet.
// ---------------------------------------------------------------------------

// TestSuKienDanhQuaiNapVaoDanhSachDaTungDanhKhiTrungDon: mot don danh TRUNG
// (khong truot, khong giet chet quai) PHAI ghi nhan nguoi danh vao
// mucTieu.DaTungDanh (dac ta brief Task 9: "nap khi trung don").
func TestSuKienDanhQuaiNapVaoDanhSachDaTungDanhKhiTrungDon(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))

	// HpHienTai: 100 (con song) — Ke hoach 6 Task 13 them chot "nguoi da chet
	// (HpHienTai<=0) khong danh duoc" vao case suKienDanhQuai (khuvuc.go); bo
	// truong nay se mac dinh HpHienTai=0 (zero value) va bi CHOT DO chan
	// TRUOC ca kiem dang test o day, lam sai lech y nghia test.
	nvA := &store.NhanVat{ID: 1, Ten: "A", Lang: 0, ChakraHienTai: 100, HpHienTai: 100}
	pA := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvA, phien: pA})
	k.xuLy(suKienChonChieu{idNhanVat: nvA.ID, mauChieuId: 1}) // can chien, tam 100

	// Hp != HpToiDa: TRANH nhanh "quai day mau khong chet boi mot don" (§1.3)
	// can thiep vao phep tinh — khong lien quan toi dieu dang kiem o day.
	k.quai = []*Quai{{Id: 0, Hp: 100, HpToiDa: 150}}

	k.xuLy(suKienDanhQuai{idNhanVat: nvA.ID, mobId: 0})

	if !k.quai[0].DaTungDanh[nvA.ID] {
		t.Fatalf("DaTungDanh[%d] = false sau mot don TRUNG, mong doi true (brief Task 9: 'danh sach da tung danh no nap khi trung don')", nvA.ID)
	}
}

// TestSuKienDanhQuaiXoaDanhSachDaTungDanhKhiQuaiChet: mot don danh GIET CHET
// quai PHAI xoa sach DaTungDanh (dac ta brief Task 9: "xoa sach khi quai
// chet") — DaTungDanh bat dau CO SAN 2 phan tu (mot trong do LA nguoi vua
// giet, mot LA nguoi khac) de chung minh viec xoa khong phu thuoc vao NOI
// DUNG danh sach, chi phu thuoc vao SU KIEN "quai chet".
func TestSuKienDanhQuaiXoaDanhSachDaTungDanhKhiQuaiChet(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))

	// HpHienTai: 100 (con song) — Ke hoach 6 Task 13 them chot "nguoi da chet
	// (HpHienTai<=0) khong danh duoc" vao case suKienDanhQuai (khuvuc.go); bo
	// truong nay se mac dinh HpHienTai=0 (zero value) va bi CHOT DO chan
	// TRUOC ca kiem dang test o day, lam sai lech y nghia test.
	nvA := &store.NhanVat{ID: 1, Ten: "A", Lang: 0, ChakraHienTai: 100, HpHienTai: 100}
	pA := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: nvA, phien: pA})
	k.xuLy(suKienChonChieu{idNhanVat: nvA.ID, mauChieuId: 1})

	// Hp=1: sat thuong (toi thieu vai diem o cap 1) chac chan >= 1 -> quai
	// chet sau DUNG mot don (kep §1.2 dam bao satThuong bi giu <= Hp = 1).
	k.quai = []*Quai{{Id: 0, Hp: 1, HpToiDa: 150, DaTungDanh: map[int64]bool{999: true, nvA.ID: true}}}

	k.xuLy(suKienDanhQuai{idNhanVat: nvA.ID, mobId: 0})

	if k.quai[0].Hp != 0 {
		t.Fatalf("quai chua chet sau don danh (Hp=%d) — sua lai gia thiet test (can Hp ban dau du nho de chac chan chet)", k.quai[0].Hp)
	}
	if len(k.quai[0].DaTungDanh) != 0 {
		t.Fatalf("DaTungDanh con %d phan tu sau khi quai chet, mong doi RONG (brief Task 9: 'xoa sach khi quai chet')", len(k.quai[0].DaTungDanh))
	}
}

// ---------------------------------------------------------------------------
// Task 9 — noi day nhipTim: aiQuai chay moi tick cho quai con song, gui goi
// -11/-10 byte-chinh-xac, VA test bat buoc "khong co muc tieu thi nhip khong
// bi tieu" (dac ta §5.4.2) bang mot kich ban CO THAT (khong chi kiem
// chonMucTieu don le).
// ---------------------------------------------------------------------------

// TestNhipTimAiQuaiKhongCoMucTieuThiNhipKhongBiTieu la TEST BAT BUOC cua
// brief Task 9 (dac ta §5.4.2 — bay ban Java Mob.java:386-390 dat lai moc
// thoi gian NGOAI nhanh "co nguoi choi"). Kich ban: tick 1 KHONG co ai trong
// tam (khong danh, khong duoc cham q.LanDanhCuoi); nguoi choi buoc vao tam
// tang 2 (100) chi 100ms sau — PHAI bi danh NGAY o tick ke tiep, KHONG phai
// cho du 2000ms (nhip tang 2) tinh tu tick 1.
func TestNhipTimAiQuaiKhongCoMucTieuThiNhipKhongBiTieu(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucNhipTimGia(&luc)
	k.nguoiTrong[1].nv.HpHienTai = 50
	k.nguoiTrong[1].nv.ViTriX = 10_000 // rat xa, ngoai CA HAI tang (300 va 100)
	k.quai = []*Quai{{Id: 0, Hp: 30, HpToiDa: 30, X: 0, Y: 0, SatThuongCoBan: 5, LanHoiMauCuoi: t0}}

	k.nhipTim() // tick 1, luc = t0: khong ai trong tam.

	if len(p.nhan) != 0 {
		t.Fatalf("nhan %d goi khi chua co ai trong tam, mong doi 0", len(p.nhan))
	}
	if !k.quai[0].LanDanhCuoi.IsZero() {
		t.Fatalf("LanDanhCuoi = %v, mong doi VAN LA zero value — khong co muc tieu thi TUYET DOI khong duoc cham truong nay (dac ta §5.4.2)", k.quai[0].LanDanhCuoi)
	}

	// Nguoi choi buoc vao tam tang 2, chi 100ms sau tick khong-muc-tieu.
	k.nguoiTrong[1].nv.ViTriX = 50
	luc = t0.Add(100 * time.Millisecond)
	k.nhipTim() // tick 2

	if len(p.nhan) != 1 {
		t.Fatalf("nhan %d goi sau khi nguoi vao tam, mong doi DUNG 1 (-11) — quai PHAI danh NGAY khi vua thay nguoi, khong phai cho them gan du 2000ms (dac ta §5.4.2)", len(p.nhan))
	}
	if p.nhan[0].Command != -11 {
		t.Fatalf("Command = %d, mong doi -11", p.nhan[0].Command)
	}
	muon := []byte{
		0,          // mobId
		0, 0, 0, 5, // satThuong = SatThuongCoBan = 5 (khong vuot HpHienTai=50)
	}
	if !bytes.Equal(p.nhan[0].Data, muon) {
		t.Fatalf("goi -11 byte sai:\n got  = %v\n muon = %v", p.nhan[0].Data, muon)
	}
	if k.nguoiTrong[1].nv.HpHienTai != 45 {
		t.Fatalf("HpHienTai nan nhan = %d, mong doi 50-5=45", k.nguoiTrong[1].nv.HpHienTai)
	}
}

// TestNhipTimAiQuaiTonTrongNhipTangHai: sau don danh dau tien (tang 2, nhip
// 2000ms), mot tick 100ms sau KHONG duoc danh tiep (con trong nhip hoi); mot
// tick du 2000ms tinh tu LAN DANH TRUOC moi duoc danh lan 2.
func TestNhipTimAiQuaiTonTrongNhipTangHai(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucNhipTimGia(&luc)
	k.nguoiTrong[1].nv.HpHienTai = 100
	k.nguoiTrong[1].nv.ViTriX = 0 // trong tam tang 2 (100) tu dau
	k.quai = []*Quai{{Id: 2, Hp: 30, HpToiDa: 30, X: 0, Y: 0, SatThuongCoBan: 3, LanHoiMauCuoi: t0}}

	k.nhipTim() // tick 1, luc=t0: danh lan 1 (LanDanhCuoi zero -> elapsed lon)
	if len(p.nhan) != 1 {
		t.Fatalf("nhan %d goi sau tick dau, mong doi DUNG 1 (danh lan 1)", len(p.nhan))
	}

	luc = t0.Add(100 * time.Millisecond) // 100ms < 2000ms nhip tang 2
	k.nhipTim()
	if len(p.nhan) != 1 {
		t.Fatalf("nhan %d goi 100ms sau lan danh dau, mong doi VAN LA 1 (chua du nhip 2000ms tang 2)", len(p.nhan))
	}

	luc = t0.Add(2000 * time.Millisecond) // DUNG du nhip, tinh tu LAN DANH TRUOC (khong phai tu t0 don thuan)
	k.nhipTim()
	if len(p.nhan) != 2 {
		t.Fatalf("nhan %d goi tick du 2000ms tu lan danh truoc, mong doi 2 (danh lan 2)", len(p.nhan))
	}
	if p.nhan[1].Command != -11 {
		t.Fatalf("Command lan 2 = %d, mong doi -11", p.nhan[1].Command)
	}
}

// TestNhipTimAiQuaiTangMotVaGoi10ChoNguoiKhac: quai co MOT nguoi da tung
// danh no (tang 1) trong tam, VA mot nguoi thu hai RAT XA (chi la nguoi
// DUNG, khong phai muc tieu — chung minh -10 chi gui NGUOI KHAC, khong gui
// nan nhan) — dac ta §3.3: -11 rieng nan nhan, -10 nguoi khac trong khu vuc.
func TestNhipTimAiQuaiTangMotVaGoi10ChoNguoiKhac(t *testing.T) {
	t0 := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	luc := t0
	k, pNanNhan := khuVucNhipTimGia(&luc)
	k.nguoiTrong[1].nv.HpHienTai = 100
	k.nguoiTrong[1].nv.ViTriX = 0
	k.nguoiTrong[1].nv.ViTriY = 0

	pNguoiDung := &phienGiaKhuVuc{}
	k.nguoiTrong[2] = &nguoiTrongKhuVuc{
		nv:    &store.NhanVat{ID: 2, Ten: "NguoiDung", HpHienTai: 100, ViTriX: 99_999},
		phien: pNguoiDung,
	}

	k.quai = []*Quai{{
		Id: 5, Hp: 20, HpToiDa: 20, X: 0, Y: 0, SatThuongCoBan: 4, LanHoiMauCuoi: t0,
		DaTungDanh: map[int64]bool{1: true}, // nguoi 1 da tung danh quai nay -> tang 1
	}}

	k.nhipTim()

	if len(pNanNhan.nhan) != 1 || pNanNhan.nhan[0].Command != -11 {
		t.Fatalf("nan nhan (id=1) phai nhan DUNG 1 goi -11, nhan %+v", pNanNhan.nhan)
	}
	if len(pNguoiDung.nhan) != 1 || pNguoiDung.nhan[0].Command != -10 {
		t.Fatalf("nguoi dung (id=2) phai nhan DUNG 1 goi -10, nhan %+v", pNguoiDung.nhan)
	}
	muonNguoiDung := []byte{
		5,          // mobId
		0, 0, 0, 1, // idNanNhan = 1
		0, 0, 0, 96, // hpConLai = 100-4=96
	}
	if !bytes.Equal(pNguoiDung.nhan[0].Data, muonNguoiDung) {
		t.Fatalf("goi -10 byte sai:\n got  = %v\n muon = %v", pNguoiDung.nhan[0].Data, muonNguoiDung)
	}
}
