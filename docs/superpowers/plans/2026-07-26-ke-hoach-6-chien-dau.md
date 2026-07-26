# Kế hoạch 6 — Chiến đấu, chiêu thức, quái: Kế hoạch thực thi

> **Cho người thực thi:** KỸ NĂNG BẮT BUỘC: dùng `superpowers:subagent-driven-development`
> để thực thi từng task. Các bước dùng checkbox (`- [ ]`).
>
> **ĐỌC TRƯỚC:** đặc tả [`../specs/2026-07-26-giao-thuc-chien-dau.md`](../specs/2026-07-26-giao-thuc-chien-dau.md)
> — mọi bảng byte nằm ở đó. Nghiên cứu gốc:
> [`../research/kh6-goi-chien-dau.md`](../research/kh6-goi-chien-dau.md). **Đừng tự đoán byte.**

**Mục tiêu:** người chơi chọn được chiêu, đánh quái, quái đánh lại, quái chết rồi hồi sinh,
người chơi đánh nhau được, chết rồi hồi sinh được. Nhiều người trong cùng khu vực nhìn
thấy toàn bộ diễn biến đó.

**Kiến trúc:** quái trở thành **một phần trạng thái của actor khu vực** (KH5, `27d66fc`) —
không luồng mới, không khoá. Móc `nhipTim` để rỗng ở KH5 chính là chỗ AI quái sẽ vào.
Luồng phiên chỉ **gửi sự kiện** vào hộp thư actor, không chạm trạng thái.

**Công nghệ:** Go 1.25+, chỉ `pgx/v5` + stdlib. File dưới `server/`. Không đụng `1-MAY-CHU-GAME/`.

---

## Ràng buộc chung (Global Constraints)

- **Byte-exact:** mọi gói phải khớp bảng ở đặc tả §3 từng byte, từng kiểu.
- **Máu và sát thương là `int32` big-endian (4 byte)**, kẹp trần `2_147_483_647`,
  kẹp sàn `0`. Không dùng 8 byte. Ngoại lệ duy nhất: `powRequire` trong catalog chiêu
  là **`long` 8 byte**.
- **BA GÓI ĐỘ DÀI BIẾN THIÊN** — rủi ro số một của kế hoạch này:
  `-60` (hai dạng, tuỳ `docTiep`), `-45/0` (hai vòng lặp lồng nhau), `-124` (ba dạng).
  Cộng thêm gói **lên** `54` (1 byte hoặc 5 byte tuỳ `mobId == -1`).
  → Mỗi dạng **một test byte-chính-xác riêng**, không gộp.
- **`-60` phải là HAI HÀM dựng gói riêng biệt**, không một hàm có cờ. Một hàm có cờ thì
  người sửa sau rất dễ ghi thừa trường; hai hàm thì lỗi đó không tồn tại được.
- **Sát thương không vượt máu còn lại**; **quái đầy máu không chết bởi một đòn**
  (`dame = hp - 1`).
- **Kỷ luật actor:** chỉ luồng `vongLap` của khu vực chạm trạng thái quái và trạng thái
  chiến đấu — **kể cả *chiêu đang chọn* của người chơi**. Máy này KHÔNG chạy được
  `-race` (thiếu gcc, không quyền admin) → thay bằng **soi mã từng đường chạm trạng thái**
  + **stress test `-count=5`**.
- **Vượt tầm → sát thương 0 nhưng VẪN PHÁT GÓI.** Không im lặng bỏ qua.
- **Dựng xong dữ liệu rồi mới ghi byte.** Lỗi giữa chừng → bỏ cả gói, không gửi nửa vời.
- **Bất tử 1500 ms sau hồi sinh**; **hồi chiêu nới 50 ms**; tầm cận chiến `100`, tầm xa `300`.
- **Go:** `D:\Go\bin\go.exe`, gofmt `D:\Go\bin\gofmt.exe`. Windows PowerShell (KHÔNG `&&`).
  Test cần CSDL: `$env:NRO_TEST_DSN='postgres://postgres:nrodev123@127.0.0.1:5432/nro_test'`
  (dòng riêng). ⚠️ **Cổng 5432**, không phải 5433 — xem §9 của `BAT-DAU-TU-DAY.md`.
  PostgreSQL phải đang chạy: `D:\PostgreSQL\bin\pg_ctl.exe -D D:\PostgreSQL\data start`.
