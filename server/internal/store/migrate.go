package store

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"sort"
)

//go:embed migrations/*.sql
var thuMucMigration embed.FS

// khoaMigration: khoa advisory CO DINH cho toan bo qua trinh migration. Hai
// server khoi dong cung luc (vd auto-restart chong nhau) se NOI TIEP nhau chu
// khong dua nhau tao bang -> tranh loi "duplicate key pg_type..." kho hieu, va
// server thu hai chi thay "da ap het" roi chay tiep.
const khoaMigration int64 = 771020304050607

// ApDungMigrationConLai ap moi migration CHUA co trong schema_migrations, theo
// thu tu ten file, moi cai mot transaction fail-closed. Xem apDungMigrationTuFS.
func (s *Store) ApDungMigrationConLai(ctx context.Context) (soApMoi int, err error) {
	return s.apDungMigrationTuFS(ctx, thuMucMigration)
}

// apDungMigrationTuFS lam viec that. Nhan fs.FS de test TIEM DUOC fs gia: vd
// migration hong -> kiem fail-closed (tra loi + KHONG ghi phien ban do); vd
// migration phu thuoc nhau -> kiem ap dung thu tu. Idempotent: da ap thi bo qua;
// migration hien tai dung IF NOT EXISTS nen ap len DB da co san schema cung an toan.
func (s *Store) apDungMigrationTuFS(ctx context.Context, fsys fs.FS) (soApMoi int, err error) {
	// Advisory lock: chi MOT server chay migration tai mot thoi diem. Server thu
	// hai CHO o day roi thay da ap het -> soApMoi=0, chay tiep binh thuong (khong
	// chet vi dua nhau tao bang). Khoa giu tren mot ket noi rieng suot ca ham.
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return 0, fmt.Errorf("lay ket noi cho migration that bai: %w", err)
	}
	defer conn.Release()
	if _, err := conn.Exec(ctx, "SELECT pg_advisory_lock($1)", khoaMigration); err != nil {
		return 0, fmt.Errorf("lay advisory lock migration that bai: %w", err)
	}
	defer func() {
		// Go khoa du thanh cong hay loi. context.Background vi ctx goc co the da huy.
		_, _ = conn.Exec(context.Background(), "SELECT pg_advisory_unlock($1)", khoaMigration)
	}()

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

	danhSach, err := fs.ReadDir(fsys, "migrations")
	if err != nil {
		return 0, fmt.Errorf("doc thu muc migration that bai: %w", err)
	}
	var ten []string
	for _, f := range danhSach {
		if !f.IsDir() {
			ten = append(ten, f.Name())
		}
	}
	sort.Strings(ten) // 0001 truoc 0002 — thu tu ap PHAI on dinh

	for _, ph := range ten {
		if daAp[ph] {
			continue
		}
		noiDung, err := fs.ReadFile(fsys, "migrations/"+ph)
		if err != nil {
			return soApMoi, fmt.Errorf("doc file migration %s that bai: %w", ph, err)
		}
		tx, err := s.pool.Begin(ctx)
		if err != nil {
			return soApMoi, fmt.Errorf("mo transaction cho %s that bai: %w", ph, err)
		}
		// Exec KHONG tham so -> pgx dung simple protocol -> chay duoc file NHIEU
		// cau lenh (CREATE TABLE + INDEX + DO block) trong mot lan goi.
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
