# 00 — KIỂM KÊ "GAME CẦN GÌ" — Đợt 4 (bản kiểm kê, chỉ đọc · 05/09/2026)

> Đây là **sổ kiểm kê nhu cầu ảnh**, không phải brief vẽ. Các brief đợt 4 (01-…, 02-…) sẽ tách từ đây, mỗi hạng mục một file MD + JSON. Đọc `_KHUON-DOT-4.md` và `../dot-3-toan-bo/_KHUON.md` trước khi viết brief.
>
> **Cách đọc:** mọi dòng "cần" đều có cột **Nguồn** trỏ `file:dòng` — đó là chỗ máy chủ hoặc client *quyết định* game cần ảnh đó. Dòng nào máy chủ **chưa chốt nội dung** thì ghi rõ **"đề xuất, chờ chốt"** — không được coi là đã chốt.
>
> **Ưu tiên:** **P1** = để Làng Lá + Khu rừng chơi trọn với máy chủ *đã có hôm nay* · **P2** = 4 làng còn lại theo lộ trình · **P3** = đẹp thêm / thương hiệu.

---

## 0. Tóm tắt một phút

| Ưu tiên | Cần vẽ MỚI (ước) | Vẽ LẠI | Đã có, chỉ cần lắp | Đề xuất, chờ máy chủ chốt |
|---|---|---|---|---|
| **P1** Làng Lá + Khu rừng | ~45 file: 1 chân dung người chơi · 7 dải boss (1 bộ 5 dải + 2 dải đầu) · 6 ảnh map 1 · ~25 ảnh map Khu rừng · ~6 dải hiệu ứng chiêu | 5 dải người chơi (chờ chủ dự án quyết V1 dữ liệu / V3-B vẽ lại) | Sói 4 dải · Iruka master + chân dung · 4 icon chiêu Lá · HUD máu/chakra · 5 icon vật phẩm (có trong 343, chưa gắn id) | Quái Khu rừng 2–3 con (8–12 dải) · NPC thêm (KH9 chưa bắt đầu) |
| **P2** 4 làng còn lại | 4 concept map · ~120 ảnh map (4 × ~30) · 4 bộ người chơi (20 dải) + 4 chân dung · 1 bản đồ thế giới | — | 16 icon chiêu (4 làng) đã vẽ, chưa vào lock | Quái ≥ 8 con (32 dải) · Boss người 4 bộ (20 dải) · NPC 4–12 master + chân dung · Hiệu ứng 4–16 dải · Icon vật phẩm 8–160 |
| **P3** đẹp / thương hiệu | ~45–50 ảnh HUD `/mainImage` · 2–3 nền bảng · 2 nền màn hình · 12 bộ trang phục (60 dải) · 20 vũ khí | — | 50 UI + 36 mảnh trang phục + 20 vũ khí trong 343 (chỉ 8 UI đã bind) | Boss vĩ thú (cần code mới) · hiệu ứng thời tiết theo làng |

**Ba phát hiện quan trọng nhất của lần kiểm kê này:**

1. **Ảnh tĩnh 96×96 / 160×160 của quái và boss trong `kich-thuoc-anh.json` KHÔNG lắp được vào game.** Client vẽ quái bằng **sprite sheet nhiều khung** (`Mob.cs:383-393`), vẽ boss bằng **bộ xương nhân vật 3 mảnh × 33 khung** như người chơi (`boss.yaml:18-31`, KH6b §3). 44 tên `quai_*` và 18 tên `boss_*` trong bảng cỡ chỉ còn giá trị **tham chiếu/chân dung**.
2. **Boss Làng Lá hôm nay "vô hình" về mỹ thuật:** pha 1 mặc đúng bộ người chơi (`30/57/58`), pha 2 mặc đúng bộ Iruka (`2189/2190/2191`) — `boss.yaml:28-31` tự ghi là *nợ khách*. Đây là nhu cầu P1 to nhất mà chưa brief nào nói tới.
3. **Hiệu ứng chiêu chưa có đường lắp rõ:** client nạp `/effectdata/{id}/data` + `img.png` (`Effect.cs:156-167`) và máy chủ nhúng `skill.bin` 70 byte (`dongbo.go:21`) — chưa ai giải mã pipeline này sang ảnh Naruto. Ghi là **chưa chắc**, phải giải mã trước khi giao vẽ.

---

## 1. Client dùng ảnh thế nào — 7 kiểu, quyết định kiểu file phải giao

Đây là bảng gốc để mọi brief đợt 4 chọn đúng **kiểu file**. Sai kiểu = vẽ đẹp mà không lắp được (bài học `quai_*`/`boss_*` ở trên).

