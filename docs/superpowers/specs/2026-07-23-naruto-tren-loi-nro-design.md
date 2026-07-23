# Thiết kế: Game Naruto trên lõi Ngọc Rồng

**Ngày:** 2026-07-23
**Trạng thái:** Chờ duyệt
**Người quyết định:** Chủ dự án (non-tech) — mọi giải thích viết cho người không lập trình

---

## 1. Tóm tắt

Xây dựng một game nhập vai trực tuyến chủ đề **Naruto**, dùng lại **lõi cơ chế và cảm giác điều khiển** của server Ngọc Rồng hiện có, nhưng:

- Thay **toàn bộ nội dung**: cốt truyện, làng, nhân vật, kỹ năng, vật phẩm, mỹ thuật
- Thay **toàn bộ nền tảng kỹ thuật**: bỏ Java và PHP, chuyển sang Go + PostgreSQL + Next.js
- Giữ **client Unity hiện có** (sửa, không viết lại) để bảo toàn cảm giác chơi

Phân công: **chủ dự án lo mỹ thuật và hoạt ảnh. Trợ lý lo toàn bộ code.**

---

## 2. Bối cảnh — hệ thống hiện tại

Các số liệu dưới đây **đo trực tiếp từ mã nguồn và cơ sở dữ liệu**, không phải ước lượng.

### 2.1 Quy mô

| Hạng mục | Số lượng |
|---|---|
| File mã nguồn Java | 489 |
| Bản đồ | 183 |
| Loại quái | 117 |
| Boss | ~110 |
| Vật phẩm | 1.771 |
| Loại chỉ số trên trang bị | 255 |
| Kỹ năng | 21 (8 mỗi chủng tộc) |
| Nội tại | 25 |
| Danh hiệu | 21 bậc |
| Nhiệm vụ chính | 36 |
| Chức năng chế tác | 37 |
| File ảnh | ~70.600 (820 MB) |

### 2.2 Kiến trúc hiện tại

- **Server game**: Java 21, bắt buộc có giao diện đồ hoạ Swing → khoá vào Windows
- **Web**: PHP 8 trên XAMPP/Apache
- **CSDL**: MariaDB, schema `z`, 48 bảng
- **Client**: Unity 2022.3.62f2, IL2CPP, 424 file C# (212 × 2 bản trùng lặp Game1/Game2)
- **Liên lạc web ↔ game**: không có API. Chỉ qua (a) CSDL dùng chung, (b) file `panel_cmd.txt` poll 2 giây, (c) file JSON trạng thái

### 2.3 Các vấn đề đã xác định

| Vấn đề | Bằng chứng |
|---|---|
| Mỗi người chơi một luồng hệ điều hành | `Controller.java:1140` → `player.start()` |
| Tranh chấp dữ liệu giữa các luồng | `CopyOnWriteArrayList` rải khắp, nhiều ghi chú `[FIX-N]` về `ConcurrentModificationException` |
| Mã vật phẩm là **vị trí trong danh sách** | `ItemService.getTemplate(id)` → `ITEM_TEMPLATES.get(id)`. Xoá 1 dòng = lệch toàn bộ rương người chơi |
| Nạp vật phẩm không có `ORDER BY` | `Manager.java:606` — đang đúng nhờ may mắn |
| Dữ liệu nhân vật là JSON nhồi trong ô text | Bảng `player`, ~40 cột `text` |
| Vật phẩm không có danh tính | Lưu dạng `[mã, số lượng, chỉ số]` → không phát hiện được nhân bản |
| Lỗi bị nuốt im lặng | `Manager.java:338` — `catch (Exception ex) { }` rỗng bọc quanh nạp clan/giftcode/ký gửi |
| **Server tin tuyệt đối vị trí client báo** | `Controller.java:413` (`case -7`) — chỉ map Ngọc Rồng Đen có kiểm khoảng cách |
| Khoá xác thực client chưa đổi | `ClientVerifier.java:7` — `"NRO_SECRET_KEY_2024_CHANGE_ME"` |
| Chặn IP mất khi khởi động lại, và chỉ chạy trên Windows | `AntiSpam.java` — dùng `netsh`, lưu trong bộ nhớ |
| Tự khoá tài khoản khi phát hiện trùng | `Client.java:114` — rủi ro khoá oan người chơi thật |
| Không có git, không có kiểm thử tự động, không có tài liệu | — |

