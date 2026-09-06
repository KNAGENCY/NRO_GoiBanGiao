# MAP — Đợt 4 (P1 Làng Lá phần còn thiếu + Khu rừng · P2 bốn làng — cần vẽ)

> Đọc `_KHUON-DOT-4.md` và `../dot-3-toan-bo/_KHUON.md` trước. **LƯU VÀO:** `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\`
>
> File cặp: `05-map.json` (manifest máy đọc, `anh[]` = đúng danh sách ở mục 5, `prompt_day_du` = đúng khối ở mục 6). Viết 06/09/2026 từ: `00-KIEM-KE-DA-CO.md` §3.5, `00-KIEM-KE-GAME-CAN.md` P1-09/P1-10/P2-03/P2-04, `_KHUON-DOT-4.md` B.5, `../KHUON-DE-BAI-ANH-MAP.md`, `../SO-TAY-LAM-MAP.md`, `../map-1-lang-la/{SOI-LAN-3, DOT-2-ANH-CAN-VE, LOI-DA-GAP}.md`, `../../superpowers/plans/2026-08-23-lo-trinh-4-map-con-lai.md`. Chỗ nào là **đề xuất** hay **chưa chắc** đều ghi thẳng.

**Tổng: 165 ảnh cần vẽ** — P1 **34** (Làng Lá 6 · Khu rừng 28) · P2 **131** (Làng Lá 2 đề xuất · Khu rừng 2 tiền cảnh · Cát 31 · Sương Mù 33 · Đá 31 · Mây 32) · P3 0 — cộng 5 **concept** (tham chiếu, không lắp). Đã có sẵn và dùng được: **51 ảnh map 1** (mục 2).

---

## 1. Game dùng ở đâu, vì sao cần

Một **map** trong game này là bốn thứ xếp chồng: (a) **4 lớp nền** phía sau cuộn chậm theo camera (trời → núi → làng xa → lớp gần), (b) **lưới ô 96 px** làm mặt đất, lòng đất, vách, bậc thang — nhân vật đứng và va chạm trên lưới này, (c) **vật thể to** cắm trên mặt đất (nhà, cây, cổng, giếng…) và (d) **tiền cảnh** sẫm che rìa. Máy chủ chỉ biết lưới va chạm; mọi thứ nhìn thấy là ảnh — vẽ sai cỡ là engine cắt mất một phần, vẽ hẹp là nền lặp lộ liễu, tile trơn là một phần tư màn hình thành mảng nâu chết (đã vấp cả ba).

Tình trạng hôm nay (đọc từ code, không đoán):

| Map | Máy chủ | Mỹ thuật | Còn thiếu / lỗi | Ưu tiên |
|---|---|---|---|---|
| **1 Làng Lá** | vào được, có cổng sang map 2 (`banDo.go:84,115-135`) | 51 ảnh hand-painted đã cắm, **6/7 chỉ số đạt** (`do-map.py`) | **Đất phẳng lì 91 % (đích ≤ 40 %)** — đã thử 3 cách bằng code, cả 3 hỏng → phải vẽ 4 tile đất có vân · 2 nền `_rong` còn **1536** nên núi lặp 1,9 lần/khung, phải **2976** · SOI-LAN-3: núi thấp hơn mái nhà (#16), lối lát như rãnh sỏi (#19), vật hở trời 3 px dưới chân (#20) | **P1** |
| **6 Khu rừng** | `catalog.json maps[5]`, sói sinh 2 con (`quai.yaml:118-121`) **nhưng `banDo.go:81-89` chỉ khai 5 map** → chưa vào được, chưa có lưới va chạm | **chưa có gì** (25 tile `map_rung_*` pixel-art trong lô 343 không dùng vì map đã chốt hand-painted) | cả bộ theo khuôn, **không nhà** — rừng rậm, cây cổ thụ, lối đá rêu | **P1** (⚠️ chỉ có nghĩa sau khi máy chủ nối map 6 — C8 của KIEM-KE-GAME-CAN) |
| **2 Cát · 3 Sương Mù · 4 Đá · 5 Mây** | vào được nhưng map 3–5 là **ngõ cụt** (`BAT-DAU-TU-DAY.md` §5 1a′) | **tile Dragon Ball** (`banDo.go:85-88` TileID 1/2/3, BgID 0/1/2) | cả bộ theo khuôn; thứ tự lộ trình **Cát → Sương Mù → Đá → Mây** (Cát khác Làng Lá nhất — thử khuôn được) | **P2** |

Vì sao vẽ map trước cả quái/boss của 4 làng: người chơi nhìn thấy map **trong một giây đầu**; và map là thứ duy nhất trong đợt này có **thước đo bằng máy** (`do-map.py` 7 chỉ số) nên vẽ xong biết ngay đạt hay chưa, không phải cãi nhau bằng mắt.

---

## 2. Đã có gì

### 2.1 Map 1 Làng Lá — 51 ảnh trong `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-da-tao\` (đo thật 05/09/2026)

| Tên | Cỡ | Nhóm | Trạng thái lắp | Dùng được không · vì sao |
|---|---|---|---|---|
| `nen_troi_x4.png` | 768×768 | nền | LẮP (client b03/b13) | trời + mây, nền đặc; cỡ 768 chấp nhận vì trời lặp không lộ |
| `nen_nuihokage_x4.png` | 1536×768 | nền | bản đợt 1, đã thay | bản đầu, thay bằng _rong |
| `nen_nuihokage_rong_x4.png` | 1536×768 | nền | LẮP (client b02/b12) — SAI CỠ | còn 1536 → núi lặp 1,9 lần/khung → vẽ lại 2976 (nen_nuihokage_rong2) |
| `nen_lang_xa_x4.png` | 1536×576 | nền | bản đợt 1, đã thay | bản đầu, thay bằng _rong |
| `nen_lang_xa_rong_x4.png` | 1536×576 | nền | LẮP (client b01/b11) — SAI CỠ | còn 1536 → vẽ lại 2976 (nen_lang_xa_rong2) |
| `nha_nho_x4.png` | 576×576 | công trình | LẮP (tile 70) | dùng được |
| `nha_vua_x4.png` | 768×576 | công trình | LẮP (tile 71) | dùng được |
| `nha_lon_x4.png` | 768×768 | công trình | LẮP (tile 72) | dùng được |
| `nha_quan_x4.png` | 768×576 | công trình | LẮP (tile 73) | dùng được (= nha_dacsan của khuôn) |
| `tile_dat_1_x4.png` | 96×96 | tile | LẮP (đợt 1) | đất trơn — nguồn của lỗi 'đất phẳng lì 91%' |
| `tile_dat_2_x4.png` | 96×96 | tile | LẮP (đợt 1) | đất biến thể, vẫn trơn |
| `tile_co_mat_x4.png` | 96×96 | tile | LẮP (đợt 1) | mặt cỏ — tham chiếu màu cho mọi tile mới |
| `tile_co_goc_trai_x4.png` | 96×96 | tile | LẮP (đợt 1) | dùng được |
| `tile_co_goc_phai_x4.png` | 96×96 | tile | LẮP (đợt 1) | dùng được |
| `tile_vach_trai_x4.png` | 96×96 | tile | LẮP (đợt 1) | dùng được |
| `tile_vach_phai_x4.png` | 96×96 | tile | LẮP (đợt 1) | dùng được |
| `tile_doc_len_x4.png` | 96×96 | tile | LẮP (đợt 1) — né dùng | dốc chéo không khớp bậc thẳng (LOI-DA-GAP) |
| `tile_doc_xuong_x4.png` | 96×96 | tile | LẮP (đợt 1) — né dùng | như trên |
| `tile_datlat_x4.png` | 96×96 | tile | LẮP (đợt 1) | cả ô là sỏi → đọc thành rãnh (SOI-LAN-3 #19); đề xuất tile_datlat_mat |
| `tile_dat_lienmach_1_x4.png` | 96×96 | tile | LẮP (tile 60) | lặp rất liền (lệch 3/96) nhưng trơn — giữ làm tham chiếu MÀU cho 4 tile vân |
| `tile_dat_lienmach_2_x4.png` | 96×96 | tile | LẮP (tile 61) | như trên |
| `tile_dat_tang_giua_x4.png` | 96×96 | tile | LẮP (tile 67) | tối hơn 8%, trơn |
| `tile_dat_tang_day_x4.png` | 96×96 | tile | LẮP (tile 68) | tối hơn 16%, trơn |
| `tile_vach_go_trai_x4.png` | 96×96 | tile | LẮP (tile 65) | thay hai tile hỏng 65/66 cũ |
| `tile_vach_go_phai_x4.png` | 96×96 | tile | LẮP (tile 66) | như trên |
| `tile_chan_go_x4.png` | 96×96 | tile | LẮP (tile 83) | dùng được |
| `tile_ke_da_x4.png` | 96×96 | tile | LẮP (tile 84) | dùng được |
| `tile_bac_thang_x4.png` | 96×96 | tile | LẮP (tile 85) | dùng được |
| `tile_bac_thang_dinh_x4.png` | 96×96 | tile | LẮP (tile 86) | dùng được |
| `tile_co_sang_da_x4.png` | 96×96 | tile | LẮP (tile 63) | dùng được |
| `tile_da_sang_co_x4.png` | 96×96 | tile | LẮP (tile 64) | dùng được |
| `tt_cong_torii_x4.png` | 480×480 | công trình | LẮP (tile 74) | chuẩn vàng vật thể rời |
| `tt_hangrao_x4.png` | 192×192 | vật nhỏ | LẮP (tile 75) | dùng được |
| `tt_denlong_x4.png` | 192×384 | vật nhỏ | LẮP (tile 76) | dùng được |
| `tt_gieng_x4.png` | 288×288 | vật nhỏ | LẮP (tile 77) | dùng được |
| `tt_thunggo_x4.png` | 192×192 | vật nhỏ | LẮP (tile 78) | dùng được |
| `tt_bienhieu_x4.png` | 192×288 | vật nhỏ | LẮP (tile 79) | dùng được |
| `tt_cay_x4.png` | 576×768 | cây | LẮP (tile 80) | dùng được — tham chiếu mức chi tiết cây |
| `tt_cay_nho_x4.png` | 384×384 | cây | LẮP (tile 87) | dùng được |
| `tt_bui_co_x4.png` | 192×96 | vật nhỏ | LẮP (tile 81) | dùng được |
| `tt_bui_ram_x4.png` | 288×192 | cây | LẮP (tile 88) | dùng được |
| `bong_cay_x4.png` | 576×96 | bóng đổ | chưa chắc đã lắp | BAT-DAU-TU-DAY §5 1a ghi 'xong'; KIEM-KE-DA-CO không thấy trong x1/t; mapBackGround/500–507 (8 file) có thể là 8 bóng — chưa xác minh |
| `bong_den_x4.png` | 192×96 | bóng đổ | chưa chắc đã lắp | như trên |
| `bong_gieng_x4.png` | 288×96 | bóng đổ | chưa chắc đã lắp | như trên |
| `bong_nha_lon_x4.png` | 768×96 | bóng đổ | chưa chắc đã lắp | như trên |
| `bong_nha_nho_x4.png` | 576×96 | bóng đổ | chưa chắc đã lắp | như trên |
| `bong_quan_x4.png` | 768×96 | bóng đổ | chưa chắc đã lắp | như trên |
| `bong_thung_x4.png` | 192×96 | bóng đổ | chưa chắc đã lắp | như trên |
| `bong_torii_x4.png` | 480×96 | bóng đổ | chưa chắc đã lắp | như trên |
| `tc_canh_la_x4.png` | 768×384 | tiền cảnh | chưa xác minh | tiền cảnh cành lá — chưa rõ đường lắp |
| `tc_bui_truoc_x4.png` | 384×192 | tiền cảnh | chưa xác minh | tiền cảnh bụi cỏ — chưa rõ đường lắp |

Kết luận 2.1: **49/51 dùng được nguyên**, chỉ 2 nền `_rong` phải vẽ lại đúng cỡ (`ve_lai`), và 4 tile đất mới phải vẽ thêm. Bản gốc trước sửa viền: `anh-da-tao\goc-truoc-khi-sua\` (8 file).

### 2.2 Những thứ khác liên quan map

| Thứ | Nơi lưu | Dùng được không | Vì sao |
|---|---|---|---|
| Concept Làng Lá đã duyệt | `map-1-lang-la\concept-map1-langla.png` 2048×1024 | **Có — là tham chiếu phong cách cho MỌI map** | tranh vẽ mịn, không nhân vật; mọi prompt dưới đây đều bảo bám theo nó |
| 79 tile pixel-art `map_*` lô 343 (`map_la_*` 26, `map_rung_*` 25, sông/cầu/biển/sương/học viện/núi 28) | `anh-nguon\343-anh\09-ban-do\` 96×96 (+ cổng 192, núi 384×192) | **Không lắp** | pixel-art ≠ hand-painted đã chốt (DOT-2 22/08; SOI-LAN-3 #5 'hai game khác nhau'); `map_nen_*` 384×192 sai cỡ nền. Chỉ dùng tra **tên/bố cục** |
| 56 tên `map_{cat,smu,da,may}_*` 96×96 + 8 `map_nen_*` 384×192 (đợt 3) | `kich-thuoc-anh.json` dòng 1598–1854 | **Chưa có ảnh; KHÔNG vẽ theo cỡ đó nữa** | khuôn cũ 16 ảnh/làng, nền 384×192 sai. Tên giữ trong bảng (append-only), đợt 4 dùng bộ tên mới ở mục 5 |
| Brief đợt 3 `lang-*/map.md` (16 ảnh/làng) | `dot-3-toan-bo\lang-*\map.md` | **Thay bằng file này** | KIEM-KE-GAME-CAN P2-04 đã ghi là khuôn cũ sai cỡ |
| Tile + nền đang chạy trong client | `D:\NRO_Project\Assets\Resources\res\x{1..4}\t\1$40–88`, `bg\b00–b03` = `b10–b13`; sao lưu `lang-la-client-assets\` | Bản chạy thật của map 1 | `x1\mapBackGround\500–507` (8 file, 27/08) **chưa rõ là gì** — có thể là 8 bóng đổ (chưa xác minh) |
| Công cụ | `tools\map-lang-la\{do-map, kiem-anh-dot-2, lap-anh-dot-2, sinh-va-cham}.py`, `tiled\` | Có, nhưng **ghi cứng map 1** | `kiem-anh-dot-2.py:11-29` và `lap-anh-dot-2.py:13-44` phải sửa `THU_MUC`/`CHO_DOI`/bảng slot cho map mới — xem mục 7 |
| Thước tỉ lệ map | `map-1-lang-la\anh-tham-chieu\THUOC-TI-LE.png` | Có — gắn kèm mọi lần gen vật thể | ghi nhân vật 2,4 ô (số cũ); số đo mới 246 px ≈ 2,56 ô — chênh không đáng kể cho nhà/cây |

---

## 3. Ràng buộc kỹ thuật riêng loại MAP

Brief này **không có nhân vật, NPC, quái, boss** — mọi ảnh là **tĩnh** (`kieu: tinh_alpha`), đi đường `KIEM-TRA-ANH.py → TAO-4-MUC-ZOOM.py` (hoặc `kiem-anh-dot-2.py` cho map 1). Vì thế 7 luật bài học 05/09 (cỡ đầu, chân chạm ô, dải chroma…) **không áp ở đây**; nếu hoạ sĩ vẽ người vào nền hay đặt vật lên nền tím là sai đề. Riêng luật *chân vật thể chạm mép dưới* thì vẫn đúng y như nhân vật — vì cùng một lý do: engine neo theo mép ảnh.

### 3.1 Hai kiểu file trong một brief

| | Nền `nen_*` (4 lớp/map) | Tile · công trình · cây · vật nhỏ · tiền cảnh |
|---|---|---|
| Nền ảnh | **ĐẶC** (opaque). Giao PNG RGB, hoặc RGBA với alpha = 255 toàn ảnh | **Trong suốt thật** (PNG RGBA); tile mặt đất có phần trên là 'không khí' trong suốt |
| Máy làm gì | `lap-anh-dot-2.py:77` `convert('RGB')` → `bg/b{typeBg}{k}.png` 4 mức zoom, chép cả bản `b1x` | thu 4 mức (premultiply) → slot tile `t/<tileID>$k` hoặc `mapBackGround/<id>` |
| `KIEM-TRA-ANH.py` | sẽ báo **`THIEU KENH TRONG SUOT`** cho file RGB (`KIEM-TRA-ANH.py:117`) — với `nen_*` **đó không phải lỗi**; muốn công cụ im thì giao RGBA alpha 255 | kiểm đủ 7 điều bình thường |
| Client đọc | `GameCanvas.cs:1694` `"/bg/b" + typeBg + k + ".png"`, `typeBg` = `BgID` của map (`Controller.cs:4983 loadBG(TileMap.bgID)`); map 1 chạy `typeBg = 1` → `b10..b13` (LOI-DA-GAP §4). Mỗi `typeBg` có một `case` riêng trong `GameCanvas.cs:1541+` (số lớp, tốc độ cuộn) → **số `typeBg` cho map 2–6 phía dự án chốt**, hoạ sĩ không cần quan tâm | tile: `TileMap.cs:327` `"/t/" + tileID + "$" + k`, 100 ô/bộ (`build.py:1052-1056`); map 1 dùng bộ 1 slot 40–88 → **map mới cần bộ tile riêng (`TileID` 2–6) hay slot 89–100 — chưa chốt** |
| Bật trong Unity | **Read/Write Enabled** (không thì `getRGB` = 0 → trời đen), `wrapMode Clamp`, `npotScale None`, mipmap off (LOI-DA-GAP §3) | như bên, `npotScale None` là lý do cỡ phải bội số 96 |

### 3.2 Bảng cỡ — đúng từng pixel, đều là bội số 96 (mức x4)

| Nhóm | Cỡ | Vì sao con số đó |
|---|---|---|
| `nen_troi` | **2976 × 384** | Ảnh giao x4, game hiển thị **x2** → bề rộng chia đôi; màn rộng 1454 px → ảnh phải ≥ 2908 x4, làm tròn bội 96 = **2976**. Bản 1536 chỉ còn 768 trên màn → lặp 1,9 lần (đã vấp, cả núi lẫn làng xa) |
| `nen_nui` | **2976 × 768** | khối núi cần cao để 'khối núi ≥ 30 % khung' (`do-map.py`); SOI-LAN-3 #16: núi hiện thấp hơn mái nhà |
| `nen_langxa` | **2976 × 576** | lớp che chân núi |
| `nen_gan` | **2976 × 288** | lớp thấp nhất, ngay sau mặt cỏ |
| tile | **96 × 96** | 1 ô = 24 đơn vị game × 4 |
| nhà nhỏ · nhà vừa · nhà lớn · nhà đặc sản | 576×576 · 768×576 · 768×768 · 768×576 | nhà phải cao **2,5–3 lần nhân vật** (246 px) — 6–8 ô; vẽ nhà 2 ô là nhân vật cao hơn nhà (đã vấp bản 1 README map 1) |
| cổng làng | 480 × 480 | 5 ô |
| cây lớn · cây nhỏ · bụi | 768×768 · 384×384 · 288×192 | 8 · 4 · 2 ô |
| vật nhỏ | 288×288 · 192×384 · 192×192 · 192×96 · 288×192 · 576×192 · 480×192 · 768×288 · 768×384 | tất cả bội 96; ghi cụ thể từng dòng ở mục 5 |

### 3.3 Bố cục dọc 4 lớp nền — chỗ sinh lỗi nhiều nhất

| Lớp | 0 → 30 % chiều cao | 30 → 80 % | 80 → 100 % |
|---|---|---|---|
| `nen_troi` | trời + mây | trời nhạt dần | **màu nối** thuần |
| `nen_nui` | trời + mây (chỉ tới 25 %) | **khối núi đặc, chiếm ≥ 55 % chiều cao ảnh** | nhoà dần về **màu nối** |
| `nen_langxa` | **mái nhà / tán cây NGAY TỪ HÀNG ĐẦU — không được có trời** | làng dày | nhoà dần về **màu nối** |
| `nen_gan` | hình gần mờ | sương dày | nhoà về **màu nối** |

Hai luật tuyệt đối (KHUON-DE-BAI-ANH-MAP): (1) `nen_langxa` **không có trời ở đỉnh** — có là nó nằm thấp hơn chân núi, che không nổi, giữa núi và nhà lộ một khoảng trời rỗng (map 1 từng phải cắt 36 % trên của ảnh làng xa); (2) **hàng pixel cuối cùng của mỗi nền = đúng mã màu nối**, không còn vân — engine điền phần dưới mỗi lớp bằng đúng màu hàng đáy (`LOI-DA-GAP.md` §8), đáy còn vân là một vạch ngang chạy suốt màn (SOI-LAN-3 #7, #25).

**Màu nối từng map** (một mã cho cả 4 lớp của map đó):

| Map | Màu nối | Trạng thái |
|---|---|---|
| Làng Lá (map 1) | `RGB(10, 86, 154)` | đã chốt — đo từ đáy b02 đang chạy (LOI-DA-GAP.md §5) |
| Khu rừng (map 6) | `RGB(52, 88, 72)` | đề xuất (sương rừng xanh xám) — chốt lại sau khi có concept |
| Làng Cát (map 2) | `RGB(226, 190, 120)` | đề xuất (trời vàng cát bụi) — chốt lại sau khi có concept |
| Làng Sương Mù (map 3) | `RGB(150, 172, 184)` | đề xuất (sương xám lam) — chốt lại sau khi có concept |
| Làng Đá (map 4) | `RGB(142, 140, 132)` | đề xuất (trời xám bụi đá) — chốt lại sau khi có concept |
| Làng Mây (map 5) | `RGB(120, 110, 160)` | đề xuất (tím lam nhạt) — chốt lại sau khi có concept |

> Với 5 map mới, màu nối là **đề xuất** để hoạ sĩ có số mà vẽ; sau khi concept (mục 6.0) được duyệt, phía dự án lấy màu trời chủ đạo của concept làm màu nối chính thức và **sửa số trong prompt trước khi gen bộ**. `SO-TAY-LAM-MAP.md` §5 nói kỹ hơn: nếu về sau muốn mỗi cặp lớp một màu riêng thì máy nhoà lại được (`lam-mem-cho-giap-nen.py`), nhưng ảnh gốc vẫn phải nhoà về một màu.

### 3.4 Tile — luật đo được

- **Lặp được 4 phía**: `kiem-anh-dot-2.py` đo chênh mép trái/phải và trên/dưới; bản đợt 2 đạt 3/96 — giữ mức đó.
- **Bốn tile lòng đất (`dat_1`, `dat_2`, `dat_giua`, `dat_day`) khác nhau ở VÂN, gần nhau về ĐỘ SÁNG** (chênh < 8 %; `giua` tối hơn ~6 %, `day` ~12 %). Chênh 48 % thì tầng đất thành **bàn cờ** — đã thử bằng code và hỏng. Vân phải **đọc ra được ở cỡ thật** (mạch đá, rễ, cuội), không phải nhiễu hạt — đó là lý do 4 tile `lienmach` đợt 2 lặp rất tốt mà vẫn bị đánh 'trơn'.
- **Không sinh 4 tile từ một mẫu** — trong game chúng nằm cạnh nhau, mắt bắt ngay cái lặp (SO-TAY §6).
- **Tile mặt đi** (`co_mat`/`cat_mat`/`da_mat`): ~35 % trên là mặt (mép **răng cưa không đều**, không phải đường thẳng), 65 % dưới là đúng vật liệu của `dat_1` để nối xuống. Phần trên mép cỏ là **trong suốt**.
- **Lối lát chỉ là MẶT** (~30 % trên), dưới vẫn là lòng đất — bài học SOI-LAN-3 #19: tile lát cả ô làm con đường đọc thành 'rãnh sỏi đào xuống'. Hai ô chuyển tiếp mỗi đầu có ranh giới lởm chởm, vật liệu mặt lấn lên đá.
- **Không viền** trắng/đen quanh ô (7 tile map 1 từng dính, lộ thành ô viền trắng giữa mặt đất; `xoa-vien-tile.py` chữa được nhưng gọt mất ảnh).
- Tile mới phải được **khai là sàn** trong `tile_set_info.bin` (bit 0x0002, `build_inputs.tile_set_overlay`, `assetpack/tileset.py:1-33`) — không khai = **rơi xuyên đất**. Việc phía dự án, nhưng hoạ sĩ nên biết vì sao có tile đẹp mà nhân vật rơi.

### 3.5 Vật thể rời — luật đo được

- **Chân chạm ĐÚNG mép dưới ảnh**, không hở một hàng pixel, không thừa (hở 2 px = lơ lửng trong game; cổng torii từng mất cột phải vì cỡ lẻ bị phóng 96→128).
- **Nhìn ngang** (side view) — không 3/4, không từ trên xuống; nguyên vật thể trong một ảnh (không tách mái/tường).
- **Không bóng đổ, không mảng đất/cỏ dưới chân, không cảnh sau lưng** — bóng đổ là ảnh riêng `bong_*` đi đường bg-item (mục 8).
- Vật **ghép ngang** (tường thành, hàng rào, rào xích): mép trái/phải nối được với bản sao của chính nó.
- Tiền cảnh `tc_*`: cành lá **neo mép TRÊN**, bụi cỏ neo mép dưới; sẫm gần đen.
- Đường lắp: (a) chiếm một slot tile với ảnh to hơn ô (`tt_cay_nho` 384 → tile 87, `lap-anh-dot-2.py:34-38`) hoặc (b) **bg-item** `mapBackGround/<AnhID>.png`, AnhID ≥ 500, **thứ tự bảng −31 là giao thức, chỉ nối cuối** (`server/internal/game/bgitem.go:50-95`).

### 3.6 Điều KHÔNG được làm (gộp 'Sáu lỗi hay mắc' + bẫy đã vấp)

1. Pixel-art hay 'trộn' hai phong cách — **hand-painted, chốt một thứ**. 2. Cỡ lẻ, không phải bội 96. 3. Viền/khung quanh ảnh. 4. Nền **ghép đôi** nửa ảnh cho đủ rộng (chính là lỗi 1536 của map 1). 5. Đáy nền cắt phẳng còn vân. 6. Trời ở đỉnh `nen_langxa`. 7. Tile không lặp / 4 tile đất lệch sáng. 8. Vật thể hở chân, có bóng, có đất dưới chân. 9. Vẽ nhân vật/quái vào nền. 10. Tự đổi tên file hay tự chế cỡ — tên là khoá tra của máy.

---

## 4. Phong cách · ảnh tham chiếu · subject block

**Hand-painted (tranh vẽ mịn, cel-shading), KHÔNG pixel art** — chốt từ 22/08 (`DOT-2-ANH-CAN-VE.md`), vì nhà cửa/cây/cổng trong game và concept đã duyệt đều là tranh mịn; SOI-LAN-3 #5 chỉ đúng lỗi 'pixel thô đứng cạnh tranh mịn nhìn 1 giây thấy hai game'. Lưu ý: master prompt cũ ở `_KHUON.md` §3.5 ghi *pixel-art* — **với MAP thì không dùng dòng đó**, dùng khối dưới đây (đây là lệch có chủ đích, đã ghi ở `KHUON-DE-BAI-ANH-MAP.md`).

**Hoạ tiết phong cách khoá cho map (đã dán vào MỌI prompt ở mục 6):**

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures
```

