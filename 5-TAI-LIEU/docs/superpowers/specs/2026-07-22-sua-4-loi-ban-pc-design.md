# Sửa 4 lỗi bản PC + chuyển đăng ký về web

Ngày: 2026-07-22 · Trạng thái: user đã duyệt hướng, chờ duyệt bản viết

## Bối cảnh

User tải bản PC từ nrokame.com về chạy thử và gặp 4 vấn đề. Khảo sát xong,
tìm ra tổng cộng **5 lỗi thật** (một vấn đề chứa hai lỗi). Quyết định lớn kèm
theo: **bỏ hẳn đăng ký trong game — người chơi đăng ký trên web, vào game chỉ
đăng nhập.**

## Quyết định đã chốt với user

1. Đăng ký **chỉ ở web** (nrokame.com/sv1/dangky). Trong game không còn màn
   đăng ký. Không thêm ô số điện thoại vào game (không còn chỗ nào cần).
2. Màn chọn server có **4 nút cố định**: Đăng nhập · Đổi tài khoản ·
   Máy chủ: Ngoc Rong Kame · Tùy chọn (user đồng ý giữ Tùy chọn).
3. Lỗi đăng nhập sửa theo **cách A**: tái hiện ở máy nhà trước, sửa, thử lại
   có trước–có–sau. Không sửa theo suy luận.
4. Tên file tải bản PC: `PC.zip` → `NROKAME.zip`.

## Lỗi và cách sửa

### 1. Đăng nhập treo "Xin chờ" rồi văng ra (2 lỗi server)

**Lỗi 1a — trả về im lặng.** `zz1\src\jdbc\daos\NDVSqlFetcher.java:80`:
khi mật khẩu không khớp, `return null` mà không gửi gì cho client. Client đợi
câu trả lời không bao giờ đến → treo → tự thoát. Khối `catch` dòng 178 cũng
im lặng y hệt với mọi lỗi bất ngờ khác.

Sửa: mọi đường ra khỏi hàm login đều phải gửi thông báo trước. Sai mật khẩu
gửi: **"Sai tài khoản hoặc mật khẩu. Chưa có tài khoản? Đăng ký tại
nrokame.com"** (câu này kiêm luôn việc dẫn người chơi tới trang đăng ký —
không cần nút mới nào trong client). Khối catch gửi câu chung "Không đăng
nhập được, vui lòng thử lại" + sendLoginFail.

**Lỗi 1b — đăng ký trong game lưu mật khẩu thô.**
`zz1\src\services\Service.java:456` INSERT mật khẩu chưa băm, phá lỗ hổng đã
vá. User đã đăng ký trong game → tài khoản có mật khẩu thô. (Vẫn đăng nhập
được nhờ nhánh so-thô-rồi-nâng-cấp trong PasswordUtil, nhưng nguồn tạo dữ
liệu bẩn phải chặn.)

Sửa: khoá hẳn lệnh đăng ký phía server — `registerAccount()` trả lời
"Vui lòng đăng ký tại nrokame.com" thay vì tạo tài khoản. Chỉ gỡ nút client
là chưa đủ vì client sửa đổi vẫn gọi thẳng được lệnh này. Đồng thời vẫn sửa
chỗ băm cho đúng (`auth.PasswordUtil.hash`) phòng sau này bật lại.

Ghi chú khảo sát: bảng `registration_ip_limits` không có trong DB nhưng KHÔNG
phải lỗi — `ensureRegistrationLimitTable()` tự tạo khi cần. Cột `admin` (khác
`is_admin`) có thật trong bảng account.

### 2. Màn chọn server nền đen + nút lệch (client)

`ServerListScreen.paint` (Game1/Game2) tô đen toàn màn khi `loadScreen` chưa
bật. Nghi ngờ: ảnh nền chưa tải về lần đầu → đen. **Chưa khẳng định nguyên
nhân — phải chạy thử ở máy nhà xem tận mắt rồi mới chọn cách sửa** (có thể là
đóng sẵn ảnh nền vào bản build).

Nút "lệch" thực ra đã căn giữa ngang (`cmd[i].x = (w - cmd[i].w)/2`); cái
lệch là dọc, do công thức vị trí phụ thuộc số nút thay đổi theo `nCmdPlay`.

Sửa: bỏ cơ chế `nCmdPlay` (0/1 tùy có tài khoản lưu) và mọi phép
`cmd[2 + nCmdPlay]`. Mảng nút cố định 4 phần tử: Đăng nhập / Đổi tài khoản /
Máy chủ / Tùy chọn. Cụm nút + logo tính chung một công thức, cân giữa dọc.

### 3. Bỏ màn đăng ký trong game (client)

`AccountRegisterScreen.cs` (Game1/Game2) bỏ. Ô "Tôi đồng ý với Điều khoản"
nằm trong màn này nên hết theo. Mọi chỗ `new AccountRegisterScreen()`
(ServerListScreen dòng ~838, ~880) bỏ hoặc trỏ về màn đăng nhập.
Trang web dieu-khoan.php đã dùng cách "vào chơi là xem như đồng ý" — nhất quán.

### 4. Tên file tải (web)

`PC.zip` → `NROKAME.zip` tại 4 chỗ: `sv1\download.php` (danh sách cho phép),
`sv1\tai-game.php` (thẻ Windows), `sv1\down\.htaccess` (FilesMatch + headers),
và đổi tên file thật trong `sv1\down\`. Làm ở máy nhà, VPS nhận qua DONG-BO.
Lưu ý `tai-game.php` tự kiểm `is_file` — đổi tên xong thẻ Windows vẫn
"Sẵn sàng" vì file mới nằm đúng chỗ.

## Kiểm chứng (cách A)

1. Dựng môi trường thử máy nhà: DB riêng từ `goi-len-vps\z-moi-nhat.sql`,
   game root riêng (mklink /J cho asset nặng), đổi port + tên DB — theo đúng
   quy trình đã ghi trong CLAUDE.md.
2. Build client trỏ localhost, đăng nhập **sai mật khẩu** → phải tái hiện
   được treo. Đây là bằng chứng "trước".
3. Sửa → thử lại: sai mật khẩu hiện đúng câu báo, quay về màn đăng nhập;
   đúng mật khẩu vào game bình thường. Bằng chứng "sau".
4. Cùng phiên chạy thử: nhìn tận mắt vụ nền đen, chốt nguyên nhân, sửa,
   chạy lại xác nhận.
5. Thử lệnh đăng ký cũ (nếu client cũ còn) → nhận câu "đăng ký tại web".

## Build lại và phát hành

Thứ tự: `zz1.jar` → bản PC (`NROKAME.zip`) → APK (**user ký lại** bằng
`KY-APK.bat`) → xuất lại gói iOS (`D:\NRO_iOS_GuiMac.zip` hiện tại ĐÃ CŨ,
đừng chép sang MacBook). Cuối cùng user chạy `DONG-BO.ps1` trên VPS.

Việc user phải tự làm: ký APK (cần mật khẩu khoá) và chạy DONG-BO trên VPS.

## Ngoài phạm vi

- Không thêm ô số điện thoại / khôi phục mật khẩu trong game.
- Không refactor gì ngoài các file nêu trên.
- Trang quenmatkhau.php giữ nguyên cơ chế "liên hệ Admin".
