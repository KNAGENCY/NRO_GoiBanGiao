# BOSS — Làng Sương Mù (Kiri, hệ Thuỷ) 🟡 chờ server chốt

> **Đọc `../_KHUON.md` trước.** 🟡 CHỜ: server chưa chốt boss Làng Sương Mù — 3 boss dưới là
> **đề xuất đại diện**, chốt lại khi có nội dung. Vẽ 3 con này trước, đừng bịa thêm.

**Số ảnh: 3.** Kích thước **160×160** (x4), tỉ lệ 1:1. Nền trong suốt.

## Phong cách
Hệ Thuỷ: **xanh dương – xanh lam – xám đá – bạc lạnh**. Boss là "người chơi do máy điều khiển" nên
vẽ như **nhân vật oai vệ, dữ dằn hơn quái thường**, chi tiết cao (≥40 màu, đổ bóng), có hơi nước/sương.
Nhìn nghiêng side-view, mặt xoay **PHẢI**, toàn thân.
Tham chiếu: `../anh-tham-chieu/TOT-boss-mau.png` (mức chi tiết boss).

**Subject block:** `a powerful water-village ninja boss, imposing pose, water and mist motif, cold blue palette, side view facing right, full body, detailed shading, floating on neutral background.`

## Danh sách

| Tên file | Cỡ | Mô tả chi tiết |
|---|---|---|
| `boss_smu_mizukage_x4.png` | 160×160 | **Thuỷ Ảnh (Mizukage)**: áo choàng Kage trắng–xanh dương, mũ Kage chóp khắc **4 đường sương lượn sóng** (ký hiệu Kiri), khăn choàng lam nhạt, dòng nước/chakra Thuỷ cuộn quanh tay. Uy nghi, sắc lạnh. *(Bản NPC thân thiện cho khung hội thoại ở `npc.md` — cùng nhân vật, khác vai.)* |
| `boss_smu_thatdao_x4.png` | 160×160 | **Kiếm sĩ Thất Đao** (Ninja Kiếm Kiri): chiến binh vai u thịt bắp, vác **đại kiếm** khổng lồ bản to có khuyết tròn (kiểu đao chặt đầu), áo giáp xanh xám ẩm, mặt nạ/băng trán 4 vạch sóng, tư thế vác kiếm dữ dằn. Xám thép + xanh lam. |
| `boss_smu_tamvi_x4.png` | 160×160 | **Quái Tam Vĩ (Isobu)**: quái rùa khổng lồ mai gai xanh xám, **ba đuôi** dạng xúc tu/tôm, một sừng nhọn trên đầu, mắt đỏ, thân nửa giáp nửa nước, hơi nước bốc quanh mai. Hung tợn. |

## Tự kiểm
```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-của-bạn> nguon
```
Tên mới (`boss_smu_*`) → báo thêm vào `../../kich-thuoc-anh.json` (đều 160×160). Đừng tự đổi tên.
