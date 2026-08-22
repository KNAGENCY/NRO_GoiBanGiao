package game

import (
	"testing"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// Goi -42 mang chi so nhan vat: mau, chakra, suc danh, thu, chi mang. May chu
// CHUA BAO GIO gui goi nay nen client khong co nguon nao khac de biet mau —
// thanh mau va thanh chakra tren HUD luon hien 0, danh nhau khong doc duoc gi.
//
// Khuon byte doc tu client Controller.cs:1704. Luu y: readIntToLongDQT() doc
// INT 4 byte vi ModFunc.isReadInt = true va isReadDouble = false.

func nhanVatMau() *store.NhanVat {
	return &store.NhanVat{
		ID: 7, Ten: "Naruto", Lang: 0, KieuToc: 30,
		HpToiDa: 1500, HpHienTai: 1200,
		ChakraToiDa: 800, ChakraHienTai: 640,
		SucManh: 50_000, Cap: 3,
	}
}

func TestGoiChiSoNhanVatDungLenh(t *testing.T) {
	m := goiChiSoNhanVat(nhanVatMau())
	if m.Command != lenhChiSoNhanVat {
		t.Fatalf("lenh = %d, mong doi %d", m.Command, lenhChiSoNhanVat)
	}
	if lenhChiSoNhanVat != -42 {
		t.Errorf("lenhChiSoNhanVat = %d, client doc o case -42", lenhChiSoNhanVat)
	}
}

func TestGoiChiSoNhanVatDungKhuonByte(t *testing.T) {
	nv := nhanVatMau()
	r := protocol.NewReader(goiChiSoNhanVat(nv).Data)

	docInt := func(ten string) int32 {
		v, err := r.Int()
		if err != nil {
			t.Fatalf("doc %s loi: %v", ten, err)
		}
		return v
	}
	docByte := func(ten string) byte {
		v, err := r.Byte()
		if err != nil {
			t.Fatalf("doc %s loi: %v", ten, err)
		}
		return v
	}

	// Bay chi so dau: goc / day / hien tai
	if v := docInt("cHPGoc"); v != nv.HpToiDa {
		t.Errorf("cHPGoc = %d, mong doi %d", v, nv.HpToiDa)
	}
	if v := docInt("cMPGoc"); v != nv.ChakraToiDa {
		t.Errorf("cMPGoc = %d, mong doi %d", v, nv.ChakraToiDa)
	}
	docInt("cDamGoc")
	if v := docInt("cHPFull"); v != nv.HpToiDa {
		t.Errorf("cHPFull = %d, mong doi %d", v, nv.HpToiDa)
	}
	if v := docInt("cMPFull"); v != nv.ChakraToiDa {
		t.Errorf("cMPFull = %d, mong doi %d", v, nv.ChakraToiDa)
	}
	if v := docInt("cHP"); v != nv.HpHienTai {
		t.Errorf("cHP = %d, mong doi %d", v, nv.HpHienTai)
	}
	if v := docInt("cMP"); v != nv.ChakraHienTai {
		t.Errorf("cMP = %d, mong doi %d", v, nv.ChakraHienTai)
	}

	if v := docByte("cspeed"); v == 0 {
		t.Error("cspeed = 0 -> nhan vat khong di chuyen duoc theo cach client tinh")
	}
	docByte("hpFrom1000TiemNang")
	docByte("mpFrom1000TiemNang")
	docByte("damFrom1000TiemNang")

	docInt("cDamFull")
	docInt("cDefull")
	docByte("cCriticalFull")

	if _, err := r.Long(); err != nil { // cTiemNang
		t.Fatalf("doc cTiemNang loi: %v", err)
	}
	if _, err := r.Short(); err != nil { // expForOneAdd
		t.Fatalf("doc expForOneAdd loi: %v", err)
	}
	if _, err := r.Short(); err != nil { // cDefGoc
		t.Fatalf("doc cDefGoc loi: %v", err)
	}
	docByte("cCriticalGoc")

	// Het goi. Phan sau (tlDef, tlPst, …) client boc trong try/catch va tu dat 0
	// neu thieu, nen ban ngan dung o day la hop le — nhung khong duoc THUA byte,
	// thua la client doc lech sang truong khac.
	if _, err := r.Byte(); err == nil {
		t.Error("goi -42 con byte thua sau cCriticalGoc — client se doc lech")
	}
}

func TestGoiChiSoNhanVatKhongGuiMauQuaDay(t *testing.T) {
	// Mau hien tai > mau day thi thanh mau client ve tran ra ngoai khung.
	nv := nhanVatMau()
	nv.HpHienTai = nv.HpToiDa + 500
	nv.ChakraHienTai = nv.ChakraToiDa + 500
	r := protocol.NewReader(goiChiSoNhanVat(nv).Data)
	for i := 0; i < 5; i++ {
		if _, err := r.Int(); err != nil { // bo qua 5 truong dau
			t.Fatalf("doc truong %d loi: %v", i, err)
		}
	}
	hp, _ := r.Int()
	mp, _ := r.Int()
	if hp > nv.HpToiDa {
		t.Errorf("cHP = %d > cHPFull = %d -> thanh mau tran khung", hp, nv.HpToiDa)
	}
	if mp > nv.ChakraToiDa {
		t.Errorf("cMP = %d > cMPFull = %d -> thanh chakra tran khung", mp, nv.ChakraToiDa)
	}
}

func TestGoiChiSoNhanVatChiuDuocNhanVatMoiToanSo0(t *testing.T) {
	// Nhan vat vua tao co the chua co chi so. Khong duoc chia cho 0 hay panic.
	m := goiChiSoNhanVat(&store.NhanVat{ID: 1, Ten: "x"})
	if len(m.Data) == 0 {
		t.Fatal("goi rong")
	}
	r := protocol.NewReader(m.Data)
	hpGoc, err := r.Int()
	if err != nil {
		t.Fatalf("doc cHPGoc loi: %v", err)
	}
	if hpGoc <= 0 {
		t.Errorf("cHPGoc = %d — nhan vat khong co mau thi client coi nhu da chet", hpGoc)
	}
}
