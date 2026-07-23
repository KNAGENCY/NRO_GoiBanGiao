# Thiết kế: Dựng chạy server NRO trên máy local (Giai đoạn 1)

**Ngày:** 2026-07-21
**Trạng thái:** Đã được chủ dự án phê duyệt

## Bối cảnh

Chủ dự án mua lại một bộ source server private Ngọc Rồng Online (server Java + web PHP + client Unity + dump database). Mục tiêu cuối là **vận hành thương mại**, nhưng giai đoạn đầu chỉ mở giới hạn cho 10–50 người test. Chủ dự án không phải lập trình viên.

Hiện chưa có VPS. Chủ dự án sẽ tự thuê và cung cấp thông tin sau.

## Mục tiêu giai đoạn này

Cho toàn bộ hệ thống chạy được **trên máy local**, chỉ chủ dự án truy cập được, không mở ra Internet.

**Tiêu chí nghiệm thu:** chủ dự án tự đăng ký tài khoản qua web trên máy mình, mở `Dragon Boy.exe`, đăng nhập thành công, tạo được nhân vật và di chuyển trong game.

## Ngoài phạm vi (cố ý hoãn, kèm lý do)

| Việc | Vì sao hoãn |
|---|---|
| Sửa mật khẩu plaintext → hash | Phải sửa đồng thời ở Java và PHP; sai là toàn bộ tài khoản mất đăng nhập. Giai đoạn này chưa có mật khẩu người thật nào để lộ, và chưa có ai để kiểm chứng sau khi sửa. Làm ở giai đoạn chuẩn bị mở cửa. |
| Đặt mật khẩu cho database | Cùng lý do — chỉ có ý nghĩa khi máy tiếp xúc mạng ngoài. |
| Xoá dấu vết chủ cũ (SĐT, email, tên miền) | Chưa ai nhìn thấy. Làm khi chuẩn bị đón người test. |
| Vá địa chỉ server trong client | Client đã trỏ sẵn `127.0.0.1:14445`, đúng bằng địa chỉ server local. Không cần đụng. |
| Mở port ra Internet | Chính việc **không** mở port là thứ giữ cho giai đoạn này an toàn tuyệt đối. |
| Build lại từ source bằng Ant | Sai thứ tự chẩn đoán — xem phần Quyết định dưới. |

## Quyết định kỹ thuật

### Chạy `zz1.jar` có sẵn, không build lại

**Chọn:** dùng nguyên file `zz1.jar` ở thư mục gốc.

**Lý do:** đây là artifact mà chủ cũ đang chạy thật (build 2026-07-20 17:07, mới hơn file source mới nhất cùng ngày 16:43 → jar khớp source). Nếu hệ thống hỏng, ta biết chắc nguyên nhân là cấu hình chứ không phải do ta build sai.

**Phương án loại bỏ:** build lại bằng Ant trước. Bị loại vì (a) build đang hỏng sẵn — `nbproject/project.properties:32` cần `lib\LIB_SALE.jar` nhưng trong `lib/` chỉ có `ga.jar` (hai file giống hệt, chỉ khác tên); (b) sản phẩm build là `dist/NROKRAI.jar`, **khác** file đang chạy; (c) nếu lỗi thì không phân biệt được do source hỏng hay do dựng sai.

**Phương án loại bỏ:** chạy `toolcaigame.exe`. Bị loại vì máy đã có sẵn JDK 21 và XAMPP nên tool không đem lại gì thêm, trong khi nó chạy quyền admin, sửa PATH hệ thống, đổi múi giờ, đổi trình duyệt mặc định, và tải một archive không kiểm chứng được từ Google Drive cá nhân. Rủi ro không đổi lấy được lợi ích.

### Thứ tự khởi động: Database → Server Java → Web → Client

Không phải tuỳ ý. `Manager.loadProperties()` gọi `System.exit()` nếu không đọc được cấu hình, và `loadDatabase()` nạp toàn bộ item/mob/npc/map/skill từ MySQL ngay lúc khởi động — **không có database thì server tự thoát**. Web và client đều vô nghĩa nếu server chưa sống.

## Các khối việc