- **Ghi chú tiếng Việt giải thích *tại sao***, không phải *cái gì*.
- **Kiểm thử đột biến bắt buộc** sau mỗi task có gói byte: cố tình phá code, test **phải
  đỏ**. Không đỏ = test vô dụng, viết lại test. *(Dự án đã 3 lần dính test xanh mà code sai.)*

---

## Thứ tự và lý do

Nhóm **E (chết & hồi sinh)** xếp **sau cùng** vì hai gói `-8` và `-17` **chưa xác minh
được** (đặc tả §7) — task đầu của nhóm E là đọc client Unity để gỡ. Mọi nhóm khác không
phụ thuộc, làm được ngay.

| Task | Nhóm | Phụ thuộc |
|---|---|---|
| 1 Catalog chiêu | A | — |
| 2 Chọn chiêu (`34`) | A | 1 |
| 3 Chốt kiểm trước khi đánh | A | 2 |
| 4 Công thức sát thương + hệ nguyên tố | A | 1 |
| 5 Mô hình quái + gói `-24` | B | — |
| 6 Gói `54` / `-9` / `-12` / `-13` | C | 5 |
| 7 Người đánh quái (nối điều phối) | C | 3, 4, 6 |
| 8 Quái hồi sinh + hồi máu | B | 6 |
| 9 AI quái + gói `-11` / `-10` | D | 5, 7 |
| 10 Gói `-60` **hai dạng** | F | 4 |
| 11 Người đánh người (nối điều phối) | F | 10 |
| 12 **Xác minh client** cho `-8`/`-17` | E | — *(làm được song song từ đầu)* |
| 13 Chết & hồi sinh | E | 12, 7 |
| 14 Test đầu-cuối 2 người chơi | — | tất cả |
| 15 Nghiệm thu tay bằng client Unity | — | 14 |

---

### Task 1: Catalog chiêu `-28/7` — điền 15 chiêu Naruto

**Files:**
- Modify: `server/internal/game/dongbo.go` (hàm `goiCatalogSkill`)
- Modify: `server/internal/game/dongbo_test.go`
- Create: `server/internal/game/chieuthuc.go` (danh mục chiêu, kiểu `MauChieu` / `CapChieu`)
- Create: `server/internal/game/chieuthuc_test.go`

**Interfaces:**
- Produces: `type MauChieu struct { Id byte; Ten string; CapToiDa byte; KieuTonChakra byte; Loai byte; IconId int16; MoTa string; CacCap []CapChieu }`
- Produces: `type CapChieu struct { SkillId int16; Diem byte; NhanLucCan int64; ChakraDung int16; HoiChieuMs int32; Dx, Dy int16; MaxFight byte; SatThuong, Gia int16; ThemThongTin string }`
- Produces: `func ChieuCuaLang(idLang byte) []MauChieu`

- [ ] **Step 1: Viết test đỏ** — dựng gói `-28/7`, so **toàn bộ mảng byte** với giá trị
      ghim cứng cho **một làng có đúng 1 chiêu, 1 cấp** (bảng byte ở đặc tả §3.7).
      Test riêng khẳng định `NhanLucCan` chiếm **8 byte** trong dòng.
- [ ] **Step 2: Chạy — fail.** `D:\Go\bin\go.exe test ./internal/game/ -run TestCatalogChieu -v`
- [ ] **Step 3: Cài đặt** khuôn ghi theo đặc tả §3.7. **`templateId` đánh số từ 1** —
      *tại sao:* bản Java đối xử đặc biệt với `templateId == 0` (nhồi 2 cấp giả 105/106);
      né bằng cách không dùng số 0 thì nhánh đó không bao giờ chạy. Ghi chú rõ trong code.
      Trường `ghiChu` ghi chuỗi rỗng — **không** phát lại `"Ngoc Rong Kame"` của bản cũ.
