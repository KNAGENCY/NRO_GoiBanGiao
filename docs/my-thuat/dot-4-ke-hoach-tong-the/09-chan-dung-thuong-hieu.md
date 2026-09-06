# CHÂN DUNG · THƯƠNG HIỆU — Đợt 4 (brief chờ duyệt 3 ảnh mẫu · viết 06/09/2026)

> Đọc `_KHUON-DOT-4.md` và `../dot-3-toan-bo/_KHUON.md` trước. **LƯU VÀO:** `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\`
>
> Hạng mục 09 của đợt 4. Hai nửa khác hẳn nhau về kỹ thuật nhưng cùng một câu hỏi "người chơi nhìn thấy MẶT ai, và thấy THƯƠNG HIỆU nào": (A) **chân dung** cho ô HUD và khung hội thoại — dải chroma 1 ô theo đúng khuôn `iruka_portrait_chroma.png` đang chạy; (B) **thương hiệu** — icon app, màn chờ, ảnh bìa cửa hàng, banner đăng nhập theo mùa — ảnh tĩnh ĐỤC, không đi qua đường ống sprite. Mọi con số dưới đây đọc từ code, kèm `file:dòng`; chỗ nào chưa chắc ghi rõ *chưa chắc*. Tên NPC/boss lấy từ kho nội dung máy chủ và `kich-thuoc-anh.json` vì hai brief `02-boss.md`, `03-npc.md` **chưa được viết** tại thời điểm này — brief nào viết sau phải dùng đúng các tên ở mục 5.

Tổng: **44 ảnh** — P1 **2** · P2 **36** · P3 **5** · đã có **1**. Thư mục lưu chưa tồn tại — agent vẽ tự tạo.

---

## 1. Game dùng ở đâu, vì sao cần

**Chân dung (`cd_` / `_portrait_chroma`) hiện lên ở ba chỗ, cùng một kiểu ảnh `Small<id>`:**

| Chỗ | Client vẽ thế nào | Cỡ thật trên màn | Ai quyết id ảnh |
|---|---|---|---|
| Nóc **khung hội thoại NPC** | `ChatPopup.cs:316` vẽ `c.avatar` neo **đáy-trái** tại `(cx+14, cy)` — tức mép DƯỚI chân dung đặt lên mép TRÊN khung thoại, nhún 1 px theo nhịp | 200×200 x4 = **50×50 đơn vị** (chuẩn Iruka `Small27078`) | Máy chủ: gói NPC (`server/internal/game/npc.go:53,68` ghi `Portrait`), gói thoại 38/32 đọc `avatar` ở đuôi gói (`Controller.cs:3618`, `:3659`), gói -70 "thông báo lớn" (`Controller.cs:3562-3568`) |
| **Ô xanh HUD** góc trên-trái | `GameScr.cs:5588-5590` vẽ `c.avatarz()` **cỡ gốc** tại `(6,6)` — không co giãn | ô cao ~28 đơn vị → ảnh phải **≈26×26 x1 = 104×104 x4** | Bảng head→avatar trong ME_INFO (`thegioi.go:141-143`), hằng `AvatarNhanVat = 27084` (`thegioi.go:165-168`) — **một id cho mọi làng** |
| Góc trên **bảng cửa hàng / menu** | `Panel.cs:10072-10270` vẽ `avatarz()` hoặc `npcFocus.avatar` tại `(X+25, 50)` neo 33 | cùng sprite với hai chỗ trên | như trên |

Vì sao cần vẽ: (1) chân dung người chơi hôm nay (`Small27084`) chỉ là **cái đầu cắt từ sprite 27041** thu về 26 đơn vị (`tools/map-lang-la/tao-chan-dung-nguoi-choi.py:1-9`) — nhìn được nhưng không phải chân dung, và mờ vì cắt từ ảnh nhỏ; (2) NPC mới nào cũng cần một chân dung đi kèm master (B.3 `_KHUON-DOT-4.md`), mà 15 chân dung `cd_*` lô 343 vẽ **pixel-art, quay TRÁI** — lệch với Iruka đang lắp; (3) boss có thoại (`boss.yaml` `thoai_hien/danh/chet`) — nhưng **đọc code thì thoại boss là bong bóng trên đầu**, không có chân dung (xem mục 3, C7 đã trả lời).

**Thương hiệu — cái người chơi thấy TRƯỚC khi vào game:**

- **Icon app đang là Dragon Ball.** `D:\NRO_Project\Assets\Resources\appicon.png` (479×480) là rồng xanh quấn ngọc 4 sao; `ProjectSettings\ProjectSettings.asset:292-335` trỏ mọi nền tảng (mặc định, iPhone, iPad) về đúng guid `b611de4c…` của file này; `productName: Ngoc Rong Kame` (`:16`). Đây là dấu vết game cũ lộ ngay ở màn hình điện thoại — vì thế icon xếp **P1** dù I2 xếp thương hiệu chung là P3 (một ảnh, ảnh hưởng lớn).
- **Logo trong game** là wordmark chữ NARUTO 599×136 (`Assets\Logo\1.png` → `res\x4\Logo\1.png`, nạp ở `ModFunc.cs:2624`, vẽ ở `LoginScr.cs:655-667`); asset-lock lại ký `ui_logo_game_x4` 512×256 là **emblem lá-lửa** (lệch #2, `00-KIEM-KE-DA-CO.md` §4). Brief này **không vẽ logo mới** — dùng emblem `TOT-naruto-logo.png` cho icon/splash, và **chừa chỗ dán wordmark** thay vì để AI vẽ chữ (AI vẽ chữ hay sai).
- **Màn chờ**: Unity còn bật splash mặc định (`m_ShowUnitySplashScreen: 1`, `m_SplashScreenLogos: []`, `androidSplashScreen: {fileID: 0}` — `ProjectSettings.asset:20,42,271`) → chưa có ảnh chờ nào của game. ⚠️ Đừng nhầm với `/e/sp0..2.png` ở `GameScr.cs:5950` — đó là **hiệu ứng toé nước** (`startSplash(x, y, dir)`), không phải màn chờ.
- **Ảnh bìa cửa hàng** 1024×500 và **banner đăng nhập theo mùa**: chưa có; client hôm nay không nạp nền đăng nhập (`LoginScr` chỉ vẽ logo) nên banner là P3.

---

## 2. Đã có gì

| Tên | Nơi lưu | Dùng được không | Vì sao |
|---|---|---|---|
| `iruka_portrait_chroma.png` 1254×1254 → `Small27078` 200×200 | `anh-nguon\vertical-slice-v1\` (+ `qa\alpha`, `qa\contact-sheets`) | ✅ **Dùng thẳng — là CHUẨN** | QA PASS (`review\vertical-slice-v1.json`: content_box [105,58,1045,1223], soft alpha 0,5 %); bán thân, 3/4 quay PHẢI, vai cắt phẳng chạm đáy ô. Đang lắp (`catalog.json` npcs[0].portrait) |
| `Small27084` (HUD người chơi) 26 đơn vị | `D:\NRO_Project\Assets\Resources\res\x1..x4\smallimage\` | ⚠️ **Thay** | Là đầu 27041 cắt + thu nhỏ (`tao-chan-dung-nguoi-choi.py`), không có trong lock (id ký tới 27082) — cần chân dung vẽ thật (A1) |
| `Small27083` 128×128 | client `smallimage` (22/08) | ❓ chưa xác minh | Không có trong lock, không khớp file nào trong repo (`00-KIEM-KE-DA-CO.md` §3.2, §4 #6) — đừng cấp id này cho ảnh mới |
| 15 `cd_*_x4.png` 800×800 (chuẩn 200×200) | `anh-nguon\343-anh\07-npc\` (bản đúng cỡ: `anh-game\x4\`) | ⚠️ **Chỉ làm tham chiếu nội dung** | Pixel-art (phong cách khoá nay là hand-painted), **quay TRÁI** (client không lật chân dung → lệch với Iruka), nền alpha AI có viền — không đi qua `artqa` được. `cd_iruka` đã thay bằng chroma |
| 15 `npc_*_x4.png` 384×512 | `anh-nguon\343-anh\07-npc\` | ⚠️ tham chiếu nội dung | Nhân vật đứng nguyên con format cũ — dùng để biết NPC đó tóc/đồ gì khi vẽ chân dung + master (brief 03) |
| 12 `cd_{cat,smu,da,may}_*` | `kich-thuoc-anh.json` (tên) | — | **Chưa có ảnh** (`00-KIEM-KE-DA-CO.md` §5: cd 27 tên, có 15) |
| Chân dung boss | — | — | Chưa có tên, chưa có ảnh; boss hôm nay mặc y bộ người chơi / Iruka (`boss.yaml:18-31`) |
| Wordmark `Logo\1.png` 599×136 | `lang-la-client-assets\Logo\1.png` (sao lưu) · nguồn `D:\NRO_Project\Assets\Logo\1.png` | ✅ dùng để **dán** vào splash/bìa | Đang chạy trong client; sửa phải sửa NGUỒN `Assets\Logo\1.png` (bộ nhớ `logo-resource-sync`) |
| Emblem `TOT-naruto-logo.png` 512×256 (= `ui_logo_game_x4` trong lock) | `dot-3-toan-bo\anh-tham-chieu\` · `my thuat\_v2-bung\01-giao-dien\ui_logo_game_x4.png` | ✅ **nguồn hình cho icon/splash** | Đã duyệt; bảng màu thương hiệu (cam-đỏ, xanh ngọc, navy) |
| `appicon.png` 479×480 | `D:\NRO_Project\Assets\Resources\` (chỉ đọc) | ❌ **Phải thay** | Rồng + ngọc 4 sao Dragon Ball — mẫu "xấu" về thương hiệu, không phải mẫu vẽ |
| `ui_nen_dangnhap_x4` · `ui_nen_chonlang_x4` 1024×576 | `anh-nguon\343-anh\01-giao-dien\` | ⚠️ bố cục gốc cho banner mùa | Pixel-art lô 343; client chưa có chỗ nạp nền đăng nhập |
| `THUOC-NHAN-VAT.png` 2400×1060 | `dot-4-ke-hoach-tong-the\anh-tham-chieu\` | ✅ thước bắt buộc | Cỡ đầu idle_0: cao mặt 96–99, tóc→cằm 266–271 px trong ô 724 |

---

## 3. Ràng buộc kỹ thuật riêng loại này

### 3.1 Chân dung = dải chroma **1 ô**, khuôn Iruka

| Mục | Bắt buộc | Vì sao / nguồn |
|---|---|---|
| Kiểu file | `<tên>_portrait_chroma.png`, **1 ô 1254×1254**, nền tím phẳng RGB(246, 4, 242) | Máy tách nền bằng flood-fill từ mép (`tools/artqa/artqa.py:112-149`); tên không có `_<n>f` → 1 ô (`artqa.py:450-452`). Ảnh alpha AI xuất thường dính viền/bóng → hộp đo sai |
| Máy làm gì | `artqa review` đo `content_boxes` → `task6b_closure.py:1004-1019` `_single_frame_spec(... target_size=(200,200))` → `fit_to_canvas` (`:632-652`: cắt sát alpha, co giữ tỉ lệ, **neo giữa-dưới**) → `Small<id>` 200×200 ×4 zoom | Vì neo giữa-dưới nên phần vai cắt phẳng sẽ nằm sát đáy ô 200 — đúng chỗ khung thoại cần |
| Khung hình | **Bán thân**: đầu + cổ + vai + ngực trên. Hộp nội dung cao **85–92 %** ô; đỉnh tóc cách mép trên ≥ 40 px; mép trái/phải chừa ≥ 60 px tím; **đáy vai là một đường cắt ngang sạch CHẠM mép dưới ô** | ChatPopup vẽ neo đáy-trái đặt lên nóc khung (`ChatPopup.cs:316`): đáy ảnh chính là đường khung — vai phải kết thúc gọn, không vát mờ. Iruka: hộp [105,58,1045,1223]/1254 |
| Cỡ mặt | Mặt chiếm **45–55 %** chiều cao ô | Ảnh hiện ở 50×50 (khung thoại) và 26×26 đơn vị (HUD) — mặt nhỏ hơn là thành chấm |
| Hướng | **3/4 quay PHẢI**, đầu thẳng, mắt nhìn trước-phải | Client **không lật** chân dung (`SmallImage.drawSmallImage` vẽ thẳng) → phải thống nhất tay với Iruka đang lắp; lô 343 quay trái là lệch |
| Đồng bộ nhân vật | Cùng mặt/tóc/băng trán/màu áo với **dải hoặc master của chính nhân vật đó** — kèm ảnh đó mỗi lần gen | Chân dung đứng cạnh sprite chạy của cùng người; khác mặt là người chơi thấy "hai người" |
| Cấm màu | **Không tím hồng/magenta trong hình** (R−G ≥ 110 và B−G ≥ 110 bị coi là nền, đảo tím ≥ ~0,04 % bị xoá — `artqa.py:60-62`, `task6b_closure.py:499-543`). **Làng Mây** dùng tím **xanh** (indigo/lavender, R−G < 100) | Áo tím hồng là mất mảng áo. Bông tai nhỏ thì sống |
| Cấm hiệu ứng | Không glow/khói mờ lan ra nền (soft alpha > 15 % = FAIL, `artqa.py:25, 245-269`), không bóng đổ, không viền sáng/tối, không gradient nền (`spread > 12` = nguồn FAIL, `artqa.py:23-24`) | v1 0/20 dải đạt `raw_source_status` vì gradient nền — đợt 4 phải đạt |
| Không vẽ | Chữ, số, logo, **khung trang trí** (khung avatar là UI riêng `ui_khung_avatar`), tay, vũ khí che mặt, nhân vật thứ hai | Khung/chữ do client vẽ |
| HUD người chơi | Từ **cùng một ảnh nguồn** A1, dự án cắt thêm bản sát mặt **104×104 x4 (26 x1)** thay nội dung `Small27084` | HUD vẽ cỡ gốc tại (6,6) trong ô cao ~28 đơn vị (`GameScr.cs:5590`, `tao-chan-dung-nguoi-choi.py:17`) — đặt 200×200 vào là tràn ô |
| Id ảnh | Dự án cấp khi lắp, **nối tiếp ≥ 27085** (27083 lạ, 27084 đang là HUD) — append-only, không dùng lại | `mapping-id.json` `phase.reserved_sprite_id_range` [27041, 27999]; client cache theo idImage (`BAT-DAU-TU-DAY.md:441-444`) |
| **Chưa có, cần sửa (phía dự án, không phải việc người vẽ)** | (a) `_single_frame_spec` khoá cứng `27078`/`iruka_portrait_chroma.png` (`task6b_closure.py:1004-1010`) → tham số hoá danh sách chân dung; (b) `fit_to_canvas` thu 1254→200 bằng **NEAREST** (`:646`) → đổi LANCZOS cho chân dung kẻo răng cưa; (c) cổng `artqa gate`/`generate_vertical_slice` đòi đúng 20 tên v1 (`artqa.py:31-52`) → nới; (d) muốn chân dung **boss** hiện thì máy chủ phải gửi thoại boss qua gói có `avatar` (38/32/-70) thay gói 44 | |

**C7 (00-KIEM-KE-GAME-CAN.md §7) — đã trả lời bằng code:** gói 44 chỉ có 2 trường `int id + UTF nội dung` (`server/internal/game/chat.go:1-2, 49-59, 70-87`; `Controller.cs:4567-4572` → `Char.addInfo` = bong bóng trên đầu). Boss chào/đánh/chết đều đi gói này (`khuvuc.go:2399-2411` `noiThoaiBoss` → `goiChat44`). **Không có đường nào vẽ chân dung boss hôm nay** → 3 chân dung boss xếp **P2 đề xuất**, không phải P1 như I2 tạm ghi.

### 3.2 Bảy luật 05/09 áp cho chân dung (thước: `anh-tham-chieu\THUOC-NHAN-VAT.png`)

| # | Luật gốc (nhân vật) | Với chân dung nghĩa là |
|---|---|---|
| 1 | Mọi khung cùng **cỡ đầu** với idle_0 (cao mặt 96–99, tóc→cằm 266–271 px trong ô 724) | Chân dung không ghép xương nên không đo bằng px này, nhưng **cùng khuôn mặt, cùng tỉ lệ đầu/cổ/vai** với master của nhân vật — kèm `THUOC-NHAN-VAT.png` + ảnh master khi gen |
| 2 | **Bàn chân chạm mép dưới ô** ở mọi khung | Không có chân → **đường cắt vai chạm mép dưới ô** (máy neo giữa-dưới) |
| 3 | Đứng thẳng, không khom | **Đầu thẳng**, không ngửa/cúi/nghiêng quá 10° |
| 4 | Nền tím phẳng, không bóng/viền/gradient | Y nguyên |
| 5 | Một ô một khung, ô đều, thứ tự cố định | **Đúng 1 ô**, tên không có số khung |
| 6 | Đầu 3/4 quay PHẢI | Y nguyên — và với chân dung còn quan trọng hơn vì client **không lật** |
| 7 | Dải hit không vẽ nằm | Không áp (không có dải) — nhưng đừng vẽ chân dung "gục/ngã" cho boss, chỉ biểu cảm |

Câu bắt buộc nằm trong mọi prompt nhân vật: *"same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242)"* — mục 6 có kèm chú giải "feet" = đường cắt vai.

### 3.3 Thương hiệu = ảnh tĩnh **ĐỤC**, không đi qua đường ống sprite

| Mục | Bắt buộc | Vì sao |
|---|---|---|
| Định dạng | PNG **RGBA, alpha toàn 255** (đục) — trừ `ui_icon_app_hinh_x4` (lớp hình adaptive, alpha thật) | App Store từ chối icon có alpha; `KIEM-TRA-ANH.py:117-118` lại báo thiếu alpha nếu file RGB → giao RGBA đục là qua cả hai |
| Icon | Vẽ **1024×1024**, giao nguyên; dự án thu về 512 (Google Play, Unity). Emblem nằm trong **vòng an toàn 66 %** giữa; **không bo góc, không viền, không chữ** | Android adaptive icon cắt tròn/vuông bo tuỳ máy; iOS tự bo góc; chữ ở 48 px không đọc được |
| Splash | 1080×1920 và 1920×1080, mọi thứ quan trọng trong **80 % giữa**; **chừa vùng trống cho wordmark**, không vẽ chữ | Tai thỏ/góc bo; dự án dán `Logo\1.png` — AI vẽ chữ hay sai nét |
| Bìa cửa hàng | 1024×500, lề 10 %, không chữ (Google Play khuyến nghị) | Feature graphic bị cắt/đè tiêu đề tuỳ chỗ hiển thị |
| Banner mùa | 1024×576 (tỉ lệ `ui_nen_dangnhap`), đục | Client chưa nạp nền đăng nhập — P3 |
| Nhân vật trên bìa | **Đúng mặt** người chơi (`player_idle_3f_chroma`) và Iruka (`master_iruka`) — kèm ảnh khi gen | Bìa cửa hàng mà nhân vật khác trong game là lừa người tải |
| Kiểm | `KIEM-TRA-ANH.py <thư mục> nguon` (tên phải có trong `kich-thuoc-anh.json` trước — `KIEM-TRA-ANH.py:97-100`); **không** chạy `TAO-4-MUC-ZOOM.py` cho nhóm này | Không phải asset 4 zoom |
| Lắp | Unity Player Settings → Icon / Splash (thay `appicon.png`, đặt `androidSplashScreen`) — **ngoài phạm vi brief này**, `D:\NRO_Project` chỉ đọc | Việc của người lắp client, ghi ở §7 |
| Tên mới | Nối CUỐI `kich-thuoc-anh.json`, tiền tố `ui_` (đã có `ui_logo_game`, `ui_nen_dangnhap`), hậu tố `_x4` chỉ là khoá tra | `_KHUON.md` §5 append-only |
| Rủi ro ngoài mỹ thuật (*chưa chắc*) | Wordmark **NARUTO** là nhãn hiệu của bên khác; ảnh bìa/icon lên cửa hàng ứng dụng có thể bị gỡ | Chủ dự án quyết trước khi dùng bìa cửa hàng; brief chỉ ghi nhận |

---

## 4. Phong cách · ảnh tham chiếu · subject block tiếng Anh

**Phong cách khoá** (dán vào mọi prompt, bản gốc): `detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, NO flat vector shapes, NO background scenery, NO drop shadow on ground`

- Với **chân dung** đổi cụm `side/3-4 view facing RIGHT, full body` → `3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body)`; với **thương hiệu** bỏ cụm hướng nhìn/toàn thân. Mục 6 đã thay sẵn.
- Màu hệ làng (`_KHUON.md` §3): Lá cam-đỏ · Cát vàng cát-lục nhạt · Sương Mù lam-xanh dương · Đá nâu-xám · Mây **tím xanh**-trắng-vàng chớp (tím phải ngả xanh, xem cấm màu §3.1).

**Ảnh tham chiếu — mỗi lần gen chân dung kèm 3 ảnh:** (1) `iruka_portrait_chroma.png` (khuôn chân dung ĐÃ LẮP — bố cục, hướng, đường cắt vai); (2) ảnh master/dải của **chính nhân vật đó** (đồng bộ mặt); (3) `THUOC-NHAN-VAT.png`. Kèm `XAU-npc.png`/`XAU-boss-zabuza.png` làm phản diện. Với thương hiệu: `TOT-naruto-logo.png` (emblem) + `concept-map1-langla.png` (tông nền) + `logo_wordmark_client` (để biết vùng chữ sẽ to bao nhiêu).

| Khoá | Đường dẫn | Dùng để |
|---|---|---|
| `TOT-naruto-logo` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-naruto-logo.png` | Emblem la-lua xoay 512x256 da duyet — bang mau thuong hieu (cam-do, xanh ngoc, navy). Gan kem MOI lan gen. |
| `TOT-chitiet-nhanvat` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-chitiet-nhanvat.png` | Muc chi tiet va do bong can dat (chi tham chieu KIEU VE, khong phai noi dung). |
| `XAU-boss-zabuza` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\XAU-boss-zabuza.png` | Mau XAU: hinh khoi phang vai mau — tuyet doi tranh. |
| `XAU-npc` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\XAU-npc.png` | Mau XAU: chan dung NPC phang 200x200 lo cu — tranh. |
| `TOT-npc-ong-tazuna` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-npc-ong-tazuna.png` | Mau TOT ve vai co nep gap, phu kien, da nguoi gia. |
| `THUOC-NHAN-VAT` | `D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\anh-tham-chieu\THUOC-NHAN-VAT.png` | Thuoc ti le nhan vat 05/09: cao mat 96-99 px, toc->cam 266-271 px trong o 724. Chan dung phai cung khuon mat/ti le dau nay. |
| `iruka_portrait_chroma` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\iruka_portrait_chroma.png` | CHAN DUNG MAU DA LAP TRONG GAME (Small27078): ban than, 3/4 quay PHAI, vai cat phang cham day o, nen tim. Moi chan dung moi phai cung khuon nay. |
| `iruka_portrait_contact` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\qa\contact-sheets\iruka_portrait_light-dark.png` | Ban da tach nen cua chan dung Iruka tren nen sang/toi — de nhin ro khung hinh. |
| `player_idle_3f_chroma` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\player_idle_3f_chroma.png` | Nhan vat nguoi choi Lang La goc (idle) — chan dung nguoi choi phai cung mat/toc/bang tran/ao voi khung idle_0 nay. |
| `master_hoc_vien_x4_chroma` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\master_hoc_vien_x4_chroma.png` | Ban dung nguyen con cua nguoi choi hoc vien — dong bo trang phuc. |
| `master_iruka_x4_chroma` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\master_iruka_x4_chroma.png` | Master Iruka dang lap — de boss pha 2 KHONG bi giong Iruka. |
| `logo_wordmark_client` | `D:\NRO_GoiBanGiao\docs\my-thuat\lang-la-client-assets\Logo\1.png` | Wordmark NARUTO 599x136 dang chay trong client (res/x4/Logo/1.png). Splash/bia chua vung trong de DAN wordmark nay, KHONG ve chu. |
| `appicon_cu_dragonball` | `D:\NRO_Project\Assets\Resources\appicon.png` | Icon app HIEN TAI cua client (479x480): rong xanh quan ngoc 4 sao — game cu Dragon Ball. La thu PHAI THAY, khong phai mau. |
| `concept-map1-langla` | `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\concept-map1-langla.png` | Concept Lang La da duyet — tong mau/net ve cho phong nen splash va bia cua hang. |
| `ui_nen_dangnhap_343` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\01-giao-dien\ui_nen_dangnhap_x4.png` | Nen dang nhap 1024x576 lo 343 (lang ninja hoang hon, trang lon) — bo cuc goc cho banner theo mua; ve lai theo phong cach hand-painted. |
| `wolf_idle_2f_chroma` | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\wolf_idle_2f_chroma.png` | Soi Rung dang lap — neu dua vao bia cua hang thi lay dang nay. |
| `npc343_<tên>` · `cd343_<tên>` (15 cặp) | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\07-npc\npc_<tên>_x4.png` · `\cd_<tên>_x4.png` | Chỉ lấy NỘI DUNG nhân vật lô 343 (tóc, đồ, phụ kiện); bỏ phong cách pixel-art và hướng quay trái |