| # | Đối tượng | Client vẽ thế nào (bằng chứng) | Kiểu file PHẢI giao | Cỡ / ô / số khung |
|---|---|---|---|---|
| K1 | **Người chơi** | Ghép 3 mảnh đầu/thân/chân theo bộ xương cứng 33 khung `Char.CharInfo` (`ChooseCharScr.cs:102-104`, `CreateCharScr.cs:410-412`, `Npc.cs:305-313` cùng công thức); part có 3 đầu / 17 thân / 14 chân (`asset-lock.json` `closure.frame_counts` 30:3 · 57:17 · 58:14). Pipeline `tools/artqa/task6b_closure.py` cắt từ **5 dải hành động** | **5 dải chroma**: `idle_3f` · `run_4f` · `melee_4f` · `cast_3f` · `hit_death_getup_3f` = 17 khung → 34 mảnh × 4 zoom + `part.bin` | Ô **724×724**, nền tím RGB(246,4,242). Dải 3 khung = 2172×724; dải 4 khung **phải 2896×724** (bản v1 để 2172 → ô 543, vi phạm luật ô đều — xem §3) |
| K2 | **Boss** | Xuất hiện bằng gói `-5` **như một nhân vật**, mang `outfit[6]` = id part đầu/thân/chân (`docs/superpowers/specs/2026-08-03-kh6b-boss-design.md` §3 dòng 28-30; `boss.yaml:18-31`). `outfit` bị kiểm phải nằm trong `closure.part_ids` của lock | **Giống K1**: 1 bộ = 5 dải chroma → 34 mảnh. Boss dạng **thú** (vĩ thú, golem to) **không vẽ được** bằng bộ xương người → cần code mới (boss dạng Mob) — ngoài phạm vi | Như K1 |
| K3 | **NPC** | Cũng ghép 3 mảnh nhưng chỉ khung đứng cf 0/1 (`Npc.cs:305-313`); assetpack cắt đầu/thân/chân từ **một ảnh master đứng** (spec vertical-slice §6.2 dòng 295-299). Chân dung hội thoại = một Small image id (`ChatPopup.cs:316` `c.avatar`; `catalog.json` `npcs[0].portrait: 27078`) | **1 master đứng chroma** (mẫu `master_iruka_x4_chroma.png` 1254×1254) + **1 chân dung chroma** (mẫu `iruka_portrait_chroma.png` 1254×1254) | Master: đứng thẳng, nhìn 3/4 quay phải, chân chạm mép dưới. Chân dung: cận mặt vuông |
| K4 | **Quái** | Mỗi mẫu quái một thư mục `Mob/<templateId>/data` + `img.png` đọc bằng `EffectData.readData` (`Mob.cs:383-393`); sói hiện là sheet 10 ô 96×96: idle 2 · run 4 · pounce 2 · hit 1 · death 1 (`asset-lock.json` `closure.mob.actions`). `anh_id` trong `quai.yaml` phải có trong `closure.mob_ids` (`quai.yaml:20-31`) | **4 dải chroma**: `idle_2f` · `run_4f` · `pounce_2f` · `hit_death_2f` (mẫu `wolf_*_chroma.png`) | Dải 2 khung 1774×887 (ô 887) · dải 4 khung 2172×724 (ô 543) trong v1 — **đợt 4 chốt một cỡ ô chung cho cả 4 dải** |
| K5 | **Map** | Tile 96×96 (x4) theo id `res/x{zoom}/t/{tileID}$V.png`, tối đa 100 ô/tileset (`TileMap.cs:287-332`; spec vertical-slice §7); nền 4 lớp `res/x{zoom}/bg/b{typeBg}{k}.png` (`GameCanvas.cs:1694`); vật thể lớn xếp qua bố cục `mymap/{id}.bytes` (`TileMap.cs:529`); ảnh nền phải rộng **2976** vì game hiển thị x2 (`KHUON-DE-BAI-ANH-MAP.md` §Vì sao 2976) | **Ảnh tĩnh alpha** theo khuôn map: 4 nền (đặc) 2976 rộng + 12–14 tile 96×96 + công trình/cây/vật nhỏ bội số 96 | Xem `KHUON-DE-BAI-ANH-MAP.md` bảng cỡ; 27–33 ảnh/map |
| K6 | **Icon vật phẩm / chiêu / UI nhỏ** | Vẽ qua `SmallImage.drawSmallImage(g, id…)`: icon vật phẩm = `item.template.iconID` (`Panel.cs:2807`, `:6301`, `:7158`), icon chiêu = `skill.template.iconId` (`Panel.cs:7006`, `:7519`), chân dung HUD = `Char.idAvatar` (`GameScr.cs:5588-5590`, máy chủ gửi cặp head→avatar trong ME_INFO `thegioi.go:156-168`). Small image được **đăng ký id trong lock** (`closure.sprite_ids` 27041–27082, `skill_ids` 27079–27082) | **Ảnh tĩnh alpha** 64×64 (vp/vk/kn) hoặc 200×200 (cd), giao theo tên `kich-thuoc-anh.json`; dự án gắn id + ký lock | Theo bảng §2 `_KHUON.md` |
| K7 | **Hiệu ứng** | `Effect.cs:156-167`: nạp `/x{zoom}/effectdata/{id}/data` (khung) + `/effectdata/{id}/img.png`; chiêu gắn hiệu ứng qua `SkillPaint` (`Char.cs:940`, `:2068-2079`); máy chủ nhúng `skill.bin` (`dongbo.go:21`) và `catalog.json` `skill_poses [0,1,2,3]` | **Dải chroma ≥ 5 khung**, ô rộng chia hết cho 4 (mẫu `hu_anthan_10f_chroma.png` 1774×887) | ⚠️ **Chưa chắc**: chưa giải mã cách đóng `effectdata` mới; đợt 4 chỉ vẽ sau khi có script đóng gói |

