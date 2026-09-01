package game

import (
	"testing"
	"time"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// ─────────────────────────────────────────────────────────────────────────────
// KH6b Task 10 — AI BOSS trong nhip 100ms (buoc HOAT_DONG cua may trang thai).
//
// File nay khoa SAU dieu, moi dieu la mot cach boss co the "song sai" ma khong
// bao mot loi nao:
//
//  1. Boss danh nguoi choi trong tam — bang goi -60 (KHONG phai -9/-12 cua
//     quai: hai goi do mang `byte mobId`, mot id boss am 4 byte khong nhet vao
//     duoc), mau nan nhan tru that, nguoi KHAC trong khu vuc nhan -30/9.
//  2. Boss KHONG danh nguoi dang bat tu sau hoi sinh, nguoi da chet, va nguoi
//     khong con trong khu vuc.
//  3. Ngoai tam thi boss DI CHUYEN lai gan (goi -7), khong danh voi tu xa.
//  4. Boss KHONG di xuyen dat — cung luoi va cham ma nguoi choi phai theo.
//  5. Nhip danh theo PhaBoss.NhipDanhMs, khong phai moi tick 100ms.
//  6. ThoaiDanh la "thinh thoang", khong phai spam moi nhip.
//
// KY THUAT: goi k.nhipTim() TRUC TIEP tren khu vuc CHUA BatDau() (khuon cua
// khuvuc_vongdoiboss_test.go) — mot goroutine duy nhat, dong ho tiem duoc,
// KHONG sleep that.
// ─────────────────────────────────────────────────────────────────────────────

// Toa do dung chung: hang y=384 cua map Lang La di duoc suot mot doan dai
// (khuVucBossGia da dat nguoi xem o day), nen moi test dat boss + muc tieu
// tren cung hang do va chi doi X — khoang cach doc thang tu hieu hai so X.
const (
	yDiDuocBoss    int16 = 384
	xBossThuNghiem int16 = 500
)

// khuVucAiBoss dung mot khu vuc trong (KHONG co san nguoi xem nhu
// khuVucBossGia — o day tung test tu dat nguoi vao dung cho no can) voi dong
// ho tiem duoc va dan quai da don sach.
func khuVucAiBoss(luc *time.Time) *KhuVuc {
	k := MoKhuVuc(TimBanDo(1))
	k.dongHo = func() time.Time { return *luc }
	k.quai = nil
	return k
}

// themNguoiAiBoss dat mot nguoi choi vao khu vuc tai (x, y) va tra ve ca ban
// ghi lan phien gia de dem goi.
func themNguoiAiBoss(k *KhuVuc, id int64, x, y int16) (*nguoiTrongKhuVuc, *phienGiaKhuVuc) {
	p := &phienGiaKhuVuc{}
	n := &nguoiTrongKhuVuc{
		nv: &store.NhanVat{
			ID: id, Ten: "Nguoi", Lang: 0,
			HpToiDa: 1000, HpHienTai: 1000,
			ViTriX: int32(x), ViTriY: int32(y),
		},
		phien: p,
	}
	k.nguoiTrong[id] = n
	return n, p
}

// bossChienDau dung mot boss MOT PHA dang HOAT_DONG, co du bo chieu (can chien
// tam 100 + tam xa tam 300) — dung dung bo chieu cua boss mau trong
// noidungboss.go, khong bia so.
func bossChienDau(luc time.Time, x int16, nhipMs int32, thoaiDanh []string) *Boss {
	return &Boss{
		Id: -7,
		DuLieu: []PhaBoss{{
			Ten:        "Kẻ Lạ Mặt",
			Lang:       0,
			Cap:        50,
			TypePk:     TypePkBoss,
			Dame:       60,
			NhipDanhMs: nhipMs,
			BoChieu:    [][2]int32{{1, 2}, {2, 2}}, // {can chien, cap 2}, {tam xa, cap 2}
			ThoaiDanh:  thoaiDanh,
		}},
		Hp: 3000, HpToiDa: 3000,
		ViTriX: x, ViTriY: yDiDuocBoss,
		TrangThai:       BossHoatDong,
		MocThoiGian:     luc,
		LucHienPha:      luc,
		NhipDanhHienTai: nhipMs,
	}
}

// docDanh60 doc nguoc mot goi -60 dang "don danh" thanh cac truong can kiem.
func docDanh60(t *testing.T, m *protocol.Message) (idNguoiDanh, idNguoiTrung, satThuong int32, daChet bool) {
	t.Helper()
	if m.Command != lenhNguoiDanhNguoi {
		t.Fatalf("command = %d, mong doi %d (-60)", m.Command, lenhNguoiDanhNguoi)
	}
	r := protocol.NewReader(m.Data)
	idNguoiDanh, _ = r.Int()
	_, _ = r.Byte() // skillId
	_, _ = r.Byte() // soMucTieu
	idNguoiTrung, _ = r.Int()
	_, _ = r.Byte() // docTiep
	_, _ = r.Byte() // loaiChieu
	satThuong, _ = r.Int()
	// Bool ghi bang 1 byte (protocol.Writer.Bool) — Reader khong co ham Bool
	// nen doc lai bang Byte va so voi 0, dung cach cac test goi khac lam.
	byteChet, _ := r.Byte()
	daChet = byteChet != 0
	return
}

// ─────────────────────────────────────────────────────────────────────────────
// 1. Boss danh nguoi trong tam
// ─────────────────────────────────────────────────────────────────────────────

// TestBossHoatDongDanhNguoiTrongTam thay the TestBossHoatDongDungYen_ChoTask10
// (Task 9 ghim nguyen trang "im lang" de Task 10 buoc phai dung lai doc).
//
// Khoa BON thu cua mot don boss danh nguoi:
//   - goi -60 mang ID BOSS (am) o o "nguoi danh" va id nan nhan o o "nguoi
//     trung" — KHONG phai -9/-12 cua quai.
//   - mau nan nhan TRU THAT tren server (khong chi la mot goi trang tri).
//   - con so trong goi = luong mau da tru that.
//   - NGUOI KHAC trong khu vuc nhan -30/9 de ve lai thanh mau nan nhan (goi
//     -60 khong mang truong hp con lai) — nan nhan thi KHONG, client cua ho tu
//     tru tu chinh goi -60, dung duong ma nhanh PvP (case suKienDanhNguoiChoi)
//     da di.
func TestBossHoatDongDanhNguoiTrongTam(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k := khuVucAiBoss(&luc)
	nanNhan, pNanNhan := themNguoiAiBoss(k, 1, xBossThuNghiem+50, yDiDuocBoss) // trong tam can chien 100
	_, pNguoiKhac := themNguoiAiBoss(k, 2, xBossThuNghiem+60, yDiDuocBoss)
	b := bossChienDau(luc, xBossThuNghiem, 1500, nil)
	k.boss = append(k.boss, b)

	mauTruoc := nanNhan.nv.HpHienTai
	k.nhipTim()

	goi := locGoi(pNanNhan, lenhNguoiDanhNguoi)
	if len(goi) != 1 {
		t.Fatalf("boss phat %d goi -60, mong doi dung 1", len(goi))
	}
	idDanh, idTrung, sat, daChet := docDanh60(t, goi[0])
	if idDanh != -7 {
		t.Errorf("goi -60 ghi nguoi danh = %d, mong doi -7 (ID BOSS)", idDanh)
	}
	if idTrung != 1 {
		t.Errorf("goi -60 ghi nan nhan = %d, mong doi 1", idTrung)
	}
	if sat <= 0 {
		t.Fatalf("sat thuong = %d, mong doi > 0", sat)
	}
	if daChet {
		t.Errorf("nan nhan 1000 mau ma bao chet ngay mot don")
	}
	if mauTruoc-nanNhan.nv.HpHienTai != sat {
		t.Errorf("mau tru that = %d, con so bao xuong client = %d — hai so PHAI bang nhau",
			mauTruoc-nanNhan.nv.HpHienTai, sat)
	}
	if n := len(locGoi(pNanNhan, lenhMeInfo)); n != 0 {
		t.Errorf("nan nhan nhan %d goi -30/9 — goi do CHI danh cho nguoi khac", n)
	}
	if n := len(locGoi(pNguoiKhac, lenhMeInfo)); n != 1 {
		t.Errorf("nguoi khac nhan %d goi -30/9, mong doi 1 (de ve lai thanh mau nan nhan)", n)
	}
}

// TestBossKhongDanhNguoiVuaHoiSinh: nguoi vua hoi sinh co 1500ms bat tu (dac ta
// §1.7) — boss phai TON TRONG khoang do y het nguoi choi va quai. Bo kiem nay
// nghia la boss don chet nguoi choi ngay khi ho vua bam nut hoi sinh.
func TestBossKhongDanhNguoiVuaHoiSinh(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k := khuVucAiBoss(&luc)
	n, p := themNguoiAiBoss(k, 1, xBossThuNghiem+50, yDiDuocBoss)
	n.lucHoiSinh = t0.Add(-1000 * time.Millisecond) // moi 1s, con bat tu
	// NGUOI THU HAI, XA HON mot chut, hoan toan hop le. Khong co anh ta thi test
	// nay YEU: bo kiem bat tu trong chonMucTieuBoss van "xanh" vi KiemDanhHopLe
	// tu choi don o buoc sau. Co anh ta thi khac biet lo ra ngay — chon nham
	// nguoi bat tu nghia la BO QUEN nguoi dang danh minh.
	nguoiThat, pThat := themNguoiAiBoss(k, 2, xBossThuNghiem+60, yDiDuocBoss)
	b := bossChienDau(luc, xBossThuNghiem, 1500, nil)
	k.boss = append(k.boss, b)

	k.nhipTim()

	if got := len(locGoi(p, lenhNguoiDanhNguoi)); got != 1 {
		t.Fatalf("boss phat %d goi -60, mong doi 1 (danh NGUOI KHONG BAT TU)", got)
	}
	if _, idTrung, _, _ := docDanh60(t, locGoi(p, lenhNguoiDanhNguoi)[0]); idTrung != 2 {
		t.Fatalf("boss danh vao id %d, mong doi 2 — id 1 dang bat tu sau hoi sinh", idTrung)
	}
	if n.nv.HpHienTai != 1000 {
		t.Errorf("mau nguoi bat tu = %d, mong doi nguyen 1000", n.nv.HpHienTai)
	}
	if nguoiThat.nv.HpHienTai >= 1000 {
		t.Errorf("nguoi hop le khong mat mau nao (%d)", nguoiThat.nv.HpHienTai)
	}
	_ = pThat

	// Het bat tu -> nguoi id 1 (GAN HON) tro lai lam muc tieu.
	luc = t0.Add(600 * time.Millisecond)
	n.lucHoiSinh = time.Time{}
	b.LanDanhCuoi = time.Time{} // bo qua nhip pha, ta chi quan tam viec CHON ai
	k.nhipTim()
	if got := len(locGoi(p, lenhNguoiDanhNguoi)); got != 2 {
		t.Fatalf("het bat tu roi ma tong so don = %d, mong doi 2", got)
	}
	if _, idTrung, _, _ := docDanh60(t, locGoi(p, lenhNguoiDanhNguoi)[1]); idTrung != 1 {
		t.Errorf("het bat tu, boss danh id %d, mong doi 1 (gan hon)", idTrung)
	}

}

// TestBossKhongDanhNguoiDaChet: mot cai xac khong phai muc tieu. Neu chon, boss
// se dam mai vao nguoi da chet trong khi nguoi con song dung ngay canh.
func TestBossKhongDanhNguoiDaChet(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k := khuVucAiBoss(&luc)
	xac, pXac := themNguoiAiBoss(k, 1, xBossThuNghiem+20, yDiDuocBoss) // gan hon
	xac.nv.HpHienTai = 0
	_, pSong := themNguoiAiBoss(k, 2, xBossThuNghiem+60, yDiDuocBoss)
	b := bossChienDau(luc, xBossThuNghiem, 1500, nil)
	k.boss = append(k.boss, b)

	k.nhipTim()

	goi := locGoi(pSong, lenhNguoiDanhNguoi)
	if len(goi) != 1 {
		t.Fatalf("boss phat %d goi -60, mong doi 1 (danh NGUOI CON SONG)", len(goi))
	}
	if _, idTrung, _, _ := docDanh60(t, goi[0]); idTrung != 2 {
		t.Errorf("boss danh vao id %d, mong doi 2 — id 1 la mot cai xac", idTrung)
	}
	if xac.nv.HpHienTai != 0 {
		t.Errorf("mau cai xac = %d, khong ai duoc dam vao no nua", xac.nv.HpHienTai)
	}
	_ = pXac
}

// TestBossKhuVucTrongThiImLang: khong co ai thi boss KHONG phat goi nao va
// KHONG dung toi LanDanhCuoi — dung bay ma aiQuai da ghi (ban Java dat lai moc
// thoi gian ngoai nhanh "co nguoi choi", khien quai vua thay nguoi phai cho
// them gan mot nhip nua moi duoc danh).
func TestBossKhuVucTrongThiImLangVaKhongChamNhip(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k := khuVucAiBoss(&luc)
	b := bossChienDau(luc, xBossThuNghiem, 1500, []string{"Không ai đến sao?"})
	k.boss = append(k.boss, b)

	for i := 0; i < 100; i++ { // 10 giay
		luc = t0.Add(time.Duration(i) * 100 * time.Millisecond)
		k.nhipTim()
	}

	if !b.LanDanhCuoi.IsZero() {
		t.Errorf("LanDanhCuoi = %v, mong doi zero — khong co ai thi khong duoc dem nhip", b.LanDanhCuoi)
	}
	if b.ViTriX != xBossThuNghiem {
		t.Errorf("boss tu di tu %d sang %d khi khu vuc trong", xBossThuNghiem, b.ViTriX)
	}
	if b.TrangThai != BossHoatDong {
		t.Errorf("trang thai = %s, mong doi HOAT_DONG", b.TrangThai)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// 2. Di chuyen
// ─────────────────────────────────────────────────────────────────────────────

// TestBossNgoaiTamThiDiChuyenChuKhongDanh: muc tieu o 400 don vi — ngoai tam
// chieu xa nhat (300) nhung trong tam duoi. Boss phai DI, va TUYET DOI khong
// duoc danh tu do (bo kiem tam = boss ban tia ca ban do).
func TestBossNgoaiTamThiDiChuyenChuKhongDanh(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k := khuVucAiBoss(&luc)
	_, p := themNguoiAiBoss(k, 1, xBossThuNghiem+400, yDiDuocBoss)
	b := bossChienDau(luc, xBossThuNghiem, 1500, nil)
	k.boss = append(k.boss, b)

	k.nhipTim()

	if got := len(locGoi(p, lenhNguoiDanhNguoi)); got != 0 {
		t.Fatalf("boss phat %d don vao muc tieu NGOAI TAM, mong doi 0", got)
	}
	di := locGoi(p, lenhDiChuyen)
	if len(di) != 1 {
		t.Fatalf("boss phat %d goi -7, mong doi 1 (di lai gan)", len(di))
	}
	if b.ViTriX <= xBossThuNghiem {
		t.Errorf("boss o X=%d, mong doi tien ve phia muc tieu (>%d)", b.ViTriX, xBossThuNghiem)
	}
	if b.ViTriX > xBossThuNghiem+buocBossMoiNhip {
		t.Errorf("boss nhay %d don vi mot nhip, toi da %d", b.ViTriX-xBossThuNghiem, buocBossMoiNhip)
	}

	// Goi -7 phai bao dung vi tri moi cua boss.
	r := protocol.NewReader(di[0].Data)
	id, _ := r.Int()
	x, _ := r.Short()
	if id != -7 || x != b.ViTriX {
		t.Errorf("goi -7 ghi (id=%d, x=%d), mong doi (-7, %d)", id, x, b.ViTriX)
	}
}

// TestBossKhongDiXuyenDat: muc tieu nam BEN KIA mot o dac (than be da duoi
// cong torii — cung luoi va cham ma nguoi choi phai theo). Boss khong duoc
// xuyen qua.
//
// Bo kiem va cham o day thi boss di xuyen nui, nguoi choi thi khong — no dung
// giua da roi ban ra, khong ai danh lai duoc.
func TestBossKhongDiXuyenDat(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0

	// Tim mot cap (vi tri boss di duoc, huong bi chan) THAT tren luoi va cham
	// — khong ghim toa do bia: doc thang collisionLangLa.
	var bx, by, mx, my int16
	timDuoc := false
	for x := int16(24); x < 1200 && !timDuoc; x += 24 {
		for y := int16(24); y < 700; y += 24 {
			if !collisionLangLa.DiDuoc(x, y) {
				continue
			}
			if collisionLangLa.DiDuoc(x, y+24) {
				continue // duoi chan khong phai dat dac -> khong dung lam tuong
			}
			bx, by = x, y
			// Muc tieu ben kia buc tuong, NGOAI tam chieu xa nhat (300) nhung
			// TRONG tam duoi (600) — bat buoc phai vay: trong tam chieu thi
			// boss danh xuyen tuong va KHONG BAO GIO thu di chuyen, test se
			// xanh ma khong he cham vao doan va cham.
			mx, my = x, y+400
			timDuoc = true
			break
		}
	}
	if !timDuoc {
		t.Skip("khong tim duoc o di duoc co dat dac ngay ben duoi tren map 1")
	}

	k := khuVucAiBoss(&luc)
	_, p := themNguoiAiBoss(k, 1, mx, my)
	b := bossChienDau(luc, bx, 1500, nil)
	b.ViTriY = by
	k.boss = append(k.boss, b)

	k.nhipTim()

	if n := len(locGoi(p, lenhDiChuyen)); n != 0 {
		t.Fatalf("boss phat %d goi -7 xuyen qua dat dac, mong doi 0", n)
	}
	if b.ViTriX != bx || b.ViTriY != by {
		t.Fatalf("boss di tu (%d,%d) sang (%d,%d) — xuyen dat", bx, by, b.ViTriX, b.ViTriY)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// 3. Nhip danh
// ─────────────────────────────────────────────────────────────────────────────

// TestBossGiuNhipDanhTheoPha: NhipDanhMs cua pha la khoang cach TOI THIEU giua
// hai don. Khong co no thi boss danh MOI TICK — 10 don/giay, khong ai song noi.
func TestBossGiuNhipDanhTheoPha(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k := khuVucAiBoss(&luc)
	_, p := themNguoiAiBoss(k, 1, xBossThuNghiem+50, yDiDuocBoss)
	b := bossChienDau(luc, xBossThuNghiem, 1500, nil)
	k.boss = append(k.boss, b)

	// 14 nhip dau (0 -> 1400ms): dung DUNG mot don, don dau tien.
	for i := 0; i < 15; i++ {
		luc = t0.Add(time.Duration(i) * 100 * time.Millisecond)
		k.nhipTim()
	}
	if got := len(locGoi(p, lenhNguoiDanhNguoi)); got != 1 {
		t.Fatalf("trong 1400ms boss danh %d don, mong doi 1 (nhip 1500ms)", got)
	}

	luc = t0.Add(1500 * time.Millisecond)
	k.nhipTim()
	if got := len(locGoi(p, lenhNguoiDanhNguoi)); got != 2 {
		t.Fatalf("du 1500ms roi ma tong so don = %d, mong doi 2", got)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// 4. Thoai danh
// ─────────────────────────────────────────────────────────────────────────────

// TestBossThoaiDanhThinhThoangChuKhongSpam: "thinh thoang" phai la mot con so.
// Trong 60 giay danh nhau lien tuc, boss noi dung so lan cua nhip
// nhipThoaiDanhBoss — khong phai 600 lan (moi tick mot cau).
func TestBossThoaiDanhThinhThoangChuKhongSpam(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k := khuVucAiBoss(&luc)
	_, p := themNguoiAiBoss(k, 1, xBossThuNghiem+50, yDiDuocBoss)
	b := bossChienDau(luc, xBossThuNghiem, 1500, []string{"A", "B"})
	k.boss = append(k.boss, b)

	for i := 0; i < 600; i++ { // 60 giay
		luc = t0.Add(time.Duration(i) * 100 * time.Millisecond)
		k.nhipTim()
	}

	cau := cauDaNoi(t, p)
	// Nhip cuoi cung o moc 59,9s (600 nhip tinh tu 0), nen so cau la
	// floor(59,9 / 15) + 1 — cong 1 vi cau DAU TIEN ra ngay o moc 0.
	mongDoi := int(59900*time.Millisecond/nhipThoaiDanhBoss) + 1
	if len(cau) != mongDoi {
		t.Fatalf("60 giay boss noi %d cau (%v), mong doi %d cau — nhip %v",
			len(cau), cau, mongDoi, nhipThoaiDanhBoss)
	}
	// Quay vong tren mang thoai: A, B, A, B...
	for i, c := range cau {
		if want := []string{"A", "B"}[i%2]; c != want {
			t.Fatalf("cau thu %d = %q, mong doi %q (quay vong tren ThoaiDanh)", i, c, want)
		}
	}
}

// TestBossKhongThoaiKhiKhongCoMucTieu: boss dung giua ban do trong ma lam bam
// mot minh la mot loi ai cung nhin thay. Thoai danh CHI ra khi dang danh nhau.
func TestBossKhongThoaiKhiKhongCoMucTieu(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k := khuVucAiBoss(&luc)
	// Nguoi choi o RAT XA, ngoai ca tam duoi.
	_, p := themNguoiAiBoss(k, 1, xBossThuNghiem+2000, yDiDuocBoss)
	b := bossChienDau(luc, xBossThuNghiem, 1500, []string{"A"})
	k.boss = append(k.boss, b)

	for i := 0; i < 300; i++ {
		luc = t0.Add(time.Duration(i) * 100 * time.Millisecond)
		k.nhipTim()
	}
	if cau := cauDaNoi(t, p); len(cau) != 0 {
		t.Fatalf("boss noi %v khi khong co ai trong tam duoi", cau)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// 5. Nan nhan chet
// ─────────────────────────────────────────────────────────────────────────────

// TestBossGietNguoiChoiDiDungDuongChet: khi don cua boss ha guc nguoi choi,
// phai di DUNG duong xuLyNguoiChoiChet da co (-17 rieng chinh chu, -8 cho
// nguoi khac) — khong duoc tu che mot duong bao chet thu hai.
func TestBossGietNguoiChoiDiDungDuongChet(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k := khuVucAiBoss(&luc)
	nanNhan, pNanNhan := themNguoiAiBoss(k, 1, xBossThuNghiem+50, yDiDuocBoss)
	nanNhan.nv.HpHienTai = 1 // mot don la di
	_, pNguoiKhac := themNguoiAiBoss(k, 2, xBossThuNghiem+60, yDiDuocBoss)
	b := bossChienDau(luc, xBossThuNghiem, 1500, nil)
	k.boss = append(k.boss, b)

	k.nhipTim()

	goi := locGoi(pNanNhan, lenhNguoiDanhNguoi)
	if len(goi) != 1 {
		t.Fatalf("phat %d goi -60, mong doi 1", len(goi))
	}
	_, _, sat, daChet := docDanh60(t, goi[0])
	if sat != 1 {
		t.Errorf("sat thuong = %d, mong doi 1 — khong duoc bao nhieu hon mau con lai", sat)
	}
	if !daChet {
		t.Errorf("goi -60 khong danh dau nan nhan da chet")
	}
	if nanNhan.nv.HpHienTai != 0 {
		t.Errorf("mau nan nhan = %d, mong doi 0", nanNhan.nv.HpHienTai)
	}
	if n := len(locGoi(pNanNhan, lenhNguoiChet)); n != 1 {
		t.Errorf("nan nhan nhan %d goi -17, mong doi 1", n)
	}
	if n := len(locGoi(pNguoiKhac, lenhAiDoChet)); n != 1 {
		t.Errorf("nguoi khac nhan %d goi -8, mong doi 1", n)
	}
}
