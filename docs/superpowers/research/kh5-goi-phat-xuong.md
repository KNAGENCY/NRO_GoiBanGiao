# KH5 — Giải mã gói PHÁT XUỐNG cho di chuyển / xuất hiện / rời khu vực

Nghiên cứu đọc mã nguồn server Java gốc tại `D:\NRO_GoiBanGiao\1-MAY-CHU-GAME\src\`.
Chỉ đọc, không sửa. Mọi dòng trích dẫn là số dòng thật tại thời điểm nghiên cứu (2026-07-24).

Quy ước Java: `DataOutputStream` — `writeByte`=1B, `writeShort`=2B big-endian,
`writeInt`=4B big-endian, `writeLong`=8B big-endian, `writeUTF`=2B độ dài (BE) + chuỗi UTF-8
kiểu "modified" (null-terminator escape kiểu Java), `writeBoolean`=1B (0/1).

Khung gói: `new Message(int command)` ép `command` thành `byte` (constructor
`network/Message.java:27-29`). Vậy `-7`, `-5`, `-6` là các giá trị byte ký hiệu (signed byte,
ví dụ -7 = 0xF9 khi đọc unsigned). Nội dung payload nằm trong `DataOutputStream` riêng
(`msg.writer()`), được gắn với `command` khi gói được đóng gói gửi đi qua
`session.sendMessage(msg)` (framing length+command nằm ngoài phạm vi câu hỏi này).

Toàn bộ `player.id` trong 3 gói dưới đây có kiểu Java là `long`
(`player/Player.java:88`) nhưng **luôn bị ép `(int)` trước khi `writeInt`** — tức trên dây
chỉ có 4 byte (int32), không phải 8 byte. Quan trọng khi viết lại bằng Go: dùng `int32` cho
playerId trên wire, dù model trong game có thể dùng kiểu rộng hơn.

Toàn bộ `x`, `y` trong Player.Location có kiểu Java `int` (`player/Location.java:11-12`)
nhưng **luôn được `writeShort`** — tức trên dây chỉ có 2 byte, giá trị bị cắt còn 16 bit thấp.

---

## 1. Gói PHÁT DI CHUYỂN — lệnh `-7`

**Luồng gọi:** `PlayerService.playerMove(player, x, y)` (`services/PlayerService.java:140-221`)
cập nhật state server-side rồi gọi `MapService.gI().sendPlayerMove(player)` ở dòng **218**:

```
services/PlayerService.java:174-176
    player.location.x = x;
    player.location.y = y;
    player.location.lastTimeplayerMove = System.currentTimeMillis();
...
services/PlayerService.java:218
    MapService.gI().sendPlayerMove(player);
```

Hàm phát gói thật nằm ở `services/MapService.java:457-469`:

```
457  public void sendPlayerMove(Player player) {
458      Message msg;
459      try {
460          msg = new Message(-7);
461          msg.writer().writeInt((int) player.id);
462          msg.writer().writeShort(player.location.x);
463          msg.writer().writeShort(player.location.y);
464          Service.gI().sendMessAnotherNotMeInMap(player, msg);
465          msg.cleanup();
466      } catch (Exception e) {
467          System.out.println("Lỗi sendPlayerMove");
468      }
469  }
```

| # | Trường | Kiểu ghi | Ý nghĩa | Bằng chứng |
|---|---|---|---|---|
| 1 | playerId | `writeInt` (4B BE) | id của người vừa di chuyển (ép từ `long` xuống `int`) | MapService.java:461 |
| 2 | x | `writeShort` (2B BE) | tọa độ x mới (server đã set ở PlayerService.java:174) | MapService.java:462 |
| 3 | y | `writeShort` (2B BE) | tọa độ y mới (server đã set ở PlayerService.java:175) | MapService.java:463 |

**Không có trường nào khác** — không tốc độ, không hướng, không animation state trong gói -7.

**Có phát lại cho chính người vừa di chuyển không?** KHÔNG. Gói được gửi qua
`Service.gI().sendMessAnotherNotMeInMap(player, msg)` (`services/Service.java:322-356`),
hàm này lặp qua `player.zone.getPlayers()` và **loại trừ chính người gửi**:

```
services/Service.java:348-351 (nhánh map online, trường hợp chính)
    for (int i = 0; i < players.size(); i++) {
        Player pl = players.get(i);
        if (pl != null && pl.getSession() != null && !pl.equals(player) && pl.isPl()) {
            pl.sendMessage(msg);
        }
    }