**Ba khối luật theo kiểu ảnh (cũng đã nằm trong từng prompt):**

```text
NEN (background layer): One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally

TILE: Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise)

OBJECT: Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges
```

**Subject block tiếng Anh của hạng mục** (dán trước phần riêng của từng ảnh nếu công cụ cần một câu mở):

```text
hand-painted 2D side-scrolling game environment asset for a Naruto ninja village map, painterly anime background quality, cel-shaded, clean dark outlines, rich shading, straight side view, 1 grid cell = 96 px, player character about 246 px tall
```

**Bảng màu 6 map** (mở rộng từ `_KHUON.md` §3; câu tiếng Anh này nằm trong mọi prompt của map đó):

| Map | Hệ | Palette (EN, đã dán vào prompt) |
|---|---|---|
| Làng Lá (map 1) | Hoả | Hidden Leaf palette: warm orange-red clay roof tiles, aged brown timber, cream plaster walls, fresh green grass, warm brown soil, clear bright blue sky |
| Khu rừng (map 6) | — (map chung) | Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight |
| Làng Cát (map 2) | Phong | Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky |
| Làng Sương Mù (map 3) | Thuỷ | Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening |
| Làng Đá (map 4) | Thổ | Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky |
| Làng Mây (map 5) | Lôi | Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky |

**Ảnh tham chiếu — gắn kèm khi gen** (công cụ nhận tối đa 3 ảnh: luôn có `concept-map1-langla` + `THUOC-TI-LE` cho vật thể / `tile_co_mat` + `tile_dat_lienmach_1` cho tile / `nen_nuihokage_rong` cho nền; ảnh thứ 3 là mẫu cùng loại đã tốt của map 1):

