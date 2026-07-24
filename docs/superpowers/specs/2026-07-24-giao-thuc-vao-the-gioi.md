# Giao thức vào thế giới — tham chiếu byte

**Mục đích:** nguồn sự thật cho Kế hoạch 3 (đồng bộ dữ liệu tĩnh) và Kế hoạch 4
(vào thế giới). Mọi bảng byte dưới đây **đã giải mã hai chiều** (Java ghi ↔ Unity
đọc) và **kiểm chéo đối kháng**. Đừng viết lại từ trí nhớ — trích từ đây.

**Phạm vi:** đưa một nhân vật đã đăng nhập vào một bản đồ để client Unity thật
rời màn "Xin chờ" và hiện lên. Client giữ nguyên (`D:\NRO_Project`), không sửa.

---

## 0. Chuỗi tối thiểu tuyệt đối

Sau khi đăng nhập có nhân vật (`vaoTheGioi`), thứ tự gói máy chủ gửi/nhận:

```
Máy chủ → Client:  -28/4  versionGame        (BẮT BUỘC — kích client so phiên bản)
Máy chủ → Client:  -82    tileSetInfo         (BẮT BUỘC trước -24 — thiếu = client THOÁT app)

  [Client so vsData/vsMap/vsSkill/vsItem với cache; client MỚI luôn xin lại hết]

Client → Máy chủ:  -87    (xin data)     → Máy chủ trả -87  (part + 5 blob rỗng)
Client → Máy chủ:  -28/6  (xin map)      → Máy chủ trả -28/6 (catalog map, rỗng được)
Client → Máy chủ:  -28/7  (xin skill)    → Máy chủ trả -28/7 (5 LÀNG NARUTO)
Client → Máy chủ:  -28/8  (xin item)     → Máy chủ trả -28/8 (subtype 2 rỗng)

Client → Máy chủ:  -28/13 (client ok)         (chỉ gửi khi ĐỦ CẢ 4 phiên bản khớp)

Máy chủ → Client:  -30/0  ME_INFO             (BẮT BUỘC)
Máy chủ → Client:  -24    MAP INFO            (BẮT BUỘC — gói kích switchToMe, rời "Xin chờ")
```

**Bỏ hẳn** (đã kiểm chứng vô hại): `-77 sendSmallVersion` (32KB dữ liệu chết,
không ai đọc lại), `-93 sendBgItemVersion` (có bug indexing trên cache RMS),
`-31 sendDataItemBG` (chỉ vật thể nền trang trí — gửi `short 0` nếu muốn).

**Ranh giới kế hoạch:**
- **Kế hoạch 3** = `-28/4` + phục vụ 4 catalog (`-87`, `-28/6`, `-28/7`, `-28/8`)
  → client gửi `-28/13`.
- **Kế hoạch 4** = `-82` + `-30/0` + `-24` + mô hình khu vực → client render.

---

## 1. Phiên bản dữ liệu tĩnh

Client so 4 phiên bản `vsData/vsMap/vsSkill/vsItem` (máy chủ gửi trong `-28/4`)
với cache `vc*` của nó. Lệch → xin cập nhật. `-28/13` chỉ gửi khi **cả 4 khớp**.

**Bẫy:** client MỚI (RMS trống) **luôn** xin lại cả 4 dù phiên bản có khớp hay
không — vì nhánh "khớp → nạp từ RMS" ném lỗi trên RMS rỗng rồi ép xin lại. Nên
máy chủ **phải phục vụ đủ cả 4** ở lần đăng nhập đầu, không né được.

**Phiên bản Naruto phải KHÁC bản Dragon Ball cũ** để ghi đè cache client đang có:
cũ là `vsData=10, vsMap=2, vsSkill=1, vsItem=4` (`DataGame.java:44-47`). Đặc biệt
`vsSkill` cũ = 1 — nếu ta gửi 1, client dùng cache 3 lớp Dragon Ball, KHÔNG lấy
5 làng của ta.

> **Chốt:** dùng `vsData=vsMap=vsSkill=vsItem = 20` cho bản Naruto (mỗi giá trị
> đều ≠ giá trị cũ và ≠ 0). Bump khi đổi nội dung tương ứng. Kiểu byte (client
> đọc `sbyte`), 20 nằm trong [0,127].

