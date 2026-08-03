package store

import (
	"context"
	"os"
	"testing"
	"testing/fstest"
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

// TestApDungMigration_FailClosed khoa thuoc tinh QUAN TRONG NHAT (reviewer F1):
// migration hong -> TRA LOI + KHONG ghi phien ban do (rollback). Dung fs gia de
// tiem migration sai cu phap. Dot bien fail-open (return -> continue) se lam
// test nay DO.
func TestApDungMigration_FailClosed(t *testing.T) {
	s := moStoreTest(t)
	ctx := context.Background()
	_, _ = s.pool.Exec(ctx, `DROP TABLE IF EXISTS _mt_a`)
	t.Cleanup(func() { _, _ = s.pool.Exec(context.Background(), `DROP TABLE IF EXISTS _mt_a`) })

	fsys := fstest.MapFS{
		"migrations/0001_tot.sql":  {Data: []byte(`CREATE TABLE _mt_a (id INT);`)},
		"migrations/0002_hong.sql": {Data: []byte(`DAY LA SQL SAI CU PHAP;`)},
	}
	n, err := s.apDungMigrationTuFS(ctx, fsys)
	if err == nil {
		t.Fatal("migration hong PHAI tra loi (fail-closed), nhung err=nil (fail-open?)")
	}
	if n != 1 {
		t.Fatalf("soApMoi=%d, mong 1 (0001 xong truoc khi 0002 hong)", n)
	}
	var co0001, co0002 bool
	if err := s.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE phien_ban='0001_tot.sql')`).Scan(&co0001); err != nil {
		t.Fatal(err)
	}
	if err := s.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE phien_ban='0002_hong.sql')`).Scan(&co0002); err != nil {
		t.Fatal(err)
	}
	if !co0001 {
		t.Fatal("0001 phai duoc ghi (da commit truoc khi 0002 hong)")
	}
	if co0002 {
		t.Fatal("0002 HONG khong duoc ghi vao schema_migrations (phai rollback)")
	}
}

// TestApDungMigration_TheoThuTu (reviewer F2): 0002 phu thuoc bang 0001 tao ra.
// Ap dung thu tu -> ca hai xong; sai thu tu -> 0002 vo vi bang chua ton tai.
func TestApDungMigration_TheoThuTu(t *testing.T) {
	s := moStoreTest(t)
	ctx := context.Background()
	_, _ = s.pool.Exec(ctx, `DROP TABLE IF EXISTS _mt_b`)
	t.Cleanup(func() { _, _ = s.pool.Exec(context.Background(), `DROP TABLE IF EXISTS _mt_b`) })

	fsys := fstest.MapFS{
		"migrations/0001_tao.sql":  {Data: []byte(`CREATE TABLE _mt_b (id INT);`)},
		"migrations/0002_dung.sql": {Data: []byte(`INSERT INTO _mt_b (id) VALUES (1);`)},
	}
	n, err := s.apDungMigrationTuFS(ctx, fsys)
	if err != nil {
		t.Fatalf("ap theo thu tu that bai: %v (0002 phu thuoc bang cua 0001)", err)
	}
	if n != 2 {
		t.Fatalf("soApMoi=%d, mong 2", n)
	}
}
