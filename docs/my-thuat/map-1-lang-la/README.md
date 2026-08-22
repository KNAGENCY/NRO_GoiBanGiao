# MAP 1 — LÀNG LÁ: đề bài tạo ảnh (dành cho AI tạo ảnh / thợ vẽ)

> **Thư mục này tự chứa đủ.** Đọc từ trên xuống. Mở `anh-tham-chieu/` ra xem **trước khi vẽ tấm đầu tiên**.
> Luật chung toàn dự án: `../dot-3-toan-bo/_KHUON.md`. File này chỉ nói **map 1 cần gì**.
>
> **Ngày:** 21/08/2026 · **Trạng thái:** map 1 đã chạy được trong game nhưng **chưa đạt mỹ thuật** — xem `anh-tham-chieu/HIEN-TRANG-map1-trong-game.png`.

---

## 0. LÀM THEO ĐÚNG 2 BƯỚC — đừng đảo thứ tự

| Bước | Việc | Vì sao |
|---|---|---|
| **BƯỚC 1** | Vẽ **MỘT bức toàn cảnh map** (concept art) → chờ **DUYỆT** | Duyệt bố cục + tông màu một lần, rẻ. Vẽ 20 mảnh rồi mới thấy không hợp nhau = bỏ hết |
| **BƯỚC 2** | Dựa vào bức ĐÃ DUYỆT, **tách ra từng ảnh rời** theo bảng §3 | Bức concept là kim chỉ nam — mọi mảnh phải cùng tông, cùng nét với nó |

**Bước 2 phải gắn bức concept đã duyệt làm ảnh tham chiếu cho mỗi lần tạo**, nếu không các mảnh sẽ lệch tông nhau.

---

## 1. BƯỚC 1 — BỨC TOÀN CẢNH ĐỂ DUYỆT

**Tên file:** `concept-map1-langla.png` · **Cỡ:** 2048×1024 (ngang) · nền đặc, không cần trong suốt.

**Nội dung:** một khu Làng Lá (Konoha) nhìn ngang kiểu game 2D side-scroll — góc nhìn thẳng từ bên hông, KHÔNG phối cảnh 3/4, KHÔNG nhìn từ trên xuống.

Bố cục từ xa tới gần:

1. **Trời** xanh trong, vài cụm mây trắng.
2. **Núi Hokage** — vách đá lớn khắc mặt các đời Hokage, nằm phía sau làng, mờ hơn (xa).
3. **Dãy nhà Làng Lá** — mái ngói đỏ-cam, tường gỗ, cửa giấy kiểu Nhật. **Cao thấp khác nhau, cụm 2-3 căn rồi chừa khoảng trống** — tuyệt đối không xếp đều tăm tắp.
4. **Cổng torii đỏ** + hàng rào gỗ + đèn lồng + giếng nước + thùng gỗ rải rác.
5. **Mặt đất** — đất nâu có cỏ xanh phủ mép trên, **gồ ghề nhẹ** (có gò, có bậc), không phẳng lì.
6. **Tiền cảnh** — vài bụi cỏ / cành lá sẫm màu ở rìa dưới, tạo chiều sâu.

**Phong cách bắt buộc:** tranh vẽ tay 2D game, **có khối, có đổ bóng, chuyển màu mượt**, viền trong sạch. Tông ấm, đậm chất ninja cổ trang Nhật. **KHÔNG** vector phẳng, KHÔNG hoạt hình trẻ em, KHÔNG 3D render.

**Bám tông màu art đã có** — xem `anh-tham-chieu/HIENCO-*.png` (đất nâu, mái ngói cam-đỏ, núi Hokage xám nâu). Bức concept nên trông như **cùng một bộ** với chúng.

> ⛔ **Dừng ở đây, gửi bức này đi duyệt.** Chưa vẽ tiếp §3.

---

## 2. ẢNH THAM CHIẾU — mở ra xem

