package game

import (
	"bytes"
	"testing"

	"nroserver/internal/protocol"
)

func TestNhungPartDungKichThuoc(t *testing.T) {
	// File part nhung phai la file goc 87707 byte. Sai kich thuoc = nhung nham
	// hoac file hong -> nhan vat khong ve duoc o Ke hoach 4.
	if len(duLieuPart) != 87707 {
		t.Errorf("duLieuPart = %d byte, mong doi 87707 (file part goc)", len(duLieuPart))
	}
}

func TestPhienBanNarutoKhacBanCu(t *testing.T) {
	// Phai khac phien ban Dragon Ball cu (10/2/1/4) de ghi de cache client.
	// Dac biet vsSkill phai != 1, khong client dung cache 3 lop cu.
	if vsData == 10 || vsMap == 2 || vsSkill == 1 || vsItem == 4 {
		t.Error("phien ban Naruto trung ban Dragon Ball cu -> client se dung cache cu")
	}
	if vsData == 0 || vsMap == 0 || vsSkill == 0 || vsItem == 0 {
		t.Error("phien ban = 0 trung mac dinh client moi, khong ep cap nhat duoc dung cach")
	}
}

func TestGoiVersionGameDungByte(t *testing.T) {
	m := goiVersionGame()
	if m.Command != -28 {
		t.Fatalf("lenh = %d, mong doi -28", m.Command)
	}
	// than: [4][vsData][vsMap][vsSkill][vsItem][0][countCaption=0]
	want := []byte{4, 20, 20, 20, 20, 0, 0}
	if !bytes.Equal(m.Data, want) {
		t.Errorf("than versionGame = %v, mong doi %v", m.Data, want)
	}
}

func TestGoiCatalogDataDungCauTruc(t *testing.T) {
	m := goiCatalogData()
	if m.Command != -87 {
		t.Fatalf("lenh = %d, mong doi -87", m.Command)
	}
	r := protocol.NewReader(m.Data)
	if b, _ := r.Byte(); b != vsData {
		t.Errorf("byte dau = %d, mong doi vsData %d", b, vsData)
	}
	// 6 blob theo thu tu dart/arrow/effect/image/part/skill; chi part co noi dung.
	doDaiBlob := func(ten string) int32 {
		n, err := r.Int()
		if err != nil {
			t.Fatalf("doc do dai %s loi: %v", ten, err)
		}
		// bo qua noi dung blob
		for i := int32(0); i < n; i++ {
			if _, err := r.Byte(); err != nil {
				t.Fatalf("doc noi dung %s loi: %v", ten, err)
			}
		}
		return n
	}
	if n := doDaiBlob("dart"); n != 2 {
		t.Errorf("dart do dai %d, mong doi 2 (rong)", n)
	}
	doDaiBlob("arrow")
	doDaiBlob("effect")
	doDaiBlob("image")
	if n := doDaiBlob("part"); n != 87707 {
		t.Errorf("part do dai %d, mong doi 87707", n)
	}
	doDaiBlob("skill")
	// Toan goi phai < 100000 byte (moc mark(100000) client).
	if len(m.Data) >= 100000 {
		t.Errorf("goi -87 = %d byte, phai < 100000", len(m.Data))
	}
}
