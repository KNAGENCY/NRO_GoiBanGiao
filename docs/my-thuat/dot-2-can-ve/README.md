# ĐỢT 2 — 25 ẢNH CẦN VẼ (dành cho AI tạo ảnh)

> **Thư mục này tự chứa đủ.** Đọc từ trên xuống, mở thư mục `anh-tham-chieu/`
> ra xem trước khi vẽ tấm đầu tiên, rồi vẽ theo danh sách §3.
>
> Luật chung đầy đủ ở `../HUONG-DAN-TAO-ANH.md` — file này chỉ nói **đợt 2 cần gì**.

**Số lượng: 25 ảnh.** Thời điểm: 27/07/2026, sau khi máy chủ xong phần chiến đấu.

---

## 1. BỐN ĐIỀU BẮT BUỘC (đọc 60 giây)

1. **VẼ TO, ĐỪNG TỰ THU NHỎ.** Vẽ ở cỡ gốc lớn nhất bạn vẽ đẹp được
   (1024px trở lên), **giữ đúng tỉ lệ khung** ghi ở §3. Đừng cố vẽ ảnh 64×64 —
   cỡ đó AI vẽ ra là bết. Việc thu nhỏ đã có công cụ lo.
2. **Nền trong suốt thật** (PNG có kênh alpha). Cả 25 ảnh đợt này đều cần.
3. **Tranh vẽ có khối, có đổ bóng** — không phải hình khối phẳng vài màu.
   Xem `anh-tham-chieu/TOT-*.png`. Nếu ảnh bạn xuất ra **nhẹ dưới 3 KB** thì sai.
4. **Báo cáo phải nói thật.** Lô trước tự nhận "QA sạch" trong khi 62 ảnh sai cỡ.
   Chỉ ghi "đã kiểm X" khi thật sự chạy công cụ kiểm X.

---

## 2. PHONG CÁCH — mở ảnh tham chiếu ra xem

| File trong `anh-tham-chieu/` | Học điều gì |
|---|---|
| `TOT-quai-soi-rung.png` | **Chuẩn vàng**: lông có lớp, cơ có khối, nền cắt sạch không viền |
| `TOT-npc-ong-tazuna.png` | Vải có nếp gấp, đổ bóng đúng hướng, phụ kiện chi tiết |
| `TOT-boss-mau.png`, `TOT-vatpham-mau.png` | Mức chi tiết cho boss / vật phẩm nhỏ |
| `TOT-naruto-khung.png`, `TOT-naruto-logo.png`, `TOT-naruto-thanh-chakra.png` | **Phong cách giao diện**: gỗ sơn mài tối + viền cam ấm |
| `TOT-chitiet-nhanvat.png`, `TOT-chitiet-hieuung.png` | Tham chiếu **kiểu vẽ** (không phải nội dung) |
| `THAM-CHIEU-cuon-ky-nang.png` | Cuộn giấy kỹ năng đã có — icon chiêu ở §3.2 phải **cùng gu** với nó |
| ❌ `XAU-boss-zabuza.png`, `XAU-npc.png`, `XAU-vatpham.png` | Lô bị trả lại: phẳng, vài màu, chân que. **Tránh tuyệt đối** |
| ❌ `XAU-thanh-gop-san-khong-dung-duoc.png` | **Lỗi phải sửa ở §3.1** — xem giải thích ngay dưới |
| ❌ `XAU-hu-anthan-qua-mo.png` | **Lỗi phải sửa ở §3.3** — mờ tới mức gần như vô hình |

**Bảng màu 5 làng** (dùng cho icon chiêu ở §3.2):

| Làng | Hệ | Màu chủ đạo |
|---|---|---|
| Lá | Hoả | cam–đỏ, ánh lửa |
| Cát | Phong | vàng cát–xanh lục nhạt, xoáy gió |
| Sương Mù | Thuỷ | xanh dương–xanh lam, nước/sương |
| Đá | Thổ | nâu đất–xám đá |
| Mây | Lôi | tím–trắng xanh, tia sét |

---

## 3. DANH SÁCH 25 ẢNH

### 3.1 — Bốn thanh máu / chakra **tách lớp** (4 ảnh, 240×32, tỉ lệ 15:2)

**Đây là lỗi phải sửa, không phải ảnh mới.** Bộ hiện có giao **ảnh gộp sẵn**
(`XAU-thanh-gop-san-khong-dung-duoc.png`) — nhìn thì đẹp, nhưng **game không
dùng được**.

**Vì sao:** game vẽ thanh máu bằng cách lấy lớp **fill** cắt theo phần trăm máu
còn lại, rồi chồng lên lớp **nền**. Với ảnh gộp sẵn thì máu **luôn đầy, không
bao giờ vơi**. Phải là **hai tấm rời**.

