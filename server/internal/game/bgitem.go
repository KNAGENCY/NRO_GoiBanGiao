// Package game — bong do tren nen dat ("bg-item").
//
// TAI SAO CO FILE NAY: map 1 (Lang La) da co 8 anh bong do ve xong nhung
// khong dung duoc, vi may chu gui `short 0` o o bg-item trong MAP INFO.
// Client Unity ho tro san bg-item qua HAI goi khac nhau — thieu MOT trong
// hai la khong hien gi ca (va khong bao loi):
//
//	-31  BANG MAU (template) — anh nao, lop nao, lech bao nhieu. Gui MOT LAN
//	     mot phien, PHAI truoc -24. Client luu vao TileMap.vItemBg.
//	-24  DAT CHO — moi bong do la 3 short: chi so mau + o tile X + o tile Y.
//
// Client tra mau bang `TileMap.getBIById(id)` (Controller.cs:5207) — id
// chinh la CHI SO trong danh sach -31 (Controller.cs:3741 `bgItem.id =
// num177`, tuc bien dem vong lap). Tra khong thay thi client BO QUA IM
// LANG, khong ve gi, khong log. Do la ly do phai gui -31 truoc, va la ly do
// thu tu cua `danhSachMauBgItem` LA MOT PHAN CUA GIAO THUC — dao thu tu la
// doi anh cua moi bong do dang dat.
package game

import "nroserver/internal/protocol"

// lenhMauBgItem: goi XUONG -31 gui bang mau bg-item.
// Nguon: 1-MAY-CHU-GAME/src/data/DataGame.java:296-311 (sendDataItemBG) doi
// chieu D:\NRO_Project\Assets\Scripts\Game1\Controller.cs:3735-3756 (case -31).
const lenhMauBgItem int8 = -31

// LopBongDo = layer 2. DAY LA LUA CHON QUAN TRONG, khong phai so tuy y —
// layer quyet dinh CA thu tu ve LAN co bi nhuom mau hay khong:
//
//	layer 1: ve TRUOC tile nen (TileMap.cs:378 paintTilemap) -> bong do se bi
//	         chinh nen dat de len tren, coi nhu tang hinh. VA bi nhuom mau
//	         theo tileId (BgItemMn.blendImage).
//	layer 2: ve SAU tile nen, TRUOC npc/quai/nhan vat (GameScr.cs:4761 roi
//	         4775) -> dung y nghia "bong do tren mat dat, nguoi di de len".
//	         VA `changeColor()` (BgItem.cs) THOAT NGAY khi layer == 2 nen anh
//	         duoc ve NGUYEN BAN, khong bi nhuom — dung cho bong do vi bong do
//	         da tu co mau/alpha rieng.
//	layer 3: ve SAU nhan vat (GameScr.cs:4889) -> bong do se phu LEN nguoi
//	         choi, sai hoan toan.
//	layer 4: ve tren cung VA bi keo theo camera (transX = -cmx/2 + 100,
//	         BgItem.cs) -> troi noi, khong bam dat.
const LopBongDo byte = 2

// LechBongDoY: bong do ve tu MEP TREN cua o tile, ma cai bong (vung alpha
// thuc) nam o giua anh — do 8..22 tren tong 24 don vi (do bang alpha bbox cua
// 8 file bong_*.png: hang 32..87 tren 96 pixel muc x4). Keo len 6 don vi de
// tam bong nam ngay duoi chan cong trinh thay vi tut xuong giua vach dat.
const LechBongDoY int16 = -6

// MauBgItem la MOT dong trong bang mau -31.
//
// AnhID la so hieu FILE ANH phia client: client nap
// `res/x{zoom}/mapBackGround/{AnhID}.png` (GameCanvas.loadImage +
// Image.__createImage -> Resources.Load). No KHONG lien quan gi toi
// mapping-id.json (khong gian ten cua sprite/part) — day la khong gian ten
// rieng cua bg-item.
type MauBgItem struct {
	Ten   string // chi de doc code va bao loi; KHONG di xuong day
	AnhID int16
	Lop   byte
	Dx    int16
	Dy    int16
}

// AnhBongDoDau: moc dau cua dai so hieu anh danh cho bong do Lang La.
//
// TAI SAO 500 chu khong phai 0-7 (nho, gon hon): client co MOT bo nho dem
// anh bg-item tren dia theo khoa `x{zoom}bgItem{AnhID}` (Rms.loadRMS,
// Controller.cs:5242). Neu khoa da ton tai thi client dung ANH TRONG CACHE va
// KHONG BAO GIO doc thu muc Resources nua. Ban Ngoc Rong goc dung so hieu
// bg-item trong khoang ~0..350, nen may nguoi choi da tung choi server cu co
// the con cache o cac so do — cam vao day la lay nham anh cu, loi IM LANG
// dung nghia. 500+ la vung ban goc chua bao gio dung -> chac chan khong co
// cache cu.
//
// Kem theo: KHONG duoc gui goi -93 (sendBgItemVersion) voi dai so hieu nay.
// Client lam `newSmallVersion[AnhID]` (Controller.cs:5244) tren mang co do dai
// bang SO LUONG mau — AnhID 500 se ban ra ngoai mang. May chu Go khong gui
// -93 (grep ca repo), va dung them.
const AnhBongDoDau int16 = 500

