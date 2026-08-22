# PROMPT SẴN — copy dán thẳng cho AI tạo ảnh

> **BẢN 2 (22/08)** — đã sửa kích thước theo tỉ lệ thật đo từ sprite trong game.
> Dùng kèm `README.md`.

## ⚠️ MỖI LẦN TẠO ẢNH PHẢI GẮN KÈM 2 ẢNH NÀY

1. `concept-map1-langla.png` — giữ đúng phong cách/tông màu đã duyệt
2. `anh-tham-chieu/THUOC-TI-LE.png` — giữ đúng tỉ lệ (nhà cao gấp 2,5–3 lần người)

---

## KHỐI PHONG CÁCH — dán kèm MỌI prompt bên dưới

```
Match the exact art style, color palette, brushwork and lighting of the attached approved concept art.
Hand-painted 2D side-scrolling game asset, volumetric shading, soft gradients, clean readable silhouette,
warm feudal-Japan ninja palette. Strict SIDE VIEW, flat orthographic — NOT top-down, NOT 3/4 isometric.
NOT flat vector, NOT 3D render, NOT pixel art.
```

---

## A. TẦNG PHÔNG (3 ảnh — xa, mờ, lặp ngang)

```
[KHỐI PHONG CÁCH]
+ Seamless horizontally-tileable sky: clear blue sky with soft white clouds.
  Left and right edges must match perfectly for looping. Upper half bright.
  SOLID background, no transparency. 768x384.
```

```
[KHỐI PHONG CÁCH]
+ Hokage Rock cliff carved with stern Hokage faces, HAZY and DESATURATED (far distance, atmospheric perspective).
  Rock occupies the LOWER half; UPPER half fully TRANSPARENT.
  Left and right edges seamlessly tileable. Transparent PNG. 1536x768.
```

```
[KHỐI PHONG CÁCH]
+ Distant village skyline: a row of Japanese tiled rooftops seen far away, HAZY, PALE, LOW CONTRAST,
  few details, as if veiled in morning mist. Rooftops occupy the LOWER portion; upper portion fully TRANSPARENT.
  Left and right edges seamlessly tileable. Transparent PNG. 1536x576.
```

---

## B. GẠCH MẶT ĐẤT (10 viên — mỗi viên 96x96)

Prompt gốc:

```
[KHỐI PHONG CÁCH]
+ A seamless 2D game terrain tile, 96x96, tiles seamlessly on all four sides,
  no border, no frame, texture fills the entire square. [MÔ TẢ]
```

Thay `[MÔ TẢ]`:

| File | MÔ TẢ |
|---|---|
| `tile_dat_1` | *rich brown packed soil* |
| `tile_dat_2` | *brown soil with small pebbles and root fragments* |
| `tile_co_mat` | *top half lush green grass with RAGGED UNEVEN blades hanging down into the bottom half of brown soil; the grass silhouette must be irregular, never a straight line* |
| `tile_co_goc_trai` | *grass covering the top surface AND wrapping down the LEFT edge; right and bottom plain soil* |
| `tile_co_goc_phai` | *grass covering the top surface AND wrapping down the RIGHT edge* |
| `tile_vach_trai` | *vertical cut-away soil wall facing LEFT, exposed earth strata* |
| `tile_vach_phai` | *vertical cut-away soil wall facing RIGHT* |
| `tile_doc_len` | *45-degree slope rising to the right, grass on the slope surface, soil below* |
| `tile_doc_xuong` | *45-degree slope descending to the right, grass on the slope surface* |
| `tile_datlat` | *packed dirt path with flat embedded stepping stones, village walkway* |

---

## C. NHÀ TẦNG GẦN (4 căn — cao 6–8 ô, RÕ NÉT, đậm màu)

```
[KHỐI PHONG CÁCH]
+ ONE complete Japanese ninja-village house drawn as a SINGLE whole object, SIDE VIEW, sharp and vivid
  (this is a foreground playable-layer building, NOT a hazy background one).
  [MÔ TẢ NHÀ]
  Isolated on a fully TRANSPARENT background. The base of the building must touch the EXACT bottom edge
  of the image with no gap. No ground, no grass patch, no cast shadow, no scenery around it.
  The building must be TALL: roughly 2.5 to 3 times the height of a human character (see attached scale ruler).
  [CỠ]
```

