# Kế hoạch tạo TOÀN BỘ ảnh game Naruto (thiết kế)

**Ngày:** 05/08/2026 · **Trạng thái:** thiết kế (chờ chủ dự án soát) · **Nhánh:** `map-lang-la`
**Bám format đợt cũ:** `docs/my-thuat/HUONG-DAN-TAO-ANH.md` + `dot-2-can-ve/`. Công cụ sẵn: `KIEM-TRA-ANH.py`, `TAO-4-MUC-ZOOM.py`, bảng `kich-thuoc-anh.json`.

## 1. Mục tiêu

Sinh ra **một bộ file MD "brief" hoàn chỉnh** để tạo mọi ảnh còn thiếu cho game Naruto (5 làng), chia **theo LOẠI × LÀNG** (nhiều file nhỏ). Mỗi file là một đề bài tự chứa: danh sách ảnh + ảnh tham chiếu + kích thước — **sẵn sàng kéo cho một agent vẽ ngoài** (AI/thợ), năng suất ~100 ảnh/agent/ngày.

## 2. Vai trò — ai làm gì (chốt với chủ dự án)

| Ai | Làm gì |
|---|---|
| **Phiên này (Claude)** | Tạo cây thư mục + **viết tất cả file MD brief** + `_KHUON.md` + gom `anh-tham-chieu/`. **KHÔNG** gen ảnh, **KHÔNG** spawn/điều phối agent, **KHÔNG** lắp ảnh. |
| **Chủ dự án** | Cầm từng file MD + thư mục tham chiếu → giao cho agent vẽ ngoài. |
| **Agent vẽ ngoài** (không thuộc phiên này) | Đọc brief → vẽ ảnh đúng khuôn → tự kiểm `KIEM-TRA-ANH.py` → giao. |
| **Đợt sau** | Lắp ảnh đã vẽ vào game (đường map/nhân vật đã dò — bộ nhớ `map-render-client`). |

## 3. Cấu trúc thư mục

```
docs/my-thuat/dot-3-toan-bo/
├── _KHUON.md              # bảng cỡ + phong cách 5 làng + quy trình (mục 6)
├── README.md              # mục lục đợt: đọc gì trước, trạng thái từng file
├── anh-tham-chieu/        # tái dùng 16 ảnh mẫu cũ + bổ sung mỗi loại 1 mẫu từ kho 343
├── chung/                 # KHÔNG theo làng
│   ├── ui.md · thanh-mau-chakra.md · hieu-ung.md · vat-pham-chung.md
├── lang-la/               # mỗi làng 7 file theo loại
│   ├── map.md · nhan-vat.md · quai.md · boss.md · npc.md · trang-phuc.md · icon-chieu.md
├── lang-cat/  ·  lang-suong-mu/  ·  lang-da/  ·  lang-may/   # (7 file mỗi làng)
```
Tổng ~**39 file MD** (4 chung + 5×7 làng). Mỗi file ≤ ~100 ảnh; loại nào > 100 (vd tile map, khung hình nhân vật) thì **tách đôi** (`map-1.md`, `map-2.md`).

## 4. Danh mục + trạng thái (đánh dấu để agent không vẽ thừa)

| Nhóm | Trạng thái | Ghi chú |
|---|---|---|
| `chung/` (UI, thanh máu, hiệu ứng, vật phẩm chung) | 🟢 **Làm ngay** | Server đã chốt; nhiều cái chỉ thiếu bản tách lớp / bổ sung |
| `lang-la/` (đủ 7 loại) | 🟢 **Làm ngay** | Nội dung Làng Lá đã rõ (map đang cắm, NPC Iruka, quái Sói...) |
| `lang-cat/`, `lang-suong-mu/`, `lang-da/`, `lang-may/` | 🟡 **Viết khung, chờ số** | Chờ server chốt (KH7 vật phẩm + nội dung 4 làng) mới điền số lượng chính xác. `icon-chieu` 4 làng này **làm ngay được** (20 chiêu server đã có). |

Mỗi file 🟡 vẫn viết đủ khung (phong cách, tham chiếu, khuôn) + một dòng cảnh báo **"CHỜ: server chưa chốt bao nhiêu bộ/quái — đừng vẽ số lượng, chỉ vẽ mẫu đại diện"**.

