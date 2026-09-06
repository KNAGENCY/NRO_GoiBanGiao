# _KHUÔN ĐỢT 4 — spec kỹ thuật từng loại ảnh (đọc trước mọi brief `dot-4`)

> **File một cửa của đợt 4.** Bổ sung cho `../dot-3-toan-bo/_KHUON.md` (luật vàng, bảng cỡ tĩnh, màu 5 làng, master prompt) — phần đó vẫn đúng, **không chép lại**. File này nói thứ `_KHUON.md` chưa nói: game *thật sự* cắt/ghép/vẽ từng loại ảnh thế nào, đọc thẳng từ code client Unity và đường ống ảnh của dự án, để người vẽ (hay agent vẽ) không phải đoán.
>
> Viết 05/09/2026, **chỉ đọc code, không suy đoán**. Mỗi con số có `file:dòng`. Chỗ nào chưa kiểm được ghi rõ *"chưa chắc"*; công cụ nào chưa có ghi *"chưa có, cần viết"*.
>
> Thước tỉ lệ nhân vật kèm theo: `anh-tham-chieu\THUOC-NHAN-VAT.png` (cùng thư mục).

---

## A. Hai kiểu file — và khi nào dùng kiểu nào

Game nhận ảnh qua **hai đường ống khác nhau**. Chọn sai kiểu là ảnh đẹp mấy cũng không lắp được.

| | **Kiểu 1 — TĨNH ALPHA** (`tinh_alpha`) | **Kiểu 2 — DẢI CHROMA** (`dai_chroma`) |
|---|---|---|
| Là gì | Một PNG **RGBA**, nền **trong suốt thật**, một hình duy nhất | Một PNG nền **tím đồng nhất RGB(246, 4, 242)**, chia thành **n ô bằng nhau xếp ngang**, mỗi ô một khung |
| Dùng cho | icon vật phẩm `vp_`, vũ khí `vk_`, icon chiêu `kn_` (bản lẻ), UI `ui_`, tile/vật thể/nền map `map_`/`tile_`/`nen_`… | người chơi, boss, NPC (master đứng), quái, hiệu ứng `hu_`, HUD master, chân dung, bộ 4 icon chiêu một làng |
| Tên file | đúng tên trong `../kich-thuoc-anh.json`, hậu tố `_x4` là **khoá tra**, cỡ giao tự do miễn đúng tỉ lệ (`_KHUON.md` §1) | `<đối tượng>_<hành động>_<n>f_chroma.png` · `master_<tên>_x4_chroma.png` · `<tên>_portrait_chroma.png` — **con số `n` trong tên chính là số ô** máy dùng để cắt (`tools/artqa/artqa.py:450-452`, regex `_(\d+)(?:f|_chroma)`; không có số → 1 ô) |
| Máy làm gì tiếp | `KIEM-TRA-ANH.py` kiểm 7 điều → `TAO-4-MUC-ZOOM.py` thu về x4/x3/x2/x1 có premultiply (`TAO-4-MUC-ZOOM.py:64-79`) | `tools/artqa/artqa.py review` tách nền bằng flood-fill từ mép (`artqa.py:112-149`), chia ô đều theo chiều rộng (`artqa.py:152-184`), ghi `content_boxes` từng ô → review JSON → `task6b_closure.py` cắt mảnh/ghép atlas theo đúng hộp đó |
| Nền | trong suốt (riêng `nen_*`/`map_nen_*` là **nền đặc RGB**) | tím **phẳng tuyệt đối**: viền ảnh ≥ 85 % là tím, độ lệch màu ≤ 12 (`artqa.py:23-24, 74-95`); có gradient là nguồn bị đánh FAIL (`hu_anthan` v1 từng dính lỗi này — `review/vertical-slice-v1.json` mục `source_issues`) |

**Vì sao nhân vật phải là dải chroma chứ không phải PNG alpha?** Vì đường ống nhân vật cần biết **bàn chân từng khung nằm ở đâu, tâm khung ở đâu** để bù toạ độ vào bộ xương client — nó đọc `content_boxes` do `artqa` đo trên ô sạch nền. Ảnh alpha do AI xuất thường có viền mờ, bóng đổ, pixel rác — hộp đo sai là nhân vật lệch chân. Nền tím phẳng thì máy đo được chính xác đến từng pixel.

**Ba điều cấm với dải chroma** (rút từ cách máy tách nền, `artqa.py:60-62, 98-109`):
1. **Không dùng màu tím hồng/magenta trong hình vẽ** (R−G ≥ 110, B−G ≥ 110 là bị coi là nền). Vùng tím lớn trong hình cũng bị xoá (`task6b_closure.py:499-543` xoá đảo tím ≥ ~0,04 % diện tích). Bông tai Iruka nhỏ thì sống, mảng áo tím thì mất.
2. **Không glow/khói mờ lan ra nền**: pixel nửa trong (soft alpha) > 15 % là FAIL (`artqa.py:25, 245-269`).
3. **Hình không được cắn sang ô kề**: máy chia ô **đều nhau theo chiều rộng** (`artqa.py:152-153`), phần nào lọt sang ô bên là thành khung bên. Chiều rộng cả dải nên **chia hết cho n** và mỗi ô **chia hết cho 4** (để x1 không lệch nửa pixel).

---

## B. Spec kỹ thuật từng loại

### B.1 Người chơi (nhân vật ghép 3 mảnh, 33 khung xương)

