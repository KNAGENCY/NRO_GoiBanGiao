# HƯỚNG DẪN TẠO ẢNH — dùng cho MỌI đợt sau

> **Đây là file dành cho AI vẽ ảnh.** Đọc hết trước khi vẽ tấm đầu tiên.
> File này thay cho `DOT-1-PROMPT-TAO-ANH.md` và `DOT-1B-LAM-LAI-343-ANH.md`
> ở phần **luật chung** (hai file kia vẫn giữ để tra danh sách 343 ảnh cũ).
>
> **Cập nhật:** 27/07/2026, sau khi nghiệm thu lô 343 ảnh.

---

## 0. ĐỌC 60 GIÂY — bốn điều quan trọng nhất

1. **VẼ TO, ĐỪNG TỰ THU NHỎ.** Vẽ ở cỡ gốc lớn nhất mà bạn vẽ đẹp được
   (1024×1024 hay hơn nữa). **Tuyệt đối đừng cố vẽ ảnh 96×96** — cỡ đó AI vẽ
   ra là bết. Việc thu nhỏ đã có công cụ lo, xem §2.
2. **Nền phải trong suốt thật** (PNG có kênh alpha), trừ ảnh nền bản đồ.
3. **Phải là tranh vẽ có khối, có đổ bóng** — không phải hình khối phẳng vài
   màu. Xem ảnh mẫu ở §4, đây là phần quan trọng nhất file này.
4. **Báo cáo phải nói thật.** Lô trước tự nhận "QA sạch, không sai kích thước"
   trong khi 62 ảnh sai cỡ — xem §3. Không chắc thì ghi "không chắc", đừng ghi
   "đã kiểm".

---

## 1. Vì sao lô 343 phải sửa — ĐỌC ĐỂ KHÔNG LẶP LẠI

Lô 343 ảnh (giao 27/07/2026) **chất lượng vẽ rất tốt** — nét, có khối, nền cắt
sạch, phong cách nhất quán. Nhưng có **5 lỗi quy trình** phải tránh:

| # | Lỗi thật đã xảy ra | Vì sao tệ | Lần sau làm gì |
|---|---|---|---|
| 1 | **62/343 ảnh để nguyên cỡ gốc (gấp 4 lần), 281 ảnh còn lại đã thu nhỏ** | Không phải sai cỡ mới đáng ngại — mà là **không nhất quán**. Người nhận không biết ảnh nào đã xử lý, ảnh nào chưa | **Giao TẤT CẢ ở cỡ gốc lớn, đừng thu nhỏ ảnh nào.** Một quy tắc, không ngoại lệ |
| 2 | **README tự nhận "Không phát hiện sai kích thước"** trong khi 62 ảnh sai | Báo cáo sai còn tệ hơn không báo cáo — nó tạo niềm tin giả, người nhận bỏ qua bước kiểm | Chỉ ghi "đã kiểm X" khi **thật sự chạy công cụ kiểm X**. Nêu rõ đã chạy lệnh nào |
| 3 | **Thanh máu / thanh chakra giao bản GỘP SẴN** (`ui_thanh_theluc_x4.png`) thay vì tách 2 lớp (`_nen` + `_fill`) | Game vẽ thanh máu bằng cách **cắt lớp fill theo phần trăm** rồi chồng lên lớp nền. Ảnh gộp sẵn thì máu **luôn đầy**, không bao giờ vơi được | Ảnh nào ghi rõ `_nen` / `_fill` trong tên thì **phải là hai tấm rời** |
| 4 | **`hu_anthan` mờ gần như vô hình** (chỉ ~57 điểm ảnh đục trên 1200×480) | Hiệu ứng mà người chơi không nhìn thấy thì bằng không | Hiệu ứng "mờ" vẫn phải **nhìn thấy được**: khung đầu tiên nên rõ ≥60% độ đục |
| 5 | Bảng cỡ có mục rác tên rỗng `_x4.png` | Công cụ đếm nhầm, báo thiếu ảnh không tồn tại | Đừng thêm dòng vào `kich-thuoc-anh.json` nếu chưa có ảnh thật |

**Lỗi cũ hơn (lô v3, đã sửa ở lô 343 — vẫn phải tránh):** vẽ ra "hình khối
phẳng vài màu, chân que". Dấu hiệu nhận biết: **file nhẹ dưới ~3 KB**. Ảnh thật
có khối, có bóng thì nặng hàng chục KB.

---

## 2. KÍCH THƯỚC — luật mới, đọc kỹ

### Luật

