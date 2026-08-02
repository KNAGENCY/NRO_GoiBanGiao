# Bàn giao phiên Codex — vertical slice Làng Lá

**Ngày chốt:** 2026-07-29, múi giờ Asia/Bangkok  
**Mã cuộc trò chuyện cũ:** `019fa9a3-b3bf-7551-bdd3-a315d68617a6`  
**Mục tiêu chưa hoàn thành:** làm được gameplay vertical slice Làng Lá thật,
review độc lập, sửa theo review, mutation test/canary và E2E hai client liên tục
cho tới khi đủ bằng chứng nghiệm thu.

> Đây là điểm tiếp tục, không phải báo cáo nghiệm thu. Không được tuyên bố gameplay
> đã xong chỉ vì các nhánh con có test xanh.

## 1. Việc phiên mới phải đọc trước

Đọc theo đúng thứ tự:

1. `D:\NRO_GoiBanGiao\AGENTS.md`
2. `D:\NRO_GoiBanGiao\BAT-DAU-TU-DAY.md`
3. `D:\NRO_GoiBanGiao\NHAT-KY-AGENT.md`
4. `D:\NRO_GoiBanGiao\CLAUDE.md`
5. `D:\NRO_GoiBanGiao\.claude\commands\viec.md`
6. File bàn giao này
7. Spec:
   `docs/superpowers/specs/2026-07-29-vertical-slice-lang-la.md`
8. Kế hoạch:
   `docs/superpowers/plans/2026-07-29-ke-hoach-vertical-slice-lang-la.md`

Sau đó chạy:

```powershell
git -C D:\NRO_GoiBanGiao worktree list
git -C D:\NRO_GoiBanGiao status --short --branch
git -C D:\NRO_GoiBanGiao\.worktrees\kh-vertical-slice status --short --branch
git -C D:\NRO_GoiBanGiao\.worktrees\kh-assets status --short --branch
git -C D:\NRO_GoiBanGiao\.worktrees\kh-actor status --short --branch
git -C D:\NRO_GoiBanGiao\.worktrees\kh-unity-e2e status --short --branch
```

## 2. Trạng thái Git chính xác lúc bàn giao

### Thư mục chính — không được dọn hay ghi đè

Đường dẫn: `D:\NRO_GoiBanGiao`  
Nhánh: `kh-mapping-anh`  
HEAD: `4c775ad`

Đang có POC cũ của người dùng:

```text
M  server/internal/game/dongbo.go
M  server/internal/game/dongbo_test.go
M  server/internal/game/part.bin
?? server/server-log.txt
```

POC này dùng sprite thử `9000..9002` và `vsData21`, chưa phải asset chính thức.
Không reset, checkout, xóa hoặc trộn nó vào vertical slice.

### Worktree tích hợp

Đường dẫn: `D:\NRO_GoiBanGiao\.worktrees\kh-vertical-slice`  
Nhánh: `kh-vertical-slice`  
HEAD: `62a0943 docs: chot vertical slice Lang La va E2E Unity`

Hai file đang sửa nhưng chưa commit:

```text
M docs/superpowers/plans/2026-07-29-ke-hoach-vertical-slice-lang-la.md
M docs/superpowers/specs/2026-07-29-vertical-slice-lang-la.md
```

Sửa đổi này là cần thiết: khảo sát Game1 cho thấy tám UI không phải tám key 1:1.
Contract mới có ba mode:

- `direct`: consumer nạp key thật;
- `slice_group`: consumer ghép nhiều key/slice;
- `overlay_required`: consumer hiện vẽ primitive/composite, phải thêm
  `future_logical_key` và SHA precondition.

Chưa commit vì lượt re-review độc lập bị ngắt khi đổi phiên. Phiên mới phải
review diff này cùng Task 0; CLEAN rồi mới commit.

### Nhánh asset/runtime key

Đường dẫn: `D:\NRO_GoiBanGiao\.worktrees\kh-assets`  
Nhánh: `kh-assets`  
Worktree sạch.

Commit:

```text
c854b34 fix(asset): khoa mode va physical path runtime UI
b8efa5a docs(asset): khoa 8 runtime resource key Unity
```

Task 0 đã làm:

- discovery tám binding Game1;
- tỷ lệ thật `1 direct / 2 slice_group / 5 overlay_required`;
- fixture có cả nút thường/nút nhấn, mỗi nút đủ left/mid/right;
- khóa mode theo source và scope case-sensitive `Assets/Scripts/Game1`;
- canonical path, cấm absolute/`..` và cấm thoát `Assets/Resources/res/xN`;
- `-RequireFutureOutputs` bắt đủ năm output overlay tương lai;
- mutation đổi mode và path traversal đều đỏ;
- self-test và discovery trên `D:\NRO_Project` đã PASS;
- Go test toàn bộ và `go vet` đã PASS ở lần agent chạy.

