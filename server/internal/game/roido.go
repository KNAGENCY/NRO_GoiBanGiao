package game

// roido.go — KH7 Task 12: quái/boss chết thì RƠI ĐỒ.
//
// Giải mã hai phía (PHA A): docs/superpowers/research/2026-09-03-kh7-roi-do.md
// Đặc tả: docs/superpowers/specs/2026-09-03-kh7-vat-pham-design.md hàng 12, §4.
// Kỷ luật actor: BAT-DAU-TU-DAY.md §6.9/§6.10b.
//
// ---------------------------------------------------------------------------
// AI CHẠM CSDL — đọc dòng này trước khi sửa bất cứ thứ gì trong file
// ---------------------------------------------------------------------------
//
//	goroutine ACTOR khu vực (vongLap)        goroutine BO ROI DO
//	---------------------------------        -------------------
//	quái Hp<=0  /  bossChet
//	  quyetDinhRoiDo(bảng rơi, k.roiDoNgauNhien)   <- THUẦN, không I/O
//	  phát -12 với soDoRoi = 0
//	  boRoi.Xep(yeuCauRoiDo{...})  ------------->
//	                                            store.SinhDoDuoiDat
//	                                              (INSERT vat_pham + lịch sử
//	                                               SINH_RA, MỘT giao dịch — BB2)
//	  <---- k.Gui(suKienDoDaSinh{...}) <---------
//	  thaDoXuongKhuVuc: cấp itemMapId + phát 68
//
// LỖI CSDL -> KHÔNG gửi suKienDoDaSinh -> actor không thả gì -> KHÔNG CÓ MÓN MA.
// Thứ tự này là bắt buộc ("ghi sổ trước — giao sau", §2.3): thả trước rồi mới
// ghi là tạo ra một món người chơi nhìn thấy nhưng không tồn tại trong CSDL,
// và đường nhặt của task 11 sẽ trả về ErrDaCoNguoiNhat một cách khó hiểu.
//
// ---------------------------------------------------------------------------
// VÌ SAO GÓI `-12` VẪN GHI soDoRoi = 0 (quyết định của task này)
// ---------------------------------------------------------------------------
//
// Client đọc `[byte soDoRoi]` rồi từng món NGAY trong gói -12 (Controller.cs
// case -12, ~dòng 4418). Muốn điền được, actor phải đứng CHỜ CSDL trả id —
// đúng cái §6.10b cấm: một câu INSERT chậm 50ms làm đứng hình cả bản đồ với
// mọi người đang ở trong đó.
//
// Thêm hai lý do độc lập, ghi đầy đủ ở PHA A §2:
//   - Client `break` NGAY (trước cả khi đọc dame) nếu nó không tìm thấy con
//     quái trong danh sách của nó -> cả danh sách đồ bị vứt, im lặng.
//   - Toàn bộ phần sau nằm trong `catch (Exception) { }` RỖNG.
//
// Gói 68 không có nhược điểm nào trong số đó: nó độc lập, không tham chiếu con
// quái nào, và task 11 đã xác minh + có test đọc ngược byte. Khác biệt người
// chơi thấy: món hiện tại chỗ thay vì bay ra từ xác quái. Thuần thẩm mỹ.
//
// Vì vậy nhánh quái chết trong khuvuc.go GIỮ NGUYÊN `goiQuaiChet(..., nil)`.
// Đừng "sửa cho đủ" chỗ đó — nó là chủ đích, không phải việc còn dở.

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"time"

	"nroserver/internal/noidung"
	"nroserver/internal/store"
)

// ---------------------------------------------------------------------------
// 1. Bảng rơi
// ---------------------------------------------------------------------------

// roiDoMauSo là mẫu số của mọi tỉ lệ trong file này: TiLe = 1500 nghĩa là 15%.
// Dùng số nguyên phần vạn thay cho float để tỉ lệ so sánh được chính xác tuyệt
// đối và test tiêm được bộ chọn giả (0 -> luôn rơi, mẫu số-1 -> không bao giờ).
const roiDoMauSo = 10000

// MucRoiDo là MỘT dòng của bảng rơi: "món mauID, soLuong cái, tỉ lệ TiLe/10000".
// Mỗi dòng gieo ĐỘC LẬP — một lần quái chết có thể rơi 0, 1 hay tất cả các dòng.
type MucRoiDo struct {
	MauID   int32
	SoLuong int32
	TiLe    int
}

