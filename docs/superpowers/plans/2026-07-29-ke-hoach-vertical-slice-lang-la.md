# Vertical slice Làng Lá — Kế hoạch thực thi

> **Cho người thực thi:** KỸ NĂNG BẮT BUỘC
> `superpowers:subagent-driven-development`, `superpowers:test-driven-development`
> và `superpowers:verification-before-completion`.
>
> **Đọc trước:** đặc tả
> [`../specs/2026-07-29-vertical-slice-lang-la.md`](../specs/2026-07-29-vertical-slice-lang-la.md),
> `AGENTS.md`, `BAT-DAU-TU-DAY.md`, `NHAT-KY-AGENT.md`, `CLAUDE.md` và
> `.claude/commands/viec.md`.

**Mục tiêu:** từ một client Unity sạch, hai tài khoản thật vào Làng Lá, nhìn
thấy nhân vật/NPC/Sói/map/HUD Naruto, di chuyển có collision, dùng đủ bốn skill,
đánh quái/PvP/chết/hồi sinh, rồi tắt–bật server vẫn giữ đúng vị trí và HP.

**Kiến trúc:** một manifest ID sinh ra cặp artifact server/client có cùng lock;
server giữ gameplay authoritative trong actor; Unity gốc chỉ được đọc rồi chép
thành staging có overlay; runner điều khiển hai executable qua telemetry
loopback và lưu evidence máy đọc được.

**Công nghệ:** Go 1.25 + pgx, Python 3 + Pillow, Unity 2022.3.62f2/C#,
PowerShell, PostgreSQL.

---

## Ràng buộc chung

- Không sửa các thư mục Java/PHP cũ. Không sửa trực tiếp `D:\NRO_Project`.
- Mỗi luồng dùng branch + worktree riêng. Không cho hai người cùng sửa
  `server/internal/game/` hoặc cùng sửa Unity staging.
- TDD nghiêm: test đỏ vì đúng hành vi thiếu → code tối thiểu → test xanh.
- Sau mỗi task phải tạo một mutation có chủ đích. Test phải đỏ đúng assertion/mã
  lỗi, không phải compile error; hoàn nguyên và chạy xanh lại.
- Mỗi commit implement phải qua hai cổng: spec review rồi quality review.
  Critical/Important phải sửa và re-review trước khi tích hợp.
- Asset/overlay sinh vào thư mục mới theo SHA, không ghi đè bản đang chạy.
- `mapping-id.json` là nguồn ID; `asset-lock.json` là bằng chứng build. Không
  hard-code lại ID/catalog ở Go, C# hoặc runner.
- Actor là nguồn trạng thái duy nhất. Không I/O, socket write hay DB write trên
  goroutine `KhuVuc.vongLap`.
- Outbound gameplay chỉ enqueue; mỗi session có đúng một writer.
- Test thời gian dùng đồng hồ tiêm được/event, không sleep đoán.
- PostgreSQL phải tự dò 5432/5433. Dùng `D:\Go\bin\go.exe`; không dùng
  `go test -race` trên máy này, bù bằng `-count=5` và test actor xác định.
- Mọi lệnh E2E chỉ dùng DB test, account prefix theo run ID và port loopback.
- Evidence/token/profile không chứa mật khẩu thật; token không ghi vào log.
- Ảnh AI phải qua alpha, crop, anchor, temporal animation và visual review.

---

## Worktree và các đợt song song

| Đợt | Luồng A | Luồng B | Luồng C | Coordinator |
|---|---|---|---|---|
| 0 | Task 0 khóa 8 runtime key | — | — | review + tích hợp |
| 1 | Task 1–3 manifest/assetpack | Task 4 actor outbound | Task 5 Unity staging/profile | Task 6A sinh/review ảnh thô |
| 2 | Task 6B bind ảnh vào manifest | Task 7 sau khi 6B xanh | — | review + tích hợp tuần tự |
| 3 | Task 8 server gameplay | Task 9 telemetry/runner | — | review + tích hợp |
| 4 | Task 10 clean/warm E2E | Task 11 canary sau E2E xanh | review độc lập | sửa đến xanh |
| 5 | Task 12 nghiệm thu | — | — | tài liệu + commit |

