# KH6 — Giải mã giao thức CHIẾN ĐẤU (đánh quái, đánh người, chiêu thức, chết & hồi sinh)

Nghiên cứu đọc mã nguồn server Java gốc tại `D:\NRO_GoiBanGiao\1-MAY-CHU-GAME\src\`.
**Chỉ đọc, không sửa.** Mọi số dòng là số dòng thật tại thời điểm nghiên cứu (2026-07-26).

Tiếp nối hai tài liệu KH5: [`kh5-goi-dichuyen-len.md`](kh5-goi-dichuyen-len.md) (gói client gửi lên)
và [`kh5-goi-phat-xuong.md`](kh5-goi-phat-xuong.md) (gói server phát xuống). Quy ước kiểu dữ liệu
Java (`writeByte`=1B, `writeShort`=2B BE, `writeInt`=4B BE, `writeBoolean`=1B) giữ nguyên,
không nhắc lại.

---

## 0. Hai phát hiện chi phối toàn bộ tài liệu này

### 0.1 `writeLongByGiaCatLuong` = `writeInt` (4 byte), không phải 8 byte

Mọi trường **máu** và **sát thương** trong giao thức đều đi qua hàm này:

```
network/Message.java:221-227
    public void writeLongByGiaCatLuong(long v, boolean is) throws IOException {
        if (is) {
            this.writeInt((int) v);      // ← nhánh đang chạy
            return;
        }
        this.writeLong(v);
    }
