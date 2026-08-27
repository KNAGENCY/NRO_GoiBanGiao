# Hàng đợi việc CODE — soi toàn bộ ngày 26/08/2026

Chủ dự án yêu cầu: *"soi toàn bộ xem cái gì cần code thì đẻ subagent ra code trước hết đi, phần hình ảnh anh với em làm sau"*.

File này là kết quả soi. **Chỉ liệt kê việc CODE** — việc chờ ảnh vẽ tay để riêng ở cuối.

---

## Luật chống đụng nhau — đọc trước khi giao bất cứ việc gì

`server/internal/game/khuvuc.go` là **63KB / ~1500 dòng** và gần như mọi việc gameplay đều chạm vào nó. Kế hoạch KH6b đã ghi thẳng:

> *"từ task 4 trở đi làm TUẦN TỰ: task 5-12 đều chạm `khuvuc.go` nên hai agent cùng vào là đụng chắc."*

Nên **không bao giờ** giao song song hai việc cùng chạm `internal/game/`. Chia theo vùng file, không chia theo "nghe có vẻ khác nhau".

| Vùng | Ai được vào |
|---|---|
| `server/internal/game/**` | **một agent duy nhất tại một thời điểm** |
| `assetpack/**` + `docs/my-thuat/mapping-id.json` | tách biệt, an toàn chạy song song |
| `tools/**` | tách biệt, an toàn chạy song song |
| `web/**` | tách biệt, an toàn chạy song song |
| `D:\NRO_Project` (client Unity) | **một agent duy nhất** — không nằm trong git này nên không có lưới an toàn |

---

## Đợt 1 — đang chạy (giao 26/08)

| # | Việc | Vùng file | Vì sao đáng làm |
|---|---|---|---|
| 1 | **Waypoint + bg-item trong MAP INFO** | `internal/game/`, client Unity | `thegioi.go:146` đang gửi `waypointCount = 0` → **người chơi không có cách nào rời map 1**. Và `w.Short(0)` bg-item → **8 ảnh bóng đổ đã vẽ xong nằm không dùng được**. Client đã có sẵn `BgItem.cs`, `BgItemMn.cs`, `Waypoint` — chỉ thiếu phía server |
| 2 | **Nợ assetpack** | `assetpack/`, `mapping-id.json` | `test_task6b_canonical_preflight.py` **chưa bao giờ chạy** (hỏng lúc thu thập) · `mapping-id.json` ghi map **5×17** trong khi map thật **52×24** → chạy `assetpack build` sẽ **đè artifact va chạm đúng bằng bản sai** |
| 3 | **Thước đo chỉ ra CHỖ hỏng** | `tools/map-lang-la/` | `do-map.py` nói "đất phẳng 91%" nhưng không nói **ở đâu** → vẫn phải sửa mò. Kèm cho `ve-bo-cuc-map.py` nghỉ hưu an toàn (chạy nhầm là thổi bay bố cục) |

---

## Đợt 2 — xếp hàng, chạy khi đợt 1 xong

Cả ba đều chạm `internal/game/` nên **phải tuần tự**, theo đúng thứ tự này.

### 2.1 — Bảo mật §7 (ưu tiên cao nhất)

**Đã xác minh còn nguyên, không tin sổ:** `grep` toàn bộ `login.go` + `internal/auth/` không có bất cứ thứ gì chống dò.

| Lỗ hổng | Bằng chứng |
|---|---|
| **Không chống dò mật khẩu** | Đo thật: **200 lần sai liên tiếp không bị chặn, 76 lần/giây** |
| **Rò rỉ thời gian ~13,5 ms** | Nhìn thấy trong code: không có tài khoản → `return` **ngay**; sai mật khẩu → phải chạy qua `auth.Verify` (băm chậm). Kẻ xấu **dò được tên tài khoản nào có thật** |
| **Không giới hạn kết nối trên một IP** | — |

Cách vá rò rỉ thời gian là chuẩn mực: khi không tìm thấy tài khoản, **vẫn chạy `auth.Verify` với một băm giả cố định** để hai nhánh tốn thời gian bằng nhau.

Chạm: `internal/game/login.go`, gói mới cho bộ đếm/khoá tạm, có thể cả `internal/net` cho giới hạn IP.

### 2.2 — KH6b Boss, task 6 → 15

Kế hoạch đã viết kỹ: [`2026-08-22-kh6b-boss-ke-hoach.md`](2026-08-22-kh6b-boss-ke-hoach.md). Task 1–5 xong (`db287a9`).

