# Cắm map Làng Lá vào game (thiết kế + kế hoạch)

**Ngày:** 04/08/2026 · **Trạng thái:** thiết kế (chờ thực thi) · **Nhánh:** `map-lang-la`
**Nguồn điều tra:** 2 subagent Explore (đường ống part.bin + định nghĩa/render map) + thực nghiệm trực tiếp trên client qua Unity MCP (04/08/2026). Bộ nhớ liên quan: `map-render-client`, `logo-resource-sync`.

## 1. Phát hiện cốt lõi + vì sao (đảo giả định ban đầu)

**Map KHÔNG đi qua `part.bin`.** `part.bin` chỉ chứa sprite **nhân vật/skill** (dải ID 27041–27082 là head/body/leg/Iruka/skill). Map được vẽ hoàn toàn ở **client Unity** từ 3 nguồn Resources:

| Thành phần | Đường dẫn client | Nội dung |
|---|---|---|
| Bố cục ô | `res/mymap/{mapID}.bytes` | byte0=`tmw`, byte1=`tmh`, rồi `tmw*tmh` chỉ số tile |
| Ảnh viên gạch | `res/x{zoom}/t/{tileID}/*.png` | ảnh từng ô, khoá theo `tileID` (**thư mục `/t/` CHƯA tồn tại**) |
| Nền parallax | `res/x{zoom}/bg/b{typeBg}{k}.png` (k=0..3) | 4 lớp nền |

**Ảnh nạp qua `loadImageRMS` (`GameCanvas.cs:2474`): thử Resources TRƯỚC, chỉ khi thiếu mới lấy từ RMS cache** → đặt ảnh Resources đúng tên sẽ **thắng cache cũ** (cùng nguyên tắc đã dùng cho logo). Map cỏ/đá hiện thấy trong game = **tàn dư RMS cache của server Dragon Ball cũ**, không phải Resources.

→ Cắm map = **đặt đúng ảnh Naruto vào Resources client + vẽ bố cục ô + (một chút) sửa server Go**. Đây là việc **client-nặng**, khác hẳn assetpack/part.bin.

## 2. Trạng thái hiện tại (đã xác minh 04/08)

