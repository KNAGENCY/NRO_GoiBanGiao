# Task 3 — Lắp thử ba mẫu UI vào Unity và chụp bằng chứng

## Kết luận

**BLOCKED riêng phần bằng chứng HUD trong game.** Phần lắp asset, backup, SHA/kích thước,
popup và button đã PASS. Unity chạy tới popup báo máy chủ mất kết nối; không có tài khoản
test/máy chủ test sẵn trong lượt này nên không thử hoặc đoán credential, không dùng production.
Vì vậy không có ảnh HUD sau login và không được báo toàn bộ Task 3 là PASS.

- Unity cuối lượt: đúng `D:/NRO_Project/Assets`, Edit Mode, không Play, không sắp Play,
  không compile, không update.
- Đã chép đúng 15 PNG runtime: 5 tệp × x2/x3/x4. Không chép x1 vào Unity.
- Backup phục hồi 15/15 ở
  `D:\NRO_GoiBanGiao\.worktrees\dot-4-thi-cong\docs\my-thuat\lang-la-client-assets\mainimage-truoc-dot-4`.
- `Small27083/27084` còn đủ x1..x4 và hash trước/sau không đổi.
- Không đổi production code, `server/internal/game/assets/asset-lock.json`, manifest/hash ký
  hoặc ba PNG canonical.
- Không dùng `kill`/`Stop-Process`. Unity được đóng một lần bằng yêu cầu đóng cửa sổ bình
  thường rồi mở lại; cấu hình auto-start MCP tạm thời đã được trả về trạng thái ban đầu
  (key không tồn tại).

## Source canonical

| Tệp | Kích thước | SHA-256 |
|---|---:|---|
| `ui_hud_khung_trai_x4.png` | 620×164 | `1eb658cc3bffeb74e2105efe178f5950db865b3d8949dfcbf98837156cec13fd` |
| `ui_khung_popup_9o_x4.png` | 40×320 | `0fc7659d319bb756ba954a51dae9b570bef171b24bfe8f4cd457be2d9e6bc2e6` |
| `ui_nut_thuong_3manh_x4.png` | 120×96 | `15edc2e5ebf57e29e05b0276d1d5cacea14952c1abb74103f0275c7228638d6c` |

Nguồn: `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\08-giao-dien`.
Hash được kiểm lại trước cả backup và mọi lần install.

## TDD và mutation

### RED ban đầu

Harness gọi tính năng chưa tồn tại; bốn test đều đỏ vì thiếu
`tools/ui-dot4/lap_ba_mau_ui.py`:

```text
4 failed
```

Các test khóa bốn cổng: backup thiếu, Unity đang Play, source sai SHA và mapping chính xác
15 target.

### RED phát hiện từ MCP thật

Lần backup thật đầu tiên bị dừng trước khi tạo backup vì MCP mới bọc state trong
`data.result`; parser cũ chỉ hiểu text thẳng:

```text
UI-DOT4-UNITY-STATE-INVALID raw='{\\"data_path\\":\\"D:/NRO_Project/Assets\\",...'
```

Test hồi quy được thêm trước khi sửa:

```text
_______________ test_doc_duoc_state_tu_lop_boc_json_cua_mcp_moi _______________
E       AssertionError: UI-DOT4-UNITY-STATE-INVALID raw=...
1 failed, 4 deselected in 0.39s
```

Sau khi parser đọc đúng JSON envelope:

```text
.                                                                        [100%]
1 passed, 4 deselected in 0.28s
```

### GREEN cuối

```text
.....                                                                    [100%]
5 passed in 1.28s
```

`git diff --check` sạch.

### Mutation cổng backup

Cố ý thay `verify_backup()` bằng kết quả rỗng trong fixture. Test
`test_preflight_tu_choi_install_khi_backup_thieu` đỏ vì cổng mong đợi
`UI-DOT4-BACKUP-MISSING` không còn chạy và chương trình đi tiếp tới lỗi staging:

```text
E assert 'UI-DOT4-BACKUP-MISSING' in 'UI-DOT4-STAGING-COUNT ...'
1 failed
```

Mutation đã được hoàn nguyên bằng patch; full suite xanh lại như trên.

## Unity preflight và cài đặt

Lần kiểm đầu không có listener MCP 8080 nên fail-closed, không backup/copy:

```text
<urlopen error [WinError 10061] No connection could be made because the target machine actively refused it>
```

Sau khi khởi động lại Unity bằng đường chính thức và bật auto-start tạm thời:

```text
MCP_8080_LISTEN=true
UI_DOT4_STATE={"data_path":"D:/NRO_Project/Assets","is_playing":false,
"is_playing_or_will_change":false,"is_compiling":false,"is_updating":false}
```

