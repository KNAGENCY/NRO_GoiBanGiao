# NHÂN VẬT (trang phục hiển thị) — Làng Đá 🟡 phần khung hình chờ

> **Đọc `../_KHUON.md` trước.** Đây là **part vẽ LÊN nhân vật người chơi** (khác `trang-bi.md`
> là icon trong túi đồ). ⚠️ **Loại khó nhất:** nhân vật thật cần nhiều khung hình
> (đứng/đi/đánh) — phần đó **tính riêng sau khi bộ đứng này ổn** (ghi chú GOOGLE-FLOW).
> Đợt này chỉ làm **bộ part TƯ THẾ ĐỨNG** đại diện. 🟡 chờ server chốt số bộ/cấp.

**Số ảnh: 6** (đề xuất) — 3 bộ × (đầu 128×128 + thân 128×160). Nền trong suốt.

## Phong cách
Ninja núi đá Iwa: **nâu đất – xám đá – be**, áo choàng thô ráp chống bụi đá, giáp đá thô kệch.
Nhìn nghiêng side-view, mặt xoay **PHẢI**, part ghép khít (đầu chồng lên thân).
Tham chiếu mức chi tiết: `../anh-tham-chieu/TOT-chitiet-nhanvat.png`.

**Subject block:** `an Iwa stone-village ninja character part (head or torso), side view facing right, earth-brown and rock-grey palette, clean seams for layering, floating on neutral background.`

## Danh sách (3 bộ theo cấp)

| Tên file | Cỡ | Mô tả |
|---|---|---|
| `tp_da_tanbinh_dau_x4.png` | 128×128 | Đầu tân binh Đá: khăn trùm nâu đất, băng trán hai phiến đá, mặt trẻ. |
| `tp_da_tanbinh_than_x4.png` | 128×160 | Thân tân binh: áo choàng nâu đất cổ cao, thắt lưng, tư thế đứng nghiêng phải. |
| `tp_da_chunin_dau_x4.png` | 128×128 | Đầu chunin Đá: khăn trùm kín hơn, dây quấn quanh trán. |
| `tp_da_chunin_than_x4.png` | 128×160 | Thân chunin: áo giáp đá Iwa ốp phiến đá xám vai, túi đồ nghề, đứng nghiêng phải. |
| `tp_da_jonin_dau_x4.png` | 128×128 | Đầu jonin Đá: mũ/khăn cao cấp, ánh mắt sắc, quai hàm cứng. |
| `tp_da_jonin_than_x4.png` | 128×160 | Thân jonin: áo choàng Iwa dày, giáp vai bằng phiến đá, đứng oai. |

## Tự kiểm
```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-của-bạn> nguon
```
Tên mới (`tp_da_*`) → báo thêm vào `../../kich-thuoc-anh.json` (đầu 128×128, thân 128×160). Đừng tự đổi tên.

**Chưa làm đợt này:** các khung hình đi/đánh (animation) — chờ bộ đứng duyệt xong.
