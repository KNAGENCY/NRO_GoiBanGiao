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

// TestTimBanDoSaoSAUCaDanhSachCong khoa dong `cp.DiemDi = append(...)` trong
// TimBanDo — cai bay ma chinh chu thich cua ham do canh bao ma truoc gio
// khong test nao chung minh.
//
// VI SAO NGUY HIEM neu mat: sao mot struct chi sao HEADER cua slice, ba mang
// (con tro/len/cap) van tro vao CUNG mot mang nen. Nguoi goi tuong minh dang
// cam mot ban rieng, sua mot cong — va sua that vao DANH MUC HE THONG, tu do
// MOI phien sau deu thay cong sai. Loi kieu nay khong co trieu chung tuc thi:
// no chi hien ra khi nguoi choi bam cong roi bi nem sang map khac.
func TestTimBanDoSaoSAUCaDanhSachCong(t *testing.T) {
	bd := TimBanDo(1)
	if len(bd.DiemDi) == 0 {
		t.Fatal("map 1 phai co it nhat mot cong de kiem duoc phep sao sau")
	}
	// Sua ban sao that nang tay: doi CA dich lan ten cong.
	bd.DiemDi[0].DenMap = 99
	bd.DiemDi[0].Ten = "cong bi sua bay"
	bd.DiemDi[0].MinX = -1

	lai := TimBanDo(1)
	if lai.DiemDi[0].DenMap == 99 || lai.DiemDi[0].Ten == "cong bi sua bay" || lai.DiemDi[0].MinX == -1 {
		t.Errorf("TimBanDo dung chung mang DiemDi voi danh muc goc: sua ban sao lam hong cong that (%+v)", lai.DiemDi[0])
	}
}

func TestCatalogMap6ChiLaTemplateChuaPlayable(t *testing.T) {
	if len(catalogLangLa.Maps) != 6 {
		t.Fatalf("catalog phải quảng bá đủ 6 template map, được %d", len(catalogLangLa.Maps))
	}
	map6 := catalogLangLa.Maps[5]
	if map6.ID != 6 || map6.Name != "Khu rừng" {
		t.Fatalf("template thứ 6 sai: %+v", map6)
	}
	if bd := TimBanDo(int32(map6.ID)); bd != nil {
		t.Fatalf("map 6 chưa có runtime tile/render contract nên chưa được playable: %+v", bd)
	}
	if got := len(DanhSachBanDo()); got != 5 {
		t.Fatalf("runtime chỉ được mở 5 map làng đã có render contract, được %d", got)
	}
}
