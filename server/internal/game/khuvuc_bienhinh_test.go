package game

import (
	"testing"
	"time"

	"nroserver/internal/protocol"
)

// ─────────────────────────────────────────────────────────────────────────────
// KH6b Task 11 — BIEN HINH QUA PHA (giu nguyen id) + NOI GIAN THEO %MAU.
//
// File nay khoa SAU dieu, moi dieu la mot cach hong ma KHONG bao mot loi nao:
//
//  1. Bien hinh phai GIU NGUYEN b.Id. Doi id thi con boss pha 2 la mot thuc the
//     la voi client: thanh mau cu treo lai, moi goi -60 dang bay toi id cu roi
//     vao khoang khong.
//  2. Thu tu goi PHAI la -6 TRUOC roi -5 SAU, trong CUNG mot luot actor. Nguoc
//     lai thi client (Controller.cs:4154) thay findCharInMap != nil nen BO goi
//     -5, roi -6 den sau xoa not — boss pha 2 bien mat khong dau vet.
//     (GATE task 1: docs/superpowers/research/2026-09-02-kh6b-kiem-client-bien-hinh.md)
//  3. Vao pha moi phai DAT LAI het trang thai chay cua pha cu: mau, nhip danh,
//     con tro noi gian, CA BA con tro thoai. Quen mot cai nao thi boss pha 2 ra
//     doi voi ky uc cua pha 1 — vd noi loi chao tu cau thu hai, hoac khong noi
//     gi ca.
//  4. Moi moc noi gian no DUNG MOT LAN. HeSoNhip nhan DON, nen no hai lan la
//     boss nhanh gap doi so du lieu ghi; no moi tick la nhip tut ve 0 sau vai
//     giay (boss danh 10 don/giay).
//  5. Don KET LIEU cham moc noi gian thi KHONG duoc phat cau noi gian: boss da
//     chet, da chia thuong, da vao TROI — mot cau "Thoi giau tay" bat ra sau do
//     la con ma noi chuyen.
//  6. HeSoDame ap o huong boss DANH ra, HeSoDef ap o huong boss NHAN vao. Dat
//     nham cho la boss vua manh gap ruoi vua trau gap ruoi.
//
// Khuon: k.nhipTim()/k.xuLy() goi TRUC TIEP tren khu vuc CHUA BatDau(), dong ho
// tiem duoc — dung khuon khuvuc_vongdoiboss_test.go / khuvuc_aiboss_test.go.
// ─────────────────────────────────────────────────────────────────────────────

// bossHaiPhaMau dung con boss mau THAT (noidungboss.go, "Kẻ Lạ Mặt" → "Mặt Nạ
// Vỡ") thay vi bia mot con hai pha trong test.
//
// VI SAO DUNG DU LIEU THAT: doan bien hinh la doan DUY NHAT trong ca he boss ma
// du lieu noi dung phai khop voi code o nhieu truong cung luc (Ten, Cap, Outfit,
// HpChon, NhipDanhMs, MocNoiGian cua pha 2). Mot con boss bia trong test se van
// xanh khi ai do lam hong pha 2 cua boss mau that.
func bossHaiPhaMau(t *testing.T, luc time.Time) *Boss {
	t.Helper()
	nd := TimNoiDungBoss(MaBossNhanGiaMatNa)
	if nd == nil {
		t.Fatalf("khong tim thay boss mau %q", MaBossNhanGiaMatNa)
	}
	if len(nd.Pha) < 2 {
		t.Fatalf("boss mau %q chi co %d pha — test bien hinh can 2", MaBossNhanGiaMatNa, len(nd.Pha))
	}
	hp := KepHp(nd.Pha[0].HpChon[0])
	return &Boss{
		Id:              -7,
		PhaHienTai:      0,
		DuLieu:          nd.Pha,
		Hp:              hp,
		HpToiDa:         hp,
		ViTriX:          100,
		ViTriY:          384,
		TrangThai:       BossHoatDong,
		MocThoiGian:     luc,
		LucHienPha:      luc,
		NhipDanhHienTai: nd.Pha[0].NhipDanhMs,
	}
}

