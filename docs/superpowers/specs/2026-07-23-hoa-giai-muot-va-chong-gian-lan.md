# Hoà giải: mượt tay và chống gian lận

**Thuộc về:** thiết kế gốc [`2026-07-23-naruto-tren-loi-nro-design.md`](2026-07-23-naruto-tren-loi-nro-design.md) §4.2 và §5.1
**Đích đến:** Global Constraints của **Kế hoạch 3**; một phần sang Kế hoạch 4 và 6
**Trạng thái:** đã chốt sau một vòng kiểm chéo độc lập

> **Lịch sử:** bản đầu của tài liệu này sai gần hết phần định lượng — bỏ sót
> `Time.timeScale = 1.5f`, nhầm vận tốc dọc thành vận tốc ngang, trích dẫn
> một dòng code không liên quan, và tự mâu thuẫn số học ở phần ước lượng tải.
> Bản này đã sửa. Mọi con số dưới đây đã được **hai bên đọc độc lập** rồi
> đối chiếu.

---

## 1. Vấn đề tưởng là một, hoá ra là hai

Ban đầu chúng tôi đặt vấn đề: *"người chơi bấm đánh phải thấy phản hồi ngay,
nhưng máy chủ không được tin client — hai điều này đá nhau."*

Đọc mã nguồn xong thì thấy **client đã tách sẵn hai chuyện này**, và giải
đúng một nửa:

| | Client làm gì | Đánh giá |
|---|---|---|
| **Chiến đấu** | Chạy hoạt ảnh ngay, nhưng **mọi con số đợi máy chủ** | ✅ Đã đúng, chỉ cần không phá |
| **Di chuyển** | Tự tính, tự vẽ, gửi **toạ độ tuyệt đối** lên máy chủ | ❌ Đây mới là lỗ hổng |

### 1.1 Chiến đấu đã đúng sẵn

Bấm chiêu → hoạt ảnh chạy ngay (`GameScr.cs:3764`), nhưng gói tin tấn công
chỉ gửi ở **khung hình cuối** của hoạt ảnh (`Char.cs:3145`). Gói đó chỉ chứa
**ID mục tiêu và hướng** (`Service.cs:1202-1274`) — không có số sát thương.
Sát thương, chí mạng, trúng/trượt, chết đều đọc từ gói máy chủ
(`Controller.cs:3961-4005`).

Đã rà cả **sáu** file client có `sendMessage` (`Service.cs`, `Session_ME.cs`,
`Session_ME2.cs`, `ISession.cs`, `QuaNapTuan.cs`, `QuayTamBao.cs`) và đối
chiếu ngược từ phía máy chủ (`Controller.java` không đọc hp/dame/def/crit/level
từ bất kỳ gói nào của client). **Không có phản ví dụ.**

**Ràng buộc:** máy chủ Go giữ nguyên phân chia này. Hoạt ảnh là việc của
client, con số là việc của máy chủ. Không cần cơ chế hoà giải nào.

### 1.2 Di chuyển là lỗ hổng thật

Client tự cộng `cx += cvx` mỗi nhịp, rồi gửi **toạ độ tuyệt đối**
(`Service.cs:886-889`). Máy chủ Java **gán thẳng, không kiểm**
(`PlayerService.java:174-175`).

Gói `checkMMove` (lệnh `-78`) đáng lẽ chống hack tốc độ, nhưng số giây trong
đó **do chính client tự đếm** (`GameScr.cs:4375-4384`), và máy chủ **đọc xong
vứt đi** (`Controller.java:815`).

*Ngoại lệ duy nhất:* có kiểm khoảng cách 500px, nhưng **chỉ trên map Hắc Long**
(`Controller.java:431-434`).

---

## 2. Ba giới hạn của client

### 2.1 Client không quay ngược được

Grep toàn bộ `Game1/` cho `predict|reconcil|rollback|interpolat`: không có gì
liên quan mạng. Không thể bảo client "tua lại". Công cụ duy nhất máy chủ có là
**kéo phịch về** (`Controller.cs:5035-5036`).

### 2.2 Client gửi theo ngưỡng khoảng cách, không theo nhịp

| Trạng thái | Gửi khi | Nguồn |
|---|---|---|
| Chạy dưới đất | `\|Δx\| + \|Δy\| ≥ 70` **và `cy − cySend ≤ 0`** | `Char.cs:2338` |
| Đang bay | `\|Δx\| > 96` **hoặc** `\|Δy\| > 24` | `Char.cs:4170` |
| Đổi trạng thái | ngay lập tức | `GameScr.cs:4232-4237` |

