package game

import (
	"testing"
	"time"

	"nroserver/internal/protocol"
)

// ─────────────────────────────────────────────────────────────────────────────
// KH6b — 3 test bit 3 lo hong ma kiem thu dot bien phat hien (M19, M22, M26).
// Moi test o day khoa DUNG mot hanh vi ma cac test cu (khuvuc_aiboss_test.go,
// khuvuc_bienhinh_test.go) khong cham toi, nen mot dot bien tuong ung van
// "song sot" (test cu van xanh du logic sai).
// ─────────────────────────────────────────────────────────────────────────────

// ─────────────────────────────────────────────────────────────────────────────
// M19 — chonMucTieuBoss bo tang 1 (nguoi gop sat thuong nhieu nhat).
//
// Dot bien: `if gopNhat != nil { return gopNhat }; return ganNhat` bi rut gon
// thanh `return ganNhat` — tuc la BO QUA HOAN TOAN tang uu tien "ai danh boss
// nhieu nhat", boss luon chon nguoi GAN NHAT bat ke ai da do mau vao no.
//
// Test nay dung THANG chonMucTieuBoss (ham thuan, khong can nhip tim) de tach
// rieng logic chon muc tieu khoi toan bo may AI — dung ket qua tra ve la du,
// khong can boc qua goi -60.
// ─────────────────────────────────────────────────────────────────────────────
func TestChonMucTieuBossUuTienNguoiGopSatThuongHonNguoiGanHon(t *testing.T) {
	bayGio := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	b := bossChienDau(bayGio, xBossThuNghiem, 1500, nil)

	// A da danh boss (co ghi cong vao SatThuongTheoNguoi) nhung DUNG XA hon B.
	// B chua danh phat nao nhung dung GAN boss hon A.
	//
	// Bay: neu chi dat A ngoai tam chieu thi test se xanh gia — boss chon A vi
	// ly do "B ngoai tam" (dung logic tang 2) chu khong phai vi tang 1 hoat
	// dong. Nen O DAY ca hai PHAI cung nam trong tam chieu 100 (can chien).
	// Dung khu vuc that (khuVucAiBoss) chi de co san map nguoiTrong cho
	// themNguoiAiBoss ghi vao — chonMucTieuBoss nhan mang []*nguoiTrongKhuVuc
	// truc tiep nen khong dung toi nhip tim hay dong ho cua k.
	luc := bayGio
	k := khuVucAiBoss(&luc)
	a, _ := themNguoiAiBoss(k, 1, xBossThuNghiem+90, yDiDuocBoss)  // xa hon, trong tam
	bb, _ := themNguoiAiBoss(k, 2, xBossThuNghiem+20, yDiDuocBoss) // gan hon, chua danh
	b.SatThuongTheoNguoi = map[int64]int64{1: 500}

	nguoi := []*nguoiTrongKhuVuc{a, bb}
	muc := chonMucTieuBoss(b, nguoi, 100, bayGio)

	if muc == nil {
		t.Fatalf("chonMucTieuBoss tra nil, mong doi nguoi 1 (A)")
	}
	if muc.nv.ID != 1 {
		t.Errorf("boss chon nguoi %d, mong doi nguoi 1 (da gop sat thuong) du B dung gan hon — "+
			"day chinh la lo hong M19: bo tang 1 se khien boss luon chon nguoi GAN NHAT",
			muc.nv.ID)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// M22 — mocNoiGianVuaVuot dung sai dau bien o bien mocs (>= thay vi >).
//
// Chu thich ham noi ro: mocNoiGian VUOT khi `hp*100 <= hpToiDa*PhanTramHP`,
// tuc CON DUNG 50% mau (hp=500 tren hpToiDa=1000) THI mocNoiGian 50% PHAI no.
// Dot bien doi dieu kien "chua toi moc" tu `>` thanh `>=` se lam hp=500 (dung
// bang nguong) bi coi la "chua toi", tuc mocNoiGian tre mot nap.
//
// Test don vi thuan tren mocNoiGianVuaVuot, khong qua may trang thai boss.
// ─────────────────────────────────────────────────────────────────────────────
func TestMocNoiGianVuaVuotBienDungMocKhiConDungPhanTramHP(t *testing.T) {
	moc := []MocNoiGian{{PhanTramHP: 50}}
	hpToiDa := int32(1000)

	// hp = 500 -> dung bang 50% -> PHAI da vuot (tra chi so 0).
	if got := mocNoiGianVuaVuot(moc, 500, hpToiDa, 0); got != 0 {
		t.Errorf("hp=500 (dung 50%% cua %d) -> mocNoiGianVuaVuot = %d, mong doi 0 (da vuot moc) — "+
			"day la bien M22: doi > thanh >= se lam dung diem bien nay bi coi la CHUA vuot",
			hpToiDa, got)
	}

	// hp = 501 -> con TREN 50% mot chut -> CHUA vuot (tra -1).
	if got := mocNoiGianVuaVuot(moc, 501, hpToiDa, 0); got != -1 {
		t.Errorf("hp=501 (tren 50%% cua %d) -> mocNoiGianVuaVuot = %d, mong doi -1 (chua vuot)",
			hpToiDa, got)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// M26 — thoaiDanhBoss khong quay vong chi so (thieu % len(pha.ThoaiDanh) khi
// cong don). Dot bien: `b.ChiSoThoaiDanh = b.ChiSoThoaiDanh + 1` (khong %).
//
// QUAN TRONG — vi sao test PHAI kiem truc tiep b.ChiSoThoaiDanh (khong chi
// kiem cau thoai phat ra): dong doc `pha.ThoaiDanh[b.ChiSoThoaiDanh%len(...)]`
// (dong TRUOC dong bi dot bien) DA CO san phep %, nen ca hai ban — dung va
// dot bien — phat ra CUNG MOT chuoi cau thoai (0,1,0,1,...) trong pham vi mot
// tran boss binh thuong; chi khac o gia tri NOI BO cua ChiSoThoaiDanh sau moi
// lan noi. Neu test chi so sanh noi dung cau thoai thi dot bien nay SONG SOT
// (da tu kiem chung: ap dot bien, test theo cau thoai van xanh). Phai doc
// truc tiep truong b.ChiSoThoaiDanh moi giet duoc no — va day cung chinh la
// nguy co that: chi so cong don khong gioi han se sat MaxInt roi TRAN SO
// (chuyen am) sau du lau, luc do dong doc %len co the panic hoac chon sai
// cau vi mot so am mod mot so duong trong Go van co the am.
//
// Test nay bom dong ho qua 3 nhip thoaiDanhBoss (hang so nhipThoaiDanhBoss =
// 15s) voi mang ThoaiDanh 2 cau, kiem cau thu 3 == cau thu 1 VA kiem
// ChiSoThoaiDanh van nam trong [0, len(ThoaiDanh)) sau 3 lan noi.
// ─────────────────────────────────────────────────────────────────────────────
func TestThoaiDanhBossQuayVongKhongPanic(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k := khuVucAiBoss(&luc)
	thoai := []string{"Cau mot", "Cau hai"}
	_, p := themNguoiAiBoss(k, 1, xBossThuNghiem+50, yDiDuocBoss)
	b := bossChienDau(luc, xBossThuNghiem, 1500, thoai)
	// nhipThoaiDanhBoss = 15s va nhip danh (1500ms) ngan hon nhieu -> boss vua
	// co the danh vua co the noi trong cung nhip tim, khong anh huong nhau vi
	// hai co che doc lap (aiBoss goi lien tiep ca hai).
	k.boss = append(k.boss, b)

	// Nhip 1: cau ngay lap tuc (MocThoaiCuoi zero -> Sub rat lon).
	k.nhipTim()
	// Nhip 2: bom qua du 15s.
	luc = luc.Add(nhipThoaiDanhBoss + 100*time.Millisecond)
	k.nhipTim()
	// Nhip 3: bom qua du 15s nua -> day la nhip chi so QUAY VONG ve 0. Neu
	// dot bien M26 duoc ap (khong %) thi ChiSoThoaiDanh o day = 2, truy cap
	// pha.ThoaiDanh[2] tren mang 2 phan tu -> panic index out of range —
	// chinh xac cai ma bao cao phai neu ro neu no xay ra.
	luc = luc.Add(nhipThoaiDanhBoss + 100*time.Millisecond)
	k.nhipTim()

	goiChat := locGoi(p, lenhChat)
	if len(goiChat) != 3 {
		t.Fatalf("so cau thoai danh phat ra = %d, mong doi 3 (mot cau moi nhip 15s)", len(goiChat))
	}
	cau1 := docNoiDungChat(t, goiChat[0])
	cau3 := docNoiDungChat(t, goiChat[2])
	if cau1 != thoai[0] {
		t.Errorf("cau thoai dau tien = %q, mong doi %q", cau1, thoai[0])
	}
	if cau3 != cau1 {
		t.Errorf("cau thoai thu 3 = %q, mong doi bang cau thoai dau tien %q (%q) — "+
			"day chinh la lo hong M26: thieu %%len se lam chi so chay thang thay vi quay vong",
			cau3, cau1, thoai)
	}

	// Kiem TRUC TIEP truong noi bo: day la phep kiem THAT SU giet dot bien
	// M26 (xem chu thich phia tren) — noi dung cau thoai KHONG PHAN BIET
	// duoc hai phien ban vi dong doc da co san %.
	if b.ChiSoThoaiDanh < 0 || b.ChiSoThoaiDanh >= len(thoai) {
		t.Errorf("b.ChiSoThoaiDanh = %d sau 3 lan noi, mong doi nam trong [0,%d) — "+
			"day la lo hong M26: thieu %%len(pha.ThoaiDanh) khi cong don se lam chi so "+
			"chay thang khong gioi han thay vi quay vong, som muon tran so",
			b.ChiSoThoaiDanh, len(thoai))
	}
}

// docNoiDungChat doc noi dung chuoi tu mot goi chat44 (lenhChat), dung de
// kiem cau thoai boss noi ra co dung noi dung va dung thu tu quay vong hay
// khong.
func docNoiDungChat(t *testing.T, m *protocol.Message) string {
	t.Helper()
	if m.Command != lenhChat {
		t.Fatalf("command = %d, mong doi %d (chat44)", m.Command, lenhChat)
	}
	r := protocol.NewReader(m.Data)
	_, _ = r.Int() // id nguoi noi (boss, am) — khong can kiem o day
	s, err := r.UTF()
	if err != nil {
		t.Fatalf("doc UTF that bai: %v", err)
	}
	return s
}