| Mục | Spec | Nguồn |
|---|---|---|
| Kiểu nộp | **5 dải chroma**, ô **724×724**: `player_idle_3f` · `player_run_4f` · `player_melee_4f` · `player_cast_3f` · `player_hit_death_getup_3f` (tên và số khung cố định) | `tools/artqa/task6b_closure.py:41-47` (`PLAYER_ACTIONS`), mẫu `anh-nguon/vertical-slice-v1/player_*.png` 2172×724 |
| Client vẽ thế nào | Ba mảnh rời **đầu / thân / chân**; mỗi khung `cf` (0..32) tra bảng xương `CharInfo[cf][k] = (slot, x, y)`, k = 0 đầu, 1 chân, 2 thân. Toạ độ: `x = cx + CharInfo[cf][k][1] + part.dx`, `y = cy − CharInfo[cf][k][2] + part.dy`, neo **góc trên-trái**, thứ tự vẽ đầu → chân → thân | `D:\NRO_Project\Assets\Scripts\Game1\Char.cs:483` (bảng), `:6011` (`paintCharBody`), `:6081-6084`; bản sao bảng: `tools/nhan-vat/charinfo.json` |
| Số mảnh | **3 đầu / 17 thân / 14 chân = 34 mảnh** (part type 0/1/2 có đúng 3/17/14 frame) | `assetpack/part.py:7`, `assetpack/animation.py:18` |
| Thứ tự khung → slot | Ghép 17 khung theo đúng thứ tự dải: idle 0-1-2 → run 0-1-2-3 → melee 0-1-2-3 → cast 0-1-2 → hit 0-1-2. **Thân** lấy cả 17, **chân** lấy 14 khung đầu (bỏ hit), **đầu** lấy 3 khung idle | `task6b_closure.py:189-195` |
| Slot → cf thật (client) | Bảng dưới (sinh từ `charinfo.json`, đúng cách `assetpack/animation.py:21-41` làm) | |
| Một tỉ lệ chung | Mọi khung nhân với **một** hệ số `s = 128 / (rộng khung idle 0)`; hiện `s = 0,3765` (idle 0 rộng 340 px) | `task6b_closure.py:69-75, 256-264` |
| Neo | gốc mỗi khung = **giữa-dưới hộp nội dung cả khung** (bàn chân); bàn chân đặt tại `cy − 11` (đơn vị x1, `PLAYER_FEET_OFFSET_X1 = -11`) — đổi hằng này là dời cả nhân vật | `task6b_closure.py:72-75, 267-291` |
| Cắt mảnh | theo % chiều cao khung: đầu 0–42 %, thân 28–76 %, chân 66–100 % (chồng 10–14 %). **Chỉ đúng khi mọi khung cùng tỉ lệ và đứng thẳng** — cắt theo mốc cằm (V1) đang chờ chủ dự án quyết | `task6b_closure.py:143-159`; `BAT-DAU-TU-DAY.md:453-454` |
| Cỡ trong game | idle 0 cao 653 px nguồn × 0,3765 = **246 px x4 = 61,5 px x1 ≈ 2,56 ô** (ô = 24 px x1 = 96 px x4). Đầu 128×103 x4. Client vẽ ảnh mọi cỡ, không nhồi hộp | `review/vertical-slice-v1.json` (`content_boxes` idle), `animation.py:51-59` |
| dx/dy | mỗi mảnh MỘT `dx/dy` kiểu **signed byte −128..127** (đơn vị x1) dù phục vụ nhiều cf → sai số dư (`charinfo_residual`) đưa vào tối ưu MILP | `assetpack/part.py:48-51`, `task6b_closure.py:446-454`, `tools/artqa/nhan_vat_toi_uu.py:1-34` |
| Đích client | `res/x{1..4}/smallimage/Small<id>.png`, id **27041–27074** (đầu 27041-43, thân 27044-60, chân 27061-74); `part.bin` id 30/57/58 | `assetpack/build.py:1051`, `task6b_closure.py:1356-1358`, `tools/nhan-vat/chep-sprite-client.py:26-28` |
| Đổi part.bin thì | **phải tăng `vsData`** (hiện 27) — không tăng client dùng cache cũ | `tools/map-lang-la/README.md:59`, `tools/nhan-vat/sinh-lai-sprite.py:164, 226` |

**Bảng slot → cf thật** (client tra `CharInfo`; "khung" = khung trong dải nộp; đây là quyết định của đường ống, hoạ sĩ cần biết để hiểu khung nào lên hình lúc nào):

| Mảnh | slot | lấy từ khung | client vẽ ở cf | cf đó là gì |
|---|---|---|---|---|
| đầu | 0 | idle 0 | 0,1,7,8,11,12,17–24,30–32 | đứng, đánh, nạp, chết, chiêu — **đầu chính** |
| đầu | 1 | idle 1 | 2–6, 9,10, 13–16, 26–29 | chạy, nhảy — **đầu nghiêng** |
| đầu | 2 | idle 2 | 25 | sau dịch chuyển |
| thân | 0 | idle 0 | 23, 24 | **chết đứng / trúng đòn** (client không có khung nằm) |
| thân | 1 | idle 1 | 0, 1 | đứng |
| thân | 2 | idle 2 | 2, 15, 28 | bước chạy đầu |
| thân | 3–6 | run 0–3 | 3, 4, 5, 6 | chạy |
| thân | 7 | melee 0 | 7, 8, 17 | đánh / nhảy đà / nạp |
| thân | 8 | melee 1 | 12, 18, 22 | nạp trên đất |
| thân | 9 | melee 2 | 11, 21, 32 | đỉnh nhảy |
| thân | 10 | melee 3 | 13, 26 | |
| thân | 11 | cast 0 | 14, 27 | |
| thân | 12 | cast 1 | 9, 10 | nhảy |
| thân | 13 | cast 2 | 16, 29 | |
| thân | 14 | hit 0 | 19, 30 | |
| thân | 15 | **hit 1** | **20, 31** | **thế bắt đầu chiêu** — v1 vẽ khung này NẰM → lỗi |
| thân | 16 | hit 2 | 25 | sau dịch chuyển |
| chân | 0 | idle 0 | 8 | |
| chân | 1 | idle 1 | 0, 1, 22 | đứng |
| chân | 2–6 | idle 2, run 0–3 | 2, 3, 4, 5, 6 | chạy |
| chân | 7 | melee 0 | 7, 23 | đánh, **chết đứng** |
| chân | 8 | melee 1 | 12, 26–32 | |
| chân | 9 | melee 2 | 13–21, 24 | |
| chân | 10–12 | melee 3, cast 0, cast 1 | 9, 10, 11 | nhảy |
| chân | 13 | cast 2 | 25 | |

Ý nghĩa cf đọc từ `Char.cs`: 0/1 đứng lắc (`:2706-2731`), 7 rơi/đà (`:3375, 3824, 4029`), 8 đang di chuyển nhanh (`:4047`), 9/10/11 các pha nhảy (`:3387-3395`), 12 nạp/hạ xuống đất (`:2029, 2542, 3908, 4003`), 17 nạp chiêu / chờ (`:2001, 2026, 2815`), 16/31/32 thế chiêu (`:2124, 2242-2246`), 23 trúng đòn / chết / đổi map (`:1750, 1967, 2770, 3249, 3276`), 25 sau dịch chuyển (`:6245`; bóng dịch chuyển vẽ cf 25/7 `:5301-5305`). Cf 2–6 là chạy (`BAT-DAU-TU-DAY.md:460`).

### B.2 Boss

