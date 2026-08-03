# KH6b — Giải mã hệ BOSS (đọc code Java cũ + đối chiếu client Unity)

Nghiên cứu đọc mã nguồn server Java gốc tại `D:\NRO_GoiBanGiao\1-MAY-CHU-GAME\src\` và
client Unity tại `D:\NRO_Project\Assets\Scripts\Game1\`. **Chỉ đọc, không sửa.** Mọi số
dòng là số dòng thật tại thời điểm nghiên cứu (2026-08-03).

Tiếp nối [`kh6-goi-chien-dau.md`](kh6-goi-chien-dau.md) (quái thường) và spec
[`../specs/2026-07-26-giao-thuc-chien-dau.md`](../specs/2026-07-26-giao-thuc-chien-dau.md).
Quy ước kiểu Java (`writeByte`=1B, `writeShort`=2B BE, `writeInt`=4B BE, `writeUTF`=2B
độ dài + chuỗi) giữ nguyên.

Đánh dấu độ chắc chắn:
- ✅ **CHẮC CHẮN** — đọc được cả hai phía (server Java + client Unity), hoặc phía nào là bên đọc thì đã kiểm phía đó.
- 🔶 **SUY LUẬN** — mới đọc một phía (thường là server Java), chưa mở client xác nhận.

---

## 0. Phát hiện chi phối toàn bộ tài liệu — BOSS LÀ MỘT "PLAYER-AI", KHÔNG PHẢI QUÁI

Đây là điều quan trọng nhất, và nó **lật ngược trực giác**. Có **hai** khái niệm "boss"
hoàn toàn khác nhau trong code, đừng lẫn:

| Khái niệm | Lớp Java | Kế thừa | Bản chất |
|---|---|---|---|
| **Boss thật** (thứ ta cần làm) | `boss.Boss` | **`extends Player`** | Một **nhân vật do máy điều khiển** — có bộ đồ, chiêu thức, máu, tên, biết chat, tự đi lại. Với client nó **là một người chơi**. |
| **Quái to** (`isBigBoss`) | `mob.BigBoss` | `extends Mob` | Chỉ là con quái thường được đánh cờ `isBigBoss=true` trong gói `-24`. Đã nằm trong mô hình quái KH6. |

*(Nguồn: `boss/Boss.java:32` `public class Boss extends Player implements IBoss, IBossOutfit`;
`mob/BigBoss.java:10` `public class BigBoss extends Mob`.)*

**Hệ quả kiến trúc khổng lồ:** vì `Boss extends Player`, boss **không có một gói mạng
nào của riêng nó**. Nó tái sử dụng **toàn bộ giao thức người chơi**:

- Boss xuất hiện → gói **`-5`** (xuất hiện nhân vật), y như một người chơi vào khu vực.
- Boss đi lại → gói di chuyển nhân vật (KH5).
- Boss nói chuyện → gói **`44`** (chat bong bóng nhân vật).
- Boss đánh người → gói **`-60`** / **`-45`** (đúng gói người-đánh-người của KH6).
- Người đánh boss → client gửi **`-60`** (đánh người chơi) với `id` của boss.
- Boss chết / biến hình / rời đi → gói **`-6`** (nhân vật rời khu vực).
- **Chỉ có đúng một gói "riêng của boss": `93`** (thông báo chạy chữ toàn server) khi boss
  xuất hiện và khi boss bị giết. Nhưng gói `93` bản thân nó cũng không phải của riêng boss.

> **Kết luận cho KH6b:** phần lớn "hệ boss" **không phải là giao thức mới**, mà là một
> **máy trạng thái AI** điều khiển một thực thể giống-người-chơi (server-side actor), cộng
> với việc phát các gói người-chơi đã có. Việc khó của KH6b nằm ở **AI + vòng đời + dữ
> liệu**, không nằm ở byte. Điều kiện tiên quyết: server Go phải đã có "người chơi khác
> hiện ra trong khu vực" (gói `-5`/`-6` cho thực thể không phải chính chủ). Xem §10.

---

## 1. Boss khác quái thường ở đâu — bảng đối chiếu

| Khía cạnh | Quái thường (`mob.Mob`, KH6) | Boss (`boss.Boss`, KH6b) |
|---|---|---|
| Lớp cơ sở | `Mob` (thực thể riêng) | **`Player`** (giống người chơi) |
| Client nhìn thấy | Con quái (gói `-24`, `-9`, `-12`, `-13`) | **Một nhân vật** (gói `-5`, `-6`, `44`, `-60`) |
| Định danh | `mobId` = **byte**, đánh số trong khu vực | `id` = **int âm** (charID), toàn cục |
| Vòng đời | sống → chết → hồi sinh 5 s tại chỗ | **máy trạng thái 9 pha** (REST→…→LEAVE_MAP), tự chọn map |
| "Pha" theo máu | không | không — **pha = biến hình sau khi chết** (đọc §4) |
| Di chuyển | AI đơn giản, đứng tại điểm | **tự đi tới người chơi, né, đuổi** như người chơi |
| Tấn công | 1 công thức, nhịp 1000/2000 ms | **dùng chiêu thức thật** (`SkillService.useSkill`), nhịp 100 ms |
| Chat | không | **có** — chào lúc hiện, khích lúc đánh, trối lúc chết |
| Thông báo | không | **announce toàn server** (gói `93`) khi hiện + khi chết |
| Phần thưởng | trong gói `-12` (KH7) | **rơi đồ riêng + nhiệm vụ + điểm săn boss + điểm phúc lợi** |
| Nhịp cập nhật | actor khu vực 100 ms | **luồng riêng `BossManager` 150 ms** (không phải 100 ms!) |
| Cấu hình | bảng `mob_template` (DB) | **hardcode trong `BossesData.java`** + override file |
| Máu tối đa | vài nghìn–vài triệu | **200–500 triệu** (vẫn < trần int32 2,1 tỷ) |

---

## 2. Mười hai (+) file Java của hệ boss + vai trò

Thư mục `boss/` có **113 file**. Chia thành **lõi khung** (~12 file — thứ cần hiểu để
thiết kế) và **manifest** (~95 file boss cụ thể, mỗi con một class).

### 2.1 Lõi khung (đây là "12 file" mà sổ dự án nhắc)

| # | File | Vai trò |
|---|---|---|
| 1 | `boss/Boss.java` (859 dòng) | **Lớp cơ sở**. Máy trạng thái `update()`, vòng đời (rest/respawn/joinMap/chatS/active/attack/die/chatE/leaveMap), `injured()`, di chuyển, chat, gồng bom. Mọi boss kế thừa nó. |
| 2 | `boss/BossData.java` | **Mô hình dữ liệu MỘT pha** (POJO Lombok `@Data`): name, gender, outfit[6], dame, hp[], mapJoin[], skillTemp[][], textS/M/E, secondsRest, typeAppear, bossesAppearTogether[]. |
| 3 | `boss/BossesData.java` (~2400 dòng) | **"Cơ sở dữ liệu" boss hardcode**. Mỗi pha là một `static final BossData`. Đây là nơi chứa toàn bộ chỉ số/chiêu/lời thoại của mọi boss. |
| 4 | `boss/BossManager.java` (836 dòng) | **Bộ quản lý**: danh sách boss, `loadBoss()` (spawn lúc khởi động), `createBoss()` (factory theo id), **ticker 150 ms** gọi `boss.update()`, xuất `boss_status.json` cho web mỗi 1 s, xử lý lệnh panel (respawn/edit boss), menu `-96`. |
| 5 | `boss/BossID.java` | **Hằng số id** cho ~120 boss (đều **âm**: BROLY=-1, BLACK_GOKU=-203…). |
| 6 | `boss/BossStatus.java` | **enum máy trạng thái**: REST, RESPAWN, JOIN_MAP, CHAT_S, ACTIVE, DIE, CHAT_E, LEAVE_MAP, AFK. |
| 7 | `boss/BossType.java` | enum **loại quản lý** boss: FINAL, SKILLSUMMONED, PHOBAN (phó bản), event… quyết định boss vào manager nào. |
| 8 | `boss/AppearType.java` | enum **kiểu xuất hiện của một pha**: DEFAULT_APPEAR, APPEAR_WITH_ANOTHER, ANOTHER_LEVEL, CALL_BY_ANOTHER. |
| 9 | `boss/iboss/IBoss.java` | **Interface vòng đời** — liệt kê mọi bước máy trạng thái (rest/respawn/joinMap/chatS/active/attack/die/chatE/leaveMap/afk…). |
| 10 | `boss/iboss/IBossDie.java` | Interface cái chết (doSomeThing/notifyDie/rewards/leaveMap) — dùng cho một số biến thể. |
| 11 | `boss/iboss/IBossOutfit.java` | Interface trang phục (getHead/getBody/getLeg/getFlagBag/getAura/getEffFront). |
| 12 | `boss/BossManager.java` con: `OtherBossManager.java`, `FinalBossManager.java`, `SkillSummonedManager.java`, `BrolyManager.java` | **Các manager phụ** cho từng nhóm boss (phó bản, boss cuối, boss triệu hồi bằng chiêu, cụm Broly). Đều `extends BossManager` hoặc tương tự. |

*(Ngoài rìa nhưng liên quan: `mob/BigBoss.java` — quái to, KHÁC boss-AI, xem §0.)*

### 2.2 Manifest (~95 file) — mỗi boss một class

Nằm trong `boss/boss_manifest/<Nhóm>/<Ten>.java`. Mỗi class `extends Boss`, hàm dựng gọi
`super(BossID.X, …, BossesData.X_pha1, BossesData.X_pha2, …)` rồi **override** vài hành vi:
`reward()` (rơi đồ riêng), `attack()` (kiểu đánh riêng), `injured()` (giảm sát thương
riêng), `joinMap()`, `doneChatE()`… Ví dụ đã đọc:

- `boss_manifest/Cell/XenBoHung.java` — boss 3 pha, có chiêu **hấp thụ** người chơi.
- `boss_manifest/Black/BlackGoku.java` — 2 pha, pha ≥1 giảm nửa sát thương nhận, tự rời map.
- `boss_manifest/Cell/XENCON1.java` — "boss con" xuất hiện cùng boss mẹ (`parentBoss`).
- `boss_manifest/Training/TrainingBoss.java` (abstract), `boss_manifest/Bossvip/TestDameBoss.java`.
- Nhóm: `Cell/`, `Black/`, `Android/`, `Earth/`, `GinyuForce/`, `MajinBuu12H/`, `MajinBuu14H/`,
  `Bossvip/`, `Training/`, `DeathOrAliveArena/`, `Wolves/`, `Cumber/`, `Commeson/` (PhanThan, NhanBan — phân thân)…

**Thiết kế Naruto:** không cần 95 class. Nên **một class boss duy nhất, hành vi lái bằng
dữ liệu** — các "override" của boss cũ (giảm sát thương pha 2, rơi đồ, hấp thụ…) đưa thành
**cờ/tham số trong dữ liệu boss**, không đẻ ra class con. Đây đúng tinh thần "công thức
trong code, nội dung trong dữ liệu" của spec KH6 §4.1.

---

## 3. Vòng đời boss — máy trạng thái (chi tiết)

Trái tim hệ boss là `Boss.update()` (`boss/Boss.java:294-350`) chạy mỗi **150 ms** (do
`BossManager.start()` `boss/BossManager.java:804-818` dùng `scheduleAtFixedRate(…, 150ms)`).

```
        ┌──────── REST ────────┐   (nghỉ secondsRest giây; chỉ thoát khi pha kế = DEFAULT_APPEAR)
        │                      │
        ▼                      │
     RESPAWN ──► JOIN_MAP ──► CHAT_S ──► ACTIVE ──► DIE ──► CHAT_E ──► LEAVE_MAP
   (chọn map,   (đổi map,    (thoại    (đánh,       (bị     (thoại    (nếu còn pha →
    initBase)    announce)    chào)     chat khích)  giết)   trối)      RESPAWN pha+1;
                                                                        hết pha → REST)
