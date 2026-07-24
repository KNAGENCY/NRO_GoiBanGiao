# Google Flow (Nano Banana Pro) — Master prompt + tạo ảnh hàng loạt cho game Naruto

> Mục tiêu: dựng **một Custom Tool** trong Google Flow để tạo toàn bộ ảnh game
> Naruto — **đồng bộ phong cách**, đúng kích thước, nền trong suốt, tạo hàng loạt.
> File này gồm: (1) master prompt tạo tool, (2) khối phong cách khoá, (3) cách
> gắn 3 ảnh tham chiếu, (4) quy trình tạo hàng loạt + mọi tính năng.

## Vì sao Flow hợp với dự án này

- **Nano Banana Pro tạo ảnh chất lượng cao thật** — khác hẳn AI cũ (vẽ hình khối
  phẳng 3 màu). Riêng việc đổi engine đã gần như hết lỗi 343 ảnh trước.
- **Nhận tối đa 3 ảnh tham chiếu / lần** (nhân vật + môi trường + phong cách) rồi
  hoà làm một. → gắn **cùng một ảnh phong cách** cho mọi lần tạo là ra **đồng bộ**.
- **Custom Tools**: gói prompt gốc + ảnh tham chiếu cố định thành một công cụ,
  bấm phát ra ảnh đúng chuẩn — không phải dán lại luật mỗi lần.

> ⚠️ Hai điểm Nano Banana KHÔNG tự lo, phải xử ở bước sau (mục 4): **kích thước
> nhỏ chính xác** (64×64…) và **nền trong suốt alpha**. Nano Banana tạo ảnh to,
> nền đặc — phải thu nhỏ + tách nền sau.

---

# PHẦN 1 — MASTER PROMPT TẠO CUSTOM TOOL

Dán khối này khi tạo Custom Tool (ô mô tả tool / prompt gốc của tool). Nó định
nghĩa tool làm gì + nhét sẵn luật game.

```
Tao mot cong cu tao SPRITE GAME 2D phong cach PIXEL-ART cho game nhap vai Naruto.

NHIEM VU: nhan mot dong mo ta ngan (ten_file + noi_dung + nhom + lang) va sinh ra
MOT sprite dung chuan game duoi day. Luon giu dung phong cach da khoa, bat ke noi
dung la gi.

PHONG CACH KHOA (ap cho MOI anh):
- Pixel-art / game sprite 2D chi tiet, CO DO BONG va chuyen mau (shading), vien
  trong dam sach. TUYET DOI khong ve hinh khoi phang vector, khong to mau det.
- Nen: mot mau phang don gian (vi du xam trung tinh) de tach nen de o buoc sau —
  KHONG ve boi canh, KHONG do bong duoi chan.
- Bang mau am, dam chat ninja co trang (xem KHOI PHONG CACH ben duoi).
- Nhan vat/boss/npc: nhin nghieng cho game side-view, mat xoay sang PHAI, toan than.
- Do chi tiet: bang muc "anh tot" trong 3 anh tham chieu dinh kem (khong duoc so
  sai/phang nhu "anh xau").

DAU RA: mot anh vuong, chu de o giua, chua nhieu khoang trong quanh mep de crop.

Khi toi gui mot dong spec, hay tao dung MOT anh theo spec do + phong cach khoa.
Neu spec thieu mau/chi tiet, suy theo LANG (mau he lang) chu dung bia noi dung la.
```

**Sau khi dán master prompt, gắn 3 ảnh tham chiếu (mục 3) làm "ingredients" cố
định của tool**, rồi lưu tool. Từ đó mỗi lần chỉ cần gõ một dòng spec.

---

# PHẦN 2 — KHỐI PHONG CÁCH KHOÁ (chép vào tool, hoặc dán kèm mỗi prompt nếu tool không giữ được)

