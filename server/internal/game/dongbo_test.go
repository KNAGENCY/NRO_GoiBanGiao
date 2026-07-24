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
	// 6 blob theo thu tu dart/arrow/effect/image/part/skill.
	// Blob rong (dart/arrow/effect/image/skill): do dai 2, noi dung [0x00, 0x00].
	// Part: do dai 87707, noi dung la duLieuPart 87707 byte.
	doKiemBlob := func(ten string, khongRong bool) int32 {
		n, err := r.Int()
		if err != nil {
			t.Fatalf("doc do dai %s loi: %v", ten, err)
		}
		if !khongRong {
			// Blob rong: do dai phai 2, noi dung [0x00, 0x00]
			if n != 2 {
				t.Errorf("%s do dai %d, mong doi 2 (blob rong)", ten, n)
			}
			b1, err1 := r.Byte()
			b2, err2 := r.Byte()
			if err1 != nil || err2 != nil || b1 != 0 || b2 != 0 {
				t.Errorf("%s noi dung khong dung, doc [%d, %d], mong doi [0, 0]", ten, b1, b2)
			}
		} else {
			// Part: chi kiem do dai, bo qua noi dung 87707 byte
			if n != 87707 {
				t.Errorf("%s do dai %d, mong doi 87707", ten, n)
			}
			for i := int32(0); i < n; i++ {
				if _, err := r.Byte(); err != nil {
					t.Fatalf("doc noi dung %s loi: %v", ten, err)
				}
			}
		}
		return n
	}
	doKiemBlob("dart", false)
	doKiemBlob("arrow", false)
	doKiemBlob("effect", false)
	doKiemBlob("image", false)
	doKiemBlob("part", true)
	doKiemBlob("skill", false)
	// Toan goi phai < 100000 byte (moc mark(100000) client).
	if len(m.Data) >= 100000 {
		t.Errorf("goi -87 = %d byte, phai < 100000", len(m.Data))
	}
	// Kiem khong co byte thua sau 6 blob.
	if _, err := r.Byte(); err == nil {
		t.Error("goi -87 con byte thua sau 6 blob — phai het buffer")
	}
}