> **Vẽ ở cỡ gốc lớn nhất bạn vẽ đẹp được. Giao nguyên cỡ đó. KHÔNG thu nhỏ.**

Cỡ gốc bao nhiêu cũng được (1024×1024, 1536×1536, tuỳ bạn) — miễn là **giữ đúng
tỉ lệ khung** ghi trong `kich-thuoc-anh.json`. Ví dụ ảnh chuẩn là 96×128 (tỉ lệ
3:4) thì bạn vẽ 768×1024 hay 1152×1536 đều được, **không** được vẽ 1024×1024.

### Vì sao

- AI vẽ ở 96×96 cho ra ảnh bết, mất chi tiết — không cứu được.
- Thu nhỏ từ ảnh lớn thì **sắc nét**, và làm được **bất kỳ lúc nào** bằng công cụ.
- Ngược lại thì không: đã nhỏ rồi phóng to lại là mờ vĩnh viễn.

### Ai lo phần thu nhỏ

Công cụ `TAO-4-MUC-ZOOM.py` (cùng thư mục này). Nó tự tra cỡ chuẩn theo tên
file rồi sinh **đủ 4 mức** game cần:

| Mức | Tỉ lệ | Dùng khi |
|---|---|---|
| `x4` | 100% cỡ chuẩn | màn hình lớn |
| `x3` | 75% | |
| `x2` | 50% | |
| `x1` | 25% | máy yếu / màn nhỏ |

```powershell
python TAO-4-MUC-ZOOM.py <thư-mục-ảnh-bạn-giao> anh-game
```

**Thiếu một mức là người chơi ở cỡ đó mất ảnh trong im lặng** — máy chủ cũ nuốt
lỗi, không báo gì. Nên đừng bao giờ giao thủ công từng mức; luôn để công cụ sinh.

### Ảnh nhiều khung hình (sprite sheet)

Ảnh hiệu ứng thường là **một dải nhiều khung** (`hu_anthan_x4.png` = 5 khung
240×480 ghép ngang). Với loại này:

- **Chiều rộng mỗi khung phải chia hết cho 4**, nếu không tới mức `x1` các khung
  sẽ lệch nửa điểm ảnh và hoạt ảnh bị giật mép.
- Ghi rõ **số khung** khi giao, đừng để người nhận đoán.

---

## 3. NỀN TRONG SUỐT

- **Mọi ảnh phải là PNG có kênh alpha**, nền trong suốt hoàn toàn.
- **Ngoại lệ duy nhất: ảnh nền bản đồ** (`map_*`) — loại này nền đặc là đúng,
  vì nó chính là nền.
- Không để **viền sáng/viền tối** quanh mép hình (dấu vết cắt nền ẩu). Kiểm
  bằng cách đặt ảnh lên nền trắng **và** nền đen, mép phải sạch ở cả hai.

---

## 4. ẢNH THAM CHIẾU — MỞ RA XEM, ĐÂY LÀ PHẦN QUAN TRỌNG NHẤT

Trong thư mục **`anh-tham-chieu/`**. Luật bằng chữ không đủ — phải nhìn.

### ✅ MẪU TỐT — chuẩn cần đạt

| File | Là gì | Học điều gì |
|---|---|---|
| `TOT-quai-soi-rung.png` | Quái sói, lô 343 | **Chuẩn vàng cho quái**: lông có lớp, cơ bắp có khối, mắt có thần, nền cắt sạch không viền |
| `TOT-npc-ong-tazuna.png` | NPC ông già, lô 343 | **Chuẩn vàng cho người**: vải có nếp gấp, đổ bóng đúng hướng, phụ kiện chi tiết (nón, bình rượu) |
| `TOT-boss-mau.png` | Boss Haku, lô 343 | Mức chi tiết cần có ở boss |
| `TOT-vatpham-mau.png` | Vật phẩm áo, lô 343 | Vật phẩm nhỏ vẫn phải có khối và bóng, không phẳng |
| `TOT-naruto-logo.png`, `TOT-naruto-khung.png`, `TOT-naruto-thanh-chakra.png` | Giao diện đã duyệt | Phong cách gỗ sơn mài tối + viền cam |
| `TOT-chitiet-nhanvat.png`, `TOT-chitiet-hieuung.png` | Tham chiếu **kiểu vẽ** (không phải nội dung) | Mức chi tiết và đổ bóng cần đạt |

### ❌ MẪU XẤU — tuyệt đối tránh

