# Map 1 Làng Lá — đợt 2: những ảnh cần vẽ

Viết ngày 22/08/2026, sau khi đã sửa hết những gì code sửa được.
Điểm hiện tại: **45/100 → đang đo lại**. Mọi lỗi còn lại đều nằm ở *ảnh*, không ở code.

> Kích thước dưới đây là **mức x4** (mức gốc). Các mức x1/x2/x3 mình tự thu nhỏ, người vẽ chỉ cần giao x4.
> Ảnh nền để **PNG không cần trong suốt**; tile và vật thể để **PNG có kênh trong suốt**.

---

## Ưu tiên 1 — Bộ tile vách gò (3 ảnh, 96×96 mỗi ảnh)

**Chữa:** vách gò hiện ghép từ 5 chất liệu khác nhau, có một đường ngang thẳng băng
chia đôi vách (độ sáng nhảy +55% trong đúng một hàng pixel, kéo dài suốt bề ngang gò).

Khi mặt đất cao lên một bậc, phần đất lộ ra ở mặt cắt đang phải mượn tạm mấy ô dùng
để lấp lòng đất — nên nhìn như vá chằng vá đụp.

| Tên file | Là gì | Yêu cầu |
|---|---|---|
| `tile_vach_go_trai_x4.png` | mặt cắt đất bên **trái** gò | mép trái là đất lộ thiên (rễ, đá nhỏ), mép phải nối liền vào lòng đất |
| `tile_vach_go_phai_x4.png` | mặt cắt đất bên **phải** gò | đối xứng gương của ảnh trên |
| `tile_chan_go_x4.png` | ô nối giữa vách và mặt đất phẳng | dùng ở chân vách để chuyển mềm |

**Bắt buộc:** tông màu **đáy ô mặt cỏ** phải nối liền vào **đỉnh ô vách** — hiện đang lệch
hẳn nên lộ đường cắt. Cả 3 ảnh phải ghép được với nhau và với ô đất thường mà không thấy mối nối.

```
English prompt: seamless 96x96 pixel-art tile, cross-section of a dirt embankment
wall, exposed soil with small roots and pebbles, warm brown earth tones matching
grass-topsoil tile, hand-painted 2D game tileset style, edges must tile seamlessly
with adjacent dirt tiles, transparent background outside the tile square, no border,
no outline, no drop shadow
```

---

## Ưu tiên 2 — Tile đất có vân liền mạch (4 ảnh, 96×96)

**Chữa:** gần **1/4 màn hình dưới** là đất. Hiện chỉ có 2 kiểu ô nên mắt bắt được lưới
ô vuông 48px. Đã thử làm tối dần bằng code — đỡ nhưng chưa hết, vì bản thân ô vẫn phẳng.

| Tên file | Là gì |
|---|---|
| `tile_dat_lienmach_1_x4.png` | đất tầng mặt |
| `tile_dat_lienmach_2_x4.png` | biến thể, khác vân nhưng **cùng độ sáng** |
| `tile_dat_tang_giua_x4.png` | tối hơn ~8% |
| `tile_dat_tang_day_x4.png` | tối hơn ~16% |

**Bài học đã trả giá:** hai ô cạnh nhau lệch quá **8% độ sáng** là mắt bắt ngay ranh ô,
tầng đất biến thành bàn cờ. Bốn ảnh này phải **rất gần nhau về độ sáng**, khác nhau ở *vân*.

Vân phải chạy **hết ra tận mép** cả 4 cạnh và **khớp với chính nó khi lặp** (tileable).

```
English prompt: seamless tileable 96x96 pixel-art dirt tile, dense fine soil
texture with tiny stones and cracks, texture runs edge to edge and tiles
seamlessly in all four directions, warm brown, hand-painted 2D game tileset,
flat lighting no directional shadow, no border, no outline
```

---

## Ưu tiên 3 — Vẽ lại nền núi Hokage (1 ảnh, 1536×768)

**Chữa:** ảnh hiện tại thật ra là bản **768px nhân đôi** thành 1536 — nên trong game
núi Hokage lặp 2 lần trên một màn, có đường nối dọc rõ ở giữa. Bốn khuôn mặt đá giống
hệt nhau xuất hiện hai lần.

| Tên file | `nen_nuihokage_rong_x4.png` |
|---|---|
| Kích thước | **1536×768**, vẽ liền một mạch, **không ghép đôi** |
| Nội dung | dãy núi đá khắc **5 khuôn mặt Hokage khác nhau**, mỗi mặt một dáng |
| Bố cục dọc | 0–55% chiều cao: trời + mây · 55–80%: khối núi đá · 80–100%: **chuyển dần sang đúng màu trời** |

**Rất quan trọng — phần 20% dưới cùng:** phải nhoà dần vào màu trời `RGB(10, 86, 154)`
cho tới khi hàng cuối cùng là **đúng màu đó**. Nếu để chân núi cắt phẳng, trong game
sẽ hiện một nhát chém ngang chạy suốt màn hình (lỗi đã vấp, xem `LOI-DA-GAP.md` §5).

