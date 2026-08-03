# KH6c — Stress test 20 client cùng khu vực (thiết kế)

**Ngày:** 03/08/2026 · **Trạng thái:** thiết kế đã duyệt (chờ viết plan) · **Nhánh:** `kh6c-stress-test`

## 1. Vì sao làm cái này

§6.10b ("actor tự gửi gói làm treo cả khu vực") **đã vá** (commit `a135979`) và có test 2 tầng
chứng minh *cơ chế*: `XepGui` không chặn socket, actor gạch client chậm mà không treo. Nhưng
các test đó chỉ dựng **1–2 client**. **Chưa có bằng chứng dưới tải thật: 20 người cùng khu vực
cùng xả chiêu.** KH6c lấp đúng khoảng đó — chứng minh dưới tải, actor vẫn nhịp đều, không treo,
và các cơ chế bảo vệ (trần hộp thư 256, gạch client chậm) vẫn chạy **và có đo đếm**.

Đây **không** phải sửa lỗi mới. Là **bài kiểm tải + thêm quan sát (observability)**.

## 2. Bối cảnh code hiện tại (đã đọc, không đoán)

- Actor mỗi khu vực = `vongLap` (`khuvuc.go`), nhịp `ticker` **100ms**, là **nơi duy nhất** chạm
  trạng thái khu vực.
- Hộp thư vào `hopThu` sức chứa **256**. `Gui(sk)` (`khuvuc.go:1420`): đầy thì **bỏ sự kiện +
  `soLanGuiRoi.Add(1)` (atomic) + log Warn**. → **Bộ đếm rơi ĐÃ CÓ SẴN.**
- Gửi ra: actor gọi `XepGui` (hàng đợi riêng mỗi phiên, 256, không chặn); client đầy hàng bị
  `loaiNhieuNguoiLoiGui` **gạch tên + `DongCham`**.
- **CHƯA CÓ:** đo **trễ nhịp** (khoảng cách giữa 2 lần `nhipTim`). Ticker của Go **bỏ tick** nếu
  actor bận không kịp nhận → cần đo mới biết actor có bị dồn dưới tải không.

## 3. Thiết kế — 3 mảnh

### 3.1. Mảnh A — Thêm đo trễ nhịp (code sản xuất, rất nhỏ)

Actor tự đo khoảng cách giữa 2 lần chạm nhịp. Vì `vongLap` là luồng **duy nhất** ghi, giữ mốc
thời gian bằng biến thường bên trong; phơi ra ngoài bằng **atomic** để test/giám sát đọc song
song an toàn.

- Trường mới trên `KhuVuc`:
  - `soNhipTre atomic.Int64` — đếm số nhịp có khoảng cách từ nhịp trước **> `nguongNhipTre`**.
  - `treNhipLauNhatMs atomic.Int64` — trễ lâu nhất (ms) từng ghi nhận (để báo cáo/giám sát).
- Hằng số `nguongNhipTre = 200 * time.Millisecond` (2× chu kỳ; kiểu **hằng** để test đột biến
  "đổi 200ms thành 2s" sửa được và bắt được).
- Đo ở đâu: trong `vongLap`, tại nhánh `case <-ticker.C`, so mốc `time.Now()` với mốc nhịp
  trước. **Nhịp đầu bỏ qua** (chưa có mốc trước).
- **Không đổi hành vi actor** — chỉ cộng đo. Không thêm khóa, không thêm goroutine.
- **Tách hàm thuần để test tất định:** `capNhatTreNhip(khoangCach time.Duration)` nhận khoảng
  cách, quyết định có tăng `soNhipTre` + cập nhật max — test gọi thẳng, **không phụ thuộc đồng
  hồ thật** (đúng tinh thần đã tách `ghiNhanNhip`).

### 3.2. Mảnh B — Bài stress test in-process (Go test)

File mới `server/internal/game/khuvuc_stress_test.go`, hàm `TestStress20ClientXaChieu...`:

1. Mở 1 khu vực (map có quái để đánh), `BatDau`.
2. Thêm **20 người chơi giả**. Phiên giả (`phienGiaKhuVuc`, mở rộng nếu cần): **đa số "khỏe"**
   (rút hàng đều) + **2–3 "chậm"** (không rút → hàng đầy → bị gạch). Trộn để test luôn cả đường
   gạch-client-chậm dưới tải.
3. **20 goroutine**, mỗi cái là 1 người, lặp gửi sự kiện broadcast vào actor qua `k.Gui(...)`,
   trong **N giây** (đề xuất 2s) với nhịp cao (cố tình làm đầy hộp thư 256). *Thi công dùng
   `suKienDiChuyen` (di chuyển) thay vì đánh quái: cùng đường tới hạn mailbox→`phatChoNguoiKhac`
   →`XepGui`→gạch client chậm, mà không cần `chieuDangChon` (nil sẽ nil-deref) — throughput
   ~5M sự kiện/2s vẫn là stress thật.*
4. **1 goroutine dò (prober)**: cứ ~5ms gọi `AnhChupNguoiChoiCoLoi(0)`, đếm số lần trả lời được
   (OK) và số lần quá hạn (timeout 100ms nội tại): actor **không treo** thì OK hàng trăm lần;
   treo kiểu §6.10b (kẹt 30s) thì OK ~0.
