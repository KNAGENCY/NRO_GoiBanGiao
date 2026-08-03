# Lưới chắn panic — Kế hoạch thực thi

> **For agentic workers:** REQUIRED SUB-SKILL: dùng `superpowers:subagent-driven-development`
> để thi công từng task. Các bước dùng checkbox (`- [ ]`).
>
> **Đọc trước:** đặc tả [`../specs/2026-08-03-luoi-chan-panic-design.md`](../specs/2026-08-03-luoi-chan-panic-design.md),
> `BAT-DAU-TU-DAY.md` §6 (lỗi đã vấp) + §9 (bẫy môi trường), `AGENTS.md`.

**Goal:** Đặt lưới đỡ `recover()` ở cả 6 goroutine của máy chủ để một panic bất kỳ không
làm sập cả tiến trình — thay vào đó ghi log + cô lập trong phạm vi hẹp nhất.

**Architecture:** Một hàm trợ giúp dùng chung `antoan.ChayCoLuoi` (recover + log kèm stack).
Caller quyết định làm gì sau khi chặn, theo 3 chính sách: actor khu vực **chạy tiếp** (bọc
từng sự kiện, kèm van an toàn chặn hot-loop); mỗi kết nối **tự đóng**; vòng lưu **giữ sống**.

**Tech Stack:** Go 1.25 + stdlib (`log/slog`, `runtime/debug`). Không thêm thư viện ngoài.

## Global Constraints

- **`slog`** là logger chuẩn (đã dùng ở `listener.go`, KH11). Mọi panic bị chặn PHẢI ghi
  log kèm `stack` — **không nuốt lỗi im lặng** (trái bẫy §6 bản Java).
- **Không đổi hành vi lúc-không-có-lỗi.** Lưới chỉ hoạt động khi panic; đường chạy bình
  thường y hệt trước — mọi test cũ phải vẫn xanh.
- **Máy này KHÔNG chạy `go test -race`** (thiếu trình biên dịch C) → dùng `-count=5`.
- **Go:** `D:\Go\bin\go.exe`, `D:\Go\bin\gofmt.exe`, `D:\Go\bin\go.exe vet`. PowerShell
  KHÔNG có `&&`.
- **Cổng PostgreSQL tự dò**, đừng ghim: `Get-NetTCPConnection -State Listen | ? { $_.LocalPort -in 5432,5433 }`;
  chưa chạy thì `D:\PostgreSQL\bin\pg_ctl.exe -D D:\PostgreSQL\data start`. DSN test:
  `$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:<CONG>/nro_test'` (dòng riêng).
- **Chú thích tiếng Việt giải thích *tại sao***. Nhánh: `luoi-chan-panic` (đã có spec commit).
- **Kiểm thử đột biến bắt buộc** sau mỗi task có lưới: gỡ lưới → test tương ứng **phải đỏ
  vì panic lan ra** (không phải lỗi biên dịch). Không đỏ = viết lại test.

---

## Bản đồ file

| File | Trách nhiệm |
|---|---|
| `server/internal/antoan/luoi.go` (tạo) | Hàm `ChayCoLuoi` — recover + log. Chỉ phụ thuộc stdlib |
| `server/internal/antoan/luoi_test.go` (tạo) | Test hàm lưới thuần |
| `server/internal/game/khuvuc.go` (sửa) | Bọc `xuLy`/`nhipTim`/`xuLyDieuKhien` trong `vongLap` + van an toàn |
| `server/internal/game/khuvuc_test.go` (sửa) | Test khu vực chạy tiếp sau panic + van |
| `server/internal/net/session.go` (sửa) | Bọc thân `PhucVu` (:196) + `vongWriter` (:142) |
| `server/internal/net/session_test.go` (sửa) | Test một kết nối lỗi tự đóng, kết nối khác sống |
| `server/internal/game/luudinhky.go` (sửa) | Bọc thân vòng lặp `vongLap` (:125) |
| `server/internal/game/thogi.go` (sửa) | Bọc thân `vongLap` (:104) + `go func(nv)` (:238) |
| `server/internal/game/thogi_test.go` (sửa) | Test một bản ghi lỗi bị bỏ, người kế tiếp vẫn lưu |

