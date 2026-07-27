package game

// Kiem thu Finding 1 cua soi toan nhanh KH11 (fix-cuoi-report.md, CRITICAL):
// xuLyDangNhap (login.go) KHONG chan dang nhap lan hai TREN CUNG MOT KET NOI.
// Moi lan -29 THANH CONG chi don gian gan lai b.taiKhoan roi DangKy tai khoan
// MOI — neu ket noi gui -29 THANH CONG hai lan (X roi Y), CHINH mot BoXuLy
// nam o HAI khoa trong SoDangNhap.phien cung luc (X->b VA Y->b). Khoa CU (X)
// khong bao gio duoc XoaNeuDung don dep dung cach (KhiDong chi biet xoa khoa
// b.taiKhoan HIEN TAI, tuc Y) — no tro thanh "rac" vinh vien, va lan da TIEP
// THEO nham dung X se goi close() len CHINH b.biDa DA DONG TU TRUOC (o lan da
// cua Y) -> double-close mot channel -> PANIC khong the recover (khong co
// recover() nao tren duong PhucVu, net/session.go) -> CHET CA TIEN TRINH.
//
// Ca hai test o day dung mang that (net.Pipe + mynet.PhucVu that, dung y het
// cach client that noi) + SoDangNhap that, nhung KHONG can CSDL that: kho gia
// khoHaiTaiKhoan duoi day tra ve ErrKhongThayNhanVat cho moi tai khoan, nen
// dang nhap thanh cong luon di theo nhanh "chua co nhan vat" (gui lenh 2,
// khong cham toi vaoTheGioi/KhoLuu that).

