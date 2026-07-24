# Kế hoạch 5 — Di chuyển + actor khu vực: Kế hoạch thực thi

> **Cho người thực thi:** KỸ NĂNG BẮT BUỘC: dùng superpowers:subagent-driven-development để thực thi từng task. Các bước dùng checkbox (`- [ ]`).

**Mục tiêu:** Nhân vật đi lại được trong bản đồ; máy chủ giữ vị trí thật + chống gian lận mức vừa phải; khu vực chạy theo mô hình actor (một luồng sở hữu trạng thái, nhịp 100ms) — nền cho quái/nhiều người ở kế hoạch sau.

**Kiến trúc:** `KhuVuc` (KH4) nâng thành actor: hộp thư (channel), vòng lặp goroutine, ticker 100ms. Phiên gửi sự kiện (`Vào/DiChuyen/Ra`) vào hộp thư; luồng khu vực xử lý tuần tự + phát tin bằng `s.Gui` (đã an toàn nhờ `muGui`). 5 khu vực cố định (mapId 1-5) tạo lúc khởi động.

**Công nghệ:** Go 1.25+, chỉ `pgx/v5` + stdlib (thêm `time`, `sync` khi cần). Giao thức nhị phân: đã giải mã + soi chéo ở `docs/superpowers/research/kh5-*.md`.

## Ràng buộc chung (Global Constraints)

- **Byte-exact:** gói `-7`/`-5`/`-6` phải khớp bảng trong đặc tả §4 từng byte, từng kiểu (writeInt/Short/Long/UTF). `-5` có 28 trường, hp/hpMax là **writeInt** (4 byte), gender ghi **2 lần**.
- **Server authoritative:** vị trí thật chỉ do luồng khu vực đổi; phiên không tự chạm trạng thái khu vực — chỉ gửi sự kiện.
- **Chống gian lận "vừa phải":** kiểm biên map + khoảng cách Euclid mỗi bước ≤ `NGUONG_BUOC = 500`, áp **mọi map**. Từ chối = **im lặng bỏ gói + giữ vị trí cũ**, KHÔNG ngắt kết nối.
- **Không dội `-7` lại cho chính người vừa đi** — chỉ phát cho người khác trong khu vực (`sendMessAnotherNotMeInMap`).
- **Sạch tranh chấp:** máy này KHÔNG có trình biên dịch C nên `-race` không chạy được (cgo báo "gcc not found"; không quyền admin để cài). Thay bằng: (1) **kỷ luật actor** — mọi trạng thái khu vực (danh sách người, vị trí) CHỈ luồng `vongLap` chạm tới; phiên chỉ gửi sự kiện, không đọc/ghi trực tiếp — kiểm bằng SOI MÃ; (2) **stress test** chạy nhiều lần liên tiếp không treo/flake. (Người dùng có thể chạy `-race` sau trên máy có gcc.)
- **Go:** `D:\Go\bin\go.exe`, gofmt `D:\Go\bin\gofmt.exe`. Windows PowerShell (KHÔNG `&&`). Test cần CSDL: `$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5433/nro_test'` (dòng riêng).
- **Ghi chú tiếng Việt** giải thích *tại sao*. Danh tính Naruto; giữ mapId/đồ hoạ cũ chỉ như bản mã.
- File dưới `server/`. Không đụng `1-MAY-CHU-GAME/`.

---

### Task 1: Bộ khung actor khu vực (hộp thư + vòng lặp + ticker)

**Files:**
- Modify: `server/internal/game/khuvuc.go`
- Modify: `server/internal/game/khuvuc_test.go`

**Interfaces:**
- Produces: kiểu sự kiện `type suKienKhuVuc interface{}` với các struct `suKienVao{nv *store.NhanVat; phien guiTin}`, `suKienDiChuyen{idNhanVat int64; x, y int16}`, `suKienRa{idNhanVat int64}`. Interface `guiTin` = `interface{ Gui(*protocol.Message) error }` (để test dùng phiên giả; `*mynet.Session` thoả mãn).
- `KhuVuc` thêm: `hopThu chan suKienKhuVuc`, `dung chan struct{}`, và (nội bộ) `map[int64]*nguoiTrongKhuVuc`.
- `func (k *KhuVuc) BatDau()` khởi động goroutine vòng lặp; `func (k *KhuVuc) Gui(sk suKienKhuVuc)` đẩy sự kiện (không chặn — `select` với default bỏ gói khi đầy); `func (k *KhuVuc) Dung()` đóng actor sạch.

- [ ] **Step 1: Viết test đỏ** — actor khởi động, nhận một sự kiện `suKienRa` của id không tồn tại (no-op an toàn), rồi `Dung()` không treo, không rò goroutine. Chạy với `-race`.