```

Từng trạng thái (`boss/Boss.java`):

| Trạng thái | Dòng | Việc |
|---|---|---|
| **REST** | 352-362 | Nghỉ. `rest()` chỉ chuyển RESPAWN khi **pha kế tiếp là `DEFAULT_APPEAR`** và đã đủ `secondsRest` giây. → **Đây là "lịch hồi sinh".** |
| **RESPAWN** | 369-377 | `currentLevel++` (vòng lại 0 nếu vượt), `initBase()` — nạp `BossData[currentLevel]`: tên, giới, outfit, dame, hp (random trong `hp[]`), **nạp lại skill**, chuyển NON_PK. |
| **JOIN_MAP** | 379-466 | Chọn khu vực (né khu >10 người, né khu đã có boss), `changeMap` (phát gói `-5` cho cả zone), **`notifyJoinMap()` → announce `93`**, đánh thức boss đi kèm, chuyển CHAT_S. |
| **CHAT_S** | 494-514 | Đọc lần lượt `textS[]`, mỗi câu phát gói `44`, chờ `len×100 ms` (trần 2000). Xong → ACTIVE. |
| **ACTIVE** | 541-547 | `chatM()` (thỉnh thoảng khích tướng, gói `44`) + `attack()`. Chuyển sang PK_ALL để ai cũng đánh được nó / nó đánh được mọi người. |
| **attack** | 551-583 | Mỗi 100 ms: chọn mục tiêu ngẫu nhiên, chọn chiêu ngẫu nhiên trong bộ, nếu trong tầm → `SkillService.useSkill(this, pl,…)` (phát `-60`/`-45`); nếu ngoài tầm → `moveToPlayer` (phát gói di chuyển). |
| **DIE** | 339-340 | Chỉ chuyển thẳng sang CHAT_E. (Cái chết thật xảy ra trong `injured()` → `die()`.) |
| **CHAT_E** | 638-657 | Đọc `textE[]` (lời trối), phát `44`. Xong → LEAVE_MAP. |
| **LEAVE_MAP** | 665-676 | **Nếu còn pha** (`currentLevel < data.length-1`): lưu zone, quay lại RESPAWN (biến hình pha kế). **Hết pha**: `exitMap` (gói `-6`), đặt `lastTimeRest`, về REST. |
| **AFK** | 365-367 | Rỗng (treo). |

Điểm mấu chốt về `die()` (`boss/Boss.java:604-610`): gọi khi máu về 0 trong `injured()`.
Nó **announce toàn server** (`ServerNotify.notify(plKill.name + " đã tiêu diệt " + name)`),
gọi `reward(plKill)`, rồi chuyển trạng thái DIE. **Cái chết boss không dùng gói `-12`** như
quái — nó đi qua CHAT_E rồi `-6` (rời map).

---

## 4. Hệ "PHA" — không phải theo mốc máu, mà là BIẾN HÌNH sau khi chết

Đây là điểm dễ hiểu sai nhất. Câu hỏi ban đầu "boss đổi hành vi khi máu xuống mốc X%?" —
**KHÔNG. Bản gốc không có phase theo % máu.** Cơ chế thật:

- Một boss mang **mảng `BossData[] data`** — mỗi phần tử là **một pha/dạng** (tên, HP,
  chiêu, outfit khác nhau). `currentLevel` trỏ pha hiện tại.
- Boss xuất hiện ở `data[0]`. Người chơi giết được → boss **chết pha 0** → `LEAVE_MAP`
  thấy còn pha → **RESPAWN pha 1** (mạnh hơn, outfit khác, thoại khác) → giết tiếp → pha 2…
- Hết pha cuối → về REST, chờ `secondsRest`, rồi lặp lại từ pha 0.

**Ví dụ thật — Xên Bọ Hung** (`XenBoHung.java:30` + `BossesData.java:933-1040`):

| Pha | Tên | HP | outfit | typeAppear |
|---|---|---|---|---|
| 0 | Xên Hoàn Thiện | 200.000.000 | 228/229/230 | DEFAULT_APPEAR |
| 1 | Xên Bọ Hung 1 | 300.000.000 | 231/232/233 | ANOTHER_LEVEL |
| 2 | Xên Hoàn Thiện | 400.000.000 | 234/235/236 | ANOTHER_LEVEL |

`ANOTHER_LEVEL` = pha này **chỉ tồn tại như bước tiến hoá của cùng con boss**, không tự
respawn độc lập (vì `rest()` chỉ nhả khi pha kế là `DEFAULT_APPEAR`).

**Bốn `AppearType`:**
- `DEFAULT_APPEAR` — pha tự xuất hiện theo lịch (pha 0 của đa số boss).
- `ANOTHER_LEVEL` — pha biến hình (pha ≥1 của boss nhiều dạng).
- `APPEAR_WITH_ANOTHER` — pha xuất hiện **kéo theo boss khác cùng lúc** (cụm boss, ví dụ đội Ginyu).
- `CALL_BY_ANOTHER` — pha chỉ xuất hiện khi **boss khác gọi** (boss con `parentBoss`).

**Boss đi theo cụm** (`bossAppearTogether[][]`, `Boss.java:105-118, 766-788`): một boss có
thể sinh ra một mảng boss con/đồng đội, liên kết `parentBoss`. Ví dụ Xên mẹ + 7 Xên con
(`XENCON1.doneChatE` chờ mọi Xên con chết mới cho mẹ ACTIVE). Đây là **cơ chế "boss triệu
hồi lính"** — rất hợp làm boss Naruto (Kage triệu hồi thú, Pain 6 thân…).

🔶 **SUY LUẬN:** toàn bộ §3-4 đọc từ server Java. Client không có khái niệm "pha" — với
client mỗi pha chỉ là một nhân vật mới (id cũ, outfit mới) xuất hiện qua `-5`. Không cần
gói mới. Nhưng **client dùng `id` để phân biệt nhân vật**, nên khi biến hình boss giữ
nguyên `id`; cần kiểm cách client xử lý một `id` cũ hiện lại với outfit mới (§11 bẫy).

---

## 5. Giao thức mạng boss — bảng byte, đối chiếu hai phía

Nhắc lại: boss **không có gói riêng** trừ `93`. Dưới đây là **các gói người-chơi mà boss
tái dùng**, đã lọc ra đúng phần liên quan boss.

### 5.1 `-5` — boss XUẤT HIỆN (client vẽ nhân vật) ✅

Boss vào khu vực qua `ChangeMapService.changeMap` → phát `-5` như một người chơi. Server Go
**đã có** hàm này: `server/internal/game/xuathien.go:15` `goiXuatHien()` (28 trường).

| Chi tiết | Nguồn |
|---|---|
| Layout 28 trường (playerId int, clanId, level, cờ, typePk, gender×2, head, name UTF, **hp int, hpMax int**, body, leg, flagbag, −1, x, y, buff…, aura, effFront, hat) | `server/internal/game/xuathien.go:15-48` (đã soi chéo `Zone.infoPlayer` bản Java) |

**Việc boss cần thêm vào `-5` hiện có:** server Go đang ghi `0`/hằng cho outfit (head/body/
leg/aura/effFront) vì người chơi KH5 chưa có trang bị. **Boss thì bắt buộc có outfit thật**
(từ `BossData.outfit[6]`) để nhìn ra hình dạng. → KH6b phải điền các trường outfit này khi
thực thể là boss. `typePk` cũng cần = PK_ALL (đánh nhau) thay vì 0.

### 5.2 `-6` — boss BIẾN MẤT (chết / biến hình / rời map) ✅

`services/func/ChangeMapService.java:834`:
```
msg = new Message(-6);
msg.writer().writeInt((int) player.id);   // chi 1 truong: id nhan vat
sendMessAnotherNotMeInMap(player, msg);
```
Client xoá nhân vật khỏi khu vực. Boss dùng khi `exitMap` (rời hẳn) và khi biến hình
(rời để hiện lại pha mới bằng `-5`).

| # | Trường | Kiểu |
|---|---|---|
| 1 | `id` | int |

### 5.3 `44` — boss CHAT (bong bóng thoại) ✅

`services/Service.java:596-607` — dùng cho chatS/chatM/chatE của boss.
```
msg = new Message(44);
msg.writer().writeInt((int) player.id);
msg.writer().writeUTF(text);
sendMessAllPlayerInMap(player, msg);
```
Client đọc (`Game1/Controller.cs:4567-4572`):
```csharp
case 44:
    int id = msg.reader().readInt();
    string text = msg.reader().readUTF();
    (Char.myCharz().charID != id ? GameScr.findCharInMap(id) : Char.myCharz())?.addInfo(text);
    break;
