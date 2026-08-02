# Task 10 — audit so sánh hai bộ oracle

Ngày audit: 31/07/2026

- Bộ A: `e6582ce11c768322f171e33ad0558b2ffe614545`
- Bộ B: `36846262febdc860544914597b4168b1d5887698`
- Base chung: `01af875cfd5b3c6403c5b95838bfb94b5346b07c`
- Worktree audit: `kh-task10-oracle` tại commit B

## Kết luận ngắn

**Chọn A làm canonical baseline**, vì A đặt quyết định PASS trong
`Test-Evidence.ps1`, runner thật gọi đúng parser đó, và test oracle đi xuyên qua
`run-unity-e2e.ps1 -OracleOnly`. Đây là ranh giới production đúng.

Không nên coi A nguyên trạng là oracle hoàn thiện cuối cùng. Trước khi dùng làm
canonical acceptance, cần **port thủ công** các guard độc lập của B: asset lock
tách khỏi server lock, clean PASS/warm FAIL, artifact/evidence rỗng, SHA sai,
file evidence zero-byte, Sói phải đúng source ID 1 và HP tối đa phải đúng 150.

**Không cherry-pick nguyên commit B vào A.** Hai commit cùng thêm
`tools/e2e/tests/Test-RunnerOracle.ps1`; `git merge-tree` báo `added in both`.
Quan trọng hơn, schema B là một JSON chứa inline aggregate/clean/warm, trong khi
A dùng cây hồ sơ thật `clean/result.json` + `warm-upgrade/result.json` và sinh
`aggregate-result.json` sau khi validate. Giải quyết xung đột bằng cách chọn một
file sẽ làm mất một trong hai mục tiêu: boundary thật của A hoặc guard bổ sung
của B.

Đây chỉ là audit oracle unit-level. Không có gameplay E2E nào được chạy hoặc
được tuyên bố PASS.

## 1. Diff và file overlap

Hai commit cùng rẽ trực tiếp từ base `01af875`.

| | Bộ A | Bộ B |
|---|---:|---:|
| File thay đổi/tạo | 22 | 26 |
| Dòng thêm | 761 | 803 |
| File overlap | `tools/e2e/tests/Test-RunnerOracle.ps1` | cùng file |

A còn sửa production boundary:

- `tools/e2e/Test-Evidence.ps1`;
- `tools/e2e/run-unity-e2e.ps1`;
- `BAT-DAU-TU-DAY.md`, `NHAT-KY-AGENT.md`;
- 17 file fixture/oracle riêng.

B không chạm runner/implementation, chỉ thêm:

- parser kiêm self-test độc lập;
- 14 fixture JSON và 5 NDJSON cùng evidence tĩnh;
- report Task 10 riêng.

Ngoài file oracle cùng tên, tên fixture không overlap. Vì vậy rủi ro conflict
Git nhìn bề ngoài chỉ có một file, nhưng rủi ro **semantic cao**: toàn bộ fixture
B phụ thuộc schema/parser B và không tự kiểm được parser production của A.

## 2. Contract và schema

### A — schema hồ sơ run thật

A đọc trực tiếp cấu trúc dự kiến của runner:

```text
<run-root>/
  clean/result.json
  clean/events.ndjson
  clean/db-before.json
  clean/db-after-shutdown.json
  clean/db-after-relogin.json
  clean/screenshots/...
  warm-upgrade/result.json
  warm-upgrade/events.ndjson
  warm-upgrade/cache-before.json
  warm-upgrade/cache-after.json
  warm-upgrade/screenshots/...
```

`result.json` dùng `status="pass"` chữ thường, có trường riêng cho combat,
persistence, database, resources và visual review. Sau khi clean và warm đều
hợp lệ, `Test-Evidence.ps1` mới sinh `aggregate-result.json`.

Điểm mạnh:

- aggregate không phải dữ liệu tự khai để parser tin; nó là output sau validate;
- DB snapshot được parse và so trạng thái qua shutdown/relogin;
- warm kiểm negotiation, map 1, applied lock và cấm Dragon Ball fallback;
- bắt buộc screenshot và visual checklist;
- sequence được theo dõi **theo từng client**, phù hợp hai luồng telemetry A/B.

### B — schema fixture gói độc lập

B đọc một JSON:

```text
schema_version
aggregate { status, executable_sha256, server_lock_sha256, asset_lock_sha256 }
clean     { status, scope, hashes, final_state, events_path, artifacts, evidence }
warm      { status, scope, hashes, events_path, artifacts, evidence }
```