Vì sao package riêng `internal/antoan`: cả `net` lẫn `game` đều dùng lưới; đặt ở `game` thì
`net` phải import `game` (vòng import). `antoan` chỉ phụ thuộc stdlib nên ai import cũng được.

---

### Task 1: Hàm lưới đỡ dùng chung `antoan.ChayCoLuoi`

**Files:**
- Create: `server/internal/antoan/luoi.go`
- Test: `server/internal/antoan/luoi_test.go`

**Interfaces:**
- Produces: `func ChayCoLuoi(cho string, f func()) (daChanPanic bool)` — chạy `f`; nếu `f`
  panic thì chặn, ghi `slog.Error` kèm `cho`/`loi`/`stack`, trả `true`; chạy trót trả `false`.

- [ ] **Step 1: Viết test đỏ** — `luoi_test.go`:

```go
package antoan

import "testing"

func TestChayCoLuoi_ChanPanicVaKhongLanRa(t *testing.T) {
	// f panic → ChayCoLuoi PHAI chan lai (test khong duoc chet) va tra true.
	daChan := ChayCoLuoi("test.panic", func() { panic("bug gia") })
	if !daChan {
		t.Fatal("panic phai bi chan va tra true")
	}
}

func TestChayCoLuoi_ChayTrotTraFalse(t *testing.T) {
	chay := false
	daChan := ChayCoLuoi("test.ok", func() { chay = true })
	if daChan {
		t.Fatal("chay trot khong duoc bao la chan panic")
	}
	if !chay {
		t.Fatal("f phai duoc chay")
	}
}
```

- [ ] **Step 2: Chạy — fail** (`ChayCoLuoi` chưa có).
  `D:\Go\bin\go.exe test ./internal/antoan/ -run TestChayCoLuoi -v`
- [ ] **Step 3: Cài đặt** `luoi.go`:

```go
// Package antoan: luoi do recover dung chung cho moi goroutine. KHONG phu thuoc
// package game/net (tranh vong import) — chi stdlib.
package antoan

import (
	"log/slog"
	"runtime/debug"
)

// ChayCoLuoi chay f; neu f panic thi CHAN lai, ghi log kem dau vet, roi tra ve
// binh thuong — KHONG cho panic lan len lam sap tien trinh. Tra ve true neu da
// chan mot panic (de caller dem/quyet dinh), false neu f chay tron.
func ChayCoLuoi(cho string, f func()) (daChanPanic bool) {
	defer func() {
		if r := recover(); r != nil {
			daChanPanic = true
			slog.Error("panic bi chan boi luoi",
				"cho", cho, "loi", r, "stack", string(debug.Stack()))
		}
	}()
	f()
	return false
}
```

- [ ] **Step 4: Chạy — PASS.**
- [ ] **Step 5: Kiểm thử đột biến** — đổi thân thành gọi thẳng `f()` (bỏ `defer recover`)
  → `TestChayCoLuoi_ChanPanicVaKhongLanRa` phải **đỏ vì panic lan ra làm test chết**.
  Hoàn nguyên, xanh lại.
- [ ] **Step 6: Commit** — `feat(antoan): luoi do recover dung chung ChayCoLuoi`

---

### Task 2: Actor khu vực — bọc từng sự kiện, khu vực chạy tiếp sau panic

**Files:**
- Modify: `server/internal/game/khuvuc.go` — hàm `vongLap` (dòng 479-505)
- Modify: `server/internal/game/khuvuc_test.go`

**Interfaces:**
- Consumes: `antoan.ChayCoLuoi` (Task 1).
- Bối cảnh: `vongLap` là `for { select {...} }`; nhánh `case sk := <-k.hopThu: k.xuLy(sk)`
  và `case <-ticker.C: k.nhipTim()` và hai nhánh `k.xuLyDieuKhien(sk)`.

