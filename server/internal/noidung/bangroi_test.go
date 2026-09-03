package noidung

// bangroi_test.go — KH8 task 9.
//
// Hai nhóm:
//
//	A. Bộ kiểm bang-roi.yaml — mỗi test là một ĐỘT BIẾN dữ liệu cụ thể phải ĐỎ.
//	B. Cửa tra cứu trên ảnh chụp — mã lạ về bảng mặc định, bản sao sâu.

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Vật liệu dùng chung
// ---------------------------------------------------------------------------

func bangRoiMauVP() []MauVatPham {
	return []MauVatPham{
		{Id: 0, Ten: "Thuoc mau", TranChong: 99},
		{Id: 1, Ten: "Ao", TranChong: 1},
	}
}

func bangRoiKhoa() []khoaIdMuc {
	return []khoaIdMuc{{Id: 0, Ma: "thuoc-mau"}, {Id: 1, Ma: "ao"}}
}

func bangRoiMauQuai() []MauQuai {
	return []MauQuai{{Ma: "soi-rung", Ten: "Soi Rung", Hp: 100, SatThuongCoBan: 5, Cap: 1, CapTinhAnh: 1}}
}

// kiemBangRoiThu chạy bộ kiểm với một file bảng rơi tuỳ ý trên bộ vật liệu hợp lệ.
func kiemBangRoiThu(f *fileBangRoi) error {
	return kiemBangRoi(f, bangRoiMauQuai(), bangRoiMauVP(), bangRoiKhoa())
}

// ---------------------------------------------------------------------------
// A. Bộ kiểm
// ---------------------------------------------------------------------------

func TestKiemBangRoiBanHopLeThiXanh(t *testing.T) {
	f := &fileBangRoi{
		QuaiMacDinh: []MucRoi{{MauId: 0, SoLuong: 1, TiLePhanTram: 15}},
		Quai: []BangRoiTheoMa{{Ma: "soi-rung", Roi: []MucRoi{
			{MauId: 0, SoLuong: 2, TiLePhanTram: 10},
		}}},
		Boss: []BangRoiTheoMa{{Ma: MaBangRoiBossChung, Roi: []MucRoi{
			{MauId: 0, SoLuong: 1, ChacChan: true},
			{MauId: 1, SoLuong: 1, TiLePhanTram: 50, ChiBossLon: true},
		}}},
	}
	if err := kiemBangRoiThu(f); err != nil {
		t.Fatalf("ban hop le bi tu choi: %v", err)
	}
}

// ⚠️ ĐỘT BIẾN: bỏ vòng đối chiếu `maQuaiCo[k.Ma]` trong kiemBangRoi.
func TestKiemBangRoiMaQuaiLaThiDo(t *testing.T) {
	f := &fileBangRoi{Quai: []BangRoiTheoMa{{Ma: "ho-mo-cau", Roi: []MucRoi{
		{MauId: 0, SoLuong: 1, TiLePhanTram: 10},
	}}}}
	err := kiemBangRoiThu(f)
	if err == nil {
		t.Fatal("ma quai khong co trong mau_quai PHAI bi tu choi nap")
	}
	if !strings.Contains(err.Error(), "ho-mo-cau") {
		t.Fatalf("thong bao loi khong noi ro ma nao sai: %v", err)
	}
}

// ⚠️ ĐỘT BIẾN: bỏ kiểm `daKhoa[m.MauId]` (chỉ còn kiểm "tồn tại").
//
// mau_id 2 KHÔNG có trong danh mục lẫn sổ khoá; để tách bạch hai điều kiện, test
// thứ hai dựng một mẫu CÓ trong danh mục nhưng CHƯA vào sổ khoá.
func TestKiemBangRoiMauIdChuaKhoaThiDo(t *testing.T) {
	f := &fileBangRoi{QuaiMacDinh: []MucRoi{{MauId: 2, SoLuong: 1, TiLePhanTram: 10}}}

	// Danh mục CÓ mẫu 2, sổ khoá thì KHÔNG.
	mauVP := append(bangRoiMauVP(), MauVatPham{Id: 2, Ten: "Mau moi chua ghi so", TranChong: 1})
	err := kiemBangRoi(f, bangRoiMauQuai(), mauVP, bangRoiKhoa())
	if err == nil {
		t.Fatal("mau_id ton tai nhung CHUA KHOA phai bi tu choi")
	}
	if !strings.Contains(err.Error(), TenFileKhoaId) {
		t.Fatalf("thong bao loi phai chi toi so khoa: %v", err)
	}

	// Và mẫu hoàn toàn không tồn tại thì cũng đỏ (kiểm còn lại vẫn sống).
	if kiemBangRoiThu(f) == nil {
		t.Fatal("mau_id khong ton tai phai bi tu choi")
	}
}

// ⚠️ ĐỘT BIẾN: bỏ kiểm `so_luong <= tran_chong`.
func TestKiemBangRoiSoLuongVuotTranChongThiDo(t *testing.T) {
	// Mẫu 1 ("Ao") có tran_chong = 1.
	f := &fileBangRoi{QuaiMacDinh: []MucRoi{{MauId: 1, SoLuong: 2, TiLePhanTram: 10}}}
	err := kiemBangRoiThu(f)
	if err == nil {
		t.Fatal("so_luong 2 > tran_chong 1 phai bi tu choi")
	}
	if !strings.Contains(err.Error(), "tran_chong") {
		t.Fatalf("thong bao loi khong noi ve tran chong: %v", err)
	}
}

