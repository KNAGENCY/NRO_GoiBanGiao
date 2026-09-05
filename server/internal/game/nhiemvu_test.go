package game

import (
	"testing"

	"nroserver/internal/protocol"
)

// TestGoiNhiemVuChinhByteDung doc lai goi 40 dung thu tu client (Controller.cs:3273-3316).
func TestGoiNhiemVuChinhByteDung(t *testing.T) {
	m, err := goiNhiemVuChinh()
	if err != nil {
		t.Fatal(err)
	}
	if m.Command != 40 || lenhNhiemVuChinh != 40 {
		t.Fatalf("lenh = %d, client doc o case 40", m.Command)
	}
	r := protocol.NewReader(m.Data)
	taskId, _ := r.Short()
	index, _ := r.Byte()
	ten, _ := r.UTF()
	chiTiet, _ := r.UTF()
	soMuc, _ := r.Byte()
	if soMuc != 1 {
		t.Fatalf("soMucCon = %d, mong doi 1 (GameScr.tasks[index] phai co phan tu)", soMuc)
	}
	tenMuc, _ := r.UTF()
	npcDich, _ := r.Byte()
	mapDich, _ := r.Short()
	ghiChu, _ := r.UTF()
	count, _ := r.Short()
	dem, _ := r.Short()
	if _, err := r.Byte(); err == nil {
		t.Fatal("goi 40 co byte thua")
	}
	if index != 0 || ten == "" || chiTiet == "" || tenMuc == "" || ghiChu != "" || count != 0 || dem != 0 {
		t.Fatalf("goi 40 sai: index=%d ten=%q chiTiet=%q tenMuc=%q ghiChu=%q count=%d dem=%d",
			index, ten, chiTiet, tenMuc, ghiChu, count, dem)
	}
	if npcDich != 0 || mapDich != 1 {
		t.Fatalf("dich = npc %d map %d, mong doi Iruka (0) o map 1", npcDich, mapDich)
	}
	// Bo loc NPC trong -24 (Controller.cs:5145): taskId >= 7 cho moi mau NPC
	// (tru 6/16) di qua; 0 la nhiem vu mo dau Dragon Ball (re nhanh NPC 3/4).
	if taskId < 7 {
		t.Fatalf("taskId = %d: phai >= 7 de NPC khong bi bo loc -24 chan", taskId)
	}
}
