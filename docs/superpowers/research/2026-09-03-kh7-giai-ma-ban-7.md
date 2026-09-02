# PHA A — Giải mã gói `7` (bán hàng cho NPC) và gói phản hồi

KH7 Task 10. Đọc trước khi viết một dòng code bán nào. Bản này là **ảnh gương** của
[giải mã mua `6`](2026-09-03-kh7-giai-ma-mua-6.md) — chỗ nào giống thì nói giống, chỗ nào
khác thì nói rõ vì sao khác.

Nguồn (chỉ đọc):
- Client Unity: `D:\NRO_Project\Assets\Scripts\Game1\Service.cs`, `Controller.cs`, `Panel.cs`
- Server Java cũ: `D:\NRO_GoiBanGiao\1-MAY-CHU-GAME\src\server\Controller.java`,
  `src\shop\ShopService.java`

---

## 1. Gói LÊN `7` — client xin bán

`Service.cs:937-956`:

```csharp
public void saleItem(sbyte action, sbyte type, short id)
{
    message = new Message((sbyte)7);
    message.writer().writeByte(action);
    message.writer().writeByte(type);
    message.writer().writeShort(id);
    session.sendMessage(message);
}
```

Server Java đọc y hệt (`Controller.java:341-361`):

```java
action = _msg.reader().readByte();
if (action == 0) {
    ShopService.gI().showConfirmSellItem(player, _msg.reader().readByte(), _msg.reader().readShort());
} else {
    ShopService.gI().sellItem(player, _msg.reader().readByte(), _msg.reader().readShort());
}
```

**Khuôn chốt — ba ô, KHÔNG có ô nào là số lượng:**

| Ô | Kiểu | Ý nghĩa |
|---|---|---|
| 1 | `byte` | `action` — `0` = **xin xác nhận** (chưa bán gì cả), khác 0 = **bán thật** |
| 2 | `byte` | `type`/`where` — **nơi chứa**: `0` = đồ đang MẶC (`itemsBody`), khác 0 = TÚI (`itemsBag`) |
| 3 | `short` | `id`/`index` — **chỉ số Ô** trong nơi chứa đó (KHÔNG phải `MauID`, khác hẳn gói `6`) |

Hai chỗ client gọi (`Panel.cs`):

- `Panel.cs:13505` (`idAction == 3002`): `saleItem(0, isBody ? 0 : 1, ô)` — bấm "Bán" trong túi
  → **xin xác nhận**.
- `Panel.cs:13511` (`idAction == 3003`): `saleItem(1, type, id)` — bấm Đồng ý trong hộp thoại
  xác nhận → **bán thật**. `type`/`id` client lấy lại đúng từ gói xác nhận máy chủ vừa gửi.

Vậy luồng là **hai vòng**: client hỏi → máy chủ trả hộp thoại → client đồng ý → máy chủ bán.

### 1.1 Khác gói mua ở đâu (đọc kỹ chỗ này)

| | Mua `6` | Bán `7` |
|---|---|---|
| Ô định danh | `MauID` (mẫu vật phẩm) | **chỉ số ô** trong túi/mặc |
| Số lượng | có ô thứ 3, có thể vắng | **không có ô nào** → luôn bán CẢ CHỒNG |
| Số vòng | một vòng | **hai vòng** (xác nhận rồi mới bán) |

Vì ô thứ ba là **chỉ số ô**, đường bán bắt buộc phải đi qua bảng phiên dịch
ô → `vat_pham.id` (task 4, `hanhtrang.go` `BangHanhTrang.TraId`) — `vat_pham.id` không bao
giờ xuống client (§4.1).

### 1.2 Bẫy giống hệt mua: máy chủ mới giữ "đang mở cửa hàng nào"

`ShopService.sellItem:899-902` từ chối ngay nếu `pl.iDMark.getShopOpen() == null`. Bản Go
làm y vậy: `BoXuLy.muaBanCuaHangDangMo == nil` → từ chối. Đây là chốt chặn "gửi gói `7`
khống để bán đồ từ xa".

---

## 2. Gói XUỐNG — client mong nhận gì

### 2.1 Gói `7` xuống — hộp thoại xác nhận

`ShopService.showConfirmSellItem:882-888` ghi, `Controller.cs:3669-3676` đọc:

```csharp
case 7:
    sbyte type  = msg.reader().readByte();
    short id    = msg.reader().readShort();
    string info = msg.reader().readUTF();
    GameCanvas.panel.saleRequest(type, info, id);
```

**Khuôn:** `[where byte][index short][text UTF]`. `where`/`index` phải **dội lại nguyên văn**
cái client vừa gửi — client giữ chúng để gói `7` vòng hai gửi lại. `text` là câu tiếng Việt
"Bạn có muốn bán\nxN <tên món>\nvới giá là G vàng?".

### 2.2 Sau khi bán thật

Java `sellItem:938-953`, theo thứ tự:

1. `sendMoney` → **gói `6`** (`[vàng long][ngọc int][ngọcKhóa int]`, 16 byte). Đây chính là
   gói tiền của đường mua — dùng lại `goiMuaBanTien`.
2. `sendThongBao` → **gói `-25`** "Đã bán X thu được N vàng" (dòng chữ chạy, không chặn màn hình).
3. `sendItemBag` → **gói `-36`** (túi mới), hoặc `sendItemBody` `-37` nếu bán đồ đang mặc.

