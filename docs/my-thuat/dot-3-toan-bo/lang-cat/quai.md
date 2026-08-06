# QUÁI — Làng Cát (hệ Phong) 🟡 chờ server chốt danh sách

> **Đọc `../_KHUON.md` trước.** File này chỉ ghi nội dung riêng của quái Làng Cát.
> 🟡 **CHỜ:** server chưa chốt chính xác Làng Cát có bao nhiêu quái / con nào.
> Bộ 8 con dưới là **đề xuất đại diện** (sa mạc/gió) — khi server chốt nội dung sẽ
> thêm/bớt. Agent **cứ vẽ 8 con này trước** (chắc chắn dùng được), đừng tự bịa thêm.

**Số ảnh: 8.** Kích thước: **96×96** (x4), tỉ lệ 1:1. Nền trong suốt.

---

## Phong cách (hệ Phong — Làng Cát)

Bảng màu: **vàng cát · nâu đất · xanh lục nhạt**, ánh nắng sa mạc gắt, bụi cát bay.
Quái sống ở sa mạc/hoang mạc gió cát. Nhìn nghiêng side-view, mặt xoay sang **PHẢI**,
toàn thân, có khối + đổ bóng (≥40 màu). **Tham chiếu chuẩn:** `../anh-tham-chieu/TOT-quai-soi-rung.png`
(lông/vảy có lớp, cơ có khối, nền cắt sạch). **Tránh:** `XAU-*.png` (phẳng, chân que).

**Subject block (dán kèm master prompt khi gen):**
`a desert wasteland monster, sand-worn hide, wind-blown dust, sun-bleached palette, side view facing right, full body, floating on neutral background.`

---

## Danh sách 8 quái

| Tên file | Cỡ | Mô tả chi tiết (màu + đặc trưng) |
|---|---|---|
| `quai_cat_bocap_x4.png` | 96×96 | **Bọ cạp sa mạc** cỡ lớn: vỏ giáp vàng cát ánh nâu, đuôi cong giương nọc, càng to bóng loáng, khớp chân bám cát. Có vệt bụi mờ dưới bụng. Dữ tợn. |
| `quai_cat_giun_cat_x4.png` | 96×96 | **Giun cát khổng lồ** trồi lên khỏi mặt cát: thân đốt nâu-vàng nhớt, miệng tròn xoáy răng, quanh miệng cát tung. Thân dưới lẫn vào bụi. |
| `quai_cat_nhen_x4.png` | 96×96 | **Nhện cát** tám cẳng dài mảnh, thân phủ lông vàng cát lẫn nâu, mắt cụm đỏ nhỏ, bụng có hoa văn xoáy gió. Tư thế rình. |
| `quai_cat_kenken_x4.png` | 96×96 | **Kền kền hoang mạc**: lông xám-nâu xơ xác, cổ trụi hồng tái, mỏ khoằm vàng, cánh dang rộng lộ lông vũ rách, đậu trên xương khô. |
| `quai_cat_ran_luc_x4.png` | 96×96 | **Rắn lục sa mạc** cuộn mình: vảy vàng cát viền nâu, mắt vàng con ngươi dọc, lưỡi chẻ, phùng mang nhỏ. Bóng đổ hình cuộn. |
| `quai_cat_golem_cat_x4.png` | 96×96 | **Golem cát**: khối cát nén hình người thô, tay to, nứt nẻ lộ lõi đá nâu, mắt phát sáng lục nhạt (chakra Phong), bụi rơi khỏi vai. |
| `quai_cat_linh_gac_x4.png` | 96×96 | **Lính gác Cát** (ninja Suna cấp thấp): áo choàng nâu cát trùm đầu, khăn che mặt để lộ mắt, băng trán khắc đồng hồ cát, tay thủ kunai. Người, dáng thủ thế. |
| `quai_cat_lac_da_x4.png` | 96×96 | **Lạc đà hoang dữ**: lông nâu vàng bết cát, bờm xù, mắt hằn, miệng nhe, yên da rách bỏ hoang trên lưng. Chân bụi mù. |

---

## Tự kiểm trước khi giao (bắt buộc)

```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-ảnh-của-bạn> nguon
```
Dán nguyên kết quả vào báo cáo. **Tên mới** (`quai_cat_*`) chưa có trong `kich-thuoc-anh.json`
— báo phía dự án thêm vào bảng (đều 96×96) trước khi công cụ tra được. Không tự đổi tên.

**Chưa cần vẽ đợt này:** biến thể tinh anh / boss của Làng Cát (xem `boss.md`), quái 3 làng còn lại.
