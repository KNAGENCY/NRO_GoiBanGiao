# Soi map 1 Lang La — lan 1

Ngay 22/08/2026. 5 lang kinh doc lap soi anh trong game, moi phat hien qua mot luot phan bien doi khang.
**42 loi song sot** sau phan bien. Diem hoan thien: **45/100**.

> Anh chung: `pano0.png`, `pano2.png`, `pano4.png`, `sach2A.png` (trong D:/NRO_Project/)

---


## CHAN DUONG

### 1. Núi Hokage bị chém ngang phẳng lỳ bằng một đường thẳng chạy suốt chiều rộng màn hình

**O dau:** pano0.png / pano2.png / pano4.png — hàng y=343 sang y=344, chạy suốt x=0..1453 (rõ nhất ở pano0 vùng x 0-580 và x 700-1454, nơi không bị nhà che)

Chân dãy núi Hokage không tan dần vào không khí mà bị cắt cụt bằng một nhát dao thẳng băng ngang hết màn hình. Ngay trên đường cắt vẫn còn thấy vân đá xám; ngay dưới đường cắt là trời xanh phẳng tuyệt đối. Nguyên nhân: cả 4 lớp nền đều là PNG chế độ RGB (KHÔNG có kênh alpha) nên khi vẽ chồng, mỗi lớp là một hình chữ nhật đặc — mép TRÊN của lớp trước (b01, bắt đầu đúng ở y=344) xén phẳng lớp sau (b02). So với concept đã duyệt thì sai hẳn: trong concept núi là một khối lớn chạy xuống khuất sau mái nhà, không có đường cắt nào.

**Bang chung:** Đo bằng python/PIL: chênh lệch RGB trung vị giữa hàng 343 và 344 là (-12, -2, +6); độ đồng nhất của bước nhảy này trên toàn bộ 1454 cột là 97% (pano0), 73% (pano2), 63% (pano4) — nghĩa là gần như mọi cột đều nhảy đúng một lượng như nhau, dấu hiệu của mép ảnh chứ không phải nội dung. Màu hàng 344 = (10, 86, 155) trùng KHỚP TUYỆT ĐỐI với hàng pixel số 0 của b02... đúng hơn là của b01.png (đo được (10, 86, 154)). Kiểm tra chế độ ảnh: b00/b01/b02/b03.png đều là mode 'RGB', không alpha. Ảnh cắt phóng to: p0_horizon.png trong thư mục scratchpad.

**Cach sua:** Chọn 1 trong 3: (a) chuyển b01/b02/b03 sang PNG có alpha và xoá phần trời của lớp đứng trước thành trong suốt để các lớp thật sự chồng nhau; (b) hạ mép trên của lớp b01 xuống dưới đường chân núi của b02 (tức dưới y≈420 màn hình) để chỗ nối rơi vào vùng trời phẳng; (c) nếu bắt buộc giữ ảnh đặc thì vẽ thêm dải chuyển mờ (fade) cao 40-60px ở mép trên mỗi lớp thay cho mép cứng.

<sub>Lang kinh `lop-nen` · phan bien: Tôi cố bác bỏ nhưng không bác nổi — đây là lỗi rõ nhất cả map. Tự crop pano0 vùng (0,300)-(600,400) phóng 3x (file v_p0_horizon.png trong scratchpad): vân đá và mặt Hokage chạy tới đúng một đường ngang tuyệt đối thẳng rồi DỪNG, dưới là trời xanh trơn, không hề có tán cây/sương/mái nhà che. Không thể nhầm với 'ý đồ nghệ thuật'. Số đo của tôi: bước nhảy R trung vị hàng 343→344 = -12 (pano0), -11 (pano2), -10 (pano4); tỉ lệ cột nằm trong ±4 của trung vị = 92.9% / 66.3% / 58.5% (người báo cáo ghi 97/73/63 — hơi thổi nhưng cùng kết luận: đây là mép ảnh chứ không phải nội dung). pano0 hàng 344 trung bình x350-1454 = (10.0, 86.0, 155.0), b01.png hàng 0 = (10, 86, 154) — khớp, đúng là mép trên lớp b01 xén lớp núi. Kiểm chế độ ảnh: b00/b01/b02/b03 đều mode 'RGB', không alpha — đúng như báo cáo. Đối chiếu concept đã duyệt (concept-map1-langla.png): núi là một khối liền chạy xuống khuất sau mái nhà, không có đường cắt nào. Giữ nguyên mức chặn đường.</sub>

### 2. Một cây cột gỗ lơ lửng giữa trời, lặp dọc từ đỉnh màn hình xuống tận đường chân núi

**O dau:** pano2.png x=582..597, y=0..343 · pano0.png x=602..617, y=0..345 · pano4.png x=572..587, y=0..~230 (phần dưới bị nhà lớn che)

Giữa vùng trời có một cột gỗ màu be rộng 16px dựng đứng từ mép trên màn hình xuống đúng đường cắt của lớp núi. Nhìn gần thấy rõ đó là một miếng sprite gỗ (thân cột + thanh ngang + chống chéo hình chữ X) được lặp lại theo chiều dọc chứ không phải vật thể có thật. Concept đã duyệt không có bất cứ cột nào trên trời. Đây là thứ dễ bị người chơi chụp màn hình đem đi chế nhất.

**Bang chung:** Lọc pixel ấm (R>120 và R>B+20) trong dải trời y 150-230: chỉ đúng một cụm cột liên tục ở mỗi ảnh — pano0 x 604-615, pano2 x 584-595, pano4 x 574-585, không có cụm nào khác trên toàn bộ 1454 cột. Tự tương quan theo chiều dọc cho chu kỳ lặp 76px (đỉnh r=0.71) và 152px (r=0.74) → đúng là một sprite lặp. Khoảng cách từ cột tới đường nối dọc của lớp nền là 225px, GIỐNG HỆT ở cả ba ảnh (609-384, 589-364, 579-354) → cột trôi cùng tốc độ với lớp nền, chu kỳ 768px. Đã dò template 16x76 này qua toàn bộ 1696 file PNG trong Assets/Resources/res: KHÔNG khớp file nào ở tỉ lệ 1:1, và cũng không có trong b00/b01/b02/b03 (quét cột màu ấm trên cả 4 file: rỗng). Ảnh cắt phóng to: pole_zoom.png, crop_bamboo.png.

**Cach sua:** Truy ngược danh sách đối tượng được vẽ đè lên nền (lớp vẽ sau b00-b03, chu kỳ 768px, offset 225px so với mép lớp nền) rồi gỡ bỏ. Sau khi gỡ, chụp lại đúng 3 vị trí camera này (X=700 và X=1238) để xác nhận dải trời sạch.

<sub>Lang kinh `lop-nen` · phan bien: Tôi định bác cái này nhất (nghi là cột cổng torii hoặc vật thể hợp lệ) nhưng phóng to thì hết cãi. Crop pano2 vùng (570,120)-(610,280) phóng 8x (v_pole_zoom.png): rõ ràng là một sprite cột gỗ quấn dây, có đai ngang và thanh chống chéo hình chữ X, LẶP LẠI y nguyên theo chiều dọc giữa nền trời và mây — không có chân, không có đỉnh, không đổ bóng, không dính vào bất cứ vật gì. Tự lọc pixel ấm (R>120 và R>B+20) trong dải y150-230 trên cả 1454 cột: đúng một cụm duy nhất mỗi ảnh — pano0 x=604..615, pano2 x=584..595, pano4 x=574..585, không cụm nào khác. Chiều dọc: chạy liên tục từ y=0 tới y=343 (pano0 và pano2, dừng đúng ở đường cắt núi), pano4 tới y≈277 rồi bị nhà lớn che. Tôi tự kiểm nguồn: quét cột ấm liên tục dài nhất trong b01 = 34px (cột 570), trong b02 = 45px (cột 389), b00 và b03 không có pixel ấm nào → cột này KHÔNG nằm trong 4 file nền, đúng như báo cáo. Tôi cũng tự dò template 12x76 lấy từ pano2 qua toàn bộ PNG trong res/x2: khớp tốt nhất chỉ đạt sai số trung bình 32/255 (invenxd.png) — tức là không khớp file nào. Concept đã duyệt không có cột nào trên trời. NÂNG MỨC từ 'nặng' lên 'chặn đường': vật thể vô nghĩa này nằm gần chính giữa khung hình ở CẢ BA vị trí camera, cao suốt từ mép trên xuống chân trời, và là thứ đầu tiên mắt bắt vào. Không thể phát hành với nó. Lưu ý cho người sửa: tốc độ trôi đo được rất chậm (camera dịch 1076px màn từ pano2 sang pano4, cột chỉ dịch 10px) — có thể là lớp nền parallax chậm HOẶC lớp parallax ~0.7 bị cuộn vòng 768px; đừng chỉ tìm trong lớp vật thể tiền cảnh.</sub>

### 3. Không nhìn thấy nhân vật trong bất kỳ ảnh nào — camera không xuống tới chỗ nhân vật đứng

**O dau:** Cả 4 ảnh (pano0/pano2/pano4/sach2A). Dòng HUD gỡ lỗi ở pano2 vùng x 177-290, y 65-100 ghi "X: 700 - Y: 576"; pano4 ghi "X: 1238 - Y: 576"; pano0 "X: 60 - Y: 576". Quét toàn khung 1454x865 không có sprite nhân vật.

Đây là khía cạnh MỚI của lỗi rơi xuyên đất đã biết: không phải nhân vật đứng ở đáy map cho ta nhìn thấy, mà camera bị kẹp (clamp) cách biên map đúng 1 ô nên nhân vật nằm HOÀN TOÀN ngoài khung hình. Tính ra: camera đáy nhất là cam_y = 119 đơn vị game, tầm nhìn cao 432 đơn vị -> điểm thấp nhất nhìn thấy được là y = 551. Nhân vật ở y = 576, tức thấp hơn mép dưới màn hình 25 đơn vị = 50 pixel. Hậu quả cho người chơi: mở game ra thấy một bức tranh phong cảnh, không thấy mình đâu, không biết mình đang ở đâu, không cách nào tự kiểm chứng là mình có đi được hay không.

**Bang chung:** Đọc map-1-collision.bin (header 52,24,24,24): đặc từ hàng 19 trở xuống, cột 0 và 51 đặc toàn bộ. Đo camera từ 3 ảnh bằng cách khớp vị trí vật thể với ô lưới trong res/mymap/1.bytes: pano0 screen_x = 2*gx-49, pano2 = 2*gx-674, pano4 = 2*gx-992; screen_y = 2*gy-238 ở cả ba. Từ đó cam_y=119, đáy nhìn thấy = 551 < 576.

**Cach sua:** Sửa gốc là va chạm client (tile 60-82 chưa khai loại trong tile_set_info.bin). Nhưng NGAY BÂY GIỜ nên nới kẹp camera theo trục Y để nó luôn bám nhân vật kể cả khi nhân vật lọt ra ngoài biên map — có thế mới nhìn thấy được lỗi mà sửa. Kèm theo: chặn cứng y của nhân vật ở mức sàn va chạm server (456) để không bao giờ rơi quá đáy.

<sub>Lang kinh `choi-duoc` · phan bien: XÁC NHẬN, và tôi có bằng chứng mạnh hơn cả người báo cáo. So sánh từng pixel pano0.png với sach2A.png (hai ảnh này camera GIỐNG HỆT nhau — dịch chuyển đo được = 0 px, sai số 0.0): toàn bộ 1454x865 pixel chỉ khác nhau ở đúng 752 pixel nằm trong hộp x 195-284, y 84-100 — tức là CHỈ dòng chữ HUD 'X: 60' đổi thành 'X: 370'. Nhân vật đi 310 đơn vị game mà màn hình không đổi một pixel nào -> sprite nhân vật hoàn toàn không được vẽ trên khung hình. Số học cũng khớp: tôi đo mặt cỏ hàng 18 (game y=432) nằm ở screen y=626 -> cam_y = 432 - 626/2 = 119 (đúng con số người báo cáo đưa ra). Đáy màn hình 865 px -> game y = 119 + 432.5 = 551.5, trong khi nhân vật ở y=576 -> nằm dưới mép màn hình 24.5 đơn vị = 49 px. Camera bị kẹp đúng 1 ô (24 đơn vị) mỗi biên: tôi đo bằng cross-correlation dải nền pano0->pano4 = 946 px = 473 đơn vị, đúng bằng (1248-727) - 2*24. Lưu ý phản biện: quy kết 'lỗi camera' là sai địa chỉ — kẹp 1 ô là chuẩn của thể loại; gốc vẫn là lỗi va chạm đã biết #4. Nhưng hệ quả 'người chơi không thấy mình đâu' là THẬT và là chặn đường.</sub>


## Nang

### 4. Cột pixel rác rộng 7px lộ ra đúng ở góc mép gò

**O dau:** pano0.png x 521–527, y 626–673 · pano2.png x 1143–1149, y 626–673 · pano4.png x 823–829, y 626–673 (cùng một chỗ trong map nhìn từ 3 vị trí camera: mép PHẢI của cả hai gò, tức cột 11 và cột 37, hàng 18). Nguồn ảnh gốc: D:/NRO_Project/Assets/Resources/res/x2/t/1$66.png

Ngay chỗ mặt đất tụt xuống một bậc — chỗ mắt người chơi bám vào nhiều nhất — có một dải dọc rộng 7px màu xám kẻ ngang như mã vạch / cái thang nhỏ. Nó không phải đất, không phải bóng, mà là pixel rác. Lỗi nằm sẵn trong file tile: 7 cột cuối của tile 66 mỗi cột chỉ là một màu phẳng lặp lại cột texel cuối cùng, nghĩa là ảnh tile chỉ vẽ thật 41/48 cột, 7 cột còn lại bị 'bôi' ra.

**Bang chung:** Đo trên ảnh chụp game: tại pano4 từ x=824 đến x=829, mọi pixel trong cùng một hàng có giá trị RGB y hệt nhau, suốt 48 hàng liền (y 626–673) — dấu hiệu kinh điển của bôi mép texture. Chạy bộ dò 'cột màu phẳng' trên cả 3 ảnh pano: mỗi ảnh tìm được đúng 1 chỗ, đúng toạ độ nói trên, không nơi nào khác. Mở thẳng 1$66.png kiểm chứng: cột 41–47 mỗi cột một màu phẳng — hàng 10 = (131,121,113) lặp 7 lần, hàng 30 = (80,57,40) lặp 7 lần. Phóng 14x vùng pano4 x776–840 y618–678 nhìn thấy rõ dải xám kẻ ngang.

**Cach sua:** Vẽ lại 7 cột bên phải của 1$66.png cho vân đất chạy hết tới mép (hoặc lật ngang 1$65.png để dựng bản mặt phải rồi xoá luôn dải phẳng của nó). Chỉ cần thay file ảnh, không phải sửa 1.bytes.

<sub>Lang kinh `dia-hinh` · phan bien: Tôi tự chạy bộ dò 'cột giống hệt cột bên trái' trên dải y626-673 của cả 3 ảnh: pano0 trúng đúng 1 chỗ x=521, pano2 x=1143, pano4 x=823 — không chỗ nào khác, khớp 100% báo cáo. Đọc pixel thô pano0 y=650: x521..527 đều (87,55,30) y hệt nhau, x520=(68,44,27), x528=(97,60,34). Mở 1$66.png bằng numpy: cột 42-47 giống hệt cột trước nó → cột 41-47 (7 cột) là bản sao ngang của cột 41. Phóng 8x pano2 x1090-1180 y600-690: NHÌN THẤY RÕ một dải xám kẻ ngang như cái thang, lạc lõng giữa vách đất nâu. ĐÍNH CHÍNH mô tả: các cột đó KHÔNG phải 'mỗi cột một màu phẳng' (theo chiều dọc chúng vẫn đổi màu: hàng 10 = 131,121,113; hàng 30 = 80,57,40) — đúng bản chất là smear ngang, cột 41 bị nhân bản 7 lần. Con số và toạ độ thì chuẩn từng pixel. Giữ mức 'nặng' vì nó đọc ra như texture hỏng, không phải xấu, và xuất hiện ở cả 2 gò.</sub>

### 5. Mặt trước gò ghép từ 5 chất liệu khác nhau, có đường cắt ngang cứng chia đôi vách

**O dau:** pano2.png x 766–1150, y 578–673 (gò cột 30–37) · pano0.png x 240–528, y 578–673 (gò cột 6–11) · pano4.png x 446–830, y 578–673. Đường cắt ngang nằm ở y ≈ 626.

Vách gò cao 2 ô (96px màn hình). Nửa trên là đất rễ nâu sẫm của ô cỏ 62; nửa dưới lại lấy nguyên các ô dùng để lấp lòng đất — lần lượt 65 (vân đá sáng) – 60 (nâu phẳng) – 61 (sỏi) – 60 – 60 – 60 – 66. Kết quả: một bức vách bị chia đôi bằng một đường ngang thẳng băng, nửa dưới lại bị chia tiếp thành các ô vuông 48px khác chất liệu, trong đó ô sỏi nổi lên như một miếng vá sáng giữa vách.

**Bang chung:** Đo độ sáng trên pano2 tại vùng vách (x 820–1100): dải y 617–623 = 42.8; dải y 627–633 = 66.1 → nhảy +55% trong đúng 1 hàng pixel, kéo dài liền 384px. Đọc D:/NRO_Project/Assets/Resources/res/mymap/1.bytes (52x24, 2 byte đầu là 52 và 24): hàng 17 cột 30–37 toàn ô 62, còn hàng 18 cột 30–37 = 65,60,61,60,60,60,60,66 — đúng bộ ô lấp lòng đất. Gò kia y hệt: hàng 18 cột 6–11 = 65,61,60,60,60,66. Ảnh phóng 2.6x pano2 x760–1180 y560–690 cho thấy đủ 5 mảng.

