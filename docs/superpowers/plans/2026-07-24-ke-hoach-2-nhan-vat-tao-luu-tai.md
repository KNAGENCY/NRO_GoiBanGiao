# Kế hoạch 2 — Nhân vật: tạo, lưu, tải

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Người chơi tạo được nhân vật Naruto (tên, làng, kiểu tóc), nhân vật lưu bền vững trong PostgreSQL, và khi đăng nhập lại máy chủ nhận ra tài khoản đã có nhân vật.

**Architecture:** Nối tiếp Kế hoạch 1. Kế hoạch 1 đưa người chơi tới màn tạo nhân vật (gửi lệnh `2`); Kế hoạch 2 nhận gói tạo nhân vật (`-28`/lệnh con `2`), kiểm tên, ghi vào bảng `nhan_vat`, và tại lúc đăng nhập rẽ nhánh: chưa có nhân vật → màn tạo; đã có → vào thế giới. **"Vào thế giới" ở Kế hoạch 2 chỉ là điểm hẹn (stub)** — gửi đủ chuỗi gói để client rời màn chờ là việc của Kế hoạch 3.

**Tech Stack:** Go 1.25+ (máy có 1.26.5), PostgreSQL 16.4, `github.com/jackc/pgx/v5` (gồm `pgconn` để nhận diện lỗi trùng), thư viện chuẩn `regexp`, `testing`.

## Môi trường đã dựng sẵn

| Thứ | Vị trí / giá trị |
|---|---|
| Go | `D:\Go\bin\go.exe` — 1.26.5 |
| gofmt | `D:\Go\bin\gofmt.exe` |
| PostgreSQL | `D:\PostgreSQL\bin\` — cổng **5433** |
| Người dùng CSDL | `postgres` / `nrodev123` |
| CSDL chạy thật | `nro` |
| CSDL kiểm thử | `nro_test` |
| Chuỗi kết nối thật | `postgres://postgres:nrodev123@127.0.0.1:5433/nro` |
| Chuỗi kiểm thử | `postgres://postgres:nrodev123@127.0.0.1:5433/nro_test` (biến `NRO_TEST_DSN`) |
| psql | `D:\PostgreSQL\bin\psql.exe -h 127.0.0.1 -p 5433 -U postgres` (mật khẩu `nrodev123`) |

## Global Constraints

- **Go 1.25 trở lên.** Không dùng thư viện ngoài, **ngoại trừ `github.com/jackc/pgx/v5`** (đã có trong `go.mod`). `pgconn` là gói con của pgx, được phép.
- Mọi ghi chú trong code viết **tiếng Việt**, giải thích *tại sao* chứ không phải *cái gì*.
- Mọi file nguồn đặt dưới `server/`. **Không đụng `1-MAY-CHU-GAME/`.**
- Tên bảng và cột đặt **tiếng Việt không dấu**, theo đúng lối bảng `tai_khoan` ở Kế hoạch 1.
- Số nguyên trên dây truyền **big-endian** (chuẩn Java), chuỗi dùng **modified UTF-8** — dùng lại `protocol.Reader`/`protocol.Writer` đã có, KHÔNG tự đọc/ghi byte.
- **Không bao giờ tin dữ liệu client gửi lên.** Tên nhân vật phải kiểm ở tầng Go trước khi chạm CSDL (giống cách Kế hoạch 1 kiểm tên đăng nhập).
- **Không ghi nhật ký cho gói/lệnh-con không xử lý** — kể cả mức Debug. Đây là lỗi ngập nhật ký 9 GB/giờ đã xảy ra ở Kế hoạch 1 (xem ghi chú trong `login.go`).
- Mỗi task kết thúc bằng một commit.

---

## Đặc tả giao thức (đã giải mã từ client + máy chủ Java)

Nguồn: `D:\NRO_Project\Assets\Scripts\Game1\Service.cs`, `CreateCharScr.cs`, `Controller.cs`; và `1-MAY-CHU-GAME\src\server\Controller.java`, `jdbc\daos\PlayerDAO.java`.

### Gói tạo nhân vật (client → server)

```
Lệnh ngoài: -28   (hộp thư "không trên bản đồ" — nhiều nghiệp vụ, phân biệt bằng byte con đầu)
Thân tin:
    byte    lệnh con = 2        (luôn là 2 cho "tạo nhân vật")
    UTF     tên nhân vật
    byte    làng                (client cũ gửi 0/1/2; biến gốc tên "gender" nhưng THỰC CHẤT là chủng tộc/làng)
    byte    kiểu tóc            (ID sprite tóc, client cũ lấy từ mảng cứng hairIDOld[làng][chỉ_số])
```
Nguồn: `Service.cs:901-916` (`createChar`).

### Máy chủ trả lời

- **Thất bại** (tên sai/trùng/cấm): lệnh **`-26`**, thân = UTF câu thông báo. Client hiện hộp thoại và tự tắt vòng xoay "Xin chờ" (`InfoDlg.hide()`). Nguồn: `Controller.cs:2963-2986`.
- **Thành công**: máy chủ Java chạy lại toàn bộ luồng đăng nhập → vào thế giới. **Kế hoạch 2 KHÔNG làm phần vào thế giới** (đó là Kế hoạch 3); chỉ ghi nhận nhân vật vào phiên và dừng.

### Ràng buộc tên (khớp máy chủ Java)

- Độ dài **5–10** ký tự (máy chủ Java: `Controller.java:1043`). Client cho gõ tới 15 nhưng máy chủ cũ chỉ nhận 5–10 — giữ đúng 5–10 để tương thích.
- Chỉ **chữ cái tiếng Anh và chữ số** (`Util.haveSpecialCharacter` chặn mọi thứ ngoài `[a-zA-Z0-9]`).
- Không trùng tên đã có (không phân biệt hoa thường).
- Không nằm trong danh sách cấm.

### Điểm quan trọng cho người viết

> Trường `gender` trong gói và cột `gender` bảng cũ **là làng/chủng tộc, KHÔNG phải giới tính sinh học**. Đặt tên biến Go là `lang` để khỏi nhầm.

> Máy chủ Java hiện chỉ hỗ trợ **một nhân vật mỗi tài khoản** (`ChooseCharScr` là code chết). Kế hoạch 2 theo đúng mô hình 1-tài-khoản-1-nhân-vật.

> Client cũ cứng **đúng 3 làng** (`hairIDOld` là mảng 3×3), nên hiện chỉ chọn được làng 0/1/2. Nhưng giao thức đọc số làng động (`new NClass[d.readByte()]`), và thiết kế §7.2 có **5 làng**. Máy chủ Kế hoạch 2 định nghĩa sẵn cả 5 làng; làng 3–4 sẽ dùng được khi client được dựng lại.

---

## Cấu trúc file

**Tạo mới:**
- `server/internal/store/migrations/0002_nhan_vat.sql` — bảng `nhan_vat`
- `server/internal/game/lang.go` — danh mục 5 làng (chỉ số khởi đầu). Kế hoạch 6 sẽ chuyển sang YAML.
- `server/internal/game/lang_test.go`
- `server/internal/game/nhanvat.go` — xử lý gói tạo nhân vật + "vào thế giới" (stub)
- `server/internal/game/nhanvat_test.go`

**Sửa:**
- `server/internal/net/session.go` — `PhucVu` nhận **hàm tạo** bộ xử lý (mỗi kết nối một bộ)
- `server/internal/net/listener.go` — `LangNghe` nhận hàm tạo
- `server/internal/net/session_test.go` — cập nhật lời gọi
- `server/cmd/nroserver/main.go` — truyền hàm tạo `BoXuLy`
- `server/internal/store/store.go` — kiểu `NhanVat`/`NhanVatMoi`, các hàm CRUD, lỗi sentinel
- `server/internal/store/store_test.go` — kiểm thử CRUD nhân vật
- `server/internal/game/login.go` — `BoXuLy` giữ trạng thái phiên; rẽ nhánh có/chưa nhân vật; điều phối lệnh `-28`
- `server/internal/game/login_test.go` — cập nhật `khoGia` cho giao diện `Kho` mới
- `server/cmd/thuclient/main.go` — thêm bước tạo nhân vật để thử tay

---

## Ranh giới Kế hoạch 2 ↔ Kế hoạch 3

| Kế hoạch 2 làm | Kế hoạch 3 làm |
|---|---|
| Nhận gói `-28/2`, kiểm tên, tạo hàng trong `nhan_vat` | Gửi `-28/7` (danh sách làng), `-30/0` (ME_INFO), `-24` (MAP INFO) để client rời màn chờ |
| Lúc đăng nhập rẽ nhánh có/chưa nhân vật | Nội dung bản đồ thật (waypoint, quái, npc, đồ rơi) |
| `vaoTheGioi` = stub (ghi nhân vật vào phiên + log) | `vaoTheGioi` gửi đủ chuỗi gói vào map |
| Cột cốt lõi của nhân vật (§6.2, câu 10) | Cột mở rộng: túi đồ, kỹ năng, nhiệm vụ, pet |

> **Kiểm thử Kế hoạch 2 bằng client giả (`thuclient`), không bằng client Unity thật.** Client Unity thật gửi được gói tạo (hàng sẽ xuất hiện trong CSDL) nhưng **sau đó sẽ treo ở "Xin chờ"** vì chưa có chuỗi vào map — đó là điều **bình thường và đúng dự kiến** ở Kế hoạch 2. Kiểm thử với client Unity thật để tới Kế hoạch 3.

---

## Task 1: Mỗi kết nối một bộ xử lý riêng