| Mục | Spec | Nguồn |
|---|---|---|
| Client vẽ | **y như người chơi** — boss là nhân vật đủ 3 mảnh, server gửi bộ `outfit` = 3 part id (đầu/thân/chân) phải nằm trong `closure.part_ids` của `asset-lock.json`; hiện chỉ có bộ "học viên" 30/57/58 và bộ Iruka 2189/2190/2191 | `server/kho-noi-dung/boss.yaml:18-25`, `server/internal/game/assets/asset-lock.json` (`closure.part_ids`) |
| Kiểu nộp | Boss có đủ hành động → **5 dải chroma như người chơi** (34 mảnh, part mới 3 id). Boss chỉ đứng-đánh → tối thiểu vẫn cần idle 3 + melee 4 + cast 3 khung đúng thứ tự (thân slot 0–13); thiếu dải nào máy hiện **không** có đường "lặp khung thay thế" ngoài cách của Iruka (một ảnh cho mọi slot, `task6b_closure.py:1359-1361`) | `task6b_closure.py:162-224` |
| Thứ tự khung / neo / tỉ lệ | như B.1 — **cùng cỡ đầu trong mọi khung của con boss đó**, chân chạm mép dưới ô | |
| Cỡ mục tiêu | ghi rõ trong brief theo **bội của người chơi** (người chơi 61,5 px x1). Giới hạn cứng: mỗi mảnh `dx/dy` ∈ [−128, 127] x1 (`part.py:48-51`) và `CharInfo` offset cố định → boss cao quá ~2× người chơi có nguy cơ vượt dải bù; **chưa thử boss to, chưa chắc** | `assetpack/part.py:48-51` |
| Part id mới | nối sau legacy (2188), hiện `assetpack/part.py:9, 146-147` **khoá cứng** chỉ nhận thêm đúng 2189/2190/2191 → boss có bộ part riêng **cần nới `PHASE_ONE_ADDITIONS`** (chưa làm) | `assetpack/part.py:8-9, 144-147` |
| Boss Làng Lá hiện có | `nhan-gia-mat-na` (2 pha: pha 1 đội lốt học viên = outfit 30/57/58; pha 2 thân thật → **cần bộ part mới**), `hoc-vien-dao-tau` | `boss.yaml:94-250` |
| Ảnh đích | `res/x{n}/smallimage/Small<id>.png` (id mới trong dải dự trữ 27083–27999) + `part.bin` | `mapping-id.json` → `phase.reserved_sprite_id_range` |

### B.3 NPC

| Mục | Spec | Nguồn |
|---|---|---|
| Client vẽ | **cùng bảng xương `CharInfo` với người chơi** (`Npc.cs:305-313`); 3 part của NPC đều 3/17/14 frame **trỏ về cùng một sprite** → NPC chỉ có một tư thế đứng | `BAT-DAU-TU-DAY.md:458-459`, `task6b_closure.py:1359-1361` |
| Kiểu nộp | **1 master đứng** `master_<tên>_x4_chroma.png` (1 ô, ví dụ 1254×1254) + **1 chân dung** `<tên>_portrait_chroma.png` (1 ô) | `task6b_closure.py:720-818, 848-866` |
| Cắt | đầu 0–42 % / thân 28–76 % / chân 66–100 % từ hộp nội dung master | `task6b_closure.py:143-159, 758-769` |
| **Lỗi ĐÃ SỬA 06/09 — giữ để không lặp** | Iruka từng bị **nhồi vào hộp 128×160** (nội dung ≤ 120×152) → cao **38 px x1 = 62 % người chơi**; khe đầu +13, chân +29 vì `dx=dy=0` không qua `component_placement`. Bản 06/09 (`75a64cc`) dùng `iruka_component_images`: cùng cỡ đầu người chơi, cắt thân từ cằm, dx/dy theo cf 0 → cao 62–63 px. **Mọi NPC mới phải đi đường này**, đừng nhồi hộp | `task6b_closure.py:743-756, 1359-1361`; `BAT-DAU-TU-DAY.md:520-521` |
| Cỡ mục tiêu đợt 4 | **NPC người lớn cùng cỡ đầu người chơi** (thước `THUOC-NHAN-VAT.png`: cao mặt 96–99 px, tóc→cằm 266–271 px trong ô 724); cao tổng ghi theo ô: người lớn ≈ 2,6–2,8 ô, trẻ em ≈ 2,2 ô. Pipeline NPC **cần sửa** để dùng `player_uniform_scale` + `component_placement` như người chơi (chưa làm) | |
| Chân dung | cắt 1 ô → **200×200 x4** (`Small<id>`), server gửi id trong gói NPC, client vẽ ở khung hội thoại | `task6b_closure.py:922`, `server/internal/game/npc.go:53`, `docs/my-thuat/unity-runtime-key-map.json` (consumer `ChatPopup.cs:316`) |
| Đích | `smallimage/Small<id>.png` ×3 mảnh + 1 chân dung; part id mới ≥ 2192 (xem B.2 về khoá `PHASE_ONE_ADDITIONS`); `catalog.json` `npcs[]` head/body/leg/portrait phải có trong closure | `server/internal/game/catalog.go:119` |

### B.4 Quái (sprite sheet riêng, không dùng bộ xương)