| Khoá (dùng trong JSON `dung_tham_chieu`) | Đường dẫn | Học gì |
|---|---|---|
| `concept-map1-langla` | `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\concept-map1-langla.png` | ⭐ phong cách, tông màu, nét vẽ — cho MỌI map |
| `THUOC-TI-LE` | `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-tham-chieu\THUOC-TI-LE.png` | ⭐ vật nào cao bao nhiêu ô — gắn mọi lần gen vật thể |
| `TOT-cong-torii-da-tach-nen` | `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-tham-chieu\TOT-cong-torii-da-tach-nen.png` | chuẩn vàng vật thể rời: nguyên vật, nền sạch, chân chạm mép |
| `GIAI-THICH-tileset-vien-chuyen-tiep` | `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-tham-chieu\GIAI-THICH-tileset-vien-chuyen-tiep.png` | vì sao ô chuyển tiếp phải lởm chởm, cỏ phủ xuống |
| `HIEN-TRANG-map1-trong-game` | `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-tham-chieu\HIEN-TRANG-map1-trong-game.png` | map 1 trong game — cái đang chạy để bám tông |
| `nha_nho_map1` | `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-da-tao\nha_nho_x4.png` | mức chi tiết nhà đã duyệt |
| `tt_cay_map1` | `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-da-tao\tt_cay_x4.png` | mức chi tiết cây đã duyệt |
| `tile_dat_lienmach_1_map1` | `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-da-tao\tile_dat_lienmach_1_x4.png` | MÀU và độ sáng gốc của đất map 1 — 4 tile vân phải cùng độ sáng với nó |
| `tile_co_mat_map1` | `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-da-tao\tile_co_mat_x4.png` | mặt cỏ đang chạy — mép răng cưa, tỉ lệ cỏ/đất |
| `nen_nuihokage_rong_1536` | `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-da-tao\nen_nuihokage_rong_x4.png` | nội dung đã duyệt của núi Hokage (chỉ sai cỡ) — bản 2976 giữ đúng mặt và màu |
| `nen_lang_xa_rong_1536` | `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-da-tao\nen_lang_xa_rong_x4.png` | nội dung đã duyệt của làng xa (chỉ sai cỡ) |
| `TOT-naruto-logo` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-naruto-logo.png` | bảng màu thương hiệu, độ đậm viền |
| `XAU-vatpham` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\XAU-vatpham.png` | PHẢN DIỆN: phẳng, vài màu — tránh |
| `XAU-boss-zabuza` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\XAU-boss-zabuza.png` | PHẢN DIỆN: cục màu chân que — tránh |

---

## 5. DANH SÁCH ẢNH

Quy ước tên: **map 1 giữ tên hệ cũ** của `map-1-lang-la` (`nen_*`, `tile_*`, `tt_*`) vì 51 file và hai công cụ `kiem-anh-dot-2.py`/`lap-anh-dot-2.py` đang ghi cứng theo tên đó; **5 map mới dùng `map_<làng>_<tên>_x4.png`** theo tiền tố `map` của `kich-thuoc-anh.json` (làng = `rung` · `cat` · `smu` · `da` · `may`; `rung` đã là mã có sẵn trong bảng — `map_rung_*`). Toàn bộ tên dưới đây **chưa có** trong `kich-thuoc-anh.json` → nối cuối bảng trước khi kiểm (danh sách kèm cỡ ở mục 7.1). Kiểu: *tĩnh alpha* = PNG trong suốt; *tĩnh, nền ĐẶC* = nền không trong suốt.

### 5.1 Làng Lá (map 1) — 8 ảnh · màu nối `RGB(10, 86, 154)`

Chỉ phần **còn thiếu/lỗi** (49 ảnh đã có ở mục 2.1 không lặp lại). 6 ảnh P1 là đúng 'hai việc còn lại' của `DOT-2-ANH-CAN-VE.md` §Đợt 3; 2 ảnh P2 là **đề xuất** từ SOI-LAN-3, chưa ai chốt.

| STT | Tên file | Kiểu | Cỡ (px) | Làng | Mô tả | Ưu tiên · trạng thái |
|---|---|---|---|---|---|---|
| 1 | `nen_nuihokage_rong2_x4.png` | tĩnh, nền ĐẶC | 2976×768 | la | Nền núi Hokage vẽ LẠI cho đúng cỡ 2976 (bản 1536 hiện lặp 1,9 lần/khung). 5 khuôn mặt đá khác nhau; khối núi phải TO (chiếm ≥55% chiều cao ảnh) để trong game núi cao hơn mái nhà (SOI-LAN-3 #16); đáy nhoà về RGB(10, 86, 154). Không nhân đôi ảnh cũ — vẽ liền một dải núi dài. | **P1** · ve_lai |
| 2 | `nen_lang_xa_rong2_x4.png` | tĩnh, nền ĐẶC | 2976×576 | la | Nền làng xa vẽ LẠI cho đúng cỡ 2976. Mái nhà phải bắt đầu NGAY TỪ HÀNG ĐẦU (không có trời ở đỉnh — lớp này có nhiệm vụ che chân núi); làng dày, cây xen kẽ, mờ vì ở xa; đáy nhoà về RGB(10, 86, 154). | **P1** · ve_lai |
| 3 | `tile_dat_van_1_x4.png` | tĩnh alpha | 96×96 | la | Đất có MẠCH ĐÁ NẰM NGANG rõ + vài viên cuội nhỏ. Cùng độ sáng với tile_dat_lienmach_1. | **P1** · can_ve |
| 4 | `tile_dat_van_2_x4.png` | tĩnh alpha | 96×96 | la | Đất có RỄ CÂY thò từ mép trên xuống, uốn lượn. Cùng độ sáng. | **P1** · can_ve |
| 5 | `tile_dat_van_3_x4.png` | tĩnh alpha | 96×96 | la | Đất lẫn ĐÁ CUỘI TO (3–4 viên cỡ 20–30 px), tối hơn ~6%. | **P1** · can_ve |
| 6 | `tile_dat_van_4_x4.png` | tĩnh alpha | 96×96 | la | Đất chặt, vân nén NGANG mảnh, ít chi tiết nhất trong 4 — ô 'nghỉ' để xen giữa cho đỡ rối, tối hơn ~12%. | **P1** · can_ve |
| 7 | `tile_datlat_mat_x4.png` | tĩnh alpha | 96×96 | la | ĐỀ XUẤT (SOI-LAN-3 #19): lối lát đá chỉ chiếm ~30% trên cùng của ô (mặt đường), 70% dưới vẫn là đất có rễ y như tile_dat_lienmach_1 — để con đường đọc ra là 'rải đá lên mặt đất', không phải 'rãnh sỏi đào xuống'. | **P2** · can_ve |
| 8 | `tt_co_phu_chan_x4.png` | tĩnh alpha | 384×96 | la | ĐỀ XUẤT (SOI-LAN-3 #20, #11): dải cỏ dại + hoa nhỏ + đá con, mép trên lởm chởm, để VẼ TRƯỚC (đè lên) chân nhà/giếng/thùng cho vật 'cắm xuống đất' thay vì hở trời 3 px. Đường lắp (slot tile hay bg-item) CHƯA có — vẽ trước, lắp sau. | **P2** · can_ve |

### 5.2 Khu rừng (map 6) — 30 ảnh · màu nối `RGB(52, 88, 72)`

Rừng rậm, **không nhà** — 5 slot 'công trình' của khuôn thay bằng 2 cây cổ thụ khác dáng + gốc/đá/thân đổ/cầu gỗ/biển. ⚠️ Máy chủ chưa nối map 6 (`banDo.go:81-89`), vẽ trước lắp sau.

| STT | Tên file | Kiểu | Cỡ (px) | Làng | Mô tả | Ưu tiên · trạng thái |
|---|---|---|---|---|---|---|
| 1 | `map_rung_nen_troi_x4.png` | tĩnh, nền ĐẶC | 2976×384 | rung | Lớp trời của Khu rừng: trời xanh + mây, nhạt dần xuống, 20% đáy là màu nối thuần RGB(52, 88, 72). Không cây, không đất. | **P1** · can_ve |
| 2 | `map_rung_nen_nui_x4.png` | tĩnh, nền ĐẶC | 2976×768 | rung | Lớp núi của Khu rừng: 2–3 dãy núi phủ rừng xanh thẫm chồng lớp, sương giữa các dãy; khối núi chiếm ≥55% chiều cao; đáy nhoà về màu nối. | **P1** · can_ve |
| 3 | `map_rung_nen_langxa_x4.png` | tĩnh, nền ĐẶC | 2976×576 | rung | Lớp 'làng xa' của Khu rừng = TƯỜNG TÁN RỪNG XA (không nhà): tán cây khổng lồ chồng nhau bắt đầu ngay hàng đầu (KHÔNG có trời ở đỉnh — lớp này che chân núi), thân cây lộ ở nửa dưới, mờ vì xa; đáy nhoà về màu nối. | **P1** · can_ve |
| 4 | `map_rung_nen_gan_x4.png` | tĩnh, nền ĐẶC | 2976×288 | rung | Lớp gần nhất của Khu rừng: thân cây và dương xỉ mờ trong sương, 20% đáy là màu nối thuần. | **P1** · can_ve |
| 5 | `map_rung_tile_co_mat_x4.png` | tĩnh alpha | 96×96 | rung | Mặt cỏ rừng: ~35% trên là cỏ xanh thẫm mép răng cưa không đều, lá rụng, rêu; 65% dưới là đất rừng có rễ (nối vào tile_dat_1). | **P1** · can_ve |
| 6 | `map_rung_tile_co_sang_da_x4.png` | tĩnh alpha | 96×96 | rung | Chuyển tiếp cỏ → lối đá rêu: nửa trái cỏ, nửa phải đá lát cũ phủ rêu, ranh giới lởm chởm, cỏ lấn lên đá. | **P1** · can_ve |
| 7 | `map_rung_tile_da_sang_co_x4.png` | tĩnh alpha | 96×96 | rung | Ảnh gương của tile_co_sang_da (đá bên trái, cỏ bên phải). | **P1** · can_ve |
| 8 | `map_rung_tile_datlat_x4.png` | tĩnh alpha | 96×96 | rung | Lối đá cổ phủ rêu: ~30% trên là phiến đá xám xanh mòn, rêu ở khe; 70% dưới là đất rừng. | **P1** · can_ve |
| 9 | `map_rung_tile_dat_1_x4.png` | tĩnh alpha | 96×96 | rung | Đất rừng có mạch đá ngang + cuội nhỏ. | **P1** · can_ve |
| 10 | `map_rung_tile_dat_2_x4.png` | tĩnh alpha | 96×96 | rung | Đất rừng có rễ to thò xuống (rừng nên rễ dày, nhiều hơn Làng Lá). | **P1** · can_ve |
| 11 | `map_rung_tile_dat_giua_x4.png` | tĩnh alpha | 96×96 | rung | Đất lẫn đá cuội to, tối hơn ~6%. | **P1** · can_ve |
| 12 | `map_rung_tile_dat_day_x4.png` | tĩnh alpha | 96×96 | rung | Đất chặt vân nén ngang, tối hơn ~12%. | **P1** · can_ve |
| 13 | `map_rung_tile_vach_trai_x4.png` | tĩnh alpha | 96×96 | rung | Mặt cắt vách đất bên TRÁI gò: mép trái là đất lộ thiên (rễ, đá nhỏ, cỏ rủ ở góc trên), mép phải nối vào lòng đất. | **P1** · can_ve |
| 14 | `map_rung_tile_vach_phai_x4.png` | tĩnh alpha | 96×96 | rung | Ảnh gương của tile_vach_trai (mặt lộ thiên bên phải). | **P1** · can_ve |
| 15 | `map_rung_tile_chan_vach_x4.png` | tĩnh alpha | 96×96 | rung | Ô chân vách: chỗ vách đứng gặp mặt đất phẳng — đất tơi trượt xuống, cuội, một túm cỏ; nối mềm giữa vách và đất thường. | **P1** · can_ve |
| 16 | `map_rung_tile_ke_da_x4.png` | tĩnh alpha | 96×96 | rung | Kè đá xếp khô phủ rêu (thay vách đất ở khu trung tâm), lặp được. | **P1** · can_ve |
| 17 | `map_rung_tile_bac_thang_x4.png` | tĩnh alpha | 96×96 | rung | Một bậc thang đá cổ phủ rêu: mặt bậc phẳng ở trên, mặt đứng bên trái, phần dưới-phải đặc để xếp chéo với chính nó. | **P1** · can_ve |
| 18 | `map_rung_tile_bac_thang_dinh_x4.png` | tĩnh alpha | 96×96 | rung | Bậc trên cùng nối vào mặt cỏ rừng. | **P1** · can_ve |
| 19 | `map_rung_cay_co_thu_1_x4.png` | tĩnh alpha | 768×768 | rung | Cây cổ thụ rừng #1: thân khổng lồ có rễ bạnh, vỏ sần, tán rộng dày nhiều lớp. 8 ô cao. | **P1** · can_ve |
| 20 | `map_rung_cay_co_thu_2_x4.png` | tĩnh alpha | 768×768 | rung | Cây cổ thụ rừng #2: dáng KHÁC hẳn #1 — thân nghiêng, có hốc, dây leo rủ, tán lệch — để rải xen không lặp. | **P1** · can_ve |
| 21 | `map_rung_cay_nho_x4.png` | tĩnh alpha | 384×384 | rung | Cây non tán tròn, thân mảnh. | **P1** · can_ve |
| 22 | `map_rung_bui_ram_x4.png` | tĩnh alpha | 288×192 | rung | Bụi dương xỉ / bụi gai thấp rậm. | **P1** · can_ve |
| 23 | `map_rung_goc_cay_x4.png` | tĩnh alpha | 288×192 | rung | Gốc cây cụt có vân tuổi, rêu, vài cái nấm. | **P1** · can_ve |
| 24 | `map_rung_tang_da_x4.png` | tĩnh alpha | 288×192 | rung | Cụm tảng đá phủ rêu. | **P1** · can_ve |
| 25 | `map_rung_cay_do_x4.png` | tĩnh alpha | 576×192 | rung | Thân cây đổ nằm ngang, rỗng ruột, rêu và nấm. | **P1** · can_ve |
| 26 | `map_rung_cau_go_x4.png` | tĩnh alpha | 480×192 | rung | Cầu gỗ nhỏ kiểu rừng (ván cũ, lan can dây thừng) — vật trang trí đặt sau lối đi. | **P1** · can_ve |
| 27 | `map_rung_bien_chi_duong_x4.png` | tĩnh alpha | 192×288 | rung | Cột biển chỉ đường gỗ hai mũi tên, khắc ký hiệu (không chữ đọc được), dây thừng buộc. | **P1** · can_ve |
| 28 | `map_rung_bui_co_x4.png` | tĩnh alpha | 192×96 | rung | Túm cỏ dại + dương xỉ nhỏ + hoa li ti. | **P1** · can_ve |
| 29 | `map_rung_tc_canh_la_x4.png` | tĩnh alpha | 768×384 | rung | Tiền cảnh: cành lá rủ từ MÉP TRÊN, xanh sẫm gần đen (ngược sáng, gần camera). | **P2** · can_ve |
| 30 | `map_rung_tc_bui_truoc_x4.png` | tĩnh alpha | 384×192 | rung | Tiền cảnh: bụi cỏ/dương xỉ cao sẫm màu mọc từ MÉP DƯỚI. | **P2** · can_ve |

### 5.3 Làng Cát (map 2) — 31 ảnh · màu nối `RGB(226, 190, 120)`

| STT | Tên file | Kiểu | Cỡ (px) | Làng | Mô tả | Ưu tiên · trạng thái |
|---|---|---|---|---|---|---|
| 1 | `map_cat_nen_troi_x4.png` | tĩnh, nền ĐẶC | 2976×384 | cat | Lớp trời: trời vàng cam nắng gắt, bụi cát mờ thấp, mây mỏng. 20% đáy là màu nối thuần RGB(226, 190, 120). | **P2** · can_ve |
| 2 | `map_cat_nen_nui_x4.png` | tĩnh, nền ĐẶC | 2976×768 | cat | Lớp núi: vách đá sa thạch tầng lớp và các mesa bao quanh làng, vài đụn cát xa. Khối núi chiếm ≥55% chiều cao (do-map.py 'khối núi cao ≥30% khung'); đáy nhoà về RGB(226, 190, 120). | **P2** · can_ve |
| 3 | `map_cat_nen_langxa_x4.png` | tĩnh, nền ĐẶC | 2976×576 | cat | Lớp làng xa: mái vòm đất nung, tháp tròn gạch bùn, mái bằng với mái hiên vải, tường thành nhiều lớp (bù dải giữa vì sa mạc không có cây). Mái/khối bắt đầu NGAY HÀNG ĐẦU (không trời ở đỉnh — lớp này che chân núi); đáy nhoà về RGB(226, 190, 120). | **P2** · can_ve |
| 4 | `map_cat_nen_gan_x4.png` | tĩnh, nền ĐẶC | 2976×288 | cat | Lớp gần: đỉnh tường thành đất nung mờ + bụi cát là là. 20% đáy là màu nối thuần. | **P2** · can_ve |
| 5 | `map_cat_tile_cat_mat_x4.png` | tĩnh alpha | 96×96 | cat | Mặt đi trên cùng: cát mịn gợn sóng gió, vài viên sỏi. ~35% trên là mặt, 65% dưới nối vào tile_dat_1. | **P2** · can_ve |
| 6 | `map_cat_tile_cat_sang_da_x4.png` | tĩnh alpha | 96×96 | cat | Chuyển tiếp cát → đá lát: nửa trái là mặt, nửa phải là lối lát, ranh giới lởm chởm tự nhiên. | **P2** · can_ve |
| 7 | `map_cat_tile_da_sang_cat_x4.png` | tĩnh alpha | 96×96 | cat | Chuyển tiếp đá lát → cát (ảnh gương của ô trên). | **P2** · can_ve |
| 8 | `map_cat_tile_datlat_x4.png` | tĩnh alpha | 96×96 | cat | Lối lát: đá sa thạch lát phẳng màu be hồng. chỉ ~30% trên là mặt lát, 70% dưới là đất/đá lòng đất (bài học SOI-LAN-3 #19: đường là mặt, không phải rãnh). | **P2** · can_ve |
| 9 | `map_cat_tile_dat_1_x4.png` | tĩnh alpha | 96×96 | cat | Lòng đất 1: đất nén khô có mạch sa thạch ngang + sỏi. | **P2** · can_ve |
| 10 | `map_cat_tile_dat_2_x4.png` | tĩnh alpha | 96×96 | cat | Lòng đất 2: đất khô nứt nẻ, vài rễ cây chết. | **P2** · can_ve |
| 11 | `map_cat_tile_dat_giua_x4.png` | tĩnh alpha | 96×96 | cat | Lòng đất giữa: đất lẫn khối sa thạch to, tối hơn ~6%. | **P2** · can_ve |
| 12 | `map_cat_tile_dat_day_x4.png` | tĩnh alpha | 96×96 | cat | Lòng đất đáy: trầm tích nén thành dải ngang, tối hơn ~12%. | **P2** · can_ve |
| 13 | `map_cat_tile_vach_trai_x4.png` | tĩnh alpha | 96×96 | cat | Vách trái: mặt cắt bờ cát-đất có tầng sa thạch. mép TRÁI lộ thiên, mép phải nối vào lòng đất. | **P2** · can_ve |
| 14 | `map_cat_tile_vach_phai_x4.png` | tĩnh alpha | 96×96 | cat | Vách phải: ảnh gương của tile_vach_trai (mặt lộ thiên bên PHẢI). | **P2** · can_ve |
| 15 | `map_cat_tile_chan_vach_x4.png` | tĩnh alpha | 96×96 | cat | Ô chân vách: chỗ vách đứng gặp mặt đất phẳng — vật liệu tơi trượt xuống, nối mềm sang ô thường. | **P2** · can_ve |
| 16 | `map_cat_tile_ke_da_x4.png` | tĩnh alpha | 96×96 | cat | Kè đá: tường kè khối sa thạch xếp, mòn gió. lặp được, dùng thay vách ở khu trung tâm. | **P2** · can_ve |
| 17 | `map_cat_tile_bac_thang_x4.png` | tĩnh alpha | 96×96 | cat | Một bậc thang: bậc đá sa thạch mòn. mặt bậc phẳng trên, mặt đứng bên trái, dưới-phải đặc để xếp chéo. | **P2** · can_ve |
| 18 | `map_cat_tile_bac_thang_dinh_x4.png` | tĩnh alpha | 96×96 | cat | Bậc trên cùng nối vào mặt đi. | **P2** · can_ve |
| 19 | `map_cat_nha_nho_x4.png` | tĩnh alpha | 576×576 | cat | Nhà nhỏ Suna: tường đất nung tròn, mái vòm, cửa gỗ, cửa sổ vuông nhỏ, mái hiên vải. 6 ô cao. | **P2** · can_ve |
| 20 | `map_cat_nha_vua_x4.png` | tĩnh alpha | 768×576 | cat | Nhà vừa Suna: rộng, mái bằng có lan can, cửa vòm, mái hiên sọc, vại nước trên mái. | **P2** · can_ve |
| 21 | `map_cat_nha_lon_x4.png` | tĩnh alpha | 768×768 | cat | Nhà lớn Suna: tháp đất nung 2 tầng, mái vòm + ống thông gió, cầu thang ngoài. | **P2** · can_ve |
| 22 | `map_cat_nha_dacsan_x4.png` | tĩnh alpha | 768×576 | cat | Đặc sản Cát: quán chợ sa mạc / xưởng rối gỗ — mặt trước mở, vải treo, bình gốm, chuông gió. | **P2** · can_ve |
| 23 | `map_cat_cong_lang_x4.png` | tĩnh alpha | 480×480 | cat | Cổng làng Cát: hai trụ sa thạch to khắc biểu tượng gió, xà gỗ, cát bám chân. | **P2** · can_ve |
| 24 | `map_cat_cay_cha_la_x4.png` | tĩnh alpha | 384×768 | cat | Cây chà là cao, thân đốt, tán lá xoè. | **P2** · can_ve |
| 25 | `map_cat_xuong_rong_x4.png` | tĩnh alpha | 192×288 | cat | Cụm xương rồng cao có hoa. | **P2** · can_ve |
| 26 | `map_cat_bui_gai_x4.png` | tĩnh alpha | 288×192 | cat | Bụi gai khô / cỏ lăn. | **P2** · can_ve |
| 27 | `map_cat_gieng_x4.png` | tĩnh alpha | 288×288 | cat | Giếng đá sa mạc, khung gỗ, mái vải che, gàu da. | **P2** · can_ve |
| 28 | `map_cat_den_dau_x4.png` | tĩnh alpha | 192×384 | cat | Cột đèn dầu đồng, chân đá. | **P2** · can_ve |
| 29 | `map_cat_binh_gom_x4.png` | tĩnh alpha | 192×192 | cat | Cụm 2–3 bình gốm/vò nước, hoa văn sọc. | **P2** · can_ve |
| 30 | `map_cat_tuong_thanh_x4.png` | tĩnh alpha | 768×384 | cat | Đoạn tường thành đất nung có lỗ châu mai, GHÉP NGANG liền mạch được — bù dải giữa vì sa mạc không cây. | **P2** · can_ve |
| 31 | `map_cat_bui_co_kho_x4.png` | tĩnh alpha | 192×96 | cat | Túm cỏ khô + cát. | **P2** · can_ve |

### 5.4 Làng Sương Mù (map 3) — 33 ảnh · màu nối `RGB(150, 172, 184)`

| STT | Tên file | Kiểu | Cỡ (px) | Làng | Mô tả | Ưu tiên · trạng thái |
|---|---|---|---|---|---|---|
| 1 | `map_smu_nen_troi_x4.png` | tĩnh, nền ĐẶC | 2976×384 | smu | Lớp trời: trời xám lam u ám, dải sương thấp. 20% đáy là màu nối thuần RGB(150, 172, 184). | **P2** · can_ve |
| 2 | `map_smu_nen_nui_x4.png` | tĩnh, nền ĐẶC | 2976×768 | smu | Lớp núi: mũi đá đảo tối, nhọn, nhô khỏi mặt biển sương, chồng lớp. Khối núi chiếm ≥55% chiều cao (do-map.py 'khối núi cao ≥30% khung'); đáy nhoà về RGB(150, 172, 184). | **P2** · can_ve |
| 3 | `map_smu_nen_langxa_x4.png` | tĩnh, nền ĐẶC | 2976×576 | smu | Lớp làng xa: mái dốc ngói xanh xám của Kiri trên các tháp đá trụ tròn, cầu treo và ống nước nối nhau, sương giăng. Mái/khối bắt đầu NGAY HÀNG ĐẦU (không trời ở đỉnh — lớp này che chân núi); đáy nhoà về RGB(150, 172, 184). | **P2** · can_ve |
| 4 | `map_smu_nen_gan_x4.png` | tĩnh, nền ĐẶC | 2976×288 | smu | Lớp gần: cọc bến và cột neo thuyền mờ trong sương dày. 20% đáy là màu nối thuần. | **P2** · can_ve |
| 5 | `map_smu_tile_co_mat_x4.png` | tĩnh alpha | 96×96 | smu | Mặt đi trên cùng: cỏ ẩm lẫn rêu, mép cỏ ướt. ~35% trên là mặt, 65% dưới nối vào tile_dat_1. | **P2** · can_ve |
| 6 | `map_smu_tile_co_sang_da_x4.png` | tĩnh alpha | 96×96 | smu | Chuyển tiếp cỏ ẩm → đá lát ướt: nửa trái là mặt, nửa phải là lối lát, ranh giới lởm chởm tự nhiên. | **P2** · can_ve |
| 7 | `map_smu_tile_da_sang_co_x4.png` | tĩnh alpha | 96×96 | smu | Chuyển tiếp đá lát ướt → cỏ ẩm (ảnh gương của ô trên). | **P2** · can_ve |
| 8 | `map_smu_tile_datlat_x4.png` | tĩnh alpha | 96×96 | smu | Lối lát: phiến đá lát xám xanh ướt, rêu ở khe. chỉ ~30% trên là mặt lát, 70% dưới là đất/đá lòng đất (bài học SOI-LAN-3 #19: đường là mặt, không phải rãnh). | **P2** · can_ve |
| 9 | `map_smu_tile_dat_1_x4.png` | tĩnh alpha | 96×96 | smu | Lòng đất 1: đất ướt có mạch đá xám ngang + cuội. | **P2** · can_ve |
| 10 | `map_smu_tile_dat_2_x4.png` | tĩnh alpha | 96×96 | smu | Lòng đất 2: đất có rễ và vệt nước rỉ. | **P2** · can_ve |
| 11 | `map_smu_tile_dat_giua_x4.png` | tĩnh alpha | 96×96 | smu | Lòng đất giữa: đất lẫn cuội to ướt bóng, tối hơn ~6%. | **P2** · can_ve |
| 12 | `map_smu_tile_dat_day_x4.png` | tĩnh alpha | 96×96 | smu | Lòng đất đáy: sét nén dải ngang, tối hơn ~12%. | **P2** · can_ve |
| 13 | `map_smu_tile_vach_trai_x4.png` | tĩnh alpha | 96×96 | smu | Vách trái: mặt cắt bờ đất ướt, rêu rủ, nước nhỏ giọt. mép TRÁI lộ thiên, mép phải nối vào lòng đất. | **P2** · can_ve |
| 14 | `map_smu_tile_vach_phai_x4.png` | tĩnh alpha | 96×96 | smu | Vách phải: ảnh gương của tile_vach_trai (mặt lộ thiên bên PHẢI). | **P2** · can_ve |
| 15 | `map_smu_tile_chan_vach_x4.png` | tĩnh alpha | 96×96 | smu | Ô chân vách: chỗ vách đứng gặp mặt đất phẳng — vật liệu tơi trượt xuống, nối mềm sang ô thường. | **P2** · can_ve |
| 16 | `map_smu_tile_ke_da_x4.png` | tĩnh alpha | 96×96 | smu | Kè đá: tường đá xếp ướt phủ rêu, nước rỉ. lặp được, dùng thay vách ở khu trung tâm. | **P2** · can_ve |
| 17 | `map_smu_tile_bac_thang_x4.png` | tĩnh alpha | 96×96 | smu | Một bậc thang: bậc đá ướt có rêu ở mép. mặt bậc phẳng trên, mặt đứng bên trái, dưới-phải đặc để xếp chéo. | **P2** · can_ve |
| 18 | `map_smu_tile_bac_thang_dinh_x4.png` | tĩnh alpha | 96×96 | smu | Bậc trên cùng nối vào mặt đi. | **P2** · can_ve |
| 19 | `map_smu_tile_nuoc_mat_x4.png` | tĩnh alpha | 96×96 | smu | Mặt nước: mép trên có gợn sóng nhỏ và vệt bọt, dưới là nước xanh lam sâu — lặp trái/phải. | **P2** · can_ve |
| 20 | `map_smu_tile_nuoc_sau_x4.png` | tĩnh alpha | 96×96 | smu | Nước sâu lặp 4 phía, xanh lam sẫm, vài vệt sáng. | **P2** · can_ve |
| 21 | `map_smu_nha_nho_x4.png` | tĩnh alpha | 576×576 | smu | Nhà nhỏ Kiri: đế đá trụ tròn, mái dốc ngói xanh xám ướt, cửa sổ tròn, rêu chân tường. 6 ô. | **P2** · can_ve |
| 22 | `map_smu_nha_vua_x4.png` | tĩnh alpha | 768×576 | smu | Nhà vừa Kiri: nhà bến trên cọc gỗ, lưới cá treo, đèn lồng xanh. | **P2** · can_ve |
| 23 | `map_smu_nha_lon_x4.png` | tĩnh alpha | 768×768 | smu | Nhà lớn Kiri: tháp cao hai tầng mái dốc chồng, ống nước ngoài, cầu nhỏ. | **P2** · can_ve |
| 24 | `map_smu_nha_dacsan_x4.png` | tĩnh alpha | 768×576 | smu | Đặc sản Sương Mù: chợ cá / nhà thuyền — mặt trước mở, cá treo, thùng, thuyền nhỏ úp. | **P2** · can_ve |
| 25 | `map_smu_cong_lang_x4.png` | tĩnh alpha | 480×480 | smu | Cổng làng Sương Mù: hai trụ đá khắc 4 đường sương lượn sóng, xà gỗ, rêu rủ. | **P2** · can_ve |
| 26 | `map_smu_cay_lieu_x4.png` | tĩnh alpha | 768×768 | smu | Cây liễu rủ trong sương, thân xoắn, lá rủ dài. | **P2** · can_ve |
| 27 | `map_smu_cay_nho_x4.png` | tĩnh alpha | 384×384 | smu | Cây nhỏ ẩm rêu, tán tròn thẫm. | **P2** · can_ve |
| 28 | `map_smu_bui_lau_x4.png` | tĩnh alpha | 288×192 | smu | Bụi lau sậy ven nước. | **P2** · can_ve |
| 29 | `map_smu_cau_go_x4.png` | tĩnh alpha | 768×288 | smu | Cầu gỗ qua nước: ván ẩm bạc màu, lan can dây thừng, rêu — vật đặc trưng của làng. | **P2** · can_ve |
| 30 | `map_smu_cot_buoc_thuyen_x4.png` | tĩnh alpha | 192×288 | smu | Cột neo thuyền quấn dây thừng, hà bám gốc. | **P2** · can_ve |
| 31 | `map_smu_den_hieu_x4.png` | tĩnh alpha | 192×384 | smu | Cột đèn hiệu bến cảng, ánh xanh lạnh. | **P2** · can_ve |
| 32 | `map_smu_thung_ca_x4.png` | tĩnh alpha | 192×192 | smu | Thùng cá + lưới + phao. | **P2** · can_ve |
| 33 | `map_smu_bui_co_uot_x4.png` | tĩnh alpha | 192×96 | smu | Túm cỏ ướt + vũng nước nhỏ. | **P2** · can_ve |

### 5.5 Làng Đá (map 4) — 31 ảnh · màu nối `RGB(142, 140, 132)`

| STT | Tên file | Kiểu | Cỡ (px) | Làng | Mô tả | Ưu tiên · trạng thái |
|---|---|---|---|---|---|---|
| 1 | `map_da_nen_troi_x4.png` | tĩnh, nền ĐẶC | 2976×384 | da | Lớp trời: trời xám xanh bụi, mây mỏng, bụi đá ấm ở thấp. 20% đáy là màu nối thuần RGB(142, 140, 132). | **P2** · can_ve |
| 2 | `map_da_nen_nui_x4.png` | tĩnh, nền ĐẶC | 2976×768 | da | Lớp núi: vách hẻm núi khổng lồ nhiều tầng trầm tích nâu đỏ–xám, chóp đá nhọn, dựng đứng. Khối núi chiếm ≥55% chiều cao (do-map.py 'khối núi cao ≥30% khung'); đáy nhoà về RGB(142, 140, 132). | **P2** · can_ve |
| 3 | `map_da_nen_langxa_x4.png` | tĩnh, nền ĐẶC | 2976×576 | da | Lớp làng xa: nhà đục trong vách đá, nhà chồng tầng bằng khối đá, cửa sổ vuông nhỏ, cầu đá nối các gờ. Mái/khối bắt đầu NGAY HÀNG ĐẦU (không trời ở đỉnh — lớp này che chân núi); đáy nhoà về RGB(142, 140, 132). | **P2** · can_ve |
| 4 | `map_da_nen_gan_x4.png` | tĩnh, nền ĐẶC | 2976×288 | da | Lớp gần: khối đá mờ + bụi. 20% đáy là màu nối thuần. | **P2** · can_ve |
| 5 | `map_da_tile_da_mat_x4.png` | tĩnh alpha | 96×96 | da | Mặt đi trên cùng: mặt đá tự nhiên mòn phẳng, vết nứt, vài túm cỏ khô. ~35% trên là mặt, 65% dưới nối vào tile_dat_1. | **P2** · can_ve |
| 6 | `map_da_tile_da_sang_lat_x4.png` | tĩnh alpha | 96×96 | da | Chuyển tiếp đá tự nhiên → đá lát: nửa trái là mặt, nửa phải là lối lát, ranh giới lởm chởm tự nhiên. | **P2** · can_ve |
| 7 | `map_da_tile_lat_sang_da_x4.png` | tĩnh alpha | 96×96 | da | Chuyển tiếp đá lát → đá tự nhiên (ảnh gương của ô trên). | **P2** · can_ve |
| 8 | `map_da_tile_datlat_x4.png` | tĩnh alpha | 96×96 | da | Lối lát: đá đục phẳng xếp khít, màu xám ấm. chỉ ~30% trên là mặt lát, 70% dưới là đất/đá lòng đất (bài học SOI-LAN-3 #19: đường là mặt, không phải rãnh). | **P2** · can_ve |
| 9 | `map_da_tile_dat_1_x4.png` | tĩnh alpha | 96×96 | da | Lòng đất 1: đất đá có dải trầm tích đỏ ngang + sỏi. | **P2** · can_ve |
| 10 | `map_da_tile_dat_2_x4.png` | tĩnh alpha | 96×96 | da | Lòng đất 2: đất đá có mạch khoáng thạch anh sáng và vết nứt. | **P2** · can_ve |
| 11 | `map_da_tile_dat_giua_x4.png` | tĩnh alpha | 96×96 | da | Lòng đất giữa: đá lẫn khối đá to góc cạnh, tối hơn ~6%. | **P2** · can_ve |
| 12 | `map_da_tile_dat_day_x4.png` | tĩnh alpha | 96×96 | da | Lòng đất đáy: trầm tích dày nhiều dải sát nhau, tối hơn ~12%. | **P2** · can_ve |
| 13 | `map_da_tile_vach_trai_x4.png` | tĩnh alpha | 96×96 | da | Vách trái: vách đá tầng trầm tích dựng đứng. mép TRÁI lộ thiên, mép phải nối vào lòng đất. | **P2** · can_ve |
| 14 | `map_da_tile_vach_phai_x4.png` | tĩnh alpha | 96×96 | da | Vách phải: ảnh gương của tile_vach_trai (mặt lộ thiên bên PHẢI). | **P2** · can_ve |
| 15 | `map_da_tile_chan_vach_x4.png` | tĩnh alpha | 96×96 | da | Ô chân vách: chỗ vách đứng gặp mặt đất phẳng — vật liệu tơi trượt xuống, nối mềm sang ô thường. | **P2** · can_ve |
| 16 | `map_da_tile_ke_da_x4.png` | tĩnh alpha | 96×96 | da | Kè đá: tường đá đục khối xếp khít, mạch vữa tối. lặp được, dùng thay vách ở khu trung tâm. | **P2** · can_ve |
| 17 | `map_da_tile_bac_thang_x4.png` | tĩnh alpha | 96×96 | da | Một bậc thang: bậc thang đục thẳng vào đá núi. mặt bậc phẳng trên, mặt đứng bên trái, dưới-phải đặc để xếp chéo. | **P2** · can_ve |
| 18 | `map_da_tile_bac_thang_dinh_x4.png` | tĩnh alpha | 96×96 | da | Bậc trên cùng nối vào mặt đi. | **P2** · can_ve |
| 19 | `map_da_nha_nho_x4.png` | tĩnh alpha | 576×576 | da | Nhà nhỏ Iwa: đục vào một khối đá, mặt tiền đá thô, cửa gỗ dày, cửa sổ tròn. 6 ô. | **P2** · can_ve |
| 20 | `map_da_nha_vua_x4.png` | tĩnh alpha | 768×576 | da | Nhà vừa Iwa: nhà khối đá xếp, mái phiến đá nặng, ống khói đá. | **P2** · can_ve |
| 21 | `map_da_nha_lon_x4.png` | tĩnh alpha | 768×768 | da | Nhà lớn Iwa: tháp đá cao nhiều tầng đục, thang gỗ nối tầng. | **P2** · can_ve |
| 22 | `map_da_nha_dacsan_x4.png` | tĩnh alpha | 768×576 | da | Đặc sản Đá: lò rèn / xưởng khai đá — đe, xe goòng, lò nung đỏ. | **P2** · can_ve |
| 23 | `map_da_cong_lang_x4.png` | tĩnh alpha | 480×480 | da | Cổng làng Đá: hai khối đá nguyên khối thô khắc biểu tượng Thổ, đá đè ngang nặng. | **P2** · can_ve |
| 24 | `map_da_cay_kho_x4.png` | tĩnh alpha | 384×768 | da | Cây khô trơ cành xoắn. | **P2** · can_ve |
| 25 | `map_da_cay_nho_x4.png` | tĩnh alpha | 384×384 | da | Cây bụi núi cứng cỏi, lá nhỏ. | **P2** · can_ve |
| 26 | `map_da_bui_ram_x4.png` | tĩnh alpha | 288×192 | da | Bụi cây khô cằn. | **P2** · can_ve |
| 27 | `map_da_duoc_x4.png` | tĩnh alpha | 192×384 | da | Cột đuốc giá sắt, lửa cam. | **P2** · can_ve |
| 28 | `map_da_tang_da_x4.png` | tĩnh alpha | 288×288 | da | Tảng đá lớn có dải trầm tích. | **P2** · can_ve |
| 29 | `map_da_xe_day_x4.png` | tĩnh alpha | 288×192 | da | Xe goòng chở quặng. | **P2** · can_ve |
| 30 | `map_da_hang_rao_da_x4.png` | tĩnh alpha | 192×192 | da | Rào đá xếp thấp, ghép ngang liền mạch được. | **P2** · can_ve |
| 31 | `map_da_bui_co_kho_x4.png` | tĩnh alpha | 192×96 | da | Túm cỏ khô giữa đá. | **P2** · can_ve |

### 5.6 Làng Mây (map 5) — 32 ảnh · màu nối `RGB(120, 110, 160)`

| STT | Tên file | Kiểu | Cỡ (px) | Làng | Mô tả | Ưu tiên · trạng thái |
|---|---|---|---|---|---|---|
| 1 | `map_may_nen_troi_x4.png` | tĩnh, nền ĐẶC | 2976×384 | may | Lớp trời: trời tím lam hoàng hôn, mây bão, ánh sét xa. 20% đáy là màu nối thuần RGB(120, 110, 160). | **P2** · can_ve |
| 2 | `map_may_nen_nui_x4.png` | tĩnh, nền ĐẶC | 2976×768 | may | Lớp núi: đỉnh núi nhọn như kim chọc lên biển mây, granite xám xanh. Khối núi chiếm ≥55% chiều cao (do-map.py 'khối núi cao ≥30% khung'); đáy nhoà về RGB(120, 110, 160). | **P2** · can_ve |
| 3 | `map_may_nen_langxa_x4.png` | tĩnh, nền ĐẶC | 2976×576 | may | Lớp làng xa: tháp trụ tròn mái nhọn xây trên đỉnh núi, cầu treo nối nhau, cột thu lôi phát sáng, mây quấn chân. Mái/khối bắt đầu NGAY HÀNG ĐẦU (không trời ở đỉnh — lớp này che chân núi); đáy nhoà về RGB(120, 110, 160). | **P2** · can_ve |
| 4 | `map_may_nen_gan_x4.png` | tĩnh, nền ĐẶC | 2976×288 | may | Lớp gần: dải mây dày là là + chóp tháp mờ. 20% đáy là màu nối thuần. | **P2** · can_ve |
| 5 | `map_may_tile_da_mat_x4.png` | tĩnh alpha | 96×96 | may | Mặt đi trên cùng: mặt granite xám xanh, rêu thưa, sương mỏng. ~35% trên là mặt, 65% dưới nối vào tile_dat_1. | **P2** · can_ve |
| 6 | `map_may_tile_da_sang_lat_x4.png` | tĩnh alpha | 96×96 | may | Chuyển tiếp đá tự nhiên → đá lát: nửa trái là mặt, nửa phải là lối lát, ranh giới lởm chởm tự nhiên. | **P2** · can_ve |
| 7 | `map_may_tile_lat_sang_da_x4.png` | tĩnh alpha | 96×96 | may | Chuyển tiếp đá lát → đá tự nhiên (ảnh gương của ô trên). | **P2** · can_ve |
| 8 | `map_may_tile_datlat_x4.png` | tĩnh alpha | 96×96 | may | Lối lát: phiến granite xám xanh lát khít, ánh tím nhẹ. chỉ ~30% trên là mặt lát, 70% dưới là đất/đá lòng đất (bài học SOI-LAN-3 #19: đường là mặt, không phải rãnh). | **P2** · can_ve |
| 9 | `map_may_tile_dat_1_x4.png` | tĩnh alpha | 96×96 | may | Lòng đất 1: granite-đất có mạch tinh thể sáng ngang. | **P2** · can_ve |
| 10 | `map_may_tile_dat_2_x4.png` | tĩnh alpha | 96×96 | may | Lòng đất 2: đất núi có rễ rêu và vết nứt. | **P2** · can_ve |
| 11 | `map_may_tile_dat_giua_x4.png` | tĩnh alpha | 96×96 | may | Lòng đất giữa: lẫn khối granite to, tối hơn ~6%. | **P2** · can_ve |
| 12 | `map_may_tile_dat_day_x4.png` | tĩnh alpha | 96×96 | may | Lòng đất đáy: granite nén dải ngang, tối hơn ~12%. | **P2** · can_ve |
| 13 | `map_may_tile_vach_trai_x4.png` | tĩnh alpha | 96×96 | may | Vách trái: vách granite dựng đứng. mép TRÁI lộ thiên, mép phải nối vào lòng đất. | **P2** · can_ve |
| 14 | `map_may_tile_vach_phai_x4.png` | tĩnh alpha | 96×96 | may | Vách phải: ảnh gương của tile_vach_trai (mặt lộ thiên bên PHẢI). | **P2** · can_ve |
| 15 | `map_may_tile_chan_vach_x4.png` | tĩnh alpha | 96×96 | may | Ô chân vách: chỗ vách đứng gặp mặt đất phẳng — vật liệu tơi trượt xuống, nối mềm sang ô thường. | **P2** · can_ve |
| 16 | `map_may_tile_ke_da_x4.png` | tĩnh alpha | 96×96 | may | Kè đá: tường khối granite xếp khít, ánh lam. lặp được, dùng thay vách ở khu trung tâm. | **P2** · can_ve |
| 17 | `map_may_tile_bac_thang_x4.png` | tĩnh alpha | 96×96 | may | Một bậc thang: bậc granite. mặt bậc phẳng trên, mặt đứng bên trái, dưới-phải đặc để xếp chéo. | **P2** · can_ve |
| 18 | `map_may_tile_bac_thang_dinh_x4.png` | tĩnh alpha | 96×96 | may | Bậc trên cùng nối vào mặt đi. | **P2** · can_ve |
| 19 | `map_may_nha_nho_x4.png` | tĩnh alpha | 576×576 | may | Nhà nhỏ Kumo: trụ đá tròn, mái nón nhọn, cột thu lôi nhỏ, cửa sổ sáng lam. 6 ô. | **P2** · can_ve |
| 20 | `map_may_nha_vua_x4.png` | tĩnh alpha | 768×576 | may | Nhà vừa Kumo: tháp rộng có ban công, hoa văn mây, cửa sổ sáng. | **P2** · can_ve |
| 21 | `map_may_nha_lon_x4.png` | tĩnh alpha | 768×768 | may | Nhà lớn Kumo: tháp cao nhiều tầng kiểu tháp Raikage, chóp nhọn, nhiều cột thu lôi. | **P2** · can_ve |
| 22 | `map_may_nha_dacsan_x4.png` | tĩnh alpha | 768×576 | may | Đặc sản Mây: lò rèn kiếm / đền sét — mặt trước mở, kiếm treo, ánh điện lam. | **P2** · can_ve |
| 23 | `map_may_cong_lang_x4.png` | tĩnh alpha | 480×480 | may | Cổng làng Mây: hai trụ đá cao khắc biểu tượng sét/mây, xích sắt, phù văn phát sáng. | **P2** · can_ve |
| 24 | `map_may_cay_thong_x4.png` | tĩnh alpha | 768×768 | may | Cây thông núi cao, thân thẳng, tán tầng. | **P2** · can_ve |
| 25 | `map_may_cay_thong_nho_x4.png` | tĩnh alpha | 384×384 | may | Thông nhỏ. | **P2** · can_ve |
| 26 | `map_may_bui_ram_x4.png` | tĩnh alpha | 288×192 | may | Bụi cây núi cao lá nhỏ. | **P2** · can_ve |
| 27 | `map_may_cot_thu_loi_x4.png` | tĩnh alpha | 192×384 | may | Cột thu lôi đầu phát sáng lam. | **P2** · can_ve |
| 28 | `map_may_den_loi_x4.png` | tĩnh alpha | 192×384 | may | Đèn lồng điện treo, ánh lam tím. | **P2** · can_ve |
| 29 | `map_may_cau_treo_x4.png` | tĩnh alpha | 768×288 | may | Đoạn cầu treo dây qua vực mây — vật đặc trưng. | **P2** · can_ve |
| 30 | `map_may_tang_da_x4.png` | tĩnh alpha | 288×288 | may | Tảng granite lớn. | **P2** · can_ve |
| 31 | `map_may_hang_rao_xich_x4.png` | tĩnh alpha | 192×192 | may | Rào xích sắt giữa hai trụ đá thấp, ghép ngang được. | **P2** · can_ve |
| 32 | `map_may_bui_co_x4.png` | tĩnh alpha | 192×96 | may | Túm cỏ núi cao. | **P2** · can_ve |

### 5.7 Tổng

| Map | P1 | P2 | P3 | Cộng |
|---|---|---|---|---|
| Làng Lá (map 1) | 6 | 2 | 0 | 8 |
| Khu rừng (map 6) | 28 | 2 | 0 | 30 |
| Làng Cát (map 2) | 0 | 31 | 0 | 31 |
| Làng Sương Mù (map 3) | 0 | 33 | 0 | 33 |
| Làng Đá (map 4) | 0 | 31 | 0 | 31 |
| Làng Mây (map 5) | 0 | 32 | 0 | 32 |
| **Tổng** | **34** | **131** | 0 | **165** |

Cộng 5 concept (mục 6.0, không tính vào bảng vì không lắp vào game) và 51 ảnh map 1 đã có.

---

## 6. PROMPT TỪNG ẢNH

Mỗi khối là prompt tiếng Anh **đầy đủ, tự đứng được** (đã gồm phong cách khoá + bảng màu map + nội dung riêng + cỡ + nền + góc nhìn + 'do not draw'). Dòng `OUTPUT:` ngay dưới là tên file, cỡ, số khung và nơi lưu. Khối trong `05-map.json` (`prompt_day_du`) trùng từng chữ với khối ở đây.

### 6.0 Concept mỗi map — tham chiếu, KHÔNG lắp vào game (làm TRƯỚC bộ ảnh của map đó)

Lộ trình 4 map bước 1–2: vẽ **một bức concept**, duyệt bằng mắt, rồi `python D:\NRO_GoiBanGiao\tools\map-lang-la\do-map.py <concept.png>` lấy 7 con số làm đích. Làng Lá đã có `concept-map1-langla.png`. Cỡ gợi ý 2048×1024, nền đặc, không cần alpha, không cần đúng bội 96 (không lắp). Trong JSON nằm ở khoá riêng `concept[]`, không nằm trong `anh[]`.

**concept_khu_rung.png** — Concept Khu rừng (map 6): rừng rậm KHÔNG NHÀ, cây cổ thụ, lối đá rêu, gò đất cao 1 ô, núi rừng xa, sương. Đo 7 chỉ số trước khi vẽ bộ 30 ảnh.

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a dense wild forest clearing used as a training and hunting ground for young ninja: giant ancient trees with buttress roots, a mossy old flagstone path on grassy ground, one-cell-high earth embankments with exposed roots, a fallen hollow log, mossy boulders, a small wooden footbridge, forested mountain ridges in mist behind, NO houses or buildings at all. one wide horizontal painting, straight side view like a side-scrolling game screen (no 3/4, no top-down), same brush quality, color depth and level of detail as the approved concept-map1-langla.png, COMPOSITION TARGETS (measured later by do-map.py): sky at most 30% of the height, the mountain/background mass at least 30% of the height and taller than every rooftop, buildings packed shoulder to shoulder with no sky gaps between them, the bottom quarter is layered underground soil/rock with readable veins, roots and stones (never a flat brown band), dark foreground leaves or grass at the top and bottom edges. Do not draw: characters, monsters, text, UI, watermark. The sky/haze color of this map should read as RGB(52, 88, 72). OUTPUT SIZE: about 2048x1024 px, solid opaque PNG.
```
OUTPUT: concept_khu_rung.png · ~2048×1024 (tự do) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**concept_lang_cat.png** — Concept Làng Cát (map 2): sa mạc, nhà đất nung mái vòm, tường thành nhiều lớp (bù dải giữa vì không có cây), vách sa thạch bao quanh, cổng khắc gió.

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a desert ninja village at the foot of huge layered sandstone cliffs: round mud-brick houses with clay domes, flat roofs with striped awnings, multiple layers of crenellated city walls closing the middle band, a carved sandstone gate, date palms and cacti, sand ripples on the ground, packed ochre earth with sandstone strata underground, dust haze. one wide horizontal painting, straight side view like a side-scrolling game screen (no 3/4, no top-down), same brush quality, color depth and level of detail as the approved concept-map1-langla.png, COMPOSITION TARGETS (measured later by do-map.py): sky at most 30% of the height, the mountain/background mass at least 30% of the height and taller than every rooftop, buildings packed shoulder to shoulder with no sky gaps between them, the bottom quarter is layered underground soil/rock with readable veins, roots and stones (never a flat brown band), dark foreground leaves or grass at the top and bottom edges. Do not draw: characters, monsters, text, UI, watermark. The sky/haze color of this map should read as RGB(226, 190, 120). OUTPUT SIZE: about 2048x1024 px, solid opaque PNG.
```
OUTPUT: concept_lang_cat.png · ~2048×1024 (tự do) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**concept_lang_suong_mu.png** — Concept Làng Sương Mù (map 3): làng bến nước trong sương dày, tháp đá mái dốc, cầu gỗ qua nước, đảo đá phía sau. Sương làm các lớp cùng tông — sẽ phải ĐO chỗ giáp lớp.

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a harbor ninja village in thick fog: cylindrical stone towers with steep dark blue-grey tiled roofs, dock houses on stilts, a long wooden footbridge over teal water, mooring posts and nets, weeping willows and reeds, jagged island rock pinnacles rising out of a misty sea behind, wet dark soil with grey stone veins underground. one wide horizontal painting, straight side view like a side-scrolling game screen (no 3/4, no top-down), same brush quality, color depth and level of detail as the approved concept-map1-langla.png, COMPOSITION TARGETS (measured later by do-map.py): sky at most 30% of the height, the mountain/background mass at least 30% of the height and taller than every rooftop, buildings packed shoulder to shoulder with no sky gaps between them, the bottom quarter is layered underground soil/rock with readable veins, roots and stones (never a flat brown band), dark foreground leaves or grass at the top and bottom edges. Do not draw: characters, monsters, text, UI, watermark. The sky/haze color of this map should read as RGB(150, 172, 184). OUTPUT SIZE: about 2048x1024 px, solid opaque PNG.
```
OUTPUT: concept_lang_suong_mu.png · ~2048×1024 (tự do) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**concept_lang_da.png** — Concept Làng Đá (map 4): hẻm núi vách đá dựng, nhà đục trong đá, nhiều tầng địa hình (bố cục nền cao phức tạp — tính kỹ cột), lò rèn, cổng đá nguyên khối.

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a ninja village carved into a canyon: dwellings cut directly into layered rust-red and grey rock walls, stacked stone block houses with heavy slab roofs, stone stairs and ledges at several heights, a glowing forge, torches, a gate of two rough rock monoliths, towering strata cliffs and rock spires behind, layered rock-earth underground. one wide horizontal painting, straight side view like a side-scrolling game screen (no 3/4, no top-down), same brush quality, color depth and level of detail as the approved concept-map1-langla.png, COMPOSITION TARGETS (measured later by do-map.py): sky at most 30% of the height, the mountain/background mass at least 30% of the height and taller than every rooftop, buildings packed shoulder to shoulder with no sky gaps between them, the bottom quarter is layered underground soil/rock with readable veins, roots and stones (never a flat brown band), dark foreground leaves or grass at the top and bottom edges. Do not draw: characters, monsters, text, UI, watermark. The sky/haze color of this map should read as RGB(142, 140, 132). OUTPUT SIZE: about 2048x1024 px, solid opaque PNG.
```
OUTPUT: concept_lang_da.png · ~2048×1024 (tự do) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**concept_lang_may.png** — Concept Làng Mây (map 5): làng trên đỉnh núi cao chọc mây, tháp trụ tròn mái nhọn, cầu treo, cột thu lôi, biển mây — KHÔNG có 'mặt đất' theo nghĩa thường, xem lại quy ước va chạm trước khi vẽ bộ.

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a ninja village built on sharp granite peaks above a sea of clouds at twilight: cylindrical stone towers with pointed roofs and glowing lightning rods, rope suspension bridges between peaks, alpine pines, a gate of tall carved granite pillars with iron chains, violet storm clouds with distant lightning behind, blue-grey granite with glowing crystal veins underground. one wide horizontal painting, straight side view like a side-scrolling game screen (no 3/4, no top-down), same brush quality, color depth and level of detail as the approved concept-map1-langla.png, COMPOSITION TARGETS (measured later by do-map.py): sky at most 30% of the height, the mountain/background mass at least 30% of the height and taller than every rooftop, buildings packed shoulder to shoulder with no sky gaps between them, the bottom quarter is layered underground soil/rock with readable veins, roots and stones (never a flat brown band), dark foreground leaves or grass at the top and bottom edges. Do not draw: characters, monsters, text, UI, watermark. The sky/haze color of this map should read as RGB(120, 110, 160). OUTPUT SIZE: about 2048x1024 px, solid opaque PNG.
```
OUTPUT: concept_lang_may.png · ~2048×1024 (tự do) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

