# Kế hoạch đưa map Làng Lá tới mức hoàn thiện

Ngày lập: 22/08/2026 · Điểm khởi đầu: **45/100** (đo bằng đợt soi 11 agent, 42 lỗi thật)

Cách làm: **sửa → cho agent khác soi → sửa tiếp**. Mỗi vòng kết thúc bằng một lượt soi
độc lập chấm điểm lại. Không tự chấm điểm cho mình.

> Danh sách 42 lỗi đầy đủ (có toạ độ pixel + bằng chứng): `docs/my-thuat/map-1-lang-la/SOI-LAN-1.md`
> Lỗi cũ và cách sửa: `docs/my-thuat/map-1-lang-la/LOI-DA-GAP.md`

---

## Vòng 1 — Chặn đường: cho người chơi *thấy mình* (không cần vẽ gì)

Ba lỗi này khiến map không thể cho ai xem. Sửa hết mới đi tiếp.

| # | Việc | Trạng thái |
|---|---|---|
| 1.1 | Khai loại va chạm cho tile 60–69 trong `tile_set_info.bin` (bit 2 sàn / 4 phải / 8 trái / 0x2000 trần) | ✅ xong, 4 test + 2 đột biến |
| 1.2 | `map-1-collision.bin` — **hoá ra không lệch**: client dừng nhân vật ở *đỉnh* ô sàn (`cy+1` mới là ô đặc), server kiểm chính `cy` nên hàng 18 rỗng là đúng. Đo thật: `cy=432`, hợp lệ cả hai phía | ✅ không phải sửa |
| 1.3 | Gỡ **cột tre lơ lửng** — client vẽ cứng `caycot.png` cho `mapID == 1` (di sản Ngọc Rồng: map 1 cũ là rừng tre) | ✅ xong |
| 1.4 | **Nhát chém y=344** (chân núi): làm chân `b02` tan dần vào trời từ 68%→84% chiều cao. Bước nhảy 36,9 → không còn trong top | ✅ xong |
| 1.4b | **Nhát chém y=561** (đáy làng xa): làm mềm đáy `b01` cho khớp đỉnh `b00`. Chỗ nối bị nhà che nên khó thấy | ✅ xong |
| 1.5 | Nhân vật **vô hình** dù đã đứng đúng chỗ: `mGraphics.drawRegion(Small,…)` chỉ có `throw new NotImplementedException()` → `paintCharBody` chết ở mảnh đầu tiên, chỉ còn cái bóng | ✅ xong |
| 1.6 | Nhân vật **lún nửa người**: `dy` của 34 frame part 30/57/58 từ 0 → −22; kèm `vsData` 21→22 và cắm 42 sprite ArtRoot chính thức vào client | ✅ xong |
| 1.7 | Tắt dòng gỡ rối `X: … Y: …` trên HUD (thêm công tắc `ModFunc.hienToaDoGoRoi`) | ✅ xong |
| 1.8 | **Thanh máu và chakra hiện 0** — server **chưa bao giờ gửi gói `-42`** (chỉ số nhân vật: HP/MP gốc, HP/MP đầy, HP/MP hiện tại, tốc độ, sát thương, thủ, chí mạng). Client không có nguồn nào khác để biết máu | ⬜ **thuộc KH6 Task 15**, xem ghi chú dưới |

**Chốt vòng 1:** nhân vật hiện ra và đi lại được → cho agent soi lần 2. ✅ đã soi.

> **Ghi chú việc 1.8.** Lộ trình ghi KH6 (chiến đấu) đã xong, nhưng còn **Task 15 — nghiệm thu tay
> bằng client Unity**, và chính lỗi này là thứ Task 15 đáng lẽ bắt được: đánh nhau mà hai bên đều
> 0 máu thì không kiểm được gì. Cấu trúc gói `-42` đọc được từ `Controller.cs:1704`; phần trong
> `try/catch` là tuỳ chọn (client tự đặt 0 nếu thiếu), nên bản ngắn tới `cCriticalGoc` là đủ dùng.
> Để riêng vì đây là việc của KH6, không phải của map — và đang có việc Boss dở trên cùng nhánh.

