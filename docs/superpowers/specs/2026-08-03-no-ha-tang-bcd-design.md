# Nợ hạ tầng B/C/D — Thiết kế

**Ngày:** 03/08/2026 · **Trạng thái:** thiết kế (chờ viết plan) · **Quyết định đã chốt với chủ dự án:** DB tự-host trên VPS (mảnh D dùng `pg_dump`).

## 1. Vì sao + 3 mảnh

Sau lưới chắn panic (A) và KH6c, nốt phần "nợ hạ tầng" còn lại để **an toàn trước khi mở cửa**:
- **C — Migration DB:** hết cảnh gõ `psql` tay; ai chạy server cũng có schema đúng, tự động.
- **B — CI:** một lệnh chạy hết kiểm tra + chặn push khi đỏ → không để hồi quy lọt.
- **D — Sao lưu DB:** `pg_dump` định kỳ, giữ N bản, tự dọn bản cũ → mất máy không mất người chơi.

Đây là **3 sub-project độc lập**, làm **theo thứ tự C → B → D** (C là code thật, làm kỹ nhất; B bao trùm; D là script). Một spec gộp vì cùng nhóm + nhỏ.

## 2. Thực trạng (đã soi)

- **Không có remote git.** → CI local (git hook) chạy ngay; kèm sẵn file GitHub Actions cho tương lai.
- **2 migration** `0001_tai_khoan.sql`, `0002_nhan_vat.sql` — **đều dùng `IF NOT EXISTS`** (áp lại an toàn). Chạy tay bằng `psql`; server **không tự áp**.
- **pytest** ở `assetpack/tests/`. Go ở `D:\Go`, Postgres `D:\PostgreSQL` (cổng tự dò 5432/5433).
- Chưa có CI, chưa có backup.

## 3. Mảnh C — Migration runner (Go, TDD)

**File mới:** `server/internal/store/migrate.go` (+ `migrate_test.go`).

- `//go:embed migrations/*.sql` — nhúng SQL vào binary (ai clone/chạy đều có, không lệ thuộc file rời).
- Hàm `ApDungMigrationConLai(ctx, db)` (chữ ký khớp handle DB của `store`):
  1. `CREATE TABLE IF NOT EXISTS schema_migrations (phien_ban TEXT PRIMARY KEY, ap_luc TIMESTAMPTZ NOT NULL DEFAULT now())`.
  2. Đọc tên file migration đã nhúng, **sắp tăng dần** (`0001…`, `0002…`).
  3. Đọc danh sách `phien_ban` đã áp.
  4. Với mỗi migration **chưa áp** (đúng thứ tự): `BEGIN` → chạy SQL → `INSERT` bản ghi phiên bản → `COMMIT`. Lỗi bất kỳ → **rollback + trả lỗi, DỪNG** (fail-closed: không chạy tiếp, không để schema nửa vời).
  5. Log mỗi migration áp + tổng kết ("áp N, bỏ qua M đã có").
- **Gắn vào khởi động** `cmd/nroserver/main.go`: sau khi nối DB, **trước khi lắng nghe**, gọi `ApDungMigrationConLai`. Lỗi → log tiếng Việt rõ + thoát mã khác 0 (**không** mở cổng với schema sai).
- **Idempotent:** chạy lần 2 áp 0 cái. Vì `schema_migrations` chặn chạy lại, và migration hiện tại `IF NOT EXISTS` nên áp lên DB đã có schema (máy dev) cũng an toàn — tự "nhận" 2 migration cũ vào bảng theo dõi.

**Test (`migrate_test.go`, cần `NRO_TEST_DSN` → tự SKIP nếu không có):**
- Dọn sạch `schema_migrations` trước mỗi test (tất định).
- Áp lần đầu → `schema_migrations` có đúng 2 phiên bản, bảng `tai_khoan`/`nhan_vat` tồn tại.
- Áp lần 2 → **0 cái mới** (idempotent).
- Fail-closed: nhét một migration SQL hỏng qua FS test → trả lỗi, **không** ghi phiên bản đó (transaction rollback).
- **Đột biến:** bỏ ghi `schema_migrations` → test idempotent đỏ (áp lại 2 cái); bỏ `ORDER` → test thứ tự đỏ.

⚠️ **Bẫy:** DO-block cuối `0001` (ràng buộc CHECK) — phải chắc nó cũng chạy-lại-an-toàn (đọc kỹ khi code; nếu không thì bọc điều kiện tồn tại). Đọc file trước khi viết runner.