```
English prompt: 1536x768 hand-painted 2D game background, distant mountain range
carved with five different giant stone Hokage faces in Naruto style, each face
distinct, blue sky with soft cumulus clouds in the upper half, the bottom 20 percent
fades smoothly into flat sky blue RGB(10,86,154), single continuous artwork not
mirrored or tiled, painterly anime background art
```

---

## Ưu tiên 4 — Bóng đổ (8 ảnh nhỏ)

**Chữa:** nhà, giếng, cây đều như sticker dán lên nền cỏ.

Đã thử vẽ bóng bằng code và **thất bại**: chân vật nằm sát mép dưới ảnh nên bóng
không có chỗ lan ra. Cần **ảnh bóng riêng**, mỗi công trình một cái:

| Tên file | Kích thước | Dùng cho |
|---|---|---|
| `bong_nha_nho_x4.png` | 576×96 | nhà nhỏ |
| `bong_nha_lon_x4.png` | 768×96 | nhà lớn |
| `bong_quan_x4.png` | 768×96 | quán mì |
| `bong_torii_x4.png` | 480×96 | cổng torii |
| `bong_cay_x4.png` | 576×96 | cây |
| `bong_gieng_x4.png` | 288×96 | giếng |
| `bong_thung_x4.png` | 192×96 | thùng gỗ |
| `bong_den_x4.png` | 192×96 | đèn lồng |

Mỗi ảnh: một **vệt elip dẹt màu nâu thẫm, mờ dần ra mép**, nền trong suốt hoàn toàn.
Đậm nhất ở giữa khoảng 40%, ra mép về 0.

---

## Ưu tiên 5 — Kè đá và bậc thang (3 ảnh, 96×96)

Concept đã duyệt có **kè đá** dưới chân làng và **cầu thang đá** dẫn lên cổng torii.
Game hiện không có cái nào, nên chỗ cao chỗ thấp nối nhau bằng vách đất trơ.

| Tên file | Là gì |
|---|---|
| `tile_ke_da_x4.png` | tường kè đá xếp, dùng thay vách gò ở khu trung tâm |
| `tile_bac_thang_x4.png` | một bậc thang đá |
| `tile_bac_thang_dinh_x4.png` | bậc trên cùng, nối vào mặt cỏ |

---

## Ưu tiên 6 — Hai ô chuyển tiếp cỏ ↔ lối lát đá (2 ảnh, 96×96)

Hiện cỏ và lối đá gặp nhau bằng **một nhát cắt dọc tuyệt đối phẳng**, màu nhảy 67%
trong đúng một pixel.

| Tên file | Là gì |
|---|---|
| `tile_co_sang_da_x4.png` | nửa trái cỏ, nửa phải đá, ranh giới **lởm chởm tự nhiên**, vài túm cỏ mọc lấn sang đá |
| `tile_da_sang_co_x4.png` | ảnh gương của trên |

---

## Ưu tiên 7 — Vẽ lại nền làng xa (1 ảnh, 1536×576)

Cùng bệnh với núi: bản 768px ghép đôi. Ngoài ra **đáy bị cắt phẳng** nên chỉ lộ ra một
dải mái ngói cao 13px rồi cụt, thành hàng mái nhà treo lơ lửng giữa trời.

| Tên file | `nen_lang_xa_rong_x4.png` |
|---|---|
| Kích thước | 1536×576, vẽ liền mạch |
| Nội dung | mái ngói làng ninja nhìn từ xa, cây xanh xen kẽ, mờ như ở xa |
| Đáy | **20% dưới nhoà dần sang** `RGB(10, 86, 154)` |

---

## Ưu tiên 8 — Thêm cây (2 ảnh)

Cả map hiện chỉ có **một** cây. Cần thêm 2 dáng khác để rải cho đỡ trống:

| Tên file | Kích thước | Là gì |
|---|---|---|
| `tt_cay_nho_x4.png` | 384×384 | cây nhỏ tán tròn |
| `tt_bui_ram_x4.png` | 288×192 | bụi rậm thấp |

---

## Sáu lỗi hay mắc (đọc trước khi vẽ)

1. **Viền quanh ảnh.** Đừng vẽ khung, đừng để viền trắng/đen. Đã từng có 7 tile dính
   viền, trong game lộ thành ô vuông viền trắng giữa mặt đất.
2. **Kích thước phải đúng từng pixel.** Sai một pixel là Unity phóng ảnh lên số chẵn
   gần nhất và game cắt mất một phần.
3. **Chân vật phải chạm đúng mép dưới ảnh.** Hở 2px thôi là trong game nhà lơ lửng.
4. **Ảnh nền không được ghép đôi.** Vẽ liền một mạch đủ bề rộng.
5. **Đáy ảnh nền phải nhoà về đúng màu đã ghi**, không được cắt phẳng.
6. **Tile phải lặp được với chính nó** — dán 4 bản cạnh nhau không thấy mối nối.
