# Thiết kế Phase 1 — Vertical slice Làng Lá có hình Naruto và E2E Unity

**Ngày:** 29/07/2026
**Trạng thái:** Chốt để lập kế hoạch thi công
**Phạm vi:** `server/`, công cụ trong `docs/my-thuat/`, bản staging của
`D:\NRO_Project`; không sửa các thư mục Java/PHP cũ.

## 1. Vì sao làm mốc này trước

Máy chủ Go đã có đăng nhập, vào bản đồ, di chuyển, chiến đấu, chết–hồi sinh và
lưu trạng thái. Các test Go và PostgreSQL thật đang xanh. Tuy nhiên client Unity
chưa thể nghiệm thu gameplay bằng mắt vì:

- `part.bin` chính thức vẫn là dữ liệu Dragon Ball;
- map 1–5 chỉ là bản sao map rỗng;
- background và cache `Big*.png` vẫn là Dragon Ball;
- POC hiện tại chỉ đổi part 30/57/58 sang ba ảnh đứng yên `9000..9002`, chép tay
  vào Unity;
- test chỉ ghim kích thước 87.707 byte nên một `part.bin` sai nội dung vẫn xanh;
- một client không đọc mạng có thể khiến actor khu vực chờ `Session.Gui()` tới
  30 giây nhiều lần.

Phase 1 tạo một lát cắt Làng Lá nhỏ nhưng chơi thật. Đây là cổng kỹ thuật trước
khi nhân rộng sang năm làng, không phải tuyên bố toàn bộ game đã hoàn thành.

## 2. Tiêu chí người chơi nhìn thấy

Một lượt nghiệm thu phải đi được trọn chuỗi sau mà không thấy tài nguyên Dragon
Ball:

1. Hai tài khoản mở hai client, vào cùng khu vực Làng Lá và thấy nhau.
2. Nhân vật Học viên/Genin có hình và hoạt ảnh đứng, chạy, đánh, trúng đòn,
   chết và hồi sinh.
3. Có một NPC Iruka đứng yên, chân dung và một đoạn hội thoại thật qua lệnh mở
   NPC của client.
4. Có một Sói Rừng với hoạt ảnh đứng/chạy/vồ/trúng đòn/chết; AI máy chủ đánh,
   chết và hồi sinh đúng. Phase 1 chưa biến vị trí quái thành chuyển động
   authoritative trên server.
5. Bốn chiêu Làng Lá có icon Naruto và tư thế ra đòn đúng; thanh HP/chakra dùng
   skin Naruto và mức fill giảm/tăng theo giá trị thật. Projectile/effect catalog
   đầy đủ được làm ở KH8, không dùng hình Dragon Ball để giả là đã xong.
6. Map có tile Naruto nhìn được, nền đứng được, vật cản/va chạm client–server
   khớp nhau và bốn lớp
   parallax Làng Lá.
7. Hai người thấy cùng diễn biến đánh quái và PvP; chết–hồi sinh đồng bộ.
8. Tắt máy chủ đúng quy trình rồi bật lại giữ chính xác map, khu, tọa độ và HP.
9. Một client chậm hoặc ngừng đọc không làm đứng actor và người chơi còn lại.
10. Hai lượt cache sạch và cache cũ đều vào được game bằng đúng cặp
    server/client staging; không mượn Big/NR cache Dragon Ball từ máy.

## 3. Ba phương án tài nguyên và lựa chọn

### 3.1 Chỉ dùng `Small<ID>.png`

Đơn giản, dễ soi từng frame và POC đã chứng minh client render được. Nhược điểm
là nhiều file, build phình và không phải đích phát hành lâu dài.

### 3.2 Đóng ngay toàn bộ `Big0..4 + NR_image`

Đúng kiến trúc legacy và hiệu quả khi có nhiều sprite. Nhưng Go chưa phân phối
Big resource hoàn chỉnh; `NR_image` và Big phải đổi nguyên tử, trong khi client
chỉ có năm atlas và giới hạn vùng crop. Làm tất cả ngay sẽ trộn vertical slice
với một dự án resource CDN/cache lớn.

### 3.3 Hybrid — lựa chọn

ID sprite là căn cước bền. Phase 1 phát hành lát cắt bằng `Small<ID>.png` rời và
`NR_image` không có mapping cho các ID đó. Client ưu tiên atlas **khi** mapping
hợp lệ; Small là fallback khi mapping thiếu/lỗi. Sau này cùng ID được chuyển
sang atlas mà không đổi `part.bin`, skill hay gameplay. Khi chuyển, bản Small
tương ứng bị loại khỏi build phát hành để không che lỗi atlas bằng fallback cũ.

Lựa chọn này không giữ cách chép tay của POC: mọi file vẫn phải được sinh và
deploy từ manifest, có lock/hash và kiểm tra tự động.

## 4. Một nguồn sự thật: `mapping-id.json`

Tạo `docs/my-thuat/mapping-id.json` với các nguyên tắc:

- `schema_version`, `source_catalog_version`, `runtime_slice_version`,
  `data_version`, `map_version` và `resource_version` rõ ràng;