```go
func TestKhuVucActorKhoiDongVaDungSach(t *testing.T) {
	k := MoKhuVuc(TimBanDo(1))
	k.BatDau()
	k.Gui(suKienRa{idNhanVat: 999}) // no-op: khong co ai
	k.Dung()                         // phai tra ve, khong treo
	// gui sau khi dung khong duoc panic:
	k.Gui(suKienRa{idNhanVat: 1})
}
```

- [ ] **Step 2: Chạy — fail** (`BatDau/Gui/Dung` chưa có). `D:\Go\bin\go.exe test ./internal/game/ -run TestKhuVucActor -race -v`
- [ ] **Step 3: Cài đặt** vòng lặp: `select { case sk := <-hopThu: xuLy(sk); case <-ticker.C: nhipTim(); case <-dung: return }`. `Gui` dùng `select` non-blocking (default = bỏ gói, ghi log debug). `Dung` đóng kênh `dung` một lần (dùng `sync.Once`), goroutine thấy kênh đóng thì thoát; sau `Dung`, `Gui` không được panic (kiểm cờ đã đóng qua `select` trên `dung` trong `Gui`). Ticker 100ms; `nhipTim` KH5 rỗng (chỗ cho KH6) — ghi chú.
- [ ] **Step 4: Chạy — PASS** với `-race`.
- [ ] **Step 5: Commit** — `feat(game): KhuVuc thanh actor (hop thu + vong lap + ticker 100ms)`

---

### Task 2: Kiểm hợp lệ di chuyển (hàm thuần)

**Files:**
- Create: `server/internal/game/dichuyen.go`
- Create: `server/internal/game/dichuyen_test.go`

**Interfaces:**
- Produces: `const NguongBuoc = 500`; `func kiemDiChuyenHopLe(bd *BanDo, cuX, cuY, toX, toY int16) bool` — true nếu đích trong biên + khoảng cách ≤ NguongBuoc.
- Biên map: KH5 dùng biên rộng an toàn từ hằng số (chưa nạp kích thước tile thật) — `const BienRongToiDa = 4096` cho cả X,Y; ghi chú TODO nạp kích thước map thật (đặc tả §5). toX,toY phải `≥ 0` và `≤ BienRongToiDa`.

- [ ] **Step 1: Viết test đỏ** — trong biên + gần: true; ngoài biên (âm / > BienRongToiDa): false; vượt NguongBuoc (khoảng cách > 500): false; đúng bằng ngưỡng: true.

```go
func TestKiemDiChuyen(t *testing.T) {
	bd := TimBanDo(1)
	if !kiemDiChuyenHopLe(bd, 100, 100, 300, 100) { t.Error("buoc ngan hop le phai true") }
	if kiemDiChuyenHopLe(bd, 100, 100, -5, 100) { t.Error("toX am phai false") }
	if kiemDiChuyenHopLe(bd, 100, 100, 5000, 100) { t.Error("ngoai bien phai false") }
	if kiemDiChuyenHopLe(bd, 0, 0, 400, 300) { t.Error("khoang cach 500 (=nguong) phai true... kiem lai") } // 400²+300²=500²
	if kiemDiChuyenHopLe(bd, 0, 0, 600, 0) { t.Error("khoang cach 600 > 500 phai false") }
}
```
*(sửa assert cho đúng: khoảng cách 500 = true, 600 = false.)*

- [ ] **Step 2: Chạy — fail.** `-run TestKiemDiChuyen -v`
- [ ] **Step 3: Cài đặt** — khoảng cách Euclid: `dx*dx+dy*dy ≤ NguongBuoc*NguongBuoc` (so bình phương, tránh sqrt/float). Kiểm biên trước.
- [ ] **Step 4: Chạy — PASS.**
- [ ] **Step 5: Commit** — `feat(game): kiemDiChuyenHopLe (bien + khoang cach/buoc, moi map)`

---

### Task 3: Dựng gói `-7` (phát) và `-6` (rời)

**Files:**
- Modify: `server/internal/game/khuvuc.go` (hoặc `dichuyen.go`), `+ _test.go`

**Interfaces:**
- Produces: `func goiDiChuyen(idNhanVat int64, x, y int16) *protocol.Message` (lệnh -7: writeInt id, writeShort x, writeShort y); `func goiRoiKhuVuc(idNhanVat int64) *protocol.Message` (lệnh -6: writeInt id). Thêm hằng `lenhDiChuyen int8 = -7`, `lenhRoiKhuVuc int8 = -6`.

- [ ] **Step 1: Test đỏ byte-đúng** — đọc lại bằng Reader: -7 → Int(id), Short(x), Short(y), hết buffer; -6 → Int(id), hết buffer.
- [ ] **Step 2: fail.** **Step 3:** cài đặt. **Step 4: PASS.**
- [ ] **Step 5: Commit** — `feat(game): goi -7 di chuyen + -6 roi khu vuc (byte-dung)`

