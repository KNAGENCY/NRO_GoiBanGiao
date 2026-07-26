# Nhật ký agent — ai đã làm gì

## ⏳ ĐANG CHỜ GÌ NGAY LÚC NÀY — cập nhật: 26/07/2026

> Ai vào cũng đọc khối này trước. **Xong việc thì phải sửa lại khối này**, nếu không nó
> thành thông tin cũ gây hại hơn là không có.

- **ĐANG CHỜ SOI LẦN HAI:** nhánh `kh6-nghien-cuu-chien-dau` (worktree
  `.claude\worktrees\vao-viec-b4ec16`) đã **sửa xong cả 5 việc** ở
  [`docs/superpowers/plans/KH6-VIEC-CAN-SUA-TRUOC-KHI-GOP.md`](docs/superpowers/plans/KH6-VIEC-CAN-SUA-TRUOC-KHI-GOP.md)
  (khối ✅ ở đầu file đó ghi rõ làm gì), đã `git merge main` tới `6c8cbec`, test toàn xanh.
  **Nhánh chưa gộp — chờ phiên điều phối soi rồi gộp `--no-ff`.**
- **Việc kế tiếp của phiên điều phối:** soi lần hai → nếu sạch thì **gộp `--no-ff` vào
  `main`** rồi xoá nhánh.
- ~~Chuyển bài học PowerShell/mojibake vào §6~~ → **đã làm**, nay là **§6.15**.
- **Rồi mới tới:** thi công KH6 theo kế hoạch 15 task. **Task 12 đã xong** (đã đọc client
  Unity xác minh `-8`/`-17`) → **không còn task nào bị chặn**, ba nhánh chạy song song được.

---

Sổ chung cho **mọi agent và mọi người** làm dự án này. Mục đích: phiên sau (hoặc
agent khác đang làm song song) biết ngay ai vừa chạm vào cái gì, để không sửa đè
lên nhau.

**Cách ghi — bắt buộc:**
- Thêm dòng mới **vào CUỐI bảng** (mới nhất nằm dưới cùng — máy tự đọc 20 dòng cuối).
- Ghi **ngay khi xong việc**, cùng lúc với commit việc đó. Đừng dồn cuối ngày.
- Cột *Ai*: ghi rõ để phân biệt, ví dụ `Claude Code (cửa sổ 1)`, `Codex`, `Cursor`, `tôi (tay)`.
- Cột *Còn dở / cảnh báo*: quan trọng nhất. Ghi file đang sửa nửa vời, nhánh chưa gộp,
  test đang đỏ. Không có gì dở thì ghi `—`.

**Quy tắc chống đụng nhau (đọc trước khi bắt đầu):**
1. Mỗi agent làm **một nhánh git riêng**: `git switch -c viec-cua-toi-<tên>`. Không hai
   agent cùng commit lên `main`.
2. Xem bảng dưới trước khi sửa. Nếu agent khác đang mở cùng vùng (`server/internal/game/`,
   `web/app/`…) thì **đổi việc khác** hoặc chờ nó gộp xong.
3. Vùng ít đụng nhau nhất: **máy chủ game** (`server/`) và **web** (`web/`) là hai vùng
   tách biệt — hai agent làm hai vùng này song song rất an toàn.

