# Kế hoạch 3 — Đồng bộ dữ liệu tĩnh Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Máy chủ Go đưa client Unity thật qua trọn "bắt tay phiên bản": gửi versionGame, phục vụ bốn catalog dữ liệu tĩnh (data/map/skill/item) với **5 làng Naruto thay cho 3 lớp Dragon Ball cũ**, để client gửi được "client ok" (`-28/13`).

**Architecture:** Nối tiếp Kế hoạch 2. Khi `vaoTheGioi` chạy (đăng nhập có nhân vật, hoặc vừa tạo nhân vật), máy chủ gửi `-28/4` versionGame. Client so phiên bản, tự xin bốn cập nhật (`-87` data, `-28/6` map, `-28/7` skill, `-28/8` item); máy chủ trả byte-đúng từng cái. Khi cả bốn phiên bản khớp, client gửi `-28/13`. Kế hoạch 3 dừng ở đó (Kế hoạch 4 sẽ gửi ME_INFO + MAP INFO để client render).

**Tech Stack:** Go 1.25+, `//go:embed` (thư viện chuẩn) để nhúng file part, `testing`. Không thư viện ngoài.

## Tài liệu tham chiếu giao thức (ĐỌC TRƯỚC)

**Mọi bảng byte lấy từ:** [`../specs/2026-07-24-giao-thuc-vao-the-gioi.md`](../specs/2026-07-24-giao-thuc-vao-the-gioi.md).
Kế hoạch này KHÔNG chép lại bảng byte — trích từ tài liệu đó. Các mục liên quan:
§0 (chuỗi tối thiểu), §1 (phiên bản), §2 (`-28/4`), §3 (`-87`), §4 (`-28/6`), §5 (`-28/7`), §6 (`-28/8`), §11 (ràng buộc cứng).

## Môi trường đã dựng sẵn

| Thứ | Vị trí / giá trị |
|---|---|
| Go | `D:\Go\bin\go.exe` — 1.26.5 · gofmt `D:\Go\bin\gofmt.exe` |
| PostgreSQL | cổng **5433**, user `postgres` / `nrodev123`, CSDL `nro`/`nro_test` |
| Test CSDL | `$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro_test'` |
| File part gốc | `1-MAY-CHU-GAME\data\update_data\part` (87.707 byte) |
| Client đích | `D:\NRO_Project` — đã xác nhận `ModFunc.isReadInt=true` |

## Global Constraints

- **Go 1.25+.** Không thư viện ngoài (chỉ `github.com/jackc/pgx/v5` đã có; task này chỉ dùng thư viện chuẩn `embed`).
- Ghi chú code **tiếng Việt**, giải thích *tại sao*.
- Mọi file dưới `server/`. **Không đụng `1-MAY-CHU-GAME/`.**
- **Không ghi nhật ký cho gói không xử lý** (lỗi ngập nhật ký đã có tiền lệ).
- Số nguyên big-endian, chuỗi modified UTF-8 — **dùng lại `protocol.Writer`**, không tự ghi byte.
- **Phiên bản Naruto `vsData=vsMap=vsSkill=vsItem=20`** (≠ Dragon Ball cũ 10/2/1/4, để ghi đè cache client).
- **`countClass` trong `-28/7` phải = số làng (5)**; mỗi làng gửi tên + `0` template.
- **Gói `-87` < 100.000 byte** (mốc `mark(100000)` client). Part 87.7KB + 5 blob rỗng ≈ 87.7KB, đạt.
- **`hp` mob = 4 byte int** (client `isReadInt=true`) — không dùng ở Kế hoạch 3 (catalog map rỗng) nhưng ghi nhớ.
- Mọi `count`/`size` đọc byte **có dấu** ở client → **≤127**.
- Mỗi task kết thúc bằng một commit.

## Khung gói `-28/X` (đã xác nhận từ code)

`Service.messageNotMap((byte)X)` = `Message(-28)` với **byte lệnh con `X` là byte đầu thân**. Mọi gói `-28/X` đều vậy. Riêng `-87` là `Message(-87)` trực tiếp, KHÔNG có byte lệnh con (thân bắt đầu bằng `vsData`).

Client **xin** cập nhật bằng cách gửi: `-87` (thân rỗng) cho data; `-28/6`, `-28/7`, `-28/8` (thân chỉ có byte lệnh con) cho map/skill/item. Máy chủ **trả** bằng cùng lệnh, thân đầy dữ liệu.

## Cấu trúc file

**Tạo mới:**
- `server/internal/game/part.bin` — bản sao file part gốc (placeholder Dragon Ball; thay Naruto khi có art). Nhúng bằng `//go:embed`.
- `server/internal/game/dongbo.go` — hằng phiên bản, các hàm dựng gói catalog, xử lý gói xin cập nhật.
- `server/internal/game/dongbo_test.go`

**Sửa:**
- `server/internal/protocol/datastream.go` — thêm `Writer.ThemMang([]byte)` để ghi khối byte thô (dùng cho part 87KB).
- `server/internal/protocol/datastream_test.go` — test `ThemMang`.
- `server/internal/game/login.go` — `XuLy` điều phối thêm `-87`; `xuLyKhongTrenBanDo` thêm lệnh con 6/7/8/13; `vaoTheGioi` gửi `-28/4`.
- `server/internal/game/login_test.go` — cập nhật nếu `vaoTheGioi` đổi hành vi (giờ gửi `-28/4`).

## Ranh giới Kế hoạch 3 ↔ Kế hoạch 4

