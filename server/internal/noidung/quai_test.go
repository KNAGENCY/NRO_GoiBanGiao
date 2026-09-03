// quai_test.go — bo kiem cua quai.yaml (KH8 task 8).
//
// Moi test o day la MOT DOT BIEN cu the tren file YAML THAT (van lay tu ban
// nhung), khong phai mot du lieu bia. Go mot cau kiem trong nap.go ra thi test
// tuong ung DO.
package noidung

import (
	"strings"
	"testing"
)

// dongSatThuong la CHUOI DUY NHAT trong quai.yaml chua gia tri sat_thuong_co_ban
// (dong khai that, khong phai dong chu thich nhac lai no). Dot bien phai va vao
// dong nay, khong duoc va vao chu thich.
const dongSatThuong = "sat_thuong_co_ban: 5 #"

// TestQuai_VanThatNapDuoc — luoi an toan: moi test duoi day deu bat dau tu van
// that, nen van that phai nap duoc, neu khong ca tep nay xanh vi ly do sai.
//
// Cac con so o day la con so CU (game/quai.go truoc task 8, va assets/
// catalog.json). Task 8 chi chuyen cho o cua chung — lech mot so la task sai.
func TestQuai_VanThatNapDuoc(t *testing.T) {
	du, err := napVan(khoThat(t))
	if err != nil {
		t.Fatalf("van that khong nap duoc: %v", err)
	}
	if len(du.MauQuai) != 1 {
		t.Fatalf("mong doi dung 1 mau quai (Soi Rung), co %d", len(du.MauQuai))
	}
	m := du.MauQuai[0]
	if m.Ma != "soi-rung" || m.Hp != 100 || m.SatThuongCoBan != 5 || m.Cap != 1 || m.CapTinhAnh != 1 || m.AnhId != 0 {
		t.Fatalf("mau quai lech so cu: %+v", m)
	}
	// Sau ban do, moi ban 2 con — y het taoQuaiMacDinh cu.
	for mapID := int32(1); mapID <= 6; mapID++ {
		var dat []DiemSinhQuai
		for _, s := range du.SinhQuai {
			if s.Map == mapID {
				dat = s.Dat
			}
		}
		if len(dat) != 2 || dat[0].Id != 0 || dat[0].X != 160 || dat[0].Y != 384 ||
			dat[1].Id != 1 || dat[1].X != 40 || dat[1].Y != 384 {
			t.Fatalf("map %d: diem sinh lech ban cu: %+v", mapID, dat)
		}
	}
}

// TestQuai_MaTrungBiTuChoi — DOT BIEN: hai mau cung `ma`. Neu qua duoc, `sinh:`
// tro toi ma do se an nham mau nao la chuyen may rui.
func TestQuai_MaTrungBiTuChoi(t *testing.T) {
	van := khoThat(t)
	them := "  - ma: soi-rung\n    ten: Soi Rung Khac\n    anh_id: 0\n    hp: 50\n" +
		"    sat_thuong_co_ban: 3\n    cap: 1\n    cap_tinh_anh: 1\n"
	doiMotLan(t, van, TenFileQuai, "\nsinh:\n", "\n"+them+"\nsinh:\n")
	err := phaiHong(t, van)
	if !strings.Contains(err.Error(), "hai lan") {
		t.Errorf("loi bao %q — mong doi noi ro ma khai hai lan", err)
	}
}

// TestQuai_HpKhongDuongBiTuChoi — DOT BIEN: hp 0 (va am). Quai sinh ra da chet
// san se hoi sinh rồi chet vong tron moi tick ma khong mot dong log.
func TestQuai_HpKhongDuongBiTuChoi(t *testing.T) {
	for _, hp := range []string{"0", "-1"} {
		van := khoThat(t)
		doiMotLan(t, van, TenFileQuai, "hp: 100 #", "hp: "+hp+" #")
		err := phaiHong(t, van)
		if !strings.Contains(err.Error(), "hp") {
			t.Errorf("hp = %s: loi bao %q — mong doi noi ve hp", hp, err)
		}
	}
}

// TestQuai_SatThuongKhongDuongBiTuChoi — DOT BIEN quan trong nhat cua task 8:
// dua `sat_thuong_co_ban` ve 0 (hay am) khong duoc phep im lang di qua.
func TestQuai_SatThuongKhongDuongBiTuChoi(t *testing.T) {
	for _, st := range []string{"0", "-5"} {
		van := khoThat(t)
		doiMotLan(t, van, TenFileQuai, dongSatThuong, "sat_thuong_co_ban: "+st+" #")
		err := phaiHong(t, van)
		if !strings.Contains(err.Error(), "sat_thuong_co_ban") {
			t.Errorf("sat_thuong = %s: loi bao %q — mong doi noi ve sat_thuong_co_ban", st, err)
		}
	}
}

