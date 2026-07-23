# Kế hoạch 1 — Lõi giao thức + Đăng nhập

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Viết máy chủ Go nói được giao thức nhị phân của client Unity hiện có, để client thật kết nối, bắt tay khoá mã, đăng nhập bằng tài khoản trong PostgreSQL, và tới được màn hình chọn/tạo nhân vật.

**Architecture:** Ba tầng tách bạch. Tầng `protocol` chỉ biết byte và khoá mã, không biết gì về game. Tầng `net` quản lý kết nối TCP và vòng đời phiên. Tầng `game` xử lý nghiệp vụ đăng nhập. Mọi thứ liên quan giao thức cũ gói kín trong `internal/protocol` để sau này thay client mới chỉ cần thêm tầng khác.

**Tech Stack:** Go 1.24+ (máy đang có 1.26.5), PostgreSQL 16.4, `github.com/jackc/pgx/v5`, thư viện chuẩn `crypto/pbkdf2`, `testing`.

## Môi trường đã dựng sẵn trên máy này

| Thứ | Vị trí / giá trị |
|---|---|
| Go | `D:\Go\bin\go.exe` — phiên bản 1.26.5, đã có trong PATH |
| PostgreSQL | `D:\PostgreSQL\bin\` — 16.4, **cổng 5433** |
| Người dùng CSDL | `postgres` / mật khẩu `nrodev123` |
| CSDL chạy thật | `nro` |
| CSDL cho kiểm thử | `nro_test` |
| Chuỗi kết nối | `postgres://postgres:nrodev123@127.0.0.1:5433/nro` |

## Global Constraints

- **Go 1.24 trở lên** (bắt buộc: `crypto/pbkdf2` chỉ có trong thư viện chuẩn từ 1.24).
- **Không dùng thư viện ngoài, ngoại trừ đúng một cái: `github.com/jackc/pgx/v5`.**
  Cụ thể: KHÔNG dùng `golang.org/x/crypto` — đã kiểm chứng thư viện chuẩn
  `crypto/pbkdf2` cho ra kết quả khớp tuyệt đối với bản Java.
- Mọi ghi chú trong code viết **tiếng Việt**, giải thích *tại sao* chứ không phải *cái gì*.
- Số nguyên trên dây truyền theo **big-endian** (chuẩn Java DataOutputStream), trừ đúng một ngoại lệ ghi rõ ở Task 3.
- Độ dài khoá mã **phải ≤ 127 byte**. Client dùng biến `sbyte` làm con trỏ khoá; dài hơn 127 là tràn số.
- Chuỗi ký tự dùng **modified UTF-8 của Java**: 2 byte độ dài big-endian rồi tới nội dung.
- Định dạng mật khẩu phải khớp bản cũ: `pbkdf2$<số vòng>$<muối base64>$<băm base64>`, PBKDF2-HMAC-SHA256, khoá 256 bit.
- Mọi file nguồn đặt dưới `server/`. Không đụng vào `1-MAY-CHU-GAME/`.
- Mỗi task kết thúc bằng một commit.

---

## Đặc tả giao thức (đọc kỹ trước khi làm)

Trích xuất từ `D:\NRO_Project\Assets\Scripts\Game1\Session_ME.cs` và `1-MAY-CHU-GAME\src\server\ServerManager.java`.

### Bắt tay

```
1. Client mở TCP, gửi ngay:  [cmd = -27] [0x00] [0x00]        (chưa mã hoá)
2. Server trả lời cmd -27, thân tin (chưa mã hoá):
       byte    độ dài khoá (n ≤ 127)
       n byte  khoá thô
       UTF     địa chỉ server phụ  (chuỗi rỗng nếu không có)
       int     cổng server phụ     (0 nếu không có)
       byte    có server phụ không (0 = không)
3. Hai bên cùng biến đổi khoá thô:  for j = 0..n-2:  key[j+1] ^= key[j]
4. Từ đây mọi byte đều XOR theo luồng khoá, hai con trỏ RIÊNG BIỆT
   cho chiều đọc và chiều ghi.
```

### Luồng khoá

```
đọc(b)  = key[curR] XOR b ;  curR++ ;  nếu curR >= n thì curR %= n
ghi(b)  = key[curW] XOR b ;  curW++ ;  nếu curW >= n thì curW %= n
```

XOR đối xứng, nên cùng một hàm dùng được cho cả mã hoá lẫn giải mã — chỉ khác con trỏ.

### Client → Server (server phải ĐỌC)

```
[1 byte cmd]                         đã XOR nếu có khoá
[2 byte độ dài]                      nếu có khoá: XOR, thứ tự (len>>8) rồi (len & 0xFF)
                                     nếu chưa: little-endian (chỉ xảy ra ở gói -27, độ dài 0)
[độ dài byte thân tin]               mỗi byte XOR nếu có khoá
```

### Server → Client (server phải GHI)

```
[1 byte cmd]                         XOR nếu có khoá

Nếu cmd thuộc {-32, -66, -74, 11, -67, -87, 66}  → độ dài 3 byte:
    ghi  XOR(len)      + 128
    ghi  XOR(len>>8)   + 128
    ghi  XOR(len>>16)  + 128

Ngược lại nếu có khoá → độ dài 2 byte:
    ghi  XOR(len>>8)
    ghi  XOR(len & 0xFF)

Ngược lại → 2 byte big-endian không mã hoá.

[thân tin]                           mỗi byte XOR nếu có khoá
```

> Ghi chú: bản Java ghi `- 128`, client đọc `+ 128`. Trên byte 8 bit, cộng 128 và trừ 128 cho kết quả giống nhau, nên Go dùng `+ 128` là đúng.

---

## Cấu trúc file

| File | Trách nhiệm |
|---|---|
| `server/go.mod` | Khai báo module, phiên bản Go |
| `server/internal/protocol/keystream.go` | Sinh khoá, biến đổi khoá, XOR luồng |
| `server/internal/protocol/keystream_test.go` | Kiểm thử luồng khoá |
| `server/internal/protocol/datastream.go` | Đọc/ghi kiểu dữ liệu tương thích Java |
| `server/internal/protocol/datastream_test.go` | Kiểm thử kiểu dữ liệu |
| `server/internal/protocol/message.go` | Cấu trúc gói tin |
| `server/internal/protocol/frame.go` | Đọc/ghi khung tin trên dây |
| `server/internal/protocol/frame_test.go` | Kiểm thử khung tin |
| `server/internal/auth/password.go` | Băm và đối chiếu mật khẩu PBKDF2 |
| `server/internal/auth/password_test.go` | Kiểm thử mật khẩu |
| `server/internal/store/store.go` | Kết nối PostgreSQL, truy vấn tài khoản |
| `server/internal/store/migrations/0001_tai_khoan.sql` | Lược đồ tài khoản |
| `server/internal/net/session.go` | Vòng đời một kết nối |
| `server/internal/net/listener.go` | Lắng nghe TCP |
| `server/internal/game/login.go` | Nghiệp vụ đăng nhập |
| `server/cmd/nroserver/main.go` | Điểm khởi động |
| `server/internal/protocol/integration_test.go` | Kiểm thử đầu-cuối bằng client giả lập |

---

### Task 1: Khởi tạo dự án Go

**Files:**
- Create: `server/go.mod`
- Create: `server/cmd/nroserver/main.go`
- Create: `server/.gitignore`

**Interfaces:**
- Consumes: không
- Produces: module `nroserver`, chạy được bằng `go run ./cmd/nroserver`

- [ ] **Step 1: Tạo module**

```bash
cd D:/NRO_GoiBanGiao
mkdir -p server/cmd/nroserver
cd server
go mod init nroserver
go mod edit -go=1.24
```

- [ ] **Step 2: Viết điểm khởi động tạm**

Tạo `server/cmd/nroserver/main.go`:

```go
package main

import "fmt"

func main() {
	fmt.Println("nroserver: chua co gi de chay")
}
```

- [ ] **Step 3: Chạy để xác nhận môi trường Go hoạt động**

Run: `cd D:/NRO_GoiBanGiao/server && go run ./cmd/nroserver`
Expected: in ra `nroserver: chua co gi de chay`, thoát mã 0

- [ ] **Step 4: Tạo .gitignore riêng cho thư mục server**

Tạo `server/.gitignore`:

```
# File bien dich cua Go
/nroserver
/nroserver.exe
# Ket qua do do bao phu kiem thu
*.out
```

- [ ] **Step 5: Commit**

```bash
cd D:/NRO_GoiBanGiao
git add server/go.mod server/cmd/nroserver/main.go server/.gitignore
git commit -m "chore: khoi tao module Go cho may chu game"
```

---

### Task 2: Luồng khoá mã

**Files:**
- Create: `server/internal/protocol/keystream.go`
- Test: `server/internal/protocol/keystream_test.go`

**Interfaces:**
- Consumes: không
- Produces:
  - `func GenerateRawKey(n int) ([]byte, error)` — sinh khoá thô ngẫu nhiên, `n` phải trong khoảng 1..127
  - `func DeriveKey(raw []byte) []byte` — áp phép biến đổi `key[j+1] ^= key[j]`
  - `type KeyStream struct` với `func NewKeyStream(key []byte) *KeyStream`
  - `func (k *KeyStream) Read(b byte) byte`
  - `func (k *KeyStream) Write(b byte) byte`

- [ ] **Step 1: Viết kiểm thử thất bại**

Tạo `server/internal/protocol/keystream_test.go`:

```go
package protocol

import (
	"bytes"
	"testing"
)

// Khoa tho [1,2,3] sau bien doi phai thanh [1, 1^2, (1^2)^3] = [1, 3, 0].
func TestDeriveKey(t *testing.T) {
	got := DeriveKey([]byte{1, 2, 3})
	want := []byte{1, 3, 0}
	if !bytes.Equal(got, want) {
		t.Fatalf("DeriveKey = %v, mong doi %v", got, want)
	}
}

// DeriveKey khong duoc sua mang goc.
func TestDeriveKeyKhongSuaMangGoc(t *testing.T) {
	raw := []byte{1, 2, 3}
	_ = DeriveKey(raw)
	if !bytes.Equal(raw, []byte{1, 2, 3}) {
		t.Fatalf("mang goc bi sua: %v", raw)
	}
}

// XOR doi xung: ghi roi doc lai bang mot luong khac phai ra byte ban dau.
func TestGhiRoiDocRaByteBanDau(t *testing.T) {
	key := DeriveKey([]byte{9, 4, 7, 2})
	ben_ghi := NewKeyStream(key)
	ben_doc := NewKeyStream(key)

	goc := []byte{0, 1, 127, 128, 255, 42, 200, 13, 99}
	for i, b := range goc {
		maHoa := ben_ghi.Write(b)
		giaiMa := ben_doc.Read(maHoa)
		if giaiMa != b {
			t.Fatalf("byte thu %d: goc %d, giai ma ra %d", i, b, giaiMa)
		}
	}
}

// Con tro phai quay vong khi vuot qua do dai khoa.
func TestConTroQuayVong(t *testing.T) {
	key := []byte{0xAA, 0xBB}
	ks := NewKeyStream(key)
	if got := ks.Write(0x00); got != 0xAA {
		t.Fatalf("byte 1: %02X, mong doi AA", got)
	}
	if got := ks.Write(0x00); got != 0xBB {
		t.Fatalf("byte 2: %02X, mong doi BB", got)
	}
	if got := ks.Write(0x00); got != 0xAA {
		t.Fatalf("byte 3 phai quay ve dau khoa: %02X, mong doi AA", got)
	}
}

// Hai con tro doc va ghi phai doc lap nhau.
func TestHaiConTroDocLap(t *testing.T) {
	ks := NewKeyStream([]byte{0x11, 0x22})
	ks.Write(0x00) // con tro ghi -> 1
	if got := ks.Read(0x00); got != 0x11 {
		t.Fatalf("con tro doc bi anh huong boi ghi: %02X, mong doi 11", got)
	}
}

func TestGenerateRawKeyChanDoDaiSai(t *testing.T) {
	if _, err := GenerateRawKey(0); err == nil {
		t.Fatal("do dai 0 phai bao loi")
	}
	if _, err := GenerateRawKey(128); err == nil {
		t.Fatal("do dai 128 phai bao loi vi client dung sbyte lam con tro")
	}
	k, err := GenerateRawKey(8)
	if err != nil {
		t.Fatalf("do dai 8 phai hop le: %v", err)
	}
	if len(k) != 8 {
		t.Fatalf("do dai khoa = %d, mong doi 8", len(k))
	}
}
```

