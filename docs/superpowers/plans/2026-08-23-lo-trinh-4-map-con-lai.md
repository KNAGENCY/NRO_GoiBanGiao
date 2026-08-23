# Lộ trình 4 map còn lại — làm nhanh gấp mấy lần map 1

Viết 23/08/2026, ngay sau khi map 1 Làng Lá đạt 6/7 chỉ số.

## Map 1 tốn bao nhiêu, và tốn vào đâu

| Việc | Thời gian | Đáng lẽ |
|---|---|---|
| 4 đợt soi bằng agent | **~4 giờ** | 1 đợt ở mốc cuối |
| 3 đợt vẽ ảnh (đợt 1 thiếu, đợt 2 sai cỡ) | 3 vòng chờ | 1 đợt |
| Sửa → chụp → nhìn → sửa lại | rất nhiều vòng | đo bằng thước, 1 giây/lần |
| 3 lần thử làm ảnh bằng code, cả 3 hỏng | ~1 giờ | biết trước là không được |

**Gốc:** không có đích đo được, nên mỗi bước đều là đoán; và không có khuôn đề bài, nên ảnh phải vẽ lại.

Cả hai thứ đó **giờ đã có**.

---

## Cách làm map 2–5

### Chuẩn bị một lần (đã xong)

- [`SO-TAY-LAM-MAP.md`](../../my-thuat/SO-TAY-LAM-MAP.md) — quy trình 7 bước, quy tắc cỡ ảnh, bẫy đã vấp
- [`KHUON-DE-BAI-ANH-MAP.md`](../../my-thuat/KHUON-DE-BAI-ANH-MAP.md) — chép ra, thay tên làng là giao được
- `tools/map-lang-la/do-map.py` — thước đo 7 chỉ số, 1 giây
- `tools/map-lang-la/kiem-anh-dot-2.py` — kiểm ảnh trước khi lắp
- `tools/map-lang-la/lap-anh-dot-2.py` · `ve-bo-cuc-map.py` · `sinh-va-cham.py`

### Mỗi map — 6 bước

| # | Việc | Ai làm | Ước |
|---|---|---|---|
| 1 | Vẽ **concept một bức** cho làng đó, duyệt bằng mắt | người vẽ | 1 vòng |
| 2 | `do-map.py <concept>` → ghi 7 con số làm đích | máy | 1 phút |
| 3 | Chép khuôn đề bài, thay tên làng + bảng màu + màu nối → giao **một đợt duy nhất** (~30 ảnh) | mình + người vẽ | 1 vòng |
| 4 | `kiem-anh-dot-2.py` → đạt hết mới lắp; `lap-anh-dot-2.py` | máy | 10 phút |
| 5 | Sửa số trong `ve-bo-cuc-map.py` (cột từng công trình, cột nền cao), chạy nó + `sinh-va-cham.py` | mình | 30 phút |
| 6 | `do-map.py <ảnh game>` → chưa đạt thì quay lại 5. Đạt 6/7 thì **một** đợt soi agent để chốt | máy + agent | 1 giờ |

**Không gọi agent trước bước 6.** Gọi sớm là đốt 45–80 phút để nghe lại điều thước đo vừa nói trong một giây.

### Chỉ số phải đạt trước khi chốt

| Chỉ số | Ngưỡng |
|---|---|
| Trời chiếm màn | ≤ 30% |
| Dải nhà cửa là trời | ≤ 3% |
| Ô trơn ở dải giữa | ≤ 10% |
| Đất phẳng lì | ≤ 40% |
| Khối núi cao | ≥ 30% khung |
| Đường kẻ ngang | < 10 |
| Nền lặp lại | < 40 |

---

## Bốn map, thứ tự đề nghị

Làm **Làng Cát trước** — nó khác Làng Lá nhất (sa mạc, không cây, không núi đá xanh), nên nếu khuôn đề bài chịu được nó thì chịu được cả ba làng còn lại.

| # | Làng | Nét riêng cần giữ | Chỗ dễ vấp riêng |
|---|---|---|---|
| 2 | **Cát** | tường thành đất nung, mái vòm, bão cát | không có cây → dải giữa dễ trống, phải bù bằng tường thành nhiều lớp |
| 3 | **Sương Mù** | cầu, nước, sương dày | sương làm mọi lớp cùng tông → chỗ giáp lớp khó thấy, phải đo chứ đừng tin mắt |
| 4 | **Đá** | vách núi dựng, nhà đục trong đá | địa hình nhiều tầng → bố cục nền cao phức tạp hơn, tính kỹ cột |
| 5 | **Mây** | mây dày, nhà trên cao | không có "mặt đất" theo nghĩa thường → xem lại quy ước va chạm trước khi vẽ |

---

## Việc còn lại của map 1

| Việc | Trạng thái |
|---|---|
| 6/7 chỉ số đạt | ✅ |
| **Đất phẳng lì 91% (đích ≤ 40%)** | ⬜ chờ 4 tile đất có vân — [đề bài đợt 3](../../my-thuat/map-1-lang-la/DOT-2-ANH-CAN-VE.md) |
| 2 ảnh nền đúng cỡ 2976 (hiện 1536 nên núi còn lặp) | ⬜ cùng đợt 3 |
| 8 ảnh bóng đổ đã vẽ nhưng chưa lắp được | ⬜ cần lớp `bgItem` riêng, server phải gửi |
| Waypoint rời map, NPC, quái | ⬜ phần server, không thuộc map |

Ba việc đầu **gộp vào một đợt ảnh cuối** cho map 1, làm cùng lúc với đợt ảnh của Làng Cát.