### 2.4 Những điều hệ thống hiện tại làm ĐÚNG (phải giữ)

- **Sát thương do server tính**, client không quyết được
- **Server kiểm tra tầm đánh, hồi chiêu, khí**
- Mật khẩu băm PBKDF2-HMAC-SHA256, 100.000 vòng
- Câu báo lỗi đăng nhập không tiết lộ tài khoản nào tồn tại
- Cơ chế `admin_checking` khoá tài khoản trước khi admin sửa dữ liệu
- Tự động lưu 5 phút/lần (`GiaCatLuongManager.startAutoSave`)
- Hệ thống thay thế biến trong lời thoại (`%1`–`%14`) cho phép một nhiệm vụ dùng chung nhiều chủng tộc

---

## 3. Mục tiêu

### 3.1 Mục tiêu

1. **Chịu tải**: vài nghìn người chơi đồng thời, không nghẽn
2. **Mượt**: thao tác không giật, không chờ cơ sở dữ liệu
3. **Ít lỗi**: loại bỏ tận gốc lỗi tranh chấp dữ liệu
4. **Nhẹ cho người chơi**: APK ~35–40 MB, kho ảnh ~180 MB
5. **Dễ nâng cấp**: thêm nội dung không cần lập trình, không cần tắt server
6. **Chống gian lận**: server không tin client bất cứ điều gì

### 3.2 Phi mục tiêu (không làm trong phạm vi này)

- Viết lại client Unity từ đầu
- Hỗ trợ nhiều server song song (sharding)
- Ứng dụng di động riêng cho quản trị
- Hệ thống nạp tiền tự động

---

## 4. Kiến trúc tổng thể

```
                    ┌──────────────────────┐
   Người chơi ─────►│   CLIENT UNITY       │
                    └───┬──────────────┬───┘
                        │              │ tải ảnh (HTTPS)
              đường dây │              ▼
              game (TCP)│      ┌───────────────┐
                        │      │ CLOUDFLARE R2 │
                        │      │  + CDN        │
                        │      └───────────────┘
                        ▼
        ┌───────────────────────────────┐
        │   MÁY CHỦ GAME  (Go, Linux)   │
        │  ├ Cổng vào — lớp phiên dịch  │
        │  ├ Khu vực 1..N — mỗi khu 1   │
        │  │   goroutine sở hữu dữ liệu │
        │  ├ Kho nội dung (chỉ đọc)     │
        │  ├ Hàng chờ ghi CSDL          │
        │  └ API quản trị (HTTP)        │
        └────┬──────────────────────┬───┘
             ▼                      ▼
   ┌──────────────────┐   ┌──────────────────────┐
   │    SUPABASE      │◄──│  WEB + PANEL         │
   │   (PostgreSQL)   │   │  (Next.js / Vercel)  │
   └──────────────────┘   └──────────────────────┘
```

### 4.1 Sáu thành phần

| Thành phần | Nơi chạy | Trách nhiệm | Khi hỏng |
|---|---|---|---|
| Client Unity | Máy người chơi | Vẽ hình, nhận phím. Không tính toán luật chơi | Chỉ người đó bị ảnh hưởng |
| Máy chủ game | VPS Linux | Toàn bộ luật chơi | Tự bật lại ~10s, client tự nối lại |
| Supabase | Đám mây | Lưu trữ bền vững | Game vẫn chạy, ghi xếp hàng chờ |
| Cloudflare R2 + CDN | Toàn cầu | Kho ảnh, hoạt ảnh | Client lấy dự phòng từ máy chủ game |
| Web + Panel | Vercel | Trang chủ, quản trị | Không ảnh hưởng game |
| Kho công cụ | GitHub / máy dev | Đóng gói nội dung và ảnh | Không ảnh hưởng người chơi |