> Hai dòng đầu **khác đơn vị**: dưới đất là tổng Manhattan một ngưỡng, khi bay
> là hai ngưỡng độc lập theo từng trục. Đừng gộp khi tính toán.

> ⚠️ Vế `cy − cySend ≤ 0` nghĩa là **chạy xuống dốc thì không bao giờ gửi** —
> chỉ có gói lúc đổi trạng thái. Máy chủ phải chịu được khoảng lặng dài bất
> thường khi người chơi đi xuống, nếu không sẽ tính sai ngân sách và đá oan.

### 2.3 Ngưỡng nhảy hình — chỗ chật nhất

Client nội suy vị trí người khác `cx = (cx + đích) / 2` mỗi nhịp
(`Char.cs:2533`), hội tụ sau ~4 nhịp. Nhưng `moveTo` (`Char.cs:6218-6230`):

```csharp
if (type == 1 || Res.abs(toX - cx) > 100 || Res.abs(toY - cy) > 300)
{
    createShadow(cx, cy, 10);
    cx = toX;  cy = toY;
    vMovePoints.removeAllElements();   // <- xoá sạch hàng đợi
    ...
}
```

Ba điều quan trọng:
1. So với `cx` = **vị trí đang vẽ**, không phải mục tiêu nhận lần trước
2. Vượt ngưỡng thì **xoá sạch hàng đợi**, không phải chỉ nhảy một điểm
3. `moveTo` **xếp hàng** chứ không thay thế (`vMovePoints`, `Char.cs:6272`)

**Hệ quả — và đây là chỗ bản đầu của tài liệu sai:**

Client gửi ở nhịp đầu tiên vượt 96px, mỗi nhịp đi `s` px, nên mỗi gói mang
đúng `s · ⌈97/s⌉` px:

| `s` (px/nhịp) | Trạng thái | Delta mỗi gói | Vượt 100? |
|---|---|---|---|
| 5 | chạy, bay thường | **100** | không (đúng bằng) |
| **6** (`cspeed+1`, `GameScr.cs:2951`) | **bay ngang duy trì** | **102** | **CÓ** |
| 8 | có đồ tăng tốc | 104 | CÓ |

**Nghĩa là ở trạng thái bay ngang, người chơi khác nhìn thấy nhân vật nhảy
phịch trên MỌI gói tin — kể cả khi máy chủ phát lại hoàn hảo một-đổi-một.**
Quy tắc "đừng gom gói" là điều kiện *cần*, không phải *đủ*. Đây là lỗi có
sẵn trong game hiện tại, không phải thứ máy chủ gây ra.

Tệ hơn: `Session_ME.update()` rút **toàn bộ** hàng đợi tin trong một lần
(`Session_ME.cs:477-494`), và nó chạy trong `OnGUI()` (`Main.cs:175`) — tức
theo tốc độ **vẽ** (30/60 fps), không theo nhịp logic. Nên hai gói tới cùng
một khung hình thì gói thứ hai so với `cx` **chưa hề dịch chuyển** → 204px →
nhảy hình chắc chắn.

---

## 3. Các con số đo được

### 3.1 Nhịp thật là 75/giây, không phải 50

```csharp
// Main.cs:136 (trong Start())      →  ModFunc.cs:632 (trong LoadGame())
ModFunc.GI().LoadGame();               Time.timeScale = 1.5f;
```

`Time.fixedDeltaTime` = 0,02 (mặc định Unity, không bị đổi ở đâu — đã kiểm cả
`ProjectSettings/TimeManager.asset`). Nhưng đó là thời gian **game**.
`timeScale = 1.5` → `FixedUpdate` chạy **75 lần mỗi giây thực**.

| Đại lượng | Tính | Kết quả |
|---|---|---|
| Chạy (`cvx = cspeed = 5`, `GameScr.cs:2737`) | 5 × 75 | **375 px/giây** |
| Bay ngang tối đa (`cvx = cspeed + 1 = 6`, `GameScr.cs:2951`) | 6 × 75 | **450 px/giây** |
| `cBonusSpeed` | chỉ từng được gán 0 | không tác dụng |

> `cspeed = 5` là giá trị máy chủ gửi thật: `NPoint.java:76` → `Service.java:739`
> → `Controller.cs:1714`. Giá trị khởi tạo 4 ở `Char.cs:105` bị ghi đè lúc
> đăng nhập.

### 3.2 Hai lỗ hổng tốc độ nằm sẵn trong client

**(a) Lệnh chat `s <n>` — hack tốc độ có sẵn, không cần sửa gì**