Pipeline thật:

```text
UI-DOT4-BACKUP-CREATED count=15
Da sinh du 4 muc cho : 3/3 anh
Tong file xuat ra    : 12
UI-DOT4-STAGED count=15
UI-DOT4-INSTALLED count=15
```

Backup gốc trước install được tạo tại
`D:\NRO_Project\.codex-backups\ui-dot4-task3`. Sau đó nó được sao nguyên byte vào
đường bàn giao bắt buộc trong repo:

```text
HANDOFF_BACKUP_FILES=17 PNG=15 JSON=2
HANDOFF_BACKUP_SHA_MISMATCH=0
```

Để sửa đúng thứ tự handoff đã nêu trong brief, lần install cuối được chạy idempotent sau
khi backup repo 15/15 đã tồn tại và được verify; chính backup repo là `--backup-root` của
lần install cuối:

```text
UI-DOT4-INSTALLED count=15
Refresh requested ... resulting_state="idle"
UI-DOT4-VERIFIED count=15
```

Lưu ý audit: lần chép đầu đã có đủ backup gốc ngoài `Assets`, nhưng bản mirror ở đúng đường
repo được tạo ngay sau lần chép đầu. Lần chép cuối đã chạy lại qua đúng backup repo để khép
cổng quy trình; không có tệp gốc nào bị mất.

## SHA backup và bản đã lắp

| Zoom / target | Cỡ backup → installed | SHA backup | SHA installed |
|---|---:|---|---|
| x2 `myTexture2dpanel.png` | 309×82 → 310×82 | `0b01a9a2f836951225d689c0bf824c3809f705092d1f79157d55f05ca664846a` | `a490055cfac42ed20f01084663b01be61af7fa96f3bb5da797030fa3fcce2a10` |
| x2 `myTexture2dimgPopup.png` | 20×160 → 20×160 | `161f893fedf5a0b7973146c7c8911b6f9aae8c93a9f23f6c8246fe0671e5034a` | `0ecaa645d5be2103d42db18d6cb4486affb40a3ffb7f0d54f8191b04c35ab7ca` |
| x2 `btn0left.png` | 20×48 → 20×48 | `ba2b42d3d7a99a0aaa9f9dab54f59c9fc12c8a2d2cc90b583f63a05cb2403fb0` | `9f4796a6ce3e850962a9f4e9ae7dfc910208a82a9d64ec1ff56afc91d09eb450` |
| x2 `btn0mid.png` | 20×48 → 20×48 | `04be88cbf9816d6ceb3cca0985ff6399248ebc57108edfeb52c1474e53f55106` | `3ec2e44367202b7094777efe92a93e348b536873c90f217fc0f43f557d7ba88e` |
| x2 `btn0right.png` | 20×48 → 20×48 | `74df3cfeb1ee2686bea659f0cf00e2f97c8c59d4793cec5a4723d719b4a1a05b` | `58f2e9bbf8f34f66aaf17e8c2dffb7b71c847bf280507e478352f2ecdc0351ba` |
| x3 `myTexture2dpanel.png` | 463×123 → 465×123 | `86b8ca87a22371c92a8a1f088c4722a61e9dbfe80c564dc30337e96ff5c431d5` | `99a583bd3c60a268edb8cb492e5db41b007f6bba89b3dbc6fd24cf6a628b62ce` |
| x3 `myTexture2dimgPopup.png` | 30×240 → 30×240 | `e605099bef3f8f7cd96d116ec7e14cf89c1e66f6601456865eedca3a0ca3a6ab` | `5d9ef9d74828f21bd914e35f01edf4081d91b7c1e0008aca15f0a99f0fb0c71f` |
| x3 `btn0left.png` | 30×72 → 30×72 | `899e92bb6ab8cca2b6fc35e316440f54e17e649e069bb4d22ba223343292ffa2` | `f14f281d298753742f9bfb485f9b794523a9869240caba0a0d9840de5fa82636` |
| x3 `btn0mid.png` | 30×72 → 30×72 | `599ebd6bba21243750b867c2a44a8443bcd65c9f9aa1b40370a58376994f663c` | `1d9e3497933203598820349f906c5097710b08ce92dffab8a06bd9e90551c208` |
| x3 `btn0right.png` | 30×72 → 30×72 | `254f5194e70b4e4d2f061fae991ed815cf09c02de1dcc6a61c3ca6b6e1984775` | `979db756b33176a4f80fa461119ab8a99d4dc965772b982d3b8d057ffd8490ae` |
| x4 `myTexture2dpanel.png` | 618×164 → 620×164 | `f049195a041545090cc1fc8d9c94ca316ba2272b3678b293da89a04a24e5d76e` | `e37f1ac4a39649fd0086709b358ec216b11ccaec28630edb6bada09e7419a80a` |
| x4 `myTexture2dimgPopup.png` | 40×320 → 40×320 | `9a9bec501a49f7059eb89d2694dfb538d5c1b6f2b66dd2c3939f6505b6cc4c6e` | `9db5d749fcef37846ac07498631143c341e58738218e5de75117fd65fe7a0b15` |
| x4 `btn0left.png` | 40×96 → 40×96 | `982c0374534e8204e0b0c54b3fbea67d09e48d6f451ce8eafacaa48dd8ad3706` | `cb18ff5e6ee33852f71f3eb84d3ee629c7c034644802a673cbeea937e05579b3` |
| x4 `btn0mid.png` | 40×96 → 40×96 | `82f5decccfdcff2703b6c982f65a30d18ee270111c564e63acb4c5edb50f37bd` | `09d040ea025f9da9389564bb714b68abba8e726341b5f5f60a21ad6b3ba6d311` |
| x4 `btn0right.png` | 40×96 → 40×96 | `7276c92ff6eb2847dc2cd98e239ddbc89267143734e32065c625860f346840f2` | `91124c7f534eb6519c94e126dde3f7710c58a981833b66cd508f9d94923ef995` |

