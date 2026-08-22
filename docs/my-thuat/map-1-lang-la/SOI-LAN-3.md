# Soi map 1 Lang La — lan 3

Ngay 22/08/2026, sau khi sua 9 viec. **35 loi song sot** (33 loi moi). Diem: **56/100** (vong 1: 45, vong 2: 52).

> Anh chung: `pano0.png` .. `pano4.png`, `v3A.png` (trong D:/NRO_Project/)

---


## Nang

### 1. Cuối bản đồ bên phải: nhân vật đi hẳn ra ngoài màn hình, người chơi mất dấu chính mình

**O dau:** D:\NRO_Project\pano4.png — cả thân nhân vật nằm ngoài khung; chỉ còn mép trái cái bóng ở x 1443–1453, y 620–631. Code: D:\NRO_Project\Assets\Scripts\Game1\GameScr.cs:1588 (cmxLim) và Char.cs:6083–6084 (mỏ neo vẽ thân).

Camera dừng ở giới hạn cmxLim = (52-1)*24 - gW = 1224 - 727 = 497 (đơn vị game), tức chỉ nhìn được tới x = 1224. Nhưng ô đi được cuối cùng (cột 50) kéo tới đúng x = 1224, mà thân nhân vật lại được vẽ SANG PHẢI điểm neo cx (anchor TOP|LEFT khi quay phải). Hệ quả: đi tới cuối bản đồ thì nhân vật trượt hẳn ra khỏi mép phải màn hình — người chơi không còn thấy nhân vật của mình, chỉ còn một mẩu bóng. Bên trái không bị vì thân vẽ về phía phải nên vẫn nằm trong khung.

**Bang chung:** pano3 và pano4 CÙNG camera: khớp mẫu ảnh tile 65 rơi đúng x=110 ở cả hai (và tile 66 ở x=302, tile 65 ở x=638) → cmx = 497 = cmxLim ở cả hai ảnh. Trừ pixel pano3 − pano4: pano3 nhân vật chiếm x 984–1053, y 542–632; pano4 chỉ còn vệt khác ở x 1443–1453, y 620–631, hình cánh trái của ellipse bóng (hàng 620: 2px, 624–626: 10px, 630: 4px) — KHÔNG còn pixel thân nào. Quy đổi: tâm bóng ≈ x 1467 → cx ≈ 1230.

**Cach sua:** Ba hướng, chọn một: (a) nới cmxLim thêm ~1 ô cộng nửa bề rộng sprite (dùng tmw*24 - gW rồi trừ bù), (b) đặt tường vô hình ở cột 50 để nhân vật dừng sớm hơn một ô, hoặc (c) căn lại sprite nhân vật cho cx nằm giữa thân (xem phát hiện về bóng đổ) — cách này chữa luôn hai lỗi.

<sub>Lang kinh `kiem-sua` · LOI MOI · phan bien: CÓ THẬT, tự đo lại còn nặng hơn báo cáo. Trừ pixel pano3 − pano4: vùng khác nhau CHỈ nằm ở x 983-1053 (pano3) và x 1443-1453 (pano4) — camera hai ảnh y hệt, chỉ nhân vật dời chỗ. Khớp mẫu tile cho cmx=497 ở CẢ HAI ảnh (tile65 tại x=110 và 638, tile66 tại x=302 và 974), sai số khớp 0.0-0.01. Ở pano4, từ y=542 tới y=618 KHÔNG có một pixel thân nào khác pano3 — thân biến mất 100%. Phần còn lại ở x 1443-1453, y 620-631: tôi đối chiếu từng hàng với ảnh gốc myTexture2dbong.png (46x14, mép trái theo hàng = 15,9,6,4,2,1,0,0,1,2,4,6,9,15); dãy quan sát 8,5,3,2,0,0,0,1,2,3,5,9 khớp đúng hàng 1..12 của ảnh bóng → đây là mép trái cái bóng, ảnh bóng bắt đầu ở x=1443, tâm x=1466, quy ra cx=1230. BỔ SUNG người báo cáo bỏ sót: 1230/24 = cột 51, tức nhân vật đứng LỌT VÀO trong cột tường vô hình (mục sửa số 4). Client lẽ ra phải kẹp cx tại 1213 (Char.cs:3611-3618 đặt cx = tileXofPixel(cx+chw)−chw, chw=11, tile 39 có cờ 4 trong tile_set_info.bin) — vậy tường vô hình cũng không giữ được nhân vật. Ngay cả ở mốc hợp lệ cx=1213 thì thân vẫn bị cắt mất ~40% bên phải (tâm bóng rơi vào x=1432 trên màn rộng 1454).</sub>

### 2. Đỉnh đầu nhân vật bị chém phẳng — mất tóc gai và băng trán Konoha (do chính bản sửa #1 gây ra)

**O dau:** D:\NRO_Project\v3A.png vùng x 668-745, y 532-570 (đường cắt ngang tại y=540, kéo dài x≈672→730). Nguồn: D:\NRO_Project\Assets\Resources\res\x2\smallimage\Small27045.png (thân) và Small27042.png (đầu). Thấy rõ nhất ở pano0.png vùng x 20-160, y 500-660 (nhân vật đứng trước tường nhà).

Sửa #1 chữa lỗi "hai khuôn mặt" bằng cách trỏ part ĐẦU sang ảnh trong suốt. Nhưng ảnh THÂN (Small27045) vốn được vẽ CẮT CỤT ở chân tóc — nó chỉ có mặt từ lông mày xuống, phần sọ + tóc gai + băng trán Konoha nằm ở part ĐẦU riêng (Small27041/42/43). Xoá part đầu = xoá luôn nóc sọ. Kết quả trong game: nhân vật như bị lia dao ngang đầu, mặt phẳng lì, và MẤT hoàn toàn băng trán lá Konoha — thứ nhận diện thương hiệu quan trọng nhất của một ninja Làng Lá. Lỗi "hai mặt" thật ra là do part đầu bị đặt SAI OFFSET chứ không phải do thừa part đầu.

**Bang chung:** 1) Small27045.png (64x80): bbox nội dung y=16..79, KHÔNG chạm mép trên khung → không phải bị khung cắt, mà là được vẽ cố ý cắt: hàng y=16 có 56 pixel đục liền mạch từ x=5 đến x=60 (một đường cắt ngang thẳng băng). 2) Small27042.png (64x64) là part đầu đầy đủ: tóc gai đen, băng trán navy có ký hiệu xoáy lá Konoha, bbox y=12..63. 3) Ghép thử bằng khớp màu vùng chồng lấn cho offset đầu = dx +3, dy -30 so với góc trên ảnh thân → ra nhân vật hoàn chỉnh, một khuôn mặt, có băng trán. Ảnh đối chiếu trước/sau: C:\Users\Admin\AppData\Local\Temp\claude\D--NRO-GoiBanGiao\8d082123-bbf0-446c-aa97-6342f0eda3a3\scratchpad\truoc_sau.png

**Cach sua:** Bỏ cách blank part đầu. Vẽ lại part đầu (27041/27042/27043 theo hướng nhìn) ĐÈ LÊN thân với offset dy ≈ -30, dx ≈ +3 (đơn vị pixel ảnh x2, tức -60/+6 đơn vị game nếu 1 đơn vị = 0.5px). Kiểm bằng cách so mắt của ảnh đầu (local y 46-53) trùng mắt của ảnh thân (local y 20-27).

<sub>Lang kinh `nhan-vat` · LOI MOI · phan bien: Tự kiểm và XÁC NHẬN. Đo alpha Small27045.png: hàng trên cùng có nội dung là y=16, và top-per-column từ cột 5 đến 60 ĐỀU bằng 16 — một đường cắt ngang thẳng băng, không phải chân tóc vẽ tự nhiên (27044/27046 y hệt, cắt tại y=15). Nghĩa là ảnh thân được thiết kế để có part đầu đè lên. Khớp mẫu trên v3A.png: thân 27045 vẽ tại (674,524) sai số 0.37 (khớp tuyệt đối), còn 27041/42/43 sai số ~126 ở mọi vị trí → part đầu KHÔNG hề được vẽ. Nhìn ảnh thật (v3A.png crop 650,500-780,650 và pano0.png crop 0,440-300,660): đỉnh đầu là một vạch ngang phẳng lì, không tóc gai, không băng trán. Ghép thử head 27042 tại dy=-30 cho ra nhân vật hoàn chỉnh có băng trán Konoha (scratchpad/headscan.png). Cũng xác nhận giả thuyết nguyên nhân: ghép head ở dy=0 tạo ra hiệu ứng 'hai mặt' (scratchpad/r3_twoface.png) → lỗi cũ là sai offset chứ không phải thừa part. HẠ MỨC từ chan-duong xuống nang: game vẫn chạy, map vẫn đúng, đây là lỗi sprite nhân vật chứ không chặn đường.</sub>

### 3. Nhân vật có BỐN bàn tay và HAI dải đai — ảnh thân và ảnh chân đè chồng lệch 11px

**O dau:** D:\NRO_Project\v3A.png: tay phải "trên" x 731-741 y 579-586, tay phải "dưới" x 729-736 y 590-596; bên trái hai mảnh đai cam + hai bàn tay ở y≈578-583 và y≈591-596. Xem crop scratchpad\hand_right.png (vùng 690,570-760,620 phóng 12x) và scratchpad\waist_x14.png (vùng 668,570-745,610 phóng 14x). Lặp lại y hệt ở pano0/pano1/pano2/pano3.

Ảnh THÂN Small27045 đã vẽ đủ hai cánh tay + hai bàn tay + thắt lưng + dải đai cam. Ảnh CHÂN Small27062 cũng vẽ hai bàn tay + đai cam riêng ở mép trên. Máy chủ/client đặt chân cách thân dy=+40px trong khi nội dung hai ảnh chỉ khớp ở dy≈+51 → hai bộ tay lệch nhau ~10-11px cùng hiện ra. Nhìn ở tỉ lệ 100% nó đọc thành hai cục da thịt kỳ dị hai bên hông, và hai lớp đai cam chồng nhau ở bên trái. Đây là lỗi cùng họ với lỗi "hai khuôn mặt" nhưng chưa được sửa.

**Bang chung:** Khớp mẫu pixel trên v3A.png: ảnh thân Small27045 vẽ tại (674,524) với sai số trung bình 0.37/765 (khớp gần như tuyệt đối); ảnh chân Small27062 vẽ tại (676,564) → dy = +40. Small27045 có bàn tay phải ở local x 55-62, y 65-72 → màn hình y 589-596. Small27062 có bàn tay phải ở local x 55-65, y 15-22 → màn hình y 579-586. Quét toàn bộ dy từ 20 đến 60 để tìm sai số màu nhỏ nhất ở vùng chồng lấn: tốt nhất là dy=+51, dx=+3 (sai số 50.9) — hơn hẳn dy=+40. Ảnh ghép thử ở dy=51 chỉ còn MỘT bộ tay, MỘT dải đai: scratchpad\compose_dy51.png

**Cach sua:** Đổi offset ảnh chân từ dy=+40 thành dy≈+51 (dx +3) so với góc trên ảnh thân, rồi hạ lại toàn thân cho đáy dép về đúng mép cỏ. Nếu vẫn còn chồng, cắt bỏ 13 hàng dưới cùng của ảnh thân (vùng đai + tay) và để ảnh chân đảm nhiệm phần đó — chọn một trong hai ảnh làm chủ vùng hông, đừng để cả hai cùng vẽ.

<sub>Lang kinh `nhan-vat` · LOI MOI · phan bien: Tự kiểm và XÁC NHẬN, kể cả bằng mắt. Zoom v3A.png vùng (720,570)-(755,605) phóng 20x thấy RÕ hai bàn tay màu da tách rời chồng nhau: tay trên y≈579-585, tay dưới y≈590-595 (scratchpad/r3_rhand.png). Bên trái vùng (665,568)-(700,610) có hai mảnh đai cam + hai mảng da (scratchpad/r3_lhand.png). Đo: thân 27045 vẽ tại (674,524), chân 27062 vẽ tại (676,564) → dy=+40, dx=+2 (43.6% pixel của 27062 khớp tuyệt đối tại đó, phần còn lại bị thân đè lên). Ghép thử dãy dy=40/45/48/51/54/56 (scratchpad/dyscan.png): dy=40 lộ rõ hai bộ tay + hai đai; dy≈48-51 cho ra đúng MỘT bộ tay, MỘT đai. Số của người báo cáo (dy≈51) khớp. Con số dx họ ghi là +3, đo thật là +2 — sai lệch không đáng kể.</sub>

### 4. Bóng đổ lệch sang trái 15-20px — chân sau hoàn toàn không có bóng

**O dau:** v3A.png: bóng chiếm x 671-712 (hàng y=627), tâm x≈691.5, trong khi tâm sprite = 706 và tâm sải chân = 712. pano0.png: bóng x 50-92 (hàng y=626), tâm 71, sprite tại x=54 → tâm 86. pano1.png: tâm bóng 712, tâm sprite 726. pano3.png: tâm bóng ≈1001, tâm sprite 1020.

Ở cả 4 ảnh chụp, vệt bóng nằm lệch hẳn về bên trái: nó ôm bàn chân trước và tràn ra vùng cỏ trống bên trái, còn bàn chân sau (v3A: x 719-737) nằm hoàn toàn NGOÀI bóng, đứng trơ trên nền cỏ. Mắt người đọc ngay ra là nhân vật không dính vào mặt đất. Ngoài ra bóng chỉ rộng 42px trong khi sải chân rộng 49px, nên kể cả căn đúng tâm nó vẫn hơi hụt. Ở pano3 còn tệ hơn: khoảng x 985-996 của bóng bị bôi lên VÁCH ĐẤT DỰNG ĐỨNG của gò cao (mặt cỏ gò ở y=578, nên tại y 620-630 chỗ đó là mặt cắt đứng chứ không phải mặt sàn) — bóng nằm trên tường.

**Bang chung:** Đo bằng cách so từng hàng pixel với màu cỏ trung vị lấy ở vùng sạch cùng hàng (v3A x 780-1000): hàng y=627 tối hơn ngưỡng 45 đơn vị liên tục từ x=671 đến x=712, hàng y=629 từ x=675 đến x=710, hàng y=630 từ x=676 đến x=706 → ellipse tâm ≈691. Vị trí sprite lấy bằng khớp mẫu Small27045 (sai số 0.37). Chênh lệch tâm bóng ↔ tâm sprite: v3A -14.5px, pano0 -15px, pano1 -14px, pano3 -19px → lệch có hệ thống, không phải nhiễu. Crop nhìn mắt thường: scratchpad\shadow_zoom_v3A.png (660,610-760,638 phóng 14x) và scratchpad\nv_pano0.png

**Cach sua:** Cộng thêm ~15px (≈30 đơn vị game) vào toạ độ x khi vẽ bóng trong Char.paintShadow, hoặc — đúng gốc hơn — kiểm lại xem sprite đang bị vẽ lệch +15px so với cx va chạm thật; nếu đúng vậy thì sửa offset vẽ sprite chứ đừng sửa bóng, kẻo va chạm và hình vẫn lệch nhau. Đồng thời nới bề rộng ellipse lên ≥ 52px cho phủ hết sải chân.