```

→ Client tự chịu trách nhiệm di chuyển avatar của chính mình (client-side prediction khi
người chơi bấm/chạm để di chuyển); server chỉ đẩy `-7` cho NHỮNG NGƯỜI KHÁC trong cùng zone
để họ vẽ lại vị trí của người vừa di chuyển.

---

## 2. Gói NGƯỜI CHƠI XUẤT HIỆN trong zone — lệnh `-5`

**Hàm ghi byte thật:** `map/Zone.java:653-736`, hàm private `infoPlayer(Player plReceive, Player plInfo)`
— `plReceive` là người NHẬN gói, `plInfo` là người ĐƯỢC mô tả (người vừa xuất hiện).

**Được gọi khi nào (trigger xuất hiện):**
- `Zone.load_Me_To_Another(Player player)` (Zone.java:577-604) — báo cho TẤT CẢ người khác
  trong zone rằng `player` vừa xuất hiện (gọi `infoPlayer(pl, player)` cho từng `pl` khác).
- `Zone.load_Another_To_Me(Player player)` (Zone.java:606-628) — báo cho `player` (người mới
  vào) về TẤT CẢ người/boss/pet khác đã có sẵn trong zone (gọi `infoPlayer(player, pl)` cho
  từng `pl` khác, lặp trên `this.humanoids` = player+boss+pet).
- Cả hai được gọi cùng lúc trong `ChangeMapService.finishLoadMap(Player player)`
  (`services/func/ChangeMapService.java:547-556`), là hàm xử lý khi **client báo "đã load xong
  map"** — lệnh client gửi lên là **`-39`**, xử lý tại `server/Controller.java:660-664`:

  ```
  660  case -39:
  661      if (player != null) {
  662          ChangeMapService.gI().finishLoadMap(player);
  663      }
  664      break;
  ```

  Đây là trigger chính cho "ai đó vào map/zone" (cả login lần đầu lẫn đổi map): client vào
  map → nhận `-24` (mapInfo, tĩnh: mob/npc/item/waypoint) → client load xong → gửi `-39` →
  server phát `-5` hai chiều.

**Nội dung gói `-5` (đúng thứ tự ghi, Zone.java:656-706):**

| # | Trường | Kiểu ghi | Ý nghĩa | Dòng |
|---|---|---|---|---|
| 1 | playerId | `writeInt` (ép từ `long`) | id người xuất hiện | 657 |
| 2 | clanId | `writeInt` | id bang hội; `-100` nếu là boss Mabu/SuperBu, `-2` nếu `isCopy`, `-1` nếu không có clan | 658-666 |
| 3 | level (caption) | `writeByte` | cấp bậc hiển thị, từ `CaptionManager.getInstance().getLevel(plInfo)` | 667 |
| 4 | (cờ cố định) | `writeBoolean` | luôn `false` — hard-code, không rõ ý nghĩa còn dùng | 668 |
| 5 | typePk | `writeByte` | loại PK hiện tại của người chơi | 669 |
| 6 | gender | `writeByte` | giới tính | 670 |
| 7 | gender (lặp lại) | `writeByte` | ghi lại đúng giá trị gender lần 2 (không phải lỗi đọc nhầm — code ghi 2 lần liên tiếp) | 671 |
| 8 | head | `writeShort` | id đầu tóc/mặt (`plInfo.getHead()`) | 672 |
| 9 | name | `writeUTF` (2B len + UTF-8) | tên nhân vật, qua `Service.gI().name(plInfo)` | 673 |
| 10 | hp | `writeInt` **hiện tại** (qua `writeLongByGiaCatLuong(hp, cn.readInt)`, `cn.readInt=true` → nhánh `writeInt`, 4B) | HP hiện tại (đã `Util.maxIntValue` để không tràn int) | 674, `consts/cn.java:17` |
| 11 | hpMax | như trên — `writeInt` 4B | HP tối đa | 675 |
| 12 | body | `writeShort` | id trang phục thân (`getBody()`) | 676 |
| 13 | leg | `writeShort` | id trang phục chân (`getLeg()`) | 677 |
| 14 | flagbag | `writeByte` | cờ túi đồ; có remap `83→205` nếu người NHẬN có `session.version >= 228` | 678-686 |
| 15 | (hằng số) | `writeByte` | luôn `-1`, không tên biến/comment | 687 |
| 16 | x | `writeShort` | tọa độ x hiện tại của `plInfo` | 688 |
| 17 | y | `writeShort` | tọa độ y hiện tại của `plInfo` | 689 |
| 18 | effBuffHp | `writeShort` | luôn `0` (hard-code, có comment) | 690 |
| 19 | effBuffMp | `writeShort` | luôn `0` (hard-code, có comment) | 691 |
| 20 | numEff | `writeByte` | luôn `0` — số hiệu ứng buff đi kèm; nếu >0 thì mỗi effect còn có `byte templateId, int timeStart, int timeLength, short param` theo comment (692-695) nhưng KHÔNG được ghi trong code hiện tại vì numEff luôn 0 | 693-695 |
| 21 | idSpaceShip | `writeByte` | id tàu vũ trụ (`iDMark.getIdSpaceShip()`) | 696 |
| 22 | isMonkey | `writeByte` | 1 nếu `effectSkill.isMonkey`, else 0 | 698 |
| 23 | mount | `writeShort` | id thú cưỡi/mount (`getMount()`) | 699 |
| 24 | cFlag | `writeByte` | cờ (quốc kỳ/team flag) hiện tại | 700 |
| 25 | (hằng số) | `writeByte` | luôn `0`, không tên biến | 702 |
| 26 | aura | `writeShort` | id hiệu ứng aura (`getAura()`) | 703 |
| 27 | effFront | `writeByte` | id set hiệu ứng phía trước (`getEffFront()`) | 704 |
| 28 | hat | `writeShort` | id mũ (`getHat()`) | 705 |

**Các gói phụ đi kèm ngay sau `-5` cho cùng người (vẫn trong `infoPlayer`, không phải phần
của gói -5, nhưng luôn đi kèm khi "người chơi xuất hiện"):**

- `Service.gI().sendFlagPlayerToMe(plReceive, plInfo)` — lệnh **`-103`**
  (`services/Service.java:1490-1501`): `writeByte(2)`, `writeByte(pl.cFlag)`,
  `writeShort(flagIconId[pl.cFlag])`.
- Nếu `plInfo.effectSkill.isChibi` → `sendChibiFollowToMe`, ngược lại → `sendPetFollowToMe`
  (Zone.java:712-720) — gói theo/pet đi kèm, không thuộc phạm vi câu hỏi này (chỉ chạy khi
  `plInfo.isPl()`).
- Nếu `plInfo.isDie()` → thêm gói lệnh **`-8`** (Zone.java:723-735):
  `writeInt((int) plInfo.id)`, `writeByte(0)`, `writeShort(x)`, `writeShort(y)` — báo người
  này đang trong trạng thái chết ngay khi xuất hiện.

---

## 3. Gói NGƯỜI CHƠI RỜI/BIẾN MẤT — lệnh `-6`

**Nơi ghi byte thật:** `services/func/ChangeMapService.java:830-847`, trong hàm `exitMap(Player player)`
(gọi khi người chơi rời zone hiện tại — đổi map, thoát game, v.v.):

```
830  player.zone.removePlayer(player);
831  if (!MapService.gI().isMapOffline(player.zone.map.mapId)) {
832      Message msg = null;
833      try {
834          msg = new Message(-6);
835          msg.writer().writeInt((int) player.id);
836          Service.gI().sendMessAnotherNotMeInMap(player, msg);
837          player.zone = null;
838      } catch (Exception e) {
...
```

| # | Trường | Kiểu ghi | Ý nghĩa | Dòng |
|---|---|---|---|---|
| 1 | playerId | `writeInt` (ép từ `long`) | id người vừa rời zone | 835 |

**Chỉ 1 trường duy nhất — id.** Không có tọa độ, không có lý do rời.
Gửi qua `sendMessAnotherNotMeInMap` — cũng loại trừ chính người rời đi (không cần thiết vì
người đó không còn xem gói của chính mình), chỉ những người CÒN LẠI trong zone nhận được.

Lưu ý: `player.zone.removePlayer(player)` (dòng 830, cập nhật danh sách zone nội bộ —
`map/Zone.java:283-292`) được gọi **TRƯỚC KHI** gửi `-6`, và chỉ gửi `-6` nếu map không phải
"offline map" (`MapService.isMapOffline`). Với map offline (map riêng tư kiểu training), không
ai khác cần biết vì không có người chơi thật khác trong đó.

Đối lệnh `-6` với gói `5`/`6` (không dấu, POSITIVE) tìm thấy tại
`services/Service.java:1390` và `boss/boss_manifest/Training/TrainingBoss.java:226` —
đây là các lệnh khác dùng cho luồng riêng của Boss training offline (map riêng tư từng
người), KHÔNG phải gói broadcast xuất hiện/rời zone multiplayer chuẩn. Không đi sâu vì
không nằm trong phạm vi câu hỏi (không liên quan luồng nhiều người chơi/zone chung).

---

## 4. Gói tối thiểu cho KH5 — chỉ MỘT người chơi (chưa có ai khác/quái)

**Kết luận: server không cần phát gói mạng nào (`-5`/`-6`/`-7`) cho chính người chơi khi họ tự
di chuyển, kể cả trong bản gốc nhiều người chơi.** Đây không phải là một tối ưu hoá riêng cho
KH5 — nó là hành vi CÓ SẴN trong code gốc, vì lý do:

1. **`-7` (move) không bao giờ echo về chính người gửi** — `sendMessAnotherNotMeInMap` luôn
   loại trừ `player` khỏi danh sách nhận (Service.java:350: `!pl.equals(player)`). Client gốc
   tự vẽ chuyển động của chính nó (client-side prediction) ngay khi người dùng bấm di chuyển,
   không chờ server xác nhận. Khi zone chỉ có 1 người (chính họ), vòng lặp gửi
   `sendMessAnotherNotMeInMap`/`sendMessAllPlayerInMap` chạy nhưng không tìm thấy người nhận
   nào khác → **0 gói thực sự được gửi ra mạng**, nhưng đó vẫn đúng hành vi gốc dù server có
   người khác hay không.

2. **`-5` (xuất hiện) qua `load_Me_To_Another`/`load_Another_To_Me`** cũng lặp trên danh sách
   zone (`players`/`humanoids`) loại trừ chính mình. Nếu solo, danh sách "người khác" rỗng →
   0 gói `-5` gửi đi. Trigger của `-5` là sự kiện **client gửi `-39` (báo đã load xong map)**,
   không phải sự kiện di chuyển — nên nó chỉ cần chạy đúng 1 lần khi vào map, không lặp lại
   mỗi lần di chuyển.

3. **`-6` (rời zone)** tương tự — chỉ có ý nghĩa khi có người khác trong zone để thông báo.

**Vậy tối thiểu server (Go, KH5) cần làm khi CHỈ MỘT người chơi tự di chuyển:**

- Vẫn phải xử lý toàn bộ logic `playerMove` phía server (tương đương
  `PlayerService.playerMove`, PlayerService.java:140-221): validate không chết
  (`!player.isDie()`), cập nhật `player.location.x/y` authoritative, cập nhật
  `lastTimeplayerMove`, các side-effect map-cụ-thể (giới hạn map, bay/fly, tiêu hao mp khi
  bay, task check) — để giữ state server đúng, chống gian lận, và để các hệ thống khác (task,
  spawn quái theo vị trí…) hoạt động đúng.
- **Gọi hàm phát `-7` như bình thường** (`sendPlayerMove` tương đương) — không cần bỏ qua nó,
  chỉ cần đảm bảo vòng lặp gửi-cho-người-khác không lỗi khi danh sách người nhận rỗng (đúng
  hành vi code gốc: `sendMessAllPlayerInMap`/`sendMessAnotherNotMeInMap` đã tự xử lý
  `players.isEmpty()` bằng cách `dispose()` message, xem Service.java:271-274, 344-347).
- **Không cần** implement gì thêm để "client mượt": client tự vẽ chuyển động cục bộ ngay khi
  người dùng thao tác, hoàn toàn không phụ thuộc phản hồi `-7` cho chính họ — đây là thiết kế
  gốc, không phải suy luận riêng cho KH5.
- Khi người chơi đó vào map (đăng nhập hoặc đổi map) vẫn nên phát đúng luồng `-24`
  (mapInfo, tĩnh) → client gửi `-39` → server chạy `load_Me_To_Another`/`load_Another_To_Me`
  (dù rỗng), để giữ đúng state machine cho các bước sau của KH (KH6+ khi có người chơi thứ 2
  hoặc quái sẽ cắm ngay vào đúng luồng có sẵn mà không cần đổi lại giao thức).

**Tóm lại: xử lý `-7` phía server luôn cần chạy (để giữ state đúng), nhưng vì danh sách
người nhận rỗng khi solo, không có byte nào thực sự đi ra dây — không cần code đường tắt
riêng, chỉ cần không crash khi zone rỗng người khác.**

---

## Phụ lục — vị trí file & dòng tham chiếu chính

| File | Dòng | Nội dung |
|---|---|---|
| `services/PlayerService.java` | 140-221 | `playerMove` — logic + trigger `sendPlayerMove` |
| `services/MapService.java` | 457-469 | `sendPlayerMove` — đóng gói `-7` |
| `services/Service.java` | 265-282 | `sendMessAllPlayerInMap(Zone,...)` |
| `services/Service.java` | 284-320 | `sendMessAllPlayerInMap(Player,...)` |
| `services/Service.java` | 322-356 | `sendMessAnotherNotMeInMap` — loại trừ chính mình |
| `services/Service.java` | 1490-1501 | `sendFlagPlayerToMe` — gói phụ `-103` |
| `map/Zone.java` | 256-292 | `addPlayer`/`removePlayer` — quản lý list nội bộ zone |
| `map/Zone.java` | 577-604 | `load_Me_To_Another` |
| `map/Zone.java` | 606-628 | `load_Another_To_Me` |
| `map/Zone.java` | 653-736 | `infoPlayer` — đóng gói `-5` (+ gói phụ `-103`, `-8`) |
| `services/func/ChangeMapService.java` | 547-566 | `finishLoadMap` — trigger `-5` hai chiều |
| `services/func/ChangeMapService.java` | 805-847 | `exitMap` — đóng gói `-6` |
| `server/Controller.java` | 660-664 | client `-39` → `finishLoadMap` |
| `network/Message.java` | 27-35 | `Message(int command)` ép byte |
| `network/Message.java` | 221-227 | `writeLongByGiaCatLuong` |
| `consts/cn.java` | 17 | `readInt = true` (quyết định hp/hpMax ghi 4B) |
| `player/Player.java` | 88 | `id` là `long` nhưng luôn ép `int` khi ghi wire |
| `player/Location.java` | 11-12 | `x`,`y` là `int` nhưng luôn ghi `writeShort` |
