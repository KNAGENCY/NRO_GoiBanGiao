# NHÂN VẬT (trang phục hiển thị) — Làng Sương Mù 🟡 phần khung hình chờ

> **Đọc `../_KHUON.md` trước.** Đây là **part vẽ LÊN nhân vật người chơi** (khác `trang-bi.md`
> là icon trong túi đồ). ⚠️ **Loại khó nhất:** nhân vật thật cần nhiều khung hình
> (đứng/đi/đánh) — phần đó **tính riêng sau khi bộ đứng này ổn** (ghi chú GOOGLE-FLOW).
> Đợt này chỉ làm **bộ part TƯ THẾ ĐỨNG** đại diện. 🟡 chờ server chốt số bộ/cấp.

**Số ảnh: 6** (đề xuất) — 3 bộ × (đầu 128×128 + thân 128×160). Nền trong suốt.

## Phong cách
Ninja đảo–biển Kiri: **xanh dương – xanh lam – xám đá**, áo khoác cổ cao chống ẩm, mặt nạ/khăn
che nửa mặt kiểu sương mù. Nhìn nghiêng side-view, mặt xoay **PHẢI**, part ghép khít (đầu chồng lên thân).
Tham chiếu mức chi tiết: `../anh-tham-chieu/TOT-chitiet-nhanvat.png`.

**Subject block:** `a Kiri water-village ninja character part (head or torso), side view facing right, deep-blue and slate-grey palette, damp misty look, clean seams for layering, floating on neutral background.`

## Danh sách (3 bộ theo cấp)

| Tên file | Cỡ | Mô tả |
|---|---|---|
| `tp_smu_tanbinh_dau_x4.png` | 128×128 | Đầu tân binh Kiri: tóc ẩm, băng trán khắc **4 đường sương lượn sóng**, mặt trẻ. |
| `tp_smu_tanbinh_than_x4.png` | 128×160 | Thân tân binh: áo khoác xanh lam cổ cao, thắt lưng, tư thế đứng nghiêng phải. |
| `tp_smu_chunin_dau_x4.png` | 128×128 | Đầu chunin Kiri: khăn/mặt nạ che nửa mặt kiểu sương, băng trán 4 vạch sóng. |
| `tp_smu_chunin_than_x4.png` | 128×160 | Thân chunin: áo lưới thép xanh lạnh + giáp nhẹ xanh dương–xám, túi cuộn, đứng nghiêng phải. |
| `tp_smu_jonin_dau_x4.png` | 128×128 | Đầu jonin Kiri: mặt nạ/khăn cao cấp, ánh mắt sắc lạnh. |
| `tp_smu_jonin_than_x4.png` | 128×160 | Thân jonin: áo giáp Jonin Kiri dày xanh dương–xám, giáp vai, dây đai chéo, đứng oai. |

## Tự kiểm
```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-của-bạn> nguon
```
Tên mới (`tp_smu_*`) → báo thêm vào `../../kich-thuoc-anh.json` (đầu 128×128, thân 128×160). Đừng tự đổi tên.

**Chưa làm đợt này:** các khung hình đi/đánh (animation) — chờ bộ đứng duyệt xong.
