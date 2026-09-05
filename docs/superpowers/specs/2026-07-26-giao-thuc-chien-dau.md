# Giao thức chiến đấu — tham chiếu byte + thiết kế Kế hoạch 6

**Mục đích:** nguồn sự thật cho **Kế hoạch 6 (chiến đấu, chiêu thức, quái)**. Mọi bảng
byte dưới đây trích từ nghiên cứu
[`../research/kh6-goi-chien-dau.md`](../research/kh6-goi-chien-dau.md) (đọc ngược
`1-MAY-CHU-GAME/src/`). **Đừng viết lại từ trí nhớ — trích từ đây.**

**Phạm vi:** cho người chơi đánh được quái, quái đánh lại, người chơi đánh nhau, chết và
hồi sinh. Client Unity giữ nguyên (`D:\NRO_Project`), không sửa.

**Tiền đề đã có (KH1–KH5):** khung gói + mã hoá, đăng nhập, nhân vật, 4 catalog tĩnh,
vào bản đồ, actor khu vực nhịp 100 ms, di chuyển + chống dịch chuyển tức thời.

---

## 0. Ranh giới — KH6 làm gì, KHÔNG làm gì

Nghiên cứu cho thấy hệ chiến đấu của bản Java trải rộng gấp nhiều lần dự tính. Chia đôi:

### ✅ Trong KH6

| Nhóm | Nội dung |
|---|---|
| **A. Chiêu thức nền** | Điền catalog `-28/7` (hiện đang gửi 0 chiêu), gói chọn chiêu `34`, kiểm mana / hồi chiêu / tầm |
| **B. Quái** | Mô hình quái trong actor khu vực, gửi quái trong `-24`, hồi sinh 5 s, hồi máu 30 s |
| **C. Người đánh quái** | Gói lên `54`, gói xuống `54` / `-9` / `-12` |
| **D. Quái đánh người** | AI chọn mục tiêu 2 tầng, gói xuống `-11` / `-10` |
| **E. Chết & hồi sinh** | Gói `-17` / `-8` / `-16` / `-30/15` / `-30/9`, gói lên `-16` |
| **F. Người đánh người** | Gói lên `-60`, gói xuống `-60` (cả **hai dạng độ dài**) |

### ❌ Ngoài KH6 (ghi rõ để không ai tưởng bị bỏ quên)

| Bị hoãn | Lý do | Đi đâu |
|---|---|---|
| Rơi đồ khi quái chết | Cần hệ vật phẩm | KH7 — KH6 ghi `soDoRoi = 0` |
| Boss | Hệ riêng, 12 file Java, phức tạp hơn quái thường nhiều | KH6b hoặc KH9 |
| Hiệu ứng trạng thái đầy đủ (trói, tàng hình, khiên) | Chưa có bảng mã `loaiHieuUng` | KH6b |
| Chiêu `type = 4` (không cần nhắm) | Chưa biết gói lên nào mang `status == 20` | KH6b |
| Thú cưng / mobMe (`54` với `mobId == -1`) | Hệ thú cưng chưa có | KH7+ |
| PvP có luật (cờ `typePk`, khu PK) | Cần thiết kế luật chơi riêng | KH9 |
| **Gói `87` gửi kèm sau `-13`** | **KHÔNG BAO GIỜ làm** — bản Java dùng sai mã lệnh, xem §0.1 | Bỏ hẳn |

### 0.1 Gói `87` — bản Java gửi, nhưng bản Go KHÔNG gửi. Lý do:

Bản Java, sau khi phát `-13` (quái hồi sinh), phát tiếp `Message(87)` mang `byte mobId +
int maxHp` (`mob/Mob.java:525` gọi `sendMobMaxHp`, định nghĩa ở `mob/Mob.java:1480-1490`).
Nhìn qua tưởng là gói bắt buộc. **Không phải.** Đã mở client kiểm tận nơi:

```csharp
// Controller.cs:3016-3018
case 87:
    ((Mob)GameScr.vMob.elementAt(msg.reader().readUnsignedByte())).isWind = msg.reader().readBool();
    break;
```

**Mã `87` của client là lệnh "quái bị gió thổi", thân gói `byte mobId + bool isWind`
— hoàn toàn không liên quan tới máu.** Bản Java dùng **cùng một mã lệnh cho hai việc
khác nhau**: `sendWind` (`Mob.java:1468-1478`, khớp client, nhưng **không chỗ nào gọi**)
và `sendMobMaxHp` (`Mob.java:1480-1490`, **không khớp client**).

Ba lý do chốt **không gửi**:

1. **Client đã tự đặt `maxHp` rồi.** Nhánh `-13` ở `Controller.cs:4274-4300` có
   `mob15.hp = readIntToLongDQT(); mob15.maxHp = mob15.hp; mob15.updateHp_bar();`
   → gửi thêm `87` là **thừa tuyệt đối**.
2. **Gửi `87` gây hại thật.** Byte thứ hai client đọc thành `isWind` chính là **byte cao
   nhất của `int maxHp`** (big-endian). Quái có `maxHp ≥ 16.777.216` → byte đó khác 0 →
   `isWind = true` → `Mob.cs:1172` (`if (isDontMove || isWind)`) làm **con quái đứng im
   không di chuyển nữa**. Đây là bug thật đang có trong bản Java đang chạy.
3. **Muốn gửi maxHp riêng thì đã có mã đúng: `-84`** (`Controller.cs:767-783`,
   `byte mobIndex + int maxHp`). Nhưng tra khắp `1-MAY-CHU-GAME/src/` **không có chỗ nào
   tạo `Message(-84)`** — đường đúng mà server chưa bao giờ dùng. KH6 cũng chưa cần.

**Nhóm F (người đánh người) vẫn ở trong KH6** dù PvP có luật bị hoãn — vì gói `-60` là
gói *nguy hiểm nhất* của cả hệ (§2.1) và nên làm đúng ngay từ đầu, khi còn ít thứ phụ thuộc.

---

## 1. Ràng buộc cứng — chốt vào Global Constraints của kế hoạch

