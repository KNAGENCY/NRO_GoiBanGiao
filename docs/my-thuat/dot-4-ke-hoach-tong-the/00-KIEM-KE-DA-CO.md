# 00 — KIỂM KÊ "ĐÃ CÓ GÌ" — Đợt 4 (kiểm kê xong 05/09/2026, chỉ đọc)

> Đây là bản đếm **bằng lệnh** (Python đọc header PNG + JSON), không đếm bằng mắt, không suy đoán.
> Phạm vi soi: `docs\my-thuat\anh-nguon\343-anh` (đối chiếu `manifest-tao-anh.json`), `anh-nguon\vertical-slice-v1`,
> `map-1-lang-la\anh-da-tao`, `lang-la-client-assets`, `anh-game`, thư mục gốc `D:\NRO_GoiBanGiao\my thuat`,
> `server\internal\game\assets\asset-lock.json` + `catalog.json`, `mapping-id.json`, kho nội dung `server\kho-noi-dung\*.yaml`,
> và (chỉ đọc) `D:\NRO_Project\Assets\Resources\res`. Không sửa, không commit gì. Chỗ nào chưa chắc ghi rõ "chưa xác minh".

> ⚠️ **Cây git đang bẩn bởi agent khác trong lúc kiểm kê** (`git status` 05/09 ~22h): 33 mảnh `lang-la-client-assets\res\x1..x4\smallimage\Small27041–27072, 27074` đang bị sửa,
> `mapping-id.json` viết lại (+2 900/−791 dòng), `asset-lock.json` data **27 → 28**, `part.bin` đổi nội dung, `dongbo.go`, `tools\nhan-vat\*`, `tools\artqa\nhan_vat_*` mới.
> Số đo mảnh người chơi ở §3.1 là **ảnh chụp tại thời điểm đếm** (ví dụ Small27044 x4 đo 117×118 rồi 15 phút sau đã là 107×92). Mọi con số khác (343, chroma, map-1, client thật) không bị việc này đụng tới.

---

## 0. Tóm tắt một phút

| Kho | Đếm được | Ghi chú |
|---|---:|---|
| Lô ảnh tĩnh **343** (`anh-nguon\343-anh`) | **343 file** = 343 mục manifest, không thiếu không thừa | 62 file **to hơn** cỡ manifest đúng 4 lần (vk 20, npc 15, cd 15, quai 12) — là "vẽ to" chứ không hỏng; bản đúng cỡ nằm ở `anh-game\x4` |
| Bản 4 mức zoom (`anh-game\x1..x4`) | **1 372 file** (343 × 4) | Sinh bằng `TAO-4-MUC-ZOOM.py`; x4 = đúng cỡ manifest |
| Dải chroma vertical slice (`anh-nguon\vertical-slice-v1`) | **20 file** nguồn (đã QA 20/20 PASS bản dẫn xuất) | + thư mục `generated` (closure v1→v11, 7 bản build) + `qa` |
| Map 1 Làng Lá (`map-1-lang-la\anh-da-tao`) | **51 file** + 8 bản gốc trước sửa | hand-painted; phần lớn đã cắm vào client tile 40–88 và nền b00–b03 |
| Sao lưu art đã cắm (`lang-la-client-assets`) | Logo 1 · nền 16 · mảnh nhân vật 136 (34×4) · bản lùi 136 · tile 196 · mymap 1 | Bản byte thật đang chạy trong client |
| UI đợt 1 (`D:\NRO_GoiBanGiao\my thuat`) | **377 png**: `dot1-nhom1-duyet` 10 UI · `_v2-bung` 12 UI · `_v3-bung` 355 (bản gốc lô 343 + 12 UI) | 8 UI trong đó là **input của asset-lock**; 12 tên `ui_*` này KHÔNG có trong `343-anh` |
| Client thật `D:\NRO_Project\Assets\Resources\res` | x1: bg 8 · mapBackGround 8 · t 48 · smallimage 47 · x4/Logo 1 | **Không có** thư mục `Mob/`, `naruto/ui/`, `mainimage/naruto-*` dù asset-lock có ghi output — xem §4 |
| Bảng cỡ `kich-thuoc-anh.json` | **691 tên**; đã có ảnh **355** (343 + 12 UI); **chưa có 336** | Chi tiết §5 |

**Ba câu kết luận:**
1. **Đang lắp thật trong client** chỉ có: 34 mảnh người chơi học viên (27041–27074, sinh lại 05/09), Iruka 3 mảnh + chân dung (27075–27078), 4 icon chiêu Làng Lá (27079–27082), tile + 4 lớp nền map 1, logo wordmark, HUD máu/chakra (qua lock). Sói: ảnh có trong lock nhưng **không có file trong Resources** — client lấy qua máy chủ/cache, chưa xác minh.
2. **Kho 343 ảnh tĩnh** phần lớn dùng lại được cho **icon** (vật phẩm 100, UI 50, vũ khí 20, chân dung 15) sau khi soi phong cách; còn **quái 12, boss 6, NPC 15, trang phục 36** là **format cũ đứng nguyên con** — client không ăn được, chỉ dùng làm **ảnh tham chiếu** khi vẽ dải chroma.
3. **Lệch lock ↔ client thật** ở 3 chỗ (nền map, logo, UI overlay) — phải chốt lại trước khi đợt 4 bind thêm ảnh (§4).

---

## 1. Kho nào ở đâu

| Thư mục (tuyệt đối) | Số file ảnh | Là gì | Kiểu ảnh |
|---|---:|---|---|
| `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\01-giao-dien` | 50 | Icon UI, khung, nền màn | tĩnh alpha |
| `…\343-anh\02-vat-pham` | 100 | Icon vật phẩm 64×64 | tĩnh alpha |
| `…\343-anh\03-trang-phuc` | 36 | 12 bộ × 3 mảnh (đầu 128×128, thân 128×160, chân 128×128) | tĩnh alpha, 1 tư thế |
| `…\343-anh\04-vu-khi` | 20 | Vũ khí rời | tĩnh alpha (vẽ to 4×) |
| `…\343-anh\05-quai` | 12 | Quái đứng nguyên con 384×384 (chuẩn 96×96) | **format cũ** |
| `…\343-anh\06-boss` | 6 | Boss 160×160 phẳng | **format cũ, lô xấu** |
| `…\343-anh\07-npc` | 30 | 15 `npc_` đứng nguyên con 384×512 (chuẩn 96×128) + 15 `cd_` chân dung 800×800 (chuẩn 200×200) | **npc_ format cũ**; cd_ dùng được |
| `…\343-anh\08-hieu-ung` | 10 | Dải hiệu ứng 1200×480 | dải, chưa cắt khung |
| `…\343-anh\09-ban-do` | 79 | Tile 96×96 (77) + cổng làng 192×192 + núi Hokage 384×192 | tĩnh, pixel-art |
| `D:\NRO_GoiBanGiao\docs\my-thuat\anh-game\x1..x4` | 1 372 | 343 ảnh × 4 mức zoom, x4 = đúng cỡ manifest | phái sinh |
| `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1` | 20 | Dải chroma tím (246,4,242): người chơi 5 dải, sói 4 dải, 3 master, 1 chân dung, HUD, hu_anthan, 5 sheet icon chiêu | **dải chroma** |
| `…\vertical-slice-v1\generated\production-closure-v11` | 642 | Closure đang ký trong asset-lock: 42 sprite (27041–27082) ×4 zoom, HUD 4, mob img ×4, tile-art 100 ×4, UI 8 ×4, evidence 2 | phái sinh đã ký SHA |
| `…\vertical-slice-v1\generated\{production-closure-v1..v10, canonical-sprites-v1, task6b-*}` | ~7 600 | Các vòng closure/build cũ | lịch sử, không dùng |
| `…\vertical-slice-v1\qa` | 20 alpha + 20 contact sheet + 20 preview + 2 tổng | QA Task 6A | bằng chứng |
| `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-da-tao` | 51 (+8 `goc-truoc-khi-sua`) | Nền 5, nhà 4, tile 22, vật trang trí 10, bóng 8, tiền cảnh 2 | hand-painted, bội số 96 |
| `D:\NRO_GoiBanGiao\docs\my-thuat\lang-la-client-assets` | Logo 1 · `res\x1..x4\bg` 4×4 · `res\x1..x4\smallimage` 34×4 · `smallimage-truoc-05-09` 34×4 · `res\x1..x4\t` 34/64/34/64 · `res\mymap\1.bytes` | Sao lưu byte thật đã cắm client (05/08 + 05/09) | bản chạy |
| `D:\NRO_GoiBanGiao\my thuat\dot1-nhom1-duyet\dot1-nhom1-duyet\01-giao-dien` | 10 | UI đợt 1 đã duyệt — **8 file là input asset-lock** | tĩnh alpha |
| `D:\NRO_GoiBanGiao\my thuat\_v2-bung\01-giao-dien` | 12 | Như trên + thanh máu/chakra tách `_nen`/`_fill` | tĩnh alpha |
| `D:\NRO_GoiBanGiao\my thuat\_v3-bung` | 355 | Bản gốc lô 343 đúng cỡ manifest + 12 UI (62 UI) | trùng nội dung 343 |
| `D:\NRO_GoiBanGiao\docs\my-thuat\{anh-tham-chieu, dot-2-can-ve, dot-3-toan-bo, map-1-lang-la}\anh-tham-chieu`, `dot-4-ke-hoach-tong-the\anh-tham-chieu` | 12 / 16 / 16 / 16 / 1 | TOT-*, XAU-*, HIENCO-*, THUOC-* | tham chiếu, không phải asset |

---

## 2. Cách đọc cột "trạng thái lắp"

| Nhãn | Nghĩa |
|---|---|
| **LẮP (client)** | Có **file thật** trong `D:\NRO_Project\Assets\Resources\res` (kiểm ngày 05/09) |
| **LẮP (lock)** | Có trong `inputs`/`outputs` của `server\internal\game\assets\asset-lock.json` (closure v11, versions data 27 / map 22 / resource 21) — nhưng có thể **chưa** có file trong client thật |
| **Kho** | Chỉ nằm trong `docs\my-thuat`, chưa ai bind |
| **Thay rồi** | Đã có bản mới hơn (dải chroma) thay thế; bản cũ chỉ còn giá trị tham chiếu |
| **Tham chiếu** | Không phải asset; dùng để so phong cách/thước |

Cột "dùng lại được không" xét theo **kiểu file client cần** (xem bối cảnh: người chơi/boss = 3 mảnh 33 khung; NPC = 3 mảnh đứng + chân dung; quái = sprite sheet 4 dải; map = 4 lớp nền + tile 96; icon = tĩnh alpha).

---

## 3. Kiểm kê theo hạng mục

### 3.1 Nhân vật người chơi

| Tên | Nơi | Cỡ | Trạng thái lắp | Dùng lại được không | Ghi chú |
|---|---|---|---|---|---|
| `player_idle_3f_chroma.png` | `anh-nguon\vertical-slice-v1` | 2172×724, ô 724, 3 khung | **LẮP (lock + client)** — nguồn mảnh 27041–27074 | Có — **là chuẩn tỉ lệ gốc** (THUOC-NHAN-VAT đo từ idle_0) | Đầu 3/4 quay phải, đứng thẳng |
| `player_run_4f_chroma.png` | 〃 | 2172×724, 4 khung | LẮP (lock + client) | **Vẽ lại (P1)** | Nhân vật vẽ **nhỏ hơn idle 12 %** và khom 13 % → phải bù toán, nhân vật thấp ~10 % khi chạy (bài học 05/09) |
| `player_melee_4f_chroma.png` | 〃 | 2172×724, 4 khung | LẮP (lock + client) | **Vẽ lại (P1)** | Nhỏ hơn idle **17 %**, khom 6 % |
| `player_cast_3f_chroma.png` | 〃 | 2172×724, 3 khung | LẮP (lock + client) | Giữ được (tỉ lệ 1,00–1,07) | Khung cast0/cast2 gán slot thân 15/14 |
| `player_hit_death_getup_3f_chroma.png` | 〃 | 2172×724, 3 khung | LẮP (lock + client) | **Vẽ lại (P1)** | Có khung **nằm** — client không bao giờ vẽ nằm (chết = đứng + vòng sáng) → 3 khung phải là đứng: ngả nhẹ / đỡ / khuỵu gối |
| `master_hoc_vien_x4_chroma.png` | 〃 | 1254×1254 | Lock: `reviewed_vertical_slice_source` (không sinh sprite) | Có — tham chiếu nhân vật | Bản đứng nguyên con để đồng bộ trang phục |
| 34 mảnh `Small27041–27074.png` | `lang-la-client-assets\res\x1..x4\smallimage` (nguồn) → `D:\NRO_Project\...\res\x1..x4\smallimage` | x4: đầu 128×103…, thân 108–186 × 59–113, chân 107–134 × 59–84 (canvas trim theo pose) | **LẮP (client, 05/09)** — **đang bị sửa tiếp** bởi agent khác (cây git bẩn, xem đầu file) | — (phái sinh) | Sinh bởi `tools\nhan-vat\sinh-lai-sprite.py` (một tỉ lệ 0,376; neo chân −11). `mapping-id.json sprites[].zooms` **vẫn trỏ closure v11 cũ** (128×128 / 128×160 canvas cố định) — sự thật mới nằm ở `build_inputs.player_component_geometry` |
| Bản lùi `smallimage-truoc-05-09\x1..x4` (136 file) | `lang-la-client-assets` | x4: 128×128 / 128×160 | Sao lưu | — | Bản closure v11 trước khi sinh lại |
| `tp_*` 36 mảnh trang phục (12 bộ) | `343-anh\03-trang-phuc` | 128×128 / 128×160 | Kho | **KHÔNG dùng thẳng** — 1 tư thế, không theo 33 khung | Dùng làm **tham chiếu thiết kế** khi vẽ dải cho bộ đồ mới (học viên, genin, chunin, jonin, ANBU, kage, lễ hội…) — bảng đủ ở §3.6 |
| `part.bin` | `server\internal\game\assets\part.bin` (87 846 byte) | 6 part: 30/57/58 (học viên) + 2189/2190/2191 (Iruka) | LẮP (server) | — | **V1 đã tích hợp 06/09** (`75a64cc`): dx/dy 34 mảnh người chơi là nghiệm MILP, 3 mảnh Iruka đặt theo `component_placement` cf 0; bản thử trước đó nằm ở `review\khe-dau-than-2026-09-05\part-V1-siet-ghim38.bin` |

