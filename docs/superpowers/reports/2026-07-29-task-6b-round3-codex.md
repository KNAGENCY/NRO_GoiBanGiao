# Task 6B round 3 — bounded generator và blocker mapping

Ngày: 29/07/2026  
Worktree: `kh-vertical-slice`

## RED/GREEN

- RED: test generator mới fail vì chưa có module `tools.artqa.generate_vertical_slice`.
- GREEN: thêm generator Pillow deterministic và 4 test; `python -m unittest
  tools.artqa.test_vertical_slice_generator -v` → **4/4 PASS**.
- Generator chỉ đọc 20 PNG nguồn, kiểm SHA source + đủ 4 evidence SHA, dùng alpha
  evidence đã review làm canvas canonical, sinh x1/x2/x3/x4 RGBA và ghi SHA/dimensions
  vào `generated-index.json` trong thư mục output atomic.
- Mutation đổi byte `iruka_portrait_chroma.png` → **ASSET-ART-001**, không có output.
- Mutation bỏ Iruka khỏi review → **ASSET-ART-002**.

## Blocker production còn lại — không giả PASS

`mapping-id.json` chưa thể bind vào alpha evidence hiện có. Cụ thể sprite `27043`
khai crop nguồn `[418, 0, 209, 209]` trên `master_hoc_vien_x4_chroma.png`, nhưng
alpha evidence đã review là canvas `542×1046`; crop đó vượt canvas (`x+width=627`).
Các sprite còn lại cần được đối chiếu lại cùng cách, không được cắt lại theo tọa độ
sheet chưa có evidence tương ứng. Generator fail-closed bằng **ASSET-ART-004** và
không tạo output mapping.

Ngoài crop binding, build canonical vẫn thiếu các input không có trong 20 PNG:
`tile_set_info.bin` production, 100 tile x4, b00..b03 x4, 8 UI binding source,
Mob bytes/data production và các legacy input thật (`NR_part`, `NR_image`, skill,
mob, Small, hard-coded IDs). Vì vậy chưa sửa `mapping-id.json`, chưa đổi ledger/hash,
chưa tuyên bố art gate/bind/build production xanh. Kiểm tra hiện tại cũng trả
`ASSET-LEGACY-009: SHA legacy lệch: assetpack/tests/fixtures/legacy/mob.bin`
(ledger khai 11 byte nhưng fixture hiện có 227 byte); không sửa hash mù.
