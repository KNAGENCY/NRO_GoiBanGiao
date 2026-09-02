package game

import (
	"encoding/binary"
	"testing"

	"nroserver/internal/store"
)

// ---------------------------------------------------------------------------
// Test cho bang phien dich (noi chua, o so) -> vat_pham.id.
//
// Test quan trong nhat la TestDungBangVaGoiKhongLech o cuoi: no doc NGUOC goi
// -36 dung trinh tu cua client, lay mau_id o dung vi tri byte cua o, roi doi
// chieu voi id ma bang tra ra o CHINH o do. Do la bang chung "dung cung luc".
// ---------------------------------------------------------------------------

func TestDungBangTraDungId(t *testing.T) {
	ds := []store.VatPham{
		mon(11, 0, 0, store.NoiChuaTui, 1),
		mon(22, 1, 3, store.NoiChuaTui, 1),
		mon(33, 2, 1, store.NoiChuaMac, 1),
	}
	b, err := DungBangHanhTrang(ds)
	if err != nil {
		t.Fatalf("DungBangHanhTrang loi: %v", err)
	}
	if got := b.SoMon(); got != 3 {
		t.Fatalf("SoMon = %d, mong doi 3", got)
	}
	for _, c := range []struct {
		noiChua string
		o       int16
		id      int64
	}{
		{store.NoiChuaTui, 0, 11},
		{store.NoiChuaTui, 3, 22},
		{store.NoiChuaMac, 1, 33},
	} {
		id, co := b.TraId(c.noiChua, c.o)
		if !co || id != c.id {
			t.Errorf("TraId(%s,%d) = (%d,%v), mong doi (%d,true)", c.noiChua, c.o, id, co, c.id)
		}
	}
}

func TestTraIdORongTraFalse(t *testing.T) {
	b, err := DungBangHanhTrang([]store.VatPham{mon(11, 0, 0, store.NoiChuaTui, 1)})
	if err != nil {
		t.Fatalf("loi: %v", err)
	}
	if id, co := b.TraId(store.NoiChuaTui, 4); co {
		t.Errorf("o 4 rong nhung TraId tra (%d,true) — nguoi goi se thao tac nham mon", id)
	}
}

// TestTuiVaRuongLaHaiKhoaKhacNhau: dot bien "bang bo qua NoiChua, chi khoa theo
// o" phai lam test nay do.
func TestTuiVaRuongLaHaiKhoaKhacNhau(t *testing.T) {
	ds := []store.VatPham{
		mon(11, 0, 2, store.NoiChuaTui, 1),
		mon(99, 1, 2, store.NoiChuaRuong, 1),
	}
	b, err := DungBangHanhTrang(ds)
	if err != nil {
		t.Fatalf("o 2 cua TUI va o 2 cua RUONG la hai o khac nhau, khong duoc bao trung: %v", err)
	}
	if id, _ := b.TraId(store.NoiChuaTui, 2); id != 11 {
		t.Errorf("TraId(TUI,2) = %d, mong doi 11", id)
	}
	if id, _ := b.TraId(store.NoiChuaRuong, 2); id != 99 {
		t.Errorf("TraId(RUONG,2) = %d, mong doi 99", id)
	}
}

// TestHaiMonCungOLaLoi: dot bien "bo kiem trung o" phai lam test nay do.
func TestHaiMonCungOLaLoi(t *testing.T) {
	ds := []store.VatPham{
		mon(11, 0, 2, store.NoiChuaTui, 1),
		mon(12, 1, 2, store.NoiChuaTui, 1),
	}
	if _, err := DungBangHanhTrang(ds); err == nil {
		t.Fatal("hai mon cung o 2 cua TUI ma khong bao loi — bang se im lang mat mot mon")
	}
}

func TestOVuotTranLaLoi(t *testing.T) {
	ds := []store.VatPham{mon(11, 0, tranSoO, store.NoiChuaTui, 1)}
	if _, err := DungBangHanhTrang(ds); err == nil {
		t.Fatalf("o %d (= tranSoO) phai bi tu choi vi bien hop le la [0,%d)", tranSoO, tranSoO)
	}
	ds = []store.VatPham{mon(11, 0, -1, store.NoiChuaTui, 1)}
	if _, err := DungBangHanhTrang(ds); err == nil {
		t.Fatal("o -1 phai bi tu choi")
	}
}

func TestMonKhongCoOSoLaLoi(t *testing.T) {
	ds := []store.VatPham{{ID: 11, MauID: 0, NoiChua: store.NoiChuaTui, OSo: nil, SoLuong: 1}}
	if _, err := DungBangHanhTrang(ds); err == nil {
		t.Fatal("mon o TUI ma o_so = nil phai la loi")
	}
}

func TestMonDuoiDatBiBoQuaChuKhongLoi(t *testing.T) {
	ds := []store.VatPham{
		{ID: 50, MauID: 0, NoiChua: store.NoiChuaDuoiDat, OSo: nil, SoLuong: 1},
		mon(11, 0, 0, store.NoiChuaTui, 1),
	}
	b, err := DungBangHanhTrang(ds)
	if err != nil {
		t.Fatalf("mon duoi dat phai duoc bo qua lang le: %v", err)
	}
	if b.SoMon() != 1 {
		t.Errorf("SoMon = %d, mong doi 1 (mon duoi dat khong nam trong hanh trang)", b.SoMon())
	}
}