- [ ] **Step 2: Chạy để xác nhận thất bại**

Run: `cd D:/NRO_GoiBanGiao/server && go test ./internal/protocol/ -run TestDeriveKey -v`
Expected: FAIL — `undefined: DeriveKey`

- [ ] **Step 3: Viết cài đặt tối thiểu**

Tạo `server/internal/protocol/keystream.go`:

```go
// Package protocol cai dat giao thuc nhi phan cua client Unity cu.
//
// Toan bo phan "noi chuyen bang byte" nam gon trong package nay. Phan con lai
// cua may chu khong duoc biet gi ve khoa ma hay thu tu byte — nho vay sau nay
// them client moi chi can them mot package tuong tu, khong dung vao luat choi.
package protocol

import (
	"crypto/rand"
	"fmt"
)

// doDaiKhoaToiDa gioi han o 127 vi client dung kieu sbyte lam con tro khoa
// (xem Session_ME.cs, ham readKey/writeKey). Khoa dai hon se lam tran so
// ben phia client va giai ma sai — loi im lang, rat kho tim.
const doDaiKhoaToiDa = 127

// GenerateRawKey sinh khoa tho ngau nhien de gui cho client trong goi bat tay.
func GenerateRawKey(n int) ([]byte, error) {
	if n < 1 || n > doDaiKhoaToiDa {
		return nil, fmt.Errorf("do dai khoa phai trong khoang 1..%d, nhan duoc %d", doDaiKhoaToiDa, n)
	}
	raw := make([]byte, n)
	if _, err := rand.Read(raw); err != nil {
		return nil, fmt.Errorf("khong sinh duoc khoa ngau nhien: %w", err)
	}
	return raw, nil
}

// DeriveKey ap dung phep bien doi ma client thuc hien sau khi nhan khoa tho:
//
//	for j := 0; j < len-1; j++ { key[j+1] ^= key[j] }
//
// Hai ben phai ra cung ket qua, neu khong moi byte sau do deu sai.
// Ham tra ve mang moi, khong sua mang dau vao.
func DeriveKey(raw []byte) []byte {
	key := make([]byte, len(raw))
	copy(key, raw)
	for j := 0; j+1 < len(key); j++ {
		key[j+1] ^= key[j]
	}
	return key
}

// KeyStream giu hai con tro doc lap: mot cho chieu doc, mot cho chieu ghi.
// Dung chung mot con tro cho ca hai chieu se lech ngay goi tin dau tien.
type KeyStream struct {
	key  []byte
	curR int
	curW int
}

func NewKeyStream(key []byte) *KeyStream {
	return &KeyStream{key: key}
}

// Read giai ma mot byte nhan tu client.
func (k *KeyStream) Read(b byte) byte {
	out := k.key[k.curR] ^ b
	k.curR++
	if k.curR >= len(k.key) {
		k.curR %= len(k.key)
	}
	return out
}

// Write ma hoa mot byte truoc khi gui cho client.
func (k *KeyStream) Write(b byte) byte {
	out := k.key[k.curW] ^ b
	k.curW++
	if k.curW >= len(k.key) {
		k.curW %= len(k.key)
	}
	return out
}
```

- [ ] **Step 4: Chạy kiểm thử để xác nhận đạt**

Run: `cd D:/NRO_GoiBanGiao/server && go test ./internal/protocol/ -v`
Expected: PASS — cả 6 kiểm thử

- [ ] **Step 5: Commit**

```bash
cd D:/NRO_GoiBanGiao
git add server/internal/protocol/keystream.go server/internal/protocol/keystream_test.go
git commit -m "feat(protocol): luong khoa ma tuong thich client cu"
```

---

### Task 3: Đọc ghi kiểu dữ liệu tương thích Java

> ## ⚠️ ĐÍNH CHÍNH — code mẫu bên dưới CÓ LỖI, đã sửa khi thi công
>
> **Đoạn `func (w *Writer) UTF` trong kế hoạch này viết sai.** Nó ghi thẳng
> chuỗi Go (UTF-8 tiêu chuẩn), trong khi Java dùng biến thể gọi là
> **modified UTF-8**. Hai chỗ khác nhau:
>
> | Ký tự | UTF-8 chuẩn (code sai) | Modified UTF-8 (đúng) |
> |---|---|---|
> | NUL (U+0000) | 1 byte `00` | 2 byte `C0 80` |
> | Emoji (ngoài BMP) | 4 byte | **6 byte** — tách thành cặp surrogate |
>
> Hậu quả đo bằng Java thật: ghi emoji theo UTF-8 chuẩn làm
> `java.io.DataInputStream.readUTF` phía client **ném `UTFDataFormatException`**
> — gói tin coi như hỏng, không phải hiện sai chữ.
>
> Nguy hiểm ở chỗ: bộ test trong kế hoạch này **vẫn đạt hết** dù code sai, vì
> nó chỉ thử chuỗi ASCII và tiếng Việt có dấu — hai loại đó trùng nhau giữa
> hai cách mã hoá. Lỗi chỉ lộ khi người chơi gõ emoji vào tên hoặc chat.
>
> **Cài đặt đúng nằm ở `server/internal/protocol/datastream.go`** (hàm
> `encodeModifiedUTF8` / `decodeModifiedUTF8`). Ai đọc lại kế hoạch này thì
> lấy code trong repo, đừng chép code mẫu bên dưới.
>
> Hai chuỗi byte chuẩn đo từ Java thật, dùng làm mốc kiểm thử:
> ```
> writeUTF("A" + NUL + "B")  ->  00 04 41 C0 80 42
> writeUTF("X" + emoji + "Y") ->  00 08 58 ED A0 BC ED BE 89 59
> ```
>
> Bài học: **test đạt không có nghĩa là code đúng** — chỉ có nghĩa là code
> đúng với những gì test thử. Ràng buộc chung của kế hoạch đã ghi rõ
> "modified UTF-8 của Java" ngay từ đầu, nhưng cả code lẫn test đều không
> thực hiện đúng ràng buộc đó.

**Files:**
- Create: `server/internal/protocol/datastream.go`
- Test: `server/internal/protocol/datastream_test.go`

**Interfaces:**
- Consumes: không
- Produces:
  - `type Writer struct` với `NewWriter() *Writer`, các phương thức `Byte(v byte)`, `Bool(v bool)`, `Short(v int16)`, `Int(v int32)`, `Long(v int64)`, `UTF(s string) error`, `Bytes() []byte`
  - `type Reader struct` với `NewReader(b []byte) *Reader`, các phương thức `Byte() (byte, error)`, `Short() (int16, error)`, `Int() (int32, error)`, `UTF() (string, error)`

- [ ] **Step 1: Viết kiểm thử thất bại**

Tạo `server/internal/protocol/datastream_test.go`:

```go
package protocol

import (
	"bytes"
	"testing"
)

// Java DataOutputStream ghi so nguyen theo big-endian.
func TestGhiSoNguyenBigEndian(t *testing.T) {
	w := NewWriter()
	w.Short(0x0102)
	w.Int(0x03040506)
	got := w.Bytes()
	want := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	if !bytes.Equal(got, want) {
		t.Fatalf("got % X, mong doi % X", got, want)
	}
}

// writeUTF cua Java: 2 byte do dai big-endian roi toi noi dung.
func TestGhiChuoiKieuJava(t *testing.T) {
	w := NewWriter()
	if err := w.UTF("AB"); err != nil {
		t.Fatalf("loi khong mong doi: %v", err)
	}
	want := []byte{0x00, 0x02, 'A', 'B'}
	if !bytes.Equal(w.Bytes(), want) {
		t.Fatalf("got % X, mong doi % X", w.Bytes(), want)
	}
}

// Chuoi tieng Viet phai giu nguyen sau khi ghi roi doc lai.
func TestChuoiTiengVietVongTron(t *testing.T) {
	w := NewWriter()
	if err := w.UTF("Làng Lá"); err != nil {
		t.Fatalf("loi ghi: %v", err)
	}
	r := NewReader(w.Bytes())
	got, err := r.UTF()
	if err != nil {
		t.Fatalf("loi doc: %v", err)
	}
	if got != "Làng Lá" {
		t.Fatalf("got %q, mong doi %q", got, "Làng Lá")
	}
}

// Chuoi qua 65535 byte khong ghi duoc — phai bao loi thay vi ghi hong.
func TestChuoiQuaDaiBaoLoi(t *testing.T) {
	w := NewWriter()
	qua_dai := string(make([]byte, 70000))
	if err := w.UTF(qua_dai); err == nil {
		t.Fatal("chuoi 70000 byte phai bao loi")
	}
}

func TestDocSoNguyen(t *testing.T) {
	r := NewReader([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06})
	s, err := r.Short()
	if err != nil || s != 0x0102 {
		t.Fatalf("Short = %d, err = %v", s, err)
	}
	i, err := r.Int()
	if err != nil || i != 0x03040506 {
		t.Fatalf("Int = %d, err = %v", i, err)
	}
}

// Doc qua cuoi bo dem phai bao loi, khong duoc panic.
func TestDocQuaCuoiBaoLoi(t *testing.T) {
	r := NewReader([]byte{0x01})
	if _, err := r.Int(); err == nil {
		t.Fatal("doc Int tu 1 byte phai bao loi")
	}
}
```

- [ ] **Step 2: Chạy để xác nhận thất bại**

Run: `cd D:/NRO_GoiBanGiao/server && go test ./internal/protocol/ -run TestGhiSoNguyenBigEndian -v`
Expected: FAIL — `undefined: NewWriter`

