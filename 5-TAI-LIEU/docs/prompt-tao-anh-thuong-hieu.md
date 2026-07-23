# Bộ prompt tạo ảnh thương hiệu — Ngọc Rồng Kame

Danh sách đầy đủ ảnh cần làm, prompt cho từng cái, kích thước thật, và **quy tắc đặt tên**.

---

# 1. QUY TẮC ĐẶT TÊN FILE

**Bắt buộc theo đúng mẫu này.** Nhìn tên là tôi biết ảnh đi đâu, khỏi phải hỏi lại.

```
kame_<nhóm>_<tên>_<rộng>x<cao>_v<số>.<đuôi>
```

| Thành phần | Giá trị hợp lệ |
|---|---|
| `<nhóm>` | `logo` · `web` · `social` · `game` · `icon` · `tex` |
| `<tên>` | chữ thường, không dấu, nối bằng `-` |
| `<rộng>x<cao>` | đúng số điểm ảnh thật của file |
| `v<số>` | `v1`, `v2`... tăng dần mỗi lần sửa |
| `<đuôi>` | `png` nếu **cần nền trong suốt**, `jpg` nếu nền đặc |

**Ví dụ đúng:**
```
kame_logo_ngang_599x136_v1.png
kame_logo_vuong_1024x1024_v1.png
kame_web_hero_1920x1080_v1.jpg
kame_web_og_1200x630_v1.jpg
kame_social_fbcover_1640x856_v1.jpg
kame_game_panel_1000x650_v1.png
kame_icon_boss_110x110_v1.png
kame_tex_giay-cu_1024x1024_v1.jpg
```

**Sai — đừng gửi kiểu này:** `logo final.png`, `anh1.jpg`, `hero_moi_sua_lan_cuoi.png`

Nếu làm nhiều phương án cho cùng một chỗ, dùng `v1`, `v2`, `v3` rồi nói tôi biết bạn chọn cái nào.

---

# 2. MASTER PROMPT ⭐ — dùng cái này

**Đây là thứ chính cần dùng.** Một prompt duy nhất, chừa một ô trống. Dán một lần, mỗi lần chạy chỉ thay đúng 3 dòng trong ô `TASK`.

> ⚠️ Không công cụ nào sinh được 15 ảnh khác kích thước trong một lần chạy. Mỗi lần chạy ra **một ảnh**. Chạy 12 lần với 12 khối TASK bên dưới.

```
=== NGỌC RỒNG KAME — MASTER ASSET PROMPT ===

Brand asset artwork for a Vietnamese Dragon Ball MMO called "Ngọc Rồng Kame".
Must belong to the SAME WORLD as the existing game client art.

WORLD & STYLE (never change this):
- Aged oriental scroll aesthetic: weathered sepia parchment, soft ink wash,
  faint halftone dots, sparse ink splatter, subtle oriental cloud and floral
  motifs at the edges.
- Hand-painted 2D anime illustration. Warm, saturated, soft airbrush
  gradients, confident dark outlines. Lantern-light mood on old paper.
- NOT pixel art. NOT 3D render. NOT flat vector. NOT photorealistic.

EXACT PALETTE (use only these):
  #1A0F08 deepest ground   #2E1A0C page ground    #4A2A12 raised surface
  #6B4522 hairline         #E8D9BC aged parchment #C9B48C parchment shade
  #F2E6CE light text       #A8916E muted text
  #F06030 PRIMARY ORANGE   #F09060 highlight      #600018 deep maroon outline
  #D80000 alert red        #60C060 green

BANNED (reject if present):
- purple/blue AI gradients, neon glow, glassmorphism, cyberpunk
- cold grey minimalism, corporate flat design, Scandinavian look
- rainbow color, gradient text as an effect
- watermarks, signatures, stock-photo look
- any real trademarked Dragon Ball character face or official logo
- text in the image unless the task line explicitly asks for it

>>> TASK — replace this single block each run: <<<

SIZE: [KÍCH THƯỚC]
BACKGROUND: [NỀN]
SUBJECT: [NỘI DUNG]

>>> end task block <<<

Render exactly one asset. Clean edges, correct aspect ratio, no borders or
padding added beyond what the task asks for.
```

## 12 khối TASK — dán lần lượt vào ô trống