// danhSachMauBgItem — 8 anh bong do da ve san o
// docs/my-thuat/map-1-lang-la/anh-da-tao/bong_*.png.
//
// THU TU LA MOT PHAN CUA GIAO THUC (xem dau file): chi so trong danh sach nay
// chinh la so hieu client dung de tra mau. Them mau moi thi THEM VAO CUOI.
//
// Tat ca deu Dx = 0: dat cho theo o tile 24 don vi da du sat, va Dx dung
// CHUNG cho moi cho dat cung mot mau nen keo ngang o day se lam lech cai kia.
//
// BANG CAM ANH VAO CLIENT (repo Unity KHONG nam trong git nay, nen bang duoi
// la ban ghi duy nhat cua viec cam). File nguon o muc x4; x1/x2/x3 thu nho
// theo ti le 1/4, 2/4, 3/4 (LANCZOS) — dung quy uoc cua
// tools/map-lang-la/lap-anh-dot-2.py. Dich:
// D:\NRO_Project\Assets\Resources\res\x{1..4}\mapBackGround\{AnhID}.png
//
//	500 <- bong_nha_nho_x4.png   576x96 px  (144 x 24 don vi game)
//	501 <- bong_nha_lon_x4.png   768x96     (192 x 24)
//	502 <- bong_quan_x4.png      768x96     (192 x 24)
//	503 <- bong_torii_x4.png     480x96     (120 x 24)
//	504 <- bong_gieng_x4.png     288x96     ( 72 x 24)
//	505 <- bong_cay_x4.png       576x96     (144 x 24)
//	506 <- bong_thung_x4.png     192x96     ( 48 x 24)
//	507 <- bong_den_x4.png       192x96     ( 48 x 24)
//
// Quy tac co anh (SO-TAY-LAM-MAP.md): px o muc x4 = don vi game x 4, vi
// Image.getWidth() ben client chia lai cho zoomLevel.
var danhSachMauBgItem = []MauBgItem{
	{Ten: "bong nha nho", AnhID: AnhBongDoDau + 0, Lop: LopBongDo, Dy: LechBongDoY},
	{Ten: "bong nha lon", AnhID: AnhBongDoDau + 1, Lop: LopBongDo, Dy: LechBongDoY},
	{Ten: "bong quan", AnhID: AnhBongDoDau + 2, Lop: LopBongDo, Dy: LechBongDoY},
	{Ten: "bong cong torii", AnhID: AnhBongDoDau + 3, Lop: LopBongDo, Dy: LechBongDoY},
	{Ten: "bong gieng", AnhID: AnhBongDoDau + 4, Lop: LopBongDo, Dy: LechBongDoY},
	{Ten: "bong cay lon", AnhID: AnhBongDoDau + 5, Lop: LopBongDo, Dy: LechBongDoY},
	{Ten: "bong thung go", AnhID: AnhBongDoDau + 6, Lop: LopBongDo, Dy: LechBongDoY},
	{Ten: "bong cot den", AnhID: AnhBongDoDau + 7, Lop: LopBongDo, Dy: LechBongDoY},
}

// Chi so mau — dat ten de bang dat cho o duoi doc duoc, va de doi thu tu
// danh sach tren khong lam sai bang duoi mot cach am tham.
const (
	mauBongNhaNho int16 = iota
	mauBongNhaLon
	mauBongQuan
	mauBongTorii
	mauBongGieng
	mauBongCay
	mauBongThung
	mauBongDen
)

// DatBgItem la MOT cho dat bong do trong mot ban do.
//
// TileX/TileY tinh theo O TILE, khong phai don vi game: client nhan lai 24
// (`bgItem.x = num8 * TileMap.size`, Controller.cs:5211, TileMap.size = 24).
type DatBgItem struct {
	Mau   int16
	TileX int16
	TileY int16
}

// Hai hang dat bong do cua map 1.
//
// Map 1 rong 52 o x 24 hang. Mat dat co HAI muc (doc tu
// assets/map-1-collision.bin):
//   - hang 18 di duoc o khap noi -> mat co o y = 432, bong do dat o hang 18.
//   - RIENG cot 16..28 hang 18 la DAC: do la be da nang cao (cong torii +
//     nha lon dung tren do), mat be o y = 408 -> bong do dat o hang 17.
//
// Dat sai hang thi bong do se noi giua khong khi hoac chim duoi dat — khong
// co canh bao nao ca.
const (
	hangDatBongDo  int16 = 18 // nen dat thap (y = 432)
	hangBeDaBongDo int16 = 17 // mat be da cot 16..28 (y = 408)
	cotBeDaDau     int16 = 16
	cotBeDaCuoi    int16 = 28
)