- `sprite_id` và `part_id` ghi tường minh, không suy ra từ sort tên file;
- ID mới chỉ nối cuối; xóa dùng `tombstone`, không tái sử dụng;
- đổi ảnh nhưng giữ cùng ý nghĩa thì giữ ID và tăng `asset_version`;
- đổi ý nghĩa thì cấp ID mới;
- mỗi sprite ghi nguồn x4, SHA-256, kích thước, alpha, anchor, cách phân phối
  `small|atlas`; sprite tách từ ảnh phẳng còn phải ghi `source_rect`, overlap,
  thứ tự lớp và crop policy;
- mỗi part ghi đủ sequence frame, `sprite_id`, `dx`, `dy`;
- mỗi slice làng ghi default head/body/leg, map, tile, bg và danh sách tài nguyên;
- frame lặp phải khai báo có chủ ý; không được trỏ toàn bộ 17/14 frame vào một
  ảnh đứng yên như POC.

Giới hạn phải được validator chặn:

- không suy số lượng namespace từ 2.997 record `NR_image`: `part.bin` legacy
  thực tế tham chiếu 13.317 sprite ID khác nhau, tới ID 32704. Scanner baseline
  ghim toàn bộ closure local (kể cả các cache `NR_part`) và dành riêng khoảng
  trống `27041..27999` cho Naruto; Phase 1 chỉ dùng `27041..27082`, các ID
  Naruto sau đó nối từ 27083. Đây là vùng an toàn theo input local đã khóa,
  không phải lời khẳng định về CDN legacy chưa có;
  ID thử nghiệm `9000..9002` của POC không được mang vào artifact chính thức;
- validator phải quét `part.bin`, `NR_image`, skill/mob data, resource Small và
  ID hard-code phía client; bất kỳ legacy reference/resource nào đụng khoảng
  dành riêng đều làm build đỏ trước khi ghi artifact;
- giữ nguyên 2.189 record part legacy; Phase 1 chỉ thay part 30/57/58 và nối thêm
  ba part Iruka ở 2189/2190/2191. Mọi record khác phải byte-identical;
- `sprite_id`, `part_id` nằm trong `0..32767`; `part_count > max(part_id)` và
  mọi khoảng trống đều có tombstone tường minh;
- `dx`, `dy` nằm trong miền `sbyte`;
- part type 0/1/2 lần lượt có đúng 3/17/14 frame;
- mỗi frame ghi cả `slot_index` và tập `client_cf` lấy từ bảng
  `Char.CharInfo[0..32]`; đảo slot phải làm golden state đỏ;
- atlas chỉ 0–4; `bigIndex` là byte không âm; `x/y/w/h` không âm, nằm trong PNG
  và nhỏ hơn 256 theo fast-path thực của client;
- mọi sprite runtime tồn tại đủ x1, x2, x3, x4.

`source_catalog_version` theo dõi cả 25 ảnh nguồn, kể cả ảnh chưa deploy.
`runtime_slice_version`/lock chỉ đổi khi closure tài nguyên Làng Lá thay đổi.
Validator dựng closure từ part, NPC, mob, skill icon/pose, UI, tile và background;
mọi ID/path trong closure phải có file local. Phase 1 không được phụ thuộc
`Big*.png`, `NR_image`, `NR_effect`, `NR_arrow` hay `NR_dart` Dragon Ball.

Allocation Phase 1 được ghim:

| ID | Nghĩa |
|---|---|
| 27041..27074 | 3 head + 17 body + 14 leg của Học viên |
| 27075..27077 | head/body/leg Iruka |
| 27078 | chân dung Iruka |
| 27079..27082 | bốn icon skill Làng Lá |

Không chèn ID mới vào giữa các khoảng này. Asset ngoài slice bắt đầu từ 27083
theo thứ tự append-only sau khi Phase 1 khóa lock.

## 5. Công cụ build tài nguyên

Tạo package Python ở `assetpack/` tại gốc repo với hai lệnh chính:

```text
python -m assetpack validate --manifest docs/my-thuat/mapping-id.json
python -m assetpack build --manifest docs/my-thuat/mapping-id.json \
  --slice lang-la --out <thu-muc-build>
```

Đầu vào:

- ảnh nguồn x4 đã duyệt;
- `part.bin` legacy có SHA-256 ghim độc lập;
- `tile_set_info` legacy có SHA-256 ghim độc lập;
- layout map, bảng `tile_index → ngữ nghĩa/solid` và tile art Làng Lá;
- bốn lớp background;
- dữ liệu frame Sói Rừng theo khuôn `EffectData.readData`;
- catalog map/NPC/mob và bốn pose skill;
- manifest.

Đầu ra xác định, build hai lần phải cùng SHA:

```text
<out>/server/part.bin
<out>/server/tile_set_info.bin
<out>/server/map-1-collision.bin
<out>/server/skill.bin
<out>/server/catalog.json
<out>/client/x1..x4/smallimage/Small<ID>.png
<out>/client/mymap/1.bytes
<out>/client/x1/t/1/t_01.png..t_100.png
<out>/client/x2..x4/t/1$1.png..1$100.png
<out>/client/x1..x4/bg/b00.png..b03.png
<out>/client/Mob/0.bytes
<out>/client/Mob/0/data.bytes
<out>/client/x1..x4/Mob/0/img.png
<out>/client/x1..x4/mainImage/naruto-hp-bg.png
<out>/client/x1..x4/mainImage/naruto-hp-fill.png
<out>/client/x1..x4/mainImage/naruto-chakra-bg.png
<out>/client/x1..x4/mainImage/naruto-chakra-fill.png
<out>/client/x1..x4/mainImage/<8-ui-runtime-overrides>.png
<out>/asset-lock.json
```

