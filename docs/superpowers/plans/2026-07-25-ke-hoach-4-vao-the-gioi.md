# Kế hoạch 4 — Vào thế giới Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Sau khi client gửi `-28/13` (đã xong ở Kế hoạch 3), máy chủ gửi nốt `-82` tileset + `-30/0` ME_INFO + `-24` MAP INFO để client Unity thật **rời màn "Xin chờ" và hiện nhân vật đứng trong một bản đồ làng Naruto**.

**Architecture:** Nối tiếp Kế hoạch 3. Điểm hẹn là `xuLyClientOk` (`login.go`) — nơi Kế hoạch 3 chỉ ghi nhận `clientSanSang`. Kế hoạch 4 gắn thêm: tra bản đồ của nhân vật theo `mapId` (danh mục `BanDo`), dựng một **khu vực (zone) tối thiểu** giữ bản đồ + danh sách người chơi, rồi gửi ba gói `-82` → `-30/0` → `-24` (đúng thứ tự: `-82` bắt buộc trước `-24`). Client nhận đủ ba gói thì `switchToMe()` và render.

**Tech Stack:** Go 1.25+, `//go:embed` (thư viện chuẩn) để nhúng file `tile_set_info`, `testing`. Không thư viện ngoài (chỉ `github.com/jackc/pgx/v5` đã có).

## Tài liệu tham chiếu giao thức (ĐỌC TRƯỚC)

**Mọi bảng byte lấy từ:** [`../specs/2026-07-24-giao-thuc-vao-the-gioi.md`](../specs/2026-07-24-giao-thuc-vao-the-gioi.md).
Kế hoạch này KHÔNG chép lại toàn bộ bảng byte — trích phần cần và trỏ tới tài liệu đó. Các mục liên quan:
§0 (chuỗi tối thiểu), §7 (`-82` tileset), §8 (`-30/0` ME_INFO), §9 (`-24` MAP INFO + byte dự phòng + cạm bẫy), §10 (định dạng file tile map), §11 (ràng buộc cứng). Thêm: `2026-07-23-hoa-giai-muot-va-chong-gian-lan.md` §3.4 (kích thước thế giới), §4.3 (tầm nhìn toàn khu vực).

> **⚠️ Hai điểm ĐÃ ĐỐI CHIẾU LẠI VỚI CODE JAVA THẬT (bổ khuyết cho bảng byte spec):**
>
> 1. **`-30/0` ME_INFO CÓ một byte lệnh con `0` DẪN ĐẦU thân** mà bảng §8 **bỏ sót**. Đã xác minh: `Service.player()` (`Service.java:782`) tạo gói bằng `messageSubCommand((byte) 0)`, và `messageSubCommand` (`Service.java:911-915`) = `new Message(-30)` + `writer().writeByte(command)`. Nên thân ME_INFO **bắt đầu bằng `byte 0`** rồi mới tới `id int` (field 1 của §8). Giống hệt quy ước `-28/X` (byte lệnh con đầu thân). Code dưới đây ghi `w.Byte(0)` trước tiên.
> 2. **`-24` MAP INFO KHÔNG có byte lệnh con** — `Zone.mapInfo()` (`Zone.java:741`) là `new Message(-24)` trực tiếp, thân bắt đầu ngay bằng `mapId`. Byte dự phòng `0` giữa mob và npc đã xác minh tại `Zone.java:805` (`msg.writer().writeByte(0);` sau vòng mob, trước vòng npc). Field order §9 khớp `Zone.java:742-861` từng dòng.

## Môi trường đã dựng sẵn

| Thứ | Vị trí / giá trị |
|---|---|
| Go | `D:\Go\bin\go.exe` (1.26.5) · gofmt `D:\Go\bin\gofmt.exe` |
| PostgreSQL | cổng **5433**, user `postgres` / `nrodev123`, CSDL `nro`/`nro_test` |
| Test CSDL | `$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro_test'` |
| File tileset gốc | `1-MAY-CHU-GAME\data\map\tile_set_info` (**1696 byte**, byte đầu `numTileMap = 39`) |
| Dữ liệu map cũ | `1-MAY-CHU-GAME\z.sql` bảng `map_template` (đã trích giá trị thật cho map 1-5, xem Task 3) |
| Client đích | `D:\NRO_Project` — đã xác nhận `ModFunc.isReadInt=true` |

## Global Constraints

*(Sao chép giá trị chính xác từ spec §11 và bối cảnh thiết kế đã chốt. Mọi task ngầm chịu cả mục này.)*

- **Go 1.25+.** Không thư viện ngoài (chỉ `github.com/jackc/pgx/v5` đã có; task này chỉ dùng thư viện chuẩn `embed`).
- Ghi chú code **tiếng Việt**, giải thích *tại sao*. Mọi file dưới `server/`. **Không đụng `1-MAY-CHU-GAME/`.**
- **Không ghi nhật ký cho gói không xử lý** (lỗi ngập nhật ký đã có tiền lệ). Sự kiện hiếm (vd nhân vật ở mapId không có bản đồ) được log — đó không phải gói rác lặp lại.
- Số nguyên **big-endian**, chuỗi **modified UTF-8** — **dùng lại `protocol.Writer`**, không tự ghi byte.
- Client build đích `D:\NRO_Project` có **`ModFunc.isReadInt=true`** → mọi `hp` ghi **4 byte int** (KH4 map trống nên chưa dùng, nhưng ghi nhớ).
- Mọi `count`/`size`/`tempId` client đọc byte **có dấu** → **≤127**; **trừ** `mapId` và `itemsBody.size` client đọc **unsigned** → ≤255.
- **`gold` trong ME_INFO là `long` 8 byte** (client luôn `readLong()`; ghi 4 byte lệch toàn bộ phần sau).
- **`-82` phải gửi TRƯỚC `-24`** trong cùng phiên; **`map.tileId ∈ [1, numTileMap]` = [1, 39]**. Thiếu `-82`, hoặc gửi sau `-24`, hoặc tileId ngoài khoảng → client `GameMidlet.exit()` (**thoát hẳn app**).
- **`-82` = gửi NGUYÊN file `tile_set_info` (1696 byte)**, không header, không biến đổi — nhúng bằng `//go:embed` như `part.bin` ở KH3.
- `countClass` trong `-28/7` **≥ 5** (đã xong KH3); **`gender`-slot trong ME_INFO = `NhanVat.Lang` (0-4)**, ghi **2 lần**, phải `< countClass`.
- ME_INFO: **3 mảng item (`itemsBody/Bag/Box`) size = 0**; `gem` = `NhanVat.Ngoc`.
- MAP INFO map trống: **0 waypoint / 0 mob / 0 npc / 0 item**; **NHỚ byte dự phòng `0` giữa mob và npc**; **bg-item = `short 0`** (2 byte, KHÔNG bỏ trắng); **effMap = `short 0`**.
- Bản đồ làng **dùng lại mapId cũ** để client Unity hiện tại có sẵn đồ hoạ. **mapId ≤ 255**; `tileId/bgId/planetId/bgType` lấy từ dữ liệu `map_template` map cũ.
- **Nhịp khu vực 100ms chỉ áp cho quái/vật phẩm** (`2026-07-23...` §4.3) — **CHƯA có ở KH4**. **Tầm nhìn toàn khu vực**, không bán kính.
- **Mỗi task một commit.**

## Cấu trúc file

**Tạo mới:**
- `server/internal/game/tileset.bin` — bản sao file `tile_set_info` gốc (1696 byte). Nhúng bằng `//go:embed`.
- `server/internal/game/thegioi.go` — hằng lệnh (`lenhTileSetInfo/lenhMeInfo/lenhMapInfo`), nhúng tileset, và 3 hàm dựng gói vào thế giới (`goiTileSetInfo`, `goiMeInfo`, `goiMapInfo`).
- `server/internal/game/thegioi_test.go` — test byte-exact 3 gói.
- `server/internal/game/banDo.go` — struct `BanDo` (danh mục bản đồ làng) + `danhSachBanDo` + `TimBanDo`/`DanhSachBanDo`.
- `server/internal/game/banDo_test.go`.
- `server/internal/game/khuvuc.go` — struct `KhuVuc` (mô hình khu vực tối thiểu).
- `server/internal/game/khuvuc_test.go`.

**Sửa:**
- `server/internal/protocol/datastream.go` — thêm `Reader.Long() (int64, error)` (nghịch đảo của `Writer.Long`, cần cho test byte-exact ME_INFO đọc 2 trường `long`).
- `server/internal/protocol/datastream_test.go` — test `Reader.Long`.
- `server/internal/game/login.go` — thêm trường `BoXuLy.khuVuc`; `xuLyClientOk` gửi `-82` + `-30/0` + `-24` và dựng khu vực.
- `server/internal/game/login_test.go` — thêm test điều phối `xuLyClientOk` gửi đúng 3 gói.
- `server/internal/game/dongbo_tichhop_test.go` — **mở rộng** `TestVersionDanceDayDu`: đi tiếp qua `-28/13` → nhận ME_INFO + MAP INFO byte-đúng.

## Quyết định thiết kế (đưa vào bản nháp để coordinator duyệt)

### QĐ1 — Mô hình khu vực (zone): YAGNI, chưa dùng actor
Kế hoạch 4 chỉ cần **một người đứng yên trong map trống**: một người, không di chuyển, không broadcast, không quái/vật phẩm. **Chưa dựng mô hình actor (goroutine mỗi zone).** Dùng `KhuVuc` **tối thiểu**: giữ `*BanDo` + slice người chơi (không mutex, không goroutine, không nhịp 100ms).

*Vì sao:* actor chỉ cần khi có **tranh chấp thật** — nhiều luồng cùng đọc/ghi trạng thái zone (di chuyển, thấy nhau, quái đi lại). KH4 chưa có thứ nào trong đó. Dựng actor bây giờ là giải bài toán chưa tồn tại. **Kế hoạch 5 (di chuyển)** sẽ nâng `KhuVuc` thành actor khi có tranh chấp thật, kèm nhịp 100ms cho quái/vật phẩm. `KhuVuc` ở KH4 chỉ là nơi *nạp bản đồ* và làm nền cho KH5.