// bgItemMap1: 11 cho dat, do bang mat tren anh dung tu
// tools/map-lang-la/tiled/xem-truoc.py --game (1 pixel = 1 don vi game).
// Cot ghi trong chu thich la khoang MAT TIEN cua cong trinh do duoc.
//
// MOI MAU chi duoc dung TOI DA 2 lan: client bat dau lat guong (trans = 2)
// khi mot mau xuat hien QUA 2 lan (TileMap.isExistMoreOne, TileMap.cs:177).
// Bong do doi xung nen lat guong khong nhin thay, nhung giu <= 2 de bang nay
// khong phu thuoc vao mot chi tiet mo cua client.
var bgItemMap1 = []DatBgItem{
	{Mau: mauBongNhaNho, TileX: 2, TileY: hangDatBongDo},   // nha nho trai, mat tien x 50..177
	{Mau: mauBongGieng, TileX: 8, TileY: hangDatBongDo},    // gieng trai, x 200..245
	{Mau: mauBongQuan, TileX: 10, TileY: hangDatBongDo},    // quan trai, x 230..383
	{Mau: mauBongTorii, TileX: 16, TileY: hangBeDaBongDo},  // cong torii tren be da, x 388..495
	{Mau: mauBongThung, TileX: 19, TileY: hangBeDaBongDo},  // thung go tren be da, x 455..505
	{Mau: mauBongNhaLon, TileX: 21, TileY: hangBeDaBongDo}, // nha lon 2 tang tren be da, x 505..695
	{Mau: mauBongGieng, TileX: 29, TileY: hangDatBongDo},   // gieng phai, x 701..762
	{Mau: mauBongDen, TileX: 31, TileY: hangDatBongDo},     // cot den phai, x 759..784
	{Mau: mauBongQuan, TileX: 32, TileY: hangDatBongDo},    // quan phai, x 790..950
	{Mau: mauBongCay, TileX: 40, TileY: hangDatBongDo},     // cay lon, goc x 1000..1070
	{Mau: mauBongNhaNho, TileX: 44, TileY: hangDatBongDo},  // nha nho phai, x 1065..1195
}

// bgItemTheoMap: map nao co bong do. Map 2-5 chua co my thuat rieng nen chua
// co bong do — tra danh sach rong, MAP INFO van ghi `short 0` nhu cu.
var bgItemTheoMap = map[int32][]DatBgItem{
	1: bgItemMap1,
}

// goiMauBgItem dung goi -31: bang mau bg-item.
//
// KHUON BYTE (DataGame.java:299-307 == Controller.cs:3737-3754):
//
//	short  soMau
//	moi mau: short AnhID | byte Lop | short Dx | short Dy | byte soO
//
// `soO` LUON BANG 0 va day KHONG phai su luoi bieng: ban Java cung ghi
// `writeByte(0)` (DataGame.java:306), va client co MOT LOI THAT o vong doc
// o do — no ghi vao `bgItem.tileX[num177]` (chi so vong NGOAI) thay vi
// `[num178]` (Controller.cs:3751). Gui soO > 0 se lam client nem
// IndexOutOfRange giua goi -31 va MAT SACH cac mau dang sau. Hai truong
// tileX/tileY do khong duoc dung o bat cu dau khac trong client (da grep).
func goiMauBgItem() *protocol.Message {
	w := protocol.NewWriter()
	w.Short(int16(len(danhSachMauBgItem)))
	for _, m := range danhSachMauBgItem {
		w.Short(m.AnhID)
		w.Byte(m.Lop)
		w.Short(m.Dx)
		w.Short(m.Dy)
		w.Byte(0) // soO = 0 — DUNG DOI, xem chu thich ham
	}
	return protocol.NewMessage(lenhMauBgItem, w)
}

// ghiBgItem ghi o bg-item cua MAP INFO (-24).
//
// KHUON BYTE (Zone.java:838-842 ghi nguyen file data/map/item_bg_map_data/
// <mapId>, doi chieu Controller.cs:5199-5205 la noi doc):
//
//	short soBong
//	moi bong: short chiSoMau | short tileX | short tileY
//
// Map khong co bong do van PHAI ghi `short 0` (2 byte) — bo trang la client
// doc lech sang o effMap phia sau (Zone.java bat exception cung ghi short 0).
func ghiBgItem(w *protocol.Writer, bd *BanDo) {
	if bd == nil {
		w.Short(0)
		return
	}
	ds := bgItemTheoMap[bd.MapID]
	w.Short(int16(len(ds)))
	for _, d := range ds {
		w.Short(d.Mau)
		w.Short(d.TileX)
		w.Short(d.TileY)
	}
}
