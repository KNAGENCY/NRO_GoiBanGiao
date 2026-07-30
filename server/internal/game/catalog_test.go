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
	partOrderChanged := strings.Replace(string(duLieuAssetLock), `2189,2190`, `2190,2189`, 1)
	if _, err := NapCatalog(duLieuCatalog, []byte(partOrderChanged)); err == nil || !strings.Contains(err.Error(), "binding") {
		t.Fatalf("dao thu tu binding part phai fail closed, got %v", err)
	}
}

func TestTaiCollisionKhoaHashVaCanarySolid(t *testing.T) {
	col, err := TaiCollision(duLieuCollision, duLieuAssetLock)
	if err != nil {
		t.Fatalf("TaiCollision() loi: %v", err)
	}
	if !col.DiDuoc(100, 384) || !col.DiDuocDoan(100, 384, 110, 384) {
		t.Fatal("spawn va doan di tren artifact phai walkable")
	}
	mutated := append([]byte(nil), duLieuCollision...)
	mutated[2] = 1
	if _, err := TaiCollision(mutated, duLieuAssetLock); err == nil || !strings.Contains(err.Error(), "ASSET-COLLISION-006") {
		t.Fatalf("canary dao solid phai do ASSET-COLLISION-006, got %v", err)
	}
	blocked := &Collision{width: 2, height: 1, cells: []byte{0, 1}}
	if blocked.DiDuoc(60, 10) || blocked.DiDuocDoan(10, 10, 70, 10) {
		t.Fatal("foot point va segment crossing phai bi chan boi solid cell")
	}
}
