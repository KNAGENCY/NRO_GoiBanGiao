# QUÁI — Làng Sương Mù (hệ Thuỷ) 🟡 chờ server chốt danh sách

> **Đọc `../_KHUON.md` trước.** File này chỉ ghi nội dung riêng của quái Làng Sương Mù.
> 🟡 **CHỜ:** server chưa chốt chính xác Làng Sương Mù có bao nhiêu quái / con nào.
> Bộ 8 con dưới là **đề xuất đại diện** (biển/đảo/sương) — khi server chốt nội dung sẽ
> thêm/bớt. Agent **cứ vẽ 8 con này trước** (chắc chắn dùng được), đừng tự bịa thêm.

**Số ảnh: 8.** Kích thước: **96×96** (x4), tỉ lệ 1:1. Nền trong suốt.

---

## Phong cách (hệ Thuỷ — Làng Sương Mù)

Bảng màu: **xanh dương · xanh lam · xám đá · bạc lạnh**, mặt da/vảy luôn ướt bóng, hơi sương
trắng xanh bám quanh. Quái sống ở biển–đảo–đầm lầy sương mù. Nhìn nghiêng side-view, mặt xoay sang **PHẢI**,
toàn thân, có khối + đổ bóng (≥40 màu). **Tham chiếu chuẩn:** `../anh-tham-chieu/TOT-quai-soi-rung.png`
(lông/vảy có lớp, cơ có khối, nền cắt sạch). **Tránh:** `XAU-*.png` (phẳng, chân que).

**Subject block (dán kèm master prompt khi gen):**
`a sea/mist monster of a water village, wet glistening hide or scales, cold blue palette, faint white mist, side view facing right, full body, floating on neutral background.`

---

## Danh sách 8 quái

| Tên file | Cỡ | Mô tả chi tiết (màu + đặc trưng) |
|---|---|---|
| `quai_smu_camap_x4.png` | 96×96 | **Cá mập biển đứng** (quái người-cá mập): da xám xanh ướt bóng, vây lưng cao, mang hằn hai bên cổ, hàm răng nhọn tua tủa, mắt trắng dữ, chân vây bám đất. Hung tợn. |
| `quai_smu_luondien_x4.png` | 96×96 | **Lươn điện**: thân dài trơn nhớt xanh đen, tia điện xanh lam chạy dọc thân, miệng nhọn nhe răng, đuôi quẫy. Ánh sáng lạnh quanh mình. |
| `quai_smu_thuyquai_x4.png` | 96×96 | **Thuỷ quái xúc tu** (mực/bạch tuộc): đầu phình xanh lam bóng, xúc tu cuộn giác hút, mắt to đờ, phun tia nước. Nhớt, ẩm. |
| `quai_smu_cua_x4.png` | 96×96 | **Cua biển khổng lồ**: mai giáp xanh xám cứng gai, hai càng to bóng nước, mắt cuống nhô, bọt nước quanh chân. Thủ thế kẹp. |
| `quai_smu_quaisuong_x4.png` | 96×96 | **Quái sương**: hình bóng ninja mờ ảo tụ từ sương mù trắng xanh, mắt phát sáng lạnh, thân dưới tan thành hơi, tay vươn mờ. Bí hiểm. |
| `quai_smu_golem_bang_x4.png` | 96×96 | **Golem băng**: khối băng/nước đóng hình người thô, trong veo ánh xanh, lõi chakra Thuỷ sáng giữa ngực, hơi lạnh bốc, mảnh băng gãy ở vai. |
| `quai_smu_linh_gac_x4.png` | 96×96 | **Lính gác Kiri** (ninja Kiri cấp thấp): áo xanh lam, mặt nạ che mặt để lộ mắt, băng trán khắc 4 đường sương, tay thủ kiếm ngắn. Người, dáng thủ thế. |
| `quai_smu_thuyxa_x4.png` | 96×96 | **Thuỷ xà** (rắn nước lớn) cuộn mình: vảy xanh lam ướt viền lam đậm, mào vây trên đầu, mắt vàng lạnh, miệng phun tia nước. Bóng đổ hình cuộn. |

---

## Tự kiểm trước khi giao (bắt buộc)

```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-ảnh-của-bạn> nguon
```
Dán nguyên kết quả vào báo cáo. **Tên mới** (`quai_smu_*`) chưa có trong `kich-thuoc-anh.json`
— báo phía dự án thêm vào bảng (đều 96×96) trước khi công cụ tra được. Không tự đổi tên.

**Chưa cần vẽ đợt này:** biến thể tinh anh / boss của Làng Sương Mù (xem `boss.md`), quái 4 làng còn lại.
