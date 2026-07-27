# Kế hoạch 11 — Lưu trạng thái sống: Kế hoạch thực thi

> **Cho người thực thi:** KỸ NĂNG BẮT BUỘC `superpowers:subagent-driven-development`.
> **ĐỌC TRƯỚC:** đặc tả [`../specs/2026-07-27-luu-trang-thai-song.md`](../specs/2026-07-27-luu-trang-thai-song.md).

**Mục tiêu:** người chơi đi lại, mất máu, chết đi sống lại — tắt máy chủ rồi bật
lại, đăng nhập vào **thấy đúng chỗ mình đứng và đúng máu mình còn**.

**Kiến trúc:** actor chụp ảnh (bản sao giá trị) → hàng đợi → thợ ghi riêng chạm
đĩa. **Actor không bao giờ chờ ổ đĩa.**

**Công nghệ:** Go 1.25+, chỉ `pgx/v5` + stdlib. Chỉ động vào `server/`.

---

## Ràng buộc chung (Global Constraints)

- **LUẬT SỐ MỘT — actor KHÔNG BAO GIỜ chờ ổ đĩa.** Mọi lời gọi cơ sở dữ liệu
  phải nằm ngoài goroutine `vongLap`. Vi phạm luật này là lặp lại đúng lỗi
  Nghiêm trọng #1 đang sống trong dự án (§6.10b: actor bị chặn 30 giây vì tự đi
  gửi gói mạng). **Mỗi task đụng actor phải tự soi và ghi bằng chứng vào báo cáo.**
- **Kỷ luật actor:** chỉ `vongLap` chạm trạng thái khu vực. Ảnh chụp phải là
  **bản sao GIÁ TRỊ** (`ban := *nguoi.nv`), tuyệt đối không đưa con trỏ ra ngoài.
  Theo khuôn `suKienXinDanhSachQuai` (`khuvuc.go:223-283`).
- **Kênh trả lời có đệm 1** + **có hạn giờ chờ** — như khuôn cũ, để không treo.
- **Nhịp lưu = 60 giây**, đặt thành **một hằng số duy nhất** có tên rõ ràng.
- **Hàng đợi đầy → gộp theo id nhân vật, giữ bản mới nhất.** Không chặn actor.
- **Mọi truy vấn ghi PHẢI có hạn giờ** (`context.WithTimeout`, đề xuất 5 giây) —
  hiện toàn bộ mã nguồn dùng `context.Background()` trần, chưa có tiền lệ.
- **Lỗi ghi KHÔNG được nuốt im lặng** — đếm, log, thử lại lượt sau. Bài học từ
  bản Java (`PlayerDAO.updatePlayer` nuốt lỗi → mất dữ liệu mà không ai biết).
- **`Gui()` rơi im lặng khi hộp thư đầy** (`khuvuc.go:1031-1039`) — yêu cầu chụp
  ảnh có thể mất. Phải **đếm và log**, không được im.
- **TDD bắt buộc** (test đỏ trước), PASS `-count=5`. **Kiểm thử đột biến bắt
  buộc** sau mỗi task: cố tình phá code, test **phải đỏ**. Không đỏ = test vô
  dụng, viết lại. *(Dự án đã 3 lần dính "test xanh mà code sai".)*
- **Test giá trị ghim cứng**, không lấy từ chính nguồn đang kiểm.
- **Go:** `D:\Go\bin\go.exe`, gofmt `D:\Go\bin\gofmt.exe`. Windows PowerShell
  (KHÔNG `&&`). Cổng PostgreSQL **TỰ KIỂM**, đừng ghim số:
  `Get-NetTCPConnection -State Listen | Where-Object { $_.LocalPort -in 5432,5433 }`
- **Ghi chú tiếng Việt giải thích *tại sao***, không phải *cái gì*.

---

## Thứ tự và lý do

| Task | Nhóm | Phụ thuộc |
|---|---|---|
| 1 Hàm `CapNhatNhanVat` trong store | Nền | — |
| 2 Sự kiện chụp ảnh trong actor | Nền | — |
| 3 Thợ ghi + hàng đợi gộp | Lõi | 1 |
| 4 Lưu định kỳ 60 giây | Nối | 2, 3 |
| 5 Lưu khi người chơi rời | Nối | 2, 3 |
| 6 Lưu khi tắt máy chủ | Nối | 2, 3 |
| 7 Một tài khoản một phiên | Chặn hỏng | 5 |
| 8 Test đầu-cuối tắt–bật | — | tất cả |

**Chạy song song được:** Task 1 và Task 2 (khác file hoàn toàn).

---