```
✅ **Khớp.** Lưu ý client `findCharInMap(id)` — **id không có char tương ứng thì chat rơi
im lặng** (củng cố bẫy §11.2: id boss phải là một char hợp lệ, không trùng).

| # | Trường | Kiểu |
|---|---|---|
| 1 | `id` | int |
| 2 | `text` | UTF |

Lời thoại boss có **tiền tố định tuyến** `|prefix|text` (`Boss.java:731-764`):
`-1` = boss tự nói · `-2` = một người chơi trong map "bị bắt nói" · `-3` = boss mẹ nói ·
`>=0` = boss thứ N trong cụm nói. **Chỉ là logic server chọn ai phát gói `44`** — client
không thấy prefix.

### 5.4 `93` — ANNOUNCE TOÀN SERVER (gói duy nhất "của boss") ✅

`server/ServerNotify.java:71-80`. Boss gọi `ServerNotify.gI().notify(text)` (đẩy vào hàng
đợi), một thread phát mỗi 1,5 s bằng gói `93` tới **TẤT CẢ người chơi** (`sendMessAllPlayer`,
không chỉ khu vực).

Server ghi:
```
msg = new Message(93);
msg.writer().writeUTF(text);
Service.gI().sendMessAllPlayer(msg);
```
Client đọc (`Game1/Assets.src.f/Controller2.cs:825-829`):
```csharp
case 93:
    string str = msg.reader().readUTF();
    str = Res.changeString(str);
    GameScr.gI().chatVip(str);   // dong chu chay
    break;