- [ ] **Step 4: Chạy — PASS.**
- [ ] **Step 5: Điền 15 chiêu** — 3 chiêu × 5 làng theo đặc tả §4.3 (1 cận chiến `type=1`,
      1 tầm xa theo hệ `type=1`, 1 tự thân `type=3`), mỗi chiêu 3 cấp. Test khẳng định
      **mỗi làng đúng 3 chiêu** và **không chiêu nào có `templateId == 0`**.
- [ ] **Step 6: Kiểm thử đột biến** — đổi `NhanLucCan` từ 8 byte sang 4 byte → test phải đỏ.
      Đổi `templateId` của một chiêu thành 0 → test phải đỏ.
- [ ] **Step 7: Commit** — `feat(game): catalog -28/7 phat 15 chieu Naruto (3 chieu x 5 lang)`

---

### Task 2: Gói lên `34` — chọn chiêu, đi qua hộp thư actor

**Files:**
- Modify: `server/internal/game/khuvuc.go` (thêm sự kiện `suKienChonChieu`)
- Modify: `server/internal/game/khuvuc_test.go`

**Interfaces:**
- Produces: `type suKienChonChieu struct { idNhanVat int64; mauChieuId byte }`
- `nguoiTrongKhuVuc` thêm trường `chieuDangChon *CapChieu`

- [ ] **Step 1: Viết test đỏ** — gửi `suKienChonChieu` vào actor, đọc lại (qua một sự
      kiện truy vấn) thấy chiêu đã đổi. Test thứ hai: chọn chiêu **không thuộc làng của
      nhân vật** → không đổi gì.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt.** Bộ điều phối đọc `short templateId` từ gói `34` rồi **gửi sự
      kiện**, KHÔNG ghi thẳng. *Tại sao:* luồng phiên ghi trong khi actor đọc = tranh chấp
      dữ liệu, đúng bẫy §6.9 của `BAT-DAU-TU-DAY.md` mà KH5 đã vấp với `nv.ViTriX/Y`.
- [ ] **Step 4: Chạy — PASS** với `-count=5`.
- [ ] **Step 5: Soi mã** — xác nhận không đường nào ngoài `vongLap` chạm `chieuDangChon`.
- [ ] **Step 6: Commit** — `feat(game): goi 34 chon chieu qua hop thu actor`

---

### Task 3: Chốt kiểm trước khi cho đánh (hàm thuần)

**Files:**
- Create: `server/internal/game/chotkiem.go`
- Create: `server/internal/game/chotkiem_test.go`

**Interfaces:**
- Produces: `func KiemDanhHopLe(nguoiDanh *nguoiTrongKhuVuc, mucTieu ViTriCoMau, bayGio time.Time) KetQuaKiem`
- Produces: `type KetQuaKiem struct { ChoPhep bool; Truot bool; LyDo string }`

Bảng kiểm (đặc tả §2 của spec, đúng thứ tự):

| Kiểm | Kết quả khi hỏng |
|---|---|
| Mục tiêu vừa hồi sinh < 1500 ms | `ChoPhep = false` |
| Mục tiêu đã chết | `ChoPhep = false` |
| Không đủ chakra | `ChoPhep = false` |
| Chiêu chưa hồi (`daTroiQua < hoiChieu − 50ms`) | `ChoPhep = false` |
| **Vượt tầm** | `ChoPhep = true`, **`Truot = true`** |

