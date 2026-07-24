package game

// NguongBuoc: mỗi bước di chuyển không vượt quá 500 đơn vị (tránh teleport chống lag).
// Đơn vị: pixel/cell trên bản đồ.
const NguongBuoc = 500

// BienRongToiDa: chiều rộng an toàn của bản đồ (tạm thời dùng hằng số).
// TODO: nạp kích thước thật của mỗi map từ dữ liệu (đặc tả §5).
// Hiện tại: x,y ∈ [0, 4096].
const BienRongToiDa = 4096

// kiemDiChuyenHopLe: kiểm tra xem di chuyển từ (cuX, cuY) tới (toX, toY) có hợp lệ không.
// Hợp lệ nếu: (1) đích nằm trong biên bản đồ, (2) khoảng cách ≤ NguongBuoc.
//
// Tham số bd *BanDo: hiện chưa sử dụng (biên là hằng số). Dự trữ cho kế hoạch 6
// khi nạp kích thước map thật từ dữ liệu thay vì hằng số.
//
// Khoảng cách được tính dưới dạng bình phương (dx²+dy²) để tránh sqrt/float,
// giảm chi phí tính toán và tránh sai số làm tròn. So sánh: dx²+dy² ≤ NguongBuoc².
func kiemDiChuyenHopLe(bd *BanDo, cuX, cuY, toX, toY int16) bool {
	// Kiểm biên: đích phải nằm trong [0, BienRongToiDa] cả X lẫn Y
	if toX < 0 || toX > BienRongToiDa || toY < 0 || toY > BienRongToiDa {
		return false
	}

	// Ép kiểu int để tránh tràn int16 khi bình phương
	// Ví dụ: nếu giữ int16, (400)² = 160000 sẽ overflow
	dx := int(toX) - int(cuX)
	dy := int(toY) - int(cuY)

	// So bình phương: khoảng cách ≤ NguongBuoc ⟺ dx²+dy² ≤ NguongBuoc²
	distSq := dx*dx + dy*dy
	thresholdSq := NguongBuoc * NguongBuoc

	return distSq <= thresholdSq
}
