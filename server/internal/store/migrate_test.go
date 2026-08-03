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

func TestApDungMigration_GhiDungPhienBan(t *testing.T) {
	s := moStoreTest(t)
	ctx := context.Background()
	if _, err := s.ApDungMigrationConLai(ctx); err != nil {
		t.Fatalf("ap migration: %v", err)
	}
	for _, ph := range []string{"0001_tai_khoan.sql", "0002_nhan_vat.sql"} {
		var co bool
		if err := s.pool.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE phien_ban = $1)`, ph).Scan(&co); err != nil {
			t.Fatalf("kiem phien ban %s: %v", ph, err)
		}
		if !co {
			t.Fatalf("schema_migrations thieu phien ban %q", ph)
		}
	}
}
