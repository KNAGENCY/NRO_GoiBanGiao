# Task 6A — fix round 2 report

Ngày: 29/07/2026 (ICT)  
Reviewer: `Codex Task 6A NEEDS WORK closure`

## Phạm vi và files

- `tools/artqa/artqa.py`: gate fail-closed, precondition `require_review_gate`,
  CLI `gate` và `bind-precondition`, kiểm exact 20 asset/source/evidence.
- `tools/artqa/test_artqa.py`: contract test TDD và mutation oracle.
- `docs/my-thuat/review/vertical-slice-v1.json`: review/evidence SHA mới sinh từ
  pipeline thật.
- Shared QA output: `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\qa\`
  (ignored, không commit).

Không sửa PNG nguồn, `mapping-id.json`, `D:\NRO_Project` hay worktree khác.

## RED / GREEN

- RED: 5 contract test mới fail đúng: thiếu exact-20/source SHA/source status,
  evidence ngoài QA/trùng, và thiếu bind precondition API.
- GREEN: `python -m unittest tools/artqa/test_artqa.py -v` — **13/13 PASS**.
- Precondition bind: `python -m tools.artqa.artqa bind-precondition ...` — **exit 0**
  với review hợp lệ; review lỗi trả **exit 2**.

## Pipeline và gate thật

- Raw source: **0/20 PASS, 20/20 FAIL**; các lỗi raw được giữ trong
  `raw_source_status`/`raw_source_*`, không bị giả thành raw PASS.
- Derived output: **20/20 PASS**, `status=PASS`, `source_status=PASS`, `issues=[]`.
- `hu_anthan_10f_chroma.png`: pipeline ghi rõ semantic
  `one_shot_effect_sequence`, `loop=false`, 10 frame/100ms; không ép loop giả
  cho một shot hiệu ứng một lần.
- `hud_hp_chakra_4_chroma.png`: flood matte qua viền trắng theo từng ô, alpha
  transparency **0.789141**, opaque **0.198155**, không duplicate adjacent frame.
- Gate thật: `python -m tools.artqa.artqa gate --review ...` — **exit 0**,
  `allowed=true`, `issues=[]`.
- Exact names: **20 total / 20 unique**; evidence ngoài QA root: **0**.
- Source SHA: **20/20 unchanged** trước/sau pipeline.

## Mutation thật

Đã chạy mutation trên bản review tạm, không sửa artifact thật; tất cả đều bị bắt
và CLI trả **exit 2**:

| Mutation | Kết quả gate |
|---|---|
| thiếu asset | bắt `exact_required_assets_mismatch` |
| sai source SHA | bắt `source_sha256_mismatch` |
| `source_status=FAIL` | bắt `source_status_not_pass` |
| còn `issues` | bắt `issues_not_empty` |
| evidence trùng/giả | bắt `evidence_duplicate`/SHA hoặc content mismatch |
| CLI gate trên review lỗi | exit **2** |

## Kết luận

Task 6A đã đóng blocker review: Task 6B/assetpack phải gọi
`require_review_gate()` hoặc CLI `bind-precondition` trước khi bind; gate fail thì
không thể chạy tiếp.
