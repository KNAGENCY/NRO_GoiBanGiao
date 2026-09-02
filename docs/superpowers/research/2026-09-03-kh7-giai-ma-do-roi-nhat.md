# PHA A — Giải mã hai phía: đồ rơi dưới đất + nhặt (`68`, `-21`, `-19`, `-20`)

KH7 Task 11. Ngày 03/09/2026.

Nguồn: máy chủ Java `D:\NRO_GoiBanGiao\1-MAY-CHU-GAME\src` (chỉ đọc) · client Unity
`D:\NRO_Project\Assets\Scripts\Game2` (chỉ đọc).
Đặc tả: `docs/superpowers/specs/2026-09-03-kh7-vat-pham-design.md` §4.6, §4.6b (L2, L4),
§4.6c, §4.8, D8.

Mọi dòng dưới đây đều có **số dòng hai phía**. Chỗ nào không tra được thì ghi thẳng là
chưa biết — không đoán.

---

## 0. Bức tranh chung — bốn gói, không phải một

| Gói | Chiều | Ý nghĩa | Gửi cho ai |
|---|---|---|---|
| `68` | xuống | **Một món xuất hiện dưới đất** | cả khu vực |
| `-20` | lên | Xin nhặt món `itemMapId` | — |
| `-20` | xuống | **Bạn nhặt được** (kèm lời nhắn + số lượng) | riêng người nhặt |
| `-19` | xuống | **Người khác** vừa nhặt món này | cả khu vực TRỪ người nhặt |
| `-21` | xuống | Món biến mất (hết giờ) | cả khu vực |

**Phát hiện quan trọng nhất — client tự xoá món, máy chủ không cần gửi `-21` sau khi nhặt.**
Cả `-20` và `-19` đều gọi `itemMap.setPoint(...)` (`Controller.cs:3407`, `:3490`).
`ItemMap.setPoint` đặt `status = 2` (`ItemMap.cs:96-103`); `ItemMap.update()` thấy
`status == 2 && x == xEnd && y == yEnd` thì tự `vItemMap.removeElement(this)`
(`ItemMap.cs:105-115`). Nghĩa là món **bay tới nhân vật rồi tự biến mất** ở phía client.
`-21` chỉ dùng cho **hết giờ** — gửi `-21` sau khi nhặt là thừa (món đã bị xoá rồi, vòng lặp
tìm không thấy `itemMapID` nên `-21` chỉ không làm gì; vô hại nhưng vô nghĩa).

---

## 1. Gói XUỐNG `68` — một món xuất hiện dưới đất

Java `Service.java:1614-1627` (`dropItemMap`) ↔ client `Controller.cs:3511-3526`.

| # | Trường | Java ghi | Client đọc | Khớp? |
|---|---|---|---|---|
| 1 | `itemMapId` | `writeShort` (:1618) | `readShort()` (:3513) | ✅ |
| 2 | `mauId` (templateId) | `writeShort` (:1619) | `readShort()` (:3514) | ✅ |
| 3 | `x` | `writeShort` (:1620) | `readShort()` (:3515) | ✅ |
| 4 | `y` | `writeShort` (:1621) | `readShort()` (:3516) | ✅ |
| 5 | `chuUuTien` | `writeInt(3)` (:1622) | `readInt()` (:3517) | ✅ |
| 6 | `hàoQuang` | **chỉ khi** trường 5 == `-2` (`Service.java:2056-2057`) | **chỉ khi** trường 5 == `-2` (:3519-3522) | ✅ — **độ dài biến thiên** |

**Trường 5 là gì.** Client gán thẳng vào `ItemMap.playerId` (`ItemMap.cs:79-88`) và chỉ
dùng để vẽ/lọc phía client. Java ghi cứng `3` ở `dropItemMap` — con số `3` **không có ý
nghĩa nào tra được**; giá trị đặc biệt duy nhất client thật sự phân nhánh là `-2`
(= món có hào quang, đọc thêm 1 short bán kính, `Service.java:2042-2058` `dropSatellite`).

⚠️ **Bẫy §6.6 (một mã lệnh, nhiều dạng độ dài).** Bản Go: **hai hàm dựng riêng** —
`goiDoXuatHien` (không hào quang) và (nếu sau này cần) một hàm riêng cho `-2`. Không dùng cờ.

**Bản Go chốt:** trường 5 = **`-1`** nghĩa là "ai cũng nhặt được". Chọn `-1` thay vì `3` vì
`-1` chính là quy ước "không có chủ" mà máy chủ Java dùng ở chỗ kiểm quyền nhặt
(`Zone.java:450`: `itemMap.playerId == -1` → cho nhặt). Client không phân nhánh giá trị này
(chỉ `-2` mới đặc biệt) nên `-1` an toàn tuyệt đối. Ưu tiên chủ (đồ rơi của riêng ai) là
việc của **task 12**, chưa làm ở đây.

---

## 2. Gói LÊN `-20` — client xin nhặt

Client `Service.cs:1276-1291` (`pickItem`) ↔ Java `Controller.java:723-728`.

| # | Trường | Client ghi | Java đọc |
|---|---|---|---|
| 1 | `itemMapId` | `writeShort(itemMapId)` (:1282) | `readShort()` (:725) |

Thân gói **đúng 2 byte**, không có gì khác. Java còn chặn nhặt liên tục < 1000ms
(`ItemMapService.java:26-33`) và bỏ qua khi người chơi đang chết (`Controller.java:724`).

---

## 3. Gói XUỐNG `-20` — "bạn nhặt được" ⚠️ **chỗ có lỗi L2**

Java `Zone.java:462-500` ↔ client `Controller.cs:3395-3434`.