> Ba màn hình khác không sinh nhu cầu ảnh mới ở P1: `LoginScr.cs` chỉ vẽ logo (`ModFunc.imgLogoBig`, dòng 655-667 — logo đã thay `Logo/1.png`); `ServerListScreen.cs:306-324` dùng logo + `imgRoomStat`; `InfoDlg.cs` không nạp ảnh riêng. `CreateCharScr.cs:127-188` vẽ **nền màn tạo nhân vật bằng một map** (`/mymap/…` + `/mapBackGround/{id}.png`) → dùng lại map Làng Lá, không cần ảnh riêng.

---

## 2. Máy chủ đã chốt gì — nguồn sự thật (đọc trước khi đếm)

| Nội dung | Đã chốt | Nguồn |
|---|---|---|
| 5 làng | id 0 Lá (Hoả) · 1 Cát (Phong) · 2 Sương Mù (Thuỷ) · 3 Đá (Thổ) · 4 Mây (Lôi); **không được thêm làng 6** | `server/kho-noi-dung/lang.yaml:27-88`, `:10-11` |
| Bản đồ | catalog 6 map: 1 Làng Lá · 2 Cát · 3 Sương Mù · 4 Đá · 5 Mây · 6 Khu rừng | `server/internal/game/assets/catalog.json:1` `maps[]` |
| Bản đồ **thật sự vào được** | `banDo.go` chỉ khai **5** dòng (map 1–5); map 1 dùng tile Naruto (TileID 1, BgID 0), **map 3–5 còn dùng tile Dragon Ball** (TileID 2/3, BgID 1/2); cổng chỉ nối 1 ↔ 2 | `server/internal/game/banDo.go:81-89`, `:115-135`; `BAT-DAU-TU-DAY.md` §5 dòng 1a′ |
| Chiêu | 20 chiêu = 5 làng × 4 slot (cận chiến · tầm xa theo hệ · hồi máu · tụ chakra), `icon_id` 1–20; **đã chốt giữ 4 chiêu/làng (03/09)** | `chieu.yaml:60-951`; `LO-TRINH.md:70` |
| Quái | **1 mẫu duy nhất** `soi-rung` (`anh_id 0`), sinh 2 con ở **cả 6 map** | `quai.yaml:67-74`, `:93-121` |
| Boss | 2 boss, đều ở **map 1**: `nhan-gia-mat-na` 2 pha (outfit học viên → outfit jounin/Iruka) · `hoc-vien-dao-tau` 1 pha (làng 2, outfit học viên) | `boss.yaml:91-251`; nợ mỹ thuật ghi ở `:28-31` |
| NPC | **1** NPC Iruka (part 2189/2190/2191, portrait 27078, menu Cửa hàng) | `catalog.json:1` `npcs[0]`; `BAT-DAU-TU-DAY.md` §5 mục (7) |
| Vật phẩm | **5 mẫu**: thuốc hồi máu nhỏ · thuốc hồi chakra nhỏ · áo ninja Làng Lá · quần ninja Làng Lá · kunai luyện tập; **mọi `icon_id` = 0, bảng SmallImage client rỗng** | `vat-pham.yaml:90-149`, `:78-81`; `khoa-id.yaml:77-88` |
| Cửa hàng / bảng rơi | 1 cửa hàng "Tạp hoá Làng Lá" bán 5 món; quái rơi thuốc + kunai; boss lớn rơi áo/quần | `cua-hang.yaml:40-79`; `bang-roi.yaml:54-101` |
| Người chơi | Mọi làng hiện dùng **một** bộ học viên 30/57/58 (`cgender` luôn 0); chân dung HUD 27084 là **đầu cắt từ sprite**, không phải chân dung vẽ | `hinhdang.go:49-57`; `thegioi.go:156-168` |
| KH9 (NPC, nhiệm vụ) · KH10 (cảm giác chơi) | **Chưa có spec** — chỉ có phác thảo | `LO-TRINH.md:61-62`, `:102-155`; không có file `kh9-*`/`kh10-*` trong `docs/superpowers/specs/` |
| KH7 mở rộng vật phẩm | Bản 1 xong 5 mẫu; danh mục 4 làng **chưa có** | `docs/superpowers/specs/2026-09-03-kh7-vat-pham-design.md` §3; `BAT-DAU-TU-DAY.md` §5 mục 2 |

---

## 3. Đã có gì trong kho ảnh — và dùng được không

### 3.1 Lô `vertical-slice-v1` (dải chroma, đúng kiểu client cần)