`server/catalog.json` là artifact xác định cho sáu map, template Iruka, template
Sói Rừng, vị trí NPC/mob và bốn pose skill; code Go nạp chính file này thay vì
chép lại hằng số từ manifest. `asset-lock.json` ghi cả hash catalog, version,
hash đầu vào/đầu ra, kích thước PNG, ID, part,
frame count, runtime closure và SHA của server/client ghép cặp. Script deploy
chỉ chép từ `<out>` vào một Unity staging snapshot theo bảng:

| Artifact | Đích trong staging |
|---|---|
| `client/xN/**` | `Assets/Resources/res/xN/**` |
| `client/mymap/*.bytes` | `Assets/Resources/res/mymap/*.bytes` |
| `client/Mob/**` | `Assets/Resources/res/Mob/**` |

Không có lệnh nào được sửa ảnh nguồn hoặc Unity gốc tại chỗ.

Server nhúng `part.bin` sinh ra và quảng bá `vsData` khớp lock. Test phải đỏ nếu
asset runtime đổi nhưng version không tăng; oracle so manifest/lock hiện tại với
ledger version đã commit, không tự tính expected bằng code đang test. `vsMap`
tăng khi catalog/map/tile/collision đổi. `skill.bin` là blob `NR_skill` tối thiểu
cho 12 levelId Làng Lá, chỉ điều khiển pose; các effect/arrow ID đều `-1`.
Unity overlay giữ `-1 = không có effect` cho data version Naruto, thay vì ép
mọi `effectHappenOnMob<=0` thành effect 80 Dragon Ball. Golden/log cấm load
effect 80 trong lát cắt.

Go `-28/6` phải phát catalog có sáu tên map (index 0..5), một NPC template Iruka
và một mob template Sói Rừng. `-24` phải phát NPC Iruka thật với template/part
đã sinh và quái dùng `MauId=0`. Lệnh client `33` mở Iruka được server trả lệnh
`32` gồm câu thoại và menu đóng. Không có catalog rỗng trong lát cắt.

Phase 1 ghim nhân vật E2E Làng Lá dùng `KieuToc=30`; client gender/làng 0 hiện
dùng body 57 và leg 58. Fixture phải fail-fast nếu ME_INFO không trả bộ
30/57/58. Luồng chọn ba kiểu tóc lúc tạo nhân vật và default part của bốn làng
còn lại chưa được tuyên bố hoàn thành; bước mở rộng năm làng sẽ sinh bảng default
part thay hard-code client.

Tám UI đã duyệt (`ui_logo_game`, khung avatar/hội thoại, ba nút, ô chọn và ô vật
phẩm) được manifest ánh xạ tới **consumer thật** của Game1, không giả định mọi
consumer đều đã có một resource key 1:1. Mỗi binding có một mode:

- `direct`: một hoặc nhiều logical key thật được call-site nạp trực tiếp;
- `slice_group`: consumer ghép nhiều key/slice của một composite;
- `overlay_required`: client hiện vẽ primitive/composite không có key thay thế;
  overlay phải thêm `future_logical_key`, có `consumer_sites[]` không rỗng và
  consumer SHA làm precondition.

Task phân tích đầu tiên phải ghim đủ tám binding, `logical_keys[]`, loader và
physical paths bằng test source/load. Binding `direct` và `slice_group` phải có
call-site evidence/SHA, nhưng được phép có `consumer_sites[]` rỗng. Binding
`overlay_required` phải có `consumer_sites[]` không rỗng, `future_logical_key`
và consumer SHA precondition. Không được đổi tên file theo phỏng đoán hoặc bịa
key không tồn tại. Golden từ splash tới gameplay không được còn logo/icon
Dragon Ball trong các checkpoint của lát cắt.

Phase 1 chỉ chạy bằng cặp staging có lock khớp; runner/preflight từ chối khởi
động nếu lock server/client khác nhau. Đây chưa phải cơ chế cập nhật client công cộng:
không phát hành server này cho client cũ trước khi KH resource distribution/gate
được làm. Phase 1 không phụ thuộc server tải Big qua `-74`.

Artifact build nằm trong thư mục theo SHA và không ghi đè. Giữ ít nhất artifact
đang chạy và một artifact trước đó. Rollback luôn khôi phục **cả cặp**
server/client lock; không hạ số version hoặc ghép `part.bin` mới với client cũ.
Fixture E2E chỉ dùng CSDL test và prefix run ID, không chạm dữ liệu tài khoản thật.

## 6. Ảnh và hoạt ảnh

### 6.1 Dùng lại

Từ bộ 343 ảnh dùng cho lát cắt:

- `map_la_*`, `map_rung_*`, `map_nen_troi_*`, `map_nen_may_*`,
  `map_nen_nuihokage`;