`status="PASS"` chữ hoa. B kiểm hash/file thật tốt, nhưng không dùng schema
runner A và không được `run-unity-e2e.ps1` gọi. Nó cũng không parse nội dung DB,
không kiểm persistence tuyệt đối, visual review hoặc warm negotiation/fallback.

Kết luận schema: **A là schema canonical; B là nguồn test pressure**, không phải
schema thứ hai cần giữ song song.

## 3. So sánh từng hợp đồng

| Hợp đồng | A | B | Audit |
|---|---|---|---|
| Damage | Đúng 5 event `5,2,2,2,2`, đúng chuỗi HP và source kind, lọc theo target client | Đúng chuỗi, bắt thêm source ID Sói = 1, nhưng gom mọi damage của file | Giữ lọc theo client của A; port check source ID 1 từ B |
| HP cuối | So telemetry, result, persistence và DB relogin đều 137 | Bắt `final_state=137/150` | A mạnh hơn về persistence; cần port absolute `hp_max=150` |
| Clean/warm scope | `full` và `warm_smoke` | `full` và `warm_smoke` | Cùng đúng; A dùng mã chung `E2E-SCOPE-014` |
| Clean + warm cùng PASS | Mỗi phase bị reject trước khi aggregate được sinh | Fixture riêng clean PASS/warm FAIL | Port fixture B để pressure-test invariant này |
| Executable | So clean/warm | So aggregate/clean/warm | A đủ cho output do chính nó sinh; B hữu ích nếu aggregate là input |
| Server lock | Có một `server_lock_sha256`, so clean/warm | Có field và fixture riêng | Giữ check A, thêm mã riêng thay vì gộp với asset |
| Asset lock | Không có field riêng; dùng `server_lock_sha256` cho cả applied resource lock | Có `asset_lock_sha256` và fixture mismatch riêng | Gap quan trọng của A; bắt buộc port |
| Evidence SHA | Mã riêng cho DB và evidence; parse DB | Hash lại mọi entry, có fixture missing/mismatch | Giữ taxonomy A; port wrong-SHA/empty/zero-byte từ B |
| Artifact rỗng | Runtime A reject mảng rỗng, nhưng oracle suite không có fixture riêng | Có fixture `artifacts-empty.json` | Port fixture theo schema A |
| Evidence/file rỗng | NDJSON không event bị reject; helper hash không chặn file zero-byte có SHA hợp lệ | `Resolve-FixtureFile` chặn file zero-byte; có fixture evidence array rỗng | Port length check và fixture zero-byte |
| Sequence NDJSON | Tăng liên tục theo từng client | Tăng liên tục toàn file | Giữ A; thuật toán B sẽ reject sai stream A1/B1 hợp lệ |
| Event fields | Task 10 reader A chỉ bắt client/sequence; chưa ghim đủ monotonic/utc/type/payload | B bắt sequence/type/payload, chưa bắt đủ monotonic/utc | Canonical cần dùng đủ contract Task 9 cho từng dòng |
| Visual review | Có | Không | Giữ A |
| Warm negotiation/fallback | Có | Không | Giữ A |
| Runner boundary | Test gọi runner `-OracleOnly`, runner gọi production `Test-Evidence.ps1` | Parser/self-test đứng riêng, runner không gọi | Đây là lý do chính chọn A |

## 4. Findings cần xử lý trước canonical acceptance

### Important 1 — A chưa phân biệt server lock và asset lock

Trong A, clean/warm chỉ có `server_lock_sha256`; warm so
`resources.applied_lock_sha256` với chính trường này. Mã `E2E-ASSET-013` đang
được dùng cho cả “server lock thiếu/sai” và “asset”. Do đó A không chứng minh
được hai khóa độc lập cùng đúng.

Khuyến nghị:

- thêm `asset_lock_sha256` vào clean, warm và aggregate;
- `server_lock_sha256` chỉ đại diện server/source lock;
- `resources.applied_lock_sha256` phải bằng `asset_lock_sha256`;
- so equality clean/warm cho từng loại;
- thêm hai fixture mismatch tách biệt.

### Important 2 — top-level manifest của A ghi sai scope full run

`run-unity-e2e.ps1` đặt mọi run không phải control smoke thành
`scope="task9_contract"`. Sau đó full Task 10 có thể đổi status thành
`full_pass` nhưng không đổi scope. Clean/warm result và aggregate có scope đúng,
song manifest ngoài cùng vẫn gây diễn giải sai.

Khuyến nghị: full path phải ghi scope Task 10 rõ ràng; `ContractOnly` mới dùng
`task9_contract`.

