# KH5 — Giải mã gói DI CHUYỂN & HÀNH ĐỘNG (client → server)

Nghiên cứu chỉ-đọc từ mã Java gốc tại `D:\NRO_GoiBanGiao\1-MAY-CHU-GAME\src\`.
Mục tiêu: làm cơ sở cho Kế hoạch 5 (di chuyển + actor) của server Go viết lại.

## 0. Khung message chung

- `Controller.onMessage(ISession s, Message _msg)` — `D:\NRO_GoiBanGiao\1-MAY-CHU-GAME\src\server\Controller.java:86`
- Lệnh (command) được đọc sẵn thành 1 byte trước switch: `byte cmd = _msg.command;` (dòng 114), sau đó `switch (cmd)`.
- `_msg.reader()` trả về `java.io.DataInputStream` chuẩn (xác nhận tại `D:\NRO_GoiBanGiao\1-MAY-CHU-GAME\src\network\Message.java:49`), nên các kiểu đọc là **big-endian, chuẩn Java**:
  - `readByte()` — 1 byte, signed (-128..127)
  - `readUnsignedByte()` — 1 byte, unsigned (0..255)
  - `readShort()` — 2 byte, BE, signed
  - `readInt()` — 4 byte, BE, signed
  - `readLong()` — 8 byte, BE, signed
  - `readUTF()` — 2 byte length (unsigned short) + modified UTF-8 bytes
- Có rate-limit chung ở tầng session (không riêng cho di chuyển): nếu > 1500 packet/5s → kick (dòng 100-108); nếu unknown cmd → `AntiSpam.handleUnknownCmd` có thể dẫn tới ban IP (dòng 837-849). Đây là chống spam nói chung, không phải chống teleport tọa độ.

---

## 1. Gói `-7` — DI CHUYỂN (di chuyển client → server)

Vị trí: `Controller.java:413-441`.

```java
case -7:
    if (player != null) {
        if (player.isDie()) { Service.gI().charDie(player); return; }
        if (player.effectSkill.isHaveEffectSkill()) { return; }
        int toX = player.location.x;
        int toY = player.location.y;
        try {
            byte b = _msg.reader().readByte();      // cờ loại di chuyển
            toX = _msg.reader().readShort();          // toX
            try {
                toY = _msg.reader().readShort();       // toY (optional!)
            } catch (IOException ex) { }
            if (player.zone != null && MapService.gI().isMapBlackBallWar(player.zone.map.mapId)
                    && Util.getDistance(player.location.x, player.location.y, toX, toY) > 500) {
                return;                                 // chặn nếu nhảy quá xa — CHỈ áp dụng map 85-91
            }
            if (b == 1) {
                AchievementService.gI().checkDoneTaskFly(player, player.location.x - toX);
            }
        } catch (IOException e) { }
        PlayerService.gI().playerMove(player, toX, toY);
    }
    break;