Worktree dự kiến:

```text
.worktrees/kh-assets
.worktrees/kh-actor
.worktrees/kh-unity-e2e
.worktrees/kh-vertical-slice
```

Mỗi nhánh bắt đầu từ cùng commit chứa spec + kế hoạch. Implementer chỉ commit
file thuộc vùng được giao. Coordinator cherry-pick sau review, không chép tay.
Chỉ luồng A sở hữu `mapping-id.json`. Task 6A chỉ sinh ảnh trong ArtRoot và
review metadata, không sửa manifest/không đòi validator đã tồn tại. Sau khi
Task 1–3 được tích hợp, luồng A làm Task 6B để bind hash/crop/anchor đã duyệt
vào manifest. Task 7 chỉ bắt đầu sau 6B và chạy Unity preflight thật sau deploy.
Task 8 bắt đầu sau khi Task 4 và Task 7 đã tích hợp; Task 8/9 được chạy song
song vì một bên chỉ sửa `server/internal/game`, bên kia sửa Unity/tools và
`server/cmd/nroserver`. Task 11 chỉ bắt đầu sau một lượt Task 10 xanh.

---

### Task 0: Khóa tám runtime resource key Unity

**Files:**

- Create `docs/my-thuat/unity-runtime-key-map.json`
- Create `tools/unity-overlay/discover-runtime-keys.ps1`
- Create `tools/unity-overlay/tests/Test-RuntimeKeys.ps1`

**Produces:** bảng tám dòng
`source_asset → Game1 call-site → logical resource key → physical xN path`,
kèm SHA của file C# chứa call-site.

- [ ] Viết self-test đỏ với fixture cố tình đổi một key/call-site SHA.
- [ ] Quét source Game1 thật cho logo, avatar/dialog frame, ba nút, ô chọn và ô
      vật phẩm; mỗi key phải có file+line/call expression làm bằng chứng.
- [ ] Test cấm placeholder và cấm suy key từ tên source art.
- [ ] Test map đủ đúng tám entry, không trùng logical key và có quy tắc x1..x4.
- [ ] Mutation: đổi một key sang tên “đoán”; test phải đỏ vì call-site không
      tham chiếu key đó.
- [ ] Commit: `docs(asset): khoa 8 runtime resource key Unity`

---

### Task 1: Khóa baseline và schema manifest

**Files:**

- Create `docs/my-thuat/mapping-id.json`
- Create `docs/my-thuat/asset-version-ledger.json`
- Create `assetpack/__init__.py`
- Create `assetpack/__main__.py`
- Create `assetpack/manifest.py`
- Create `assetpack/tests/test_manifest.py`
- Create `assetpack/tests/fixtures/manifest-minimal.json`
- Create `assetpack/tests/fixtures/legacy-baseline.json`
- Create `docs/my-thuat/legacy-input-ledger.json`

**Produces:**

- `python -m assetpack validate --manifest <path>`
- mã lỗi `ASSET-SPRITE-002`, `ASSET-ZOOM-003`
- model typed cho sprite, part, animation, map/tile/bg, mob, catalog và version

- [ ] Viết test đỏ cho: ID trùng; sprite thiếu x3; part trỏ sprite không tồn
      tại; `dx/dy` vượt sbyte; source hash sai; version runtime không tăng; crop
      Iruka thiếu overlap/layer; animation thiếu slot/cf.
- [ ] Viết scanner đọc baseline `part.bin`, mọi cache `NR_part`, `NR_image`,
      skill/mob legacy, Small local và bảng ID hard-code đã khóa.
