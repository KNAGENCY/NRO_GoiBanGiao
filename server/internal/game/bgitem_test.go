package game

import (
	"testing"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// docGoi gom ba ham doc quen thuoc lai mot cho, de cac test byte-chinh-xac
// duoi day khong phai lap ba closure giong het nhau.
type docGoi struct {
	t *testing.T
	r *protocol.Reader
}

func moDocGoi(t *testing.T, m *protocol.Message) *docGoi {
	t.Helper()
	return &docGoi{t: t, r: protocol.NewReader(m.Data)}
}

func (d *docGoi) b(ten string, muon byte) {
	d.t.Helper()
	g, err := d.r.Byte()
	if err != nil || g != muon {
		d.t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
	}
}

func (d *docGoi) sh(ten string, muon int16) {
	d.t.Helper()
	g, err := d.r.Short()
	if err != nil || g != muon {
		d.t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
	}
}

func (d *docGoi) utf(ten string, muon string) {
	d.t.Helper()
	g, err := d.r.UTF()
	if err != nil || g != muon {
		d.t.Fatalf("%s = %q (loi %v), mong doi %q", ten, g, err, muon)
	}
}

func (d *docGoi) hetGoi(ten string) {
	d.t.Helper()
	if _, err := d.r.Byte(); err == nil {
		d.t.Errorf("%s: con byte thua — goi phai het o day", ten)
	}
}

// TestGoiMauBgItemByteDung khoa TUNG BYTE cua goi -31.
//
// Gia tri o day GHIM BANG TAY (khong doc tu danhSachMauBgItem) — dung luat
// §6.12 cua BAT-DAU-TU-DAY: test lay ky vong tu chinh nguon minh dang kiem
// thi doi gi cung "dung", thanh diem mu.
func TestGoiMauBgItemByteDung(t *testing.T) {
	m := goiMauBgItem()
	if m.Command != -31 {
		t.Fatalf("lenh = %d, mong doi -31", m.Command)
	}
	d := moDocGoi(t, m)

	d.sh("soMau", 8)
	for i := 0; i < 8; i++ {
		d.sh("anhId", int16(500+i))
		d.b("lop", 2)  // layer 2: ve sau nen dat, truoc nhan vat, KHONG nhuom mau
		d.sh("dx", 0)  //
		d.sh("dy", -6) // keo bong len sat chan cong trinh
		d.b("soO", 0)  // BAT BUOC 0 — client co loi chi so o vong doc o nay
	}
	d.hetGoi("-31")
}

// TestMauBgItemKhongDinhBayCuaClient khoa lai cac gia tri ma client Unity
// doi xu DAC BIET. Day khong phai test "code lam dung nhu code" — moi so
// trong danh sach duoi deu doc ra tu mot dong `if` cu the trong client, va
// neu ai do doi AnhBongDoDau sang mot dai khac thi test nay la thu duy nhat
// bao truoc.
func TestMauBgItemKhongDinhBayCuaClient(t *testing.T) {
	// BgItem.paint (BgItem.cs) bo qua han 96 va 279; ve them nuoc chay cho 11;
	// keo camera cho 28 va 67..70. BgItem.isMiniBgz keo camera cho 79,80,81,
	// 85,86,90,91,92 va 99..108.
	bay := map[int16]string{
		11: "ve them hieu ung nuoc chay",
		28: "bi keo theo camera (transX)",
		67: "bi keo theo camera", 68: "bi keo theo camera",
		69: "bi keo theo camera", 70: "bi keo theo camera",
		96: "bi client bo qua han", 279: "bi client an khi hieu ung Mabu",
		79: "isMiniBg — bi keo theo camera", 80: "isMiniBg", 81: "isMiniBg",
		85: "isMiniBg", 86: "isMiniBg", 90: "isMiniBg", 91: "isMiniBg", 92: "isMiniBg",
		99: "isMiniBg", 100: "isMiniBg", 101: "isMiniBg", 102: "isMiniBg",
		103: "isMiniBg", 104: "isMiniBg", 105: "isMiniBg", 106: "isMiniBg",
		107: "isMiniBg", 108: "isMiniBg",
	}
	daDung := map[int16]bool{}
	for i, mau := range danhSachMauBgItem {
		if ly, dinh := bay[mau.AnhID]; dinh {
			t.Errorf("mau %d (%s) dung anhId %d — client doi xu dac biet voi so nay: %s",
				i, mau.Ten, mau.AnhID, ly)
		}
		if mau.AnhID == -1 {
			t.Errorf("mau %d dung anhId -1 — client dung -1 lam co 'khong ve'", i)
		}
		if daDung[mau.AnhID] {
			t.Errorf("anhId %d dung hai lan — moi mau phai co file anh rieng", mau.AnhID)
		}
		daDung[mau.AnhID] = true
		// >= 500: ban Ngoc Rong goc dung so hieu bg-item trong khoang ~0..350.
		// Client co bo nho dem anh tren dia theo khoa `x{zoom}bgItem{AnhID}`
		// (Rms.loadRMS) va KHONG kiem lai phien ban khi may chu khong gui -93 —
		// tai su dung mot so hieu cu la lay nham anh Dragon Ball trong cache cua
		// nguoi tung choi server cu, mot loi hoan toan im lang.
		if mau.AnhID < 500 {
			t.Errorf("mau %d (%s) dung anhId %d < 500 — de trung cache RMS cua ban Ngoc Rong cu",
				i, mau.Ten, mau.AnhID)
		}
		if mau.Lop != 2 {
			t.Errorf("mau %d (%s) o lop %d — bong do phai o lop 2 (sau nen dat, truoc nhan vat, khong bi nhuom mau)",
				i, mau.Ten, mau.Lop)
		}
	}
	if len(danhSachMauBgItem) != 8 {
		t.Errorf("co %d mau, mong doi 8 (dung 8 file bong_*.png da ve)", len(danhSachMauBgItem))
	}
}

// rongMauTheoO la CHIEU RONG (tinh bang O TILE 24 don vi) cua tung file anh
// bong do, GHIM BANG TAY theo dung luat §6.12 — day la kich thuoc THAT cua
// cac file docs/my-thuat/map-1-lang-la/anh-da-tao/bong_*.png o muc x4, chia
// cho 4 roi chia cho 24:
//
//	500 bong nha nho  576px x4 -> 144 don vi -> 6 o
//	501 bong nha lon  768px     -> 192       -> 8 o
//	502 bong quan     768px     -> 192       -> 8 o
//	503 bong torii    480px     -> 120       -> 5 o
//	504 bong gieng    288px     ->  72       -> 3 o
//	505 bong cay      576px     -> 144       -> 6 o
//	506 bong thung    192px     ->  48       -> 2 o
//	507 bong den      192px     ->  48       -> 2 o
//
// VI SAO may chu KHONG giu con so nay trong bgitem.go: giao thuc khong gui
// chieu rong xuong client (client tu do kich thuoc file anh), nen mot truong
// nhu vay se la du lieu chet ma khong ai buoc phai giu dung. O day thi nguoc
// lai: no LA ky vong doc lap ma test dung de kiem bang dat cho.
var rongMauTheoO = map[int16]int16{
	mauBongNhaNho: 6,
	mauBongNhaLon: 8,
	mauBongQuan:   8,
	mauBongTorii:  5,
	mauBongGieng:  3,
	mauBongCay:    6,
	mauBongThung:  2,
	mauBongDen:    2,
}

// TestBongDoMap1NamDungTrenMatDat la test co GIA TRI THAT SU o day: no doi
// chieu bang dat cho bong do voi mot nguon DOC LAP — file
// assets/map-1-collision.bin.
//
// Luat: MOI O TILE ma anh bong do PHU LEN (khong chi o neo!) phai DI DUOC, va
// o ngay DUOI no phai DAC. Do chinh la dinh nghia "mat dat". Dat cao mot o ->
// bong lo lung; thap mot o -> bong chim trong dat. Ca hai deu la loi IM LANG
// (client van ve, chi nhin sai).
//
// VI SAO PHAI KIEM CA CHIEU RONG chu khong chi o neo: map 1 co HAI muc mat
// dat (nen thap y=432 va be da cot 16..28 y=408). Mot bong do rong 8 o neo o
// cot 10 keo dai toi cot 17 — hai o cuoi cua no roi vao VACH DA DUNG DUNG cua
// be, tuc la duoc ve o do cao khong co mat dat nao. Ban truoc cua test nay
// chi kiem o neo nen "quan trai" (cot 10) di qua im lang dung loi do.
func TestBongDoMap1NamDungTrenMatDat(t *testing.T) {
	const oTile = 24
	if len(rongMauTheoO) != len(danhSachMauBgItem) {
		t.Fatalf("bang rongMauTheoO co %d dong, bang mau co %d — them mau moi thi phai ghim ca chieu rong o day",
			len(rongMauTheoO), len(danhSachMauBgItem))
	}
	for i, d := range bgItemMap1 {
		if int(d.Mau) < 0 || int(d.Mau) >= len(danhSachMauBgItem) {
			t.Errorf("bong do #%d tro toi mau %d — ngoai bang mau (client se bo qua im lang)", i, d.Mau)
			continue
		}
		rong, co := rongMauTheoO[d.Mau]
		if !co {
			t.Errorf("bong do #%d dung mau %d ma bang rongMauTheoO khong co dong nao — khong kiem duoc chieu rong", i, d.Mau)
			continue
		}
		y := d.TileY * oTile
		for o := int16(0); o < rong; o++ {
			cot := d.TileX + o
			x := cot * oTile
			if !collisionLangLa.DiDuoc(x, y) {
				t.Errorf("bong do #%d (neo o cot %d, rong %d o): o (%d,%d) nam trong vach da — anh se de len da",
					i, d.TileX, rong, cot, d.TileY)
			}
			if collisionLangLa.DiDuoc(x, y+oTile) {
				t.Errorf("bong do #%d (neo o cot %d, rong %d o): o ngay duoi (%d,%d) van di duoc — cho do khong co mat dat, bong se lo lung",
					i, d.TileX, rong, cot, d.TileY+1)
			}
		}
	}
	if len(bgItemMap1) == 0 {
		t.Fatal("map 1 khong con bong do nao — 8 anh bong do lai thanh vo dung")
	}
}

// TestBeDaKhopVoiCollision: hai hang mat dat cua map 1 (hangDatBongDo va
// hangBeDaBongDo) chi dung khi be da THAT SU nam o dung cac cot 16..28. Doc
// lai bien do tu chinh file collision — sua bo cuc map ma quen sua bang bong
// do thi test nay bao truoc, thay vi de bong do lo lung trong game.
func TestBeDaKhopVoiCollision(t *testing.T) {
	const oTile = 24
	for cot := int16(2); cot <= 49; cot++ {
		x := cot * oTile
		treBe := !collisionLangLa.DiDuoc(x, hangDatBongDo*oTile) // hang 18 dac = day la be da
		trongBe := cot >= cotBeDaDau && cot <= cotBeDaCuoi
		if treBe != trongBe {
			t.Errorf("cot %d: collision noi %v con hang so cotBeDa noi %v — be da da doi cho, bang bong do phai sua theo",
				cot, treBe, trongBe)
		}
	}
}

// TestMoiMauDungToiDaHaiLan: qua 2 lan la client bat dau lat guong anh
// (TileMap.isExistMoreOne -> trans = 2). Bong do doi xung nen mat thuong
// khong thay, nhung day la mot hanh vi ngam cua client — giu duoi nguong de
// bang dat cho khong phu thuoc vao no.
func TestMoiMauDungToiDaHaiLan(t *testing.T) {
	dem := map[int16]int{}
	for _, d := range bgItemMap1 {
		dem[d.Mau]++
	}
	for mau, n := range dem {
		if n > 2 {
			t.Errorf("mau %d duoc dat %d lan (> 2) — client se lat guong xen ke cac ban sao", mau, n)
		}
	}
}

// TestGoiMapInfoMap1CoBongDo khoa tung byte cua o bg-item trong MAP INFO cua
// map 1. Gia tri ghim tay (§6.12).
func TestGoiMapInfoMap1CoBongDo(t *testing.T) {
	nv := &store.NhanVat{MapID: 1, ViTriX: 100, ViTriY: 384}
	m, err := goiMapInfo(nv, TimBanDo(1), nil)
	if err != nil {
		t.Fatal(err)
	}
	d := moDocGoi(t, m)

	// Nhay qua phan dau toi truoc o bg-item.
	d.b("mapId", 1)
	d.b("planetId", 0)
	d.b("tileId", 1)
	d.b("bgId", 0)
	d.b("type", 0)
	d.utf("mapName", "Làng Lá")
	d.b("zoneId", 0)
	d.sh("x", 100)
	d.sh("y", 384)
	nhayWaypointMap1(d)
	d.b("mobCount", 0)
	d.b("byte du phong", 0)
	d.b("npcCount", 1)
	d.b("npc.status", 0)
	d.sh("npc.x", 172)
	d.sh("npc.y", 384)
	d.b("npc.template", 0)
	d.sh("npc.portrait", 27078)
	d.b("itemCount", 0)

	// --- o bg-item ---
	d.sh("soBong", 11)
	cho := [][3]int16{
		{0, 2, 18},  // bong nha nho — nha nho trai
		{4, 8, 18},  // bong gieng   — gieng trai
		{2, 8, 18},  // bong quan    — quan trai (neo o cot 8, khong phai 10: xem bgItemMap1)
		{3, 16, 17}, // bong torii   — tren be da
		{6, 19, 17}, // bong thung   — tren be da
		{1, 21, 17}, // bong nha lon — tren be da
		{4, 29, 18}, // bong gieng   — gieng phai
		{7, 31, 18}, // bong den     — cot den phai
		{2, 32, 18}, // bong quan    — quan phai
		{5, 40, 18}, // bong cay     — cay lon
		{0, 44, 18}, // bong nha nho — nha nho phai
	}
	for i, c := range cho {
		d.sh("bong.mau", c[0])
		d.sh("bong.tileX", c[1])
		d.sh("bong.tileY", c[2])
		_ = i
	}

	d.sh("effMap", 0)
	d.b("bgType", 0)
	d.b("idSpaceShip", 0)
	d.b("isMapDouble", 0)
	d.hetGoi("MAP INFO map 1")
}

// nhayWaypointMap1 doc qua dung mot cong cua map 1 (1 byte dem + 8 byte toa
// do + 2 byte co + chuoi ten). Dung o cac test chi quan tam phan phia sau.
func nhayWaypointMap1(d *docGoi) {
	d.t.Helper()
	d.b("waypointCount", 1)
	d.sh("wp.minX", 408)
	d.sh("wp.minY", 384)
	d.sh("wp.maxX", 472)
	d.sh("wp.maxY", 432)
	d.b("wp.laCong", 1)
	d.b("wp.ngoaiTuyen", 0)
	d.utf("wp.ten", "Cổng Làng Lá")
}
