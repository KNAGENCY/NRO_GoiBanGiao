# Thư mục này là gì

Bốn file cấu hình đã chuẩn bị sẵn để đưa server **Ngọc Rồng Kame** từ máy bạn lên VPS
`103.116.38.60`. Đây mới chỉ là **phần cấu hình** — chưa phải toàn bộ quy trình cài đặt.

Tất cả nội dung viết bằng tiếng Việt. Trong file cấu hình tôi bỏ dấu tiếng Việt
(viết "may chu" thay vì "máy chủ") để Apache và Java đọc không bị lỗi phông chữ.

---

## Từng file dùng làm gì

### 1. `config.properties` — cấu hình server game

**Chép đè lên:** `<thư mục game trên VPS>\data\config\config.properties`

Dựa trên file gốc ở máy bạn, chỉ khác **đúng 2 chỗ**:

| Dòng | Máy nhà | Bản VPS |
|---|---|---|
| `server.sv1` | `...:127.0.0.1:14445:...` | `...:103.116.38.60:14445:...` |
| `database.pass` | để trống | `__MAT_KHAU_MYSQL__` |

**`__MAT_KHAU_MYSQL__` là chỗ trống, không phải mật khẩu.** Script cài đặt sẽ hỏi bạn
rồi tự điền. Nếu tự sửa tay thì thay cả cụm chữ đó bằng mật khẩu thật.

Dòng `server.sv1` là **dòng quan trọng nhất trong cả bộ cài**. Đây là địa chỉ mà server
gửi xuống cho điện thoại người chơi. Để `127.0.0.1` thì máy người chơi tự gọi vào chính
nó — **không ai vào được game, dù server trên VPS chạy hoàn toàn bình thường**. Lỗi này
đặc biệt dễ lọt vì máy bạn test thì vẫn vào được (với máy bạn, `127.0.0.1` đúng là máy bạn).

Hai điều tôi cố tình **không** đổi, để bạn tự quyết:

- `server.debug=true` → **nên đổi thành `false` trước ngày mở cửa.** Chế độ này in rất
  nhiều thông tin nội bộ, vừa lộ cấu trúc hệ thống, vừa làm file log phình to ăn hết
  ổ C: 70 GB. Tôi để nguyên vì yêu cầu là "giữ nguyên mọi thứ khác".
- `database.host=localhost` → giữ nguyên. Chỉ đổi thành `127.0.0.1` **nếu** bạn làm
  theo khuyến nghị bỏ dấu `#` ở dòng `bind-address` trong `C:\xampp\mysql\bin\my.ini`.
  Lý do ghi ngay trong file.

---

### 2. `nrokame-vhost.conf` — cấu hình tên miền cho Apache

**Không chép đè.** Đây là một đoạn văn bản để bạn **chép thêm vào cuối** file:

```
C:\xampp\apache\conf\extra\httpd-vhosts.conf
```

Mở bằng Notepad, kéo xuống cuối, chép phần nằm giữa hai dòng đánh dấu
`=== BAT DAU CHEP ===` và `=== HET CHEP ===`, dán vào, lưu lại.

Không cần sửa `httpd.conf` để nạp file này — dòng `Include conf/extra/httpd-vhosts.conf`
đã có sẵn ở `httpd.conf` dòng 521 (tôi đã kiểm tra trên máy bạn).

**Nhưng có một chỗ khác trong `httpd.conf` phải sửa tay:** dòng 228, đổi
`ServerName localhost:80` thành `ServerName nrokame.com:80`.

Cấu hình này làm cả `nrokame.com` và `www.nrokame.com` cùng mở ra thư mục
`C:\xampp\htdocs`. Vì đây là vhost duy nhất nên nó cũng là vhost mặc định — người ta gõ
thẳng `http://103.116.38.60` cũng ra đúng trang này, điều đó là bình thường.

**Cảnh báo quan trọng ghi trong file:** đừng tự thêm khối `<Directory>` vào vhost.
`httpd.conf` dòng 253-273 đã có sẵn khối `<Directory "C:/xampp/htdocs">` với
`AllowOverride All`, vhost này dùng chung nên tự được hưởng. Nếu bạn thêm một khối
`<Directory>` khác mà quên `AllowOverride All` thì Apache **bỏ qua hoàn toàn file
`.htaccess`**: đường dẫn sạch (`/sv1/trangchu`) hỏng hết **và** lớp chặn spam im lặng
không chạy. Cả hai đều không báo lỗi gì.

