# Thiết kế: Lưu trạng thái sống (KH11)

**Ngày:** 2026-07-27 · **Trạng thái:** Chờ chủ dự án chốt 3 con số ở §9
**Viết cho người không lập trình** — nói *tại sao* trước khi nói *làm thế nào*.

---

## 1. Vấn đề — nói thẳng

**Máy chủ hiện không lưu gì cả sau lúc tạo nhân vật.**

Đo thật, không phải suy đoán: `server/internal/store/store.go` có đúng **4 hàm**
— `LayTaiKhoan`, `TaoNhanVat`, `LayNhanVatTheoTaiKhoan`, `TenNhanVatDaDung`.
**Không có một hàm cập nhật nào.** Interface `Kho` mà tầng game dùng
(`login.go:64-69`) cũng chỉ khai 4 hàm đó — nghĩa là tầng game **không có cách
nào** ghi xuống cơ sở dữ liệu kể cả khi muốn.

Người chơi đi lại, mất máu, hồi sinh — tất cả chỉ nằm trong bộ nhớ. Tắt máy chủ
là **bay sạch**, người chơi đăng nhập lại thấy mình y như lúc mới tạo nhân vật.

Tệ hơn: `cmd/nroserver/main.go` có bắt tín hiệu tắt (Ctrl+C / SIGTERM) nhưng
kết thúc bằng `os.Exit(0)` — mà `os.Exit` **giết tiến trình ngay**, mọi việc dọn
dẹp đang chờ ở các luồng khác **không chạy**. Nên hiện tại tắt máy chủ "đàng
hoàng" cũng mất dữ liệu y như rút điện.

**Đây là lỗ hổng kiến trúc, không nhét sau được** — vì cách lưu quyết định hình
dạng của tầng khu vực. Phải làm **trước KH7 (vật phẩm)**, nếu không mỗi hệ sau
lại tự chế một kiểu lưu, cuối cùng chắp vá mà vẫn mất dữ liệu.

---

## 2. Cái gì cần lưu

### 2.1 Nguồn sự thật

Trạng thái sống của một người chơi nằm trong **một** đối tượng `store.NhanVat`.
Có **hai nơi giữ con trỏ tới cùng đối tượng đó**:

- `BoXuLy.nhanVat` — thuộc luồng phiên (mỗi kết nối một cái)
- `KhuVuc.nguoiTrong[id].nv` — thuộc luồng actor của khu vực

Chúng là **cùng một ô nhớ**, không phải hai bản sao. Từ khi người chơi vào khu
vực, **actor là bên duy nhất còn ghi**; luồng phiên chỉ đọc `ID` (bất biến).

### 2.2 Trường thay đổi khi chơi

Hiện tại (KH6) chỉ **4 trường** thật sự đổi:

| Trường | Đổi khi nào |
|---|---|
| `ViTriX`, `ViTriY` | Di chuyển |
| `HpHienTai` | Bị đánh, hồi sinh |
| `ChakraHienTai` | Hồi sinh |

**Nhưng đừng chỉ lưu 4 trường này.** KH7 (vật phẩm) và KH8 (số liệu) sắp nối
dây cho `Cap`, `KinhNghiem`, `Vang`, `Ngoc`, `MapID`, `HpToiDa`… — danh sách sẽ
phình ra. Ghi **cả struct** thì thêm trường mới không phải sửa chỗ lưu.

### 2.3 Cái KHÔNG lưu

Bốn trường riêng của `nguoiTrongKhuVuc` — `chieuDangChon`, `loaiChieuDangChon`,
`lanDanhCuoi`, `lucHoiSinh` — đều là **trạng thái tạm của phiên chiến đấu**
(hồi chiêu vài giây, bất tử 1,5 giây). Hết ý nghĩa ngay khi tắt máy chủ. Không lưu.

---

## 3. Nguyên tắc số một: actor KHÔNG BAO GIỜ chờ ổ đĩa

Mỗi khu vực có **một** goroutine ("người quản lý") sở hữu toàn bộ trạng thái.
Nếu nó dừng lại chờ ghi cơ sở dữ liệu mà đĩa chậm, **cả bản đồ đứng hình** —
không ai đi được, không ai đánh được, quái ngừng hồi sinh.

Đây **không phải lo xa**: dự án đã có đúng lỗi này đang sống — `§6.10b` của
`BAT-DAU-TU-DAY.md` ghi actor bị chặn tới **30 giây** vì tự đi gửi gói mạng cho
một client ngừng đọc. Lặp lại kiểu đó với cơ sở dữ liệu là chắc chắn hỏng.

**Luật:** actor chỉ **chụp ảnh** trạng thái (bản sao giá trị) rồi **ném vào hàng
đợi**, xong đi tiếp. Một luồng riêng ("thợ ghi") lo phần chạm đĩa.

```
   Actor khu vực                Hàng đợi              Thợ ghi
   ─────────────                ────────              ───────
   chụp ảnh (bản sao)  ──đẩy──►  [ ][ ][ ]  ──lấy──►  ghi PostgreSQL
   đi tiếp ngay                                       (chờ đĩa ở đây, không sao)
   KHÔNG BAO GIỜ CHỜ
```

