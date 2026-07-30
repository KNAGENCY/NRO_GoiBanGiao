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

## Fix round 1 — catalog explicit và collision semantic

RED: collision test mới gọi `tile_size` trên API cũ và nhận `TypeError`. Mutation
handler NPC (`XepGui` đổi thành `Gui`) làm test client `net.Pipe` không đọc output
đỏ sau 250ms.

GREEN: collision artifact là `[width,height,tileWidth,tileHeight,0|1...]`, ghim
24x24, spawn walkable và có solid; build từ `client_solid`, không serialize tile
index. Catalog chứa map spawn, NPC (parts/portrait/toạ độ/dialog/menu) và mob
(hp/range/speed/dart/spawns); Go kiểm từng ID với closure, không suy diễn index.
Command 33 dùng `XepGui`, test chứng minh trả nhanh, FIFO và byte chính xác.
`vsMap=22` khớp lock `map=22`, `runtime_slice=22`.

Đã xanh: Python `test_map` + `test_build_contract`: 25 tests/26.602s; focused Go
catalog/collision/NPC/queue/version `-count=5`: OK 0.180s; `go test ./internal/net
-count=5`: OK 0.284s; `go vet ./...`: OK. Hai lần build production deterministic:
catalog `f24db071ce8708d4fb8e318f934a702cbed7090bb709fc76a42c752d4617a973`,
collision `4fe176b20d67dd98bd7b96421ab741d07f0d045b9e61b7f15efcbbeaadfd043b`, lock
`f72a20b2bd50408dab472f5baa2c9d7ed3edd523e5b38784e11f64c697dea29f`.

Fix collision placement: canary solid dời từ cell `(0,0)` sang `(4,0)`, tránh
đường actor test bắt đầu tại `(0,0)`. `TestTaiCollision...` mutation cập nhật cell
index 4 (offset 8 sau header). Full `D:\Go\bin\go.exe test ./internal/game -count=1`
đã xanh, 30.451s.

Artifact sau khi dời canary: collision
`39eee2c743a6a27c1b8afe3e640c9dc1f611685b3a7e0b81e32538234ad7cd5f`; lock
`54a638f0e7c008586dedbf13eaf45aff42edcdf63f6708fd6138b5219c773d82`.

## Concern

Artifact catalog hiện chỉ chứa tên/spawn; câu thoại/menu Iruka không nằm trong
catalog/lock. Runtime phải giữ text/menu tối thiểu ở `npc.go`; artifact nên được
mở rộng để đưa dialog vào nguồn dữ liệu nếu muốn áp dụng tuyệt đối quy tắc
"không hard-code catalog" cho nội dung thoại.