**Vì sao task này đứng đầu:** `BoXuLy` hiện được tạo một lần (`main.go:99`) và dùng chung cho mọi kết nối, nên không thể giữ trạng thái riêng của phiên (đã đăng nhập bằng tài khoản nào). Tạo nhân vật bắt buộc phải biết điều đó. Cách sạch nhất: cho mỗi kết nối một `BoXuLy` riêng — khi đó trạng thái phiên chính là **trường của `BoXuLy`**, không cần bản đồ phiên→trạng thái, không cần dọn dẹp khi ngắt (kết nối đóng thì `BoXuLy` được thu gom theo).

**Files:**
- Modify: `server/internal/net/session.go`
- Modify: `server/internal/net/listener.go`
- Modify: `server/internal/net/session_test.go`
- Modify: `server/cmd/nroserver/main.go`

**Interfaces:**
- Produces: `net.LangNghe(diaChi string, taoXuLy func() XuLyGoiTin) error` và `net.PhucVu(conn net.Conn, taoXuLy func() XuLyGoiTin)` — mỗi kết nối gọi `taoXuLy()` một lần để lấy bộ xử lý mới.

- [ ] **Step 1: Xem test hiện có gọi `PhucVu`/`LangNghe` thế nào**

Run: `D:\Go\bin\go.exe test ./internal/net/ -run . -v` (trong thư mục `server/`)
Đọc `session_test.go` để biết chỗ nào dựng handler. Mục tiêu: sau khi đổi, các test này vẫn xanh.

- [ ] **Step 2: Đổi `PhucVu` để nhận hàm tạo**

Trong `server/internal/net/session.go`, đổi chữ ký và dòng dựng handler:

```go
// PhucVu chay vong doi mot ket noi cho toi khi dong.
//
// Nhan MOT HAM TAO bo xu ly, khong phai mot bo xu ly dung chung: moi ket noi
// can trang thai rieng (dang nhap bang tai khoan nao, dang giu nhan vat nao).
// Cho moi ket noi mot bo xu ly rieng nghia la trang thai do chi la truong cua
// bo xu ly — khong can ban do phien->trang thai, va khi ket noi dong thi bo xu
// ly duoc thu gom theo, khong ro ri.
func PhucVu(conn net.Conn, taoXuLy func() XuLyGoiTin) {
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

	xuLy := taoXuLy()

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
```

- [ ] **Step 3: Đổi `LangNghe` để truyền hàm tạo**

Trong `server/internal/net/listener.go`:

```go
// LangNghe mo cong TCP va phuc vu tung ket noi tren mot goroutine rieng.
// taoXuLy duoc goi mot lan cho MOI ket noi (xem PhucVu).
func LangNghe(diaChi string, taoXuLy func() XuLyGoiTin) error {
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
		go PhucVu(conn, taoXuLy)
	}
}
```

- [ ] **Step 4: Cập nhật MỌI nơi gọi `PhucVu` trong test**

Có **hai** file gọi `PhucVu`, cả hai phải sửa (nếu không sẽ vỡ biên dịch):

1. `server/internal/net/session_test.go` — mọi lời gọi `PhucVu(conn, xuLy)` đổi thành `PhucVu(conn, func() XuLyGoiTin { return xuLy })`.
2. `server/internal/game/login_test.go` — hàm `dungPhienThat(t, b *BoXuLy)` (khoảng dòng 194) gọi `mynet.PhucVu(benServer, b)`. Đổi thành:
   ```go
   go mynet.PhucVu(benServer, func() mynet.XuLyGoiTin { return b })
   ```
   Giữ đúng `b` sẵn có (một kết nối, một handler có trạng thái) — không dựng handler mới, vì test cần đọc lại trạng thái của chính `b` đó.

Với test một-kết-nối, bọc handler sẵn có trong hàm tạo trả về chính nó là đủ. Chỉ dựng handler mới trong hàm tạo nếu test thật sự mở nhiều kết nối cần trạng thái riêng.

- [ ] **Step 5: Cập nhật `main.go`**

Trong `server/cmd/nroserver/main.go`, đổi dòng 99:

```go
		loiMayChu <- mynet.LangNghe(diaChi, func() mynet.XuLyGoiTin {
			return game.NewBoXuLy(kho)
		})
```

- [ ] **Step 6: Chạy toàn bộ test — phải vẫn xanh**

Run (trong `server/`):
```
$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro_test'
D:\Go\bin\go.exe test ./...
```
Expected: tất cả PASS (70 test của Kế hoạch 1 không được vỡ). `D:\Go\bin\gofmt.exe -l .` rỗng. `D:\Go\bin\go.exe vet ./...` sạch.

- [ ] **Step 7: Commit**

```bash
git add server/internal/net/ server/cmd/nroserver/main.go
git commit -m "refactor(net): moi ket noi mot bo xu ly rieng de giu trang thai phien"
```

---

## Task 2: Bảng `nhan_vat`

**Files:**
- Create: `server/internal/store/migrations/0002_nhan_vat.sql`

**Interfaces:**
- Produces: bảng `nhan_vat` với các cột cốt lõi; chỉ mục duy nhất `nhan_vat_ten_duy_nhat` trên `lower(ten)`; ràng buộc **một nhân vật mỗi tài khoản** qua `UNIQUE (tai_khoan_id)`.

- [ ] **Step 1: Viết file migration**

Tạo `server/internal/store/migrations/0002_nhan_vat.sql`:

```sql
-- Bang nhan vat. Theo nguyen tac §6.2 cua thiet ke: CAC COT RIENG BIET, tra cuu
-- duoc — KHONG gop nhieu thu vao mot chuoi JSON nhu bang `player` cu (bang cu
-- co ~40 cot text JSON giai ma bang vi tri mang, rat gion, doi thu tu ghi la
-- hong het du lieu). Cot mo rong (tui do, ky nang, nhiem vu) se la bang rieng o
-- cac ke hoach sau, khong nhet vao day.
CREATE TABLE IF NOT EXISTS nhan_vat (
    id              BIGSERIAL PRIMARY KEY,
    tai_khoan_id    BIGINT      NOT NULL REFERENCES tai_khoan(id),
    ten             TEXT        NOT NULL,
    lang            SMALLINT    NOT NULL,   -- id lang/chung toc (KHONG phai gioi tinh)
    kieu_toc        SMALLINT    NOT NULL,   -- id sprite toc
    -- vi tri trong the gioi
    map_id          INTEGER     NOT NULL,
    vi_tri_x        INTEGER     NOT NULL,
    vi_tri_y        INTEGER     NOT NULL,
    -- chi so cot loi
    cap             INTEGER     NOT NULL DEFAULT 1,
    kinh_nghiem     BIGINT      NOT NULL DEFAULT 0,
    suc_manh        BIGINT      NOT NULL DEFAULT 0,   -- nhan luc
    hp_toi_da       INTEGER     NOT NULL,
    hp_hien_tai     INTEGER     NOT NULL,
    chakra_toi_da   INTEGER     NOT NULL,
    chakra_hien_tai INTEGER     NOT NULL,
    -- tien
    vang            BIGINT      NOT NULL DEFAULT 0,
    ngoc            BIGINT      NOT NULL DEFAULT 0,
    -- thoi gian
    tao_luc         TIMESTAMPTZ NOT NULL DEFAULT now(),
    dang_nhap_cuoi  TIMESTAMPTZ
);

-- Ten nhan vat khong phan biet hoa thuong: khong cho hai nguoi ten "Naruto" va
-- "naruto" cung ton tai (chong gia mao). Cung ly do dung lower() nhu ten dang
-- nhap o Ke hoach 1 — va vi charset chi la [A-Za-z0-9] (kiem o tang Go) nen
-- lower() duoi locale C hoat dong dung, khong dinh loi dau tieng Viet.
CREATE UNIQUE INDEX IF NOT EXISTS nhan_vat_ten_duy_nhat
    ON nhan_vat (lower(ten));

-- Mot tai khoan chi mot nhan vat: dung voi mo hinh server cu (client khong co
-- man chon nhan vat that su). Rang buoc o CSDL de khong the tao hai nhan vat
-- cho cung tai khoan du co dua tranh (race) o tang tren.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'nhan_vat_mot_tai_khoan'
    ) THEN
        ALTER TABLE nhan_vat
            ADD CONSTRAINT nhan_vat_mot_tai_khoan UNIQUE (tai_khoan_id);
    END IF;
END $$;

-- Tra nhan vat theo tai khoan la truy van nong nhat (moi lan dang nhap). Chi
-- muc nay khong that su can vi da co UNIQUE(tai_khoan_id) o tren tu dong tao
-- chi muc — giu lai lam tai lieu y dinh, IF NOT EXISTS nen khong trung.
CREATE INDEX IF NOT EXISTS nhan_vat_theo_tai_khoan
    ON nhan_vat (tai_khoan_id);
```

- [ ] **Step 2: Áp migration vào cả hai CSDL**

Run (PowerShell — đặt mật khẩu qua biến môi trường để psql không hỏi):
```
$env:PGPASSWORD='nrodev123'
D:\PostgreSQL\bin\psql.exe -h 127.0.0.1 -p 5433 -U postgres -d nro      -f server\internal\store\migrations\0002_nhan_vat.sql
D:\PostgreSQL\bin\psql.exe -h 127.0.0.1 -p 5433 -U postgres -d nro_test -f server\internal\store\migrations\0002_nhan_vat.sql
```
Expected: `CREATE TABLE`, `CREATE INDEX`, `DO`, `CREATE INDEX` — không lỗi. Chạy lại lần hai phải vẫn không lỗi (idempotent).

- [ ] **Step 3: Kiểm bảng đã tồn tại**

Run:
```
$env:PGPASSWORD='nrodev123'
D:\PostgreSQL\bin\psql.exe -h 127.0.0.1 -p 5433 -U postgres -d nro_test -c "\d nhan_vat"
```
Expected: liệt kê đủ các cột trên, thấy `nhan_vat_ten_duy_nhat`, `nhan_vat_mot_tai_khoan`, khóa ngoại tới `tai_khoan`.

