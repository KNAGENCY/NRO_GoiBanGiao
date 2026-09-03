// lichboss_napnong_test.go — "BOSS DANG SONG KHONG DOI SO KHI NAP NONG"
// (KH8 task 7, thiet ke §4 "cai gi khong nap nong duoc").
//
// LOI HUA DANG DUOC KHOA:
//
//	con boss DANG DUNG trong the gioi giu so cua THE HE LUC NO SINH RA;
//	con boss SINH SAU lan nap nong moi lay so moi.
//
// Vi sao khong the de "nap nong doi luon ca boss dang song": HpToiDa va Ten cua
// boss DA DI XUONG CLIENT trong goi -5 luc no hien ra. Doi giua chung ma khong
// gui lai goi -5 la lam thanh mau va ten o client sai vinh vien — khong mot dong
// log nao. Va doi Dame giua tran la doi luat choi cua nhung nguoi dang danh do.
//
// CACH QUAN SAT: dem va doc goi -5 ma mot nguoi quan sat trong khu vuc nhan duoc
// (dung khuon lichboss_test.go) — do dung la thu nguoi choi that thay, nen test
// khoa cai dang duoc hua chu khong khoa mot chi tiet noi bo.
package game

import (
	"bytes"
	"testing"
	"time"

	"nroserver/internal/noidung"
	"nroserver/internal/protocol"
)

// goi5CuaPhien tra TAT CA goi -5 phien gia da nhan, theo thu tu.
func goi5CuaPhien(p *phienGiaKhuVuc) []*protocol.Message {
	var ra []*protocol.Message
	for _, m := range p.banTin() {
		if m.Command == lenhXuatHien {
			ra = append(ra, m)
		}
	}
	return ra
}

// tenMoiSauNapNong / dameMoiSauNapNong: hai con so cua "the he thu hai".
const (
	tenMoiSauNapNong  = "TEN-SAU-NAP-NONG"
	dameMoiSauNapNong = 12345
)

// duLieuTheHeHai nap kho noi dung that roi doi TEN + DAME cua pha dau con boss
// lon. Doi TEN vi ten la thu DI XUONG CLIENT trong goi -5 (quan sat duoc tu
// ngoai); doi DAME vi dame la thu bo dieu phoi doc lai (kiem duoc tren muc lich).
func duLieuTheHeHai(t *testing.T) *noidung.DuLieu {
	t.Helper()
	du, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("nap kho noi dung: %v", err)
	}
	thay := false
	for i := range du.Boss {
		if du.Boss[i].Ma != MaBossNhanGiaMatNa {
			continue
		}
		du.Boss[i].Pha[0].Ten = tenMoiSauNapNong
		du.Boss[i].Pha[0].Dame = dameMoiSauNapNong
		thay = true
	}
	if !thay {
		t.Fatalf("khong tim thay boss %q trong kho noi dung — dot bien khong va vao dau", MaBossNhanGiaMatNa)
	}
	return du
}