```csharp
// ModFunc.cs:2039-2050, gọi từ GameScr.cs:6242
if (int.TryParse(strSpeed, out var speed) && speed > 0 && speed <= 10)
    Time.timeScale = speed;
```

Gõ `s 10` trong khung chat → 500 nhịp/giây → **3.000 px/giây khi bay**.

**(b) Alt-tab làm tốc độ mỗi người một khác**

```csharp
// Main.cs:148 / 152 — khi mất và lấy lại tiêu điểm cửa sổ
Time.timeScale = 0f;   ...   Time.timeScale = 1f;
```

`LoadGame()` không chạy lại, nên sau một lần alt-tab `timeScale` thành **1,0
vĩnh viễn**. Hai người chơi cùng cấp chạy 375 px/s và 250 px/s **tuỳ vào việc
họ đã alt-tab hay chưa**.

### 3.3 Công thức tăng tốc của bản Java là hàm bậc thang

```java
// NPoint.java:658, với tlSpeed là int (NPoint.java:214)
speed = (byte) (5 + 3 * (tlSpeed / 100));
```

Chia số nguyên → hàm bậc thang, **không phải luôn bằng 0**:

| `tlSpeed` | `speed` | px/giây (chạy, ở 75 nhịp/s) |
|---|---|---|
| 0–99 | 5 | 375 |
| 100–199 | **8** | **600** |
| 200–299 | 11 | 825 |

