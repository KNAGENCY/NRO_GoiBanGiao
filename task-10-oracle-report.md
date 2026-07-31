# Task 10 — báo cáo oracle/fixture độc lập

Ngày: 31/07/2026

Nhánh: `kh-task10-oracle`

Base đã kiểm tra: `01af875cfd5b3c6403c5b95838bfb94b5346b07c`

## Phạm vi thật

Phiên này chỉ tạo oracle/fixture unit-level để kiểm tra kết quả aggregate của
Task 10. Không mở Unity, không khởi động server gameplay, không chạy hai client
thật và không tạo bằng chứng giả rằng gameplay E2E đã PASS.

Không sửa `run-unity-e2e.ps1`, `New-Fixture.ps1`, `Test-Evidence.ps1`,
`server/` hay `tools/unity-overlay/`.

File `task-9-report.md` được yêu cầu đọc nhưng không tồn tại trong worktree hoặc
tree của base. Bối cảnh Task 9 được đối chiếu từ commit base, spec, plan,
`BAT-DAU-TU-DAY.md`, `NHAT-KY-AGENT.md`, `CLAUDE.md` và ba script E2E đã có.

## Hợp đồng oracle

`tools/e2e/tests/Test-RunnerOracle.ps1` chỉ PASS khi:

- fixture JSON schema 1 có đủ aggregate, clean và warm;
- clean là `scope=full`, warm là `scope=warm_smoke`;
- aggregate, clean và warm cùng `status=PASS`;
- ba nơi cùng executable SHA, server-lock SHA và asset-lock SHA;
- artifacts/evidence đều không thiếu, không rỗng; file thật không rỗng và
  SHA-256 tồn tại, đúng định dạng, khớp nội dung;
- clean có DB evidence sau shutdown, sau relogin và NDJSON; warm có cache
  evidence và NDJSON;
- mỗi NDJSON bắt đầu từ sequence 1, tăng đúng từng đơn vị, không lùi/trùng/đứt;
- clean có đúng một damage Sói `150→145` và bốn damage PvP
  `145→143→141→139→137`, không có damage thừa;
- trạng thái cuối clean đúng `137/150`.

Mã lỗi ghim:

| Mã | Ý nghĩa |
|---|---|
| `E2E-DAMAGE-010` | Damage thiếu/thừa/sai chuỗi |
| `E2E-HP-011` | HP cuối không phải 137/150 |
| `E2E-EXECUTABLE-012` | Executable SHA không đồng nhất |
| `E2E-ASSET-LOCK-013` | Asset-lock SHA không đồng nhất |
| `E2E-WARM-SCOPE-014` | Warm bị ghi sai scope |
| `E2E-CLEAN-WARM-016` | Aggregate/clean/warm không cùng PASS |
| `E2E-SERVER-LOCK-017` | Server-lock SHA không đồng nhất |
| `E2E-SHA-006` | Evidence thiếu/sai SHA, file thiếu/rỗng hoặc hash lệch |
| `E2E-SEQUENCE-003` | Sequence lùi hoặc trùng |
| `E2E-SEQUENCE-004` | Sequence đứt |
| `E2E-ARTIFACT-002` | Danh sách artifact thiếu/rỗng |
| `E2E-EVIDENCE-002` | Danh sách evidence/NDJSON thiếu/rỗng |

## TDD — RED trước implementation

Lúc chạy RED, toàn bộ fixture JSON đã tự parse được; hàm oracle vẫn cố ý chưa
được cài.

Lệnh:

```powershell
powershell.exe -NoProfile -ExecutionPolicy Bypass `
  -File tools/e2e/tests/Test-RunnerOracle.ps1 -SelfTest