```
BANG MAU CHUAN (dung xuyen suot, khong lech tong):
- Nen toi: nau go / xam da / den lam am.
- Diem nhan am: cam chay (#E8720C), do (#C0392B).
- Kim loai: xam thep sang toi.
- Giay/cuon: kem nga vang (#E8D8B0).

MAU RIENG TUNG LANG (dung dung khi anh thuoc lang do):
- Lang La (He Hoa):     cam / do.
- Lang Cat (He Phong):  vang cat / nau.
- Lang Suong Mu (Thuy): xanh duong / xam.
- Lang Da (He Tho):     nau dat / xam da.
- Lang May (He Loi):    tim / vang chop.

CHAT LUONG BAT BUOC:
- Anh nhan vat/boss/quai/vat pham/trang phuc: >= 40 mau khac nhau, co do bong.
- Giao dien (nut/khung/thanh): duoc it mau hon nhung PHAI co khoi/vat/vien, khong phang.
- Mep sach, khong vien sang (glow halo) quanh mep ngoai.
```

---

# PHẦN 3 — BA ẢNH THAM CHIẾU (INGREDIENTS) ĐỂ ĐỒNG BỘ

Nano Banana nhận 3 ảnh tham chiếu. Dùng như sau (lấy từ thư mục `anh-tham-chieu/`):

| Khe | Vai trò | Dùng ảnh nào |
|---|---|---|
| **Ảnh 1 — Phong cách** | Khoá tông màu + độ chi tiết + kiểu vẽ | `TOT-naruto-logo.png` (đúng màu hệ làng, đổ bóng, viền đậm) — **gắn CỐ ĐỊNH cho mọi lần tạo** |
| **Ảnh 2 — Mức chi tiết** | Nhắc mức đổ bóng/khối cần đạt | `TOT-chitiet-nhanvat.png` (nhân vật chi tiết) — gắn khi tạo nhân vật/boss/npc |
| **Ảnh 3 — Nhân vật gốc** | Đồng bộ MỘT nhân vật qua nhiều ảnh | Ảnh nhân vật đã tạo tốt → dùng làm tham chiếu cho các ảnh liên quan (trang phục, khung hình của chính nhân vật đó) |

**Mẹo đồng bộ mạnh nhất (đúng cái bạn cần):**
- Tạo xong **một ảnh chuẩn cho mỗi làng/nhân vật** trước → dùng chính nó làm "ảnh
  nhân vật gốc" (khe 3) cho mọi ảnh cùng làng/nhân vật đó → cả bộ trông một style.
- Với 5 làng: tạo 5 "ảnh gốc" (mỗi làng một nhân vật đại diện) → mọi trang phục/
  vũ khí/npc của làng đó tham chiếu ảnh gốc tương ứng.

> Kèm cả `XAU-boss-zabuza.png` vào prompt như **ảnh phản diện** với câu "KHONG ve
> phang nhu anh nay" — Nano Banana hiểu tham chiếu-để-tránh khá tốt.

---

# PHẦN 4 — TẠO HÀNG LOẠT + TẤT CẢ TÍNH NĂNG

## 4.1 Nguồn danh sách 343 ảnh
Toàn bộ tên + nội dung từng ảnh đã có ở **`DOT-1B-LAM-LAI-343-ANH.md`** (9 nhóm).
Mỗi dòng spec cho tool theo mẫu:
```
<ten_file> | <noi dung> | nhom=<nhom> | lang=<lang neu co>
```
Ví dụ: `boss_zabuza_x4.png | kiem si nua duoi mat quan bang, vac dai dao lo khoet, chan co suong mu | nhom=boss`

## 4.2 Thứ tự làm theo lô (để đồng bộ)
Làm **theo nhóm và theo làng**, không nhảy lung tung:
1. Trước tiên tạo **5 "ảnh gốc" đại diện 5 làng** → duyệt tay → dùng làm tham chiếu.
2. Rồi từng nhóm: giao diện → vật phẩm → trang phục → vũ khí → npc → quái → boss →
   hiệu ứng → bản đồ. Trong mỗi nhóm, làm hết một làng rồi sang làng khác.
3. Nano Banana làm việc theo kiểu hội thoại: tạo → xem → sửa nhỏ ("dam mau hon",
   "them chi tiet giap"). Đừng đòi hoàn hảo ngay lần đầu.

