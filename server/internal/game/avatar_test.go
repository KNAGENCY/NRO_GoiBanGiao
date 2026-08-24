package game

import (
	"encoding/json"
	"testing"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// O chan dung tren HUD lay anh qua bang anh xa headId -> avatarId ma may chu
// gui trong ME_INFO. Truoc day may chu gui headAvatar.size = 0 nen bang rong,
// Char.getAvatar tra -1 va client bo qua khong ve gi — o chan dung trong tron.
//
// Client doc bang nay o Controller.cs:5978 (short soCap, roi tung cap
// short headId, short avatarId).

// docDenHeadAvatar doc luot qua cac truong dau cua ME_INFO de con tro dung
// ngay truoc headAvatar.size. Thu tu chep tu goiMeInfo — doi thu tu ben do ma
// quen doi o day thi test se bao loi doc, khong am tham qua.
func docDenHeadAvatar(t *testing.T, data []byte) *protocol.Reader {
	t.Helper()
	r := protocol.NewReader(data)
	buoc := func(ten string, f func() error) {
		if err := f(); err != nil {
			t.Fatalf("doc %s loi: %v", ten, err)
		}
	}
	b := func(ten string) { buoc(ten, func() error { _, err := r.Byte(); return err }) }
	sh := func(ten string) { buoc(ten, func() error { _, err := r.Short(); return err }) }
	in := func(ten string) { buoc(ten, func() error { _, err := r.Int(); return err }) }
	lo := func(ten string) { buoc(ten, func() error { _, err := r.Long(); return err }) }

	b("lenh con")
	in("id")
	b("taskMain.id")
	b("gender")
	sh("head")
	buoc("ten", func() error { _, err := r.UTF(); return err })
	b("cPK")
	b("typePk")
	lo("power")
	sh("buffHp")
	sh("buffMp")
	b("gender lan 2")
	b("skillCount")
	lo("gold")
	in("ruby")
	in("gem")
	b("itemsBody.size")
	b("itemsBag.size")
	b("itemsBox.size")
	return r
}

func TestMeInfoGuiBangAnhXaChanDung(t *testing.T) {
	nv := &store.NhanVat{ID: 1, Ten: "Naruto", Lang: 0, KieuToc: 30}
	m, err := goiMeInfo(nv)
	if err != nil {
		t.Fatalf("goiMeInfo: %v", err)
	}
	r := docDenHeadAvatar(t, m.Data)

	soCap, err := r.Short()
	if err != nil {
		t.Fatalf("doc headAvatar.size loi: %v", err)
	}
	if soCap == 0 {
		t.Fatal("headAvatar.size = 0 -> bang anh xa rong -> o chan dung tren HUD trong tron")
	}

	thay := false
	for i := int16(0); i < soCap; i++ {
		headID, err1 := r.Short()
		avatarID, err2 := r.Short()
		if err1 != nil || err2 != nil {
			t.Fatalf("doc cap thu %d loi: %v %v", i, err1, err2)
		}
		if headID == HeadNhanVat {
			thay = true
			if avatarID != AvatarNhanVat {
				t.Errorf("head %d anh xa toi avatar %d, mong doi %d", headID, avatarID, AvatarNhanVat)
			}
		}
	}
	if !thay {
		t.Errorf("bang anh xa khong co head %d (dau nhan vat nguoi choi)", HeadNhanVat)
	}
}

func TestAnhChanDungKhacAnhDau(t *testing.T) {
	// Neu tro nham chan dung ve chinh anh part dau thi o HUD se hien mot cai dau
	// bi cat cup, khong phai chan dung. Hai id phai khac nhau.
	if AvatarNhanVat == HeadNhanVat {
		t.Error("AvatarNhanVat trung HeadNhanVat")
	}
	if AvatarNhanVat <= 0 || HeadNhanVat <= 0 {
		t.Error("id chan dung / dau phai duong")
	}
}

func TestChanDungNguoiChoiKhongTrungChanDungNPC(t *testing.T) {
	// Nguoi choi va NPC khong duoc dung chung mot khuon mat. Doc thang catalog
	// de bat duoc ngay khi ai do doi portrait cua NPC thanh trung voi nguoi choi.
	var cat struct {
		NPCs []struct {
			Portrait int16  `json:"portrait"`
			Ten      string `json:"name"`
		} `json:"npcs"`
	}
	if err := json.Unmarshal(duLieuCatalog, &cat); err != nil {
		t.Fatalf("doc catalog loi: %v", err)
	}
	if len(cat.NPCs) == 0 {
		t.Skip("catalog chua co NPC nao de doi chieu")
	}
	for _, npc := range cat.NPCs {
		if npc.Portrait == AvatarNhanVat {
			t.Errorf("chan dung nguoi choi (%d) trung chan dung NPC %q — hai ben se cung mot khuon mat",
				AvatarNhanVat, npc.Ten)
		}
	}
}
