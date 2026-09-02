# KH7 — Vật phẩm, hành trang, cửa hàng NPC, giao dịch (thiết kế)

**Ngày:** 03/09/2026 · **Trạng thái:** thiết kế (chờ duyệt) · **Chỉ tài liệu, chưa có dòng code nào.**

**Điều kiện tiên quyết bắt buộc** ([BO-SUNG-LO-TRINH](../plans/BO-SUNG-LO-TRINH.md) Khối 3):
*"phải chốt căn cước vật phẩm + chống nhân bản TRƯỚC khi viết code"*. Tài liệu này chính là
cái chốt đó. Không được mở task thi công nào của KH7 khi §1 và §2 dưới đây chưa được duyệt.

**Nguồn đối chiếu:** server Java `1-MAY-CHU-GAME/src/` **và** client Unity
`D:\NRO_Project\Assets\Scripts\Game1\` — đọc **cả hai phía** theo luật §6.16 của
`BAT-DAU-TU-DAY.md`. Mọi dòng trong §4 đều có `file:dòng` của cả hai bên. Việc đối chiếu
được làm hai lượt độc lập (một lượt trực tiếp, một lượt bằng agent riêng) rồi hợp nhất —
**13 chỗ lệch** và **11 lỗ hổng giao dịch** tìm được đều ghi nguyên trong §4.6b và §2.2.

---

## 0. Tóm tắt cho chủ dự án — 5 câu

1. **Bản Java không hề đánh số món đồ.** Một món chỉ được nhận diện bằng *"ô thứ mấy trong
   túi"*. Đổi thứ tự ô là đổi danh tính món — đây là gốc của mọi lỗi nhân bản đồ.
2. **Bản Go sẽ cấp cho mỗi món một số căn cước riêng, không bao giờ tái sử dụng**, và mọi
   thao tác đổi chủ đều đi qua một giao dịch cơ sở dữ liệu khoá theo đúng số đó.
3. **Đường lưu hiện tại (`ThoGhi`) TUYỆT ĐỐI không được dùng cho vật phẩm và tiền** — nó ghi
   đè cả dòng bằng ảnh chụp cũ, tức là có thể "hồi sinh" số vàng đã tiêu.
4. **KH7 bản 1 chỉ làm 2 trong 4 thứ**: cửa hàng NPC và kho/hành trang. Giao dịch 2 người và
   ký gửi hoãn sang bản 2 — vì chúng là hai nguồn nhân bản nguy hiểm nhất và cần nền móng ở
   bản 1 đã chạy thật rồi mới an toàn.
5. Đối chiếu byte đã tìm ra **13 chỗ hai bên lệch nhau** — trong đó có một **lỗ hổng lừa
   đảo trực tiếp** (người chơi gỡ đồ khỏi bàn giao dịch, màn hình xoá nhưng máy chủ vẫn
   giữ) và một chỗ bản Java **không hề kiểm người chơi có đủ vàng hay không** trước khi
   trừ. Chi tiết §4.6b và §2.2.

---

## 1. Căn cước vật phẩm

### 1.1 Bản Java làm gì — và vì sao đó là lỗ hổng

**Vì sao phải nói chuyện này trước:** nếu không biết bản gốc sai ở đâu, người thi công sẽ
"chép cho giống" và bê nguyên lỗ hổng sang bản Go.

`Item.java:12-43` — toàn bộ một món đồ trong bản Java:

| Trường | Kiểu | Ghi chú |
|---|---|---|
| `template` | `ItemTemplate` | *loại* đồ (kiếm, giáp…) |
| `quantity` | `int` | số lượng |
| `itemOptions` | `List<ItemOption>` | chỉ số phụ đã roll |
| `createTime` | `long` | thời điểm tạo |
| `info` / `content` | `String` | chuỗi mô tả để hiển thị |

**Không có trường `id`.** (`idtv` ở dòng 22 là id thành viên bang hội, không phải id món.)

Hệ quả: danh tính một món = **chỉ số ô trong mảng**. Bằng chứng thẳng trong code:

- `UseItem.java:140-150` — dùng đồ: `readByte type, readByte where, readByte index` rồi
  `player.inventory.itemsBag.get(index)`. Client gửi **số ô**, server tin luôn.
- `Trade.java:110-112` — đưa đồ vào giao dịch: `itemsBag1.get(index)`.
- `Player.java:878-882` — trang bị: `itemsBody.get(0)` là áo, `get(1)` là quần… **vị trí là ý nghĩa.**

Đây đúng cái bẫy §6 mà `BAT-DAU-TU-DAY.md` gọi là **"chỉ số = danh tính"**. Ba hậu quả:

1. **Sắp xếp lại túi là đổi danh tính mọi món.** `InventoryService.sendItemBag` gọi
   `sortItems()` (`InventoryService.java:562`) **mỗi lần gửi** — nghĩa là số ô của một món
   có thể đổi giữa hai thao tác của người chơi. Một lệnh "dùng ô 3" gửi trước lúc sắp xếp,
   tới server sau lúc sắp xếp, sẽ dùng **nhầm món**.
2. **Không thể truy vết.** Không có gì để trả lời "món này từ đâu ra" — nên khi nhân bản
   xảy ra, không thu hồi chính xác được, chỉ còn cách xoá sạch cả tài khoản.
3. **Không khoá được.** Muốn chặn "hai giao dịch cùng lúc một món" thì phải khoá *một cái
   gì đó*. Không có id thì không có gì để khoá.

### 1.2 Quyết định cho bản Go

| Câu hỏi | Chốt | Vì sao |
|---|---|---|
| Món có id bền toàn cục không? | **CÓ.** Mỗi *món* (instance) một dòng trong bảng `vat_pham`, khoá chính `id BIGSERIAL`. | Không có id thì không khoá được, không truy vết được, không chống nhân bản được (§1.1). |
| `BIGSERIAL` hay `UUID`? | **`BIGSERIAL`** | (a) Chỉ có **một** máy chủ game ghi — không có kịch bản nhiều nút sinh id độc lập, tức là lý do chính để dùng UUID **không tồn tại**. (b) 8 byte thay vì 16, chỉ mục B-tree nhỏ hơn và chèn tuần tự nên không phân mảnh. (c) Đọc bằng mắt được, tiện hỗ trợ người chơi và điều tra gian lận. (d) Dự án **đã dùng `BIGSERIAL`** cho `tai_khoan` và `nhan_vat` (`0001`, `0002`) — thêm một kiểu id thứ hai là thêm một thứ phải nhớ mà không được gì. **Đánh đổi chấp nhận:** id đoán được (1, 2, 3…). Không sao — id **không bao giờ đi xuống client** (§4.1), nó chỉ sống trong máy chủ và CSDL. |
| Id có tái sử dụng không? | **KHÔNG BAO GIỜ.** Món bị tiêu/huỷ thì dòng **không xoá**, chỉ đổi `noi_chua = 'DA_HUY'`. | Tái dùng id làm nhật ký đổi chủ nói dối: "món 42" trong nhật ký năm ngoái thành một món khác hẳn hôm nay. |
| Chỉ số phụ (option) lưu ở đâu? | **`JSONB` một cột trên chính dòng món.** | Xem bảng đánh đổi ngay dưới. |

**Option: bảng con hay JSONB — cân nhắc thật, không chọn bừa**

| | Bảng con `vat_pham_option` | **JSONB `chi_so_phu` (CHỌN)** |
|---|---|---|
| Đọc một món | 2 truy vấn (hoặc JOIN) | 1 truy vấn, không JOIN |
| Đọc cả túi 30 món | JOIN + gom nhóm ở tầng Go | đọc thẳng 30 dòng |
| Ghi khi đổi chủ | phải đổi khoá ngoại của N dòng con | **1 lệnh `UPDATE` một dòng** — quan trọng nhất: giữ được bất biến §2.3 |
| Truy vấn "món nào có option 47 > 100" | dễ (SQL thường) | cần chỉ mục GIN, cú pháp lạ hơn |
| Ràng buộc kiểu chặt | CSDL ép được | phải ép ở tầng Go + `CHECK` jsonb |

Chọn JSONB vì **đơn vị nguyên tử của KH7 là "một món"**, không phải "một option". Mọi thao
tác nguy hiểm (đổi chủ, bán, giao dịch) đều tác động lên trọn món; để option ở bảng riêng
là tự tay chẻ một hành động nguyên tử thành nhiều dòng ở nhiều bảng — đúng thứ đẻ ra trạng
thái nửa vời khi đứt giữa chừng. Truy vấn phân tích option là việc của panel/admin, hiếm và
không nhạy cảm thời gian; khi cần thì thêm chỉ mục GIN.

Khuôn JSONB chốt cứng (ép ở tầng Go, có `CHECK` chống rác):

```json
{"opt": [{"id": 47, "param": 120}, {"id": 6, "param": 30}]}
```

`id` ∈ [0, 255] (client đọc 1 byte, §4.5) · `param` ∈ [0, 65535] (client đọc **unsigned
short**, §4.5 — **ghi số âm là hỏng**).

### 1.3 Món cộng dồn (stackable) — chỗ dễ sai

Đồ tiêu hao (thuốc hồi máu…) không đáng một dòng cho mỗi cái. **Luật chốt:**

- Món **có option** (`chi_so_phu` khác rỗng) → **luôn một dòng một cái**, `so_luong = 1`.
  Vì hai món cùng loại nhưng option khác nhau **không phải cùng một thứ**.
- Món **không option** và mẫu khai `cong_don = true` → một dòng mang `so_luong > 1`.
- **Ranh giới không được nhập nhèm:** một mẫu là cộng dồn hay không do **mẫu** quyết định,
  không do món cụ thể. Nếu một mẫu cộng dồn bỗng được gắn option (do chế tác sau này) thì
  phải **tách** thành dòng riêng — quy tắc này phải có test.

---

## 2. Chống nhân bản (dupe)

### 2.1 Bất biến — thứ phải luôn đúng, mọi lúc

Ba câu này là **hợp đồng của cả KH7**. Mỗi câu phải có ít nhất một test và một đột biến:

| # | Bất biến | Nói bằng lời thường |
|---|---|---|
| **BB1** | Với mọi `id`: đúng **một** dòng trong `vat_pham`, và dòng đó có đúng **một** chỗ chứa (`chu_so_huu_id` + `noi_chua` + `o_so`). | Một món chỉ nằm ở **một** chỗ, không bao giờ ở hai nơi cùng lúc. |
| **BB2** | Tổng `so_luong` của một `id` không tự tăng. Chỉ tăng khi có dòng `vat_pham_lich_su` kiểu `SINH_RA` tương ứng. | Đồ không tự đẻ. Mọi món mới đều có giấy khai sinh. |
| **BB3** | Mọi thay đổi chủ sở hữu đều để lại **đúng một** dòng `vat_pham_lich_su`, ghi trong **cùng** giao dịch CSDL. | Không có lần đổi chủ nào không được ghi sổ. |

Thêm một bất biến về tiền, vì tiền cũng bị nhân bản y như đồ:

| **BB4** | `vang`/`ngoc` chỉ đổi qua đường `UPDATE ... SET vang = vang + $delta WHERE id = $1 AND vang + $delta >= 0`, **không bao giờ** bằng cách ghi một giá trị tuyệt đối đã đọc trước đó. | Cộng/trừ tại chỗ, không "đọc rồi ghi đè" — vì giữa đọc và ghi số dư có thể đã đổi. |

### 2.2 Mọi đường nhân bản — liệt kê đầy đủ, mỗi đường một chốt chặn

| # | Đường nhân bản | Chuyện gì xảy ra | Chốt chặn |
|---|---|---|---|
| **D1** | **Giao dịch 2 người đứt giữa chừng** | Bản Java: `Trade.java:365-398` — `itemsBag1`/`itemsBag2` là **bản sao** túi lúc mở giao dịch; lúc chốt thì `player1.inventory.itemsBag = itemsBag1` **ghi đè cả túi**. Nếu trong lúc giao dịch mở, người chơi bán/dùng một món ở túi thật (tab khác), lúc chốt bản sao cũ **hồi sinh** món đó — nhân bản, và người chơi đã cầm luôn tiền bán. | **Một giao dịch CSDL duy nhất** cho cả hai chiều. Không có "bản sao túi". Khoá `SELECT ... FOR UPDATE` theo **danh sách `id` món** (sắp xếp tăng dần, xem D5) + `id` hai nhân vật. Kiểm lại **tại thời điểm chốt** rằng mỗi món vẫn `chu_so_huu_id` đúng người và `noi_chua='TUI'`; lệch một món → **huỷ cả giao dịch**. |
| **D2** | **Rớt mạng lúc bán** | Đã trừ đồ, chưa cộng tiền (hoặc ngược lại) → mất đồ hoặc được tiền không. | Trừ đồ và cộng tiền nằm trong **cùng một `BEGIN…COMMIT`**. Postgres đảm bảo hoặc cả hai, hoặc không gì cả. Máy chủ **không** báo cho client trước khi `COMMIT` trả về thành công ("**ghi sổ trước — giao sau**", §2.4). |
| **D3** | **Mua trùng gói** (bấm 2 lần / client gửi lại / mạng lặp gói) | Trừ tiền hai lần, hoặc tệ hơn: một lần trừ thất bại nhưng đồ cộng hai lần. | **Idempotency key**: mỗi hành động kinh tế mang một khoá `(nhan_vat_id, loai_hanh_dong, so_thu_tu)` với `so_thu_tu` là bộ đếm tăng dần **do máy chủ giữ trong phiên**. Bảng `vat_pham_giao_dich` có `UNIQUE(khoa_idem)`; gói trùng đụng UNIQUE → nhận diện là **lặp**, trả lại **đúng kết quả lần đầu**, không làm lại. |
| **D4** | **Đăng nhập trùng** (hai phiên cùng tài khoản) | Hai phiên cùng cầm một bản sao trạng thái trong bộ nhớ → cùng tiêu một món. | Đã có nền: `SoDangNhap` (`sodangnhap.go`) — **một tài khoản một phiên**, đá phiên cũ, chờ nó lưu xong rồi phiên mới mới đọc. KH7 **thêm** một lớp nữa vì lớp trên là bộ nhớ: mọi `UPDATE` món đều có `WHERE chu_so_huu_id = $1 AND noi_chua = $2` — phiên "ma" đọc dữ liệu cũ sẽ **không khớp** và ghi 0 dòng, tức là thất bại chứ không âm thầm nhân bản. |
| **D5** | **Hai giao dịch cùng lúc chạm cùng một món** (A giao dịch với B, đồng thời bán món đó cho NPC) | Cả hai cùng đọc "món còn đó", cả hai cùng giao → một món thành hai. | `SELECT ... FOR UPDATE` theo `id` món. **Khoá theo thứ tự `id` tăng dần** ở *mọi* nơi — nếu không thì hai giao dịch khoá ngược thứ tự sẽ **deadlock** (Postgres tự huỷ một bên, người chơi thấy lỗi ngẫu nhiên không tái hiện được). |
| **D6** | **Actor và luồng phiên cùng ghi** | Đúng bẫy §6.9: chỉ luồng `vongLap` của khu vực được chạm trạng thái khu vực. Nếu handler mua hàng (chạy trên goroutine phiên) sửa `nv.Vang` trong khi actor cũng đang sửa → tranh chấp dữ liệu thật. | **Kỷ luật actor §6.9/§6.10b:** mọi thay đổi vật phẩm/tiền do người chơi khởi xướng đi qua **một sự kiện gửi vào actor** hoặc — tốt hơn cho KH7 — chạy trên goroutine phiên **nhưng chỉ chạm CSDL**, rồi *sau khi `COMMIT` xong* mới báo actor cập nhật bản trong bộ nhớ. **Tuyệt đối không chạm CSDL từ `vongLap`** (§6.10b: actor không bao giờ chờ ổ đĩa). |
| **D7** | **`ThoGhi` ghi đè cả dòng** ⚠️ **nguy hiểm nhất, đã tồn tại sẵn trong code hôm nay** | `thogi.go:120-129` xếp một ảnh chụp `store.NhanVat` (gồm `Vang`, `Ngoc`) và `store.go:194-206` `UPDATE ... SET vang = $13, ngoc = $14` — **giá trị tuyệt đối**. Kịch bản chết người: actor chụp ảnh lúc 12:00:00 (vàng = 1000) → người chơi mua đồ hết 400 (CSDL còn 600) → thợ ghi xả hàng đợi lúc 12:00:03 và ghi đè **1000**. Người chơi **có cả đồ lẫn vàng**. Chính `khuvuc.go:1742-1744` đã ghi nhận đúng mối nguy này bằng lời. | **Tách tiền ra khỏi đường `ThoGhi`.** Xem §5.3 — đây là một task riêng, làm **trước** mọi task mua bán. |
| **D8** | **Nhặt đồ dưới đất hai lần** | Hai người bấm nhặt cùng lúc, hoặc một người gửi gói `-20` hai lần. | Món dưới đất là **một dòng `vat_pham`** với `noi_chua='DUOI_DAT'`. Nhặt = `UPDATE ... SET chu_so_huu_id=$1, noi_chua='TUI' WHERE id=$2 AND noi_chua='DUOI_DAT'`. Người thứ hai nhận `RowsAffected() == 0` → thua cuộc, báo "món đã có người nhặt". **Không dùng kiểm-rồi-ghi hai bước.** |
| **D9** | **Admin/panel web ghi thẳng** | Web dùng chung CSDL. Panel "tặng đồ" viết cẩu thả có thể chèn dòng không có lịch sử → thủng BB2. | Cấp đồ **chỉ** qua một hàm duy nhất trong Go (`CapVatPham`) ghi cả `vat_pham` lẫn `vat_pham_lich_su` trong một giao dịch. Panel gọi hàm đó (qua hàng đợi lệnh), **không** `INSERT` thẳng. Ràng buộc CSDL đỡ phần còn lại (§5.2). |
| **D10** | **Không kiểm số dư trước khi trừ** ⚠️ *(lỗ hổng thật của bản Java — T1)* | `Trade.java:369-396`: chỉ kiểm **trần trên** (`gold + goldTrade > LIMIT_GOLD`), rồi trừ thẳng `player1.inventory.gold -= goldTrade1`. **Không có một dòng nào kiểm "có đủ vàng không".** Người chơi hứa 500k vàng khi trong túi có 0 → vàng **âm**, hoặc tuỳ kiểu số mà thành **in vàng**. | BB4: điều kiện `vang + $delta >= 0` nằm **trong chính câu `UPDATE`** (§2.1), không phải một `if` ở tầng Go. Số dư không đủ → `RowsAffected() == 0` → cả giao dịch `ROLLBACK`. Không có đường nào lách. |
| **D11** | **Bộ đếm đồng ý dùng chung cho hai người** *(T2)* | `Trade.java:344-349`: `this.accept++`, đủ `2` là chốt. **Một người gửi gói chốt hai lần là giao dịch chạy** dù người kia chưa bấm gì. | Trạng thái đồng ý phải là **hai cờ riêng theo id nhân vật** (`daDongY[idA]`, `daDongY[idB]`), không phải một con số. Test đột biến: đổi lại thành bộ đếm → phải đỏ. |
| **D12** | **Thêm/đổi đồ sau khi đã khoá bàn** *(T4)* | `TransactionService.java:129-162`: không có cờ "đã khoá"; lệnh thêm đồ vẫn được nhận sau lệnh khoá. Đối phương đang nhìn danh sách **cũ** → tráo đồ ngay trước lúc chốt. | Phiên giao dịch là một **máy trạng thái** rõ ràng: `MO → DANG_BAY → DA_KHOA → DA_CHOT`. Lệnh thêm/gỡ đồ chỉ hợp lệ ở `DANG_BAY`; tới `DA_KHOA` thì **mọi thay đổi làm mất khoá của cả hai bên** và phải bày lại. |
| **D13** | **Client gỡ đồ mà máy chủ không biết** *(T5 / L6)* | Client có gửi lệnh gỡ (`Service.cs:180-184`, `Panel.cs:13304`) nhưng máy chủ **không có nhánh xử lý** (`TransactionService.java:77-177` chỉ có 0,1,2,3,5,7). Màn hình xoá món, máy chủ vẫn giữ trong danh sách → **chốt xong là mất đồ**. Đây là lỗ hổng lừa đảo trực tiếp, không phải lỗi hiển thị. | Bản 2 **bắt buộc** cài nhánh gỡ đồ, và **danh sách trên bàn phải do máy chủ gửi lại sau mỗi thay đổi** — client không bao giờ được là nguồn sự thật về "đang bày gì". |
| **D14** | **Mua trùng ở chợ ký gửi** *(K4)* | `ConsignShopService.java:156/189`: đọc cờ `isBuy` rồi mới đặt `isBuy = true` — **hai bước không nguyên tử**, không khoá. Hai người bấm mua cùng lúc → **món nhân đôi**. | Bản 3: chính là D8 nhưng cho `noi_chua='KY_GUI'` — `UPDATE ... WHERE id=$1 AND noi_chua='KY_GUI'` một bước, người thứ hai nhận 0 dòng. |
| **D15** | **Id ký gửi tái sử dụng** *(K5)* | `ConsignShopService.java:395-405`: id mới = `getMaxId() + 1`. Xoá món có id lớn nhất thì id đó **được cấp lại** → lệnh "nhận tiền món 37" tra nhầm món. | Đúng lý do §1.2 chốt **không bao giờ tái sử dụng id**. Bản 3 dùng thẳng `vat_pham.id`, không có bộ đếm riêng. |
| **D16** | **Trừ phí trước khi kiểm dữ liệu** *(K3)* | `ConsignShopService.java:439-447`: trừ 2 thỏi vàng **trước** khi kiểm giá và số lượng hợp lệ. Gửi gói rác là mất phí. | §2.3 bước 4: kiểm **đầy đủ** trên dữ liệu đã khoá **rồi mới** động vào bất cứ thứ gì. Phí là một thay đổi, phải nằm sau bước kiểm. |

### 2.3 "Ghi sổ trước — giao sau"

**Vì sao:** thứ tự sai là mất đồ hoặc đẻ đồ, và cả hai đều im lặng.

Thứ tự bắt buộc cho **mọi** hành động kinh tế:

```
1. Kiểm sơ bộ ở bộ nhớ (nhanh, để từ chối sớm — KHÔNG phải chốt chặn)
2. BEGIN
3.   SELECT ... FOR UPDATE  (theo id tăng dần)
4.   Kiểm lại ĐẦY ĐỦ trên dữ liệu vừa khoá  ← đây mới là chốt chặn thật
5.   UPDATE/INSERT vat_pham + UPDATE tiền (cộng/trừ tại chỗ) + INSERT vat_pham_lich_su
6. COMMIT
7. CHỈ SAU KHI commit trả về nil: cập nhật bản trong bộ nhớ + gửi gói cho client
```

Bước 1 và bước 4 **cố ý trùng nhau**. Bước 1 để không phải mở giao dịch CSDL cho mọi cú
bấm bậy; bước 4 là cái duy nhất được tin. Ai bỏ bước 4 vì "đã kiểm ở trên rồi" là mở lại
toàn bộ D1/D5.

### 2.4 Cái KH7 **không** chống được — nói thẳng

- **Bơm đồ do lỗi cân bằng** (quái rơi đồ quá tay, giá NPC lệch) — đó là bài toán số liệu,
  thuộc KH8. Căn cước không cứu được.
- **Người chơi bị lừa** (đưa đồ rồi bên kia không đưa lại). Hệ thống làm đúng, người chơi
  vẫn mất đồ. Chỉ có nhật ký đổi chủ giúp GM xử lý sau — thuộc Khối 5.
- **Nhân bản qua đường phục hồi sao lưu** — khôi phục CSDL về mốc cũ sẽ hồi sinh đồ đã tiêu.
  Thuộc quy trình vận hành, không phải code.

---

## 3. Tách rõ bốn thứ — cái nào làm, cái nào hoãn

`BO-SUNG` Khối 3 cảnh báo đúng: cái tên "KH7" che mất **bốn cơ chế khác hẳn nhau**.

| Cơ chế | Bản Java | Bản 1 (KH7-1) | Vì sao |
|---|---|---|---|
| **Kho / hành trang** | `InventoryService.java`, `Inventory.java` | ✅ **LÀM** | Là nền móng. Không có túi thì không có gì để mua, bán, trao. Gói `-36`/`-37`/`-35`. |
| **Cửa hàng NPC** | `ShopService.java` (`-44`, lệnh lên `6`/`7`) | ✅ **LÀM** | Rủi ro **thấp nhất**: chỉ một người chơi, đối tác là máy chủ nên không có "bên kia rớt mạng". Là nơi thử toàn bộ kỷ luật §2.3 trong điều kiện dễ. |
| **Giao dịch 2 người** | `Trade.java` 508 dòng | ⛔ **HOÃN sang bản 2** | Rủi ro **cao nhất** (D1, D5 cùng lúc, hai phiên, hai actor). Làm nó khi §2.3 chưa chạy thật là chắc chắn thủng. Nó cũng cần thêm một thứ bản 1 chưa có: **trạng thái phiên giao dịch sống giữa hai người**. |
| **Ký gửi / chợ đồ** | `ConsignShopService.java` 598 dòng (`-100`) | ⛔ **HOÃN sang bản 3** | Là **chợ bất đồng bộ**: món rời khỏi chủ nhưng chưa có chủ mới, sống trong CSDL nhiều ngày, có tiền treo. Cần thêm `noi_chua='KY_GUI'`, đường hoàn trả khi huỷ, và chống "mua trùng" giữa nhiều người mua — mỗi thứ là một nguồn nhân bản mới. Chỉ an toàn khi bản 1 + bản 2 đã chạy thật. |

**Ranh giới cố tình không nhập nhèm:** ký gửi **không phải** cửa hàng NPC dù client dùng gói
gần giống. Cửa hàng NPC bán hàng **do máy chủ sinh ra** (nguồn vô hạn, không có bên bán thật).
Ký gửi chuyển **món có căn cước** giữa hai người thật. Gộp hai thứ vào một luồng code là lỗi
thiết kế, không phải tiết kiệm.

**Ngoài phạm vi KH7 hoàn toàn:** chế tác/cường hoá (`CombineService.java` — `BO-SUNG` Khối 9
yêu cầu ghi rõ: **thuộc KH8**, vì nó là *số liệu* nhiều hơn *luồng*), hòm thư, rương gỗ,
thú cưng mang đồ, đồ theo bang hội.

---

## 4. Bảng byte — đối chiếu **cả hai phía**

**Luật §6.16:** đọc server Java thôi là chưa đủ; client là trọng tài duy nhất. Mọi dòng dưới
đây có `file:dòng` cả hai bên. Client build đích: `D:\NRO_Project`.

### 4.0 Hai cờ client quyết định kiểu đọc — đọc trước, nếu không sai hết bảng

| Cờ | Giá trị thật | File | Ảnh hưởng |
|---|---|---|---|
| `ModFunc.isReadInt` | **`true`** | `ModFunc.cs:45` | `ReadParamDQT()` đọc **unsigned short (2 byte)** (`myReader.cs:222-229`) |
| `ModFunc.isShortOptionTemp` | **`false`** | `ModFunc.cs:300` | `readShortOptionTemp()` đọc **unsigned byte (1 byte)** (`myReader.cs:213-220`) |

→ **Khối option chốt cứng: `byte optionId` + `unsigned short param`.** Trùng khớp bản Java
(`InventoryService.java:592-593` `writeByte(id)` + `writeShort(param)`). ⚠️ `param` phía
client là **unsigned** → miền hợp lệ **0…65535**; ghi số âm sẽ ra số khổng lồ trên màn hình
người chơi, **không có lỗi nào báo** (đúng §6.18).

### 4.1 Nguyên tắc số một của KH7 về giao thức

> **Số căn cước (`vat_pham.id`) KHÔNG BAO GIỜ đi xuống client.**

Client cũ không có ô nào chứa nó — mọi gói đều chỉ mang `template.id` (`short`) và định vị
món bằng **thứ tự trong gói**. Bản Go vì thế giữ một **bảng phiên dịch mỗi phiên**:
`o_so → vat_pham.id`, dựng lại mỗi lần gửi hành trang. Client gửi "ô 3", máy chủ tra ra
món thật rồi mới khoá theo `id`. **Bảng này phải được dựng lại cùng lúc với gói gửi đi,
trong cùng một hàm** — nếu tách ra, một lần sắp xếp xen vào giữa là đủ để dùng nhầm món
(chính là lỗi §1.1 điểm 1 của bản Java).

### 4.2 Gói XUỐNG `-36` — hành trang (túi)

Java ghi: `InventoryService.java:561-602` · Client đọc: `Controller.cs:2202-2264`

| # | Trường | Java ghi | Client đọc | Khớp? |
|---|---|---|---|---|
| 1 | `cAction` | `writeByte(0)` (:566) | `readByte()` (:2204) | ✅ |
| 2 | `soMon` | `writeByte(size)` (:567) | `readUnsignedByte()` (:2208) | ⚠️ Java `writeByte` có dấu → **≤127**; client đọc unsigned ≤255. Bản Go: **kẹp 127** cho an toàn |
| | *lặp `soMon` lần:* | | | |
| 3 | `templateId` | `writeShort` (:573) | `readShort()`, **`== -1` thì `continue`** (:2214-2218) | ⚠️ **xem §4.6 lỗi L1** |
| 4 | `soLuong` | `writeInt` (:574) | `readInt()` (:2221) | ✅ |
| 5 | `info` | `writeUTF` (:575) | `readUTF()` (:2222) | ✅ |
| 6 | `content` | `writeUTF` (:576) | `readUTF()` (:2223) | ✅ |
| 7 | `soOption` | `writeByte` (:577) | `readUnsignedByte()` (:2225) | ✅ |
| | *lặp `soOption` lần:* | | | |
| 8 | `optionId` | `writeByte` (:592) | `readShortOptionTemp()` = **1 byte** (§4.0) | ✅ |
| 9 | `param` | `writeShort` (:593) | `ReadParamDQT()` = **unsigned short** (§4.0) | ✅ (miền 0…65535) |

Lệnh con khác của `-36`: `cAction = 2` — đổi số lượng một ô:
`byte oSo, int soLuongMoi` (client `Controller.cs:2248-2262`). Client tự đặt ô về `null` khi
số lượng về 0. Bản Go **dùng gói này** cho "vừa dùng hết một món" thay vì gửi lại cả túi.

### 4.3 Gói XUỐNG `-37` — đồ đang mặc

Java: `InventoryService.java:604-646` · Client: `Controller.cs:2005+`

| # | Trường | Java | Ghi chú |
|---|---|---|---|
| 1 | `cAction` | `writeByte(0)` (:608) | client `!= 0` thì `break` (:2007) |
| 2 | `head` | `writeShort` (:609) | |
| 3 | `soO` | `writeByte(size)` (:610) | |
| | *mỗi ô:* ô rỗng → `writeShort(-1)` **và dừng ô đó** (:613); ô có đồ → như §4.2 trường 3-9 (:615-637) | | ✅ **Đúng** — đây là khuôn chuẩn mà `-36` **lẽ ra phải theo** |

**Vị trí ô mang ý nghĩa cố định** (`Player.java:878-882`): 0 áo · 1 quần · 2 găng · 3 giày ·
4 nhẫn. Bản Go giữ nguyên bố cục này (client vẽ theo nó) nhưng **không** dùng nó làm danh
tính — danh tính vẫn là `vat_pham.id` (§4.1).

### 4.4 Gói XUỐNG `-35` — rương

Java: `InventoryService.java:648-697` · Client: `Controller.cs:2265+`
Cùng khuôn `-37` (có ghi `-1` cho ô rỗng, :655). `cAction=1` = "mở giao diện rương", thân rỗng
(:688-697). Bản 1 **có** rương (là kho), nhưng **không** có "rương gỗ" và các loại hộp đặc biệt.

### 4.5 Gói XUỐNG `-44` — mở cửa hàng NPC

Java: `ShopService.java:165-218` (`openShopType0`, `NORMAL_SHOP = 0`) ·
Client: `Controller.cs:2572-2685`

| # | Trường | Java | Client | Khớp? |
|---|---|---|---|---|
| 1 | `loaiShop` | `writeByte(NORMAL_SHOP)` (:172) | `readByte()` (:2579) | ✅ |
| 2 | `soTab` | `writeByte` (:173) | `readUnsignedByte()` (:2580) | ✅ (kẹp 127) |
| | *mỗi tab:* | | | |
| 3 | `tenTab` | `writeUTF` (:175) | `readUTF()`, **tách bằng `\n` thành tối đa 2 dòng** (:2598) | ⚠️ dùng `\n` để có tab hai dòng; đừng vô tình nhét `\n` vào tên |
| 4 | `soMon` | `writeByte` (:176) | `readUnsignedByte()` (:2612) | ✅ |
| | *mỗi món:* | | | |
| 5 | `templateId` | `writeShort` (:178) | `readShort()`, `-1` → `continue` (:2621-2625) | ✅ |
| 6 | `giaVang` | `writeInt(cost)` khi `COST_GOLD` (:180) | `buyCoin = readInt()` (:2640) | ✅ |
| 7 | `giaNgoc` | `writeInt(cost)` khi `COST_GEM` (:184) | `buyGold = readInt()` (:2641) | ✅ **Chú ý tên gây lẫn:** ô thứ hai client gọi là `buyGold` nhưng thực chất là **ngọc**. Bản Go đặt tên `giaNgoc`, không chép tên client |
| 8 | `soOption` | `writeByte` (:192) | `readUnsignedByte()` (:2660) | ✅ |
| 9 | *option* | `byte id` + `short param` (:194-195) | §4.0 | ✅ |
| 10 | `laMoi` | `writeByte(0/1)` (:197) | `readByte()` (:2675) | ✅ |
| 11 | `coPart` | `writeByte(1)` nếu `type==5`, ngược lại `0` (:198-205) | `readByte()==1` (:2677) | ✅ |
| 12 | *(nếu `coPart`)* | `short head, body, leg, -1` (:200-203) | `readShort()` ×4 (:2679-2682) | ✅ |

⚠️ **Chỉ làm `loaiShop = 0`.** Các loại 1/2/3/4/8 (`Controller.cs:2631-2658`) đọc **thân
gói khác hẳn** — loại 2 còn chèn thêm `maxPage` sau tên tab (:2601). Đây đúng bẫy §6.6
"một mã lệnh, nhiều dạng độ dài". Bản Go: **mỗi loại một hàm dựng riêng, không dùng cờ.**

### 4.6 Gói XUỐNG `-20` — nhặt đồ

Java: `Zone.java:463-500` · Client: `Controller.cs:3395-3430`

| # | Trường | Java ghi | Client đọc | Khớp? |
|---|---|---|---|---|
| 1 | `itemMapId` | `writeShort` (:464) | `readShort()` (:3399) | ✅ |
| 2 | `loiNhan` | `writeUTF` (:468/477/486) | `readUTF()` (:3408) | ✅ |
| 3 | `soLuong` | `writeShort` (:498) — **đúng một lần** | `readShort()` (:3412), **rồi nếu `type==9` hoặc `10` đọc THÊM một `readShort()` nữa** (:3412-3423) | ❌ **LỖI L2 — xem dưới** |

### 4.6b Hai lỗi thật của bản Java — **đừng chép**

| Mã | Gói | Lỗi | Triệu chứng người chơi |
|---|---|---|---|
| **L1** | `-36` | Java khai `soMon = size` (`InventoryService.java:567`) rồi **`continue` bỏ qua ô rỗng mà KHÔNG ghi `-1`** (:570-572). Client thì luôn đọc đủ `soMon` phần tử và chỉ bỏ qua khi gặp `-1` (`Controller.cs:2214-2218`). Túi có ô trống → **ghi thiếu byte**. | Client đọc quá phần thân → ném lỗi → bị `catch` rỗng nuốt (§6.18). **Hành trang không hiện, không lỗi, không log.** Cùng file `-37` và `-35` lại làm **đúng** (ghi `-1`) — bằng chứng đây là lỗi sót, không phải cố ý. **Bản Go: luôn ghi `-1` cho ô rỗng.** |
| **L2** | `-20` | Java ghi **một** `short soLuong` (`Zone.java:498`). Client với `template.type == 9` (vàng) hoặc `10` (ngọc) đọc **hai** short liên tiếp (`Controller.cs:3412-3423`). | Nhặt vàng/ngọc: client đọc hụt → **số tiền không cộng**, im lặng hoàn toàn. **Bản Go: với type 9/10 ghi ĐỦ hai short** (short thứ hai là số tiền thật). |

**Các chỗ lệch còn lại — ghi đủ để không phải giải mã lại**

| Mã | Chỗ | Lệch gì | Bản Go làm gì |
|---|---|---|---|
| **L3** | Mọi `param` option | Java `writeShort` (có dấu), client `readUnsignedShort` | Kẹp **0…65535**; số cần lớn hơn thì phải đổi *đơn vị*, không phải tràn số. Chính vì giới hạn này mà bản Java phải bịa mẹo đổi option 213→222/223 và chia 1000/1000000 (`InventoryService.java:579-590`) — **bản Go không chép mẹo đó**, thiết kế đơn vị cho đúng ngay từ đầu |
| **L4** | `-20` số lượng | Java chốt cứng `9999` khi `> 32767` (`Zone.java:498`) | Nhặt chồng lớn hiện sai số. Bản Go: chia thành nhiều lần nhặt, hoặc dùng gói cập nhật số lượng riêng |
| **L5** | `-100` lật trang ký gửi | Máy chủ Java gửi (`ConsignShopService.java:233-281`), **client không có nhánh nào đọc** | Tính năng phân trang chợ **đã chết sẵn trong bản gốc**. Bản 3 phải thiết kế phân trang bằng gói client thật sự có đọc, đừng chép |
| **L6** | `-86` gỡ đồ | Client gửi, máy chủ không xử lý | = **D13**, lỗ hổng lừa đảo. Bản 2 bắt buộc cài |
| **L7** | `-35` lệnh con 1 | Java gửi đúng 1 byte rồi hết (`InventoryService.java:688-697`); client cố đọc thêm một byte `isBoxClan` trong `try/catch` rỗng | Cờ rương bang **luôn** sai. Bản 1 không có rương bang nên vô hại, nhưng phải **ghi đủ byte** nếu sau này dùng |
| **L8** | `-36` lệnh con 2 | Client có đọc, máy chủ Java **không bao giờ gửi** | Đây là **cơ hội**: bản Go **dùng** lệnh con này để cập nhật một ô thay vì gửi lại cả túi (§4.2) |
| **L9** | Nhánh theo `version` | Java ghi `byte` hay `int` tuỳ `session.version` (`ConsignShopService.java:256-260`; `Service.java:1391-1394`); client **luôn** đọc kiểu lớn | Lệch 3 byte, hỏng phần còn lại của gói. **Bản Go không có khái niệm version client** — luôn ghi kiểu client đọc. Đây là một dấu vết phải xoá, không phải tính năng |
| **L10** | `-44` loại tiền | Java nhét cả ngọc, hồng ngọc và phiếu vào **cùng một ô** (`ShopService.java:185-191`); client chỉ có hai ô | Bản 1 **chỉ hai loại tiền**: vàng và ngọc. Thêm loại thứ ba là phải đổi giao thức, không nhét chung |
| **L11** | `ItemOption` phía client tự đổi id | `ItemOption.cs:29-38` viết lại 22→6, 23→7 và nhân param ×1000 | **Tránh hẳn dải id 22/23** trong thiết kế option Naruto. Dùng vào là số hiển thị khác số gửi, không có gì báo |
| **L12** | Số bản mẫu option | Client cấp phát mảng bằng `readUnsignedByte` (≤256, `Controller.cs:4632-4637`); option id vượt số mẫu → `IndexOutOfRange` → nuốt | **Trần cứng 255 mẫu option**, và mọi `optionId` gửi đi phải có mẫu tương ứng trong `-28/8` subtype 0 |
| **L13** | Nuốt lỗi ở **cả hai đầu** | Máy chủ Java cũng bọc `try {} catch {}` rỗng khi ghi gói (`InventoryService.java:600-601, 643-644`) | Đúng ràng buộc §1.10 của spec chiến đấu: **dựng xong toàn bộ dữ liệu rồi mới ghi byte**; lỗi giữa chừng thì **bỏ cả gói và log**, không gửi gói thiếu |

**Về việc client nuốt lỗi — đọc kỹ, nó đổi cách chẩn đoán:**
`Controller.cs:118-122` mở một `try` bao trọn bộ phân giải gói, đóng ở `:4600-4602` bằng
`catch (Exception) { }` **rỗng hoàn toàn**. Cộng với `myReader.cs:51-93` đọc `buffer[posRead++]`
**không kiểm biên**, hệ quả là: gói thiếu byte → ném `IndexOutOfRangeException` → **bị nuốt
sạch**, và **trạng thái đã gán một nửa vẫn giữ nguyên**. Nghĩa là một gói hành trang thiếu
byte không chỉ "không làm gì" — nó để lại **túi đồ hiển thị dở dang**. Đây là lý do §6.18
nói chạy thử bằng mắt không phát hiện được, và là lý do mọi gói ở §4 phải có test
byte-chính-xác **kèm ca ô rỗng**.

Cả L1 và L2 đều là lỗi **im lặng tuyệt đối**. Mỗi lỗi phải có một test đối chiếu số byte.

### 4.6c Gói XUỐNG `-21` / `-19` — đồ dưới đất biến mất / người khác nhặt

| Lệnh | Ý nghĩa | Nguồn |
|---|---|---|
| `-21` | Món dưới đất biến mất (hết giờ) | Java `ItemMapService.java:41-51` ↔ client `Controller.cs:3381-3394` |
| `-19` | **Người khác** vừa nhặt món | Java `Service.java:1404-1415` ↔ client `Controller.cs:3475-3500` |

Bản 1 cần cả hai: thiếu `-21` thì món hết hạn vẫn nằm trên màn hình người chơi mãi mãi;
thiếu `-19` thì hai người cùng thấy một món, một người nhặt được, người kia **vẫn thấy món
nằm đó** và bấm mãi không được — đúng loại lỗi §6.32 (gửi một thứ mà không xử lý phần còn lại).

### 4.6d Gói XUỐNG `6` và `-30/4` — cập nhật tiền

| Lệnh | Thân gói | Nguồn |
|---|---|---|
| `6` | `long vang, int ngoc, int hongNgoc` | Java `Service.java:1387-1402` ↔ client `Controller.cs:3677-3685` |
| `-30/4` | `long vang, int ngoc, int hp, int chakra, int hongNgoc` | Java `PlayerService.java:114-138` ↔ client `Controller.cs:6016-6025` |

⚠️ **`vang` là `long` 8 byte** — cùng ràng buộc đã chốt cho ME_INFO (spec vào thế giới §11.4).
Bản 1 ghi `hongNgoc = 0` (chưa có loại tiền này). ⚠️ Bản Java chỉ ghi `long` khi
`session.version >= 214`, ngược lại ghi `int` (L9) — **bản Go luôn ghi `long`**.

### 4.7 Gói XUỐNG `-86` — giao dịch 2 người *(bản 2, ghi lại để không phải giải mã lại)*

Java: `Trade.java:77-90, 180-199, 256-259, 285-330` · Client: `Controller.cs:628-720`

| Lệnh con | Ý nghĩa | Thân gói | Ghi chú |
|---|---|---|---|
| `0` | mời giao dịch | `int idNguoiMoi` | client `Controller.cs:631-635` |
| `1` | mở bảng giao dịch | `int idDoiTac` (Java :79/:84) | ⚠️ client gọi `findCharInMap(id)` và **`return` nếu không thấy** (:640-643) → đối tác **phải ở cùng khu vực** |
| `2` | gỡ một món khỏi bàn | `byte oSo` (Java :182) | client tìm theo `indexUI` (:650-659) |
| `6` | khoá bàn + gửi danh sách bên kia | `int vang`, `byte soMon`, *mỗi món:* `short templateId`, `int soLuong`, `byte soOption`, *option* | Java :288-302 / client :673-694 — ✅ khớp |
| `7` | huỷ giao dịch | *(rỗng)* | Java :256-259 |

### 4.8 Gói LÊN — client → server

| Lệnh | Ý nghĩa | Thân gói | Nguồn |
|---|---|---|---|
| `6` | **Mua ở cửa hàng NPC** | `byte typeBuy, short templateId`, **`short soLuong` chỉ khi còn ≥2 byte** | `Controller.java:329-331` — ⚠️ **độ dài biến thiên**, đúng bẫy §6.6. Bản Go: đọc theo số byte còn lại, thiếu thì mặc định 1 |
| `7` | **Bán cho NPC** | `byte action`, `byte where`, `short oSo` | `Controller.java:352-359`; `action == 0` = xin xác nhận, khác 0 = bán thật |
| `-43` | **Dùng / bỏ đồ** | `byte type, byte where, byte oSo`, **+ `short templateId` nếu `oSo == -1`** | `UseItem.java:140-142, 190` ↔ client `Service.cs:415-421`. `type`: 0 dùng · 1 vứt · 2 xác nhận vứt · 3 xác nhận dùng (`UseItem.java:52-55`). Gói **xuống** cùng mã: `byte type, byte where, byte oSo, UTF cauHoi` (`UseItem.java:153-158` ↔ `Controller.cs:1248-1252`). ⚠️ **độ dài biến thiên** — mỗi dạng một hàm riêng (§6.6) |
| `-40` | **Chuyển ô** (rương ↔ túi ↔ mặc) | `byte loai, byte oSo` | `Controller.java:622-630` → `UseItem.java:44-50, 81-109` ↔ client `Service.cs:605-608`. `loai`: 0 rương→túi/mặc · 1 túi→rương · 3 mặc→rương · 4 túi→mặc · 5 mặc→túi · 6/7 đệ tử (bản 1 **bỏ 6/7**) |
| `-20` | **Nhặt đồ** | `short itemMapId` | `Controller.java:723-727` |
| `-86` | Giao dịch *(bản 2)* | `byte lenhCon` + tuỳ | `Controller.java:302-306` |
| `-100` | Ký gửi *(bản 3)* | `byte action` + tuỳ: `0` ký gửi = `short idItem, byte moneyType, int gia, int soLuong`; `3` mua = `short idItem, byte, int` (**hai trường sau bị đọc rồi vứt**) | `Controller.java:128-167` |

### 4.9 Chỗ nối sẵn có trong bản Go

| Chỗ | Hiện tại | KH7 làm gì |
|---|---|---|
| `thegioi.go:79-81` (ME_INFO `-30/0`) | ghi `0, 0, 0` cho `itemsBody/Bag/Box.size` | **Điền số thật.** ⚠️ `itemsBody.size` client đọc **unsigned** (≤255), hai ô kia có dấu (≤127) — spec vào thế giới §11.2 |
| `dongbo.go:142-153` (catalog item `-28/8`) | gửi **đúng một** message subtype 2 **rỗng** | **Điền `ItemTemplate` thật.** Không có mẫu thì client không biết vẽ icon gì, tên gì |
| `login.go:282-318` (`XuLy`) | chưa có case nào của KH7 | Thêm `6`, `7`, `-43`, `-20` |
| `goichiendau.go` (`-12` quái chết) | `soDoRoi = 0` cố định (KH6 hoãn có chủ đích) | **Bản 1 giữ nguyên 0.** Rơi đồ nối vào sau khi `vat_pham` chạy thật — vì mỗi món rơi là một dòng CSDL sinh ra **từ trong actor**, đúng chỗ §6.10b cấm chạm ổ đĩa |

---

## 5. Mô hình dữ liệu

### 5.1 DDL phác (migration `0004_vat_pham.sql`)

Máy chủ **tự áp migration** khi khởi động (`store/migrate.go`) — chỉ cần thêm file, không
cần chạy tay.

```sql
-- Moi MON (instance) mot dong. id BIGSERIAL, KHONG BAO GIO tai su dung.
CREATE TABLE IF NOT EXISTS vat_pham (
    id              BIGSERIAL   PRIMARY KEY,
    mau_id          INTEGER     NOT NULL,          -- ItemTemplate.id (short xuong client)
    chu_so_huu_id   BIGINT      REFERENCES nhan_vat(id),  -- NULL khi duoi dat / da huy
    noi_chua        TEXT        NOT NULL,          -- TUI | MAC | RUONG | DUOI_DAT | KY_GUI | DA_HUY
    o_so            SMALLINT,                      -- vi tri hien thi; NULL khi khong nam trong o nao
    so_luong        INTEGER     NOT NULL DEFAULT 1,
    chi_so_phu      JSONB       NOT NULL DEFAULT '{"opt":[]}'::jsonb,
    tao_luc         TIMESTAMPTZ NOT NULL DEFAULT now(),
    sua_luc         TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT vat_pham_noi_chua_hop_le
        CHECK (noi_chua IN ('TUI','MAC','RUONG','DUOI_DAT','KY_GUI','DA_HUY')),
    -- BB2: so luong khong bao gio am, va mon co option thi luon la 1 cai
    CONSTRAINT vat_pham_so_luong_duong CHECK (so_luong > 0),
    CONSTRAINT vat_pham_co_option_thi_don_chiec
        CHECK (jsonb_array_length(chi_so_phu->'opt') = 0 OR so_luong = 1),
    -- Mon co chu thi phai o mot noi cua nguoi do; mon khong chu thi phai duoi dat/da huy
    CONSTRAINT vat_pham_chu_va_noi_khop CHECK (
        (chu_so_huu_id IS NOT NULL AND noi_chua IN ('TUI','MAC','RUONG','KY_GUI'))
     OR (chu_so_huu_id IS NULL     AND noi_chua IN ('DUOI_DAT','DA_HUY'))
    )
);