| File | Cỡ đo thật | Dùng được? | Ghi chú |
|---|---|---|---|
| `player_idle_3f` · `cast_3f` · `hit_death_getup_3f` | 2172×724 (ô 724) | ⚠️ đang lắp, **lệch tỉ lệ** | Bài học 05/09: 5 dải vẽ nhân vật khác cỡ (chạy nhỏ hơn 12%, đánh 17%), tư thế khom → khe đầu–thân. Chờ chủ dự án quyết V1 (sửa dữ liệu) hay V3-B (vẽ lại) — `BAT-DAU-TU-DAY.md` §5 dòng "Khe đầu–thân" |
| `player_run_4f` · `player_melee_4f` | **2172×724** (ô chỉ 543 rộng) | ⚠️ như trên | Ô không đều với dải 3 khung → vi phạm luật (5). Nếu vẽ lại: 2896×724 |
| `master_hoc_vien_x4_chroma` | 1254×1254 | ✅ tham chiếu đồng bộ | — |
| `master_iruka_x4_chroma` + `iruka_portrait_chroma` | 1254×1254 | ✅ đang lắp | Lỗi hiển thị (đầu lửng +13 px, thân chìm +29 px) là **dx/dy part = 0** — **đã sửa bằng dữ liệu 06/09** (`75a64cc`), không vẽ lại |
| `master_soi_rung_x4_chroma` + `wolf_{idle_2f, run_4f, pounce_2f, hit_death_2f}` | 1774×887 / 2172×724 | ✅ đang lắp (mob sheet 10 ô) | Ô hai cỡ khác nhau giữa các dải — chấp nhận vì đã đóng gói xong |
| `hud_hp_chakra_4_chroma` | 1774×887 | ✅ đang lắp | 4 lớp nền/fill máu + chakra (spec §6.3) |
| `skill_icons_leaf_4` | 1983×793 | ✅ đã vào lock (`skill_ids` 27079–27082) | — |
| `skill_icons_{sand, mist, stone, cloud}_4` | 1983×793 | ✅ đã vẽ, **chưa vào lock** | Việc lắp (P2), không vẽ |
| `hu_anthan_10f_chroma` | 1774×887 | ⚠️ đã vẽ, chưa vào runtime | Là hiệu ứng ẩn thân — không ứng với chiêu nào trong `chieu.yaml` |

### 3.2 Lô 343 ảnh tĩnh alpha (`anh-nguon/343-anh/`)

| Nhóm | Số | Dùng được theo kiểu client (§1)? |
|---|---|---|
| 01 giao diện `ui_*` | 50 | ✅ **K6** đúng kiểu. 8 UI đã bind (`asset-lock.json` `runtime_bindings`: logo, khung avatar, khung hội thoại, 3 nút, 2 ô); 42 còn lại chưa có chỗ lắp (P3) |
| 02 vật phẩm `vp_*` | 100 | ✅ **K6** đúng kiểu, cần QC + gắn id. Có sẵn cho 5 mẫu P1: `vp_ao_t1`, `vp_quan_t1`, `vp_vien_theluc_1` (thuốc máu), `vp_vien_chakra_1` (thuốc chakra) + `vk_kunai_thuong` (04) |
| 03 trang phục `tp_*` | 36 (12 bộ × 3 mảnh tĩnh) | ❌ **không phải** 34 khung K1 — chỉ là tham chiếu phong cách (spec vertical-slice §6.2 dòng 294) |
| 04 vũ khí `vk_*` | 20 | ✅ K6 làm icon túi; client không có ô vẽ vũ khí lên người → P3 |
| 05 quái `quai_*` | 12 (bảng cỡ có 44 tên) | ❌ **sai kiểu** — client cần sheet K4. Dùng làm silhouette tham chiếu (như `quai_soi_rung` đã làm) |
| 06 boss `boss_*` | 6 (bảng cỡ 18 tên) | ❌ **sai kiểu** — boss là nhân vật K2. Chỉ dùng làm tham chiếu/chân dung |
| 07 NPC `npc_*` 96×128 + `cd_*` 200×200 | 15 + 15 | ⚠️ `npc_*` cắt 3 mảnh được như Iruka nhưng **đã lộ lỗi dx/dy**; đợt 4 nên vẽ master chroma theo mẫu `master_iruka`. `cd_*` ✅ K6 |
| 08 hiệu ứng `hu_*` 1200×480 (5 khung 240) | 10 | ⚠️ khuôn cũ đợt 1, chưa vào runtime; cần giải mã K7 trước |
| 09 bản đồ `map_*` | 79 | ✅ tile 96×96 đúng kiểu (30 tile map 1 đã lắp `t/1$V`); **`map_nen_*` 384×192 SAI CỠ** (phải 2976 rộng) — không dùng cho nền lớn nữa |

### 3.3 Đã lắp thật trong client (`lang-la-client-assets/`, kiểm 05/09)

`res/x{1..4}/bg/b00–b03` (4 nền map 1) · `res/x{1..4}/t/1$V.png` (64 tile x2) · `res/x{1..4}/smallimage/` (34 mảnh người chơi) · `res/mymap/1.bytes` · `Logo/1.png`. Sói/Iruka/HUD đi qua `assetpack` (lock `versions.data 27`). **Không có** ảnh `/mainImage/*`, `/effectdata/*`, icon vật phẩm nào là Naruto.

---

## 4. P1 — Làng Lá + Khu rừng chơi trọn (máy chủ đã có)

