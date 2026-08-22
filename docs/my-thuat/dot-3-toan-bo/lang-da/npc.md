# NPC — Làng Đá (Iwa, hệ Thổ) 🟡 chờ server chốt

> **Đọc `../_KHUON.md` trước.** 🟡 CHỜ: server chưa chốt NPC Làng Đá — bộ dưới là
> **đề xuất đại diện** (thầy + Tsuchikage + người giao NV), chốt lại khi có nội dung.

**Số ảnh: 6** — 3 NPC toàn thân (96×128) + 3 chân dung (200×200). Nền trong suốt.

## Phong cách
Hệ Thổ: **nâu đất – xám đá – be**, áo choàng thô ráp núi đá, khăn trùm chống bụi.
Nhìn nghiêng side-view, mặt xoay **PHẢI**, toàn thân, có khối + đổ bóng (≥40 màu).
Tham chiếu người: `../anh-tham-chieu/TOT-npc-ong-tazuna.png` (vải có nếp gấp, đổ bóng, phụ kiện).
Chân dung `cd_` là **cận mặt** cho khung hội thoại, cùng nhân vật với npc tương ứng.

**Subject block:** `an Iwagakure stone-village ninja NPC, earth-brown robes and head wrap, side view facing right, full body, floating on neutral background.`

## Danh sách

| Tên file | Cỡ | Mô tả chi tiết |
|---|---|---|
| `npc_da_tsuchikage_x4.png` | 96×128 | **Thổ Ảnh (Tsuchikage)**: lão nhân nhỏ người, mũi to, ria mép và râu cằm nhọn trắng, áo choàng Kage cam đất viền đỏ, mũ Kage chóp có chữ "Thổ". Dáng uy nghi, hơi lơ lửng. Lớn tuổi. |
| `npc_da_thay_x4.png` | 96×128 | **Thầy huấn luyện Iwa**: áo vest nâu đất, ốp phiến đá vai, băng trán khắc hai phiến đá, tay cầm cuộn giấy. Vạm vỡ, trung niên. |
| `npc_da_giao_nv_x4.png` | 96×128 | **Người giao nhiệm vụ**: thợ mỏ/thợ đá Iwa áo choàng nâu sờn bụi, khăn trùm, cuốc đá đeo lưng, túi da đeo chéo, tay chỉ đường. |
| `cd_da_tsuchikage_x4.png` | 200×200 | Chân dung Tsuchikage: cận mặt nghiêm, mũi to, râu trắng, mũ Kage, tông nâu đất, ánh mắt sắc. |
| `cd_da_thay_x4.png` | 200×200 | Chân dung thầy Iwa: cận mặt cứng cỏi, băng trán hai phiến đá, quai hàm vuông. |
| `cd_da_giao_nv_x4.png` | 200×200 | Chân dung thợ đá: cận mặt khắc khổ, khăn trùm phủ bụi, da sạm ráp. |

## Tự kiểm
```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-của-bạn> nguon
```
Tên mới (`npc_da_*`, `cd_da_*`) → báo thêm vào `../../kich-thuoc-anh.json` (npc 96×128, cd 200×200). Đừng tự đổi tên.