Review đầu tìm thấy ba vấn đề trên và `c854b34` đã sửa. Re-review sau sửa chưa
kịp trả kết quả vì phiên bị ngắt. Không gọi Task 0 là CLEAN trước khi re-review.
Task 1–3 assetpack chưa bắt đầu.

### Nhánh actor outbound

Đường dẫn: `D:\NRO_GoiBanGiao\.worktrees\kh-actor`  
Nhánh: `kh-actor`  
Worktree sạch.

Commit:

```text
2a95a92 fix(server): tach writer socket khoi actor gameplay
```

Agent báo đã làm:

- FIFO outbound dung lượng 256 và chỉ một socket writer;
- message immutable;
- `XepGui` không block actor gameplay;
- write error đi đúng lifecycle;
- đuổi client chậm;
- control mailbox và snapshot timeout `ErrControlQuaHan`;
- mutation/canary liên quan;
- full Go test và vet xanh.

Đã có một lần pet E2E timeout lẻ; agent chạy lại test cô lập, prefix-order nhiều
lần và full suite ba lần thì xanh. Một reviewer độc lập đang soi commit nhưng
chưa trả kết quả khi phiên bị ngắt. Phải review lại từ đầu, chú ý đặc biệt:

- quan hệ giữa outbound queue, disconnect và pet loop;
- FIFO khi client chậm;
- không có hai goroutine cùng ghi socket;
- snapshot timeout không làm actor treo;
- test có thật sự đỏ khi cố tình phá các invariant trên.

Không cherry-pick `2a95a92` trước khi reviewer báo CLEAN hoặc mọi finding đã sửa
và được re-review.

### Nhánh Unity staging/E2E

Đường dẫn: `D:\NRO_GoiBanGiao\.worktrees\kh-unity-e2e`  
Nhánh: `kh-unity-e2e`  
HEAD vẫn là `62a0943`; toàn bộ việc đang **chưa commit**.

Các file đang có:

```text
M BAT-DAU-TU-DAY.md
M NHAT-KY-AGENT.md
A tools/unity-overlay/Assets/Editor/NroResourceImportSettings.cs
A tools/unity-overlay/Assets/Scripts/NroE2E/NroE2EConfig.cs
A tools/unity-overlay/Assets/Scripts/NroE2E/NroUiOverlay.cs
A tools/unity-overlay/Assets/Tests/EditMode/NroResourcePreflightTests.cs
A tools/unity-overlay/apply-overlay.ps1
A tools/unity-overlay/overlay-manifest.json
A tools/unity-overlay/patches/rms-profile-root.patch
A tools/unity-overlay/patches/skill-no-effect-80.patch
A tools/unity-overlay/patches/skill-no-effect-render.patch
A tools/unity-overlay/patches/ui-chat-popup.patch
A tools/unity-overlay/patches/ui-command-disabled.patch
A tools/unity-overlay/patches/ui-panel-slots.patch
A tools/unity-overlay/tests/Test-Overlay.ps1
A tools/unity-overlay/tests/Test-ResourcePreflight.ps1
A tools/unity-overlay/tests/Test-RuntimeBindingOverlay.ps1
A tools/unity-overlay/verify-overlay.ps1
A tools/unity-overlay/verify-runtime-bindings.ps1
```

Agent đang làm Task 5: overlay có precondition SHA, profile A/B tách RMS,
resource import/preflight, giữ `-1 = no effect`, fixture test cho ba loại binding.
Phiên mới phải:

1. đọc diff và chạy test hiện có;
2. không tự ý bỏ thay đổi chưa commit;
3. yêu cầu agent cũ/agent mới hoàn thành TDD và mutation;
4. review độc lập;
5. sửa finding và re-review;
6. chỉ commit khi worktree sạch và bằng chứng test đầy đủ.

## 3. Ảnh nguồn đã tạo bằng Codex

Ảnh nằm ở vùng dùng chung bị Git ignore:

```text
D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1
```

Có 20 file nguồn:

```text
hu_anthan_10f_chroma.png
hud_hp_chakra_4_chroma.png
iruka_portrait_chroma.png
master_hoc_vien_x4_chroma.png
master_iruka_x4_chroma.png
master_soi_rung_x4_chroma.png
player_cast_3f_chroma.png
player_hit_death_getup_3f_chroma.png
player_idle_3f_chroma.png
player_melee_4f_chroma.png
player_run_4f_chroma.png
skill_icons_cloud_4_chroma.png
skill_icons_leaf_4_chroma.png
skill_icons_mist_4_chroma.png
skill_icons_sand_4_chroma.png
skill_icons_stone_4_chroma.png
wolf_hit_death_2f_chroma.png
wolf_idle_2f_chroma.png
wolf_pounce_2f_chroma.png
wolf_run_4f_chroma.png
```