// chiSoGoi tra vi tri (theo dung thu tu phat ra) cua moi goi mang ma lenh da
// cho. Khac locGoi: o day can BIET THU TU giua hai loai goi khac nhau, ma
// locGoi thi da loc mat.
func chiSoGoi(p *phienGiaKhuVuc, lenh int8) []int {
	var ra []int
	for i, m := range p.banTin() {
		if m.Command == lenh {
			ra = append(ra, i)
		}
	}
	return ra
}

// protocolReaderId doc truong id (int 4 byte dau) cua mot goi — dung cho goi -6.
func protocolReaderId(t *testing.T, m *protocol.Message) int32 {
	t.Helper()
	id, err := protocol.NewReader(m.Data).Int()
	if err != nil {
		t.Fatalf("doc id cua goi %d loi: %v", m.Command, err)
	}
	return id
}

// chayToiBienHinh dua con boss tu HOAT_DONG qua chet -> TROI -> RA_DI -> (bien
// hinh), roi tra ve so nhip da chay. Dung nhipTim that (khong goi tay tung buoc)
// de thu tu goi tin trong ket qua dung bang thu tu client se nhan.
func chayToiBienHinh(t *testing.T, k *KhuVuc, b *Boss, luc *time.Time) {
	t.Helper()
	// Don ket lieu: di qua bossNhanSatThuong — DUONG DUY NHAT tru mau boss.
	k.bossNhanSatThuong(b, 1, int64(b.Hp), *luc)
	if b.TrangThai != BossTroi {
		t.Fatalf("sau don ket lieu, trang thai = %v, mong doi TROI", b.TrangThai)
	}
	// Xac nam thoiGianXacBoss (3s) + con phai noi het loi trang troi. 60 nhip =
	// 6 giay, du rong cho ca hai.
	for i := 0; i < 60; i++ {
		*luc = luc.Add(100 * time.Millisecond)
		k.nhipTim()
		if b.PhaHienTai == 1 {
			return
		}
	}
	t.Fatalf("chay 60 nhip ma boss van chua sang pha 2 (PhaHienTai=%d, trang thai=%v)",
		b.PhaHienTai, b.TrangThai)
}

// ─────────────────────────────────────────────────────────────────────────────
// 1. Bien hinh — id, thu tu goi, ngoai hinh pha moi
// ─────────────────────────────────────────────────────────────────────────────