Bảng `options` có sẵn các dòng cộng vào `tlSpeed` với số lớn (id 146 "+100%
tốc độ chạy", id 145 "+66%", id 176 "+50%"). Chỉ cần tổng chạm 100 là tốc độ
**nhảy một phát từ 5 lên 8**.

Ngoài ra `EffectSkillService.java:405,414` đặt `speed = 1` / `speed = 8` cho
hiệu ứng làm chậm, nhưng **bị `calPoint()` xoá ngay** (`Service.java:723` →
`NPoint.java:460`) — chiêu làm chậm của bản cũ cũng hỏng. Riêng
`sendSpeedPlayer(player, 0)` (biến khỉ) thì đi thẳng và **có** hiệu lực.

> Bản Go tự đặt công thức mới, không port lại. Nhưng phải mô hình hoá được
> tốc độ **biến thiên theo hiệu ứng tạm thời**, không chỉ theo trang bị.

### 3.4 Kích thước thế giới

| Đại lượng | Giá trị | Nguồn |
|---|---|---|
| Bản đồ rộng nhất | 1920 px (map 160) | đọc nhị phân 176 file `tile_map_data/` |
| Bản đồ cao nhất | 1296 px (map 5, 29, 30, 89, 147) | như trên |
| Một ô tile | 24 px | `Map.java:50` |
| Người / khu vực | 15 (181 map), 30 (map 168, 172) | SQL `map_template` |
| Người / máy chủ | 2000 | `Manager.java:88` |
| Tầm nhìn | không có — gửi toàn khu vực | `Service.java:322-354` |

Khu vực nhỏ và ít người → **không cần chia tầm nhìn theo bán kính**.

---

## 4. Cách hoà giải

### 4.1 Sửa tận gốc ở client, thay vì bù ở máy chủ

Bản đầu của tài liệu định giải bài toán nhảy hình bằng cách bắt máy chủ **chia
nhỏ** mỗi cập nhật thành nhiều gói ≤50px. Cách đó **không chạy được**: `moveTo`
so với vị trí *đang vẽ* và các gói xếp hàng chứ không thay thế, nên ba gói tới
cùng một khung hình vẫn cộng dồn thành 100 → 200 → 300px và nhảy hình.

Vì game mới **có xây lại client** (đổi toàn bộ mỹ thuật, giao diện, cốt truyện),
sửa ba dòng ở client rẻ hơn nhiều so với mọi thủ thuật phía máy chủ:

| # | Sửa gì | Ở đâu | Được gì |
|---|---|---|---|
| **C1** | `Time.timeScale = 1.5f` → **`1.0f`** | `ModFunc.cs:632` | Tốc độ do máy chủ quyết hoàn toàn. **Đồng thời hết luôn lỗi alt-tab** (§3.2b), vì `Main.cs:152` vốn đã đặt lại `1f`. Bù lại bằng cách nâng chỉ số `speed` máy chủ gửi xuống, cảm giác chơi không đổi. |
| **C2** | **Xoá** lệnh chat `s <n>` | `ModFunc.cs:2039-2050`, `GameScr.cs:6242` | Bịt hack tốc độ có sẵn (§3.2a) |
| **C3** | Ngưỡng nhảy hình `> 100` → **`> 250`** | `Char.cs:6219` | Hết nhảy hình ở mọi tốc độ hợp lệ, kể cả khi hai gói dồn một khung. Dịch chuyển tức thời (5000px) vẫn nhảy đúng như phải thế. |

Sau ba sửa này, nhịp logic về **50/giây** và tốc độ tính thẳng bằng px/giây do
máy chủ đặt — không còn hệ số ẩn nào.

### 4.2 Nguyên tắc máy chủ: tin nhưng có kiểm, sai thì kéo về

Vì client không quay ngược được (§2.1), máy chủ **không mô phỏng lại chuyển
động người chơi**. Nó **thẩm định** báo cáo của client:

```
Mỗi người chơi có một "bình tín dụng di chuyển".
  · Bình tự hồi:  tốc độ trần CỦA CHÍNH NGƯỜI ĐÓ × 1,3
  · Sức chứa:     400 px

Nhận gói vị trí:
  chi phí = khoảng cách từ vị trí đã duyệt lần trước
  đủ bình  →  trừ bình, duyệt, phát lại cho cả khu vực
  hết bình →  KHÔNG duyệt, giữ vị trí cũ, kéo kẻ vi phạm về
```

**Vì sao bình tín dụng, không so từng gói?** Mạng giật thì gói dồn cục — so
từng gói sẽ đá oan. Bình cộng dồn theo thời gian nên chịu được dồn cục, nhưng
vẫn chặn kẻ đi nhanh **liên tục**.

**Vì sao trần riêng từng người, không phải một số chung?** Máy chủ tự đặt chỉ
số `speed` cho từng người nên nó **biết chính xác** ai được đi nhanh tới đâu.
Người mới không có đồ tăng tốc thì ngân sách chật, gian lận lộ ngay. Một trần
chung sẽ rộng bằng người mạnh nhất máy chủ — vô dụng với 99% còn lại.

### 4.3 Bảng số chốt cho Kế hoạch 3

| Thứ | Giá trị | Vì sao |
|---|---|---|
| Nhịp mô phỏng khu vực | **100 ms** | Khớp Java (`GameLoop.java:8`). Client không giả định nhịp nào (`lastUpdateTime` chỉ ghi, không ai đọc) → tự do chọn; giữ 100ms vì đã chạy tốt nhiều năm. |
| Nhịp này áp cho ai | **Chỉ quái, boss, vật phẩm rơi** | Vị trí người chơi đi theo sự kiện. Hai thứ khác nhau, đừng gộp. |
| Phát lại vị trí người chơi | **Tức thời, một-đổi-một, KHÔNG gom** | Gom là cộng dồn delta → vượt cả ngưỡng 250 mới. |
| Nhịp logic client (sau C1) | **50/giây** | `fixedDeltaTime` 0,02 với `timeScale` 1,0 |
| Tốc độ chạy gốc | **400 px/giây** (chỉ số `speed` = 8) | Giữ nguyên cảm giác hiện tại (375 px/s) sau khi bỏ hệ số 1,5 |
| Trần tốc độ thiết kế | **1000 px/giây** | 2,5× gốc. Bậc thang cũ chạm 825 px/s (§3.3) nên trần phải trên mức đó. |
| Bình tín dụng | **400 px**, hồi = trần riêng × 1,3 | ≈ 0,4–1 giây giật mạng tuỳ tốc độ người chơi |
| Vượt ngân sách | **Từ chối + kéo về** | Client không có lựa chọn khác (§2.1) |
| Chiêu dịch chuyển | **Cấp tín dụng riêng một lần** | `SkillService.java:674-694`, 5000px là hợp lệ. Máy chủ tự cấp khi **chính nó** xử lý chiêu, không phải khi client xin. |
| Client tự nhảy khi đánh | **Cấp tín dụng riêng** | `GameScr.cs:2050-2062` gán `cx = mục tiêu ± 60` rồi gửi ngay — hợp lệ, phải cho qua |
| Kiểm địa hình | **Có** — tra `tileMap`, chặn xuyên tường khi không bay | Java không kiểm. Tra một ô mảng, gần như miễn phí. |
| Tầm nhìn | **Toàn khu vực, không bán kính** | ≤1920px, ≤30 người. Chia bán kính là phức tạp hoá vô ích. |
| Thân gói `-7` từ client | **Dài thay đổi: 3 hoặc 5 byte** | `cy` là **tuỳ chọn** (`Service.cs:887`). Đọc sai là lệch cả luồng. |
| Không port lại | **Phát máu ra toàn khu vực mỗi lần bay** | `PlayerService.java:212-217` broadcast `Send_Info_NV` mỗi gói bay — gần gấp đôi lưu lượng vô ích |

### 4.4 Vá lỗ hổng §5.1 — kiểm tầm đánh

Bình tín dụng chặn *dịch chuyển tới sát boss*. Còn một đường nữa: đứng yên và
**đánh từ xa**.

Bản Java kiểm khoảng cách cho đấm thường (`SkillService.java:577-584`, 100px),
nhưng khối `case` đó **rơi xuyên** vào `KAMEJOKO/MASENKO/ANTOMIC` — ba chiêu
chưởng nhảy thẳng vào nhánh, bỏ qua kiểm. Makankosappo
(`SkillService.java:642-659`) không kiểm gì, mà tầm ghi trong dữ liệu là
**20000px** trong khi bản đồ rộng nhất 1920px → bắn khắp bản đồ.

Hai điều tài liệu bản đầu nói sai, đã sửa:
- Hằng số `RANGE_ATTACK_CHIEU_CHUONG = 300` **đã có sẵn** (`Skill.java:14`),
  chỉ là dùng cho boss (`Boss.java:595`) mà quên áp cho người chơi. Đây là
  **bật lại thứ có sẵn**, không phải thêm mới.
- **Nhiều chiêu diện rộng CÓ kiểm khoảng cách**: Thái Dương Hạ San
  (`SkillService.java:776, 792`), Tự phát nổ (`921, 935`), Trị thương 300px
  (`968`), Phân tâm 200px (`1082`), hệ chiêu mới dùng `skillSelect.dx`
  (`291-321`). Kế hoạch 4 cần biết cái nào đã có để không viết lại.

**Ràng buộc:** máy chủ Go kiểm khoảng cách cho **mọi** chiêu, lấy tầm từ dữ
liệu. Mọi tầm đánh trong dữ liệu game mới phải **nhỏ hơn chiều rộng bản đồ**.

---

## 5. Tải khi đông người

Kịch bản xấu nhất: một khu vực 30 người cùng bay tốc độ trần.

Ngưỡng gửi khi bay có **hai** vế, và vế dọc mới là vế nặng: bay lên có
`cvy = -7` (`GameScr.cs:2849`) → 350 px/s dọc, vượt ngưỡng 24px sau ~0,07 giây.

```
Bay lên:   350 ÷ 24  ≈ 15 gói/giây/người      (vế nặng)
Bay ngang: 450 ÷ 96  ≈ 4,7 gói/giây/người
Lấy 15    →  30 người  ≈ 450 gói/giây vào
Phát cho 29 người khác  ≈ 13.100 gói/giây ra
Gói -7 = 1 lệnh + 2 kích thước + 8 thân = 11 byte   (ServerManager.java:742-777)
                        →  ≈ 145 KB/giây cho một khu vực
```

Bản Java còn kèm một broadcast máu 17 byte mỗi gói bay
(`PlayerService.java:216`), đẩy con số lên gần gấp đôi. **Bản Go không port
lại chỗ này** — máu chỉ gửi khi thật sự đổi đáng kể.

Một goroutine sở hữu một khu vực thừa sức mức này. Phần tốn kém là **số lần
ghi socket** (Java gọi `flush()` mỗi gói, `ServerManager.java:775`), giải
quyết bằng gộp nhiều gói vào một lần ghi cho mỗi người nhận trong cùng vòng.

> Gộp **nhiều gói khác nhau vào một lần ghi socket**: nên làm.
> Gộp **nhiều lần cập nhật vị trí của cùng một người thành một gói**: cấm.
> Hai chuyện khác nhau.

---

## 6. Việc phải làm

- [ ] Đưa **bảng §4.3** vào Global Constraints của Kế hoạch 3
- [ ] Ba sửa đổi client **C1, C2, C3** (§4.1) — làm cùng đợt dựng lại client
- [ ] Ràng buộc kiểm tầm đánh mọi chiêu (§4.4) vào Kế hoạch 4, kèm danh sách chiêu **đã có** kiểm
- [ ] Kế hoạch 6: mọi tầm đánh trong dữ liệu < chiều rộng bản đồ
- [ ] Kế hoạch 6: `Manager.java:725` đọc `max_player` bằng `getByte()` — vỡ nếu >127 người/khu vực
- [ ] Kế hoạch 9 kế thừa các con số này, không tự đặt lại