| Mục | Spec | Nguồn |
|---|---|---|
| Client vẽ | `Mob.getData()` đọc `/Mob/<templateId>` + `/data` + `/img.png` (định dạng `EffectData`), không có thì xin máy chủ | `D:\NRO_Project\Assets\Scripts\Game1\Mob.cs:383-410` |
| Định dạng nhị phân | `EffectData`: danh sách ảnh con `(id, x, y, w, h)` mỗi trường **1 byte không dấu ≤ 255** (đơn vị x1); danh sách khung, mỗi khung = các component `(dx, dy, imageId)`; `arr_frames`; 2 byte đuôi | `assetpack/mob.py:41-91, 124-199` |
| Kiểu nộp | **4 dải chroma**: `<tên>_idle_2f` · `<tên>_run_4f` · `<tên>_pounce_2f` · `<tên>_hit_death_2f` → **đúng 10 khung**, thứ tự slot 0–9: idle 0,1 · run 2–5 · pounce 6,7 · hit 8 · death 9 | `task6b_closure.py:80-85, 952-989, 1491-1498`; `mob.py:10-11` (5 hành động, số khung khoá cứng) |
| Cỡ | mỗi khung nhồi **96×96 x4 = 24×24 x1 ≈ 1 ô** (sói); atlas 10 ô ngang **960×96 x4**, ảnh con `x = i×24, w = h = 24`, component `dx = −12, dy = −24` (neo giữa-dưới) | `task6b_closure.py:972-974, 1452-1472, 1487-1491`; `asset-lock.json` output `client/x4/Mob/0/img.png` 960×96 |
| **Bẫy tỉ lệ** | pipeline quái vẫn dùng `fit_to_canvas` **từng khung** (`task6b_closure.py:972-974`) — đúng lỗi "mỗi khung một tỉ lệ" của người chơi trước 05/09. Hoạ sĩ phải vẽ **10 khung cùng tỉ lệ, chân chạm mép dưới ô**; và pipeline **cần sửa** sang một tỉ lệ chung (chưa làm) | `BAT-DAU-TU-DAY.md:534-537` |
| Quái to hơn 1 ô | ảnh con ≤ 255 px x1 mỗi chiều → atlas 10 khung × 48 px (quái 2 ô) = 480 px **vượt byte** → phải xếp **2 hàng** (`x, y` khác 0). Encoder hiện xếp 1 hàng (`task6b_closure.py:1454-1466`) → cần sửa | `mob.py:60, 139-146` |
| Nhiều quái | `assetpack build` hiện chỉ nhận **một** `build_inputs.mob` (`build.py:1064-1091`); `quai.yaml.anh_id` phải nằm trong `closure.mob_ids` (hiện chỉ `[0]`) → thêm quái thứ 2 **cần mở rộng build** (chưa có) | `server/kho-noi-dung/quai.yaml:20-31`, `build.py:1064-1091` |
| Chọn khung khi chạy | client chọn frame theo bảng riêng (`Mob.cs:606-640, 1182`) — **chưa đối chiếu** 10 slot với bảng đó; *chưa chắc* khung nào hiện lúc nào ngoài idle/run |
| Đích | `client/Mob/<id>.bytes`, `client/Mob/<id>/data.bytes`, `client/x{n}/Mob/<id>/img.png` → deploy chép vào `Assets/Resources/res/...` | `build.py:1071-1091`, `assetpack/deploy.py:241-248` |

### B.5 Map — nền / tile / vật thể

| Mục | Spec | Nguồn |
|---|---|---|
| Tile | **96×96 x4 = 24 x1**, PNG alpha; 100 tile một bộ; đích `res/x{n}/t/1$<k>.png` (k = 1..100; riêng x1 là `t/1/t_<kk>.png`) | `assetpack/build.py:1052-1056`, `tools/artqa/generate_tile_sheet.py:19-22` |
| Tile lắp lẻ (map 1) | ảnh 96 chép thẳng vào slot 60–68, 83–88; tile mới phải khai là **sàn** trong `tile_set_info.bin` (bit 0x0002) qua `build_inputs.tile_set_overlay`, không khai = **rơi xuyên đất** | `tools/map-lang-la/lap-anh-dot-2.py:18-32`, `assetpack/tileset.py:1-33, 45-47` |
| Tile phải | lặp được 4 phía; không viền; 4 tile đất **cùng độ sáng (chênh < 8 %) khác vân** | `../KHUON-DE-BAI-ANH-MAP.md` §"Sáu lỗi", §"Riêng bốn tile đất" |
| Nền 4 lớp | `res/x{n}/bg/b{typeBg}{k}.png`, k = 0..3; **client thật nạp `typeBg = 1` → `b10..b13`**, phải chép cả `b0x` và `b1x`; **nền đặc RGB**; Unity phải bật *Read/Write Enabled* | `GameCanvas.cs:1694`, `lap-anh-dot-2.py:77-82`, `../dot-3-toan-bo/_KHUON.md` §7 |
| Cỡ nền | `nen_troi` **2976×384** · `nen_nui` **2976×768** · `nen_langxa` **2976×576** · `nen_gan` **2976×288** (x4); rộng ≥ 2908 vì game hiển thị x2 | `../KHUON-DE-BAI-ANH-MAP.md` bảng cỡ; `../SO-TAY-LAM-MAP.md` §1 |
| Bố cục dọc nền | đáy mỗi lớp nhoà về **đúng màu nối**; `nen_langxa` không có trời ở đỉnh | `../KHUON-DE-BAI-ANH-MAP.md` §"Bố cục dọc" |
| Vật thể rời | cỡ **bội số 96**, alpha, **chân chạm mép dưới**; hai đường lắp: (a) chiếm một slot tile với ảnh to hơn ô (`tt_cay_nho` 384 → tile 87); (b) **bg-item**: `res/x{n}/mapBackGround/<AnhID>.png`, AnhID **≥ 500**, thứ tự bảng −31 là một phần giao thức, thêm mẫu **chỉ nối cuối** | `lap-anh-dot-2.py:34-38`, `server/internal/game/bgitem.go:50-95` |
| Bảng cỡ nhanh | 1 ô 96 · nhân vật ≈ 246 (2,56 ô) · nhà nhỏ 576 · nhà lớn 768 · cổng 480 · cây 768 | `../SO-TAY-LAM-MAP.md` §1 (bảng ghi nhân vật 232 = số cũ 2,4 ô; số đo mới 246) |
| asset-lock hiện tại | `bg/b00..b02` trong lock vẫn là **96×96 giữ chỗ**, `b03` 384×192 — khác ảnh đang chạy trong client (lắp tay bằng `lap-anh-dot-2.py`). Đường `assetpack build` cho nền map thật **chưa nối** | `asset-lock.json` outputs `client/x4/bg/*` |

### B.6 Icon vật phẩm / vũ khí

| Mục | Spec | Nguồn |
|---|---|---|
| Kiểu | tĩnh alpha; `vp_*` **64×64** (260 tên), `vk_*` 96×48 (13) · 64×64 (5) · 96×96 (2) | `../kich-thuoc-anh.json` |
| Client | vẽ qua `SmallImage` theo **số hiệu ảnh** (`Small<id>.png`); client **cache theo idImage trên máy người chơi** → id mới phải nằm ngoài vùng bản gốc (~0..350): dùng dải dự trữ **27083–27999** | `BAT-DAU-TU-DAY.md:441-444`, `mapping-id.json` `phase.reserved_sprite_id_range` |
| Trạng thái | 5 vật phẩm mẫu trong `server/kho-noi-dung/vat-pham.yaml`; KH7 sẽ mở rộng — icon vẽ theo `kich-thuoc-anh.json`, gắn id khi lắp |

### B.7 Icon chiêu