- [ ] **Step 4: Commit**

```bash
git add server/internal/store/migrations/0002_nhan_vat.sql
git commit -m "feat(store): bang nhan_vat voi cot rieng biet, ten duy nhat, 1 nhan vat/tai khoan"
```

---

## Task 3: Danh mục 5 làng

**Files:**
- Create: `server/internal/game/lang.go`
- Create: `server/internal/game/lang_test.go`

**Interfaces:**
- Produces:
  - `type Lang struct { ID int16; Ten string; HpKhoiDau, ChakraKhoiDau, SatThuongKhoiDau int32; MapKhoiDau, ViTriXKhoiDau, ViTriYKhoiDau int32; VangKhoiDau int64 }`
  - `func TimLang(id int16) *Lang` — trả `nil` nếu không có làng đó
  - `func DanhSachLang() []Lang`

- [ ] **Step 1: Viết test cho danh mục làng**

Tạo `server/internal/game/lang_test.go`:

```go
package game

import "testing"

func TestTimLangHopLe(t *testing.T) {
	l := TimLang(0)
	if l == nil {
		t.Fatal("lang 0 (Lang La) phai ton tai")
	}
	if l.Ten != "Làng Lá" {
		t.Errorf("ten lang 0 = %q, mong doi \"Làng Lá\"", l.Ten)
	}
	// Moi lang phai co chi so khoi dau duong — nhan vat khong the co 0 mau.
	for _, l := range DanhSachLang() {
		if l.HpKhoiDau <= 0 || l.ChakraKhoiDau <= 0 {
			t.Errorf("lang %d (%s) co hp/chakra khoi dau khong duong", l.ID, l.Ten)
		}
	}
}

func TestTimLangKhongTonTai(t *testing.T) {
	if TimLang(99) != nil {
		t.Error("lang 99 khong ton tai, phai tra nil")
	}
	if TimLang(-1) != nil {
		t.Error("lang -1 khong ton tai, phai tra nil")
	}
}

func TestCoDuNamLang(t *testing.T) {
	// Thiet ke §7.2 co 5 lang. Client cu chi chon duoc 0/1/2, nhung may chu
	// dinh nghia san ca 5 de client dung lai sau khong phai sua may chu.
	if n := len(DanhSachLang()); n != 5 {
		t.Errorf("co %d lang, mong doi 5", n)
	}
}
```

- [ ] **Step 2: Chạy test — phải fail vì chưa có code**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run TestTimLang -v`
Expected: FAIL (không biên dịch được — `TimLang` chưa định nghĩa).

- [ ] **Step 3: Viết danh mục làng**

Tạo `server/internal/game/lang.go`:

```go
package game

// Lang la mot lang nhan gia (thay cho "chung toc" cua ban Dragon Ball cu).
// Chi so khoi dau la GIA TRI DAU — se can bang lai o Ke hoach 4 (chien dau).
// Map va vi tri khoi dau la CHO DAT TRUOC — ban do that dung o Ke hoach 3/6;
// o day chi can mot toa do hop le de nhan vat ton tai.
type Lang struct {
	ID               int16
	Ten              string
	HpKhoiDau        int32
	ChakraKhoiDau    int32
	SatThuongKhoiDau int32 // dung o Ke hoach 4
	MapKhoiDau       int32
	ViTriXKhoiDau    int32
	ViTriYKhoiDau    int32
	VangKhoiDau      int64
}

// danhSachLang: 5 lang theo thiet ke §7.2. Chi so phan anh tinh cach tung lang
// (Cat chakra cao danh xa, Suong Mu can chien sat thuong, Da mau trau phong thu,
// May toc do chi mang, La can bang). Client cu CHI gui duoc lang 0/1/2 vi cung
// dung mang toc 3x3 — lang 3/4 dung duoc khi client duoc dung lai.
//
// Ke hoach 6 se chuyen bang nay sang file YAML nap nong. Gio de trong code cho
// don gian, khong dung som he thong nap YAML.
var danhSachLang = []Lang{
	{ID: 0, Ten: "Làng Lá", HpKhoiDau: 150, ChakraKhoiDau: 150, SatThuongKhoiDau: 12, MapKhoiDau: 1, ViTriXKhoiDau: 100, ViTriYKhoiDau: 384, VangKhoiDau: 1000},
	{ID: 1, Ten: "Làng Cát", HpKhoiDau: 100, ChakraKhoiDau: 250, SatThuongKhoiDau: 10, MapKhoiDau: 2, ViTriXKhoiDau: 100, ViTriYKhoiDau: 384, VangKhoiDau: 1000},
	{ID: 2, Ten: "Làng Sương Mù", HpKhoiDau: 130, ChakraKhoiDau: 100, SatThuongKhoiDau: 16, MapKhoiDau: 3, ViTriXKhoiDau: 100, ViTriYKhoiDau: 384, VangKhoiDau: 1000},
	{ID: 3, Ten: "Làng Đá", HpKhoiDau: 250, ChakraKhoiDau: 100, SatThuongKhoiDau: 10, MapKhoiDau: 4, ViTriXKhoiDau: 100, ViTriYKhoiDau: 384, VangKhoiDau: 1000},
	{ID: 4, Ten: "Làng Mây", HpKhoiDau: 120, ChakraKhoiDau: 130, SatThuongKhoiDau: 13, MapKhoiDau: 5, ViTriXKhoiDau: 100, ViTriYKhoiDau: 384, VangKhoiDau: 1000},
}

// TimLang tra ve lang theo id, hoac nil neu khong co.
func TimLang(id int16) *Lang {
	for i := range danhSachLang {
		if danhSachLang[i].ID == id {
			return &danhSachLang[i]
		}
	}
	return nil
}

// DanhSachLang tra ban sao danh sach lang (chi doc).
func DanhSachLang() []Lang {
	ds := make([]Lang, len(danhSachLang))
	copy(ds, danhSachLang)
	return ds
}
```

- [ ] **Step 4: Chạy test — phải PASS**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run "TestTimLang|TestCoDuNamLang" -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add server/internal/game/lang.go server/internal/game/lang_test.go
git commit -m "feat(game): danh muc 5 lang voi chi so khoi dau"
```

---

## Task 4: Kho dữ liệu nhân vật (CRUD)

**Files:**
- Modify: `server/internal/store/store.go`
- Modify: `server/internal/store/store_test.go`

**Interfaces:**
- Produces (trong package `store`):
  - `var ErrKhongThayNhanVat = errors.New(...)`, `var ErrTenNhanVatDaDung = errors.New(...)`, `var ErrDaCoNhanVat = errors.New(...)`
  - `type NhanVat struct { ID, TaiKhoanID int64; Ten string; Lang, KieuToc int16; MapID, ViTriX, ViTriY, Cap int32; KinhNghiem, SucManh int64; HpToiDa, HpHienTai, ChakraToiDa, ChakraHienTai int32; Vang, Ngoc int64 }`
  - `type NhanVatMoi struct { TaiKhoanID int64; Ten string; Lang, KieuToc int16; MapID, ViTriX, ViTriY int32; HpToiDa, ChakraToiDa int32; Vang int64 }`
  - `func (s *Store) TaoNhanVat(ctx context.Context, nv NhanVatMoi) (*NhanVat, error)` — cấp cap=1, kinh_nghiem=0, suc_manh=0, hp_hien_tai=HpToiDa, chakra_hien_tai=ChakraToiDa, ngoc=0. Trả `ErrTenNhanVatDaDung` hoặc `ErrDaCoNhanVat` khi vi phạm ràng buộc.
  - `func (s *Store) LayNhanVatTheoTaiKhoan(ctx context.Context, taiKhoanID int64) (*NhanVat, error)` — `ErrKhongThayNhanVat` nếu chưa có.
  - `func (s *Store) TenNhanVatDaDung(ctx context.Context, ten string) (bool, error)`

- [ ] **Step 1: Viết test CRUD (dùng CSDL thật `nro_test`)**

Trong `server/internal/store/store_test.go`, thêm. **Đọc phần đầu file** để biết cách các test hiện có mở kết nối `nro_test` (biến `NRO_TEST_DSN`), cách bỏ qua khi không có DSN, và cách tạo tài khoản thử. Dùng lại đúng lối đó. Mẫu:

```go
func TestTaoVaLayNhanVat(t *testing.T) {
	s := moKhoThu(t) // helper co san trong store_test.go cua Ke hoach 1
	ctx := context.Background()

	// Nhan vat can mot tai khoan co that (khoa ngoai). Tao tai khoan thu.
	taiKhoanID := taoTaiKhoanThu(t, s, "chunin"+hauToNganNgau(t))

	nv, err := s.TaoNhanVat(ctx, NhanVatMoi{
		TaiKhoanID: taiKhoanID, Ten: "Naruto" + hauToNganNgau(t),
		Lang: 0, KieuToc: 5, MapID: 1, ViTriX: 100, ViTriY: 384,
		HpToiDa: 150, ChakraToiDa: 150, Vang: 1000,
	})
	if err != nil {
		t.Fatalf("TaoNhanVat loi: %v", err)
	}
	if nv.ID == 0 {
		t.Error("nhan vat moi phai co ID")
	}
	// Gia tri suy ra phai dung: mau/chakra hien tai bang toi da, cap = 1.
	if nv.HpHienTai != 150 || nv.ChakraHienTai != 150 {
		t.Errorf("hp/chakra hien tai = %d/%d, mong doi 150/150", nv.HpHienTai, nv.ChakraHienTai)
	}
	if nv.Cap != 1 {
		t.Errorf("cap = %d, mong doi 1", nv.Cap)
	}

	// Lay lai theo tai khoan phai ra dung nhan vat do, dung tung truong.
	lay, err := s.LayNhanVatTheoTaiKhoan(ctx, taiKhoanID)
	if err != nil {
		t.Fatalf("LayNhanVatTheoTaiKhoan loi: %v", err)
	}
	if lay.ID != nv.ID || lay.Ten != nv.Ten || lay.Lang != 0 || lay.MapID != 1 {
		t.Errorf("lay lai khong khop: %+v", lay)
	}
}

func TestLayNhanVatChuaCo(t *testing.T) {
	s := moKhoThu(t)
	taiKhoanID := taoTaiKhoanThu(t, s, "genin"+hauToNganNgau(t))
	_, err := s.LayNhanVatTheoTaiKhoan(context.Background(), taiKhoanID)
	if !errors.Is(err, ErrKhongThayNhanVat) {
		t.Errorf("mong doi ErrKhongThayNhanVat, nhan %v", err)
	}
}

func TestTenNhanVatTrungBiTuChoi(t *testing.T) {
	s := moKhoThu(t)
	ctx := context.Background()
	ten := "Sasuke" + hauToNganNgau(t)

	tk1 := taoTaiKhoanThu(t, s, "acc1"+hauToNganNgau(t))
	tk2 := taoTaiKhoanThu(t, s, "acc2"+hauToNganNgau(t))

	if _, err := s.TaoNhanVat(ctx, NhanVatMoi{TaiKhoanID: tk1, Ten: ten, Lang: 0, KieuToc: 5, MapID: 1, ViTriX: 100, ViTriY: 384, HpToiDa: 150, ChakraToiDa: 150, Vang: 1000}); err != nil {
		t.Fatalf("tao nhan vat 1 loi: %v", err)
	}
	// Ten y het (khac hoa thuong) tu tai khoan khac -> phai bi tu choi.
	_, err := s.TaoNhanVat(ctx, NhanVatMoi{TaiKhoanID: tk2, Ten: strings.ToUpper(ten), Lang: 0, KieuToc: 5, MapID: 1, ViTriX: 100, ViTriY: 384, HpToiDa: 150, ChakraToiDa: 150, Vang: 1000})
	if !errors.Is(err, ErrTenNhanVatDaDung) {
		t.Errorf("mong doi ErrTenNhanVatDaDung, nhan %v", err)
	}
}

func TestMotTaiKhoanChiMotNhanVat(t *testing.T) {
	s := moKhoThu(t)
	ctx := context.Background()
	tk := taoTaiKhoanThu(t, s, "solo"+hauToNganNgau(t))

	if _, err := s.TaoNhanVat(ctx, NhanVatMoi{TaiKhoanID: tk, Ten: "First" + hauToNganNgau(t), Lang: 0, KieuToc: 5, MapID: 1, ViTriX: 100, ViTriY: 384, HpToiDa: 150, ChakraToiDa: 150, Vang: 1000}); err != nil {
		t.Fatalf("tao nhan vat dau loi: %v", err)
	}
	_, err := s.TaoNhanVat(ctx, NhanVatMoi{TaiKhoanID: tk, Ten: "Second" + hauToNganNgau(t), Lang: 1, KieuToc: 5, MapID: 2, ViTriX: 100, ViTriY: 384, HpToiDa: 100, ChakraToiDa: 250, Vang: 1000})
	if !errors.Is(err, ErrDaCoNhanVat) {
		t.Errorf("mong doi ErrDaCoNhanVat, nhan %v", err)
	}
}
```

> **Ghi chú cho người cài đặt:** file `store_test.go` của Kế hoạch 1 đã có helper mở kho thử. Nếu **chưa** có `taoTaiKhoanThu` và `hauToNganNgau`, thêm chúng vào file test (đừng để trùng tên với helper sẵn có):
> ```go
> // hauToNganNgau tra mot chuoi ngan duy nhat de ten khong dung nhau giua cac
> // lan chay test (khong dung Date/rand — dung dia chi bo nho cua bien cuc bo,
> // luon khac nhau va hop le [a-z0-9]).
> func hauToNganNgau(t *testing.T) string {
> 	t.Helper()
> 	var x int
> 	return fmt.Sprintf("%p", &x)[2:8] // 6 ky tu hex tu con tro
> }
>
> func taoTaiKhoanThu(t *testing.T, s *Store, ten string) int64 {
> 	t.Helper()
> 	var id int64
> 	err := s.pool.QueryRow(context.Background(),
> 		`INSERT INTO tai_khoan (ten_dang_nhap, mat_khau_bam) VALUES ($1, 'x') RETURNING id`,
> 		ten,
> 	).Scan(&id)
> 	if err != nil {
> 		t.Fatalf("tao tai khoan thu loi: %v", err)
> 	}
> 	t.Cleanup(func() {
> 		_, _ = s.pool.Exec(context.Background(), `DELETE FROM nhan_vat WHERE tai_khoan_id = $1`, id)
> 		_, _ = s.pool.Exec(context.Background(), `DELETE FROM tai_khoan WHERE id = $1`, id)
> 	})
> 	return id
> }
> ```
> Tên tài khoản thử phải khớp CHECK `[A-Za-z0-9_]{4,20}` — các mẫu trên (`chunin<hex>`...) đều hợp lệ. Nhớ `import ("fmt"; "strings")` nếu chưa có.

- [ ] **Step 2: Chạy test — phải fail vì chưa có code**

Run (trong `server/`):
```
$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro_test'
D:\Go\bin\go.exe test ./internal/store/ -run "TestTaoVaLayNhanVat|TestLayNhanVatChuaCo|TestTenNhanVatTrung|TestMotTaiKhoan" -v
```
Expected: FAIL (không biên dịch — `TaoNhanVat` chưa có).

- [ ] **Step 3: Thêm kiểu và lỗi sentinel vào `store.go`**

Trong `server/internal/store/store.go`, thêm import `pgconn` và các khai báo:

```go
import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrKhongThayNhanVat  = errors.New("khong tim thay nhan vat")
	ErrTenNhanVatDaDung  = errors.New("ten nhan vat da ton tai")
	ErrDaCoNhanVat       = errors.New("tai khoan da co nhan vat")
)

// NhanVat la nhan vat da luu, du cot cot loi de vao the gioi.
type NhanVat struct {
	ID            int64
	TaiKhoanID    int64
	Ten           string
	Lang          int16
	KieuToc       int16
	MapID         int32
	ViTriX        int32
	ViTriY        int32
	Cap           int32
	KinhNghiem    int64
	SucManh       int64
	HpToiDa       int32
	HpHienTai     int32
	ChakraToiDa   int32
	ChakraHienTai int32
	Vang          int64
	Ngoc          int64
}

// NhanVatMoi la du lieu can de tao mot nhan vat. Cac chi so suy ra (cap ban dau,
// mau/chakra hien tai = toi da) do TaoNhanVat tu dat, khong nhan tu ben ngoai.
type NhanVatMoi struct {
	TaiKhoanID  int64
	Ten         string
	Lang        int16
	KieuToc     int16
	MapID       int32
	ViTriX      int32
	ViTriY      int32
	HpToiDa     int32
	ChakraToiDa int32
	Vang        int64
}
```

- [ ] **Step 4: Viết `TaoNhanVat`**

Thêm vào `store.go`:

```go
// TaoNhanVat them mot nhan vat moi. Mau/chakra hien tai khoi tao bang toi da,
// cap = 1. Neu ten trung (khong phan biet hoa thuong) tra ErrTenNhanVatDaDung;
// neu tai khoan da co nhan vat tra ErrDaCoNhanVat — de tang tren doi hai loi
// nay ra hai cau bao khac nhau cho nguoi choi.
func (s *Store) TaoNhanVat(ctx context.Context, nv NhanVatMoi) (*NhanVat, error) {
	kq := &NhanVat{}
	err := s.pool.QueryRow(ctx,
		`INSERT INTO nhan_vat
		    (tai_khoan_id, ten, lang, kieu_toc, map_id, vi_tri_x, vi_tri_y,
		     hp_toi_da, hp_hien_tai, chakra_toi_da, chakra_hien_tai, vang)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$8,$9,$9,$10)
		 RETURNING id, tai_khoan_id, ten, lang, kieu_toc, map_id, vi_tri_x, vi_tri_y,
		           cap, kinh_nghiem, suc_manh, hp_toi_da, hp_hien_tai,
		           chakra_toi_da, chakra_hien_tai, vang, ngoc`,
		nv.TaiKhoanID, nv.Ten, nv.Lang, nv.KieuToc, nv.MapID, nv.ViTriX, nv.ViTriY,
		nv.HpToiDa, nv.ChakraToiDa, nv.Vang,
	).Scan(
		&kq.ID, &kq.TaiKhoanID, &kq.Ten, &kq.Lang, &kq.KieuToc, &kq.MapID,
		&kq.ViTriX, &kq.ViTriY, &kq.Cap, &kq.KinhNghiem, &kq.SucManh,
		&kq.HpToiDa, &kq.HpHienTai, &kq.ChakraToiDa, &kq.ChakraHienTai,
		&kq.Vang, &kq.Ngoc,
	)
	if err != nil {
		// Phan biet hai rang buoc duy nhat khac nhau bang TEN rang buoc, khong
		// chi bang ma loi 23505 (ca hai deu 23505). Nho vay tang tren bao dung
		// "ten da co" hay "da co nhan vat".
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "nhan_vat_ten_duy_nhat":
				return nil, ErrTenNhanVatDaDung
			case "nhan_vat_mot_tai_khoan":
				return nil, ErrDaCoNhanVat
			}
		}
		return nil, fmt.Errorf("them nhan vat that bai: %w", err)
	}
	return kq, nil
}
```

