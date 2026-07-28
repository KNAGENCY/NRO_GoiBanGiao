# Task 3 re-review — đóng golden closure và fail-closed validation

## Kết quả

Đã sửa các điểm NEEDS WORK trong `assetpack/`:

- Mob parse binary thật theo `EffectData.readData(myReader)` legacy (big-endian image/frame/arrFrame và 2 byte reserved trailer), dùng fixture độc lập trích từ `D:\NRO_Project\Assets\Resources\res\x2\effectdata\106\data.bytes`; header/length/payload đổi dù cập nhật SHA bị `ASSET-MOB-003` chặn.
- Build bắt buộc đủ closure animation `3/17/14/10`, exact name duy nhất, cell bounds, slot và `client_cf` golden độc lập từ `Char.CharInfo[0..32]`; `client_cf=[99]`, malformed object, duplicate hoặc extra bị `ASSET-ANIMATION-008` chặn.
- UI bind runtime-key map đủ đúng 8 binding duy nhất, đúng tỷ lệ `1 direct / 2 slice_group / 5 overlay_required`, logical/source/physical không duplicate; logical key/future key phải có call-site/consumer precondition và SHA file khớp.
- `tile_set_info.bin` được đối chiếu bytes golden và semantic header độc lập; payload thay đổi dù SHA manifest đổi vẫn bị `ASSET-TILESET-011`. Tile size 0 và tileset malformed trả `BuildError` ổn định.
- Lock bổ sung provenance manifest/runtime-key map, dimensions, IDs, frame counts và tách `review_evidence` khỏi `generated_evidence`; output vẫn build vào thư mục tạm rồi replace một lần, không ghi đè output có sẵn.

## Kiểm thử

- `python -m unittest discover -s assetpack/tests -v`: **68/68 PASS**.
- Chạy lặp `python -m unittest discover -s assetpack/tests -q`: **68/68 PASS**.
- `python -m compileall -q assetpack`: **PASS**.
- Mutation thật đã làm đỏ đúng oracle cho Mob binary/header/SHA, animation `client_cf`, duplicate closure, duplicate runtime binding và tileset golden; mọi mutation đã khôi phục và suite xanh lại.

## Known blocker

Build production vẫn dừng fail-closed ở closure manifest nếu Task 6B chưa bind đủ source/evidence/build inputs thật. Test build dùng fixture canonical đầy đủ để chứng minh pipeline không bỏ qua required closure; chưa tự ý sửa PNG hay manifest production.