| Mục | Spec | Nguồn |
|---|---|---|
| Kiểu đã dùng | **dải chroma 4 ô** `skill_icons_<lang>_4_chroma.png` (1983×793, 4 ô ≈ 496) → mỗi ô nhồi **64×64 x4** → `Small27079..27082` | `task6b_closure.py:926-946`, `asset-lock.json` |
| Kiểu lẻ | tĩnh alpha `kn_<lang>_<tên>_x4.png` 64×64 (20 tên, đủ 5 làng × 4) | `../kich-thuoc-anh.json` |
| Server | `chieu.yaml` `icon_id = 1 + lang×4 + slot` (1..20) là **số vị trí, chưa phải id ảnh thật**; `skill.bin` hiện chỉ 4 pose 27079–27082 (Làng Lá), `effect_id = −1` | `server/kho-noi-dung/chieu.yaml:22, 52`, `task6b_closure.py:1365-1375` |
| Đã có ảnh | đủ 5 dải (leaf/sand/mist/stone/cloud) trong `vertical-slice-v1`, **chỉ leaf đã bind** | `anh-nguon/vertical-slice-v1/skill_icons_*` |

### B.8 Hiệu ứng (`hu_`)

| Mục | Spec | Nguồn |
|---|---|---|
| Client | `Effect` đọc `/x{zoom}/effectdata/<id>/data` + `/effectdata/<id>/img.png`, **cùng định dạng `EffectData`** như quái (id 101–199 dùng `readData2`); client hiện có 21 thư mục effectdata (100…) | `D:\NRO_Project\Assets\Scripts\Game1\Effect.cs:156-167`, `Assets/Resources/res/x4/effectdata/` |
| Kiểu nộp | dải chroma `hu_<tên>_<n>f_chroma.png`, **n ≥ 5**, mỗi ô rộng **chia hết cho 4**, cả dải chia hết cho n; tâm hiệu ứng **cố định giữa ô** mọi khung (client đặt bằng dx/dy component, không tự căn) | `_KHUON.md` §2 ghi chú; `mob.py:41-91` |
| Cỡ tĩnh cũ | `kich-thuoc-anh.json`: `hu_*` **1200×480 = 5 khung 240×480** (10 tên) — chỉ là cỡ khai, chưa có đường lắp | `../kich-thuoc-anh.json` |
| Mẫu đã có | `hu_anthan_10f_chroma.png` 1774×887 → 10 ô **177,4 px** (không chia hết 4; alpha 1570×694); chính sách **one-shot** (không lặp) | `review/vertical-slice-v1.json`, `artqa.py:237-242` |
| Giới hạn byte | mỗi ảnh con ≤ **255 px x1** (≈ 1020 px x4) mỗi chiều; hiệu ứng to phải xếp nhiều hàng trong atlas | `mob.py:60, 139-146` |
| Lặp / một lần | one-shot: khung cuối được tan biến; **loop** (aura, hồi máu): khung cuối phải nối về khung đầu — máy kiểm `loop_boundary_motion_jump` và **cấm 2 khung liền giống nhau** | `artqa.py:203-234` |
| **Trạng thái lắp** | **CHƯA LẮP**: máy chủ gửi blob effect rỗng (`dongbo.go:82-85`), chiêu `effect_id −1`, `assetpack build` không sinh `effectdata`. Bộ đóng gói `EffectData` cho `hu_` **chưa có, cần viết** (`encode_mob` cùng định dạng nhưng khoá 5 hành động quái, `mob.py:98-121`) | `server/internal/game/dongbo.go:72-91`, `assetpack/build.py` |
| Tham chiếu xấu | `../dot-3-toan-bo/anh-tham-chieu/XAU-hu-anthan-qua-mo.png` (quá mờ, halo) |

### B.9 UI / HUD

| Mục | Spec | Nguồn |
|---|---|---|
| Kiểu | tĩnh alpha theo `kich-thuoc-anh.json` (62 tên `ui_`); riêng HUD thanh máu/chakra đã làm bằng **master chroma 2×2** (`hud_hp_chakra_4_chroma.png` 1774×887: trên-trái nền HP · trên-phải fill HP · dưới-trái nền chakra · dưới-phải fill chakra). Lưu ý: `artqa review` cắt theo tên `_4_` thành **4 cột dọc** (không khớp 2×2, vẫn PASS), còn closure cắt 2×2 theo tâm ảnh — đợt 4 nên tách 4 file tĩnh `ui_thanh_*` 240×32 cho khỏi lệch | `task6b_closure.py:1127-1167` (`:1136-1142` cắt 2×2), `review/vertical-slice-v1.json` |
| HUD | 4 ảnh **240×32 x4** → `res/x{n}/mainImage/naruto-hp-bg|naruto-hp-fill|naruto-chakra-bg|naruto-chakra-fill.png`; fill bị **cắt theo % từ mép trái** (25/50/100 %), không co → vẽ fill **đầy thanh, đẹp khi cắt ở bất kỳ x** | `assetpack/visual.py:41-66`, `build.py:1060-1062` |
| 8 UI đã khoá (x4) | `ui_logo_game` 512×256 → `res/x4/Logo/1.png` (mọi zoom dùng bản x4, `GameCanvas.cs:2517-2519`) · `ui_khung_avatar` 96×96 → `naruto/ui/avatar-frame.png` · `ui_khung_hoithoai` 480×160 → `naruto/ui/dialog-frame.png` · `ui_nut_khoa` 96×48 → `naruto/ui/button-disabled.png` · `ui_o_vatpham` 64×64 → `naruto/ui/item-slot.png` · `ui_o_dangchon` 64×64 → `naruto/ui/item-slot-selected.png` · `ui_nut_thuong` / `ui_nut_nhan` 96×48 → `mainImage/btn0left|mid|right`, `btn1*` | `task6b_closure.py:86-95`, `docs/my-thuat/unity-runtime-key-map.json`, `asset-lock.json` outputs |
| **Bẫy** | `naruto/ui/*` là `overlay_required`: client gốc **không có dòng code nào nạp** khoá đó (call_sites rỗng) → vẽ xong vẫn chưa hiện cho tới khi có overlay code. Nút `slice_group`: 3 file left/mid/right hiện là **cùng một ảnh chép 3 lần** (SHA bằng nhau trong lock) — chưa cắt 3 mảnh thật | `unity-runtime-key-map.json`, `asset-lock.json` (`btn0*` 3841 byte ×3) |
| Cỡ UI theo bảng | 64×64 ×35 (nút, tiền, làng, hệ, bậc) · 96×48 ×3 · 240×32 ×4 · 48×48 ×4 (mũi tên) · 96×96 ×2 · 128×128 ×1 · 192×192 ×1 · 240×24 · 400×48 · 400×120 ×2 · 400×480 · 480×400 · 480×480 ×2 · 480×160 · 512×256 · 1024×576 ×2 | `../kich-thuoc-anh.json` |