| File | MÔ TẢ NHÀ | CỠ |
|---|---|---|
| `nha_nho` | *small single-storey house: sloped red-orange tile roof, wooden plank walls, one paper sliding door, one lattice window* | 576x576 |
| `nha_vua` | *medium house: wider, front porch with cloth curtain, two lattice windows* | 768x576 |
| `nha_lon` | *large two-storey house: wooden balcony, big tiled roof, hanging lanterns by the door* | 768x768 |
| `nha_quan` | *ramen shop: blue noren cloth curtains hanging in the doorway, wooden signboard with Japanese lettering, bench seats* | 768x576 |

---

## D. VẬT TRANG TRÍ (8 ảnh)

```
[KHỐI PHONG CÁCH]
+ A single [VẬT THỂ], SIDE VIEW, isolated on a fully TRANSPARENT background.
  Base touches the EXACT bottom edge of the image — no gap, no ground patch, no drop shadow.
  [CỠ]
```

| File | VẬT THỂ | CỠ |
|---|---|---|
| `tt_cong_torii` | *large red wooden torii gate with shimenawa rope and small bells* | 480x480 |
| `tt_hangrao` | *sparse wooden fence section, seamlessly tileable horizontally* | 192x192 |
| `tt_denlong` | *standing red paper lantern on a wooden post with base* | 192x384 |
| `tt_gieng` | *stone well with small wooden roof and hanging bucket* | 288x288 |
| `tt_thunggo` | *stack of two or three wooden barrels* | 192x192 |
| `tt_bienhieu` | *wooden signpost planted in the ground, Japanese lettering, rope binding* | 192x288 |
| `tt_cay` | *broad-canopy Japanese tree, brown trunk, deep green foliage* | 576x768 |
| `tt_bui_co` | *clump of wild grass with small flowers* | 192x96 |

---

## E. TIỀN CẢNH (2 ảnh — sẫm gần đen)

```
[KHỐI PHONG CÁCH]
+ Foreground silhouette layer, VERY DARK desaturated green (backlit, close to camera):
  hanging tree branches with leaves drooping from the TOP edge of the frame.
  Fully TRANSPARENT background. 768x384.
```

```
[KHỐI PHONG CÁCH]
+ Foreground silhouette layer, VERY DARK desaturated green:
  tall grass tufts rooted at the BOTTOM edge of the frame.
  Fully TRANSPARENT background. 384x192.
```

---

## 6 LỖI AI HAY MẮC — kiểm lại trước khi giao

**5. Nền "trong suốt" mà thật ra chỉ MỜ** (đã xảy ra thật với `nen_nuihokage`: alpha 226 thay vì 0).
   Vùng phải trong suốt thì alpha phải **đúng bằng 0**, không phải "trắng nhạt" hay "mờ 90%".
   Tự kiểm: mở ảnh, vùng đó phải thấy ô caro, KHÔNG phải màu trắng/xám.

**6. Ảnh "lặp ngang" mà hai mép không khớp** (đã xảy ra thật với `nen_troi`: lệch 157/765).
   Mép trái và mép phải phải nối liền như một dải vô tận. Tự kiểm: ghép 3 bản cạnh nhau,
   KHÔNG được thấy đường sọc dọc ở chỗ nối.


1. **Vẽ luôn nền đất/bóng dưới vật thể** → sai. Vật thể phải đứng riêng trên nền trong suốt.
2. **Chừa khoảng trống dưới chân** → vật thể lơ lửng trong game. Chân phải chạm mép dưới.
3. **Vẽ nhà quá thấp** → nhân vật cao hơn nhà. Nhà phải cao gấp 2,5–3 lần người (xem thước).
4. **Trả về cỡ lẻ** (1024x1024) → phải đúng cỡ trong bảng (bội số 96). AI không xuất đúng thì cắt/resize lại trước khi giao.