| File | Học gì |
|---|---|
| `HIEN-TRANG-map1-trong-game.png` | **Map hiện tại trong game** — chính là thứ cần thay. Thấy rõ: nhà bị chôn nửa dưới đất, 3 nhà giống hệt cách đều, đất phẳng lì, núi Hokage lặp thành dải tường |
| `GIAI-THICH-tileset-vien-chuyen-tiep.png` | **Vì sao cần viên viền** — trái là đất-cỏ cắt thẳng đuột (sai), phải là cỏ lưa tưa phủ xuống (đúng) |
| `TOT-cong-torii-da-tach-nen.png` | **Chuẩn vàng cho vật thể rời**: vẽ nguyên vật thể, nền trong suốt sạch, chân chạm mép dưới ảnh |
| `HIENCO-*.png` | Art đã có trong game — **bám đúng tông màu này** |
| `TOT-naruto-logo.png`, `TOT-chitiet-nhanvat.png` | Mức chi tiết + đổ bóng cần đạt |
| ❌ `XAU-boss-zabuza.png`, `XAU-vatpham.png` | Lô từng bị trả lại: phẳng, vài màu, chân que. **Tránh tuyệt đối** |

---

## 3. BƯỚC 2 — TÁCH THÀNH ẢNH RỜI (chỉ làm sau khi §1 được duyệt)

### Luật kỹ thuật — SAI LÀ GAME HIỂN THỊ HỎNG

1. **Cỡ ảnh phải là bội số của 96** (96, 192, 288, 384…). Cỡ lẻ → Unity tự phóng méo, mất một phần ảnh. *(Lỗi này đã xảy ra thật: ảnh 96px bị phóng lên 128 làm mất hẳn cột phải của cổng torii.)*
2. **Nền trong suốt thật** (PNG có kênh alpha) cho mọi vật thể. Ngoại lệ: nhóm NỀN XA (§3.1) nền đặc.
3. **Chân vật thể chạm ĐÚNG mép dưới ảnh** — không hở khoảng trống, không thừa. Vật thể sẽ được đặt đứng trên mặt đất, hở là nó lơ lửng.
4. **Không viền sáng/tối quanh mép** (dấu cắt nền ẩu). Kiểm bằng cách đặt lên nền trắng VÀ nền đen.
5. Vật thể vẽ **nhìn ngang**, mặt hướng ra trước hoặc sang phải.

### 3.1 — NỀN XA (nền đặc, lặp ngang được)

| Tên file | Cỡ | Nội dung |
|---|---|---|
| `nen_troi_x4.png` | 768×384 | Trời xanh + mây. **Mép trái và mép phải phải nối liền nhau** (lặp vòng không thấy đường ghép). Nửa trên sáng — game lấy màu ở đây làm màu trời |
| `nen_nuihokage_x4.png` | 768×384 | Vách núi khắc mặt Hokage. **Nửa dưới là núi, nửa trên TRONG SUỐT** (để lộ trời). Mép trái-phải nối liền được. Màu nhạt/mờ hơn tiền cảnh (xa) |

### 3.2 — GẠCH MẶT ĐẤT (mỗi viên 96×96, ghép liền mạch)

**Đây là phần map hiện tại yếu nhất** — cần đủ bộ viền, không chỉ 2 viên.

| Tên file | Cỡ | Nội dung |
|---|---|---|
| `tile_dat_1_x4.png` | 96×96 | Đất nâu thân, ghép 4 phía liền mạch |
| `tile_dat_2_x4.png` | 96×96 | Đất nâu biến thể (lẫn sỏi/rễ) — rải xen kẽ chống lặp |
| `tile_co_mat_x4.png` | 96×96 | **Mặt cỏ trên cùng**: nửa trên cỏ xanh lưa tưa, nửa dưới đất — viên nằm ở hàng trên cùng của khối đất |
| `tile_co_goc_trai_x4.png` | 96×96 | Góc trái khối đất: cỏ phủ mặt trên **và** vòng xuống mép trái |
| `tile_co_goc_phai_x4.png` | 96×96 | Góc phải: cỏ phủ mặt trên và vòng xuống mép phải |
| `tile_vach_trai_x4.png` | 96×96 | Vách đất bên trái (mặt cắt dựng đứng) |
| `tile_vach_phai_x4.png` | 96×96 | Vách đất bên phải |
| `tile_doc_len_x4.png` | 96×96 | Dốc thoải đi lên (đất, mặt cỏ chéo 45°) |
| `tile_doc_xuong_x4.png` | 96×96 | Dốc thoải đi xuống |
| `tile_datlat_x4.png` | 96×96 | Đường đất lát đá — lối đi trong làng |

