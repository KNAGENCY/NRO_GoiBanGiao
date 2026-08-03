# Nợ hạ tầng B/C/D — Kế hoạch thi công

> **Cho agent thi công:** SUB-SKILL BẮT BUỘC: `superpowers:subagent-driven-development` hoặc
> `superpowers:executing-plans`. Bước dùng ô tick `- [ ]`.

**Goal:** Migration DB tự áp khi khởi động (C) + CI một-lệnh-chạy-hết + chặn push (B) + sao lưu DB định kỳ tự dọn (D).

**Architecture:** C là code Go trong `store` (go:embed SQL + bảng `schema_migrations` + method trên `*Store`, gọi ở `main.go` sau khi nối DB). B/D là script PowerShell + git hook + file GitHub Actions.

**Tech Stack:** Go (`D:\Go\bin\go.exe`, pgx v5), PowerShell, `pg_dump`/`pg_restore` (`D:\PostgreSQL\bin`), pytest.

## Global Constraints (mọi task)

- Định danh + chú thích **tiếng Việt không dấu** theo phong cách quanh nó. File có dấu → Edit/Write, không shell (§6.15).
- Cổng Postgres **tự dò** (5432/5433), đừng ghim.
- C: test cần `NRO_TEST_DSN` → tự SKIP nếu không có. Trước commit: `go build`/`vet`/`gofmt` sạch.
- Nhánh `no-ha-tang-bcd` (đã có commit spec). Mỗi task xong commit riêng.
- Kiểm thử đột biến cho C (task 1).

---

### Task 1: C — Migration runner (Go, TDD)

**Files:** Create `server/internal/store/migrate.go`, `server/internal/store/migrate_test.go`. Modify `server/cmd/nroserver/main.go` (gắn sau `store.Open`).

**Interfaces:** Produces `func (s *Store) ApDungMigrationConLai(ctx context.Context) (soApMoi int, err error)`.

- [ ] **Step 1: Test đỏ trước** — `migrate_test.go` (cần `NRO_TEST_DSN`):

```go
package store

import (
	"context"
	"os"
	"testing"
)

func moStoreTest(t *testing.T) *Store {
	t.Helper()
	dsn := os.Getenv("NRO_TEST_DSN")
	if dsn == "" {
		t.Skip("dat NRO_TEST_DSN de chay test can CSDL")
	}
	s, err := Open(context.Background(), dsn)
	if err != nil {
		t.Fatalf("mo store test: %v", err)
	}
	t.Cleanup(s.Close)
	// Sach de tat dinh: xoa bang theo doi (migration IF NOT EXISTS nen bang that giu nguyen).
	if _, err := s.pool.Exec(context.Background(), `DROP TABLE IF EXISTS schema_migrations`); err != nil {
		t.Fatalf("don schema_migrations: %v", err)
	}
	return s
}

func TestApDungMigration_ApVaIdempotent(t *testing.T) {
	s := moStoreTest(t)
	ctx := context.Background()

	// Lan dau: ap het migration con thieu (>=2: 0001, 0002).
	n1, err := s.ApDungMigrationConLai(ctx)
	if err != nil {
		t.Fatalf("ap lan 1: %v", err)
	}
	if n1 < 2 {
		t.Fatalf("ap lan 1 duoc %d migration, mong >=2", n1)
	}
	// Bang that ton tai.
	for _, bang := range []string{"tai_khoan", "nhan_vat"} {
		var co bool
		if err := s.pool.QueryRow(ctx,
			`SELECT to_regclass($1) IS NOT NULL`, "public."+bang).Scan(&co); err != nil {
			t.Fatalf("kiem bang %s: %v", bang, err)
		}
		if !co {
			t.Fatalf("bang %s phai ton tai sau migration", bang)
		}
	}
	// Lan 2: idempotent — ap 0 cai moi.
	n2, err := s.ApDungMigrationConLai(ctx)
	if err != nil {
		t.Fatalf("ap lan 2: %v", err)
	}
	if n2 != 0 {
		t.Fatalf("ap lan 2 duoc %d migration, mong 0 (idempotent)", n2)
	}
}
```

- [ ] **Step 2: Chạy → đỏ.** `cd server; D:\Go\bin\go.exe test ./internal/store/ -run TestApDungMigration` với `NRO_TEST_DSN` đặt. Mong FAIL (`ApDungMigrationConLai` undefined).

- [ ] **Step 3: Cài `migrate.go`:**