### 4.2 Quyết định kiến trúc 1 — Mỗi khu vực một chủ sở hữu

**Vấn đề:** hiện tại nhiều luồng cùng sửa một cấu trúc dữ liệu → lỗi tranh chấp, đã vá chắp vá bằng `CopyOnWriteArrayList` và `synchronized`.

**Giải pháp:** mỗi khu vực (zone) do **đúng một goroutine sở hữu**. Mọi thay đổi trạng thái khu vực đi qua hộp thư (channel) của goroutine đó, xử lý tuần tự.

**Hệ quả:**
- Lỗi tranh chấp dữ liệu **không thể xảy ra về mặt thiết kế**, không phụ thuộc kỷ luật lập trình
- Các khu vực chạy song song thật sự, không chờ nhau
- Một khu vực panic → `recover`, khởi tạo lại khu vực đó từ trạng thái đã lưu; các khu khác không ảnh hưởng

**Ràng buộc:** người chơi di chuyển giữa khu vực = bàn giao quyền sở hữu qua tin nhắn, không phải truy cập chéo.

### 4.3 Quyết định kiến trúc 2 — Nội dung tách khỏi code

Nội dung game chuyển từ code sang **file dữ liệu trong GitHub**:

```yaml
# noi-dung/boss/zabuza.yaml
boss:
  ma: zabuza
  ten: "Zabuza Momochi"
  lang_xuat_hien: [lang_suong_mu, lang_la]
  mau: 500000000
  chieu: [suong_mu_an_than, thuy_long_dan]
  roi_do:
    - { vat_pham: kiem_kubikiribocho, ti_le: 2 }
```

Sửa file → đẩy GitHub → **máy chủ nạp lại nóng**, không tắt server.

**Ngoại lệ cố ý:** các con số cần chỉnh nóng từ panel (tỉ lệ EXP, tỉ lệ rơi, giá shop, bật/tắt sự kiện) nằm trong CSDL.

### 4.4 Quyết định kiến trúc 3 — Ghi CSDL không chặn game

| Loại dữ liệu | Cách ghi | Lý do |
|---|---|---|
| Vị trí, máu, chakra, kinh nghiệm | Ghi dồn, 1–5 phút/lần | Mất vài phút chấp nhận được |
| Giao dịch, mua bán, nhận đồ hiếm, ký gửi | **Ghi ngay trong giao dịch CSDL** | Mất là mất đồ thật |

Ghi hỏng → trả lại hàng chờ, tự thử lại. Không mất dữ liệu, không chặn game.

### 4.5 Quyết định kiến trúc 4 — Lớp phiên dịch giao thức

Giao thức nhị phân cũ (byte lệnh + mã hoá khoá chuỗi) **gói gọn trong một lớp duy nhất ở cổng vào**. Bên trong máy chủ chỉ làm việc với khái niệm nghiệp vụ.

**Lợi ích:** ngày muốn viết client mới (Godot, web, Unity mới) chỉ cần thêm lớp phiên dịch thứ hai. Luật chơi không đổi một dòng. Client cũ và mới chạy song song được.

**Cơ sở khả thi:** toàn bộ tầng mạng phía client có mã nguồn C# (`Session_ME.cs` — `getKey`, `readKey`, `writeKey`). Client chính là bản đặc tả giao thức.

---

## 5. Bảo mật, chống gian lận, chống phá game

### 5.1 Lỗ hổng nghiêm trọng cần sửa

**Server tin tuyệt đối vị trí client báo.** `Controller.java:413` đọc toạ độ từ client và dùng luôn. Chuỗi hệ quả:

```
Client sửa đổi báo toạ độ giả (sát boss)
  → Server tin
  → Kiểm tra tầm đánh ĐẠT (dựa trên toạ độ giả)
  → Đánh xong, báo toạ độ khác → biến mất
```

Kiểm tra tầm đánh (vốn được cài đúng) trở nên vô nghĩa. Cho phép: dịch chuyển tức thời, chạy nhanh, đánh ngoài tầm, cày tự động.

### 5.2 Năm lớp phòng thủ

