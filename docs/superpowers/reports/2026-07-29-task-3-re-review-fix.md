# Task 3 re-review — đóng golden closure và fail-closed validation

## Kết quả

Đã sửa các điểm NEEDS WORK trong `assetpack/`:

- Mob đọc và parse cả hai payload, đối chiếu độc lập với golden frame/action contract. Đổi payload rồi chỉ cập nhật SHA bị `ASSET-MOB-003` chặn.
- Build bắt buộc đủ closure animation `3/17/14/10`, đúng sheet, slot, `client_cf`, frame count, cell size và bounds; thiếu/xóa sheet hoặc cell sai kích thước bị `ASSET-ANIMATION-008` chặn.
- UI bind runtime-key map đủ 8 entry, đúng tỷ lệ `1 direct / 2 slice_group / 5 overlay_required`, logical key phải xuất hiện trong call-site, physical path canonical và SHA call-site/consumer phải khớp precondition.
- `tile_set_info.bin` được đối chiếu bytes golden và semantic header độc lập; payload thay đổi dù SHA manifest đổi vẫn bị `ASSET-TILESET-011`. Tile size 0 và tileset malformed trả `BuildError` ổn định.
- Lock bổ sung provenance manifest/runtime-key map, dimensions, IDs, frame counts và tách `review_evidence` khỏi `generated_evidence`; output vẫn build vào thư mục tạm rồi replace một lần, không ghi đè output có sẵn.

## Kiểm thử

- `python -m unittest discover -s assetpack/tests -v`: **63/63 PASS**.
- Chạy lặp `python -m unittest discover -s assetpack/tests -q`: **63/63 PASS**.
- `python -m compileall -q assetpack`: **PASS**.
- Mutation thật đã làm đỏ đúng oracle cho Mob, animation closure, UI runtime-key và tileset golden; mọi mutation đã khôi phục và suite xanh lại.

## Known blocker

Build production vẫn dừng fail-closed ở closure manifest nếu Task 6B chưa bind đủ source/evidence/build inputs thật. Test build dùng fixture canonical đầy đủ để chứng minh pipeline không bỏ qua required closure; chưa tự ý sửa PNG hay manifest production.