// TestBienHinhGiuNguyenIdVaGuiSauTruocNam la test DAU-CUOI cua ca task.
//
// Bon khang dinh, moi cai giet mot dot bien khac nhau:
//   - dung MOT goi -6 va MOT goi -5, va -6 di TRUOC (dot bien "gui -5 truoc -6")
//   - ca hai mang CUNG mot id (dot bien "doi id khi sang pha")
//   - goi -5 mo ta pha 2 that (ten/cap/outfit doi) — khong phai anh chup pha 1
//   - boss VAN NAM trong k.boss (khong bi nhipBoss go nhu con het pha)
func TestBienHinhGiuNguyenIdVaGuiSauTruocNam(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucBossGia(&luc)
	b := bossHaiPhaMau(t, luc)
	k.boss = append(k.boss, b)
	p.nhan = nil

	// GHI LAI id TRUOC khi bien hinh. So voi b.Id SAU do la mot cai bay that:
	// neu code doi id, ca hai ve cua phep so sanh doi theo nhau va test van
	// xanh — dot bien "doi id khi sang pha" song sot.
	idTruoc := b.Id

	chayToiBienHinh(t, k, b, &luc)

	if b.Id != idTruoc {
		t.Errorf("id boss doi tu %d thanh %d — bien hinh PHAI giu nguyen id", idTruoc, b.Id)
	}

	viTri6, viTri5 := chiSoGoi(p, lenhRoiKhuVuc), chiSoGoi(p, lenhXuatHien)
	if len(viTri6) != 1 || len(viTri5) != 1 {
		t.Fatalf("phat %d goi -6 va %d goi -5, mong doi dung 1 moi loai", len(viTri6), len(viTri5))
	}
	if viTri6[0] > viTri5[0] {
		t.Errorf("goi -6 o vi tri %d, goi -5 o vi tri %d — -6 PHAI di TRUOC, nguoc lai client bo goi -5",
			viTri6[0], viTri5[0])
	}

	g := docGoi5(t, p.banTin()[viTri5[0]])
	if g.id != int32(idTruoc) {
		t.Errorf("goi -5 pha 2 mang id %d, mong doi %d — bien hinh PHAI giu nguyen id", g.id, idTruoc)
	}
	if viTri6[0] < len(p.banTin()) {
		if g6 := p.banTin()[viTri6[0]]; g6 != nil {
			r := protocolReaderId(t, g6)
			if r != int32(idTruoc) {
				t.Errorf("goi -6 mang id %d, mong doi %d — hai goi phai cung mot id", r, idTruoc)
			}
		}
	}

	pha2 := b.DuLieu[1]
	if g.ten != pha2.Ten {
		t.Errorf("ten trong goi -5 = %q, mong doi %q (pha 2)", g.ten, pha2.Ten)
	}
	if g.cap != pha2.Cap {
		t.Errorf("cap trong goi -5 = %d, mong doi %d (pha 2)", g.cap, pha2.Cap)
	}
	if g.head != pha2.Outfit[0] || g.body != pha2.Outfit[1] || g.leg != pha2.Outfit[2] {
		t.Errorf("outfit trong goi -5 = (%d,%d,%d), mong doi (%d,%d,%d) — pha 2",
			g.head, g.body, g.leg, pha2.Outfit[0], pha2.Outfit[1], pha2.Outfit[2])
	}

	muonHp := KepHp(pha2.HpChon[0])
	if g.hp != muonHp || g.hpToiDa != muonHp {
		t.Errorf("goi -5 ghi hp=%d/%d, mong doi %d/%d — pha moi phai DAY MAU theo HpChon cua CHINH no",
			g.hp, g.hpToiDa, muonHp, muonHp)
	}
	if b.Hp != muonHp || b.HpToiDa != muonHp {
		t.Errorf("mau boss sau bien hinh = %d/%d, mong doi %d/%d", b.Hp, b.HpToiDa, muonHp, muonHp)
	}
	if len(k.boss) != 1 || k.boss[0] != b {
		t.Errorf("boss bi go khoi k.boss sau bien hinh — con pha thi phai o lai")
	}
	if b.TrangThai != BossHien && b.TrangThai != BossChao {
		t.Errorf("trang thai sau bien hinh = %v, mong doi HIEN (hoac CHAO o nhip ke tiep)", b.TrangThai)
	}
}

