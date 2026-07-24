package game

import "testing"

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
