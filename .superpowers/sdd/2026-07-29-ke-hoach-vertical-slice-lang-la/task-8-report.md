# Task 8 report — Catalog, NPC Sói và collision Làng Lá

## RED trước implementation

`D:\Go\bin\go.exe test ./internal/game -run 'TestNapCatalog|TestTaiCollision|TestGoiMoIruka' -count=1`

Kết quả đỏ đúng mục tiêu: build fail vì `NapCatalog`, `TaiCollision` và
`goiMoNPC` chưa tồn tại.

## GREEN

- `catalog.json` và `asset-lock.json` được parse lúc khởi động; catalog rỗng,
  thiếu binding, đổi thứ tự binding và hash không khớp đều fail-closed.
- `-28/6` gửi 6 map, Iruka (parts từ lock) và Sói template 0; test đọc đủ body
  đến EOF.
- `-24` map 1 có NPC Iruka; mob mặc định dùng `MauId=0`.
- Command 33 trả 32 (text, menu `Đóng`, portrait 27078) và test EOF.
- Movement chỉ kiểm collision trong actor; foot point/segment qua solid bị chặn.

Focused GREEN:

`D:\Go\bin\go.exe test ./internal/game -run 'TestNapCatalog|TestTaiCollision|TestGoiMoIruka|TestGoiCatalogMapLangLa|TestGoiMapInfoVoiQuaiByteDung|TestKiemDiChuyen' -count=1`

Kết quả: `ok nroserver/internal/game 0.038s`.

Game suite GREEN trước các assert EOF cuối:

`D:\Go\bin\go.exe test ./internal/game -count=1`

Kết quả: `ok nroserver/internal/game 30.410s`.

## Mutation/canary

- Đổi một byte collision rồi gọi loader: đỏ `ASSET-COLLISION-006`.
- Bỏ portrait khỏi `sprite_ids`: catalog fail-closed (`binding`).
- Đảo thứ tự `part_ids`: catalog fail-closed (`binding`).
- Collision synthetic có solid: cả foot point và segment crossing bị chặn.

## Verification còn treo

Đã dò PostgreSQL thực tế ở cổng 5432. Lệnh bắt buộc
`go test ./internal/game ./internal/net -count=5` chạy quá 364 giây và timeout
(exit 124, không có failure output), nên chưa thể khẳng định pass lặp 5/full/vet.
Cần chạy lại tách từng package hoặc điều tra timeout trước khi merge.

## Concern

Artifact catalog hiện chỉ chứa tên/spawn; câu thoại/menu Iruka không nằm trong
catalog/lock. Runtime phải giữ text/menu tối thiểu ở `npc.go`; artifact nên được
mở rộng để đưa dialog vào nguồn dữ liệu nếu muốn áp dụng tuyệt đối quy tắc
"không hard-code catalog" cho nội dung thoại.