### 6.1 Làng Lá (map 1)

**1. `nen_nuihokage_rong2_x4.png`** — P1 · ve_lai

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Leaf palette: warm orange-red clay roof tiles, aged brown timber, cream plaster walls, fresh green grass, warm brown soil, clear bright blue sky. SUBJECT: a wide distant mountain range whose central cliff is carved with FIVE different giant stone faces of past village leaders (each face a different hairstyle and expression, one wearing a forehead protector), weathered grey-brown rock with green shrubs on the ledges, seen from the village below. VERTICAL LAYOUT (top to bottom): 0-25% blue sky with soft cumulus clouds; 25-80% the solid rock mass with the five carved faces — the mountain must fill at least 55% of the image height so it towers over any rooftop; 80-100% the rock foot dissolves into haze and the flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(10, 86, 154) so that the very last pixel row is exactly RGB(10, 86, 154) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. Keep the same faces, colors and mood as the approved 1536-wide version nen_nuihokage_rong_x4.png but paint the range wider on both sides with new rock, do NOT duplicate or mirror that image. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x768 px, solid opaque PNG.
```
OUTPUT: nen_nuihokage_rong2_x4.png · 2976×768 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**2. `nen_lang_xa_rong2_x4.png`** — P1 · ve_lai

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Leaf palette: warm orange-red clay roof tiles, aged brown timber, cream plaster walls, fresh green grass, warm brown soil, clear bright blue sky. SUBJECT: a dense distant ninja village seen from ground level: rows of orange-red and dark grey tiled rooftops, small towers, a few temple roofs, green trees squeezed between houses, everything softened and desaturated by distance haze. VERTICAL LAYOUT (top to bottom): 0-30% rooftops and treetops starting at the VERY FIRST pixel row — absolutely no sky at the top; 30-80% dense village with walls, balconies and trees; 80-100% the lower village dissolves into haze and the flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(10, 86, 154) so that the very last pixel row is exactly RGB(10, 86, 154) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. Same mood and roof colors as the approved 1536-wide nen_lang_xa_rong_x4.png but painted anew, wider, not duplicated. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x576 px, solid opaque PNG.
```
OUTPUT: nen_lang_xa_rong2_x4.png · 2976×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**3. `tile_dat_van_1_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Leaf palette: warm orange-red clay roof tiles, aged brown timber, cream plaster walls, fresh green grass, warm brown soil, clear bright blue sky. SUBJECT: a hand-painted underground soil tile for a side-scrolling ninja village: warm brown soil with one clearly readable horizontal vein of pale grey rock running across the middle and 3-4 small rounded pebbles, subtle soil grain elsewhere. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: same brightness as tile_dat_lienmach_1_x4.png (within 8%); the four soil tiles of this set must differ in TEXTURE, not in brightness, otherwise the ground turns into a checkerboard. The texture must be readable at real size (a viewer should think 'ah, a rock vein / roots' not 'a brown square'). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: tile_dat_van_1_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**4. `tile_dat_van_2_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Leaf palette: warm orange-red clay roof tiles, aged brown timber, cream plaster walls, fresh green grass, warm brown soil, clear bright blue sky. SUBJECT: a hand-painted underground soil tile for a side-scrolling ninja village: warm brown soil with two or three thick tree roots entering from the top edge and winding down, root tips fading into the soil, a few root hairs. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: same brightness as tile_dat_lienmach_1_x4.png (within 8%); the four soil tiles of this set must differ in TEXTURE, not in brightness, otherwise the ground turns into a checkerboard. The texture must be readable at real size (a viewer should think 'ah, a rock vein / roots' not 'a brown square'). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: tile_dat_van_2_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**5. `tile_dat_van_3_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Leaf palette: warm orange-red clay roof tiles, aged brown timber, cream plaster walls, fresh green grass, warm brown soil, clear bright blue sky. SUBJECT: a hand-painted underground soil tile for a side-scrolling ninja village: warm brown soil packed with three or four large rounded cobblestones (each 20-30 px wide) in grey and ochre, soil filling the gaps. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: about 6% darker than tile_dat_lienmach_1_x4.png; the four soil tiles of this set must differ in TEXTURE, not in brightness, otherwise the ground turns into a checkerboard. The texture must be readable at real size (a viewer should think 'ah, a rock vein / roots' not 'a brown square'). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: tile_dat_van_3_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**6. `tile_dat_van_4_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Leaf palette: warm orange-red clay roof tiles, aged brown timber, cream plaster walls, fresh green grass, warm brown soil, clear bright blue sky. SUBJECT: a hand-painted underground soil tile for a side-scrolling ninja village: dense compacted soil with faint thin horizontal compression bands and very few tiny stones — the calmest of the four soil tiles, used to rest the eye between busier tiles. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: about 12% darker than tile_dat_lienmach_1_x4.png; the four soil tiles of this set must differ in TEXTURE, not in brightness, otherwise the ground turns into a checkerboard. The texture must be readable at real size (a viewer should think 'ah, a rock vein / roots' not 'a brown square'). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: tile_dat_van_4_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**7. `tile_datlat_mat_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Leaf palette: warm orange-red clay roof tiles, aged brown timber, cream plaster walls, fresh green grass, warm brown soil, clear bright blue sky. SUBJECT: a ground-surface tile for a village path: the top 30% is a flat paved path of worn beige-grey flagstones with thin moss lines between them and an irregular lower edge, the lower 70% is the same warm brown soil with small roots as tile_dat_lienmach_1_x4.png so the tile sits on top of the soil column. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Only the left/right edges must tile; the bottom edge must continue into the soil tiles. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: tile_datlat_mat_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**8. `tt_co_phu_chan_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Leaf palette: warm orange-red clay roof tiles, aged brown timber, cream plaster walls, fresh green grass, warm brown soil, clear bright blue sky. SUBJECT: a long thin strip of wild meadow grass tufts with a few tiny white and yellow flowers, small pebbles and fallen leaves; the top silhouette is irregular and lively (tufts of different heights), the bottom edge is perfectly flush with the canvas bottom; the strip is about 40-70 px tall and fills the whole width. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Meant to be drawn IN FRONT of the feet of houses and props, so keep it low and airy. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 384x96 px, PNG with fully transparent background.
```
OUTPUT: tt_co_phu_chan_x4.png · 384×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

### 6.2 Khu rừng (map 6)

**1. `map_rung_nen_troi_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a calm blue sky above a deep forest, soft white cumulus clouds and a few high thin clouds, warm late-morning light. VERTICAL LAYOUT (top to bottom): 0-30% blue sky with clouds; 30-80% the sky pales toward a soft green-grey forest haze; 80-100% flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(52, 88, 72) so that the very last pixel row is exactly RGB(52, 88, 72) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. No trees, no ground, no sun disc. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x384 px, solid opaque PNG.
```
OUTPUT: map_rung_nen_troi_x4.png · 2976×384 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**2. `map_rung_nen_nui_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: two or three overlapping ridges of forested mountains completely covered with dense dark-green canopy, morning mist pooling between the ridges, a few pale rock outcrops, varied silhouette with no repeating shapes. VERTICAL LAYOUT (top to bottom): 0-25% sky with clouds; 25-80% the mountain ridges — the mass must fill at least 55% of the height; 80-100% the lowest ridge dissolves into mist and the flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(52, 88, 72) so that the very last pixel row is exactly RGB(52, 88, 72) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. No buildings anywhere — this is wild forest. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x768 px, solid opaque PNG.
```
OUTPUT: map_rung_nen_nui_x4.png · 2976×768 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**3. `map_rung_nen_langxa_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a wall of giant forest tree crowns seen from inside the forest, overlapping canopies of many greens, thick trunks and hanging vines visible in the lower half, all softened by distance haze. VERTICAL LAYOUT (top to bottom): 0-30% tree crowns starting at the VERY FIRST pixel row — absolutely no sky at the top; 30-80% dense canopy and trunks; 80-100% the undergrowth dissolves into mist and the flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(52, 88, 72) so that the very last pixel row is exactly RGB(52, 88, 72) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. No houses, no paths, no characters. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x576 px, solid opaque PNG.
```
OUTPUT: map_rung_nen_langxa_x4.png · 2976×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**4. `map_rung_nen_gan_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: the nearest background layer of a misty forest: blurred dark tree trunks, fern silhouettes and low branches half hidden in soft mist. VERTICAL LAYOUT (top to bottom): 0-30% blurred trunks and branches; 30-80% dense mist with faint ferns; 80-100% flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(52, 88, 72) so that the very last pixel row is exactly RGB(52, 88, 72) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. Soft, low contrast — it sits right behind the playable ground. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x288 px, solid opaque PNG.
```
OUTPUT: map_rung_nen_gan_x4.png · 2976×288 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**5. `map_rung_tile_co_mat_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a hand-painted forest ground tile for a side-scrolling ninja game: the top ~35% is lush dark-green forest grass with an uneven ragged edge, a few fallen brown leaves and moss patches, the lower 65% is dark forest soil with fine roots (must continue into the soil tiles below). Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the top area above the grass edge is open air and transparent.
```
OUTPUT: map_rung_tile_co_mat_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**6. `map_rung_tile_co_sang_da_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a hand-painted forest ground tile for a side-scrolling ninja game: left half forest grass surface, right half an old mossy flagstone path surface, the boundary irregular with grass tufts creeping over the stones; lower 65% forest soil. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the top area above the surface edge is open air and transparent.
```
OUTPUT: map_rung_tile_co_sang_da_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**7. `map_rung_tile_da_sang_co_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a hand-painted forest ground tile for a side-scrolling ninja game: left half old mossy flagstone path surface, right half forest grass, irregular boundary with grass creeping over the stones; lower 65% forest soil. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the top area above the surface edge is open air and transparent.
```
OUTPUT: map_rung_tile_da_sang_co_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**8. `map_rung_tile_datlat_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a hand-painted forest ground tile for a side-scrolling ninja game: the top 30% is an old worn flagstone path of grey-green stones with moss in the cracks and an irregular lower edge, the lower 70% is dark forest soil with roots. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four soil tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_rung_tile_datlat_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**9. `map_rung_tile_dat_1_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a hand-painted forest ground tile for a side-scrolling ninja game: dark forest soil with one readable horizontal vein of pale rock and a few small pebbles. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four soil tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_rung_tile_dat_1_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**10. `map_rung_tile_dat_2_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a hand-painted forest ground tile for a side-scrolling ninja game: dark forest soil with several thick winding tree roots coming down from the top edge, root hairs, a hint of moss. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four soil tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_rung_tile_dat_2_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**11. `map_rung_tile_dat_giua_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a hand-painted forest ground tile for a side-scrolling ninja game: dark forest soil packed with large rounded cobblestones, about 6% darker than tile_dat_1. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four soil tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_rung_tile_dat_giua_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**12. `map_rung_tile_dat_day_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a hand-painted forest ground tile for a side-scrolling ninja game: dense compacted deep soil with faint horizontal bands and almost no stones, about 12% darker than tile_dat_1. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four soil tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_rung_tile_dat_day_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**13. `map_rung_tile_vach_trai_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a hand-painted forest ground tile for a side-scrolling ninja game: the vertical cut face of an earth embankment: the LEFT edge is the exposed face (roots poking out, small stones, a lip of grass hanging over the top-left corner), the right side continues into solid soil. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_rung_tile_vach_trai_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**14. `map_rung_tile_vach_phai_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a hand-painted forest ground tile for a side-scrolling ninja game: the vertical cut face of an earth embankment: the RIGHT edge is the exposed face with roots and small stones and a lip of grass at the top-right corner, the left side continues into solid soil. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_rung_tile_vach_phai_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**15. `map_rung_tile_chan_vach_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a hand-painted forest ground tile for a side-scrolling ninja game: the foot of an embankment where the vertical face meets flat ground: a small slope of loose crumbled soil, a few pebbles and one grass tuft, blending softly into plain soil on the other side. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_rung_tile_chan_vach_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**16. `map_rung_tile_ke_da_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a hand-painted forest ground tile for a side-scrolling ninja game: a dry-stacked retaining wall of mossy grey forest stones of varied sizes, damp shadows in the joints. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_rung_tile_ke_da_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**17. `map_rung_tile_bac_thang_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a hand-painted forest ground tile for a side-scrolling ninja game: a single old mossy stone stair block seen from the side: a flat tread along the top, a vertical riser on the left side, solid stone and soil filling the lower-right so copies stack diagonally into a staircase. Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the area above the tread and left of the riser is open air and transparent.
```
OUTPUT: map_rung_tile_bac_thang_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**18. `map_rung_tile_bac_thang_dinh_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a hand-painted forest ground tile for a side-scrolling ninja game: the topmost stair block: a mossy stone tread that merges on its right side into forest grass and soil (matching tile_co_mat). Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the area above the tread is open air and transparent.
```
OUTPUT: map_rung_tile_bac_thang_dinh_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**19. `map_rung_cay_co_thu_1_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a giant ancient forest tree, massive gnarled trunk with buttress roots spreading at the base, rough bark with moss, a broad dense multi-layered dark-green crown; 8 cells tall = the full 768 px. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x768 px, PNG with fully transparent background.
```
OUTPUT: map_rung_cay_co_thu_1_x4.png · 768×768 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**20. `map_rung_cay_co_thu_2_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a second giant forest tree with a clearly different silhouette: a leaning twisted trunk with a dark hollow, vines hanging from the branches, an asymmetric crown that leans to one side; 8 cells tall = the full 768 px. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x768 px, PNG with fully transparent background.
```
OUTPUT: map_rung_cay_co_thu_2_x4.png · 768×768 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**21. `map_rung_cay_nho_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a young forest tree with a slender trunk and a round leafy crown, a few leaves falling; 4 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 384x384 px, PNG with fully transparent background.
```
OUTPUT: map_rung_cay_nho_x4.png · 384×384 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**22. `map_rung_bui_ram_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a low dense clump of ferns and bramble with a few red berries; 2 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 288x192 px, PNG with fully transparent background.
```
OUTPUT: map_rung_bui_ram_x4.png · 288×192 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**23. `map_rung_goc_cay_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a cut tree stump showing growth rings on top, moss on one side and a cluster of small brown mushrooms at its base; 2 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 288x192 px, PNG with fully transparent background.
```
OUTPUT: map_rung_goc_cay_x4.png · 288×192 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**24. `map_rung_tang_da_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a cluster of two or three mossy grey boulders of different sizes, damp cracks, a fern growing in between; 2 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 288x192 px, PNG with fully transparent background.
```
OUTPUT: map_rung_tang_da_x4.png · 288×192 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**25. `map_rung_cay_do_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a fallen hollow tree trunk lying horizontally, bark peeling, moss and small fungi growing on it, broken branch stubs; 6 cells wide, 2 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 576x192 px, PNG with fully transparent background.
```
OUTPUT: map_rung_cay_do_x4.png · 576×192 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**26. `map_rung_cau_go_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a small rustic wooden footbridge with weathered planks, two low posts and rope railings, moss on the posts; 5 cells wide, 2 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 480x192 px, PNG with fully transparent background.
```
OUTPUT: map_rung_cau_go_x4.png · 480×192 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**27. `map_rung_bien_chi_duong_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a wooden signpost with two arrow-shaped boards pointing different ways, carved marks only (no readable letters), rope lashing, moss at the base; 3 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x288 px, PNG with fully transparent background.
```
OUTPUT: map_rung_bien_chi_duong_x4.png · 192×288 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**28. `map_rung_bui_co_x4.png`** — P1 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a tuft of wild forest grass with a small fern and a few tiny white flowers; 1 cell tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x96 px, PNG with fully transparent background.
```
OUTPUT: map_rung_bui_co_x4.png · 192×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**29. `map_rung_tc_canh_la_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: dark foreground tree branches with leaves hanging DOWN from the top edge of the canvas, very dark green almost black silhouette with a few rim-lit leaves, the branches attached to the top edge, the lower part of the canvas empty and transparent; 8 cells wide. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. This is a foreground overlay, so the object touches the TOP edge instead of the bottom edge. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x384 px, PNG with fully transparent background.
```
OUTPUT: map_rung_tc_canh_la_x4.png · 768×384 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**30. `map_rung_tc_bui_truoc_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Deep forest palette: dark emerald and moss greens, mottled olive light, damp brown bark and soil, cool blue-green mist in the distance, dappled warm sunlight. SUBJECT: a tall dark clump of ferns and grass in the foreground, very dark green almost black silhouette with a few rim-lit edges, growing up from the bottom edge; 4 cells wide. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 384x192 px, PNG with fully transparent background.
```
OUTPUT: map_rung_tc_bui_truoc_x4.png · 384×192 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

### 6.3 Làng Cát (map 2)

**1. `map_cat_nen_troi_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a hazy sun-bleached desert sky, pale yellow-orange near the horizon, a few thin wispy clouds and a faint veil of blowing dust low in the frame. VERTICAL LAYOUT (top to bottom): 0-30% sky with clouds; 30-80% the sky pales toward the haze color; 80-100% flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(226, 190, 120) so that the very last pixel row is exactly RGB(226, 190, 120) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. No ground, no buildings. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x384 px, solid opaque PNG.
```
OUTPUT: map_cat_nen_troi_x4.png · 2976×384 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**2. `map_cat_nen_nui_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: massive layered sandstone cliffs and flat-topped mesas encircling a desert village, ochre and terracotta strata, a few distant sand dunes at their feet, heat haze. VERTICAL LAYOUT (top to bottom): 0-25% sky with clouds; 25-80% the solid mountain mass — it must fill at least 55% of the image height so it towers over every rooftop; 80-100% the foot dissolves into haze and the flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(226, 190, 120) so that the very last pixel row is exactly RGB(226, 190, 120) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x768 px, solid opaque PNG.
```
OUTPUT: map_cat_nen_nui_x4.png · 2976×768 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**3. `map_cat_nen_langxa_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a dense desert ninja village seen from ground level: rounded mud-brick domes, cylindrical clay towers with small dark windows, flat roofs with striped cloth awnings, multiple layers of mud-brick city walls, all softened by dust haze. VERTICAL LAYOUT (top to bottom): 0-30% rooftops starting at the VERY FIRST pixel row — absolutely no sky at the top; 30-80% the dense village; 80-100% the lower village dissolves into haze and the flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(226, 190, 120) so that the very last pixel row is exactly RGB(226, 190, 120) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x576 px, solid opaque PNG.
```
OUTPUT: map_cat_nen_langxa_x4.png · 2976×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**4. `map_cat_nen_gan_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: the nearest background layer of a desert village: blurred tops of mud-brick walls and a low drifting dust haze. VERTICAL LAYOUT (top to bottom): 0-30% blurred near shapes; 30-80% dense haze with faint details; 80-100% flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(226, 190, 120) so that the very last pixel row is exactly RGB(226, 190, 120) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. Soft and low contrast — it sits right behind the playable ground. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x288 px, solid opaque PNG.
```
OUTPUT: map_cat_nen_gan_x4.png · 2976×288 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**5. `map_cat_tile_cat_mat_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a ground-surface tile: the top ~35% is pale rippled desert sand with wind ripples and a few tiny pebbles with an uneven natural edge, the lower 65% is packed sun-baked desert earth in warm ochre (must continue into the underground tiles) Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the top area above the surface edge is open air and transparent.
```
OUTPUT: map_cat_tile_cat_mat_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**6. `map_cat_tile_cat_sang_da_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a transition ground tile: left half pale rippled desert sand with wind ripples and a few tiny pebbles, right half flat pinkish-beige sandstone flag paving with fine sand in the joints, the boundary irregular and natural with the surface material creeping over the paving; lower 65% packed sun-baked desert earth in warm ochre Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the top area above the surface edge is open air and transparent.
```
OUTPUT: map_cat_tile_cat_sang_da_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**7. `map_cat_tile_da_sang_cat_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a transition ground tile: left half flat pinkish-beige sandstone flag paving with fine sand in the joints, right half pale rippled desert sand with wind ripples and a few tiny pebbles, irregular natural boundary; lower 65% packed sun-baked desert earth in warm ochre Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the top area above the surface edge is open air and transparent.
```
OUTPUT: map_cat_tile_da_sang_cat_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**8. `map_cat_tile_datlat_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a paved-path tile: the top 30% is flat pinkish-beige sandstone flag paving with fine sand in the joints with an irregular lower edge, the lower 70% is packed sun-baked desert earth in warm ochre Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_cat_tile_datlat_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**9. `map_cat_tile_dat_1_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: an underground tile: packed sun-baked desert earth in warm ochre with one horizontal band of layered sandstone and a few pebbles Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four underground tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_cat_tile_dat_1_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**10. `map_cat_tile_dat_2_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: an underground tile: packed sun-baked desert earth in warm ochre with sun-crack lines and two or three dead dry roots Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four underground tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_cat_tile_dat_2_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**11. `map_cat_tile_dat_giua_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: an underground tile: packed sun-baked desert earth in warm ochre with embedded large chunks of sandstone, about 6% darker than tile_dat_1 Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four underground tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_cat_tile_dat_giua_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**12. `map_cat_tile_dat_day_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: an underground tile: packed sun-baked desert earth in warm ochre with dense compressed sediment bands and almost no stones, about 12% darker than tile_dat_1 Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four underground tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_cat_tile_dat_day_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**13. `map_cat_tile_vach_trai_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: the vertical cut face of a sand-and-earth embankment with thin sandstone layers: the LEFT edge is the exposed face (crumbling sand, exposed sandstone ledges, a lip of sand at the top corner), the right side continues into solid packed sun-baked desert earth in warm ochre Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_cat_tile_vach_trai_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**14. `map_cat_tile_vach_phai_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: the vertical cut face of a sand-and-earth embankment with thin sandstone layers: the RIGHT edge is the exposed face (crumbling sand, exposed sandstone ledges, a lip of sand at the top corner), the left side continues into solid packed sun-baked desert earth in warm ochre Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_cat_tile_vach_phai_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**15. `map_cat_tile_chan_vach_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: the foot of a sand-and-earth embankment with thin sandstone layers where the vertical face meets flat ground: a small slope of loose crumbled material, a few small stones, blending softly into plain packed sun-baked desert earth in warm ochre Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_cat_tile_chan_vach_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**16. `map_cat_tile_ke_da_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a retaining wall tile: stacked wind-worn sandstone blocks of varied size, sand collected in the joints Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_cat_tile_ke_da_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**17. `map_cat_tile_bac_thang_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a single stair block seen from the side: a flat tread of worn pinkish sandstone along the top, a vertical riser on the left, solid material filling the lower-right so copies stack diagonally into a staircase Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the area above the tread and left of the riser is open air and transparent.
```
OUTPUT: map_cat_tile_bac_thang_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**18. `map_cat_tile_bac_thang_dinh_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: the topmost stair block: a tread of worn pinkish sandstone that merges on its right side into pale rippled desert sand with wind ripples and a few tiny pebbles (matching the surface tile) Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the area above the tread is open air and transparent.
```
OUTPUT: map_cat_tile_bac_thang_dinh_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**19. `map_cat_nha_nho_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a small round mud-brick desert house with a smooth clay dome roof, a wooden door, two tiny square windows, a striped cloth awning over the door and a clay water jar beside it; 6 cells tall = 576 px. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 576x576 px, PNG with fully transparent background.
```
OUTPUT: map_cat_nha_nho_x4.png · 576×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**20. `map_cat_nha_vua_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a wider flat-roofed desert clay house with a low parapet, an arched doorway, a long striped awning, wooden beams poking out of the walls and water jars on the roof; 6 cells tall, 8 cells wide. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x576 px, PNG with fully transparent background.
```
OUTPUT: map_cat_nha_vua_x4.png · 768×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**21. `map_cat_nha_lon_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a tall two-tier desert clay tower house topped by a dome with wind vents, an external clay staircase, small arched windows, faded teal painted trim; 8 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x768 px, PNG with fully transparent background.
```
OUTPUT: map_cat_nha_lon_x4.png · 768×768 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**22. `map_cat_nha_dacsan_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a desert bazaar stall and puppet workshop: an open front under a wide cloth canopy, hanging fabrics, shelves of clay pots, a wooden marionette hanging from a beam, wind chimes; 6 cells tall, 8 cells wide. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x576 px, PNG with fully transparent background.
```
OUTPUT: map_cat_nha_dacsan_x4.png · 768×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**23. `map_cat_cong_lang_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: the village gate of the Hidden Sand: two massive carved sandstone pillars with a large carved wind-swirl symbol, a heavy wooden crossbeam with hanging cloth strips, sand piled at the feet; 5 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 480x480 px, PNG with fully transparent background.
```
OUTPUT: map_cat_cong_lang_x4.png · 480×480 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**24. `map_cat_cay_cha_la_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a tall date palm with a ringed slender trunk and a spread of long arching fronds, a cluster of dates; 8 cells tall, 4 cells wide. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 384x768 px, PNG with fully transparent background.
```
OUTPUT: map_cat_cay_cha_la_x4.png · 384×768 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**25. `map_cat_xuong_rong_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a cluster of tall desert cacti of different heights with spines and one small red flower; 3 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x288 px, PNG with fully transparent background.
```
OUTPUT: map_cat_xuong_rong_x4.png · 192×288 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**26. `map_cat_bui_gai_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a dry thorny desert bush, brittle brown-grey twigs with a few dusty green leaves; 2 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 288x192 px, PNG with fully transparent background.
```
OUTPUT: map_cat_bui_gai_x4.png · 288×192 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**27. `map_cat_gieng_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a desert stone well with a wooden frame, a small cloth sunshade and a leather bucket on a rope; 3 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 288x288 px, PNG with fully transparent background.
```
OUTPUT: map_cat_gieng_x4.png · 288×288 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**28. `map_cat_den_dau_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a desert oil lamp post: a stone base, a wooden pole and a brass lamp with warm glowing glass; 4 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x384 px, PNG with fully transparent background.
```
OUTPUT: map_cat_den_dau_x4.png · 192×384 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**29. `map_cat_binh_gom_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a cluster of two or three clay water jars of different sizes with painted stripe patterns, one with a cloth lid; 2 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x192 px, PNG with fully transparent background.
```
OUTPUT: map_cat_binh_gom_x4.png · 192×192 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**30. `map_cat_tuong_thanh_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a segment of sun-baked mud-brick city wall with crenellations on top, weathered cracks and a small arched niche; the left and right edges must join seamlessly with a copy of itself; 4 cells tall, 8 cells wide. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x384 px, PNG with fully transparent background.
```
OUTPUT: map_cat_tuong_thanh_x4.png · 768×384 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**31. `map_cat_bui_co_kho_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Sand (Wind) palette: warm sand yellow, ochre and terracotta sun-baked clay, pale beige, dusty faded teal accents, hazy sun-bleached sky. SUBJECT: a tuft of dry yellow desert grass with a little sand drift at its base; 1 cell tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x96 px, PNG with fully transparent background.
```
OUTPUT: map_cat_bui_co_kho_x4.png · 192×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

