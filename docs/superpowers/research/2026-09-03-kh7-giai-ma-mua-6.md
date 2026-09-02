# PHA A — Giải mã gói `6` (mua hàng NPC) và gói phản hồi

KH7 Task 9. Đọc trước khi viết một dòng code mua bán nào.

Nguồn (chỉ đọc):
- Client Unity: `D:\NRO_Project\Assets\Scripts\Game1\Service.cs`, `Controller.cs`, `Panel.cs`, `ModFunc.cs`
- Server Java cũ: `D:\NRO_GoiBanGiao\1-MAY-CHU-GAME\src\server\Controller.java`,
  `src\shop\ShopService.java`, `src\shop\Shop.java`, `src\services\Service.java`

---

## 1. Gói LÊN `6` — client xin mua

`Service.cs:958-971`:

```csharp
public void buyItem(sbyte type, int id, int quantity)
{
    message = new Message((sbyte)6);
    message.writer().writeByte(type);
    message.writer().writeShort(id);
    if (quantity > 1) { message.writer().writeShort(quantity); }
    session.sendMessage(message);
}
```

Server Java đọc y hệt (`Controller.java:318-339`):

```java
byte typeBuy   = reader.readByte();
int  tempId    = reader.readShort();
int  quantity  = reader.available() >= 2 ? reader.readShort() : 1;
if (quantity < 1) quantity = 1;
```

**Khuôn chốt:**

| Ô | Kiểu | Ý nghĩa |
|---|---|---|
| 1 | `byte` (có dấu) | `typeBuy` — **loại thao tác**, KHÔNG phải chỉ số tab |
| 2 | `short` (có dấu) | `tempId` — **id mẫu vật phẩm** (`MauID`, chính con số máy chủ đã gửi xuống trong gói `-44`) |
| 3 | `short`, **có thể vắng** | số lượng. Vắng ⇒ 1. `<1` ⇒ 1 |

### 1.1 `typeBuy` là gì — điểm dễ hiểu nhầm nhất

Nó **không** phải chỉ số tab. Các chỗ client gọi (`Panel.cs:13468-13521`, `ModFunc.cs:1977-1982`):

| Giá trị | Chỗ gọi | Nghĩa |
|---|---|---|
| `0` | `idAction == 3000` "mua đồ" — cửa hàng NPC thường | mua 1 món, không hỏi số lượng |
| `1` | `idAction == 3001` | luồng khác (có `pleasewait` dialog) |
| `2` | `idAction == 30003`, `ModFunc.cs:994` | luồng hộp/nhận tất |
| `3` | `idAction == 3004/3005`, `AutoBuyItem` | **mua nhiều cái** (có ô số lượng) |

Server Java **bỏ qua hoàn toàn** `typeBuy` cho cửa hàng NPC thường: `takeItem` phân nhánh
theo `tagName` của cửa hàng đang mở (trạng thái phía **máy chủ**), rồi `buyItem(player,
tempId)` tra món bằng `Shop.getItemShop(tempId)` — **duyệt hết mọi tab, so `is.temp.id ==
tempId`** (`Shop.java:56-65`).

⇒ **Bản Go làm y vậy:** món được tra bằng `MauID` trên **toàn bộ** cửa hàng đang mở, không
dùng `typeBuy` để định vị. `typeBuy` chỉ có ý nghĩa "kiểu thao tác" và bản 1 **chấp nhận
0 và 3** (hai luồng mua hàng NPC thật), các giá trị khác từ chối — vì chúng thuộc những
màn hình chưa cài (hộp may mắn, hộp thư).

### 1.2 Bẫy: máy chủ mới giữ "đang mở cửa hàng nào"

Client **không** gửi id cửa hàng trong gói `6`. Java lấy từ `player.iDMark.getShopOpen()`.
Bản Go phải có cùng thứ đó: phiên nhớ cửa hàng vừa mở qua gói `-44`. Không có cửa hàng
đang mở ⇒ từ chối. Đây cũng là chốt chặn chống "gửi gói `6` khống từ giữa đồng cỏ".

---

## 2. Gói XUỐNG — client mong nhận gì

Java `ShopService.buyItem` + `takeItem` gửi, theo đúng thứ tự:

1. **`-36`** — nạp lại cả túi (`InventoryService.sendItemBag`). Bộ dựng đã có: `GoiTui`.
2. **`-25`** — dòng chữ chạy "Mua thành công &lt;tên món&gt;" (`Service.java:1354-1367`,
   `writeUTF`). Client: `Controller.cs:2988-2990` → `GameScr.info1.addInfo(...)` — chữ
   bay, không chặn màn hình.
