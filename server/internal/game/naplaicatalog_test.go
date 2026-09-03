// naplaicatalog_test.go — NO CUA TASK 3, tra o task 10: goi catalog -28/7 luc
// DANG NHAP phai doc kho noi dung DA TIEM, khong phai ban NHUNG.
//
// Vi sao dieu nay quan trong: nap nong (noidung/kich.go) goi Kho.Thay tren kho
// TIEM. Neu duong dang nhap con doc ban nhung thi sau mot lan nap nong, nguoi
// dang choi (actor doc kho tiem) an so MOI trong khi nguoi vua dang nhap nhan
// catalog voi so CU — hai su that cung luc tren cung mot may chu.
package game

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"nroserver/internal/noidung"
	"nroserver/internal/protocol"
)

// chakraDauTienTuGoiSkill doc NGUOC BYTE goi -28/7 va tra ve manaUse cua cap
// DAU TIEN cua chieu DAU TIEN cua lang 0 — tuc `chakra_dung` cua chieu do
// trong chieu.yaml.
func chakraDauTienTuGoiSkill(t *testing.T, m *protocol.Message) int16 {
	t.Helper()
	r := protocol.NewReader(m.Data)
	for _, ten := range []string{"lenh con", "vsSkill", "countOption"} {
		if _, err := r.Byte(); err != nil {
			t.Fatalf("doc %s that bai: %v", ten, err)
		}
	}
	if _, err := r.Byte(); err != nil { // countClass
		t.Fatalf("doc countClass that bai: %v", err)
	}
	if _, err := r.UTF(); err != nil { // ten lang 0
		t.Fatalf("doc ten lang 0 that bai: %v", err)
	}
	if _, err := r.Byte(); err != nil { // countTemplate
		t.Fatalf("doc countTemplate that bai: %v", err)
	}
	// Mot MauChieu: templateId, ten, maxPoint, manaUseType, type, iconId,
	// damInfo, ghiChu, countSkill (cung khuon voi boQuaMotMauChieu).
	if _, err := r.Byte(); err != nil {
		t.Fatalf("doc templateId that bai: %v", err)
	}
	if _, err := r.UTF(); err != nil {
		t.Fatalf("doc ten chieu that bai: %v", err)
	}
	for i := 0; i < 3; i++ {
		if _, err := r.Byte(); err != nil {
			t.Fatalf("doc byte thu %d cua mau chieu that bai: %v", i, err)
		}
	}
	if _, err := r.Short(); err != nil { // iconId
		t.Fatalf("doc iconId that bai: %v", err)
	}
	if _, err := r.UTF(); err != nil { // damInfo
		t.Fatalf("doc damInfo that bai: %v", err)
	}
	if _, err := r.UTF(); err != nil { // ghiChu
		t.Fatalf("doc ghiChu that bai: %v", err)
	}
	if _, err := r.Byte(); err != nil { // countSkill
		t.Fatalf("doc countSkill that bai: %v", err)
	}
	if _, err := r.Short(); err != nil { // skillId
		t.Fatalf("doc skillId that bai: %v", err)
	}
	if _, err := r.Byte(); err != nil { // point
		t.Fatalf("doc point that bai: %v", err)
	}
	if _, err := r.Long(); err != nil { // powRequire
		t.Fatalf("doc powRequire that bai: %v", err)
	}
	mana, err := r.Short()
	if err != nil {
		t.Fatalf("doc manaUse that bai: %v", err)
	}
	return mana
}

// thuMucChieuDaSua do kho noi dung nhung ra thu muc tam roi doi chakra_dung cua
// cap dau tien thanh chakraMoi.
func thuMucChieuDaSua(t *testing.T, chakraMoi string) string {
	t.Helper()
	d := t.TempDir()
	for _, ten := range noidung.TenCacFile() { // nguon duy nhat — xem chu thich noidung.TenCacFile
		b, err := os.ReadFile(filepath.Join("..", "..", "kho-noi-dung", ten))
		if err != nil {
			t.Fatalf("doc %s: %v", ten, err)
		}
		van := string(b)
		if ten == noidung.TenFileChieu {
			const cu = "chakra_dung: 10"
			if !strings.Contains(van, cu) {
				t.Fatalf("khong tim thay %q trong chieu.yaml — dot bien khong va vao dau", cu)
			}
			van = strings.Replace(van, cu, "chakra_dung: "+chakraMoi, 1)
		}
		if err := os.WriteFile(filepath.Join(d, ten), []byte(van), 0o600); err != nil {
			t.Fatalf("ghi %s: %v", ten, err)
		}
	}
	return d
}

// ⚠️ DOT BIEN: doi BoXuLy.anhChupNoiDung thanh "luon tra khoChieuMacDinh()"
// (tuc duong dang nhap doc lai ban NHUNG) -> test nay DO o dong cuoi.
func TestCatalogDangNhapDocKhoDaTiemSauNapNong(t *testing.T) {
	du, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("NapNhung: %v", err)
	}
	kho := noidung.MoKho(du)

	b := NewBoXuLy(nil, nil)
	b.DatKhoNoiDung(kho)

	m, err := b.goiCatalogSkillPhien()
	if err != nil {
		t.Fatalf("goiCatalogSkillPhien: %v", err)
	}
	truoc := chakraDauTienTuGoiSkill(t, m)
	if truoc != 10 {
		t.Fatalf("chakra cap 1 cua chieu dau tien = %d, mong 10 (theo chieu.yaml dang commit)", truoc)
	}

	duMoi, err := noidung.NapTuThuMuc(thuMucChieuDaSua(t, "77"))
	if err != nil {
		t.Fatalf("NapTuThuMuc sau khi sua: %v", err)
	}
	kho.Thay(duMoi)

	// Nguoi dang nhap MOI (mot BoXuLy moi tinh, cung kho tiem) phai thay so MOI.
	bMoi := NewBoXuLy(nil, nil)
	bMoi.DatKhoNoiDung(kho)
	mMoi, err := bMoi.goiCatalogSkillPhien()
	if err != nil {
		t.Fatalf("goiCatalogSkillPhien sau nap nong: %v", err)
	}
	if sau := chakraDauTienTuGoiSkill(t, mMoi); sau != 77 {
		t.Fatalf("sau nap nong, catalog dang nhap van tra chakra = %d, mong 77 — duong dang nhap dang doc ban NHUNG chu khong phai kho da tiem", sau)
	}
}

// Chua tiem kho thi phai roi ve ban NHUNG, khong panic — hang chuc test cu goi
// NewBoXuLy truc tiep va khong biet gi ve kho noi dung.
func TestCatalogDangNhapChuaTiemKhoThiDungBanNhung(t *testing.T) {
	b := NewBoXuLy(nil, nil)
	m, err := b.goiCatalogSkillPhien()
	if err != nil {
		t.Fatalf("goiCatalogSkillPhien khi chua tiem kho: %v", err)
	}
	if got := chakraDauTienTuGoiSkill(t, m); got != 10 {
		t.Fatalf("ban nhung cho chakra = %d, mong 10", got)
	}
}
