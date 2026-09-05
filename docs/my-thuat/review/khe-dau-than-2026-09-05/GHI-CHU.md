# Khe đầu–thân nhân vật — bằng chứng chẩn đoán và phương án (05/09/2026, chiều)

Câu hỏi của chủ dự án: *"tại sao không fix được cho nhân vật hoàn chỉnh, có cần cài thêm
công cụ gì không?"* Trả lời ngắn: **fix được hoàn toàn bằng dữ liệu, không cần cài gì**;
cái chưa làm là vì trước đó chẩn đoán sai (tưởng chỉ là khe, hoá ra là hai khuôn mặt) và
vì có một giá thẩm mỹ cần anh quyết (nhân vật thấp đi ~10 % khi chạy).

## 1. Nhìn bằng mắt

| Ảnh | Nội dung |
|---|---|
| `so-sanh-ghep.png` | 4 hàng × 9 cf. Hàng 1 hiện trạng: cf 4/5/7/12/17 có **mặt thứ hai trên ngực**, cf 20 thân nằm. Hàng 4: cắt thân theo cằm → hết hai mặt, lộ khe thật đúng bằng thiếu hụt tư thế. |
| `zoom-co-x6.png` | Vùng cổ ×6, cf 0/3/4/7/12/16/24: hàng trên hiện trạng (hai mặt), hàng dưới V1 (đầu ngồi lên vai). |
| `vong-chay-truoc-sau.gif` / `.png` | Vòng chạy cf 2→6 đúng thứ tự client, trước (trái) / sau V1-siết (phải), zoom ×4. Thấy rõ: sau sửa **thấp hơn ~10 %** và thân trượt nhẹ ở cf 3. |
| `sprite-than-truoc-sau.png` | Bốn mảnh thân 27048/27055/27058/27053 cũ (có mặt) và mới (từ cằm xuống). |
| `iruka-cf0-x5.png`, `kiem-iruka.png` | **NPC Iruka cũng hỏng, nặng hơn**: 3 mảnh 2189/2190/2191 có dx = dy = 0 → đầu lửng +13 px, thân chìm dưới đất +29 px. Rất có thể đây là cái anh thấy "đầu tách khỏi thân" khi đứng cạnh Iruka. |
| `ghep-now.png` | 33 cf hiện trạng (x1), bảng đo trong `BAT-DAU-TU-DAY.md`. |
| `kiem-V1-siet-x1.png`, `kiem-V1-siet-x4.png`, `kiem-V1-siet-ghim38-x1.png` | 33 cf sau sửa, điều phối tự chạy lại `ghep-khung.py`. |

## 2. Nguyên nhân gốc — ba tầng lệch, đã đo

1. **Ảnh nguồn khác tỉ lệ.** Dải chạy vẽ nhân vật nhỏ hơn idle ~12 % (0,86–0,90), dải đánh
   ~17 % (0,78–0,86); đo bằng 6 thước cứng của đầu (cao/rộng mặt, mắt, mắt→cằm, băng đô)
   trong `moc-giai-phau-do.json`. Dải chạy còn khom thêm ~13 %, dải đánh chỉ ~6 %.
2. **Lát cắt theo phần trăm chiều cao khung** (đầu 0–42 %, thân 28–76 %) đúng với idle
   nhưng ở khung thấp/nhỏ thì 28 % rơi lên trán → mảnh thân **chứa mặt**. Thân vẽ sau
   đầu nên mặt ấy đè lên đầu idle. Khe bbox âm vẫn có thể là hai mặt — số đo cũ không
   bắt được.
3. **Bộ xương cứng trong client** (`Char.cs:483`, 33 cf, 3 đầu dùng chung, một dx/dy cho
   mỗi mảnh) đặt đầu gần cố định (dy 32–34 cho cf 2–6) trong khi khung chạy nguồn thấp
   hơn 5–12 px x1. Chỉ đổi 3 dy đầu là **vô nghiệm** (cửa sổ 9 px, biên độ cần 12–14);
   nhưng 17 dy thân (mỗi slot chỉ phục vụ 1–3 cf) đủ bậc tự do → giải bằng MILP.