```
✅ **Khớp hoàn toàn.** Thân gói = **1 UTF**.

| # | Trường | Kiểu |
|---|---|---|
| 1 | `text` | UTF |

Boss phát `93` ở hai thời điểm:
- **Xuất hiện** (`Boss.notifyJoinMap()` `Boss.java:477-481`): `"BOSS <ten> vừa xuất hiện tại <map>"`
  — nhưng **có bộ lọc** `canSendNotify()` (`Boss.java:483-492`): tắt cho boss phó bản, map
  ma bư, offline, một số mapId. Cờ `isNotifyDisabled` cũng tắt được.
- **Bị giết** (`Boss.die()` `Boss.java:607`): `"<tên người>: Đã tiêu diệt được <boss>, mọi người đều ngưỡng mộ."`

### 5.5 `-60` / `-45` — boss ĐÁNH người & người ĐÁNH boss ✅ (tái dùng nguyên KH6)

Vì boss là Player, **không có gì mới**:
- **Boss đánh người:** `Boss.attack()` gọi `SkillService.useSkill(boss, pl,…)` → phát
  `-60` (dạng đòn đánh) hoặc `-45/4,7` (khởi chiêu) **y hệt người chơi đánh người chơi**
  (spec KH6 §2.1, §3.5). Boss được truyền vào chỗ `plAtt`.
- **Người đánh boss:** client gửi **`-60`** (đánh người chơi) với `int id` = **id của boss**
  (số âm). Server `attackPlayer` → `playerAttackPlayer` (`SkillService.java:1113`) →
  `boss.injured(...)` → phát `-60` xuống. Client cập nhật máu boss như máu một nhân vật.

⚠️ Có vài nhánh xử lý riêng khi mục tiêu là boss (`plInjure.isBoss`), **chỉ đổi con số, không
đổi byte**: `SkillService.java:1120-1122` (skin xDame chia 3 khi đánh boss),
`:936` (chia 2)… Đây là **cân bằng**, không phải giao thức.

### 5.6 `-96` / `108` / `-118` — menu DANH SÁCH BOSS (phụ trợ, admin) 🔶

Tính năng "xem danh sách boss / dịch tới boss", không thuộc vòng đời chiến đấu:
- **`108`** (lên): client xin danh sách → `BossManager.showListBoss(player, true)`.
- **`-96`** (xuống): danh sách boss (`BossManager.java:712-774`). ✅ Client đọc ở
  `Game1/Controller.cs:233-264` (dùng chung khung "bảng xếp hạng" `vTop`): byte type, UTF
  tên bảng, byte đếm, lặp {int rank, int pId, short head, short headIcon(nếu ver≥214),
  short body, short leg, UTF tên, UTF info, UTF info2}.
- **`-118`** (lên): client bấm một boss → gửi `int index`; server (menuType 3, **chỉ admin**)
  `getBossByIndex` rồi **dịch chuyển người chơi tới chỗ boss** (`Controller.java:765-787`).
  Không phải "thách đấu spawn boss riêng" — chỉ là công cụ admin.

→ **KH6b có thể bỏ qua menu này**, làm sau. Vòng đời boss không cần nó.

### 5.7 Gói boss KHÔNG dùng (để khỏi nhầm)

Boss **không** phát `-9`, `-12`, `-13`, `-24`, `-84`, `87` (đó là gói **quái**). Nếu ai
định "gửi máu boss như gửi máu quái" là sai hướng — máu boss đi qua gói người chơi (`-60`,
`-5` mang hp/hpMax, và `-30/9` cập nhật máu nếu cần).

---

## 6. Chiêu ĐẶC BIỆT của boss (khác đòn đánh thường)

Ngoài đánh chiêu thường, boss có mấy "chiêu kịch bản" viết tay:

### 6.1 Gồng bom tự sát — `setBom()` (`Boss.java:806-857`)

Khi bị dồn (một số boss), boss đặt `hp=1`, phát **`-45` sub `7`** (khởi bom):
```
msg = new Message(-45);
msg.writer().writeByte(7);
msg.writer().writeInt((int) id);
msg.writer().writeShort(104);    // skillId hoat anh bom
msg.writer().writeShort(2000);   // tre 2000 ms
```
Sau 2500 ms: boss chết + **gây sát thương = maxHp lên MỌI quái và MỌI người chơi trong
khu vực** (đánh diện rộng toàn map). Đây là "đòn cuối" khét tiếng.

### 6.2 Hấp thụ người chơi — ví dụ `XenBoHung.hapThu()` (`XenBoHung.java:80-101`)

Boss "nuốt" một người chơi ngẫu nhiên: dịch chuyển tới họ, **tự tăng dame/hp/crit của mình**
theo chỉ số nạn nhân, hồi đầy máu mình, **giết nạn nhân**, chat thông báo. Cơ chế viết tay,
không qua skill catalog.

### 6.3 Phân thân / triệu hồi (`boss_manifest/Commeson/PhanThan.java`, `NhanBan.java`, `SkillSummonedManager`)

Boss tạo bản sao/lính bằng chính hệ `bossAppearTogether` + `BossType.SKILLSUMMONED`. Hợp làm
Naruto (Kage Bunshin, triệu hồi thú).

> **Cho Naruto:** ba thứ trên nên mô hình hoá thành **"kịch bản pha" trong dữ liệu** (một pha
> có thể gắn "hành động đặc biệt: bom / hấp thụ / triệu hồi" + tham số), không hardcode class.

---

## 7. Mô hình dữ liệu boss — `BossData` (một pha)

`boss/BossData.java` — POJO Lombok. Một boss = `BossData[]` (mảng pha).

| Trường | Kiểu | Ý nghĩa |
|---|---|---|
| `name` | String | Tên hiển thị (có thể chứa `%d` để chèn số ngẫu nhiên, `Boss.java:159`) |
| `gender` | byte | Tộc/giới (ConstPlayer.XAYDA…) — với Naruto là **làng** |
| `outfit` | short[6] | **head, body, leg, flagBag, aura, effFront** (id vật phẩm hiển thị) |
| `dame` | long | Sát thương gốc (VD 3.000–10.000) |
| `hp` | long[] | Danh sách HP; lúc spawn **chọn ngẫu nhiên một giá trị** (`Boss.java:163`). VD `{200000000}` |
| `mapJoin` | int[] | Danh sách mapId boss có thể xuất hiện (chọn ngẫu nhiên) |
| `skillTemp` | int[][] | Bộ chiêu: mỗi dòng `{templateId, level}` hoặc `{templateId, level, cooldownMs}` |
| `textS` | String[] | Lời thoại lúc **xuất hiện** (mỗi câu `|prefix|text`) |
| `textM` | String[] | Lời thoại lúc **đang đánh** (chọn ngẫu nhiên) |
| `textE` | String[] | Lời thoại lúc **chết** |
| `secondsRest` | int | Giây nghỉ trước khi hồi sinh (pha DEFAULT_APPEAR) |
| `typeAppear` | AppearType | Kiểu xuất hiện (§4) |
| `bossesAppearTogether` | int[] | Id các boss sinh ra cùng (cụm) |

Chú ý `initBase()` (`Boss.java:156-169`): `mpg = 31_07_2002` (một số lớn cố định) → boss
**không bao giờ hết mp** (mỗi update `mp = mpg`, `Boss.java:302`). Boss dùng chiêu vô hạn.

**Cấu hình boss nằm ở đâu:**
- **Chính: hardcode** trong `BossesData.java` (~2400 dòng) + danh sách spawn trong
  `BossManager.loadBoss()` (`BossManager.java:478-526`). **KHÔNG có bảng DB định nghĩa boss.**
- **Override runtime qua file** (panel web): `panel_boss_overrides.properties` (hp/dame/name
  theo id, `BossManager.java:194-254`) — bật bằng `panel_boss_persist_status.txt`.
- **Xuất trạng thái ra file** cho web đọc: `boss_status.json` (vào `htdocs/sv1/`), mỗi 1 s
  (`BossManager.java:304-362, 820-824`).
- **Điểm săn boss ghi DB**: bảng `phuc_loi_progress` (cột `broly`/`so4`/`xen`, +1 mỗi lần
  giết, trần 10) — `PhucLoiService.recordBossKill` (`PhucLoiService.java:142-145`). Đây là
  bảng **tiến độ người chơi**, không phải định nghĩa boss.

---

## 8. Phần thưởng khi giết boss — bốn kênh, KHÔNG qua gói `-12`

Khi boss chết (`die()` → `reward(plKill)`), mỗi boss override `reward()` riêng, nhưng khung chung:

1. **Rơi đồ ra map** qua `Service.dropItemMap(zone, ItemMap)` — item **nằm trên đất**, người
   chơi nhặt bằng gói `-20` (KH7). Không nhét vào gói chết như quái. Ví dụ:
   - `BlackGoku.reward`: 15% rơi đồ thần (`randDoTL`), 60% rơi 3× item 457 (ai cũng nhặt).
   - `XENCON1.reward`: random item 17-20/702.
   - `XenBoHung.reward`: 20% đồ thần, else 2× item 457.
2. **Nhiệm vụ**: `TaskService.checkDoneTaskKillBoss(plKill, boss)` — hoàn thành nhiệm vụ theo
   `boss.id` (`TaskService.java:469+`).
3. **Điểm săn boss (event)**: `pointBoss()` (`Boss.java:619-635`) — nếu boss có HP > 490 triệu
   thì +1 điểm BHM.
4. **Điểm phúc lợi**: `PhucLoiService.recordBossKill` → cột trong `phuc_loi_progress`.

Cộng với **announce `93`** (§5.4) — cả server biết ai giết boss.

⚠️ **"Ai được thưởng" = `plKill`** — người tung đòn kết liễu (`injured()` phát hiện `isDie()`
thì gọi `setDie(plAtt); die(plAtt)`). **Bản gốc không chia damage theo đóng góp** — kẻ chốt
hạ ăn tất. Đây là một quyết định thiết kế cần chốt lại cho Naruto (§12).

---

## 9. Sát thương & phòng thủ boss — `injured()` (khác quái)

`Boss.injured()` (`Boss.java:678-710`, mỗi boss có thể override):
```
if (!piercing && Util.isTrue(nPoint.tlNeDon, 1000)) { chat("Xí hụt"); return 0; }  // NÉ theo % (phần nghìn)
if (plAtt != null && plAtt.idNRNM != -1) return 1;                                  // mien voi 1 loai skill
damage = nPoint.subDameInjureWithDeff(damage);                                      // tru theo giap
if (!piercing && effectSkill.isShielding) { …; damage = 1; }                        // KHIEN: chi an 1
nPoint.subHP(damage);
if (isDie()) { setDie(plAtt); die(plAtt); }
```

Khác quái thường:
- Boss **né đòn theo tỉ lệ** (`tlNeDon`), chat "Xí hụt" khi né.
- Boss có **khiên** (`isShielding`) — chặn đòn về còn 1 sát thương, vỡ khiên nếu đòn > maxHp.
- Nhiều boss override để **giảm sát thương pha sau** (BlackGoku `/2`, XenBoHung `/2` rồi
  `/4` khi có khiên) — cân bằng theo pha.
- Boss dùng **`nPoint` như người chơi** (hp/hpMax/dame/def/crit), không phải `MobPoint`.

---

## 10. Đối chiếu nền tảng Go hiện có — boss cần gì trước

| Thứ boss cần | Trạng thái Go hiện tại | Việc |
|---|---|---|
| Gói `-5` xuất hiện nhân vật cho thực thể **không phải chính chủ** | ✅ `xuathien.go` đã có (28 trường) | **Điền outfit thật + typePk=PK_ALL** khi là boss |
| Gói `-6` rời khu vực | 🔶 Cần kiểm Go đã phát chưa (Java `ChangeMapService:834`) | Thêm nếu thiếu (1 int) |
| Gói `44` chat bong bóng | 🔶 Chưa thấy trong game dir | Thêm (int id + UTF) |
| Gói `93` announce toàn server | 🔶 Chưa có; cần `sendMessAllPlayer` (toàn server, không chỉ zone) | Thêm; cần cơ chế gửi mọi phiên |
| Gói `-60`/`-45` combat | ✅ KH6 đang làm | Tái dùng, truyền boss vào `plAtt` |
| "Người chơi khác" hiện trong khu vực (actor thấy nhau) | 🔶 KH5 có actor khu vực; cần xác nhận đã đồng bộ char giữa các phiên | **Điều kiện tiên quyết** — boss là một "char" trong khu vực |
| Ticker riêng cho boss | ✅ mô hình actor 100 ms có sẵn | Boss update trong actor khu vực (KHÔNG cần luồng 150 ms riêng như Java) |

Go `quai.go` hiện chỉ có `isBigBoss` (quái to, luôn `false`) — **chưa có char-AI boss**.

> **Kết luận nền tảng:** boss chỉ làm được **sau khi** server Go biết cho một thực thể
> phi-người-chơi (NPC/actor) hiện ra trong khu vực và các client khác vẽ được nó (gói `-5`
> + đồng bộ char giữa phiên). Nếu KH6 chưa có "nhiều người chơi thấy nhau trong khu vực"
> thì đó là việc phải xong trước KH6b.

---

## 11. Bẫy nguy hiểm — chỗ dễ làm sai

1. **Nhầm hai loại "boss".** `mob.BigBoss`/`isBigBoss` (quái) ≠ `boss.Boss` (Player-AI). Gói,
   định danh, vòng đời khác hẳn. Bắt đầu sai chỗ này là làm lại từ đầu.

2. **Định danh boss trùng khi spawn nhiều con cùng loại.** `BossManager.loadBoss()` gọi
   `createBoss(BROLY, 5)` tạo **5 boss cùng `id = -1`** (`BossManager.java:528-532`, 507).
   Client `findCharInMap(id)` phân biệt nhân vật **theo id** → id trùng = client nhầm con
   này với con kia (chọn sai máu/vị trí). Bản gốc "thoát" nhờ chúng ở khác khu vực.
   → **Bản Go phải cấp `id` runtime DUY NHẤT cho mỗi instance boss** (ví dụ bộ đếm âm toàn
   cục), không dùng id-loại làm id-mạng.

3. **`id` giữ nguyên khi biến hình pha.** Boss pha 0 chết → pha 1 hiện lại **cùng id** với
   outfit mới. Client vẽ theo `-5` mới, nhưng cần chắc client **chấp nhận một id đã `-6`
   quay lại bằng `-5`**. 🔶 Chưa kiểm client — **phải xác minh** trước khi làm biến hình.

4. **Máu boss sát trần int32.** HP tới 500 triệu, panel cho tới **2 tỷ** (`PANEL_MAX_BOSS_HP
   = 2_000_000_000`, `BossManager.java:119`) — cách trần int32 (2,147 tỷ) chỉ ~7%. Cộng
   dồn khi hấp thụ (`hpg += ...`) có thể **tràn**. Phải kẹp `[0, 2_147_483_647]` như spec
   KH6 §1.1, và chặn override > 2 tỷ (bản gốc đã chặn, đừng bỏ).

5. **Boss tự di chuyển → tranh chấp trạng thái.** Boss đi lại/đổi zone trong luồng riêng
   (Java 150 ms). Trong Go **phải nằm trong actor khu vực** (luật §6.9/§6.10b), không luồng
   riêng — nếu không sẽ đúng bẫy KH5 (`nv.ViTriX/Y`).

6. **Boss đánh diện rộng (bom, choáng).** `setBom` gây sát thương lên **mọi** người + quái
   trong zone, lặp danh sách. Đúng loại "gói/lệnh độ dài biến thiên + nhiều mục tiêu" mà
   spec KH6 §2.2 cảnh báo — nếu làm choáng diện rộng (`-45/0`) phải test byte từng dạng.

7. **Announce `93` gửi TOÀN SERVER.** Khác mọi gói combat (chỉ trong khu vực). Cần đường
   `sendMessAllPlayer`. Spam announce (nhiều boss chết cùng lúc) → phải có hàng đợi/nhịp
   như Java (1,5 s/lần) kẻo ngập chat.

8. **`reward` chốt-hạ ăn tất.** Không chia theo sát thương. Nếu Naruto muốn công bằng (chia
   thưởng theo đóng góp) thì **phải thiết kế mới**, bản gốc không có mẫu.

9. **Boss reset MP mỗi tick** (`mpg = 31_07_2002`). Đừng bê nguyên con số ma thuật; mô hình
   Go nên có cờ "boss vô hạn chakra" tường minh.

10. **Bộ lọc announce âm thầm** (`canSendNotify`): boss ở một số map (phó bản, ma bư, offline,
    mapId 140/111) **không announce**. Nếu bê thẳng danh sách mapId Dragon Ball vào Naruto
    sẽ lọc nhầm — phải map lại theo map Naruto.

11. **Panel đọc/ghi boss qua FILE, không qua socket** — nhất quán với ba kênh web↔game của
    dự án (DB chung, `panel_cmd.txt`, JSON). Boss Go cần giữ tương thích panel: xuất
    `boss_status.json`, nhận lệnh respawn/edit. (Có thể hoãn, nhưng ghi để không quên.)

---

## 12. Quyết định thiết kế lớn cần chủ dự án chốt

1. **Boss theo LỊCH hay theo TRIGGER?** Bản gốc: **theo lịch** (spawn lúc server bật, tự
   hồi sinh sau `secondsRest`). Không có "người chơi kích hoạt boss". Naruto muốn giữ lịch,
   hay thêm boss theo sự kiện/nhiệm vụ (trigger)? (Ảnh hưởng máy trạng thái REST↔RESPAWN.)

2. **"Pha" là biến-hình-sau-khi-chết (bản gốc) hay đổi-hành-vi-theo-%máu (game hiện đại)?**
   Bản gốc chỉ có kiểu 1. Nếu Naruto muốn "boss xuống 50% máu đổi chiêu/nổi giận" thì đó là
   **cơ chế MỚI**, phải thiết kế thêm (theo dõi mốc HP trong actor). Nên chốt sớm vì nó nằm
   trong vòng lặp cập nhật boss.

3. **Chia thưởng thế nào?** Chốt-hạ ăn tất (bản gốc) / chia theo sát thương đóng góp / chia
   đều cho mọi người tham gia? Ảnh hưởng cấu trúc lưu "ai đánh boss".

4. **Announce toàn server: có hay không, tới mức nào?** Bản gốc announce cả lúc hiện lẫn lúc
   chết (có lọc theo map). Giữ nguyên / chỉ boss lớn / tắt hẳn? (Ảnh hưởng cần kênh
   `sendMessAllPlayer`.)

5. **Một class boss lái-bằng-dữ-liệu, hay class con mỗi boss?** Bản gốc: ~95 class con.
   Khuyến nghị mạnh: **một mô hình dữ liệu** (như `BossData`) + vài "hành động đặc biệt"
   khai báo trong dữ liệu (bom/hấp thụ/triệu hồi). Cần chốt vì nó định hình toàn bộ KH6b.

6. **Boss dùng chung bộ chiêu Naruto (KH6) hay có chiêu riêng?** Bản gốc boss dùng chính
   catalog chiêu người chơi (`skillTemp` trỏ vào `Skill.*`). Nếu boss Naruto cần chiêu độc
   quyền (không cho người chơi học) thì catalog chiêu phải phân biệt "chiêu chỉ-boss".

---

## 13. Phụ lục — tham chiếu mã nguồn đã kiểm

**Server Java** (`1-MAY-CHU-GAME/src/`):

| Sự kiện | Vị trí |
|---|---|
| `Boss extends Player` | `boss/Boss.java:32` |
| Máy trạng thái `update()` | `boss/Boss.java:294-350` |
| rest / respawn / joinMap / leaveMap | `boss/Boss.java:352-362 / 369-377 / 379-466 / 665-676` |
| attack (dùng skill thật) | `boss/Boss.java:551-583` |
| die + announce | `boss/Boss.java:604-610` |
| injured (né/khiên) | `boss/Boss.java:678-710` |
| gồng bom `-45/7` | `boss/Boss.java:806-857` |
| chat prefix định tuyến | `boss/Boss.java:727-764` |
| notify xuất hiện + lọc | `boss/Boss.java:477-492` |
| Mô hình 1 pha | `boss/BossData.java` |
| "DB" boss hardcode | `boss/BossesData.java` (VD Xên 3 pha `:933-1040`) |
| Manager: ticker 150 ms, factory, panel, JSON | `boss/BossManager.java:478-526, 528-706, 804-824, 304-362` |
| Override HP/dame qua file | `boss/BossManager.java:119-254` |
| Enum trạng thái / loại / xuất hiện | `boss/BossStatus.java`, `BossType.java`, `AppearType.java` |
| Boss con `parentBoss`, cụm | `boss_manifest/Cell/XENCON1.java`, `XenBoHung.java` |
| Hấp thụ | `boss_manifest/Cell/XenBoHung.java:80-101` |
| Announce `93` (chat vip) | `server/ServerNotify.java:71-84` |
| Chat `44` | `services/Service.java:596-607` |
| Rời map `-6` | `services/func/ChangeMapService.java:834` |
| Người đánh boss (`isBoss`) | `services/SkillService.java:1113-1168, 1120, 1396` |
| Menu boss `-96`/`108`/`-118` | `boss/BossManager.java:708-774`, `server/Controller.java:754-810` |
| Reward/task/điểm | `TaskService.java:469`, `func/PhucLoiService.java:142`, `boss/Boss.java:619-635` |

**Client Unity** (`D:\NRO_Project\Assets\Scripts\Game1\`):

| Sự kiện | Vị trí |
|---|---|
| `93` = chat vip (đọc UTF, dòng chữ chạy) | `Assets.src.f/Controller2.cs:825-829` |
| `44` = chat bong bóng (int id + UTF → addInfo) | `Controller.cs:4567-4572` |
| `-6` khung chung xoá char khỏi map | (đối chiếu qua `findCharInMap`, khớp `-8`/`-5`) |
| `-96` = bảng danh sách (dùng khung vTop) | `Controller.cs:233-264` |
| `mob.isBoss` (cờ **quái to**, KHÁC boss-AI) | `Controller.cs:5071, 5101` |

**Server Go** (`server/internal/game/`):

| Sự kiện | Vị trí |
|---|---|
| Gói `-5` xuất hiện nhân vật (boss tái dùng) | `xuathien.go:15-48` |
| Quái `isBigBoss` luôn false (chưa có boss-AI) | `quai.go:125, 153` |
