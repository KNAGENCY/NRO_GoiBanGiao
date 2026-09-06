# QUÁI — Đợt 4 (P1: 3 quái mới Khu rừng + Sói Rừng đã có · P2: 16 quái 4 làng — đề xuất chờ chốt)

> Đọc `_KHUON-DOT-4.md` và `../dot-3-toan-bo/_KHUON.md` trước. **LƯU VÀO:** `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\`
>
> Viết 06/09/2026 từ `00-KIEM-KE-DA-CO.md` §3.3, `00-KIEM-KE-GAME-CAN.md` P1-04/P1-11/P2-05, `_KHUON-DOT-4.md` §B.4/§C/§D, `server/kho-noi-dung/quai.yaml`, `bang-roi.yaml`, `catalog.json`, `asset-lock.json` (`closure.mob`), `tools/artqa/task6b_closure.py`, `assetpack/mob.py`, và đọc thẳng `D:\NRO_Project\Assets\Scripts\Game1\Mob.cs` + `EffectData.cs` (chỉ đọc). Chỗ nào máy chủ chưa chốt ghi **đề xuất**; chỗ nào chưa chạy thử ghi **chưa chắc**. Manifest máy đọc cùng tên: `03-quai.json`.

---

## 1. Game dùng ở đâu, vì sao cần

Quái là thứ người chơi đánh để lên cấp và nhặt đồ (thuốc máu/chakra, kunai — `bang-roi.yaml:54-72`). Máy chủ hôm nay chỉ có **một mẫu quái duy nhất, Sói Rừng** (`quai.yaml:67-74`), sinh 2 con ở **cả 6 bản đồ** với cùng toạ độ (`quai.yaml:93-121` — chính file đó ghi "thiết kế nội dung thật sẽ cho mỗi làng một đàn quái khác nhau"). Nghĩa là người chơi vào Làng Cát vẫn gặp đúng con sói Làng Lá — game chưa có cảm giác "mỗi vùng một loài".

Client vẽ quái **khác hẳn người chơi**: không ghép 3 mảnh theo bộ xương, mà đọc **một sprite sheet riêng cho mỗi mẫu** (`Mob/<templateId>/data` + `img.png`, `Mob.cs:383-393`) — nên ảnh tĩnh 96×96 `quai_*` trong lô 343 **không lắp được** (`00-KIEM-KE-GAME-CAN.md` §0 phát hiện 1). Mỗi quái mới = **4 dải chroma** (đứng 2 · chạy 4 · tấn công 2 · trúng đòn + chết 2 = 10 khung) theo đúng mẫu sói đã lắp, cộng **1 ảnh master** để duyệt thiết kế trước khi vẽ 10 khung.

Đợt này:
- **P1 — Làng Lá + Khu rừng (map 6):** giữ Sói Rừng (đã có), thêm **3 quái theo cấp** cho Khu rừng: Lợn Rừng (cấp ~3) · Nhện Độc (~5) · Khỉ Rừng (~8) — đúng 3 con `00-KIEM-KE-GAME-CAN.md` P1-11 đề xuất, đã có silhouette trong lô 343. ⚠️ Map 6 chưa vào được (`banDo.go:81-89` chỉ khai map 1–5, mục C8) — vẽ trước, lắp sau khi máy chủ nối.
- **P2 — 4 làng còn lại:** mỗi làng **4 con** (1 thú nhỏ · 1 thú lớn/rắn · 1 lính gác dạng người · 1 golem) chọn từ 8 con/làng của brief đợt 3 `../dot-3-toan-bo/lang-*/quai.md`. Tất cả là **đề xuất chờ chốt** — `quai.yaml` chưa có con nào ngoài sói (mục C2). Cấp đề xuất 10/13/16/20 chỉ để xếp thứ tự vẽ, không phải số cân bằng.

Đề xuất khối `mau_quai` cho máy chủ (chỉ là **đề xuất**, người viết `quai.yaml` chốt): `ma` như cột "mã" ở mục 5, `anh_id` = id catalog mới nối tiếp sau 0 (1, 2, 3…), `hp`/`sat_thuong_co_ban` để vòng cân bằng KH8 quyết; `anh_id` chỉ trỏ được sau khi ảnh vào `closure.mob_ids` của asset-lock (`quai.yaml:20-31`).

## 2. Đã có gì

| Tên | Nơi lưu | Dùng được không | Vì sao |
|---|---|---|---|
| `wolf_idle_2f_chroma.png` 1774×887 (ô 887) | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\` | ✅ **da_co** — đang lắp (`asset-lock.json` `closure.mob` id 0, output `client/x4/Mob/0/img.png` 960×96) | Hộp nội dung 2 khung 812×530 / 822×516, chân gần chạm đáy (744/887) |
| `wolf_run_4f_chroma.png` 2172×724 (ô **543**) | 〃 | ✅ da_co, nhưng **không làm mẫu ô** | Ô khác cỡ với 3 dải kia; 4 khung **bay** (đáy hộp 501–535/724 — hở 190+ px) nên pipeline phải `fit_to_canvas` từng khung (`task6b_closure.py:1069`) — đúng lỗi "mỗi khung một tỉ lệ" (`_KHUON-DOT-4.md` §B.4) |
| `wolf_pounce_2f_chroma.png` 1774×887 | 〃 | ✅ da_co | Khung 1 vồ, chân sau rời đáy (642/887) |
| `wolf_hit_death_2f_chroma.png` 1774×887 | 〃 | ✅ da_co | Khung 0 trúng đòn đứng chồm · khung 1 **nằm** — với quái thú chấp nhận (sheet có slot chết riêng), nhưng xem §3 về việc client có vẽ slot này không |
| `master_soi_rung_x4_chroma.png` 1254×1254 | 〃 | ✅ tham chiếu (lock: `reviewed_vertical_slice_source`, không sinh sprite) | Mẫu "ảnh master đứng" cho mọi quái mới |
| `Mob/0/img.png` x1–x4 (240×24 … 960×96) + `Mob/0.bytes` 117 byte | `asset-lock.json` outputs (closure v11 → data 28) | Bản đang ký | ⚠️ `D:\NRO_Project\Assets\Resources\res` **không có thư mục `Mob/`** — sói lên hình qua máy chủ stream / cache, đường đi thực chưa xác minh (`00-KIEM-KE-DA-CO.md` §3.3) |
| 12 ảnh tĩnh `quai_*_x4.png` 384×384 | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\05-quai\` | ❌ **không lắp** — 1 khung đứng, sai kiểu | Dùng làm **silhouette tham chiếu** (`SIL-*` trong manifest): lợn rừng, nhện độc, khỉ rừng, bọ cạp, rắn lớn, ninja giả, lính đánh thuê… `quai_soi_rung` = `TOT-quai-soi-rung.png` |
| 32 tên `quai_{cat,smu,da,may}_*` 96×96 trong `kich-thuoc-anh.json` | chỉ có tên, **chưa có ảnh** | — | Tên đợt 3; đợt 4 dùng lại **gốc tên** (bỏ `quai_`) làm mã đối tượng của dải chroma (mục 3.6) |
| Brief đợt 3 `lang-{cat,suong-mu,da,may}/quai.md` (8 con/làng, 96×96 tĩnh) | `../dot-3-toan-bo/` | ❌ sai kiểu file, đúng ý tưởng | Lấy mô tả loài/màu; bỏ cỡ 96×96 |

## 3. Ràng buộc kỹ thuật riêng loại quái

### 3.1 Kiểu file: 4 dải chroma + 1 master, ô 724×724 đều nhau

| Dải | Tên file | Số khung `n` | Cỡ cả dải | Ghi chú |
|---|---|---|---|---|
| đứng | `<obj>_idle_2f_chroma.png` | 2 | **1448×724** | lặp 0→1→0 |
| chạy | `<obj>_run_4f_chroma.png` | 4 | **2896×724** | lặp 0→1→2→3→0 (nhưng xem 3.4: client cũ chỉ lặp 3 khung) |
| tấn công | `<obj>_pounce_2f_chroma.png` | 2 | **1448×724** | tên slot là `pounce` kể cả khi đòn là chém/đập/húc/cắn (khoá cứng `assetpack/mob.py:10-11`) |
| trúng đòn + chết | `<obj>_hit_death_2f_chroma.png` | 2 | **1448×724** | khung 0 = trúng đòn (đứng) · khung 1 = chết |
| master | `master_<obj>_x4_chroma.png` | 1 | **1254×1254** | không lắp; để duyệt thiết kế và làm ảnh gốc đồng bộ cho 4 dải |

- **Số trong tên (`_2f`, `_4f`) chính là số ô máy dùng để chia** (`tools/artqa/artqa.py:450-452`); máy chia **đều theo chiều rộng** → cả dải phải đúng `n × 724`, hình không được cắn sang ô kề (`_KHUON-DOT-4.md` §A ba điều cấm). 724 chia hết cho 4 ✓.
- Vì sao 724 cho **mọi** dải và **mọi** quái (kể cả quái 0,75 ô lẫn lính gác 2,56 ô): pipeline mới sẽ đặt cỡ trong game bằng **một hệ số chung** `s = cao_mục_tiêu / cao_khung_idle_0` (như người chơi, `task6b_closure.py:69-75`) — hoạ sĩ chỉ cần vẽ **10 khung cùng tỉ lệ với nhau**, to vừa ô; cỡ thật trong game lấy từ cột "cao mục tiêu" ở mục 5, không lấy từ cỡ vẽ. Lô sói v1 để ô hai cỡ (887/543) là chuyện đã đóng gói, **không lấy làm mẫu**.

### 3.2 Nền chroma và ba điều cấm (rút từ `artqa.py:23-27, 60-62, 98-109, 112-149`)

1. Nền **tím RGB(246,4,242) phẳng tuyệt đối** — viền ảnh ≥ 85 % tím, lệch màu ≤ 12. Cả 20 file v1 đều FAIL `raw_source_status` vì nền có gradient (`review/vertical-slice-v1.json` `source_issues`) — đợt 4 phải PASS ngay bản nguồn.
2. **Không màu tím hồng/magenta trên quái** (R−G ≥ 110 và B−G ≥ 110 bị coi là nền; đảo tím trong hình bị xoá `task6b_closure.py:499-543`). Hệ quả: **Nhện Độc** (lô 343 vẽ tím) và **mọi quái Làng Mây** phải dùng **xanh chàm / blue-violet** (#3B3F8C, #5A4FBF), tuyệt đối không #FF00FF, #C040C0.
3. **Không glow, khói mờ, tia sét mờ lan ra nền** (pixel nửa trong > 15 % là FAIL). Tia điện của quái Lôi vẽ thành **nét cel-shading sắc**, không hào quang. Xem phản diện `XAU-hu-anthan-qua-mo.png`.

### 3.3 Bảy luật bài học 05/09 áp cho quái (bắt buộc, kèm `D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\anh-tham-chieu\THUOC-NHAN-VAT.png`)

Chuyện đã xảy ra: 5 dải người chơi v1 vẽ nhân vật **khác cỡ nhau** (chạy nhỏ hơn 12 %, đánh 17 %) và khom → mảnh thân mang theo mặt, client vẽ đè lên đầu → "hai khuôn mặt", phải bù toán, nhân vật thấp đi khi chạy. Sói v1 cũng mắc lỗi anh em: 4 khung chạy vẽ bay cách đáy 190 px, pipeline phải co từng khung. Với quái, luật là:

| # | Luật | Vì sao (đời thường) |
|---|---|---|
| 1 | **Mọi khung của một quái cùng CỠ** (cùng cỡ đầu, cùng dài thân) với khung `idle_0`. Quái dạng người (4 lính gác, cá mập người): **cùng cỡ đầu người chơi** theo thước — cao mặt 96–99 px, tóc→cằm 266–271 px, cằm→bàn chân 375–395 px trong ô 724 | Pipeline mới nhân **một** hệ số cho cả 10 khung; khung vẽ nhỏ hơn là quái co lại khi chạy/đánh |
| 2 | **Bàn chân / bụng / đế chạm mép dưới ô ở MỌI khung** — khung vồ vẫn để chân sau chạm đáy, không có khung bay | Máy lấy mép dưới vùng đục làm chân và neo `dy = −cao` (component sói `[-12, -24]`, `asset-lock.json` `closure.mob`). Hở 20 px là quái bay 20 px |
| 3 | **Đứng thẳng, không khom quá**; thú bốn chân chạy chỉ hạ thân vừa phải | Hộp nội dung thấp hơn 15 % là quái nhỏ đi 15 % ở khung đó |
| 4 | **Nền tím phẳng, không bóng đổ, không viền, không gradient, không sàn** | Tách nền bằng đổ nước từ mép (`artqa.py:112-149`); bóng dính chân → đo chân sai |
| 5 | **Một ô một khung, ô đều nhau, thứ tự khung cố định** như bảng 3.4 | Máy gán slot theo thứ tự; đảo hai khung là quái đánh khi đang đứng |
| 6 | **Đầu nhìn 3/4, quái quay PHẢI ở mọi khung** | Client chỉ lật gương (`EffectData.cs:296` `trans 1`) |
| 7 | **Dải hit: khung 0 trúng đòn ĐỨNG; khung 1 chết — thú gục sát đất được, quái dạng người và golem KHUỴU GỐI, không nằm** | Người chơi/boss client không bao giờ vẽ nằm (cf 23 đứng + vòng sáng); quái có slot chết riêng nên thú được gục — nhưng slot đó hiện chưa chắc lên hình (3.4). Dạng người giữ luật 7 nguyên bản để sau này có thể dùng lại làm boss |

### 3.4 Thứ tự khung, slot pipeline và **client thật tra khung nào** (đọc `Mob.cs`, chưa chạy thử — *chưa chắc*)

| Dải · khung | Slot atlas (pipeline `task6b_closure.py:1588-1595`, `mob.py:10-11`) | Ý nghĩa hoạ sĩ vẽ | Client cũ (`Mob.cs`, mẫu quái `templateId < 58`, dữ liệu `readData`) tra frame nào |
|---|---|---|---|
| idle 0 | 0 | đứng thở A | đứng/chờ: `stand[] = {0×8, 1×4}` (`Mob.cs:189-193`, `:1010`) → **frame 0** ✓ |
| idle 1 | 1 | đứng thở B | `stand[]` → **frame 1** ✓ — nhưng **cũng là khung đi số 1** (`move[]` :195) |
| run 0..3 | 2..5 | chạy vươn · thu · đạp · tiếp đất | di chuyển: `move[] = {1,1,1,1,2,2,2,2,3,3,3,3,2,2,2}`, `moveFast[] = {1,1,2,2,3,3,2}` (`:195-201`, `:1230`, `:1266`) → chỉ **frame 1, 2, 3** → hiện lên hình = idle 1 · run 0 · run 1; **run 2, run 3 chưa bao giờ hiện** |
| pounce 0, 1 | 6, 7 | lấy đà · ra đòn | đánh: **frame 10 rồi 11** (`:634-641`, `:963-970`, cờ `templateId < 58 && frame == 11` :973) → **vượt 10 khung** của atlas hiện tại; `EffectData.paintFrame` lấy `this.frame[f]` ngoài `try` (`EffectData.cs:281`) → lỗi chỉ số, khung không vẽ |
| hit | 8 | trúng đòn | client cũ **không có nhánh trúng đòn riêng** cho quái (`hurt[]` chỉ dùng khi `data.typeData == 2`, `:604-607`, `:1604-1608`) |
| death | 9 | chết | chết (`status 1`): **frame 11** + nhấp nháy `gameTick % 3` rồi biến mất (`:700-711`, `:1446`) → cùng chỉ số với khung đánh thứ hai |

Kết luận cho hoạ sĩ: **vẫn vẽ đủ 10 khung theo 4 dải** (format khoá trong `mob.py`, rẻ, và để dành khi client được sửa). Kết luận cho đường ống (**chưa có, cần viết** — nối vào danh sách `_KHUON-DOT-4.md` §D.3 mục 4): (a) bảng ánh xạ **12 frame** trong `data.bytes` (frame 0–1 đứng, 1–3 đi, 10–11 đánh/chết) trỏ về 10 ảnh con — không cần ảnh mới, chỉ đổi `frames[]`; (b) một tỉ lệ chung cho 10 khung thay `fit_to_canvas`; (c) atlas **nhiều hàng** vì `x,y,w,h` mỗi ảnh con là 1 byte ≤ 255 (`mob.py:60`, `EffectData.cs:99-102`) — 10 khung × > 25 px x1 là tràn (sói 10 × 24 = 240 vừa lọt); (d) `assetpack build` nhận nhiều `mob` (`build.py:1064-1091`), `closure.mob_ids` thêm 1, 2, 3… Trước khi bind quái thứ hai phải **kiểm bằng Unity MCP** xem sói đang đánh có vẽ khung không — nếu không, đây là lỗi P1 đang sống mà chưa ai thấy vì sói đánh nhanh.

### 3.5 Cỡ trong game (đơn vị ô 96 px x4 · người chơi 246 px = 2,56 ô)

| Nhóm | Cao mục tiêu (idle_0) | Dài tối đa | Ghi chú |
|---|---|---|---|
| thú nhỏ sát đất (nhện, bọ cạp, cua) | 72 px = 0,75 ô | 144 | rộng hơn cao — atlas cần nhiều hàng |
| thú vừa (sói, lợn, sói núi, rắn lục) | 96 px = 1 ô | 144 | sói đang ở canvas 96×96 |
| thú lớn (dê, thuỷ xà, hổ) | 120 px = 1,25 ô | 144–168 | |
| khỉ rừng (khom) | 144 px = 1,5 ô | 144 | đứng dậy khi đánh vẫn ≤ 1,9 ô |
| golem (4 làng) | 192 px = 2 ô | 144 | thấp hơn người chơi một chút, bè ngang |
| lính gác / cá mập người | 246 px = 2,56 ô | 120–144 | **cùng cỡ đầu người chơi** (luật 1) |

Giới hạn cứng: mỗi ảnh con ≤ 255 px x1 mỗi chiều (≈ 1020 px x4) — mọi con trên đều lọt. Boss dạng thú to hơn 2,6 ô **không thuộc brief này** (cần code boss dạng Mob, `00-KIEM-KE-GAME-CAN.md` P3-09).

### 3.6 Tên file và ID (append-only)

- Mã đối tượng `<obj>` = gốc tên tĩnh trong `kich-thuoc-anh.json` bỏ tiền tố `quai_`: `lon_rung`, `nhen_doc`, `khi_rung`, `cat_bocap`, `smu_camap`, `da_golem`, `may_ho_loi`… (bảng đủ ở mục 5). Sói giữ tên v1 `wolf_*` / `master_soi_rung` — **không đổi tên file đã review** (review ghim `source_sha256`, `artqa.py:629-630`).
- Dải chroma **không** cần thêm vào `kich-thuoc-anh.json` (bảng đó cho ảnh tĩnh); nhưng cổng `artqa gate`/`task6b_closure` hiện khoá đúng 20 tên v1 (`artqa.py:31-52`) → danh sách tên mới ở mục 5 là đầu vào để **nới cổng** (việc dự án).
- Mob id (`catalog.json.mobs[].id`) mới nối tiếp **1, 2, 3…** theo thứ tự lắp; `quai.yaml.anh_id` chỉ trỏ id đã ký (`quai.yaml:20-31`). `tempId` mỗi map 0..127 (`quai.yaml:90-91`). Đích: `client/Mob/<id>.bytes`, `client/Mob/<id>/data.bytes`, `client/x{n}/Mob/<id>/img.png` (`build.py:1071-1091`).

## 4. Phong cách · ảnh tham chiếu · subject block

**Phong cách khoá (dán vào mọi prompt):** `detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground` + màu hệ làng (`_KHUON.md` §3). Không còn "pixel-art" như master prompt đợt 3 — lô v1 (sói, người chơi) là hand-painted cel-shading và đã lắp, mọi quái mới phải cùng bộ với sói.

| Khoá trong manifest | Đường dẫn | Dùng để |
|---|---|---|
| `THUOC-NHAN-VAT` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\anh-tham-chieu\THUOC-NHAN-VAT.png` | thước cỡ đầu / chân — gắn kèm MỌI lần gen quái dạng người, và để hiểu 724 px ô là gì |
| `PLAYER-IDLE-GOC` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\player_idle_3f_chroma.png` | nhân vật gốc để đồng bộ cỡ đầu lính gác |
| `WOLF-IDLE-2F` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\wolf_idle_2f_chroma.png` | mẫu dải quái đã lắp: nét, độ chi tiết, chân chạm đáy |
| `WOLF-RUN-4F` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\wolf_run_4f_chroma.png` | mẫu nét chạy — nhưng **phản diện về ô/chân bay** (xem §2) |
| `WOLF-POUNCE-2F` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\wolf_pounce_2f_chroma.png` | mẫu vồ 2 khung |
| `WOLF-HIT-DEATH-2F` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\wolf_hit_death_2f_chroma.png` | mẫu trúng đòn + chết của thú |
| `MASTER-SOI-RUNG` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\master_soi_rung_x4_chroma.png` | mẫu ảnh master đứng 1254×1254 |
| `TOT-quai-soi-rung` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-quai-soi-rung.png` | mức chi tiết lông/khối cần đạt |
| `TOT-chitiet-nhanvat` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-chitiet-nhanvat.png` | mức đổ bóng, viền sạch (không vẽ lại nhân vật này) |
| `TOT-naruto-logo` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-naruto-logo.png` | bảng màu Naruto đã duyệt |
| `TOT-boss-mau` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-boss-mau.png` | độ chi tiết cho golem / quái lớn |
| `XAU-boss-zabuza` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\XAU-boss-zabuza.png` | phản diện: phẳng, chân que — tránh |
| `XAU-npc` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\XAU-npc.png` | phản diện dạng người |
| `XAU-hu-anthan-qua-mo` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\XAU-hu-anthan-qua-mo.png` | phản diện: glow/halo lan nền — tránh cho quái Lôi |
| `SIL-quai_lon_rung` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\05-quai\quai_lon_rung_x4.png` | silhouette tham chiếu dáng/màu (lô 343, ảnh tĩnh, không lắp) |
| `SIL-quai_nhen_doc` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\05-quai\quai_nhen_doc_x4.png` | silhouette tham chiếu dáng/màu (lô 343, ảnh tĩnh, không lắp) |
| `SIL-quai_khi_rung` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\05-quai\quai_khi_rung_x4.png` | silhouette tham chiếu dáng/màu (lô 343, ảnh tĩnh, không lắp) |
| `SIL-quai_bo_cap` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\05-quai\quai_bo_cap_x4.png` | silhouette tham chiếu dáng/màu (lô 343, ảnh tĩnh, không lắp) |
| `SIL-quai_ran_lon` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\05-quai\quai_ran_lon_x4.png` | silhouette tham chiếu dáng/màu (lô 343, ảnh tĩnh, không lắp) |
| `SIL-quai_ninja_gia` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\05-quai\quai_ninja_gia_x4.png` | silhouette tham chiếu dáng/màu (lô 343, ảnh tĩnh, không lắp) |
| `SIL-quai_linh_danhthue` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\05-quai\quai_linh_danhthue_x4.png` | silhouette tham chiếu dáng/màu (lô 343, ảnh tĩnh, không lắp) |