### 3.3 — NHÀ NGUYÊN CĂN (nền trong suốt, chân chạm mép dưới)

⚠️ **Vẽ NGUYÊN CĂN thành 1 ảnh** — đừng tách mái/tường/cửa rời. Nhà hiện tại vụn vì ghép từ mảnh.

| Tên file | Cỡ | Nội dung |
|---|---|---|
| `nha_nho_x4.png` | 192×192 | Nhà nhỏ 1 tầng: mái ngói đỏ-cam dốc, tường gỗ nâu, 1 cửa giấy, 1 cửa sổ |
| `nha_vua_x4.png` | 288×192 | Nhà vừa: rộng hơn, mái 2 tầng nhỏ, hiên trước, rèm vải |
| `nha_lon_x4.png` | 288×288 | Nhà lớn 2 tầng: ban công gỗ, mái ngói lớn, đèn lồng treo trước cửa |
| `nha_quan_x4.png` | 288×192 | Quán ăn (kiểu quán ramen): **rèm vải noren** treo cửa, biển hiệu chữ Nhật, ghế băng |

### 3.4 — VẬT TRANG TRÍ (nền trong suốt, chân chạm mép dưới)

| Tên file | Cỡ | Nội dung |
|---|---|---|
| `tt_cong_torii_x4.png` | 192×192 | Cổng torii gỗ đỏ — **to hơn hiện tại**, cao gần bằng nhà nhỏ |
| `tt_hangrao_x4.png` | 96×96 | Hàng rào gỗ thưa, ghép ngang liền mạch được |
| `tt_denlong_x4.png` | 96×192 | Cột đèn lồng đỏ đứng, có chân đế |
| `tt_gieng_x4.png` | 192×192 | Giếng nước đá, mái che gỗ nhỏ, gàu treo |
| `tt_thunggo_x4.png` | 96×96 | Cụm 2-3 thùng/vại gỗ xếp chồng |
| `tt_bienhieu_x4.png` | 96×192 | Bảng gỗ cắm đất, chữ Nhật, dây thừng buộc |
| `tt_cay_x4.png` | 192×288 | Cây xanh tán rộng kiểu Nhật, thân nâu, tán xanh đậm |
| `tt_bui_co_x4.png` | 96×96 | Bụi cỏ dại + hoa nhỏ |

### 3.5 — TIỀN CẢNH (nền trong suốt) — map hiện chưa có, thêm vào là khác hẳn

Vẽ **sẫm màu hơn** phần còn lại (vì gần camera, ngược sáng).

| Tên file | Cỡ | Nội dung |
|---|---|---|
| `tc_canh_la_x4.png` | 384×192 | Cành lá rủ xuống từ mép trên, xanh sẫm gần đen |
| `tc_bui_truoc_x4.png` | 192×96 | Bụi cỏ cao sẫm màu, dùng che rìa dưới màn hình |

---

## 4. TỔNG KẾT SỐ LƯỢNG

| Nhóm | Số ảnh |
|---|---|
| Bước 1 — concept duyệt | **1** |
| 3.1 Nền xa | 2 |
| 3.2 Gạch mặt đất | 10 |
| 3.3 Nhà nguyên căn | 4 |
| 3.4 Vật trang trí | 8 |
| 3.5 Tiền cảnh | 2 |
| **Tổng bước 2** | **26** |

---

## 5. GIAO THẾ NÀO

- Một thư mục phẳng chứa các file PNG, **đúng tên**, **đúng cỡ** (bội số 96).
- Kèm báo cáo trung thực: đã vẽ bao nhiêu ảnh · ảnh nào không chắc và vì sao.
- **Không suy đoán là "đạt"** — không chắc thì ghi "không chắc".

Sau khi nhận, phía dự án sẽ: cắm vào game → **chụp màn hình trong game thật** gửi lại để duyệt → ảnh nào lệch thì vẽ lại.