`XAU-boss-zabuza.png`, `XAU-npc.png`, `XAU-vatpham.png` — chính là lô v3 bị trả
lại. Phẳng, vài màu, chân que. Nếu ảnh bạn vẽ ra **trông giống mấy tấm này**,
hoặc **file nhẹ dưới 3 KB**, thì làm lại.

---

## 5. ĐẶT TÊN FILE

- Lấy **đúng tên** ghi trong `manifest-tao-anh.json` / `kich-thuoc-anh.json`,
  không tự đổi. Sai tên = game không tìm thấy ảnh (game tra theo tên).
- Giữ hậu tố `_x4` trong tên kể cả khi bạn giao ảnh cỡ gốc lớn — đó là **khoá
  tra cứu**, không phải mô tả cỡ thật. Công cụ tự đổi hậu tố khi sinh các mức.
- Tên ảnh mới (chưa có trong bảng): báo trước để thêm vào `kich-thuoc-anh.json`
  **kèm cỡ chuẩn**, đừng tự thêm.

---

## 6. TỰ KIỂM TRƯỚC KHI GIAO — bắt buộc

Chạy công cụ kiểm, **dán nguyên kết quả vào báo cáo**:

```powershell
python KIEM-TRA-ANH.py <thư-mục-ảnh-của-bạn> nguon
```

Nó kiểm thật 7 điều: tên có trong manifest · thiếu/thừa · **kích thước** ·
kênh trong suốt · ảnh phẳng dưới 40 màu · ảnh trống · file nhẹ bất thường.

> ⚠️ **Đừng dùng `KIEM-TRA-ANH.ps1` (bản PowerShell cũ)** — phần chú thích của
> nó ghi là có kiểm kích thước nhưng **mã thì không hề kiểm**. Chính nó đã báo
> "sạch" cho lô 343 trong khi 62 ảnh sai cỡ. Giữ lại chỉ để tham khảo lịch sử.

**Báo cáo khi giao phải có:** số ảnh giao · lệnh kiểm đã chạy + kết quả nguyên
văn · ảnh nào không chắc chắn và vì sao. **Không suy đoán là "đạt".**

---

## 7. VIỆC CÒN THIẾU CỦA LÔ 343 — làm trước tiên nếu bạn nhận đợt tiếp

| Việc | Chi tiết |
|---|---|
| **4 ảnh thanh máu/chakra tách lớp** | Cần `ui_thanh_theluc_nen_x4.png` + `ui_thanh_theluc_fill_x4.png` + `ui_thanh_chakra_nen_x4.png` + `ui_thanh_chakra_fill_x4.png`, mỗi tấm 240×32 (vẽ ở cỡ gốc lớn theo tỉ lệ 15:2). Hiện chỉ có bản **gộp** không dùng được — xem §1 lỗi #3. **Lớp nền**: rãnh lõm tối bo tròn hai đầu, ruột TRỐNG. **Lớp fill**: dải màu đặc (thể lực cam-đỏ, chakra xanh cyan), hai mép **cắt phẳng vuông** để cắt theo % không bị hụt |
| **Vẽ lại `hu_anthan`** | Hiện mờ gần như vô hình — xem §1 lỗi #4 |
| **Dọn `kich-thuoc-anh.json`** | Xoá mục rác tên rỗng `_x4.png` |

*(8 ảnh giao diện đã duyệt còn lại nằm ở `my thuat/dot1-nhom1-duyet/` — đã kiểm
đúng cỡ, đúng alpha, **không cần vẽ lại**.)*

---

## 8. TÓM TẮT QUY TRÌNH

```
   AI vẽ ở cỡ GỐC LỚN, đúng tỉ lệ khung, nền trong suốt
                    │
                    ▼
   python KIEM-TRA-ANH.py <thư-mục> nguon      ← AI tự kiểm, dán kết quả vào báo cáo
                    │
                    ▼
   python TAO-4-MUC-ZOOM.py <thư-mục> anh-game ← sinh đủ x4/x3/x2/x1
                    │
                    ▼
   python KIEM-TRA-ANH.py anh-game/x4          ← kiểm lại sau khi thu nhỏ
                    │
                    ▼
   Đưa vào game (gán ID ảnh — xem quy tắc mapping ở BAT-DAU-TU-DAY.md)
```

**Quy tắc mapping quan trọng khi đưa vào game:** ảnh được game gọi bằng **số ID**,
và số đó là **vị trí trong danh sách**. Nên khi thêm ảnh mới: **luôn thêm ID ở
CUỐI, không chèn vào giữa, không đổi thứ tự** — chèn giữa là lệch toàn bộ ảnh
phía sau.