> ⚠️ `$8,$8` và `$9,$9` là cố ý: `hp_hien_tai` khởi tạo bằng `hp_toi_da` (`$8`), `chakra_hien_tai` bằng `chakra_toi_da` (`$9`). Đừng "sửa" thành tham số riêng.

- [ ] **Step 5: Viết `LayNhanVatTheoTaiKhoan` và `TenNhanVatDaDung`**

```go
// LayNhanVatTheoTaiKhoan tra nhan vat cua mot tai khoan, hoac ErrKhongThayNhanVat.
func (s *Store) LayNhanVatTheoTaiKhoan(ctx context.Context, taiKhoanID int64) (*NhanVat, error) {
	kq := &NhanVat{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, tai_khoan_id, ten, lang, kieu_toc, map_id, vi_tri_x, vi_tri_y,
		        cap, kinh_nghiem, suc_manh, hp_toi_da, hp_hien_tai,
		        chakra_toi_da, chakra_hien_tai, vang, ngoc
		   FROM nhan_vat
		  WHERE tai_khoan_id = $1
		  LIMIT 1`, taiKhoanID,
	).Scan(
		&kq.ID, &kq.TaiKhoanID, &kq.Ten, &kq.Lang, &kq.KieuToc, &kq.MapID,
		&kq.ViTriX, &kq.ViTriY, &kq.Cap, &kq.KinhNghiem, &kq.SucManh,
		&kq.HpToiDa, &kq.HpHienTai, &kq.ChakraToiDa, &kq.ChakraHienTai,
		&kq.Vang, &kq.Ngoc,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrKhongThayNhanVat
	}
	if err != nil {
		return nil, fmt.Errorf("tra cuu nhan vat that bai: %w", err)
	}
	return kq, nil
}

// TenNhanVatDaDung kiem nhanh mot ten da co ai dung chua (khong phan biet hoa
// thuong). Dung de bao nguoi choi som; rang buoc duy nhat trong CSDL moi la
// chot chan that su khi co dua tranh.
func (s *Store) TenNhanVatDaDung(ctx context.Context, ten string) (bool, error) {
	var co bool
	err := s.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM nhan_vat WHERE lower(ten) = lower($1))`, ten,
	).Scan(&co)
	if err != nil {
		return false, fmt.Errorf("kiem ten nhan vat that bai: %w", err)
	}
	return co, nil
}
```

- [ ] **Step 6: Chạy test — phải PASS**

Run (trong `server/`):
```
$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro_test'
D:\Go\bin\go.exe test ./internal/store/ -v
```
Expected: tất cả PASS (gồm test Kế hoạch 1 và 4 test mới).

- [ ] **Step 7: Commit**

```bash
git add server/internal/store/store.go server/internal/store/store_test.go
git commit -m "feat(store): CRUD nhan vat — tao, lay theo tai khoan, kiem ten trung"
```

---

## Task 5: Đăng nhập rẽ nhánh có/chưa nhân vật

**Files:**
- Modify: `server/internal/game/login.go`
- Modify: `server/internal/game/login_test.go`

**Interfaces:**
- Produces:
  - Giao diện `Kho` (mở rộng từ `KhoTaiKhoan`) thêm ba hàm nhân vật; `BoXuLy.kho` đổi sang kiểu `Kho`.
  - `BoXuLy` giữ trạng thái phiên: trường `taiKhoan *store.TaiKhoan`, `nhanVat *store.NhanVat`.
  - `func (b *BoXuLy) vaoTheGioi(s *mynet.Session, nv *store.NhanVat) error` — **stub Kế hoạch 2**: gán `b.nhanVat`, ghi log. Kế hoạch 3 sẽ gửi chuỗi gói vào map.
- Consumes: `store.LayNhanVatTheoTaiKhoan`, `store.ErrKhongThayNhanVat` (Task 4).

- [ ] **Step 1: Cập nhật test đăng nhập cho luồng mới**

Trong `server/internal/game/login_test.go`:

1. `khoGia` hiện chỉ cài `LayTaiKhoan`. Thêm ba hàm để nó thỏa giao diện `Kho` mới. Cho phép test nạp sẵn kết quả:

```go
type khoGia struct {
	taiKhoan       *store.TaiKhoan
	loiLayTaiKhoan error

	nhanVat          *store.NhanVat // nil = tai khoan chua co nhan vat
	loiLayNhanVat    error          // dat = store.ErrKhongThayNhanVat de gia lap "chua co"
	tenDaDung        bool
	nhanVatVuaTao    *store.NhanVat
	loiTaoNhanVat    error
}

func (k *khoGia) LayTaiKhoan(ctx context.Context, ten string) (*store.TaiKhoan, error) {
	if k.loiLayTaiKhoan != nil {
		return nil, k.loiLayTaiKhoan
	}
	return k.taiKhoan, nil
}
func (k *khoGia) LayNhanVatTheoTaiKhoan(ctx context.Context, taiKhoanID int64) (*store.NhanVat, error) {
	if k.loiLayNhanVat != nil {
		return nil, k.loiLayNhanVat
	}
	return k.nhanVat, nil
}
func (k *khoGia) TenNhanVatDaDung(ctx context.Context, ten string) (bool, error) {
	return k.tenDaDung, nil
}
func (k *khoGia) TaoNhanVat(ctx context.Context, nv store.NhanVatMoi) (*store.NhanVat, error) {
	if k.loiTaoNhanVat != nil {
		return nil, k.loiTaoNhanVat
	}
	return k.nhanVatVuaTao, nil
}
```

> Nếu `khoGia` của Kế hoạch 1 dùng con trỏ hàm/field khác, giữ tương thích: điều bắt buộc là nó cài đủ bốn phương thức của `Kho`. Các test cũ vốn chỉ đặt `taiKhoan` vẫn chạy vì các field mới mặc định nil/false.

2. Thêm test cho hai nhánh mới. Các test cũ của Kế hoạch 1 giả định đăng nhập xong luôn gửi lệnh `2`; **giờ điều đó chỉ đúng khi CHƯA có nhân vật**. Sửa các test cũ dựng `khoGia` với `loiLayNhanVat: store.ErrKhongThayNhanVat` (mô phỏng tài khoản chưa có nhân vật) để chúng vẫn kỳ vọng lệnh `2`. Rồi thêm:

```go
func TestDangNhapChuaCoNhanVat_GuiManTao(t *testing.T) {
	kho := &khoGia{
		taiKhoan:      taiKhoanHopLeChoTest(), // helper co san: tra *store.TaiKhoan mat khau khop "matkhau"
		loiLayNhanVat: store.ErrKhongThayNhanVat,
	}
	goiGuiRa := chayDangNhap(t, kho, "nguoidung", "matkhau") // helper dựng phiên qua net.Pipe, tra danh sach goi may chu gui ra
	if !coLenh(goiGuiRa, 2) {
		t.Error("chua co nhan vat -> phai gui lenh 2 (man tao nhan vat)")
	}
}

func TestDangNhapDaCoNhanVat_KhongGuiManTao(t *testing.T) {
	kho := &khoGia{
		taiKhoan: taiKhoanHopLeChoTest(),
		nhanVat:  &store.NhanVat{ID: 7, Ten: "Naruto", Lang: 0},
	}
	goiGuiRa := chayDangNhap(t, kho, "nguoidung", "matkhau")
	if coLenh(goiGuiRa, 2) {
		t.Error("da co nhan vat -> KHONG duoc gui lenh 2 (khong quay lai man tao)")
	}
}
```

> **Ghi chú:** `login_test.go` của Kế hoạch 1 (Task 8) đã có bộ khung end-to-end qua `net.Pipe` bắt các gói máy chủ gửi ra (đó là 3 test end-to-end đã thêm). Dùng lại đúng khung đó — đặt tên helper cho khớp cái đã có (`chayDangNhap`, `coLenh`, `taiKhoanHopLeChoTest` chỉ là tên gợi ý; nếu tên thật khác thì theo tên thật). Nếu helper cũ chưa trả **danh sách** gói mà chỉ gói đầu, mở rộng để đọc hết các gói tới khi kết nối rảnh.

- [ ] **Step 2: Chạy test — phải fail**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run "TestDangNhap" -v`
Expected: FAIL (chưa có logic rẽ nhánh; và có thể lỗi biên dịch do `khoGia` chưa đủ phương thức — đúng như mong đợi).

- [ ] **Step 3: Mở rộng giao diện `Kho` và thêm trạng thái phiên**

Trong `server/internal/game/login.go`, đổi khối giao diện + `BoXuLy`:

```go
// Kho gom moi thao tac CSDL ma tang game can. Tach thanh giao dien de kiem thu
// khong can CSDL that. (Ke hoach 1 chi co LayTaiKhoan; Ke hoach 2 them nhan vat.)
type Kho interface {
	LayTaiKhoan(ctx context.Context, ten string) (*store.TaiKhoan, error)
	LayNhanVatTheoTaiKhoan(ctx context.Context, taiKhoanID int64) (*store.NhanVat, error)
	TenNhanVatDaDung(ctx context.Context, ten string) (bool, error)
	TaoNhanVat(ctx context.Context, nv store.NhanVatMoi) (*store.NhanVat, error)
}

// BoXuLy giu trang thai cua MOT phien (mot ket noi). Xem Task 1: moi ket noi
// mot BoXuLy rieng, nen cac truong duoi day khong bao gio bi hai phien tranh.
type BoXuLy struct {
	kho      Kho
	taiKhoan *store.TaiKhoan // nil khi chua dang nhap
	nhanVat  *store.NhanVat  // nil khi chua co / chua vao the gioi
}

