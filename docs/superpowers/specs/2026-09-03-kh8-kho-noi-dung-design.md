# KH8 — Kho nội dung YAML + nạp nóng + cân bằng thực nghiệm (thiết kế)

**Ngày:** 03/09/2026 · **Trạng thái:** thiết kế (chờ duyệt) · **Chỉ tài liệu, chưa có dòng code nào.**

**Phụ thuộc:** KH6 (chiến đấu), KH6b (boss), KH7 (vật phẩm) — cả ba đã có code trên `main`.
KH8 **không thêm tính năng chơi được nào**; nó chuyển chỗ ở của những con số đã có, rồi đo
lại chúng. Đó là lý do nó nguy hiểm theo một kiểu khác: mọi lỗi của KH8 đều là lỗi *đổi hành
vi mà không ai để ý*.

**Nguồn đối chiếu:** toàn bộ khẳng định về code hiện tại trong tài liệu này đều kèm
`file:dòng` đã đọc thật trong cây làm việc `D:\NRO_wt_kh8spec` (nhánh `kh8-spec`, gốc
`d8c568a`). KH8 **không** cần đối chiếu client Unity như KH6/KH7 — vì KH8 cố ý **không đụng
vào một khuôn gói nào**; xem §1.2. Chỗ nào đụng tới trần byte của client thì trích lại số đã
được KH7 xác minh, không giải mã lại.

---

## 0. Tóm tắt cho chủ dự án — 6 câu

1. **Hôm nay muốn đổi giá một món thuốc là phải sửa mã nguồn Go, chạy test, biên dịch lại
   máy chủ, tắt máy chủ, bật lại.** Đó là lý do KH8 tồn tại — không phải vì "YAML đẹp hơn".
2. KH8 đưa **số liệu nội dung** (chiêu, quái, boss, mẫu vật phẩm, cửa hàng, bảng rơi, hằng
   cân bằng) ra **file YAML người thường đọc được**, còn **khuôn gói tin và trần byte của
   client vẫn nằm trong Go** — vì đó là sự thật của client, không phải nội dung game.
3. **Nạp nóng** nghĩa là sửa file rồi bảo máy chủ đọc lại **mà không đá ai ra khỏi game**.
   Cách làm an toàn duy nhất: đọc file → kiểm tra toàn bộ → dựng một **bản mới hoàn chỉnh**
   → **tráo con trỏ**. Bản cũ vẫn sống cho tới khi người đang dùng nó dùng xong.
4. **Có những thứ không nạp nóng được** và tài liệu này nói thẳng ra: số ô túi (đã gửi xuống
   client lúc vào game), boss đang sống (đang giữa trận), lịch sinh boss (dựng một lần lúc
   khởi động). Chúng chỉ đổi ở **lần sau**.
5. **Căn cước là chỗ chết người:** id mẫu vật phẩm hiện nay = **vị trí trong danh sách**, và
   id đó **đang nằm trong cơ sở dữ liệu của người chơi**. Đổi thứ tự trong YAML = đổi món đồ
   của mọi người. YAML vì thế phải **khai `id` tường minh**, và bộ kiểm phải **từ chối** mọi
   thay đổi id.
6. **Task 0 là một cái cổng:** trước khi đổi bất cứ con số nào, phải chứng minh "YAML sinh ra
   **đúng từng byte** như dữ liệu Go hôm nay". Chứng minh xong mới được chỉnh số. Không có
   cổng này thì không ai phân biệt được "số mới" với "bộ nạp sai".

---

## 1. Phạm vi bản 1

### 1.1 Ra YAML

| # | Nội dung | Chỗ ở hôm nay | Vì sao ra YAML |
|---|---|---|---|
| N1 | **Chiêu thức** (20 chiêu × 3 cấp) | `chieuthuc.go:92-110` `taoCapChieu` sinh bằng công thức `dameGoc + (cap-1)*20`; bảng tên `chieuthuc.go:131-137` | Chính chú thích trong code đã nói: *"Gia tri chi la KHOI DAU can chinh o KH8"* (`chieuthuc.go:89-91`). Thêm một chiêu hiện nay phải sửa hàm sinh, tức sửa **tất cả** chiêu |
| N2 | **Quái** | `quai.go:180-190` `taoQuaiMacDinh` đọc `catalogLangLa.Mob` (một mẫu duy nhất, `SatThuongCoBan: 5` gán cứng tại chỗ) | Code tự khai là tạm: *"CHUA CO catalog quai theo lang/theo map ... Danh dau ro dat de coordinator biet can thay bang catalog that"* (`quai.go:172-179`) |
| N3 | **Boss** | `noidungboss.go:170-364` — hai con boss, dữ liệu thuần Go | Chú thích đầu file: *"thiet ke §8 chot chuyen kho noi dung sang YAML o KH8"* (`noidungboss.go:10-11`) |
| N4 | **Mẫu vật phẩm + mẫu option** | `mauvatpham.go:97-106` (8 option) và `:133-164` (5 mẫu) | KH7 §9 chốt thẳng: đơn giá/danh mục/tỉ lệ/chỉ số option **thuộc KH8** |
| N5 | **Cửa hàng + giá** | `cuahang.go:116-137` — "Tạp hoá Làng Lá", 5 món, giá 100/1000/5 | Giá là dữ liệu **riêng của từng cửa hàng** (`cuahang.go:70-73`), đúng thứ phải chỉnh thường xuyên. Giá **bán lại** hiện là `max(1, Gia/4)` — công thức tạm, xem §5.4 |
| N6 | **Bảng rơi** | `roido.go:86-116` — quái 15%/10%/2%, boss/boss lớn | Code tự khai: *"thuộc KH8. Chúng cố tình rộng rãi để nghiệm thu tay không phải..."* (`roido.go:80-84`) |
| N7 | **Hằng cân bằng làng** | `lang.go:25-32` — `HpKhoiDau` 100..250, `SatThuongKhoiDau` 10..16, `VangKhoiDau` 1000 | Chú thích ngay trên bảng: *"Ke hoach 6 se chuyen bang nay sang file YAML nap nong"* (`lang.go:22-23`) — lời hứa chưa trả |
| N8 | **Sức chứa hành trang** | `goihanhtrang.go:104-106` — `SoOMacChuan=5`, `SoOTuiChuan=20`, `SoORuongChuan=30` | Là số cân bằng thuần tuý. **Nhưng nạp nóng có điều kiện** — xem §3.4 |

### 1.2 CỐ Ý giữ trong Go — và vì sao

Ranh giới này là phần quan trọng nhất của §1. Nói sai chỗ này thì KH8 biến thành "cho phép
người vận hành sửa giao thức bằng file text" — tức là mở lại mọi lỗi im lặng mà KH6/KH7 đã
trả giá để bịt.

**Nguyên tắc chia:** *YAML chứa những gì **chúng ta** quyết định. Go chứa những gì **client**
quyết định.*

| Thứ | Ở đâu | Vì sao KHÔNG ra YAML |
|---|---|---|
| Mã lệnh gói (`-36`, `-44`, `-28/8`, `93`…) | Go, hằng | Là sự thật của client Unity đã biên dịch. Sửa được bằng file text nghĩa là một cú gõ nhầm làm client đọc lệch dòng byte — lỗi *im lặng tuyệt đối* (§6.18 BAT-DAU) |
| Thứ tự trường trong gói (`ghiQuai` hp→level→hpMax, `chieuthuc.go:71-85`) | Go | Đảo hai trường vẫn ra đúng tổng độ dài, client **không báo lỗi**, giá trị sai hoàn toàn (`quai.go:127-132`). Không được để dữ liệu quyết định |
| Trần byte: `tranMauOption=255`, `tranMauVatPham=32767`, `tranParamOption=65535`, `tranThanGoiItem=100000` (`mauvatpham.go:73-88`), trần 127 ô của `-36` | Go, hằng | Đây là **sức đọc của client**, không phải lựa chọn thiết kế. YAML **bị kiểm** theo chúng, không được **khai** chúng |
| Công thức `TinhSatThuong` (`satthuong.go:82-115`), vòng khắc chế ngũ hành (`satthuong.go:34-50`) | Go | Là *luật chơi*, không phải *số liệu*. YAML cấp tham số cho công thức; đổi công thức là đổi code + test |
| Nhịp actor 100ms, `tamNhatBinhPhuong = 120*120` (`dovuong.go:105-111`) | Go | Nhịp là ràng buộc kiến trúc (mọi `GiayNghi`/`NhipDanhMs` phải là bội số của nó). Tầm nhặt gắn với `kiemDiChuyenHopLe` — đổi bằng file là mở đường lách chống gian lận. **Đề xuất: tầm nhặt ở lại Go bản 1**, xem §8.5 |
| **Kiểu chiêu** (`loaiDanh=1`, `loaiBuff=2`, `loaiTuThan=3`, `chieuthuc.go:120-127`) | Go | Đây chính là cảnh báo của `LO-TRINH.md` (khung "Phải chốt trước KH8"): thêm chiêu **cùng kiểu** là thêm YAML, rẻ; thêm **kiểu mới** (choáng, đẩy lùi, độc, khiên, hút máu) là **viết code máy chủ mới**. YAML chỉ được chọn trong các kiểu Go đã cài; kiểu lạ → validator **đỏ**, không im lặng bỏ qua |
| Ảnh, part, tile, collision | `asset-lock.json` (đã ký) | Xem §6 |

