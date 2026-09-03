package game

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Catalog là dữ liệu runtime tường minh: không suy diễn ID từ thứ tự closure.
type Catalog struct {
	Maps []MapCatalog
	NPC  NPCatalog
	Mob  MobCatalog

	// MobIDs la TOAN BO id anh quai ma goi asset DA KY co (asset-lock.json,
	// closure.mob_ids) — khong phai chi con quai dang duoc dung. KH8 task 8
	// doi chieu `anh_id` trong kho-noi-dung/quai.yaml voi danh sach nay
	// (khoquai.go): file YAML sua duoc bang tay, danh sach nay thi khong.
	MobIDs []int
}
type MapCatalog struct {
	ID             byte
	Name           string
	SpawnX, SpawnY int16
}
type NPCatalog struct {
	ID                        byte
	Name                      string
	Head, Body, Leg, Portrait int16
	X, Y                      int16
	Dialog                    string
	Menu                      []string
}
type MobSpawnCatalog struct {
	ID   byte
	X, Y int16
}
type MobCatalog struct {
	ID                 byte
	Name               string
	HP                 int32
	Range, Speed, Dart byte
	Spawns             []MobSpawnCatalog
}

type catalogArtifact struct {
	Maps []struct {
		ID    byte   `json:"id"`
		Name  string `json:"name"`
		Spawn struct {
			X int16 `json:"x"`
			Y int16 `json:"y"`
		} `json:"spawn"`
	} `json:"maps"`
	NPCs []struct {
		ID       byte     `json:"id"`
		Name     string   `json:"name"`
		Head     int16    `json:"head"`
		Body     int16    `json:"body"`
		Leg      int16    `json:"leg"`
		Portrait int16    `json:"portrait"`
		X        int16    `json:"x"`
		Y        int16    `json:"y"`
		Dialog   string   `json:"dialog"`
		Menu     []string `json:"menu"`
	} `json:"npcs"`
	Mobs []struct {
		ID     byte   `json:"id"`
		Name   string `json:"name"`
		HP     int32  `json:"hp"`
		Range  byte   `json:"range"`
		Speed  byte   `json:"speed"`
		Dart   byte   `json:"dart"`
		Spawns []struct {
			ID byte  `json:"id"`
			X  int16 `json:"x"`
			Y  int16 `json:"y"`
		} `json:"spawns"`
	} `json:"mobs"`
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
	var l lockArtifact
	if err := json.Unmarshal(lockData, &l); err != nil {
		return nil, fmt.Errorf("asset lock malformed: %w", err)
	}
	if !lockHashKhop(data, l, "server/catalog.json") {
		return nil, fmt.Errorf("catalog hash khong khop asset lock")
	}
	n, m := a.NPCs[0], a.Mobs[0]
	if n.Name == "" || n.Dialog == "" || len(n.Menu) == 0 || m.Name == "" || m.HP <= 0 || len(m.Spawns) == 0 {
		return nil, fmt.Errorf("catalog co template rong")
	}
	if !coID(l.Closure.PartIDs, int(n.Head)) || !coID(l.Closure.PartIDs, int(n.Body)) || !coID(l.Closure.PartIDs, int(n.Leg)) || !coID(l.Closure.SpriteIDs, int(n.Portrait)) || !coID(l.Closure.MobIDs, int(m.ID)) {
		return nil, fmt.Errorf("asset lock thieu binding Iruka/mob")
	}
	c := &Catalog{NPC: NPCatalog{ID: n.ID, Name: n.Name, Head: n.Head, Body: n.Body, Leg: n.Leg, Portrait: n.Portrait, X: n.X, Y: n.Y, Dialog: n.Dialog, Menu: append([]string(nil), n.Menu...)}, Mob: MobCatalog{ID: m.ID, Name: m.Name, HP: m.HP, Range: m.Range, Speed: m.Speed, Dart: m.Dart}}
	for _, item := range a.Maps {
		if item.Name == "" {
			return nil, fmt.Errorf("catalog map rong")
		}
		c.Maps = append(c.Maps, MapCatalog{ID: item.ID, Name: item.Name, SpawnX: item.Spawn.X, SpawnY: item.Spawn.Y})
	}
	for _, item := range m.Spawns {
		c.Mob.Spawns = append(c.Mob.Spawns, MobSpawnCatalog{ID: item.ID, X: item.X, Y: item.Y})
	}
	c.MobIDs = append([]int(nil), l.Closure.MobIDs...)
	return c, nil
}
func coID(ids []int, want int) bool {
	for _, id := range ids {
		if id == want {
			return true
		}
	}
	return false
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
	width, height, tileWidth, tileHeight int
	cells                                []byte
}

func TaiCollision(data, lockData []byte) (*Collision, error) {
	var l lockArtifact
	if err := json.Unmarshal(lockData, &l); err != nil {
		return nil, fmt.Errorf("asset lock malformed: %w", err)
	}
	if !lockHashKhop(data, l, "server/map-1-collision.bin") {
		return nil, fmt.Errorf("ASSET-COLLISION-006: collision hash khong khop lock")
	}
	if len(data) < 5 {
		return nil, fmt.Errorf("ASSET-COLLISION-006: collision header sai")
	}
	w, h, tw, th := int(data[0]), int(data[1]), int(data[2]), int(data[3])
	if w == 0 || h == 0 || tw != 24 || th != 24 || len(data) != 4+w*h {
		return nil, fmt.Errorf("ASSET-COLLISION-006: collision cell sai")
	}
	hasSolid := false
	for _, cell := range data[4:] {
		if cell > 1 {
			return nil, fmt.Errorf("ASSET-COLLISION-006: semantic collision sai")
		}
		hasSolid = hasSolid || cell == 1
	}
	if !hasSolid {
		return nil, fmt.Errorf("ASSET-COLLISION-006: collision khong co o chan")
	}
	return &Collision{width: w, height: h, tileWidth: tw, tileHeight: th, cells: append([]byte(nil), data[4:]...)}, nil
}
func (c *Collision) DiDuoc(x, y int16) bool {
	if c == nil || x < 0 || y < 0 {
		return false
	}
	col, row := int(x)/c.tileWidth, int(y)/c.tileHeight
	if col >= c.width || row >= c.height {
		return false
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
