# QUÁI — Làng Mây (hệ Lôi) 🟡 chờ server chốt danh sách

> **Đọc `../_KHUON.md` trước.** File này chỉ ghi nội dung riêng của quái Làng Mây.
> 🟡 **CHỜ:** server chưa chốt chính xác Làng Mây có bao nhiêu quái / con nào.
> Bộ 8 con dưới là **đề xuất đại diện** (núi cao/sấm sét) — khi server chốt nội dung sẽ
> thêm/bớt. Agent **cứ vẽ 8 con này trước** (chắc chắn dùng được), đừng tự bịa thêm.

**Số ảnh: 8.** Kích thước: **96×96** (x4), tỉ lệ 1:1. Nền trong suốt.

---

## Phong cách (hệ Lôi — Làng Mây)

Bảng màu: **tím · trắng xanh · vàng chớp**, ánh sáng lạnh trên đỉnh núi, mây dày,
tia sét lóe. Quái sống ở núi cao/vùng bão sét. Nhìn nghiêng side-view, mặt xoay sang **PHẢI**,
toàn thân, có khối + đổ bóng (≥40 màu). **Tham chiếu chuẩn:** `../anh-tham-chieu/TOT-quai-soi-rung.png`
(lông/vảy có lớp, cơ có khối, nền cắt sạch). **Tránh:** `XAU-*.png` (phẳng, chân que).

**Subject block (dán kèm master prompt khi gen):**
`a high-mountain thunder beast, storm-charged hide, crackling lightning, cloud-wreathed peak, cold blue light, side view facing right, full body, floating on neutral background.`

---

## Danh sách 8 quái

| Tên file | Cỡ | Mô tả chi tiết (màu + đặc trưng) |
|---|---|---|
| `quai_may_daibang_set_x4.png` | 96×96 | **Đại bàng sét** khổng lồ: lông xám-trắng viền tím, sải cánh rộng, tia điện lách tách quanh lông vũ, móng vuốt vàng phóng điện, mắt vàng chớp, mỏ khoằm. Dữ tợn, đang lao. |
| `quai_may_ho_loi_x4.png` | 96×96 | **Hổ lôi** cơ bắp: bộ lông trắng xanh vằn tím-đen, hồ quang điện chạy dọc sống lưng, nanh trắng nhe, mắt sáng lam, chân sau chùng rình. Khối cơ rõ. |
| `quai_may_bo_tot_x4.png` | 96×96 | **Bò tót sừng điện**: thân to lực lưỡng, lông nâu-xám bết, đôi sừng dài cong **tích điện xanh**, vai u cơ bắp, mũi phì hơi, móng nện đất tung bụi. Tư thế húc. |
| `quai_may_soi_nui_x4.png` | 96×96 | **Sói núi lôi**: sói xám lông dựng, viền điện tím dọc lưng, răng nhe, mắt vàng, đuôi xù. Chân sau chùng, dáng săn mồi. |
| `quai_may_son_duong_x4.png` | 96×96 | **Sơn dương sừng điện**: dê núi cao lông trắng-xám dày, **đôi sừng xoắn tích điện lam**, móng chẻ bám đá, mắt ngang, râu cằm. Đứng trên mỏm đá. |
| `quai_may_golem_set_x4.png` | 96×96 | **Golem sét**: khối đá granite hình người thô, khe nứt lộ **lõi phát điện tím**, hồ quang chạy giữa các khối, mắt sáng lam, mảnh đá rơi khỏi vai. |
| `quai_may_linh_gac_x4.png` | 96×96 | **Lính gác Mây** (ninja Kumo cấp thấp): áo Kumo cổ cao trắng xanh, băng trán khắc khối mây, tay thủ đoản đao ánh điện, cơ bắp. Người, dáng thủ thế. |
| `quai_may_tinhmay_x4.png` | 96×96 | **Tinh linh mây**: khối mây trắng-xám tụ hình thú, lõi phát sáng tím-lam, tia sét lóe bên trong, tua mây bay lơ lửng. Bồng bềnh, hung. |

---

## Tự kiểm trước khi giao (bắt buộc)

```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-ảnh-của-bạn> nguon
```
Dán nguyên kết quả vào báo cáo. **Tên mới** (`quai_may_*`) chưa có trong `kich-thuoc-anh.json`
— báo phía dự án thêm vào bảng (đều 96×96) trước khi công cụ tra được. Không tự đổi tên.

**Chưa cần vẽ đợt này:** biến thể tinh anh / boss của Làng Mây (xem `boss.md`), quái 3 làng còn lại.