<sub>Lang kinh `nhan-vat` · LOI MOI · phan bien: Tự kiểm và XÁC NHẬN trên 4 ảnh độc lập. Đo bóng bằng cách so từng hàng với trung vị cỏ sạch cùng hàng: v3A.png y=626 bóng x 670..713 (tâm 691.5) trong khi tâm sprite = 706 → lệch -14.5. pano0: bóng x 50..93 tâm 71.5, sprite tâm 86 → -14.5. pano2: bóng x 704..747 tâm 725.5, sprite tâm 740 → -14.5. pano3: bóng x 984..1021 tâm 1002.5, sprite tâm 1020 → -17.5. Lệch có hệ thống. Bàn chân thật: chân trước x 688-701, chân sau x 719-737 → sải chân 688-737 (tâm 712.5, rộng 50px) nhưng bóng chỉ rộng 44px và mép phải dừng ở x=713 → chân sau nằm HOÀN TOÀN ngoài bóng. Nhìn mắt thường ở scratchpad/r3_feet1x4.png (crop 650,595-780,645 phóng 4x) thấy ngay. Chi tiết pano3 cũng đúng: tại y=626-630, x 984-995 nền là ĐẤT NÂU (89,53,27 ở y=600/620) tức mặt cắt đứng của gò cao (cỏ gò ở y=578), và chỗ đó bị bôi tối (49,44,40) → bóng đang nằm trên vách đứng.</sub>

### 5. Nét vẽ nhân vật lệch hẳn phong cách với nhà cửa/đèn/torii — pixel art thô đặt trên tranh vẽ mịn

**O dau:** So sánh 3 vùng cùng tỉ lệ 1:1 trong D:\NRO_Project\v3A.png: nhân vật (668,532)-(742,632); đèn lồng (540,440)-(614,540); cửa nhà + rèm noren (120,470)-(194,570). Ảnh ghép: scratchpad\style_cmp.png

Đặt cạnh nhau ở cùng độ phóng, hai thứ này rõ ràng thuộc hai game khác nhau. Nhà, đèn lồng, torii, cây, thùng gỗ: tranh vẽ tay/AI độ phân giải cao — chuyển sắc mượt, bóng mềm, highlight bóng gỗ, vân gỗ, giấy shoji trong mờ, đường cong khử răng cưa. Nhân vật: pixel art thô, 1 pixel ảnh = 1 pixel màn hình, cạnh răng cưa cứng, bảng màu hẹp, alpha chỉ có 0 hoặc 255. Thêm nữa tỉ lệ nhân vật là chibi (~3 đầu) trong khi kiến trúc vẽ theo tỉ lệ thực → nhân vật đọc thành hình dán đè lên tranh nền. Đây là vấn đề mỹ thuật lớn nhất của nhân vật, lớn hơn mọi lỗi ghép hình ở trên.

**Bang chung:** Kênh alpha của Small27045.png: 2741 pixel alpha=255, chỉ 45 pixel alpha=192, 0 pixel ở các mức trung gian khác → hoàn toàn không khử răng cưa. Small27062.png và Small27042.png tương tự (51 và 42 pixel trung gian). Ngược lại vùng đèn lồng/nhà trong ảnh chụp có hàng trăm sắc độ chuyển liên tục. Khớp mẫu xác nhận sprite được vẽ đúng tỉ lệ 1:1 (sai số 0.37), tức không phải do phóng to làm vỡ.

**Cach sua:** Hai hướng, chọn một và làm dứt khoát: (a) nâng nhân vật lên cùng độ phân giải nét vẽ với nền — dùng bộ x4 (128x160, đã có sẵn và có thêm chi tiết thật, sai lệch trung bình 7.3 so với bản phóng x2) rồi thu về, hoặc vẽ lại nhân vật kiểu tranh mịn; (b) hạ nền xuống chuẩn pixel art — pixelate lại toàn bộ asset môi trường về cùng cỡ pixel với nhân vật. Đừng để nguyên trạng, vì mắt người nhận ra ngay trong 1 giây.

<sub>Lang kinh `nhan-vat` · LOI MOI · phan bien: XÁC NHẬN bằng mắt, đây là vấn đề dễ thấy nhất trong ảnh. Xem pano0.png crop (0,440)-(300,660) phóng 3x (scratchpad/r3_p0_char.png): ngôi nhà có vân gỗ liên tục, giấy shoji trong mờ, ngói chuyển sắc, đường cong khử răng cưa; nhân vật ngay cạnh là pixel 1:1, cạnh bậc thang cứng, bảng màu 3-4 sắc/vùng. Cùng hiện tượng với cây khổng lồ ở v3A.png (600,470)-(850,660). Kiểm alpha xác nhận sprite không khử răng cưa: Small27045 có 2741 pixel alpha=255, 45 pixel alpha=192, 0 pixel ở mức trung gian khác; 27062/27042 tương tự (51/42). Khớp mẫu sai số 0.37 chứng minh sprite vẽ đúng tỉ lệ 1:1, không phải do phóng to làm vỡ. Tỉ lệ cũng lệch: nhân vật cao ~87px với đầu ~30px (chibi ~3 đầu) đứng cạnh cửa nhà cao ~200px vẽ theo tỉ lệ thật. Concept đã duyệt (docs/my-thuat/map-1-lang-la/concept-map1-langla.png) là tranh vẽ hoàn toàn và KHÔNG có nhân vật nào → đây là quyết định mỹ thuật chưa ai chốt, không phải đặc điểm bình thường của thể loại.</sub>

### 6. Nhân vật đi lọt ra ngoài mép phải màn hình, mất hút

**O dau:** pano4.png — góc phải, x 1444-1453, y 617-630 (chỉ còn mỏm bóng). So với pano3.png nhân vật ở x 1003-1050.

Người chơi đi sang phải tới cuối làng thì biến mất khỏi màn hình — chỉ còn một mẩu bóng đen thò ra ở mép phải. Camera đã cuộn kịch nhưng vẫn không theo kịp chỗ nhân vật được phép đứng. Với người mới, đây là lỗi đáng sợ nhất: đang điều khiển thì mất luôn nhân vật, không biết mình ở đâu.

**Bang chung:** Ghép nền pano3 và pano4 khớp nhau 100% (chỉ khác ở vùng nhân vật) → camera của cả hai ảnh đều đã ở vị trí cuộn tối đa. Đo độ lệch nền pano0→pano4 = 946px, pano0→pano2 = 626px, nên camera chỉ cuộn được 0..946. Vùng thế giới nhìn thấy được = 946+1454 = 2400px = 1200 đơn vị = 50 ô, trong khi map khai báo 52 ô (1248 đơn vị). Bóng nhân vật trong pano4 có mỏm trái tại x=1444; bóng rộng ~45px (đo ở v3A: 669-714) nên tâm bóng ≈1466 màn hình ≈ 1206 đơn vị thế giới — tức ô số 50, nằm ngoài tầm camera vĩnh viễn.

**Cach sua:** Hoặc kéo giới hạn cuộn camera thêm ~2 ô về bên phải (cho cuộn tới 1042px), hoặc dời tường vô hình vào sớm hơn 2 ô (đặt ở cột 49-50 thay vì 51) để chỗ đứng cuối cùng luôn nằm trong khung hình. Kiểm lại bằng cách đi sát mép phải: nhân vật phải luôn còn nhìn thấy toàn thân.

<sub>Lang kinh `tong-the` · LOI MOI · phan bien: CÓ THẬT, tự kiểm được. Tôi dò độ lệch nền giữa các pano bằng khớp dải đất y630-670: pano0/pano1/v3A lệch 0, pano2 lệch 626, pano3 VÀ pano4 đều lệch 946 với sai số 0.0 tuyệt đối. Trừ ảnh pano3-pano4 chỉ khác nhau ở x 983-1453, y 542-632 — nền khớp 100%, tức camera đã kẹt ở 946 trong cả hai ảnh dù nhân vật dịch sang phải ~420px. Cắt pano4 vùng (1390,590)-(1454,650) phóng 8 lần: chỉ còn một nửa elip đen (bóng) thò ra ở x 1443-1453, y 620-631; toàn bộ thân nhân vật nằm ngoài màn hình. Tôi tính lại và bác bỏ được cả cách sửa mà người báo cáo đề xuất theo hướng dời tường: map 52 ô x 48px = 2496px, màn 1454px, cuộn tối đa ĐÚNG phải là 2496-1454 = 1042px, nhưng thực tế kẹt ở 946 — thiếu đúng 96px = 2 ô. Tức camera clamp sai 2 ô chứ không phải tường đặt sai. Nếu clamp đúng 1042 thì nhân vật ở vị trí này sẽ nằm ở x màn hình ~1370, vẫn nhìn thấy toàn thân. Mức nặng giữ nguyên: đây là lỗi chơi game thật (mất nhân vật), không chỉ là lỗi nhìn.</sub>

### 7. Một đường cắt ngang toàn màn hình ở y=561 chặt đôi làng xa

**O dau:** Mọi ảnh pano, đường thẳng tại y=561→562 chạy hết bề ngang 1454px. Nhìn rõ nhất ở pano0.png x 980-1300 và pano2.png x 240-1000.

Có một đường kẻ ngang cứng đơ vắt ngang cả màn hình, đúng tầm đầu nhân vật. Phía dưới đường đó bầu trời đổi thành một dải bệt nhợt nhạt. Hậu quả: tất cả mái ngói, tán cây, ngọn núi nhỏ của 'làng xa' đều bị chặt cụt ngang thân, trông như đang chìm xuống một mặt biển phẳng. Đây là thứ đập vào mắt ngay ở lần nhìn đầu tiên và làm cả bức tranh trông dở dang.

**Bang chung:** Đo độ chênh trung bình giữa hai hàng pixel liên tiếp trên toàn bề ngang: hàng 561→562 cho giá trị 35.7 — mạnh thứ nhì cả ảnh, chỉ sau chính đường mặt đất (82.6). Lấy màu tại x=1300: y=560 là [51,163,232] (xanh đậm) nhảy sang y=562 là [173,208,230] (xám nhạt bệt). Ảnh cắt pano0 (980,530)-(1300,590) cho thấy mái chùa, mái ngói đỏ, tán cây và núi xa đều bị xén phẳng đúng tại đường này.

**Cach sua:** Lớp nền xa (b01/lớp mây) đang bị cắt cụt đáy — hoặc ảnh nguồn thiếu phần dưới, hoặc đang vẽ đè một dải màu phẳng lên. Kéo dài ảnh nền xuống tới tận đường cỏ, hoặc chèn một dải chuyển màu mềm (sương mù) cao 40-60px thay cho mép cứng. Kiểm bằng cách chụp lại và đo lại độ chênh hàng 561: phải xuống dưới 10.

<sub>Lang kinh `tong-the` · LOI MOI · phan bien: CÓ THẬT, đây là lỗi nhìn thấy đầu tiên khi mở ảnh. Đo độ chênh hai hàng liên tiếp trên toàn bề ngang: pano0 hàng 561→562 = 35.70, mạnh thứ nhì cả ảnh (chỉ sau mép cỏ 626 = 82.64) và mạnh hơn cả mép nền ở hàng 673. Kiểm tỉ lệ cột bị ảnh hưởng: pano0 66.4% số cột chênh >10, pano2 60.7%, pano4 55.5% — chạy hết bề ngang, chỗ thiếu là do vật thể che. Màu tại x=1300: y560 [51,163,232] xanh đậm → y562 [173,208,230] xám nhạt bệt. Ảnh cắt pano0 (960,510)-(1420,610) phóng 3 lần cho thấy rất rõ: mái chùa, mái ngói đỏ, tán cây, ngọn núi nhỏ của lớp làng xa đều bị xén phẳng đúng tại đường này, bên dưới là một dải bệt xanh trắng cao ~65px chạy tới mặt cỏ. Đây KHÔNG trùng lỗi đã biết (lỗi đã biết nói về ảnh nền 384px ghép đôi gây LẶP ngang và ĐƯỜNG NỐI DỌC — tôi cũng thấy các đường nối dọc đó ở x≈1153 và x≈1345, nhưng đường cắt NGANG này là khía cạnh khác). Giữ mức nặng.</sub>

### 8. Hai cái gò KHÔNG chặn ngang ở client — 13 ô client cho đứng mà server chặn

**O dau:** D:/NRO_Project/Assets/Resources/res/mymap/1.bytes hàng 17+18, cột 23-27 và 34-41 · D:/NRO_GoiBanGiao/server/internal/game/assets/map-1-collision.bin hàng 18 · D:/NRO_Project/Assets/Scripts/Game1/Char.cs:3611 và 3624

Đợt này thêm bit 4/8 (tường trái/phải) cho tile 65 và 66 — nhưng hai tile đó nằm ở HÀNG 18, trong khi client kiểm tường ở (cx±chw, cy−chh) = (cx±11, 432−16=416), tức HÀNG 17. Tile ở hàng 17 mép gò là 62, chỉ có bit 2, không có bit 4/8. Kết quả: đi bộ trên mặt đất phẳng không hề bị chặn ở chân gò — nhân vật đi thẳng VÀO trong lòng gò ở cy=432, còn máy chủ coi đúng 13 ô đó là đặc và từ chối gói. khuvuc.go:831 bỏ gói, giữ vị trí cũ, không phát → người chơi khác thấy đứng chết một chỗ. Đúng cái mà commit 3627c8f định tránh. File va chạm tự nó ĐÚNG theo quy ước của nó; lỗi nằm ở chỗ client không có gì bắt buộc quy ước đó.

**Bang chung:** tile_set_info.bin nhóm 1: bit 4 = [1,23,28,39,29,30,60,61,65,66], bit 8 = [18,24,29,39,28,30,60,61,65,66] — KHÔNG có 62. 1.bytes: (17,23)=(17,27)=(17,34)=(17,41)=62. Đối chiếu ô-với-ô: đúng 13 ô vừa có bit 2 (client cho đứng) vừa =1 trong map-1-collision.bin, là hàng 18 cột 23-27 (5 ô) và 34-41 (8 ô). CHỨNG THỰC TRONG ẢNH: D:/NRO_Project/pano3.png — camera cmx=497, cmy=119 (khớp template sai số 0.0). Bóng nhân vật (rộng đúng 46 px) có tâm ở pano x≈1004 → world x≈999 → cột 41; chân ở pano y 626 → world y 432 → hàng 18. Ô (18,41) trong map-1-collision.bin = 1. Xem thêm ảnh cắt scratchpad/p3_hill.png: thân nhân vật đè lên vách gò.

**Cach sua:** Cho bit 4 vào tile ở HÀNG 17 cột 23 và 34 (mép trái gò), bit 8 vào hàng 17 cột 27 và 41 (mép phải) — hoặc đổi bốn ô đó sang một tile mới vừa có bit 2 vừa có bit 4/8 (giống tile 1/18 trong bộ gốc, hoặc tile 39). Sau đó chạy lại tools/map-lang-la/sinh-va-cham.py. Thêm một test đối chiếu: không ô nào được vừa có bit 2 ở client vừa =1 ở server trong hàng mà nhân vật đi bộ tới được.

<sub>Lang kinh `loi-moi` · LOI MOI · phan bien: TỰ KIỂM LẠI TỪ ĐẦU, ĐÚNG HOÀN TOÀN. (a) Tôi parse lại tile_set_info.bin (39 nhóm, tiêu thụ đủ 1716/1716 byte): nhóm 0 có bit 2 = [1,2,3,4,5,6,18,31,32,33,34,39,60,61,62,63,64,65,66,67,68,69]; bit 4 = [1,23,28,29,30,39,60,61,65,66]; bit 8 = [18,24,28,29,30,39,60,61,65,66]. Tile 62 CÓ bit 2 nhưng KHÔNG có bit 4/8 — đúng như báo cáo. So với tile_set_info.bin.bak-goc thì đợt này đúng là mới thêm 60,61,65,66 vào bit 4 và 8. (b) 1.bytes (header 2 byte, 52x24): hàng 17 cột 23-27 và 34-41 = 62; hàng 18 cột 23=65,24=60,25=61,26=61,27=66 và 34=65,35..40=60/61,41=66. (c) Char.cs:209/211 chw=11, chh=16 → kiểm tường ở cy-chh = 432-16 = 416 → HÀNG 17, nơi chỉ có tile 62 ⇒ không có bit 4/8 ⇒ không chặn. Kiểm rơi ở cy+1=433 → hàng 18, tile 65/60/61/66 đều có bit 2 ⇒ nhân vật đi thẳng vào lòng gò. (d) LƯU Ý: map-1-collision.bin có header 4 BYTE (52,24,24,24) chứ không phải 2 — người báo cáo đọc đúng, tôi đọc sai lần đầu rồi sửa. Đọc đúng ra hàng 18 đặc tại cột 23-27 và 34-41 = ĐÚNG 13 ô, trùng khít 13 ô client cho đứng. (e) catalog.go:177 DiDuoc điểm-hoá (x/24, y/24), khuvuc.go 'return // giu vi tri cu, KHONG phat' ⇒ đúng là bỏ gói, đứng chết với người khác. (f) BẰNG CHỨNG ẢNH TÔI TỰ ĐO: pano3.png — tôi xác nhận cmx=497 độc lập (mép phải gò ở world 42*24=1008 → pano x 1022; ảnh cho 1023). Bóng nhân vật rộng đúng 46 px, tâm pano x≈1012 → cx≈1003 → cột 41, chân ở pano y 627 → world 432 → hàng 18. Ô (18,41)=1 trên server. Ảnh cắt scratchpad/r_p3_char.png và r_p3_tile66.png cho thấy RẤT rõ nhân vật + bóng được vẽ ĐÈ lên khối đất của gò, đứng lơ lửng giữa lòng gò. Giữ mức 'nang': vừa là lỗi nhìn thấy được (đi xuyên vào sườn gò 8 ô liền), vừa là desync máy chủ.</sub>