### 6.4 Làng Sương Mù (map 3)

**1. `map_smu_nen_troi_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a cold overcast grey-blue sky with low horizontal bands of fog and soft diffused light. VERTICAL LAYOUT (top to bottom): 0-30% sky with clouds; 30-80% the sky pales toward the haze color; 80-100% flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(150, 172, 184) so that the very last pixel row is exactly RGB(150, 172, 184) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. No ground, no buildings. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x384 px, solid opaque PNG.
```
OUTPUT: map_smu_nen_troi_x4.png · 2976×384 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**2. `map_smu_nen_nui_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: jagged dark island rock pinnacles rising out of a misty sea, several overlapping layers of wet slate-grey cliffs with moss and stunted pines, fog pooling between them. VERTICAL LAYOUT (top to bottom): 0-25% sky with clouds; 25-80% the solid mountain mass — it must fill at least 55% of the image height so it towers over every rooftop; 80-100% the foot dissolves into haze and the flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(150, 172, 184) so that the very last pixel row is exactly RGB(150, 172, 184) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x768 px, solid opaque PNG.
```
OUTPUT: map_smu_nen_nui_x4.png · 2976×768 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**3. `map_smu_nen_langxa_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a dense mist-village seen from ground level: steep dark blue-grey tiled roofs on cylindrical stone towers, rope bridges and pipes connecting them, lanterns glowing cold blue, everything veiled in fog. VERTICAL LAYOUT (top to bottom): 0-30% rooftops starting at the VERY FIRST pixel row — absolutely no sky at the top; 30-80% the dense village; 80-100% the lower village dissolves into haze and the flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(150, 172, 184) so that the very last pixel row is exactly RGB(150, 172, 184) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x576 px, solid opaque PNG.
```
OUTPUT: map_smu_nen_langxa_x4.png · 2976×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**4. `map_smu_nen_gan_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: the nearest background layer of a harbor village: blurred dock piles, mooring posts and hanging nets inside dense white fog. VERTICAL LAYOUT (top to bottom): 0-30% blurred near shapes; 30-80% dense haze with faint details; 80-100% flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(150, 172, 184) so that the very last pixel row is exactly RGB(150, 172, 184) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. Soft and low contrast — it sits right behind the playable ground. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x288 px, solid opaque PNG.
```
OUTPUT: map_smu_nen_gan_x4.png · 2976×288 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**5. `map_smu_tile_co_mat_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a ground-surface tile: the top ~35% is damp dark-green grass mixed with moss, wet and glistening with an uneven natural edge, the lower 65% is wet dark soil in cool brown-grey (must continue into the underground tiles) Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the top area above the surface edge is open air and transparent.
```
OUTPUT: map_smu_tile_co_mat_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**6. `map_smu_tile_co_sang_da_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a transition ground tile: left half damp dark-green grass mixed with moss, wet and glistening, right half wet dark grey-blue flagstones with moss in the joints and small puddles, the boundary irregular and natural with the surface material creeping over the paving; lower 65% wet dark soil in cool brown-grey Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the top area above the surface edge is open air and transparent.
```
OUTPUT: map_smu_tile_co_sang_da_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**7. `map_smu_tile_da_sang_co_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a transition ground tile: left half wet dark grey-blue flagstones with moss in the joints and small puddles, right half damp dark-green grass mixed with moss, wet and glistening, irregular natural boundary; lower 65% wet dark soil in cool brown-grey Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the top area above the surface edge is open air and transparent.
```
OUTPUT: map_smu_tile_da_sang_co_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**8. `map_smu_tile_datlat_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a paved-path tile: the top 30% is wet dark grey-blue flagstones with moss in the joints and small puddles with an irregular lower edge, the lower 70% is wet dark soil in cool brown-grey Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_smu_tile_datlat_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**9. `map_smu_tile_dat_1_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: an underground tile: wet dark soil in cool brown-grey with one horizontal vein of grey stone and a few wet pebbles Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four underground tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_smu_tile_dat_1_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**10. `map_smu_tile_dat_2_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: an underground tile: wet dark soil in cool brown-grey with winding roots and a thin seep of water running down Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four underground tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_smu_tile_dat_2_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**11. `map_smu_tile_dat_giua_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: an underground tile: wet dark soil in cool brown-grey with large wet glistening cobblestones, about 6% darker than tile_dat_1 Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four underground tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_smu_tile_dat_giua_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**12. `map_smu_tile_dat_day_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: an underground tile: wet dark soil in cool brown-grey with dense compressed clay bands and almost no stones, about 12% darker than tile_dat_1 Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four underground tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_smu_tile_dat_day_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**13. `map_smu_tile_vach_trai_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: the vertical cut face of a wet earth embankment: the LEFT edge is the exposed face (dripping moss, exposed roots, a lip of damp grass at the top corner), the right side continues into solid wet dark soil in cool brown-grey Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_smu_tile_vach_trai_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**14. `map_smu_tile_vach_phai_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: the vertical cut face of a wet earth embankment: the RIGHT edge is the exposed face (dripping moss, exposed roots, a lip of damp grass at the top corner), the left side continues into solid wet dark soil in cool brown-grey Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_smu_tile_vach_phai_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**15. `map_smu_tile_chan_vach_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: the foot of a wet earth embankment where the vertical face meets flat ground: a small slope of loose crumbled material, a few small stones, blending softly into plain wet dark soil in cool brown-grey Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_smu_tile_chan_vach_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**16. `map_smu_tile_ke_da_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a retaining wall tile: stacked wet mossy grey stones with water seeping through the joints Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_smu_tile_ke_da_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**17. `map_smu_tile_bac_thang_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a single stair block seen from the side: a flat tread of wet mossy grey stone along the top, a vertical riser on the left, solid material filling the lower-right so copies stack diagonally into a staircase Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the area above the tread and left of the riser is open air and transparent.
```
OUTPUT: map_smu_tile_bac_thang_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**18. `map_smu_tile_bac_thang_dinh_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: the topmost stair block: a tread of wet mossy grey stone that merges on its right side into damp dark-green grass mixed with moss, wet and glistening (matching the surface tile) Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the area above the tread is open air and transparent.
```
OUTPUT: map_smu_tile_bac_thang_dinh_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**19. `map_smu_tile_nuoc_mat_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a water tile: the top edge is the water surface with small ripples and a thin foam line, below it deep teal water with soft light streaks; must tile left/right and continue into tile_nuoc_sau below Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the thin area above the water surface is open air and transparent.
```
OUTPUT: map_smu_tile_nuoc_mat_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**20. `map_smu_tile_nuoc_sau_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a deep-water tile: dark teal water with faint wavy light streaks, tileable in all four directions Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_smu_tile_nuoc_sau_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**21. `map_smu_nha_nho_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a small mist-village house: a cylindrical mossy stone base, a steep dark blue-grey tiled roof glistening with water, a round window, a wooden door, water dripping from the eaves; 6 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 576x576 px, PNG with fully transparent background.
```
OUTPUT: map_smu_nha_nho_x4.png · 576×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**22. `map_smu_nha_vua_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a dock house standing on wooden stilts, damp plank walls, a steep tiled roof, fishing nets and buoys hanging on the side, a cold blue lantern by the door; 6 cells tall, 8 cells wide. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x576 px, PNG with fully transparent background.
```
OUTPUT: map_smu_nha_vua_x4.png · 768×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**23. `map_smu_nha_lon_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a tall mist-village tower house with two stacked steep tiled roofs, external iron pipes, a small balcony bridge stub, glowing round windows, moss climbing the stone walls; 8 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x768 px, PNG with fully transparent background.
```
OUTPUT: map_smu_nha_lon_x4.png · 768×768 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**24. `map_smu_nha_dacsan_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a fish market and boat house: an open front under a plank roof, hanging fish and squid, barrels, a small upturned rowing boat leaning on the wall; 6 cells tall, 8 cells wide. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x576 px, PNG with fully transparent background.
```
OUTPUT: map_smu_nha_dacsan_x4.png · 768×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**25. `map_smu_cong_lang_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: the village gate of the Hidden Mist: two tall wet stone pillars carved with four wavy mist lines, a dark wooden crossbeam, dripping moss and hanging bells; 5 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 480x480 px, PNG with fully transparent background.
```
OUTPUT: map_smu_cong_lang_x4.png · 480×480 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**26. `map_smu_cay_lieu_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a large weeping willow with a twisted mossy trunk and long trailing leafy branches, drops of water on the leaves; 8 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x768 px, PNG with fully transparent background.
```
OUTPUT: map_smu_cay_lieu_x4.png · 768×768 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**27. `map_smu_cay_nho_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a small damp tree with a moss-covered trunk and a dense dark round crown; 4 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 384x384 px, PNG with fully transparent background.
```
OUTPUT: map_smu_cay_nho_x4.png · 384×384 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**28. `map_smu_bui_lau_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a clump of tall reeds and cattails with brown seed heads; 2 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 288x192 px, PNG with fully transparent background.
```
OUTPUT: map_smu_bui_lau_x4.png · 288×192 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**29. `map_smu_cau_go_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a long wooden footbridge over water with weathered silver-grey planks, rope railings on wooden posts, moss and barnacles at the base of the posts; 3 cells tall, 8 cells wide. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x288 px, PNG with fully transparent background.
```
OUTPUT: map_smu_cau_go_x4.png · 768×288 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**30. `map_smu_cot_buoc_thuyen_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a wooden mooring post wrapped with thick rope, barnacles and seaweed clinging to its base; 3 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x288 px, PNG with fully transparent background.
```
OUTPUT: map_smu_cot_buoc_thuyen_x4.png · 192×288 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**31. `map_smu_den_hieu_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a harbor lantern post: a stone base, an iron pole and a glass lantern glowing cold blue-white in the fog; 4 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x384 px, PNG with fully transparent background.
```
OUTPUT: map_smu_den_hieu_x4.png · 192×384 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**32. `map_smu_thung_ca_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: two wooden fish barrels with a draped fishing net and a couple of cork floats; 2 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x192 px, PNG with fully transparent background.
```
OUTPUT: map_smu_thung_ca_x4.png · 192×192 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**33. `map_smu_bui_co_uot_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Mist (Water) palette: slate blue-grey, deep teal water, wet dark wood, mossy grey stone, cold white fog, every surface damp and glistening. SUBJECT: a tuft of wet dark grass with a small shallow puddle reflecting light; 1 cell tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x96 px, PNG with fully transparent background.
```
OUTPUT: map_smu_bui_co_uot_x4.png · 192×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