// TestBienHinhDatLaiTrangThaiChayCuaPhaCu khoa DANH SACH DAY DU cac truong phai
// dat lai — chinh danh sach ma task 9 de lai o cho noi.
//
// Ba con tro thoai la cho de quen nhat: chung khong lam gi hong ngay, chi lam
// boss pha 2 "nuot" may cau dau cua no. Test dat CA BA ve gia tri khac 0 truoc
// khi chet, roi doi ca ba ve 0.
func TestBienHinhDatLaiTrangThaiChayCuaPhaCu(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucBossGia(&luc)
	b := bossHaiPhaMau(t, luc)

	// Ban tay ban: pha 1 da chay het co — noi gian 2 lan, noi het ba mang thoai.
	b.MocNoiGianDaQua = 2
	b.HeSoDameHienTai = 1.6
	b.HeSoDefHienTai = 0.7
	b.ChiSoThoaiHien = 2
	b.ChiSoThoaiChet = 1
	b.ChiSoThoaiDanh = 3
	b.NhipDanhHienTai = 900 // 1500 da bi nhan don hai lan

	k.boss = append(k.boss, b)
	p.nhan = nil
	chayToiBienHinh(t, k, b, &luc)

	if b.MocNoiGianDaQua != 0 {
		t.Errorf("MocNoiGianDaQua = %d sau bien hinh, mong doi 0 — pha moi co bang moc RIENG",
			b.MocNoiGianDaQua)
	}
	if heSoHoacMot(b.HeSoDameHienTai) != 1 || heSoHoacMot(b.HeSoDefHienTai) != 1 {
		t.Errorf("he so dame/def sau bien hinh = %v/%v, mong doi 1/1 — mang he so pha cu sang lam boss pha cuoi manh gap boi",
			heSoHoacMot(b.HeSoDameHienTai), heSoHoacMot(b.HeSoDefHienTai))
	}
	if b.ChiSoThoaiHien != 0 || b.ChiSoThoaiChet != 0 || b.ChiSoThoaiDanh != 0 {
		t.Errorf("con tro thoai sau bien hinh = (hien %d, chet %d, danh %d), mong doi (0,0,0)",
			b.ChiSoThoaiHien, b.ChiSoThoaiChet, b.ChiSoThoaiDanh)
	}
	if muon := b.DuLieu[1].NhipDanhMs; b.NhipDanhHienTai != muon {
		t.Errorf("NhipDanhHienTai sau bien hinh = %d, mong doi %d (NhipDanhMs cua PHA MOI, khong phai nhip da nhan don o pha cu)",
			b.NhipDanhHienTai, muon)
	}
	if !b.LucHienPha.Equal(luc) {
		t.Errorf("LucHienPha = %v, mong doi %v — pha moi phai duoc khoang bat tu ngan cua chinh no",
			b.LucHienPha, luc)
	}
	if b.SatThuongTheoNguoi != nil {
		t.Errorf("so dong gop chua duoc xoa sau khi chia thuong pha 1 — pha 2 se chia lai cho nguoi cu")
	}
}

// TestBienHinhNoiLaiLoiChaoTuCauDauTien la dot bien cua rieng ChiSoThoaiHien:
// neu con tro khong ve 0, boss pha 2 se bo qua cau chao dau tien cua no.
func TestBienHinhNoiLaiLoiChaoTuCauDauTien(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucBossGia(&luc)
	b := bossHaiPhaMau(t, luc)
	b.ChiSoThoaiHien = len(b.DuLieu[0].ThoaiHien) // pha 1 da noi het loi chao
	k.boss = append(k.boss, b)
	p.nhan = nil

	chayToiBienHinh(t, k, b, &luc)
	p.nhan = nil // bo het goi cua doan chet + bien hinh, chi nghe pha 2 noi

	// Chay du lau cho ca hai cau chao cua pha 2 treo xong.
	for i := 0; i < 60; i++ {
		luc = luc.Add(100 * time.Millisecond)
		k.nhipTim()
	}
	cau := cauDaNoi(t, p)
	if len(cau) == 0 || cau[0] != b.DuLieu[1].ThoaiHien[0] {
		t.Fatalf("cau dau pha 2 noi = %q, mong doi %q — con tro ThoaiHien chua ve 0",
			cau, b.DuLieu[1].ThoaiHien[0])
	}
}