### Important 3 — A chưa fail-closed với evidence zero-byte

`Assert-E2EArtifactSha` kiểm tồn tại và SHA nhưng không kiểm `Length > 0`.
Screenshot/cache rỗng có SHA của file rỗng vẫn có thể qua. B có guard này.

Khuyến nghị: chặn zero-byte trước khi hash cho DB, NDJSON, cache và artifact;
thêm fixture riêng, không chỉ kiểm mảng `[]`.

### Important 4 — assertion combat/HP của A còn thiếu hai literal trong spec

A không kiểm `source_id=1` cho damage Sói đầu tiên, và không ghim tuyệt đối
`hp_max=150`; DB chỉ được so bằng nhau qua các snapshot. B đã bắt source ID 1
và final state 137/150.

Khuyến nghị: port hai assertion literal vào production oracle A và thêm fixture
Sói sai ID cùng HP max sai.

### Important 5 — fixture suite A chưa pressure-test đủ các nhánh fail-closed

A có 10 case; B có 14 case. A thiếu case trực tiếp cho:

- clean PASS nhưng warm FAIL;
- artifacts rỗng;
- evidence rỗng;
- SHA đúng định dạng nhưng sai nội dung;
- server lock và asset lock sai độc lập.

Một phần logic runtime A có thể đã reject các trường hợp này, nhưng chưa có
fixture chạy qua chính runner boundary để chứng minh.

### Warning — thuật toán sequence B không được port nguyên trạng

B dùng một biến `$previous` chung toàn NDJSON. Với log hợp nhất hai client,
sequence `A:1, B:1` sẽ bị coi là lùi/trùng. A dùng `$lastByClient` và phù hợp hơn.

## 5. Mã lỗi canonical đề xuất

Giữ nguyên taxonomy A nơi đã có consumer/test, chỉ bổ sung mã cho invariant thật
sự chưa có:

| Invariant | Mã canonical | Xử lý mã B |
|---|---|---|
| Damage count/value/path | `E2E-COMBAT-010` | Không giữ `E2E-DAMAGE-010`; remap fixture B |
| HP cuối/HP max | `E2E-HP-011` | Giữ nguyên, hai bộ cùng dùng |
| Executable clean/warm | `E2E-EXECUTABLE-012` | Giữ nguyên |
| Asset lock | `E2E-ASSET-013` | Không tạo alias `E2E-ASSET-LOCK-013` |
| Scope clean/warm | `E2E-SCOPE-014` | Remap `E2E-WARM-SCOPE-014` |
| DB snapshot/SHA | `E2E-DB-SHA-015` | Không dùng generic `E2E-SHA-006` |
| Evidence SHA/file | `E2E-EVIDENCE-SHA-016` | Dùng cho missing/mismatch SHA |
| Server lock | `E2E-SERVER-LOCK-017` | Port mã B vì đây là domain tách biệt |
| Warm negotiation/resource | `E2E-WARM-017` | Giữ A |
| Empty artifact/evidence | `E2E-EVIDENCE-002` | Dùng mã đã tồn tại trong A; không giữ `E2E-ARTIFACT-002` |
| Sequence lùi/trùng | `E2E-SEQUENCE-003` | Giữ nguyên |
| Sequence đứt | `E2E-SEQUENCE-004` | Giữ nguyên |
| Persistence | `E2E-PERSIST-004` | Giữ A |
| Visual review | `E2E-VISUAL-018` | Giữ A |
| Clean/warm không cùng PASS | `E2E-CLEAN-WARM-016` | Port từ B vì invariant aggregate riêng, không trùng full code evidence |

Việc remap phải diễn ra trong fixture/harness A. Không nên duy trì hai mã đồng
nghĩa vì sẽ làm runner/report downstream khó ổn định.

## 6. Fixture bắt buộc sau tích hợp

Giữ/port theo **ý nghĩa**, nhưng viết lại theo directory schema của A:

1. valid clean full + warm smoke;
2. extra damage;
3. damage sai source ID Sói;
4. HP cuối 139;
5. HP max khác 150;
6. executable mismatch;
7. server-lock mismatch;
8. asset-lock mismatch;
9. warm scope bị ghi `full`;
10. clean PASS/warm FAIL;
11. thiếu DB SHA;
12. thiếu artifact/evidence SHA;
13. SHA đúng định dạng nhưng sai nội dung;
14. artifacts array rỗng;
15. evidence record/mảng rỗng;
16. evidence file zero-byte;
17. sequence backward/trùng theo một client;
18. sequence gap theo một client;
19. interleave `A:1, B:1, A:2, B:2` phải PASS;
20. warm negotiation false hoặc Dragon Ball fallback true phải đỏ.