```

### Bảng byte, đọc theo đúng thứ tự

| # | Trường | Kiểu đọc | Kích thước | Ý nghĩa | Dòng bằng chứng |
|---|--------|----------|-------------|---------|------------------|
| 1 | `b` | `readByte()` | 1B signed | Cờ loại di chuyển. `b == 1` → di chuyển bằng "khinh công / bay" (dùng để cộng dồn thành tích `KHINH_CONG_THANH_THAO`, xem `AchievementService.checkDoneTaskFly` ở `D:\NRO_GoiBanGiao\1-MAY-CHU-GAME\src\models\Achievement\AchievementService.java:141-148` — `length = abs((x_cũ - toX)/10)`, cộng dồn nếu length < 10). Các giá trị khác (0, …) không được xử lý riêng — coi là đi bộ thường. | Controller.java:425, 435-437 |
| 2 | `toX` | `readShort()` | 2B BE signed | Tọa độ X đích | Controller.java:426 |
| 3 | `toY` | `readShort()` | 2B BE signed, **optional** (bọc trong `try/catch` riêng — nếu client không gửi đủ byte, `toY` giữ nguyên giá trị mặc định = `player.location.y` hiện tại, KHÔNG lỗi toàn gói) | Tọa độ Y đích | Controller.java:427-430 |

Ghi chú quan trọng cho Go rewrite:
- Toàn bộ khối đọc byte nằm trong 1 `try/catch(IOException)` bên ngoài — nếu đọc `b` hoặc `toX` lỗi (thiếu byte), server **không throw**, chỉ bỏ qua và gọi `playerMove` với tọa độ hiện tại (coi như không di chuyển gì / no-op). Đây là hành vi "graceful fallback", không phải reject cứng.
- `toY` có try/catch RIÊNG lồng bên trong — nghĩa là gói `-7` về mặt lý thuyết có thể chỉ gửi `b + toX` (thiếu toY) và vẫn được chấp nhận, chỉ Y không đổi. Trên thực tế client gốc luôn gửi đủ 5 byte (1+2+2), nhưng server code chấp nhận thiếu.
- Có 2 điều kiện chặn sớm trước khi đọc byte nào: `player.isDie()` (chết thì không xử lý di chuyển, gọi lại `charDie`) và `player.effectSkill.isHaveEffectSkill()` (đang bị hiệu ứng khống chế/skill effect thì bỏ qua toàn bộ gói, kể cả không đọc byte — **cẩn thận**: nếu điều kiện này đúng, hàm return ngay TRƯỚC khi đọc payload, nghĩa là các byte của gói `-7` bị bỏ lại chưa đọc trong buffer của message đó. Vì Java xử lý mỗi `Message` là 1 buffer độc lập (`_msg.cleanup()`/`dispose()` ở `finally` của `onMessage`, dòng 861-864), việc không đọc hết không làm lệch state của gói tiếp theo — an toàn để port sang Go miễn là mỗi gói cũng được framing/đọc độc lập theo kiểu tương tự (length-prefixed).

---

## 2. Các gói hành động liên quan (-71, -72, -74, -76, -78, -79)

| Lệnh | Dòng | Tên/ý nghĩa suy được | Trường đọc | Có phải di chuyển tối thiểu? |
|------|------|------------------------|------------|-------------------------------|
| `-72` | Controller.java:282-286 | Chat riêng (private chat) — gọi `FriendAndEnemyService.gI().chatPrivate(player, _msg)` (đọc byte tiếp bên trong `chatPrivate`, không đọc gì ở Controller) | (đọc trong `chatPrivate`, không thuộc phạm vi di chuyển) | KHÔNG — phụ (chat) |
| `-71` | Controller.java:373-381 | Chat toàn cục (global chat) | `_msg.reader().readUTF()` — 1 chuỗi UTF | KHÔNG — phụ (chat) |
| `-79` | Controller.java:382-386 | Mở menu thông tin người chơi khác (context menu khi click vào 1 player) — `Service.gI().getPlayerMenu(player, id)` | `_msg.reader().readInt()` — 4B, ID của player mục tiêu | KHÔNG — phụ (UI/tương tác, không phải nhặt đồ/tấn công nhưng cũng không phải di chuyển) |
| `-74` | Controller.java:443-498 | Xác thực client + yêu cầu tải dữ liệu ban đầu (anti-bot token + data request), KHÔNG liên quan di chuyển | 1) `byte[40] token` qua `readFully(token)` — xác minh bằng `ClientVerifier.verify(token)`, sai thì kick; 2) sau đó (cùng case, có rate-limit tải dữ liệu riêng: >5 lần/10s thì kick) đọc tiếp `byte type` — `type==1` gọi `DataGame.sendSizeRes`, `type==2` gọi `DataGame.sendRes` | KHÔNG — phụ (bootstrap/anti-cheat, không phải hành động gameplay) |
| `-78` | Controller.java:814-816 | Comment ghi `// checkMMove` (có vẻ từng dùng để client báo "còn đang di chuyển"/heartbeat chống speed-hack) nhưng implementation HIỆN TẠI chỉ đọc và **bỏ (discard)** | `_msg.reader().readInt()` — 4B, biến gọi là "second" (có thể là timestamp/giây), không dùng vào đâu cả — no-op thực sự | KHÔNG cần cho di chuyển tối thiểu (code hiện hành không làm gì với giá trị này) |
| `-76` | Controller.java:825-836 | Nhận thưởng thành tích (claim achievement/gem) — tùy `player.typeRecvieArchiment` (0/1/2) gọi 1 trong 3 hàm nhận thưởng khác nhau | `_msg.reader().readByte()` — 1B, index/id thành tích cần nhận (chỉ đọc nếu `_msg.reader().available() >= 1`) | KHÔNG — phụ (nhặt/nhận thưởng, không phải di chuyển) |

**Kết luận mục 2**: Trong nhóm 6 gói được liệt kê, KHÔNG có gói nào là phần bắt buộc của "đi bộ tối thiểu" — tất cả đều là phụ trợ (chat, UI menu, auth/bootstrap, heartbeat bỏ không, nhận thưởng). Gói di chuyển thật sự chỉ có `-7`.

---

## 3. Gói TỐI THIỂU cho đi bộ trong map (walk-only, KH5 scope hẹp nhất)

Chỉ cần đúng **1 gói client→server**:

- **`-7`**: `byte b` (cờ bay/khinh công, dùng cho achievement — có thể bỏ qua ở bản KH5 tối giản nếu chưa làm hệ thống thành tích) + `short toX` + `short toY`.