### Task 1: Hàm `CapNhatNhanVat` trong store

**Files:** Modify `server/internal/store/store.go`, `store_test.go`

**Interfaces:**
- Produces: `func (s *Store) CapNhatNhanVat(ctx context.Context, nv NhanVat) error`

- [ ] **Step 1: Test đỏ** — tạo nhân vật, đổi vị trí + máu, gọi `CapNhatNhanVat`,
      đọc lại bằng `LayNhanVatTheoTaiKhoan`, khẳng định **từng trường** đúng giá
      trị mới (ghim cứng). Test thứ hai: cập nhật nhân vật **không tồn tại** →
      trả `ErrKhongThayNhanVat` (không im lặng coi như xong).
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt.** Ghi **mọi trường có thể đổi** (không chỉ 4 trường đang
      đổi hôm nay) — *tại sao:* KH7/KH8 sắp nối dây `Cap`, `KinhNghiem`, `Vang`,
      `MapID`…; ghi cả struct thì thêm trường mới không phải sửa chỗ lưu.
      **KHÔNG** ghi `ID`, `TaiKhoanID`, `Ten` (bất biến — ghi lại là mở cửa cho
      lỗi đổi tên/đổi chủ). Ghi luôn `dang_nhap_cuoi = now()` (cột đã có sẵn
      trong schema, chưa ai dùng). Theo khuôn SQL của `TaoNhanVat`.
- [ ] **Step 4: Chạy — PASS.**
- [ ] **Step 5: Kiểm thử đột biến** — bỏ một cột khỏi câu UPDATE → test phải đỏ.
- [ ] **Step 6: Commit** — `feat(store): ham CapNhatNhanVat luu trang thai xuong CSDL`

---

### Task 2: Sự kiện chụp ảnh người chơi trong actor

**Files:** Modify `server/internal/game/khuvuc.go`, `khuvuc_test.go`

**Interfaces:**
- Produces: `type suKienXinAnhChup struct { idNhanVat int64; traLoi chan []store.NhanVat }`
- Produces: `func (k *KhuVuc) AnhChupNguoiChoi(idNhanVat int64) []store.NhanVat`
  *(idNhanVat = 0 → xin tất cả người trong khu vực)*

- [ ] **Step 1: Test đỏ** — hai người trong khu vực; xin ảnh chụp tất cả → được
      2 bản; xin theo id → được đúng 1 bản đúng người. Test thứ ba **(quan trọng
      nhất)**: sửa bản chụp trả về **KHÔNG** làm đổi trạng thái trong actor —
      khẳng định đây là **bản sao giá trị**, không phải con trỏ.
      Test thứ tư: actor chưa chạy → trả nil ngay, không treo.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt** theo đúng khuôn `suKienXinDanhSachQuai` (kênh đệm 1,
      hạn giờ chờ, sao chép giá trị `ban := *nguoi.nv`).
      **Đếm số lần `Gui()` rơi** (hộp thư đầy) và log — *tại sao:* rơi im lặng
      nghĩa là một khu vực bị bỏ qua trong lượt lưu mà không ai biết.
- [ ] **Step 4: Chạy — PASS** với `-count=5`.
- [ ] **Step 5: Soi mã** — xác nhận không đường nào ngoài `vongLap` chạm
      `nguoiTrong`. Ghi bằng chứng vào báo cáo.
- [ ] **Step 6: Kiểm thử đột biến** — đổi sao-chép-giá-trị thành trả thẳng con
      trỏ → test thứ ba phải đỏ.
- [ ] **Step 7: Commit** — `feat(game): su kien chup anh nguoi choi qua hop thu actor`

---

### Task 3: Thợ ghi + hàng đợi gộp

**Files:** Create `server/internal/game/thogi.go`, `thogi_test.go`

**Interfaces:**
- Produces: `type ThoGhi struct{...}`
- Produces: `func MoThoGhi(kho KhoLuu, soThoGhi int) *ThoGhi`
- Produces: `func (t *ThoGhi) Xep(nv store.NhanVat)` — **không bao giờ chặn**
- Produces: `func (t *ThoGhi) DongVaChoXong(ctx context.Context) error`
- Produces: `type KhoLuu interface { CapNhatNhanVat(context.Context, store.NhanVat) error }`
  *(interface hẹp để test tiêm được kho giả — theo khuôn interface `Kho` ở login.go)*