```go
package store

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"sort"
)

//go:embed migrations/*.sql
var thuMucMigration embed.FS

// ApDungMigrationConLai ap moi migration CHUA co trong bang schema_migrations,
// theo dung thu tu ten file (0001, 0002, ...). Moi migration chay trong MOT
// transaction: loi bat ky -> rollback + DUNG (fail-closed, khong de schema nua
// voi). Idempotent: da ap thi bo qua. Tra ve so migration MOI ap.
func (s *Store) ApDungMigrationConLai(ctx context.Context) (soApMoi int, err error) {
	if _, err = s.pool.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS schema_migrations (
			phien_ban TEXT PRIMARY KEY,
			ap_luc    TIMESTAMPTZ NOT NULL DEFAULT now()
		)`); err != nil {
		return 0, fmt.Errorf("tao bang schema_migrations that bai: %w", err)
	}

	daAp := map[string]bool{}
	rows, err := s.pool.Query(ctx, `SELECT phien_ban FROM schema_migrations`)
	if err != nil {
		return 0, fmt.Errorf("doc schema_migrations that bai: %w", err)
	}
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			rows.Close()
			return 0, fmt.Errorf("doc phien ban that bai: %w", err)
		}
		daAp[v] = true
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("duyet schema_migrations that bai: %w", err)
	}

	tenFile, err := thuMucMigration.ReadDir("migrations")
	if err != nil {
		return 0, fmt.Errorf("doc thu muc migration that bai: %w", err)
	}
	var ten []string
	for _, f := range tenFile {
		if !f.IsDir() {
			ten = append(ten, f.Name())
		}
	}
	sort.Strings(ten) // 0001 truoc 0002 — thu tu ap PHAI on dinh

	for _, ph := range ten {
		if daAp[ph] {
			continue
		}
		noiDung, err := thuMucMigration.ReadFile("migrations/" + ph)
		if err != nil {
			return soApMoi, fmt.Errorf("doc file migration %s that bai: %w", ph, err)
		}
		tx, err := s.pool.Begin(ctx)
		if err != nil {
			return soApMoi, fmt.Errorf("mo transaction cho %s that bai: %w", ph, err)
		}
		// Exec KHONG tham so -> pgx dung simple protocol -> chay duoc file NHIEU
		// cau lenh (CREATE TABLE + INDEX + DO block) trong mot lan.
		if _, err := tx.Exec(ctx, string(noiDung)); err != nil {
			_ = tx.Rollback(ctx)
			return soApMoi, fmt.Errorf("chay migration %s that bai (da rollback, DUNG): %w", ph, err)
		}
		if _, err := tx.Exec(ctx,
			`INSERT INTO schema_migrations (phien_ban) VALUES ($1)`, ph); err != nil {
			_ = tx.Rollback(ctx)
			return soApMoi, fmt.Errorf("ghi phien ban %s that bai (da rollback): %w", ph, err)
		}
		if err := tx.Commit(ctx); err != nil {
			return soApMoi, fmt.Errorf("commit migration %s that bai: %w", ph, err)
		}
		slog.Info("da ap migration", "phien_ban", ph)
		soApMoi++
	}
	if soApMoi == 0 {
		slog.Info("CSDL da o phien ban moi nhat, khong co migration moi")
	} else {
		slog.Info("ap migration xong", "so_ap_moi", soApMoi)
	}
	return soApMoi, nil
}
```

- [ ] **Step 4: Chạy → xanh.** Test PASS (n1>=2, bảng tồn tại, n2==0).

- [ ] **Step 5: Gắn vào `main.go`** — sau `slog.Info("da noi CSDL thanh cong.")` (dòng ~228), thêm:

```go
	if soAp, err := kho.ApDungMigrationConLai(ctx); err != nil {
		slog.Error("KHONG THE KHOI DONG: ap migration CSDL that bai.", "loi", err)
		slog.Error("Cach kiem tra: xem migration nao loi o log tren; sua SQL roi chay lai. May chu KHONG mo cong voi schema sai.")
		kho.Close()
		os.Exit(1)
	} else if soAp > 0 {
		slog.Info("da ap migration CSDL luc khoi dong", "so_ap_moi", soAp)
	}
