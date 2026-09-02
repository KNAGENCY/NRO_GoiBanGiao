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

// roiDoTheoQuai: bảng rơi theo MauId của quái (Quai.MauId).
//
// ⚠️ CÁC CON SỐ DƯỚI ĐÂY LÀ KHỞI ĐẦU DO NGƯỜI VIẾT ĐẶT, không trích từ bản
// Java — bản Java là game Dragon Ball, đây là game Naruto. Cân bằng tỉ lệ thật
// thuộc KH8. Chúng cố tình rộng rãi để nghiệm thu tay (task 13) không phải
// đứng đánh quái nửa tiếng mới thấy một món.
//
// Mẫu vật phẩm tra ở mauvatpham.go: 0 = thuốc hồi máu nhỏ, 1 = thuốc hồi
// chakra nhỏ, 2 = áo, 3 = quần, 4 = kunai luyện tập.
var roiDoTheoQuai = map[byte][]MucRoiDo{
	// Quái mẫu 0 (loại yếu nhất) rơi ít hơn — nó chết nhanh và đông.
	0: {
		{MauID: 0, SoLuong: 1, TiLe: 1000}, // 10% thuốc máu
		{MauID: 1, SoLuong: 1, TiLe: 500},  // 5% thuốc chakra
	},
}

// roiDoQuaiMacDinh dùng cho mọi MauId chưa có dòng riêng ở trên. Có mặc định
// (thay vì "không có bảng thì không rơi gì") vì thêm một mẫu quái mới mà quên
// bảng rơi là lỗi im lặng khó thấy nhất trong cả hệ thống này.
var roiDoQuaiMacDinh = []MucRoiDo{
	{MauID: 0, SoLuong: 1, TiLe: 1500}, // 15% thuốc máu
	{MauID: 1, SoLuong: 1, TiLe: 1000}, // 10% thuốc chakra
	{MauID: 4, SoLuong: 1, TiLe: 200},  // 2% kunai
}

// roiDoTheoBoss là bảng rơi của boss THƯỜNG (mỗi pha chết đều chạy qua đây).
var roiDoTheoBoss = []MucRoiDo{
	{MauID: 0, SoLuong: 1, TiLe: roiDoMauSo}, // chắc chắn
	{MauID: 1, SoLuong: 1, TiLe: roiDoMauSo}, // chắc chắn
	{MauID: 4, SoLuong: 1, TiLe: 3000},       // 30%
}

// roiDoTheoBossLon là phần THÊM cho pha `LaBossLon` — trang bị, thứ không bao
// giờ rơi từ quái thường.
var roiDoTheoBossLon = []MucRoiDo{
	{MauID: 2, SoLuong: 1, TiLe: 5000}, // 50% áo
	{MauID: 3, SoLuong: 1, TiLe: 5000}, // 50% quần
}

// roiDoCuaQuai tra bảng rơi cho một mẫu quái. Hàm THUẦN.
func roiDoCuaQuai(mauId byte) []MucRoiDo {
	if bang, co := roiDoTheoQuai[mauId]; co {
		return bang
	}
	return roiDoQuaiMacDinh
}

// roiDoCuaBoss tra bảng rơi cho MỘT PHA boss. Hàm THUẦN.
//
// pha == nil (dữ liệu pha hỏng — bossChet chịu được nil ở mọi bước) -> không
// rơi gì. Thà không rơi còn hơn rơi theo một pha không biết là pha nào.
func roiDoCuaBoss(pha *PhaBoss) []MucRoiDo {
	if pha == nil {
		return nil
	}
	if !pha.LaBossLon {
		return roiDoTheoBoss
	}
	// Nối hai bảng thành slice MỚI — KHÔNG append thẳng vào roiDoTheoBoss:
	// append lên một biến gói dùng chung là cách kinh điển để hai lần gọi liên
	// tiếp ghi đè lẫn nhau (nếu cap dư) và bảng gói bị bẩn vĩnh viễn.
	gop := make([]MucRoiDo, 0, len(roiDoTheoBoss)+len(roiDoTheoBossLon))
	gop = append(gop, roiDoTheoBoss...)
	gop = append(gop, roiDoTheoBossLon...)
	return gop
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