| # | Trường | Java ghi | Client đọc | Khớp? |
|---|---|---|---|---|
| 1 | `itemMapId` | `writeShort` (:464) | `readShort()` (:3399) | ✅ |
| 2 | `loiNhan` | `writeUTF` (:468/477/479/481/486) | `readUTF()` (:3408) | ✅ |
| 3 | `soLuong` | `writeShort` **một lần** (:498) | `readShort()` (:3412) | ✅ |
| 4 | `soLuongThat` | ❌ **KHÔNG ghi** | `readShort()` **thêm** nếu `template.type ∈ {9, 10, 34}` (:3413-3430) | ❌ **L2** |

**L2 nói bằng lời thường:** nhặt vàng (type 9), ngọc (type 10) hoặc lượng khoá (type 34),
client đọc **hai** short liên tiếp và **chỉ cộng tiền bằng short thứ hai**
(`Char.myCharz().xu += num` ở :3416 — `num` vừa bị gán lại bởi lần đọc thứ hai). Java chỉ
ghi một → client đọc hụt → `IndexOutOfRangeException` → bị `catch (Exception)` **rỗng**
(:3432-3434) nuốt sạch → **tiền không cộng, không lỗi, không log**.

**Bản Go chốt:** với `Loai ∈ {9, 10, 34}` ghi **ĐỦ HAI short**, cả hai bằng số lượng thật.
Với mọi loại khác ghi **đúng một** short.

**L4 — số lượng tràn short.** Java ghi cứng `9999` khi `quantity > Short.MAX_VALUE`
(:498) → người chơi thấy sai số. **Bản Go không chép:** `goiNhat` **trả lỗi** khi
`soLuong` ngoài `[0, 32767]` và không gửi gì cả (§1.10/L13 — thà không gửi còn hơn gửi số
sai). Chồng lớn hơn thì phải chia thành nhiều lần rơi — việc của task 12.

**Lời nhắn rỗng có ý nghĩa.** `loiNhan == ""` → client hiện **chữ bay `+N`** màu vàng
(type 9) / xanh (type 10) / đỏ (type 34) (:3435-3451). `loiNhan != ""` → hiện dòng chữ đó.
Bản Go: tiền để rỗng (chữ bay đẹp hơn), đồ thường ghi `"Bạn nhận được <tên món>"`.

---

## 4. Gói XUỐNG `-19` — người khác vừa nhặt

Java `Service.java:1404-1415` (`sendToAntherMePickItem`) ↔ client `Controller.cs:3475-3502`.

| # | Trường | Java ghi | Client đọc | Khớp? |
|---|---|---|---|---|
| 1 | `itemMapId` | `writeShort` (:1408) | `readShort()` (:3477) | ✅ |
| 2 | `idNguoiNhat` | `writeInt` (:1409) | `readInt()` → `findCharInMap` (:3478) | ✅ |

⚠️ Client `return` **ngay** nếu không tìm thấy nhân vật đó trong map (:3486-3489) —
người nhặt **phải còn trong khu vực** lúc gói tới, nếu không món **kẹt lại trên màn hình
người khác vĩnh viễn**. Vì vậy actor phải phát `-19` **ngay trong cùng nhịp** với lúc gỡ
món, chứ không phải sau một vòng chờ CSDL. Gửi cho **mọi người trừ người nhặt**
(`sendMessAnotherNotMeInMap`) — người nhặt đã có `-20` riêng rồi.

---

## 5. Gói XUỐNG `-21` — món biến mất (hết giờ)

Java `ItemMapService.java:41-51` ↔ client `Controller.cs:3381-3393`.

| # | Trường | Java ghi | Client đọc | Khớp? |
|---|---|---|---|---|
| 1 | `itemMapId` | `writeShort` (:45) | `readShort()` (:3384) | ✅ |

Client xoá thẳng khỏi `vItemMap`, **không hoạt ảnh**. Gửi cho cả khu vực.

---

## 6. Điều KHÔNG giải mã được / cố ý chưa làm

| Điều | Trạng thái |
|---|---|
| Ý nghĩa con số `3` ở trường 5 gói `68` | **Chưa tra được** trong client. Bản Go dùng `-1` (xem §1) |
| Danh sách đồ dưới đất khi **mới vào khu vực** | Java gửi từng gói `68` cho người mới vào — chưa tra dòng cụ thể. **Ngoài phạm vi task 11** (thuộc đường vào thế giới, task 5). Hệ quả: người vào sau hiện chưa thấy đồ rơi trước đó |
| Gói `-18` (tự vứt đồ xuống đất) | Client `Controller.cs:3503-3510` có đọc; thuộc **task 7** (`-43` vứt đồ), không làm ở đây |
| Ưu tiên chủ sở hữu đồ rơi (rơi từ quái ai giết) | Thuộc **task 12** |
| `type 34` (lượng khoá) | Bản 1 **chỉ có vàng (9) và ngọc (10)** — L10. Bộ dựng gói vẫn xử lý đúng 34 để sau này thêm không phải sửa |

---

## 7. Kết luận cho tầng Go

1. `68` = 5 trường cố định (`-1` ở trường 5), **không** ghi trường hào quang.
2. `-20` xuống: `short id`, `UTF lời nhắn`, `short soLuong`, **cộng thêm `short soLuong` lần
   hai khi `Loai ∈ {9,10,34}`**.
3. `-19`: `short id`, `int idNguoiNhat` — phát **cùng nhịp** với lúc gỡ món khỏi khu vực.
4. `-21`: `short id` — **chỉ** cho hết giờ, không dùng sau khi nhặt.
5. Số lượng ngoài `[0, 32767]` → **lỗi, không gửi gói**.