**Cach sua:** Vẽ một bộ ô 'vách gò' riêng (thân vách + chân vách) khớp tông với đáy ô cỏ 62, rồi thay các ô 60/61 ở hàng 18 nằm dưới gò bằng bộ đó trong 1.bytes. Không dùng ô lấp lòng đất làm mặt vách.

<sub>Lang kinh `dia-hinh` · phan bien: Đo lại độ sáng pano2 x820-1100 từng hàng: y=625 → 39.2, y=626 → 63.6 (nhảy +62% trong 1 hàng pixel), dải 617-623 = 42.7 vs 627-633 = 66.1 — gần đúng số báo cáo (42.8/66.1). Đường nhảy chạy liền 384px đúng bề ngang gò. Đọc 1.bytes (52x24, khớp header): hàng 17 cột 30-37 toàn 62; hàng 18 cột 30-37 = 65,60,61,60,60,60,60,66; gò kia hàng 18 cột 6-11 = 65,61,60,60,60,66 — trùng khít báo cáo. Đo độ sáng từng ô hàng 18 của gò: 65=60.8, 60=67.1 (4 ô giống hệt), 61=62.4, 66=56.1 → 4 giá trị rời rạc, biên 48px. Phóng 2x pano2 x740-1180 y560-700: NHÌN THẤY nửa trên là đất rễ nâu sẫm, nửa dưới là mảng nâu nhạt phẳng, cắt nhau bằng một đường ngang thẳng băng — rất chối. Concept đã duyệt thì vách là tường đá xếp có cây cỏ rủ, càng lệch. ĐÍNH CHÍNH nhỏ: ô sỏi 61 và ô 65 thực ra TỐI hơn ô 60 (62.4 và 60.8 so với 67.1), không phải 'miếng vá sáng' như mô tả.</sub>

### 6. Làng xa bị chém ngang: mái ngói và ngọn cây đứt lửng giữa thân, dưới đường cắt là trời trống

**O dau:** pano4.png — hàng y=561 sang y=562, rõ nhất vùng x 150-470 (một mái ngói đỏ, một tán cây xanh và mấy nóc nhà nhỏ bị cắt ngang thân). Cùng đường đó ở pano2.png vùng x 480-720 (ngọn cây làng xa bị cắt)

Lớp 'làng xa' (b01) chỉ lộ ra một dải mỏng khoảng 20px rồi bị mép trên của lớp mây đứng trước xén ngang. Kết quả là người chơi thấy một hàng mái nhà và ngọn cây lơ lửng, đáy bị cắt phẳng, ngay dưới là trời sáng hơn hẳn — đọc ra như một dải sticker dán lên trời chứ không phải làng ở xa.

**Bang chung:** Đo pano4 vùng x 180-360: hàng 561 = (133, 139, 145) → hàng 562 = (61, 161, 224); vùng x 850-930 bước nhảy R = -117. Đo pano2 vùng x 500-700: hàng 561 = (103, 148, 177) → hàng 562 = (114, 181, 225), chênh kênh Blue 48 đơn vị. Cùng một hàng y=561/562 ở cả pano2 lẫn pano4 → là mép ảnh cố định, không phải nội dung. Ảnh cắt phóng to 4x: p4_village_cut.png và p2_village.png trong scratchpad.

**Cach sua:** Cùng gốc với lỗi trên. Nếu chưa kịp làm alpha cho toàn bộ lớp thì tối thiểu hạ mép trên của lớp mây xuống dưới y≈640 (dưới chân làng xa), hoặc cắt bớt phần trời thừa ở đỉnh b01 để làng xa nằm gọn trong dải còn hiện.

<sub>Lang kinh `lop-nen` · phan bien: Tự crop pano4 vùng (100,520)-(500,600) phóng 4x (v_p4_village.png): thấy rõ một mái chùa nhỏ, một mái ngói đỏ lớn và một tán cây xanh đều bị một đường ngang cắt phẳng đáy ở y=561/562, dưới là trời xanh trơn — đúng như mô tả, không cần suy diễn. Crop pano2 (400,520)-(800,600) (v_p2_village.png) cũng thấy ngọn cây/mây bị xén phẳng cùng hàng y đó. Số đo: pano0 hàng 561 = (96.0,107.0,120.0) → 562 = (86.3,120.1,145.9); pano4 561 = (107.8,99.4,93.2) → 562 = (93.0,99.2,102.4). Cùng một hàng y ở cả 3 ảnh → mép ảnh cố định. ĐIỂM TÔI BÁC ĐƯỢC: phần quy kết nguyên nhân sai/chưa chứng minh. Người báo cáo nói do 'mép trên của lớp mây'. Tôi kiểm chuỗi hàng 562-564 (pano0: 86.3→88.4→89.5, tăng đều) KHÔNG có vết ringing tối-sáng như hàng 0-3 của b02/b03, nên hàng 562 không phải hàng 0 của b02/b03. Lỗi hình ảnh có thật 100%, nhưng ai sửa thì phải tự dò lại đúng lớp nào đang xén, đừng tin ngay là lớp mây. Mức 'nặng' là đúng — dải bị ảnh hưởng mỏng hơn lỗi số 1 nhưng nội dung bị cắt ngang thân thì đọc ra rất giả.</sub>

### 7. Chân cột đá góc phải nhà lớn treo lơ lửng giữa trời, hở 47px thấy nguyên nền trời

**O dau:** pano2.png vùng x 1150–1200, y 578–626 (rõ nhất ở x 1150–1170). Cùng lỗi đó trong pano4.png ở x 831–879, y 578–626. Ảnh soi: scratchpad/z_chan_lo_lung.png (phóng 8x) và z_nhalon_phai.png

Nhà lớn 2 tầng rộng 8 ô, đặt ở cột 31 nên chiếm cột 31–38. Nhưng gò cao (mặt đất ở hàng 17) chỉ chạy từ cột 30 tới cột 37. Cột 38 mặt đất tụt xuống hàng 18, thấp hơn 1 ô. Kết quả: cạnh phải nhà — gồm cột gỗ góc và HÒN ĐÁ KÊ CHÂN màu xám — chìa ra khỏi mép vách đất và lơ lửng trên không, dưới nó là trời xanh có mây. Đây là thứ người chơi nhìn thấy đầu tiên khi đi từ giữa map sang phải: một căn nhà 2 tầng đứng hẫng một góc.

**Bang chung:** Giải mã 1.bytes (52x24): slot 72 = nhà lớn ở hàng 9 cột 31, ảnh 1$72.png cỡ 384x384 = 8 ô x 8 ô -> chiếm cột 31–38, đáy ở hết hàng 16. Bảng độ cao trong ve-bo-cuc-map.py: cột 30–37 = hàng 17, cột 38 = hàng 18. Đo pixel pano2.png: mép vách đất kết thúc đúng x=1150; hòn đá kê chân nằm y 562–579 (màu xám 104,88,77); từ y 584 tới y 625 toàn màu trời (94,185,231); mặt đất mới bắt đầu lại ở y=626. Khoảng hẫng đo được = 47 px màn hình = 23,5 đơn vị game ≈ 1 ô.

**Cach sua:** Dịch nhà lớn sang cột 30 (khi đó nó phủ đúng cột 30–37 = đúng gò), hoặc kéo dài gò thêm 1 cột (đổi điều kiện `30 <= c <= 37` thành `30 <= c <= 38` trong ve-bo-cuc-map.py). Sửa xong phải chụp lại đúng khung này để xác nhận không còn trời dưới chân đá.

<sub>Lang kinh `vat-the` · phan bien: Tự kiểm và ĐÚNG, nhưng hạ mức. Giải mã lại 1.bytes: slot 72 ở (hàng 9, cột 31), 1$72.png = 384x384 -> cột 31-38; hàng 17 (mặt gò) chỉ có tile ở cột 30-37, cột 38 là bụi cỏ chứ không phải đất. Đo camera pano2 độc lập: bậc gò trái cột 30 (game x=720) rơi đúng màn x=767 -> camLeft=336.5; camTop=119 (mặt đất hàng 18 = màn y=626). Crop 4x vùng pano2 (1090,500)-(1260,680) [file v1_nhalon_phai.png]: nhìn RÕ hòn đá xám ở x 1152-1170, y 562-577, ngay dưới nó y 578-625 là trời xanh có mây, đất mới bắt đầu y=626. Lặp lại y hệt trên pano4 x 832-850 [v1b_nhalon_phai_p4.png]. Khoảng hẫng đo được là 48px (1 ô) chứ không phải 47px. Hạ 'chan-duong' -> 'nang': lỗi thị giác rất lộ nhưng thuần thẩm mỹ, không chặn đường đi (va chạm client vốn đã = 0, đó mới là thứ chặn đường).</sub>

### 8. Lớp "làng xa" chỉ lộ ra dải mái nhà cao 15px rồi bị chặt phẳng — bên dưới lại là trời, làng xa lơ lửng

**O dau:** pano4.png vùng x 0–520, y 545–565 (rõ nhất); pano2.png vùng x 440–820, y 540–565. Đối chiếu concept-map1-langla.png vùng trung cảnh x 600–1000, y 560–760.

File nguồn res/x2/bg/b01.png (768x288) chứa nguyên một khu làng xa rất đẹp — mái ngói đỏ, cổng torii, cây xanh — nằm ở dải y 200–288 của ảnh. Nhưng trong game chỉ lộ ra ~13–15px trên cùng của khu làng đó (y 547–562), tức chỉ thấy chóp mái, rồi bị cắt ngang phẳng lì. Dưới đường cắt đó lại là TRỜI + MÂY, nên cả khu làng xa trông như đang treo lơ lửng giữa không trung, cách mặt cỏ tiền cảnh 63px toàn trời. Concept thì trung cảnh đầy nhà xa mờ dần chồng lớp, nối liền chân núi xuống mái nhà tiền cảnh, tạo chiều sâu. Hệ quả: dải giữa (y 344–560) của game rỗng 50.9% là trời trơn, trong khi dải tương đương của concept (y 400–660) chỉ rỗng 21.7%.

**Bang chung:** Read pano2/pano4 + crop phóng 2x vùng pano4 (0,470)-(520,620): thấy rõ hàng mái nhà tí hon bị một đường ngang chặt đôi, dưới đó là mây. Quét cột pano4 x=300: y548 (45,151,223) trời → y550–560 màu làng (137,152,116)/(154,158,112) → y562 (62,161,223) trời trở lại. Mở res/x2/bg/b01.png xác nhận khu làng chiếm dải y 200–288 của ảnh gốc. Đo tỉ lệ trời/mây trơn: game dải y 344–560 = 50.9%, concept dải y 400–660 = 21.7%.

**Cach sua:** Kéo lớp b01 xuống thấp hơn (hoặc tăng chiều cao khu làng trong ảnh) sao cho khu làng xa nằm đè lên vùng ngay sau mái nhà tiền cảnh, không còn khoảng trời hở giữa làng xa và mặt đất. Đồng thời để đáy b01 tắt dần vào màu nền thay vì cắt cụt.

<sub>Lang kinh `doi-chieu` · phan bien: Cố bác bỏ nhưng không bác được — đây là phát hiện chắc nhất trong cả danh sách. Quét cột pano4: x=150 trời (44,150,220) tới y548, nội dung lạ y549-561 (228,228,231)/(153,122,118), rồi y=562 trở lại trời (57,163,226) và giữ trời tới tận mặt cỏ. x=300 và x=900 y hệt: nội dung y549-561, cắt tại y=562. Tôi crop pano4 (0,520)-(520,620) phóng 3x và NHÌN THẤY: mấy chóp mái ngói đỏ + một ngọn cây bị một đường ngang thẳng tuyệt đối chặt đôi, bên dưới là trời có mây trắng — chúng treo lơ lửng giữa không trung. Mở res/x2/bg/b01.png (768x288) xác nhận ảnh gốc có nguyên khu làng xa rất đẹp ở dải y≈195-288 (mái ngói, cổng torii, cây), tức trong game chỉ lộ ~13px trên cùng. Đúng cả về vị trí lẫn cơ chế. Giữ mức nặng.</sub>

### 9. Một cột thang tre dựng đứng lơ lửng giữa trời, cắt cụt ở y=344 — concept không hề có vật này

**O dau:** pano2.png x 570–600, y 0–344; pano0.png x 598–625, y 0–345; pano4.png x 565–590, y 0–240.

Có một cột thang/giàn tre chạy dọc từ mép trên màn hình xuống đúng y=344 rồi biến mất — nó không chạm đất, không nối vào công trình nào, và bị cắt cụt đúng tại đường nối của lớp nền núi. Nó dịch chuyển rất chậm theo camera (chỉ ~5px trong khi tiền cảnh dịch 368px) nên thuộc một lớp nền parallax, không phải HUD. Concept không có bất cứ thứ gì tương tự; danh sách 18 vật thể dự kiến của map cũng không có "thang". Đây gần như chắc chắn là vật thể sót lại từ map Ngọc Rồng cũ.

**Bang chung:** Crop phóng 3x vùng pano2 (540,0)-(660,380) thấy rõ hình thang tre nhiều nấc, bị cắt phẳng ở đáy. So vị trí giữa 3 ảnh: pano0 x≈610, pano2 x≈580, pano4 x≈575 trong khi tiền cảnh dịch 626px (pano0→pano2) và 320px (pano2→pano4) — chứng tỏ nó nằm ở lớp nền xa. Kiểm tra res/x2/bg/b02.png và b03.png không có hình thang này, nên nó là vật thể đặt riêng.

**Cach sua:** Tìm và gỡ vật thể này khỏi dữ liệu map (nó không nằm trong danh sách 18 vật thể dự kiến). Nếu muốn giữ, phải hạ xuống chạm đất và ghép vào một công trình.

<sub>Lang kinh `doi-chieu` · phan bien: Cố bác bỏ nhưng ngược lại còn thấy nó tệ hơn mô tả. Crop pano2 (560,0)-(620,360) phóng 4x: đó là một cột/giàn tre nhiều nấc có đai buộc chéo, vẽ rất chi tiết, chạy từ mép trên màn hình xuống và bị chặt phẳng lì ở đáy — không chân, không nối vào gì. Dò tự động vị trí ở hàng y=180: pano0 x607-614, pano1 x605-612, pano2 x587-594, pano3/pano4 x577-584 — dịch tổng cộng 30px trong khi tiền cảnh dịch 946px, tức nằm ở lớp nền parallax rất chậm, không phải HUD. Đáy: pano0/1/2 cắt đúng y=343 (trùng đường nối ngang ở phát hiện 1), pano3/4 xuống y=440. Tôi kiểm tất cả ảnh nền có sẵn — b00, b01, b02, b03, b10, b11, b12, b13, fog0 — không file nào chứa hình này (đếm pixel màu tre: chỉ light.png và b12 có, mà mở ra xem thì không phải cột này). Nghĩa là nó là một vật thể được đặt riêng, gần chắc là rác sót từ map Ngọc Rồng cũ. Nó hiện trên MỌI khung hình, giữa trời, không thể giải thích được. Nâng vừa → nặng.</sub>

### 10. Hàng rào cắm xuyên cột cổng torii và lấn vào lòng cổng

**O dau:** pano2.png vùng x 382-477, y 583-625 (cột phải cổng torii ở x ~400-440). Cùng vật thể đó trong pano0.png ở x 1007-1103, y 582-626.

Một cái hàng rào gỗ (tile 75) được đặt chồng đúng lên cột phải của cổng torii: các thanh ngang của rào đi xuyên qua thân cột đỏ và qua cả bệ đá, rồi chìa ra hai bên. Nhìn ở mức phóng to 4x thấy rất rõ rào nằm ĐÈ LÊN cột chứ không phải đứng cạnh. Với người chơi, đây là cái cổng chính của làng — vật thể bắt mắt nhất map — mà lại có một khúc rào mọc xuyên qua chân cổng và chặn một phần lối bước qua cổng. Nhìn là biết lỗi ngay.

**Bang chung:** res/mymap/1.bytes: tile 75 ở hàng 16 cột 22 (x 528-576 đơn vị game); cổng torii tile 74 ở hàng 13 cột 19 (vùng đặc x 458-574). Hai vùng chồng nhau 46 đơn vị. Đã crop pano2 (350,540)-(560,660) phóng 4x để nhìn tận mắt.

**Cach sua:** Dời hàng rào này ra khỏi vùng cổng — ví dụ đưa về cột 25-26 (sau lối lát đá) hoặc bỏ hẳn. Trước khi đặt vật thể, nên có một bước kiểm tra tự động: hai vật thể có hình chữ nhật đặc chồng lên nhau > 20% thì báo lỗi.

<sub>Lang kinh `choi-duoc` · phan bien: XÁC NHẬN, giữ nguyên mức nặng — tôi cố bác bỏ mà không bác được. Crop pano2.png (330,520)-(560,660) phóng 4x nhìn tận mắt: các thanh ngang gỗ của hàng rào chạy THẲNG QUA thân cột đỏ và qua cả bệ đá đen của cổng, không có bóng đổ, không có khoảng hở — rõ ràng là hai ảnh đè lên nhau chứ không phải rào đứng cạnh cột. Dữ liệu khớp: 1.bytes có tile 75 tại hàng 16 cột 22 -> x 528-576, y 384-432; tile 74 (torii, ảnh 240x240, opaque 4..236) tại hàng 13 cột 19 -> x 458-574, y 313-432. Chồng lấn x 528-574 và y 384-432. Vị trí trên màn hình đúng như báo cáo (x 382-477, y 583-625).</sub>

### 11. Có một cột thừng/tre nhạt màu treo lơ lửng giữa trời trong mọi khung hình

