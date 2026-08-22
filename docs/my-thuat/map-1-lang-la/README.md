# MAP 1 — LÀNG LÁ: đề bài tạo ảnh (dành cho AI tạo ảnh / thợ vẽ)

> **Thư mục này tự chứa đủ.** Mở `anh-tham-chieu/` xem **trước khi vẽ tấm đầu tiên** — đặc biệt là `THUOC-TI-LE.png`.
> Luật chung toàn dự án: `../dot-3-toan-bo/_KHUON.md`.
>
> **Cập nhật 22/08/2026 — BẢN 2.** Bản 1 ghi **sai kích thước** (nhà 192px thấp hơn nhân vật 230px). Bản này đã sửa theo tỉ lệ đo từ chính sprite trong game.

---

## 0. BỐ CỤC ĐÃ CHỐT — map có 2 TẦNG

```
          ☁️  trời + mây                          ← lặp ngang, xa nhất
       🗿 NÚI HOKAGE (mờ, xám nhạt)                ← lặp ngang
    🏘️ LÀNG XA — nhà mờ, nhạt màu, KHÔNG đi tới    ← TẦNG PHÔNG
   ▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔
    🏠 NHÀ GẦN 🏠  ⛩️ cổng  🏮 đèn  🪣 giếng        ← TẦNG CHƠI ĐƯỢC
   ▔▔▔▔▔▔🧍▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔  ← nhân vật đi ở đây (mặt đất)
   ░░░░░ đất đặc ░░░░░░░░░░
    🌿 cỏ/cành sẫm màu che rìa                     ← TIỀN CẢNH
```

- **Tầng phông** (làng xa + núi + trời): mờ, nhạt, lặp ngang được. Chỉ để ngắm.
- **Tầng chơi được**: nhà + đồ vật đứng **trên cùng mặt đất nhân vật đi**. Rõ nét, đậm màu.
- **Tiền cảnh**: sẫm gần đen, che rìa trên/dưới màn hình.

---

## 1. TỈ LỆ — ĐỌC KỸ, ĐÂY LÀ CHỖ DỄ SAI NHẤT

Mở `anh-tham-chieu/THUOC-TI-LE.png`. Mọi cỡ tính theo **ô**: **1 ô = 96px** (mức x4).

| Vật | Cao bao nhiêu ô | = px |
|---|---|---|
| **Nhân vật người chơi** | 2,4 ô | ~230px |
| Quái sói | 1 ô | 96px |
| NPC (Iruka) | 1,33 ô | 128px |
| **Nhà nhỏ** | **6 ô** | **576px** |
| **Nhà lớn** | **8 ô** | **768px** |
| Cổng torii | 5 ô | 480px |
| Đèn lồng | 4 ô | 384px |
| Cây | 8 ô | 768px |

> Nhà phải cao **gấp 2,5–3 lần nhân vật**. Nếu vẽ nhà chỉ 2 ô thì nhân vật cao hơn cả nhà — sai hoàn toàn.

**Luật kỹ thuật — sai là game hiển thị hỏng:**

1. **Cỡ ảnh luôn là bội số của 96** (96, 192, 288, 384, 480, 576, 768…). Cỡ lẻ → Unity tự phóng méo, mất một phần ảnh. *(Đã xảy ra thật: ảnh 96px bị phóng lên 128 làm mất hẳn cột phải của cổng torii.)*
2. **Nền trong suốt thật** (PNG có alpha) cho mọi vật thể. Ngoại lệ: `nen_troi` (§3.1) nền đặc.
3. **Chân vật thể chạm ĐÚNG mép dưới ảnh** — không hở, không thừa. Hở là vật thể lơ lửng trong game.
4. **Không viền sáng/tối quanh mép**. Kiểm bằng cách đặt lên nền trắng VÀ nền đen.
5. Vẽ **nhìn ngang** (side view), không phối cảnh 3/4, không nhìn từ trên xuống.

---

## 2. ẢNH THAM CHIẾU — mở ra xem

| File | Học gì |
|---|---|
| **`THUOC-TI-LE.png`** | ⭐ **Quan trọng nhất** — vật nào cao bao nhiêu ô. Gắn kèm MỌI lần tạo ảnh |
| `concept-map1-langla.png` (thư mục cha) | **Concept đã duyệt về phong cách** — bám đúng tông màu, nét vẽ, chất liệu này |
| `HIEN-TRANG-map1-trong-game.png` | Map hiện tại — thứ cần thay. Thấy rõ: nhà bị chôn, 3 nhà giống hệt cách đều, đất phẳng lì |
| `GIAI-THICH-tileset-vien-chuyen-tiep.png` | Vì sao cần viên viền — trái sai (cắt thẳng), phải đúng (cỏ lưa tưa phủ xuống) |
| `TOT-cong-torii-da-tach-nen.png` | Chuẩn vàng vật thể rời: nguyên vật thể, nền trong suốt sạch, chân chạm mép dưới |
| `HIENCO-*.png` | Art đã có trong game — bám tông màu |
| ❌ `XAU-*.png` | Lô từng bị trả lại: phẳng, vài màu, chân que. Tránh tuyệt đối |

---

## 3. DANH SÁCH ẢNH CẦN TẠO — 27 ảnh

### 3.1 — TẦNG PHÔNG (lặp ngang được)