func NewBoXuLy(kho Kho) *BoXuLy { return &BoXuLy{kho: kho} }
```

> Xóa khai báo `type KhoTaiKhoan interface {...}` cũ (thay bằng `Kho`). Nếu còn chỗ nào tham chiếu `KhoTaiKhoan`, đổi hết sang `Kho`.

- [ ] **Step 4: Sửa phần cuối `KiemTraDangNhap`/`XuLy` để rẽ nhánh**

Trong `XuLy` của `login.go`, thay đoạn cuối (từ chỗ `slog.Info("dang nhap thanh cong"...)` trở đi):

```go
	slog.Info("dang nhap thanh cong",
		"ip", s.DiaChiIP, "ten", kq.TaiKhoan.TenDangNhap, "id", kq.TaiKhoan.ID)
	b.taiKhoan = kq.TaiKhoan // ghi nhan phien nay da dang nhap bang tai khoan nao

	nv, err := b.kho.LayNhanVatTheoTaiKhoan(context.Background(), kq.TaiKhoan.ID)
	if errors.Is(err, store.ErrKhongThayNhanVat) {
		return guiChuyenSangTaoNhanVat(s) // chua co -> mo man tao
	}
	if err != nil {
		slog.Error("loi he thong khi tai nhan vat", "ip", s.DiaChiIP, "loi", err)
		return guiTuChoi(s, "Không tải được nhân vật, vui lòng thử lại sau.")
	}
	return b.vaoTheGioi(s, nv) // da co -> vao the gioi (stub Ke hoach 2)
}
```

- [ ] **Step 5: Thêm `vaoTheGioi` (stub)**

Sửa ghi chú của `guiChuyenSangTaoNhanVat` (bỏ dòng "Ke hoach 2 se thay bang...") và thêm hàm mới:

```go
// vaoTheGioi la diem hen chung cho ca hai duong: dang nhap co san nhan vat, va
// vua tao nhan vat xong. Ke hoach 2 chi ghi nhan nhan vat vao phien va log —
// Ke hoach 3 se gui du chuoi goi (-28/7 danh sach lang, -30/0 ME_INFO, -24 MAP
// INFO) de client roi man "Xin cho". Voi Ke hoach 2, client Unity that se treo
// o day: dung du kien, kiem thu bang client gia.
func (b *BoXuLy) vaoTheGioi(s *mynet.Session, nv *store.NhanVat) error {
	b.nhanVat = nv
	slog.Info("nhan vat vao the gioi (stub Ke hoach 2)",
		"ten", nv.Ten, "lang", nv.Lang, "id", nv.ID, "map", nv.MapID)
	return nil
}
```

- [ ] **Step 6: Chạy test — phải PASS**

Run (trong `server/`):
```
$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro_test'
D:\Go\bin\go.exe test ./internal/game/ -v
```
Expected: tất cả PASS. `gofmt -l .` rỗng, `go vet ./...` sạch.

- [ ] **Step 7: Commit**

```bash
git add server/internal/game/login.go server/internal/game/login_test.go
git commit -m "feat(game): dang nhap re nhanh co/chua nhan vat; BoXuLy giu trang thai phien"
```

---

## Task 6: Xử lý gói tạo nhân vật

**Files:**
- Create: `server/internal/game/nhanvat.go`
- Create: `server/internal/game/nhanvat_test.go`
- Modify: `server/internal/game/login.go` (điều phối lệnh `-28`)

**Interfaces:**
- Consumes: `b.taiKhoan`, `b.kho`, `TimLang` (Task 3), `store.NhanVatMoi`, `store.TaoNhanVat`, `store.TenNhanVatDaDung`, `store.ErrTenNhanVatDaDung`, `store.ErrDaCoNhanVat`, `b.vaoTheGioi` (Task 5).
- Produces: `func (b *BoXuLy) xuLyTaoNhanVat(s *mynet.Session, r *protocol.Reader) error`; hàm `guiHopThoai(s, cau string) error` (chỉ gửi `-26`, KHÔNG kèm `-102`).

- [ ] **Step 1: Viết test cho xử lý tạo nhân vật (dùng `khoGia`)**

Tạo `server/internal/game/nhanvat_test.go`. Dùng lại khung `net.Pipe` như `login_test.go`. Các ca:

```go
package game

import (
	"testing"

	"nroserver/internal/store"
)

// dungGoiTaoNhanVat dong goi than tin gói -28/2 (khong gom byte lenh con -28,
// vi handler nhan Reader da tro sau byte lenh con). Tra []byte than sau lenh con.
func thanTaoNhanVat(t *testing.T, ten string, lang, toc byte) []byte {
	t.Helper()
	w := protocol.NewWriter()
	if err := w.UTF(ten); err != nil {
		t.Fatal(err)
	}
	w.Byte(lang)
	w.Byte(toc)
	return w.Bytes() // dung ham xuat byte cua Writer (xem datastream.go)
}

func TestTaoNhanVat_TenQuaNgan_TuChoi(t *testing.T) {
	b := boXuLyDaDangNhap(t) // helper: BoXuLy voi b.taiKhoan da set, kho gia
	s, docGoi := phienGia(t)  // helper: net.Pipe, tra ham doc goi may chu gui
	r := protocol.NewReader(thanTaoNhanVat(t, "abc", 0, 5)) // 3 ky tu < 5
	if err := b.xuLyTaoNhanVat(s, r); err != nil {
		t.Fatalf("xuLyTaoNhanVat loi: %v", err)
	}
	goi := docGoi()
	if len(goi) == 0 || goi[0].Command != -26 {
		t.Errorf("ten qua ngan phai tra -26, nhan %+v", goi)
	}
}

func TestTaoNhanVat_KyTuLa_TuChoi(t *testing.T) {
	b := boXuLyDaDangNhap(t)
	s, docGoi := phienGia(t)
	r := protocol.NewReader(thanTaoNhanVat(t, "na ru to", 0, 5)) // co dau cach
	_ = b.xuLyTaoNhanVat(s, r)
	if goi := docGoi(); len(goi) == 0 || goi[0].Command != -26 {
		t.Errorf("ky tu la phai tra -26, nhan %+v", goi)
	}
}

func TestTaoNhanVat_LangKhongHopLe_TuChoi(t *testing.T) {
	b := boXuLyDaDangNhap(t)
	s, docGoi := phienGia(t)
	r := protocol.NewReader(thanTaoNhanVat(t, "Naruto", 99, 5)) // lang 99 khong ton tai
	_ = b.xuLyTaoNhanVat(s, r)
	if goi := docGoi(); len(goi) == 0 || goi[0].Command != -26 {
		t.Errorf("lang khong hop le phai tra -26, nhan %+v", goi)
	}
}

func TestTaoNhanVat_TenDaDung_TuChoi(t *testing.T) {
	b := boXuLyDaDangNhap(t)
	b.kho.(*khoGia).tenDaDung = true
	s, docGoi := phienGia(t)
	r := protocol.NewReader(thanTaoNhanVat(t, "Naruto", 0, 5))
	_ = b.xuLyTaoNhanVat(s, r)
	if goi := docGoi(); len(goi) == 0 || goi[0].Command != -26 {
		t.Errorf("ten da dung phai tra -26, nhan %+v", goi)
	}
}

func TestTaoNhanVat_ChuaDangNhap_BoQua(t *testing.T) {
	b := NewBoXuLy(&khoGia{}) // b.taiKhoan == nil: chua dang nhap
	s, docGoi := phienGia(t)
	r := protocol.NewReader(thanTaoNhanVat(t, "Naruto", 0, 5))
	if err := b.xuLyTaoNhanVat(s, r); err != nil {
		t.Fatalf("khong duoc tra loi khi chua dang nhap, chi bo qua: %v", err)
	}
	if goi := docGoi(); len(goi) != 0 {
		t.Errorf("chua dang nhap ma gui create -> phai IM LANG, khong gui gi; nhan %+v", goi)
	}
}

func TestTaoNhanVat_ThanhCong_VaoTheGioi(t *testing.T) {
	b := boXuLyDaDangNhap(t)
	b.kho.(*khoGia).nhanVatVuaTao = &store.NhanVat{ID: 42, Ten: "Naruto", Lang: 0, MapID: 1}
	s, _ := phienGia(t)
	r := protocol.NewReader(thanTaoNhanVat(t, "Naruto", 0, 5))
	if err := b.xuLyTaoNhanVat(s, r); err != nil {
		t.Fatalf("xuLyTaoNhanVat loi: %v", err)
	}
	// vaoTheGioi (stub) phai gan nhan vat vao phien.
	if b.nhanVat == nil || b.nhanVat.ID != 42 {
		t.Errorf("tao thanh cong phai gan nhan vat vao phien, b.nhanVat=%+v", b.nhanVat)
	}
}
```

> **Ghi chú helper:** `boXuLyDaDangNhap`, `phienGia`, `docGoi` dựng trên khung `net.Pipe` đã có ở `login_test.go`. `boXuLyDaDangNhap` = `NewBoXuLy(&khoGia{...})` rồi gán `b.taiKhoan = &store.TaiKhoan{ID: 1, TenDangNhap: "nguoidung"}`. `phienGia` mở `net.Pipe`, dựng `*mynet.Session` phía máy chủ, và trả hàm đọc **tất cả** gói máy chủ ghi ra (giải khung bằng `protocol.Framer` phía client giả). Nếu Kế hoạch 1 chưa tách các helper này ra dạng tái dùng được, trích chúng ra trong task này để cả `login_test.go` và `nhanvat_test.go` cùng dùng.

- [ ] **Step 2: Chạy test — phải fail**

Run: `D:\Go\bin\go.exe test ./internal/game/ -run TestTaoNhanVat -v`
Expected: FAIL (chưa có `xuLyTaoNhanVat`).

- [ ] **Step 3: Viết `nhanvat.go`**

Tạo `server/internal/game/nhanvat.go`:

```go
package game