// TestLichBoss_BossDangSongKhongDoiSoKhiNapNong chay ca ba chang cua loi hua:
//
//  1. sinh lan mot   -> boss mang so cua THE HE 1
//  2. nap nong       -> KHONG co goi -5 nao moi, boss dang song khong doi gi
//  3. boss roi di, sinh lan hai -> con moi mang so cua THE HE 2
//
// Test nay DO neu lichboss quay ve doc bang Go cung (khong con bang nao), neu no
// khong doc lai kho luc sinh (chang 3 do), hoac neu no doc lai kho GIUA CHUNG va
// dung toi con boss dang song (chang 2 do).
func TestLichBoss_BossDangSongKhongDoiSoKhiNapNong(t *testing.T) {
	s := moSoKhuVucCoDon(t)
	k := s.Tra(1)
	p := datNguoiQuanSat(t, k, 1)

	gen1, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("nap kho noi dung: %v", err)
	}
	kho := noidung.MoKho(gen1)

	// Chi lay MOT con boss vao lich — con boss le kia se lam nhieu phep dem goi -5.
	var lich []MucLichBoss
	for _, nd := range danhSachNoiDungBossTu(kho.AnhChup()) {
		if nd.Ma == MaBossNhanGiaMatNa {
			lich = MucLichTuNoiDungBoss([]NoiDungBoss{nd})
		}
	}
	if len(lich) != 1 {
		t.Fatalf("dung duoc %d muc lich, mong doi 1", len(lich))
	}
	tenGoc := lich[0].DuLieu[0].Ten
	dameGoc := lich[0].DuLieu[0].Dame
	if tenGoc == tenMoiSauNapNong || dameGoc == dameMoiSauNapNong {
		t.Fatalf("du lieu goc da mang gia tri cua the he hai (%q/%d) — test khong chung minh duoc gi",
			tenGoc, dameGoc)
	}

	l := MoBoLichBoss(s, lich)
	l.chon = func(int) int { return 0 } // luon map dau, mau dau — khong phu thuoc may man
	// NGUON TIEM: doc tu `kho` (anh chup dang co hieu luc), khong phai kho nhung
	// toan cuc — de test dieu khien duoc "khi nao the he moi co hieu luc".
	l.nguonBoss = func(ma string) *NoiDungBoss { return timNoiDungBossTu(kho.AnhChup(), ma) }

	goc := time.Date(2026, 9, 3, 10, 0, 0, 0, time.UTC)
	l.muc[0].lanSinhCuoi = goc

	// ── Chang 1: sinh lan mot, the he 1 ───────────────────────────────────────
	l.nhip(goc.Add(2000 * time.Second)) // giay_nghi cua pha dau la 1800
	chotActor(t, k)
	goi := goi5CuaPhien(p)
	if len(goi) != 1 {
		t.Fatalf("chang 1: %d goi -5, mong doi dung 1", len(goi))
	}
	if !bytes.Contains(goi[0].Data, []byte(tenGoc)) {
		t.Fatalf("chang 1: goi -5 khong mang ten %q — boss khong sinh tu kho noi dung?", tenGoc)
	}

	// ── Chang 2: nap nong GIUA CHUNG ─────────────────────────────────────────
	kho.Thay(duLieuTheHeHai(t))

	// Nhip tiep: con boss the he 1 VAN dang song nen khong duoc sinh them con
	// nao, va TUYET DOI khong duoc phat lai goi -5 voi ten moi (do se la doi ten
	// mot con boss dang tren man hinh nguoi choi).
	l.nhip(goc.Add(4000 * time.Second))
	chotActor(t, k)
	goi = goi5CuaPhien(p)
	if len(goi) != 1 {
		t.Fatalf("chang 2: nap nong xong da co %d goi -5, mong doi van dung 1 — boss dang song "+
			"KHONG duoc sinh lai/doi so khi nap nong", len(goi))
	}
	if bytes.Contains(goi[0].Data, []byte(tenMoiSauNapNong)) {
		t.Fatal("chang 2: goi -5 CU bi doi sang ten moi — goi da phat di roi, khong the sua nguoc")
	}

	// ── Chang 3: boss cu roi khoi the gioi -> con sinh SAU lay so moi ────────
	l.muc[0].the.Tat() // dung soi day ma may trang thai dung khi boss ra di
	l.nhip(goc.Add(4000 * time.Second))
	chotActor(t, k)
	goi = goi5CuaPhien(p)
	if len(goi) != 2 {
		t.Fatalf("chang 3: %d goi -5, mong doi 2 (con thu hai phai duoc sinh)", len(goi))
	}
	if !bytes.Contains(goi[1].Data, []byte(tenMoiSauNapNong)) {
		t.Fatalf("chang 3: con boss sinh SAU nap nong van mang ten cu — bo dieu phoi khong doc lai " +
			"kho noi dung luc sinh, nap nong boss.yaml se khong bao gio co hieu luc")
	}
	if got := l.muc[0].DuLieu[0].Dame; got != dameMoiSauNapNong {
		t.Errorf("chang 3: muc lich van cam dame %d, mong doi %d", got, dameMoiSauNapNong)
	}
}