Nội dung đã có:

- 20 icon kỹ năng, mỗi làng một strip bốn icon;
- bốn lớp HUD HP/chakra trong một sheet;
- hiệu ứng `hu_anthan` 10 frame;
- master Học viên, Iruka và Sói;
- player 17 pose tổng cộng;
- Sói 10 pose tổng cộng;
- portrait Iruka.

### Cảnh báo QA ảnh — rất quan trọng

Ảnh nhìn đẹp hơn placeholder cũ, nhưng **chưa PASS pipeline**:

- imagegen không tạo nền `#FF00FF` phẳng tuyệt đối; góc ảnh là nhiều sắc magenta;
- `hud_hp_chakra_4_chroma.png` có nền gần trắng, không phải chroma;
- một số strip có chiều rộng không chia hết cho số ô, ví dụ icon `1983/4`;
- phải tách nền bằng helper có flood/alpha an toàn, chuẩn hóa canvas và ranh giới;
- phải xem x1 trên nền sáng lẫn tối, kiểm halo, anchor chân và loop;
- phải sinh contact sheet/animation preview;
- phải có reviewer độc lập PASS/FAIL từng asset;
- chỉ ảnh PASS mới được Task 6B bind vào manifest.

Kiểm tra ban đầu cho thấy các PNG vẫn là RGB, không có alpha đã tách. Đừng dùng
thẳng vào assetpack chỉ vì tên file có chữ `chroma`.

Ba SHA master đã ghi:

```text
master_hoc_vien_x4_chroma.png
ea3e32732485861fe323e8c191543a42dc6cbab74ea56b0f9c851228368d47cd

master_iruka_x4_chroma.png
a53d149ca17557f365dd3d2b8f9f1b9d4fcb6c14dbfb130876157e7e27533622

master_soi_rung_x4_chroma.png
c8002663b8c1cd2ccbe42ba0ea4d2fbc1ed54dd4ec1c001c5f22f6c341d64827
```

Hai ảnh cuối phiên:

```text
hu_anthan_10f_chroma.png
ba1280259b8818d99ce53648e3b5fe3583ccbd10abc4ed3cf94f48bbfc55ecba

iruka_portrait_chroma.png
f6b84bde6409bd3e8de1a82f25a3329d135453bafdba00a4297cd5877b897251
```

Task 6A còn thiếu:

- helper + test đỏ trước;
- alpha/canvas/cell report;
- contact sheet sáng/tối và preview;
- review JSON `docs/my-thuat/review/vertical-slice-v1.json`;
- reviewer độc lập và vòng sửa;
- commit metadata `art: review source vertical slice Lang La`.

## 4. Spec và kế hoạch đã được chốt

Commit gốc:

```text
62a0943 docs: chot vertical slice Lang La va E2E Unity
```

Phạm vi Phase 1:

- Học viên có idle/run/melee/cast/hit/death/get-up;
- NPC Iruka có composite và hội thoại;
- Sói có idle/run/pounce/hit/death;
- map Làng Lá, collision client/server và parallax;
- bốn skill Làng Lá + HUD;
- hai client cùng thấy nhau;
- restart vẫn giữ trạng thái;
- actor không bị block bởi socket chậm.

ID sprite dự trữ: `27041..27999`; Phase 1 dùng `27041..27082`.

Tile native Game1:

```text
x1: res/x1/t/1/t_01.png ... t_100.png
x2..x4: res/xN/t/1$1.png ... 1$100.png
```

Kích thước tile là `24*N`, lọc nearest. Unity preflight phải kiểm x1..x4.

E2E phải có cả clean full run và warm smoke. Checkpoint damage chính xác:

```text
HP ban đầu 150
Sói đánh một lần: 150 -> 145
Bốn hit PvP, mỗi hit 2: 145 -> 143 -> 141 -> 139 -> 137
```

Phải có event `skill_ready`, damage event, safe point sau Sói, restart
persistence và bằng chứng hai client.

## 5. Thứ tự tiếp tục đề nghị

Không được nhảy thẳng tới E2E giả lập khi asset/gameplay chưa tích hợp.