- [ ] `legacy-input-ledger.json` ghi danh sách input path tương đối/loại/SHA,
      occupied IDs đã review và nguồn nào bắt buộc ở cổng rescan. Build thường
      đọc ledger versioned để tái lập; lệnh `verify-legacy-inputs` trên máy nguồn
      phải fail nếu thiếu input hoặc SHA lệch, tuyệt đối không lặng lẽ quét một
      tập cache khác. Fixture ghim số record/min/max/SHA và expected occupied
      set độc lập, không tái sinh expected bằng scanner đang test.
- [ ] Validator phải từ chối mọi va chạm với legacy và mọi ID vượt signed short.
- [ ] Chỉ thêm parser/validator đủ làm test xanh; không build ảnh ở task này.
- [ ] Mutation: bỏ kiểm tra ID trùng và bỏ kiểm tra x3; hai test phải đỏ đúng mã.
- [ ] Chạy:

```powershell
python -m unittest discover -s assetpack/tests -v
python -m assetpack validate --manifest docs/my-thuat/mapping-id.json
python -m assetpack verify-legacy-inputs --ledger docs/my-thuat/legacy-input-ledger.json
```

- [ ] Commit: `feat(asset): khoa manifest va baseline id Naruto`

---

### Task 2: Packer part.bin xác định và round-trip legacy

**Files:**

- Create `assetpack/part.py`
- Create `assetpack/tests/test_part.py`
- Create `assetpack/tests/fixtures/part-corrupt-*.bin`
- Modify `docs/my-thuat/mapping-id.json`

**Produces:**

- parser/encoder part big-endian
- `build_part(manifest, legacy_part) -> bytes`
- mã lỗi `ASSET-PART-001`

- [ ] Test đỏ ghim 2.189 record legacy, type 0/1/2 có 3/17/14 frame, tổng byte
      baseline và SHA.
- [ ] Test round-trip baseline phải byte-identical.
- [ ] Test Phase 1 chỉ thay part 30/57/58, nối 2189/2190/2191; mọi record còn
      lại byte-identical.
- [ ] Test corruption riêng cho endian, type, frame count và sprite ID thiếu.
- [ ] Cài encoder tối thiểu; không tự “sửa” dữ liệu legacy.
- [ ] Mutation: đảo endian một short và đổi frame count; test phải đỏ
      `ASSET-PART-001`.
- [ ] Commit: `feat(asset): sinh part bin round trip va append Iruka`

---

### Task 3: Packer ảnh, map, mob, catalog và lock

**Files:**

- Create `assetpack/images.py`
- Create `assetpack/animation.py`
- Create `assetpack/map.py`
- Create `assetpack/mob.py`
- Create `assetpack/skill.py`
- Create `assetpack/catalog.py`
- Create `assetpack/lock.py`
- Create `assetpack/visual.py`
- Create `assetpack/build.py`
- Create `assetpack/tests/test_images.py`
- Create `assetpack/tests/test_animation.py`
- Create `assetpack/tests/test_map.py`
- Create `assetpack/tests/test_mob.py`
- Create `assetpack/tests/test_skill.py`
- Create `assetpack/tests/test_lock.py`
- Create `assetpack/tests/test_visual.py`
- Create `assetpack/requirements.txt`

**Produces:**

- `python -m assetpack build --manifest ... --slice lang-la --out ...`
- server `part.bin`, `tile_set_info.bin`, `map-1-collision.bin`, `skill.bin`,
  `catalog.json`
- client Small, mymap, tile, bg, Mob, HUD/UI
- `asset-lock.json`

- [ ] Test đỏ cho 4 zoom: Small/crop/alpha; tile x1 `t_01..t_100`; tile
      x2–x4 `1$1..1$100`; đúng 24×N; nearest-neighbor; không thiếu key.
