package game

// roido_test.go — KH7 Task 12.
//
// Bốn nhóm:
//
//	A. Bảng rơi + gieo tỉ lệ (hàm thuần, bộ chọn ngẫu nhiên tiêm được).
//	B. ACTOR: quái/boss chết -> xếp đúng số yêu cầu, và KHÔNG có lệnh CSDL nào
//	   chạy trong actor (kho giả đếm).
//	C. Đường về: suKienDoDaSinh -> món có trong doDuoiDat + gói 68 đúng byte.
//	D. Kho lỗi -> KHÔNG có món ma.
//
// Các test chốt chặn đột biến được đánh dấu ⚠️ ngay trên hàm.

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"nroserver/internal/noidung"
	"nroserver/internal/store"
)

// ---------------------------------------------------------------------------
// Kho giả
// ---------------------------------------------------------------------------

// khoRoiGia đếm số lần SinhDoDuoiDat được gọi và ghi lại tham số. Bộ đếm là
// atomic vì goroutine BoRoiDo có thể gọi nó song song với goroutine test.
type khoRoiGia struct {
	soLan  atomic.Int64
	idKe   atomic.Int64
	loi    error
	mauCuo atomic.Int32
	soLuoc atomic.Int32
}

func (k *khoRoiGia) SinhDoDuoiDat(_ context.Context, mauID int32, soLuong int32,
	_ store.ChiSoPhu) (int64, error) {
	k.soLan.Add(1)
	k.mauCuo.Store(mauID)
	k.soLuoc.Store(soLuong)
	if k.loi != nil {
		return 0, k.loi
	}
	return k.idKe.Add(1), nil
}

// luonRoi / khongRoi là hai bộ chọn ngẫu nhiên tiêm được, ép hai đầu mút.
func luonRoi(int) int    { return 0 }
func khongRoi(n int) int { return n - 1 }

// ---------------------------------------------------------------------------
// A. Bảng rơi + gieo tỉ lệ
// ---------------------------------------------------------------------------

// ⚠️ Tỉ lệ 0% KHÔNG BAO GIỜ rơi, 100% LUÔN rơi. Đổi `<` thành `<=` trong
// quyetDinhRoiDo là mục TiLe=0 bắt đầu rơi -> test này đỏ.
func TestQuyetDinhRoiDoHaiDauMut(t *testing.T) {
	bang := []MucRoiDo{
		{MauID: 0, SoLuong: 1, TiLe: 0},
		{MauID: 1, SoLuong: 1, TiLe: roiDoMauSo},
	}
	// Bộ chọn luôn trả 0 (may mắn nhất có thể).
	ra := quyetDinhRoiDo(bang, luonRoi)
	if len(ra) != 1 || ra[0].MauID != 1 {
		t.Fatalf("ti le 0%% khong duoc roi, 100%% phai roi; nhan %+v", ra)
	}
	// Bộ chọn luôn trả mẫu số-1 (xui nhất có thể) — 100% vẫn phải rơi.
	ra = quyetDinhRoiDo(bang, khongRoi)
	if len(ra) != 1 || ra[0].MauID != 1 {
		t.Fatalf("ti le 100%% phai roi ke ca khi gieo xui; nhan %+v", ra)
	}
}

func TestQuyetDinhRoiDoBoQuaSoLuongKhongDuong(t *testing.T) {
	bang := []MucRoiDo{{MauID: 0, SoLuong: 0, TiLe: roiDoMauSo}}
	if ra := quyetDinhRoiDo(bang, luonRoi); len(ra) != 0 {
		t.Fatalf("so luong 0 khong duoc sinh mon nao, nhan %+v", ra)
	}
}

func TestQuyetDinhRoiDoBoChonNilThiKhongRoi(t *testing.T) {
	if ra := quyetDinhRoiDo(bangRoiQuaiMacDinhChoTest(t), nil); ra != nil {
		t.Fatalf("bo chon nil phai tra nil, nhan %+v", ra)
	}
}