`installation-manifest.json` kiểm lại 15/15: SHA thực tế khớp `installed_sha256`.

## Bằng chứng game thật

Ảnh được chụp bằng Game View trong Play Mode, không dựng/mô phỏng:

```text
D:\NRO_Project\Captures\ui-dot4-task3\ui-dot4-task3-game.png
1454x865
sha256=495c2497e65979b23aa41e1cf1afc4db7be9566032f1ce21f6f3e2fdfe41827b
```

Console ngay trong lượt Play:

```text
Retrieved 0 log entries. (error + warning)
```

Đánh giá trực quan từ ảnh:

- PASS popup x2: nền giấy hiển thị liền mạch; không thấy seam giữa tám lát và không có viền
  trắng trên nền cảnh.
- PASS button x2: hai nút `Có`/`Không` dùng đủ left/mid/right, chữ đọc được, không hở mối.
- BLOCKED HUD: popup mất kết nối che màn và không vào được nhân vật; không có evidence thật
  để đánh giá HP/chakra nằm trong rãnh hoặc anchor HUD.
- Menu Iruka không mở được; ảnh popup đã là evidence rõ ràng thay thế cho ba mảnh button.

Không có screenshot giả và không có credential nào được ghi vào report.

## Bảo vệ Small27083/27084

Hash trước và sau giống nhau:

| Tệp | SHA-256 trước = sau |
|---|---|
| x1 `Small27083.png` | `b0775cba09eadf2fde241b35daa28cdcd03af6d23e99e0decd32b42f13f6d220` |
| x1 `Small27084.png` | `0a447e754f64bb1fb3cf809c9cbd22ab9292a06e44b546048c42f8169696c836` |
| x2 `Small27083.png` | `052c48aa8b5e4f2c40bcd0d1bb7a4b93276276285dfa3e8a971b435dc88fdfee` |
| x2 `Small27084.png` | `766cbe1c226e1098bca94ad5f756a0e7059a48df628615698b62605c1a4442df` |
| x3 `Small27083.png` | `e93cf7cf69e59a49eaaa011e5f0ec7b6e385408ba1b4025b9231832d7c87ec46` |
| x3 `Small27084.png` | `7359eaaa4f99f52ab55ef789204ee299285fdf5e230f5963667a6ca4fd1004f1` |
| x4 `Small27083.png` | `1594258d650464372a84556d099c6f0fad9551bbf0efe58e7403336fee183eee` |
| x4 `Small27084.png` | `ee47e805225f2ff19df563f1f1c75298e4ae8511f50abd3e376f773211cf9448` |

## Phục hồi

Backup giữ nguyên cây đường dẫn `Assets/Resources/res/x2..x4/mainimage`. Khi cần rollback,
chỉ sao 15 PNG dưới
`docs/my-thuat/lang-la-client-assets/mainimage-truoc-dot-4/Assets/Resources/` trở lại
`D:\NRO_Project\Assets\Resources\` lúc Unity đang Edit Mode/idle, rồi refresh asset.
`backup-manifest.json` chứa SHA và kích thước bản cũ; `installation-manifest.json` chứa SHA
bản đã lắp. Không cần và không được đụng `Small27083/27084` hay asset-lock khi rollback.