- [ ] **Step 3: Viết cài đặt**

Tạo `server/internal/protocol/datastream.go`:

```go
package protocol

import (
	"encoding/binary"
	"fmt"
)

// Writer ghi kieu du lieu theo dung cach java.io.DataOutputStream ghi,
// vi client duoc chuyen the tu ban Java ME va van doc theo quy uoc do.
// Moi so nguyen deu big-endian.
type Writer struct {
	buf []byte
}

func NewWriter() *Writer { return &Writer{buf: make([]byte, 0, 64)} }

func (w *Writer) Byte(v byte) { w.buf = append(w.buf, v) }

func (w *Writer) Bool(v bool) {
	if v {
		w.buf = append(w.buf, 1)
	} else {
		w.buf = append(w.buf, 0)
	}
}

func (w *Writer) Short(v int16) {
	w.buf = binary.BigEndian.AppendUint16(w.buf, uint16(v))
}

func (w *Writer) Int(v int32) {
	w.buf = binary.BigEndian.AppendUint32(w.buf, uint32(v))
}

func (w *Writer) Long(v int64) {
	w.buf = binary.BigEndian.AppendUint64(w.buf, uint64(v))
}

// UTF ghi chuoi theo dinh dang writeUTF cua Java: 2 byte do dai roi noi dung.
// Go luu chuoi san o dang UTF-8 nen khong can chuyen doi, chi can kiem do dai.
func (w *Writer) UTF(s string) error {
	if len(s) > 65535 {
		return fmt.Errorf("chuoi dai %d byte, vuot gioi han 65535 cua writeUTF", len(s))
	}
	w.Short(int16(len(s)))
	w.buf = append(w.buf, s...)
	return nil
}

func (w *Writer) Bytes() []byte { return w.buf }

// Reader doc nguoc lai nhung gi Writer ghi.
type Reader struct {
	buf []byte
	pos int
}

func NewReader(b []byte) *Reader { return &Reader{buf: b} }

func (r *Reader) con(n int) ([]byte, error) {
	if r.pos+n > len(r.buf) {
		return nil, fmt.Errorf("doc %d byte tai vi tri %d nhung bo dem chi co %d byte", n, r.pos, len(r.buf))
	}
	out := r.buf[r.pos : r.pos+n]
	r.pos += n
	return out, nil
}

func (r *Reader) Byte() (byte, error) {
	b, err := r.con(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func (r *Reader) Short() (int16, error) {
	b, err := r.con(2)
	if err != nil {
		return 0, err
	}
	return int16(binary.BigEndian.Uint16(b)), nil
}

func (r *Reader) Int() (int32, error) {
	b, err := r.con(4)
	if err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(b)), nil
}

func (r *Reader) UTF() (string, error) {
	n, err := r.Short()
	if err != nil {
		return "", err
	}
	b, err := r.con(int(uint16(n)))
	if err != nil {
		return "", err
	}
	return string(b), nil
}
```

- [ ] **Step 4: Chạy kiểm thử**

Run: `cd D:/NRO_GoiBanGiao/server && go test ./internal/protocol/ -v`
Expected: PASS — tất cả kiểm thử của Task 2 và Task 3

- [ ] **Step 5: Commit**

```bash
cd D:/NRO_GoiBanGiao
git add server/internal/protocol/datastream.go server/internal/protocol/datastream_test.go
git commit -m "feat(protocol): doc ghi kieu du lieu tuong thich java DataStream"
```

---

### Task 4: Khung tin trên dây

> ## ⚠️ ĐÍNH CHÍNH THỨ HAI — code mẫu bên dưới CÓ LỖI NGHIÊM TRỌNG, đã sửa
>
> Trong `WriteMessage`, câu `switch` chọn nhánh độ dài 3 byte viết là:
>
> ```go
> case n > 0 && lenhDoDaiBaByte[m.Command]:   // SAI
> ```
>
> **Điều kiện `n > 0` là thừa và gây lệch dòng byte.** Bằng chứng lấy từ client
> (`Session_ME.cs`, hàm `readMessage`):
>
> ```csharp
> if (b == -32 || b == -66 || b == 11 || b == -67 || b == -74 || b == -87 || b == 66)
>     return readMessage2(b);   // readMessage2 LUÔN đọc 3 byte độ dài
> ```
>
> Client quyết định đọc 3 byte **chỉ dựa vào byte lệnh**, không xét độ dài thân
> tin. Nên nếu server gửi một trong 7 lệnh đó với thân tin rỗng: server ghi 2
> byte, client đọc 3 byte → **mọi gói tin sau đó đều đọc sai**. Triệu chứng:
> client treo im lặng, không có thông báo lỗi nào.
>
> **Sửa đúng:** bỏ `n > 0`, chỉ còn `case lenhDoDaiBaByte[m.Command]:`
>
> Điều đáng chú ý: bản Java cũ cũng có cùng lỗ hổng này cho trường hợp
> `data == null` (ghi `writeShort(0)` 2 byte không xét lệnh). Nó chưa nổ vì
> trên thực tế chưa ai gửi 7 lệnh đó với thân tin rỗng. Bản Go **sửa cho đúng
> hơn bản gốc**, thay vì chép nguyên lỗi.
>
> ### Và còn HAI lỗi nữa cùng loại, phát hiện ở vòng soi sau
>
> **Lỗi 2 — nguy hiểm hơn lỗi 1.** Nhánh còn lại cũng có cùng điều kiện thừa:
> ```go
> case n > 0 && f.coKhoa():   // SAI y hệt
> ```
> Lệnh **thường** gửi với thân tin rỗng **sau khi đã có khoá** sẽ ghi 2 byte
> thô, trong khi client vẫn XOR chúng → vừa ra độ dài rác, vừa **đẩy lệch con
> trỏ khoá vĩnh viễn**. Nguy hiểm hơn lỗi 1 vì "lệnh thường gửi rỗng" là
> thông báo và xác nhận — xảy ra liên tục khi chơi, còn 7 lệnh đặc biệt thì hiếm.
>
> **Lỗi 3 — sập cả máy chủ.** Nhánh lệnh đặc biệt gọi `f.ks.Write(...)` không
> kiểm `coKhoa()`. Gọi trước `SetKey` là panic con trỏ rỗng, mà panic không
> `recover()` trong Go **làm sập cả tiến trình** — đá văng toàn bộ người chơi
> đang online, không chỉ một kết nối.
>
> ### Nguyên tắc đúng, áp dụng cho CẢ BA nhánh
>
> | Quyết định | Dựa vào |
> |---|---|
> | Mã hoá hay không | **chỉ** `coKhoa()` |
> | 2 byte hay 3 byte độ dài | **chỉ** byte lệnh |
> | Độ dài thân tin | **không tham gia quyết định nào** |
>
> Bằng chứng dứt điểm cho nguyên tắc này: **chính client, ở hàm
> `doSendMessage` nhánh `data == null`, cũng mã hoá 2 byte độ dài khi đã có
> khoá** (`Session_ME.cs`, gọi `writeKey((sbyte)0)` hai lần). Xác nhận từ cả
> chiều đọc lẫn chiều gửi của client.
>
> Cài đặt đúng nằm ở `server/internal/protocol/frame.go`. Bộ test bổ sung:
> `TestGoiDacBietThanRongKhongLamLechGoiSau`, `TestLenhThuongThanRongVanMaHoaDoDai`,
> `TestLenhDacBietChuaCoKhoaTraLoiKhongPanic` — hai cái đầu ghi hai gói liên
> tiếp rồi đọc lại, bắt được hiện tượng lệch dòng chứ không chỉ đếm độ dài.
>
> **Bài học lớn nhất của Task 4:** sửa xong lỗi 1 mà không soi nhánh còn lại
> nên bỏ sót lỗi 2 — cùng một kiểu sai, cách đó vài dòng. Khi sửa một lỗi,
> phải tìm xem cùng kiểu sai đó còn ở đâu nữa.

**Files:**
- Create: `server/internal/protocol/message.go`
- Create: `server/internal/protocol/frame.go`
- Test: `server/internal/protocol/frame_test.go`

**Interfaces:**
- Consumes: `KeyStream` (Task 2)
- Produces:
  - `type Message struct { Command int8; Data []byte }`
  - `type Framer struct` với `NewFramer(r io.Reader, w io.Writer) *Framer`
  - `func (f *Framer) SetKey(key []byte)` — bật mã hoá sau bắt tay
  - `func (f *Framer) ReadMessage() (*Message, error)`
  - `func (f *Framer) WriteMessage(m *Message) error`

- [ ] **Step 1: Viết kiểm thử thất bại**

Tạo `server/internal/protocol/frame_test.go`:

```go
package protocol

import (
	"bytes"
	"testing"
)

// Truoc khi bat tay xong, goi tin di nguyen ban khong ma hoa.
func TestDocGoiBatTayChuaMaHoa(t *testing.T) {
	// Client gui: cmd = -27, do dai 0 (hai byte 0)
	vao := bytes.NewReader([]byte{0xE5, 0x00, 0x00})
	f := NewFramer(vao, &bytes.Buffer{})
	m, err := f.ReadMessage()
	if err != nil {
		t.Fatalf("loi doc: %v", err)
	}
	if m.Command != -27 {
		t.Fatalf("Command = %d, mong doi -27", m.Command)
	}
	if len(m.Data) != 0 {
		t.Fatalf("Data dai %d, mong doi 0", len(m.Data))
	}
}

// Sau khi co khoa, ghi roi doc lai bang mot Framer khac phai ra goi ban dau.
func TestVongTronSauKhiCoKhoa(t *testing.T) {
	raw := []byte{5, 9, 13, 21}
	var day bytes.Buffer

	ben_gui := NewFramer(&bytes.Buffer{}, &day)
	ben_gui.SetKey(DeriveKey(raw))

	goc := &Message{Command: -29, Data: []byte("xin chao")}
	if err := ben_gui.WriteMessage(goc); err != nil {
		t.Fatalf("loi ghi: %v", err)
	}

	// Ben nhan dung KeyStream rieng nhung cung khoa. Vi ben gui dung con tro
	// ghi con ben nhan dung con tro doc, ta phai gia lap dung chieu: tao Framer
	// moi va doc bang con tro DOC.
	ben_nhan := NewFramer(bytes.NewReader(day.Bytes()), &bytes.Buffer{})
	ben_nhan.SetKey(DeriveKey(raw))

	nhan, err := ben_nhan.ReadMessage()
	if err != nil {
		t.Fatalf("loi doc: %v", err)
	}
	if nhan.Command != goc.Command {
		t.Fatalf("Command = %d, mong doi %d", nhan.Command, goc.Command)
	}
	if !bytes.Equal(nhan.Data, goc.Data) {
		t.Fatalf("Data = %q, mong doi %q", nhan.Data, goc.Data)
	}
}

// Cac lenh trong danh sach dac biet dung do dai 3 byte thay vi 2.
func TestLenhDacBietDungDoDaiBaByte(t *testing.T) {
	var day bytes.Buffer
	f := NewFramer(&bytes.Buffer{}, &day)
	f.SetKey(DeriveKey([]byte{1, 2, 3, 4}))

	// -87 nam trong danh sach dac biet
	if err := f.WriteMessage(&Message{Command: -87, Data: []byte{0xAA}}); err != nil {
		t.Fatalf("loi ghi: %v", err)
	}
	// 1 byte lenh + 3 byte do dai + 1 byte than tin = 5
	if day.Len() != 5 {
		t.Fatalf("goi dai %d byte, mong doi 5 (lenh dac biet dung 3 byte do dai)", day.Len())
	}
}

// Lenh thuong dung do dai 2 byte.
func TestLenhThuongDungDoDaiHaiByte(t *testing.T) {
	var day bytes.Buffer
	f := NewFramer(&bytes.Buffer{}, &day)
	f.SetKey(DeriveKey([]byte{1, 2, 3, 4}))

	if err := f.WriteMessage(&Message{Command: -29, Data: []byte{0xAA}}); err != nil {
		t.Fatalf("loi ghi: %v", err)
	}
	if day.Len() != 4 {
		t.Fatalf("goi dai %d byte, mong doi 4", day.Len())
	}
}

// Goi khong co than tin van phai ghi 2 byte do dai bang 0.
func TestGoiRongVanGhiDoDai(t *testing.T) {
	var day bytes.Buffer
	f := NewFramer(&bytes.Buffer{}, &day)
	if err := f.WriteMessage(&Message{Command: -27}); err != nil {
		t.Fatalf("loi ghi: %v", err)
	}
	if day.Len() != 3 {
		t.Fatalf("goi dai %d byte, mong doi 3", day.Len())
	}
}

// Goi qua lon phai bi tu choi de tranh ke xau bao do dai khong lo lam can RAM.
func TestTuChoiGoiQuaLon(t *testing.T) {
	// cmd -29, do dai 0xFFFF, nhung khong co du lieu theo sau
	vao := bytes.NewReader([]byte{0xE3, 0xFF, 0xFF})
	f := NewFramer(vao, &bytes.Buffer{})
	if _, err := f.ReadMessage(); err == nil {
		t.Fatal("goi bao do dai 65535 ma khong co du lieu phai bao loi")
	}
}
```