**Subject block chân dung (EN, dùng chung — mục 6 nhúng sẵn):**

```text
bust portrait of <character>, head and shoulders, 3/4 view facing RIGHT, head upright, eyes looking forward-right, same face / hair / forehead protector / outfit colors as this character's reference sprite, shoulders ending in a clean straight horizontal cut exactly on the bottom edge, one single 1254x1254 cell on a perfectly flat magenta background RGB(246,4,242), hand-painted cel-shading, clean dark outlines, 40+ colors, no text, no frame, no glow, no shadow, no magenta inside the character
```

**Subject block thương hiệu (EN):**

```text
hand-painted cel-shaded game key art in the approved Leaf-Fire emblem palette (orange-red flame, teal accents, deep navy), fully opaque PNG, no text or lettering anywhere (the real wordmark is composited later), important content inside the safe zone, Leaf Village dusk skyline motif (cliff with blank carved shapes, red torii gate, tiled rooftops), no real-world logos
```

---

## 5. DANH SÁCH ẢNH

Ưu tiên theo `00-KIEM-KE-GAME-CAN.md` (P1 Làng Lá chơi trọn · P2 4 làng và đề xuất chờ chốt · P3 đẹp thêm). **Lệch có chủ đích so với I2:** icon app lên P1 (icon đang là Dragon Ball); chân dung boss xuống P2 (gói 44 không có chân dung). Thứ tự làm trong P2: boss + 14 NPC Làng Lá (khi `02-boss.md`/KH9 chốt) → splash/bìa → 4 làng.

