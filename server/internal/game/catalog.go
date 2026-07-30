package game

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Catalog là bản đọc duy nhất của artifact build. Các ID hiển thị được rút từ
// lock để server không thể tự chạy với một bộ asset khác client.
type Catalog struct {
	Maps []string
	NPC  NPCatalog
	Mob  MobCatalog
}
type NPCatalog struct {
	ID                        byte
	Name                      string
	Head, Body, Leg, Portrait int16
}
type MobCatalog struct {
	ID   byte
	Name string
}

type catalogArtifact struct {
	Maps []string `json:"maps"`
	Mobs []string `json:"mob_templates"`
	NPCs []string `json:"npc_templates"`
}
type lockArtifact struct {
	Closure struct {
		PartIDs   []int `json:"part_ids"`
		SpriteIDs []int `json:"sprite_ids"`
		MobIDs    []int `json:"mob_ids"`
	} `json:"closure"`
	Outputs []struct{ Path, SHA256 string } `json:"outputs"`
}

func NapCatalog(data, lockData []byte) (*Catalog, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("catalog rong")
	}
	var a catalogArtifact
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, fmt.Errorf("catalog malformed: %w", err)
	}
	if len(a.Maps) != 6 || len(a.NPCs) != 1 || len(a.Mobs) != 1 {
		return nil, fmt.Errorf("catalog thieu map/npc/mob")
	}
	if a.NPCs[0] == "" || a.Mobs[0] == "" {
		return nil, fmt.Errorf("catalog co template rong")
	}
	var l lockArtifact
	if err := json.Unmarshal(lockData, &l); err != nil {
		return nil, fmt.Errorf("asset lock malformed: %w", err)
	}
	if !lockHashKhop(data, l, "server/catalog.json") {
		return nil, fmt.Errorf("catalog hash khong khop asset lock")
	}
	parts := l.Closure.PartIDs
	sprites := l.Closure.SpriteIDs
	if len(parts) != 6 || len(sprites) != 42 || len(l.Closure.MobIDs) != 1 || !tangDan(parts) || !tangDan(sprites) {
		return nil, fmt.Errorf("asset lock thieu binding Iruka/mob")
	}
	// Ba part Iruka là ba ID lớn nhất trong nhóm part runtime; portrait là ID
	// đứng ngay trước bốn icon skill ở cuối nhóm sprite closure.
	c := &Catalog{Maps: append([]string(nil), a.Maps...), NPC: NPCatalog{ID: 0, Name: a.NPCs[0], Head: int16(parts[len(parts)-3]), Body: int16(parts[len(parts)-2]), Leg: int16(parts[len(parts)-1]), Portrait: int16(sprites[len(sprites)-5])}, Mob: MobCatalog{ID: byte(l.Closure.MobIDs[0]), Name: a.Mobs[0]}}
	return c, nil
}

func tangDan(ids []int) bool {
	for i := 1; i < len(ids); i++ {
		if ids[i-1] >= ids[i] {
			return false
		}
	}
	return true
}

func lockHashKhop(data []byte, l lockArtifact, path string) bool {
	want := ""
	for _, o := range l.Outputs {
		if o.Path == path {
			want = o.SHA256
			break
		}
	}
	if want == "" {
		return false
	}
	got := sha256.Sum256(data)
	return hex.EncodeToString(got[:]) == want
}

type Collision struct {
	width, height int
	cells         []byte
}

func TaiCollision(data, lockData []byte) (*Collision, error) {
	var l lockArtifact
	if err := json.Unmarshal(lockData, &l); err != nil {
		return nil, fmt.Errorf("asset lock malformed: %w", err)
	}
	if !lockHashKhop(data, l, "server/map-1-collision.bin") {
		return nil, fmt.Errorf("ASSET-COLLISION-006: collision hash khong khop lock")
	}
	if len(data) < 3 {
		return nil, fmt.Errorf("ASSET-COLLISION-006: collision header sai")
	}
	w, h := int(data[0]), int(data[1])
	if w == 0 || h == 0 || len(data) != 2+w*h {
		return nil, fmt.Errorf("ASSET-COLLISION-006: collision cell sai")
	}
	for _, tile := range data[2:] {
		if tile > 99 {
			return nil, fmt.Errorf("ASSET-COLLISION-006: tile collision sai")
		}
	}
	return &Collision{width: w, height: h, cells: append([]byte(nil), data[2:]...)}, nil
}
func (c *Collision) DiDuoc(x, y int16) bool {
	if c == nil || x < 0 || y < 0 {
		return false
	}
	// Toa do gameplay dung o 48 px, con tile art x1 la 24 px; map artifact
	// 5x17 vi vay phu duoc spawn (100,384) va khu choi thuc, khong phai 120px.
	col, row := int(x)/48, int(y)/48
	// Artifact chi mang semantic cua phan map slice. Toa do nam ngoai luoi
	// khong duoc tu y coi la tuong: bien map cu va NguongBuoc van xu ly no.
	if col >= c.width || row >= c.height {
		return true
	}
	return c.cells[row*c.width+col] == 0
}
func (c *Collision) DiDuocDoan(x0, y0, x1, y1 int16) bool {
	dx, dy := int(x1)-int(x0), int(y1)-int(y0)
	n := abs(dx)
	if abs(dy) > n {
		n = abs(dy)
	}
	if n == 0 {
		return c.DiDuoc(x0, y0)
	}
	for i := 0; i <= n; i++ {
		if !c.DiDuoc(int16(int(x0)+dx*i/n), int16(int(y0)+dy*i/n)) {
			return false
		}
	}
	return true
}
func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

var catalogLangLa = phaiNapCatalog()
var collisionLangLa = phaiTaiCollision()

func phaiNapCatalog() *Catalog {
	c, err := NapCatalog(duLieuCatalog, duLieuAssetLock)
	if err != nil {
		panic(err)
	}
	return c
}
func phaiTaiCollision() *Collision {
	c, err := TaiCollision(duLieuCollision, duLieuAssetLock)
	if err != nil {
		panic(err)
	}
	return c
}
