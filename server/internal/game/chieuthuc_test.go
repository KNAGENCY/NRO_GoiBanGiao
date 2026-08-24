package game

import (
	"bytes"
	"testing"

	"nroserver/internal/protocol"
)

// TestCatalogChieuMotChieuMotCap la RED TEST cua Task 1 (buoc 1-2): mot chieu
// TU DUNG lam fixture (khong lay tu chieuTheoLang that), 1 cap duy nhat, so
// TOAN BO mang byte voi gia tri ghim cung theo khuon dac ta §3.7:
//
//	byte   templateId
//	UTF    ten
//	byte   maxPoint
//	byte   manaUseType
//	byte   type
//	short  iconId
//	UTF    damInfo
//	UTF    ghiChu               (luon rong)
//	byte   countSkill
//	  short skillId
//	  byte  point
//	  long  powRequire          (8 byte)
//	  short manaUse
//	  int   coolDown
//	  short dx
//	  short dy
//	  byte  maxFight
//	  short damage
//	  short price
//	  UTF   moreInfo
func TestCatalogChieuMotChieuMotCap(t *testing.T) {
	fixture := MauChieu{
		Id:            1, // templateId tu 1, KHONG dung 0 (xem chieuthuc.go)
		Ten:           "Chieu Thu",
		CapToiDa:      1,
		KieuTonChakra: 0,
		Loai:          1,
		IconId:        7,
		MoTa:          "Mo ta",
		CacCap: []CapChieu{
			{
				SkillId:      42,
				Diem:         1,
				NhanLucCan:   1000,
				ChakraDung:   30,
				HoiChieuMs:   1500,
				Dx:           100,
				Dy:           100,
				MaxFight:     1,
				SatThuong:    50,
				Gia:          200,
				ThemThongTin: "Thong tin",
			},
		},
	}

	w := protocol.NewWriter()
	if err := ghiMauChieu(w, fixture); err != nil {
		t.Fatalf("ghiMauChieu loi: %v", err)
	}

	want := []byte{}
	want = append(want, 1)                        // templateId = 1
	want = append(want, 0, 9)                     // do dai UTF "Chieu Thu" = 9 byte
	want = append(want, []byte("Chieu Thu")...)   //
	want = append(want, 1)                        // maxPoint = 1
	want = append(want, 0)                        // manaUseType = 0
	want = append(want, 1)                        // type = 1
	want = append(want, 0, 7)                     // iconId = 7 (short)
	want = append(want, 0, 5)                     // do dai UTF "Mo ta" = 5 byte
	want = append(want, []byte("Mo ta")...)       //
	want = append(want, 0, 0)                     // ghiChu = "" (do dai 0)
	want = append(want, 1)                        // countSkill = 1
	want = append(want, 0, 42)                    // skillId = 42 (short)
	want = append(want, 1)                        // point = 1
	want = append(want, 0, 0, 0, 0, 0, 0, 3, 232) // powRequire = 1000, LONG 8 byte
	want = append(want, 0, 30)                    // manaUse = 30 (short)
	want = append(want, 0, 0, 5, 220)             // coolDown = 1500, int 4 byte
	want = append(want, 0, 100)                   // dx = 100 (short)
	want = append(want, 0, 100)                   // dy = 100 (short)
	want = append(want, 1)                        // maxFight = 1
	want = append(want, 0, 50)                    // damage = 50 (short)
	want = append(want, 0, 200)                   // price = 200 (short)
	want = append(want, 0, 9)                     // do dai UTF "Thong tin" = 9 byte
	want = append(want, []byte("Thong tin")...)   //

	if !bytes.Equal(w.Bytes(), want) {
		t.Errorf("ghiMauChieu(mot chieu mot cap) =\n%v\nmong doi\n%v", w.Bytes(), want)
	}
}