Sao chép nguyên khối này vào phần Global Constraints. Sai một dòng là hỏng client.

1. **Mọi trường máu và sát thương là `int32` big-endian (4 byte)**, kẹp trần
   `2_147_483_647`, kẹp sàn `0`. Không dùng 8 byte. *(Vì `cn.readInt = true`;
   `ModFunc.isReadInt = true` phía client — ràng buộc §11.1 của spec vào thế giới.)*
2. **Sát thương không bao giờ vượt máu còn lại của mục tiêu** — client hiện số âm nếu vượt.
3. **Quái đang đầy máu không chết bởi một đòn**: nếu `hp == hpMax` và `dame >= hp` thì
   `dame = hp - 1`.
4. **Mọi `count` / `size` / `mobId` / `tempId` đọc byte có dấu → ≤ 127.**
5. **Nhịp khu vực 100 ms** — giữ nguyên actor của KH5, không tạo luồng mới cho quái.
6. **Chỉ luồng `vongLap` của actor khu vực được chạm trạng thái quái và trạng thái
   chiến đấu.** Bao gồm cả *chiêu đang chọn* của người chơi (§3.2).
7. **Bất tử 1500 ms sau khi hồi sinh** — mọi đòn vào mục tiêu vừa hồi sinh bị từ chối.
8. **Hồi chiêu nới 50 ms**: cho phép dùng khi `daTroiQua >= coolDown - 50ms`.
9. **Vượt tầm đánh → `truot = true`, sát thương 0, NHƯNG VẪN PHÁT GÓI.** Không được im
   lặng bỏ qua — im lặng làm client thấy như lag.
10. **Dựng xong toàn bộ dữ liệu rồi mới ghi byte.** Lỗi giữa chừng → **bỏ cả gói**,
    không gửi gói thiếu byte. *(Bản Java nuốt lỗi giữa lúc ghi — xem §5.2 nghiên cứu.)*
11. **Tầm đánh:** cận chiến `100`, tầm xa `300` (đơn vị toạ độ, như bản gốc).

---

## 2. Ba gói ĐỘ DÀI BIẾN THIÊN — rủi ro số một của KH6

### 2.0 Sai trường thì hỏng thế nào — ĐỌC TRƯỚC, đừng chẩn đoán nhầm

Client đọc gói theo khung **có độ dài**: 1 byte lệnh → 2 byte độ dài → đọc **trọn** ngần
ấy byte vào một mảng riêng, rồi mới phân tích mảng đó
(`Session_ME.cs:181-219`). Hệ quả rất quan trọng:

> **Ghi sai trường bên trong một gói KHÔNG làm hỏng dòng byte của các gói sau.**
> Chỉ **gói đó** hỏng. Không có chuyện "client hỏng vĩnh viễn".

Ba kiểu sai và triệu chứng thật:

| Kiểu sai | Chuyện gì xảy ra | Triệu chứng người chơi thấy |
|---|---|---|
| **Ghi thiếu** (client đọc quá phần thân) | Client ném lỗi, bị `catch (Exception) {}` ở `Controller.cs:4600` nuốt | **Gói im lặng không làm gì** — quái không chết, sát thương không hiện |
| **Ghi thừa** (client đọc ít hơn) | Byte dư bị vứt lúc `msg.cleanup()` | **Không triệu chứng** — im lặng đúng, nhưng phí băng thông |
| **Sai thứ tự / sai kiểu, tổng độ dài vẫn khớp** | Client đọc giá trị rác vào đúng ô | **Số liệu vô lý** — máu âm, sát thương khổng lồ, quái đứng im |

Cả ba đều **im lặng, không có log** — đó mới là chỗ nguy hiểm thật: không ai biết mà sửa.

⚠️ **Một ngoại lệ cho hàng "ghi thừa": gói xuống `54`.** Gói này **không có trường đếm** —
client đọc `mobId` trong vòng lặp **cho tới khi hết thân gói** (tối đa 10 lần,
`Controller.cs:3851-3872`). Ghi thừa 1 byte ở đây **không im lặng**: client hiểu thành
"người này đánh 2 con quái", nhân vật quay sai hướng và nhắm sai mục tiêu.
→ **Với `54`, tuyệt đối không ghi thêm byte nào sau `mobId`.**

**Ngoại lệ duy nhất làm hỏng cả dòng byte** là ghi sai **số byte của trường độ dài**
(2 hay 3 byte) — đó là lỗi tầng khung, không phải tầng trường, và đã ghi ở §6.2 của
`BAT-DAU-TU-DAY.md`. Client dùng khung **3 byte** cho đúng 7 mã lệnh
(`Session_ME.cs:190`): `-32, -66, 11, -67, -74, -87, 66`. **Không mã chiến đấu nào nằm
trong danh sách này** → toàn bộ gói của KH6 dùng khung **2 byte**.

### 2.1 trở đi — ba gói biến thiên

Đọc kỹ trước khi viết dòng code nào.

### 2.1 `-60` có **hai dạng**, không phải một gói có trường tuỳ chọn

Byte thứ **11** (`docTiep`) quyết định gói dừng ở đó hay còn **7** byte nữa:

| Dạng | `docTiep` | Tổng thân | Dùng khi |
|---|---|---|---|
| **Đòn đánh** | `1` | **18 byte** | Gây sát thương lên người chơi |
| **Hồi máu** | `0` | **11 byte** | Chiêu hồi máu / hồi sinh đồng đội |

Cộng byte để đối chiếu (với `soMucTieu = 1`):
`4 + 1 + 1 + 4 + 1` = **11** (dừng ở `docTiep`), rồi `+ 1 + 4 + 1 + 1` = **18**.
Chênh lệch đúng **7 byte**.

**Yêu cầu thiết kế:** bản Go mô hình hoá thành **hai hàm dựng gói riêng biệt**
(`goiDanhNguoiChoi` và `goiHoiMauNguoiChoi`), **không** một hàm với tham số `docTiep`.
Lý do: một hàm với cờ thì người sửa sau rất dễ ghi thừa trường; hai hàm thì lỗi đó
không tồn tại được.

