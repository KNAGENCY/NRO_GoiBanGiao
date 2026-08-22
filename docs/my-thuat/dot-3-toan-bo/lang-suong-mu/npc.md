# NPC — Làng Sương Mù (Kiri, hệ Thuỷ) 🟡 chờ server chốt

> **Đọc `../_KHUON.md` trước.** 🟡 CHỜ: server chưa chốt NPC Làng Sương Mù — bộ dưới là
> **đề xuất đại diện** (Mizukage + thầy + người giao NV), chốt lại khi có nội dung.

**Số ảnh: 6** — 3 NPC toàn thân (96×128) + 3 chân dung (200×200). Nền trong suốt.

## Phong cách
Hệ Thuỷ: **xanh dương – xanh lam – xám đá**, áo choàng chống ẩm, mặt nạ/khăn che nửa mặt kiểu sương.
Nhìn nghiêng side-view, mặt xoay **PHẢI**, toàn thân, có khối + đổ bóng (≥40 màu).
Tham chiếu người: `../anh-tham-chieu/TOT-npc-ong-tazuna.png` (vải có nếp gấp, đổ bóng, phụ kiện).
Chân dung `cd_` là **cận mặt** cho khung hội thoại, cùng nhân vật với npc tương ứng.

**Subject block:** `a Kirigakure water-village ninja NPC, deep-blue and slate-grey robes, damp misty look, side view facing right, full body, floating on neutral background.`

## Danh sách

| Tên file | Cỡ | Mô tả chi tiết |
|---|---|---|
| `npc_smu_mizukage_x4.png` | 96×128 | **Thuỷ Ảnh (Mizukage)**: áo choàng Kage trắng–xanh dương, mũ Kage chóp khắc **4 đường sương lượn sóng**, khăn choàng lam nhạt, dáng uy nghi, tay khoanh. *(Bản boss chiến đấu ở `boss.md` — cùng nhân vật, khác vai.)* |
| `npc_smu_thay_x4.png` | 96×128 | **Thầy huấn luyện Kiri**: áo giáp Jonin xanh dương–xám ẩm, băng trán khắc 4 đường sương, mặt nạ che nửa mặt, tay cầm cuộn giấy. Trung niên. |
| `npc_smu_giao_nv_x4.png` | 96×128 | **Người giao nhiệm vụ**: ngư dân/thương nhân biển áo choàng xanh xám sờn, nón lá rộng ven biển, túi lưới đeo chéo, tay chỉ đường. |
| `cd_smu_mizukage_x4.png` | 200×200 | Chân dung Mizukage: cận mặt nghiêm, mũ Kage, tông xanh dương, ánh mắt sắc lạnh. |
| `cd_smu_thay_x4.png` | 200×200 | Chân dung thầy Kiri: cận mặt, mặt nạ che nửa mặt để lộ mắt, băng trán 4 đường sương. |
| `cd_smu_giao_nv_x4.png` | 200×200 | Chân dung ngư dân: cận mặt dày dạn sương gió biển, nón lá, da tái lạnh. |

## Tự kiểm
```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-của-bạn> nguon
```
Tên mới (`npc_smu_*`, `cd_smu_*`) → báo thêm vào `../../kich-thuoc-anh.json` (npc 96×128, cd 200×200). Đừng tự đổi tên.