**Bắt buộc tuần tự**, mỗi task một agent, TDD. Hai chốt chặn kế hoạch đã ghi rõ, đừng phá:
- Gọi `SoKhuVuc.PhatToanServer` **từ trong actor là tự khoá chắc chắn** → phải qua bộ phát thanh riêng (`loathongbao.go`) theo khuôn `ThoGhi`
- `MoSoKhuVuc()` **cấm tuyệt đối** nhét khởi động lịch boss vào — sẽ làm đổ test tích hợp

### 2.3 — Nợ vặt §8 (gom một đợt)

Nhỏ nhưng thật, gom lại làm một lần cho đỡ tốn vòng:

- `byte(nv.Cap)` **cắt cụt** nếu cấp > 255
- `context.Background()` ở `login.go` — truy vấn CSDL **không hạn giờ, không huỷ** khi client ngắt
- Vòng `Accept()` **không nghỉ giãn** → cạn file descriptor sẽ quay nóng CPU
- `Dung()` gọi trước `BatDau()` sẽ **panic**
- Mã chết: `Writer.Bool`, `Writer.Long`, `TaiKhoan.LaAdmin`
- Client gửi `-39` trước `-28/13` → không làm gì (đúng, nhưng **chưa có test khoá**)
- Hai test chập chờn dưới tải: cho dùng khu vực **không quái mặc định**

---

## Đợt 3 — cần đặc tả TRƯỚC khi viết code

Đừng giao thẳng, sẽ phải làm lại.

| Việc | Chặn ở đâu |
|---|---|
| **KH7 — vật phẩm, hành trang, cửa hàng, giao dịch** | 🧱 **phải chốt căn cước vật phẩm + chống nhân bản trong đặc tả TRƯỚC** ([BO-SUNG](BO-SUNG-LO-TRINH.md) Khối 3). Tách rõ 4 thứ: cửa hàng NPC / giao dịch 2 người / ký gửi / kho |
| **Hồi máu PvP** | `goiHoiMauNguoiChoi` đã dựng + test byte nhưng **cố ý chưa nối dây**: chưa đặc tả lượng hồi bao nhiêu, người được hồi biết máu mới qua gói nào |
| **KH8 — kho nội dung YAML + nạp nóng** | đơn vị `SatThuong` phải chốt (TODO trong `satthuong.go`) |
| **W3 — web:** thu hồi phiên JWT | JWT hiệu lực 7 ngày, **không thu hồi** khi tài khoản bị khoá/hạ quyền |

---

## Đang chờ ẢNH — cố ý để lại, không giao đợt này

Chủ dự án đã nói rõ *"phần hình ảnh anh với em làm sau"*.

| Việc | Chờ gì |
|---|---|
| Map 1: chỉ số **đất phẳng lì 91%** (đích ≤40%) | 4 tile đất có vân đọc được ở cỡ thật. **Đã thử làm bằng code 3 lần, hỏng cả 3** — máy không tạo được nội dung mới |
| Map 1: 2 ảnh nền đúng cỡ **2976** | hiện 1536 nên núi lặp ~1,9 lần trong một khung hình |
| Map 2–5 (Cát, Sương Mù, Đá, Mây) | concept + một đợt ảnh mỗi làng, theo [`KHUON-DE-BAI-ANH-MAP.md`](../../my-thuat/KHUON-DE-BAI-ANH-MAP.md) |
| `part.bin` Naruto thật, atlas `Big0..4.png` | 343 ảnh đã có ở `docs/my-thuat/anh-game/` nhưng thiếu mapping ID → tên file |

---

## Bẫy nghiệm thu — dán vào mọi brief

1. **PostgreSQL tắt thì test tích hợp TỰ BỎ QUA trong im lặng** và bảng tổng kết vẫn `[OK]`. Ngày 26/08 bật lên thì lộ **15 test đỏ** nằm im nhiều ngày. Bật rồi mới được kết luận là xanh:
   ```
   D:\PostgreSQL\bin\pg_ctl.exe -D D:\PostgreSQL\data start
   export NRO_TEST_DSN="postgres://postgres:nrodev123@127.0.0.1:5432/nro_test"
   ```
2. **Đổi tài nguyên client mà quên tăng `vsData`** → client dùng bản cache cũ, **trong game không đổi gì cả**. Đã mắc hai lần.
3. **Sai byte không làm hỏng cả dòng, chỉ hỏng IM LẶNG** — lệch sang trường khác.
4. **`asset-lock.json` là JSON một dòng 444KB** — vá bằng regex đúng một mục, `json.dump` sẽ format lại và làm vỡ test bám chuỗi nguyên văn.
5. **`pytest` dùng `.pyc` cũ** nếu sửa cùng file hai lần trong một giây → `find assetpack -name __pycache__ -type d -exec rm -rf {} +`.