**Lớp 1 — Chặn từ ngoài (DDoS)**

| Phần | Cơ chế |
|---|---|
| Website + panel | Cloudflare (gói hiện có) |
| Cổng game (TCP) | Chống DDoS của nhà cung cấp VPS. *Cloudflare gói thường KHÔNG chắn được TCP game* |
| Kho ảnh | Cloudflare R2 + CDN — tách khỏi VPS, giảm tải và giảm bề mặt tấn công |

Trong server: giới hạn kết nối/IP, ngắt kết nối chưa xác thực sau vài giây, giới hạn gói tin/giây, **danh sách chặn lưu xuống CSDL** (không mất khi khởi động lại), chặn **có thời hạn** thay vì vĩnh viễn.

**Lớp 2 — Server không tin client**

Nguyên tắc: client gửi **ý định**, không gửi **kết quả**.

| Client gửi | Server làm |
|---|---|
| "muốn đi tới (x,y)" | Tự tính đường; nếu quãng đường vượt tốc độ cho phép trong khoảng thời gian đó → kéo về vị trí thật, ghi nhận vi phạm |
| "đánh quái #5" | Tự kiểm vị trí thật, tầm, hồi chiêu, chakra → tự tính sát thương |
| "nhặt món đồ" | Tự kiểm món có thật, khoảng cách, quyền sở hữu |
| "mua món này" | Tự kiểm giá, số dư, tự trừ |

Client vẫn vẽ chuyển động ngay khi bấm (dự đoán phía client) để giữ cảm giác mượt; server kiểm tra ở phía sau.

**Lớp 3 — Vật phẩm có danh tính**

Mỗi món đồ có **số căn cước duy nhất, không tái sử dụng**, kèm nhật ký đổi chủ đầy đủ. Nhân bản gần như không thể; nếu xảy ra thì truy ngược và thu hồi chính xác được.

**Lớp 4 — Phát hiện bất thường bằng số liệu**

Đánh dấu (không tự khoá): tốc độ giết quái bất khả thi, thu nhập vượt mặt bằng, online quá dài liên tục, khoảng cách thao tác đều tăm tắp (dấu hiệu bot), bị kéo về vị trí nhiều lần.

**Kết quả đưa lên panel kèm bằng chứng, người quản trị quyết định.** Cố ý bỏ cơ chế tự khoá hiện tại vì rủi ro khoá oan.

**Lớp 5 — Chống phá game giữa người chơi**

| Kiểu phá | Xử lý |
|---|---|
| Spam chat, ngôn từ xấu | Giới hạn tần suất, lọc từ, cấm chat tạm thời |
| Cướp boss | Boss ghi nhận đóng góp sát thương; đồ rơi ưu tiên theo đóng góp |
| Nhiều tài khoản một người | Giới hạn theo IP + dấu vân máy; đánh dấu để admin xem |
| Lừa đảo giao dịch | Xác nhận hai bên + nhật ký đầy đủ để tra khi khiếu nại |

### 5.3 Điều không thể làm được

**Không thể ngăn người chơi sửa client.** Client Unity IL2CPP đều dịch ngược được. Mọi biện pháp xác thực phía client chỉ nâng rào cản, không chặn được người quyết tâm.

Chiến lược đúng: **làm cho việc sửa client trở nên vô ích** — server tính hết, không tin gì, ghi nhật ký mọi thứ có giá trị.

---

## 6. Mô hình dữ liệu

### 6.1 Bốn nguyên tắc

1. **Tách nội dung game (file GitHub) khỏi dữ liệu người chơi (CSDL)**
2. **Mã định danh là chuỗi chữ, không phải vị trí** — kèm bảng phiên dịch chữ ↔ số cho client cũ
3. **Mỗi món đồ một số căn cước duy nhất** + nhật ký đổi chủ
4. **Thứ quý ghi ngay, thứ rẻ ghi dồn**

### 6.2 Cấu trúc bảng

**Nhóm người chơi**