**Quan hệ với đợt 2 (chống trùng):** `dot-2-can-ve/` đã liệt kê **4 thanh máu tách lớp + 20 icon chiêu (5 làng)**. Nếu đợt 2 CHƯA vẽ xong, đợt 3 **trỏ về đợt 2, KHÔNG lặp** — `chung/thanh-mau-chakra.md` và các `*/icon-chieu.md` chỉ ghi *"xem `../dot-2-can-ve/`"* + bổ sung phần còn thiếu. Bước đầu khi thực thi phải **kiểm `anh-game/` xem đợt 2 đã có ảnh chưa** rồi mới quyết trỏ hay điền.

## 5. Format mỗi file MD brief (đúng khuôn `dot-2-can-ve`)

1. Tiêu đề + số ảnh + trạng thái 🟢/🟡.
2. **4 luật bắt buộc** (vẽ to đừng thu nhỏ · nền trong suốt · có khối có bóng · báo cáo thật) — trỏ về `../_KHUON.md`.
3. **Phong cách** + bảng ảnh tham chiếu (TOT/XAU) áp cho loại này + màu hệ của làng.
4. **Bảng danh sách ảnh:** `tên_x4.png` · nội dung/gợi ý hình · kích thước · tỉ lệ khung.
5. **Tự kiểm:** lệnh `KIEM-TRA-ANH.py` + yêu cầu dán kết quả.
6. **Giao thế nào** + "chưa cần vẽ gì đợt này".

## 6. `_KHUON.md` chứa (phần chủ dự án dặn "làm khuôn")

- **Bảng kích thước chuẩn từng loại** (gom từ `kich-thuoc-anh.json`): nhân vật 96×128, tile map 96 (x4), icon chiêu 64×64, boss 160×160, npc, vật phẩm, thanh máu 240×32... + loại mới bổ sung.
- **Bảng màu 5 làng** (Lá=hoả cam-đỏ · Cát=phong vàng-lục · Sương Mù=thuỷ xanh · Đá=thổ nâu-xám · Mây=lôi tím).
- **Quy trình 1 dòng:** vẽ to → `KIEM-TRA-ANH.py` → `TAO-4-MUC-ZOOM.py` → lắp game.
- **Cách lắp vào game** (khuôn map dò hôm nay): tile = 100 ảnh `/t/{id}$V.png` 48px, nền = `bg/b0k`, bố cục `mymap`, ảnh ưu tiên Resources, nền cần `isReadable`. (Bộ nhớ `map-render-client`, `logo-resource-sync`.)

## 7. Ảnh tham chiếu

- Tái dùng 16 ảnh trong `dot-2-can-ve/anh-tham-chieu/` (TOT/XAU + THAM-CHIEU).
- **Bổ sung**: mỗi loại chưa có mẫu thì trích 1 ảnh đại diện từ kho 343 (`anh-game/x4/`) làm "mẫu tốt" (vd `TOT-map-la-conglang.png`, `TOT-nhanvat-mau.png`). Chỉ copy, không vẽ mới.

## 8. Ràng buộc (bất biến)

- Mỗi file brief ≤ ~100 ảnh (1 agent/ngày). Vượt thì tách.
- Tên ảnh **lấy đúng** `kich-thuoc-anh.json`; tên mới phải **bổ sung vào bảng kèm cỡ**, không tự chế.
- ID ảnh **append-only** (thêm cuối, không chèn giữa — lệch toàn bộ ảnh sau).
- Nền trong suốt (trừ `map_*`). Vẽ ở cỡ gốc lớn, đừng thu nhỏ.

## 9. Ngoài phạm vi (YAGNI)

- Phiên này **không vẽ/gen ảnh**, không chạy agent vẽ, không lắp ảnh vào game.
- Không chốt số lượng nội dung 4 làng (chờ server) — chỉ viết khung.
- Không đụng code server/client.

## 10. Thứ tự làm (khi được duyệt)

1. `_KHUON.md` + `README.md` + gom `anh-tham-chieu/` (nền tảng).
2. `chung/` (4 file) — 🟢.
3. `lang-la/` (7 file) — 🟢.
4. `lang-cat/ · lang-suong-mu/ · lang-da/ · lang-may/` (28 file khung) — 🟡, riêng `icon-chieu` điền luôn.
