# Prompt giao việc cho agent tạo ảnh

Copy khối dưới đây, dán cho agent đã được cấp quyền đọc thư mục `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la`.

---

## GIAI ĐOẠN 1 — làm 3 ảnh mẫu (dán cái này TRƯỚC)

```
Bạn là hoạ sĩ tạo ảnh cho game 2D Naruto. Thư mục làm việc:
D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la

ĐỌC TRƯỚC KHI VẼ (bắt buộc, theo thứ tự):
1. README.md — đề bài đầy đủ, bảng tỉ lệ, danh sách 27 ảnh
2. PROMPT-COPY-DAN.md — prompt tiếng Anh cho từng ảnh
3. concept-map1-langla.png — concept ĐÃ DUYỆT, bám đúng phong cách/tông màu/nét vẽ này
4. anh-tham-chieu\THUOC-TI-LE.png — QUAN TRỌNG NHẤT: vật nào cao bao nhiêu ô
5. anh-tham-chieu\HIEN-TRANG-map1-trong-game.png — map hiện tại, thứ cần thay
6. anh-tham-chieu\GIAI-THICH-tileset-vien-chuyen-tiep.png — vì sao cần viên viền
7. anh-tham-chieu\TOT-cong-torii-da-tach-nen.png — chuẩn vàng cho vật thể rời
8. anh-tham-chieu\XAU-*.png — lô từng bị trả lại, TRÁNH TUYỆT ĐỐI

VIỆC GIAI ĐOẠN 1: chỉ tạo ĐÚNG 3 ảnh sau (làm mẫu để duyệt trước khi làm loạt):
  1. nha_nho_x4.png        576x576   (README §3.3)
  2. tile_co_mat_x4.png     96x96    (README §3.2)
  3. tt_cong_torii_x4.png  480x480   (README §3.4)

MỖI LẦN TẠO ẢNH PHẢI GẮN KÈM 2 ẢNH THAM CHIẾU:
  - concept-map1-langla.png  (giữ đúng phong cách)
  - anh-tham-chieu\THUOC-TI-LE.png  (giữ đúng tỉ lệ)

LUẬT KHÔNG ĐƯỢC VI PHẠM:
- Cỡ ảnh phải ĐÚNG y như ghi trên (đều là bội số của 96). Nếu công cụ xuất ra cỡ
  khác (vd 1024x1024) thì phải cắt/resize lại về đúng cỡ trước khi lưu.
- Nền TRONG SUỐT thật (PNG có kênh alpha) cho nha_nho và tt_cong_torii.
  KHÔNG vẽ nền đất, KHÔNG vẽ bóng đổ dưới chân, KHÔNG vẽ cảnh xung quanh.
- CHÂN vật thể phải chạm ĐÚNG mép dưới ảnh — không hở khoảng trống, không thừa.
- Nhà phải CAO gấp 2,5-3 lần nhân vật (xem THUOC-TI-LE.png). Vẽ nhà thấp là sai.
- Nhìn NGANG (side view) phẳng, KHÔNG phối cảnh 3/4, KHÔNG nhìn từ trên xuống.
- tile_co_mat_x4.png phải ghép LIỀN MẠCH 4 phía khi lặp, mép cỏ răng cưa không đều
  (không được là đường thẳng).

LƯU VÀO: D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-da-tao\
Đặt đúng tên file như trên.

BÁO CÁO khi xong (nói thật, không suy đoán là "đạt"):
- Đã tạo mấy ảnh, cỡ thực tế từng ảnh là bao nhiêu
- Ảnh nào có nền trong suốt, ảnh nào không
- Ảnh nào bạn KHÔNG CHẮC đạt yêu cầu và vì sao
DỪNG LẠI sau 3 ảnh này, chờ duyệt. Không tự làm tiếp 24 ảnh còn lại.
```

---

## GIAI ĐOẠN 2 — làm nốt 24 ảnh (chỉ dán SAU khi 3 ảnh mẫu được duyệt)

```
3 ảnh mẫu đã được duyệt. Giờ tạo nốt 24 ảnh còn lại theo README.md §3,
dùng đúng các prompt trong PROMPT-COPY-DAN.md.

Thứ tự làm:
  A. Tầng phông (README §3.1) — 3 ảnh
  B. Gạch mặt đất (§3.2) — 9 ảnh còn lại (đã có tile_co_mat)
  C. Nhà (§3.3) — 3 căn còn lại (đã có nha_nho)
  D. Vật trang trí (§3.4) — 7 ảnh còn lại (đã có tt_cong_torii)
  E. Tiền cảnh (§3.5) — 2 ảnh

Giữ nguyên mọi luật ở giai đoạn 1: đúng cỡ (bội số 96), nền trong suốt,
chân chạm mép dưới, đúng tỉ lệ theo THUOC-TI-LE.png, gắn kèm concept mỗi lần tạo.

Lưu vào cùng thư mục anh-da-tao\. Báo cáo trung thực như lần trước.
```
