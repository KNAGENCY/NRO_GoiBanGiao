# _KHUÔN — luật + kích thước chuẩn cho MỌI ảnh game Naruto

> **File một cửa.** Bất kỳ agent nào nhận việc vẽ ảnh: đọc hết file này **trước**,
> rồi mới mở file brief của loại mình làm (ví dụ `lang-la/quai.md`).
> Cập nhật: 05/08/2026. Gộp luật từ `../HUONG-DAN-TAO-ANH.md` + `../dot-2-can-ve/` +
> cỡ thật từ `../kich-thuoc-anh.json` + cách lắp vào game (dò 05/08/2026).

---

## 0. BỐN LUẬT VÀNG — đọc 60 giây

1. **VẼ TO, ĐỪNG TỰ THU NHỎ.** Vẽ ở cỡ gốc lớn nhất vẽ đẹp được (1024px trở lên),
   **giữ đúng TỈ LỆ khung** ở bảng §2. Đừng cố vẽ 64×64 — cỡ đó AI vẽ ra bết.
   Việc thu nhỏ đã có công cụ `TAO-4-MUC-ZOOM.py` lo (§4).
2. **Nền trong suốt thật** (PNG có kênh alpha) — trừ ảnh nền bản đồ `map_nen_*`.
3. **Tranh vẽ có KHỐI, có ĐỔ BÓNG** — không phải hình phẳng vài màu, chân que.
   Nếu file xuất ra **nhẹ dưới 3 KB** là sai, vẽ lại. Xem `anh-tham-chieu/TOT-*.png`.
4. **BÁO CÁO NÓI THẬT.** Lô 343 từng tự nhận "QA sạch" trong khi 62 ảnh sai cỡ.
   Chỉ ghi "đã kiểm X" khi **thật sự chạy công cụ kiểm** (§5). Không chắc thì ghi "không chắc".

---

## 1. KÍCH THƯỚC — luật

> **Vẽ ở cỡ gốc lớn, giao nguyên cỡ đó, KHÔNG thu nhỏ.** Miễn giữ đúng **tỉ lệ khung**
> ghi ở §2. Ví dụ khung chuẩn 96×128 (tỉ lệ 3:4) thì vẽ 768×1024 hay 1152×1536 đều được,
> **không** vẽ 1024×1024. Hậu tố `_x4` trong tên là **khoá tra cứu**, không phải cỡ thật —
> giữ nguyên kể cả khi giao ảnh lớn hơn.

## 2. BẢNG CỠ CHUẨN TỪNG LOẠI (cỡ x4 — lấy từ `kich-thuoc-anh.json`)

| Loại (tiền tố tên) | Cỡ chuẩn x4 | Tỉ lệ | Là gì |
|---|---|---|---|
| `map_` (tile) | **96×96** | 1:1 | viên gạch bản đồ (đất, cỏ, tường...) |
| `map_` (vật thể) | 192×192 | 1:1 | vật thể lớn (cổng làng, cây to) |
| `map_nen_` (nền) | 384×192 | 2:1 | nền trời/núi — **nền ĐẶC, không trong suốt** |
| `tp_` đầu | **128×128** | 1:1 | trang phục nhân vật — phần ĐẦU |
| `tp_` thân | **128×160** | 4:5 | trang phục nhân vật — phần THÂN |
| `npc_` | **96×128** | 3:4 | NPC đứng nguyên con |
| `cd_` | **200×200** | 1:1 | chân dung (mặt) cho khung hội thoại |
| `quai_` | **96×96** | 1:1 | quái thường |
| `boss_` | **160×160** | 1:1 | boss |
| `vp_` | **64×64** | 1:1 | vật phẩm (áo, đồ, nguyên liệu) |
| `vk_` kunai/phi tiêu | **96×48** | 2:1 | vũ khí dài |
| `vk_` shuriken | 64×64 | 1:1 | vũ khí tròn nhỏ |
| `kn_` | **64×64** | 1:1 | icon chiêu thức (kỹ năng) |
| `hu_` | **1200×480** | dải 5 khung 240×480 | hiệu ứng (sprite sheet) |
| `ui_` ô/nút nhỏ | 64×64 · 96×48 | | ô vật phẩm, nút bấm |
| `ui_` thanh máu/chakra | **240×32** | 15:2 | tách 2 lớp `_nen` + `_fill` (§7) |
| `ui_` khung/nền lớn | 480×160 · 480×480 · 1024×576 | | khung hội thoại, túi đồ, nền đăng nhập |

> **Ảnh nhiều khung hình** (`hu_`): chiều rộng **mỗi khung phải chia hết cho 4**
> (240 ✓), nếu không tới mức x1 các khung lệch nửa điểm ảnh, hoạt ảnh giật mép.
> Ghi rõ **số khung** khi giao.
>
> **Tên mới chưa có trong bảng:** báo để thêm vào `../kich-thuoc-anh.json` **kèm cỡ chuẩn**,
> đừng tự chế tên/cỡ — game tra ảnh theo đúng tên.

