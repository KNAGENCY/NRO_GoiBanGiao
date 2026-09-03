package game

// bangroi_noidung_test.go — KH8 task 9: bảng rơi đọc TỪ KHO NỘI DUNG.
//
// Ba nhóm, mỗi nhóm nhắm một đột biến cụ thể:
//
//	A. Số trong YAML là số THẬT được dùng — đổi `ti_le_phan_tram` trong một bản
//	   nạp nóng thì số món yêu cầu rơi đổi theo. (Đột biến: roiDoCuaQuai đọc lại
//	   một bảng cứng nào đó / bỏ qua tỉ lệ.)
//	B. Hợp đồng "một lần đầu nhịp" vẫn đúng khi quái/boss CHẾT. (Đột biến:
//	   roiDoCuaQuai/roiDoCuaBoss gọi AnhChup()/layAnhChupNhip() thêm.)
//	C. Tra bảng rơi bằng `Ma`, KHÔNG bằng `MauId` (anh_id). (Đột biến: quay về
//	   tra theo MauId — hai mẫu quái cùng anh_id sẽ dùng chung một bảng.)

import (
	"testing"

	"nroserver/internal/noidung"
	"nroserver/internal/store"
)

// ---------------------------------------------------------------------------
// Vật liệu: một kho nội dung thật, có sửa bảng rơi
// ---------------------------------------------------------------------------

// khoBangRoiSua nạp bản NHÚNG rồi để người gọi sửa bảng rơi trước khi mở kho.
func khoBangRoiSua(t *testing.T, sua func(*noidung.DuLieu)) *noidung.Kho {
	t.Helper()
	du, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("nap kho noi dung: %v", err)
	}
	sua(du)
	return noidung.MoKho(du)
}

// ---------------------------------------------------------------------------
// A. Số trong YAML là số thật
// ---------------------------------------------------------------------------

// ⚠️ ĐỘT BIẾN NHẮM TỚI: roiDoCuaQuai trả một bảng cố định (hoặc bỏ qua
// ti_le_phan_tram). Ở đây tỉ lệ của MỘT dòng bị hạ về 0 -> số món ép-rơi-100%
// phải giảm đúng 1.
func TestBangRoiDoiTiLeTrongYamlThiSoMonYeuCauDoi(t *testing.T) {
	kho := khoBangRoiSua(t, func(du *noidung.DuLieu) {
		ds := du.BangRoi.QuaiMacDinh
		if len(ds) < 2 {
			t.Fatalf("bang mac dinh chi co %d dong, test can >= 2", len(ds))
		}
		ds[0].TiLePhanTram = 0 // tắt đúng một dòng
		ds[0].ChacChan = false
	})

	goc := len(roiDoCuaQuai(khoChieuMacDinh(), "ma-la-hoan-toan"))
	if goc < 2 {
		t.Fatalf("ban nhung: bang mac dinh %d dong, test can >= 2", goc)
	}

	k := MoKhuVuc(TimBanDo(1))
	k.GanKhoNoiDung(kho)
	k.roiDoNgauNhien = luonRoi
	bo := MoBoRoiDo(&khoRoiGia{}, 32)
	k.GanBoRoi(bo)

	bang := roiDoCuaQuai(k.layAnhChupNhip(), "ma-la-hoan-toan")
	n := k.roiDoKhiChet(0, 0, bang)
	if n != goc-1 {
		t.Fatalf("ha ti le mot dong ve 0 -> phai xep %d yeu cau, xep %d "+
			"(so trong bang-roi.yaml khong duoc dung toi?)", goc-1, n)
	}
}

// Chắc chắn = 100%: gieo XUI nhất vẫn phải rơi. Đây là nửa còn lại của việc
// `chac_chan: true` được dịch đúng sang tỉ lệ mẫu số.
func TestBangRoiChacChanThiGieoXuiVanRoi(t *testing.T) {
	ac := khoChieuMacDinh()
	bang := roiDoCuaBoss(ac, noidung.MaBangRoiBossChung, &PhaBoss{})
	soChacChan := 0
	for _, m := range bang {
		if m.TiLe == roiDoMauSo {
			soChacChan++
		}
	}
	if soChacChan == 0 {
		t.Fatal("bang boss trong bang-roi.yaml phai co it nhat mot muc chac_chan")
	}
	ra := quyetDinhRoiDo(bang, khongRoi)
	if len(ra) != soChacChan {
		t.Fatalf("gieo xui nhat: chi %d muc chac chan duoc roi, nhan %d", soChacChan, len(ra))
	}
}

