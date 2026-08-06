# NHÂN VẬT (trang phục hiển thị) — Làng Cát 🟡 phần khung hình chờ

> **Đọc `../_KHUON.md` trước.** Đây là **part vẽ LÊN nhân vật người chơi** (khác `trang-bi.md`
> là icon trong túi đồ). ⚠️ **Loại khó nhất:** nhân vật thật cần nhiều khung hình
> (đứng/đi/đánh) — phần đó **tính riêng sau khi bộ đứng này ổn** (ghi chú GOOGLE-FLOW).
> Đợt này chỉ làm **bộ part TƯ THẾ ĐỨNG** đại diện. 🟡 chờ server chốt số bộ/cấp.

**Số ảnh: 6** (đề xuất) — 3 bộ × (đầu 128×128 + thân 128×160). Nền trong suốt.

## Phong cách
Ninja sa mạc Suna: **vàng cát – nâu – lục nhạt**, áo choàng chống cát, khăn trùm.
Nhìn nghiêng side-view, mặt xoay **PHẢI**, part ghép khít (đầu chồng lên thân).
Tham chiếu mức chi tiết: `../anh-tham-chieu/TOT-chitiet-nhanvat.png`.

**Subject block:** `a Suna desert ninja character part (head or torso), side view facing right, sand palette, clean seams for layering, floating on neutral background.`

## Danh sách (3 bộ theo cấp)

| Tên file | Cỡ | Mô tả |
|---|---|---|
| `tp_cat_tanbinh_dau_x4.png` | 128×128 | Đầu tân binh Cát: khăn trùm vàng cát, băng trán đồng hồ cát, mặt trẻ. |
| `tp_cat_tanbinh_than_x4.png` | 128×160 | Thân tân binh: áo choàng cát cổ cao, thắt lưng, tư thế đứng nghiêng phải. |
| `tp_cat_chunin_dau_x4.png` | 128×128 | Đầu chunin Cát: khăn trùm kín hơn, khăn che nửa mặt. |
| `tp_cat_chunin_than_x4.png` | 128×160 | Thân chunin: áo giáp gió xanh lục nhạt, túi cuộn, đứng nghiêng phải. |
| `tp_cat_jonin_dau_x4.png` | 128×128 | Đầu jonin Cát: mũ/khăn cao cấp, ánh mắt sắc. |
| `tp_cat_jonin_than_x4.png` | 128×160 | Thân jonin: áo choàng Suna dày, giáp vai, đứng oai. |

## Tự kiểm
```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-của-bạn> nguon
```
Tên mới (`tp_cat_*`) → báo thêm vào `../../kich-thuoc-anh.json` (đầu 128×128, thân 128×160). Đừng tự đổi tên.

**Chưa làm đợt này:** các khung hình đi/đánh (animation) — chờ bộ đứng duyệt xong.