// TestQuai_DiemSinhTroToiMaLaBiTuChoi — DOT BIEN: doi `ma` cua mau ma quen sua
// khoi `sinh:`. Khong co kiem nay thi ban do ra RONG QUAI, im lang.
func TestQuai_DiemSinhTroToiMaLaBiTuChoi(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileQuai, "  - ma: soi-rung\n    ten:", "  - ma: soi-nui\n    ten:")
	err := phaiHong(t, van)
	if !strings.Contains(err.Error(), "khong co trong `mau_quai`") {
		t.Errorf("loi bao %q — mong doi noi diem sinh tro toi ma khong ton tai", err)
	}
}

// TestQuai_TempIdTrungTrongMotMap — DOT BIEN: hai con cung tempId. Goi 54 (xin
// danh) se luon trung con dau; con kia bat tu vinh vien.
func TestQuai_TempIdTrungTrongMotMap(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileQuai,
		"{ ma: soi-rung, id: 1, x: 40, y: 384 }",
		"{ ma: soi-rung, id: 0, x: 40, y: 384 }")
	err := phaiHong(t, van)
	if !strings.Contains(err.Error(), "tempId") {
		t.Errorf("loi bao %q — mong doi noi ve tempId trung", err)
	}
}

// TestQuai_MapKhaiHaiKhoiSinh — DOT BIEN: chep dan mot khoi map. Khoi thu hai
// bi bo qua im lang neu khong co kiem.
func TestQuai_MapKhaiHaiKhoiSinh(t *testing.T) {
	van := khoThat(t)
	van[TenFileQuai] += "  - map: 1\n    dat:\n      - { ma: soi-rung, id: 5, x: 10, y: 20 }\n"
	err := phaiHong(t, van)
	if !strings.Contains(err.Error(), "hai khoi") {
		t.Errorf("loi bao %q — mong doi noi map khai hai khoi sinh", err)
	}
}

// TestQuai_TenRongBiTuChoi
func TestQuai_TenRongBiTuChoi(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileQuai, "    ten: Sói Rừng\n", "    ten: \"\"\n")
	if err := phaiHong(t, van); !strings.Contains(err.Error(), "ten rong") {
		t.Errorf("loi bao %q — mong doi noi ten rong", err)
	}
}

// TestQuai_CapKhongBiTuChoi
func TestQuai_CapKhongBiTuChoi(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileQuai, "    cap: 1 #", "    cap: 0 #")
	if err := phaiHong(t, van); !strings.Contains(err.Error(), "cap") {
		t.Errorf("loi bao %q — mong doi noi ve cap", err)
	}
}

// TestQuai_TruongGoNhamTenBiTuChoi: KnownFields(true) phai ap dung cho ca file
// nay — `sat_thuog_co_ban` khong duoc lang le thanh 0.
func TestQuai_TruongGoNhamTenBiTuChoi(t *testing.T) {
	van := khoThat(t)
	doiMotLan(t, van, TenFileQuai, dongSatThuong, "sat_thuog_co_ban: 5 #")
	if err := phaiHong(t, van); !strings.Contains(err.Error(), TenFileQuai) {
		t.Errorf("loi bao %q — mong doi noi ro file quai.yaml sai cu phap", err)
	}
}

// ---------------------------------------------------------------------------
// KiemIdAnhQuai — ham THUAN, kiem duoc ma khong can toi catalog that
// ---------------------------------------------------------------------------

func TestKiemIdAnhQuai_IdLaBiTuChoi(t *testing.T) {
	err := KiemIdAnhQuai([]MauQuai{{Ma: "soi-rung", AnhId: 999}}, []int{0, 1, 2})
	if err == nil {
		t.Fatal("anh_id 999 khong co trong catalog ma van duoc chap nhan")
	}
	if !strings.Contains(err.Error(), "asset-lock") {
		t.Errorf("loi bao %q — mong doi nhac toi asset-lock (nguon su that cua id anh)", err)
	}
}

func TestKiemIdAnhQuai_IdCoThatThiQua(t *testing.T) {
	if err := KiemIdAnhQuai([]MauQuai{{Ma: "soi-rung", AnhId: 2}}, []int{0, 1, 2}); err != nil {
		t.Fatalf("anh_id co that ma bi tu choi: %v", err)
	}
}

// phaiHong nap van da doi va doi hoi mot loi. Fatal neu nap duoc — mot dot bien
// "song" phai lam test DO ngay tai day, khong duoc lot xuong phep so chuoi.
func phaiHong(t *testing.T, van map[string]string) error {
	t.Helper()
	du, err := napVan(van)
	if err == nil {
		t.Fatalf("van da bi dot bien VAN NAP DUOC (%d mau quai) — bo kiem khong bat", len(du.MauQuai))
	}
	return err
}