// ---------------------------------------------------------------------------
// B. Hợp đồng "một lần đầu nhịp" khi quái/boss chết
// ---------------------------------------------------------------------------

// ⚠️ ĐỘT BIẾN NHẮM TỚI: thêm một lời gọi AnhChup()/layAnhChupNhip() vào đường
// rơi đồ (trong roiDoCuaQuai, roiDoKhiChet, hay ở chỗ gọi trong khuvuc.go).
// Một SỰ KIỆN giết chết một con quái vẫn phải lấy ĐÚNG MỘT ảnh chụp.
func TestKhoBangRoiMotLanKhiQuaiChet(t *testing.T) {
	d := moKhoChieuDem(t)

	k := MoKhuVuc(TimBanDo(1))
	k.ganNguonNoiDung(d.nguon())
	k.roiDoNgauNhien = luonRoi
	k.GanBoRoi(MoBoRoiDo(&khoRoiGia{}, 64))

	nv := &store.NhanVat{ID: 1, Ten: "A", Lang: 0, ChakraHienTai: 100, HpHienTai: 100}
	k.xuLy(suKienVao{nv: nv, phien: &phienGiaKhuVuc{}})
	k.xuLy(suKienChonChieu{idNhanVat: nv.ID, mauChieuId: 1})
	// Hp=1: chắc chắn chết sau đúng một đòn.
	k.quai = []*Quai{{Id: 0, Ma: "soi-rung", Hp: 1, HpToiDa: 150}}

	d.soLan = 0
	k.xuLy(suKienDanhQuai{idNhanVat: nv.ID, mobId: 0})

	if k.quai[0].Hp != 0 {
		t.Fatalf("chuan bi test sai: quai chua chet (Hp=%d)", k.quai[0].Hp)
	}
	if d.soLan != 1 {
		t.Fatalf("su kien giet quai lay anh chup %d lan, hop dong la DUNG 1 "+
			"(xem khoi chu thich dau internal/noidung/kho.go)", d.soLan)
	}
}

// Cùng hợp đồng, phía boss: một NHỊP giết boss vẫn đúng một ảnh chụp.
func TestKhoBangRoiMotLanKhiBossChet(t *testing.T) {
	luc := mocDanhBoss
	d := moKhoChieuDem(t)
	k := khuVucDemChieu(t, d, &luc)
	k.roiDoNgauNhien = luonRoi
	k.GanBoRoi(MoBoRoiDo(&khoRoiGia{}, 64))

	d.soLan = 0
	k.acNhip = k.layAnhChupNhip() // mở nhịp, đúng như nhipTim làm
	b := &Boss{Id: 1, PhaHienTai: 0, DuLieu: []PhaBoss{{Ten: "Thu linh", Cap: 1}},
		Hp: 0, HpToiDa: 100, TrangThai: BossHoatDong}
	k.bossChet(b, 0, k.dongHo())

	if d.soLan != 1 {
		t.Fatalf("nhip giet boss lay anh chup %d lan, mong doi 1", d.soLan)
	}
}

// ---------------------------------------------------------------------------
// C. Tra bằng `Ma`, KHÔNG bằng `MauId` (anh_id)
// ---------------------------------------------------------------------------

// ⚠️ ĐỘT BIẾN NHẮM TỚI: quay lại tra bảng rơi bằng `Quai.MauId`.
//
// Fixture dựng HAI mẫu quái CÙNG `anh_id` nhưng khác `ma`, và cho chúng hai
// bảng rơi khác hẳn nhau về số dòng. Nếu tra theo MauId thì hai con dùng chung
// một bảng -> hai số bằng nhau -> test đỏ.
func TestBangRoiTraTheoMaKhongTheoAnhId(t *testing.T) {
	kho := khoBangRoiSua(t, func(du *noidung.DuLieu) {
		mau := du.MauQuai[0]
		a := mau
		a.Ma = "sinh-doi-a"
		b := mau
		b.Ma = "sinh-doi-b"
		// CÙNG anh_id (chép nguyên từ mẫu gốc) — đây là điểm mấu chốt.
		if a.AnhId != b.AnhId {
			t.Fatalf("chuan bi test sai: hai mau phai cung anh_id")
		}
		du.MauQuai = append(du.MauQuai, a, b)

		du.BangRoi.Quai["sinh-doi-a"] = []noidung.MucRoi{
			{MauId: 0, SoLuong: 1, TiLePhanTram: 10},
		}
		du.BangRoi.Quai["sinh-doi-b"] = []noidung.MucRoi{
			{MauId: 0, SoLuong: 1, TiLePhanTram: 10},
			{MauId: 1, SoLuong: 1, TiLePhanTram: 10},
			{MauId: 4, SoLuong: 1, TiLePhanTram: 10},
		}
	})
	ac := kho.AnhChup()

	a := roiDoCuaQuai(ac, "sinh-doi-a")
	b := roiDoCuaQuai(ac, "sinh-doi-b")
	if len(a) != 1 || len(b) != 3 {
		t.Fatalf("hai mau CUNG anh_id phai co hai bang roi RIENG: %d va %d muc "+
			"(dang tra theo MauId/anh_id thay vi theo Ma?)", len(a), len(b))
	}
}

