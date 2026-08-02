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

Python version gate sau khi bump: `ProductionVersionContractTests` xanh (1 test,
0.002s), xác nhận `runtime_slice=22`, `map=22`, `data/resource=21`. Full
`python -m unittest discover -s assetpack/tests -v` xanh **124/124** trong
163.649s.

## Independent audit fix — collision ngoài footprint fail-open

Audit sau commit bàn giao `902b5a8` phát hiện `Collision.DiDuoc` trả `true` khi
`col >= width || row >= height`. Artifact thật chỉ có 5×17 cell, tile 24×24,
nên footprint hợp lệ là `x=0..119`, `y=0..407`; trong khi map 1 vẫn dùng biên
fallback 4096. Từ spawn `(100,384)`, bước tới `(121,384)` vì vậy từng thoát
artifact nhưng vẫn được nhận.

### RED và mutation

Thêm test trên artifact thật:

- `(119,407)` dùng đúng cell cuối; `(120,384)`, `(121,384)`, `(100,408)` và
  tọa độ âm đều bị chặn.
- segment `(100,384) -> (121,384)` bị chặn.
- `kiemDiChuyenHopLe` map 1 chặn bước `(100,384) -> (121,384)`.
- catalog vẫn quảng bá 6 template, nhưng map 6 `Khu rừng` được ghim là chưa
  playable: runtime chỉ tạo 5 map và `TimBanDo(6) == nil`.

Raw RED:

```text
--- FAIL: TestCollisionFailClosedNgoaiFootprintArtifact
    catalog_test.go: DiDuoc(120,384) = true, muốn false
--- FAIL: TestKiemDiChuyenLangLaKhongThoatKhoiCollisionGrid
    dichuyen_test.go: kiemDiChuyenHopLe((100,384) -> (121,384)) = true, muốn false
FAIL
```

Mutation: tạm đổi nhánh ngoài grid về `return true`; đúng hai test collision và
movement trên đỏ. Khôi phục `return false`; focused GREEN `-count=5`:

```text
ok  	nroserver/internal/game	0.191s
```

### Bản vá

- `Collision.DiDuoc` fail-closed với mọi tọa độ âm hoặc cell ngoài grid.
- Map 1 lấy footprint collision làm biên authoritative; `BienRongToiDa=4096`
  chỉ còn là fallback cho map chưa gắn artifact.
- Các fixture test cũ của map 1 được đưa về footprint thật, không nới production
  boundary để chiều test.
- Comment và test ghim chủ ý catalog có 6 template nhưng runtime chỉ có 5 map
  playable; map 6 chưa có tile/render/collision contract.

### Fresh verification sau fix

```text
python -m unittest discover -s assetpack/tests -v
Ran 124 tests in 285.481s
OK

D:\Go\bin\go.exe test ./internal/game -count=1
ok  	nroserver/internal/game	83.940s

D:\Go\bin\go.exe test ./internal/game -run 'Test(NapCatalog|TaiCollision|Collision|KiemDiChuyenLangLa|CatalogMap6|GoiCatalogMapLangLa|GoiMoIruka|MoNPCXepGui|TaoQuaiMacDinh)' -count=5
ok  	nroserver/internal/game	0.191s

D:\Go\bin\go.exe test ./internal/net -count=5
ok  	nroserver/internal/net	0.280s

D:\Go\bin\go.exe vet ./...
(exit 0, không có output)
```

Ba thư mục `.task8-build-2..4`, hai `__pycache__` và log kiểm thử tạm đã được
kiểm tra là untracked output rồi dọn. Task 8 đang **chờ independent re-review**;
không được mở Task 9 hoặc merge trước kết luận mới.