```

Tham số `is` **luôn** là hằng `cn.readInt`, và `consts/cn.java:17` đặt:

```
public static boolean readInt = true;
```

→ **Kết luận: mọi trường hp / hpMax / dame trên dây là int32 big-endian (4 byte).**

Kèm theo đó, `utils/Util.java:153-161` chặn trần:

```
public static long maxIntValue(double a) {
    if (cn.readInt) {
        if (a > Integer.MAX_VALUE) { a = Integer.MAX_VALUE; }
        return (int) a;
    }
    return (long) a;
}
```

→ Máu và sát thương **bão hoà ở 2.147.483.647**, không tràn âm. Bản Go phải kẹp
(clamp) y hệt, **không được để tràn `int32`** — tràn thì client hiện sát thương âm.

Điều này khớp với code Go đã có: `server/internal/game/xuathien.go:28-29` đang ghi
`w.Int(nv.HpToiDa)` kèm chú thích "writeInt vì cn.readInt=true". **Nhất quán, không cần đổi.**

### 0.2 Client chỉ **đề nghị mục tiêu**, server quyết định tất cả

Gói client gửi lên khi đánh **không chứa sát thương, không chứa toạ độ đòn đánh** — chỉ
có *id mục tiêu*. Toàn bộ sát thương, né, chí mạng, chết, rơi đồ do server tính.

Đây là tin tốt: nó hợp sẵn với ràng buộc §5.1 của thiết kế gốc ("máy chủ là bên quyết
định thật") và với mô hình actor của KH5. **Không phải hoà giải gì thêm ở tầng giao thức.**

Chiêu thức đang chọn **không nằm trong gói đánh** — nó là *trạng thái phiên*, đặt trước
bằng gói `34` (xem §2). Hệ quả quan trọng cho mô hình actor: **chiêu đang chọn phải là dữ
liệu của actor khu vực**, không phải của luồng phiên — nếu để ở luồng phiên thì luồng
phiên ghi trong khi actor đọc, đúng vào cái bẫy đã vấp ở KH5 (§6.9 của `BAT-DAU-TU-DAY.md`).

---

## 1. Gói LÊN — client → server

### 1.1 `54` — đánh quái

`server/Controller.java:697-708`:

```
697  case 54:
698      if (player != null) {
699          int mobId = _msg.reader().readByte();
700          int masterId = -1;
701          boolean isMobMe = mobId == -1;
702          if (isMobMe) {
703              masterId = _msg.reader().readInt();
704          }
705          // _msg.reader().readByte();
706          Service.gI().attackMob(player, mobId, isMobMe, masterId);
707      }
708      break;
```

| # | Trường | Kiểu | Ý nghĩa |
|---|---|---|---|
| 1 | `mobId` | byte | Chỉ số quái **trong khu vực** (không phải tempId). `-1` = đánh "quái của người khác" (mobMe/thú cưng triệu hồi) |
| 2 | `masterId` | int | **Chỉ có khi `mobId == -1`** — id người chơi sở hữu con mobMe |

**Bẫy độ dài biến thiên:** gói này dài 1 byte *hoặc* 5 byte tuỳ giá trị byte đầu. Đọc sai
thì phần thân còn lại bị hiểu sai (nhưng **không** làm hỏng các gói sau — xem §10).

Dòng 705 bị comment: bản gốc từng đọc thêm 1 byte nữa rồi bỏ. **Bản Go không đọc byte này.**

**Phạm vi KH6:** mobMe (`mobId == -1`) thuộc hệ thú cưng — **hoãn sang KH7 trở đi**.
Bản Go giai đoạn KH6 vẫn phải *đọc đúng 5 byte* khi gặp `-1` (để không lệch dòng) rồi
lặng lẽ bỏ qua.

### 1.2 `-60` — đánh người chơi

`server/Controller.java:709-715`:

```
709  case -60:
710      if (player != null) {
711          int playerId = _msg.reader().readInt();
712          // _msg.reader().readByte();
713          Service.gI().attackPlayer(player, playerId);
714      }
715      break;
```

| # | Trường | Kiểu | Ý nghĩa |
|---|---|---|---|
| 1 | `playerId` | int | Id người chơi bị đánh |

### 1.3 `34` — chọn chiêu

`server/Controller.java:691-696`:

```
691  case 34:
692      if (player != null) {
693          int selectSkill = _msg.reader().readShort();
694          SkillService.gI().selectSkill(player, selectSkill);
695      }
696      break;
```

| # | Trường | Kiểu | Ý nghĩa |
|---|---|---|---|
| 1 | `selectSkill` | short | **`template.id`** của chiêu, KHÔNG phải `skillId` |

Bằng chứng cho "template.id chứ không phải skillId" — `services/SkillService.java:1462-1469`:

```
1462  public void selectSkill(Player player, int skillId) {
1463      for (Skill skill : player.playerSkill.skills) {
1464          if (skill.skillId != -1 && skill.template.id == skillId) {   // ← so với template.id
1465              player.playerSkill.skillSelect = skill;
```

Nhưng khi **phát xuống** thì lại ghi `skillSelect.skillId` (§3.1, §3.2). **Hai trường khác
nhau, đừng lẫn.**

### 1.4 `-16` — hồi sinh

`server/Controller.java:744-748`:

```
744  case -16: // hồi sinh
745      if (player != null && !player.isPKDHVT) {
746          PlayerService.gI().hoiSinh(player);
747      }
```

Không có thân gói. Server tự trừ phí và quyết định.

### 1.5 `-20` — nhặt đồ rơi

`server/Controller.java:723-728`:

```
723  case -20:
724      if (player != null && !player.isDie()) {
725          int itemMapId = _msg.reader().readShort();
726          ItemMapService.gI().pickItem(player, itemMapId, false);
```

| # | Trường | Kiểu | Ý nghĩa |
|---|---|---|---|
| 1 | `itemMapId` | short | Id món đồ đang nằm dưới đất trong khu vực |

Lưu ý guard `!player.isDie()` — chết thì không nhặt được.

---

## 2. Cổng kiểm tra trước khi cho đánh (server-side)

`services/SkillService.java:58-136` là hàm `useSkill`. Các chốt chặn, theo đúng thứ tự
chúng chạy:

| Chốt | Dòng | Nội dung |
|---|---|---|
| Mục tiêu vừa hồi sinh | 69-71 | `!Util.canDoWithTime(plTarget.lastTimeRevived, 1500)` → từ chối. **Bất tử 1,5 giây sau hồi sinh** |
| Đang dính hiệu ứng khống chế | 90-92 | `player.effectSkill.isHaveEffectSkill()` → từ chối |
| Không đủ mana | 96/107 | `canUseSkillWithMana(player)` |
| Chiêu chưa hồi | 96/107 | `canUseSkillWithCooldown(player)` |
| Mục tiêu đã chết | 106 | `mobTarget.isDie()` |
| Không được phép PK | 105 | `canAttackPlayer(player, plTarget)` |

Hồi chiêu — `services/SkillService.java:1273-1276`:

```
public boolean canUseSkillWithCooldown(Player player) {
    return Util.canDoWithTime(player.playerSkill.skillSelect.lastTimeUseThisSkill,
            player.playerSkill.skillSelect.coolDown - 50);
}
```

→ **Nới 50 ms** so với cooldown danh nghĩa, để bù độ trễ mạng. Con số này phải giữ
nguyên tinh thần trong bản Go (§4.3 của spec hoà giải mượt đã chốt nhịp 100 ms —
nới 50 ms là nửa nhịp, hợp lý).

Tầm đánh — `skill/Skill.java:13-14`:

```
public static final int RANGE_ATTACK_CHIEU_DAM = 100;      // chiêu cận chiến
public static final int RANGE_ATTACK_CHIEU_CHUONG = 300;   // chiêu chưởng/xa
```

Vượt tầm **không bị từ chối** mà bị đánh dấu `miss = true` → vẫn phát gói, sát thương 0
(`services/SkillService.java:578-584`). Đây là lựa chọn thiết kế có chủ ý: client thấy
hoạt ảnh đánh, chỉ không mất máu. Bản Go **nên giữ**, vì từ chối thẳng sẽ làm client
"nuốt" thao tác của người chơi, cảm giác lag.

Thể lực — `services/SkillService.java:540-552`: mỗi **500** đòn đánh trừ 1 thể lực;
hết thể lực thì báo "Thể lực đã cạn kiệt" và chặn đánh.

---

## 3. Gói XUỐNG — server → client

### 3.1 `54` — báo "người chơi X đánh quái Y"

`services/SkillService.java:1446-1460`:

```
1449  msg = new Message(54);
1450  msg.writer().writeInt((int) plAtt.id);
1451  msg.writer().writeByte(plAtt.playerSkill.skillSelect.skillId);
1452  msg.writer().writeByte(mob.id);
1453  Service.gI().sendMessAllPlayerInMap(plAtt, msg);
```

| # | Trường | Kiểu | Ý nghĩa |
|---|---|---|---|
| 1 | `idNguoiDanh` | int | |
| 2 | `skillId` | byte | `skillSelect.skillId` (**không** phải template.id) |
| 3 | `mobId` | byte | Chỉ số quái trong khu vực |

Gửi cho **tất cả** trong khu vực, kể cả người đánh. Gói này chỉ để **vẽ hoạt ảnh** —
máu quái đi bằng gói `-9`/`-12` riêng, gửi ngay sau.

### 3.2 `-60` — báo "người chơi X đánh người chơi Y"

`services/SkillService.java:1140-1150`:

```
1140  msg = new Message(-60);
1141  msg.writer().writeInt((int) plAtt.id);              // id pem
1142  msg.writer().writeByte(plAtt.playerSkill.skillSelect.skillId);  // skill pem
1143  msg.writer().writeByte(1);                          // số người pem
1144  msg.writer().writeInt((int) plInjure.id);           // id ăn pem
1145  msg.writer().writeByte(1);                          // read continue
1146  msg.writer().writeByte(0);                          // type skill
1147  msg.writeLongByGiaCatLuong(Util.maxIntValue(dameHit), cn.readInt);  // dame ăn
1148  msg.writer().writeBoolean(plInjure.isDie());        // is die
1149  msg.writer().writeBoolean(plAtt.nPoint.isCrit);     // crit
```

| # | Trường | Kiểu | Ý nghĩa |
|---|---|---|---|
| 1 | `idNguoiDanh` | int | |
| 2 | `skillId` | byte | |
| 3 | `soMucTieu` | byte | **Luôn ghi 1** ở bản gốc — nhưng tên biến cho thấy client hỗ trợ nhiều mục tiêu |
| 4 | `idNguoiTrung` | int | |
| 5 | **`docTiep`** | byte | **CỜ QUYẾT ĐỊNH ĐỘ DÀI GÓI** — xem cảnh báo dưới |
| 6 | `loaiChieu` | byte | ↓ **chỉ có khi `docTiep == 1`** |
| 7 | `satThuong` | int | Đã kẹp trần int32 |
| 8 | `daChet` | bool | |
| 9 | `chiMang` | bool | |

⚠️⚠️ **`-60` LÀ GÓI ĐỘ DÀI BIẾN THIÊN — đây là cái bẫy nguy hiểm nhất của KH6.**

Trường `docTiep` **không phải hằng số 1** như tên biến gốc ("read continue") gợi ý. Có
đúng hai dạng gói `-60`, dùng chung mã lệnh:

| Dạng | `docTiep` | Có trường 6-9? | Dùng khi | Bằng chứng |
|---|---|---|---|---|
| **Đòn đánh** | `1` | Có (**18 byte thân**) | Đánh gây sát thương | `SkillService.java:1145-1149` |
| **Chiêu hồi máu** | `0` | **Không** (**11 byte thân**) | Chiêu Trị Thương hồi máu/hồi sinh đồng đội | `SkillService.java:980-981` |

Cộng tay: `4 + 1 + 1 + 4 + 1` = **11** (dừng ở `docTiep`, tức `docTiep` là **byte thứ 11**),
rồi `+ 1 + 4 + 1 + 1` = **18**. Chênh lệch đúng **7 byte**.

Đoạn chứng minh — `services/SkillService.java:975-981`, chiêu Trị Thương:

```
975   msg = new Message(-60);
976   msg.writer().writeInt((int) player.id);      // id pem
977   msg.writer().writeByte(player.playerSkill.skillSelect.skillId);
978   msg.writer().writeByte(1);                   // số người pem
979   msg.writer().writeInt((int) pl.id);          // id ăn pem
980   msg.writer().writeByte(0);                   // read continue = 0
981   Service.gI().sendMessAllPlayerInMap(pl, msg);   // ← GỬI LUÔN, hết gói
```

Viết bản Go mà mặc định `docTiep = 1` rồi luôn ghi đủ 9 trường thì **chiêu hồi máu sẽ
đẩy thừa 7 byte vào gói**. Client vứt byte dư **im lặng, không lỗi, không triệu chứng**
(§10) — nghĩa là chạy thử sẽ **không** phát hiện được. Bản Go phải mô hình hoá `-60`
như **hai gói khác nhau**, không phải một gói có trường tuỳ chọn, và chặn bằng test byte.

⚠️ **Chưa chắc chắn (mức nhẹ):** trường 3 luôn = 1 nên không quan sát được vòng lặp
nhiều mục tiêu chạy thật ra sao — trường 9 (`chiMang`) nằm trong hay ngoài vòng lặp là
**suy đoán từ vị trí code**, chưa có bằng chứng từ client. Giai đoạn KH6 **chỉ ghi 1 mục
tiêu** (đúng như bản gốc) nên không ảnh hưởng; ghi lại để ai làm đánh diện rộng sau này
biết là **phải xác minh bằng client thật trước khi ghi > 1**.

Ngay sau gói này, bản gốc gọi `Service.gI().reload_HP_NV(plInjure)` (dòng 1151) → phát
tiếp gói cập nhật máu (§3.6).

### 3.3 `-9` — quái **còn sống** sau khi trúng đòn

`mob/Mob.java:1361-1374`:

```
1364  msg = new Message(-9);
1365  msg.writer().writeByte(this.id);
1366  msg.writeLongByGiaCatLuong(Util.maxIntValue(this.point.gethp()), cn.readInt);
1367  msg.writeLongByGiaCatLuong(Util.maxIntValue(dameHit), cn.readInt);
1368  msg.writer().writeBoolean(crit);
1369  msg.writer().writeInt(-1);
1370  Service.gI().sendMessAllPlayerInMap(this.zone, msg);
```

| # | Trường | Kiểu | Ý nghĩa |
|---|---|---|---|
| 1 | `mobId` | byte | |
| 2 | `hpConLai` | int | |
| 3 | `satThuong` | int | |
| 4 | `chiMang` | bool | |
| 5 | `hieuUngId` | **byte** | `-1` = không hiệu ứng. Bản Java ghi `int` 4 byte là **LỖI** — xem §6.3 |

### 3.4 `-12` — quái **chết** + rơi đồ

`mob/Mob.java:574-587` (khung) và `mob/Mob.java:589-612` (phần đồ rơi):

```
577  msg = new Message(-12);
578  msg.writer().writeByte(this.id);
579  msg.writeLongByGiaCatLuong(Util.maxIntValue(dameHit), cn.readInt);
580  msg.writer().writeBoolean(plKill.nPoint.isCrit);
581  List<ItemMap> items = mobReward(plKill, this.dropItemTask(plKill), msg);   // ← ghi tiếp vào msg
582  Service.gI().sendMessAllPlayerInMap(this.zone, msg);
...
600  msg.writer().writeByte(itemReward.size());        // sl item roi
601  for (ItemMap itemMap : itemReward) {
602      msg.writer().writeShort(itemMap.itemMapId);
603      msg.writer().writeShort(itemMap.itemTemplate.id);
604      msg.writer().writeShort(itemMap.x);
605      msg.writer().writeShort(itemMap.y);
606      msg.writer().writeInt((int) itemMap.playerId);
607  }
```

| # | Trường | Kiểu | Ý nghĩa |
|---|---|---|---|
| 1 | `mobId` | byte | |
| 2 | `satThuongDonCuoi` | int | |
| 3 | `chiMang` | bool | |
| 4 | `soDoRoi` | byte | |
| 5 | `itemMapId` | short | ↓ lặp `soDoRoi` lần |
| 6 | `templateId` | short | Loại vật phẩm |
| 7 | `x` | short | Toạ độ rơi |
| 8 | `y` | short | |
| 9 | `idNguoiDuocNhat` | int | Khoá đồ cho riêng người này một khoảng thời gian |

⚠️ **Bẫy nghiêm trọng ở bản gốc:** `mobReward` bọc toàn bộ trong `try/catch` in stack
trace rồi **trả về bình thường** (dòng 608-610). Nếu ném lỗi *sau khi* đã ghi
`writeByte(itemReward.size())` nhưng *trước khi* ghi đủ số món → gói phát ra **thiếu
byte**. Client đọc quá phần thân → ném lỗi → bị `catch` rỗng nuốt → **gói im lặng không
làm gì**: quái chết mà người chơi không thấy chết, đồ rơi không hiện, không log ở đâu cả
(§10). Bản Go **phải dựng xong danh sách đồ rơi trước, rồi mới bắt đầu ghi byte** —
không vừa tính vừa ghi.

**Phạm vi KH6:** hệ vật phẩm là KH7. Giai đoạn KH6 ghi `soDoRoi = 0` (quái chết, không
rơi gì) — gói vẫn đúng khuôn, chỉ chưa có nội dung.

### 3.5 `-13` — quái hồi sinh

`mob/Mob.java:516-524` (và bản sao ở 1376-1390, 1392-1401 cho các chế độ đặc biệt):

```
519  msg = new Message(-13);
520  msg.writer().writeByte(this.id);
521  msg.writer().writeByte(this.tempId);
522  msg.writer().writeByte(lvMob());
523  msg.writeLongByGiaCatLuong(Util.maxIntValue(this.point.hp), cn.readInt);
```

| # | Trường | Kiểu | Ý nghĩa |
|---|---|---|---|
| 1 | `mobId` | byte | Chỉ số trong khu vực |
| 2 | `tempId` | byte | Loại quái (tra bảng mẫu quái) |
| 3 | `lvMob` | byte | |
| 4 | `hp` | int | |

### 3.6 `-30 / 9` — cập nhật máu người chơi khác

`services/Service.java:542-556`:

```
545  msg = messageSubCommand((byte) 9);
546  msg.writer().writeInt((int) pl.id);
547  msg.writeLongByGiaCatLuong(Util.maxIntValue(pl.nPoint.hp), cn.readInt);
548  msg.writeLongByGiaCatLuong(Util.maxIntValue(pl.nPoint.hpMax), cn.readInt);
549  sendMessAnotherNotMeInMap(pl, msg);
```

Đây là **lệnh `-30` sub-command 9** — cùng họ với `-30/0` (ME_INFO) mà KH4 đã làm.

| # | Trường | Kiểu |
|---|---|---|
| 1 | `idNguoiChoi` | int |
| 2 | `hp` | int |
| 3 | `hpMax` | int |

Gửi cho **người khác**, không gửi cho chính chủ (chính chủ nhận qua `sendInfoHpMpMoney`).

### 3.7 `-30 / 15` — người chơi hồi sinh (báo cho cả khu vực)

`services/Service.java:1122-1129`:

```
1122  msg = messageSubCommand((byte) 15);
1123  msg.writer().writeInt((int) pl.id);
1124  msg.writeLongByGiaCatLuong(Util.maxIntValue(hp), cn.readInt);
1125  msg.writeLongByGiaCatLuong(Util.maxIntValue(mp), cn.readInt);
1126  msg.writer().writeShort(pl.location.x);
1127  msg.writer().writeShort(pl.location.y);
1128  sendMessAllPlayerInMap(pl, msg);
```

| # | Trường | Kiểu |
|---|---|---|
| 1 | `idNguoiChoi` | int |
| 2 | `hp` | int |
| 3 | `mp` | int |
| 4 | `x` | short |
| 5 | `y` | short |

### 3.8 `-16` — báo cho **chính chủ** là đã hồi sinh

`services/Service.java:1116-1118` — **gói rỗng**, không có thân:

```
1116  msg = new Message(-16);
1117  pl.sendMessage(msg);
```

Trùng mã lệnh với gói lên `-16` (§1.4) nhưng khác chiều. Bình thường trong giao thức này.

### 3.9 `-17` — báo cho **chính chủ** là đã chết

`services/Service.java:1146-1152`:

```
1146  msg = new Message(-17);
1147  msg.writer().writeByte((int) pl.id);      // ⚠️ xem §6
1148  msg.writer().writeShort(pl.location.x);
1149  msg.writer().writeShort(pl.location.y);
1150  // msg.writer().writeLong(-1); send Power
1151  pl.sendMessage(msg);
```

### 3.10 `-8` — báo cho **người khác** là ai đó đã chết

Hai chỗ phát gói này, **ghi trường đầu KHÁC KIỂU NHAU** — xem §6.

`map/Zone.java:723-732` (khi người mới vào khu vực, đồng bộ ai đang nằm chết):

```
725  msg = new Message(-8);
726  msg.writer().writeInt((int) plInfo.id);
727  msg.writer().writeByte(0);
728  msg.writer().writeShort(plInfo.location.x);
729  msg.writer().writeShort(plInfo.location.y);
```

`services/Service.java:1156-1162` (khi có người vừa chết):

```
1156  msg = new Message(-8);
1157  msg.writer().writeShort((int) pl.id);     // ⚠️ SHORT, không phải INT
1158  msg.writer().writeByte(cPk);              // cPk luôn = 0
1159  msg.writer().writeShort(pl.location.x);
1160  msg.writer().writeShort(pl.location.y);
```

### 3.11 `-11` — quái đánh trúng mình (chỉ gửi cho nạn nhân)

`mob/Mob.java:481-485`:

```
481  msg = new Message(-11);
482  msg.writer().writeByte(this.id);
483  msg.writeLongByGiaCatLuong(Util.maxIntValue(dame), cn.readInt);
484  player.sendMessage(msg);
```

| # | Trường | Kiểu |
|---|---|---|
| 1 | `mobId` | byte |
| 2 | `satThuong` | int |

⚠️ Có **bản sao thứ hai** ở `mob/Mob.java:88-92` ghi `writeInt(dame)` trực tiếp thay vì
qua `writeLongByGiaCatLuong`. Với `cn.readInt = true` thì **hai bản ra byte giống hệt
nhau**, nên không phải lỗi thật — chỉ là code trùng lặp. Ghi lại để người sau đọc code
cũ không hoang mang.

### 3.12 `-10` — quái đánh trúng ai đó (gửi cho người khác trong khu vực)

`mob/Mob.java:491-499`:

```
494  msg = new Message(-10);
495  msg.writer().writeByte(this.id);
496  msg.writer().writeInt((int) player.id);
497  msg.writeLongByGiaCatLuong(Util.maxIntValue(player.nPoint.hp), cn.readInt);
498  Service.gI().sendMessAnotherNotMeInMap(player, msg);
```

| # | Trường | Kiểu | Ý nghĩa |
|---|---|---|---|
| 1 | `mobId` | byte | |
| 2 | `idNanNhan` | int | |
| 3 | `hpConLai` | int | Máu còn lại, **không phải** sát thương |

### 3.13 `-45 / 4` và `-45 / 7` — khởi chiêu (báo trước hoạt ảnh)

`services/SkillService.java:1208-1223` (khởi chiêu thường) và `1225-1240` (khởi bom):

```
1211  msg = new Message(-45);
1212  msg.writer().writeByte(4);          // 7 với bom
1213  msg.writer().writeInt((int) player.id);
1214  msg.writer().writeShort(player.playerSkill.skillSelect.skillId);
1215  msg.writer().writeShort(affterMiliseconds);
```

| # | Trường | Kiểu | Ý nghĩa |
|---|---|---|---|
| 1 | `subCmd` | byte | `4` = khởi chiêu, `7` = khởi bom |
| 2 | `idNguoiChoi` | int | |
| 3 | `skillId` | short | ⚠️ **short** ở đây, trong khi gói `54`/`-60` ghi **byte** |
| 4 | `treMiliGiay` | short | Bao nhiêu ms nữa thì chiêu nổ |

Đây chính là mảnh ghép cho ràng buộc "client dự đoán trước, server xác nhận sau" ở KH10:
server báo trước mốc thời gian, client chạy hoạt ảnh, server chốt kết quả sau.

**Bảng đầy đủ sub-command của `-45`** (gom từ `services/EffectSkillService.java` và
`services/SkillService.java`). Byte đầu là sub-command, ba trường sau gần như luôn là
`int idNguoiChoi` + `short skillId`:

| sub | Ý nghĩa | Thân gói sau byte sub | Nguồn |
|---|---|---|---|
| `0` | Choáng diện rộng (Thái Dương Hạ San) | int id, short skillId, **byte soQuai** + lặp{byte mobId, byte giay}, **byte soNguoi** + lặp{int plId, byte giay} | `EffectSkillService.java:550-563` |
| `1` | Đang gồng | int id, short skillId | `:592-595` |
| `3` | Dừng gồng | int id, **short −1** | `:606-609` |
| `4` | Khởi chiêu | int id, short skillId, **short treMs** | `SkillService.java:1211-1215` |
| `5` | Nổ kết thúc gồng | int id, short skillId | `:621-624` |
| `6` | Bắt đầu gồng | int id, short skillId | `:576-579` |
| `7` | Khởi bom | int id, short skillId, **short treMs** | `SkillService.java:1228-1232` |
| `8` | Dùng chiêu (hiệu ứng chung) | int id, short skillId | `:50-53` |

Sub `0` cũng là gói **độ dài biến thiên** (hai vòng lặp lồng nhau) — cùng loại bẫy với `-60`.
Sub `2` chưa tìm thấy chỗ phát; có thể không dùng hoặc nằm ở phần chưa đọc.

### 3.14 `-124` — bật/tắt hiệu ứng trạng thái (trói, choáng, tàng hình…)

`services/EffectSkillService.java:61-79` (cho người chơi) và `:81-95` (cho quái).
Gói này **ba dạng khác nhau** tuỳ hai byte đầu:

| Dạng | `toggle` | `loaiDoiTuong` | Thân còn lại |
|---|---|---|---|
| Tắt **hết** hiệu ứng của một người | `TURN_OFF_ALL_EFFECT` | `0` | int idNguoiDinh |
| Bật/tắt một hiệu ứng trên người | khác | `0` | byte loaiHieuUng, int idNguoiDinh, int idNguoiGay |
| Bật/tắt một hiệu ứng trên quái | bất kỳ | `1` | byte loaiHieuUng, **byte** mobId, int idNguoiGay |

Chú ý: id đối tượng là **int** với người chơi nhưng **byte** với quái — nhất quán với
mọi gói khác (quái đánh số trong phạm vi khu vực nên 1 byte là đủ).

Lại thêm một gói độ dài biến thiên. **Ba gói `-60`, `-45/0`, `-124` chiếm phần lớn rủi
ro lệch dòng byte của KH6** — kế hoạch thi công phải có test byte-chính-xác riêng cho
từng dạng, không gộp.

### 3.15 Danh sách quái trong gói MAP INFO (`-24`)

`map/Zone.java:779-800` — phần này KH4 đã làm nhưng khi ấy chưa có quái nên ghi 0. Bảng
đầy đủ để KH6 điền vào:

| # | Trường | Kiểu | Ý nghĩa |
|---|---|---|---|
| 0 | `soQuai` | byte | |
| 1 | `isDisable` | bool | ↓ lặp `soQuai` lần. Bản gốc **luôn ghi false** cho cả 5 cờ |
| 2 | `isDontMove` | bool | |
| 3 | `isFire` | bool | |
| 4 | `isIce` | bool | |
| 5 | `isWind` | bool | |
| 6 | `tempId` | byte | Loại quái |
| 7 | `sys` | byte | Luôn 0 |
| 8 | `hp` | int | |
| 9 | `level` | byte | |
| 10 | `hpMax` | int | |
| 11 | `x` | short | |
| 12 | `y` | short | |
| 13 | `status` | byte | |
| 14 | `lvMob` | byte | |
| 15 | `isBigBoss` | bool | |

Lưu ý thứ tự lạ: **`hp` (8) đứng trước `level` (9), rồi mới tới `hpMax` (10)** — không
phải hp/hpMax liền nhau. Chép sai thứ tự này là lệch dòng.

---

## 4. Công thức sát thương của bản gốc (tham khảo, KHÔNG chép nguyên)

### 4.1 Quái đánh người

`mob/MobPoint.java:42-46`:

```
public long getDameAttack() {
    return this.dame != 0
        ? this.dame + Util.nextLong(-(this.dame / 100), (this.dame / 100))
        : this.getHpFull() * Util.nextInt(mob.pDame - 1, mob.pDame + 1) / 100
              + Util.nextInt(-(mob.level * 10), mob.level * 10);
}
```

Hai chế độ: nếu quái có `dame` riêng thì dùng nó ±1 %; nếu không thì suy ra từ **máu tối
đa** × `pDame` % ± nhiễu theo cấp.

### 4.2 Người đánh quái / đánh người

`player/NPoint.java:1715` trở đi. Cấu trúc: `dame` gốc của nhân vật, rồi cộng dồn
`percentDameSkill` (theo chiêu), `percentDameIntrinsic` (nội tại), `percentXDame` (bộ đồ
đủ 5 món)…

**Đánh giá:** hàm này dài hơn 200 dòng, là một `switch` khổng lồ nhét cứng từng chiêu
Dragon Ball (Kamejoko, Galick, Masenko, Kaioken, Makankosappo…) trộn với nội tại, bộ đồ,
bùa, thú cưng. Nó **không phải công thức**, nó là một đống ngoại lệ.

→ **Khuyến nghị mạnh cho spec: KHÔNG chép cái này.** Game Naruto viết mới nên có công
thức sát thương **riêng, ngắn, kiểm thử được**, và mọi thứ đặc thù từng chiêu đẩy sang
**dữ liệu** (chuẩn bị sẵn cho KH8 — kho nội dung YAML) chứ không nhét vào code. Lý do
không chỉ là gọn: hàm gốc **không thể viết test tử tế** vì mỗi nhánh phụ thuộc trạng
thái toàn cục khác nhau — mà §6.11 của `BAT-DAU-TU-DAY.md` nói rõ dự án này đã ba lần
dính "test xanh mà code sai".

Cái **phải giữ** từ bản gốc chỉ là ba thứ, vì client trông vào chúng:
1. Trần `int32` (§0.1).
2. Cờ chí mạng (bool) đi kèm mỗi đòn.
3. Cờ "trượt" = sát thương 0 nhưng **vẫn phát gói** (§2).

### 4.3 Chốt chặn máu quái

`mob/Mob.java:139-152`:

```
139  public void injured(Player plAtt, long damage, boolean dieWhenHpFull) {
140      if (!this.isDie()) {
141          if (damage >= this.point.hp) {
142              damage = this.point.hp;          // không đánh quá máu còn lại
143          }
144          if (!dieWhenHpFull) {
145              if (this.point.hp == this.point.maxHp && damage >= this.point.hp) {
146                  damage = this.point.hp - 1;  // đầy máu thì không chết ngay đòn đầu
147              }
```

Hai quy tắc đáng giữ: **sát thương không vượt máu còn lại** (để client không hiện số âm),
và **quái đang đầy máu không chết bởi một đòn** (chống one-shot làm hỏng cảm giác chơi).

---

## 5. Ba thứ ở bản gốc KHÔNG được chép lại

1. **`MobPoint.sethp` kiểm tra nhầm biến** (`mob/MobPoint.java:34-40`): kiểm
   `this.hp < 0` (giá trị **cũ**) rồi mới gán giá trị mới → không bao giờ kẹp được sàn 0.
   Bản gốc thoát nạn nhờ `Mob.injured` gán thẳng `this.point.hp -= damage`, bỏ qua setter.
   Bản Go phải kẹp `[0, hpMax]` ở **một chỗ duy nhất**.

2. **Nuốt lỗi giữa lúc ghi gói** (`mob/Mob.java:589-612`, và khắp nơi trong `Service.java`):
   `try/catch` rỗng ôm cả đoạn ghi byte. Lỗi ở giữa = gói thiếu byte = **gói đó**
   gói đó **im lặng không làm gì** ở phía client (client ném lỗi rồi bị `catch` nuốt — §10),
   **không có log ở cả hai đầu**. Bản Go: dựng xong dữ liệu → ghi gói một mạch → lỗi thì
   **bỏ cả gói**, không gửi nửa vời.

3. **Sát thương nhét cứng theo mapId** (`mob/Mob.java:180-189`): `mapId == 155/156/157`
   → `damage = 1`; map Ngũ Hành Sơn → `damage = 80_000`. Đây là luật chơi giấu trong code.
   Bản Naruto: đưa vào **dữ liệu khu vực**, không phải `if` theo số.

---

## 6. Ba chỗ từng CHƯA CHẮC — nay ĐÃ GỠ bằng client (26/07/2026)

Đã mở `D:\NRO_Project\Assets\Scripts\Game1\Controller.cs` đọc tận nơi. **Client là bên
đọc nên client là trọng tài duy nhất.** Cả ba đều lộ ra **lỗi thật của bản Java**.

*(`Game2/Controller.cs` giống hệt `Game1`, chỉ khác `namespace`.)*

### 6.1 Gói `-8`: `writeInt` ĐÚNG, `writeShort` là LỖI

```csharp
// Controller.cs:4548-4554
case -8:
    @char = GameScr.findCharInMap(msg.reader().readInt());   // ← readInt
    if (@char != null) {
        @char.cPk = msg.reader().readByte();
        @char.waitToDie(msg.reader().readShort(), msg.reader().readShort());
    }
    break;
```

`map/Zone.java:726` (`writeInt`) đúng. `services/Service.java:1157` (`writeShort`) **sai
thật**: client lấy nhầm id → `findCharInMap` trả `null` → cả gói bị bỏ. **Hậu quả trong
game đang chạy: người chơi khác KHÔNG thấy ai đó vừa chết**, cho tới khi họ ra vào khu
vực lại (lúc đó `Zone.java` phát bản đúng).

### 6.2 Gói `-17`: KHÔNG CÓ trường id — byte đầu là `cPk`

```csharp
// Controller.cs:4533-4547
case -17:
    Char.myCharz().meDead = true;
    Char.myCharz().cPk = msg.reader().readByte();                               // ← cPk
    Char.myCharz().startDie(msg.reader().readShort(), msg.reader().readShort());
    try {
        Char.myCharz().cPower = msg.reader().readLong();                        // ← tuỳ chọn
        Char.myCharz().applyCharLevelPercent();
    } catch (Exception) { Cout.println("Loi tai ME_DIE " + msg.command); }
    Char.myCharz().countKill = 0;
    break;
```

Gói chỉ gửi cho chính chủ nên client dùng thẳng `Char.myCharz()`, **không cần và không
đọc id**. `services/Service.java:1147` viết `writeByte((int) pl.id)` là **nhét id nhân
vật vào ô `cPk`** — sai dữ liệu (không lệch byte vì vẫn 1 byte). Dòng bị comment
`// writeLong(-1); send Power` ở `Service.java:1150` khớp trường `cPower` tuỳ chọn.

**Khuôn đúng:** `byte cPk, short x, short y` (+ `long power` tuỳ chọn, bản Go bỏ).

### 6.3 Trường cuối gói `-9`: là `byte hieuUngId`, KHÔNG phải `int`

```csharp
// Controller.cs:4352-4356
sbyte b78 = msg.reader().readByte();            // ← 1 BYTE
if (b78 != -1) {
    EffecMn.addEff(new Effect(b78, mob11.x, mob11.getY(), 3, 1, -1));
}
```

Trường bí ẩn là **id hiệu ứng hình ảnh, 1 byte, `-1` = không có hiệu ứng**.
`mob/Mob.java:1369` ghi `writeInt(-1)` = **4 byte** → client đọc byte đầu `0xFF` = −1
→ "không hiệu ứng" → **tình cờ đúng ý**, 3 byte sau bị vứt. Bản Go ghi **1 byte `-1`**.

**Hai điều khác lộ ra khi đọc `-9`:**

- **Client BỎ QUA phần còn lại nếu `dame == 1`** (`Controller.cs:4338`, `if (dame2 != 1.0)`).
  Sát thương đúng bằng 1 là một giá trị **đặc biệt** ở client — không vẽ gì cả.
  Bản Go nên tránh sinh ra sát thương đúng bằng 1, hoặc chấp nhận là nó vô hình.
- **`dame == 0` là "trượt"**: client hiện chữ *miss* và **kéo quái về vị trí gốc**
  (`mob.x = mob.xFirst`). Khớp đúng thiết kế "vượt tầm → sát thương 0 nhưng vẫn phát gói".

---

## 7. Nhịp máy chủ, AI quái và hồi sinh quái

### 7.1 Nhịp 100 ms — trùng khớp với KH5

`server/GameLoop.java` (cả file chỉ 36 dòng):

```
8   private static final int TICK = 100; // ms
...
14  GameExecutor.EXECUTOR.scheduleAtFixedRate(() -> {
...
25          for (Zone zone : map.zones) {
27              zone.update(now);
...
34  }, 0, TICK, java.util.concurrent.TimeUnit.MILLISECONDS);
```

→ **Server Java gốc cũng chạy nhịp 100 ms, quét từng khu vực một.** Đây là xác nhận
độc lập cho con số đã chốt ở §4.3 của
[`../specs/2026-07-23-hoa-giai-muot-va-chong-gian-lan.md`](../specs/2026-07-23-hoa-giai-muot-va-chong-gian-lan.md)
và đã hiện thực ở KH5 (actor khu vực, ticker 100 ms, commit `27d66fc`).

Hệ quả kiến trúc: `Zone.update()` gọi `mob.update()` cho từng con quái
(`map/Zone.java:148-151`). Nghĩa là **toàn bộ AI quái chạy trong cùng luồng với khu vực**.
Bản Go đã có sẵn khung này — quái chỉ việc trở thành một phần trạng thái của actor khu
vực, cập nhật trong `vongLap`. **Không cần luồng mới, không cần khoá.** Đây đúng là luật
§6.9 của `BAT-DAU-TU-DAY.md`: chỉ luồng `vongLap` của khu vực được chạm trạng thái khu vực.

### 7.2 Hồi sinh quái và hồi máu quái

`mob/Mob.java:320-361`:

| Việc | Điều kiện | Con số |
|---|---|---|
| Hồi sinh sau khi chết | `Util.canDoWithTime(lastTimeDie, 5000)` | **5 giây** |
| Hồi máu khi chưa đầy | `Util.canDoWithTime(lastTimePhucHoi, 30000)` | mỗi **30 giây**, hồi **hpMax/10** |

Map đặc biệt (doanh trại, bản đồ kho báu, con đường rắn độc, khí gas, tay Karin) có luật
riêng — bản Naruto không có các map này, bỏ qua.

### 7.3 Quái chọn mục tiêu thế nào

Hai tầng ưu tiên (`mob/Mob.java:393-443`):

| Ưu tiên | Hàm | Chọn ai | Tầm | Nhịp đánh |
|---|---|---|---|---|
| 1 | `getFirstPlayerCanAttack` | **Người đã từng đánh nó** (`temporaryEnemies`), gần nhất | 300 | `timeAttack = 1000` ms |
| 2 | `getPlayerCanAttack` | Người bất kỳ gần nhất | 100 | `timeAttack = 2000` ms |

Nghĩa là: **quái trả thù nhanh gấp đôi và với tầm gấp ba so với khi nó tự gây sự.**
Đây là luật chơi tốt, đáng giữ nguyên: người chơi đi ngang qua thì tương đối an toàn,
đã đánh thì phải chịu trận.

`temporaryEnemies` được nạp ở `mob/Mob.java:195` (`addTemporaryEnemies(plAtt)` mỗi lần
trúng đòn) và xoá sạch khi quái chết (dòng 199).

### 7.4 Ba lỗi trong AI quái — không chép lại

1. **Đếm ngược nhịp đánh chạy cả khi không có mục tiêu** (`mob/Mob.java:386-390`):
   `lastTimeAttackPlayer = System.currentTimeMillis()` nằm **ngoài** `if (player != null)`.
   Kết quả: quái không có ai để đánh vẫn "tiêu" mất nhịp, nên lúc người chơi vừa bước
   vào tầm, quái phải chờ thêm một nhịp nữa. Bản Go: chỉ đặt lại mốc thời gian **khi
   thật sự đánh**.

2. **Thứ tự `&&` / `||` sai** (`mob/Mob.java:402-404`): điều kiện viết là
   `A && B && ... && (C || D) || isBigBoss()`. Trong Java `&&` ràng chặt hơn `||`, nên
   nó thành `(A && B && ... && (C||D)) || isBigBoss()` — với big boss thì **mọi bộ lọc
   phía trước bị bỏ qua**, kể cả bộ lọc "người chơi chưa chết". Big boss sẽ nhắm cả xác
   chết. Bản Go: viết ngoặc tường minh, và có test cho trường hợp mục tiêu đã chết.

3. **Đặt `timeAttack` như tác dụng phụ của hàm "tìm mục tiêu"** (dòng 412 và 438):
   hàm tên là *lấy người chơi* nhưng lại lặng lẽ đổi nhịp đánh của quái. Bản Go: hàm tìm
   mục tiêu **trả về cả mục tiêu lẫn nhịp**, không sửa trạng thái ngầm.

### 7.5 Quái đánh người — luồng phát gói

`mob/Mob.java:445-475`, đúng thứ tự:

1. Tính sát thương (`MobPoint.getDameAttack`, §4.1) rồi trừ máu người chơi.
2. Phát `-11` cho **nạn nhân** (§3.11).
3. Phát `-10` cho **những người khác** trong khu vực (§3.12).
4. Phản sát thương (nếu nạn nhân có nội tại phản đòn).

Lưu ý ở bước 1: nếu quái có `lvMob > 0` (quái tinh anh) thì sát thương **không** theo
công thức mà bị ghi đè thành `hpMax của người chơi × 10 %` (dòng 457-459) — tức đánh
theo phần trăm máu, luôn 10 đòn là chết bất kể trang bị. Bản Naruto muốn giữ ý này thì
đưa vào **dữ liệu quái** (một cờ "sát thương theo % máu"), không nhét `if` vào code.

---

## 8. Những gì tài liệu này CHƯA phủ

Ghi ra để không ai tưởng đã đủ:

**Bốn loại chiêu** (`template.type`, phân nhánh ở `SkillService.java:123-133`):

| type | Loại | Hàm xử lý | Đã đọc? |
|---|---|---|---|
| 1 | Chiêu đánh (có mục tiêu) | `useSkillAttack` | ✅ §2, §3.1, §3.2 |
| 2 | Chiêu buff lên người khác | `useSkillBuffToPlayer` | ⚠️ mới đọc nhánh Trị Thương |
| 3 | Chiêu tự thân | `useSkillAlone` | ⚠️ mới đọc nhánh Thái Dương Hạ San |
| 4 | Chiêu không cần nhắm | `useNewSkillNotFocus` | ❌ chưa đọc |

Còn thiếu:

- **Chiêu type 4** (`useNewSkillNotFocus`): là nhánh duy nhất đọc thêm dữ liệu từ gói
  lên (`status == 20` → byte skillId, short dx, short dy, byte dir, short x, short y —
  `SkillService.java:79-89`). **Cần biết gói lên nào mang `status == 20`** trước khi
  làm chiêu type 4. **Có thể hoãn khỏi KH6.**
- **Các nhánh còn lại của type 2 / type 3**: khiên, tăng sức, biến hình, tàng hình.
- **Hệ hiệu ứng trạng thái đầy đủ** (`EffectSkillService.java`, 701 dòng): đã có khuôn
  gói `-124` (§3.14) nhưng **chưa có bảng mã `loaiHieuUng`**. Cần trước khi làm choáng/trói.
- **Quái di chuyển**: §7 mới phủ *chọn mục tiêu* và *đánh*; phần quái **đi lại** chưa đọc.
  **Cần cho KH6.**
- **Boss** (`boss/`, 12 file): hệ riêng, phức tạp hơn quái thường nhiều. **Nên tách khỏi
  KH6** — nếu gộp vào thì kế hoạch phình gấp đôi.
- **Rơi đồ chi tiết** (`getItemMobReward`, tỷ lệ rơi): thuộc KH7.
- **PvP / cờ `typePk`**: `canAttackPlayer` mới đọc phần khung.

**Đủ để bắt đầu KH6:** vòng chiến đấu cơ bản (đánh quái, đánh người, quái đánh lại, chết,
hồi sinh, hồi sinh quái) — §1 đến §7 đã phủ kín. Phần chiêu đặc biệt và hiệu ứng trạng
thái cần một đợt nghiên cứu nữa, nên tách thành giai đoạn sau trong cùng kế hoạch.

---

## 9. Gói `87` — bản Java gửi nhầm mã lệnh

Sau `-13` (quái hồi sinh), `mob/Mob.java:525` gọi tiếp `sendMobMaxHp(this.point.hp)`:

```
mob/Mob.java:1480-1490
public void sendMobMaxHp(long maxHp) {
    msg = new Message(87);
    msg.writer().writeByte(this.id);
    msg.writeLongByGiaCatLuong(Util.maxIntValue(maxHp), cn.readInt);   // 4 byte
    Service.gI().sendMessAllPlayerInMap(this.zone, msg);
```

Nhưng mã `87` ở client **không phải** lệnh máu:

```csharp
// Controller.cs:3016-3018
case 87:
    ((Mob)GameScr.vMob.elementAt(msg.reader().readUnsignedByte())).isWind = msg.reader().readBool();
    break;
```

**`87` = "quái bị gió thổi", thân `byte mobId + bool isWind`.** Server gửi 5 byte, client
đọc 2. Và byte thứ hai client đọc thành `isWind` chính là **byte cao nhất của `int maxHp`**:

| maxHp | Byte cao | `isWind` | Hậu quả |
|---|---|---|---|
| < 16.777.216 | 0 | `false` | Vô hại |
| ≥ 16.777.216 | ≠ 0 | **`true`** | `Mob.cs:1172` → **quái đứng im, không di chuyển** |

Bản Java dùng **một mã cho hai việc**: `sendWind` (`Mob.java:1468-1478`, khớp client
nhưng **không chỗ nào gọi**) và `sendMobMaxHp` (**không khớp client**). `sendMobMaxHp` gọi 1 lần trực tiếp ở
`Mob.java:525`, và gián tiếp qua `sendMobHoiSinh()` từ **cả hai** nhánh: hồi sinh 5 giây
(`Mob.java:349`) và hồi máu 30 giây (`Mob.java:357`) — tức nó **bắn ra mỗi lần bất kỳ con
quái nào hồi sinh**, không phải trường hợp hiếm.

**Mã đúng để gửi maxHp là `-84`** (`Controller.cs:767-783`, `byte mobIndex + int maxHp`) —
nhưng tra khắp `1-MAY-CHU-GAME/src/` **không có chỗ nào tạo `Message(-84)`**.

→ **Bản Go không gửi `87`.** Client đã tự đặt `maxHp = hp` khi xử lý `-13`
(`Controller.cs:4294-4295`), nên không thiếu gì cả.

---

## 10. Khung gói: sai trường KHÔNG làm hỏng dòng byte

Điều này sửa lại một cách nói sai trong bản đầu của tài liệu này.

```csharp
// Session_ME.cs:181-219 (rút gọn)
sbyte b = dis.ReadSByte();                        // ma lenh
if (b == -32 || b == -66 || b == 11 || b == -67
    || b == -74 || b == -87 || b == 66) return readMessage2(b);   // khung 3 byte
sbyte b2 = dis.ReadSByte();
sbyte b3 = dis.ReadSByte();
num = ((readKey(b2) & 0xFF) << 8) | (readKey(b3) & 0xFF);          // do dai 2 byte
sbyte[] array = new sbyte[num];
Buffer.BlockCopy(dis.ReadBytes(num), 0, array, 0, num);            // doc TRON goi
return new Message(b, array);                                      // roi moi phan tich
```

Client đọc **trọn** thân gói theo độ dài rồi mới phân tích mảng riêng đó. Vậy nên:

| Kiểu sai | Chuyện gì xảy ra | Triệu chứng |
|---|---|---|
| Ghi **thiếu** | Client đọc quá → ném lỗi → bị `catch (Exception) {}` (`Controller.cs:4600`) nuốt | Gói **im lặng không làm gì** |
| Ghi **thừa** | Byte dư bị vứt lúc `msg.cleanup()` | **Không triệu chứng gì** |
| Sai **thứ tự**, độ dài vẫn khớp | Giá trị rác vào đúng ô | Số liệu vô lý |

**Không có chuyện "client hỏng vĩnh viễn" vì sai trường.** Cái làm hỏng cả dòng byte là
sai **số byte của trường độ dài** (2 hay 3) — lỗi tầng khung, và danh sách 7 mã dùng
khung 3 byte nằm ngay ở `Session_ME.cs:190`: `-32, -66, 11, -67, -74, -87, 66`.
**Không mã chiến đấu nào trong đó** → KH6 dùng khung 2 byte cho tất cả.

Điều nguy hiểm thật không phải "hỏng to" mà là **hỏng im lặng ở cả hai đầu**: server
không log, client nuốt lỗi. Chạy thử bằng mắt sẽ **không** phát hiện được. Đó là lý do
kế hoạch bắt buộc test byte-chính-xác + kiểm thử đột biến cho từng dạng gói.
