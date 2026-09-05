package game

import (
	"testing"

	"nroserver/internal/protocol"
)

// TestGoiMoIrukaByteDung: goi 32 XUONG cua Iruka. Menu tu 05/09/2026 la
// ["Cửa hàng", "Đóng"] (menuNPC, npc.go) — KHONG con la ["Đóng"] cua catalog.json,
// vi catalog bi ghim SHA va menu la noi dung ta quyet.
func TestGoiMoIrukaByteDung(t *testing.T) {
	m, err := goiMoNPC(0)
	if err != nil {
		t.Fatal(err)
	}
	if m.Command != 32 {
		t.Fatalf("lenh = %d, muon 32", m.Command)
	}
	r := protocol.NewReader(m.Data)
	id, _ := r.Short()
	text, _ := r.UTF()
	n, _ := r.Byte()
	if id != 0 || text != "Chào mừng em đến Làng Lá. Hãy vững tin trên con đường ninja của mình." || n != 2 {
		t.Fatalf("goi 32 Iruka sai: id=%d text=%q n=%d", id, text, n)
	}
	dong0, _ := r.UTF()
	dong1, _ := r.UTF()
	portrait, _ := r.Short()
	if dong0 != "Cửa hàng" || dong1 != "Đóng" || portrait != 27078 {
		t.Fatalf("menu/portrait sai: %q %q %d", dong0, dong1, portrait)
	}
	if _, err := r.Byte(); err == nil {
		t.Fatal("goi 32 co byte thua")
	}
}

// TestMenuNPC: dong 0 cua Iruka la "Cửa hàng" (xuLyChonMenuNPC dua vao chi so
// nay); NPC la -> nil.
func TestMenuNPC(t *testing.T) {
	menu := menuNPC(0)
	if len(menu) != 2 || menu[0] != "Cửa hàng" || menu[1] != "Đóng" {
		t.Fatalf("menu Iruka = %v", menu)
	}
	if menuNPC(9) != nil {
		t.Fatal("NPC la phai khong co menu")
	}
	if _, err := goiMoNPC(9); err == nil {
		t.Fatal("goiMoNPC(9) phai loi: NPC khong co trong catalog")
	}
}
