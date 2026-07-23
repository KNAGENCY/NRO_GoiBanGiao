# PROMPT KEY — Khoá nhận diện thị giác Ngọc Rồng Kame

Dán **nguyên khối này vào đầu mọi prompt sinh ảnh**, rồi mới nối phần mô tả từng khối phía sau.

Mục đích: web và game phải trông như **cùng một sản phẩm**, không phải hai thứ rời nhau.

---

## Cơ sở: bảng màu đo từ chính tài nguyên game

Không phải màu tự nghĩ. Đo bằng cách lấy mẫu điểm ảnh từ:

| Nguồn | Màu trội đo được |
|---|---|
| `Resources\res\x4\IconChucNang\bgshowicon.png` (nền bảng giao diện) | `#481800` 22% · `#301800` 19% · `#603000` 16% · `#300000` 10% · `#D80000` 6% |
| `Resources\res\x4\Logo\1.png` (logo game) | `#600018` 27% · `#F09060` 19% · `#F06030` 13% · `#60C060` 10% |
| `data\item_bg_temp\x4\0.png` (cảnh vật) | `#00A818` · `#00C030` · `#007800` |

**Kết luận: game có bản sắc GIẤY CŨ SEPIA + CAM ĐỎ + XANH LÁ, phong cách cuộn giấy phương Đông vẽ tay.**
Không phải tối lạnh hiện đại. Không phải neon. Không phải tối giản Bắc Âu.

---

## Bảng màu web chính thức

Giữ đúng sắc độ của game, nhưng dựng thành thang dùng được cho web.

| Vai trò | Mã | Gốc từ game |
|---|---|---|
| Nền sâu nhất | `#1A0F08` | tối hơn của `#301800` |
| Nền chính | `#2E1A0C` | `#301800` |
| Bề mặt nổi | `#4A2A12` | `#603000` |
| Viền mảnh | `#6B4522` | trung gian |
| **Giấy cũ** (khối nội dung sáng) | `#E8D9BC` | nền cuộn giấy |
| Giấy cũ đậm | `#C9B48C` | bóng trên giấy |
| Chữ trên nền tối | `#F2E6CE` | kem |
| Chữ phụ | `#A8916E` | nâu nhạt |
| **Nhấn chính — nút bấm** | `#F06030` | cam logo |
| Nhấn sáng | `#F09060` | cam nhạt logo |
| Viền / bóng chữ | `#600018` | mận viền logo |
| **Cảnh báo · Boss sống** | `#D80000` | đỏ viền giao diện |
| Thành công · Xanh | `#60C060` | xanh chữ ONLINE |
| Xanh cảnh vật | `#00A818` | dùng rất ít, cho minh hoạ |

**Quy tắc dùng màu:** nền tối nâu là chủ đạo → khối nội dung dùng giấy cũ `#E8D9BC` → nút bấm cam `#F06030` → đỏ `#D80000` **chỉ** cho trạng thái sống/khẩn.

---

## Chữ

Logo game là **chữ tròn mập, viền dày, đổ bóng** — web phải nối tiếp tinh thần đó ở tiêu đề, nhưng nội dung vẫn phải dễ đọc.

| Vai trò | Font | Vì sao |
|---|---|---|
| Tiêu đề lớn | **Baloo 2** ExtraBold | tròn mập giống logo, **có đủ dấu tiếng Việt** |
| Nội dung | **Be Vietnam Pro** | thiết kế riêng cho tiếng Việt |
| Số liệu | **Be Vietnam Pro** hoặc mono | căn cột thẳng |

⚠️ **Tuyệt đối tránh**: Anton, Bebas Neue, Monument, Oswald bản cũ — **không có dấu tiếng Việt**, chữ vỡ thành ô vuông.

---

## PROMPT KEY — dán khối này vào đầu mọi prompt

