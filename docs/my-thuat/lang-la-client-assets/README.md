# Art client Làng Lá — bản sao lưu (map 1)

Client Unity (`D:\NRO_Project`) **KHÔNG nằm trong git repo này**, nên toàn bộ art đã cắm vào game được sao lưu ở đây để không mất khi đổi máy / Unity reimport.

Ngày: 05/08/2026 · Nhánh: `map-lang-la` · Thực hiện qua Unity MCP + tự soi screenshot.

## Nội dung (byte thật, khôi phục là chạy)

| Thư mục | File | Là gì |
|---|---|---|
| `Logo/1.png` | 1 | Wordmark **NARUTO** (thay logo "Ngọc Rồng Kame") — xem bộ nhớ `logo-resource-sync` |
| `res/x2/bg/`, `res/x4/bg/` | b00–b03 | **Nền map 1**: b03=trời (`map_nen_troi_1`), b02=núi Hokage (`map_nen_nuihokage` + padding đẩy lên), b00/b01=đất (màu fill vực) |
| `res/x2/t/`, `res/x4/t/` | `1$V.png` (30 mỗi mức) | **Tile mặt đất map 1** (Lát 2): các value đất đã thay = `map_la_dat_1`; cỏ (7/9/11/12) giữ nguyên DB |
| `res/x{1..4}/t/` | `1$40..1$50.png` | **Tile NHÀ CỬA (Lát 3)**: 40-42 mái ngói · 43-44 tường gỗ · 45 cửa · 46 cửa sổ · 47-48 hàng rào · 49 đèn lồng · 50 cổng làng (96×96 ở x2, tràn 2×2 ô) |
| `res/mymap/1.bytes` | 1 | **Bố cục map 1 (Lát 3)** — 52×24; đã đặt 28 ô nhà cửa ở hàng 15-16 (3 nhà + cổng torii + hàng rào). Bản trước Lát 3: `1.bytes.bak-truoc-lat3` trong client |

## Khôi phục vào client

1. Chạy `KHOI-PHUC.ps1` (copy các file này vào đúng chỗ trong `D:\NRO_Project\Assets`).
2. Mở Unity. **BẮT BUỘC bật Read/Write cho 8 file `res/x{2,4}/bg/b0*.png`** (nếu không, trời sẽ ĐEN — code lấy màu trời từ pixel ảnh qua `getRGB`, texture không readable trả 0). Cách: chọn 8 file → Inspector → tick **Read/Write Enabled** → Apply. (Hoặc qua MCP `TextureImporter.isReadable=true`.)
3. Vào map 1 (Làng Lá) — nền + mặt đất Naruto hiện.

## Cơ chế (tóm tắt — chi tiết ở bộ nhớ `map-render-client`)

- Map vẽ ở **client**, KHÔNG qua `part.bin`: bố cục `res/mymap/1.bytes` + tile `res/x{zoom}/t/1$V.png` + nền `res/x{zoom}/bg/b{typeBg}{k}.png`.
- Ảnh **ưu tiên Resources** (thắng RMS cache cũ). Client hiện `zoom=2` → dùng `res/x2/…`.
- Đây là **Lát 1 + 2** của kế hoạch `docs/superpowers/specs/2026-08-04-cam-map-lang-la.md`. Còn Lát 3 (nhà cửa) + Lát 4 (va chạm).