---

### Task 4: Dựng gói `-5` xuất hiện (28 trường, byte-đúng)

**Files:**
- Modify: `server/internal/game/khuvuc.go` (hoặc file mới `xuathien.go`), `+ _test.go`

**Interfaces:**
- Produces: `func goiXuatHien(nv *store.NhanVat) (*protocol.Message, error)` — lệnh -5, 28 trường theo đặc tả §4.4. Hằng `lenhXuatHien int8 = -5`. Giá trị: clanId=-1, level=byte(nv.Cap) nếu có else 1, gender×2=byte(nv.Lang), head=nv.KieuToc, name=nv.Ten, hp=hpMax=**writeInt** nv.HpToiDa, còn lại 0/hằng theo bảng. Bỏ `-103`/`-8`.

- [ ] **Step 1: Test đỏ byte-đúng** — đọc lại toàn bộ 28 trường đúng thứ tự/kiểu, khẳng định gender ghi 2 lần, hp/hpMax là Int, hết buffer sau `hat`.
- [ ] **Step 2: fail.** **Step 3:** cài đặt (chép thứ tự bảng §4.4 chính xác). **Step 4: PASS.**
- [ ] **Step 5: Commit** — `feat(game): goi -5 xuat hien nguoi choi (28 truong byte-dung)`

---

### Task 5: Actor xử lý Vào / DiChuyen / Ra (phát tin)

**Files:**
- Modify: `server/internal/game/khuvuc.go`, `khuvuc_test.go`

**Interfaces:**
- `nguoiTrongKhuVuc{ nv *store.NhanVat; phien guiTin }`. `xuLy(sk)`:
  - `suKienVao`: thêm vào map; gửi `goiXuatHien` của TỪNG người cũ cho người mới; phát `goiXuatHien` người mới cho người cũ.
  - `suKienDiChuyen`: cập nhật `nv.ViTriX/Y` của người đó (authoritative); phát `goiDiChuyen` cho NGƯỜI KHÁC (không gửi lại chính họ).
  - `suKienRa`: gỡ khỏi map; phát `goiRoiKhuVuc` cho người còn lại.
- Hàm phát trợ giúp: `phatChoNguoiKhac(trừ id, m)` và `guiCho(id, m)`.

- [ ] **Step 1: Test đỏ** (phiên giả ghi vào slice, chạy `-race`): 2 người Vào cùng khu vực → mỗi người nhận `-5` của người kia. DiChuyen của A → B nhận `-7 [idA][x][y]`, A KHÔNG nhận. Ra của A → B nhận `-6 [idA]`.

```go
type phienGia struct{ nhan []*protocol.Message }
func (p *phienGia) Gui(m *protocol.Message) error { p.nhan = append(p.nhan, m); return nil }
```
*(Lưu ý: vì actor xử lý bất đồng bộ, test phải đồng bộ hoá — gửi sự kiện rồi chờ xử lý xong. Dùng một sự kiện "đồng bộ" gửi vào actor trả về qua channel, hoặc cho `Gui` của actor nhận thêm callback báo xong. Đơn giản nhất: thêm hàm test-only `xuLyDongBo(sk)` gọi trực tiếp `xuLy` không qua channel — nhưng vẫn phải sạch race. Quyết định: test gọi `k.xuLy(sk)` trực tiếp trong luồng test (actor CHƯA BatDau), khẳng định trạng thái + gói phiên giả nhận. Cách này test logic xử lý tách khỏi vòng lặp channel — vòng lặp đã test ở Task 1.)*

- [ ] **Step 2: fail.** **Step 3:** cài đặt `xuLy`. **Step 4: PASS** với `-race`.
- [ ] **Step 5: Commit** — `feat(game): actor khu vuc xu ly Vao/DiChuyen/Ra + phat -5/-7/-6`

---

### Task 6: Sổ khu vực + tạo 5 actor lúc khởi động

**Files:**
- Create: `server/internal/game/sokhuvuc.go`, `+ _test.go`
- Modify: nơi khởi tạo máy chủ (`cmd/nroserver` hoặc chỗ dựng `BoXuLy`/Kho) — task này CHỈ tạo sổ + tra; nối vào phiên ở Task 7.

**Interfaces:**
- Produces: `type SoKhuVuc struct{ theoMap map[int32]*KhuVuc }`; `func MoSoKhuVuc() *SoKhuVuc` tạo 5 khu vực mapId 1-5 (dùng `TimBanDo`), gọi `BatDau()` mỗi cái; `func (s *SoKhuVuc) Tra(mapId int32) *KhuVuc`; `func (s *SoKhuVuc) DungTatCa()`.

