# BẮT ĐẦU TỪ ĐÂY

> **File này là bản đồ + nhật ký sống của dự án.** Ai mới vào — người hay agent AI —
> đọc file này trước tiên, đọc xong là hiểu: dự án là gì, đã làm được gì, còn gì
> chưa làm, đã vấp những lỗi nào để đừng vấp lại.
>
> **Quy tắc bắt buộc:** làm xong bất cứ việc gì đáng kể (một task, một kế hoạch,
> một bản vá, một phát hiện lỗi) thì **phải cập nhật file này ngay** rồi mới commit.
> Cách cập nhật xem mục [11. Quy ước cập nhật](#11-quy-ước-cập-nhật-file-này).

**Cập nhật lần cuối:** 08/09/2026 (**Task 3 UI Đợt 4 — đã lắp 15 PNG và có backup phục hồi**: popup/button x2 đã chụp thật trong Game View; HUD sau login còn BLOCKED vì máy chủ không hoạt động và lượt này không có tài khoản test được cấp. Không dùng production credential.)

**Trước đó:** 06/09/2026 (**Bước 0 phía repo cho Đợt 4 ĐÃ XONG trên nhánh `dot-4-thi-cong`**: D1 theo hạng mục 09, C1 vẽ lại run 5 khung + hit, nối 216 tên/cỡ canonical vào registry 907 tên; gate không còn khóa 20 tên V1 và closure chặn đầu run/hit lệch quá ±5 %. Không tạo/lắp ảnh, không đụng asset-lock hay client Unity.)

**Xa hơn:** 05/09/2026 (**Chơi thử thật qua Unity MCP → 8 lỗi code/dữ liệu đã vá trên `main`; khe đầu–thân CHẨN ĐOÁN XONG: thực chất là HAI KHUÔN MẶT (mảnh thân cắt 28 % chứa mặt) + 3 tầng lệch; phương án dữ liệu V1 đo 33/33 khung sạch, KHÔNG cần cài gì, chờ chủ dự án quyết** — bằng chứng `docs/my-thuat/review/khe-dau-than-2026-09-05/`)

---

## 1. Dự án này là gì

Xây **game Naruto online** (2D, client Unity, chơi trên điện thoại/PC) bằng cách
**viết lại từ đầu** phần máy chủ, dựa trên phần lõi của một server private
**Ngọc Rồng Online** (Dragon Ball) đã mua.

- **Không phải** vá hay đổi tên game cũ. Là **game mới hoàn toàn**, thương hiệu Naruto:
  5 làng (Lá, Cát, Sương Mù, Đá, Mây) thay 3 hành tinh Dragon Ball.
- Phần code Java/PHP cũ giữ lại **chỉ để tra cứu giao thức** — client Unity cũ nói
  chuyện bằng một giao thức nhị phân không có tài liệu, phải đọc ngược từ code cũ.
  Dấu vết Dragon Ball được **xoá dần** khỏi mọi thứ làm mới.
- Client Unity **giữ nguyên tạm thời**, dùng lại mã bản đồ / mã đồ hoạ cũ để có sẵn
  hình ảnh; mỹ thuật Naruto thay dần (xem `docs/my-thuat/`).

**Công nghệ mới:** máy chủ **Go** + **PostgreSQL** · web & panel **Next.js (App Router) + TypeScript + Tailwind**.
**Công nghệ cũ (chỉ tham chiếu):** Java 21 + MariaDB + PHP 8.

**Chủ dự án không phải lập trình viên** → mọi giải thích phải bằng **tiếng Việt đời
thường**, nói *tại sao* trước khi nói *làm thế nào*.

---

## 2. Bản đồ thư mục — cái nào sống, cái nào là di sản

| Thư mục | Là gì | Trạng thái |
|---|---|---|
| `server/` | **Máy chủ game mới, Go.** Toàn bộ công sức đang đổ vào đây | 🟢 Đang sống |
| `web/` | **Web + panel quản trị mới, Next.js.** Dùng chung CSDL với `server/` | 🟢 Đang sống |
| `docs/superpowers/` | Thiết kế (`specs/`), kế hoạch thi công (`plans/`), nghiên cứu giao thức (`research/`) | 🟢 Đang sống |
| `docs/my-thuat/` | Công cụ + prompt sinh ảnh Naruto (manifest, script thu nhỏ ảnh) | 🟢 Đang sống |
| `1-MAY-CHU-GAME/` | Server Java cũ (489 file `.java` + `zz1.jar` + `data/`) | 📖 Chỉ tra cứu |
| `2-WEB/htdocs/` | Web PHP cũ | 📖 Chỉ tra cứu |
| `3-CO-SO-DU-LIEU/z-moi-nhat.sql` | Dump MariaDB cũ, 48 bảng | 📖 Chỉ tra cứu |
| `4-BO-CAI-DAT/` | Script cài server cũ lên VPS | 📖 Chỉ tra cứu |
| `5-TAI-LIEU/` | Tài liệu kỹ thuật server cũ (rất chi tiết, hữu ích khi giải mã giao thức) | 📖 Chỉ tra cứu |
| `.superpowers/sdd/` | Sổ tiến độ chi tiết từng task + báo cáo review (**không vào git**) | 🟢 Nội bộ |
| `_BAN-LUU-*` | Ảnh chụp trạng thái cũ trước mỗi lần đổi lớn | 🧊 Đóng băng |

**Source client Unity KHÔNG nằm trong repo này** — nó ở `D:\NRO_Project` (nhiều GB).

---

## 3. Chạy và kiểm thử — lệnh cụ thể

Môi trường máy đang dùng: **Go ở `D:\Go`** · **PostgreSQL 16.4 ở `D:\PostgreSQL`**,
CSDL `nro` (chạy thật) và `nro_test` (kiểm thử) · Windows + PowerShell.

### ⚠️ Cổng PostgreSQL — TỰ KIỂM, đừng ghim con số

Cổng **không cố định giữa các lần khởi động**: `postgresql.conf` để mặc định `5432`,
nhưng bật bằng `pg_ctl ... -o "-p 5433"` thì nó nghe `5433`. Đã có agent mất thời gian
vì tin con số ghi sẵn trong tài liệu. **Luôn tự kiểm trước khi chạy:**

```powershell
Get-NetTCPConnection -State Listen | Where-Object { $_.LocalPort -in 5432,5433 } |
  Select-Object LocalPort
```

Không in ra gì tức là PostgreSQL **chưa chạy** (nó không chạy như dịch vụ Windows):

```powershell
D:\PostgreSQL\bin\pg_ctl.exe -D D:\PostgreSQL\data start
```

Dưới đây viết `CONG` — thay bằng số vừa tìm được. Chi tiết: `server/README.md` mục 0.

### 🎮 Chơi thử thật bằng client Unity — đường ngắn nhất (05/09/2026)

1. **Bấm đúp `server\CHOI-THU.bat`** — nó bật PostgreSQL nếu chưa chạy, build lại `nroserver.exe`,
   chạy máy chủ (cổng 14445) và **tự áp migration còn thiếu** vào CSDL `nro` lúc khởi động
   (05/09 điều phối kiểm: `nro` đang thiếu 0004/0005, máy chủ tự áp xong trong 0,2 giây). **Để cửa sổ đó mở.**
2. Mở Unity **`D:\Unity\2022.3.62f2\Editor\Unity.exe`** → project `D:\NRO_Project` → bấm **Play**.
   Client đã trỏ sẵn `127.0.0.1:14445`. Nếu báo *connect fail* là máy chủ chưa chạy, không phải lỗi code.
3. Đăng nhập **`naruto` / `1`** (đã có nhân vật "narruto", Làng Lá, map 1) hoặc **`thutest` / `1`** (chưa có
   nhân vật → sẽ qua màn tạo). Tạo thêm tài khoản: web `npm run dev` → `http://localhost:3000/dang-ky`,
   hoặc SQL ở `server/README.md` mục 2.
4. Tắt: **Ctrl+C** trong cửa sổ máy chủ (lưu nhân vật rồi mới tắt — đó chính là bài nghiệm thu KH11).
5. Có gì lạ → chép nguyên dòng log trong cửa sổ máy chủ cho agent; log cũng in ra màn hình, không ghi file.

Danh sách việc phải thử tay (KH4/6/6b/7/11) ở §5 mục A của lần soát 05/09 trong `NHAT-KY-AGENT.md`.

```powershell
# --- Máy chủ game ---
cd D:\NRO_GoiBanGiao\server
$env:NRO_DSN = "postgres://postgres:nrodev123@127.0.0.1:CONG/nro"
D:\Go\bin\go.exe run ./cmd/nroserver          # lắng nghe :14445

# --- Client giả để thử, không cần mở Unity ---
D:\Go\bin\go.exe run ./cmd/thuclient -dia-chi 127.0.0.1:14445 -ten thutest -mat-khau 1

# --- Chạy toàn bộ test (PHẢI xanh trước khi commit) ---
$env:NRO_TEST_DSN = "postgres://postgres:nrodev123@127.0.0.1:CONG/nro_test"
D:\Go\bin\go.exe test ./... -count=1
D:\Go\bin\gofmt.exe -l internal/ cmd/         # phải không in gì (đã chuẩn hoá LF, §9)
D:\Go\bin\go.exe vet ./...

# --- Web ---
cd D:\NRO_GoiBanGiao\web
npm run dev                                   # http://localhost:3000
npm run build                                 # kiểm biên dịch trước khi commit
```

Chi tiết hơn: `server/README.md` và `web/README.md`.

---

## 4. Đã làm được những gì

Lộ trình đầy đủ: `docs/superpowers/plans/LO-TRINH.md`.
Mỗi kế hoạch đều đi qua: thiết kế → kế hoạch → thi công từng task bằng subagent
(TDD) → review đối kháng → gộp vào `main`.

| # | Kế hoạch | Kết quả | Commit gộp |
|---|---|---|---|
| — | **Xoá vệt hồng trên nhân vật (chroma còn sót) — 06/09/2026** | Chơi thử thấy chấm hồng ở khe nách/thắt lưng và dưới nút khăn trán Iruka. Đo: **293 px** màu nền còn sót trên 37 mảnh × 4 zoom (x1 11 · x2 41 · x3 85 · x4 156). **Hai tầng cùng bỏ sót:** (1) `artqa.extract_alpha` chỉ flood nền magenta **chạm mép ảnh**, nên khe KÍN trong lòng hình giữ nguyên magenta đục (alpha 255); (2) `remove_large_chroma_islands` chỉ xoá đảo ≥ 598 px — và **chú thích trong code nói đảo 586 px là "accent bông tai" là SAI**: cắt ảnh nguồn ra xem thì đó là **khe giữa dải khăn trán và tóc Iruka**. Sửa: thêm `NGUON_NHAN_VAT` (5 dải người chơi + master Iruka) dùng ngưỡng khe 8 px + vá màu đốm < 8 px; đường closure v11 cũ giữ tham số mặc định nên **bất động từng byte** (kiểm: 20/20 ảnh nguồn khác không lệch). Tự kiểm trước khi xoá: histogram cho thấy trang phục **không có tím/hồng thật** (mọi pixel tím đều có kênh lục 1–5 = viền răng cưa của nền). Kết quả: **0 px chroma** ở cả 4 mức, `part.bin` **không đổi byte** (dx/dy giữ nguyên → `versions.data`/`vsData` giữ 28), 67 ảnh đổi, mapping chỉ đổi trường sha256. Hai agent kiểm đối kháng: đạt, nhưng bắt được **2 đột biến sống** — điều phối vá lưới (test tái sinh nối mã sinh với ảnh trong repo; `artqa.chroma_mask` làm một nguồn sự thật + bảng màu ghim tay) rồi tự áp lại 3 đột biến, cả 3 đều chết. Ảnh: `docs/my-thuat/review/khe-dau-than-2026-09-05/in-game/vet-hong-truoc-sau.png`. | `main` 06/09 |
| — | **Kế hoạch mỹ thuật ĐỢT 4 — 05–06/09/2026** | Trả lời câu "cần tạo những ảnh gì": 3 bản kiểm kê (đã có gì · game cần gì · spec kỹ thuật đọc từ code) + **9 brief** (nhân vật · NPC · quái · boss · map · vật phẩm · kỹ năng-hiệu ứng · giao diện · chân dung-thương hiệu), mỗi brief một `.md` cho người đọc và một `.json` manifest cho tool: **652 mục ảnh** (cần vẽ P1 97 · P2 364 · P3 170; đã có dùng lại 250). Mỗi ảnh có **prompt tiếng Anh tự đứng**, dòng OUTPUT (tên · cỡ · số khung) và **thư mục LƯU VÀO tuyệt đối** — agent tạo ảnh đọc file là tự vẽ và tự lưu. `PROMPT-CHO-AGENT.md` có 9 × 2 khối copy-dán (3 ảnh mẫu → duyệt → làm loạt). Điểm khác đợt 1–3: chọn **kiểu file theo cách client thật ăn ảnh** (tĩnh alpha cho icon/UI/map; **dải chroma nền tím** cho mọi thứ chuyển động) nên 44 tên `quai_*` và 18 tên `boss_*` của đợt 3 bị loại (client không lắp được); **7 luật 05/09** (cùng cỡ đầu, chân chạm đáy ô, đứng thẳng, nền tím phẳng, một ô một khung, quay phải, không vẽ nằm) nhúng vào mọi brief nhân vật/NPC/quái/boss + thước `anh-tham-chieu/THUOC-NHAN-VAT.png`. Điều phối tự kiểm: **631/631 prompt khớp từng chữ giữa MD và JSON**, 652/652 `luu_vao` tồn tại, 152/152 ảnh tham chiếu tồn tại, 0 lỗi cứng (`kiem-ke-hoach-dot-4.py`). Agent soát đối kháng đạt, tự sửa 9 điểm; điều phối xử nốt 9 điểm còn lại (cỡ 4 dải sói khai sai so file thật, hai bộ mã chân dung boss → chốt D1 theo brief 09, mã làng về `la|cat|smu|da|may`, ô chân dung 1254, giữ ô menu 230 vì thay thẳng file client). | (chỉ tài liệu, không đụng code) |
| — | **Đợt 4 — Bước 0 phía repo — 06/09/2026** | D1 chốt theo 09 nên chỉ nối 3 mã boss đang dùng (7 `cd_` dự phòng của 04 bị loại); C1 chốt vẽ lại run 5 khung + hit. Append-only `kich-thuoc-anh.json` thêm đúng **216** tên/cỡ từ manifest 04/05/06/08/09 (registry 907 tên). `artqa gate` nhận tập review không rỗng/không trùng thay vì 20 tên V1; closure ưu tiên run C1, giữ fallback tái lập V1 đã ký, và từ chối cỡ đầu run/hit lệch quá ±5 %. Test RED→GREEN + mutation nới 1,05→1,50 bị bắt; không tạo/lắp ảnh, không đổi lock/client. | `dot-4-thi-cong` / `HEAD` |
| — | **Đợt 4 — Task 2 UI, fix breaker round 5/5 — 06/09/2026** | Bịt nốt finding detector chỉ đếm transition ở biên: học cadence ngang/dọc, thử pha 2D, rồi đòi đúng màu + đúng pha ≥85%, coverage ≥10% và phủ ≥8/16 vùng trên toàn mặt ảnh. Fixture 120×96 viền `242/255` thật nhưng ruột `242` đặc bị từ chối trước khi đổi pixel; nền caro thật vẫn xoá, highlight `(244,244,244)` chạm mép giữ nguyên, trắng/backup lệch SHA vẫn reject. RED→GREEN **8/8**, mutation vô hiệu riêng ratio spatial bị test mới giết. Ba PNG final chạy idempotent `changed=False`, SHA không đổi; không chạm geometry/Unity/asset-lock. | `dot-4-thi-cong` / `HEAD` |
| — | **Đợt 4 — Task 3 UI, lắp thử Unity — 08/09/2026** | Xác minh SHA 3 source canonical; tạo backup phục hồi **15/15** ở `docs/my-thuat/lang-la-client-assets/mainimage-truoc-dot-4/`; sinh 12 zoom, stage/copy đúng 15 PNG x2–x4 rồi Unity refresh idle. Harness TDD **5/5**, mutation cổng backup bị bắt. Game View thật chứng minh popup liền mạch và nút Có/Không đọc được; `Small27083/27084` x1–x4 giữ nguyên hash, asset-lock không đổi. **Chưa đủ PASS toàn task:** máy chủ mất kết nối nên không vào HUD để kiểm HP/chakra/anchor; không đoán credential. | `dot-4-thi-cong` / `HEAD` |
| — | **Sửa khe đầu–thân nhân vật bằng dữ liệu (V1) + Iruka — 05–06/09/2026** | Gốc lỗi là **hai khuôn mặt**: mảnh thân cắt ở 28 % chiều cao khung chứa cả mặt, vẽ đè lên đầu (client vẽ đầu→chân→thân); cộng 5 dải nguồn khác tỉ lệ (chạy nhỏ hơn idle ~12 %, đánh ~17 %) và khung NẰM gán nhầm slot thân 15 (cf 20/31). Sửa toàn bộ bằng dữ liệu, không đụng client: `tools/artqa/nhan_vat_moc.py` (đo cằm/tâm mặt/6 thước cứng → tỉ lệ từng khung theo cỡ đầu), `_component_rect_v1` (cắt thân từ cằm − 2 px), bảng slot V1 (cf 2..6 = 5 khung chạy, slot 15 = cast0, bỏ 3 khung hit), `tools/artqa/nhan_vat_toi_uu.py` (MILP scipy giải dx/dy một lần cho 34 mảnh, ghim dy đầu −38, bàn chân −11 giữ nguyên), `REQUIRED_CLIENT_CF` sinh từ CharInfo (hết lệch một bậc). Iruka: `iruka_component_images` cùng cỡ đầu người chơi (r = 1,664, s = 0,226), dx/dy đặt theo cf 0 → hết đầu lửng +13/chân chìm +29. `sinh-lai-sprite.py` một lệnh sinh 37 mảnh ×4 zoom + `part.bin` + mapping + lock (`versions.data` 28, `vsData` 28); `ghep-khung.py` thêm cột `so mat` (bộ dò hai khuôn mặt) + cờ `--kiem`/`--cf` làm cửa nghiệm thu. **Đo lại (điều phối tự chạy):** người chơi x1 33/33 khe đầu–thân ∈ [−3,−1], thân–chân ∈ [−11,−2], bàn chân −11/−10, một mặt; x4 không khe dương; Iruka cf 0/1: −3/−6/−11, cao 62–63 (người chơi 61–62). pytest 256 xanh / 8 đỏ sẵn (test_vertical_slice_generator, thiếu nguồn 27078); `go test ./...` + `-race` xanh. 3 agent kiểm đối kháng tự chạy lại mọi cửa + 20 đột biến tự áp đều bị bắt. Chép 148 sprite vào client (`chep-sprite-client.py --that`, sao lưu bản cũ chỉ ghi một lần), máy chủ build lại (part.bin nhúng), chụp trong game qua MCP: đứng/chạy phải/chạy trái/cạnh Iruka — một mặt, đầu ngồi vai. **Giá đã chấp nhận (chủ dự án chốt tối 05/09):** nhân vật thấp đi ~10 % khi chạy (53–55 so 61–62 px), thân trượt ≤ 5 px ở cf 3/14/16, đầu idle hạ 1 px. | nhánh `sua-khe-dau-than-v1` → `main` 06/09 |
| 1 | Lõi giao thức + đăng nhập | Bắt tay, mã hoá XOR, khung gói, PBKDF2, đăng nhập | `6529d35` |
| 2 | Nhân vật: tạo / lưu / tải | Tạo nhân vật, chọn làng, lưu PostgreSQL | `0816c89` |
| 3 | Đồng bộ dữ liệu tĩnh | 4 catalog (data/map/skill/item), **5 làng Naruto thay 3 hành tinh** | `7a85431` |
| 4 | Vào thế giới | Nhân vật hiện ra đứng trong làng (tileset + ME_INFO + MAP INFO) | `aef5f71` |
| 5 | Di chuyển + mô hình actor | Nhiều người **thấy nhau**, đi lại, chống dịch chuyển tức thời | `e1ddc44` |
| 6 | **Chiến đấu, chiêu thức, quái** | 20 chiêu Naruto (4/làng) + sát thương ngũ hành (×1,25/×1,00/×0,80); đánh quái, quái đánh lại + hồi sinh, **PvP**, chết & hồi sinh; nhiều người cùng khu vực **thấy toàn bộ diễn biến** | `162bbe9` |
| 11 | **Lưu trạng thái sống** | Trước đó máy chủ **không lưu gì** sau lúc tạo nhân vật. Nay lưu **3 thời điểm**: định kỳ 60 giây · khi người chơi rời · khi tắt máy chủ. Kiến trúc: actor **chụp ảnh** → hàng đợi gộp → **thợ ghi riêng** chạm đĩa (actor không bao giờ chờ ổ đĩa). Kèm **một tài khoản một phiên** (đá phiên cũ → chờ lưu xong → mới đọc) | `6ea467f` |
| W1 | Web nền tảng | Đăng ký, đăng nhập, phiên cookie, dùng chung CSDL với game | `af5a2e4` |
| W2 | Web tính năng | Trang cá nhân, đổi mật khẩu, bảng xếp hạng, panel admin (chặn quyền 3 lớp) | trên `main` |
| VS-6A | QA ảnh nguồn vertical slice Lang La | Fix round 2: semantic one-shot GIF, raw provenance, canonical evidence và gate fail-closed; raw **0/20 PASS**, derived **20/20 PASS**, gate exit 0 | nhánh `kh-vertical-slice` |
| VS-6B | Bind ảnh đã duyệt vào manifest Làng Lá | Exact **20/20 provenance** fail-closed; closure v11 có 42 sprite `27041..27082`, HUD 25/50/100 và Iruka sáng/tối được reviewer độc lập ký SHA; mutation/full test/build A/B đều xanh | `bf91ee8`, `545b9c6`, `0fd6fc5` |
| VS-7 | **Build artifact + deploy staging theo lock** | Artifact server/client bất biến theo lock 21; bắt buộc previous lock + protected Unity root; canonical resource đủ x1–x4 cho small/skill/mob/background/HUD; Unity preflight thật marker `RESOURCE-PREFLIGHT-OK`; full Python **123/123**, Go toàn repo + vet xanh | `2a77a43`, `9627680`, `56fa302` |
| VS-8 | **Catalog, NPC Iruka, Sói và collision Làng Lá** | Đã vá audit Important: grid 5×17 tile 24×24 fail-closed; map 1 lấy footprint 120×408 làm biên authoritative, map 6 ghim là template chưa playable. Full Python **124/124**, Go game **-count=1**, focused collision/net **-count=5**, vet xanh; **independent audit APPROVE** | `de1f0c6`…`b82d950` |
| VS-9 | **Telemetry Unity, server control và runner** | Telemetry JSONL loopback có token/command ID/sequence và main-thread dispatch; session cũ không thể lọt lệnh sang client mới. Server control chỉ mở bằng env rõ ràng, bind loopback và tắt qua đúng `donDepKhiTat`. Fixture chỉ dùng `nro_test`, đúng 2 account/run và cleanup fail-safe. Runner có preflight/timeout/evidence SHA. Unity contract thật **50 assertion**, Go toàn repo + vet, PowerShell self-test và control smoke DB thật đều xanh; **independent audit APPROVE** | `01af875` |
| VS-10 | **Clean full E2E và warm-upgrade smoke** | Canonical A đã harden thủ công theo audit B, không cherry-pick: 24 case qua runner boundary; tách server/asset/applied lock, file zero-byte/SHA/array rỗng fail-closed, Sói `source_id=1`, HP `137/150`, clean+warm cùng PASS, sequence per-client. Bốn mutation đều bị bắt. Asset blocker đã được xử đúng provenance: build A/B **669/669 byte-identical**, publish release Task 8 lock `54a638...`, pair `afd03...`, version `22/21/22/21`, Unity resource preflight thật PASS. **Gameplay thật vẫn chưa chạy** | `e6582ce`, `0eb4c60`, `f56341f` |
| A | **Lưới chắn panic (nợ hạ tầng A)** | Trước đó **KHÔNG goroutine nào có `recover()`** — một panic bất kỳ đá văng TOÀN BỘ người chơi. Nay **6/6 goroutine** qua lưới chung `antoan.ChayCoLuoi` (recover + log kèm `stack`, không nuốt im lặng): actor khu vực **chạy tiếp** sau panic (bọc từng sự kiện) + **van an toàn** dừng riêng khu vực panic đại dẳng (≥100 nhịp liên tiếp, reset khi có nhịp trót); mỗi kết nối (`PhucVu`+`vongWriter`) **tự đóng**; hệ lưu (`luudinhky`+`thogi`+mỗi bản ghi) **bỏ đúng việc lỗi, giữ vòng sống**. 6 task TDD + kiểm thử đột biến (gỡ lưới → đỏ vì panic lan ra) | **gộp `main` `3c0a539`** (đã soi đối kháng) |
| 6c | **KH6c — Stress test 20 client** (nợ hạ tầng) | Chứng minh §6.10b **vững dưới tải**: 20 người dồn **~5 triệu** sự kiện/2s vào 1 khu vực — actor **không treo lần nào** (probe luôn <100ms, `probeTimeout=0`), **giữ nhịp** (`soNhipTre=0`), **gạch đúng 3 client chậm**, ~1,7 triệu gói rơi **đều được đếm**. Thêm 2 chỉ số quan sát cho server: `soNhipTre` + `treNhipLauNhatMs` (tách `danhDauNhip` để test khóa được đường nối). 2 task TDD, soi đối kháng độc lập (APPROVE-WITH-NITS) → đã vá. Nhân tiện phát hiện **§6.10b đã vá sẵn từ `a135979`**, sổ sửa bỏ báo động giả | **gộp `main` `322b140`** (đã soi đối kháng) |
| B/C/D | **Nợ hạ tầng: Migration + CI + Backup** | **C** — server **tự áp migration** khi khởi động (`store/migrate.go`: `go:embed` + bảng `schema_migrations` + **advisory lock** chống đua 2 server + **fail-closed** rollback; tách `fs.FS` để test được). **B** — `kiem-tra-tat-ca.ps1` **một lệnh** (build/vet/gofmt/go test/pytest) + git hook `pre-push` + GitHub Actions sẵn (chưa có remote). **D** — `sao-luu-db.ps1` pg_dump định kỳ **chống tự huỷ** (dump `.partial`→đổi tên khi xong, retention lọc >0 byte, ghi log) + Scheduled Task tự-host. Soi đối kháng **REQUEST-CHANGES** → vá 2 lỗ hổng **nặng** (fail-closed không có test → đột biến fail-open sống; backup tự huỷ khi task fail) + 3 vừa; 2 nhỏ (test_task6b worktree, CRLF `.gitattributes`) tách chip riêng. CI full xanh, smoke server thật | **gộp `main` `a51b601`** (đã soi đối kháng + vá) |
| G1 | **Cổng rời Làng Lá + bóng đổ (bg-item)** | Gói MAP INFO đang gửi **số 0** cho hai thứ client đã hỗ trợ sẵn, hậu quả: người chơi **không có cách nào rời map 1**, và **8 ảnh bóng đổ vẽ xong nằm không dùng được**. Nay: **cổng torii** ở giữa map 1 (`isEnter=true` — bấm mới đi, đi ngang không bị hút) dẫn sang Làng Cát, Làng Cát có cổng **về lại** (không nhốt người chơi); **11 bóng đổ** đặt dưới chân công trình, ảnh cắm đủ x1–x4 vào `Resources/res/x*/mapBackGround/500..507`. Kèm hai gói mới: **`-31`** bảng mẫu bg-item (thiếu nó thì bóng đổ bị client **bỏ qua im lặng**) và **`-23`** đổi map (thiếu nó thì bấm cổng là nhân vật **đứng hình vĩnh viễn**). Vòng soi đối kháng độc lập sau đó vá thêm 3 lỗi và bịt 3/4 lỗ hổng test — kiểm thử đột biến **21/22 bị giết** | nhánh `worktree-agent-a6c5f42987d2a805f` (chờ gộp) |
| F9 | **Dọn nợ nhiễu CRLF của gofmt (chip F9)** | `.gitattributes` đã ép `* text=auto eol=lf` (blob toàn LF); nay **bỏ danh sách hardcode 9 file** trong `kiem-tra-tat-ca.ps1` → gofmt soi **mọi** file, **thêm bước gofmt vào CI** (`ci.yml` job go), thêm ngoại lệ `*.bat`/`*.cmd = eol=crlf` bảo vệ script Windows. Đột biến: thả 1 file `.go` sai format → gofmt bắt đỏ. `kiem-tra-tat-ca.ps1` full **XANH** (gofmt trọn bộ, go test 8 gói, pytest 118) | nhánh `claude/blissful-haibt-1a23fa` (chờ gộp) |

> **VERTICAL SLICE (VS-6A…VS-10, tức Task 0–9) ĐÃ GỘP `main`** — merge `47b8f77` (02/08/2026),
> `--no-ff`. Các dòng VS ở trên trước ghi "nhánh `kh-vertical-slice`" nay đã nằm trên `main`.
> **Còn dở: Task 10–12 — chơi thật 2 client Unity + nghiệm thu.** Runner đã gộp kèm (`tools/e2e`,
> `unity-overlay`) nhưng **gameplay chưa chạy trọn**; làm thành đợt riêng.

**Tình trạng chất lượng hiện tại (kiểm tay 02/08 trên `main` sau gộp):** `go test ./...` —
**toàn bộ 7 gói xanh, 0 đỏ** (gói `game` ~83s), `go build`+`go vet` sạch, `gofmt` sạch (nhiễu CRLF file cũ **đã dọn gốc** 04/08, §9). `pytest assetpack` **124/124 xanh** (máy sinh `part.bin`/catalog/collision).
KH6 còn được **soi toàn nhánh (opus)** + **2 canary** (thả bug thật để kiểm: reviewer bắt
được lỗi byte im lặng, test E2E bắt được lỗi tích hợp).

**Người chơi hiện đã làm được gì:** đăng ký trên web → đăng nhập vào game → tạo nhân
vật, chọn làng → vào bản đồ làng → đi lại, thấy người khác → **chọn chiêu, đánh quái,
bị quái đánh, đánh nhau (PvP), chết rồi hồi sinh**. Còn chờ **nghiệm thu tay bằng client
Unity thật** (Task 15) — đó là việc chỉ chủ dự án làm được.

**KH7 Task 2 — ĐÃ DUYỆT TAY + GỘP `752d6a1`** (điều phối tự đối chiếu `loadItemNew` ↔ `goiCatalogItemMau` 9 trường, xác minh `readShortOptionTemp`=1 byte, tự áp đột biến đảo trường → 2 test đỏ): GATE giải mã hai phía gói
`-28/8` xong ([`docs/superpowers/research/2026-09-03-kh7-giai-ma-28-8.md`](docs/superpowers/research/2026-09-03-kh7-giai-ma-28-8.md))
— client và Java khớp 1-1, mọi "lệch" đều là tràn số im lặng phía Java. Rồi
`server/internal/game/mauvatpham.go`: bảng 8 mẫu option + 5 mẫu vật phẩm Naruto, gửi **ba**
message `-28/8` theo đúng thứ tự 0 → 1 → 2 (subtype 2 phải cuối vì chỉ nó kích `clientOk`).
Dữ liệu để **trong Go** như `noidungboss.go` — **không đụng `catalog.json`/`asset-lock.json`**.

**KH7 Task 3 — ĐÃ DUYỆT TAY + GỘP `b4f83e8`** (điều phối tự đối chiếu khuôn ghi với bảng client, tự áp đột biến tái hiện L1 → 6 test đỏ). **KH7 Task 8 (cửa hàng NPC, gói `-44`) cũng ĐÃ GỘP `d952cfe`** — `cuahang.go`: `MonHang`/`TabCuaHang`/`CuaHang`, tiền chỉ vàng+ngọc (bẫy L10), kiểm hết rồi mới ghi byte, 8 test, 10/10 đột biến agent + đảo vàng↔ngọc điều phối tự áp → 4 test đỏ; giải mã tại [`2026-09-03-kh7-giai-ma-44.md`](docs/superpowers/research/2026-09-03-kh7-giai-ma-44.md). ⚠️ Gộp cả hai → main đứt build vì **trùng hằng `tranParamOption`**, vá `ba3a56d` (xem §8). Pha A task 3: giải mã hai phía
ba gói hành trang ([`docs/superpowers/research/2026-09-03-kh7-giai-ma-36-37-35.md`](docs/superpowers/research/2026-09-03-kh7-giai-ma-36-37-35.md))
— **đủ rõ, không chỗ mập mờ**; cả ba gói dùng chung một khuôn "một món", dấu ô rỗng là
`short -1` ở trường `templateId`. Pha B: `server/internal/game/goihanhtrang.go` —
`GoiTui` (`-36`), `GoiDoDangMac` (`-37`, có `head` chèn giữa), `GoiRuong` (`-35`).
**Sửa lỗi L1 của bản Java: ô rỗng LUÔN ghi `-1` ở cả ba gói.** Dựng xong dữ liệu mới ghi byte;
dữ liệu hỏng → trả lỗi, không gửi gói thiếu. Trần ô = **127** (vì ME_INFO ghi số ô bằng byte
CÓ DẤU). 12 test byte-chính-xác, 3/3 đột biến bị giết. **Chưa nối vào `login.go`** — đó là task 4.

---

## 5. Còn gì chưa làm

Xếp theo thứ tự nên làm:

| Ưu tiên | Việc | Ghi chú |
|---|---|---|
| ⏳ | **Vẽ ảnh theo kế hoạch đợt 4** | Bước 0 + Task 2 đã xong. Task 3 đã lắp thử 15 PNG x2–x4, backup đủ và chụp thật popup/button; **còn đúng bằng chứng HUD sau login** vì máy chủ không hoạt động trong lượt 08/09. Cần bật máy chủ test + dùng tài khoản test được cấp, chụp HUD rồi mới quyết giữ/rollback. Sau đó mới làm tiếp 7 ảnh giao diện P1, rồi **2 dải người chơi vẽ lại** → icon kunai + chân dung + icon app → **6 ảnh map 1 còn thiếu** → boss Lá → hiệu ứng Lá → rừng + quái → NPC Lá. `Small27083/27084` không bị sửa. ⚠️ `docs/my-thuat/anh-nguon/` ngoài git — ảnh vẽ xong phải tự sao lưu. |
| ✅ | ~~KH6 — Chiến đấu~~ | **Xong, gộp `162bbe9`.** Còn **Task 15 — nghiệm thu tay bằng client Unity** (chỉ chủ dự án làm). Hồi máu PvP (`goiHoiMauNguoiChoi`) cố ý hoãn — xem §8 |
| ✅ | ~~KH11 — Lưu trạng thái sống~~ | **Xong, gộp `6ea467f`.** Còn **Task 9 — nghiệm thu tay**: chơi vài phút → Ctrl+C tắt máy chủ → bật lại → đăng nhập xem có đúng chỗ, đúng máu không |
| ✅ | ~~Vertical slice Task 6A — fix round 2~~ | Derived output **20/20 PASS**; raw source vẫn **0/20 PASS** nhưng được ghi rõ và không bị tráo. Gate kiểm semantic loop, raw summary, reviewer/timestamp và nội dung evidence trước khi bind |
| ✅ | ~~Vertical slice Task 6B — bind ảnh đã duyệt~~ | **Xong, review sạch.** 109/109 test; A/B 669 file không khác; exact-20 provenance và evidence HUD/Iruka có trust anchor reviewer |
| ✅ | ~~Vertical slice Task 7 — build artifact + deploy staging theo lock~~ | **Xong trên nhánh `kh-vertical-slice`**, review độc lập sạch. Artifact lock 21, deploy fail-closed (previous lock + protected root + canonical x1–x4), Unity preflight thật, full Python 123/123; chưa gộp `main`. **Tiếp: Task 8** |
| ✅ | ~~Vertical slice Task 8 — bản vá audit collision~~ | Important fail-open ngoài footprint đã có RED/mutation và được sửa fail-closed; **independent audit APPROVE**. Chưa gộp `main`. |
| ✅ | ~~Vertical slice Task 9 — telemetry/control/runner~~ | **Independent audit APPROVE.** Commit `01af875`; chưa gộp `main`. |
| ⏳ | **Vertical slice Task 10 — clean full + warm smoke** | Canonical oracle đã harden **24 case**, 4 mutation; stale lock owner=0 đã đi qua. Full ArtRoot Task 8 v22 đã được rebuild và publish immutable đúng lock `54a638...`/pair `afd03...`; Unity resource preflight có `RESOURCE-PREFLIGHT-OK`. `E2E-ASSET-006` đã được gỡ ở tầng provenance. **Clean full/warm gameplay vẫn NOT RUN và cần phê duyệt bước riêng**; không tự mở Task 11. |
| 1a | **Map 1 Làng Lá — hoàn thiện mỹ thuật** | **6/7 chỉ số đạt** (đo bằng `tools/map-lang-la/do-map.py`, so với concept đã duyệt). Còn: *đất phẳng lì 91% (đích ≤40%)* — **phải vẽ**, đã thử 3 cách làm bằng code, cả 3 hỏng · 2 ảnh nền đúng cỡ **2976** (hiện 1536 nên núi lặp ~1,9 lần/khung). ~~8 ảnh bóng đổ chưa lắp được~~ **xong** (dòng G1 §4). Quy trình + bẫy: [`SO-TAY-LAM-MAP.md`](docs/my-thuat/SO-TAY-LAM-MAP.md) |
| 1a′ | **Cổng rời map — mới làm nửa đường** | Cổng + đổi map đã chạy (dòng G1 §4), nhưng **chỉ map 1 ↔ map 2** có cổng. Map 3/4/5 vẫn là **ngõ cụt**: ai chọn làng Sương Mù/Đá/Mây thì sinh ra ở đó và không đi đâu được. Ba map ấy cũng chưa có mỹ thuật, chưa có `collision` riêng (đang dùng biên dự phòng 4096) và chưa có NPC/quái — làm cổng cho chúng chỉ có nghĩa sau khi có nội dung. Bảng cổng nằm ở `server/internal/game/banDo.go` (`diemDiMap1`/`diemDiMap2`) |
| ✅ | ~~**KH6b — Boss**~~ | **Bản 1 đủ code 02/09 (`ea43237`)**: 15 task, tất cả trên `main`. Boss sinh theo lịch ngoài actor, xuất hiện `-5`, chào/đánh/thoại qua 44, người đánh boss qua `-60` id âm, boss đánh trả bằng chiêu KH6, chia thưởng theo đóng góp vào `KinhNghiem`, chết → xác 3s → biến hình giữ id hoặc `-6` rời hẳn, nổi giận theo %máu, chạy chữ 93 qua loa. **Đã soi** (28 đột biến, 0 lỗi code, 3 lỗ hổng test đã bịt — `1b57a55`). **NGHIỆM THU 03/09** — xem dòng ngay dưới bảng. |
| ✅ | ~~KH6c — Stress test 20 client cùng khu vực~~ | **Xong, gộp `322b140`.** Actor không treo dưới ~5 triệu sự kiện/2s, giữ nhịp, gạch đúng client chậm; thêm đo trễ nhịp (`soNhipTre`/`treNhipLauNhatMs`). Dùng test in-process (`suKienDiChuyen` + `phienStress`), không cần TCP thật |
| 2 | KH7 — Vật phẩm, hành trang, cửa hàng, giao dịch | ✅ **Đặc tả xong + gộp** (`f57cd1a`, [spec](docs/superpowers/specs/2026-09-03-kh7-vat-pham-design.md)): căn cước món = `id BIGSERIAL` bền; 4 bất biến chống nhân bản; 16 đường nhân bản. ✅ **Task 0 (GATE chặn D7) XONG + gộp `4283759`** — gỡ vàng/ngọc khỏi `CapNhatNhanVat`, thêm `store.CongTruTien` (cộng/trừ tại chỗ, CTE+`FOR UPDATE`, chặn `>=0`); đã kiểm tay: đột biến ghi-đè giết cả 2 test chốt. ✅ **Task 1 XONG + gộp** (`5723c32`): migration **`0005_vat_pham.sql`** (KHÔNG phải 0004 như spec — số đó đã là `thuong_cho_nhan`) gồm 3 bảng `vat_pham` / `vat_pham_lich_su` / `vat_pham_giao_dich` + partial UNIQUE `vat_pham_mot_o_mot_mon` (BB1) + 4 CHECK; kiểu `store.VatPham` + hằng `NoiChua*`; 4 test ràng buộc CSDL THẬT (`vatpham_test.go`), đột biến xóa index BB1 → đỏ. ⚠️ Bài học ghi trong test: PostgreSQL coi NULL khác nhau trong UNIQUE nên "bỏ WHERE" khỏi partial index KHÔNG bị test bắt — đừng tin chú thích, phải áp đột biến. ✅ **Task 2 XONG + gộp `752d6a1`** (`mauvatpham.go`: 8 mẫu option + 5 mẫu Naruto, gói `-28/8` subtype 0→1→2 thật, subtype 2 phải cuối vì chỉ nó kích `clientOk`; id mẫu = chỉ số, CHỈ THÊM VÀO CUỐI; không đụng `asset-lock.json`). ✅ **Task 3 XONG + gộp `b4f83e8`** (`goihanhtrang.go`: -36/-37/-35, ô rỗng luôn ghi `-1` — sửa lỗi L1 của Java; trần 127 ô). ✅ **Task 8 XONG + gộp `d952cfe`** (`cuahang.go`: gói `-44` loaiShop=0 + "Tạp hoá Làng Lá" mẫu). Vá trùng hằng `ba3a56d`. ✅ **Task 4 XONG + gộp `18c407a`** (`hanhtrang.go`: `BangHanhTrang` khoá `(nơi chứa, ô)`→id, sống ở goroutine PHIÊN không phải actor khu vực; `DungBangVaGoi` dựng bảng + 3 gói từ MỘT lát cắt — đột biến dựng lệch → đỏ). ✅ **Task 6 XONG + gộp `652315c`** (`store/vatpham_kho.go`: `LayTui`/`CapVatPham`/`ChuyenO`/`HuyVatPham`, mỗi hàm một giao dịch, lịch sử cùng tx, guard chủ sở hữu D4 ở **hai lớp** (SELECT FOR UPDATE + UPDATE), khoá theo id tăng dần D5, hủy = `DA_HUY` không DELETE; bản 1 KHÔNG hoán đổi ô, chuyển ô cùng chủ không ghi lịch sử — có lý do trong file). ✅ **Task 9 XONG + gộp `9320312`** (`muaban.go` + `store/muaban_kho.go`: gói lên `6` = `[typeBuy][mauID short][soLuong short — có thể vắng ⇒1]`, `typeBuy` KHÔNG phải tab, client KHÔNG gửi id cửa hàng → máy chủ nhớ `muaBanCuaHangDangMo` (chưa ai gọi `DatCuaHangDangMo` — task mở NPC cắm); `MuaMotGiaoDich`: idem `ON CONFLICT DO NOTHING` → trừ tiền → cấp món → `ket_qua`, MỘT tx; `CongTruTienTx`/`CapVatPhamTx` dùng chung thân SQL với bản không-tx; phản hồi `-36` → `-25` → `6` (gói `6` gọi `endDlg` nên PHẢI gửi cả khi thất bại). Đột biến tách 2 giao dịch → đỏ — agent tự phát hiện bộ test đầu KHÔNG bắt được, phải thêm `TestMuaCapMonHongKhongMatTien`). ✅ **Task 11 XONG + gộp `0517af6`** (`dovuong.go` + `store/vatpham_dat.go` + `khuvuc.go` +55: gói `68`/`-20`/`-19`/`-21`/`-25`; **L2 xác nhận** type 9/10/34 cần HAI short; nhặt đúng §6.10b: actor gỡ món khỏi `doDuoiDat` NGAY rồi xếp yêu cầu cho goroutine `BoNhatDo` chạm CSDL (`NhatLen` = MỘT UPDATE `WHERE noi_chua='DUOI_DAT'`, D8), hỏng → trả món lại; client tự xoá món sau `-20`/`-19` nên KHÔNG gửi `-21` sau nhặt. ⚠️ `main.go` chưa nối `MoBoNhatDo`/`GanBoNhat` → hiện món nằm nguyên dưới đất, không panic; người vào sau chưa thấy đồ đã rơi (đường vào thế giới = task 5)). ✅ **Task 5 XONG + gộp `7d089b5`** (`thegioi.go` + đoạn vào thế giới `login.go`: ME_INFO `-30/0` ghi **số Ô** (không phải số món) rồi **đủ thân từng ô** — client `Controller.cs:5860-5968` đọc thân ô ngay trong ME_INFO, spec §4.9 nói thiếu; `itemsBody` unsigned ≤255, `itemsBag/Box` **có dấu** ≤127 và ngoại lệ tràn ở đó KHÔNG trong try/catch → kẹt "Xin chờ"; số ô là HẰNG `SoOMacChuan=5/SoOTuiChuan=20/SoORuongChuan=30` (Java lưu sức chứa trong CSDL vì có vật phẩm mở rộng túi — bản 1 chưa có nên giữ hằng, ghi rõ chỗ phải đổi); nạp túi `LayTui`→`DungBangVaGoi` TRƯỚC gói đầu tiên của chuỗi vào thế giới, bảng giữ ở `BoXuLy.hanhTrang`, lỗi CSDL → không gửi byte nào, đóng phiên; `-36/-37/-35` phải SAU ME_INFO. Đột biến số ô túi len+1 → 4 test đỏ). ✅ **Task 12 XONG + gộp `583377e`** (`roido.go` + `khuvuc.go` +33 + `sokhuvuc.go` + `main.go`: quái/boss chết → actor `quyetDinhRoiDo` (thuần, ngẫu nhiên tiêm) rồi `boRoi.Xep`; goroutine `BoRoiDo` gọi `SinhDoDuoiDat` (1 tx, lịch sử SINH_RA) rồi `Gui(suKienDoDaSinh)` → actor `thaDoXuongKhuVuc` phát `68`; `-12` giữ `soDoRoi=0` vì điền id buộc actor chờ CSDL (§6.10b) và client `break` sớm nếu không thấy quái; kho lỗi → không món ma. **Phát hiện: `BoNhatDo` task 11 chưa ai nối** → nối cả hai vào `main.go`, thứ tự tắt có test. Bảng rơi là số khởi đầu (quái 15/10/2%, boss 2 món chắc) — KH8 cân bằng. **Treo nhỏ:** đồ rơi chưa có hạn giờ tự dọn (`-21`/`suKienDonDoHetGio` có sẵn, chưa ai đặt bộ đếm); `chuUuTien=-1` ai cũng nhặt được). ✅ **Task 10 XONG + gộp `6ea49a0`** (gói `7` = `[action][where][oSo short]` — KHÔNG có số lượng, ô 3 là **chỉ số ô** khác gói mua; hai vòng xác nhận→bán; giá bán `max(1, MonHang.Gia/4)` của cửa hàng đang mở (KH8 thêm giá gốc vào mẫu); bản 1 chỉ bán đồ trong TÚI, đồ đang mặc từ chối; `BanMotGiaoDich` idem → `FOR UPDATE` guard D4 trong WHERE → bớt/hủy → cộng tiền → `ket_qua`; gói `6` gửi cả khi thất bại vì client `endDlg` ở đó). ✅ **Task 7 XONG + gộp `c854376`** (`dungvatpham.go`: `-43` = `[type][where][oSo]`(+`short mauId` khi oSo=-1), vòng **2 bước** xác nhận (gói XUỐNG cùng mã `-43` mở hộp Có/Không — bỏ là nút vứt chết im lặng); `-40` = **chỉ 2 byte** `[loai][oSo]`, client KHÔNG gửi ô đích → máy chủ tự tìm ô trống; ô mặc = chính `MauVatPham.Loai` 0..4; kiểm biên `oSo` qua cửa DUY NHẤT `dungDoKiemBien`; `dungDoLamMoiVaGui` dùng chung cho task sau. ⚠️ Bẫy §7.11 làm rõ: test biên đi QUA handler không giết được đột biến bỏ kiểm biên vì bảng phiên là lưới thứ hai → phải test thẳng loại lỗi. **Hoãn có lý do sang bản 2:** "dùng thuốc" (Loai 6 đổi HP → thuộc actor §6.9, cần sự kiện `khuvuc.go`; gói type 0/3 trả `-25` "Chưa dùng được", không im lặng); mặc đè trang bị cũ bị từ chối (`ChuyenO` không hoán đổi); mặc/tháo chưa đổi hình nhân vật (thiếu `Send_Caitrang`). `Loai 6` = tiêu hao vẫn là suy luận (Java `case 6: eatPea`), task 13 xác nhận). ✅ **Bổ sung: hạn giờ tự dọn đồ rơi — gộp** (`b269e9c`): `hetGioDoRoi = 50s` (Java `ItemMap.java:133` `canDoWithTime(createTime, 50000)` — không bịa); `nhipDonDoHetGio` chạy cuối `nhipTim` với đồng hồ tiêm; ba bước **gỡ khỏi map → phát `-21` → xếp huỷ CSDL** cho `BoNhatDo` (cờ `huyHetGio`); `HuyDoDuoiDatHetGio` một bước `WHERE id AND noi_chua='DUOI_DAT'` + lịch sử **`HET_GIO`** mới (tách khỏi `BO` để nhật ký không nói dối); món vừa bị nhặt → 0 dòng → không huỷ nhầm đồ trong túi; huỷ lỗi chỉ log (chấp nhận 1 dòng CSDL mồ côi hơn "nhặt ma"). 6/6 đột biến agent + đột biến bỏ guard của điều phối → đúng test đỏ. Món trả lại sau nhặt lỗi giữ `luc` cũ (cố ý). **KH7 code 0–12 + bổ sung XONG HẾT. Còn 13 = nghiệm thu tay** (chủ dự án). Tiền chỉ đổi qua `CongTruTien`. Giao dịch 2 người hoãn bản 2, ký gửi bản 3 |
| 3 | KH8 — Kho nội dung YAML + nạp nóng | ✅ **Đặc tả xong + gộp** (`05ad2bd`, [spec](docs/superpowers/specs/2026-09-03-kh8-kho-noi-dung-design.md), 627 dòng; điều phối xác minh 4 dẫn chứng `file:dòng`). 5 quyết định: **YAML chứa cái TA quyết định, Go chứa cái CLIENT quyết định** (mã lệnh, trần byte, công thức, tập KIỂU chiêu ở lại Go); id mẫu **tường minh, liên tục từ 0, chỉ thêm cuối** + `khoa-id.yaml` chỉ-thêm; nạp nóng = **ảnh chụp bất biến + `atomic.Pointer`**, actor lấy con trỏ một lần đầu nhịp (loại RWMutex vì actor phải chờ — §6.10b); kích nạp qua **HTTP nội bộ `127.0.0.1`** (VPS Windows không SIGHUP); **`SatThuong` = phần trăm** (giữ `satthuong.go:92`). **14 task, task 0 = GATE tương đương byte**: chuyển dữ liệu Go → YAML → dựng lại `-28/7`, `-28/8`, `-44` phải y hệt byte — vì phát hiện bộ test hiện có **ghim cấu trúc, không ghim giá trị** (`mauvatpham_test.go:47` dựng `want` từ chính `DanhSachMauVatPham()`). **3 chưa xác minh chặn:** số chiêu/số KIỂU chiêu (chủ dự án quyết — LO-TRINH đề 8 chiêu/làng, hiện 4); `go test -race` chưa chạy được — **điều phối kiểm 03/09: máy KHÔNG có gcc/clang/zig ở đâu cả, `CGO_ENABLED=0`** → muốn `-race` thì chủ dự án phải cài một bộ (đề xuất w64devkit, giải nén là chạy, không cần cài đặt hệ thống) — đột biến nạp-nóng-sửa-tại-chỗ chỉ `-race` bắt được; ~~thư viện YAML chưa thử~~ **đã thử 03/09: `go list -m -versions gopkg.in/yaml.v3` → v3.0.1 qua proxy.golang.org, tải được**. **CHỦ DỰ ÁN ĐÃ CHỐT 03/09:** giữ **4 chiêu/làng** (không mở rộng kiểu chiêu), **cài w64devkit** cho `-race` (xong, xem §7), **code KH8 ngay** theo thứ tự lộ trình. ✅ **Task 0 (GATE tương đương byte) XONG + gộp** (`cc7ddb2`): gói `server/internal/noidung/` (`Nap(fs.FS)`/`NapNhung`/`NapTuThuMuc`, `yaml.v3` `KnownFields(true)` — gõ sai tên trường là LỖI, kiểm fail-closed), `server/kho-noi-dung/{chieu,vat-pham,cua-hang}.yaml` (chú thích "vì sao" chép sang; `chieu.yaml` viết tường minh 60 cấp) + gói cầu `khonoidung` để `//go:embed`; **hướng import `game`→`noidung`**; 7 hằng trần client (`tranMauOption`, `tranParamOption`, `tranSoTab`…) **dời về `noidung`**, `game` chỉ bí danh; 3 hàm dựng gói tách `…Tu(dữ liệu)` để GATE chạy qua CHÍNH bộ ghi byte; `kh8gate_test.go` so 4 gói (`-28/7` 3483 B, `-28/8` s0/s1, `-44`) **hai lớp**: Go-vs-YAML và vs **bản vàng hex ghim tay** (`kh8gate_vang_test.go`). 6/6 đột biến agent + điều phối đổi giá 100→101 → GATE `-44` đỏ đúng byte 15. ⚠️ Sau task 0 YAML **chưa chạy trong game** — mới chứng minh tương đương; `ten_lang` trong `chieu.yaml` là bản sao thứ hai của tên làng (task 6 xoá); `go mod tidy` kéo 2 phụ thuộc test gián tiếp của yaml.v3. ✅ **Task 2 XONG + gộp** (`56d9fbf`): `noidung/kho.go` — `Kho{atomic.Pointer[anhChup]}`, `MoKho`/`Thay` (CAS, thế hệ tăng đơn điệu)/`AnhChup`/`TheHe`; **hợp đồng "gọi `AnhChup()` MỘT lần đầu nhịp, giữ suốt nhịp"** (hai lần giữa nhịp có thể rơi vào hai thế hệ — chưa cưỡng chế, task 3 khoá bằng test actor); cửa tra cứu trả **bản sao sâu** (đo: `AnhChup` 0,24ns/0 cấp phát, `TraCapChieu` 2,1ns/0, `CuaHang` 54ns/3, `ChieuCuaLang` 2,95µs/31 — ~0,3% nhịp cho 100 khu vực; chọn sao chép vì "chứng minh không ai ghi" đã phá hai lần: `TimBanDo`, `saoNoiDungBoss`); `anhChup` **không xuất khẩu** có chủ ý (chặn lách bản sao). `-race` 0 cảnh báo; 5/5 đột biến agent (M12 bị bắt bởi CẢ test nhất quán lẫn race detector) + điều phối M14 → 3 test đỏ. ✅ **Task 1 XONG + gộp** (`34a0b14`): `noidung/kiemtra.go` — K1 căn cước (đối chiếu **sổ khoá** `kho-noi-dung/khoa-id.yaml` chỉ-thêm: id biến mất / mã lệch / mục mới thiếu dòng khoá; `skill_id` duy nhất toàn cục), K2 tham chiếu chéo (món có option → mẫu `tran_chong: 1`), K3 miền giá trị (tên option có `#`, tab ≤2 dòng, `hoi_chieu_ms` bội số 100ms, không âm, ước thân `-28/8` < `mark(100000)`); đầu file có bảng "kiểm nào chặn hỏng ngầm nào". Sổ khoá **tự khai** → review bằng `git log -p -- khoa-id.yaml` chỉ được có dòng `+`. Giết M1–M6, M9, M14 (kể cả biến thể độc "đổi chỗ rồi sửa id cho liên tục" chỉ sổ khoá bắt). ⚠️ `ma` sinh từ `ten` (chưa có trường `ma:` riêng — đổi tên hiển thị mẫu đã phát hành làm kiểm đỏ, cố ý); cờ `ngung_dung` (§2.6) chưa làm; M7/M8 (`GiayNghi` boss) chờ boss ra YAML (task 7). ✅ **Task 3 XONG + gộp** (`6988365`): chiêu đọc từ **ảnh chụp** — `Kho` là giá trị **được tiêm** (`main.go` → `NapNhung` → `MoKho` → `soKhuVuc.GanKhoNoiDung` → `KhuVuc.khoNoiDung atomic.Pointer`, đúng khuôn `loa`/`soNo`); nạp thất bại = **dừng máy chủ**; chưa gắn kho (test cũ) → mặc định nạp YAML **nhúng** một lần (`sync.Once`), **bảng Go cứng của chiêu đã XOÁ** (`taoCapChieu`/`chieuTheoLang`…); actor gán `k.acNhip` **đúng một lần** đầu `nhipTim` và đầu `xuLy` (chọn trường thay tham số vì đường gọi 5 tầng); `khochieu_hopdong_test.go` — kho giả đếm số lần lấy ảnh chụp = 1/nhịp (đột biến gọi thêm trong `aiBoss` → đỏ; điều phối gọi thêm trong `suKienChonChieu` → đỏ); bản vàng `-28/7` không đổi một byte; `-race` game+noidung 0 race. ⚠️ **Nợ:** catalog `-28/7` lúc đăng nhập (`goiCatalogSkill`) vẫn đọc kho **nhúng**, không đọc kho tiêm → nạp nóng chưa tới người mới vào (nối ở đường đăng nhập, task 10); `nguoi.chieuDangChon` giữ qua thế hệ (cố ý, spec "không nạp nóng được"). ✅ **Task 10 XONG + gộp** (`10e143d`): `noidung/kich.go` — kênh nạp nóng `net/http` chuẩn **chỉ loopback**, kiểm **hai lớp** (chuỗi trước bind + địa chỉ thật của listener sau bind; `":9090"` cũng bị từ chối); `POST /admin/nap-lai` → `NapTuThuMuc` → lỗi validator **422 + nguyên văn thông điệp, giữ bản cũ**; thành công 200 `{the_he, luc, tom_tat}` + log; `GET /admin/the-he`; sai khoá `X-Nro-Khoa` → 403; chưa có thư mục → 501 kèm cách sửa; mutex là lớp thứ hai (CAS của task 2 đã che — nói thẳng, không giả vờ có test giết). **Biến môi trường mới:** `NRO_NAP_LAI_DIA_CHI` (rỗng = tắt kênh, không loopback = máy chủ dừng), `NRO_NAP_LAI_KHOA` (bắt buộc khi mở), `NRO_KHO_NOI_DUNG` (thư mục YAML — **nếu đặt thì nạp từ đó NGAY LÚC KHỞI ĐỘNG**, không chỉ lúc nạp nóng, để khởi động lại không lặng lẽ quay về bản nhúng — mở rộng của agent, hợp lý). Thứ tự tắt `boLuu → kenhNapLai → lichBoss → loa → TatMayChu → …` (kênh tắt SỚM để không trả "200 thế hệ mới" cho máy chủ đã hết actor). **Nợ task 3 trả xong:** catalog `-28/7` lúc đăng nhập đọc kho tiêm qua `BoXuLy.DatKhoNoiDung` (test đọc ngược byte sau `Thay`). ⚠️ Khoá so `==` không constant-time, chưa phải xác thực thật → **Kế hoạch B §7**; `-28/8` vẫn đọc nhúng chờ task 4. ✅ **Task 8 XONG + gộp** (`5742455` → `672640c`): `kho-noi-dung/quai.yaml` (mẫu quái tra bằng **mã chuỗi**, `anh_id` kiểm với `asset-lock.json` `closure.mob_ids` qua `Catalog.MobIDs` — gói `noidung` không import `game` nên kiểm chia đôi: hàm thuần nhận danh sách id + `game/khoquai.go` gọi với danh sách thật, lệch = panic lúc sinh quái; `sinh:` khai tường minh 6 map), `sat_thuong_co_ban` từng mẫu thay hằng 5 (`quai.go:187`), quái sinh từ ảnh chụp **nhúng** lúc `MoKhuVuc` (chạy TRƯỚC `GanKhoNoiDung`), **quái đang sống giữ số thế hệ lúc sinh** (HpToiDa đã xuống client trong -24 — có test), bảng cứng đã xoá; tầm/nhịp AI + hồi sinh 5s vẫn hằng Go (dùng chung mọi quái, bội số 100ms — agent chủ động thu hẹp, ghi lý do đầu `quai.yaml`). ⚠️ `Quai.MauId` nay mang `anh_id` và cũng là khoá `roiDoCuaQuai` — task 9 phải biết hai vai chung một trường. **Gộp xong main ĐỎ 4 test nạp nóng của task 10** (fixture thư mục tạm chép danh sách file cứng, thiếu `quai.yaml`) → vá `b54a51a`: `noidung.TenCacFile()` là nguồn duy nhất, 3 fixture duyệt nó (bài học §8). ✅ **Task 4 XONG + gộp** (`894c535` → `fb3ce68`): mẫu vật phẩm + option đọc từ ảnh chụp — interface `khoVatPhamAnhChup` **nhúng vào** `khoChieuAnhChup` (một nguồn tiêm, `khuvuc.go` không đổi); cửa không tham số (`DanhSachMauVatPham`…) đọc mặc định nhúng — **mọi chỗ gọi đều trên goroutine phiên/thợ, không trong actor** (bảng chỗ gọi trong báo cáo agent; `TestCanCuocKhuVucDocKhoDaGan` khoá đường đúng); bảng Go cứng xoá; bản vàng `-28/8` không đổi. ⚠️ Gọi cửa không tham số **trong actor** sẽ KHÔNG bị test hợp đồng bắt (đọc nhúng, không qua nguồn đếm) — vô hại hôm nay, cảnh báo ghi tại chỗ. Xung đột gộp ở `khochieu.go` (task 8 + 4 cùng mở rộng interface) → giữ cả hai. ✅ **Task 6 XONG + gộp** (`9ee215a` → `aec6702`): `lang.yaml` + kiểu `Lang` + kiểm + cửa `Lang(id)`; `lang.go` đọc ảnh chụp, bảng cứng xoá; `ten_lang` trùng trong `chieu.yaml` xoá; `khoNhungMotLan()` dùng chung một `sync.Once`; chỗ gọi `TimLang` trong `khuvuc.go` **chưa đổi sang `k.acNhip`** (chưa có nạp nóng cho làng — ghi trong `lang.go`, task sau); kiểm "vị trí trong biên" chỉ là ≤ int16 (gói `noidung` không import `game`). Xung đột gộp 3 file `noidung` (task 8 vs 6 cùng thêm kiểu/field) → giữ cả hai + `TenFileLang` vào `TenCacFile()`. Test chớp tắt `…ThuCungDocDu5Byte…` đỏ 1 lần ở full test, chạy riêng 3/3 xanh (đúng luật §8). ✅ **Task 5 XONG + gộp** (`9b35ce8` → `41a1973`): cửa hàng đọc ảnh chụp, **`gia_ban_lai` tường minh** trong `cua-hang.yaml` (= `gia/4` sàn 1 của công thức cũ, chép nguyên văn: 25/25/250/250/1), kiểm `0 ≤ gia_ban_lai ≤ gia` ở CẢ `noidung` lẫn `game/cuahang.go` (bán đắt hơn mua = máy in tiền); `banDonGia` đọc thẳng, `banChiaGia` xoá; **cửa hàng đang mở giữ bản sao lúc mở** — nạp nóng không đổi giá giữa phiên mua (cố ý: client đã vẽ bảng giá). Điều phối đột biến YAML 25→30 → 2 test đỏ. ✅ **Task 9 XONG + gộp** (`42ee96e` → `282570a`): `bang-roi.yaml` (tỉ lệ **phần trăm** cho người đọc, `roido.go` đổi phần vạn; `chac_chan`/`chi_boss_lon` cờ riêng), **tách hai vai** `Quai.Ma` (khoá bảng rơi) / `Quai.MauId` (= `anh_id`, chỉ ảnh) — hai mẫu chung ảnh không còn rơi giống nhau; `khuvuc.go` 2 chỗ đọc `k.acNhip` (test đếm ảnh chụp bắt gọi thêm); kiểm `mau_id` đã khoá + `ma` quái có thật. ⚠️ Bảng boss chỉ MỘT mục `boss-chung` vì `game.Boss` chưa mang `Ma` (nối ở `lichboss.go` — task 7 nay đã thêm `MucLichBoss.Ma`, việc nối bảng rơi theo mã boss thật để lại); `KiemMaBangRoiBoss` có sẵn nhưng chưa ai gọi trong đường chạy thật. ✅ **Task 7 XONG + gộp** (`d5addc6` → `4e207dc`): `boss.yaml` (`bo_chieu` là struct có tên trường — cặp số trần gõ ngược không ai thấy), kiểm **M7** (nhịp bội số 100ms, cả SAU khi nhân dồn hệ số nổi giận), **M8** (`giay_nghi > 0`), **M10** (chiêu có thật, cấp ≤ catalog, loại ĐÁNH), **M11** (miền `vi_tri`, ↔ `map_vao`), `KiemIdPartBoss` đối chiếu `outfit` với `asset-lock` `part_ids` (`Catalog.PartIDs` mới); lịch boss **đọc lại pha từ ảnh chụp đúng lúc sắp sinh** (sau đủ `GiayNghi`, không phải mỗi nhịp), boss đang sống giữ số thế hệ lúc sinh; mục lịch không `Ma` giữ hành vi cũ (KH6b). Điều phối: đối chiếu tay TỪNG số YAML ↔ Go cũ (cấp/dame/nhịp/hp/nghỉ/5 mốc nổi giận) — đúng; đột biến `nhip_danh_ms` 1500→1550 → M7 đỏ. Xung đột gộp 3 file (`khochieu.go`, `noidung/{nap,noidung}.go`) → giữ cả hai + `TenFileBoss` vào `TenCacFile()`. **⚠️ Điểm mù oracle lặp lại ở task 7 và 9:** đổi số trong YAML (`15→99`, `chac_chan→30%`, `dame 60→61`) **không test nào đỏ** — bộ kiểm chỉ chặn MIỀN, test đọc YAML rồi so với chính nó. Điều phối thêm 2 test ghim số theo bảng Go cũ (`noidung/bangroi_ghimso_test.go` `e1d0cd3`, `boss_ghimso_test.go`), mỗi cái tự kiểm bằng đột biến → đỏ. **Task 13 (vòng chỉnh số) đổi số thì phải đổi cả hai test này — có chủ đích.** **KH8 lớp "chuyển dữ liệu" XONG HẾT: task 0–10 trên `main`.** ✅ **Task 11 XONG + gộp** (`edd4a41`): `game/canbang_test.go` — 4 phép đo Đ1–Đ4 là test có tên, đồng hồ tiêm qua `k.dongHo` sẵn có (KH6 task 8), **không đổi một dòng code sản xuất, không đổi giá trị YAML**; Đ3 PvP là **mô phỏng công thức** `TinhSatThuong` (không qua actor); Đ4 vàng = `gia_ban_lai` món rơi (bảng rơi không rơi vàng); test golden ghim số hiện tại (task 13 được cập nhật). **Chủ dự án bấm đúp `server/DO-CAN-BANG.bat` → đọc `server/bang-can-bang.txt`**, không cần CSDL. Điều phối đột biến sói rừng 10→20% → golden đỏ. **⚠️ PHÉP ĐO LỘ RA (chưa sửa — task 13):** (1) **không làng nào giết nổi một con Sói Rừng** ở kịch bản thật — chết sau 12–32 giây, vì đòn cấp 1 chỉ 2–3 máu (`sat_thuong: 20` là %, `SatThuongKhoiDau` 10..16 × 20% làm tròn hết về 2) còn quái cứng 5 máu và **hai** con cùng trong tầm điểm xuất phát (`quai.yaml` `sinh:` x=160 và x=40, người ở x=100); (2) **Đ2 boss thắng tuyệt đối**: 5 người gục sau 51 giây, boss còn 2708/3000, **chưa từng nổi giận** — hai mốc trong `boss.yaml` chưa bao giờ chạy thật; (3) chênh lệch Đ1 giữa làng **48%** (ngưỡng ≤25%) do làm tròn số nguyên → chỉ số 10..16 của `lang.yaml` gần vô nghĩa; (4) **chakra không bao giờ bị trừ** — `suKienDanhQuai`/PvP chỉ *kiểm* qua `KiemDanhHopLe`, `chakra_dung` là số trang trí (lỗi code thật, không phải số — cần task riêng); (5) PvP: Làng Đá 4/0, Làng Mây 0/4. Chớp tắt mới: `TestTatMayChuVanDungKhuVucKhiHetGioLuu` (ctx 1ms + Sleep thật). ✅ **Task 12 XONG + gộp** (`9ebbdec` → `f1227fd`): đơn vị sát thương chiêu **chốt = phần trăm** — YAML `sat_thuong` → **`he_so_phan_tram`** (60 cấp, `KnownFields` từ chối tên cũ, có test), `noidung.CapChieu.HeSoPhanTram` (còn `game.CapChieu.SatThuong` giữ vì là khuôn gói `-28/7` — spec §5.1), K3 `0 < he_so_phan_tram ≤ 1000`, TODO `satthuong.go` xoá, **công thức `/100` không đổi**, bản vàng `-28/7` không đổi byte. Điều phối chứng minh không đổi số: chuẩn hoá tên trên bản cũ rồi so 831 dòng dữ liệu = giống hệt; 6 YAML khác 0 dòng dữ liệu đổi; đột biến `2000` → K3 chặn. ⚠️ Agent lại để **bản đổi tên dở (60 dòng) trong cây chính** — lần 3 (§8), bỏ bằng `git checkout --`. `sat_thuong_co_ban` của quái giữ tên (đã là số tuyệt đối). ✅ **Task 13 XONG + gộp** (`3b448d3`…`3221026` → `f3235f8`): 3 vòng đo→chỉnh→đo, **cả 4 ngưỡng đề xuất §5.3 đạt** — Đ1 8,1–9,1 giây (cả 5 làng nay GIẾT ĐƯỢC hai con sói, không ai chết), Đ2 hạ boss 2 phút 02, nổi giận 2 lần, 2/5 sống, Đ3 chênh 12,3% + không làng nào toàn thắng/thua, Đ4 27–30 phút. **Số đổi (chỉ 3 file YAML):** `lang.yaml` sát thương gốc 10..16 → **50..60** (vì `× he_so_phan_tram 20 / 100` làm tròn xuống nuốt sạch chênh lệch — 10..16 đều ra 2), máu 100..250 → 160..220 (Đá 250 thắng cả 4, Mây 120 thua cả 4); `boss.yaml` pha 1 `dame` 60 → **15**, `hp_chon` [2400,3000,3600] → [3600,4500,5400]; `bang-roi.yaml` sói rừng 10/5% → 12/8%. Điểm sinh hai sói, giá áo 1000, công thức, `khoa-id`/`anh_id`/`outfit`: không đổi. 0 dòng code sản xuất; 10 file test đổi con số dẫn xuất (điều phối đếm: không bớt một `t.Fatal`/`t.Error` hay test nào), 2 test ghim số cập nhật có chú thích lý do. Điều phối tự chạy bộ đo → cùng bảng. **⚠️ CHỦ DỰ ÁN QUYẾT:** (1) **boss pha 1 giờ nhẹ đòn phi lý** (5–6 máu/đòn, người 200 máu chịu ~36 đòn) — là số học ép vì bộ đo **không có hồi máu trong trận**; đường ra thật là cho uống thuốc trong trận (bảng rơi đã có thuốc) rồi nâng dame lại — chấp nhận boss nhẹ, hay làm "dùng thuốc" (KH7 hoãn bản 2) rồi cân lại?; (2) ngưỡng §5.3 vẫn là đề xuất — anh bấm `DO-CAN-BANG.bat` đọc bảng rồi chốt; (3) PvP khít hơn 2T/1B/1H thì phải đổi công thức khắc chế 125/100/80 (`satthuong.go`), không đổi số được; (4) `chakra_dung` chưa có tác dụng (code không trừ). **KH8 XONG HẾT 14 task (0–13) trên `main`** |
| ⏳ | **Chơi thử thật bằng client Unity — 05/09/2026, điều phối tự lái qua MCP** | Chủ dự án vào game được, báo "đầu nhân vật tách khỏi thân, map lỗi vài điểm". Điều phối nối **Unity MCP qua HTTP `127.0.0.1:8080/mcp`** (script `unity_mcp.py` trong scratchpad phiên: gọi JSON-RPC thẳng vì MCP `UnityMCP` của phiên Claude Code mất kết nối lúc mở) → đọc trạng thái nhân vật, tự đăng nhập (`Rms.saveRMSString` + `LoginScr.doLogin`), chụp màn hình (`ScreenCapture`), đọc Console. **Phát hiện + đã vá trên `main`:** (1) **nhân vật chết trong CSDL đăng nhập vẫn chết** (`narruto` hp 0/150; `-42` kẹp máu về đầy nên HUD nói sống, actor nói chết) → `hoiSinhLucVao` trước ME_INFO (`b2ef289`, test tích hợp CSDL thật + đột biến); (2) gói `-5` ghi **thân/chân 0/0** cho người chơi → client khác vẽ part 0 (part ĐẦU 3 khung) làm thân → `IndexOutOfRange`; nay 57/58 (`hinhdang.go`); (3) client Unity chỉ biết **3 hành tinh**: `cgender` = làng 1 → đồ Namek 59/60, làng 3/4 → không đặt thân/chân → nay byte cgender luôn 0, byte lớp/làng vẫn = làng; (4) màn tạo nhân vật gửi **id tóc Dragon Ball** (`hairIDOld`) → lưu `dauHopLe` = 30; (5) **tắt máy chủ panic** (`Shutdown` trên `*http.Server` nil bọc interface khi E2E control tắt = mặc định) → **mọi Ctrl+C đều bỏ qua bước lưu** → `coControlThat` (`6dc72a7`, đột biến → panic). Kiểm thật: sau vá, đăng nhập qua MCP `hp=150/150 đầu 30 thân 57 chân 58`, CSDL lưu 150 lúc 16:23. **Nguyên nhân "đầu tách thân" tìm ra bằng công cụ mới `tools/nhan-vat/ghep-khung.py`** (ghép 33 khung `cf` đúng phép toán `Char.paintCharBody` + `CharInfo`, không cần Unity): pipeline `tools/artqa/task6b_closure.py` **co MỖI KHUNG riêng cho vừa hộp 128×160** (`fit_to_canvas`) → nhân vật to nhỏ khác nhau tuỳ khung (thân rộng 32 nhưng cao 6..32 px), chân lơ lửng 10–26 px, khung 20/31 (nằm) thân còn 6 px nên đầu hở +7 px; đầu chỉ có 3 khung idle nên luôn cỡ idle. **✅ Đã sửa + gộp (`a28750f`, agent Opus + điều phối duyệt tay):** một tỉ lệ chung `128/340`, neo bàn chân `cy−11`, `dx/dy` từng khung bù `CharInfo` (cf neo nhỏ nhất), canvas trim theo pose; 34 sprite ×4 zoom sinh lại bằng `tools/nhan-vat/sinh-lai-sprite.py` (dùng hàm toán của `task6b_closure.py`), `part.bin` + `asset-lock` + `vsData` 27; **đã chép vào client** (`chep-sprite-client.py --that`, bản cũ lưu `docs/my-thuat/lang-la-client-assets/smallimage-truoc-05-09/`). Đo bằng `ghep-khung.py` (điều phối chạy lại, khớp từng số): bàn chân 33 khung **cùng một cao độ** (trước lệch 16 px), cao tổng ±3 (trước −14..+1), hết hở thân–chân; **khe đầu–thân còn +1..+5 px ở 8 khung và +19 ở 2 khung nằm** vì 5 dải ảnh nguồn vẽ nhân vật cao thấp khác nhau (idle 653 px, run 459 px) mà đầu chỉ có 3 khung idle → **anh quyết**: (a) dịch đầu xuống 3–5 px (hết hở, idle chồng sâu hơn), hay (b) vẽ lại 5 strip cùng chiều cao (triệt để). Điều phối áp đột biến neo −11→−5: **25/25 test vẫn xanh** → thêm test ghim hằng vào mapping + số đo (`12f3c0d`). ⚠️ Closure `production-closure-v11` **không sinh lại được** (reviewer độc lập phải ký lại `ASSET-ARTIFACT-REVIEW-016`) → `sprites[].zooms` trong mapping vẫn trỏ ảnh cũ, sự thật mới ở `build_inputs.player_component_geometry`; **`assetpack build` lúc này sinh `client/smallimage` từ ảnh CŨ** — đừng deploy đầu ra đó cho tới khi closure được ký lại. `REQUIRED_CLIENT_CF` (`assetpack/animation.py`) **lệch một bậc** so với `CharInfo` ở thân/chân — việc riêng. **Y sinh 384→432 đã làm** (`3abd395`, script `sua-y-sinh-432.py`): sói `y=432`, NPC `cy=432` kiểm qua MCP. `pytest assetpack` 10 test `test_manifest` đỏ sẵn trên máy này = `subprocess` giải mã GBK → vá `encoding="utf-8"` → 155/155. Công cụ vào repo: `tools/unity-mcp/` (cầu HTTP + kịch bản C#), `tools/goi-tin/giai-ma-24.py`. **Map:** nền núi lặp 2 lần = ảnh 1536 thay 2976 (hàng 1a, anh vẽ); vết nối đất giữa màn = bậc gò 1 ô có chủ đích (SOI-LAN-3 #17/#18/#22 là tile art). **Thêm 2 lỗi code tìm ra tiếp trong ngày:** (6) **NPC Iruka chưa bao giờ hiện, 11 bóng đổ (G1) chưa bao giờ vẽ** — hai nguyên nhân chồng nhau: client deref `taskMaint` (null vì máy chủ không gửi gói 40) khi đọc NPC trong `-24`, VÀ ô 6 của mỗi con quái ghi `Quai.Id` trong khi client đọc là **`templateId`** → sói thứ hai `tpl=1` vượt mảng 1 mẫu → `loadInfoMap` ném → bỏ nốt NPC/đồ rơi/bóng đổ. Tìm ra bằng **log Debug hex `-24` một dòng/lần vào map** + `scratchpad/giai-ma-24.py` giải mã đúng theo code client (không đoán). Vá: `nhiemvu.go` gói 40 giữ chỗ (taskId 100) trước `-24`; `quai.go` ô 6 = `MauId`; spec §3.6 đính chính. (7) **Mở được cửa hàng từ client thật**: menu Iruka `["Cửa hàng","Đóng"]` (Go, thay catalog bị ghim SHA — cầu tới KH9 `npc.yaml`) + handler cmd 32 → `DatCuaHangDangMo` + `-44`. Kiểm qua MCP sau vá: `vNpc=1`, 2 sói, `vCurrItem=11`, hộp thoại Iruka hiện 2 nút. (8) **Đứng yên 2 phút là bị đá**: `internal/net` hạn đọc mỗi gói 120s, client Unity không có gói ping định kỳ → log Debug `i/o timeout` đúng 2 phút sau đăng nhập, 3 lần liên tiếp; nâng lên 10 phút (`aa17ddd`, kết nối chết thật do TCP keep-alive lo) — đo lại: 3 phút 22 giây không thao tác vẫn trong game. **Dữ liệu còn treo (làm sau khi gộp nhánh sprite vì đụng `asset-lock.json`; script sẵn `tools/map-lang-la/sua-y-sinh-432.py`):** sói/người/NPC sinh ở `y=384` (hàng 16, giữa trời) trong khi mặt cỏ map 1 là `y=432` — quái không rơi (`Mob.y` kéo về `yFirst`) nên **hai con sói lơ lửng 2 ô**; cần đổi `quai.yaml`/`lang.yaml`/`catalog.json` maps.spawn + npcs.y (429) + lock + ~8 test ghim 384. |
| ⏳ | **Sau khi sửa khe đầu–thân (06/09): việc còn lại quanh nhân vật/NPC** | (1) **Vẽ lại 5 khung chạy + 3 khung trúng đòn đứng thẳng theo thước** `docs/my-thuat/dot-4-ke-hoach-tong-the/anh-tham-chieu/THUOC-NHAN-VAT.png` để bỏ giá "thấp đi ~10 % khi chạy" (brief `01-nhan-vat` đợt 4). (2) Closure v11 vẫn là hình học CŨ trong `sprites[].zooms`/asset-lock outputs — cần reviewer độc lập ký closure mới (ASSET-ARTIFACT-REVIEW-016) rồi đưa 37 sprite thật vào lock. (3) ~~Vết hồng trên nhân vật~~ **XONG 06/09** (293 px → 0; xem §4). Còn **31 px chroma ở ảnh MAP TILE** `x*/t/1$70..1$80.png` (nặng nhất x4/1$80 18 px) — nếu thấy chấm hồng ở MẶT ĐẤT thì là chỗ này, chưa sửa vì ngoài phạm vi đợt vá nhân vật. (4) Sprite sói không có trong client (`res/Mob/` trống) — sói tàng hình trong game; đường ảnh quái → client chưa nối (kiểm kê đợt 4 §4). (5) `ch` = 34 + dy đầu + 10 = 6 px (`Char.cs:6087`) — chữ sát thương/chat có thể neo thấp; chưa nhìn thấy trong game vì chưa bị đánh lúc chụp. (6) Tài khoản thử `naruto` sinh ra ngay điểm boss: boss xuất hiện 10 phút sau khi máy chủ khởi động và giết trong ~9 s; sau chết client hiện "Xin chờ" rồi rớt về màn chọn máy chủ sau ~1–2 phút (luồng chết chưa có ở máy chủ Go: `statusMe 14` không có gói hồi sinh/về làng). |
| 4 | KH9 — Nhiệm vụ, cốt truyện, NPC | |
| 5 | KH10 — Cảm giác chơi (hoạt ảnh, chuyển màn, hướng dẫn người mới) | Ba con số ràng buộc đã chốt ở KH5 |
| — | **W3 — Web:** cấu hình / nhật ký / nạp thẻ, phân trang, thu hồi phiên JWT | Một phần phụ thuộc hệ thống game chưa có |
| — | **B — Bảo mật** (xem §7) | **Phải xong trước khi mở công khai** |
| — | **V — Triển khai VPS**, giám sát, phát hành | |

**Nghiệm thu KH6b Boss (03/09, qua Unity MCP + máy chủ tạm):**
- **Chứng minh SỐNG trong client Unity thật:** máy chủ sinh boss đúng lịch (log `da gui lenh sinh boss`),
  **chạy chữ toàn máy chủ** (log `phat thanh toan may chu noiDung="Boss Kẻ Lạ Mặt vừa xuất hiện!"`),
  client **nhận và hiển thị** hai boss trong danh sách boss trên màn với ĐÚNG tên Naruto + máu khớp dữ
  liệu (`Kẻ Lạ Mặt [2.400]`, `Học Viên Đào Tẩu [600]`), map Làng Lá + nhân vật Naruto vẽ đúng (ảnh chụp).
- **Combat (đánh trúng / boss đánh trả / lột xác giữ id / nổi giận):** chứng minh **tất định bằng bộ test
  tích hợp máy chủ** — `go test ./internal/game/` xanh, gồm `TestBossHoatDongDanhNguoiTrongTam` (boss đánh
  người trong tầm), `TestNguoiMoiVaoThayBoss` (người vào SAU khi boss sinh vẫn nhận boss),
  `TestBienHinhGiuNguyenIdVaGuiSauTruocNam` (lột xác giữ id), `TestMocNoiGian*` (nổi giận theo %máu),
  `TestChiaThuongTheoDongGop*` (chia thưởng).
- **Giới hạn công cụ (KHÔNG phải lỗi code):** chọc combat tương tác qua Unity MCP reflection không tái hiện
  ổn định được vì (a) vị trí client em set qua reflection **desync** với vị trí authoritative của máy chủ,
  (b) nước đi vượt `NguongBuoc=500` hoặc vướng va chạm bị máy chủ từ chối im lặng → đòn thành "vượt tầm"/
  "miss", (c) gói gửi bất đồng bộ. Máy chủ chỉ tính sát thương theo vị trí THẬT của nó, nên poke tay không
  tới tầm boss. Đây là ma sát của bộ công cụ; đường code combat đã được bộ test trên phủ.

**Hai việc chỉ chủ dự án làm được (nghiệm thu tay):**
1. Mở client Unity thật, đăng nhập → xem nhân vật có đứng trong làng không *(nghiệm thu KH4)*.
2. Mở **2 client cùng lúc, cùng làng** → xem hai nhân vật có thấy nhau và đi mượt không *(nghiệm thu KH5)*.

---

## 6. Lỗi sai đã gặp — đọc để đừng lặp lại

Đây là phần đắt giá nhất của file này. Mỗi dòng là một lần đã mất thời gian.

### Về giao thức nhị phân

1. **UTF-8 chuẩn ≠ modified UTF-8 của Java.** Dùng UTF-8 thường thì emoji làm client
   Java ném `UTFDataFormatException`. Phải dùng đúng biến thể của Java.
2. **Chiều dài gói: 2 byte hay 3 byte tuỳ lệnh.** Nguyên tắc đúng:
   *mã hoá* quyết định bởi `coKhoa()`, *số byte chiều dài* quyết định bởi mã lệnh,
   **độ dài thân gói không quyết định gì cả**. Làm sai → lệch dòng byte, client hỏng vĩnh viễn.
3. **Framer bất đối xứng byte-order là CỐ Ý** (đọc big-endian, ghi little-endian) —
   vì client thật làm vậy. Ai thấy "lỗi" này đừng sửa; đã có chú thích cảnh báo trong code.
4. **Chống dịch chuyển tức thời của bản gốc gần như không có** (chỉ map 85–91, khoảng
   cách > 500). Đây là lỗ hổng thật, bản Go đã vá chặt.
5. **Máy chủ Java cũ có lỗi treo "Xin chờ"**: tài khoản bị khoá chỉ gửi `-26`, không gửi
   `-102` và không `return` (`NDVSqlFetcher.java:127-129`). Bản Go **không** lặp lại lỗi này.
6. **Một mã lệnh có thể là HAI gói khác độ dài** (phát hiện 26/07 khi giải mã chiến đấu).
   `-60` có byte thứ **11** (`docTiep`): bằng `1` thì còn 7 byte nữa (đòn đánh, thân 18 byte),
   bằng `0` thì **hết gói ngay** (chiêu hồi máu, thân 11 byte). Tên biến gốc là
   "read continue" nên rất dễ tưởng là hằng số. Tương tự: gói **lên** `54` dài 1 byte
   *hoặc* 5 byte tuỳ `mobId == -1`; `-45/0` và `-124` cũng biến thiên.
   **Cách tránh: dựng mỗi dạng bằng một hàm riêng, không dùng một hàm có cờ** — hàm có cờ
   thì người sửa sau ghi thừa trường lúc nào không hay, mà ghi thừa là lỗi **im lặng
   tuyệt đối** (client vứt byte dư, không lỗi, không triệu chứng — xem §6.18).

### Về cơ sở dữ liệu

7. **`lower()` của PostgreSQL dưới locale C không hạ được chữ có dấu** → `Anguoi` và
   `anguoi` (có dấu) thành hai tài khoản khác nhau. Đã vá bằng ràng buộc
   `CHECK '^[A-Za-z0-9_]{4,20}$'` — chặn được cả trò giả mạo bằng chữ Kirin/Hy Lạp
   trông giống Latin, thứ mà đổi collation **không** chặn nổi.
8. **Vòng lặp băm mật khẩu phải có trần.** Chuỗi băm hỏng làm `Verify` treo hàng phút
   (Java tình cờ an toàn nhờ số 32-bit tràn; Go 64-bit thì chạy thật). Đã chặn ở 1.000.000 vòng.

### Về luồng / đồng thời (quan trọng nhất từ KH5)

9. **Kiểm tra hợp lệ phải nằm TRONG actor.** Kế hoạch ban đầu định kiểm nước đi ở luồng
   phiên (đọc `nv.ViTriX/Y`) — nhưng actor **ghi** đúng hai trường đó → tranh chấp dữ liệu.
   Luật: **chỉ luồng `vongLap` của khu vực được chạm vào trạng thái khu vực.**
10. **Gói `-28/13` (clientOk) bị gửi lại lần hai** cũng gây tranh chấp tương tự. Đã chốt
    guard chống chạy lại; tương tự cho `-39`.
10b. **[✅ ĐÃ VÁ — commit `a135979`, xác minh chạy test lại 03/08/2026]** Giữ nguyên mô tả
    vì là bài học kiến trúc quan trọng nhất. **Actor tự đi phát gói mạng = một kết nối treo
    làm đứng cả khu vực** (phát hiện 27/07 khi AI quái của Task 9 thành nguồn lưu lượng tự
    phát). *Trước khi vá:* `phatChoNguoiKhac` / `phatChoCaKhuVuc` gọi `phien.Gui()` **đồng bộ
    ngay trên goroutine `vongLap`**; `Gui()` có hạn giờ ghi **30 giây** (`net/session.go`).
    Một client ngừng đọc (mạng chập chờn, app xuống nền, hoặc cố ý) làm `vongLap` đứng tối đa
    30s → cả bản đồ đứng hình, quái ngừng hồi sinh. Tệ hơn: `Gui()` lỗi chỉ ghi log, **không
    gạch tên** người đó khỏi `nguoiTrong` → lần phát sau lại treo tiếp, **lặp vô hạn**. Đây
    từng là **lỗ hổng Nghiêm trọng #1** của review toàn dự án 27/07.
    **Đã vá (`a135979` "tach writer socket khoi actor gameplay"):** mỗi phiên một hàng đợi
    gửi riêng `hangGui` (sức chứa 256) + goroutine `vongWriter` riêng; actor gọi `XepGui`
    **chỉ ném gói vào queue rồi đi tiếp, không bao giờ chạm socket**; client đầy hàng bị
    `loaiNhieuNguoiLoiGui` **gạch tên khỏi `nguoiTrong` + `DongCham`** (đóng socket không
    chờ) rồi báo `-6` cho người còn lại. Chứng minh 2 tầng: `TestXepGuiKhongChanFIFO256...`
    (net — `XepGui` không chặn, chặn ở 256) và
    `TestActorLoaiClientChamKhongChanClientKhoe_SERVER_ACTOR_005` (game — probe actor
    **<100ms** khi 1 client đầy hàng → actor KHÔNG treo, gạch đúng client chậm). Đã chạy lại
    03/08, cả hai xanh (`-count=5`). **KH6c (stress 20 client) sẽ chứng minh thêm dưới tải thật.**

### Về kiểm thử (bài học lặp lại nhiều lần nhất)

11. **Bộ test có thể xanh trong khi code sai.** Đã xảy ra ít nhất 3 lần:
    - Xoá hẳn gói `-102` → 6/6 test vẫn đạt.
    - Xoá cả đoạn điều phối `-28/2` (tạo nhân vật) → 24/24 test vẫn xanh.
    - Test đếm số làng bằng chính nguồn dữ liệu nó phải kiểm → điểm mù oracle.
    **Cách chữa duy nhất hiệu quả: kiểm thử đột biến** — cố tình phá code, test phải đỏ.
    Không đỏ nghĩa là test vô dụng, phải viết lại test chứ không phải mừng.
12. **Test dùng giá trị ghim (hard-code), đừng lấy từ chính nguồn đang kiểm.**
12b. **Hai tập tự khai có thể cùng sai mà vẫn “khớp”.** Task 6B từng cho phép xoá
    cùng lúc 1 record `source_provenance` và SHA closure tương ứng; build vẫn xanh với
    19/20 nguồn. Validator phải so với **tập độc lập tuyệt đối** trong review, không so
    manifest với chính manifest.
12c. **Tên reviewer nằm trong chính JSON review không phải trust anchor.** Chỉ kiểm
    “khác reviewer nguồn” vẫn cho alias tự ký lọt. Artifact production phải pin identity
    reviewer độc lập ở ngoài payload đang được kiểm, và có mutation alias tự ký.

### Về quy trình

13. **Đừng commit việc ngoài phạm vi task xen vào giữa một task** — gói review bị lẫn,
    người soi đọc nhầm. Commit việc lẻ **trước khi** bắt đầu task hoặc **sau khi** review xong.
14. Khi người thi công thấy thứ trông như lỗi nhưng không chắc → **dừng lại hỏi**, đừng tự
    sửa. Đã cứu được một lần (vụ byte-order §6.3).

### Về công cụ

15. **KHÔNG sửa file tiếng Việt bằng PowerShell `Get-Content -Raw` + `Set-Content`.**
    PowerShell 5.1 đọc file UTF-8 bằng codepage ANSI rồi ghi lại thành UTF-8 → **mojibake
    sạch cả file**, và `git diff` phình ra như vừa viết lại toàn bộ. Đã vấp khi sửa
    `server/README.md`, phải `git checkout --` rồi làm lại. **File có dấu thì luôn dùng
    công cụ Edit/Write**, không dùng shell.
15b. **Đừng import bytecode duy nhất trước khi phục hồi source cùng module.** Python sẽ
    tự ghi đè `.pyc` khi thấy file `.py` mới. Task 6B từng mất bản bytecode generator đầy
    đủ vì vấp đúng việc này; cuối cùng phải phục hồi nguyên source từ JSONL của phiên Codex.
    Cách tránh: sao lưu bytecode ra path tạm trước, hoặc trích source lịch sử rồi mới import.

### Về giao thức nhị phân — bổ sung sau khi đọc CLIENT (26/07)

16. **Client là trọng tài duy nhất — đọc code Java thôi là CHƯA ĐỦ.** Đọc server Java rồi
    đối chiếu client Unity (`D:\NRO_Project\Assets\Scripts\Game1\Controller.cs`) đã lộ ra
    **ba lỗi thật** mà đọc một phía không thể thấy:
    - `-8`: `Service.java:1157` ghi `short` trong khi client `readInt` → **người chơi khác
      không thấy ai vừa chết** (client tra id hụt, bỏ cả gói).
    - `-17`: `Service.java:1147` nhét **id nhân vật vào ô `cPk`** — gói này không có
      trường id, client dùng thẳng `Char.myCharz()`.
    - `-9`: `Mob.java:1369` ghi `writeInt(-1)` (4 byte) trong khi client đọc **1 byte**
      `hieuUngId`; tình cờ vẫn ra `-1` nên chưa ai phát hiện.
    **Luật rút ra: gói nào quan trọng thì phải đọc CẢ HAI PHÍA trước khi chốt bảng byte.**
17. **Một mã lệnh có thể bị bản gốc dùng cho HAI việc khác nhau.** `Message(87)` vừa là
    "quái bị gió thổi" (`byte + bool`, khớp client) vừa là "gửi maxHp" (`byte + int`,
    **không** khớp client). Chép mù theo bản Java sẽ khiến quái có máu ≥ 16.777.216
    **đứng im không di chuyển**, vì byte cao của maxHp rơi đúng vào ô `isWind`.
    Bản Go: **một mã lệnh, một ý nghĩa.**
18. **Sai trường bên trong một gói KHÔNG làm hỏng dòng byte** — client đọc trọn gói theo
    độ dài rồi mới phân tích (`Session_ME.cs:181-219`), và nuốt lỗi bằng `catch` rỗng.
    Triệu chứng vì thế **im lặng**: gói không làm gì, hoặc hiện số vô lý — không treo,
    không log. **Chạy thử bằng mắt sẽ không phát hiện được; phải test byte-chính-xác.**
    (Thứ làm hỏng cả dòng byte là sai **số byte của trường độ dài** — đó là §6.2, khác hẳn.)
19. **Tên file `chroma` không chứng minh được alpha an toàn.** Nền magenta chuyển sắc có thể
    nhìn giống chroma nhưng vượt sai số soft-matte tĩnh; HUD gần trắng lại không hề là chroma.
    Vì vậy review phải FAIL theo evidence (không được tự suy PASS), và Task bind chỉ nhận source
    có PASS + SHA khớp.
20. **Hash/kích thước evidence chưa đủ để chứng minh nội dung.** GIF `loop=0` có thể làm
    one-shot chạy vô hạn; alpha/contact sheet/crop report có thể bị thay bằng file màu khác
    rồi cập nhật SHA. Gate phải tái tạo pipeline, đối chiếu nội dung canonical và fail closed
    với JSON null/malformed; raw source FAIL và derived output PASS là hai trạng thái khác nhau.
21. **Deploy có thể xanh giả nếu tham số bảo vệ là tuỳ chọn hoặc predicate quá rộng.** Vòng
    đầu Task 7 cho phép caller bỏ `previous_lock`/`protected_unity_root`, và preflight chỉ
    tìm một file bất kỳ trong nhóm `Small*`; vì vậy runtime đổi cùng version hoặc thiếu
    `Small27041.png` nhưng còn `Small27042.png` vẫn có đường lọt. Cách tránh: API/CLI bắt
    buộc hai tham số, lock phải có closure ID độc lập, và mỗi resource canonical phải tồn
    tại đúng x1–x4 trước khi publish; mutation decoy phải đỏ `ASSET-PREFLIGHT-022`.
22. **Collision phải tách semantic, không dùng thẳng tile index.** Artifact phải ghim header
    `width,height,tileWidth,tileHeight` và cell `0/1`; server phải dùng cùng tile 24x24 với assetpack,
    có canary solid trên artifact thật, và test foot-point/segment trên chính artifact.
23. **Gameplay sau khi vào actor không ghi socket trực tiếp.** Command NPC phải `XepGui` qua outbound
    queue; test net.Pipe phải chứng minh handler không bị chặn khi client không đọc, và giữ FIFO.
24. **Collision artifact và biên movement không được là hai nguồn sự thật.** Grid 5×17,
    tile 24×24 chỉ phủ 120×408; nếu ngoài grid trả walkable trong khi movement còn biên
    4096 thì người chơi thoát artifact ngay từ spawn. Map có collision phải dùng chính
    footprint đó làm biên authoritative và fail-closed ở mọi tọa độ ngoài grid; biên
    fallback chỉ dành cho map chưa gắn artifact. Test phải có mép trong, mép ngoài,
    segment crossing, tọa độ âm và đường actor thật.
25. **Reconnect telemetry phải cô lập theo thế hệ phiên, không chỉ theo socket hiện tại.**
    Sau khi client A rớt, action/capture đã xếp hàng của A vẫn có thể chạy trên main thread
    rồi gửi kết quả nhầm sang client B mới kết nối. Mỗi phiên phải có generation riêng và
    kiểm lại generation ở cả lúc dequeue lẫn lúc gửi response; test phải dựng đúng kịch bản
    A enqueue → A rớt → B nối → main thread drain.
26. **Evidence “đúng schema” vẫn có thể xanh rỗng.** Validator chỉ kiểm từng phần tử sẽ
    vô tình chấp nhận `events=[]` hoặc `artifacts=[]`. Evidence E2E phải bắt buộc hai mảng
    không rỗng, kiểm sequence/UTC/client/type/payload của mọi event và tự băm lại SHA của
    mọi artifact; screenshot thiếu hoặc sai hash phải fail-closed.
27. **Fixture E2E không được giả lập mật khẩu production bằng plaintext.** Chỉ nhận DSN
    trỏ chính xác `nro_test`, sinh đúng hai account có prefix từ run ID, dùng cùng
    PBKDF2-SHA256 với server và chỉ nhận mật khẩu thô qua biến môi trường. Cleanup phải
    xóa đúng con trước cha trong `finally`, rồi truy vấn lại chứng minh còn 0 account.
28. **Clean và warm không được tự chấm điểm ở hai nhánh runner khác nhau.** Nếu mỗi phase
    tự viết assertion, một nhánh có thể bỏ kiểm damage thừa/HP 137 còn nhánh kia dùng
    executable hoặc asset lock khác mà aggregate vẫn xanh. Cả runner thật và fixture test
    phải gọi cùng một oracle fail-closed: parser tự đọc NDJSON, băm lại DB/screenshot,
    kiểm scope `full|warm_smoke`, rồi chỉ ghi aggregate PASS khi executable và server lock
    bằng tuyệt đối. Unity asset lock là blocker thật; không xóa lock để “làm test chạy”.
29. **File lock tồn tại không đồng nghĩa Unity đang giữ lock.** Chẩn đoán Task 10 thấy cả
    `ArtifactDB-lock` và `SourceAssetDB-lock` còn trên đĩa nhưng không có Unity/Hub/worker,
    Windows Restart Manager báo 0 owner và mở read-only độc quyền đều thành công. Preflight
    chỉ dùng `Test-Path` vì vậy chặn nhầm file stale. Khi sửa ở một task được phép, phải
    phân biệt file hiện diện với handle đang hoạt động; tuyệt đối không chữa bằng xóa lock
    hoặc kill Unity.
30. **Server lock, asset lock và file evidence phải là ba cổng độc lập.** Một SHA 64 ký tự
    vẫn có thể là SHA đúng của file rỗng, và dùng server lock thay cho asset lock làm warm
    resource “khớp” sai nguồn. Oracle phải kiểm length dương trước khi băm, tách
    `server_lock_sha256`/`asset_lock_sha256`, rồi so `applied_lock_sha256` với asset lock.
    Fixture phải đi xuyên runner boundary; parser đứng riêng không chứng minh runner thật.
31. **Build deterministic chưa đủ nếu artifact canonical bị dọn sau test.** Task 8 đã sinh
    lock/server runtime version 22 đúng và commit lock, nhưng dọn `.task8-build-2..4` mà không
    publish release immutable tương ứng. Task 10 vì vậy chỉ còn client release version 21 và
    buộc dừng `E2E-ASSET-006`. Mỗi lần runtime/map bump phải giữ nguyên full
    `asset-lock.json + server/ + client/`, publish theo thư mục lock SHA và ghi mapping
    lock→release path trong report trước khi dọn build tạm. Lượt phục hồi 01/08 đã áp dụng
    đúng cổng này: giữ cả A/B, publish release `task8-releases-22/54a638...` trước mọi cleanup,
    rồi băm lại toàn bộ release và marker snapshot.

32. **Gửi cho client một thứ mà không xử lý phần trả lời của nó là TỆ HƠN không gửi.**
    Waypoint (cổng rời map) là ví dụ hoàn hảo: client tự khoá mình lại **ngay khi
    bấm cổng** — `Char.ischangingMap = true`, `isLockKey = true`, hiện hộp thoại
    "Xin chờ" (`Waypoint.cs` / `Char.cs:2322`) — rồi **đợi máy chủ**. Nếu máy chủ
    không xử lý gói `-23`, nhân vật **đứng hình vĩnh viễn**, phải thoát game.
    Nghĩa là "chỉ gửi waypoint xuống thôi, đổi map làm sau" sẽ đẻ ra một lỗi nặng
    hơn cái nó định sửa. Cách chữa: mọi nhánh của handler `-23` đều phải trả một
    gói — thành công thì `-22` + `-24`, thất bại thì `46` (đặt lại vị trí tại chỗ,
    client tự mở khoá). Nhớ tìm cái "van mở khoá" đó **trước khi** viết tính năng.
33. **Một tính năng của client có thể cần HAI gói ở hai thời điểm khác nhau.**
    Bóng đổ (bg-item) cần `-31` (**bảng mẫu**: ảnh nào, lớp nào — gửi một lần mỗi
    phiên) *và* ô bg-item trong `-24` (**chỗ đặt**: chỉ số mẫu + ô tile). Thiếu
    `-31` thì `TileMap.getBIById` trả null và client `continue` — **không vẽ,
    không lỗi, không log**. Nhìn game chỉ thấy "bóng đổ không hiện", không có gì
    chỉ về nguyên nhân. Luật: khi client tra cứu theo *id*, phải tìm xem cái bảng
    id đó được nạp từ gói nào.
34. **`vsData` KHÔNG phải công tắc chung cho mọi tài nguyên client.** Nó chỉ gác
    `part.bin` (khoá RMS `NRdataVersion`). Ảnh bg-item đi đường khác hẳn: cache RMS
    theo khoá `x{zoom}bgItem{idImage}` và cache đó **chỉ được ghi khi máy chủ gửi
    gói `-32`** (máy chủ Go không gửi) — nên ảnh mới nạp thẳng từ `Resources`,
    **không cần tăng `vsData`**. Quan trọng hơn: `vsData` đang bị **ghim vào
    `asset-lock.json`** (`TestRuntimeServerNhungDungByteVaSHATuLock`,
    `server/internal/game/dongbo_test.go:22` — phần ghim `vsData` ở dòng 55-64) nên tăng bừa
    là test đỏ và phải sửa artifact đã ký. **Nhưng** chính cái cache theo `idImage`
    đó lại là lý do phải chọn số hiệu ảnh **≥ 500**: bản Ngọc Rồng gốc dùng
    ~0..350, trùng vào là lấy nhầm ảnh Dragon Ball trong cache của người từng chơi
    server cũ — lỗi im lặng tuyệt đối.

---

### Về sprite nhân vật nhiều mảnh (05/09 chiều — khe đầu–thân)

- **Khe bbox âm vẫn có thể là "hai khuôn mặt".** `ghep-khung.py` chỉ đo mép hộp đục; mảnh thân chứa mặt của khung
  khác đè lên đầu cho khe −1…−4 mà mắt thấy hai mặt. Cửa nghiệm thu phải thêm bộ dò số vùng da liên thông trong hộp đầu
  (= 1) và nhìn tấm 33 cf bằng mắt; đo cả x4 (x1 tròn 0 có thể là +2 ở x4).
- **Cắt mảnh theo % chiều cao khung chỉ đúng khi mọi khung cùng tỉ lệ và cùng tư thế.** Ảnh AI sinh từng dải rời
  thì không bao giờ như vậy → cắt theo mốc giải phẫu (cằm, tâm mặt) đo từ ảnh, và chuẩn hoá tỉ lệ từng khung trước.
- **Đừng kết luận "không bù được bằng dữ liệu" khi chỉ nhìn một loại biến.** Điều phối chỉ xét 3 dy đầu (vô nghiệm)
  và bỏ qua 17 dy thân; agent phản biện giải MILP ra nghiệm. Bài học: liệt kê HẾT bậc tự do trước khi đẩy việc sang
  "quyết định mỹ thuật".
- **NPC dùng cùng bảng xương `CharInfo` với người chơi** (`Npc.cs:305–313`) — sửa sprite người chơi thì phải đo cả NPC;
  Iruka đang đầu lửng +13 px vì part 2189–2191 dx=dy=0.
- **Đọc ý nghĩa cf trong `Char.cs` trước khi gán khung:** chạy là cf 2..6 (không phải 3..6), cf 20/31 là thế bắt đầu skill
  (không phải nằm), chết = đứng + vòng sáng, slot đầu 1 = đầu nghiêng. Gán khung theo tên đoán (`REQUIRED_CLIENT_CF`)
  đã sai ở 3 chỗ.

### Về ảnh khoá màu (chroma) — 06/09

- **Tách nền chỉ flood từ MÉP ảnh thì khe kín không bao giờ sạch.** `artqa.extract_alpha` lan từ viền vào;
  khe giữa cánh tay và thân, giữa dải khăn và tóc là ĐẢO KÍN nên giữ nguyên magenta đục — lên tới màn hình
  người chơi thành chấm hồng. Muốn sạch phải quét thêm đảo trong lòng hình (`remove_large_chroma_islands`).
- **Đừng tin chú thích cũ, hãy cắt ảnh ra xem.** Chú thích trong `task6b_closure` nói đảo 586 px là "accent
  bông tai" nên cố ý giữ; cắt vùng đó ra xem thì là khe khăn trán Iruka. Một câu chú thích sai giữ nguyên lỗi
  suốt nhiều đợt.
- **Trước khi xoá màu hàng loạt, chứng minh màu đó không phải mỹ thuật thật.** Histogram: mọi pixel tím trong
  ảnh nguồn đều có kênh lục 1–5 (viền răng cưa của nền), vải tím thật sẽ có lục 40–80 và thành mảng.
- **Test đọc ẢNH ĐÃ SINH không bảo vệ được MÃ SINH.** Đột biến tắt hẳn bước làm sạch mà 33 test vẫn xanh, vì
  không test nào sinh lại rồi so với file trong repo. Cửa kiểu này phải có một test tái sinh (bọc `skipUnless`
  khi máy không có kho ảnh nguồn).
- **Test chép cứng ngưỡng là test tự lừa.** Bản đầu chép bốn ngưỡng chroma vào test rồi "chốt" bằng cách gõ
  lại đúng biểu thức đó — đổi ngưỡng trong test thì test vẫn xanh dù ảnh có pixel magenta thật. Nay ngưỡng nằm
  ở `artqa.chroma_mask` (một nguồn sự thật) và oracle là **bảng màu chép tay** với kết quả mong đợi.
- **“Sáng + trung tính + nối từ biên” vẫn chưa đủ để gọi là nền caro.** Task 2 UI round 3 đã có guard
  trắng/xám toàn biên nhưng mask vẫn flood cả highlight `(244,244,244)` của artwork chạm mép. Muốn xoá an toàn
  phải chứng minh **hai màu trội lặp theo cả hai trục**, rồi chỉ mask đúng cặp màu đó; pixel còn mơ hồ thì giữ
  lại (fail-closed), không đánh cược bằng việc xoá mỹ thuật thật.
- **“Có nhịp caro trên bốn mép” cũng chưa chứng minh ruột ảnh là nền caro.** Round 4 giữ được highlight nhưng
  detector chỉ đếm transition ở biên, nên ảnh 120×96 có viền `242/255` và ruột `242` đặc vẫn qua rồi bị xoá sạch
  11.520 pixel. Cửa mask phải kiểm **cadence + pha + tỉ lệ đúng màu + độ phủ trên toàn mặt ảnh**; kiểm biên chỉ
  đủ cho backup cũ đã resize dùng đọc lineage, không đủ quyền sửa một input đục mới.
- **Unity MCP mới bọc kết quả `execute_code` trong JSON `data.result`; tìm marker thẳng trên stdout sẽ đọc trúng
  chuỗi đã escape và báo state hỏng dù Unity trả đúng.** Phải parse envelope JSON trước, lấy `data.result`, rồi
  mới parse marker; giữ fixture đúng shape của MCP thật để lần nâng package sau không làm cổng backup hỏng lại.

## 7. Nợ bảo mật — phải xử trước khi mở công khai

Gom lại thành **Kế hoạch B**, sửa lẻ từng cái là vô nghĩa:

- ~~**Chưa chống dò mật khẩu** (đo thật: 200 lần sai liên tiếp không bị chặn, 76 lần/giây)~~ ·
  ~~**rò rỉ thời gian ~13,5 ms** giữa nhánh "không có tài khoản" và "sai mật khẩu"~~ ·
  ~~**chưa giới hạn số kết nối trên một IP**~~ → **ĐÃ VÁ 28/08** (`3aea99d`, **chỉ phía game**).
  Gói mới `server/internal/chanip` (chỉ thư viện chuẩn, khoá ngắn, tự dọn rác — theo khuôn `antoan`).
  - **Chống dò:** túi thẻ đếm theo **địa chỉ IP**, **KHÔNG theo tên tài khoản** — đếm theo tên thì
    kẻ xấu cố tình gõ sai là **khoá được tài khoản người khác**, biến lớp bảo vệ thành vũ khí bắt nạt.
    10 thẻ, hồi 1 thẻ/20 giây, chặn **trước** khi chạm CSDL/hàm băm. Vẫn trừ thẻ **kể cả khi tên không
    tồn tại** — nếu không, *số thẻ còn lại* thành chuông báo "tên này có thật", tức mở lại lỗ hổng
    thời gian bằng cửa sau. Đo lại: **200 lần sai → chỉ 10 lần tới được CSDL**; dò đủ 200 lần mất
    **1h03m** thay vì 2,1 giây.
  - **Rò rỉ thời gian:** nhánh "không thấy tài khoản" chạy thêm một lần băm **giả** (hằng số sinh sẵn,
    `auth.CanBangThoiGian`). Phép đo **tự hiệu chuẩn nhiễu** bằng cách chạy trên hai nhánh giống hệt nhau:
    **nhiễu 478 µs · chênh lệch thật 106 µs** — nhỏ hơn nhiễu 4,5 lần, dụng cụ không phân biệt được nữa.
  - **Kết nối/IP:** mặc định 32, trả chỗ qua `defer`. Thiếu gác thì `net.PhucVuTren` **từ chối chạy**.
  - Ba con số chỉnh được bằng biến môi trường (`NRO_CHONGDO_SUCCHUA`, `NRO_CHONGDO_HOIPHUC_GIAY`,
    `NRO_TOIDA_KETNOI_MOI_IP`) — nhà mạng VN dùng NAT lớn, ngày mở cửa có thể phải nới ngay.
- ⚠️ **CÒN HỞ, đã biết:** (a) **tấn công phân tán** (botnet nhiều IP) vẫn lọt — muốn chặn cần captcha
  hoặc xác minh email, game chưa có cả hai.
- ~~(b) chưa chạy được `go test -race` (máy không có `gcc`)~~ → **ĐÃ CÓ 03/09**: chủ dự án duyệt, điều phối cài
  **w64devkit v2.9.1** (gcc 16.2.0) tại `D:\w64devkit\bin` — giải nén, không đụng hệ thống, **không nằm trong PATH
  toàn máy**. Chạy race: `$env:PATH="D:\w64devkit\bin;"+$env:PATH; $env:CGO_ENABLED="1"; go test -race ./...`.
  Lần chạy đầu 03/09: 7 gói ngoài `game` **0 data race**; gói `game` lần 1 có **1 race — ở phía test** (mã sản
  xuất đúng, xem §8), sửa `afda280` → lần 2 **0 race, ok 176s**. **Toàn bộ 8 gói sạch race** tính đến `afda280`.
  Hai test đồng thời viết sẵn từ 28/08 nay đã chạy thật dưới race.
- ~~**rò rỉ thời gian phía `web/`**~~ → **ĐÃ VÁ 03/09** (gộp `7184b95`). `web/lib/matkhau.ts`: thêm
  `BAM_GIA` (băm PBKDF2 THẬT của một mật khẩu ngẫu nhiên, 100.000 vòng, dán cứng) + `canBangThoiGian()`;
  `web/app/dang-nhap/actions.ts`: nhánh "không tìm thấy tài khoản" nay LUÔN chạy `canBangThoiGian()` thay
  vì `||` đoản mạch trả ngay — y hệt `server/internal/auth/password.go`. Đo bằng `rothoigian.test.ts`:
  **trước ~10,5 ms chênh lệch, sau ~0,06–0,10 ms**. Khác bên Go: ở tầng web+PostgreSQL tín hiệu sau vá
  (~60–100 µs) SÁT nhiễu đo (~30–90 µs), không luôn dưới nhiễu như Go (106/478 µs); test dùng thêm
  ngưỡng tuyệt đối 0,5 ms làm tiêu chí chính. Test đột biến (khôi phục `||`) giết được ngay (10,5 ms).
- ~~Web: JWT **không thu hồi phiên** khi tài khoản bị khoá / hạ quyền; đổi mật khẩu chưa
  kiểm cờ `bi_khoa`~~ → **ĐÃ VÁ 27/08** (`9492315`). `docPhien()` nay hỏi CSDL ở **mọi**
  lần đọc phiên và trả dữ liệu **từ CSDL**, không từ payload token; thêm cột `the_he_phien`
  (migration `0003`) để đăng xuất/đổi mật khẩu giết mọi bản sao token cũ. Giá: **đúng 1 truy
  vấn cho mỗi request có cookie**, 0 cho khách vãng lai (đo thật bằng `idx_scan`, không ước
  lượng). Phép kiểm `laAdmin` ở `proxy.ts` **đã gỡ có chủ ý** — tầng Edge không đọc được CSDL
  nên nó chỉ sai về phía *cho qua*; chốt thật là `AdminLayout` + từng server action.
  Đánh đổi đã chọn: đăng xuất một thiết bị = đăng xuất mọi thiết bị.
- `npm audit` 3 lỗi high (nằm trong bundle Next, chưa phải bề mặt tấn công).
- Di sản (nếu còn dùng phần cũ): chưa có HTTPS, chưa sao lưu DB tự động, **không được mở
  cổng 3306 ra internet**.

---

## 8. Nợ kỹ thuật lặt vặt (không chặn việc gì)

- **Client sau khi chết không có đường về (06/09):** máy chủ chưa gửi gói hồi sinh/về làng khi hp = 0 → client `statusMe = 14`, hộp "Xin chờ" và ~1–2 phút sau rớt về màn chọn máy chủ (watchdog `timeBreakLoading` 30 s của `GameCanvas`). Cần gói chết/hồi sinh theo giao thức (KH8 boss đã có sát thương nhưng chưa có luồng chết người chơi).
- **Chụp màn hình qua MCP khi Unity không có focus:** Game view không tự vẽ lại → `ScreenCapture` trả khung cũ (đã thấy màn "Xin chờ" cũ thay vì map). `cs/chup.cs` và `chup-loat.py` nay gọi `InternalEditorUtility.RepaintAllViews()` trước khi chụp. Bấm Play ngay sau `refresh_unity` khi còn nhập lại 148 ảnh cũng làm phiên chơi bị reset về màn chọn máy chủ — chờ nhập xong rồi mới Play.
- **Di chuyển nhân vật qua MCP:** trên PC client đọc `keyHold[24]`/`[23]` (phải/trái), `[21]` (nhảy) trong chuỗi else-if (`GameScr.cs:2708–2845`); đang giữ nhảy thì phím trái/phải không được xử lý; lần giữ đầu chỉ quay mặt, lần sau mới đi; kẹt vách phải nhảy. Kịch bản: `cs/chay-phai.cs`, `chay-trai.cs`, `nhay.cs`, `tha-phim.cs`.
- **Client Unity đã bị sửa tay và không có git (05/09):** `Game1/Char.cs:6013` có `Debug.Log("[POC] ...")` chèn 22/08;
  `Game2/Char.cs` là bản 16/07 (BOM/CRLF khác nhau). Mọi phương án đụng `Char.cs` phải sao lưu vào repo trước.
- **`ch` = 34 + dy_đầu0 + 10 = 5 px (`Char.cs:6087`)** — chiều cao nhân vật client dùng để neo số sát thương/chat/mũi tên
  giả định sprite Dragon Ball; với đầu dy −39 nó thành 5 px → chữ có thể neo gần bàn chân. Chưa nhìn trong game.
- **`hatInfo` chỉ 32 phần tử (`Char.cs:1183`)** → cf 32 lỗi chỉ số bị `try/catch` nuốt (chỉ lộ khi có nón).
- **Iruka part 2189/2190/2191 dx=dy=0** → khe +13/+2, bàn chân +29 px dưới neo; pipeline Iruka (`task6b_closure.py:679–730`)
  ép vừa 128×160 không qua `component_placement`. 10 sprite 27075–27084 trong client không khớp SHA file nào trong repo.
- **Baseline pytest `assetpack tools`: 183 xanh / 8 đỏ SẴN (tools/artqa/test_vertical_slice_generator.py, 'ASSET-ART-004: sprite 27078 raw source missing'), 358 s — đo 05/09** — ghi để không lẫn đỏ cũ với đỏ mới.
- **Interface Go bọc con trỏ nil KHÔNG phải nil (05/09):** `moE2EControl` trả `(*http.Server)(nil)` khi
  control tắt; gán vào interface `mayChuCoTheTat` thì `!= nil` vẫn đúng và `Shutdown` panic → cả trình tự tắt
  (LƯU NGƯỜI CHƠI) bị bỏ qua, im lặng, **ở mặc định**. Luật: hàm trả con trỏ có thể nil thì chỗ gọi chỉ gán vào
  interface khi đã kiểm `!= nil`; và hàm nhận interface kiểm thêm kiểu cụ thể nil (`coControlThat`). Test phải
  truyền đúng kiểu nil-có-kiểu, không phải `nil` trần.
- **Client Unity là bản port Dragon Ball: 3 "gender", không phải 5 làng (05/09).** Mọi bảng tra theo `cgender`
  bên client dài 3 (`setDefaultBody/Leg`, `info1.charId[gender]`). Byte "gender" và byte "lớp/làng" trong ME_INFO
  và `-5` là HAI vai khác nhau — gender luôn 0, lớp = làng. Ai thêm làng/đồ mới phải đọc `hinhdang.go` trước.
- **"Chơi được" ≠ "đúng": -42 kẹp máu về đầy che mất nhân vật chết trong CSDL.** Kẹp ở tầng gói tin để HUD
  không 0 là hợp lý, nhưng phải sửa DỮ LIỆU (hồi sinh lúc vào) chứ không chỉ sửa hiển thị — hai bên nói hai sự
  thật là lỗi khó thấy nhất khi chơi thử.
- **Pipeline ảnh: co từng khung cho vừa hộp = nhân vật đổi cỡ theo khung (05/09).** `fit_to_canvas` với
  `min(w/W, h/H)` cho MỖI khung sprite nhân vật → khung đấm/đá rộng bị co nhỏ, khung đứng thì không → đầu (3 khung
  idle) không khớp thân/chân. Sprite của cùng một nhân vật phải dùng **một tỉ lệ**, neo **một điểm** (bàn chân), lệch
  còn lại đưa vào `dx/dy` từng khung. Đo bằng `tools/nhan-vat/ghep-khung.py` trước khi tin mắt.
- **Client "im lặng bỏ qua" = phải giải mã gói THẬT theo code client, không tin spec lẫn test byte (05/09).**
  Spec §3.6 gọi ô 6 của quái là `tempId`, test byte ghim đúng theo spec, máy chủ ghi đúng theo test — cả ba
  đồng ý với nhau và cùng SAI so với `Controller.cs:5068` (client đọc là `templateId`). Hậu quả ẩn suốt 6 tuần
  (NPC không hiện, bóng đổ không vẽ) vì client nuốt ngoại lệ trong `try/catch`. Cách bắt: ghi hex gói ở mức Debug
  (một dòng mỗi lần vào map, không phải mỗi gói) rồi viết bộ giải mã lặp lại ĐÚNG thứ tự đọc của client
  (`scratchpad/giai-ma-24.py`) — sai lệch lộ ra ở đúng byte. Áp dụng cho mọi gói client "không kêu": -24, -30, -5, -44.
- **Test tích hợp CSDL: tiền tố tên tài khoản phải RIÊNG cho từng file test.** Hàng rác `chet8469` (03/09, dọn
  hỏng lúc PostgreSQL tắt) làm test mới cùng tiền tố `chet` đỏ "duplicate key" ở lần full đầu 05/09.

- **Hai agent làm song song cùng package → gộp xong PHẢI build + test lại main TRƯỚC khi coi là xong.**
  03/09: task 3 (`goihanhtrang.go`) và task 8 (`cuahang.go`) mỗi nhánh tự xanh, gộp cả hai vào `main` thì
  **`tranParamOption redeclared`** — hai agent cùng khai một hằng 65535 cho cùng một sự thật của client.
  Main đứt build ở `b4f83e8`, vá ở `ba3a56d` (gom về `mauvatpham.go`). Luật rút ra: hằng số mô tả *giao
  thức client* (trần byte, unsigned…) chỉ khai **một nơi** ở file gốc chung; brief cho agent song song phải
  liệt kê tên hằng/kiểu đã có để dùng lại. Cùng ngày một agent còn **lỡ commit vào nhánh của agent kia** vì
  hai bên đổi nhánh trong cùng cây làm việc → agent làm song song phải dùng **worktree riêng** ngay từ đầu.
- **`TestTichHop_DanhQuaiConSong` (`khuvuc_tichhop_test.go:682`) chớp tắt dưới tải:** agent task 8 thấy nó
  trượt hạn 5s khi chạy `go test ./...` toàn bộ, chạy riêng 3 lần đều xanh. Test nhạy thời gian, chưa sửa.
  **Cùng họ:** `TestTichHop_DanhQuaiThuCungDocDu5ByteKhongLechDong` đỏ một lần (9,3s) ngay sau gộp task 12 (`583377e`),
  chạy riêng `-count=3` xanh cả ba (3,4s) → KHÔNG phải hồi quy của task 12. Hai test này cần một việc nhỏ riêng: nới hạn
  chờ hoặc thay đợi-thời-gian bằng đồng hồ tiêm. Người gộp: **gặp đỏ ở hai test này thì chạy riêng trước khi kết luận.**
- **Đột biến tay phải SẠCH và phải bỏ ĐỦ các lớp chặn.** 03/09 điều phối áp đột biến D4 lên `ChuyenO` bằng
  cách thay `chu_so_huu_id = $1` thành `$1 = $1` → test đỏ, nhưng đỏ ở **happy-path** (pgx không suy được
  kiểu `$1`) chứ không phải test D4 — vì D4 còn một lớp nữa ở câu `SELECT … FOR UPDATE`. Làm lại: giữ kiểu
  (`$1::bigint IS NOT NULL`) và bỏ guard ở **cả hai** câu → đúng test "người khác không đụng được đồ" đỏ.
  Luật: một đột biến "giết được" chỉ có nghĩa khi **đúng test mang tính chất đó** đỏ và các test khác vẫn xanh.
- **Đừng sửa file `.go` có tiếng Việt bằng `Set-Content -Encoding utf8` của PowerShell** — nó làm hỏng dấu
  (agent task 4 gặp, phải khôi phục). Dùng Edit hoặc Python `io.open(..., encoding="utf-8")`.
- **Điểm mù oracle trong test byte KH7 (agent KH8 phát hiện):** `mauvatpham_test.go:47` dựng `want` từ chính
  `DanhSachMauVatPham()` → test chỉ ghim **cấu trúc** gói, không ghim **giá trị** — đổi số trong danh mục test vẫn
  xanh. `muaban_tichhop_test.go:117` cố ý dùng cửa hàng riêng của test nên cũng không khoá số thật. Không phải lỗi
  code, nhưng là lý do KH8 task 0 phải là GATE "YAML sinh ra đúng byte với dữ liệu Go hiện tại". Khi làm KH8, ghim
  một bản byte vàng của danh mục hiện tại trước khi chuyển.
- **Script Python gọi `subprocess` trên Windows phải đặt `encoding="utf-8"`** — mặc định GBK, output git có tiếng
  Việt làm sập luồng đọc (03/09, script nhật ký của điều phối). Đã xảy ra một lần, bỏ sót một dòng nhật ký phải bù.
- **Test KHÔNG được `time.Sleep` rồi đọc thẳng trường của `BoXuLy`** — lần `go test -race` đầu tiên (03/09, sau khi
  có w64devkit) bắt `TestNhanClientOk` đọc `b.clientSanSang` trong khi goroutine phiên đang ghi (`login.go:490`).
  Mã sản xuất ĐÚNG (chỉ goroutine phiên ghi); lỗi ở test. Cách đúng: `dungPhienThat` nay trả `dong()` **đóng ống
  rồi ĐỢI `PhucVu` thoát** (hạn 5s, `sync.Once`) — `net.Pipe` đồng bộ nên `WriteMessage` trả về là phiên đã đọc
  xong khung; gọi `dong()` trước khi đọc trường = happens-before thật. Sửa ở `afda280`. Test đo thời gian
  (`rothoigian_test.go`) **bỏ qua dưới race** qua cặp `race_bat.go`/`race_tat.go` (`raceEnabled`) — dưới race
  PBKDF2 chậm ~10× nên nhiễu 2,8ms > ngưỡng 1ms, là dụng cụ đo bị nhiễu chứ không phải rò rỉ.
- **Agent làm trong worktree vẫn có thể để lại file lạc ở CÂY CHÍNH** — 03/09 gộp KH8 task 3 bị git chặn vì
  `server/internal/game/khochieu_hopdong_test.go` **0 byte, chưa theo dõi** nằm ở `D:\NRO_GoiBanGiao` (agent lỡ
  tạo rồi bỏ; file rỗng làm cả gói `game` ở main không biên dịch được). **Lặp lại ngay task 10** (`kich.go` 0 byte
  + `dondep_test.go` sửa dở 3 dòng ở cây chính). Nguyên nhân: agent gọi Write/Edit bằng đường dẫn **tương đối**
  trước khi `cd` vào worktree → tạo file ở cây chính rồi mới làm thật ở worktree. **Luật cho brief:** mọi Write/Edit
  dùng **đường dẫn tuyệt đối vào worktree** (`D:\NRO_wt_…\server\…`). Người gộp: trước khi merge chạy
  `git status --short` ở cây chính, file `??`/` M` không phải của mình → so với bản trên nhánh (hash/kích cỡ) rồi
  `git checkout --` / xoá; KHÔNG commit vào main.
- **"Hai nhánh tự xanh, gộp lại đỏ" lần 2 — danh sách file bắt buộc nằm rải trong fixture test.** KH8 task 8 thêm
  `quai.yaml` vào `Nap`; task 10 (song song) viết fixture "chép kho nhúng ra thư mục tạm" với danh sách file **cứng**
  → mỗi nhánh xanh, gộp xong 4 test nạp nóng đỏ vì thư mục tạm thiếu `quai.yaml` (`672640c` → vá `b54a51a`).
  Luật: thứ gì là "danh sách đầy đủ" (file YAML, cột CSDL, hằng client) thì **một hàm duy nhất** trả ra
  (`noidung.TenCacFile()`), mọi fixture/test duyệt hàm đó — thêm mục mới chỉ sửa một chỗ. Người gộp: sau khi gộp
  agent thứ hai trở đi trong cùng đợt song song, **luôn chạy full test trước khi kết luận** (đã là luật, nay có ví
  dụ thứ hai).

- **Bộ chạy đột biến PHẢI có timeout từng con, và đừng chờ "mọi python.exe tắt".** 02/09 agent soi (Opus) chạy 28 đột biến bằng `subprocess.run` **không timeout** + vòng `until ! tasklist … python.exe` → kẹt **19 giờ** trong khi 8/28 kết quả đã nằm sẵn trong file output. Treo **không** do đột biến M9 (chạy lại xong 104s) mà nghi do PostgreSQL rớt lúc 21:50 + kết nối CSDL `context.Background()` không hạn giờ (nợ ngay dưới). Người điều phối: **đọc output tiến trình con trước khi kết luận agent kẹt**.
- **Agent model rẻ (Sonnet / Explore) hay làm xong mà KHÔNG commit** — 02/09 hai lần liên tiếp: agent Explore không có quyền `git commit` (đúng vai chỉ-đọc), agent Sonnet task 14 thì quên. Việc không mất (nằm trong worktree) nhưng người điều phối phải **kiểm `git status` trong worktree** trước khi tin báo cáo, và commit thay. Brief cho model rẻ nên nói thẳng *"bước cuối: `git add -A && git commit`, báo lại hash"*.
- ~~**`assetpack/tests/test_task6b_canonical_preflight.py` chưa bao giờ chạy**~~ → **ĐÃ VÁ 27/08.**
  `SHARED_ROOT` bây giờ hỏi git `rev-parse --git-common-dir` rồi lấy thư mục cha, thay vì đếm
  `ROOT.parents[1]`. **Tại sao phải hỏi git chứ không dùng luôn `ROOT`:** `docs/my-thuat/anh-nguon/`
  và `anh-game/` bị `.gitignore` loại, nên một **worktree không có ảnh nguồn** — chỉ bản checkout
  chính (nơi chứa `.git` thật) mới có. `--git-common-dir` trỏ đúng chỗ đó ở **cả hai** kiểu checkout.
  Chạy được rồi thì **cả 6 chốt kiểm đều XANH** (không cái nào ghim thứ đã trôi), `kiem-tra-tat-ca.ps1`
  đã bỏ `--ignore`.
- ~~**`server/internal/game/assets/tile_set_info.bin` là quả mìn CÙNG LOẠI**~~ → **ĐÃ VÁ 27/08**
  (worktree `agent-a31f4c5affe76a227`, chờ gộp). Build sinh bảng gốc **1696 byte** (chỉ khai tile
  1–39) trong khi bản đang chạy là **1734 byte** đã sửa tay để thêm tile Làng Lá 60–69/83–88 →
  build rồi chép sang máy chủ = `types[]` = 0 → **rơi xuyên đất xuống y=576**. **Cách vá:** phép
  "thêm tile Làng Lá" đưa hẳn vào đường ống — khai báo trong `mapping-id.json`
  (`build_inputs.tile_set_overlay`: 4 nhóm bit 2/4/8/0x2000 + `result` golden 1734 byte) và
  `assetpack/tileset.py` thi hành. **Chọn hướng này chứ không thêm file .bin nguồn mới** vì
  `docs/my-thuat/anh-nguon/` bị gitignore (file mới sẽ không được commit), và vì để nguồn vẫn là
  MỘT: bảng gốc đã qua review giữ nguyên, `tile_set_baseline` + `production_closure_sha256`
  không đổi. Chốt chặn cuối `validate_floor_closure` đỏ **kể cả khi `tile_set_overlay` bị xoá
  hẳn**, nên quả mìn không dựng lại được. Nghiệm thu: build thật → **cả 5 artifact giống hệt
  từng byte** bản đang chạy; 4 đột biến đều bị giết.
- ⚠️ **`versions.data` lệch giữa manifest và bản đang chạy — PHÁT HIỆN 27/08, CHƯA VÁ.**
  `mapping-id.json` ghi `versions.data = 21`; `server/internal/game/assets/asset-lock.json` và hằng
  `vsData` trong `server/internal/game/dongbo.go` đều ghi **26**. Ai build + deploy sẽ **lùi version
  dữ liệu 26 → 21**, mà chú thích ngay tại `dongbo.go` nói rõ "mỗi lần `part.bin` đổi phải tăng,
  không client dùng bản cũ trong RMS". **Khác quả mìn tile ở chỗ nó KÊU:**
  `TestRuntimeServerNhungDungByteVaSHATuLock` so hằng Go với lock nên sẽ đỏ ngay — hỏng ồn chứ
  không hỏng im. Vá đúng cần chốt version nào là thật rồi sửa manifest + version-ledger cùng lúc.
- ⚠️ **`docs/my-thuat/lang-la-client-assets/res/mymap/1.bytes` đã cũ 51 ô** so với bố cục client
  thật (`D:\NRO_Project\...\mymap\1.bytes`). Đây là "bản sao byte thật" chụp hồi 22/08, còn bố cục
  bị sửa tiếp sau đó (`6fac494` xếp bố cục dày đặc). Kích thước vẫn đúng 52×24 nên chốt kiểm mới
  chỉ dám so **header**, không so ruột. Muốn so được ruột thì phải chụp lại bản sao này.
- **Một đột biến của `doimap.go` KHÔNG test nào giết được — và không thể giết bằng test
  giá trị.** Đột biến: giữ nguyên bước chụp ảnh từ actor rồi vẫn đọc `x, y` từ
  `b.nhanVat.ViTriX/Y`. Lý do bất khả: actor và luồng phiên **dùng chung một con trỏ**
  `*store.NhanVat` (`k.nguoiTrong[id].nv` chính là `b.nhanVat`), nên sau vòng hỏi-đáp chụp
  ảnh hai cách đọc **luôn ra cùng một giá trị**. Nó thuần tuý là **data race**, chỉ máy dò
  đua (`go test -race`) mới bắt được — mà máy này **không chạy được `-race`**:
  `CGO_ENABLED=0` và không có `gcc`. Đã có test giết được dạng *mạnh* của đột biến này
  (`TestTichHop_DoiMapXetCongBangViTriLayTuActor`: bỏ hẳn bước chụp ảnh thì đỏ). Muốn bịt
  nốt thì phải cài trình biên dịch C rồi thêm một lượt `go test -race` vào CI.
- `Dung()` gọi trước `BatDau()` sẽ panic (cạnh biên, chưa ai gặp).
- Client gửi `-39` trước `-28/13` → không làm gì (im lặng, đúng nhưng chưa có test).
- Đăng nhập trùng tài khoản ghi đè người cũ trong khu vực (có từ trước).
- `byte(nv.Cap)` cắt cụt nếu cấp > 255.
- `context.Background()` ở `login.go`: truy vấn CSDL không hạn giờ, không huỷ khi client ngắt.
- Vòng `Accept()` không nghỉ giãn → cạn file descriptor sẽ quay nóng CPU.
- Mã chết: ~~`Writer.Bool`~~ (**hết chết** — nhánh waypoint dùng nó ghi hai cờ của cổng,
  `thegioi.go:206`), `Writer.Long`; `TaiKhoan.LaAdmin` đọc từ CSDL nhưng chưa ai dùng.
- Test package `game` chậm (~70 giây) do nhiều test đầu-cuối chờ hết hạn giờ đọc.
- ~~KH6 Task 13: quái không tôn trọng bất tử 1500ms~~ → **đã vá ở Task 11** (`chonMucTieu`
  thêm tham số `bayGio` + lọc `lucHoiSinh`, có test thuần + test nhịp thật).
- **KH6 — hồi máu PvP hoãn có kỷ luật:** `goiHoiMauNguoiChoi` (gói `-60` dạng hồi máu, 11
  byte) đã dựng + test byte nhưng **cố ý CHƯA nối dây** — cơ chế hồi máu chưa đặc tả đủ
  (lượng hồi bao nhiêu; người được hồi biết máu mới qua gói nào — `-30/9` §3.4 chỉ gửi
  "người khác"). Chiêu type=2 (hồi máu) hiện có **chốt không xử như đòn đánh** để không gây
  sát thương sai (`khuvuc.go`, `case suKienDanhNguoiChoi`). Khi có đặc tả hồi máu (KH riêng
  hoặc cùng PvP-có-luật KH9) thì nối `chieuDangChon.Loai==2` → `goiHoiMauNguoiChoi`.
- **KH6 — `iconId` sinh theo vị trí** (`chieuthuc.go:166` `int16(1+lang*4+slot)`): bug lớp
  "chỉ số = danh tính" (§2.3 thiết kế tổng) — chèn làng/đổi thứ tự chiêu làm lệch icon các
  chiêu sau. Chưa cắn (5 làng cố định). KH7/KH8 thay bằng bảng tra ID bền, thêm ID ở cuối.
- **KH6 — hai test chập chờn dưới tải cao**: `TestTichHop_DanhQuaiVuotTam...` và
  `TestTichHop_DanhQuaiThuCung...` (`khuvuc_tichhop_test.go`) thỉnh thoảng hết giờ khi
  chạy `-count=5` trên máy đang bận — quái AI thời gian thực giết nhân vật thử nghiệm
  trước khi nó kịp đánh trả. Đã xác minh **có từ trước KH11** (chạy trên commit gốc tái
  hiện y hệt). Cách chữa: cho hai test đó dùng khu vực không quái mặc định.
- **KH11 — khe hở lúc tắt máy chủ (~1 ms)**: client rớt mạng đúng lúc `DungTatCa()` đang
  lặp dừng 5 khu vực, mà khu vực của họ **chưa tới lượt** → `KhiDong()` gọi `Xep()` sau khi
  thợ ghi đã dừng → **mất bản ghi cuối của MỘT người** (không phải mất sạch). Vá đúng cần
  hạ tầng "rút hết phiên đang mở trước khi dừng".
- **KH11 — khe hở "rớt mạng tự nhiên trùng đăng nhập trùng"**: phiên cũ tự đứt đúng lúc có
  đăng nhập trùng → `XoaNeuDung` giành khoá trước `DangKy` → phiên mới **không chờ ai**, có
  thể đọc trước khi bản ghi trôi xuống đĩa. **Cùng bản chất** lỗi bản Java (đọc trước, đá
  sau) nhưng khác hẳn mức độ: bản Java sai **100% mỗi lần** trùng phiên; khe hở này cần
  **hai sự kiện độc lập trùng nhau trong vài mili-giây**.
- **KH11 — actor bận > 100 ms thì cả khu vực bị bỏ qua một lượt lưu**: `AnhChupNguoiChoi`
  chờ tối đa 100 ms rồi trả rỗng; `LuuTatCa` bỏ qua nguyên khu vực đó (có ghi `slog.Warn`,
  không câm hẳn). *(Cập nhật 03/08: nhánh actor-kẹt-30s vì tự gửi gói mạng đã được vá — xem
  **§6.10b ĐÃ VÁ** — nên actor giờ chỉ có thể bận vì tính toán gameplay nặng, không vì
  socket. Vẫn nên đo dưới tải thật bằng KH6c stress test.)*
- **KH6 Task 11**: chưa chặn tự đánh chính mình (`-60` với `idMucTieu==idNhanVat`) — vô hại
  ở KH6 (chưa có luật đội/PvP), hoãn tới KH9.
- **Nghiệm thu Unity 28/07 — client Unity đã bị tinker tay, KHÔNG REPRODUCE ĐƯỢC trên máy khác.**
  Nếu bạn (hoặc agent khác) mở dự án Unity `D:\NRO_Project` trên máy sạch, sẽ vấp lại toàn
  bộ chuỗi 4 lỗi kể trong §12 dòng 28/07. Ba việc phải làm để "chính thức hoá" nghiệm thu:
  (a) tạo **5 file tile map** cho làng 1-5 vào `Assets/Resources/res/mymap/{1..5}.bytes` —
  hiện đang là bản sao của `0.bytes` (map trống 52×24, hiển thị cột thang lạ giữa màn);
  (b) tạo **`part.bin` Naruto** thay placeholder Dragon Ball ở `server/internal/game/part.bin`
  — hiện client render nhân vật bằng sprite Dragon Ball nên narruto không hiển thị đúng
  (`dongbo.go:15` ghi rõ "PLACEHOLDER cua ban Dragon Ball cu"); (c) **đóng atlas `Big0..4.png`
  Naruto + `NR_image`** để client render npc/boss/quái/hu/vk bằng ảnh mới. Cả ba đều đòi
  **1 KH riêng có subagent + TDD** (như KH6/KH11) — bộ 343 ảnh Naruto ở `docs/my-thuat/anh-game/x1..x4/`
  đã có sẵn nhưng thiếu (i) mapping ID số → tên file, (ii) script deploy, (iii) format
  binary chưa reverse-engineer đủ chi tiết. Xem thêm ghi chú "Khối 10" ở
  [`BO-SUNG-LO-TRINH.md`](docs/superpowers/plans/BO-SUNG-LO-TRINH.md) khi được viết.
- **Client Unity đã có 3 sửa đổi tay ở `D:\NRO_Project` — GHI VÀO ĐÂY vì repo Unity không
  thuộc gói `NRO_GoiBanGiao`**, phiên sau tra `git log` không thấy:
  (1) `Assets/Scripts/Game1/SplashScr.cs` — thêm `ServerListScreen.waitToLogin = true;` sau
  `GameCanvas.connect()` ở cả hai nhánh của `loadIP()` (line 109-110, 124-125) để bỏ qua
  click "Đăng nhập" tay. Cần rollback trước khi build APK phát hành nếu muốn user gõ tài
  khoản tay;
  (2) `Assets/Resources/res/mymap/{1,2,3,4,5}.bytes` — bản sao của `0.bytes` (map trống);
  (3) RMS `C:\Users\Admin\AppData\LocalLow\Kame\Ngoc Rong Kame\` đã được reset: xoá 16 file
  cache `NR*` (NRdataVersion, NRmap, NRitem*, NR_arrow, NR_dart, NR_effect, NR_image, NR_part,
  NR_skill) + tạo `acc="naruto"`, `pass="1"`, `ResVersion="1"`, `svselect=0`. Nếu chạy trên
  máy khác cần **repeat toàn bộ 3 bước trên** hoặc mở KH mapping ảnh sửa gốc.
- **Server `-74` handler là stub tối thiểu** (`login.go xuLyXinTaiNguyen`): trả `[b42=0][ResVersion=1]`
  bỏ qua token 40 byte client gửi (HMAC-SHA256 anti-bot). Client thật sự **KHÔNG gửi -74 trong
  luồng đăng nhập bình thường** (chỉ gửi khi user bấm nút "Download resource" ở màn chọn
  server — mà luồng auto-login bỏ qua) → handler này thực tế **chưa bao giờ được gọi trong
  runtime hiện tại**. Giữ lại như bảo hiểm, không vá thêm cho tới khi có KH resource sync.

---

## 9. Bẫy môi trường trên máy này

- **Muốn nhìn cấu trúc `server/` nhanh:** `graphify update D:/NRO_GoiBanGiao/server` (skill `dothi`, tree-sitter local, **0 token**, ~10 giây) → mở `server/graphify-out/graph.html` hoặc đọc `GRAPH_REPORT.md`. Thư mục này đã `.gitignore`. Lần đo 02/09: **1.536 node · 4.718 cạnh · 64 cụm · không vòng import**. Hub thật của game là `KhuVuc`/`khuvuc.go` (84 cạnh) — lý do hai agent cùng vào là đụng. `TimBanDo()` 91 cạnh **không phải điểm nóng**: gần hết là test/helper gọi.
- **⚠️ Test cần CSDL TỰ BỎ QUA khi PostgreSQL tắt — và bảng tổng kết vẫn xanh.**
  `kiem-tra-tat-ca.ps1` chỉ đặt `NRO_TEST_DSN` khi thấy cổng 5432/5433 đang nghe. Máy này
  thường tắt PostgreSQL, nên **toàn bộ test tích hợp** (`*_tichhop_test.go`, `store`,
  `migration`) im lặng bỏ qua. Script *có* in cảnh báo ở cuối — **đọc nó**, đừng chỉ nhìn
  dòng `[OK] go test`. Ngày 26/08 bật CSDL lên thì lộ ra **15 test đỏ** đã nằm im nhiều ngày.
  Bật CSDL rồi mới kết luận là xanh:
  ```
  D:\PostgreSQL\bin\pg_ctl.exe -D D:\PostgreSQL\data start
  ```
  Sau lần tắt đột ngột nó cần ~30 giây phục hồi mới nhận kết nối — chờ `pg_isready`.
- **PowerShell 5.1 có thể để `$PSScriptRoot` RỖNG lúc gán giá trị mặc định cho `param`.**
  `param([string]$Manifest = (Join-Path $PSScriptRoot ...))` làm script chết **trước khi**
  chạy dòng đầu tiên. Cách đúng: để tham số mặc định rỗng, rồi gán bằng `$PSScriptRoot`
  **ngay sau** khối `param`. Đã có test chạy script mà không truyền `-ManifestPath`.
  *(Cứu từ nhánh `kh-unity-e2e` — đây là thứ duy nhất nhánh đó có mà `main` chưa có.)*
- **Không chạy được `go test -race`** — thiếu trình biên dịch C (cgo), không có quyền admin
  để cài. Bù lại bằng: **kỷ luật actor** (soi từng đường chạm trạng thái) + **stress test**
  (`-count=3` hoặc `-count=5`). Với mô hình actor thì cách này đủ chắc, vì tính đúng đến từ
  thiết kế "một chủ sở hữu" chứ không từ khoá.
- **Nhiễu CRLF của `gofmt` — đã dọn gốc (04/08).** `.gitattributes` ép `* text=auto eol=lf`
  nên **blob luôn LF** và mọi checkout mới (worktree, CI) ra LF sạch → `gofmt`/CI soi được
  **mọi** file, không cần lọc tay. Nếu một **working tree cũ** còn CRLF trên đĩa (checkout
  từ trước khi có `.gitattributes`), sửa một lần rồi commit: `git checkout -- <file>` (hoặc
  `git add --renormalize .`). Đã thêm ngoại lệ `*.bat`/`*.cmd = eol=crlf` (script Windows).
- **Cổng PostgreSQL KHÔNG cố định giữa các lần khởi động — đừng tin con số ghi sẵn ở
  bất cứ đâu, kể cả file này.** `postgresql.conf` để mặc định `5432`, nhưng ai bật bằng
  `pg_ctl ... -o "-p 5433"` thì nó nghe `5433`. Cách kiểm 1 dòng ở §3 (hoặc
  `server/README.md` mục 0).
  **Triệu chứng khi dùng nhầm cổng — nhận ra ngay:** các gói `auth` / `net` / `protocol`
  **vẫn xanh**, còn mọi kiểm thử cần CSDL **đỏ hàng loạt** với `connectex: ... actively
  refused it`. Rất dễ tưởng mình vừa làm hỏng code, trong khi thực ra chỉ sai cổng
  hoặc PostgreSQL chưa bật (nó **không** chạy như dịch vụ Windows).
  *Vì sao ghi lại:* ngày 26/07 một agent thấy tài liệu ghi `5433` mà máy nghe `5432`, đã
  kết luận nhầm là "tài liệu sai từ đầu" rồi ghim `5432` vào — chỉ dời quả mìn sang lần
  khởi động sau. Bài học thật là **ghim số nào cũng sai; phải dạy cách tự kiểm.**
- **PowerShell không có `&&`** — dùng `;` hoặc `if ($?) { ... }`.
- **Windows dùng `python`, không phải `python3`.** Thứ cần GPU/CUDA thì dùng `py -3.12`.
- Đặt biến môi trường phải là **dòng riêng**, không nối cùng dòng lệnh.

---

## 10. Cách làm việc đã chốt (giữ nguyên cho các kế hoạch sau)

1. **Nghiên cứu giao thức trước** — hai agent đọc độc lập code Java cũ, một agent thứ ba
   soi đối kháng, chỉ chấp nhận khi 0 sai lệch.
2. **Viết thiết kế (spec) có bảng byte** → `docs/superpowers/specs/`.
3. **Viết kế hoạch chia task nhỏ** → `docs/superpowers/plans/`.
4. **Thi công bằng subagent, mỗi task một agent, bắt buộc TDD** (test đỏ trước, code sau).
5. **Mỗi task xong → gói review → một agent khác soi** (không phải agent vừa làm).
   Có lỗi thì vòng sửa, xong mới ghi sổ.
6. **Cả nhánh xong → một agent opus soi toàn nhánh** → phán quyết "gộp được" mới gộp
   `--no-ff` vào `main`, rồi xoá nhánh.
7. **Ghi sổ mọi thứ** vào `.superpowers/sdd/progress*.md` và **cập nhật file này**.

Mã nguồn viết **chú thích tiếng Việt giải thích *tại sao***, không phải *cái gì*.

---

## 10b. Nhiều agent làm song song

| Việc | Dùng |
|---|---|
| Bắt đầu **và** kết thúc một việc | Gõ **`/viec`** — lệnh duy nhất, tự nhận biết: cây git sạch → **nhận việc** (đọc bối cảnh, xem ai đang làm gì, tạo nhánh riêng, báo cáo); có thay đổi chưa commit → **ghi sổ** (chạy test, ghi nhật ký, cập nhật tài liệu, commit). Ép chế độ: `/viec ghi so` hoặc `/viec <việc muốn làm>` |
| Xem ai vừa chạm vào gì | [`NHAT-KY-AGENT.md`](NHAT-KY-AGENT.md) — dòng mới nhất ở **cuối** bảng |
| Mở agent đúng thư mục | Nhấn đúp `MO-AGENT.bat` |
| Agent không phải Claude Code | Cho nó đọc [`AGENTS.md`](AGENTS.md) |

`/viec` **dừng lại hỏi** nếu cây làm việc bẩn mà không phải phiên đó sửa — để không bao giờ
commit lẫn việc dở của agent khác.

**Chống đụng nhau:** mỗi agent một **nhánh git riêng** (`git switch -c <ten-viec>`); không hai
agent cùng commit lên `main`. `server/` (game) và `web/` là hai vùng tách biệt → làm song song
an toàn. Hai agent cùng vào `server/internal/game/` thì **chắc chắn** đụng.

Hook `SessionStart` (`.claude/settings.json`) tự nạp `git log` + nhật ký vào đầu mỗi phiên,
nên phiên mới luôn biết người/agent khác vừa làm gì — **chỉ hoạt động khi mở Claude Code
trong thư mục dự án này.**

---

## 11. Quy ước cập nhật file này

Làm xong một việc đáng kể → sửa **đúng 3 chỗ** rồi commit chung với việc đó:

1. **§4 Đã làm** — thêm dòng, kèm mã commit.
2. **§5 Còn lại** — xoá việc vừa xong, thêm việc mới lộ ra.
3. **§6 Lỗi sai** — nếu vấp lỗi nào mới, ghi lại **kèm lý do vì sao vấp**. Đây là phần
   giá trị nhất, đừng bỏ qua vì "lỗi nhỏ".
4. Thêm một dòng vào **§12 Nhật ký** và sửa dòng "Cập nhật lần cuối" ở đầu file.
5. Thêm một dòng vào cuối bảng **`NHAT-KY-AGENT.md`** (ai làm, nhánh nào, còn dở gì).

Cách nhanh nhất: gõ `/viec`, nó làm cả 5 bước trên.

Commit riêng cho việc cập nhật tài liệu thì dùng tiền tố `docs:`.

---

## 12. Nhật ký

| Ngày | Việc |
|---|---|
| 23/07/2026 | Chốt thiết kế tổng "Naruto trên lõi NRO". Xong **KH1** (giao thức + đăng nhập), gộp `6529d35` |
| 24/07/2026 | Xong **KH2** (nhân vật) gộp `0816c89`; **KH3** (đồng bộ tĩnh, 5 làng Naruto) gộp `7a85431`; giải mã xong giao thức di chuyển |
| 25/07/2026 | Xong **KH4** (vào thế giới) gộp `aef5f71`; web **W1** + **W2** lên `main` |
| 26/07/2026 | Xong **KH5** (di chuyển + actor) gộp `e1ddc44`. Tạo file `BAT-DAU-TU-DAY.md` này làm điểm vào duy nhất cho người/agent mới |
| 26/07/2026 | Dựng cơ chế **nhiều agent song song**: sổ `NHAT-KY-AGENT.md`, lệnh `/viec` (một lệnh tự nhận biết nhận việc / ghi sổ), `AGENTS.md` cho agent ngoài, `MO-AGENT.bat`, hook `SessionStart` tự nạp bối cảnh (xem §10b) |
| 26/07/2026 | **KH6 xong phần thiết kế**: nghiên cứu giao thức chiến đấu (gói đánh quái/đánh người/chiêu/chết/hồi sinh, AI quái, nhịp 100 ms) → đặc tả có bảng byte → kế hoạch thi công 15 task. **Chưa viết code.** Bảng byte đã đối chiếu **cả hai phía** (Java + client Unity). Phát hiện mới: gói `-60` có **hai dạng độ dài** (§6.6); client lộ ra **3 lỗi thật** của bản Java ở `-8`/`-17`/`-9` (§6.16); mã `87` bị dùng cho hai việc (§6.17); sai trường **không** làm hỏng cả dòng byte, chỉ hỏng im lặng (§6.18); cổng PostgreSQL **không cố định giữa các lần khởi động** (§9) |
| 27/07/2026 | **KH6 thiết kế gộp vào `main`** (`c812617`, sau soi lần hai: gói 87 đối chiếu client khớp, test 7/7 xanh). **Thi công KH6 trên nhánh `kh6-thi-cong`**: Task 1 (catalog 20 chiêu) → 2 (chọn chiêu qua actor) → 3 (chốt kiểm) → 4 (sát thương ngũ hành) → 5 (quái trong actor + `-24`, vá kỷ luật actor bằng snapshot qua hộp thư) → 6 (gói 54/-9/-12/-13) → 10 (gói `-60` hai dạng) → 7 (điều phối đánh quái + 2 test lưới an toàn) — **mỗi task đều qua review đối kháng + kiểm thử đột biến**. Chủ dự án **chốt hệ số khắc chế ×1,25/×1,00/×0,80**; thêm **KH6b (Boss)** và **KH6c (stress 20 client)** vào lộ trình |
| 27/07/2026 | **KH6 tiếp tục**: Task 8 (quái hồi sinh 5s + hồi máu 30s trong nhịp actor, `e41678b`) → Task 9 (AI quái 2 tầng ưu tiên + gói `-11`/`-10`, `b22152c`) → Task 13 (chết & hồi sinh người chơi, gói `-17`/`-8`/`-16`/`-30/15`, `67e6650`). Task 13 chỉ cần thêm **một** trường `lucHoiSinh` chứ không phải Hp/HpToiDa như rủi ro đã lường (Task 9 đã dùng `nv.HpHienTai`) |
| 27/07/2026 | **KH6 XONG — gộp `main` `162bbe9`.** Task 11 (PvP `-60` + gói `-30/9` + **vá gap bất tử quái** của Task 13) → Task 14 (**test đầu-cuối 2 người qua mạng + CSDL thật**, diễn trọn chuỗi chiến đấu). **Soi toàn nhánh bằng opus** → bắt được `goiHoiMauNguoiChoi` là mã chết (chiêu type=2 gây sát thương sai) → đợt sửa cuối chốt hoãn có kỷ luật. **2 canary** (thả bug thật): reviewer bắt được lỗi byte im lặng, test E2E bắt được lỗi tích hợp — pipeline review đáng tin. Phân tích **bổ sung lộ trình** (9 khối thiếu, [`BO-SUNG-LO-TRINH.md`](docs/superpowers/plans/BO-SUNG-LO-TRINH.md)) + **mapping ảnh** đã lưu. Còn **Task 15 nghiệm thu tay bằng Unity** (chủ dự án) |
| 27/07/2026 | **KH11 XONG — gộp `main` `6ea467f`.** Máy chủ trước đó **không lưu gì** sau lúc tạo nhân vật; nay lưu 3 thời điểm (định kỳ 60 giây · khi rời · khi tắt máy) theo kiến trúc *actor chụp ảnh → hàng đợi gộp → thợ ghi riêng*, **actor không bao giờ chờ ổ đĩa**. Kèm **một tài khoản một phiên** (đá phiên cũ → chờ lưu xong → mới đọc; bản Java làm ngược và chưa bao giờ vá). 8 task, mỗi task TDD + đột biến + review đối kháng. **Soi toàn nhánh (opus) tìm ra 1 lỗi Nghiêm trọng thật**: client đăng nhập chồng trên cùng kết nối → mục rác trong sổ phiên → đóng channel hai lần → **panic chết cả tiến trình** (DoS từ xa) — đã vá, có test dựng lại trọn kịch bản 4 bước. Cũng vá bẫy §7.1 ở `XepVaCho`. Bốn khoản nợ mới ghi ở §8 |
| 28/07/2026 | **Nghiệm thu tay Unity — client VÀO ĐƯỢC GAME.** Chủ dự án bấm Play, tự động login `naruto/1`, vào bản đồ. Nhưng vấp 4 chỗ, mỗi chỗ đều là lỗi thật cần nhớ: (1) client Unity chỉ có file `mymap/{0,39,40,41}.bytes`, server gửi mapID=1-5 → **kẹt vô hình** — vá tạm bằng `cp 0.bytes {1..5}.bytes`; (2) client không gửi `-74` tự động khi mở, nhưng `ServerListScreen.switchToMe` chỉ set `loadScreen=true` khi RMS `ResVersion > 0` — thiếu, client kẹt màn "Xin chờ"; (3) client cần user click 2 lần (chọn server + đăng nhập), auto-login mặc định tắt → **sửa `SplashScr.cs`** thêm `waitToLogin = true` sau `GameCanvas.connect()` để bỏ qua cả 2 click; (4) RMS cache Java cũ (`NRdataVersion=20`) khớp `vsData=20` server → client tự đọc cache Java thay vì xin mới → parse fail âm thầm — xoá 16 file `NR*` cache. Thêm handler stub `-74` ở server (`login.go`), rebuild. **Server 100% đúng — testfull đi trọn từ bắt tay → clientOk → tileset+ME+MAP**. Còn 3 vấn đề mỹ thuật vẫn chưa map: sprite nhân vật (`part.bin` = Dragon Ball placeholder), tile map (1.bytes rỗng), atlas ảnh (`Big*.png` = Dragon Ball) — xem §8 mới thêm |
| 29/07/2026 | **Vertical slice Task 6A review nguồn ảnh:** helper QA Pillow đã qua TDD và hai mutation; review có SHA source/evidence, contact sheet sáng/tối và preview frame. **20/20 FAIL**: 19 nền magenta chuyển sắc vượt ngưỡng soft-matte tĩnh, HUD nền gần trắng không nhận diện chroma; Task 6B bị chặn đến khi có nguồn PASS. |
| 29/07/2026 | **Task 6A fix round 1:** thêm crop/cell edge thật, kiểm loop/timing, GIF temporal, gate executable và SHA đủ bốn evidence/asset. Raw source vẫn 0/20 PASS; output thích nghi 18/20 PASS, còn `hu_anthan` và HUD FAIL nên Task 6B vẫn bị chặn trung thực. |
| VS-6A-fix | Gate exact 20/source SHA/evidence fail-closed, precondition bind; derived output 20/20 PASS | `kh-vertical-slice` / pending commit | Raw 0/20; gate exit 0; report `docs/superpowers/reports/2026-07-29-task-6a-fix-report.md` |
| 29/07/2026 | **Task 6A fix round 2 hoàn tất:** TDD RED **19 test (9 fail + 2 error)** → GREEN **19/19**; review thật **20 PNG**, raw **0/20 PASS**, derived **20/20 PASS**; gate và bind exit **0**. Mutation SHA/raw summary/reviewer/evidence kind + alpha/crop/contact/animation/null asset đều exit **2**. Không sửa PNG nguồn hay worktree khác. |
| 29/07/2026 | **Vertical slice Task 6B hoàn tất, review sạch.** Closure canonical v11 khóa exact 20 nguồn review, 42 sprite `27041..27082`, part/animation/Mob/tile/bg/HUD/UI/runtime snapshot. Reviewer độc lập ký HUD golden 25/50/100 và composite Iruka sáng/tối; validator pin path/SHA/digest/identity, mutation 19/20 và alias tự ký đều đỏ. Full **109/109**, build A/B **669/669 file, DIFF_COUNT=0**. Commit `bf91ee8` → `545b9c6` → `0fd6fc5`; việc tiếp theo là Task 7 build/deploy staging theo lock. |
| 30/07/2026 | **Vertical slice Task 7 hoàn tất, review độc lập sạch.** Build/deploy artifact Naruto lock 21; fix hai guard fail-open (previous lock, protected Unity root) và canonical resource exact x1–x4 sau khi reviewer bắt được decoy `Small27042`. Mutation đều đỏ đúng mã; full Python **123/123**, Go toàn repo + `go vet` xanh; Unity staging preflight thật có `RESOURCE-PREFLIGHT-OK`. Commit `2a77a43` → `9627680` → `56fa302`. Việc tiếp: Task 8 catalog/NPC/Sói/collision, chưa bắt đầu. |
| 30/07/2026 | **Vertical slice Task 8 — vá Important collision sau audit mốc `902b5a8`.** `DiDuoc` ngoài grid 5×17 từng fail-open, khiến bước `(100,384)→(121,384)` thoát artifact. Đã viết RED/mutation, đổi map 1 sang footprint 120×408 authoritative, ghim map 6 là template chưa playable và sửa fixture cũ về vùng thật. Fresh Python **124/124**, Go game -count=1, focused collision/net -count=5, vet xanh. **Chờ independent re-review; Task 9 chưa bắt đầu.** |
| 31/07/2026 | **Vertical slice Task 9 hoàn tất sau khi Task 8 được independent audit APPROVE.** TDD bao phủ token sai/thiếu, command ID thiếu/trùng, sequence lùi/đứt, main-thread dispatch, screenshot SHA, graceful shutdown/persistence và fixture prefix/cleanup. Hai vòng review đã vá evidence rỗng, fixture hash, common shutdown path, reconnect và stale-session generation. Fresh Unity batch thật **50 assertion**, Go focused `-count=5` + toàn repo + vet, PowerShell parser/self-test và control smoke DB thật đều xanh; fixture còn **0 account** sau cleanup. Full gameplay hai client chưa chạy vì thuộc Task 10. |
| 31/07/2026 | **Task 9 được independent audit APPROVE; Task 10 hoàn tất phần không phụ thuộc Unity lock.** Viết 10 fixture/oracle TDD cho damage `5,2,2,2,2`, HP 137, clean/warm executable+lock+scope, DB/evidence SHA và NDJSON sequence. Ba mutation bỏ damage thừa/HP/lock đều làm suite đỏ đúng case, hoàn nguyên xanh 10/10. Runner thật gọi chính parser đã test và không truyền DSN qua command line. Clean full/warm Unity **BLOCKED** vì cả hai asset lock của người dùng vẫn tồn tại; không xóa lock, không chạy production account, không bắt đầu Task 11/12. |
| 31/07/2026 | **Task 10 lock diagnosis read-only:** hai lock file vẫn tồn tại nhưng **stale/không có owner**: không có Unity/Hub/AssetImportWorker, Restart Manager `0` owner, exclusive read-only open thành công; Editor.log cho thấy đúng `D:\NRO_Project` đã shutdown sạch. Full E2E vẫn chưa chạy vì preflight hiện chỉ `Test-Path`. Không xóa lock, không sửa Unity, không bắt đầu Task 11. |
| 31/07/2026 | **Task 10 canonical hardening theo audit B, không cherry-pick B:** port thủ công guard vào schema A và runner boundary thật. Suite 24 case tách server/asset/applied lock, zero-byte/SHA/array rỗng, Sói source 1, HP 137/150, clean+warm cùng PASS và sequence per-client; 4 mutation đều đỏ đúng. Runner dùng Restart Manager read-only; preflight thật cho hai stale lock đi qua, ghi `scope=task10_full`, rồi dừng đúng `E2E-ASSET-006` vì ArtRoot/embedded lock lệch. Go full + vet xanh; gameplay chưa chạy. |
| 31/07/2026 | **Task 10 — chẩn đoán asset provenance `E2E-ASSET-006`:** Task 7 A/B/release cùng lock `e5b093...` version 21; server HEAD đã chuyển có chủ ý ở Task 8 sang lock `54a638...`, pair `afd03...`, runtime/map 22. Các ArtRoot tạm Task 8 `.task8-build-2..4` đã bị dọn, không còn release/full client pair version 22; runner nhận build A từ tham số caller chứ không tự chọn. Dừng BLOCKED, không ghép tay, không build lại, không chạy gameplay/Task 11. |
| 01/08/2026 | **Rebuild canonical ArtRoot Task 8 v22 theo phê duyệt:** preflight 684/684 input, hai Unity lock owner=0; build A/B đều 669 file và byte-identical; lock `54a638...`, pair `afd03...`, version `22/21/22/21`, 668 output đúng record. Tạo fresh snapshot bằng overlay chính thức, publish immutable release theo lock SHA và Unity batch có đúng một `RESOURCE-PREFLIGHT-OK`. Giữ nguyên A/B/release; không sửa `D:\NRO_Project`, không gameplay, không Task 11. |
| 02/08/2026 | **VERTICAL SLICE (Task 0–9) GỘP `main`** — merge `47b8f77` (`--no-ff`). Sau khi review + tự kiểm chứng (nhánh chỉ chồng thêm lên `main`, gộp 0 xung đột; `go build`+`vet` sạch; `go test ./...` xanh **7 gói** — game 83s; `pytest assetpack` **124/124**; `assets/part.bin` Naruto thật 87.846B nhúng + khóa hash). Đưa về `main`: **code server NPC Iruka + va chạm Làng Lá + `assets/` nhúng** và **đường ống `assetpack`** — thay hoàn toàn POC sửa tay `part.bin` (đã bỏ, stash). Chặn `e2e-runs/`+`.artifacts/`+`artqa-canonical-run/`+`server/server-log.txt` vào `.gitignore`. **Còn dở: Task 10–12 chơi thật 2 client Unity + nghiệm thu** (đợt riêng); 11 worktree `.worktrees/` chưa dọn |
| 03/08/2026 | **Lưới chắn panic — mảnh A "vá nợ hạ tầng" — GỘP `main` `3c0a539`** (agent thi công, phiên điều phối soi đối kháng: đọc code 6 chỗ + tự chạy test + đột biến độc lập rồi mới gộp). Trước đó máy chủ có 6 điểm sinh goroutine mà **không điểm nào `recover()`** → một panic bất kỳ sập cả tiến trình. Thi công 6 task TDD nghiêm ngặt (test đỏ → code → xanh → **kiểm thử đột biến gỡ lưới thấy đỏ vì panic lan ra** → commit): **Task 1** hàm lưới chung `internal/antoan.ChayCoLuoi` (recover + `slog.Error` kèm `stack`, không nuốt im lặng); **Task 2** actor khu vực bọc từng sự kiện → chạy tiếp sau panic (dùng điểm chèn `suKienChay`); **Task 3** van an toàn `soPanicLienTiep`/`nguongPanicDungKhuVuc=100` — panic đại dẳng dừng riêng khu vực, **reset khi có nhịp trót** nên lỗi thoáng qua không chạm ngưỡng; **Task 4** `PhucVu`+`vongWriter` mỗi kết nối tự đóng; **Task 5** `luudinhky.vongLap`+`thogi.vongLap`+`go func(nv)` bỏ đúng bản lỗi, giữ vòng sống. 3 chính sách phục hồi theo loại luồng (spec §3.2 a/b/c). Nghiệm thu: `go build ./...` sạch, `go test ./... -count=5` xanh **8 gói** (game 415s), `go vet ./...` sạch, `gofmt -l internal/` không in gì. **Cố ý để ngoài phạm vi:** tách gửi mạng khỏi actor (§6.10b — lỗi KẸT, khác SẬP), 2 goroutine khởi động ở `cmd/nroserver/main.go` (LangNghe/E2E-control, ngoài internal/) |
| 04/08/2026 | **Dọn nợ nhiễu CRLF của `gofmt` (chip F9).** Trước: `gofmt -l` báo 9 file `.go` cũ vì checkout chính còn CRLF trên đĩa → `kiem-tra-tat-ca.ps1` phải **hardcode danh sách 9 file** để lọc, và `ci.yml` **không dám** chạy gofmt. Xác minh gốc: `.gitattributes` đã ép `* text=auto eol=lf` từ trước nên **blob toàn LF** và worktree mới ra LF sạch (`git ls-files --eol` xác nhận) — nhiễu chỉ còn ở working tree cũ trên đĩa. Vá: **bỏ danh sách hardcode** trong `kiem-tra-tat-ca.ps1` → gofmt soi mọi file; **thêm bước gofmt vào `ci.yml`** (job go, fail nếu `gofmt -l` có output); thêm ngoại lệ `*.bat`/`*.cmd = eol=crlf` (không đổi blob, không làm bẩn cây). **Kiểm thử đột biến**: thả 1 file `.go` sai format → `gofmt -l` bắt đỏ → xoá → sạch lại. `kiem-tra-tat-ca.ps1` full **XANH**: build/vet/**gofmt**/go test 8 gói (game 84.7s, store 1.2s DB thật)/pytest 118. Không đụng code Go nên `go test` **không đổi** | nhánh `claude/blissful-haibt-1a23fa`, chờ gộp |
| 22–23/08/2026 | **Map 1 Làng Lá + rút kinh nghiệm làm map** (nhánh `kh6b-boss-code`, `cb20bd8`→`ab865ce`). Chủ dự án nhận xét đúng: *"làm như chơi xếp hình, cứ chắp vá"*. **Gốc: làm ngược thứ tự** — dựng trước, đối chiếu concept sau, nên nhiều bản "sửa lỗi" lại **đẩy xa concept thêm** (giãn nhà cho khỏi chạm nhau ↔ concept muốn nhà san sát; fade chân núi ↔ concept **giấu** chân núi sau mái nhà; gò cao thành 2 mảnh rời ↔ concept có 1 nền đá liền). Ba thứ để lại: **(1)** `tools/map-lang-la/do-map.py` — thước đo 7 chỉ số so concept, **1 giây**, thay cho đợt soi agent 45–80 phút; kiểm chứng trên 3 mốc thật, bắt đúng lúc có đường kẻ (25,0) và lúc hết (5,3); 3 bẫy khi đo đã mắc và ghi trong chính file. **(2)** `docs/my-thuat/SO-TAY-LAM-MAP.md` — quy trình 7 bước + **quy tắc vàng cỡ ảnh: px ở x4 = đơn vị game × 4** (nền phải 2976 chứ không phải 1536, vì game hiển thị ở x2 — đã hỏng ảnh 2 lần vì quy tắc này). **(3)** `docs/my-thuat/KHUON-DE-BAI-ANH-MAP.md` + `docs/superpowers/plans/2026-08-23-lo-trinh-4-map-con-lai.md` — gộp 3 đợt ảnh của map 1 thành 1 đợt cho map 2–5. Sửa client Unity: `drawRegion(Small,…)` đang `throw NotImplementedException` nên **nhân vật vô hình**; ghép lại đầu+thân+chân (bỏ ảnh đầu là mất tóc + băng trán Konoha — agent soi vòng 3 bắt được); ô chân dung trên HUD. Server: gói **-42** gửi chỉ số nhân vật (thiếu nó thì thanh máu nằm im ở 0). **Quy luật đắt nhất, đã ghi vào sổ tay:** máy **sắp xếp** được, **hoà màu** được, nhưng **không tạo được nội dung mới** — thứ gì cần *mỗi ô một khác* thì bắt buộc phải vẽ; 3 lần thử ngược lại (bóng đổ, tăng tương phản đất, ghép mạch đá) đều hỏng, lần cuối còn tệ hơn lúc chưa làm, tốn ~1 giờ, đã hoàn nguyên 16 ảnh. Điểm agent soi qua 3 vòng: **45 → 52 → 56/100**, sau đó chuyển sang đo bằng thước thì 6/7 đạt |
| 24/08/2026 | **Map 1 Làng Lá + bộ công cụ làm map GỘP `main`** — merge `8c558e8` (`--no-ff`), 41 commit từ `kh6b-boss-code`. Trước khi gộp, pytest **bắt được một lỗi thật do chính nhánh này gây ra**: oracle độc lập `assetpack/tests/test_part.py` vẫn ghim `dx=dy=0` như thời ba mảnh nhân vật chưa ghép, còn `mapping-id.json` đã đổi sang đầu **−41** / thân **−50** / chân **−33**. Vá `dba8a9d`: ghim **cứng** `DY_DAU/DY_THAN/DY_CHAN` trong chính oracle — **không** đọc từ manifest, để sửa nhầm manifest vẫn bị bắt; kiểm thử đột biến hai chiều (lệch dy 1 đơn vị trong manifest → SHA khác · đổi hằng số trong oracle → test đỏ). Nghiệm thu **trên kết quả đã gộp**: `go test ./...` xanh 7 gói, `go vet`+`gofmt` sạch, pytest **118 passed** + 94 subtests. Đưa về `main`: map 1 Làng Lá (6/7 chỉ số), thước đo `do-map.py`, `SO-TAY-LAM-MAP.md`, `KHUON-DE-BAI-ANH-MAP.md`, lộ trình 4 map, các bản vá client Unity (nhân vật hết vô hình, ghép lại đầu+thân+chân, ô chân dung), gói **-42** chỉ số nhân vật, và nền móng KH6b Boss task 4+5. **Nợ mới phát hiện, CÓ SẴN trên `main`:** `assetpack/tests/test_task6b_canonical_preflight.py` hỏng ngay lúc thu thập (`SHARED_ROOT = ROOT.parents[1]` giả định repo nằm sâu 2 tầng — viết hồi còn chạy trong worktree) → **cả file chưa bao giờ chạy**, phải `--ignore` mới có 118 xanh; ghi ở §8. Boss KH6b task 8–15 làm tiếp trên nhánh mới cắt từ `main` |
| 24/08/2026 | **Dọn nhánh cũ: `kh-unity-e2e` + `kh-assets` ghi nhận đã gộp bằng `-s ours`.** Cả hai tách từ `62a0943` (29/07) và đã bị vertical slice thay thế khi gộp `main` ở `47b8f77`. Soi trước khi quyết: không file nào chỉ hai nhánh có mà `main` thiếu; mọi tên hàm/test của `kh-assets` đều có trên `main`; file bằng số dòng thì giống hệt từng byte; số còn lại `main` đều lớn hơn. **Gộp nội dung sẽ là bước lùi** — 15 xung đột add/add, và bản `mapping-id.json` cũ kéo ngược độ dời dọc ba mảnh nhân vật về 0, làm vỡ golden `part.bin` vừa vá ở `dba8a9d`. Cứu riêng bẫy `$PSScriptRoot` vào §9 (`25ce316`) — thứ duy nhất còn giá trị. Sau gộp: Go xanh 7 gói, pytest 118 passed. Còn 3 nhánh chưa soi: `kh-actor`, `kh-task10-oracle`, `claude/blissful-haibt-1a23fa` |
| 26/08/2026 | **Bật PostgreSQL lên thì lộ 15 test đỏ nằm im nhiều ngày.** Gốc của sự vô hình: `kiem-tra-tat-ca.ps1` chỉ đặt `NRO_TEST_DSN` khi thấy cổng 5432/5433 đang nghe — máy này thường tắt CSDL nên **mọi test tích hợp tự bỏ qua trong im lặng**, bảng tổng kết vẫn `[OK] go test`. Ghi vào §9. Ba lỗi tách bạch: **(1) 13 test do CHÍNH MÌNH làm đỏ** (`317b068`) — `ca44983` thêm gói **-42** (chỉ số nhân vật) gửi ngay sau MAP INFO, nhưng ba helper bắt tay trong test chỉ đọc `-82/-30/-24` rồi dừng; gói -42 nằm lại làm **nghẽn ống pipe** → máy chủ **hết giờ ghi 30 giây** mỗi test, trông như treo chứ không ra lỗi logic (`login_test.go` đã vá đúng từ trước và còn ghi rõ triệu chứng — ba file tích hợp bị bỏ sót). Kèm `dongbo_tichhop_test.go` còn ghim `headAvatar.size = 0`, khuôn cũ trước khi ME_INFO mang cặp ảnh chân dung. **(2) Lỗi có sẵn — gói -9 vượt tầm báo sai byte** (`1338c1c`): test cho A đi tới `x=0` để ra ngoài tầm, nhưng từ khi map 1 có **tường vô hình ở hai cột biên** thì `x=0..36` bị chặn → `-7` bị từ chối → A **đứng yên** ở x=100, cách quái 60 < tầm 100 → trúng đòn 2 sát thương. Triệu chứng hiện ra là *sai byte*, không giống nguyên nhân thật. Đổi sang `x=48` (ô đi được đầu tiên, |48−160| = 112 > 100). **(3) Lỗi có sẵn — duplicate key**: `demHauToTichHop` bắt đầu từ **0 mỗi lần chạy** nên sinh đúng một dãy tên, mà chú thích lại **hứa sai** là *giữa các lần chạy lại không đụng nhau*; cộng với `t.Cleanup` xoá hàng bằng `_, _ =` nên **nuốt lỗi** — PostgreSQL tắt giữa chừng là để lại một hàng, hàng đó làm mọi lần chạy sau báo duplicate key. Nay bắt đầu từ mốc ngẫu nhiên (trần 9000 vì tên nhân vật tối đa 10 ký tự, tiền tố dài nhất 6) và `t.Errorf` khi dọn dẹp hỏng. **Đo lường:** trước 14 đỏ / gói chạy 450s → sau **xanh cả 7 gói** / 88,9s, chạy lại lượt hai vẫn xanh, `nro_test` sạch 0/0 sau khi chạy |
| 26/08/2026 | **Dựng cầu nối Tiled — chữa đúng cái gốc của "chắp vá, xếp hình".** Chủ dự án nhận xét đúng, và gốc không nằm ở tay nghề mà ở **cách làm**: bố cục map sống trong `ve-bo-cuc-map.py` dưới dạng **93 con số cứng / 167 dòng**, nên dịch một căn nhà = sửa một số = script **xoá sạch 1248 ô rồi vẽ lại từ đầu**. Đó là *sinh lại*, không phải *sửa*. Ba thứ thiếu nay đã có: bố cục nằm trong `.tmx` sửa được bằng chuột; `xem-truoc.py` vẽ map ra ảnh trong **0,6 giây** (trước phải mở Unity → Play → chờ nạp → chụp, mỗi vòng vài phút, nên mỗi lần sửa là một lần đoán mò); **ảnh concept nằm ngay trong `.tmx` làm lớp giấy can mờ 45%** — trước đây concept một nơi bố cục một nơi, không gì buộc chúng khớp nên chúng trôi dần, đó chính là lý do nhiều bản "sửa lỗi" lại đẩy XA concept. **Bẫy quy ước neo, đã tự mắc:** client neo vật thể to ở góc **TRÊN**-trái, Tiled neo góc **DƯỚI**-trái, lệch đúng `cao_ô − 1` hàng — bản `xem-truoc` đầu tiên vẽ sai quy ước rồi thấy *16 vật thể treo cách mặt đất 1–7 ô*, suýt báo map hỏng; đo lại đúng quy ước thì **16/16 chạm đất chính xác** (map luôn đúng, công cụ mới sai). **Ba phép kiểm chứng:** `map → .tmx → map` **khớp từng byte** (1250/1250) · vẽ từ `.tmx` và vẽ từ map game ra cùng ảnh trừ **10 ô** chồng lấn (0,065% pixel — game vẽ đèn đè lên nhà, Tiled vẽ ngược) · đổi mà chồng lấn/tràn map thì **báo lỗi và dừng**. Hạn chế ghi rõ trong README, không giấu. Tiled **chưa cài** trên máy này |
| 27/08/2026 | **Trả 2 nợ trong đường ống tài sản (`assetpack`), cả hai đều là "im lặng không ai biết".** (1) `test_task6b_canonical_preflight.py` hỏng ngay lúc *thu thập* nên **chưa bao giờ chạy**: `SHARED_ROOT = ROOT.parents[1]` đếm bậc thư mục, đúng khi file được viết trong worktree `<gốc>/.worktrees/<tên>`, `IndexError` ở checkout thẳng. Nay hỏi `git rev-parse --git-common-dir` — **phải hỏi git chứ không dùng thẳng `ROOT`** vì `anh-nguon/`+`anh-game/` bị gitignore nên worktree không có ảnh, chỉ checkout chính mới có. Chạy được rồi thì **cả 6 chốt kiểm đều XANH** — không cái nào ghim thứ đã trôi. (2) `mapping-id.json` còn ghi map **5×17 (85 ô)** trong khi map 1 Làng Lá thật là **52×24 (1248 ô)** → lần `assetpack build` kế tiếp sẽ đè lưới va chạm đúng bằng bản sai và tái sinh lỗi "đứng ngoài lưới bị từ chối âm thầm mọi di chuyển". Đồng bộ `build_inputs.map` từ hai nguồn sự thật (bố cục client `1.bytes` + artifact `map-1-collision.bin`), đặt `tile_art[].solid` + `tiles[]` theo 6 mã tile đặc `{39,60,61,67,68,84}` — đối chiếu 1248 ô thấy **không mã tile nào vừa đặc vừa đi được** nên `validate_collision_semantics` qua sạch. **Nghiệm thu bằng chạy thật `assetpack build`**: `map-1-collision.bin` / `part.bin` / `skill.bin` / `catalog.json` giống hệt bản đang chạy. Chốt kiểm mới `assetpack/tests/test_map1_manifest_khop_artifact.py` (4 test; đột biến: trả manifest về 5×17 → 2/4 đỏ ngay). pytest **128 passed** (118 → +6 preflight +4 mới), `go test ./...` 8/8 gói. **Nợ mới phát hiện, ghi §8:** `tile_set_info.bin` lệch y hệt kiểu (build 1696 byte gốc vs repo 1734 byte đã thêm tile Làng Lá) và `lang-la-client-assets/res/mymap/1.bytes` đã cũ 51 ô |
| 27/08/2026 | **Gỡ quả mìn cuối của bộ artifact: `assetpack build` không còn lùi bảng tile về bản gốc.** Bảng gốc Ngọc Rồng chỉ khai tile 1–39; map Làng Lá dùng 60–69/83–88 nên client đặt `types[]`=0 cho cả màn → rơi xuyên đất. Bản 1734 byte đang chạy là **sửa tay** bằng `tools/map-lang-la/them-tile-vao-tileset.py`, manifest không hề biết → build lại là tái sinh lỗi. **Cân ba hướng:** (a) thêm file .bin nguồn 1734 — **loại**, `docs/my-thuat/anh-nguon/` bị gitignore nên file sẽ không được commit (`git check-ignore -v` xác nhận), mà đặt chỗ khác thì lại là một blob nhị phân không truy được nguồn — đúng cái vừa gây ra chuyện này; (b) suy lại từ `tile_art[].solid` — **loại**, một cờ `solid` không tả nổi ba vai trò khác nhau (sàn / chặn ngang / chặn nhảy lên: tile 62 chặn ngang nhưng `solid=false`, tile 83 thì ngược lại); (c) **chọn** — đưa phép biến đổi vào đường ống, khai trong `mapping-id.json` (`build_inputs.tile_set_overlay`, git theo dõi được, đã tự kiểm), `assetpack/tileset.py` thi hành. Nguồn vẫn là MỘT: bảng gốc đã review giữ nguyên, `tile_set_baseline` + `production_closure_sha256` không đổi. **Chốt chặn cuối** `validate_floor_closure`: tile nào máy chủ coi là đặc mà client không khai là sàn thì build đỏ — đỏ **kể cả khi `tile_set_overlay` bị xoá hẳn**, nên mìn không dựng lại được. **Nghiệm thu chạy thật:** build ra rồi so từng byte — `part.bin` `skill.bin` `map-1-collision.bin` `catalog.json` và `tile_set_info.bin` **1734 byte** đều giống hệt bản đang chạy; mọi `outputs[]` sha trong asset-lock build ra khớp bản đã commit. **Đột biến (build thật, 4/4 bị giết):** bỏ tile 60 khỏi nhóm sàn → `ASSET-TILESET-013` lệch golden · bỏ tile 60 **rồi sửa luôn golden** để vòng qua chốt hash → vẫn đỏ ở chốt ngữ nghĩa · **xoá hẳn `tile_set_overlay`** → đỏ · khai lại tile đã có trong bảng gốc → `ASSET-TILESET-012`. pytest **137 passed** + 98 subtest (128 → +9), Go 8/8 gói (PostgreSQL 5432 bật) |
| 27/08/2026 | **Bốn việc code chạy song song bằng subagent — ba đã gộp `main`.** Chủ dự án giao: soi toàn bộ, việc gì cần code thì đẻ subagent làm trước, phần ảnh để sau. Chia việc theo **vùng file** chứ không theo "nghe có vẻ khác nhau" — `khuvuc.go` nặng 63KB nên hai agent cùng vào `internal/game/` là đụng chắc. **(1) Thước đo chỉ ra CHỖ hỏng** (`55b576b`): `do-map.py` xuất thêm ảnh chẩn đoán 8 ô, mỗi ô tô đúng vùng làm ra con số, kèm mốc concept + ngưỡng + phán quyết, và **nhắc lại chính bẫy đo** ngay trên ô đó; kèm cho `ve-bo-cuc-map.py` nghỉ hưu an toàn (chạy nhầm là dừng trước khi chạm file). **(2) Nợ assetpack** (`155c27c`): `test_task6b` **chưa bao giờ chạy** (`SHARED_ROOT` giả định repo nằm sâu 2 tầng — hoá ra nó phải trỏ vào gốc checkout CHÍNH vì ảnh nguồn bị `.gitignore` loại nên worktree không có); manifest ghi map **5×17** trong khi thật **52×24**, tức cả khối 85 ô thay vì 1248. Sau vá: `assetpack build` sinh ra artifact **giống hệt từng byte** bản đang chạy. pytest 118 → **128**. **(3) Gỡ mìn `tile_set_info.bin`** (`53c84c2`): build ra **1696** byte (bảng gốc Ngọc Rồng chỉ tile 1–39) trong khi bản chạy **1734** (đã thêm tile Làng Lá) → build rồi chép = **nhân vật rơi xuyên đất**, đúng lỗi đã vá 22/08. Đưa phép biến đổi vào chính đường ống thay vì thêm file `.bin` sửa tay — vì `anh-nguon/` bị `.gitignore` loại nên file mới đặt ở đó **sẽ không bao giờ được commit**, chính là cái bẫy đã gây ra chuyện này. pytest **137**. Tự chạy 2 đột biến: xoá hẳn khai báo (đúng tình trạng mìn cũ) và bỏ một tile khỏi nhóm sàn — **cả hai bị bắt**, lỗi báo đúng hậu quả *"nhân vật sẽ rơi xuyên đất"*. **(4) Thu hồi phiên JWT** (`9492315`): xem §7. **Nợ mới phát hiện:** `versions.data` manifest ghi 21 còn asset-lock + hằng Go ghi 26 — nhưng lỗi này **kêu** (có test so hằng Go với lock) nên để lại. **Còn chờ:** nhánh waypoint + bg-item đang được soi đối kháng, chưa gộp |
| 27/08/2026 | **Cổng rời Làng Lá + 8 bóng đổ đã lắp được.** Gói MAP INFO đang gửi **số 0** cho hai chỗ mà client Unity đã hỗ trợ sẵn. Hậu quả thật: (1) **không có đường nào rời map 1** — ai vào Làng Lá là ở đó mãi; (2) **8 ảnh bóng đổ vẽ xong từ 24/08 nằm không dùng được**. Đọc cả hai phía (Java `Zone.mapInfo` / `DataGame.sendDataItemBG` ↔ Unity `Controller.loadInfoMap` / `case -31` / `BgItem.cs`) thì lộ ra **việc lớn hơn đề bài**: waypoint chỉ là cái nút, muốn nó chạy phải có **gói `-23` đổi map**; và bóng đổ cần **hai** gói — `-31` bảng mẫu (một lần mỗi phiên) + ô bg-item trong `-24`. Thiếu `-31` thì client **bỏ qua im lặng**; thiếu `-23` thì bấm cổng là nhân vật **đứng hình vĩnh viễn** (client tự khoá rồi đợi máy chủ) — tức là làm nửa vời còn **tệ hơn không làm**. Đã làm trọn: cổng **torii** giữa map 1 (`isEnter=true` nên **bấm mới đi**, đi ngang không bị hút) → Làng Cát, và Làng Cát có cổng **về lại** để không nhốt người chơi; đổi map dùng đúng `RoiKhuVuc` đồng bộ (chụp trạng thái → rời khu vực cũ → *mới* được ghi vào nhân vật → `-22` + `-24` → client `-39` vào khu vực mới). **11 bóng đổ** đặt bằng cách đo trên chính ảnh `xem-truoc.py --game`, kiểm lại bằng nguồn độc lập là `map-1-collision.bin` (ô đặt phải đi được, ô ngay dưới phải đặc). Ảnh cắm đủ **x1–x4** vào `Resources/res/x*/mapBackGround/500..507`. **Kiểm thử đột biến 15/15 bị giết** (con số này là **tự đo, nói quá** — một vòng soi độc lập sau đó tự nghĩ 22 đột biến và chỉ giết được **18**, xem dòng dưới) — một đột biến lộ ra dòng `*b.nhanVat = hienTai` vừa vô dụng vừa **có hại** (đẩy ngược mọi thay đổi xảy ra giữa lúc chụp và lúc rời khu vực), đã bỏ. Ba bài học mới ghi ở §6.32–34. `go test ./...` xanh 7 gói, `vet`+`gofmt` sạch | nhánh `worktree-agent-a6c5f42987d2a805f` |
| 27/08/2026 | **Vá 6 việc do một vòng soi đối kháng độc lập tìm ra** (cùng nhánh cổng + bóng đổ). Nặng nhất: `xuLyDoiMap` **tuyên bố** hai đường ra đều gửi gói, thực tế có **bốn** và **hai im lặng** (`TimBanDo` trả nil · không chụp được trạng thái). Client tự khoá mình ngay khi bấm cổng rồi đợi máy chủ, nên im lặng = **đứng hình vĩnh viễn, phải thoát game** — và đường thứ hai xảy ra thật khi van an toàn dừng hẳn actor (100 nhịp panic liên tiếp), lúc đó **cả khu vực** kẹt. Nay cả hai trả gói `46` (mở khoá tại chỗ); bước chụp ảnh dời lên **trước mọi chốt** để mỗi đường ra đều có sẵn toạ độ authoritative. Thứ hai: `-22`/`-24`/`46` đổi từ `s.Gui` (ghi thẳng socket) sang `s.XepGui` (hàng FIFO) — `Gui` không vỡ khung nhưng **nhảy hàng** trước gói `-5` actor vừa xếp, đẻ ra **nhân vật ma của làng cũ đứng giữa làng mới**; sau khi đổi thì `doimap.go` không còn chỗ nào ghi thẳng socket sau lúc vào thế giới. Thứ ba: bóng **quán trái** rộng 8 ô neo ở cột 10 nên hai ô cuối **vẽ đè lên vách đá dựng đứng** (dời về cột 8; test nay kiểm **cả dải ô** chứ không chỉ ô neo). Thêm **6 test** bịt các lỗ hổng bản soi chỉ ra: lưu ngay khi đổi map, vị trí xét cổng lấy từ actor, sao sâu `DiemDi`, chốt map đích tồn tại — **đột biến 21/22 bị giết** (cái còn lại là data race thuần, §8). Sửa hai chỗ sai trong chính sổ này (tên test `TestVersionKhopAssetLock` là **bịa**) và 6 số dòng lệch trong chú thích tham chiếu client — trong đó **2 cái bản soi báo sai**, đã đối chiếu tay với `D:\NRO_Project` | nhánh `worktree-agent-a6c5f42987d2a805f` |
| 28/08/2026 | **Vá ba lỗ hổng bảo mật đăng nhập — nợ chặn ngày mở cửa** (`3aea99d`). Đây là việc thứ 6 trong đợt giao subagent. **(1) Chống dò mật khẩu:** trước đó **không có gì cả** (tự grep xác minh, không tin sổ). Gói mới `internal/chanip`, túi thẻ đếm theo **IP** — *không* theo tên tài khoản, vì đếm theo tên thì kẻ xấu cố tình gõ sai là **khoá được tài khoản người khác**. Chi tiết tinh ý: vẫn trừ thẻ **kể cả khi tên không tồn tại**, vì nếu không thì *số thẻ còn lại* thành chuông báo "tên này có thật" — tức mở lại lỗ hổng 2 bằng cửa sau. **(2) Rò rỉ thời gian:** nhánh "không có tài khoản" trả về **ngay** còn "sai mật khẩu" phải qua hàm băm chậm → **liệt kê được tài khoản nào có thật**. Vá bằng băm giả sinh sẵn. **(3)** Giới hạn 32 kết nối/IP, thiếu gác thì **từ chối chạy** thay vì âm thầm mở cổng. **Tự chạy lại cả hai phép đo, không tin báo cáo:** 200 lần sai → chỉ **10 lần** tới được CSDL, dò đủ 200 lần mất **1h03m** thay vì 2,1 giây · rò rỉ thời gian **nhiễu đo 478 µs vs chênh lệch thật 106 µs** (phép đo tự hiệu chuẩn nhiễu bằng hai nhánh giống hệt nhau — đúng phương pháp), trước là ~13,5 ms |
| 28/08/2026 | **KH6b Task 6 — hoá ra đã xong từ trước, nhưng có lỗi thật** (`0a82053`). Bộ phát thanh `BoPhatThanh` không nằm ở `loathongbao.go` như kế hoạch ghi mà ở **`thongbao.go`**, có từ `6e51dd8` — lần soát của tôi tìm theo **tên file** nên kết luận nhầm là chưa làm. Bài học: soát trạng thái phải grep theo **thứ code làm được**, không theo tên file kế hoạch đoán trước. Nhờ chuyển sang nghiệm thu bằng **kiểm thử đột biến** mà lỗi lòi ra: **loa đảo thứ tự thông báo**. `nhip()` lấy yêu cầu ra khỏi hàng đợi **trước** rồi mới xem đồng hồ; chưa đủ 1,5 giây thì xếp lại **cuối hàng** — mà vòng lặp đập 50ms/nhịp nên giữa hai lần phát có tới **~29 nhịp xoay hàng đợi**. Bốn thông báo A,B,C,D lên sóng thành **A,D,C,B**: người chơi đọc *"Boss đã bị hạ"* **trước** *"Boss xuất hiện"*. **Tự tái hiện được** — lấy nguyên bản `thongbao.go` trên `main` rồi chạy test mới lên nó, ra đúng `[A D C B]`. Vá bằng **kiểm đồng hồ trước, lấy ra sau**: không động vào hàng đợi thì không có gì để xoay. Kèm một lỗ cùng dạng: `BatDau()` gọi hai lần sinh hai vòng lặp, cả hai `close(b.xong)` → **panic giết cả tiến trình lúc tắt máy chủ**. Ba test mới, đột biến **5/5 bị giết**. Còn lại KH6b: **12 → 13 → 14 → 15** |
| 28/08/2026 | **KH6b — gộp 4 task cùng lúc, `922eed8`** (task 12 chia thưởng · 13 lịch sinh boss · 14 nối bộ phát thanh · 15 hai boss mẫu). Bốn agent chạy song song, chia theo **vùng file**; task 12 và 13 cùng sửa `khuvuc.go` nhưng khác vùng cách nhau ~600 dòng nên git tự hoà, **không một xung đột nào**. **Phát hiện lớn nhất: task 8–11 THỰC TẾ CHƯA LÀM** — hai agent độc lập cùng báo, tôi tự xác minh bằng grep *trong* `khuvuc.go`: `BossChao`/`BossTroi`/`BossRaDi` **0 lần**, `MocNoiGian` **0 lần**, `CongSatThuong`/`ChiaThuongTheoDongGop` **0 lần**. Những tên đó chỉ nằm ở `boss.go` — sản phẩm của task 3. Đây là lần soát nhầm **thứ hai** trong ngày, cùng một kiểu (lần 1 grep theo *tên file*, lần 2 grep theo *tên kiểu* trên toàn `*.go`). Cách đúng đã ghi vào kế hoạch: **grep trong đúng file mà task đó phải sửa**. **Việc từng task:** task 12 nối `CongSatThuong`/`ChiaThuongTheoDongGop` vào actor, thưởng = `KinhNghiem` (bền nhờ đường lưu sẵn có, và đã đi xuống client qua `cTiemNang`), chia ở **mỗi pha chết** rồi xoá sổ **ngay trong cùng hàm** để không hở cho chia hai lần; 13 đột biến — **3 cái sống sót ở lượt đầu**, phải viết thêm test mới bịt. task 13 dựng `lichboss.go` theo khuôn `BoLuuDinhKy`, sửa `khuvuc.go` đúng 3 chỗ; **test chốt CHAN-3** khẳng định `MoSoKhuVuc()` không tự sinh boss nào. task 14 phát hiện `main.go` **chưa bao giờ gọi `GanLoa`** nên mọi thông báo boss **biến mất không một dòng log**; đặt `loa.Dung()` **trước** `TatMayChu` vì `LuuTatCa` chỉ chờ actor 100ms mà actor phải xả hết hộp thư gameplay trước khi trả lời — bão thông báo có thể làm **mất tiến độ cả một khu vực**. task 15 làm 2 boss mẫu; agent bị dừng **trước lúc commit** nên tôi commit thay và **tự chạy 4 đột biến, cả bốn bị giết**. Nghiệm thu trên `main`: Go xanh 8 gói, vet+gofmt+build sạch, `nro_test` sạch 0 dòng. **Hai mối nối cố ý treo:** `bossNhanSatThuong` chưa ai gọi (chờ task 8), `TheSongBoss.Tat()` chưa ai gọi (chờ task 9), lịch boss chưa bật ở `main.go` |
| 28/08/2026 | **KH6b Task 8 — người chơi đánh được boss** (`15be8df`, gộp `main`). Client **không hề có khái niệm "đánh boss"**: boss vào thế giới bằng chính gói `-5` của người chơi nên với client nó là một *char*; nút đánh gửi đúng gói `-60`, chỉ khác con số id. **Dấu ÂM của id boss là thứ duy nhất tách hai đường** — id nhân vật là `BIGSERIAL` nên luôn dương. Agent đọc client **tận nơi** thay vì suy: `-60` tra mục tiêu bằng `findCharInMap` (`Controller.cs:3989`) — hàm so int thường, không giả định dấu — nên id âm chạy đúng; còn `-9`/`-12` mang `byte mobId`, **về khuôn byte không nhét id boss vào được**. Sửa gọn **+139 dòng, 0 xoá** trong `case suKienDanhNguoiChoi`, dùng lại nguyên `KiemDanhHopLe` (task 5) + `TinhSatThuong` + `bossNhanSatThuong` (task 12). **Hai quyết định đáng nhớ:** con số báo xuống client là **hiệu máu trước/sau** chứ không phải con số vừa tính — chỉ còn **một** nơi kẹp máu, và nếu đòn bị từ chối thì gói là cú *miss* thật thà thay vì báo con số máy chủ không hề trừ; và **không gửi kèm `-30/9`** như nhánh PvP, vì lý do khiến PvP cần nó (gói `-5` người chơi ghi `hp = HpToiDa`) không tồn tại ở boss. Trạng thái không đánh được **không cần nhánh mới**: lọc `DangHien()` loại NGHI/RA_DI — chốt này còn che cho **client**, vì đoạn `char3.cx` (`Controller.cs:3934`) nằm **ngoài** try/catch nên bắn `-60` vào con boss client chưa nhận `-5` là gửi null-deref cho cả khu vực. **Tôi tự chạy đột biến quan trọng nhất**: thay `bossNhanSatThuong` bằng `boss.Hp -= satThuong` → **4 test đỏ**, gồm cả hai test chia thưởng — đúng cái bẫy task 12 cảnh báo. Còn **1 đột biến sống sót**, agent nói thật: ép `chiMang=false` khi trượt là no-op vì KH6 chưa có công thức chí mạng. **Hôm nay boss mới chỉ ĐỨNG YÊN CHỊU ĐÒN** — chết xong nằm mãi ở TROI (task 9), không đánh trả (task 10), không biến hình (task 11) |
| 28/08/2026 | **KH6b Task 9 — máy trạng thái vòng đời boss** (`9e6f6c9`, gộp `main`). Trước đó boss chết xong **nằm mãi ở `TROI`**, xác đứng giữa bản đồ vĩnh viễn, và mục lịch đẻ ra nó **chết theo** (không bao giờ gọi `Tat()` nên không sinh được con thứ hai). Nay chạy trọn `NGHI → HIEN → CHAO → HOAT_DONG → TROI → RA_DI → rời hẳn`, trong `nhipTim` 100ms ở **vòng lặp riêng** cạnh đàn quái. **Đã nối mối treo của task 13**: boss mang cờ sống của mục lịch qua trường mới `Boss.The`, máy trạng thái gạt nó lúc phát `-6` — nhờ vậy **lịch boss sinh được con thứ hai**. **Vài con số có nguồn, không bịa:** nhịp thoại 100ms/ký tự kẹp 2 giây lấy từ bản Java (`Boss.java:507-510`) nhưng **đếm ký tự chứ không đếm byte** — *"Đủ rồi."* là 7 ký tự mà 10 byte, đếm byte thì mọi câu tiếng Việt treo lâu hơn ~40% vì một chi tiết mã hoá; xác nằm tối thiểu **3 giây = hai nhịp loa 1,5s**, để dòng chữ *"vừa hạ gục Boss X"* kịp tới nơi khi người chơi **còn thấy xác** (bản Java không có con số này). **Tôi tự chạy đột biến**: bỏ `TheSongBoss.Tat()` → **4 test đỏ**, gồm test đầu-cuối chứng minh bộ điều phối *ngoài* actor và máy trạng thái *trong* actor nối được vào nhau. **Đột biến 12/12 bị giết, 0 sống sót.** **Còn treo cố ý:** bước `HOAT_DONG` để trống chờ task 10; nhánh biến hình ở `RA_DI` chưa làm vì **GATE task 1 vẫn chưa có** — hậu quả hôm nay là boss nhiều pha chết ở pha 1 sẽ **rời hẳn** thế giới thay vì lột xác |
| 27/08/2026 | **`web/.env.local.example` vào git — hết cảnh đọc README rồi đi tìm một file không tồn tại.** Ba chỗ (`web/lib/db.ts`, `web/lib/xacthucphien.ts`, `web/README.md` mục 2) bảo copy từ file mẫu đó, nhưng nó **chưa bao giờ được git theo dõi** — ai clone repo hoặc mở worktree mới đều gặp lỗi thiếu `DATABASE_URL`/`AUTH_SECRET` rồi tìm một file ma. **Thủ phạm không phải `web/.gitignore`** (mẫu `.env*.local` đòi chuỗi kết thúc bằng `.local` nên không khớp) **mà là `.gitignore` gốc dòng 42 `.env.*`**; ngoại lệ `!.env.example` ngay dưới chỉ khớp đúng cái tên đó. Gỡ bằng `!.env.local.example` đặt trong `web/.gitignore` — `.gitignore` thư mục con thắng thư mục cha nên **không phải sửa file gốc** (yêu cầu: chỉ đụng `web/`). **Cổng giữ nguyên `5433`, chỉ thêm chú thích dạy tự kiểm** (§3/§9): máy này đang nghe `5432`, nhưng sửa số cho khớp hôm nay đúng là việc agent 26/07 đã làm — chỉ dời quả mìn sang lần khởi động sau; chú thích kèm cả triệu chứng khi dùng sai cổng vì nó trông hệt lỗi code. **Kiểm chứng:** `git status` thấy file là untracked *trước khi* add (hết bị ignore), `web/.env.local` thật vẫn bị chặn, `@next/env` đọc file mẫu ra đúng hai biến (chú thích không làm hỏng parse). Bài học `git check-ignore -v` → §6.15c; nợ nhỏ còn lại → §8. Nhánh `claude/sad-rhodes-d515e4`, commit `ef57a61` |
| 02/09/2026 | **Chia agent theo model + GATE task 1 + task 14 xong.** Chủ dự án yêu cầu không phải agent nào cũng dùng model đắt. Chia: việc *đọc-và-trích-dẫn* (GATE task 1) → **Sonnet, agent chỉ đọc**, ~108k token — bằng 1/3 một agent Opus viết code; việc *nối dây theo mẫu có sẵn* (task 14 phần 2) → **Sonnet**; việc *logic actor lõi* (task 10, 11) → **Opus**. **GATE task 1** (`4dbc8db`): client **chấp nhận** id đã `-6` quay lại bằng `-5` — `-6` xoá hẳn khỏi `vCharInMap`, `-5` tạo `Char` mới hoàn toàn, không cache nào theo id. Bẫy duy nhất (`isUsePlane` làm `-6` không xoá) **không chạm boss**: cờ chỉ bật trong `case -65` (teleport) mà server Go không bao giờ gửi (grep 0 lần). **Task 11 hết bị chặn.** **Task 14 phần 2** (`0709caf`): `noiLichBoss()` bật lịch boss ngay sau bộ phát thanh; tắt `lichBoss.Dung()` **trước** `loa.Dung()`; thêm `MucLichTuNoiDungBoss` nối dữ liệu task 15 vào bộ điều phối task 13. **Hai agent Sonnet liên tiếp làm xong nhưng QUÊN COMMIT** (Explore không có quyền commit; agent task 14 thì quên) — tôi commit thay, đã ghi §8. Bớt việc thừa: không soi đối kháng từng task nữa (để một lượt cuối KH6b), test full một lần/nhánh, gộp doc bằng script. Dọn sạch mọi nhánh cũ đã gộp — kho chỉ còn `main` + `so-no-thuong` |
| 02/09/2026 | **KH6b Task 10 — AI boss** (`ec937f5`, Opus). Boss hết đứng yên: chọn mục tiêu **2 tầng** (ai đánh nó nhiều nhất theo sổ `SatThuongTheoNguoi` của task 12 → gần nhất; loại chết/bất tử/ngoài tầm; **hoà thì id nhỏ hơn** vì map Go duyệt ngẫu nhiên, không có chốt này boss quay qua quay lại), đánh bằng chiêu KH6 qua `KiemDanhHopLe` (`VoHanChakra=true`), phát **`-60` + `-30/9` đúng nhánh PvP** — không dùng `-9/-12` vì hai gói đó mang `byte mobId`, id boss âm không nhét vào được. Di chuyển **24 đơn vị/nhịp = đúng 1 ô lưới va chạm** qua chính `kiemDiChuyenHopLe` của người chơi; vướng tường thì **đứng im** (không trượt, không tìm đường — rung lắc trông như game hỏng). Agent **tự khai hai con số không có nguồn** (15s/câu thoại đánh; tầm đuổi = 2× tầm chiêu xa nhất) — ghi rõ trong chú thích là lựa chọn. **Bài học đột biến:** hai đột biến sống sót lần đầu vì test đặt **sai khoảng cách** — mục tiêu trong tầm chiêu thì nhánh di chuyển không bao giờ chạy; kiểm bất tử ở tầng chọn mục tiêu bị `KiemDanhHopLe` che. *Muốn test cắn một nhánh thì phải dựng đúng tình huống buộc nhánh đó chạy.* Agent **quên commit** (lần 3 liên tiếp) — tôi commit thay. `main` xanh 8 gói với 10 + 14 cùng lúc |
| 02/09/2026 | **KH6b Task 11 xong → KH6b BẢN 1 ĐỦ CODE** (`ea43237`, Opus). Biến hình: `case BossRaDi` còn pha → phát `-6` **trước** `-5` cùng lượt actor, **giữ id**, đặt lại sạch trạng thái pha cũ (máu `HpChon[0]` pha mới, nhịp, mốc, hai hệ số, cả ba con trỏ thoại, `LucHienPha`). Nổi giận: kiểm mốc ngay sau trừ máu trong `bossNhanSatThuong`, **chỉ khi còn sống** → đòn kết liễu không phát chat ma; `HeSoDame/HeSoDef` ghi đè tuyệt đối, `HeSoNhip` nhân dồn **một lần mỗi mốc**; 0 = chưa nổi giận = ×1. Hai quyết định có lý do: HP pha mới lấy `HpChon[0]` không random (CHAN-7 cấm actor chọn ngẫu nhiên — giá: các giá trị sau của `HpChon` pha 2+ không dùng); thông báo 93 phát lại lúc HIỆN pha mới theo `LaBossLon` của pha đó. **Agent tự bắt 2 test xanh giả** ở vòng đột biến lần 1 (so id sau biến hình với chính nó; boss mẫu không có `BoChieu` nên AI thoát ngay dòng đầu) — sửa xong 7/7 đột biến bị giết. **Lần đầu agent tự commit** nhờ brief có bước cuối bắt buộc. `main` xanh 8 gói, CSDL thử sạch. Kế: một lượt soi đối kháng toàn KH6b (Opus) rồi nghiệm thu tay client Unity |
| 02/09/2026 | **Rà soát codebase bằng đồ thị + dọn rác.** Dùng skill `dothi` (graphify, tree-sitter local, 0 token) trên `server/`: **1.536 node · 4.718 cạnh · 64 cụm · 0 vòng import**. Hub thật là `KhuVuc`/`khuvuc.go` (84 cạnh) — xác nhận bằng số nỗi đau "cái gì cũng chạm khuvuc.go"; các cụm boss (`khuVucDanhBoss` 0,26, `MoBoLichBoss` 0,13) là tách gọn nhất. `TimBanDo()` 91 cạnh trông như điểm nóng nhưng soi ra **toàn test/helper gọi** — không sửa. Chặn `graphify-out/` vào `.gitignore` (`a1ad0ca`) trước khi 4MB lọt kho. **Dọn rác:** 6 file `.bak` artifact (đều untracked; artifact đã tái tạo được bằng `assetpack build`, kiểm byte-identical 28/08), `__pycache__` ×4, `.pytest_cache`, `.artifacts`, 3 worktree `wf_*` chết ở commit cũ. **Giữ** `.superpowers/sdd/*.diff` (lịch sử review, cố ý ignore). Kho chỉ còn `main` + worktree agent soi KH6b đang chạy |
| 03/09/2026 | **KH6b — soi đối kháng xong, ĐẠT BẢN 1** (`1b57a55`). Agent soi (Opus) kẹt 19 giờ: bộ chạy đột biến không timeout + vòng chờ mọi python.exe; em dừng nó, vớt 8 kết quả sẵn có, chạy nốt 20 con **có timeout 420s** (0 token, ~35 phút CPU). **28 đột biến: 22 giết, 6 sống** → soi code: M2/M16/M24 **tương đương** (lớp phòng thủ thừa), M19/M22/M26 **lỗ hổng test** → Sonnet viết 3 test (`khuvuc_dotbien_test.go`, 124k token), mỗi test áp lại đột biến → đỏ. Sonnet tự bắt xanh giả ở M26 (dòng đọc đã có `%` nên so nội dung thì y hệt — phải kiểm chỉ số). **Không có lỗi code sản xuất.** `main` xanh 8 gói. **KH6b còn duy nhất nghiệm thu tay client Unity** (chủ dự án) |
| 12b | **Sổ nợ phần thưởng boss (KH6b)** | Trước đó ai đóng góp sát thương cho boss rồi **rời khu vực / rớt mạng** đúng lúc boss chết thì **mất trắng** phần của mình — `traThuongBoss` chỉ ghi một dòng log rồi bỏ qua. Lý do khi đó đúng (actor không sở hữu con trỏ nhân vật của người đã rời — §6.9; mà đẩy bản chụp cũ qua `ThoGhi` thì **ghi đè cả dòng**, kéo lùi mọi thứ họ làm từ lúc rời đi). Nay khoản đó thành **món nợ bền** trong bảng `thuong_cho_nhan`, đi qua hàng đợi `SoNoThuong` (khuôn `ThoGhi` — actor chỉ chạm map trong bộ nhớ, **không** chạm CSDL) và được trả ở **lần đăng nhập sau**. Đánh dấu + cộng nằm trong **một giao dịch** nên không thể phát hai lần. 30 test, **đột biến 8/8 bị giết** | nhánh `so-no-thuong` |
| 28/08/2026 | **Sổ nợ phần thưởng boss — người rớt mạng đúng lúc boss chết không còn mất trắng** (nhánh `so-no-thuong`, dựng trên `cf18340` chưa gộp). Task 12 vừa xong đã chốt một hạn chế **có chủ đích**: ai đóng góp sát thương rồi rời khu vực trước lúc boss chết thì `traThuongBoss` chỉ ghi `slog.Warn` kèm số điểm rồi **bỏ**. Lý do đúng và đến giờ vẫn đúng — actor không sở hữu con trỏ `*store.NhanVat` của người đã rời (§6.9), mà đẩy một bản chụp cũ qua `ThoGhi` thì **ghi đè CẢ DÒNG**, kéo lùi mọi thứ họ làm từ lúc rời đi (vị trí, máu, vàng) để đổi lấy một khoản thưởng. **Đường thứ ba:** khoản đó thành **món nợ bền** — bảng `thuong_cho_nhan` (migration 0004) + hàng đợi `SoNoThuong` đúng khuôn `ThoGhi` (actor chỉ chạm map trong bộ nhớ vài chục nano-giây, **không** chạm CSDL) + đường trả ở lần đăng nhập sau. **Bẫy lớn nhất là cùng một cái bẫy, rình ở đầu bên kia:** `NhanNoThuong` cộng vào *cột* `kinh_nghiem` dưới CSDL, nhưng bản `*store.NhanVat` trong RAM mà actor sắp sở hữu vẫn mang số **cũ** — `ThoGhi` ghi đè cả dòng nên chỉ cần một lần lưu định kỳ 60 giây là con số cũ đó **xoá sạch phần vừa cộng**. Vì thế `traNoThuongDangCho` cộng luôn vào `nv` ngay sau khi CSDL cộng xong. **Khác `ThoGhi` một điểm cốt tử:** hàng đợi nợ **cộng dồn**, không ghi đè — `ThoGhi` ghi đè vì "bản chụp mới luôn tốt hơn bản cũ", còn ở đây mỗi khoản là một món tiền riêng. Đánh dấu `da_nhan` + cộng điểm nằm trong **một giao dịch** nên gọi hai lần ra `500` rồi `0`, không bao giờ `500` rồi `500`. **Nối dây có chủ đích tránh `main.go`** (đang thuộc task 14, có agent khác làm): mở + cắm sổ nợ trong `MoLuuTru` (điểm lắp ráp duy nhất mà dự án đã chốt), `TatMayChu` lấy nó từ `SoKhuVuc` thay vì đổi chữ ký — **không sửa `main.go` một dòng nào**. Ba test khoá riêng việc nối dây, vì bài học của chính task 12 là `CongSatThuong` từng nằm trong `boss.go` suốt từ task 3 mà không ai gọi tới. **30 test mới, mỗi mảnh thấy ĐỎ trước khi viết code; đột biến 8/8 bị giết** (bỏ `da_nhan = false` · bỏ lọc `nhan_vat_id` · hàng đợi ghi đè · vứt lô khi ghi lỗi · quay về bỏ qua người vắng mặt · cộng CSDL mà quên cộng RAM · quên cắm sổ nợ vào khu vực · tắt máy chủ không xả nốt). `go test ./...` **8/8 gói xanh**, `gofmt` + `vet` sạch; `-race` vẫn không chạy được (thiếu `gcc`, §9) nên bù `-count=5` | nhánh `so-no-thuong` |
| 28/08/2026 | **Hai agent soi đối kháng sổ nợ — tìm ra một lỗi thật mà 30 test không thấy, và một sự thật khó chịu hơn.** Theo §10 mục 5, gọi hai agent độc lập: một soi tính đúng đắn, một tấn công chính bộ test bằng đột biến tự nghĩ (27 cái mới, **11 bị giết, 16 sống sót**). **Lỗi thật:** lưới chống panic đặt cao hơn một tầng so với chỗ cần. `ghiHetHangDoi` LẤY HẲN lô nợ ra khỏi hàng đợi rồi mới gọi kho; `xepLai` chỉ chạy trên đường **trả lỗi**. `GhiNoThuong` mà **panic** thì lưới bắt được ở tận `vongLap` — lúc đó lô đã rời hàng đợi và không ai xếp lại, **cả lô nợ bốc hơi im lặng**, ngược hẳn điều chú thích của chính hàm đó cam kết ("dữ liệu vẫn nằm trong bộ nhớ nên không mất"). `ThoGhi` không dính vì nó có lưới **hai tầng** (`thogi.xa` bọc vòng lặp VÀ `thogi.ghi` bọc từng bản ghi) — tôi chỉ chép tầng ngoài. **Ba lỗi mất/nhân đôi tiền khác:** commit mập mờ (PostgreSQL commit thật rồi mới trả lỗi) làm lần thử lại chèn hai lần → thêm `khoa_lo` + ràng buộc duy nhất + `ON CONFLICT DO NOTHING`, kéo theo hàng đợi phải tách **hai tầng** (lô hỏng thử lại NGUYÊN VẸN, khoản mới đi lô riêng — trộn lại là khoá cũ mang số mới, `ON CONFLICT` nuốt mất phần mới); vòng thử lại **không có van xả** nên một dòng rác (nhân vật đã xoá → vi phạm khoá ngoại) làm nợ của **mọi người** kẹt tới lúc tắt máy rồi mất sạch → bỏ lô sau 5 lần kèm log Error; `DongVaChoXong` báo "tắt sạch" trong khi vừa đánh rơi tiền → nay trả lỗi lần xả cuối, và `TatMayChu` cho nó **hạn giờ riêng** thay vì chia chung ngân sách 30 giây. **Chú thích `GanSoNo` chỉ dẫn NGƯỢC thứ tự tắt máy** (bảo đóng sổ nợ TRƯỚC `DungTatCa`, code làm SAU và SAU mới đúng) — tôi chép từ `GanLoa` mà không sửa, đúng vố §12 dòng 27/08 lặp lại. **Điểm mù oracle** (§6.12): test lấy kỳ vọng từ chính hằng `LoaiThuongKinhNghiem` mà code đang dùng — đổi tên hằng thì **cả gói `game` vẫn xanh**; nay ghim chuỗi cứng. Thêm **20 test** (tổng 50) (rò kết nối khi giao dịch hỏng, bão hoà phía Go — bài cũ chỉ có MỘT dòng nợ nên chưa bao giờ chạm bộ cộng dồn, hai ràng buộc CHECK, index một phần, loại rỗng, nhịp không hợp lệ, đóng hai lần). **Vòng đột biến thứ hai (13 cái, tự chạy trên chính bản vá): 10 bị giết, 3 sống sót, đã bịt cả ba.** Đáng nói nhất là cái đầu — **test van xả của chính tôi viết sai**: kho giả từ chối theo *lần gọi* ("lỗi 5 lần đầu rồi thôi") nên lần thứ 6 thành công ở **cả hai** bản, có van xả hay không cũng thế. Phải có một lô hỏng **vĩnh viễn** (kho giả từ chối theo *nội dung*) mới phân biệt được. Viết lại nó còn lộ ra một điều: **hậu quả của việc thiếu van xả nhẹ hơn bản soi mô tả** — người soi nói một dòng rác làm "không khoản nợ nào của bất kỳ ai xuống được đĩa nữa", đúng với thiết kế cũ nơi lô hỏng trộn chung với khoản mới; nhưng việc tách hàng đợi hai tầng (làm vì lý do khác — chống nhân đôi) đã vô tình chữa luôn chỗ đó, nên van xả giờ chỉ còn ngăn `choThuLai` phình mãi. Đã thêm test khoá riêng điều này để nó không lặng lẽ mất đi. Hai cái sống sót còn lại: khoá lô bỏ mã phiên ngẫu nhiên (máy chủ khởi động lại là trùng khoá với lô cũ trong bảng → `ON CONFLICT` nuốt mất một khoản nợ thật, im lặng tuyệt đối) và `TatMayChu` dùng chung `ctx`. **Sự thật khó chịu nhất, và nó không phải lỗi của bản vá:** `bossNhanSatThuong` **không có ai gọi ngoài test** — task 9/10/11 chưa làm, nên máy chủ thật chưa có con boss nào ăn được đòn và sổ nợ chưa từng chạy một lần. Tôi đã viết §4 như thể tính năng đang sống; đã đính chính. Xem §8 | nhánh `so-no-thuong` |
| 03/09/2026 | **KH7 Task 2** (nhánh kh7-task2-mauvatpham): GATE giải mã hai phía -28/8 → kết luận **đủ rõ**, client ↔ Java khớp 1-1; 6 chỗ lệch đều là **tràn số im lặng phía Java** (số mẫu option là 1 byte không dấu, số mẫu vật phẩm là short **có dấu**). Ba sự thật đắt: 	emplateId/optionId **không đi qua mạng** (= chỉ số vòng lặp), client **bám cứng** option id 22→6 và 23→7 với param ×1000 và iOptionTemplates[id] **không kiểm biên**, và Service.clientOk() **chỉ** kích ở subtype 2 → subtype 2 phải gửi **cuối cùng**. Viết mauvatpham.go (8 mẫu option + 5 mẫu Naruto), 12 test byte-chính-xác + ca biên, 4/4 đột biến bị giết. **Không đụng sset-lock.json** |
| 05/09/2026 | **Khe đầu–thân: chẩn đoán lại = HAI KHUÔN MẶT + 3 tầng lệch; V1 (dữ liệu) đo 33/33 sạch; Iruka +13/+29; không cần cài công cụ; chờ chủ dự án quyết** (chi tiết §5 hàng "Khe đầu–thân", bằng chứng `docs/my-thuat/review/khe-dau-than-2026-09-05/`, script `tools/nhan-vat/thu-nghiem-khe-dau-than/`). |
| 06/09/2026 | **Khe đầu–thân + Iruka SỬA XONG bằng dữ liệu (V1), gộp `main`.** Agent A tích hợp V1 (đêm 05/09, bị cắt hạn mức giữa chừng nhưng cây đủ), B thêm bộ dò hai mặt + `--kiem`, C thử Iruka, D gộp Iruka; 3 kiểm đối kháng đạt (20 đột biến đều bị bắt), điều phối tự chạy lại mọi cửa, sửa 8 điểm lặt vặt (README, docstring, ngưỡng theo hằng, test part.bin đã commit, 3 script mẫu trỏ repo, hoàn `mapping.versions` về khối closure vì `test_deploy` ghim). Chép 148 sprite vào client, build lại máy chủ, chụp trong game: một mặt, đầu ngồi vai, Iruka bằng người chơi (`docs/my-thuat/review/khe-dau-than-2026-09-05/in-game/`). Phát hiện kèm: chết → "Xin chờ" → rớt; sói tàng hình vì client thiếu `res/Mob`. |
| 06/09/2026 | **Kế hoạch mỹ thuật đợt 4 xong** — 3 kiểm kê + 9 brief + 9 manifest JSON + README + `PROMPT-CHO-AGENT.md` (652 mục ảnh, mỗi ảnh một prompt + thư mục tự lưu). 12 agent viết, 1 agent soát đối kháng; điều phối tự kiểm chéo prompt MD↔JSON 631/631 và xử 9 điểm soát để lại. Thư mục: `docs/my-thuat/dot-4-ke-hoach-tong-the/`. |
| 06/09/2026 | **Xoá vệt hồng trên nhân vật**: 293 px màu nền còn sót (khe kín không được flood + ngưỡng đảo 598 px giữ nhầm khe khăn trán Iruka) → 0 px; `part.bin` không đổi nên `vsData` giữ 28. Kiểm đối kháng bắt 2 đột biến sống, điều phối vá lưới (test tái sinh + `chroma_mask` một nguồn sự thật + bảng màu ghim tay), áp lại 3 đột biến đều chết. |
| 06/09/2026 | **Task 2 UI Đợt 4 fix round 4/5:** test RED tái hiện highlight `(244,244,244)` chạm mép bị mask thành alpha 0 giữa nền caro `255/242`; detector mới học hai màu trội + bắt buộc nhịp lặp ngang/dọc, giữ highlight và từ chối hai dải màu không tạo caro. 7/7 test xanh; 2 mutation bị giết; ba PNG final idempotent, SHA giữ nguyên. |
| 06/09/2026 | **Task 2 UI Đợt 4 breaker fix round 5/5:** test RED fixture đúng 120×96 với viền caro `242/255`, ruột `242` đặc từng qua detector và bị xoá 11.520 pixel. Detector nay học cadence hai trục + pha 2D, đòi đúng màu/pha và coverage phân bố trên toàn ảnh trước khi mask. 8/8 test xanh; mutation bỏ riêng ratio spatial bị giết; ba PNG final idempotent và SHA giữ nguyên. |
| 08/09/2026 | **Task 3 UI Đợt 4 lắp Unity:** preflight MCP đúng project/idle; backup 15/15; sinh 12 zoom và cài đúng 15 PNG x2–x4; 5/5 test + mutation backup. Screenshot Game View thật PASS popup/button, nhưng HUD sau login BLOCKED vì máy chủ không hoạt động; không dùng production credential. `Small27083/27084` và asset-lock không đổi. |