### 9. Sửa 'hai khuôn mặt' làm nhân vật MẤT NÓC ĐẦU — ảnh đầu tốt vẫn còn trong bộ mà bị bỏ

**O dau:** D:/NRO_GoiBanGiao/server/internal/game/assets/part.bin part 30 · D:/NRO_Project/Assets/Resources/res/x2/smallimage/Small27044.png .. Small27060.png · ảnh đầu đầy đủ: Small27041/27042/27043.png

part 30 (đầu) bị trỏ từ 27041/27042/27043 sang 27083 (ảnh trong suốt) với lý do 'ảnh thân đã có đầu'. Nhưng ảnh thân chỉ có NỬA DƯỚI của đầu: nó bị cắt phẳng ngang ngay trên chân mày, không có sọ, không có tóc. Trong game nhân vật giờ có cái đầu bằng phẳng như bị vạt, và mắt phải nằm sát mép ảnh. Ba ảnh đầu ĐẦY ĐỦ (tóc dựng + băng trán lá) vẫn nằm nguyên trong bộ, chỉ là không được vẽ nữa. Vòng 2 báo 'hai khuôn mặt' là do hai mảnh lệch nhau, không phải do thừa mảnh — bỏ hẳn mảnh đầu là chữa triệu chứng.

**Bang chung:** part.bin so với part.bin.bak-goc: part 30 pi[0..2] (27041,0,0)/(27042,0,0)/(27043,0,0) -> (27083,0,-22) cả ba. Đo Small27044.png (x2, 64x80): hàng pixel không trong suốt trên cùng là hàng 15, rộng 56/64 px, toàn màu tóc (20,23,34)/(61,78,102) — một cái đầu bình thường thì hàng trên cùng chỉ vài pixel. Bbox chạm cả x=63 (mép phải khung). Trong game: pano1.png vùng (690,525)-(770,575) và pano2.png (690,510)-(790,570) — đỉnh đầu là một đường cắt ngang phẳng lì. Ảnh Small27041.png có đầu nguyên vẹn.

**Cach sua:** Trả part 30 về 27041/27042/27043 rồi chỉnh dx/dy cho khớp với thân (thân dy=-46), đồng thời xoá phần mặt đã vẽ sẵn trong 27044-27060 để hai mảnh không chồng. Hoặc ngược lại: vẽ nốt nóc đầu vào 17 ảnh thân. Không có cách nào giữ nguyên hiện trạng mà nhân vật có đầu đủ.

<sub>Lang kinh `loi-moi` · LOI MOI · phan bien: ĐÚNG, và bằng chứng còn mạnh hơn báo cáo. Tôi parse part.bin vs part.bin.bak-goc bằng đúng định dạng client đọc (GameScr.readPart: short count, byte type, mỗi PartImage = short id + byte dx + byte dy; 2192 part, tiêu thụ đủ 87846 byte). Chỉ 3 part đổi: part 30 (đầu) pi[0..2] (27041/27042/27043, 0, 0) -> (27083, 0, -22) cả ba; part 57 dy 0->-46; part 58 dy 0->-33. Small27083.png có alpha lớn nhất = 0 (trong suốt tuyệt đối) ⇒ part đầu KHÔNG vẽ gì. Tôi mở Small27044.png (thân, 64x80): hàng pixel đầu tiên có màu là hàng 15 và bbox chạm x=0..63 — nhìn bằng mắt (scratchpad/r_parts.png) thì thân chỉ có NỬA DƯỚI của đầu, cắt phẳng ngay trên chân mày, không sọ không tóc. Small27041.png là một cái đầu HOÀN CHỈNH (tóc dựng + băng trán lá) vẫn nằm nguyên trong bộ. TRONG GAME: tôi cắt zoom đầu nhân vật ở cả pano3.png (scratchpad/r_p3_head.png) và v3A.png (scratchpad/r_v3a_head.png) — đỉnh đầu là một đường cắt ngang phẳng lì, không thể nhầm. Giữ 'nang'.</sub>


## Vua

### 10. Nửa trên của bóng nổi lơ lửng TRÊN NỀN TRỜI, phía trên mép cỏ (bản sửa #7 bù quá tay)

**O dau:** Mép cỏ ở y=626; bóng chiếm y≈620-632. pano0.png cột x=55: y 622-625 có màu (57,65,72)/(52,59,65) trong khi trời cùng hàng ở x=44 là (157,196,230). v3A.png cột x=710: y 623-625 ≈ (47,60,66) so với trời (99,188,234).

Sửa #7 nâng bóng lên 12 đơn vị để nó thôi chìm dưới đất. Nhưng nâng quá tay: giờ tâm ellipse rơi đúng vào mép cỏ, nên gần nửa trên của ellipse được vẽ lên vùng TRỜI. Vì bản đồ này có đường chân đất là một vạch ngang sắc lẻm với trời xanh sáng ngay bên trên, cái nửa ellipse đó hiện thành một vệt xám đậm treo giữa không khí ngay dưới bụng nhân vật — nhìn còn khó chịu hơn lỗi cũ. Nó cũng che mất một mảng trời giữa hai chân.

**Bang chung:** pano0.png hàng y=624 quét theo x: x=49 (153,192,224) là trời, x=51 đã là (58,66,72), rồi giữ nguyên (52,59,65) suốt tới x=59 — tức 4-5 hàng pixel (y 621→625) nằm HOÀN TOÀN trên nền trời, phía trên mép cỏ y=626. Cùng hiện tượng ở v3A.png (x 702-712, y 620-625) và pano1.png. Nhìn mắt thường: scratchpad\shadow_zoom_v3A.png, scratchpad\shadow_zoom_pano1.png, scratchpad\ground_pano1.png

**Cach sua:** Hạ tâm ellipse xuống thêm khoảng nửa chiều cao của nó (~6px ≈ 3 đơn vị) để toàn bộ bóng nằm DƯỚI mép cỏ, hoặc cắt (clip) bóng theo đường chân đất của ô đang đứng. Cách chắc nhất: vẽ bóng vào lớp mặt đất và cho clip theo y >= mặt cỏ.

<sub>Lang kinh `nhan-vat` · LOI MOI · phan bien: XÁC NHẬN nhưng người báo cáo THỔI PHỒNG mức độ. Mép cỏ = y626 (v3A x=790: y625 = trời (142,197,235), y626 = cỏ (86,88,22)). Cột trống giữa hai chân x=712: y622 (86,145,141), y623 (66,102,104), y624 (49,64,72), y625 (47,57,62) — đều tối hơn hẳn trời (101,188,234) → có 4 hàng bóng nằm trên trời. x=708 và x=710 tương tự; x=716 sạch hoàn toàn. Nhưng ellipse trải y622..632 = 11 hàng, phần trên trời là 4/11 ≈ 36%, không phải 'gần nửa'. 4 pixel trên nền trời sáng thì thấy được nhưng đọc thành vệt mờ bám mép đất chứ không phải 'treo giữa không khí' — hạ từ nang xuống vua.</sub>

### 11. Ô chân dung HUD VẪN TRỐNG — bản sửa #8 không hề ăn

**O dau:** D:\NRO_Project\v3A.png ô x 8-70, y 10-58 (góc trên trái). Giống hệt ở pano0.png, pano1.png, pano2.png, pano3.png, pano4.png.

Sửa #8 nói máy chủ đã gửi ánh xạ head 30 → avatar 27078. Thực tế trong ảnh chụp mới nhất (18:42-18:43 cùng phiên) ô chân dung vẫn chỉ là tấm nền xanh lá cùng vài khối placeholder xanh đậm hơn — không có một pixel nào của mặt người. Ảnh 27078 tồn tại và rất đẹp (100x100, tóc nâu, băng trán Konoha) nhưng không được vẽ ra. Lưu ý: bản sửa #9 thì ĐÃ ăn — thanh máu/chakra hiện 150/150 đúng.

**Bang chung:** Đếm màu trong ô (8,10)-(70,58) của v3A.png: 271 màu nhưng toàn bộ nằm trong dải xanh lá — 4 màu nhiều nhất là (137,231,170) 1001px, (101,210,112) 731px, (138,232,172) 100px, (0,83,37) 90px. Không có màu da (~230,180,120) cũng không có màu tóc nâu (90,57,27) vốn chiếm phần lớn Small27078.png. Crop: scratchpad\avatar_slot.png (phóng 14x) và scratchpad\hud_x4.png

**Cach sua:** Kiểm lại phía CLIENT: xem hàm vẽ HUD có đọc đúng id 27078 và tìm ảnh ở đúng thư mục Resources/res/x2/smallimage không (Small27078.png là 100x100, KHÔNG phải 64x64 như các part khác — nhiều khả năng client cắt/khung theo 64x64 nên vẽ trượt hoặc bỏ qua). Log id ảnh mà HUD đang yêu cầu ra console rồi đối chiếu.

<sub>Lang kinh `nhan-vat` · phan bien: XÁC NHẬN. Đếm màu vùng (8,10)-(70,58) của v3A.png: 278 màu, 4 màu nhiều nhất đều là xanh lá ((137,231,170) 1019px, (101,210,112) 792px, (138,232,172) 101px, (0,83,37) 90px); lọc màu da (R>190, 130<G<210, B<160) chỉ ra 4 màu tổng vài pixel, thuộc viền cam của khung chứ không phải mặt người. Nhìn ảnh phóng 5x (scratchpad/r3_hud.png, r3_hud3.png): chỉ có khung xanh lá + chữ 'narruto' + số 0, tuyệt đối không có chân dung. Giống hệt ở pano3. Đồng thời xác nhận bản sửa #9 ĐÃ ăn: hai thanh hiện 150/150.</sub>

### 12. Bóng đổ lệch hẳn sang trái ~20px, chân phải không có bóng

**O dau:** v3A.png: bóng x 669-714, hai bàn chân x 689-737. pano0.png: bóng x 49-94, chân x 69-116. pano3.png: bóng x 983-1018, chân x 1003-1050. Tất cả ở y 619-632.

Việc sửa số 7 mới chỉ chỉnh được ĐỘ CAO của bóng (bóng không còn chìm dưới đất). Nhưng bóng vẫn nằm lệch sang trái khoảng 20px so với hai bàn chân, nên bàn chân phải hoàn toàn không có bóng, còn bóng thì thò ra một khoảng trống bên trái. Nhìn như nhân vật đang đứng cạnh cái bóng của mình chứ không phải đứng trên nó.

**Bang chung:** Trừ ảnh (v3A trừ pano0, pano3 trừ pano4) để tách riêng nhân vật: tâm cụm hai bàn chân lần lượt là x=712 / 92 / 1026; tâm elip bóng (đo ở các hàng 629-632, chỉ còn bóng, không còn chân) là x=691 / 71 / 1005. Chênh lệch ổn định 18-21px ở cả 3 ảnh, tương đương ~10 đơn vị map. Ngoài ra mép phải bóng (714) nằm hẳn bên trái mép phải bàn chân (737).

**Cach sua:** Bóng đang vẽ đúng ở toạ độ thật của nhân vật (kiểm chứng: ở tường trái pano0 tâm bóng = 35.5 đơn vị, hợp với tường cột 0 rộng 24 + nửa thân ~11). Vậy chính THÂN nhân vật mới đang bị vẽ lệch phải ~10 đơn vị. Sửa dx của các part thân/chân trong bảng xếp sprite (trừ đi ~10 đơn vị), rồi chụp lại kiểm: tâm bóng phải trùng tâm hai bàn chân. Nhớ thử cả lúc quay mặt sang trái để chắc không phải lỗi lật hình.

<sub>Lang kinh `tong-the` · LOI MOI · phan bien: CÓ THẬT nhưng người báo cáo thổi mức độ. Tôi tách nhân vật bằng cách trừ v3A - pano0 (cùng camera 0) rồi liệt kê từng hàng: ở v3A hàng 617-619 có hai cụm rời — bàn chân trái x 689-703, bàn chân phải x 719-736 (tâm cụm chân = 712.5); từ hàng 625 trở xuống chỉ còn một elip liền x 669-714 (hàng 632 thu về 683-700), tâm = 691.5. Lệch 21px. Làm lại với nhân vật ở pano0: bóng 49-94 (tâm 71.5), chân trái 69-80 + chân phải 99-117 (tâm 93) → lệch 21.5px. Hai ảnh cho cùng con số, không phải nhiễu. Nhìn ảnh cắt (655,530)-(755,645) phóng 6 lần thấy rõ: elip bóng nằm dưới bàn chân trái và thò hẳn ra bên trái, bàn chân phải đứng trên nền cỏ trắng trơn. Cũng thấy lại ở pano2 quanh x 700-745. Hạ từ 'nặng' xuống 'vừa': đây là lỗi thẩm mỹ, ở cỡ 1:1 phải nhìn kỹ mới thấy, không ảnh hưởng điều khiển. Lưu ý cách sửa người báo cáo đề xuất (dịch thân trái 10 đơn vị) chỉ là SUY ĐOÁN — bằng chứng chỉ chứng minh bóng và chân không trùng nhau, chưa chứng minh bên nào sai.</sub>

### 13. Nhà bên trái treo lơ lửng, hở trời có mây ở khe dưới chân

**O dau:** pano0.png / pano1.png — dải x 100-320, y 622-625 (ngay dưới chân căn nhà đầu map). Cũng thấy ở nhà lớn: pano4.png x 660-1020, y 575-577.

Chân căn nhà kết thúc ở y=622 nhưng mặt cỏ bắt đầu ở y=626, chừa ra một khe cao 3-4px. Qua khe đó nhìn thấy nguyên trời xanh và cả mây trắng. Ở kích thước thật, đây là một vệt sáng chói nằm dưới một khối gỗ tối màu — rất bắt mắt, người chơi sẽ nghĩ 'nhà chưa đặt xong'.

**Bang chung:** Lấy màu theo cột dọc trong pano0: x=150 các hàng 623/624/625 = [117,190,232] (trời xanh), hàng 626 mới là cỏ [102,101,26]. x=250 hàng 622-625 = [141,144,148],[224,229,234],[205,219,231],[178,206,230] — đúng màu mây trắng. Ảnh phóng to (100,590)-(340,650) thấy rõ dải mây chạy suốt bề ngang căn nhà, kể cả dưới hai bệ đá trước cửa.

**Cach sua:** Hạ toạ độ đặt các vật thể xuống 2 đơn vị (4px ở zoom 2), hoặc sửa điểm neo của sprite nhà từ 'giữa thân' sang 'đáy'. Cách kiểm nhanh sau khi sửa: quét cột dọc, không được còn pixel trời nào ở hàng 623-625 tại vùng có nhà.