Ghi nhầm ở đây rơi vào ô "ghi thừa" của bảng §2.0: **im lặng, không triệu chứng gì** —
tức là sẽ không ai phát hiện cho tới khi có người ngồi so byte. Đó là lý do phải chặn
bằng test đột biến chứ không trông vào chạy thử.

### 2.2 `-45 / 0` — choáng diện rộng, hai vòng lặp lồng nhau

`byte 0, int idNguoiDung, short skillId, byte soQuai, {byte mobId, byte giay}×soQuai,
byte soNguoi, {int plId, byte giay}×soNguoi`

### 2.3 `-124` — hiệu ứng trạng thái, ba dạng

Byte 1 (`toggle`) và byte 2 (`loaiDoiTuong`) cùng quyết định phần còn lại. Id đối tượng
là **int** với người chơi, **byte** với quái.

> **Yêu cầu kiểm thử:** mỗi dạng của ba gói này phải có **một test byte-chính-xác riêng**,
> so với mảng byte ghim cứng. Không gộp chung. Và phải có **test đột biến**: đổi
> `docTiep` từ 0 thành 1 trong gói hồi máu → test **phải đỏ**. Không đỏ = test vô dụng
> (§6.11 của `BAT-DAU-TU-DAY.md` — dự án đã ba lần dính test xanh mà code sai).

---

## 3. Bảng byte

Quy ước: `byte`=1B, `short`=2B BE, `int`=4B BE, `bool`=1B (0/1), `UTF`=2B độ dài + chuỗi
modified-UTF8. Mọi `idNguoiChoi` là **int**; mọi `mobId` là **byte** (đánh số trong khu vực).

### 3.1 Gói LÊN — client → server

| Lệnh | Tên | Thân gói |
|---|---|---|
| `54` | Đánh quái | `byte mobId`; **nếu `mobId == -1`** thì đọc thêm `int masterId` |
| `-60` | Đánh người chơi | **`int idNguoiBiDanh` × N + `byte cdir`** ⚠️ *(sửa 22/08/2026 — xem ghi chú dưới bảng)* |
| `34` | Chọn chiêu | `short templateId` — **là `template.id`, KHÔNG phải `skillId`** |
| `-16` | Hồi sinh | *(rỗng)* |
| `-20` | Nhặt đồ | `short itemMapId` *(KH7)* |

> ### ⚠️ Đính chính gói LÊN `-60` (22/08/2026 — bằng chứng từ client thật)
>
> Bảng trên trước đây ghi thân gói `-60` chỉ có **một** `int idNguoiBiDanh`. **Sai.**
> Client Unity (`D:\NRO_Project\Assets\Scripts\Game1\Service.cs:1256-1265`) ghi một **vòng lặp** `writeInt(charID)` cho
> mọi mục tiêu đang chọn, rồi **luôn** ghi thêm `writeSByte(cdir)` ở cuối (dòng 1265,
> nằm ngoài cả ba nhánh `54`/`-60`).
>
> Server Go hiện đọc đúng **một** `r.Int()` (`login.go` `xuLyDanhNguoiChoi`) nên **vô tình
> đúng** ở trường hợp phổ biến (đánh 1 mục tiêu) và bỏ qua phần dư — không crash vì khung
> gói có độ dài riêng. Server Java gốc cũng chỉ đọc 1 int
> (`1-MAY-CHU-GAME/src/server/Controller.java:709`, `readByte()` bị chú thích).
>
> **Hệ quả cần nhớ:** khi làm KH6b (boss) hoặc PvP nhiều mục tiêu, muốn xử lý đủ danh sách
> thì phải sửa parser + thêm test byte cho gói lên. Trước khi có nhu cầu đó, giữ nguyên
> hành vi hiện tại nhưng **đừng tin dòng bảng cũ**.


⚠️ **`54` là gói độ dài biến thiên** (1 byte hoặc 5 byte). KH6 chưa làm thú cưng nhưng
**vẫn phải đọc đúng 5 byte** khi gặp `-1` rồi bỏ qua, nếu không sẽ lệch dòng.

### 3.2 Chọn chiêu — hai trường khác nhau, đừng lẫn

| Trường | Kiểu | Dùng ở |
|---|---|---|
| `template.id` | byte trong catalog, **short** trong gói `34` | Định danh *loại* chiêu (Hoả Độn, Lôi Độn…) |
| `skillId` | **short** trong catalog, **byte** trong gói `54`/`-60`, **short** trong `-45` | Định danh *một cấp cụ thể* của chiêu đó |

Client gửi `34` với `template.id`; server tra ra `Skill` (cấp người chơi đang có) rồi
mọi gói phát xuống dùng `skillId`.

**Quyết định kiến trúc:** *chiêu đang chọn* là trạng thái thuộc **actor khu vực**, không
thuộc luồng phiên. Lý do: luồng phiên nhận gói `34` và ghi, actor đọc khi tính sát
thương → tranh chấp dữ liệu, đúng bẫy §6.9 của `BAT-DAU-TU-DAY.md` (KH5 đã vấp một lần
với `nv.ViTriX/Y`). Luồng phiên **gửi thông điệp `ChonChieu` vào hộp thư actor**, không
ghi thẳng.

### 3.3 Gói XUỐNG — chiến đấu

| Lệnh | Tên | Thân gói | Gửi cho |
|---|---|---|---|
| `54` | Người đánh quái | `int idNguoiDanh, byte skillId, byte mobId` — ⚠️ **danh sách quái kết thúc bằng ĐỘ DÀI GÓI, không có byte đếm**; tuyệt đối không ghi thêm gì sau `mobId` (§2.0) | Cả khu vực |
| `-9` | Quái còn sống | `byte mobId, int hpConLai, int satThuong, bool chiMang, byte hieuUngId` | Cả khu vực |
| `-12` | Quái chết | `byte mobId, int satThuong, bool chiMang, byte soDoRoi, {short itemMapId, short templateId, short x, short y, int idNguoiNhat}×soDoRoi` | Cả khu vực |
| `-13` | Quái hồi sinh | `byte mobId, byte sys, byte levelBoss, int hp` — **KHÔNG gửi `87` theo sau**, xem §0.1 | Cả khu vực |
| `-11` | Quái đánh trúng **mình** | `byte mobId, int satThuong` | Riêng nạn nhân |
| `-10` | Quái đánh trúng **ai đó** | `byte mobId, int idNanNhan, int hpConLai` | Người khác trong khu vực |
| `-60` | Người đánh người | xem §2.1 | Cả khu vực |