---

## 2. `-28/4` versionGame

Ghi: `DataGame.java:54-72`. Đọc: `Controller.cs:5441-5578`.

> Gói dùng `messageNotMap((byte)4)` = `Message(-28)` + byte lệnh con `4` **là byte
> đầu thân**. Mọi gói `-28/X` đều vậy (byte lệnh con đầu thân). Riêng `-87` là
> `Message(-87)` trực tiếp, KHÔNG có byte lệnh con.

| # | Trường | Kiểu | Giá trị |
|---|---|---|---|
| 0 | **lệnh con** | byte | **4** |
| 1 | vsData | byte | 20 |
| 2 | vsMap | byte | 20 |
| 3 | vsSkill | byte | 20 |
| 4 | vsItem | byte | 20 |
| 5 | đệm | byte | 0 (client đọc rồi bỏ) |
| 6 | countCaption | byte | **0** (bảng mốc kinh nghiệm — rỗng an toàn cho vào-map) |
| 7.. | caption.power | long × count | (0 vòng nếu count=0) |

Client đọc count+bảng long vào `GameScr.exps` **sau** khi chạy logic so phiên bản.
count=0 → `exps` rỗng, không sập khi đứng trong map.

---

## 3. `-87` data (client xin, máy chủ trả)

Client xin: `Service.cs:1355` gửi `-87` thân rỗng. Máy chủ nhận: `Controller.java:513`.
Máy chủ trả: `DataGame.java:85-100`. Client nhận: `Controller.cs:4746-4754` (chỉ
**lưu raw vào RMS, không parse ngay** — lỗi part chỉ lộ khi VẼ nhân vật ở Kế hoạch 4).

| # | Trường | Kiểu | Nội dung |
|---|---|---|---|
| 1 | vsData | byte | 20 (khớp mục 2) |
| 2 | len_dart | int (4, BE) | 2 |
| 3 | dart | byte[2] | `00 00` (short 0 = rỗng) |
| 4 | len_arrow | int | 2 |
| 5 | arrow | byte[2] | `00 00` |
| 6 | len_effect | int | 2 |
| 7 | effect | byte[2] | `00 00` |
| 8 | len_image | int | 2 |
| 9 | image | byte[2] | `00 00` |
| 10 | len_part | int | **87707** |
| 11 | part | byte[87707] | **nội dung file `1-MAY-CHU-GAME/data/update_data/part` nguyên bản** |
| 12 | len_skill | int | 2 |
| 13 | skill | byte[2] | `00 00` |

**Tổng ≈ 87.742 byte < 100.000** (mốc `mark(100000)` client — vượt là client
`reset()` lỗi, đứt chuỗi, kẹt "Xin chờ").

**Vì sao part bắt buộc, 5 blob kia rỗng được:** client vẽ nhân vật bằng
`GameScr.parts[head/body/leg]` **không kiểm biên** (`Char.cs:5998-6000`) — part
rỗng hoặc thiếu chỉ số = `IndexOutOfRange` khi vẽ (Kế hoạch 4). Dart/arrow/effect/
image/skill chỉ dùng khi thứ đó thật sự được vẽ (chưa có ở giai đoạn này); icon
thiếu thì client tự xin qua `-67 requestIcon`.

**part là placeholder Dragon Ball:** đây là dữ liệu sprite nhân vật cũ, dùng tạm
để nhân vật hiện được; thay bằng part Naruto khi đường ống mỹ thuật giao. **Tuyệt
đối không tự dựng part type 3 hoặc ≥4** (server ghi 3 frame, client đọc 2 → lệch
stream) — nhưng ở đây ta forward nguyên file cũ nên không dựng gì.

---

## 4. `-28/6` map catalog (client xin, máy chủ trả)

Ghi: `DataGame.java:107-138`. Đọc: `Controller.cs:4703-4742`.

