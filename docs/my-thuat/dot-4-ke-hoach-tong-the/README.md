# ĐỢT 4 — KẾ HOẠCH TỔNG THỂ MỸ THUẬT (đọc file này trước)

> **Cập nhật 06/09/2026 (trưa):** phần sửa lỗi "đầu tách thân" bằng dữ liệu (V1) **và** hình học NPC Iruka **đã xong, gộp `main`** (`75a64cc`) — mọi câu trong kế hoạch nói "đang hỏng / chưa tích hợp / chờ quyết" về hai việc đó đã được sửa lại. Kế hoạch vẽ dưới đây không đổi.
>
> Viết 06/09/2026 cho **chủ dự án** (không phải lập trình viên). Mọi con số dưới đây được đếm bằng máy từ 9 file manifest JSON trong thư mục này — muốn đếm lại, chạy `python D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\kiem-ke-hoach-dot-4.py`. Nguồn sự thật là từng `NN-*.json`; README chỉ tóm lại và xếp thứ tự.
>
> Ba file nền đã viết trước: [`_KHUON-DOT-4.md`](_KHUON-DOT-4.md) (luật kỹ thuật — game *thật sự* cắt/ghép ảnh thế nào), [`00-KIEM-KE-DA-CO.md`](00-KIEM-KE-DA-CO.md) (kho ảnh có gì), [`00-KIEM-KE-GAME-CAN.md`](00-KIEM-KE-GAME-CAN.md) (game cần gì). Khối copy-dán để giao cho agent tạo ảnh: [`PROMPT-CHO-AGENT.md`](PROMPT-CHO-AGENT.md).

**Đọc trong 2 phút:** đợt 4 có **9 hạng mục**, mỗi hạng mục một brief `.md` + một manifest `.json` cùng tên. Tổng cần vẽ **P1 97 · P2 364 · P3 170** ảnh/dải (cộng thô — chân dung bị đếm hai lần ở 4 brief, xem §6 D1); **đã có 250** dùng lại được. Việc đáng làm ngay: **10 ảnh giao diện** (không cần code), **2 dải người chơi vẽ lại** (sửa gốc lỗi "đầu tách thân"), **boss Làng Lá ra dáng boss**, **map Khu rừng**, **hiệu ứng 4 chiêu Làng Lá**. Mọi hạng mục làm theo một nhịp: **3 ảnh mẫu → lắp thử + chụp game → duyệt → làm loạt**.

---

## 1. Mục tiêu đợt 4 — nói đời thường

**Hôm nay game có gì:** một làng chơi được (Làng Lá "vertical slice"): 1 nhân vật học viên dùng chung cho cả 5 làng, 1 NPC Iruka (**đã sửa sáng 06/09** — cao bằng người chơi, ghép liền), 1 quái Sói Rừng sinh ở cả 6 map, 2 boss **mặc nhầm đồ** (pha 1 y hệt người chơi, pha 2 y hệt Iruka), map 1 hand-painted nhưng đất phẳng lì 91 %, 4 icon chiêu, HUD/nút/hộp thoại **vẫn là đồ Ngọc Rồng**, icon app vẫn là rồng Dragon Ball. Bốn làng còn lại vào được nhưng lát tile Dragon Ball và là ngõ cụt; Khu rừng (map 6) chưa vào được. Bấm chiêu chỉ thấy con số, không thấy lửa/nước/sét.

**Đợt 4 muốn ba việc:**

1. **Sửa gốc lỗi "đầu tách thân" bằng ảnh đúng thước** thay vì bù toán. Ngày 05/09 phát hiện 5 dải người chơi vẽ nhân vật **khác cỡ nhau** (chạy nhỏ hơn 12 %, đánh 17 %) và khom lưng → máy cắt mảnh thân lên tận trán → trong game hiện **hai khuôn mặt**. Phương án dữ liệu (V1) **đã làm xong và gộp `main` sáng 06/09** (commit `75a64cc`: 33/33 khung một mặt, đầu ngồi lên vai, đã chụp trong game) — nhưng vẫn phải trả giá **nhân vật thấp đi ~10 % khi chạy**, vì ảnh nguồn vẽ khom. Đợt 4 viết sẵn brief **vẽ lại đúng thước** (`THUOC-NHAN-VAT.png`), và **bảy luật 05/09** được nhúng vào mọi brief nhân vật/NPC/quái/boss để không lặp lại lỗi này với 4 làng.
2. **Làng Lá + Khu rừng chơi trọn** với máy chủ đang có: boss ra dáng boss, 3 quái rừng, NPC đứng đúng cỡ, hiệu ứng 4 chiêu, HUD Naruto, chân dung vẽ tay, icon 5 vật phẩm, map 1 hết đất phẳng.
3. **Bốn làng còn lại mỗi làng một gói trọn** (map → người chơi → NPC → quái → boss → trang bị → hiệu ứng) theo thứ tự lộ trình **Cát → Sương Mù → Đá → Mây**; cộng thương hiệu (icon app, màn chờ, bìa cửa hàng).

**Khác đợt 1–3 ở chỗ nào:** (a) chọn **kiểu file theo cách client ăn** — ảnh tĩnh alpha cho icon/UI/map, **dải chroma nền tím** cho mọi thứ chuyển động (`_KHUON-DOT-4.md` §A); vì thế 44 tên `quai_*` 96×96 và 18 tên `boss_*` 160×160 của đợt 3 **không vẽ nữa** (client không lắp được); (b) mỗi brief có **manifest JSON** để agent tạo ảnh tự chạy và **tự lưu** đúng thư mục; (c) mọi thứ chưa được máy chủ chốt đều ghi **"đề xuất"** — không brief nào giả vờ đã chốt.

