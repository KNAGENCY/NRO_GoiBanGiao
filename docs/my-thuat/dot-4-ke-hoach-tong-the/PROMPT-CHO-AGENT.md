# Prompt giao việc cho agent tạo ảnh — Đợt 4 (9 hạng mục)

Giữ đúng cách đã dùng thật ở `D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\PROMPT-CHO-AGENT.md`: mỗi hạng mục có **khối GIAI ĐOẠN 1** (đúng 3 ảnh mẫu, dừng chờ duyệt) và **khối GIAI ĐOẠN 2** (làm nốt, chỉ dán sau khi 3 ảnh mẫu được duyệt **trên ảnh chụp game**). Copy nguyên khối trong dấu ``` và dán cho agent tạo ảnh đã được cấp quyền đọc `D:\NRO_GoiBanGiao\docs\my-thuat\`.

**Trước khi dán bất kỳ khối nào** (phía dự án, xem `README.md` §3.0):
- Chốt **D1** (chân dung vẽ theo 09 hay theo brief chủ thể) — ảnh hưởng khối 01, 02, 04, 09. Dưới mỗi khối có ghi *"nếu chốt theo 09 thì bỏ dòng…"*.
- Chốt **C1** (người chơi vẽ lại hay giữ V1; chạy 5 khung) — ảnh hưởng khối 01.
- **Nối tên mới vào cuối** `D:\NRO_GoiBanGiao\docs\my-thuat\kich-thuoc-anh.json` cho 05 (165), 08 (30), 06 (1), 09 (17), 04 (10 mã) — không nối thì `KIEM-TRA-ANH.py` bỏ qua file.
- Thư mục LƯU VÀO đã tạo sẵn (có `.gitkeep`): `D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\<hạng mục>\`.

Thứ tự giao đề nghị (README §3.1): **08 → 01 → 06 → 09 → 05 (Làng Lá) → 04 → 07 → 05 (Khu rừng) + 03 → 02**.

---

## 01 — NHÂN VẬT NGƯỜI CHƠI

> Chỉ dán khi **C1 đã chốt vẽ lại** (nếu chọn V1 sửa dữ liệu thì bỏ khối này, chỉ còn chân dung ở 09). Nếu **D1 chốt theo 09**: ở giai đoạn 2 bỏ 5 dòng `player_*_portrait_chroma.png`.

### GIAI ĐOẠN 1 — 3 ảnh mẫu

```
Bạn là hoạ sĩ tạo ảnh cho game 2D Naruto. Thư mục brief:
D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the

ĐỌC TRƯỚC KHI VẼ (bắt buộc, theo thứ tự):
1. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\_KHUON-DOT-4.md — phần A (hai kiểu file), B.1 (người chơi), C (7 luật 05/09)
2. D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\_KHUON.md — luật vàng, màu 5 làng
3. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\01-nhan-vat.md — mục 3 (ràng buộc), 4 (tham chiếu), 5 (danh sách), 6 (prompt từng ảnh)
4. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\01-nhan-vat.json — manifest: anh[].prompt_day_du trùng từng chữ với mục 6
5. MỞ VÀ NHÌN KỸ:
   - D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\anh-tham-chieu\THUOC-NHAN-VAT.png  (THƯỚC: cỡ đầu, đường chân)
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\player_idle_3f_chroma.png  (nhân vật GỐC — mọi khung mới phải cùng cỡ đầu với khung idle_0 này)
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\player_cast_3f_chroma.png  (mẫu dải 3 khung tỉ lệ đúng)
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\master_hoc_vien_x4_chroma.png  (chi tiết trang phục)
   - D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-chitiet-nhanvat.png  (mức chi tiết cần đạt)
   - PHẢN DIỆN, TRÁNH: D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\player_run_4f_chroma.png (vẽ nhỏ 12 % + khom),
     D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\player_hit_death_getup_3f_chroma.png (khung NẰM),
     D:\NRO_GoiBanGiao\docs\my-thuat\review\khe-dau-than-2026-09-05\so-sanh-ghep.png (hậu quả: hai khuôn mặt trong game)

VIỆC GIAI ĐOẠN 1: chỉ tạo ĐÚNG 3 dải sau, dùng nguyên prompt ở 01-nhan-vat.md mục 6:
  1. player_run_5f_chroma.png               3620x724  (5 ô 724 — chạy, khó nhất: 5 khung cùng cỡ đầu, không khom)
  2. player_hit_death_getup_3f_chroma.png   2172x724  (3 ô 724 — trúng đòn, CẢ 3 KHUNG ĐỨNG: ngả nhẹ · đỡ/gồng · khuỵu gối)
  3. player_cat_idle_3f_chroma.png          2172x724  (3 ô 724 — học viên Làng Cát, thử trang phục làng khác nhưng đúng cỡ đầu Lá)

MỖI LẦN TẠO ẢNH PHẢI GẮN KÈM: THUOC-NHAN-VAT.png + player_idle_3f_chroma.png (+ một ảnh XAU-* làm phản diện nếu công cụ cho thêm).

LUẬT KHÔNG ĐƯỢC VI PHẠM (7 luật 05/09 + 3 cấm chroma):
- Cỡ ảnh ĐÚNG từng pixel: n ô x 724, mỗi ô 724x724 vuông, ô đều nhau. Số trong tên (_5f, _3f) = số ô máy sẽ cắt.
- MỌI khung cùng CỠ ĐẦU với idle_0 của player_idle_3f: cao mặt 96–99 px, tóc→cằm 266–271 px trong ô 724. Thân to/nhỏ hơn là lộ mặt thứ hai / đầu lửng trong game.
- BÀN CHÂN CHẠM MÉP DƯỚI Ô ở mọi khung, không khung bay, không đệm trống dưới chân.
- Đứng thẳng như sprite Dragon Ball gốc; chạy nghiêng ≤ 10°, KHÔNG khom.
- Nền tím RGB(246,4,242) PHẲNG TUYỆT ĐỐI: không gradient, không bóng đổ, không viền, không sàn.
- Một ô một khung, thứ tự khung cố định, hình cách mép trái/phải ô ≥ 40 px, không cắn sang ô kề.
- Đầu 3/4 quay PHẢI ở mọi khung (client lật gương cho hướng trái).
- Dải hit KHÔNG vẽ nằm (client không bao giờ vẽ nằm).
- Không dùng màu tím hồng/magenta trong hình (bị coi là nền và xoá); không glow/khói lan ra nền; không vẽ vũ khí lên tay.
- Không tự thu nhỏ về 128 px, không đổi tên/số khung. Câu này phải có trong mọi prompt: "same head size as reference idle frame; feet touching bottom edge; uniform magenta background RGB(246,4,242)".

LƯU VÀO: D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\01-nhan-vat\
Đặt đúng tên file như trên.

SAU KHI LƯU, CHẠY VÀ DÁN NGUYÊN KẾT QUẢ:
python D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\kiem-dai-nhan-vat.py D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\01-nhan-vat
(kiểm tên _nf, cỡ n×724, viền tím ≥ 85 % lệch ≤ 12, ô rỗng, lề ≥ 40 px, chân hở đáy ≤ 4 px, cao hình 600–700, đảo tím trong hình. Script không đo cỡ đầu — cỡ đầu bạn tự đặt khung mới cạnh idle_0 trên THUOC-NHAN-VAT.png và ghi số đo tóc→cằm từng khung.)

BÁO CÁO khi xong (nói thật, không suy đoán "đạt"): cỡ thực từng file; số đo tóc→cằm từng khung; khung nào chân không chạm đáy;
khung nào bạn KHÔNG CHẮC và vì sao. DỪNG LẠI sau 3 dải này, chờ duyệt. Không tự làm tiếp.
```

### GIAI ĐOẠN 2 — làm nốt (chỉ dán sau khi 3 dải mẫu được duyệt)

```
3 dải mẫu đã được duyệt. Làm nốt theo 01-nhan-vat.md mục 5–6, giữ nguyên MỌI luật giai đoạn 1, lưu cùng thư mục
D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\01-nhan-vat\ , chạy lại kiem-dai-nhan-vat.py sau mỗi nhóm và báo cáo trung thực từng file.

Thứ tự:
  A. P1 còn lại: player_la_portrait_chroma.png 1254x1254 (chân dung 1 ô — BỎ dòng này nếu chân dung đã giao cho hạng mục 09)
  B. P2 Làng Lá: player_melee_4f_chroma.png 2896x724 (vẽ lại dải đánh, tay/vũ khí không vượt vai)
  C. P2 bốn làng, từng làng trọn bộ rồi mới sang làng kế, thứ tự Cát → Sương Mù → Đá → Mây:
     player_<làng>_idle_3f · run_5f · melee_4f · cast_3f · hit_death_getup_3f (+ player_<làng>_portrait nếu chân dung thuộc 01)
     làng = cat | smu | da | may. Trang phục theo subject block của làng đó ở mục 4. Làng Mây: chỉ dùng chàm ngả xanh, CẤM hồng tím.
  D. P3 (chỉ khi được bảo): player_la_nu_idle_3f, player_la_toc2_idle_3f — thân/chân y hệt bản nam, chỉ đổi đầu.