<sub>Lang kinh `tong-the` · LOI MOI · phan bien: CÓ THẬT ở chỗ chính, nhưng chỗ phụ thì người báo cáo bịa thêm. Quét cột dọc pano0: x=150 hàng 623/624/625 = [100,160,195],[117,190,232],[112,186,228] (trời xanh) rồi hàng 626 mới là cỏ [102,101,26]; x=200 hàng 623-625 = [187,213,232],[174,206,231],[149,193,227]; x=250 = [224,229,234],[205,219,231],[178,206,230] (mây trắng); x=300 hàng 624-625 cũng là trời. Ảnh cắt (120,600)-(360,640) phóng 5 lần: thấy rõ một sọc sáng 3px chạy suốt bề ngang căn nhà, kể cả dưới hai bệ đá. NHƯNG chỗ thứ hai người báo cáo nêu (pano4 nhà lớn, x 660-1020, y 575-577) KHÔNG đúng: tôi đo x=700/800/900 chỉ có ĐÚNG MỘT hàng 577 màu pha tối [37,63,79] / [62,69,78] / [56,75,86] — là viền khử răng cưa, không phải khe hở trời. Hạ từ 'nặng' xuống 'vừa': khe chỉ 3px, ở cỡ 1:1 tôi xem ảnh đầy đủ mà không tự phát hiện ra. Lưu ý phát hiện này và phát hiện 'không vật nào lún vào đất' cùng một gốc (vật thể bị neo cao hơn mặt đất ~2 đơn vị) — nên gộp khi sửa.</sub>

### 14. Ô avatar trên HUD vẫn là mảng hình rác, không phải mặt nhân vật

**O dau:** Cả 5 ảnh pano — góc trên trái, x 14-46, y 20-45 (bên trong khung xanh lá, phía trước chữ 'narruto').

Việc sửa số 8 KHÔNG đạt. Ô đựng mặt nhân vật vẫn chỉ hiện vài khối chữ nhật xanh sẫm nằm rời rạc trên nền xanh lá — không ra mặt mũi gì. Người chơi mới nhìn vào sẽ thấy phần quan trọng nhất của giao diện (ảnh đại diện của chính mình) đang hỏng.

**Bang chung:** Phóng to vùng (8,10)-(60,62) của pano0 lên 10 lần: chỉ thấy 4-5 hình chữ nhật màu xanh đậm hơn nền, xếp thành hình giống chữ 'L' và vài chấm — không có đường nét khuôn mặt, không có tóc, không có màu da. Giống hệt nhau ở cả pano0..pano4, nghĩa là không phải lỗi chớp nhoáng.

**Cach sua:** Kiểm lại xem client có thật sự tải được ảnh 27078 không (bật log phía client khi nạp icon HUD). Nhiều khả năng id vẫn trỏ vào một ảnh trống/sai kích thước. Trước mắt nên vẽ tạm một ảnh mặt 32x32 riêng cho HUD thay vì tái dùng part đầu của sprite thân — vừa chắc ăn vừa nhìn rõ hơn ở cỡ nhỏ.

<sub>Lang kinh `tong-the` · phan bien: CÓ THẬT, việc sửa số 8 thất bại. Cắt pano0 (0,0)-(180,75) phóng 8 lần: trong khung xanh lá, vùng x≈14-46 / y≈20-45 chỉ có một hình góc kiểu dấu ngoặc '⌐' và vài ô vuông nhỏ màu xanh sẫm hơn nền — không có tóc, không có màu da, không có nét mặt. Trông đúng như biểu tượng 'ảnh hỏng / ảnh thiếu'. Giống hệt nhau ở cả pano0, pano2, pano3, pano4 nên không phải lỗi chớp nhoáng. Hạ từ 'nặng' xuống 'vừa': vùng hỏng chỉ 32x26px ở góc màn, không cản trở gì; nhưng vì đây là việc VỪA SỬA MÀ KHÔNG ĂN nên phải báo lại.</sub>

### 15. 22% đáy màn hình là khối đất nâu phẳng lì, không có gì

**O dau:** Mọi ảnh pano — dải y 673-862, chạy hết bề ngang 1454px.

Đây là thứ làm màn hình mất cân bằng nặng nhất. Từ y=673 xuống hết đáy (189px, gần 1/4 màn hình) chỉ là một mảng nâu đều đều: không đá, không rễ cây, không mạch nước, không hang. Cộng thêm dải trời gần như trống ở trên (y 70-222, 18%), tức là khoảng 40% màn hình không mang thông tin gì, còn cả ngôi làng bị ép vào một dải hẹp ở giữa. Đây là khác biệt lớn nhất so với tranh concept — trong concept phần đất dưới có vách đá, cỏ mọc, tầng lớp rõ ràng. (Ghi chú: khác với lỗi đã biết 'tầng đất sâu còn lộ ô vuông' — ở đây nói về TỈ LỆ bố cục, không phải chuyện ô vuông.)

**Bang chung:** Đo độ chênh giữa hai hàng pixel liên tiếp trên toàn bề ngang pano4: hàng 673 nhảy vọt lên 19.9, các hàng 665-672 và 674-681 chỉ 3-4 → đúng ranh giới giữa lớp đất có vân và lớp bệt. Độ lệch chuẩn màu của vùng 676-862 chỉ 27.2, trong khi vùng trời 300-560 là 66.2. Tỉ lệ chiều cao: (862-673)/862 = 21.9%. Đỉnh mái cao nhất là y=222 (pano4, x=696), nên dải trời trống phía trên = 222-70 = 152px = 17.6%.

**Cach sua:** Hai hướng, nên làm cả hai: (1) đưa camera xuống thấp hơn hoặc cắt bớt số hàng đất phía dưới để dải đất chỉ còn ~10% màn hình; (2) nhét nội dung vào phần đất — mạch đá, rễ cây to xuyên xuống, một đường hầm/cống, vài viên đá lớn — đúng như tranh concept có vách đất nhiều tầng.

<sub>Lang kinh `tong-the` · LOI MOI · phan bien: CÓ THẬT về mặt số liệu, tôi đo lại và ra đúng. pano4: độ chênh hàng 672→673 = 19.87 trong khi các hàng 665-672 và 674-684 chỉ 3.1-4.4 → ranh giới rõ giữa lớp đất có vân và lớp bệt. (862-673)/862 = 21.9%. Đỉnh mái ngói cao nhất pano4 = y 223 tại x 698 (tôi lọc pixel ngói đỏ R>150, R-G>60, R-B>70, bỏ vùng HUD/logo) → dải trời gần trống phía trên đúng như báo cáo. So với concept art (D:/NRO_GoiBanGiao/docs/my-thuat/map-1-lang-la/concept-map1-langla.png) thì đúng: concept có vách đất nhiều tầng, cỏ mọc, đá, lối mòn ở phần dưới, còn trong game là mảng nâu trơn. NHƯNG bằng chứng 'độ lệch chuẩn 27.2' của người báo cáo là YẾU và tôi bác bỏ: tôi đo dải đất mặt 626-672 (dải CÓ vân rễ nhìn rõ) cho std 29.3 — gần bằng, tức std không phân biệt được 'có vân' và 'bệt'. Cái đúng nằm ở số 19.87 và ở mắt thường. Giữ mức 'vừa'.</sub>

### 16. Núi Hokage thấp hơn cả mái nhà — không đọc ra là Konoha

**O dau:** pano4.png — đỉnh dãy núi ở y≈236; đỉnh mái nhà lớn ở y=222 (x=696). Vùng núi: y 236-337, cao chỉ ~100px trên màn 862px.

Người chơi mới nhìn vào không nhận ra đây là làng Lá. Vách đá Hokage — biểu tượng số một của Naruto — đang bé hơn một căn nhà hai tầng bình thường: mái nhà còn nhô cao hơn cả đỉnh núi. Màu núi lại nhợt, xám xanh, tương phản thấp, ở cỡ thật trông như một vệt bẩn ngang trời chứ không phải tượng đá khổng lồ. Trong tranh concept núi chiếm gần 40% chiều cao và là điểm neo thị giác của cả cảnh. (Khác với lỗi đã biết về việc núi lặp 2 lần / có đường nối dọc — đây là chuyện TỈ LỆ và độ tương phản.)

**Bang chung:** Quét pixel màu ngói đỏ (loại trừ HUD và logo): đỉnh mái cao nhất pano4 = y 222. Ảnh cắt phóng to pano0 (150,230)-(700,345) cho thấy đá núi bắt đầu ở y≈236 và tan dần vào trời quanh y 333-345. Vậy núi cao 100px, nhà cao (626-222)=404px — nhà cao gấp 4 lần núi. Trong concept-map1-langla.png núi cao ~390/1024 px còn nhà ~330/1024 px (núi cao hơn nhà).

**Cach sua:** Phóng lớp nền núi (b02) lên ít nhất 2.5-3 lần chiều cao, đẩy đỉnh núi lên khoảng y=90-120 để nó vượt hẳn mọi mái nhà, và tăng độ tương phản/độ đậm màu đá. Vì lớp nền đang là ảnh 384px ghép đôi, tiện thể vẽ lại một ảnh núi rộng đủ 1 màn để bỏ luôn chuyện lặp.

<sub>Lang kinh `tong-the` · LOI MOI · phan bien: Sự việc CÓ THẬT nhưng phần diễn giải thì SAI, phải sửa lại. Đo được: đá núi trong pano0 bắt đầu ở y≈235 và tan vào trời quanh y≈340 → cao ~105px trên màn 865px (12%). Đỉnh mái nhà lớn pano4 = y 223, mặt đất 626 → nhà cao 403px, gấp ~3.8 lần chiều cao núi, và đỉnh mái nhô cao hơn đỉnh núi 12px. Concept art thì núi trải từ y≈100 đến y≈700 trên 1024 và trùm hẳn lên mọi mái nhà. Vậy vấn đề TỈ LỆ là thật. NHƯNG tôi bác bỏ hai ý: (a) 'màu nhợt, tương phản thấp, trông như vệt bẩn' — cắt (150,215)-(760,360) phóng 2 lần thấy vách đá vẽ khá kỹ, và ngay ở ảnh 1:1 tôi vẫn nhận ra 5 khuôn mặt Hokage ngay lập tức; (b) 'người chơi không nhận ra đây là làng Lá' — sai, cổng torii đỏ + vách mặt Hokage + mái ngói Nhật là quá đủ dấu hiệu. Giữ 'vừa', nhưng phải mô tả lại là lỗi TỈ LỆ/bố cục, không phải lỗi nhận diện.</sub>

### 17. Vệt sọc xám rộng 7px lòi ra ở mép mọi gò cao

**O dau:** pano4.png: x 343-349 và x 1015-1021, y 626-670. pano3.png: cùng vị trí. pano2.png: x 663-670 (và một vệt sẫm hơn ở x 431-438), y 626-672.

Ngay dưới mép phải của mỗi gò cao có một dải dọc rộng đúng 7px gồm các vạch xám ngang xếp chồng — trông như một cái thang đá hay khối đá bị bóp méo còn sót lại. Nó nằm chình ình giữa mặt đất nâu nên rất dễ thấy. Đây là rác hiển thị, không phải chi tiết cố ý.

**Bang chung:** Lọc pixel xám (chênh lệch R-G-B < 25, độ sáng 95-235) trong vùng y 620-690: pano4 trả về đúng các cột 343,344,345,346,347,348,349 và 1015,1016,1017,1018,1019,1020,1021 — trùng khít mép phải hai gò cao (gò 1 kết thúc ở 349, gò 2 ở 1021). Ảnh phóng to (300,600)-(420,700) và (620,560)-(720,700) cho thấy rõ dải vạch ngang xám trắng.

**Cach sua:** Đây gần như chắc chắn là một ô tile bị vẽ sai chiều rộng (7px thay vì 48px) hoặc một tile 'kè đá / thang' bị lấy nhầm chỉ số ở ô chuyển tiếp cuối gò. Rà lại bảng tile ở cột cuối của mỗi gò (cột 27 và cột 41) trong hàng ngay dưới mặt đất, đổi về tile đất bình thường hoặc vẽ hẳn một tile kè đá đúng khổ 24 đơn vị.

<sub>Lang kinh `tong-the` · LOI MOI · phan bien: CÓ THẬT và tôi tìm được bằng chứng mạnh hơn cả người báo cáo. Cắt pano4 (300,600)-(420,700) phóng 6 lần: thấy rõ một dải dọc hẹp gồm các vạch ngang nhạt xếp chồng, nằm đúng chỗ tiếp giáp giữa ô đất dưới gò và ô đất thường. Đo màu hàng 645: x 342 = [87,51,22], rồi x 343,344,345,346,347,348,349 đều GIỐNG HỆT NHAU [115,93,76], rồi x 350 = [82,50,26]. Bảy cột màu y hệt nhau = một cột pixel bị kéo giãn 7 lần, tức tile bị vẽ sai chiều rộng chứ không phải hoạ tiết cố ý. Độ bão hoà của dải này là 39 so với 53-83 của đất xung quanh nên nó đọc ra 'xám' đúng như mô tả. Vị trí xác nhận: pano4 x 343-349 và x 1015-1021 (y ~627-673), pano2 x ~663-670 — trùng mép phải của từng gò. Giữ mức 'vừa'.</sub>

### 18. Dưới gò cao là hai mảng đá rời và một khoảng nâu trơn — trông như chưa lấp xong

**O dau:** pano2.png — dải y 625-673: mảng đá sáng ở x 433-480 và x 620-667, khoảng nâu trơn không vân ở x 480-620.

Cái gò cao trông như một tấm ván nổi: dưới bụng nó không phải đất liền mà là hai miếng đá vụn nằm rời hai đầu, ở giữa là một mảng nâu trơn hoàn toàn không có vân. Nhìn ra ngay là 'chưa lấp xong', và làm cả gò mất cảm giác nặng, mất cảm giác có thật.

**Bang chung:** Ảnh cắt pano2 (400,560)-(780,700) phóng 3 lần: thấy rõ ba vùng khác chất liệu nằm cạnh nhau trong cùng một hàng — đá lốm đốm sáng (x 433-480), nâu bệt không vân (x 480-620), đá lốm đốm sáng (x 620-667). Trong khi đó vùng đất dưới mặt đất thường (x 672-780) lại là đất có rễ màu sẫm, khác hẳn.

**Cach sua:** Lấp toàn bộ hàng ngay dưới gò bằng cùng loại tile đất-có-rễ đang dùng cho mặt đất thường, hoặc vẽ một tile 'vách đá dưới gò' liền mạch. Nguyên tắc chung: mọi ô đất nằm ngay dưới một mặt đứng phải cùng một chất liệu, không được lẫn 3 loại trong 5 ô.

<sub>Lang kinh `tong-the` · LOI MOI · phan bien: CÓ THẬT. Cắt pano2 (400,560)-(790,700) phóng 3 lần thấy rõ ba chất liệu khác nhau nằm liền kề trong cùng một hàng ngay dưới gò: mảng đá lốm đốm sáng ở x≈433-480, mảng nâu trơn hoàn toàn không vân ở x≈480-620, mảng đá lốm đốm sáng ở x≈620-667; trong khi đất dưới mặt đất thường (x>672) là đất nâu sẫm có rễ, khác hẳn. Ngay trong cùng ảnh cắt còn thấy vệt sọc 7px (phát hiện phía trên) dính vào mép phải mảng đá thứ hai. Ghi nhận thêm: mặt đứng của gò (từ y 578 xuống 622) thì lại vẽ đúng chất liệu đất-có-rễ, nên lỗi chỉ nằm ở hàng ô NGAY DƯỚI gò. Giữ mức 'vừa'.</sub>

### 19. Lối lát đá đọc ra như một rãnh sỏi, không phải mặt đường