---

## 4. Ba thời điểm lưu

| # | Khi nào | Vì sao cần | Mất tối đa bao nhiêu |
|---|---|---|---|
| **A** | **Định kỳ** (auto-save) | Máy chủ sập/mất điện thì vẫn còn bản gần nhất | Bằng nhịp lưu (xem §9) |
| **B** | **Người chơi rời** (mất kết nối / thoát) | Đây là lúc dữ liệu quan trọng nhất — người chơi vừa chơi xong | ~0 |
| **C** | **Tắt máy chủ** (Ctrl+C / SIGTERM) | Bảo trì hằng ngày không được làm mất tiến độ ai | ~0 |

Bản Java cũ có A và B, **thiếu C** — nút "Thoát" trên cửa sổ gọi thẳng
`System.exit(0)` không lưu ai. Còn hàm `saveData()` làm đúng thì **tồn tại nhưng
không gắn vào nút nào** (mã chết). Bản Go phải có đủ cả ba.

### Chi tiết từng điểm móc

- **A** — một goroutine hẹn giờ, cứ tới nhịp thì xin ảnh chụp **tất cả** người
  trong **tất cả** khu vực, đẩy vào hàng đợi.
- **B** — `BoXuLy.KhiDong()` (`login.go:405-413`) là nơi **duy nhất** biết một
  phiên vừa đứt. Hiện nó chỉ gửi sự kiện "rời khu vực". Thêm bước xin ảnh chụp
  **trước** khi gửi sự kiện rời — vì sau khi rời, actor đã xoá người đó khỏi
  danh sách, không chụp được nữa.
- **C** — trong `main.go`, nhánh nhận tín hiệu tắt: xin ảnh chụp toàn bộ → chờ
  thợ ghi làm xong → **rồi mới** `DungTatCa()`. Thứ tự này bắt buộc: gọi
  `DungTatCa()` trước thì actor đã chết, không ai trả lời yêu cầu chụp ảnh nữa.

---

## 5. Cách lấy ảnh chụp — theo khuôn đã có

Dự án đã có sẵn khuôn "sự kiện truy vấn" ở `suKienXinDanhSachQuai`
(`khuvuc.go:223-283`). Làm y hệt cho người chơi:

```go
type suKienXinAnhChup struct {
    idNhanVat int64                 // 0 = xin tất cả
    traLoi    chan []store.NhanVat  // BẢN SAO GIÁ TRỊ, không phải con trỏ
}
```

Ba điểm bắt buộc, đều học từ khuôn cũ:

1. **Trả bản sao giá trị** (`ban := *nguoi.nv`), tuyệt đối không đưa con trỏ ra
   ngoài — bên nhận chạy ở luồng khác, giữ con trỏ là tranh chấp dữ liệu.
2. **Kênh trả lời có đệm 1** — để actor gửi xong đi tiếp ngay, kể cả khi bên hỏi
   đã bỏ đi vì hết giờ chờ.
3. **Có hạn giờ chờ** — actor đã dừng thì bên hỏi không treo vĩnh viễn.

⚠️ **Một cạm bẫy phải xử:** `Gui()` hiện **rơi im lặng** khi hộp thư đầy
(`khuvuc.go:1031-1039`). Nghĩa là một yêu cầu chụp ảnh có thể **bị mất** mà
không ai biết → khu vực đó bị bỏ qua trong lượt lưu. Phải **đếm và ghi log** số
lần rơi, không được im.

---

## 6. Ghi xuống cơ sở dữ liệu

### 6.1 Hàm mới trong `store`

```go
func (s *Store) CapNhatNhanVat(ctx context.Context, nv NhanVat) error
```

Viết theo đúng khuôn có sẵn: SQL thô với `$1,$2…`, lỗi bọc
`fmt.Errorf("...that bai: %w", err)`, tên tiếng Việt không dấu.

**Không cần migration mới** cho bước tối thiểu — mọi cột cần ghi đã có trong
bảng `nhan_vat`. Cột `dang_nhap_cuoi` cũng **đã có sẵn** nhưng chưa ai dùng;
tiện thì ghi luôn thời điểm lưu gần nhất vào đó.

### 6.2 Bắt buộc có hạn giờ

Hiện **toàn bộ mã nguồn không có một `context.WithTimeout` nào** — mọi truy vấn
dùng `context.Background()` trần, tức là **không gì huỷ được nó** nếu cơ sở dữ
liệu treo. Với việc lưu định kỳ thì đó là hỏng: một lần ghi treo sẽ làm nghẽn
hàng đợi vô hạn.

→ `CapNhatNhanVat` **phải** được gọi kèm hạn giờ (đề xuất 5 giây). Đây sẽ là chỗ
đầu tiên trong dự án dùng `context.WithTimeout` — chưa có tiền lệ để bắt chước.

### 6.3 Lỗi ghi KHÔNG được nuốt im lặng