- [ ] Test action strip ghép xác định theo `slot_index` thành đúng sheet
      3/17/14/10 ô, đúng canvas/thứ tự/cell boundary; render golden trực tiếp
      các `Char.CharInfo.cf` đại diện. Mutation tráo hai slot phải đỏ dù tổng
      kích thước sheet không đổi.
- [ ] Test map: header 2 byte, đúng `2+w*h`, index `0..99`, spawn walkable;
      collision client/server cùng semantic hash.
- [ ] Test background đủ b00..b03 ở bốn zoom, đúng anchor/repeat.
- [ ] Test Mob `EffectData.readData` bằng golden byte độc lập và đúng frame
      idle/run/pounce/hit/death.
- [ ] Test `skill.bin` bằng golden byte độc lập theo parser Unity: đúng 12
      levelId, pose slot, endian/count, effect/arrow đều `-1`; mutation đổi
      count/endian hoặc ép effect 80 phải đỏ.
- [ ] Test `catalog.json` có 6 map name, Iruka, Sói, spawn và bốn pose skill.
- [ ] Dùng `unity-runtime-key-map.json` từ Task 0 để sinh đúng tám UI output;
      test mỗi key có bằng chứng call-site SHA và cấm placeholder/tên phỏng đoán.
- [ ] Test lock chứa mọi input/output SHA, dimension, ID/version và build hai
      lần khác thư mục cho cùng SHA.
- [ ] Build phải từ chối nếu review JSON thiếu PASS, source SHA/evidence SHA
      lệch hoặc evidence file thiếu. `asset-lock.json` phải chứa SHA của review
      JSON và từng contact sheet/animation capture/golden đã duyệt.
- [ ] Test crop Iruka composite trên nền sáng/tối và HUD 25/50/100%.
- [ ] `visual.py` bắt buộc sinh contact sheet tile 10×10 và full-map render
      x1/x2/x3/x4; test output tồn tại, hash vào review/lock và đúng logical
      viewport, không chỉ kiểm dimension từng tile.
- [ ] Mutation: đổi một solid tile nhưng không đổi collision; đỏ
      `ASSET-COLLISION-006`. Xóa b03; đỏ `ASSET-BG-005`.
- [ ] Commit: `feat(asset): build cap artifact server client co lock`

---

### Task 4: Outbound queue và control path không chặn actor

**Files:**

- Modify `server/internal/net/session.go`
- Modify `server/internal/net/session_test.go`
- Modify `server/internal/game/khuvuc.go`
- Modify `server/internal/game/khuvuc_test.go`
- Modify `server/internal/game/login.go`
- Modify `server/internal/game/login_test.go`
- Modify `server/internal/game/khuvuc_tichhop_test.go`

**Produces:**

- `(*Session).XepGui(*protocol.Message) error`
- `(*Session).DongCham()`
- `ErrHangGuiDay`
- control mailbox cho snapshot/leave/stop

- [ ] Test đỏ cho FIFO 256, immutable message, một writer, không xen byte và
      enqueue trả dưới 10ms khi conn giả chặn.
- [ ] Test writer gặp `Write` error: đóng idempotent, bỏ queue, goroutine thoát,
      `KhiDong`/cleanup đúng một lần.
- [ ] Test actor với client A không đọc: A bị loại, B vẫn đủ event, probe actor
      dưới 100ms, không ghost player.
- [ ] Test control snapshot/leave không rơi khi mailbox gameplay đầy; actor đã
      dừng trả lỗi rõ.
- [ ] Cài `Gui` đồng bộ chỉ cho handshake/login trước actor; thay mọi broadcast
      gameplay bằng `XepGui`.
- [ ] Chạy scan cấm `phien.Gui(` trong vùng actor/broadcast.
- [ ] Mutation: đổi enqueue thành write trực tiếp; test phải đỏ
      `SERVER-ACTOR-005`. Bỏ `DongCham` ở write error; leak test phải đỏ.