```

Raw output, exit `1`:

```text
JSON-PARSE failures=0
SELFTEST-FAIL valid.json: expected=PASS actual=ORACLE-NOT-IMPLEMENTED
SELFTEST-FAIL damage-extra.json: expected=E2E-DAMAGE-010 actual=ORACLE-NOT-IMPLEMENTED
SELFTEST-FAIL hp-final-not-137.json: expected=E2E-HP-011 actual=ORACLE-NOT-IMPLEMENTED
SELFTEST-FAIL executable-mismatch.json: expected=E2E-EXECUTABLE-012 actual=ORACLE-NOT-IMPLEMENTED
SELFTEST-FAIL asset-lock-mismatch.json: expected=E2E-ASSET-LOCK-013 actual=ORACLE-NOT-IMPLEMENTED
SELFTEST-FAIL server-lock-mismatch.json: expected=E2E-SERVER-LOCK-017 actual=ORACLE-NOT-IMPLEMENTED
SELFTEST-FAIL warm-scope-full.json: expected=E2E-WARM-SCOPE-014 actual=ORACLE-NOT-IMPLEMENTED
SELFTEST-FAIL db-evidence-missing-sha.json: expected=E2E-SHA-006 actual=ORACLE-NOT-IMPLEMENTED
SELFTEST-FAIL sequence-backward.json: expected=E2E-SEQUENCE-003 actual=ORACLE-NOT-IMPLEMENTED
SELFTEST-FAIL sequence-gap.json: expected=E2E-SEQUENCE-004 actual=ORACLE-NOT-IMPLEMENTED
SELFTEST-FAIL clean-pass-warm-fail.json: expected=E2E-CLEAN-WARM-016 actual=ORACLE-NOT-IMPLEMENTED
SELFTEST-FAIL artifacts-empty.json: expected=E2E-ARTIFACT-002 actual=ORACLE-NOT-IMPLEMENTED
SELFTEST-FAIL evidence-empty.json: expected=E2E-EVIDENCE-002 actual=ORACLE-NOT-IMPLEMENTED
SELFTEST-FAIL evidence-sha-mismatch.json: expected=E2E-SHA-006 actual=ORACLE-NOT-IMPLEMENTED
ORACLE-SELFTEST-FAIL: 14/14 fixture sai
```

## GREEN

Sau implementation tối thiểu, raw output, exit `0`:

```text
SELFTEST-PASS valid.json: PASS
SELFTEST-PASS damage-extra.json: E2E-DAMAGE-010
SELFTEST-PASS hp-final-not-137.json: E2E-HP-011
SELFTEST-PASS executable-mismatch.json: E2E-EXECUTABLE-012
SELFTEST-PASS asset-lock-mismatch.json: E2E-ASSET-LOCK-013
SELFTEST-PASS server-lock-mismatch.json: E2E-SERVER-LOCK-017
SELFTEST-PASS warm-scope-full.json: E2E-WARM-SCOPE-014
SELFTEST-PASS db-evidence-missing-sha.json: E2E-SHA-006
SELFTEST-PASS sequence-backward.json: E2E-SEQUENCE-003
SELFTEST-PASS sequence-gap.json: E2E-SEQUENCE-004
SELFTEST-PASS clean-pass-warm-fail.json: E2E-CLEAN-WARM-016
SELFTEST-PASS artifacts-empty.json: E2E-ARTIFACT-002
SELFTEST-PASS evidence-empty.json: E2E-EVIDENCE-002
SELFTEST-PASS evidence-sha-mismatch.json: E2E-SHA-006
ORACLE-SELFTEST-PASS: 14/14
```

## Mutation

Mỗi mutation được áp dụng trực tiếp vào oracle, chạy suite, rồi hoàn nguyên.
Ba lượt đều exit `1`, chỉ đúng fixture bảo vệ chốt đó bị lọt:

```text
# Bỏ check HP 137
SELFTEST-FAIL hp-final-not-137.json: expected=E2E-HP-011 actual=PASS
ORACLE-SELFTEST-FAIL: 1/14 fixture sai

# Bỏ check asset-lock equality
SELFTEST-FAIL asset-lock-mismatch.json: expected=E2E-ASSET-LOCK-013 actual=PASS
ORACLE-SELFTEST-FAIL: 1/14 fixture sai

# Bỏ check clean + warm cùng PASS
SELFTEST-FAIL clean-pass-warm-fail.json: expected=E2E-CLEAN-WARM-016 actual=PASS
ORACLE-SELFTEST-FAIL: 1/14 fixture sai
```

Sau hoàn nguyên, suite được chạy xanh lại `14/14`.

## Kiểm chứng cuối trước commit

```text
POWERSHELL-PARSE errors=0
ORACLE-SELFTEST-PASS: 14/14
ORACLE-PASS
DIRECT-CHECK suite_exit=0 valid_exit=0 bad_exit=1 bad_code_hits=3
```

`bad_code_hits=3` là do PowerShell in cùng mã `E2E-HP-011` ở message, vị trí
throw và FullyQualifiedErrorId; điều quan trọng là fixture xấu exit `1` và mã
ổn định có mặt.

## Trạng thái acceptance

Phần oracle/fixture unit-level trong phạm vi worktree này đã có RED, GREEN và
mutation. Acceptance gameplay Task 10 vẫn chưa được chứng minh: clean full E2E,
warm-upgrade smoke, Unity screenshot/reviewer, graceful restart và DB thật còn
thuộc luồng implementation/runner khác.