| # | Trường | Kiểu | Giá trị (rỗng) |
|---|---|---|---|
| 1 | lệnh con | byte | 6 |
| 2 | vsMap | byte | 20 |
| 3 | countMap | byte (client `readUnsignedByte`) | **0** |
| 4.. | tên map | UTF × count | (0) |
| | countNpc | byte (client `readByte` CÓ DẤU, ≤127) | **0** |
| | *mỗi npc* | UTF tên + short×3 head/body/leg + byte menuCount | (0) |
| | countMob | byte (có dấu, ≤127) | **0** |
| | *mỗi mob* | byte type + UTF tên + **int hp (4 byte)** + byte rangeMove + byte speed + byte dartType | (0) |

Rỗng hoàn toàn (0/0/0) hợp lệ cho bắt tay, client `createMap` không sập. Khi bản
đồ có npc/mob thật (Kế hoạch sau), countNpc/countMob phải là **mảng dày đặc từ
index 0 tới id lớn nhất** (id template = chỉ số mảng, không bỏ trống).

> **Ràng buộc cứng:** `hp` mob ghi **4 byte int** (không phải long/double). Đúng vì
> client build này có `ModFunc.isReadInt=true` (`ModFunc.cs:45`, đã xác nhận),
> `isReadDouble=false` (`ModFunc.cs:47`). Sai kiểu → lệch stream → kẹt "Xin chờ".

---

## 5. `-28/7` skill catalog — **5 LÀNG NARUTO**

Ghi: `DataGame.java:141-216`. Đọc: `Controller.cs:4648-4701`.

| # | Trường | Kiểu | Giá trị |
|---|---|---|---|
| 1 | lệnh con | byte | 7 |
| 2 | vsSkill | byte | 20 |
| 3 | countOption | byte | 0 (không skill-option) |
| 4.. | (option) | — | (0 vòng) |
| | **countClass** | byte | **5** (5 làng) |
| | *mỗi làng (i=0..4)* | | |
| | tên | UTF | tên làng (Làng Lá/Cát/Sương Mù/Đá/Mây — từ danh mục `lang.go`) |
| | countTemplate | byte | **0** (chưa có skill — Kế hoạch skill sau) |

**Điểm quyết định:** client dùng byte `gender` trong ME_INFO (`-30/0`) làm chỉ số
`GameScr.nClasss[gender]` — nên `countClass` **phải ≥ 5** để id làng 0-4 hợp lệ.
Nếu client còn cache 3 lớp Dragon Ball mà nhân vật ở làng 3/4 → `IndexOutOfRange`.
Đây là lý do phải ghi đè bằng `-28/7` của ta (mục 1: vsSkill ≠ 1).

Trường mô tả skill trong bản Java bị ghi cứng chuỗi `"Ngoc Rong Kame"`
(`DataGame.java:163`) — **dấu vết game cũ, ta không phát lại** (không có skill nên
không tới trường này).

---

## 6. `-28/8` item catalog (client xin, máy chủ trả)

Ghi: `ItemData.java:18-119`. Đọc: `Controller.cs:6453-6536` (`loadItemNew`).

Gửi **4 message `-28` con** phân biệt bằng byte subtype (sau `8` và `vsItem`):
subtype 0 (option template), 1 (item 0..750), 2 (item thêm), 100 (arrHead2Frames).

**Tối thiểu tuyệt đối = đúng MỘT message subtype 2 rỗng:**

| # | Trường | Kiểu | Giá trị |
|---|---|---|---|
| 1 | lệnh con | byte | 8 |
| 2 | vsItem | byte | 20 |
| 3 | subtype | byte | 2 |
| 4 | start | short | 0 |
| 5 | end | short | 0 |

Chỉ subtype 2 mới gọi `clientOk` (`Controller.cs:6517-6534`), và chỉ khi đủ 4
phiên bản khớp. **Thiếu subtype 2 = kẹt "Xin chờ" vĩnh viễn** (soft-lock, không
sập). Nhân vật mới không có đồ (`itemsBody/Bag/Box size=0` trong ME_INFO) vẫn vẽ
được với catalog item rỗng (part/thân/chân lấy từ `-87` + `setDefaultPart`, không
qua ItemTemplate).

> Ghi chú: bỏ subtype 0/1/100 khiến cache item không tái dùng được ở lần đăng
> nhập sau (client tự xin lại mỗi lần — tự phục hồi, không treo). Chấp nhận được
> ở giai đoạn này.

---

