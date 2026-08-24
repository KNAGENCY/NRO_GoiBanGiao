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

---

## 6. Nhà và cây biến mất đột ngột khi cuộn camera

**Triệu chứng.** Đi sang phải, đến một điểm nào đó thì cả cái nhà (hoặc cả cây cổ thụ) **biến mất trong một khung hình**, dù còn gần 200px thân của nó đáng lẽ vẫn nằm trong màn.

**Nguyên nhân.** Vật thể lớn được nhét vào lưới tile như **một ô đơn**, nhưng ảnh của nó trải rộng 6–8 ô sang phải. Vòng vẽ tile chỉ chạy từ `gssx`, mà `gssx = cmx / 24 - 1` — chỉ nới **một ô** sang trái. Khi ô neo trượt ra khỏi mép trái, cả vật thể bị bỏ qua.

**Cách sửa.** Nới mép trái đủ bọc vật rộng nhất, và phải tính **sau** `gssxe` để mép phải giữ nguyên:

```csharp
gssxe = gssx + gssw;      // tính mép phải TRƯỚC
gssx -= OBOC_VAT_LON;     // rồi mới lùi mép trái (9 ô)
if (gssx < 0) gssx = 0;
```

Sửa ở **cả hai** chỗ tính `gssx` trong `GameScr.cs` (dòng ~1613 và ~1809). Kiểm chứng: ở `cmx=497`, trước khi sửa `gssx=19` > cột cây 15 nên cây mất sạch; sau khi sửa `gssx=10` nên cây vẫn vẽ.

---

## 7. Bóng đổ nằm chìm trong lòng đất

**Triệu chứng.** Nhân vật đứng đúng trên cỏ, nhưng cái bóng lại nằm tách hẳn xuống dưới, trong tầng đất.

**Cách đo (dùng lại được cho mọi thứ vẽ sai chỗ).** Chụp hai ảnh với nhân vật ở hai vị trí khác nhau **trong lúc camera đứng yên** (đứng gần mép trái thì `cmx` bị kẹp về một giá trị), rồi trừ pixel hai ảnh. Vùng khác nhau chính là nhân vật + bóng, tách bạch khỏi nền:

- thân nhân vật: ảnh y 542–626 (đáy trùng mép cỏ — đúng)
- bóng: y 646–656, **tách rời hẳn** — tâm lệch +25px ảnh = +12,5 đơn vị

**Cách sửa.** Bù đúng lượng đo được trong `Char.paintShadow`: `g.drawImage(TileMap.bong, xSd, ySd - BU_BONG_Y, 3)` với `BU_BONG_Y = 12`.

**Chưa truy ra gốc**, và ghi rõ điều đó trong mã: `xSd`/`ySd` đều bằng `cx`/`cy` và đều đúng (432), bóng lại đi qua đúng hàm `drawRegion` như thân nhân vật. Nghi do `g.translate(0, GameCanvas.transY)` gọi trước `paintShadow` trong `GameScr`, nhưng `transY = 7` không khớp con số đo được. Ai truy ra gốc thì bỏ lượng bù này đi.


---

## 8. Đường kẻ ngang chạy suốt màn hình ở đáy mỗi lớp nền

**Triệu chứng.** Một (hoặc nhiều) đường kẻ ngang vắt qua cả khung hình, **đứng yên khi cuộn camera** — khác với mối nối *dọc* do ảnh nền ghép đôi.

**Nguyên nhân.** Các lớp nền là PNG **đặc**, vẽ với chiều cao cố định. Dưới đáy mỗi lớp, engine điền nốt màn hình bằng **một màu đơn** (`fillRect` với màu hàng đáy của ảnh). Nếu đáy ảnh còn nội dung — mây, mái ngói, tán cây — thì trên đường đáy là nội dung, dưới là mảng màu phẳng.

**Cách sửa.** Cho phần dưới mỗi lớp nhoà dần về đúng màu **đỉnh lớp kế dưới**, và 20% cuối cùng là màu đó **thuần** (hết vân). Script: `tools/map-lang-la/lam-mem-day-lop-nen.py` — xử cả ba cặp b03→b02, b02→b01, b01→b00 trong một lần.

Đo kiểm bằng **độ gợn sóng hàng đáy**: `b01` 13,7 → 0,0 và `b03` 1,6 → 0,0.

> **Bẫy khi đo:** "độ lệch chuẩn của hàng" phải tính **theo chiều ngang, cho từng kênh màu**.
> Tính `std` trên cả 3 kênh thì một hàng đồng màu `(12,95,169)` cho ra std ≈ 64 — đọc nhầm
> thành "còn nhiều vân" trong khi hàng đó đã phẳng tuyệt đối.

**Không phải mọi đường ngang đều là lỗi.** Đo ở *ba cột khác nhau* trước khi kết luận: nếu chỉ một cột nhảy còn các cột kia chuyển mượt thì đó là ranh giới tự nhiên của vật (đỉnh núi, mép mái nhà), không phải mối nối lớp.

---

## 9. Ô chân dung trên HUD

**Không phải "ô bị hỏng".** HUD gốc Ngọc Rồng **không hề vẽ chân dung** — ô xanh bên trái chỉ có tên và sức mạnh, phần còn lại bỏ trống; cái hình nhìn như "biểu tượng ảnh hỏng" thực ra là **hoạ tiết vẽ sẵn trong chính ảnh khung** `myTexture2dpanel.png`.

Nên việc phải làm là **thêm hẳn chỗ vẽ**, không phải "nối bảng ánh xạ vào HUD":

1. Máy chủ gửi bảng ánh xạ `headId → avatarId` trong ME_INFO (trước gửi `size = 0` nên `Char.getAvatar` trả −1).
2. Thêm `SmallImage.drawSmallImage(g, c.avatarz(), 6, 6, …)` vào `GameScr.paintImageBar`.
3. Dịch tên và sức mạnh sang phải 15 đơn vị, nếu không chữ đè lên mặt.

**Đừng dùng ảnh 27078 làm chân dung người chơi** — đó là `iruka-portrait`, chân dung NPC thầy Iruka đã đăng ký trong `catalog.json`. Dùng chung thì người chơi và NPC có cùng một khuôn mặt. Ảnh riêng: `Small27084`, cắt từ chính ảnh đầu nhân vật (`tools/map-lang-la/tao-chan-dung-nguoi-choi.py`).

---

## 10. Tường vô hình không được chiếm cả cột

Đặt tile tường (ảnh trong suốt) **cả cột** thì phần dưới mặt đất không có tile nào — mép màn hình hở ra một dải màu trời, nhìn như map bị khuyết.

Chỉ đặt tường ở phần **trên mặt đất**; phần dưới dùng tile đất bình thường. Chặn ngang của client kiểm ở hàng *(mặt đất − 1)* nên thế là đủ giữ người chơi lại. Đo kiểm: mép trái dưới mặt đất còn **0%** pixel màu trời.
