# Task 3 re-review — đóng golden closure và fail-closed validation

## Kết quả

Đã sửa các điểm NEEDS WORK trong `assetpack/`:

- Mob parse binary thật theo `EffectData.readData(myReader)` legacy (big-endian image/frame/arrFrame và 2 byte reserved trailer), dùng fixture độc lập trích từ `D:\NRO_Project\Assets\Resources\res\x2\effectdata\106\data.bytes`; header/length/payload đổi dù cập nhật SHA bị `ASSET-MOB-003` chặn.
- Build bắt buộc đủ closure animation `3/17/14/10`, exact name duy nhất, cell bounds, slot và `client_cf` golden độc lập từ `Char.CharInfo[0..32]`; `client_cf=[99]`, malformed object, duplicate hoặc extra bị `ASSET-ANIMATION-008` chặn.
- UI bind runtime-key map đủ đúng 8 binding duy nhất, đúng tỷ lệ `1 direct / 2 slice_group / 5 overlay_required`, logical/source/physical không duplicate; logical key/future key phải có call-site/consumer precondition và SHA file khớp.
- Fix4 khóa collision bằng fingerprint của `logical_keys`, `future_logical_key`, `physical_paths`, `call_sites` và `consumer_sites`, không dùng `source_asset` để che duplicate; copy entry 3 sang 4 bị `ASSET-UI-007` chặn.
- `validate_animation_closure` kiểm tra kiểu `name` và slot trước mọi `set`/lookup; list, dict và input malformed trả ổn định `ASSET-ANIMATION-008`, không traceback `TypeError`.
- `tile_set_info.bin` được đối chiếu bytes golden và semantic header độc lập; payload thay đổi dù SHA manifest đổi vẫn bị `ASSET-TILESET-011`. Tile size 0 và tileset malformed trả `BuildError` ổn định.
- Lock bổ sung provenance manifest/runtime-key map, dimensions, IDs, frame counts và tách `review_evidence` khỏi `generated_evidence`; output vẫn build vào thư mục tạm rồi replace một lần, không ghi đè output có sẵn.

## Kiểm thử

- `python -m unittest discover -s assetpack/tests -q`: **70/70 PASS**, chạy lặp 3 lần.
- Hai reproduction reviewer: **2/2 PASS**.
- `python -m compileall -q assetpack`: **PASS**.
- Mutation thật fix4: thêm lại `source_asset` làm test collision đỏ; tắt guard `name` làm 4 case `TypeError`; đã khôi phục và suite xanh lại.

## Task 3 fix5 — final animation and binding guards

- Runtime UI collision gate now rejects each duplicate `logical_keys`, non-empty `future_logical_key`, and normalized `physical_paths` independently, so changing `source_asset`, `call_sites`, or `consumer_sites` cannot hide a collision; all failures remain `ASSET-UI-007`.
- `validate_animation_closure` now requires `slot_count` on both sheet and animation contracts, accepts only canonical non-boolean integers, and reports malformed values as `ASSET-ANIMATION-008` without raw `KeyError`.
- Added exactly two reproduction tests: one covers all three runtime collision dimensions and one covers missing/bool/float/wrong slot counts for both closures. Minimal reproduction tests pass; real mutations made both reproductions fail.

## Known blocker

Build production vẫn dừng fail-closed ở closure manifest nếu Task 6B chưa bind đủ source/evidence/build inputs thật. Test build dùng fixture canonical đầy đủ để chứng minh pipeline không bỏ qua required closure; chưa tự ý sửa PNG hay manifest production.