**O dau:** pano0.png x 604-613, y 0-343; pano2.png x 584-590, y 0-343; pano4.png x 577-581, y 0-276 (bị mái nhà lớn che nốt). Có mặt cả ở sach2A.png.

Một cột dọc màu be-vàng nhạt, rộng khoảng 7-10 pixel, chạy từ MÉP TRÊN màn hình xuống tới ngang chân dãy núi Hokage rồi cụt ngang. Phóng to 8x thấy nó có vân như dây thừng quấn nút và các nếp chéo hình chữ X, giống một cây tre hoặc sợi dây thừng. Nó không có đầu, không có gốc, không dính vào cái gì. Vì nó gần như đứng yên trên màn hình (chỉ trôi ~3% theo camera) nên người chơi đi đâu cũng thấy nó lơ lửng giữa trời, ngay khu vực trung tâm — trông như một sợi dây leo được, nhưng chẳng leo được và cũng chẳng phải cái gì cả.

**Bang chung:** Đã mở và quét toàn bộ 8 ảnh nền res/x2/bg/b00-b03 và b10-b13, cộng light.png và fog0.png — KHÔNG file nào có cột này (b01/b02 chỉ có một đường kẻ mảnh ở x=384). Nó cũng không nằm trong danh sách 18 vật thể của map (không có món nào là cột/dây). Đo dịch chuyển: cột nằm cách đường nối nền đúng 223 pixel ở cả 3 ảnh (pano0: 384 và 607; pano2: 364 và 587; pano4: 354 và 577) -> nó được vẽ cùng lớp với nền xa. Ở pano4 nó bị mái nhà lớn che, tức được vẽ SAU nền và TRƯỚC vật thể map.

**Cach sua:** Cần tìm ra nó do đâu: kiểm tra đoạn code vẽ lớp nền của client (chỗ lặp b01/b02) xem có sprite nào bị vẽ với chiều rộng/UV sai, hoặc có một ảnh thứ 5 lọt vào danh sách nền. Đây là thứ nên chặn trước khi cho người chơi thấy.

<sub>Lang kinh `choi-duoc` · phan bien: XÁC NHẬN, và tôi NÂNG mức từ 'vừa' lên 'nặng' — đây là thứ khó chối nhất trong cả danh sách. Crop pano0 (560,0)-(660,360) phóng 4x: rõ ràng là một cây tre/cột gỗ có các nút thừng quấn và nếp chéo hình chữ X, chạy từ mép trên màn hình xuống rồi cụt ngang ở y=343, không đầu không gốc. Tôi đo vị trí ở cả 4 ảnh: pano0 x 605-613, sach2A x 605-613, pano2 x 585-593, pano4 x 575-583 — camera đi 946 px mà cột chỉ trôi 30 px (parallax 3.2%) nên người chơi đi đâu cũng thấy nó ngay giữa trời. Tôi cũng đã tự kiểm phần 'không có trong file nền': mở xem b01/b02/b03 phóng nguyên cỡ (chỉ có đường nối dọc ở x=384/192, không có cột nào) VÀ chạy khớp mẫu dải cột (180x18 px) lên toàn bộ b00-b03, b10-b13, fog0, light ở CẢ 4 mức zoom x1-x4 — sai số thấp nhất là 46.9, tức không khớp file nào. Quét thêm 150 ảnh ứng viên trong Assets/Resources (lọc theo màu be) cũng không ra nguồn. Nó nằm ở lớp sau nền, trước vật thể map (bị mái nhà lớn ở pano4 che). Một vật thể lạ, to, giữa màn hình, không rõ nguồn gốc -> nặng.</sub>


## Vua

### 12. Cỏ và lối lát đá gặp nhau bằng một nhát cắt dọc, không có ô chuyển tiếp

**O dau:** pano2.png x 286–288 và x 622–624, y 626–673 (hai đầu lối lát đá, biên cột 19|20 và 26|27). Cùng chỗ đó nhìn từ pano0.png là x ≈ 910 và x ≈ 1248; từ pano4.png là x ≈ 300.

Đi từ bãi cỏ sang lối đá, màu nhảy từ nâu sẫm sang vàng đất sáng trong đúng 1 pixel, ranh giới là một đường thẳng đứng tuyệt đối cao 48px. Không có ô cỏ lấn sang đá, không có viên đá lẻ rơi ra cỏ, mép cỏ cũng không có túm cỏ nào phủ qua. Mặt lối đá lại bị cắt phẳng lì ở đỉnh trong khi ô cỏ bên cạnh có ngọn cỏ nhô lên, nên đường chân trời của mặt đất đứt quãng: cỏ – trơ – cỏ.

**Bang chung:** Đo pano2: thân ô cỏ (x 240–286, y 642–670) độ sáng 52.5; ô lát đá (x 290–336, cùng dải y) độ sáng 87.9 → chênh 67%. Bước nhảy trung bình giữa x=285 và x=289 là 36 đơn vị sáng. Phóng 8x hai vùng x250–350 y615–700 và x580–680 y615–700 thấy ranh giới là một đường dọc hoàn toàn phẳng, không một pixel pha trộn. 1.bytes hàng 18: cột 19 = 62, cột 20–26 = 69, cột 27 = 62 — không có ô chuyển tiếp nào.

**Cach sua:** Vẽ thêm 2 ô chuyển tiếp (cỏ-lấn-đá bên trái và bên phải) đặt vào cột 20 và cột 26; hoặc rải vài túm cỏ / viên đá lẻ đè lên đúng đường ranh để phá đường thẳng.

<sub>Lang kinh `dia-hinh` · phan bien: Đo pano2 dải y642-670: thân cỏ x240-285 = 52.1, lát đá x290-335 = 87.7 (chênh 68%). Bước nhảy 1 pixel: x=285 → 44.8, x=286 → 83.7. Đầu kia: x=621 → 81.9, x=622 → 45.7. Cùng chỗ đó ở pano0 nhảy tại x=912 và x=1248, ở pano4 tại x=302 — đều khớp lưới cột 20 và 27 (báo cáo ghi '≈910' và '≈300', lệch 2px, không đáng kể). 1.bytes hàng 18: cột 19=62, 20-26=69, 27=62, không ô chuyển tiếp — đúng. Phóng 6x pano2 x250-350 y610-700: NHÌN THẤY ranh giới là đường dọc tuyệt đối phẳng, và đỉnh ô lát đá bị cắt trơ trong khi ô cỏ bên cạnh có ngọn cỏ nhô lên → đường chân trời mặt đất đứt quãng. Giữ 'vừa' chứ không nâng: biên tile cứng vốn là chuyện thường của thể loại này, phần thật sự lỗi chỉ là silhouette đỉnh bị hụt.</sub>

### 13. Có một đường chân đất thứ hai (giả) chạy suốt bề ngang map ở y ≈ 673

**O dau:** cả 3 ảnh pano, dải y 672–675, chạy liền từ x=0 tới x=1454 (chỗ nào không bị nhà che)

Đáy ô mặt đất (hàng 18) và đỉnh ô đất sâu (hàng 19) không cùng tông, nên sinh ra một vạch ngang thẳng băng chạy hết map, nằm đúng 48px dưới đường cỏ thật. Nhìn ra thành mặt đất có hai tầng nền chồng lên nhau. Bằng chứng nó là lỗi ghép ô chứ không phải tầng địa chất: đi qua đoạn lối lát đá thì vạch này đảo chiều sáng-tối.

**Bang chung:** Đo pano2: hiệu độ sáng giữa dải (y 674–677) và dải (y 669–672) = +24.7 ở các đoạn có cỏ (x 120–290, 630–760, 1160–1330), nhưng = −16.8 ở đoạn dưới lối lát đá (x 300–620). Cùng một đường, hai dấu ngược nhau. Ảnh phóng 2x pano2 x600–1100 y660–870 thấy vạch chạy thẳng không đứt.

**Cach sua:** Chỉnh tông đáy các ô 62 và 69 cho khớp đỉnh ô 60, hoặc thêm một hàng ô chuyển tiếp mờ dần giữa hàng 18 và hàng 19.

<sub>Lang kinh `dia-hinh` · phan bien: Đo pano2 quanh biên hàng 18/19 (thực ra hàng 19 bắt đầu ở y=674, không phải 673 — lệch 1 pixel so với báo cáo): đoạn cỏ x120-290 y673=40.1 → y674=63.8; x630-760 chênh +25.2; x1160-1330 chênh +25.2. Đoạn dưới lối lát đá x300-620: y673=80.3 → y674=63.2, tức −16.8. ĐÚNG hiện tượng đảo dấu mà báo cáo nêu, và con số −16.8 trùng khớp. Nhìn ảnh tổng thể thu nhỏ pano2 và ảnh phóng x250-350: có một vạch ngang thẳng, không đứt, chạy hết khung, nằm đúng 48px dưới đường cỏ. Lập luận 'đảo dấu nên không phải tầng địa chất' là hợp lý — một tầng địa chất thật không thể vừa sáng hơn vừa tối hơn lớp trên.</sub>

### 14. Tầng đất sâu 4 hàng gần như phẳng một màu, chỉ điểm vài ô sỏi hình vuông cạnh sắc

**O dau:** cả 3 ảnh pano, dải y 673–861 (hàng 19–22), suốt chiều ngang màn hình — chiếm khoảng 22% chiều cao khung hình

Gần một phần tư dưới màn hình là một mảng nâu gần như một màu duy nhất. Cả vùng chỉ có 2 mẫu ô: ô 60 (phẳng trơn) và ô 61 (sỏi). Các ô sỏi nằm rời rạc thành từng hình vuông 48x48 có cạnh vuông vắn nhìn thấy được, nên vùng này đọc ra thành 'bàn cờ ô vuông' chứ không phải một khối đất. Riêng hàng 21 gần như trống trơn: một dải cao 48px chạy hết map hoàn toàn không có chi tiết nào.

**Bang chung:** Đo độ lệch chuẩn màu của ô 60 = 4.7 / 3.0 / 1.5 trên thang 255 (tức gần như đơn sắc), so với ô cỏ 62 là 22.3 / 23.2 / 7.0. Đếm trong 1.bytes: hàng 19–23 chỉ dùng ô 60 và ô 61, không có mẫu thứ ba; hàng 21 chỉ có 2 ô sỏi trên 50 cột (cột 27 và 40). Ảnh phóng 2x pano2 x600–1100 y660–870 nhìn rõ từng ô sỏi là hình chữ nhật viền sắc, không hoà vào nền.

**Cach sua:** Thêm 3–4 biến thể ô đất sâu (vân đất, rễ dài, đá tảng, ô tối dần) rồi rải xen kẽ; hoặc phủ một lớp gradient tối dần theo độ sâu lên toàn vùng để đáy map không bị phẳng.

<sub>Lang kinh `dia-hinh` · phan bien: Đo trực tiếp file tile: 1$60.png std = 4.7/3.0/1.5 trên thang 255 → gần như đơn sắc, khớp y hệt báo cáo. (Ô 62 tôi đo được std 18.5/23.7/5.3, báo cáo ghi 22.3/23.2/7.0 — lệch chút nhưng cùng kết luận.) Đếm 1.bytes: hàng 19-23 chỉ dùng 60 và 61, không có mẫu thứ ba; hàng 21 chỉ có 2 ô 61 ở cột 27 và 40 trên 50 cột — đếm tay, đúng. Vùng này chiếm y674-865 = 192/865 ≈ 22% khung hình, đúng. Phóng pano2 x600-1100 y655-865: NHÌN THẤY rõ một mảng nâu phẳng lì với vài hình chữ nhật sỏi viền vuông vắn nằm rời rạc — đọc ra đúng như 'bàn cờ', và có một dải trống hoàn toàn chạy ngang giữa. Concept đã duyệt không hề có lớp cắt ngang lòng đất này.</sub>

### 15. Bụi cỏ ở chân gò bị mép gò xén đôi theo một đường thẳng đứng

**O dau:** pano0.png x 192–240, y 578–626 (bụi tile 81 đặt ở hàng 17 cột 5, gò bắt đầu ngay cột 6)

Bụi cỏ rộng 96px nhưng nửa phải bị khối gò vẽ đè lên, cắt ngang thân lá thành một đường dọc phẳng lì đúng tại x=240 — nhìn ra là bị 'xén' chứ không phải mọc sau vách. Thêm nữa gốc bụi hở khoảng 3px thấy trời ở y 623–625 trước khi chạm mặt cỏ.

**Bang chung:** Phóng 7x vùng pano0 x170–320 y555–640: lá bị cắt cụt tuyệt đối thẳng tại x=240, trùng khít biên ô. Mở 1$81.png: ảnh 96x48, tán lá phủ gần hết chiều ngang nên phần bị mất là thật. 1.bytes hàng 17: cột 5 = 81, cột 6 = 62 (bắt đầu gò). Đo pixel tại x=200: y622 = (74,105,53) là lá, y623–625 = xanh trời, y626 = cỏ.

**Cach sua:** Dời bụi sang cột 3–4 để nó nằm trọn trên nền thấp, hoặc thay bằng bụi nhỏ 48px ở cột 5, hoặc chọn một bụi có sẵn hình bám mép vách.

<sub>Lang kinh `dia-hinh` · phan bien: Xác nhận cả hai vế bằng pixel: pano0 y=600, x=239 còn là lá (116,101,60) thì x=240 đã nhảy sang nâu đất (81,50,27) — cắt đúng biên cột 6, tức mép gò. Và tại x=200: y=622 là lá (74,105,53), y=623-625 là TRỜI XANH (148,171,179 / 174,206,231 / 147,192,224), y=626 mới là cỏ → đúng khe hở 3px thấy trời dưới gốc bụi. 1.bytes hàng 17 cột 5 = 81, cột 6 = 62, khớp. Phóng 6x pano0 x170-320 y555-645: NHÌN THẤY tán lá bị chặt cụt tuyệt đối thẳng đứng giữa không trung, và sọc trời xanh len dưới gốc bụi. NÂNG mức từ 'nhẹ' lên 'vừa': người báo cáo tự hạ thấp — lá cây bị cắt phẳng giữa trời là loại lỗi mắt người tự động bắt được, khác hẳn một mảng màu phẳng ở góc tối.</sub>

### 16. Vạch tối 4px chạy suốt màn hình ở y=34, do 4 hàng pixel hỏng ở mép trên file b02.png và b03.png

**O dau:** cả 3 ảnh — hàng y=33 xuống y=37, chạy suốt x≈330..1453 (bên trái bị HUD che). Nguồn lỗi: D:/NRO_Project/Assets/Resources/res/x2/bg/b02.png hàng pixel 0-3, và b03.png hàng 0-3

Ở sát đỉnh màn hình có một vạch xanh navy đậm nằm ngang, dày 4px, kéo hết bề ngang. Đây không phải chỗ nối lớp thông thường mà là VIỀN HỎNG có sẵn trong file: mép trên b02.png bị một hàng tối thui rồi một hàng sáng quá đà (kiểu ringing khi resize ảnh). Khi lớp này được vẽ, viền hỏng đó hiện nguyên xi thành cái gạch chân đen giữa trời. b03.png bị y hệt — hiện chưa lộ vì mép trên của nó nằm ngoài khung, nhưng chỉ cần dịch lớp là lòi ra.

**Bang chung:** Đo file gốc: b02.png hàng 0 = (6, 58, 104), hàng 1 = (9, 74, 131), hàng 2 = (13, 96, 171), hàng 3 = (13, 100, 178), hàng 5 ổn định = (10, 90, 162) → hàng 0 tối hơn 45%, hàng 3 sáng vượt 10%. b03.png hàng 0-3 giống hệt. b01.png và b00.png KHÔNG bị (b01 hàng 0-5 đều = (10,86,154)). Trên màn hình, chuỗi chênh lệch giữa các hàng 34→35→36→37→38 là +(3,15,26), +(4,23,41), +(0,4,6), -(2,9,15) — trùng khít với chuỗi của b02.png hàng 0→1→2→3→4, chứng minh hàng màn hình 34 chính là hàng 0 của b02.png. Vạch này hiện ở 86% số cột (pano0), 91% (pano2 và pano4); ví dụ tại x=800: Blue hàng 33 = 228, hàng 34 = 104, hàng 40 = 165. Ảnh cắt: p0_flare.png, p4_topzoom.png.

**Cach sua:** Cắt bỏ 4 hàng pixel trên cùng của b02.png và b03.png rồi lưu lại (hoặc vẽ đè 4 hàng đó bằng màu hàng thứ 5), và kiểm tra lại toàn bộ 4 mép của cả 4 file nền — lỗi viền resize thường có ở nhiều cạnh.

<sub>Lang kinh `lop-nen` · phan bien: Nguồn lỗi trong file: XÁC NHẬN CHÍNH XÁC. Tôi đọc trực tiếp b02.png hàng 0 = (6.3,58.3,104.5), hàng 1 = (8.8,73.5,131.1), hàng 2 = (12.8,96.4,171.0), hàng 3 = (12.9,100.4,177.9), hàng 4-7 ổn định ~(10,91,163) → đúng kiểu ringing khi resize: 1 hàng tối sụt, rồi 2 hàng sáng vượt. b03.png hàng 0-3 giống hệt tới từng đơn vị. b00 và b01 sạch (b01 hàng 0-7 đều ~(10,86,154)). Chứng minh hàng màn hình 34 = hàng 0 của b02: tôi khớp được một phép biến đổi tuyến tính duy nhất (do lớp loé sáng phủ lên) — screen ≈ 0.76*b02 + (61,56,28) đúng cho CẢ hàng 34 lẫn hàng 38, sai số dưới 1 đơn vị. Không thể trùng hợp. Trên màn: cột có B(38)-B(34) > 20 chiếm 90.9% (pano0), 93.4% (pano2), 93.4% (pano4) — cao hơn cả con số 86/91% người báo cáo đưa. Nhìn crop 1:1 (v_top1x.png, pano4 vùng x700-1454, y0-120) thấy rõ một gạch ngang tối chạy hết bề ngang. HẠ MỨC từ 'nặng' xuống 'vừa' vì hai lý do: (a) 'dày 4px' hơi thổi — chỉ hàng 34 tối rõ (B=106 so với 156 ở hàng 38), hàng 35 tối nhẹ, hàng 36-37 thực ra SÁNG hơn nền, nên mắt đọc là một vạch mảnh ~2px kèm quầng sáng, không phải dải đen 4px; (b) nó nằm sát mép trên khung hình, phần trái bị HUD che, sức phá hình ảnh kém xa lỗi số 1 và 2. Vẫn phải sửa, và cảnh báo b03 sẽ lòi ra khi dịch lớp là hợp lý.</sub>

