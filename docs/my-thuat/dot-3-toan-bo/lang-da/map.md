# MAP — Làng Đá (Iwa, núi đá – hẻm núi – hang động) 🟢 làm ngay

> **Đọc `../_KHUON.md` trước** (đặc biệt §7 cách lắp tile/nền vào game).
> Bản đồ núi đá Làng Đá: **tile mặt đất/kiến trúc** (nền TRONG SUỐT) + **nền xa** (nền ĐẶC).

**Số ảnh: 16** — 14 tile 96×96 + 2 nền 384×192.

## Phong cách
Núi đá Iwa: **nâu đất – xám đá – be**, đá xám xù xì, tường đá xếp khối, hẻm núi dựng đứng,
hang động tối. Tile phải **ghép liền mạch** (tileable) khi lặp. Nền xa (`map_nen_*`) vẽ đặc,
phần trên sáng (màu trời). Tham chiếu tile Làng Lá đã có ở `../../anh-game/x4/map_la_*` để
đồng mức chi tiết.

**Subject block (tile):** `a seamless tileable rocky-mountain ground/architecture tile, Iwagakure style, earth-brown and stone-grey palette, top-down game tile.`

## Tile mặt đất + kiến trúc (96×96, nền trong suốt)

| Tên file | Cỡ | Mô tả |
|---|---|---|
| `map_da_da_1_x4.png` | 96×96 | Nền đá xám lát tự nhiên, gợn vân, xù xì. |
| `map_da_da_2_x4.png` | 96×96 | Đá xám lẫn rêu và vết nứt, tông tối hơn. |
| `map_da_soi_x4.png` | 96×96 | Sỏi đá vụn nâu xám rải đều, bóng đổ lấm tấm. |
| `map_da_datnui_x4.png` | 96×96 | Đất núi nâu đỏ khô cứng, lẫn mảnh đá nhỏ. |
| `map_da_tangda_x4.png` | 96×96 | Tảng đá sa thạch xám lớn nhô lên, cạnh sắc. |
| `map_da_vachda_x4.png` | 96×96 | Mặt vách đá dựng đứng, lớp trầm tích ngang nâu–xám. |
| `map_da_tuongda_1_x4.png` | 96×96 | Tường đá xếp khối màu xám nâu, mạch vữa, vết nứt. |
| `map_da_tuongda_2_x4.png` | 96×96 | Tường đá xếp khối có ô cửa vòm hẹp. |
| `map_da_hang_x4.png` | 96×96 | Cửa hang động tối trong vách đá, miệng hang gồ ghề. |
| `map_da_conglang_x4.png` | 96×96 | Cổng làng Đá: hai trụ đá lớn khắc chữ "Thổ", khe hẹp giữa hai vách đá dựng. |
| `map_da_cauda_x4.png` | 96×96 | Cầu đá bắc qua hẻm núi sâu. |
| `map_da_thachnhu_x4.png` | 96×96 | Măng đá / thạch nhũ trong hang, nâu xám ẩm. |
| `map_da_denduoc_x4.png` | 96×96 | Đuốc gắn giá sắt trên vách hang, ánh cam ấm hắt lên đá. |
| `map_da_bacthang_x4.png` | 96×96 | Bậc thang đá đục thẳng vào sườn núi. |

## Nền xa (384×192, nền ĐẶC — phần trên sáng màu trời)

| Tên file | Cỡ | Mô tả |
|---|---|---|
| `map_nen_nuida_x4.png` | 384×192 | Nền dãy núi đá trập trùng trải xa tới chân trời, các đỉnh đá xám nâu chồng lớp, trời xám xanh nhạt phía trên (mép trên sáng). Tileable ngang. |
| `map_nen_henui_x4.png` | 384×192 | Nền hẻm núi đá dựng đứng hai bên bao quanh làng Iwa, vách đá nâu–xám lớp trầm tích, khe trời sáng hẹp phía trên. |

## Tự kiểm
```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-của-bạn> nguon
```
Tên mới → báo thêm vào `../../kich-thuoc-anh.json` (tile 96×96, nền 384×192). ⚠️ Nền `map_nen_*`
là ảnh ĐẶC, không cần alpha; tile thì nền trong suốt.

> ⚠️ **Tên nền không trùng Làng Cát:** Cát đã dùng `map_nen_samac` + `map_nen_vachda`. Làng Đá
> dùng `map_nen_nuida` + `map_nen_henui` — tên riêng, không đụng. Đừng đặt lại `map_nen_vachda`.