⚠️ **Hai byte giữa của `-13` KHÔNG phải "loại quái" và "cấp quái"** dù bản Java đặt tên
thế. Client ghi byte 2 vào `mob.sys` và byte 3 vào `mob.levelBoss`
(`Controller.cs:4284-4288`) — **client không đổi loại quái khi hồi sinh**. Và
`levelBoss != 0` sẽ **bật hiệu ứng siêu cấp ngẫu nhiên**. → **Bản Go ghi `0` cho cả hai.**
Đây đúng kiểu lỗi "đặt đúng byte, sai ý nghĩa" đã bắt được ở `-17` (§7.2).

**Hai gói `-11` / `-10` còn một trường tuỳ chọn ở cuối** mà bản Java không gửi và bản Go
cũng không gửi (client bọc `try/catch`): `-11` có thêm `int satThuongChakra`, `-10` có
thêm `int mpConLai`. Ghi ra đây để ai đọc log gói của bản Java không bối rối.

**Ba chỗ của `-9` đã xác minh bằng client — bản Java sai, đừng chép:**

1. **Trường cuối là `byte hieuUngId` (1 byte), `-1` = không hiệu ứng.** Bản Java ghi
   `writeInt(-1)` = **4 byte** (`Mob.java:1369`); client chỉ đọc 1 byte
   (`Controller.cs:4352`), 3 byte sau bị vứt. Tình cờ vẫn ra `-1` nên không ai phát hiện.
   **Bản Go ghi 1 byte.**
2. **Sát thương đúng bằng `1` là giá trị đặc biệt:** client bỏ qua toàn bộ phần vẽ
   (`Controller.cs:4338`, `if (dame2 != 1.0)`). Tránh sinh ra sát thương bằng đúng 1.
3. **Sát thương `0` = "trượt":** client hiện chữ *miss* và **kéo quái về vị trí gốc** —
   nhưng **chỉ khi `chiMang = false`** (`Controller.cs:4356-4363`; nhánh chí mạng được
   kiểm trước và hiện "-0" màu đỏ thay vì *miss*).
   → **Ràng buộc bổ sung cho §1.9: gói vượt tầm phải ghi `chiMang = false`.**

Chi tiết `-60` dạng đòn đánh:

| # | Trường | Kiểu |
|---|---|---|
| 1 | `idNguoiDanh` | int |
| 2 | `skillId` | byte |
| 3 | `soMucTieu` | byte (**luôn 1** ở KH6) |
| 4 | `idNguoiTrung` | int |
| 5 | `docTiep` | byte = **1** |
| 6 | `loaiChieu` | byte = 0 |
| 7 | `satThuong` | int |
| 8 | `daChet` | bool |
| 9 | `chiMang` | bool |

Dạng hồi máu: các trường 1–4 giống hệt, trường 5 = **0**, **hết gói**.

### 3.4 Gói XUỐNG — chết và hồi sinh

| Lệnh | Tên | Thân gói | Gửi cho |
|---|---|---|---|
| `-17` | Bạn đã chết | `byte cPk, short x, short y` — **KHÔNG có id** | Riêng chính chủ |
| `-8` | Ai đó đã chết | `int idNguoiChoi, byte cPk, short x, short y` | Người khác |
| `-16` | Bạn đã hồi sinh | *(rỗng)* | Riêng chính chủ |
| `-30/15` | Ai đó hồi sinh | `int idNguoiChoi, int hp, int mp, short x, short y` | Cả khu vực |
| `-30/9` | Cập nhật máu | `int idNguoiChoi, int hp, int hpMax` | Người khác |

✅ **Đã xác minh bằng client** (§7). Hai điều dễ chép sai:
- `-17` **không có trường id** — byte đầu là `cPk`. Bản Java nhét id vào đó là **lỗi**.
- `-8` dùng `int` cho id; bản Java có một chỗ ghi `short` — **lỗi**, đừng chép.

### 3.5 Gói XUỐNG — hoạt ảnh chiêu (`-45`)

| sub | Ý nghĩa | Thân sau byte sub |
|---|---|---|
| `4` | Khởi chiêu | `int id, short skillId, short treMs` |
| `7` | Khởi bom | `int id, short skillId, short treMs` |
| `8` | Dùng chiêu | `int id, short skillId` |
| `6` / `1` / `5` / `3` | Bắt đầu gồng / đang gồng / nổ / dừng gồng | `int id, short skillId` (sub 3 ghi `short −1`) |
| `0` | Choáng diện rộng | xem §2.2 |

Sub `4` và `7` là mảnh ghép cho ràng buộc "client dự đoán trước, server xác nhận sau"
(§4.1 của [spec hoà giải mượt](2026-07-23-hoa-giai-muot-va-chong-gian-lan.md)): server
báo trước *bao nhiêu ms nữa chiêu nổ*, client chạy hoạt ảnh ngay, server chốt kết quả sau.
**KH6 chỉ cần sub `8`**; sub `4`/`7` để dành khi làm chiêu có thời gian khởi.

### 3.6 Danh sách quái trong `-24` MAP INFO

KH4 đã làm gói `-24` nhưng ghi `soQuai = 0` vì chưa có quái. KH6 điền vào. **Thứ tự
trường rất dễ chép sai** — `hp` đứng trước `level`, rồi mới tới `hpMax`:

| # | Trường | Kiểu |
|---|---|---|
| 0 | `soQuai` | byte |
| | *lặp `soQuai` lần:* | |
| 1–5 | `isDisable`, `isDontMove`, `isFire`, `isIce`, `isWind` | bool ×5 (KH6 ghi `false` cả 5) |
| 6 | **`templateId`** — mẫu quái (`MauId`/`anh_id`), **KHÔNG phải `tempId`** | byte |
| 7 | `sys` | byte = 0 |

> ⚠️ **Đính chính 05/09/2026:** bảng cũ ghi ô 6 là `tempId` và máy chủ đã ghi `Quai.Id` vào đó suốt từ KH6.
> Client (`Controller.cs:5068`) đọc ô 6 vào tham số `templateId` của `Mob(...)` rồi tra
> `Mob.arrMobTemplate[templateId]`; còn `mobId` (id tạm) client **tự lấy bằng chỉ số vòng lặp**. Với hai con
> sói cùng mẫu 0, con thứ hai (`Id=1`) làm client tra mẫu 1 không có → ngoại lệ → **bỏ nốt NPC, đồ rơi và 11
> bóng đổ** phía sau trong cùng gói. Hệ quả kéo theo: danh sách quái mỗi khu vực phải xếp đúng `Id = 0,1,2…`
> (client dùng chỉ số làm `mobId` cho `-9/-10/-11/-12`). Tìm ra bằng cách ghi hex `-24` thật ở mức Debug rồi
> giải mã đúng theo code client (`scratchpad/giai-ma-24.py`), không đoán.
| 8 | **`hp`** | int |
| 9 | **`level`** | byte |
| 10 | **`hpMax`** | int |
| 11 | `x` | short |
| 12 | `y` | short |
| 13 | `status` | byte |
| 14 | `lvMob` | byte |
| 15 | `isBigBoss` | bool (KH6 ghi `false`) |

### 3.7 Catalog chiêu `-28/7` — phần hiện đang rỗng

Khuôn một chiêu (`data/DataGame.java:150-209`), nối tiếp cấu trúc đã có ở
[spec vào thế giới §5](2026-07-24-giao-thuc-vao-the-gioi.md):

```
byte countClass = 5
  mỗi làng:
    UTF tenLang
    byte countTemplate                    ← KH6 đổi từ 0 thành số chiêu thật
      mỗi chiêu (SkillTemplate):
        byte   templateId
        UTF    ten
        byte   maxPoint                   ← số cấp tối đa của chiêu
        byte   manaUseType                ← 0 = tốn số chakra cố định, 1 = tốn % chakra tối đa
        byte   type                       ← 1 đánh, 2 buff người khác, 3 tự thân, 4 không nhắm
        short  iconId
        UTF    damInfo                    ← mô tả sát thương
        UTF    ghiChu                     ← ⚠️ bản Java ghi cứng "Ngoc Rong Kame" — TA GHI KHÁC
        byte   countSkill
          mỗi cấp (Skill):
            short skillId
            byte  point
            long  powRequire              ← 8 byte, KHÔNG phải int
            short manaUse
            int   coolDown                ← mili-giây
            short dx
            short dy
            byte  maxFight
            short damage
            short price
            UTF   moreInfo
```

⚠️ **Ba cái bẫy trong khuôn này:**

1. **`powRequire` là `long` 8 byte** — trường duy nhất trong toàn bộ hệ chiêu dùng 8 byte.
   Nhầm thành `int` là lệch 4 byte mỗi cấp chiêu.
2. **Chiêu có `templateId == 0` được đối xử đặc biệt** (`DataGame.java:179-208`): bản Java
   ghi `countSkill + 2` rồi nhồi thêm hai cấp giả `skillId` 105 và 106 toàn số 0. Đây là
   yêu cầu của client, không phải rác. **Nếu làng nào có chiêu `templateId = 0` thì phải
   nhồi y hệt.** Cách né sạch sẽ: **đánh số templateId của chiêu Naruto bắt đầu từ 1**,
   không dùng 0 — thì nhánh đặc biệt này không bao giờ chạy. **Spec này chọn cách né.**
3. **`ghiChu` là dấu vết thương hiệu cũ** — bản Java ghi cứng `"Ngoc Rong Kame"`. Ta ghi
   chuỗi rỗng hoặc tên làng. Đây là một trong những "dấu vết Dragon Ball xoá dần".

---

## 4. Chiêu thức Naruto — cấu trúc, không phải danh sách

### 4.1 Nguyên tắc: công thức trong code, nội dung trong dữ liệu

Hàm tính sát thương của bản Java (`player/NPoint.java:1715+`) dài hơn 200 dòng, là một
`switch` khổng lồ nhét cứng từng chiêu Dragon Ball trộn với nội tại, bộ đồ, bùa, thú cưng.
**Nó không phải công thức, nó là một đống ngoại lệ.**

Không chép. Lý do không chỉ là gọn: hàm đó **không viết test tử tế được** vì mỗi nhánh
phụ thuộc một mảng trạng thái toàn cục khác nhau — mà §6.11 của `BAT-DAU-TU-DAY.md` ghi
rõ dự án này đã **ba lần** dính "test xanh trong khi code sai".

Thay bằng:

```
satThuong = satThuongGoc(nhanVat) × heSoChieu(capChieu) × heSoChiMang × heSoHeNguyenTo
```

- `satThuongGoc` — từ chỉ số nhân vật.
- `heSoChieu` — đọc từ **dữ liệu chiêu** (`damage` trong catalog), không phải `if` trong code.
- `heSoChiMang` — 2× nếu chí mạng, ngược lại 1×.
- `heSoHeNguyenTo` — **mới của Naruto**, xem §4.2.

Mọi thứ đặc thù từng chiêu nằm ở **dữ liệu**, chuẩn bị sẵn cho KH8 (kho nội dung YAML +
nạp nóng). Code chiến đấu **không được biết tên chiêu nào cả**.

### 4.2 Hệ nguyên tố — thứ Dragon Ball không có

Thiết kế gốc đã gán mỗi làng một hệ ([thiết kế tổng §7.2](2026-07-23-naruto-tren-loi-nro-design.md)):
Lá–Hoả, Cát–Phong, Sương Mù–Thuỷ, Đá–Thổ, Mây–Lôi.