// ⚠️ ĐỘT BIẾN NHẮM TỚI (nửa còn lại, ở CHỖ GỌI trong khuvuc.go): tra bảng bằng
// `mucTieu.MauId` (anh_id) thay vì `mucTieu.Ma`.
//
// Con quái ở đây là "soi-rung" — mẫu CÓ bảng riêng, và bảng riêng đó ÍT dòng
// hơn bảng mặc định. Tra nhầm bằng MauId thì mã không khớp gì cả -> rơi về
// bảng mặc định -> số yêu cầu xếp nhiều hơn -> test đỏ.
func TestQuaiChetXepTheoBangRiengCuaMaChuKhongPhaiBangMacDinh(t *testing.T) {
	ac := khoChieuMacDinh()
	soRieng := len(roiDoCuaQuai(ac, "soi-rung"))
	soMacDinh := len(roiDoCuaQuai(ac, "ma-la-hoan-toan"))
	if soRieng == 0 || soRieng == soMacDinh {
		t.Fatalf("chuan bi test sai: bang rieng cua soi-rung (%d muc) phai khac bang mac dinh (%d muc)",
			soRieng, soMacDinh)
	}

	k := MoKhuVuc(TimBanDo(1))
	k.roiDoNgauNhien = luonRoi
	bo := MoBoRoiDo(&khoRoiGia{}, 64)
	k.GanBoRoi(bo)

	nv := &store.NhanVat{ID: 1, Ten: "A", Lang: 0, ChakraHienTai: 100, HpHienTai: 100}
	k.xuLy(suKienVao{nv: nv, phien: &phienGiaKhuVuc{}})
	k.xuLy(suKienChonChieu{idNhanVat: nv.ID, mauChieuId: 1})
	k.quai = []*Quai{{Id: 0, Ma: "soi-rung", MauId: 0, Hp: 1, HpToiDa: 150}}

	k.xuLy(suKienDanhQuai{idNhanVat: nv.ID, mobId: 0})

	if k.quai[0].Hp != 0 {
		t.Fatalf("chuan bi test sai: quai chua chet (Hp=%d)", k.quai[0].Hp)
	}
	if n := len(bo.hang); n != soRieng {
		t.Fatalf("quai %q chet phai xep %d yeu cau (bang RIENG cua ma), xep %d "+
			"— dang tra bang MauId/anh_id nen roi ve bang mac dinh (%d muc)?",
			"soi-rung", soRieng, n, soMacDinh)
	}
}

// taoQuaiTuAnhChup phải nạp `Ma` vào từng con quái — không có nó thì mọi con
// quái đều rơi về bảng mặc định và đột biến ở trên sống nhăn.
func TestTaoQuaiTuAnhChupNapMa(t *testing.T) {
	quai := taoQuaiTuAnhChup(khoChieuMacDinh(), TimBanDo(1))
	if len(quai) == 0 {
		t.Fatal("ban do 1 phai co quai")
	}
	for i, q := range quai {
		if q.Ma == "" {
			t.Fatalf("quai thu %d khong co Ma — bang roi se roi ve mac dinh cho MOI con", i)
		}
	}
}

// Danh sách mã boss của gói này phải khác rỗng — tham số cho
// noidung.KiemMaBangRoiBoss (xem chú thích MaBossHopLeChoBangRoi).
func TestMaBossHopLeChoBangRoiKhongRong(t *testing.T) {
	ma := MaBossHopLeChoBangRoi()
	if len(ma) == 0 {
		t.Fatal("khong co ma boss nao")
	}
	du, err := noidung.NapNhung()
	if err != nil {
		t.Fatalf("nap kho noi dung: %v", err)
	}
	if err := noidung.KiemMaBangRoiBoss(du, ma); err != nil {
		t.Fatalf("bang roi boss cua ban nhung lech kho noi dung boss: %v", err)
	}
}
