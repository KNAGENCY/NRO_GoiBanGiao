# Lưới chắn panic — thiết kế (mảnh A của "vá nợ hạ tầng")

**Ngày:** 2026-08-03 · **Trạng thái:** đã duyệt thiết kế, chờ chủ dự án soi spec

> **Cho người thi công:** đọc [`BAT-DAU-TU-DAY.md`](../../../BAT-DAU-TU-DAY.md) §6 (lỗi
> đã vấp) + §9 (bẫy môi trường) trước. Kế hoạch thi công sẽ nằm ở
> `docs/superpowers/plans/` cạnh spec này.

---

## 1. Mục đích

Đặt **lưới đỡ (recover)** ở mọi luồng chạy song song (goroutine) của máy chủ, để **một
lỗi lạ (panic) ở bất kỳ đâu không làm sập cả tiến trình, đá văng TOÀN BỘ người chơi**.
Thay vì sập, lỗi được **ghi log rõ + cô lập trong phạm vi hẹp nhất**.

Đây là mảnh **A** trong bốn mảnh "vá nợ hạ tầng" (A lưới panic · B CI · C migration DB ·
D backup DB — xem [`BO-SUNG-LO-TRINH.md`](../plans/BO-SUNG-LO-TRINH.md) Khối 2 và Khối 8).
Nó là 🔴 vì đang hở thật và hậu quả là **cả server chết**.

## 2. Vì sao gấp — bằng chứng đã kiểm (2026-08-03)

Máy chủ có **6 điểm sinh goroutine**, và **KHÔNG điểm nào có `recover()`**
(`grep recover() internal/` chỉ ra 3 kết quả — cả 3 đều là **lời chú thích cảnh báo**,
không phải lưới thật):

| # | Goroutine | Vị trí | Panic ở đây → |
|---|---|---|---|
| 1 | **Actor mỗi khu vực** `vongLap` | `internal/game/khuvuc.go:473` | Bug lúc đánh nhau / xử sự kiện → **sập cả server** |
| 2 | **Nhận gói mỗi kết nối** `PhucVu` | `internal/net/listener.go:26` | Gói tin hỏng → **sập cả server** |
| 3 | **Gửi gói mỗi kết nối** `vongWriter` | `internal/net/session.go:117` | Lỗi lúc gửi → **sập cả server** |
| 4 | **Lưu định kỳ** `vongLap` | `internal/game/luudinhky.go:125` | Lỗi trong nhịp lưu → **sập cả server** |
| 5 | **Thợ ghi đĩa** `vongLap` | `internal/game/thogi.go:104` | DB/ổ đĩa trục trặc → **sập cả server** |
| 6 | **Mỗi lần lưu** `go func(nv)` | `internal/game/thogi.go:238` | Một bản ghi hỏng → **sập cả server** |

Chính code tự cảnh báo — `internal/protocol/frame.go:178`:
> *"không được recover() sẽ làm sập cả tiến trình máy chủ, đá văng TOÀN BỘ người chơi"*

…và `internal/game/login.go:156,743` ghi chú `net/session.go` **không có recover()**.
Nghĩa là lưới **chưa từng được đặt** — mới chỉ có người ghi chú là nó thiếu.

**Tầng mạng lõi (`network.*`, `utils.Logger`) nằm trong `lib/ga.jar` không có source** —
nhưng đó là bản Java cũ. Bản Go này tự viết toàn bộ 6 goroutine trên, nên **vá được trọn**.

## 3. Thiết kế

### 3.1 Một "lưới đỡ" dùng chung

Một hàm trợ giúp duy nhất, đặt ở package mới `internal/antoan` (an toàn):

```go
// ChayCoLuoi chay f; neu f panic thi CHAN lai, ghi log kem dau vet (stack),
// roi tra ve binh thuong — KHONG cho panic lan len lam sap tien trinh.
// Tra ve true neu da chan mot panic (de caller dem/quyet dinh), false neu chay tron.
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

Dùng chung một hàm — **không copy-paste 6 lần `recover()`**. Ai đọc log thấy trường `cho`
là biết panic ở đâu (`khuvuc.xuLy`, `phucvu:127.0.0.1`, `thogi.ghi`…).

Lý do đặt ở package riêng `internal/antoan`: cả `net` lẫn `game` đều dùng, không được
để `net` phụ thuộc `game` (vòng import). Package này chỉ phụ thuộc stdlib (`log/slog`,
`runtime/debug`).

### 3.2 Ba chính sách phục hồi — theo từng loại luồng

Lưới thì dùng chung, nhưng **caller quyết định làm gì SAU khi chặn** — khác nhau theo loại:

**(a) Actor khu vực (#1) — CÔ LẬP TỚI TỪNG SỰ KIỆN, khu vực CHẠY TIẾP.**
Đặt lưới **bên trong vòng lặp**, bọc từng lần xử sự kiện / nhịp tim — **không** bọc cả
`vongLap` (bọc ngoài thì panic vẫn thoát vòng lặp, làm **chết cả khu vực**). Cụ thể ở
`khuvuc.go:vongLap`:

```go
case sk := <-k.hopThu:
    antoan.ChayCoLuoi("khuvuc.xuLy", func() { k.xuLy(sk) })
