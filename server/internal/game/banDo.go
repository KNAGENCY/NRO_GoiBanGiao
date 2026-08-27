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

	// DiemDi: cac cong roi ban do cua map nay (xem kieu DiemDi ben duoi).
	// Rong = map kin, nguoi choi khong co duong ra — do CHINH LA tinh trang
	// cua ca 5 map truoc ban va nay.
	DiemDi []DiemDi
}

// DiemDi la MOT cong roi ban do ("waypoint" ben client).
//
// KHUON BYTE gui xuong client (Zone.java:753-767 == Controller.cs:5055):
//
//	short MinX | short MinY | short MaxX | short MaxY | bool LaCong |
//	bool NgoaiTuyen | UTF Ten
//
// DenMap/DenX/DenY KHONG di xuong day — chung la BI MAT CUA MAY CHU
// (WayPoint.java co goMap/goX/goY nhung Zone.mapInfo khong ghi chung).
// Client chi biet "co mot cong o day, ten the nay"; khi nguoi choi kich hoat
// no thi client gui goi -23 tay khong, va may chu tu tra ra nguoi do dang
// dung trong cong nao (xem doimap.go). Nho vay client khong the tu chon di
// dau — mot nguoi sua client khong the nhay sang map bat ky.
//
// LaCong (isEnter) — QUYET DINH CACH KICH HOAT, doc rat ky truoc khi doi:
//
//	true : chi doi map khi nguoi choi TU BAM vao bang ten cong
//	       (Waypoint.perform idAction=2, Waypoint.cs). Di ngang qua KHONG
//	       bi hut di.
//	false: client TU DONG doi map ngay khi nhan vat buoc vao khung
//	       (Char.isInWaypoint -> requestChangeMap, Char.cs:2322 va 3787).
//	       Chi hop voi cong nam sat MEP ban do; dat giua ban do la nguoi
//	       choi di ngang la bi keo di.
//
// NgoaiTuyen (isOffline) doi map "offline" cua ban Ngoc Rong cu (goi -33),
// he thong do KHONG duoc port sang ban Naruto — luon de false.
type DiemDi struct {
	MinX, MinY int16
	MaxX, MaxY int16
	LaCong     bool
	NgoaiTuyen bool
	Ten        string

	DenMap int32
	DenX   int16
	DenY   int16
}

// Chua trong: DiemDi kiem nhan vat co dung trong khung khong. Bien la KHEP
// KIN ca hai dau, dung het nhu client (Char.cs:1664 `cx >= minX && cx <=
// maxX`) va nhu ban Java (MapService.isPlayerInWaypoint) — lech mot don vi o
// day la nguoi choi bam duoc nut nhung may chu bao "khong o trong cong".
func (d DiemDi) ChuaTrong(x, y int16) bool {
	return x >= d.MinX && x <= d.MaxX && y >= d.MinY && y <= d.MaxY
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
		{MapID: int32(maps[0].ID), Ten: maps[0].Name, PlanetID: 0, TileID: 1, BgID: 0, BgType: 0, Loai: 0, SpawnX: maps[0].SpawnX, SpawnY: maps[0].SpawnY, DiemDi: diemDiMap1},
		{MapID: int32(maps[1].ID), Ten: maps[1].Name, PlanetID: 0, TileID: 1, BgID: 0, BgType: 0, Loai: 0, SpawnX: maps[1].SpawnX, SpawnY: maps[1].SpawnY, DiemDi: diemDiMap2},
		{MapID: int32(maps[2].ID), Ten: maps[2].Name, PlanetID: 0, TileID: 2, BgID: 1, BgType: 0, Loai: 0, SpawnX: maps[2].SpawnX, SpawnY: maps[2].SpawnY},
		{MapID: int32(maps[3].ID), Ten: maps[3].Name, PlanetID: 0, TileID: 2, BgID: 1, BgType: 0, Loai: 0, SpawnX: maps[3].SpawnX, SpawnY: maps[3].SpawnY},
		{MapID: int32(maps[4].ID), Ten: maps[4].Name, PlanetID: 0, TileID: 3, BgID: 2, BgType: 0, Loai: 0, SpawnX: maps[4].SpawnX, SpawnY: maps[4].SpawnY},
	}
}()

