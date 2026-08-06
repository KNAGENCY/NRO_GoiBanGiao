# MAP — Làng Cát (Suna, sa mạc) 🟢 làm ngay

> **Đọc `../_KHUON.md` trước** (đặc biệt §7 cách lắp tile/nền vào game).
> Bản đồ sa mạc Làng Cát: **tile mặt đất/kiến trúc** (nền TRONG SUỐT) + **nền xa** (nền ĐẶC).

**Số ảnh: 16** — 14 tile 96×96 + 2 nền 384×192.

## Phong cách
Sa mạc Suna: **vàng cát – nâu đất – be**, nắng gắt, tường đất nung, mái vòm. Tile phải
**ghép liền mạch** (tileable) khi lặp. Nền xa (`map_nen_*`) vẽ đặc, phần trên sáng (màu trời).
Tham chiếu tile Làng Lá đã có ở `../../anh-game/x4/map_la_*` để đồng mức chi tiết.

**Subject block (tile):** `a seamless tileable desert ground/architecture tile, Sunagakure style, sandy palette, top-down game tile.`

## Tile mặt đất + kiến trúc (96×96, nền trong suốt)

| Tên file | Cỡ | Mô tả |
|---|---|---|
| `map_cat_cat_1_x4.png` | 96×96 | Cát mịn sa mạc, gợn sóng nhẹ, vàng ấm. |
| `map_cat_cat_2_x4.png` | 96×96 | Cát lẫn sỏi nhỏ, vệt gió. |
| `map_cat_datnut_x4.png` | 96×96 | Đất nứt nẻ khô hạn, nâu xám. |
| `map_cat_da_x4.png` | 96×96 | Tảng đá sa thạch nâu vàng nhô lên. |
| `map_cat_concat_x4.png` | 96×96 | Đụn cát cao, bóng đổ mềm. |
| `map_cat_tuongdat_1_x4.png` | 96×96 | Tường đất nung màu nâu cát, vết nứt. |
| `map_cat_tuongdat_2_x4.png` | 96×96 | Tường đất có ô cửa vòm. |
| `map_cat_maivom_x4.png` | 96×96 | Mái vòm nhà Suna tròn, nâu đất. |
| `map_cat_conglang_x4.png` | 96×96 | Cổng làng Cát: trụ đá lớn khắc chữ "Phong", khe hẹp giữa hai vách đá. |
| `map_cat_cauda_x4.png` | 96×96 | Cầu đá bắc qua khe cát. |
| `map_cat_xuongrong_x4.png` | 96×96 | Bụi xương rồng sa mạc, xanh lục nhạt. |
| `map_cat_binhgom_x4.png` | 96×96 | Bình gốm/vò nước đặt cạnh tường (trang trí). |
| `map_cat_dennlong_x4.png` | 96×96 | Đèn lồng sa mạc treo, ánh cam ấm. |
| `map_cat_bacthang_x4.png` | 96×96 | Bậc thang đá dẫn lên nhà. |

## Nền xa (384×192, nền ĐẶC — phần trên sáng màu trời)

| Tên file | Cỡ | Mô tả |
|---|---|---|
| `map_nen_samac_x4.png` | 384×192 | Nền sa mạc: đụn cát trải xa tới chân trời, trời vàng cam nắng gắt phía trên (mép trên sáng). Tileable ngang. |
| `map_nen_vachda_x4.png` | 384×192 | Nền vách đá sa thạch cao bao quanh làng Suna, nâu vàng, trời xanh nhạt phía trên. |

## Tự kiểm
```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-của-bạn> nguon
```
Tên mới → báo thêm vào `../../kich-thuoc-anh.json` (tile 96×96, nền 384×192). ⚠️ Nền `map_nen_*` là ảnh ĐẶC, không cần alpha; tile thì nền trong suốt.