| Kế hoạch 3 | Kế hoạch 4 |
|---|---|
| `-28/4` versionGame | `-82` tileSetInfo |
| Phục vụ `-87`/`-28/6`/`-28/7`/`-28/8` | `-30/0` ME_INFO + `-24` MAP INFO |
| `-28/7` = 5 làng Naruto | Mô hình khu vực, nạp tile map |
| `-28/13` → ghi nhận "client sẵn sàng" (stub) | `-28/13` → gửi ME_INFO + MAP INFO |

> **Kiểm thử Kế hoạch 3 bằng client giả** (`thuclient` mở rộng) đi trọn version
> dance và khẳng định catalog byte-đúng. Client Unity thật sẽ gửi `-28/13` rồi
> **kẹt "Xin chờ"** (chưa có ME_INFO/MAP INFO) — đúng dự kiến, Kế hoạch 4 hoàn tất.

---

## Task 1: Ghi khối byte thô vào Writer

**Vì sao:** gói `-87` chứa 87KB dữ liệu part. Ghi từng byte qua `Writer.Byte` được nhưng bẩn; thêm một hàm ghi cả khối.

**Files:**
- Modify: `server/internal/protocol/datastream.go`
- Modify: `server/internal/protocol/datastream_test.go`

**Interfaces:**
- Produces: `func (w *Writer) ThemMang(p []byte)` — nối `p` vào cuối bộ đệm, không thêm độ dài (khác `UTF` có tiền tố độ dài).

- [ ] **Step 1: Viết test đỏ**

Thêm vào `server/internal/protocol/datastream_test.go`:

```go
func TestThemMangGhiThoKhongTienTo(t *testing.T) {
	w := NewWriter()
	w.Byte(0xAA)
	w.ThemMang([]byte{0x01, 0x02, 0x03})
	w.Byte(0xBB)
	got := w.Bytes()
	want := []byte{0xAA, 0x01, 0x02, 0x03, 0xBB}
	if !bytes.Equal(got, want) {
		t.Errorf("ThemMang sai: got %v, want %v", got, want)
	}
}

func TestThemMangRong(t *testing.T) {
	w := NewWriter()
	w.Byte(0xAA)
	w.ThemMang(nil)
	if got := w.Bytes(); !bytes.Equal(got, []byte{0xAA}) {
		t.Errorf("ThemMang(nil) phai khong doi bo dem, got %v", got)
	}
}
```

Nếu `bytes` chưa import trong file test, thêm `import "bytes"`.

- [ ] **Step 2: Chạy — phải fail**

Run: `D:\Go\bin\go.exe test ./internal/protocol/ -run TestThemMang -v`
Expected: FAIL (`ThemMang` chưa có).

- [ ] **Step 3: Viết hàm**

Thêm vào `server/internal/protocol/datastream.go` (cạnh các hàm `Writer` khác):

```go
// ThemMang ghi ca khoi byte THO vao cuoi bo dem, KHONG kem tien to do dai (khac
// UTF). Dung cho cac blob nhi phan da co san (vi du du lieu part 87KB trong goi
// -87), nơi do dai da duoc ghi rieng bang Int() truoc do.
func (w *Writer) ThemMang(p []byte) {
	w.buf = append(w.buf, p...)
}
```

- [ ] **Step 4: Chạy — phải PASS**

Run: `D:\Go\bin\go.exe test ./internal/protocol/ -run TestThemMang -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add server/internal/protocol/datastream.go server/internal/protocol/datastream_test.go
git commit -m "feat(protocol): Writer.ThemMang ghi khoi byte tho (cho blob part)"
```

---

## Task 2: Nhúng file part + hằng phiên bản

**Files:**
- Create: `server/internal/game/part.bin` (sao chép từ file gốc)
- Create: `server/internal/game/dongbo.go`
- Create: `server/internal/game/dongbo_test.go`

**Interfaces:**
- Produces: `var duLieuPart []byte` (nhúng); hằng `vsData/vsMap/vsSkill/vsItem byte = 20`; hằng lệnh con `lenhConVersionGame/CatalogMap/CatalogSkill/CatalogItem/ClientOk`; hằng `lenhDuLieu int8 = -87`.

- [ ] **Step 1: Sao chép file part**

Run (PowerShell — sao chép nhị phân nguyên vẹn):
```
Copy-Item "D:\NRO_GoiBanGiao\1-MAY-CHU-GAME\data\update_data\part" "D:\NRO_GoiBanGiao\server\internal\game\part.bin"
(Get-Item "D:\NRO_GoiBanGiao\server\internal\game\part.bin").Length
```
Expected: `87707`.

- [ ] **Step 2: Viết test đỏ cho hằng + nhúng**

Tạo `server/internal/game/dongbo_test.go`:

```go
package game

import "testing"

func TestNhungPartDungKichThuoc(t *testing.T) {
	// File part nhung phai la file goc 87707 byte. Sai kich thuoc = nhung nham
	// hoac file hong -> nhan vat khong ve duoc o Ke hoach 4.
	if len(duLieuPart) != 87707 {
		t.Errorf("duLieuPart = %d byte, mong doi 87707 (file part goc)", len(duLieuPart))
	}
}

func TestPhienBanNarutoKhacBanCu(t *testing.T) {
	// Phai khac phien ban Dragon Ball cu (10/2/1/4) de ghi de cache client.
	// Dac biet vsSkill phai != 1, khong client dung cache 3 lop cu.
	if vsData == 10 || vsMap == 2 || vsSkill == 1 || vsItem == 4 {
		t.Error("phien ban Naruto trung ban Dragon Ball cu -> client se dung cache cu")
	}
	if vsData == 0 || vsMap == 0 || vsSkill == 0 || vsItem == 0 {
		t.Error("phien ban = 0 trung mac dinh client moi, khong ep cap nhat duoc dung cach")
	}
}
```