### 1 → `kame_logo_ngang_599x136_v1.png`
```
SIZE: 599x136 pixels, horizontal
BACKGROUND: transparent
SUBJECT: Game logo wordmark reading "NGỌC RỒNG KAME". Chunky rounded cartoon
game-logo lettering, thick #600018 maroon outline, soft drop shadow, letters
filled with a warm vertical gradient #F09060 to #F06030. Optional small amber
dragon ball orb tucked beside the text. Full Vietnamese diacritics rendered
correctly. Bold and readable, clean edges.
```

### 2 → `kame_logo_vuong_1024x1024_v1.png`
```
SIZE: 1024x1024 pixels, square
BACKGROUND: transparent
SUBJECT: A single bold emblem, no long text — a stylized amber dragon ball orb
with four red stars, wrapped by a small coiled green oriental dragon. Thick
#600018 maroon outline, #F06030 to #F09060 gradient fill, subtle parchment
texture inside the shape only. Big simple shapes with a strong silhouette that
still reads at 32x32 pixels.
```

### 3 → `kame_web_hero_1920x1080_v1.jpg`
```
SIZE: 1920x1080 pixels, horizontal
BACKGROUND: full-bleed painted scene
SUBJECT: Vast warm golden-hour sky over floating rocky islands, a distant
oriental temple silhouette, one glowing amber dragon ball orb in the sky.
Keep the LEFT THIRD and BOTTOM THIRD calm and low-detail — text will overlay
there. Visual interest in the upper right. Darkens naturally toward the
bottom-left. No faces, no text.
```

### 4 → `kame_web_og_1200x630_v1.jpg`
```
SIZE: 1200x630 pixels, horizontal
BACKGROUND: aged parchment over warm dark brown
SUBJECT: The Ngọc Rồng Kame logo centered with a soft amber glow behind it,
faint oriental cloud motifs in the corners, and a thin #D80000 red border inset
a few pixels from the edge like a game UI panel. Small clean text
"nrokame.com" along the bottom in #F2E6CE. High contrast, readable as a small
thumbnail.
```

### 5 → `kame_tex_giay-cu_1024x1024_v1.jpg`
```
SIZE: 1024x1024 pixels, square, SEAMLESSLY TILEABLE
BACKGROUND: fills the frame
SUBJECT: Aged parchment paper — warm #E8D9BC base, irregular sepia staining,
subtle fiber grain, faint water marks, a few light ink speckles. Soft and even,
no focal point, no vignette, no dark corners. Must tile with no visible seam.
Low contrast so text stays readable on top.
```

### 6 → `kame_social_fbcover_1640x856_v1.jpg`
```
SIZE: 1640x856 pixels, horizontal
BACKGROUND: painted scene fading to dark brown
SUBJECT: Logo lockup in the LEFT HALF. Right half shows warm painterly scenery
— floating islands and an amber orb — fading into #2E1A0C. Keep the bottom-left
circular area about 180x180 pixels clear (profile picture overlaps there) and
keep the outer 100 pixels on all sides free of critical detail.
```

### 7 → `kame_game_panel_1000x650_v1.png`
```
SIZE: 1000x650 pixels, horizontal
BACKGROUND: fills the frame
SUBJECT: Game UI panel — aged parchment with a bold #D80000 red rounded border
frame. Inside: warm sepia paper, faint oriental cloud and floral motifs, light
ink splatter, and a large coiled oriental dragon in muted olive green and tan
sitting BEHIND everything at low contrast like a watermark. Center stays calm
and readable; darkest at the edges, lighter toward the middle.
```

### 8, 9, 10 → `kame_icon_boss_110x110_v1.png` · `kame_icon_tuido_110x110_v1.png` · `kame_icon_vongquay_110x110_v1.png`
```
SIZE: 110x110 pixels, square
BACKGROUND: transparent
SUBJECT: A single game UI icon of [ĐỔI: a red demon head | a treasure chest |
a spinning prize wheel]. Chunky hand-painted 2D game icon, thick #600018
maroon outline, warm #F06030 and #F09060 fills with soft airbrush shading.
Bold simple silhouette readable at 40x40 pixels. Centered with a few pixels of
padding. No text.
```
Chạy 3 lần, đổi phần trong ngoặc.