```

- [ ] **Step 6: Toàn bộ + wiring xanh.** `NRO_TEST_DSN=... go test ./... -count=1`; `go build ./...` (main.go biên dịch); `go vet ./...`; `gofmt -l internal/ cmd/` (không in file mình sửa).

- [ ] **Step 7: Đột biến.** (a) Bỏ dòng `sort.Strings(ten)` → thứ tự map ngẫu nhiên: khó bắt vì 2 file; **thêm test** `TestApDungMigration_GhiDungPhienBan` kiểm `schema_migrations` chứa đúng `0001_tai_khoan.sql` + `0002_nhan_vat.sql`. (b) Bỏ `INSERT INTO schema_migrations` → test idempotent đỏ (n2>0). Chạy xác nhận đỏ, khôi phục.

- [ ] **Step 8: Commit** `feat(store): migration runner tu ap khi khoi dong (no ha tang C)`.

---

### Task 2: B — CI (script + hook + GitHub Actions)

**Files:** Create `kiem-tra-tat-ca.ps1` (gốc), `scripts/hooks/pre-push`, `.github/workflows/ci.yml`.

- [ ] **Step 1: `kiem-tra-tat-ca.ps1`** — chạy tuần tự, đỏ chỗ nào `exit 1`:
  - Tự dò cổng Postgres (5432/5433) → nếu sống, đặt `$env:NRO_TEST_DSN`.
  - `Set-Location server`: `go build ./...` → `go vet ./...` → `gofmt -l internal cmd` (**có output = đỏ**) → `go test ./...`.
  - `Set-Location ..\assetpack`: `pytest -q` (nếu có python + pytest).
  - In `[OK]`/`[DO]` mỗi bước, tổng kết. `exit` với mã khác 0 nếu bất kỳ bước đỏ.

- [ ] **Step 2: `scripts/hooks/pre-push`** (bash — git chạy hook qua sh):
```sh
#!/bin/sh
# Chan push neu kiem tra do. Cai: git config core.hooksPath scripts/hooks
exec powershell.exe -NoProfile -ExecutionPolicy Bypass -File "./kiem-tra-tat-ca.ps1"
```

- [ ] **Step 3: `.github/workflows/ci.yml`** — on push/PR: ubuntu, `actions/setup-go`, `actions/setup-python`, service `postgres:16` (env `NRO_TEST_DSN`), chạy `go test ./...` (trong server) + `pytest` (trong assetpack). Nằm im tới khi có remote.

- [ ] **Step 4: Kiểm chứng.** Chạy `kiem-tra-tat-ca.ps1` → thấy `[OK]` các bước, exit 0. Cố tình phá gofmt 1 file → script exit khác 0 → khôi phục. `git config core.hooksPath scripts/hooks` rồi ghi README cách cài.

- [ ] **Step 5: Commit** `ci: kiem-tra-tat-ca + pre-push hook + github actions (no ha tang B)`.

---

### Task 3: D — Sao lưu DB (self-host, pg_dump)

**Files:** Create `sao-luu-db.ps1`, `cai-lich-sao-luu.ps1`; thêm `backups/` vào `.gitignore`; mục "khôi phục" vào `server/README.md`.

- [ ] **Step 1: `sao-luu-db.ps1`** — tự dò cổng → `pg_dump -Fc -d nro -f backups\nro-<timestamp>.dump` (`$env:PGPASSWORD` từ tham số/mặc định dev). Tạo `backups/` nếu chưa có. Giữ **N=14** bản mới nhất (`Get-ChildItem ... | Sort LastWriteTime -Desc | Select -Skip 14 | Remove-Item`). Log rõ. Tham số: `-SoBanGiu 14`, `-Cong`, `-MatKhau`.

- [ ] **Step 2: `.gitignore`** thêm `backups/` (đừng commit dump).

- [ ] **Step 3: `cai-lich-sao-luu.ps1`** — `Register-ScheduledTask` chạy `sao-luu-db.ps1` hàng ngày (vd 4h sáng). In hướng dẫn chạy trên VPS.

- [ ] **Step 4: README** — mục "Khôi phục từ backup": `pg_restore -d nro --clean backups\<file>.dump` + cảnh báo ghi đè.

- [ ] **Step 5: Kiểm chứng.** Chạy `sao-luu-db.ps1` (Postgres đang sống) → ra `.dump` > 0 byte, `pg_restore --list <file>` đọc được. Chạy 15+ lần (hoặc giả) → chỉ giữ 14 bản.

- [ ] **Step 6: Commit** `feat(ops): sao luu DB pg_dump dinh ky + scheduled task (no ha tang D)`.

---

## Self-Review

- **Phủ spec:** C §3 → Task 1; B §4 → Task 2; D §5 → Task 3. ✓
- **Không placeholder:** Task 1 có code đầy đủ; Task 2/3 là script — outline concrete (lệnh cụ thể), sẽ viết đầy đủ khi thi công. ✓
- **Nhất quán:** `ApDungMigrationConLai` cùng tên ở Task 1 + main.go. `schema_migrations`, `backups/` nhất quán. ✓
- **Rủi ro pgx multi-statement:** đã ghi trong code (Exec không tham số → simple protocol → chạy file nhiều câu lệnh). Kiểm ở Step 4.