---

## Vòng 2 — Dọn chỗ đặt sai (không cần vẽ gì)

| # | Việc |
|---|---|
| 2.1 | Hàng rào cắm xuyên cột cổng torii — dời sang cột trống |
| 2.2 | Góc phải nhà lớn treo lơ lửng ngoài mép gò (nhà cột 31–38, gò chỉ tới 37) — nới gò hoặc dời nhà |
| 2.3 | Nhà nhỏ bị gò nuốt mất nửa phải; bụi cỏ cạnh đó bị xẻ dọc |
| 2.4 | Giếng bị vách gò xén mất một chân; bỏ bớt 1 trong 2 giếng gần giống nhau |
| 2.5 | Cây bị mép map cắt mất 1/3 tán — dời vào trong |
| 2.6 | Biển hiệu đang cắm bên trong quán ramen (quán vốn đã có 4 biển) — dời ra |
| 2.7 | Dọn viền rác tile 65/66 (cột sọc xám như mã vạch, lộ ở 4 chỗ) |
| 2.8 | Tắt dòng debug `X: … Y: …` trên HUD; sửa gói gửi HP/chakra đang hiện 0 |

**Chốt vòng 2:** không vật thể nào bị cắt / treo / chồng → soi lần 3.

---

## Vòng 3 — Đặt hàng ảnh mới (phải vẽ)

Viết đề bài kèm ảnh tham chiếu, giao cho agent vẽ. Xếp theo mức đáng tiền:

| Ưu tiên | Món | Chữa được lỗi nào |
|---|---|---|
| 1 | **Bộ tile vách gò** (thân + chân) khớp tông đáy ô cỏ | vách 5 chất liệu · tile 65/66 hỏng · đường cắt ngang y=626 |
| 2 | **3–4 biến thể tile đất sâu** + lớp tối dần theo độ sâu | gần 1/4 màn hình là đất chết kẻ ô bàn cờ |
| 3 | **Vẽ lại b02 núi Hokage** — bản 768px thật, to gấp 3–4 lần | núi lặp 2 lần/màn · đường nối dọc |
| 4 | **Sprite bóng đổ elip mờ** cho mọi công trình | nhà/giếng/cây như sticker dán |
| 5 | **Kè đá + bậc thang đá** dưới torii | concept đã duyệt có, game không có |
| 6 | **2 ô chuyển tiếp cỏ ↔ lát đá** + túm cỏ, đá lẻ rải mép | nhát cắt dọc giữa cỏ và lối đá |
| 7 | **Vẽ lại b01 làng xa** cao hơn, đáy tắt dần | mái nhà treo lơ lửng ở y=561 |
| 8 | Thêm **1–2 kiểu cây** cho khoảng giữa map | cả map chỉ có 1 cây, dính sát biên phải |

---

## Vòng 4 — Lắp ảnh mới, soi lại

Lắp theo thứ tự ưu tiên trên, mỗi 2–3 món chụp lại và soi một lượt.

---

## Vòng 5 — Lối chơi: map sống

| # | Việc |
|---|---|
| 5.1 | Waypoint rời map (hiện **không có cái nào** — vào rồi không ra được) |
| 5.2 | NPC đứng ở quán ramen, cổng torii, giếng |
| 5.3 | Quái cấp thấp rải ở khoảng trống giữa map |
| 5.4 | Vật phẩm rơi, rương |
| 5.5 | Điểm xuất hiện đặt đúng chỗ đẹp (hiện đang ở đáy map) |

---

## Cách đo tiến bộ

Mỗi vòng chạy lại đúng đợt soi cũ (5 lăng kính + phản biện + chấm điểm) trên ảnh mới.
Ghi điểm vào bảng dưới. Mục tiêu: **≥ 85/100** rồi mới tính là xong.

| Vòng | Ngày | Điểm | Lỗi chặn đường còn |
|---|---|---|---|
| 0 (khởi đầu) | 22/08/2026 | 45 | 3 |