Bài học đắt nhất từ bản Java: `PlayerDAO.updatePlayer` bọc toàn bộ trong một
`try/catch`, lỗi chỉ ghi log rồi thôi — **không ném lại, không thử lại, không
đánh dấu "cần lưu lại"**. Cơ sở dữ liệu chập chờn đúng lúc lưu là người chơi mất
sạch thay đổi mà **không ai biết**.

Bản Go: ghi hỏng thì **đếm, log mức cảnh báo, và thử lại ở lượt sau** (dữ liệu
vẫn còn trong bộ nhớ nên lượt sau chụp lại là có bản mới hơn). Nếu hỏng liên
tiếp N lần thì log mức lỗi để người vận hành thấy.

---

## 7. Hai cái bẫy phải chặn

### 7.1 Ghi đè bản mới bằng bản cũ

Bản Java có lỗi này **và chưa bao giờ vá**: khi có phiên thứ hai đăng nhập cùng
tài khoản, nó **đọc dữ liệu từ cơ sở dữ liệu TRƯỚC** rồi mới đá phiên cũ (mà việc
đá phiên cũ mới là lúc phiên cũ ghi xuống). Kết quả: phiên mới chạy trên dữ liệu
**cũ hơn**, rồi ghi đè lên bản mới → mất đồ, mất tiến độ.

Bản Go tránh bằng thứ tự: **đá phiên cũ và chờ nó lưu xong → rồi mới đọc dữ liệu
cho phiên mới.**

### 7.2 Một tài khoản, hai phiên cùng lưu

Hiện hai phiên cùng tài khoản tạo ra **hai đối tượng `store.NhanVat` độc lập** —
cả hai cùng được lưu, đè nhau. Bản Java vá bằng **ba lớp chắp vá** chồng lên
nhau (khoá theo tên, quét dọn, và cuối cùng là *khoá luôn tài khoản* khi phát
hiện trùng) — dấu hiệu điển hình của việc **không có một nguồn sự thật duy nhất**.

Bản Go làm đúng từ đầu: **một chỗ duy nhất** biết "tài khoản X đang ở phiên nào".
Đăng nhập mới thì đá phiên cũ, chờ lưu xong, rồi mới cho vào.

**Phạm vi:** đây là việc riêng nhưng **phải làm cùng KH11**, vì không có nó thì
chính việc lưu sẽ làm hỏng dữ liệu.

---

## 8. Cái KHÔNG làm trong KH11 (để sau)

- **Ghi ngay khi có sự kiện quý** (nhặt đồ hiếm) — chờ KH7 có vật phẩm đã.
- **Nhật ký giao dịch / chống nhân bản** — thuộc KH7.
- **Công cụ quản lý phiên bản cơ sở dữ liệu (migration)** — đáng làm sớm nhưng
  là việc riêng, không chặn KH11.
- **Sao lưu tự động** — thuộc Kế hoạch V (vận hành).
- **Tách gửi gói mạng khỏi actor** (§6.10b) — thuộc Kế hoạch B. KH11 không được
  làm nó tệ hơn.

---

## 9. ✅ BA CON SỐ — ĐÃ CHỐT (27/07/2026, chủ dự án)

| # | Quyết định | Vì sao |
|---|---|---|
| 1 | **Nhịp lưu = 60 giây** | Sập thì mất tối đa 1 phút tiến độ. Với 100 người online chỉ ~1,7 lần ghi/giây — rất nhẹ. *(Bản Java cũ để 5 phút và người chơi từng mất đồ vì nó.)* **Đặt thành hằng số một chỗ** để đổi được dễ |
| 2 | **Hàng đợi đầy → gộp, giữ bản mới nhất của mỗi người** | Bản mới luôn tốt hơn bản cũ nên bỏ bản cũ không mất gì. Cách khác (chặn actor lại chờ) thì vi phạm §3 — điều tuyệt đối không được làm |
| 3 | **Đăng nhập trùng → đá phiên cũ, chờ lưu xong, rồi mới cho phiên mới vào** | Không đá thì hai phiên ghi đè nhau, hỏng dữ liệu. Bản Java từng phải *khoá luôn tài khoản* khi phát hiện trùng vì không xử được gốc — ta không lặp lại |

---

## 10. Nghiệm thu

**Tự động (test):**
1. Người chơi đi lại → tới nhịp lưu → dữ liệu trong cơ sở dữ liệu đúng vị trí mới.
2. Người chơi mất máu rồi ngắt kết nối → đăng nhập lại thấy đúng máu đã mất.
3. Tắt máy chủ (tín hiệu) khi có người đang chơi → dữ liệu của họ đã xuống đĩa.
4. **Actor không bao giờ bị chặn**: cơ sở dữ liệu chậm (giả lập) mà khu vực vẫn
   chạy nhịp bình thường, người chơi vẫn đi lại được.
5. Đăng nhập trùng: phiên cũ bị đá, dữ liệu phiên cũ được lưu **trước** khi phiên
   mới đọc — không mất tiến độ.

**Tay (chủ dự án):** chơi vài phút, tắt máy chủ, bật lại, đăng nhập — thấy đúng
chỗ mình đứng và đúng máu.