**O dau:** pano0.png / pano1.png — x 355-640, y 623-672 (đoạn giữa nhà đầu map và cây cột đèn).

Đoạn lát đá duy nhất của cả map không đọc ra là con đường. Ở đó lớp cỏ biến mất hoàn toàn và cả khối đất một ô bị thay bằng một mảng sỏi sáng màu — kể cả phần mặt cắt đứng. Kết quả: nhìn như người ta đào một cái rãnh sỏi xuống đất, chứ không phải rải đá lên mặt đường. Hai đầu thì cỏ bị cắt đứt bằng một đường dọc thẳng đơ, không có chuyển tiếp mềm nào. Ngoài ra cả map 50 ô chỉ có đúng đoạn đường này rồi thôi — không có đường nối các nhà với nhau.

**Bang chung:** Ảnh cắt pano0 (300,590)-(700,700) phóng 3 lần: bên trái x<355 và bên phải x>640 là cỏ xanh mỏng phủ trên đất sẫm có rễ; ở giữa hoàn toàn không còn vệt cỏ nào, cả chiều cao 623-672 là một chất liệu sỏi sáng đồng nhất. Đường ranh hai bên là đường thẳng đứng, không có ô pha trộn.

**Cach sua:** Tách làm hai tile khác nhau: một tile 'mặt đường lát đá' chỉ chiếm 6-10px trên cùng (thay cho lớp cỏ), phần dưới vẫn giữ nguyên đất có rễ như chỗ khác. Thêm 1 ô chuyển tiếp mỗi đầu có cỏ mọc lấn lên mép đá. Và kéo dài con đường chạy suốt map, đi qua chân các nhà — đó chính là thứ làm nó đọc ra 'một ngôi làng'.

<sub>Lang kinh `tong-the` · LOI MOI · phan bien: CÓ THẬT. Cắt pano0 (280,610)-(720,690) phóng 3 lần: bên trái x<355 và bên phải x>640 là lớp cỏ xanh mỏng nằm trên đất sẫm có rễ; ở giữa thì lớp cỏ biến mất hoàn toàn và CẢ chiều cao ô (626 tới ~672, kể cả phần mặt cắt đứng) bị thay bằng một chất liệu sỏi sáng đồng nhất. Hai đường ranh là đường thẳng đứng tuyệt đối, không có ô pha trộn nhìn thấy được — dù bối cảnh ghi 'hai đầu có ô chuyển tiếp', trên ảnh không phân biệt được ô chuyển tiếp nào. Kết quả đúng là đọc ra 'rãnh sỏi đào xuống' chứ không phải mặt đường. Cũng xác nhận cả 5 pano chỉ có đúng đoạn lát này, phần còn lại của map không có đường nối các nhà. Giữ mức 'vừa'.</sub>

### 20. Không vật thể nào lún vào đất — tất cả bị cắt phẳng trên một đường kẻ

**O dau:** pano0.png x 370-620 (chum, vại, thùng gỗ), pano4.png x 400-560 (giếng) và x 545-635 (hàng rào) — đều bị xén tại y≈624-626.

Chum sành và thùng gỗ là vật tròn nhưng đáy bị xén phẳng bằng một đường thẳng tuyệt đối; chân giếng đá, chân cột hàng rào cũng vậy. Không có một ngọn cỏ, một hòn đá hay một vệt bóng nào phủ lên chân bất cứ vật gì. Vì thế cả làng trông như đồ chơi bày trên một cái kệ chứ không phải cắm xuống đất. Trong tranh concept thì cỏ, bụi cây, hàng rào đều mọc chồng lên chân nhà.

**Bang chung:** Ảnh cắt pano0 (370,590)-(620,650) phóng 6 lần: đáy chum, đáy vại và đáy thùng gỗ đều kết thúc đúng cùng một hàng y≈624.7 bằng một đường ngang phẳng, mất hẳn phần cong đáy. Ảnh cắt pano4 (400,590)-(640,650): chân giếng đá và chân 4 cột hàng rào cũng cụt tại y≈625, phía dưới là dải sáng rồi tới cỏ. Quét toàn ảnh: pixel trời thấp nhất nằm ở hàng 625 tại 874/1454 cột của pano0 — nghĩa là ở phần lớn bề ngang, trời nhìn thấy được tới sát mặt cỏ.

**Cach sua:** Vẽ thêm một lớp 'cỏ tiền cảnh' mỏng (bụi cỏ, hoa dại, đá nhỏ) vẽ SAU vật thể, rắc ngẫu nhiên dọc mặt đất và cho phủ lên 4-8px chân mỗi vật. Kèm theo đó thêm bóng đổ tối mềm dưới chân nhà/giếng/cây (lỗi này đã biết là chưa có) — hai thứ này cùng làm mới đủ để vật 'chạm đất'.

<sub>Lang kinh `tong-the` · LOI MOI · phan bien: CÓ THẬT về phần chính. Cắt pano0 (380,600)-(480,635) phóng 10 lần: đáy tròn của chum sành lớn, của vại nhỏ có nắp vải và của thùng gỗ đều bị xén cụt bằng một đường ngang phẳng tuyệt đối tại y≈625, mất hẳn đoạn cong đáy; giữa hai chum còn nhìn thấy sọc trời/mây sáng lọt qua trước khi tới mặt đất. Cắt (370,595)-(660,645) cũng cho thấy chân cột đèn lồng bị cắt y hệt. NHƯNG câu tuyệt đối 'không một ngọn cỏ nào phủ lên chân bất cứ vật gì' thì SAI: trong chính ảnh cắt đó, bụi cỏ + hoa trắng ở x≈370-392 có phủ lên mép trái chum lớn (chum bắt đầu ~x 389). Vậy có che nhưng lác đác, không thành lớp. Hạ nhẹ: giữ 'vừa'. Lưu ý gốc lỗi trùng với phát hiện 'nhà treo lơ lửng' — cùng là vật thể bị neo cao hơn mặt cỏ ~4px, nên sửa một lần là hết cả hai.</sub>

### 21. Sprite nhân vật lệch phải 7,5 đơn vị so với vị trí thật (cx) — bóng không nằm dưới thân

**O dau:** D:/NRO_Project/Assets/Scripts/Game1/Char.cs:5816 (vẽ bóng tại xSd=cx) và paintCharBody ~Char.cs:6057-6077 · part.bin part 57/58 (dx đều =0)

Bóng được vẽ đúng tâm cx (drawImage anchor 3 = HCENTER|VCENTER, xSd = cx), còn cả thân lẫn chân lại nằm lệch hẳn sang phải. Hộp va chạm là cx±11, nhưng hình vẽ chiếm từ cx−9 đến cx+23,5 — tức là hình thò ra khỏi hộp va chạm 12,5 đơn vị bên phải và hụt 2 đơn vị bên trái. Hậu quả nhìn thấy ngay: bóng nằm dưới bàn chân SAU chứ không dưới thân; nhân vật sẽ đâm hình vào tường bên phải trước khi bị chặn, và dừng cách tường bên trái một quãng. Đợt này đã sửa dy của thân/chân (−46/−33) nhưng dx vẫn để 0 cho mọi khung.

**Bang chung:** pano0.png và pano1.png cùng camera (cmx=24, cmy=119; khớp template render sai số 0.14) → trừ pixel tách được nhân vật khỏi nền. Bóng (các hàng y≥625, không dính chân): pano x 689..734 = đúng 46 px, bằng đúng bề rộng myTexture2dbong.png (46x14) → tâm world x = 379,75, tức cx≈380. Sprite (các hàng y<615): pano x 694..759 → world x 371,0..403,5, tâm 387,25. Lệch = +7,5 đơn vị. Xem scratchpad/charmask.png và charmask0.png.

**Cach sua:** Đặt dx âm cho part 57 (thân) và 58 (chân) trong part.bin sao cho tâm hình trùng cx — đo lại bằng đúng cách trên (chụp hai ảnh cùng camera, trừ pixel, so tâm bóng với bbox sprite). Nhớ tăng vsData sau khi đổi part.bin.

<sub>Lang kinh `loi-moi` · LOI MOI · phan bien: CON SỐ +7,5 ĐÚNG, nhưng cách đo của người báo cáo KHÔNG dùng được và tôi phải làm lại. pano0/pano1 gần như trùng nhau (chỉ 7529 px khác trên 1,25 triệu) — trừ hai ảnh đó chỉ ra khác biệt KHUNG HOẠT ẢNH, không tách được silhouette nhân vật. Tôi làm cách khác: pano2 và v3A lệch camera đúng 626 px (dò tương quan cho sai số 0.000 tuyệt đối) nên trừ được nền sạch. Kết quả đo: bóng chiếm pano x 669..714 = đúng 46 px (bằng myTexture2dbong.png 46x14) ⇒ tâm pano 692 ⇒ với cmx=24 (tôi tự xác nhận cmx=24 bằng mép gò ở world 552 → pano 1056) thì cx = 370 CHẴN. Kiểm chứng bằng suy luận mã: Char.cs:6083-6085 vẽ thân với x = cx + (CharInfo[0][2][1] + dx)*num2 và neo TOP_RIGHT/TOP_LEFT; CharInfo[0] = head[0,-13,34], leg[1,-8,10], body[1,-9,16]; ảnh thân 64 px = 32 đơn vị. Hướng phải (cdir==1, num2=+1, neo LEFT): thân chiếm [cx-9, cx+23] = world [361,393] → pano [674,738]; đo trên ảnh được đúng 674..740. Tâm hình = cx+7,5. Hộp va chạm cx±11 ⇒ hình thò phải 12 đơn vị, hụt trái 2 — khớp báo cáo. Nhìn thấy được: bóng (pano 669..714) KHÔNG phủ bàn chân trước (pano 720..738) — xem scratchpad/r_v3a_feet.png. Hạ từ 'nang' xuống 'vua': không phá gameplay, triệu chứng là lệch bóng và hình cấn tường, không phải kẹt.</sub>

### 22. Tile 65 và 66 (mép vách gò) là ảnh hỏng: bệt tối + sọc xám dọc

**O dau:** D:/NRO_Project/Assets/Resources/res/x2/t/1$65.png và 1$66.png · hiện ra ở 1.bytes hàng 18 cột 23, 27, 34, 41

Hai tile dùng làm góc vách gò không phải ảnh đất mà là ảnh hỏng. 1$65 có một thanh TỐI cứng 9 px ở mép trái; 1$66 có một thanh XÁM kiểu mã vạch 7 px ở mép phải. Cả thân hai tile cũng tối và bợt màu hơn hẳn tile đất 60/67 xung quanh. Trong game, mỗi đầu của hai cái gò hiện ra một khối chữ nhật tối khác hẳn kèm một vạch dọc sáng — 4 chỗ như vậy trên map. Khác với lỗi đã biết 'vách gò ghép từ tile lặp lòng đất, có đường cắt ngang': đây là hai tile riêng, hỏng theo chiều DỌC.

**Bang chung:** Đo từng cột pixel: 1$65 cột 0-8 có độ sáng TB 39 và chênh lệch kênh 25, cột 9-47 là 64/50. 1$66 cột 41-47 có độ sáng 76 và chênh 34, cột 0-40 là 54/51. Tile 60 và 67 đối chứng: 66/66 và 62/63 đều nhau suốt 48 cột. Trong game: pano2.png vùng (622,625)-(670,674) = hàng 18 cột 27 — khối tối lạ + sọc xám ở mép phải (ảnh cắt scratchpad/z_seam.png); pano3.png x≈1018 y≈630-670 cũng vậy.

**Cach sua:** Vẽ lại 1$65 / 1$66 thành mép vách đất tử tế, hoặc tạm thời thay hai ô góc bằng tile 60/61 rồi vẽ mép sau. Kiểm nhanh bằng cách so độ sáng và độ chênh kênh trung bình của tile mới với tile 60/67 — lệch quá ~10% là còn thấy được bằng mắt.

<sub>Lang kinh `loi-moi` · LOI MOI · phan bien: ĐÚNG, và nhìn bằng mắt còn tệ hơn số đo. Tôi đo lại độ sáng trung bình từng dải cột (x2, 48x48): 1$65 cột 0-8 = 40 so với thân 64; 1$66 cột 41-47 = 76 so với thân 55. Đối chứng 1$60 = 66/67/66, 1$67 = 62/63/63, 1$68 = 58/59/58 — đều tăm tắp suốt 48 cột. Nhìn ảnh (scratchpad/r_tiles.png): 65 và 66 không phải cùng chất liệu với đất 60/61/67/68 — chúng là vân đá xếp lớp ngang, 65 có thanh tối cứng bên trái, 66 có dải xám kiểu MÃ VẠCH bên phải. TRONG GAME tôi cắt pano3.png vùng (950,615)-(1060,685) = hàng 18 cột 41 (scratchpad/r_p3_tile66.png): hiện ra đúng một khối chữ nhật tối vân ngang kèm một dải sọc trắng-xám dọc rõ mồn một, chỏi hẳn với đất xung quanh. Có 4 chỗ như vậy (cột 23, 27, 34, 41). Giữ 'vua' (chỉ 4 ô, nhưng cực chướng mắt).</sub>

### 23. Sửa ô chân dung HUD KHÔNG có tác dụng — bảng ánh xạ head→avatar không nuôi ô đó

**O dau:** D:/NRO_GoiBanGiao/server/internal/game/thegioi.go:86-89 và 108-110 · D:/NRO_Project/Assets/Scripts/Game1/Char.cs:1482 (getAvatar) · Panel.cs:9759 …

Máy chủ đã gửi cặp head 30 → avatar 27078, nhưng ô chân dung trên HUD góc trái vẫn trống y hệt trước khi sửa. Lý do: Char.avatarz() / Char.idAvatar CHỈ được dùng ở Panel.cs (bảng thông tin nhân vật mở từ menu) và ChatPopup.cs — không có chỗ nào trong HUD góc trái đọc bảng này. Ngoài ra bảng chỉ có 1 cặp trong khi CreateCharScr.hairIDOld cho phép 9 giá trị kiểu tóc (64,30,31,9,29,32,6,27,28) → ai chọn kiểu tóc khác vẫn nhận -1. Test avatar_test.go tự dựng NhanVat{KieuToc:30} nên không thể bắt được cả hai vấn đề.

**Bang chung:** So pixel vùng ô chân dung (x 10-55, y 10-60) giữa pano5.png (chụp 14:21, TRƯỚC khi sửa) và v3A.png (SAU khi sửa): 0 pixel khác nhau. grep 'idAvatar|avatarz()' trong Assets/Scripts/Game1 chỉ ra Char.cs (định nghĩa), Controller.cs:5980-5984 (đọc gói) và Panel.cs — không có GameScr/InfoMe/HUD. Ảnh Small27078.png là chân dung đẹp, đầy đủ, tức lỗi không nằm ở ảnh.

**Cach sua:** Trước hết xác định ĐÚNG hàm vẽ ô chân dung trên HUD (mở bảng thông tin nhân vật trong game xem chân dung có hiện không — nếu hiện thì bảng ánh xạ đã chạy và ô HUD là chuyện khác hẳn). Song song, gửi đủ 9 cặp headId→27078 cho cả 9 kiểu tóc, và đổi test để dựng nhân vật bằng đúng đường tạo nhân vật thật thay vì gán KieuToc: 30 bằng tay.

