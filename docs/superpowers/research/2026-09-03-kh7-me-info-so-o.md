# PHA A — ME_INFO (`-30/0`) đọc ba mảng vật phẩm thế nào

KH7 Task 5. Nguồn duy nhất: client Unity `D:\NRO_Project\Assets\Scripts\Game1\Controller.cs`
(chỉ đọc) và server Java cũ `1-MAY-CHU-GAME\src\`.

## 1. Phát hiện quan trọng nhất — ME_INFO KHÔNG chỉ mang "ba con số"

Đặc tả §4.9 viết "điền số thật vào `itemsBody/Bag/Box.size`". Đọc client mới thấy: **ngay
sau mỗi con số, client đọc luôn NỘI DUNG từng ô, ngay trong gói ME_INFO**.

`Controller.cs:5860-5968`:

```csharp
Char.myCharz().arrItemBody = new Item[msg.reader().readUnsignedByte()];   // :5860
for (int k = 0; k < arrItemBody.Length; k++) {
    short num7 = msg.reader().readShort();      // :5866  mẫu, -1 = ô rỗng
    if (num7 == -1) continue;                   // :5867
    ... readInt (số lượng), readUTF (info), readUTF (content),
        readUnsignedByte (số option), rồi mỗi option: readShortOptionTemp + ReadParamDQT
}
Char.myCharz().arrItemBag = new Item[msg.reader().readByte()];            // :5907  CÓ DẤU
   ... y hệt, số option đọc bằng readByte (:5922)
Char.myCharz().arrItemBox = new Item[msg.reader().readByte()];            // :5942  CÓ DẤU
   ... y hệt, số option đọc bằng readByte (:5956)