- [ ] **Step 2: Chạy để xác nhận thất bại**

Run: `cd D:/NRO_GoiBanGiao/server && go test ./internal/protocol/ -run TestDocGoiBatTay -v`
Expected: FAIL — `undefined: NewFramer`

- [ ] **Step 3: Viết cấu trúc gói tin**

Tạo `server/internal/protocol/message.go`:

```go
package protocol

// Message la mot goi tin tren day: mot byte lenh va phan than tin.
// Than tin co the rong (nil hoac do dai 0).
type Message struct {
	Command int8
	Data    []byte
}

// NewMessage tao goi tin voi than tin duoc dung san bang Writer.
func NewMessage(cmd int8, w *Writer) *Message {
	return &Message{Command: cmd, Data: w.Bytes()}
}
```

- [ ] **Step 4: Viết bộ đọc ghi khung tin**

Tạo `server/internal/protocol/frame.go`:

```go
package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

// doDaiThanToiDa chan ke xau bao mot do dai khong lo de ep may chu cap phat
// het bo nho. Goi lon nhat that su trong game la du lieu vat pham, duoi 1 MB.
const doDaiThanToiDa = 4 << 20 // 4 MB

// lenhDoDaiBaByte la cac lenh dung 3 byte do dai thay vi 2.
// Danh sach lay tu ServerManager.java (ham doSendMessage) va Session_ME.cs
// (ham readMessage). Hai ben phai khop tuyet doi, lech mot lenh la lech ca dong.
var lenhDoDaiBaByte = map[int8]bool{
	-32: true, -66: true, -74: true, 11: true, -67: true, -87: true, 66: true,
}

// Framer doc va ghi goi tin tren mot ket noi.
// Truoc khi goi SetKey thi moi thu di nguyen ban — do la giai doan bat tay.
type Framer struct {
	r  io.Reader
	w  io.Writer
	ks *KeyStream
}

func NewFramer(r io.Reader, w io.Writer) *Framer {
	return &Framer{r: r, w: w}
}

// SetKey bat che do ma hoa. Goi dung mot lan, ngay sau khi gui goi bat tay.
func (f *Framer) SetKey(key []byte) {
	f.ks = NewKeyStream(key)
}

func (f *Framer) coKhoa() bool { return f.ks != nil }

func (f *Framer) docMotByte() (byte, error) {
	var b [1]byte
	if _, err := io.ReadFull(f.r, b[:]); err != nil {
		return 0, err
	}
	return b[0], nil
}

// ReadMessage doc mot goi tin do client gui len.
//
// Dinh dang (xem Session_ME.cs, ham doSendMessage):
//
//	[1 byte lenh]  [2 byte do dai]  [than tin]
//
// Khi da co khoa, do dai ghi theo thu tu (len>>8) roi (len & 0xFF), moi byte
// deu XOR. Khi chua co khoa, client dung BinaryWriter.Write(ushort) cua .NET
// — day la little-endian. Truong hop nay chi xay ra o goi bat tay -27 voi do
// dai bang 0, nen hai byte deu la 0 va thu tu khong anh huong.
func (f *Framer) ReadMessage() (*Message, error) {
	rawCmd, err := f.docMotByte()
	if err != nil {
		return nil, err
	}
	if f.coKhoa() {
		rawCmd = f.ks.Read(rawCmd)
	}
	cmd := int8(rawCmd)

	b1, err := f.docMotByte()
	if err != nil {
		return nil, err
	}
	b2, err := f.docMotByte()
	if err != nil {
		return nil, err
	}

	var n int
	if f.coKhoa() {
		hi := f.ks.Read(b1)
		lo := f.ks.Read(b2)
		n = int(hi)<<8 | int(lo)
	} else {
		n = int(binary.LittleEndian.Uint16([]byte{b1, b2}))
	}

	if n < 0 || n > doDaiThanToiDa {
		return nil, fmt.Errorf("goi tin bao do dai %d, vuot gioi han %d", n, doDaiThanToiDa)
	}

	data := make([]byte, n)
	if n > 0 {
		if _, err := io.ReadFull(f.r, data); err != nil {
			return nil, fmt.Errorf("doc than tin %d byte that bai: %w", n, err)
		}
		if f.coKhoa() {
			for i := range data {
				data[i] = f.ks.Read(data[i])
			}
		}
	}
	return &Message{Command: cmd, Data: data}, nil
}

// WriteMessage gui mot goi tin xuong client.
//
// Dinh dang (xem ServerManager.java, ham doSendMessage): lenh, do dai, than tin.
// Do dai dung 3 byte cho cac lenh trong lenhDoDaiBaByte, 2 byte cho phan con lai.
//
// Ban Java ghi "XOR(x) - 128" con client doc "readKey(b) + 128". Tren byte 8 bit,
// cong 128 va tru 128 cho cung ket qua, nen o day dung cong cho de doc.
func (f *Framer) WriteMessage(m *Message) error {
	out := make([]byte, 0, len(m.Data)+8)

	if f.coKhoa() {
		out = append(out, f.ks.Write(byte(m.Command)))
	} else {
		out = append(out, byte(m.Command))
	}

	n := len(m.Data)
	switch {
	case n > 0 && lenhDoDaiBaByte[m.Command]:
		out = append(out, f.ks.Write(byte(n))+128)
		out = append(out, f.ks.Write(byte(n>>8))+128)
		out = append(out, f.ks.Write(byte(n>>16))+128)
	case n > 0 && f.coKhoa():
		out = append(out, f.ks.Write(byte(n>>8)))
		out = append(out, f.ks.Write(byte(n&0xFF)))
	default:
		// Chua co khoa, hoac than tin rong: ghi 2 byte big-endian khong ma hoa.
		out = binary.BigEndian.AppendUint16(out, uint16(n))
	}

	if n > 0 {
		than := make([]byte, n)
		copy(than, m.Data)
		if f.coKhoa() {
			for i := range than {
				than[i] = f.ks.Write(than[i])
			}
		}
		out = append(out, than...)
	}

	_, err := f.w.Write(out)
	return err
}
```

- [ ] **Step 5: Chạy kiểm thử**

Run: `cd D:/NRO_GoiBanGiao/server && go test ./internal/protocol/ -v`
Expected: PASS — tất cả kiểm thử

> Nếu `TestVongTronSauKhiCoKhoa` thất bại: nguyên nhân gần như chắc chắn là hai con trỏ đọc/ghi bị dùng lẫn. Bên gửi phải dùng con trỏ **ghi**, bên nhận dùng con trỏ **đọc**, và mỗi bên có `KeyStream` riêng.

- [ ] **Step 6: Commit**

```bash
cd D:/NRO_GoiBanGiao
git add server/internal/protocol/message.go server/internal/protocol/frame.go server/internal/protocol/frame_test.go
git commit -m "feat(protocol): doc ghi khung tin tren day"
```

---

### Task 5: Đối chiếu mật khẩu tương thích bản cũ

**Files:**
- Create: `server/internal/auth/password.go`
- Test: `server/internal/auth/password_test.go`

**Interfaces:**
- Consumes: không
- Produces:
  - `func Hash(matKhau string) (string, error)` — sinh chuỗi `pbkdf2$100000$<muối>$<băm>`
  - `func Verify(matKhau, luuTru string) bool` — đối chiếu, hỗ trợ cả dạng cũ chưa băm
  - `func IsHashed(luuTru string) bool`

- [ ] **Step 1: Viết kiểm thử thất bại**

Tạo `server/internal/auth/password_test.go`:

```go
package auth

import "testing"

// Chuoi bam do may chu Java cu sinh ra phai doi chieu duoc.
// Gia tri nay lay tu tai khoan that trong CSDL `z`, mat khau la "1".
const bamTuBanJava = "pbkdf2$100000$5Q53KdKJ0c559SGZQOh2fw==$CtbK9GXUDmUYhDEg4AOdFM3cLT7s08mEPGgCoyUOAzU="

func TestDoiChieuBamCuaBanJava(t *testing.T) {
	if !Verify("1", bamTuBanJava) {
		t.Fatal("mat khau dung nhung doi chieu that bai — dinh dang bam khong khop ban Java")
	}
	if Verify("2", bamTuBanJava) {
		t.Fatal("mat khau sai nhung van cho qua")
	}
}

func TestBamRoiDoiChieu(t *testing.T) {
	h, err := Hash("matkhau-thu")
	if err != nil {
		t.Fatalf("loi bam: %v", err)
	}
	if !Verify("matkhau-thu", h) {
		t.Fatal("bam roi doi chieu lai that bai")
	}
	if Verify("matkhau-khac", h) {
		t.Fatal("mat khau khac ma van cho qua")
	}
}

// Hai lan bam cung mot mat khau phai ra hai chuoi khac nhau (muoi ngau nhien).
func TestMuoiNgauNhien(t *testing.T) {
	a, _ := Hash("abc")
	b, _ := Hash("abc")
	if a == b {
		t.Fatal("hai lan bam ra ket qua giong nhau — muoi khong ngau nhien")
	}
}

// Tai khoan cu luu mat khau tho van phai dang nhap duoc.
func TestTaiKhoanCuLuuTho(t *testing.T) {
	if !Verify("123456", "123456") {
		t.Fatal("tai khoan cu luu mat khau tho phai dang nhap duoc")
	}
	if Verify("sai", "123456") {
		t.Fatal("mat khau tho sai ma van cho qua")
	}
}

func TestNhanBietDangBam(t *testing.T) {
	if !IsHashed(bamTuBanJava) {
		t.Fatal("chuoi bam phai duoc nhan biet")
	}
	if IsHashed("123456") {
		t.Fatal("mat khau tho khong phai dang bam")
	}
}

// Chuoi bam hong khong duoc lam sap may chu.
func TestChuoiBamHong(t *testing.T) {
	for _, x := range []string{"pbkdf2$", "pbkdf2$abc$x$y", "pbkdf2$100000$!!!$!!!", ""} {
		if Verify("bat ky", x) {
			t.Fatalf("chuoi hong %q ma van cho qua", x)
		}
	}
}
```

- [ ] **Step 2: Chạy để xác nhận thất bại**

Run: `cd D:/NRO_GoiBanGiao/server && go test ./internal/auth/ -v`
Expected: FAIL — `undefined: Verify`

- [ ] **Step 3: Viết cài đặt**

Tạo `server/internal/auth/password.go`:

```go
// Package auth lo viec bam va doi chieu mat khau.
//
// Dinh dang PHAI khop tuyet doi voi ban Java cu (auth/PasswordUtil.java) va
// ban PHP (sv1/config.php). Sua mot ben ma quen ben kia se khien toan bo
// nguoi choi khong dang nhap duoc.
package auth

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

const (
	tienTo   = "pbkdf2$"
	soVong   = 100000
	soByteMu = 16
	soBitKey = 256
)

// IsHashed cho biet chuoi luu trong CSDL da o dang bam hay con la mat khau tho
// cua he thong cu.
func IsHashed(luuTru string) bool {
	return strings.HasPrefix(luuTru, tienTo)
}

// Hash bam mat khau moi voi muoi ngau nhien.
func Hash(matKhau string) (string, error) {
	muoi := make([]byte, soByteMu)
	if _, err := rand.Read(muoi); err != nil {
		return "", fmt.Errorf("khong sinh duoc muoi: %w", err)
	}
	// Chu y chu ky ham cua thu vien CHUAN khac voi golang.org/x/crypto:
	//   pbkdf2.Key(ham_bam, mat_khau_string, muoi, so_vong, so_byte) (bytes, error)
	bam, err := pbkdf2.Key(sha256.New, matKhau, muoi, soVong, soBitKey/8)
	if err != nil {
		return "", fmt.Errorf("bam that bai: %w", err)
	}
	return fmt.Sprintf("%s%d$%s$%s",
		tienTo, soVong,
		base64.StdEncoding.EncodeToString(muoi),
		base64.StdEncoding.EncodeToString(bam)), nil
}

// Verify doi chieu mat khau nguoi dung nhap voi gia tri luu trong CSDL.
// Ho tro ca tai khoan cu con luu mat khau tho — nhung van so sanh bang
// thuat toan hang thoi gian de khong lo thong tin qua do tre.
func Verify(matKhau, luuTru string) bool {
	if luuTru == "" {
		return false
	}
	if !IsHashed(luuTru) {
		return subtle.ConstantTimeCompare([]byte(luuTru), []byte(matKhau)) == 1
	}

	phan := strings.Split(luuTru, "$")
	if len(phan) != 4 {
		return false
	}
	vong, err := strconv.Atoi(phan[1])
	if err != nil || vong < 1 {
		return false
	}
	muoi, err := base64.StdEncoding.DecodeString(phan[2])
	if err != nil {
		return false
	}
	muon, err := base64.StdEncoding.DecodeString(phan[3])
	if err != nil || len(muon) == 0 {
		return false
	}

	duoc, err := pbkdf2.Key(sha256.New, matKhau, muoi, vong, len(muon))
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(muon, duoc) == 1
}
```

- [ ] **Step 4: Không cần thêm thư viện nào**

`crypto/pbkdf2` nằm trong thư viện chuẩn từ Go 1.24. Đã kiểm chứng trên máy này rằng nó cho ra **đúng chuỗi băm mà bản Java sinh ra** với cùng dữ liệu vào.

Run: `cd D:/NRO_GoiBanGiao/server && go build ./...`
Expected: không lỗi, `go.mod` **không** có `golang.org/x/crypto`

- [ ] **Step 5: Chạy kiểm thử**

Run: `cd D:/NRO_GoiBanGiao/server && go test ./internal/auth/ -v`
Expected: PASS — cả 6 kiểm thử. Đặc biệt `TestDoiChieuBamCuaBanJava` phải đạt: đó là bằng chứng bản Go đọc được mật khẩu do bản Java sinh ra.

- [ ] **Step 6: Commit**

```bash
cd D:/NRO_GoiBanGiao
git add server/internal/auth/ server/go.mod server/go.sum
git commit -m "feat(auth): doi chieu mat khau PBKDF2 tuong thich ban Java"
```

---

### Task 6: Lược đồ tài khoản và kho dữ liệu

**Files:**
- Create: `server/internal/store/migrations/0001_tai_khoan.sql`
- Create: `server/internal/store/store.go`
- Test: `server/internal/store/store_test.go`

**Interfaces:**
- Consumes: không
- Produces:
  - `type TaiKhoan struct { ID int64; TenDangNhap string; MatKhauBam string; BiKhoa bool; LaAdmin bool }`
  - `type Store struct` với `func Open(ctx context.Context, dsn string) (*Store, error)`
  - `func (s *Store) LayTaiKhoan(ctx context.Context, ten string) (*TaiKhoan, error)` — trả `ErrKhongThayTaiKhoan` nếu không có
  - `var ErrKhongThayTaiKhoan = errors.New(...)`

- [ ] **Step 1: Viết lược đồ**

Tạo `server/internal/store/migrations/0001_tai_khoan.sql`:

```sql
-- Bang tai khoan. Co tinh dat toi thieu o buoc nay: ke hoach 2 se mo rong
-- day du va chuyen du lieu cu sang.
CREATE TABLE IF NOT EXISTS tai_khoan (
    id            BIGSERIAL PRIMARY KEY,
    ten_dang_nhap TEXT        NOT NULL,
    mat_khau_bam  TEXT        NOT NULL,
    bi_khoa       BOOLEAN     NOT NULL DEFAULT FALSE,
    la_admin      BOOLEAN     NOT NULL DEFAULT FALSE,
    tao_luc       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Ten dang nhap khong phan biet hoa thuong: nguoi choi go "Abc" hay "abc"
-- deu phai vao dung mot tai khoan. Chi muc duy nhat chan viec tao trung.
CREATE UNIQUE INDEX IF NOT EXISTS tai_khoan_ten_duy_nhat
    ON tai_khoan (lower(ten_dang_nhap));
```

- [ ] **Step 2: Viết kiểm thử thất bại**

Tạo `server/internal/store/store_test.go`:

```go
package store

import (
	"context"
	"errors"
	"os"
	"testing"
)

// Kiem thu nay can mot CSDL PostgreSQL that. Dat bien moi truong:
//
//	NRO_TEST_DSN=postgres://postgres:nrodev123@127.0.0.1:5433/nro_test
//
// Khong dat thi bo qua, de nguoi khac chay `go test ./...` khong bi do.
func moKetNoiThu(t *testing.T) *Store {
	t.Helper()
	dsn := os.Getenv("NRO_TEST_DSN")
	if dsn == "" {
		t.Skip("bo qua: chua dat NRO_TEST_DSN")
	}
	s, err := Open(context.Background(), dsn)
	if err != nil {
		t.Fatalf("khong mo duoc CSDL: %v", err)
	}
	t.Cleanup(s.Close)
	return s
}

func TestLayTaiKhoanKhongTonTai(t *testing.T) {
	s := moKetNoiThu(t)
	_, err := s.LayTaiKhoan(context.Background(), "khong-he-ton-tai-xyz")
	if !errors.Is(err, ErrKhongThayTaiKhoan) {
		t.Fatalf("err = %v, mong doi ErrKhongThayTaiKhoan", err)
	}
}

func TestLayTaiKhoanKhongPhanBietHoaThuong(t *testing.T) {
	s := moKetNoiThu(t)
	ctx := context.Background()

	_, err := s.pool.Exec(ctx,
		`INSERT INTO tai_khoan (ten_dang_nhap, mat_khau_bam) VALUES ($1, $2)
		 ON CONFLICT DO NOTHING`, "NguoiThu", "bam-gia")
	if err != nil {
		t.Fatalf("khong chen duoc du lieu thu: %v", err)
	}
	t.Cleanup(func() {
		s.pool.Exec(ctx, `DELETE FROM tai_khoan WHERE lower(ten_dang_nhap) = 'nguoithu'`)
	})

	for _, ten := range []string{"NguoiThu", "nguoithu", "NGUOITHU"} {
		tk, err := s.LayTaiKhoan(ctx, ten)
		if err != nil {
			t.Fatalf("tim %q that bai: %v", ten, err)
		}
		if tk.MatKhauBam != "bam-gia" {
			t.Fatalf("lay nham tai khoan cho %q", ten)
		}
	}
}
```

- [ ] **Step 3: Chạy để xác nhận thất bại**

Run: `cd D:/NRO_GoiBanGiao/server && go test ./internal/store/ -v`
Expected: FAIL — `undefined: Open`

- [ ] **Step 4: Viết cài đặt**

Tạo `server/internal/store/store.go`:

```go
// Package store lo viec noi chuyen voi PostgreSQL.
package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrKhongThayTaiKhoan = errors.New("khong tim thay tai khoan")

type TaiKhoan struct {
	ID          int64
	TenDangNhap string
	MatKhauBam  string
	BiKhoa      bool
	LaAdmin     bool
}

type Store struct {
	pool *pgxpool.Pool
}

func Open(ctx context.Context, dsn string) (*Store, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("tao pool that bai: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("khong ket noi duoc CSDL: %w", err)
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close() { s.pool.Close() }

// LayTaiKhoan tim theo ten dang nhap, khong phan biet hoa thuong.
func (s *Store) LayTaiKhoan(ctx context.Context, ten string) (*TaiKhoan, error) {
	var tk TaiKhoan
	err := s.pool.QueryRow(ctx,
		`SELECT id, ten_dang_nhap, mat_khau_bam, bi_khoa, la_admin
		   FROM tai_khoan
		  WHERE lower(ten_dang_nhap) = lower($1)
		  LIMIT 1`, ten,
	).Scan(&tk.ID, &tk.TenDangNhap, &tk.MatKhauBam, &tk.BiKhoa, &tk.LaAdmin)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrKhongThayTaiKhoan
	}
	if err != nil {
		return nil, fmt.Errorf("truy van tai khoan that bai: %w", err)
	}
	return &tk, nil
}
```