- [ ] **Step 1: Test đỏ** — dùng kho giả:
      (a) xếp 3 bản của **cùng một người** → chỉ ghi **bản mới nhất** (gộp);
      (b) `Xep` **không bao giờ chặn** kể cả khi kho giả cố tình chậm — đo bằng
          thời gian gọi `Xep` (ghim cứng: dưới 50ms dù kho chậm 2 giây);
      (c) kho trả lỗi → **thử lại** ở lượt sau, và **đếm lỗi**, không nuốt;
      (d) `DongVaChoXong` ghi hết những gì còn trong hàng đợi rồi mới trả về.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt.** Hàng đợi là **map id → bản mới nhất** (gộp tự nhiên),
      không phải danh sách — *tại sao:* bản mới luôn tốt hơn bản cũ nên giữ bản
      cũ là vô nghĩa, và map làm hàng đợi không thể phình vô hạn.
      Mỗi lần ghi bọc `context.WithTimeout` (5 giây).
- [ ] **Step 4: Chạy — PASS** với `-count=5`.
- [ ] **Step 5: Kiểm thử đột biến** — bỏ hạn giờ ghi → test (b) phải đỏ;
      bỏ gộp (đổi map thành slice) → test (a) phải đỏ.
- [ ] **Step 6: Commit** — `feat(game): tho ghi + hang doi gop, actor khong cho o dia`

---

### Task 4: Lưu định kỳ 60 giây

**Files:** Modify `server/internal/game/sokhuvuc.go` (hoặc nơi hợp lý), test kèm

- [ ] **Step 1: Test đỏ** — đồng hồ tiêm được (như KH6 Task 8): chưa tới 60 giây
      thì **chưa** lưu; đúng/qua 60 giây thì mọi người trong mọi khu vực được xếp
      vào thợ ghi. Ghim cứng mốc 60 giây, **không** đọc lại hằng số từ code.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt.** Một goroutine hẹn giờ, xin ảnh chụp từng khu vực rồi
      `Xep`. **Không** chạy trong `nhipTim` của actor — *tại sao:* `nhipTim` chạy
      100ms/lần trong `vongLap`, nhét việc lưu vào đó là kéo actor vào đường ghi đĩa.
- [ ] **Step 4: Chạy — PASS** với `-count=5`.
- [ ] **Step 5: Kiểm thử đột biến** — đổi 60 giây thành 600 → test biên phải đỏ.
- [ ] **Step 6: Commit** — `feat(game): luu dinh ky 60 giay`

---

### Task 5: Lưu khi người chơi rời

**Files:** Modify `server/internal/game/login.go` (`KhiDong`), test tích hợp

- [ ] **Step 1: Test đỏ** — client giả vào khu vực, di chuyển, **ngắt kết nối** →
      kho giả nhận đúng bản ghi với **vị trí mới nhất**. Test thứ hai: người chưa
      vào khu vực mà ngắt → không lưu bừa, không panic.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt.** Xin ảnh chụp **TRƯỚC** khi gửi `suKienRa` — *tại sao:*
      sau khi actor xử lý `suKienRa`, người đó đã bị xoá khỏi danh sách, không
      chụp được nữa. Thứ tự này là bắt buộc, ghi chú rõ trong code.
- [ ] **Step 4: Chạy — PASS** với `-count=5`.
- [ ] **Step 5: Kiểm thử đột biến** — đảo thứ tự (gửi `suKienRa` trước rồi mới
      chụp) → test phải đỏ.
- [ ] **Step 6: Commit** — `feat(game): luu khi nguoi choi roi khu vuc`

---

### Task 6: Lưu khi tắt máy chủ

**Files:** Modify `server/cmd/nroserver/main.go`, test kèm

- [ ] **Step 1: Test đỏ** — gọi hàm tắt (tách ra được để test) khi có người trong
      khu vực → kho giả nhận đủ bản ghi của họ **trước khi** khu vực bị dừng.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt.** Trong nhánh nhận tín hiệu: xin ảnh chụp toàn bộ → xếp
      → `DongVaChoXong` (có hạn giờ, đề xuất 30 giây) → **rồi mới** `DungTatCa()`
      → `kho.Close()` → thoát. *Tại sao thứ tự này:* gọi `DungTatCa()` trước thì
      actor đã chết, không ai trả lời yêu cầu chụp ảnh — đúng lỗi bản Java cũ
      (nút "Thoát" gọi thẳng `System.exit` không lưu ai, còn hàm lưu đúng thì tồn
      tại nhưng **quên gắn vào nút nào**).
      Tách phần dọn dẹp ra một hàm gọi được từ test, đừng để dính chặt `os.Exit`.
- [ ] **Step 4: Chạy — PASS.**
- [ ] **Step 5: Kiểm thử đột biến** — chuyển `DungTatCa()` lên trước bước lưu →
      test phải đỏ.