Char.myCharz().statusMe = 4;                                             // :5969
```

Hệ quả **cứng**: ghi `size = 20` mà không ghi 20 thân ô thì client đọc lấn sang phần sau
(`headAvatar`, `conMeo`, `fusion`…) → lệch toàn bộ đuôi gói. Vì vậy Task 5 **bắt buộc**
phải ghi cả thân ô vào ME_INFO, không thể chỉ sửa ba byte. Đây là lý do bản hiện tại ghi
`0,0,0` (không thân) mà vẫn chạy được.

Khuôn thân ô **giống hệt** khuôn của `-36/-37/-35` (`game/goihanhtrang.go`, hàm `ghiCacO`):
`short mauId | -1`, `int soLuong`, `UTF info`, `UTF content`, `byte soOption`, rồi mỗi
option `byte optionId` + `unsigned short param`. Không có trường thừa, không có trường
thiếu. Đã đối chiếu từng dòng `:5866-5890` với `:2210-2237` (gói `-36`).

## 2. Số Ô hay số MÓN?

**Số Ô.** Ba chỗ chứng minh:

1. Client cấp mảng đúng bằng con số đó rồi **lặp đủ con số đó lần**, mỗi lần đọc một ô;
   ô rỗng là `short -1` rồi `continue` (`:5866-5870`, `:5911-5915`, `:5946-5950`). Nếu
   con số là "số món" thì không có cách nào diễn tả "ô 0 có món, ô 1 rỗng, ô 2 có món" —
   mà `indexUI = m` (`:5921`) chính là **chỉ số ô** client gửi ngược lên khi người chơi
   bấm vào một ô.
2. `arrItemBag[m].indexUI = m` — client dùng đúng chỉ số vòng lặp làm định danh ô. Máy chủ
   tra `o_so` → `vat_pham.id` bằng chính con số này (`game/hanhtrang.go`, `TraId`). Hai bên
   chỉ khớp khi vòng lặp chạy qua **mọi ô**, kể cả ô rỗng.
3. Java lưu số ô trong CSDL: `Inventory.java:36-38` khởi tạo danh sách rỗng, rồi
   `NDVSqlFetcher.java:467/530` `add(...)` **một phần tử cho mỗi phần tử của mảng JSON**,
   phần tử `[-1,0,"[]",…]` là **ô rỗng** (`createItemNull()`, `:466`). Vật phẩm mở rộng túi
   (`InventoryService.java:789-800`, mẫu 1627/517/518) `add` thêm **một ô rỗng** — bằng
   chứng trực tiếp rằng `size()` là **sức chứa**, không phải số món.

Trần: `MAX_ITEMS_BAG = MAX_ITEMS_BOX = 100` (`Inventory.java:17-18`) — đó là trần mở rộng
của bản Java, nằm dưới trần giao thức nên không phải là ràng buộc ở đây.

## 3. Dấu của ba con số

| Mảng | Client đọc bằng | Trần thật |
|---|---|---|
| `arrItemBody` (đồ mặc) | `readUnsignedByte()` `:5860` | 255 |
| `arrItemBag` (túi) | `readByte()` `:5907` — **CÓ DẤU** | 127 |
| `arrItemBox` (rương) | `readByte()` `:5942` — **CÓ DẤU** | 127 |

Ghi 200 vào hai ô có dấu → client nhận `-56` → `new Item[-56]` → `OverflowException`. Ở
túi/rương ngoại lệ đó **không nằm trong `try`** (`try/catch` chỉ bọc khối `arrItemBody`,
`:5861-5906`) → nổ giữa hàm xử lý gói, `statusMe = 4` (`:5969`) **không bao giờ chạy**,
người chơi kẹt màn "Xin chờ". Đây là ca xấu nhất và là lý do trần chung của cả dự án chốt
ở **127** (`goihanhtrang.go`, hằng `tranSoO`).

Ngược lại, `arrItemBody` đọc unsigned nên nếu có ngày muốn >127 ô đồ mặc thì **chỉ** ô đó
được, hai ô kia thì không.

## 4. Thứ tự gói bắt buộc

- **`-36/-37/-35` phải đến SAU `-30/0`.** ME_INFO **cấp lại** cả ba mảng
  (`new Item[...]` ở `:5860/:5907/:5942`) — mọi thứ nạp trước nó bị xoá sạch không báo.
  Ngược lại `-36` đến sau chỉ nạp lại đúng mảng đó, vô hại.
- `-36` (`:2202`), `-37` (`:2005`), `-35` (`:2265`) **không đòi** ME_INFO trước về mặt cấu
  trúc: mỗi gói tự cấp mảng của mình, tự đọc `soO` riêng. Không có phụ thuộc nào khác.
- `-37` có thêm `short head` chen giữa lệnh con và `soO` (`:2011`); ME_INFO **không** có
  trường đó ở khối đồ mặc (head của ME_INFO nằm tận đầu gói, trường `KieuToc`).
- Ba gói này đến trước hay sau `-24` MAP INFO đều được (client không đọc `GameScr` khi xử
  lý chúng). Bản Go gửi **sau `-42` chỉ số**, cuối chuỗi vào thế giới, để chuỗi cũ
  (`-31 → -82 → -30 → -24 → -42`) không đổi thứ tự — mọi test bắt tay hiện có vẫn đúng, chỉ
  thêm đuôi.

## 5. Chốt cho Task 5

1. ME_INFO ghi **số ô thật + đủ thân ô** cho cả ba mảng, dựng từ **đúng một** lát cắt
   `LayTui` (cùng lát cắt dựng `BangHanhTrang` và ba gói) — không có đường nào lệch.
2. Sau `-42`, gửi thêm `-36/-37/-35` từ **cùng lát cắt đó**. Dư thừa về dữ liệu nhưng
   không mâu thuẫn được, và là mặt cắt mà test tích hợp quan sát được.
3. Lỗi CSDL khi lấy hành trang ⇒ **không cho vào thế giới nửa vời**: trả lỗi lên
   `XuLy` (đóng phiên) sau khi ghi log. Gửi túi rỗng "có chủ đích" bị loại: người chơi sẽ
   thấy hành trang trống và tưởng mất đồ, còn bảng tra ô rỗng khiến mọi thao tác của task
   7/9 im lặng không làm gì (§1.10 "không gửi gói thiếu").
