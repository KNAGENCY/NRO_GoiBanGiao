# Ngọc Rồng Kame — Tình trạng dự án & Hướng dẫn mở server

Cập nhật: 2026-07-22. Viết cho người **không phải lập trình viên**.

---

## 1. Đây là gì

Server private Ngọc Rồng Online (NRO / "Dragon Boy"), thương hiệu **Ngọc Rồng
Kame**, tên miền **nrokame.com**. Gồm **4 mảnh dính chặt nhau**:

| Mảnh | Công nghệ | Vai trò |
|---|---|---|
| Server game | Java 21 (project `NROOCEAN`) | Lõi game, chạy trên VPS Windows |
| Web + Panel | PHP 8 (XAMPP/Apache) | Trang chủ, đăng ký, admin panel |
| Client | Unity 2022.3.62f2 | Bản chơi PC / Android / iOS |
| Cơ sở dữ liệu | MariaDB 10.4, schema `z` | Dữ liệu người chơi + cấu hình game |

**Cách 4 mảnh nói chuyện:** không có API. Chỉ qua (1) **DB dùng chung**, (2) file
`panel_cmd.txt` (web ghi lệnh → game đọc), (3) file JSON trạng thái (game ghi →
web đọc). Đây là kiến trúc mong manh — xem mục Điểm yếu.

---

## 2. Đã làm xong (tính đến 22/07/2026)

- **Xoá sạch dấu vết chủ cũ** (187 chỗ: bot quảng cáo, cửa hậu, watermark, SĐT).
- **Đổi thương hiệu** sang Ngọc Rồng Kame / nrokame.com — server, web, client.
- **Bảo mật mật khẩu**: chuyển sang băm PBKDF2 (trước lưu thô).
- **Bản cài người chơi**: PC (`NROKAME.zip`), Android (`NRO.apk`, đã ký khoá riêng),
  iOS (đã xuất dự án Xcode, chờ tài khoản Apple Developer).
- **Sửa 4 lỗi bản PC**: đăng nhập sai không còn treo; đăng ký chỉ ở web; màn chọn
  server 4 nút cố định hết nền đen; icon app mới.
- **VPS đang chạy**: web + game online tại nrokame.com, cổng game 14445.
- **Panel admin**: bịt lỗ hổng quyền; đang restyle đồng nhất theo hệ thiết kế Kame
  (một phần đã xong, một phần đang chạy).
- **Bộ đo giữ chân người chơi** (analytics): đã dựng xong tầng lõi trong sandbox,
  chờ test + deploy trước khi mở server.

---

## 3. Chưa làm / đang treo

- **Chưa có HTTPS** cho web (đang chạy HTTP — mật khẩu đi qua mạng dạng thô).
- **Chưa có sao lưu DB tự động** — rủi ro mất dữ liệu người chơi.
- **iOS**: còn thiếu tài khoản Apple Developer (99 USD/năm).
- **Bộ analytics**: chưa build vào bản server thật + chưa tạo bảng trên VPS.
- **Restyle panel**: còn vài trang lớn đang hoàn tất.
- **Ngày mở server** (`$SITE_OPEN_TIME`) chưa chốt.
- **Windows Server** trên VPS là bản dùng thử, hết hạn khoảng 12/2026.

---

## 4. Điểm yếu khi vận hành (đọc kỹ)

**Nặng — nên xử trước khi đông người:**
- Web HTTP, chưa HTTPS → mật khẩu đăng nhập lộ được trên đường truyền.
- Không sao lưu DB tự động → server hỏng là mất dữ liệu người chơi.
- Một DB dùng chung cho cả game lẫn web → nghẽn cổ chai, lỗi lan chéo.
- Tầng mạng game (`ga.jar`) không có mã nguồn → điểm mù, không tự vá được.

**Vừa:**
- Server bắt buộc có màn hình (Swing) → khoá vào VPS Windows, không headless.
- Không scale ngang (1 tiến trình Java, 1 VPS) → đông quá là nghẽn.
- Web↔game nối bằng file + DB chung, ~15 chỗ hardcode đường dẫn Windows.
- Panel quyền mạnh, chỉ chặn `is_admin` + session, không 2FA.

**Nhẹ:**
- Code cũ chắp vá, có N+1 query, file lớn — bảo trì khó (đang dọn).
- Không test tự động, không staging tách bạch.
- Bản quyền Dragon Ball (Toei/Shueisha) → rủi ro pháp lý, nhất là lên chợ app.
- Mất `nrokame.keystore` = người chơi phải cài lại; mất DB = mất nhân vật.
  **Hai thứ này cần bản dự phòng ngoài máy (USB).**

---

## 5. Mở một server khác từ gói này

Gói zip chia theo khu vực:

```
1-MAY-CHU-GAME/    zz1.jar + mã nguồn Java + config MẪU (mật khẩu trống)
2-WEB/             toàn bộ htdocs (trang chủ + panel)
3-CO-SO-DU-LIEU/   z-moi-nhat.sql (schema + dữ liệu gốc)
4-BO-CAI-DAT/      script cài đặt VPS (bo-cai-vps)
5-TAI-LIEU/        CLAUDE.md, các spec, file này
```

**Các bước tối thiểu để chạy (trên máy Windows có XAMPP + JDK 21):**

1. Import DB: tạo database rồi nạp `3-CO-SO-DU-LIEU/z-moi-nhat.sql`.
2. Điền cấu hình: mở `1-MAY-CHU-GAME/data/config/config.properties` — điền
   `database.name`, `database.user`, `database.pass` của bạn.
3. Điền cấu hình web: `2-WEB/sv1/config.php` và `boss_status.php` — điền mật khẩu DB.
4. Đặt `htdocs` vào thư mục web của XAMPP, bật Apache + MySQL.
5. Chạy server game: đặt `1-MAY-CHU-GAME` làm thư mục gốc, chạy `zz1.jar` bằng
   JDK 21 (`java -jar zz1.jar`), working directory phải là thư mục gốc đó.
6. Client: sửa địa chỉ server trong client rồi build lại (xem 5-TAI-LIEU/CLAUDE.md).

**BẮT BUỘC đọc `5-TAI-LIEU/CLAUDE.md`** — trong đó có mọi cái bẫy: server phải có
màn hình, phải Java 21, kho asset stream theo 4 mức zoom, đổi tên thư mục `sv1` là
đứt liên kết, v.v.

---

## 6. Những thứ CỐ Ý KHÔNG có trong gói (vì bí mật/không nên phát tán)

- **Mật khẩu DB thật** — config trong gói để trống, bạn tự điền.
- **`nrokame.keystore`** (khoá ký APK) — giữ riêng ở máy, không đem đi.
- **Kho asset game nặng** (icon 545MB, img_by_name 137MB...) — xem mục dưới.
- **Source client Unity** (`D:\NRO_Project`, nhiều GB) — copy riêng nếu cần build client.

> Nếu chọn gói "đầy đủ" thì kho asset ĐƯỢC kèm; gói "gọn" thì không.
> Kho asset gốc nằm ở `zz1\data\` (icon, img_by_name, effect, res, mob...).

---

## 7. Liên hệ kỹ thuật

Mọi chi tiết kỹ thuật + cạm bẫy đã ghi trong `5-TAI-LIEU/CLAUDE.md`. Đây là tài
liệu quan trọng nhất để một người mới (hoặc một trợ lý AI khác) tiếp quản dự án.