### 17. Nhà nhỏ bị gò đất nuốt mất nửa phải, hòn đá kê chân bên phải biến mất

**O dau:** pano0.png vùng x 241–385, y 578–626 (nửa phải căn nhà). Ảnh soi: scratchpad/z_nhanho_chon.png (phóng 4x) và z_buico_cut.png

Nhà nhỏ rộng 6 ô đặt ở cột 3 nên chiếm cột 3–8, đáy nằm ở hàng 18. Nhưng cột 6, 7, 8 lại là gò cao (mặt đất hàng 17), cao hơn đáy nhà đúng 1 ô. Nên 3 cột bên phải của căn nhà bị cỏ gò phủ lên, ăn mất 48px cuối. So sánh với ảnh gốc 1$70.png: căn nhà có HAI hòn đá kê chân ở hai góc dưới; trong game chỉ còn nhìn thấy hòn bên trái, hòn bên phải mất sạch. Chân tường và mép cửa sổ bên phải bị cắt ngang bởi vệt cỏ — trông như căn nhà bị chôn xuống đồi chứ không phải đứng trên đất. Cùng chỗ đó, bụi cỏ ở cột 5–6 cũng bị cắt dọc mất nửa phải (pano0 x 241–289).

**Bang chung:** 1.bytes: slot 70 hàng 12 cột 3; ảnh 1$70.png 288x288 = 6x6 ô -> cột 3–8. Bảng độ cao: cột 6–11 = hàng 17. Đo pixel pano0.png (camera trái = 23,5 đơn vị game, đỉnh = 119): quét mép cỏ theo cột cho ra y=578 liên tục từ x=240 tới x=528, và y=626 ở x<=230 và x>=530 — tức bậc lên đúng cột 6, bậc xuống đúng cột 11. Đáy nhà lẽ ra ở y=626 nên 48px cuối bị che. Đối chiếu sprite 1$70.png: hòn đá phải nằm ở toạ độ ảnh x 245–270 (= cột 8,1–8,6) — nằm trọn trong vùng bị chôn.

**Cach sua:** Hai lựa chọn: (a) dời nhà nhỏ sang cột 0 để nó nằm gọn trên nền thấp cột 0–5; (b) giữ cột 3 nhưng thu gò lại còn cột 9–11, để cột 6–8 hạ về hàng 18. Nhớ dời luôn bụi cỏ cột 5–6 ra khỏi mép bậc.

<sub>Lang kinh `vat-the` · phan bien: Tự kiểm và ĐÚNG cả chi tiết hòn đá. 1.bytes: slot 70 ở (12,3), 1$70.png 288x288 -> cột 3-8, đáy = hàng 18 = màn y 626; gò cột 6-11 mặt ở y 578. Render riêng sprite 1$70 [v2_sprite70.png]: có 2 hòn đá kê chân, hòn phải ở toạ độ sprite x 225-260, y 252-285 -> quy ra màn x 322-357, y 590-623 => nằm TRỌN dưới mép gò 578 nên mất sạch thật (ban đầu tôi tưởng nó còn, nhưng thứ còn nhìn thấy ở x 320-354 y 544-577 là bệ ĐÈN LỒNG slot 76 cột 7-8, không phải chân nhà). Crop 5x pano0 (180,540)-(400,650) [v2b_junction.png]: thấy rõ tường gỗ nhà bị cắt ngang ở y=577, và bụi cỏ cột 5-6 bị xẻ dọc đúng x=241 (chỉ còn nửa trái, lá đứt thẳng đứng). Hạ 'nang' -> 'vua': ở cỡ xem thường nó vẫn đọc được là nhà nép vào sườn đồi; thứ lộ nhất là bụi cỏ bị xẻ đôi.</sub>

### 18. Cổng torii cắm xuyên qua nhà vừa — thanh xà đâm vào tường tầng 2, chân cột đứng trên bậc cửa

**O dau:** pano2.png vùng chồng lấn x 239–335, y 386–626. Ảnh soi: scratchpad/z_torii_nhavua.png (phóng 4x)

Torii rộng 5 ô đặt ở cột 19 nên chiếm cột 19–23; nhà vừa rộng 8 ô đặt ở cột 13 nên chiếm cột 13–20. Hai cột 19 và 20 chồng nhau suốt từ hàng 13 xuống hàng 17. Nhìn trong game: thanh xà ngang dưới của torii chạy thẳng vào mảng tường trát trắng và ô cửa sổ tầng 2 rồi cụt luôn trong đó; cột đỏ chạy đè lên mái hiên tầng 1; bệ đá to của cột torii ngồi ngay trên nền hiên/bậc cửa của nhà. Không đọc được là cổng đứng trước nhà — đọc ra là cổng bị nhét vào trong nhà.

**Bang chung:** 1.bytes: slot 74 hàng 13 cột 19, ảnh 1$74.png 240x240 = 5x5 ô -> cột 19–23, hàng 13–17. slot 71 hàng 12 cột 13, ảnh 1$71.png 384x288 = 8x6 ô -> cột 13–20, hàng 12–17. Giao nhau: cột 19–20 x hàng 13–17. Kiểm bằng mắt trên ảnh phóng 4x: rõ chỗ thanh xà cụt trong tường và bệ đá đè lên sàn hiên.

**Cach sua:** Dời torii sang cột 21 (chiếm cột 21–25). Khi đó nó đứng trọn trên lối lát đá (cột 20–26), hở ra một ô so với nhà vừa (kết thúc ở cột 20), và vẫn nằm giữa lối đi như ý đồ ban đầu.

<sub>Lang kinh `vat-the` · phan bien: Chồng lấn CÓ THẬT nhưng mô tả thổi phồng. 1.bytes xác nhận torii cột 19-23 hàng 13-17, nhà vừa cột 13-20 hàng 12-17 -> giao cột 19-20. Crop 3x pano2 (150,360)-(420,640) [v3_torii_house.png]: cột đỏ trái của torii (màn x 280-315) đè lên mái hiên và tường phải của nhà, mũi xà cong màu xám của torii kết thúc lửng ngay trên mái ngói đỏ (x ~242-265) — đây là chỗ xấu thật. NHƯNG hai chi tiết người báo nói SAI: (a) thanh xà đỏ dưới nằm ở y 437-470, mảng tường trát trắng ở y 470-497 — xà nằm PHÍA TRÊN mảng tường, không 'đâm vào' và không cụt trong đó; (b) bệ đá cột torii (x 268-317, đáy y 623) đứng sát tường nhà trên nền đất, không đứng trên bậc cửa/sàn hiên. Hạ 'nang' -> 'vua'.</sub>

### 19. Cả 3 hàng rào đều dán đè lên vật khác, không cái nào đứng trên đất trống

**O dau:** pano2.png: (1) x 143–239, y 530–626 dán lên tường nhà vừa; (2) x 1007–1103, y 482–578 dán lên tường tầng trệt nhà lớn; (3) x 383–479, y 530–626 chắn ngang cột và bệ đá torii. Ảnh soi: z_hangrao_nhavua.png, z_fence_nhalon.png, z_fence_torii.png

Ba mảnh hàng rào đặt ở cột 17–18, 22–23 và 35–36. Cột 17–18 nằm gọn trong chân nhà vừa, cột 35–36 nằm gọn trong chân nhà lớn, còn cột 22–23 trùng chân phải cổng torii. Kết quả trong game: hai mảnh rào trông như tấm gỗ đóng phẳng lên mặt tường nhà (nhìn xuyên thấy vân gỗ tường phía sau), mảnh thứ ba thì thanh ngang chạy xuyên qua thân cột đỏ và bệ đá torii. Không mảnh nào rào cái gì cả — chỉ là rác thị giác. Đã đối chiếu ảnh gốc 1$72.png để chắc chắn hàng rào KHÔNG phải chi tiết vẽ sẵn trong nhà: mặt tiền nhà lớn trong sprite chỉ có cửa lùa, rèm đỏ và 2 lồng đèn, không có rào.

**Bang chung:** 1.bytes: slot 75 xuất hiện 3 lần — hàng 16 cột 17, hàng 16 cột 22, hàng 15 cột 35; ảnh 1$75.png 96x96 = 2x2 ô. Vùng phủ: cột 17–18 (trong nhà vừa cột 13–20), cột 22–23 (trong torii cột 19–23), cột 35–36 (trong nhà lớn cột 31–38). Nhìn ảnh phóng 4x thấy rõ thanh rào nằm trước vân gỗ tường và trước bệ đá torii.

**Cach sua:** Chuyển 3 hàng rào ra chỗ đất trống thật sự: cột 11–12 (khoảng hở giữa thùng gỗ và nhà vừa), cột 24–25 (mép lối lát đá), và bỏ hẳn cái thứ ba hoặc để nó nối tiếp thành một đoạn rào 2 mảnh liền nhau cho có nghĩa.

<sub>Lang kinh `vat-the` · phan bien: Vị trí ĐÚNG hết, nhưng chỉ 1 trong 3 cái thật sự xấu. 1.bytes: slot 75 (96x96 = 2x2 ô) ở (16,17), (16,22), (15,35) -> cột 17-18, 22-23, 35-36; đối chiếu vùng phủ nhà vừa 13-20, torii 19-23, nhà lớn 31-38 => cả 3 đều nằm trong chân vật khác. Kiểm bằng crop 5x: (a) pano2 (120,500)-(260,640) [v4_fence_house.png] — rào ở x 142-236 dán phẳng lên vách gỗ nhà vừa, nhưng nhìn ra vẫn giống lan can hiên, KHÔNG gây khó chịu; (b) pano2 (980,440)-(1130,600) [v4c_fence_bighouse.png] — rào x 1006-1104 đè lên cửa lùa nhà lớn, cũng chỉ hơi thừa; (c) pano2 (370,500)-(500,640) [v4b_fence_torii.png] — cái này XẤU RÕ: hai thanh ngang chạy xuyên thẳng qua bệ đá và thân cột đỏ của torii, đọc ra ngay là hai vật cắm vào nhau. Giữ mức 'vua' nhưng thực chất chỉ cái ở cột 22-23 đáng sửa gấp.</sub>

### 20. Núi Hokage — nhân vật chính của concept — trong game chỉ còn dải đá mỏng bị dao cắt phẳng ở chân

**O dau:** pano0.png vùng x 0–1454, y 240–345 (đỉnh đá y 253–290 tuỳ cột, chân núi cắt phẳng tại y 344). Đối chiếu concept-map1-langla.png vùng x 200–1900, y 130–515.

Trong concept, khối núi Hokage cao từ y=130 tới y≈515 trên khung 1024 — chiếm ~38% chiều cao khung, lấp kín toàn bộ nền sau làng, 5 khuôn mặt to và đọc được rõ. Trong game, dải đá chỉ cao 54–101px trên khung 865px (6–12% chiều cao), tức nhỏ hơn khoảng 4 lần về tỉ lệ; mặt Hokage chỉ cao ~40px, mờ nhoè gần như không nhận ra ở tỉ lệ 1:1. Tệ hơn: chân núi không bị nhà cửa che như concept mà bị CẮT NGANG bằng một đường thẳng tuyệt đối ở y=344 chạy hết chiều ngang màn hình, phía dưới đường đó là trời trơn — núi trông như một miếng dán bị kéo cắt. Đây là khía cạnh MỚI, không phải lỗi lặp 2 lần/đường nối dọc đã ghi trong sổ.

**Bang chung:** Đọc concept-map1-langla.png và pano0/2/4.png bằng Read. Quét cột dọc: đỉnh đá tại pano0 x=200→y253, x=400→y287, x=600→y284, x=900→y275, x=1200→y259, x=1400→y290; chân cắt phẳng y=344 ở mọi cột. Crop phóng 2x vùng pano0 (150,230)-(780,350) thấy rõ đường cắt ngang và lớp mù xanh đè lên mặt Hokage. Đo bước nhảy màu qua y=343 trên các cột chỉ có trời (x=10,20,30,40,1440,1445,1450): tổng chênh RGB 21–25 đơn vị, trên (22,88,148) → dưới (10,86,155). Ảnh gốc res/x2/bg/b02.png (768x384) xác nhận đá chỉ nằm ở dải y 205–300, phần dưới là gradient xanh tắt dần.

**Cach sua:** Vẽ lại b02 cho khối núi cao gấp 3–4 lần (chiếm ~35–40% chiều cao khung) và kéo chân núi xuống thấp hơn đường mái nhà tiền cảnh để nhà che chân núi như concept, thay vì để chân núi lơ lửng giữa trời. Nếu chưa vẽ lại kịp thì tối thiểu kéo dài phần gradient đáy b02 xuống tận đường cỏ để không còn đường cắt ngang, và giảm lớp mù phủ mặt Hokage để 5 khuôn mặt đọc được.

<sub>Lang kinh `doi-chieu` · phan bien: Tôi tự đo lại, số liệu tái lập được, nhưng cách mô tả bị thổi phồng ở hai điểm. ĐÚNG: quét cột pano0 (dò đá S<0.30, max>110) cho đỉnh đá x=200→y253, x=400→y246, x=600→y273, x=900→y275, x=1200→y260, x=1400→y243; đáy đá y≈310-315. Dải đá cao ~72px trên khung 865 = 8.3%. Concept cùng phép đo: x=1000→y135..481, x=1300→y197..475, tức ≥346px/1024 = ≥34% (còn bị nhà che nên thực tế cao hơn). Chênh ~4 lần — đúng. ĐÚNG: có bước nhảy màu cứng full chiều ngang tại y=343→344 kể cả ở cột chỉ có trời — pano0 x=10 (23,91,152)→(10,86,155), x=30 (26,92,152)→(10,86,155), x=1440 (21,89,150)→(10,86,155); y hệt trên pano2/pano4. Crop phóng 2x vùng (0,230)-(1454,370) nhìn thấy rõ đường ngang chạy suốt màn hình. SAI 1: không phải 'chân núi bị cắt'. Đá kết thúc y≈315 rồi tan dần vào mù — tôi mở res/x2/bg/b02.png xem: đá ở y205-300, dưới là gradient tắt về (10,86,154) từ y≈322. Cái bị chặt ở y=344 là LỚP MÙ/TRỜI đè lên trên (màu 26,92,152), tức mép dưới một layer phủ, không phải thân núi. Vẫn là lỗi (đường nối ngang chạy suốt mọi khung hình) nhưng khác bản chất. SAI 2: 'mờ nhoè gần như không nhận ra' — tôi crop 1:1 vùng pano0 (150,235)-(600,330) và đọc được: 5 khuôn mặt vẫn nhận ra rõ là mặt người, có mắt/mũi/tóc. Nhỏ thì đúng, 'không nhận ra' thì không. Vì vậy hạ nặng → vừa.</sub>

### 21. Khối đất nâu phẳng lì chiếm gần 23% mỗi màn hình — concept là vách đất nhiều tầng phủ cỏ dại và kè đá

**O dau:** Cả 3 ảnh, vùng y 660–855, toàn chiều rộng x 0–1454 (rõ nhất pano2.png x 0–1454, y 665–861).

Toàn bộ phần dưới mặt cỏ trong game là một khối nâu gần như đồng nhất, chỉ điểm vài ô "nhiễu" hình vuông rải rác ngẫu nhiên, biên ô 48px lộ rõ. Nó chiếm 22.7% chiều cao khung hình ở mọi màn nhưng không mang thông tin thị giác nào. Concept ở cùng vị trí là vách đất/đá nhiều tầng, mỗi tầng có mép đá lởm chởm, cỏ dại mọc trùm xuống, có lối mòn đất chạy chéo, có bụi cỏ đen tiền cảnh ở hai góc dưới làm khung. Chênh lệch đo được rất lớn: độ lệch chuẩn màu của vùng đất concept là (38.3, 29.1, 10.6) với 1349 nhóm màu, còn game chỉ (10.0, 6.6, 3.3) với 100 nhóm màu — kém phong phú ~13 lần.

**Bang chung:** Read + crop 2 dải đất đặt cạnh nhau: game (0,590)-(1454,865) của pano2 và concept (0,700)-(2048,1024). Tính bằng numpy: concept vùng y800–1000/x200–1800 → mean (70.2,53.8,31.6), std (38.3,29.1,10.6), 1349 nhóm màu (chia 8); game vùng y665–861/x0–1454 → mean (99.5,63.1,34.0), std (10.0,6.6,3.3), 100 nhóm màu. Tỉ lệ pixel vật liệu trung tính (đá/gỗ mộc, S<0.28): concept 26.4% vs game 13.4%.

**Cach sua:** Chia vách đất thành 2–3 tầng như concept, mỗi tầng có một hàng tile mép đá + cỏ dại rủ xuống; thêm 2–3 biến thể tile đất khác nhau về sắc độ để phá đơn điệu; thêm một lớp bụi cỏ tối làm tiền cảnh ở đáy màn hình.