### QĐ2 — Nội dung bản đồ làng: danh mục `BanDo` dùng lại mapId cũ 1-5
Định nghĩa **5 map khởi đầu** (mỗi làng 0-4 một map) làm dữ liệu Go trong `banDo.go`, dùng lại **mapId cũ 1-5** để client có sẵn đồ hoạ + `tile_map_data`. Các trường render (`tileId/bgId/planetId/bgType/type`) lấy **đúng từ `map_template`** của map cũ 1-5 (xem bảng Task 3) — đảm bảo `tileId ∈ [1,39]` và client vẽ đúng tileset có sẵn. Tên map đặt theo **tên làng Naruto** cho nhập vai (nhãn UI, coordinator có thể tinh chỉnh). Vị trí spawn khớp `ViTriKhoiDau` trong `lang.go` (100, 384 — nằm trên nền đất của các map này).

## Điểm neo giao thức (trích, không chép hết)

**`-82` tileSetInfo (§7)** — `DataGame.sendTileSetInfo` (`DataGame.java:424-434`): `new Message(-82)` + ghi **nguyên file** `tile_set_info`. KHÔNG byte lệnh con. `-82` không nằm trong `lenhDoDaiBaByte` → độ dài 2 byte (1696 < 65535, vừa).

**`-30/0` ME_INFO (§8)** — `Service.player` (`Service.java:776-897`) + `DataGame.sendHeadAvatar` (`:455-464`). Thứ tự thân (đã đối chiếu code, **có byte lệnh con 0 dẫn đầu**):

```
byte 0 (lệnh con) · int id · byte taskMain(0) · byte gender · short head · UTF name
· byte cPK(0) · byte typePk(0) · long power · short buffHp(0) · short buffMp(0)
· byte gender(lần 2) · byte skillCount(0)[+ short×count] · long gold(8 byte!)
· int ruby(0) · int gem · byte itemsBody.size(0) · byte itemsBag.size(0)
· byte itemsBox.size(0) · short headAvatar.size(0) · short conMeo0 · short conMeo1
· short conMeo2 · byte fusion(0) · int deltaTime(0) · byte isNewMember(0/1)
[· short aura · byte effFront · short hat]  ← 3 trường cuối BỎ HẲN (client bắt lỗi êm, §8 note d)
```

**`-24` MAP INFO (§9)** — `Zone.mapInfo` (`Zone.java:738-869`). KHÔNG byte lệnh con. Thứ tự thân (map trống):

```
byte mapId(≤255) · byte planetId · byte tileId(∈[1,39]) · byte bgId · byte type(0)
· UTF mapName · byte zoneId(0) · short x · short y · byte waypointCount(0)
· byte mobCount(0) · byte DỰ_PHÒNG(0)  ← Zone.java:805, giữa mob và npc, ĐỪNG QUÊN
· byte npcCount(0) · byte itemCount(0) · short bgItem(0) · short effMapCount(0)
· byte bgType · byte idSpaceShip(0) · byte isMapDouble(mapId==148 ? 1 : 0)
```

## Ranh giới Kế hoạch 4 ↔ Kế hoạch 5

| Kế hoạch 4 (dừng ở đây) | Kế hoạch 5 |
|---|---|
| Gửi `-82` + `-30/0` + `-24` → client render đứng yên | Di chuyển (`-7`), bình tín dụng, kéo về |
| `KhuVuc` tối thiểu (map + list người chơi) | Nâng `KhuVuc` thành actor, nhịp 100ms cho quái/vật phẩm |
| Một người, tầm nhìn toàn khu vực | Nhiều người, broadcast một-đổi-một |

> **Kiểm thử KH4 bằng client giả** đi trọn version dance rồi **đi tiếp qua `-28/13`** và khẳng định nhận ME_INFO + MAP INFO byte-đúng. Task cuối: chạy tay client Unity thật, xác nhận nhân vật **hiện lên đứng trong làng** (rời "Xin chờ"). KH4 dừng ở "render đứng yên"; di chuyển là KH5.

---

## Task 1: `Reader.Long()` trong protocol

**Vì sao:** test byte-exact ME_INFO phải đọc lại 2 trường `long` (power, gold). `Writer` đã có `Long`, nhưng `Reader` mới chỉ có `Byte/Short/Int/UTF`. Thêm `Reader.Long` — nghịch đảo tự nhiên, tái dùng cho mọi test/parse sau này.

**Files:**
- Modify: `server/internal/protocol/datastream.go`
- Modify: `server/internal/protocol/datastream_test.go`

**Interfaces:**
- Produces: `func (r *Reader) Long() (int64, error)` — đọc 8 byte big-endian thành `int64`.

- [ ] **Step 1: Viết test đỏ**

Thêm vào `server/internal/protocol/datastream_test.go`:

```go
func TestWriterLongReaderLongVongTron(t *testing.T) {
	// Ghi bang Writer.Long roi doc lai bang Reader.Long phai ra dung so ban dau,
	// ke ca so am va so lon (long 64-bit). Big-endian, khop java DataOutputStream.
	cacGiaTri := []int64{0, 1, -1, 250, 1_000_000, 9_223_372_036_854_775_807, -9_223_372_036_854_775_808}
	for _, muon := range cacGiaTri {
		w := NewWriter()
		w.Long(muon)
		r := NewReader(w.Bytes())
		got, err := r.Long()
		if err != nil {
			t.Fatalf("Long(%d): doc loi %v", muon, err)
		}
		if got != muon {
			t.Errorf("Long vong tron: ghi %d, doc %d", muon, got)
		}
		if _, err := r.Byte(); err == nil {
			t.Errorf("Long(%d): con byte thua sau 8 byte, phai het buffer", muon)
		}
	}
}

func TestReaderLongThieuByteBaoLoi(t *testing.T) {
	// Bo dem chi 7 byte -> Long phai bao loi (khong doc lem sang vung khong co).
	r := NewReader([]byte{0, 0, 0, 0, 0, 0, 0})
	if _, err := r.Long(); err == nil {
		t.Error("Long tren bo dem 7 byte phai bao loi thieu byte")
	}
}
```

- [ ] **Step 2: Chạy — phải fail**

Run: `D:\Go\bin\go.exe test ./internal/protocol/ -run "TestWriterLongReaderLong|TestReaderLong" -v`
Expected: FAIL (biên dịch lỗi — `Reader.Long` chưa có).

- [ ] **Step 3: Viết hàm**

Thêm vào `server/internal/protocol/datastream.go` (cạnh `Reader.Int`, dưới nó):

```go
// Long doc 8 byte big-endian thanh int64 — nghich dao cua Writer.Long, khop
// java.io.DataInputStream.readLong. Dung cho cac truong long trong ME_INFO
// (power, gold): client luon readLong() nen ta ghi/doc dung 8 byte.
func (r *Reader) Long() (int64, error) {
	b, err := r.con(8)
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(b)), nil
}
```

(`encoding/binary` đã được import trong `datastream.go`.)

- [ ] **Step 4: Chạy — phải PASS**

Run: `D:\Go\bin\go.exe test ./internal/protocol/ -run "TestWriterLongReaderLong|TestReaderLong" -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add server/internal/protocol/datastream.go server/internal/protocol/datastream_test.go
git commit -m "feat(protocol): Reader.Long doc 8 byte big-endian (cho ME_INFO power/gold)"
```

---

## Task 2: Nhúng `tile_set_info` + gói `-82` tileSetInfo

**Vì sao:** `-82` bắt buộc trước `-24`; thiếu = client thoát app. Gửi nguyên file 1696 byte, nhúng như `part.bin`.

**Files:**
- Create: `server/internal/game/tileset.bin` (sao chép nhị phân từ file gốc)
- Create: `server/internal/game/thegioi.go`
- Create: `server/internal/game/thegioi_test.go`

**Interfaces:**
- Consumes: `protocol.NewWriter`, `Writer.ThemMang` (đã có từ KH3), `protocol.NewMessage`.
- Produces: `var duLieuTileSet []byte`; hằng `lenhTileSetInfo int8 = -82`, `lenhMeInfo int8 = -30`, `lenhMapInfo int8 = -24`; `func goiTileSetInfo() *protocol.Message`.

- [ ] **Step 1: Sao chép file tileset**

Run (PowerShell — sao chép nhị phân nguyên vẹn):
```
Copy-Item "D:\NRO_GoiBanGiao\1-MAY-CHU-GAME\data\map\tile_set_info" "D:\NRO_GoiBanGiao\server\internal\game\tileset.bin"
(Get-Item "D:\NRO_GoiBanGiao\server\internal\game\tileset.bin").Length
```
Expected: `1696`.

- [ ] **Step 2: Viết test đỏ**

Tạo `server/internal/game/thegioi_test.go`:

```go
package game

import (
	"bytes"
	"testing"
)

func TestNhungTileSetDungKichThuoc(t *testing.T) {
	// File tileset nhung phai la file goc 1696 byte, byte dau (numTileMap) = 39.
	// Sai kich thuoc = nhung nham -> tileId ngoai [1,numTileMap] -> client thoat app.
	if len(duLieuTileSet) != 1696 {
		t.Fatalf("duLieuTileSet = %d byte, mong doi 1696 (file tile_set_info goc)", len(duLieuTileSet))
	}
	if duLieuTileSet[0] != 39 {
		t.Errorf("numTileMap (byte dau) = %d, mong doi 39", duLieuTileSet[0])
	}
}

func TestGoiTileSetInfoGuiNguyenFile(t *testing.T) {
	m := goiTileSetInfo()
	if m.Command != -82 {
		t.Fatalf("lenh = %d, mong doi -82", m.Command)
	}
	// Than = NGUYEN file, khong header, khong byte lenh con (§7).
	if !bytes.Equal(m.Data, duLieuTileSet) {
		t.Errorf("than -82 khong bang nguyen file tileset (dai %d, mong doi %d)", len(m.Data), len(duLieuTileSet))
	}
}
```

