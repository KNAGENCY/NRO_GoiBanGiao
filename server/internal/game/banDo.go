package game

// BanDo la MOT ban do lang Naruto. Ke hoach 4 DUNG LAI mapId cu (1-5) de client
// Unity hien tai co san do hoa + tile_map_data (hinh anh Naruto se thay khi dung
// lai client). Cac truong render (PlanetID/TileID/BgID/BgType) lay DUNG tu du lieu
// map cu that (bang map_template trong 1-MAY-CHU-GAME/z.sql, map 1-5) de client ve
// dung tileset co san — tileId ngoai [1,numTileMap]=[1,39] lam client THOAT HAN app
// (giao thuc §7).
type BanDo struct {
	MapID    int32  // dung lai mapId cu (client co san). Client doc unsigned byte -> <=255.
	Ten      string // ten hien thi trong MAP INFO (dat theo lang Naruto — nhan UI)
	PlanetID byte
	TileID   byte // ∈ [1,39] cua tile_set_info (§7)
	BgID     byte
	BgType   byte
	Loai     byte  // map type: 0 = MAP_NORMAL
	SpawnX   int16 // vi tri khoi dau mac dinh — hien khop ViTriKhoiDau trong lang.go
	SpawnY   int16 // (nam tren nen dat cua map). Ke hoach 5 dung khi hoi sinh.
}

// Catalog quảng bá 6 template để giữ mảng client dày đặc, nhưng template thứ 6
// "Khu rừng" chưa có runtime tile/render/collision contract nên cố ý chưa
// playable: TimBanDo(6) phải trả nil. danhSachBanDo chỉ chứa 5 map khởi đầu,
// mỗi làng (gender 0-4) một map. mapId/tileId/bgId/
// planetId/bgType lay DUNG tu map_template map cu 1-5 (Doi hoa cuc/Thung lung tre/
// Rung nam/Rung xuong/Dao Kame) de client ve dung do hoa co san. Ten map dat theo
// lang Naruto cho nhap vai (chi la nhan UI trong MAP INFO — coordinator co the
// tinh chinh). Ke hoach 6 se chuyen bang nay sang file nap nong (nhu ghi chu o
// lang.go). Gio de trong code cho don gian.
var danhSachBanDo = func() []BanDo {
	maps := catalogLangLa.Maps
	return []BanDo{
		{MapID: int32(maps[0].ID), Ten: maps[0].Name, PlanetID: 0, TileID: 1, BgID: 0, BgType: 0, Loai: 0, SpawnX: maps[0].SpawnX, SpawnY: maps[0].SpawnY},
		{MapID: int32(maps[1].ID), Ten: maps[1].Name, PlanetID: 0, TileID: 1, BgID: 0, BgType: 0, Loai: 0, SpawnX: maps[1].SpawnX, SpawnY: maps[1].SpawnY},
		{MapID: int32(maps[2].ID), Ten: maps[2].Name, PlanetID: 0, TileID: 2, BgID: 1, BgType: 0, Loai: 0, SpawnX: maps[2].SpawnX, SpawnY: maps[2].SpawnY},
		{MapID: int32(maps[3].ID), Ten: maps[3].Name, PlanetID: 0, TileID: 2, BgID: 1, BgType: 0, Loai: 0, SpawnX: maps[3].SpawnX, SpawnY: maps[3].SpawnY},
		{MapID: int32(maps[4].ID), Ten: maps[4].Name, PlanetID: 0, TileID: 3, BgID: 2, BgType: 0, Loai: 0, SpawnX: maps[4].SpawnX, SpawnY: maps[4].SpawnY},
	}
}()

// TimBanDo tra ban do theo mapId, hoac nil neu khong co. Tra ban SAO (cung ly do
// nhu TimLang trong lang.go: nguoi goi lo tay sua khong dung toi danh muc he thong).
func TimBanDo(mapID int32) *BanDo {
	for i := range danhSachBanDo {
		if danhSachBanDo[i].MapID == mapID {
			cp := danhSachBanDo[i]
			return &cp
		}
	}
	return nil
}

// DanhSachBanDo tra ban sao danh sach ban do (chi doc).
func DanhSachBanDo() []BanDo {
	ds := make([]BanDo, len(danhSachBanDo))
	copy(ds, danhSachBanDo)
	return ds
}