<sub>Lang kinh `doi-chieu` · phan bien: Tôi chạy lại đúng phép đo bằng numpy và ra y hệt con số họ báo: game pano2 vùng y665-861/x0-1454 → mean (99.5,63.1,34.0), std (10.0,6.6,3.3), 100 nhóm màu; concept y800-1000/x200-1800 → mean (70.2,53.8,31.6), std (38.3,29.1,10.6), 1349 nhóm màu. Tỉ lệ chiều cao 196/865 = 22.7% cũng đúng. Crop pano2 (0,610)-(740,865) phóng 2x cho thấy thêm một thứ họ tả chưa đủ: các mảng 'nhiễu sỏi' là những HÌNH CHỮ NHẬT biên sắc lẻm 48px nổi trên nền nâu phẳng, trông như texture placeholder chưa xử lý. Điểm phải trừ: so sánh hơi lệch tay — vùng concept y800-1000 gồm cả cỏ, bụi, mép vách, tiền cảnh chứ không thuần đất, nên tỉ số '13 lần' là so vùng đầy chi tiết với vùng đất trần. Ngoài ra thư mục anh-da-tao ĐÃ có sẵn tile_dat_1 và tile_dat_2 (2 biến thể), nên đề xuất 'thêm 2-3 biến thể' là thừa một phần. Đây là phần dưới mặt đất, thể loại này chấp nhận đất trơn ở mức nào đó → hạ nặng → vừa.</sub>

### 22. Nhà nhỏ bị gò đất chồng lên, chôn mất đúng 1 ô (48px) chân nhà và bệ đá cột đèn lồng

**O dau:** pano0.png vùng x 100–380, y 540–640. Mép gò tại x≈241; mặt gò cao y=577; chân nhà (đế đá) y=622; mặt cỏ thấp y=625.

Ngôi nhà nhỏ ở đầu map đứng trên mặt đất thấp (đế đá chạm cỏ ở y≈622–625, khớp đúng). Nhưng gò đất cao hơn 1 ô bắt đầu tại x=241 lại được vẽ ĐÈ LÊN nửa phải của ngôi nhà: từ x=241 trở đi, cánh cửa gỗ và bệ đá của cột đèn lồng bị mặt cỏ gò (y=577) chặt phẳng, tức bị chôn 45–48px. Nhìn vào thì ngôi nhà như bị lún một nửa xuống đất. Concept không có tình huống này — mọi công trình đứng trọn vẹn trên một mặt bằng, và mỗi bậc cao độ đều được kè bằng tường đá.

**Bang chung:** Crop phóng 4x vùng panorama ghép (100,540)-(420,700) (= pano0 cùng toạ độ, vì pano0 là gốc panorama). Thấy rõ: bên trái x<241 đế đá nhà nằm trên cỏ thấp (đáy đế y=622, cỏ y=625 — chạm đúng); bên phải x>241 cánh cửa gỗ và bệ đá đèn lồng bị mép cỏ gò y=577 cắt ngang phẳng. Chênh 622−577 = 45px ≈ đúng 1 ô 48px ở zoom 2.

**Cach sua:** Hoặc dời ngôi nhà sang trái để nó nằm trọn trên mặt thấp, hoặc nâng nó lên đứng trọn trên gò, hoặc dịch mép gò sang phải khỏi thân nhà. Kèm theo nên đặt tile kè đá ở mọi mép gò để bậc cao độ trông có chủ đích.

<sub>Lang kinh `doi-chieu` · phan bien: Đo lại từng con số và khớp gần như tuyệt đối. Dò mặt cỏ pano0 (điều kiện g>b+35, g>50, r<g+35): mặt gò = y578 liên tục cho mọi cột x từ 245 tới ~545; mặt cỏ thấp = y626 cho x=20,60,100,140,180. Chênh đúng 48px = đúng 1 ô. Mép gò: quét cột x=239→240 còn là vân gỗ/đá của nhà (180,154,131), x=242 đã sang tối (75,61,56) → mép rơi vào x≈241, khớp báo cáo. Crop pano0 (80,480)-(440,680) phóng 4x nhìn rõ: nửa trái ngôi nhà đứng trên đế đá hiện đầy đủ chạm cỏ thấp, còn từ x≈241 sang phải thì tường gỗ + cánh cửa hông bị mặt cỏ gò y=578 chặt ngang. Bệ đá cột đèn lồng đáy đúng y≈577 tức nó đứng TRÊN gò, trong khi ngôi nhà nó gắn vào lại đứng thấp hơn 48px — hai bộ phận của cùng một công trình ở hai cao độ. Chỗ trừ điểm duy nhất: nhìn tổng thể ở tỉ lệ 1:1 nó đọc ra như 'nhà đứng sau một thềm đất' hơn là 'nhà bị lún', chưa phải thứ đập vào mắt ngay. Hạ nặng → vừa.</sub>

### 23. Không có kè đá và bậc thang đá — concept dùng chúng ở mọi bậc cao độ và ở lối lên torii

**O dau:** pano2.png vùng x 250–620, y 600–670 (chân torii và mép sân lát đá); pano0.png x 200–260, y 570–630 (mép gò). Đối chiếu concept-map1-langla.png x 640–860, y 720–800 (bậc thang lên torii) và x 0–2048, y 760–820 (tường kè đá chạy suốt).

Trong concept, cả dải làng đứng trên một bờ kè xây đá, và có một cầu thang đá 5–6 bậc dẫn từ mặt đất lên cổng torii — đó là chi tiết nói lên "đây là lối vào làng". Trong game hoàn toàn không có kè đá lẫn bậc thang: torii đứng thẳng trên một mảng lát đá phẳng, và mọi bậc cao độ đều là vách đất trần cắt vuông góc. Chỗ chuyển từ cỏ sang sân lát đá (pano2 x≈287) là một đường dọc cứng không có tile chuyển tiếp, không có cỏ mọc lên mép đá.

**Bang chung:** Đặt cạnh nhau crop game panorama (860,380)-(1180,680) phóng 3x và crop concept (480,420)-(900,830). Concept có bậc thang đá + tường kè đá xây khối; game có mặt lát đá phẳng, mép cỏ/đá cắt dọc không blend. Đo tỉ lệ pixel vật liệu trung tính (đá/gỗ mộc, S<0.28) trên toàn khung: concept 26.4% vs game 13.4%.

**Cach sua:** Bổ sung bộ tile kè đá (mặt tường + mép trên) và một cụm bậc thang đá đặt ngay dưới torii. Đây cũng là cách rẻ nhất để "lấp" bớt mảng đất nâu chết ở phát hiện 4.

<sub>Lang kinh `doi-chieu` · phan bien: Đọc cả 3 ảnh game: không có một viên kè xây nào, không có bậc thang nào; mọi bậc cao độ (y578 và y626) đều là vách đất trần cắt vuông góc, torii đứng thẳng trên mảng lát đá phẳng. Concept thì có cả tường kè đá chạy suốt dưới dãy nhà lẫn cầu thang đá lên torii — tôi nhìn thấy rõ khi đọc concept-map1-langla.png. Chi tiết 'mép cỏ sang sân lát đá là đường dọc cứng' cũng đúng và tôi đo được: pano2 tại y=628 màu nhảy từ (76,78,20) cỏ ở x=285 sang (133,85,40) đá ở x=286, không có pixel chuyển tiếp nào; y=635 và y=650 y hệt. Trừ điểm: con số '26.4% vs 13.4% pixel vật liệu trung tính' là chỉ báo gián tiếp yếu, đo trên hai ảnh khác hẳn bố cục thì không nói lên gì về kè đá; tôi không dựa vào nó. Đây là 'chưa vẽ' hơn là 'vẽ sai'. Giữ mức vừa.</sub>

### 24. Cả map chỉ có 1 cây, đặt ở sát biên phải; concept có 3 cụm cây xen giữa nhà và cành lá đen làm khung góc trên

**O dau:** Cây duy nhất: pano4.png x 1290–1454, y 350–620 (tán cây chạm cột cuối x=1453 — bị biên map cắt). Vùng thiếu cây: pano0.png x 370–1300, y 400–620. Đối chiếu concept-map1-langla.png: cây tại x 350–560 y 380–620, x 1640–1790 y 480–620, và cành lá đen khung tại x 0–140 y 20–230.

Concept dùng cây xanh làm nhịp thở giữa các dãy nhà: một cây lớn chen giữa nhà trái và cổng torii, một cây sau nhà phải, cộng thêm một cành lá sẫm màu ở góc trên trái để đóng khung bố cục. Trong game chỉ có đúng 1 cây, và nó nằm ở mép phải cùng của map, tán bị biên màn hình cắt (đo được lá cây vẫn kín tới cột x=1453). Toàn bộ khoảng giữa map (pano0 x 370–1300) không có một mảng xanh cao nào — chỉ nhà, trời và vài bụi cỏ nhỏ ngang mặt đất. Bảng màu vì thế thiếu hẳn mảng xanh lá trung gian giữa xanh trời và nâu đất. Lưu ý file cành lá đã được vẽ sẵn (docs/my-thuat/map-1-lang-la/anh-da-tao/tc_canh_la_x4.png) nhưng không thấy xuất hiện trong game.

**Bang chung:** Read pano0/2/4 và concept; crop panorama (1900,300)-(2400,720). Quét cột pano4 đếm số dòng có màu lá (H 60–160, S>0.2) trong y 350–620: x=1440→53, x=1446→66, x=1450→60, x=1452→56, x=1453→51 — tán cây kín tới cột cuối cùng, tức bị biên map cắt. Liệt kê thư mục anh-da-tao thấy có tc_canh_la_x4.png và tc_bui_truoc_x4.png nhưng không nhận ra chúng trong bất kỳ ảnh game nào.

**Cach sua:** Đặt thêm 2 cây nữa vào khoảng trống giữa map (quanh pano0 x 450–560 và x 1150–1260), dời cây hiện tại vào trong khoảng 60–80px để tán không bị biên cắt, và cắm lớp cành lá tc_canh_la ở góc trên trái như concept.

<sub>Lang kinh `doi-chieu` · phan bien: Tôi định bác cái 'bị biên map cắt' vì ảnh chụp có thể chỉ là biên màn hình, nên kiểm bằng cách ghép panorama: tương quan dải y470-620 giữa các pano cho pano0→1 dịch 26px, 1→2 dịch 600px, 2→3 dịch 320px, còn 3→4 dịch 0px với sai khác trung bình ĐÚNG BẰNG 0.0 (pano3 và pano4 trùng khít). Camera đã chạm mép phải ở pano3, nên biên phải pano4 chính là biên map. Đếm pixel lá (H 60-160°, S>0.2) trong y350-620 của pano4: x=1440→53 dòng, x=1446→66, x=1450→60, x=1452→56, x=1453→51 — tán vẫn kín tới cột cuối, tức bị biên cắt thật. Đọc pano0 và pano2 (cộng lại phủ hết map): không có cây cao nào, chỉ vài bụi cỏ ngang mặt đất → đúng là 1 cây duy nhất. Liệt kê docs/my-thuat/map-1-lang-la/anh-da-tao/ có tc_canh_la_x4.png và tc_bui_truoc_x4.png, không thấy chúng trong bất kỳ ảnh game nào — đúng. Chỉ trừ điểm nhẹ: thân cây vẫn hiện đủ, chỉ mất phần rìa tán, nên 'bị cắt' không nghiêm trọng như giọng văn. Giữ mức vừa.</sub>

### 25. Không có bóng đổ và không có lớp hoà chân công trình — đường cỏ thẳng tuyệt đối suốt cả map

**O dau:** pano4.png x 830–1454, y 616–632; pano0.png x 700–1454, y 620–635; nhìn tổng thể ở panorama ghép y 618–632 chạy suốt x 0–2400.

Mọi công trình trong game đặt lên một đường cỏ thẳng như kẻ bằng thước, chạy liên tục hết chiều dài map, và không vật thể nào có bóng đổ hay cỏ mọc chờm lên chân. Kết quả là nhà, giếng, cây trông như sticker dán lên nền. Concept làm ngược lại: mép đất lồi lõm, cỏ dại cao thấp khác nhau phá vỡ đường chân trời, chân mỗi công trình chìm trong một vùng tối và có cỏ/bụi mọc trùm lên nền móng.

**Bang chung:** Crop phóng 2x panorama (1900,300)-(2400,720): quán ramen và cây đều tiếp đất bằng một đường cỏ ngang tuyệt đối, không thấy vùng tối nào dưới mái hiên hay dưới tán cây. Crop concept (0,700)-(2048,1024) cho thấy mép đất gợn sóng, cỏ cao thấp và chân nhà tối dần.

**Cach sua:** Thêm sprite bóng đổ elip mờ dưới chân mỗi công trình lớn, và rắc thêm 8–12 bụi cỏ/đá nhỏ dọc đường cỏ để phá đường thẳng — hai việc rẻ nhưng đổi cảm giác "dán" thành "đứng".

<sub>Lang kinh `doi-chieu` · phan bien: Xác nhận bằng mắt lẫn bằng đo. Đo mặt cỏ: y=626 (đất thấp) và y=578 (gò) không đổi một pixel nào suốt chiều ngang, chỉ nhảy tại mép gò. Crop pano4 (830,590)-(1454,680) phóng 2x: thùng gỗ, ghế băng, chậu cây, biển hiệu quán ramen và cả gốc cây to đều tiếp đất bằng một đường cỏ kẻ thước, không một vùng tối nào dưới mái hiên hay dưới tán cây, không cọng cỏ nào chờm lên chân vật thể — chân vật thể bị cắt đúng tại đường cỏ. Đối chiếu concept: mép đất gợn sóng, cỏ dại cao thấp, chân nhà chìm trong vùng tối. Lưu ý công bằng: mặt đất phẳng tuyệt đối là bình thường của thể loại tile-based ngang này (chính Ngọc Rồng cũng vậy), nên phần 'đường thẳng' không tính là lỗi; phần thật sự là lỗi là KHÔNG CÓ bóng đổ nào cả, khiến vật thể trông như sticker dán. Giữ mức vừa.</sub>

### 26. Mép phải map có vùng chết: nhân vật đi ra khỏi khung hình, cây to bị cắt mất 1/3 vĩnh viễn

**O dau:** pano4.png — cây ở vùng x 1264-1454 (bị cắt ở đúng mép phải màn hình), y 340-700. Vị trí nhân vật theo HUD: X = 1238.

Camera chạy hết cỡ sang phải chỉ nhìn được tới x = 1223 đơn vị game, trong khi nhân vật đi được tới ít nhất x = 1238 (HUD pano4). Nghĩa là ở đoạn cuối map người chơi đẩy nhân vật ra ngoài khung hình bên phải rồi mất dấu. Cùng lý do đó, cây cổ thụ (vật thể to nhất map, ảnh 288x384) được đặt ở ô cột 47 nên trải dài tới x = 1272 — vượt cả chiều rộng map (1248) — nên 98 trên 288 pixel của cây (34%) không bao giờ hiện ra. Người chơi luôn thấy một cái cây bị dao xén dọc ở rìa màn hình.

**Bang chung:** Đo mốc trong pano4: quán mì (tile 73, ô hàng 12 cột 40, ảnh 384x288, vùng đặc x 966-1148 đơn vị game) hiện ở screen x 940-1304 -> screen_x = 2*gx - 992 -> camera trái = 496, phải = 496+727 = 1223. Cây (tile 80, hàng 10 cột 47, ảnh 288x384, vùng đặc 1128-1272) -> screen 1264-1552 trong khi màn hình rộng 1454.

**Cach sua:** Dời cây sang trái ít nhất 2 ô (đặt ở cột 44-45 thay vì 47) để cả tán nằm gọn trong vùng camera thấy được. Đồng thời rà lại kẹp camera trục X (đang thụt vào 1 ô mỗi bên) so với giới hạn di chuyển của nhân vật — hai con số này phải khớp, nếu không nhân vật sẽ luôn chạy được ra ngoài khung.

<sub>Lang kinh `choi-duoc` · phan bien: XÁC NHẬN phần đo, nhưng HẠ mức từ 'nặng' xuống 'vừa'. Tôi khớp mẫu (template match) file 1$80.png (288x384) lên pano4.png: khớp chính xác tại offset (0,0), vị trí vẽ = screen x 1262..1550, trong khi màn hình rộng 1454 -> mất 96/288 px = 33.3% (người báo cáo nói 34%, gần đúng). Đã crop pano4 (1180,220)-(1454,700) xem tận mắt: thân cây và tán bị một đường thẳng đứng xén cụt ở rìa phải, rất rõ. Camera phải nhất = 497 (không phải 496) -> mép nhìn thấy xa nhất là game x 1224; HUD pano4 ghi X: 1238 -> nhân vật ra ngoài khung 14 đơn vị = 28 px, có thật. Lý do hạ mức: (a) phần 'nhân vật mất dấu' chỉ là hệ quả của cùng gốc lỗi #1, không phải lỗi độc lập; (b) cây bị xén là lỗi THẨM MỸ ở rìa map, không chặn chơi. Bổ sung: cột 0 và cột 51 của 1.bytes trống hoàn toàn và kẹp camera 1 ô che đúng hai cột này — đó là thiết kế CÓ CHỦ Ý, không phải lỗi; chỉ riêng cái cây đặt ở cột 47 (tràn tới cột 52, quá cả mép map 1248) mới là đặt sai.</sub>

### 27. Hai ô đất mép gò bị lỗi ảnh: một cột sọc trắng-xám như mã vạch và một khối đất lệch màu, xuất hiện 4 chỗ

**O dau:** Cột sọc (tile 66): pano2.png x 1143-1150, y 626-673 và pano0.png x 519-527, y 626-673. Khối đất lệch (tile 65): pano0.png x 239-287, y 626-673 (dải đen đặc ở x 239-249) và pano2.png x 766-814.