// ⚠️ BỐN BIẾN GÓI `roiDoTheoQuai` / `roiDoQuaiMacDinh` / `roiDoTheoBoss` /
// `roiDoTheoBossLon` ĐÃ BỊ XOÁ Ở KH8 TASK 9.
//
// Số của chúng nay nằm trong `kho-noi-dung/bang-roi.yaml`, tra cứu bằng MÃ
// CHUỖI của mẫu quái (`Quai.Ma`) chứ không phải bằng `Quai.MauId` — từ task 8
// trường đó mang `anh_id`, nên tra theo nó là hai mẫu quái dùng chung một ảnh
// sẽ dùng chung một bảng rơi, im lặng. Chi tiết ở đầu `internal/noidung/bangroi.go`.
//
// Đột biến "actor đọc bảng Go cũ" KHÔNG THỂ viết được nữa — không còn bảng nào.

// mucRoiTuNoiDung đổi các dòng bảng rơi của kho nội dung (tỉ lệ PHẦN TRĂM, cho
// người đọc file) sang MucRoiDo của gói này (tỉ lệ PHẦN VẠN, cho phép so sánh
// số nguyên chính xác tuyệt đối).
//
// `laBossLon` lọc: false (quái, và boss thường) = BỎ mọi dòng `chi_boss_lon`;
// true = lấy hết (boss lớn hưởng cả bảng thường lẫn phần thêm).
//
// Hàm THUẦN. Luôn dựng slice MỚI — không bao giờ trả về slice nằm trong ảnh
// chụp, kể cả khi lấy hết: người gọi `append` vào đó là sửa ảnh chụp đã công bố
// cho mọi người chơi (bẫy đã dính hai lần trong dự án này, xem đầu kho.go).
func mucRoiTuNoiDung(ds []noidung.MucRoi, laBossLon bool) []MucRoiDo {
	ra := make([]MucRoiDo, 0, len(ds))
	for _, m := range ds {
		if m.ChiBossLon && !laBossLon {
			continue
		}
		tiLe := int(m.TiLePhanTram / 100 * roiDoMauSo)
		if m.ChacChan {
			tiLe = roiDoMauSo
		}
		ra = append(ra, MucRoiDo{MauID: int32(m.MauId), SoLuong: int32(m.SoLuong), TiLe: tiLe})
	}
	if len(ra) == 0 {
		return nil
	}
	return ra
}

// roiDoCuaQuai tra bảng rơi cho MỘT mẫu quái, đọc TỪ MỘT ẢNH CHỤP đã cầm sẵn.
//
// ⚠️ NHẬN ẢNH CHỤP LÀM THAM SỐ, KHÔNG TỰ GỌI AnhChup(). Trong actor phải truyền
// `k.anhChupNoiDung()` (ảnh chụp của nhịp/sự kiện hiện tại). Gọi `AnhChup()`
// hay `layAnhChupNhip()` thêm một lần ở đây là phá hợp đồng "một lần đầu nhịp"
// (internal/noidung/kho.go) — có test giết được: TestKhoBangRoiMotLanKhiQuaiChet.
//
// Mã lạ (hoặc rỗng — quái dựng tay trong test) -> bảng MẶC ĐỊNH, không rỗng.
func roiDoCuaQuai(ac khoChieuAnhChup, ma string) []MucRoiDo {
	return mucRoiTuNoiDung(ac.BangRoiCuaQuai(ma), false)
}

// roiDoCuaBoss tra bảng rơi cho MỘT PHA của MỘT con boss, đọc TỪ MỘT ẢNH CHỤP.
//
// pha == nil (dữ liệu pha hỏng — bossChet chịu được nil ở mọi bước) -> không
// rơi gì. Thà không rơi còn hơn rơi theo một pha không biết là pha nào.
//
// `ma` hôm nay luôn là noidung.MaBangRoiBossChung tại chỗ gọi trong khuvuc.go:
// `game.Boss` chưa mang mã con boss (lichboss.go thả trường đó khi dựng lịch).
// Chữ ký vẫn nhận `ma` để ngày Boss mang được Ma thì chỉ phải đổi ĐÚNG một dòng
// ở chỗ gọi. Lý do đầy đủ ở đầu internal/noidung/bangroi.go.
func roiDoCuaBoss(ac khoChieuAnhChup, ma string, pha *PhaBoss) []MucRoiDo {
	if pha == nil {
		return nil
	}
	return mucRoiTuNoiDung(ac.BangRoiCuaBoss(ma), pha.LaBossLon)
}

