# BẮT ĐẦU TỪ ĐÂY

> **File này là bản đồ + nhật ký sống của dự án.** Ai mới vào — người hay agent AI —
> đọc file này trước tiên, đọc xong là hiểu: dự án là gì, đã làm được gì, còn gì
> chưa làm, đã vấp những lỗi nào để đừng vấp lại.
>
> **Quy tắc bắt buộc:** làm xong bất cứ việc gì đáng kể (một task, một kế hoạch,
> một bản vá, một phát hiện lỗi) thì **phải cập nhật file này ngay** rồi mới commit.
> Cách cập nhật xem mục [11. Quy ước cập nhật](#11-quy-ước-cập-nhật-file-này).

**Cập nhật lần cuối:** 27/08/2026 (**Cổng rời Làng Lá + 8 bóng đổ đã lắp được** — người chơi có đường ra khỏi map 1; và **gỡ quả mìn `tile_set_info.bin`** — build lại không còn làm nhân vật rơi xuyên đất)

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

---

## 5. Còn gì chưa làm

Xếp theo thứ tự nên làm:

| Ưu tiên | Việc | Ghi chú |
|---|---|---|
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
| 1b | **KH6b — Boss** (hệ riêng, tách khỏi KH6) | Bản Java là 12 file, phức tạp hơn quái thường nhiều. Thiết kế sau khi KH6 gộp |
| ✅ | ~~KH6c — Stress test 20 client cùng khu vực~~ | **Xong, gộp `322b140`.** Actor không treo dưới ~5 triệu sự kiện/2s, giữ nhịp, gạch đúng client chậm; thêm đo trễ nhịp (`soNhipTre`/`treNhipLauNhatMs`). Dùng test in-process (`suKienDiChuyen` + `phienStress`), không cần TCP thật |
| 2 | KH7 — Vật phẩm, hành trang, cửa hàng, giao dịch | 🧱 **Phải chốt căn cước vật phẩm + chống nhân bản trong đặc tả TRƯỚC khi viết code** ([BO-SUNG](docs/superpowers/plans/BO-SUNG-LO-TRINH.md) Khối 3). Tách rõ 4 thứ: cửa hàng NPC / giao dịch 2 người / ký gửi / kho |
| 3 | KH8 — Kho nội dung YAML + nạp nóng | Để thêm chiêu/quái mà không phải sửa code. **Kèm cân bằng thực nghiệm** (đơn vị `SatThuong` phải chốt — TODO trong `satthuong.go`) |
| 4 | KH9 — Nhiệm vụ, cốt truyện, NPC | |
| 5 | KH10 — Cảm giác chơi (hoạt ảnh, chuyển màn, hướng dẫn người mới) | Ba con số ràng buộc đã chốt ở KH5 |
| — | **W3 — Web:** cấu hình / nhật ký / nạp thẻ, phân trang, thu hồi phiên JWT | Một phần phụ thuộc hệ thống game chưa có |
| — | **B — Bảo mật** (xem §7) | **Phải xong trước khi mở công khai** |
| — | **V — Triển khai VPS**, giám sát, phát hành | |

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
  hoặc xác minh email, game chưa có cả hai; (b) **rò rỉ thời gian phía `web/` còn nguyên** — mới vá phía
  game; cách chữa giống hệt: băm giả ở nhánh không tìm thấy tài khoản; (c) chưa chạy được `go test -race`
  (máy không có `gcc`), hai test đồng thời đã viết sẵn chờ máy có trình biên dịch C.
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