| Bảng | Nội dung |
|---|---|
| `tai_khoan` | Tên đăng nhập, mật khẩu băm, email, tình trạng khoá |
| `nhan_vat` | Tên, làng, vị trí, chỉ số, tiền — **các cột riêng biệt, tra cứu được** |
| `nhan_vat_ky_nang` | Mỗi kỹ năng đã học một dòng, kèm cấp |
| `nhan_vat_nhiem_vu` | Tiến độ từng nhiệm vụ |
| `nhan_vat_quan_he` | Bạn bè / kẻ thù |

**Nhóm vật phẩm**

| Bảng | Nội dung |
|---|---|
| `vat_pham` | Mỗi **món đồ** một dòng: số căn cước, mã mẫu, chủ sở hữu, vị trí chứa (túi/rương/hòm thư/chợ), số lượng, chỉ số đã roll (JSONB) |
| `vat_pham_lich_su` | Mỗi lần đổi chủ một dòng: từ ai, sang ai, lý do, thời điểm |

**Nhóm vận hành**

| Bảng | Nội dung |
|---|---|
| `cau_hinh_game` | Tỉ lệ EXP/rơi đồ/giá/sự kiện — panel sửa nóng |
| `nhat_ky_quan_tri` | Ai làm gì trên panel, khi nào |
| `canh_bao_gian_lan` | Nghi vấn kèm bằng chứng, chờ duyệt |
| `phien_dich_ma` | Mối nối `mã_chữ ↔ số` cho client cũ, **vĩnh viễn không đổi** |
| `chan_ip` | Danh sách chặn, có thời hạn |

**Không còn bảng nào chứa** nhiệm vụ, boss, bản đồ, định nghĩa vật phẩm — chuyển sang file GitHub.

### 6.3 Chuyển dữ liệu cũ

1. Chụp bản sao CSDL hiện tại, khoá lại
2. **Đếm trước:** tổng vàng, tổng ngọc, số món đồ theo loại, số nhân vật, tổng sức mạnh
3. Chuyển: bung JSON thành dòng, cấp số căn cước cho từng món
4. **Đếm lại và đối chiếu tuyệt đối.** Lệch bất kỳ con số nào → dừng, tìm nguyên nhân
5. Chạy song song cũ/mới trên cùng dữ liệu, so sánh hành vi trước khi chuyển hẳn

Kèm **bảng quy đổi** vật phẩm Dragon Ball → Naruto: đổi tên và hình, giữ nguyên chỉ số và giá trị.

---

## 7. Nội dung Naruto

### 7.1 Ánh xạ khung NRO → Naruto

| Khung NRO | Naruto |
|---|---|
| Đánh 5 mộc nhân | Ném kunai vào bia ở Học viện |
| Học bay | **Tụ chakra dưới chân** (leo cây, chạy trên nước) |
| Bái sư | **Nhận đội 3 người + Jonin hướng dẫn** |
| Tàu vũ trụ | **Cổng làng + đường nhẫn giả** |
| Bang hội | **Gia tộc / tổ đội** |
| Sức mạnh | **Nhẫn lực** |
| Tiềm năng | **Điểm tu luyện** |
| KI | **Chakra** |
| Danh hiệu | **Cấp bậc nhẫn giả** |
| Đại hội võ thuật | **Kỳ thi Chunin** |
| Ngọc Rồng (7 viên) | **Bảy Cuộn Thần Chú** |
| Đệ tử / hợp thể | **Thú triệu hồi / hợp nhất yêu thú** |

### 7.2 Năm làng

| Làng | Hệ | Vũ khí | Thiên về | Kế thừa |
|---|---|---|---|---|
| Làng Lá (Konoha) | Hỏa | Kunai, thủ ấn | Cân bằng | Trái Đất |
| Làng Cát (Suna) | Phong | Con rối | Chakra cao, đánh xa | Namếc |
| Làng Sương Mù (Kiri) | Thủy | Đại kiếm | Sát thương cận chiến | Xayda |
| Làng Đá (Iwa) | Thổ | Giáp đá, chuỳ | Thể lực, phòng thủ | *mới* |
| Làng Mây (Kumo) | Lôi | Song đao | Tốc độ, chí mạng | *mới* |