// MaBossHopLeChoBangRoi là danh sách mã của MỌI con boss có trong kho nội dung
// boss của gói này — tham số cho `noidung.KiemMaBangRoiBoss`.
//
// VÌ SAO CHỈ CÓ HÀM NÀY MÀ CHƯA CÓ CHỖ GỌI TRONG ACTOR: bộ kiểm mã boss là hàm
// THUẦN nhận `*noidung.DuLieu`, mà đường khởi động của gói này chỉ cầm
// `*noidung.Kho` (ảnh chụp), không cầm `DuLieu`. Đường duy nhất có `DuLieu`
// trong tay là bộ nạp nóng (KH8 task 10) và `main.go` — cả hai đều là file
// người khác đang sở hữu ở nhịp này. Hàm để sẵn đây, có test riêng, chờ nối.
//
// Hậu quả của việc chưa nối: một mục `boss:` trỏ tới mã không có thật sẽ KHÔNG
// bị bắt lúc nạp. Hậu quả thật của lỗi đó là "bảng đó không bao giờ được dùng
// tới" — boss không rơi gì, chứ không phải máy chủ sai. Chấp nhận được cho tới
// khi task 10 nối vào.
func MaBossHopLeChoBangRoi() []string {
	ds := DanhSachNoiDungBoss()
	ma := make([]string, 0, len(ds))
	for _, nd := range ds {
		ma = append(ma, nd.Ma)
	}
	return ma
}

// quyetDinhRoiDo gieo từng dòng của bảng và tra danh sách món THỰC SỰ rơi.
//
// Hàm THUẦN (không I/O, không time.Now) + bộ chọn ngẫu nhiên TIÊM ĐƯỢC — cùng
// khuôn với `chon` của lichboss.go. `chon(n)` phải tra một số trong [0, n).
// Test tiêm `func(int) int { return 0 }` để ép rơi 100%, và
// `func(int) int { return n-1 }` để ép không rơi gì.
//
// So sánh là `chon(roiDoMauSo) < TiLe`, nên TiLe = 0 KHÔNG BAO GIỜ rơi (0 < 0
// sai) và TiLe = roiDoMauSo LUÔN rơi. Hai đầu mút đều đúng, có test riêng.
func quyetDinhRoiDo(bang []MucRoiDo, chon func(int) int) []MucRoiDo {
	if chon == nil {
		return nil
	}
	var ra []MucRoiDo
	for _, m := range bang {
		if m.TiLe <= 0 || m.SoLuong <= 0 {
			continue
		}
		if chon(roiDoMauSo) < m.TiLe {
			ra = append(ra, m)
		}
	}
	return ra
}

// ---------------------------------------------------------------------------
// 2. Bộ rơi đồ — goroutine DUY NHẤT được chạm CSDL trong luồng này
// ---------------------------------------------------------------------------

// KhoRoiDo là phần tối thiểu của *store.Store mà bộ rơi đồ cần — interface HẸP,
// cùng khuôn với KhoNhat (dovuong.go) và Kho (login.go), để test được bằng kho
// giả đồng bộ, không cần CSDL thật.
type KhoRoiDo interface {
	SinhDoDuoiDat(ctx context.Context, mauID int32, soLuong int32, chiSoPhu store.ChiSoPhu) (int64, error)
}

// yeuCauRoiDo là MỘT yêu cầu ghi CSDL do actor xếp ra. Nó KHÔNG mang id vật
// phẩm (chưa có) và KHÔNG mang itemMapId (chỉ actor mới cấp được).
type yeuCauRoiDo struct {
	khu      *KhuVuc
	mauID    int32
	soLuong  int32
	x, y     int16
	chiSoPhu store.ChiSoPhu
}

// suKienDoDaSinh là đường VỀ từ goroutine CSDL: "món này đã có dòng thật trong
// CSDL rồi, thả xuống đất đi". Chỉ được gửi khi giao dịch ĐÃ COMMIT.
//
// Tách riêng với suKienThaDo (task 11) dù thân gần giống: suKienThaDo có kênh
// traLoi cho người gọi đồng bộ trong test, còn đường rơi đồ là một chiều và
// không ai chờ. Giữ hai kiểu để test phân biệt được "đồ rơi tự sinh" với "ai
// đó gọi ThaDoVaoKhuVuc bằng tay".
type suKienDoDaSinh struct {
	mon doVuong
}

// BoRoiDo là goroutine ghi CSDL cho đường rơi đồ. Cùng khuôn với BoNhatDo:
// hàng đợi có sức chứa, Xep KHÔNG BAO GIỜ chặn, Dung() gọi đúng một lần.
type BoRoiDo struct {
	kho  KhoRoiDo
	hang chan yeuCauRoiDo
	dung chan struct{}
	// hetGio là hạn cho MỘT yêu cầu CSDL. Vượt hạn -> bỏ món đó (không có món
	// ma) còn hơn để yêu cầu treo mãi và nghẽn cả hàng.
	hetGio time.Duration
}