## 7. Có thể tích hợp B vào A không?

**Có, nhưng chỉ bằng port thủ công có chọn lọc.**

Không cherry-pick `3684626` nguyên commit vì:

- add/add conflict tại `Test-RunnerOracle.ps1` đã được `git merge-tree` xác nhận;
- hai file cùng tên có vai trò khác nhau: A là harness đi qua runner, B là parser;
- fixture B không tương thích schema A;
- chọn version B khi resolve sẽ làm mất runner boundary;
- chọn version A sẽ khiến hầu hết fixture B được add nhưng không được chạy.

Cách tích hợp an toàn:

1. lấy A làm nhánh đích;
2. giữ `Test-Evidence.ps1` là production oracle duy nhất;
3. giữ `Test-RunnerOracle.ps1` của A làm harness gọi `-OracleOnly`;
4. port guard/fixture B vào schema A, không copy parser B;
5. sửa asset/server lock tách biệt và manifest scope;
6. chạy lại toàn bộ suite/mutation trước khi tạo commit tích hợp riêng.

Rủi ro conflict:

- **Git/text:** trung bình — một add/add conflict trực tiếp;
- **schema:** cao — uppercase/lowercase status, inline/nested result khác nhau;
- **mã lỗi:** cao — nhiều mã đồng nghĩa nhưng tên khác;
- **runner:** rất cao nếu resolve nhầm về B, vì B không được runner production gọi;
- **fixture:** trung bình — tên không trùng nhưng không có nghĩa là chạy được.

## 8. Kiểm chứng lại bộ B trong worktree hiện tại

Không có file oracle/fixture nào bị sửa.

```text
POWERSHELL-PARSE errors=0
ORACLE-SELFTEST-PASS: 14/14
ORACLE-PASS
B-RECHECK parse_errors=0 suite_exit=0 valid_exit=0 bad_exit=1 bad_has_E2E-HP-011=True
```

Mutation được thực hiện trên source biến đổi **trong bộ nhớ** ở tiến trình
PowerShell con; không ghi bản mutated ra đĩa:

```text
MUTATION remove_hp_137 exit=1 detected=True
SELFTEST-FAIL hp-final-not-137.json: expected=E2E-HP-011 actual=PASS

MUTATION remove_asset_lock_equality exit=1 detected=True
SELFTEST-FAIL asset-lock-mismatch.json: expected=E2E-ASSET-LOCK-013 actual=PASS

MUTATION remove_clean_warm_both_pass exit=1 detected=True
SELFTEST-FAIL clean-pass-warm-fail.json: expected=E2E-CLEAN-WARM-016 actual=PASS
```

Cả ba mutation đều làm đúng một case lọt và suite exit 1.

## 9. Test bắt buộc sau khi hợp nhất

Tối thiểu phải chạy, theo thứ tự:

1. PowerShell parse cho `run-unity-e2e.ps1`, `Test-Evidence.ps1` và
   `Test-RunnerOracle.ps1`;
2. `Test-Evidence.ps1 -SelfTest`;
3. `run-unity-e2e.ps1 ... -SelfTest`;
4. suite oracle canonical qua `run-unity-e2e.ps1 -OracleOnly`, đủ toàn bộ fixture
   ở mục 6; valid exit 0, từng fixture xấu exit 1 đúng mã;
5. kiểm aggregate chỉ được sinh sau clean full + warm smoke cùng PASS;
6. mutation tối thiểu: bỏ extra-damage check, HP 137/150, executable equality,
   server-lock equality, asset-lock equality, clean+warm đồng thời, evidence
   length/SHA và sequence per-client;
7. hoàn nguyên mutation và chạy xanh lại;
8. `git diff --check` và xác nhận không có fixture orphan không được harness gọi.

Các test trên chỉ chứng minh oracle. Clean/warm gameplay thật, Unity build,
graceful restart và DB thật vẫn là cổng riêng; chỉ sau khi chúng thực sự chạy
mới được đánh giá Task 10 full E2E.

## Verdict

**Canonical baseline: A, REQUEST CHANGES trước acceptance.**

B không nên trở thành parser production, nhưng là bộ pressure-test có giá trị và
nên được port thủ công. Không cherry-pick nguyên commit B, không duy trì hai
schema song song, và không bắt đầu Task 11 từ kết quả oracle unit-level này.