KHÔNG vẽ lại player_idle_3f và player_cast_3f Làng Lá (đã có, giữ nguyên).
```

---

## 02 — NPC

> Nếu **D1 chốt theo 09**: thay ảnh mẫu 3 (`cd_hokage3_x4.png`) bằng `master_nguoi_giao_nv_x4_chroma.png` và bỏ toàn bộ dòng `cd_*` ở giai đoạn 2 (chân dung sẽ vẽ theo khối 09, cùng lúc với master). Mọi NPC ngoài Iruka là **đề xuất** chờ KH9 — vẽ để KH9 không phải chờ vẽ.

### GIAI ĐOẠN 1 — 3 ảnh mẫu

```
Bạn là hoạ sĩ tạo ảnh cho game 2D Naruto. Thư mục brief:
D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the

ĐỌC TRƯỚC KHI VẼ (bắt buộc, theo thứ tự):
1. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\_KHUON-DOT-4.md — phần A, B.3 (NPC), B.10 (chân dung), C (7 luật)
2. D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\_KHUON.md — luật vàng, màu 5 làng
3. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\02-npc.md — mục 3, 4, 5, 6 (đặc biệt bảng "chuẩn hoá cỡ" ở mục 3)
4. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\02-npc.json — manifest (anh[].luu_vao khác nhau giữa master và chân dung)
5. MỞ VÀ NHÌN KỸ:
   - D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\anh-tham-chieu\THUOC-NHAN-VAT.png  (THƯỚC cỡ đầu/chân — bắt buộc)
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\player_idle_3f_chroma.png  (người chơi gốc: NPC người lớn = ĐẦU BẰNG, thân chân dài hơn)
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\master_iruka_x4_chroma.png  (mẫu KIỂU FILE nền tím 1 ô — KHÔNG lấy cỡ, Iruka v1 vẽ sai thước)
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\iruka_portrait_chroma.png  (bố cục chân dung đầu + vai)
   - D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-npc-ong-tazuna.png và TOT-chitiet-nhanvat.png  (mức chi tiết)
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\07-npc\  (ảnh cũ của từng NPC — CHỈ lấy nhận diện tóc/đồ/phụ kiện, KHÔNG lấy tỉ lệ, tư thế, đạo cụ chạm đất)
   - PHẢN DIỆN: D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\XAU-npc.png, XAU-boss-zabuza.png

VIỆC GIAI ĐOẠN 1: chỉ tạo ĐÚNG 3 file sau, dùng nguyên prompt ở 02-npc.md mục 6:
  1. master_hokage3_x4_chroma.png   724x868  (1 ô nền tím — Hokage Đệ Tam đội nón, khó nhất về chiều cao)
  2. master_tho_ren_x4_chroma.png   724x868  (1 ô nền tím — thợ rèn vai rộng, búa nằm trong silhouette, không chạm đất)
  3. cd_hokage3_x4.png              chuẩn 200x200, VẼ 800x800  (chân dung PNG nền trong suốt thật, mặt phải khớp master 1; máy thu về 200 bằng TAO-4-MUC-ZOOM.py — bạn KHÔNG tự thu)

MỖI LẦN TẠO ẢNH PHẢI GẮN KÈM: THUOC-NHAN-VAT.png + player_idle_3f_chroma.png (master) hoặc iruka_portrait_chroma.png (chân dung) + ảnh 343 của NPC đó.

LUẬT KHÔNG ĐƯỢC VI PHẠM:
- Master: canvas ĐÚNG 724 rộng x 868 cao, nền tím RGB(246,4,242) phẳng tuyệt đối (không gradient/bóng/sàn). Tên không có _nf → máy hiểu 1 ô.
- CÙNG CỠ ĐẦU người chơi: cao mặt 96–99 px, tóc→cằm 266–271 px. NPC cao/thấp chỉ được kéo THÂN + CHÂN, không đụng đầu.
  Người lớn cao 2,65–2,8 ô (676–714 px), Kage đội nón 2,95–3,05 ô (752–778 px).
- CHUẨN HOÁ Cỡ bắt buộc: vẽ to tuỳ ý → đo cao mặt → scale TOÀN ảnh để mặt = 96–99 px → dán lên canvas tím 724x868, bàn chân chạm mép dưới,
  căn giữa ngang, đỉnh tóc/nón cách mép trên ≥ 40 px. Tuyệt đối không co riêng đầu hay thân.
- Một khung đứng thẳng, hai chân chạm đất, đầu 3/4 quay PHẢI. Không tựa cột, không ngồi, không lơ lửng. Đạo cụ trong silhouette, KHÔNG chạm đất.
- Không chữ/kanji/logo trên nón, băng trán, áo. Không màu tím hồng trên người (Làng Mây dùng chàm/xanh-tím). Không glow lan nền.
- Chân dung: VẼ 800x800 (cỡ chuẩn x4 là 200x200 — phía dự án thu bằng TAO-4-MUC-ZOOM.py, bạn không tự thu), nền TRONG SUỐT thật (không tím), đầu + vai, 3/4 quay PHẢI, cùng mặt/màu với master, đọc được khi thu về 200x200.
- Không vẽ lại Iruka.

LƯU VÀO (hai thư mục con, KHÔNG để chung):
  master  → D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\02-npc\master\
  chân dung → D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\02-npc\chan-dung\
Đặt đúng tên file như trên.

SAU KHI LƯU, CHẠY VÀ DÁN NGUYÊN KẾT QUẢ (chạy từ thư mục D:\NRO_GoiBanGiao, cần Pillow + numpy):
python -m tools.artqa.artqa review --art-root "D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\02-npc\master" --review "D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\02-npc\qa\review-02-npc-master.json" --reviewer "<tên bạn>"
python D:\NRO_GoiBanGiao\docs\my-thuat\KIEM-TRA-ANH.py "D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\02-npc\chan-dung" nguon
(review: cần raw_source_status = PASS, issues = [], hộp nội dung mép dưới = 868, mép trên ≥ 40. KIEM-TRA-ANH: bỏ qua dòng "THIEU so voi manifest", 5 mục còn lại phải trống.)
Đo tay trên qa\alpha\<tên>_alpha.png: đỉnh tóc→cằm (266–271) và cao mặt (96–99) — GHI SỐ vào báo cáo cho từng master.

BÁO CÁO khi xong (nói thật): cỡ thực từng file; số đo tóc→cằm / cao mặt từng master; raw_source_status; chân dung có alpha thật không;
file nào bạn KHÔNG CHẮC và vì sao. DỪNG LẠI sau 3 file, chờ duyệt. Không tự làm tiếp.
```

### GIAI ĐOẠN 2 — làm nốt

```
3 file mẫu đã được duyệt. Làm nốt theo 02-npc.md mục 5–6, giữ nguyên MỌI luật giai đoạn 1 (đặc biệt chuẩn hoá cỡ theo thước),
lưu đúng hai thư mục con master\ và chan-dung\, chạy lại hai lệnh kiểm sau mỗi nhóm, báo cáo trung thực từng file kèm số đo tóc→cằm.

Thứ tự:
  A. P1 Làng Lá còn lại (11 file): master_nguoi_giao_nv · master_thay_hocvien · master_y_ta · master_linh_gac · master_ong_tazuna (724x868)
     + cd_nguoi_giao_nv · cd_thay_hocvien · cd_tho_ren · cd_y_ta · cd_linh_gac · cd_ong_tazuna (VẼ 800x800, chuẩn 200x200)
     (BỎ mọi dòng cd_* nếu chân dung đã giao cho hạng mục 09)
  B. P2 bốn làng, thứ tự Cát → Sương Mù → Đá → Mây, mỗi làng 3 master + 3 chân dung: <làng>_kazekage|mizukage|tsuchikage|raikage · <làng>_thay · <làng>_giao_nv
     (tên đầy đủ ở mục 5; Kage đội nón → cao 2,95–3,05 ô; Tsuchikage lão nhân 2,6 ô nhưng ĐẦU KHÔNG ĐỔI; Làng Mây chàm/xanh-tím, cấm hồng tím)
```

---

## 03 — QUÁI

> Sói Rừng đã có (5 file v1, không vẽ lại). 3 quái Khu rừng và 16 quái 4 làng là **đề xuất** chờ máy chủ chốt `quai.yaml`; map 6 chưa vào được — vẽ trước, lắp sau.

### GIAI ĐOẠN 1 — 3 ảnh mẫu

```
Bạn là hoạ sĩ tạo ảnh cho game 2D Naruto. Thư mục brief:
D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the

ĐỌC TRƯỚC KHI VẼ (bắt buộc, theo thứ tự):
1. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\_KHUON-DOT-4.md — phần A, B.4 (quái), C (7 luật)
2. D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\_KHUON.md — luật vàng, màu 5 làng
3. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\03-quai.md — mục 3, 4, 5, 6
4. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\03-quai.json — manifest (thu_tu_khung, cao_muc_tieu_x4 từng file)
5. MỞ VÀ NHÌN KỸ:
   - D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\anh-tham-chieu\THUOC-NHAN-VAT.png  (thước: ô 724 to bằng gì, cỡ đầu cho quái dạng người)
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\wolf_idle_2f_chroma.png và master_soi_rung_x4_chroma.png  (mẫu dải quái ĐÃ LẮP: nét, chi tiết, chân chạm đáy)
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\wolf_run_4f_chroma.png  (mẫu nét chạy — nhưng PHẢN DIỆN về ô lệch và chân bay 190 px, đừng lặp lại)
   - D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-quai-soi-rung.png  (mức chi tiết lông/khối cần đạt)
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\05-quai\quai_lon_rung_x4.png, quai_nhen_doc_x4.png, quai_khi_rung_x4.png  (silhouette tham chiếu dáng/màu — ảnh tĩnh cũ, không lắp)
   - PHẢN DIỆN: D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\XAU-boss-zabuza.png (phẳng, chân que), XAU-hu-anthan-qua-mo.png (glow lan nền)

