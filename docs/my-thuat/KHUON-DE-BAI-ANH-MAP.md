# Khuôn đề bài ảnh cho một map — điền vào là dùng được

Chép file này, thay `<TÊN LÀNG>` và bảng màu, giao cho người vẽ. Mọi con số kỹ thuật đã tính sẵn — **đừng sửa chúng**, vì mỗi con số ở đây đều đổi bằng một lỗi đã vấp thật.

> Map 1 Làng Lá đi qua **3 đợt ảnh** vì đề bài đợt đầu thiếu và đợt hai sai cỡ. Khuôn này gộp tất cả vào **một đợt**.

---

## Phong cách — chốt một thứ, không mô tả hai kiểu

**Hand-painted (tranh vẽ mịn), KHÔNG pixel art.**

*(Đã vấp: brief cũ cùng một câu vừa ghi `pixel-art tile` vừa ghi `hand-painted tileset style` — người vẽ phải dừng lại hỏi.)*

Cùng bảng màu, cùng độ chi tiết cho **toàn bộ** ảnh trong một map. Ảnh nhân vật, nhà cửa, nền phải nhìn ra là cùng một bàn tay.

---

## Bảng cỡ — đúng từng pixel, đều là bội số của 96

| Nhóm | Tên file | Cỡ (px, mức x4) |
|---|---|---|
| **Nền** | `nen_troi_<làng>_x4.png` | **2976 × 384** |
| | `nen_nui_<làng>_x4.png` | **2976 × 768** |
| | `nen_langxa_<làng>_x4.png` | **2976 × 576** |
| | `nen_gan_<làng>_x4.png` | **2976 × 288** |
| **Tile nền** (12 ảnh) | `tile_co_mat` · `tile_co_sang_da` · `tile_da_sang_co` · `tile_datlat` | 96 × 96 |
| | `tile_dat_1` · `tile_dat_2` · `tile_dat_giua` · `tile_dat_day` | 96 × 96 |
| | `tile_vach_trai` · `tile_vach_phai` · `tile_chan_vach` · `tile_ke_da` | 96 × 96 |
| **Tile bậc** | `tile_bac_thang` · `tile_bac_thang_dinh` | 96 × 96 |
| **Công trình** | `nha_nho` | 576 × 576 |
| | `nha_vua` | 768 × 576 |
| | `nha_lon` | 768 × 768 |
| | `nha_dacsan` (quán/tiệm riêng của làng) | 768 × 576 |
| | `cong_lang` (torii hoặc tương đương) | 480 × 480 |
| **Cây cối** | `cay_lon` | 768 × 768 |
| | `cay_nho` | 384 × 384 |
| | `bui_ram` | 288 × 192 |
| **Vật nhỏ** | `gieng` 288×288 · `den` 192×384 · `thung` 192×192 · `hang_rao` 192×192 · `bui_co` 192×96 | |

> ### Vì sao nền phải 2976 chứ không phải 1536
> Ảnh giao ở mức **x4** nhưng game hiển thị ở **x2** — bề rộng bị chia đôi. Ảnh 1536 chỉ còn 768px trên màn rộng 1454px, nên nó **lặp gần 2 lần trong cùng một khung hình**, người xem đọc ra ngay là hình dán. *(Đã vấp: cả núi lẫn làng xa của map 1 phải vẽ lại vì lỗi này.)*

---

## Bố cục dọc từng ảnh nền — phần hay bị bỏ sót nhất

Bốn lớp nền chồng lên nhau. **Chỗ giáp giữa chúng là nơi sinh lỗi**, nên mỗi ảnh phải tuân đúng tỉ lệ dọc dưới đây.

| Lớp | 0 → 30% chiều cao | 30 → 80% | 80 → 100% |
|---|---|---|---|
| `nen_troi` | trời + mây | trời nhạt dần | **màu nối** thuần |
| `nen_nui` | trời + mây | khối núi đặc | nhoà dần về **màu nối** |
| `nen_langxa` | **mái nhà ngay từ hàng đầu** — không được có trời | làng dày, cây xen kẽ | nhoà dần về **màu nối** |
| `nen_gan` | mái nhà mờ + sương | sương dày | nhoà về **màu nối** |

**Hai luật tuyệt đối:**

1. **`nen_langxa` không được có trời ở trên.** Nhiệm vụ của nó là **che chân núi**. Có trời ở đỉnh là nó nằm thấp hơn chân núi, che không nổi, và giữa núi với nhà hiện ra một khoảng trời rỗng — làng trông như treo lơ lửng. *(Đã vấp: ảnh làng xa map 1 có 36% trên là trời thuần, phải cắt bỏ.)*

2. **Đáy mỗi ảnh nền phải nhoà về đúng màu nối**, hàng cuối cùng là **đúng mã màu đó**, không còn vân. Cắt phẳng là hiện một vạch ngang chạy suốt màn hình.

**Màu nối của map này:** `RGB( ___ , ___ , ___ )` ← *điền màu trời chủ đạo của làng*

---

## Sáu lỗi hay mắc

1. **Viền quanh ảnh.** Không khung, không viền trắng/đen. *(Đã có 7 tile dính viền, trong game lộ thành ô vuông viền trắng giữa mặt đất.)*
2. **Cỡ phải đúng từng pixel** và là bội số 96. Cỡ lẻ thì engine phóng lên số chẵn gần nhất và **cắt mất một phần** — đã mất hẳn cột phải của một cái cổng.
3. **Chân vật thể chạm đúng mép dưới ảnh.** Hở 2px là trong game nó lơ lửng.
4. **Ảnh nền vẽ liền một mạch**, không nhân đôi nửa ảnh cho đủ bề rộng.
5. **Tile phải lặp được với chính nó** — dán 4 bản cạnh nhau không thấy mối nối.
6. **Bốn tile đất phải gần nhau về độ sáng** (chênh dưới 8%), khác nhau ở **vân**. Chênh nhiều là tầng đất thành bàn cờ.

---

## Riêng bốn tile đất — chỗ khó nhất, đọc kỹ

Lòng đất chiếm **hơn một phần tư màn hình**. Nếu bốn tile này trơn thì cả vùng đó là một mảng nâu chết, và **không cứu được bằng code** — đã thử ba cách khác nhau, cả ba đều hỏng.

Mỗi tile phải có **vân đọc ra được ở cỡ thật**, không phải nhiễu hạt mịn:

| Tên | Nội dung |
|---|---|
| `tile_dat_1` | đất, có **mạch đá nằm ngang** rõ và vài viên cuội |
| `tile_dat_2` | đất, có **rễ cây** thò xuống |
| `tile_dat_giua` | đất lẫn **đá cuội to**, tối hơn ~6% |
| `tile_dat_day` | đất chặt, **vân nén ngang**, tối hơn ~12% |

Nhìn một ô 96px ở khoảng cách bình thường phải thấy *"à, có mạch đá"* — chứ không phải *"một mảng nâu"*.

**Bốn tile phải khác nhau rõ.** Sinh cả bốn từ một mẫu là vô ích: trong game chúng nằm cạnh nhau và mắt bắt ngay cái lặp.

---

## Giao xong thì kiểm

```bash
python tools/map-lang-la/kiem-anh-dot-2.py     # sửa đường dẫn sang thư mục map mới
```

Kiểm cỡ, viền bẩn, đáy nền, ảnh có bị ghép đôi không, chân vật có chạm mép không. **Đạt hết mới lắp.**