-- BB1: mot o cua mot nguoi chi chua DUNG MOT mon. Day la rang buoc CHONG NHAN BAN
-- manh nhat trong ca KH7 — no lam cho "hai mon cung o" KHONG THE ton tai, ke ca khi
-- tang Go co bug. Partial index vi mon duoi dat/da huy khong co o.
CREATE UNIQUE INDEX IF NOT EXISTS vat_pham_mot_o_mot_mon
    ON vat_pham (chu_so_huu_id, noi_chua, o_so)
    WHERE chu_so_huu_id IS NOT NULL AND o_so IS NOT NULL;

-- Truy van nong nhat: "lay ca tui cua nguoi X" moi lan dang nhap / moi lan doi tui.
CREATE INDEX IF NOT EXISTS vat_pham_theo_chu
    ON vat_pham (chu_so_huu_id, noi_chua)
    WHERE chu_so_huu_id IS NOT NULL;

-- BB3: moi lan doi chu mot dong. CHI THEM, khong bao gio sua/xoa.
CREATE TABLE IF NOT EXISTS vat_pham_lich_su (
    id            BIGSERIAL   PRIMARY KEY,
    vat_pham_id   BIGINT      NOT NULL REFERENCES vat_pham(id),
    tu_ai         BIGINT      REFERENCES nhan_vat(id),   -- NULL = tu he thong sinh ra
    sang_ai       BIGINT      REFERENCES nhan_vat(id),   -- NULL = bi tieu huy
    ly_do         TEXT        NOT NULL,   -- SINH_RA|MUA_NPC|BAN_NPC|NHAT|DUNG|BO|GIAO_DICH|ADMIN
    so_luong      INTEGER     NOT NULL,
    luc           TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS vat_pham_lich_su_theo_mon ON vat_pham_lich_su (vat_pham_id);

-- D3: chong lam hai lan cung mot hanh dong kinh te.
CREATE TABLE IF NOT EXISTS vat_pham_giao_dich (
    id          BIGSERIAL   PRIMARY KEY,
    khoa_idem   TEXT        NOT NULL UNIQUE,   -- "<nhanVatId>:<loai>:<soThuTu>"
    ket_qua     JSONB       NOT NULL,          -- de tra lai DUNG ket qua lan dau khi goi lap
    luc         TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### 5.2 Vì sao đặt ràng buộc ở **CSDL** chứ không chỉ ở Go

Đây là bài học đã trả giá của chính dự án này: `nhan_vat` có `CHECK` tên hợp lệ **ở CSDL**
(`0002_nhan_vat.sql:64-73`) với lý do ghi thẳng trong file — *"để web/panel/script import
chèn thẳng cũng không lách được"*. Vật phẩm còn cần điều đó hơn: web panel, script chuyển
dữ liệu, và tay người vận hành đều chạm được bảng này. `vat_pham_mot_o_mot_mon` là tuyến
phòng thủ **cuối cùng** của BB1 — nó đúng ngay cả khi toàn bộ tầng Go có bug.

### 5.3 Đường lưu: **không** đi qua `ThoGhi`

**Chốt: vật phẩm và tiền ghi thẳng bằng giao dịch CSDL trên goroutine phiên, KHÔNG qua `ThoGhi`.**

| | `ThoGhi` (đang dùng cho nhân vật) | Vật phẩm/tiền (KH7) |
|---|---|---|
| Ngữ nghĩa ghi | **ghi đè cả dòng** bằng ảnh chụp (`store.go:194-206`) | **cộng/trừ tại chỗ**, có điều kiện |
| Mất mát khi sập | chấp nhận mất tới 60 giây tiến độ | **không được mất một xu** |
| Gộp nhiều lần ghi | có (map id → bản mới nhất, `thogi.go:120-129`) | **cấm** — mỗi hành động là một sự thật riêng |
| Chạy ở đâu | goroutine thợ ghi riêng | goroutine **phiên** (được phép chờ ổ đĩa — `thogi.go:196-201` đã nói rõ luật §6.10b **chỉ** cấm `KhuVuc.vongLap`) |

**Task bắt buộc làm trước mọi thứ khác (chặn D7):** gỡ `vang`/`ngoc` khỏi câu `UPDATE` của
`CapNhatNhanVat`. Sau khi gỡ, `ThoGhi` không còn đường nào ghi đè tiền, và tiền chỉ đổi qua
đúng một hàm cộng/trừ có điều kiện. **Nếu bỏ qua task này thì mọi chốt chặn khác của §2 đều
vô nghĩa** — thợ ghi sẽ lặng lẽ hoàn tiền cho người chơi sau mỗi lần mua.

---

## 6. Chia task thi công

Cột **"chạm `khuvuc.go`?"** là để biết task nào **phải chạy tuần tự**: `khuvuc.go` hiện
~63KB (`khuvuc.go` §10b `BAT-DAU-TU-DAY.md`), hai agent cùng vào là chắc chắn đụng.

| # | Task | File mục tiêu | Phụ thuộc | Chạm `khuvuc.go`? |
|---|---|---|---|---|
| **0** | ⛔ **GATE — chặn D7:** gỡ `vang`/`ngoc` khỏi `CapNhatNhanVat`; thêm `CongTruTien(ctx, nvID, dVang, dNgoc) error` cộng/trừ tại chỗ, từ chối khi âm. Test đột biến: trả `UPDATE` về ghi đè → phải đỏ | `server/internal/store/store.go` (+ test) | không | **không** — chạy song song ngay từ đầu |
| **1** | Migration `0004_vat_pham.sql` + kiểu `store.VatPham` + test ràng buộc CSDL thật (chèn 2 món cùng ô → phải lỗi UNIQUE) | `server/internal/store/migrations/0004_vat_pham.sql`, `store/vatpham.go` | không | **không** — song song |
| **2** | Danh mục mẫu vật phẩm (`MauVatPham`) + gói `-28/8` subtype 0/1 có nội dung thật. ⚠️ đụng `asset-lock.json` — xem bẫy §7.6 | `server/internal/game/mauvatpham.go` | không | **không** — song song |
| **3** | Bộ dựng gói **`-36`** (túi), **`-37`** (mặc), **`-35`** (rương) — test byte-chính-xác, **có ca ô rỗng** (bắt lỗi L1) | `server/internal/game/goihanhtrang.go` | 1, 2 | **không** — song song |
| **4** | Bảng phiên dịch `o_so → vat_pham.id` mỗi phiên, dựng **cùng lúc** với gói ở task 3 | `server/internal/game/hanhtrang.go` | 3 | **không** |
| **5** | Điền `itemsBody/Bag/Box.size` thật vào ME_INFO | `server/internal/game/thegioi.go` | 3, 4 | **không** (file khác, nhưng **tuần tự** với task 3) |
| **6** | Kho vật phẩm: `LayTui`, `ChuyenO`, `CapVatPham`, `HuyVatPham` — **tất cả trong giao dịch, ghi lịch sử cùng giao dịch**, khoá `FOR UPDATE` theo id tăng dần | `server/internal/store/vatpham.go` | 1 | **không** |
| **7** | Handler `-43` (dùng/bỏ đồ) + `-40` (chuyển ô túi↔rương↔mặc) theo đúng trình tự §2.3; **kiểm biên `oSo`** ở mọi nhánh (bẫy §7.11) | `server/internal/game/dungvatpham.go` + `login.go` (thêm case) | 4, 6 | **không** |
| **8** | Cửa hàng NPC: gói `-44` `loaiShop=0` (test byte) + danh mục hàng | `server/internal/game/cuahang.go` | 2 | **không** |
| **9** | Handler mua `6` — **có idempotency key** (D3), trừ tiền + cấp món trong **một** giao dịch. Đột biến: tách thành 2 giao dịch → test phải đỏ | `server/internal/game/muaban.go` + `login.go` | 0, 6, 8 | **không** |
| **10** | Handler bán `7` — cùng kỷ luật | `server/internal/game/muaban.go` | 9 | **không** |
| **11** | Đồ dưới đất + handler nhặt `-20` — `UPDATE ... WHERE noi_chua='DUOI_DAT'` một bước (D8); **ghi đủ 2 short cho type 9/10** (lỗi L2) | `server/internal/game/dovuong.go` + `khuvuc.go` | 6 | ✅ **CÓ — tuần tự, không song song với bất kỳ task nào khác chạm `khuvuc.go`** |
| **12** | Nối rơi đồ khi quái/boss chết: `-12` `soDoRoi > 0`. Actor **chỉ xếp yêu cầu**, một goroutine khác mới chạm CSDL (§6.10b) | `server/internal/game/khuvuc.go` + `goichiendau.go` | 11 | ✅ **CÓ — tuần tự, sau task 11** |
| **13** | Nghiệm thu tay bằng client Unity thật: mở túi, mua, bán, dùng, nhặt | — | 12 | — (chỉ chủ dự án làm được) |

**Kỷ luật TDD + đột biến** (§10 `BAT-DAU-TU-DAY.md`, làm như KH6b): mỗi task **test đỏ trước,
code sau**; xong task → agent **khác** soi. Cuối nhánh chạy **một đợt đột biến** như KH6b
(28 đột biến). Bốn đột biến **bắt buộc phải bị giết**, vì mỗi cái tương ứng một bất biến:

| Đột biến | Bất biến bị phá | Test nào phải đỏ |
|---|---|---|
| Bỏ `FOR UPDATE` trong luồng mua | BB1 / D5 | test hai goroutine cùng mua một món |
| Tách trừ tiền và cấp món thành 2 giao dịch | BB2 / D2 | test huỷ giữa chừng |
| Bỏ ghi `vat_pham_lich_su` | BB3 | test đếm dòng lịch sử sau mỗi thao tác |
| Trả `CapNhatNhanVat` về ghi đè `vang` | BB4 / D7 | test "mua rồi thợ ghi xả hàng đợi" |

---

## 7. Bẫy — chỉ đọc client mới thấy

1. **`param` của option là UNSIGNED short phía client** (`myReader.cs:222-229` + `isReadInt=true`).
   Ghi số âm → người chơi thấy chỉ số khổng lồ, **không lỗi, không log** (§6.18).
2. **`optionId` chỉ 1 byte** vì `isShortOptionTemp = false` (`ModFunc.cs:300`). Nếu ai đó
   sau này bật cờ đó trong client thì **mọi** gói item lệch — phải ghim giá trị cờ này vào
   test, đúng tinh thần §6.12 ("test dùng giá trị ghim").
3. **`-36` phải ghi `-1` cho ô rỗng** dù bản Java không ghi (lỗi L1, §4.6b). Người thi công
   sẽ có xu hướng "chép cho giống bản gốc" — chỗ này **cố tình khác**, phải có chú thích
   trong code nói rõ vì sao, kẻo người sau "sửa lại cho giống".
4. **`-20` với vàng/ngọc cần HAI short** (lỗi L2). Chỉ nhìn code Java thì không đời nào thấy.
5. **`-86/1` yêu cầu đối tác ở cùng khu vực** — client `return` ngay nếu `findCharInMap` trả
   `null` (`Controller.cs:640-643`). Bản 2 phải kiểm điều kiện này **phía máy chủ trước khi
   gửi**, nếu không người chơi bấm giao dịch mà **không có gì xảy ra**, không thông báo gì.
6. **`catalog.json` bị ghim SHA trong `asset-lock.json`** (`catalog.go:100`, `lockHashKhop`).
   Thêm mẫu vật phẩm vào đó là **đỏ test** và phải sửa artifact đã ký. Task 2 phải quyết
   trước: (a) để mẫu vật phẩm trong file Go riêng như `noidungboss.go`, hay (b) đi qua quy
   trình lock. **Khuyến nghị (a)** — cùng cách KH6b đã làm với dữ liệu boss.
7. **Chưa có icon cho vật phẩm.** `asset-lock.json` closure không có sprite icon nào cho đồ.
   Bản 1 phải hoặc dùng lại icon sẵn có, hoặc bổ sung qua quy trình mỹ thuật + lock. **Sai
   `iconId` là lỗi im lặng** (ô trống trên giao diện), giống hệt bẫy bg-item §6.33.
8. **Cấm gửi cho client một thứ mà không xử lý phần trả lời của nó** (§6.32). Gửi `-44` mở
   cửa hàng mà không xử lý gói lên `6` → người chơi bấm mua và **không có gì xảy ra**. Phải
   tìm "van mở khoá" **trước khi** viết tính năng.
9. **Gói lên `6` có độ dài biến thiên** (`Controller.java:331`: `available() >= 2`). Theo
   §6.6, dựng **mỗi dạng một hàm riêng**, đừng dùng một hàm có cờ.
10. **`sortItems` của bản Java là cái bẫy gốc** (`InventoryService.java:562`). Bản Go
    **không** được sắp xếp lại túi ở phía máy chủ giữa hai lần người chơi thao tác; nếu có
    sắp xếp thì phải dựng lại bảng phiên dịch §4.1 **trong cùng một hàm** với lúc gửi gói.
    Đáng sợ hơn: chính `sortItems` là lý do **duy nhất** lỗi L1 chưa làm vỡ game thật —
    nó dồn ô trống về cuối nên gói thường không thiếu byte. Tức là bản gốc đang sống nhờ
    một sự trùng hợp, không phải nhờ đúng.
11. **Chỉ số ô từ client KHÔNG BAO GIỜ được dùng thẳng.** Bản Java `get(index)` không kiểm
    biên ở ít nhất bốn chỗ (`Trade.java:109-113`, `UseItem.java:150`,
    `ShopService.java:903-914`, `InventoryService.java:116-136`) và dựa vào `catch` tổng để
    khỏi sập. Bản Go: mọi `oSo` phải qua bảng phiên dịch §4.1, ngoài miền thì **từ chối và
    trả gói lỗi**, không im lặng.
12. **Một gói tin có thể ép máy chủ chạy 32.767 vòng.** `Controller.java:329-337`: mua NPC
    lặp `quantity` lần (đọc bằng `short`), mỗi vòng trừ tiền + thêm đồ. Bản Go: **kẹp số
    lượng mua mỗi gói** (đề xuất ≤99, đúng với trần chồng đồ) và xử lý **một giao dịch cho
    cả lô**, không lặp N giao dịch.
13. **Món "không có option" vẫn bị bản Java nhét option giả `(73, 0)`** khi vào túi
    (`InventoryService.java:811-813`). Nếu bản Go không nhét mà client lại phụ thuộc, giao
    diện có thể khác. **Chưa kiểm client có phụ thuộc hay không** — task 3 phải thử cả hai.

---

## 8. Chưa xác minh được — nói thẳng

1. **Chưa chạy thử một gói item nào với client Unity thật.** Toàn bộ §4 là đối chiếu *đọc
   code* hai phía. Đúng là mức chứng cứ cao nhất từng dùng cho KH6/KH6b, nhưng vẫn **không
   thay được nghiệm thu tay** — và §6.18 nói rõ lỗi trường bên trong gói là **im lặng tuyệt
   đối**. Task 13 là bắt buộc.
2. **Không biết `ItemTemplate` cần những trường gì cho gói `-28/8` subtype 0/1.** Tôi đã xác
   minh subtype 2 rỗng là đủ để **không kẹt "Xin chờ"** (spec vào thế giới §6), nhưng **chưa**
   giải mã khuôn đầy đủ của subtype 0 (bảng option template) và 1 (mẫu vật phẩm 0..750).
   **Task 2 phải bắt đầu bằng một bước giải mã hai phía** (`ItemData.java` ↔
   `Controller.cs` `loadItemNew` ~:6453-6536) rồi mới viết code — giống GATE task 1 của KH6b.
3. **Chưa biết client xử lý ra sao khi túi vượt 127 ô.** Java `writeByte` có dấu, client
   `readUnsignedByte`. Bản 1 kẹp **≤127** cho an toàn, nhưng nếu muốn túi lớn hơn thì phải
   kiểm client trước. Chưa kiểm.
4. **Chưa đo chi phí một giao dịch CSDL cho mỗi lần mua.** Kỷ luật §2.3 đúng về mặt an toàn,
   nhưng chưa biết Postgres chịu bao nhiêu lượt mua/giây với 20+ người chơi (KH6c mới chỉ
   đo actor, không đụng ổ đĩa). Nếu chậm thì phương án là gộp theo lô — **nhưng không được
   gộp bằng cách nới ràng buộc**.
5. **Chưa quyết bản 1 có "rương" thật hay không.** Gói `-35` đã có khuôn, nhưng rương cần
   NPC giữ rương và luồng chuyển đồ túi ↔ rương. Task 3 dựng **gói**, còn có bật tính năng
   hay không thì chủ dự án chốt.
6. **Số ô túi khởi điểm và luật mở rộng túi chưa có.** Là quyết định cân bằng (KH8), không
   phải kỹ thuật. Bản 1 tạm ghim một con số trong mã và ghi chú rõ là tạm.
7. **Chưa kiểm client có phụ thuộc option giả `(73, 0)` hay không** (bẫy §7.13).
8. **Chưa kiểm cái bẫy khung gói ở `Session_ME.cs:203-205`.** Khi chưa trao khoá, client
   tính độ dài bằng `(num2 & 0xFF00) | (b5 & 0xFF)` với `num2` kiểu `sbyte` — phép `& 0xFF00`
   **luôn ra 0**, tức byte cao bị mất. Mọi gói **dài hơn 255 byte gửi trước lúc trao khoá**
   sẽ đọc lệch. Gói hành trang và gói mở cửa hàng **rất dễ vượt 255 byte**. Bản Go hiện gửi
   chúng **sau** khi đã trao khoá nên không dính, nhưng đây là một quả mìn: ai đó dời thứ tự
   gửi là vỡ. **Phải có một test ghim thứ tự "trao khoá trước, gói lớn sau".**

---

## 9. Nội dung game — cố ý chưa chốt

Đơn giá, danh mục hàng NPC, tỉ lệ rơi đồ, chỉ số option: **thuộc KH8** (kho nội dung YAML +
cân bằng thực nghiệm). KH7 chỉ dựng **khung**: căn cước, chống nhân bản, giao thức, luồng.
Bản 1 dùng vài mẫu vật phẩm tối thiểu để chạy được đầu-cuối, đúng cách KH6b làm với 2 boss mẫu.

Một ghi chú để không bất ngờ sau này: chợ ký gửi bản Java **không dùng vàng làm tiền tệ** mà
dùng một vật phẩm ("thỏi vàng", `ConsignShopService.java:31`) — tức là tiền tệ của chợ chính
là một dòng trong `vat_pham`. Bản 3 phải chốt sớm: giữ cách đó (nhất quán, mọi thứ đều là
món có căn cước) hay dùng thẳng cột `vang`. **Chưa quyết.**
