# SDD ledger — plan: docs/superpowers/plans/2026-07-29-ke-hoach-vertical-slice-lang-la.md

Task 0: complete (commits d2576a6..faec0fc, review/merge đã hoàn tất trước phiên Codex này)
Task 1: complete (commits a7b889e..736e32b, review/merge đã hoàn tất trước phiên Codex này)
Task 2: complete (commits cfde9fe..25b2bd7, review/merge đã hoàn tất trước phiên Codex này)
Task 3: complete (commits 21f05b4..153e627, review/merge đã hoàn tất trước phiên Codex này)
Task 4: complete (commits a135979..0a4cc80, review/merge đã hoàn tất trước phiên Codex này)
Task 5: complete (commits d4879b6..0ba0ece, review/merge đã hoàn tất trước phiên Codex này)
Task 6A: complete (commits d4b35e3..ee7ef8d, 20/20 ảnh PASS theo bàn giao người dùng)
Task 6B: complete (commits 0f2805f..05b84f2, production closure/canonical generator/crop đã commit)

Task 6B: reopened — baseline tại 05b84f2 có 94 test, 3 fail + 1 error:
manifest thiếu `skills`, `ui_bindings`, các SHA `build_inputs`, và contract
đường dẫn không resolve được shared ArtRoot. Task 7 bị khóa cho tới khi suite
Task 6B xanh thật.
Task 6B: fix round 4/5 BLOCKED tại visual gate HUD: nội dung x4 chỉ rộng
59/240 px. Commit an toàn đã giữ lại: 4198e70.
Task 6B: fix round 5/5 in progress — truy vết phép scale/crop HUD trước khi sửa.
Task 6B: fix round 5/5 implementation commit bf91ee8; fresh verification
100/100 PASS, build 667 files, lock
19b84a6da6fb2ca7bfe2352bb0d4066a3a18e680f8ec6a8e74b8c362a646b0bf.
Task 6B: fix round 5/5 BLOCKED after independent review — 2 Important:
validator accepts 19/20 source provenance; HUD/Iruka artifact evidence and
reviewer-signed SHA are not bound into review JSON. Task 7 remains locked.
Task 6B: fix round 6 explicitly authorized by user — scope is exactly the
two open Important findings above; TDD/mutation/full suite/build A-B and
independent scoped re-review are mandatory. Task 7 remains locked until
the re-review returns all findings addressed with no new Critical/Important.
Task 6B: fix round 6 re-review 1 (1 addressed, 1 open; commit 545b9c6) —
exact-20 provenance ADDRESSED; artifact reviewer identity NOT ADDRESSED
because a self-signed alias is still accepted. Same implementer resumed;
Task 7 remains locked.
Task 6B: fix round 6 follow-up (1 addressed, 0 open; commit 0fd6fc5) —
trusted artifact reviewer identity pinned outside review JSON; self-signed
alias mutation rejected; scoped re-review clean.
Task 6B: complete (commits bf91ee8..0fd6fc5, review clean) — 109/109 tests,
CLI validate PASS, A/B 669 files with DIFF_COUNT=0, visual HUD/Iruka PASS.
Task 7: in progress from base a7062e9 — fresh implementer `task7_impl`;
Task 8 remains locked until Task 7 independent review PASS.
Task 7: implementation commit 2a77a43; independent review FAIL — 3 Important:
previous lock optional lets same-version mutation pass; protected Unity root
optional; missing x2–x4 resource falls back to x1. Fix round 1/5 commit
9627680 addressed the first two and tightened zoom/dimension validation; scoped
re-review left one Important about canonical resource selection.
Task 7: fix round 2/5 commit 56fa302 added closure-driven canonical paths for
small/skill/mob/background/HUD, plus a decoy Small ID regression and mutation;
scoped re-review approved, no new Critical/Important.
Task 7: complete (commits 2a77a43..56fa302, review clean). Task 8 remains
locked until the coordinator explicitly starts the next task.

Task 8: complete (commits de1f0c6..f53a255, review clean). Catalog/NPC/mob
bindings are explicit; collision artifact is semantic 0/1 with 24x24 header;
NPC command 33 uses XepGui FIFO. Full Python 124/124, Go game -count=1,
net -count=5, and vet pass. Next is Task 9/telemetry; branch not merged to main.

Resume: Task 8, after fresh verification of Task 7 and bookkeeping.
Preflight: hai sửa đổi dở ngoài Task 7 được bảo toàn tại stash
`codex-pre-task7-stale-task6b`; cache Python tại stash
`codex-pre-task7-python-cache`. Worktree sạch ở HEAD 05b84f2 trước baseline.