VIỆC GIAI ĐOẠN 1: chỉ tạo ĐÚNG 3 file sau, dùng nguyên prompt ở 03-quai.md mục 6:
  1. master_lon_rung_x4_chroma.png      1254x1254  (1 ô — thiết kế Lợn Rừng để duyệt trước khi vẽ dải)
  2. lon_rung_run_4f_chroma.png         2896x724   (4 ô 724 — dải chạy, khó nhất: 4 khung cùng tỉ lệ, chân chạm đáy TỪNG ô)
  3. nhen_doc_hit_death_2f_chroma.png   1448x724   (2 ô 724 — thú thấp bẹp: khung chết vẫn tựa đáy; KHÔNG tím trên nhện, KHÔNG glow)

MỖI LẦN TẠO ẢNH PHẢI GẮN KÈM: THUOC-NHAN-VAT.png + wolf_idle_2f_chroma.png + TOT-quai-soi-rung.png (+ ảnh 343 của con đó; XAU-boss-zabuza.png làm phản diện).

LUẬT KHÔNG ĐƯỢC VI PHẠM:
- Cỡ ĐÚNG từng pixel: dải = n ô x 724, mỗi ô 724x724; master 1254x1254. Số trong tên (_2f, _4f) = số ô máy cắt. Công cụ ra cỡ khác → cắt/ghép lại TRƯỚC khi lưu.
- Nền tím RGB(246,4,242) PHẲNG (không gradient, không bóng, không sàn, không cỏ).
- Một ô một khung, ô đều nhau, thứ tự khung cố định: idle 2 (đứng, thở) · run 4 · pounce 2 (đòn: chém/đập/húc/cắn) · hit_death 2 (trúng đòn, chết).
- 10 khung của MỘT con CÙNG TỈ LỆ với khung idle_0 của nó; chân/bụng chạm mép dưới ở MỌI khung (không khung bay).
- Quái dạng người (lính gác): ĐÚNG cỡ đầu người chơi (mặt 96–99, tóc→cằm 266–271) và chết = khuỵu gối, KHÔNG nằm. Thú bốn chân được gục sát đất nhưng vẫn tựa đáy ô.
- Quay PHẢI mọi khung. Hình cách mép trái/phải ô ≥ 24 px, không cắn sang ô kề.
- KHÔNG màu tím hồng/magenta trên quái (Nhện Độc và cả Làng Mây chuyển sang xanh chàm #3B3F8C / #5A4FBF). Tia sét vẽ nét cel sắc, không glow (alpha mềm > 15 % là FAIL).
- Không vẽ nhân vật người chơi, không vẽ mặt đất, không chữ.

LƯU VÀO: D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\
Đặt đúng tên file như trên.

SAU KHI LƯU, KIỂM VÀ DÁN NGUYÊN KẾT QUẢ:
(a) lưu đoạn Python "kiem-o-quai.py" ở 03-quai.md mục 7.2 vào thư mục ảnh rồi chạy  python kiem-o-quai.py  — in cỡ, ô, màu 4 góc, đáy hộp từng ô (chuẩn ≥ 716/724)
(b) chạy từ D:\NRO_GoiBanGiao (cần Pillow + numpy):
python -m tools.artqa.artqa review --art-root D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai --review D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\qa\review-03-quai.json --reviewer "<tên bạn>"
    ĐẠT khi raw_source_status = PASS, issues = [], content_boxes mỗi ô: đáy ≥ 716, mép trái ≥ 24, mép phải ≤ 700, cao hộp các khung lệch ≤ 10 % so idle_0.

BÁO CÁO khi xong (nói thật, không suy đoán "đạt"): cỡ thực từng file; màu nền đo ở 4 góc; khung nào chân không chạm đáy;
khung nào bạn KHÔNG CHẮC và vì sao. DỪNG LẠI sau 3 file này, chờ duyệt. Không tự làm tiếp.
```

### GIAI ĐOẠN 2 — làm nốt

```
3 file mẫu đã được duyệt. Làm nốt theo 03-quai.md mục 5–6, giữ nguyên MỌI luật giai đoạn 1, lưu cùng thư mục
D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\03-quai\ , chạy lại kiem-o-quai.py + artqa review sau mỗi con, báo cáo trung thực từng file.

Thứ tự (từng con trọn bộ master → idle → run → pounce → hit_death rồi mới sang con kế):
  A. P1 Khu rừng còn lại (12 file): lon_rung idle/pounce/hit_death · nhen_doc master/idle/run/pounce · khi_rung master/idle/run/pounce/hit_death
  B. P2 bốn làng, thứ tự Cát → Sương Mù → Đá → Mây, mỗi làng 4 con x 5 file (mã con ở mục 5: cat_bocap, cat_ran_luc, cat_linh_gac, cat_golem_cat;
     smu_cua, smu_thuyxa, smu_linh_gac, smu_camap; da_nhen_da, da_de_nui, da_linh_gac, da_golem; may_soi_nui, may_ho_loi, may_linh_gac, may_golem_set).
     Lính gác = cỡ đầu người chơi, chết khuỵu gối. Golem 2 ô. Làng Mây: chàm/xanh-tím, cấm hồng tím, sét nét cel không glow.
KHÔNG vẽ lại 5 file Sói Rừng (wolf_*, master_soi_rung — đã có).
```

---

## 04 — BOSS

> 3 ảnh mẫu ở đây là **thí nghiệm** (boss cao 1,15× trên ô 880 chưa ai chạy thật) — phía dự án phải chồng khung idle_0 lên `player_idle_3f` đo cỡ đầu và chạy `ghep-khung.py --kiem` trước khi cho vẽ loạt. Nếu **D1 chốt theo 09** hoặc **C7** xác nhận client không vẽ chân dung boss: bỏ 3 dòng `*_portrait_chroma.png` ở giai đoạn 2. Boss 4 làng là **đề xuất** (`boss.yaml` chưa có).

### GIAI ĐOẠN 1 — 3 ảnh mẫu

```
Bạn là hoạ sĩ tạo ảnh cho game 2D Naruto. Thư mục brief:
D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the

ĐỌC TRƯỚC KHI VẼ (bắt buộc, theo thứ tự):
1. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\_KHUON-DOT-4.md — phần A, B.1 (người chơi), B.2 (boss), C (7 luật)
2. D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\_KHUON.md — luật vàng, màu 5 làng
3. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\04-boss.md — mục 3 (đặc biệt 3.4 "vùng cắt mảnh"), 4, 5, 6
4. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\04-boss.json — manifest
5. MỞ VÀ NHÌN KỸ:
   - D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\anh-tham-chieu\THUOC-NHAN-VAT.png  (THƯỚC — bắt buộc)
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\player_idle_3f_chroma.png  (nhân vật gốc: boss "đổi đầu" phải chép Y HỆT thân/chân này; mọi boss cùng cỡ đầu)
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\master_iruka_x4_chroma.png  (tỉ lệ "người lớn = đầu bằng, thân chân dài hơn")
   - D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-boss-mau.png và TOT-chitiet-nhanvat.png  (mức chi tiết)
   - PHẢN DIỆN: D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\XAU-boss-zabuza.png, XAU-npc.png

VIỆC GIAI ĐOẠN 1: chỉ tạo ĐÚNG 3 dải sau, dùng nguyên prompt ở 04-boss.md mục 6:
  1. boss_la_mat_na_vo_idle_3f_chroma.png    2640x880  (3 ô 880 — Mặt Nạ Vỡ pha 2, người lớn 1,15x, thử ô 880 + cùng cỡ đầu)
  2. boss_la_ke_la_mat_idle_3f_chroma.png    2172x724  (3 ô 724 — Kẻ Lạ Mặt pha 1: thân/chân Y HỆT player_idle_3f, chỉ đổi đầu mặt nạ)
  3. boss_la_mat_na_vo_melee_4f_chroma.png   3520x880  (4 ô 880 — dải đánh, thử luật "tay/vũ khí dưới vai")

MỖI LẦN TẠO ẢNH PHẢI GẮN KÈM: THUOC-NHAN-VAT.png + player_idle_3f_chroma.png (+ với dải 3: khung idle_0 của chính boss đó đã duyệt).

LUẬT KHÔNG ĐƯỢC VI PHẠM:
- Cỡ ĐÚNG từng pixel: boss đổi đầu ô 724 (3 khung 2172x724); boss người lớn ô 880 (3 khung 2640x880, 4 khung 3520x880). Số trong tên = số ô.
- CÙNG CỠ ĐẦU ở mọi boss, mọi khung: cao mặt 96–99 px, tóc→cằm 266–271 px (bằng người chơi). Boss cao hơn (≤ 1,20x) là do THÂN + CHÂN dài hơn, không phải đầu to hơn.
- Boss đổi đầu: thân/chân phải CHỒNG KHÍT lên player_idle_3f — khác nhau CHỈ ở đầu.
- Bàn chân chạm mép dưới ô ở MỌI khung; không bụi, không bóng, không sàn dưới chân.
- Đứng thẳng; chạy nghiêng ≤ 10°, không khom. Dải hit: 3 khung ĐỨNG (ngả nhẹ · đỡ · khuỵu gối), không nằm.
- Ở mọi khung KHÔNG PHẢI idle, phần vẽ phía trên cằm bị máy bỏ → tay, vũ khí, khí chakra KHÔNG vượt lên trên vai (đại đao/kiếm cầm ngang hông, chém về trước).
- Ở khung idle, vùng trên cằm CHỈ có đầu (tóc, mặt nạ, băng trán): không kiếm vác lưng nhô lên, không cổ áo cao, không khăn choàng phủ vai.
- Không áo choàng dài quá hông. Không màu tím hồng/magenta (Làng Mây dùng chàm/tím-xanh). Không glow/khói mềm.
- Nền tím RGB(246,4,242) phẳng tuyệt đối. Một ô một khung, quay PHẢI, cách mép ô ≥ 24 px. Không đổi tên/số khung.

LƯU VÀO: D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\04-boss\
Đặt đúng tên file như trên.

SAU KHI LƯU, KIỂM VÀ DÁN NGUYÊN KẾT QUẢ:
(a) kiểm nhanh cỡ + ô + màu góc (cần Pillow):
python -c "import sys,glob;from PIL import Image;import re;[print(f, Image.open(f).size, 'o=',Image.open(f).size[0]//int((re.search(r'_(\d+)(?:f|_chroma)',f) or [None,'1'])[1]), 'goc=',Image.open(f).convert('RGB').getpixel((0,0))) for f in glob.glob(r'D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\04-boss\*.png')]"
    mọi dòng phải đúng cỡ ở mục 5, o= bằng cỡ ô, goc= (246, 4, 242)
(b) chạy từ D:\NRO_GoiBanGiao (cần Pillow + numpy):
python -m tools.artqa.artqa review --art-root D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\04-boss --review D:\NRO_GoiBanGiao\docs\my-thuat\review\dot-4-boss.json --reviewer "<tên bạn>"
    ĐẠT khi raw_source_status = PASS, issues = [], content_boxes: đáy = đáy ô, không chạm mép trái/phải, cao hộp các ô cùng dải chênh ≤ 3 %.
(c) đo tay trên qa\alpha: tóc→cằm mỗi khung 266–271, cao mặt 96–99 — ghi số từng khung.

BÁO CÁO khi xong (nói thật): cỡ thực từng file, số ô; số đo tóc→cằm từng khung; khung nào tay/vũ khí cao quá vai;
khung nào bạn KHÔNG CHẮC và vì sao. DỪNG LẠI sau 3 dải, chờ duyệt (phía dự án sẽ chồng khung lên player_idle_3f để đo). Không tự làm tiếp.
```

### GIAI ĐOẠN 2 — làm nốt

```
3 dải mẫu đã được duyệt (đã đo cỡ đầu và ghép thử 33 khung). Làm nốt theo 04-boss.md mục 5–6, giữ nguyên MỌI luật giai đoạn 1,
lưu cùng thư mục D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\04-boss\ , chạy lại (a)(b)(c) sau mỗi boss, báo cáo trung thực từng file.

Thứ tự:
  A. P1 Làng Lá còn lại: boss_la_mat_na_vo_run_4f (3520x880) · cast_3f (2640x880) · hit_death_getup_3f (2640x880, 3 khung ĐỨNG)
     · boss_smu_hoc_vien_dao_tau_idle_3f (2172x724 — đổi đầu, băng trán Sương Mù 4 vạch sóng, thân/chân y hệt player_idle_3f)
     · 3 chân dung boss_la_ke_la_mat / boss_la_mat_na_vo / boss_smu_hoc_vien_dao_tau _portrait_chroma.png 1256x1256
       (BỎ 3 dòng chân dung nếu đã giao cho hạng mục 09 hoặc client không vẽ chân dung boss)
  B. P2 bốn làng, mỗi boss trọn bộ 5 dải + chân dung rồi mới sang boss kế, thứ tự: boss_cat_kiemphong (Kiếm Khách Gió, 1,15x)
     → boss_smu_thatdao (Kiếm Sĩ Thất Đao, 1,20x — đại đao cầm ngang hông) → boss_da_golem_lon (Golem Đá hình người, 1,20x — khối đầu đá cao 266–271 px từ đỉnh tới hàm)
     → boss_may_songdao (Kiếm Sĩ Song Đao Lôi, 1,15x — chàm/tím-xanh, cấm hồng tím). Tất cả ô 880.
```

---

## 05 — MAP (Làng Lá còn thiếu · Khu rừng · 4 làng)

> Phía dự án phải **nối 165 tên** (khoá `ten_moi_can_them_kich_thuoc` trong `05-map.json`) vào cuối `kich-thuoc-anh.json` trước. Khu rừng chưa vào được (map 6 chưa khai trong máy chủ) — vẽ trước, lắp sau. P2 bốn làng **chỉ bắt đầu sau khi concept của làng đó được duyệt và đo** (prompt concept ở 05-map.md mục 6.0).

### GIAI ĐOẠN 1 — 3 ảnh mẫu (P1 Làng Lá + Khu rừng)

```
Bạn là hoạ sĩ tạo ảnh MÔI TRƯỜNG cho game 2D Naruto. Thư mục brief:
D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the

ĐỌC TRƯỚC KHI VẼ (bắt buộc, theo thứ tự):
1. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\_KHUON-DOT-4.md — phần A, B.5 (map)
2. D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\_KHUON.md — luật vàng (với MAP: hand-painted, KHÔNG pixel art — xem 05-map.md mục 4)
3. D:\NRO_GoiBanGiao\docs\my-thuat\KHUON-DE-BAI-ANH-MAP.md — khuôn 27–33 ảnh một map, bảng cỡ, 6 lỗi phải tránh
4. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\05-map.md — mục 3, 4, 5, 6
5. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\05-map.json — manifest (khoá nen: "dac" | "alpha", dung_tham_chieu từng ảnh)
6. MỞ VÀ NHÌN KỸ:
   - D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\concept-map1-langla.png  (PHONG CÁCH đã duyệt — bám cho mọi map)
   - D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-tham-chieu\THUOC-TI-LE.png  (vật nào cao bao nhiêu ô; nhân vật ≈ 246 px ≈ 2,5 ô)
   - D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-tham-chieu\HIEN-TRANG-map1-trong-game.png  (map 1 đang chạy)
   - D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-da-tao\tile_dat_lienmach_1_x4.png  (MÀU + ĐỘ SÁNG đất gốc — 4 tile vân phải cùng độ sáng)
   - D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-da-tao\nen_nuihokage_rong_x4.png  (nội dung núi Hokage đã duyệt — chỉ sai cỡ 1536)
   - D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-tham-chieu\TOT-cong-torii-da-tach-nen.png  (chuẩn vàng vật thể rời)
   - D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-tham-chieu\GIAI-THICH-tileset-vien-chuyen-tiep.png  (vì sao ô chuyển tiếp phải lởm chởm)
   - PHẢN DIỆN: D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\XAU-vatpham.png, XAU-boss-zabuza.png

VIỆC GIAI ĐOẠN 1: chỉ tạo ĐÚNG 3 ảnh sau, dùng nguyên prompt ở 05-map.md mục 6:
  1. nen_nuihokage_rong2_x4.png       2976x768  (nền núi Hokage vẽ LẠI đúng cỡ 2976 — mục 6.1 #1; PNG ĐẶC; hàng cuối đúng màu nối RGB(10,86,154))
  2. tile_dat_van_1_x4.png            96x96     (tile đất có vân — mục 6.1 #3; lặp 4 phía; cùng độ sáng với tile_dat_lienmach_1)
  3. map_rung_cay_co_thu_1_x4.png     768x768   (cây cổ thụ Khu rừng — mục 6.2 #19; alpha thật; gốc chạm đúng mép dưới)

MỖI LẦN GEN gắn kèm concept-map1-langla.png + các ảnh ghi ở dung_tham_chieu của ảnh đó trong 05-map.json
(vật thể: THUOC-TI-LE.png · tile: tile_dat_lienmach_1 + tile_co_mat · nền: nen_nuihokage_rong 1536).

LUẬT KHÔNG ĐƯỢC VI PHẠM:
- Cỡ ĐÚNG từng pixel như ghi (đều là bội số của 96; nền rộng đúng 2976). Công cụ ra cỡ khác → cắt/resize về đúng cỡ TRƯỚC khi lưu.
- Nền nen_* / map_*_nen_*: PNG ĐẶC (không trong suốt), một bức liền mạch cả chiều rộng — KHÔNG ghép đôi/lật gương hai nửa; mép trái–phải nối liền khi lặp; hàng pixel cuối đúng mã màu nối.
- Tile: lặp liền mạch 4 phía, không viền, không khung, ánh sáng đều; 4 tile đất khác VÂN nhưng độ sáng chênh < 8 %; lối lát chỉ là MẶT (~30 % trên), dưới vẫn là lòng đất.
- Vật thể rời (cây, nhà, vật nhỏ): nền TRONG SUỐT thật, nhìn ngang phẳng (KHÔNG 3/4, KHÔNG từ trên xuống), chân chạm ĐÚNG mép dưới không hở một hàng pixel, không bóng đổ, không mảng đất/cỏ dưới chân, không viền mờ.
- Đúng tỉ lệ THUOC-TI-LE.png: cây cổ thụ 8 ô, nhà 2,5–3 lần nhân vật.
- Không vẽ nhân vật, quái, NPC, chữ vào ảnh map. Hand-painted cel-shading, KHÔNG pixel art.

LƯU VÀO: D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\
Đặt đúng tên file như trên.

SAU KHI LƯU, CHẠY VÀ DÁN NGUYÊN KẾT QUẢ:
python D:\NRO_GoiBanGiao\docs\my-thuat\KIEM-TRA-ANH.py D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\ nguon
(dòng "THIEU KENH TRONG SUOT" cho nen_* là BÌNH THƯỜNG — nền đặc; bỏ qua "THIEU so voi manifest"; các mục [X] khác phải = 0)

BÁO CÁO nói thật theo mẫu 05-map.md mục 7.5: cỡ thực từng ảnh, ảnh nào đặc/ảnh nào alpha, kết quả lệnh, ảnh nào KHÔNG CHẮC và vì sao.
DỪNG LẠI sau 3 ảnh, chờ duyệt. Không tự làm tiếp.
```

### GIAI ĐOẠN 2 — làm nốt

```
3 ảnh mẫu đã được duyệt (đã cắm thử/đặt lên concept và đo bằng do-map.py). Làm nốt theo 05-map.md mục 5–6, giữ nguyên MỌI luật giai đoạn 1,
lưu cùng thư mục D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\05-map\ , chạy lại KIEM-TRA-ANH.py sau mỗi nhóm, báo cáo theo mẫu 7.5.

Thứ tự:
  A. P1 Làng Lá còn lại (5): nen_lang_xa_rong2_x4 2976x576 · tile_dat_van_2/3/4_x4 96x96 · (tile_dat_van_1 đã có)
  B. P1 Khu rừng còn lại (27), theo nhóm: 4 nền map_rung_nen_troi/nui/langxa/gan (2976 x 384/768/576/288, ĐẶC; đáy nhoà về màu nối Khu rừng ghi trong prompt)
     → 14 tile map_rung_tile_* 96x96 → 3 cây còn lại (cay_co_thu_2 768, cay_nho 384, bui_ram 288x192) → 6 vật nhỏ (goc_cay, tang_da, cay_do, cau_go, bien_chi_duong, bui_co)
  C. P2 không theo làng (4): tile_datlat_mat_x4 · tt_co_phu_chan_x4 384x96 (Làng Lá, đề xuất) · map_rung_tc_canh_la 768x384 · map_rung_tc_bui_truoc 384x192 (tiền cảnh)
  D. P2 bốn làng — CHỈ khi được bảo và SAU khi concept làng đó (prompt mục 6.0: concept_lang_cat/suong_mu/da/may.png) được duyệt + đo.
     Thứ tự Cát → Sương Mù → Đá → Mây. Mỗi làng làm 3 mẫu trước: map_<làng>_nha_nho_x4 576x576 · tile mặt đi (map_<làng>_tile_cat_mat / co_mat / da_mat) · map_<làng>_cong_lang_x4 480x480
     → dừng chờ duyệt → rồi nền 4 → tile 14 (Sương Mù thêm 2 tile nước) → nhà/cây → vật nhỏ. Màu nối mỗi làng lấy đúng mã trong prompt của làng đó.
```

---

## 06 — ICON VẬT PHẨM · TRANG BỊ · VŨ KHÍ

> Phía dự án nối `"vp_kunai_luyentap_x4.png": [64, 64]` vào cuối `kich-thuoc-anh.json` trước. Icon 4 làng (P2 8, P3 152) là **đề xuất** chờ KH8; 100 `vp_*` + 20 `vk_*` lô 343 dùng thẳng, **không vẽ lại**.

### GIAI ĐOẠN 1 — 3 ảnh mẫu

```
Bạn là hoạ sĩ tạo icon vật phẩm cho game 2D Naruto. Thư mục brief:
D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the

ĐỌC TRƯỚC KHI VẼ (bắt buộc, theo thứ tự):
1. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\_KHUON-DOT-4.md — phần A, B.6, D.1, E (kiểu tĩnh alpha, cách nộp)
2. D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\_KHUON.md — luật vàng, bảng cỡ, màu 5 làng
3. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\06-vat-pham.md — mục 3 (ràng buộc icon), 4 (ITEM ICON RULES + tham chiếu), 6 (prompt)
4. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\06-vat-pham.json — manifest: lọc uu_tien, chỉ vẽ trang_thai can_ve / ve_lai
5. MỞ XEM: D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-vatpham-mau.png (đạt) và XAU-vatpham.png (tránh);
   5–6 icon trong D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\02-vat-pham\ (phong cách đang có — hand-painted cel-shading);
   D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\04-vu-khi\vk_kunai_thuong_x4.png (hình dáng kunai)

VIỆC GIAI ĐOẠN 1: chỉ tạo ĐÚNG 3 ảnh sau, theo đúng khối prompt ở mục 6:
  1. vp_kunai_luyentap_x4.png   chuẩn 64x64, VẼ 1024x1024  — mục 6.1  (kunai vuông, khó nhất: ép vật dài vào ô vuông, đọc được ở 16 px)
  2. vp_cat_ao_t1_x4.png        chuẩn 64x64, VẼ 1024x1024  — mục 6.7  (áo t1 Làng Cát — thử nhân bố cục từ vp_ao_t1 Làng Lá sang làng khác)
  3. vp_may_quan_t1_x4.png      chuẩn 64x64, VẼ 1024x1024  — mục 6.14 (quần t1 Làng Mây — bảng màu khó nhất: chàm/tím-xanh + trắng)

MỖI LẦN TẠO ẢNH PHẢI GẮN KÈM các file ghi ở dòng THAM CHIẾU của khối prompt đó
(icon Làng Lá CÙNG bộ phận + CÙNG cấp trong 343-anh\02-vat-pham + TOT-vatpham-mau.png + XAU-vatpham.png).

LUẬT KHÔNG ĐƯỢC VI PHẠM:
- Vuông 1024x1024, vật lấp ~80 % khung, canh giữa, nhìn 3/4 từ trên, sáng trên-trái, lơ lửng. Máy sẽ thu về 64/48/32/16 — không tự thu.
- Nền TRONG SUỐT thật (PNG alpha). Nếu công cụ không xuất alpha: nền xám phẳng RGB(128,128,128) — KHÔNG dùng tím chroma.
- Không bóng đổ dưới, không nền cảnh, không tay người, không bàn/kệ, không chữ, không viền sáng quanh mép. Một ảnh một món.
- Icon 4 làng: CÙNG bố cục với icon Làng Lá cùng bộ phận/cấp, chỉ đổi chất liệu/màu hệ làng/biểu tượng băng trán
  (Cát đồng hồ cát · Sương Mù 4 vạch sóng · Đá 2 phiến đá · Mây khối mây) — KHÔNG lá Konoha trên đồ làng khác.
- Đặt ĐÚNG tên file. Không tự đổi tên, không tự thêm tên vào kich-thuoc-anh.json.

LƯU VÀO: D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\06-vat-pham\

SAU KHI LƯU, CHẠY VÀ DÁN NGUYÊN KẾT QUẢ:
python D:\NRO_GoiBanGiao\docs\my-thuat\KIEM-TRA-ANH.py D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\06-vat-pham\ nguon
(bỏ qua mục "THIEU so voi manifest" — đó là các tên khác của bảng 691; các mục [X] còn lại phải trống)

BÁO CÁO khi xong (nói thật): đã tạo mấy ảnh, cỡ thực tế từng ảnh, ảnh nào alpha thật / phải tách nền,
ảnh nào KHÔNG CHẮC đạt và vì sao (ví dụ kunai vẫn mảnh ở 16 px). DỪNG LẠI sau 3 ảnh, chờ duyệt. Không tự làm tiếp.
```

### GIAI ĐOẠN 2 — làm nốt

```
3 ảnh mẫu đã được duyệt. Làm nốt theo 06-vat-pham.md mục 5–6, giữ nguyên MỌI luật giai đoạn 1, lưu cùng thư mục
D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\06-vat-pham\ , chạy lại KIEM-TRA-ANH.py, báo cáo trung thực từng file.

Thứ tự:
  A. P2 còn lại (6): vp_cat_quan_t1 · vp_smu_ao_t1 · vp_smu_quan_t1 · vp_da_ao_t1 · vp_da_quan_t1 · vp_may_ao_t1 (mỗi cái gắn icon Làng Lá cùng bộ phận/cấp)
  B. P3 (154) — CHỈ khi được bảo rõ là KH8 đã chốt danh mục: 152 icon 4 làng x 5 bộ phận x 8 cấp còn lại + vẽ lại vp_giay_t1, vp_gang_t1 Làng Lá.
KHÔNG vẽ lại 4 icon vp_vien_theluc_1 / vp_vien_chakra_1 / vp_ao_t1 / vp_quan_t1 và vk_kunai_thuong (đã có, dùng thẳng).
```

---

## 07 — KỸ NĂNG · HIỆU ỨNG

> 20 icon chiêu đã có (không vẽ). Bộ đóng gói hiệu ứng **chưa có** — vẽ vẫn được, lắp sau; 3 dải mẫu dùng để viết bộ đóng gói song song. Cảnh báo `loop_boundary_motion_jump` trên dải **một lần** là báo giả của chính sách ghi cứng — ghi vào báo cáo, không sửa dải cho "xanh".

### GIAI ĐOẠN 1 — 3 dải mẫu

```
Bạn là hoạ sĩ tạo hiệu ứng (VFX) cho game 2D Naruto. Thư mục brief:
D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the

ĐỌC TRƯỚC KHI VẼ (bắt buộc, theo thứ tự):
1. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\_KHUON-DOT-4.md — phần A (3 điều cấm chroma), B.7, B.8 (hiệu ứng)
2. D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\_KHUON.md — luật vàng, màu 5 làng
3. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\07-ky-nang-hieu-ung.md — mục 3.2 (12 luật dải hiệu ứng), 4, 5 (cột Neo/Lặp), 6
4. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\07-ky-nang-hieu-ung.json — manifest (khoá neo A/B/C, lap true/false, y_nghia_khung từng khung)
5. MỞ VÀ NHÌN KỸ:
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\343-anh\08-hieu-ung\hu_hoadon_x4.png  (MỨC CHI TIẾT VFX cần đạt: lửa nhiều lớp, khói đặc, tàn lửa — ảnh cũ dạng lưới, chỉ học nét)
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\qa\contact-sheets\hu_anthan_10f_light-dark.png  (khói/loé kiểu cel đúng chroma)
   - D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-naruto-logo.png  (bảng màu)
   - D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\anh-tham-chieu\THUOC-NHAN-VAT.png và D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\player_idle_3f_chroma.png  (thước: nhân vật cao 246 px cùng thang)
   - PHẢN DIỆN: D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\XAU-hu-anthan-qua-mo.png (quá mờ, halo, alpha mềm)

VIỆC GIAI ĐOẠN 1: chỉ tạo ĐÚNG 3 dải sau, dùng nguyên prompt ở 07-ky-nang-hieu-ung.md mục 6:
  1. hu_la_haohoacau_no_6f_chroma.png    2304x384  (6 ô 384 — Hào Hoả Cầu NỔ trên mục tiêu, một lần, neo B tâm ô; khó nhất về khói đặc không alpha mềm)
  2. hu_la_ylieu_8f_chroma.png           3072x384  (8 ô 384 — Y Liệu hồi máu, LẶP, neo A giữa-đáy; khó nhất về ánh sáng không glow)
  3. hu_la_haohoacau_bay_3f_chroma.png   576x192   (3 ô 192 — đạn cầu lửa bay, neo C tâm ô, nằm ngang mũi hướng PHẢI; khó nhất về cân tâm)

MỖI LẦN TẠO ẢNH PHẢI GẮN KÈM: hu_hoadon_x4.png + TOT-naruto-logo.png + XAU-hu-anthan-qua-mo.png (phản diện); với dải neo A thêm player_idle_3f_chroma.png để đúng thước.

LUẬT KHÔNG ĐƯỢC VI PHẠM:
- MỘT HÀNG ngang, một ô một khung, ô đều nhau, trái→phải theo đúng thứ tự khung ở y_nghia_khung. Số trong tên (_6f) = số ô. Ô 384x384 hoặc 192x192, cả dải = n x ô.
- Nền tím RGB(246,4,242) PHẲNG tuyệt đối — không gradient, không vignette, không nhiễu, không bóng, không sàn.
- KHÔNG màu tím/hồng/fuchsia trong hình (bị coi là nền và xoá). Làng Mây dùng chàm/xanh-tím RGB≈(110,90,230) + trắng-cyan.
- KHÔNG glow/khói mờ lan ra nền, KHÔNG alpha mềm: vẽ lửa/nước/ánh sáng bằng dải màu ĐẶC kiểu cel (lõi sáng → trung → viền tối). Khung cuối dải một lần NHỎ đi, không mờ đi.
- Hình không chạm mép trái/phải ô. Tâm/neo CỐ ĐỊNH ở mọi khung theo cột Neo (A giữa-đáy · B tâm · C tâm, đạn nằm ngang mũi phải).
- Mọi chuyển động hướng PHẢI. KHÔNG vẽ nhân vật, tay, chân, bóng người, mặt đất, đường chân trời.
- Không hai khung liền giống nhau. Dải LẶP: khung cuối nối về khung 1 (gần giống, không trùng). Dải MỘT LẦN: khung 1 yếu → giữa đỉnh → cuối tàn dư nhỏ.
- Chỉ vẽ khung chính (3–8); tốc độ do dữ liệu quyết định, không vẽ thêm khung. Không đổi tên/số khung sau khi nộp.

LƯU VÀO: D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\07-ky-nang-hieu-ung\
Đặt đúng tên file như trên.

SAU KHI LƯU, KIỂM VÀ DÁN NGUYÊN KẾT QUẢ:
(a) kiểm cỡ từ header PNG (không cần cài gì):
python -c 'import struct,glob,os,re;[print(os.path.basename(p), struct.unpack(">II", open(p,"rb").read(24)[16:24]), "n="+re.search(r"_(\d+)f_", p).group(1)) for p in glob.glob(r"D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\07-ky-nang-hieu-ung\*.png")]'
    đối chiếu từng dòng với cột Cỡ ở mục 5: rộng = n x ô, cao = ô.
(b) chạy từ D:\NRO_GoiBanGiao (cần Pillow + numpy):
python -m tools.artqa.artqa review --art-root "D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\07-ky-nang-hieu-ung" --review "D:\NRO_GoiBanGiao\docs\my-thuat\review\dot-4-07-ky-nang-hieu-ung.json" --reviewer "<tên bạn>"
    ĐẠT khi raw_source_status = PASS, không empty_animation_cell, không duplicate_adjacent_frame, soft alpha ≤ 15 %; mở qa\previews\*_animation.gif xem chuyển động,
    qa\crop-reports\*_crop.json xem tâm hộp không đổi giữa các ô (neo B/C) hoặc đáy hộp = đáy ô (neo A).
    Với dải MỘT LẦN, cảnh báo loop_boundary_motion_jump là báo giả — ghi rõ trong báo cáo, KHÔNG sửa dải để cho xanh.

BÁO CÁO khi xong (nói thật): cỡ thực từng file; dải nào không chắc và vì sao; dán nguyên summary + issues từng dải của artqa review.
DỪNG LẠI sau 3 dải, chờ duyệt. Không tự làm tiếp.
```

### GIAI ĐOẠN 2 — làm nốt

```
3 dải mẫu đã được duyệt. Làm nốt theo 07-ky-nang-hieu-ung.md mục 5–6, giữ nguyên MỌI luật giai đoạn 1, lưu cùng thư mục
D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\07-ky-nang-hieu-ung\ , chạy lại (a)(b) sau mỗi nhóm, báo cáo trung thực từng dải.

Thứ tự:
  A. P1 còn lại (5): hu_la_cuoclienhoan_6f 2304x384 (neo B, một lần) · hu_la_tuchakra_8f 3072x384 (neo A, lặp, hai nhịp đập)
     · hu_trungdon_5f 960x192 (neo B, một lần) · hu_chet_vongsang_8f 3072x384 (neo A, lặp) · hu_hoisinh_6f 2304x384 (neo A, một lần)
  B. P2 bốn làng, thứ tự Cát → Sương Mù → Đá → Mây, mỗi làng 5 dải: <cận chiến>_6f · <tầm xa>_bay_3f (192) · <tầm xa>_no_6f · <hồi máu>_8f · <làng>_tuchakra_8f
     (tên đầy đủ ở mục 5: hu_cat_tramphongcuoc, hu_cat_cuongphongtram_bay/no, hu_cat_sathoylieu … hu_may_loianhcuoc, hu_may_causet_bay/no, hu_may_vanylieu);
     hồi máu mọi làng giữ xanh lục y liệu làm chính, tụ chakra giữ xanh dương chakra, điểm nhấn màu làng; Làng Mây chàm, cấm hồng tím.
     + hu_ketan_6f 2304x384 (kết ấn chung, neo B).
  C. P3 (4) — chỉ khi được bảo: hu_lencap_8f · hu_boss_xuathien_8f · hu_napchakra_thuoc_6f · hu_hoimau_thuoc_6f (hai cái sau ô 192).
KHÔNG vẽ icon chiêu (20 icon đã có), KHÔNG vẽ lại hu_anthan_10f.
```

---

## 08 — GIAO DIỆN (UI / HUD)

> Phía dự án phải **nối 30 tên** (cuối §5 của `08-giao-dien.md`) vào `kich-thuoc-anh.json` trước. Đây là hạng mục rẻ và an toàn nhất: ảnh thay **thẳng file client** theo tên cũ, không cần code — nhưng vì thế **cỡ phải đúng số**, không "vẽ to rồi thu".

### GIAI ĐOẠN 1 — 3 ảnh mẫu

```
Bạn là hoạ sĩ vẽ giao diện (UI) cho game 2D Naruto. Thư mục brief:
D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the

ĐỌC TRƯỚC KHI VẼ (bắt buộc, theo thứ tự):
1. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\_KHUON-DOT-4.md — phần A, B.9 (UI/HUD), D.1
2. D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\_KHUON.md — luật vàng
3. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\08-giao-dien.md — mục 3 (hình học phải giữ), 4 (chất liệu bộ UI), 5 (cột "Thay cho" file client), 6
4. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\08-giao-dien.json — manifest (khoá anh_cu_x2_chi_doc = file client cũ để giữ hình học, cat_manh)
5. MỞ VÀ NHÌN KỸ:
   - D:\NRO_GoiBanGiao\my thuat\dot1-nhom1-duyet\dot1-nhom1-duyet\01-giao-dien\ui_khung_hoithoai_x4.png  (chất liệu gỗ sơn mài + dây cam + giấy — ẢNH GỐC đồng bộ cả bộ)
   - D:\NRO_GoiBanGiao\my thuat\dot1-nhom1-duyet\dot1-nhom1-duyet\01-giao-dien\ui_nut_thuong_x4.png  (màu/độ vát nút)
   - D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\anh-tham-chieu\HIEN-TRANG-hud-cu-x2.png  (34 file client cũ ghép một tấm có nhãn — GIỮ ĐÚNG HÌNH HỌC từng file)
   - D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\anh-tham-chieu\HIEN-TRANG-map1-trong-game.png  (HUD cũ nằm đâu trên màn hình)
   - D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-naruto-khung.png, TOT-naruto-thanh-chakra.png (chỉ để hiểu rãnh), TOT-naruto-logo.png
   - PHẢN DIỆN: D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\XAU-thanh-gop-san-khong-dung-duoc.png (đừng gộp phần game tự vẽ vào ảnh)
   - CHỈ ĐỌC, KHÔNG SỬA: file client cũ tương ứng trong D:\NRO_Project\Assets\Resources\res\x2\mainimage\ (tên ở cột "Thay cho" mục 5)

VIỆC GIAI ĐOẠN 1: chỉ tạo ĐÚNG 3 ảnh sau, theo khối prompt mục 6 (stt 1, 4, 5):
  1. ui_hud_khung_trai_x4.png       620x164  (khung HUD góc trên-trái: bảng tên + ô chân dung 26x26 x1 + hai RÃNH máu/chakra là vùng ĐẶC tối, KHÔNG vẽ thanh)
  2. ui_khung_popup_9o_x4.png       40x320   (dải 8 ô 40x40 xếp dọc đúng thứ tự cho hộp thoại 9-slice; ruột hộp do code đổ màu kem #FFE5C7)
  3. ui_nut_thuong_3manh_x4.png     120x96   (nút chữ 3 mảnh trái/giữa/phải 40x96, mảnh giữa lặp liền mạch; mặt nút sáng vừa để chữ tối đọc được)

MỖI LẦN TẠO ẢNH PHẢI GẮN KÈM: ui_khung_hoithoai_x4.png + ui_nut_thuong_x4.png (chất liệu) + HIEN-TRANG-hud-cu-x2.png (hình học) — và các ảnh ghi ở dung_tham_chieu của ảnh đó.

LUẬT KHÔNG ĐƯỢC VI PHẠM:
- Cỡ = ĐÚNG SỐ như ghi (không vẽ to rồi thu, không lề trống thừa — lề làm lệch neo). Mọi cỡ chia hết cho 4.
- PNG alpha thật: ngoài hình trong suốt hoàn toàn, mép cắt cứng, không glow/bóng lan ra ngoài, không nền cảnh, không chữ, không watermark, không viền quanh canvas.
- HÌNH HỌC PHẢI GIỮ vì client vẽ đè lên: cửa sổ ĐẶC (lòng ô chiêu, ô chân dung, hai rãnh máu/chakra = mảng tối không trong suốt);
  cửa sổ TRONG (nút thuốc chừa vùng x4 36..124, 40..120); vùng CHỮ TRỐNG (client in chữ lên nút/ô menu/bảng tên);
  mảnh ghép đúng THỨ TỰ (8 ô dọc, 3 mảnh ngang); mặt trong viền kết thúc ĐÚNG mã màu ruột do code đổ (kem #FFE5C7 / xanh #B7E8B0 / giấy #F0E6C3 / trắng #FFFFFF).
- Chất liệu chung 5 làng: gỗ sơn mài tối #3A2418–#5A3A22 · dây/viền cam ấm #E8722A–#F2A445 · giấy da #F1E2C0 · viền tối #1B1410 · một dấu lá xoáy nhỏ. Có vát/bóng/vân gỗ thật, không phẳng 3 màu.
- KHÔNG vẽ thanh máu/chakra vào khung; KHÔNG chữ mẫu lên nút; KHÔNG đồ Dragon Ball (ngọc rồng, đậu thần, tiền cũ); KHÔNG đổi thứ tự 8 ô / 3 mảnh; KHÔNG bo góc mảnh lặp; KHÔNG tím hồng làm màu chủ.

LƯU VÀO: D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\08-giao-dien\
Đặt đúng tên file như trên.

SAU KHI LƯU, CHẠY VÀ DÁN NGUYÊN KẾT QUẢ:
python D:\NRO_GoiBanGiao\docs\my-thuat\KIEM-TRA-ANH.py D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\08-giao-dien nguon
(nút trơn có thể bị báo "ít màu" — được phép nếu có vát/bóng thật, ghi rõ trong báo cáo thay vì tô nhiễu; bỏ qua "THIEU so voi manifest")
Tự soi: đặt ảnh lên nền đen và nền trắng — vùng đặc/trong đúng chỗ chưa, có quầng sáng ở mép không.

BÁO CÁO khi xong (nói thật): cỡ thực từng ảnh; vùng đặc/trong suốt nằm đúng chỗ chưa; ảnh nào KHÔNG CHẮC và vì sao.
DỪNG LẠI sau 3 ảnh, chờ duyệt (phía dự án sẽ chép vào client một máy và chụp màn hình). Không tự làm 27 ảnh còn lại.
```

### GIAI ĐOẠN 2 — làm nốt

```
3 ảnh mẫu đã được duyệt trên ảnh chụp game. Làm nốt theo 08-giao-dien.md mục 5–6, giữ nguyên MỌI luật giai đoạn 1 (đặc biệt: ĐÚNG SỐ cỡ, hình học theo file cũ),
lưu cùng thư mục D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\08-giao-dien\ , chạy lại KIEM-TRA-ANH.py, báo cáo trung thực từng ảnh.

Thứ tự:
  A. P1 còn lại (7): ui_hud_o_kynang 112x112 · ui_hud_o_kynang_chon 112x112 · ui_nut_nhan_3manh 120x96 · ui_menu_o 230x230 · ui_menu_o_chon 230x230 · ui_hud_nut_chat 124x124 · ui_hud_nut_chat_nhan 124x124
  B. P2 (15): ui_khung_popup_9o_sang 40x320 · ui_hud_nut_menu 56x120 · ui_hud_muiten_chi_duong 52x64 · 6 góc cuộn giấy ui_khung_goc_1..6 (256x216 · 132x52 · 244x72 · 204x88 · 256x32 · 256x40)
     · ui_khung_thoai_goc 28x24 · ui_nen_nhanvat_tuido 800x824 (phông túi đồ, được vignette mép) · ui_nut_phim_mem 272x100 · ui_nut_phim_mem_nhan 272x100 · ui_hud_nut_thuoc 160x160 · ui_hud_nut_thuoc_nhan 160x160
  C. P3 (5) — chỉ khi được bảo (màn chọn làng 5 ô đã có code): ui_lang_la/cat/suong/da/may_lon 192x192 theo màu hệ từng làng.
```

---

## 09 — CHÂN DUNG · THƯƠNG HIỆU

> Nếu **D1 chốt theo 09** (đề nghị của README): đây là nơi duy nhất vẽ chân dung — nhưng mỗi chân dung phải vẽ **cùng lúc** với master/dải idle của nhân vật đó (01/02/04) để khớp mặt. Nếu D1 chốt ngược lại: bỏ các dòng chân dung ở đây, chỉ giữ thương hiệu. Phía dự án nối 17 tên thương hiệu + mã chân dung vào `kich-thuoc-anh.json` trước khi chạy `KIEM-TRA-ANH.py` cho nhóm thương hiệu.

### GIAI ĐOẠN 1 — 3 ảnh mẫu

```
Bạn là hoạ sĩ tạo ảnh cho game 2D Naruto. Thư mục brief:
D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the

ĐỌC TRƯỚC KHI VẼ (bắt buộc, theo thứ tự):
1. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\_KHUON-DOT-4.md — phần A, B.10 (chân dung), C
2. D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\_KHUON.md — luật vàng, màu 5 làng
3. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\09-chan-dung-thuong-hieu.md — mục 3, 4, 5, 6
4. D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\09-chan-dung-thuong-hieu.json — manifest (kich_thuoc_dich, ve_o_co từng ảnh)
5. MỞ VÀ NHÌN KỸ:
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\iruka_portrait_chroma.png  (KHUÔN chân dung ĐÃ LẮP: bố cục đầu+vai, quay phải, vai cắt phẳng chạm đáy)
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\qa\contact-sheets\iruka_portrait_light-dark.png  (bản đã tách nền trên nền sáng/tối)
   - D:\NRO_GoiBanGiao\docs\my-thuat\dot-4-ke-hoach-tong-the\anh-tham-chieu\THUOC-NHAN-VAT.png  (tỉ lệ đầu — chân dung phải cùng khuôn mặt này)
   - D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\player_idle_3f_chroma.png và master_hoc_vien_x4_chroma.png  (mặt/tóc/băng trán/áo người chơi Lá)
   - D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\TOT-naruto-logo.png  (EMBLEM lá-lửa + bảng màu thương hiệu — gắn kèm MỌI lần gen)
   - D:\NRO_GoiBanGiao\docs\my-thuat\map-1-lang-la\concept-map1-langla.png  (tông nền cho splash/bìa)
   - D:\NRO_GoiBanGiao\docs\my-thuat\lang-la-client-assets\Logo\1.png  (wordmark 599x136 đang chạy — để biết vùng chữ sẽ to bao nhiêu; KHÔNG vẽ chữ)
   - PHẢN DIỆN: D:\NRO_GoiBanGiao\docs\my-thuat\dot-3-toan-bo\anh-tham-chieu\XAU-npc.png, XAU-boss-zabuza.png; D:\NRO_Project\Assets\Resources\appicon.png (icon Dragon Ball PHẢI THAY — chỉ đọc)

VIỆC GIAI ĐOẠN 1: chỉ tạo ĐÚNG 3 ảnh sau, dùng nguyên prompt ở mục 6 (6.1, 6.36, 6.6):
  1. nguoi_choi_la_portrait_chroma.png       1254x1254  (1 ô nền tím — chân dung người chơi Lá, mặt phải khớp khung idle_0 đang chạy)
  2. ui_icon_app_x4.png                      chuẩn 512x512, VẼ 1024x1024  (PNG ĐỤC, emblem lá-lửa, đọc được ở 48 px, KHÔNG chữ)
  3. nhan_gia_mat_na_p1_portrait_chroma.png  1254x1254  (1 ô nền tím — boss pha 1: mặt nạ sứ trắng trên bộ học viên, phải giống người chơi mà vẫn ra "kẻ lạ")

MỖI LẦN TẠO CHÂN DUNG gắn kèm: iruka_portrait_chroma.png + ảnh master/dải của CHÍNH nhân vật đó + THUOC-NHAN-VAT.png (+ XAU-npc.png).
MỖI LẦN TẠO THƯƠNG HIỆU gắn kèm: TOT-naruto-logo.png + concept-map1-langla.png + Logo\1.png.

LUẬT KHÔNG ĐƯỢC VI PHẠM:
- Chân dung: ĐÚNG 1 ô 1254x1254, nền tím RGB(246,4,242) phẳng tuyệt đối; đầu + vai, 3/4 quay PHẢI, đầu thẳng; vai cắt PHẲNG chạm đúng mép dưới;
  cùng mặt/tóc/băng trán/màu áo với sprite của nhân vật đó; không chữ, không khung, không glow, không bóng; KHÔNG tím hồng trong hình (Làng Mây: tím ngả xanh, R−G < 100).
  Hình phải đậm, đọc được khi máy thu về 200x200 (khung thoại) và 104x104 (ô HUD).
- Thương hiệu: PNG RGBA ĐỤC (alpha toàn 255), đúng cỡ ghi ở dòng OUTPUT của prompt; KHÔNG chữ/lettering (wordmark dán sau); nội dung quan trọng trong vùng an toàn;
  không logo thật ngoài đời; emblem lá-lửa theo TOT-naruto-logo.png; không dùng rồng/ngọc Dragon Ball.
- Không vẽ lại chân dung Iruka. Không vẽ logo mới.

LƯU VÀO: D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\
(được phép tách 2 thư mục con chan-dung\ và thuong-hieu\ để lệnh review không báo FAIL giả cho ảnh đục — nếu tách thì ghi rõ trong báo cáo)
Đặt đúng tên file như trên.

SAU KHI LƯU, KIỂM VÀ DÁN NGUYÊN KẾT QUẢ:
(a) chân dung — chạy từ D:\NRO_GoiBanGiao (cần Pillow + numpy):
python -m tools.artqa.artqa review --art-root D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\ --review D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\qa\review-09-chan-dung.json --reviewer "<tên bạn>"
    ĐẠT khi raw_source_status = PASS, issues = [], soft_alpha ≤ 0.15, content_boxes: đáy ≈ 1254 (vai chạm đáy), cao hộp ≥ 0,85 x 1254, không chạm mép trái/phải.
    (ảnh thương hiệu đục sẽ bị báo FAIL nguồn — đó là mong đợi, chỉ đọc kết quả của *_portrait_chroma.png)
(b) xem thử đúng cỡ trong game từ bản alpha: thu 200 và 104 px (lệnh Pillow ở 09-chan-dung-thuong-hieu.md mục 7.2 (b)) — còn ra mặt không.
(c) thương hiệu:  python D:\NRO_GoiBanGiao\docs\my-thuat\KIEM-TRA-ANH.py D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\ nguon
    KHÔNG chạy TAO-4-MUC-ZOOM.py cho nhóm thương hiệu.

BÁO CÁO khi xong (nói thật): cỡ thực + mode từng ảnh; nền tím phẳng hay có gradient (spread=?); quay phải; vai chạm đáy (đáy hộp = ?/1254);
ảnh nào KHÔNG CHẮC và vì sao; dán nguyên kết quả lệnh. DỪNG LẠI sau 3 ảnh, chờ duyệt. Không tự làm tiếp.
```

### GIAI ĐOẠN 2 — làm nốt

```
3 ảnh mẫu đã được duyệt. Làm nốt theo 09-chan-dung-thuong-hieu.md mục 5–6, giữ nguyên MỌI luật giai đoạn 1, lưu cùng thư mục
D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\09-chan-dung-thuong-hieu\ , chạy lại (a)(c) sau mỗi nhóm, báo cáo trung thực từng ảnh.

Thứ tự (chân dung chỉ vẽ khi master/dải idle của nhân vật đó ĐÃ được duyệt ở hạng mục 01/02/04 — để khớp mặt):
  A. P2 boss còn lại (2): nhan_gia_mat_na_p2_portrait_chroma · hoc_vien_dao_tau_portrait_chroma (1254x1254)
  B. P2 thương hiệu (3): ui_splash_doc_x4 1080x1920 · ui_splash_ngang_x4 1920x1080 · ui_bia_cuahang_x4 1024x500 — PNG đục, chừa vùng trống để dán wordmark, KHÔNG chữ
  C. P2 NPC Làng Lá vẽ lại (14): hokage3 · kakashi · chu_ramen · ba_vukhi · thay_hocvien · linh_gac · nguoi_giao_nv · y_ta · ong_tazuna · co_be_hoa · tho_ren · nguoi_dua_thu · trongtai_thi · nha_chep_cuon _portrait_chroma.png
     (mặt khớp master mới của 02; ảnh 343 chỉ lấy nội dung tóc/đồ, bỏ pixel-art và hướng quay trái)
  D. P2 NPC bốn làng (12), cùng nhịp với master ở 02, thứ tự Cát → Sương Mù → Đá → Mây: <làng>_kazekage|mizukage|tsuchikage|raikage · <làng>_thay · <làng>_giao_nv _portrait_chroma.png
  E. P2 người chơi bốn làng (4): nguoi_choi_cat/smu/da/may_portrait_chroma.png — CHỈ sau khi 01 chốt dải đầu theo làng.
  F. P3 (5) — chỉ khi được bảo: ui_icon_app_hinh_x4 432x432 (alpha thật) · ui_icon_app_nen_x4 432x432 (đục) · ui_nen_dangnhap_tet/he/trungthu_x4 1024x576 (đục).
KHÔNG vẽ lại iruka_portrait_chroma.
```

---

*Các khối trên chỉ trỏ tới file có thật (kiểm bằng Test-Path 06/09/2026: 100 đường dẫn `D:\…` tồn tại; 27/27 tên ảnh mẫu và mọi cỡ ghi trong khối khớp đúng manifest `NN-*.json`) và lệnh có thật trong `_KHUON-DOT-4.md` §D. Năm đường dẫn `--review …\*.json` (02, 03, 04, 07, 09) là **file đầu ra** — chưa có là đúng, `artqa review` tự tạo thư mục cha (`tools/artqa/artqa.py:445`). Thư mục LƯU VÀO đã tạo sẵn kèm `.gitkeep`. Khi một brief đổi tên file hay số khung, sửa khối tương ứng ở đây cùng lúc — agent tạo ảnh chỉ đọc khối này và brief, không đọc lịch sử.*
