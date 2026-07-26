---
description: Ghi lại việc vừa làm vào nhật ký + cập nhật tài liệu trạng thái, rồi commit
argument-hint: "[mô tả ngắn việc vừa làm, để trống thì tự suy từ git diff]"
---

Ghi lại việc vừa làm để phiên sau và agent khác hiểu được. Làm đúng thứ tự:

## 1. Xem lại mình vừa làm gì

Chạy `git status --short`, `git log --oneline -5`, và `git diff --stat` (cả staged và chưa)
để biết chính xác đã chạm vào file nào. Đừng ghi theo trí nhớ.

Việc người dùng mô tả: $ARGUMENTS

## 2. Kiểm chứng trước khi ghi "xong"

**Không được ghi là xong nếu chưa chạy và thấy kết quả.** Chạy đúng bộ kiểm phù hợp:

```powershell
# Máy chủ game
cd D:\NRO_GoiBanGiao\server
$env:NRO_TEST_DSN = "postgres://postgres:nrodev123@127.0.0.1:5433/nro_test"
D:\Go\bin\go.exe test ./... -count=1
D:\Go\bin\go.exe vet ./...

# Web
cd D:\NRO_GoiBanGiao\web
npm run build
```

Nếu test đỏ → ghi vào sổ là **đang dở, test đỏ ở đâu**, đừng ghi là xong.

## 3. Ghi `NHAT-KY-AGENT.md`

Thêm **một dòng vào CUỐI bảng** (mới nhất nằm dưới cùng):

`| <ngày dd/mm/yyyy> | <ai: Claude Code / Codex / …> | <làm gì, một câu> | <nhánh> / <commit> | <còn dở gì, hoặc —> |`

Cột "còn dở / cảnh báo" là cột quan trọng nhất — file đang sửa nửa vời, nhánh chưa gộp,
test đang đỏ, quyết định còn treo. Không có thì ghi `—`.

## 4. Cập nhật `BAT-DAU-TU-DAY.md`

Theo đúng §11 của chính file đó:
- **§4 Đã làm** — thêm dòng nếu vừa xong một mốc, kèm mã commit.
- **§5 Còn lại** — xoá việc vừa xong, thêm việc mới lộ ra.
- **§6 Lỗi sai** — nếu vấp lỗi mới: ghi lỗi + **vì sao vấp** + cách tránh. Đây là phần
  giá trị nhất của cả tài liệu, đừng bỏ qua vì "lỗi nhỏ".
- **§12 Nhật ký** + dòng "Cập nhật lần cuối" ở đầu file.

## 5. Commit

Commit code + tài liệu **cùng nhau** (để lịch sử khớp nhau), thông điệp không dấu:

```
<loai>(<vung>): <viec da lam>
```

`<loai>` = `feat` / `fix` / `test` / `docs` / `refactor`. `<vung>` = `game` / `web` / `docs`.

Nếu đang ở nhánh riêng và việc đã xong hẳn + test xanh → nói cho người dùng biết là
**gộp được rồi**, nhưng đừng tự gộp lên `main` nếu người dùng chưa nói, vì có thể agent
khác đang làm dở trên nhánh khác.

## 6. Báo lại bằng tiếng Việt

Ba câu: vừa làm gì · kiểm chứng bằng kết quả nào (con số thật, không nói chung chung) ·
việc tiếp theo nên là gì.