## 4.3 Ba bước xử lý SAU khi Flow tạo (bắt buộc — Flow không tự làm)
Flow cho ảnh to, nền đặc. Phải:
1. **Tách nền → trong suốt (alpha):** dùng công cụ xoá nền (remove.bg, hoặc chính
   Flow/Whisk có inpainting để xoá nền), xuất PNG có alpha. Nền phẳng đã dặn ở
   prompt giúp bước này dễ.
2. **Thu nhỏ về đúng kích thước:** 64×64 cho đa số, 128×128 nút đánh, 160×160 boss,
   240×32 thanh, v.v. (kích thước từng ảnh ở `DOT-1B`). Thu nhỏ bằng thuật toán
   giữ nét (nearest/bilinear tuỳ ảnh).
3. **Đặt đúng tên file** như trong `DOT-1B` (game gọi theo tên).

## 4.4 Tự kiểm chất lượng (đã có sẵn công cụ)
Sau mỗi lô, chạy `docs/my-thuat/KIEM-TRA-ANH.ps1` trỏ vào thư mục ảnh → nó báo ảnh
nào "phẳng"/"thiếu nền trong suốt"/"viền sáng" để làm lại. **Chỉ nhận lô khi báo
"LÔ NÀY DÙNG ĐƯỢC".** Đây là bộ lọc tự động, đỡ phải soi tay 343 ảnh.

## 4.5 Mọi tính năng cần dùng của Flow, gom lại
| Tính năng Flow | Dùng để |
|---|---|
| **Custom Tool** | Gói master prompt + phong cách + ảnh tham chiếu → tạo nhanh, đồng bộ |
| **3 ảnh tham chiếu (ingredients)** | Đồng bộ phong cách + đồng bộ một nhân vật qua nhiều ảnh |
| **Số lượng đầu ra / lần** | Tạo nhiều biến thể một lúc, chọn cái tốt nhất |
| **Aspect ratio = vuông (1:1)** | Sprite vuông cho game 2D |
| **Nano Banana editing / inpainting** | Sửa nhỏ (đổi màu giáp, xoá nền, thêm chi tiết) mà không tạo lại từ đầu |
| **Prompt expander** (nếu có) | Nở prompt ngắn thành mô tả đủ chi tiết — nhưng vẫn giữ phong cách khoá |

---

# QUY TRÌNH GỌN (in ra dán cạnh màn hình)

```
1. Tao Custom Tool: dan MASTER PROMPT (Phan 1) + gan 3 anh tham chieu (Phan 3).
2. Tao 5 "anh goc" 5 lang → duyet tay → luu lam tham chieu.
3. Chay tool theo tung dong spec trong DOT-1B (theo nhom, theo lang).
4. Moi anh: tach nen (alpha) → thu nho dung px → dat dung ten.
5. Moi lo: chay KIEM-TRA-ANH.ps1 → chi nhan khi bao DUNG DUOC.
6. Gui anh cho tro ly nap vao game.
```

---

---

# PHẦN 5 — CÁCH LÀM TỐI ƯU (cập nhật: tách nền OK, chỉ vướng kích thước)

Bạn xác nhận: **ảnh tách nền được, chỉ không thu nhỏ được.** Vậy chia việc thế này
là tối ưu nhất cho bạn — bạn lo phần sáng tạo (Flow, thứ bạn thích), tôi lo phần
cơ học (thu nhỏ + kiểm chất lượng, thứ máy làm nhanh):

```
BẠN (trong Flow):  tạo ảnh TO, đã tách nền, ĐẶT ĐÚNG TÊN (theo DOT-1B)
        │  gửi tôi folder / zip
        ▼
TÔI (tự động):     thu nhỏ đúng size + lọc ảnh phẳng  →  trả ảnh chuẩn game
```

## 5.1 Kích thước — ĐÃ GIẢI QUYẾT

Có công cụ **`THU-NHO-ANH.ps1`** (cạnh file này) + bảng **`kich-thuoc-anh.json`**
(356 ảnh, tôi trích sẵn từ DOT-1B). Nó **tự tra kích thước đích của từng ảnh theo
tên file** rồi thu nhỏ, **giữ nguyên nền trong suốt**. Bạn không phải nhập size.