- [ ] **Step 1: Viết test đỏ** — một test cho **mỗi dòng** trong bảng, giá trị **ghim
      cứng** (không lấy từ chính hằng số đang kiểm — §6.12 của `BAT-DAU-TU-DAY.md`).
      Test biên: `daTroiQua == hoiChieu − 50ms` **đúng bằng** → cho phép.
      Test biên: khoảng cách == đúng tầm → cho phép, không trượt.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt.** Hàm **thuần** (không chạm trạng thái, nhận `bayGio` làm tham
      số) — *tại sao:* để test được mọi mốc thời gian mà không cần `sleep`.
- [ ] **Step 4: Chạy — PASS.**
- [ ] **Step 5: Kiểm thử đột biến** — đổi `− 50ms` thành `+ 50ms` → test biên phải đỏ.
      Đổi "vượt tầm" từ `Truot` sang `!ChoPhep` → test phải đỏ.
- [ ] **Step 6: Commit** — `feat(game): chot kiem truoc khi danh (bat tu/chakra/hoi chieu/tam)`

---

### Task 4: Công thức sát thương + hệ nguyên tố ngũ hành

**Files:**
- Create: `server/internal/game/satthuong.go`
- Create: `server/internal/game/satthuong_test.go`

**Interfaces:**
- Produces: `func TinhSatThuong(goc int64, chieu *CapChieu, heNguoiDanh, heMucTieu byte, chiMang bool) int32`
- Produces: `func HeSoNguyenTo(heNguoiDanh, heMucTieu byte) float64`

Vòng khắc chế (đặc tả §4.2): `Hoả → Phong → Lôi → Thổ → Thuỷ → Hoả`.
Khắc chế ×1,25 · trung tính ×1,00 · bị khắc ×0,80.

- [ ] **Step 1: Viết test đỏ** — bảng 5×5 đầy đủ 25 cặp hệ, **giá trị ghim cứng**.
      Test kẹp trần: đầu vào lớn → ra đúng `2_147_483_647`, **không tràn âm**.
      Test kẹp sàn: không bao giờ ra số âm.
      Test chí mạng nhân đôi.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt.** Tính bằng `int64` rồi mới kẹp xuống `int32` — *tại sao:* tính
      thẳng trên `int32` thì phép nhân tràn **trước khi** kịp kẹp.
      Code **không được biết tên chiêu nào** — mọi thứ đặc thù đọc từ dữ liệu chiêu
      (chuẩn bị cho KH8). Không `switch` theo tên chiêu.
- [ ] **Step 4: Chạy — PASS.**
- [ ] **Step 5: Kiểm thử đột biến** — bỏ bước kẹp trần → test tràn phải đỏ.
      Đảo chiều một cặp khắc chế → test bảng 5×5 phải đỏ.
- [ ] **Step 6: Commit** — `feat(game): cong thuc sat thuong + he nguyen to ngu hanh`

---

### Task 5: Mô hình quái trong actor + phát quái trong `-24`

**Files:**
- Create: `server/internal/game/quai.go`
- Create: `server/internal/game/quai_test.go`
- Modify: `server/internal/game/banDo.go` hoặc nơi dựng `-24` (phần `soQuai`)

**Interfaces:**
- Produces: `type Quai struct { Id byte; MauId byte; Cap byte; CapTinhAnh byte; Hp, HpToiDa int32; X, Y int16; TrangThai byte; ... }`
- Produces: `func (k *KhuVuc) danhSachQuai() []*Quai` *(chỉ gọi từ trong `vongLap`)*

- [ ] **Step 1: Viết test đỏ** — dựng gói `-24` với **2 con quái**, so toàn bộ mảng byte
      với giá trị ghim cứng. Test riêng khẳng định **thứ tự `hp` → `level` → `hpMax`**
      (đặc tả §3.6) — *tại sao có test riêng:* thứ tự này phản trực giác và là chỗ dễ
      chép sai nhất.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt.** Quái là trạng thái **thuộc `KhuVuc`**, tạo lúc `MoKhuVuc`.
      `soQuai ≤ 127` (byte có dấu).