| STT | Đối tượng | Kiểu (§1) | Số khung / cỡ | Làng | Số lượng | Trạng thái | Nguồn quyết định |
|---|---|---|---|---|---|---|---|
| P1-01 | **Người chơi học viên** — 5 dải hành động | K1 dải | idle 3 · run 4 · melee 4 · cast 3 · hit 3 = 17 khung; ô 724×724 | Lá (dùng chung 5 làng hôm nay) | 5 dải | **ve_lai — CHỜ QUYẾT** (V1 sửa dữ liệu, nhân vật thấp đi ~10 % khi chạy · V3-B vẽ lại 5 dải cùng cỡ đầu theo `THUOC-NHAN-VAT.png`) | `Char.cs` CharInfo; `asset-lock.json` `closure.frame_counts`; `BAT-DAU-TU-DAY.md` §5 dòng "Khe đầu–thân"; `hinhdang.go:49-57` |
| P1-02 | **Chân dung người chơi** (ô HUD + chat) | K6 tĩnh / chroma | 200×200 (cd) — vẽ lớn 1254×1254 như `iruka_portrait` | Lá | 1 | **can_ve** — hiện 27084 là đầu cắt từ sprite (`tools/map-lang-la/tao-chan-dung-nguoi-choi.py`), không phải chân dung | `GameScr.cs:5588-5590`; `ChatPopup.cs:316`; `thegioi.go:161-167` |
| P1-03 | **NPC Iruka** master + chân dung | K3 | — | Lá | 0 | **da_co** — chỉ sửa dx/dy part (dữ liệu) | `catalog.json:1` `npcs[0]`; `BAT-DAU-TU-DAY.md` §5 "NPC Iruka hỏng nặng hơn" |
| P1-04 | **Quái Sói Rừng** 4 dải | K4 | 10 khung | dùng chung 6 map | 0 | **da_co** (mob sheet trong lock) | `quai.yaml:67-74`; `asset-lock.json` `closure.mob` |
| P1-05 | **Boss "Kẻ Lạ Mặt"** (Nhân Gia Mặt Nạ pha 1) — đầu mặt nạ riêng | K2 (chỉ mảnh ĐẦU) | 1 dải đầu 3 khung idle (dùng lại thân/chân học viên 57/58) | Lá | 1 dải | **can_ve** — hiện outfit `[30,57,58]` y hệt người chơi | `boss.yaml:113-119`; `boss.yaml:18-31` (nợ khách) |
| P1-06 | **Boss "Mặt Nạ Vỡ"** (pha 2) — bộ nhân vật riêng | K2 | 5 dải × ô 724 (17 khung) → 34 mảnh, thêm 3 part id vào lock | Lá | 5 dải | **can_ve** — hiện outfit `[2189,2190,2191]` y hệt Iruka | `boss.yaml:167-171`; KH6b design §3 dòng 28-30 |
| P1-07 | **Boss "Học Viên Đào Tẩu"** — đầu riêng (băng trán Sương Mù 4 vạch sóng) | K2 (chỉ mảnh ĐẦU) | 1 dải đầu 3 khung | Sương Mù (đứng ở map 1) | 1 dải | **can_ve** | `boss.yaml:212-225` |
| P1-08 | **Chân dung 3 boss** cho khung thoại | K6 | 200×200 | Lá/Sương Mù | 0–3 | **chưa chắc** — chưa xác minh client vẽ avatar khi boss thoại qua gói 44 (`ChatPopup.cs:316` dùng `c.avatar`); kiểm trước khi vẽ | `boss.yaml` `thoai_*`; `ChatPopup.cs:316` |
| P1-09 | **Map 1 Làng Lá — phần còn thiếu** | K5 tĩnh | 4 tile đất có vân 96×96 (`tile_dat_van_1..4`) + `nen_nuihokage_rong2` **2976×768** + `nen_lang_xa_rong2` **2976×576** | Lá | 6 | **can_ve** (đất phẳng lì 91 %, đích ≤ 40 %; núi lặp 1,9 lần/khung) | `map-1-lang-la/DOT-2-ANH-CAN-VE.md` §Đợt 3 A/B; `BAT-DAU-TU-DAY.md` §5 dòng 1a; `2026-08-23-lo-trinh-4-map-con-lai.md` §Việc còn lại |
| P1-10 | **Map 6 Khu rừng** — bộ ảnh theo khuôn map | K5 tĩnh | 4 nền 2976 (trời 384 · núi/rừng xa 768 · tán rừng 576 · gần 288) + 12 tile nền + 2 tile bậc + 3 cây (768/384/288×192) + ~4 vật nhỏ (gốc cây, tảng đá, cầu gỗ, biển) — **không có nhà** | Khu rừng | ~25 | **can_ve** — 25 tile `map_rung_*` trong 343 chỉ là tile 96, chưa có nền/cây theo khuôn. ⚠️ **Chưa chắc phần server:** `banDo.go:81-89` chỉ khai map 1–5, map 6 chưa có `BanDo`/cổng → lắp ảnh chỉ có nghĩa sau khi server nối | `catalog.json:1` `maps[5]`; `quai.yaml:118-121`; `KHUON-DE-BAI-ANH-MAP.md` |
| P1-11 | **Quái cho Khu rừng** (ngoài sói) | K4 | mỗi con 4 dải (10 khung) | Khu rừng | 2–3 con = 8–12 dải | **đề xuất, chờ chốt** — `quai.yaml` chỉ có 1 mẫu; đề xuất lợn rừng / nhện độc / khỉ rừng (đã có silhouette `quai_lon_rung`, `quai_nhen_doc`, `quai_khi_rung` trong 343) | `quai.yaml:67-74`, `:80-88` ("thiết kế nội dung thật sẽ cho mỗi làng một đàn quái khác nhau") |
| P1-12 | **Icon 5 vật phẩm** đang bán/rơi | K6 tĩnh 64×64 | thuốc máu · thuốc chakra · áo Lá · quần Lá · kunai | Lá / chung | 5 | **da_co (cần QC + gắn id + ký lock)** — 343 có `vp_vien_theluc_1`, `vp_vien_chakra_1`, `vp_ao_t1`, `vp_quan_t1`, `vk_kunai_thuong`; vẽ lại nếu không qua `KIEM-TRA-ANH.py` | `vat-pham.yaml:90-149`, `:78-81`; `Panel.cs:2807`, `:6301` |
| P1-13 | **Icon 4 chiêu Làng Lá** | K6 | 64×64 | Lá | 0 | **da_co** (lock `skill_ids` 27079–27082) | `chieu.yaml:65-238`; `asset-lock.json` `closure.skill_ids` |
| P1-14 | **Hiệu ứng 4 chiêu Làng Lá** + trúng đòn + hồi sinh | K7 dải | mỗi dải ≥ 5 khung, ô rộng chia hết 4: Cước Liên Hoàn (cận chiến) · Hào Hoả Cầu (bay + nổ) · Y Liệu (vòng sáng hồi máu) · Tụ Chakra (khí bốc) · trúng đòn (chớp) · chết/hồi sinh (vòng sáng dưới chân — client không vẽ nằm) | Lá | ~6 dải | **can_ve — CHƯA CHẮC đường lắp**: phải giải mã `EffectData`/`SkillPaint` + `skill.bin` trước; 343 có 10 `hu_*` khuôn cũ | `Effect.cs:156-167`; `Char.cs:940`, `:2068-2079`; `dongbo.go:21`; `catalog.json:1` `skill_poses`; `LO-TRINH.md:112-118` (KH10 A) |
| P1-15 | **HUD máu/chakra** 4 lớp | K6 | 240×32 tách nền/fill | chung | 0 | **da_co** | spec vertical-slice §6.3 |
| P1-16 | **NPC thêm cho Làng Lá** (nhiệm vụ, y tá, thợ rèn…) | K3 | master + chân dung mỗi NPC | Lá | ? | **đề xuất, chờ KH9** — máy chủ chỉ có Iruka; 343 có 15 NPC tĩnh làm tham chiếu | `LO-TRINH.md:61`; `catalog.json:1` `npcs` |