// TestCatalogChieuNhanLucCanChiem8Byte khang dinh RIENG rieng rang powRequire chiem DUNG 8
// byte trong dong — day la bay §3.7 diem 1 (nham thanh int 4 byte se lech het
// cac truong phia sau). Dem do dai buffer TRUOC/SAU khi ghi 1 CapChieu de bat
// duoc dot bien doi Long -> Int ma khong phu thuoc vao cac truong khac.
func TestCatalogChieuNhanLucCanChiem8Byte(t *testing.T) {
	cap := CapChieu{
		SkillId:      1,
		Diem:         1,
		NhanLucCan:   123456789012345, // > int32 max, chi bieu dien dung duoc trong 8 byte
		ChakraDung:   10,
		HoiChieuMs:   1000,
		Dx:           100,
		Dy:           100,
		MaxFight:     1,
		SatThuong:    10,
		Gia:          10,
		ThemThongTin: "",
	}
	w := protocol.NewWriter()
	w.Short(cap.SkillId)
	w.Byte(cap.Diem)
	truocPowRequire := len(w.Bytes())
	w.Long(cap.NhanLucCan)
	sauPowRequire := len(w.Bytes())
	if sauPowRequire-truocPowRequire != 8 {
		t.Fatalf("powRequire chiem %d byte, mong doi DUNG 8 (long)", sauPowRequire-truocPowRequire)
	}

	// Doc lai bang Reader.Long() phai ra dung gia tri (kiem ca chieu ghi/doc).
	r := protocol.NewReader(w.Bytes()[truocPowRequire:sauPowRequire])
	got, err := r.Long()
	if err != nil {
		t.Fatalf("doc lai powRequire loi: %v", err)
	}
	if got != cap.NhanLucCan {
		t.Errorf("doc lai powRequire = %d, mong doi %d", got, cap.NhanLucCan)
	}
}

// --- Tu day tro xuong: kiem du lieu THAT cua ChieuCuaLang (buoc 5 cua Task 1) ---

// TestCatalogChieuMoiLangDung4Chieu khang dinh moi lang trong 5 lang co DUNG 4 MauChieu
// (dac ta §4.3: 4 chieu x 5 lang = 20).
func TestCatalogChieuMoiLangDung4Chieu(t *testing.T) {
	for idLang := byte(0); idLang < 5; idLang++ {
		chieu := ChieuCuaLang(idLang)
		if len(chieu) != 4 {
			t.Errorf("lang %d co %d chieu, mong doi DUNG 4", idLang, len(chieu))
		}
	}
}

// TestCatalogChieuKhongTemplateIdBang0 khang dinh KHONG chieu nao co templateId ==
// 0 — day la cach ne nhanh nhoi 2 cap gia 105/106 cua ban Java (dac ta §3.7
// diem 2, brief Task 1 buoc 3).
func TestCatalogChieuKhongTemplateIdBang0(t *testing.T) {
	for idLang := byte(0); idLang < 5; idLang++ {
		for _, m := range ChieuCuaLang(idLang) {
			if m.Id == 0 {
				t.Errorf("lang %d co chieu %q templateId == 0 — se kich nhanh nhoi cap gia 105/106 cua ban Java", idLang, m.Ten)
			}
		}
	}
}

// TestCatalogChieuMoiLangDungMotHoiMau khang dinh moi lang co DUNG 1 chieu type=2
// (hoi mau dong doi). Day la khang dinh QUAN TRONG NHAT cua Task 1: goi -60
// dang hoi mau (Task 10, ham goiHoiMauNguoiChoi) CHI sinh ra tu chieu type=2 —
// thieu no thi ham do khong co duong chay that (§4.3 dac ta).
func TestCatalogChieuMoiLangDungMotHoiMau(t *testing.T) {
	const loaiHoiMau byte = 2
	for idLang := byte(0); idLang < 5; idLang++ {
		dem := 0
		for _, m := range ChieuCuaLang(idLang) {
			if m.Loai == loaiHoiMau {
				dem++
			}
		}
		if dem != 1 {
			t.Errorf("lang %d co %d chieu type=2 (hoi mau), mong doi DUNG 1", idLang, dem)
		}
	}
}

// TestCatalogChieuMoiChieuDungBaCap khang dinh moi chieu co dung 3 cap (dac ta §4.3: "moi
// chieu 3 cap").
func TestCatalogChieuMoiChieuDungBaCap(t *testing.T) {
	for idLang := byte(0); idLang < 5; idLang++ {
		for _, m := range ChieuCuaLang(idLang) {
			if len(m.CacCap) != 3 {
				t.Errorf("lang %d chieu %q co %d cap, mong doi 3", idLang, m.Ten, len(m.CacCap))
			}
		}
	}
}