**Số làng không bị giới hạn ở 3.** Ràng buộc "đúng 3 chủng tộc" là của server Java (690 chỗ dùng biến `gender` trong 74 file). Client chỉ dùng 19 chỗ trong 3 file, và giao thức đã hỗ trợ N chủng tộc (`new NClass[d.readByte()]`). Server mới thiết kế cho N làng ngay từ đầu.

### 7.3 Cấu trúc 76 nhiệm vụ

**Chặng 1–10: mỗi làng một câu chuyện riêng (5 × 10 = 50)**

*Làng Lá:* Ném kunai → Thi tốt nghiệp (Phân Thân) → Nhận đội và Jonin → Bắt mèo Tora → **Leo cây (mở tụ chakra)** → Hộ tống tới Xứ Sóng → **Zabuza & Haku trên cầu** → Huấn luyện nâng cao → Điều tra Rừng Tử Thần → Đăng ký Kỳ thi Chunin

*Làng Cát:* Điều khiển cát → Sống sót bão cát → Nhận đội → Hộ tống thương đoàn → **Học Thuật Con Rối** → Dẹp cướp Hố Cát → **Chạm trán Gaara** → Bí mật của Kazekage → Chặn quái sa mạc → Lên đường dự thi

*Làng Sương Mù:* Rèn kiếm trong sương → Sống sót Đảo Sương → Nhận đội → Truy bắt nhẫn giả đào tẩu → **Học Kiếm thuật Thất Đao** → Cứu trẻ huyết kế → **Truy đuổi Zabuza** *(cùng nhân vật, góc nhìn ngược Làng Lá)* → Mizukage bất thường → Dẹp loạn nội bộ → Lên đường dự thi

*Làng Đá:* Rèn trong hầm mỏ → Nâng tảng đá → Nhận đội → Bảo vệ mỏ → **Học Thổ độn** → Nghe kể Đại chiến lần 3 và "Tia Chớp Vàng" → Bắt gián điệp → Onoki tuyên bố cử người → Vượt thử thách Onoki → Lên đường mang định kiến

*Làng Mây:* Chạy vách núi → Bắt tia sét → Nhận đội → Hộ tống đoàn xe → **Học Lôi độn** → **Gặp Killer Bee và Bát Vĩ** → Điều tra vụ mất tích Hyuga → Raikage A thử sức → Chặn bão sét → Lên đường dự thi

**Chặng 11–36: chung cho cả năm làng (26)**

| # | Nhiệm vụ | Boss |
|---|---|---|
| 11 | Kỳ thi Chunin vòng 1 | |
| 12 | Rừng Tử Thần — tranh cuộn Trời/Đất | |
| 13 | Orochimaru xuất hiện — Dấu Nguyền | Orochimaru |
| 14 | Vòng loại sơ bộ | |
| 15 | Một tháng huấn luyện | |
| 16 | Vòng chung kết đấu trường | |
| 17 | Konoha bị đánh úp | Gaara–Shukaku |
| 18 | Hokage đệ tam hy sinh | Nhị Kage (Edo Tensei) |
| 19 | Truy tìm Tsunade | Kabuto |
| 20 | Sasuke đào tẩu — đội truy bắt | Tứ Âm |
| 21 | Thung Lũng Tận Cùng | Sasuke Dấu Nguyền 2 |
| — | *Gián đoạn 2,5 năm* | |
| 22 | Akatsuki lộ diện | |
| 23 | Giải cứu Kazekage | Sasori, Deidara |
| 24 | Lần theo dấu Sasuke | Orochimaru |
| 25 | Cái chết của Asuma | Hidan, Kakuzu |
| 26 | Itachi và sự thật Uchiha | Itachi |
| 27 | Jiraiya về Làng Mưa | Sáu Đạo Pain |
| 28 | Pain san phẳng Làng Lá | Tendo Pain |
| 29 | Hiền Nhân Thuật | Pain hoàn chỉnh |
| 30 | Hội nghị Ngũ Kage bị phá | Sasuke, Tobi |
| 31 | Liên minh Nhẫn giả | |
| 32 | Đại chiến bắt đầu | Edo Tensei quân đoàn |
| 33 | Mười Vĩ hồi sinh | Obito Lục Đạo |
| 34 | Madara Uchiha | Madara Lục Đạo |
| 35 | Kaguya Otsutsuki | Kaguya |
| 36 | Thời đại mới | |