## 7. `-82` tileSetInfo (Kế hoạch 4)

Ghi: `DataGame.java:424-434` — **gửi nguyên file `1-MAY-CHU-GAME/data/map/tile_set_info`
(1696 byte), không header, không biến đổi.** Đọc: `Controller.cs:843-865`.

Định dạng file (đã parse xác nhận 1696/1696 byte, 39 tileset):
```
byte numTileMap                          (≤127, client đọc sbyte)
  × numTileMap:
    byte numTileOfMap                    (≤127)
      × numTileOfMap:
        int  tileType   (4 byte BE)
        byte numIndex                    (≤127)
          × numIndex: byte typeIndex
```

**Bẫy chí mạng:** thiếu `-82` (hoặc gửi sau `-24`, hoặc `map.tileId` ngoài
`[1, numTileMap]`) → client `GameMidlet.instance.exit()` — **thoát hẳn ứng dụng**
(`TileMap.cs:255-259`). Gửi 1 lần/phiên, trước `-24`. Với làng Naruto dùng tileId
có sẵn trong file này thì forward nguyên file là đủ.

---

## 8. `-30/0` ME_INFO (Kế hoạch 4)

Ghi: `Service.java:776-897` (+ `DataGame.sendHeadAvatar:455-464`). Đọc:
`Controller.cs:5820-6009`.

| # | Trường | Kiểu | Nguồn (nhân vật) |
|---|---|---|---|
| 1 | id | int | `NhanVat.ID` |
| 2 | taskMain.id | byte | 0 (chưa có nhiệm vụ) |
| 3 | gender | byte | **`NhanVat.Lang`** (0-4) |
| 4 | head | short | `NhanVat.KieuToc` |
| 5 | name | UTF | `NhanVat.Ten` |
| 6 | cPK | byte | 0 |
| 7 | typePk | byte | 0 |
| 8 | power | long | `NhanVat.SucManh` |
| 9 | buffHp | short | 0 |
| 10 | buffMp | short | 0 |
| 11 | **gender (lần 2)** | byte | **`NhanVat.Lang`** — chỉ số `nClasss[]`, PHẢI < countClass (=5) |
| 12 | skillCount | byte | 0 |
| 12a | (skillId) | short × count | (0) |
| 13 | gold | **long (8 byte!)** | `NhanVat.Vang` |
| 14 | ruby | int | 0 |
| 15 | gem | int | `NhanVat.Ngoc` |
| 16 | itemsBody.size | **unsigned byte** | **0** |
| 17 | itemsBag.size | byte có dấu (≤127) | **0** |
| 18 | itemsBox.size | byte có dấu (≤127) | **0** |
| 19 | headAvatar.size | short | 0 |
| 20 | conMeo[gender][0..2] | short × 3 | 0,0,0 |
| 21 | fusion | byte | 0 |
| 22 | deltaTime | int | 0 |
| 23 | isNewMember | byte | 0/1 |
| 24 | aura | short | (bỏ được — client bắt lỗi êm) |
| 25 | effFront | byte | (bỏ được) |
| 26 | hat | short | (bỏ được) |

**Bẫy:** (a) **gold phải 8 byte long** — client luôn `readLong()`; ghi 4 byte lệch
toàn bộ phần sau. (b) **gender (byte 11) phải < countClass** — nếu không
`IndexOutOfRange` bị nuốt, nhân vật kẹt "Xin chờ" (không set được `statusMe=4`).
(c) item id lạ (không có trong catalog `-28/8`) → NPE giữa gói → gửi cả 3 mảng
item size=0 cho nhân vật mới. (d) 3 trường cuối (24-26) bỏ hẳn được.

---

## 9. `-24` MAP INFO (Kế hoạch 4)

Ghi: `Zone.java:738-869`. Đọc: `Controller.cs:5027-5346` → `switchToMe()` (rời
"Xin chờ").