Vòng khắc chế đề xuất (đúng ngũ hành Naruto):

```
Hoả → Phong → Lôi → Thổ → Thuỷ → Hoả
```

(đọc: Hoả bị Thuỷ dập, Thuỷ bị Thổ hút, Thổ bị Lôi xuyên, Lôi bị Phong tán, Phong bị Hoả đốt)

| Quan hệ | Hệ số |
|---|---|
| Khắc chế | ×1,25 |
| Trung tính | ×1,00 |
| Bị khắc | ×0,80 |

**✅ ĐÃ CHỐT (27/07/2026, chủ dự án): giữ nguyên ×1,25 / ×1,00 / ×0,80.** Lý do phải chốt trong KH6 chứ không để sau:
hệ số này nằm trong công thức sát thương lõi; thêm sau nghĩa là viết lại và kiểm thử lại
toàn bộ tầng chiến đấu — đúng kiểu sai lầm mà §"phần phải quyết sớm" của KH10 cảnh báo.

⚠️ **Ràng buộc client:** hệ nguyên tố **không có chỗ nào trong giao thức** — client không
biết khái niệm này. Nó thuần tuý là phép nhân phía server. **Không cần sửa client**, và
đó là lý do có thể làm ngay trong KH6.

### 4.3 Tập chiêu tối thiểu cho KH6

Đề xuất **4 chiêu mỗi làng** (20 chiêu), phủ đủ **cả ba `type` nằm trong phạm vi KH6**:

| type | Vai trò | Ví dụ Làng Lá | Tầm | Gói phát xuống |
|---|---|---|---|---|
| 1 | Đòn cận chiến | Cước Liên Hoàn | 100 | `54` / `-60` dạng đòn đánh |
| 1 | Đòn tầm xa theo hệ | Hoả Độn · Hào Hoả Cầu | 300 | `54` / `-60` dạng đòn đánh |
| **2** | **Hồi máu đồng đội** | **Y Liệu Nhẫn Thuật** | **300** | **`-60` dạng hồi máu** |
| 3 | Tự thân | Tụ Chakra (hồi chakra) | — | *(không phát gói chiến đấu)* |

⚠️ **Chiêu `type = 2` là bắt buộc, không phải cho đẹp.** Gói `-60` **dạng hồi máu**
(§2.1) chỉ sinh ra từ chiêu `type = 2`. Không có chiêu nào type 2 thì hàm
`goiHoiMauNguoiChoi` **không có đường chạy thật** — test của nó chỉ là mô phỏng, và lỗi
sẽ chỉ lộ ra ở kế hoạch sau. Đây đúng loại bẫy "test xanh mà code sai" (§6.11 của
`BAT-DAU-TU-DAY.md`).

Bốn chiêu × 5 làng = 20 `SkillTemplate`, mỗi chiêu 3 cấp = 60 `Skill`. Vừa đủ nhỏ để
kiểm thử tay, vừa đủ rộng để lộ mọi lỗi khuôn byte.

*(`type = 4` — chiêu không cần nhắm — nằm ngoài KH6, xem §0.)*

**Mở rộng lên bộ đầy đủ là việc của KH8** (thêm chiêu = thêm dòng YAML, không sửa code).
Đó chính là lý do KH8 tồn tại.

---

## 5. Mô hình quái

### 5.1 Quái sống trong actor khu vực

`GameLoop.java:8` cho thấy server Java cũng chạy **nhịp 100 ms**, và `Zone.update()` gọi
`mob.update()` cho từng con — **AI quái chạy trong cùng luồng với khu vực**.

Bản Go đã có sẵn khung này từ KH5 (actor khu vực + ticker 100 ms, commit `27d66fc`).
Quái chỉ việc trở thành **một phần trạng thái của actor**, cập nhật trong `vongLap`.
**Không luồng mới, không khoá.** Đây đúng là luật §6.9 của `BAT-DAU-TU-DAY.md`.

### 5.2 Con số hành vi (giữ nguyên từ bản gốc)

| Việc | Con số | Nguồn |
|---|---|---|
| Hồi sinh sau khi chết | **5 giây** | `Mob.java:347` |
| Hồi máu khi chưa đầy | mỗi **30 giây**, hồi **hpMax/10** | `Mob.java:351-355` |
| Quái **trả thù** (đã bị đánh) | tầm **300**, nhịp đánh **1000 ms** | `Mob.java:423, 438` |
| Quái **tự gây sự** | tầm **100**, nhịp đánh **2000 ms** | `Mob.java:398, 412` |

Luật "trả thù nhanh gấp đôi, tầm gấp ba" là luật chơi tốt — **giữ nguyên**: người chơi
đi ngang qua thì tương đối an toàn, đã đánh thì phải chịu trận.

### 5.3 Sát thương quái

```
neu quai co dame rieng:  dame ± 1%
nguoc lai:               hpMax × pDame% ± nhieu theo cap
```
*(`mob/MobPoint.java:42-46`)*

Quái tinh anh (`lvMob > 0`) ở bản gốc bị ghi đè thành **10 % máu tối đa của người chơi**
(`Mob.java:457-459`) — tức luôn 10 đòn là chết bất kể trang bị. Giữ ý này nhưng đưa vào
**dữ liệu quái** (một cờ `satThuongTheoPhanTramMau`), không nhét `if` vào code.

### 5.4 Sáu lỗi của bản gốc — KHÔNG chép

1. **Kẹp máu sai biến** (`MobPoint.java:34-40`): kiểm giá trị **cũ** rồi mới gán giá trị
   mới → không bao giờ kẹp được sàn 0. Bản Go kẹp `[0, hpMax]` ở **một chỗ duy nhất**.
2. **Nhịp đánh chạy cả khi không có mục tiêu** (`Mob.java:386-390`): mốc thời gian đặt
   lại ngoài `if (player != null)` → quái vừa thấy người phải chờ thêm một nhịp.
   Bản Go: chỉ đặt lại mốc **khi thật sự đánh**.