- [ ] **Step 4: Chạy — PASS.**
- [ ] **Step 5: Kiểm thử đột biến** — hoán đổi `level` và `hpMax` → test phải đỏ.
- [ ] **Step 6: Commit** — `feat(game): mo hinh quai trong actor khu vuc + phat trong -24`

---

### Task 6: Dựng gói `54` / `-9` / `-12` / `-13`

**Files:**
- Create: `server/internal/game/goichiendau.go`
- Create: `server/internal/game/goichiendau_test.go`

- [ ] **Step 1: Viết test đỏ** — bốn test byte-chính-xác, mỗi gói một test, mảng byte
      ghim cứng theo đặc tả §3.3. `-12` test với `soDoRoi = 0`.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt.** Trường cuối `-9` ghi hằng `int -1` **kèm chú thích rõ là chưa
      biết ý nghĩa** — đừng giả vờ hiểu (đặc tả §3.3).
- [ ] **Step 4: Chạy — PASS.**
- [ ] **Step 5: Kiểm thử đột biến** — đổi `mobId` từ byte sang short trong `-9` → test phải đỏ.
- [ ] **Step 6: Commit** — `feat(game): goi 54/-9/-12/-13 (byte-dung)`

---

### Task 7: Nối điều phối — người chơi đánh quái

**Files:**
- Modify: bộ điều phối gói (nơi xử `-7`, `-39` của KH5)
- Modify: `server/internal/game/khuvuc.go` (sự kiện `suKienDanhQuai`)
- Modify: `server/internal/game/khuvuc_tichhop_test.go`

- [ ] **Step 1: Viết test đỏ** — client giả gửi `54`, nhận về đúng chuỗi `54` rồi `-9`.
      Test thứ hai: đánh đủ để quái chết → nhận `54` rồi `-12`.
      **Test thứ ba (bắt buộc):** gửi `54` với `mobId = -1` → server đọc đủ **5 byte** và
      gói **kế tiếp trên cùng kết nối vẫn đọc đúng** — *tại sao:* đây là chỗ lệch dòng byte.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt.** Điều phối chỉ đọc gói rồi **gửi sự kiện**; toàn bộ kiểm tra và
      tính toán nằm **trong `vongLap`**. Chốt "quái đầy máu không chết bởi một đòn".
- [ ] **Step 4: Chạy — PASS** với `-count=5`.
- [ ] **Step 5: Kiểm thử đột biến** — xoá hẳn nhánh điều phối `54` → test phải đỏ.
      *(§6.11: đã có tiền lệ xoá cả đoạn điều phối mà 24/24 test vẫn xanh.)*
- [ ] **Step 6: Commit** — `feat(game): dieu phoi 54 danh quai`

---

### Task 8: Quái hồi sinh 5 giây + hồi máu 30 giây

**Files:**
- Modify: `server/internal/game/quai.go`, `khuvuc.go` (móc `nhipTim`)
- Modify: `server/internal/game/quai_test.go`

- [ ] **Step 1: Viết test đỏ** — dùng đồng hồ tiêm được (không `sleep` thật):
      quái chết → sau 4,9 s **chưa** hồi sinh; sau 5,0 s hồi sinh và phát `-13`.
      Quái mất máu → sau 30 s hồi đúng `hpToiDa/10`, không vượt `hpToiDa`.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt** trong `nhipTim` (móc để rỗng từ KH5 — nay dùng tới).
      Đồng hồ là **tham số/trường tiêm được**, không gọi `time.Now()` thẳng — *tại sao:*
      để test được mốc thời gian mà không làm bộ test chậm thêm 35 giây.
- [ ] **Step 4: Chạy — PASS.**
- [ ] **Step 5: Kiểm thử đột biến** — đổi 5 s thành 50 s → test biên phải đỏ.
- [ ] **Step 6: Commit** — `feat(game): quai hoi sinh 5s + hoi mau 30s trong nhip actor`

---

### Task 9: AI quái + gói `-11` / `-10`

**Files:**
- Modify: `server/internal/game/quai.go`
- Modify: `server/internal/game/goichiendau.go`
- Modify: `server/internal/game/quai_test.go`

