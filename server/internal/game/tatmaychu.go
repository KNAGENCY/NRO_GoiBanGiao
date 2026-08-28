package game

import (
	"context"
	"errors"
	"time"
)

// TatMayChu thuc hien DUNG THU TU BAT BUOC cua dac ta
// 2026-07-27-luu-trang-thai-song.md §4 diem C (tat may chu) va
// task-6-brief.md Step 3:
//
//  1. soKhuVuc.LuuTatCa(thoGhi) — xin anh chup TOAN BO nguoi choi o TAT CA
//     khu vuc, xep vao tho ghi.
//  2. thoGhi.DongVaChoXong(ctx) — cho tho ghi ghi HET hang doi xuong dia
//     (co han gio qua ctx, dac ta de xuat 30 giay o tang goi — xem main.go).
//  3. RỒI MỚI soKhuVuc.DungTatCa() — dung actor cua tung khu vuc.
//
// VI SAO THU TU NAY TUYET DOI KHONG DUOC DAO NGUOC (dac ta §4 diem C, lap
// lai o task-6-brief.md Step 3): neu DungTatCa() chay TRUOC, actor CUA TUNG
// KHU VUC da dung han — AnhChupNguoiChoi (goi ben trong LuuTatCa) gui su
// kien qua hop thu actor NHUNG khong con ai doc hop thu do de tra loi nua,
// nen sau 100ms het gio no tra ve RONG (xem AnhChupNguoiChoi, khuvuc.go).
// Ket qua: LuuTatCa Xep du lieu RONG cho MOI khu vuc, thoGhi khong co gi de
// ghi, MAT SACH tien do cua TOAN BO nguoi dang choi luc tat may chu. Day
// CHINH XAC la loi da xay ra tren ban Java cu: nut "Thoát" goi thang
// System.exit(0) khong luu ai, con ham luu dung (saveData()) thi ton tai
// nhung khong duoc gan vao nut nao. Khac o Task 6 nay la LuuTatCa/DungTatCa
// deu la ham that da co san — chi can GOI DUNG THU TU la du, khong phai viet
// lai logic luu.
//
// Loi tra ve gop ca loi cua DongVaChoXong (thuong la ctx.Err() khi vuot han
// gio — vd CSDL treo that su dung luc tat may chu) VA loi dung actor. Ham nay
// VAN GOI DungTatCa() du DongVaChoXong co loi hay khong: mot lan tat "danh
// hong" (mat vai thay doi GAN NHAT chua kip ghi) con hon treo ca tien trinh
// may chu khong bao gio thoat duoc — dung tinh than context.WithTimeout da
// giai thich o hanGioGhi (thogi.go).
// hanGioDongSoNo la khoan thoi gian RIENG cho lan xa cuoi cung cua so no luc
// tat may chu — xem ly do trong than TatMayChu.
const hanGioDongSoNo = 10 * time.Second

func TatMayChu(ctx context.Context, soKhuVuc *SoKhuVuc, thoGhi *ThoGhi) error {
	soKhuVuc.LuuTatCa(thoGhi)
	err := thoGhi.DongVaChoXong(ctx)
	errDung := soKhuVuc.DungTatCa()

	// So no dong SAU CUNG, khi khong con actor nao co the xep them khoan nao
	// nua — dong som hon la vut di dung nhung khoan cuoi cung.
	//
	// HAN GIO RIENG, khong dung chung ctx: neu thoGhi vua tieu het ngan sach
	// (mot CSDL cham la du), ctx da het han va DongSoNo se tra ngay ctx.Err()
	// ma khong kip xa lan nao — tien cua nguoi choi bi vut di vi mot buoc TRUOC
	// no cham. Cho no mot khoan rieng, van co tran de khong treo tien trinh.
	ctxNo, huy := context.WithTimeout(context.Background(), hanGioDongSoNo)
	defer huy()
	return errors.Join(err, errDung, soKhuVuc.DongSoNo(ctxNo))
}
