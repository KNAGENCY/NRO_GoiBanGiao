package store

// roido_test.go — KH7 Task 12, phần CSDL THẬT.
//
// Task 11 đã có TestSinhDoDuoiDatGhiCaMonLanLichSu (đếm số dòng lịch sử). Test
// dưới đây kiểm thứ mà task 12 dựa vào và task 11 chưa chốt: dòng lịch sử đó
// phải là `SINH_RA` với CẢ HAI đầu `tu_ai`/`sang_ai` đều NULL (BB2 + §5.1).
//
// VÌ SAO PHẢI CHỐT: món rơi từ quái không đến từ tay ai và cũng chưa vào tay
// ai. Nếu ai đó "sửa cho gọn" bằng cách điền id người hạ quái vào `sang_ai`,
// bảng lịch sử sẽ nói món đã thuộc về người đó trong khi nó vẫn nằm dưới đất —
// mọi truy vấn kiểm toán sau này (ai đang giữ gì) đọc sai mà không có lỗi nào.

import (
	"context"
	"testing"
)

func TestSinhDoDuoiDatGhiLichSuSinhRaHaiDauNull(t *testing.T) {
	s := moKhoVatPham(t)
	ctx := context.Background()

	id, err := s.SinhDoDuoiDat(ctx, 3, 2, khongOpt())
	if err != nil {
		t.Fatalf("sinh do duoi dat that bai: %v", err)
	}
	donDoDuoiDat(t, s, id)

	var lyDo string
	var tuAi, sangAi *int64
	var soLuong int32
	err = s.pool.QueryRow(ctx,
		`SELECT ly_do, tu_ai, sang_ai, so_luong FROM vat_pham_lich_su WHERE vat_pham_id = $1`,
		id).Scan(&lyDo, &tuAi, &sangAi, &soLuong)
	if err != nil {
		t.Fatalf("doc lich su that bai: %v", err)
	}

	if lyDo != LyDoSinhRa {
		t.Fatalf("ly_do = %q, mong doi %q", lyDo, LyDoSinhRa)
	}
	if tuAi != nil {
		t.Fatalf("tu_ai phai NULL (mon khong den tu tay ai), nhan %d", *tuAi)
	}
	if sangAi != nil {
		t.Fatalf("sang_ai phai NULL (mon chua vao tay ai), nhan %d", *sangAi)
	}
	if soLuong != 2 {
		t.Fatalf("so_luong trong lich su = %d, mong doi 2", soLuong)
	}
}