**Phương án bị loại — "đưa hết ra YAML, kể cả khuôn gói".** Nghe hấp dẫn ("sau này thêm gói
mới không cần agent"). Bị loại vì: khuôn gói **không phải thứ chúng ta được quyền chọn** —
nó do client đã biên dịch quyết định, và mọi sai lệch đều im lặng. Đưa nó ra file chỉ chuyển
một lỗi biên dịch bắt được thành một lỗi chỉ lộ ra khi có người chơi thật vào nhìn.

**Phương án bị loại — "để nội dung trong CSDL, sửa qua panel web".** Đây là cách bản Java cũ
làm (logic + chỉ số trong MySQL, panel sửa game đang chạy). Bị loại cho bản 1 vì bốn lý do:
(a) nội dung mất khả năng **lần vết bằng git** — hôm nay không ai trả lời được "ai đổi giá
lúc nào" trong bản cũ; (b) nội dung không đi cùng mã nguồn nên **triển khai lệch phiên bản**
âm thầm; (c) review đối kháng không soi được một hàng SQL; (d) đã có một đường ghi CSDL
nguy hiểm (panel web ghi thẳng, D9 của KH7) — không mở thêm. **Đánh đổi chấp nhận:** chủ dự
án phải sửa file text thay vì bấm nút web. Panel web đọc/hiển thị nội dung là chuyện của W3,
và nếu sau này muốn sửa qua web thì web **ghi file YAML + commit**, không ghi CSDL.

### 1.3 Ngoài phạm vi KH8 hoàn toàn

Chế tác/cường hoá **luồng** (thuộc KH7 bản sau, `BO-SUNG` Khối 9 yêu cầu ghi rõ — ghi ở đây:
**luồng chế tác thuộc KH7 bản 2/3, số liệu chế tác thuộc KH8**); nhiệm vụ và thoại NPC (KH9 —
tuy chúng cũng là nội dung, nhưng khuôn dữ liệu nhiệm vụ chưa tồn tại nên đưa vào KH8 là
thiết kế mù); bản đồ và cổng (`banDo.go:78-90` — gắn với collision đã ký, xem §6).

---

## 2. Căn cước + bất biến

### 2.1 Vì sao mục này đứng trước mọi mục khác

KH7 §1.1 đã trả giá để rút ra một câu: **"chỉ số = danh tính" là gốc của mọi lỗi im lặng.**
KH8 làm đúng cái việc nguy hiểm nhất với căn cước: **chuyển danh sách từ code sang file**.
Một danh sách trong file thì ai cũng thấy "sắp xếp lại cho gọn" là vô hại. Nó không vô hại.

Bằng chứng thẳng trong code hôm nay (`mauvatpham.go:14-17`):

> *"id KHONG DI QUA MANG. Client dat id = CHI SO vong lap (Controller.cs:6479 va :6498).
> Chen mot mau vao GIUA danh sach se doi id cua moi mau sau no — va do la id da nam trong
> CSDL cua nguoi choi. CHI DUOC THEM VAO CUOI."*

Và `vat_pham.mau_id` là một cột thật trong CSDL người chơi (KH7 §5.1). Đổi thứ tự mẫu trong
YAML = **áo ninja của mọi người biến thành thuốc hồi máu**, không một dòng log.

### 2.2 Bảng căn cước — mỗi loại một quy tắc

| Loại | Căn cước | Hình thức trong YAML | Quy tắc |
|---|---|---|---|
| **Mẫu vật phẩm** | `id` = chỉ số liên tục từ 0 (`mauvatpham.go:133`, `TimMauVatPham` kiểm biên `:171-176`) | `id: 0` **khai tường minh** trên từng mục | Liên tục, bắt đầu 0, **chỉ thêm vào cuối**. Cấm đổi, cấm thủng lỗ, cấm hoán vị |
| **Mẫu option** | `id` = chỉ số (`mauvatpham.go:97-106`) | `id: 0` tường minh | Như trên. **Thêm:** id 6 và 7 **không được đổi ý nghĩa** — client băm cứng 22→6, 23→7 với param ×1000 (`mauvatpham.go:94-96`, `ItemOption.cs:29-38`) |
| **Chiêu** | hai số khác nhau: `templateId` (byte, 1..4 theo slot, `chieuthuc.go:162`) và `skillId` (int16, mỗi chiêu chiếm 3 số liên tiếp, `:165`) | **cả hai** khai tường minh | `templateId` **không bao giờ = 0** (`chieuthuc.go:159-161`: tránh nhánh nhồi 2 cấp giả của bản Java). `skillId` không trùng giữa các chiêu/làng |
| **Quái** | `MauId` byte + `Id` (id thể hiện trong khu vực) | `ma:` chuỗi (khoá tra cứu) **+** `mau_id:` số | Xem §2.3 |
| **Boss** | `Ma` chuỗi (`MaBossNhanGiaMatNa = "nhan-gia-mat-na"`, `noidungboss.go:158-161`) | `ma:` chuỗi | Đã đúng sẵn — chuỗi bền, không phụ thuộc thứ tự. **Giữ nguyên** |
| **Cửa hàng** | `ID` (`cuahang.go:145-147` tra theo ID) | `id:` số tường minh | Không đi xuống client; nhưng máy chủ nhớ `muaBanCuaHangDangMo` theo id → đổi id giữa phiên là bán nhầm bảng giá |
| **Làng** | `ID` 0..4 (`lang.go:25-31`) | `id:` tường minh | Client cũ **chỉ gửi được làng 0/1/2** (`lang.go:20-21`). Không được thêm làng thứ 6 bằng YAML |

### 2.3 Quái — chỗ phải quyết, không được lờ

Quái hôm nay có **hai** số: `MauId` (loại quái) và `Id` (con thứ mấy trong khu vực,
`quai.go:187` lấy từ `spawn.ID`). `Id` đi **xuống client** trong gói `-24` bằng **một byte**
(`quai.go:140`) và bảng rơi tra theo `MauId` (`roido.go:118`).

**Chốt:** YAML dùng **`ma:` chuỗi** làm khoá tra cứu giữa các file (bảng rơi trỏ tới quái
bằng `ma`, không bằng số), **và** khai `mau_id:` tường minh vì nó đi xuống client. Bộ kiểm
đối chiếu hai chiều: mọi `ma` trong bảng rơi phải tồn tại; mọi `mau_id` phải duy nhất.

**Phương án bị loại — dùng số làm khoá tra cứu chéo giữa các file.** Nó biến mỗi tham chiếu
chéo thành một cơ hội gõ nhầm không ai đọc ra ("mau_id: 3" không tự nói nó là con gì). Chuỗi
`ma: soi-rung` gõ nhầm thì **validator bắt được ngay** vì không tra thấy; số gõ nhầm thì
**tra thấy con khác** và chạy êm.

### 2.4 Bất biến — thứ phải luôn đúng

Bốn câu này là hợp đồng của KH8. Mỗi câu **phải có ít nhất một test và một đột biến**.

| # | Bất biến | Nói bằng lời thường |
|---|---|---|
| **K1** | Với mọi loại có id-theo-chỉ-số (mẫu vật phẩm, mẫu option): tập id trong YAML = `{0, 1, …, n-1}`, **không lỗ, không trùng**, và **thứ tự trong file = thứ tự id**. | Danh sách chỉ được dài thêm ở cuối. |
| **K2** | Mọi tham chiếu chéo giữa các file YAML (bảng rơi → mẫu vật phẩm, cửa hàng → mẫu, boss → chiêu, boss → map, chiêu → làng) **tra được**. | Không có con trỏ trỏ vào hư không. |
| **K3** | Mọi số trong YAML **nằm trong miền client đọc được** — kiểm **lúc nạp**, không phải lúc gửi gói. | Số vượt trần bị chặn ở cửa, không lọt vào tới mạng. |
| **K4** | Sau khi nạp nóng, **không goroutine nào đang đọc dữ liệu bị sửa dở**. Ảnh chụp là bất biến; đổi = tráo cả ảnh, không sửa tại chỗ. | Không ai đang đọc một quyển sách trong lúc có người tẩy xoá nó. |

### 2.5 Đột biến bắt buộc

Khuôn KH7 §6: mỗi bất biến một đột biến, đột biến nào **không** giết được test là **lỗ hổng
test**, phải viết lại test chứ không phải mừng (§6.11 BAT-DAU).

| # | Đột biến (cố tình phá) | Bất biến bị phá | Test nào phải đỏ — và **chỉ** test đó |
|---|---|---|---|
| **M1** | Hoán vị hai mẫu vật phẩm trong YAML (đổi chỗ id 2 và 3, giữ nguyên `id:`) | K1 | Test "id phải bằng vị trí" đỏ |
| **M2** | Đổi `id:` của một mẫu giữa danh sách (2 → 9) | K1 | Test "tập id liên tục từ 0" đỏ |
| **M3** | Xoá một mẫu đang được CSDL tham chiếu | K1 + K2 | Test "danh mục chỉ được dài thêm" đỏ (xem §2.6 — đây là đột biến khó nhất) |
| **M4** | `param` option = 70000 (> 65535) | K3 | Test miền giá trị đỏ **lúc nạp**; test dựng gói **vẫn xanh** — chứng minh chặn ở đúng cửa |
| **M5** | `param` option = -1 | K3 | Như M4. (Client đọc **unsigned** — số âm hiện thành số khổng lồ, KH7 §7.1) |
| **M6** | Thêm mẫu option thứ 256 | K3 | Test trần 255 (`mauvatpham.go:75`) đỏ |
| **M7** | `NhipDanhMs: 1450` (không phải bội số 100) | K3 | Test bội số nhịp actor đỏ — hiện đã có `TestNoiDungBoss_NhipDanhHopVoiNhipActor` (`noidungboss_test.go:446`), KH8 mở rộng sang YAML |
| **M8** | `GiayNghi` của một pha boss = 0 giữa hai pha | K3 | Test "khoảng lặng giữa pha > 0" đỏ (lý do: `noidungboss.go:267-270`) |
| **M9** | Bảng rơi trỏ tới `mau_id` không tồn tại | K2 | Test tham chiếu chéo đỏ |
| **M10** | Boss cầm chiêu `templateId: 3` (loại BUFF) | K2 | `TestNoiDungBoss_ChieuCuaBossPhaiLaChieuDANH` (`noidungboss_test.go:151`) đỏ — boss cầm chiêu buff sẽ "vung tay vào không khí" (`noidungboss.go:86-92`) |
| **M11** | Vị trí sinh boss dời vào ô đặc / lơ lửng | K2 | `TestNoiDungBoss_ViTriSinhDungTrenMatDat` (`:226`) đỏ |
| **M12** | Bộ nạp **sửa tại chỗ** ảnh chụp cũ thay vì dựng ảnh mới | K4 | Test đua (hai goroutine đọc + một lần nạp nóng) đỏ. ⚠️ Cần `-race` — xem §8.2 |
| **M13** | Trả `TimMauVatPham` về trỏ thẳng vào slice gốc (bỏ bản sao) | K4 | Test "người gọi sửa kết quả không làm bẩn kho" đỏ — cùng loại bẫy `TimBanDo` đã phải vá (`banDo.go:145-148`) và `saoNoiDungBoss` (`noidungboss.go:370-379`) |
| **M14** | Chiêu khai `loai: 9` (kiểu chưa cài) | K3 | Test "kiểu chiêu phải thuộc tập Go cài đặt" đỏ |
| **M15** | Bỏ một trong 5 làng khỏi YAML | K2 | `TestCoDuNamLang` (`lang_test.go:30`) đỏ |
| **M16** | Nạp nóng một file YAML **hỏng cú pháp** | K4 | Test "nạp hỏng thì giữ nguyên bản cũ, máy chủ không sập, có log" đỏ. **Fail-closed**, xem §3.5 |

### 2.6 Đột biến M3 — vì sao khó, và chốt gì

"Xoá mẫu đang được CSDL tham chiếu" là đột biến duy nhất mà **file YAML tự nó không đủ để
phán quyết** — sự thật nằm trong CSDL của người chơi. Ba phương án:

| | Truy vấn CSDL lúc nạp | **Danh sách chỉ-thêm ghim trong repo (CHỌN)** | Không kiểm |
|---|---|---|---|
| Bắt được | đúng thật 100% | mọi lần xoá/đổi id | không gì |
| Chi phí | mỗi lần nạp nóng đụng CSDL, nạp nóng phụ thuộc CSDL sống | 0 | 0 |
| Chạy được trong test đơn vị | không (cần DSN — bẫy §G-0: `go test` báo `ok` khi 27 test bị SKIP) | có | — |

**Chốt:** giữ một file **`kho-noi-dung/khoa-id.yaml`** — bảng "id nào đã từng được cấp phát,
cho tên gì" — **chỉ được thêm dòng**, và validator so danh mục hiện tại với nó. Xoá một mẫu
= xoá một dòng trong file khoá = **test đỏ**. Muốn ngừng bán một món thì đặt cờ
`ngung_dung: true`, **không xoá dòng**. Đây đúng logic KH7 §1.2 đã chốt cho `vat_pham.id`:
*"KHÔNG BAO GIỜ tái sử dụng"*.

**Đánh đổi đã chọn:** file khoá là *tự khai*, nên hai file cùng sai vẫn "khớp" — đúng bẫy
§6.12b của `BAT-DAU-TU-DAY.md`. Vì vậy trong review đối kháng, người soi phải đối chiếu file
khoá với **`git log` của chính nó** (chỉ có dòng thêm, không có dòng sửa/xoá), chứ không đối
chiếu file khoá với danh mục.

---

## 3. Nạp nóng an toàn với actor

### 3.1 Vấn đề, nói bằng lời thường

Trạng thái khu vực chạy trong một **actor**: một goroutine `vongLap` nhịp 100ms, và **chỉ nó**
được chạm trạng thái khu vực (§6.9 BAT-DAU). Nhưng actor **đọc kho nội dung ngay trong nhịp**:

- `khuvuc.go:965` — `ChieuCuaLang(byte(nguoi.nv.Lang))` khi người chơi chọn chiêu
- `khuvuc.go:2545` và `:2660` — `TraCapChieu(...)` khi boss ra đòn
- `khuvuc.go:1129` — `roiDoCuaQuai(mucTieu.MauId)` khi quái chết
- `khuvuc.go:1052`, `:1265`, `:1384` — `TimLang(nguoi.nv.Lang)`

Nạp nóng = có người ghi vào chính những bảng đó **trong khi actor đang đọc**. Trong Go, đọc và
ghi một `map`/`slice` cùng lúc không phải "giá trị cũ hay mới" — nó là **hành vi không xác
định**, có thể sập tiến trình. Và một tiến trình sập là **đá văng toàn bộ người chơi** (cảnh
báo của chính `protocol/frame.go`, ghi ở `BO-SUNG` Khối 2).

### 3.2 Ba phương án — cân nhắc thật

| | **Ảnh chụp bất biến + con trỏ atomic (CHỌN)** | Khoá đọc-ghi (`sync.RWMutex`) | Gửi sự kiện vào từng actor |
|---|---|---|---|
| Cách làm | Nạp file → kiểm toàn bộ → dựng struct mới hoàn chỉnh → `atomic.Pointer.Store`. Actor lấy con trỏ **một lần ở đầu nhịp** rồi dùng suốt nhịp | Mỗi lần đọc `RLock`, nạp nóng `Lock` | Bộ nạp gửi `suKienDoiNoiDung` vào hộp thư từng khu vực; actor tự thay bản của mình |
| Actor có bao giờ chờ không | **Không** — đọc con trỏ là một lệnh máy | **Có** — `RLock` chờ nếu người ghi đang giữ. Đúng thứ §6.10b cấm | Không |
| Nhất quán trong một nhịp | **Có** — cả nhịp dùng một ảnh | Không: hai lần `RLock` trong cùng nhịp có thể ra hai bản khác nhau | Có |
| Rủi ro | Bản cũ sống thêm tối đa một nhịp (100ms). Bộ nhớ giữ hai bản một lúc — vài chục KB | Deadlock nếu ai đó `RLock` lồng nhau; và người ghi có thể bị bỏ đói | **Hộp thư actor sức chứa 256, `Gui()` rơi im lặng khi đầy** (KH6c). Một khu vực bận có thể **âm thầm không nhận** bản mới |
| Kiểm bằng test | được: hai goroutine đọc + một lần tráo, `-race` | được | khó — phải dựng 5 khu vực |

**Chốt: ảnh chụp bất biến + `atomic.Pointer`.**

Lý do quyết định không phải hiệu năng mà là **kỷ luật actor**: hai phương án kia đều thêm một
đường để actor bị chặn hoặc bị bỏ sót, mà đó chính là hai lỗ hổng nặng nhất dự án đã vá
(§6.10b: actor chờ socket 30s; KH6c: hộp thư đầy rơi im lặng). Phương án ảnh chụp **không
thêm đường nào** — actor vẫn chỉ đọc.

**Điều kiện đi kèm, bắt buộc:**

1. **Ảnh chụp là bất biến sau khi công bố.** Sau `Store`, không ai được sửa nó. Cưỡng chế
   bằng: mọi cửa tra cứu **trả bản sao sâu**, đúng như `saoNoiDungBoss` (`noidungboss.go:370-396`)
   và `TimBanDo` (`banDo.go:145-148`) đang làm — cả hai đã phải vá đúng lỗi này. Đột biến M13
   khoá điều đó lại.
2. **Actor lấy con trỏ MỘT LẦN ở đầu nhịp**, giữ trong biến cục bộ. Không gọi `Load()` rải rác
   giữa nhịp — nếu không, nửa đầu nhịp dùng bảng cũ, nửa sau dùng bảng mới, và một trận đánh
   có thể tính sát thương bằng hai bộ số khác nhau.
3. **Một con trỏ cho cả kho**, không phải mỗi loại một con trỏ. Vì tham chiếu chéo (bảng rơi
   trỏ mẫu vật phẩm) chỉ đúng khi hai bảng cùng thế hệ. Tráo riêng lẻ = có khoảnh khắc bảng
   rơi mới trỏ vào danh mục cũ.
4. **Mỗi ảnh chụp mang một số thế hệ** (`TheHe int64`, tăng dần) và số đó vào log mỗi lần
   tráo. Không có nó thì không ai trả lời được "lúc 21:47 máy chủ đang chạy bộ số nào".

### 3.3 Trình tự nạp nóng

```
1. Đọc TẤT CẢ file YAML vào bộ nhớ            (chưa đụng gì tới bản đang chạy)
2. Phân tích cú pháp                           lỗi → DỪNG, giữ bản cũ, log
3. Kiểm K1 (căn cước) + K2 (tham chiếu chéo) + K3 (miền giá trị)
                                               lỗi → DỪNG, giữ bản cũ, log
4. Đối chiếu khoa-id.yaml (§2.6)               lỗi → DỪNG, giữ bản cũ, log
5. Dựng struct KhoNoiDung mới, gán TheHe = cũ + 1
6. atomic.Pointer.Store(bản mới)               ← khoảnh khắc DUY NHẤT có thay đổi
7. Log: thế hệ, giờ, ai kích, thống kê (bao nhiêu chiêu/quái/mẫu)
```

Bước 1–5 **hoàn toàn không ảnh hưởng** máy chủ đang chạy: dù có lỗi ở đâu, người chơi không
biết gì cả. Đây là "ghi sổ trước — giao sau" của KH7 §2.3 áp cho nội dung.

### 3.4 Cái KHÔNG nạp nóng được — nói thẳng

| Thứ | Vì sao | Cách xử |
|---|---|---|
| **Số ô túi/mặc/rương** (`SoOMacChuan/SoOTuiChuan/SoORuongChuan`) | Đã gửi xuống client **trong ME_INFO** lúc vào thế giới, và client cấp phát mảng theo con số đó (KH7 task 5). Đổi giữa phiên = client có 20 ô, máy chủ nghĩ 25 → gói `-36` dài hơn client đọc → **lệch, im lặng** | Nạp nóng **được phép đọc** giá trị mới nhưng **chỉ áp cho phiên đăng nhập SAU**. Ảnh chụp mỗi phiên giữ số ô của **thế hệ lúc phiên đó vào game**. Phải có test |
| **Boss đang sống** | Đổi `Dame`, `MocNoiGian`, `HpChon` của một con đang giữa trận = người chơi thấy máu boss nhảy hoặc nộ vô lý | `taoBoss` đã sao chép sâu dữ liệu pha lúc sinh (`lichboss.go:239-243` ghi rõ ba lớp sao chép). Vì vậy boss đang sống **tự nhiên** giữ bản cũ. **Chỉ cần không phá tính chất đó** — thêm test khoá lại |
| **Lịch sinh boss** | Dựng **một lần lúc khởi động**: `main.go:169` `MucLichTuNoiDungBoss(DanhSachNoiDungBoss())` rồi `MoBoLichBoss` + `BatDau()` | Bản 1: thêm/bớt **con boss** cần khởi động lại. Đổi **số liệu** của con đã có thì có tác dụng ở **lần sinh sau**. Ghi rõ trong sổ tay vận hành |
| **Quái đã sinh trong khu vực** | `taoQuaiMacDinh` chạy lúc `MoKhuVuc` (`sokhuvuc.go:24`) | Áp cho lần hồi sinh sau. Không "sửa máu quái đang đứng" |
| **Chiêu người chơi đang giữ** | Người chơi đã chọn chiêu; `-60` đang bay | Ảnh chụp đầu nhịp (§3.2 điều kiện 2) làm cho một đòn luôn tính trọn bằng một bộ số |

**Luật chung rút ra, ghi thành một câu để không phải nhớ bảng trên:**
*nạp nóng đổi được **luật của lần sau**, không đổi được **thứ đang diễn ra**.*

### 3.5 Fail-closed

Nạp hỏng → **giữ nguyên bản cũ**, ghi log rõ ràng (file nào, dòng nào, sai gì), **không** áp
một phần. Đây là ngược với bản Java cũ: `panel_cmd.txt` được đọc rồi **xoá file** (nêu ở
`CLAUDE.md` gốc) — lệnh hỏng biến mất không dấu vết. Bản Go **không xoá file nguồn**, và
không có trạng thái "nạp được một nửa".

---

## 4. Kênh kích nạp nóng

### 4.1 Bốn phương án

| | Đánh giá |
|---|---|
| **`SIGHUP`** | ❌ **Không dùng được.** Máy chủ đích là **VPS Windows** (`BAT-DAU-TU-DAY.md` §2 + `CLAUDE.md`); Windows không có `SIGHUP`. `main.go:467` hiện chỉ bắt `os.Interrupt` + `SIGTERM` |
| **Gói admin qua giao thức game** | ❌ Bản 1. Cần một mã lệnh mới, tức đụng client — mà client là thứ **đã biên dịch**, không có gói admin nào sẵn. Đó là mở một khuôn gói mới cho một việc vận hành |
| **File cờ kiểu `panel_cmd.txt`** | ⚠️ Chạy được nhưng kế thừa đúng nhược điểm bản cũ: poll 2 giây, xoá file, không có danh tính người kích, không có phản hồi "nạp thành công hay hỏng ở dòng nào" |
| **HTTP nội bộ, chỉ nghe `127.0.0.1` (CHỌN)** | ✅ `POST /admin/nap-lai` trên một cổng riêng, `net/http` thư viện chuẩn, **chỉ bind `127.0.0.1`** |

### 4.2 Chốt và lý do

**Chốt: một điểm cuối HTTP nội bộ, mặc định `127.0.0.1:<cổng cấu hình>`, tắt được bằng biến
môi trường.**

Lý do quyết định là **phản hồi**: nạp nóng mà không biết kết quả thì tệ hơn không có. Với
HTTP, người kích nhận ngay thân trả lời — "thế hệ 7, 20 chiêu, 5 mẫu vật phẩm" hoặc
"`cua-hang.yaml` dòng 14: mau_id 9 không tồn tại". Với file cờ, chủ dự án phải đi mò log.
Đây đúng bài học §6.32: *gửi cho ai đó một thứ mà không xử lý phần trả lời là tệ hơn không gửi*.

Thêm: HTTP nội bộ **dùng lại được** cho W3 (panel web chạy cùng máy gọi `localhost`), và cho
một điểm cuối `/admin/tinh-trang` mà `BO-SUNG` Khối 5 đang đòi ("không có bảng đèn báo
sống/chết server").

### 4.3 Rủi ro bảo mật — §7

Nạp nóng là **quyền đổi luật chơi**. Ba rủi ro và ba chốt:

1. **Lộ ra internet.** Chốt: bind **`127.0.0.1`**, không phải `0.0.0.0`. Phải có test khẳng
   định địa chỉ bind, vì gõ nhầm một chữ là mở cổng quản trị ra đời. Muốn kích từ xa thì đi
   qua SSH tunnel, **không** nới bind.
2. **Bất kỳ tiến trình nào trên máy cũng gọi được.** Trên VPS Windows chạy AutoAdminLogon,
   "cùng máy" không phải một biên giới mạnh. Chốt: thêm một **khoá dùng chung** trong header,
   đọc từ biến môi trường; thiếu/sai → 403 và **ghi log**. Đây là lớp mỏng, ghi thẳng vào §8
   là *chưa phải xác thực thật*.
3. **Ai kích, lúc nào, thế hệ nào.** Bản Java cũ không trả lời được câu này. Chốt: mỗi lần nạp
   ghi một dòng log có giờ, thế hệ mới, và **thống kê tóm tắt**. Nhật ký sâu hơn (ai đổi số gì)
   là việc của `git log` trên chính các file YAML — thêm một lý do nữa để nội dung ở file chứ
   không ở CSDL (§1.2).

**Không làm bản 1:** đăng nhập admin, phân quyền, nạp nóng từ xa. Chúng thuộc **Kế hoạch B**.

---

## 5. Cân bằng thực nghiệm

### 5.1 Chốt đơn vị `SatThuong` — món nợ to nhất của KH8

`satthuong.go:85-92` để lại nguyên văn:

> *"TODO(KH8 - can bang so lieu): day la DIEN GIAI TAM THOI, CHUA duoc chu du an chot. Du lieu
> Task 1 (chieuthuc.go, taoCapChieu: "dameGoc := 20 + slot*10") trong GIONG sat thuong TUYET
> DOI hon la phan tram"*

Hôm nay code tính `ketQua := goc * chieu.SatThuong / 100` (`satthuong.go:92`) — tức coi
`SatThuong` là **phần trăm**. Còn dữ liệu (20/30/40…) **trông như** sát thương tuyệt đối. Hai
cách đọc cùng một con số.

| | **Phần trăm (CHỌN — giữ nguyên công thức)** | Tuyệt đối |
|---|---|---|
| Ý nghĩa | `SatThuong: 60` = "chiêu này gây 60% sát thương gốc của người dùng" | `SatThuong: 60` = "chiêu này gây 60 điểm, bất kể ai dùng" |
| Người chơi mạnh lên thì | chiêu mạnh theo — chỉ số nhân vật **có ý nghĩa** | chiêu đứng yên; lên cấp không đổi gì cả trừ máu |
| Cân bằng | một trục (chỉ số nhân vật) × một trục (hệ số chiêu) | hai trục rời nhau, dễ mâu thuẫn |
| Rủi ro tràn số | nhân trước chia sau → đã xử bằng int64 xuyên suốt (`satthuong.go:74-81`) | thấp |
| Đổi code | **không đổi một dòng nào** | phải sửa `satthuong.go:92` + mọi test sát thương |
| Dữ liệu hiện tại | phải **đọc lại** cho đúng nghĩa (§5.3) | giữ nguyên |

**Chốt: `SatThuong` là PHẦN TRĂM của sát thương gốc, 100 = ×1,0.** YAML **đổi tên trường**
thành `he_so_phan_tram:` để không ai đọc nhầm lần nữa — tên trường là tài liệu rẻ nhất.
Trường Go giữ tên `SatThuong` (khuôn gói `-28/7` đã ghim, `chieuthuc.go:80`); ánh xạ nằm ở bộ
nạp.

**Vì sao không chọn tuyệt đối:** nó làm chỉ số nhân vật thành trang trí. `lang.go` đã thiết kế
5 làng khác nhau bằng `SatThuongKhoiDau` 10..16 — với sát thương tuyệt đối, khác biệt đó biến
mất khỏi mọi đòn đánh bằng chiêu, và cả năm làng đánh giống hệt nhau.

### 5.2 Bảng công thức hiện có — nguồn sự thật để đo

| Thứ | Nơi | Nội dung |
|---|---|---|
| Sát thương | `satthuong.go:82-115` | `goc × heSoChieu/100 × heSoNguyenTo/100 × (2 nếu chí mạng)`, kẹp `[0, MaxInt32]` |
| Ngũ hành | `satthuong.go:34-50` | Hoả→Phong→Lôi→Thổ→Thuỷ→Hoả. Khắc ×1,25 · trung tính ×1,00 · bị khắc ×0,80 |
| Nộ boss | `khuvuc.go` + `noidungboss.go:224-231` | Mốc theo %HP, nhân **dồn** `HeSoDame`/`HeSoDef`/`HeSoNhip` |
| Sát thương quái | `quai.go:111-120` | **Chưa có công thức** — một hằng `SatThuongCoBan`, hiện gán cứng `5` (`quai.go:187`) |
| Nhịp AI quái | `quai.go:200-213` | Trả thù: tầm 300, nhịp 1000ms · Tự gây sự: tầm 100, nhịp 2000ms |

### 5.3 Cách ĐO — kịch bản giả lập, không phải "chơi thử rồi cảm thấy"

Nền đã có sẵn: KH6c đã dựng test tích hợp trong tiến trình (`khuvuc_stress_test.go`,
`suKienDiChuyen` + `phienStress`) chạy được **không cần TCP thật**. KH8 dùng lại nền đó, đổi
câu hỏi từ "actor có treo không" sang "trận đánh dài bao lâu".

**Bốn phép đo bản 1** (mỗi phép là một test có tên, chạy được bằng một lệnh, **in ra số**):

| Mã | Câu hỏi | Cách dựng | Con số muốn thấy (đề xuất, chủ dự án chốt) |
|---|---|---|---|
| **Đ1** | Người mới cấp 1 giết một con Sói Rừng mất bao lâu? | 1 người `HpKhoiDau` của làng, chiêu cận chiến cấp 1, quái `HP=100` (`catalog.json`), đồng hồ **tiêm** (không `time.Sleep`) | **8–15 giây** |
| **Đ2** | Năm người hạ boss lớn pha 1 mất bao lâu? | 5 người chỉ số khởi đầu, `HpChon` giữa (3000) | **2–3 phút** — con số `noidungboss.go:143-146` đã dự tính |
| **Đ3** | Làng nào thắng? | Chạy Đ1 **cho cả 5 làng**, so thời gian; và ma trận PvP 5×5 theo hệ | Chênh lệch giữa làng nhanh nhất và chậm nhất **≤ 25%** |
| **Đ4** | Đi bao lâu mới mua nổi một món trong cửa hàng? | Đếm vàng rơi ra mỗi phút từ bảng rơi (§N6) so với giá (§N5) | Áo 1000 vàng ≈ **20–40 phút** cho người mới |

**Hai điều kiện làm phép đo có nghĩa** (thiếu cái nào thì con số là rác):

1. **Đồng hồ tiêm, không ngủ thật.** Dự án đã vấp: `TestTichHop_DanhQuaiConSong` chớp tắt dưới
   tải (§8 BAT-DAU) vì nhạy thời gian. Phép đo cân bằng mà phụ thuộc tốc độ máy thì mỗi lần
   chạy ra một kết quả.
2. **Ngẫu nhiên tiêm.** `roido.go` đã làm sẵn đúng cách (`quyetDinhRoiDo(bang, chon)` nhận hàm
   chọn, test tiêm `func(int) int { return 0 }` — `roido.go:149-154`). Sát thương/chí mạng phải
   theo cùng khuôn đó.

**Đầu ra để chủ dự án tự chỉnh số:** mỗi lần chạy sinh một **bảng tóm tắt** (một file text,
một bảng, không cần đọc code) dạng:

```
BANG CAN BANG — the he 7 — 03/09/2026
lang        Đ1 giet Soi Rung   Đ3 thang/thua PvP (5 tran)
Làng Lá             11,2 giây              3 / 2
Làng Cát            14,8 giây              2 / 3
...
Đ2 boss lon pha 1, 5 nguoi:  2 phút 41 giây
Đ4 mua áo 1000 vàng:         31 phút
```

Chủ dự án đọc bảng → sửa một số trong YAML → chạy lại một lệnh → đọc bảng mới. **Không cần
agent, không cần biên dịch.** Đây mới là "cân bằng thực nghiệm"; mọi thứ khác trong KH8 chỉ
là hạ tầng cho vòng lặp này.

### 5.4 Bộ số bản 1 — KHÔNG được đổi hành vi test hiện tại

**Nguyên tắc: task chuyển sang YAML giữ số **y hệt**; đổi số là task riêng, sau cổng §7.**

Đã rà bộ test hiện có để biết cái gì ghim cái gì:

| Test | Ghim gì | Đổi số có làm đỏ không |
|---|---|---|
| `mauvatpham_test.go:47` `TestGoiCatalogItemSubtype1DungByte` | Dựng `want` **từ chính `DanhSachMauVatPham()`** | ❌ Không — đây là **điểm mù oracle** (§6.11 BAT-DAU: *"test đếm số làng bằng chính nguồn dữ liệu nó phải kiểm"*). Nó bắt lỗi **thứ tự trường**, không bắt lỗi **giá trị** |
| `cuahang_test.go:348` `TestDanhSachCuaHangHopLe` | Chỉ kiểm "dựng được gói", không kiểm giá | ❌ Không |
| `noidungboss_test.go` (18 test) | **Tính chất**: mốc nộ giảm dần không trùng (`:398`), chỉ mạnh lên (`:421`), nhịp là bội số 100ms (`:446`), vị trí sinh trên mặt đất (`:226`), chiêu tồn tại thật (`:114`) | ❌ Không với số hợp lệ · ✅ **Có** với số vi phạm tính chất — **đúng cái ta muốn** |
| `muaban_tichhop_test.go:117-128` | Dựng **cửa hàng riêng của test** (`{MauID: 0, Gia: 100}`), có chú thích *"test không vỡ khi ai đó đổi giá cửa hàng mẫu"* | ❌ Không — đã phòng sẵn |
| `meinfo_hanhtrang_test.go:168-199`, `meinfo_tichhop_test.go:161-176`, `thegioi_test.go:119-121`, `dongbo_tichhop_test.go:501-503` | So ME_INFO với **hằng** `SoOTuiChuan`… (không phải số trần) | ❌ Không nếu đổi hằng · ⚠️ **NHƯNG** xem §3.4: đổi số ô là chuyện khác hẳn |
| `dungvatpham_test.go` (nhiều chỗ) | Dùng `SoOTuiChuan` làm **biên** để test kiểm biên | ❌ Không |
| `lang_test.go:30` `TestCoDuNamLang` | **Đúng 5 làng** | ✅ Có nếu thêm/bớt làng — đúng ý (M15) |
| `chieuthuc_test.go:138/164/181` | Mỗi làng **4 chiêu**, **1 chiêu hồi máu**, mỗi chiêu **3 cấp** | ✅ **Có** nếu YAML khai 8 chiêu/làng |

**Kết luận quan trọng:** bộ test hiện tại **hầu như không ghim giá trị**, nó ghim **cấu trúc**.
Đó vừa là tin tốt (đổi số không làm đỏ oan) vừa là tin xấu (**đổi nhầm số cũng không đỏ**).
Vì vậy cổng task 0 (§7) là bắt buộc — nó là thứ duy nhất phát hiện được "bộ nạp làm lệch số".

**Ba chỗ phải quyết trước khi viết YAML:**

1. **Số chiêu mỗi làng.** `chieuthuc_test.go:138` ghim đúng 4. `LO-TRINH.md` đề xuất **8
   chiêu/làng, ~6 kiểu** và ghi rõ *"Chủ dự án chốt trước khi thiết kế KH8"*. **Chưa chốt** —
   §8.1. Bản 1 của KH8 nên giữ **4 chiêu** (chuyển chỗ ở, không đổi nội dung), việc thêm chiêu
   là một task sau, có sửa test.
2. **Giá bán lại cho NPC.** Hiện `max(1, MonHang.Gia/4)` (KH7 task 10, ghi trong §5 BAT-DAU:
   *"KH8 thêm giá gốc vào mẫu"*). Chốt: YAML thêm `gia_ban_lai:` **tường minh** cho từng món;
   `Gia/4` chỉ còn là **giá trị mặc định khi thiếu**. Vì `/4` là một luật ngầm không ai đọc
   được từ dữ liệu.
3. **Sát thương quái.** `quai.go:187` gán cứng `5` cho **mọi** con quái. YAML khai
   `sat_thuong_co_ban:` cho từng mẫu quái. Công thức đầy đủ của đặc tả §5.3 (dame riêng ±1%
   hoặc `hpMax×pDame%`) **hoãn** — bản 1 chỉ đưa con số ra ngoài, chưa đổi luật.

---

## 6. Quan hệ với `catalog.json` / `asset-lock.json`

### 6.1 Ranh giới

| | **Mỹ thuật đã ký** | **Nội dung (KH8)** |
|---|---|---|
| Là gì | icon, part (mảnh nhân vật), tile, ảnh nền, lưới va chạm | số liệu: máu, sát thương, giá, tỉ lệ, tên, thoại |
| Ở đâu | `assets/asset-lock.json` (444 KB) + `assets/catalog.json`, nhúng bằng `go:embed` (`dongbo.go:18-30`) | `kho-noi-dung/*.yaml` |
| Đổi thế nào | qua **quy trình mỹ thuật + ký lại lock** | sửa file, nạp nóng |
| Vì sao khác nhau | Ảnh là **artifact nhị phân**, sai một byte là client vẽ sai hoặc `IndexOutOfRange` nuốt im lặng. Nó cần SHA ghim | Số là thứ **cố ý muốn đổi hằng ngày** |

`catalog.json` bị ghim SHA thật: `catalog.go:100` `lockHashKhop(data, l, "server/catalog.json")`
— sai hash thì `NapCatalog` trả lỗi và `phaiNapCatalog` **panic lúc khởi động**
(`catalog.go:213-219`). Đây chính là bẫy §7.6 của KH7, và KH7 task 2 đã né bằng cách để mẫu
vật phẩm trong file Go riêng (`mauvatpham.go:4-7`).

### 6.2 Chốt

> **YAML tham chiếu id ảnh trong catalog/lock. YAML KHÔNG chứa ảnh, KHÔNG chứa hash, và
> KHÔNG được đổi `catalog.json`.**

Validator kiểm **hai chiều**:

1. Mọi `icon_id`, `phan`, `outfit` trong YAML phải **tồn tại trong closure của lock đang chạy**
   — `closure.part_ids` / `closure.sprite_ids` / `closure.mob_ids` (`catalog.go:76-83`).
   `NapCatalog` đã làm đúng việc này cho Iruka và mob (`catalog.go:107-109`); KH8 mở rộng cho
   nội dung YAML.
2. `catalog.json` giữ nguyên SHA — nếu một task KH8 đụng vào nó, test khớp lock **đỏ**.

**Vì sao kiểm id ảnh là quan trọng chứ không phải hình thức:** `noidungboss.go:54-57` ghi
thẳng — *"Dien mot id ngoai tap do vao Outfit la con boss VO HINH — client khong tim thay anh
thi khong ve gi, va may chu khong co mot dong log nao."* Với YAML, gõ nhầm một số dễ hơn hẳn
so với sửa hằng Go. Validator là thứ duy nhất chặn.

**Một no khách phải nói ra:** hôm nay **chưa có icon nào cho vật phẩm** — mọi mẫu để
`IconID: 0` và bảng `SmallImage` phía client đang **rỗng** (`mauvatpham.go:125-129`,
`dongbo.go:79-82`). Nghĩa là ô icon **trống**, không phải sai. KH8 **không sửa được** việc
này — nó cần bộ icon đi qua quy trình mỹ thuật + lock. YAML chỉ cần **sẵn sàng** nhận
`icon_id` khi bộ icon có.

**Phương án bị loại — nhét nội dung vào `catalog.json`.** Sẽ phải ký lại một artifact đã ký
**mỗi lần đổi một con số giá**. Quy trình ký tồn tại để bảo vệ thứ *không được đổi tuỳ tiện*;
đưa vào đó thứ *phải đổi tuỳ tiện* làm hỏng ý nghĩa của cả chữ ký.

---

## 7. Bảng task

Cột **"chạm `khuvuc.go`?"** để biết task nào phải chạy tuần tự: `khuvuc.go` hiện **134 KB**
(đo thật trong cây làm việc), hai agent cùng vào là chắc chắn đụng.

| # | Task | File mục tiêu | Phụ thuộc | Chạm `khuvuc.go`? |
|---|---|---|---|---|
| **0** | ⛔ **GATE — tương đương byte.** Bộ nạp + bộ kiểm tối thiểu, rồi test: **chuyển dữ liệu Go hiện tại sang YAML, nạp lại, dựng lại TẤT CẢ gói catalog (`-28/7` chiêu, `-28/8` mẫu vật phẩm + option, `-44` cửa hàng) và so BYTE-CHÍNH-XÁC với gói dựng từ dữ liệu Go.** Phải **y hệt**. Kèm quyết định thư viện YAML (§8.3) | `server/internal/noidung/` (gói mới) + test | không | **không** |
| **1** | Khuôn YAML + validator K1/K2/K3 + `khoa-id.yaml` (§2.6). Đột biến M1–M6, M9, M14 | `noidung/kiemtra.go` | 0 | **không** — song song |
| **2** | Ảnh chụp bất biến + `atomic.Pointer` + số thế hệ + cửa tra cứu trả **bản sao sâu**. Đột biến M12, M13 (cần `-race`, §8.2) | `noidung/kho.go` | 0 | **không** — song song |
| **3** | Chuyển **chiêu** sang YAML (giữ y số). Cửa `ChieuCuaLang`/`TraCapChieu` đọc từ ảnh chụp | `game/chieuthuc.go` + `kho-noi-dung/chieu.yaml` | 0, 1, 2 | **không** |
| **4** | Chuyển **mẫu vật phẩm + mẫu option**. ⚠️ Đây là task **rủi ro căn cước cao nhất** — id đang nằm trong CSDL người chơi (§2.1) | `game/mauvatpham.go` + `kho-noi-dung/vat-pham.yaml` | 0, 1, 2 | **không** |
| **5** | Chuyển **cửa hàng + giá**, thêm `gia_ban_lai:` tường minh (§5.4 điểm 2) | `game/cuahang.go`, `game/muaban.go` + `kho-noi-dung/cua-hang.yaml` | 4 | **không** |
| **6** | Chuyển **hằng làng** (`lang.go:25-31`) — trả lời lời hứa ghi ở `lang.go:22-23` | `game/lang.go` + `kho-noi-dung/lang.yaml` | 0, 1, 2 | **không** |
| **7** | Chuyển **boss**. Giữ nguyên tính chất sao chép sâu ba lớp (`lichboss.go:239-243`); thêm test "boss đang sống không đổi số khi nạp nóng". Đột biến M7, M8, M10, M11 | `game/noidungboss.go` + `kho-noi-dung/boss.yaml` | 3, 6 | **không** (file khác, nhưng **tuần tự** với task 3) |
| **8** | Chuyển **quái** + `sat_thuong_co_ban` cho từng mẫu (thay hằng `5` ở `quai.go:187`) | `game/quai.go`, `game/khuvuc.go` + `kho-noi-dung/quai.yaml` | 3 | ✅ **CÓ — tuần tự** |
| **9** | Chuyển **bảng rơi**, tham chiếu quái bằng `ma:` (§2.3) | `game/roido.go`, `game/khuvuc.go` + `kho-noi-dung/bang-roi.yaml` | 4, 8 | ✅ **CÓ — tuần tự, sau task 8** |
| **10** | Kênh nạp nóng: HTTP `127.0.0.1` + khoá dùng chung + log thế hệ. Test **khẳng định địa chỉ bind** (§4.3) | `cmd/nroserver/main.go`, `noidung/kich.go` | 2 | **không** |
| **11** | **Bộ đo cân bằng** Đ1–Đ4 (§5.3): đồng hồ tiêm, ngẫu nhiên tiêm, xuất **bảng tóm tắt** đọc được | `game/canbang_test.go` + `tools/` | 3, 6, 7, 8, 9 | **không** |
| **12** | **Chốt đơn vị `SatThuong`** (§5.1): đổi tên trường YAML thành `he_so_phan_tram`, xoá TODO ở `satthuong.go:85-92`, cập nhật chú thích suy luận số của `noidungboss.go:126-131` | `game/satthuong.go`, `noidung/` | 11 | **không** |
| **13** | **Vòng cân bằng thật**: chạy Đ1–Đ4, đọc bảng, chỉnh số YAML tới khi đạt ngưỡng §5.3. **Task duy nhất được đổi giá trị** | `kho-noi-dung/*.yaml` | 12 | **không** |
| **14** | Nghiệm thu tay bằng client Unity thật: sửa một giá trong YAML → nạp nóng → vào game xem đổi chưa, **không ai bị đá** | — | 13 | — (chỉ chủ dự án làm được) |

### 7.1 Vì sao task 0 là cái cổng đó

Ba cổng khác đã được cân nhắc và loại:

- *"Cổng = chốt số chiêu/kiểu chiêu"* — đúng là phải chốt (§8.1), nhưng đó là **quyết định của
  chủ dự án**, không phải thứ agent chứng minh được. Nó chặn **task 3**, không chặn cả KH8.
- *"Cổng = giải mã lại gói catalog"* — không cần: KH7 đã giải mã và có test byte-chính-xác. Làm
  lại là tốn công và tạo **bản sao sự thật thứ hai**.
- *"Cổng = validator"* — validator không trả lời được câu hỏi quan trọng nhất: *"bộ nạp có làm
  lệch số không?"* Một validator hoàn hảo vẫn cho qua một bộ nạp đọc nhầm `dx` thành `dy`.

**Tương đương byte** là cổng đúng vì nó là điều kiện duy nhất khiến câu này thành thật: *"KH8
không đổi hành vi game, chỉ đổi chỗ ở của số liệu"*. Chứng minh xong, mọi thay đổi hành vi về
sau đều **cố ý** và truy được về task 13.

### 7.2 Kỷ luật nhiều agent — bài học 03/09

- **Worktree riêng cho từng agent, ngay từ đầu.** 03/09 hai agent đổi nhánh trong **cùng một
  cây làm việc** → một agent commit vào nhánh của agent kia (§8 BAT-DAU).
- **Không khai trùng hằng.** Cùng ngày, task 3 và task 8 của KH7 mỗi bên tự khai
  `tranParamOption = 65535` → gộp vào `main` là **đứt build** (`b4f83e8`, vá `ba3a56d`). KH8
  còn dễ vấp hơn vì nhiều task cùng cần các trần byte. **Chốt: mọi trần byte khai một nơi duy
  nhất — `mauvatpham.go` hôm nay, và task 0 chuyển nguyên về gói `noidung/`, không sao chép.**
  Brief cho mỗi agent song song **phải liệt kê tên hằng/kiểu đã có**.
- **Gộp xong phải build + test lại `main`** trước khi coi là xong.
- **Kiểm `git status` trong worktree** trước khi tin báo cáo của agent — model rẻ hay quên
  commit (§8 BAT-DAU).
- **Bộ chạy đột biến phải có timeout từng con** — 02/09 kẹt 19 giờ vì thiếu (§8 BAT-DAU).

---

## 8. Chưa xác minh được — nói thẳng

1. **Số chiêu và số KIỂU chiêu chưa được chủ dự án chốt.** `LO-TRINH.md` (khung "⚠️ Phải chốt
   trước KH8", nêu 27/07/2026) đề xuất **8 chiêu/làng phủ ~6 kiểu**, hiện có **4 chiêu/làng, 3
   kiểu** (`chieuthuc.go:120-127`). Đây là câu hỏi *"phải viết thêm bao nhiêu code"*, không
   phải *"đổ thêm bao nhiêu dữ liệu"*. Chặn task 3 mở rộng; **không** chặn task 0–2.
2. **Chưa chạy được `go test -race`** — máy thiếu trình biên dịch C (§7 BAT-DAU). Mà đột biến
   **M12** (nạp nóng sửa tại chỗ) là loại lỗi mà **chỉ** `-race` bắt được đáng tin. Bản 1 phải
   viết sẵn test và ghi rõ **chưa chạy được**, đúng cách hai test đồng thời đang chờ. **Đây là
   lỗ hổng nghiệm thu lớn nhất của KH8.**
3. **Chưa quyết thư viện YAML.** `go.mod` hiện chỉ có **một** phụ thuộc trực tiếp:
   `github.com/jackc/pgx/v5`. Thư viện chuẩn Go **không có** bộ phân tích YAML, nên KH8 buộc
   phải chọn: (a) thêm `gopkg.in/yaml.v3` — một phụ thuộc mới cho một dự án đang giữ kỷ luật
   "chỉ thư viện chuẩn" (`chanip` được khen vì điều đó, §7 BAT-DAU); (b) dùng **JSON** (stdlib,
   0 phụ thuộc) — nhưng JSON không có chú thích, mà chú thích *"vì sao con số này"* chính là
   thứ quý nhất trong `noidungboss.go`; (c) tự viết bộ đọc YAML tối giản — **loại**, tự viết
   trình phân tích là nguồn lỗi mới. **Nghiêng về (a)** vì mục tiêu KH8 là *chủ dự án sửa được*
   và YAML + chú thích thắng JSON rõ ràng ở đó. **Chưa kiểm** `go mod` tải được thư viện trên
   máy này (có thể không có mạng lúc build) — task 0 phải thử **trước** mọi thứ khác.
4. **Chưa đo chi phí nạp nóng.** Dựng lại toàn bộ ảnh chụp mỗi lần nạp; hôm nay dữ liệu bé
   (20 chiêu, 5 mẫu, 2 boss) nên chắc chắn không sao. Nhưng chưa biết ngưỡng — nếu sau này có
   750 mẫu vật phẩm thì phải đo lại. Không phải vấn đề của bản 1.
5. **Chưa quyết tầm nhặt 120px có ra YAML hay không.** `dovuong.go:105-111` ghi rõ nó gắn với
   `kiemDiChuyenHopLe` (chống gian lận). Đề xuất **giữ trong Go bản 1** — nhưng nếu chủ dự án
   coi đó là số cảm giác chơi cần chỉnh thường xuyên thì phải bàn lại, vì nới nó bằng file là
   nới chống gian lận bằng file.
6. **Chưa biết panel web sẽ đọc nội dung thế nào.** W3 có thể muốn hiển thị bảng giá. Nếu web
   đọc thẳng YAML thì có hai bộ phân tích cho một định dạng (Go và TypeScript) — hai nguồn sự
   thật. Chưa thiết kế. Có thể máy chủ game phơi `GET /admin/noi-dung` trả JSON, nhưng chưa đo
   nhu cầu.
7. **Chưa xác minh giả định "client cũ chỉ gửi được làng 0/1/2"** (`lang.go:20-21`) bằng cách
   đọc client. Nếu sai thì ràng buộc "không thêm làng bằng YAML" (§2.2) chặt hơn cần thiết —
   nhưng chặt hơn thì an toàn, nên không chặn KH8.
8. **Chưa biết `Loai 6` của mẫu vật phẩm có thật là "tiêu hao" không** — KH7 task 7 ghi đó vẫn
   là **suy luận** (`case 6: eatPea` bên Java), chờ nghiệm thu tay task 13 của KH7. KH8 chuyển
   nguyên con số 6 sang YAML mà **không** hợp thức hoá nó; nếu KH7 task 13 lật lại thì sửa dữ
   liệu, không sửa code.
9. **Chưa nghiệm thu tay bất cứ gì.** Toàn bộ §5.3 là phép đo trong tiến trình. Nó trả lời
   "trận đánh dài bao lâu", **không** trả lời "chơi có vui không". Task 14 là bắt buộc.

---

## 9. Bẫy — đọc trước khi viết dòng code đầu tiên

1. **Đổi thứ tự trong file YAML trông vô hại hơn đổi thứ tự trong code Go rất nhiều.** Đây là
   rủi ro **duy nhất** mà KH8 *tạo ra* chứ không phải thừa hưởng. Mọi file YAML có id-theo-chỉ-số
   phải mở đầu bằng một khối chú thích in hoa nói điều đó (§2.1).
2. **`git diff` của file YAML tiếng Việt.** Đừng sửa bằng PowerShell `Get-Content -Raw` +
   `Set-Content` — mojibake sạch cả file (§6.15 BAT-DAU, đã vấp hai lần). Dùng Edit/Write hoặc
   Python `io.open(..., encoding="utf-8")`.
3. **Chú thích *vì sao* phải theo số sang YAML.** `noidungboss.go:120-149` có nguyên một đoạn
   suy luận *"Dame boss 35 → đòn 10..21: người chơi 150 máu chịu ~10 đòn"*. Bỏ lại phía sau là
   mất phần đắt nhất. YAML giữ được chú thích — đó là **lý do chọn YAML thay vì JSON** (§8.3).
4. **Số giây/mili-giây phải là bội số của nhịp actor 100ms.** `noidungboss.go:207-210` nói
   thẳng: *"1450ms tren thuc te van la 1500ms — phan le la con so noi doi"*. Validator chặn
   (M7).
5. **Bộ test hiện tại KHÔNG ghim giá trị, chỉ ghim cấu trúc** (§5.4). Đừng suy ra "test xanh =
   số đúng". Chỉ cổng task 0 mới nói được điều đó.
6. **Đừng để bộ nạp làm việc lúc tra cứu.** Mọi tính toán (dựng cấp chiêu từ công thức, tính
   `Gia/4`) phải làm **một lần lúc nạp**, kết quả nằm sẵn trong ảnh chụp. Tính lúc tra cứu =
   tính **trong nhịp actor**.
7. **`atomic.Pointer` không làm dữ liệu bên trong bất biến.** Nó chỉ đảm bảo con trỏ đọc/ghi
   nguyên tử. Nếu một cửa tra cứu trả `[]MucRoiDo` trỏ thẳng vào ảnh chụp và người gọi `append`
   vào đó — hỏng. Đúng bẫy `TimBanDo` (`banDo.go:145-148`) và `saoNoiDungBoss`
   (`noidungboss.go:370-379`) đã phải vá. Đột biến M13.
8. **Nạp nóng thành công cũng phải nói cho người kích biết nó làm gì** — thế hệ mới và thống kê.
   "Không có lỗi" và "đã nạp" là hai chuyện khác nhau; im lặng thì chủ dự án sẽ sửa nhầm file
   rồi tưởng đã áp dụng.
9. **Đột biến tay phải bỏ ĐỦ các lớp chặn.** 03/09 một đột biến D4 "giết được" nhưng đỏ ở
   happy-path chứ không đỏ ở test D4 (§8 BAT-DAU). Với KH8, nhiều bất biến có **hai** lớp
   (validator lúc nạp + kiểm biên lúc dựng gói) — một đột biến chỉ có nghĩa khi **đúng** test
   mang tính chất đó đỏ, các test khác vẫn xanh.
10. **`catalog.json` panic lúc khởi động nếu hash lệch** (`catalog.go:213-219`). Task nào lỡ
    format lại file đó là máy chủ **không chạy được**, và thông báo lỗi (*"catalog hash khong
    khop asset lock"*) không gợi gì về nguyên nhân thật.