### 6.5 Làng Đá (map 4)

**1. `map_da_nen_troi_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a dusty pale grey-blue sky with thin streaky clouds and a warm dust haze near the bottom. VERTICAL LAYOUT (top to bottom): 0-30% sky with clouds; 30-80% the sky pales toward the haze color; 80-100% flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(142, 140, 132) so that the very last pixel row is exactly RGB(142, 140, 132) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. No ground, no buildings. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x384 px, solid opaque PNG.
```
OUTPUT: map_da_nen_troi_x4.png · 2976×384 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**2. `map_da_nen_nui_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: towering canyon walls of layered rock strata in rust-red, brown and grey, sheer vertical faces and tall rock spires, a few hardy shrubs on ledges. VERTICAL LAYOUT (top to bottom): 0-25% sky with clouds; 25-80% the solid mountain mass — it must fill at least 55% of the image height so it towers over every rooftop; 80-100% the foot dissolves into haze and the flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(142, 140, 132) so that the very last pixel row is exactly RGB(142, 140, 132) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x768 px, solid opaque PNG.
```
OUTPUT: map_da_nen_nui_x4.png · 2976×768 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**3. `map_da_nen_langxa_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a dense stone-village seen from ground level: dwellings carved directly into the cliff, stacked block houses with small square windows, stone bridges connecting the ledges, dust haze. VERTICAL LAYOUT (top to bottom): 0-30% rooftops starting at the VERY FIRST pixel row — absolutely no sky at the top; 30-80% the dense village; 80-100% the lower village dissolves into haze and the flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(142, 140, 132) so that the very last pixel row is exactly RGB(142, 140, 132) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x576 px, solid opaque PNG.
```
OUTPUT: map_da_nen_langxa_x4.png · 2976×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**4. `map_da_nen_gan_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: the nearest background layer of a stone village: blurred rough rock blocks and boulders in a warm dust haze. VERTICAL LAYOUT (top to bottom): 0-30% blurred near shapes; 30-80% dense haze with faint details; 80-100% flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(142, 140, 132) so that the very last pixel row is exactly RGB(142, 140, 132) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. Soft and low contrast — it sits right behind the playable ground. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x288 px, solid opaque PNG.
```
OUTPUT: map_da_nen_gan_x4.png · 2976×288 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**5. `map_da_tile_da_mat_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a ground-surface tile: the top ~35% is flat worn warm-grey natural rock with cracks and two or three dry grass tufts with an uneven natural edge, the lower 65% is brown rock-earth with rust-red strata (must continue into the underground tiles) Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the top area above the surface edge is open air and transparent.
```
OUTPUT: map_da_tile_da_mat_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**6. `map_da_tile_da_sang_lat_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a transition ground tile: left half flat worn warm-grey natural rock with cracks and two or three dry grass tufts, right half cut stone paving of tightly fitted warm-grey blocks with chisel marks, the boundary irregular and natural with the surface material creeping over the paving; lower 65% brown rock-earth with rust-red strata Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the top area above the surface edge is open air and transparent.
```
OUTPUT: map_da_tile_da_sang_lat_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**7. `map_da_tile_lat_sang_da_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a transition ground tile: left half cut stone paving of tightly fitted warm-grey blocks with chisel marks, right half flat worn warm-grey natural rock with cracks and two or three dry grass tufts, irregular natural boundary; lower 65% brown rock-earth with rust-red strata Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the top area above the surface edge is open air and transparent.
```
OUTPUT: map_da_tile_lat_sang_da_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**8. `map_da_tile_datlat_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a paved-path tile: the top 30% is cut stone paving of tightly fitted warm-grey blocks with chisel marks with an irregular lower edge, the lower 70% is brown rock-earth with rust-red strata Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_da_tile_datlat_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**9. `map_da_tile_dat_1_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: an underground tile: brown rock-earth with rust-red strata with one horizontal rust-red sediment band and a few angular stones Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four underground tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_da_tile_dat_1_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**10. `map_da_tile_dat_2_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: an underground tile: brown rock-earth with rust-red strata with a pale quartz mineral vein and a few cracks Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four underground tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_da_tile_dat_2_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**11. `map_da_tile_dat_giua_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: an underground tile: brown rock-earth with rust-red strata with embedded large angular rocks, about 6% darker than tile_dat_1 Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four underground tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_da_tile_dat_giua_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**12. `map_da_tile_dat_day_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: an underground tile: brown rock-earth with rust-red strata with dense tightly packed strata bands, about 12% darker than tile_dat_1 Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four underground tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_da_tile_dat_day_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**13. `map_da_tile_vach_trai_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: the vertical cut face of a layered rock cliff: the LEFT edge is the exposed face (sheer strata, small ledges, loose stones at the top corner), the right side continues into solid brown rock-earth with rust-red strata Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_da_tile_vach_trai_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**14. `map_da_tile_vach_phai_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: the vertical cut face of a layered rock cliff: the RIGHT edge is the exposed face (sheer strata, small ledges, loose stones at the top corner), the left side continues into solid brown rock-earth with rust-red strata Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_da_tile_vach_phai_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**15. `map_da_tile_chan_vach_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: the foot of a layered rock cliff where the vertical face meets flat ground: a small slope of loose crumbled material, a few small stones, blending softly into plain brown rock-earth with rust-red strata Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_da_tile_chan_vach_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**16. `map_da_tile_ke_da_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a retaining wall tile: tightly fitted cut stone blocks of warm grey with dark mortar joints and chisel marks Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_da_tile_ke_da_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**17. `map_da_tile_bac_thang_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a single stair block seen from the side: a flat tread of rough chiselled mountain rock along the top, a vertical riser on the left, solid material filling the lower-right so copies stack diagonally into a staircase Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the area above the tread and left of the riser is open air and transparent.
```
OUTPUT: map_da_tile_bac_thang_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**18. `map_da_tile_bac_thang_dinh_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: the topmost stair block: a tread of rough chiselled mountain rock that merges on its right side into flat worn warm-grey natural rock with cracks and two or three dry grass tufts (matching the surface tile) Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the area above the tread is open air and transparent.
```
OUTPUT: map_da_tile_bac_thang_dinh_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**19. `map_da_nha_nho_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a small dwelling carved into a rough rock mass: a chiselled stone facade with a thick wooden door, a round window, an iron lantern bracket, rock strata visible on the sides; 6 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 576x576 px, PNG with fully transparent background.
```
OUTPUT: map_da_nha_nho_x4.png · 576×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**20. `map_da_nha_vua_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a stacked stone block house with a heavy flat slab roof, a stone chimney with smoke, small square windows with wooden shutters; 6 cells tall, 8 cells wide. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x576 px, PNG with fully transparent background.
```
OUTPUT: map_da_nha_vua_x4.png · 768×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**21. `map_da_nha_lon_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a tall carved rock tower residence with several stepped levels, wooden ladders and balconies between them, warm light in the window slits; 8 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x768 px, PNG with fully transparent background.
```
OUTPUT: map_da_nha_lon_x4.png · 768×768 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**22. `map_da_nha_dacsan_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a stone forge and quarry workshop: an open front with an anvil, a glowing furnace mouth, ore carts and hammers hanging on the wall; 6 cells tall, 8 cells wide. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x576 px, PNG with fully transparent background.
```
OUTPUT: map_da_nha_dacsan_x4.png · 768×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**23. `map_da_cong_lang_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: the village gate of the Hidden Stone: two huge rough rock monoliths carved with a large earth symbol, a massive stone lintel laid across, rubble at the feet; 5 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 480x480 px, PNG with fully transparent background.
```
OUTPUT: map_da_cong_lang_x4.png · 480×480 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**24. `map_da_cay_kho_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a dead gnarled tree with twisted bare branches and cracked grey bark; 8 cells tall, 4 cells wide. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 384x768 px, PNG with fully transparent background.
```
OUTPUT: map_da_cay_kho_x4.png · 384×768 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**25. `map_da_cay_nho_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a hardy small mountain tree with a thick short trunk and sparse small dusty-green leaves; 4 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 384x384 px, PNG with fully transparent background.
```
OUTPUT: map_da_cay_nho_x4.png · 384×384 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**26. `map_da_bui_ram_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a low dry scrub bush with grey-green leaves and thorny twigs; 2 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 288x192 px, PNG with fully transparent background.
```
OUTPUT: map_da_bui_ram_x4.png · 288×192 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**27. `map_da_duoc_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a torch post: a stone base, an iron bracket and a burning torch with a warm orange flame; 4 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x384 px, PNG with fully transparent background.
```
OUTPUT: map_da_duoc_x4.png · 192×384 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**28. `map_da_tang_da_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a large boulder showing rust-red and grey strata, cracked, a small stone at its foot; 3 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 288x288 px, PNG with fully transparent background.
```
OUTPUT: map_da_tang_da_x4.png · 288×288 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**29. `map_da_xe_day_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a wooden mine cart with iron bands full of rough ore rocks; 2 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 288x192 px, PNG with fully transparent background.
```
OUTPUT: map_da_xe_day_x4.png · 288×192 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**30. `map_da_hang_rao_da_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a low dry-stacked stone fence segment whose left and right edges join seamlessly with copies of itself; 2 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x192 px, PNG with fully transparent background.
```
OUTPUT: map_da_hang_rao_da_x4.png · 192×192 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**31. `map_da_bui_co_kho_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Stone (Earth) palette: earth brown, warm grey and rust-red rock strata, dusty beige, deep cave shadows, dusty pale grey-blue sky. SUBJECT: a tuft of dry grass growing between small rocks; 1 cell tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x96 px, PNG with fully transparent background.
```
OUTPUT: map_da_bui_co_kho_x4.png · 192×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

### 6.6 Làng Mây (map 5)

**1. `map_may_nen_troi_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a twilight violet-blue sky with towering storm clouds and a faint glow of distant lightning inside them. VERTICAL LAYOUT (top to bottom): 0-30% sky with clouds; 30-80% the sky pales toward the haze color; 80-100% flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(120, 110, 160) so that the very last pixel row is exactly RGB(120, 110, 160) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. No ground, no buildings. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x384 px, solid opaque PNG.
```
OUTPUT: map_may_nen_troi_x4.png · 2976×384 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**2. `map_may_nen_nui_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: sharp needle-like granite peaks piercing a sea of white-blue clouds, several overlapping peaks, snow and lightning-scarred rock. VERTICAL LAYOUT (top to bottom): 0-25% sky with clouds; 25-80% the solid mountain mass — it must fill at least 55% of the image height so it towers over every rooftop; 80-100% the foot dissolves into haze and the flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(120, 110, 160) so that the very last pixel row is exactly RGB(120, 110, 160) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x768 px, solid opaque PNG.
```
OUTPUT: map_may_nen_nui_x4.png · 2976×768 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**3. `map_may_nen_langxa_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a dense cloud-village seen from ground level: cylindrical stone towers with pointed roofs built on peaks, rope bridges between them, glowing lightning rods, clouds wrapping their bases. VERTICAL LAYOUT (top to bottom): 0-30% rooftops starting at the VERY FIRST pixel row — absolutely no sky at the top; 30-80% the dense village; 80-100% the lower village dissolves into haze and the flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(120, 110, 160) so that the very last pixel row is exactly RGB(120, 110, 160) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x576 px, solid opaque PNG.
```
OUTPUT: map_may_nen_langxa_x4.png · 2976×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**4. `map_may_nen_gan_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: the nearest background layer of a cloud village: a dense soft cloud bank with a few blurred tower tips poking out. VERTICAL LAYOUT (top to bottom): 0-30% blurred near shapes; 30-80% dense haze with faint details; 80-100% flat join color. The bottom 20% of the image must fade smoothly into the flat join color RGB(120, 110, 160) so that the very last pixel row is exactly RGB(120, 110, 160) with no texture left; never end with a hard horizontal cut. One single continuous painting across the full width (NOT mirrored, NOT tiled, NOT two halves pasted together), solid opaque background with no transparency, the left and right edges must join seamlessly when the image is repeated horizontally. Match the painterly style, brush quality and color depth of the approved reference concept-map1-langla.png. Soft and low contrast — it sits right behind the playable ground. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, borders, mirrored or duplicated halves, transparency, a visible seam. OUTPUT SIZE: exactly 2976x288 px, solid opaque PNG.
```
OUTPUT: map_may_nen_gan_x4.png · 2976×288 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**5. `map_may_tile_da_mat_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a ground-surface tile: the top ~35% is flat blue-grey granite with sparse alpine moss and a thin wisp of mist with an uneven natural edge, the lower 65% is blue-grey granite mixed with dark alpine earth (must continue into the underground tiles) Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the top area above the surface edge is open air and transparent.
```
OUTPUT: map_may_tile_da_mat_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**6. `map_may_tile_da_sang_lat_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a transition ground tile: left half flat blue-grey granite with sparse alpine moss and a thin wisp of mist, right half fitted blue-grey granite flagstones with a faint lavender sheen, the boundary irregular and natural with the surface material creeping over the paving; lower 65% blue-grey granite mixed with dark alpine earth Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the top area above the surface edge is open air and transparent.
```
OUTPUT: map_may_tile_da_sang_lat_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**7. `map_may_tile_lat_sang_da_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a transition ground tile: left half fitted blue-grey granite flagstones with a faint lavender sheen, right half flat blue-grey granite with sparse alpine moss and a thin wisp of mist, irregular natural boundary; lower 65% blue-grey granite mixed with dark alpine earth Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the top area above the surface edge is open air and transparent.
```
OUTPUT: map_may_tile_lat_sang_da_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**8. `map_may_tile_datlat_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a paved-path tile: the top 30% is fitted blue-grey granite flagstones with a faint lavender sheen with an irregular lower edge, the lower 70% is blue-grey granite mixed with dark alpine earth Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_may_tile_datlat_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**9. `map_may_tile_dat_1_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: an underground tile: blue-grey granite mixed with dark alpine earth with one horizontal vein of pale glowing crystals Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four underground tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_may_tile_dat_1_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**10. `map_may_tile_dat_2_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: an underground tile: blue-grey granite mixed with dark alpine earth with alpine moss roots and thin cracks Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four underground tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_may_tile_dat_2_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**11. `map_may_tile_dat_giua_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: an underground tile: blue-grey granite mixed with dark alpine earth with embedded large angular granite blocks, about 6% darker than tile_dat_1 Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four underground tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_may_tile_dat_giua_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**12. `map_may_tile_dat_day_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: an underground tile: blue-grey granite mixed with dark alpine earth with dense compressed granite bands, about 12% darker than tile_dat_1 Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). BRIGHTNESS: the four underground tiles (dat_1, dat_2, dat_giua, dat_day) must differ in TEXTURE and stay within the stated brightness steps, otherwise the ground turns into a checkerboard. Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_may_tile_dat_day_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**13. `map_may_tile_vach_trai_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: the vertical cut face of a granite cliff: the LEFT edge is the exposed face (sheer rock, thin cracks, a wisp of cloud at the top corner), the right side continues into solid blue-grey granite mixed with dark alpine earth Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_may_tile_vach_trai_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**14. `map_may_tile_vach_phai_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: the vertical cut face of a granite cliff: the RIGHT edge is the exposed face (sheer rock, thin cracks, a wisp of cloud at the top corner), the left side continues into solid blue-grey granite mixed with dark alpine earth Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_may_tile_vach_phai_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**15. `map_may_tile_chan_vach_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: the foot of a granite cliff where the vertical face meets flat ground: a small slope of loose crumbled material, a few small stones, blending softly into plain blue-grey granite mixed with dark alpine earth Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_may_tile_chan_vach_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**16. `map_may_tile_ke_da_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a retaining wall tile: fitted blue-grey granite blocks with a faint cool glow in the joints Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, fully opaque inside the square.
```
OUTPUT: map_may_tile_ke_da_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**17. `map_may_tile_bac_thang_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a single stair block seen from the side: a flat tread of blue-grey granite along the top, a vertical riser on the left, solid material filling the lower-right so copies stack diagonally into a staircase Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the area above the tread and left of the riser is open air and transparent.
```
OUTPUT: map_may_tile_bac_thang_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**18. `map_may_tile_bac_thang_dinh_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: the topmost stair block: a tread of blue-grey granite that merges on its right side into flat blue-grey granite with sparse alpine moss and a thin wisp of mist (matching the surface tile) Seamless tileable on all four edges (the texture must continue perfectly when the tile is repeated left, right, up and down), flat even lighting with no directional cast shadow, texture runs edge to edge, no border, no frame, no outline around the square, no vignette, readable texture at real size (not fine noise). Match the detail level and brush quality of the existing map-1 tiles (tile_co_mat_x4.png, tile_dat_lienmach_1_x4.png) and the approved concept-map1-langla.png. Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, objects casting shadows, a visible square frame, blurry noise instead of readable texture. OUTPUT SIZE: exactly 96x96 px, PNG with alpha channel, the area above the tread is open air and transparent.
```
OUTPUT: map_may_tile_bac_thang_dinh_x4.png · 96×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**19. `map_may_nha_nho_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a small cylindrical stone house with a pointed conical roof, a small lightning rod on top, a round door and windows glowing cool blue; 6 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 576x576 px, PNG with fully transparent background.
```
OUTPUT: map_may_nha_nho_x4.png · 576×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**20. `map_may_nha_vua_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a wider stone tower house with a wrap-around balcony, carved cloud motifs, a conical roof and glowing blue windows; 6 cells tall, 8 cells wide. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x576 px, PNG with fully transparent background.
```
OUTPUT: map_may_nha_vua_x4.png · 768×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**21. `map_may_nha_lon_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a tall multi-tier stone tower with a pointed cap, several lightning rods, balconies and small bridges stubs, lit windows; 8 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x768 px, PNG with fully transparent background.
```
OUTPUT: map_may_nha_lon_x4.png · 768×768 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**22. `map_may_nha_dacsan_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a sword forge and lightning shrine: an open front with hanging blades, a glowing forge, chains and a crackling blue electric lantern; 6 cells tall, 8 cells wide. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x576 px, PNG with fully transparent background.
```
OUTPUT: map_may_nha_dacsan_x4.png · 768×576 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**23. `map_may_cong_lang_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: the village gate of the Hidden Cloud: two tall granite pillars carved with lightning and cloud symbols, iron chains between them, faint glowing runes; 5 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 480x480 px, PNG with fully transparent background.
```
OUTPUT: map_may_cong_lang_x4.png · 480×480 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**24. `map_may_cay_thong_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a tall alpine pine with a straight trunk and layered dark-green needle branches, a little snow on the upper branches; 8 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x768 px, PNG with fully transparent background.
```
OUTPUT: map_may_cay_thong_x4.png · 768×768 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**25. `map_may_cay_thong_nho_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a small alpine pine with a compact conical crown; 4 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 384x384 px, PNG with fully transparent background.
```
OUTPUT: map_may_cay_thong_nho_x4.png · 384×384 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**26. `map_may_bui_ram_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a low alpine shrub with small blue-green leaves and a few white flowers; 2 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 288x192 px, PNG with fully transparent background.
```
OUTPUT: map_may_bui_ram_x4.png · 288×192 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**27. `map_may_cot_thu_loi_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a lightning rod post: a granite base, an iron pole and a glowing blue-white tip with a small electric crackle; 4 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x384 px, PNG with fully transparent background.
```
OUTPUT: map_may_cot_thu_loi_x4.png · 192×384 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**28. `map_may_den_loi_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: an electric lantern post with a glass lantern glowing lavender-blue and a coil of copper wire; 4 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x384 px, PNG with fully transparent background.
```
OUTPUT: map_may_den_loi_x4.png · 192×384 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**29. `map_may_cau_treo_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a rope suspension bridge segment: wooden planks on thick ropes between two stone anchor posts, a wisp of cloud passing under it; 3 cells tall, 8 cells wide. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 768x288 px, PNG with fully transparent background.
```
OUTPUT: map_may_cau_treo_x4.png · 768×288 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**30. `map_may_tang_da_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a large blue-grey granite boulder with a glowing crystal crack; 3 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 288x288 px, PNG with fully transparent background.
```
OUTPUT: map_may_tang_da_x4.png · 288×288 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**31. `map_may_hang_rao_xich_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a low fence segment of iron chains hung between two short granite posts; left and right edges join seamlessly with copies; 2 cells tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x192 px, PNG with fully transparent background.
```
OUTPUT: map_may_hang_rao_xich_x4.png · 192×192 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

**32. `map_may_bui_co_x4.png`** — P2 · can_ve

```text
detailed 2D game environment art for a Naruto ninja world, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), painterly anime background quality, NO flat vector shapes, NOT pixel art, NO photo textures. Hidden Cloud (Lightning) palette: violet-blue storm clouds, cool white-blue cloud banks, blue-grey granite, glowing lavender lightning accents, twilight lilac sky. SUBJECT: a tuft of short alpine grass with tiny blue flowers; 1 cell tall. SCALE: 1 grid cell = 96 px; the player character is about 246 px tall (2.5 cells), a small house is 6 cells (576 px), a big house 8 cells (768 px) — see THUOC-TI-LE.png. Straight side view at eye level (NO 3/4 perspective, NO top-down), the whole object in one image, its base touching the exact bottom edge of the canvas with zero empty rows below it, fully transparent background (PNG alpha), NO ground patch, NO grass or sand under it, NO drop shadow on the ground, NO background scenery, NO frame, no soft halo at the edges. Match the style of the approved reference concept-map1-langla.png and the existing map-1 objects (nha_nho_x4.png, tt_cay_x4.png, TOT-cong-torii-da-tach-nen.png). Do not draw: player characters, monsters, NPCs, readable text or letters, logos, UI, watermark, signature, ground, shadow, sky, scenery behind the object, a frame, empty transparent rows under the base. OUTPUT SIZE: exactly 192x96 px, PNG with fully transparent background.
```
OUTPUT: map_may_bui_co_x4.png · 192×96 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\

---

## 7. Tự kiểm và giao

### 7.1 Trước khi vẽ (phía dự án, 5 phút)

Nối **165 tên mới** vào cuối `D:\NRO_GoiBanGiao\docs\my-thuat\kich-thuoc-anh.json` kèm cỡ (không chèn giữa, không đổi tên cũ) — không nối thì `KIEM-TRA-ANH.py` và `TAO-4-MUC-ZOOM.py` **bỏ qua file** (`KIEM-TRA-ANH.py:97-100`, `TAO-4-MUC-ZOOM.py:126-128`). Danh sách máy đọc: khoá `ten_moi_can_them_kich_thuoc` trong `05-map.json`. Ví dụ 3 dòng đầu:

```json
 "nen_nuihokage_rong2_x4.png": [2976, 768],
 "nen_lang_xa_rong2_x4.png": [2976, 576],
 "tile_dat_van_1_x4.png": [96, 96],