Ảnh tile 66 (mép phải gò) có khoảng 8 pixel cuối bên phải là một dải sọc ngang xám-trắng xen nâu, không phải đất. Trên màn hình nó thành một cột kẻ sọc cao 47px nằm trong lòng đất, trông y như một cái thang nhỏ hoặc mã vạch cắm xuống đất. Ảnh tile 65 (mép trái gò) thì có 10 pixel đầu bên trái là một dải nâu đen phẳng lì, và phần thân còn lại là loại đất kẻ vằn ngang xám nhạt, khác hẳn đất xung quanh — nhìn thành một miếng vá 48x48 dán lệch màu. Vì mỗi gò có 1 mép trái + 1 mép phải và map có 2 gò, người chơi gặp 4 vệt lỗi này khi đi hết map.

**Bang chung:** Mở thẳng res/x2/t/1$66.png và 1$65.png phóng 10x: thấy rõ dải hỏng ở rìa ảnh. Đo pixel trên pano2 tại cột x=1146, y 626-672: các giá trị (131,121,113), (127,119,114), (107,101,97) — gần như xám trung tính (chênh lệch R-B chỉ 13-18), trong khi đất bên cạnh ở x=1170 là (86,52,27) nâu bão hoà (chênh lệch 59). res/mymap/1.bytes: tile 65 ở hàng 18 cột 6 và 30; tile 66 ở hàng 18 cột 11 và 37.

**Cach sua:** Vẽ lại hoặc cắt lại hai file res/x2/t/1$65.png và 1$66.png cho đủ 48x48 sạch (bỏ dải rác ở rìa, làm phần thân đồng chất với tile 60/61/62). Đây là sửa 2 file ảnh, không đụng logic — nên làm sớm vì nó lộ liễu nhất.

<sub>Lang kinh `choi-duoc` · phan bien: XÁC NHẬN có lỗi ảnh, nhưng HẠ mức từ 'nặng' xuống 'vừa' và SỬA số đo bị thổi phồng. Mở 1$60/61/62/65/66.png cạnh nhau phóng 9x: tile 66 đúng là có dải sọc ngang xám-trắng ~8 px ở rìa phải, tile 65 có dải nâu-đen phẳng ~10 px ở rìa trái, và THÂN cả hai tile là loại đất vằn ngang hoàn toàn khác loại đất lấm tấm của 60/61 -> nhìn ra ngay là miếng vá lệch. Crop pano0 (440,600)-(620,700) phóng 6x thấy rõ cột sọc kiểu mã vạch tại screen x 521-528, y 626-674. Số đo màu của người báo cáo bị THỔI PHỒNG: họ nói dải này gần như xám trung tính, chênh R-B chỉ 13-18; tôi đo trung bình cột x=1148 (pano2) và x=523 (pano0) trên y 630-670 ra (95.6, 75.4, 60.0), chênh R-B = 35.6 — nhạt màu hơn đất xung quanh (R-B 55-68) nhưng vẫn là nâu, không phải xám. Vị trí trong 1.bytes đúng: tile 65 hàng 18 cột 6 và 30, tile 66 hàng 18 cột 11 và 37 -> 4 chỗ. Là lỗi thật, lộ liễu, nhưng chỉ là 2 ô 48x48 -> 'vừa'.</sub>

### 28. Giếng nước bị vách đất của gò cắt đôi

**O dau:** pano2.png vùng x 766-806, y 578-626 (cả giếng trải từ x 678 đến 806). Ở pano4.png là vùng x 452-497.

Giếng được đặt ở cột 28 nhưng ảnh giếng rộng tới cột 30, mà cột 30 đã là chân gò cao. Kết quả: mái giếng và một phần bệ đá bị một vách đất dựng đứng, thẳng băng, xén mất. Nhìn thì giống như giếng bị chôn vào vách núi, nhưng vì vách cắt đúng theo đường thẳng của ô lưới (không có viền, không có bóng đổ) nên nó đọc ra là lỗi che khuất chứ không phải ý đồ.

**Bang chung:** res/mymap/1.bytes: giếng tile 77 (ảnh 144x144) ở hàng 15 cột 28, vùng đặc x 676-740 đơn vị game; gò 2 là các ô 62 ở hàng 17 cột 30-37 tức bắt đầu tại x 720. Chồng lấn 20 đơn vị = 40 pixel màn hình. Đã crop pano2 (660,530)-(880,670) phóng 4x, thấy rõ mái giếng bị cắt phẳng.

**Cach sua:** Dời giếng sang trái 2 ô (về cột 26) để cả ảnh nằm trên nền đất thấp, hoặc kéo gò 2 bắt đầu muộn hơn 2 ô (từ cột 32).

<sub>Lang kinh `choi-duoc` · phan bien: XÁC NHẬN, giữ nguyên mức vừa. Crop pano2.png (640,520)-(900,690) phóng 4x: mái giếng và bệ đá bị một vách đất dựng đứng xén phẳng lì bằng một đường thẳng đứng tại screen x=766, không viền, không bóng — đọc ra là lỗi che khuất chứ không phải chủ ý. Dữ liệu khớp: tile 77 (giếng, ảnh 144x144, opaque 9..135) hàng 15 cột 28 -> x 676.5-739.5, y 360-432; gò 2 là dãy tile 62 hàng 17 cột 30-37 -> bắt đầu x 720, y 408-432. Chồng lấn x 720-739.5 tại y 408-432, đúng 40x48 px trên màn hình. Thứ tự vẽ theo hàng (hàng 17 vẽ sau hàng 15) nên gò đè lên giếng.</sub>

### 29. Nhà nhỏ bị gò đất nuốt mất nửa dưới, chân nhà nằm ở hai độ cao khác nhau

**O dau:** pano0.png vùng x 240-383, y 578-626 (cả ngôi nhà trải từ x 95 đến 383).

Ngôi nhà nhỏ nằm vắt qua đúng mép trái của gò 1. Nửa trái của nhà (x 95-239) đứng bình thường trên nền cỏ thấp, chân nhà chạm đất đúng ở y 626. Nhưng từ x 240 trở đi, gò cao lên 48px và ăn mất 48 pixel dưới cùng của nhà: cửa lùa, tường gỗ và bậc đá bị một đường cỏ-đất cắt ngang. Người chơi nhìn thấy một ngôi nhà có nền nhà ở hai cao độ khác nhau — bên trái đứng trên mặt đất, bên phải chìm xuống đất.

**Bang chung:** res/mymap/1.bytes: nhà nhỏ tile 70 (ảnh 288x288) ở hàng 12 cột 3, đáy phần đặc ở y 432 (mặt đất thấp); gò 1 là các ô 62 ở hàng 17 cột 6-11, mặt gò ở y 408 — cao hơn đáy nhà đúng 1 ô. Đã crop pano0 (170,540)-(420,690) phóng 3.5x.

**Cach sua:** Hoặc dời nhà nhỏ sang trái 3 ô cho nó nằm trọn trên nền thấp, hoặc kéo gò 1 bắt đầu từ cột 9 trở đi. Nguyên tắc chung nên áp cho cả map: mọi vật thể phải nằm trọn trong một dải nền cùng cao độ, không được vắt qua bậc.

<sub>Lang kinh `choi-duoc` · phan bien: XÁC NHẬN, giữ nguyên mức vừa. Crop pano0.png (80,540)-(420,700) phóng 3.5x: từ screen x=240 trở sang phải, mặt cỏ của gò 1 cắt ngang 48 px dưới cùng của ngôi nhà — cửa lùa, chân tường và cả viên đá kê bị đường cỏ xén ngang theo một đường thẳng tuyệt đối, trong khi nửa trái nhà (x 96-240) vẫn đứng trên nền thấp y=626. Nhà đúng là có hai cao độ chân. Dữ liệu khớp chính xác: tile 70 (288x288) hàng 12 cột 3 -> x 72-216 game = screen 96-384; gò 1 là tile 62 hàng 17 cột 6-11 -> bắt đầu x=144 game = screen 240. Tôi khớp mẫu ảnh 1$70.png lên pano0 và ra đúng offset (0,0) nên vị trí không có nghi ngờ.</sub>


## Nhe

### 30. Tile 65 có 9 cột màu chết ở mép trái, thành vạch đen ở góc dưới bên trái mỗi gò

**O dau:** pano0.png x 240–248 · pano2.png x 766–774 · pano4.png x 446–454, đều ở y 626–673 (mép TRÁI của hai gò, cột 6 và cột 30, hàng 18). Nguồn: D:/NRO_Project/Assets/Resources/res/x2/t/1$65.png

Cùng loại lỗi với tile 66 nhưng ở mép bên kia: một vạch dọc rộng 9px màu (54,40,28) hoàn toàn phẳng, không vân, không hạt, nằm ngay góc dưới bên trái mỗi gò. Ở xa thì trông như bóng đổ nên đỡ chối hơn cái mã vạch, đứng gần thì rõ là một mảng màu chết.

**Bang chung:** Đo pano2 tại y=650: x 766 → 774 đều ≈ (54,40,28) không đổi, trong khi x 765 = (75,44,25) và x 775 = (63,46,31). Lặp lại y hệt ở pano0 x240–248 và pano4 x446–454. Mở 1$65.png: cột 0–7 phẳng một màu (52–54, 39–41, 26–28) trên toàn bộ 48 hàng — tức ảnh tile chỉ vẽ thật 40/48 cột.

**Cach sua:** Vẽ tiếp vân đá vào 8 cột trái của 1$65.png. Sửa chung một lượt với 1$66.png vì hai ô là một cặp.

<sub>Lang kinh `dia-hinh` · phan bien: Số liệu đúng từng pixel: pano2 y=650 x766..774 = (52,39,27)…(50,39,26) gần như bất biến, x765=(75,44,25), x775=(63,46,31) — trùng khít báo cáo, và pano0 x240-248, pano4 x446-454 cho dãy số Y HỆT. Trong 1$65.png, cột 0-7 có std chỉ 0.70/0.64/0.49, biên độ 52-55/38-42/26-29 trên cả 48 hàng → đúng là một dải phẳng. Phóng 8x pano2 x750-840 y600-690: có thấy một dải dọc tối, phẳng, mép cứng ở góc dưới trái gò. NHƯNG BÁC BỎ phần lập luận: đây KHÔNG cùng loại lỗi với tile 66. Tile 66 là smear thật (7 cột giống hệt nhau từng byte); tile 65 thì các cột vẫn khác nhau ±3 mức, tức có nhiễu vẽ/render, và nó nằm đúng vị trí một bóng đổ ở mép cắt bên trái — nhiều khả năng là chủ ý của người vẽ chứ không phải bôi mép. Vì vậy giữ mức 'nhẹ' và coi là điểm thẩm mỹ đáng bàn, không phải lỗi kỹ thuật chắc chắn.</sub>

### 31. Cả 4 chỗ đổi độ cao đều là vách dựng đứng 90 độ, không có một đoạn dốc nào

**O dau:** pano0.png x=240 và x=528 · pano2.png x=766 và x=1150 · pano4.png x=446 và x=830 (biên cột 5|6, 11|12, 29|30, 37|38), dải y 578–673

Mặt đất chỉ có đúng hai độ cao và mọi lần chuyển giữa chúng là một bậc dựng đứng 48px — không dốc, không bậc đá, không cỏ rủ xuống mép. So với concept đã duyệt (bờ đất bậc thang, cây cỏ tràn qua từng mép) thì địa hình trong game cứng và 'ô vuông' hơn hẳn. Đáng nói là bộ tile CÓ sẵn hai ô dốc chéo vẽ sạch mà map không dùng ô nào.

**Bang chung:** Quét đỉnh mặt đất theo từng cột trên cả 3 ảnh pano: chỉ ra đúng 2 giá trị — y=578 (mặt gò) và y=626 (mặt nền) — và chuyển giữa chúng luôn gọn trong 1 pixel ngang, không có giá trị trung gian nào. Đọc 1.bytes: không có ô 63, 64, 67 hay 68 (các ô dốc/góc) ở bất kỳ ô nào trong 1248 ô. Mở ảnh tile kiểm tra: 1$67.png và 1$68.png là hai ô dốc chéo vẽ sạch, viền không lỗi; chỉ 1$63.png mới hỏng (viền rác kẻ sọc xám-lục ở cả 4 cạnh, trung bình cạnh 130–180 so với ruột 94,71,38).

**Cach sua:** Dùng 67/68 làm dốc ở ít nhất một đầu mỗi gò (đầu người chơi hay đi lên), hoặc thêm bậc đá/cầu thang. Nếu quyết giữ vách đứng thì phải rải cỏ rủ + đá lẻ đè lên mép cho mềm. Riêng 1$63.png thì cắt bỏ viền rác trước khi dùng.

<sub>Lang kinh `dia-hinh` · phan bien: Quét đỉnh mặt đất từng cột trên cả 3 ảnh: chỉ ra 2 giá trị y=578 và y=626, chuyển gọn trong 1 pixel (mép gò pano0 tại x=528, pano2 tại x=1150, pano4 tại x=830 — đo bằng dải nâu ở y=580: pano2 765-1149, pano4 445-829, đúng 384px = 8 ô). Đọc hết 1248 ô trong 1.bytes: chỉ dùng 60,61,62,65,66,69,70-81 — KHÔNG có 63, 64, 67, 68. Đúng. Dán 10 tile cạnh nhau xem tận mắt: 1$67 và 1$68 là hai ô dốc chéo cỏ/đất vẽ sạch, viền không lỗi; 1$63 thì viền rác kẻ sọc xám ở cả 4 cạnh (đo được cạnh trên 152, dưới 145, trái 130, phải 180 so với ruột 90,67,30 — khớp báo cáo). Concept đã duyệt cũng đúng là bờ đá bậc thang có cây cỏ tràn mép. NHƯNG hạ mức xuống 'nhẹ' và lưu ý: (a) vách đứng 48px là chuẩn mực bình thường của dòng game này, đây là ý kiến thẩm mỹ chứ không phải lỗi; (b) nó trùng phần lớn với lỗi ĐÃ BIẾT số 2 trong bối cảnh ('tile dốc 63/64 không khớp, đang né bằng 62'). Phần thật sự mới chỉ là: 67/68 sạch mà bị bỏ không dùng, và 63 có viền rác.</sub>

### 32. Gradient trời bị đảo ngược: dải 33px trên đỉnh màn hình sáng hơn hẳn phần trời ngay dưới nó

**O dau:** cả 3 ảnh — dải y=0..33 so với dải y=38..343, đo trên x 350-1453

Trời thật thì đậm nhất ở đỉnh đầu và nhạt dần xuống chân trời. Ở đây ngược lại: 33 hàng trên cùng là trời nhạt trắng, rồi tụt cái rụp xuống xanh đậm bão hoà cho tới y=343, rồi lại nhạt trở lại từ y=562. Tông sáng đi theo kiểu nhạt–đậm–nhạt nên dải trên cùng đọc ra như một thanh băng dán ngang đầu màn hình chứ không phải bầu trời. Đây là hệ quả của việc xếp các lớp nền như những băng ngang rời nhau, không phải là chuyển tiếp parallax.

**Bang chung:** pano4, trung bình theo hàng trên x 350-1453: hàng 5 = (96, 167, 209), hàng 20 = (92, 165, 206) — nhạt; hàng 45 = (40, 109, 159), hàng 60 = (39, 109, 162) — đậm. Chênh khoảng 55 đơn vị ở kênh Red và 47 ở Blue chỉ qua 4 hàng. pano2 tại x 1200-1300: hàng 33 = (128, 185, 186) → hàng 38 = (74, 128, 140). Ảnh cắt 1:1: p4_top1x.png.

**Cach sua:** Sau khi xử lý alpha/chồng lớp ở lỗi số 1, sắp lại thứ tự và vị trí dọc để tông trời chạy liên tục từ đậm ở đỉnh xuống nhạt ở chân trời; nếu vẫn dùng lớp trời riêng cho dải trên cùng thì chỉnh màu lớp đó tối lại cho khớp mép nối (đưa Red về khoảng 40-50 tại chỗ giáp y=34).

<sub>Lang kinh `lop-nen` · phan bien: Đo được thật. pano4 trung bình x350-1454: hàng 32 = (97.9,169.9,206.7) và hàng 33 = (97.2,170.1,207.0) — sáng; hàng 40 = (42.1,109.3,158.0), hàng 45 = (40.0,108.6,158.9) — đậm hẳn. Chênh ~57 ở kênh Red chỉ qua vài hàng. pano0: hàng 33 = (111.9,175.5,196.0) → hàng 45 = (67.5,123.8,152.8). Đúng là nhạt–đậm–nhạt, ngược quy luật trời thật. Crop 1:1 (v_top1x.png) cho thấy dải trên còn bị cắt cụt cả đám mây ở y=34, đọc ra như băng dán. NHƯNG TÔI HẠ MẠNH TỪ 'VỪA' XUỐNG 'NHẸ': đây KHÔNG phải một lỗi độc lập mà là mặt còn lại của đúng cái mối nối y=34 đã tính ở phát hiện số 3 — cùng một đường, cùng một nguyên nhân, sửa một lần là hết cả hai. Tính thành hai đầu việc là đếm trùng. Ngoài ra vùng ảnh hưởng chỉ 33 hàng sát mép trên, phần trái bị HUD và loé sáng che, người chơi ít nhìn tới. Ghi nhận để sửa kèm, không nên xếp thành mục riêng cùng hạng với lỗi số 1-2.</sub>

### 33. Biển hiệu cắm bên trong quán ramen — quán đã có sẵn 4 tấm biển của nó

**O dau:** pano4.png vùng x 1167–1263, y 482–626. Ảnh soi: scratchpad/z_bienhieu.png (phóng 4x)

