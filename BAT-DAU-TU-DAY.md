# BẮT ĐẦU TỪ ĐÂY

> **File này là bản đồ + nhật ký sống của dự án.** Ai mới vào — người hay agent AI —
> đọc file này trước tiên, đọc xong là hiểu: dự án là gì, đã làm được gì, còn gì
> chưa làm, đã vấp những lỗi nào để đừng vấp lại.
>
> **Quy tắc bắt buộc:** làm xong bất cứ việc gì đáng kể (một task, một kế hoạch,
> một bản vá, một phát hiện lỗi) thì **phải cập nhật file này ngay** rồi mới commit.
> Cách cập nhật xem mục [11. Quy ước cập nhật](#11-quy-ước-cập-nhật-file-này).

**Cập nhật lần cuối:** 31/07/2026 (Task 8 đã được independent audit APPROVE; Vertical slice Task 9 telemetry/control/runner đã hoàn tất và review độc lập sạch) · **Nhánh đang làm:** `kh-vertical-slice` · **Mốc trước Task 9:** `b82d950` · **main:** `485b095`

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
D:\Go\bin\gofmt.exe -l internal/ cmd/         # phải không in gì (trừ nhiễu CRLF, xem §9)
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
| VS-9 | **Telemetry Unity, server control và runner** | Telemetry JSONL loopback có token/command ID/sequence và main-thread dispatch; session cũ không thể lọt lệnh sang client mới. Server control chỉ mở bằng env rõ ràng, bind loopback và tắt qua đúng `donDepKhiTat`. Fixture chỉ dùng `nro_test`, đúng 2 account/run và cleanup fail-safe. Runner có preflight/timeout/evidence SHA. Unity contract thật **50 assertion**, Go toàn repo + vet, PowerShell self-test và control smoke DB thật đều xanh; review độc lập sạch | commit Task 9 này |

**Tình trạng chất lượng hiện tại:** `go test ./...` — **toàn bộ 6 gói xanh, 0 đỏ**
(gói `game` ~70s), chạy lặp `-count=5` không chập chờn. `gofmt` sạch (trừ nhiễu CRLF file
cũ, §9), `go vet` sạch. Vertical slice Task 7 có full `python -m unittest discover -s assetpack/tests -v`
**123/123 xanh** và preflight Unity staging thật. KH6 còn được **soi toàn nhánh (opus)** + **2 canary** (thả bug thật
để kiểm: reviewer bắt được lỗi byte im lặng, test E2E bắt được lỗi tích hợp).

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
| ✅ | ~~Vertical slice Task 9 — telemetry/control/runner~~ | **Xong trên `kh-vertical-slice`, review độc lập sạch.** Unity contract 50 assertion; Go toàn repo + vet; PowerShell parser/self-test; control smoke dùng DB `nro_test`, graceful shutdown và cleanup fixture đều xanh. Full gameplay hai client **chưa chạy** vì thuộc Task 10. |
| 🔒 | **Vertical slice Task 10 — gameplay E2E hai client** | Việc tiếp theo theo kế hoạch; **chưa bắt đầu**. Runner Task 9 đã chuẩn bị cổng telemetry/control/fixture/evidence, nhưng không được coi control smoke là gameplay PASS. |
| 1b | **KH6b — Boss** (hệ riêng, tách khỏi KH6) | Bản Java là 12 file, phức tạp hơn quái thường nhiều. Thiết kế sau khi KH6 gộp |
| 1c | **KH6c — Stress test 20 client cùng khu vực** | Cùng xả chiêu, đo gói rơi hộp thư actor (trần 256, rơi im lặng) + trễ nhịp 100ms. Dựng trên `thuclient`/`clientGia` |
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
10b. **Actor tự đi phát gói mạng = một kết nối treo làm đứng cả khu vực** (phát hiện
    27/07 khi AI quái của Task 9 thành nguồn lưu lượng tự phát). `phatChoNguoiKhac` /
    `phatChoCaKhuVuc` gọi `phien.Gui()` **đồng bộ ngay trên goroutine `vongLap`**; `Gui()`
    có hạn giờ ghi **30 giây** (`net/session.go`). Một client ngừng đọc (mạng chập chờn,
    app xuống nền, hoặc cố ý) làm `vongLap` đứng tối đa 30s → cả bản đồ đứng hình, quái
    ngừng hồi sinh. Tệ hơn: `Gui()` lỗi chỉ ghi log, **không gạch tên** người đó khỏi
    `nguoiTrong` → lần phát sau lại treo tiếp, **lặp vô hạn**. Thiết kế KH5 đã yêu cầu gỡ
    phiên hỏng nhưng phần đó **chưa vào code**. Đây là **lỗ hổng Nghiêm trọng #1** của
    review toàn dự án 27/07 — **Kế hoạch B phải tách gửi mạng khỏi actor** (mỗi phiên một
    hàng đợi gửi riêng, actor chỉ ném gói vào rồi đi tiếp). Task 9 né tạm trong test, chưa
    sửa gốc.

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

---

## 7. Nợ bảo mật — phải xử trước khi mở công khai

Gom lại thành **Kế hoạch B**, sửa lẻ từng cái là vô nghĩa:

- **Chưa chống dò mật khẩu.** Đo thật: **200 lần sai liên tiếp không bị chặn, 76 lần/giây**.
  Bản Java cũ *có* AntiLogin (15 lần sai → khoá IP 15 phút) → đây là **bước lùi**.
- **Rò rỉ thời gian ~13,5 ms** giữa nhánh "không có tài khoản" và "sai mật khẩu" → đoán được
  tên tài khoản nào tồn tại. Có ở **cả web lẫn game**. Cách chữa: băm giả ở nhánh không tìm thấy.
- **Chưa giới hạn số kết nối trên một IP.**
- Web: JWT **không thu hồi phiên** khi tài khoản bị khoá / hạ quyền (hiệu lực 7 ngày);
  đổi mật khẩu chưa kiểm cờ `bi_khoa`; `npm audit` 3 lỗi high (nằm trong bundle Next, chưa
  phải bề mặt tấn công).
- Di sản (nếu còn dùng phần cũ): chưa có HTTPS, chưa sao lưu DB tự động, **không được mở
  cổng 3306 ra internet**.

---

## 8. Nợ kỹ thuật lặt vặt (không chặn việc gì)

- `Dung()` gọi trước `BatDau()` sẽ panic (cạnh biên, chưa ai gặp).
- Client gửi `-39` trước `-28/13` → không làm gì (im lặng, đúng nhưng chưa có test).
- Đăng nhập trùng tài khoản ghi đè người cũ trong khu vực (có từ trước).
- `byte(nv.Cap)` cắt cụt nếu cấp > 255.
- `context.Background()` ở `login.go`: truy vấn CSDL không hạn giờ, không huỷ khi client ngắt.
- Vòng `Accept()` không nghỉ giãn → cạn file descriptor sẽ quay nóng CPU.
- Mã chết: `Writer.Bool`, `Writer.Long`; `TaiKhoan.LaAdmin` đọc từ CSDL nhưng chưa ai dùng.
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
  không câm hẳn). Kết hợp với **§6.10b** (actor bị chặn tới 30 giây vì tự gửi gói mạng) thì
  lúc tắt máy chủ, một khu vực đang dính lỗi đó sẽ mất tới 60 giây tiến độ. **Thêm một lý do
  để Kế hoạch B ưu tiên tách gửi gói mạng khỏi actor.**
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

- **Không chạy được `go test -race`** — thiếu trình biên dịch C (cgo), không có quyền admin
  để cài. Bù lại bằng: **kỷ luật actor** (soi từng đường chạm trạng thái) + **stress test**
  (`-count=3` hoặc `-count=5`). Với mô hình actor thì cách này đủ chắc, vì tính đúng đến từ
  thiết kế "một chủ sở hữu" chứ không từ khoá.
- **`gofmt -l` báo bẩn ở vài file cũ do CRLF cục bộ** — vô hại. Đã kiểm:
  `git add --renormalize` không stage gì (blob trong repo đã là LF nhờ `.gitattributes`).
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