| STT | Tên file | Kiểu | Cỡ | Làng | Mô tả | Ưu tiên | Trạng thái |
|---|---|---|---|---|---|---|---|
| 1 | `nguoi_choi_la_portrait_chroma.png` → `cd_nguoi_choi_la_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | Chân dung người chơi Làng Lá (học viên) — CÙNG KHUÔN MẶT với dải `player_idle_3f_chroma.png` đang chạy trong game: tóc đen-nâu dựng, mắt nâu, băng trán vải navy gắn tấm thép khắc lá xoáy, áo giáp vest xanh ô-liu, cổ áo navy, thoáng nút khăn cam ở mép dưới. Máy sẽ cắt ra HAI bản: 104×104 (26 đơn vị) thay nội dung `Small27084` cho ô xanh HUD và bảng cửa hàng; 200×200 (id mới) cho khung thoại khi máy chủ cần. Hiện 27084 chỉ là cái đầu cắt từ sprite, không phải chân dung vẽ. | **P1** | can_ve |
| 2 | `nguoi_choi_cat_portrait_chroma.png` → `cd_nguoi_choi_cat_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Cát | Chân dung người chơi mặc định Làng Cát (đề xuất). Cùng tuổi, cùng tỉ lệ đầu với học viên Làng Lá, đổi tóc/da/băng trán/màu cổ áo theo làng. CHỜ: brief `01-nhan-vat-nguoi-choi.md` chốt dải đầu theo làng và máy chủ gửi đầu/chân dung theo làng (`hinhdang.go:56-58` hiện bỏ qua `lang`; `thegioi.go:165-168` hằng `AvatarNhanVat` = 27084 dùng cho mọi làng). Vẽ SAU khi 01 chốt, dùng đúng mô tả đầu của 01. | **P2** | can_ve |
| 3 | `nguoi_choi_smu_portrait_chroma.png` → `cd_nguoi_choi_smu_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Sương Mù | Chân dung người chơi mặc định Làng Sương Mù (đề xuất). Cùng tuổi, cùng tỉ lệ đầu với học viên Làng Lá, đổi tóc/da/băng trán/màu cổ áo theo làng. CHỜ: brief `01-nhan-vat-nguoi-choi.md` chốt dải đầu theo làng và máy chủ gửi đầu/chân dung theo làng (`hinhdang.go:56-58` hiện bỏ qua `lang`; `thegioi.go:165-168` hằng `AvatarNhanVat` = 27084 dùng cho mọi làng). Vẽ SAU khi 01 chốt, dùng đúng mô tả đầu của 01. | **P2** | can_ve |
| 4 | `nguoi_choi_da_portrait_chroma.png` → `cd_nguoi_choi_da_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Đá | Chân dung người chơi mặc định Làng Đá (đề xuất). Cùng tuổi, cùng tỉ lệ đầu với học viên Làng Lá, đổi tóc/da/băng trán/màu cổ áo theo làng. CHỜ: brief `01-nhan-vat-nguoi-choi.md` chốt dải đầu theo làng và máy chủ gửi đầu/chân dung theo làng (`hinhdang.go:56-58` hiện bỏ qua `lang`; `thegioi.go:165-168` hằng `AvatarNhanVat` = 27084 dùng cho mọi làng). Vẽ SAU khi 01 chốt, dùng đúng mô tả đầu của 01. | **P2** | can_ve |
| 5 | `nguoi_choi_may_portrait_chroma.png` → `cd_nguoi_choi_may_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Mây | Chân dung người chơi mặc định Làng Mây (đề xuất). Cùng tuổi, cùng tỉ lệ đầu với học viên Làng Lá, đổi tóc/da/băng trán/màu cổ áo theo làng. CHỜ: brief `01-nhan-vat-nguoi-choi.md` chốt dải đầu theo làng và máy chủ gửi đầu/chân dung theo làng (`hinhdang.go:56-58` hiện bỏ qua `lang`; `thegioi.go:165-168` hằng `AvatarNhanVat` = 27084 dùng cho mọi làng). Vẽ SAU khi 01 chốt, dùng đúng mô tả đầu của 01. | **P2** | can_ve |
| 6 | `nhan_gia_mat_na_p1_portrait_chroma.png` → `cd_boss_nhan_gia_mat_na_p1_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | Chân dung boss lớn **Nhân Gia Mặt Nạ — pha 1 "Kẻ Lạ Mặt"** (`boss.yaml` ma `nhan-gia-mat-na`, pha 1, outfit học viên 30/57/58): mặc đúng bộ học viên như người chơi (vest ô-liu, cổ navy) nhưng mặt che sau mặt nạ sứ trắng trơn, hai khe mắt hẹp, một vết nứt đỏ mảnh chéo má, tóc đen dài rối lòi ra sau mặt nạ. Cảm giác "đứng giữa đám học viên mà không ai ngoảnh lại". Boss thoại qua gói 44 = bong bóng trên đầu (server/internal/game/chat.go:1-2,70-87; Controller.cs:4567-4572) — KHÔNG có trường chân dung. Chân dung boss chỉ hiện khi máy chủ dùng gói 38/32 (Controller.cs:3618, 3659) hoặc -70 (Controller.cs:3562-3568) có `avatar`. Đề xuất; đồng bộ với dải ĐẦU boss trong 02-boss.md (chưa viết) — bên nào vẽ trước, bên sau theo. | **P2** | can_ve |
| 7 | `nhan_gia_mat_na_p2_portrait_chroma.png` → `cd_boss_nhan_gia_mat_na_p2_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | Chân dung boss lớn **Nhân Gia Mặt Nạ — pha 2 "Mặt Nạ Vỡ"** (chân thân, hôm nay đang mặc bộ Iruka 2189/2190/2191 nên KHÔNG được giống Iruka): đàn ông ~30, mặt gầy sắc, tóc đen dài có vệt bạc, sẹo chéo má trái, băng trán Làng Lá bị **gạch ngang xoá biểu tượng lá** (kẻ đào tẩu), nửa mặt nạ vỡ còn treo bên hàm, cổ áo choàng đỏ thẫm-đen lót cam than hồng (hệ Hoả). Thoại pha 2: "Ta từng đeo băng trán của làng này". Boss thoại qua gói 44 = bong bóng trên đầu (server/internal/game/chat.go:1-2,70-87; Controller.cs:4567-4572) — KHÔNG có trường chân dung. Chân dung boss chỉ hiện khi máy chủ dùng gói 38/32 (Controller.cs:3618, 3659) hoặc -70 (Controller.cs:3562-3568) có `avatar`. Đề xuất; đồng bộ với dải ĐẦU boss trong 02-boss.md (chưa viết) — bên nào vẽ trước, bên sau theo. | **P2** | can_ve |
| 8 | `hoc_vien_dao_tau_portrait_chroma.png` → `cd_boss_hoc_vien_dao_tau_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Sương Mù | Chân dung boss lẻ **Học Viên Đào Tẩu** (`boss.yaml` ma `hoc-vien-dao-tau`, làng 2 Sương Mù, đứng ở map 1, outfit học viên 30/57/58): thiếu niên ~14 bỏ trốn giữa kỳ thi Sương Mù — da tái, tóc xanh đen ướt rũ che một mắt, quầng mắt thâm, băng trán Sương Mù (4 vạch sóng) đeo lệch lỏng ở cổ, bộ học viên (vest ô-liu, cổ navy) bẩn sờn, miếng băng dán trên má. Boss thoại qua gói 44 = bong bóng trên đầu (server/internal/game/chat.go:1-2,70-87; Controller.cs:4567-4572) — KHÔNG có trường chân dung. Chân dung boss chỉ hiện khi máy chủ dùng gói 38/32 (Controller.cs:3618, 3659) hoặc -70 (Controller.cs:3562-3568) có `avatar`. Đề xuất; đồng bộ với dải ĐẦU boss trong 02-boss.md (chưa viết) — bên nào vẽ trước, bên sau theo. | **P2** | can_ve |
| 9 | `iruka_portrait_chroma.png` → `cd_iruka_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | NPC Iruka — **ĐÃ CÓ, KHÔNG VẼ**: `anh-nguon\vertical-slice-v1\iruka_portrait_chroma.png` đã QA PASS, đang lắp là `Small27078` (catalog.json npcs[0].portrait). Đây là CHUẨN cho mọi chân dung mới. Prompt dưới chỉ để dùng nếu chủ dự án muốn vẽ lại đồng bộ sau. | **P3** | da_co |
| 10 | `hokage3_portrait_chroma.png` → `cd_hokage3_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | Chân dung NPC **Hokage đệ Tam** (đề xuất Làng Lá, chờ KH9 chốt danh sách NPC — máy chủ hôm nay chỉ có Iruka). Ông già thủ lĩnh làng: mặt nhăn hiền, râu cằm bạc, nón chóp Kage đỏ-trắng (KHÔNG chữ), áo choàng trắng viền đỏ, tẩu dài, sợi khói mảnh. Bản cũ `cd_hokage3_x4.png` trong lô 343 là pixel-art 800×800 và **quay TRÁI** — lệch với Iruka đang lắp, nên vẽ lại theo khuôn chroma; ảnh cũ chỉ dùng làm tham chiếu nội dung. Vẽ CÙNG LÚC với `master_hokage3_x4_chroma.png` ở brief 03-npc (chưa viết) để cùng một khuôn mặt. | **P2** | ve_lai |
| 11 | `kakashi_portrait_chroma.png` → `cd_kakashi_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | Chân dung NPC **Kakashi** (đề xuất Làng Lá, chờ KH9 chốt danh sách NPC — máy chủ hôm nay chỉ có Iruka). Jonin tóc bạc dựng lệch, băng trán navy kéo xuống che một mắt, khẩu trang vải navy che mũi-miệng, một mắt đen lười biếng, vest ô-liu. Bản cũ `cd_kakashi_x4.png` trong lô 343 là pixel-art 800×800 và **quay TRÁI** — lệch với Iruka đang lắp, nên vẽ lại theo khuôn chroma; ảnh cũ chỉ dùng làm tham chiếu nội dung. Vẽ CÙNG LÚC với `master_kakashi_x4_chroma.png` ở brief 03-npc (chưa viết) để cùng một khuôn mặt. | **P2** | ve_lai |
| 12 | `chu_ramen_portrait_chroma.png` → `cd_chu_ramen_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | Chân dung NPC **Chú bán mì** (đề xuất Làng Lá, chờ KH9 chốt danh sách NPC — máy chủ hôm nay chỉ có Iruka). Chú bếp quán mì trung niên: mặt tròn phúc hậu, khăn trắng buộc trán, cười tít mắt, cổ áo bếp trắng, tạp dề trắng. Bản cũ `cd_chu_ramen_x4.png` trong lô 343 là pixel-art 800×800 và **quay TRÁI** — lệch với Iruka đang lắp, nên vẽ lại theo khuôn chroma; ảnh cũ chỉ dùng làm tham chiếu nội dung. Vẽ CÙNG LÚC với `master_chu_ramen_x4_chroma.png` ở brief 03-npc (chưa viết) để cùng một khuôn mặt. | **P2** | ve_lai |
| 13 | `ba_vukhi_portrait_chroma.png` → `cd_ba_vukhi_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | Chân dung NPC **Bà chủ tiệm vũ khí** (đề xuất Làng Lá, chờ KH9 chốt danh sách NPC — máy chủ hôm nay chỉ có Iruka). Bà lão chủ tiệm vũ khí: tóc bạc búi cao, mắt sắc, nếp nhăn sâu, cổ kimono tím mận tối. Bản cũ `cd_ba_vukhi_x4.png` trong lô 343 là pixel-art 800×800 và **quay TRÁI** — lệch với Iruka đang lắp, nên vẽ lại theo khuôn chroma; ảnh cũ chỉ dùng làm tham chiếu nội dung. Vẽ CÙNG LÚC với `master_ba_vukhi_x4_chroma.png` ở brief 03-npc (chưa viết) để cùng một khuôn mặt. | **P2** | ve_lai |
| 14 | `thay_hocvien_portrait_chroma.png` → `cd_thay_hocvien_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | Chân dung NPC **Thầy học viện** (đề xuất Làng Lá, chờ KH9 chốt danh sách NPC — máy chủ hôm nay chỉ có Iruka). Giáo viên học viện nghiêm khắc: kính vuông có chấm phản sáng nhỏ, mày nhíu, tóc đen buộc, băng trán navy, vest ô-liu. Bản cũ `cd_thay_hocvien_x4.png` trong lô 343 là pixel-art 800×800 và **quay TRÁI** — lệch với Iruka đang lắp, nên vẽ lại theo khuôn chroma; ảnh cũ chỉ dùng làm tham chiếu nội dung. Vẽ CÙNG LÚC với `master_thay_hocvien_x4_chroma.png` ở brief 03-npc (chưa viết) để cùng một khuôn mặt. | **P2** | ve_lai |
| 15 | `linh_gac_portrait_chroma.png` → `cd_linh_gac_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | Chân dung NPC **Lính gác cổng** (đề xuất Làng Lá, chờ KH9 chốt danh sách NPC — máy chủ hôm nay chỉ có Iruka). Lính gác cổng chán chường đang ngáp: mắt lim dim, râu lún phún, miếng băng nhỏ trên má, băng trán navy, vest ô-liu. Bản cũ `cd_linh_gac_x4.png` trong lô 343 là pixel-art 800×800 và **quay TRÁI** — lệch với Iruka đang lắp, nên vẽ lại theo khuôn chroma; ảnh cũ chỉ dùng làm tham chiếu nội dung. Vẽ CÙNG LÚC với `master_linh_gac_x4_chroma.png` ở brief 03-npc (chưa viết) để cùng một khuôn mặt. | **P2** | ve_lai |
| 16 | `nguoi_giao_nv_portrait_chroma.png` → `cd_nguoi_giao_nv_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | Chân dung NPC **Người giao nhiệm vụ** (đề xuất Làng Lá, chờ KH9 chốt danh sách NPC — máy chủ hôm nay chỉ có Iruka). Nhân viên bàn giao nhiệm vụ điềm đạm: tóc đen chải gọn, băng trán navy, vest ô-liu, đầu hai cuộn giấy nhô sau vai. Bản cũ `cd_nguoi_giao_nv_x4.png` trong lô 343 là pixel-art 800×800 và **quay TRÁI** — lệch với Iruka đang lắp, nên vẽ lại theo khuôn chroma; ảnh cũ chỉ dùng làm tham chiếu nội dung. Vẽ CÙNG LÚC với `master_nguoi_giao_nv_x4_chroma.png` ở brief 03-npc (chưa viết) để cùng một khuôn mặt. | **P2** | ve_lai |
| 17 | `y_ta_portrait_chroma.png` → `cd_y_ta_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | Chân dung NPC **Y tá** (đề xuất Làng Lá, chờ KH9 chốt danh sách NPC — máy chủ hôm nay chỉ có Iruka). Y tá bệnh viện dịu dàng: mũ y tá trắng có dấu lá xanh nhỏ, tóc đen búi thấp, cổ áo đồng phục hồng nhạt. Bản cũ `cd_y_ta_x4.png` trong lô 343 là pixel-art 800×800 và **quay TRÁI** — lệch với Iruka đang lắp, nên vẽ lại theo khuôn chroma; ảnh cũ chỉ dùng làm tham chiếu nội dung. Vẽ CÙNG LÚC với `master_y_ta_x4_chroma.png` ở brief 03-npc (chưa viết) để cùng một khuôn mặt. | **P2** | ve_lai |
| 18 | `ong_tazuna_portrait_chroma.png` → `cd_ong_tazuna_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | chung | Chân dung NPC **Ông Tazuna** (đề xuất Làng Lá, chờ KH9 chốt danh sách NPC — máy chủ hôm nay chỉ có Iruka). Ông già dựng cầu dày dạn: râu ria lởm chởm bạc, vành nón rơm, kính tròn, khăn quàng cổ, da rám nắng. Bản cũ `cd_ong_tazuna_x4.png` trong lô 343 là pixel-art 800×800 và **quay TRÁI** — lệch với Iruka đang lắp, nên vẽ lại theo khuôn chroma; ảnh cũ chỉ dùng làm tham chiếu nội dung. Vẽ CÙNG LÚC với `master_ong_tazuna_x4_chroma.png` ở brief 03-npc (chưa viết) để cùng một khuôn mặt. | **P2** | ve_lai |
| 19 | `co_be_hoa_portrait_chroma.png` → `cd_co_be_hoa_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | Chân dung NPC **Cô bé bán hoa** (đề xuất Làng Lá, chờ KH9 chốt danh sách NPC — máy chủ hôm nay chỉ có Iruka). Bé gái bán hoa: má bầu, mắt to, tóc búi hai bên buộc ruy-băng, vài bông hoa từ giỏ nhô lên ở mép dưới. Bản cũ `cd_co_be_hoa_x4.png` trong lô 343 là pixel-art 800×800 và **quay TRÁI** — lệch với Iruka đang lắp, nên vẽ lại theo khuôn chroma; ảnh cũ chỉ dùng làm tham chiếu nội dung. Vẽ CÙNG LÚC với `master_co_be_hoa_x4_chroma.png` ở brief 03-npc (chưa viết) để cùng một khuôn mặt. | **P2** | ve_lai |
| 20 | `tho_ren_portrait_chroma.png` → `cd_tho_ren_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | Chân dung NPC **Thợ rèn** (đề xuất Làng Lá, chờ KH9 chốt danh sách NPC — máy chủ hôm nay chỉ có Iruka). Thợ rèn vạm vỡ: mặt vuông rắn rỏi, vài chấm mồ hôi, khăn vải quấn đầu, vai trần với dây tạp dề da, vết bụi than. Bản cũ `cd_tho_ren_x4.png` trong lô 343 là pixel-art 800×800 và **quay TRÁI** — lệch với Iruka đang lắp, nên vẽ lại theo khuôn chroma; ảnh cũ chỉ dùng làm tham chiếu nội dung. Vẽ CÙNG LÚC với `master_tho_ren_x4_chroma.png` ở brief 03-npc (chưa viết) để cùng một khuôn mặt. | **P2** | ve_lai |
| 21 | `nguoi_dua_thu_portrait_chroma.png` → `cd_nguoi_dua_thu_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | Chân dung NPC **Người đưa thư** (đề xuất Làng Lá, chờ KH9 chốt danh sách NPC — máy chủ hôm nay chỉ có Iruka). Giao liên nhanh nhẹn: mắt sắc tỉnh, tóc ngắn rối vì gió, băng trán navy, quai túi da vắt chéo vai. Bản cũ `cd_nguoi_dua_thu_x4.png` trong lô 343 là pixel-art 800×800 và **quay TRÁI** — lệch với Iruka đang lắp, nên vẽ lại theo khuôn chroma; ảnh cũ chỉ dùng làm tham chiếu nội dung. Vẽ CÙNG LÚC với `master_nguoi_dua_thu_x4_chroma.png` ở brief 03-npc (chưa viết) để cùng một khuôn mặt. | **P2** | ve_lai |
| 22 | `trongtai_thi_portrait_chroma.png` → `cd_trongtai_thi_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | Chân dung NPC **Trọng tài kỳ thi** (đề xuất Làng Lá, chờ KH9 chốt danh sách NPC — máy chủ hôm nay chỉ có Iruka). Trọng tài kỳ thi Chunin nghiêm nghị mệt mỏi: băng trán navy đội kiểu bandana, một cây kim senbon ngậm bên mép, cổ áo đồng phục đen-xanh than. Bản cũ `cd_trongtai_thi_x4.png` trong lô 343 là pixel-art 800×800 và **quay TRÁI** — lệch với Iruka đang lắp, nên vẽ lại theo khuôn chroma; ảnh cũ chỉ dùng làm tham chiếu nội dung. Vẽ CÙNG LÚC với `master_trongtai_thi_x4_chroma.png` ở brief 03-npc (chưa viết) để cùng một khuôn mặt. | **P2** | ve_lai |
| 23 | `nha_chep_cuon_portrait_chroma.png` → `cd_nha_chep_cuon_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Lá | Chân dung NPC **Nhà chép cuộn** (đề xuất Làng Lá, chờ KH9 chốt danh sách NPC — máy chủ hôm nay chỉ có Iruka). Thư lại gầy trầm ngâm: má hõm, bút lông gài sau tai, vệt mực nhỏ trên má, cổ áo học giả xám-be. Bản cũ `cd_nha_chep_cuon_x4.png` trong lô 343 là pixel-art 800×800 và **quay TRÁI** — lệch với Iruka đang lắp, nên vẽ lại theo khuôn chroma; ảnh cũ chỉ dùng làm tham chiếu nội dung. Vẽ CÙNG LÚC với `master_nha_chep_cuon_x4_chroma.png` ở brief 03-npc (chưa viết) để cùng một khuôn mặt. | **P2** | ve_lai |
| 24 | `cat_kazekage_portrait_chroma.png` → `cd_cat_kazekage_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Cát | Chân dung NPC **Phong Ảnh (Kazekage)** (đề xuất đợt 3 `dot-3-toan-bo\lang-*\npc.md`, chờ KH9 và máy chủ chốt NPC 4 làng — hôm nay chưa có làng nào ngoài Lá có NPC). Lãnh đạo Làng Cát lớn tuổi, râu, mũ Kage chóp trắng-vàng cát (KHÔNG chữ), khăn choàng xanh nhạt, áo choàng trắng-vàng cát, mắt sắc nghiêm. Vẽ CÙNG LÚC với master của NPC đó ở brief 03-npc để cùng khuôn mặt. Tên tra cứu `cd_cat_kazekage_x4` đã có trong kich-thuoc-anh.json (200×200) nhưng chưa có ảnh. | **P2** | can_ve |
| 25 | `cat_thay_portrait_chroma.png` → `cd_cat_thay_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Cát | Chân dung NPC **Thầy huấn luyện Suna** (đề xuất đợt 3 `dot-3-toan-bo\lang-*\npc.md`, chờ KH9 và máy chủ chốt NPC 4 làng — hôm nay chưa có làng nào ngoài Lá có NPC). Thầy trung niên Làng Cát: vest nâu cát, băng trán khắc đồng hồ cát, khăn che nửa mặt dưới, mắt hõm sắc, da rám. Vẽ CÙNG LÚC với master của NPC đó ở brief 03-npc để cùng khuôn mặt. Tên tra cứu `cd_cat_thay_x4` đã có trong kich-thuoc-anh.json (200×200) nhưng chưa có ảnh. | **P2** | can_ve |
| 26 | `cat_giao_nv_portrait_chroma.png` → `cd_cat_giao_nv_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Cát | Chân dung NPC **Người giao nhiệm vụ Suna** (đề xuất đợt 3 `dot-3-toan-bo\lang-*\npc.md`, chờ KH9 và máy chủ chốt NPC 4 làng — hôm nay chưa có làng nào ngoài Lá có NPC). Thương nhân sa mạc: khăn trùm kín đầu nâu sờn, mặt khắc khổ rám nắng, mắt nheo, áo choàng nâu, dây túi da chéo vai. Vẽ CÙNG LÚC với master của NPC đó ở brief 03-npc để cùng khuôn mặt. Tên tra cứu `cd_cat_giao_nv_x4` đã có trong kich-thuoc-anh.json (200×200) nhưng chưa có ảnh. | **P2** | can_ve |
| 27 | `smu_mizukage_portrait_chroma.png` → `cd_smu_mizukage_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Sương Mù | Chân dung NPC **Thuỷ Ảnh (Mizukage)** (đề xuất đợt 3 `dot-3-toan-bo\lang-*\npc.md`, chờ KH9 và máy chủ chốt NPC 4 làng — hôm nay chưa có làng nào ngoài Lá có NPC). Lãnh đạo Làng Sương Mù: mũ Kage chóp trắng-xanh dương với 4 đường sương lượn sóng (không chữ), khăn choàng lam nhạt, mặt lạnh nghiêm, da tái. Vẽ CÙNG LÚC với master của NPC đó ở brief 03-npc để cùng khuôn mặt. Tên tra cứu `cd_smu_mizukage_x4` đã có trong kich-thuoc-anh.json (200×200) nhưng chưa có ảnh. | **P2** | can_ve |
| 28 | `smu_thay_portrait_chroma.png` → `cd_smu_thay_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Sương Mù | Chân dung NPC **Thầy huấn luyện Kiri** (đề xuất đợt 3 `dot-3-toan-bo\lang-*\npc.md`, chờ KH9 và máy chủ chốt NPC 4 làng — hôm nay chưa có làng nào ngoài Lá có NPC). Thầy trung niên Sương Mù: giáp jonin xanh dương-xám ẩm, băng trán 4 đường sương, mặt nạ che nửa mặt dưới xám, mắt lạnh. Vẽ CÙNG LÚC với master của NPC đó ở brief 03-npc để cùng khuôn mặt. Tên tra cứu `cd_smu_thay_x4` đã có trong kich-thuoc-anh.json (200×200) nhưng chưa có ảnh. | **P2** | can_ve |
| 29 | `smu_giao_nv_portrait_chroma.png` → `cd_smu_giao_nv_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Sương Mù | Chân dung NPC **Người giao nhiệm vụ Kiri** (đề xuất đợt 3 `dot-3-toan-bo\lang-*\npc.md`, chờ KH9 và máy chủ chốt NPC 4 làng — hôm nay chưa có làng nào ngoài Lá có NPC). Ngư dân ven biển: nón lá rộng vành, mặt dày dạn sương gió, da tái lạnh, râu bạc lún phún, áo choàng xanh xám sờn, dây túi lưới chéo vai. Vẽ CÙNG LÚC với master của NPC đó ở brief 03-npc để cùng khuôn mặt. Tên tra cứu `cd_smu_giao_nv_x4` đã có trong kich-thuoc-anh.json (200×200) nhưng chưa có ảnh. | **P2** | can_ve |
| 30 | `da_tsuchikage_portrait_chroma.png` → `cd_da_tsuchikage_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Đá | Chân dung NPC **Thổ Ảnh (Tsuchikage)** (đề xuất đợt 3 `dot-3-toan-bo\lang-*\npc.md`, chờ KH9 và máy chủ chốt NPC 4 làng — hôm nay chưa có làng nào ngoài Lá có NPC). Lão nhân Làng Đá: nhỏ người, mũi to, ria mép và râu cằm nhọn trắng, mũ Kage chóp cam đất viền đỏ (không chữ), áo choàng cam đất, mắt nghiêm. Vẽ CÙNG LÚC với master của NPC đó ở brief 03-npc để cùng khuôn mặt. Tên tra cứu `cd_da_tsuchikage_x4` đã có trong kich-thuoc-anh.json (200×200) nhưng chưa có ảnh. | **P2** | can_ve |
| 31 | `da_thay_portrait_chroma.png` → `cd_da_thay_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Đá | Chân dung NPC **Thầy huấn luyện Iwa** (đề xuất đợt 3 `dot-3-toan-bo\lang-*\npc.md`, chờ KH9 và máy chủ chốt NPC 4 làng — hôm nay chưa có làng nào ngoài Lá có NPC). Thầy vạm vỡ Làng Đá: hàm vuông, băng trán hai phiến đá trên dải nâu đỏ, ốp phiến đá xám trên vai, vest nâu đất, da sạm. Vẽ CÙNG LÚC với master của NPC đó ở brief 03-npc để cùng khuôn mặt. Tên tra cứu `cd_da_thay_x4` đã có trong kich-thuoc-anh.json (200×200) nhưng chưa có ảnh. | **P2** | can_ve |
| 32 | `da_giao_nv_portrait_chroma.png` → `cd_da_giao_nv_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Đá | Chân dung NPC **Người giao nhiệm vụ Iwa** (đề xuất đợt 3 `dot-3-toan-bo\lang-*\npc.md`, chờ KH9 và máy chủ chốt NPC 4 làng — hôm nay chưa có làng nào ngoài Lá có NPC). Thợ đá: khăn trùm phủ bụi, mặt khắc khổ da sạm ráp, mắt nhỏ, râu lún phún, áo choàng nâu sờn, cán cuốc đá nhô sau vai. Vẽ CÙNG LÚC với master của NPC đó ở brief 03-npc để cùng khuôn mặt. Tên tra cứu `cd_da_giao_nv_x4` đã có trong kich-thuoc-anh.json (200×200) nhưng chưa có ảnh. | **P2** | can_ve |
| 33 | `may_raikage_portrait_chroma.png` → `cd_may_raikage_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Mây | Chân dung NPC **Lôi Ảnh (Raikage)** (đề xuất đợt 3 `dot-3-toan-bo\lang-*\npc.md`, chờ KH9 và máy chủ chốt NPC 4 làng — hôm nay chưa có làng nào ngoài Lá có NPC). Lãnh đạo Làng Mây: to lớn cơ bắp, da sẫm, tóc bạc dựng, mũ Kage chóp trắng-vàng chớp (không chữ), khăn choàng tím xanh nhạt, mắt nghiêm. Vẽ CÙNG LÚC với master của NPC đó ở brief 03-npc để cùng khuôn mặt. Tên tra cứu `cd_may_raikage_x4` đã có trong kich-thuoc-anh.json (200×200) nhưng chưa có ảnh. | **P2** | can_ve |
| 34 | `may_thay_portrait_chroma.png` → `cd_may_thay_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Mây | Chân dung NPC **Thầy huấn luyện Kumo** (đề xuất đợt 3 `dot-3-toan-bo\lang-*\npc.md`, chờ KH9 và máy chủ chốt NPC 4 làng — hôm nay chưa có làng nào ngoài Lá có NPC). Thầy dày dạn Làng Mây: hàm bạnh, băng trán khối mây trên dải xanh chàm, vest chiến thuật xám xanh, cơ bắp cổ vai, tóc bạch kim ngắn. Vẽ CÙNG LÚC với master của NPC đó ở brief 03-npc để cùng khuôn mặt. Tên tra cứu `cd_may_thay_x4` đã có trong kich-thuoc-anh.json (200×200) nhưng chưa có ảnh. | **P2** | can_ve |
| 35 | `may_giao_nv_portrait_chroma.png` → `cd_may_giao_nv_x4.png` | dải chroma (1 ô) | ô 1254×1254 × 1 khung → 200×200 | Mây | Chân dung NPC **Người giao nhiệm vụ Kumo** (đề xuất đợt 3 `dot-3-toan-bo\lang-*\npc.md`, chờ KH9 và máy chủ chốt NPC 4 làng — hôm nay chưa có làng nào ngoài Lá có NPC). Sứ giả núi cao: mặt khắc khổ da rám gió, khăn quàng dày, áo choàng xám sờn, dây túi da đựng cuộn thư chéo vai. Vẽ CÙNG LÚC với master của NPC đó ở brief 03-npc để cùng khuôn mặt. Tên tra cứu `cd_may_giao_nv_x4` đã có trong kich-thuoc-anh.json (200×200) nhưng chưa có ảnh. | **P2** | can_ve |
| 36 | `ui_icon_app_x4.png` | tĩnh (đục) | 512×512 (vẽ 1024×1024) | chung | **Icon ứng dụng** — thay icon hiện tại của client vẫn là **Rồng + ngọc 4 sao Dragon Ball** (`D:\NRO_Project\Assets\Resources\appicon.png` 479×480, `ProjectSettings.asset:292-335` mọi nền tảng trỏ guid b611de4c…). Emblem lá-lửa xoáy (TOT-naruto-logo) đặt giữa nền navy-than tối, KHÔNG chữ, KHÔNG bo góc (hệ điều hành tự bo), đục hoàn toàn. Vẽ 1024×1024 (App Store cần 1024), dự án thu về 512 (Google Play/Unity). | **P1** | can_ve |
| 37 | `ui_icon_app_hinh_x4.png` | tĩnh (đục) | 432×432 | chung | Lớp **hình** của icon thích ứng Android (adaptive icon foreground): chỉ emblem, nền TRONG SUỐT, nằm trong vòng an toàn 66 % giữa (Android cắt tròn/vuông bo tuỳ máy). Cặp với `ui_icon_app_nen_x4`. | **P3** | can_ve |
| 38 | `ui_icon_app_nen_x4.png` | tĩnh (đục) | 432×432 | chung | Lớp **nền** của icon thích ứng Android (adaptive icon background): chỉ nền navy-than có ánh sáng toả nhẹ ở giữa, không hình, đục hoàn toàn. Cặp với `ui_icon_app_hinh_x4`. | **P3** | can_ve |
| 39 | `ui_splash_doc_x4.png` | tĩnh (đục) | 1080×1920 | Lá | **Màn chờ dọc** 1080×1920 khi mở game (Unity Player Settings: `m_ShowUnitySplashScreen: 1`, `m_SplashScreenLogos: []`, `androidSplashScreen` = 0 — chưa có ảnh nào). Trên: emblem lá-lửa lớn. Giữa (45–60 % chiều cao): **chừa trống nền tối** để dự án dán wordmark `Logo\1.png` — KHÔNG vẽ chữ. Dưới: bóng Làng Lá lúc chạng vạng (núi Hokage không mặt người thật, cổng torii đỏ, mái ngói), trời cam ấm → navy đậm ở đỉnh. Không nhân vật. Đục. Lưu ý: `/e/sp0..2.png` (GameScr.cs:5950) là hiệu ứng toé nước, không phải màn chờ. | **P2** | can_ve |
| 40 | `ui_splash_ngang_x4.png` | tĩnh (đục) | 1920×1080 | Lá | **Màn chờ ngang** 1920×1080 (bản PC / máy tính bảng xoay ngang). Cùng chất liệu với bản dọc: emblem ở trái-giữa, vùng trống cho wordmark ở bên phải emblem, bóng Làng Lá chạy dọc mép dưới. Không chữ, không nhân vật, đục. | **P2** | can_ve |
| 41 | `ui_bia_cuahang_x4.png` | tĩnh (đục) | 1024×500 | Lá | **Ảnh bìa cửa hàng** 1024×500 (Google Play feature graphic; cũng dùng làm ảnh bìa fanpage). Người chơi học viên Làng Lá (đúng mặt/tóc/vest của `player_idle_3f_chroma`) ở 1/3 trái, thầy Iruka mờ hơn phía sau, một con Sói Rừng nhỏ; nền núi Hokage + cổng torii; 5 vệt chakra màu 5 làng (cam · vàng cát · lam · nâu · tím xanh) quét ngang. Chừa vùng tối yên ở giữa-trên cho wordmark. KHÔNG chữ. Đục. Lưu ý ngoài mỹ thuật (chưa chắc): wordmark NARUTO là nhãn hiệu của bên khác — đưa lên cửa hàng ứng dụng có rủi ro bị gỡ; chủ dự án quyết trước khi dùng. | **P2** | can_ve |
| 42 | `ui_nen_dangnhap_tet_x4.png` | tĩnh (đục) | 1024×576 | Lá | **Banner đăng nhập theo mùa — Tết (xuân)** 1024×576: cùng bố cục với `ui_nen_dangnhap_x4` lô 343 (Làng Lá nhìn từ sườn đồi, mái ngói đổ bóng, trăng lớn) nhưng vẽ lại hand-painted và đổi không khí: đèn lồng đỏ treo dọc phố, cành đào-mai nở hồng-vàng ở tiền cảnh, pháo hoa nhỏ trên trời đêm, tuyết nhẹ trên mái ngói. KHÔNG chữ. Đục. Chưa có chỗ lắp: `LoginScr.cs:655-667` hôm nay chỉ vẽ logo, cần code nạp nền trước — P3. | **P3** | can_ve |
| 43 | `ui_nen_dangnhap_he_x4.png` | tĩnh (đục) | 1024×576 | Lá | **Banner đăng nhập theo mùa — Hè** 1024×576: cùng bố cục với `ui_nen_dangnhap_x4` lô 343 (Làng Lá nhìn từ sườn đồi, mái ngói đổ bóng, trăng lớn) nhưng vẽ lại hand-painted và đổi không khí: trưa nắng gắt, trời xanh cao có mây trắng lớn, cờ hội hè phấp phới, hướng dương ven đường, bóng đổ ngắn đậm. KHÔNG chữ. Đục. Chưa có chỗ lắp: `LoginScr.cs:655-667` hôm nay chỉ vẽ logo, cần code nạp nền trước — P3. | **P3** | can_ve |
| 44 | `ui_nen_dangnhap_trungthu_x4.png` | tĩnh (đục) | 1024×576 | Lá | **Banner đăng nhập theo mùa — Trung thu** 1024×576: cùng bố cục với `ui_nen_dangnhap_x4` lô 343 (Làng Lá nhìn từ sườn đồi, mái ngói đổ bóng, trăng lớn) nhưng vẽ lại hand-painted và đổi không khí: trăng tròn vàng lớn ngay giữa, đèn ông sao và đèn cá chép treo trước hiên, bàn cỗ nhỏ, trời chàm tím xanh. KHÔNG chữ. Đục. Chưa có chỗ lắp: `LoginScr.cs:655-667` hôm nay chỉ vẽ logo, cần code nạp nền trước — P3. | **P3** | can_ve |

