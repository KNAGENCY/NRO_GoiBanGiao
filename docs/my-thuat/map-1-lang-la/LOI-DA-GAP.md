# Map 1 Làng Lá — những lỗi đã vấp và cách sửa

Ghi ngày 22/08/2026, sau đợt soi kỹ map trong game thật (zoom từng vùng, đo pixel).
Mục đích: phiên sau **không vấp lại**. Mỗi lỗi ghi theo dạng *triệu chứng → nguyên nhân → cách sửa*.

> Ảnh chứng: `D:\NRO_Project\goc4.png` (trước), `D:\NRO_Project\lienA.png` (sau).

---

## 1. Nhà bị cắt ngang, có 2 dải mặt đất chồng nhau

**Triệu chứng.** Quán ramen và mấy nhà khác bị một dải cỏ + đất cắt ngang thân, trông như bị chôn nửa người.

**Nguyên nhân.** File `res/mymap/1.bytes` là bản đồ **kế thừa từ Ngọc Rồng**. Các đợt trước chỉ *ghi đè* các ô mình dùng (slot ≥ 60) mà **không xoá tile cũ**. Hàng 17 còn nguyên 69 tile Dragon Ball → thành dải mặt đất thứ hai nằm giữa thân nhà.

**Cách sửa.** Trước khi vẽ bố cục, **xoá sạch mọi tile < 60**:

```python
for r in range(tmh):
    for c in range(tmw):
        if 0 < lay(r, c) < 60:
            dat(r, c, 0)
```

Kiểm chứng: sau khi vẽ, đếm lại `sum(1 for v in cells if 0 < v < 60)` — **phải bằng 0**.

---

## 2. Ô vuông viền trắng nằm giữa mặt đất

**Triệu chứng.** Vài ô 48px có khung trắng rõ rệt, như ô bị thiếu ảnh.

**Nguyên nhân.** Không phải thiếu ảnh — **viền trắng nằm sẵn trong chính file PNG** do AI tạo ảnh vẽ kèm khung. Tile 63 dính trắng cả 4 cạnh; 66 dính cạnh phải; 62/64/65/67/68 dính viền đen.

**Cách sửa.** Không cắt ảnh (sẽ lệch tỉ lệ) mà **kéo hàng/cột sạch gần nhất đè ra mép**. Ngưỡng: sáng > 212 là trắng, tối < 30 là đen — **đừng đặt ngưỡng đen cao hơn 30**, vì tile vách đất vốn tối thật, đặt 42 thì script gọt mất 1/4 ảnh.

Script: `scratchpad/xoa_vien_tile.py`.

---

## 3. Đường kẻ mảnh quanh mọi ô (texture bleeding)

**Triệu chứng.** Sau khi xoá viền vẫn còn khung mờ quanh ô, và các vệt kẻ mảnh ở nền.

**Nguyên nhân.** Cả **242 ảnh** trong `res/x*/t/` và `res/x*/bg/` đang để `wrapMode = Repeat`. Khi Unity lọc bilinear ở mép, nó hút màu từ **cạnh đối diện** của chính ảnh đó.

**Cách sửa.** Đặt `ti.wrapMode = TextureWrapMode.Clamp` cho mọi ảnh, kèm 3 cài đặt đã biết từ trước:

| Cài đặt | Giá trị | Không đặt thì sao |
|---|---|---|
| `wrapMode` | `Clamp` | viền kẻ quanh mọi ô |
| `npotScale` | `None` | ảnh 96px bị phóng lên 128 → game cắt mất 25% |
| `isReadable` | `true` (thư mục `bg`) | `getRGB` trả 0 → trời đen |
| `mipmapEnabled` | `false` | ảnh nhoè khi thu nhỏ |

---

## 4. Game không dùng ảnh nền mình sửa

**Triệu chứng.** Sửa `bg/b02.png` bao nhiêu lần cũng không thấy đổi trong game.

**Nguyên nhân.** Client nạp nền theo đường dẫn `"/bg/b" + typeBg + k + ".png"`, mà **`typeBg = 1`** chứ không phải 0 → nó tìm `b10..b13`. Không thấy trong Resources thì nó **lặng lẽ lấy bản cũ trong bộ nhớ đệm RMS** (nền Dragon Ball 266×176).

Cách nhận biết: hỏi Unity kích thước ảnh đang dùng, so với kích thước file thật.

```csharp
Game1.mGraphics.getRealImageWidth(Game1.GameCanvas.imgBG[i])
```

Lệch nhau ⇒ game **không** dùng file của mình.

**Cách sửa.** Chép `b00..b03` thành `b10..b13` ở cả 4 mức zoom. Mọi lần sửa nền sau này phải **chép sang bản `b1x`**, nếu không công sức đổ sông.

---

## 5. Một đường ngang chạy suốt màn hình, chỗ núi giáp trời

**Triệu chứng.** Chân núi Hokage bị cắt phẳng lỳ, dưới nó là mảng trời màu khác hẳn.

**Nguyên nhân.** Hai chỗ cộng lại:
- Đáy ảnh núi (`b02`) vẫn còn màu đá tại đúng chỗ lớp `b01` đè lên.
- Mép trên `b01` là **(6, 60, 107)** trong khi đáy `b02` là **(10, 86, 154)** → chênh màu thấy rõ.

Chỗ giáp nhau **luôn rơi vào hàng 324/384 của ảnh núi** (≈ 84% chiều cao) bất kể camera lên xuống, vì hai lớp dịch theo cùng một hệ số.

**Cách sửa.**
1. `b02` chuyển dần sang màu trời từ **68% → 84%** chiều cao, dưới 84% là trời thuần. (Fade tới tận đáy **không có tác dụng** — chỗ giáp nằm ở 84%, không phải ở đáy.)
2. Mép trên `b01` (30% trên) chuyển dần từ đúng màu đáy `b02` về màu vốn có.

Đo kiểm: bước nhảy độ sáng lớn nhất ở vùng giáp — **36,9 → 8,3** (mắt thường không còn thấy).

---

## Còn lại, chưa sửa được bằng code

| Việc | Vì sao code không làm được |
|---|---|
| **Nền lặp lộ liễu** — `b01`, `b02`, `b03` đều là bản 384px ghép đôi thành 768px, nên núi Hokage lặp 2 lần trong một màn, có đường nối dọc ở giữa | Cần **vẽ lại ảnh nền rộng, không lặp** — việc của người tạo ảnh |
| **Tile dốc chéo (63, 64) không khớp bậc thẳng** — hiện đang né bằng cách dùng mặt cỏ 62 cho bậc vuông | Cần **bộ ảnh dốc khớp nhau** (dốc lên/xuống + chân dốc) |
| **Map trống** — không NPC, không quái, không vật phẩm | Thuộc phần server, chưa tới lượt |