Game cần thêm (chưa có ảnh): **không có bộ trang phục nào khác ngoài học viên** được vẽ thành dải; 4 làng còn lại chưa có nhân vật khởi đầu riêng (brief `dot-3-toan-bo\lang-*\nhan-vat.md` chỉ là đề xuất).

### 3.2 NPC

Game cần (catalog.json + cua-hang.yaml): **1 NPC** — Iruka (id 0, head 2189 / body 2190 / leg 2191, portrait 27078, cửa hàng "Tạp hoá Làng Lá").

| Tên | Nơi | Cỡ | Trạng thái lắp | Dùng lại được không | Ghi chú |
|---|---|---|---|---|---|
| `master_iruka_x4_chroma.png` | `anh-nguon\vertical-slice-v1` | 1254×1254 | **LẮP (lock + client)** → Small27075/27076/27077 (x4 128×160) | Có | **ĐÃ SỬA 06/09** (`75a64cc`): sinh lại 3 mảnh cùng cỡ đầu người chơi (r 1,664 · s 0,226) + dx/dy qua `component_placement` → cf 0/1 khe đầu–thân −3, thân–chân −6/−5, bàn chân −11, cao 62–63 px x1 (người chơi 61–62). Trước đó: dx=dy=0 → đầu lửng +13, thân chìm +29 |
| `iruka_portrait_chroma.png` | 〃 | 1254×1254 → sprite 27078 200×200 | LẮP (lock + client) | Có | Chân dung khung hội thoại |
| 15 `npc_*` (bà vũ khí, chú ramen, cô bé hoa, Hokage 3, Iruka, Kakashi, lính gác, người đưa thư, người giao NV, nhà chép cuộn, ông Tazuna, thầy học viện, thợ rèn, trọng tài, y tá) | `343-anh\07-npc` | 384×512 (chuẩn 96×128) | Kho (npc_iruka: thay rồi) | **KHÔNG dùng thẳng** — client cần 3 mảnh + chân dung | Dùng làm **ảnh tham chiếu nhân vật** khi vẽ `master_<npc>_x4_chroma`. Bảng đủ §3.9 kèm chân dung |
| `Small27083.png`, `Small27084.png` | `D:\NRO_Project\...\smallimage` (x4 128×128, 104×104; ngày 22/08) | | LẮP (client) | ? | **Không có trong lock (id ký chỉ tới 27082) và không khớp file nào trong repo** — chưa xác minh là gì |
| `Small9000/9001/9002.png` | `D:\NRO_Project\...\x1\smallimage` (28/07) | 32×32 / 32×40 | LẮP (client) | Không | Là `forbidden_poc_ids` trong mapping-id — **rác POC còn nằm trong client** |

NPC 4 làng còn lại (kich-thuoc-anh có 12 tên `npc_{cat,smu,da,may}_*` + 12 `cd_*`): **chưa có ảnh nào**.

### 3.3 Quái

Game cần (quai.yaml + catalog.json): **1 quái** — Sói Rừng (`anh_id 0`), sinh ở cả 6 map. Client vẽ quái bằng `Mob/<templateId>/img.png + data` (`Mob.cs:388`), thiếu thì xin máy chủ.

