---
description: Nhận việc dự án Naruto — đọc bối cảnh, xem ai đang làm gì, đề xuất việc tiếp
argument-hint: "[việc muốn làm, để trống thì tự đề xuất]"
---

Bạn vừa được gọi vào dự án game Naruto này. Làm đúng các bước sau **trước khi
sửa bất cứ dòng code nào**:

## 1. Nạp bối cảnh (đọc, đừng bỏ bước nào)

- `BAT-DAU-TU-DAY.md` — dự án là gì, đã xong tới đâu, còn gì, **§6 các lỗi đã vấp**
  (đọc kỹ mục này, đừng lặp lại), §9 bẫy môi trường.
- `NHAT-KY-AGENT.md` — **ai đang làm gì**, có ai đang mở cùng vùng với việc bạn định làm không.
- `git log --oneline -15` và `git status --short --branch` — trạng thái thật của cây làm việc.
- `docs/superpowers/plans/LO-TRINH.md` — thứ tự các kế hoạch.
- Nếu việc liên quan giao thức nhị phân: `docs/superpowers/specs/` có bảng byte đã giải mã.
  **Đừng tự đoán byte** — tra spec, không có thì đọc code Java cũ trong `1-MAY-CHU-GAME/src/`.

## 2. Chống đụng nhau với agent khác

- Tạo nhánh riêng ngay: `git switch -c <tên-việc>`. **Không commit trực tiếp lên `main`.**
- Nếu `NHAT-KY-AGENT.md` cho thấy có agent khác đang sửa cùng vùng → nói rõ cho người
  dùng biết và đề nghị đổi việc, đừng làm chồng.

## 3. Báo cáo rồi mới làm

Trả lời **bằng tiếng Việt đời thường** (chủ dự án không phải lập trình viên), gồm:

1. Hiện trạng: đã xong gì, đang có ai làm gì.
2. Việc bạn sẽ làm và **tại sao** nó là bước hợp lý tiếp theo.
3. Cạm bẫy liên quan bạn đã đọc được ở §6/§9 và sẽ tránh thế nào.

Việc người dùng muốn làm: $ARGUMENTS
(Nếu để trống: đề xuất việc tiếp theo hợp lý nhất theo `LO-TRINH.md` và §5 của `BAT-DAU-TU-DAY.md`.)

## 4. Cách làm việc bắt buộc của dự án này

- **TDD**: test đỏ trước, code sau.
- **Kiểm thử đột biến**: xong việc, cố tình phá code xem test có bắt. Không bắt = test vô dụng,
  phải viết lại test. (Dự án này đã 3 lần bị test xanh trong khi code sai.)
- Chú thích code bằng **tiếng Việt, giải thích *tại sao*** chứ không phải *cái gì*.
- Trước khi báo xong: `go test ./... -count=1` toàn xanh (game) hoặc `npm run build` sạch (web).
- Xong việc → chạy `/ghi-so` để ghi lại. **Không ghi sổ = phiên sau mất phương hướng.**
