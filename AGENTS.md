# AGENTS.md — hướng dẫn cho mọi agent AI làm dự án này

File này dành cho các agent **không phải Claude Code** (Codex, Cursor, Copilot, Gemini
CLI…). Claude Code đọc `CLAUDE.md`; nội dung bắt buộc thì giống nhau.

## Đọc 3 file này trước, theo đúng thứ tự

1. **`BAT-DAU-TU-DAY.md`** — dự án là gì, đã xong tới đâu, còn gì chưa làm, **§6 danh sách
   lỗi đã vấp** (đọc kỹ, đừng lặp lại), **§9 bẫy môi trường của máy này**.
2. **`NHAT-KY-AGENT.md`** — ai đang làm gì, có ai đang mở cùng vùng code với bạn không.
3. **`CLAUDE.md`** — kiến trúc, đường ống dữ liệu, bẫy của phần code cũ.

Rồi chạy `git log --oneline -15` và `git status --short --branch` để biết trạng thái thật.

## Bốn quy tắc không được vi phạm

1. **Trả lời bằng tiếng Việt đời thường.** Chủ dự án không phải lập trình viên — nói *tại sao*
   trước khi nói *làm thế nào*.
2. **Làm trên nhánh git riêng**, không commit trực tiếp lên `main`:
   `git switch -c <ten-viec>`. Có nhiều agent chạy song song.
3. **TDD + kiểm thử đột biến.** Test đỏ trước, code sau. Xong việc thì cố tình phá code xem
   test có bắt — không bắt nghĩa là test vô dụng, phải viết lại test. Dự án này đã 3 lần bị
   "test xanh mà code sai".
4. **Xong việc phải ghi sổ**: thêm một dòng vào cuối bảng `NHAT-KY-AGENT.md` và cập nhật
   `BAT-DAU-TU-DAY.md` theo §11 của file đó, commit chung với code. Không ghi = phiên sau
   và agent khác mất phương hướng.

Quy trình đầy đủ của dự án (nhận việc + ghi sổ) viết ở `.claude/commands/viec.md` — đọc file
đó và làm theo. Claude Code gọi nó bằng lệnh `/viec`; agent khác thì đọc trực tiếp.

## Chia vùng để không đụng nhau

| Vùng | Đường dẫn | Ghi chú |
|---|---|---|
| Máy chủ game (Go) | `server/` | Vùng chính, đang phát triển nhanh |
| Web + panel (Next.js) | `web/` | Tách biệt hẳn với `server/` — an toàn làm song song |
| Tài liệu | `docs/` | |
| Code Java/PHP cũ | `1-MAY-CHU-GAME/`, `2-WEB/`, `3-CO-SO-DU-LIEU/`, `5-TAI-LIEU/` | **Chỉ đọc để tra cứu giao thức.** Không sửa |

Hai agent làm `server/` và `web/` song song thì gần như không bao giờ đụng nhau. Hai agent
cùng làm `server/internal/game/` thì **chắc chắn** sẽ đụng — đừng làm vậy.

## Lệnh kiểm chứng (phải xanh trước khi báo xong)

```powershell
# Máy chủ game
cd server
$env:NRO_TEST_DSN = "postgres://postgres:nrodev123@127.0.0.1:5433/nro_test"
D:\Go\bin\go.exe test ./... -count=1
D:\Go\bin\go.exe vet ./...

# Web
cd web
npm run build
```

Trên máy này: **PowerShell không có `&&`** (dùng `;`) · dùng `python` chứ không phải
`python3` · **không chạy được `go test -race`** (thiếu trình biên dịch C) → bù bằng
`-count=5` và soi kỹ kỷ luật actor.