- [ ] **Step 3: Chạy — phải fail**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run "TestNhungPart|TestPhienBan" -v`
Expected: FAIL (biên dịch lỗi — `duLieuPart`, `vsData`... chưa có).

- [ ] **Step 4: Viết `dongbo.go` (phần đầu)**

Tạo `server/internal/game/dongbo.go`:

```go
// Package game — phan dong bo du lieu tinh (Ke hoach 3).
//
// Khi vao the gioi, client so 4 phien ban (data/map/skill/item) may chu gui
// trong -28/4 voi cache cua no; lech thi tu xin cap nhat. Client MOI (RMS trong)
// luon xin lai ca 4 du phien ban co khop hay khong, nen may chu PHAI phuc vu du
// bon o lan dau. Day la cho ta thay 3 lop Dragon Ball cu bang 5 lang Naruto.
package game

import (
	_ "embed"

	mynet "nroserver/internal/net"
	"nroserver/internal/protocol"
)

// duLieuPart: du lieu sprite nhan vat (dau/than/chan). Day la PLACEHOLDER cua
// ban Dragon Ball cu — nhan vat hien ra bang sprite cu; thay bang part Naruto
// khi duong ong my thuat giao. Client ve nhan vat bang parts[head/body/leg]
// KHONG kiem bien, nen part rong = sap khi ve (Ke hoach 4). File nhung phai phu
// du chi so head/body/leg cua moi nhan vat -> forward nguyen file cu la an toan.
//
//go:embed part.bin
var duLieuPart []byte

// Phien ban du lieu tinh cua ban Naruto. PHAI khac ban Dragon Ball cu
// (vsData=10, vsMap=2, vsSkill=1, vsItem=4) de client khong dung cache cu — nhat
// la vsSkill != 1 de khong dung cache 3 lop cu. 20 la gia tri tuy chon != cu va
// != 0. Bump khi doi noi dung tuong ung o cac ke hoach sau.
const (
	vsData  byte = 20
	vsMap   byte = 20
	vsSkill byte = 20
	vsItem  byte = 20
)

// lenhDuLieu: goi data. Client gui -87 than rong de xin, may chu tra -87 day du.
const lenhDuLieu int8 = -87

// Cac lenh con cua khung -28 lien quan dong bo tinh.
const (
	lenhConVersionGame  byte = 4  // may chu -> client
	lenhConCatalogMap   byte = 6  // client xin map / may chu tra
	lenhConCatalogSkill byte = 7  // client xin skill / may chu tra
	lenhConCatalogItem  byte = 8  // client xin item / may chu tra
	lenhConClientOk     byte = 13 // client bao san sang
)
```

- [ ] **Step 5: Chạy — phải PASS**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run "TestNhungPart|TestPhienBan" -v`
Expected: PASS (2 test).

- [ ] **Step 6: Commit**

```bash
git add server/internal/game/part.bin server/internal/game/dongbo.go server/internal/game/dongbo_test.go
git commit -m "feat(game): nhung file part (placeholder), hang phien ban Naruto=20"
```

---

## Task 3: Dựng gói versionGame `-28/4` và data `-87`

**Files:**
- Modify: `server/internal/game/dongbo.go`
- Modify: `server/internal/game/dongbo_test.go`

**Interfaces:**
- Consumes: `protocol.NewWriter`, `w.Byte/Int/ThemMang`, `protocol.NewMessage`, `lenhKhongTrenBanDo int8 = -28` (đã có trong `login.go`), hằng Task 2.
- Produces: `func goiVersionGame() *protocol.Message`; `func goiCatalogData() *protocol.Message`.

- [ ] **Step 1: Viết test đỏ**

Thêm vào `dongbo_test.go` (đối chiếu tài liệu §2, §3):

```go
import "nroserver/internal/protocol" // them vao import neu chua co

func TestGoiVersionGameDungByte(t *testing.T) {
	m := goiVersionGame()
	if m.Command != -28 {
		t.Fatalf("lenh = %d, mong doi -28", m.Command)
	}
	// than: [4][vsData][vsMap][vsSkill][vsItem][0][countCaption=0]
	want := []byte{4, 20, 20, 20, 20, 0, 0}
	if !bytesEqual(m.Data, want) {
		t.Errorf("than versionGame = %v, mong doi %v", m.Data, want)
	}
}

func TestGoiCatalogDataDungCauTruc(t *testing.T) {
	m := goiCatalogData()
	if m.Command != -87 {
		t.Fatalf("lenh = %d, mong doi -87", m.Command)
	}
	r := protocol.NewReader(m.Data)
	if b, _ := r.Byte(); b != vsData {
		t.Errorf("byte dau = %d, mong doi vsData %d", b, vsData)
	}
	// 6 blob theo thu tu dart/arrow/effect/image/part/skill; chi part co noi dung.
	doDaiBlob := func(ten string) int32 {
		n, err := r.Int()
		if err != nil {
			t.Fatalf("doc do dai %s loi: %v", ten, err)
		}
		// bo qua noi dung blob
		for i := int32(0); i < n; i++ {
			if _, err := r.Byte(); err != nil {
				t.Fatalf("doc noi dung %s loi: %v", ten, err)
			}
		}
		return n
	}
	if n := doDaiBlob("dart"); n != 2 {
		t.Errorf("dart do dai %d, mong doi 2 (rong)", n)
	}
	doDaiBlob("arrow")
	doDaiBlob("effect")
	doDaiBlob("image")
	if n := doDaiBlob("part"); n != 87707 {
		t.Errorf("part do dai %d, mong doi 87707", n)
	}
	doDaiBlob("skill")
	// Toan goi phai < 100000 byte (moc mark(100000) client).
	if len(m.Data) >= 100000 {
		t.Errorf("goi -87 = %d byte, phai < 100000", len(m.Data))
	}
}
```

