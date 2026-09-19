package game

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"

	"nroserver/internal/protocol"
)

type khoaArtifactTest struct {
	Versions map[string]byte `json:"versions"`
	Outputs  []struct {
		Path   string `json:"path"`
		SHA256 string `json:"sha256"`
		Bytes  int    `json:"bytes"`
	} `json:"outputs"`
}

func TestRuntimeServerNhungDungByteVaSHATuLock(t *testing.T) {
	var lock khoaArtifactTest
	if err := json.Unmarshal(duLieuAssetLock, &lock); err != nil {
		t.Fatalf("asset-lock.json nhúng không đọc được: %v", err)
	}
	runtime := map[string][]byte{
		"server/part.bin":            duLieuPart,
		"server/skill.bin":           duLieuSkill,
		"server/tile_set_info.bin":   duLieuTileSet,
		"server/map-1-collision.bin": duLieuCollision,
		"server/catalog.json":        duLieuCatalog,
	}
	for path, data := range runtime {
		var expectedSHA string
		var expectedBytes int
		for _, output := range lock.Outputs {
			if output.Path == path {
				expectedSHA = output.SHA256
				expectedBytes = output.Bytes
				break
			}
		}
		if expectedSHA == "" {
			t.Fatalf("asset lock thiếu %s", path)
		}
		sum := sha256.Sum256(data)
		if got := hex.EncodeToString(sum[:]); got != expectedSHA {
			t.Errorf("%s SHA = %s, lock = %s", path, got, expectedSHA)
		}
		if len(data) != expectedBytes {
			t.Errorf("%s = %d byte, lock = %d", path, len(data), expectedBytes)
		}
	}
	if vsData != lock.Versions["data"] ||
		vsMap != lock.Versions["map"] ||
		vsSkill != lock.Versions["resource"] {
		t.Fatalf(
			"version server data/map/skill = %d/%d/%d, lock = %d/%d/%d",
			vsData, vsMap, vsSkill,
			lock.Versions["data"], lock.Versions["map"], lock.Versions["resource"],
		)
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
	//
	// So viet cung co chu dich: doi mot byte o day nghia la client se coi cache
	// cua no da cu va tai lai du lieu tinh — phai la quyet dinh co y, khong duoc
	// troi theo mot thay doi khac. vsData 21 -> 22 khi part.bin doi (nang 3 manh
	// nhan vat cho chan cham dat); khong tang thi client van dung part.bin cu
	// trong RMS va nhan vat tiep tuc lun nua nguoi.
	//
	// 26 -> 27 (05/09/2026): sprite nguoi choi sinh lai bang MOT ti le chung cho
	// moi khung + neo ban chan, nen dx/dy cua part 30/57/58 doi het.
	//
	// 27 -> 28 (05/09/2026, chieu — V1 khe dau-than): mang than cat tu cam xuong
	// (het hai khuon mat), ti le TUNG khung theo co dau, slot than 2..6 = 5 khung
	// chay, slot 15 het khung nam, dx/dy giai bang MILP -> part 30/57/58 doi het lan nua.
	want := []byte{4, 29, 22, 21, 20, 0, 0}
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
	// Part và skill là đúng nguyên blob đã khóa; bốn blob còn lại rỗng.
	doKiemBlob := func(ten string, want []byte) int32 {
		n, err := r.Int()
		if err != nil {
			t.Fatalf("doc do dai %s loi: %v", ten, err)
		}
		if want == nil {
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
			if n != int32(len(want)) {
				t.Errorf("%s do dai %d, mong doi %d", ten, n, len(want))
			}
			got := make([]byte, n)
			for i := int32(0); i < n; i++ {
				value, err := r.Byte()
				if err != nil {
					t.Fatalf("doc noi dung %s loi: %v", ten, err)
				}
				got[i] = value
			}
			if !bytes.Equal(got, want) {
				t.Errorf("%s không byte-identical artifact đã nhúng", ten)
			}
		}
		return n
	}
	doKiemBlob("dart", nil)
	doKiemBlob("arrow", nil)
	doKiemBlob("effect", duLieuHieuUngTrong)
	doKiemBlob("image", nil)
	doKiemBlob("part", duLieuPart)
	doKiemBlob("skill", duLieuSkill)
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

func TestGoiCatalogMapLangLa(t *testing.T) {
	m := goiCatalogMap()
	if m.Command != -28 {
		t.Fatalf("lenh = %d, mong doi -28", m.Command)
	}
	r := protocol.NewReader(m.Data)
	if sub, _ := r.Byte(); sub != 6 {
		t.Fatalf("sub = %d", sub)
	}
	if v, _ := r.Byte(); v != vsMap {
		t.Fatalf("vsMap = %d", v)
	}
	n, _ := r.Byte()
	if n != 6 {
		t.Fatalf("countMap = %d", n)
	}
	wantMaps := []string{"Làng Lá", "Làng Cát", "Làng Sương Mù", "Làng Đá", "Làng Mây", "Khu rừng"}
	for i, want := range wantMaps {
		got, err := r.UTF()
		if err != nil || got != want {
			t.Fatalf("map %d = %q (%v), muon %q", i, got, err, want)
		}
	}
	n, _ = r.Byte()
	if n != 1 {
		t.Fatalf("countNpc = %d", n)
	}
	if name, _ := r.UTF(); name != "Iruka" {
		t.Fatalf("NPC = %q", name)
	}
	for i := 0; i < 3; i++ {
		if _, err := r.Short(); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := r.Byte(); err != nil {
		t.Fatal(err)
	}
	n, _ = r.Byte()
	if n != 1 {
		t.Fatalf("countMob = %d", n)
	}
	if typ, _ := r.Byte(); typ != 0 {
		t.Fatalf("mob type = %d", typ)
	}
	if name, _ := r.UTF(); name != "Sói Rừng" {
		t.Fatalf("mob name = %q", name)
	}
	if hp, _ := r.Int(); hp != 100 {
		t.Fatalf("mob hp = %d", hp)
	}
	for i := 0; i < 3; i++ {
		if _, err := r.Byte(); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := r.Byte(); err == nil {
		t.Fatal("-28/6 co byte thua")
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

// TestBangHieuUngDuChoLiveFromDead: client tra efs[id-1] khong kiem bien —
// liveFromDead goi hieu ung 109 truoc khi go khoa di chuyen. Bang phai >= 109
// o, moi o dung 1 khung idImg = -1 (khong ve, khong xin anh). Gap 19/09/2026.
func TestBangHieuUngDuChoLiveFromDead(t *testing.T) {
	r := protocol.NewReader(duLieuHieuUngTrong)
	n, _ := r.Short()
	if n < 109 {
		t.Fatalf("so o hieu ung = %d, phai >= 109 (Char.liveFromDead goi 109)", n)
	}
	for i := int16(0); i < n; i++ {
		id, _ := r.Short()
		k, _ := r.Byte()
		img, _ := r.Short()
		dx, _ := r.Byte()
		dy, err := r.Byte()
		if err != nil || id != i+1 || k != 1 || img != -1 || dx != 0 || dy != 0 {
			t.Fatalf("o %d = id %d khung %d img %d dx %d dy %d (loi %v), muon id %d khung 1 img -1", i, id, k, img, dx, dy, err, i+1)
		}
	}
	if _, err := r.Byte(); err == nil {
		t.Fatal("con byte thua sau bang hieu ung")
	}
}