- [ ] **Step 1: Viết test đỏ** — khu vực nhận một sự kiện làm `xuLy` panic, rồi PHẢI xử
  được sự kiện KẾ TIẾP (chứng minh vòng lặp còn sống). Thêm vào `khuvuc_test.go`:

```go
// Dung mot su kien "test" lam xuLy panic. Neu goi trong game khong co san su
// kien nay, them mot loai su kien noi bo chi dung trong test (build tag hoac
// mot truong co san). Cach don gian: mot suKien co ham func() de vongLap goi.
func TestKhuVuc_PanicTrongSuKien_KhuVucVanChayTiep(t *testing.T) {
	k := khuVucTest(t) // helper dung san trong package (xem cac test khuvuc hien co)
	k.BatDau()
	defer k.Dung()

	daGoiKeTiep := make(chan struct{}, 1)
	// Su kien 1: panic. Su kien 2: bao hieu con song.
	k.Gui(suKienChay{f: func() { panic("bug danh nhau gia") }})
	k.Gui(suKienChay{f: func() { daGoiKeTiep <- struct{}{} }})

	select {
	case <-daGoiKeTiep:
		// PASS: vong lap song sot qua panic, xu duoc su kien ke tiep.
	case <-time.After(2 * time.Second):
		t.Fatal("khu vuc chet sau panic — vong lap khong chay tiep")
	}
}
```

*Ghi chú thiết kế cho người thi công:* cần một **điểm chèn test** để bơm panic mà không
bẩn code thật. Đề xuất: thêm một loại sự kiện `suKienChay struct{ f func() }` mà `xuLy`
gọi `sk.f()` — hữu ích chung cho test, không phải nhánh `if test`. Nếu package đã có cách
bơm hành vi vào actor thì dùng lại, đừng đẻ trùng.

- [ ] **Step 2: Chạy — fail** (chưa bọc lưới → panic làm cả test binary chết / khu vực chết).
- [ ] **Step 3: Cài đặt** — trong `vongLap`, bọc mỗi lần xử lý bằng `antoan.ChayCoLuoi`:

```go
case sk := <-k.hopThuDieuKhien:
	var dung bool
	antoan.ChayCoLuoi("khuvuc.dieuKhien", func() { dung = k.xuLyDieuKhien(sk) })
	if dung {
		return
	}
case sk := <-k.hopThu:
	antoan.ChayCoLuoi("khuvuc.xuLy", func() { k.xuLy(sk) })
case <-ticker.C:
	antoan.ChayCoLuoi("khuvuc.nhipTim", k.nhipTim)
```

(Áp cho cả nhánh `hopThuDieuKhien` trong khối `select` đầu — drain control — y hệt.)
**Không** bọc ngoài cả `vongLap` — bọc ngoài thì panic vẫn thoát vòng, giết khu vực.

- [ ] **Step 4: Chạy — PASS** với `-count=5`.
- [ ] **Step 5: Kiểm thử đột biến** — bỏ `ChayCoLuoi` ở nhánh `xuLy` (gọi thẳng `k.xuLy(sk)`)
  → test phải **đỏ vì panic lan ra**. Hoàn nguyên.
- [ ] **Step 6: Soi mã** — xác nhận không nhánh xử lý nào trong `vongLap` còn gọi trần.
- [ ] **Step 7: Commit** — `feat(game): actor khu vuc chay tiep sau panic (boc tung su kien)`

---

### Task 3: Van an toàn — khu vực panic dai dẳng thì dừng riêng nó

**Files:**
- Modify: `server/internal/game/khuvuc.go` — thêm bộ đếm panic liên tiếp vào `KhuVuc` + `vongLap`
- Modify: `server/internal/game/khuvuc_test.go`

**Interfaces:**
- Consumes: `antoan.ChayCoLuoi` trả `bool` (Task 1) để biết nhịp vừa rồi có panic không.
- `KhuVuc` thêm trường `soPanicLienTiep int` (chỉ luồng `vongLap` chạm — an toàn).