## 3. BẢNG MÀU 5 LÀNG (áp cho map, nhân vật, quái, icon chiêu của làng đó)

| Làng | Hệ | Màu chủ đạo |
|---|---|---|
| **Lá** | Hoả | cam–đỏ, ánh lửa ấm |
| **Cát** | Phong | vàng cát–xanh lục nhạt, xoáy gió |
| **Sương Mù** | Thuỷ | xanh dương–xanh lam, nước/sương |
| **Đá** | Thổ | nâu đất–xám đá |
| **Mây** | Lôi | tím–trắng xanh, tia sét |

## 4. NỀN TRONG SUỐT

- Mọi ảnh là **PNG có kênh alpha**, nền trong hoàn toàn. Ngoại lệ: `map_nen_*` (nền đặc).
- **Không viền sáng/tối** quanh mép (dấu cắt nền ẩu). Kiểm bằng cách đặt ảnh lên nền
  trắng **và** nền đen — mép phải sạch cả hai.

## 5. TÊN FILE + ID (game tra theo tên & số thứ tự)

- Lấy **đúng tên** trong `../kich-thuoc-anh.json`, không tự đổi. Sai tên = game không thấy ảnh.
- Giữ hậu tố `_x4` (khoá tra cứu). Công cụ tự đổi hậu tố khi sinh các mức.
- **ID ảnh append-only**: thêm ảnh mới luôn ở **CUỐI** danh sách, **không chèn giữa,
  không đổi thứ tự** — chèn giữa là lệch toàn bộ ảnh phía sau.

## 6. QUY TRÌNH — vẽ xong làm gì

```
  Vẽ ở cỡ GỐC LỚN, đúng tỉ lệ, nền trong suốt
                 │
                 ▼
  python ..\KIEM-TRA-ANH.py <thư-mục-của-bạn> nguon   ← tự kiểm, DÁN kết quả vào báo cáo
                 │  (kiểm 7 điều: tên · thiếu/thừa · cỡ · alpha · ảnh phẳng · trống · file nhẹ)
                 ▼
  Giao thư mục phẳng PNG đúng tên + báo cáo thật
                 │
                 ▼   (phía dự án chạy)
  python ..\TAO-4-MUC-ZOOM.py <thư-mục-giao> ..\anh-game   ← sinh đủ x4/x3/x2/x1
                 │
                 ▼
  Lắp vào game (§7)
```
> ⚠️ **Đừng dùng `KIEM-TRA-ANH.ps1`** (bản cũ) — chú thích ghi có kiểm cỡ nhưng mã
> KHÔNG kiểm; chính nó báo "sạch" cho lô 343 khi 62 ảnh sai cỡ. Dùng bản `.py`.

## 7. CÁCH LẮP VÀO GAME (khuôn dò 05/08/2026 — để hiểu ảnh đi đâu)

Ảnh **ưu tiên đọc từ Unity Resources** (thắng cache cũ). Client hiện chạy **zoom 2** → đặt ảnh vào `res/x2/…` (và đủ x1–x4). Chi tiết ở bộ nhớ dự án `map-render-client`, `logo-resource-sync`.

- **Tile map** → `Assets/Resources/res/x{zoom}/t/{tileID}$V.png` (100 ảnh rời/1 tileset, mỗi ô 48px ở x2). Ô nào đặt gạch nào do sơ đồ `res/mymap/{mapID}.bytes`.
- **Nền map** → `res/x{zoom}/bg/b{typeBg}{k}.png` (4 lớp b0..b3, xếp dọc, lặp ngang). ⚠️ **Ảnh nền phải bật `Read/Write Enabled`** trong Unity, nếu không màu trời tính ra 0 = **trời đen**.
- **Logo** → sửa **NGUỒN** `Assets/Logo/1.png` (KHÔNG sửa bản Resources — có script tự copy đè).
- Người chơi thật chỉ thấy sau khi **build lại client** (APK/exe).

## 8. ẢNH THAM CHIẾU — mở `anh-tham-chieu/` ra xem

Luật bằng chữ không đủ, phải nhìn. `TOT-*.png` = chuẩn cần đạt (lông/vải có khối, nền cắt sạch). `XAU-*.png` = lô bị trả (phẳng, chân que) — **tránh tuyệt đối**. Mỗi file brief sẽ chỉ đúng ảnh tham chiếu cho loại của nó.

---

**Tóm tắt cho agent:** đọc file này → mở brief loại mình → vẽ TO đúng tỉ lệ §2, nền trong,
có khối → `KIEM-TRA-ANH.py` → giao kèm báo cáo thật. Hết.