3. **`6`** — tiền (`Service.java:1387-1399`). **Cùng mã lệnh `6` nhưng chiều XUỐNG.**

Khuôn gói `6` xuống, client đọc ở `Controller.cs:3677-3685`:

```csharp
case 6:
    Char.myCharz().xu       = msg.reader().readLong();   // 8 byte — VÀNG
    Char.myCharz().luong    = msg.reader().readInt();    // 4 byte — NGỌC
    Char.myCharz().luongKhoa= msg.reader().readInt();    // 4 byte — ngọc khoá
    ...
    GameCanvas.endDlg();                                  // ĐÓNG hộp thoại chờ
```

**16 byte, không có ô nào khác.** Ba điều phải nhớ:

- Java ghi `writeLong` khi `session.version >= 214`, `writeInt` nếu không (`Service.java:
  1391-1394`). Client **luôn** đọc `readLong`. Đây đúng bẫy **L9** của đặc tả: bản Go
  **luôn ghi long**, không có khái niệm version.
- Ô thứ ba client gọi là `luongKhoa` (ngọc khoá), Java nhét `ruby` vào. Bản 1 chỉ có
  **vàng + ngọc** (bẫy **L10**) ⇒ ô thứ ba **luôn ghi 0**.
- `GameCanvas.endDlg()` nằm trong chính nhánh này. Luồng `typeBuy == 1` client bật
  `pleasewait` trước khi gửi (`Panel.cs:13496`). ⇒ **Gói `6` phải được gửi cả khi mua
  THẤT BẠI**, nếu không người chơi kẹt màn hình chờ. Java làm đúng vậy: `takeItem` gọi
  `sendMoney` ở cuối, **ngoài** mọi nhánh `return` của `buyItem`.

### 2.1 Kết luận khuôn phản hồi bản Go

| Trường hợp | Gửi |
|---|---|
| Mua được | `-36` (túi mới) → `-25` "Mua thành công X" → `6` (tiền mới) |
| Từ chối (thiếu tiền / túi đầy / gói rác) | `-25` (lý do) → `6` (tiền **không đổi**) |

Gói `6` luôn là gói **cuối cùng**: nó vừa cập nhật số tiền vừa đóng hộp thoại chờ.

---

## 3. Những chỗ bản Go CỐ Ý làm khác bản Java

| # | Java | Go bản 1 | Vì sao |
|---|---|---|---|
| 1 | Vòng `for (i=0; i<quantity; i++) takeItem(...)` — **mua N cái = N lần mua rời** | **MỘT** lần mua, một giao dịch CSDL, một khoá idempotency | Vòng lặp của Java là D3 dạng thuần khiết: đứt giữa chừng ở cái thứ 7 thì đã trừ tiền 7 lần và không ai biết. |
| 2 | `subMoneyByItemShop` kiểm số dư trong bộ nhớ Java rồi trừ trong bộ nhớ | `UPDATE ... WHERE vang + $delta >= 0` (BB4/D10) | Kiểm ở tầng Go không phải chốt chặn. |
| 3 | Kiểm "túi đầy" xong mới trừ tiền, nhưng cả hai ở bộ nhớ | Kiểm đầy đủ trước, rồi **một** giao dịch trừ tiền + cấp món | §2.3 / D16. |
| 4 | 4 loại tiền (vàng, ngọc, hồng ngọc, phiếu) | 2 loại | L10. |
| 5 | `typeBuy` bị bỏ qua hoàn toàn | Chỉ nhận `0` và `3` | Giá trị khác thuộc màn hình chưa cài; nhận bừa là mở cửa cho gói rác. |

---

## 4. Điều CHƯA giải mã được — ghi ra để không ai tưởng đã xong

- **Số lượng tối đa client cho nhập** ở ô "mua N": chưa tìm được chỗ chặn trên trong
  `Panel.cs`. Bản Go tự đặt trần = `TranChong` của mẫu (một ô túi chứa được bao nhiêu),
  vì đó là trần **có thật** ở phía dữ liệu.
- **Cửa hàng nào gắn NPC nào**: thuộc task mở NPC, không phải task này.
- `typeBuy == 1` và `== 2` dẫn tới luồng nào chính xác trong bản Java thì mã nguồn cũ
  trộn nhiều tính năng riêng của server private (hộp may mắn, hộp thư, "đá bán") —
  **không chép**, bản 1 không có những màn hình đó.
