# Thử nghiệm khe đầu–thân nhân vật (05/09/2026) — ĐÃ TÍCH HỢP 05–06/09/2026 (V1 người chơi + Iruka); thư mục này chỉ còn là lịch sử

Thư mục này giữ **các script mẫu đã chạy được** do một vòng agent (3 phản biện → 3 thiết kế
→ 2 giám khảo → 1 soát thiếu) viết trong thư mục tạm của phiên, được điều phối chép về đây
để không mất. **Đường sinh asset chính nay là** `tools/artqa/task6b_closure.py` (phần V1) +
`tools/artqa/nhan_vat_moc.py` (mốc giải phẫu) + `tools/artqa/nhan_vat_toi_uu.py` (MILP dx/dy) →
`tools/nhan-vat/sinh-lai-sprite.py` (một lệnh sinh 37 mảnh người chơi + Iruka, part.bin, mapping,
lock). Các script ở đây **không còn dùng để sinh asset**; giữ để tra lại cách chẩn đoán. Bằng chứng
ảnh/số ở `docs/my-thuat/review/khe-dau-than-2026-09-05/`.

## Chẩn đoán (đã tự kiểm lại bằng `tools/nhan-vat/ghep-khung.py`)

"Đầu tách khỏi thân" thực chất là **hai khuôn mặt chồng nhau**: mảnh THÂN bị cắt ở 28 %
chiều cao khung nguồn; ở khung chạy/đánh (vẽ nhỏ hơn idle 12–27 % và tư thế thấp hơn) lát
cắt ấy nằm trên trán, nên mảnh thân chứa cả mắt–mũi–miệng và client vẽ nó **đè lên** đầu
idle (thứ tự vẽ đầu → chân → thân). Cộng thêm: khung NẰM (hit_death_getup 1) bị gán vào
slot thân 15 (cf 20/31 — client dùng cho thế bắt đầu skill, không bao giờ vẽ nằm), và cf 2
là khung CHẠY đầu tiên (`Char.cs:2537`, `cf = gameTick % 5 + 2`) nhưng mapping gán thân idle.

## File

| File | Nguồn gốc | Làm gì |
|---|---|---|
| `do-strip.py` | điều phối | Đo cao khung + bề rộng đầu 17 khung nguồn (số liệu đầu tiên chỉ ra khác tỉ lệ) |
| `do-moc-giai-phau.py` | phản biện ảnh | Tách màu da → cao mặt, mắt, mắt→cằm, băng đô từng khung; hệ số tỉ lệ theo dải (run ×1,13, melee ×1,20, cast ×1,03) |
| `phan-tich-cua-so-dy.py` | phản biện toán | Chứng minh chỉ đổi 3 dy đầu là vô nghiệm; 17 dy thân mới đủ bậc tự do (MILP) |
| `sinh-manh-v1.py` | thiết kế V1 | Sinh 34 mảnh x1/x4 theo V1: tỉ lệ từng khung theo cỡ đầu, cắt thân tại cằm−2 px, gán lại slot (cf 2 = run0, slot thân 15 = cast0, slot 16 = idle2) → `geometry.json` |
| `toi-uu-dxdy.py` | thiết kế V1 + giám khảo 1 + điều phối | MILP (scipy) chọn dx/dy nguyên cho 34 mảnh; `KL_MAX=-2` siết khe thân–chân; `PIN_HEAD0=-38` ghim dy đầu slot 0 để cf 0 chỉ đổi 1 px |
| `chuan-hoa-dai-thu.py` | thiết kế ảnh | Thử chuẩn hoá tỉ lệ dải + lát cắt giải phẫu; đo phần thiếu hụt do TƯ THẾ (run thiếu 5–12 px x1) |
| `sinh-charinfo-phuong-an-client.py` | thiết kế client | Phương án dự phòng V2: sinh bảng `CharInfo` mới cho `Char.cs` — **không dùng** khi spec còn cấm sửa client |
| ~~`kiem-bang-khe.py`~~ | soát thiếu | **Đã xoá 05/09 (chiều)** — việc "đọc bảng, đếm cf hở" nay là cờ `--kiem` của chính `tools/nhan-vat/ghep-khung.py` (thoát mã ≠ 0 + in cf vi phạm, có thêm cột `so mat` bắt HAI KHUÔN MẶT mà khe bbox âm không lộ). Script cũ trỏ cứng vào thư mục tạm của phiên đã mất, không chạy lại được |