```
=== VISUAL LOCK: NGỌC RỒNG KAME ===

Website design comp, horizontal 16:9, for a Vietnamese Dragon Ball MMO portal
called "Ngọc Rồng Kame". The site must feel like it belongs to the SAME WORLD
as the game client.

WORLD & MATERIAL (this is the identity — do not drift):
- Aged oriental scroll aesthetic: weathered parchment, sepia-stained paper,
  soft ink wash, faint halftone dots, subtle ink splatter, oriental floral
  and cloud motifs used sparingly as edge ornament.
- Hand-painted 2D anime illustration, saturated but warm, soft airbrush
  gradients, dark confident outlines. NOT pixel art. NOT 3D render.
  NOT flat vector.
- Deep warm brown environment with glowing warm light, like lantern light
  on old paper.

EXACT PALETTE — use only these:
- Deepest ground   #1A0F08
- Page ground      #2E1A0C
- Raised surface   #4A2A12
- Hairline         #6B4522
- Aged parchment   #E8D9BC   (content panels sit on this)
- Parchment shade  #C9B48C
- Text on dark     #F2E6CE
- Muted text       #A8916E
- PRIMARY ACTION   #F06030   (orange — buttons, key emphasis)
- Highlight        #F09060
- Outline / shadow #600018   (deep maroon)
- Live / alert     #D80000   (only for "boss alive" and urgent state)
- Success green    #60C060

TYPOGRAPHY:
- Headings: heavy rounded display with thick dark maroon outline and a soft
  drop shadow, echoing a cartoon game logo (Baloo 2 ExtraBold character).
- Body: clean humanist sans, highly readable (Be Vietnam Pro character).
- FULL Vietnamese diacritics required on every glyph.
- Strong size contrast. Short lines. Never all-caps paragraphs.

COMPONENT LANGUAGE:
- Content panels look like parchment cards with a thin dark maroon border
  and a soft inner shadow, 12px radius.
- Primary buttons: solid orange #F06030, dark maroon 2px outline, 8px radius,
  slight bevel — like a game button, not a flat SaaS pill.
- Secondary buttons: parchment fill, maroon outline, no glow.
- Dividers: thin ink-drawn hairlines, not crisp 1px tech lines.

STRICTLY BANNED:
- purple/blue AI gradients, neon glow, glassmorphism, cyberpunk
- cold grey minimalism, Scandinavian flat, corporate SaaS look
- pixel art, voxel, 3D render
- rainbow color, gradient text
- marquee strips, fake KPI blocks, countdown timers

Layout must stay clean, breathable, and codeable — generous spacing, obvious
hierarchy, one clear primary action per section.

=== END VISUAL LOCK ===
```

---

## Cách kiểm tra ảnh sinh ra có đúng khoá không

Sau khi có 7 ảnh, soi 5 điểm:

1. **Nền có phải nâu ấm không?** Nếu ra xám lạnh hoặc đen tuyền → sai, sinh lại.
2. **Nút chính có phải cam `#F06030` không?** Nếu ra xanh dương hoặc tím → sai.
3. **Khối nội dung có cảm giác giấy cũ không?** Nếu là thẻ trắng bo góc kiểu SaaS → sai.
4. **Tiêu đề có mập, có viền không?** Nếu mảnh và phẳng → lệch khỏi logo game.
5. **Đặt cạnh ảnh `bgshowicon.png` của game — có thấy cùng một thế giới không?** Đây là phép thử cuối cùng.

Nếu 7 ảnh không nhất quán với nhau, sinh lại **cả bộ** thay vì vá từng cái — lệch màu giữa các khối là lỗi nặng hơn một khối xấu.

---

## Ghi chú cho việc code sau này

Khi dựng CSS, khai báo bảng màu này thành biến ngay từ đầu:

```css
:root{
  --ground-deep:#1A0F08;  --ground:#2E1A0C;    --surface:#4A2A12;
  --hairline:#6B4522;     --parchment:#E8D9BC; --parchment-2:#C9B48C;
  --text:#F2E6CE;         --text-muted:#A8916E;
  --action:#F06030;       --action-hi:#F09060; --outline:#600018;
  --alert:#D80000;        --success:#60C060;
  --r-card:12px;          --r-btn:8px;
}
```

Đổi thương hiệu sau này chỉ cần sửa khối này, không phải lục từng file.