// TestHetPhaThiRoiHanTheGioi ghim rang nhanh CU khong bi task nay lam hong: con
// boss dang o pha CUOI van phai phat -6 va bien mat, khong duoc "bien hinh" ra
// mot pha khong ton tai.
func TestHetPhaThiRoiHanTheGioi(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucBossGia(&luc)
	b := bossHaiPhaMau(t, luc)
	b.PhaHienTai = 1 // dang o PHA CUOI
	b.Hp, b.HpToiDa = 100, 100
	k.boss = append(k.boss, b)
	p.nhan = nil

	k.bossNhanSatThuong(b, 1, 100, luc)
	for i := 0; i < 60; i++ {
		luc = luc.Add(100 * time.Millisecond)
		k.nhipTim()
	}
	if len(k.boss) != 0 {
		t.Errorf("boss het pha van con trong k.boss (%d con) — phai duoc go", len(k.boss))
	}
	if n := len(chiSoGoi(p, lenhRoiKhuVuc)); n != 1 {
		t.Errorf("phat %d goi -6, mong doi dung 1", n)
	}
	if n := len(chiSoGoi(p, lenhXuatHien)); n != 0 {
		t.Errorf("phat %d goi -5 khi HET pha — khong duoc phat cai nao", n)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// 2. Noi gian theo %mau
// ─────────────────────────────────────────────────────────────────────────────

// bossNoiGian dung mot boss MOT PHA co bang moc noi gian do nguoi goi dat.
//
// PHAI CO BoChieu du that ra khong don nao trung: thieu no thi aiBoss THOAT NGAY
// o dong dau ("pha khong co chieu -> khong co tam -> dung im"), tuc nhipDanhBoss
// khong bao gio chay — va moi khang dinh ve nhip trong file nay se xanh gia du
// ai do nhan HeSoNhip moi tick o trong do. Da vap dung cai bay nay mot lan khi
// chay kiem thu dot bien.
func bossNoiGian(luc time.Time, hp int32, moc []MocNoiGian) *Boss {
	return &Boss{
		Id:         -7,
		PhaHienTai: 0,
		DuLieu: []PhaBoss{{
			Ten:        "Kẻ Lạ Mặt",
			Lang:       0,
			Cap:        50,
			TypePk:     TypePkBoss,
			Dame:       60,
			NhipDanhMs: 1000,
			BoChieu:    [][2]int32{{1, 2}, {2, 2}}, // {can chien, cap 2}, {tam xa, cap 2}
			MocNoiGian: moc,
		}},
		Hp: hp, HpToiDa: hp,
		ViTriX: 100, ViTriY: 384,
		TrangThai:       BossHoatDong,
		MocThoiGian:     luc,
		LucHienPha:      luc,
		NhipDanhHienTai: 1000,
	}
}

// TestMocNoiGianNoDungMotLan la dot bien "moc noi gian phat 2 lan" VA dot bien
// "nhan HeSoNhip moi tick".
//
// Kich ban: boss 1000 mau, moc o 50%. Danh mot don 600 -> qua moc. Roi danh
// them hai don nua (van duoi 50%) va chay them nhip — cau chat chi duoc phat
// MOT lan, va nhip chi duoc nhan MOT lan.
func TestMocNoiGianNoDungMotLan(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucBossGia(&luc)
	b := bossNoiGian(luc, 1000, []MocNoiGian{
		{PhanTramHP: 50, HeSoDame: 1.3, HeSoDef: 0.85, HeSoNhip: 0.8, Chat: "Đủ rồi."},
	})
	k.boss = append(k.boss, b)
	p.nhan = nil

	k.bossNhanSatThuong(b, 1, 600, luc) // 1000 -> 400 = 40% < 50%
	k.bossNhanSatThuong(b, 1, 50, luc)  // van duoi moc
	k.bossNhanSatThuong(b, 1, 50, luc)

	cau := cauDaNoi(t, p)
	if len(cau) != 1 || cau[0] != "Đủ rồi." {
		t.Fatalf("boss noi %q, mong doi DUNG MOT cau %q — moc noi gian phai no dung mot lan",
			cau, "Đủ rồi.")
	}
	if b.MocNoiGianDaQua != 1 {
		t.Errorf("MocNoiGianDaQua = %d, mong doi 1", b.MocNoiGianDaQua)
	}
	if b.NhipDanhHienTai != 800 {
		t.Errorf("NhipDanhHienTai = %d, mong doi 800 (1000 x 0,8 DUNG MOT LAN)", b.NhipDanhHienTai)
	}
	if heSoHoacMot(b.HeSoDameHienTai) != 1.3 || heSoHoacMot(b.HeSoDefHienTai) != 0.85 {
		t.Errorf("he so dame/def = %v/%v, mong doi 1,3/0,85",
			heSoHoacMot(b.HeSoDameHienTai), heSoHoacMot(b.HeSoDefHienTai))
	}

	// Chay them 20 nhip: khong duoc no lai, va nhip KHONG duoc nhan them.
	for i := 0; i < 20; i++ {
		luc = luc.Add(100 * time.Millisecond)
		k.nhipTim()
	}
	if b.NhipDanhHienTai != 800 {
		t.Errorf("sau 20 nhip, NhipDanhHienTai = %d, mong doi van 800 — nhan HeSoNhip MOI TICK la cong don vo tan",
			b.NhipDanhHienTai)
	}
	if got := len(locGoi(p, lenhChat)); got != 1 {
		t.Errorf("tong so cau chat = %d, mong doi 1", got)
	}
}

// TestMocNoiGianKhongPhatKhiDonKetLieu: don giet boss cham moc noi gian thi moc
// do KHONG duoc no. Boss da chet, da chia thuong, da vao TROI — mot cau noi gian
// bat ra sau do la con ma noi chuyen.
func TestMocNoiGianKhongPhatKhiDonKetLieu(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucBossGia(&luc)
	b := bossNoiGian(luc, 1000, []MocNoiGian{
		{PhanTramHP: 90, HeSoDame: 2, HeSoDef: 0.5, HeSoNhip: 0.5, Chat: "Đủ rồi."},
	})
	k.boss = append(k.boss, b)
	p.nhan = nil

	k.bossNhanSatThuong(b, 1, 5000, luc) // mot don xuyen qua moc 90% xuong 0

	for _, c := range cauDaNoi(t, p) {
		if c == "Đủ rồi." {
			t.Fatal("don ket lieu van phat cau noi gian — boss da chet roi con noi")
		}
	}
	if b.MocNoiGianDaQua != 0 {
		t.Errorf("MocNoiGianDaQua = %d sau don ket lieu, mong doi 0", b.MocNoiGianDaQua)
	}
}

// TestMocNoiGianNhieuMocMotDon: mot don cuc manh vuot qua CA HAI moc cung luc
// thi ca hai he so deu phai duoc ap (khong bo sot moc), nhung chi noi MOT cau —
// hai cau chat trong cung mot khoanh khac se de len nhau, nguoi choi chi doc
// duoc cau cuoi.
func TestMocNoiGianNhieuMocMotDon(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	luc := t0
	k, p := khuVucBossGia(&luc)
	b := bossNoiGian(luc, 1000, []MocNoiGian{
		{PhanTramHP: 50, HeSoDame: 1.3, HeSoDef: 0.85, HeSoNhip: 0.8, Chat: "Moc mot"},
		{PhanTramHP: 20, HeSoDame: 1.6, HeSoDef: 0.70, HeSoNhip: 0.75, Chat: "Moc hai"},
	})
	k.boss = append(k.boss, b)
	p.nhan = nil

	k.bossNhanSatThuong(b, 1, 900, luc) // 1000 -> 100 = 10%, vuot ca hai moc

	if b.MocNoiGianDaQua != 2 {
		t.Errorf("MocNoiGianDaQua = %d, mong doi 2 — khong duoc bo sot moc nao", b.MocNoiGianDaQua)
	}
	if b.NhipDanhHienTai != 600 {
		t.Errorf("NhipDanhHienTai = %d, mong doi 600 (1000 x0,8 x0,75 — HeSoNhip nhan DON)", b.NhipDanhHienTai)
	}
	if heSoHoacMot(b.HeSoDameHienTai) != 1.6 {
		t.Errorf("HeSoDame = %v, mong doi 1,6 (gia tri TUYET DOI cua moc CUOI, khong nhan don)",
			heSoHoacMot(b.HeSoDameHienTai))
	}
	cau := cauDaNoi(t, p)
	if len(cau) != 1 || cau[0] != "Moc hai" {
		t.Errorf("boss noi %q, mong doi dung mot cau %q (cau cua moc CUOI vua vuot)", cau, "Moc hai")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// 3. Hai he so ap dung CHO
// ─────────────────────────────────────────────────────────────────────────────

// TestHeSoDameApVaoDonBossDanhRa: boss noi gian danh dau hon. Do bang cach chay
// HAI lan cung mot kich ban, chi khac HeSoDameHienTai.
func TestHeSoDameApVaoDonBossDanhRa(t *testing.T) {
	t0 := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	do := func(heSo float64) int32 {
		luc := t0
		k := khuVucAiBoss(&luc)
		nanNhan, _ := themNguoiAiBoss(k, 1, xBossThuNghiem+50, yDiDuocBoss)
		b := bossChienDau(luc, xBossThuNghiem, 1000, nil)
		b.HeSoDameHienTai = heSo
		k.boss = append(k.boss, b)
		mauTruoc := nanNhan.nv.HpHienTai
		k.nhipTim()
		return mauTruoc - nanNhan.nv.HpHienTai
	}
	thuong := do(0) // 0 = chua noi gian, phai duoc hieu la x1
	if thuong <= 0 {
		t.Fatalf("don thuong = %d, mong doi > 0 (he so 0 PHAI duoc hieu la x1, khong phai x0)", thuong)
	}
	noiGian := do(2)
	if noiGian != thuong*2 {
		t.Errorf("don khi noi gian x2 = %d, mong doi %d (gap doi don thuong %d)",
			noiGian, thuong*2, thuong)
	}
}

// TestHeSoDefApVaoDonBossNhanVao: boss noi gian trau hon. Don goc trong khuon
// khuvuc_danhboss_test.go la 2 (xem chu thich bossDanhDuoc) — he so 0,5 dua no
// ve 1.
//
// KHONG dung chung khuon voi test tren duoc: huong nguoi-danh-boss di qua
// k.xuLy(suKienDanhNguoiChoi), mot duong hoan toan khac.
func TestHeSoDefApVaoDonBossNhanVao(t *testing.T) {
	b := bossDanhDuoc(BossHoatDong, 1000, 0)
	b.HeSoDefHienTai = 0.5
	k, _, _, nvA := khuVucDanhBoss(t, b)

	k.xuLy(suKienDanhNguoiChoi{idNhanVat: nvA.ID, idMucTieu: b.Id})

	if b.Hp != 999 {
		t.Errorf("Hp boss sau don = %d, mong doi 999 (don goc 2 x HeSoDef 0,5 = 1)", b.Hp)
	}
}

// TestHeSoKhongDuocLamSatThuongVeKhong: he so nho (vd 0,01) van phai de lai it
// nhat 1 sat thuong khi don goc > 0. Mot boss "trau" ma khong the bi tru mau la
// mot boss bat tu — nguoi choi danh mai, thanh mau khong nhuc nhich.
func TestHeSoKhongDuocLamSatThuongVeKhong(t *testing.T) {
	if got := nhanHeSo(1, 0.01); got != 1 {
		t.Errorf("nhanHeSo(1, 0.01) = %d, mong doi 1 — don co that khong duoc lam tron ve 0", got)
	}
	if got := nhanHeSo(0, 5); got != 0 {
		t.Errorf("nhanHeSo(0, 5) = %d, mong doi 0 — khong co don thi nhan bao nhieu cung khong co", got)
	}
	if got := nhanHeSo(100, 0); got != 100 {
		t.Errorf("nhanHeSo(100, 0) = %d, mong doi 100 — he so 0 nghia la CHUA DAT, tuc x1", got)
	}
	if got := nhanHeSo(int32(satThuongToiDa), 10); got != int32(satThuongToiDa) {
		t.Errorf("nhanHeSo(tran, 10) = %d, mong doi kep o %d", got, satThuongToiDa)
	}
}