**Interfaces:**
- Produces: `func (q *Quai) chonMucTieu(nguoi []*nguoiTrongKhuVuc) (*nguoiTrongKhuVuc, time.Duration)`
  — **trả về cả mục tiêu lẫn nhịp đánh**, không sửa trạng thái ngầm.

Hai tầng ưu tiên (đặc tả §5.2):

| Ưu tiên | Chọn ai | Tầm | Nhịp |
|---|---|---|---|
| 1 | Người **đã từng đánh nó**, gần nhất | 300 | 1000 ms |
| 2 | Người bất kỳ gần nhất | 100 | 2000 ms |

- [ ] **Step 1: Viết test đỏ** — test tầng 1 và tầng 2 riêng, kèm nhịp trả về.
      **Test bắt buộc: mục tiêu đã chết KHÔNG được chọn** — *tại sao:* bản Java sai
      thứ tự `&&`/`||` nên big boss nhắm cả xác chết (đặc tả §5.4.3).
      **Test bắt buộc: không có mục tiêu thì nhịp KHÔNG bị tiêu** — quái vừa thấy người
      phải đánh được ngay, không phải chờ thêm nhịp (đặc tả §5.4.2).
      Hai test byte-chính-xác cho `-11` và `-10`.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt.** Danh sách "đã từng đánh nó" nạp khi trúng đòn, **xoá sạch khi
      quái chết**. Ngoặc điều kiện lọc **tường minh**, không dựa vào thứ tự toán tử.
- [ ] **Step 4: Chạy — PASS** với `-count=5`.
- [ ] **Step 5: Kiểm thử đột biến** — bỏ bộ lọc "chưa chết" → test phải đỏ.
      Chuyển việc đặt lại nhịp ra ngoài nhánh "thật sự đánh" → test phải đỏ.
- [ ] **Step 6: Commit** — `feat(game): AI quai 2 tang uu tien + goi -11/-10`

---

### Task 10: Gói `-60` — HAI hàm dựng riêng biệt

**Files:**
- Modify: `server/internal/game/goichiendau.go`
- Modify: `server/internal/game/goichiendau_test.go`

**Interfaces:**
- Produces: `func goiDanhNguoiChoi(...) *protocol.Message` — `docTiep = 1`, **18 byte thân**
- Produces: `func goiHoiMauNguoiChoi(...) *protocol.Message` — `docTiep = 0`, **11 byte thân**

- [ ] **Step 1: Viết test đỏ** — hai test byte-chính-xác riêng, mảng byte ghim cứng.
      Test thứ ba khẳng định **hai thân chênh nhau đúng 7 byte**
      (`loaiChieu` 1 + `satThuong` 4 + `daChet` 1 + `chiMang` 1).
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt.** **Hai hàm tách hẳn**, không một hàm có cờ — *tại sao:* một hàm
      có cờ thì người sửa sau rất dễ ghi thừa trường vào nhánh hồi máu, và triệu chứng là
      client lệch dòng byte vĩnh viễn, cực khó chẩn đoán (đặc tả §2.1).
- [ ] **Step 4: Chạy — PASS.**
- [ ] **Step 5: Kiểm thử đột biến (QUAN TRỌNG NHẤT CỦA CẢ KẾ HOẠCH)** — trong
      `goiHoiMauNguoiChoi` đổi `docTiep` từ 0 thành 1 → test **phải đỏ**. Ghi thêm 1 byte
      vào cuối gói hồi máu → test **phải đỏ**. Không đỏ = test vô dụng, viết lại.
- [ ] **Step 6: Commit** — `feat(game): goi -60 hai dang (don danh 14B / hoi mau 10B)`

---

### Task 11: Nối điều phối — người chơi đánh người chơi

**Files:**
- Modify: bộ điều phối gói, `server/internal/game/khuvuc.go`
- Modify: `server/internal/game/khuvuc_tichhop_test.go`

