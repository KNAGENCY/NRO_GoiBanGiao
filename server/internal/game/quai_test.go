package game

import (
	"bytes"
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

// TestGhiDanhSachQuaiByteDung: so TOAN BO mang byte ma ghiDanhSachQuai ghi ra
// cho 2 con quai voi gia tri ghim cung, dung dac ta §3.6.
func TestGhiDanhSachQuaiByteDung(t *testing.T) {
	w := protocol.NewWriter()
	ghiDanhSachQuai(w, haiQuaiMau())

	muon := []byte{
		2, // soQuai

		// --- quai 1: Id=0, Cap=3, Hp=80, HpToiDa=100, X=120, Y=200 ---
		0, 0, 0, 0, 0, // isDisable/isDontMove/isFire/isIce/isWind = false
		0,           // tempId = 0
		0,           // sys = 0
		0, 0, 0, 80, // hp = 80 (int32, TRUOC level)
		3,            // level = 3 (GIUA hp va hpMax)
		0, 0, 0, 100, // hpMax = 100 (int32, SAU level)
		0, 120, // x = 120
		0, 200, // y = 200 (200 < 256, 1 byte du dung 2 byte short)
		0, // status = 0
		3, // lvMob = 3
		0, // isBigBoss = false

		// --- quai 2: Id=1, Cap=10, Hp=500, HpToiDa=500, X=300, Y=150 ---
		0, 0, 0, 0, 0,
		1,            // tempId = 1
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

	// 5 co bool + tempId + sys.
	for i := 0; i < 5; i++ {
		if _, err := r.Byte(); err != nil {
			t.Fatalf("doc co bool thu %d loi: %v", i+1, err)
		}
	}
	if _, err := r.Byte(); err != nil { // tempId
		t.Fatalf("doc tempId loi: %v", err)
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
	b("waypointCount", 0)

	b("mobCount", 2)

	// --- quai 1 ---
	b("q1.isDisable", 0)
	b("q1.isDontMove", 0)
	b("q1.isFire", 0)
	b("q1.isIce", 0)
	b("q1.isWind", 0)
	b("q1.tempId", 0)
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
	b("q2.tempId", 1)
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
	b("npcCount", 0)
	b("itemCount", 0)
	sh("bgItem", 0)
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
		if *q != *goc {
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

	nvA := &store.NhanVat{ID: 1, Ten: "A", Lang: 0, ChakraHienTai: 100}
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
