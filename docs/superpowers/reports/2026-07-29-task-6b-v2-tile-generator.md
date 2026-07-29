# Task 6B — v2 tile generator closure

Ngày: 29/07/2026
Worktree: `kh-vertical-slice`

## Đã chốt

- Nguồn chỉ đọc: `D:/NRO_GoiBanGiao/docs/my-thuat/anh-nguon/vertical-slice-v1/generated/naruto-village-tile-sheet-10x10-v2.png`.
- SHA-256: `2031edef163fc766b9f44df6c30724c0e41c57730b4f98749e3940fae157036b`.
- QA nguồn: PNG `1254x1254`, lattice `10x10`, đủ 100 ô; `nonblack_ratio_min=0.825856`, `nonblack_ratio_max=0.930624`.
- Generator sinh atomically 100 tile, mỗi tile có `x1/x2/x3/x4`, ghi SHA/dimensions và `generated-index.json`.
- Mutation append byte vào nguồn bị chặn bằng `ASSET-TILE-001` trước khi tạo output.
- Output probe nằm ngoài worktree tại `.artifacts/task6b-v2-tiles`; không stage output và không sửa nguồn v1/raw.

## Kiểm thử

`python -m unittest tools.artqa.test_tile_sheet_generator -v` — **2/2 PASS**.

## Production status

V2 generator đã đạt closure riêng, nhưng **chưa được bind vào `mapping-id.json` trong commit này**. Vì vậy production gate/build vẫn fail-closed ở manifest baseline: mapping hiện chưa có `skills`, `ui_bindings`, `build_inputs.source_base` và `tile_art`; các sprite zoom `planned/...` cũng chưa có output thật. Bước kế tiếp là bind manifest sau khi có đủ input production độc lập cho các blocker này.