- `npc_iruka`, `cd_iruka`;
- `quai_soi_rung`;
- một bộ `tp_hocvien_*` hoặc `tp_genin_co_*`;
- hiệu ứng hồi phục, tụ chakra, hỏa độn, cận chiến;
- UI thực sự xuất hiện trong lát cắt.

Asset của bốn làng còn lại, boss, NPC, biome và phần lớn vật phẩm chưa đi vào
runtime Phase 1. Chúng vẫn nằm trong kho nguồn và sẽ vào các slice sau.

### 6.2 Tạo thêm bằng Codex image generation

Sinh cùng một đợt 25 ảnh để giữ phong cách:

- 20 icon kỹ năng;
- bốn lớp riêng: nền/fill HP và nền/fill chakra;
- vẽ lại `hu_anthan`.

Phase 1 tích hợp ngay bốn icon Làng Lá và bốn lớp thanh. 17 ảnh còn lại được
duyệt và khóa nguồn cho các phase sau.

Ngoài 25 ảnh catalog, cần bốn sheet runtime mới, ghép xác định từ các action
strip sinh từ cùng master reference đã khóa:

| File nguồn x4 | Khuôn | Nội dung |
|---|---|---|
| `nv_hocvien_head_sheet_x4.png` | 3 ô 128×128 | head slot 0..2 |
| `nv_hocvien_body_sheet_x4.png` | 17 ô 128×160 | body slot 0..16 |
| `nv_hocvien_leg_sheet_x4.png` | 14 ô 128×128 | leg slot 0..13 |
| `quai_soi_rung_sheet_x4.png` | 10 ô 96×96 | idle 2, chạy 4, vồ 2, trúng 1, chết 1 |

Imagegen không phải tự tạo một hàng 17/14/10 ô hoàn hảo trong một lần. Mỗi lần
sinh một action strip hoặc một nhóm pose nhỏ từ **cùng master reference**, giữ
chung trang phục, tỷ lệ và bảng màu; assetpack chuẩn hóa canvas rồi ghép theo
`slot_index` thành bốn sheet runtime ở trên. Validator phải kiểm đúng số ô,
kích thước từng ô, thứ tự slot và không có pixel nhìn thấy vượt ranh giới ô.
Không được sinh từng frame độc lập bằng các prompt không chung reference.

Ba ảnh `tp_hocvien_*` hiện có là reference phong cách/tỷ lệ, không phải 34 frame
runtime. `quai_soi_rung` hiện có là reference silhouette. Iruka được assetpack
cắt đầu/thân/chân từ ảnh tĩnh; manifest ghim ba `source_rect`, vùng overlap ở
cổ/hông, anchor bàn chân và thứ tự vẽ chân→thân→đầu. Golden composite Iruka hoàn
chỉnh phải xanh trên cả nền sáng và tối, không khe/mất nét ở biên crop, trước khi
chấp nhận part 2189..2191. Lặp frame đứng yên là hành vi NPC có chủ ý.

Image generation chỉ tạo ảnh nguồn lớn trên nền chroma-key phẳng. Công cụ
imagegen chuẩn của Codex tách nền bằng helper rồi phải qua QA alpha; pipeline
không được coi tách nền tự động là luôn đúng. Assetpack chỉ nhận ảnh đã duyệt,
cắt/ghép frame, thu nhỏ bốn zoom, đặt anchor và đóng gói.

Sequence player phải có silhouette khác cho idle/thở, chạy, đánh cận chiến,
thủ ấn/phát chiêu, trúng đòn, chết và đứng dậy. Manifest ghi `slot_index` cùng
mọi `cf` sử dụng slot đó; golden render trực tiếp các `cf` đại diện từ
`Char.CharInfo`, không tự đặt tên rồi suy đoán thứ tự.

QA bắt buộc:

- xem x1 trên nền sáng và tối; không halo xanh/trắng/đen;
- alpha có vùng nhìn thấy và góc trong suốt;
- anchor bàn chân của idle/run lệch không quá 1 px ở x1;
- chiều cao/tỷ lệ trang phục không trôi giữa frame cùng action;
- silhouette đòn đánh đọc được ở x1 và pose impact trùng checkpoint gây damage;
- capture video/GIF/contact strip toàn bộ frame ở nhịp runtime cố định cho từng
  action player và Sói, gồm cả chuyển frame cuối→đầu; reviewer phải xác nhận
  không đảo/lặp frame ngoài thiết kế, không giật anchor và nhịp chạy/đánh/impact
  khớp event gameplay;
- contact sheet/golden có reviewer và PASS/FAIL được ghi vào hồ sơ E2E.

### 6.3 Tích hợp thanh HP/chakra

Unity staging patch thay đúng đường vẽ HUD trong `GameScr`:

- load bốn key `/mainImage/naruto-hp-bg`, `naruto-hp-fill`,
  `naruto-chakra-bg`, `naruto-chakra-fill`;
- nền và fill của từng thanh cùng kích thước; fill được clip từ trái sang phải
  theo `clamp(current/max, 0..1)`, không co giãn texture;