| Ngày | Ai | Làm gì | Nhánh / commit | Còn dở / cảnh báo |
|---|---|---|---|---|
| 26/07/2026 | Claude Code (phiên chính) | Xong KH5 (di chuyển + actor), gộp về main | `main` / `e1ddc44` | — |
| 26/07/2026 | Claude Code (phiên chính) | Tạo `BAT-DAU-TU-DAY.md` làm điểm vào duy nhất | `main` / `24514d1` | — |
| 26/07/2026 | Claude Code (phiên chính) | Dựng cơ chế nhiều agent: sổ này + hook tự nạp bối cảnh + `AGENTS.md` + `MO-AGENT.bat` | `main` / `db34f6c` | — |
| 26/07/2026 | Claude Code (phiên chính) | Gộp `/vao-viec` + `/ghi-so` thành **một lệnh `/viec`** tự nhận biết chế độ (theo yêu cầu người dùng cho nhanh) | `main` / `280a774` | — |
| 26/07/2026 | Claude Code (worktree `vao-viec-b4ec16`) | **KH6 — thiết kế**: nghiên cứu giao thức chiến đấu + đặc tả bảng byte + kế hoạch thi công. Không sửa code Go | `kh6-nghien-cuu-chien-dau` / `3f7d23a` | Nhánh chưa gộp. Ai lấy KH6 làm tiếp thì **cả `server/internal/game/` là vùng của người đó** — đừng để agent khác vào cùng |
| 26/07/2026 | Claude Code (phiên chính) | **Soi commit `3f7d23a`** (KH6 thiết kế, của worktree `vao-viec-b4ec16`): đối chiếu từng bảng byte + 12 con số hành vi + 4 "lỗi bản gốc" với code Java → **khớp hết**; chạy lại `go test ./...` toàn xanh trên cổng 5432. Sửa `server/README.md` dạy cách **tự kiểm cổng** thay vì ghim số | `main` | ⚠️ **Nhánh `kh6-nghien-cuu-chien-dau` rẽ từ `db34f6c`, CHƯA có `280a774`** → gộp thẳng sẽ xung đột `BAT-DAU-TU-DAY.md` + sổ này và **có thể xoá mất lệnh `/viec`**. Phải `git merge main` vào nhánh đó trước, giải xung đột theo hướng GIỮ `/viec`. **5 việc phải sửa trước khi gộp — đã viết đầy đủ, kèm bằng chứng, ở [`docs/superpowers/plans/KH6-VIEC-CAN-SUA-TRUOC-KHI-GOP.md`](docs/superpowers/plans/KH6-VIEC-CAN-SUA-TRUOC-KHI-GOP.md)**: sót gói `87`, số byte 14B/10B sai, Task 5 chỉ sai file, §9 kết luận sai về cổng, và việc merge `main` để không mất `/viec` |
| 26/07/2026 | Claude Code (phiên chính) | **Bài học mới, chưa kịp đưa vào §6** (đợi gộp xong để không rối số mục): **KHÔNG dùng PowerShell `Get-Content -Raw` + `Set-Content` để sửa file tiếng Việt** — PS 5.1 đọc UTF-8 bằng codepage ANSI rồi ghi lại UTF-8 → mojibake sạch cả file, và diff phình ra như viết lại toàn bộ. Đã tự vấp khi sửa `server/README.md`, phải `git checkout --` rồi làm lại bằng công cụ Edit. Sửa file có dấu thì **luôn dùng Edit/Write**, không dùng shell | `main` | ✅ Đã chuyển thành **§6.15** khi gộp KH6 |
| 26/07/2026 | Claude Code (worktree `vao-viec-b4ec16`) | **Sửa xong 5 việc soi ra** + gộp `main` (tới `6c8cbec`). Đọc client Unity nên **Task 12 xong luôn** → nhóm E hết bị chặn. Lộ thêm: gói `87` là "gió thổi" chứ không phải maxHp (bản Go **không gửi**), `-9` trường cuối là **1 byte** không phải 4, `-13` hai byte giữa là `sys`/`levelBoss`, gói xuống `54` **không có trường đếm**. Tự sửa một chỗ mình viết sai: "sai trường → client hỏng vĩnh viễn" là **SAI**, chỉ hỏng gói đó và hỏng **im lặng** | `kh6-nghien-cuu-chien-dau` | **Chưa gộp — chờ soi lần hai.** Test 7/7 gói xanh, vet + gofmt sạch. Hai thứ chủ dự án cần chốt trước khi thi công: **hệ số ngũ hành** (đề xuất 1,25 / 1,00 / 0,80) và **20 chiêu (4 mỗi làng)** |