<sub>Lang kinh `loi-moi` · LOI MOI · phan bien: ĐÚNG, và tôi kiểm được cả 'ảnh sau đúng là ảnh sau'. grep 'idAvatar|avatarz()' trong Assets/Scripts/Game1 chỉ ra: Char.cs:367/1477/1490 (định nghĩa), Controller.cs:5978-5985 (đọc gói -30), và toàn bộ chỗ dùng còn lại nằm ở Panel.cs — không có GameScr/InfoMe/HUD nào đọc. Tôi tự so pixel vùng (10..55, 10..60) giữa pano5.png (14:21) và v3A.png (18:42): 0 pixel khác. Quan trọng: đợt sửa CÓ hiệu lực thật trong bản chụp v3A vì thanh máu/chakra đổi từ 0/0 (pano5) sang 150/150 (v3A) — cùng commit ca44983 (18:32) gửi cả -42 lẫn bảng ánh xạ. Vậy bảng ánh xạ tới nơi mà ô chân dung vẫn trống trơn (phóng to scratchpad/r_hud_ava.png: chỉ có dấu placeholder xanh). Xác nhận nốt: CreateCharScr.cs:28-33 hairIDOld = {64,30,31},{9,29,32},{6,27,28} = đúng 9 kiểu tóc, mà máy chủ chỉ gửi 1 cặp (thegioi.go: w.Short(1)). Bố cục gói khớp với Controller.cs nên không phải lỗi lệch byte. Giữ 'vua'.</sub>


## Nhe

### 24. Ảnh chân dung của người chơi trùng đúng ảnh chân dung NPC (thầy Iruka)

**O dau:** D:\NRO_GoiBanGiao\server\internal\game\thegioi.go:110 (AvatarNhanVat = 27078) so với D:\NRO_GoiBanGiao\server\internal\game\assets\catalog.json (npcs[0].portrait = 27078). Tên sprite trong D:\NRO_GoiBanGiao\docs\my-thuat\mapping-id.json: id 27078 = "iruka-portrait".

Bản vá mục 8 chọn ảnh 27078 làm chân dung người chơi, nhưng 27078 chính là chân dung NPC đã đăng ký trong catalog và trong bảng mapping tên nó là "iruka-portrait". Khi chân dung thật sự hiện ra (menu nhân vật, hoặc HUD sau khi sửa), người chơi sẽ thấy khuôn mặt thầy Iruka làm ảnh đại diện của chính mình — trùng hệt mặt NPC đang đứng trong làng. Ngoài ra ánh xạ cứng head = 30 chỉ đúng vì kiểu tóc mặc định đang là 30; người chơi chọn kiểu tóc khác là getAvatar trả -1 và ô chân dung lại trống.

**Bang chung:** catalog.json npcs[0] = {head 2189, body 2190, leg 2191, portrait 27078, x 172, y 384}. mapping-id.json sprites: {id 27078, name 'iruka-portrait'}. thegioi.go: const AvatarNhanVat int16 = 27078. Xem thẳng ảnh res\x2\smallimage\Small27078.png: một ninja tóc nâu buộc đuôi, băng trán Lá, áo giáp jounin — đúng tạo hình Iruka, khác hẳn nhân vật người chơi trong game (tóc đen, khăn xanh).

**Cach sua:** Vẽ/cắt một ảnh chân dung riêng cho nhân vật người chơi (id còn trống trong dải 27041–27082 đã dành sẵn), rồi trỏ AvatarNhanVat vào ảnh đó. Và thay ánh xạ cứng head=30 bằng ánh xạ theo đúng kiểu tóc nhân vật đang có (gửi cặp cho mọi kiểu tóc hợp lệ).

<sub>Lang kinh `kiem-sua` · LOI MOI · phan bien: DỮ LIỆU ĐÚNG, xác nhận từng chỗ: thegioi.go:110 const AvatarNhanVat int16 = 27078; assets/catalog.json npcs[0].portrait = 27078; mapping-id.json dòng 1977-1979 ghi {id 27078, name 'iruka-portrait'} với source_path = docs/my-thuat/anh-nguon/vertical-slice-v1/iruka_portrait_chroma.png. Mở thẳng res/x2/smallimage/Small27078.png (100x100): ninja tóc nâu buộc đuôi, băng trán Lá, áo giáp jounin — khác hẳn nhân vật người chơi trong ảnh chụp (tóc đen, khăn sẫm). Vậy người chơi và NPC Iruka đang dùng chung một ảnh đại diện. HẠ từ 'vua' xuống 'nhe': hiện KHÔNG có triệu chứng nào trên màn hình — HUD không vẽ chân dung (xem mục trên) và NPC chưa được sinh ra (đã nằm trong danh sách lỗi đã biết) — đây là lỗi dữ liệu tiềm ẩn, chỉ lộ khi mở menu nhân vật hoặc khi NPC chạy. Ý phụ về 'head=30 chỉ đúng với kiểu tóc mặc định' hợp lý nhưng chưa quan sát được vì chưa có hệ thống chọn tóc.</sub>

### 25. Nền: một đường cắt ngang sắc nét chạy hết bề ngang màn hình ở y=344

**O dau:** Có ở mọi ảnh, ví dụ D:\NRO_Project\pano0.png hàng y=344, chạy từ x=0 tới x=1453 (kiểm ở x=400, 500, 600, 1100 đều giống nhau). Lớp nền liên quan: D:\NRO_Project\Assets\Resources\res\x2\bg\b12.png. Ảnh cắt: ...\scratchpad\seam344.png

Chân dãy núi Hokage tan dần thành sương rồi ĐỘT NGỘT biến mất tại đúng một hàng pixel: từ y=343 trở lên còn vệt sương xám (khoảng (25,90,150)), từ y=344 trở xuống là màu xanh phẳng lì (10,86,155) không còn chút vân nào. Vì đường cắt nằm ở toạ độ MÀN HÌNH cố định (không đổi khi cuộn camera) nên nó luôn hiện ra như một đường kẻ ngang vắt ngang cả khung hình. Khác với lỗi đã biết (ảnh nền 384px ghép đôi → mối nối DỌC): đây là mối nối NGANG.

**Bang chung:** Quét chênh lệch màu trung bình giữa hai hàng liên tiếp trên dải x 200–1400 của pano0: nhảy vọt ở hàng 344. Đọc pixel: x=400 → y343 (25,91,151), y344 (10,86,155) rồi giữ nguyên (10,86,155) xuống dưới; x=500, 600, 1100 cho kết quả y hệt. Bản thân ảnh b12.png thì chuyển mượt (độ lệch chuẩn theo hàng giảm dần 44 → 0 từ hàng 255 tới hàng 320, không có bậc), nên bậc thang này sinh ra ở khâu ghép lớp lúc vẽ chứ không phải trong ảnh gốc. Đường cắt đứng yên ở y=344 trong cả pano0, pano2, pano3 (camera khác nhau).

**Cach sua:** Rà lại toạ độ y và chiều cao khi vẽ lớp b12 (và b11/b13) — nhiều khả năng lớp bị vẽ chồng thêm một lần hoặc bị cắt cụt ở đúng hàng đó. Cách chữa nhanh: kéo dài vùng chuyển màu của b12 xuống hết ảnh hoặc cho lớp nền trên tràn qua chỗ nối.

<sub>Lang kinh `kiem-sua` · LOI MOI · phan bien: CÓ THẬT. Đọc pixel pano0: x=400 y343 (25,91,151) -> y344 (10,86,155); x=500 (19,86,147)->(10,86,155); x=600 (25,91,150)->(10,86,155); x=1100 (21,88,149)->(10,86,155). Trên y=343 còn vân sương lăn tăn, từ y=344 xuống là gradient trời mượt hết vân. Đường này ở đúng y=344 trong pano0, pano2 và pano3 dù cmx khác nhau (24 / 337 / 497) nên nó đứng yên theo màn hình. Kiểm thêm ảnh nguồn b12.png (768x384): phần dưới chuyển mượt rồi thành xanh phẳng (10,86,154) từ khoảng hàng 320, không có bậc — bậc thang sinh ra ở khâu ghép lớp lúc vẽ, đúng như báo cáo. Crop x 300-1100 y 320-370 nhìn mắt thường thấy được vạch; tăng tương phản 4 lần thì rất rõ. Giữ 'nhe': biên độ chỉ ~15/255 ở kênh đỏ. Không trùng lỗi đã biết (lỗi đã biết là đường nối DỌC do ảnh 384px ghép đôi — đường nối dọc đó tôi cũng thấy ngay trong chính b12.png tại x=384).</sub>

### 26. NPC và điểm hồi sinh đặt ở y=384, cao hơn mặt đất đúng 2 ô

**O dau:** D:\NRO_GoiBanGiao\server\internal\game\assets\catalog.json — npcs[0] có x=172, y=384. D:\NRO_GoiBanGiao\server\internal\game\lang.go:27–31 — cả 5 làng đều ViTriYKhoiDau = 384. Mặt đất map 1 là hàng 18 → y = 432.

Trên map 1, mặt cỏ nằm ở hàng 18 nên chân nhân vật/NPC phải ở y=432 (đã kiểm: nhân vật trong game đứng đúng y=432). Nhưng NPC được đặt ở y=384 và điểm khởi đầu của cả 5 làng cũng là y=384 — cao hơn mặt đất 48 đơn vị = 2 ô = 96 pixel màn hình. Client còn cộng thêm +3 khi dựng NPC (Controller.cs:5155) nên NPC sẽ lơ lửng giữa trời. Với người chơi thì mỗi lần vào game nhân vật xuất hiện trên không rồi rơi xuống. Ghi chú trong banDo.go:18 nói rõ ý định là "nằm trên nền đất của map" — vậy đây là số sai chứ không phải cố ý.

**Bang chung:** 1.bytes hàng 18 cột 7 = tile 62 (mặt cỏ), hàng 17 và 16 cột 7 = 0 (trống) → mặt đất ở cột 7 là y=432. NPC x=172 rơi vào cột 7. 384/24 = hàng 16 → cách mặt đất 2 hàng. Kiểm chéo: catalog_test.go:51 khẳng định DiDuoc(100, 384) = true, tức y=384 là KHÔNG KHÍ đi được, không phải mặt đất. Trong game ở v3A vị trí NPC sẽ rơi vào màn hình x≈296, y≈530 — chỗ đó chỉ có trời (NPC hiện chưa vẽ ra, đây là lỗi tiềm ẩn sẽ lộ ngay khi NPC chạy).

**Cach sua:** Đổi npcs[0].y trong catalog.json thành 429 (để client cộng +3 ra đúng 432) và đổi ViTriYKhoiDau của làng Lá thành 432; các làng 2–5 sửa theo đúng hàng mặt đất của map tương ứng khi có map. Nên thêm một test khẳng định mọi toạ độ trong catalog/lang đều nằm ngay trên ô đặc của map-1-collision.bin.

<sub>Lang kinh `kiem-sua` · LOI MOI · phan bien: CÓ THẬT, số liệu kiểm được hết. catalog.json: npcs[0] = {x:172, y:384}, cả 6 map spawn {x:100, y:384}; lang.go:27-31 cả 5 làng ViTriYKhoiDau = 384. Đọc res/mymap/1.bytes: hàng 18 cột 7 (x=172 rơi vào cột 7) = tile 62 (mặt cỏ), hàng 17 và 16 cột 7 = 0 → mặt đất ở đó là y=432. 384 = hàng 16, cách mặt đất đúng 48 đơn vị = 2 ô. Đo trong ảnh chụp cũng khớp: cmy=119 ở mọi ảnh, đỉnh hàng 18 ở màn y=626, chân nhân vật kết thúc y=624 → nhân vật thực tế đứng ở cy=432 chứ không phải 384. Controller.cs:5155 dựng NPC với cy = y+3 nên NPC sẽ ở 387, vẫn lơ lửng 45 đơn vị; chiếu sang v3A (cmx=24, cmy=119) thì NPC rơi vào màn x=296, y=536 — đúng chỗ trống trên trời. Giữ 'nhe': với người chơi chỉ là rơi nhẹ 2 ô lúc vào game (client có trọng lực, ảnh chụp cho thấy đã tiếp đất đúng); với NPC là lỗi tiềm ẩn vì NPC chưa được vẽ ra. Ghi chú banDo.go:18 nói SpawnY phải 'nằm trên nền đất của map' nên đây đúng là số sai chứ không phải cố ý.</sub>

### 27. Chân dung 27078 và sprite trong game không phải cùng một nhân vật (tóc nâu + băng trán vs tóc xanh đen không băng trán)

**O dau:** D:\NRO_Project\Assets\Resources\res\x2\smallimage\Small27078.png (chân dung) so với Small27042.png (part đầu) và Small27045.png (thân) — thể hiện trong v3A.png vùng x 668-745, y 532-570.

Ảnh chân dung dùng cho HUD vẽ một ninja tóc NÂU buộc đuôi, đeo băng trán Konoha. Sprite chạy trong game lại tóc XANH ĐEN và (sau khi part đầu bị blank) không có băng trán. Kể cả khi khôi phục part đầu thì tóc vẫn là đen chứ không nâu. Người chơi nhìn chân dung một đằng, nhìn nhân vật một nẻo — hỏng tính nhất quán của nhân vật chính.

**Bang chung:** Lấy màu nổi trội: vùng tóc của Small27078 (x 20-70, y 10-40) = (90,57,27), (53,33,13), (89,57,28) — nâu ấm. Vùng tóc của sprite thân Small27045 (x 2-22, y 18-34) = (15,25,37), (36,50,70), (35,49,72) — xanh navy đen. Băng trán chỉ có ở 27078 và ở part đầu 27041/42/43, không có ở thân.

**Cach sua:** Chọn một bản thiết kế duy nhất rồi chỉnh bên kia theo: hoặc đổi màu tóc chân dung 27078 sang xanh đen, hoặc đổi bảng màu tóc của bộ sprite 27041-27074 sang nâu. Và bắt buộc khôi phục part đầu để sprite có băng trán như chân dung.

<sub>Lang kinh `nhan-vat` · LOI MOI · phan bien: CÓ THẬT nhưng người báo cáo vừa sai bằng chứng vừa thổi mức độ. Đo màu tóc: Small27078 vùng (20-75, 10-45) = (90,57,27)/(53,33,13) — nâu ấm, có đuôi tóc buộc; part đầu Small27042 vùng (10-60, 12-40) = (54,50,46)/(57,53,48) — xám đen trung tính, tóc gai. Khác thật. NHƯNG: (a) người báo cáo lấy màu 'tóc thân' ở Small27045 (2-22, 18-34) ra (15,25,37)/(36,50,70) và gọi là 'xanh navy đen' — vùng đó là DẢI BĂNG TRÁN bay ngược, không phải tóc, nên bằng chứng đó sai; (b) quần áo hai bên khớp hoàn toàn (áo giáp xanh rêu, áo lót navy, băng trán navy có phù hiệu lá) nên rõ ràng cùng một thiết kế, chỉ lệch sắc tóc; (c) quan trọng nhất — chân dung HIỆN KHÔNG ĐƯỢC VẼ RA (mục trên), nên người chơi không hề thấy sự bất nhất này. Đây là nợ tiềm ẩn, hạ từ vua xuống nhe.</sub>

### 28. Bàn chân sau hở trời 2-4px so với mép cỏ — nhân vật đọc thành lơ lửng

**O dau:** D:\NRO_Project\v3A.png bàn chân sau x 719-737, đáy y 622-624; mép cỏ y=626. Nhìn rõ ở scratchpad\ground_v3A.png (crop 640,608-810,645 tăng sáng, phóng 8x) và scratchpad\shadow_zoom_pano1.png

Ảnh chân Small27062 vẽ bàn chân sau cao hơn bàn chân trước 3px (thủ pháp tạo chiều sâu, bản thân nó không sai). Nhưng vì mặt đất ở đây là một vạch ngang sắc lẻm có TRỜI XANH SÁNG ngay bên trên, và vì vệt bóng lại lệch trái không với tới bàn chân sau, nên khe 2-4px đó hiện thành một dải trời xanh rõ mồn một dưới gót chân sau. Kết quả: nhân vật trông như đang nhấc một chân lơ lửng dù đang đứng yên.

**Bang chung:** Quét đáy từng cột của Small27062.png: cột 12-25 (chân trước) đáy = 62-63; cột 43-61 (chân sau) đáy = 58-60. Ảnh chân được vẽ tại y=564 trong v3A → gót trước ở y 626-627 (trùng mép cỏ 626, đúng), gót sau ở y 622-624 (cao hơn mép cỏ 2-4px). Ảnh tăng sáng cho thấy dải màu trời (100,188,234) chen giữa đế dép sau và mép cỏ.

