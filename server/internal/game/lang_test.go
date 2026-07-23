package game

import "testing"

func TestTimLangHopLe(t *testing.T) {
	l := TimLang(0)
	if l == nil {
		t.Fatal("lang 0 (Lang La) phai ton tai")
	}
	if l.Ten != "Làng Lá" {
		t.Errorf("ten lang 0 = %q, mong doi \"Làng Lá\"", l.Ten)
	}
	// Moi lang phai co chi so khoi dau duong — nhan vat khong the co 0 mau.
	for _, l := range DanhSachLang() {
		if l.HpKhoiDau <= 0 || l.ChakraKhoiDau <= 0 {
			t.Errorf("lang %d (%s) co hp/chakra khoi dau khong duong", l.ID, l.Ten)
		}
	}
}

func TestTimLangKhongTonTai(t *testing.T) {
	if TimLang(99) != nil {
		t.Error("lang 99 khong ton tai, phai tra nil")
	}
	if TimLang(-1) != nil {
		t.Error("lang -1 khong ton tai, phai tra nil")
	}
}

func TestCoDuNamLang(t *testing.T) {
	// Thiet ke §7.2 co 5 lang. Client cu chi chon duoc 0/1/2, nhung may chu
	// dinh nghia san ca 5 de client dung lai sau khong phai sua may chu.
	if n := len(DanhSachLang()); n != 5 {
		t.Errorf("co %d lang, mong doi 5", n)
	}
}