### 11, 12 → `kame_web_band-boss_1920x600_v1.jpg` · `kame_web_band-congdong_1920x600_v1.jpg`
```
SIZE: 1920x600 pixels, wide horizontal
BACKGROUND: full-bleed, heavily darkened
SUBJECT: [ĐỔI: a dark cavern lit by amber torchlight | a serene mountain
training ground at dawn]. Low contrast, no focal subject in the center — text
will sit on top. Warm brown and amber grade only.
```

## Hai ảnh KHÔNG cần sinh

- **Icon ứng dụng** `479x480` → lấy ảnh số 2, đặt lên nền `#2E1A0C`
- **Avatar Zalo/Facebook** → cũng lấy ảnh số 2

Sinh riêng chỉ làm hai thứ lệch nhau.

---

# 3. ASSET LOCK — bản chi tiết (tham khảo)

Phần bên dưới là bản **tách rời từng ảnh**, viết trước khi có Master Prompt ở trên. Giữ lại để tra cứu ngữ cảnh và giải thích từng ảnh dùng ở đâu. **Để tạo ảnh thì dùng Master Prompt ở mục 2.**

Khác với `prompt-key-khoa-visual.md` (dùng cho **ảnh mockup website**), khối này dùng cho **ảnh đơn lẻ**: logo, icon, banner, hoạ tiết.

```
=== ASSET LOCK: NGỌC RỒNG KAME ===

Brand asset artwork for a Vietnamese Dragon Ball MMO called "Ngọc Rồng Kame".
Must belong to the SAME WORLD as the existing game client art.

STYLE:
- Aged oriental scroll world: weathered sepia parchment, soft ink wash,
  faint halftone dots, sparse ink splatter, subtle oriental cloud and floral
  motifs at the edges.
- Hand-painted 2D anime illustration. Warm, saturated, soft airbrush
  gradients, confident dark outlines.
- NOT pixel art. NOT 3D render. NOT flat vector. NOT photorealistic.

EXACT PALETTE — use only these:
  #1A0F08 deepest ground   #2E1A0C page ground    #4A2A12 raised surface
  #6B4522 hairline         #E8D9BC aged parchment #C9B48C parchment shade
  #F2E6CE light text       #A8916E muted text
  #F06030 PRIMARY ORANGE   #F09060 highlight      #600018 deep maroon outline
  #D80000 alert red        #60C060 green

BANNED:
- purple/blue AI gradients, neon glow, glassmorphism, cyberpunk
- cold grey minimalism, corporate flat design
- rainbow color, gradient text as an effect
- watermarks, signatures, stock-photo look
- any real trademarked Dragon Ball character face or official logo

=== END ASSET LOCK ===
```

⚠️ **Dòng cuối quan trọng**: đừng sinh mặt nhân vật Dragon Ball chính thức (Goku, Vegeta...). Dùng hình bóng, dáng người chung chung, hoặc sinh vật/cảnh vật. Bạn đã ở vùng xám bản quyền rồi, đừng làm nặng thêm.

---

# 4. DANH SÁCH ẢNH CẦN LÀM

Xếp theo **mức cần gấp**.

## 🔴 Nhóm 1 — Không có là không mở server được

### 1.1 Logo ngang

- **Tên file:** `kame_logo_ngang_599x136_v1.png`
- **Kích thước:** `599 x 136` — **đúng bằng logo cũ**, để thả thẳng vào client
- **Nền:** trong suốt (PNG)
- **Sẽ đặt vào:** `Assets\Logo\1.png` + `Assets\Resources\res\x4\Logo\1.png` (cần build lại client) và dùng trên web

```
[ASSET LOCK]

Horizontal game logo wordmark, transparent background, 599x136 pixels.

Text: two lines — "NGỌC RỒNG" large on top, "KAME" smaller beneath, or
"NGỌC RỒNG KAME" on one line if it reads better at this width.

Style: chunky rounded cartoon game-logo lettering with thick #600018 deep
maroon outline and a soft drop shadow. Main letters filled with a warm
vertical gradient from #F09060 down to #F06030. Optional small accent: a
single amber dragon ball orb with stars, or a stylized turtle-shell mark,
tucked beside the text — small, not dominating.

Bold, readable, friendly. Full Vietnamese diacritics rendered correctly
(Ọ with dot below, Ồ with circumflex and grave). Clean edges, no texture
noise on the letters themselves.
```

