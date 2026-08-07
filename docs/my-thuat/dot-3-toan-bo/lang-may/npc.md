# NPC — Làng Mây (Kumo, hệ Lôi) 🟡 chờ server chốt

> **Đọc `../_KHUON.md` trước.** 🟡 CHỜ: server chưa chốt NPC Làng Mây — bộ dưới là
> **đề xuất đại diện** (thầy + Raikage + người giao NV), chốt lại khi có nội dung.

**Số ảnh: 6** — 3 NPC toàn thân (96×128) + 3 chân dung (200×200). Nền trong suốt.

## Phong cách
Hệ Lôi: **tím – trắng xanh – vàng chớp**, áo choàng núi cao, cơ bắp lực lưỡng.
Nhìn nghiêng side-view, mặt xoay **PHẢI**, toàn thân, có khối + đổ bóng (≥40 màu).
Tham chiếu người: `../anh-tham-chieu/TOT-npc-ong-tazuna.png` (vải có nếp gấp, đổ bóng, phụ kiện).
Chân dung `cd_` là **cận mặt** cho khung hội thoại, cùng nhân vật với npc tương ứng.

**Subject block:** `a Kumogakure mountain ninja NPC, purple and white-blue robes, cloud-emblem forehead protector, side view facing right, full body, floating on neutral background.`

## Danh sách

| Tên file | Cỡ | Mô tả chi tiết |
|---|---|---|
| `npc_may_raikage_x4.png` | 96×128 | **Lôi Ảnh (Raikage)**: áo choàng Kage trắng-vàng chớp, mũ Kage chóp có chữ "Lôi", thân hình cơ bắp lực lưỡng, khăn choàng tím nhạt, dáng uy nghi, tay khoanh. Lớn tuổi, tóc bạc dựng. |
| `npc_may_thay_x4.png` | 96×128 | **Thầy huấn luyện Kumo**: áo vest chiến thuật xám xanh, băng trán khắc khối mây, thân cơ bắp, tay cầm cuộn giấy. Trung niên, dày dạn. |
| `npc_may_giao_nv_x4.png` | 96×128 | **Người giao nhiệm vụ**: sứ giả Kumo áo choàng xám sờn, khăn quàng, túi da đeo chéo đựng cuộn thư, tay chỉ đường. |
| `cd_may_raikage_x4.png` | 200×200 | Chân dung Raikage: cận mặt nghiêm, mũ Kage/tóc bạc dựng, tông tím-vàng chớp, ánh mắt sắc. |
| `cd_may_thay_x4.png` | 200×200 | Chân dung thầy Kumo: cận mặt dày dạn, băng trán khối mây, hàm bạnh. |
| `cd_may_giao_nv_x4.png` | 200×200 | Chân dung sứ giả: cận mặt khắc khổ, khăn quàng, da rám gió núi. |

## Tự kiểm
```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-của-bạn> nguon
```
Tên mới (`npc_may_*`, `cd_may_*`) → báo thêm vào `../../kich-thuoc-anh.json` (npc 96×128, cd 200×200). Đừng tự đổi tên.