**Hằng số:** `nguongPanicDungKhuVuc = 100` (≈ 10 giây panic liền ở nhịp 100ms).

- [ ] **Step 1: Viết test đỏ** — (a) panic **liên tiếp quá ngưỡng** → khu vực dừng; (b)
  panic **thoáng qua** (vài lần rồi thôi) → khu vực **không** dừng:

```go
func TestKhuVuc_PanicDaiDang_DungRiengKhuVuc(t *testing.T) {
	k := khuVucTest(t)
	k.BatDau()
	// Nhip tim luon panic → sau <= nguong+ vai nhip, vong lap phai thoat.
	k.datNhipTimTest(func() { panic("hong dai dang") }) // helper bom nhipTim test
	select {
	case <-k.xong: // vongLap dong k.xong khi return
		// PASS: khu vuc da tu dung.
	case <-time.After(15 * time.Second):
		t.Fatal("khu vuc panic dai dang nhung khong dung — van an toan khong chay")
	}
}

func TestKhuVuc_PanicThoangQua_KhongDung(t *testing.T) {
	k := khuVucTest(t)
	k.BatDau()
	defer k.Dung()
	dem := 0
	k.datNhipTimTest(func() {
		dem++
		if dem <= 3 { // panic 3 nhip dau roi thoi
			panic("thoang qua")
		}
	})
	// Sau vai tram ms, khu vuc phai VAN song → gui su kien va nhan phan hoi.
	daSong := make(chan struct{}, 1)
	time.Sleep(500 * time.Millisecond)
	k.Gui(suKienChay{f: func() { daSong <- struct{}{} }})
	select {
	case <-daSong: // PASS
	case <-time.After(2 * time.Second):
		t.Fatal("loi thoang qua khong duoc lam dung khu vuc")
	}
}
```

- [ ] **Step 2: Chạy — fail** (van chưa có / trường chưa có).
- [ ] **Step 3: Cài đặt** — trong `vongLap`, sau mỗi lần `ChayCoLuoi` cho `xuLy`/`nhipTim`:

```go
case <-ticker.C:
	if antoan.ChayCoLuoi("khuvuc.nhipTim", k.nhipTim) {
		k.soPanicLienTiep++
		if k.soPanicLienTiep >= nguongPanicDungKhuVuc {
			slog.Error("khu vuc hong dai dang, dung rieng no de khong ghi chet CPU",
				"map", k.MapID, "so_panic", k.soPanicLienTiep)
			return // thoat vongLap → defer close(k.xong) chay
		}
	} else {
		k.soPanicLienTiep = 0 // mot nhip trot lot → reset, loi thoang qua khong cham nguong
	}
```

Áp cùng nguyên tắc cho nhánh `xuLy` (đếm chung `soPanicLienTiep`, reset khi có nhịp trót).

- [ ] **Step 4: Chạy — PASS** với `-count=5`.
- [ ] **Step 5: Kiểm thử đột biến** — bỏ dòng reset `k.soPanicLienTiep = 0`
  → `TestKhuVuc_PanicThoangQua_KhongDung` phải **đỏ** (lỗi thoáng qua tích luỹ tới ngưỡng
  làm dừng nhầm). Hoàn nguyên.
- [ ] **Step 6: Commit** — `feat(game): van an toan dung khu vuc panic dai dang (reset khi trot)`

---

### Task 4: Mỗi kết nối — bọc `PhucVu` + `vongWriter`, một kết nối lỗi tự đóng

**Files:**
- Modify: `server/internal/net/session.go` — `PhucVu` (:196) và `vongWriter` (:142)
- Modify: `server/internal/net/session_test.go`

**Interfaces:**
- Consumes: `antoan.ChayCoLuoi` (Task 1).
- Bối cảnh: `PhucVu` đã có `defer s.Dong()` + `defer c.KhiDong()` → recover ở TRONG `PhucVu`
  (sau các defer) vẫn để defer chạy = đóng kết nối sạch. `vongWriter` có `defer close(s.writerXong)`.