// ---------------------------------------------------------------------------
// Test "dung cung luc": doc nguoc goi -36 va doi chieu voi bang.
// ---------------------------------------------------------------------------

// docGoiTui doc nguoc than goi -36 theo DUNG trinh tu client, tra ve mang soO
// phan tu: -1 = o rong, con lai = mau_id doc duoc tai dung vi tri byte cua o do.
func docGoiTui(t *testing.T, than []byte) []int32 {
	t.Helper()
	i := 0
	u8 := func() byte {
		if i >= len(than) {
			t.Fatalf("goi het byte som tai %d", i)
		}
		v := than[i]
		i++
		return v
	}
	i16 := func() int16 {
		if i+2 > len(than) {
			t.Fatalf("goi het byte som tai %d", i)
		}
		v := int16(binary.BigEndian.Uint16(than[i:]))
		i += 2
		return v
	}
	i32 := func() int32 {
		if i+4 > len(than) {
			t.Fatalf("goi het byte som tai %d", i)
		}
		v := int32(binary.BigEndian.Uint32(than[i:]))
		i += 4
		return v
	}
	utf := func() {
		n := int(uint16(i16()))
		i += n
	}

	if c := u8(); c != 0 {
		t.Fatalf("cAction = %d, mong doi 0", c)
	}
	soO := int(u8())
	ra := make([]int32, soO)
	for o := 0; o < soO; o++ {
		mau := i16()
		if mau == -1 {
			ra[o] = -1
			continue
		}
		ra[o] = int32(mau)
		i32() // so luong
		utf() // info
		utf() // content
		n := int(u8())
		for k := 0; k < n; k++ {
			u8()  // option id
			i16() // param
		}
	}
	if i != len(than) {
		t.Fatalf("doc xong con thua %d byte — goi va bo doc lech nhau", len(than)-i)
	}
	return ra
}

// TestDungBangVaGoiKhongLech: voi MOI o co mon trong goi -36, id ma bang tra ra
// phai la id cua mon co dung mau_id doc duoc tu goi tai o do. Dot bien "ham
// tien ich dung bang tu lat cat da loc khac voi lat cat ghi goi" lam test do.
func TestDungBangVaGoiKhongLech(t *testing.T) {
	ds := []store.VatPham{
		mon(11, 0, 0, store.NoiChuaTui, 5),
		mon(22, 1, 2, store.NoiChuaTui, 1),
		mon(33, 4, 3, store.NoiChuaTui, 1, store.TuyChon{Id: 0, Param: 7}),
		mon(44, 2, 0, store.NoiChuaMac, 1),
		mon(55, 3, 1, store.NoiChuaRuong, 1),
	}
	// mauTheoId: tra cuu nguoc de doi chieu.
	mauTheoId := map[int64]int32{11: 0, 22: 1, 33: 4, 44: 2, 55: 3}

	ht, err := DungBangVaGoi(ds, 100, 4, SoOMacChuan, 6)
	if err != nil {
		t.Fatalf("DungBangVaGoi loi: %v", err)
	}
	if ht.GoiTui.Command != -36 || ht.GoiMac.Command != -37 || ht.GoiRuong.Command != -35 {
		t.Fatalf("ma lenh sai: %d/%d/%d", ht.GoiTui.Command, ht.GoiMac.Command, ht.GoiRuong.Command)
	}

	oTrongGoi := docGoiTui(t, ht.GoiTui.Data)
	if len(oTrongGoi) != 4 {
		t.Fatalf("goi -36 co %d o, mong doi 4", len(oTrongGoi))
	}
	for o, mau := range oTrongGoi {
		id, co := ht.Bang.TraId(store.NoiChuaTui, int16(o))
		if mau == -1 {
			if co {
				t.Errorf("o %d rong trong goi nhung bang tra id=%d — client va may chu lech", o, id)
			}
			continue
		}
		if !co {
			t.Errorf("o %d co mau %d trong goi nhung bang khong tra id nao", o, mau)
			continue
		}
		if mauTheoId[id] != mau {
			t.Errorf("o %d: goi noi mau=%d, bang tra id=%d (mau=%d) — DUNG LECH",
				o, mau, id, mauTheoId[id])
		}
	}
}

func TestDungBangVaGoiLoiThiKhongTraGiCa(t *testing.T) {
	ds := []store.VatPham{
		mon(11, 0, 2, store.NoiChuaTui, 1),
		mon(12, 1, 2, store.NoiChuaTui, 1), // trung o
	}
	ht, err := DungBangVaGoi(ds, 0, 5, SoOMacChuan, 5)
	if err == nil {
		t.Fatal("trung o phai tra loi")
	}
	if ht != nil {
		t.Error("co loi ma van tra ket qua — nguoi goi co the gui goi nua vo")
	}
}