Và ở phía server, hàm lõi cần port là `PlayerService.playerMove(Player player, int x, int y)` (`D:\NRO_GoiBanGiao\1-MAY-CHU-GAME\src\services\PlayerService.java:140-221`), việc nó làm:

1. Guard: `player.zone == null` → return (chưa vào map nào thì bỏ qua).
2. Một số rule đặc thù theo `mapId` (map 176 chặn clan, map 127 chỉ mở cửa sổ giờ 22h-23h) — đá người chơi ra map khác nếu vi phạm (KHÔNG bắt buộc cho KH5 tối giản, đây là rule nghiệp vụ riêng của map).
3. Guard: chỉ xử lý nếu `!player.isDie()`.
4. Dừng các hiệu ứng đang chạy dở nếu có (`isCharging`, `useTroi`) — phụ.
5. **Cập nhật vị trí thật**: `player.location.x = x; player.location.y = y; player.location.lastTimeplayerMove = System.currentTimeMillis();`
6. Rule biên bản đồ riêng cho map 85-91 (BlackBallWar) — nếu ra ngoài biên thì đá về map khác (chi tiết ở mục 4 bên dưới).
7. Cho pet đi theo (`pet.followMaster()`) — phụ.
8. Kiểm tra ô tile hiện tại có phải ô "bay được" không (`tileMap[y/24][x/24] == 0` → `isFly = true`) và trừ mana nếu đang bay không cưỡi thú — phụ, không cần cho KH5 tối giản.
9. **Phát gói di chuyển cho người chơi khác trong cùng zone**: `MapService.gI().sendPlayerMove(player)` — **đây là phần bắt buộc để client khác thấy được di chuyển**.
10. `TaskService.gI().checkDoneTaskGoToMap(...)` — phụ (nhiệm vụ).

→ **Gói tối thiểu cho đi bộ = client gửi `-7` (b, toX, toY) → server cập nhật `player.location` → server phát lại gói `-7` (id, x, y) cho các player khác cùng zone.**

### Gói PHÁT (server → client khác) sau khi di chuyển

`MapService.sendPlayerMove(Player player)` — `D:\NRO_GoiBanGiao\1-MAY-CHU-GAME\src\services\MapService.java:457-469`:

```java
msg = new Message(-7);
msg.writer().writeInt((int) player.id);
msg.writer().writeShort(player.location.x);
msg.writer().writeShort(player.location.y);
Service.gI().sendMessAnotherNotMeInMap(player, msg);
```

- **Cùng số lệnh `-7`** được tái sử dụng cho cả 2 chiều (lên: b+toX+toY; xuống: playerId+x+y) — client và server dùng chung command id nhưng format payload khác nhau tùy hướng.
- Phạm vi phát: hàm `Service.sendMessAnotherNotMeInMap` (`D:\NRO_GoiBanGiao\1-MAY-CHU-GAME\src\services\Service.java:322-356`) lặp qua `player.zone.getPlayers()` và gửi cho tất cả player khác (loại trừ chính người vừa di chuyển) đang **trong cùng Zone hiện tại** của người chơi — KHÔNG phải toàn bộ Map (Map có thể có nhiều Zone). Có nhánh riêng cho map offline (chỉ gửi cho các đối tượng liên quan tới `TrainingBoss`/`NonInteractiveNPC`).
- (Chi tiết định dạng đầy đủ + các trường mở rộng của gói phát này để agent khác giải mã kỹ hơn nếu cần — ở đây chỉ ghi tên lệnh theo yêu cầu.)

---

## 4. Ràng buộc / kiểm tra hợp lệ đã có (chống gian lận)

### 4.1 Chống teleport / dịch chuyển tức thời