- [ ] **Step 1: Viết test đỏ** — hai client giả trong cùng khu vực; A gửi `-60` đánh B →
      **cả hai** nhận gói `-60`, và B nhận thêm `-30/9` cập nhật máu.
      Test: đánh người vừa hồi sinh < 1500 ms → không có gói nào phát ra.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt** qua sự kiện vào hộp thư actor.
- [ ] **Step 4: Chạy — PASS** với `-count=5`.
- [ ] **Step 5: Kiểm thử đột biến** — xoá nhánh điều phối `-60` → test phải đỏ.
- [ ] **Step 6: Commit** — `feat(game): dieu phoi -60 danh nguoi choi`

---

### Task 12: ⚠️ XÁC MINH CLIENT cho gói `-8` và `-17`

> **Task này KHÔNG viết code.** Nó gỡ hai chỗ chưa chắc ở đặc tả §7. Làm được **song
> song từ đầu kế hoạch** — không phụ thuộc task nào.

**Files:**
- Modify: `docs/superpowers/specs/2026-07-26-giao-thuc-chien-dau.md` (§3.4 và §7)
- Modify: `docs/superpowers/research/kh6-goi-chien-dau.md` (§6)

**Vấn đề:** bản Java phát `-8` ở hai chỗ với **kiểu trường đầu khác nhau**
(`writeInt` ở `map/Zone.java:726`, `writeShort` ở `services/Service.java:1157`), và phát
`-17` với id nhân vật nhét vào **1 byte** (`services/Service.java:1147`). Không thể cùng
đúng — client là trọng tài duy nhất.

- [ ] **Step 1:** Mở source client Unity ở `D:\NRO_Project`, tìm chỗ đọc lệnh `-8`
      (nhiều khả năng trong `Controller.cs`, cùng file đã dùng ở KH3/KH4).
- [ ] **Step 2:** Ghi lại **kiểu đọc thật** của trường đầu: `readInt` / `readShort` / `readByte`.
- [ ] **Step 3:** Làm y hệt cho `-17`.
- [ ] **Step 4:** Cập nhật đặc tả §3.4 (bỏ dấu ⚠️) và §7 (ghi kết luận + bằng chứng số dòng).
- [ ] **Step 5:** Nếu phát hiện bản Java **có lỗi thật**, ghi thêm một dòng vào
      §6 của `BAT-DAU-TU-DAY.md` — đây là loại phát hiện đắt giá nhất.
- [ ] **Step 6: Commit** — `docs(kh6): xac minh kieu truong id trong goi -8 va -17 tu client`

> **Nếu không mở được client Unity:** dừng lại **hỏi chủ dự án**, đừng đoán.
> §6.14 của `BAT-DAU-TU-DAY.md`: "thấy thứ trông như lỗi nhưng không chắc → dừng lại hỏi".

---

### Task 13: Chết và hồi sinh

**Files:**
- Modify: `server/internal/game/goichiendau.go`, bộ điều phối, `khuvuc.go`
- Modify: `server/internal/game/khuvuc_tichhop_test.go`

**Phụ thuộc: Task 12 phải xong.** Không viết `-8` / `-17` trước khi biết kiểu trường thật.

- [ ] **Step 1: Viết test đỏ** — chuỗi đầy đủ theo đặc tả §6.4: máu về 0 → chính chủ nhận
      `-17`, người khác nhận `-8`. Client gửi `-16` → chính chủ nhận `-16`, cả khu vực
      nhận `-30/15`. Byte-chính-xác cho từng gói.
      Test: người đã chết gửi `54` → **không** đánh được.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Cài đặt.** Đặt mốc `lucHoiSinh` để chốt bất tử 1500 ms (Task 3 dùng tới).
- [ ] **Step 4: Chạy — PASS** với `-count=5`.
- [ ] **Step 5: Kiểm thử đột biến** — bỏ chốt "đã chết thì không đánh" → test phải đỏ.
- [ ] **Step 6: Commit** — `feat(game): chet va hoi sinh nguoi choi (-17/-8/-16/-30/15)`