// MoBoRoiDo tạo bộ rơi đồ với hàng đợi sức chứa `sucChua`. Chưa chạy gì — gọi
// Chay() (thường là `go bo.Chay()`) để bắt đầu.
func MoBoRoiDo(kho KhoRoiDo, sucChua int) *BoRoiDo {
	if sucChua <= 0 {
		sucChua = 256
	}
	return &BoRoiDo{
		kho:    kho,
		hang:   make(chan yeuCauRoiDo, sucChua),
		dung:   make(chan struct{}),
		hetGio: 5 * time.Second,
	}
}

// Xep đưa một yêu cầu vào hàng. KHÔNG BAO GIỜ CHẶN — hàng đầy thì tra false và
// món đó không rơi. Actor chờ ở đây một phần nghìn giây cũng là treo cả bản đồ
// (§6.10b). Mất một món đồ rơi chấp nhận được; treo một actor thì không.
func (b *BoRoiDo) Xep(yc yeuCauRoiDo) bool {
	if b == nil {
		return false
	}
	select {
	case b.hang <- yc:
		return true
	default:
		return false
	}
}

// Chay là vòng lặp của goroutine ghi CSDL. Trả về khi Dung() được gọi.
func (b *BoRoiDo) Chay() {
	for {
		select {
		case <-b.dung:
			return
		case yc := <-b.hang:
			b.lam(yc)
		}
	}
}

// Dung dừng goroutine. Gọi ĐÚNG MỘT LẦN (chủ sở hữu là main) — cùng quy ước
// với BoNhatDo.Dung.
func (b *BoRoiDo) Dung() { close(b.dung) }

// lam chạy MỘT yêu cầu: ghi CSDL -> báo ngược về actor.
//
// LỖI -> LOG RỒI THÔI, KHÔNG gửi gì về actor. Đây là điểm khác căn bản với
// BoNhatDo.lam (nó LUÔN báo ngược, kể cả lỗi, vì actor phải trả món lại mặt
// đất). Ở đây chưa có gì trên mặt đất để trả: không có dòng CSDL thì không có
// món, chấm hết. Gửi suKienDoDaSinh khi lỗi chính là cách tạo ra MÓN MA.
func (b *BoRoiDo) lam(yc yeuCauRoiDo) {
	ctx, huy := context.WithTimeout(context.Background(), b.hetGio)
	defer huy()

	id, err := b.kho.SinhDoDuoiDat(ctx, yc.mauID, yc.soLuong, yc.chiSoPhu)
	if err != nil {
		slog.Error("sinh do roi duoi dat that bai, bo mon nay",
			"mauID", yc.mauID, "soLuong", yc.soLuong, "loi", err)
		return
	}

	var loai int8
	var ten string
	if mau := TimMauVatPham(int16(yc.mauID)); mau != nil {
		loai, ten = mau.Loai, mau.Ten
	}
	yc.khu.Gui(suKienDoDaSinh{mon: doVuong{
		vatPhamID: id, mauID: yc.mauID, soLuong: yc.soLuong,
		x: yc.x, y: yc.y, loai: loai, ten: ten,
	}})
}

// ---------------------------------------------------------------------------
// 3. Đường vào từ actor
// ---------------------------------------------------------------------------

// roiDoKhiChet là TOÀN BỘ phần rơi đồ mà actor được làm. CHẠY TRONG ACTOR.
//
// Ba việc, tất cả đều chỉ chạm bộ nhớ: gieo tỉ lệ, xếp yêu cầu, đếm. KHÔNG
// một lời gọi CSDL nào — có test riêng dùng kho giả đếm để chốt điều đó.
//
// Trả về số yêu cầu ĐÃ XẾP ĐƯỢC (dùng cho test và log; nhánh gọi bỏ qua).
func (k *KhuVuc) roiDoKhiChet(x, y int16, bang []MucRoiDo) int {
	bo := k.boRoi.Load()
	if bo == nil {
		return 0 // chưa nối bộ rơi đồ -> im lặng không rơi gì, KHÔNG panic
	}
	chon := k.roiDoNgauNhien
	if chon == nil {
		chon = rand.IntN
	}
	n := 0
	for _, m := range quyetDinhRoiDo(bang, chon) {
		if bo.Xep(yeuCauRoiDo{khu: k, mauID: m.MauID, soLuong: m.SoLuong, x: x, y: y}) {
			n++
		}
	}
	return n
}

// GanBoRoi nối bộ ghi CSDL cho đường rơi đồ vào khu vực này. Gọi TỪ goroutine
// main (sau MoKhuVuc, có thể sau cả BatDau) — an toàn nhờ atomic.Pointer, cùng
// lý do với GanBoNhat/GanLoa.
func (k *KhuVuc) GanBoRoi(b *BoRoiDo) { k.boRoi.Store(b) }
