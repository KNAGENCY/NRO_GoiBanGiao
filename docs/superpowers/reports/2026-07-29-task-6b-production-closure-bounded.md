# Task 6B production closure — bounded ArtRoot base fix

Ngày: 29/07/2026
Worktree: `kh-vertical-slice`

## Phạm vi vòng này

- `assetpack.build.resolve_manifest_base()` là cổng dùng chung để resolve
  `build_inputs.source_base` từ `--base-dir` thật của repository.
- Source production phải là đường dẫn tương đối từ build base; source ngắn theo
  thư mục mapping hoặc source không tồn tại bị từ chối trước khi staging output.
- `assetpack.bind.bind_manifest()` dùng cùng cổng base này khi manifest đã khai
  báo `source_base`.
- Không sửa PNG nguồn, `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1`,
  fixture Dragon Ball, BAT hoặc NHAT-KY.

## RED / GREEN

- RED: 2 test canonical base lỗi vì `resolve_manifest_base` chưa tồn tại.
- GREEN: 2/2 test pass sau khi thêm cổng resolve:
  - đường dẫn `docs/my-thuat/anh-nguon/vertical-slice-v1/iruka_portrait_chroma.png`
    được resolve từ `D:\NRO_GoiBanGiao`;
  - `iruka_portrait_chroma.png` dạng short path bị trả `ASSET-REVIEW-010`.
- QA/art generator unit: **36/38 PASS** trong nhóm
  `tools.artqa.test_artqa`, `tools.artqa.test_vertical_slice_generator`,
  `assetpack.tests.test_task6b_round2`; 2 test production hiện vẫn đỏ đúng vì
  mapping chưa canonical.

## Fresh production probe

Với `--base-dir D:\NRO_GoiBanGiao`:

| Lệnh | Exit | Kết quả |
|---|---:|---|
| Art QA gate | 0 | `allowed=true`, `issues=[]` |
| Generator 42 sprite | 1 | `ASSET-ART-004: sprite 27041 crop outside alpha evidence` |
| Bind | 1 | `ASSET-REVIEW-010: thiếu source đã review: master_hoc_vien_x4_chroma.png` |
| Legacy verify | 1 | `ASSET-LEGACY-009: SHA legacy lệch: assetpack/tests/fixtures/legacy/mob.bin` |
| Build lần 1 | 1 | `ASSET-REVIEW-010: thiếu source đã review: iruka_portrait_chroma.png` |
| Build lần 2 | 1 | cùng blocker; không tạo output |

Generator, bind và cả hai build đều không để lại output staging. Probe ID:
`5923a906d70f4d04a14bfa61484ca267`.

## Blocker production chưa được phép giả PASS

`docs/my-thuat/mapping-id.json` hiện vẫn là mapping baseline: thiếu top-level
`skills`, `ui_bindings`, `build_inputs`, thiếu `source_base`; 42 sprite zoom
`x1..x3` còn trỏ `planned/...`, còn `x4` trỏ basename raw. Các crop sheet hiện
không chứng minh được trong alpha canvas canonical (alpha evidence của master
Học viên là `542×1046`, không phải sheet `1254×1254`).

Ngoài sprite binding, repository/shared ArtRoot chưa có input Naruto độc lập
đã chứng minh cho 100 tile × 4 zoom, b00..b03, 8 UI source, Mob bytes/data và
legacy `NR_part`/`NR_image`/skill/mob/Small/C# ID. Legacy ledger hiện còn chứa
fixture; không được đổi SHA để che mismatch.

Vì vậy vòng này chỉ commit cổng base + test + báo cáo blocker. Chưa sinh hoặc
ghi manifest/artifact production, không tuyên bố bind/build xanh.
