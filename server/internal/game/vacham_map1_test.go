package game

import "testing"

// Luoi va cham cua map 1 phai khop bo cuc client — nhung khop theo dung QUY UOC
// cua tung ben, khong phai chep y nguyen:
//
//	client: o co bit 2 la SAN, nhan vat dung TREN no. Dung tren mat co hang 18
//	        thi cy = 18*24 = 432.
//	server: o dac = KHONG VAO DUOC, kiem chinh diem (cx, cy). Voi cy = 432 thi
//	        432/24 = 18 -> hang 18 phai RONG, neu khong moi buoc di deu bi tu choi.
//
// Vi vay: hang MAT DAT de rong, chi phan DUOI mat dat moi dac. Sinh lai bang
// tools/map-lang-la/sinh-va-cham.py doc thang tu res/mymap/1.bytes.
const (
	hangMatDatThuong = 18
	hangMatDatGo     = 17

	// So cot TUONG moi ben. Hai chu khong phai mot: camera dung truoc mep map
	// nen neu chi mot cot thi nhan vat van di duoc toi sat mep va truot han ra
	// ngoai man hinh — nguoi choi mat dau chinh minh.
	soCotTuong = 2
)

// goCuaMap1 tra ve tap cot co mat dat nhich len hang 17 (duoi cong torii va
// duoi nha lon). Giu dung mot cho khai bao de doi bo cuc chi sua o day.
func goCuaMap1() map[int]bool {
	// Nen cao cua khu trung tam: cong torii + nha lon dung tren do (bo cuc day
	// dac theo concept — xem tools/map-lang-la/ve-bo-cuc-map.py).
	g := map[int]bool{}
	for c := 16; c <= 28; c++ {
		g[c] = true
	}
	return g
}

func taiCollisionMap1(t *testing.T) *Collision {
	t.Helper()
	c, err := TaiCollision(duLieuCollision, duLieuAssetLock)
	if err != nil {
		t.Fatalf("TaiCollision: %v", err)
	}
	return c
}

func TestVaChamMap1DungKichThuocLuoi(t *testing.T) {
	c := taiCollisionMap1(t)
	if c.width != 52 || c.height != 24 {
		t.Fatalf("luoi va cham %dx%d, client dung 52x24", c.width, c.height)
	}
	if c.tileWidth != 24 || c.tileHeight != 24 {
		t.Fatalf("o va cham %dx%d, client dung 24x24", c.tileWidth, c.tileHeight)
	}
}

func TestVaChamMap1ChoDungTrenMatDat(t *testing.T) {
	c := taiCollisionMap1(t)
	cacGo := goCuaMap1()
	dac := func(cot, hang int) bool { return c.cells[hang*c.width+cot] == 1 }

	for cot := soCotTuong; cot < c.width-soCotTuong; cot++ {
		hangMat := hangMatDatThuong
		if cacGo[cot] {
			hangMat = hangMatDatGo
		}
		// O mat dat phai RONG: day la o chua chan nhan vat khi dung.
		if dac(cot, hangMat) {
			t.Errorf("cot %d: hang %d (mat dat) dang DAC -> nhan vat dung o do bi coi la vi tri khong hop le",
				cot, hangMat)
		}
		// Ngay duoi mat dat phai DAC: khong duoc dao xuong long dat.
		if !dac(cot, hangMat+1) {
			t.Errorf("cot %d: hang %d (ngay duoi mat dat) dang rong -> chui xuong long dat duoc",
				cot, hangMat+1)
		}
	}
}

func TestVaChamMap1DacHetPhanLongDat(t *testing.T) {
	c := taiCollisionMap1(t)
	cacGo := goCuaMap1()
	for cot := soCotTuong; cot < c.width-soCotTuong; cot++ {
		hangMat := hangMatDatThuong
		if cacGo[cot] {
			hangMat = hangMatDatGo
		}
		for hang := hangMat + 1; hang < c.height; hang++ {
			if c.cells[hang*c.width+cot] != 1 {
				t.Errorf("cot %d hang %d: long dat phai dac nhung dang rong", cot, hang)
			}
		}
	}
}

func TestVaChamMap1ChanNgangHaiCaiGo(t *testing.T) {
	c := taiCollisionMap1(t)
	// Hang 18 duoi go phai DAC: do la than go, di ngang vao la chui vao trong go.
	for _, cot := range []int{16, 20, 24, 28} {
		if c.cells[hangMatDatThuong*c.width+cot] != 1 {
			t.Errorf("cot %d hang %d (than go) dang rong -> di xuyen ngang vao trong go duoc",
				cot, hangMatDatThuong)
		}
	}
	// Cot ngay ben canh go thi hang 18 phai rong (mat dat thuong).
	for _, cot := range []int{14, 15, 29, 30} {
		if c.cells[hangMatDatThuong*c.width+cot] == 1 {
			t.Errorf("cot %d hang %d dang dac nhung day la mat dat thuong", cot, hangMatDatThuong)
		}
	}
}

func TestVaChamMap1CoTuongHaiBien(t *testing.T) {
	c := taiCollisionMap1(t)
	for hang := 0; hang < c.height; hang++ {
		for _, cot := range []int{0, 1, c.width - 2, c.width - 1} {
			if c.cells[hang*c.width+cot] != 1 {
				t.Errorf("tuong bien cot %d hang %d rong -> chay tuot ra ngoai map", cot, hang)
			}
		}
	}
}

func TestVaChamMap1KhongDacTrenTroi(t *testing.T) {
	c := taiCollisionMap1(t)
	// Tren mat dat cao nhat (hang 17) phai la khong khi, tru hai tuong bien.
	for hang := 0; hang < hangMatDatGo; hang++ {
		for cot := soCotTuong; cot < c.width-soCotTuong; cot++ {
			if c.cells[hang*c.width+cot] == 1 {
				t.Errorf("o dac lo lung tren troi tai cot %d hang %d", cot, hang)
			}
		}
	}
}