> ⚠️ Dấu tiếng Việt trên logo là chỗ AI hay làm hỏng nhất. Nếu sinh 5 lần vẫn sai dấu, **để tôi ghép chữ bằng code** — bạn chỉ cần gửi phần hoạ tiết/biểu tượng.

### 1.2 Logo vuông (biểu tượng)

- **Tên file:** `kame_logo_vuong_1024x1024_v1.png`
- **Kích thước:** `1024 x 1024`
- **Nền:** trong suốt
- **Dùng cho:** icon ứng dụng, favicon web, avatar Zalo/Facebook

```
[ASSET LOCK]

Square app icon mark, 1024x1024, transparent background.

A single bold emblem — no long text. Choose ONE concept and execute it
cleanly:
(a) a stylized amber dragon ball orb with four red stars, wrapped by a small
    coiled green oriental dragon;
(b) a turtle-shell "Kame" roundel in aged parchment tones with a red seal
    border;
(c) a bold "K" monogram built from dragon-scale forms.

Thick #600018 maroon outline, warm #F06030 to #F09060 gradient fills, subtle
parchment texture inside the shape only. Must stay legible when shrunk to
32x32 pixels — big shapes, no fine detail, strong silhouette.
```

### 1.3 Icon ứng dụng

- **Tên file:** `kame_logo_appicon_480x480_v1.jpg`
- **Kích thước:** `480 x 480` — thay `Assets\Resources\appicon.jpg`
- **Nền:** đặc (JPG)

> Lấy thẳng từ `kame_logo_vuong` đặt lên nền `#2E1A0C` là xong, không cần sinh riêng.

---

## 🟠 Nhóm 2 — Cần cho website

### 2.1 Ảnh nền hero

- **Tên file:** `kame_web_hero_1920x1080_v1.jpg`
- **Kích thước:** `1920 x 1080`

```
[ASSET LOCK]

Wide cinematic key art, 1920x1080, for a website hero background.

Scene: a vast warm-toned sky at golden hour over floating rocky islands, a
distant oriental temple silhouette, and one glowing amber dragon ball orb
suspended in the sky. Painterly anime landscape, warm sepia and amber grade
matching the palette.

Composition: keep the LEFT THIRD and BOTTOM THIRD visually calm and low-detail
— website headline and buttons will be overlaid there. Put the visual interest
in the upper right. Darken naturally toward the bottom-left so light text
stays readable.

No characters with recognizable faces. No text in the image.
```

### 2.2 Ảnh chia sẻ liên kết

- **Tên file:** `kame_web_og_1200x630_v1.jpg`
- **Kích thước:** `1200 x 630` — chuẩn xem trước khi dán link lên Zalo/Facebook

```
[ASSET LOCK]

Social link preview card, 1200x630.

Centered: the Ngọc Rồng Kame logo lockup. Background: aged parchment texture
over a warm dark brown field, with a soft amber glow behind the logo and faint
oriental cloud motifs in the corners. A thin #D80000 red border inset a few
pixels from the edge, like a game UI panel.

Bottom edge: small clean text "nrokame.com" in #F2E6CE.

Simple, high contrast, readable as a small thumbnail.
```

### 2.3 Hoạ tiết giấy cũ (lặp liền mạch)

- **Tên file:** `kame_tex_giay-cu_1024x1024_v1.jpg`
- **Kích thước:** `1024 x 1024`, **lặp liền mạch**

```
[ASSET LOCK]

Seamless tileable texture, 1024x1024.

Aged parchment paper: warm #E8D9BC base with irregular sepia staining, subtle
fiber grain, faint water marks, a few very light ink speckles. Soft and even
— no strong focal point, no dark corners, no vignette.

Must tile seamlessly with no visible seam. Low contrast so text remains
readable when laid on top.
```

### 2.4 Dải nền phụ (tuỳ chọn, làm sau)

- **Tên file:** `kame_web_band-<tên>_1920x600_v1.jpg`
- Ví dụ: `kame_web_band-boss_1920x600_v1.jpg`

```
[ASSET LOCK]

Wide atmospheric band image, 1920x600, for a website section background.

Subject: [ĐỔI PHẦN NÀY — ví dụ "a dark cavern with amber torchlight" cho khối
boss, hoặc "a serene mountain training ground at dawn" cho khối cộng đồng].

Heavily darkened and low contrast — this sits BEHIND text. No focal subject in
the center. Warm brown and amber grade only.
```