Ngoài ra: khung NẰM bị gán vào slot thân 15 (cf 20/31 = thế bắt đầu skill, client không
bao giờ vẽ nằm — chết = đứng + vòng sáng); cf 2 là khung chạy đầu tiên nhưng mapping cho
thân idle; slot đầu 1 về ngữ nghĩa là đầu NGANG (tập `isFrNgang`) nhưng pipeline nạp idle
chính diện.

## 3. Phương án được hai giám khảo chọn: V1 (chỉ đổi dữ liệu)

- Nhân tỉ lệ **từng khung** theo cỡ đầu (run ×1,12–1,17, melee ×1,17–1,27, cast ×1,00–1,07).
- Cắt đầu tới cằm + 6 px nguồn; cắt **thân từ cằm − 2 px x1** (không còn mặt); chân giữ.
- Gán lại slot: thân/chân 2..6 = run0..3,run0; slot thân 14 = cast2; 15 = cast0; 16 + chân 13 = idle2.
- MILP (`scipy.optimize.milp`) chọn dx/dy nguyên: dy chân cố định (bàn chân −11 giữ
  nguyên), khe đầu–thân ∈ [−9, 0], thân–chân ∈ [−11, **−2**] (siết của giám khảo 1 để x4
  không hở), ghim dy đầu slot 0 = **−38** (điều phối thử: −39 vô nghiệm, −38 có nghiệm).
- Kết quả tự đo lại: 33/33 cf khe đầu–thân ∈ [−3, −1], thân–chân ∈ [−11, −2], bàn chân
  −11/−10 như cũ; x4 không cf nào hở. File: `part-V1-siet-ghim38.bin` (+ `.dxdy.json`),
  `geometry-V1.json`.
- **Giá:** nhân vật thấp đi ~10 % khi chạy (53–55 so 61–62 px), thân trượt ≤ 5 px ở cf 3/14/16,
  đầu idle hạ 1 px (`ch` 5 → 6 px). Không sửa client, không cài gì, quay lui bằng git.

Dự phòng: **V2** sinh lại `CharInfo` trong client (`charinfo-B-phuong-an-client.json`) —
hình học sạch nhất nhưng sửa 2 file `Char.cs` (Game1/Game2), client không có git, tên/HP
nhấp 14–20 px, trái spec "client giữ nguyên". **V3-B** vẽ lại 5 khung chạy đứng thẳng
(spec: cùng nhân vật idle_0; ở thang nguồn 655 px: cao mặt 96–99, mắt 33–34, băng đô 25–26,
cằm→bàn chân 375–395 px, tóc→bàn chân 640–670, ít nhất một bàn chân chạm mép dưới, nền
chroma (246,4,242), đầu 3/4 như idle) — chỉ khi anh không chấp nhận "thấp đi khi chạy".

## 4. Công cụ

Đường V1 **không cần cài gì**: Python 3.14 + Pillow + numpy + scipy (đã có), Unity MCP
đang chạy, Photoshop (Beta) đã cài để vá tay nếu cần. Nếu chọn vẽ lại: dùng MCP sinh ảnh
đã nối (Higgsfield: `nano_banana_pro` có ảnh tham chiếu, model `autosprite`; còn ~398
credit) — không cài gì; chỉ khi thất bại mới cân nhắc ComfyUI + Qwen-Image-Edit (~30 GB,
RTX 4070 chạy được). Không khuyến nghị PixelLab/Aseprite/ImageMagick/GIMP cho việc này.

## 5. Còn chưa kiểm (nói thẳng)

- Chưa có ảnh **trong game** nào của biến thể mới (mọi số đo qua `ghep-khung.py`, bản mô
  phỏng `Char.cs:6057–6084`); hướng nhìn trái (lật gương) chưa đo.
- `ch` = 5 px (`Char.cs:6087`) làm số sát thương/chat neo gần chân — suy từ mã, chưa nhìn game.
- Iruka: số đo +13/+29 là từ `ghep-khung.py`; nguồn 10 sprite 27075–27084 trong client
  không khớp file nào trong repo.