1. Re-review Task 0 (`b8efa5a`, `c854b34`) cùng diff contract đang dở.
2. CLEAN thì commit contract ở nhánh tích hợp.
3. Cho nhánh `kh-assets` làm Task 1–3 theo TDD; review/sửa/re-review.
4. Re-review actor `2a95a92`; finding nào có thì trả lại sửa và review tiếp.
5. Hoàn thành/commit/review Task 5 Unity từ worktree đang dở.
6. Song song hoàn tất Task 6A ảnh: helper, QA, contact sheet, review JSON.
7. Task 6B bind chỉ những ảnh PASS vào manifest.
8. Cherry-pick theo thứ tự phụ thuộc vào `kh-vertical-slice`, giải quyết log
   `BAT-DAU-TU-DAY.md`/`NHAT-KY-AGENT.md` cẩn thận, không lấy bên này đè bên kia.
9. Task 7 build/embed/deploy cặp artifact server/client có lock khớp.
10. Task 8 gameplay server.
11. Task 9 telemetry/control/runner.
12. Task 10 E2E hai client, clean + warm.
13. Task 11 thả toàn bộ canary/mutation đã ghi trong spec; từng canary phải làm
    đúng test đỏ rồi revert, chạy lại xanh.
14. Task 12 chạy acceptance đầy đủ, review cuối và cập nhật hai sổ dự án.

Mỗi phần đều theo vòng:

```text
test đỏ -> code -> test xanh -> mutation đỏ -> revert -> test xanh
-> review độc lập -> sửa -> re-review
```

## 6. Lệnh kiểm chứng nền

Tự dò cổng PostgreSQL, không ghim số:

```powershell
Get-NetTCPConnection -State Listen |
  Where-Object { $_.LocalPort -in 5432,5433 }
```

Server, thay `CONG` bằng cổng tìm được:

```powershell
cd D:\NRO_GoiBanGiao\.worktrees\kh-vertical-slice\server
$env:NRO_TEST_DSN = "postgres://postgres:nrodev123@127.0.0.1:CONG/nro_test"
D:\Go\bin\go.exe test ./... -count=1
D:\Go\bin\go.exe test ./... -count=5
D:\Go\bin\go.exe vet ./...
```

Web:

```powershell
cd D:\NRO_GoiBanGiao\.worktrees\kh-vertical-slice\web
$env:DATABASE_URL = "postgres://postgres:nrodev123@127.0.0.1:CONG/nro_test"
$env:AUTH_SECRET = "nro-e2e-local-only"
npm run build
```

Máy này không chạy được `go test -race`; bù bằng `-count=5` và review actor.

## 7. Prompt để dán nguyên văn vào phiên Codex mới

```text
Tiếp quản vertical slice Làng Lá tại D:\NRO_GoiBanGiao.

Trước khi làm gì, hãy đọc D:\NRO_GoiBanGiao\AGENTS.md rồi đọc đúng thứ tự các
file mà AGENTS.md yêu cầu. Sau đó đọc toàn bộ:
D:\NRO_GoiBanGiao\.worktrees\kh-vertical-slice\docs\ban-giao\2026-07-29-codex-vertical-slice.md

Đây là phiên tiếp tục, không lập kế hoạch lại từ đầu. Hãy kiểm tra trạng thái
thật của cả bốn worktree. Giữ nguyên POC bẩn ở thư mục chính. Tiếp tục dùng
subagent song song theo đúng vùng sở hữu; phần nào code xong phải review độc lập,
sửa finding rồi re-review. Luôn TDD, mutation/canary và E2E; không tuyên bố hoàn
thành nếu chưa có bằng chứng gameplay hai client, restart persistence và toàn bộ
acceptance trong spec.

Việc đầu tiên:
1) re-review Task 0 cùng diff contract chưa commit;
2) re-review commit actor 2a95a92;
3) tiếp tục worktree Unity đang có thay đổi chưa commit;
4) sau khi Task 0 CLEAN, cho nhánh asset bắt đầu Task 1–3;
5) hoàn thành QA ảnh nguồn, tuyệt đối không coi nền chroma hiện tại là đã PASS.

Hãy làm liên tục tới khi đạt mục tiêu hoặc gặp blocker thật cần quyền/người dùng;
gửi cập nhật ngắn bằng tiếng Việt trong lúc làm.
```

## 8. Điều tuyệt đối không làm

- Không reset/xóa thay đổi ở thư mục chính.
- Không sửa code Java/PHP cũ; chỉ đọc để tra giao thức.
- Không sửa trực tiếp `D:\NRO_Project`; chỉ quét/copy snapshot được phép.
- Không cherry-pick nhánh chưa review CLEAN.
- Không coi ảnh đã sinh là ảnh runtime đã duyệt.
- Không coi test xanh một lần là đủ; phải có mutation và lặp.
- Không tự ghim cổng PostgreSQL.
- Không báo “xong gameplay” nếu chưa có E2E hai client thật và evidence.
