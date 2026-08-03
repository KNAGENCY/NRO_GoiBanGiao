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
// voi). Idempotent: da ap thi bo qua (va migration hien tai dung IF NOT EXISTS
// nen ap len DB da co san schema cung an toan). Tra ve so migration MOI ap.
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

	danhSach, err := thuMucMigration.ReadDir("migrations")
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
		noiDung, err := thuMucMigration.ReadFile("migrations/" + ph)
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