- [ ] Chạy hẹp `-count=5`, sau đó `go test ./... -count=1` và `go vet ./...`.
- [ ] Commit: `fix(server): tach writer socket khoi actor gameplay`

---

### Task 5: Snapshot Unity, overlay, profile và resource preflight

**Files:**

- Create `tools/unity-overlay/overlay-manifest.json`
- Create `tools/unity-overlay/apply-overlay.ps1`
- Create `tools/unity-overlay/verify-overlay.ps1`
- Create `tools/unity-overlay/Assets/Scripts/NroE2E/NroE2EConfig.cs`
- Create `tools/unity-overlay/Assets/Tests/EditMode/NroResourcePreflightTests.cs`
- Create `tools/unity-overlay/patches/rms-profile-root.patch`
- Create `tools/unity-overlay/patches/skill-no-effect-80.patch`
- Create `tools/unity-overlay/tests/Test-Overlay.ps1`

**Produces:**

- snapshot chỉ gồm Assets/Packages/ProjectSettings
- overlay có precondition SHA
- `-profile-root` chỉ khi `NRO_E2E`
- Unity batch preflight bốn zoom

- [ ] Self-test đỏ khi source SHA sai, overlay áp hai lần, hoặc đích là
      `D:\NRO_Project`.
- [ ] EditMode test đỏ khi thiếu bất kỳ tile x1–x4, sai dimension, Small ID
      thiếu, bg thiếu hoặc HUD key không load.
- [ ] EditMode test đọc `unity-runtime-key-map.json` của Task 0 và
      `Resources.Load` đủ tám key; không chấp nhận danh sách key viết lại trong
      C# test. Ở Task 5 chạy bằng fixture tối thiểu; lượt preflight artifact
      thật bắt buộc nằm ở Task 7 sau deploy.
- [ ] Test profile A/B sentinel không đọc chéo; production không có `NRO_E2E`
      giữ nguyên `persistentDataPath`.
- [ ] Patch skill Naruto giữ `-1=no effect`; test cấm load effect 80.
- [ ] Mutation: xóa x3 tile; Unity batch test phải đỏ `ASSET-ZOOM-003`.
- [ ] Commit: `feat(unity): staging overlay profile va resource preflight`

---

### Task 6A: Sinh và duyệt ảnh nguồn bằng Codex

**Files/data:**