- [ ] **Step 1: Viết test đỏ** — một `XuLyGoiTin` giả panic khi nhận một gói → `PhucVu`
  của kết nối đó kết thúc sạch (không làm test chết), và có thể mở kết nối khác chạy bình
  thường. Thêm vào `session_test.go`:

```go
func TestPhucVu_XuLyPanic_DongKetNoiKhongLamSap(t *testing.T) {
	// taoXuLy tra ve bo xu ly panic ngay khi xu goi dau tien.
	taoXuLy := func() XuLyGoiTin { return xuLyPanicGia{} }
	srv, dia := moServerTest(t, taoXuLy) // helper dung san trong package net
	defer srv.Dong()

	c := noiToiTest(t, dia)
	batTayClientTest(t, c)
	guiGoiBatKyTest(t, c) // kich xuLy panic

	// Ket noi bi dong (doc tra EOF), va TIEN TRINH TEST VAN SONG.
	if !ketNoiDaDongTrongHan(t, c, 2*time.Second) {
		t.Fatal("ket noi gap panic phai duoc dong")
	}
	// Mo ket noi thu hai — phai bat tay duoc binh thuong (server con song).
	c2 := noiToiTest(t, dia)
	batTayClientTest(t, c2) // khong duoc panic/timeout
}
```

*Ghi chú:* nếu package `net` chưa có helper dựng server test, tạo tối thiểu (một
`net.Listen` + `LangNghe` trên cổng ngẫu nhiên). `xuLyPanicGia` implement `XuLyGoiTin`,
`XuLy` panic.

- [ ] **Step 2: Chạy — fail** (panic trong `PhucVu` làm cả test binary chết).
- [ ] **Step 3: Cài đặt** — đầu `PhucVu`, ngay sau các `defer` đã có, thêm:

```go
defer func() {
	if r := recover(); r != nil {
		slog.Error("panic bi chan trong PhucVu, dong ket noi",
			"ip", diaChi, "loi", r, "stack", string(debug.Stack()))
	}
}()
```

(Đặt SAU `defer s.Dong()` để khi panic, `s.Dong()` vẫn chạy — Go chạy defer theo thứ tự
ngược, recover ở defer trên cùng vẫn bắt được panic từ thân.) Với `vongWriter`, bọc thân
vòng `for` bằng cùng mẫu recover (log rồi return → `defer close(s.writerXong)` chạy, kết
nối được dọn qua đường đóng phiên đã có).

- [ ] **Step 4: Chạy — PASS** với `-count=5`.
- [ ] **Step 5: Kiểm thử đột biến** — bỏ `defer recover` trong `PhucVu` → test phải **đỏ**.
- [ ] **Step 6: Commit** — `feat(net): luoi panic moi ket noi (PhucVu + vongWriter tu dong)`

---

### Task 5: Hệ lưu — bọc vòng lưu định kỳ + thợ ghi + mỗi lần lưu

**Files:**
- Modify: `server/internal/game/luudinhky.go` — `vongLap` (:125)
- Modify: `server/internal/game/thogi.go` — `vongLap` (:104) + `go func(nv)` (:238)
- Modify: `server/internal/game/thogi_test.go`

**Interfaces:**
- Consumes: `antoan.ChayCoLuoi` (Task 1).

- [ ] **Step 1: Viết test đỏ** — thợ ghi nhận hai bản: bản đầu làm hàm ghi panic, bản hai
  bình thường → **bản hai vẫn được lưu**. Thêm vào `thogi_test.go`:

```go
func TestThoGhi_MotBanGhiPanic_NguoiKeTiepVanLuu(t *testing.T) {
	kho := khoLuuGia{ // kho gia: ghi id X thi panic, con lai luu binh thuong
		ghi: func(nv store.NhanVat) error {
			if nv.ID == 1 {
				panic("ghi hong gia")
			}
			return nil
		},
	}
	t2 := MoThoGhi(kho, 1)
	defer t2.Dung()
	t2.Xep(store.NhanVat{ID: 1}) // panic
	t2.Xep(store.NhanVat{ID: 2}) // phai luu duoc

	if !khoLuuGhiDuocTrongHan(&kho, 2, 2*time.Second) {
		t.Fatal("ban ghi ke tiep phai duoc luu du ban truoc panic")
	}
}
```

