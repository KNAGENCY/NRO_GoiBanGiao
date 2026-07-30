package game

import (
	"strings"
	"testing"
)

func TestNapCatalogArtifactLangLa(t *testing.T) {
	c, err := NapCatalog(duLieuCatalog, duLieuAssetLock)
	if err != nil {
		t.Fatalf("NapCatalog() loi: %v", err)
	}
	if got, want := len(c.Maps), 6; got != want {
		t.Fatalf("so map = %d, muon %d", got, want)
	}
	if c.NPC.Name != "Iruka" || c.Mob.Name != "Sói Rừng" || c.Mob.ID != 0 {
		t.Fatalf("catalog NPC/mob sai: %+v %+v", c.NPC, c.Mob)
	}
	if c.NPC.Portrait != 27078 || c.NPC.Head != 2189 || c.NPC.Body != 2190 || c.NPC.Leg != 2191 {
		t.Fatalf("binding Iruka trong lock sai: %+v", c.NPC)
	}
	if c.NPC.Dialog != "Chào mừng em đến Làng Lá. Hãy vững tin trên con đường ninja của mình." || len(c.NPC.Menu) != 1 || c.NPC.Menu[0] != "Đóng" || c.NPC.X != 172 || c.NPC.Y != 384 {
		t.Fatalf("NPC phai dung fixture catalog: %+v", c.NPC)
	}
	if c.Mob.HP != 100 || len(c.Mob.Spawns) != 2 || c.Mob.Spawns[0].X != 160 || c.Mob.Spawns[1].X != 40 {
		t.Fatalf("mob phai dung spawn catalog: %+v", c.Mob)
	}
}

func TestNapCatalogFailClosedKhiThieuHoacUnknownID(t *testing.T) {
	if _, err := NapCatalog(nil, duLieuAssetLock); err == nil || !strings.Contains(err.Error(), "catalog") {
		t.Fatalf("catalog thieu phai loi ro, got %v", err)
	}
	badLock := strings.Replace(string(duLieuAssetLock), `"mob_ids":[0]`, `"mob_ids":[]`, 1)
	if _, err := NapCatalog(duLieuCatalog, []byte(badLock)); err == nil || !strings.Contains(err.Error(), "binding") {
		t.Fatalf("mob ID unknown phai fail closed, got %v", err)
	}
	withoutPortrait := strings.Replace(string(duLieuAssetLock), `,27078,27079`, `,27079`, 1)
	if _, err := NapCatalog(duLieuCatalog, []byte(withoutPortrait)); err == nil || !strings.Contains(err.Error(), "binding") {
		t.Fatalf("bo portrait lock phai fail closed, got %v", err)
	}
}

func TestTaiCollisionKhoaHashVaCanarySolid(t *testing.T) {
	col, err := TaiCollision(duLieuCollision, duLieuAssetLock)
	if err != nil {
		t.Fatalf("TaiCollision() loi: %v", err)
	}
	if !col.DiDuoc(100, 384) || !col.DiDuocDoan(100, 384, 110, 384) || col.DiDuoc(100, 10) || col.DiDuocDoan(10, 10, 100, 10) {
		t.Fatal("artifact phai co spawn/segment walkable va o solid chan duong")
	}
	mutated := append([]byte(nil), duLieuCollision...)
	// Ô cản production nằm ở cell index 4 (sau header 4 byte).
	mutated[8] = 0
	if _, err := TaiCollision(mutated, duLieuAssetLock); err == nil || !strings.Contains(err.Error(), "ASSET-COLLISION-006") {
		t.Fatalf("canary dao solid phai do ASSET-COLLISION-006, got %v", err)
	}
	blocked := &Collision{width: 2, height: 1, tileWidth: 24, tileHeight: 24, cells: []byte{0, 1}}
	if blocked.DiDuoc(30, 10) || blocked.DiDuocDoan(10, 10, 30, 10) {
		t.Fatal("foot point va segment crossing phai bi chan boi solid cell")
	}
}