**Cach sua:** Sửa xong lỗi bóng lệch trái (mục 3) là khe này gần như tự lành, vì bóng phủ tới chân sau sẽ lấp dải trời. Nếu vẫn thấy hở, hạ nhân vật thêm 2px hoặc nới ellipse bóng cao thêm 2 hàng về phía trên.

<sub>Lang kinh `nhan-vat` · LOI MOI · phan bien: XÁC NHẬN, hơi bị nói quá về số. Đáy từng cột của Small27062 vẽ tại (676,564): cột 43-49 → đáy màn hình y=622, cột 50-54 → y=623, cột 55-57 → y=624. Mép cỏ y=626. Đọc pixel v3A.png dưới gót sau: x=723 y623=(102,189,234), y624=(104,191,236), y625=(104,191,235) — trời nguyên chất, rồi y626 mới là cỏ. Vậy khe hở thật là 1-3 hàng chứ không phải 2-4. Chân trước (cột 12-25, đáy 62-63 → y 626-627) chạm đúng mép cỏ. Lỗi có thật nhưng chỉ 1-3px, và ghi chú thêm: NGÔI NHÀ ở pano0 cũng hở trời tương tự ở đáy (pano0 crop 100,590-400,650) nên đây có thể là vấn đề chung của tầng đất chứ không riêng nhân vật.</sub>

### 29. Cả làng xếp thẳng hàng, không vật nào che vật nào — thiếu chiều sâu

**O dau:** Toàn map (ảnh ghép pano0 + pano4, rộng 2400px). Các cụm vật thể: x 37-374 (nhà 1), 389-446 (chum thùng), 468-627 (cột đèn), 676-958 (cây), 1042-1285 (cổng torii), 1296-1472 (giếng), 1587-1970 (nhà lớn), 2038-2386 (quán mì).

Ghép cả map lại nhìn một lượt: 8 cụm vật thể nằm cạnh nhau như 8 món hàng trên kệ, mỗi món có khoảng trống bao quanh, không món nào chồng lên món nào, tất cả cùng một tầng sâu. Mắt không có gì để bám, không biết nên nhìn vào đâu. Tranh concept thì nhà chồng nhà, hàng rào và bụi cây nằm trước nhà, có nhà lùi ra sau mờ đi — đó là thứ tạo cảm giác 'một ngôi làng đông đúc'.

**Bang chung:** Quét cột dọc trên ảnh ghép trong dải y 350-570 (dải chứa thân nhà, không tính trời và đất): chỉ 14% bề ngang là trống, nhưng KHÔNG có cụm nào giao nhau — mọi cụm cách nhau bởi khoảng trời hoặc chạm nhau đúng ở mép, không có vùng chồng lấp. Cũng không có vật thể nào nằm ở tầng gần (tiền cảnh) hay tầng xa mờ.

**Cach sua:** Thêm một lớp tiền cảnh và một lớp trung cảnh: vài mái nhà nhỏ hơn, mờ hơn, nhô lên phía sau giữa các nhà lớn; vài bụi tre / cột đèn / góc mái vẽ ở tiền cảnh phủ lên 1/4 màn hình bên trái hoặc phải. Cho ít nhất 3 cặp vật thể chồng mép lên nhau (hàng rào chạy cắt ngang chân nhà, cây che một góc mái).

<sub>Lang kinh `tong-the` · LOI MOI · phan bien: Quan sát ĐÚNG nhưng người báo cáo thổi mức. Tôi xem lần lượt pano0-pano4 ở 1:1: nhà đầu map, cụm chum-thùng, cột đèn, cây, cổng torii, giếng, nhà lớn, quán mì — mỗi cụm đứng riêng, cách nhau bằng khoảng trời, không cụm nào chồng mép lên cụm nào, tất cả cùng một tầng sâu, không có lớp tiền cảnh. So concept art thì concept có lá cây tiền cảnh góc trên trái, cỏ silhouette đáy khung, nhà chồng nhà và nhà lùi ra sau mờ đi. NHƯNG tôi bác một ý: 'cũng không có vật thể nào ở tầng xa mờ' là SAI — lớp làng xa b01 có tồn tại (mái chùa, mái ngói, cây, núi nhỏ nhìn thấy ở y 530-561 trong pano0/pano2), chỉ có điều nó đang bị cắt cụt bởi lỗi đường ngang y=561. Hạ từ 'vừa' xuống 'nhẹ': đây là danh sách mong muốn về mỹ thuật (cần vẽ thêm asset), không phải thứ đang hỏng, và trùng phần lớn với hai mục đã biết 'cả map chỉ có 1 cây' + 'không có bóng đổ dưới nhà/giếng/cây'.</sub>

### 30. Địa hình 50 ô chỉ có 2 bậc cao 1 ô, còn lại phẳng tuyệt đối

**O dau:** Toàn map: mặt đất nền ở y=626 suốt 50 ô; chỉ 2 gò cao ở y=578 — ô 22-27 (rộng 5 ô, dưới cổng torii) và ô 33-41 (rộng 8 ô, dưới nhà lớn).

Nhìn ra thì cả ngôi làng nằm trên một mặt phẳng dài 50 ô, đứt đúng 2 chỗ bởi hai bậc cao vỏn vẹn 1 ô. Với người chơi mới, màn chơi không gợi ra chỗ nào đáng leo, đáng nhảy, đáng khám phá — đi từ đầu tới cuối là hết. Cộng với việc chưa có bậc thang (lỗi đã biết), hai landmark đẹp nhất là cổng torii và nhà lớn lại đứng trên bệ mà người chơi không thấy cách nào lên.

**Bang chung:** Dò mép cỏ theo từng cột trên ảnh ghép 2400px: chỉ tồn tại đúng 2 mức độ cao là y=578 và y=626 (chênh 48px = 1 ô). Các đoạn ở mức 578 là x 1056-1295 (ô 22.0-27.0) và x 1584-1967 (ô 33.0-41.0); toàn bộ phần còn lại ở mức 626.

**Cach sua:** Thêm ít nhất 2-3 mức cao nữa: một mái nhà nhảy được lên, một bậc thềm 2 ô cao dẫn tới cổng torii, một hố/rãnh nhỏ cần nhảy qua. Ưu tiên làm bậc thang lên hai gò hiện có trước — vừa sửa lỗi đã biết vừa lập tức làm địa hình đọc ra 'có đường đi'.

<sub>Lang kinh `tong-the` · LOI MOI · phan bien: CÓ THẬT. Xem cả 5 pano: mặt cỏ nền luôn ở y=626, chỉ có hai đoạn ở y=578 (chênh đúng 48px = 1 ô). Quy về toạ độ thế giới bằng độ lệch camera đã đo: gò 1 ở pano2 x 430-670 với camera 626 → thế giới 1056-1296px → ô 22.0-27.0; gò 2 ở pano4 x 640-1030 với camera 946 → thế giới 1586-1976px → ô 33.0-41.2. Khớp với bối cảnh (cột 23-27 và 34-41). Ngoài hai gò đó không còn mức cao nào khác. Giữ mức 'nhẹ' như người báo cáo: đây là nhận xét thiết kế màn chơi, và nó dính chặt với lỗi đã biết 'không có bậc thang' — bản thân việc thiếu bậc thang mới là cái phải sửa trước.</sub>

### 31. HUD bị xén mất mép trên, số chakra xanh trên nền xanh, và một vạch sáng ngang chạy hết màn ở y=34

**O dau:** Góc trên trái mọi ảnh pano: khung HUD x 0-310, y 0-70; số '150' chakra ở x 210-255, y 38-52; vạch ngang tại y=33→34 chạy hết 1454px.

Ba chuyện nhỏ ở giao diện nhưng làm phần trên màn hình trông cẩu thả. Một là khung viền HUD bị màn hình cắt cụt ở mép trên, thanh máu không có viền trên. Hai là số chakra màu xanh lơ nằm trên chính thanh chakra màu xanh lơ nên gần như không đọc được (số máu màu hồng trên nền cam cũng chìm). Ba là có một vạch sáng ngang chạy suốt bề ngang màn hình ở y=34, cắt cả bầu trời làm hai.

**Bang chung:** Ảnh cắt (150,0)-(400,80) phóng 4 lần: khung cam của thanh máu bắt đầu ngay tại y=0 bằng một nhát cắt thẳng, không có bo góc trên. Số '150' chakra và thanh chakra cùng tông xanh lơ. Đo độ chênh hai hàng liên tiếp trên toàn bề ngang pano0: hàng 33→34 = 34.9 (mạnh thứ ba cả ảnh); lấy màu tại x=1300 thấy [151,189,162] ở y=33 nhảy xuống [130,159,137] ở y=34 — tức có một dải phủ sáng cao 34px nằm đè lên đỉnh màn hình.

**Cach sua:** Dời cả cụm HUD xuống 6-8px cho lọt viền. Đổi chữ số trên thanh chakra sang trắng có viền đen (và số máu sang trắng viền đen) để đọc được ở mọi nền. Kiểm lại dải phủ ở y 0-33: nếu là hiệu ứng ánh sáng thì cho nó nhoè dần (gradient) thay vì cắt cứng, nếu là lớp thừa thì bỏ hẳn.

<sub>Lang kinh `tong-the` · LOI MOI · phan bien: ĐÚNG 2 trong 3 ý, ý còn lại tôi bác. (1) Vạch ngang y=33→34: CÓ THẬT và mạnh — độ chênh hai hàng liên tiếp trên toàn bề ngang là 34.91 (pano0), 38.07 (pano2), 41.24 (pano4), và 79.6% số cột chênh >10 ở cả ba ảnh. Ảnh cắt (1100,0)-(1454,60) phóng 3 lần cho thấy mép cứng: phía trên y=34 là trời sáng có mây, phía dưới là một lớp phủ tối bệt. (2) Số liệu tương phản: CÓ THẬT — cắt (0,0)-(320,80) phóng 4 lần thấy số '150' chakra màu xanh lơ nằm trên chính thanh chakra xanh lơ, và số '150' máu màu hồng nhạt trên thanh cam, cả hai đều chìm. (3) 'Khung HUD bị xén mép trên, thanh máu không có viền trên': SAI. Đo hàng 0 tại x=60/150/250 ra [30,30,30] (viền đen) và hàng 2 ra [109,76,32] (ruột cam) — viền TRÊN có tồn tại, chỉ là khung dán sát y=0 không chừa lề; và trong ảnh phóng to thanh máu cam có viền đen đủ bốn phía. Hạ xuống 'nhẹ' vì hai ý đúng đều là chuyện chỉnh màu chữ và một dải phủ, không ảnh hưởng chơi.</sub>

### 32. Bù bóng 12 đơn vị quá tay — nửa cái bóng nổi lên trời

**O dau:** D:/NRO_Project/Assets/Scripts/Game1/Char.cs:997 (BU_BONG_Y = 12) và Char.cs:5816

Trước đây bóng chìm dưới đất, đợt này nâng lên 12 đơn vị. Nay bóng nằm cưỡi lên đường cỏ: đúng một nửa chiều cao bóng (7 px ảnh = 3,5 đơn vị) nằm TRÊN mặt cỏ và được vẽ đè lên nền trời, thành một vệt elip xám nổi giữa trời xanh, thấy rõ khoảng hở giữa hai chân. paintShadow chỉ setClip trong 3 nhánh đặc biệt (ô bên cạnh rỗng / có bit 8); đứng trên mặt đất phẳng thì không nhánh nào chạy nên bóng không bị cắt.

**Bang chung:** myTexture2dbong.png (x2) = 46x14 px = 23x7 đơn vị. Trừ pano0-pano1: hàng cuối cùng của bóng ở pano y=632, cao 14 px → đỉnh bóng ở pano y=619. Mặt cỏ (world 432) ở pano y=626. Vậy 7 px ảnh (y 619..625) nằm trên mặt cỏ. Ảnh so sánh trước/sau cùng vùng: scratchpad/shadow_cmp.png (pano0 phía trên, pano1 phía dưới, vùng (680,595)-(780,640)).

**Cach sua:** Hạ BU_BONG_Y xuống còn ~7-8, HOẶC (chắc hơn) setClip theo hàng ô mặt đất trước khi drawImage bóng để bóng không bao giờ vẽ được lên phía trên bề mặt. Nên truy nốt nguyên nhân gốc của con số bù thực nghiệm này — comment trong Char.cs:5809 tự nhận là chưa truy ra.

<sub>Lang kinh `loi-moi` · LOI MOI · phan bien: HIỆN TƯỢNG CÓ THẬT, NHƯNG CHẨN ĐOÁN SAI — đây là chỗ tôi bác lại. Tôi đo trên v3A.png (trừ nền bằng pano2, xem trên): bóng chiếm pano y 619..632, mặt cỏ (world 432, hàng 18) ở pano y 626 ⇒ đúng 7 px ảnh (3,5 đơn vị) nằm trên mặt cỏ, vẽ đè lên trời/bụi cỏ. Ảnh phóng to scratchpad/r_shadow_zoom.png cho thấy rõ vệt xám tràn lên trên đường cỏ. Đến đây báo cáo đúng. NHƯNG: Char.cs:5816 dùng anchor 3 = HCENTER|VCENTER, tức THIẾT KẾ GỐC là đặt TÂM bóng đúng tại cy. Sau khi bù 12, tâm bóng đo được đúng ở world y≈432 = đúng chỗ thiết kế muốn. Vậy 12 KHÔNG 'quá tay' — nó đưa bóng về đúng vị trí danh nghĩa; cái nửa trên tràn lên trời là hệ quả của anchor VCENTER cộng với mặt đất bị cắt ngang cứng, chứ không phải bù dư. Cách sửa đúng là setClip theo hàng ô mặt đất (đề xuất thứ hai của báo cáo), không phải hạ BU_BONG_Y xuống 7-8 (làm thế là kéo bóng chìm lại). Hạ mức từ 'vua' xuống 'nhe'.</sub>

### 33. Vệt cắt ngang toàn màn hình ở nền, tại pano y=344

**O dau:** Mọi ảnh chụp: pano0-pano4.png, v3A.png, hàng y=343 → y=344, suốt chiều ngang 0..1453

Có một đường ngang cứng chạy hết bề rộng màn hình, phía trên là dải trời hơi xám có vệt sọc dọc (đáy lớp núi Hokage tan dần), phía dưới là màu trời phẳng lì. Đây là mép dưới của lớp nền, không hoà vào lớp trời sau nó. Khác với lỗi đã biết (đường nối DỌC do ảnh 384px ghép đôi) — đây là đường NGANG.

**Bang chung:** Đo màu trung bình từng hàng pano0.png: y=343 = (29,2 / 88,7 / 133,3), y=344 = (19,9 / 86,7 / 137,5), và bên dưới giữ nguyên (10,86,155) suốt nhiều chục hàng. Lấy mẫu tại x = 150, 300, 600, 1100, 1400 đều nhảy đúng ở y=344, chênh kênh đỏ 11-19. Ảnh cắt: scratchpad/seam_h.png (dải y 320-370).

**Cach sua:** Kéo dài phần tan dần ở đáy b02/b12 xuống thêm, hoặc cho màu hàng cuối của lớp này trùng đúng màu lớp trời phía sau ở đúng độ cao đó (hiện lớp trời ở đó là (10,86,155)).