- [ ] **Step 2: Chạy — fail** (panic trong thợ ghi làm chết vòng / test).
- [ ] **Step 3: Cài đặt** — trong `ThoGhi.vongLap`, bọc phần xử một lượt ghi bằng
  `antoan.ChayCoLuoi("thogi.ghi", ...)`; trong `go func(nv)` (:238) bọc thân bằng cùng mẫu
  recover; trong `luudinhky.vongLap` bọc thân mỗi vòng bằng `antoan.ChayCoLuoi("luudinhky.nhip", ...)`.
  Bọc **thân vòng**, KHÔNG bọc ngoài (bọc ngoài → panic giết cả vòng, mất tính năng lưu).

- [ ] **Step 4: Chạy — PASS** với `-count=5`.
- [ ] **Step 5: Kiểm thử đột biến** — bỏ lưới ở `ThoGhi` → test phải **đỏ**.
- [ ] **Step 6: Commit** — `feat(game): luoi panic he luu (luu dinh ky + tho ghi + moi ban ghi)`

---

### Task 6: Nghiệm thu toàn bộ + ghi sổ

**Files:**
- Modify: `BAT-DAU-TU-DAY.md` (§4, §8), `NHAT-KY-AGENT.md`

- [ ] **Step 1: Quét sót** — `grep -rn "go func\|go [a-z].*vongLap\|go PhucVu\|go s\.vongWriter" server/internal/`
  → xác nhận **6/6** điểm goroutine đều đã qua lưới. Không còn điểm trần.
- [ ] **Step 2: Toàn bộ test** — dò cổng, `$env:NRO_TEST_DSN=...`, rồi
  `D:\Go\bin\go.exe test ./... -count=5` **toàn xanh**; `D:\Go\bin\go.exe vet ./...` sạch;
  `D:\Go\bin\gofmt.exe -l internal/ cmd/` không in gì (trừ nhiễu CRLF §9);
  `D:\Go\bin\go.exe build ./...` sạch.
- [ ] **Step 3: Kiểm thử đột biến tổng** — gỡ lưới ở actor khu vực (Task 2) → chạy test khu
  vực thấy **đỏ**; hoàn nguyên. Ghi lại kết quả đột biến từng task vào commit/ghi chú.
- [ ] **Step 4: Ghi sổ** — `BAT-DAU-TU-DAY.md`: §4 thêm dòng "Lưới chắn panic (nợ hạ tầng A)";
  §8 gỡ/đánh dấu mục liên quan panic nếu có; §12 nhật ký + dòng "Cập nhật lần cuối".
  `NHAT-KY-AGENT.md`: khối "ĐANG CHỜ" + một dòng cuối bảng.
- [ ] **Step 5: Commit** — `docs: ghi so xong luoi chan panic (no ha tang manh A)`

---

## Tự soi (đối chiếu spec)

- **Phủ spec:** §3.1 helper → Task 1. §3.2(a) actor chạy tiếp → Task 2; van an toàn → Task 3.
  §3.2(b) mỗi kết nối → Task 4. §3.2(c) hệ lưu → Task 5. §4 TDD+đột biến → mỗi task có bước
  đột biến; §7 tiêu chí xong → Task 6. ✅ đủ.
- **Ranh giới:** §6 (tách writer khỏi actor, Accept nghỉ giãn, tự bật lại) — **không** có
  task nào đụng, đúng chủ ý.
- **Nhất quán tên:** `ChayCoLuoi(cho string, f func()) bool`, `soPanicLienTiep`,
  `nguongPanicDungKhuVuc=100`, package `antoan` — dùng thống nhất mọi task.