import (
	"context"
	"errors"
	"log/slog"
	"regexp"

	mynet "nroserver/internal/net"
	"nroserver/internal/protocol"
	"nroserver/internal/store"
)

// Lenh con cua khung -28 ma Ke hoach 2 xu ly. Cac lenh con khac (versionGame,
// clientOk...) la Ke hoach 3.
const lenhConTaoNhanVat byte = 2

// reTenNhanVat: 5-10 ky tu, chi chu cai tieng Anh va chu so. Khop dung rang
// buoc cua may chu Java cu (Controller.java:1043 + Util.haveSpecialCharacter).
// Bien dich mot lan o muc goi.
var reTenNhanVat = regexp.MustCompile(`^[A-Za-z0-9]{5,10}$`)

// danhSachTenCam: ten khong duoc dat de tranh gia mao vai tro quan tri. So sanh
// khong phan biet hoa thuong. Danh sach ngan, mo rong sau khi can.
var danhSachTenCam = map[string]bool{
	"admin": true, "administrator": true, "gm": true, "system": true, "quantri": true,
}

// xuLyTaoNhanVat xu ly than tin goi -28/2 (Reader da tro qua byte lenh con).
//
// KHONG tin client: phai kiem ten o day truoc khi cham CSDL, va phai kiem phien
// da dang nhap chua — client co the gui create truoc khi dang nhap.
func (b *BoXuLy) xuLyTaoNhanVat(s *mynet.Session, r *protocol.Reader) error {
	if b.taiKhoan == nil {
		// Chua dang nhap ma doi tao nhan vat: bo qua im lang, KHONG log (cung
		// ly do ngap nhat ky nhu cac lenh khong ho tro o login.go). Client that
		// khong bao gio lam vay; chi ke gia lap moi gui.
		return nil
	}
	if b.nhanVat != nil {
		return guiHopThoai(s, "Bạn đã có nhân vật.")
	}

	ten, err := r.UTF()
	if err != nil {
		// Than goi hong: bo ket noi (tra loi de tang net dong phien).
		return err
	}
	langByte, err := r.Byte()
	if err != nil {
		return err
	}
	kieuTocByte, err := r.Byte()
	if err != nil {
		return err
	}

	if !reTenNhanVat.MatchString(ten) {
		return guiHopThoai(s, "Tên nhân vật dài 5–10 ký tự, chỉ gồm chữ cái và chữ số.")
	}
	if danhSachTenCam[toLowerASCII(ten)] {
		return guiHopThoai(s, "Tên này không được phép dùng.")
	}

	lang := TimLang(int16(langByte))
	if lang == nil {
		return guiHopThoai(s, "Làng không hợp lệ.")
	}

	ctx := context.Background()
	daDung, err := b.kho.TenNhanVatDaDung(ctx, ten)
	if err != nil {
		slog.Error("loi kiem ten nhan vat", "ip", s.DiaChiIP, "loi", err)
		return guiHopThoai(s, "Lỗi hệ thống, vui lòng thử lại.")
	}
	if daDung {
		return guiHopThoai(s, "Tên nhân vật đã tồn tại.")
	}

	nv, err := b.kho.TaoNhanVat(ctx, store.NhanVatMoi{
		TaiKhoanID:  b.taiKhoan.ID,
		Ten:         ten,
		Lang:        lang.ID,
		KieuToc:     int16(kieuTocByte),
		MapID:       lang.MapKhoiDau,
		ViTriX:      lang.ViTriXKhoiDau,
		ViTriY:      lang.ViTriYKhoiDau,
		HpToiDa:     lang.HpKhoiDau,
		ChakraToiDa: lang.ChakraKhoiDau,
		Vang:        lang.VangKhoiDau,
	})
	if errors.Is(err, store.ErrTenNhanVatDaDung) {
		// Dua tranh: co nguoi khac vua lay ten nay giua luc kiem va luc chen.
		return guiHopThoai(s, "Tên nhân vật đã tồn tại.")
	}
	if errors.Is(err, store.ErrDaCoNhanVat) {
		return guiHopThoai(s, "Bạn đã có nhân vật.")
	}
	if err != nil {
		slog.Error("loi tao nhan vat", "ip", s.DiaChiIP, "loi", err)
		return guiHopThoai(s, "Lỗi hệ thống, vui lòng thử lại.")
	}

	slog.Info("tao nhan vat thanh cong",
		"ip", s.DiaChiIP, "ten", nv.Ten, "lang", nv.Lang, "tai_khoan", b.taiKhoan.ID)
	return b.vaoTheGioi(s, nv)
}

// guiHopThoai gui MOT hop thoai co nut OK (lenh -26). Khac guiTuChoi o login.go:
// guiTuChoi con gui kem -102 de thoat trang thai "dang cho" cua man DANG NHAP;
// o man TAO NHAN VAT, client tu tat vong xoay khi nhan -26 (InfoDlg.hide), nen
// KHONG duoc gui -102 — gui thua se lam sai trang thai client.
func guiHopThoai(s *mynet.Session, cau string) error {
	w := protocol.NewWriter()
	if err := w.UTF(cau); err != nil {
		return err
	}
	return s.Gui(protocol.NewMessage(lenhThongBaoOK, w))
}

// toLowerASCII ha chu thuong CHI cho [A-Z]. Du dung vi ten nhan vat da qua
// reTenNhanVat (chi con [A-Za-z0-9]) — khong can xu ly unicode.
func toLowerASCII(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 'a' - 'A'
		}
	}
	return string(b)
}
```

- [ ] **Step 4: Điều phối lệnh `-28` trong `XuLy`**

Trong `server/internal/game/login.go`, thêm hằng số lệnh và tách `XuLy` thành bộ điều phối. Thêm vào khối hằng số lệnh:

```go
	lenhKhongTrenBanDo int8 = -28 // hop thu "khong tren ban do" — Controller.java messageNotMap
```

Đổi đầu hàm `XuLy` (phần kiểm `m.Command`) thành điều phối:

```go
// XuLy nhan goi tin da giai ma tu tang mang va dieu phoi theo lenh.
func (b *BoXuLy) XuLy(s *mynet.Session, m *protocol.Message) error {
	switch m.Command {
	case lenhChuaDangNhap: // -29
		return b.xuLyDangNhap(s, m)
	case lenhKhongTrenBanDo: // -28
		return b.xuLyKhongTrenBanDo(s, m)
	default:
		// Lenh chua ho tro: bo qua, KHONG log (ly do ngap nhat ky — xem duoi).
		return nil
	}
}