**Tên MỚI cần nối cuối `kich-thuoc-anh.json` (kèm cỡ x4):**

- `cd_nguoi_choi_la_x4.png = 200 200`
- `cd_nguoi_choi_cat_x4.png = 200 200`
- `cd_nguoi_choi_smu_x4.png = 200 200`
- `cd_nguoi_choi_da_x4.png = 200 200`
- `cd_nguoi_choi_may_x4.png = 200 200`
- `cd_boss_nhan_gia_mat_na_p1_x4.png = 200 200`
- `cd_boss_nhan_gia_mat_na_p2_x4.png = 200 200`
- `cd_boss_hoc_vien_dao_tau_x4.png = 200 200`
- `ui_icon_app_x4.png = 512 512`
- `ui_icon_app_hinh_x4.png = 432 432`
- `ui_icon_app_nen_x4.png = 432 432`
- `ui_splash_doc_x4.png = 1080 1920`
- `ui_splash_ngang_x4.png = 1920 1080`
- `ui_bia_cuahang_x4.png = 1024 500`
- `ui_nen_dangnhap_tet_x4.png = 1024 576`
- `ui_nen_dangnhap_he_x4.png = 1024 576`
- `ui_nen_dangnhap_trungthu_x4.png = 1024 576`

Các `cd_*` còn lại (Iruka, 14 NPC Làng Lá, 12 NPC 4 làng) **đã có tên** trong bảng cỡ (200×200) — không thêm.