`GameCanvas.endDlg()` nằm trong nhánh `case 6` phía client (`Controller.cs:3684`) → **gói `6`
phải luôn được gửi, kể cả khi bán thất bại**, đúng cùng lý do như mua: không gửi thì người
chơi kẹt hộp thoại chờ.

### 2.3 Kết luận khuôn phản hồi bản Go

| Trường hợp | Gửi |
|---|---|
| `action == 0` (xin xác nhận), món hợp lệ | `7` (`[where][index][câu hỏi]`) |
| `action == 0`, món không bán được | `-25` (lý do) → `6` (tiền không đổi) |
| `action != 0`, bán được | `-36` (túi mới) → `-25` "Đã bán X thu được N vàng" → `6` (tiền mới) |
| `action != 0`, từ chối | `-25` (lý do) → `6` (tiền **không đổi**) |

Giữ đúng thứ tự của đường mua (túi → chữ chạy → tiền) chứ không theo thứ tự Java
(tiền → chữ → túi): gói `6` phải là gói **cuối** vì nó đóng hộp thoại chờ. Java gửi tiền
trước được vì client cũ không mở hộp thoại chờ ở luồng bán — nhưng đặt nó cuối thì đúng
trong cả hai luồng, còn đặt đầu thì chỉ đúng một luồng.

---

## 3. Giá bán lấy từ đâu — QUYẾT ĐỊNH CỦA GATE

Java (`ShopService.java:922-932`, cả `showConfirmSellItem` lẫn `sellItem`, **hai bản sao
giống hệt nhau**):

```java
int quantity = item.quantity;          // bán CẢ CHỒNG
int cost     = item.template.gold;
cost /= 4;                             // 1/4 giá gốc (trừ mẫu 457 — thỏi vàng, bán nguyên giá)
if (cost == 0) cost = 1;               // sàn 1
cost *= quantity;
```

Bản Go **không có** trường giá trên `MauVatPham` (`mauvatpham.go:52-67`: Loại, GiớiTính, Tên,
MôTả, Cấp, YêuCầuSức, IconID, Phần, NângCấpĐược, TrầnChồng — không có `gold`). Thêm trường
vào đó là sửa file của agent khác.

**Quyết định bản 1:** giá bán lấy từ **món tương ứng trong cửa hàng ĐANG MỞ**:

```
đơn giá bán = max(1, MonHang.Gia / 4)     (chia nguyên, sàn 1)
loại tiền   = MonHang.Tien                (vàng hay ngọc — đúng loại đã mua)
tổng        = đơn giá bán × số lượng
```

Hệ quả phải nói thẳng ra: **chỉ bán được mẫu mà cửa hàng đang mở CÓ BÁN**. Món rơi từ quái
mà không cửa hàng nào bán thì bản 1 từ chối ("Cửa hàng không mua món này"). Đó là mất mát so
với bản Java, đổi lại không có con số giá nào bịa ra từ hư không. Khi KH8 (cân bằng kinh tế)
thêm giá gốc vào `MauVatPham`, đổi **một hàm** `banDonGia` là xong.

Mẫu 457 (thỏi vàng) và 570 (không bán được) là dữ liệu Dragon Ball — **không chép**.

---

## 4. Những chỗ bản Go CỐ Ý làm khác bản Java

| # | Java | Go bản 1 | Vì sao |
|---|---|---|---|
| 1 | Bán được cả đồ đang MẶC (`where == 0`) | **Chỉ bán từ TÚI**; `where == 0` bị từ chối | Bán đồ đang mặc còn phải dựng lại gói `-37` + `Send_Caitrang` + gỡ chỉ số cộng thêm — chưa có đường nào trong bản 1. Từ chối rõ ràng tốt hơn nửa vời. |
| 2 | Cộng tiền trong bộ nhớ rồi `subQuantityItemsBag` — **hai bước rời** | **MỘT** giao dịch CSDL (`store.BanMotGiaoDich`) | D2: đứt giữa chừng thì mất đồ mà không có tiền. |
| 3 | Không có khoá chống lặp | Khoá idem `<nvID>:BAN:<soThuTu>` | D3: bấm hai lần / gói lặp lại không được cộng tiền hai lần. |
| 4 | `item.template.gold` | `MonHang.Gia / 4` của cửa hàng đang mở | §3 trên. |
| 5 | Chặn `LIMIT_GOLD` (vàng vượt trần) | không có | Bản 1 chưa có trần vàng; `vang` là `bigint`, cộng vài nghìn không tràn. Ghi ra để KH8 nhớ. |
| 6 | Bán thì cả chồng, không kiểm gì | Hàm thuần nhận `soLuong`, **kiểm ≤ số đang có** trước khi chạm tiền | Handler luôn truyền cả chồng (đúng giao thức), nhưng chốt chặn nằm ở tầng dưới chứ không dựa vào người gọi. |

---

## 5. Điều CHƯA giải mã được — ghi ra để không ai tưởng đã xong

- `Panel.saleRequest(type, info, id)` dựng hộp thoại thế nào (cỡ chữ, nút) — không cần cho
  máy chủ, chỉ cần đúng ba ô.
- Cửa hàng nào gắn NPC nào: thuộc task mở NPC.
- Java còn `VatPhamDaBan.gI().addItem` ("mua lại đồ vừa bán") và mốc thành tựu theo
  `tagNameShop` — **không chép**, bản 1 không có hai màn hình đó.