```

Rồi chốt **màu nối** của map sắp vẽ (mục 3.3) và sửa mã màu trong prompt nếu concept cho màu khác.

### 7.2 Thứ tự làm — 3 ảnh mẫu trước, duyệt, rồi mới làm loạt

| Đợt | 3 ảnh mẫu (mỗi ảnh một kiểu khó nhất) | Vì sao chọn |
|---|---|---|
| P1 Làng Lá + Khu rừng | `nen_nuihokage_rong2_x4.png` · `tile_dat_van_1_x4.png` · `map_rung_cay_co_thu_1_x4.png` | nền rộng không ghép đôi + đáy nhoà đúng màu · tile vân đọc được mà vẫn lặp và cùng độ sáng · vật thể 8 ô chân chạm mép |
| P2 mỗi làng (sau concept) | `map_<làng>_nha_nho_x4.png` · tile mặt đi (`tile_co_mat` / `tile_cat_mat` / `tile_da_mat`) · `map_<làng>_cong_lang_x4.png` | y như cách map 1 đã làm (`PROMPT-CHO-AGENT.md` giai đoạn 1) |

Duyệt = cắm thử vào game hoặc ít nhất đặt lên concept, **nhìn thật**; ưng rồi mới gen phần còn lại theo thứ tự: nền 4 → tile 14 → công trình/cây → vật nhỏ → tiền cảnh.

### 7.3 Lệnh tự kiểm (người vẽ chạy, DÁN nguyên kết quả vào báo cáo)

```powershell
# 1. Tên · thiếu/thừa · cỡ (mức nguon = cỡ tự do) · alpha · ảnh phẳng (<40 màu) · trống · nhẹ (<3 KB)
python D:\NRO_GoiBanGiao\docs\my-thuat\KIEM-TRA-ANH.py D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\ nguon
#    -> 'THIEU KENH TRONG SUOT' cho các file nen_* / map_*_nen_* là BÌNH THƯỜNG (nền đặc). Mọi [X] khác phải = 0.
#    -> 'THIEU so voi manifest' sẽ liệt kê cả 691+ tên khác — bỏ qua dòng đó, chỉ nhìn SAI TEN / SAI KICH THUOC / ANH PHANG / ANH TRONG.

