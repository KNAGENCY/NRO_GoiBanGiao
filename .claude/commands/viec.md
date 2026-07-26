---
description: Lệnh duy nhất của dự án Naruto — tự biết là nhận việc hay ghi sổ việc vừa làm
argument-hint: "[để trống = tự đoán | mô tả việc muốn làm | \"ghi so\"]"
---

Lệnh duy nhất của dự án này. Tự nhận biết bạn đang ở **đầu việc** (cần nhận việc) hay
**cuối việc** (cần ghi sổ + commit). Làm đúng thứ tự dưới đây, đừng bỏ bước.

Người dùng nói thêm: $ARGUMENTS

---

# PHẦN A — Luôn làm trước, dù ở chế độ nào

Nạp bối cảnh. Đọc thật, đừng nhớ theo trí nhớ:

1. `git status --short --branch` và `git log --oneline -15` — trạng thái thật.
2. `NHAT-KY-AGENT.md` — **ai vừa làm gì**, có ai đang mở cùng vùng code không.
3. `BAT-DAU-TU-DAY.md` — §4 đã làm, §5 còn lại, **§6 các lỗi đã vấp** (đọc kỹ, đừng lặp
   lại), §9 bẫy môi trường của máy này.

---

# PHẦN B — Chọn chế độ

| Tình huống | Chế độ |
|---|---|
| Người dùng nói rõ `ghi so` / "ghi lại" / "xong rồi" | **GHI SỔ** |
| Người dùng mô tả một việc muốn làm | **NHẬN VIỆC** |
| Để trống + `git status` **sạch** | **NHẬN VIỆC** |
| Để trống + có thay đổi **do chính phiên này** vừa tạo ra | **GHI SỔ** |
| Để trống + cây bẩn nhưng **không phải phiên này sửa** | ⛔ **DỪNG** — xem dưới |

⛔ **Bắt buộc dừng lại hỏi** nếu cây làm việc có thay đổi mà phiên này không tạo ra
(agent khác hoặc người dùng đang sửa dở). Liệt kê file đang bẩn, nói rõ "cái này không
phải tôi sửa", hỏi người dùng muốn gì. **Tuyệt đối không commit công việc dở của người
khác** — đó là cách nhanh nhất để mất việc của họ.

Nói ra chế độ đã chọn và vì sao, rồi làm phần tương ứng.

---

# PHẦN C1 — Chế độ NHẬN VIỆC

## 1. Chống đụng nhau với agent khác
- Tạo nhánh riêng ngay: `git switch -c <ten-viec>`. **Không commit trực tiếp lên `main`.**
- Nếu sổ cho thấy agent khác đang sửa cùng vùng → nói cho người dùng biết, đề nghị đổi việc.
- Vùng an toàn để làm song song: `server/` (game) và `web/` là hai vùng tách biệt. Hai
  agent cùng vào `server/internal/game/` thì **chắc chắn** đụng.

## 2. Nếu việc liên quan giao thức nhị phân
Tra `docs/superpowers/specs/` (bảng byte đã giải mã + kiểm chéo). Không có thì đọc code
Java cũ trong `1-MAY-CHU-GAME/src/`. **Đừng tự đoán byte** — sai một byte là client hỏng
dòng dữ liệu vĩnh viễn, và lỗi chỉ lộ ra rất muộn.

## 3. Báo cáo rồi mới sửa code
Trả lời **tiếng Việt đời thường** (chủ dự án không phải lập trình viên):
1. Hiện trạng: đã xong gì, ai đang làm gì.
2. Việc sẽ làm và **tại sao** nó là bước hợp lý tiếp theo.
3. Cạm bẫy liên quan đã đọc ở §6/§9 và sẽ tránh thế nào.

Nếu người dùng không nói làm gì: đề xuất việc hợp lý nhất theo
`docs/superpowers/plans/LO-TRINH.md` và §5 của `BAT-DAU-TU-DAY.md`.

## 4. Cách làm việc bắt buộc
- **TDD**: test đỏ trước, code sau.
- **Kiểm thử đột biến**: xong thì cố tình phá code xem test có bắt. Không bắt = test vô
  dụng, phải viết lại test. Dự án này đã **3 lần** bị "test xanh mà code sai".
- Chú thích code bằng **tiếng Việt, giải thích *tại sao*** chứ không phải *cái gì*.
- Xong việc → gọi lại `/viec` (nó sẽ tự vào chế độ GHI SỔ).

---

# PHẦN C2 — Chế độ GHI SỔ

## 1. Xem chính xác vừa làm gì
`git status --short`, `git diff --stat` (cả staged và chưa) để biết đã chạm file nào.

## 2. Kiểm chứng — không có kết quả thì không được ghi "xong"

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

Test đỏ → ghi vào sổ là **đang dở, đỏ ở đâu**. Đừng ghi là xong.

## 3. Ghi `NHAT-KY-AGENT.md`
Thêm **một dòng vào CUỐI bảng** (mới nhất nằm dưới cùng — hook đọc 12 dòng cuối):

`| <dd/mm/yyyy> | <ai> | <làm gì, một câu> | <nhánh> / <commit> | <còn dở gì, hoặc —> |`

Cột "còn dở / cảnh báo" quan trọng nhất: nhánh chưa gộp, test đang đỏ, file sửa nửa vời,
quyết định còn treo. Không có thì ghi `—`.

## 4. Cập nhật `BAT-DAU-TU-DAY.md` (theo §11 của chính nó)
- **§4 Đã làm** — thêm dòng nếu vừa xong một mốc, kèm mã commit.
- **§5 Còn lại** — xoá việc vừa xong, thêm việc mới lộ ra.
- **§6 Lỗi sai** — vấp lỗi mới thì ghi lỗi + **vì sao vấp** + cách tránh. Phần giá trị nhất
  của cả tài liệu; đừng bỏ qua vì "lỗi nhỏ".
- **§12 Nhật ký** + dòng "Cập nhật lần cuối" ở đầu file.

## 5. Commit
Code + tài liệu **cùng một commit** (để lịch sử khớp nhau). Thông điệp **không dấu**:
`<loai>(<vung>): <viec da lam>` — `<loai>` = `feat`/`fix`/`test`/`docs`/`refactor`,
`<vung>` = `game`/`web`/`docs`.

Nếu đang ở nhánh riêng, việc xong hẳn và test xanh → nói với người dùng là **gộp được rồi**,
nhưng **đừng tự gộp lên `main`** khi chưa được bảo, vì agent khác có thể đang làm dở.

## 6. Báo lại bằng tiếng Việt
Ba câu: vừa làm gì · kiểm chứng bằng **con số thật** (không nói chung chung) · việc tiếp
theo nên là gì.