Biển hiệu (cột gỗ + bảng mũi tên + lồng đèn nhỏ) đặt ở cột 45–46, nằm trọn trong chân quán ramen (cột 40–47). Trong game: tấm bảng mũi tên có hoa văn xoáy đỏ đè lên mép mái ngói của quán, cái lồng đèn của biển treo ngay cạnh lồng đèn sẵn có của quán, cây cột thì lẫn vào cột hiên. Mà bản thân sprite quán (1$73.png) đã có biển tô mì lớn dưới mái, một biển treo bên hông, một bảng thực đơn dựng đất và hai chậu cây — nhét thêm biển hiệu vào chỉ thành mớ rối. Bụi cỏ ở cột 44–45 cũng rơi vào trong quán, đè lên quầy hàng, làm nặng thêm.

**Bang chung:** 1.bytes: slot 79 hàng 15 cột 45 (ảnh 1$79.png 96x144 = 2x3 ô -> cột 45–46, hàng 15–17); slot 73 hàng 12 cột 40 (ảnh 384x288 = 8x6 ô -> cột 40–47). Bụi cỏ slot 81 hàng 17 cột 44 -> cột 44–45. Đã render riêng 1$73.png (scratchpad/sprite_nha_quan.png) để xác nhận quán tự có sẵn 4 biển + 2 chậu cây.

**Cach sua:** Dời biển hiệu ra cột 38–39 (khoảng hở giữa nhà lớn và quán) hoặc cột 24–25 ở đầu lối lát đá — chỗ biển chỉ đường mới có nghĩa. Dời bụi cỏ cột 44 ra cột 48–49.

<sub>Lang kinh `vat-the` · phan bien: Hình học ĐÚNG, kết luận thị giác thổi phồng. 1.bytes: slot 79 (96x144) ở (15,45) -> cột 45-46, nằm trong quán slot 73 cột 40-47; bụi cỏ slot 81 ở (17,44) -> cột 44-45 cũng trong quán. Crop 3x pano4 (1030,400)-(1330,650) [v5_quan.png]: biển hiệu ở màn x 1183-1257, đúng là đứng trước quán, tấm bảng mũi tên có chạm mép mái hiên, lồng đèn của biển nằm cạnh lồng đèn của quán (2 lồng đèn đỏ gần y hệt sát nhau) — chỗ này hơi rối thật. Nhưng KHÔNG có vật nào bị cắt, bị hẫng hay bị chôn; nhìn tổng thể nó đọc ra là một quán mì đông đúc, hợp thể loại. Bụi cỏ mọc trước quầy/ghế là chi tiết vô lý nhẹ. Hạ 'vua' -> 'nhe' vì đây là ý kiến thẩm mỹ chứ không phải lỗi dựng cảnh.</sub>

### 34. Cây bị cắt dọc ở mép phải map, khoảng 1/3 tán và thân không bao giờ nhìn thấy được

**O dau:** pano4.png vùng x 1263–1454 (ảnh dừng ở 1454, cây còn kéo dài tới x 1551). Ảnh soi: scratchpad/z_cay.png

Cây rộng 6 ô đặt ở cột 47 nên chiếm cột 47 tới 52. Nhưng lưới map chỉ có 52 cột đánh số 0–51, tức cột 52 KHÔNG tồn tại. Thêm nữa camera đẩy hết cỡ sang phải cũng chỉ tới đơn vị game 1223,5 (đo được từ pano4). Nên khoảng 97px bên phải của sprite cây — gồm rìa tán và một phần thân — vĩnh viễn nằm ngoài màn hình. Người chơi đi tới cuối map thấy một cái cây bị dao xẻ dọc, mặt cắt thẳng đứng ngay mép màn hình. Đây là điểm dừng chân cuối map nên ai chơi cũng gặp.

**Bang chung:** 1.bytes: slot 80 hàng 10 cột 47; ảnh 1$80.png 288x384 = 6x8 ô -> cột 47–52 = đơn vị game 1128–1272, trong khi map rộng 52*24 = 1248. Đo camera pano4 bằng so khớp mẫu vùng đất với pano2 (sai số trung bình 0,21 mức xám -> khớp tuyệt đối): camera trái = 496,5, camera phải = 496,5+727 = 1223,5. Sprite cây kéo tới 1272 -> hụt 48,5 đơn vị game = 97px màn.

**Cach sua:** Dời cây về cột 45 (chiếm cột 45–50, nằm gọn trong vùng camera tới được) — kèm dời biển hiệu đi chỗ khác để khỏi đụng. Hoặc thêm 2 cột đất nữa ở cuối map.

<sub>Lang kinh `vat-the` · phan bien: Tự kiểm và ĐÚNG, kể cả điểm khó nhất là 'camera không tới được'. 1.bytes: slot 80 ở (10,47), 1$80.png 288x384 -> cột 47-52 = game x 1128-1272. camLeft(pano4) tôi đo độc lập bằng bậc gò cột 38 (game x=912 rơi đúng màn x=830) -> camLeft=497, mép phải = 497+727 = 1224. Bằng chứng pano4 ĐÚNG là vị trí cuộn tối đa: HUD ghi 'X: 1238' tức nhân vật đã ở game x=1238 mà vẫn không lọt khung, và dải đất cuối (cột 50, hết ở game x=1224) kết thúc đúng bằng mép ảnh x=1454. Vậy 1272-1224 = 48 đơn vị = 96px sprite cây vĩnh viễn ngoài màn. Đo alpha 1$80.png dải x>=192: 44,6% pixel đặc -> phần bị cắt là nội dung thật. Crop [v6_cay.png] thấy thân và tán bị xẻ thẳng đứng ở mép khung. Hạ 'vua' -> 'nhe': cắt cảnh ở rìa map là chuyện thường của game cuộn ngang, và cả dải đất cũng đứt ngang ở đó.</sub>

### 35. Hai cái giếng gần như y hệt nhau hiện cùng lúc trên một màn hình

**O dau:** pano2.png: giếng thứ nhất ở x 0–107 (vẽ liền trong sprite nhà vừa), giếng thứ hai ở x 671–815 (vật thể rời slot 77). Cả hai cùng khung hình. Ảnh soi: scratchpad/sprite_nha_vua.png và z_gieng.png

Sprite nhà vừa (1$71.png) đã có sẵn một cái giếng vẽ liền bên hông trái: bệ đá xếp, mái gỗ lợp ván nhọn, dây thừng và gàu gỗ treo. Vật thể giếng rời (1$77.png) đặt ở cột 28–30 dùng đúng cùng ngôn ngữ tạo hình: bệ đá xếp, mái gỗ lợp ván nhọn, gàu gỗ treo. Hai cái cách nhau 13 ô nên lọt chung một khung hình — mắt đọc ngay ra là copy-paste, làm ngôi làng trông nghèo chi tiết.

**Bang chung:** Render riêng 1$71.png (scratchpad/sprite_nha_vua.png) thấy rõ giếng nằm ở 1/3 trái sprite = cột 13–15 trong map. slot 77 ở hàng 15 cột 28 (1$77.png 144x144 = 3x3 ô -> cột 28–30). Quy chiếu ra màn hình pano2 (camera trái 336,5): giếng-trong-nhà ở x -49..111, giếng-rời ở x 671..815 — cùng nằm trong khung 0..1454.

**Cach sua:** Chọn một trong hai: hoặc bỏ vật thể giếng rời ở cột 28 (thay bằng thùng gỗ / bụi cây), hoặc vẽ lại giếng rời khác kiểu hẳn (giếng tròn không mái, hoặc bể nước đá). Đừng để hai cái cùng khung.

<sub>Lang kinh `vat-the` · phan bien: ĐÚNG, tôi tự nhìn cạnh nhau. Cắt 2 vùng pano0: giếng vẽ liền trong sprite nhà vừa ở màn x 577-721 [v7a_well_in_house.png] và giếng rời slot 77 (cột 28-30) ở x 1297-1441 [v7b_well_free.png] — cả hai cùng nằm trong 1 khung 0-1454, và cùng công thức tạo hình: mái gỗ lợp ván dốc, 2 cột vuông, thanh ngang, dây + gàu gỗ, bệ đá xếp khối. Nhìn liền hai ảnh đúng là ra cảm giác lặp. Nhưng khác nhau đủ để không phải copy-paste thật (góc mái, kiểu xếp đá, hướng gàu khác nhau) và một làng có hai giếng không vô lý. Hạ 'vua' -> 'nhe': đây là chuyện đánh giá gu, không phải lỗi.</sub>

### 36. Giếng bị bậc gò cắt cụt cột chống bên phải

**O dau:** pano2.png vùng x 767–815, y 578–626 (cột 30 của sprite giếng). Ảnh soi: scratchpad/z_gieng.png (phóng 4x)

Giếng rộng 3 ô đặt ở cột 28 nên chiếm cột 28–30, đáy ở hàng 18. Cột 30 lại đã là gò cao (mặt đất hàng 17). Nên phần dưới cột chống phải và mép phải bệ đá của giếng bị vạt cỏ gò che mất 48px cuối. Nhìn ra là cái giếng bị cụt một chân — cột trái chạy xuống tới bệ, cột phải hụt lửng ở mép đồi.

**Bang chung:** 1.bytes: slot 77 hàng 15 cột 28; 1$77.png 144x144 = 3x3 ô -> cột 28–30. Bảng độ cao: gò chạy từ cột 30. Đo alpha sprite 1$77.png: ô cột thứ 3 (x 96–144) vẫn còn 51% pixel đặc ở 6 dòng cuối — tức phần bị chôn là nội dung thật, không phải vùng trong suốt.

**Cach sua:** Dời giếng sang cột 26 hoặc 27 (chiếm cột 26–28/27–29, lùi khỏi bậc gò), hoặc bắt đầu gò từ cột 31 thay vì 30.

<sub>Lang kinh `vat-the` · phan bien: ĐÚNG, mức 'nhe' cũng đúng. 1.bytes: slot 77 (144x144) ở (15,28) -> cột 28-30, đáy hàng 18 (màn y 626); gò bắt đầu đúng cột 30 (mặt y 578). Crop 4x pano0 (1280,480)-(1454,640) [v7b_well_free.png]: khối đất+cỏ bắt đầu ở màn x=1393 (đúng biên cột 30, game x=720) và che phần dưới cột chống phải cùng mép phải bệ đá từ y=578 xuống. Nhìn ra được là chân phải giếng hụt so với chân trái, nhưng vì bệ đá vốn thấp nên nó vẫn đọc tạm được là giếng đứng nép mô đất.</sub>

### 37. Ô đất 66 (vách phải) lòi ra dải sọc trắng xám như cái thang cắm giữa lớp đất

**O dau:** pano2.png x 1139–1151, y 626–674 (cột 37 hàng 18) và pano0.png x 517–529, y 626–674 (cột 11 hàng 18). Ảnh soi: scratchpad/z_chan_lo_lung.png góc dưới trái, và tiles.png (bảng ảnh 10 viên gạch)

Ảnh viên gạch 1$66.png bị dính một dải viền sọc sáng-tối rộng ~4px ở mép phải (viền thừa lúc cắt ảnh, chưa dọn). Viên này được đặt ở 2 chỗ trong map (cột 11 và cột 37, hàng 18) làm mặt cắt vách đất. Trong game nó hiện ra thành một dải sọc trắng xám thẳng đứng như cái thang gỗ / cái lược cắm giữa khối đất nâu, chỗ nào cũng nhìn thấy. Ngoài ra thân viên 66 sáng và phẳng hơn hẳn đất xung quanh nên còn thêm một mảng vá màu lạ. Viên 1$65.png (vách trái) dính lỗi tương tự ở mép trái, và 1$63.png (dốc lên) dính viền sọc cả 4 phía. Ghi chú: đây có thể thuộc phần soi mặt đất chứ không phải phần vật thể, nhưng nó lộ ngay cạnh chỗ chân nhà lớn treo lơ lửng nên báo kèm.

**Bang chung:** Dựng bảng ảnh 10 viên tile 60–69 (scratchpad/tiles.png) thấy rõ viền bẩn trên 63, 65, 66. Giải mã 1.bytes: tile 65 ở (hàng 18, cột 6) và (18, 30); tile 66 ở (18, 11) và (18, 37). Đối chiếu ảnh game phóng 8x tại pano2 x 1120–1220 y 540–640: dải sọc xám hiện đúng vị trí mép phải cột 37.

**Cach sua:** Chạy tools/map-lang-la/do-vien-tile.py rồi xoa-vien-tile.py cho 1$63, 1$65, 1$66 (script có sao lưu .bak-vien). Sau đó chỉnh lại tông màu thân viên 65/66 cho khớp với 1$60/1$61 để hết mảng vá sáng.

<sub>Lang kinh `vat-the` · phan bien: Hiện tượng chính ĐÚNG, hai chi tiết phụ SAI. Đúng: đếm ô trong 1.bytes cho thấy tile 66 đặt ở (18,11) và (18,37), tile 65 ở (18,6) và (18,30). Crop 7x pano2 (1090,610)-(1210,690) [v9b_t66.png] và pano0 (470,610)-(590,690) [v9c_t66_p0.png]: cả hai chỗ đều lòi ra dải sọc sáng-tối dọc rộng ~6px ở mép phải (màn x 1143-1149 và x 522-529, y 626-673), trông đúng như cái thang/cái lược cắm trong đất. SAI 1: tile 63 KHÔNG hề được đặt ở đâu trong map (đếm được 0 ô) — viền bẩn của nó không ảnh hưởng gì trong game, đưa vào báo cáo là thừa. SAI 2: người báo nói thân viên 65/66 'sáng và phẳng hơn hẳn đất xung quanh' — đo độ sáng trung bình cho kết quả ngược: 65 = 59, 66 = 58, trong khi 60 = 67, 61 = 62, tức nó TỐI hơn. Vẫn là mảng vá lệch tông, nhưng lệch theo hướng tối.</sub>

### 38. Phía sau cổng torii — điểm nhấn của map — chỉ có trời trơn, concept là vách núi + nhà xa nhiều lớp

**O dau:** pano2.png vùng x 250–620, y 350–600 (cổng torii ở x 250–470). Đối chiếu concept-map1-langla.png vùng x 480–900, y 420–830.

Cổng torii là vật thể to nhất và bắt mắt nhất map, nhưng trong game nó đứng trước một nền trời xanh trống trơn: qua khung torii nhìn thấy mây, hai bên torii cũng chỉ có trời. Trong concept, ngay sau torii là vách đá núi Hokage lấp đầy khung, cộng thêm 2–3 lớp nhà xa mờ dần hai bên và một tán cây xanh bên trái. Vì mất hẳn lớp nền sau lưng, torii mất cảm giác "cổng dẫn vào làng" và trông như một mô hình đặt trên phông xanh.

**Bang chung:** Crop phóng 3x vùng ghép panorama (860,380)-(1180,680) và crop concept (480,420)-(900,830), xem cạnh nhau. Trong crop game, toàn bộ vùng bên trong và bên phải khung torii là màu trời (44,148,217)→(62,161,223) không có chi tiết nào ngoài mây. Trong crop concept cùng vị trí là đá núi + 4 mái nhà xa + tán cây.

**Cach sua:** Sau khi phóng to núi (phát hiện 1) và hạ lớp làng xa (phát hiện 2), kiểm tra lại riêng khung nhìn quanh torii — đây là chỗ người chơi dừng lại lâu nhất. Cân nhắc đặt thêm 1 cây và 1 cụm nhà xa ngay hai bên torii.

<sub>Lang kinh `doi-chieu` · phan bien: Sự kiện thì đúng: tôi crop pano2 (240,360)-(660,660) phóng 2x, toàn bộ khung torii và hai bên chỉ có trời xanh + mây, cộng đúng cái dải làng xa bị chặt ở phát hiện 2. Concept cùng vị trí là vách đá núi + nhà xa + tán cây. NHƯNG đây không phải một lỗi độc lập — nó chỉ là hệ quả nhìn từ góc khác của phát hiện 1 (núi quá nhỏ/quá cao) và phát hiện 2 (làng xa bị cắt). Chính người báo cáo cũng viết cách sửa là 'sau khi sửa 1 và 2 thì kiểm tra lại'. Đếm nó thành một mục nặng riêng là nhân đôi cùng một khuyết điểm. Hạ nặng → nhẹ (giữ lại như một điểm cần kiểm tra sau khi sửa 1+2, không phải lỗi mới).</sub>

### 39. Hai cái giếng giống hệt nhau cùng lộ trong một màn hình; concept chỉ có một

**O dau:** pano0.png x 585–670 và x 1300–1425, y 490–620 (cả hai cùng xuất hiện trong MỘT ảnh chụp). Đối chiếu concept-map1-langla.png: giếng duy nhất tại x 830–950, y 600–780.

Trong cùng một khung hình pano0 có hai cái giếng dùng chung một sprite, cùng kích thước, cùng góc nhìn, cách nhau 731px. Người chơi lia camera sẽ thấy ngay là bản sao. Concept chỉ có một cái giếng, đặt làm điểm nhấn cạnh torii. Danh sách vật thể dự kiến của map cũng chỉ ghi "1 giếng".

**Bang chung:** Read pano0.png: giếng thứ nhất ở x 585–670 (bị nhà quán che một phần), giếng thứ hai ở x 1300–1425. Cắt và phóng 3x cả hai (z_gieng1 / z_gieng2) — cùng mái gỗ, cùng thân đá, cùng gàu. Dò offset khớp nhất giữa hai vùng: dx=731, dy=−6.

**Cach sua:** Bỏ một cái, hoặc thay bằng vật thể khác (sạp hàng, đống củi, bảng tin làng) để không lộ lặp sprite trong cùng một màn.

