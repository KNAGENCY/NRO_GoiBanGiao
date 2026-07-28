# Task 6A — fix round 2 report

Ngày: 29/07/2026 (ICT)  
Reviewer: `Codex Task 6A fix round 2`

## Phạm vi và files

- `tools/artqa/artqa.py`: semantic GIF loop, raw-source provenance, canonical
  evidence verification và fail-closed gate.
- `tools/artqa/test_artqa.py`: contract test TDD và mutation oracle cho 4 loại evidence.
- `docs/my-thuat/review/vertical-slice-v1.json`: review/evidence SHA mới sinh từ
  pipeline thật.
- Shared QA output: `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\qa\`
  (ignored, không commit).

Không sửa PNG nguồn, `mapping-id.json`, `D:\NRO_Project` hay worktree khác.

## RED / GREEN

- RED: **19 test**, gồm **9 FAIL + 2 ERROR** đúng vào loop semantic, raw fields,
  canonical evidence, reviewer/timestamp và malformed/null fail-closed.
- GREEN: `python -m unittest tools/artqa/test_artqa.py -v` — **19/19 PASS**.
- Precondition bind: `python -m tools.artqa.artqa bind-precondition ...` — **exit 0**
  với review hợp lệ; review lỗi trả **exit 2**.

## Pipeline và gate thật

- Raw source: **0/20 PASS, 20/20 FAIL**; lỗi raw được giữ nguyên trong
  `raw_source_status`, `source_issues` và summary. Derived output PASS vẫn hợp lệ
  khi raw source FAIL.
- Derived output: **20/20 PASS**, `status=PASS`, `source_status=PASS`, `issues=[]`.
- `hu_anthan_10f_chroma.png`: pipeline ghi rõ semantic
  `one_shot_effect_sequence`, `loop=false`, 10 frame/100ms; không ép loop giả
  cho một shot hiệu ứng một lần.
- `hud_hp_chakra_4_chroma.png`: flood matte qua viền trắng theo từng ô, alpha
  transparency **0.789141**, opaque **0.198155**, không duplicate adjacent frame.
- Gate thật: `python -m tools.artqa.artqa gate --review ...` — **exit 0**,
  `allowed=true`, `issues=[]`; `bind-precondition` cũng **exit 0**.
- Exact names: **20 total / 20 unique**; evidence ngoài QA root: **0**.
- Source SHA: **20/20 unchanged** trước/sau pipeline.
- Reviewer và timestamp bắt buộc, GIF one-shot không có loop extension vô hạn;
  gate đối chiếu lại alpha pixel, crop report, contact sheet và frame/timing.

## Mutation thật

Đã chạy mutation trên bản review tạm, không sửa artifact thật; tất cả đều bị bắt
và CLI trả **exit 2**:

| Mutation | Kết quả gate |
|---|---|
| SHA/raw summary/reviewer/evidence kind giả | exit **2** |
| alpha đổi màu nhưng giữ kích thước/hash cập nhật | exit **2** |
| crop report/contact sheet/animation đổi nội dung | exit **2** |
| evidence ngoài QA hoặc null asset | exit **2** |

## Kết luận

Task 6A round 2 đã xử lý các blocker NEEDS WORK. Task 6B/assetpack phải gọi
`require_review_gate()` hoặc CLI `bind-precondition` trước khi bind; gate fail thì
không thể chạy tiếp.