**Tổng P1 cần vẽ mới:** 1 + 7 dải + 6 + ~25 + ~6 dải ≈ **45 file**; **vẽ lại (chờ quyết):** 5 dải; **đề xuất chờ chốt:** 8–12 dải quái + NPC KH9.

---

## 5. P2 — 4 làng còn lại (Cát · Sương Mù · Đá · Mây) theo lộ trình

Thứ tự đề nghị theo `2026-08-23-lo-trinh-4-map-con-lai.md`: **Cát → Sương Mù → Đá → Mây** (Cát khác Làng Lá nhất, thử khuôn được).

| STT | Đối tượng | Kiểu | Số khung / cỡ | Làng | Số lượng | Trạng thái | Nguồn quyết định |
|---|---|---|---|---|---|---|---|
| P2-01 | **Người chơi mặc định mỗi làng** | K1 | 5 dải × 17 khung mỗi bộ | Cát · SMù · Đá · Mây | 4 bộ = 20 dải (**rẻ hơn:** giữ thân/chân học viên, chỉ vẽ 4 dải ĐẦU 3 khung — băng trán/khăn theo làng — nhưng server phải gửi đầu theo làng, `hinhdang.go:57` hiện `_ = lang`) | **can_ve** | `hinhdang.go:15,49-57`; `lang.yaml:27-88`; đợt 3 `lang-*/nhan-vat.md` (3 bộ tĩnh/làng — sai kiểu) |
| P2-02 | **Chân dung người chơi 4 làng** | K6 | 200×200 | 4 làng | 4 | **can_ve** | `GameScr.cs:5588-5590` |
| P2-03 | **Concept 1 bức mỗi map** (đo 7 chỉ số trước khi vẽ bộ) | tham chiếu | tự do, ~1454 rộng | 4 làng | 4 | **can_ve trước** | `lo-trinh-4-map` bước 1–2; `SO-TAY-LAM-MAP.md` §2 |
| P2-04 | **Map 2–5** — bộ ảnh theo khuôn | K5 | mỗi map: 4 nền 2976 + 14 tile 96 + 5 công trình (576–768) + 3 cây + 5 vật nhỏ ≈ 30 | 4 làng | ~120 | **can_ve** — brief đợt 3 `lang-*/map.md` (16 ảnh, nền 384×192) là **khuôn cũ sai cỡ**, thay bằng `KHUON-DE-BAI-ANH-MAP.md`. Nét riêng + bẫy từng làng: bảng ở `lo-trinh-4-map` §Bốn map | `banDo.go:86-88` (tile Dragon Ball); `boss.yaml:52-53`; `BAT-DAU-TU-DAY.md` §5 1a′ |
| P2-05 | **Quái 4 làng** | K4 | 4 dải/con (10 khung) | 4 làng | tối thiểu 2 con/làng (1 thú + 1 lính gác) = 8 con = **32 dải**; đợt 3 đề xuất 8 con/làng (32 con) — quá sức | **đề xuất, chờ chốt** — `quai.yaml` chưa có quái nào ngoài sói; danh sách đề xuất ở `dot-3-toan-bo/lang-*/quai.md` | `quai.yaml:67-74`, `:80-88` |
| P2-06 | **Boss 4 làng** (dạng người) | K2 | 5 dải/bộ | 4 làng | 1 boss người/làng = 4 bộ = **20 dải** (+ chân dung nếu P1-08 xác nhận) | **đề xuất, chờ chốt** — `boss.yaml` chỉ có 2 boss Làng Lá. Đợt 3 đề xuất Kage/kiếm sĩ (vẽ được) + **vĩ thú Shukaku/Isobu/Son Gokū/Ngưu Quỷ (KHÔNG vẽ được bằng bộ xương người → cần code boss dạng Mob → P3/hoãn)** | `boss.yaml:91-251`; KH6b design §3; đợt 3 `lang-*/boss.md` |
| P2-07 | **NPC 4 làng** | K3 | master + chân dung | 4 làng | tối thiểu 1 NPC cửa hàng/làng = 4 + 4; đợt 3 đề xuất 3/làng (Kage · thầy · giao NV) = 12 + 12 | **đề xuất, chờ KH9** | `LO-TRINH.md:61`; đợt 3 `lang-*/npc.md` |
| P2-08 | **Icon 16 chiêu 4 làng** | K6 | 64×64 | 4 làng | 0 | **da_co, chưa vào lock** (`skill_icons_{sand,mist,stone,cloud}`) — việc lắp | `chieu.yaml:243-951` `icon_id 5–20`; `asset-lock.json` `closure.skill_ids` |
| P2-09 | **Hiệu ứng chiêu 4 làng** | K7 | dải ≥ 5 khung | 4 làng | tối thiểu 4 dải tầm xa (Cuồng Phong Trảm · Đại Thuỷ Long Đạn · Nham Thạch Đạn · Cầu Sét) + đổi màu cận chiến/hồi máu/tụ chakra = 4–16 dải | **chờ giải mã K7** (P1-14) | `chieu.yaml` slot 1 mỗi làng |
| P2-10 | **Icon vật phẩm 4 làng** | K6 | 64×64 | 4 làng | tối thiểu áo + quần làng mình = 8; đợt 3 đề xuất 40/làng (5 bộ phận × 8 cấp) = 160 | **đề xuất, chờ KH7 mở rộng** — danh mục chưa có, `khoa-id.yaml` chỉ thêm cuối | `vat-pham.yaml:14-15`; KH7 spec §3; đợt 3 `lang-*/trang-bi.md` |
| P2-11 | **Bản đồ thế giới / minimap** `/img/map{planetID}.png` | tĩnh | theo ảnh cũ (kiểm cỡ trong client) | chung (planetID luôn 0) | 1 | **can_ve (nhỏ)** | `Panel.cs:959`, `:5587`; `banDo.go:84-88` `PlanetID: 0` |