| # | Trường | Kiểu | Giá trị (map trống) |
|---|---|---|---|
| 1 | mapId | **unsigned byte (≤255)** | id map làng |
| 2 | planetId | byte | (của map) |
| 3 | tileId | byte | **∈ [1, numTileMap]** của `-82` |
| 4 | bgId | byte | (của map) |
| 5 | type | byte | 0 (MAP_NORMAL) |
| 6 | mapName | UTF | tên làng |
| 7 | zoneId | byte | 0 (< số zone của map) |
| 8 | x | short | vị trí spawn X |
| 9 | y | short | vị trí spawn Y |
| 10 | waypointCount | byte (≤127) | **0** |
| 11 | mobCount | byte (≤127) | **0** |
| 12 | **byte dự phòng** | byte | **0** ← ĐỪNG QUÊN, nằm giữa mob và npc |
| 13 | npcCount | byte (≤127) | **0** |
| 14 | itemCount | byte (≤127) | **0** |
| 15 | bg-item | short count + N×6 byte | **`00 00`** (short 0 — KHÔNG bỏ 0 byte) |
| 16 | effMapCount | short | **0** |
| 17 | bgType | byte | (của map) |
| 18 | idSpaceShip | byte | 0 |
| 19 | isMapDouble | byte | 0 (trừ mapId=148) |

**Bẫy:** (a) **quên byte dự phòng #12** giữa mob và npc → lệch toàn bộ phần sau.
(b) mọi count ≥128 → client đọc âm → lệch. (c) mapId ≤255, tempId mob/npc ≤127.
(d) **`playerId=-2` là sentinel** (client đọc thêm 1 short) — không bao giờ gửi
`-2` cho item. (e) mob.hp = **int 4 byte** (như `-28/6`). (f) bg-item rỗng vẫn
phải `short 0` (2 byte), không bỏ trắng. Client KHÔNG sập với 0 mob/npc/waypoint/item.

---

## 10. Định dạng file tile map (Kế hoạch 4, cho kiểm địa hình)

File `1-MAY-CHU-GAME/data/map/tile_map_data/<mapId>` (đọc: `Manager.java:1116-1132`,
xác nhận bằng dump byte thật):
```
byte tmw   (unsigned, số ô ngang)
byte tmh   (unsigned, số ô dọc)
tmw×tmh byte, row-major: mỗi byte = chỉ số tile
```
`mapWidth = tmw×24`, `mapHeight = tmh×24` (SIZE=24 px/ô). Ô giá trị **0 = trống**
(bay lơ lửng), khác 0 = có nền. Máy chủ Java **không kiểm va chạm tường** khi di
chuyển (chỉ dùng để bật cờ `isFly`). Kiểm địa hình là việc thiết kế mới ở Kế hoạch
di chuyển (xem `2026-07-23-hoa-giai-muot-va-chong-gian-lan.md` §4.3).

Gói `-28/10` (map template) chỉ cần khi client thiếu tile map của một mapId trong
resource local `/mymap/`. Client build hiện tại có sẵn mọi map cũ; **map Naruto
dùng mapId cũ thì không cần `-28/10`.** Nếu thêm mapId mới ngoài client build,
**bắt buộc** trả `-28/10` (`DataGame.java:437-452`: byte tmw + byte tmh + tile
data nguyên), nếu không client kẹt "Xin chờ" vĩnh viễn dù `-24` đúng.

---

## 11. Ràng buộc cứng (chốt vào Global Constraints của kế hoạch)

1. Client build đích: `D:\NRO_Project`, có `ModFunc.isReadInt=true` → mọi `hp`
   (mob, `-28/6`, `-24`) ghi **4 byte int**.
2. Mọi `count`/`size`/`tempId` đọc byte **có dấu** → **≤127** (trừ `mapId` và
   `itemsBody.size` đọc unsigned → ≤255).
3. Gói `-87` **< 100.000 byte**.
4. `gold` trong ME_INFO là **long 8 byte**.
5. `-82` phải gửi **trước** `-24`; `map.tileId ∈ [1, numTileMap]`.
6. `countClass` trong `-28/7` **≥ 5**; `gender`-slot trong ME_INFO là **id làng 0-4**.
7. Chuỗi UTF là **modified UTF-8 Java** — với tiếng Việt (trong BMP, không NUL)
   giống UTF-8 chuẩn Go; dùng lại `protocol.Writer.UTF` đã có.
8. Phiên bản Naruto `vsData=vsMap=vsSkill=vsItem=20` (≠ cũ 10/2/1/4).
