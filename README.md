# Game Naruto online (viết lại trên lõi Ngọc Rồng Online)

MMORPG 2D chủ đề Naruto, 5 làng ninja. Client Unity tái sử dụng từ Ngọc Rồng Online
(giữ nguyên giao thức mạng nhị phân); máy chủ viết lại bằng **Go + PostgreSQL** (`server/`),
web + panel quản trị bằng **Next.js** (`web/`).

> ⚠️ **Repo CÔNG KHAI** (chủ dự án quyết 22/10/2026). Ai cũng đọc được — **tuyệt đối không
> commit mật khẩu thật, khoá, dump CSDL, IP máy chủ đang chạy**. Mật khẩu trong script
> (`nrodev123`) chỉ dành cho CSDL chạy thử trên máy mình.

## Người mới vào: đọc theo thứ tự

1. [`BAT-DAU-TU-DAY.md`](BAT-DAU-TU-DAY.md) — **điểm vào duy nhất**: dự án là gì, xong tới đâu
   (§4), còn gì (§5), lỗi đã vấp (§6), cách chạy (§3), cách làm việc (§10).
2. [`NHAT-KY-AGENT.md`](NHAT-KY-AGENT.md) — ai vừa làm gì, dòng mới nhất ở cuối.
3. Agent AI không phải Claude Code: [`AGENTS.md`](AGENTS.md). Claude Code: [`CLAUDE.md`](CLAUDE.md).

## Những thứ KHÔNG có trong repo — xin chủ dự án gửi riêng

| Thứ | Ở máy chủ dự án | Cần khi |
|---|---|---|
| Client Unity (nhiều GB, không có git) | `D:\NRO_Project` | Chơi thử, sửa client, lắp ảnh |
| Ảnh nguồn mỹ thuật (~270 MB) | `docs/my-thuat/anh-nguon/` | Sinh lại sprite / làm ảnh mới |
| Kho asset server Java cũ (~820 MB) | `1-MAY-CHU-GAME/data/` | Chỉ khi tra cứu bản cũ |
| CSDL, file cấu hình có mật khẩu | (không chia sẻ) | Tự tạo CSDL trống theo §3 |

## Cách làm việc chung

- **Mỗi người / agent một nhánh riêng**, không commit thẳng lên `main`; xong thì mở Pull Request.
- `server/` và `web/` tách biệt — làm song song an toàn. Hai người cùng sửa
  `server/internal/game/` thì chắc chắn đụng nhau, báo trước.
- Làm xong việc đáng kể → cập nhật `BAT-DAU-TU-DAY.md` (§11) + một dòng `NHAT-KY-AGENT.md`.
- Kiểm toàn bộ trước khi đẩy: `powershell -File kiem-tra-tat-ca.ps1`. CI GitHub chạy lại mỗi lần push.
- Mật khẩu `nrodev123` trong script chỉ là CSDL **chạy thử trên máy mình** — không dùng cho máy chủ thật.