- `0%` không vẽ fill, `100%` vẽ đủ chiều rộng; chia cho `max<=0` trả 0;
- HP và chakra đọc trực tiếp state client vừa nhận, không dùng biến demo;
- golden chụp 25%, 50%, 100% và telemetry phải có đúng tỷ lệ tương ứng.

Các thay đổi client, profile và telemetry được lưu thành overlay/patch có
precondition SHA trong repo, rồi áp vào snapshot; không chỉnh tay
`D:\NRO_Project`. Iruka dùng part 2189..2191 và avatar Small ID từ manifest để
chat popup thực sự hiển thị chân dung `cd_iruka`.

## 7. Map và background

`mymap/1.bytes` phải do build tạo với khuôn `width byte`, `height byte`, rồi
`width*height` tile. Validator kiểm:

- độ dài đúng `2 + width*height`;
- map dùng `tileID=1`, `bgID=0`, `bgType=0`; server `BanDo` và lock phải khớp;
- mỗi tile index nằm trong `0..99`;
- spawn `(100,384)` đứng trên nền, không nằm trong vật cản;
- đường chơi đủ cho hai người, NPC và quái;
- collision khớp hình nhìn thấy.

Phase 1 giữ bảng collision/type của tileID 1 trong `tile_set_info` legacy, nhưng
thay **hình** của đúng 100 index bằng tile Naruto tại
layout native của `Main.typeClient=4`: zoom x1 ở
`res/x1/t/1/t_01.png..t_100.png`, còn zoom x2..x4 ở
`res/xN/t/1$1.png..1$100.png`. Đây là đúng hai nhánh `TileMap.getTile()` hiện
tại; dấu `$`, prefix `t_` và padding x1 đều là một phần của contract. Nếu đổi
loader phải có overlay kèm precondition SHA và test resource-load ở cả bốn zoom.
Manifest ghi ngữ nghĩa từng index được dùng
(`ground|wall|decor|empty`) và `solid`. Packer phải chứng minh `tile_set_info`
đầu ra byte-identical với baseline trừ khi một task riêng có test client parser
cho thay đổi đó.

Mỗi tile có kích thước vật lý x1/x2/x3/x4 lần lượt 24×24, 48×48, 72×72 và
96×96 px, cùng crop logical 24×24; scale bằng nearest-neighbor, cấm blur/bilinear
và cấm pixel nhìn thấy tràn mép. Unity batch/EditMode preflight bắt buộc đặt
`Main.typeClient=4`, thử load đủ 100 resource key ở từng zoom, kiểm non-null và
đúng kích thước. Visual gate sinh contact sheet 10×10 và render toàn map ở cả
bốn zoom để bắt tile mờ, lệch cạnh và đường nối; canary `ASSET-ZOOM-003` phải
đỏ qua đúng preflight này.

Cùng bảng `solid` sinh `server/map-1-collision.bin`. Máy chủ nạp artifact này để
kiểm điểm chân và đoạn di chuyển, thay cho biên 4096 đơn thuần; client và server
phải cho cùng kết quả trên golden walkable/blocked. Canary đổi một ô client mà
không đổi collision server (và ngược lại) phải đỏ vì hash/semantic mismatch.

Background giữ renderer parallax hiện có, dựng bốn lớp `b00..b03` ở đủ bốn
zoom. `bgID=0` khiến renderer nạp theo thứ tự xa→gần `b00`, `b01`, `b02`,
`b03`; tốc độ/offset tiếp tục theo `GameCanvas.paintBackgroundtLayer` cho type
0. Mỗi lớp ghi anchor/repeat trong manifest, deploy vào `res/xN/bg/`, và không
được để cache/RMS Dragon Ball thắng resource staging. Không đổi một ảnh map đơn
thành `b1.png` như thử nghiệm cũ. Test kiểm đúng resource key, đủ lớp, đúng tỷ
lệ và không có khe rõ ở biên lặp.

## 8. Không để client chậm chặn actor

Actor khu vực vẫn là nơi duy nhất thay đổi trạng thái gameplay, nhưng không
được gọi thao tác ghi socket có thể chặn.

Thiết kế:

1. Tách interface gameplay thành `XepGui(*protocol.Message) error` và
   `DongCham()`. `XepGui` chỉ enqueue, tuyệt đối không ghi socket.
2. Mỗi session có đúng một writer goroutine và hàng FIFO giới hạn 256 message.
   Message đã enqueue là immutable; không tạo goroutine theo broadcast và không
   có hai writer ghi xen byte.
3. Actor dựng message rồi gọi `XepGui`. Hàng đầy trả `ErrHangGuiDay` ngay; actor
   gọi `DongCham()` idempotent, chỉ đóng connection/đánh dấu, không chờ writer.
4. Connection đóng làm vòng phục vụ chạy `KhiDong()` đúng một lần. Thứ tự bắt
   buộc: xin snapshot có timeout → xếp lưu → gửi control event rời khu vực.
5. `KhuVuc` có kênh control riêng cho snapshot/rời/dừng. Control event không đi
   qua `Gui()` best-effort và không được rơi khi mailbox gameplay đầy; API trả
   lỗi rõ khi actor đã dừng. Vòng actor ưu tiên drain control trước gameplay.