### B.10 Chân dung (`cd_`)

| Mục | Spec | Nguồn |
|---|---|---|
| Kiểu | tĩnh alpha `cd_<tên>_x4.png` **200×200** (27 tên) **hoặc** `<tên>_portrait_chroma.png` 1 ô → máy nhồi 200×200 | `../kich-thuoc-anh.json`, `task6b_closure.py:907-924` |
| Dùng ở | khung hội thoại NPC (`Small<portrait>` gửi trong gói NPC); avatar người chơi trên HUD cắt từ ảnh đầu (`tao-chan-dung-nguoi-choi.py`) | `npc.go:53`, `tools/map-lang-la/README.md:55` |
| Nội dung | đầu + vai, nhìn 3/4 quay phải, cùng cỡ đầu/màu với master của NPC đó |

---

## C. Bài học 05/09 → 7 luật cho người vẽ nhân vật / NPC / quái / boss

**Chuyện đã xảy ra:** 5 dải người chơi v1 vẽ nhân vật **khác cỡ nhau** (dải chạy nhỏ hơn idle ~12 %, dải đánh ~17 %) và **khom lưng**. Đường ống cắt mảnh thân ở 28 % chiều cao khung — trên khung idle đó là dưới cằm, nhưng trên khung chạy/đánh thấp hơn thì lát cắt rơi lên **trán**: mảnh thân mang theo cả mắt–mũi–miệng, client vẽ nó **đè lên đầu idle** → trong game thấy "đầu tách thân", thực ra là **hai khuôn mặt**. Phải bù bằng toán (`part.bin` dx/dy, MILP) và nhân vật vẫn thấp đi khi chạy. Nguồn: `tools/artqa/nhan_vat_moc.py:4-15`, `BAT-DAU-TU-DAY.md:448-462`, hồ sơ `docs/my-thuat/review/khe-dau-than-2026-09-05/`.

Bảy luật dưới đây **bắt buộc chép vào mục 3 của mọi brief nhân vật**, kèm `THUOC-NHAN-VAT.png`:

| # | Luật | Vì sao (nói đời thường) |
|---|---|---|
| 1 | **Mọi khung của một nhân vật cùng CỠ ĐẦU** với khung `idle_0`: trong ô 724×724, cao mặt **96–99 px**, tóc→cằm **266–271 px**, rộng mặt ~168 px; cằm→bàn chân **375–395 px** khi đứng (thước: `THUOC-NHAN-VAT.png`) | Client chỉ có **3 cái đầu** (idle) dùng chung cho 33 khung. Đầu to mà thân nhỏ thì đầu lửng; thân to hơn thì mặt thứ hai lộ ra. Máy có thể co giãn cả khung nhưng không "vẽ lại" được cổ vai cho khớp |
| 2 | **Bàn chân chạm mép dưới ô ở MỌI khung** — không khung nào bay, không đệm trống dưới chân | Máy lấy **mép dưới vùng đục** làm bàn chân rồi đặt vào `cy − 11` (`task6b_closure.py:72-75`). Chân hở 20 px là nhân vật bay 20 px trong game; khung nhảy cũng vậy vì client tự dời `cy` khi nhảy |
| 3 | **Đứng thẳng như sprite Dragon Ball gốc**; chạy hơi nghiêng người là đủ, **không khom** | Cắt mảnh theo % chiều cao (`:147-159`): khom là cằm tụt xuống dưới lát cắt 28 % → mặt lọt vào mảnh thân. Đứng thẳng thì cằm của mọi khung nằm cùng độ cao |
| 4 | **Nền tím RGB(246, 4, 242) phẳng tuyệt đối** — không bóng đổ, không viền sáng/tối, không gradient, không "sàn" | Máy tách nền bằng cách "đổ nước" từ mép vào những pixel tím (`artqa.py:112-149`). Bóng đổ = pixel không tím dính vào chân → bàn chân đo sai; gradient = nguồn FAIL (`spread > 12`) |
| 5 | **Một ô một khung, các ô đều nhau, thứ tự khung cố định** như B.1 (idle 3 · run 4 · melee 4 · cast 3 · hit 3; quái idle 2 · run 4 · pounce 2 · hit_death 2) | Máy chia ô **đều theo chiều rộng** rồi gán slot **theo thứ tự** — đảo hai khung là trong game nhân vật đánh khi đang đứng. Con số trong tên file (`_4f`) là số ô máy dùng để chia |
| 6 | **Đầu nhìn 3/4, quay PHẢI ở mọi khung** | Client chỉ có hình quay phải; quay trái là **lật gương** (`Char.cs` neo `TOP|LEFT` / `TOP|RIGHT`). Vẽ một khung quay trái là trong game nó quay ngược đồng đội |
| 7 | **Dải hit KHÔNG vẽ nằm** → 3 khung đều **đứng**: trúng đòn ngả nhẹ · đỡ/loạng choạng · khuỵu gối | Client **không bao giờ vẽ nằm**: chết = cf 23 đứng + vòng sáng (thân lấy khung idle 0!). Ba khung "hit" thật ra lên hình ở cf 19/30, **20/31 (bắt đầu chiêu)** và **25 (sau dịch chuyển)** — xem bảng B.1. Khung nằm (hit 1, rộng 665 px) của v1 hiện ra đúng lúc người chơi tung chiêu; pipeline cũ còn co nó còn 1/5 khung đứng (`task6b_closure.py:57-60`) |

**Ghi chiều cao mục tiêu trong mọi brief** (đơn vị ô 96 px x4 / px x1): người chơi 246 px x4 = 61,5 px x1 ≈ **2,56 ô** · NPC người lớn **cùng cỡ đầu** người chơi · sói ≈ **1 ô** · boss: ghi rõ bội số, ≤ 2× người chơi cho tới khi thử được dải bù dx/dy (B.2).

---

## D. Quy trình nộp và nghiệm thu

### D.1 Ảnh tĩnh (`tinh_alpha`) — công cụ có sẵn, chạy được ngay