- Client Editor **nối đúng server Go `:14445`** (không phải Java cũ; cổng 3306 tắt). `nroserver.exe` đang chạy.
- Hỏi client đang chạy: `mapID=1, tmw=52, tmh=24, tileID=1, bgType=0` (⇒ `typeBg=0`), **`zoom=2`** (client nạp từ `res/x2/…`).
- Server Go `banDo.go:33`: map 1 = `TileID:1, BgID:0, PlanetID:0`, **tên đã "Làng Lá"** (spawn 100,384; NPC Iruka; quái Sói). Chỉ thiếu đồ hoạ.
- `res/x{1..4}/bg/` chỉ có `fog0.png`, `light.png` — **không có `b0k.png`** ⇒ nền hiện đến từ RMS/skyColor.
- Collision server = lưới **5×17 = footprint 120×408** (`map-1-collision.bin`), nhưng `mymap/1.bytes` là **map rỗng 52×24** ⇒ **đã lệch sẵn** (cần thống nhất, xem bẫy #5).
- Kho art Làng Lá sẵn (x1–x4): 26 mảnh `map_la_*` (cỏ/đất/cổng/hàng rào/mái ngói/tường/đèn lồng…), 10 mảnh `map_hv_*` (học viện), nền `map_nen_nuihokage/troi/may`.

## 3. Phạm vi — chia 4 lát, mỗi lát TỰ SOI được trong Unity

**Nguyên tắc (bài học logo):** mỗi lát làm xong phải **tự vào Unity chụp màn hình kiểm bằng mắt**, không tin report. Không để game ở trạng thái xấu giữa các lát.

| Lát | Việc | Kết quả nhìn thấy | Đụng |
|---|---|---|---|
| **1. Nền Làng Lá** | Dựng 4 lớp `bg/b0k.png` đúng cách (canh màu trời + vị trí + đủ mức zoom) | Nền sau lưng = trời/mây/**núi Hokage** Naruto | Client Resources |
| **2. Mặt đất** | Xác minh khuôn `/t/{tileID}/`, cắt `map_la_co/dat` thành bộ gạch, vẽ `mymap/1.bytes` phủ cỏ/đất | Đi trên đất/cỏ Làng Lá (thay cỏ DB) | Client + có thể `banDo.go` TileID |
| **3. Nhà cửa** | Thêm mảnh trang trí (cổng làng, hàng rào, mái ngói, đèn lồng) vào bố cục | Làng có kiến trúc Naruto | Client (bố cục + tile) |
| **4. Va chạm** | Thống nhất footprint `mymap ↔ collision ↔ biên di chuyển`, regen `asset-lock.json` | Đi lại đúng, không rơi/kẹt | Server Go + assets |

**Bản 1 = Lát 1 + 2** (nền + mặt đất) — đủ để game "ra dáng Làng Lá". Lát 3–4 làm tiếp sau khi bản 1 chạy ngon.

## 4. Quyết định kỹ thuật

1. **Bắt đầu bằng NỀN (Lát 1)** — rủi ro thấp nhất (chỉ client Resources, không đụng server/collision), thấy ngay, kiểm chứng trọn đường "Resources thắng cache".
2. **Đặt ảnh đủ 4 mức zoom** `x1..x4` (client hiện zoom=2 nhưng người chơi đổi được) — kho art có sẵn cả 4 mức.
3. **Giữ `tileID`/`bgType` map 1 nếu có thể** (dùng lại số hiện có để tránh đụng `tile_set_info.bin` + regen lock). Chỉ đổi `banDo.go` khi bắt buộc, và giữ `tileId ∈ [1,39]`.
4. **Không đụng assetpack/part.bin** cho map (map không qua đó). Nếu Lát 4 đổi collision → **phải regen `asset-lock.json`** kèm.
5. **Thực nghiệm nhỏ trước mỗi lát có ẩn số**: Lát 2 phải xác minh khuôn `/t/{id}/` thật (mấy ảnh/ô, tên file) bằng một tile thử trước khi cắt cả bộ.

## 5. Chi tiết Lát 1 — Nền Làng Lá (làm trước)

Nền `typeBg==0` (`GameCanvas.cs:1185`) vẽ 4 lớp qua `paintBackgroundtLayer(g, layer, …)` dùng `imgBG[0..3]` = `bg/b00..b03`. **Bẫy đã kiểm chứng:** hàm này **lấy màu trời từ pixel top/bottom của ảnh nền** (đặt núi Hokage tối màu vào b03 → trời thành ĐEN). Nên:
- Lớp xa `b03`: dải trời/mây sáng (dùng `map_nen_troi` + `map_nen_may` ghép ngang), **mép trên sáng** để màu trời đúng.
- Lớp giữa: dãy núi Hokage (`map_nen_nuihokage`) đặt thấp, nền trong suốt để không phủ trời.
- Ghép/чải các mảnh 96×96 (x4) thành ảnh nền đủ rộng (cuộn ngang) + xuất 4 mức zoom.
- Đặt vào `res/x{1..4}/bg/b0{k}.png`, reimport, `stop→play`, chụp soi. Lặp tinh chỉnh tới khi nền đẹp + trời đúng màu.

## 6. Bẫy — đọc trước khi code (từ điều tra + thực nghiệm)

1. **Màu trời = pixel ảnh nền** → ảnh nền phải sáng phần trời, hoặc chấp nhận trời theo màu ảnh. (Đã vấp: Hokage tối → trời đen.)
2. **Đúng mức zoom**: client hiện `zoom=2` → phải có `res/x2/…`; đặt thiếu x2 mà chỉ x4 thì không hiện. Đặt đủ x1–x4.
3. **`tileId ∈ [1,39]`** (số tileset trong `tile_set_info.bin`); ngoài khoảng → **client `GameMidlet.exit()`** (crash cứng, `thegioi.go:14`). Sai thứ tự gói `-82`→`-24` cũng thoát app.
4. **Regen `asset-lock.json`** nếu đổi `collision/catalog/*.bin` — lệch hash thì **server panic fail-closed** lúc khởi động (`catalog.go:213`).
5. **Footprint phải khớp 3 nguồn**: `mymap/1.bytes` (bố cục) ↔ `map-1-collision.bin` (5×17=120×408) ↔ biên di chuyển. Hiện `mymap` là 52×24 → **đang lệch**; Lát 2/4 phải chọn một footprint và đồng bộ cả ba, nếu không "người chơi rơi khỏi map ngay từ spawn" (BAT-DAU §6.24).
6. **Khuôn `/t/{tileID}/` chưa reverse đủ** (BAT-DAU §8c): chưa rõ mỗi tileset gồm mấy ảnh/đặt tên sao. **Lát 2 phải thực nghiệm xác minh khuôn** trước khi cắt cả bộ.
7. **RMS cache** có thể còn tàn dư che Resources ở vài đường — nếu đặt Resources mà không hiện, kiểm tra/ xoá RMS tương ứng.
8. **Client Unity KHÔNG nằm trong git** repo này → sửa không version được; cần một **script deploy** chép art vào Unity để lặp lại trên máy khác (BAT-DAU §8b). Người chơi thật **phải build lại client** mới thấy.

## 7. Rủi ro chính

- **Khuôn tile `/t/` là ẩn số lớn nhất** → Lát 2 có thể tốn vòng thực nghiệm; giảm rủi ro bằng cách thử 1 tile trước.
- **Footprint lệch sẵn** → phải quyết một kích thước map Làng Lá thật và đồng bộ; đây là gốc nhiều lỗi "rơi khỏi map".
- **Nền parallax nhiều tham số** (màu, vị trí Y, tiling, tốc độ lớp) → Lát 1 cần vài vòng canh chỉnh.

## 8. Thứ tự thi công + kiểm thử

1. **Lát 1 — Nền:** ghép nền Làng Lá (script Python, dùng `map_nen_*`) → 4 mức zoom → đặt `bg/b0k` → **chụp soi Unity** (kiểm: trời đúng màu, núi Hokage hiện, không đen). Lặp tới đẹp.
2. **Lát 2a — Thực nghiệm khuôn tile:** đặt 1 tile thử vào `res/x2/t/{tileID}/` → vào map → chụp: có thay ô đất không? Xác định khuôn (số ảnh, tên).
3. **Lát 2b — Mặt đất:** cắt `map_la_co/dat` theo khuôn → vẽ `mymap/1.bytes` phủ nền → **chụp soi**. Nếu cần `tileID` mới → sửa `banDo.go` (giữ [1,39]) + **test Go** (`banDo`/`thegioi`).
4. **Lát 3 — Nhà cửa:** thêm mảnh trang trí vào bố cục → **chụp soi**.
5. **Lát 4 — Va chạm:** đồng bộ footprint (mymap ↔ collision ↔ biên) + regen `asset-lock.json` → **test Go** (collision, di chuyển) + đi thử trong client (không rơi/kẹt).

**Kiểm thử:** phần server Go (Lát 2b/4) làm **TDD** như các KH khác (`go test`); phần client là **thực nghiệm-trong-Unity có chụp soi** (client không test tự động được). Mỗi lát cập nhật `BAT-DAU-TU-DAY.md` + `NHAT-KY-AGENT.md`.

## 9. Ngoài phạm vi (YAGNI)

- Các map khác (2–6) — làm sau khi map 1 xong khuôn.
- Map Học Viện trong nhà (`map_hv_*`) — lát riêng sau Làng Lá ngoài trời.
- Đấu nối art vào assetpack/part.bin — không cần cho map (map dùng client Resources trực tiếp).
- Script deploy tự động + đưa client vào git — ghi nợ, làm khi ổn định khuôn.