## Tái lập số đo biến thể được chọn (V1-siết, ghim đầu −38)

```powershell
$ev  = "docs\my-thuat\review\khe-dau-than-2026-09-05"
$tam = "<thư mục tạm>"
# 1. sinh lại 34 mảnh V1 (x1 ở $tam\V1\Small*.png, x4 ở $tam\V1\x4\) — đúng cấu hình của
#    geometry-V1.json: slot thân 14 = cast2, 15 = cast0 (mặc định trong script là hit0/cast1)
python tools\nhan-vat\thu-nghiem-khe-dau-than\sinh-manh-v1.py --ten V1 --gan "14=cast2,15=cast0" --out $tam
# 2. đo + cửa kiểm (thoát mã 0 = đạt; ≠ 0 thì in cf vi phạm)
python tools\nhan-vat\ghep-khung.py --part-bin $ev\part-V1-siet-ghim38.bin --sprites $tam\V1 --out $ev\kiem.png --kiem "dau-than=-9,0;than-chan=-11,-2;so-mat=1;chan=-11,-10"
python tools\nhan-vat\ghep-khung.py --zoom 4 --part-bin $ev\part-V1-siet-ghim38.bin --sprites $tam\V1\x4 --out $ev\kiem-x4.png --kiem "dau-than=-36,0;than-chan=-44,0;so-mat=1"
```

Kết quả phải là 33/33 cf: khe đầu–thân ∈ [−3, −1], thân–chân ∈ [−11, −2], bàn chân −11
(cf 0–25) / −10 (cf 26–32), **`so mat` = 1**; ở `--zoom 4` không cf nào có khe dương và
cũng 33/33 một mặt. Đối chứng hiện trạng client 05/09 (chạy `ghep-khung.py` không tham số
`--sprites`): 23/33 cf có `so mat` = 2 ở x1 (cf 3–13, 16–19, 21–22, 25–26, 29–30, 32 — gồm
cf 4/5/7/12/17 mà `so-sanh-ghep.png` chỉ ra), x4 còn bắt thêm cf 20/23/24/31.

Ghi chú tái lập: 34 mảnh sinh lại ngày 05/09 (chiều) trùng sha 34/34 (x1 lẫn x4) với bản
đã đo. `sinh-manh-v1.py` đọc số đo mốc giải phẫu qua hằng `DO_JSON` trỏ vào thư mục tạm
của phiên; file ấy trùng sha với `moc-giai-phau-do.json` trong `$ev` — thư mục tạm mất
thì trỏ `DO_JSON` sang đó. (Ảnh x1 của V1 không chép vào repo vì sinh lại được; bảng đã đo
lưu ở `kiem-V1-siet-ghim38-x1.png`.)

## Giá phải trả của phương án dữ liệu (chủ dự án quyết)

- Nhân vật **thấp đi ~10 % khi chạy** (53–55 px so 61–62 px khi đứng) vì bộ xương client
  giữ đầu gần cố định còn khung chạy nguồn vẽ khom.
- Thân **trượt so với chân ≤ 5 px** ở 1/5 vòng chạy (cf 3) và vài khung kỹ năng (cf 14/16).
- Đầu idle hạ 1 px (dy −39 → −38) → `ch` trong `Char.cs:6087` đổi 5 → 6 px.

Muốn hết hai giá đầu phải **vẽ lại 5 khung chạy đứng thẳng** (spec số trong
`docs/my-thuat/review/khe-dau-than-2026-09-05/GHI-CHU.md`) hoặc đổi bảng xương client.