- [ ] **Step 1: Test đỏ** — `MoSoKhuVuc` có đủ 5 khu vực; `Tra(1..5)` không nil, đúng bản đồ; `Tra(99)` nil; `DungTatCa` không treo.
- [ ] **Step 2: fail.** **Step 3:** cài đặt. **Step 4: PASS.**
- [ ] **Step 5: Commit** — `feat(game): SoKhuVuc tao 5 actor mapId 1-5 luc khoi dong`

---

### Task 7: Nối bộ điều phối — `-39` (vào khu vực) + `-7` (di chuyển)

**Files:**
- Modify: `server/internal/game/login.go` (switch xử lý gói + `BoXuLy`), `login_test.go`
- Modify: nơi dựng `BoXuLy` để truyền `*SoKhuVuc` vào (thay `MoKhuVuc` per-phiên của KH4 bằng tra sổ chung).

**Interfaces:**
- `BoXuLy` thêm `soKhuVuc *SoKhuVuc`. `xuLyClientOk` (KH4): BỎ `b.khuVuc = MoKhuVuc(bd)`, THAY bằng `b.khuVuc = b.soKhuVuc.Tra(nv.MapID)` (khu vực CHUNG). Chưa gửi `suKienVao` ở đây — đợi `-39`.
- `case -39`: gửi `suKienVao{nv, s}` vào `b.khuVuc`. `case -7`: đọc byte cờ + short toX + short toY (toY optional); `kiemDiChuyenHopLe`; nếu hợp lệ gửi `suKienDiChuyen`; nếu không, bỏ qua im lặng.
- Khi phiên đóng (defer trong vòng phục vụ): nếu `b.khuVuc != nil && b.nhanVat != nil`, gửi `suKienRa{b.nhanVat.ID}`. (Cần điểm móc "khi phiên kết thúc" — xác định trong login.go/net.)

- [ ] **Step 1: Test đỏ** — dùng `dungPhienThat` + phiên thứ 2 vào cùng mapId 1: sau `-39` cả hai, mỗi phiên nhận `-5` của người kia; `-7` của A → B nhận, A không; `-7` vượt ngưỡng của A → B không nhận, `nv.ViTriX/Y` của A không đổi.
- [ ] **Step 2: fail.** **Step 3:** nối switch + đổi `xuLyClientOk` + móc đóng phiên. **Step 4: PASS** (kể cả test KH1-4 cũ; cập nhật test KhuVuc/xuLyClientOk cũ nếu chữ ký đổi).
- [ ] **Step 5: Commit** — `feat(game): dieu phoi -39 vao khu vuc + -7 di chuyen, dung so khu vuc chung`

---

### Task 8: Test đầu-cuối 2 người chơi (mạng + CSDL thật)

**Files:**
- Modify: `server/internal/game/dongbo_tichhop_test.go` (hoặc file test tích hợp mới)

- [ ] **Step 1: Viết test** — 2 client giả tạo 2 nhân vật ở mapId 1, đi trọn version dance + `-28/13` + `-39`. Khẳng định: mỗi client nhận `-5` của người kia (byte-đúng, đối chiếu nv từ DB). Client A gửi `-7 [cờ][x][y]` → client B nhận `-7 [idA][x][y]`. Client A gửi `-7` đích cách > 500 → B KHÔNG nhận trong hạn ngắn; đọc lại nhân vật A từ DB (hoặc trạng thái server) xác nhận vị trí không nhảy.
- [ ] **Step 2: Chạy** với DSN → PASS. Chạy `go test ./... -count=1` toàn xanh + `-race` trên gói game.
- [ ] **Step 3: Commit** — `test(game): dau-cuoi 2 nguoi choi thay nhau di chuyen + chan teleport`

---

### Task 9 (THỦ CÔNG — người dùng): xác nhận bằng client Unity thật

Không có code. Người dùng mở 2 client thật (hoặc 1), đăng nhập vào cùng làng, đi lại và xác nhận thấy nhân vật di chuyển mượt. Báo lại nếu lệch để tinh chỉnh.

---

## Tự soi (sau khi viết)

- **Phủ đặc tả:** actor (T1), kiểm gian lận (T2), gói -7/-6 (T3), -5 (T4), xử lý sự kiện (T5), sổ+khởi động (T6), điều phối -39/-7 (T7), đầu-cuối (T8). ✔ mọi mục §4-§9 có task.
- **Điểm rủi ro:** T5 đồng bộ hoá test actor (đã ghi quyết định gọi `xuLy` trực tiếp); T7 đổi `xuLyClientOk` KH4 + móc đóng phiên (có thể cần đọc net/session để tìm điểm móc — implementer báo NEEDS_CONTEXT nếu không rõ). T6/T7 nối vào khởi động máy chủ — cần xem `cmd/nroserver` dựng `BoXuLy` thế nào.
- **Không placeholder** ngoài các "sửa assert" đã chỉ rõ trong test mẫu.
