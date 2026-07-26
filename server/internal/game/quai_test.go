package game

import (
	"bytes"
	"testing"

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