---

## 6. P3 — đẹp thêm / thương hiệu (game chạy được mà không cần)

| STT | Đối tượng | Kiểu | Số lượng | Ghi chú | Nguồn |
|---|---|---|---|---|---|
| P3-01 | **HUD/giao diện `/mainImage/*`** còn Dragon Ball (nút, khung, thanh, mũi tên, analog, chat, focus, panel, capsule…) | tĩnh theo **tên file client** | ~45–50 | Cần bảng ánh xạ tên client ↔ tên `ui_*`; 42 UI trong 343 chưa có chỗ lắp vì tên khác | `GameScr.cs:866-968`; `Panel.cs:952-974` |
| P3-02 | **Nền bảng túi đồ / panel** (`inventd.png`, `myTexture2dpanel`, `panel2`) | tĩnh | 2–3 | Theo gender → giờ 1 bản Naruto | `Panel.cs:876-877`, `:995-1026` |
| P3-03 | **Nền màn đăng nhập / chọn làng** | tĩnh | 0–2 | `ui_nen_dangnhap` 1024×576, `ui_nen_chonlang` đã có trong 343; client hiện chỉ vẽ logo + nền GameCanvas | `LoginScr.cs:655-667` |
| P3-04 | **Splash `/e/sp{i}.png`** | tĩnh | ? | Chưa đếm số khung | `GameScr.cs:5950` |
| P3-05 | **Hiệu ứng nền thời tiết** (`/bg/mua*`, `lacay*`, `fog0`, `sao`, `tuyet`) | tĩnh nhỏ | 0–15 | Dùng lại được ngay; bản Naruto theo làng (lá rơi Lá · cát bay Cát · sương SMù) là đẹp thêm | `BackgroudEffect.cs:143-404` |
| P3-06 | **Trang phục đổi hình khi mặc** (12 bộ `tp_*`) | K1 | 12 bộ × 5 dải = 60 dải | Chỉ có ý nghĩa khi KH7 bản 2 làm `Send_Caitrang` (hoãn) | `BAT-DAU-TU-DAY.md` §5 mục 2 (Task 7 hoãn); 343 `tp_*` |
| P3-07 | **Vũ khí `vk_*` 20** | K6 | 0 (đã có) | Chưa có ô vũ khí trên người; chờ danh mục KH7 | `vat-pham.yaml:139-149` (kunai là `loai 2`) |
| P3-08 | **Icon bậc `ui_bac_*` 9 · hệ `ui_he_*` 5 · tiền `ui_tien_*` 4 · làng `ui_lang_*` 5** | K6 | 0 (đã có) | Chưa có chỗ dùng trong client hiện tại | 343 `01-giao-dien` |
| P3-09 | **Boss vĩ thú 4 làng** | Mob sheet lớn | 4 | Cần code boss dạng Mob trước (KH6b ngoài phạm vi) | KH6b design §8 |
| P3-10 | **Ảnh `/u/{id}` cho bảng hiệu ứng nâng cấp** | tĩnh | ? | Tính năng cũ, chưa rõ có giữ | `Panel.cs:9457-9461` |

