package game

import (
	"testing"

	"nroserver/internal/protocol"
)

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
	menu, _ := r.UTF()
	portrait, _ := r.Short()
	if id != 0 || text != "Chào mừng em đến Làng Lá. Hãy vững tin trên con đường ninja của mình." || n != 1 || menu != "Đóng" || portrait != 27078 {
		t.Fatalf("goi 32 Iruka sai: id=%d text=%q n=%d menu=%q portrait=%d", id, text, n, menu, portrait)
	}
	if _, err := r.Byte(); err == nil {
		t.Fatal("goi 32 co byte thua")
	}
}