5. Kết thúc: dừng bơm, `Dung()` khu vực sạch (không rò goroutine).

### 3.3. Mảnh C — Khẳng định + báo số + đột biến

**Khẳng định (test đỏ nếu sai):**
- **Không treo (bằng chứng chính):** (a) actor trả lời **≥10 probe OK** trong 2s (treo §6.10b → ~0);
  (b) sau khi ngừng tải, actor hồi phục + **gạch đủ client chậm** — chỉ actor đang chạy mới gạch
  được, nên đây là bằng chứng "còn sống" mạnh nhất. *Timeout probe lẻ dưới tải cực nặng KHÔNG
  fail (§6) — chỉ báo số.*
- **Đóng socket đúng:** mỗi client chậm bị `DongCham` **đúng 1 lần**; client khỏe **không** bị đóng.
- **Trễ nhịp có kiểm soát:** `treNhipLauNhatMs` **< 500ms** (ngưỡng đậu; thực tế ~100ms).
- **Rơi hộp thư: chỉ BÁO SỐ, không hard-assert.** (`soLanGuiRoi > 0` phụ thuộc tốc độ máy →
  hard-assert sẽ chập chờn.) Đường đếm-rơi đã có test riêng `TestGuiHopThuDayDemVaLog`; KH6c chỉ
  đọc + in `soLanGuiRoi` để quan sát, không lấy làm điều kiện đậu/rớt.
- **Client chậm bị gạch:** sau test, số người còn trong khu vực = số "khỏe".

**Báo số (`t.Log`):** tổng sự kiện bơm, tổng rơi hộp thư (`soLanGuiRoi`), số nhịp trễ
(`soNhipTre`), trễ lâu nhất (`treNhipLauNhatMs`), max thời gian probe.

**Kiểm thử đột biến (bắt buộc — kỷ luật dự án §6.11) — đã chạy thật:**
- `capNhatTreNhip` no-op → `TestCapNhatTreNhip` **đỏ**. ✓
- `danhDauNhip` quên gán lại mốc → `TestDanhDauNhip` **đỏ** (khóa đường nối đo nhịp — reviewer bắt). ✓
- Hạ `nguongNhipTre` xuống 1ms (positive control) → `soNhipTre=19`, `treNhip=104ms` → chứng minh
  vòng lặp đo THẬT khoảng cách ~100ms (không đo mù). ✓
- Bỏ `loaiNhieuNguoiLoiGui` → client chậm không bị gạch → `conLai=20` → stress test **đỏ**. ✓

## 4. Ngưỡng đậu (tùy chỉnh trong test)

| Chỉ số | Ngưỡng |
|---|---|
| Probe OK trong 2s | **≥ 10** (actor treo §6.10b → ~0) |
| `DongCham` client chậm | **đúng 1 lần** mỗi người; client khỏe = 0 |
| Trễ nhịp lâu nhất | **< 500ms** (thực tế ~100ms) |
| Rơi hộp thư | chỉ **báo số** (không hard-assert, tránh chập chờn) |
| Người còn lại sau test | = số người "khỏe" |

Số client (20), thời lượng (2s), nhịp bơm: hằng số trong test, chỉnh được.

## 5. Ngoài phạm vi (YAGNI)

- **Test tải qua TCP thật** (mở rộng `thuclient` vào thế giới + xả chiêu): để dành. Đường TCP
  thật đã smoke bằng **Unity thật** 03/08 (login → vào map → render).
- **Không** đổi hành vi actor (trần 256, gạch client chậm giữ nguyên) — chỉ **thêm đo**.
- **Không** tối ưu hiệu năng — KH6c chỉ **đo + chứng minh**, chưa phải chỗ tăng tốc.

## 6. Rủi ro / bẫy (đọc trước khi code)

- **Test tải dễ chập chờn trên máy bận** (§8 đã có 2 test chập chờn tương tự). Giảm bằng: ngưỡng
  probe rộng (100ms, không phải 10ms), thời lượng vừa phải, **không khẳng định throughput tuyệt
  đối** (chỉ khẳng định "không treo + có đếm + trễ nhịp trong ngưỡng").
- `AnhChupNguoiChoiCoLoi` hiện chờ tối đa 100ms rồi trả rỗng (§8). Prober phải phân biệt "trả
  rỗng vì quá 100ms" (dấu hiệu bận, vẫn tính là **treo** nếu lặp lại) với "trả kịp có dữ liệu".
  Dùng cả giá trị trả về lẫn thời gian đo.
- Đo trễ nhịp bằng đồng hồ thật → **tách hàm thuần** (`capNhatTreNhip`) để test tất định.
- `time.Now()` trong actor là code sản xuất bình thường (khác script workflow) — dùng được.

## 7. Thứ tự thi công (sẽ chi tiết ở plan)

1. Mảnh A trước, **TDD**: viết test thuần `capNhatTreNhip` (đỏ) → cài đo trễ nhịp (xanh) → đột biến.
2. Mảnh B+C: viết stress test (đỏ vì chưa có metric/chưa đủ) → chạy → chỉnh ngưỡng → đột biến 3 kiểu.
3. `go test ./... -count=5` + `vet` + `gofmt` sạch, rồi soi đối kháng, rồi gộp `main`.