<sub>Lang kinh `doi-chieu` · phan bien: Đúng. Tôi cắt hai vùng pano0 (575,480)-(690,640) và (1295,480)-(1440,640), phóng 3x, ghép cạnh nhau và nhìn: cùng một sprite — cùng mái gỗ xếp thanh, cùng thanh quay ngang có dây buộc, cùng gàu gỗ đai kim loại, cùng thân đá xếp khối, cùng góc nhìn. Cái trái bị nhà quán che một phần nên khớp pixel chỉ ra sai khác trung bình 44 (do nền và phần bị che khác nhau), nhưng bằng mắt thì rõ ràng là bản sao. Kiểm chéo bằng panorama: cái phải của pano0 (x1300) chính là cái giếng thấy ở pano4 x≈370 (lệch camera 946px), còn cái trái x585 là vật thể khác — vậy đúng 2 giếng thật, không phải nền lặp. Danh sách 18 vật thể của map chỉ ghi 1 giếng. Giữ mức nhẹ: hai giếng cách nhau ~715px, người chơi có thể thấy trùng nhưng không phá hỏng gì.</sub>

### 40. Map không có cửa/waypoint nào — không có chỗ nào để đi tới

**O dau:** server/internal/game/thegioi.go:126 — `w.Byte(0) // waypointCount = 0`. Biểu hiện trên ảnh: lối lát đá pano2 x 285-623, y 626-674, bắt đầu dưới cổng torii rồi cụt giữa bãi cỏ.

Gói MAP INFO gửi cho client khai báo cứng số waypoint = 0, nghĩa là không có một cái cửa nào để rời map. Cộng với thứ duy nhất trên map gợi ý 'đường đi' là dải lát đá dài 7 ô: nó bắt đầu ngay dưới cổng torii — đúng chỗ người chơi kỳ vọng 'qua cổng là đi đâu đó' — rồi dừng đột ngột giữa bãi cỏ, không dẫn tới đâu. Với người chơi, cả màn chơi là một hành lang thẳng dài 1194 đơn vị, đi hết rồi quay lại, không có lối ra, không có đích.

**Bang chung:** grep waypoint trong server/internal/game/thegioi.go chỉ có đúng một dòng, giá trị 0. res/mymap/1.bytes: tile 69 (lát đá) chỉ ở hàng 18 cột 20-26, hai đầu cắt vuông góc không có ô chuyển tiếp — đã crop pano2 (250,605)-(700,700) phóng 3x để xem: đúng là một hình chữ nhật màu sáng dán lên nền đất, hai mép là hai đường thẳng đứng.

**Cach sua:** Thêm ít nhất một waypoint ở một đầu map (ví dụ mép phải, dưới cây cổ thụ) trỏ sang map kế tiếp; và cho lối lát đá kết thúc ở waypoint đó thay vì cụt giữa bãi cỏ. Nếu chưa có map 2 để đi, ít nhất đặt lối đá dẫn tới cổng torii rồi dừng ở đó cho có lý.

<sub>Lang kinh `choi-duoc` · phan bien: Dữ kiện ĐÚNG nhưng tôi HẠ mức từ 'vừa' xuống 'nhẹ'. Kiểm code: server/internal/game/thegioi.go:126 ghi `w.Byte(0) // waypointCount = 0`, chỉ đúng một chỗ, hằng số cứng — xác nhận. Lối lát đá cũng đúng: 1.bytes có tile 69 tại hàng 18 cột 20-26 -> game x 480-648 -> screen pano2 x 286-622, y 626-674 (khớp con số 285-623 báo cáo đưa). Crop pano2 (240,600)-(700,700) phóng 3x: đúng là một hình chữ nhật lát đá sáng màu, hai đầu cắt vuông góc, không có ô chuyển tiếp, và bắt đầu ngay dưới cổng torii rồi cụt giữa cỏ. Lý do hạ mức: waypoint = 0 KHÔNG phải lỗi mà là việc chưa làm — dự án mới có map 1, chưa có map 2 để trỏ sang; gán nó là 'lỗi map' là thổi phồng. Phần thật sự đáng sửa chỉ là hai đầu lối đá cắt vuông góc, và đó là chuyện thẩm mỹ nhỏ.</sub>

### 41. Cột đá chân nhà lớn treo lơ lửng ngoài mép gò

**O dau:** pano2.png vùng x 1153-1168, y 562-578 (mép phải gò 2 kết thúc ở x 1150).

Nhà lớn được đặt sao cho góc phải của nó nhô ra khỏi mép gò. Viên đá kê chân cột ở góc phải (khối xám khoảng 15x16 pixel) nằm hoàn toàn bên ngoài vách đất, bên dưới nó là bầu trời. Nhìn kỹ vào góc nhà là thấy một cục đá bay lơ lửng.

**Bang chung:** res/mymap/1.bytes: nhà lớn tile 72 (ảnh 384x384) ở hàng 9 cột 31, vùng đặc x 744-936 đơn vị game; gò 2 kết thúc ở cột 37 tức x 912. Phần nhà từ 912 đến 936 không có gò đỡ. Đã crop pano2 (1100,545)-(1210,690) phóng 7x: thấy rõ cột gỗ + viên đá kê nằm ngoài vách, phía dưới là trời và bụi cỏ ở nền thấp.

**Cach sua:** Kéo dài gò 2 thêm 1 ô (tới cột 38), hoặc dời nhà lớn sang trái 1 ô (về cột 30).

<sub>Lang kinh `choi-duoc` · phan bien: XÁC NHẬN, giữ nguyên mức nhẹ. Crop pano2.png (1080,530)-(1230,700) phóng 5x: nhìn rõ khối đá kê màu xám ở góc phải nhà lớn nằm hẳn ngoài vách gò, BÊN DƯỚI nó là trời xanh rồi mới tới bụi cỏ ở nền thấp. Dữ liệu khớp: tile 72 (384x384) hàng 9 cột 31 -> x 744-936, đáy y=408; gò 2 (tile 62 hàng 17) kết thúc ở cột 37 -> x=912. Đoạn x 912-936 (screen 1150-1198) không có gò đỡ, treo cao hơn mặt cỏ thấp (y=432) đúng 24 đơn vị = 48 px. Khớp mẫu 1$72.png lên pano2 ra offset (0,0) nên vị trí chắc chắn.</sub>

### 42. HUD còn dòng gỡ lỗi X/Y, hai thanh máu/chakra đều bằng 0, và trên màn hình không có nút nào ngoài bánh răng

**O dau:** pano2.png: hộp nhân vật x 3-157 y 5-63; thanh máu x 163-300 y 10-33; thanh chakra x 163-293 y 35-55; dòng gỡ lỗi "1 - Khu: 0 / X: 700 - Y: 576" ở x 177-290 y 65-100; bánh răng x 320-373 y 5-53.

Về câu hỏi 'HUD có che mất phần quan trọng của map không' thì câu trả lời là KHÔNG — toàn bộ HUD nằm ở góc trên trái, chỗ đó chỉ có bầu trời, dải chơi thật (y 370-680) hoàn toàn không bị che. Nhưng ba thứ khác đập vào mắt: (1) dòng toạ độ gỡ lỗi vẫn đang bật, ở pano0 và sach2A nó còn vỡ font thành các khối vuông không đọc được, ở pano2/pano4 thì đọc được; (2) cả thanh máu lẫn thanh chakra đều hiển thị số 0 và rỗng hoàn toàn, nhìn như nhân vật đã chết; (3) ngoài cái bánh răng cài đặt thì không có nút đánh, nút kỹ năng, nút menu hay khung chat nào — người chơi mở game ra không có gì để bấm.

**Bang chung:** Crop pano0 (0,0)-(380,120) phóng 3x: chữ 'narruto', số 0 trong hộp xanh, số 0 đỏ trên thanh máu, số 0 xanh lơ trên thanh chakra, hai dòng chữ trắng vỡ nét. So với server/internal/game/lang.go:27 thì Làng Lá có HpKhoiDau 150 và ChakraKhoiDau 150 — tức HUD đang không nhận được chỉ số, chứ không phải nhân vật thật sự có 0 máu.

**Cach sua:** Tắt dòng gỡ lỗi X/Y trong bản cho người chơi xem (giữ sau một cờ debug). Kiểm tra gói gửi HP/chakra xuống client vì giá trị khởi đầu bên server là 150 chứ không phải 0. Phần thiếu nút điều khiển thì ghi vào việc còn phải làm, không phải lỗi map.

<sub>Lang kinh `choi-duoc` · phan bien: XÁC NHẬN cả ba ý, giữ mức nhẹ. Crop pano0 (0,0)-(380,120) phóng 3x nhìn tận mắt: (1) dòng 'X: 60 - Y: 576' hiện rõ ở x 195-284 y 84-100, ở pano0/sach2A thì font vỡ nham nhở, ở pano2/pano4 thì đọc được ('X: 700 - Y: 576' / 'X: 1238 - Y: 576'); (2) thanh máu và thanh chakra đều hiển thị số 0 và lòng thanh rỗng hoàn toàn màu nâu sẫm, đúng như mô tả; (3) ngoài bánh răng ở x 320-373 không có nút nào khác trên toàn khung. Ý 'HUD không che map' cũng đúng — HUD nằm gọn ở góc trên trái trên nền trời, dải chơi y 570-680 không bị che. Giữ mức nhẹ vì đây là HUD chứ không phải map, và câu hỏi của chủ dự án là về map.</sub>


---

## Tong ket

## Trả lời ngắn: map **chưa hoàn thiện** — đẹp nhưng còn 3 thứ không thể cho người chơi thấy

Tôi đã xem tận mắt cả 3 ảnh chụp game. Ấn tượng chung: nhà cửa, cổng torii, quán mì vẽ rất đẹp, bố cục làng đọc ra ngay. Nhưng có mấy thứ đập vào mắt trước cả cái đẹp.

---

### Ba lỗi CHẶN ĐƯỜNG (phải sửa trước khi cho ai xem)

**1. Cây cột tre lơ lửng giữa trời, ngay chính giữa màn hình.**
Nó có ở cả 3 ảnh (pano0 x≈607, pano2 x≈587, pano4 x≈577), chạy từ mép trên màn hình xuống tận chân núi rồi cụt ngang. Không chân, không gốc, không dính vào gì. Phóng to thấy rõ là sprite tre có nút thừng quấn, lặp lại theo chiều dọc. Đã quét toàn bộ 4 file nền và ~1700 ảnh trong Resources — **không khớp file nào**, tức nó là rác được vẽ đè từ đâu đó (gần chắc sót lại từ map Ngọc Rồng cũ). Đây là thứ người chơi chụp màn hình đem đi chế đầu tiên.

**2. Cả bầu trời bị một nhát dao chém ngang ở y=344, chạy suốt màn hình.**
Núi Hokage đang có vân đá thì đột ngột hết, dưới là trời xanh trơn. Nguyên nhân: 4 lớp nền đều là PNG **không có kênh trong suốt**, nên mỗi lớp là một hình chữ nhật đặc, mép trên lớp này xén phẳng lớp kia. Cùng bệnh đó gây thêm một nhát chém thứ hai ở y=561 — chỗ đó là **làng xa**: chỉ lộ ra dải mái ngói cao 13px rồi bị cắt đáy, thành ra một hàng mái nhà treo lơ lửng giữa trời.

**3. Không nhìn thấy nhân vật ở bất kỳ ảnh nào.**
Tôi so pano0 với sach2A từng pixel: nhân vật đi 310 đơn vị mà màn hình **không đổi một pixel nào** ngoài dòng chữ toạ độ. Nhân vật rơi xuyên đất xuống y=576, trong khi mép dưới màn hình chỉ tới y=551 — tức người chơi mở game ra thấy một bức tranh phong cảnh, không thấy mình đâu. Gốc là lỗi va chạm đã biết, nhưng hệ quả thị giác thì mới.

---

### Lỗi NẶNG (thấy ngay, phá cảm giác "hoàn thiện")

- **Vách gò ghép từ 5 chất liệu khác nhau**, có một đường ngang thẳng băng chia đôi vách ở y=626 (độ sáng nhảy +55% trong đúng 1 hàng pixel, kéo dài 384px). Nửa trên là đất rễ, nửa dưới lấy nguyên các ô dùng để lấp lòng đất.
- **Hàng rào cắm xuyên cột cổng torii** — thanh ngang chạy thẳng qua thân cột đỏ và bệ đá. Đây là vật thể bắt mắt nhất map.
- **Góc phải nhà lớn treo lơ lửng**: hòn đá kê chân chìa ra ngoài mép gò, dưới nó là trời xanh. Nhà đặt cột 31–38 nhưng gò chỉ tới cột 37.

### Lỗi VỪA (gộp lại thì rất chối)

- **Gần 1/4 màn hình dưới là đất nâu phẳng lì** (22.7% chiều cao), chỉ có 2 mẫu ô, các mảng sỏi là hình chữ nhật 48px viền sắc lẻm nằm rời rạc như bàn cờ. Riêng một dải cao 48px chạy hết map hoàn toàn trống trơn.
- Có **một đường chân đất giả thứ hai** ở y=673 chạy hết map (chỗ nối hàng 18/19 lệch tông).
- **Hai ô tile hỏng ảnh** (1$65, 1$66) làm mép gò lòi ra một cột sọc xám như mã vạch và một dải màu chết — xuất hiện ở 4 chỗ.
- **Nhà nhỏ bị gò nuốt mất nửa phải** 48px, hòn đá kê chân bên phải biến mất; bụi cỏ cạnh đó bị xẻ dọc mất nửa, gốc còn hở 3px thấy trời.
- **Giếng bị vách gò xén mất một chân**; **cây bị mép map cắt mất 1/3 tán** vĩnh viễn; **cổng torii chồng 2 cột lên nhà vừa**.
- **Không có bóng đổ dưới bất cứ vật gì** — nhà, giếng, cây đều như sticker dán lên nền cỏ kẻ thước.
- **Không có kè đá, không có bậc thang** — concept đã duyệt có cả hai, và có cầu thang dẫn lên torii.
- **Cả map chỉ có 1 cây**, đặt sát biên phải; khoảng giữa map không có mảng xanh nào.

### Lỗi NHẸ (ghi để nhớ)

Hai cái giếng gần giống hệt nhau lọt chung một khung · biển hiệu cắm bên trong quán ramen (quán vốn đã có 4 tấm biển) · 2 hàng rào còn lại dán phẳng lên tường nhà · vạch tối mảnh ở sát đỉnh màn hình do 4 hàng pixel hỏng trong b02/b03 · dải trời trên cùng lệch tông (nhạt–đậm–nhạt, ngược quy luật) · HUD còn dòng debug "X: 700 - Y: 576", thanh máu và chakra đều hiện số 0 · chưa có waypoint nào để rời map.

---

## Điểm hoàn thiện: **45/100**

Lý do ngắn: phần **vẽ vật thể** (nhà, quán, torii, giếng) đạt mức 80 — thật sự đẹp. Nhưng phần **ghép cảnh** kéo tụt xuống: 3 lỗi chặn đường còn nguyên, gần 1/4 màn hình là đất chết, và cứ chỗ nào gò gặp vật thể là có một cái bị cắt hoặc bị treo. Đang ở mức "dùng để test được, chưa cho người ngoài xem" và còn dính lỗi chặn đường nên phải trừ xuống dưới 50.

---

## Việc còn lại chia làm hai nhóm

**Nhóm 1 — sửa bằng code / chỉnh dữ liệu, KHÔNG cần vẽ gì (làm được ngay, chiếm ~60% số lỗi):**

Truy ra và gỡ cây cột tre giữa trời · chuyển 4 lớp nền sang PNG có kênh trong suốt (hoặc hạ mép trên xuống dưới đường chân núi) để hết hai nhát chém ngang y=344 và y=561 · cắt 4 hàng pixel hỏng ở đỉnh b02/b03 · khai loại va chạm cho tile 60–82 trong `tile_set_info.bin` và sửa `map-1-collision.bin` lệch 1 hàng để nhân vật đứng được trên đất → tự khắc thấy nhân vật · dời vị trí trong `1.bytes` cho 8 vật thể đặt sai (nhà nhỏ, nhà lớn, giếng, torii, 3 hàng rào, biển hiệu, cây, mấy bụi cỏ) hoặc nới bảng độ cao gò · bỏ bớt 1 giếng · chạy script `xoa-vien-tile.py` dọn viền rác tile 65/66 · đặt sẵn `tc_canh_la_x4.png` (đã vẽ rồi mà chưa dùng) vào góc trên trái · tắt dòng debug X/Y và sửa gói gửi HP/chakra.

**Nhóm 2 — phải vẽ ảnh mới (cần hoạ sĩ hoặc AI vẽ):**

Bộ tile **vách gò** riêng (thân vách + chân vách) khớp tông với đáy ô cỏ — đây là món quan trọng nhất, chữa luôn cả lỗi vách 5 chất liệu lẫn tile 65/66 · **3–4 biến thể tile đất sâu** + một lớp tối dần theo độ sâu để lấp mảng nâu chết · **bộ tile kè đá + cụm bậc thang đá** dưới torii (concept có, game không có) · **2 ô chuyển tiếp cỏ ↔ lát đá** và vài túm cỏ/đá lẻ rải mép · **sprite bóng đổ elip mờ** cho mọi công trình · **vẽ lại b02 (núi Hokage)** to gấp 3–4 lần và làm bản 768px thật thay vì ghép đôi 384px — chữa luôn lỗi núi lặp 2 lần đã biết · **vẽ lại b01 (làng xa)** cao hơn, đáy tắt dần · thêm 1–2 kiểu cây nữa cho khoảng giữa map.

Thứ tự tôi đề nghị: làm hết Nhóm 1 trước (rẻ, nhanh, gỡ được cả 3 lỗi chặn đường), chụp lại 3 ảnh, rồi mới đặt hàng Nhóm 2 — vì sau khi dời vật thể xong, một số lỗi "vẽ" có thể tự hết hoặc đổi yêu cầu.