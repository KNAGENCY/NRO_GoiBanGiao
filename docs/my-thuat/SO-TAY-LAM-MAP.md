# Sổ tay làm map — đọc trước khi bắt đầu map mới

Viết ngày 23/08/2026, rút từ trọn một ngày làm map 1 Làng Lá: **4 đợt soi bằng agent, ~120 lỗi, 15 commit**. Map cuối cùng đạt, nhưng đi vòng rất nhiều.

Sổ tay này để **map 2–5 không phải đi lại con đường đó**.

---

## 0. Vì sao map 1 "vỡ" — một câu

**Làm ngược thứ tự: dựng trước, đối chiếu sau.**

Mỗi vòng chỉ sửa cái lỗi trước mắt, không lần nào đặt màn hình cạnh concept đã duyệt. Kết quả là nhiều bản "sửa lỗi" lại **đẩy xa concept thêm**:

| Sửa lỗi kỹ thuật | Nhưng concept muốn |
|---|---|
| Bỏ bớt nhà, giãn ra cho "khỏi chạm nhau" | Nhà **san sát**, chồng lấn nhau là bình thường |
| Fade chân núi cho hết đường kẻ ngang | Núi **không fade** — nó bị **giấu sau mái nhà** làng xa |
| Đặt gò cao thành hai mảnh rời | Làng đứng trên **một nền đá liền** |

Chủ dự án gọi đó là "chơi xếp hình". Đúng.

**Cách tránh:** bước 1 và bước 6 của quy trình dưới đây — đo đối chiếu concept **ngay từ đầu** và **sau mỗi thay đổi**.

---

## 1. QUY TẮC VÀNG về cỡ ảnh

> ### Số pixel ở mức x4 = số đơn vị game × 4

Sai quy tắc này là hỏng ảnh, và **đã hỏng thật hai lần**.

Vì sao: ảnh giao ở mức **x4**, nhưng game hiển thị ở mức **x2** — bề rộng bị chia đôi. Ai quen nghĩ "ảnh 1536 là to lắm rồi" sẽ giao thiếu **một nửa**.

### Bảng cỡ tính sẵn (mức x4)

| Vật | Đơn vị game | **px ở x4** | Ghi chú |
|---|---|---|---|
| 1 ô lưới | 24 | **96** | mọi tile đều 96×96 |
| Nhân vật | ~58 | 232 | cao 2,4 ô |
| Nhà nhỏ | 144 | **576** | 6 ô |
| Nhà lớn / quán | 192 | **768** | 8 ô |
| Cổng torii | 120 | **480** | 5 ô |
| Cây | 192 | **768** | 8 ô |
| **Ảnh nền phủ hết 1 màn** | **727** | **2 976** | ⚠️ chỗ hay sai nhất |
| Màn hình | 727 × 433 | 2 908 × 1 732 | |

**Ảnh nền phải rộng ít nhất 2 908px.** Hẹp hơn là nó **lặp ngang** trong cùng một khung hình — và người xem đọc ra ngay là hình dán. Lấy tròn theo bội số 96: **2 976**.

Mọi cỡ đều phải là **bội số của 96**. Cỡ lẻ thì Unity tự phóng lên số chẵn gần nhất và game **cắt mất một phần** — đã mất hẳn cột phải của cổng torii vì lỗi này.

---

## 2. Quy trình 7 bước

### Bước 1 — Đo concept TRƯỚC khi vẽ gì

```bash
python tools/map-lang-la/do-map.py <ảnh-concept.png>
```

Ghi lại 7 con số. **Đó là đích.** Không có đích thì mọi bản sửa đều là đoán.

### Bước 2 — Viết đề bài ảnh, có 5 phần bắt buộc

1. **Phong cách** — chốt một thứ. *(Đã vấp: brief cùng một câu vừa ghi `pixel-art` vừa ghi `hand-painted`; agent vẽ phải hỏi lại.)*
2. **Bảng cỡ** — lấy từ §1, ghi px cụ thể từng file, không ghi "to vừa".
3. **Bố cục dọc từng ảnh nền** — bao nhiêu % là trời, bao nhiêu % là nội dung, đáy nhoà về màu nào. *(Đã vấp: ảnh làng xa có 36% trên là trời thuần → nó nằm thấp hơn chân núi, không che được gì.)*
4. **Màu nối** — mọi ảnh nền phải nhoà về đúng một mã màu ghi rõ.
5. **Sáu lỗi hay mắc** — chép từ `map-1-lang-la/DOT-2-ANH-CAN-VE.md`.

### Bước 3 — Kiểm ảnh TRƯỚC khi lắp

```bash
python tools/map-lang-la/kiem-anh-dot-2.py
```

Kiểm cỡ từng pixel, viền bẩn, đáy nền, ảnh có bị ghép đôi không, chân vật thể có chạm mép không. **Đừng duyệt vì đẹp** — cũng đừng duyệt chỉ vì đo xong: nhìn thật một lượt nữa.

### Bước 4 — Lắp

```bash
python tools/map-lang-la/lap-anh-dot-2.py     # tự thu nhỏ ra 4 mức zoom
```

Ảnh nền còn phải chép sang bản `b1x` — client nạp theo `typeBg = 1`, không phải `b0x`. Quên là sửa mãi không thấy đổi.

### Bước 5 — Xếp bố cục **dày**

Đếm ô trước khi xếp: map 52 ô, trừ 4 ô tường biên còn **48 ô dùng được**. Công trình phải **phủ kín 48 ô**, cho chồng lấn 1–2 ô. Dùng lại một mẫu nhà nhiều lần là bình thường — concept cũng vậy.

### Bước 6 — Đo lại, so với bước 1