> **Ghi chú:** `bytesEqual` — nếu chưa có helper trong package test, dùng `bytes.Equal` (thêm `import "bytes"`) hoặc so từng byte. `protocol.NewReader` và `Reader.Int()` đã có từ Kế hoạch 1.

- [ ] **Step 2: Chạy — phải fail**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run "TestGoiVersionGame|TestGoiCatalogData" -v`
Expected: FAIL (hàm chưa có).

- [ ] **Step 3: Viết hai hàm dựng gói**

Thêm vào `dongbo.go`:

```go
// goiVersionGame dung goi -28/4: bao client 4 phien ban du lieu tinh + bang moc
// kinh nghiem rong. Client so 4 phien ban nay voi cache; lech thi tu xin cap nhat.
// Xem tai lieu giao thuc §2.
func goiVersionGame() *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(lenhConVersionGame) // 4
	w.Byte(vsData)
	w.Byte(vsMap)
	w.Byte(vsSkill)
	w.Byte(vsItem)
	w.Byte(0) // byte dem, client doc roi bo
	w.Byte(0) // countCaption = 0 (bang moc kinh nghiem rong, an toan cho vao map)
	return protocol.NewMessage(lenhKhongTrenBanDo, w)
}

// goiCatalogData dung goi -87: vsData roi 6 blob (dart/arrow/effect/image/part/
// skill) theo dung thu tu. Chi part co noi dung (du lieu sprite nhan vat); 5 blob
// kia rong (do dai 2 = mot short 0). Client chi LUU raw vao RMS, khong parse ngay
// — loi part chi lo khi VE nhan vat (Ke hoach 4). Xem tai lieu giao thuc §3.
func goiCatalogData() *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(vsData)
	// blob rong = do dai 2 byte, noi dung la mot short 0 (count = 0 phan tu).
	blobRong := func() {
		w.Int(2)
		w.Short(0)
	}
	blobRong() // dart
	blobRong() // arrow
	blobRong() // effect
	blobRong() // image
	w.Int(int32(len(duLieuPart))) // part: do dai that
	w.ThemMang(duLieuPart)        // part: noi dung
	blobRong() // skill
	return protocol.NewMessage(lenhDuLieu, w)
}
```

- [ ] **Step 4: Chạy — phải PASS**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run "TestGoiVersionGame|TestGoiCatalogData" -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add server/internal/game/dongbo.go server/internal/game/dongbo_test.go
git commit -m "feat(game): dung goi -28/4 versionGame va -87 data (part + 5 blob rong)"
```

---

## Task 4: Dựng gói skill `-28/7` (5 làng Naruto) + map `-28/6` + item `-28/8`

**Files:**
- Modify: `server/internal/game/dongbo.go`
- Modify: `server/internal/game/dongbo_test.go`

**Interfaces:**
- Consumes: `DanhSachLang() []Lang` (Kế hoạch 2, `lang.go`), `Lang.Ten`.
- Produces: `func goiCatalogSkill() (*protocol.Message, error)`; `func goiCatalogMap() *protocol.Message`; `func goiCatalogItem() *protocol.Message`.

- [ ] **Step 1: Viết test đỏ**

Thêm vào `dongbo_test.go`:

```go
func TestGoiCatalogSkillCoDu5Lang(t *testing.T) {
	m, err := goiCatalogSkill()
	if err != nil {
		t.Fatalf("goiCatalogSkill loi: %v", err)
	}
	if m.Command != -28 {
		t.Fatalf("lenh = %d, mong doi -28", m.Command)
	}
	r := protocol.NewReader(m.Data)
	sub, _ := r.Byte()
	if sub != lenhConCatalogSkill { // 7
		t.Errorf("lenh con = %d, mong doi 7", sub)
	}
	if v, _ := r.Byte(); v != vsSkill {
		t.Errorf("vsSkill = %d, mong doi %d", v, vsSkill)
	}
	if opt, _ := r.Byte(); opt != 0 {
		t.Errorf("countOption = %d, mong doi 0", opt)
	}
	countClass, _ := r.Byte()
	dsLang := DanhSachLang()
	if int(countClass) != len(dsLang) {
		t.Fatalf("countClass = %d, mong doi %d (so lang)", countClass, len(dsLang))
	}
	if countClass < 5 {
		t.Errorf("countClass = %d, phai >= 5 de id lang 0-4 hop le trong ME_INFO", countClass)
	}
	// Doc tung lang: ten (UTF) + countTemplate (0). Ten phai khop danh muc lang.
	for i := 0; i < int(countClass); i++ {
		ten, err := r.UTF()
		if err != nil {
			t.Fatalf("doc ten lang %d loi: %v", i, err)
		}
		if ten != dsLang[i].Ten {
			t.Errorf("ten lang %d = %q, mong doi %q", i, ten, dsLang[i].Ten)
		}
		if ct, _ := r.Byte(); ct != 0 {
			t.Errorf("countTemplate lang %d = %d, mong doi 0", i, ct)
		}
	}
}

func TestGoiCatalogMapRong(t *testing.T) {
	m := goiCatalogMap()
	if m.Command != -28 {
		t.Fatalf("lenh = %d, mong doi -28", m.Command)
	}
	// [6][vsMap][countMap=0][countNpc=0][countMob=0]
	want := []byte{6, 20, 0, 0, 0}
	if !bytes.Equal(m.Data, want) {
		t.Errorf("than catalog map = %v, mong doi %v", m.Data, want)
	}
}

func TestGoiCatalogItemSubtype2Rong(t *testing.T) {
	m := goiCatalogItem()
	if m.Command != -28 {
		t.Fatalf("lenh = %d, mong doi -28", m.Command)
	}
	// [8][vsItem][subtype=2][start short=0][end short=0]
	want := []byte{8, 20, 2, 0, 0, 0, 0}
	if !bytes.Equal(m.Data, want) {
		t.Errorf("than catalog item = %v, mong doi %v", m.Data, want)
	}
}
```