```powershell
# 1. Người vẽ tự kiểm (tên, cỡ, alpha, ảnh phẳng, trống, file nhẹ) — DÁN kết quả vào báo cáo
python D:\NRO_GoiBanGiao\docs\my-thuat\KIEM-TRA-ANH.py <thư-mục-ảnh> nguon
# 2. Phía dự án sinh 4 mức zoom (premultiply, không viền bẩn) rồi kiểm lại đúng cỡ x4
python D:\NRO_GoiBanGiao\docs\my-thuat\TAO-4-MUC-ZOOM.py <thư-mục-ảnh> <thư-mục-xuất>
python D:\NRO_GoiBanGiao\docs\my-thuat\KIEM-TRA-ANH.py <thư-mục-xuất>\x4 x4
```
Nguồn: `KIEM-TRA-ANH.py:30-35`, `TAO-4-MUC-ZOOM.py:30-35`. **Tên mới** phải được thêm vào `kich-thuoc-anh.json` trước, không thì cả hai lệnh bỏ qua file (`KIEM-TRA-ANH.py:97-100`, `TAO-4-MUC-ZOOM.py:126-128`). Đừng dùng `KIEM-TRA-ANH.ps1` (bản cũ không kiểm cỡ, `_KHUON.md` §6).

Riêng map: `python D:\NRO_GoiBanGiao\tools\map-lang-la\kiem-anh-dot-2.py` (cỡ từng pixel, viền bẩn, đáy nền, ghép đôi, chân chạm mép) — **đường dẫn và bảng cỡ ghi cứng trong file** (`kiem-anh-dot-2.py:11-29`), map mới phải sửa `THU_MUC`/`CHO_DOI`; `do-map.py <ảnh>` đo 7 chỉ số so concept (`tools/map-lang-la/README.md:11`).

### D.2 Dải chroma (`dai_chroma`) — có công cụ review, nhưng cổng nghiệm thu đang khoá theo lô v1

```powershell
# Sinh alpha / crop-report / contact-sheet / GIF xem thử cho MỌI *.png trong thư mục + file review JSON
python -m tools.artqa.artqa review --art-root <thư-mục-dải> --review <đường-dẫn-review.json> --reviewer "<tên>"
```
Nguồn: `tools/artqa/artqa.py:718-736` (chạy từ `D:\NRO_GoiBanGiao`, cần Pillow + numpy). Lệnh `review` **chạy được cho thư mục bất kỳ**: kết quả ở `<thư-mục-dải>\qa\{alpha,crop-reports,contact-sheets,previews}` — mở GIF trong `previews` để xem hoạt ảnh 100 ms/khung, mở `crop-reports\*_crop.json` xem `content_boxes` (ô nào rỗng → `empty_animation_cell`, ô nào dính khung kề → hộp tràn sát mép ô).

**Khoá đang có:** `artqa gate` / `bind-precondition`, `generate_vertical_slice.py`, `task6b_closure.py` đều **đòi đúng 20 tên file của lô v1** và `review_id = "vertical-slice-v1"` (`artqa.py:31-52, 582, 607`; `generate_vertical_slice.py:66-73`; `task6b_closure.py:34-35, 947-948`). Nghĩa là **chưa có cổng tự động nào nhận lô đợt 4** — *cần nới* (tham số hoá danh sách tên và dải ID) trước khi bind. Ghi vào kế hoạch, đừng tự "chế" review JSON tay.

Đọc kết quả review thế nào là ĐẠT (theo ngưỡng trong `artqa.py:23-27, 203-269`):
- `raw_source_status = PASS` (viền ≥ 85 % tím, lệch ≤ 12) — v1 **0/20** đạt điều này, chỉ đạt bản alpha dẫn xuất;
- `issues = []`: không ô rỗng, soft alpha ≤ 15 %, không 2 khung liền giống nhau, loop không nhảy;
- `content_boxes` mỗi ô: mép dưới hộp ≈ đáy ô (chân chạm), hộp không chạm mép trái/phải ô (không dính khung kề).

### D.3 Nhân vật ghép mảnh — đo trước khi tin mắt

```powershell
# In bảng 34 mảnh + dx/dy (không ghi gì) — HIỆN chỉ chạy cho 5 dải người chơi v1 (đường dẫn ghi cứng)
python D:\NRO_GoiBanGiao\tools\nhan-vat\sinh-lai-sprite.py --thu
# Ghép 33 cf đúng phép vẽ client, in khe đầu–thân / thân–chân / chân so cy / số mặt; thoát mã 1 nếu vi phạm
python D:\NRO_GoiBanGiao\tools\nhan-vat\ghep-khung.py --out kq.png --kiem "dau-than=-9,0;than-chan=-11,-2;so-mat=1;chan=-11,-10"
#   tuỳ chọn: --sprites-closure <thư mục sprites/<id>/x1/sprite.png> · --head/--body/--leg <part id> · --zoom 4
# Đo mốc đầu (cằm, cỡ đầu) từng khung — HIỆN khoá 5 tên dải người chơi v1
python -m tools.artqa.nhan_vat_moc --out <file.json>
```
Nguồn: `sinh-lai-sprite.py:28-30, 56-64`; `ghep-khung.py:15-36`; `nhan_vat_moc.py:26-27, 45-51`. Cửa `--kiem` phải **đủ 33/33 cf đạt và `so-mat = 1`** (khe âm vẫn có thể là hai mặt — `BAT-DAU-TU-DAY.md:450-452`); nhìn tấm ghép bằng mắt ở cả x1 và `--zoom 4`.

**Chưa có, cần viết:** (1) kiểm "cùng cỡ đầu mọi khung" cho **dải mới** (đo `nhan_vat_moc` đang ghi cứng 5 tên v1; `thu-nghiem-khe-dau-than/do-strip.py` là bản thử, chưa là cửa kiểm); (2) `sinh-lai-sprite.py` cho nhân vật/boss thứ hai (đang ghi cứng part 30/57/58, id 27041–27074); (3) pipeline NPC qua `component_placement`; (4) pipeline quái một tỉ lệ chung + atlas nhiều hàng + nhiều mob.

### D.4 Lắp và phát hành (phía dự án, sau khi mọi cửa xanh)

```powershell
python -m assetpack validate --manifest docs\my-thuat\mapping-id.json
python -m assetpack build --manifest docs\my-thuat\mapping-id.json --slice lang-la --out <thư mục> --review docs\my-thuat\review\<review>.json
python D:\NRO_GoiBanGiao\tools\nhan-vat\chep-sprite-client.py          # xem trước; thêm --that để ghi (Unity đã dừng Play)
```
Nguồn: `assetpack/__main__.py:14-36`, `chep-sprite-client.py:5-15`. `build` đòi `production_artifact_review` do reviewer độc lập ký (`sinh-lai-sprite.py:12-16`, mã `ASSET-ARTIFACT-REVIEW-016`) — đổi một byte sprite là chữ ký hết hiệu lực; **không tự sửa chữ ký**.

