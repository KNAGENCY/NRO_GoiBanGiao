# ĐỢT 3 — TOÀN BỘ ẢNH 5 LÀNG (mục lục)

> **Mở rộng từ 343 ảnh đợt 1.** Đợt 1 (`../DOT-1B-LAM-LAI-343-ANH.md`) đã phủ **Làng Lá + đồ dùng chung**. Đợt 2 (`../dot-2-can-ve/`) = 20 icon chiêu + 4 thanh máu. **Đợt 3 = 4 làng còn lại + phần thiếu.**

## Agent vẽ đọc gì trước

1. **`_KHUON.md`** — luật vàng, bảng cỡ chuẩn, màu 5 làng, **master prompt phong cách khoá**, cách lắp game. ĐỌC TRƯỚC.
2. **`anh-tham-chieu/`** — ảnh mẫu TỐT/XẤU (mở ra xem).
3. File brief của loại mình làm (ví dụ `lang-cat/quai.md`).

## Cấu trúc — chia LOẠI × LÀNG

```
chung/         ← không theo làng (phần lớn đã có ở 343; đợt 3 chỉ bổ sung chỗ thiếu)
lang-la/       ← đã có ở 343; đợt 3 chỉ bổ sung
lang-cat/  lang-suong-mu/  lang-da/  lang-may/   ← KHỐI CHÍNH đợt 3
   mỗi làng: map · nhan-vat · trang-bi · quai · boss · npc · icon-chieu
```

## Trạng thái (🟢 làm ngay / 🟡 chờ server chốt nội dung)

| Phần | Trạng thái | Ghi chú |
|---|---|---|
| `_KHUON.md`, mẫu `lang-cat/quai.md` | ✅ xong | nền tảng + mẫu format |
| **Bảng cỡ `../kich-thuoc-anh.json`** | ✅ **đã thêm 316 tên đợt 3** (11/08) | 375 → **691 mục**, append-only. `KIEM-TRA-ANH.py` tra được mọi tên đợt 3 ngay. Bản cũ giữ ở `.bak-truoc-dot3`. Các file brief còn ghi "báo thêm vào kich-thuoc-anh.json" — **việc đó ĐÃ XONG**, agent vẽ khỏi làm |
| `lang-*/trang-bi`, `nhan-vat`, `map` | 🟢 | nhân cấu trúc 343 (chắc chắn cần) |
| `lang-*/quai`, `boss`, `npc` | 🟡 | bộ **đề xuất đại diện** — server chưa chốt nội dung 4 làng; chốt lại sau |
| `lang-*/icon-chieu` | ↪ trỏ đợt 2 | 20 icon `kn_*` đã ở `../dot-2-can-ve/`, không lặp |

## Mỗi file brief gồm

Tiêu đề + trạng thái → phong cách hệ làng + ảnh tham chiếu → **subject block** tiếng Anh → **bảng ảnh** (`tên_x4.png` · cỡ · mô tả chi tiết) → tự kiểm `KIEM-TRA-ANH.py`. Đúng format `../DOT-1B-LAM-LAI-343-ANH.md`.

## Quy trình một agent vẽ

`_KHUON.md` → mở brief loại mình → gen ảnh (master prompt + subject block + dòng mô tả) → tách nền + thu nhỏ (`../TAO-4-MUC-ZOOM.py`) → `KIEM-TRA-ANH.py` → giao. Tên đúng, ID thêm ở cuối.