// xuLyKhongTrenBanDo doc byte lenh con roi dieu phoi. Ke hoach 2 chi lam lenh
// con 2 (tao nhan vat); cac lenh con khac la Ke hoach 3.
func (b *BoXuLy) xuLyKhongTrenBanDo(s *mynet.Session, m *protocol.Message) error {
	r := protocol.NewReader(m.Data)
	lenhCon, err := r.Byte()
	if err != nil {
		return fmt.Errorf("goi -28 thieu byte lenh con: %w", err)
	}
	switch lenhCon {
	case lenhConTaoNhanVat: // 2
		return b.xuLyTaoNhanVat(s, r)
	default:
		return nil // lenh con khac la Ke hoach 3, bo qua im lang
	}
}
```

Đổi tên phần thân đăng nhập cũ (từ `r := protocol.NewReader(m.Data)` sau khi đã bỏ nhánh `if m.Command != ...`) thành hàm riêng `xuLyDangNhap`:

```go
// xuLyDangNhap xu ly khung -29 (dang nhap). Goi long hai tang: byte lenh con
// roi moi toi ten/mat khau.
func (b *BoXuLy) xuLyDangNhap(s *mynet.Session, m *protocol.Message) error {
	r := protocol.NewReader(m.Data)
	lenhCon, err := r.Byte()
	if err != nil {
		return fmt.Errorf("goi -29 thieu byte lenh con: %w", err)
	}
	if lenhCon != lenhConDangNhap {
		return nil
	}
	// ... phan con lai giu nguyen: doc ten, mat khau, KiemTraDangNhap, re nhanh ...
}
```

> Giữ nguyên toàn bộ phần đọc tên/mật khẩu và rẽ nhánh có/chưa nhân vật đã viết ở Task 5 — chỉ di chuyển nó vào `xuLyDangNhap`. Bỏ khối ghi chú cũ "Ke hoach 1 chi lam dang nhap" ở đầu (giờ đã có -28).

- [ ] **Step 5: Chạy toàn bộ test package game — phải PASS**

Run (trong `server/`):
```
$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro_test'
D:\Go\bin\go.exe test ./internal/game/ -v
```
Expected: tất cả PASS (test đăng nhập Kế hoạch 1, rẽ nhánh Task 5, tạo nhân vật Task 6). `gofmt -l .` rỗng, `go vet ./...` sạch.

- [ ] **Step 6: Commit**

```bash
git add server/internal/game/nhanvat.go server/internal/game/nhanvat_test.go server/internal/game/login.go
git commit -m "feat(game): xu ly goi tao nhan vat -28/2, kiem ten, vao the gioi (stub)"
```

---

## Task 7: Thử đầu-cuối bằng client giả

**Files:**
- Modify: `server/cmd/thuclient/main.go`
- Create: `server/internal/game/tichhop_test.go` (test tích hợp thật với CSDL `nro_test`)

**Interfaces:**
- Consumes: toàn bộ chuỗi đã dựng (net factory, store CRUD, handler tạo/đăng nhập).

- [ ] **Step 1: Viết test tích hợp thật (đăng ký → tạo nhân vật → đăng nhập lại)**

Tạo `server/internal/game/tichhop_test.go`. Test này dùng **CSDL thật `nro_test`** và `store.Store` thật (không phải `khoGia`), chạy qua `net.Pipe` để giống đường thật nhất. Bỏ qua nếu không có `NRO_TEST_DSN`.

Luồng:
1. Mở `store.Store` thật; tạo một tài khoản thử (mật khẩu băm bằng `auth.Hash`).
2. Dựng phiên máy chủ (`net.Pipe` + `PhucVu` với `func() net.XuLyGoiTin { return NewBoXuLy(khoThat) }`), phía client giả tự bắt tay, giải khung.
3. Client giả gửi gói đăng nhập → kỳ vọng nhận lệnh `2` (chưa có nhân vật).
4. Client giả gửi gói `-28/2` tạo nhân vật tên hợp lệ, làng 0 → kỳ vọng **không** nhận `-26` (không lỗi).
5. Kiểm CSDL: `LayNhanVatTheoTaiKhoan` trả đúng nhân vật vừa tạo.
6. Ngắt, mở phiên mới, đăng nhập lại cùng tài khoản → kỳ vọng **không** nhận lệnh `2` (đã có nhân vật).

> Dùng lại khung `net.Pipe` + bắt gói của `login_test.go`. Bọc chuỗi bắt tay của client giả thành helper nếu tiện. Mọi tên tài khoản/nhân vật thêm hậu tố ngẫu nhiên ngắn để chạy lại không đụng. `t.Cleanup` xóa nhân vật rồi tài khoản.

Khung khẳng định chính:

```go
func TestTaoRoiDangNhapLai_ThayNhanVat(t *testing.T) {
	khoThat := moKhoThatChoTest(t) // *store.Store toi nro_test, skip neu thieu DSN
	ctx := context.Background()

	tenTK := "nguoi" + hauToNganNgau(t)
	matKhau := "matkhau123"
	taoTaiKhoanCoMatKhau(t, khoThat, tenTK, matKhau) // INSERT voi auth.Hash(matKhau), Cleanup xoa

	// Lan 1: dang nhap -> chua co nhan vat -> lenh 2
	{
		cl := moPhienClientGia(t, khoThat)
		defer cl.dong()
		cl.dangNhap(tenTK, matKhau)
		if !cl.choLenh(t, 2) {
			t.Fatal("lan dau phai nhan lenh 2 (man tao nhan vat)")
		}
		tenNV := "Naru" + hauToNganNgau(t)
		cl.taoNhanVat(tenNV, 0, 5)
		if cl.nhanLenhTrongKhoang(t, -26) {
			t.Fatal("tao nhan vat hop le khong duoc tra -26")
		}
		nv, err := khoThat.LayNhanVatTheoTaiKhoan(ctx, layIDTaiKhoan(t, khoThat, tenTK))
		if err != nil || nv.Ten != tenNV {
			t.Fatalf("nhan vat khong duoc luu dung: nv=%+v err=%v", nv, err)
		}
	}

	// Lan 2: dang nhap lai -> da co nhan vat -> KHONG lenh 2
	{
		cl := moPhienClientGia(t, khoThat)
		defer cl.dong()
		cl.dangNhap(tenTK, matKhau)
		if cl.nhanLenhTrongKhoang(t, 2) {
			t.Fatal("da co nhan vat -> khong duoc nhan lenh 2")
		}
	}
}
```

> Các helper `moPhienClientGia`, `dangNhap`, `taoNhanVat`, `choLenh`, `nhanLenhTrongKhoang` là lớp mỏng quanh `net.Pipe` + `protocol.Framer`. Nếu Kế hoạch 1 đã có phần lớn (bắt tay + đăng nhập) trong `login_test.go`, nâng chúng thành helper dùng chung ở đây.

- [ ] **Step 2: Chạy test tích hợp — phải PASS**

Run (trong `server/`):
```
$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro_test'
D:\Go\bin\go.exe test ./internal/game/ -run TestTaoRoiDangNhapLai -v
```
Expected: PASS.

- [ ] **Step 3: Thêm bước tạo nhân vật vào `thuclient` (thử tay)**

Trong `server/cmd/thuclient/main.go`, sau bước đăng nhập thành công, thêm gửi gói `-28/2` với tên/làng/tóc lấy từ cờ dòng lệnh (hoặc hằng số), rồi in ra các gói máy chủ trả về. Mục đích: chạy tay đối chiếu với máy chủ thật. Giữ đúng lối `thuclient` đã có (dùng `protocol.Writer`, `Framer`).

```go
// Sau khi dang nhap thanh cong va nhan duoc lenh 2 (man tao nhan vat), thu tao:
{
	w := protocol.NewWriter()
	w.Byte(2) // lenh con = tao nhan vat
	if err := w.UTF(*tenNhanVat); err != nil { // *tenNhanVat: co dong lenh, mac dinh vi du "Naruto01"
		log.Fatal(err)
	}
	w.Byte(byte(*lang)) // co dong lenh, mac dinh 0
	w.Byte(5)           // kieu toc vi du
	if err := framer.WriteMessage(protocol.NewMessage(-28, w)); err != nil {
		log.Fatal(err)
	}
	log.Printf("da gui goi tao nhan vat: ten=%s lang=%d", *tenNhanVat, *lang)
	// doc va in cac goi tra ve (vi du -26 neu loi, hoac im lang neu thanh cong o Ke hoach 2)
}
```

> Thêm cờ dòng lệnh `tenNhanVat` và `lang` bằng `flag` (giống các cờ sẵn có của `thuclient`). Nếu `thuclient` chưa dùng `flag`, khai báo tối thiểu.

- [ ] **Step 4: Chạy `thuclient` đối chiếu máy chủ thật (thủ công, ghi lại kết quả)**

Mở hai cửa sổ. Cửa sổ 1 chạy máy chủ:
```
$env:NRO_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro'
$env:NRO_MUC_NHAT_KY='info'
D:\Go\bin\go.exe run ./cmd/nroserver
```
Cửa sổ 2 chạy client giả (đăng ký tài khoản trước nếu cần bằng psql/insert), rồi:
```
D:\Go\bin\go.exe run ./cmd/thuclient -tenNhanVat=Naruto01 -lang=0   # cùng các cờ tài khoản sẵn có
```
Expected: máy chủ log `tao nhan vat thanh cong`; kiểm `psql ... -c "select ten,lang from nhan_vat"` thấy hàng mới. Gửi lại tên đó lần nữa (tài khoản khác) → máy chủ trả `-26` "Tên nhân vật đã tồn tại." Ghi kết quả quan sát vào báo cáo task.

- [ ] **Step 5: Chạy TOÀN BỘ test lần cuối**

Run (trong `server/`):
```
$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro_test'
D:\Go\bin\go.exe test ./...
D:\Go\bin\gofmt.exe -l .
D:\Go\bin\go.exe vet ./...
```
Expected: mọi test PASS; `gofmt -l` rỗng; `vet` không in gì.

- [ ] **Step 6: Commit**

```bash
git add server/cmd/thuclient/main.go server/internal/game/tichhop_test.go
git commit -m "test(game): tich hop tao nhan vat -> dang nhap lai thay nhan vat; thuclient tao nhan vat"
```

---

## Tự soi lại kế hoạch

**1. Bám thiết kế:**
- §6.2 "các cột riêng biệt, tra cứu được" → Task 2 dùng cột thật, không JSON gộp. ✅
- §6.2 nhóm người chơi (`nhan_vat`) → Task 2. Các bảng `nhan_vat_ky_nang`, `nhan_vat_nhiem_vu` để kế hoạch sau (đúng ranh giới). ✅
- §7.2 năm làng → Task 3 (đủ 5, client cũ dùng được 3). ✅
- §6.1 "một nhân vật một tài khoản" theo mô hình client cũ → ràng buộc `UNIQUE(tai_khoan_id)`. ✅
- Nguyên tắc "không tin client" → kiểm tên ở tầng Go trước CSDL (Task 6). ✅

**2. Không còn chỗ trống:** mọi bước code đều có mã thật; không có "TODO"/"xử lý lỗi phù hợp".

**3. Nhất quán kiểu:**
- `Lang.ID int16` ↔ `store.NhanVat.Lang int16` ↔ tham số `lang SMALLINT`. ✅
- Handler đọc `lang byte` từ gói rồi ép `int16(langByte)` khi tra `TimLang`. ✅
- `Kho` (game) khớp các phương thức `*store.Store` cài đặt. ✅
- `guiHopThoai` gửi `-26` (dùng `lenhThongBaoOK` đã khai ở `login.go`), KHÔNG `-102`. ✅

**4. Điểm rủi ro cần người soi để mắt:**
- Task 5 sửa test Kế hoạch 1: các test cũ giả định "đăng nhập xong luôn gửi lệnh 2" nay phải nạp `loiLayNhanVat: ErrKhongThayNhanVat`. Nếu bỏ sót, test cũ sẽ đỏ — đó là tín hiệu đúng, không phải hồi quy.
- Helper test (`phienGia`, `chayDangNhap`...) tên chỉ là gợi ý; phải khớp tên thật trong `login_test.go` của Kế hoạch 1.
- `w.Bytes()` ở test giả định `Writer` có hàm xuất byte — kiểm `datastream.go` để dùng đúng tên (có thể là `.Data()` hoặc trả `[]byte` qua cách khác).

---

## Bàn giao thực thi

**Kế hoạch xong, lưu tại `docs/superpowers/plans/2026-07-24-ke-hoach-2-nhan-vat-tao-luu-tai.md`. Hai lựa chọn thực thi:**

**1. Subagent-Driven (khuyến nghị)** — mỗi task một subagent mới, soi giữa các task, lặp nhanh, có kiểm chéo.

**2. Inline** — chạy tuần tự trong phiên này theo executing-plans.

**Chọn cách nào?**