---

### Task 14: Test đầu-cuối 2 người chơi (mạng + CSDL thật)

**Files:**
- Modify: `server/internal/game/tichhop_test.go`

- [ ] **Step 1: Viết test đỏ** — hai client giả đăng nhập thật, vào cùng khu vực:
      1. A chọn chiêu, đánh quái → **B nhìn thấy** gói `54` và `-9`.
      2. A đánh quái tới chết → B thấy `-12`; sau 5 s cả hai thấy `-13`.
      3. Quái đánh A → A nhận `-11`, B nhận `-10`.
      4. A đánh B tới chết → cả hai thấy `-60`, A thấy `-8` của B.
      5. B hồi sinh → cả hai thấy `-30/15`.
- [ ] **Step 2: Chạy — fail.**
- [ ] **Step 3: Sửa** những gì lộ ra.
- [ ] **Step 4: Chạy — PASS.** Sau đó **toàn bộ bộ test**: `D:\Go\bin\go.exe test ./... -count=5`
      phải xanh, không chập chờn. `gofmt -l` sạch, `go vet` sạch.
- [ ] **Step 5: Commit** — `test(game): dau-cuoi 2 nguoi choi chien dau + quai`

---

### Task 15 (THỦ CÔNG — chủ dự án): nghiệm thu bằng client Unity thật

> Không agent nào làm thay được. Đây là lần đầu tiên game **đánh nhau được**.

- [ ] Mở client Unity thật, đăng nhập, vào làng.
- [ ] Chọn chiêu → đánh quái → **có thấy hoạt ảnh và số sát thương bay lên không?**
- [ ] Đánh quái tới chết → **quái có biến mất rồi 5 giây sau hiện lại không?**
- [ ] Đứng yên gần quái → **quái có tự đánh mình không? Máu có tụt không?**
- [ ] Mở **2 client cùng lúc, cùng làng** → **người kia có thấy mình đánh nhau không?**
- [ ] Để nhân vật chết → **có hiện màn hình chết và hồi sinh được không?**

**Nếu có bất kỳ gói nào sai byte, triệu chứng sẽ là: client đứng hình, thoát app, hoặc
hiện số sát thương vô lý.** Ghi lại triệu chứng chính xác rồi báo — đừng đoán nguyên nhân.

---

## Tự soi (sau khi viết kế hoạch)

**Chỗ dễ hỏng nhất, theo thứ tự:**

1. **Gói `-60` hai dạng** (Task 10) — sai là client lệch dòng byte vĩnh viễn, không có
   thông báo lỗi nào. Đã đặt kiểm thử đột biến riêng, gọi tên là "quan trọng nhất cả kế hoạch".
2. **Gói lên `54` độ dài biến thiên** (Task 7) — KH6 không làm thú cưng nhưng vẫn phải
   đọc đủ 5 byte. Đã có test riêng cho việc "gói kế tiếp vẫn đọc đúng".
3. **Thứ tự `hp` → `level` → `hpMax`** trong `-24` (Task 5) — phản trực giác. Đã có test riêng.
4. **`powRequire` 8 byte** giữa một rừng trường 2/4 byte (Task 1). Đã có test riêng.
5. **Tràn `int32`** khi nhân hệ số (Task 4). Đã bắt buộc tính bằng `int64` rồi mới kẹp.

**Đã cố ý bỏ ra ngoài phạm vi:** rơi đồ (KH7), boss, hiệu ứng trạng thái đầy đủ, chiêu
`type = 4`, thú cưng, PvP có luật. Lý do và đích đến ghi ở đặc tả §0.

**Rủi ro lịch trình:** Task 12 (xác minh client) có thể **chặn** Task 13. Nếu không mở
được source Unity thì nhóm E treo — nhưng 12/15 task còn lại vẫn chạy được, nên không
chặn cả kế hoạch. Đó là lý do xếp nhóm E sau cùng.