- Create under the shared ignored art store
  `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\` (không đặt
  trong worktree vì file ignored không đi theo cherry-pick)
- Create `docs/my-thuat/review/vertical-slice-v1.json`

**Produces:**

- 20 icon skill, HP/chakra bg/fill, `hu_anthan`
- master reference Học viên và Sói
- action strips player/wolf, crop Iruka

- [ ] Dùng skill `imagegen`; mỗi asset/strip là một lần sinh riêng khi mục tiêu
      khác nhau, nhưng mọi strip dùng cùng master reference.
- [ ] Tách chroma bằng helper, không coi kết quả tự động là duyệt.
- [ ] Dùng helper độc lập tối thiểu kiểm alpha/canvas/cell boundary trước khi
      assetpack tồn tại; lưu report cạnh ảnh, không sửa `mapping-id.json`.
- [ ] Sinh contact sheet nguồn sáng/tối và animation preview đủ loop boundary.
- [ ] Reviewer độc lập ghi PASS/FAIL từng asset; FAIL thì sinh/sửa lại, không
      cho phép Task 6B bind.
- [ ] Review JSON ghi source SHA, evidence path/SHA, reviewer, timestamp và
      PASS/FAIL; file ảnh lớn không commit theo policy.
- [ ] Commit metadata: `art: review source vertical slice Lang La`

---

### Task 6B: Bind ảnh đã duyệt vào manifest

**Phụ thuộc:** Task 1–3 và Task 6A đã review/merge.

**Files:**

- Modify `docs/my-thuat/mapping-id.json`
- Modify `docs/my-thuat/review/vertical-slice-v1.json` chỉ để thêm evidence do
  assetpack sinh; không đổi PASS nguồn mà không có reviewer

- [ ] Luồng A duy nhất điền source SHA, crop/overlap/layer, anchor, action
      strip/slot/cf và tám runtime key đã khóa vào manifest.
- [ ] Chạy validator thật; build phải đỏ nếu review PASS/evidence/source SHA
      không khớp.
- [ ] Sinh sheet runtime 3/17/14/10, composite Iruka, HUD golden, animation
      capture và tile contact sheet/full-map x1–x4.
- [ ] Reviewer ký PASS artifact; cập nhật evidence SHA vào review JSON rồi
      validate lại.
- [ ] Mutation: tráo hai slot và thay source sau review; lần lượt phải đỏ ở
      golden cf và review SHA gate.
- [ ] Commit: `art: bind asset da duyet vao manifest Lang La`

---

### Task 7: Build artifact và deploy staging theo lock

**Phụ thuộc:** Task 0, Task 1–3, Task 5 và Task 6B đã review/merge vào cùng
integration branch.

**Files:**

- Create `assetpack/deploy.py`
- Create `assetpack/tests/test_deploy.py`
- Modify `server/internal/game/dongbo.go`
- Move generated runtime binaries to `server/internal/game/assets/`
- Modify `server/internal/game/dongbo_test.go`
- Modify `server/internal/game/dongbo_tichhop_test.go`

**Produces:**

- artifact immutable theo SHA
- deploy server/client chỉ khi cùng lock
- `vsData/vsMap/vsSkill` đúng ledger

- [ ] Test đỏ khi server/client lock lệch, version không bump hoặc deploy định
      ghi đè artifact đang chạy.
- [ ] Build hai lần, so toàn bộ SHA.
- [ ] Embed đúng `part.bin`, `skill.bin`, tile/collision/catalog; gói `-87`
      chứa blob đúng version/length/hash.
- [ ] Deploy client vào snapshot, không vào Unity gốc.
- [ ] Chạy Unity batch preflight thật trên snapshot đã deploy: đủ tám UI key,
      100 tile ở từng x1–x4, Small/bg/Mob/HUD/skill; lưu log vào artifact.
- [ ] Mutation: đổi một byte runtime nhưng giữ version; gate phải đỏ.
- [ ] Commit: `feat(runtime): nhung artifact Naruto va khoa cap server client`

---

### Task 8: Catalog, NPC, mob và collision server

**Files:**

- Create `server/internal/game/catalog.go`
- Create `server/internal/game/catalog_test.go`
- Create `server/internal/game/npc.go`
- Create `server/internal/game/npc_test.go`
- Modify `server/internal/game/banDo.go`
- Modify `server/internal/game/banDo_test.go`
- Modify `server/internal/game/quai.go`
- Modify `server/internal/game/quai_test.go`
- Modify `server/internal/game/dichuyen.go`
- Modify `server/internal/game/dichuyen_test.go`
- Modify `server/internal/game/dongbo.go`
- Modify `server/internal/game/dongbo_tichhop_test.go`
- Modify `server/internal/game/login.go`
- Modify `server/internal/game/khuvuc.go`

**Produces:**

- loader `catalog.json` + collision
- `-28/6` catalog thật; `-24` Iruka + Sói
- command 33 → 32 hội thoại Iruka
- server movement collision theo segment/foot point

- [ ] Test đỏ parser catalog với expected ghim độc lập; thiếu/unknown ID fail
      fast khi startup.
- [ ] Test protocol byte chính xác cho sáu map, Iruka, Sói và map info.
- [ ] Test NPC 33→32 có tên/text/menu/portrait đúng manifest.
- [ ] Test collision walkable/blocked/segment crossing khớp golden asset.
- [ ] Test fixture map 1 dùng head/body/leg 30/57/58, mob template 0.
- [ ] Mutation: đảo một solid cell server; đỏ `ASSET-COLLISION-006`.
- [ ] Mutation: bỏ portrait hoặc catalog mob; protocol/E2E contract test đỏ.
- [ ] Chạy Go hẹp `-count=5`, suite và vet.
- [ ] Commit: `feat(game): noi catalog NPC soi va collision Lang La`

---

### Task 9: Telemetry Unity, server control và runner

**Files:**

- Create `tools/unity-overlay/Assets/Scripts/NroE2E/NroE2EBridge.cs`
- Create `tools/unity-overlay/Assets/Scripts/NroE2E/NroE2EState.cs`
- Create `tools/unity-overlay/Assets/Scripts/NroE2E/NroE2ECommands.cs`
- Create `tools/unity-overlay/Assets/Tests/EditMode/NroE2EContractTests.cs`
- Create `tools/e2e/run-unity-e2e.ps1`
- Create `tools/e2e/New-Fixture.ps1`
- Create `tools/e2e/Test-Evidence.ps1`
- Create `server/cmd/nroserver/e2e_control.go`
- Create `server/cmd/nroserver/e2e_control_test.go`
- Modify `server/cmd/nroserver/main.go`
- Modify `server/cmd/nroserver/main_test.go`

**Produces:**

- loopback JSONL telemetry tokenized
- command/state schema trong spec
- E2E-only graceful shutdown gọi cùng path production

- [ ] Contract test đỏ cho auth token, command_id, sequence, main-thread
      dispatch, screenshot flush+SHA và state fields.
- [ ] Test server control không bật nếu thiếu env; chỉ bind loopback; token sai
      bị từ chối; shutdown gọi đúng `donDepKhiTat`.
- [ ] Fixture tạo/xóa đúng hai account prefix run ID và ghim stats/parts/skills.
- [ ] Runner preflight source hash, DB, port, lock, Unity version, profile.
- [ ] Runner có timeout từng phase; không dùng `Stop-Process` làm graceful.
- [ ] Mutation: control endpoint gọi exit trực tiếp; persistence test phải đỏ.
- [ ] Commit: `test(e2e): telemetry hai client va graceful control`

---

### Task 10: Clean full E2E và warm-upgrade smoke

**Files/output:**

- Modify `tools/e2e/run-unity-e2e.ps1`
- Modify `tools/e2e/Test-Evidence.ps1`
- Create `tools/e2e/tests/fixtures/*.ndjson`
- Create `tools/e2e/tests/fixtures/*-result.json`
- Create `tools/e2e/tests/Test-RunnerOracle.ps1`
- Create ignored `e2e-runs/<run-id>/...`

- [ ] Viết test oracle đỏ trước bằng fixture: damage thừa; HP cuối 139 thay vì
      137; clean/warm lệch executable hoặc lock; warm gắn nhầm `scope=full`;
      thiếu DB/evidence SHA; sequence NDJSON lùi/đứt.
- [ ] Cài parser/aggregate tối thiểu để từng fixture đỏ đúng mã và fixture chuẩn
      xanh; runner thật phải gọi chính parser đã test, không viết lại assertion.
- [ ] Unity cold snapshot/import/build; chạy resource preflight x1–x4.
- [ ] Clean full: hai client login, thấy nhau, Iruka, collision/parallax, đủ
      bốn skill, quái, PvP, chết/hồi sinh, HUD.
- [ ] Đưa A về 137 bằng đúng một đòn Sói 5 + bốn đòn PvP 2; mọi event/damage
      phải khớp.
- [ ] Graceful shutdown, DB snapshot, restart/relogin, so tuyệt đối map/x/y/HP.
- [ ] Warm smoke: seed cache cũ, xác nhận negotiation/update/lock mới, vào map
      bằng resource mới và không fallback Dragon Ball.
- [ ] `aggregate-result.json` chỉ PASS khi clean full + warm smoke cùng
      executable/server lock.
- [ ] Reviewer xem toàn bộ screenshots/animation capture, contact sheet tile
      10×10 và full-map render x1/x2/x3/x4; checklist ghi rõ từng evidence SHA
      và PASS/FAIL cho blur, seam, edge alignment, anchor, loop và impact timing.
- [ ] Mutation runner: bỏ kiểm damage thừa, bỏ so HP 137 và bỏ lock equality;
      ba fixture xấu phải lần lượt đỏ vì test oracle phát hiện mutation. Hoàn
      nguyên rồi chạy fixture suite xanh trước full E2E.

---

### Task 11: Canary/mutation wave

**Files/output:**

- Create `tools/e2e/run-canaries.ps1`
- Create ignored `e2e-runs/<run-id>/canaries/...`

- [ ] Từng mutation chạy trên artifact/worktree/snapshot riêng, ghi hash trước.
- [ ] Chạy đủ:
      `ASSET-PART-001`, `ASSET-SPRITE-002`, `ASSET-ZOOM-003`,
      `ASSET-MAP-004`, `ASSET-BG-005`, `ASSET-COLLISION-006`,
      `E2E-CACHE-002`, `E2E-BROADCAST-003`, `E2E-PERSIST-004`,
      `SERVER-ACTOR-005`.
- [ ] Mỗi canary phải đỏ đúng mã và expected/actual; compile error không tính.
- [ ] Hoàn nguyên, hash phải trở lại; chạy clean/warm xanh lại.
- [ ] Canary nào còn xanh: đánh dấu test vô dụng, bổ sung oracle rồi lặp lại.
- [ ] Commit: `test: canary matrix vertical slice Lang La`

---

### Task 12: Nghiệm thu, sổ dự án và bàn giao

**Files:**

- Modify `BAT-DAU-TU-DAY.md`
- Modify `NHAT-KY-AGENT.md`
- Modify `README.md` hoặc `server/README.md` nếu lệnh chạy thay đổi
- Add `docs/superpowers/reports/2026-07-29-nghiem-thu-vertical-slice-lang-la.md`

- [ ] Chạy cuối:

```powershell
cd server
$env:NRO_TEST_DSN = "postgres://postgres:nrodev123@127.0.0.1:<CONG>/nro_test"
D:\Go\bin\go.exe test ./... -count=1
D:\Go\bin\go.exe test ./internal/game ./internal/net -count=5
D:\Go\bin\go.exe vet ./...

cd ..\web
$env:DATABASE_URL = "postgres://postgres:nrodev123@127.0.0.1:<CONG>/nro"
$env:AUTH_SECRET = "<ephemeral-test-secret>"
npm run build
```

- [ ] Chạy clean full + warm smoke sau canary revert.
- [ ] Đối chiếu evidence SHA, DB, screenshots, visual checklist.
- [ ] Reviewer cuối không còn Critical/Important.
- [ ] Cập nhật §4/§5/§8/§12 và §11 trong `BAT-DAU-TU-DAY.md`; thêm dòng sổ
      `NHAT-KY-AGENT.md`.
- [ ] Commit: `feat: hoan thanh vertical slice gameplay Lang La`

---

## Cổng hoàn thành

Không được báo “xong” nếu thiếu bất kỳ điều nào:

- mọi test Go/Python/Unity/Web xanh ở lệnh cuối;
- hai client thật hoàn thành clean full và warm smoke;
- asset server/client cùng lock, không fallback Dragon Ball;
- actor không bị client chậm chặn và không ghost/leak;
- persistence giữ đúng 137/150 + vị trí qua graceful restart;
- toàn bộ canary đỏ đúng lý do rồi xanh sau hoàn nguyên;
- visual reviewer PASS animation/map/HUD/NPC/player/wolf;
- reviewer code cuối sạch Critical/Important;
- sổ dự án và hướng dẫn chạy đã cập nhật trong cùng commit.
