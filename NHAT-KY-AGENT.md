# Nhật ký agent — ai đã làm gì

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
| 26/07/2026 | Claude Code (phiên chính) | Dựng cơ chế nhiều agent: sổ này + `/vao-viec` + `/ghi-so` + hook tự nạp bối cảnh | `main` | — |