// diemDiMap1 — CONG TORII cua Lang La.
//
// VI SAO DAT O TORII chu khong o mep ban do (quy uoc thuong thay cua the
// loai): torii la mot CAI CONG THAT trong my thuat da ve (nhin thay ro o
// tools/map-lang-la/tiled/xem-truoc.py --game, x 388..495) — nguoi choi hieu
// ngay do la loi ra. Hai mep ban do thi deu bi nha o de: nha nho trai chiem
// x 50..177, nha nho phai chiem x 1065..1195, dat cong o do thi cai bang ten
// se nam de len mai nha.
//
// VI SAO O DAY DI DUOC: torii dung tren be da (cot 16..28 cua
// assets/map-1-collision.bin). O do hang 18 la DAC nen nhan vat dung tren
// MAT BE, y = 408 (hang 17, di duoc). Khung 384..432 theo truc y bao tron ca
// y = 408 do va con du cho nhan vat nhun/nhay mot chut.
//
// LaCong = true (bam moi di) chu khong phai false (buoc vao la di) VI cong
// nay nam GIUA ban do: de false thi ai di ngang qua torii cung bi keo sang
// map khac, khong the di lai binh thuong trong lang duoc nua.
//
// Dich = map 2 (Lang Cat) tai diem xuat phat cua chinh no (100, 384). Trong
// 5 map hien co, moi map la mot lang, khong co map "duong noi" nao ca —
// nen cong chi co the noi thang sang mot lang khac. Chon map 2 vi no la map
// ke tiep theo so, va PHAI co cong quay ve (diemDiMap2) neu khong nguoi
// choi lai ket o Lang Cat dung cai loi vua sua o day.
var diemDiMap1 = []DiemDi{{
	MinX: 408, MinY: 384, MaxX: 472, MaxY: 432,
	LaCong: true, NgoaiTuyen: false,
	Ten:    "Cổng Làng Lá",
	DenMap: 2, DenX: 100, DenY: 384,
}}

// diemDiMap2 — cong quay ve Lang La, dat trum len chinh diem xuat phat cua
// Lang Cat (100, 384) de nguoi choi vua toi noi la thay ngay duong ve.
//
// Lang Cat chua co my thuat/collision rieng (chua co artifact) nen no dung
// bien du phong (kiemDiChuyenHopLe trong dichuyen.go) — khung nay chi can
// nam trong bien do, khong co vach da nao de tranh.
//
// Diem ha canh khi ve Lang La la (441, 408): giua hai cot torii, tren mat be
// da — collisionLangLa.DiDuoc(441, 408) = true (cot 18, hang 17).
var diemDiMap2 = []DiemDi{{
	MinX: 48, MinY: 336, MaxX: 168, MaxY: 456,
	LaCong: true, NgoaiTuyen: false,
	Ten:    "Về Làng Lá",
	DenMap: 1, DenX: 441, DenY: 408,
}}

// TimBanDo tra ban do theo mapId, hoac nil neu khong co. Tra ban SAO (cung ly do
// nhu TimLang trong lang.go: nguoi goi lo tay sua khong dung toi danh muc he thong).
func TimBanDo(mapID int32) *BanDo {
	for i := range danhSachBanDo {
		if danhSachBanDo[i].MapID == mapID {
			cp := danhSachBanDo[i]
			// Sao CA slice DiemDi: sao mot struct chi sao con tro cua slice,
			// nguoi goi van sua duoc bang cong goc qua ban "sao chep" — dung
			// cai bay ma dong tren dinh cam.
			cp.DiemDi = append([]DiemDi(nil), danhSachBanDo[i].DiemDi...)
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