- [ ] **Step 2: Chạy — phải fail**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run "TestGoiCatalog" -v`
Expected: FAIL (3 hàm chưa có).

- [ ] **Step 3: Viết ba hàm**

Thêm vào `dongbo.go`:

```go
// goiCatalogSkill dung goi -28/7: DANH SACH 5 LANG NARUTO thay cho 3 lop Dragon
// Ball cu. Client dung byte "gender" trong ME_INFO lam chi so nClasss[] nen
// countClass PHAI >= so lang (5) — id lang 3/4 se sap neu client con cache 3 lop
// cu. Moi lang gui ten + 0 template (chua co skill). Xem tai lieu §5.
func goiCatalogSkill() (*protocol.Message, error) {
	w := protocol.NewWriter()
	w.Byte(lenhConCatalogSkill) // 7
	w.Byte(vsSkill)
	w.Byte(0) // countOption = 0
	dsLang := DanhSachLang()
	w.Byte(byte(len(dsLang))) // countClass = 5
	for _, l := range dsLang {
		if err := w.UTF(l.Ten); err != nil {
			return nil, err
		}
		w.Byte(0) // countTemplate = 0 (chua co skill)
	}
	return protocol.NewMessage(lenhKhongTrenBanDo, w), nil
}

// goiCatalogMap dung goi -28/6 RONG: 0 map, 0 npc, 0 mob. Du de qua bat tay;
// noi dung map/npc/mob that se them o cac ke hoach sau. Xem tai lieu §4.
func goiCatalogMap() *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(lenhConCatalogMap) // 6
	w.Byte(vsMap)
	w.Byte(0) // countMap
	w.Byte(0) // countNpc
	w.Byte(0) // countMob
	return protocol.NewMessage(lenhKhongTrenBanDo, w)
}

// goiCatalogItem dung goi -28/8 toi thieu: dung MOT message con subtype 2 rong.
// Chi subtype 2 moi kich clientOk (khi du 4 phien ban khop). Nhan vat moi khong
// do van ve duoc voi catalog item rong. Xem tai lieu §6.
func goiCatalogItem() *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(lenhConCatalogItem) // 8
	w.Byte(vsItem)
	w.Byte(2)  // subtype 2 (item them)
	w.Short(0) // start
	w.Short(0) // end
	return protocol.NewMessage(lenhKhongTrenBanDo, w)
}
```

> Nếu `bytes` chưa import trong `dongbo_test.go`, thêm `import "bytes"`.

- [ ] **Step 4: Chạy — phải PASS**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run "TestGoiCatalog" -v`
Expected: PASS (3 test).

- [ ] **Step 5: Commit**

```bash
git add server/internal/game/dongbo.go server/internal/game/dongbo_test.go
git commit -m "feat(game): dung goi catalog skill (5 lang Naruto), map rong, item subtype-2 rong"
```

---

## Task 5: Điều phối — gửi versionGame + phục vụ gói xin cập nhật

**Files:**
- Modify: `server/internal/game/login.go`
- Modify: `server/internal/game/login_test.go`

**Interfaces:**
- Consumes: `goiVersionGame`, `goiCatalogData`, `goiCatalogSkill`, `goiCatalogMap`, `goiCatalogItem` (Task 3-4); `s.Gui`, `b.nhanVat`, `b.taiKhoan`.
- Produces: `vaoTheGioi` gửi `-28/4`; `XuLy` xử lý `-87`; `xuLyKhongTrenBanDo` xử lý lệnh con 6/7/8/13. Trạng thái phiên: `BoXuLy.clientSanSang bool`.

- [ ] **Step 1: Cập nhật/viết test đỏ**

Trong `login_test.go`:

1. `vaoTheGioi` giờ gửi `-28/4`. Test Kế hoạch 2 `TestDangNhapDaCoNhanVat_KhongGuiManTao` kỳ vọng "không gửi lệnh 2" — vẫn đúng, nhưng giờ có gửi `-28/4`. Nếu test đó khẳng định "không gửi GÌ", nới thành "không gửi lệnh 2, có gửi -28/4". Đọc test hiện tại và điều chỉnh khẳng định cho khớp hành vi mới (đăng nhập có nhân vật → gửi `-28/4`).