// bangRoiQuaiMacDinhChoTest trả bảng rơi MẶC ĐỊNH đọc từ kho nội dung NHÚNG —
// thay cho biến gói `roiDoQuaiMacDinh` đã bị xoá ở KH8 task 9.
func bangRoiQuaiMacDinhChoTest(t *testing.T) []MucRoiDo {
	t.Helper()
	// "ma-quai-khong-co-that" cố tình không tồn tại -> rơi về bảng mặc định.
	return roiDoCuaQuai(khoChieuMacDinh(), "ma-quai-khong-co-that")
}

func TestRoiDoCuaQuaiCoBangMacDinh(t *testing.T) {
	ac := khoChieuMacDinh()
	if len(roiDoCuaQuai(ac, "ma-quai-khong-co-that")) == 0 {
		t.Fatalf("mau quai la -> phai roi ve bang mac dinh, khong duoc rong")
	}
	if len(roiDoCuaQuai(ac, "")) == 0 {
		t.Fatalf("ma rong (quai dung tay) -> cung phai roi ve bang mac dinh")
	}
	// soi-rung CÓ bảng riêng trong bang-roi.yaml, và bảng đó KHÁC mặc định.
	rieng := roiDoCuaQuai(ac, "soi-rung")
	if len(rieng) == 0 {
		t.Fatalf("soi-rung phai co bang rieng")
	}
	macDinh := roiDoCuaQuai(ac, "ma-quai-khong-co-that")
	if len(rieng) == len(macDinh) {
		t.Fatalf("bang rieng cua soi-rung dang giong het bang mac dinh (%d muc) — "+
			"test khong con phan biet duoc hai duong tra cuu", len(rieng))
	}
}

// ⚠️ roiDoCuaBoss phải trả slice MỚI mỗi lần. Nếu ai đó tối ưu thành "trả thẳng
// slice trong ảnh chụp" rồi `append` vào đó thì ảnh chụp đã công bố bị bẩn —
// test này gọi hai lần và kiểm bảng thường không đổi số mục.
func TestRoiDoCuaBossKhongLamBanBangGoc(t *testing.T) {
	ac := khoChieuMacDinh()
	ma := noidung.MaBangRoiBossChung

	soThuong := len(roiDoCuaBoss(ac, ma, &PhaBoss{}))
	if soThuong == 0 {
		t.Fatalf("boss thuong phai co bang roi")
	}

	lon := &PhaBoss{LaBossLon: true}
	a := roiDoCuaBoss(ac, ma, lon)
	b := roiDoCuaBoss(ac, ma, lon)
	if len(a) <= soThuong || len(b) != len(a) {
		t.Fatalf("boss lon phai co THEM muc so voi boss thuong: %d / %d (thuong %d)",
			len(a), len(b), soThuong)
	}
	if n := len(roiDoCuaBoss(ac, ma, &PhaBoss{})); n != soThuong {
		t.Fatalf("bang boss thuong bi ban: %d -> %d", soThuong, n)
	}
	if roiDoCuaBoss(ac, ma, nil) != nil {
		t.Fatalf("pha nil -> khong roi gi")
	}
	if roiDoCuaBoss(ac, "ma-boss-khong-co-that", &PhaBoss{}) != nil {
		t.Fatalf("ma boss la -> khong roi gi (khac quai: boss khong co bang mac dinh)")
	}
}

// ---------------------------------------------------------------------------
// B. Actor
// ---------------------------------------------------------------------------

// khuVucRoiGia dựng khu vực CHƯA BatDau() đã nối bộ rơi đồ (chưa Chay()) và bộ
// chọn ngẫu nhiên ép rơi 100%.
func khuVucRoiGia(t *testing.T, kho KhoRoiDo, chon func(int) int) (*KhuVuc, *BoRoiDo) {
	t.Helper()
	k := MoKhuVuc(TimBanDo(1))
	k.roiDoNgauNhien = chon
	bo := MoBoRoiDo(kho, 32)
	k.GanBoRoi(bo)
	return k, bo
}