# 2. Kiểm ĐÚNG CỠ x4 sau khi sinh 4 mức zoom (phía dự án)
python D:\NRO_GoiBanGiao\docs\my-thuat\TAO-4-MUC-ZOOM.py D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\ D:\NRO_GoiBanGiao\docs\my-thuat\anh-game
python D:\NRO_GoiBanGiao\docs\my-thuat\KIEM-TRA-ANH.py D:\NRO_GoiBanGiao\docs\my-thuat\anh-game\x4 x4

# 3. Kiểm riêng map: viền bẩn, tile lặp, đáy nền nhoà đúng màu, ảnh có ghép đôi không, chân vật chạm mép
python D:\NRO_GoiBanGiao\tools\map-lang-la\kiem-anh-dot-2.py
#    ⚠️ ĐANG GHI CỨNG map 1: THU_MUC (dòng 11), CHO_DOI (14-29), MAU_TROI (36) = (10,86,154).
#    Cho lô này phía dự án phải sửa THU_MUC -> thư mục LƯU VÀO, thêm 24 tên map 1 + tên map mới vào CHO_DOI,
#    và MAU_TROI theo màu nối của map đang kiểm. CHƯA CÓ bản tham số hoá (kiem-anh-map.py <thư mục> <json cỡ> <màu nối>) — cần viết.

# 4. Đo 7 chỉ số: concept trước (đích), rồi ảnh game sau khi lắp
python D:\NRO_GoiBanGiao\tools\map-lang-la\do-map.py <concept.png>
python D:\NRO_GoiBanGiao\tools\map-lang-la\do-map.py <ảnh-chụp-game.png>   # -> tools\map-lang-la\chan-doan\*.png chỉ đúng chỗ hỏng
```

Ngưỡng đạt (`SO-TAY-LAM-MAP.md` §3): trời ≤ 30 % · dải nhà là trời ≤ 3 % · ô trơn dải giữa ≤ 10 % · **đất phẳng lì ≤ 40 %** · khối núi ≥ 30 % khung · đường kẻ ngang < 10 · nền lặp lại < 40. Đạt 6/7 mới gọi agent soi (gọi sớm là đốt 45–80 phút nghe lại điều máy nói trong 1 giây).

### 7.4 Lắp (phía dự án — công cụ có thật và chỗ còn thiếu)

- **Map 1:** `python D:\NRO_GoiBanGiao\tools\map-lang-la\lap-anh-dot-2.py` — thêm 4 tile vân vào bảng `TILE` (slot trống kế tiếp, hiện dùng tới 88) và đổi `NEN` sang `nen_nuihokage_rong2` → `b02`, `nen_lang_xa_rong2` → `b01` (script tự chép `b1x`); khai sàn 4 tile mới qua `tile_set_overlay`; xếp lại tầng đất trong Tiled (`tools\map-lang-la\tiled\`) cho 4 tile vân xen nhau, không lặp một ô liên tiếp; chạy `sinh-va-cham.py`; `do-map.py` đo lại — đích **đất phẳng lì ≤ 40 %**.
- **Map 6 và map 2–5:** **chưa có đường lắp** — cần (1) máy chủ khai `BanDo` map 6 + cổng + lưới va chạm (`banDo.go`), (2) chốt `TileID`/bộ tile riêng và `typeBg` cho từng map (mục 3.1), (3) một `lap-anh-map.py` tham số hoá thay bản ghi cứng. Ghi vào kế hoạch, **không tự chế**.
- Mọi ảnh nền phải bật *Read/Write Enabled*; đổi byte runtime thì tăng phiên bản nhóm tương ứng (`assetpack/deploy.py:197-239`).

### 7.5 Báo cáo nói thật (mẫu)

```text
Đã tạo: <n>/<tổng> ảnh. Cỡ thực từng ảnh: <tên: WxH> ... (đo bằng công cụ, không đoán)
Nền trong suốt: <danh sách có> / <danh sách không (nen_* là cố ý)>
Kết quả KIEM-TRA-ANH.py: <dán nguyên>
Ảnh KHÔNG CHẮC đạt và vì sao: <tên — lý do> (ví dụ: tile_dat_van_3 mép trên/dưới lệch 6/96; nen_troi hàng cuối chưa đúng mã màu)
Đã dừng sau 3 ảnh mẫu, chờ duyệt: có/không
```

Lô 343 từng tự nhận 'QA sạch' khi 62 ảnh sai cỡ (`_KHUON.md` §0.4) — chỉ ghi 'đã kiểm' khi thật sự chạy lệnh.

### 7.6 Khối copy-dán giao việc cho agent vẽ (giai đoạn 1 — P1)

```text
Bạn là hoạ sĩ tạo ảnh môi trường cho game 2D Naruto. Thư mục brief:
D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\

ĐỌC TRƯỚC KHI VẼ (theo thứ tự): _KHUON-DOT-4.md → ../dot-3-toan-bo/_KHUON.md → 05-map.md (mục 3, 4 rồi 5, 6)
MỞ XEM: ../map-1-lang-la/concept-map1-langla.png (phong cách) · ../map-1-lang-la/anh-tham-chieu/THUOC-TI-LE.png (tỉ lệ)
        · ../map-1-lang-la/anh-da-tao/tile_dat_lienmach_1_x4.png (màu đất gốc) · ../dot-3-toan-bo/anh-tham-chieu/XAU-*.png (TRÁNH)

VIỆC GIAI ĐOẠN 1: tạo ĐÚNG 3 ảnh mẫu, dùng nguyên prompt ở 05-map.md mục 6:
  1. nen_nuihokage_rong2_x4.png       2976x768  (mục 6.1 #1)
  2. tile_dat_van_1_x4.png            96x96     (mục 6.1 #3)
  3. map_rung_cay_co_thu_1_x4.png     768x768   (mục 6.2 #19)
MỖI LẦN GEN gắn kèm concept-map1-langla.png + ảnh tham chiếu ghi ở dòng dung_tham_chieu của ảnh đó trong 05-map.json.
Cỡ phải ĐÚNG từng pixel như ghi (công cụ ra cỡ khác thì cắt/resize về đúng cỡ trước khi lưu). Nền: nen_* ĐẶC, còn lại trong suốt thật.
Chân vật thể chạm đúng mép dưới. Nền không ghép đôi, hàng cuối đúng mã màu nối. Tile lặp 4 phía, không viền.

LƯU VÀO: D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\
Đặt đúng tên file như trên. Chạy: python D:\NRO_GoiBanGiao\docs\my-thuat\KIEM-TRA-ANH.py <thư mục trên> nguon và dán kết quả.
BÁO CÁO nói thật theo mẫu 05-map.md §7.5. DỪNG LẠI sau 3 ảnh, chờ duyệt. Không tự làm tiếp.
```

Giai đoạn 2 (sau duyệt): làm nốt P1 theo thứ tự nền → tile → cây → vật nhỏ, cùng thư mục, cùng cách báo cáo. P2 chỉ bắt đầu **sau khi concept của làng đó được duyệt và đo**.

---

## 8. Chưa cần vẽ đợt này

| Thứ | Vì sao chưa |
|---|---|
| **Bóng đổ** `bong_*` cho Khu rừng và 4 làng (mỗi công trình một vệt 96 cao) | map 1 đã có 8 ảnh nhưng đường lắp bg-item **chưa xác minh** (mục 2.1); vẽ khi đường đó chắc chắn — mẫu cỡ: `bong_nha_nho` 576×96, `bong_cay` 576×96, `bong_gieng` 288×96 |
| Hiệu ứng nền thời tiết theo làng (lá rơi Lá · cát bay Cát · sương Sương Mù · bụi Đá · sét Mây) | client đang dùng lại `/bg/mua*, lacay*, fog0, sao, tuyet` gốc (`BackgroudEffect.cs:143-404`) — chạy được; là P3 'đẹp thêm' (KIEM-KE-GAME-CAN P3-05) |
| Bản đồ thế giới / minimap `/img/map0.png` | thuộc brief UI (KIEM-KE-GAME-CAN P2-11), không phải ảnh map |
| Nội thất học viện (`map_hv_*`), sông–cầu–biển chung, `map_suongmu_*` bán trong suốt | không có map nào trong `catalog.json` dùng; vẽ khi máy chủ khai map |
| Tile nước cho Làng Lá, dốc chéo mới (`tile_doc_*` khớp bậc) | LOI-DA-GAP: đang né dốc chéo bằng bậc vuông; chỉ vẽ khi bố cục Tiled cần |
| Bố cục (.tmx), lưới va chạm, cổng, điểm sinh của map 6 và 2–5 | việc máy chủ/Tiled, không phải ảnh — `SO-TAY-LAM-MAP.md` bước 5 |
| Người chơi, NPC, quái, boss đứng trong map | brief 01–04 của đợt 4 |

*Sinh bởi script từ một nguồn dữ liệu duy nhất (MD và JSON đồng bộ). Con số kích thước đã kiểm bội số 96 và nền 2976 bằng máy trước khi ghi.*