2. Thêm test điều phối gói xin cập nhật (dùng `dungPhienThat` + `khoGia`, đặt `b.nhanVat` sẵn để mô phỏng đã vào thế giới):

```go
func TestPhucVuGoiXinData(t *testing.T) {
	b := NewBoXuLy(&khoGia{})
	b.nhanVat = &store.NhanVat{ID: 1, Lang: 0}
	fClient, dong := dungPhienThat(t, b)
	defer dong()
	// Client xin data: gui -87 than rong.
	if err := fClient.WriteMessage(&protocol.Message{Command: -87}); err != nil {
		t.Fatal(err)
	}
	m := docMotGoi(t, fClient)
	if m.Command != -87 {
		t.Errorf("tra loi xin data = lenh %d, mong doi -87", m.Command)
	}
}

func TestPhucVuGoiXinSkill5Lang(t *testing.T) {
	b := NewBoXuLy(&khoGia{})
	b.nhanVat = &store.NhanVat{ID: 1, Lang: 0}
	fClient, dong := dungPhienThat(t, b)
	defer dong()
	// Client xin skill: -28 lenh con 7.
	w := protocol.NewWriter()
	w.Byte(7)
	if err := fClient.WriteMessage(protocol.NewMessage(-28, w)); err != nil {
		t.Fatal(err)
	}
	m := docMotGoi(t, fClient)
	r := protocol.NewReader(m.Data)
	sub, _ := r.Byte()
	if m.Command != -28 || sub != 7 {
		t.Fatalf("tra loi xin skill = lenh %d/%d, mong doi -28/7", m.Command, sub)
	}
	r.Byte() // vsSkill
	r.Byte() // countOption
	countClass, _ := r.Byte()
	if countClass < 5 {
		t.Errorf("catalog skill tra ve %d lop, phai >= 5 lang", countClass)
	}
}
```

> **Helper `docMotGoi`:** đọc đúng một gói máy chủ gửi ra qua `fClient` (framer client giả). Nếu Kế hoạch 1/2 đã có helper đọc gói (ví dụ trong `chayDangNhap`/`docGoi`), tái dùng; nếu chưa có dạng "đọc một gói", trích ra `docMotGoi(t, f) *protocol.Message { m, err := f.ReadMessage(); ... }`.

- [ ] **Step 2: Chạy — phải fail**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run "TestPhucVuGoi" -v`
Expected: FAIL (điều phối chưa xử lý -87 / lệnh con 6/7/8).

- [ ] **Step 3: Thêm trạng thái + gửi versionGame trong `vaoTheGioi`**

Trong `login.go`, thêm trường vào `BoXuLy` (cạnh `taiKhoan`, `nhanVat`):

```go
	clientSanSang bool // client da gui -28/13 (bat tay xong) — Ke hoach 4 se gui ME_INFO/MAP INFO
```

Sửa `vaoTheGioi` (thêm gửi `-28/4` sau khi gán `b.nhanVat`):

```go
func (b *BoXuLy) vaoTheGioi(s *mynet.Session, nv *store.NhanVat) error {
	b.nhanVat = nv
	slog.Info("nhan vat vao the gioi, bat dau bat tay phien ban",
		"ten", nv.Ten, "lang", nv.Lang, "id", nv.ID)
	// Gui versionGame de kich client so phien ban va tu xin 4 catalog. Client
	// gui -28/13 khi ca 4 phien ban khop; luc do Ke hoach 4 se gui ME_INFO + MAP
	// INFO. Ke hoach 3 dung o buoc client gui -28/13.
	return s.Gui(goiVersionGame())
}
```

- [ ] **Step 4: Điều phối `-87` và lệnh con 6/7/8/13**

Trong `XuLy`, thêm nhánh `-87`:

```go
	switch m.Command {
	case lenhChuaDangNhap: // -29
		return b.xuLyDangNhap(s, m)
	case lenhKhongTrenBanDo: // -28
		return b.xuLyKhongTrenBanDo(s, m)
	case lenhDuLieu: // -87 — client xin data
		return b.xuLyXinData(s)
	default:
		return nil
	}
```

Thêm `xuLyXinData` (trong `login.go` hoặc `dongbo.go`):

```go
// xuLyXinData tra goi -87 khi client xin du lieu tinh (part...). Chi phuc vu khi
// phien da vao the gioi (vaoTheGioi da chay); truoc do bo qua im lang (khong log
// — cung ly do chong ngap nhat ky).
func (b *BoXuLy) xuLyXinData(s *mynet.Session) error {
	if b.nhanVat == nil {
		return nil
	}
	return s.Gui(goiCatalogData())
}
```

Trong `xuLyKhongTrenBanDo`, thêm các lệnh con:

```go
	switch lenhCon {
	case lenhConTaoNhanVat: // 2
		return b.xuLyTaoNhanVat(s, r)
	case lenhConCatalogMap: // 6
		return b.phucVuCatalog(s, goiCatalogMap())
	case lenhConCatalogSkill: // 7
		m, err := goiCatalogSkill()
		if err != nil {
			return err
		}
		return b.phucVuCatalog(s, m)
	case lenhConCatalogItem: // 8
		return b.phucVuCatalog(s, goiCatalogItem())
	case lenhConClientOk: // 13
		return b.xuLyClientOk(s)
	default:
		return nil
	}
```

Thêm `phucVuCatalog` và `xuLyClientOk`:

```go
// phucVuCatalog gui mot goi catalog khi client xin, chi khi da vao the gioi.
func (b *BoXuLy) phucVuCatalog(s *mynet.Session, m *protocol.Message) error {
	if b.nhanVat == nil {
		return nil
	}
	return s.Gui(m)
}