| Tên file | Cỡ | Nội dung |
|---|---|---|
| `nen_troi_x4.png` | 768×384 | Trời xanh + mây. **Mép trái–phải nối liền nhau**. Nền ĐẶC. Nửa trên sáng (game lấy màu ở đây làm màu trời) |
| `nen_nuihokage_x4.png` | 1536×768 | Vách núi khắc mặt Hokage. **Nửa dưới là núi, nửa trên TRONG SUỐT**. Mép trái–phải nối liền. Xám nhạt, mờ (ở xa) |
| `nen_lang_xa_x4.png` | 1536×576 | **Dãy mái nhà làng nhìn từ xa** — mờ, nhạt màu, ít chi tiết, như phủ sương. Phần trên TRONG SUỐT. Mép trái–phải nối liền |

### 3.2 — GẠCH MẶT ĐẤT (mỗi viên 96×96, ghép liền mạch 4 phía)

Đây là phần map hiện tại yếu nhất — cần đủ bộ viền, không chỉ 2 viên.

| Tên file | Cỡ | Nội dung |
|---|---|---|
| `tile_dat_1_x4.png` | 96×96 | Đất nâu thân, lặp liền mạch |
| `tile_dat_2_x4.png` | 96×96 | Đất nâu biến thể (lẫn sỏi/rễ) — rải xen kẽ chống lặp |
| `tile_co_mat_x4.png` | 96×96 | **Mặt cỏ trên cùng**: nửa trên cỏ xanh lưa tưa, nửa dưới đất. Mép cỏ phải **răng cưa không đều** |
| `tile_co_goc_trai_x4.png` | 96×96 | Cỏ phủ mặt trên **và** vòng xuống mép trái |
| `tile_co_goc_phai_x4.png` | 96×96 | Cỏ phủ mặt trên **và** vòng xuống mép phải |
| `tile_vach_trai_x4.png` | 96×96 | Vách đất cắt dựng đứng, mặt quay trái |
| `tile_vach_phai_x4.png` | 96×96 | Vách đất cắt dựng đứng, mặt quay phải |
| `tile_doc_len_x4.png` | 96×96 | Dốc 45° đi lên, cỏ trên mặt dốc |
| `tile_doc_xuong_x4.png` | 96×96 | Dốc 45° đi xuống |
| `tile_datlat_x4.png` | 96×96 | Đường đất lát đá — lối đi trong làng |

### 3.3 — NHÀ TẦNG GẦN (nền trong suốt, chân chạm mép dưới)

⚠️ **Vẽ NGUYÊN CĂN thành 1 ảnh** — đừng tách mái/tường/cửa rời.

| Tên file | Cỡ | Nội dung |
|---|---|---|
| `nha_nho_x4.png` | 576×576 | Nhà nhỏ 1 tầng: mái ngói đỏ-cam dốc, tường gỗ, 1 cửa giấy, 1 cửa sổ lưới |
| `nha_vua_x4.png` | 768×576 | Nhà vừa: rộng hơn, có hiên trước, rèm vải, 2 cửa sổ |
| `nha_lon_x4.png` | 768×768 | Nhà 2 tầng: ban công gỗ, mái ngói lớn, đèn lồng treo trước cửa |
| `nha_quan_x4.png` | 768×576 | Quán ramen: **rèm noren xanh** treo cửa, biển gỗ chữ Nhật, ghế băng |

### 3.4 — VẬT TRANG TRÍ (nền trong suốt, chân chạm mép dưới)

| Tên file | Cỡ | Nội dung |
|---|---|---|
| `tt_cong_torii_x4.png` | 480×480 | Cổng torii gỗ đỏ, có dây shimenawa + chuông |
| `tt_hangrao_x4.png` | 192×192 | Hàng rào gỗ thưa, **ghép ngang liền mạch được** |
| `tt_denlong_x4.png` | 192×384 | Cột đèn lồng đỏ đứng, có chân đế |
| `tt_gieng_x4.png` | 288×288 | Giếng đá, mái che gỗ nhỏ, gàu treo |
| `tt_thunggo_x4.png` | 192×192 | Cụm 2–3 thùng/vại gỗ xếp chồng |
| `tt_bienhieu_x4.png` | 192×288 | Bảng gỗ cắm đất, chữ Nhật, dây thừng buộc |
| `tt_cay_x4.png` | 576×768 | Cây tán rộng kiểu Nhật, thân nâu, tán xanh đậm |
| `tt_bui_co_x4.png` | 192×96 | Bụi cỏ dại + hoa nhỏ |

### 3.5 — TIỀN CẢNH (nền trong suốt) — map hiện chưa có, thêm vào là khác hẳn

Vẽ **sẫm màu gần đen** (gần camera, ngược sáng).

| Tên file | Cỡ | Nội dung |
|---|---|---|
| `tc_canh_la_x4.png` | 768×384 | Cành lá rủ xuống từ **mép trên**, xanh sẫm gần đen |
| `tc_bui_truoc_x4.png` | 384×192 | Bụi cỏ cao sẫm màu, mọc từ **mép dưới** |

---

## 4. TỔNG KẾT

| Nhóm | Số ảnh |
|---|---|
| 3.1 Tầng phông | 3 |
| 3.2 Gạch mặt đất | 10 |
| 3.3 Nhà tầng gần | 4 |
| 3.4 Vật trang trí | 8 |
| 3.5 Tiền cảnh | 2 |
| **TỔNG** | **27** |

**Làm thử 3 ảnh trước** (`nha_nho`, `tile_co_mat`, `tt_cong_torii`) → gửi cắm vào game chụp duyệt → ưng rồi mới làm loạt.

---

## 5. GIAO THẾ NÀO

- Thư mục phẳng chứa file PNG, **đúng tên**, **đúng cỡ** (bội số 96).
- Báo cáo trung thực: vẽ bao nhiêu ảnh · ảnh nào không chắc và vì sao. **Không suy đoán là "đạt"**.

Sau khi nhận: cắm vào game → **chụp màn hình trong game thật** gửi lại duyệt → lệch thì vẽ lại.