6. `suKienRa` xóa người chơi khỏi actor trước rồi phát `-6`. Test phải chứng
   minh không còn “người chơi ma” và snapshot cuối đã được xếp lưu dù outbound
   queue của chính phiên đó đầy.
7. Một client chậm không làm mất event của client khỏe và actor phải xử lý một
   event probe trong 100 ms ở test deterministic; không dùng timeout 30 giây
   làm oracle.
8. Nếu writer gặp lỗi `Write`/reset socket, nó gọi `DongCham()` idempotent,
   ngừng nhận message mới, bỏ phần queue chưa gửi, thoát goroutine và báo hoàn
   tất để vòng phục vụ join có timeout. Cleanup/`KhiDong()` vẫn chỉ chạy đúng
   một lần; test kiểm không rò goroutine và không còn session đăng ký sau lỗi
   ghi thật.

Các gói login trước khi người chơi vào actor tiếp tục dùng gửi đồng bộ có hạn
giờ. Thay đổi tập trung vào `phienKhuVuc`/broadcast gameplay và control path
disconnect; mọi call site `phien.Gui()` bên trong actor phải bị test/scan cấm.

## 9. E2E Unity mức 2

Runner duy nhất:

```powershell
powershell -File tools/e2e/run-unity-e2e.ps1 `
  -UnityProject D:\NRO_Project -ServerRoot <worktree> -ArtRoot <art-root>
```

Exit `0` khi mọi assertion PASS; exit khác `0` khi preflight/build/fixture/
gameplay/evidence fail. Không dùng hai cửa sổ chung `persistentDataPath`.

### 9.1 Snapshot, build và profile

Mỗi run:

1. Khóa nguồn trong lúc copy, tạo snapshot riêng chỉ gồm `Assets`, `Packages`,
   `ProjectSettings`; loại `Library`, `Temp`, `Logs`, `UserSettings`.
2. Ghi SHA-256 từng file, Unity `2022.3.62f2`, server commit/status, DB schema,
   source/runtime lock. Deploy chỉ vào snapshot có run ID.
3. Build một executable `NRO-E2E` với scripting define `NRO_E2E`; release build
   không chứa telemetry.
4. Mở cùng executable hai lần với:

```text
-run-id <id> -profile-root <absolute A|B> -telemetry-port <port>
-telemetry-token <random-32-byte-token> -logFile <client-A|B.log>
```

`Rms.GetiPhoneDocumentsPath()` dùng `-profile-root` khi và chỉ khi có
`NRO_E2E`, tạo directory trước lần đọc RMS đầu tiên. Canary ghi `ONLY_A` và
`ONLY_B`; mỗi client phải đọc được sentinel của mình và không đọc được sentinel
của client kia.

Chạy hai biến thể cache:

- **clean:** profile không có file `NR*`/Big;
- **warm upgrade:** profile được seed `vsData/vsMap` cũ và cache cũ có hash biết
  trước; server buộc tải artifact mới.

Clean chạy toàn bộ chuỗi gameplay + graceful restart ở §9.4. Warm-upgrade là
smoke riêng: kiểm version negotiation phát hiện cache cũ, tải đúng lock/hash
mới, không dùng Big/NR Dragon Ball, login và vào gameplay/map 1 thành công,
resource player/tile/bg/HUD/skill đều là bản mới. Warm không được báo như một
full gameplay/restart run. `aggregate-result.json` chỉ xanh khi full clean và
warm smoke cùng xanh trên đúng một executable/server lock.

### 9.2 Telemetry contract

Mỗi client bind TCP loopback duy nhất, yêu cầu token ở message đầu. Mọi command
được xếp và thực thi trên Unity main thread. JSON-lines request/response có
`command_id`; event có `monotonic_ms`, `utc`, `client`, `sequence`, `type`,
`payload`.

Command tối thiểu:

```text
login, wait_ready, move_to, select_skill, target_mob, target_player,
attack, respawn, open_npc, choose_npc_menu, capture, get_state, shutdown
```

`get_state` trả:

```text
account, character_id, map_id, zone_id, x, y, hp, hp_max,
chakra, chakra_max, alive, visible_player_ids, visible_mob_ids,
selected_skill, selected_skill_ready_at_ms, target_id,
background_layers[{resource,offset_x,offset_y}],
last_event, sequence
```

`selected_skill_ready_at_ms` dùng cùng clock monotonic của event. Event tối thiểu
cho đường damage xác định là
`damage_received {source_kind,source_id,amount,hp_before,hp_after}` và
`skill_ready {skill_id,ready_at_ms}`; không suy damage từ screenshot.

Ảnh `capture` chỉ được báo hoàn tất sau khi file đã flush và ghi SHA. Port chỉ
bind `127.0.0.1`; token không được ghi vào screenshot hoặc log công khai.

### 9.3 Fixture xác định

CLI fixture versioned tạo hai account/character duy nhất theo run ID, không dùng
tài khoản `naruto` thủ công:

- cùng `Lang=0`, `MapID=1`, zone 0, `KieuToc=30`;
- head/body/leg phía client phải là 30/57/58;
- HP/chakra tối đa 150, đầy lúc bắt đầu;
- A tại `(100,384)`, B tại `(148,384)`;
- hai Sói Rừng template 0, mob ID 0/1, HP 100, vị trí từ catalog spawn;
- bốn skill Làng Lá level 1 đã gán;
- không trang bị/buff làm thay đổi damage.

Fixture fail-fast nếu DB hoặc gói ME_INFO/catalog khác các giá trị trên. Cleanup
chỉ xóa account mang đúng prefix/run ID.

### 9.4 Evidence và chuỗi E2E

Một run tạo hồ sơ:

```text
e2e-runs/<run-id>/
  manifest.json
  preflight.log
  unity-build.log
  build-sha256.txt
  clean/
    server-run-1.log
    server-run-2.log
    client-A.log
    client-B.log
    events.ndjson
    db-before.json
    db-after-shutdown.json
    db-after-relogin.json
    screenshots/
    result.json
  warm-upgrade/
    server.log
    client-A.log
    client-B.log
    events.ndjson
    cache-before.json
    cache-after.json
    screenshots/
    result.json
  aggregate-result.json
