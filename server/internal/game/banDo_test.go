package game

import "testing"

func TestDanhMucBanDoDu5Lang(t *testing.T) {
	ds := DanhSachBanDo()
	if len(ds) != 5 {
		t.Fatalf("co %d ban do, mong doi 5 (mot map khoi dau moi lang)", len(ds))
	}
	// Ghim gia tri DOC LAP voi danhSachBanDo (khong doc tu chinh no) — bat duoc
	// neu ai do sua nham tileId/bgId/planetId; day la gia tri render THAT tu
	// map_template map cu 1-5, sai la client ve sai tileset hoac thoat app.
	type mong struct {
		mapID                          int32
		ten                            string
		planet, tile, bg, bgType, loai byte
	}
	muon := []mong{
		{1, "Làng Lá", 0, 1, 0, 0, 0},
		{2, "Làng Cát", 0, 1, 0, 0, 0},
		{3, "Làng Sương Mù", 0, 2, 1, 0, 0},
		{4, "Làng Đá", 0, 2, 1, 0, 0},
		{5, "Làng Mây", 0, 3, 2, 0, 0},
	}
	for i, w := range muon {
		g := ds[i]
		if g.MapID != w.mapID || g.Ten != w.ten || g.PlanetID != w.planet ||
			g.TileID != w.tile || g.BgID != w.bg || g.BgType != w.bgType || g.Loai != w.loai {
			t.Errorf("ban do %d = %+v, mong doi mapId=%d ten=%q planet=%d tile=%d bg=%d bgType=%d loai=%d",
				i, g, w.mapID, w.ten, w.planet, w.tile, w.bg, w.bgType, w.loai)
		}
	}
}

func TestTileIdMoiBanDoTrongKhoang(t *testing.T) {
	// tileId PHAI thuoc [1, numTileMap] = [1, 39]; ngoai khoang -> client thoat app.
	for _, bd := range DanhSachBanDo() {
		if bd.TileID < 1 || bd.TileID > 39 {
			t.Errorf("ban do mapId=%d co tileId=%d, phai thuoc [1,39]", bd.MapID, bd.TileID)
		}
		if bd.MapID < 0 || bd.MapID > 255 {
			t.Errorf("ban do mapId=%d ngoai [0,255] (client doc unsigned byte)", bd.MapID)
		}
	}
}

func TestTimBanDoTheoMapId(t *testing.T) {
	if bd := TimBanDo(3); bd == nil || bd.Ten != "Làng Sương Mù" || bd.TileID != 2 {
		t.Errorf("TimBanDo(3) = %+v, mong doi Lang Suong Mu tileId 2", bd)
	}
	if bd := TimBanDo(999); bd != nil {
		t.Errorf("TimBanDo(999) = %+v, mong doi nil (khong co)", bd)
	}
	// Tra ban SAO: sua ket qua khong duoc anh huong danh muc goc.
	bd := TimBanDo(1)
	bd.TileID = 99
	if lai := TimBanDo(1); lai.TileID != 1 {
		t.Errorf("TimBanDo tra con tro vao slice goc: sua ban sao lam doi danh muc (tileId=%d)", lai.TileID)
	}
}
