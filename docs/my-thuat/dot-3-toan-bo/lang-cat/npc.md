# NPC — Làng Cát (Suna, hệ Phong) 🟡 chờ server chốt

> **Đọc `../_KHUON.md` trước.** 🟡 CHỜ: server chưa chốt NPC Làng Cát — bộ dưới là
> **đề xuất đại diện** (thầy + Kazekage + người giao NV), chốt lại khi có nội dung.

**Số ảnh: 6** — 3 NPC toàn thân (96×128) + 3 chân dung (200×200). Nền trong suốt.

## Phong cách
Hệ Phong: **vàng cát – nâu – lục nhạt**, áo choàng sa mạc, khăn trùm chống nắng.
Nhìn nghiêng side-view, mặt xoay **PHẢI**, toàn thân, có khối + đổ bóng (≥40 màu).
Tham chiếu người: `../anh-tham-chieu/TOT-npc-ong-tazuna.png` (vải có nếp gấp, đổ bóng, phụ kiện).
Chân dung `cd_` là **cận mặt** cho khung hội thoại, cùng nhân vật với npc tương ứng.

**Subject block:** `a Sunagakure desert ninja NPC, sand-colored robes and head wrap, side view facing right, full body, floating on neutral background.`

## Danh sách

| Tên file | Cỡ | Mô tả chi tiết |
|---|---|---|
| `npc_cat_kazekage_x4.png` | 96×128 | **Phong Ảnh (Kazekage)**: áo choàng trắng-vàng cát, mũ Kage chóp có chữ "Phong", khăn choàng xanh nhạt, dáng uy nghi, tay khoanh. Lớn tuổi, râu. |
| `npc_cat_thay_x4.png` | 96×128 | **Thầy huấn luyện Suna**: áo vest nâu cát, băng trán khắc đồng hồ cát, khăn che nửa mặt, tay cầm cuộn giấy. Trung niên. |
| `npc_cat_giao_nv_x4.png` | 96×128 | **Người giao nhiệm vụ**: thương nhân sa mạc áo choàng nâu sờn, khăn trùm kín, túi da đeo chéo, tay chỉ đường. |
| `cd_cat_kazekage_x4.png` | 200×200 | Chân dung Kazekage: cận mặt nghiêm, mũ Kage, tông vàng cát, ánh mắt sắc. |
| `cd_cat_thay_x4.png` | 200×200 | Chân dung thầy Suna: cận mặt, khăn che nửa mặt để lộ mắt, băng trán đồng hồ cát. |
| `cd_cat_giao_nv_x4.png` | 200×200 | Chân dung thương nhân: cận mặt khắc khổ, khăn trùm, da rám nắng sa mạc. |

## Tự kiểm
```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-của-bạn> nguon
```
Tên mới (`npc_cat_*`, `cd_cat_*`) → báo thêm vào `../../kich-thuoc-anh.json` (npc 96×128, cd 200×200). Đừng tự đổi tên.