**Subject block chung (EN):** `a Naruto-world wild creature or low-rank enemy ninja, hand-painted cel-shaded game sprite, side view facing right, full body, village palette, uniform flat magenta background RGB(246,4,242), feet touching the bottom edge, same scale in every frame.` Subject riêng từng quái nằm trong prompt mục 6.

## 5. DANH SÁCH ẢNH

Tổng **100** file: da_co 5 (sói) · P1 cần vẽ 15 · P2 đề xuất 80 · P3 0. Cột "mã · cấp · nơi sinh" ghi *(đề xuất)* là chưa chốt.

| STT | Tên file | Kiểu | Cỡ (px · ô×khung) | Làng | Mô tả tiếng Việt | Ưu tiên · trạng thái |
|---|---|---|---|---|---|---|
| 1 | `master_soi_rung_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Lá / Khu rừng | Sói Rừng — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. sói xám gầy, dáng thấp săn mồi, lông xám xanh bụng sáng, mắt vàng, tai dựng · mã `soi-rung (đã chốt trong quai.yaml)` · cấp 1 (đã chốt) · map 1 + mọi map (đã chốt) · cao mục tiêu 96 px x4 (1 ô), dài ≤ 144 · **đã có tại** `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\master_soi_rung_x4_chroma.png` | P1 · da_co |
| 2 | `wolf_idle_2f_chroma.png` | dai_chroma | **1774×887 (ô 887 × 2 khung — cỡ FILE THẬT, đo 06/09)** | Lá / Khu rừng | Sói Rừng — dải đứng 2 khung: đứng thở 2 nhịp (0 đứng yên · 1 ngực phồng, đầu hạ nhẹ). sói xám gầy, dáng thấp săn mồi, lông xám xanh bụng sáng, mắt vàng, tai dựng · **đã có tại** `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\wolf_idle_2f_chroma.png` | P1 · da_co |
| 3 | `wolf_run_4f_chroma.png` | dai_chroma | **2172×724 (ô 543 × 4 khung — cỡ FILE THẬT, đo 06/09)** | Lá / Khu rừng | Sói Rừng — dải chạy/di chuyển 4 khung: phi 4 nhịp: vươn · thu · đạp · tiếp đất — chân luôn chạm đáy ô. sói xám gầy, dáng thấp săn mồi, lông xám xanh bụng sáng, mắt vàng, tai dựng · **đã có tại** `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\wolf_run_4f_chroma.png` | P1 · da_co |
| 4 | `wolf_pounce_2f_chroma.png` | dai_chroma | **1774×887 (ô 887 × 2 khung — cỡ FILE THẬT, đo 06/09)** | Lá / Khu rừng | Sói Rừng — dải tấn công (tên slot: pounce) 2 khung: tấn công 2 khung: thu mình lấy đà · vồ ra (chân sau vẫn chạm đáy). sói xám gầy, dáng thấp săn mồi, lông xám xanh bụng sáng, mắt vàng, tai dựng · **đã có tại** `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\wolf_pounce_2f_chroma.png` | P1 · da_co |
| 5 | `wolf_hit_death_2f_chroma.png` | dai_chroma | **1774×887 (ô 887 × 2 khung — cỡ FILE THẬT, đo 06/09)** | Lá / Khu rừng | Sói Rừng — dải trúng đòn + chết 2 khung: 0 trúng đòn giật lùi, còn đứng · 1 chết: gục sát đất, vẫn tựa đáy ô. sói xám gầy, dáng thấp săn mồi, lông xám xanh bụng sáng, mắt vàng, tai dựng · **đã có tại** `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\wolf_hit_death_2f_chroma.png` | P1 · da_co |
| 6 | `master_lon_rung_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Lá / Khu rừng | Lợn Rừng — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. heo rừng lùn chắc, lông cứng nâu sẫm, bờm đen dựng dọc lưng, cặp nanh ngà cong sứt mẻ, mắt đỏ nhỏ · mã `lon-rung (đề xuất)` · cấp 3 (đề xuất) · Khu rừng — bìa rừng (đề xuất) · cao mục tiêu 96 px x4 (1 ô), dài ≤ 144 | P1 · can_ve |
| 7 | `lon_rung_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Lá / Khu rừng | Lợn Rừng — dải đứng 2 khung: đứng thở 2 nhịp (0 đứng yên · 1 ngực phồng, đầu hạ nhẹ). heo rừng lùn chắc, lông cứng nâu sẫm, bờm đen dựng dọc lưng, cặp nanh ngà cong sứt mẻ, mắt đỏ nhỏ | P1 · can_ve |
| 8 | `lon_rung_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Lá / Khu rừng | Lợn Rừng — dải chạy/di chuyển 4 khung: phi 4 nhịp: vươn · thu · đạp · tiếp đất — chân luôn chạm đáy ô. heo rừng lùn chắc, lông cứng nâu sẫm, bờm đen dựng dọc lưng, cặp nanh ngà cong sứt mẻ, mắt đỏ nhỏ | P1 · can_ve |
| 9 | `lon_rung_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Lá / Khu rừng | Lợn Rừng — dải tấn công (tên slot: pounce) 2 khung: húc 2 khung: cúi đầu lấy đà, chân trước cào đất · hất đầu/sừng/nanh lên (chân sau chạm đáy). heo rừng lùn chắc, lông cứng nâu sẫm, bờm đen dựng dọc lưng, cặp nanh ngà cong sứt mẻ, mắt đỏ nhỏ | P1 · can_ve |
| 10 | `lon_rung_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Lá / Khu rừng | Lợn Rừng — dải trúng đòn + chết 2 khung: 0 trúng đòn giật lùi, còn đứng · 1 chết: gục sát đất, vẫn tựa đáy ô. heo rừng lùn chắc, lông cứng nâu sẫm, bờm đen dựng dọc lưng, cặp nanh ngà cong sứt mẻ, mắt đỏ nhỏ | P1 · can_ve |
| 11 | `master_nhen_doc_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Lá / Khu rừng | Nhện Độc — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. nhện rừng khổng lồ, bụng tròn xanh chàm/đen (KHÔNG tím hồng — kỵ nền), hoa văn vàng đất, cụm mắt vàng, tám chân khớp dài · mã `nhen-doc (đề xuất)` · cấp 5 (đề xuất) · Khu rừng — rừng sâu (đề xuất) · cao mục tiêu 72 px x4 (0.75 ô), dài ≤ 144 | P1 · can_ve |
| 12 | `nhen_doc_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Lá / Khu rừng | Nhện Độc — dải đứng 2 khung: nằm rình 2 nhịp (0 sát đất · 1 nhấc thân nhẹ, hai chân trước giơ cao hơn). nhện rừng khổng lồ, bụng tròn xanh chàm/đen (KHÔNG tím hồng — kỵ nền), hoa văn vàng đất, cụm mắt vàng, tám chân khớp dài | P1 · can_ve |
| 13 | `nhen_doc_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Lá / Khu rừng | Nhện Độc — dải chạy/di chuyển 4 khung: bò 4 nhịp, chân so le — đầu chân chạm đáy ô. nhện rừng khổng lồ, bụng tròn xanh chàm/đen (KHÔNG tím hồng — kỵ nền), hoa văn vàng đất, cụm mắt vàng, tám chân khớp dài | P1 · can_ve |
| 14 | `nhen_doc_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Lá / Khu rừng | Nhện Độc — dải tấn công (tên slot: pounce) 2 khung: tấn công 2 khung: dựng thân, giơ chân trước · chồm xuống cắn. nhện rừng khổng lồ, bụng tròn xanh chàm/đen (KHÔNG tím hồng — kỵ nền), hoa văn vàng đất, cụm mắt vàng, tám chân khớp dài | P1 · can_ve |
| 15 | `nhen_doc_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Lá / Khu rừng | Nhện Độc — dải trúng đòn + chết 2 khung: 0 trúng đòn co chân giật lùi · 1 chết: bẹp xuống, chân cụp quanh thân. nhện rừng khổng lồ, bụng tròn xanh chàm/đen (KHÔNG tím hồng — kỵ nền), hoa văn vàng đất, cụm mắt vàng, tám chân khớp dài | P1 · can_ve |
| 16 | `master_khi_rung_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Lá / Khu rừng | Khỉ Rừng — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. khỉ mặt đỏ nhe răng, lông xám nâu xù, mông đỏ, tay dài chống đất · mã `khi-rung (đề xuất)` · cấp 8 (đề xuất) · Khu rừng — rừng sâu (đề xuất) · cao mục tiêu 144 px x4 (1.5 ô), dài ≤ 144 | P1 · can_ve |
| 17 | `khi_rung_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Lá / Khu rừng | Khỉ Rừng — dải đứng 2 khung: chống tay khom lưng 2 nhịp, nhe răng. khỉ mặt đỏ nhe răng, lông xám nâu xù, mông đỏ, tay dài chống đất | P1 · can_ve |
| 18 | `khi_rung_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Lá / Khu rừng | Khỉ Rừng — dải chạy/di chuyển 4 khung: phi bằng cả tay chân 4 nhịp. khỉ mặt đỏ nhe răng, lông xám nâu xù, mông đỏ, tay dài chống đất | P1 · can_ve |
| 19 | `khi_rung_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Lá / Khu rừng | Khỉ Rừng — dải tấn công (tên slot: pounce) 2 khung: đánh 2 khung: đứng dậy giơ tay hú · vung tay đập xuống. khỉ mặt đỏ nhe răng, lông xám nâu xù, mông đỏ, tay dài chống đất | P1 · can_ve |
| 20 | `khi_rung_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Lá / Khu rừng | Khỉ Rừng — dải trúng đòn + chết 2 khung: 0 trúng đòn ôm mặt lùi · 1 chết: ngồi sụp, đầu gục. khỉ mặt đỏ nhe răng, lông xám nâu xù, mông đỏ, tay dài chống đất | P1 · can_ve |
| 21 | `master_cat_bocap_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Cát | Bọ Cạp Sa Mạc — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. bọ cạp lớn vỏ giáp vàng cát ánh nâu, càng to bóng, đuôi cong giương ngòi · mã `cat-bocap (đề xuất)` · cấp 10 (đề xuất) · map 2 Làng Cát (đề xuất) · cao mục tiêu 72 px x4 (0.75 ô), dài ≤ 144 | P2 · can_ve |
| 22 | `cat_bocap_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Cát | Bọ Cạp Sa Mạc — dải đứng 2 khung: đứng 2 nhịp, đuôi giương (0 · 1 càng hé, đuôi lắc nhẹ). bọ cạp lớn vỏ giáp vàng cát ánh nâu, càng to bóng, đuôi cong giương ngòi | P2 · can_ve |
| 23 | `cat_bocap_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Cát | Bọ Cạp Sa Mạc — dải chạy/di chuyển 4 khung: bò 4 nhịp, chân so le. bọ cạp lớn vỏ giáp vàng cát ánh nâu, càng to bóng, đuôi cong giương ngòi | P2 · can_ve |
| 24 | `cat_bocap_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Cát | Bọ Cạp Sa Mạc — dải tấn công (tên slot: pounce) 2 khung: tấn công 2 khung: cuộn đuôi ra sau · quật đuôi cắm ngòi xuống trước. bọ cạp lớn vỏ giáp vàng cát ánh nâu, càng to bóng, đuôi cong giương ngòi | P2 · can_ve |
| 25 | `cat_bocap_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Cát | Bọ Cạp Sa Mạc — dải trúng đòn + chết 2 khung: 0 trúng đòn thu càng · 1 chết: bẹp, đuôi xụi. bọ cạp lớn vỏ giáp vàng cát ánh nâu, càng to bóng, đuôi cong giương ngòi | P2 · can_ve |
| 26 | `master_cat_ran_luc_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Cát | Rắn Lục Sa Mạc — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. rắn cuộn mình, vảy vàng cát viền nâu, mắt vàng con ngươi dọc, lưỡi chẻ · mã `cat-ran-luc (đề xuất)` · cấp 13 (đề xuất) · map 2 Làng Cát (đề xuất) · cao mục tiêu 96 px x4 (1 ô), dài ≤ 144 | P2 · can_ve |
| 27 | `cat_ran_luc_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Cát | Rắn Lục Sa Mạc — dải đứng 2 khung: cuộn mình ngóc đầu 2 nhịp (1 đầu lắc, lưỡi thè). rắn cuộn mình, vảy vàng cát viền nâu, mắt vàng con ngươi dọc, lưỡi chẻ | P2 · can_ve |
| 28 | `cat_ran_luc_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Cát | Rắn Lục Sa Mạc — dải chạy/di chuyển 4 khung: trườn 4 pha sóng chữ S, bụng chạm đáy suốt chiều dài. rắn cuộn mình, vảy vàng cát viền nâu, mắt vàng con ngươi dọc, lưỡi chẻ | P2 · can_ve |
| 29 | `cat_ran_luc_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Cát | Rắn Lục Sa Mạc — dải tấn công (tên slot: pounce) 2 khung: mổ 2 khung: rụt đầu thành chữ S · phóng đầu cắn về trước. rắn cuộn mình, vảy vàng cát viền nâu, mắt vàng con ngươi dọc, lưỡi chẻ | P2 · can_ve |
| 30 | `cat_ran_luc_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Cát | Rắn Lục Sa Mạc — dải trúng đòn + chết 2 khung: 0 trúng đòn giật đầu · 1 chết: cuộn bung lỏng, đầu gối lên thân. rắn cuộn mình, vảy vàng cát viền nâu, mắt vàng con ngươi dọc, lưỡi chẻ | P2 · can_ve |
| 31 | `master_cat_linh_gac_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Cát | Lính Gác Cát — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. ninja Suna cấp thấp: áo choàng nâu cát trùm đầu, khăn che mặt hở mắt, băng trán khắc đồng hồ cát, tay thủ kunai; cùng cỡ đầu người chơi · mã `cat-linh-gac (đề xuất)` · cấp 16 (đề xuất) · map 2 Làng Cát (đề xuất) · cao mục tiêu 246 px x4 (2.56 ô), dài ≤ 120 | P2 · can_ve |
| 32 | `cat_linh_gac_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Cát | Lính Gác Cát — dải đứng 2 khung: đứng gác 2 nhịp, lưng thẳng, vũ khí hạ thấp. ninja Suna cấp thấp: áo choàng nâu cát trùm đầu, khăn che mặt hở mắt, băng trán khắc đồng hồ cát, tay thủ kunai; cùng cỡ đầu người chơi | P2 · can_ve |
| 33 | `cat_linh_gac_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Cát | Lính Gác Cát — dải chạy/di chuyển 4 khung: chạy 4 nhịp: chân phải trước · lướt · chân trái trước · lướt — đầu cùng cao, không khom. ninja Suna cấp thấp: áo choàng nâu cát trùm đầu, khăn che mặt hở mắt, băng trán khắc đồng hồ cát, tay thủ kunai; cùng cỡ đầu người chơi | P2 · can_ve |
| 34 | `cat_linh_gac_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Cát | Lính Gác Cát — dải tấn công (tên slot: pounce) 2 khung: đánh 2 khung: rút vũ khí ra sau lấy đà · chém vòng về trước. ninja Suna cấp thấp: áo choàng nâu cát trùm đầu, khăn che mặt hở mắt, băng trán khắc đồng hồ cát, tay thủ kunai; cùng cỡ đầu người chơi | P2 · can_ve |
| 35 | `cat_linh_gac_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Cát | Lính Gác Cát — dải trúng đòn + chết 2 khung: 0 trúng đòn ngả lưng 10–15°, còn đứng · 1 chết: KHUỴU MỘT GỐI, không nằm. ninja Suna cấp thấp: áo choàng nâu cát trùm đầu, khăn che mặt hở mắt, băng trán khắc đồng hồ cát, tay thủ kunai; cùng cỡ đầu người chơi | P2 · can_ve |
| 36 | `master_cat_golem_cat_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Cát | Golem Cát — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. khối cát nén hình người thô, tay to, nứt nẻ lộ lõi đá nâu, mắt sáng lục nhạt (chakra Phong) · mã `cat-golem-cat (đề xuất)` · cấp 20 (đề xuất) · map 2 Làng Cát (đề xuất) · cao mục tiêu 192 px x4 (2 ô), dài ≤ 144 | P2 · can_ve |
| 37 | `cat_golem_cat_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Cát | Golem Cát — dải đứng 2 khung: đứng khom 2 nhịp, vai nhấc, khe nứt sáng lên. khối cát nén hình người thô, tay to, nứt nẻ lộ lõi đá nâu, mắt sáng lục nhạt (chakra Phong) | P2 · can_ve |
| 38 | `cat_golem_cat_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Cát | Golem Cát — dải chạy/di chuyển 4 khung: đi nặng nề 4 nhịp. khối cát nén hình người thô, tay to, nứt nẻ lộ lõi đá nâu, mắt sáng lục nhạt (chakra Phong) | P2 · can_ve |
| 39 | `cat_golem_cat_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Cát | Golem Cát — dải tấn công (tên slot: pounce) 2 khung: đập 2 khung: giơ hai nắm đấm lên đầu · đập xuống trước mặt. khối cát nén hình người thô, tay to, nứt nẻ lộ lõi đá nâu, mắt sáng lục nhạt (chakra Phong) | P2 · can_ve |
| 40 | `cat_golem_cat_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Cát | Golem Cát — dải trúng đòn + chết 2 khung: 0 trúng đòn ngửa thân · 1 chết: sụp gối, thân rũ về trước, không nằm bẹp. khối cát nén hình người thô, tay to, nứt nẻ lộ lõi đá nâu, mắt sáng lục nhạt (chakra Phong) | P2 · can_ve |
| 41 | `master_smu_cua_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Sương Mù | Cua Biển Khổng Lồ — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. cua mai xanh xám gai, hai càng to bóng nước (một càng lớn hơn), mắt cuống · mã `smu-cua (đề xuất)` · cấp 10 (đề xuất) · map 3 Làng Sương Mù (đề xuất) · cao mục tiêu 72 px x4 (0.75 ô), dài ≤ 144 | P2 · can_ve |
| 42 | `smu_cua_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Sương Mù | Cua Biển Khổng Lồ — dải đứng 2 khung: đứng giơ càng 2 nhịp. cua mai xanh xám gai, hai càng to bóng nước (một càng lớn hơn), mắt cuống | P2 · can_ve |
| 43 | `smu_cua_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Sương Mù | Cua Biển Khổng Lồ — dải chạy/di chuyển 4 khung: bò ngang 4 nhịp, chân so le. cua mai xanh xám gai, hai càng to bóng nước (một càng lớn hơn), mắt cuống | P2 · can_ve |
| 44 | `smu_cua_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Sương Mù | Cua Biển Khổng Lồ — dải tấn công (tên slot: pounce) 2 khung: kẹp 2 khung: mở càng lớn lấy đà · thọc càng kẹp về trước. cua mai xanh xám gai, hai càng to bóng nước (một càng lớn hơn), mắt cuống | P2 · can_ve |
| 45 | `smu_cua_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Sương Mù | Cua Biển Khổng Lồ — dải trúng đòn + chết 2 khung: 0 trúng đòn thu càng che mặt · 1 chết: chân xoạc bẹp, càng mở. cua mai xanh xám gai, hai càng to bóng nước (một càng lớn hơn), mắt cuống | P2 · can_ve |
| 46 | `master_smu_thuyxa_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Sương Mù | Thuỷ Xà — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. rắn nước lớn vảy xanh lam ướt, mào vây trên đầu, mắt vàng lạnh · mã `smu-thuyxa (đề xuất)` · cấp 13 (đề xuất) · map 3 Làng Sương Mù (đề xuất) · cao mục tiêu 120 px x4 (1.25 ô), dài ≤ 144 | P2 · can_ve |
| 47 | `smu_thuyxa_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Sương Mù | Thuỷ Xà — dải đứng 2 khung: cuộn mình ngóc đầu 2 nhịp (1 đầu lắc, lưỡi thè). rắn nước lớn vảy xanh lam ướt, mào vây trên đầu, mắt vàng lạnh | P2 · can_ve |
| 48 | `smu_thuyxa_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Sương Mù | Thuỷ Xà — dải chạy/di chuyển 4 khung: trườn 4 pha sóng chữ S, bụng chạm đáy suốt chiều dài. rắn nước lớn vảy xanh lam ướt, mào vây trên đầu, mắt vàng lạnh | P2 · can_ve |
| 49 | `smu_thuyxa_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Sương Mù | Thuỷ Xà — dải tấn công (tên slot: pounce) 2 khung: mổ 2 khung: rụt đầu thành chữ S · phóng đầu cắn về trước. rắn nước lớn vảy xanh lam ướt, mào vây trên đầu, mắt vàng lạnh | P2 · can_ve |
| 50 | `smu_thuyxa_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Sương Mù | Thuỷ Xà — dải trúng đòn + chết 2 khung: 0 trúng đòn giật đầu · 1 chết: cuộn bung lỏng, đầu gối lên thân. rắn nước lớn vảy xanh lam ướt, mào vây trên đầu, mắt vàng lạnh | P2 · can_ve |
| 51 | `master_smu_linh_gac_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Sương Mù | Lính Gác Sương Mù — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. ninja Kiri cấp thấp: áo xanh lam, mặt nạ che nửa mặt, băng trán 4 vạch sóng, tay thủ kiếm ngắn; cùng cỡ đầu người chơi · mã `smu-linh-gac (đề xuất)` · cấp 16 (đề xuất) · map 3 Làng Sương Mù (đề xuất) · cao mục tiêu 246 px x4 (2.56 ô), dài ≤ 120 | P2 · can_ve |
| 52 | `smu_linh_gac_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Sương Mù | Lính Gác Sương Mù — dải đứng 2 khung: đứng gác 2 nhịp, lưng thẳng, vũ khí hạ thấp. ninja Kiri cấp thấp: áo xanh lam, mặt nạ che nửa mặt, băng trán 4 vạch sóng, tay thủ kiếm ngắn; cùng cỡ đầu người chơi | P2 · can_ve |
| 53 | `smu_linh_gac_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Sương Mù | Lính Gác Sương Mù — dải chạy/di chuyển 4 khung: chạy 4 nhịp: chân phải trước · lướt · chân trái trước · lướt — đầu cùng cao, không khom. ninja Kiri cấp thấp: áo xanh lam, mặt nạ che nửa mặt, băng trán 4 vạch sóng, tay thủ kiếm ngắn; cùng cỡ đầu người chơi | P2 · can_ve |
| 54 | `smu_linh_gac_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Sương Mù | Lính Gác Sương Mù — dải tấn công (tên slot: pounce) 2 khung: đánh 2 khung: rút vũ khí ra sau lấy đà · chém vòng về trước. ninja Kiri cấp thấp: áo xanh lam, mặt nạ che nửa mặt, băng trán 4 vạch sóng, tay thủ kiếm ngắn; cùng cỡ đầu người chơi | P2 · can_ve |
| 55 | `smu_linh_gac_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Sương Mù | Lính Gác Sương Mù — dải trúng đòn + chết 2 khung: 0 trúng đòn ngả lưng 10–15°, còn đứng · 1 chết: KHUỴU MỘT GỐI, không nằm. ninja Kiri cấp thấp: áo xanh lam, mặt nạ che nửa mặt, băng trán 4 vạch sóng, tay thủ kiếm ngắn; cùng cỡ đầu người chơi | P2 · can_ve |
| 56 | `master_smu_camap_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Sương Mù | Cá Mập Người — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. quái người-cá mập đi hai chân: da xám xanh ướt bóng, vây lưng cao, mang hằn hai bên cổ, hàm răng tua tủa, mắt trắng dữ · mã `smu-camap (đề xuất)` · cấp 20 (đề xuất) · map 3 Làng Sương Mù (đề xuất) · cao mục tiêu 246 px x4 (2.56 ô), dài ≤ 144 | P2 · can_ve |
| 57 | `smu_camap_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Sương Mù | Cá Mập Người — dải đứng 2 khung: đứng gác 2 nhịp, lưng thẳng, vũ khí hạ thấp. quái người-cá mập đi hai chân: da xám xanh ướt bóng, vây lưng cao, mang hằn hai bên cổ, hàm răng tua tủa, mắt trắng dữ | P2 · can_ve |
| 58 | `smu_camap_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Sương Mù | Cá Mập Người — dải chạy/di chuyển 4 khung: chạy 4 nhịp: chân phải trước · lướt · chân trái trước · lướt — đầu cùng cao, không khom. quái người-cá mập đi hai chân: da xám xanh ướt bóng, vây lưng cao, mang hằn hai bên cổ, hàm răng tua tủa, mắt trắng dữ | P2 · can_ve |
| 59 | `smu_camap_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Sương Mù | Cá Mập Người — dải tấn công (tên slot: pounce) 2 khung: cắn 2 khung: ngửa đầu há hàm răng · chồm cắn, vuốt cào về trước. quái người-cá mập đi hai chân: da xám xanh ướt bóng, vây lưng cao, mang hằn hai bên cổ, hàm răng tua tủa, mắt trắng dữ | P2 · can_ve |
| 60 | `smu_camap_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Sương Mù | Cá Mập Người — dải trúng đòn + chết 2 khung: 0 trúng đòn ngả lưng 10–15°, còn đứng · 1 chết: KHUỴU MỘT GỐI, không nằm. quái người-cá mập đi hai chân: da xám xanh ướt bóng, vây lưng cao, mang hằn hai bên cổ, hàm răng tua tủa, mắt trắng dữ | P2 · can_ve |
| 61 | `master_da_nhen_da_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Đá | Nhện Đá Hang — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. nhện hang thân phủ mai đá xám nứt, mắt cụm đỏ nhỏ, chân khớp như đá · mã `da-nhen-da (đề xuất)` · cấp 10 (đề xuất) · map 4 Làng Đá (đề xuất) · cao mục tiêu 72 px x4 (0.75 ô), dài ≤ 144 | P2 · can_ve |
| 62 | `da_nhen_da_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Đá | Nhện Đá Hang — dải đứng 2 khung: nằm rình 2 nhịp (0 sát đất · 1 nhấc thân nhẹ, hai chân trước giơ cao hơn). nhện hang thân phủ mai đá xám nứt, mắt cụm đỏ nhỏ, chân khớp như đá | P2 · can_ve |
| 63 | `da_nhen_da_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Đá | Nhện Đá Hang — dải chạy/di chuyển 4 khung: bò 4 nhịp, chân so le — đầu chân chạm đáy ô. nhện hang thân phủ mai đá xám nứt, mắt cụm đỏ nhỏ, chân khớp như đá | P2 · can_ve |
| 64 | `da_nhen_da_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Đá | Nhện Đá Hang — dải tấn công (tên slot: pounce) 2 khung: tấn công 2 khung: dựng thân, giơ chân trước · chồm xuống cắn. nhện hang thân phủ mai đá xám nứt, mắt cụm đỏ nhỏ, chân khớp như đá | P2 · can_ve |
| 65 | `da_nhen_da_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Đá | Nhện Đá Hang — dải trúng đòn + chết 2 khung: 0 trúng đòn co chân giật lùi · 1 chết: bẹp xuống, chân cụp quanh thân. nhện hang thân phủ mai đá xám nứt, mắt cụm đỏ nhỏ, chân khớp như đá | P2 · can_ve |
| 66 | `master_da_de_nui_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Đá | Dê Núi Hoang — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. sơn dương lông xám nâu xù, cặp sừng cong dày, móng bám đá, đầu cúi thủ thế húc · mã `da-de-nui (đề xuất)` · cấp 13 (đề xuất) · map 4 Làng Đá (đề xuất) · cao mục tiêu 120 px x4 (1.25 ô), dài ≤ 144 | P2 · can_ve |
| 67 | `da_de_nui_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Đá | Dê Núi Hoang — dải đứng 2 khung: đứng thở 2 nhịp (0 đứng yên · 1 ngực phồng, đầu hạ nhẹ). sơn dương lông xám nâu xù, cặp sừng cong dày, móng bám đá, đầu cúi thủ thế húc | P2 · can_ve |
| 68 | `da_de_nui_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Đá | Dê Núi Hoang — dải chạy/di chuyển 4 khung: phi 4 nhịp: vươn · thu · đạp · tiếp đất — chân luôn chạm đáy ô. sơn dương lông xám nâu xù, cặp sừng cong dày, móng bám đá, đầu cúi thủ thế húc | P2 · can_ve |
| 69 | `da_de_nui_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Đá | Dê Núi Hoang — dải tấn công (tên slot: pounce) 2 khung: húc 2 khung: cúi đầu lấy đà, chân trước cào đất · hất đầu/sừng/nanh lên (chân sau chạm đáy). sơn dương lông xám nâu xù, cặp sừng cong dày, móng bám đá, đầu cúi thủ thế húc | P2 · can_ve |
| 70 | `da_de_nui_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Đá | Dê Núi Hoang — dải trúng đòn + chết 2 khung: 0 trúng đòn giật lùi, còn đứng · 1 chết: gục sát đất, vẫn tựa đáy ô. sơn dương lông xám nâu xù, cặp sừng cong dày, móng bám đá, đầu cúi thủ thế húc | P2 · can_ve |
| 71 | `master_da_linh_gac_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Đá | Lính Gác Đá — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. ninja Iwa cấp thấp: áo choàng nâu đất trùm đầu, ốp phiến đá trên vai, băng trán khắc hai phiến đá, tay thủ kunai; cùng cỡ đầu người chơi · mã `da-linh-gac (đề xuất)` · cấp 16 (đề xuất) · map 4 Làng Đá (đề xuất) · cao mục tiêu 246 px x4 (2.56 ô), dài ≤ 120 | P2 · can_ve |
| 72 | `da_linh_gac_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Đá | Lính Gác Đá — dải đứng 2 khung: đứng gác 2 nhịp, lưng thẳng, vũ khí hạ thấp. ninja Iwa cấp thấp: áo choàng nâu đất trùm đầu, ốp phiến đá trên vai, băng trán khắc hai phiến đá, tay thủ kunai; cùng cỡ đầu người chơi | P2 · can_ve |
| 73 | `da_linh_gac_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Đá | Lính Gác Đá — dải chạy/di chuyển 4 khung: chạy 4 nhịp: chân phải trước · lướt · chân trái trước · lướt — đầu cùng cao, không khom. ninja Iwa cấp thấp: áo choàng nâu đất trùm đầu, ốp phiến đá trên vai, băng trán khắc hai phiến đá, tay thủ kunai; cùng cỡ đầu người chơi | P2 · can_ve |
| 74 | `da_linh_gac_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Đá | Lính Gác Đá — dải tấn công (tên slot: pounce) 2 khung: đánh 2 khung: rút vũ khí ra sau lấy đà · chém vòng về trước. ninja Iwa cấp thấp: áo choàng nâu đất trùm đầu, ốp phiến đá trên vai, băng trán khắc hai phiến đá, tay thủ kunai; cùng cỡ đầu người chơi | P2 · can_ve |
| 75 | `da_linh_gac_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Đá | Lính Gác Đá — dải trúng đòn + chết 2 khung: 0 trúng đòn ngả lưng 10–15°, còn đứng · 1 chết: KHUỴU MỘT GỐI, không nằm. ninja Iwa cấp thấp: áo choàng nâu đất trùm đầu, ốp phiến đá trên vai, băng trán khắc hai phiến đá, tay thủ kunai; cùng cỡ đầu người chơi | P2 · can_ve |
| 76 | `master_da_golem_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Đá | Golem Đá — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. khối đá nén hình người thô kệch, tay to nặng, khớp nứt lộ mạch khoáng cam nâu, rêu bám vai · mã `da-golem (đề xuất)` · cấp 20 (đề xuất) · map 4 Làng Đá (đề xuất) · cao mục tiêu 192 px x4 (2 ô), dài ≤ 144 | P2 · can_ve |
| 77 | `da_golem_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Đá | Golem Đá — dải đứng 2 khung: đứng khom 2 nhịp, vai nhấc, khe nứt sáng lên. khối đá nén hình người thô kệch, tay to nặng, khớp nứt lộ mạch khoáng cam nâu, rêu bám vai | P2 · can_ve |
| 78 | `da_golem_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Đá | Golem Đá — dải chạy/di chuyển 4 khung: đi nặng nề 4 nhịp. khối đá nén hình người thô kệch, tay to nặng, khớp nứt lộ mạch khoáng cam nâu, rêu bám vai | P2 · can_ve |
| 79 | `da_golem_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Đá | Golem Đá — dải tấn công (tên slot: pounce) 2 khung: đập 2 khung: giơ hai nắm đấm lên đầu · đập xuống trước mặt. khối đá nén hình người thô kệch, tay to nặng, khớp nứt lộ mạch khoáng cam nâu, rêu bám vai | P2 · can_ve |
| 80 | `da_golem_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Đá | Golem Đá — dải trúng đòn + chết 2 khung: 0 trúng đòn ngửa thân · 1 chết: sụp gối, thân rũ về trước, không nằm bẹp. khối đá nén hình người thô kệch, tay to nặng, khớp nứt lộ mạch khoáng cam nâu, rêu bám vai | P2 · can_ve |
| 81 | `master_may_soi_nui_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Mây | Sói Núi Lôi — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. sói xám lông dựng, viền điện xanh chàm dọc lưng (nét sắc, không glow), răng nhe, mắt vàng · mã `may-soi-nui (đề xuất)` · cấp 10 (đề xuất) · map 5 Làng Mây (đề xuất) · cao mục tiêu 96 px x4 (1 ô), dài ≤ 144 | P2 · can_ve |
| 82 | `may_soi_nui_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Mây | Sói Núi Lôi — dải đứng 2 khung: đứng thở 2 nhịp (0 đứng yên · 1 ngực phồng, đầu hạ nhẹ). sói xám lông dựng, viền điện xanh chàm dọc lưng (nét sắc, không glow), răng nhe, mắt vàng | P2 · can_ve |
| 83 | `may_soi_nui_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Mây | Sói Núi Lôi — dải chạy/di chuyển 4 khung: phi 4 nhịp: vươn · thu · đạp · tiếp đất — chân luôn chạm đáy ô. sói xám lông dựng, viền điện xanh chàm dọc lưng (nét sắc, không glow), răng nhe, mắt vàng | P2 · can_ve |
| 84 | `may_soi_nui_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Mây | Sói Núi Lôi — dải tấn công (tên slot: pounce) 2 khung: tấn công 2 khung: thu mình lấy đà · vồ ra (chân sau vẫn chạm đáy). sói xám lông dựng, viền điện xanh chàm dọc lưng (nét sắc, không glow), răng nhe, mắt vàng | P2 · can_ve |
| 85 | `may_soi_nui_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Mây | Sói Núi Lôi — dải trúng đòn + chết 2 khung: 0 trúng đòn giật lùi, còn đứng · 1 chết: gục sát đất, vẫn tựa đáy ô. sói xám lông dựng, viền điện xanh chàm dọc lưng (nét sắc, không glow), răng nhe, mắt vàng | P2 · can_ve |
| 86 | `master_may_ho_loi_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Mây | Hổ Lôi — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. hổ cơ bắp lông trắng xanh vằn chàm đen, hồ quang điện chạy dọc sống lưng (nét sắc), nanh trắng, mắt lam · mã `may-ho-loi (đề xuất)` · cấp 13 (đề xuất) · map 5 Làng Mây (đề xuất) · cao mục tiêu 120 px x4 (1.25 ô), dài ≤ 168 | P2 · can_ve |
| 87 | `may_ho_loi_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Mây | Hổ Lôi — dải đứng 2 khung: đứng thở 2 nhịp (0 đứng yên · 1 ngực phồng, đầu hạ nhẹ). hổ cơ bắp lông trắng xanh vằn chàm đen, hồ quang điện chạy dọc sống lưng (nét sắc), nanh trắng, mắt lam | P2 · can_ve |
| 88 | `may_ho_loi_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Mây | Hổ Lôi — dải chạy/di chuyển 4 khung: phi 4 nhịp: vươn · thu · đạp · tiếp đất — chân luôn chạm đáy ô. hổ cơ bắp lông trắng xanh vằn chàm đen, hồ quang điện chạy dọc sống lưng (nét sắc), nanh trắng, mắt lam | P2 · can_ve |
| 89 | `may_ho_loi_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Mây | Hổ Lôi — dải tấn công (tên slot: pounce) 2 khung: tấn công 2 khung: thu mình lấy đà · vồ ra (chân sau vẫn chạm đáy). hổ cơ bắp lông trắng xanh vằn chàm đen, hồ quang điện chạy dọc sống lưng (nét sắc), nanh trắng, mắt lam | P2 · can_ve |
| 90 | `may_ho_loi_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Mây | Hổ Lôi — dải trúng đòn + chết 2 khung: 0 trúng đòn giật lùi, còn đứng · 1 chết: gục sát đất, vẫn tựa đáy ô. hổ cơ bắp lông trắng xanh vằn chàm đen, hồ quang điện chạy dọc sống lưng (nét sắc), nanh trắng, mắt lam | P2 · can_ve |
| 91 | `master_may_linh_gac_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Mây | Lính Gác Mây — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. ninja Kumo cấp thấp: áo cổ cao trắng xanh, băng trán khắc khối mây, tay thủ đoản đao ánh điện; cùng cỡ đầu người chơi · mã `may-linh-gac (đề xuất)` · cấp 16 (đề xuất) · map 5 Làng Mây (đề xuất) · cao mục tiêu 246 px x4 (2.56 ô), dài ≤ 120 | P2 · can_ve |
| 92 | `may_linh_gac_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Mây | Lính Gác Mây — dải đứng 2 khung: đứng gác 2 nhịp, lưng thẳng, vũ khí hạ thấp. ninja Kumo cấp thấp: áo cổ cao trắng xanh, băng trán khắc khối mây, tay thủ đoản đao ánh điện; cùng cỡ đầu người chơi | P2 · can_ve |
| 93 | `may_linh_gac_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Mây | Lính Gác Mây — dải chạy/di chuyển 4 khung: chạy 4 nhịp: chân phải trước · lướt · chân trái trước · lướt — đầu cùng cao, không khom. ninja Kumo cấp thấp: áo cổ cao trắng xanh, băng trán khắc khối mây, tay thủ đoản đao ánh điện; cùng cỡ đầu người chơi | P2 · can_ve |
| 94 | `may_linh_gac_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Mây | Lính Gác Mây — dải tấn công (tên slot: pounce) 2 khung: đánh 2 khung: rút vũ khí ra sau lấy đà · chém vòng về trước. ninja Kumo cấp thấp: áo cổ cao trắng xanh, băng trán khắc khối mây, tay thủ đoản đao ánh điện; cùng cỡ đầu người chơi | P2 · can_ve |
| 95 | `may_linh_gac_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Mây | Lính Gác Mây — dải trúng đòn + chết 2 khung: 0 trúng đòn ngả lưng 10–15°, còn đứng · 1 chết: KHUỴU MỘT GỐI, không nằm. ninja Kumo cấp thấp: áo cổ cao trắng xanh, băng trán khắc khối mây, tay thủ đoản đao ánh điện; cùng cỡ đầu người chơi | P2 · can_ve |
| 96 | `master_may_golem_set_x4_chroma.png` | dai_chroma | 1254×1254 (1 ô) | Mây | Golem Sét — ảnh master đứng (1 ô): mẫu thiết kế, màu, tỉ lệ, cỡ đầu cho 4 dải; KHÔNG lắp vào game, chỉ để duyệt trước khi vẽ dải. khối đá granite hình người thô, khe nứt lộ lõi điện xanh chàm (nét sắc, không glow), hồ quang chạy giữa các khối · mã `may-golem-set (đề xuất)` · cấp 20 (đề xuất) · map 5 Làng Mây (đề xuất) · cao mục tiêu 192 px x4 (2 ô), dài ≤ 144 | P2 · can_ve |
| 97 | `may_golem_set_idle_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Mây | Golem Sét — dải đứng 2 khung: đứng khom 2 nhịp, vai nhấc, khe nứt sáng lên. khối đá granite hình người thô, khe nứt lộ lõi điện xanh chàm (nét sắc, không glow), hồ quang chạy giữa các khối | P2 · can_ve |
| 98 | `may_golem_set_run_4f_chroma.png` | dai_chroma | 2896×724 (ô 724 × 4 khung) | Mây | Golem Sét — dải chạy/di chuyển 4 khung: đi nặng nề 4 nhịp. khối đá granite hình người thô, khe nứt lộ lõi điện xanh chàm (nét sắc, không glow), hồ quang chạy giữa các khối | P2 · can_ve |
| 99 | `may_golem_set_pounce_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Mây | Golem Sét — dải tấn công (tên slot: pounce) 2 khung: đập 2 khung: giơ hai nắm đấm lên đầu · đập xuống trước mặt. khối đá granite hình người thô, khe nứt lộ lõi điện xanh chàm (nét sắc, không glow), hồ quang chạy giữa các khối | P2 · can_ve |
| 100 | `may_golem_set_hit_death_2f_chroma.png` | dai_chroma | 1448×724 (ô 724 × 2 khung) | Mây | Golem Sét — dải trúng đòn + chết 2 khung: 0 trúng đòn ngửa thân · 1 chết: sụp gối, thân rũ về trước, không nằm bẹp. khối đá granite hình người thô, khe nứt lộ lõi điện xanh chàm (nét sắc, không glow), hồ quang chạy giữa các khối | P2 · can_ve |

Thứ tự vẽ đề nghị: master → idle → run → pounce → hit_death cho **từng con**, xong con này mới sang con kia (để 10 khung cùng một phiên, cùng tỉ lệ). P1 theo cấp: Lợn Rừng → Nhện Độc → Khỉ Rừng. P2 theo làng: Cát → Sương Mù → Đá → Mây (thứ tự `2026-08-23-lo-trinh-4-map-con-lai.md`).

## 6. PROMPT TỪNG ẢNH

Mỗi khối là một prompt tiếng Anh **tự đứng được** (đã gồm phong cách khoá + màu làng + subject + nội dung từng khung + cỡ/ô/số khung + nền + góc nhìn + điều không vẽ). Gửi kèm ảnh tham chiếu theo cột `dung_tham_chieu` của file đó trong `03-quai.json`. Với sói (da_co) prompt chỉ để **vẽ lại nếu chủ dự án quyết đồng bộ ô 724** — không vẽ đợt này.

### Sói Rừng — `soi_rung` · Lá / Khu rừng · P1 · ĐÃ CÓ (không vẽ)

**1. `master_soi_rung_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a lean grey forest wolf - grey-blue coat with lighter cream belly and muzzle, bushy tail, upright pointed ears, sharp yellow eyes, black nose, dark claws. Mob id/name for the game: Sói Rừng (soi-rung).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Sói Rừng (a lean grey forest wolf) standing alert in profile, all four paws planted on the bottom edge, weight even, head level, ears up, tail relaxed and low, mouth closed. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.

PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 96 px at x4 (= 1 tile(s) of 96 px; player = 246 px), length up to 144 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_soi_rung_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**2. `wolf_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a lean grey forest wolf - grey-blue coat with lighter cream belly and muzzle, bushy tail, upright pointed ears, sharp yellow eyes, black nose, dark claws. Mob id/name for the game: Sói Rừng (soi-rung).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Sói Rừng (a lean grey forest wolf) standing alert in profile, all four paws planted on the bottom edge, weight even, head level, ears up, tail relaxed and low, mouth closed.
  cell 2 (frame 1) - Same stance, subtle breathing: chest slightly fuller, head dips 5-8 px, tail lifts a little - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 96 px at x4 zoom (= 1 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: wolf_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**3. `wolf_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a lean grey forest wolf - grey-blue coat with lighter cream belly and muzzle, bushy tail, upright pointed ears, sharp yellow eyes, black nose, dark claws. Mob id/name for the game: Sói Rừng (soi-rung).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Gallop 1 (reach): front legs stretched forward, hind legs extended back, body long and low, paws on the bottom edge.
  cell 2 (frame 1) - Gallop 2 (gather): all four legs pulled under the body, back arched high, paws bunched together on the bottom edge.
  cell 3 (frame 2) - Gallop 3 (push-off): hind legs driving back, front legs reaching up and forward, chest lifted, hind paws on the bottom edge.
  cell 4 (frame 3) - Gallop 4 (land): front paws planted, hind legs swinging forward, tail streaming behind - leads back into frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 96 px at x4 zoom (= 1 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: wolf_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**4. `wolf_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a lean grey forest wolf - grey-blue coat with lighter cream belly and muzzle, bushy tail, upright pointed ears, sharp yellow eyes, black nose, dark claws. Mob id/name for the game: Sói Rừng (soi-rung).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Attack wind-up: crouched low, hindquarters coiled, head lowered, ears back, jaws opening, hackles raised, all paws on the bottom edge.
  cell 2 (frame 1) - Attack strike: lunging forward and up with front paws and claws spread, jaws wide showing fangs, hind paws still on the bottom edge (no fully airborne frame).
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 96 px at x4 zoom (= 1 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: wolf_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**5. `wolf_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a lean grey forest wolf - grey-blue coat with lighter cream belly and muzzle, bushy tail, upright pointed ears, sharp yellow eyes, black nose, dark claws. Mob id/name for the game: Sói Rừng (soi-rung).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction: recoiling, head thrown back, eyes squeezed shut, one front paw lifted, body leaning away but still standing on the bottom edge.
  cell 2 (frame 1) - Death: collapsed low on its belly and side, legs folded, head resting on the paws, eyes closed - the whole body still rests on the bottom edge, nothing floats.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 96 px at x4 zoom (= 1 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: wolf_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Lợn Rừng — `lon_rung` · Lá / Khu rừng · P1

**6. `master_lon_rung_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a wild boar - stocky low body, coarse dark-brown bristly coat with a raised black mane along the spine, grey-brown belly, short thick legs with dark hooves, small angry red eyes, two curved chipped ivory tusks jutting from the snout, short tufted tail. Mob id/name for the game: Lợn Rừng (lon-rung).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Lợn Rừng (a wild boar) standing alert in profile, all four paws planted on the bottom edge, weight even, head level, ears up, tail relaxed and low, mouth closed. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.

PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 96 px at x4 (= 1 tile(s) of 96 px; player = 246 px), length up to 144 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_lon_rung_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**7. `lon_rung_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a wild boar - stocky low body, coarse dark-brown bristly coat with a raised black mane along the spine, grey-brown belly, short thick legs with dark hooves, small angry red eyes, two curved chipped ivory tusks jutting from the snout, short tufted tail. Mob id/name for the game: Lợn Rừng (lon-rung).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Lợn Rừng (a wild boar) standing alert in profile, all four hooves planted on the bottom edge, weight even, head level, ears up, tail relaxed and low, mouth closed.
  cell 2 (frame 1) - Same stance, subtle breathing: chest slightly fuller, head dips 5-8 px, tail lifts a little - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 96 px at x4 zoom (= 1 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: lon_rung_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**8. `lon_rung_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a wild boar - stocky low body, coarse dark-brown bristly coat with a raised black mane along the spine, grey-brown belly, short thick legs with dark hooves, small angry red eyes, two curved chipped ivory tusks jutting from the snout, short tufted tail. Mob id/name for the game: Lợn Rừng (lon-rung).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Gallop 1 (reach): front legs stretched forward, hind legs extended back, body long and low, hooves on the bottom edge.
  cell 2 (frame 1) - Gallop 2 (gather): all four legs pulled under the body, back arched high, hooves bunched together on the bottom edge.
  cell 3 (frame 2) - Gallop 3 (push-off): hind legs driving back, front legs reaching up and forward, chest lifted, hind hooves on the bottom edge.
  cell 4 (frame 3) - Gallop 4 (land): front hooves planted, hind legs swinging forward, tail streaming behind - leads back into frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 96 px at x4 zoom (= 1 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: lon_rung_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**9. `lon_rung_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a wild boar - stocky low body, coarse dark-brown bristly coat with a raised black mane along the spine, grey-brown belly, short thick legs with dark hooves, small angry red eyes, two curved chipped ivory tusks jutting from the snout, short tufted tail. Mob id/name for the game: Lợn Rừng (lon-rung).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Charge wind-up: head lowered until the snout almost touches the ground, horns/tusks pointing forward, one front hoof scraping back, bristles or mane standing up, all hooves on the bottom edge.
  cell 2 (frame 1) - Head-butt: head and horns/tusks thrust up and forward to the right, front hooves lifted just off the ground, hind hooves planted on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 96 px at x4 zoom (= 1 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: lon_rung_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**10. `lon_rung_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a wild boar - stocky low body, coarse dark-brown bristly coat with a raised black mane along the spine, grey-brown belly, short thick legs with dark hooves, small angry red eyes, two curved chipped ivory tusks jutting from the snout, short tufted tail. Mob id/name for the game: Lợn Rừng (lon-rung).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction: recoiling, head thrown back, eyes squeezed shut, one front paw lifted, body leaning away but still standing on the bottom edge.
  cell 2 (frame 1) - Death: collapsed low on its belly and side, legs folded, head resting on the hooves, eyes closed - the whole body still rests on the bottom edge, nothing floats.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 96 px at x4 zoom (= 1 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: lon_rung_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Nhện Độc — `nhen_doc` · Lá / Khu rừng · P1

**11. `master_nhen_doc_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a giant venomous forest spider - bulbous abdomen in deep indigo and blue-black (NOT purple, NOT magenta) with a swirling yellow-ochre venom pattern, eight long jointed legs with dark segments and pale joints, a cluster of glowing yellow eyes, small dark fangs. Mob id/name for the game: Nhện Độc (nhen-doc).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Nhện Độc (a giant venomous forest spider) resting low with all eight legs spread wide, abdomen close to the ground, front pair of legs slightly raised, fangs closed, leg tips on the bottom edge. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.

PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 72 px at x4 (= 0.75 tile(s) of 96 px; player = 246 px), length up to 144 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_nhen_doc_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**12. `nhen_doc_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a giant venomous forest spider - bulbous abdomen in deep indigo and blue-black (NOT purple, NOT magenta) with a swirling yellow-ochre venom pattern, eight long jointed legs with dark segments and pale joints, a cluster of glowing yellow eyes, small dark fangs. Mob id/name for the game: Nhện Độc (nhen-doc).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Nhện Độc (a giant venomous forest spider) resting low with all eight legs spread wide, abdomen close to the ground, front pair of legs slightly raised, fangs closed, leg tips on the bottom edge.
  cell 2 (frame 1) - Same pose, body lifted 6-10 px on the legs, front legs raised a little higher - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 72 px at x4 zoom (= 0.75 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: nhen_doc_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**13. `nhen_doc_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a giant venomous forest spider - bulbous abdomen in deep indigo and blue-black (NOT purple, NOT magenta) with a swirling yellow-ochre venom pattern, eight long jointed legs with dark segments and pale joints, a cluster of glowing yellow eyes, small dark fangs. Mob id/name for the game: Nhện Độc (nhen-doc).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Scuttle 1: first and third legs on the near side forward, second and fourth back, body level, leg tips on the bottom edge.
  cell 2 (frame 1) - Scuttle 2: mid-stride, legs passing under the body.
  cell 3 (frame 2) - Scuttle 3: the opposite legs forward, body level.
  cell 4 (frame 3) - Scuttle 4: mid-stride again - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 72 px at x4 zoom (= 0.75 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: nhen_doc_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**14. `nhen_doc_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a giant venomous forest spider - bulbous abdomen in deep indigo and blue-black (NOT purple, NOT magenta) with a swirling yellow-ochre venom pattern, eight long jointed legs with dark segments and pale joints, a cluster of glowing yellow eyes, small dark fangs. Mob id/name for the game: Nhện Độc (nhen-doc).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Attack wind-up: rears up on the back legs, front two pairs of legs raised high, fangs spread wide, abdomen tilted down, back legs on the bottom edge.
  cell 2 (frame 1) - Attack strike: lunges forward and down to the right, fangs biting, front legs clamping forward, rear legs still on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 72 px at x4 zoom (= 0.75 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: nhen_doc_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**15. `nhen_doc_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a giant venomous forest spider - bulbous abdomen in deep indigo and blue-black (NOT purple, NOT magenta) with a swirling yellow-ochre venom pattern, eight long jointed legs with dark segments and pale joints, a cluster of glowing yellow eyes, small dark fangs. Mob id/name for the game: Nhện Độc (nhen-doc).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction: body jerks back, legs curl inward, eyes squint, still on the bottom edge.
  cell 2 (frame 1) - Death: body flattened to the ground, legs curled under and around the body, eyes dull - resting on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 72 px at x4 zoom (= 0.75 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: nhen_doc_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Khỉ Rừng — `khi_rung` · Lá / Khu rừng · P1

**16. `master_khi_rung_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a fierce forest macaque - shaggy grey-brown fur, bright red face with bared white fangs, red rump, long muscular arms, knuckle-walking stance, small angry dark eyes. Mob id/name for the game: Khỉ Rừng (khi-rung).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Khỉ Rừng (a fierce forest macaque) in a hunched knuckle-walk stance: both hands on the ground, back arched, head up, teeth bared, hands and feet on the bottom edge. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.

PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 144 px at x4 (= 1.5 tile(s) of 96 px; player = 246 px), length up to 144 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_khi_rung_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**17. `khi_rung_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a fierce forest macaque - shaggy grey-brown fur, bright red face with bared white fangs, red rump, long muscular arms, knuckle-walking stance, small angry dark eyes. Mob id/name for the game: Khỉ Rừng (khi-rung).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Khỉ Rừng (a fierce forest macaque) in a hunched knuckle-walk stance: both hands on the ground, back arched, head up, teeth bared, hands and feet on the bottom edge.
  cell 2 (frame 1) - Same stance, chest puffs up and the head lifts 8 px, lips pulled back further - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 144 px at x4 zoom (= 1.5 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: khi_rung_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**18. `khi_rung_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a fierce forest macaque - shaggy grey-brown fur, bright red face with bared white fangs, red rump, long muscular arms, knuckle-walking stance, small angry dark eyes. Mob id/name for the game: Khỉ Rừng (khi-rung).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Knuckle-gallop 1 (reach): arms forward with hands on the ground, legs pushing back, hands on the bottom edge.
  cell 2 (frame 1) - Knuckle-gallop 2 (gather): legs tucked under the body, back arched high, hands and feet close together on the bottom edge.
  cell 3 (frame 2) - Knuckle-gallop 3 (push-off): legs driving, arms lifted forward, feet on the bottom edge.
  cell 4 (frame 3) - Knuckle-gallop 4 (land): hands down, legs swinging forward - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 144 px at x4 zoom (= 1.5 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: khi_rung_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**19. `khi_rung_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a fierce forest macaque - shaggy grey-brown fur, bright red face with bared white fangs, red rump, long muscular arms, knuckle-walking stance, small angry dark eyes. Mob id/name for the game: Khỉ Rừng (khi-rung).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Attack wind-up: rising onto the hind legs, one arm raised high overhead, mouth wide open screaming, feet on the bottom edge.
  cell 2 (frame 1) - Attack strike: the raised arm slamming down and forward to the right, body lunging, feet still on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 144 px at x4 zoom (= 1.5 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: khi_rung_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**20. `khi_rung_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Leaf Village / Forest (Fire): warm earthy browns, moss and leaf greens, grey-blue fur, warm orange-red accents only.
SUBJECT: a fierce forest macaque - shaggy grey-brown fur, bright red face with bared white fangs, red rump, long muscular arms, knuckle-walking stance, small angry dark eyes. Mob id/name for the game: Khỉ Rừng (khi-rung).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction: staggering back, both arms up covering the face, eyes shut, feet on the bottom edge.
  cell 2 (frame 1) - Death: slumped sitting on its haunches, head hanging, arms limp on the ground - resting on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 144 px at x4 zoom (= 1.5 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: khi_rung_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Bọ Cạp Sa Mạc — `cat_bocap` · Cát · P2

**21. `master_cat_bocap_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a large desert scorpion - sand-yellow armoured shell with brown edges, two big glossy claws, segmented tail arched over the back ending in a dark curved stinger, eight jointed legs, small black eyes. Mob id/name for the game: Bọ Cạp Sa Mạc (cat-bocap).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Bọ Cạp Sa Mạc (a large desert scorpion) standing on eight legs, both claws (pedipalps) held forward, segmented tail arched up and forward over the back with the stinger pointing down, legs on the bottom edge. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.

PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 72 px at x4 (= 0.75 tile(s) of 96 px; player = 246 px), length up to 144 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_cat_bocap_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**22. `cat_bocap_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a large desert scorpion - sand-yellow armoured shell with brown edges, two big glossy claws, segmented tail arched over the back ending in a dark curved stinger, eight jointed legs, small black eyes. Mob id/name for the game: Bọ Cạp Sa Mạc (cat-bocap).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Bọ Cạp Sa Mạc (a large desert scorpion) standing on eight legs, both claws (pedipalps) held forward, segmented tail arched up and forward over the back with the stinger pointing down, legs on the bottom edge.
  cell 2 (frame 1) - Same pose, claws open a little and the tail sways back 8-10 px - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 72 px at x4 zoom (= 0.75 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: cat_bocap_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**23. `cat_bocap_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a large desert scorpion - sand-yellow armoured shell with brown edges, two big glossy claws, segmented tail arched over the back ending in a dark curved stinger, eight jointed legs, small black eyes. Mob id/name for the game: Bọ Cạp Sa Mạc (cat-bocap).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Scuttle 1: first and third legs on the near side forward, second and fourth back, body level, leg tips on the bottom edge.
  cell 2 (frame 1) - Scuttle 2: mid-stride, legs passing under the body.
  cell 3 (frame 2) - Scuttle 3: the opposite legs forward, body level.
  cell 4 (frame 3) - Scuttle 4: mid-stride again - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 72 px at x4 zoom (= 0.75 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: cat_bocap_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**24. `cat_bocap_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a large desert scorpion - sand-yellow armoured shell with brown edges, two big glossy claws, segmented tail arched over the back ending in a dark curved stinger, eight jointed legs, small black eyes. Mob id/name for the game: Bọ Cạp Sa Mạc (cat-bocap).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Attack wind-up: tail coiled far back, claws spread wide open, body leaning back, legs on the bottom edge.
  cell 2 (frame 1) - Attack strike: tail whips forward and down with the stinger plunging toward the ground on the right, claws snapping shut, legs on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 72 px at x4 zoom (= 0.75 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: cat_bocap_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**25. `cat_bocap_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a large desert scorpion - sand-yellow armoured shell with brown edges, two big glossy claws, segmented tail arched over the back ending in a dark curved stinger, eight jointed legs, small black eyes. Mob id/name for the game: Bọ Cạp Sa Mạc (cat-bocap).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction: body jerks back, claws pulled in over the head, tail drops, still on the bottom edge.
  cell 2 (frame 1) - Death: flattened to the ground, tail drooping behind, claws lying open on the ground - on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 72 px at x4 zoom (= 0.75 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: cat_bocap_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Rắn Lục Sa Mạc — `cat_ran_luc` · Cát · P2

**26. `master_cat_ran_luc_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a desert viper - sand-yellow scales with brown diamond edging and a pale belly, broad triangular head, yellow eyes with vertical pupils, forked tongue, faint dust on the lower coils. Mob id/name for the game: Rắn Lục Sa Mạc (cat-ran-luc).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Rắn Lục Sa Mạc (a desert viper) coiled in two loops on the ground, head raised above the coils, neck relaxed, tongue in, coils resting on the bottom edge. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.

PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 96 px at x4 (= 1 tile(s) of 96 px; player = 246 px), length up to 144 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_cat_ran_luc_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**27. `cat_ran_luc_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a desert viper - sand-yellow scales with brown diamond edging and a pale belly, broad triangular head, yellow eyes with vertical pupils, forked tongue, faint dust on the lower coils. Mob id/name for the game: Rắn Lục Sa Mạc (cat-ran-luc).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Rắn Lục Sa Mạc (a desert viper) coiled in two loops on the ground, head raised above the coils, neck relaxed, tongue in, coils resting on the bottom edge.
  cell 2 (frame 1) - Same coil, head sways 10-12 px to the right and the forked tongue flicks out - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 96 px at x4 zoom (= 1 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: cat_ran_luc_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**28. `cat_ran_luc_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a desert viper - sand-yellow scales with brown diamond edging and a pale belly, broad triangular head, yellow eyes with vertical pupils, forked tongue, faint dust on the lower coils. Mob id/name for the game: Rắn Lục Sa Mạc (cat-ran-luc).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Slither 1: body in an S-curve moving right, head forward, the curve peaking just behind the head; belly touching the bottom edge along its whole length.
  cell 2 (frame 1) - Slither 2: the S-curve travelled to mid-body, head level.
  cell 3 (frame 2) - Slither 3: the curve near the tail, head reaching forward again.
  cell 4 (frame 3) - Slither 4: body nearly straight before the next curve - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 96 px at x4 zoom (= 1 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: cat_ran_luc_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**29. `cat_ran_luc_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a desert viper - sand-yellow scales with brown diamond edging and a pale belly, broad triangular head, yellow eyes with vertical pupils, forked tongue, faint dust on the lower coils. Mob id/name for the game: Rắn Lục Sa Mạc (cat-ran-luc).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Strike wind-up: head pulled back into a tight S, neck compressed, mouth open showing fangs, coils on the bottom edge.
  cell 2 (frame 1) - Strike: head shooting forward and down to the right, mouth wide, fangs bared, tail coils still on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 96 px at x4 zoom (= 1 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: cat_ran_luc_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**30. `cat_ran_luc_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a desert viper - sand-yellow scales with brown diamond edging and a pale belly, broad triangular head, yellow eyes with vertical pupils, forked tongue, faint dust on the lower coils. Mob id/name for the game: Rắn Lục Sa Mạc (cat-ran-luc).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction: head snapped back, eyes shut, coils tightened, still on the bottom edge.
  cell 2 (frame 1) - Death: coils loosened flat on the ground, head lying on the loops, mouth slightly open - on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 96 px at x4 zoom (= 1 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: cat_ran_luc_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Lính Gác Cát — `cat_linh_gac` · Cát · P2

**31. `master_cat_linh_gac_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a low-rank Sand Village guard ninja (adult, same head size as the player reference) - hooded sand-brown cloak over a dark vest, face scarf showing only the eyes, forehead protector with an hourglass symbol, bandaged forearms, a kunai held low in the right hand, dark sandals. Mob id/name for the game: Lính Gác Cát (cat-linh-gac).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Lính Gác Cát (a low-rank Sand Village guard ninja (adult, same head size as the player reference)) standing guard in profile facing right: feet shoulder-width apart on the bottom edge, knees straight, back upright, weapon held low in the right hand, left hand slightly forward, head in 3/4 view. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.
HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 246 px at x4 (= 2.56 tile(s) of 96 px; player = 246 px), length up to 120 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_cat_linh_gac_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**32. `cat_linh_gac_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a low-rank Sand Village guard ninja (adult, same head size as the player reference) - hooded sand-brown cloak over a dark vest, face scarf showing only the eyes, forehead protector with an hourglass symbol, bandaged forearms, a kunai held low in the right hand, dark sandals. Mob id/name for the game: Lính Gác Cát (cat-linh-gac).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Lính Gác Cát (a low-rank Sand Village guard ninja (adult, same head size as the player reference)) standing guard in profile facing right: feet shoulder-width apart on the bottom edge, knees straight, back upright, weapon held low in the right hand, left hand slightly forward, head in 3/4 view.
  cell 2 (frame 1) - Same stance with a subtle weight shift: shoulders rise 4-6 px with the breath, weapon tip lifts slightly - loops back to frame 0. Head size identical to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 120 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: cat_linh_gac_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**33. `cat_linh_gac_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a low-rank Sand Village guard ninja (adult, same head size as the player reference) - hooded sand-brown cloak over a dark vest, face scarf showing only the eyes, forehead protector with an hourglass symbol, bandaged forearms, a kunai held low in the right hand, dark sandals. Mob id/name for the game: Lính Gác Cát (cat-linh-gac).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Run 1 (contact): right foot forward and planted, left leg trailing back, torso upright with only a slight forward lean, arms swinging opposite to the legs, both feet on or touching the bottom edge.
  cell 2 (frame 1) - Run 2 (pass): legs passing under the body, front foot flat, back heel lifting, torso upright.
  cell 3 (frame 2) - Run 3 (contact opposite): left foot forward and planted, right leg back.
  cell 4 (frame 3) - Run 4 (pass opposite): legs passing, weapon arm swinging back - loops back to frame 0. Head height and head size identical in all four frames; no airborne frame.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 120 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: cat_linh_gac_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**34. `cat_linh_gac_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a low-rank Sand Village guard ninja (adult, same head size as the player reference) - hooded sand-brown cloak over a dark vest, face scarf showing only the eyes, forehead protector with an hourglass symbol, bandaged forearms, a kunai held low in the right hand, dark sandals. Mob id/name for the game: Lính Gác Cát (cat-linh-gac).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Attack wind-up: weapon drawn back over the shoulder, rear leg bent, torso twisted back, eyes narrowed, both feet on the bottom edge.
  cell 2 (frame 1) - Attack strike: weapon slashed forward and down in a clean arc toward the right, front leg lunging, torso still upright - both feet on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 120 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: cat_linh_gac_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**35. `cat_linh_gac_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a low-rank Sand Village guard ninja (adult, same head size as the player reference) - hooded sand-brown cloak over a dark vest, face scarf showing only the eyes, forehead protector with an hourglass symbol, bandaged forearms, a kunai held low in the right hand, dark sandals. Mob id/name for the game: Lính Gác Cát (cat-linh-gac).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction (standing): staggering backward, torso leaning back 10-15 degrees, one arm up guarding the face, eyes shut, both feet on the bottom edge.
  cell 2 (frame 1) - Death (kneeling, NOT lying down): dropped to one knee, weapon hand planted on the ground, head bowed, other hand on the knee - knee and foot on the bottom edge, silhouette still upright.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 120 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: cat_linh_gac_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Golem Cát — `cat_golem_cat` · Cát · P2

**36. `master_cat_golem_cat_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a sand golem - a rough humanoid mass of packed sand with cracked crust, big heavy fists, a brown rock core showing in the cracks, two pale-green glowing eye-sockets, sand grains crumbling from the shoulders (kept inside the cell). Mob id/name for the game: Golem Cát (cat-golem-cat).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Golem Cát (a sand golem) standing hunched forward, massive arms hanging with knuckles near the ground, feet wide apart on the bottom edge, head sunk between the shoulders, eye-sockets glowing dully. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.

PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 192 px at x4 (= 2 tile(s) of 96 px; player = 246 px), length up to 144 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_cat_golem_cat_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**37. `cat_golem_cat_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a sand golem - a rough humanoid mass of packed sand with cracked crust, big heavy fists, a brown rock core showing in the cracks, two pale-green glowing eye-sockets, sand grains crumbling from the shoulders (kept inside the cell). Mob id/name for the game: Golem Cát (cat-golem-cat).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Golem Cát (a sand golem) standing hunched forward, massive arms hanging with knuckles near the ground, feet wide apart on the bottom edge, head sunk between the shoulders, eye-sockets glowing dully.
  cell 2 (frame 1) - Same stance, slow heave: shoulders rise 6-8 px, the glow in the cracks a little brighter, a few pebbles shift on the shoulders - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 192 px at x4 zoom (= 2 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: cat_golem_cat_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**38. `cat_golem_cat_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a sand golem - a rough humanoid mass of packed sand with cracked crust, big heavy fists, a brown rock core showing in the cracks, two pale-green glowing eye-sockets, sand grains crumbling from the shoulders (kept inside the cell). Mob id/name for the game: Golem Cát (cat-golem-cat).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Heavy walk 1: right foot forward, left arm swinging forward, torso tilted with the weight, both feet on the bottom edge.
  cell 2 (frame 1) - Heavy walk 2: legs passing, weight over the front foot, arms hanging.
  cell 3 (frame 2) - Heavy walk 3: left foot forward, right arm forward.
  cell 4 (frame 3) - Heavy walk 4: legs passing - loops back to frame 0. Both feet always touch the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 192 px at x4 zoom (= 2 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: cat_golem_cat_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**39. `cat_golem_cat_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a sand golem - a rough humanoid mass of packed sand with cracked crust, big heavy fists, a brown rock core showing in the cracks, two pale-green glowing eye-sockets, sand grains crumbling from the shoulders (kept inside the cell). Mob id/name for the game: Golem Cát (cat-golem-cat).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Attack wind-up: both fists raised high overhead, chest open, back arched, feet planted on the bottom edge.
  cell 2 (frame 1) - Attack slam: both fists smashed down in front toward the ground on the right, torso bent forward, a few small crisp rock chips near the fists (kept inside the cell, no dust cloud).
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 192 px at x4 zoom (= 2 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: cat_golem_cat_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**40. `cat_golem_cat_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Sand Village (Wind): sun-bleached sand yellow, ochre, dusty brown, pale mint-green accents, bone white.
SUBJECT: a sand golem - a rough humanoid mass of packed sand with cracked crust, big heavy fists, a brown rock core showing in the cracks, two pale-green glowing eye-sockets, sand grains crumbling from the shoulders (kept inside the cell). Mob id/name for the game: Golem Cát (cat-golem-cat).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction: torso jerks back, arms flung out, a crack flashes brighter, feet still planted on the bottom edge.
  cell 2 (frame 1) - Death: crumbling to its knees, upper body slumping forward, arms hanging, glow dimmed - knees on the bottom edge, no rubble outside the cell, not lying flat.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 192 px at x4 zoom (= 2 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: cat_golem_cat_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Cua Biển Khổng Lồ — `smu_cua` · Sương Mù · P2

**41. `master_smu_cua_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a giant sea crab - spiky blue-grey shell with teal sheen, one oversized claw and one smaller claw, wet glossy highlights, stalked eyes, eight jointed legs. Mob id/name for the game: Cua Biển Khổng Lồ (smu-cua).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Cua Biển Khổng Lồ (a giant sea crab) standing on eight walking legs, both claws raised and held toward the right, eyestalks up, legs on the bottom edge. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.

PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 72 px at x4 (= 0.75 tile(s) of 96 px; player = 246 px), length up to 144 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_smu_cua_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**42. `smu_cua_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a giant sea crab - spiky blue-grey shell with teal sheen, one oversized claw and one smaller claw, wet glossy highlights, stalked eyes, eight jointed legs. Mob id/name for the game: Cua Biển Khổng Lồ (smu-cua).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Cua Biển Khổng Lồ (a giant sea crab) standing on eight walking legs, both claws raised and held toward the right, eyestalks up, legs on the bottom edge.
  cell 2 (frame 1) - Same pose, claws open and the body dips 6 px - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 72 px at x4 zoom (= 0.75 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: smu_cua_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**43. `smu_cua_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a giant sea crab - spiky blue-grey shell with teal sheen, one oversized claw and one smaller claw, wet glossy highlights, stalked eyes, eight jointed legs. Mob id/name for the game: Cua Biển Khổng Lồ (smu-cua).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Side-scuttle 1: near-side legs stretched right, far-side legs tucked, body level, leg tips on the bottom edge.
  cell 2 (frame 1) - Side-scuttle 2: legs passing under the shell.
  cell 3 (frame 2) - Side-scuttle 3: near-side legs tucked, far-side legs stretched.
  cell 4 (frame 3) - Side-scuttle 4: legs passing - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 72 px at x4 zoom (= 0.75 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: smu_cua_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**44. `smu_cua_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a giant sea crab - spiky blue-grey shell with teal sheen, one oversized claw and one smaller claw, wet glossy highlights, stalked eyes, eight jointed legs. Mob id/name for the game: Cua Biển Khổng Lồ (smu-cua).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Attack wind-up: the big claw pulled back and opened wide, body leaning back, legs on the bottom edge.
  cell 2 (frame 1) - Attack strike: the big claw thrust forward and snapping shut toward the right, body lunging, legs still on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 72 px at x4 zoom (= 0.75 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: smu_cua_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**45. `smu_cua_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a giant sea crab - spiky blue-grey shell with teal sheen, one oversized claw and one smaller claw, wet glossy highlights, stalked eyes, eight jointed legs. Mob id/name for the game: Cua Biển Khổng Lồ (smu-cua).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction: body jerks back, both claws pulled in over the face, still on the bottom edge.
  cell 2 (frame 1) - Death: legs splayed flat, claws lying open on the ground, shell settled on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 72 px at x4 zoom (= 0.75 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: smu_cua_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Thuỷ Xà — `smu_thuyxa` · Sương Mù · P2

**46. `master_smu_thuyxa_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a large water serpent - wet deep-blue scales with teal highlights and a pale belly, a finned crest along the head and neck, cold yellow eyes, fangs, water droplets on the scales (crisp, no glow). Mob id/name for the game: Thuỷ Xà (smu-thuyxa).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Thuỷ Xà (a large water serpent) coiled in two loops on the ground, head raised above the coils, neck relaxed, tongue in, coils resting on the bottom edge. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.

PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 120 px at x4 (= 1.25 tile(s) of 96 px; player = 246 px), length up to 144 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_smu_thuyxa_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**47. `smu_thuyxa_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a large water serpent - wet deep-blue scales with teal highlights and a pale belly, a finned crest along the head and neck, cold yellow eyes, fangs, water droplets on the scales (crisp, no glow). Mob id/name for the game: Thuỷ Xà (smu-thuyxa).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Thuỷ Xà (a large water serpent) coiled in two loops on the ground, head raised above the coils, neck relaxed, tongue in, coils resting on the bottom edge.
  cell 2 (frame 1) - Same coil, head sways 10-12 px to the right and the forked tongue flicks out - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 120 px at x4 zoom (= 1.25 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: smu_thuyxa_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**48. `smu_thuyxa_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a large water serpent - wet deep-blue scales with teal highlights and a pale belly, a finned crest along the head and neck, cold yellow eyes, fangs, water droplets on the scales (crisp, no glow). Mob id/name for the game: Thuỷ Xà (smu-thuyxa).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Slither 1: body in an S-curve moving right, head forward, the curve peaking just behind the head; belly touching the bottom edge along its whole length.
  cell 2 (frame 1) - Slither 2: the S-curve travelled to mid-body, head level.
  cell 3 (frame 2) - Slither 3: the curve near the tail, head reaching forward again.
  cell 4 (frame 3) - Slither 4: body nearly straight before the next curve - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 120 px at x4 zoom (= 1.25 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: smu_thuyxa_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**49. `smu_thuyxa_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a large water serpent - wet deep-blue scales with teal highlights and a pale belly, a finned crest along the head and neck, cold yellow eyes, fangs, water droplets on the scales (crisp, no glow). Mob id/name for the game: Thuỷ Xà (smu-thuyxa).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Strike wind-up: head pulled back into a tight S, neck compressed, mouth open showing fangs, coils on the bottom edge.
  cell 2 (frame 1) - Strike: head shooting forward and down to the right, mouth wide, fangs bared, tail coils still on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 120 px at x4 zoom (= 1.25 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: smu_thuyxa_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**50. `smu_thuyxa_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a large water serpent - wet deep-blue scales with teal highlights and a pale belly, a finned crest along the head and neck, cold yellow eyes, fangs, water droplets on the scales (crisp, no glow). Mob id/name for the game: Thuỷ Xà (smu-thuyxa).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction: head snapped back, eyes shut, coils tightened, still on the bottom edge.
  cell 2 (frame 1) - Death: coils loosened flat on the ground, head lying on the loops, mouth slightly open - on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 120 px at x4 zoom (= 1.25 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: smu_thuyxa_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Lính Gác Sương Mù — `smu_linh_gac` · Sương Mù · P2

**51. `master_smu_linh_gac_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a low-rank Mist Village guard ninja (adult, same head size as the player reference) - deep-blue sleeveless vest with grey pinstripe undersuit, half face mask, forehead protector with four wavy lines, a short sword held low in the right hand, bandaged shins, dark sandals. Mob id/name for the game: Lính Gác Sương Mù (smu-linh-gac).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Lính Gác Sương Mù (a low-rank Mist Village guard ninja (adult, same head size as the player reference)) standing guard in profile facing right: feet shoulder-width apart on the bottom edge, knees straight, back upright, weapon held low in the right hand, left hand slightly forward, head in 3/4 view. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.
HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 246 px at x4 (= 2.56 tile(s) of 96 px; player = 246 px), length up to 120 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_smu_linh_gac_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**52. `smu_linh_gac_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a low-rank Mist Village guard ninja (adult, same head size as the player reference) - deep-blue sleeveless vest with grey pinstripe undersuit, half face mask, forehead protector with four wavy lines, a short sword held low in the right hand, bandaged shins, dark sandals. Mob id/name for the game: Lính Gác Sương Mù (smu-linh-gac).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Lính Gác Sương Mù (a low-rank Mist Village guard ninja (adult, same head size as the player reference)) standing guard in profile facing right: feet shoulder-width apart on the bottom edge, knees straight, back upright, weapon held low in the right hand, left hand slightly forward, head in 3/4 view.
  cell 2 (frame 1) - Same stance with a subtle weight shift: shoulders rise 4-6 px with the breath, weapon tip lifts slightly - loops back to frame 0. Head size identical to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 120 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: smu_linh_gac_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**53. `smu_linh_gac_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a low-rank Mist Village guard ninja (adult, same head size as the player reference) - deep-blue sleeveless vest with grey pinstripe undersuit, half face mask, forehead protector with four wavy lines, a short sword held low in the right hand, bandaged shins, dark sandals. Mob id/name for the game: Lính Gác Sương Mù (smu-linh-gac).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Run 1 (contact): right foot forward and planted, left leg trailing back, torso upright with only a slight forward lean, arms swinging opposite to the legs, both feet on or touching the bottom edge.
  cell 2 (frame 1) - Run 2 (pass): legs passing under the body, front foot flat, back heel lifting, torso upright.
  cell 3 (frame 2) - Run 3 (contact opposite): left foot forward and planted, right leg back.
  cell 4 (frame 3) - Run 4 (pass opposite): legs passing, weapon arm swinging back - loops back to frame 0. Head height and head size identical in all four frames; no airborne frame.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 120 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: smu_linh_gac_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**54. `smu_linh_gac_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a low-rank Mist Village guard ninja (adult, same head size as the player reference) - deep-blue sleeveless vest with grey pinstripe undersuit, half face mask, forehead protector with four wavy lines, a short sword held low in the right hand, bandaged shins, dark sandals. Mob id/name for the game: Lính Gác Sương Mù (smu-linh-gac).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Attack wind-up: weapon drawn back over the shoulder, rear leg bent, torso twisted back, eyes narrowed, both feet on the bottom edge.
  cell 2 (frame 1) - Attack strike: weapon slashed forward and down in a clean arc toward the right, front leg lunging, torso still upright - both feet on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 120 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: smu_linh_gac_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**55. `smu_linh_gac_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a low-rank Mist Village guard ninja (adult, same head size as the player reference) - deep-blue sleeveless vest with grey pinstripe undersuit, half face mask, forehead protector with four wavy lines, a short sword held low in the right hand, bandaged shins, dark sandals. Mob id/name for the game: Lính Gác Sương Mù (smu-linh-gac).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction (standing): staggering backward, torso leaning back 10-15 degrees, one arm up guarding the face, eyes shut, both feet on the bottom edge.
  cell 2 (frame 1) - Death (kneeling, NOT lying down): dropped to one knee, weapon hand planted on the ground, head bowed, other hand on the knee - knee and foot on the bottom edge, silhouette still upright.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 120 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: smu_linh_gac_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Cá Mập Người — `smu_camap` · Sương Mù · P2

**56. `master_smu_camap_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a bipedal shark-man brute - wet blue-grey skin with a pale belly, tall dorsal fin, gill slits on the neck, wide jaws full of triangular teeth, small white eyes, webbed clawed hands, broad hunched shoulders, tail fin. Mob id/name for the game: Cá Mập Người (smu-camap).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Cá Mập Người (a bipedal shark-man brute) standing in profile facing right, bipedal, broad shoulders hunched forward, webbed clawed hands open at the sides, dorsal fin up, feet on the bottom edge, head in 3/4 view. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.
HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 246 px at x4 (= 2.56 tile(s) of 96 px; player = 246 px), length up to 144 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_smu_camap_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**57. `smu_camap_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a bipedal shark-man brute - wet blue-grey skin with a pale belly, tall dorsal fin, gill slits on the neck, wide jaws full of triangular teeth, small white eyes, webbed clawed hands, broad hunched shoulders, tail fin. Mob id/name for the game: Cá Mập Người (smu-camap).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Cá Mập Người (a bipedal shark-man brute) standing in profile facing right, bipedal, broad shoulders hunched forward, webbed clawed hands open at the sides, dorsal fin up, feet on the bottom edge, head in 3/4 view.
  cell 2 (frame 1) - Same stance, gills flare and the chest rises 4-6 px, jaw parts slightly - loops back to frame 0. Head size identical to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: smu_camap_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**58. `smu_camap_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a bipedal shark-man brute - wet blue-grey skin with a pale belly, tall dorsal fin, gill slits on the neck, wide jaws full of triangular teeth, small white eyes, webbed clawed hands, broad hunched shoulders, tail fin. Mob id/name for the game: Cá Mập Người (smu-camap).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Run 1 (contact): right foot forward and planted, left leg trailing back, torso upright with only a slight forward lean, arms swinging opposite to the legs, both feet on or touching the bottom edge.
  cell 2 (frame 1) - Run 2 (pass): legs passing under the body, front foot flat, back heel lifting, torso upright.
  cell 3 (frame 2) - Run 3 (contact opposite): left foot forward and planted, right leg back.
  cell 4 (frame 3) - Run 4 (pass opposite): legs passing, weapon arm swinging back - loops back to frame 0. Head height and head size identical in all four frames; no airborne frame.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: smu_camap_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**59. `smu_camap_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a bipedal shark-man brute - wet blue-grey skin with a pale belly, tall dorsal fin, gill slits on the neck, wide jaws full of triangular teeth, small white eyes, webbed clawed hands, broad hunched shoulders, tail fin. Mob id/name for the game: Cá Mập Người (smu-camap).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Bite wind-up: head reared back, jaws opening wide with rows of teeth, both clawed hands raised, rear leg bent, feet on the bottom edge.
  cell 2 (frame 1) - Bite strike: lunging forward to the right, jaws snapping shut, claws raking forward, torso upright, both feet on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: smu_camap_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**60. `smu_camap_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Mist Village (Water): deep blue, teal, slate grey, silver highlights, wet glistening skin or scales.
SUBJECT: a bipedal shark-man brute - wet blue-grey skin with a pale belly, tall dorsal fin, gill slits on the neck, wide jaws full of triangular teeth, small white eyes, webbed clawed hands, broad hunched shoulders, tail fin. Mob id/name for the game: Cá Mập Người (smu-camap).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction (standing): staggering backward, torso leaning back 10-15 degrees, one arm up guarding the face, eyes shut, both feet on the bottom edge.
  cell 2 (frame 1) - Death (kneeling, NOT lying down): dropped to one knee, weapon hand planted on the ground, head bowed, other hand on the knee - knee and foot on the bottom edge, silhouette still upright.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: smu_camap_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Nhện Đá Hang — `da_nhen_da` · Đá · P2

**61. `master_da_nhen_da_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a cave rock spider - body armoured in cracked grey stone plates with brown dust, eight jointed stony legs, a cluster of small red eyes, dark fangs. Mob id/name for the game: Nhện Đá Hang (da-nhen-da).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Nhện Đá Hang (a cave rock spider) resting low with all eight legs spread wide, abdomen close to the ground, front pair of legs slightly raised, fangs closed, leg tips on the bottom edge. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.

PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 72 px at x4 (= 0.75 tile(s) of 96 px; player = 246 px), length up to 144 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_da_nhen_da_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**62. `da_nhen_da_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a cave rock spider - body armoured in cracked grey stone plates with brown dust, eight jointed stony legs, a cluster of small red eyes, dark fangs. Mob id/name for the game: Nhện Đá Hang (da-nhen-da).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Nhện Đá Hang (a cave rock spider) resting low with all eight legs spread wide, abdomen close to the ground, front pair of legs slightly raised, fangs closed, leg tips on the bottom edge.
  cell 2 (frame 1) - Same pose, body lifted 6-10 px on the legs, front legs raised a little higher - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 72 px at x4 zoom (= 0.75 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: da_nhen_da_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**63. `da_nhen_da_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a cave rock spider - body armoured in cracked grey stone plates with brown dust, eight jointed stony legs, a cluster of small red eyes, dark fangs. Mob id/name for the game: Nhện Đá Hang (da-nhen-da).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Scuttle 1: first and third legs on the near side forward, second and fourth back, body level, leg tips on the bottom edge.
  cell 2 (frame 1) - Scuttle 2: mid-stride, legs passing under the body.
  cell 3 (frame 2) - Scuttle 3: the opposite legs forward, body level.
  cell 4 (frame 3) - Scuttle 4: mid-stride again - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 72 px at x4 zoom (= 0.75 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: da_nhen_da_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**64. `da_nhen_da_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a cave rock spider - body armoured in cracked grey stone plates with brown dust, eight jointed stony legs, a cluster of small red eyes, dark fangs. Mob id/name for the game: Nhện Đá Hang (da-nhen-da).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Attack wind-up: rears up on the back legs, front two pairs of legs raised high, fangs spread wide, abdomen tilted down, back legs on the bottom edge.
  cell 2 (frame 1) - Attack strike: lunges forward and down to the right, fangs biting, front legs clamping forward, rear legs still on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 72 px at x4 zoom (= 0.75 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: da_nhen_da_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**65. `da_nhen_da_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a cave rock spider - body armoured in cracked grey stone plates with brown dust, eight jointed stony legs, a cluster of small red eyes, dark fangs. Mob id/name for the game: Nhện Đá Hang (da-nhen-da).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction: body jerks back, legs curl inward, eyes squint, still on the bottom edge.
  cell 2 (frame 1) - Death: body flattened to the ground, legs curled under and around the body, eyes dull - resting on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 72 px at x4 zoom (= 0.75 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: da_nhen_da_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Dê Núi Hoang — `da_de_nui` · Đá · P2

**66. `master_da_de_nui_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a wild mountain goat - shaggy grey-brown coat, thick curved ridged horns, black cloven hooves, pale beard, yellow eyes with horizontal pupils, sturdy legs. Mob id/name for the game: Dê Núi Hoang (da-de-nui).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Dê Núi Hoang (a wild mountain goat) standing alert in profile, all four paws planted on the bottom edge, weight even, head level, ears up, tail relaxed and low, mouth closed. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.

PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 120 px at x4 (= 1.25 tile(s) of 96 px; player = 246 px), length up to 144 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_da_de_nui_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**67. `da_de_nui_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a wild mountain goat - shaggy grey-brown coat, thick curved ridged horns, black cloven hooves, pale beard, yellow eyes with horizontal pupils, sturdy legs. Mob id/name for the game: Dê Núi Hoang (da-de-nui).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Dê Núi Hoang (a wild mountain goat) standing alert in profile, all four hooves planted on the bottom edge, weight even, head level, ears up, tail relaxed and low, mouth closed.
  cell 2 (frame 1) - Same stance, subtle breathing: chest slightly fuller, head dips 5-8 px, tail lifts a little - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 120 px at x4 zoom (= 1.25 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: da_de_nui_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**68. `da_de_nui_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a wild mountain goat - shaggy grey-brown coat, thick curved ridged horns, black cloven hooves, pale beard, yellow eyes with horizontal pupils, sturdy legs. Mob id/name for the game: Dê Núi Hoang (da-de-nui).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Gallop 1 (reach): front legs stretched forward, hind legs extended back, body long and low, hooves on the bottom edge.
  cell 2 (frame 1) - Gallop 2 (gather): all four legs pulled under the body, back arched high, hooves bunched together on the bottom edge.
  cell 3 (frame 2) - Gallop 3 (push-off): hind legs driving back, front legs reaching up and forward, chest lifted, hind hooves on the bottom edge.
  cell 4 (frame 3) - Gallop 4 (land): front hooves planted, hind legs swinging forward, tail streaming behind - leads back into frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 120 px at x4 zoom (= 1.25 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: da_de_nui_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**69. `da_de_nui_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a wild mountain goat - shaggy grey-brown coat, thick curved ridged horns, black cloven hooves, pale beard, yellow eyes with horizontal pupils, sturdy legs. Mob id/name for the game: Dê Núi Hoang (da-de-nui).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Charge wind-up: head lowered until the snout almost touches the ground, horns/tusks pointing forward, one front hoof scraping back, bristles or mane standing up, all hooves on the bottom edge.
  cell 2 (frame 1) - Head-butt: head and horns/tusks thrust up and forward to the right, front hooves lifted just off the ground, hind hooves planted on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 120 px at x4 zoom (= 1.25 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: da_de_nui_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**70. `da_de_nui_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a wild mountain goat - shaggy grey-brown coat, thick curved ridged horns, black cloven hooves, pale beard, yellow eyes with horizontal pupils, sturdy legs. Mob id/name for the game: Dê Núi Hoang (da-de-nui).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction: recoiling, head thrown back, eyes squeezed shut, one front paw lifted, body leaning away but still standing on the bottom edge.
  cell 2 (frame 1) - Death: collapsed low on its belly and side, legs folded, head resting on the hooves, eyes closed - the whole body still rests on the bottom edge, nothing floats.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 120 px at x4 zoom (= 1.25 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: da_de_nui_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Lính Gác Đá — `da_linh_gac` · Đá · P2

**71. `master_da_linh_gac_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a low-rank Stone Village guard ninja (adult, same head size as the player reference) - earth-brown hooded cloak, a single stone shoulder plate, red-brown undersuit, forehead protector with two stacked rocks symbol, a kunai held low in the right hand, dark sandals. Mob id/name for the game: Lính Gác Đá (da-linh-gac).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Lính Gác Đá (a low-rank Stone Village guard ninja (adult, same head size as the player reference)) standing guard in profile facing right: feet shoulder-width apart on the bottom edge, knees straight, back upright, weapon held low in the right hand, left hand slightly forward, head in 3/4 view. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.
HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 246 px at x4 (= 2.56 tile(s) of 96 px; player = 246 px), length up to 120 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_da_linh_gac_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**72. `da_linh_gac_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a low-rank Stone Village guard ninja (adult, same head size as the player reference) - earth-brown hooded cloak, a single stone shoulder plate, red-brown undersuit, forehead protector with two stacked rocks symbol, a kunai held low in the right hand, dark sandals. Mob id/name for the game: Lính Gác Đá (da-linh-gac).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Lính Gác Đá (a low-rank Stone Village guard ninja (adult, same head size as the player reference)) standing guard in profile facing right: feet shoulder-width apart on the bottom edge, knees straight, back upright, weapon held low in the right hand, left hand slightly forward, head in 3/4 view.
  cell 2 (frame 1) - Same stance with a subtle weight shift: shoulders rise 4-6 px with the breath, weapon tip lifts slightly - loops back to frame 0. Head size identical to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 120 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: da_linh_gac_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**73. `da_linh_gac_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a low-rank Stone Village guard ninja (adult, same head size as the player reference) - earth-brown hooded cloak, a single stone shoulder plate, red-brown undersuit, forehead protector with two stacked rocks symbol, a kunai held low in the right hand, dark sandals. Mob id/name for the game: Lính Gác Đá (da-linh-gac).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Run 1 (contact): right foot forward and planted, left leg trailing back, torso upright with only a slight forward lean, arms swinging opposite to the legs, both feet on or touching the bottom edge.
  cell 2 (frame 1) - Run 2 (pass): legs passing under the body, front foot flat, back heel lifting, torso upright.
  cell 3 (frame 2) - Run 3 (contact opposite): left foot forward and planted, right leg back.
  cell 4 (frame 3) - Run 4 (pass opposite): legs passing, weapon arm swinging back - loops back to frame 0. Head height and head size identical in all four frames; no airborne frame.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 120 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: da_linh_gac_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**74. `da_linh_gac_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a low-rank Stone Village guard ninja (adult, same head size as the player reference) - earth-brown hooded cloak, a single stone shoulder plate, red-brown undersuit, forehead protector with two stacked rocks symbol, a kunai held low in the right hand, dark sandals. Mob id/name for the game: Lính Gác Đá (da-linh-gac).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Attack wind-up: weapon drawn back over the shoulder, rear leg bent, torso twisted back, eyes narrowed, both feet on the bottom edge.
  cell 2 (frame 1) - Attack strike: weapon slashed forward and down in a clean arc toward the right, front leg lunging, torso still upright - both feet on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 120 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: da_linh_gac_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**75. `da_linh_gac_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a low-rank Stone Village guard ninja (adult, same head size as the player reference) - earth-brown hooded cloak, a single stone shoulder plate, red-brown undersuit, forehead protector with two stacked rocks symbol, a kunai held low in the right hand, dark sandals. Mob id/name for the game: Lính Gác Đá (da-linh-gac).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction (standing): staggering backward, torso leaning back 10-15 degrees, one arm up guarding the face, eyes shut, both feet on the bottom edge.
  cell 2 (frame 1) - Death (kneeling, NOT lying down): dropped to one knee, weapon hand planted on the ground, head bowed, other hand on the knee - knee and foot on the bottom edge, silhouette still upright.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 120 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: da_linh_gac_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Golem Đá — `da_golem` · Đá · P2

**76. `master_da_golem_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a stone golem - a rough humanoid of grey granite blocks, heavy fists, cracks glowing dull orange-brown (Earth chakra), patches of moss on the shoulders, deep-set glowing eye-sockets. Mob id/name for the game: Golem Đá (da-golem).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Golem Đá (a stone golem) standing hunched forward, massive arms hanging with knuckles near the ground, feet wide apart on the bottom edge, head sunk between the shoulders, eye-sockets glowing dully. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.

PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 192 px at x4 (= 2 tile(s) of 96 px; player = 246 px), length up to 144 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_da_golem_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**77. `da_golem_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a stone golem - a rough humanoid of grey granite blocks, heavy fists, cracks glowing dull orange-brown (Earth chakra), patches of moss on the shoulders, deep-set glowing eye-sockets. Mob id/name for the game: Golem Đá (da-golem).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Golem Đá (a stone golem) standing hunched forward, massive arms hanging with knuckles near the ground, feet wide apart on the bottom edge, head sunk between the shoulders, eye-sockets glowing dully.
  cell 2 (frame 1) - Same stance, slow heave: shoulders rise 6-8 px, the glow in the cracks a little brighter, a few pebbles shift on the shoulders - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 192 px at x4 zoom (= 2 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: da_golem_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**78. `da_golem_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a stone golem - a rough humanoid of grey granite blocks, heavy fists, cracks glowing dull orange-brown (Earth chakra), patches of moss on the shoulders, deep-set glowing eye-sockets. Mob id/name for the game: Golem Đá (da-golem).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Heavy walk 1: right foot forward, left arm swinging forward, torso tilted with the weight, both feet on the bottom edge.
  cell 2 (frame 1) - Heavy walk 2: legs passing, weight over the front foot, arms hanging.
  cell 3 (frame 2) - Heavy walk 3: left foot forward, right arm forward.
  cell 4 (frame 3) - Heavy walk 4: legs passing - loops back to frame 0. Both feet always touch the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 192 px at x4 zoom (= 2 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: da_golem_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**79. `da_golem_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a stone golem - a rough humanoid of grey granite blocks, heavy fists, cracks glowing dull orange-brown (Earth chakra), patches of moss on the shoulders, deep-set glowing eye-sockets. Mob id/name for the game: Golem Đá (da-golem).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Attack wind-up: both fists raised high overhead, chest open, back arched, feet planted on the bottom edge.
  cell 2 (frame 1) - Attack slam: both fists smashed down in front toward the ground on the right, torso bent forward, a few small crisp rock chips near the fists (kept inside the cell, no dust cloud).
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 192 px at x4 zoom (= 2 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: da_golem_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**80. `da_golem_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Stone Village (Earth): rock grey, umber brown, beige, dusty ochre, dull orange-brown glow in cracks.
SUBJECT: a stone golem - a rough humanoid of grey granite blocks, heavy fists, cracks glowing dull orange-brown (Earth chakra), patches of moss on the shoulders, deep-set glowing eye-sockets. Mob id/name for the game: Golem Đá (da-golem).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction: torso jerks back, arms flung out, a crack flashes brighter, feet still planted on the bottom edge.
  cell 2 (frame 1) - Death: crumbling to its knees, upper body slumping forward, arms hanging, glow dimmed - knees on the bottom edge, no rubble outside the cell, not lying flat.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 192 px at x4 zoom (= 2 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: da_golem_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Sói Núi Lôi — `may_soi_nui` · Mây · P2

**81. `master_may_soi_nui_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a storm mountain wolf - bristling slate-grey coat with an indigo-blue tint along the back, crisp small cel-shaded lightning sparks drawn as sharp white-blue zigzags on the spine (no soft glow), yellow eyes, bared fangs, bushy tail. Mob id/name for the game: Sói Núi Lôi (may-soi-nui).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Sói Núi Lôi (a storm mountain wolf) standing alert in profile, all four paws planted on the bottom edge, weight even, head level, ears up, tail relaxed and low, mouth closed. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.

PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 96 px at x4 (= 1 tile(s) of 96 px; player = 246 px), length up to 144 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_may_soi_nui_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**82. `may_soi_nui_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a storm mountain wolf - bristling slate-grey coat with an indigo-blue tint along the back, crisp small cel-shaded lightning sparks drawn as sharp white-blue zigzags on the spine (no soft glow), yellow eyes, bared fangs, bushy tail. Mob id/name for the game: Sói Núi Lôi (may-soi-nui).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Sói Núi Lôi (a storm mountain wolf) standing alert in profile, all four paws planted on the bottom edge, weight even, head level, ears up, tail relaxed and low, mouth closed.
  cell 2 (frame 1) - Same stance, subtle breathing: chest slightly fuller, head dips 5-8 px, tail lifts a little - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 96 px at x4 zoom (= 1 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: may_soi_nui_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**83. `may_soi_nui_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a storm mountain wolf - bristling slate-grey coat with an indigo-blue tint along the back, crisp small cel-shaded lightning sparks drawn as sharp white-blue zigzags on the spine (no soft glow), yellow eyes, bared fangs, bushy tail. Mob id/name for the game: Sói Núi Lôi (may-soi-nui).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Gallop 1 (reach): front legs stretched forward, hind legs extended back, body long and low, paws on the bottom edge.
  cell 2 (frame 1) - Gallop 2 (gather): all four legs pulled under the body, back arched high, paws bunched together on the bottom edge.
  cell 3 (frame 2) - Gallop 3 (push-off): hind legs driving back, front legs reaching up and forward, chest lifted, hind paws on the bottom edge.
  cell 4 (frame 3) - Gallop 4 (land): front paws planted, hind legs swinging forward, tail streaming behind - leads back into frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 96 px at x4 zoom (= 1 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: may_soi_nui_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**84. `may_soi_nui_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a storm mountain wolf - bristling slate-grey coat with an indigo-blue tint along the back, crisp small cel-shaded lightning sparks drawn as sharp white-blue zigzags on the spine (no soft glow), yellow eyes, bared fangs, bushy tail. Mob id/name for the game: Sói Núi Lôi (may-soi-nui).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Attack wind-up: crouched low, hindquarters coiled, head lowered, ears back, jaws opening, hackles raised, all paws on the bottom edge.
  cell 2 (frame 1) - Attack strike: lunging forward and up with front paws and claws spread, jaws wide showing fangs, hind paws still on the bottom edge (no fully airborne frame).
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 96 px at x4 zoom (= 1 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: may_soi_nui_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**85. `may_soi_nui_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a storm mountain wolf - bristling slate-grey coat with an indigo-blue tint along the back, crisp small cel-shaded lightning sparks drawn as sharp white-blue zigzags on the spine (no soft glow), yellow eyes, bared fangs, bushy tail. Mob id/name for the game: Sói Núi Lôi (may-soi-nui).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction: recoiling, head thrown back, eyes squeezed shut, one front paw lifted, body leaning away but still standing on the bottom edge.
  cell 2 (frame 1) - Death: collapsed low on its belly and side, legs folded, head resting on the paws, eyes closed - the whole body still rests on the bottom edge, nothing floats.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 96 px at x4 zoom (= 1 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: may_soi_nui_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Hổ Lôi — `may_ho_loi` · Mây · P2

**86. `master_may_ho_loi_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a lightning tiger - muscular body with white-blue fur and indigo-black stripes, sharp cel-shaded white-blue lightning arcs along the spine (crisp, no glow), pale blue eyes, white fangs, long striped tail. Mob id/name for the game: Hổ Lôi (may-ho-loi).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Hổ Lôi (a lightning tiger) standing alert in profile, all four paws planted on the bottom edge, weight even, head level, ears up, tail relaxed and low, mouth closed. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.

PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 120 px at x4 (= 1.25 tile(s) of 96 px; player = 246 px), length up to 168 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_may_ho_loi_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**87. `may_ho_loi_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a lightning tiger - muscular body with white-blue fur and indigo-black stripes, sharp cel-shaded white-blue lightning arcs along the spine (crisp, no glow), pale blue eyes, white fangs, long striped tail. Mob id/name for the game: Hổ Lôi (may-ho-loi).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Hổ Lôi (a lightning tiger) standing alert in profile, all four paws planted on the bottom edge, weight even, head level, ears up, tail relaxed and low, mouth closed.
  cell 2 (frame 1) - Same stance, subtle breathing: chest slightly fuller, head dips 5-8 px, tail lifts a little - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 120 px at x4 zoom (= 1.25 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 168 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: may_ho_loi_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**88. `may_ho_loi_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a lightning tiger - muscular body with white-blue fur and indigo-black stripes, sharp cel-shaded white-blue lightning arcs along the spine (crisp, no glow), pale blue eyes, white fangs, long striped tail. Mob id/name for the game: Hổ Lôi (may-ho-loi).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Gallop 1 (reach): front legs stretched forward, hind legs extended back, body long and low, paws on the bottom edge.
  cell 2 (frame 1) - Gallop 2 (gather): all four legs pulled under the body, back arched high, paws bunched together on the bottom edge.
  cell 3 (frame 2) - Gallop 3 (push-off): hind legs driving back, front legs reaching up and forward, chest lifted, hind paws on the bottom edge.
  cell 4 (frame 3) - Gallop 4 (land): front paws planted, hind legs swinging forward, tail streaming behind - leads back into frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 120 px at x4 zoom (= 1.25 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 168 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: may_ho_loi_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**89. `may_ho_loi_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a lightning tiger - muscular body with white-blue fur and indigo-black stripes, sharp cel-shaded white-blue lightning arcs along the spine (crisp, no glow), pale blue eyes, white fangs, long striped tail. Mob id/name for the game: Hổ Lôi (may-ho-loi).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Attack wind-up: crouched low, hindquarters coiled, head lowered, ears back, jaws opening, hackles raised, all paws on the bottom edge.
  cell 2 (frame 1) - Attack strike: lunging forward and up with front paws and claws spread, jaws wide showing fangs, hind paws still on the bottom edge (no fully airborne frame).
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 120 px at x4 zoom (= 1.25 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 168 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: may_ho_loi_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**90. `may_ho_loi_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a lightning tiger - muscular body with white-blue fur and indigo-black stripes, sharp cel-shaded white-blue lightning arcs along the spine (crisp, no glow), pale blue eyes, white fangs, long striped tail. Mob id/name for the game: Hổ Lôi (may-ho-loi).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction: recoiling, head thrown back, eyes squeezed shut, one front paw lifted, body leaning away but still standing on the bottom edge.
  cell 2 (frame 1) - Death: collapsed low on its belly and side, legs folded, head resting on the paws, eyes closed - the whole body still rests on the bottom edge, nothing floats.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 120 px at x4 zoom (= 1.25 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 168 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a fully airborne frame, a frame facing left.
```
OUTPUT: may_ho_loi_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Lính Gác Mây — `may_linh_gac` · Mây · P2

**91. `master_may_linh_gac_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a low-rank Cloud Village guard ninja (adult, same head size as the player reference) - white-and-blue high-collared vest over an indigo undersuit, forehead protector with a cloud symbol, a short blade with a thin white-blue edge highlight held low in the right hand, bandaged forearms, dark sandals. Mob id/name for the game: Lính Gác Mây (may-linh-gac).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Lính Gác Mây (a low-rank Cloud Village guard ninja (adult, same head size as the player reference)) standing guard in profile facing right: feet shoulder-width apart on the bottom edge, knees straight, back upright, weapon held low in the right hand, left hand slightly forward, head in 3/4 view. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.
HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 246 px at x4 (= 2.56 tile(s) of 96 px; player = 246 px), length up to 120 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_may_linh_gac_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**92. `may_linh_gac_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a low-rank Cloud Village guard ninja (adult, same head size as the player reference) - white-and-blue high-collared vest over an indigo undersuit, forehead protector with a cloud symbol, a short blade with a thin white-blue edge highlight held low in the right hand, bandaged forearms, dark sandals. Mob id/name for the game: Lính Gác Mây (may-linh-gac).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Lính Gác Mây (a low-rank Cloud Village guard ninja (adult, same head size as the player reference)) standing guard in profile facing right: feet shoulder-width apart on the bottom edge, knees straight, back upright, weapon held low in the right hand, left hand slightly forward, head in 3/4 view.
  cell 2 (frame 1) - Same stance with a subtle weight shift: shoulders rise 4-6 px with the breath, weapon tip lifts slightly - loops back to frame 0. Head size identical to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 120 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: may_linh_gac_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**93. `may_linh_gac_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a low-rank Cloud Village guard ninja (adult, same head size as the player reference) - white-and-blue high-collared vest over an indigo undersuit, forehead protector with a cloud symbol, a short blade with a thin white-blue edge highlight held low in the right hand, bandaged forearms, dark sandals. Mob id/name for the game: Lính Gác Mây (may-linh-gac).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Run 1 (contact): right foot forward and planted, left leg trailing back, torso upright with only a slight forward lean, arms swinging opposite to the legs, both feet on or touching the bottom edge.
  cell 2 (frame 1) - Run 2 (pass): legs passing under the body, front foot flat, back heel lifting, torso upright.
  cell 3 (frame 2) - Run 3 (contact opposite): left foot forward and planted, right leg back.
  cell 4 (frame 3) - Run 4 (pass opposite): legs passing, weapon arm swinging back - loops back to frame 0. Head height and head size identical in all four frames; no airborne frame.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 120 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: may_linh_gac_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**94. `may_linh_gac_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a low-rank Cloud Village guard ninja (adult, same head size as the player reference) - white-and-blue high-collared vest over an indigo undersuit, forehead protector with a cloud symbol, a short blade with a thin white-blue edge highlight held low in the right hand, bandaged forearms, dark sandals. Mob id/name for the game: Lính Gác Mây (may-linh-gac).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Attack wind-up: weapon drawn back over the shoulder, rear leg bent, torso twisted back, eyes narrowed, both feet on the bottom edge.
  cell 2 (frame 1) - Attack strike: weapon slashed forward and down in a clean arc toward the right, front leg lunging, torso still upright - both feet on the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 120 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: may_linh_gac_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**95. `may_linh_gac_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a low-rank Cloud Village guard ninja (adult, same head size as the player reference) - white-and-blue high-collared vest over an indigo undersuit, forehead protector with a cloud symbol, a short blade with a thin white-blue edge highlight held low in the right hand, bandaged forearms, dark sandals. Mob id/name for the game: Lính Gác Mây (may-linh-gac).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction (standing): staggering backward, torso leaning back 10-15 degrees, one arm up guarding the face, eyes shut, both feet on the bottom edge.
  cell 2 (frame 1) - Death (kneeling, NOT lying down): dropped to one knee, weapon hand planted on the ground, head bowed, other hand on the knee - knee and foot on the bottom edge, silhouette still upright.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. HUMANOID SCALE RULE (THUOC-NHAN-VAT.png): draw this adult at exactly the player's scale - face height 96-99 px, hair-top to chin 266-271 px, chin to sole 375-395 px when standing, inside the 724 px cell; upright posture like classic Dragon Ball sprites, no hunching in any frame; head drawn 3/4 toward the viewer and identical in size in all frames. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 246 px at x4 zoom (= 2.56 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 120 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: may_linh_gac_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

### Golem Sét — `may_golem_set` · Mây · P2

**96. `master_may_golem_set_x4_chroma.png`** — master đứng

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a lightning golem - a rough humanoid of dark granite blocks, cracks revealing an indigo-blue electric core drawn with crisp cel-shaded arcs (no soft glow, no magenta), pale-blue glowing eye-sockets, small rock chips on the shoulders. Mob id/name for the game: Golem Sét (may-golem-set).
FORMAT: ONE square 1254x1254 px image, a single standing pose (this is the design master, identical to idle frame 0 of the animation strips): Golem Sét (a lightning golem) standing hunched forward, massive arms hanging with knuckles near the ground, feet wide apart on the bottom edge, head sunk between the shoulders, eye-sockets glowing dully. The creature fills about 85% of the cell height (humanoids: the player scale below instead), centred horizontally, feet touching bottom edge, at least 60 px of empty magenta on the left and right.

PURPOSE: this master is the colour / proportion / head-size reference for the four animation strips (idle 2f, run 4f, pounce 2f, hit_death 2f) - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242).
VIEW: side view, head turned 3/4 toward the viewer, facing RIGHT.
BACKGROUND: flat uniform magenta RGB(246,4,242) edge to edge - no gradient, no vignette, no shadow, no glow bleeding into the magenta. Never use magenta, pink or red-violet on the creature.
IN-GAME SIZE (metadata only, do not draw): standing height about 192 px at x4 (= 2 tile(s) of 96 px; player = 246 px), length up to 144 px.
DO NOT DRAW: scenery, ground, text, rulers, borders, extra creatures, a frame facing left.
```
OUTPUT: master_may_golem_set_x4_chroma.png · 1254×1254 · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**97. `may_golem_set_idle_2f_chroma.png`** — đứng · khung 0 = đứng A (slot 0) · khung 1 = đứng B (slot 1)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a lightning golem - a rough humanoid of dark granite blocks, cracks revealing an indigo-blue electric core drawn with crisp cel-shaded arcs (no soft glow, no magenta), pale-blue glowing eye-sockets, small rock chips on the shoulders. Mob id/name for the game: Golem Sét (may-golem-set).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: idle):
  cell 1 (frame 0) - Golem Sét (a lightning golem) standing hunched forward, massive arms hanging with knuckles near the ground, feet wide apart on the bottom edge, head sunk between the shoulders, eye-sockets glowing dully.
  cell 2 (frame 1) - Same stance, slow heave: shoulders rise 6-8 px, the glow in the cracks a little brighter, a few pebbles shift on the shoulders - loops back to frame 0.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 192 px at x4 zoom (= 2 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: may_golem_set_idle_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**98. `may_golem_set_run_4f_chroma.png`** — chạy/di chuyển · khung 0..3 = chạy 1..4 (slot 2..5)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a lightning golem - a rough humanoid of dark granite blocks, cracks revealing an indigo-blue electric core drawn with crisp cel-shaded arcs (no soft glow, no magenta), pale-blue glowing eye-sockets, small rock chips on the shoulders. Mob id/name for the game: Golem Sét (may-golem-set).
FORMAT: one horizontal sprite strip PNG, exactly 4 equal cells of 724x724 px side by side = 2896x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: run):
  cell 1 (frame 0) - Heavy walk 1: right foot forward, left arm swinging forward, torso tilted with the weight, both feet on the bottom edge.
  cell 2 (frame 1) - Heavy walk 2: legs passing, weight over the front foot, arms hanging.
  cell 3 (frame 2) - Heavy walk 3: left foot forward, right arm forward.
  cell 4 (frame 3) - Heavy walk 4: legs passing - loops back to frame 0. Both feet always touch the bottom edge.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 192 px at x4 zoom (= 2 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: may_golem_set_run_4f_chroma.png · 2896×724 · 4 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**99. `may_golem_set_pounce_2f_chroma.png`** — tấn công (tên slot: pounce) · khung 0 = lấy đà (slot 6) · khung 1 = ra đòn (slot 7)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a lightning golem - a rough humanoid of dark granite blocks, cracks revealing an indigo-blue electric core drawn with crisp cel-shaded arcs (no soft glow, no magenta), pale-blue glowing eye-sockets, small rock chips on the shoulders. Mob id/name for the game: Golem Sét (may-golem-set).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: pounce):
  cell 1 (frame 0) - Attack wind-up: both fists raised high overhead, chest open, back arched, feet planted on the bottom edge.
  cell 2 (frame 1) - Attack slam: both fists smashed down in front toward the ground on the right, torso bent forward, a few small crisp rock chips near the fists (kept inside the cell, no dust cloud).
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 192 px at x4 zoom (= 2 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: may_golem_set_pounce_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

**100. `may_golem_set_hit_death_2f_chroma.png`** — trúng đòn + chết · khung 0 = trúng đòn (slot 8 pipeline) · khung 1 = chết (slot 9 pipeline)

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground. Village palette - Cloud Village (Lightning): INDIGO and blue-violet (e.g. #3B3F8C, #5A4FBF) with white-blue and pale-yellow lightning accents; NEVER magenta, pink or red-violet (keep red below green + 100 on every violet).
SUBJECT: a lightning golem - a rough humanoid of dark granite blocks, cracks revealing an indigo-blue electric core drawn with crisp cel-shaded arcs (no soft glow, no magenta), pale-blue glowing eye-sockets, small rock chips on the shoulders. Mob id/name for the game: Golem Sét (may-golem-set).
FORMAT: one horizontal sprite strip PNG, exactly 2 equal cells of 724x724 px side by side = 1448x724 px total; ONE frame per cell; frames in this fixed order, left to right (action: hit_death):
  cell 1 (frame 0) - Hit reaction: torso jerks back, arms flung out, a crack flashes brighter, feet still planted on the bottom edge.
  cell 2 (frame 1) - Death: crumbling to its knees, upper body slumping forward, arms hanging, glow dimmed - knees on the bottom edge, no rubble outside the cell, not lying flat.
SCALE / ANCHOR (strict): every frame is drawn at the SAME scale as this creature's idle frame 0 - same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). No frame floats above the bottom edge, no ground line, no shadow under the feet. Keep at least 24 px of empty magenta on the left and right inside each cell; nothing crosses into the neighbouring cell. 
VIEW: side view, head turned 3/4 toward the viewer, the creature faces RIGHT in every frame (the game mirrors it for the left direction).
BACKGROUND: flat uniform magenta RGB(246,4,242) filling each cell edge to edge - no gradient, no vignette, no texture, no soft glow or semi-transparent smoke bleeding into the magenta. Never use magenta, pink or red-violet anywhere on the creature (indigo / blue-violet is fine).
IN-GAME SIZE (metadata for the asset pipeline, do not draw any ruler or text): standing height about 192 px at x4 zoom (= 2 tile(s) of 96 px; the player character is 246 px = 2.56 tiles), body length up to 144 px at x4.
DO NOT DRAW: scenery, ground, text, numbers, cell borders, motion blur, dust or debris crossing a cell edge, extra creatures, a lying-down death pose (humanoids kneel), a fully airborne frame, a frame facing left.
```
OUTPUT: may_golem_set_hit_death_2f_chroma.png · 1448×724 · 2 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\

## 7. Tự kiểm và giao

### 7.1 Khối copy-dán giao việc cho agent vẽ (giai đoạn 1 — 3 ảnh mẫu, DỪNG chờ duyệt)

```
Bạn là hoạ sĩ tạo ảnh cho game 2D Naruto. Thư mục brief: D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the
ĐỌC TRƯỚC (theo thứ tự): _KHUON-DOT-4.md (§A, §B.4, §C) → 03-quai.md (mục 3, 4) → mở anh-tham-chieu\THUOC-NHAN-VAT.png
→ mở D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\wolf_idle_2f_chroma.png và master_soi_rung_x4_chroma.png (mẫu dải quái đã lắp)
→ mở D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-quai-soi-rung.png (mức chi tiết) và XAU-boss-zabuza.png (phản diện).

VIỆC GIAI ĐOẠN 1: tạo ĐÚNG 3 file sau bằng prompt trong 03-quai.md mục 6 (hoặc 03-quai.json anh[].prompt_day_du):
  1. master_lon_rung_x4_chroma.png     1254×1254 (thiết kế con lợn rừng)
  2. lon_rung_run_4f_chroma.png         2896×724  (dải 4 khung khó nhất: cùng tỉ lệ, chân chạm đáy từng ô)
  3. nhen_doc_hit_death_2f_chroma.png   1448×724  (thú thấp bẹp: khung chết vẫn tựa đáy, không glow, KHÔNG tím)
MỖI LẦN TẠO gắn kèm: THUOC-NHAN-VAT.png + wolf_idle_2f_chroma.png + TOT-quai-soi-rung.png (và XAU-boss-zabuza.png làm phản diện).

LUẬT KHÔNG ĐƯỢC VI PHẠM: cỡ đúng từng pixel (n×724 × 724; master 1254×1254); nền tím RGB(246,4,242) PHẲNG (không gradient, không bóng, không sàn);
một ô một khung, hình không cắn sang ô kề, chừa ≥ 24 px tím hai bên trong ô; 10 khung cùng tỉ lệ với idle_0; chân/bụng chạm mép dưới ở MỌI khung;
quay PHẢI mọi khung; không màu tím hồng trên quái; không glow lan nền. Nếu công cụ xuất cỡ khác thì cắt/ghép lại đúng cỡ TRƯỚC khi lưu.

LƯU VÀO: D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\
Đặt đúng tên file như trên.

BÁO CÁO khi xong (nói thật, không suy đoán "đạt"): cỡ thực từng file; màu nền đo được ở 4 góc; khung nào chân không chạm đáy;
khung nào bạn KHÔNG CHẮC và vì sao. DỪNG LẠI sau 3 file này, chờ duyệt. Không tự làm tiếp.
```

Giai đoạn 2 (chỉ dán sau khi 3 mẫu được duyệt): "3 mẫu đã duyệt. Làm nốt các file P1 còn lại theo mục 5 (thứ tự master → idle → run → pounce → hit_death từng con), rồi P2 theo làng Cát → Sương Mù → Đá → Mây. Giữ nguyên mọi luật giai đoạn 1, lưu cùng thư mục, báo cáo trung thực từng file."

### 7.2 Người vẽ tự kiểm trước khi giao (lệnh có thật)

Bước 1 — kiểm nhanh cỡ, số ô, màu nền 4 góc, chân chạm đáy từng ô (lưu đoạn dưới thành `kiem-o-quai.py` cạnh ảnh rồi chạy `python kiem-o-quai.py`; dán nguyên kết quả vào báo cáo):

```python
# kiem-o-quai.py — chay trong thu muc anh (can Pillow). In: ten, cỡ, ô, 4 góc, đáy hộp từng ô.
import glob, re
from PIL import Image
TIM = (246, 4, 242)
for p in sorted(glob.glob('*.png')):
    im = Image.open(p).convert('RGB'); w, h = im.size
    m = re.search(r'_(\d+)f_', p); n = int(m.group(1)) if m else 1
    goc = [im.getpixel(q) for q in ((0, 0), (w - 1, 0), (0, h - 1), (w - 1, h - 1))]
    ok_o = (w % n == 0) and ((w // n) % 4 == 0)
    day = []
    for i in range(n):
        o = im.crop((i * w // n, 0, (i + 1) * w // n, h)); px = o.load(); ow, oh = o.size; d = -1
        for y in range(oh - 1, -1, -1):
            if any(max(abs(px[x, y][0] - TIM[0]), abs(px[x, y][1] - TIM[1]), abs(px[x, y][2] - TIM[2])) > 40 for x in range(0, ow, 4)):
                d = y; break
        day.append(d)
    print(p, w, 'x', h, '| o=', w // n, 'x', h, 'OK' if ok_o else 'SAI', '| 4 goc', goc, '| day hop tung o', day, '(chuan >=', h - 8, ')')
```

Bước 2 — review chuẩn của dự án (chạy từ `D:\NRO_GoiBanGiao`, cần Pillow + numpy):

```powershell
python -m tools.artqa.artqa review --art-root D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai --review D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\qa\review-03-quai.json --reviewer "<tên>"
```

Đọc kết quả là ĐẠT khi (ngưỡng `artqa.py:23-27, 203-269`): `raw_source_status = PASS` (viền ≥ 85 % tím, lệch ≤ 12) · `issues = []` (không ô rỗng, soft alpha ≤ 15 %, không 2 khung liền giống nhau, loop không nhảy) · trong `qa\crop-reports\*_crop.json` mỗi `content_boxes[i]`: đáy hộp `≥ 716` (chân chạm đáy ô 724, hở ≤ 8 px), mép trái ≥ 24 và mép phải ≤ 700 trong ô (không dính ô kề), và **chiều cao hộp các khung cùng dải lệch nhau ≤ 10 %** so với idle_0 (kiểm cùng cỡ — công cụ tự động cho quái *chưa có*, so tay bằng số trong crop report). Mở GIF trong `qa\previews` xem 100 ms/khung.

### 7.3 Phía dự án (sau khi mọi cửa xanh) — việc còn thiếu

1. Nới cổng `artqa gate` / `task6b_closure` khỏi 20 tên v1 (`artqa.py:31-52`, `task6b_closure.py:34-35, 947-948`) — **chưa có**.
2. Viết pipeline quái mới: một tỉ lệ chung + hộp từng khung + atlas nhiều hàng + bảng ánh xạ 12 frame (mục 3.4) + nhiều mob (`build.py:1064-1091`) — **chưa có**.
3. `catalog.json` thêm `mobs[]` id 1, 2, 3… → ký lại `asset-lock.json` (reviewer độc lập, mã ASSET-ARTIFACT-REVIEW-016, không tự sửa chữ ký) → `quai.yaml` thêm `mau_quai` + `sinh` cho map 6 → `python -m assetpack validate/build` (`_KHUON-DOT-4.md` §D.4) → tăng `versions.data`.
4. Kiểm bằng Unity MCP: sói đang đánh có vẽ khung không (mục 3.4) — quyết định bảng ánh xạ trước khi bind con thứ hai.

## 8. Chưa cần vẽ đợt này

- **Vẽ lại 4 dải sói** cho ô 724 đều nhau: chỉ làm nếu chủ dự án quyết đồng bộ khi đổi pipeline (prompt đã để sẵn ở mục 6, trạng thái da_co).
- **4 con còn lại mỗi làng** của brief đợt 3: Cát — giun cát, nhện cát, kền kền, lạc đà · Sương Mù — lươn điện, thuỷ quái xúc tu, quái sương, golem băng · Đá — dơi hang, gấu núi, kỳ đà, ưng núi · Mây — đại bàng sét, bò tót, sơn dương, tinh linh mây. Trong đó **loài bay / lơ lửng / tan thành sương** (chim ưng, kền kền, dơi, đại bàng, tinh linh mây, quái sương, giun cát trồi lên) **hoãn có lý do**: xung đột luật chân chạm đáy + cấm glow, và nhánh bay `type 4/5` của client (`Mob.cs:1233-1268`) chưa kiểm.
- **Bia gỗ tập luyện** (`quai_bia_go`, lô 343): quái đứng yên cho sân Học viện — cần máy chủ gửi `isDontMove` (hiện `ghiQuai` luôn ghi false, `quai.go:154`) — chờ máy chủ.
- **Biến thể tinh anh** (đổi màu, `cap_tinh_anh`), **boss dạng thú** (vĩ thú, rắn khổng lồ) — cần code boss dạng Mob (KH6b ngoài phạm vi, P3-09).
- **Icon/chân dung quái** cho bảng nhiệm vụ — client chưa có chỗ dùng.

---
*Sinh bởi script từ một nguồn dữ liệu; `03-quai.json` chứa đúng các prompt ở mục 6 (`anh[].prompt_day_du`), cùng `kieu`, `so_khung`, `o_kich_thuoc`, `luu_vao`, `uu_tien`, `trang_thai`, `thu_tu_khung`, `cao_muc_tieu_x4`, `ma_quai`, `cap`, `noi_sinh`, `dung_tham_chieu`.*
