# PROMPT SẴN — copy dán thẳng cho AI tạo ảnh

> Dùng kèm `README.md`. Mỗi khối dưới đây copy nguyên, dán vào AI tạo ảnh (Nano Banana / Midjourney / Flux…).
> **Luôn gắn ảnh tham chiếu** trong `anh-tham-chieu/` vào mỗi lần tạo.

---

## BƯỚC 1 — Bức toàn cảnh để duyệt (làm TRƯỚC, chỉ 1 ảnh)

**Gắn kèm ảnh:** `HIENCO-map_nen_nuihokage.png`, `HIENCO-map_la_maingoi_2.png`, `TOT-naruto-logo.png`

```
A hand-painted 2D side-scrolling game background of a Naruto-style hidden leaf ninja village.
Strict SIDE VIEW (flat orthographic, camera level with the ground) — NOT top-down, NOT 3/4 isometric.

Layers from back to front:
- Clear blue sky with soft white clouds
- Hokage Rock in the background: a large cliff face carved with stern Hokage faces, hazy and desaturated (distance)
- A row of Japanese ninja-village houses: red-orange tile roofs, warm wooden walls, paper sliding doors.
  Houses vary in HEIGHT and WIDTH, grouped in clusters of 2-3 with gaps between — never evenly spaced
- A red wooden torii gate, wooden fences, standing red lanterns, a stone well, stacked barrels
- Ground: brown soil with green grass edging on top, gently uneven with small mounds and steps — not a flat line
- Foreground: dark silhouetted grass tufts and hanging leaves at the bottom edge for depth

Style: rich hand-painted game art, volumetric shading, soft color gradients, clean readable silhouettes,
warm feudal-Japan ninja palette. Detailed and polished like a premium 2D platformer.
NOT flat vector, NOT cartoon for kids, NOT 3D render, NOT pixel art.
Wide horizontal composition, 2048x1024.
```

---

## BƯỚC 2 — Các ảnh rời (chỉ làm SAU khi bước 1 được duyệt)

⚠️ **Mọi prompt bước 2 phải gắn kèm bức concept đã duyệt** để giữ đồng bộ tông màu.

### Khối phong cách chung — dán kèm MỌI ảnh bước 2

```
Match the exact art style, color palette and lighting of the attached approved concept art.
Hand-painted 2D game asset, volumetric shading, clean edges, warm feudal-Japan ninja palette.
NOT flat vector, NOT 3D render.
```

### 2.1 Nền xa

```
[khối phong cách chung]
+ Seamless horizontally-tileable sky background: clear blue sky with soft white clouds.
  Left and right edges must match perfectly for looping. Upper half bright.
  Solid background (no transparency). 768x384.
```

```
[khối phong cách chung]
+ Hokage Rock cliff face carved with stern Hokage faces, hazy desaturated distant tone.
  Rock occupies the LOWER half; UPPER half fully TRANSPARENT.
  Left and right edges seamlessly tileable. Transparent PNG. 768x384.
```

### 2.2 Gạch mặt đất (10 viên — mỗi viên một prompt)

```
[khối phong cách chung]
+ A seamless 2D game terrain tile: rich brown soil texture, tiles seamlessly on all 4 sides,
  no border, no frame, fills the entire square. 96x96.
```

Các viên còn lại đổi phần mô tả:
- `tile_dat_2` → *brown soil with small pebbles and root fragments*
- `tile_co_mat` → *top half lush green grass with ragged uneven blades hanging down into the bottom half of brown soil; grass silhouette must be irregular, not a straight line*
- `tile_co_goc_trai` → *grass covering the top surface AND wrapping down the LEFT edge; right and bottom are plain soil*
- `tile_co_goc_phai` → *grass covering the top surface AND wrapping down the RIGHT edge*
- `tile_vach_trai` → *vertical cut-away soil wall facing left, exposed earth strata*
- `tile_vach_phai` → *vertical cut-away soil wall facing right*
- `tile_doc_len` → *45-degree diagonal slope rising to the right, grass on the slope surface, soil below*
- `tile_doc_xuong` → *45-degree diagonal slope descending to the right, grass on the slope surface*
- `tile_datlat` → *packed dirt path with flat embedded stepping stones, village walkway*

### 2.3 Nhà nguyên căn

```
[khối phong cách chung]
+ A complete single small Japanese ninja-village house, SIDE VIEW, drawn as ONE whole object:
  sloped red-orange tile roof, warm wooden plank walls, one paper sliding door, one lattice window.
  Isolated on fully TRANSPARENT background. The base of the house must touch the exact bottom edge
  of the image with no gap. No ground, no shadow patch, no scenery. 192x192.
```

Đổi mô tả cho 3 căn còn lại: *medium house with a small porch and cloth curtain (288x192)* · *large two-storey house with wooden balcony and hanging lanterns (288x288)* · *ramen shop with hanging noren cloth curtains, wooden signboard, bench seats (288x192)*.

### 2.4 Vật trang trí

```
[khối phong cách chung]
+ A single [VẬT THỂ], SIDE VIEW, isolated on fully TRANSPARENT background.
  Base touches the exact bottom edge of the image, no gap, no ground patch, no drop shadow.
  [CỠ].
```

Thay `[VẬT THỂ]` lần lượt: *large red wooden torii gate (192x192)* · *sparse wooden fence section, tileable horizontally (96x96)* · *standing red paper lantern on a post (96x192)* · *stone well with small wooden roof and hanging bucket (192x192)* · *stack of wooden barrels (96x96)* · *wooden signpost with Japanese lettering (96x192)* · *broad-canopy Japanese tree (192x288)* · *clump of wild grass with small flowers (96x96)*.

### 2.5 Tiền cảnh

```
[khối phong cách chung]
+ Foreground silhouette layer: hanging tree branches with leaves drooping from the TOP edge,
  very dark desaturated green (backlit foreground). Fully TRANSPARENT background. 384x192.
```

```
[khối phong cách chung]
+ Foreground silhouette layer: tall dark grass tufts, very dark desaturated green,
  rooted at the BOTTOM edge. Fully TRANSPARENT background. 192x96.
```

---

## Nhắc cuối — 3 lỗi AI hay mắc

1. **Vẽ luôn cả nền đất dưới vật thể** → sai. Vật thể phải đứng riêng trên nền trong suốt.
2. **Chừa khoảng trống dưới chân** → vật thể sẽ lơ lửng trong game. Chân phải chạm mép dưới.
3. **Trả về cỡ lẻ** (1024x1024) → phải đúng cỡ ghi trong bảng (bội số 96). Nếu AI không xuất đúng cỡ, cắt/resize lại trước khi giao.