## 4. Mảnh B — CI (script + hook + GitHub Actions sẵn)

**File mới:** `kiem-tra-tat-ca.ps1` (gốc repo) · `scripts/hooks/pre-push` · `.github/workflows/ci.yml`.

- **`kiem-tra-tat-ca.ps1`** — trái tim CI, chạy TẤT CẢ, đỏ chỗ nào thoát mã khác 0:
  1. `go build ./...` (trong `server/`)
  2. `go vet ./...`
  3. `gofmt -l internal/ cmd/` → **có tên file = đỏ** (trừ nhiễu CRLF, xem §9 BAT-DAU)
  4. `go test ./...` — **tự dò cổng Postgres**, đặt `NRO_TEST_DSN` nếu Postgres sống (không thì test cần DB tự SKIP)
  5. `pytest` trong `assetpack/`
  In tóm tắt PASS/FAIL từng bước + tổng.
- **`scripts/hooks/pre-push`** — gọi `kiem-tra-tat-ca.ps1`; đỏ thì **chặn push**. Cài bằng `git config core.hooksPath scripts/hooks` (hook nằm trong repo, version-controlled, chia sẻ được — khác `.git/hooks` cục bộ). Kèm dòng lệnh cài + ghi `README`.
- **`.github/workflows/ci.yml`** — chạy khi push/PR: `ubuntu`, `setup-go`, `setup-python`, service `postgres`, chạy `go test` + `pytest`. **Nằm im tới khi có remote GitHub** rồi tự hoạt động. YAGNI: không tự tạo remote.

**Kiểm chứng:** chạy `kiem-tra-tat-ca.ps1` thấy xanh; cố tình phá (gofmt lệch / test đỏ) → script thoát khác 0. Hook: thử `git push` giả (dry-run) thấy nó gọi kiểm tra.

## 5. Mảnh D — Sao lưu DB (self-host, pg_dump)

**File mới:** `sao-luu-db.ps1` · `cai-lich-sao-luu.ps1` · mục "khôi phục" trong README.

- **`sao-luu-db.ps1`**: tự dò cổng Postgres → `pg_dump` DB `nro` ra `backups/nro-YYYYMMDD-HHMMSS.dump` (định dạng custom `-Fc`, nén + `pg_restore` được). Giữ **N bản gần nhất** (mặc định 14), **xóa bản cũ hơn**. Log rõ. `backups/` vào `.gitignore` (không commit dump).
- **`cai-lich-sao-luu.ps1`**: đăng ký **Windows Scheduled Task** chạy `sao-luu-db.ps1` **hàng ngày** (giờ thấp điểm). Chạy một lần trên VPS.
- **Khôi phục** (ghi README): `pg_restore -d nro <file.dump>` (hoặc tạo DB mới rồi restore). Nêu rõ cảnh báo ghi đè.

**Kiểm chứng:** chạy `sao-luu-db.ps1` tay → ra file `.dump` > 0 byte, `pg_restore --list` đọc được; chạy nhiều lần → chỉ giữ N bản mới nhất.

## 6. Ngoài phạm vi (YAGNI)

- **Không** tạo remote GitHub / không đẩy code lên đâu (chỉ chuẩn bị sẵn file CI).
- **Không** làm rollback migration (chỉ tiến; rollback là phức tạp, chưa cần khi chưa mở cửa).
- **Không** backup sang cloud/offsite (chỉ local trên VPS — đủ cho giai đoạn này; offsite là việc sau).
- **Không** đụng Supabase (đã chốt tự-host).

## 7. Rủi ro / bẫy

- **Migration tự áp khi khởi động** có thể chặn server lên nếu một migration hỏng — **đúng ý muốn** (fail-closed hơn là chạy với schema sai), nhưng phải log tiếng Việt cực rõ để biết đường sửa.
- **DO-block `0001`** phải chạy-lại-an-toàn (xem §3 bẫy).
- **Script PowerShell sửa file tiếng Việt**: tự thân script không sinh file có dấu, nhưng nhớ §6.15 (đừng để PowerShell đọc/ghi file dấu bằng codepage ANSI).
- **Backup không được chặn game**: `pg_dump` đọc nhất quán không khóa ghi lâu (Postgres MVCC) — an toàn chạy lúc server đang chạy.
- **`gofmt -l` trên Windows** hay báo nhiễu CRLF file cũ — script chỉ nên soi file mình quản (`internal/`, `cmd/`), và lọc nhiễu CRLF như §9 BAT-DAU.
