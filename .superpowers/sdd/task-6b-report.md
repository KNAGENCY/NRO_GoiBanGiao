# Task 6B — báo cáo bind/closure

Ngày kiểm tra: 29/07/2026. Worktree: `kh-vertical-slice`.

## Kết quả

- Kiểm tra manifest hiện tại: `mapping-id.json` có 12 nhóm trường, thiếu đúng
  `skills`, `ui_bindings`, `build_inputs` so với `assetpack.build.REQUIRED_TOP_LEVEL`.
  Vì vậy đây là mapping baseline, chưa phải canonical build manifest.
- Review gate thật:
  `python -m tools.artqa.artqa gate --review docs/my-thuat/review/vertical-slice-v1.json`
  → exit `0`, `allowed: true`, `issues: []`.
- Test RED: import `assetpack.bind` thất bại vì helper chưa tồn tại.
- Test GREEN: `python -m unittest assetpack.tests.test_bind -v` → `6/6` xanh.
  Các mutation metadata gồm trùng slot, đổi source SHA, crop, anchor, client_cf,
  thiếu evidence và thiếu runtime key đều bị từ chối; gate fail không tạo output.
- Build thử với mapping/review hiện tại → exit `1`, không tạo output:
  `ASSET-REVIEW-010: thiếu source_sha256/evidence_path`.

## Blocker không được phép giả PASS

1. Cả 20 asset review đều không có anchor; review chỉ có crop box/animation
   metrics. Bind dừng fail-closed ở `ASSET-ANCHOR-009` (asset đầu tiên:
   `hu_anthan_10f_chroma.png`). Không thể điền anchor `0,0` thay thế.
2. Review đang lưu 4 evidence SHA lồng trong `evidence{}`; contract build hiện
   tại còn yêu cầu `evidence_path` và `evidence_sha256` dạng phẳng cho từng asset.
   Do đó `assetpack.build` chưa nhận trực tiếp review hiện tại.
3. Mapping chưa khai báo closure canonical: 4 animation sheets (3/17/14/10),
   4 skill, Mob bytes/data/image x1..x4, map/tile/background, HUD và 8 UI
   binding/runtime-key. Không có dữ liệu đủ để tự suy ra các trường này.
4. `legacy-input-ledger.json` đang trỏ các input ngoài `part.bin` tới
   `assetpack/tests/fixtures/legacy/*`; đó là fixture kiểm thử, không phải closure
   production. Không được dùng chúng để tạo PASS production.

Vì các blocker trên, chưa sửa `mapping-id.json`, chưa sinh artifact staging và
không ghi nhận PASS giả. Commit này chỉ chứa guard/test/helper và báo cáo để
phiên sau bổ sung anchor + canonical review/closure thật.
