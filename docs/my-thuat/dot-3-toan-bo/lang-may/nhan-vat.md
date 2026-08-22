# NHÂN VẬT (trang phục hiển thị) — Làng Mây 🟡 phần khung hình chờ

> **Đọc `../_KHUON.md` trước.** Đây là **part vẽ LÊN nhân vật người chơi** (khác `trang-bi.md`
> là icon trong túi đồ). ⚠️ **Loại khó nhất:** nhân vật thật cần nhiều khung hình
> (đứng/đi/đánh) — phần đó **tính riêng sau khi bộ đứng này ổn** (ghi chú GOOGLE-FLOW).
> Đợt này chỉ làm **bộ part TƯ THẾ ĐỨNG** đại diện. 🟡 chờ server chốt số bộ/cấp.

**Số ảnh: 6** (đề xuất) — 3 bộ × (đầu 128×128 + thân 128×160). Nền trong suốt.

## Phong cách
Ninja núi cao Kumo: **tím – trắng xanh – vàng chớp**, cơ bắp lực lưỡng, áo cổ cao,
nhiều người đeo **song đao lôi** sau lưng. Nhìn nghiêng side-view, mặt xoay **PHẢI**,
part ghép khít (đầu chồng lên thân). Tham chiếu mức chi tiết: `../anh-tham-chieu/TOT-chitiet-nhanvat.png`.

**Subject block:** `a Kumo mountain ninja character part (head or torso), side view facing right, thunder palette, clean seams for layering, floating on neutral background.`

## Danh sách (3 bộ theo cấp)

| Tên file | Cỡ | Mô tả |
|---|---|---|
| `tp_may_tanbinh_dau_x4.png` | 128×128 | Đầu tân binh Mây: tóc dựng, băng trán khối mây, mặt trẻ. |
| `tp_may_tanbinh_than_x4.png` | 128×160 | Thân tân binh: áo Kumo cổ cao trắng xanh, thắt lưng, tư thế đứng nghiêng phải. |
| `tp_may_chunin_dau_x4.png` | 128×128 | Đầu chunin Mây: tóc/khăn gọn, ánh mắt cứng cỏi. |
| `tp_may_chunin_than_x4.png` | 128×160 | Thân chunin: áo giáp lôi tím ánh điện, song đao lôi chéo sau lưng, đứng nghiêng phải. |
| `tp_may_jonin_dau_x4.png` | 128×128 | Đầu jonin Mây: nét dày dạn, ánh mắt sắc. |
| `tp_may_jonin_than_x4.png` | 128×160 | Thân jonin: áo choàng Kumo dày, thân hình cơ bắp, giáp vai, đứng oai. |

## Tự kiểm
```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-của-bạn> nguon
```
Tên mới (`tp_may_*`) → báo thêm vào `../../kich-thuoc-anh.json` (đầu 128×128, thân 128×160). Đừng tự đổi tên.

**Chưa làm đợt này:** các khung hình đi/đánh (animation) — chờ bộ đứng duyệt xong.