---

## 2. Bảng 9 hạng mục

Số = **file cần vẽ** (không tính `da_co`). Cột "đã có" = ảnh dùng lại được theo cách brief đó đếm (gồm cả ảnh nằm ngoài manifest như 51 ảnh map 1, 100 icon `vp_*` lô 343).

| # | Hạng mục | Brief | Manifest | P1 | P2 | P3 | Đã có | Kiểu file chính | LƯU VÀO (ảnh tạo ra) |
|---|---|---|---|---:|---:|---:|---:|---|---|
| 01 | Nhân vật người chơi | [`01-nhan-vat.md`](01-nhan-vat.md) | `01-nhan-vat.json` | 3 | 25 | 2 | 2 | dải chroma ô 724×724 | `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\01-nhan-vat\` |
| 02 | NPC | [`02-npc.md`](02-npc.md) | `02-npc.json` | 14 | 24 | 0 | 2 | master chroma 724×868 + chân dung alpha 800→200 | `…\anh-nguon\dot-4\02-npc\master\` · `…\02-npc\chan-dung\` |
| 03 | Quái | [`03-quai.md`](03-quai.md) | `03-quai.json` | 15 | 80 | 0 | 5 | 4 dải chroma ô 724 + 1 master 1254 mỗi con | `…\anh-nguon\dot-4\03-quai\` |
| 04 | Boss | [`04-boss.md`](04-boss.md) | `04-boss.json` | 10 | 24 | 0 | 0 | dải chroma ô 724 (đổi đầu) / ô 880 (boss người lớn) + chân dung | `…\anh-nguon\dot-4\04-boss\` |
| 05 | Map (Làng Lá còn thiếu · Khu rừng · 4 làng) | [`05-map.md`](05-map.md) | `05-map.json` | 34 | 131 | 0 | 51 | tĩnh: 4 nền đặc rộng 2976 + tile 96 + vật thể alpha bội 96 | `…\anh-nguon\dot-4\05-map\` |
| 06 | Icon vật phẩm · trang bị · vũ khí | [`06-vat-pham.md`](06-vat-pham.md) | `06-vat-pham.json` | 1 | 8 | 154 | 120 | tĩnh alpha 64×64 (vẽ 1024 rồi máy thu) | `…\anh-nguon\dot-4\06-vat-pham\` |
| 07 | Kỹ năng · hiệu ứng | [`07-ky-nang-hieu-ung.md`](07-ky-nang-hieu-ung.md) | `07-ky-nang-hieu-ung.json` | 8 | 21 | 4 | 6 | dải chroma ô 384 (quanh người) / 192 (đạn bay) | `…\anh-nguon\dot-4\07-ky-nang-hieu-ung\` |
| 08 | Giao diện (UI / HUD) | [`08-giao-dien.md`](08-giao-dien.md) | `08-giao-dien.json` | 10 | 15 | 5 | 63 | tĩnh alpha **đúng cỡ** (thay thẳng file client) | `…\anh-nguon\dot-4\08-giao-dien\` |
| 09 | Chân dung · thương hiệu | [`09-chan-dung-thuong-hieu.md`](09-chan-dung-thuong-hieu.md) | `09-chan-dung-thuong-hieu.json` | 2 | 36 | 5 | 1 | chân dung chroma 1 ô 1254 + thương hiệu PNG đục | `…\anh-nguon\dot-4\09-chan-dung-thuong-hieu\` |
| | **Tổng (cộng thô)** | | | **97** | **364** | **170** | **250** | | gốc: `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\` (đã tạo sẵn 9 thư mục + 2 thư mục con NPC, mỗi cái có `.gitkeep`) |

> ⚠️ **Cộng thô vì trùng chân dung:** chân dung người chơi (5), NPC (19) và boss (3) nằm **cả** trong brief chủ thể (01/02/04) **và** trong 09, với tên và kiểu file khác nhau — 27 dòng. Chốt vẽ mỗi chân dung **một lần** ở §6 D1 rồi số thật sẽ thấp hơn (P1 ≈ 86, P2 ≈ 344 nếu chốt theo 09).

---

## 3. THỨ TỰ LÀM đề nghị — theo giá trị và rủi ro

Nguyên tắc xếp: **cái người chơi thấy nhiều nhất và lắp được ngay** lên trước; cái **phụ thuộc máy chủ/công cụ chưa có** vẫn vẽ (để không chờ nhau) nhưng xếp sau và ghi rõ "lắp sau". Mỗi bước ghi **"làm xong thì trong game thấy gì"** để biết mình đang mua cái gì.

### 3.0 Bước 0 — chốt trước khi vẽ (1–2 buổi, phía dự án, không tốn tiền vẽ)

| Việc | Ai | Vì sao phải làm trước |
|---|---|---|
| **D1 ĐÃ CHỐT:** chân dung theo 09 | chủ dự án | 09 là nguồn duy nhất; 01/02/04 không tạo bộ mã chân dung thứ hai |
| **C1 ĐÃ CHỐT:** vẽ lại run + hit; run **5 khung** | chủ dự án + người viết pipeline | bước 2 tồn tại; closure kiểm cỡ đầu C1 trước khi lắp |
| **ĐÃ NỐI** tên mới vào cuối `D:\NRO_GoiBanGiao\docs\my-thuat\kich-thuoc-anh.json`: 165 map (05), 30 UI (08), 1 kunai (06), 17 thương hiệu + mã chân dung (09), **3** mã tra cứu boss đang dùng (04; 7 mã dự phòng bị D1 loại) | agent điều phối | registry có 216 tên canonical mới; test khóa tên/cỡ/trùng |
| Dọn 2 sprite id lạ `Small27083/27084` trong client (00-KIEM-KE-DA-CO §4 #6) | agent điều phối | id mới của đợt 4 nối từ 27083 |
| (song song với vẽ) **ĐÃ XONG:** nới cổng `artqa gate` / `task6b_closure` khỏi 20 tên v1; kiểm "cùng cỡ đầu" cho dải mới | agent code | run C1 lệch quá ±5 % so với idle bị từ chối (`_KHUON-DOT-4.md` D.2–D.3) |

### 3.1 Đợt 4A — P1: Làng Lá + Khu rừng chơi trọn (≈ 97 ảnh, làm theo thứ tự này)

| Bước | Làm gì | Số ảnh | Cần gì trước / rủi ro | **Làm xong thì trong game thấy gì** |
|---|---|---:|---|---|
| 1 | **08 Giao diện P1** — khung HUD, 2 ô chiêu, dải hộp thoại 9 ô, 2 nút chữ 3 mảnh, 2 ô menu NPC, 2 nút chat | 10 | Nối 30 tên vào bảng cỡ. **Không cần code** (chép đúng tên cũ vào `res\x2,x3,x4\mainimage`). Rủi ro thấp nhất của cả đợt | Mở game là **hết dấu Ngọc Rồng** ở HUD, hộp đăng nhập, nút bấm, menu Iruka. Người chơi nhìn thấy 100 % thời gian |
| 2 | **01 Nhân vật P1** — vẽ lại dải **chạy 5 khung** + dải **trúng đòn 3 khung đứng** | 2 | **C1 đã chốt**: V1 là đường cũ cho tới khi lắp ảnh mới; pipeline nhận run 5 khung, khóa cỡ đầu ±5% cho cả run/hit, và chân slot 13 dùng cast 2. Kiểm bằng `kiem-dai-nhan-vat.py` rồi `ghep-khung.py --kiem` (33/33 cf, `so-mat=1`) | Nhân vật **chạy không thấp đi**, không còn "hai khuôn mặt", **lúc tung chiêu không hiện khung nằm**. Đây là lỗi lộ nhất hôm nay |
| 3 | **06 Vật phẩm P1** — icon kunai vuông; gắn id cho 4 icon đã có | 1 | Nối 1 tên; cấp Small id + ký lock (phía dự án) | Mở túi / cửa hàng Iruka **thấy icon 5 món** thay ô trống. Cần cho nghiệm thu KH7 task 13 |
| 4 | **09 P1** — chân dung người chơi Lá + **icon app** | 2 | Chân dung: sửa `_single_frame_spec` để nhận tên mới. Icon: đổi trong Unity Player Settings | Ô xanh HUD có **mặt vẽ tay** thay đầu cắt từ sprite; **icon trên điện thoại là Naruto**, không còn rồng |
| 5 | **05 Map 1 Làng Lá còn thiếu** — 2 nền rộng 2976 + 4 tile đất có vân | 6 | Có công cụ đo (`do-map.py` 7 chỉ số) và đường lắp có sẵn (`lap-anh-dot-2.py`) → **biết đạt ngay** | Đất map 1 **có vân** (hết phẳng lì 91 %), núi Hokage **không lặp 1,9 lần/khung**, núi cao hơn mái nhà |
| 6 | **04 Boss Làng Lá P1** — Kẻ Lạ Mặt (đổi đầu), Mặt Nạ Vỡ (5 dải người lớn), Học Viên Đào Tẩu (đổi đầu) | 10 (7 nếu chân dung về 09) | **Chưa ai thử boss to**: ô 880, cùng cỡ đầu, dải bù dx/dy; cần part id ≥ 2192 (nới `PHASE_ONE_ADDITIONS`), pipeline nhân vật thứ hai (chưa có). **3 ảnh mẫu là thí nghiệm bắt buộc** trước khi vẽ loạt | Boss **ra dáng boss**: pha 1 đội mặt nạ trên thân học viên, pha 2 người lớn đồ đen-đỏ, đào tẩu đeo băng trán Sương Mù — không còn giống người chơi/Iruka |
| 7 | **07 Hiệu ứng Làng Lá P1** — 5 dải chiêu Lá + trúng đòn + chết/vòng sáng + hồi sinh | 8 | **Bộ đóng gói `EffectData` / blob effect–arrow chưa có, cần viết** (07 §3.4); effect_id chưa chốt. Vẽ trước được; 3 dải mẫu dùng để thử bộ đóng gói | Bấm Hào Hoả Cầu **thấy cầu lửa bay và nổ**, Y Liệu có vòng sáng xanh, trúng đòn có chớp, chết có vòng sáng dưới chân |
| 8 | **05 Khu rừng** (28) + **03 Quái rừng** (Lợn Rừng · Nhện Độc · Khỉ Rừng, 15) | 43 | **Map 6 chưa vào được** (`banDo.go` chỉ khai 5 map — C8); pipeline quái **chưa nhận con thứ hai** (một tỉ lệ chung, atlas nhiều hàng, nhiều mob). Vẽ trước, lắp sau; đặt sau bước 5 để dùng bài học map 1 | Đi cổng sang **Khu rừng riêng** (không nhà, cây cổ thụ, lối đá rêu) gặp **lợn, nhện, khỉ** thay vì con sói dùng chung |
| 9 | **02 NPC Làng Lá P1** — Hokage Đệ Tam, người giao NV, thầy học viện, thợ rèn, y tá, lính gác, ông Tazuna (Khu rừng) | 14 (7 nếu chân dung về 09) | **KH9 chưa có spec** — máy chủ chỉ nhận đúng 1 NPC; pipeline NPC qua `component_placement` chưa có. Phụ thuộc máy chủ nhiều nhất nên xếp cuối P1 | Làng có **người**: bấm Hokage, thợ rèn… mở menu; NPC **đúng cỡ người chơi** (không lùn như Iruka hôm nay) |

Sau bước 9, **Làng Lá + Khu rừng là một game nhỏ trọn**: đăng nhập → HUD Naruto → nhân vật đúng thước → NPC → mua đồ có icon → đánh quái rừng → boss ra dáng → chiêu có hiệu ứng.

### 3.2 Đợt 4B — P2: từng làng trọn gói, **Cát → Sương Mù → Đá → Mây**

Mỗi làng làm **theo đúng chuỗi này** (mỗi mắt xích 3 ảnh mẫu → duyệt → loạt), ≈ 76–78 ảnh/làng:

| Mắt xích | Brief | Số ảnh/làng | Làm xong thấy gì |
|---|---|---:|---|
| a. **Concept map** 1 bức (prompt ở 05 §6.0) → đo `do-map.py` → duyệt | 05 | 1 (không lắp) | Có "ảnh đích" để đo mọi ảnh sau; chốt màu nối của làng (hiện là **đề xuất**) |
| b. **Bộ map** 4 nền + tile + nhà + cây + vật nhỏ | 05 | Cát 31 · Sương Mù 33 · Đá 31 · Mây 32 | Vào làng **thấy làng đó** (tường thành/cầu nước/vách đá/núi mây) thay tile Dragon Ball |
| c. **Người chơi làng** 5 dải + chân dung | 01 (+09) | 6 | Nhân vật làng đó **mặc đồ làng mình** (cần máy chủ gửi outfit theo làng — `hinhdang.go` hiện bỏ qua `lang`) |
| d. **NPC** Kage · thầy · người giao NV kiêm cửa hàng | 02 (+09) | 6 | Có người để mua đồ / học chiêu ở làng đó |
| e. **Quái** 4 con (thú nhỏ · thú lớn · lính gác · golem) × 5 file | 03 | 20 | Mỗi vùng **một đàn quái riêng** |
| f. **Boss** dạng người 5 dải + chân dung | 04 (+09) | 6 | Boss làng: Kiếm Khách Gió / Thất Đao / Golem Đá / Song Đao Lôi |
| g. **Trang bị** áo + quần t1 | 06 | 2 | Cửa hàng làng bán đồ làng mình |
| h. **Hiệu ứng** 5 dải (cận chiến · đạn bay · nổ · hồi máu · tụ chakra) | 07 | 5 | 4 chiêu của làng **nhìn khác** 4 chiêu Làng Lá |
| i. Icon chiêu | (07 → lắp) | 0 | Đã vẽ đủ 16 icon 4 làng — chỉ cấp id |

P2 **không theo làng** (làm xen khi rảnh tay): 01 dải đánh Lá vẽ lại (1) · 05 Làng Lá 2 đề xuất + Khu rừng 2 tiền cảnh (4) · 07 kết ấn chung (1) · 08 15 UI chỉ hiện khi mở bảng/vào map/bản PC · 09 splash 2 + bìa cửa hàng 1 + 14 chân dung NPC Lá vẽ lại (chỉ nếu D1 theo 09).

### 3.3 Đợt 4C — P3: đẹp thêm, chỉ vẽ khi được bảo

06: 152 icon trang bị 4 làng + 2 vẽ lại Lá (**chờ KH8 chốt số cấp/bộ phận**) · 07: lên cấp, boss xuất hiện, 2 hiệu ứng dùng thuốc (chưa có gói máy chủ) · 08: 5 huy hiệu làng 192 (chờ màn chọn làng 5 ô) · 09: 2 lớp icon adaptive + 3 banner mùa · 01: 2 biến thể đầu nữ/tóc (cần part đầu mới + máy chủ mở `dauHopLe`).

---

## 4. Quy trình một hạng mục (lặp lại đúng nhịp này cho cả 9)

| Bước | Ai | Làm gì | Bằng chứng để đi tiếp |
|---|---|---|---|
| 1 | chủ dự án | Đọc mục 1–2 của brief (5 phút) để biết đang mua gì; chốt các mục "đề xuất" liên quan ở §6 | — |
| 2 | agent điều phối | Làm việc bước 0: nối tên vào `kich-thuoc-anh.json` (nếu brief có `ten_moi_can_them_kich_thuoc`); thư mục LƯU VÀO đã có sẵn | `python kiem-ke-hoach-dot-4.py` không báo thiếu |
| 3 | chủ dự án | Copy **khối GIAI ĐOẠN 1** của hạng mục trong [`PROMPT-CHO-AGENT.md`](PROMPT-CHO-AGENT.md), dán cho agent tạo ảnh (agent phải đọc được `D:\NRO_GoiBanGiao\docs\my-thuat\`) | agent báo cáo 3 ảnh mẫu + kết quả lệnh kiểm |
| 4 | agent tạo ảnh | Đọc `_KHUON-DOT-4.md` → `../dot-3-toan-bo/_KHUON.md` → brief → JSON; mở ảnh tham chiếu; tạo **đúng 3 ảnh mẫu**, gắn ảnh tham chiếu mỗi lần gen; **tự chạy lệnh kiểm** ghi trong khối; **báo cáo nói thật**; **dừng** | file đúng tên trong thư mục LƯU VÀO |
| 5 | agent điều phối | Chạy lại lệnh kiểm phía dự án (`KIEM-TRA-ANH.py` / `python -m tools.artqa.artqa review` / `kiem-dai-nhan-vat.py`), **lắp thử một máy** và **chụp game**: Unity mở `D:\NRO_Project`, dùng `python D:\NRO_GoiBanGiao\tools\unity-mcp\unity_mcp.py` (Play/Stop/đăng nhập/chụp — xem `tools\unity-mcp\README.md`), máy chủ bật bằng `D:\NRO_GoiBanGiao\server\CHOI-THU.bat`; với nhân vật/boss chạy thêm `ghep-khung.py --kiem` | ảnh chụp trong game + kết quả lệnh dán vào báo cáo |
| 6 | chủ dự án | **Duyệt bằng mắt trên ảnh chụp game** (không duyệt trên ảnh nguồn): đúng cỡ so người chơi? đúng phong cách? có bay/lún/hai mặt? Ưng → bước 7; không → sửa prompt/ảnh mẫu, quay lại bước 3 | một câu "duyệt" hoặc danh sách sửa |
| 7 | chủ dự án | Dán **khối GIAI ĐOẠN 2** → agent tạo ảnh làm loạt theo thứ tự ghi trong khối, cùng luật, cùng thư mục | báo cáo từng file |
| 8 | agent điều phối | Kiểm cả lô bằng lệnh; ảnh nào FAIL trả lại đúng file đó (không nhận "gần đạt") | 0 lỗi trong kết quả lệnh |
| 9 | agent code + reviewer | **Lắp**: `python -m assetpack validate/build`, cấp id (sprite 27083+, part 2192+, mob, effect), **reviewer độc lập ký lại `asset-lock.json`**, tăng `vsData`, `chep-sprite-client.py --that`; map: `lap-anh-dot-2.py` (map 1) — map mới cần công cụ tham số hoá (chưa có) | game chạy với ảnh mới, chụp lại |
| 10 | agent điều phối | Cập nhật `BAT-DAU-TU-DAY.md` + `NHAT-KY-AGENT.md`, commit trên nhánh riêng | — |

**Ba điều không được bỏ:** (1) không làm loạt trước khi 3 ảnh mẫu được duyệt **trên ảnh chụp game**; (2) không tin báo cáo "QA sạch" nếu không dán kết quả lệnh (lô 343 từng tự nhận sạch khi 62 ảnh sai cỡ); (3) không tự sửa chữ ký asset-lock.

---

## 5. Cách đọc file JSON manifest (cho tool tự động)

Mỗi hạng mục có một `NN-*.json` cạnh brief, schema **mở rộng từ** `D:\NRO_GoiBanGiao\docs\my-thuat\manifest-tao-anh.json`. Khoá chung của cả 9 file:

```text
meta                 { du_an, hang_muc, ngay, brief_md | brief | file_brief, luu_vao, mo_ta, tong_anh, theo_uu_tien | so_anh ... }
phong_cach_khoa      chuỗi master prompt tiếng Anh (một số file có bản riêng: _chan_dung, _thuong_hieu, luat_icon, khoi_dai_chroma)
mau_lang             bảng màu 5 làng (EN)
anh_tham_chieu       { khoá: "đường dẫn tuyệt đối" }   ← 9/9 file, 100 % đường dẫn tồn tại (đã kiểm 06/09)
anh[]                mỗi phần tử = MỘT file ảnh phải nộp:
   ten_file          tên đúng để lưu
   kieu              "tinh_alpha" | "dai_chroma"
   so_khung          số ô (dải) — con số này cũng nằm trong tên file `_<n>f`
   o_kich_thuoc      [w, h] một ô (dải); null với ảnh tĩnh
   kich_thuoc        [w, h] cả file phải đúng (dải: = [so_khung × ô rộng, ô cao])
   noi_dung          mô tả tiếng Việt (cho người)
   dung_tham_chieu   [khoá…] → tra anh_tham_chieu để lấy ảnh gắn kèm khi gen
   prompt_day_du     prompt tiếng Anh tự đứng — TRÙNG TỪNG CHỮ với khối ```text ở mục 6 của brief
                     (ngoại lệ: 6 mục `da_co` của 07 để chuỗi RỖNG và không có khối ở mục 6 — ảnh đã có, không vẽ)
   luu_vao           thư mục tuyệt đối để lưu (02-npc: khác nhau giữa master\ và chan-dung\)
   uu_tien           "P1" | "P2" | "P3"
   trang_thai        "can_ve" | "ve_lai" | "da_co"   ← da_co = KHÔNG vẽ (dù có prompt)
```