| Tên | Nơi | Cỡ | Trạng thái lắp | Dùng lại được không | Ghi chú |
|---|---|---|---|---|---|
| `wolf_idle_2f_chroma.png` | `anh-nguon\vertical-slice-v1` | 1774×887, 2 khung | **LẮP (lock)** → `client/x1..x4/Mob/0/img.png` (10 khung 24×24 ở x1) + `Mob/0.bytes` | Có | idle 2 / run 4 / pounce 2 / hit 1 / death 1 |
| `wolf_run_4f_chroma.png` | 〃 | 2172×724, 4 khung | LẮP (lock) | Có | |
| `wolf_pounce_2f_chroma.png` | 〃 | 1774×887, 2 khung | LẮP (lock) | Có | |
| `wolf_hit_death_2f_chroma.png` | 〃 | 1774×887, 2 khung | LẮP (lock) | Có | |
| `master_soi_rung_x4_chroma.png` | 〃 | 1254×1254 | Lock: nguồn đã review, không sinh sprite | Tham chiếu | |
| Ảnh sói trong client thật | `D:\NRO_Project\Assets\Resources\res\x*\Mob\` | — | **KHÔNG CÓ thư mục** | — | Sói vẫn hiện trong game (SOI-LAN-3, 05/09) → client nhận qua RMS cache hoặc máy chủ stream; **đường đi thực chưa xác minh** |
| 12 `quai_*` (bia gỗ, bọ cạp, cá sấu, chim ưng, khỉ rừng, lính đánh thuê, lợn rừng, người rơm, nhện độc, ninja giả, rắn lớn, sói rừng) | `343-anh\05-quai` | 384×384 (chuẩn 96×96) | Kho | **KHÔNG dùng thẳng** — 1 khung đứng | Tham chiếu dáng/màu để vẽ 4 dải mỗi con. `quai_soi_rung` = TOT-quai-soi-rung (mẫu tốt) |

Bảng đủ 12 quái:

| # | Tên file | Nơi lưu | Cỡ file (px) | Làng | Trạng thái lắp | Dùng lại được không | Nội dung (rút gọn) |
|---|---|---|---|---|---|---|---|
| 1 | `quai_bia_go_x4.png` | 343-anh/05-quai | 384×384 (chuẩn 96×96) | chung | Kho | KHÔNG dùng thẳng — ảnh đứng 1 khung; client cần sprite sheet idle 2/run 4/pounce 2/hit_death 2. Dùng làm tham chiếu dáng/màu khi vẽ dải | Bia gỗ tập luyện — cọc gỗ hình người dùng để luyện đòn ở sân Học viện N… |
| 2 | `quai_bo_cap_x4.png` | 〃 | 384×384 (chuẩn 96×96) | chung | Kho | 〃 | Bọ cạp — bọ cạp đen, đuôi giương cao với ngòi độc cong ở cuối. Thân giá… |
| 3 | `quai_ca_sau_x4.png` | 〃 | 384×384 (chuẩn 96×96) | chung | Kho | 〃 | Cá sấu sông — cá sấu bò thấp, thân dài sát đất. Da xanh lục/xám rêu, lư… |
| 4 | `quai_chim_ung_x4.png` | 〃 | 384×384 (chuẩn 96×96) | chung | Kho | 〃 | Chim ưng săn mồi — ưng nâu, cánh dang nửa chừng như chuẩn bị vồ. Lông n… |
| 5 | `quai_khi_rung_x4.png` | 〃 | 384×384 (chuẩn 96×96) | chung | Kho | 〃 | Khỉ rừng — khỉ mặt đỏ nhe răng hung dữ. Lông nâu/xám, mặt và mông đỏ đặ… |
| 6 | `quai_linh_danhthue_x4.png` | 〃 | 384×384 (chuẩn 96×96) | chung | Kho | 〃 | Lính đánh thuê — chiến binh mặc giáp chắp vá không đồng bộ (mảnh da, mả… |
| 7 | `quai_lon_rung_x4.png` | 〃 | 384×384 (chuẩn 96×96) | chung | Kho | 〃 | Lợn rừng — heo rừng lông cứng dựng, thân nâu sẫm/xám, cặp nanh cong sứt… |
| 8 | `quai_nguoi_rom_x4.png` | 〃 | 384×384 (chuẩn 96×96) | chung | Kho | 〃 | Người rơm — hình nộm rơm biết cử động, được luồng chakra xanh mờ điều k… |
| 9 | `quai_nhen_doc_x4.png` | 〃 | 384×384 (chuẩn 96×96) | chung | Kho | 〃 | Nhện độc — nhện thân tím, cụm mắt vàng (nhiều mắt xếp thành chùm) sáng,… |
| 10 | `quai_ninja_gia_x4.png` | 〃 | 384×384 (chuẩn 96×96) | chung | Kho | 〃 | Ninja giả (cướp bịt mặt) — tên cướp/ninja rởm đeo khăn bịt mặt tối màu,… |
| 11 | `quai_ran_lon_x4.png` | 〃 | 384×384 (chuẩn 96×96) | chung | Kho | 〃 | Rắn lớn — rắn xanh thân to cuộn tròn, ngóc đầu cao trong tư thế tấn côn… |
| 12 | `quai_soi_rung_x4.png` | 〃 | 384×384 (chuẩn 96×96) | chung | Thay rồi | Đã có dải wolf_*_chroma (đang lắp Mob 0). Ảnh này còn dùng làm TOT-quai-soi-rung (mẫu tốt) | Sói rừng — sói xám gầy, dáng thấp săn mồi, đang nhe nanh gầm gừ. Lông x… |

Quái 4 làng (kich-thuoc-anh: 32 tên `quai_{cat,smu,da,may}_*`, 8/làng): **chưa có ảnh nào**.

### 3.4 Boss

Game cần (boss.yaml): 2 boss Làng Lá — `nhan-gia-mat-na` (2 pha: "Kẻ Lạ Mặt", "Mặt Nạ Vỡ", hệ Hoả) và `hoc-vien-dao-tau` (1 pha, hệ Thuỷ). Boss vẽ bằng **part id trong `closure.part_ids`** → hiện **mượn bộ học viên 30/57/58** (ghi trong boss.yaml). **Không có ảnh boss riêng nào đang lắp.**

| Tên | Nơi | Cỡ | Trạng thái lắp | Dùng lại được không | Ghi chú |
|---|---|---|---|---|---|
| 6 `boss_*` (Haku, Mizuki, Orochimaru, rắn khổng lồ, trùm trộm bịt mặt, Zabuza) | `343-anh\06-boss` | 160×160 (không vẽ to) | Kho | **KHÔNG** — phẳng, chính là lô XAU; boss client = 3 mảnh 33 khung | Chỉ lấy **ý tưởng nhân vật**; không dùng làm mẫu nét |
| `TOT-boss-mau.png` | `dot-3-toan-bo\anh-tham-chieu` | 160×160 | Tham chiếu | — | Mẫu tốt về độ chi tiết boss |

| # | Tên file | Nơi lưu | Cỡ file (px) | Làng | Trạng thái lắp | Dùng lại được không | Nội dung (rút gọn) |
|---|---|---|---|---|---|---|---|
| 1 | `boss_haku_x4.png` | 343-anh/06-boss | 160×160 | chung | Kho | KHÔNG — 160×160 phẳng, chính là lô XAU-*; boss client = 3 mảnh 33 khung như người chơi. Tham chiếu yếu (chỉ lấy ý tưởng nhân vật) | Haku — Nhân vật mảnh khảnh, đeo mặt nạ thợ săn màu trắng có hoa văn/vạc… |
| 2 | `boss_mizuki_x4.png` | 〃 | 160×160 | chung | Kho | 〃 | Mizuki — Giáo viên học viện phản bội. Tóc bạc trắng dài ngang vai, khuô… |
| 3 | `boss_orochimaru_x4.png` | 〃 | 160×160 | chung | Kho | 〃 | Orochimaru — Nhân vật da tái nhợt (pale), tóc đen dài thẳng. Mắt vàng k… |
| 4 | `boss_ran_khonglo_x4.png` | 〃 | 160×160 | chung | Kho | 〃 | Rắn khổng lồ (Boss triệu hồi) — Một con mãng xà (rắn) khổng lồ màu tím … |
| 5 | `boss_trom_bitmat_x4.png` | 〃 | 160×160 | chung | Kho | 〃 | Trùm trộm bịt mặt — Một tên trùm cướp to con vạm vỡ (hulking) đeo mặt n… |
| 6 | `boss_zabuza_x4.png` | 〃 | 160×160 | Sương Mù | Kho | Chính là ảnh XAU-boss-zabuza.png (mẫu XẤU cần tránh) | Zabuza — Kiếm sĩ cao lớn, nửa dưới mặt quấn băng vải trắng (che miệng c… |

Boss 4 làng (kich-thuoc-anh: 12 tên `boss_{cat,smu,da,may}_*`): chưa có ảnh.

### 3.5 Map

Game cần: 6 map (1 Lá, 2 Cát, 3 Sương Mù, 4 Đá, 5 Mây, 6 Khu rừng). **Chỉ map 1 có mỹ thuật**; map 2–6 dùng bộ tile Ngọc Rồng cũ / biên dự phòng (BAT-DAU-TU-DAY §5 dòng 1a′).

**3.5.1 Map 1 Làng Lá — hand-painted (`map-1-lang-la\anh-da-tao`, 51 file)**

| Nhóm | Tên (cỡ x4) | Trạng thái lắp | Dùng lại | Ghi chú |
|---|---|---|---|---|
| Nền (5) | `nen_troi_x4` 768×768 · `nen_nuihokage_x4` 1536×768 · `nen_nuihokage_rong_x4` 1536×768 · `nen_lang_xa_x4` 1536×576 · `nen_lang_xa_rong_x4` 1536×576 | **LẮP (client)**: b03 = trời 768×768, b02 = núi 1536×768, b01 = làng xa 1536×576 (khớp cỡ x4 trong `lang-la-client-assets\res\x4\bg`); b00 768×192 là lớp đất fill (không có file nguồn tương ứng) | Có | ⚠️ Hai bản `_rong` **vẫn 1536**, chưa đúng 2976 theo DOT-2/KHUON-DE-BAI → núi lặp ~1,9 lần/khung. **Lock vẫn trỏ nền 343** (`map_nen_troi_1/may_1/may_2/nuihokage`) — lệch, xem §4 |
| Nhà (4) | `nha_nho` 576×576 · `nha_vua` 768×576 · `nha_lon` 768×768 · `nha_quan` 768×576 | LẮP (client) — tile 70–73 (x1 144×144, 192×144, 192×192, 192×144) | Có | Bản gốc trước sửa ở `goc-truoc-khi-sua` |
| Tile mặt đất (22) | `tile_dat_1/2`, `tile_co_mat`, `tile_co_goc_trai/phai`, `tile_vach_trai/phai`, `tile_doc_len/xuong`, `tile_datlat`, `tile_dat_lienmach_1/2`, `tile_dat_tang_giua/day`, `tile_vach_go_trai/phai`, `tile_chan_go`, `tile_ke_da`, `tile_bac_thang`, `tile_bac_thang_dinh`, `tile_co_sang_da`, `tile_da_sang_co` — đều 96×96 | LẮP (client) — tile 60–69 (đợt 1) và 83–88 (đợt 2) theo `tile_set_overlay` | Có | Client x1/t có 1$60–1$69, 1$83–1$88 (22–23/08) |
| Vật trang trí (10) | `tt_cong_torii` 480×480 · `tt_hangrao` 192×192 · `tt_denlong` 192×384 · `tt_gieng` 288×288 · `tt_thunggo` 192×192 · `tt_bienhieu` 192×288 · `tt_cay` 576×768 · `tt_cay_nho` 384×384 · `tt_bui_co` 192×96 · `tt_bui_ram` 288×192 | LẮP (client) — tile 74–82 (x1 120×120, 48×48, 48×96, 72×72, 48×48, 48×72, 144×192, 48×24, 96×48) | Có | Ánh xạ tile→tên chi tiết ở `SO-TAY-LAM-MAP.md` / `SOI-LAN-3.md` (chưa dò từng ô) |
| Bóng đổ (8) | `bong_cay` 576×96 · `bong_den` 192×96 · `bong_gieng` 288×96 · `bong_nha_lon` 768×96 · `bong_nha_nho` 576×96 · `bong_quan` 768×96 · `bong_thung` 192×96 · `bong_torii` 480×96 | Chưa xác minh đã cắm (không thấy cỡ tương ứng trong x1/t) | Có | DOT-2 ưu tiên 4 |
| Tiền cảnh (2) | `tc_canh_la` 768×384 · `tc_bui_truoc` 384×192 | Chưa xác minh | Có | DOT-2 §3.5 |
| Concept | `map-1-lang-la\concept-map1-langla.png` 2048×1024 | Tham chiếu | — | Concept đã duyệt về phong cách |

Client thật (`x1\t`) hiện có **48 png**: 1$7, 9, 11, 12, 13, 17 (cỏ Ngọc Rồng giữ lại), 1$39 (83 byte, trống), 1$40–51 (Lát 3 nhà cửa cũ), 1$60–88 (Naruto). `x2\t`, `x4\t` còn nguyên 1$1–1$38 tile Ngọc Rồng cũ (75 png). `x1\bg` có **b00–b03 và b10–b13 trùng nhau** (client tìm `typeBg=1`, LOI-DA-GAP). `x1\mapBackGround\500–507` (144×24… ngày 27/08) — **chưa rõ là gì**.

**3.5.2 Tile pixel-art lô 343 (`343-anh\09-ban-do`, 79 file)**

| # | Tên file | Nơi lưu | Cỡ file (px) | Làng | Trạng thái lắp | Dùng lại được không | Nội dung (rút gọn) |
|---|---|---|---|---|---|---|---|
| 1 | `map_bien_song_1_x4.png` | 343-anh/09-ban-do | 96×96 | chung | Kho | Hạn chế — tile pixel-art 96×96; map 1 đã chốt hand-painted (DOT-2 22/08) nên không trộn. Dùng làm tham chiếu bố cục/tên tile cho 4 làng còn lại | Sóng biển — sóng biển xanh lam cuộn có bọt trắng ở đầu ngọn, kiểu biển … |
| 2 | `map_bien_song_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Sóng biển — nhịp sóng và mảng bọt khác biến thể 1, tông lam sâu hơn. Bi… |
| 3 | `map_cau_cot_1_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Cột cầu — trụ/cột gỗ đứng đỡ cầu, nâu sẫm, đai buộc ngang, thớ gỗ dọc. … |
| 4 | `map_cau_cot_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Cột cầu — độ dày trụ và vị trí đai khác biến thể 1, rêu chân cột. Biến … |
| 5 | `map_cau_van_1_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Ván cầu gỗ — các thanh ván gỗ nâu ngang lát mặt cầu, vân gỗ rõ, khe hở … |
| 6 | `map_cau_van_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Ván cầu gỗ — bề rộng ván và vị trí khe khác biến thể 1. Biến thể 2. |
| 7 | `map_cau_van_3_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Ván cầu gỗ — ván ngả sậm/bạc màu, vài vết đinh và nứt. Biến thể 3. |
| 8 | `map_hv_ban_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Bàn học — mặt bàn gỗ học viện nâu vàng nhìn từ trên, cạnh viền bo, vài … |
| 9 | `map_hv_bang_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Bảng đen — mặt bảng trong lớp xanh đen (#1f2d24), viền khung gỗ, bề mặt… |
| 10 | `map_hv_bia_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Bia gỗ tập ném — tấm bia gỗ tròn tập ném shuriken/kunai, các vòng tròn … |
| 11 | `map_hv_chieu_1_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Chiếu tatami — chiếu rơm xanh cỏ úa dệt sọc, viền vải sẫm bao mép, kiểu… |
| 12 | `map_hv_chieu_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Chiếu tatami — hướng dệt và màu viền khác biến thể 1 để ghép ô so le. B… |
| 13 | `map_hv_ke_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Kệ cuộn giấy — kệ gỗ nhiều ngăn xếp các cuộn giấy/cuộn thuật (scroll) b… |
| 14 | `map_hv_sango_1_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Sàn gỗ — ván sàn gỗ nâu vàng bóng láng, thớ gỗ chạy đều, khe ghép mảnh;… |
| 15 | `map_hv_sango_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Sàn gỗ — hướng/bề rộng ván khác biến thể 1, vài ván sậm màu. Biến thể 2. |
| 16 | `map_hv_tuong_1_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Tường trong nhà — vách nội thất vữa/giấy sáng, khung gỗ nâu chia ô, tôn… |
| 17 | `map_hv_tuong_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Tường trong nhà — vị trí thanh khung và sắc độ khác biến thể 1. Biến th… |
| 18 | `map_la_co_1_x4.png` | 〃 | 96×96 | chung | Kho | Tham chiếu HIENCO-map_la_co_1 trong map-1-lang-la/anh-tham-chieu | Cỏ xanh ngắn — bãi cỏ trong làng, xanh lá tươi (#5fa34a) pha xanh đậm, … |
| 19 | `map_la_co_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Cỏ xanh ngắn — mật độ lá dày hơn, vài cụm cỏ cao nhô lên. Biến thể 2. |
| 20 | `map_la_co_3_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Cỏ xanh ngắn — xen mảng cỏ ngả vàng úa, hoa cỏ li ti. Biến thể 3. |
| 21 | `map_la_conglang_x4.png` | 〃 | 192×192 | Lá | Kho | 〃 | Cổng làng lớn — cổng gỗ hai trụ đỏ son to, xà ngang trên nối hai cột, k… |
| 22 | `map_la_cua_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Cửa gỗ trượt (fusuma/shoji) — khung gỗ nâu chia ô, ô giấy trắng ngà, ta… |
| 23 | `map_la_cuaso_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Cửa sổ giấy — ô cửa sổ shoji lưới ô vuông gỗ nâu, giấy trắng ngà mờ đục… |
| 24 | `map_la_dat_1_x4.png` | 〃 | 96×96 | Lá | LẮP (client, Lát 2) | README lang-la-client-assets: các ô đất map 1 = map_la_dat_1 | Đất nện nâu — sân đất khô Làng Lá, nâu đất (#6b4f34 → #8a6a45), lấm tấm… |
| 25 | `map_la_dat_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Đất nện nâu — cùng tông nâu, bố cục hạt/vết khác biến thể 1 để lát xen … |
| 26 | `map_la_dat_3_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Đất nện nâu — thêm vài mảng đất sáng/tối loang nhẹ, cụm sỏi con. Biến t… |
| 27 | `map_la_datlat_1_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Đường lát đá phẳng — phiến đá xám be lát khít làm đường chính trong làn… |
| 28 | `map_la_datlat_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Đường lát đá phẳng — kích cỡ/khớp mạch phiến đá khác biến thể 1, vài vi… |
| 29 | `map_la_datlat_3_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Đường lát đá phẳng — thêm rêu mọc kẽ đá, vết ẩm loang. Biến thể 3. |
| 30 | `map_la_denlong_x4.png` | 〃 | 96×96 | Lá | Kho | 〃 | Đèn lồng cam treo — đèn lồng giấy hình trụ cam đỏ treo dây, nan tre chi… |
| 31 | `map_la_hangrao_1_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Hàng rào tre — cọc tre vàng đứng buộc dây thừng ngang, đốt tre rõ, kiểu… |
| 32 | `map_la_hangrao_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Hàng rào tre — khoảng cách cọc/nút buộc khác biến thể 1, vài cọc ngả mà… |
| 33 | `map_la_maingoi_1_x4.png` | 〃 | 96×96 | Lá | Kho | 〃 | Mái ngói đỏ cam — ngói ống đỏ cam (#d9622b) đặc trưng mái nhà Làng Lá, … |
| 34 | `map_la_maingoi_2_x4.png` | 〃 | 96×96 | chung | Kho | Tham chiếu HIENCO-map_la_maingoi_2 | Mái ngói đỏ cam — góc/độ cong ngói khác biến thể 1, vài viên ngói ngả s… |
| 35 | `map_la_maingoi_3_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Mái ngói đỏ cam — thêm điểm bạc màu, rêu mép ngói cho biến thể mái cũ. … |
| 36 | `map_la_soi_1_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Sỏi rời — sỏi tròn nhỏ nhiều sắc xám/nâu rải ven đường, có bóng nhỏ dướ… |
| 37 | `map_la_soi_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Sỏi rời — cỡ sỏi to nhỏ pha trộn khác, mật độ dày hơn. Biến thể 2. |
| 38 | `map_la_soi_3_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Sỏi rời — xen vài viên sáng màu và đất lộ giữa các viên. Biến thể 3. |
| 39 | `map_la_tuonggo_1_x4.png` | 〃 | 96×96 | chung | Kho | Tham chiếu HIENCO-map_la_tuonggo_1 | Tường gỗ đứng — ván gỗ dọc nâu vàng, vân gỗ rõ, đầu đinh nhỏ, kiểu vách… |
| 40 | `map_la_tuonggo_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Tường gỗ đứng — bề rộng ván và vân gỗ khác biến thể 1. Biến thể 2. |
| 41 | `map_la_tuonggo_3_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Tường gỗ đứng — ván ngả nâu sậm/bạc, vài khe nứt. Biến thể 3. |
| 42 | `map_la_tuongvua_1_x4.png` | 〃 | 96×96 | Lá | Kho | 〃 | Tường vữa trắng — mảng vữa trắng ngà (#efe6d4) trong khung gỗ nâu viền,… |
| 43 | `map_la_tuongvua_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Tường vữa trắng — vị trí thanh gỗ khung khác, vài vết ố nhẹ trên vữa. B… |
| 44 | `map_nen_may_1_x4.png` | 〃 | 96×96 | chung | LẮP (lock b01) | Nguồn lớp nền b01 trong asset-lock | Mây xa — cụm mây trắng bồng bềnh trên nền trời lam, viền mây mềm nhiều … |
| 45 | `map_nen_may_2_x4.png` | 〃 | 96×96 | chung | LẮP (lock b02) | Nguồn lớp nền b02 trong asset-lock | Mây xa — hình khối mây và độ dày khác biến thể 1. Biến thể 2. |
| 46 | `map_nen_may_3_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Mây xa — mây thưa/tơi hơn, vài dải mây mỏng. Biến thể 3. |
| 47 | `map_nen_nuihokage_x4.png` | 〃 | 384×192 | Lá | LẮP (lock b03) | Nguồn lớp nền b03 trong asset-lock; client THẬT đã thay bằng nen_nuihokage (map-1 anh-da-tao) — xem §4 | Núi khắc mặt các đời thủ lĩnh (Núi Hokage) — vách núi đá xám khắc silho… |
| 48 | `map_nen_troi_1_x4.png` | 〃 | 96×96 | chung | LẮP (lock b00) | Đang là nguồn lớp nền b00 trong asset-lock (anh-game x1–x4) | Nền trời (chuyển sắc nhẹ) — trời xanh lam chuyển sắc nhẹ dùng làm nền x… |
| 49 | `map_nen_troi_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Nền trời (chuyển sắc nhẹ) — dải chuyển màu ngả xanh nhạt/ngả trắng khác… |
| 50 | `map_nen_troi_3_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Nền trời (chuyển sắc nhẹ) — tông ngả vàng cam nhẹ (hoàng hôn/bình minh)… |
| 51 | `map_nui_da_1_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Vách đá — vách đá xám nâu dựng đứng, tầng lớp địa chất ngang, khe nứt t… |
| 52 | `map_nui_da_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Vách đá — hướng tầng đá và vết nứt khác biến thể 1. Biến thể 2. |
| 53 | `map_nui_da_3_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Vách đá — điểm rêu/cỏ bám kẽ đá, tông ẩm hơn. Biến thể 3. |
| 54 | `map_nui_da_4_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Vách đá — đá sáng màu hơn, bề mặt lởm chởm gồ ghề. Biến thể 4. |
| 55 | `map_nui_dinh_1_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Đỉnh núi — chóp núi xa xanh xám phủ mờ hơi sương, sườn dốc, có thể loán… |
| 56 | `map_nui_dinh_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Đỉnh núi — dáng đỉnh và độ phủ tuyết/sương khác biến thể 1. Biến thể 2. |
| 57 | `map_rung_bui_1_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Bụi rậm — bụi cây thấp xanh, cành lá lộn xộn đan nhau, vài chấm quả/hoa… |
| 58 | `map_rung_bui_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Bụi rậm — hình khối bụi và mật độ cành khác biến thể 1. Biến thể 2. |
| 59 | `map_rung_bui_3_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Bụi rậm — xanh sẫm hơn, có cành khô nhô ra. Biến thể 3. |
| 60 | `map_rung_da_1_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Tảng đá rêu — đá xám (#7d7d78) phủ rêu xanh ẩm, bề mặt lồi lõm, kẽ nứt … |
| 61 | `map_rung_da_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Tảng đá rêu — hình khối đá và mảng rêu khác biến thể 1. Biến thể 2. |
| 62 | `map_rung_da_3_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Tảng đá rêu — rêu phủ nhiều hơn, vệt nước ẩm chảy. Biến thể 3. |
| 63 | `map_rung_re_1_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Rễ cây nổi — rễ cây nâu to nổi trên mặt đất rừng, uốn lượn vắt chéo, đấ… |
| 64 | `map_rung_re_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Rễ cây nổi — hướng uốn và độ dày rễ khác biến thể 1, điểm rêu bám. Biến… |
| 65 | `map_rung_tanla_1_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Tán lá dày — vòm lá xanh đậm nhiều lớp, cụm lá tròn chồng nhau, khe hở … |
| 66 | `map_rung_tanla_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Tán lá dày — bố cục cụm lá khác, vài điểm lá xanh sáng. Biến thể 2. |
| 67 | `map_rung_tanla_3_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Tán lá dày — xen lá ngả vàng, độ dày lớp khác. Biến thể 3. |
| 68 | `map_rung_tanla_4_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Tán lá dày — tông xanh lục sẫm, mảng bóng sâu hơn. Biến thể 4. |
| 69 | `map_rung_thancay_1_x4.png` | 〃 | 96×96 | Lá | Kho | 〃 | Thân cây lớn vỏ sần — thân cây cổ thụ nâu (#5a4327), vỏ nứt sần thành m… |
| 70 | `map_rung_thancay_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Thân cây lớn vỏ sần — hướng vân vỏ và vết nứt khác biến thể 1. Biến thể… |
| 71 | `map_rung_thancay_3_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Thân cây lớn vỏ sần — thêm hốc/u bướu gỗ và rêu dày hơn. Biến thể 3. |
| 72 | `map_rung_thancay_4_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Thân cây lớn vỏ sần — tông nâu sáng hơn, vỏ bong tróc lộ gỗ. Biến thể 4. |
| 73 | `map_song_bo_1_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Bờ đất ướt — mép nước giao đất: đất nâu ẩm sẫm chuyển sang bùn/cỏ, vệt … |
| 74 | `map_song_bo_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Bờ đất ướt — đường ranh nước-đất và sỏi ven bờ khác biến thể 1. Biến th… |
| 75 | `map_song_nuoc_1_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Mặt nước gợn — mặt nước xanh lam (#3a7fb0) gợn sóng lăn tăn, vệt phản c… |
| 76 | `map_song_nuoc_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Mặt nước gợn — hoa văn gợn sóng khác biến thể 1, điểm sáng lấp lánh dịc… |
| 77 | `map_song_nuoc_3_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Mặt nước gợn — tông lam sẫm hơn (nước sâu), gợn thưa. Biến thể 3. |
| 78 | `map_suongmu_1_x4.png` | 〃 | 96×96 | Sương Mù | Kho | 〃 | Lớp sương mờ (bán trong suốt) — màn sương trắng xám alpha thấp phủ mờ, … |
| 79 | `map_suongmu_2_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Lớp sương mờ (bán trong suốt) — mật độ/độ mờ sương khác biến thể 1 để c… |

Map 4 làng + nền (kich-thuoc-anh: 56 tên `map_{cat,smu,da,may}_*` 14/làng + 8 `map_nen_*` bien/dao/henui/nuida/nuimay/samac/vachda/vachnui 384×192): **chưa có ảnh**. Lưu ý khuôn mới `KHUON-DE-BAI-ANH-MAP.md` (27–33 ảnh/map, nền 2976 rộng, 4 tile đất khác vân) **khác hẳn** cỡ 96×96/384×192 cũ trong bảng này — đợt 4 nên đi theo khuôn mới.

### 3.6 Vật phẩm · trang bị · vũ khí

Game cần (vat-pham.yaml): **5 mẫu** — id 0 Thuốc hồi máu nhỏ, 1 Thuốc hồi chakra nhỏ, 2 Áo ninja Làng Lá, 3 Quần ninja Làng Lá, 4 Kunai luyện tập. **Tất cả `icon_id: 0`, `phan: -1`** (bảng SmallImage client rỗng, mảnh đồ chưa có trong part.bin). KH7 sẽ mở rộng.

**Vật phẩm `vp_` (100, 64×64) — dùng thẳng được làm icon túi đồ sau khi cấp id append-only:**

| # | Tên file | Nơi lưu | Cỡ file (px) | Làng | Trạng thái lắp | Dùng lại được không | Nội dung (rút gọn) |
|---|---|---|---|---|---|---|---|
| 1 | `vp_ao_t1_x4.png` | 343-anh/02-vat-pham | 64×64 | chung | Kho | Ứng với vat-pham.yaml id 2 'Áo ninja Làng Lá' — icon_id còn 0 | Áo khoác cổ cao vải thô, cam xỉn + navy đậm, sờn vai. 64×64. |
| 2 | `vp_ao_t2_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Áo cổ cao vải bền, cam tươi hơn, đường may chắc. 64×64. |
| 3 | `vp_ao_t3_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Áo cam-navy có nẹp da nâu ở vai và cổ. 64×64. |
| 4 | `vp_ao_t4_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Áo cam-navy lót lớp lưới thép ánh xám ở tay. 64×64. |
| 5 | `vp_ao_t5_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Áo cam-navy gắn mảnh giáp kim loại mỏng ngực. 64×64. |
| 6 | `vp_ao_t6_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Áo phối kiểu giáp Jonin xanh rêu, cổ cao, túi cuộn nhỏ. 64×64. |
| 7 | `vp_ao_t7_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Áo mảng giáp xám ANBU, khớp nối kim loại, cổ cao. 64×64. |
| 8 | `vp_ao_t8_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Áo choàng Kage trắng-cam tay rộng, viền chỉ vàng. 64×64. |
| 9 | `vp_bando_xusong_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Bản đồ Xứ Sóng: bản đồ ven biển cuốn tròn, mở hé một phần (a rolled coa… |
| 10 | `vp_bang_gac_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Cuộn băng gạc trắng có nhãn giấy trơn (rolled white bandage, plain tag,… |
| 11 | `vp_bangtran_t1_x4.png` | 〃 | 64×64 | Lá | Kho | 〃 | Băng trán dải vải thô xanh xỉn, tấm thép xỉn khắc lá xoáy Konoha. 64×64. |
| 12 | `vp_bangtran_t2_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Băng trán dải vải bền xanh, tấm thép sáng hơn, biểu tượng lá rõ. 64×64. |
| 13 | `vp_bangtran_t3_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Băng trán dải da nâu, tấm thép khắc lá. 64×64. |
| 14 | `vp_bangtran_t4_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Băng trán nền lưới thép, tấm biển kim loại khắc lá. 64×64. |
| 15 | `vp_bangtran_t5_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Băng trán tấm giáp nhẹ bóng, biểu tượng lá nổi. 64×64. |
| 16 | `vp_bangtran_t6_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Băng trán kiểu Jonin, dải xanh rêu, biển thép khắc lá. 64×64. |
| 17 | `vp_bangtran_t7_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Băng trán ANBU xám, biển kim loại tinh xảo khắc lá. 64×64. |
| 18 | `vp_bangtran_t8_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Băng trán Kage, biển mạ vàng khắc lá xoáy, dải trắng-cam. 64×64. |
| 19 | `vp_binh_luong_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Viên binh lương (soldier pill): viên tròn sẫm màu trong hộp thiếc nhỏ. … |
| 20 | `vp_binh_nuoc_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Bình nước bằng da (leather water flask). Da nâu, nắp và dây đeo. 64×64. |
| 21 | `vp_binh_tayky_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Bình tẩy ký: bình xịt xóa dấu ấn, màu trắng (a wash-bottle for erasing … |
| 22 | `vp_bot_no_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Bột nổ: túi nhỏ chứa bột đen đang đổ ra (small pouch of black powder, s… |
| 23 | `vp_bua_hosinh_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Bùa hồi sinh: dải giấy bùa có ký tự đen (revival talisman, paper strip … |
| 24 | `vp_chia_khoa_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Chìa khóa sắt cũ (old iron key). Kim loại xám, đầu chìa cổ điển. 64×64. |
| 25 | `vp_cuon_anthan_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Cuộn Ẩn thân: biểu tượng vòng xoáy đang mờ dần, dây xám (a fading swirl… |
| 26 | `vp_cuon_dat_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Cuộn Đất: cuộn giấy có dấu "Địa", tông nâu và xanh lá sẫm (scroll with … |
| 27 | `vp_cuon_hoadon_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Cuộn Hỏa độn: biểu tượng ngọn lửa, dây đỏ (a flame, red cord). 64×64. |
| 28 | `vp_cuon_hoiphuc_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Cuộn Hồi phục: biểu tượng bàn tay chữa lành, dây xanh dương nhạt (a hea… |
| 29 | `vp_cuon_kunaibay_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Cuộn Kunai bay: biểu tượng hai kunai bắt chéo, dây xám thép (crossed ku… |
| 30 | `vp_cuon_lienhoan_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Cuộn Liên hoàn: biểu tượng vòng xoáy ba nắm đấm, dây cam (a spiral of t… |
| 31 | `vp_cuon_phanthan_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Cuộn Phân thân: biểu tượng hai bóng người chồng lên nhau, dây trắng (tw… |
| 32 | `vp_cuon_thaythe_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Cuộn Thế thân: biểu tượng khúc gỗ bắt chéo với chiếc lá, dây nâu (a log… |
| 33 | `vp_cuon_trang_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Cuộn giấy trắng cuốn tròn, chưa viết (blank paper scroll, rolled). Kem-… |
| 34 | `vp_cuon_troi_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Cuộn Trời: cuộn giấy có dấu "Thiên", tông trắng và vàng kim (scroll wit… |
| 35 | `vp_cuon_tuchakra_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Cuộn Tụ chakra: biểu tượng dấu chân trên thân cây, dây xanh lá (a footp… |
| 36 | `vp_da_baove_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Đá bảo vệ: đá xám khắc rune (protective stone, grey with a rune). 64×64. |
| 37 | `vp_da_chakra_tho_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Đá chakra thô: khối pha lê xanh dương xỉn, đục (raw dull blue crystal c… |
| 38 | `vp_da_chakra_tinh_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Đá chakra tinh: pha lê xanh dương đã mài, phát sáng (polished glowing b… |
| 39 | `vp_da_maylan_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Đá may mắn: cỏ bốn lá bọc trong nhựa thông (four-leaf clover encased in… |
| 40 | `vp_da_nangcap_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Đá nâng cấp: cục đá có lõi cam phát sáng (upgrade stone, glowing orange… |
| 41 | `vp_da_thanhtay_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Đá thanh tẩy: mảnh pha lê trong, phớt trắng (a clear crystal shard, fai… |
| 42 | `vp_da_thu_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Da thú gấp lại (folded piece of animal hide). Nâu, vân da, mép lởm chởm… |
| 43 | `vp_dango_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Xiên dango 3 viên (three-ball dango skewer). Ba viên bột hồng/trắng/xan… |
| 44 | `vp_day_thep_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Dây thép: cuộn dây thép mảnh (coil of thin steel wire). Xám ánh kim, cu… |
| 45 | `vp_gang_t1_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Đôi găng hở ngón vải thô xanh xỉn, tấm mu bằng vải cứng. 64×64. |
| 46 | `vp_gang_t2_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Găng hở ngón vải bền xanh, tấm mu vải chắc. 64×64. |
| 47 | `vp_gang_t3_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Găng hở ngón da nâu, mu tay ốp da khâu viền. 64×64. |
| 48 | `vp_gang_t4_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Găng hở ngón, mu tay ốp lưới thép xám. 64×64. |
| 49 | `vp_gang_t5_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Găng hở ngón, tấm thép mu tay bắt sáng nhẹ. 64×64. |
| 50 | `vp_gang_t6_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Găng Jonin hở ngón xanh rêu, tấm thép mu tay. 64×64. |
| 51 | `vp_gang_t7_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Găng ANBU xám, ốp giáp mu tay khớp kim loại tinh xảo. 64×64. |
| 52 | `vp_gang_t8_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Găng Kage trắng-cam, tấm mu mạ vàng viền chỉ vàng. 64×64. |
| 53 | `vp_giay_t1_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Dép ninja hở ngón vải thô navy xỉn, quai sờn. 64×64. |
| 54 | `vp_giay_t2_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Dép hở ngón navy vải bền, quai chắc. 64×64. |
| 55 | `vp_giay_t3_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Dép hở ngón đế da nâu, quai da khâu viền. 64×64. |
| 56 | `vp_giay_t4_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Dép hở ngón có nẹp lưới thép quanh cổ chân. 64×64. |
| 57 | `vp_giay_t5_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Dép hở ngón gắn tấm giáp nhẹ ống chân. 64×64. |
| 58 | `vp_giay_t6_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Bốt kiểu Jonin xanh rêu hở ngón, quai đôi. 64×64. |
| 59 | `vp_giay_t7_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Bốt ANBU xám, ốp giáp ống chân kim loại. 64×64. |
| 60 | `vp_giay_t8_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Dép/bốt Kage trắng-cam, viền chỉ vàng cao cấp. 64×64. |
| 61 | `vp_giayphep_thi_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Giấy phép thi: đơn đăng ký chính thức có con dấu đỏ (an official applic… |
| 62 | `vp_go_cung_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Gỗ cứng: khúc gỗ ngắn sẫm màu, đặc (short log of dense dark wood). Vân … |
| 63 | `vp_hop_com_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Hộp cơm bento đóng kín, bọc khăn vải (a bento box, closed, with a cloth… |
| 64 | `vp_hop_gacha_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Hộp gacha: hộp triệu hồi niêm phong có dấu xoáy (a sealed summoning box… |
| 65 | `vp_hop_qua_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Hộp quà buộc nơ cam (gift box with orange ribbon). 64×64. |
| 66 | `vp_kiem_gay_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Kiếm gãy: phần mũi gãy của thanh đại đao khổng lồ (the broken tip of an… |
| 67 | `vp_long_chim_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Lông chim sẫm màu (dark feather). Đen/xám chuyển sắc. 64×64. |
| 68 | `vp_manh_matna_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Mảnh mặt nạ: mảnh vỡ của mặt nạ thợ săn trắng (a broken piece of a whit… |
| 69 | `vp_manh_thep_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Mảnh thép răng cưa lởm chởm (jagged steel shard). Xám kim loại, cạnh bé… |
| 70 | `vp_meo_tora_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Mèo Tora: mèo tam thể mặt cáu kỉnh, đeo nơ đỏ một bên tai (a grumpy cal… |
| 71 | `vp_mi_ramen_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Tô ramen bốc khói có lát naruto-maki (khoanh xoáy hồng-trắng). Nước dùn… |
| 72 | `vp_muc_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Nghiên mực nhỏ kèm bút lông (small ink stone with brush). Đá nghiên đen… |
| 73 | `vp_ngoc_gio_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Ngọc gió: viên ngọc xanh lá nhạt vân xoáy (pale green swirl gem). 64×64. |
| 74 | `vp_ngoc_lua_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Ngọc lửa: viên ngọc đỏ như than hồng (red ember gem). Đỏ-cam rực, ánh s… |
| 75 | `vp_ngoc_nuoc_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Ngọc nước: viên ngọc hình giọt nước màu xanh dương (blue droplet gem). … |
| 76 | `vp_ngoc_set_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Ngọc sét: viên ngọc vàng cạnh lởm chởm hình tia chớp (yellow jagged gem… |
| 77 | `vp_quan_t1_x4.png` | 〃 | 64×64 | chung | Kho | Ứng với vat-pham.yaml id 3 'Quần ninja Làng Lá' — icon_id còn 0 | Quần vải thô navy xỉn, băng trắng quấn cổ chân sờn. 64×64. |
| 78 | `vp_quan_t2_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Quần vải bền navy, băng quấn cổ chân gọn gàng. 64×64. |
| 79 | `vp_quan_t3_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Quần navy có nẹp da nâu bên hông, băng cổ chân. 64×64. |
| 80 | `vp_quan_t4_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Quần navy lót lưới thép xám ở đùi, băng cổ chân. 64×64. |
| 81 | `vp_quan_t5_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Quần navy gắn tấm giáp nhẹ che đùi, băng cổ chân. 64×64. |
| 82 | `vp_quan_t6_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Quần kiểu Jonin xanh rêu, túi hông, băng cổ chân. 64×64. |
| 83 | `vp_quan_t7_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Quần giáp ANBU xám, ốp gối kim loại, băng cổ chân. 64×64. |
| 84 | `vp_quan_t8_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Quần Kage trắng-cam, viền vàng, băng cổ chân sạch đẹp. 64×64. |
| 85 | `vp_ruong_go_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Rương gỗ nhỏ, đóng nắp (small wooden chest, closed). Gỗ nâu + đai/khóa … |
| 86 | `vp_ruong_sat_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Rương sắt bọc đai sắt, đóng kín (iron-banded chest, closed). Thân sẫm +… |
| 87 | `vp_soi_chakra_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Sợi chakra: cuộn chỉ xanh dương phát sáng (a coil of glowing blue threa… |
| 88 | `vp_the_dangky_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Thẻ đăng ký: thẻ gỗ khắc treo trên dây (a wooden registration tag on a … |
| 89 | `vp_the_doiten_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Thẻ đổi tên: thẻ trống treo trên dây (a blank name tag on a cord). 64×6… |
| 90 | `vp_thu_hokage_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Thư Hokage: lá thư niêm phong có dấu sáp (a sealed letter with a wax st… |
| 91 | `vp_thuoc_giai_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Chai thuốc giải màu xanh lá nút bần (green antidote bottle with cork). … |
| 92 | `vp_tui_ninja_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Túi ninja da (leather ninja pouch). Túi da nâu đeo hông đựng dụng cụ. 6… |
| 93 | `vp_tui_shuriken_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Túi shuriken: bao đùi (thigh holster pouch). Bao da/vải đeo đùi kiểu ni… |
| 94 | `vp_ve_dichuyen_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Vé dịch chuyển: bùa giấy có dấu đỏ (teleport ticket, paper talisman wit… |
| 95 | `vp_vien_chakra_1_x4.png` | 〃 | 64×64 | chung | Kho | Ứng với vat-pham.yaml id 1 'Thuốc hồi chakra nhỏ' — icon_id còn 0 | Viên thuốc xanh dương nhỏ bọc giấy (small blue pill in paper). Xanh cha… |
| 96 | `vp_vien_chakra_2_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Viên thuốc xanh dương cỡ vừa trong lọ thủy tinh (medium blue pill, glas… |
| 97 | `vp_vien_chakra_3_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Viên thuốc xanh dương lớn trong hũ chạm trổ (large blue pill, ornate ja… |
| 98 | `vp_vien_theluc_1_x4.png` | 〃 | 64×64 | chung | Kho | Ứng với vat-pham.yaml id 0 'Thuốc hồi máu nhỏ' — icon_id còn 0 | Viên thuốc đỏ nhỏ bọc trong giấy gói (small red pill in a paper wrapper… |
| 99 | `vp_vien_theluc_2_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Viên thuốc đỏ cỡ vừa trong lọ thủy tinh (medium red pill, glass vial). … |
| 100 | `vp_vien_theluc_3_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Viên thuốc đỏ lớn trong hũ chạm trổ (large red pill, ornate jar). Đỏ + … |

**Trang phục `tp_` (36 = 12 bộ × 3 mảnh) — KHÔNG dùng thẳng (xem §3.1):**

| # | Tên file | Nơi lưu | Cỡ file (px) | Làng | Trạng thái lắp | Dùng lại được không | Nội dung (rút gọn) |
|---|---|---|---|---|---|---|---|
| 1 | `tp_anbu_chan_x4.png` | 343-anh/03-trang-phuc | 128×128 | chung | Kho | KHÔNG dùng thẳng — mảnh vẽ rời 1 tư thế; client cần 3 đầu/17 thân/14 chân theo 33 khung. Dùng làm tham chiếu THIẾT KẾ bộ đồ | chân — Quần tối, giáp ống chân (shin guard) xám, giày/dép ANBU. |
| 2 | `tp_anbu_dau_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | đầu — Mặt nạ động vật màu trắng sứ có vằn đỏ/đen (mặt nạ ANBU, kiểu mèo… |
| 3 | `tp_anbu_than_x4.png` | 〃 | 128×160 | chung | Kho | 〃 | thân+tay — Giáp ngực tấm màu xám kim loại, miếng bảo hộ cánh tay (arm g… |
| 4 | `tp_anthan_chan_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | chân — Quần xám than tối, băng quấn cổ chân màu tối, dép ninja sẫm — to… |
| 5 | `tp_anthan_dau_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | đầu — Mũ trùm đầu (hood) kéo lên màu xám than + khăn vải che kín nửa dư… |
| 6 | `tp_anthan_than_x4.png` | 〃 | 128×160 | chung | Kho | 〃 | thân+tay — Áo choàng/áo bó toàn màu xám than tối, kín cổ, kín tay, khôn… |
| 7 | `tp_aomua_chan_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | chân — Gấu áo tơi rơm rủ xuống, quần tối bên trong, guốc/dép đi mưa. |
| 8 | `tp_aomua_dau_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | đầu — Nón lá rộng vành (kasa) màu rơm/vàng nâu che kín đầu, chỉ thấy nử… |
| 9 | `tp_aomua_than_x4.png` | 〃 | 128×160 | chung | Kho | 〃 | thân+tay — Áo tơi rơm (mino) rủ nhiều sợi rơm khô màu nâu vàng phủ vai … |
| 10 | `tp_chunin_giap_chan_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | chân — Quần xanh đen, băng quấn cổ chân, dép ninja. |
| 11 | `tp_chunin_giap_dau_x4.png` | 〃 | 128×128 | Lá | Kho | 〃 | đầu — Băng trán Làng Lá, tóc sẫm. |
| 12 | `tp_chunin_giap_than_x4.png` | 〃 | 128×160 | Lá | Kho | 〃 | thân+tay — Áo giáp chống đạn (flak jacket) màu xanh lá đặc trưng Làng L… |
| 13 | `tp_chunin_luoi_chan_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | chân — Quần tối màu, băng bó quấn quanh bắp/cổ chân, dép ninja hở ngón. |
| 14 | `tp_chunin_luoi_dau_x4.png` | 〃 | 128×128 | Lá | Kho | 〃 | đầu — Băng trán Làng Lá (kim loại + dây xanh đậm), tóc sẫm. |
| 15 | `tp_chunin_luoi_than_x4.png` | 〃 | 128×160 | chung | Kho | 〃 | thân+tay — Lớp áo lưới đen (mesh) ôm sát lộ ở cánh tay và cổ, khoác ngo… |
| 16 | `tp_genin_cai_chan_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | chân — Quần cam, bao đùi đựng vũ khí buộc ở đùi phải (túi vải nâu), dép… |
| 17 | `tp_genin_cai_dau_x4.png` | 〃 | 128×128 | Lá | Kho | 〃 | đầu — Băng trán Làng Lá (tấm kim loại khắc lá xoáy), tóc sẫm — giống bộ… |
| 18 | `tp_genin_cai_than_x4.png` | 〃 | 128×160 | chung | Kho | 〃 | thân+tay — Áo cam–xanh dương + một miếng giáp vai (kim loại/da xám ở va… |
| 19 | `tp_genin_co_chan_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | chân — Quần cam dài, băng vải quấn nhẹ cổ chân, dép ninja hở ngón xanh. |
| 20 | `tp_genin_co_dau_x4.png` | 〃 | 128×128 | Lá | Kho | 〃 | đầu — Băng trán Làng Lá (tấm kim loại bạc khắc biểu tượng lá xoáy, dây … |
| 21 | `tp_genin_co_than_x4.png` | 〃 | 128×160 | chung | Kho | 〃 | thân+tay — Áo khoác genin đặc trưng: cam là chính, phối xanh dương ở va… |
| 22 | `tp_hocvien_chan_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | chân — Quần vải lửng màu tối/xanh chàm, dép ninja hở ngón màu xanh dươn… |
| 23 | `tp_hocvien_dau_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | đầu — Tóc đen/nâu sẫm ngắn, mặt trẻ con học viện, KHÔNG đeo băng trán (… |
| 24 | `tp_hocvien_than_x4.png` | 〃 | 128×160 | chung | Kho | 〃 | thân+tay — Áo thun màu trung tính (be/xanh xám nhạt) đơn giản, tay ngắn… |
| 25 | `tp_jonin_giap_chan_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | chân — Quần xanh đen, bao đùi đựng vũ khí ở đùi, băng cổ chân, dép ninj… |
| 26 | `tp_jonin_giap_dau_x4.png` | 〃 | 128×128 | Lá | Kho | 〃 | đầu — Băng trán Làng Lá + khẩu trang vải sẫm che nửa dưới mặt (kiểu Kak… |
| 27 | `tp_jonin_giap_than_x4.png` | 〃 | 128×160 | chung | Kho | 〃 | thân+tay — Giáp chống đạn xanh lá có thêm túi đựng cuộn giấy (scroll po… |
| 28 | `tp_kage_chan_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | chân — Phần dưới áo choàng dài trắng phủ, thấp thoáng quần trắng và dép… |
| 29 | `tp_kage_dau_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | đầu — Tóc + mặt uy nghi lớn tuổi, KHÔNG đội mũ (mũ được cầm ở tay — nằm… |
| 30 | `tp_kage_than_x4.png` | 〃 | 128×160 | chung | Kho | 〃 | thân+tay — Áo choàng Kage trắng, viền/lót cam, tay áo rất rộng, một tay… |
| 31 | `tp_lehoi_chan_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | chân — Phần dưới yukata cam dài, guốc gỗ geta kiểu Nhật thò ra dưới gấu. |
| 32 | `tp_lehoi_dau_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | đầu — Tóc chải gọn cho lễ hội, không băng trán chiến đấu; có thể điểm p… |
| 33 | `tp_lehoi_than_x4.png` | 〃 | 128×160 | Lá | Kho | 〃 | thân+tay — Áo yukata cam phối họa tiết lá xoáy Làng Lá, cổ chéo vạt tru… |
| 34 | `tp_tuluyen_chan_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | chân — Quần tập rộng, vòng tạ nặng (weighted band) màu kim loại tối ở h… |
| 35 | `tp_tuluyen_dau_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | đầu — Tóc sẫm, khăn/dây vải buộc đầu kiểu tập luyện (hoặc băng trán vải… |
| 36 | `tp_tuluyen_than_x4.png` | 〃 | 128×160 | chung | Kho | 〃 | thân+tay — Áo ba lỗ/không tay màu trơn (trắng ngà/xanh xám), hai cẳng t… |

**Vũ khí `vk_` (20):**

| # | Tên file | Nơi lưu | Cỡ file (px) | Làng | Trạng thái lắp | Dùng lại được không | Nội dung (rút gọn) |
|---|---|---|---|---|---|---|---|
| 1 | `vk_bom_dinh_x4.png` | 343-anh/04-vu-khi | 256×256 (chuẩn 64×64) | chung | Kho | Có làm icon túi đồ/tham chiếu; vũ khí vẽ trên người phải là mảnh part riêng | Bom đinh / chông (makibishi) — Cụm chông thép 4 gai nhọn chĩa ra bốn hư… |
| 2 | `vk_bom_khoi_x4.png` | 〃 | 256×256 (chuẩn 64×64) | chung | Kho | 〃 | Bom khói — Quả cầu nhỏ màu đen, bề mặt kim loại có highlight tròn. Ngòi… |
| 3 | `vk_bua_lon_x4.png` | 〃 | 384×192 (chuẩn 96×48) | Sương Mù | Kho | 〃 | Đại đao bản lớn (Kubikiribōchō — của Zabuza) — Kiếm chặt khổng lồ dáng … |
| 4 | `vk_chuy_x4.png` | 〃 | 384×192 (chuẩn 96×48) | chung | Kho | 〃 | Chùy đá — Cán gỗ, đầu là khối đá thô màu xám nặng (stone head), bề mặt … |
| 5 | `vk_con_roi_tay_x4.png` | 〃 | 384×384 (chuẩn 96×96) | chung | Kho | 〃 | Đoạn tay con rối (rối của Sasori/Kankurō) — Một khúc cánh tay con rối b… |
| 6 | `vk_cung_x4.png` | 〃 | 384×192 (chuẩn 96×48) | chung | Kho | 〃 | Cung ngắn (recurve) — Cung gỗ ngắn cong ngược (recurve), dáng không căn… |
| 7 | `vk_cuon_no_x4.png` | 〃 | 256×256 (chuẩn 64×64) | chung | Kho | 〃 | Cuộn bùa nổ (explosive scroll) — Cuộn giấy trải hé, bề mặt in ký hiệu/c… |
| 8 | `vk_dao_doi_x4.png` | 〃 | 384×192 (chuẩn 96×48) | chung | Kho | 〃 | Song đao (cặp đao đôi) — Một cặp đoản đao giống hệt nhau bắt chéo chữ X… |
| 9 | `vk_day_thep_cuon_x4.png` | 〃 | 256×256 (chuẩn 64×64) | chung | Kho | 〃 | Cuộn dây thép ninja — Cuộn/lô dây thép quấn tròn nhiều vòng, ánh xám bạ… |
| 10 | `vk_giao_x4.png` | 〃 | 384×192 (chuẩn 96×48) | chung | Kho | 〃 | Giáo ngắn — Cán gỗ nâu dài, đầu là mũi giáo hình lá bằng thép (leaf-sha… |
| 11 | `vk_kiem_ngan_x4.png` | 〃 | 384×192 (chuẩn 96×48) | chung | Kho | 〃 | Kiếm ngắn cong — Lưỡi cong nhẹ kiểu wakizashi, thép sáng. Bên cạnh đặt … |
| 12 | `vk_kunai_ba_x4.png` | 〃 | 384×192 (chuẩn 96×48) | chung | Kho | 〃 | Kunai ba chấu (Hiraishin — của Minato) — Kunai ba lưỡi chĩa: lưỡi giữa … |
| 13 | `vk_kunai_no_x4.png` | 〃 | 384×192 (chuẩn 96×48) | chung | Kho | 〃 | Kunai gắn bùa nổ — Kunai thường + lá bùa giấy trắng ngà (explosive tag)… |
| 14 | `vk_kunai_thep_x4.png` | 〃 | 384×192 (chuẩn 96×48) | chung | Kho | 〃 | Kunai thép (bản tinh luyện) — Cùng dáng kunai thường nhưng thép sáng bó… |
| 15 | `vk_kunai_thuong_x4.png` | 〃 | 384×192 (chuẩn 96×48) | chung | Kho | Ứng với vat-pham.yaml id 4 'Kunai luyện tập' — icon_id còn 0; cỡ 96×48 không vừa ô 64×64 | Kunai thường — Dao ném hình thoi/lá cổ điển của ninja. Lưỡi thép xám bạ… |
| 16 | `vk_quat_gio_x4.png` | 〃 | 384×192 (chuẩn 96×48) | chung | Kho | 〃 | Quạt gió chiến (quạt của Temari) — Quạt gấp lớn xòe ra, nền màu kem/trắ… |
| 17 | `vk_senbon_x4.png` | 〃 | 384×192 (chuẩn 96×48) | chung | Kho | 〃 | Kim châm senbon — Bó ba cây kim thép mảnh, dài, hai đầu vót nhọn như ki… |
| 18 | `vk_shuriken_lon_x4.png` | 〃 | 384×384 (chuẩn 96×96) | chung | Kho | 〃 | Đại phi tiêu (Fūma shuriken — loại gấp lớn) — Shuriken khổng lồ 4 lưỡi … |
| 19 | `vk_shuriken_x4.png` | 〃 | 256×256 (chuẩn 64×64) | chung | Kho | 〃 | Phi tiêu (shuriken 4 cánh) — Ngôi sao ném thép xám, 4 cánh nhọn đối xứn… |
| 20 | `vk_tanto_x4.png` | 〃 | 384×192 (chuẩn 96×48) | chung | Kho | 〃 | Đoản kiếm tantō — Kiếm thẳng, ngắn, lưỡi thép thẳng có sống rõ. Chắn ta… |

Vật phẩm 4 làng (kich-thuoc-anh: 160 tên `vp_{cat,smu,da,may}_*`, 40/làng; 24 `tp_{cat,smu,da,may}_{tanbinh,chunin,jonin}_{dau,than}`): chưa có ảnh.

### 3.7 Chiêu + hiệu ứng

Game cần (chieu.yaml): **20 chiêu** (4/làng: cận chiến, tầm xa, hồi máu, tụ chakra), `icon_id 1–20`. mapping-id `skills` chỉ có **4 id 27079–27082** (pose_slot 0–3, `effect_id -1`, `arrow_id -1`) — tức **chưa chiêu nào có hiệu ứng**.

| Tên | Nơi | Cỡ | Trạng thái lắp | Dùng lại được không | Ghi chú |
|---|---|---|---|---|---|
| `skill_icons_leaf_4_chroma.png` | `anh-nguon\vertical-slice-v1` | 1983×793, 4 icon | **LẮP (lock + client)** → Small27079–27082 (x4 64×64) | Có | Cước Liên Hoàn · Hào Hoả Cầu · Y Liệu · Tụ Chakra (icon_id 1–4) |
| `skill_icons_sand_4_chroma.png` | 〃 | 1983×793 | Nguồn đã review, **chưa cấp id** | Có | icon_id 5–8 Làng Cát |
| `skill_icons_mist_4_chroma.png` | 〃 | 1983×793 | chưa cấp id | Có | icon_id 9–12 Sương Mù |
| `skill_icons_stone_4_chroma.png` | 〃 | 1983×793 | chưa cấp id | Có | icon_id 13–16 Đá |
| `skill_icons_cloud_4_chroma.png` | 〃 | 1983×793 | chưa cấp id | Có | icon_id 17–20 Mây |
| `hu_anthan_10f_chroma.png` | 〃 | 1774×887, 10 khung | Nguồn đã review (QA PASS), **không có output** nào dùng | Có | Hiệu ứng ẩn thân — chưa có effect_id |
| 10 `hu_*` (1200×480) | `343-anh\08-hieu-ung` | 1200×480 | Kho | Có thể, sau khi cắt khung + kiểm số khung | Mô tả manifest: "10 khung 240×240"; `_KHUON.md` §2: "5 khung 240×480" — **mâu thuẫn, phải mở ảnh kiểm**. Client cần `effectdata/<id>/img.png + data.bytes` |
| 20 tên `kn_*` (64×64) | `kich-thuoc-anh.json` | 64×64 | — | — | **Chưa có file `kn_*` nào**; nội dung đã có trong 5 sheet chroma ở trên — chỉ thiếu bước xuất tên + cấp id |

| # | Tên file | Nơi lưu | Cỡ file (px) | Làng | Trạng thái lắp | Dùng lại được không | Nội dung (rút gọn) |
|---|---|---|---|---|---|---|---|
| 1 | `hu_anthan_x4.png` | 343-anh/08-hieu-ung | 1200×480 | chung | Thay rồi | Đã có hu_anthan_10f_chroma (đã QA PASS, chưa gắn effect_id). Bản này là XAU-hu-anthan-qua-mo (mẫu xấu) | Ẩn thân / tàng hình. Một vùng lung linh hình dáng cơ thể bị bóp méo như… |
| 2 | `hu_bangkinh_x4.png` | 〃 | 1200×480 | chung | Kho | 〃 | Băng kính (Ma Kính Băng Tinh — Haku). Những tấm gương băng mỏng hình th… |
| 3 | `hu_chimang_x4.png` | 〃 | 1200×480 | chung | Kho | 〃 | Chí mạng (đòn critical). Một luồng loé va chạm trắng–cam sắc nét với cá… |
| 4 | `hu_choang_x4.png` | 〃 | 1200×480 | chung | Kho | 〃 | Choáng (stun). Ba ngôi sao vàng nhỏ quay quanh thành vòng tròn, lặp mượ… |
| 5 | `hu_hoadon_x4.png` | 〃 | 1200×480 | chung | Kho | Ứng với chiêu icon_id 2 'Hoả Độn: Hào Hoả Cầu' — hiệu ứng chiêu chưa có effect_id | Hỏa độn — Cầu lửa (Katon / Goukakyuu). Một đốm than nhỏ lớn dần thành q… |
| 6 | `hu_hoiphuc_x4.png` | 〃 | 1200×480 | chung | Kho | Ứng với 5 chiêu hồi máu (loai 2) — chưa gắn | Hồi phục (trị liệu). Ánh sáng xanh lá dịu dâng lên theo xoắn ốc nhẹ, kè… |
| 7 | `hu_kunaibay_x4.png` | 〃 | 1200×480 | Lá | Kho | Ứng với chiêu tầm xa — chưa gắn | Kunai bay. Một vệt mờ chuyển động (motion-blur) của lưỡi dao ném đi, kh… |
| 8 | `hu_phanthan_x4.png` | 〃 | 1200×480 | chung | Kho | 〃 | Phân thân (Kage Bunshin — Nhân đôi bóng). Một cụm khói trắng bung ra từ… |
| 9 | `hu_thaythe_x4.png` | 〃 | 1200×480 | chung | Kho | 〃 | Thuật Thay thế (Kawarimi). Một xoáy khói trắng cuộn lên hé lộ hình một … |
| 10 | `hu_tuchakra_x4.png` | 〃 | 1200×480 | chung | Kho | Ứng với 5 chiêu tụ chakra (loai 3) — chưa gắn | Tụ chakra. Chakra xanh dương tụ lại tại một điểm, cuộn xoáy dâng lên, đ… |

### 3.8 Giao diện

| Tên | Nơi | Cỡ | Trạng thái lắp | Dùng lại được không | Ghi chú |
|---|---|---|---|---|---|
| `hud_hp_chakra_4_chroma.png` | `anh-nguon\vertical-slice-v1` | 1774×887, 4 ô | **LẮP (lock)** → `mainImage/naruto-hp-bg/fill, naruto-chakra-bg/fill` ×4 zoom; evidence `hud-golden-25-50-100-x4.png` | Có | Client thật **không có** `mainimage\naruto-*` (x2/x4 mainimage chỉ có btn0/btn1 Ngọc Rồng cũ 2024–03/2026) — xem §4 |
| `ui_logo_game_x4.png` 512×256 · `ui_khung_avatar_x4` 96×96 · `ui_khung_hoithoai_x4` 480×160 · `ui_nut_thuong_x4` 96×48 · `ui_nut_nhan_x4` 96×48 · `ui_nut_khoa_x4` 96×48 · `ui_o_vatpham_x4` 64×64 · `ui_o_dangchon_x4` 64×64 | `D:\NRO_GoiBanGiao\my thuat\dot1-nhom1-duyet\dot1-nhom1-duyet\01-giao-dien` | như tên | **LẮP (lock)**: logo `direct` → `/Logo/1.png`; 2 nút `slice_group` → `btn0/btn1 left-mid-right`; 5 cái `overlay_required` → `/naruto/ui/*.png` (client **phải sửa code** mới dùng) | Có | Client thật: `x4\Logo\1.png` là **wordmark 599×136** (05/08), không phải bản 512×256 này; `naruto\ui\` không tồn tại |
| `ui_thanh_chakra_x4`, `ui_thanh_theluc_x4` 240×32 (gộp) | 〃 | 240×32 | Kho | Không — bản gộp `XAU-thanh-gop-san-khong-dung-duoc` | Bản tách `_nen`/`_fill` ở `my thuat\_v2-bung\01-giao-dien` (4 file 240×32); HUD chroma đã thay |
| 50 `ui_*` | `343-anh\01-giao-dien` | 64×64 (33), 48×48 (4), 96×96, 128×128, 192×192, 240×24, 400×48, 400×120 (2), 400×480, 480×400, 480×480 (2), 1024×576 (2) | Kho | Có — icon/khung tĩnh; soi phong cách | Chưa có màn hình nào của client bind tới |

| # | Tên file | Nơi lưu | Cỡ file (px) | Làng | Trạng thái lắp | Dùng lại được không | Nội dung (rút gọn) |
|---|---|---|---|---|---|---|---|
| 1 | `ui_bac_anbu_x4.png` | 343-anh/01-giao-dien | 64×64 | chung | Kho | Có — icon/khung tĩnh đúng kiểu client cần; soi lại phong cách (lô 343 vẽ pixel-art, phong cách khoá nay là hand-painted) | Tier 6 — ANBU: bóng mặt nạ thú màu trắng (white animal mask silhouette). |
| 2 | `ui_bac_chunin_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Tier 3 — Chunin: tấm kim loại với cổ áo gile xanh lá (green vest collar… |
| 3 | `ui_bac_genin_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Tier 2 — Genin: tấm kim loại có một khía (one notch). |
| 4 | `ui_bac_hiennhan_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Tier 8 — Hiền nhân (Sage): mô-típ vạch quanh mắt cóc màu cam (toad-eye … |
| 5 | `ui_bac_hocvien_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Tier 1 — Học viên: dải vải băng trán trơn, không có tấm kim loại. |
| 6 | `ui_bac_jonin_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Tier 4 — Jonin: tấm kim loại viền bạc (silver trim). |
| 7 | `ui_bac_kage_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Tier 7 — Kage: bóng mũ chóp nón Kage màu vàng (conical hat silhouette, … |
| 8 | `ui_bac_lucdao_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Tier 9 — Lục Đạo: chín viên magatama xếp thành vòng, vàng trên nền đen. |
| 9 | `ui_bac_tinhanh_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Tier 5 — Tinh anh: tấm kim loại với hai cánh hai bên (two side wings). |
| 10 | `ui_he_hoa_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Hệ Hỏa: ngọn lửa cách điệu, đỏ-cam. |
| 11 | `ui_he_loi_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Hệ Lôi: tia sét (lightning bolt), vàng-trắng. |
| 12 | `ui_he_phong_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Hệ Phong: xoáy gió lượn (wind curl), xanh lá nhạt. |
| 13 | `ui_he_tho_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Hệ Thổ: cụm đá (rock cluster), nâu. |
| 14 | `ui_he_thuy_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Hệ Thủy: giọt nước kèm gợn sóng, xanh dương. |
| 15 | `ui_joystick_nen_x4.png` | 〃 | 192×192 | chung | Kho | 〃 | Đế joystick ảo: vòng tròn tối bán trong suốt (translucent), có texture … |
| 16 | `ui_joystick_num_x4.png` | 〃 | 96×96 | chung | Kho | 〃 | Núm joystick: vòm tròn màu cam (orange dome) có highlight, đổ bóng, đặt… |
| 17 | `ui_khung_canhbao_x4.png` | 〃 | 400×120 | chung | Kho | 〃 | Banner cảnh báo. Cùng khổ với thông báo nhưng viền sáng đỏ, có dấu hiệu… |
| 18 | `ui_khung_nhiemvu_x4.png` | 〃 | 400×480 | chung | Kho | 〃 | Panel nhiệm vụ hình cuộn giấy trải ra (unrolled scroll): hai trục cuộn … |
| 19 | `ui_khung_ruong_x4.png` | 〃 | 480×480 | chung | Kho | 〃 | Panel kho/rương chứa đồ. Giống túi đồ nhưng nặng đô hơn: bốn góc bọc sắ… |
| 20 | `ui_khung_shop_x4.png` | 〃 | 480×400 | chung | Kho | 〃 | Panel cửa hàng. Khung gỗ tối viền cam, phía trên có dải vải/băng-rôn tr… |
| 21 | `ui_khung_thongbao_x4.png` | 〃 | 400×120 | chung | Kho | 〃 | Banner thông báo nhỏ. Dải ngang gỗ/giấy, viền phát sáng cam dịu quanh m… |
| 22 | `ui_khung_tuido_x4.png` | 〃 | 480×480 | chung | Kho | 〃 | Panel túi đồ (inventory). Khung gỗ sơn mài tối viền cam, ruột nền giấy … |
| 23 | `ui_lang_cat_x4.png` | 〃 | 64×64 | Cát | Kho | 〃 | Làng Cát (Suna): đồng hồ cát / xoáy hình đồng hồ cát (hourglass-spiral). |
| 24 | `ui_lang_da_x4.png` | 〃 | 64×64 | Đá | Kho | 〃 | Làng Đá (Iwa): hai phiến đá (two rock slabs). |
| 25 | `ui_lang_la_x4.png` | 〃 | 64×64 | Lá | Kho | 〃 | Làng Lá (Konoha): chiếc lá xoáy cách điệu trên tấm kim loại. |
| 26 | `ui_lang_may_x4.png` | 〃 | 64×64 | Mây | Kho | 〃 | Làng Mây (Kumo): hình khối mây (cloud shape). |
| 27 | `ui_lang_suong_x4.png` | 〃 | 64×64 | Sương Mù | Kho | 〃 | Làng Sương (Kiri): bốn đường sương lượn sóng (four wavy mist lines). |
| 28 | `ui_muiten_len_x4.png` | 〃 | 48×48 | chung | Kho | 〃 | Mũi tên lên, dày, cam. |
| 29 | `ui_muiten_phai_x4.png` | 〃 | 48×48 | chung | Kho | 〃 | Mũi tên phải, dày, cam. |
| 30 | `ui_muiten_trai_x4.png` | 〃 | 48×48 | chung | Kho | 〃 | Mũi tên trái, dày, cam. |
| 31 | `ui_muiten_xuong_x4.png` | 〃 | 48×48 | chung | Kho | 〃 | Mũi tên xuống, dày, cam. |
| 32 | `ui_nen_chonlang_x4.png` | 〃 | 1024×576 | chung | Kho | 〃 | Nền màn chọn làng: tấm bản đồ năm nước lớn trải trên bàn gỗ tối, ánh đè… |
| 33 | `ui_nen_dangnhap_x4.png` | 〃 | 1024×576 | chung | Kho | 〃 | Nền màn đăng nhập: một làng ninja yên tĩnh lúc hoàng hôn nhìn từ sườn đ… |
| 34 | `ui_nut_bando_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Nút bản đồ: tấm bản đồ gấp (folded map). |
| 35 | `ui_nut_caidat_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Nút cài đặt: bánh răng lai shuriken (gear-shuriken hybrid) — vừa là răn… |
| 36 | `ui_nut_chat_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Nút chat: bong bóng thoại (speech bubble). |
| 37 | `ui_nut_danh_x4.png` | 〃 | 128×128 | chung | Kho | 〃 | Nút tấn công lớn (round attack button, 128×128): nút tròn to màu cam nổ… |
| 38 | `ui_nut_dong_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Nút đóng: dấu X trên nút tròn gỗ tối. |
| 39 | `ui_nut_kynang_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Nút kỹ năng: dấu kết ấn tay (hand-seal) — bàn tay bắt ấn ninjutsu. |
| 40 | `ui_nut_nhiemvu_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Nút nhiệm vụ: cuộn giấy nhỏ (scroll). |
| 41 | `ui_nut_thu_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Nút thư: lá thư niêm phong (sealed letter) có dấu sáp. |
| 42 | `ui_nut_toidoi_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Nút tổ đội: ba hình người nhỏ (three small figures) — đội nhóm. |
| 43 | `ui_nut_tuido_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Nút túi đồ: túi/bao ninja (ninja pouch) đeo hông. |
| 44 | `ui_nut_xephang_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Nút xếp hạng: bục podium (bậc 1-2-3). |
| 45 | `ui_thanh_kinhnghiem_x4.png` | 〃 | 240×24 | chung | Kho | 〃 | Thanh kinh nghiệm mảnh. Đế rãnh tối, phần đầy màu vàng gold sáng lấp lá… |
| 46 | `ui_thanh_nap_x4.png` | 〃 | 400×48 | chung | Kho | 〃 | Thanh loading/nạp. Đế rãnh tối dài, phần đầy màu cam, ở mép của phần đầ… |
| 47 | `ui_tien_ngoc_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Ngọc (tiền nạp): viên đá quý cyan cắt giác (faceted gemstone), xanh lam… |
| 48 | `ui_tien_ryo_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Tiền Ryo: đồng xu vàng cổ hình bầu dục có lỗ vuông ở giữa, xếp chồng ng… |
| 49 | `ui_tien_thoivang_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Thỏi vàng: thanh vàng nén (gold ingot) hình thang, ánh kim, bóng đổ. |
| 50 | `ui_tien_xu_x4.png` | 〃 | 64×64 | chung | Kho | 〃 | Xu sự kiện: token nhỏ màu tím có dấu xoáy (spiral mark) — vòng xoáy kiể… |

### 3.9 Chân dung

| Tên | Nơi | Cỡ | Trạng thái lắp | Dùng lại được không | Ghi chú |
|---|---|---|---|---|---|
| `iruka_portrait_chroma.png` → Small27078 | `vertical-slice-v1` → client | 200×200 (x4) | **LẮP (lock + client)** | Có | Chuẩn cỡ chân dung client = 200×200 |
| 15 `cd_*` | `343-anh\07-npc` | 800×800 (chuẩn 200×200; `anh-game\x4` đúng 200) | Kho (cd_iruka: thay rồi) | Có về kỹ thuật; soi phong cách | Cùng bộ với 15 `npc_*` |

| # | Tên file | Nơi lưu | Cỡ file (px) | Làng | Trạng thái lắp | Dùng lại được không | Nội dung (rút gọn) |
|---|---|---|---|---|---|---|---|
| 1 | `cd_ba_vukhi_x4.png` | 343-anh/07-npc | 800×800 (chuẩn 200×200) | chung | Kho | Có về kỹ thuật — chân dung 200×200 đúng cỡ client (như 27078); soi phong cách pixel-art vs hand-painted | Cận mặt bà lão: nếp nhăn sâu, búi tóc bạc, ánh mắt sắc, cổ áo kimono tố… |
| 2 | `cd_chu_ramen_x4.png` | 〃 | 800×800 (chuẩn 200×200) | chung | Kho | Có về kỹ thuật — chân dung 200×200 đúng cỡ client (như 27078); soi phong cách pixel-art vs hand-painted | Cận mặt chú bếp: mặt tròn phúc hậu, khăn trắng buộc trán, cười rạng rỡ,… |
| 3 | `cd_co_be_hoa_x4.png` | 〃 | 800×800 (chuẩn 200×200) | chung | Kho | Có về kỹ thuật — chân dung 200×200 đúng cỡ client (như 27078); soi phong cách pixel-art vs hand-painted | Cận mặt bé gái: má bầu, mắt to, vài bông hoa nhô lên từ giỏ ở cạnh dưới… |
| 4 | `cd_hokage3_x4.png` | 〃 | 800×800 (chuẩn 200×200) | chung | Kho | Có về kỹ thuật — chân dung 200×200 đúng cỡ client (như 27078); soi phong cách pixel-art vs hand-painted | Cận mặt ông già: nếp nhăn, râu bạc, ánh mắt hiền, vành nón đỏ-trắng che… |
| 5 | `cd_iruka_x4.png` | 〃 | 800×800 (chuẩn 200×200) | chung | Thay rồi | Client dùng iruka_portrait_chroma → sprite 27078 (200×200) | Cận mặt thầy Iruka: vết sẹo ngang mũi rõ nét, băng trán bạc dải xanh, b… |
| 6 | `cd_kakashi_x4.png` | 〃 | 800×800 (chuẩn 200×200) | chung | Kho | Có về kỹ thuật — chân dung 200×200 đúng cỡ client (như 27078); soi phong cách pixel-art vs hand-painted | Cận mặt Kakashi: một mắt đen lộ ra dưới băng trán nghiêng, khẩu trang c… |
| 7 | `cd_linh_gac_x4.png` | 〃 | 800×800 (chuẩn 200×200) | chung | Kho | Có về kỹ thuật — chân dung 200×200 đúng cỡ client (như 27078); soi phong cách pixel-art vs hand-painted | Cận mặt lính gác: biểu cảm ngáp/buồn ngủ, băng trán bạc, cổ vest xanh. … |
| 8 | `cd_nguoi_dua_thu_x4.png` | 〃 | 800×800 (chuẩn 200×200) | chung | Kho | Có về kỹ thuật — chân dung 200×200 đúng cỡ client (như 27078); soi phong cách pixel-art vs hand-painted | Cận mặt giao liên: nét nhanh nhẹn, băng trán, quai túi vắt chéo vai lộ … |
| 9 | `cd_nguoi_giao_nv_x4.png` | 〃 | 800×800 (chuẩn 200×200) | chung | Kho | Có về kỹ thuật — chân dung 200×200 đúng cỡ client (như 27078); soi phong cách pixel-art vs hand-painted | Cận mặt nhân viên: nét điềm đạm, băng trán, sau vai thấp thoáng cuộn gi… |
| 10 | `cd_nha_chep_cuon_x4.png` | 〃 | 800×800 (chuẩn 200×200) | chung | Kho | Có về kỹ thuật — chân dung 200×200 đúng cỡ client (như 27078); soi phong cách pixel-art vs hand-painted | Cận mặt nhà chép cuộn: gương mặt gầy trầm ngâm, bút lông gài sau tai lộ… |
| 11 | `cd_ong_tazuna_x4.png` | 〃 | 800×800 (chuẩn 200×200) | chung | Kho | Có về kỹ thuật — chân dung 200×200 đúng cỡ client (như 27078); soi phong cách pixel-art vs hand-painted | Cận mặt ông Tazuna: mặt khắc khổ, râu ria lởm chởm, vành nón rơm, kính … |
| 12 | `cd_thay_hocvien_x4.png` | 〃 | 800×800 (chuẩn 200×200) | chung | Kho | Có về kỹ thuật — chân dung 200×200 đúng cỡ client (như 27078); soi phong cách pixel-art vs hand-painted | Cận mặt thầy: kính vuông phản sáng nhẹ (chấm highlight, không phải viền… |
| 13 | `cd_tho_ren_x4.png` | 〃 | 800×800 (chuẩn 200×200) | chung | Kho | Có về kỹ thuật — chân dung 200×200 đúng cỡ client (như 27078); soi phong cách pixel-art vs hand-painted | Cận mặt thợ rèn: mặt vuông rắn rỏi, mồ hôi (chấm highlight nhỏ), khăn q… |
| 14 | `cd_trongtai_thi_x4.png` | 〃 | 800×800 (chuẩn 200×200) | chung | Kho | Có về kỹ thuật — chân dung 200×200 đúng cỡ client (như 27078); soi phong cách pixel-art vs hand-painted | Cận mặt trọng tài: vẻ nghiêm nghị mệt mỏi, băng trán, cổ áo đồng phục t… |
| 15 | `cd_y_ta_x4.png` | 〃 | 800×800 (chuẩn 200×200) | chung | Kho | Có về kỹ thuật — chân dung 200×200 đúng cỡ client (như 27078); soi phong cách pixel-art vs hand-painted | Cận mặt y tá: gương mặt hiền, mũ y tá, cổ áo đồng phục nhạt. Đổ bóng nh… |

15 NPC đứng nguyên con (format cũ — tham chiếu để vẽ master):

| # | Tên file | Nơi lưu | Cỡ file (px) | Làng | Trạng thái lắp | Dùng lại được không | Nội dung (rút gọn) |
|---|---|---|---|---|---|---|---|
| 1 | `npc_ba_vukhi_x4.png` | 343-anh/07-npc | 384×512 (chuẩn 96×128) | chung | Kho | KHÔNG dùng thẳng — NPC client = 3 mảnh đầu/thân/chân (cf 0/1) + chân dung; ảnh này làm THAM CHIẾU để vẽ master_<npc>_x4_chroma | Bà già chủ tiệm vũ khí, tóc bạc búi cao, hai tay khoanh trước ngực, mặc… |
| 2 | `npc_chu_ramen_x4.png` | 〃 | 384×512 (chuẩn 96×128) | chung | Kho | KHÔNG dùng thẳng — NPC client = 3 mảnh đầu/thân/chân (cf 0/1) + chân dung; ảnh này làm THAM CHIẾU để vẽ master_<npc>_x4_chroma | Đầu bếp quán mì trung niên vui tính, khăn vải trắng buộc trán, mặc tạp … |
| 3 | `npc_co_be_hoa_x4.png` | 〃 | 384×512 (chuẩn 96×128) | chung | Kho | KHÔNG dùng thẳng — NPC client = 3 mảnh đầu/thân/chân (cf 0/1) + chân dung; ảnh này làm THAM CHIẾU để vẽ master_<npc>_x4_chroma | Bé gái nhỏ ôm giỏ hoa trước bụng, váy/kimono trẻ em màu tươi, tóc buộc,… |
| 4 | `npc_hokage3_x4.png` | 〃 | 384×512 (chuẩn 96×128) | chung | Kho | KHÔNG dùng thẳng — NPC client = 3 mảnh đầu/thân/chân (cf 0/1) + chân dung; ảnh này làm THAM CHIẾU để vẽ master_<npc>_x4_chroma | Ông già thủ lĩnh mặc áo choàng trắng viền đỏ/cam, đội nón chóp (conical… |
| 5 | `npc_iruka_x4.png` | 〃 | 384×512 (chuẩn 96×128) | chung | Thay rồi | Client dùng master_iruka_x4_chroma → part 2189/2190/2191. Ảnh này chỉ còn là HIENCO-npc_iruka tham chiếu | Thầy giáo trẻ, tóc nâu buộc đuôi ngựa dựng, vết sẹo ngang sống mũi, mặc… |
| 6 | `npc_kakashi_x4.png` | 〃 | 384×512 (chuẩn 96×128) | chung | Kho | KHÔNG dùng thẳng — NPC client = 3 mảnh đầu/thân/chân (cf 0/1) + chân dung; ảnh này làm THAM CHIẾU để vẽ master_<npc>_x4_chroma | Jonin tóc bạc dựng lệch, khẩu trang che nửa mặt dưới, băng trán nghiêng… |
| 7 | `npc_linh_gac_x4.png` | 〃 | 384×512 (chuẩn 96×128) | Lá | Kho | KHÔNG dùng thẳng — NPC client = 3 mảnh đầu/thân/chân (cf 0/1) + chân dung; ảnh này làm THAM CHIẾU để vẽ master_<npc>_x4_chroma | Lính gác cổng vẻ chán chường, mặc vest xanh rêu chunin, tựa người vào c… |
| 8 | `npc_nguoi_dua_thu_x4.png` | 〃 | 384×512 (chuẩn 96×128) | chung | Kho | KHÔNG dùng thẳng — NPC client = 3 mảnh đầu/thân/chân (cf 0/1) + chân dung; ảnh này làm THAM CHIẾU để vẽ master_<npc>_x4_chroma | Giao liên đang sải bước (mid-stride), đeo túi vải đựng cuộn giấy chéo v… |
| 9 | `npc_nguoi_giao_nv_x4.png` | 〃 | 384×512 (chuẩn 96×128) | chung | Kho | KHÔNG dùng thẳng — NPC client = 3 mảnh đầu/thân/chân (cf 0/1) + chân dung; ảnh này làm THAM CHIẾU để vẽ master_<npc>_x4_chroma | Nhân viên bàn giao nhiệm vụ ngồi/đứng giữa các chồng cuộn giấy (scroll)… |
| 10 | `npc_nha_chep_cuon_x4.png` | 〃 | 384×512 (chuẩn 96×128) | chung | Kho | KHÔNG dùng thẳng — NPC client = 3 mảnh đầu/thân/chân (cf 0/1) + chân dung; ảnh này làm THAM CHIẾU để vẽ master_<npc>_x4_chroma | Thư lại gầy, ngón tay dính mực, cây bút lông gài sau tai, áo dài kiểu h… |
| 11 | `npc_ong_tazuna_x4.png` | 〃 | 384×512 (chuẩn 96×128) | chung | Kho | Là TOT-npc-ong-tazuna.png (mẫu tốt về độ chi tiết) | Ông già dựng cầu dày dạn phong sương, đội nón rơm, tay cầm chai rượu sa… |
| 12 | `npc_thay_hocvien_x4.png` | 〃 | 384×512 (chuẩn 96×128) | chung | Kho | KHÔNG dùng thẳng — NPC client = 3 mảnh đầu/thân/chân (cf 0/1) + chân dung; ảnh này làm THAM CHIẾU để vẽ master_<npc>_x4_chroma | Giáo viên học viện nghiêm khắc, đeo kính vuông, tay cầm bảng kẹp giấy (… |
| 13 | `npc_tho_ren_x4.png` | 〃 | 384×512 (chuẩn 96×128) | chung | Kho | KHÔNG dùng thẳng — NPC client = 3 mảnh đầu/thân/chân (cf 0/1) + chân dung; ảnh này làm THAM CHIẾU để vẽ master_<npc>_x4_chroma | Thợ rèn vạm vỡ, tạp dề da, tay cầm búa rèn, cơ bắp lộ, băng trán vải qu… |
| 14 | `npc_trongtai_thi_x4.png` | 〃 | 384×512 (chuẩn 96×128) | chung | Kho | KHÔNG dùng thẳng — NPC client = 3 mảnh đầu/thân/chân (cf 0/1) + chân dung; ảnh này làm THAM CHIẾU để vẽ master_<npc>_x4_chroma | Trọng tài kỳ thi Chunin mặc đồng phục tối màu (đen/xanh than), một tay … |
| 15 | `npc_y_ta_x4.png` | 〃 | 384×512 (chuẩn 96×128) | chung | Kho | KHÔNG dùng thẳng — NPC client = 3 mảnh đầu/thân/chân (cf 0/1) + chân dung; ảnh này làm THAM CHIẾU để vẽ master_<npc>_x4_chroma | Y tá bệnh viện mặc đồng phục màu nhạt (trắng/hồng nhạt), tay cầm bảng k… |

### 3.10 Thương hiệu · logo

| Tên | Nơi | Cỡ | Trạng thái lắp | Dùng lại | Ghi chú |
|---|---|---|---|---|---|
| Wordmark NARUTO `Logo\1.png` | `lang-la-client-assets\Logo\1.png` → `D:\NRO_Project\Assets\Logo\1.png` (nguồn) → `Resources\res\x4\Logo\1.png` | 599×136 | **LẮP (client)** 04–05/08 | Có | Sửa phải sửa **nguồn `Assets\Logo\1.png`** (bộ nhớ `logo-resource-sync`) |
| `ui_logo_game_x4.png` | `my thuat\dot1-nhom1-duyet\...\01-giao-dien` | 512×256 | LẮP (lock, mode direct) — **nhưng client đang là wordmark** | Có | Lệch lock ↔ client (§4) |
| `TOT-naruto-logo.png` | `anh-tham-chieu\`, `dot-3-toan-bo\anh-tham-chieu\`, `map-1-lang-la\anh-tham-chieu\` | 512×256 | Tham chiếu | — | Mẫu phong cách + bảng màu, gắn kèm mọi lần gen (manifest `anh_tham_chieu.TOT-naruto-logo`) |
| Chưa có | — | — | — | — | Icon app, splash, ảnh nền đăng nhập/chọn làng đã có 2 ảnh 1024×576 trong `ui_nen_dangnhap/chonlang` (kho) |

### 3.11 Ảnh tham chiếu & thước (không phải asset)

| File | Nơi | Cỡ | Dùng cho |
|---|---|---|---|
| `THUOC-NHAN-VAT.png` | `dot-4-ke-hoach-tong-the\anh-tham-chieu` | 2400×1060 | Thước nhân vật (idle_0: cao mặt 96–99, tóc→cằm 266–271, cằm→bàn chân 375–395 trong ô 724) — **bắt buộc** cho mọi brief nhân vật/NPC/quái/boss |
| `THUOC-TI-LE.png` | `map-1-lang-la\anh-tham-chieu` | 1096×638 | Thước map (1 ô = 96; người 2,4 ô; sói 1 ô; NPC 1,33 ô; nhà 6–8 ô) |
| `TOT-*.png` (9 + 2 THAM-CHIEU) | `dot-3-toan-bo\anh-tham-chieu` (16 file, bản đủ nhất) | 64…512 | Mẫu tốt: logo, khung, thanh chakra, chi tiết nhân vật/hiệu ứng, NPC Tazuna, sói, vật phẩm, boss mẫu |
| `XAU-*.png` (5) | 〃 | | Mẫu xấu: boss Zabuza phẳng, NPC 200×200, vật phẩm, hu_anthan quá mờ, thanh gộp |
| `HIENCO-*.png` (8), `HIEN-TRANG-map1-trong-game.png`, `GIAI-THICH-tileset-vien-chuyen-tiep.png`, `TOT-cong-torii-da-tach-nen.png` | `map-1-lang-la\anh-tham-chieu` | | Art đang có trong game để bám tông; lỗi cần tránh |
| `review\khe-dau-than-2026-09-05\*.png/.gif` (14) | `docs\my-thuat\review` | | Bằng chứng chẩn đoán khe đầu–thân; `so-sanh-ghep.png`, `zoom-co-x6.png` nên gắn vào brief nhân vật làm "lỗi phải tránh" |
| `vertical-slice-v1\qa\contact-sheet-light/dark.png` + 20 contact sheet | | 1000×950… | Kiểm nền chroma sáng/tối |

---

## 4. Lệch giữa asset-lock (đang ký) và client thật (đang chạy) — phát hiện khi kiểm kê

| # | Chỗ lệch | Lock nói gì | Client thật có gì | Hệ quả cho đợt 4 |
|---|---|---|---|---|
| 1 | **4 lớp nền map 1** | `build_inputs.background_sources`: b00 `map_nen_troi_1`, b01 `map_nen_may_1`, b02 `map_nen_may_2`, b03 `map_nen_nuihokage` (từ `anh-game`, 96×96 và 384×192) | `res\x4\bg`: b00 768×192, b01 1536×576 (= `nen_lang_xa`), b02 1536×768 (= `nen_nuihokage`), b03 768×768 (= `nen_troi`) — bản hand-painted 22–23/08 cắm qua Unity MCP, **không qua lock** | Bind lại nền theo bản thật; đồng thời vẽ bản 2976 rộng |
| 2 | **Logo** | `ui_logo_game_x4.png` 512×256 → `/Logo/1.png` (direct) | `x4\Logo\1.png` 599×136 wordmark | Chốt một trong hai |
| 3 | **UI overlay + HUD + Mob** | outputs có `client/x*/naruto/ui/*` (5 file), `client/x*/mainImage/naruto-*` (4), `client/x*/Mob/0/img.png`, `client/Mob/0.bytes` | **Không tồn tại** trong `D:\NRO_Project\Assets\Resources\res` (x1–x4 không có thư mục `Mob`, `naruto`; mainimage không có `naruto-*`) | Đây là artifact build/staging (Task 7–10) chưa chép vào client; sói/HUD hiện đến client bằng đường khác (máy chủ/cache) — **phải xác minh trước khi thêm quái/HUD mới** |
| 4 | **Mảnh người chơi** | `sprites[].zooms` trỏ `production-closure-v11/sprites/27041…/x4/sprite.png` (128×128, 128×160) | `smallimage` 05/09 là bản sinh lại (128×103, 108×84…) từ `lang-la-client-assets\res` | Closure mới cần reviewer ký lại (ghi trong `player_component_geometry.ghi_chu`) |
| 5 | **Tile map 1** | `tile_art` 100 ô từ closure v11 `tile-art/tiles/000–099` + `tile_set_overlay` khai 60–69, 83–86 | client x1/t có thêm 87, 88 và bộ 70–82 vật thể hand-painted | Kiểm `tile_set_info.bin` khớp ô 87/88 chưa |
| 6 | **ID lạ trong client** | id ký: 27041–27082; `forbidden_poc_ids` 9000–9002 | client có Small27083, 27084, 9000, 9001, 9002 | Dọn hoặc ghi sổ nguồn gốc |

---

## 5. Tên trong `kich-thuoc-anh.json` chưa có ảnh (336 / 691)

| Tiền tố | Tổng tên | Đã có ảnh | **Chưa có** | Chia theo làng (chưa có) |
|---|---:|---:|---:|---|
| `vp` | 260 | 100 | **160** | cat 40 · smu 40 · da 40 · may 40 |
| `map` | 143 | 79 | **64** | cat 14 · smu 14 · da 14 · may 14 · `map_nen_*` 8 (bien, dao, henui, nuida, nuimay, samac, vachda, vachnui) |
| `ui` | 62 | 62 (50 trong 343 + 12 trong `my thuat`) | 0 | — |
| `tp` | 60 | 36 | **24** | 4 làng × (tanbinh, chunin, jonin) × (dau, than) |
| `quai` | 44 | 12 | **32** | 8/làng: cat (bocap, giun_cat, golem_cat, kenken, lac_da, linh_gac, nhen, ran_luc…), smu, da, may |
| `npc` | 27 | 15 | **12** | 3/làng (giao_nv, thay, kage của làng) |
| `cd` | 27 | 15 | **12** | như npc |
| `kn` | 20 | 0 (nội dung có trong 5 sheet chroma) | **20** | 4/làng — chỉ thiếu xuất tên |
| `boss` | 18 | 6 | **12** | 3/làng (cat: bunhin, kiemphong, shukaku; da: golem_lon, tsuchikage, tuvi; may: batvi, raikage…; smu 3) |
| `hu` | 10 | 10 | 0 | — |
| `vk` | 20 | 20 | 0 | — |
| **Tổng** | **691** | **355** | **336** | |

> Bảng cỡ là **append-only** và các tên `_x4` này là **khoá tra cứu** của hệ cũ (ảnh tĩnh). Với nhân vật/NPC/quái/boss, đợt 4 chuyển sang **dải chroma** — nên các tên `quai_*`, `npc_*`, `boss_*`, `tp_*` chưa vẽ **không nên vẽ theo cỡ cũ** nữa; giữ tên làm mã tra cứu, đổi kiểu file.

---

## 6. Kết luận cho đợt 4 — dùng lại được gì, phải vẽ gì

**Dùng lại ngay (chỉ cần cấp id / bind, soi phong cách):**
- 100 icon vật phẩm `vp_*` 64×64 (đủ cho 5 mẫu vat-pham.yaml hiện tại và KH7).
- 15 chân dung `cd_*` (200×200 sau thu), 50 icon/khung UI, 20 vũ khí (icon), 8 UI đợt 1 đã trong lock.
- 16 icon chiêu 4 làng còn lại (đã có trong `skill_icons_{sand,mist,stone,cloud}_4_chroma`) — chỉ thiếu cấp id 20 chiêu.
- `hu_anthan_10f_chroma` + 10 `hu_*` sau khi cắt khung và gắn effect_id.
- Map 1: 51 ảnh hand-painted; cần thêm 2 nền rộng 2976.

**Dùng làm tham chiếu (không lắp thẳng):**
- 12 quái, 6 boss, 15 NPC đứng nguyên con, 36 mảnh trang phục, 79 tile pixel-art, `master_hoc_vien/soi_rung/iruka`.

**Phải vẽ lại / vẽ mới (theo bài học 05/09, thước THUOC-NHAN-VAT):**
- Người chơi: dải **run 4**, **melee 4**, **hit 3 (đứng)** — cùng cỡ đầu idle_0, chân chạm mép dưới, nền chroma đồng nhất. (idle 3, cast 3 giữ.)
- NPC: mọi NPC mới = `master_<tên>_x4_chroma` + `<tên>_portrait_chroma` (không vẽ 96×128 nữa). Iruka **không cần vẽ lại** — sửa dx/dy part.
- Quái mới: 4 dải như sói. Boss: 5 dải như người chơi (ghép 3 mảnh).
- Map 2–6: theo `KHUON-DE-BAI-ANH-MAP.md` (27–33 ảnh/map).

**Việc dọn trước khi bind thêm:** 6 chỗ lệch ở §4; xác minh đường đi ảnh sói/HUD vào client; dọn Small9000–9002/27083/27084.

---

## 7. Lệnh đã dùng để đếm (tái lập được)

Chạy bằng `python` (Python 3.14 trên máy này), chỉ đọc header PNG (8 byte chữ ký + IHDR) — không cần Pillow:

```text
# Đếm file theo thư mục 343-anh và đối chiếu manifest (tên, cỡ)
python -c "import json,io,os,struct;..."   → 343/343 khớp tên; 62 file to hơn manifest 4×
# Cỡ dải chroma: player 2172×724 (ô 724), wolf/hud/hu 1774×887 hoặc 2172×724, master/portrait 1254×1254, skill 1983×793
# asset-lock.json: inputs 684 · outputs 668 (client 663 + server 5) · closure sprites 27041–27082 · versions data 27/map 22/resource 21/runtime 22
# mapping-id.json: sprites 42 · parts 6 (30/57/58, 2189/2190/2191) · tiles 100 · maps 1 · backgrounds 1 (4 lớp) · mobs 1 · skills 4 · ui_bindings 8
# kich-thuoc-anh.json: 691 tên → chưa có ảnh 348 (trong đó 12 ui_ có ở "my thuat") → 336 thật sự chưa có
# Client (chỉ đọc): res\x1\{bg 8, mapBackGround 8, t 48, smallimage 47}; res\x4\Logo\1.png 599×136; không có Mob/, naruto/
```

Script sinh các bảng 343 trong file này: `scratchpad\sinh-bang-343.py` (đọc `manifest-tao-anh.json` + header PNG thật; không sửa kho ảnh).

*Người kiểm: agent I1, 05/09/2026. Chưa mở ảnh bằng mắt — mọi nhận xét về "phong cách pixel-art vs hand-painted" lấy từ `phong_cach_khoa` trong manifest và `DOT-2-ANH-CAN-VE.md`, cần chủ dự án nhìn xác nhận.*