### 1. Database
Bật MySQL trong XAMPP. Tạo schema `z` với charset `utf8mb4`. Import `zz1\z.sql` (856KB, 46 bảng). Xác nhận đủ bảng và bảng `player`/`account` có đúng 1 dòng dữ liệu mẫu.

### 2. Server Java
Sửa `zz1\run.bat` để gọi `C:\Program Files\Java\jdk-21\bin\java.exe` thay vì `java` trên PATH (hiện PATH trỏ JRE 1.8.0_441, trong khi `zz1.jar` biên dịch bằng Java 21 — bytecode major version 65). Chạy, xác nhận cửa sổ Swing hiện lên, socket chiếm được cổng 14445, và log báo nạp xong dữ liệu.

Đổi `server.sv1` trong `data\config\config.properties` từ `NRO:fw.vpsttt.earth:10473:0,0,0` thành `NRO:127.0.0.1:14445:0,0,0`. Lý do: đây là danh sách server mà server gửi cho client hiển thị; để nguyên địa chỉ VPS chủ cũ vừa vô nghĩa vừa gây nhầm lẫn khi chẩn đoán. Client vốn đã có entry `Dragon Boy:127.0.0.1:14445:0,0,0` nhúng sẵn nên hai bên sẽ khớp nhau.

### 3. Web
Bật Apache. Mở `http://localhost/sv1/trangchu` (KHÔNG phải `/` — trang gốc là site nguỵ trang "EduVerse"). Đăng ký một tài khoản. Set `is_admin = 1` cho tài khoản đó để vào được `/sv1/admin/panel`.

### 4. Client
Chạy `LocalNew\Dragon Boy.exe`, đăng nhập bằng tài khoản vừa tạo, tạo nhân vật, di chuyển. Không sửa gì.

## Rủi ro đã lường trước

| Rủi ro | Biểu hiện | Cách xử lý |
|---|---|---|
| MariaDB 10.4 không nhận dump | Lỗi syntax khi import | Đọc thông báo lỗi cụ thể; dump vốn được tạo bởi chính MariaDB 10.4.32 nên khả năng thấp |
| Cổng 14445 bị chiếm | `BindException` lúc khởi động | `netstat` tìm tiến trình chiếm cổng |
| Sai thư mục làm việc | Không tìm thấy `data/config/config.properties` | `run.bat` đã có `cd /d "%~dp0"`, giữ nguyên |
| Driver `com.mysql.jdbc.Driver` cũ với MariaDB | Lỗi kết nối lúc khởi động | Driver được đóng gói sẵn trong fat jar; nếu lỗi thì đổi sang chuỗi kết nối tương thích |
| Thư mục game không ghi được | Lỗi khi ghi `data/update_data/part` | Kiểm tra quyền thư mục |

## Vấn đề đã biết, không xử lý ở giai đoạn này

- **Vòng lặp khởi động lại bị kẹt:** `ServerManager.close()` gọi `run.bat` trước `System.exit(0)`, mà `run.bat` kết thúc bằng `PAUSE` → mỗi lần bảo trì đẻ thêm một console chờ bấm phím. Chỉ kích hoạt khi dùng chức năng bảo trì; giai đoạn này không dùng.
- **Server bắt buộc có màn hình đồ hoạ** (`main()` chỉ mở JFrame). Trên máy local Windows không sao. Sẽ là vấn đề lớn khi lên VPS — quyết định ở giai đoạn sau.

## Giai đoạn sau (ghi lại để không quên)

Khi có VPS, các việc bắt buộc trước khi mời người test:
1. Quyết định VPS Windows (server cần màn hình) hoặc dựng X ảo trên Linux
2. Chuyển mật khẩu sang hash — **đồng thời** ở Java và PHP
3. Đặt mật khẩu database, bỏ `root`
4. Xoá toàn bộ dấu vết chủ cũ (danh sách đầy đủ trong `CLAUDE.md`)
5. Vá địa chỉ server trong client — **dùng IP trần**, vì chuỗi trong `global-metadata.dat` chỉ vừa 32 byte và không nới dài được
6. Đóng gói client vào `htdocs\sv1\down\` (hiện đang trống, mọi nút tải đang 404)
7. Sửa `.htaccess` — lớp chống DDoS nạp qua đường dẫn Windows tuyệt đối, trên Linux sẽ im lặng không chạy