case <-ticker.C:
    antoan.ChayCoLuoi("khuvuc.nhipTim", k.nhipTim)
```

→ Một cú đánh lỗi chỉ **bỏ đúng cú đó**; nhịp sau khu vực chạy bình thường, người chơi
khác không hề hấn. **Đây là quyết định chủ dự án đã chốt: khu vực gặp lỗi cứ chạy tiếp.**

⚠️ **Cạm bẫy phải chặn: panic dai dẳng mỗi nhịp.** Nếu một bug làm `nhipTim` panic **mọi
nhịp** (10 lần/giây), khu vực sẽ quay nóng CPU + ngập log. Vì vậy thêm **một van an toàn
nhẹ**: đếm số nhịp panic **liên tiếp** của khu vực; nếu vượt ngưỡng (đề xuất **100 nhịp
= 10 giây panic liền không dứt**) → ghi log **mức CRITICAL** ("khu vực X hỏng dai dẳng")
rồi **dừng riêng khu vực đó** (fail-safe), các khu vực khác + toàn server vẫn sống.
Bộ đếm **reset về 0 ngay khi có một nhịp chạy trót lọt** — nên lỗi *tạm thời* (đúng ý
"chạy tiếp") không bao giờ chạm ngưỡng; chỉ lỗi *dai dẳng thật* mới dừng. Van này giữ
đúng tinh thần "cứ chạy tiếp" cho lỗi thoáng qua, đồng thời không để một khu vực hỏng
vĩnh viễn ghì chết cả máy.

**(b) Mỗi kết nối (#2 `PhucVu`, #3 `vongWriter`) — ĐÓNG ĐÚNG KẾT NỐI ĐÓ.**
Bọc **toàn bộ thân goroutine**. Panic ở đâu trong lúc phục vụ một client → chặn + log +
**đóng sạch đúng kết nối đó** (giải phóng phiên, để hệ lưu-trạng-thái của KH11 chạy bình
thường), người chơi khác không sao. Với `PhucVu`, việc dọn kết nối phải đi qua đúng
đường đóng phiên đã có (không bỏ ngang giữa chừng làm rò tài nguyên).

**(c) Vòng lưu + thợ ghi (#4, #5, #6) — GIỮ VÒNG SỐNG, bỏ đúng việc lỗi.**
- #4 `luudinhky.vongLap`, #5 `thogi.vongLap`: bọc **thân mỗi vòng lặp** (không bọc ngoài,
  kẻo panic giết luôn vòng → mất hẳn tính năng lưu). Một nhịp lưu lỗi → bỏ nhịp đó, vòng
  vẫn quay nhịp sau.
- #6 `thogi.go:238 go func(nv)`: bọc cả thân → một bản ghi người chơi hỏng chỉ **bỏ đúng
  người đó** (đã có log), thợ ghi vẫn lưu người tiếp theo.

## 4. Cách kiểm thử — bắt buộc TDD + kiểm thử đột biến

Theo §10 của `BAT-DAU-TU-DAY.md` (dự án **đã 3 lần** test xanh mà code sai):

1. **Test hàm lưới `ChayCoLuoi`** (thuần, nhanh): f panic → trả `true`, **không** để panic
   lan ra (test tự nó không chết); f chạy trót → trả `false`. Có ghi log (bắt qua handler
   slog giả).
2. **Test từng chính sách** — dựng panic có chủ đích rồi khẳng định **cô lập đúng phạm vi**:
   - Actor: nhét một sự kiện làm `xuLy` panic → **khu vực vẫn nhận và xử sự kiện KẾ TIẾP**
     (chứng minh vòng lặp sống), người khác trong khu vực vẫn nhận gói.
   - Van an toàn: ép panic **liên tiếp** quá ngưỡng → khu vực dừng, **các khu vực khác vẫn
     chạy**; và ép panic **thoáng qua** (1–2 lần rồi thôi) → khu vực **không** dừng.
   - Kết nối: một client làm `PhucVu` panic → **đúng kết nối đó đóng**, client khác vẫn
     gửi/nhận bình thường.
   - Thợ ghi: một bản ghi panic → người **kế tiếp vẫn được lưu** (đọc lại DB thấy).
3. **Kiểm thử đột biến (bắt buộc):** gỡ lưới ở từng chỗ → test tương ứng **phải đỏ** (và
   đỏ vì *panic lan ra*, không phải vì lỗi biên dịch). Không đỏ = lưới vô dụng, viết lại test.

**Cách dựng panic để test** (không làm bẩn code thật): dùng **điểm chèn test** — ví dụ một
`XuLyGoiTin` giả panic khi gặp mã lệnh "ma thuật", hoặc một hàm lưu giả panic. Không thêm
nhánh `if test` vào code sản xuất.

## 5. Ràng buộc chung (đưa nguyên vào Global Constraints của kế hoạch)

- **Máy này KHÔNG chạy được `go test -race`** (thiếu trình biên dịch C) → bù bằng
  `-count=5` + soi kỷ luật goroutine. Test panic **không** cần `-race`.
- **`slog`** là logger chuẩn của dự án (đã dùng ở `listener.go`, KH11). Dùng đúng nó.
- **Không nuốt lỗi im lặng** — trái với bẫy §6 bản Java. Mọi panic bị chặn **phải để lại
  log kèm stack**; van an toàn khu vực ghi mức CRITICAL.
- **Không đổi hành vi lúc-không-có-lỗi.** Lưới chỉ hoạt động khi có panic; đường chạy
  bình thường phải y hệt trước (test cũ vẫn xanh).
- **Go:** `D:\Go\bin\go.exe`, gofmt `D:\Go\bin\gofmt.exe`. Cổng PostgreSQL **tự dò**
  (§9 BAT-DAU / `server/README.md` mục 0). Chú thích tiếng Việt giải thích *tại sao*.

## 6. Cố ý ĐỂ NGOÀI phạm vi

- **Tách gửi gói mạng khỏi actor (§6.10b BAT-DAU).** Actor hiện **tự gửi gói mạng trong
  vòng lặp** → có thể **kẹt tới 30 giây**. Đó là bug **KẸT (blocking)**, khác hẳn **SẬP
  (panic)**, và là thay đổi kiến trúc lớn hơn (đúng thứ nhánh `kh-actor` cũ định làm).
  **Không nhét vào mảnh A** — để thành mảnh riêng, rà lại nhánh `kh-actor` lúc đó.
- **Vòng `Accept()` không nghỉ giãn** (§8 BAT-DAU) — lỗi khác, không phải panic.
- **Tự bật lại khi sập** (Khối 2 điểm 2, dùng systemd/OS) — thuộc KH V, không phải code Go.

## 7. Tiêu chí "xong"

- Cả 6 goroutine đều qua lưới; `grep "go func\|go [a-z].*(" internal/` không còn điểm nào
  không được bọc.
- Bộ test panic phủ đủ 3 chính sách + van an toàn, **đột biến bắt được** (gỡ lưới → đỏ).
- `go test ./... -count=5` xanh, `go vet` + `gofmt` sạch, `go build` sạch.
- Cập nhật `BAT-DAU-TU-DAY.md` (§4 xong, §8 gỡ mục "chưa có lưới panic" nếu có) + sổ nhật ký.

## 8. Phụ lục — tham chiếu mã nguồn đã kiểm

| Việc | Vị trí |
|---|---|
| Actor khu vực `vongLap` (for-select) | `internal/game/khuvuc.go:479-505`, spawn ở `:473` |
| Nhận gói mỗi kết nối `PhucVu` | spawn ở `internal/net/listener.go:26` |
| Gửi gói `vongWriter` | spawn ở `internal/net/session.go:117`, thân `:142` |
| Lưu định kỳ | `internal/game/luudinhky.go:125` |
| Thợ ghi + mỗi-lần-lưu | `internal/game/thogi.go:104`, `:238` |
| Cảnh báo "panic sập cả tiến trình" | `internal/protocol/frame.go:178` |
| Chú thích "session.go không recover" | `internal/game/login.go:156,743` |