// ⚠️ CHỐT CHẶN ĐỘT BIẾN 1: "actor gọi thẳng SinhDoDuoiDat".
//
// Sau khi actor xử lý xong một con quái chết, kho giả phải CHƯA hề bị gọi —
// vì goroutine BoRoiDo chưa Chay(). Nếu ai đó nhét lời gọi CSDL vào thẳng
// roiDoKhiChet (hay vào case suKienDanhQuai), số lần gọi > 0 -> test đỏ.
func TestActorKhongGoiCSDLKhiQuaiChet(t *testing.T) {
	kho := &khoRoiGia{}
	k, bo := khuVucRoiGia(t, kho, luonRoi)

	bang := bangRoiQuaiMacDinhChoTest(t)
	n := k.roiDoKhiChet(100, 200, bang)

	if n != len(bang) {
		t.Fatalf("ep roi 100%% -> phai xep %d yeu cau, xep %d", len(bang), n)
	}
	if len(bo.hang) != n {
		t.Fatalf("hang doi phai co %d yeu cau, co %d", n, len(bo.hang))
	}
	if got := kho.soLan.Load(); got != 0 {
		t.Fatalf("ACTOR KHONG DUOC CHAM CSDL: SinhDoDuoiDat bi goi %d lan", got)
	}
}

func TestQuaiChetTiLeKhongThiKhongXepGi(t *testing.T) {
	kho := &khoRoiGia{}
	k, bo := khuVucRoiGia(t, kho, khongRoi)
	// Bảng toàn tỉ lệ nhỏ (< mẫu số) -> gieo xui thì không món nào rơi.
	if n := k.roiDoKhiChet(0, 0, []MucRoiDo{{MauID: 0, SoLuong: 1, TiLe: 9999}}); n != 0 {
		t.Fatalf("gieo xui phai khong roi gi, xep %d", n)
	}
	if len(bo.hang) != 0 {
		t.Fatalf("hang doi phai rong")
	}
}

func TestChuaNoiBoRoiThiKhongPanicVaKhongRoiGi(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.roiDoNgauNhien = luonRoi
	if n := k.roiDoKhiChet(0, 0, bangRoiQuaiMacDinhChoTest(t)); n != 0 {
		t.Fatalf("chua noi bo roi -> phai tra 0, tra %d", n)
	}
}

// ⚠️ CHỐT CHẶN ĐỘT BIẾN 3: "bossChet gọi rơi đồ hai lần".
//
// Gọi bossChet hai lần liên tiếp trên cùng con boss: lần hai phải bị chốt
// "đã chết rồi thì thôi" chặn -> tổng số yêu cầu đúng bằng MỘT lượt bảng rơi.
// Nếu ai đó thêm một lời gọi roiDoKhiChet thứ hai (ví dụ ở máy trạng thái
// RA_DI) thì số yêu cầu gấp đôi -> test đỏ.
func TestBossChetChiRoiDoMotLan(t *testing.T) {
	kho := &khoRoiGia{}
	k, bo := khuVucRoiGia(t, kho, luonRoi)

	pha := PhaBoss{Ten: "Thu linh", Cap: 1, LaBossLon: false, DiemThuong: 0}
	b := &Boss{Id: 1, PhaHienTai: 0, DuLieu: []PhaBoss{pha},
		Hp: 0, HpToiDa: 100, TrangThai: BossHoatDong}

	bayGio := k.dongHo()
	k.bossChet(b, 0, bayGio)
	sauLanMot := len(bo.hang)
	k.bossChet(b, 0, bayGio) // lần hai: phải bị chặn hoàn toàn

	bangBossThuong := roiDoCuaBoss(k.anhChupNoiDung(), noidung.MaBangRoiBossChung, &pha)
	if sauLanMot != len(bangBossThuong) {
		t.Fatalf("boss thuong phai xep %d yeu cau, xep %d", len(bangBossThuong), sauLanMot)
	}
	if len(bo.hang) != sauLanMot {
		t.Fatalf("bossChet lan hai KHONG duoc roi them: %d -> %d", sauLanMot, len(bo.hang))
	}
	if got := kho.soLan.Load(); got != 0 {
		t.Fatalf("bossChet chay trong actor, khong duoc cham CSDL; bi goi %d lan", got)
	}
}