| Tên file | Nội dung |
|---|---|
| `ui_thanh_theluc_nen_x4.png` | **Lớp nền thanh máu**: rãnh lõm tối bo tròn hai đầu, viền sáng mảnh, **RUỘT TRỐNG HOÀN TOÀN** (không có màu đỏ nào bên trong) |
| `ui_thanh_theluc_fill_x4.png` | **Lớp đổ đầy máu**: dải cam–đỏ đặc, highlight mềm phía trên, **hai mép trái/phải cắt PHẲNG VUÔNG** |
| `ui_thanh_chakra_nen_x4.png` | Như nền thanh máu, **tái dùng đúng hình học** (cùng rãnh, cùng bo góc) |
| `ui_thanh_chakra_fill_x4.png` | Dải **xanh cyan** đặc, cùng kiểu highlight và cắt phẳng như fill máu |

⚠️ **Hai mép của lớp fill phải cắt phẳng vuông**, không bo tròn — vì game cắt
nó theo chiều ngang; bo tròn thì lúc máu vơi sẽ thấy mép cụt xấu.
⚠️ Lớp nền và lớp fill phải **khớp nhau chính xác** khi chồng lên (cùng 240×32,
cùng vị trí rãnh) — vẽ cùng một lượt để không lệch.

### 3.2 — Hai mươi icon chiêu thức (20 ảnh, 64×64, tỉ lệ 1:1)

**Đây là thứ hoàn toàn chưa có.** Máy chủ vừa làm xong 20 chiêu (4 chiêu × 5
làng) nhưng **chưa có một icon nào** — người chơi mở bảng kỹ năng sẽ thấy ô trống.

**Yêu cầu chung:** icon vuông, nhìn rõ ở cỡ nhỏ, **một hình chính duy nhất**
(đừng nhồi nhiều chi tiết vụn), nền trong suốt, có viền/ánh sáng nhẹ theo màu
hệ của làng. Cùng gu với `THAM-CHIEU-cuon-ky-nang.png`.

**Làng Lá — hệ Hoả (cam–đỏ)**

| Tên file | Chiêu | Gợi ý hình |
|---|---|---|
| `kn_la_cuoclienhoan_x4.png` | Cước Liên Hoàn *(cận chiến)* | Cú đá xoay có vệt chuyển động |
| `kn_la_haohoacau_x4.png` | Hoả Độn: Hào Hoả Cầu *(tầm xa)* | Quả cầu lửa lớn đang bay |
| `kn_la_ylieu_x4.png` | Y Liệu Nhẫn Thuật *(hồi máu)* | Bàn tay toả ánh sáng xanh lục dịu |
| `kn_la_tuchakra_x4.png` | Tụ Chakra *(tự thân)* | Thủ ấn hai tay, luồng khí cam bốc lên |

**Làng Cát — hệ Phong (vàng cát–xanh lục nhạt)**

| Tên file | Chiêu | Gợi ý hình |
|---|---|---|
| `kn_cat_tramphongcuoc_x4.png` | Trảm Phong Cước *(cận chiến)* | Cú đá kèm lưỡi gió sắc |
| `kn_cat_cuongphongtram_x4.png` | Phong Độn: Cuồng Phong Trảm *(tầm xa)* | Xoáy gió hình lưỡi liềm |
| `kn_cat_sathoylieu_x4.png` | Sa Thổ Y Liệu Thuật *(hồi máu)* | Cát vàng xoáy quanh bàn tay sáng |
| `kn_cat_tuchakra_x4.png` | Tụ Chakra Sa Mạc *(tự thân)* | Thủ ấn, hạt cát bay lên quanh người |

**Làng Sương Mù — hệ Thuỷ (xanh dương–lam)**

| Tên file | Chiêu | Gợi ý hình |
|---|---|---|
| `kn_smu_thuyanhtram_x4.png` | Thuỷ Ảnh Trảm *(cận chiến)* | Nhát chém để lại vệt nước |
| `kn_smu_daithuylong_x4.png` | Thuỷ Độn: Đại Thuỷ Long Đạn *(tầm xa)* | Đầu rồng nước đang lao tới |
| `kn_smu_ylieu_x4.png` | Sương Mù Y Liệu Thuật *(hồi máu)* | Hơi sương trắng xanh quanh bàn tay |
| `kn_smu_tuchakra_x4.png` | Tụ Chakra Sương Mù *(tự thân)* | Thủ ấn, sương cuộn quanh |

**Làng Đá — hệ Thổ (nâu đất–xám đá)**

| Tên file | Chiêu | Gợi ý hình |
|---|---|---|
| `kn_da_nhamthachquyen_x4.png` | Nham Thạch Quyền *(cận chiến)* | Nắm đấm bọc đá nứt |
| `kn_da_nhamthachdan_x4.png` | Thổ Độn: Nham Thạch Đạn *(tầm xa)* | Tảng đá bay có vệt bụi |
| `kn_da_nhamylieu_x4.png` | Nham Y Liệu Thuật *(hồi máu)* | Vòng đá tròn phát sáng ấm |
| `kn_da_tuchakra_x4.png` | Tụ Chakra Nham Thạch *(tự thân)* | Thủ ấn, mảnh đá lơ lửng quanh |