// TestCatalogChieuGoiCatalogSkillGhiDungHet khang dinh goiCatalogSkill() (du lieu
// THAT, 20 chieu) ghi ra khong loi va DOC LAI duoc het buffer bang Reader —
// bat duoc loi khuon byte tren toan bo du lieu that, khong chi tren fixture
// nho o TestCatalogChieuMotChieuMotCap.
func TestCatalogChieuGoiCatalogSkillGhiDungHet(t *testing.T) {
	m, err := goiCatalogSkill()
	if err != nil {
		t.Fatalf("goiCatalogSkill loi: %v", err)
	}
	r := protocol.NewReader(m.Data)
	if _, err := r.Byte(); err != nil { // lenh con
		t.Fatalf("doc lenh con loi: %v", err)
	}
	if _, err := r.Byte(); err != nil { // vsSkill
		t.Fatalf("doc vsSkill loi: %v", err)
	}
	if _, err := r.Byte(); err != nil { // countOption
		t.Fatalf("doc countOption loi: %v", err)
	}
	countClass, err := r.Byte()
	if err != nil {
		t.Fatalf("doc countClass loi: %v", err)
	}
	for i := 0; i < int(countClass); i++ {
		if _, err := r.UTF(); err != nil {
			t.Fatalf("doc ten lang %d loi: %v", i, err)
		}
		countTemplate, err := r.Byte()
		if err != nil {
			t.Fatalf("doc countTemplate lang %d loi: %v", i, err)
		}
		for j := 0; j < int(countTemplate); j++ {
			if err := boQuaMotMauChieu(r); err != nil {
				t.Fatalf("doc chieu %d cua lang %d loi: %v", j, i, err)
			}
		}
	}
	if _, err := r.Byte(); err == nil {
		t.Error("goi -28/7 that con byte thua sau khi doc het 5 lang — sai khuon byte")
	}
}

// TestTraCapChieuFailClosed khoa cua tra cuu chieu dung chung cho nguoi choi va
// boss (Ke hoach 6b). Diem quan trong: MOI dau vao khong hop le deu tra nil chu
// KHONG am tham roi ve cap 1 — vi du lieu boss quen dien cap phai lo ra ngay.
func TestTraCapChieuFailClosed(t *testing.T) {
	// Lay mot chieu THAT tu danh muc lang 0 lam moc so sanh.
	danhMuc := ChieuCuaLang(0)
	if len(danhMuc) == 0 || len(danhMuc[0].CacCap) == 0 {
		t.Fatal("danh muc chieu lang 0 rong — fixture sai")
	}
	mau := danhMuc[0]
	soCap := byte(len(mau.CacCap))

	// Hop le: cap 1 phai tra DUNG con tro toi CacCap[0] (khong phai ban sao).
	if got := TraCapChieu(0, mau.Id, 1); got != &danhMuc[0].CacCap[0] {
		t.Fatalf("TraCapChieu(0, %d, 1) khong tra dung con tro CacCap[0]", mau.Id)
	}
	// Hop le: cap cuoi cung.
	if got := TraCapChieu(0, mau.Id, soCap); got != &danhMuc[0].CacCap[soCap-1] {
		t.Fatalf("TraCapChieu(0, %d, %d) khong tra dung cap cuoi", mau.Id, soCap)
	}

	xau := []struct {
		ten        string
		lang       byte
		templateId byte
		cap        byte
	}{
		{"cap 0 KHONG duoc hieu la cap 1", 0, mau.Id, 0},
		{"cap vuot so cap co that", 0, mau.Id, soCap + 1},
		{"cap cuc lon", 0, mau.Id, 255},
		{"lang ngoai bang (5)", 5, mau.Id, 1},
		{"lang ngoai bang (255)", 255, mau.Id, 1},
		{"templateId khong ton tai", 0, 250, 1},
		{"templateId = 0 (id bat dau tu 1)", 0, 0, 1},
	}
	for _, c := range xau {
		t.Run(c.ten, func(t *testing.T) {
			if got := TraCapChieu(c.lang, c.templateId, c.cap); got != nil {
				t.Fatalf("TraCapChieu(%d, %d, %d) = %+v, phai la nil (fail-closed)",
					c.lang, c.templateId, c.cap, got)
			}
		})
	}

	// Moi lang hop le (0..4) deu phai tra duoc chieu dau tien cua no — chan
	// truong hop bang chieu chi nap cho mot lang.
	for lang := byte(0); lang <= 4; lang++ {
		dm := ChieuCuaLang(lang)
		if len(dm) == 0 {
			t.Fatalf("lang %d khong co chieu nao", lang)
		}
		if got := TraCapChieu(lang, dm[0].Id, 1); got == nil {
			t.Fatalf("lang %d: TraCapChieu(%d, 1) tra nil, phai co", lang, dm[0].Id)
		}
	}
}