// ---------------------------------------------------------------------------
// C. Đường về: suKienDoDaSinh -> doDuoiDat + gói 68
// ---------------------------------------------------------------------------

// ⚠️ CHỐT CHẶN ĐỘT BIẾN 2 (một nửa): đây là nhánh THÀNH CÔNG. Cặp với
// TestKhoLoiThiKhongCoMonMa ngay dưới — nếu ai thả món TRƯỚC khi CSDL trả về
// thì test kia đỏ, còn test này vẫn xanh; phải đọc CẢ HAI.
func TestDoDaSinhThiCoMonDuoiDatVaPhatGoi68(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	p := &phienGiaKhuVuc{}
	k.xuLy(suKienVao{nv: &store.NhanVat{ID: 1, Ten: "A", HpHienTai: 100}, phien: p})
	p.nhan = nil

	k.xuLy(suKienDoDaSinh{mon: doVuong{
		vatPhamID: 77, mauID: 4, soLuong: 1, x: 100, y: 384, loai: 2, ten: "Kunai",
	}})

	if len(k.doDuoiDat) != 1 {
		t.Fatalf("phai co dung mot mon duoi dat, co %d", len(k.doDuoiDat))
	}
	var id int16
	for i, m := range k.doDuoiDat {
		id = i
		if m.vatPhamID != 77 {
			t.Fatalf("vatPhamID = %d, mong doi 77", m.vatPhamID)
		}
	}
	tin := p.banTin()
	if len(tin) != 1 || tin[0].Command != 68 {
		t.Fatalf("phai phat dung mot goi 68, nhan %+v", tin)
	}
	// Đọc ngược byte: [short itemMapId][short mauId][short x][short y][int -1]
	than := tin[0].Data
	if len(than) != 12 {
		t.Fatalf("than goi 68 phai 12 byte, co %d", len(than))
	}
	if int16(than[0])<<8|int16(than[1]) != id {
		t.Fatalf("itemMapId trong goi khong khop id trong map")
	}
	if than[2] != 0 || than[3] != 4 {
		t.Fatalf("mauId trong goi = %d %d, mong doi 0 4", than[2], than[3])
	}
	if than[6] != 1 || than[7] != 128 {
		t.Fatalf("y trong goi = %d %d, mong doi 1 128 (=384)", than[6], than[7])
	}
	for i := 8; i < 12; i++ {
		if than[i] != 255 {
			t.Fatalf("chu uu tien phai la -1 (ai cung nhat duoc), byte %d = %d", i, than[i])
		}
	}
}

// ---------------------------------------------------------------------------
// D. Kho lỗi -> KHÔNG có món ma
// ---------------------------------------------------------------------------

