# Task 7 report — build artifact và deploy staging theo lock

- Base: `a7062e9e6d0120b7c5cc0b996ab212db180b0413`
- Worktree/branch: `.worktrees/kh-vertical-slice` / `kh-vertical-slice`
- Trạng thái: `DONE`
- Commit: `2a77a43 feat(runtime): nhung artifact Naruto va khoa cap server client`

## Phạm vi

- Tạo `assetpack/deploy.py`, `assetpack/tests/test_deploy.py`.
- Thêm CLI `python -m assetpack deploy`.
- Bump bốn version manifest từ 1 lên 21 theo quyết định coordinator: vừa đúng
  lock, vừa tăng đơn điệu qua staging 20 và không trùng legacy 10/2/1/4.
- Giữ nguyên legacy input `server/internal/game/part.bin` và `tileset.bin`;
  generated runtime nằm riêng trong `server/internal/game/assets/`.
- Go nhúng đúng `part.bin`, `skill.bin`, `tile_set_info.bin`,
  `map-1-collision.bin`, `catalog.json`, `asset-lock.json`.
- Gói `-87` mang đúng version, length và nguyên byte part + skill.
- Mở tối thiểu Task 5 để sửa exact whitespace precondition của
  `ui-command-disabled.patch`, theo cấp phép của coordinator.

## RED trước code

1. Python deploy:

   `python -m unittest assetpack.tests.test_deploy -v`

   RED: `ModuleNotFoundError: No module named 'assetpack.deploy'`.

2. Go runtime:

   `go test ./internal/game -run
   'TestRuntimeServerNhungDungByteVaSHATuLock|TestGoiVersionGameDungByte|
   TestGoiCatalogDataDungCauTruc' -count=1`

   RED compile đúng vì thiếu `duLieuAssetLock`, `duLieuSkill`,
   `duLieuCollision`, `duLieuCatalog`.

3. Version production:

   `python -m unittest
   assetpack.tests.test_deploy.ProductionVersionContractTests -v`

   RED: manifest đang là `1/1/1/1`, expected `21/21/21/21`.

4. Unity overlay thật:

   Lượt apply đầu trên `D:\NRO_Project` đỏ
   `OVERLAY-PATCH-001: Command.cs ... actual 0`. Source SHA vẫn đúng
   `b053ac...`; nguyên nhân là source decompiler có `\n    \t{`, còn patch/fixture
   Task 5 chỉ ghim `\n\t{`.

5. TDD sửa Task 5:

   Đổi fixture thành exact whitespace source thật rồi chạy
   `Test-Overlay.ps1`: RED `OVERLAY-PATCH-001 actual 0`.
   Sửa patch exact whitespace, không regex, không bỏ SHA/precondition.

6. Deploy thật lượt đầu:

   RED `ASSET-PREFLIGHT-022` vì runtime key map lịch sử ghi `mainimage`, output
   lock trên Windows giữ casing `mainImage`. Sửa bằng lookup case-fold nhưng
   cấm hai output collision sau case-fold.

7. Review fail-closed trước commit phát hiện deploy chưa đối chiếu SHA của
   runtime key map với provenance trong lock. Thêm test đổi một byte map sau
   build: RED vì deploy chưa bắt; sau đó ghim SHA lock và GREEN.

## GREEN hẹp

- `python -m unittest assetpack.tests.test_deploy -v`: 9/9.
- Part + manifest + deploy: 45/45.
- Go runtime focused: PASS.
- `tools/unity-overlay/tests/Test-Overlay.ps1`: 34 assertions PASS, gồm mutation
  SHA source.
- `python -m assetpack validate ...`: exit 0.

## Build A/B

- Artifact A:
  `.superpowers/.../task7-build-21-a`
- Artifact B:
  `.superpowers/.../task7-build-21-b`
- A files: 669
- B files: 669
- Diff toàn bộ path/SHA: 0
- `asset-lock.json` SHA:
  `e5b09320cf257bf260e3ef7a0de8d62acd751182ea0b6cda8cc828df4bb8cb07`
- Server/client pair SHA:
  `547157a1fa5aaac454720666b5a30f22806e93e14f4fa105df18a0db12eb2d79`
- Versions: runtime/data/map/resource = `21/21/21/21`
- Output trong lock: 668 + lock = 669 file.

## Unity snapshot và preflight thật

- Source read-only: `D:\NRO_Project`
- Snapshot:
  `.superpowers/.../task7-unity-snapshot-21b`
- Overlay apply: `OVERLAY-OK`.
- Release immutable:
  `.superpowers/.../task7-releases-21/e5b09320cf257bf260e3ef7a0de8d62acd751182ea0b6cda8cc828df4bb8cb07`
- Unity:
  `D:\Unity\2022.3.62f2\Editor\Unity.exe -batchmode -nographics ...`
- Log:
  `<release>/preflight.log`
- Marker log: dòng 6693 `RESOURCE-PREFLIGHT-OK`.
- Contract: đúng 8 UI binding, 100 tile cho x1–x4, năm nhóm
  `small/background/mob/hud/skill`.
- Snapshot `.nro-asset-lock.json` ghim đúng lock SHA, pair SHA và version 21.