```bash
python tools/map-lang-la/do-map.py <ảnh-game.png>
```

**Một giây.** Chưa đạt thì sửa rồi đo lại. Chỉ khi 6/7 chỉ số đạt mới sang bước 7.

### Bước 7 — Bấy giờ mới gọi agent soi

Agent chỉ còn hai việc mà máy không làm được: **nhận xét thẩm mỹ** ("nét vẽ chỏi với nền") và **đối chiếu ý đồ** ("concept muốn gì mà game chưa có").

Gọi agent khi bước 6 chưa đạt là **đốt 45–80 phút** để nghe lại điều máy vừa nói trong một giây.

---

## 3. Bảy con số và ngưỡng

| Chỉ số | Đạt khi | Nói lên điều gì |
|---|---|---|
| Trời chiếm màn | ≤ 30% | trống quá thì bức tranh rỗng |
| Dải nhà cửa là trời | ≤ 3% | làng có liền khối không |
| Ô trơn ở dải giữa | ≤ 10% | chân trời có bị **đóng** không |
| Đất phẳng lì | ≤ 40% | lòng đất có nội dung không |
| Khối núi cao | ≥ 30% khung | núi có còn là khối chủ đạo không |
| Đường kẻ ngang | < 10 | có vạch cắt ngang màn không |
| Nền lặp lại | < 40 | ảnh nền có đủ rộng không |

---

## 4. Ba bẫy khi ĐO — đã mắc cả ba

1. **Mặt nạ "màu trời" bắt cả làng xa mờ** (cũng xanh nhạt) → số không nhúc nhích dù ảnh đổi rõ. Đo bằng **độ chi tiết**, đừng đo bằng màu.
2. **Độ lệch chuẩn của một hàng** phải tính theo chiều ngang **cho từng kênh màu**. Tính trên cả 3 kênh thì hàng đồng màu `(12,95,169)` ra std ≈ 64 — đọc nhầm thành "còn nhiều vân" trong khi hàng đó phẳng tuyệt đối.
3. **Không phải đường ngang nào cũng là lỗi.** Chỉ tính là đường kẻ khi **đa số cột cùng nhảy**. Vài cột nhảy là ranh giới tự nhiên (đỉnh núi, mép mái). Đã đọc nhầm đỉnh núi thành đường kẻ rồi đi "sửa" nhầm, làm hỏng chân núi.

---

## 5. Bốn lớp nền — hiểu đúng thì không mất buổi nào

Bốn lớp nền **không phủ hết chiều cao màn**. Dưới đáy mỗi lớp, engine điền nốt bằng **một màu đơn**.

Hệ quả:
- đáy ảnh còn nội dung (mây, mái ngói) → thành **đường kẻ ngang** chạy suốt màn
- giữa lớp thấp nhất và mặt cỏ có khoảng hở → **trời rơi thẳng xuống cỏ**, làng trông lơ lửng

Cách làm đúng:
- đáy mỗi lớp **nhoà về đúng màu đỉnh lớp kế dưới** (mỗi cặp một màu riêng, không dùng chung một màu cho tất cả)
- lớp thấp nhất phải **có nội dung**, không để là trời trơn
- **đừng fade mép trên** của lớp có nhiệm vụ che lớp trên nó — fade là xoá mất cái nó phải che

---

## 6. Việc gì máy làm, việc gì phải vẽ

**Máy làm được:** xếp bố cục, sinh va chạm từ bố cục, nhoà mép nối các lớp, xoá viền bẩn tile, tạo ô chuyển tiếp bằng cách trộn hai tile, cắt trời thừa, làm tối dần theo độ sâu, tạo chân dung từ ảnh đầu nhân vật.

**Phải vẽ:** mọi thứ cần **nội dung mới** — vân đất, kè đá, bậc thang, nền rộng không lặp, bóng đổ.

**Đã thử bằng máy và thất bại, đừng thử lại:**

| Đã thử | Vì sao hỏng |
|---|---|
| Vẽ **bóng đổ** bằng code | Chân vật nằm sát mép dưới ảnh nên bóng không có chỗ lan; trong game gần như không thấy gì |
| **Tăng tương phản** để lòng đất bớt phẳng | Chênh 48% thì tầng đất thành **bàn cờ** — tệ hơn mảng phẳng ban đầu. Hai ô cạnh nhau lệch quá **8%** là mắt bắt ngay ranh ô |
| **Ghép mạch đá** vào tile đất từ ảnh đá có sẵn | Sinh từ một mẫu nên **mạch lặp y hệt trong mọi ô** → lưới ô vuông rõ hơn cả lúc chưa làm |

> **Quy luật rút ra từ cả ba lần:** máy **sắp xếp** được, **hoà màu** được, nhưng
> **không tạo được nội dung mới**. Thứ gì cần *mỗi ô một khác* — vân đất, mạch đá,
> rễ cây — thì bắt buộc phải vẽ, vì sinh từ một mẫu là lặp, mà lặp thì mắt bắt ngay.
>
> Ba lần thử này tốn khoảng một tiếng. Nhận ra quy luật sớm thì tiết kiệm được.

---

## 7. Đọc thêm

| Cần gì | Đọc file nào |
|---|---|
| 10 lỗi đã vấp và cách sửa từng cái | `map-1-lang-la/LOI-DA-GAP.md` |
| Mẫu đề bài ảnh đầy đủ | `map-1-lang-la/DOT-2-ANH-CAN-VE.md` |
| Kết quả soi chi tiết (có toạ độ pixel) | `map-1-lang-la/SOI-LAN-1.md`, `SOI-LAN-3.md` |
| Công cụ | `tools/map-lang-la/README.md` |
