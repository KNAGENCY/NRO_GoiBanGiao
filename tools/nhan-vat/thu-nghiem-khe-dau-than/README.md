# Thử nghiệm khe đầu–thân nhân vật (05/09/2026) — CHƯA tích hợp vào pipeline

Thư mục này giữ **các script mẫu đã chạy được** do một vòng agent (3 phản biện → 3 thiết kế
→ 2 giám khảo → 1 soát thiếu) viết trong thư mục tạm của phiên, được điều phối chép về đây
để không mất. Chúng **không** nằm trong đường sinh asset chính (`tools/artqa/task6b_closure.py`
→ `tools/nhan-vat/sinh-lai-sprite.py`); muốn dùng thật phải tích hợp theo kế hoạch ở
`BAT-DAU-TU-DAY.md` §5 (hàng "Khe đầu–thân nhân vật"). Bằng chứng ảnh/số ở
`docs/my-thuat/review/khe-dau-than-2026-09-05/`.

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
| `kiem-bang-khe.py` | soát thiếu | Đọc bảng `ghep-khung.py` của nhiều biến thể, đếm cf hở ở x1 và x4 |

## Tái lập số đo biến thể được chọn (V1-siết, ghim đầu −38)

```powershell
$ev = "docs\my-thuat\review\khe-dau-than-2026-09-05"
python tools\nhan-vat\ghep-khung.py --part-bin $ev\part-V1-siet-ghim38.bin --sprites <thư mục Small*.png x1 sinh bởi sinh-manh-v1.py> --out $ev\kiem.png
```

Kết quả phải là 33/33 cf: khe đầu–thân ∈ [−3, −1], thân–chân ∈ [−11, −2], bàn chân −11
(cf 0–25) / −10 (cf 26–32); ở `--zoom 4` không cf nào có khe dương. (Ảnh x1 của V1 không
chép vào repo vì sinh lại được bằng `sinh-manh-v1.py`; bảng đã đo lưu ở
`kiem-V1-siet-ghim38-x1.png`.)

## Giá phải trả của phương án dữ liệu (chủ dự án quyết)

- Nhân vật **thấp đi ~10 % khi chạy** (53–55 px so 61–62 px khi đứng) vì bộ xương client
  giữ đầu gần cố định còn khung chạy nguồn vẽ khom.
- Thân **trượt so với chân ≤ 5 px** ở 1/5 vòng chạy (cf 3) và vài khung kỹ năng (cf 14/16).
- Đầu idle hạ 1 px (dy −39 → −38) → `ch` trong `Char.cs:6087` đổi 5 → 6 px.

Muốn hết hai giá đầu phải **vẽ lại 5 khung chạy đứng thẳng** (spec số trong
`docs/my-thuat/review/khe-dau-than-2026-09-05/GHI-CHU.md`) hoặc đổi bảng xương client.