- [ ] **Step 5: Cài thư viện và tạo CSDL thử**

Run:
```bash
cd D:/NRO_GoiBanGiao/server
go get github.com/jackc/pgx/v5@latest
go mod tidy
```

Tạo CSDL thử và chạy lược đồ (thay `matkhau` bằng mật khẩu Postgres của bạn):
```bash
# CSDL nro_test da tao san tren may nay
D:\PostgreSQL\bin\psql.exe -h 127.0.0.1 -p 5433 -U postgres -d nro_test -f internal/store/migrations/0001_tai_khoan.sql
```

- [ ] **Step 6: Chạy kiểm thử**

Run:
```bash
cd D:/NRO_GoiBanGiao/server
set NRO_TEST_DSN=postgres://postgres:nrodev123@127.0.0.1:5433/nro_test
go test ./internal/store/ -v
```
Expected: PASS — 2 kiểm thử. Nếu chưa đặt `NRO_TEST_DSN` thì thấy SKIP, cũng chấp nhận được nhưng phải chạy được ít nhất một lần với CSDL thật trước khi qua task sau.

- [ ] **Step 7: Commit**

```bash
cd D:/NRO_GoiBanGiao
git add server/internal/store/ server/go.mod server/go.sum
git commit -m "feat(store): ket noi PostgreSQL va truy van tai khoan"
```

---

### Task 7: Phiên kết nối và bắt tay

**Files:**
- Create: `server/internal/net/session.go`
- Create: `server/internal/net/listener.go`
- Test: `server/internal/net/session_test.go`

**Interfaces:**
- Consumes: `protocol.Framer`, `protocol.GenerateRawKey`, `protocol.DeriveKey`, `protocol.Writer`, `protocol.Message` (Task 2–4)
- Produces:
  - `type XuLyGoiTin interface { XuLy(s *Session, m *protocol.Message) error }`
  - `type Session struct` với `func (s *Session) Gui(m *protocol.Message) error`, trường `DiaChiIP string`
  - `func PhucVu(conn net.Conn, xuLy XuLyGoiTin)` — chạy vòng đời một kết nối
  - `func LangNghe(diaChi string, xuLy XuLyGoiTin) error`

- [ ] **Step 1: Viết kiểm thử thất bại**

Tạo `server/internal/net/session_test.go`:

```go
package net

import (
	"net"
	"testing"
	"time"

	"nroserver/internal/protocol"
)

type xuLyGhiNhan struct {
	nhanDuoc chan *protocol.Message
}

func (x *xuLyGhiNhan) XuLy(s *Session, m *protocol.Message) error {
	x.nhanDuoc <- m
	return nil
}

// Kich ban day du: client gui -27, may chu tra khoa, roi ca hai noi chuyen
// bang goi da ma hoa.
func TestBatTayRoiNhanGoiDaMaHoa(t *testing.T) {
	ben_client, ben_server := net.Pipe()
	defer ben_client.Close()

	xuLy := &xuLyGhiNhan{nhanDuoc: make(chan *protocol.Message, 4)}
	go PhucVu(ben_server, xuLy)

	fClient := protocol.NewFramer(ben_client, ben_client)

	// 1. Client gui goi bat tay
	if err := fClient.WriteMessage(&protocol.Message{Command: -27}); err != nil {
		t.Fatalf("gui goi bat tay that bai: %v", err)
	}

	// 2. Client nhan khoa
	ben_client.SetReadDeadline(time.Now().Add(3 * time.Second))
	traLoi, err := fClient.ReadMessage()
	if err != nil {
		t.Fatalf("khong nhan duoc goi khoa: %v", err)
	}
	if traLoi.Command != -27 {
		t.Fatalf("lenh tra loi = %d, mong doi -27", traLoi.Command)
	}

	r := protocol.NewReader(traLoi.Data)
	doDai, err := r.Byte()
	if err != nil {
		t.Fatalf("doc do dai khoa that bai: %v", err)
	}
	if doDai < 1 || doDai > 127 {
		t.Fatalf("do dai khoa = %d, phai trong 1..127", doDai)
	}
	khoaTho := make([]byte, doDai)
	for i := range khoaTho {
		khoaTho[i], _ = r.Byte()
	}
	if _, err := r.UTF(); err != nil {
		t.Fatalf("thieu truong dia chi server phu: %v", err)
	}
	if _, err := r.Int(); err != nil {
		t.Fatalf("thieu truong cong server phu: %v", err)
	}
	if _, err := r.Byte(); err != nil {
		t.Fatalf("thieu truong co server phu: %v", err)
	}

	// 3. Ca hai bat ma hoa
	fClient.SetKey(protocol.DeriveKey(khoaTho))

	// 4. Client gui goi da ma hoa, may chu phai giai ma dung
	if err := fClient.WriteMessage(&protocol.Message{Command: -29, Data: []byte("test")}); err != nil {
		t.Fatalf("gui goi da ma hoa that bai: %v", err)
	}

	select {
	case m := <-xuLy.nhanDuoc:
		if m.Command != -29 {
			t.Fatalf("lenh = %d, mong doi -29", m.Command)
		}
		if string(m.Data) != "test" {
			t.Fatalf("than tin = %q, mong doi \"test\"", m.Data)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("may chu khong giai ma duoc goi sau bat tay")
	}
}

// Goi dau tien khong phai -27 thi phai bi ngat ket noi.
func TestGoiDauKhongPhaiBatTayThiNgat(t *testing.T) {
	ben_client, ben_server := net.Pipe()
	defer ben_client.Close()

	go PhucVu(ben_server, &xuLyGhiNhan{nhanDuoc: make(chan *protocol.Message, 1)})

	fClient := protocol.NewFramer(ben_client, ben_client)
	if err := fClient.WriteMessage(&protocol.Message{Command: 99}); err != nil {
		t.Fatalf("gui that bai: %v", err)
	}

	ben_client.SetReadDeadline(time.Now().Add(2 * time.Second))
	if _, err := fClient.ReadMessage(); err == nil {
		t.Fatal("may chu phai ngat ket noi khi goi dau khong phai -27")
	}
}
```

- [ ] **Step 2: Chạy để xác nhận thất bại**

Run: `cd D:/NRO_GoiBanGiao/server && go test ./internal/net/ -v`
Expected: FAIL — `undefined: PhucVu`

- [ ] **Step 3: Viết phiên kết nối**

Tạo `server/internal/net/session.go`:

```go
// Package net quan ly vong doi mot ket noi TCP.
//
// Package nay khong biet gi ve luat choi. No chi lo: bat tay, doc goi tin,
// va giao goi tin cho tang tren xu ly.
package net

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"nroserver/internal/protocol"
)

const (
	// doDaiKhoa 8 byte la du. Client gioi han 127, nhung khoa dai khong lam
	// an toan hon dang ke — day chi la lop xao tron, khong phai ma hoa that.
	doDaiKhoa = 8

	// thoiGianChoBatTay: ke tan cong mo hang loat ket noi roi khong gui gi
	// se an het bo nho. Khong bat tay trong 5 giay thi ngat.
	thoiGianChoBatTay = 5 * time.Second

	// thoiGianChoGoiTin: ket noi im lang qua lau coi nhu da chet.
	thoiGianChoGoiTin = 120 * time.Second
)

// XuLyGoiTin la thu tang tren cai dat de nhan goi tin da giai ma.
type XuLyGoiTin interface {
	XuLy(s *Session, m *protocol.Message) error
}

// Session la mot ket noi dang song.
type Session struct {
	conn     net.Conn
	framer   *protocol.Framer
	DiaChiIP string

	// muGui bao dam hai goroutine khong ghi chen nhau vao cung mot ket noi.
	// Luong khoa ma co con tro chay tuan tu — ghi chen nhau la lech khoa
	// vinh vien, va trieu chung la client nhan rac chu khong bao loi.
	muGui sync.Mutex

	dongMotLan sync.Once
}

// Gui gui mot goi tin xuong client. An toan khi goi tu nhieu goroutine.
func (s *Session) Gui(m *protocol.Message) error {
	s.muGui.Lock()
	defer s.muGui.Unlock()
	return s.framer.WriteMessage(m)
}

func (s *Session) Dong() {
	s.dongMotLan.Do(func() { s.conn.Close() })
}

// PhucVu chay vong doi mot ket noi cho toi khi dong.
func PhucVu(conn net.Conn, xuLy XuLyGoiTin) {
	diaChi := "khong-ro"
	if a, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		diaChi = a.IP.String()
	}

	s := &Session{
		conn:     conn,
		framer:   protocol.NewFramer(conn, conn),
		DiaChiIP: diaChi,
	}
	defer s.Dong()

	if err := batTay(s); err != nil {
		slog.Warn("bat tay that bai", "ip", diaChi, "loi", err)
		return
	}

	for {
		_ = conn.SetReadDeadline(time.Now().Add(thoiGianChoGoiTin))
		m, err := s.framer.ReadMessage()
		if err != nil {
			slog.Debug("ket noi ket thuc", "ip", diaChi, "loi", err)
			return
		}
		if err := xuLy.XuLy(s, m); err != nil {
			slog.Warn("xu ly goi tin that bai", "ip", diaChi, "lenh", m.Command, "loi", err)
			return
		}
	}
}

// batTay thuc hien buoc trao khoa. Client mo ket noi la gui ngay goi -27
// rong; may chu tra ve khoa tho roi ca hai bat ma hoa.
func batTay(s *Session) error {
	_ = s.conn.SetReadDeadline(time.Now().Add(thoiGianChoBatTay))

	dau, err := s.framer.ReadMessage()
	if err != nil {
		return fmt.Errorf("khong doc duoc goi dau tien: %w", err)
	}
	if dau.Command != -27 {
		return fmt.Errorf("goi dau tien co lenh %d, mong doi -27", dau.Command)
	}

	khoaTho, err := protocol.GenerateRawKey(doDaiKhoa)
	if err != nil {
		return fmt.Errorf("sinh khoa that bai: %w", err)
	}

	w := protocol.NewWriter()
	w.Byte(byte(len(khoaTho)))
	for _, b := range khoaTho {
		w.Byte(b)
	}
	// Ba truong sau la dia chi may chu phu. Du an nay chi chay mot may chu,
	// nhung client van doc du ba truong — thieu la no ket o vong doc.
	if err := w.UTF(""); err != nil {
		return err
	}
	w.Int(0)
	w.Byte(0)

	if err := s.Gui(protocol.NewMessage(-27, w)); err != nil {
		return fmt.Errorf("gui khoa that bai: %w", err)
	}

	// Bat ma hoa SAU khi da gui goi khoa: chinh goi do phai di nguyen ban.
	s.framer.SetKey(protocol.DeriveKey(khoaTho))
	return nil
}
```