### 7.4 Thang bậc nhẫn giả (thay 21 danh hiệu)

| Bậc | Cấp | Nhẫn lực |
|---|---|---|
| Học viên | 1–3 | 1.000 → 15.000 |
| Genin | 4–7 | 40.000 → 340.000 |
| Chunin | 8–9 | 700.000 → 1,5 triệu |
| Jonin | 10–11 | 15 triệu → 150 triệu |
| Jonin Tinh Anh | 12–13 | 1,5 tỷ → 5 tỷ |
| ANBU | 14–15 | 10 tỷ → 40 tỷ |
| Kage | 16–18 | 50 tỷ → 70 tỷ |
| Hiền Nhân | 19–20 | 80 tỷ → 90 tỷ |
| Lục Đạo Tiên Nhân | 21 | 100 tỷ |

---

## 8. Hệ thống mỹ thuật

### 8.1 Quy tắc kỹ thuật (đo từ game)

Mọi ảnh tồn tại ở các cỡ **gốc × 1, × 2, × 3, × 4** (icon chỉ dùng x2–x4). Kiểm chứng trên 7 mã icon: icon 674 = gốc 15px → 30×32, 45×48, 60×64.

| Quy tắc | Chi tiết |
|---|---|
| Vẽ một lần | Chỉ vẽ bản **x4**; công cụ tự sinh x3 (×0,75), x2 (×0,5), x1 (×0,25) |
| Nền | PNG 32-bit trong suốt |
| Hoạt ảnh | Nhiều khung trên một tấm sprite + bảng toạ độ |

### 8.2 Phong cách bắt buộc

Icon rất nhỏ (x4 chỉ 40–64px). Yêu cầu: **hình khối lớn, ít chi tiết, viền dày, tương phản mạnh, bảng màu hẹp.** Phong cách chọn: **anime-pixel lai** — nét dày kiểu manga, tô phẳng 3–4 sắc độ, viền tối rõ.

### 8.3 Khối lượng đợt 1

| Nhóm | Số tấm |
|---|---|
| Giao diện | ~60 |
| Vật phẩm | ~120 |
| Trang phục Làng Lá | 36 |
| Vũ khí | ~20 |
| Quái | 12 |
| Boss | 6 |
| NPC | 15 |
| Hiệu ứng chiêu | 10 |
| Mảnh ghép bản đồ | ~80 |
| **Tổng** | **~360** |

Ba đợt: (1) Làng Lá ~400 tấm → chơi thử tới Kỳ thi Chunin; (2) bốn làng còn lại ~800; (3) Shippuden ~700.

### 8.4 Đường ống sản xuất

```
1. Dán prompt vào AI            → ảnh
2. Tách nền, cắt sát viền       → công cụ
3. Chuẩn hoá về cỡ x4           → công cụ
4. Sinh x3, x2, x1              → công cụ
5. Gộp sprite + bảng toạ độ     → công cụ
6. Đẩy lên Cloudflare R2        → công cụ
7. Game nhận ngay               → không build lại
```

Bước 5 khả thi vì **định dạng file hiệu ứng của client đã giải mã được** từ `EffectData.cs`:

```
byte   số ảnh con
lặp:   ID, x0, y0, rộng, cao
short  số khung
lặp:   số mảnh; lặp: dx, dy, id_ảnh_con
```

Đọc được thì ghi được → công cụ đóng gói hiệu ứng tự động.

### 8.5 Tối ưu dung lượng

| | Hiện tại | Sau |
|---|---|---|
| Số bản phải vẽ mỗi thứ | 3–4 cỡ | **1 cỡ** |
| Định dạng | PNG | WebP (−40%) |
| Tổ chức | 21.000 file lẻ | Gộp sprite |
| Kho ảnh | 820 MB | **~180 MB** |
| APK | 71,6 MB | mục tiêu **35–40 MB** |
| Bản PC | 274 MB | mục tiêu **~120 MB** |