### D.5 Cách giao việc cho agent vẽ (giữ như `map-1-lang-la/PROMPT-CHO-AGENT.md`)

1. Khối copy-dán có dòng **`LƯU VÀO: <đường dẫn tuyệt đối>`** và danh sách ảnh **đúng tên**.
2. **3 ảnh mẫu trước** (một ảnh mỗi kiểu khó nhất) → duyệt → mới làm loạt.
3. Mỗi lần gen gắn kèm: `THUOC-NHAN-VAT.png` (nhân vật) hoặc `map-1-lang-la\anh-tham-chieu\THUOC-TI-LE.png` (map) + một ảnh "gốc" cùng bộ để đồng bộ (`player_idle_3f_chroma.png` cho người) + `TOT-*.png`; kèm `XAU-*.png` làm phản diện.
4. **Báo cáo nói thật**: cỡ thực từng ảnh, ảnh nào không chắc, dán nguyên kết quả `KIEM-TRA-ANH.py` / `artqa review`. Lô 343 từng tự nhận "QA sạch" khi 62 ảnh sai cỡ (`_KHUON.md` §0.4).

---

## E. Tên file và ID — append-only

| Không gian | Luật | Nguồn |
|---|---|---|
| Ảnh tĩnh | tên **đúng** trong `../kich-thuoc-anh.json` (691 tên; tiền tố `vp` 260 · `map` 143 · `ui` 62 · `tp` 60 · `quai` 44 · `npc` 27 · `cd` 27 · `vk` 20 · `kn` 20 · `boss` 18 · `hu` 10); mẫu `<tiền tố>_<làng>_<tên>_x4.png`, làng = `la` · `cat` · `smu` · `da` · `may`; tên mới **nối CUỐI** file kèm cỡ, không chèn giữa, không đổi tên cũ | `../dot-3-toan-bo/_KHUON.md` §5 |
| Dải chroma | `<đối tượng>_<hành động>_<n>f_chroma.png` · `master_<tên>_x4_chroma.png` · `<tên>_portrait_chroma.png` · `skill_icons_<lang>_4_chroma.png` · `hu_<tên>_<n>f_chroma.png`; **số trong tên = số ô** (`artqa.py:450-452`); **không đổi tên/số khung sau khi review** — review ghim `source_sha256` theo file (`artqa.py:629-630`) | mẫu `anh-nguon/vertical-slice-v1/` |
| Sprite id (`Small<id>`) | lô Làng Lá dùng **27041–27082**, dải dự trữ **27041–27999**, cấm 9000–9002 (POC); id mới **nối tiếp 27083**, không dùng lại id cũ (client cache theo idImage, `BAT-DAU-TU-DAY.md:441-444`) | `docs/my-thuat/mapping-id.json` → `phase` |
| Part id | legacy tới 2188; mới nối **2189+** theo thứ tự; `assetpack/part.py:9, 146-147` hiện khoá đúng {2189, 2190, 2191} → thêm part phải nới hằng | `assetpack/part.py:8-9, 144-147` |
| Mob (quái) id | `catalog.json.mobs[].id` ∈ `asset-lock.closure.mob_ids` (hiện `[0]`); `quai.yaml.anh_id` chỉ được trỏ id đã ký | `server/kho-noi-dung/quai.yaml:20-31`, `server/internal/game/catalog.go:119` |
| Bg-item (bóng đổ, vật thể nền) | `AnhID` **≥ 500**, đích `mapBackGround/<AnhID>.png`; **thứ tự bảng −31 là giao thức**, chỉ thêm cuối | `server/internal/game/bgitem.go:50-95` |
| Tile | slot 1..100 của bộ tile 1 (`t/1$<k>.png`); slot mới khai sàn qua `tile_set_overlay` | `build.py:1055`, `assetpack/tileset.py:145-194` |
| Effect id | client đã có 100–… trong `res/x4/effectdata/`; id 42–46 bị ánh xạ về 106, 101–199 dùng `readData2` (`Effect.cs:152-166`) → chọn id cho `hu_` **khi viết bộ đóng gói**, tránh trùng cache cũ | `Effect.cs:152-166` |
| Phiên bản | đổi `part.bin` → tăng `vsData`; đổi byte runtime → `runtime_slice` và nhóm tương ứng phải tăng, không được hạ | `assetpack/deploy.py:197-239`, `tools/map-lang-la/README.md:59` |

**Manifest đợt 4** (một file JSON cạnh mỗi brief, schema mở rộng từ `../manifest-tao-anh.json`): mỗi phần tử `anh[]` thêm `"kieu": "tinh_alpha" | "dai_chroma"`, `"so_khung"`, `"o_kich_thuoc": [w, h]` (dải), `"luu_vao"` (đường dẫn tuyệt đối), `"uu_tien": "P1"|"P2"|"P3"`, `"trang_thai": "can_ve"|"da_co"|"ve_lai"`. Ảnh tạo ra lưu tại `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\<hạng mục>\`.

---

> **Ngoại lệ tên (ghi 06/09):** ảnh map của lô `map-1-lang-la` đang chạy còn dùng thêm tiền tố `nen_` · `tile_` · `tt_` · `bong_` · `nha_` · `tc_` (ví dụ `nen_nuihokage_rong_x4.png`, `tile_dat_van_1_x4.png`). Đây là nếp cũ đã lắp vào client, **hợp lệ, giữ nguyên** — brief 05 dùng tiếp; 11 tiền tố ở bảng trên là của `kich-thuoc-anh.json`.

## Tóm tắt một hơi cho người vẽ

Ảnh **tĩnh** (icon, UI, tile, vật thể): vẽ to đúng tỉ lệ, nền trong, đúng tên `kich-thuoc-anh.json`, chạy `KIEM-TRA-ANH.py`. Ảnh **chuyển động** (người, boss, NPC, quái, hiệu ứng): **dải nền tím phẳng**, một ô một khung, số ô ghi trong tên, **cùng cỡ đầu mọi khung**, **chân chạm đáy ô**, đứng thẳng, quay phải, hit không nằm, không màu tím trong hình. Giao **3 ảnh mẫu trước**. Báo cáo nói thật.