<sub>Lang kinh `loi-moi` · LOI MOI · phan bien: ĐÚNG. Tôi đo màu trung bình từng hàng: pano0.png y=343 = (29,2 / 88,7 / 133,3) → y=344 = (19,9 / 86,7 / 137,5), rồi giữ phẳng nhiều chục hàng. pano2.png (camera khác hẳn) nhảy ở ĐÚNG y=344 luôn: (42,2/80,0/107,0) → (36,6/79,6/110,6) ⇒ không phải trùng hợp theo nội dung mà là mép LỚP, cố định theo màn hình. v3A giống hệt pano0. Cắt dải y 310-380 (scratchpad/r_seam_h.png) thấy rõ một đường ngang cứng chạy hết 1454 px: trên là đáy lớp núi Hokage tan dần, dưới là trời phẳng. Khác lỗi đã biết (đường nối DỌC do ảnh 384 px ghép đôi). Đáy b12.png là (10,86,154) đồng nhất — khớp con số báo cáo nêu. Giữ 'nhe'.</sub>

### 34. Bàn chân trước hở đất khoảng 2 đơn vị

**O dau:** Nhân vật ở tư thế đứng, pano1.png vùng (690,600)-(790,640); part.bin part 58 (chân) dy=-33

Chân sau chạm đúng mặt cỏ, chân trước lơ lửng cách mặt cỏ khoảng 2 đơn vị (4 px ảnh). Nhìn kỹ thấy bàn chân trước không dính đất. Một phần là do bộ sprite chân vẽ hai bàn lệch nhau ~0,75 đơn vị, phần còn lại là do làm tròn dy.

**Bang chung:** pano1.png: đế chân sau ở pano y≈627 → world y 432,4 (mặt cỏ 432, khớp). Đế chân trước ở pano y≈623 → world y 430,3. Ảnh cắt scratchpad/feet1.png. Ảnh nguồn Small27061.png: hai đế chân trong ảnh chỉ lệch nhau ~1,5 px (0,75 đơn vị), tức phần lớn 2 đơn vị này đến từ cách xếp mảnh chứ không phải nét vẽ.

**Cach sua:** Hạ chân thêm 2 đơn vị (dy = -31) rồi đo lại; nếu chân sau chìm thì phải sửa trong ảnh sprite chứ không sửa được bằng dy.

<sub>Lang kinh `loi-moi` · LOI MOI · phan bien: HIỆN TƯỢNG ĐÚNG, NGUYÊN NHÂN BÁO CÁO SAI. Tôi phóng to v3A.png vùng (680,605)-(750,635) kèm lưới hàng (scratchpad/r_feet_grid.png): đế chân sau xuống tới pano y≈627 (world 432,5 — chạm/lún nhẹ, đúng), đế chân trước dừng ở pano y≈623 (world 430,5) ⇒ hở ~2 đơn vị. Đến đây khớp. NHƯNG người báo cáo đo nhầm ảnh nguồn: khung chân dùng ở tư thế đứng là CharInfo[0][1][0] = 1 → part 58 pi[1] = Small27062.png, KHÔNG phải 27061. Tôi quét alpha từng cột của Small27062.png: một bàn chân chạm tới hàng 63 (cột 14-22), bàn kia chỉ tới hàng 58-60 (cột 44-60) ⇒ lệch 4-5 px ảnh = 2-2,5 đơn vị. Tức là TOÀN BỘ khoảng hở đến từ NÉT VẼ, không phải từ làm tròn dy. Hệ quả: cách sửa đề xuất (hạ dy xuống -31) là SAI — sẽ chôn chân sau xuống đất mà chân trước vẫn lệch. Phải sửa trong ảnh sprite. Giữ 'nhe'.</sub>

### 35. Điểm hồi sinh của máy chủ nằm lơ lửng 2 ô trên mặt đất

**O dau:** D:/NRO_GoiBanGiao/server/internal/game/lang.go:27 — ViTriYKhoiDau: 384

Mặt đất map 1 ở hàng 18, tức y = 432. Điểm hồi sinh vẫn để y = 384 = hàng 16, cao hơn mặt đất đúng 2 ô. Nhân vật mới tạo vào map là rơi tự do 48 đơn vị rồi mới đứng. Không gây kẹt (ô (16,4) không đặc nên máy chủ vẫn chấp nhận) nhưng là số cũ chưa theo bố cục mới.

**Bang chung:** lang.go dòng 27: {ID: 0, Ten: "Làng Lá", ..., ViTriXKhoiDau: 100, ViTriYKhoiDau: 384}. 1.bytes: hàng đầu tiên có tile nền ở cột 4 là hàng 18 (tile 62) → mặt đất world y = 18*24 = 432. Cả 5 làng đều để 384.

**Cach sua:** Đổi ViTriYKhoiDau của Làng Lá thành 432. Bốn làng còn lại chưa có map nên để nguyên cũng được, nhưng nên ghi chú rõ là giá trị tạm.

<sub>Lang kinh `loi-moi` · LOI MOI · phan bien: ĐÚNG, kiểm được cả hai đầu. lang.go:27-31: cả 5 làng đều ViTriXKhoiDau: 100, ViTriYKhoiDau: 384; nhanvat.go:87 dùng thẳng giá trị này. 1.bytes: cột 4 có tile nền đầu tiên ở hàng 18 (tile 62) ⇒ mặt đất world y = 432. 384 = hàng 16, cao hơn đúng 2 ô = 48 đơn vị. map-1-collision.bin hàng 16 không đặc (chỉ cột 0 và 51) nên máy chủ vẫn nhận vị trí, không kẹt — đúng như báo cáo tự nói. Triệu chứng chỉ là một cú rơi ngắn lúc vào map. Giữ 'nhe'.</sub>


---

## Tong ket

Đã tự kiểm lại bằng ảnh (đo pixel + xem thật) trước khi tổng kết.

---

# TỔNG KẾT VÒNG 3 — map 1 Làng Lá

## 1. Chín việc vừa sửa: cái nào ăn, cái nào không

| # | Việc sửa | Kết quả |
|---|---|---|
| 1 | Hai khuôn mặt → xoá mảnh đầu | **PHẢN TÁC DỤNG.** Hết hai mặt, nhưng ảnh thân vốn chỉ vẽ nửa dưới cái đầu (cắt phẳng ngay trên chân mày). Xoá mảnh đầu = xoá luôn sọ, tóc gai và **băng trán lá Konoha**. Nhìn ảnh là thấy ngay đỉnh đầu bị vạt phẳng. Gốc lỗi cũ là mảnh đầu đặt lệch, không phải thừa mảnh. |
| 2 | Nhân vật không có chân → xếp lại thân/chân | **ĂN một nửa.** Có chân rồi, nhưng chân đè lệch thân 11px nên nhân vật giờ có **bốn bàn tay và hai dải đai cam**. |
| 3 | Nhân vật lún → hạ 11 đơn vị | **ĂN gần đủ.** Chân sau chạm đất đúng; chân trước còn hở trời 1–3px (do nét vẽ trong ảnh sprite, không sửa bằng số được). |
| 4 | Tường vô hình hai cột biên | **KHÔNG GIỮ ĐƯỢC.** Ở pano4 nhân vật đứng lọt vào đúng cột 51 — cột tường. |
| 5 | Va chạm máy chủ biết hai cái gò | **Máy chủ ĐÚNG, nhưng client không biết.** 13 ô ở hai gò: client cho đi vào, máy chủ chặn → gói bị bỏ, người khác thấy đứng chết một chỗ. Đây là lỗi nặng mới. |
| 6 | Nhà/cây biến mất khi cuộn | **ĂN.** Xem cả 5 ảnh không thấy vật nào biến mất nữa. |
| 7 | Bóng chìm dưới đất → bù 12 | **ĂN, nhưng lộ lỗi khác.** Bóng hết chìm, giờ 4 hàng pixel trên cùng của bóng nổi lên trời. Và bóng vẫn **lệch trái 21px** so với hai bàn chân — chân sau không có bóng. |
| 8 | Ô chân dung HUD | **KHÔNG ĂN.** Ô vẫn là biểu tượng "ảnh hỏng". Máy chủ gửi bảng ánh xạ đúng, nhưng HUD góc trái **không hề đọc bảng đó** — chỉ bảng thông tin nhân vật (mở từ menu) mới dùng. |
| 9 | Thanh máu/chakra hiện 0 | **ĂN.** Nay hiện 150/150 đúng. |

Tóm lại: **4 ăn, 2 ăn một nửa, 3 hỏng hoặc phản tác dụng.**

## 2. Lỗi xếp theo mức ảnh hưởng người chơi

**Chặn đường / phá trải nghiệm**
1. **Đi sang phải tới cuối làng là mất nhân vật.** Camera dừng sớm hơn đúng 2 ô so với chỗ nhân vật được phép đứng — cả thân ra ngoài màn, chỉ còn mẩu bóng ở mép. Người mới sẽ hoảng.
2. **13 ô ở hai gò lệch nhau giữa client và máy chủ.** Client cho đi thẳng vào sườn gò, máy chủ từ chối → desync, đứng chết với người chơi khác. Đúng cái mà bản sửa số 5 định tránh.

**Nhìn thấy ngay, hỏng hình nhân vật**
3. Nhân vật **mất nóc đầu + băng trán** (dấu hiệu nhận diện ninja Làng Lá) — do bản sửa #1.
4. Nhân vật **bốn bàn tay, hai dải đai** — thân và chân đè lệch 11px.
5. **Bóng lệch trái 21px**, nửa trên bóng nổi lên trời, chân sau không có bóng. Gốc là cả hình nhân vật bị vẽ lệch phải 7,5 đơn vị so với vị trí thật.

**Nhìn thấy ngay, hỏng bức tranh**
6. **Hai đường kẻ ngang chạy hết màn hình**: y=344 (đáy núi Hokage cắt cứng) và y=561 (chặt cụt cả lớp làng xa — mái chùa, tán cây, núi nhỏ đều bị xén phẳng). Cái ở y=561 là thứ đập vào mắt đầu tiên khi mở ảnh.
7. **22% đáy màn hình là khối đất nâu trơn**, cộng dải trời trống ở trên → khoảng 40% màn hình không có gì, cả ngôi làng bị ép vào một dải hẹp giữa.
8. **Nét vẽ nhân vật chỏi hẳn với nền**: pixel art thô đứng cạnh tranh vẽ mịn, lại là tỉ lệ chibi cạnh kiến trúc tỉ lệ thật. Nhìn 1 giây là thấy hai game khác nhau.
9. **Ô chân dung HUD vẫn hỏng** — góc trên trái, chỗ quan trọng nhất giao diện.

**Vặt nhưng chướng mắt**
10. Sọc xám rộng 7px ở mép 4 cái gò (một cột pixel bị kéo giãn 7 lần) + hai tile 65/66 là ảnh hỏng, khác chất liệu hẳn đất xung quanh.
11. Nhà, chum, thùng, giếng, cột đèn đều **hở trời 3px dưới chân** — trông như đồ chơi bày trên kệ.
12. Đoạn lát đá duy nhất đọc ra như một **rãnh sỏi đào xuống**, không phải mặt đường. Cả map chỉ có đúng đoạn đó.
13. Núi Hokage thấp hơn cả mái nhà (nhà cao gấp ~3,8 lần núi) — lỗi tỉ lệ, không phải lỗi nhận diện.
14. Địa hình 50 ô chỉ có 2 bậc cao 1 ô, còn lại phẳng lì; không có bậc thang lên gò.
15. Điểm hồi sinh và NPC đặt ở y=384, cao hơn mặt đất đúng 2 ô (nhân vật rơi nhẹ lúc vào; NPC sẽ lơ lửng khi bật lên).
16. Chân dung người chơi (27078) đang dùng **chung ảnh với NPC thầy Iruka**, và bảng ánh xạ chỉ gửi 1 trong 9 kiểu tóc.
17. Số trên thanh máu/chakra cùng màu với thanh → gần như không đọc được; một dải phủ tối cắt cứng ở y=34.

## 3. Điểm hoàn thiện: **56/100** (vòng 2 là 52, tăng 4)

**Được cộng:** nhân vật đã có chân, đứng đúng mặt đất, thanh máu/chakra chạy, nhà cây không còn biến mất khi cuộn, dữ liệu va chạm phía máy chủ đã đúng. Đó là những thứ trước đây làm bản build không dùng để test được.

**Bị trừ gần hết phần cộng:** đẻ ra hai lỗi nặng mới (đi hết map là mất nhân vật; 13 ô desync ở gò), một lỗi lùi về mỹ thuật (nhân vật mất nóc đầu và băng trán — trước còn có đầu, giờ không), và ba bản sửa không ăn.

**Vì sao chưa tới 70:** ở mức 70 thì người lạ mở game ra phải thấy "một ngôi làng đẹp, còn vài chỗ chưa xong". Hiện tại họ sẽ thấy hai đường kẻ ngang chặt đôi bức tranh, một phần tư màn hình là đất trơn, và một nhân vật pixel bị vạt đầu dán lên tranh vẽ. Đó là ấn tượng đầu tiên, chưa phải chi tiết nhỏ.

## 4. Năm việc đáng làm nhất tiếp theo

**1. Giữ nhân vật luôn trong màn hình (SỬA BẰNG CODE)**
Nới giới hạn cuộn camera thêm 2 ô (phải cuộn được tới 1042px chứ không phải 946px), hoặc dời tường vô hình vào sớm 2 ô. Kiểm: đi sát mép phải, nhân vật phải còn nhìn thấy toàn thân. Đây là lỗi duy nhất khiến người chơi mất dấu chính mình.

**2. Lắp lại cái đầu và bỏ hai bàn tay thừa (SỬA BẰNG CODE, có thể cần sửa nhẹ ảnh)**
Trả mảnh đầu 27041/42/43 về, đặt lệch lên trên ~30px (đừng để 0 — để 0 chính là "hai khuôn mặt" của vòng 2). Đổi lệch của mảnh chân từ 40 lên ~51 để hết bốn tay hai đai. Nếu vẫn chồng thì phải cắt 13 hàng dưới cùng của ảnh thân — đó là phần cần động vào ảnh.

**3. Cho client và máy chủ đồng ý về hai cái gò (SỬA BẰNG CODE/DỮ LIỆU)**
Client kiểm tường ở **hàng 17** chứ không phải hàng 18 — nên phải cắm cờ tường vào tile hàng 17 ở cột 23, 27, 34, 41 (hoặc đổi bốn ô đó sang tile khác đã có cờ), rồi chạy lại `tools/map-lang-la/sinh-va-cham.py`. Thêm một bài kiểm tự động: không ô nào được vừa "client cho đứng" vừa "server chặn".

**4. Vá hai đường cắt ngang và lấp đáy màn hình (PHẢI VẼ ẢNH MỚI)**
Đây là việc đổi ấn tượng đầu tiên nhiều nhất. Cần: kéo dài lớp làng xa và lớp núi xuống tới đường cỏ (hết đường kẻ y=561 và y=344), vẽ lại 4 tile mép gò 65/66, và nhét nội dung vào 22% đất dưới — mạch đá, rễ cây, đường hầm — như tranh concept đã duyệt. Kèm một lớp cỏ tiền cảnh mỏng phủ lên chân nhà/chum/giếng để hết hở trời 3px và để mọi thứ "cắm xuống đất".

**5. Căn bóng cho đúng và cắt bóng theo mặt đất (SỬA BẰNG CODE)**
Hình nhân vật đang vẽ lệch phải 7,5 đơn vị so với vị trí thật — sửa chỗ đó thì bóng tự về đúng dưới chân. Đồng thời cắt bóng theo hàng ô mặt đất để nửa trên không tràn lên trời (đừng hạ số bù xuống, làm thế bóng chìm lại). Một lần sửa chữa được ba triệu chứng.

**Ghi chú kèm:** ô chân dung HUD nên để riêng thành một việc nhỏ — trước hết mở bảng thông tin nhân vật trong game xem chân dung có hiện không, để biết lỗi nằm ở bảng ánh xạ hay ở chỗ vẽ HUD (nhiều khả năng là chỗ vẽ HUD, vì HUD không hề đọc bảng đó). Và cần **vẽ một ảnh chân dung riêng cho người chơi** — hiện đang dùng chung ảnh với NPC thầy Iruka.