*Mục tiêu APK/PC là ước tính, phải đo lại sau khi làm.*

---

## 9. Hạ tầng và triển khai

| Thành phần | Dịch vụ | Chi phí/tháng |
|---|---|---|
| Máy chủ game | VPS Linux (4 nhân, 8 GB) | 10–20 USD |
| CSDL | Supabase | 0 → 25 USD |
| Web + panel | Vercel | 0 USD |
| Mã nguồn | GitHub | 0 USD |
| Kho ảnh | Cloudflare R2 (**không tính phí băng thông ra**) | ~0–5 USD |

### 9.1 Cách cập nhật từng loại thay đổi

| Thay đổi | Cách lên | Người chơi |
|---|---|---|
| Web, panel | Đẩy GitHub → Vercel tự build | Tải lại trang |
| Nội dung game | Sửa file → nạp nóng | Không làm gì |
| Ảnh, hoạt ảnh | Đẩy lên R2 | Không làm gì |
| Code máy chủ | Đẩy GitHub → build → khởi động lại | Ngắt ~10s, tự nối lại |
| Cơ chế mới trong client | Build lại Unity | **Phải cài lại** |

**Nguyên tắc thiết kế:** tối đa hoá số thay đổi rơi vào ba dòng đầu.

### 9.2 Quy trình phát hành

Đẩy GitHub → tự động triển khai **server thử** → chủ dự án kiểm tra → phê duyệt → triển khai server thật.

---

## 10. Rủi ro

| Rủi ro | Mức | Giảm thiểu |
|---|---|---|
| Chuyển dữ liệu cũ làm mất/lệch tài sản người chơi | **Cao** | Đếm–đối chiếu tuyệt đối; chạy song song trước khi chuyển hẳn |
| Khối lượng mỹ thuật vượt sức | **Cao** | Chia 3 đợt; đợt 1 chỉ 1 làng để kiểm chứng đường ống |
| Lệch giao thức giữa Go và client C# | Trung bình | Một file mô tả giao thức duy nhất, sinh code cho cả hai phía |
| Bản quyền Naruto (Shueisha/Pierrot) | Trung bình | Cùng mức rủi ro với Dragon Ball hiện tại; thiết kế "lấy cảm hứng", không sao chép nguyên tạo hình |
| Phong cách ảnh AI không đồng nhất | Trung bình | Khối phong cách bắt buộc dùng chung; duyệt theo lô |
| Supabase mất kết nối | Thấp | Hàng chờ ghi + tự thử lại |

---

## 11. Phụ lục — Tham chiếu mã nguồn đã kiểm chứng

| Sự kiện | Vị trí |
|---|---|
| Một luồng mỗi người chơi | `Controller.java:1140` |
| Vòng lặp game 100ms | `GameLoop.java:8` |
| Mã vật phẩm là vị trí | `ItemService.java:234-236` |
| Nạp vật phẩm thiếu ORDER BY | `Manager.java:606` |
| Nuốt lỗi im lặng | `Manager.java:338` |
| Không kiểm tra di chuyển | `Controller.java:413-441` |
| Có kiểm tra tầm đánh | `SkillService.java:579-582` |
| Có kiểm tra hồi chiêu | `SkillService.java:1273` |
| Khoá xác thực chưa đổi | `ClientVerifier.java:7` |
| Chặn IP chỉ chạy Windows | `AntiSpam.java` (dùng `netsh`) |
| Tự khoá tài khoản trùng | `Client.java:114` |
| Tự động lưu 5 phút | `GiaCatLuongManager.java:216-236` |
| Client đọc số chủng tộc động | `Controller.cs:4658` |
| Client chỉ dùng `gender` 19 lần | 3 file: `Panel.cs`, `Service.cs`, `ItemTemplate.cs` |
| Định dạng file hiệu ứng | `EffectData.cs:82-125` |
| Tầng mạng client có mã nguồn | `Session_ME.cs:129-160` |