- [ ] **Step 6: Commit** — `feat(server): luu toan bo nguoi choi khi tat may chu`

---

### Task 7: Một tài khoản — một phiên sống

**Files:** Modify `server/internal/game/login.go` + nơi giữ danh sách phiên, test kèm

**Vì sao cần trong KH11:** hai phiên cùng tài khoản tạo **hai đối tượng
`store.NhanVat` độc lập**, cả hai cùng được lưu → **ghi đè nhau, hỏng dữ liệu**.
Không chặn thì chính việc lưu làm hỏng thứ nó định bảo vệ.

- [ ] **Step 1: Test đỏ** — tài khoản A đăng nhập, di chuyển; tài khoản A đăng
      nhập lần hai → phiên cũ **bị đá**, dữ liệu phiên cũ **đã được lưu**, rồi
      phiên mới mới đọc dữ liệu. Khẳng định phiên mới đọc được **vị trí mới nhất**
      (không phải vị trí cũ trong cơ sở dữ liệu).
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt.** **Một** chỗ duy nhất biết "tài khoản X đang ở phiên
      nào" (map có khoá, hoặc actor riêng). Thứ tự bắt buộc: **đá phiên cũ → chờ
      nó lưu xong → rồi mới đọc dữ liệu cho phiên mới.**
      *Tại sao:* bản Java làm ngược (đọc trước, đá sau) → phiên mới chạy trên dữ
      liệu cũ rồi ghi đè lên bản mới → mất đồ. Lỗi đó **chưa bao giờ được vá** ở
      bản Java. Cũng đừng lặp lại cách bản Java "vá" bằng ba lớp chồng nhau — một
      nguồn sự thật là đủ.
- [ ] **Step 4: Chạy — PASS** với `-count=5`.
- [ ] **Step 5: Kiểm thử đột biến** — đảo thứ tự (đọc trước, đá sau) → test phải đỏ.
- [ ] **Step 6: Commit** — `feat(game): mot tai khoan mot phien, da phien cu va cho luu xong`

---

### Task 8: Test đầu-cuối tắt–bật

**Files:** Modify `server/internal/game/tichhop_test.go`

- [ ] **Step 1: Test đỏ** — cơ sở dữ liệu **thật**: client giả đăng nhập, vào khu
      vực, **di chuyển tới vị trí ghim cứng**, bị mất máu (ghim cứng) → chạy
      đường tắt máy chủ → **mở lại** kho và đọc nhân vật → đúng vị trí, đúng máu.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Sửa** những gì lộ ra.
- [ ] **Step 4: Chạy — PASS.** Sau đó **toàn bộ** `go test ./... -count=5` phải
      xanh không chập chờn; `gofmt -l` sạch; `go vet` sạch.
- [ ] **Step 5: Commit** — `test(game): dau-cuoi luu trang thai qua tat-bat may chu`

---

### Task 9 (THỦ CÔNG — chủ dự án): nghiệm thu tay

- [ ] Chạy máy chủ, đăng nhập bằng client thật, đi tới một chỗ dễ nhớ, đánh quái
      cho mất máu.
- [ ] **Tắt máy chủ bằng Ctrl+C**, bật lại, đăng nhập lại.
- [ ] **Có đứng đúng chỗ cũ không? Máu có đúng như lúc thoát không?**

---

## Tự soi (sau khi viết kế hoạch)

**Chỗ dễ hỏng nhất, theo thứ tự:**

1. **Actor lỡ chờ ổ đĩa** (Task 3, 4) — lỗi này im lặng cho tới khi đông người
   thì cả bản đồ đứng hình. Test (b) của Task 3 là lưới an toàn chính.
2. **Thứ tự chụp-ảnh-trước-khi-rời** (Task 5) — làm ngược thì lưu ra dữ liệu cũ,
   mà test hời hợt vẫn xanh vì *có* lưu.
3. **Thứ tự đá-phiên-cũ-trước-khi-đọc** (Task 7) — chính là lỗi chưa vá của bản
   Java. Sai thì mất đồ, và chỉ lộ khi có người đăng nhập hai nơi.
4. **Lỗi ghi bị nuốt** (Task 3) — hỏng mà không ai biết, đúng bài học bản Java.
5. **`Gui()` rơi im lặng** (Task 2) — một khu vực bị bỏ qua trong lượt lưu.

**Đã cố ý bỏ ra ngoài phạm vi:** ghi ngay khi nhặt đồ hiếm (KH7), nhật ký giao
dịch (KH7), công cụ migration, sao lưu tự động (Kế hoạch V), tách gửi gói mạng
khỏi actor (Kế hoạch B).
