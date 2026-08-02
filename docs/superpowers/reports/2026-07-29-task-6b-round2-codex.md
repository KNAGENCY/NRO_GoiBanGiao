# Task 6B round 2 — canonical base hardening và blocker production

Ngày: 29/07/2026  
Worktree: `kh-vertical-slice`

## RED trước khi sửa

- Test production canonical mới đỏ vì `docs/my-thuat/mapping-id.json` thiếu đúng
  `skills`, `ui_bindings`, `build_inputs` so với `assetpack.build.REQUIRED_TOP_LEVEL`.
- Test build thật đỏ vì mapping không có `build_inputs.source_base` và các đường
  `planned/27041/x1.png` không tồn tại.
- Test bind base đỏ vì bind chỉ so basename, chưa resolve `source_path` từ build base.

## GREEN đã làm

- `assetpack` CLI nhận `build --base-dir` và truyền cùng base vào `build()`.
- Bind kiểm tra path tương đối, không vượt base, file tồn tại và SHA khớp review
  khi được truyền `art_root`.
- Test wiring CLI và test bind source-base xanh.
- Bộ test assetpack hiện hữu: `84/84 PASS`.
- Art QA: `19/19 PASS`; review gate thật: `allowed=true`, `issues=[]`.

## Mutation/production evidence

- Review gate vẫn được gọi trước schema/output; build lỗi thì không để lại output.
- Build thật với art root `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1`
  trả `ASSET-REVIEW-010: thiếu source đã review: planned/27041/x1.png` và
  `OUTPUT_EXISTS=False` ở cả hai lần chạy staging độc lập.
- Bind thật trả `ASSET-SCHEMA-001: manifest thiếu/thừa trường canonical` và
  `OUTPUT_EXISTS=False`.
- `verify-legacy-inputs` trả `ASSET-LEGACY-009`.
- Ledger hiện trỏ `assetpack/tests/fixtures/legacy/mob.bin`, khai báo 11 byte và
  SHA `d38c22…b5c6d`; byte hiện tại là 227 byte, SHA
  `b7602310dbb1d0c668b972d27713b24a1d74feea7a9bb1b8398673df8273e59c`.

## Blocker không được phép giả PASS

Chưa thể tạo canonical production manifest/build vì:

1. `mapping-id.json` vẫn là baseline, chưa có canonical closure và các nguồn x1–x3
   thật; `planned/27041/x1.png` là đường không tồn tại.
2. `legacy-input-ledger.json` còn trỏ fixture cho `NR_part`, `NR_image`, skill,
   mob, Small và C# hard-code. Chưa có các byte canonical production tương ứng;
   không đổi SHA khai báo để che mismatch.
3. Vì vậy chưa sửa mapping/ledger và chưa sinh artifact server/client/lock. Không
   dùng fixture placeholder hoặc sửa PNG nguồn để làm build xanh.

`BAT-DAU-TU-DAY.md`, `NHAT-KY-AGENT.md` và PNG nguồn không bị sửa.