- [ ] **Step 4: Viết bộ lắng nghe**

Tạo `server/internal/net/listener.go`:

```go
package net

import (
	"fmt"
	"log/slog"
	"net"
)

// LangNghe mo cong TCP va phuc vu tung ket noi tren mot goroutine rieng.
func LangNghe(diaChi string, xuLy XuLyGoiTin) error {
	ln, err := net.Listen("tcp", diaChi)
	if err != nil {
		return fmt.Errorf("khong mo duoc cong %s: %w", diaChi, err)
	}
	defer ln.Close()

	slog.Info("may chu game dang lang nghe", "dia_chi", diaChi)

	for {
		conn, err := ln.Accept()
		if err != nil {
			slog.Error("nhan ket noi that bai", "loi", err)
			continue
		}
		go PhucVu(conn, xuLy)
	}
}
```

- [ ] **Step 5: Chạy kiểm thử**

Run: `cd D:/NRO_GoiBanGiao/server && go test ./internal/net/ -v`
Expected: PASS — 2 kiểm thử

- [ ] **Step 6: Commit**

```bash
cd D:/NRO_GoiBanGiao
git add server/internal/net/
git commit -m "feat(net): vong doi ket noi va bat tay trao khoa"
```

---

### Task 8: Nghiệp vụ đăng nhập

**Files:**
- Create: `server/internal/game/login.go`
- Test: `server/internal/game/login_test.go`

**Interfaces:**
- Consumes: `net.Session`, `net.XuLyGoiTin`, `protocol.Message`, `auth.Verify`, `store.TaiKhoan`, `store.ErrKhongThayTaiKhoan`
- Produces:
  - `type KhoTaiKhoan interface { LayTaiKhoan(ctx context.Context, ten string) (*store.TaiKhoan, error) }`
  - `type BoXuLy struct` với `func NewBoXuLy(kho KhoTaiKhoan) *BoXuLy`
  - `func (b *BoXuLy) XuLy(s *net.Session, m *protocol.Message) error`

- [ ] **Step 1: Viết kiểm thử thất bại**

Tạo `server/internal/game/login_test.go`:

```go
package game

import (
	"context"
	"testing"

	"nroserver/internal/auth"
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

type khoGia struct {
	tk  *store.TaiKhoan
	err error
}

func (k *khoGia) LayTaiKhoan(ctx context.Context, ten string) (*store.TaiKhoan, error) {
	if k.err != nil {
		return nil, k.err
	}
	return k.tk, nil
}

// dungGoiDangNhap dung goi dang nhap dung y het client gui:
// khung lenh -29, ben trong la byte lenh con 0 roi hai chuoi.
func dungGoiDangNhap(ten, matKhau string) *protocol.Message {
	w := protocol.NewWriter()
	w.Byte(0) // lenh con: dang nhap
	w.UTF(ten)
	w.UTF(matKhau)
	return protocol.NewMessage(-29, w)
}

// Goi dung dinh dang phai duoc boc tach ra dung ten va mat khau.
func TestBocTachGoiDangNhap(t *testing.T) {
	m := dungGoiDangNhap("nguoi", "matkhau")
	r := protocol.NewReader(m.Data)

	lenhCon, err := r.Byte()
	if err != nil || lenhCon != 0 {
		t.Fatalf("lenh con = %d, err = %v, mong doi 0", lenhCon, err)
	}
	ten, _ := r.UTF()
	mk, _ := r.UTF()
	if ten != "nguoi" || mk != "matkhau" {
		t.Fatalf("boc tach ra %q / %q, mong doi \"nguoi\" / \"matkhau\"", ten, mk)
	}
}

func TestDangNhapDungMatKhau(t *testing.T) {
	bam, _ := auth.Hash("matkhau")
	b := NewBoXuLy(&khoGia{tk: &store.TaiKhoan{ID: 7, TenDangNhap: "nguoi", MatKhauBam: bam}})

	kq, err := b.KiemTraDangNhap(context.Background(), "nguoi", "matkhau")
	if err != nil {
		t.Fatalf("loi khong mong doi: %v", err)
	}
	if !kq.ThanhCong {
		t.Fatalf("phai dang nhap duoc, ly do tu choi: %q", kq.LyDo)
	}
	if kq.TaiKhoan.ID != 7 {
		t.Fatalf("ID = %d, mong doi 7", kq.TaiKhoan.ID)
	}
}

func TestDangNhapSaiMatKhau(t *testing.T) {
	bam, _ := auth.Hash("matkhau")
	b := NewBoXuLy(&khoGia{tk: &store.TaiKhoan{ID: 7, MatKhauBam: bam}})

	kq, _ := b.KiemTraDangNhap(context.Background(), "nguoi", "sai")
	if kq.ThanhCong {
		t.Fatal("mat khau sai ma van cho vao")
	}
	if kq.LyDo != LyDoSaiThongTin {
		t.Fatalf("ly do = %q, mong doi %q", kq.LyDo, LyDoSaiThongTin)
	}
}

// Tai khoan khong ton tai va sai mat khau phai tra ve CUNG MOT cau bao,
// de ke xau khong do duoc tai khoan nao co that.
func TestTaiKhoanKhongTonTaiBaoGiongSaiMatKhau(t *testing.T) {
	b := NewBoXuLy(&khoGia{err: store.ErrKhongThayTaiKhoan})

	kq, _ := b.KiemTraDangNhap(context.Background(), "khong-co", "gi-do")
	if kq.ThanhCong {
		t.Fatal("tai khoan khong ton tai ma van cho vao")
	}
	if kq.LyDo != LyDoSaiThongTin {
		t.Fatalf("ly do = %q, phai giong het truong hop sai mat khau", kq.LyDo)
	}
}

func TestTaiKhoanBiKhoa(t *testing.T) {
	bam, _ := auth.Hash("matkhau")
	b := NewBoXuLy(&khoGia{tk: &store.TaiKhoan{ID: 7, MatKhauBam: bam, BiKhoa: true}})

	kq, _ := b.KiemTraDangNhap(context.Background(), "nguoi", "matkhau")
	if kq.ThanhCong {
		t.Fatal("tai khoan bi khoa ma van cho vao")
	}
	if kq.LyDo != LyDoBiKhoa {
		t.Fatalf("ly do = %q, mong doi %q", kq.LyDo, LyDoBiKhoa)
	}
}

// Ten dang nhap rong khong duoc cham toi CSDL.
func TestTenRongBiTuChoiNgay(t *testing.T) {
	b := NewBoXuLy(&khoGia{err: context.Canceled}) // kho se loi neu bi goi
	kq, err := b.KiemTraDangNhap(context.Background(), "", "gi-do")
	if err != nil {
		t.Fatalf("khong duoc goi toi kho: %v", err)
	}
	if kq.ThanhCong {
		t.Fatal("ten rong ma van cho vao")
	}
}
```

- [ ] **Step 2: Chạy để xác nhận thất bại**

Run: `cd D:/NRO_GoiBanGiao/server && go test ./internal/game/ -v`
Expected: FAIL — `undefined: NewBoXuLy`

- [ ] **Step 3: Viết cài đặt**

Tạo `server/internal/game/login.go`:

```go
// Package game lo nghiep vu choi. O ke hoach 1 moi chi co dang nhap.
package game

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"nroserver/internal/auth"
	mynet "nroserver/internal/net"
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// Cac cau tu choi. Sai mat khau va tai khoan khong ton tai dung CHUNG mot cau:
// bao khac nhau la lo cho ke xau biet tai khoan nao co that, tu do tap trung
// do mat khau dung tai khoan do. Ban Java cu cung lam vay, giu nguyen.
const (
	LyDoSaiThongTin = "Sai tài khoản hoặc mật khẩu."
	LyDoBiKhoa      = "Tài khoản đang bị khóa. Liên hệ quản trị viên."
)

// Lenh trong giao thuc. Cac gia tri nay lay tu ban Java cu, DUNG DOI:
//   -29  khung "chua dang nhap" — Controller.java:733, ben trong con mot
//        byte lenh con nua (0 = dang nhap) — Controller.java:876
//   -26  hien hop thoai co nut OK — Service.sendThongBaoOK
//  -102  bao client thoat trang thai "dang cho" — Service.sendLoginFail
//     2  chuyen sang man tao nhan vat — Service.switchToCreateChar
const (
	lenhChuaDangNhap  int8 = -29
	lenhConDangNhap   byte = 0
	lenhThongBaoOK    int8 = -26
	lenhDangNhapThatBai int8 = -102
	lenhTaoNhanVat    int8 = 2
)

type KetQuaDangNhap struct {
	ThanhCong bool
	LyDo      string
	TaiKhoan  *store.TaiKhoan
}

// KhoTaiKhoan tach ra thanh giao dien de kiem thu khong can CSDL that.
type KhoTaiKhoan interface {
	LayTaiKhoan(ctx context.Context, ten string) (*store.TaiKhoan, error)
}

type BoXuLy struct {
	kho KhoTaiKhoan
}

func NewBoXuLy(kho KhoTaiKhoan) *BoXuLy { return &BoXuLy{kho: kho} }

// KiemTraDangNhap chua toan bo luat dang nhap, tach khoi phan mang de kiem thu.
func (b *BoXuLy) KiemTraDangNhap(ctx context.Context, ten, matKhau string) (KetQuaDangNhap, error) {
	ten = strings.TrimSpace(ten)
	if ten == "" || matKhau == "" {
		return KetQuaDangNhap{LyDo: LyDoSaiThongTin}, nil
	}

	tk, err := b.kho.LayTaiKhoan(ctx, ten)
	if errors.Is(err, store.ErrKhongThayTaiKhoan) {
		return KetQuaDangNhap{LyDo: LyDoSaiThongTin}, nil
	}
	if err != nil {
		return KetQuaDangNhap{}, fmt.Errorf("tra cuu tai khoan that bai: %w", err)
	}

	if !auth.Verify(matKhau, tk.MatKhauBam) {
		return KetQuaDangNhap{LyDo: LyDoSaiThongTin}, nil
	}
	if tk.BiKhoa {
		return KetQuaDangNhap{LyDo: LyDoBiKhoa}, nil
	}
	return KetQuaDangNhap{ThanhCong: true, TaiKhoan: tk}, nil
}

// XuLy nhan goi tin da giai ma tu tang mang.
//
// Goi dang nhap long hai tang: khung -29, ben trong la byte lenh con roi
// moi toi du lieu. Doc thieu mot tang la boc tach ra rac.
func (b *BoXuLy) XuLy(s *mynet.Session, m *protocol.Message) error {
	if m.Command != lenhChuaDangNhap {
		// Ke hoach 1 chi lam dang nhap. Cac lenh khac bo qua, khong ngat
		// ket noi — client gui vai goi phu ngay sau khi noi.
		slog.Debug("bo qua lenh chua ho tro", "lenh", m.Command)
		return nil
	}

	r := protocol.NewReader(m.Data)
	lenhCon, err := r.Byte()
	if err != nil {
		return fmt.Errorf("goi -29 thieu byte lenh con: %w", err)
	}
	if lenhCon != lenhConDangNhap {
		slog.Debug("bo qua lenh con chua ho tro", "lenh_con", lenhCon)
		return nil
	}

	ten, err := r.UTF()
	if err != nil {
		return fmt.Errorf("goi dang nhap thieu ten: %w", err)
	}
	matKhau, err := r.UTF()
	if err != nil {
		return fmt.Errorf("goi dang nhap thieu mat khau: %w", err)
	}

	kq, err := b.KiemTraDangNhap(context.Background(), ten, matKhau)
	if err != nil {
		// Loi he thong: bao chung chung, khong lo chi tiet ky thuat ra ngoai.
		slog.Error("loi he thong khi dang nhap", "ip", s.DiaChiIP, "loi", err)
		return guiTuChoi(s, "Không đăng nhập được, vui lòng thử lại sau.")
	}
	if !kq.ThanhCong {
		slog.Info("dang nhap that bai", "ip", s.DiaChiIP, "ten", ten, "ly_do", kq.LyDo)
		return guiTuChoi(s, kq.LyDo)
	}

	slog.Info("dang nhap thanh cong",
		"ip", s.DiaChiIP, "ten", kq.TaiKhoan.TenDangNhap, "id", kq.TaiKhoan.ID)
	return guiChuyenSangTaoNhanVat(s)
}

// guiTuChoi bao loi roi bao client thoat khoi trang thai "dang cho".
//
// Ban Java cu tung QUEN buoc thu hai (goi -102), khien client treo o man hinh
// "Xin cho" cho toi khi tu thoat. Day la loi da tung xay ra that va da phai
// sua — giu ca hai buoc de khong lap lai.
func guiTuChoi(s *mynet.Session, cau string) error {
	w := protocol.NewWriter()
	if err := w.UTF(cau); err != nil {
		return err
	}
	if err := s.Gui(protocol.NewMessage(lenhThongBaoOK, w)); err != nil {
		return err
	}
	// Than tin 1 byte: 0 = khong con dang trong tien trinh dang nhap.
	return s.Gui(&protocol.Message{Command: lenhDangNhapThatBai, Data: []byte{0}})
}

// guiChuyenSangTaoNhanVat bao client mo man hinh tao nhan vat (lenh 2, than rong).
// Ke hoach 2 se thay bang: neu tai khoan da co nhan vat thi vao thang the gioi.
func guiChuyenSangTaoNhanVat(s *mynet.Session) error {
	return s.Gui(&protocol.Message{Command: lenhTaoNhanVat})
}
```