- [ ] **Step 3: Chạy — phải fail**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run "TestNhungTileSet|TestGoiTileSetInfo" -v`
Expected: FAIL (biên dịch lỗi — `duLieuTileSet`, `goiTileSetInfo` chưa có).

- [ ] **Step 4: Viết `thegioi.go`**

Tạo `server/internal/game/thegioi.go`:

```go
// Package game — cac goi "vao the gioi" (Ke hoach 4): -82 tileset, -30/0 ME_INFO,
// -24 MAP INFO. Gui sau khi client bao san sang (-28/13, xem xuLyClientOk trong
// login.go) de client roi man "Xin cho" va hien nhan vat dung trong lang.
package game

import (
	_ "embed"

	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// duLieuTileSet: NGUYEN file 1-MAY-CHU-GAME/data/map/tile_set_info (1696 byte,
// 39 tileset). Gui nguyen trong goi -82, khong header, khong bien doi (§7). Thieu
// -82, gui SAU -24, hoac map.tileId ngoai [1,numTileMap] -> client GameMidlet.exit()
// (THOAT HAN app). Nhung bang go:embed nhu part.bin o Ke hoach 3.
//
//go:embed tileset.bin
var duLieuTileSet []byte

// Lenh cac goi vao the gioi. Gia tri tu ban Java cu, DUNG DOI:
//
//	-82  DataGame.sendTileSetInfo (DataGame.java:424) — Message(-82) + nguyen file,
//	     KHONG byte lenh con.
//	-30  Service.player (Service.java:782) — messageSubCommand((byte)0) =
//	     Message(-30) + writeByte(0). Byte lenh con 0 la byte DAU than (bang §8 bo sot).
//	-24  Zone.mapInfo (Zone.java:741) — Message(-24) truc tiep, KHONG byte lenh con.
const (
	lenhTileSetInfo int8 = -82
	lenhMeInfo      int8 = -30
	lenhMapInfo     int8 = -24
)

// goiTileSetInfo dung goi -82: gui NGUYEN file tile_set_info. -82 KHONG nam trong
// lenhDoDaiBaByte (protocol/frame.go) nen dung do dai 2 byte — 1696 < 65535, vua.
// Xem tai lieu giao thuc §7.
func goiTileSetInfo() *protocol.Message {
	w := protocol.NewWriter()
	w.ThemMang(duLieuTileSet)
	return protocol.NewMessage(lenhTileSetInfo, w)
}
```

> **Ghi chú:** import `store` được đưa vào sẵn ở đây vì Task 5/6 (`goiMeInfo`/`goiMapInfo`) dùng `store.NhanVat`. Nếu Task 2 chạy độc lập trước khi có Task 5/6, Go sẽ báo `"store" imported and not used`. Để tránh, **Task 2 tạm bỏ dòng import `store`** rồi Task 5 thêm lại khi viết `goiMeInfo`. (Hoặc gộp: nếu thực thi tuần tự sát nhau, giữ import và viết luôn stub — nhưng cách sạch là thêm import đúng lúc dùng.) Bản mã trên là trạng thái **cuối Task 6**; ở Task 2 hãy bỏ `"nroserver/internal/store"` khỏi khối import.

Trạng thái import ở **Task 2**:
```go
import (
	_ "embed"

	"nroserver/internal/protocol"
)
```

- [ ] **Step 5: Chạy — phải PASS**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run "TestNhungTileSet|TestGoiTileSetInfo" -v`
Expected: PASS (2 test).

- [ ] **Step 6: Commit**

```bash
git add server/internal/game/tileset.bin server/internal/game/thegioi.go server/internal/game/thegioi_test.go
git commit -m "feat(game): nhung tile_set_info, dung goi -82 tileSetInfo (nguyen file 1696 byte)"
```

---

## Task 3: Danh mục bản đồ làng (`banDo.go`)

**Vì sao:** MAP INFO cần `mapId/planetId/tileId/bgId/bgType/type/mapName` của bản đồ. Tách thành danh mục Go tra theo `mapId` (như `lang.go`). Giá trị render lấy **đúng từ `map_template`** map cũ 1-5 để client vẽ đồ hoạ có sẵn.

**Bảng giá trị thật (trích `1-MAY-CHU-GAME/z.sql` bảng `map_template`, cột `planet_id, bg_type, tile_id, bg_id`):**

| Làng (gender) | mapId cũ | map_template gốc | planetId | tileId | bgId | bgType | type | spawn X/Y |
|---|---|---|---|---|---|---|---|---|
| 0 Làng Lá | 1 | Đồi hoa cúc | 0 | 1 | 0 | 0 | 0 | 100 / 384 |
| 1 Làng Cát | 2 | Thung lũng tre | 0 | 1 | 0 | 0 | 0 | 100 / 384 |
| 2 Làng Sương Mù | 3 | Rừng nấm | 0 | 2 | 1 | 0 | 0 | 100 / 384 |
| 3 Làng Đá | 4 | Rừng xương | 0 | 2 | 1 | 0 | 0 | 100 / 384 |
| 4 Làng Mây | 5 | Đảo Kamê | 0 | 3 | 2 | 0 | 0 | 100 / 384 |

Mọi `tileId ∈ [1,39]` ✅. Mọi `mapId ≤ 255` ✅. `type=0` (MAP_NORMAL), `is_map_double=0` cho cả 5. Spawn 100/384 khớp `ViTriKhoiDau` trong `lang.go` và nằm trên nền đất (mob/npc các map này ở y≈384-432).

**Files:**
- Create: `server/internal/game/banDo.go`
- Create: `server/internal/game/banDo_test.go`

**Interfaces:**
- Produces: `type BanDo struct{ MapID int32; Ten string; PlanetID, TileID, BgID, BgType, Loai byte; SpawnX, SpawnY int16 }`; `func TimBanDo(mapID int32) *BanDo`; `func DanhSachBanDo() []BanDo`.

- [ ] **Step 1: Viết test đỏ**

Tạo `server/internal/game/banDo_test.go`:

```go
package game

import "testing"

func TestDanhMucBanDoDu5Lang(t *testing.T) {
	ds := DanhSachBanDo()
	if len(ds) != 5 {
		t.Fatalf("co %d ban do, mong doi 5 (mot map khoi dau moi lang)", len(ds))
	}
	// Ghim gia tri DOC LAP voi danhSachBanDo (khong doc tu chinh no) — bat duoc
	// neu ai do sua nham tileId/bgId/planetId; day la gia tri render THAT tu
	// map_template map cu 1-5, sai la client ve sai tileset hoac thoat app.
	type mong struct {
		mapID                          int32
		ten                            string
		planet, tile, bg, bgType, loai byte
	}
	muon := []mong{
		{1, "Làng Lá", 0, 1, 0, 0, 0},
		{2, "Làng Cát", 0, 1, 0, 0, 0},
		{3, "Làng Sương Mù", 0, 2, 1, 0, 0},
		{4, "Làng Đá", 0, 2, 1, 0, 0},
		{5, "Làng Mây", 0, 3, 2, 0, 0},
	}
	for i, w := range muon {
		g := ds[i]
		if g.MapID != w.mapID || g.Ten != w.ten || g.PlanetID != w.planet ||
			g.TileID != w.tile || g.BgID != w.bg || g.BgType != w.bgType || g.Loai != w.loai {
			t.Errorf("ban do %d = %+v, mong doi mapId=%d ten=%q planet=%d tile=%d bg=%d bgType=%d loai=%d",
				i, g, w.mapID, w.ten, w.planet, w.tile, w.bg, w.bgType, w.loai)
		}
	}
}

func TestTileIdMoiBanDoTrongKhoang(t *testing.T) {
	// tileId PHAI thuoc [1, numTileMap] = [1, 39]; ngoai khoang -> client thoat app.
	for _, bd := range DanhSachBanDo() {
		if bd.TileID < 1 || bd.TileID > 39 {
			t.Errorf("ban do mapId=%d co tileId=%d, phai thuoc [1,39]", bd.MapID, bd.TileID)
		}
		if bd.MapID < 0 || bd.MapID > 255 {
			t.Errorf("ban do mapId=%d ngoai [0,255] (client doc unsigned byte)", bd.MapID)
		}
	}
}

func TestTimBanDoTheoMapId(t *testing.T) {
	if bd := TimBanDo(3); bd == nil || bd.Ten != "Làng Sương Mù" || bd.TileID != 2 {
		t.Errorf("TimBanDo(3) = %+v, mong doi Lang Suong Mu tileId 2", bd)
	}
	if bd := TimBanDo(999); bd != nil {
		t.Errorf("TimBanDo(999) = %+v, mong doi nil (khong co)", bd)
	}
	// Tra ban SAO: sua ket qua khong duoc anh huong danh muc goc.
	bd := TimBanDo(1)
	bd.TileID = 99
	if lai := TimBanDo(1); lai.TileID != 1 {
		t.Errorf("TimBanDo tra con tro vao slice goc: sua ban sao lam doi danh muc (tileId=%d)", lai.TileID)
	}
}
```

- [ ] **Step 2: Chạy — phải fail**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run "TestDanhMucBanDo|TestTileId|TestTimBanDo" -v`
Expected: FAIL (`BanDo`, `TimBanDo`... chưa có).

- [ ] **Step 3: Viết `banDo.go`**

Tạo `server/internal/game/banDo.go`:

```go
package game

// BanDo la MOT ban do lang Naruto. Ke hoach 4 DUNG LAI mapId cu (1-5) de client
// Unity hien tai co san do hoa + tile_map_data (hinh anh Naruto se thay khi dung
// lai client). Cac truong render (PlanetID/TileID/BgID/BgType) lay DUNG tu du lieu
// map cu that (bang map_template trong 1-MAY-CHU-GAME/z.sql, map 1-5) de client ve
// dung tileset co san — tileId ngoai [1,numTileMap]=[1,39] lam client THOAT HAN app
// (giao thuc §7).
type BanDo struct {
	MapID    int32  // dung lai mapId cu (client co san). Client doc unsigned byte -> <=255.
	Ten      string // ten hien thi trong MAP INFO (dat theo lang Naruto — nhan UI)
	PlanetID byte
	TileID   byte // ∈ [1,39] cua tile_set_info (§7)
	BgID     byte
	BgType   byte
	Loai     byte  // map type: 0 = MAP_NORMAL
	SpawnX   int16 // vi tri khoi dau mac dinh — hien khop ViTriKhoiDau trong lang.go
	SpawnY   int16 // (nam tren nen dat cua map). Ke hoach 5 dung khi hoi sinh.
}

// danhSachBanDo: 5 map khoi dau, moi lang (gender 0-4) mot map. mapId/tileId/bgId/
// planetId/bgType lay DUNG tu map_template map cu 1-5 (Doi hoa cuc/Thung lung tre/
// Rung nam/Rung xuong/Dao Kame) de client ve dung do hoa co san. Ten map dat theo
// lang Naruto cho nhap vai (chi la nhan UI trong MAP INFO — coordinator co the
// tinh chinh). Ke hoach 6 se chuyen bang nay sang file nap nong (nhu ghi chu o
// lang.go). Gio de trong code cho don gian.
var danhSachBanDo = []BanDo{
	{MapID: 1, Ten: "Làng Lá", PlanetID: 0, TileID: 1, BgID: 0, BgType: 0, Loai: 0, SpawnX: 100, SpawnY: 384},
	{MapID: 2, Ten: "Làng Cát", PlanetID: 0, TileID: 1, BgID: 0, BgType: 0, Loai: 0, SpawnX: 100, SpawnY: 384},
	{MapID: 3, Ten: "Làng Sương Mù", PlanetID: 0, TileID: 2, BgID: 1, BgType: 0, Loai: 0, SpawnX: 100, SpawnY: 384},
	{MapID: 4, Ten: "Làng Đá", PlanetID: 0, TileID: 2, BgID: 1, BgType: 0, Loai: 0, SpawnX: 100, SpawnY: 384},
	{MapID: 5, Ten: "Làng Mây", PlanetID: 0, TileID: 3, BgID: 2, BgType: 0, Loai: 0, SpawnX: 100, SpawnY: 384},
}

// TimBanDo tra ban do theo mapId, hoac nil neu khong co. Tra ban SAO (cung ly do
// nhu TimLang trong lang.go: nguoi goi lo tay sua khong dung toi danh muc he thong).
func TimBanDo(mapID int32) *BanDo {
	for i := range danhSachBanDo {
		if danhSachBanDo[i].MapID == mapID {
			cp := danhSachBanDo[i]
			return &cp
		}
	}
	return nil
}

// DanhSachBanDo tra ban sao danh sach ban do (chi doc).
func DanhSachBanDo() []BanDo {
	ds := make([]BanDo, len(danhSachBanDo))
	copy(ds, danhSachBanDo)
	return ds
}
```

- [ ] **Step 4: Chạy — phải PASS**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run "TestDanhMucBanDo|TestTileId|TestTimBanDo" -v`
Expected: PASS (3 test).

- [ ] **Step 5: Commit**

```bash
git add server/internal/game/banDo.go server/internal/game/banDo_test.go
git commit -m "feat(game): danh muc BanDo 5 lang (dung lai mapId cu 1-5, gia tri map_template that)"
```

---

## Task 4: Mô hình khu vực tối thiểu (`khuvuc.go`)

**Vì sao:** QĐ1 — KH4 cần một nơi *nạp bản đồ* + giữ danh sách người chơi, làm nền cho KH5. Chưa dùng actor (YAGNI).

**Files:**
- Create: `server/internal/game/khuvuc.go`
- Create: `server/internal/game/khuvuc_test.go`

**Interfaces:**
- Consumes: `BanDo` (Task 3), `store.NhanVat`.
- Produces: `type KhuVuc`; `func MoKhuVuc(bd *BanDo) *KhuVuc`; `func (k *KhuVuc) BanDo() *BanDo`; `func (k *KhuVuc) ThemNguoiChoi(nv *store.NhanVat)`; `func (k *KhuVuc) SoNguoiChoi() int`.

- [ ] **Step 1: Viết test đỏ**

Tạo `server/internal/game/khuvuc_test.go`:

```go
package game

import (
	"testing"

	"nroserver/internal/store"
)

func TestMoKhuVucGiuBanDo(t *testing.T) {
	bd := TimBanDo(1)
	k := MoKhuVuc(bd)
	if k.BanDo() != bd {
		t.Error("MoKhuVuc phai giu dung con tro ban do da cho")
	}
	if k.SoNguoiChoi() != 0 {
		t.Errorf("khu vuc moi mo co %d nguoi, mong doi 0", k.SoNguoiChoi())
	}
}

func TestThemNguoiChoiVaoKhuVuc(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.ThemNguoiChoi(&store.NhanVat{ID: 1, Ten: "Naruto"})
	if k.SoNguoiChoi() != 1 {
		t.Fatalf("sau khi them 1 nguoi, SoNguoiChoi = %d, mong doi 1", k.SoNguoiChoi())
	}
	k.ThemNguoiChoi(&store.NhanVat{ID: 2, Ten: "Sasuke"})
	if k.SoNguoiChoi() != 2 {
		t.Errorf("sau khi them 2 nguoi, SoNguoiChoi = %d, mong doi 2", k.SoNguoiChoi())
	}
}
```

- [ ] **Step 2: Chạy — phải fail**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run "TestMoKhuVuc|TestThemNguoiChoi" -v`
Expected: FAIL (`KhuVuc`, `MoKhuVuc`... chưa có).

- [ ] **Step 3: Viết `khuvuc.go`**

Tạo `server/internal/game/khuvuc.go`:

```go
package game

import "nroserver/internal/store"

// KhuVuc (zone) la don vi khong gian nho nhat cua mot ban do: giu ban do dang
// choi + danh sach nguoi choi dang o trong do. Ke hoach 4 CHI can mot nguoi dung
// yen trong map trong, nen day la mo hinh TOI THIEU — khong goroutine rieng,
// khong mutex, khong nhip 100ms.
//
// VI SAO CHUA DUNG MO HINH ACTOR (goroutine moi zone) NGAY: actor chi can khi co
// TRANH CHAP that — nhieu luong cung doc/ghi trang thai zone (di chuyen, nhin thay
// nhau, quai di lai). Ke hoach 4 chua co thu nao trong so do: mot nguoi, dung yen,
// khong broadcast. Dung actor bay gio la giai bai toan chua ton tai (YAGNI). Ke
// hoach 5 (di chuyen) se nang KhuVuc thanh actor khi co tranh chap that su, kem
// nhip 100ms cho quai/vat pham (xem 2026-07-23-hoa-giai-muot §4.3). O KH4, KhuVuc
// chi lo NAP ban do va lam nen cho KH5.
type KhuVuc struct {
	banDo     *BanDo
	nguoiChoi []*store.NhanVat
}

// MoKhuVuc tao mot khu vuc giu ban do da cho, chua co nguoi choi nao.
func MoKhuVuc(bd *BanDo) *KhuVuc {
	return &KhuVuc{banDo: bd}
}

// BanDo tra ban do cua khu vuc (client can no de dung MAP INFO).
func (k *KhuVuc) BanDo() *BanDo { return k.banDo }

// ThemNguoiChoi dua mot nhan vat vao khu vuc. Ke hoach 5 se broadcast su xuat
// hien nay cho nguoi khac; Ke hoach 4 chi giu danh sach de lam nen.
func (k *KhuVuc) ThemNguoiChoi(nv *store.NhanVat) {
	k.nguoiChoi = append(k.nguoiChoi, nv)
}

// SoNguoiChoi tra so nguoi dang trong khu vuc (dung cho kiem thu + Ke hoach 5).
func (k *KhuVuc) SoNguoiChoi() int { return len(k.nguoiChoi) }
```

- [ ] **Step 4: Chạy — phải PASS**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run "TestMoKhuVuc|TestThemNguoiChoi" -v`
Expected: PASS (2 test).

- [ ] **Step 5: Commit**

```bash
git add server/internal/game/khuvuc.go server/internal/game/khuvuc_test.go
git commit -m "feat(game): KhuVuc toi thieu (giu ban do + danh sach nguoi choi; actor de KH5)"
```

---

## Task 5: Gói `-30/0` ME_INFO (`goiMeInfo`)

**Vì sao:** ME_INFO là gói BẮT BUỘC để client dựng nhân vật của chính mình. Byte-exact theo §8 + phát hiện byte lệnh con `0` dẫn đầu.

**Files:**
- Modify: `server/internal/game/thegioi.go`
- Modify: `server/internal/game/thegioi_test.go`

**Interfaces:**
- Consumes: `store.NhanVat` (các trường `ID/Ten/Lang/KieuToc/SucManh/Vang/Ngoc`), `protocol.Writer`, `Reader.Long` (Task 1).
- Produces: `func goiMeInfo(nv *store.NhanVat) (*protocol.Message, error)` — trả `error` vì `Writer.UTF` (tên) có thể lỗi.

- [ ] **Step 1: Viết test đỏ (byte-exact)**

Thêm vào `server/internal/game/thegioi_test.go` (thêm `import "nroserver/internal/protocol"` và `import "nroserver/internal/store"` vào khối import):

```go
func TestGoiMeInfoByteDung(t *testing.T) {
	nv := &store.NhanVat{
		ID:      12345,
		Ten:     "Naruto",
		Lang:    3, // lang Da -> gender-slot 3, phai < countClass (5)
		KieuToc: 7,
		SucManh: 1_000_000,
		Vang:    250,
		Ngoc:    9,
		MapID:   4,
		ViTriX:  100,
		ViTriY:  384,
	}
	m, err := goiMeInfo(nv)
	if err != nil {
		t.Fatalf("goiMeInfo loi: %v", err)
	}
	if m.Command != -30 {
		t.Fatalf("lenh = %d, mong doi -30", m.Command)
	}

	r := protocol.NewReader(m.Data)
	b := func(ten string, muon byte) {
		g, err := r.Byte()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}
	sh := func(ten string, muon int16) {
		g, err := r.Short()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}
	in := func(ten string, muon int32) {
		g, err := r.Int()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}
	lo := func(ten string, muon int64) {
		g, err := r.Long()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}

	b("byte lenh con", 0) // §8 bo sot — messageSubCommand((byte)0)
	in("id", 12345)
	b("taskMain.id", 0)
	b("gender", 3) // = NhanVat.Lang
	sh("head", 7)  // = NhanVat.KieuToc
	if ten, err := r.UTF(); err != nil || ten != "Naruto" {
		t.Fatalf("name = %q (loi %v), mong doi \"Naruto\"", ten, err)
	}
	b("cPK", 0)
	b("typePk", 0)
	lo("power", 1_000_000) // = NhanVat.SucManh
	sh("buffHp", 0)
	sh("buffMp", 0)
	b("gender lan 2", 3) // = NhanVat.Lang, chi so nClasss[]
	b("skillCount", 0)
	lo("gold", 250) // LONG 8 byte = NhanVat.Vang
	in("ruby", 0)
	in("gem", 9) // = NhanVat.Ngoc
	b("itemsBody.size", 0)
	b("itemsBag.size", 0)
	b("itemsBox.size", 0)
	sh("headAvatar.size", 0)
	sh("conMeo0", 0)
	sh("conMeo1", 0)
	sh("conMeo2", 0)
	b("fusion", 0)
	in("deltaTime", 0)
	b("isNewMember", 0)

	// Bo 3 truong cuoi (aura/effFront/hat) -> phai HET buffer o day.
	if _, err := r.Byte(); err == nil {
		t.Error("ME_INFO con byte thua sau isNewMember — 3 truong cuoi (aura/effFront/hat) phai bo han (§8 note d)")
	}
}
```

- [ ] **Step 2: Chạy — phải fail**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run TestGoiMeInfoByteDung -v`
Expected: FAIL (`goiMeInfo` chưa có).

- [ ] **Step 3: Viết `goiMeInfo`**

Trong `server/internal/game/thegioi.go`: đảm bảo khối import có `"nroserver/internal/store"` (Task 2 đã tạm bỏ — thêm lại đây), rồi thêm hàm:

```go
// goiMeInfo dung goi -30/0 ME_INFO: du lieu nhan vat cua chinh nguoi choi de client
// dung "nhan vat cua toi" va set statusMe. Byte-exact theo tai lieu §8 (da doi chieu
// Service.player, Service.java:776-897). LUU Y bay cam:
//
//	(a) THAN BAT DAU BANG byte lenh con 0 (messageSubCommand((byte)0)) — bang §8 bo
//	    sot byte nay; giong quy uoc -28/X.
//	(b) gold la LONG 8 BYTE (client luon readLong) — ghi 4 byte lech toan bo phan sau.
//	(c) gender ghi 2 LAN, deu = NhanVat.Lang (0-4); lan 2 la chi so nClasss[], phai
//	    < countClass (=5, xem -28/7). Sai -> IndexOutOfRange bi nuot -> ket "Xin cho".
//	(d) 3 mang item (Body/Bag/Box) size = 0 cho nhan vat moi — item id la khong co
//	    trong catalog -28/8 se NPE giua goi.
//	(e) 3 truong cuoi (aura short, effFront byte, hat short) BO HAN — client bat loi em.
func goiMeInfo(nv *store.NhanVat) (*protocol.Message, error) {
	w := protocol.NewWriter()
	w.Byte(0) // (a) byte lenh con 0

	w.Int(int32(nv.ID))
	w.Byte(0)              // taskMain.id — chua co nhiem vu
	w.Byte(byte(nv.Lang))  // (c) gender = id lang
	w.Short(nv.KieuToc)    // head
	if err := w.UTF(nv.Ten); err != nil {
		return nil, err
	}
	w.Byte(0)          // cPK
	w.Byte(0)          // typePk
	w.Long(nv.SucManh) // power (long)
	w.Short(0)         // buffHp
	w.Short(0)         // buffMp
	w.Byte(byte(nv.Lang)) // (c) gender lan 2 — chi so nClasss[], phai < countClass
	w.Byte(0)          // skillCount = 0 (khong ghi short skillId nao)

	w.Long(nv.Vang)         // (b) gold LONG 8 byte
	w.Int(0)                // ruby
	w.Int(int32(nv.Ngoc))   // gem

	w.Byte(0) // (d) itemsBody.size
	w.Byte(0) // (d) itemsBag.size
	w.Byte(0) // (d) itemsBox.size

	w.Short(0) // headAvatar.size (DataGame.sendHeadAvatar rong)
	w.Short(0) // conMeo[gender][0]
	w.Short(0) // conMeo[gender][1]
	w.Short(0) // conMeo[gender][2]

	w.Byte(0) // fusion (khong nhap the)
	w.Int(0)  // deltaTime
	w.Byte(0) // isNewMember — store chua co truong nay, ghi 0 an toan (§8 field 23)

	// (e) BO 3 truong cuoi (aura/effFront/hat).
	return protocol.NewMessage(lenhMeInfo, w), nil
}
```

> **Lưu ý kiểu:** `nv.KieuToc` là `int16` → `w.Short` nhận `int16` trực tiếp. `nv.ID/Ngoc` là `int64` → ép `int32` (id nhân vật + ngọc trong khoảng int32). `nv.SucManh/Vang` là `int64` → `w.Long`. `nv.Lang` là `int16` (0-4) → `byte`.

- [ ] **Step 4: Chạy — phải PASS**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run TestGoiMeInfoByteDung -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add server/internal/game/thegioi.go server/internal/game/thegioi_test.go
git commit -m "feat(game): dung goi -30/0 ME_INFO (byte lenh con 0 dan dau, gold long 8 byte)"
```

---

## Task 6: Gói `-24` MAP INFO (`goiMapInfo`)

**Vì sao:** MAP INFO là gói kích `switchToMe()` — client rời "Xin chờ". Byte-exact theo §9, đặc biệt **byte dự phòng** giữa mob và npc.

**Files:**
- Modify: `server/internal/game/thegioi.go`
- Modify: `server/internal/game/thegioi_test.go`

**Interfaces:**
- Consumes: `store.NhanVat` (`ViTriX/ViTriY`), `BanDo` (Task 3), `protocol.Writer`.
- Produces: `func goiMapInfo(nv *store.NhanVat, bd *BanDo) (*protocol.Message, error)`.

- [ ] **Step 1: Viết test đỏ (byte-exact)**

Thêm vào `server/internal/game/thegioi_test.go`:

```go
func TestGoiMapInfoByteDung(t *testing.T) {
	nv := &store.NhanVat{MapID: 4, ViTriX: 100, ViTriY: 384}
	bd := TimBanDo(4) // Lang Da: planet 0, tile 2, bg 1, bgType 0, loai 0, ten "Làng Đá"
	m, err := goiMapInfo(nv, bd)
	if err != nil {
		t.Fatalf("goiMapInfo loi: %v", err)
	}
	if m.Command != -24 {
		t.Fatalf("lenh = %d, mong doi -24", m.Command)
	}

	r := protocol.NewReader(m.Data)
	b := func(ten string, muon byte) {
		g, err := r.Byte()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}
	sh := func(ten string, muon int16) {
		g, err := r.Short()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}

	b("mapId", 4)
	b("planetId", 0)
	b("tileId", 2) // ∈ [1,39]
	b("bgId", 1)
	b("type", 0) // MAP_NORMAL
	if ten, err := r.UTF(); err != nil || ten != "Làng Đá" {
		t.Fatalf("mapName = %q (loi %v), mong doi \"Làng Đá\"", ten, err)
	}
	b("zoneId", 0)
	sh("x", 100) // = NhanVat.ViTriX
	sh("y", 384) // = NhanVat.ViTriY
	b("waypointCount", 0)
	b("mobCount", 0)
	b("byte du phong", 0) // Zone.java:805 — giua mob va npc, DUNG QUEN
	b("npcCount", 0)
	b("itemCount", 0)
	sh("bgItem", 0)  // short 0 (KHONG bo trang)
	sh("effMap", 0)  // short 0
	b("bgType", 0)
	b("idSpaceShip", 0)
	b("isMapDouble", 0) // mapId != 148

	if _, err := r.Byte(); err == nil {
		t.Error("MAP INFO con byte thua sau isMapDouble — phai het buffer")
	}
}

func TestGoiMapInfoDungViTriNhanVat(t *testing.T) {
	// x/y trong MAP INFO lay tu VI TRI NHAN VAT (nv.ViTriX/Y), khong phai spawn
	// mac dinh cua ban do — de KH5 (di chuyen) van dung sau khi nguoi choi da dich.
	nv := &store.NhanVat{MapID: 1, ViTriX: 512, ViTriY: 240}
	m, err := goiMapInfo(nv, TimBanDo(1))
	if err != nil {
		t.Fatal(err)
	}
	r := protocol.NewReader(m.Data)
	for i := 0; i < 5; i++ {
		r.Byte() // mapId,planet,tile,bg,type
	}
	r.UTF()   // mapName
	r.Byte()  // zoneId
	x, _ := r.Short()
	y, _ := r.Short()
	if x != 512 || y != 240 {
		t.Errorf("x/y = %d/%d, mong doi 512/240 (vi tri nhan vat, khong phai spawn ban do)", x, y)
	}
}
```

- [ ] **Step 2: Chạy — phải fail**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run TestGoiMapInfo -v`
Expected: FAIL (`goiMapInfo` chưa có).

- [ ] **Step 3: Viết `goiMapInfo`**

Thêm vào `server/internal/game/thegioi.go`:

```go
// goiMapInfo dung goi -24 MAP INFO: gio kich switchToMe() ben client (roi man
// "Xin cho"). Byte-exact theo tai lieu §9 (da doi chieu Zone.mapInfo, Zone.java:
// 738-869). Map TRONG: 0 waypoint/mob/npc/item. LUU Y bay cam:
//
//	(a) BYTE DU PHONG 0 giua mob va npc (Zone.java:805) — QUEN la lech toan bo phan sau.
//	(b) bg-item RONG van phai short 0 (2 byte), KHONG bo trang.
//	(c) x/y lay tu VI TRI NHAN VAT (nv.ViTriX/Y) — Zone dung pl.location, khong phai
//	    spawn ban do; giu dung de KH5 van dung sau khi nguoi choi da di chuyen.
//	(d) -82 PHAI da duoc gui truoc goi nay (xuLyClientOk dam bao thu tu), va tileId
//	    phai ∈ [1,numTileMap] — neu khong client THOAT HAN app.
func goiMapInfo(nv *store.NhanVat, bd *BanDo) (*protocol.Message, error) {
	w := protocol.NewWriter()
	w.Byte(byte(bd.MapID))    // mapId (client doc unsigned byte, <=255)
	w.Byte(bd.PlanetID)       // planetId
	w.Byte(bd.TileID)         // tileId ∈ [1,39]
	w.Byte(bd.BgID)           // bgId
	w.Byte(bd.Loai)           // type = 0 MAP_NORMAL
	if err := w.UTF(bd.Ten); err != nil {
		return nil, err
	}
	w.Byte(0) // zoneId = 0 (map mot zone)

	w.Short(int16(nv.ViTriX)) // (c) x = vi tri nhan vat
	w.Short(int16(nv.ViTriY)) // (c) y

	w.Byte(0) // waypointCount = 0
	w.Byte(0) // mobCount = 0
	w.Byte(0) // (a) BYTE DU PHONG — giua mob va npc, DUNG QUEN
	w.Byte(0) // npcCount = 0
	w.Byte(0) // itemCount = 0

	w.Short(0) // (b) bg-item: short 0 (map trong, khong bo trang)
	w.Short(0) // effMapCount = 0

	w.Byte(bd.BgType) // bgType
	w.Byte(0)         // idSpaceShip = 0
	// isMapDouble: rieng mapId 148 la map doi (Zone.java:861). Lang 1-5 khong phai
	// 148 nen luon 0, nhung tinh theo mapId de dung ca khi them map 148 sau nay.
	if bd.MapID == 148 {
		w.Byte(1)
	} else {
		w.Byte(0)
	}
	return protocol.NewMessage(lenhMapInfo, w), nil
}
```

- [ ] **Step 4: Chạy — phải PASS**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run TestGoiMapInfo -v`
Expected: PASS (2 test).

- [ ] **Step 5: Commit**

```bash
git add server/internal/game/thegioi.go server/internal/game/thegioi_test.go
git commit -m "feat(game): dung goi -24 MAP INFO (byte du phong, bg-item/effMap short 0)"
```

---

## Task 7: Điều phối — `xuLyClientOk` gửi `-82` → `-30/0` → `-24`

**Vì sao:** đây là điểm hẹn KH4. Client gửi `-28/13` → máy chủ tra bản đồ, dựng khu vực, gửi ba gói (đúng thứ tự `-82` trước `-24`) → client render.

**Files:**
- Modify: `server/internal/game/login.go`
- Modify: `server/internal/game/login_test.go`

**Interfaces:**
- Consumes: `TimBanDo` (Task 3), `MoKhuVuc`/`KhuVuc` (Task 4), `goiTileSetInfo` (Task 2), `goiMeInfo` (Task 5), `goiMapInfo` (Task 6); `s.Gui`, `b.nhanVat`.
- Produces: trường `BoXuLy.khuVuc *KhuVuc`; `xuLyClientOk` gửi 3 gói vào thế giới.

- [ ] **Step 1: Viết test đỏ (điều phối, dùng `khoGia`)**

Thêm vào `server/internal/game/login_test.go`:

```go
func TestClientOkGuiVaoTheGioi(t *testing.T) {
	// Client gui -28/13 khi da vao the gioi -> may chu gui DUNG 3 goi theo thu tu
	// -82 (tileset) -> -30 (ME_INFO) -> -24 (MAP INFO). -82 PHAI truoc -24 (§7).
	b := NewBoXuLy(&khoGia{})
	b.nhanVat = &store.NhanVat{ID: 7, Ten: "Naruto", Lang: 0, KieuToc: 5,
		SucManh: 100, Vang: 1000, Ngoc: 0, MapID: 1, ViTriX: 100, ViTriY: 384}
	fClient, dong := dungPhienThat(t, b)
	defer dong()

	// Gui -28/13 (client ok).
	w := protocol.NewWriter()
	w.Byte(lenhConClientOk) // 13
	if err := fClient.WriteMessage(protocol.NewMessage(lenhKhongTrenBanDo, w)); err != nil {
		t.Fatal(err)
	}

	m82 := docMotGoi(t, fClient)
	if m82.Command != -82 {
		t.Fatalf("goi 1 = lenh %d, mong doi -82 (tileset, BAT BUOC truoc -24)", m82.Command)
	}
	m30 := docMotGoi(t, fClient)
	if m30.Command != -30 {
		t.Fatalf("goi 2 = lenh %d, mong doi -30 (ME_INFO)", m30.Command)
	}
	m24 := docMotGoi(t, fClient)
	if m24.Command != -24 {
		t.Fatalf("goi 3 = lenh %d, mong doi -24 (MAP INFO)", m24.Command)
	}
	// MAP INFO byte dau (sau khong co lenh con) = mapId = 1.
	if len(m24.Data) == 0 || m24.Data[0] != 1 {
		t.Errorf("MAP INFO mapId = %v, mong doi 1 (mapId cua nhan vat)", m24.Data)
	}
}

func TestClientOkChuaVaoTheGioiKhongGui(t *testing.T) {
	// Chua co nhan vat (b.nhanVat == nil): -28/13 bi bo qua im lang, khong gui gi,
	// khong sap. Client that khong bao gio lam vay; chi ke gia lap moi gui.
	b := NewBoXuLy(&khoGia{}) // b.nhanVat = nil
	fClient, dong := dungPhienThat(t, b)
	defer dong()

	w := protocol.NewWriter()
	w.Byte(lenhConClientOk)
	if err := fClient.WriteMessage(protocol.NewMessage(lenhKhongTrenBanDo, w)); err != nil {
		t.Fatal(err)
	}
	// Khong duoc gui goi nao: dat han doc ngan, mong doi timeout.
	_ = fClientHetHanNgan(t, fClient)
}
```

> **Helper `fClientHetHanNgan`:** nếu `login_test.go` chưa có cách khẳng định "không gói nào tới" trên `*protocol.Framer` (khác `clientGia` giữ `net.Conn`), thêm helper nhỏ đọc một gói với kỳ vọng lỗi/timeout. Vì `dungPhienThat` đặt sẵn read deadline 5s trên `benClient` nhưng KHÔNG trả ra `net.Conn`, cách gọn nhất là **đọc một gói và khẳng định `ReadMessage` trả lỗi** (kết nối rảnh → hết hạn 5s). Chấp nhận chờ tối đa 5s cho một test này, hoặc bỏ `TestClientOkChuaVaoTheGioiKhongGui` nếu muốn tránh chờ (nhánh `b.nhanVat == nil` đã được `xuLyClientOk` giữ nguyên từ KH3 và có test bao phủ gián tiếp). **Quyết định gợi ý:** giữ test, viết helper:

```go
// fClientHetHanNgan khang dinh KHONG con goi nao may chu gui ra: doc mot goi va
// mong doi loi (ket noi ranh -> ReadMessage het han theo deadline da dat trong
// dungPhienThat). Tra ve true neu dung la khong co goi.
func fClientHetHanNgan(t *testing.T, f *protocol.Framer) bool {
	t.Helper()
	if _, err := f.ReadMessage(); err == nil {
		t.Error("nhan duoc goi NGOAI DU KIEN khi mong doi im lang")
		return false
	}
	return true
}
```

- [ ] **Step 2: Chạy — phải fail**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run "TestClientOkGuiVaoTheGioi|TestClientOkChuaVaoTheGioi" -v`
Expected: FAIL (`xuLyClientOk` chưa gửi 3 gói; `BoXuLy.khuVuc` chưa có).

- [ ] **Step 3: Thêm trường `khuVuc` vào `BoXuLy`**

Trong `server/internal/game/login.go`, thêm trường (cạnh `clientSanSang`):

```go
	clientSanSang bool    // client da gui -28/13 (bat tay xong)
	khuVuc        *KhuVuc // khu vuc nhan vat dang o (dung o KH4; nen mong cho KH5)
```

- [ ] **Step 4: Sửa `xuLyClientOk` gửi 3 gói vào thế giới**

Thay thân `xuLyClientOk` (giữ nguyên guard `b.nhanVat == nil`):

```go
// xuLyClientOk: client bao da nhan du 4 catalog va san sang (-28/13). Ke hoach 4
// gui NOT 3 goi de client render: -82 tileset -> -30/0 ME_INFO -> -24 MAP INFO.
// Thu tu quan trong: -82 PHAI truoc -24 (§7), khong thi client thoat han app.
func (b *BoXuLy) xuLyClientOk(s *mynet.Session) error {
	if b.nhanVat == nil {
		return nil
	}
	b.clientSanSang = true

	bd := TimBanDo(b.nhanVat.MapID)
	if bd == nil {
		// Nhan vat o mapId khong co trong danh muc ban do: khong the gui MAP INFO
		// hop le. Ghi log (SU KIEN HIEM — khong phai goi rac lap lai, nen log o day
		// KHONG vi pham luat chong ngap nhat ky) va dung. Client se ket "Xin cho"
		// nhung KHONG sap — tot hon la gui -24 rac lam client thoat app.
		slog.Error("nhan vat o mapId khong co ban do trong danh muc",
			"mapId", b.nhanVat.MapID, "ten", b.nhanVat.Ten, "id", b.nhanVat.ID)
		return nil
	}

	// Dung khu vuc toi thieu cho phien nay va dua nhan vat vao (nen mong cho KH5).
	b.khuVuc = MoKhuVuc(bd)
	b.khuVuc.ThemNguoiChoi(b.nhanVat)

	// -82 tileset — BAT BUOC truoc -24.
	if err := s.Gui(goiTileSetInfo()); err != nil {
		return err
	}

	// -30/0 ME_INFO.
	meInfo, err := goiMeInfo(b.nhanVat)
	if err != nil {
		return err
	}
	if err := s.Gui(meInfo); err != nil {
		return err
	}

	// -24 MAP INFO — gio kich switchToMe(), client roi "Xin cho".
	mapInfo, err := goiMapInfo(b.nhanVat, bd)
	if err != nil {
		return err
	}
	if err := s.Gui(mapInfo); err != nil {
		return err
	}

	slog.Info("da gui the gioi (tileset + ME_INFO + MAP INFO), client se render",
		"ten", b.nhanVat.Ten, "mapId", bd.MapID, "lang", b.nhanVat.Lang, "id", b.nhanVat.ID)
	return nil
}
```

- [ ] **Step 5: Chạy toàn package game — phải PASS**

Run (trong `server/`):
```
$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro_test'
D:\Go\bin\go.exe test ./internal/game/ -run "TestClientOk" -v
D:\Go\bin\go.exe test ./internal/game/
```
Expected: các test `TestClientOk*` PASS; toàn package PASS (test cũ KH1-3 vẫn xanh). `TestVersionDanceDayDu` cũ sẽ ĐỎ ở bước sau — nó khẳng định "sau `-28/13` máy chủ chỉ phục vụ lại skill" (hành vi KH3), giờ máy chủ gửi `-82/-30/-24`. Task 8 sửa test đó. Nếu chạy `TestVersionDanceDayDu` ở bước này thấy đỏ là **đúng dự kiến**.

- [ ] **Step 6: Commit**

```bash
git add server/internal/game/login.go server/internal/game/login_test.go
git commit -m "feat(game): xuLyClientOk gui -82+-30/0+-24 (vao the gioi), dung khu vuc toi thieu"
```

---

## Task 8: Mở rộng `TestVersionDanceDayDu` — đi tiếp qua `-28/13` → nhận ME_INFO + MAP INFO

**Vì sao:** test tích hợp đầu-cuối phải chứng minh client giả đi trọn version dance **rồi nhận đủ 3 gói vào thế giới byte-đúng** qua mạng thật + CSDL thật.

**Files:**
- Modify: `server/internal/game/dongbo_tichhop_test.go`

**Interfaces:**
- Consumes: `doiGoiLenh(t, lenh)` (đã có — khớp theo lệnh, dùng được cho `-82/-30/-24` vì chúng dùng độ dài 2 byte), `khoThat.LayNhanVatTheoTaiKhoan`, `TimBanDo`, `duLieuTileSet`, `Reader.Long`.

- [ ] **Step 1: Thêm hàm khẳng định byte cho 3 gói vào thế giới**

Thêm vào `server/internal/game/dongbo_tichhop_test.go` (thêm `import "context"` và `import "nroserver/internal/store"` nếu chưa có; `bytes` đã có):

```go
// kiemTileSetInfo khang dinh goi -82: lenh dung, than = NGUYEN file tileset.
func kiemTileSetInfo(t *testing.T, m *protocol.Message) {
	t.Helper()
	if m.Command != -82 {
		t.Fatalf("lenh goi tileset = %d, mong doi -82", m.Command)
	}
	if !bytes.Equal(m.Data, duLieuTileSet) {
		t.Fatalf("than -82 khong bang nguyen file tileset (dai %d, mong doi %d)", len(m.Data), len(duLieuTileSet))
	}
}

// kiemMeInfo khang dinh goi -30/0 khop nhan vat da tao: byte lenh con 0, id/gender/
// name/gold/gem dung, gender lan 2 < 5. Doi chieu voi nv lay tu CSDL that.
func kiemMeInfo(t *testing.T, m *protocol.Message, nv *store.NhanVat) {
	t.Helper()
	if m.Command != -30 {
		t.Fatalf("lenh ME_INFO = %d, mong doi -30", m.Command)
	}
	r := protocol.NewReader(m.Data)
	if sub, err := r.Byte(); err != nil || sub != 0 {
		t.Fatalf("byte lenh con ME_INFO = %d (loi %v), mong doi 0", sub, err)
	}
	if id, err := r.Int(); err != nil || id != int32(nv.ID) {
		t.Fatalf("id = %d (loi %v), mong doi %d", id, err, nv.ID)
	}
	if _, err := r.Byte(); err != nil { // taskMain
		t.Fatal(err)
	}
	if g, err := r.Byte(); err != nil || g != byte(nv.Lang) {
		t.Fatalf("gender = %d (loi %v), mong doi %d (= lang)", g, err, nv.Lang)
	}
	if _, err := r.Short(); err != nil { // head
		t.Fatal(err)
	}
	if ten, err := r.UTF(); err != nil || ten != nv.Ten {
		t.Fatalf("name = %q (loi %v), mong doi %q", ten, err, nv.Ten)
	}
	r.Byte()  // cPK
	r.Byte()  // typePk
	r.Long()  // power
	r.Short() // buffHp
	r.Short() // buffMp
	g2, err := r.Byte()
	if err != nil || g2 != byte(nv.Lang) {
		t.Fatalf("gender lan 2 = %d (loi %v), mong doi %d (= lang)", g2, err, nv.Lang)
	}
	if g2 >= 5 {
		t.Fatalf("gender lan 2 = %d, phai < countClass (5) — khong client IndexOutOfRange", g2)
	}
	if sc, err := r.Byte(); err != nil || sc != 0 {
		t.Fatalf("skillCount = %d (loi %v), mong doi 0", sc, err)
	}
	if gold, err := r.Long(); err != nil || gold != nv.Vang {
		t.Fatalf("gold = %d (loi %v), mong doi %d (long 8 byte = Vang)", gold, err, nv.Vang)
	}
	r.Int() // ruby
	if gem, err := r.Int(); err != nil || gem != int32(nv.Ngoc) {
		t.Fatalf("gem = %d (loi %v), mong doi %d (= Ngoc)", gem, err, nv.Ngoc)
	}
}

// kiemMapInfo khang dinh goi -24 khop ban do cua nhan vat: mapId/tileId/name +
// byte du phong (decode het buffer den cuoi khong loi = du phong dung cho).
func kiemMapInfo(t *testing.T, m *protocol.Message, nv *store.NhanVat, bd *BanDo) {
	t.Helper()
	if m.Command != -24 {
		t.Fatalf("lenh MAP INFO = %d, mong doi -24", m.Command)
	}
	r := protocol.NewReader(m.Data)
	must := func(ten string, muon byte) {
		g, err := r.Byte()
		if err != nil || g != muon {
			t.Fatalf("%s = %d (loi %v), mong doi %d", ten, g, err, muon)
		}
	}
	must("mapId", byte(bd.MapID))
	must("planetId", bd.PlanetID)
	must("tileId", bd.TileID)
	must("bgId", bd.BgID)
	must("type", bd.Loai)
	if ten, err := r.UTF(); err != nil || ten != bd.Ten {
		t.Fatalf("mapName = %q (loi %v), mong doi %q", ten, err, bd.Ten)
	}
	must("zoneId", 0)
	if x, err := r.Short(); err != nil || x != int16(nv.ViTriX) {
		t.Fatalf("x = %d (loi %v), mong doi %d", x, err, nv.ViTriX)
	}
	if y, err := r.Short(); err != nil || y != int16(nv.ViTriY) {
		t.Fatalf("y = %d (loi %v), mong doi %d", y, err, nv.ViTriY)
	}
	must("waypointCount", 0)
	must("mobCount", 0)
	must("byte du phong", 0) // giua mob va npc
	must("npcCount", 0)
	must("itemCount", 0)
	if bg, err := r.Short(); err != nil || bg != 0 {
		t.Fatalf("bgItem = %d (loi %v), mong doi short 0", bg, err)
	}
	if ef, err := r.Short(); err != nil || ef != 0 {
		t.Fatalf("effMap = %d (loi %v), mong doi short 0", ef, err)
	}
	must("bgType", bd.BgType)
	must("idSpaceShip", 0)
	must("isMapDouble", 0)
	if _, err := r.Byte(); err == nil {
		t.Error("MAP INFO con byte thua sau isMapDouble — byte du phong co the sai vi tri")
	}
}
```

- [ ] **Step 2: Sửa phần cuối `TestVersionDanceDayDu`**

Trong `TestVersionDanceDayDu`:

1. Bắt `taiKhoanID` khi tạo tài khoản và thêm `ctx`:
```go
	ctx := context.Background()
	tenTK := "vaogame" + hauToNganNgau(t)
	matKhau := "matkhau123"
	taiKhoanID := taoTaiKhoanCoMatKhau(t, tenTK, matKhau)
```

2. **Thay** phần sau `cl.guiClientOk()` (KH3: xin lại skill để xác nhận máy chủ còn sống) **bằng** khối nhận 3 gói vào thế giới:

```go
	// Gui client ok -> Ke hoach 4: may chu gui NOT 3 goi vao the gioi theo thu tu
	// -82 (tileset) -> -30/0 (ME_INFO) -> -24 (MAP INFO). Doc va khang dinh byte-dung.
	cl.guiClientOk()

	// Lay nhan vat vua tao tu CSDL de doi chieu ME_INFO/MAP INFO byte-dung.
	nv, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, taiKhoanID)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan sau khi tao that bai: %v", err)
	}
	bd := TimBanDo(nv.MapID)
	if bd == nil {
		t.Fatalf("nhan vat o mapId %d khong co ban do trong danh muc", nv.MapID)
	}

	// doiGoiLenh khop THEO LENH (bo qua goi khac), dung duoc cho -82/-30/-24 vi ca
	// ba dung do dai 2 byte (khong nam trong lenhDoDaiBaByte). Thu tu: -82 truoc -24.
	kiemTileSetInfo(t, cl.doiGoiLenh(t, -82))
	kiemMeInfo(t, cl.doiGoiLenh(t, -30), nv)
	kiemMapInfo(t, cl.doiGoiLenh(t, -24), nv, bd)