## Mutation

1. Python đổi byte runtime, cập nhật SHA/pair trong lock nhưng giữ version:
   `test_tu_choi_runtime_doi_byte_nhung_version_khong_tang` PASS vì deploy bắt
   `ASSET-VERSION-015`.
2. Go lật bit cuối `assets/part.bin`, giữ nguyên lock:
   `TestRuntimeServerNhungDungByteVaSHATuLock` RED đúng SHA
   (`14c62c...` != `9f9c8e...`), sau đó restore byte-identical SHA
   `9f9c8ee290d8ff367ad7450584d7c1544adc41916924d9ec4a41a93229e6f327`.

## Fix round review: deploy fail-closed

### Finding

1. `previous_lock_path` va `protected_unity_root` la optional, nen caller Python/CLI
   co the bo qua khoa version hoac bao ve Unity source.
2. Preflight cho fallback x2-x4 ve x1 va chap nhan metadata dimension thieu/sai.

### Thay doi

- API `deploy()` va CLI `assetpack deploy` bat buoc `previous_lock` va
  `protected_unity_root`; API cung fail closed neu caller truyen `None`.
- Luon cam snapshot bang hoac nam trong Unity root duoc bao ve.
- Moi nhom small/background/Mob/HUD/skill phai co output dung x1-x4; metadata
  `width`/`height` phai la `int` duong va path dung zoom. Contract duoc tao va
  kiem truoc khi bat dau chep vao snapshot.
- Sua setup test dung lai `self.protected_root`, khong tao lai `unity-source`.

### RED/GREEN

- RED dau vao: `python -m unittest assetpack.tests.test_deploy -v` -> 22 failures,
  1 setup error (tao trung `unity-source`), dung cac finding review.
- GREEN: cung lenh -> 13/13 PASS.

### Mutation do

- Noi `previous_lock_path` thanh optional: API guard test do (thieu TypeError,
  runtime cho `ASSET-VERSION-015`).
- Noi `protected_unity_root` thanh optional: API guard test do (thieu TypeError,
  runtime cho `ASSET-DEPLOY-SOURCE-020`).
- Khoi phuc fallback x1 va bo check zoom: exact-zoom test do 15 subcase, subcase
  dau tien deploy duoc thay vi `ASSET-PREFLIGHT-022`.
- Khoi phuc fallback dimension: metadata test do 3 subcase, thieu width deploy duoc
  thay vi `ASSET-PREFLIGHT-022`.
- Da hoan nguyen tung mutation truoc khi kiem chung cuoi.

### Full suite

- `python -m unittest discover -s assetpack/tests -v` -> 122/122 PASS trong
  293.692s.

### Concern

- Caller deploy cu can bo sung ca `previous lock` va `protected Unity root`; day la
  thay doi co chu dich de staging khong the chay fail-open.

## Fix round 2: canonical resource per zoom

### Finding

- `small` chi tim mot file bat ky trong `smallimage` o tung zoom. Vi vay khi x2
  thieu `Small27041.png` nhung con `Small27042.png`, deploy van qua gate.

### Thay doi

- Contract preflight lay sprite, skill va mob canonical tu `closure` trong
  `asset-lock.json`, roi yeu cau dung duong dan canonical o x1-x4.
- Background va HUD cung dung exact path thay vi predicate/first match.
- Fixture them `Small27042.png` o x2 lam decoy; test xoa `Small27041.png` va yeu
  cau `ASSET-PREFLIGHT-022`.

### RED/GREEN va mutation

- RED truoc code: test decoy x2 FAIL vi `DeployError` khong duoc raise.
- GREEN: `python -m unittest assetpack.tests.test_deploy -v` -> 14/14 PASS.
- Mutation: doi rieng small ve generic first-match; test decoy do dung ly do
  `DeployError not raised`. Da hoan nguyen exact path guard truoc khi test cuoi.

### Full suite va concern

- `python -m unittest discover -s assetpack/tests -v` -> 123/123 PASS trong
  164.920s.
- Contract nay gia dinh `closure.sprite_ids`, `closure.skill_ids`, va
  `closure.mob_ids` la danh sach ID int khong rong; lock production do build tao
  da co ba closure nay, lock thieu se fail closed `ASSET-PREFLIGHT-022`.

## Full verification

- `go vet ./...`: PASS.
- Full Python lượt đầu: RED vì đã nhầm MOVE legacy input; root `part.bin` bị
  thiếu. Đã restore đúng baseline SHA `09349f...`, giữ generated ở `assets/`.
- Full Go lượt đầu: chỉ RED oracle map literal 20; đã đổi intentional sang 21.
- Full Python lượt 2: **116/116 PASS trong 170.065s**.
- Sau khi thêm cổng SHA runtime key map và rollback marker/client nếu publish
  release lỗi, full Python source cuối: **118/118 PASS trong 170.751s**.
- Full Go lượt 2: **PASS**, `internal/game` 83.334s; các package còn lại xanh.

## Bảo toàn thay đổi người khác

Không sửa/stage/commit `NHAT-KY-AGENT.md`, `assetpack/bind.py` hay pycache.