> ⚠️ **Cẩn thận khi đọc `kich_thuoc` của dòng `da_co`:** ba brief hiểu khác nhau — 03 ghi **cỡ nếu vẽ lại theo chuẩn mới** (4 dải `wolf_*` ghi ô 724 trong khi file thật là 1774×887 / 2172×724, xem 03 §2), 06 ghi **cỡ chuẩn trong `kich-thuoc-anh.json`** (nguồn vẽ to gấp 4), 02/09 ghi **cỡ file thật**. Muốn biết cỡ thật của ảnh đã có thì đo file, đừng tin ô này.

**Thuật toán cho tool:** với mỗi manifest → duyệt `anh[]` → **bỏ `trang_thai == "da_co"`** → làm 3 ảnh mẫu (`INDEX-DOT-4.json` → `anh_mau_3_truoc`, hoặc `meta.anh_mau` / `meta.mau_truoc` / `meta.giai_doan.1_mau` nếu có) → dừng chờ duyệt → rồi theo `uu_tien` P1 → P2 → P3. Mỗi ảnh: gửi `prompt_day_du` + các ảnh trong `dung_tham_chieu` tới AI vẽ → lưu `luu_vao` + `ten_file` → **đo lại cỡ thật** so với `kich_thuoc` (dải: chia hết cho `so_khung`, mỗi ô chia hết cho 4; nền tím RGB(246,4,242) ở 4 góc).