- **CHỈ tồn tại cho nhóm map "BlackBallWar" (mapId 85-91)**, tại `Controller.java:431-434`:
  ```java
  if (player.zone != null && MapService.gI().isMapBlackBallWar(player.zone.map.mapId)
          && Util.getDistance(player.location.x, player.location.y, toX, toY) > 500) {
      return; // hủy gói di chuyển này
  }
  ```
  - Ngưỡng: **500** đơn vị, đo bằng khoảng cách Euclid `Util.getDistance` (`D:\NRO_GoiBanGiao\1-MAY-CHU-GAME\src\utils\Util.java:106-108`: `sqrt((x1-x2)^2 + (y1-y2)^2)`), so giữa vị trí HIỆN TẠI của player (`player.location.x/y`, tức vị trí đã lưu ở lần move trước) và tọa độ đích `toX/toY` mới nhận được.
  - `isMapBlackBallWar(mapId)` định nghĩa tại `MapService.java:484-486`: `mapId >= 85 && mapId <= 91`.
  - **QUAN TRỌNG**: đây là kiểm tra DUY NHẤT trong toàn bộ luồng `-7` → `playerMove` giới hạn khoảng-cách-di-chuyển-mỗi-gói. Với TẤT CẢ các map khác (ngoài 85-91), server **KHÔNG có kiểm tra chống teleport theo khoảng cách** — client có thể gửi toX/toY bất kỳ và server sẽ chấp nhận (miễn còn nằm trong logic biên map ở mục 4.2, vốn cũng chỉ áp dụng cho 85-91). Đây là lỗ hổng thật sự trong code gốc — nên ghi chú cho KH5 Go rewrite cân nhắc bổ sung kiểm tra khoảng-cách-tối-đa-mỗi-tick chung cho MỌI map (không chỉ 85-91) nếu muốn chống cheat triệt để hơn bản gốc.
  - Không có kiểm tra dựa trên thời gian (không so `lastTimeplayerMove` với thời gian hiện tại để tính tốc độ tối đa cho phép) — dù trường `player.location.lastTimeplayerMove` có được cập nhật (`PlayerService.java:176`), nó KHÔNG được đọc lại ở đâu trong `playerMove` hay case `-7` để tính vận tốc. Có vẻ từng dự tính dùng cho anti-speedhack nhưng chưa hoàn thiện (khớp với gói `-78 checkMMove` cũng bị bỏ không dùng, xem mục 2).

### 4.2 Chống ra ngoài biên bản đồ

- Chỉ áp dụng cho map 85-91 (cùng nhóm BlackBallWar), tại `PlayerService.java:177-199`:
  - Nếu KHÔNG phải boss/pet, và `x < 24 || x > mapWidth-24 || y < 0 || y > mapHeight-24` → nếu không đứng ở waypoint hợp lệ thì bị đẩy về map khác (`ChangeMapService.changeMap(player, 21+gender, 0, 200, 336)`), không phải "kick" mà là "trục xuất khỏi map hiện tại".
  - Có thêm kiểm tra địa hình top-collision `yTop = map.yPhysicInTop(x,y)`.
  - **Các map khác không có bất kỳ giới hạn biên nào trong `playerMove`** — nghĩa là về lý thuyết client có thể gửi tọa độ âm hoặc vượt map size cho các map thường mà server vẫn set `player.location.x/y` theo đúng giá trị đó (không clamp). Rủi ro tương tự mục 4.1.

### 4.3 Các guard khác (không phải chống teleport nhưng liên quan hợp lệ)

- `player.isDie()` → chặn xử lý di chuyển khi đã chết (Controller.java:415-418).
- `player.effectSkill.isHaveEffectSkill()` → chặn khi đang bị khống chế bởi hiệu ứng skill (Controller.java:419-421).
- `player.zone == null` → chặn trong `playerMove` nếu chưa gắn vào Zone nào (PlayerService.java:141-143).
- Rule nghiệp vụ theo mapId cụ thể (176 chặn clan, 127 giới hạn khung giờ) — không phải chống cheat, là luật chơi.
- Rate-limit tầng session chung (packet/5s, dữ liệu tải/10s) — chống spam gói nói chung, không riêng cho `-7`.

---

## 5. Tóm tắt cho Kế hoạch 5 (Go rewrite)

- Gói lên tối thiểu: `-7` = `[byte flag][short toX][short toY]` (flag: 1 = đang bay/khinh công, dùng optional cho thành tích).
- Gói xuống tối thiểu để đồng bộ vị trí cho người khác: cũng lệnh `-7` = `[int playerId][short x][short y]`, phát cho toàn bộ player khác **trong cùng zone**.
- Server gốc gần như KHÔNG chống teleport ngoại trừ 1 rule hẹp (map 85-91, ngưỡng khoảng cách 500/gói, không có yếu tố thời gian). Không có clamp biên bản đồ chung, không có giới hạn tốc độ (speed cap) dựa trên thời gian giữa 2 lần move. → Bản Go rewrite nên chủ động làm chặt hơn bản gốc: (a) clamp toạ độ trong biên map cho MỌI map, (b) thêm kiểm tra vận tốc tối đa dựa trên `Δdistance / Δtime` từ `lastTimeplayerMove` (trường đã có sẵn ý tưởng trong code gốc nhưng chưa dùng) thay vì ngưỡng khoảng cách cố định chỉ áp dụng 1 nhóm map.
- Các gói `-71, -72, -74, -76, -78, -79` đều KHÔNG cần thiết cho "đi bộ tối thiểu" — để dành cho các kế hoạch sau (chat, UI, auth/bootstrap, nhận thưởng thành tích). `-78` hiện là no-op thực sự trong code gốc, có thể bỏ qua hoàn toàn hoặc tái sử dụng cho mục đích khác (vd. client heartbeat) trong bản Go.