// TestLichBoss_MucKhongCoMaGiuNguyenDuLieuTuCam: hang chuc test KH6b dung muc
// lich TU CAM du lieu (mucMau, khong co Ma). Chung phai giu nguyen hanh vi cu —
// bo dieu phoi KHONG duoc di tra cuu kho noi dung cho nhung muc do.
func TestLichBoss_MucKhongCoMaGiuNguyenDuLieuTuCam(t *testing.T) {
	s := moSoKhuVucCoDon(t)
	k := s.Tra(1)
	p := datNguoiQuanSat(t, k, 1)

	l := MoBoLichBoss(s, []MucLichBoss{mucMau(30, []int32{1}, []int64{500})})
	l.chon = func(int) int { return 0 }
	l.nguonBoss = func(string) *NoiDungBoss {
		t.Error("bo dieu phoi di tra cuu kho noi dung cho mot muc KHONG CO ma")
		return nil
	}

	goc := time.Date(2026, 9, 3, 10, 0, 0, 0, time.UTC)
	l.muc[0].lanSinhCuoi = goc
	l.nhip(goc.Add(30 * time.Second))
	chotActor(t, k)
	if n := demGoi5(p); n != 1 {
		t.Fatalf("%d goi -5, mong doi 1 — muc khong co ma phai sinh boss tu DuLieu tu cam", n)
	}
	if l.muc[0].DuLieu[0].Ten != "CuuVi" {
		t.Errorf("du lieu tu cam bi ghi de: %q", l.muc[0].DuLieu[0].Ten)
	}
}

// TestLichBoss_TaoBossSaoDuLieuRiengKhongDungChungVoiMucLich: day la CO CHE
// khien chang 2 o tren dung duoc. Con boss da ban giao cho actor phai cam ban
// sao RIENG cua no; neu no dung chung lat cat voi muc lich thi mot lan doc lai
// kho (hay bat ky thao tac nao cua bo dieu phoi) se sua so cua con boss dang
// dung tren man hinh nguoi choi — va do la data race that su, vi hai luong khac
// nhau dang cham cung mot mang.
func TestLichBoss_TaoBossSaoDuLieuRiengKhongDungChungVoiMucLich(t *testing.T) {
	l := MoBoLichBoss(nil, []MucLichBoss{mucMau(30, []int32{1}, []int64{500})})
	l.chon = func(int) int { return 0 }
	m := l.muc[0]
	bayGio := time.Date(2026, 9, 3, 10, 0, 0, 0, time.UTC)

	b := l.taoBoss(m, m.phaDau(), &BanDo{MapID: 1, SpawnX: 10, SpawnY: 20}, bayGio)
	if len(b.DuLieu) == 0 {
		t.Fatal("boss vua tao khong co pha nao")
	}
	tenGoc := b.DuLieu[0].Ten

	// Bo dieu phoi doi du lieu cua muc (dung cai xay ra khi nap nong): con boss
	// DA BAN GIAO khong duoc doi theo.
	m.DuLieu[0].Ten = "BI-GHI-DE"
	if b.DuLieu[0].Ten != tenGoc {
		t.Errorf("boss dang song doi ten theo muc lich (%q) — taoBoss khong sao du lieu rieng",
			b.DuLieu[0].Ten)
	}

	// Va chieu nguoc lai: actor sua du lieu pha cua con boss (may trang thai lam
	// dung the qua PhaHienTaiData) khong duoc dung toi muc lich.
	b.DuLieu[0].Cap = 199
	if m.DuLieu[0].Cap == 199 {
		t.Error("actor sua pha cua boss da dung toi muc lich cua bo dieu phoi")
	}
}