**Làng Mây — hệ Lôi (tím–trắng xanh)**

| Tên file | Chiêu | Gợi ý hình |
|---|---|---|
| `kn_may_loianhcuoc_x4.png` | Lôi Ảnh Cước *(cận chiến)* | Cú đá kèm tia điện |
| `kn_may_causet_x4.png` | Lôi Độn: Cầu Sét *(tầm xa)* | Quả cầu sét nổ tia |
| `kn_may_vanylieu_x4.png` | Vân Y Liệu Thuật *(hồi máu)* | Mây sáng ôm lấy bàn tay |
| `kn_may_tuchakra_x4.png` | Tụ Chakra Vân *(tự thân)* | Thủ ấn, tia điện nhỏ chạy quanh |

> **Bốn chiêu của mỗi làng phải nhìn ra được vai trò**: cận chiến (động tác đấm/đá),
> tầm xa (vật thể đang bay), hồi máu (ánh sáng dịu, ấm), tự thân (thủ ấn/tụ khí).
> Nhìn 4 icon cạnh nhau phải phân biệt được ngay, đừng vẽ cả 4 na ná nhau.

### 3.3 — Vẽ lại một ảnh hiệu ứng (1 ảnh, 1200×480 = dải 5 khung 240×480)

| Tên file | Vấn đề |
|---|---|
| `hu_anthan_x4.png` | **Hiện quá mờ, gần như vô hình** (xem `XAU-hu-anthan-qua-mo.png`) — chỉ ~57 điểm ảnh đục trên toàn bộ 1200×480. Hiệu ứng mà người chơi không thấy thì bằng không |

**Yêu cầu:** vẫn là hiệu ứng ẩn thân/tàng hình mờ dần qua 5 khung, nhưng
**khung 1 phải rõ (độ đục ≥ 60%)** rồi mới nhạt dần tới khung 5. Chiều rộng mỗi
khung 240 (chia hết cho 4 — bắt buộc để các mức thu nhỏ không lệch khung).

---

## 4. TỰ KIỂM TRƯỚC KHI GIAO — bắt buộc

```powershell
python ..\KIEM-TRA-ANH.py <thư-mục-ảnh-của-bạn> nguon
```

**Dán nguyên kết quả vào báo cáo.** Công cụ kiểm thật 7 điều: tên có trong bảng ·
thiếu/thừa · kích thước · kênh trong suốt · ảnh phẳng dưới 40 màu · ảnh trống ·
file nhẹ bất thường.

> ⚠️ **Đừng dùng `KIEM-TRA-ANH.ps1`** (bản PowerShell cũ) — chú thích ghi là có
> kiểm kích thước nhưng **mã không hề kiểm**; chính nó đã báo "sạch" cho lô 343
> trong khi 62 ảnh sai cỡ.

Tên 25 ảnh đợt này **đã có sẵn trong `../kich-thuoc-anh.json`** nên công cụ tra
được ngay. Đặt đúng tên, đừng tự đổi.

---

## 5. GIAO THẾ NÀO

- Một thư mục phẳng chứa 25 file PNG, **đúng tên**, **cỡ gốc lớn** (không thu nhỏ).
- Kèm báo cáo: số ảnh giao · lệnh kiểm đã chạy + kết quả nguyên văn · ảnh nào
  không chắc và vì sao.
- **Không suy đoán là "đạt".** Không chắc thì ghi "không chắc".

Sau khi nhận, phía dự án chạy:
```powershell
python ..\TAO-4-MUC-ZOOM.py <thư-mục-bạn-giao> ..\anh-game
```
để sinh đủ 4 mức zoom (x4/x3/x2/x1) mà game cần.

---

## 6. CHƯA CẦN VẼ ĐỢT NÀY (biết để khỏi làm thừa)

- **Trang bị 4 làng còn lại** (Cát, Sương Mù, Đá, Mây). Bộ hiện có 100 vật phẩm
  nhưng gần như **chỉ phủ Làng Lá** (`vp_ao_t1..t8`, `vp_bangtran_t1..`). Bốn
  làng kia sẽ cần bộ riêng — **nhưng chờ máy chủ xong phần vật phẩm (KH7)** đã,
  vì lúc đó mới chốt được mỗi làng bao nhiêu bộ, bao nhiêu cấp. Vẽ bây giờ rất
  dễ phải bỏ đi.
- **Quái / boss / NPC cho 4 làng kia** — cùng lý do.
- **8 ảnh giao diện đã duyệt** (nút, khung, logo, ô vật phẩm) nằm ở
  `my thuat/dot1-nhom1-duyet/` — đã kiểm đúng cỡ, đúng alpha, **không vẽ lại**.