```

> **Ghi chú:** `nv` tạo bằng `cl.taoNhanVat("Naru"+suffix, 0, 5)` → `Lang=0` (Làng Lá), `MapID=1`, `Vang=1000`, `Ngoc=0`, `ViTriX/Y=100/384` (từ `lang.go`). `kiemMeInfo`/`kiemMapInfo` đối chiếu với `nv` thật lấy từ CSDL nên không phụ thuộc giá trị cứng. `bd = TimBanDo(1)` = Làng Lá (tile 1, bg 0). Cập nhật chữ ký `taoTaiKhoanCoMatKhau` không đổi — nó **đã** trả `int64` (xem `tichhop_test.go:91`); chỉ cần bắt giá trị trả về (KH3 gọi nó mà bỏ giá trị).

- [ ] **Step 3: Chạy — phải PASS**

Run (trong `server/`):
```
$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro_test'
D:\Go\bin\go.exe test ./internal/game/ -run TestVersionDanceDayDu -v
```
Expected: PASS.

- [ ] **Step 4: Chạy TOÀN BỘ test + gofmt + vet**

Run (trong `server/`):
```
$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro_test'
D:\Go\bin\go.exe test ./...
D:\Go\bin\gofmt.exe -l internal/ cmd/
D:\Go\bin\go.exe vet ./...
```
Expected: mọi test PASS; `gofmt` chỉ còn khác biệt CRLF (nội dung sạch — nếu in tên file, chạy `D:\Go\bin\gofmt.exe -w` các file đó); `vet` không in gì.

- [ ] **Step 5: Commit**

```bash
git add server/internal/game/dongbo_tichhop_test.go
git commit -m "test(game): TestVersionDanceDayDu di tiep qua -28/13 -> nhan -82+-30/0+-24 byte-dung"
```

---

## Task 9: Xác minh với client Unity thật (render đứng yên trong làng)

**Files:** không sửa mã (chỉ chạy + ghi nhận).

**Interfaces:** toàn hệ thống.

- [ ] **Step 1: Chạy máy chủ thật**

Run (cửa sổ 1):
```
$env:NRO_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro'
$env:NRO_MUC_NHAT_KY='info'
D:\Go\bin\go.exe run ./cmd/nroserver
```

- [ ] **Step 2: Nối client Unity thật, quan sát nhật ký + màn hình**

Cấu hình client (`D:\NRO_Project`) trỏ về máy chủ local, đăng nhập tài khoản đã có nhân vật (hoặc tạo mới). Kỳ vọng nhật ký máy chủ tuần tự:
```
nhan vat vao the gioi, bat dau bat tay phien ban
client san sang (da gui -28/13), bat tay phien ban xong   ← KH3
da gui the gioi (tileset + ME_INFO + MAP INFO), client se render   ← KH4
```
Kỳ vọng màn hình client: **rời màn "Xin chờ"**, hiện nhân vật **đứng trong bản đồ làng** (đồ hoạ map cũ tương ứng mapId — vd Làng Lá dùng đồ hoạ "Đồi hoa cúc").

- [ ] **Step 3: Ghi nhận kết quả**

Ghi vào báo cáo task: (a) có thấy dòng "da gui the gioi" không; (b) client có rời "Xin chờ" và render nhân vật đứng yên không; (c) nếu client **thoát hẳn app** → nghi `-82` sai/thiếu hoặc tileId ngoài [1,39] (xem §7); nếu client **kẹt "Xin chờ"** → nghi ME_INFO lệch (gold 4 byte? quên byte lệnh con 0?) hoặc MAP INFO lệch (quên byte dự phòng?). KH4 dừng ở "render đứng yên"; **di chuyển là KH5** — nhân vật chưa đi được là đúng dự kiến.

- [ ] **Step 4: Commit (nếu có ghi chú)**

Không có mã để commit. Ghi kết quả xác minh vào báo cáo task.

---

## Tự soi lại kế hoạch

**1. Bám tài liệu giao thức:**
- §7 `-82`: gửi nguyên file 1696 byte, không byte lệnh con, trước `-24` → Task 2 (`goiTileSetInfo`) + Task 7 (thứ tự gửi). ✅
- §8 `-30/0` ME_INFO: **byte lệnh con 0 dẫn đầu** (đối chiếu `Service.java:782,911`), gold long 8 byte, gender ×2, 3 mảng item size 0, bỏ 3 trường cuối → Task 5. ✅
- §9 `-24` MAP INFO: **byte dự phòng** (đối chiếu `Zone.java:805`), bg-item short 0, effMap short 0, isMapDouble theo mapId → Task 6. ✅
- §10 tile map: map dùng lại mapId cũ (1-5) → client có sẵn `tile_map_data`, **không cần `-28/10`** → QĐ2 + Task 3. ✅
- §11 ràng buộc cứng → Global Constraints (sao chép giá trị). ✅
- Mô hình khu vực + tầm nhìn toàn khu vực (`2026-07-23...` §3.4/§4.3): `KhuVuc` tối thiểu, nhịp 100ms để KH5 → QĐ1 + Task 4. ✅

**2. Không còn chỗ trống:** mọi bước có mã Go thật, byte-exact; không "TODO"/"xử lý phù hợp". Chỉ một chỗ có nhánh lựa chọn (helper `fClientHetHanNgan` ở Task 7 Step 1) — đã nêu rõ quyết định gợi ý + mã kèm.

**3. Nhất quán kiểu:**
- Hằng lệnh: `lenhTileSetInfo/lenhMeInfo/lenhMapInfo int8` (Task 2), dùng ở Task 5/6/7. ✅
- `goiTileSetInfo() *Message`; `goiMeInfo(nv) (*Message, error)`; `goiMapInfo(nv, bd) (*Message, error)` — hai hàm sau trả `error` vì `UTF`; Task 7 xử lý error. ✅
- `TimBanDo(int32) *BanDo`, `MoKhuVuc(*BanDo) *KhuVuc`, `KhuVuc.BanDo() *BanDo`, `ThemNguoiChoi(*store.NhanVat)`, `SoNguoiChoi() int` — nhất quán Task 3/4/7. ✅
- `Reader.Long() (int64, error)` (Task 1) dùng trong test Task 5 + Task 8. ✅
- `BoXuLy.khuVuc *KhuVuc` (Task 7) — kiểu khớp Task 4. ✅

**4. Rủi ro cần người soi để mắt:**
- **Byte lệnh con 0 của ME_INFO** (Task 5): §8 bỏ sót; nếu quên `w.Byte(0)` đầu tiên → client kẹt "Xin chờ". Đã đối chiếu `Service.java:782,911-915`.
- **Byte dự phòng MAP INFO** (Task 6): quên → lệch toàn bộ phần sau. Đã đối chiếu `Zone.java:805`; test khẳng định decode hết buffer.
- **Thứ tự `-82` trước `-24`** (Task 7): sai → client thoát app. Test `TestClientOkGuiVaoTheGioi` khẳng định gói 1 = -82, gói 3 = -24.
- **`gold` long 8 byte** (Task 5): test byte-exact + `kiemMeInfo` đọc `r.Long()`.
- **`doiGoiLenh` khớp theo lệnh** (Task 8): dùng cho -82/-30/-24 (2-byte) — không dùng `doiGoi(cmd,sub)` vì -82/-24 không có byte lệnh con.
- **`tileId ∈ [1,39]`** (Task 3): test `TestTileIdMoiBanDoTrongKhoang`. Giá trị lấy từ `map_template` thật.

---

## Chỗ byte-table còn thiếu / cần coordinator lưu ý (báo cáo)

1. **§8 (ME_INFO) bỏ sót byte lệnh con `0` dẫn đầu.** Bảng bắt đầu từ `id int` (field 1), nhưng `Service.player` dùng `messageSubCommand((byte)0)` = `Message(-30)` + `writeByte(0)`. Kế hoạch này ghi `w.Byte(0)` trước tiên. **Đề nghị coordinator cập nhật §8** thêm "field 0: byte lệnh con = 0" cho khớp thực tế code.
2. **`isNewMember` (§8 field 23)** không có nguồn trong `store.NhanVat` — kế hoạch ghi cứng `0`. Nếu muốn phân biệt "thành viên mới" (hiện hiệu ứng chào), cần thêm cột. Ngoài phạm vi KH4.
3. **`conMeo[gender][0..2]` (§8 field 20)** — spec chốt `0,0,0` (Naruto không dùng "con chim thông báo" của DBZ). Bản Java đọc từ mảng tĩnh `conMeo`; ta cố định 0. Đúng cho map trống; nếu sau muốn thú cưng làng, cần dữ liệu riêng.
4. **`BanDo.Ten` (mapName)** đặt theo tên làng Naruto, nhưng **đồ hoạ vẫn là map DBZ cũ** (mapId 1-5 = Đồi hoa cúc/Thung lũng tre/Rừng nấm/Rừng xương/Đảo Kamê). Đây là nhãn UI tạm — coordinator quyết có đổi tên hay không tới khi dựng lại client. Không ảnh hưởng byte-exactness.

## Bàn giao thực thi

**Kế hoạch xong, lưu tại `docs/superpowers/plans/2026-07-25-ke-hoach-4-vao-the-gioi.md`. Hai lựa chọn thực thi:**

**1. Subagent-Driven (khuyến nghị)** — mỗi task một subagent, kiểm chéo + đột biến giữa các task (như Kế hoạch 1-3). REQUIRED SUB-SKILL: superpowers:subagent-driven-development.

**2. Inline Execution** — thực thi tuần tự trong phiên này, chốt kiểm tại mỗi checkpoint. REQUIRED SUB-SKILL: superpowers:executing-plans.

**Chọn cách nào?**