3. **Thứ tự `&&`/`||` sai** (`Mob.java:402-404`): với big boss thì mọi bộ lọc phía trước
   bị bỏ qua, kể cả bộ lọc "chưa chết" → boss nhắm cả xác chết. Bản Go: ngoặc tường minh
   **và có test cho mục tiêu đã chết**.
4. **Đổi nhịp đánh như tác dụng phụ của hàm tìm mục tiêu** (`Mob.java:412, 438`).
   Bản Go: hàm tìm mục tiêu **trả về cả mục tiêu lẫn nhịp**, không sửa trạng thái ngầm.
5. **Hàm `lvMob()` trông như "đọc cấp quái" nhưng thật ra GHI hai trường**
   (`Mob.java:510-514`):
   ```java
   public int lvMob() {
       this.lvMob = 0;                          // ← đặt lại cấp
       this.point.hp = this.point.maxHp;        // ← HỒI ĐẦY MÁU
       return 0;                                // ← luôn trả 0
   }
   ```
   Nó được gọi ngay trong lúc **dựng gói** `-13` (`Mob.java:522`). Hai hệ quả:
   trường thứ 3 của gói `-13` **luôn bằng 0** trong bản gốc, và **bất cứ ai gọi
   `sendMobHoiSinh()` đều vô tình hồi đầy máu con quái**.
   ⚠️ Đừng vì thế mà tưởng byte đó vô nghĩa: **client CÓ dùng nó** — khác 0 thì bật hiệu
   ứng siêu cấp (`Controller.cs:4285-4288`). Bản Go ghi đúng `0`.
   Bản Go: hàm đọc thì **chỉ đọc**; hồi máu là việc riêng, gọi tường minh.
6. **Dùng lại một mã lệnh cho hai gói khác nghĩa** — `Message(87)` vừa là "gió thổi"
   vừa là "maxHp" (§0.1). Bản Go: **một mã lệnh, một ý nghĩa.**

---

## 6. Thứ tự gói — chuỗi đầy đủ một đòn đánh

### 6.1 Người chơi đánh quái, quái còn sống

```
Client → Máy chủ:  54     {byte mobId}
  [actor khu vực: kiểm bất tử/mana/hồi chiêu/tầm → tính sát thương]
Máy chủ → cả khu vực:  54   {idNguoiDanh, skillId, mobId}      ← hoạt ảnh
Máy chủ → cả khu vực:  -9   {mobId, hpConLai, satThuong, chiMang, −1}
```

### 6.2 Người chơi đánh quái, quái chết

```
Client → Máy chủ:  54
Máy chủ → cả khu vực:  54    ← hoạt ảnh
Máy chủ → cả khu vực:  -12   {mobId, satThuong, chiMang, soDoRoi=0}
  [5 giây sau, trong vòng lặp actor]
Máy chủ → cả khu vực:  -13   {mobId, tempId, lvMob, hp}
                             ← HET. Ban Java gui them goi 87 o day; ta KHONG gui (§0.1)
```

Client tự đặt `maxHp = hp` khi xử lý `-13` (`Controller.cs:4294-4295`), nên không cần
gói nào bổ sung máu tối đa.

### 6.3 Quái đánh người chơi

```
  [vòng lặp actor, mỗi 1000/2000 ms tuỳ tầng ưu tiên]
Máy chủ → nạn nhân:      -11  {mobId, satThuong}
Máy chủ → người khác:    -10  {mobId, idNanNhan, hpConLai}
```

### 6.4 Người chơi chết rồi hồi sinh

```
Máy chủ → chính chủ:   -17  {cPk=0, x, y}          ← KHONG co id
Máy chủ → người khác:  -8   {id, cPk=0, x, y}     ← id la int 4 byte
  [người chơi bấm hồi sinh]
Client → Máy chủ:      -16  (rỗng)
  [server trừ phí]
Máy chủ → chính chủ:   -16  (rỗng)
Máy chủ → cả khu vực:  -30/15  {id, hp, mp, x, y}
```

Ngoài ra, mỗi lần máu người chơi đổi vì bị đánh: `-30/9` cho những người khác.

---

## 7. ✅ ĐÃ XÁC MINH BẰNG CLIENT (26/07/2026) — hai chỗ nghi ngờ nay đã rõ

Đã mở source client Unity ở `D:\NRO_Project\Assets\Scripts\Game1\Controller.cs` và đọc
tận nơi. **Client là bên đọc, nên client là trọng tài duy nhất** — kết luận dưới đây
thay thế mọi suy đoán trước đó.

*(`Game2/Controller.cs` giống hệt `Game1`, chỉ khác `namespace`. Không cần soi riêng.)*

### 7.1 Gói `-8`: `writeInt` ĐÚNG — bản Java đang chạy có LỖI THẬT

`Controller.cs:4548-4554`:

```csharp
case -8:
    @char = GameScr.findCharInMap(msg.reader().readInt());   // ← readInt, 4 byte
    if (@char != null)
    {
        @char.cPk = msg.reader().readByte();
        @char.waitToDie(msg.reader().readShort(), msg.reader().readShort());
    }
    break;
```

→ `map/Zone.java:726` (`writeInt`) **đúng**. `services/Service.java:1157` (`writeShort`)
là **lỗi thật của bản Java đang chạy**: nó ghi 2 byte vào chỗ client đọc 4 byte, nên
client lấy nhầm id → `findCharInMap` trả `null` → **cả gói bị bỏ qua**. Hậu quả trong
game thật: **người chơi khác KHÔNG thấy ai đó vừa chết** (nhân vật vẫn đứng như thường),
cho tới khi họ rời khu vực rồi vào lại — lúc đó `Zone.java` mới phát bản `writeInt` đúng.

**Bản Go dùng `int` 4 byte.**

### 7.2 Gói `-17`: KHÔNG CÓ trường id — trường đầu là `cPk`

`Controller.cs:4533-4547`:

```csharp
case -17:
    Char.myCharz().meDead = true;
    Char.myCharz().cPk = msg.reader().readByte();                              // ← cPk
    Char.myCharz().startDie(msg.reader().readShort(), msg.reader().readShort());
    try {
        Char.myCharz().cPower = msg.reader().readLong();                       // ← TUỲ CHỌN
        Char.myCharz().applyCharLevelPercent();
    } catch (Exception) { Cout.println("Loi tai ME_DIE " + msg.command); }
    Char.myCharz().countKill = 0;
    break;
```

Gói `-17` chỉ gửi cho **chính chủ**, nên client dùng thẳng `Char.myCharz()` — **không cần
id, và không đọc id**. Byte đầu là **`cPk`** (cờ trạng thái PK).

→ `services/Service.java:1147` viết `writeByte((int) pl.id)` là **đặt nhầm dữ liệu vào
đúng ô**: nó nhét id nhân vật vào ô `cPk`. Không làm lệch byte (vẫn 1 byte), nhưng
**đặt sai giá trị `cPk` cho người vừa chết**. Dòng chú thích `// msg.writer().writeLong(-1);
send Power` bị comment ở `Service.java:1150` khớp đúng với trường `cPower` tuỳ chọn ở
client — client bọc `try/catch` nên thiếu cũng không sao.

**Bảng đúng cho bản Go:**

| # | Trường | Kiểu | Ghi chú |
|---|---|---|---|
| 1 | `cPk` | byte | **Cờ PK, KHÔNG phải id.** Ghi `0` khi chưa có hệ PK |
| 2 | `x` | short | |
| 3 | `y` | short | |
| 4 | `power` | long (8B) | **Tuỳ chọn** — client bọc `try/catch`. Bản Go **bỏ**, khớp bản Java |

### 7.3 Hệ quả cho kế hoạch

Task 12 (xác minh) **đã xong ngay trong đợt này** → **nhóm E không còn bị chặn**, làm
song song với A–D và F được. Xem Task 12 trong kế hoạch thi công.

---

## 8. Đối chiếu với những gì đã có

| Thứ | Trạng thái | Việc của KH6 |
|---|---|---|
| Actor khu vực, nhịp 100 ms | ✅ KH5 (`27d66fc`) | Thêm quái vào trạng thái actor |
| Gói `-24` MAP INFO | ✅ KH4 | Điền phần quái (§3.6), hiện ghi `soQuai = 0` |
| Catalog chiêu `-28/7` | ✅ KH3, nhưng **0 chiêu** | Điền 15 chiêu (§3.7, §4.3) |
| `w.Int` cho hp/hpMax | ✅ `xuathien.go:28-29` | Dùng lại y hệt cho sát thương |
| Kiểm khoảng cách di chuyển | ✅ KH5 `kiemDiChuyenHopLe` | Dùng lại cho kiểm **tầm đánh** |
| Chống gọi lại `-28/13`, `-39` | ✅ `bc3a55f` | Cùng kiểu chốt cho gói chiến đấu |

---

## 9. Phụ lục — tham chiếu mã nguồn đã kiểm chứng

| Sự kiện | Vị trí |
|---|---|
| `writeLongByGiaCatLuong` = writeInt | `network/Message.java:221-227` + `consts/cn.java:17` |
| Kẹp trần int32 | `utils/Util.java:153-161` |
| Điều phối gói đánh quái / đánh người | `server/Controller.java:697-715` |
| Gói `-60` dạng đòn đánh | `services/SkillService.java:1140-1149` |
| Gói `-60` dạng hồi máu (`docTiep=0`) | `services/SkillService.java:975-981` |
| Gói `54` phát xuống | `services/SkillService.java:1446-1460` |
| Gói `-9` / `-12` / `-13` | `mob/Mob.java:1361-1374` / `574-612` / `516-524` |
| Gói `-11` / `-10` | `mob/Mob.java:481-489` / `491-502` |
| Chốt kiểm trước khi đánh | `services/SkillService.java:58-136` |
| Hồi chiêu nới 50 ms | `services/SkillService.java:1273-1276` |
| Tầm đánh 100 / 300 | `skill/Skill.java:13-14` |
| Nhịp máy chủ 100 ms | `server/GameLoop.java:8` |
| Hồi sinh quái 5 s, hồi máu 30 s | `mob/Mob.java:347, 351-355` |
| AI quái hai tầng ưu tiên | `mob/Mob.java:393-443` |
| Khuôn catalog chiêu | `data/DataGame.java:141-216` |
| Danh sách quái trong `-24` | `map/Zone.java:779-800` |
| Bất tử 1500 ms sau hồi sinh | `services/SkillService.java:69-71` |

**Phía client Unity** (`D:\NRO_Project\Assets\Scripts\Game1\`) — đọc trực tiếp 26/07/2026:

| Sự kiện | Vị trí |
|---|---|
| Khung gói: đọc trọn theo độ dài rồi mới phân tích | `Session_ME.cs:181-219` |
| Danh sách 7 mã lệnh dùng khung **3 byte** | `Session_ME.cs:190` |
| `catch (Exception) {}` nuốt lỗi phân tích gói | `Controller.cs:4600` |
| `87` = **gió thổi** (`byte mobId + bool`), KHÔNG phải maxHp | `Controller.cs:3016-3018` |
| `-84` = maxHp thật (`byte mobIndex + int`), server chưa từng dùng | `Controller.cs:767-783` |
| `-13` client tự đặt `maxHp = hp` | `Controller.cs:4294-4295` |
| `-8` đọc id bằng **`readInt`** | `Controller.cs:4548-4554` |
| `-17` **không đọc id**, byte đầu là `cPk`; `power` là trường tuỳ chọn | `Controller.cs:4533-4547` |
| `-12` thứ tự đồ rơi: 4×short + int playerId | `Controller.cs:4392-4434` |
| `isWind` làm quái đứng im | `Mob.cs:1172` |
| `ModFunc.isReadInt = true` (khớp `cn.readInt` bên server) | `ModFunc.cs:45` |

*(`Game2/Controller.cs` giống hệt `Game1`, chỉ khác `namespace` — đã kiểm bằng diff.)*