- Cách dùng: chuột phải `THU-NHO-ANH.ps1` → Run with PowerShell → dán đường dẫn
  folder ảnh Flow xuất → xong, ảnh nhỏ nằm trong `da-thu-nho\`.
- Đã kiểm chạy thật: boss→160×160, logo→512×256, alpha giữ nguyên.
- **Điều kiện duy nhất: ĐẶT ĐÚNG TÊN FILE** (theo DOT-1B). Ảnh nào sai tên nó sẽ
  báo "không có trong bảng" để bạn sửa tên. Nên khi tạo trong Flow, **đặt tên ảnh
  đúng ngay** là mọi thứ tự động khớp.
- Hoặc đơn giản hơn: **bạn cứ gửi tôi folder ảnh to, tôi chạy và trả về ảnh chuẩn**
  (đúng cách mình vẫn làm — "AI làm xong gửi file zip").

## 5.2 "Nhét file của tôi + ảnh tham chiếu vào tool tự làm" — làm được tới đâu

Sự thật về Flow: Nano Banana làm **từng ảnh một, kiểu hội thoại** — nó KHÔNG tự
đọc cả file 343 dòng rồi chạy hết một mạch không cần bạn. Nên có hai mức:

**Mức A — Bán tự động, trong Flow (miễn phí, hợp với bạn):**
1. Tạo Custom Tool = Master prompt (Phần 1) + gắn 3 ảnh tham chiếu (Phần 3) cố định.
2. Mở file `DOT-1B` → **copy từng NHÓM** (ví dụ 6 dòng boss) → dán vào tool → nó ra
   6 ảnh → bạn chọn/sửa nhanh. Rồi nhóm tiếp.
   → "nhét file + ảnh tham chiếu" = tool giữ sẵn phong cách + tham chiếu, bạn chỉ
   dán danh sách theo lô. Đây là mức Flow hỗ trợ, và đủ tốt.

**Mức B — Tự động hoàn toàn (nếu sau này muốn):**
Dùng **Nano Banana qua API (Google AI Studio)** + một script tôi viết: script đọc
thẳng file DOT-1B, với mỗi dòng tự gửi (phong cách khoá + nội dung dòng + 3 ảnh
tham chiếu) → nhận ảnh → đặt đúng tên → thu nhỏ → tự kiểm. **Nhét file một phát,
chạy hết 343 ảnh không cần ngồi canh.** Đổi lại: cần khoá API (có thể tốn phí mỗi
ảnh, khác bản Flow miễn phí). Khi nào bạn muốn, nói tôi dựng cái này.

## 5.3 Khuyến nghị

Bắt đầu **Mức A** (Flow miễn phí): tạo tool + 5 ảnh gốc 5 làng → thấy ổn thì nhân
ra theo nhóm. Thu nhỏ + kiểm chất lượng để tôi lo (Mức 5.1). Nếu thấy làm tay 343
ảnh mệt, chuyển **Mức B** (tôi dựng script API) để chạy hàng loạt thật sự.

---

# Ghi chú thành thật

- Tôi xác nhận từ tài liệu Google (2/2026): Flow đã gộp Whisk+ImageFX+Flow, engine
  ảnh là Nano Banana Pro, nhận 3 ảnh tham chiếu, có Custom Tools. Nhưng **giao diện
  cụ thể của Custom Tool có thể khác chút** — chỗ nào Flow hỏi "prompt gốc" thì dán
  Phần 1, chỗ nào cho "reference/ingredient" thì gắn ảnh Phần 3.
- Nano Banana làm việc từng ảnh (hội thoại), không tự chạy cả danh sách — nên
  "tự làm hết một mạch" chỉ có ở Mức B (API + script, mục 5.2).
- Phần **nhân vật người chơi chuyển động** (đầu/thân/chân nhiều khung hình) là định
  dạng khó hơn ảnh đứng lẻ — để tính riêng sau khi bộ ảnh đứng này ổn.
