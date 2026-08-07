# MAP — Làng Mây (Kumo, núi cao chọc mây) 🟢 làm ngay

> **Đọc `../_KHUON.md` trước** (đặc biệt §7 cách lắp tile/nền vào game).
> Bản đồ núi cao Làng Mây: **tile mặt đất/kiến trúc** (nền TRONG SUỐT) + **nền xa** (nền ĐẶC).

**Số ảnh: 16** — 14 tile 96×96 + 2 nền 384×192.

## Phong cách
Núi cao Kumo: **tím – trắng xanh – xám đá**, mây dày là là, tia sét xa, đá granite,
tháp trụ cao, cầu treo. Tile phải **ghép liền mạch** (tileable) khi lặp. Nền xa
(`map_nen_*`) vẽ đặc, phần trên sáng (màu trời). Tham chiếu tile Làng Lá đã có ở
`../../anh-game/x4/map_la_*` để đồng mức chi tiết.

**Subject block (tile):** `a seamless tileable high-mountain ground/architecture tile, Kumogakure style, thunder-cloud palette (grey stone, white-blue, purple), top-down game tile.`

## Tile mặt đất + kiến trúc (96×96, nền trong suốt)

| Tên file | Cỡ | Mô tả |
|---|---|---|
| `map_may_da_1_x4.png` | 96×96 | Đá núi granite xám xanh, gợn vân mờ. |
| `map_may_da_2_x4.png` | 96×96 | Đá núi lẫn khoáng, vệt tinh thể ánh tím. |
| `map_may_datnui_x4.png` | 96×96 | Đất núi cao ẩm, xám nâu, rêu thưa. |
| `map_may_dadung_x4.png` | 96×96 | Tảng đá granite dựng nhô lên, xám xanh. |
| `map_may_manmay_x4.png` | 96×96 | Mảng mây dày là mặt đất, trắng xanh bồng bềnh, bóng mềm. |
| `map_may_tuongda_1_x4.png` | 96×96 | Tường đá núi xám xanh, vết nứt. |
| `map_may_tuongda_2_x4.png` | 96×96 | Tường đá có ô cửa vuông. |
| `map_may_maithap_x4.png` | 96×96 | Mái tháp trụ Kumo, xám xanh, đỉnh nhọn. |
| `map_may_conglang_x4.png` | 96×96 | Cổng làng Mây: trụ đá lớn khắc chữ "Lôi"/biểu tượng mây, lối vào giữa hai vách núi cao. |
| `map_may_cautreo_x4.png` | 96×96 | Cầu treo dây bắc qua vực mây. |
| `map_may_thongnui_x4.png` | 96×96 | Cây thông/tùng núi cao, xanh thẫm bám vách. |
| `map_may_cotloi_x4.png` | 96×96 | Cột thu lôi kim loại cạnh tường, ánh điện lam (trang trí). |
| `map_may_denloi_x4.png` | 96×96 | Đèn lồng phóng điện treo, ánh lam-tím. |
| `map_may_bacthang_x4.png` | 96×96 | Bậc thang đá dẫn lên tháp. |

## Nền xa (384×192, nền ĐẶC — phần trên sáng màu trời)

> ⚠️ **KHÔNG đặt tên `map_nen_may_*`** — tên đó đã dùng cho **mây xa Làng Lá**
> (`map_nen_may_1/2/3`). Nền Làng Mây dùng tiền tố khác như dưới.

| Tên file | Cỡ | Mô tả |
|---|---|---|
| `map_nen_nuimay_x4.png` | 384×192 | Nền núi cao chọc mây: dãy núi nhọn trập trùng vươn lên biển mây dày, trời tím-lam phía trên (mép trên sáng), tia sét xa. Tileable ngang. |
| `map_nen_vachnui_x4.png` | 384×192 | Nền vách núi đá cao bao quanh làng Kumo, xám xanh, mây cuộn chân vách, trời lam nhạt phía trên. |

## Tự kiểm
```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-của-bạn> nguon
```
Tên mới → báo thêm vào `../../kich-thuoc-anh.json` (tile 96×96, nền 384×192). ⚠️ Nền `map_nen_*` là ảnh ĐẶC, không cần alpha; tile thì nền trong suốt.