func TestKiemBangRoiTiLeNgoaiMienThiDo(t *testing.T) {
	for _, ti := range []float64{-1, 100.5, 1000} {
		f := &fileBangRoi{QuaiMacDinh: []MucRoi{{MauId: 0, SoLuong: 1, TiLePhanTram: ti}}}
		if kiemBangRoiThu(f) == nil {
			t.Fatalf("ti_le_phan_tram = %g phai bi tu choi", ti)
		}
	}
}

func TestKiemBangRoiSoLuongKhongDuongThiDo(t *testing.T) {
	f := &fileBangRoi{QuaiMacDinh: []MucRoi{{MauId: 0, SoLuong: 0, TiLePhanTram: 10}}}
	if kiemBangRoiThu(f) == nil {
		t.Fatal("so_luong = 0 phai bi tu choi")
	}
}

func TestKiemBangRoiKhaiCaChacChanLanTiLeThiDo(t *testing.T) {
	f := &fileBangRoi{QuaiMacDinh: []MucRoi{{MauId: 0, SoLuong: 1, TiLePhanTram: 50, ChacChan: true}}}
	if kiemBangRoiThu(f) == nil {
		t.Fatal("khai ca chac_chan lan ti_le_phan_tram phai bi tu choi")
	}
}

func TestKiemBangRoiChiBossLonTrongKhoiQuaiThiDo(t *testing.T) {
	f := &fileBangRoi{QuaiMacDinh: []MucRoi{{MauId: 0, SoLuong: 1, TiLePhanTram: 10, ChiBossLon: true}}}
	if kiemBangRoiThu(f) == nil {
		t.Fatal("chi_boss_lon trong khoi quai la mot dong khong bao gio chay -> phai bi tu choi")
	}
}

func TestKiemBangRoiMaKhaiHaiLanThiDo(t *testing.T) {
	muc := []MucRoi{{MauId: 0, SoLuong: 1, TiLePhanTram: 10}}
	f := &fileBangRoi{Quai: []BangRoiTheoMa{{Ma: "soi-rung", Roi: muc}, {Ma: "soi-rung", Roi: muc}}}
	if kiemBangRoiThu(f) == nil {
		t.Fatal("mot ma khai hai lan phai bi tu choi (khoi sau de len khoi truoc im lang)")
	}
}

// KiemMaBangRoiBoss: mã chung luôn hợp lệ, mã lạ thì đỏ.
func TestKiemMaBangRoiBoss(t *testing.T) {
	du := &DuLieu{BangRoi: BangRoi{Boss: map[string][]MucRoi{MaBangRoiBossChung: nil}}}
	if err := KiemMaBangRoiBoss(du, nil); err != nil {
		t.Fatalf("ma %q phai luon hop le: %v", MaBangRoiBossChung, err)
	}

	du.BangRoi.Boss["nhan-gia-mat-na"] = nil
	if err := KiemMaBangRoiBoss(du, []string{"nhan-gia-mat-na"}); err != nil {
		t.Fatalf("ma boss that phai hop le: %v", err)
	}
	if KiemMaBangRoiBoss(du, nil) == nil {
		t.Fatal("ma boss khong co trong danh sach phai bi bao loi")
	}
}

// ---------------------------------------------------------------------------
// B. Bản NHÚNG thật + cửa tra cứu
// ---------------------------------------------------------------------------

func TestBanNhungCoBangRoi(t *testing.T) {
	du, err := NapNhung()
	if err != nil {
		t.Fatalf("nap ban nhung: %v", err)
	}
	if len(du.BangRoi.QuaiMacDinh) == 0 {
		t.Fatal("ban nhung khong co bang roi quai mac dinh")
	}
	if len(du.BangRoi.Quai) == 0 {
		t.Fatal("ban nhung khong co bang roi quai rieng nao")
	}
	if len(du.BangRoi.Boss[MaBangRoiBossChung]) == 0 {
		t.Fatalf("ban nhung khong co bang roi boss %q", MaBangRoiBossChung)
	}
	// Mọi mã boss trong bản nhúng phải hợp lệ với danh sách rỗng (chỉ có mã chung).
	if err := KiemMaBangRoiBoss(du, nil); err != nil {
		t.Fatalf("ban nhung co ma boss la: %v", err)
	}
}

func TestAnhChupBangRoiMaLaThiVeBangMacDinh(t *testing.T) {
	du, err := NapNhung()
	if err != nil {
		t.Fatalf("nap ban nhung: %v", err)
	}
	ac := MoKho(du).AnhChup()

	la := ac.BangRoiCuaQuai("khong-co-that")
	if len(la) != len(du.BangRoi.QuaiMacDinh) {
		t.Fatalf("ma la phai ve bang mac dinh (%d muc), nhan %d", len(du.BangRoi.QuaiMacDinh), len(la))
	}
	if ac.BangRoiCuaBoss("khong-co-that") != nil {
		t.Fatal("boss khong co bang -> nil (khong co bang mac dinh cho boss)")
	}
}

// ⚠️ Ảnh chụp phải trả BẢN SAO: sửa kết quả không được chạm vào ảnh chụp.
func TestAnhChupBangRoiTraBanSao(t *testing.T) {
	du, err := NapNhung()
	if err != nil {
		t.Fatalf("nap ban nhung: %v", err)
	}
	ac := MoKho(du).AnhChup()

	a := ac.BangRoiCuaQuai("soi-rung")
	if len(a) == 0 {
		t.Fatal("soi-rung phai co bang roi rieng")
	}
	goc := a[0].SoLuong
	a[0].SoLuong = goc + 999

	b := ac.BangRoiCuaQuai("soi-rung")
	if b[0].SoLuong != goc {
		t.Fatalf("anh chup bi ban: so_luong %d -> %d", goc, b[0].SoLuong)
	}
}