- [ ] **Step 4: Chạy kiểm thử**

Run: `cd D:/NRO_GoiBanGiao/server && go test ./internal/game/ -v`
Expected: PASS — 5 kiểm thử

- [ ] **Step 5: Commit**

```bash
cd D:/NRO_GoiBanGiao
git add server/internal/game/
git commit -m "feat(game): nghiep vu dang nhap"
```

---

### Task 9: Ghép lại và chạy thật

**Files:**
- Modify: `server/cmd/nroserver/main.go`
- Create: `server/internal/protocol/integration_test.go`
- Create: `server/README.md`

**Interfaces:**
- Consumes: tất cả các task trước
- Produces: chương trình chạy được, đọc cấu hình từ biến môi trường `NRO_DSN` và `NRO_DIA_CHI`

- [ ] **Step 1: Viết kiểm thử đầu-cuối**

Tạo `server/internal/protocol/integration_test.go`:

```go
package protocol

import (
	"bytes"
	"testing"
)

// Kiem thu nay bat loi nguy hiem nhat: hai chieu ma hoa lech nhau.
// No gia lap dung trinh tu that — may chu ghi, client doc — bang cach dung
// hai KeyStream rieng tren cung mot khoa, giong het hai dau ket noi that.
func TestMayChuGhiClientDocDungThuTu(t *testing.T) {
	khoaTho := []byte{0x3A, 0x7F, 0x01, 0xC4, 0x55, 0x09, 0xE1, 0x88}
	khoa := DeriveKey(khoaTho)

	var day bytes.Buffer

	// Ben may chu: chi ghi
	mayChu := NewFramer(&bytes.Buffer{}, &day)
	mayChu.SetKey(khoa)

	// Gui nhieu goi lien tiep de bat loi lech con tro tich luy dan.
	goiGui := []*Message{
		{Command: -26, Data: []byte("Sai tài khoản hoặc mật khẩu.")},
		{Command: -30},
		{Command: -87, Data: bytes.Repeat([]byte{0xAB}, 300)}, // lenh do dai 3 byte
		{Command: 0, Data: []byte{1, 2, 3}},
	}
	for i, m := range goiGui {
		if err := mayChu.WriteMessage(m); err != nil {
			t.Fatalf("goi %d: loi ghi: %v", i, err)
		}
	}

	// Ben client: chi doc
	client := NewFramer(bytes.NewReader(day.Bytes()), &bytes.Buffer{})
	client.SetKey(khoa)

	for i, mong := range goiGui {
		nhan, err := client.ReadMessage()
		if err != nil {
			t.Fatalf("goi %d: loi doc: %v", i, err)
		}
		if nhan.Command != mong.Command {
			t.Fatalf("goi %d: lenh = %d, mong doi %d", i, nhan.Command, mong.Command)
		}
		if !bytes.Equal(nhan.Data, mong.Data) && !(len(nhan.Data) == 0 && len(mong.Data) == 0) {
			t.Fatalf("goi %d: than tin dai %d, mong doi %d", i, len(nhan.Data), len(mong.Data))
		}
	}
}
```

- [ ] **Step 2: Chạy kiểm thử**

Run: `cd D:/NRO_GoiBanGiao/server && go test ./internal/protocol/ -run TestMayChuGhiClientDoc -v`
Expected: PASS

> Nếu thất bại ở gói thứ 3 (lệnh `-87`): kiểm tra lại nhánh độ dài 3 byte trong `WriteMessage` — thứ tự phải là `len`, `len>>8`, `len>>16`, và cả ba đều cộng 128.

- [ ] **Step 3: Viết điểm khởi động thật**

Thay toàn bộ nội dung `server/cmd/nroserver/main.go`:

```go
package main

import (
	"context"
	"log/slog"
	"os"

	"nroserver/internal/game"
	mynet "nroserver/internal/net"
	"nroserver/internal/store"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	dsn := os.Getenv("NRO_DSN")
	if dsn == "" {
		slog.Error("thieu bien moi truong NRO_DSN (chuoi ket noi PostgreSQL)")
		os.Exit(1)
	}

	diaChi := os.Getenv("NRO_DIA_CHI")
	if diaChi == "" {
		diaChi = ":14445"
	}

	ctx := context.Background()
	kho, err := store.Open(ctx, dsn)
	if err != nil {
		slog.Error("khong mo duoc CSDL", "loi", err)
		os.Exit(1)
	}
	defer kho.Close()

	if err := mynet.LangNghe(diaChi, game.NewBoXuLy(kho)); err != nil {
		slog.Error("may chu dung", "loi", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 4: Viết hướng dẫn chạy**

Tạo `server/README.md`:

```markdown
# Máy chủ game (Go)

## Chạy lần đầu

1. Tạo cơ sở dữ liệu và nạp lược đồ:

```bash
# CSDL nro da tao san tren may nay
D:\PostgreSQL\bin\psql.exe -h 127.0.0.1 -p 5433 -U postgres -d nro -f internal/store/migrations/0001_tai_khoan.sql
```

2. Tạo một tài khoản để thử (mật khẩu là `1`, dùng chuỗi băm sẵn):

```sql
INSERT INTO tai_khoan (ten_dang_nhap, mat_khau_bam) VALUES
('thu', 'pbkdf2$100000$5Q53KdKJ0c559SGZQOh2fw==$CtbK9GXUDmUYhDEg4AOdFM3cLT7s08mEPGgCoyUOAzU=');
```

3. Chạy:

```bash
set NRO_DSN=postgres://postgres:nrodev123@127.0.0.1:5433/nro
set NRO_DIA_CHI=:14445
go run ./cmd/nroserver
```

## Chạy toàn bộ kiểm thử

```bash
go test ./... -v
```

Kiểm thử trong `internal/store` cần biến `NRO_TEST_DSN`; không đặt thì bỏ qua.
```

- [ ] **Step 5: Chạy máy chủ và kiểm chứng bằng client Unity thật**

Run:
```bash
cd D:/NRO_GoiBanGiao/server
set NRO_DSN=postgres://postgres:nrodev123@127.0.0.1:5433/nro
go run ./cmd/nroserver
```
Expected: in ra `may chu game dang lang nghe dia_chi=:14445`

Sau đó mở client (bản PC ở `D:\NRO_Build\PC`), đăng nhập bằng `thu` / `1`.

Expected trong log máy chủ:
```
level=INFO msg="dang nhap thanh cong" ip=127.0.0.1 ten=thu id=1
```

Thử tiếp với mật khẩu sai. Expected:
```
level=INFO msg="dang nhap that bai" ... ly_do="Sai tài khoản hoặc mật khẩu."
```
Và client phải **hiện thông báo rồi thoát khỏi trạng thái chờ**, không treo ở màn "Xin chờ".

> Đây là bước kiểm chứng quan trọng nhất của cả kế hoạch. Nếu client kết nối và đăng nhập được, nghĩa là giao thức đã dựng lại đúng và toàn bộ hướng đi của dự án được xác nhận.

- [ ] **Step 6: Chạy toàn bộ kiểm thử lần cuối**

Run: `cd D:/NRO_GoiBanGiao/server && go test ./... -v`
Expected: PASS toàn bộ (`internal/store` có thể SKIP nếu chưa đặt `NRO_TEST_DSN`)

- [ ] **Step 7: Commit**

```bash
cd D:/NRO_GoiBanGiao
git add server/cmd/nroserver/main.go server/internal/protocol/integration_test.go server/README.md
git commit -m "feat: may chu chay duoc, client Unity that dang nhap thanh cong"
```

---

## Tiêu chí hoàn thành Kế hoạch 1

- [ ] `go test ./...` đạt toàn bộ
- [ ] Client Unity **thật** kết nối được vào máy chủ Go
- [ ] Đăng nhập đúng mật khẩu → vào được màn tạo nhân vật
- [ ] Đăng nhập sai → hiện thông báo, **không treo** ở màn chờ
- [ ] Mật khẩu do máy chủ Java cũ băm ra vẫn đăng nhập được (`TestDoiChieuBamCuaBanJava`)
- [ ] Tài khoản không tồn tại và sai mật khẩu trả về **cùng một câu báo**

## Việc CHƯA làm trong kế hoạch này (thuộc kế hoạch sau)

Bảng nhân vật và tạo nhân vật · Mô hình khu vực và di chuyển · Chống hack vị trí · Nạp nội dung · Chiến đấu · Giới hạn kết nối theo IP và chống DDoS · Chuyển dữ liệu từ CSDL cũ