---

## 7. Những chỗ CHƯA CHỐT / CHƯA CHẮC — đừng brief như đã chốt

| # | Điều chưa chắc | Ai chốt / việc gì phải làm trước |
|---|---|---|
| C1 | Người chơi: **sửa dữ liệu (V1)** hay **vẽ lại 5 dải (V3-B)** | Chủ dự án — `BAT-DAU-TU-DAY.md` §5 "Khe đầu–thân". Brief đợt 4 nên viết sẵn bản vẽ lại theo `THUOC-NHAN-VAT.png` để sẵn dùng |
| C2 | Danh sách **quái** 4 làng + Khu rừng | Máy chủ (`quai.yaml`) — mọi con trong đợt 3 là đề xuất |
| C3 | Danh sách **boss** 4 làng; boss vĩ thú cần code mới | Máy chủ (`boss.yaml`) + KH6b bản sau |
| C4 | **NPC** ngoài Iruka, nhiệm vụ, thoại | KH9 chưa có spec |
| C5 | **Danh mục vật phẩm** 4 làng, số cấp | KH7 mở rộng chưa có spec |
| C6 | **Đường lắp hiệu ứng chiêu** (`effectdata` + `skill.bin` + `SkillPaint`) | Cần một task giải mã riêng (đọc `EffectData.cs`, `SkillPaint.cs`, `Effect.cs`) trước khi giao vẽ P1-14/P2-09 |
| C7 | Client có vẽ **chân dung boss** khi boss thoại không | Kiểm `ChatPopup.cs:316` với gói 44 của boss qua Unity MCP |
| C8 | **Khu rừng (map 6)** vào được không | `banDo.go:81-89` chỉ khai 5 map; chưa có `BanDo`, cổng, lưới va chạm cho map 6 |
| C9 | Cỡ ô chung cho **dải 4 khung** (2172 hay 2896) | Đợt 4 chốt trong `_KHUON-DOT-4.md`: ô 724×724 đều nhau → 4 khung = 2896×724 |
| C10 | Thêm part/sprite id mới (boss, chân dung, icon) phải **ký lại `asset-lock.json`** (closure v11 chưa ký lại — `BAT-DAU-TU-DAY.md` §5) | Reviewer độc lập; mỹ thuật vẽ trước, lắp sau |

---

## 8. Gợi ý tách brief đợt 4 từ bản kiểm kê này

| File brief dự kiến | Gom các mục | Ưu tiên chính |
|---|---|---|
| `01-nhan-vat-nguoi-choi.md` | P1-01, P1-02, P2-01, P2-02 | P1 (vẽ lại chờ quyết) + P2 |
| `02-boss.md` | P1-05 → P1-08, P2-06 | P1 |
| `03-npc.md` | P1-03, P1-16, P2-07 | P2 (chờ KH9) |
| `04-quai.md` | P1-04, P1-11, P2-05 | P1 đề xuất + P2 |
| `05-map-lang-la-va-khu-rung.md` | P1-09, P1-10 | P1 |
| `06-map-4-lang.md` | P2-03, P2-04 | P2 |
| `07-icon-vat-pham.md` | P1-12, P2-10, P3-07 | P1 nhỏ |
| `08-icon-chieu.md` | P1-13, P2-08 | lắp, không vẽ |
| `09-hieu-ung.md` | P1-14, P2-09 | P1 — **chờ giải mã C6** |
| `10-ui-hud.md` | P1-15, P3-01 → P3-05, P3-08 | P3 |
| `11-trang-phuc-doi-hinh.md` | P3-06 | P3 |

Mỗi brief: dòng **LƯU VÀO:** `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\<hạng mục>\`; tên file theo tiền tố `kich-thuoc-anh.json` (ảnh tĩnh) hoặc mẫu `vertical-slice-v1` (dải chroma); luôn kèm bài học 05/09 cho mọi nhân vật/NPC/quái/boss.

---

*Kiểm kê từ: `server/kho-noi-dung/*.yaml`, `server/internal/game/assets/{catalog.json, asset-lock.json}`, `server/internal/game/{banDo,hinhdang,thegioi,dongbo}.go`, `docs/superpowers/plans/LO-TRINH.md`, `docs/superpowers/specs/{2026-07-29-vertical-slice-lang-la, 2026-08-03-kh6b-boss-design, 2026-09-03-kh7-vat-pham-design, 2026-09-03-kh8-kho-noi-dung-design}.md`, `BAT-DAU-TU-DAY.md` §5, `docs/my-thuat/dot-3-toan-bo/lang-*/*.md`, `docs/my-thuat/map-1-lang-la/*.md`, `docs/my-thuat/KHUON-DE-BAI-ANH-MAP.md`, client `D:\NRO_Project\Assets\Scripts\Game1\{LoginScr, CreateCharScr, ChooseCharScr, GameScr, Panel, ChatPopup, InfoDlg, ServerListScreen, Mob, Npc, Effect, TileMap, GameCanvas, SmallImage}.cs` (chỉ đọc). Cỡ ảnh v1 đo trực tiếp bằng System.Drawing ngày 05/09/2026.*