> **Khoá `anh_tham_chieu` có mục không ảnh nào trỏ tới** (19 khoá, ví dụ `GOC-master-hoc-vien`, `HIEN-TRANG-map1-trong-game`, `XAU-npc`): đó là **ảnh nền tảng của cả hạng mục** — mục 4 của brief bảo gắn kèm mọi lần gen. Tool chỉ đọc JSON thì cứ gắn kèm mọi khoá `anh_tham_chieu` không được `dung_tham_chieu` nào nhắc tới.

**Chỗ khác nhau giữa 9 file — tool phải đọc phòng thủ:**

| Khác | Ở đâu | Đọc thế nào |
|---|---|---|
| Tên khoá đường dẫn brief | `meta.brief_md` (01/02/06/07/08/09) · `meta.brief` (04) · `meta.file_brief` (05) | thử lần lượt |
| Thư mục lưu theo **từng ảnh** | 02-npc: master → `…\02-npc\master\`, chân dung → `…\02-npc\chan-dung\`; mục `da_co` trỏ `vertical-slice-v1` | luôn ưu tiên `anh[].luu_vao`, rơi về `meta.luu_vao` |
| Ảnh đã có nằm **ngoài** `anh[]` | 05: `da_co_map1[]` (51) · 08: `anh_da_co[]` (14 nhóm) | không vẽ, chỉ để biết |
| Concept (không lắp) | 05: `concept[]` (5 bức, `kieu: tham_chieu_khong_lap`) | vẽ trước khi vẽ bộ map làng đó |
| Tên mới phải nối vào `kich-thuoc-anh.json` | 05: `ten_moi_can_them_kich_thuoc` (165) · 08: `ten_moi_kich_thuoc` · 04/06/09: trong `meta` | phía dự án làm trước khi chạy `KIEM-TRA-ANH.py` |
| Khoá riêng theo loại | 07: `neo` (A/B/C) + `lap` (true/false) + `y_nghia_khung[]` · 03: `thu_tu_khung[]`, `cao_muc_tieu_x4`, `ma_quai` · 08: `che_do_lap`, `thay_cho_client`, `anh_cu_x2_chi_doc` (file client cũ, chỉ đọc), `cat_manh` · 09: `kich_thuoc_dich`, `ve_o_co`, `ma_tra_cuu` · 06: `server_mau_id` · 01: `khung[]`, `thay_the` | đọc nếu có, không bắt buộc |
| Cỡ vẽ ≠ cỡ chuẩn | 06 vẽ 1024×1024 rồi máy thu về 64 · 02 chân dung vẽ 800 → 200 · 09 icon app vẽ 1024 → 512 · **08 vẽ ĐÚNG số** (thay file client) · 05 vẽ đúng số (bội 96) | theo `kich_thuoc` + ghi chú `ve_o_co`/`kieu_file` của file đó |
| Nền | dải chroma: tím RGB(246,4,242) phẳng · tĩnh: alpha thật · riêng `nen_*`/`map_*_nen_*` (05) và thương hiệu (09): **PNG đục** | 05 có khoá `nen: "alpha" | "dac"` |

**Mục lục máy đọc cả đợt:** `INDEX-DOT-4.json` (cùng thư mục) — liệt kê 9 hạng mục: brief, manifest, `luu_vao`, số cần vẽ theo P, 3 ảnh mẫu, tên trùng giữa manifest, đường dẫn tham chiếu thiếu. Sinh lại bất cứ lúc nào bằng `python D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\kiem-ke-hoach-dot-4.py` (thêm `--khong-ghi` nếu chỉ muốn xem).

---

## 6. Những gì CHƯA CHỐT — và ai quyết

Đánh dấu **[I2]** = đã nêu trong `00-KIEM-KE-GAME-CAN.md` §7; **[mới]** = phát hiện khi viết 9 brief. Brief đều viết sẵn bản "nếu chốt theo hướng đề nghị" để gật là làm ngay.

| # | Điều chưa chốt | Ảnh hưởng | Ai quyết | Đề nghị của kế hoạch này |
|---|---|---|---|---|
| **D1** [đã chốt 06/09] | **Chân dung theo 09.** 09 là nguồn duy nhất 1 ô chroma 1254; tên `cd_` dự phòng của 04 không nối registry | hết 27 ảnh vẽ đôi / hai kiểu file | chủ dự án | Vẽ chân dung cùng master/dải idle để khớp mặt; khi dán khối 01/02/04 bỏ dòng chân dung |
| **C1** [đã chốt 06/09] | **Vẽ lại run + hit; run 5 khung.** V1 dữ liệu `75a64cc` vẫn là đường đang chạy cho tới khi lắp ảnh mới | bước 2 §3.1; closure nhận `player_run_5f_chroma.png`, gán cf 2→6, và chặn cỡ đầu lệch quá ±5 % | chủ dự án + pipeline | Giữ idle/cast; thân 14→16 dùng hit 0→2, chân không dùng hit |
| **C2** [I2] | Danh sách **quái** Khu rừng + 4 làng — `quai.yaml` chỉ có sói | 03 toàn bộ là đề xuất | máy chủ (`quai.yaml`) | Dùng 3 con rừng + 4 con/làng như 03 §5; mã quái = gốc tên `quai_*` |
| **C3** [I2] | **Boss 4 làng**; boss vĩ thú cần code Mob | 04 P2 | máy chủ (`boss.yaml`) + KH6b | 1 boss dạng người/làng như 04; vĩ thú hoãn |
| **C4** [I2] | **NPC** ngoài Iruka, vai trò, menu, thoại — KH9 chưa có spec | 02 toàn bộ là đề xuất | KH9 | vẽ 7 NPC Lá theo 02 để KH9 không phải chờ vẽ |
| **C5** [I2] | **Danh mục vật phẩm** 4 làng, số cấp (8 hay 3), ô thứ 5 là băng trán hay nhẫn | 06 P3 152 icon | KH7/KH8 | chỉ vẽ áo+quần t1 (P2), phần còn lại chờ |
| **C6** [I2] | **Đường lắp hiệu ứng**: 07 đã đọc code — có 2 đường (Effect/EffectData và bảng efs/arrs của gói −87) nhưng **bộ đóng gói chưa có**, effect_id/arrow_id chưa chốt | 07 vẽ được, lắp chưa | agent code + KH10 | vẽ 3 dải mẫu Lá để viết bộ đóng gói song song |
| **C7** [I2 → 09 đã trả lời một phần] | **Chân dung boss có hiện không**: 09 đọc code thấy thoại boss đi gói 44 = bong bóng trên đầu, **không có avatar** | 7 chân dung boss của 04 (P1) | chủ dự án xác nhận | hạ chân dung boss xuống **P2/P3** (như 09), chỉ vẽ khi máy chủ đổi gói thoại |
| **C8** [I2] | **Khu rừng (map 6)** chưa vào được (`banDo.go` chỉ 5 map) | 28 ảnh map + 15 file quái rừng lắp sau | máy chủ | vẽ trước, lắp khi nối map 6 |
| **C9** [I2 → đã chốt] | Ô chung 724 cho dải 4 khung (2896×724) | — | `_KHUON-DOT-4.md` | **02 đề xuất thêm ô 724×868** (NPC đội nón) và **04 đề xuất ô 880×880** (boss người lớn) — pipeline phải chấp nhận ô khác 724; chưa chạy thử |
| **C10** [I2] | Mọi id mới (sprite 27083+, part 2192+, mob, effect) phải **ký lại `asset-lock.json`** | mọi hạng mục lắp | reviewer độc lập | mỹ thuật vẽ trước, lắp sau |
| **D2** [mới] | Boss pha 1 **đội mặt nạ trắng** hay chỉ trùm khăn (boss.yaml nói "tưởng học viên bình thường") | 04 bước 6 | chủ dự án | mặt nạ (đúng tên "Nhân Gia Mặt Nạ") |
| **D3** [mới] | **Logo**: wordmark 599×136 đang chạy hay emblem 512×256 đang ký trong lock | 08/09 | chủ dự án | không vẽ logo mới; icon/splash dùng emblem, chữ dán wordmark |
| **D4** [mới] | 08 nâng 10 UI từ P3 (I2) lên **P1** | thứ tự §3.1 | chủ dự án | giữ P1 — không cần code, thấy 100 % thời gian |
| **D5** [mới] | Trang phục người chơi 4 làng là **đề xuất** (máy chủ chưa gửi outfit theo làng) | 01 P2 | máy chủ (`hinhdang.go`) | vẽ bộ "tân binh" theo brief đợt 3 |
| **D6** [mới] | Kunai: **vẽ icon vuông mới** hay **đệm ảnh 96×48 cũ bằng công cụ** | 06 P1 | chủ dự án | vẽ mới (1 ảnh, rẻ) |
| **D7** [mới] | **Màu nối** 5 map mới (rừng/cát/sương/đá/mây) | 05 P1/P2 | chốt sau concept | dùng số đề xuất trong 05 §3.3, sửa prompt nền khi concept khác |
| **D8** [mới] | Công cụ **chưa có, cần viết** (không cản vẽ, cản lắp): bộ đóng gói EffectData/blob effect · pipeline NPC (`component_placement`) · pipeline quái (một tỉ lệ, atlas nhiều hàng, nhiều mob, bảng ánh xạ 12 frame) · `sinh-lai-sprite` cho nhân vật thứ hai · `lap-anh-map.py`/`kiem-anh-map.py` tham số hoá · script cắt 9 ô/3 mảnh UI | bước 9 §4 | agent code | cổng artqa và kiểm cỡ đầu C1 đã có; làm phần còn lại song song với vẽ |

---

## 7. Mục lục file trong thư mục này

| File | Là gì | Đọc khi |
|---|---|---|
| `README.md` | Kế hoạch tổng thể (file này) | bắt đầu |
| `PROMPT-CHO-AGENT.md` | 9 × 2 khối copy-dán giao việc cho agent tạo ảnh (giai đoạn 1 = 3 ảnh mẫu · giai đoạn 2 = làm nốt) | mỗi lần giao việc |
| `_KHUON-DOT-4.md` | Spec kỹ thuật đọc từ code: 2 kiểu file, 10 bảng spec, 7 luật 05/09, lệnh nghiệm thu có thật, luật tên/ID | trước khi vẽ bất cứ gì |
| `00-KIEM-KE-DA-CO.md` | Kho ảnh có gì (343 + 20 dải v1 + 51 map 1 + client thật), 6 chỗ lệch lock ↔ client | khi hỏi "cái này có chưa" |
| `00-KIEM-KE-GAME-CAN.md` | Game cần gì theo P1/P2/P3, 7 kiểu client dùng ảnh, 10 điều chưa chốt | khi hỏi "vì sao cần" |
| `01-nhan-vat.md` / `.json` | Người chơi: 32 file (P1 3 · P2 25 · P3 2 · đã có 2) | |
| `02-npc.md` / `.json` | NPC: 40 mục (P1 14 · P2 24 · đã có 2) | |
| `03-quai.md` / `.json` | Quái: 100 file (P1 15 · P2 80 · đã có 5) | |
| `04-boss.md` / `.json` | Boss: 34 file (P1 10 · P2 24) | |
| `05-map.md` / `.json` | Map: 165 ảnh (P1 34 · P2 131) + 5 concept + 51 đã có | |
| `06-vat-pham.md` / `.json` | Icon vật phẩm: 168 mục (P1 1 · P2 8 · P3 154 · đã có 5 trong manifest + 113 ngoài; 2 icon lô 343 nữa nằm trong 154 P3 vì vẽ lại → tổng 120 icon `vp_`+`vk_` của lô 343) | |
| `07-ky-nang-hieu-ung.md` / `.json` | Hiệu ứng: 39 mục (P1 8 · P2 21 · P3 4 · đã có 6) | |
| `08-giao-dien.md` / `.json` | UI/HUD: 30 ảnh (P1 10 · P2 15 · P3 5) + 63 đã có | |
| `09-chan-dung-thuong-hieu.md` / `.json` | Chân dung + thương hiệu: 44 (P1 2 · P2 36 · P3 5 · đã có 1) | |
| `INDEX-DOT-4.json` | Mục lục máy đọc 9 hạng mục (sinh bởi script dưới) | tool tự động |
| `manifest-dot-4.json` | **Gộp cả 9 manifest thành một** — 657 mục (P1 97 · P2 364 · P3 170 · đã có 21 · concept 5), mã làng đã chuẩn hoá về `la\|cat\|smu\|da\|may\|rung\|chung` (mã gốc giữ ở `lang_goc`), 152 ảnh tham chiếu gộp một bảng | agent tạo ảnh chỉ muốn đọc MỘT file |
| `kiem-ke-hoach-dot-4.py` | Đếm lại kế hoạch bằng máy, kiểm đường dẫn tham chiếu + thư mục LƯU VÀO, tìm tên trùng, sinh INDEX | trước mỗi lần giao việc |
| `gop-manifest-dot-4.py` | Gộp 9 manifest thành `manifest-dot-4.json` (chuẩn hoá mã làng, thêm `hang_muc`/`stt_gop`) | chạy lại mỗi khi một manifest hạng mục đổi |
| `kiem-dai-nhan-vat.py` | Kiểm hình học dải chroma ô 724 (tên `_nf`, cỡ, viền tím, ô rỗng, lề, chân chạm đáy, đảo tím) — dùng cho 01/04 và master 02 | sau khi agent vẽ dải nhân vật |
| `anh-tham-chieu\THUOC-NHAN-VAT.png` | Thước cỡ đầu / đường chân người chơi (gắn kèm mọi lần gen nhân vật/NPC/quái/boss) | |
| `anh-tham-chieu\HIEN-TRANG-hud-cu-x2.png` | 34 file UI client cũ ghép một tấm có nhãn (cho 08) | |

> ⚠️ **Thư mục ảnh nằm NGOÀI git.** `.gitignore:83` bỏ qua cả `docs/my-thuat/anh-nguon/` — nên 9 thư mục `anh-nguon\dot-4\<hạng mục>\` (và cả 0,3 GB ảnh nguồn hiện có: 343 ảnh lô 1, 20 dải chroma v1, 51 ảnh map 1) **không có bản sao trong git**. Ảnh mới vẽ xong phải tự sao lưu (ổ ngoài / đám mây) trước khi dọn máy; chỉ ảnh đã lắp vào client mới được chép vào `docs/my-thuat/lang-la-client-assets/` và commit.

**File liên quan ngoài thư mục:** `D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\_KHUON.md` (luật vàng, bảng cỡ tĩnh, màu 5 làng) · `D:\NRO_GoiBanGiao\docs\my-thuat\KHUON-DE-BAI-ANH-MAP.md` (khuôn 27–33 ảnh/map) · `D:\NRO_GoiBanGiao\docs\my-thuat\kich-thuoc-anh.json` (bảng cỡ 907 tên, append-only) · `D:\NRO_GoiBanGiao\docs\my-thuat\manifest-tao-anh.json` (schema gốc) · `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\PROMPT-CHO-AGENT.md` (mẫu khối giao việc đã dùng thật) · `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\` (20 dải chroma gốc) · `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\` (lô 343 ảnh tĩnh) · `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\` (nơi ảnh đợt 4 sẽ về).

---

*Kế hoạch này không sửa brief nào; mọi số đếm bằng `kiem-ke-hoach-dot-4.py` ngày 06/09/2026 (Python 3.14). Mọi đường dẫn `D:\…` trong README và PROMPT-CHO-AGENT đã kiểm `Test-Path` cùng ngày (100 tồn tại; 5 đường dẫn `review\*.json` là file đầu ra của lệnh review, chưa có là đúng); 9 thư mục LƯU VÀO + 2 thư mục con NPC đã tạo kèm `.gitkeep`. Không commit. Khi một brief đổi, chạy lại script và sửa bảng §2 theo số mới.*