```

`manifest.json` ghi hash source/lock/server commit+dirty status/DB schema/Unity/
executable/profile roots/ports. `result.json` ghi schema version, từng
assertion, expected, actual, mã lỗi và đường dẫn evidence.
`aggregate-result.json` chỉ xanh khi clean full và warm-upgrade smoke dùng cùng
executable hash/server lock và cả hai xanh; schema ghi rõ `scope=full|warm_smoke`
để không diễn giải warm như full E2E. Mỗi checkpoint cần
assertion trạng thái máy đọc được và screenshot 1280×720, zoom x2, locale
Vietnamese, tên cố định. Screenshot đơn lẻ không phải bằng chứng giao thức hay
persistence. Golden/crop có ngưỡng làm cổng tự động; checklist reviewer hình ảnh
có tên và PASS/FAIL là cổng cuối.

Chuỗi bắt buộc:

1. A/B có cache riêng và hai tài khoản khác nhau.
2. Cùng map/khu, thấy nhau và đồng bộ di chuyển.
3. A mở Iruka bằng `open_npc`; nhận đúng server command 32, tên/text/menu và
   portrait 27078; chọn menu đóng bằng `choose_npc_menu`.
4. A thử một tọa độ walkable rồi một tọa độ blocked. Điểm đầu cập nhật ở cả
   client/server; điểm blocked bị cả hai từ chối và tọa độ không đổi. Trước/sau
   bước đi, bốn `background_layers` phải giữ đúng resource và thay offset theo
   parallax; có screenshot hai mốc.
5. A chọn/chạy lần lượt skill template 1..4. Mỗi skill phải trả đúng
   `selected_skill`, icon 27079..27082, pose checkpoint và event gameplay phù hợp;
   không load effect 80.
6. A đánh quái; cả hai thấy đòn, sát thương, chết và hồi sinh.
7. Sói ID 1 đánh A đúng một lần: event phải là `amount=5`,
   `hp_before=150`, `hp_after=145`. Ngay sau event, runner đưa A tới điểm an
   toàn ngoài aggro bằng `move_to` và khẳng định không có damage từ mob nào nữa
   cho tới khi kết thúc bốn đòn PvP ở bước 10.
8. A đánh B chết; B thấy mình chết và A thấy B chết.
9. B hồi sinh đầy HP/chakra. Đòn tại `+0,5s` không gây damage; đòn sau `+2,0s`
   gây damage. Quái hồi sinh được chấp nhận trong 4,5–6,5 giây sau event chết.
10. A tới tọa độ không mặc định và HP được đưa xác định về `137/150`: fixture
    dùng checkpoint Sói ở bước 7 (`150→145`), rồi B dùng skill cận chiến cấp 1
    bốn lần. Runner chỉ gửi đòn tiếp theo sau `skill_ready`/mốc
    `selected_skill_ready_at_ms`; mỗi đòn phải sinh đúng một event damage 2 với
    HP sau đòn lần lượt `143,141,139,137`, không sleep đoán và fail ngay nếu
    có event thừa hoặc damage khác. Lưu expected từ cả telemetry và DB.
11. Gửi shutdown qua control endpoint E2E của server, đợi log lưu hoàn tất và
   port đóng trong 15s;
   query DB sau shutdown phải bằng expected.
12. Khởi động lại, relogin trong 30s; telemetry và DB sau login phải bằng tuyệt
   đối map/khu/x/y/HP trước shutdown. B đăng nhập và thấy A đúng vị trí.

HUD helper có golden riêng ở 25/50/100% trước build; E2E gameplay còn đối chiếu
chiều rộng fill thực với `hp/hp_max` và `chakra/chakra_max` trong telemetry,
sai quá 1 px là đỏ.

Server control chỉ bật khi có env E2E, bind loopback + token, và lệnh `shutdown`
phải gọi đúng cùng hàm graceful shutdown mà signal Ctrl+C dùng. Evidence ghi
thời điểm nhận lệnh, save hoàn tất, exit code và thời điểm port đóng; `kill`/
`Stop-Process` không được tính.

Timeout mặc định: command/state 10s, server start/shutdown 15s, login 30s;
Unity cold import 20 phút và build 10 phút, đều cấu hình được nhưng luôn hữu hạn
và ghi duration theo phase. Runner chờ event/state, không `sleep` cố định ngoài
các mốc thời gian gameplay có tolerance ở trên.

## 10. TDD, mutation và review

Mọi task tuân theo:

1. Test mới đỏ đúng vì hành vi chưa có.
2. Code tối thiểu làm test xanh.
3. Chạy test hẹp, suite server `-count=1`, `go vet`, rồi E2E phù hợp.
4. Implementer tự review và commit.
5. Reviewer độc lập kiểm cả đúng spec lẫn chất lượng.
6. Critical/Important phải sửa và re-review trước task tiếp theo.
7. Chạy canary trên worktree/staging riêng; hoàn nguyên canary rồi chạy xanh lại.

Canary tối thiểu có mã lỗi ghim:

| Mutation | Chỗ áp dụng | Cổng phải đỏ |
|---|---|---|
| đổi endian/type/frame count part | artifact tạm | `ASSET-PART-001` |
| trỏ sprite ID thiếu | manifest tạm | `ASSET-SPRITE-002` |
| xóa x3 | artifact tạm | `ASSET-ZOOM-003` |
| phá map header / xóa `b03` | artifact tạm | `ASSET-MAP-004` / `ASSET-BG-005` |
| lệch tile và collision server | artifact tạm | `ASSET-COLLISION-006` |
| B đọc profile A | Unity snapshot | `E2E-CACHE-002` |
| bỏ broadcast chết tới A | server canary worktree | `E2E-BROADCAST-003`: A không thấy B chết |
| bỏ lưu HP shutdown | server canary worktree | `E2E-PERSIST-004`: DB/relogin sai `137/150` |
| writer A không đọc, queue đầy | test server | `SERVER-ACTOR-005`: A bị loại, probe actor <100ms, B đủ event |

Canary Unity chỉ sửa snapshot; canary server chỉ sửa worktree riêng. Mỗi lượt
ghi hash trước mutation, chạy đỏ đúng mã/expected–actual, hoàn nguyên, xác nhận
hash trở lại rồi chạy xanh. Một canary đỏ vì compile error hoặc lý do khác không
được tính là test bắt đúng bug.

Nếu canary vẫn xanh thì test liên quan chưa có giá trị và phải viết lại.

## 11. Chia việc song song

Chạy theo worktree/đường dẫn tách biệt:

- luồng A: manifest, parser, validator, assetpack và artifact lock;
- luồng B: Unity staging, profile/cache, telemetry và E2E runner;
- luồng C: outbound queue/writer và test actor không bị chặn;
- coordinator: image generation, tích hợp artifact, review package và nghiệm thu.

Không cho hai agent cùng sửa `server/internal/game/` hoặc cùng sửa Unity
snapshot. Luồng A không sửa gameplay Go trong lúc luồng C sửa actor. Sau khi
luồng C được review/merge, một task riêng mới nối catalog, NPC, collision,
`part.bin`/`skill.bin` vào `server/internal/game/`. Các giao diện giữa luồng
được ghim trong spec/plan trước khi giao.

## 12. Cổng hoàn thành Phase 1

Phase 1 chỉ hoàn thành khi đồng thời có bằng chứng:

- assetpack build lặp lại cùng hash, closure runtime không phụ thuộc Big/NR
  Dragon Ball và mọi mutation asset bị bắt;
- catalog map/NPC/mob, part, skill pose, tileset và collision có test byte/semantic
  độc lập; clean-cache và warm-cache đều render đủ;
- contact sheet x1 sáng/tối, golden state và checklist reviewer hình ảnh xanh;
- test Go, `go vet` và web build regression toàn repo xanh (web không phải
  deliverable gameplay của Phase 1);
- E2E Go hai người + PostgreSQL xanh;
- E2E Unity A/B tạo đủ hồ sơ bằng chứng và `result.json` xanh;
- chạy lại sau graceful restart giữ đúng trạng thái;
- canary actor, broadcast, persistence, cache và asset đều đỏ đúng lý do, sau
  hoàn nguyên đều xanh;
- reviewer cuối không còn Critical/Important;
- `NHAT-KY-AGENT.md` và `BAT-DAU-TU-DAY.md` được cập nhật cùng commit.

Sau cổng này mới nhân cùng pipeline sang bốn làng còn lại, atlas hóa tài nguyên
thường dùng và tiếp tục boss/vật phẩm/nhiệm vụ theo roadmap.

## 13. Việc cố ý chưa tuyên bố hoàn thành

- hình ảnh/default part của Làng Cát, Sương Mù, Đá và Mây;
- ba lựa chọn tóc trong màn tạo nhân vật;
- atlas/CDN/update resource công cộng và tương thích client cũ;
- projectile/effect catalog hoàn chỉnh của skill;
- vị trí quái di chuyển authoritative trên server;
- vật phẩm rơi/nhặt, hành trang, cửa hàng, giao dịch;
- boss, nhiệm vụ, bang hội và nội dung dài hạn.

Những việc này không được dùng asset placeholder để “làm như đã xong”. Chúng là
phase/KH tiếp theo sau khi cổng Làng Lá chứng minh pipeline và E2E đáng tin.