---

## 6. PROMPT TỪNG ẢNH

Mỗi khối là prompt **tự đứng được** (đã gồm phong cách khoá + subject + nội dung riêng + cỡ/ô + nền + góc nhìn + điều không vẽ). Dán nguyên khối, gắn kèm các ảnh tham chiếu ghi ở dòng *Tham chiếu*. Dòng `OUTPUT` là tên file, cỡ, số khung và thư mục lưu.

### 6.1 `nguoi_choi_la_portrait_chroma.png` — P1 · can_ve

*Tham chiếu:* `player_idle_3f_chroma`, `master_hoc_vien_x4_chroma`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the Leaf Village academy-student PLAYER character exactly as in the reference idle frame (player_idle_3f_chroma.png, frame 0): a boy about 12 years old, tan skin, large brown eyes with a determined brow, dark brown-black spiky hair swept back and up, dark navy cloth forehead protector with a silver metal plate engraved with the Leaf spiral symbol (plate centered above the eyes), olive-green flak vest with two chest pockets over a dark navy high-collar shirt, the knot of an orange sash just visible at the bottom cut-line. Calm, friendly, slightly stubborn expression, mouth closed. Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: this face will sit next to the running sprite of the same boy — copy hair shape, eye shape, skin tone and headband from the reference; do not restyle him, do not age him, do not add a scar.
```
OUTPUT: nguoi_choi_la_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.2 `nguoi_choi_cat_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `player_idle_3f_chroma`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the default PLAYER character of the Sand Village: a ninja academy student about 12 years old with the SAME head size, face proportions and body build as the Leaf student in the reference (player_idle_3f_chroma.png) but a different look: sandy-blond short hair with a cream desert head-wrap, tan sun-weathered skin, amber eyes, forehead protector plate engraved with the Sand hourglass symbol stitched onto the head-wrap, a sand-yellow scarf bunched at the neck over a khaki high-collar shirt with an ochre shoulder strap. Calm, friendly, slightly stubborn expression, mouth closed. Sand Village / Wind nature palette: sand-yellow, ochre, cream and pale desert green.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: keep the child-like proportions (big eyes, short nose, round cheek) identical to the Leaf student so all five starting characters read as the same age; only hair, skin, headband symbol and clothing colors change.
```
OUTPUT: nguoi_choi_cat_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.3 `nguoi_choi_smu_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `player_idle_3f_chroma`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the default PLAYER character of the Mist Village: a ninja academy student about 12 years old with the SAME head size, face proportions and body build as the Leaf student in the reference (player_idle_3f_chroma.png) but a different look: dark blue-black straight hair hanging over one eyebrow with a damp look, pale cool skin, grey-blue eyes, forehead protector plate engraved with the four wavy Mist lines on a slate-blue band, a deep-blue high collar with a teal striped inner lining, a thin strip of bandage on one cheek. Calm, friendly, slightly stubborn expression, mouth closed. Mist Village / Water nature palette: deep blue, teal and slate grey.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: keep the child-like proportions (big eyes, short nose, round cheek) identical to the Leaf student so all five starting characters read as the same age; only hair, skin, headband symbol and clothing colors change.
```
OUTPUT: nguoi_choi_smu_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.4 `nguoi_choi_da_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `player_idle_3f_chroma`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the default PLAYER character of the Stone Village: a ninja academy student about 12 years old with the SAME head size, face proportions and body build as the Leaf student in the reference (player_idle_3f_chroma.png) but a different look: dark brown cropped hair, sturdy square jaw, earth-tan skin, dark brown eyes, forehead protector plate engraved with the two Stone slabs on a rust-red band, a heavy earth-brown scarf over a rock-grey padded collar with brown leather trim. Calm, friendly, slightly stubborn expression, mouth closed. Stone Village / Earth nature palette: earth brown, rust and rock grey.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: keep the child-like proportions (big eyes, short nose, round cheek) identical to the Leaf student so all five starting characters read as the same age; only hair, skin, headband symbol and clothing colors change.
```
OUTPUT: nguoi_choi_da_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.5 `nguoi_choi_may_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `player_idle_3f_chroma`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the default PLAYER character of the Cloud Village: a ninja academy student about 12 years old with the SAME head size, face proportions and body build as the Leaf student in the reference (player_idle_3f_chroma.png) but a different look: short white-silver hair swept to one side, darker brown skin, sharp gold eyes, forehead protector plate engraved with the Cloud symbol on a bluish indigo band, a white high collar with indigo-violet shoulders and a thin lightning-yellow trim. Calm, friendly, slightly stubborn expression, mouth closed. Cloud Village / Lightning nature palette: BLUISH indigo-violet, white and lightning yellow (keep every purple bluish: R-G below 100 — no magenta, no hot pink, no red-violet).
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: keep the child-like proportions (big eyes, short nose, round cheek) identical to the Leaf student so all five starting characters read as the same age; only hair, skin, headband symbol and clothing colors change.
```
OUTPUT: nguoi_choi_may_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.6 `nhan_gia_mat_na_p1_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `player_idle_3f_chroma`, `master_hoc_vien_x4_chroma`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the boss "The Stranger" (Masked Impostor, phase 1) disguised as a Leaf academy student: wears EXACTLY the student uniform of the player reference (olive-green flak vest with chest pockets, dark navy high-collar shirt) but the face is completely hidden behind a smooth matte off-white ceramic mask with two narrow dark eye slits and a single thin red crack running diagonally across the left cheek of the mask; long unkempt black hair escapes behind and below the mask; no forehead protector; shoulders slightly hunched, still, menacing. Adult-sized head and shoulders (about 15% larger than the student player, this is a grown man in a student's clothes). Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy, with the red crack as the only red accent.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: no eyes visible through the slits (pure dark), no glow, no fire effect; the mask is painted ceramic with soft shading, not flat white.
```
OUTPUT: nhan_gia_mat_na_p1_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.7 `nhan_gia_mat_na_p2_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `master_iruka_x4_chroma`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the boss "Broken Mask" (Masked Impostor, phase 2, true form revealed): a gaunt sharp-featured man about 30 years old, pale-tan skin, long straight black hair with one grey streak falling past the shoulders, narrow cold dark-red eyes, a diagonal old scar across the left cheek (NOT across the nose), a Leaf Village forehead protector whose metal plate has the Leaf symbol SCRATCHED THROUGH by one deep horizontal gash (rogue ninja), the broken lower half of a white ceramic mask still hanging from a cord at the side of his jaw, a high dark-crimson-black cloak collar with an ember-orange inner lining. Expression: bitter, contemptuous half-smile. Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: he must NOT resemble the NPC Iruka (reference master_iruka_x4_chroma.png): no ponytail, no scar across the nose, no olive vest, different face shape. No fire effect, no glow.
```
OUTPUT: nhan_gia_mat_na_p2_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.8 `hoc_vien_dao_tau_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `player_idle_3f_chroma`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the boss "Runaway Student" (Mist Village academy defector hiding in Leaf territory): a thin teenage boy about 14, pale cool skin, wet-looking dark teal-black hair hanging over his right eye, tired eyes with dark circles, a Mist Village forehead protector (metal plate with four wavy lines on a slate-blue band) worn loose and crooked around the NECK instead of the forehead, the same academy student uniform as the player reference (olive-green flak vest, dark navy high collar) but dirty, frayed and darkened by rain, a small strip of bandage on his left cheek. Expression: haunted, defensive. Mist Village / Water nature palette: deep blue, teal and slate grey for hair, headband and shadows.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: no water droplets or glow effects, no scar across the nose (that belongs to Iruka).
```
OUTPUT: hoc_vien_dao_tau_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.9 `iruka_portrait_chroma.png` — P3 · da_co · **KHÔNG VẼ đợt này, prompt chỉ để dự phòng**

*Tham chiếu:* `master_iruka_x4_chroma`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the NPC Iruka, Leaf Village academy teacher: brown hair pulled up into a short spiky ponytail, tan skin, warm brown eyes, a horizontal scar across the bridge of the nose, dark navy forehead protector with the Leaf plate, olive-green flak vest over a dark navy high-collar shirt, small goatee. Kind, encouraging smile. Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: only redraw if the project owner explicitly asks — the current iruka_portrait_chroma.png is the approved standard.
```
OUTPUT: iruka_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.10 `hokage3_portrait_chroma.png` — P2 · ve_lai

*Tham chiếu:* `npc343_hokage3`, `cd343_hokage3`, `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the Third Hokage, the elderly village leader: deeply wrinkled kind face, white goatee and thin white moustache, heavy-lidded wise eyes, a red-and-white conical Kage hat with a white cloth veil at the sides (no writing on the hat), a white robe collar with red trim over a dark inner robe, a long thin wooden pipe held at the corner of the mouth with one short opaque wisp of painted grey smoke. Expression: tired but warm. Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: use the two reference images of this NPC from the old batch (npc_hokage3_x4.png, cd_hokage3_x4.png) ONLY for who the character is (hair, clothes, accessories); ignore their pixel-art style and ignore that they face left — this portrait faces RIGHT like iruka_portrait_chroma.png.
```
OUTPUT: hokage3_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.11 `kakashi_portrait_chroma.png` — P2 · ve_lai

*Tham chiếu:* `npc343_kakashi`, `cd343_kakashi`, `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the jonin Kakashi: silver-white spiky hair leaning to one side, a dark navy forehead protector pulled down slanted to cover the left eye, a dark navy cloth mask covering the nose and mouth, one visible half-lidded dark eye with a relaxed bored expression, olive-green flak vest with a high dark navy collar. Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: use the two reference images of this NPC from the old batch (npc_kakashi_x4.png, cd_kakashi_x4.png) ONLY for who the character is (hair, clothes, accessories); ignore their pixel-art style and ignore that they face left — this portrait faces RIGHT like iruka_portrait_chroma.png.
```
OUTPUT: kakashi_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.12 `chu_ramen_portrait_chroma.png` — P2 · ve_lai

*Tham chiếu:* `npc343_chu_ramen`, `cd343_chu_ramen`, `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the ramen-stand cook: a round-faced cheerful middle-aged man, ruddy cheeks, eyes squeezed shut by a big warm smile, a white cloth headband tied around the forehead, short brown hair, a white chef's tunic collar and the top of a white apron strap. Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: use the two reference images of this NPC from the old batch (npc_chu_ramen_x4.png, cd_chu_ramen_x4.png) ONLY for who the character is (hair, clothes, accessories); ignore their pixel-art style and ignore that they face left — this portrait faces RIGHT like iruka_portrait_chroma.png.
```
OUTPUT: chu_ramen_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.13 `ba_vukhi_portrait_chroma.png` — P2 · ve_lai

*Tham chiếu:* `npc343_ba_vukhi`, `cd343_ba_vukhi`, `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the weapon-shop owner, an old woman: silver-grey hair pulled into a tight high bun with a plain wooden hairpin, sharp narrow appraising eyes, deep wrinkles, thin pressed lips, a dark plum kimono collar layered over a charcoal under-robe. Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: use the two reference images of this NPC from the old batch (npc_ba_vukhi_x4.png, cd_ba_vukhi_x4.png) ONLY for who the character is (hair, clothes, accessories); ignore their pixel-art style and ignore that they face left — this portrait faces RIGHT like iruka_portrait_chroma.png.
```
OUTPUT: ba_vukhi_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.14 `thay_hocvien_portrait_chroma.png` — P2 · ve_lai

*Tham chiếu:* `npc343_thay_hocvien`, `cd343_thay_hocvien`, `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the strict academy instructor: a man in his thirties with black hair tied back, square-framed glasses with one small painted highlight dot on each lens (not a glowing rim), furrowed brows, thin stern mouth, dark navy forehead protector with the Leaf plate, olive-green flak vest over a dark navy collar. Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: use the two reference images of this NPC from the old batch (npc_thay_hocvien_x4.png, cd_thay_hocvien_x4.png) ONLY for who the character is (hair, clothes, accessories); ignore their pixel-art style and ignore that they face left — this portrait faces RIGHT like iruka_portrait_chroma.png.
```
OUTPUT: thay_hocvien_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.15 `linh_gac_portrait_chroma.png` — P2 · ve_lai

*Tham chiếu:* `npc343_linh_gac`, `cd343_linh_gac`, `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the bored gate guard mid-yawn: half-closed drowsy eyes, mouth slightly open in a yawn, light stubble on the chin, a small bandage strip on the cheek (NOT across the nose), messy dark hair under a dark navy forehead protector with the Leaf plate, olive-green flak vest with a slouched collar. Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: use the two reference images of this NPC from the old batch (npc_linh_gac_x4.png, cd_linh_gac_x4.png) ONLY for who the character is (hair, clothes, accessories); ignore their pixel-art style and ignore that they face left — this portrait faces RIGHT like iruka_portrait_chroma.png.
```
OUTPUT: linh_gac_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.16 `nguoi_giao_nv_portrait_chroma.png` — P2 · ve_lai

*Tham chiếu:* `npc343_nguoi_giao_nv`, `cd343_nguoi_giao_nv`, `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the mission-desk clerk: a calm composed man with neatly combed short black hair, even brows, a polite closed-mouth expression, dark navy forehead protector with the Leaf plate, olive-green flak vest over a dark navy collar, the tips of two rolled paper scrolls peeking above his far shoulder. Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: use the two reference images of this NPC from the old batch (npc_nguoi_giao_nv_x4.png, cd_nguoi_giao_nv_x4.png) ONLY for who the character is (hair, clothes, accessories); ignore their pixel-art style and ignore that they face left — this portrait faces RIGHT like iruka_portrait_chroma.png.
```
OUTPUT: nguoi_giao_nv_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.17 `y_ta_portrait_chroma.png` — P2 · ve_lai

*Tham chiếu:* `npc343_y_ta`, `cd343_y_ta`, `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the hospital nurse: a gentle young woman with black hair in a low bun, soft eyes and a small reassuring smile, a white nurse's cap with a small green leaf mark (no red cross), a pale pink uniform collar with a white trim. Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: use the two reference images of this NPC from the old batch (npc_y_ta_x4.png, cd_y_ta_x4.png) ONLY for who the character is (hair, clothes, accessories); ignore their pixel-art style and ignore that they face left — this portrait faces RIGHT like iruka_portrait_chroma.png.
```
OUTPUT: y_ta_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.18 `ong_tazuna_portrait_chroma.png` — P2 · ve_lai

*Tham chiếu:* `npc343_ong_tazuna`, `cd343_ong_tazuna`, `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the old bridge builder Tazuna: a weathered sun-tanned old man with grey stubble beard, round wire glasses, a wide straw hat brim shading the brow, a rough cotton towel-scarf wrapped around the neck, a faded grey work jacket collar. Expression: gruff, squinting. neutral warm palette, no village accent.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: use the two reference images of this NPC from the old batch (npc_ong_tazuna_x4.png, cd_ong_tazuna_x4.png) ONLY for who the character is (hair, clothes, accessories); ignore their pixel-art style and ignore that they face left — this portrait faces RIGHT like iruka_portrait_chroma.png.
```
OUTPUT: ong_tazuna_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.19 `co_be_hoa_portrait_chroma.png` — P2 · ve_lai

*Tham chiếu:* `npc343_co_be_hoa`, `cd343_co_be_hoa`, `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the little flower girl: a small child about 7 years old with chubby cheeks, huge bright eyes, hair in two side buns tied with pale yellow ribbons, a pastel green kimono collar, and three or four painted flowers (pink, yellow, white) rising from a woven basket just above the bottom cut-line. Sweet cheerful expression. Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: use the two reference images of this NPC from the old batch (npc_co_be_hoa_x4.png, cd_co_be_hoa_x4.png) ONLY for who the character is (hair, clothes, accessories); ignore their pixel-art style and ignore that they face left — this portrait faces RIGHT like iruka_portrait_chroma.png.
```
OUTPUT: co_be_hoa_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.20 `tho_ren_portrait_chroma.png` — P2 · ve_lai

*Tham chiếu:* `npc343_tho_ren`, `cd343_tho_ren`, `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the blacksmith: a burly man with a square rugged face, thick neck and shoulders, a few small painted sweat highlight dots on the brow, a dark cloth wrapped around the head, bare tanned shoulders with a brown leather apron strap over one shoulder, soot smudges on the cheek. Expression: focused, brows down. Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: use the two reference images of this NPC from the old batch (npc_tho_ren_x4.png, cd_tho_ren_x4.png) ONLY for who the character is (hair, clothes, accessories); ignore their pixel-art style and ignore that they face left — this portrait faces RIGHT like iruka_portrait_chroma.png.
```
OUTPUT: tho_ren_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.21 `nguoi_dua_thu_portrait_chroma.png` — P2 · ve_lai

*Tham chiếu:* `npc343_nguoi_dua_thu`, `cd343_nguoi_dua_thu`, `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the courier: a lean quick young man with wind-tousled short dark hair, sharp alert eyes, a light smirk, dark navy forehead protector with the Leaf plate, a brown leather messenger-bag strap crossing the chest over a dark navy jacket collar. Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: use the two reference images of this NPC from the old batch (npc_nguoi_dua_thu_x4.png, cd_nguoi_dua_thu_x4.png) ONLY for who the character is (hair, clothes, accessories); ignore their pixel-art style and ignore that they face left — this portrait faces RIGHT like iruka_portrait_chroma.png.
```
OUTPUT: nguoi_dua_thu_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.22 `trongtai_thi_portrait_chroma.png` — P2 · ve_lai

*Tham chiếu:* `npc343_trongtai_thi`, `cd343_trongtai_thi`, `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the chunin-exam proctor: a serious tired-looking man with a dark navy forehead protector worn bandana-style over shoulder-length brown hair, a thin steel senbon needle held between his lips at the corner of the mouth, heavy eyelids, a black and dark-teal uniform collar. Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: use the two reference images of this NPC from the old batch (npc_trongtai_thi_x4.png, cd_trongtai_thi_x4.png) ONLY for who the character is (hair, clothes, accessories); ignore their pixel-art style and ignore that they face left — this portrait faces RIGHT like iruka_portrait_chroma.png.
```
OUTPUT: trongtai_thi_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.23 `nha_chep_cuon_portrait_chroma.png` — P2 · ve_lai

*Tham chiếu:* `npc343_nha_chep_cuon`, `cd343_nha_chep_cuon`, `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the scroll scribe: a thin scholarly man with hollow cheeks, thoughtful downcast eyes, a bamboo writing brush tucked behind his ear, a small black ink smudge on the cheekbone, a grey-beige scholar's robe collar with a darker inner layer. Leaf Village / Fire nature palette: warm orange-red with olive green and dark navy.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: use the two reference images of this NPC from the old batch (npc_nha_chep_cuon_x4.png, cd_nha_chep_cuon_x4.png) ONLY for who the character is (hair, clothes, accessories); ignore their pixel-art style and ignore that they face left — this portrait faces RIGHT like iruka_portrait_chroma.png.
```
OUTPUT: nha_chep_cuon_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.24 `cat_kazekage_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the Kazekage, elderly leader of the Sand Village: a stern lined face with a short grey beard, sharp narrow eyes, a white-and-sand-yellow conical Kage hat with a pale cream veil at the sides (no writing), a pale desert-green shawl over a white-and-ochre robe collar. Sand Village / Wind nature palette: sand-yellow, ochre, cream and pale desert green.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: no writing or kanji on any hat or plate — village identity comes only from the engraved symbol and the palette.
```
OUTPUT: cat_kazekage_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.25 `cat_thay_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the Sand Village instructor: a middle-aged man with sun-darkened skin, deep-set sharp eyes, a sand-beige cloth mask covering the lower half of the face, a forehead protector with the Sand hourglass symbol on a khaki band, a sand-brown flak vest with a cream high collar. Sand Village / Wind nature palette: sand-yellow, ochre, cream and pale desert green.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: no writing or kanji on any hat or plate — village identity comes only from the engraved symbol and the palette.
```
OUTPUT: cat_thay_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.26 `cat_giao_nv_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the Sand Village mission clerk, a desert merchant: a gaunt sun-weathered face with squinting eyes and dry cracked lips, a frayed brown head-wrap covering hair and neck, a rough ochre travelling cloak collar, a worn leather bag strap across the chest. Sand Village / Wind nature palette: sand-yellow, ochre, cream and pale desert green.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: no writing or kanji on any hat or plate — village identity comes only from the engraved symbol and the palette.
```
OUTPUT: cat_giao_nv_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.27 `smu_mizukage_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the Mizukage, leader of the Mist Village: a cold composed face with pale skin and piercing blue-grey eyes, a white-and-deep-blue conical Kage hat marked with four wavy mist lines (no writing) and a pale blue side veil, a light blue shawl over a white-and-navy robe collar. Mist Village / Water nature palette: deep blue, teal and slate grey.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: no writing or kanji on any hat or plate — village identity comes only from the engraved symbol and the palette.
```
OUTPUT: smu_mizukage_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.28 `smu_thay_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the Mist Village instructor: a middle-aged man with pale skin and cold steady eyes, a slate-grey half-mask covering nose and mouth, a forehead protector with the four wavy Mist lines on a dark blue band, a damp-looking blue-grey flak vest with a high teal collar. Mist Village / Water nature palette: deep blue, teal and slate grey.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: no writing or kanji on any hat or plate — village identity comes only from the engraved symbol and the palette.
```
OUTPUT: smu_thay_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.29 `smu_giao_nv_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the Mist Village mission clerk, a coastal fisherman: a weather-beaten face with pale cold skin and grey stubble, a wide woven reed hat, a frayed blue-grey oilcloth cloak collar, the strap of a net bag across the chest, a few painted water stains on the cloak. Mist Village / Water nature palette: deep blue, teal and slate grey.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: no writing or kanji on any hat or plate — village identity comes only from the engraved symbol and the palette.
```
OUTPUT: smu_giao_nv_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.30 `da_tsuchikage_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the Tsuchikage, elderly leader of the Stone Village: a small wiry old man with a large prominent nose, a pointed white goatee and thin white moustache, stern hooded eyes, an earth-orange conical Kage hat with red trim and a cream side veil (no writing), an earth-orange robe collar edged in dark red. Stone Village / Earth nature palette: earth brown, rust and rock grey.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: no writing or kanji on any hat or plate — village identity comes only from the engraved symbol and the palette.
```
OUTPUT: da_tsuchikage_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.31 `da_thay_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the Stone Village instructor: a broad-shouldered middle-aged man with a square jaw, heavy brows and earth-tan skin, a forehead protector with the two Stone slabs on a rust-brown band, a rock-grey stone shoulder plate strapped over an earth-brown flak vest with a beige collar. Stone Village / Earth nature palette: earth brown, rust and rock grey.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: no writing or kanji on any hat or plate — village identity comes only from the engraved symbol and the palette.
```
OUTPUT: da_thay_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.32 `da_giao_nv_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the Stone Village mission clerk, a quarry worker: a dusty rough-skinned face with small squinting eyes and grey stubble, a dust-covered brown head cloth, a frayed earth-brown cloak collar, the wooden handle of a stone pick visible above the far shoulder. Stone Village / Earth nature palette: earth brown, rust and rock grey.
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: no writing or kanji on any hat or plate — village identity comes only from the engraved symbol and the palette.
```
OUTPUT: da_giao_nv_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.33 `may_raikage_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the Raikage, leader of the Cloud Village: a massive muscular man with dark brown skin, swept-back white hair, a hard commanding face with a thick neck, a white-and-lightning-yellow conical Kage hat with a pale side veil (no writing), a pale bluish-lavender shawl over a white robe collar with indigo trim. Cloud Village / Lightning nature palette: BLUISH indigo-violet, white and lightning yellow (keep every purple bluish: R-G below 100 — no magenta, no hot pink, no red-violet).
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: no writing or kanji on any hat or plate — village identity comes only from the engraved symbol and the palette.
```
OUTPUT: may_raikage_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.34 `may_thay_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the Cloud Village instructor: a rugged man with a wide jaw, short platinum hair, dark skin and a thick muscular neck, a forehead protector with the Cloud symbol on a bluish indigo band, a grey-blue tactical flak vest with a white high collar and a thin lightning-yellow trim. Cloud Village / Lightning nature palette: BLUISH indigo-violet, white and lightning yellow (keep every purple bluish: R-G below 100 — no magenta, no hot pink, no red-violet).
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: no writing or kanji on any hat or plate — village identity comes only from the engraved symbol and the palette.
```
OUTPUT: may_thay_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.35 `may_giao_nv_portrait_chroma.png` — P2 · can_ve

*Tham chiếu:* `TOT-npc-ong-tazuna`, `iruka_portrait_chroma`, `iruka_portrait_contact`, `THUOC-NHAN-VAT`, `TOT-naruto-logo`, `TOT-chitiet-nhanvat`, `XAU-npc`, `XAU-boss-zabuza`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, 3/4 view facing RIGHT, bust portrait crop (head and shoulders — NOT full body), NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the Cloud Village mission clerk, a mountain messenger: a wind-burnt weathered face with narrowed eyes, a thick bluish-grey wool scarf wrapped around the neck, a frayed grey travelling cloak collar, a leather scroll-case strap across the chest. Cloud Village / Lightning nature palette: BLUISH indigo-violet, white and lightning yellow (keep every purple bluish: R-G below 100 — no magenta, no hot pink, no red-violet).
CONSISTENCY: same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242). For a bust portrait the "feet" rule means the straight horizontal shoulder cut-line sits exactly on the bottom edge of the cell (in game the portrait stands on top of the dialog frame). Face, hairstyle, forehead protector and outfit colors must match this character's own reference image exactly — a player cannot mistake it for another character.
CANVAS: one single cell 1254x1254 px, exactly one frame (the file name has no frame number, the pipeline treats it as 1 cell). The bust fills 85-92% of the cell height; top of the hair at least 40 px below the top edge; at least 60 px of pure magenta on the left and on the right side; face height 45-55% of the cell so the face still reads at 50x50 px in the dialog box and at 26x26 px on the HUD.
BACKGROUND: every background pixel exactly RGB(246,4,242), perfectly flat — no gradient, no vignette, no noise, no glow around the figure, no shadow, no outline halo; hard opaque edges only.
VIEW: 3/4 view facing RIGHT (the viewer's right), head upright (not tilted up or down), eyes looking forward-right.
DO NOT DRAW: text, letters, numbers, logos, decorative frame or border, background scenery, drop shadow, glow, motion lines, hands, weapons crossing the face, a second character, semi-transparent or blurred edges; and no magenta / hot-pink / red-violet color anywhere inside the character (any pixel with R-G >= 110 and B-G >= 110 is treated as background and erased).
IMPORTANT: no writing or kanji on any hat or plate — village identity comes only from the engraved symbol and the palette.
```
OUTPUT: may_giao_nv_portrait_chroma.png · ô 1254×1254 (máy nhồi về 200×200) · 1 khung · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.36 `ui_icon_app_x4.png` — P1 · can_ve

*Tham chiếu:* `TOT-naruto-logo`, `appicon_cu_dragonball`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the approved Leaf-Fire emblem exactly as in reference TOT-naruto-logo.png: a stylized leaf whose stem curls inward into a tight spiral, wrapped in sweeping orange-red flame strokes with a few teal-blue accent flicks and dark navy shadow strokes, bold clean dark outlines, cel-shaded, centered, occupying about 62% of the canvas width (inside the central safe circle so a round or squircle OS mask never cuts it), on a deep navy-charcoal background with a soft radial lift behind the emblem (slightly lighter navy at the center, near-black at the corners) and a subtle hand-painted paper texture. Read as a bold app icon at 48 px: one silhouette, high contrast orange on navy.
CANVAS: square 1024x1024 px master (deliver at 1024; the project downsizes to 512x512 for Google Play and Unity). Everything important inside the central 66% (adaptive-icon safe zone).
OUTPUT FORMAT: PNG, RGBA with the alpha channel fully opaque (every alpha = 255) — this image is NOT a sprite, it never goes through the chroma or zoom pipeline. No transparency, no magenta background.
DO NOT DRAW: any text, letters, numbers, kanji or fake writing (the project composites the real wordmark afterwards), watermark, border, rounded-corner mask, UI buttons, real-world logos, photographic elements, blurry or noisy AI artifacts.
DO NOT add a border, a rounded-corner frame, a bevel, a glossy lens flare, letters or the word NARUTO.
```
OUTPUT: ui_icon_app_x4.png · 1024×1024 (thu về 512×512) · 1 khung (ảnh tĩnh đục) · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.37 `ui_icon_app_hinh_x4.png` — P3 · can_ve

*Tham chiếu:* `TOT-naruto-logo`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: the approved Leaf-Fire emblem exactly as in reference TOT-naruto-logo.png: a stylized leaf whose stem curls inward into a tight spiral, wrapped in sweeping orange-red flame strokes with a few teal-blue accent flicks and dark navy shadow strokes, bold clean dark outlines, cel-shaded, centered, occupying about 60% of the canvas width, everything inside the central 66% safe circle; the rest of the canvas is fully TRANSPARENT (this is the Android adaptive-icon foreground layer).
CANVAS: 432x432 px (Android adaptive icon layer, 108 dp at 4x), PNG RGBA with a real transparent background.
DO NOT DRAW: background color, drop shadow, glow, border, text or letters; no semi-transparent halo around the emblem.
```
OUTPUT: ui_icon_app_hinh_x4.png · 432×432 · 1 khung (ảnh tĩnh đục) · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.38 `ui_icon_app_nen_x4.png` — P3 · can_ve

*Tham chiếu:* `TOT-naruto-logo`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: a plain background layer for the Android adaptive icon: deep navy-charcoal with a soft radial lift at the center (slightly lighter navy) fading to near-black at the edges, a very subtle hand-painted paper texture, NO emblem, NO shapes, NO text — the emblem lives on the separate foreground layer.
CANVAS: 432x432 px, edge-to-edge (the OS may shift and mask this layer, so the texture must be uniform enough to survive cropping).
OUTPUT FORMAT: PNG, RGBA with the alpha channel fully opaque (every alpha = 255) — this image is NOT a sprite, it never goes through the chroma or zoom pipeline. No transparency, no magenta background.
DO NOT DRAW: any text, letters, numbers, kanji or fake writing (the project composites the real wordmark afterwards), watermark, border, rounded-corner mask, UI buttons, real-world logos, photographic elements, blurry or noisy AI artifacts.
```
OUTPUT: ui_icon_app_nen_x4.png · 432×432 · 1 khung (ảnh tĩnh đục) · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.39 `ui_splash_doc_x4.png` — P2 · can_ve

*Tham chiếu:* `TOT-naruto-logo`, `logo_wordmark_client`, `concept-map1-langla`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: a vertical mobile launch screen for the game. Top third: the approved Leaf-Fire emblem exactly as in reference TOT-naruto-logo.png: a stylized leaf whose stem curls inward into a tight spiral, wrapped in sweeping orange-red flame strokes with a few teal-blue accent flicks and dark navy shadow strokes, bold clean dark outlines, cel-shaded, large, centered horizontally. Middle band from 45% to 60% of the height: EMPTY quiet dark navy area with no shapes (the project will composite the game wordmark here later). Bottom 40%: a painted silhouette skyline of the Leaf Village at dusk — a tall cliff with four blank carved faces reduced to abstract dark shapes (no recognizable real people), a red torii gate, tiled rooftops with a few warm lit windows, a large pale moon low on the left — warm orange-amber sky near the horizon fading upward into deep navy that becomes near-black at the very top. Match the palette and brush feel of the approved concept (concept-map1-langla.png). No characters.
CANVAS: 1080x1920 px, portrait. Keep every important element inside the central 80% width and away from the top/bottom 8% (notches, rounded corners, home indicator). Fully opaque.
OUTPUT FORMAT: PNG, RGBA with the alpha channel fully opaque (every alpha = 255) — this image is NOT a sprite, it never goes through the chroma or zoom pipeline. No transparency, no magenta background.
DO NOT DRAW: any text, letters, numbers, kanji or fake writing (the project composites the real wordmark afterwards), watermark, border, rounded-corner mask, UI buttons, real-world logos, photographic elements, blurry or noisy AI artifacts.
DO NOT DRAW: any text or lettering, characters or faces, UI buttons, loading bar, sparkles or lens flares.
```
OUTPUT: ui_splash_doc_x4.png · 1080×1920 · 1 khung (ảnh tĩnh đục) · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.40 `ui_splash_ngang_x4.png` — P2 · can_ve

*Tham chiếu:* `TOT-naruto-logo`, `logo_wordmark_client`, `concept-map1-langla`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: a horizontal launch screen for the game. Left-center: the approved Leaf-Fire emblem exactly as in reference TOT-naruto-logo.png: a stylized leaf whose stem curls inward into a tight spiral, wrapped in sweeping orange-red flame strokes with a few teal-blue accent flicks and dark navy shadow strokes, bold clean dark outlines, cel-shaded, large (about 40% of the height). To its right, from 45% to 85% of the width and vertically centered: an EMPTY quiet dark navy zone with no shapes (the game wordmark will be composited there). Along the bottom 30%: a painted silhouette skyline of the Leaf Village at dusk — cliff with four blank abstract carved shapes (no real faces), red torii gate, tiled rooftops with a few warm lit windows, a large pale moon at the far right — warm orange-amber sky near the horizon fading upward into deep navy, near-black at the top. Match the palette and brush feel of the approved concept (concept-map1-langla.png). No characters.
CANVAS: 1920x1080 px, landscape. Keep every important element inside the central 85% width. Fully opaque.
OUTPUT FORMAT: PNG, RGBA with the alpha channel fully opaque (every alpha = 255) — this image is NOT a sprite, it never goes through the chroma or zoom pipeline. No transparency, no magenta background.
DO NOT DRAW: any text, letters, numbers, kanji or fake writing (the project composites the real wordmark afterwards), watermark, border, rounded-corner mask, UI buttons, real-world logos, photographic elements, blurry or noisy AI artifacts.
DO NOT DRAW: any text or lettering, characters or faces, UI buttons, loading bar, sparkles or lens flares.
```
OUTPUT: ui_splash_ngang_x4.png · 1920×1080 · 1 khung (ảnh tĩnh đục) · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.41 `ui_bia_cuahang_x4.png` — P2 · can_ve

*Tham chiếu:* `player_idle_3f_chroma`, `master_iruka_x4_chroma`, `wolf_idle_2f_chroma`, `TOT-naruto-logo`, `concept-map1-langla`, `logo_wordmark_client`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: a store feature banner. Left third: the Leaf Village academy-student player character copied exactly from the reference (player_idle_3f_chroma.png — same boy, dark spiky hair, navy Leaf headband, olive flak vest, orange sash), three-quarter pose from the waist up, facing RIGHT toward the center with a confident grin, cel-shaded, clean outlines. Behind him slightly smaller and softer in value: the teacher Iruka (reference master_iruka_x4_chroma.png — brown ponytail, scar across the nose, olive vest) with arms crossed, and at the lower right a small grey forest wolf (reference wolf_idle_2f_chroma.png) snarling. Background: the Leaf Village at golden hour — cliff with four blank abstract carved shapes, red torii gate, rooftops — with five sweeping translucent chakra streaks crossing the sky from left to right in orange, sand-yellow, deep blue, earth-brown and bluish indigo. Keep the upper-center 40% of the canvas visually quiet and darker so the wordmark can be composited there.
CANVAS: 1024x500 px, landscape. Keep faces and the emblem-free quiet zone inside a 10% margin on every side. Fully opaque.
OUTPUT FORMAT: PNG, RGBA with the alpha channel fully opaque (every alpha = 255) — this image is NOT a sprite, it never goes through the chroma or zoom pipeline. No transparency, no magenta background.
DO NOT DRAW: any text, letters, numbers, kanji or fake writing (the project composites the real wordmark afterwards), watermark, border, rounded-corner mask, UI buttons, real-world logos, photographic elements, blurry or noisy AI artifacts.
DO NOT DRAW: any text, price tags, badges, UI, sparkles, lens flares, extra characters, magenta/pink anywhere.
```
OUTPUT: ui_bia_cuahang_x4.png · 1024×500 · 1 khung (ảnh tĩnh đục) · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.42 `ui_nen_dangnhap_tet_x4.png` — P3 · can_ve

*Tham chiếu:* `ui_nen_dangnhap_343`, `concept-map1-langla`, `TOT-naruto-logo`, `logo_wordmark_client`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: a seasonal login-screen background: the Leaf Village seen from a hillside road — tiled rooftops stepping down toward a red torii gate, the carved cliff with four blank abstract shapes in the distance (no real faces), a big sky occupying the upper half — reinterpreted for the season: Lunar New Year: strings of red paper lanterns hung along the street and the torii gate, blooming peach-pink and apricot-yellow blossom branches in the foreground, a few small painted fireworks bursting in the deep night sky, a light dusting of snow on the tiled rooftops, warm red-gold light spilling from doorways. Hand-painted, cel-shaded, clean outlines, matching the approved concept (concept-map1-langla.png). Keep the upper-center third of the sky calm and uncluttered so the wordmark can sit there. No characters.
CANVAS: 1024x576 px, landscape, fully opaque.
OUTPUT FORMAT: PNG, RGBA with the alpha channel fully opaque (every alpha = 255) — this image is NOT a sprite, it never goes through the chroma or zoom pipeline. No transparency, no magenta background.
DO NOT DRAW: any text, letters, numbers, kanji or fake writing (the project composites the real wordmark afterwards), watermark, border, rounded-corner mask, UI buttons, real-world logos, photographic elements, blurry or noisy AI artifacts.
DO NOT DRAW: any text or lettering, characters, UI buttons, real-world brand items, magenta/pink tones.
```
OUTPUT: ui_nen_dangnhap_tet_x4.png · 1024×576 · 1 khung (ảnh tĩnh đục) · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.43 `ui_nen_dangnhap_he_x4.png` — P3 · can_ve

*Tham chiếu:* `ui_nen_dangnhap_343`, `concept-map1-langla`, `TOT-naruto-logo`, `logo_wordmark_client`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: a seasonal login-screen background: the Leaf Village seen from a hillside road — tiled rooftops stepping down toward a red torii gate, the carved cliff with four blank abstract shapes in the distance (no real faces), a big sky occupying the upper half — reinterpreted for the season: high summer noon: a bright high blue sky with big white cumulus clouds, festival streamers fluttering between the rooftops, tall sunflowers along the road in the foreground, short dark hard shadows, heat-shimmer warm light on the tiled roofs and the red torii gate. Hand-painted, cel-shaded, clean outlines, matching the approved concept (concept-map1-langla.png). Keep the upper-center third of the sky calm and uncluttered so the wordmark can sit there. No characters.
CANVAS: 1024x576 px, landscape, fully opaque.
OUTPUT FORMAT: PNG, RGBA with the alpha channel fully opaque (every alpha = 255) — this image is NOT a sprite, it never goes through the chroma or zoom pipeline. No transparency, no magenta background.
DO NOT DRAW: any text, letters, numbers, kanji or fake writing (the project composites the real wordmark afterwards), watermark, border, rounded-corner mask, UI buttons, real-world logos, photographic elements, blurry or noisy AI artifacts.
DO NOT DRAW: any text or lettering, characters, UI buttons, real-world brand items, magenta/pink tones.
```
OUTPUT: ui_nen_dangnhap_he_x4.png · 1024×576 · 1 khung (ảnh tĩnh đục) · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

### 6.44 `ui_nen_dangnhap_trungthu_x4.png` — P3 · can_ve

*Tham chiếu:* `ui_nen_dangnhap_343`, `concept-map1-langla`, `TOT-naruto-logo`, `logo_wordmark_client`

```text
detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, rich shading (40+ colors), Naruto ninja aesthetic, NO flat vector shapes, NO background scenery, NO drop shadow on ground.
SUBJECT: a seasonal login-screen background: the Leaf Village seen from a hillside road — tiled rooftops stepping down toward a red torii gate, the carved cliff with four blank abstract shapes in the distance (no real faces), a big sky occupying the upper half — reinterpreted for the season: Mid-Autumn Festival: a huge golden full moon centered above the village, star-shaped and carp-shaped paper lanterns glowing warm orange along the eaves and the torii gate, a small festival table with round cakes in the foreground, a bluish indigo night sky (keep the purple bluish, no magenta). Hand-painted, cel-shaded, clean outlines, matching the approved concept (concept-map1-langla.png). Keep the upper-center third of the sky calm and uncluttered so the wordmark can sit there. No characters.
CANVAS: 1024x576 px, landscape, fully opaque.
OUTPUT FORMAT: PNG, RGBA with the alpha channel fully opaque (every alpha = 255) — this image is NOT a sprite, it never goes through the chroma or zoom pipeline. No transparency, no magenta background.
DO NOT DRAW: any text, letters, numbers, kanji or fake writing (the project composites the real wordmark afterwards), watermark, border, rounded-corner mask, UI buttons, real-world logos, photographic elements, blurry or noisy AI artifacts.
DO NOT DRAW: any text or lettering, characters, UI buttons, real-world brand items, magenta/pink tones.
```
OUTPUT: ui_nen_dangnhap_trungthu_x4.png · 1024×576 · 1 khung (ảnh tĩnh đục) · LƯU VÀO D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\

---

## 7. Tự kiểm và giao

### 7.1 Ba ảnh mẫu trước — DỪNG chờ duyệt

1. `nguoi_choi_la_portrait_chroma.png` (P1 — kiểm đồng bộ mặt với sprite đang chạy).
2. `ui_icon_app_x4.png` (P1 — kiểm thương hiệu, đọc được ở 48 px).
3. `nhan_gia_mat_na_p1_portrait_chroma.png` (P2 — kiểu khó nhất: mặt nạ + bộ học viên phải giống người chơi mà vẫn ra "kẻ lạ").

Khối copy-dán giao việc (giữ khuôn `map-1-lang-la\PROMPT-CHO-AGENT.md`):

```
Bạn là hoạ sĩ tạo ảnh cho game 2D Naruto. Đọc theo thứ tự, KHÔNG bỏ:
1. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\_KHUON-DOT-4.md  (phần A, B.10, C)
2. D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\_KHUON.md  (luật vàng, màu 5 làng)
3. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\09-chan-dung-thuong-hieu.md  (file này) — mục 3, 4, 6
4. Mở và NHÌN: anh-nguon\vertical-slice-v1\iruka_portrait_chroma.png · dot-4-ke-hoach-tong-the\anh-tham-chieu\THUOC-NHAN-VAT.png ·
   dot-3-toan-bo\anh-tham-chieu\TOT-naruto-logo.png · XAU-npc.png · XAU-boss-zabuza.png
VIỆC GIAI ĐOẠN 1: tạo ĐÚNG 3 ảnh mẫu ở §7.1 bằng prompt 6.1, 6.36, 6.6 — mỗi lần gen gắn kèm các ảnh ghi ở dòng "Tham chiếu" của prompt đó.
LUẬT KHÔNG ĐƯỢC VI PHẠM: chân dung = 1 ô 1254×1254 nền tím RGB(246,4,242) phẳng tuyệt đối, quay PHẢI, vai cắt phẳng chạm mép dưới,
  không chữ/khung/glow/bóng, không tím hồng trong hình; thương hiệu = PNG đục, không chữ, đúng cỡ ghi ở dòng OUTPUT.
LƯU VÀO: D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\   (tự tạo thư mục; đặt ĐÚNG tên file như dòng OUTPUT)
BÁO CÁO khi xong (nói thật): cỡ thực từng ảnh; ảnh nào nền tím có gradient/không phẳng; ảnh nào bạn KHÔNG CHẮC và vì sao;
  dán nguyên kết quả lệnh kiểm ở §7.2. DỪNG sau 3 ảnh, chờ duyệt.
```

### 7.2 Lệnh kiểm (có thật, chạy được hôm nay)

```powershell
# (a) CHÂN DUNG — tách nền, đo hộp, sinh alpha/contact-sheet/GIF; chạy từ D:\NRO_GoiBanGiao (cần Pillow + numpy)
cd D:\NRO_GoiBanGiao
python -m tools.artqa.artqa review --art-root D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\ --review D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\qa\review-09-chan-dung.json --reviewer "<tên bạn>"
#   ĐẠT khi trong JSON: raw_source_status = PASS (viền ≥ 85 % tím, spread ≤ 12) · issues = [] · soft_alpha_ratio ≤ 0.15 ·
#   content_boxes: đáy hộp ≈ 1254 (vai chạm đáy), hộp cao ≥ 0.85×1254, không chạm mép trái/phải. Mở qa\contact-sheets\*.png nhìn nền sáng/tối.

# (b) Xem thử đúng cỡ trong game (50×50 khung thoại, 26×26 HUD) từ bản alpha máy tách — nhìn còn ra mặt không
python -c "from PIL import Image; import sys; p=sys.argv[1]; im=Image.open(p); im.thumbnail((200,200)); im.save(p[:-4]+'_200.png'); im2=Image.open(p); im2.thumbnail((104,104)); im2.save(p[:-4]+'_104.png')" D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\qa\alpha\nguoi_choi_la_portrait_alpha.png

# (c) THƯƠNG HIỆU — sau khi tên mới đã nối vào kich-thuoc-anh.json (§5); chế độ nguon không kiểm cỡ, kiểm tên/alpha/ảnh phẳng/trống/nhẹ
python D:\NRO_GoiBanGiao\docs\my-thuat\KIEM-TRA-ANH.py D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\ nguon
#   Cỡ thật tự kiểm bằng mắt hoặc: python -c "from PIL import Image; print(Image.open(r'<file>').size, Image.open(r'<file>').mode)"
#   KHÔNG chạy TAO-4-MUC-ZOOM.py cho nhóm thương hiệu.
```

Lưu ý khi chạy (a): thư mục có cả ảnh thương hiệu (đục, không nền tím) → các file `ui_*` sẽ bị `artqa` báo FAIL nguồn — **đó là mong đợi**, chỉ đọc kết quả của `*_portrait_chroma.png`. Nếu muốn sạch báo cáo, để ảnh thương hiệu vào thư mục con `thuong-hieu\` rồi trỏ `--art-root` vào thư mục con `chan-dung\`.

### 7.3 Báo cáo nói thật (mẫu)

```
Đã tạo N/44. Cỡ thực: <tên> WxH mode RGBA ...
Nền tím: <tên> phẳng / có gradient (spread=?) ...   Quay phải: có/không ...   Vai chạm đáy: có/không (đáy hộp = ?/1254)
Không chắc: <tên> — vì ...
Kết quả artqa review / KIEM-TRA-ANH.py dán nguyên bên dưới:
```

### 7.4 Phía dự án sau khi duyệt (không phải việc người vẽ)

1. Cấp `Small<id>` ≥ 27085 cho từng chân dung, ghi vào `mapping-id.json` (append-only), nới `_single_frame_spec` (`task6b_closure.py:1004-1019`) và cổng 20 tên v1 (`artqa.py:31-52`); đổi `fit_to_canvas` NEAREST → LANCZOS cho chân dung; reviewer độc lập ký lại lock (`ASSET-ARTIFACT-REVIEW-016`).
2. Người chơi: cắt bản 104×104 sát mặt từ `nguoi_choi_la_portrait_chroma` thay nội dung `Small27084` (giữ id, `thegioi.go:167`); bản 200×200 nhận id mới, gắn khi máy chủ có gói thoại dùng chân dung người chơi.
3. NPC/boss: `catalog.json` `npcs[].portrait` trỏ id mới; boss chỉ có chân dung khi máy chủ chuyển thoại boss sang gói 38/32/-70 (đề xuất, ngoài brief).
4. Thương hiệu: thay `Assets\Resources\appicon.png` + Player Settings Icon (mọi nền tảng), đặt Splash (`androidSplashScreen`, iOS launch), đổi `productName` — làm trong `D:\NRO_Project` theo quy trình client, **không** chép ảnh vào Resources bằng tay.

---

## 8. Chưa cần vẽ đợt này

| Việc | Vì sao chưa |
|---|---|
| Chân dung Iruka | Đã có `iruka_portrait_chroma` → `Small27078`, QA PASS; lỗi hiển thị Iruka là dx/dy part (dữ liệu), không phải ảnh |
| Chân dung boss 4 làng (`boss_{cat,smu,da,may}_*` 12 tên trong bảng cỡ) | `boss.yaml` chỉ có 2 boss Làng Lá; boss vĩ thú không vẽ được bằng bộ xương người (KH6b) — chờ máy chủ chốt |
| Chân dung NPC ngoài 27 tên ở §5 | KH9 chưa có spec; thêm NPC = thêm master + chân dung cùng lúc ở brief 03 |
| Chốt `ui_logo_game` (emblem 512×256 trong lock) ↔ wordmark 599×136 (client) | Là **quyết định**, không phải vẽ (`00-KIEM-KE-DA-CO.md` §4 #2). Nếu chủ dự án muốn một lockup emblem+chữ mới thì mở brief riêng |
| `ui_nen_chonlang_x4` bản hand-painted | Màn chọn làng client chưa nạp nền; đợi 10-ui-hud |
| Icon theo mùa, icon thông báo (notification 96×96 trắng), ảnh chụp màn hình cửa hàng | Cần bản game chơi được trước; ảnh chụp lấy từ game thật, không vẽ |
| Chân dung theo trang phục đang mặc | Chỉ có nghĩa khi KH7 bản 2 "mặc đồ đổi hình" chạy (`BAT-DAU-TU-DAY.md` §5) |
| Video/trailer, ảnh mạng xã hội | Ngoài phạm vi mỹ thuật asset |

*Viết 06/09/2026 từ: `_KHUON-DOT-4.md`, `00-KIEM-KE-DA-CO.md`, `00-KIEM-KE-GAME-CAN.md`, `server/kho-noi-dung/boss.yaml`, `catalog.json`, `server/internal/game/{chat,npc,thegioi,hinhdang,khuvuc}.go`, client (chỉ đọc) `ChatPopup.cs`, `GameScr.cs`, `Panel.cs`, `Controller.cs`, `ModFunc.cs`, `LoginScr.cs`, `ProjectSettings.asset`, `Assets\Resources\appicon.png`; `tools/artqa/{artqa,task6b_closure}.py`; `tools/map-lang-la/tao-chan-dung-nguoi-choi.py`; `docs/my-thuat/{kich-thuoc-anh,manifest-tao-anh,mapping-id,unity-runtime-key-map}.json`; `dot-3-toan-bo/lang-*/npc.md`. Chưa mở Unity, chưa chạy game — mọi số đo là đọc file.*
