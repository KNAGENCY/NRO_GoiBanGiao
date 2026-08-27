package game

import (
	"testing"

	"nroserver/internal/store"
)

// TestGoiMapInfoMap1CoCongRoiMap khoa tung byte cua o waypoint trong MAP INFO
// map 1. Gia tri ghim tay (§6.12).
//
// Truoc ban nay o day luon la `byte 0` — nghia la nguoi choi KHONG CO duong
// nao roi Lang La. Test nay la thu chan viec quay lai tinh trang do.
func TestGoiMapInfoMap1CoCongRoiMap(t *testing.T) {
	nv := &store.NhanVat{MapID: 1, ViTriX: 100, ViTriY: 384}
	m, err := goiMapInfo(nv, TimBanDo(1), nil)
	if err != nil {
		t.Fatal(err)
	}
	d := moDocGoi(t, m)

	d.b("mapId", 1)
	d.b("planetId", 0)
	d.b("tileId", 1)
	d.b("bgId", 0)
	d.b("type", 0)
	d.utf("mapName", "Làng Lá")
	d.b("zoneId", 0)
	d.sh("x", 100)
	d.sh("y", 384)

	d.b("waypointCount", 1)
	d.sh("wp.minX", 408)
	d.sh("wp.minY", 384)
	d.sh("wp.maxX", 472)
	d.sh("wp.maxY", 432)
	d.b("wp.laCong", 1)     // true: bam moi di, di ngang qua KHONG bi keo
	d.b("wp.ngoaiTuyen", 0) // false: khong dung he "map offline" cua ban cu
	d.utf("wp.ten", "Cổng Làng Lá")
}

// TestCongMap1DungTrongVungNguoiChoiToiDuoc doi chieu khung cong voi mot
// nguon DOC LAP — file assets/map-1-collision.bin. Mot cong dat vao vach da
// thi khong ai bam duoc, va loi do khong co trieu chung nao ngoai "bam mai
// khong di".
//
// Nhan vat dung tren mat be da o y = 408 (hang 17); hang 18 o cot 16..28 la
// DAC nen do la mat be. Vay diem giua khung o do y = 408 phai di duoc.
func TestCongMap1DungTrongVungNguoiChoiToiDuoc(t *testing.T) {
	bd := TimBanDo(1)
	if len(bd.DiemDi) != 1 {
		t.Fatalf("map 1 co %d cong, mong doi dung 1", len(bd.DiemDi))
	}
	c := bd.DiemDi[0]

	giuaX := (c.MinX + c.MaxX) / 2
	const yMatBeDa int16 = 408
	if !collisionLangLa.DiDuoc(giuaX, yMatBeDa) {
		t.Errorf("giua cong (%d, %d) khong di duoc — nguoi choi khong bao gio dung duoc trong cong",
			giuaX, yMatBeDa)
	}
	if !c.ChuaTrong(giuaX, yMatBeDa) {
		t.Errorf("khung cong khong bao duoc cho nguoi choi thuc su dung (%d, %d)", giuaX, yMatBeDa)
	}
	// Ca hai mep ngang cung phai la o di duoc, khong thi mot nua cong nam
	// trong da.
	for _, x := range []int16{c.MinX, c.MaxX} {
		if !collisionLangLa.DiDuoc(x, yMatBeDa) {
			t.Errorf("mep cong x=%d o y=%d nam trong vach da", x, yMatBeDa)
		}
	}
}

// TestCongDoiXungKhongNhotNguoiChoi: moi cong phai co duong VE. Day chinh la
// loi vua sua — Lang La khong co cong nao nen nguoi choi ket trong do. Them
// mot cong mot chieu sang map khac la day nguoi choi tu nha tu nay sang nha
// tu khac.
func TestCongDoiXungKhongNhotNguoiChoi(t *testing.T) {
	coDuongRa := map[int32]bool{}
	for _, bd := range DanhSachBanDo() {
		for _, c := range bd.DiemDi {
			coDuongRa[bd.MapID] = true
			if TimBanDo(c.DenMap) == nil {
				t.Errorf("cong %q cua map %d tro toi map %d khong co trong danh muc",
					c.Ten, bd.MapID, c.DenMap)
				continue
			}
			// Map dich phai co it nhat mot cong tro NGUOC lai.
			ve := false
			for _, cVe := range TimBanDo(c.DenMap).DiemDi {
				if cVe.DenMap == bd.MapID {
					ve = true
				}
			}
			if !ve {
				t.Errorf("cong %q dua nguoi choi tu map %d sang map %d, ma map %d khong co cong nao ve lai — nhot nguoi choi",
					c.Ten, bd.MapID, c.DenMap, c.DenMap)
			}
		}
	}
	if !coDuongRa[1] {
		t.Error("map 1 (Lang La) khong co cong nao — nguoi choi khong roi duoc lang khoi dau")
	}
}

// TestDiemHaCanhVeMap1DiDuoc: cho ha canh khi ve Lang La phai la o di duoc
// theo collision that. Ha canh vao vach da thi may chu se tu choi moi buoc
// di sau do (kiemDiChuyenHopLe do doan duong tu vi tri hien tai) — nhan vat
// dung im, khong hieu tai sao.
func TestDiemHaCanhVeMap1DiDuoc(t *testing.T) {
	for _, bd := range DanhSachBanDo() {
		for _, c := range bd.DiemDi {
			if c.DenMap != 1 {
				continue // chi map 1 co artifact collision de kiem
			}
			if !collisionLangLa.DiDuoc(c.DenX, c.DenY) {
				t.Errorf("cong %q (map %d) ha canh vao (%d,%d) — o do khong di duoc",
					c.Ten, bd.MapID, c.DenX, c.DenY)
			}
		}
	}
}

// TestChuaTrongKhepKinHaiDau: bien cua khung phai KHEP KIN, dung het nhu
// client (Char.cs:1663 `cx >= minX && cx <= maxX`). Lech mot don vi la nguoi
// choi bam duoc nut ma may chu bao "khong o trong cong" — vong lap vo nghia.
func TestChuaTrongKhepKinHaiDau(t *testing.T) {
	c := DiemDi{MinX: 10, MinY: 20, MaxX: 30, MaxY: 40}
	trong := [][2]int16{{10, 20}, {30, 40}, {10, 40}, {30, 20}, {20, 30}}
	ngoai := [][2]int16{{9, 30}, {31, 30}, {20, 19}, {20, 41}}
	for _, p := range trong {
		if !c.ChuaTrong(p[0], p[1]) {
			t.Errorf("(%d,%d) phai nam TRONG khung (bien khep kin)", p[0], p[1])
		}
	}
	for _, p := range ngoai {
		if c.ChuaTrong(p[0], p[1]) {
			t.Errorf("(%d,%d) phai nam NGOAI khung", p[0], p[1])
		}
	}
}

// TestMoiBanDoToiDaMotCong: xem chu thich ghiDiemDi (thegioi.go) — client tra
// cuu song song TileMap.vGo va PopUp.vPopups bang cung mot chi so, gia dinh do
// chi chac chan dung khi moi ban do co dung mot cong.
func TestMoiBanDoToiDaMotCong(t *testing.T) {
	for _, bd := range DanhSachBanDo() {
		if len(bd.DiemDi) > 1 {
			t.Errorf("map %d co %d cong — client (Char.isInEnterOnlinePoint) khop vGo voi vPopups bang cung chi so, nhieu cong la cho de vo am tham",
				bd.MapID, len(bd.DiemDi))
		}
	}
}
