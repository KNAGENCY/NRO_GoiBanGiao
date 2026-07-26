# KH6 — 5 việc phải sửa trước khi gộp nhánh `kh6-nghien-cuu-chien-dau`

**Bối cảnh:** commit `3f7d23a` (nghiên cứu + đặc tả + kế hoạch 15 task cho Kế hoạch 6)
đã được soi ngày 26/07/2026. **Nội dung ĐÚNG** — người soi đối chiếu từng bảng byte,
12 con số hành vi và cả 4 "lỗi bản gốc" với code Java gốc, khớp hết; `go test ./...`
toàn xanh. Không phải làm lại gì. Nhưng còn 5 việc dưới đây, **chưa xong thì chưa gộp**.

Nhánh nằm ở worktree `D:\NRO_GoiBanGiao\.claude\worktrees\vao-viec-b4ec16`.

---

## 1. Sót một gói — quan trọng nhất, ảnh hưởng Task 8

Khi quái hồi sinh, bản Java gửi `-13` **rồi gửi tiếp gói `87`**:

- `mob/Mob.java:1480-1490` — `sendMobMaxHp()` = `Message(87)`, thân: `byte mobId` +
  `int maxHp` (qua `writeLongByGiaCatLuong` nên 4 byte).
- Gọi ở `mob/Mob.java:525`, ngay sau `sendMobHoiSinh()`.
- Gọi lần nữa ở `mob/Mob.java:357` (nhánh quái đã đầy máu trong chu kỳ 30 giây).

Đặc tả §3.3 và chuỗi gói §6.2 không nhắc gói này, mà cũng **không** xếp nó vào §0 "ngoài
phạm vi" — khác với cách đặc tả xử lý mọi thứ bị hoãn.

**Việc:** đọc client để quyết. Nếu client cần `87` → thêm vào §3.3 + §6.2 + Task 8. Nếu
kết luận không cần (vì `-24` đã mang `hpMax` và hồi sinh không đổi `hpMax`) → ghi rõ vào
§0 **kèm lý do**. Không để im lặng.

## 2. Kế hoạch, Task 10, dòng 332 — số byte sai trong commit message mẫu

Ghi `(don danh 14B / hoi mau 10B)`. Sai. Đặc tả §2.1 và Interfaces của chính Task 10 ghi
**18B / 11B** — đã cộng tay xác nhận: `4+1+1+4+1+1+4+1+1 = 18` và `4+1+1+4+1 = 11`.

## 3. Kế hoạch, Task 5, dòng 201 — chỉ sai file

Ghi "Modify: `banDo.go` hoặc nơi dựng `-24`". Gói `-24` dựng ở
**`server/internal/game/thegioi.go:105`**, hàm `goiMapInfo`. Ghi đúng file, bỏ chữ "hoặc".

## 4. Sửa lại diễn giải về cổng PostgreSQL — kết luận hiện tại sai

§9 của `BAT-DAU-TU-DAY.md` (bản trong nhánh) viết: *"§3 của file này từng ghi `5433` —
**sai**, đã sửa"*. Bằng chứng ngược lại:

- `go test ./internal/store/` với DSN cổng `5433` → **đỏ hàng loạt, không phải SKIP**.
- Mọi test CSDL của KH1–KH5 đều từng xanh.

⇒ máy **thật sự từng nghe `5433`**; cổng đã **đổi** qua một lần khởi động lại. Tài liệu
không sai từ đầu. Vì `postgresql.conf` để mặc định `5432` nhưng `pg_ctl -o "-p 5433"` thì
nghe `5433`, **ghim con số mới chỉ dời quả mìn chứ không tháo nó**.

**Việc:** sửa §3 và §9 sang hướng **dạy cách tự kiểm**:

```powershell
Get-NetTCPConnection -State Listen | Where-Object { $_.LocalPort -in 5432,5433 }
```

Giữ lại phần triệu chứng nhận dạng đã viết (`auth`/`net`/`protocol` xanh mà mọi test CSDL
đỏ hàng loạt) — phần đó rất tốt. Mẫu đã làm sẵn: `server/README.md` mục 0, commit `4f61dea`.

## 5. Nhánh rẽ từ `db34f6c`, thiếu `280a774` — gộp thẳng sẽ **xoá mất lệnh `/viec`**

Trong lúc nhánh này đang làm, `main` đã gộp `/vao-viec` + `/ghi-so` thành **một lệnh
`/viec`** tự nhận biết chế độ. Diff của nhánh đang **lùi** §10b và §11 của
`BAT-DAU-TU-DAY.md` về hai lệnh cũ.

**Việc:** `git merge main` vào nhánh này rồi giải xung đột:

| Chỗ | Giải thế nào |
|---|---|
| `BAT-DAU-TU-DAY.md` §10b và §11 | **Giữ bản của `main`** (`/viec`), bỏ `/vao-viec` + `/ghi-so` |
| §12 nhật ký và `NHAT-KY-AGENT.md` | **Giữ cả hai bên**, đừng bỏ dòng nào |
| §6 | Giữ mục `-60` hai dạng mà nhánh thêm, đánh số lại cho liền mạch |

---

## Xong thì

`go test ./... -count=1` toàn xanh (cổng đúng — tự kiểm theo việc 4), `gofmt` sạch,
`go vet` sạch; cập nhật `NHAT-KY-AGENT.md`; commit. **Đừng tự gộp lên `main`** — báo lại
để soi lần hai.
