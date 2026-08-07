# MAP — Làng Sương Mù (Kiri, đảo–biển–sương) 🟢 làm ngay

> **Đọc `../_KHUON.md` trước** (đặc biệt §7 cách lắp tile/nền vào game).
> Bản đồ đảo–biển Làng Sương Mù: **tile mặt đất/kiến trúc** (nền TRONG SUỐT) + **nền xa** (nền ĐẶC).

**Số ảnh: 16** — 14 tile 96×96 + 2 nền 384×192.

## Phong cách
Đảo–biển Kiri: **xanh dương – xanh lam – xám đá – nâu gỗ ẩm**, sương mù dày, mọi thứ ướt
nước, ánh sáng lạnh xám. Gỗ bến/cầu bạc màu ẩm, đá phủ rêu, ngói mái dốc bóng nước. Tile phải
**ghép liền mạch** (tileable) khi lặp. Nền xa (`map_nen_*`) vẽ đặc, phần trên phủ sương (màu trời).
Tham chiếu tile Làng Lá đã có ở `../../anh-game/x4/map_la_*` để đồng mức chi tiết.

**Subject block (tile):** `a seamless tileable water-village ground/architecture tile, Kirigakure island style, deep-blue and slate-grey palette, wet damp surfaces, misty, top-down game tile.`

## Tile mặt đất + kiến trúc (96×96, nền trong suốt)

| Tên file | Cỡ | Mô tả |
|---|---|---|
| `map_smu_nuoc_1_x4.png` | 96×96 | Mặt nước biển nông, gợn sóng lăn tăn, xanh lam sáng, ánh phản chiếu. |
| `map_smu_nuoc_2_x4.png` | 96×96 | Nước sâu hơn, xanh dương đậm, đốm sáng lấp lánh. |
| `map_smu_datuot_x4.png` | 96×96 | Đất ẩm ven biển, nâu xám bóng nước, vũng đọng nhỏ. |
| `map_smu_dauot_x4.png` | 96×96 | Tảng đá ướt rêu phong xám xanh nhô lên, mặt bóng nước. |
| `map_smu_catuot_x4.png` | 96×96 | Bãi cát ướt ven bờ, xám nhạt, vệt sóng vừa rút. |
| `map_smu_govan_x4.png` | 96×96 | Sàn gỗ ván bến tàu ẩm, khe ván hở, đầu đinh gỉ. |
| `map_smu_tuongda_1_x4.png` | 96×96 | Tường đá xám ướt, mạch vữa, rêu bám chân. |
| `map_smu_tuongda_2_x4.png` | 96×96 | Tường đá có ô cửa sổ vòm nhỏ, khung gỗ ẩm. |
| `map_smu_maidoc_x4.png` | 96×96 | Mái nhà dốc ngói xanh xám kiểu Kiri, ướt nước, giọt đọng ở rìa. |
| `map_smu_conglang_x4.png` | 96×96 | Cổng làng Sương Mù: hai trụ đá lớn khắc **4 đường sương lượn sóng**, khe giữa mù dày. |
| `map_smu_caugo_x4.png` | 96×96 | Cầu gỗ bắc qua nước, cọc gỗ + dây thừng, ván ẩm bạc màu. |
| `map_smu_cotgo_x4.png` | 96×96 | Cọc gỗ buộc thuyền ven bến, dây thừng quấn, hà bám gốc. |
| `map_smu_denbao_x4.png` | 96×96 | Đèn hiệu/đèn lồng ven biển treo, ánh xanh lạnh trong sương. |
| `map_smu_bacthang_x4.png` | 96×96 | Bậc thang đá ướt dẫn lên nhà, rêu ở mép, giọt nước. |

## Nền xa (384×192, nền ĐẶC — phần trên phủ sương màu trời)

| Tên file | Cỡ | Mô tả |
|---|---|---|
| `map_nen_bien_x4.png` | 384×192 | Nền biển: mặt biển trải xa tới chân trời, sương mù dày phủ phía trên, trời xám xanh lạnh (mép trên sáng). Tileable ngang. |
| `map_nen_dao_x4.png` | 384×192 | Nền đảo đá: vách đá xám xanh ẩm bao quanh làng Kiri, sương giăng lưng chừng, trời lam nhạt phía trên. |

## Tự kiểm
```powershell
python ..\..\KIEM-TRA-ANH.py <thư-mục-của-bạn> nguon
```
Tên mới → báo thêm vào `../../kich-thuoc-anh.json` (tile 96×96, nền 384×192). ⚠️ Nền `map_nen_*` là ảnh ĐẶC, không cần alpha; tile thì nền trong suốt.
