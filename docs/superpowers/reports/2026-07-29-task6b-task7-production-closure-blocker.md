# Task6B/Task7 canonical production closure — blocker

Ngày: 29/07/2026
Worktree: `kh-vertical-slice`
Trạng thái: **BLOCKED, không tuyên bố production PASS**

## Đã kiểm tra

- Art gate: `python -m tools.artqa.artqa gate --review docs/my-thuat/review/vertical-slice-v1.json`
  trả `allowed=true`, `issues=[]`, exit `0`.
- Shared ArtRoot có đủ 20 PNG đã review. SHA nguồn được giữ nguyên trong review.
- Tile atlas thật: `D:/NRO_GoiBanGiao/docs/my-thuat/anh-nguon/vertical-slice-v1/generated/naruto-village-tile-sheet-10x10-v2.png`, 1254×1254,
  SHA-256 `2031edef163fc766b9f44df6c30724c0e41c57730b4f98749e3940fae157036b`.
- Input thật trên đĩa: `server/internal/game/part.bin` 87707 byte,
  SHA-256 `09349f28a95b46469b691abf87d6b341f3395f53ca8d6d704b48ab706992826b`;
  `server/internal/game/tileset.bin` 1696 byte,
  SHA-256 `27aa6258210b5de9e79e164a711065c4cef4dad679b68db2ac4430bddc969027`.
- Mob data 227 byte trên máy có cùng SHA với fixture:
  `D:/NRO_Project/Assets/Resources/res/x2/effectdata/106/data.bytes`,
  SHA-256 `b7602310dbb1d0c668b972d27713b24a1d74feea7a9bb1b8398673df8273e59c`.
  Bản sao trong repo được ghi rõ là `committed_test_fixture`, `production_approved=false`;
  không được coi là input production.
- Legacy ledger sau khi sửa provenance chạy được `verify-legacy-inputs` exit `0`,
  nhưng đó chỉ là kiểm tra byte/SHA, không phải phê duyệt nguồn production.

## Blocker còn lại

1. `mapping-id.json` hiện vẫn là baseline, thiếu canonical `skills`, `ui_bindings`,
   `build_inputs`, `build_inputs.source_base`; 42 sprite zoom còn `planned/...`.
2. Generator fail-closed ở sprite `27041` (`ASSET-ART-004`). Alpha evidence của
   `master_hoc_vien_x4_chroma.png` chỉ là canvas `542×1046`, trong khi các crop 209×209
   hiện tại lấy tọa độ sheet 1254×1254. Không có evidence độc lập cho 34 crop
   `27041..27074`, nên không được tự suy ra hay bind.
3. Chưa có nguồn Naruto production đã duyệt cho `b00..b03` ở đủ `x1..x4`,
   cũng như 8 UI binding có call-site/consumer-site và SHA C# chính xác. Các file
   `D:/NRO_Project/Assets/Resources/res/x2..x4` là tài nguyên client legacy Dragon Ball,
   không được dùng để tạo PASS Naruto.
4. Chưa có skill/`NR_part`/`NR_image` Naruto production tương ứng với mapping. Các
   cache `NR_*` trên máy là cache game legacy khác; fixture skill/part vẫn chỉ là fixture.

Vì các blocker trên, chưa chạy được bind production và chưa thể build hai lần để
so sánh `output/asset-lock.json` byte-identical. Không tạo artifact giả, không đổi SHA
để ép xanh, không sửa PNG/BAT/NHAT-KY.

## Test đỏ đã ghi

`assetpack/tests/test_task6b_canonical_preflight.py` cố tình đỏ với manifest hiện tại,
ở các assertion canonical top-level và SHA zoom/legacy. Đây là acceptance test cho lần
đóng tiếp theo; không được đổi thành PASS giả khi input thật chưa có.