import (
	"context"
	"net"
	"testing"
	"time"

	"nroserver/internal/auth"
	mynet "nroserver/internal/net"
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// khoHaiTaiKhoan la mot Kho gia CHO PHEP tra ve HAI TAI KHOAN KHAC NHAU tuy
// vao ten dang nhap duoc truyen — khac han khoGia (login_test.go), von LUON
// tra ve DUY NHAT truong tk bat ke ten truyen vao. Finding 1 CAN hai tai
// khoan THAT SU khac ID dang nhap tren CUNG mot ket noi de phoi bay dung loi
// (mot khoGia don le se lam ca hai lan dang nhap cung tro ve MOT ID, khong du
// de bay "mot BoXuLy nam o hai khoa trong so").
type khoHaiTaiKhoan struct {
	theoTen map[string]*store.TaiKhoan
}

func (k *khoHaiTaiKhoan) LayTaiKhoan(ctx context.Context, ten string) (*store.TaiKhoan, error) {
	tk, ok := k.theoTen[ten]
	if !ok {
		return nil, store.ErrKhongThayTaiKhoan
	}
	return tk, nil
}

func (k *khoHaiTaiKhoan) LayNhanVatTheoTaiKhoan(ctx context.Context, taiKhoanID int64) (*store.NhanVat, error) {
	return nil, store.ErrKhongThayNhanVat // moi tai khoan trong file nay deu "chua co nhan vat"
}

func (k *khoHaiTaiKhoan) TenNhanVatDaDung(ctx context.Context, ten string) (bool, error) {
	return false, nil
}

func (k *khoHaiTaiKhoan) TaoNhanVat(ctx context.Context, nv store.NhanVatMoi) (*store.NhanVat, error) {
	return nil, nil // khong dung toi trong cac test cua file nay
}

// TestXuLyDangNhap_LanHaiTrenCungKetNoiKhongDangKyLai la kiem chung TRUC TIEP
// cho chot cua Finding 1: hai goi -29 THANH CONG lien tiep tren CUNG mot
// BoXuLy (mo phong dung mot ket noi vat ly) — lan HAI PHAI bi bo qua HOAN
// TOAN (khong goi tin tra loi, khong dang ky vao SoDangNhap, khong ghi de
// b.taiKhoan) thay vi dang ky de tao mot BoXuLy nam o hai khoa trong so.
func TestXuLyDangNhap_LanHaiTrenCungKetNoiKhongDangKyLai(t *testing.T) {
	bamX, err := auth.Hash("matkhauX")
	if err != nil {
		t.Fatalf("bam mat khau X that bai: %v", err)
	}
	bamY, err := auth.Hash("matkhauY")
	if err != nil {
		t.Fatalf("bam mat khau Y that bai: %v", err)
	}
	kho := &khoHaiTaiKhoan{theoTen: map[string]*store.TaiKhoan{
		"taikhoanX": {ID: 501, TenDangNhap: "taikhoanX", MatKhauBam: bamX},
		"taikhoanY": {ID: 502, TenDangNhap: "taikhoanY", MatKhauBam: bamY},
	}}
	soDangNhap := MoSoDangNhap()
	b := NewBoXuLy(kho, nil)
	b.DatSoDangNhap(soDangNhap)

	benClient, benServer := net.Pipe()
	go mynet.PhucVu(benServer, func() mynet.XuLyGoiTin { return b })
	defer benClient.Close()
	benClient.SetReadDeadline(time.Now().Add(5 * time.Second))

	f := protocol.NewFramer(benClient, benClient)
	if err := f.WriteMessage(&protocol.Message{Command: -27}); err != nil {
		t.Fatalf("gui goi bat tay that bai: %v", err)
	}
	khoaTho := docGoiKhoaThuCong(t, benClient)
	f.SetKey(protocol.DeriveKey(khoaTho))

	// Dang nhap LAN DAU: tai khoan X -> thanh cong, dang ky vao so, nhan lenh
	// 2 (chua co nhan vat).
	if err := f.WriteMessage(dungGoiDangNhap("taikhoanX", "matkhauX")); err != nil {
		t.Fatalf("gui dang nhap X that bai: %v", err)
	}
	m1 := docMotGoi(t, f)
	if m1.Command != lenhTaoNhanVat {
		t.Fatalf("lan dau dang nhap X: lenh = %d, mong doi %d (chuyen sang man tao nhan vat)", m1.Command, lenhTaoNhanVat)
	}

	soDangNhap.mu.Lock()
	soLuongSauLanDau := len(soDangNhap.phien)
	phienChoX := soDangNhap.phien[501]
	soDangNhap.mu.Unlock()
	if soLuongSauLanDau != 1 || phienChoX != b {
		t.Fatalf("sau lan dang nhap DAU: so co %d muc (mong doi 1), phien[501]=%p (mong doi CHINH b=%p)",
			soLuongSauLanDau, phienChoX, b)
	}

	// Dang nhap LAN HAI TREN CUNG KET NOI: tai khoan Y — day CHINH LA Finding
	// 1. Voi chot da sua (b.taiKhoan != nil), xuLyDangNhap phai tra ve nil
	// NGAY, KHONG duoc gui goi tin nao, KHONG duoc goi DangKy them.
	if err := f.WriteMessage(dungGoiDangNhap("taikhoanY", "matkhauY")); err != nil {
		t.Fatalf("gui dang nhap Y (lan hai, cung ket noi) that bai: %v", err)
	}

	benClient.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	if _, err := f.ReadMessage(); err == nil {
		t.Fatal("dang nhap lan hai tren CUNG ket noi van duoc phuc vu (nhan duoc goi tra loi) — " +
			"thieu chot chan dang nhap chong (Finding 1): xuLyDangNhap phai tra ve nil ngay khi b.taiKhoan != nil")
	}
	benClient.SetReadDeadline(time.Now().Add(5 * time.Second))

	soDangNhap.mu.Lock()
	soLuongSauLanHai := len(soDangNhap.phien)
	_, coY := soDangNhap.phien[502]
	soDangNhap.mu.Unlock()

	if coY {
		t.Fatal("dang nhap lan hai (tai khoan Y) tren CUNG ket noi da duoc dang ky vao so — " +
			"day chinh la Finding 1: mot BoXuLy nam o HAI khoa trong so, tao muc rac vinh vien " +
			"(va mo duong toi double-close panic khi ai do sau nay da tai khoan CU)")
	}
	if soLuongSauLanHai != 1 {
		t.Fatalf("so co %d muc sau lan dang nhap thu hai, mong doi VAN LA 1 (khong duoc them muc nao)", soLuongSauLanHai)
	}
	if b.taiKhoan == nil || b.taiKhoan.ID != 501 {
		t.Fatalf("b.taiKhoan sau lan dang nhap thu hai = %+v, mong doi VAN LA tai khoan X (ID=501) — khong duoc ghi de bang tai khoan Y", b.taiKhoan)
	}
}

// moKetNoiChoTest dung mot phien XuLy that (net.Pipe + mynet.PhucVu that)
// tren mot *BoXuLy MOI, da DatSoDangNhap voi so dung chung — dung cho kich
// ban 4 buoc duoi day, noi can NHIEU ket noi that su rieng biet cung chia se
// mot SoDangNhap.
func moKetNoiChoTest(t *testing.T, kho Kho, soDangNhap *SoDangNhap) (b *BoXuLy, f *protocol.Framer, conn net.Conn) {
	t.Helper()
	b = NewBoXuLy(kho, nil)
	b.DatSoDangNhap(soDangNhap)

	benClient, benServer := net.Pipe()
	go mynet.PhucVu(benServer, func() mynet.XuLyGoiTin { return b })
	benClient.SetReadDeadline(time.Now().Add(5 * time.Second))

	f = protocol.NewFramer(benClient, benClient)
	if err := f.WriteMessage(&protocol.Message{Command: -27}); err != nil {
		t.Fatalf("bat tay that bai: %v", err)
	}
	khoaTho := docGoiKhoaThuCong(t, benClient)
	f.SetKey(protocol.DeriveKey(khoaTho))

	return b, f, benClient
}

// TestDangNhapChong_KichBan4Buoc_KhongTaoRacKhongPanic dung LAI dung 4 buoc
// mo ta trong Finding 1 ("kich ban giet may chu"):
//
//  1. Ket noi 1 dang nhap X -> so: X->b1
//  2. CUNG ket noi 1 gui tiep -29 dang nhap Y -> VOI chot da sua, buoc nay
//     PHAI bi chan hoan toan (so KHONG duoc them muc Y->b1) — day la diem
//     ma chot chinh phai bat duoc TRUOC KHI cac buoc sau co co hoi xay ra.
//  3. Ket noi 2 dang nhap Y -> vi Y CHUA TUNG dang ky truoc do (buoc 2 da bi
//     chan), day la dang ky LAN DAU cho Y, khong da ai ca — ket noi 1 (X)
//     KHONG duoc dung cham toi.
//  4. Ket noi 3 dang nhap X -> da PHIEN 1 (b1, van dang giu X tu buoc 1,
//     chua he bi dung cham) — day la lan da HOP LE DAU TIEN cho tai khoan
//     X, PHAI hoan tat khong panic (khac han kich ban KHONG co fix, noi day
//     se la lan close() THU HAI len CHINH b1.biDa).
func TestDangNhapChong_KichBan4Buoc_KhongTaoRacKhongPanic(t *testing.T) {
	bamX, err := auth.Hash("matkhauX")
	if err != nil {
		t.Fatalf("bam mat khau X that bai: %v", err)
	}
	bamY, err := auth.Hash("matkhauY")
	if err != nil {
		t.Fatalf("bam mat khau Y that bai: %v", err)
	}
	kho := &khoHaiTaiKhoan{theoTen: map[string]*store.TaiKhoan{
		"kb4X": {ID: 601, TenDangNhap: "kb4X", MatKhauBam: bamX},
		"kb4Y": {ID: 602, TenDangNhap: "kb4Y", MatKhauBam: bamY},
	}}
	soDangNhap := MoSoDangNhap()

	// Buoc 1: ket noi 1 dang nhap X.
	b1, f1, conn1 := moKetNoiChoTest(t, kho, soDangNhap)
	defer conn1.Close()
	if err := f1.WriteMessage(dungGoiDangNhap("kb4X", "matkhauX")); err != nil {
		t.Fatalf("buoc 1: gui dang nhap X that bai: %v", err)
	}
	m1 := docMotGoi(t, f1)
	if m1.Command != lenhTaoNhanVat {
		t.Fatalf("buoc 1: lenh = %d, mong doi %d", m1.Command, lenhTaoNhanVat)
	}

	// Buoc 2: CUNG ket noi 1, gui tiep -29 dang nhap Y.
	if err := f1.WriteMessage(dungGoiDangNhap("kb4Y", "matkhauY")); err != nil {
		t.Fatalf("buoc 2: gui dang nhap Y (cung ket noi 1) that bai: %v", err)
	}
	conn1.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	if _, err := f1.ReadMessage(); err == nil {
		t.Fatal("buoc 2: dang nhap Y tren CUNG ket noi voi X van duoc phuc vu — thieu chot chan dang nhap chong (Finding 1)")
	}
	conn1.SetReadDeadline(time.Now().Add(5 * time.Second))

	soDangNhap.mu.Lock()
	_, coY := soDangNhap.phien[602]
	soDangNhap.mu.Unlock()
	if coY {
		t.Fatal("buoc 2: tai khoan Y da bi dang ky (qua ket noi 1) — dung Finding 1 mo ta: mot BoXuLy nam o HAI khoa trong so")
	}

	// Ket noi 1 (X) khong duoc dung cham toi boi buoc 2 (bi chan hoan toan).
	select {
	case <-b1.biDa:
		t.Fatal("buoc 2: ket noi 1 (X) bi da mot cach vo ly boi chinh buoc dang nhap Y bi chan tren no")
	default:
	}

	// Buoc 3: ket noi 2 dang nhap Y — day la dang ky LAN DAU THAT SU cho Y
	// (buoc 2 khong dang ky duoc gi ca), nen KHONG da ai.
	b2, f2, conn2 := moKetNoiChoTest(t, kho, soDangNhap)
	defer conn2.Close()
	if err := f2.WriteMessage(dungGoiDangNhap("kb4Y", "matkhauY")); err != nil {
		t.Fatalf("buoc 3: gui dang nhap Y (ket noi 2) that bai: %v", err)
	}
	m3 := docMotGoi(t, f2)
	if m3.Command != lenhTaoNhanVat {
		t.Fatalf("buoc 3: lenh = %d, mong doi %d", m3.Command, lenhTaoNhanVat)
	}
	_ = b2

	select {
	case <-b1.biDa:
		t.Fatal("buoc 3: ket noi 1 (X) bi da mot cach vo ly — Y chua tung dang ky truoc do (buoc 2 bi chan) nen khong co gi de da")
	default:
	}

	// Buoc 4: ket noi 3 dang nhap X — X van dang duoc so tro DUNG boi b1 (tu
	// buoc 1, chua he bi dung cham qua cac buoc 2/3), nen day la LAN DA HOP
	// LE DAU TIEN cho tai khoan X. Diem quan trong nhat cua ca kich ban: neu
	// code co loi (thieu chot, nhu truoc khi sua), ket noi 1 se DA BI DA tu
	// buoc 3 (khi do so[Y] dang tro sai ve b1 do buoc 2 dang ky nham) — buoc
	// nay se la lan close() THU HAI len CHINH b1.biDa -> panic double-close
	// -> ca tien trinh test SAP (bang chung ro rang nhat cho Finding 1).
	b3, f3, conn3 := moKetNoiChoTest(t, kho, soDangNhap)
	defer conn3.Close()
	if err := f3.WriteMessage(dungGoiDangNhap("kb4X", "matkhauX")); err != nil {
		t.Fatalf("buoc 4: gui dang nhap X (ket noi 3) that bai: %v", err)
	}
	m4 := docMotGoi(t, f3)
	if m4.Command != lenhTaoNhanVat {
		t.Fatalf("buoc 4: lenh = %d, mong doi %d", m4.Command, lenhTaoNhanVat)
	}
	_ = b3

	// Ket noi 1 (X, phien CU) PHAI bi da boi dang nhap X lan hai tren ket
	// noi 3 — day la lan da HOP LE DAU TIEN (khong phai double-close).
	select {
	case <-b1.biDa:
	case <-time.After(2 * time.Second):
		t.Fatal("buoc 4: ket noi 1 (X, phien CU) phai bi da boi dang nhap X lan hai tren ket noi 3")
	}

	// Kiem tra sau cung: so PHAI o trang thai gon gang — dung hai muc (X->b3,
	// Y->b2), khong con "rac" nao sot lai tu buoc 2.
	soDangNhap.mu.Lock()
	soLuongCuoiCung := len(soDangNhap.phien)
	phienCuoiChoX := soDangNhap.phien[601]
	phienCuoiChoY := soDangNhap.phien[602]
	soDangNhap.mu.Unlock()
	if soLuongCuoiCung != 2 {
		t.Fatalf("so cuoi cung co %d muc, mong doi dung 2 (X va Y, khong con rac)", soLuongCuoiCung)
	}
	if phienCuoiChoX != b3 {
		t.Fatalf("so[X] cuoi cung = %p, mong doi b3 (%p, phien dang nhap X lan hai)", phienCuoiChoX, b3)
	}
	if phienCuoiChoY != b2 {
		t.Fatalf("so[Y] cuoi cung = %p, mong doi b2 (%p)", phienCuoiChoY, b2)
	}
}
