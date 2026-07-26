# BẮT ĐẦU TỪ ĐÂY

> **File này là bản đồ + nhật ký sống của dự án.** Ai mới vào — người hay agent AI —
> đọc file này trước tiên, đọc xong là hiểu: dự án là gì, đã làm được gì, còn gì
> chưa làm, đã vấp những lỗi nào để đừng vấp lại.
>
> **Quy tắc bắt buộc:** làm xong bất cứ việc gì đáng kể (một task, một kế hoạch,
> một bản vá, một phát hiện lỗi) thì **phải cập nhật file này ngay** rồi mới commit.
> Cách cập nhật xem mục [11. Quy ước cập nhật](#11-quy-ước-cập-nhật-file-này).

**Cập nhật lần cuối:** 26/07/2026 · **Nhánh chính:** `main` · **Commit HEAD:** `18269a4`

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

Môi trường máy đang dùng: **Go ở `D:\Go`** · **PostgreSQL 16.4 cổng 5433**, CSDL `nro`
(chạy thật) và `nro_test` (kiểm thử) · Windows + PowerShell.

```powershell
# --- Máy chủ game ---
cd D:\NRO_GoiBanGiao\server
$env:NRO_DSN = "postgres://postgres:nrodev123@127.0.0.1:5433/nro"
D:\Go\bin\go.exe run ./cmd/nroserver          # lắng nghe :14445

# --- Client giả để thử, không cần mở Unity ---
D:\Go\bin\go.exe run ./cmd/thuclient -dia-chi 127.0.0.1:14445 -ten thutest -mat-khau 1

# --- Chạy toàn bộ test (PHẢI xanh trước khi commit) ---
$env:NRO_TEST_DSN = "postgres://postgres:nrodev123@127.0.0.1:5433/nro_test"
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
| W1 | Web nền tảng | Đăng ký, đăng nhập, phiên cookie, dùng chung CSDL với game | `af5a2e4` |
| W2 | Web tính năng | Trang cá nhân, đổi mật khẩu, bảng xếp hạng, panel admin (chặn quyền 3 lớp) | trên `main` |

**Tình trạng chất lượng hiện tại:** `go test ./...` — **125 test xanh, 0 đỏ**, chạy lặp
5 lần không chập chờn. `gofmt` sạch, `go vet` sạch.

**Người chơi hiện đã làm được gì:** đăng ký trên web → đăng nhập vào game → tạo nhân
vật, chọn làng → vào bản đồ làng → đi lại và nhìn thấy người chơi khác đi lại.
Chưa đánh nhau được.

---

## 5. Còn gì chưa làm

Xếp theo thứ tự nên làm:

| Ưu tiên | Việc | Ghi chú |
|---|---|---|
| 1 | **KH6 — Chiến đấu, chiêu thức, quái, boss** | Phần nặng nhất còn lại. Chưa có thiết kế |
| 2 | KH7 — Vật phẩm, hành trang, cửa hàng, giao dịch | |
| 3 | KH8 — Kho nội dung YAML + nạp nóng | Để thêm chiêu/quái mà không phải sửa code |
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

### Về cơ sở dữ liệu

6. **`lower()` của PostgreSQL dưới locale C không hạ được chữ có dấu** → `Anguoi` và
   `anguoi` (có dấu) thành hai tài khoản khác nhau. Đã vá bằng ràng buộc
   `CHECK '^[A-Za-z0-9_]{4,20}$'` — chặn được cả trò giả mạo bằng chữ Kirin/Hy Lạp
   trông giống Latin, thứ mà đổi collation **không** chặn nổi.
7. **Vòng lặp băm mật khẩu phải có trần.** Chuỗi băm hỏng làm `Verify` treo hàng phút
   (Java tình cờ an toàn nhờ số 32-bit tràn; Go 64-bit thì chạy thật). Đã chặn ở 1.000.000 vòng.

### Về luồng / đồng thời (quan trọng nhất từ KH5)

8. **Kiểm tra hợp lệ phải nằm TRONG actor.** Kế hoạch ban đầu định kiểm nước đi ở luồng
   phiên (đọc `nv.ViTriX/Y`) — nhưng actor **ghi** đúng hai trường đó → tranh chấp dữ liệu.
   Luật: **chỉ luồng `vongLap` của khu vực được chạm vào trạng thái khu vực.**
9. **Gói `-28/13` (clientOk) bị gửi lại lần hai** cũng gây tranh chấp tương tự. Đã chốt
   guard chống chạy lại; tương tự cho `-39`.

### Về kiểm thử (bài học lặp lại nhiều lần nhất)

10. **Bộ test có thể xanh trong khi code sai.** Đã xảy ra ít nhất 3 lần:
    - Xoá hẳn gói `-102` → 6/6 test vẫn đạt.
    - Xoá cả đoạn điều phối `-28/2` (tạo nhân vật) → 24/24 test vẫn xanh.
    - Test đếm số làng bằng chính nguồn dữ liệu nó phải kiểm → điểm mù oracle.
    **Cách chữa duy nhất hiệu quả: kiểm thử đột biến** — cố tình phá code, test phải đỏ.
    Không đỏ nghĩa là test vô dụng, phải viết lại test chứ không phải mừng.
11. **Test dùng giá trị ghim (hard-code), đừng lấy từ chính nguồn đang kiểm.**

### Về quy trình

12. **Đừng commit việc ngoài phạm vi task xen vào giữa một task** — gói review bị lẫn,
    người soi đọc nhầm. Commit việc lẻ **trước khi** bắt đầu task hoặc **sau khi** review xong.
13. Khi người thi công thấy thứ trông như lỗi nhưng không chắc → **dừng lại hỏi**, đừng tự
    sửa. Đã cứu được một lần (vụ byte-order §6.3).

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
- Test package `game` chậm (~54 giây) do nhiều test đầu-cuối chờ hết hạn giờ đọc.

---

## 9. Bẫy môi trường trên máy này

- **Không chạy được `go test -race`** — thiếu trình biên dịch C (cgo), không có quyền admin
  để cài. Bù lại bằng: **kỷ luật actor** (soi từng đường chạm trạng thái) + **stress test**
  (`-count=3` hoặc `-count=5`). Với mô hình actor thì cách này đủ chắc, vì tính đúng đến từ
  thiết kế "một chủ sở hữu" chứ không từ khoá.
- **`gofmt -l` báo bẩn ở vài file cũ do CRLF cục bộ** — vô hại. Đã kiểm:
  `git add --renormalize` không stage gì (blob trong repo đã là LF nhờ `.gitattributes`).
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
| Bắt đầu một phiên/agent mới | Gõ `/vao-viec` (tự đọc bối cảnh, xem ai đang làm gì, đề xuất việc tiếp) |
| Kết thúc một việc | Gõ `/ghi-so` (kiểm chứng, ghi nhật ký, cập nhật tài liệu, commit) |
| Xem ai vừa chạm vào gì | [`NHAT-KY-AGENT.md`](NHAT-KY-AGENT.md) — dòng mới nhất ở **cuối** bảng |
| Mở agent đúng thư mục | Nhấn đúp `MO-AGENT.bat` |
| Agent không phải Claude Code | Cho nó đọc [`AGENTS.md`](AGENTS.md) |

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

Cách nhanh nhất: gõ `/ghi-so`, nó làm cả 5 bước trên.

Commit riêng cho việc cập nhật tài liệu thì dùng tiền tố `docs:`.

---

## 12. Nhật ký

| Ngày | Việc |
|---|---|
| 23/07/2026 | Chốt thiết kế tổng "Naruto trên lõi NRO". Xong **KH1** (giao thức + đăng nhập), gộp `6529d35` |
| 24/07/2026 | Xong **KH2** (nhân vật) gộp `0816c89`; **KH3** (đồng bộ tĩnh, 5 làng Naruto) gộp `7a85431`; giải mã xong giao thức di chuyển |
| 25/07/2026 | Xong **KH4** (vào thế giới) gộp `aef5f71`; web **W1** + **W2** lên `main` |
| 26/07/2026 | Xong **KH5** (di chuyển + actor) gộp `e1ddc44`. Tạo file `BAT-DAU-TU-DAY.md` này làm điểm vào duy nhất cho người/agent mới |
| 26/07/2026 | Dựng cơ chế **nhiều agent song song**: sổ `NHAT-KY-AGENT.md`, lệnh `/vao-viec` + `/ghi-so`, `AGENTS.md` cho agent ngoài, `MO-AGENT.bat`, hook `SessionStart` tự nạp bối cảnh (xem §10b) |