// ⚠️ CHỐT CHẶN ĐỘT BIẾN 2: "thả món trước khi CSDL trả về".
//
// Kho giả trả lỗi. BoRoiDo.lam chạy đồng bộ ở đây (gọi thẳng, không qua
// goroutine) -> nếu nó vẫn Gui(suKienDoDaSinh) thì actor sẽ có món; test kiểm
// hộp thư RỖNG và doDuoiDat RỖNG.
// khuVucCoHopThu dựng khu vực CHƯA chạy actor nhưng ĐÃ CÓ hộp thư.
//
// ⚠️ BẪY ĐÃ VẤP: MoKhuVuc KHÔNG tạo `hopThu` — BatDau() mới tạo (khuvuc.go).
// Gọi k.Gui() trên một khu vực chưa BatDau thì gửi vào kênh nil, rơi vào nhánh
// `default` và bị BỎ IM LẶNG. Nghĩa là một test "kiểm hộp thư rỗng" trên khu
// vực chưa BatDau sẽ XANH DÙ CÓ BUG — nó không kiểm gì cả. Hàm này cấp hộp thư
// bằng tay để hai test dưới đây kiểm THẬT, mà vẫn không chạy goroutine actor
// (để đọc doDuoiDat/hopThu trực tiếp không bị đua dữ liệu).
func khuVucCoHopThu() *KhuVuc {
	k := MoKhuVuc(TimBanDo(1))
	k.hopThu = make(chan suKienKhuVuc, 16)
	return k
}

func TestKhoLoiThiKhongCoMonMa(t *testing.T) {
	kho := &khoRoiGia{loi: errors.New("csdl sap")}
	k := khuVucCoHopThu()
	bo := MoBoRoiDo(kho, 4)

	bo.lam(yeuCauRoiDo{khu: k, mauID: 0, soLuong: 1, x: 10, y: 20})

	if got := kho.soLan.Load(); got != 1 {
		t.Fatalf("phai thu ghi CSDL dung mot lan, %d", got)
	}
	if n := len(k.hopThu); n != 0 {
		t.Fatalf("CSDL loi -> KHONG duoc gui su kien nao ve actor, hop thu co %d", n)
	}
	// Cho chắc: bơm hộp thư qua actor rồi kiểm mặt đất vẫn trống.
	for len(k.hopThu) > 0 {
		k.xuLy(<-k.hopThu)
	}
	if len(k.doDuoiDat) != 0 {
		t.Fatalf("KHONG duoc co mon ma: doDuoiDat co %d mon", len(k.doDuoiDat))
	}
}

// Kho chạy ngon: lam() phải gửi đúng MỘT suKienDoDaSinh mang id CSDL thật.
func TestKhoNgonThiGuiDungMotSuKienMangIdThat(t *testing.T) {
	kho := &khoRoiGia{}
	k := khuVucCoHopThu()
	bo := MoBoRoiDo(kho, 4)

	bo.lam(yeuCauRoiDo{khu: k, mauID: 4, soLuong: 2, x: 10, y: 20})

	if n := len(k.hopThu); n != 1 {
		t.Fatalf("phai gui dung mot su kien, gui %d", n)
	}
	sk, ok := (<-k.hopThu).(suKienDoDaSinh)
	if !ok {
		t.Fatalf("su kien sai kieu")
	}
	if sk.mon.vatPhamID != 1 {
		t.Fatalf("id vat pham = %d, mong doi id CSDL that (1)", sk.mon.vatPhamID)
	}
	if sk.mon.soLuong != 2 || sk.mon.mauID != 4 {
		t.Fatalf("mon sai: %+v", sk.mon)
	}
	// loai/ten phải được tra từ danh mục mẫu, không để rỗng.
	mau := TimMauVatPham(4)
	if sk.mon.ten != mau.Ten || sk.mon.loai != mau.Loai {
		t.Fatalf("loai/ten phai tra tu danh muc mau: %q/%d", sk.mon.ten, sk.mon.loai)
	}
}

func TestBoRoiDoXepKhiHangDayTraFalse(t *testing.T) {
	bo := MoBoRoiDo(&khoRoiGia{}, 1)
	if !bo.Xep(yeuCauRoiDo{}) {
		t.Fatalf("yeu cau dau phai vao duoc hang")
	}
	if bo.Xep(yeuCauRoiDo{}) {
		t.Fatalf("hang day phai tra false, KHONG duoc chan")
	}
	var nil0 *BoRoiDo
	if nil0.Xep(yeuCauRoiDo{}) {
		t.Fatalf("bo nil phai tra false, khong panic")
	}
}