// xuLyClientOk: client bao da nhan du 4 catalog va san sang. Ke hoach 3 chi ghi
// nhan; Ke hoach 4 se gui ME_INFO + MAP INFO tai day de client render.
func (b *BoXuLy) xuLyClientOk(s *mynet.Session) error {
	if b.nhanVat == nil {
		return nil
	}
	b.clientSanSang = true
	slog.Info("client san sang (da gui -28/13), bat tay phien ban xong",
		"ten", b.nhanVat.Ten, "id", b.nhanVat.ID)
	return nil
}
```

- [ ] **Step 5: Chạy toàn package game — phải PASS**

Run (trong `server/`):
```
$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro_test'
D:\Go\bin\go.exe test ./internal/game/ -v
```
Expected: PASS hết (test Kế hoạch 1-2 + điều phối mới). `gofmt -l internal/game/` sạch (CRLF thì `gofmt -w`), `go vet ./internal/game/` sạch.

- [ ] **Step 6: Commit**

```bash
git add server/internal/game/login.go server/internal/game/dongbo.go server/internal/game/login_test.go
git commit -m "feat(game): vaoTheGioi gui -28/4; phuc vu goi xin -87/-28/6/-28/7/-28/8; ghi nhan -28/13"
```

---

## Task 6: Thử đầu-cuối — client giả đi trọn version dance

**Files:**
- Modify: `server/cmd/thuclient/main.go`
- Create: `server/internal/game/dongbo_tichhop_test.go`

**Interfaces:**
- Consumes: toàn chuỗi (đăng nhập → tạo/tải nhân vật → versionGame → phục vụ catalog → -28/13).

- [ ] **Step 1: Viết test tích hợp — mô phỏng version dance đầy đủ**

Tạo `server/internal/game/dongbo_tichhop_test.go`. Dùng CSDL thật (`t.Skip` nếu thiếu `NRO_TEST_DSN`), kho `store.Store` thật, `net.Pipe`. Kịch bản: tài khoản mới → tạo nhân vật → nhận `-28/4` → client giả xin `-87`/`-28/6`/`-28/7`/`-28/8` → nhận 4 catalog byte-đúng → gửi `-28/13` → máy chủ ghi nhận sẵn sàng.

```go
func TestVersionDanceDayDu(t *testing.T) {
	khoThat := moKhoThatChoTest(t) // *store.Store toi nro_test, skip neu thieu DSN
	ctx := context.Background()

	tenTK := "vaogame" + hauToNganNgau(t)
	matKhau := "matkhau123"
	taoTaiKhoanCoMatKhau(t, khoThat, tenTK, matKhau)

	cl := moPhienClientGia(t, khoThat)
	defer cl.dong()
	cl.dangNhap(tenTK, matKhau)
	if !cl.choLenh(t, 2) {
		t.Fatal("chua co nhan vat -> phai nhan lenh 2")
	}
	// Tao nhan vat -> vaoTheGioi -> nhan -28/4.
	cl.taoNhanVat("Naru"+hauToNganNgau(t), 0, 5)
	m4 := cl.doiGoi(t, -28, 4) // doc toi khi gap -28 lenh con 4
	kiemVersionGame(t, m4)     // khang dinh vsData/map/skill/item = 20

	// Client gia xin ca 4 cap nhat, nhan tra loi byte-dung.
	cl.xinData()
	kiemDataCoPartVaDuoi100KB(t, cl.doiGoiLenh(t, -87))

	cl.xinCatalog(6)
	kiemCatalogMapRong(t, cl.doiGoi(t, -28, 6))

	cl.xinCatalog(7)
	kiemCatalogSkill5Lang(t, cl.doiGoi(t, -28, 7)) // KHANG DINH 5 lang, dung ten

	cl.xinCatalog(8)
	kiemCatalogItemSubtype2(t, cl.doiGoi(t, -28, 8))

	// Gui client ok -> may chu ghi nhan (khong sap, khong loi).
	cl.guiClientOk()
	// Ke hoach 3 dung o day: may chu KHONG gui ME_INFO/MAP INFO (do la Ke hoach 4).
	// Xac nhan may chu con song (khong ngat ket noi) bang cach xin lai skill 1 lan
	// nua va van nhan duoc tra loi.
	cl.xinCatalog(7)
	if m := cl.doiGoi(t, -28, 7); m == nil {
		t.Fatal("sau -28/13 may chu phai con phuc vu, nhung khong tra loi")
	}
}
```

> Các helper `moPhienClientGia`, `dangNhap`, `taoNhanVat`, `choLenh`, `doiGoi`, `xinData`, `xinCatalog`, `guiClientOk` mở rộng từ khung `clientGia`/`net.Pipe` Kế hoạch 2 (`tichhop_test.go`). `doiGoi(t, cmd, sub)` đọc các gói tới khi gặp gói có `Command==cmd` và byte lệnh con `==sub` (bỏ qua gói khác), có hạn giờ tránh treo. `kiem*` là các hàm khẳng định byte-đúng theo tài liệu §2-§6. **Khẳng định quan trọng nhất: `-28/7` chứa đúng 5 làng với đúng tên** (đây là mục tiêu lõi của kế hoạch — thay 3 lớp Dragon Ball bằng 5 làng Naruto).

- [ ] **Step 2: Chạy — phải PASS**

Run (trong `server/`):
```
$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro_test'
D:\Go\bin\go.exe test ./internal/game/ -run TestVersionDanceDayDu -v
```
Expected: PASS.

- [ ] **Step 3: Mở rộng `thuclient` cho version dance (thử tay với client thật/giả)**

Trong `server/cmd/thuclient/main.go`, sau bước tạo nhân vật (Kế hoạch 2), thêm: đọc `-28/4`, tự gửi `-87`/`-28/6`/`-28/7`/`-28/8` để xin, in các catalog nhận về (đặc biệt in số làng và tên làng trong `-28/7`), rồi gửi `-28/13`. Mục đích: chạy tay đối chiếu. Giữ đúng lối `thuclient` hiện có.

- [ ] **Step 4: Chạy TOÀN BỘ test lần cuối**

Run (trong `server/`):
```
$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro_test'
D:\Go\bin\go.exe test ./...
D:\Go\bin\gofmt.exe -l internal/ cmd/
D:\Go\bin\go.exe vet ./...
```
Expected: mọi test PASS; `gofmt` chỉ còn khác biệt CRLF (nội dung sạch); `vet` không in gì.

- [ ] **Step 5: Commit**

```bash
git add server/cmd/thuclient/main.go server/internal/game/dongbo_tichhop_test.go
git commit -m "test(game): tich hop version dance day du (5 lang Naruto); thuclient di het bat tay"
```

---

## Task 7: Xác minh với client Unity thật

**Files:** không sửa mã (chỉ chạy + ghi nhận).

**Interfaces:** toàn hệ thống.

- [ ] **Step 1: Chạy máy chủ thật**

Run (cửa sổ 1):
```
$env:NRO_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro'
$env:NRO_MUC_NHAT_KY='info'
D:\Go\bin\go.exe run ./cmd/nroserver
```

- [ ] **Step 2: Nối client Unity thật, quan sát nhật ký**

Cấu hình client trỏ về máy chủ (địa chỉ local), tạo nhân vật (hoặc dùng tài khoản có sẵn), theo dõi nhật ký máy chủ. Kỳ vọng thấy tuần tự:
```
nhan vat vao the gioi, bat dau bat tay phien ban
(client xin -87/-28/6/-28/7/-28/8 — may chu phuc vu)
client san sang (da gui -28/13), bat tay phien ban xong
```
Đây là bằng chứng client thật đã **nhận và chấp nhận 5 làng Naruto** (nếu client từ chối catalog, nó không bao giờ gửi `-28/13`).

- [ ] **Step 3: Ghi nhận kết quả**

Client thật sẽ **kẹt ở "Xin chờ" sau `-28/13`** — đúng dự kiến (Kế hoạch 4 gửi ME_INFO + MAP INFO). Ghi vào báo cáo: có thấy dòng "client san sang" không (bắt tay thành công), và xác nhận đây là điểm dừng đúng của Kế hoạch 3.

- [ ] **Step 4: Commit (nếu có ghi chú)**

Không có mã để commit. Ghi kết quả xác minh vào báo cáo task.

---

## Tự soi lại kế hoạch

**1. Bám tài liệu giao thức:**
- §2 `-28/4`: byte con 4 dẫn đầu → Task 3 code có `w.Byte(lenhConVersionGame)`. ✅
- §3 `-87`: 6 blob, chỉ part có nội dung, <100KB → Task 3. ✅
- §4 `-28/6` rỗng, §5 `-28/7` 5 làng, §6 `-28/8` subtype-2 → Task 4. ✅
- §1 phiên bản Naruto ≠ cũ → Task 2 hằng + test `TestPhienBanNarutoKhacBanCu`. ✅
- Chuỗi §0: `vaoTheGioi` gửi `-28/4`, phục vụ 4 catalog, ghi nhận `-28/13` → Task 5. ✅

**2. Không còn chỗ trống:** mọi bước có mã thật; không "TODO"/"xử lý phù hợp".

**3. Nhất quán kiểu:**
- Hằng lệnh con `byte`; `lenhDuLieu int8 = -87`; `lenhKhongTrenBanDo int8 = -28` (đã có). ✅
- `goiCatalogSkill` trả `(*Message, error)` (vì `UTF` trả error); ba hàm kia trả `*Message`. Điều phối Task 5 xử lý error của skill. ✅
- `DanhSachLang()` (Kế hoạch 2) → `Lang.Ten` dùng trong `-28/7`. ✅

**4. Rủi ro cần người soi để mắt:**
- Byte con `-28/4` = 4 (Task 3): thiếu là lệch cả gói. Đã xác nhận `messageNotMap((byte)4)` từ code thật.
- Part 87.7KB + 5 blob rỗng phải < 100KB: test `TestGoiCatalogData` khẳng định `len < 100000`.
- Helper test tên gợi ý (`moPhienClientGia`, `doiGoi`...) phải khớp tên thật Kế hoạch 2.
- `vaoTheGioi` đổi hành vi (giờ gửi `-28/4`) → test Kế hoạch 2 phụ thuộc `vaoTheGioi` phải cập nhật khẳng định (Task 5 Step 1).

---

## Bàn giao thực thi

**Kế hoạch xong, lưu tại `docs/superpowers/plans/2026-07-24-ke-hoach-3-dong-bo-du-lieu-tinh.md`.**

Thực thi: **Subagent-Driven** (mỗi task một subagent, kiểm chéo + đột biến giữa các task) — như Kế hoạch 1-2.
