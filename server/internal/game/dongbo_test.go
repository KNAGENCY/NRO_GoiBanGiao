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

func TestGoiCatalogSkillCoDu5Lang(t *testing.T) {
	m, err := goiCatalogSkill()
	if err != nil {
		t.Fatalf("goiCatalogSkill loi: %v", err)
	}
	if m.Command != -28 {
		t.Fatalf("lenh = %d, mong doi -28", m.Command)
	}
	r := protocol.NewReader(m.Data)
	sub, _ := r.Byte()
	if sub != lenhConCatalogSkill { // 7
		t.Errorf("lenh con = %d, mong doi 7", sub)
	}
	if v, _ := r.Byte(); v != vsSkill {
		t.Errorf("vsSkill = %d, mong doi %d", v, vsSkill)
	}
	if opt, _ := r.Byte(); opt != 0 {
		t.Errorf("countOption = %d, mong doi 0", opt)
	}
	countClass, _ := r.Byte()
	// Ghim 5 ten lang doc lap voi DanhSachLang() de bat duoc neu ai do go sai/dao
	// thu tu/dich nham ten lang trong lang.go — day la noi dung loi cua ke hoach
	// (5 lang Naruto thay 3 lop Dragon Ball).
	tenLangMongDoi := []string{"Làng Lá", "Làng Cát", "Làng Sương Mù", "Làng Đá", "Làng Mây"}
	if int(countClass) != len(tenLangMongDoi) {
		t.Fatalf("countClass = %d, mong doi %d (so lang Naruto)", countClass, len(tenLangMongDoi))
	}
	if countClass < 5 {
		t.Errorf("countClass = %d, phai >= 5 de id lang 0-4 hop le trong ME_INFO", countClass)
	}
	// Doc tung lang: ten (UTF) + countTemplate (4, tu Task 1 KH6) + noi dung tung
	// chieu. Ten phai khop danh muc lang (ghim doc lap). Kiem noi dung chi tiet
	// (4 chieu/lang, khong templateId==0, dung 1 chieu type=2) o
	// chieuthuc_test.go — o day chi doc du de khong con byte thua.
	for i := 0; i < int(countClass); i++ {
		ten, err := r.UTF()
		if err != nil {
			t.Fatalf("doc ten lang %d loi: %v", i, err)
		}
		if ten != tenLangMongDoi[i] {
			t.Errorf("ten lang %d = %q, mong doi %q", i, ten, tenLangMongDoi[i])
		}
		countTemplate, err := r.Byte()
		if err != nil {
			t.Fatalf("doc countTemplate lang %d loi: %v", i, err)
		}
		if countTemplate != 4 {
			t.Errorf("countTemplate lang %d = %d, mong doi 4 (4 chieu x 5 lang, dac ta §4.3)", i, countTemplate)
		}
		for j := 0; j < int(countTemplate); j++ {
			if err := boQuaMotMauChieu(r); err != nil {
				t.Fatalf("doc chieu %d cua lang %d loi: %v", j, i, err)
			}
		}
	}
	// Het du lieu — khong con byte thua sau khi doc du 5 lang.
	if _, err := r.Byte(); err == nil {
		t.Error("goi -28/7 con byte thua sau khi doc du 5 lang x 4 chieu — sai khuon byte")
	}
}

// boQuaMotMauChieu doc va bo qua mot MauChieu theo dung khuon dac ta §3.7, chi
// de xac nhan khuon byte tu-nhat-quan (doc het khong loi, khong lech). Kiem
// gia tri that duoc lam o TestCatalogChieu* (chieuthuc_test.go).
func boQuaMotMauChieu(r *protocol.Reader) error {
	if _, err := r.Byte(); err != nil { // templateId
		return err
	}
	if _, err := r.UTF(); err != nil { // ten
		return err
	}
	if _, err := r.Byte(); err != nil { // maxPoint
		return err
	}
	if _, err := r.Byte(); err != nil { // manaUseType
		return err
	}
	if _, err := r.Byte(); err != nil { // type
		return err
	}
	if _, err := r.Short(); err != nil { // iconId
		return err
	}
	if _, err := r.UTF(); err != nil { // damInfo
		return err
	}
	if _, err := r.UTF(); err != nil { // ghiChu
		return err
	}
	countSkill, err := r.Byte()
	if err != nil {
		return err
	}
	for i := 0; i < int(countSkill); i++ {
		if _, err := r.Short(); err != nil { // skillId
			return err
		}
		if _, err := r.Byte(); err != nil { // point
			return err
		}
		if _, err := r.Long(); err != nil { // powRequire (8 byte)
			return err
		}
		if _, err := r.Short(); err != nil { // manaUse
			return err
		}
		if _, err := r.Int(); err != nil { // coolDown
			return err
		}
		if _, err := r.Short(); err != nil { // dx
			return err
		}
		if _, err := r.Short(); err != nil { // dy
			return err
		}
		if _, err := r.Byte(); err != nil { // maxFight
			return err
		}
		if _, err := r.Short(); err != nil { // damage
			return err
		}
		if _, err := r.Short(); err != nil { // price
			return err
		}
		if _, err := r.UTF(); err != nil { // moreInfo
			return err
		}
	}
	return nil
}

func TestGoiCatalogMapRong(t *testing.T) {
	m := goiCatalogMap()
	if m.Command != -28 {
		t.Fatalf("lenh = %d, mong doi -28", m.Command)
	}
	// [6][vsMap][countMap=0][countNpc=0][countMob=0]
	want := []byte{6, 20, 0, 0, 0}
	if !bytes.Equal(m.Data, want) {
		t.Errorf("than catalog map = %v, mong doi %v", m.Data, want)
	}
}

func TestGoiCatalogItemSubtype2Rong(t *testing.T) {
	m := goiCatalogItem()
	if m.Command != -28 {
		t.Fatalf("lenh = %d, mong doi -28", m.Command)
	}
	// [8][vsItem][subtype=2][start short=0][end short=0]
	want := []byte{8, 20, 2, 0, 0, 0, 0}
	if !bytes.Equal(m.Data, want) {
		t.Errorf("than catalog item = %v, mong doi %v", m.Data, want)
	}
}