---

## 🟡 Nhóm 3 — Mạng xã hội

### 3.1 Ảnh bìa Facebook

- **Tên file:** `kame_social_fbcover_1640x856_v1.jpg`
- **Kích thước:** `1640 x 856`

```
[ASSET LOCK]

Facebook page cover, 1640x856.

Logo lockup positioned in the LEFT HALF. Right half: painterly warm scenery —
floating islands and an amber orb — fading toward the palette's dark brown.

IMPORTANT: keep the bottom-left circular area (about 180x180 pixels) clear of
important content — the profile picture overlaps there. Keep the outer 100
pixels on all sides free of critical detail — mobile crops them.

Below the logo, room for one short tagline in #F2E6CE.
```

### 3.2 Ảnh đại diện

- **Tên file:** `kame_social_avatar_1024x1024_v1.jpg`
- Lấy `kame_logo_vuong` đặt trên nền `#2E1A0C`. Không cần sinh riêng.

---

## 🟢 Nhóm 4 — Trong game

⚠️ Nhóm này **phải làm đủ 3 kích thước** `x2` `x3` `x4` theo tỉ lệ `1 : 1.5 : 2`. Thiếu mức nào thì máy có màn hình tương ứng bị **trống ảnh**.

### 4.1 Nền bảng giao diện

Thay `bgshowicon.png` — bảng nền hiện khi mở chức năng trong game.

| File | Kích thước |
|---|---|
| `kame_game_panel_1000x650_v1.png` | `1000 x 650` (x4) |
| `kame_game_panel_750x487_v1.png` | `750 x 487` (x3) |
| `kame_game_panel_500x325_v1.png` | `500 x 325` (x2) |

> Sinh **một bản 1000×650** rồi thu nhỏ ra 2 bản kia. Đừng sinh 3 lần — sẽ ra 3 hình khác nhau.

```
[ASSET LOCK]

Game UI panel background, 1000x650, landscape.

Aged parchment panel with a bold #D80000 red rounded border frame. Inside:
warm sepia paper with faint oriental cloud and floral motifs, light ink
splatter, and a large coiled oriental dragon rendered in muted olive green
and tan — sitting BEHIND everything at low contrast, like a watermark.

The center must stay calm and readable — game text and icons will be drawn on
top. Darkest at the edges, lighter toward the middle. Warm lantern-light mood.
```

### 4.2 Icon chức năng

Kích thước thật đang dùng: `110 x 110` (phần lớn), `88 x 88`, `75 x 75`.

- **Tên file:** `kame_icon_<tên>_110x110_v1.png` — nền trong suốt
- Tên gợi ý: `boss`, `tuido`, `vongquay`, `phucloi`, `gopy`, `auto`

```
[ASSET LOCK]

Single game UI icon, 110x110, transparent background.

Subject: [ĐỔI — ví dụ "a treasure chest" cho túi đồ, "a spinning prize wheel"
cho vòng quay, "a red demon head" cho boss].

Style: chunky hand-painted 2D game icon, thick #600018 maroon outline, warm
#F06030/#F09060 fills with soft airbrush shading, small parchment-tone
highlights. Bold simple silhouette that reads clearly at 40x40 pixels. Centered
with a few pixels of padding on all sides. No text.
```

---

# 5. THỨ TỰ NÊN LÀM

1. **Logo vuông** → dùng được ngay cho avatar Zalo/Facebook khi lập nhóm
2. **Logo ngang** → cần cho web và client
3. **Ảnh chia sẻ liên kết** → cần ngay khi bắt đầu gửi link mời người test
4. **Ảnh nền hero** + **hoạ tiết giấy cũ** → cần khi code web
5. **Ảnh bìa Facebook** → khi lập fanpage
6. **Nền bảng giao diện + icon** → gom lại, làm cùng lúc với lần build client

---

# 6. GỬI CHO TÔI THẾ NÀO

Cứ để nguyên tên theo quy tắc rồi gửi. Tôi tự biết file nào đi đâu.

Nếu ảnh **không đúng kích thước yêu cầu** thì nói tôi biết — tôi cắt/co lại được, nhưng cần biết trước để không làm méo.

Nếu sinh nhiều phương án, gửi cả `v1` `v2` `v3` và nói bạn thích cái nào. Tôi so được cái nào hợp với game hơn.