**Chưa cấu hình HTTPS ở bước này** — cố ý. Apache của bạn đã nghe sẵn cổng 443, nhưng
chứng chỉ đang dùng là chứng chỉ tự ký mặc định của XAMPP: trình duyệt người chơi sẽ
hiện trang cảnh báo đỏ "Kết nối của bạn không riêng tư" và phần lớn sẽ quay đi. Giai
đoạn đầu **chỉ mở cổng 80**. Phần cuối file hướng dẫn 2 cách lấy khoá xanh sau này
(Cloudflare — dễ nhất, hoặc Let's Encrypt).

---

### 3. `htaccess-vps.txt` — bản `.htaccess` cho VPS

**Chép vào `C:\xampp\htdocs\` rồi ĐỔI TÊN thành `.htaccess`** (có dấu chấm ở đầu,
không có đuôi `.txt`).

Windows hay tự ẩn đuôi file. Nếu đổi tên xong vẫn thấy `.htaccess.txt` thì vào File
Explorer → thẻ View → tick "File name extensions" rồi đổi lại.

**Kết luận sau khi kiểm tra: dòng `auto_prepend_file` không phải sửa gì cả.** Nó trỏ
tới `C:/xampp/htdocs/ddos_guard.php` — là đường dẫn Windows tuyệt đối, và trên VPS
Windows Server nó **vẫn đúng, với điều kiện bạn cài XAMPP vào đúng `C:\xampp`**.

Nội dung file giống hệt bản đang chạy, tôi chỉ thêm phần chú thích giải thích.

Chỗ nguy hiểm ở đây: nếu đường dẫn sai, **Apache không báo lỗi gì cả** — web vẫn chạy
bình thường, chỉ là lớp bảo vệ biến mất mà bạn không hề biết. Trong file có hướng dẫn
cách thử để biết chắc lớp bảo vệ đang chạy hay không.

Nếu bạn cài XAMPP sang ổ khác (ví dụ `D:\xampp`) thì phải sửa **4 chỗ cùng lúc** —
danh sách ghi ngay đầu file. **Lời khuyên: cứ cài đúng `C:\xampp` cho khỏi phải đụng.**

---

### 4. `DANH-SACH-CHEP-LEN-VPS.md` — chép những gì, nặng bao nhiêu

Danh sách đầy đủ, dung lượng **đo thật** chứ không ước lượng. **Tổng khoảng 862 MB.**

Tóm tắt: `data\` 822,1 MB · `lib\` 11,1 MB · `zz1.jar` 7,58 MB · `src\` 4,9 MB ·
web 15,4 MB · dump CSDL ~1 MB.

Bỏ qua được ~17,4 MB: `out\`, `dist\`, `nbproject\`, `build\`, `logs\`, `zz4.rar`,
`run.bat.goc`, `sources.txt`, và file `.bak-update-shop` trong `src\`. Tiết kiệm dung
lượng không đáng kể — **lý do bỏ là để sau này không ai mở nhầm file cũ ra sửa.**

File này cũng có sẵn lệnh PowerShell nén và cách kiểm tra chép có thiếu file không.

> **Hai điều dễ bất ngờ, nói trước ở đây:**
>
> **Một — phải chép cả thư mục `src\` (mã nguồn Java) lên VPS.** Nghe vô lý vì server
> chạy bằng file `.jar`, nhưng trang quản trị đọc thẳng mã nguồn bằng regex
> (`game_bridge.php` dòng 194-236) để hiện tỉ lệ nâng sao Pha Lê. Không có `src\` thì
> trang đó hiện toàn số 0.
>
> **Hai — KHÔNG được dùng `zz1\z.sql` để tạo cơ sở dữ liệu trên VPS.** Tôi đã kiểm
> chứng: `z.sql` có **47 bảng**, cơ sở dữ liệu đang chạy có **48 bảng**. `z.sql` là bản
> dump từ 20/07, cũ hơn một ngày so với các thay đổi lớn hôm 21/07 — nó **thiếu hẳn
> bảng `login_security`** (phục vụ việc mã hoá mật khẩu người chơi) và không hề chứa
> chữ "nrokame" (nghĩa là các link tải game trong bảng `adminpanel` vẫn là bản cũ).
> **Dùng `z.sql` = quay ngược toàn bộ công sức đổi thương hiệu và mã hoá mật khẩu, mà
> không có thông báo lỗi nào.** Phải xuất bản dump mới — lệnh có sẵn trong file.

---

## Ba việc quan trọng KHÔNG nằm trong thư mục này

Bốn file trên chỉ lo phần cấu hình. Còn ba việc nữa, thiếu bất kỳ việc nào thì
server vẫn không chạy được hoặc không ai vào được game.

### a) Mật khẩu MySQL phải đổi ở đúng 3 nơi, cùng một lần

Hiện mật khẩu root đang **rỗng**. Lên VPS bắt buộc phải đặt. Ba nơi:

1. `<thư mục game>\data\config\config.properties` — dòng `database.pass`
2. `C:\xampp\htdocs\sv1\config.php` — dòng 5
3. `C:\xampp\htdocs\sv1\boss_status.php` — dòng 35

Thiếu chỗ nào thì chỗ đó chết, và **triệu chứng khác hẳn nhau** nên rất dễ chẩn đoán nhầm:

- Sót **(1)**: server game không nối được CSDL → cửa sổ báo lỗi rồi tự thoát → `run.bat`
  chạy lại vô tận, cửa sổ hiện lên tắt quá nhanh để đọc lỗi.
- Sót **(2)**: toàn bộ web chết ngay, mọi trang hiện chữ đỏ "Lỗi kết nối Database".
- Sót **(3)**: web vẫn vào bình thường, **chỉ mất ảnh mặt boss** trên trang chủ. Lỗi bị
  nuốt im lặng trong khối `catch` nên không có thông báo gì. Đây là chỗ dễ sót nhất vì
  mật khẩu nằm lọt thỏm giữa dòng lệnh chứ không phải một biến riêng.

Nhớ: **đặt mật khẩu MySQL sau khi đã nhập xong dữ liệu**, làm ngược lại thì lệnh nhập
bị từ chối.

**Tôi không viết mật khẩu thật vào bất kỳ file nào trong thư mục này** — và bạn cũng
đừng viết vào tài liệu, ghi chú hay tin nhắn. Chỉ gõ thẳng vào 3 file trên, ngay trên VPS.

### b) Tường lửa: mở 80 và 14445, tuyệt đối không mở 3306

Chạy PowerShell **bằng quyền Administrator** (chuột phải → "Run as administrator"):

```
netsh advfirewall set allprofiles state on
netsh advfirewall firewall add rule name="Kame - Web HTTP 80" dir=in action=allow protocol=TCP localport=80
netsh advfirewall firewall add rule name="Kame - Game 14445" dir=in action=allow protocol=TCP localport=14445
netsh advfirewall firewall add rule name="Kame - Chan MySQL 3306" dir=in action=block protocol=TCP localport=3306
```

Cổng 14445 chỉ dùng TCP, không cần luật UDP.

**Vì sao 3306 nguy hiểm đến thế:** mật khẩu root đang rỗng, **và** MariaDB đang nghe
trên mọi card mạng (dòng `bind-address` trong `C:\xampp\mysql\bin\my.ini` đang bị vô
hiệu bởi dấu `#`). Ghép hai điều đó lại: mở cổng 3306 là bất kỳ ai trên thế giới chỉ
cần gõ `103.116.38.60:3306` với tài khoản `root` không mật khẩu là vào thẳng cơ sở dữ
liệu — tự cộng ngọc, đọc trộm toàn bộ tài khoản người chơi, hoặc xoá sạch 48 bảng
trong vài giây.

Web và server game **không cần** cổng này mở ra ngoài, vì cả hai đều chạy trên chính
con VPS đó và nối qua `localhost`.

**Nên làm thêm:** bỏ dấu `#` ở dòng `bind-address="127.0.0.1"` trong
`C:\xampp\mysql\bin\my.ini` (khoảng dòng 46), rồi khởi động lại MySQL. Như vậy dù
tường lửa có bị tắt, MariaDB vẫn tự từ chối mọi kết nối từ bên ngoài.

### c) Server game phải tự bật lại sau khi VPS khởi động lại

Đây là hạn chế thật của phần mềm, không phải lựa chọn: `main()` ở
`ServerManager.java:167` chỉ mở cửa sổ Swing, toàn bộ khởi động nằm trong constructor
`ServerManagerUI`. **Server bắt buộc phải có màn hình đồ hoạ.**

Hệ quả: **không chạy được dạng Windows Service** (service chạy ở session 0, không có
màn hình). Trên Windows Server, sau khi máy khởi động lại mà không ai đăng nhập thì
server **không tự lên**.

Cách duy nhất chạy được: bật **tự động đăng nhập** (AutoAdminLogon) + đặt lối tắt
`run.bat` vào thư mục Startup.

Ba điều bắt buộc khi tạo lối tắt đó:

1. Trỏ **thẳng vào `run.bat`**, đừng gọi `java.exe` trực tiếp. `run.bat` có dòng
   `cd /d "%~dp0"` để đặt đúng thư mục làm việc — mà **mọi đường dẫn trong Java đều là
   đường dẫn tương đối**. Sai thư mục làm việc thì server không đọc được
   `config.properties` (rơi về mật khẩu rỗng), không nạp được bản đồ. Triệu chứng lúc
   đó **rất giống lỗi cơ sở dữ liệu** nên cực kỳ dễ chẩn đoán nhầm.
2. Tick **"Run as administrator"**. Chức năng tự chặn IP phá hoại gọi lệnh `netsh` của
   Windows. Không đủ quyền thì lệnh này **thất bại lặng lẽ** — server vẫn chạy, log vẫn
   in `[BAN IP VINH VIEN WINDOWS]`, nhưng tường lửa không chặn ai cả. Bạn sẽ tưởng mình
   có lớp chống phá hoại trong khi thực tế không có.
3. **Dùng lại `run.bat` có sẵn, đừng viết lại.** Nó đã có vòng canh gác đúng: mã thoát
   2 = bảo trì, chạy lại sau 3 giây; mã 0 = dừng hẳn; mã khác = lỗi, tự chạy lại sau 10
   giây, ghi vào `logs\restart.log`. Chỉ cần kiểm tra **dòng 20** (`WEB_ROOT`) trỏ đúng
   chỗ cài XAMPP. Dòng 19 (`GAME_ROOT`) dùng `%CD%` nên tự đúng, không phải đụng.

---

## Còn client Unity thì sao

Không nằm trong thư mục này, nhưng **thiếu là hỏng cả buổi mở server**. Trong
`D:\NRO_Project` có **4 chỗ** chứa `127.0.0.1`, phải đổi hết sang `103.116.38.60`:

- `Assets\Scripts\Game1\ServerListScreen.cs` — `ListIP` và `linkDefault`
- `Assets\Scripts\Game2\ServerListScreen.cs` — `ListIP` và `linkDefault`

**Tìm theo chữ `127.0.0.1`, đừng sửa mù theo số dòng** — file Game2 lệch một dòng so
với Game1. Giữ nguyên phần đuôi `:14445:0,0,0`.

Sửa xong **bắt buộc phải build lại** (khoảng 2 phút), vì đây là hằng số biên dịch cứng
vào file thành phẩm. Kiểm chứng bằng cách quét chuỗi trong file
`<bản build>_Data\il2cpp_data\Metadata\global-metadata.dat`: phải tìm thấy
`103.116.38.60` và **không** còn `127.0.0.1`.

---

## Thứ tự làm cho đỡ rối

1. Xuất bản dump CSDL mới ở máy nhà (**đừng dùng `z.sql`**).
2. Nén `zz1` và `htdocs` theo hướng dẫn, đưa lên VPS, giải nén, **đếm file kiểm tra**.
3. Cài XAMPP vào đúng `C:\xampp`, cài JDK 21 vào `C:\Program Files\Java\jdk-21`.
4. Tạo cơ sở dữ liệu `z`, nhập bản dump mới, kiểm tra ra đúng **48 bảng**.
5. Đặt mật khẩu MySQL, rồi sửa **cả 3 chỗ** cùng một lần.
6. Chép `config.properties`, `.htaccess`, thêm đoạn vhost, sửa `ServerName` dòng 228.
7. Mở tường lửa 80 + 14445, chặn 3306, bật `bind-address` trong `my.ini`.
8. Bật Apache + MySQL, mở `http://nrokame.com` **từ điện thoại đã tắt wifi (dùng 4G)** —
   phải ra trang web.
9. Chạy `run.bat` bằng quyền Administrator. Kiểm tra server đã nghe cổng chưa:
   `netstat -ano | findstr ":14445"` — phải thấy chữ `LISTENING`.
10. Build lại client với IP mới, test vào game thật.
11. Xong xuôi mới bật tự động đăng nhập + lối tắt Startup, rồi **khởi động lại VPS một
    lần để thử** xem server có tự lên không.

---

## Còn hai việc chưa ai chốt

- **Ngày mở server.** `sv1\config.php` dòng 232 vẫn là `2026-07-01` — đã là quá khứ.
  Đồng hồ đếm ngược trên trang chủ đang hiện số âm hoặc "đã kết thúc". Cần chốt ngày thật.
- **Link Zalo và Facebook.** `sv1\config.php` dòng 233-235 đang trỏ tạm về trang liên hệ
  của chính website. Người chơi bấm vào sẽ quay vòng lại web chứ không tới nhóm cộng
  đồng — mất kênh chăm sóc người chơi ngay ngày mở server.
