# CLAUDE.md

Hướng dẫn cho Claude Code khi làm việc trong repo này.

## Đây là gì

Server private **Ngọc Rồng Online (NRO)** / "Dragon Boy" — mua lại từ người khác, đang chuẩn bị launch. Chủ sở hữu **không phải lập trình viên**: hãy giải thích bằng tiếng Việt đời thường, tránh thuật ngữ không cần thiết, và luôn nói rõ *tại sao* trước khi nói *làm thế nào*.

Đây **không phải một project** mà là **bốn mảnh ghép** nằm ở hai gốc thư mục khác nhau, dính chặt vào nhau:

| Mảnh | Đường dẫn | Công nghệ |
|---|---|---|
| Server game | `C:\Users\Admin\Music\zz1` | Java 21, project NetBeans/Ant tên `NROOCEAN` |
| Web | `C:\xampp\htdocs` | PHP 8 + XAMPP (Apache) |
| Client (bản build) | `C:\Users\Admin\Music\LocalNew` | Unity 2022.3.62f2, IL2CPP |
| Client (source) | `1_Project_HanhTrangMoi.rar` (408MB, chưa giải nén) | Unity project đầy đủ — 431 file C#, cùng phiên bản 2022.3.62f2 |
| Database | MariaDB 10.4, schema `z` | Dump gốc: `zz1\z.sql` |

`C:\Users\Admin\Music\htdocs` là **bản sao y hệt** của `C:\xampp\htdocs`. Bản đang thực sự chạy là bản trong XAMPP — sửa ở đó. Nếu sửa cả hai, phải đồng bộ.

## Lệnh thường dùng

```bat
REM Chạy server (từ C:\Users\Admin\Music\zz1)
run.bat

REM Import database lần đầu
C:\xampp\mysql\bin\mysql.exe -u root -e "CREATE DATABASE z CHARACTER SET utf8mb4;"
C:\xampp\mysql\bin\mysql.exe -u root z < z.sql

REM Build lại từ source — KHÔNG cần Ant, dùng JDK 21 có sẵn. Đã kiểm chứng chạy được.
REM Xem chi tiết 3 điểm bắt buộc ở mục "Build lại server" bên dưới.
```

Apache và MySQL bật/tắt qua **XAMPP Control Panel**.

## Kiến trúc: ba mảnh nói chuyện với nhau thế nào

Không có socket, không có HTTP API giữa web và server. Chỉ ba kênh:

1. **Database `z` dùng chung** — cả Java lẫn PHP đọc/ghi cùng bảng. Đây là lý do admin panel sửa được game đang chạy.
2. **`panel_cmd.txt`** (web ghi → Java đọc, poll mỗi 2 giây rồi **xoá file**). Lệnh phân cách bằng `|`: `BUFF_VND`, `BUFF_ITEM`, `KICK_PLAYER`, `ADMIN_MODE`, `XNXX_MODE`, `MAINTENANCE_NOW`, `SET_EXP`, `SET_DROP`, `SPAWN_BOSS`, `EDIT_BOSS`, `BOSS_PERSIST_MODE`, `LS`.
3. **File trạng thái JSON** (Java ghi → web đọc, mỗi 1 giây): `boss_status.json` ghi vào `C:\xampp\htdocs\sv1\`, `panel_server_status.json` ghi vào game root.

Phân chia dữ liệu: **logic và chỉ số game nằm trong MySQL**; thư mục `zz1\data\` là **kho asset client** mà server stream qua socket theo 4 mức zoom `x1`–`x4`.

## Bẫy — đọc trước khi sửa bất cứ thứ gì

- **Server bắt buộc có giao diện đồ hoạ.** `main()` (`zz1\src\server\ServerManager.java:167`) chỉ mở cửa sổ Swing; toàn bộ khởi động nằm trong constructor `ServerManagerUI`. **Không chạy headless được** → VPS Linux thường sẽ không chạy nổi nếu không có X ảo. Chủ cũ dự định deploy lên **VPS Windows**.
- **Phải dùng Java 21.** `zz1.jar` biên dịch bằng Java 21 (bytecode major 65), dù file cấu hình ghi `javac.target=17`. Lệnh `java` trên PATH hiện là **JRE 8** → chạy `run.bat` bây giờ sẽ crash.
- **Có hai file jar, đừng nhầm.** `zz1.jar` ở gốc là bản đang chạy (fat jar, tự chứa MySQL driver + Hikari + ga.jar). `dist\NROKRAI.jar` là bản build Ant cũ 2 tháng. `run.bat` dùng `zz1.jar`.
- **Build lại được — đã kiểm chứng 2026-07-21.** Nhưng có 3 cái bẫy, xem mục "Build lại server" bên dưới.
- **Tầng mạng không có source.** `network.Network`, `network.Session`, `network.SessionManager`, `utils.Logger` nằm trong `lib\ga.jar` đã biên dịch. Thư mục `src\network\` là **bản decompile cũ, code chết** — sửa vào đó không có tác dụng gì.
- **Vòng lặp khởi động lại bị kẹt.** `ServerManager.close()` gọi `run.bat` *trước* `System.exit(0)`, mà `run.bat` kết thúc bằng `PAUSE` → mỗi lần bảo trì đẻ thêm một console chờ bấm phím.
- **Đổi tên thư mục `sv1` là đứt toàn bộ liên kết** — tên thư mục web là khoá tra cứu game root (`game_bridge.php`).
- **Admin panel đọc source Java bằng regex** (`game_bridge.php:194-236` parse `CombineService.java`) → phải deploy cả `src\` cạnh jar, và sẽ vỡ nếu format lại code đó.
- **Working directory phải là game root** — mọi đường dẫn trong Java đều tương đối. `run.bat` xử lý bằng `cd /d "%~dp0"`.
- **Game root phải ghi được** — `data\update_data\part` bị ghi đè lại mỗi lần khởi động từ bảng `part`.

## Build lại server

**Đã kiểm chứng chạy được** (2026-07-21). Không cần Ant — dùng `javac` của JDK 21. Ba điểm bắt buộc, thiếu cái nào cũng fail:

1. **`copy lib\ga.jar lib\LIB_SALE.jar`** — đã làm rồi, giờ `lib\` có cả hai file (giống hệt nhau).
2. **Đừng dùng `sources.txt`** — nó là danh sách cũ, **thiếu 6 file** (`PhucLoiService`, `TamBaoService`, `AmbientBot`, `NangCapBongTai3/4`, `NangChiSoBongTai3`). Hãy liệt kê toàn bộ 488 file `.java` trong `src\`. Ghi danh sách **không có BOM**, nếu không `javac` báo `Invalid filename`.
3. **Chép `src\server\icon.png` vào thư mục class** trước khi đóng gói. Thiếu là server crash ngay: `NullPointerException ... ImageIcon ... ServerManagerUI.java:53`.

```powershell
cd C:\Users\Admin\Music\zz1
$files = Get-ChildItem src -Recurse -Filter *.java | ForEach-Object { $_.FullName }
[IO.File]::WriteAllLines($list, $files)          # không BOM
& "C:\Program Files\Java\jdk-21\bin\javac.exe" -nowarn -encoding UTF-8 --release 17 -cp "lib\LIB_SALE.jar" -d $out "@$list"
Copy-Item src\server\icon.png $out\server\
# gộp toàn bộ ga.jar vào $out (bỏ META-INF\*.SF,*.DSA,*.RSA), rồi jar cfm với "Main-Class: server.ServerManager"
```

`ga.jar` đã chứa sẵn mọi thư viện: MySQL driver, HikariCP, org/json, Apache commons, Lombok. Không cần tải thêm gì.

**Test bản build mới mà không đụng server đang chạy:** tạo DB riêng từ `z.sql`, dựng game root riêng dùng `mklink /J` trỏ vào các thư mục asset nặng (`icon` 545MB, `img_by_name` 137MB...), chỉ copy thật `config\` và `update_data\`, rồi đổi `server.port` và `database.name` trong bản config đó.

## Bảo mật — tình trạng hiện tại

Những điều dưới đây là **sự thật của code hiện tại**, không phải đề xuất. Đừng vô tình làm nặng thêm.

- **Mật khẩu người chơi lưu plaintext**, cả Java lẫn PHP đều so sánh chuỗi trực tiếp trên cột `account.password`. Muốn chuyển sang hash thì **phải sửa cả hai phía cùng lúc**, không thì toàn bộ tài khoản mất đăng nhập.
- **DB dùng `root` mật khẩu rỗng**. Credential lặp ở hai chỗ: `sv1\config.php:2-5` và `sv1\boss_status.php:35`.
- **`index.php` là trang ngụy trang** — site học tập giả "EduVerse". Trang game thật ở `/sv1/trangchu`.
- `ddos_guard.php` nạp qua `auto_prepend_file` với đường dẫn Windows tuyệt đối (`.htaccess:6`) → **trên Linux im lặng không chạy**, mất sạch lớp bảo vệ mà không báo lỗi.
- Số người online hiển thị công khai **bị cộng thêm 45** (`game_bridge.php:361-370`).

## Thương hiệu: Ngọc Rồng Kame

Tên chính thức **Ngọc Rồng Kame**, tên miền **nrokame.com** (user tự đăng ký). Đã thay xong toàn bộ 2026-07-21.

**Dấu vết chủ cũ đã xoá sạch** — 187 chỗ, gồm: bot chat quảng cáo SĐT chủ cũ (`BotManager.java`), cửa hậu miễn lọc chat cho nick `GiaCatLuong` (`Functions.java:23`), watermark `srcnrofree.online` hiện trong game, 3 số điện thoại, email, credential VPS `14.225.203.122`, và 2 thư mục build cũ.

**Cố ý KHÔNG đổi:** tên package/class nội bộ `GiaCatLuongManager`, `GiaCatLuongServer` (545 chỗ). Người chơi không thấy; đổi là refactor toàn bộ, rủi ro cao, vô ích.

**Các link đang để tạm, chờ thông tin thật của user:**

| Chỗ | Giá trị tạm |
|---|---|
| `sv1\config.php` — Zalo, Facebook, TestFlight | `nrokame.com/lien-he`, `nrokame.com/tai-game` |
| `sv1\admin\panel.php` — 2 nút hỗ trợ | `nrokame.com/lien-he` |
| `privacy-policy.php` — email | `support@nrokame.com` (**chưa tạo hòm thư**) |
| Bảng `adminpanel` — 7 cột link tải | `nrokame.com/tai-game` |
| Client `Panel.cs:14945` — nạp thẻ | `nrokame.com/nap-the?username=` |

`$SITE_OPEN_TIME` trong `sv1\config.php` vẫn là `2026-07-01` (quá khứ) — cần chốt ngày mở thật.

## Client — build lại và đổi địa chỉ server

**Đừng vá `global-metadata.dat` nữa.** Đó là cách chữa cháy khi chưa có source. Giờ có source đầy đủ ở `D:\NRO_Project` nên cứ sửa mã rồi build lại — nhanh hơn (1 phút) và không bị giới hạn 32 byte của chuỗi cũ.

**Mọi thứ về thương hiệu nằm trong `ServerListScreen.cs`**, và có **hai bản song song** `Game1` + `Game2` (mỗi bản 213 file, cùng được build vào 2 scene `Nro1.unity` / `NRO2.unity`). **Sửa một bản mà quên bản kia là lỗi âm thầm.**

| Hằng số | Giá trị hiện tại | Ảnh hưởng |
|---|---|---|
| `ListIP`, `linkDefault` | `Ngoc Rong Kame:nrokame.com:14445:0,0,0` | địa chỉ máy chủ |
| `ModProductName` | `Ngoc Rong Kame` | tên file .exe, thư mục lưu |
| `ModApplicationIdentifier` | `com.nrokame.game` | **vĩnh viễn** sau khi lên chợ ứng dụng |
| `ModFunc.homeUrl` | `nrokame.com` | chữ hiện góc màn hình đăng nhập |

`ModBuildIdentitySync.cs` tự đọc 3 hằng đầu và áp vào PlayerSettings mỗi lần build — không sửa tay trong Unity.

**Dùng tên miền chứ không IP trần**: tầng kết nối là `TcpClient.Connect(host, port)`, tự phân giải tên miền. Đổi VPS chỉ cần sửa DNS, client đã cài trên máy người chơi vẫn chạy.

**Lệnh build** (đã kiểm chứng — PC ~1 phút, APK ~4 phút):
```
D:\Unity\2022.3.62f2\Editor\Unity.exe -batchmode -quit -nographics
  -projectPath D:\NRO_Project -buildTarget Android
  -executeMethod KameBuild.BuildAndroid -logFile <log>
```
Đổi `Android` → `Win64` + `BuildWindows` để ra bản PC. **Luôn truyền `-buildTarget` khớp với hàm gọi** — dự án dùng chung cho 3 nền tảng, để lệch thì Unity phải đổi nền tảng giữa chừng (chậm, và với iOS thì hỏng ngầm, xem mục iOS bên dưới). `KameBuild.cs` (`Assets\Editor\`) cố định thư mục ra `D:\NRO_Build`, trả mã thoát khác 0 khi hỏng.

**Bẫy khi đóng gói**: thư mục `*_BackUpThisFolder_ButDontShipItWithYourGame` nặng 233 MB, Unity ghi rõ đừng gửi kèm. Bỏ ra thì còn 274 MB → nén zip còn 69 MB.

**Unity Hub không quản lý được bản 2022.3.62f2** vì nó cài tay (thiếu `modules.json`). Muốn thêm module phải tải gói riêng từ API của Unity:
`https://services.api.unity.com/unity/editor/release/v1/releases?version=2022.3.62f2`
API trả về đủ địa chỉ các gói con (OpenJDK 11, SDK, NDK r23b, build-tools…). Module Android nằm ở `Editor\Data\PlaybackEngines\AndroidPlayer\` (5,4 GB, đã cài đủ).

**Bộ Android thiếu 3 file giấy phép thì Gradle treo im lặng** giữa chừng chờ bấm phím mà không hiện câu hỏi. Đã tạo sẵn trong `AndroidPlayer\SDK\licenses\`: `android-sdk-license`, `android-sdk-preview-license`, `android-sdk-arm-dbt-license`.

**Bản APK 2026-07-22** — 71,3 MB, đã mở ra kiểm từng thứ bên trong:
`com.nrokame.game` · nhãn `Ngoc Rong Kame` · Android 5.1+ (minSdk 22, target 34) · cả `arm64-v8a` lẫn `armeabi-v7a` · chuỗi máy chủ nhúng đúng `Ngoc Rong Kame:nrokame.com:14445:0,0,0` · **không còn** `Dragon Boy`, `srcnrofree`, `nrocat`, `127.0.0.1`.

**Khoá ký APK — ĐÃ TẠO 2026-07-22.** `D:\NRO_Khoa\nrokame.keystore` (2.715 byte, bí danh `nrokame`, RSA 2048, hạn 100 năm). APK đã ký lại: `CN=Ngoc Rong Kame, OU=Game, O=Kame, L=Ha Noi, ST=Ha Noi, C=VN`, đủ v1+v2+v3, **không còn** `CN=Android Debug`. Ba bản `D:\NRO_Build\Android\NgocRongKame.apk`, `sv1\down\NRO.apk`, `goi-len-vps\NRO.apk` cùng mã băm MD5 `5E82B731F511C967…`.

Mỗi lần build APK mới **phải chạy lại** `bo-cai-vps\KY-APK.bat` (bấm đúp) hoặc `KY-VA-DANG-APK.ps1` — Unity luôn ký bằng khoá gỡ lỗi. Script tự phát hiện đã có khoá thì chỉ hỏi mật khẩu cũ, không đè khoá mới. Nếu chữ ký ra sai nó dừng và không chép file đi đâu cả.

⚠️ **Mất `nrokame.keystore` hoặc quên mật khẩu = mất vĩnh viễn**, mọi người chơi phải gỡ game cài lại. Không ai khôi phục được.

**Lỗi `ArgumentNullException: shader` — ĐÃ SỬA 2026-07-22.** `mGraphics.CreateLineMaterial()` gọi `Shader.Find("Lines/Colored Blended")`, mỗi lần mở game bản build đều ném lỗi.

Nguyên nhân **không phải** "Unity cắt shader khỏi build" như ghi chú cũ. Shader đó là đồ thời Unity 4, **Unity 2022.3 đã gỡ hẳn** — tìm chuỗi tên trong cả `unity default resources` lẫn `unity_builtin_extra` đều không còn (đối chứng: `Hidden/Internal-Colored`, `Sprites/Default`, `UI/Default` thì tìm thấy). Nên `Shader.Find` trả `null` **ngay trong Editor**, và không thể "thêm vào Always Included Shaders" vì không có gì để thêm.

Cách chữa gồm hai mảnh, thiếu một là lỗi quay lại:
- `Assets\Shaders\LinesColoredBlended.shader` — tự cấp lại shader mang **đúng tên** `Lines/Colored Blended` (viết bằng vertex/fragment, vì bản gốc kiểu fixed-function `BindChannels` Unity mới không dịch được).
- `Assets\Editor\KameAlwaysIncludedShaders.cs` — tự đưa shader đó vào Graphics Settings > Always Included Shaders qua `SerializedObject` (không sửa tay YAML). Chạy tự động trong `KameBuild.BuildWindows/BuildAndroid`, hoặc tay qua menu `Mod > Them shader bat buoc vao Always Included`. Chạy lại nhiều lần không sao. Không có bước này Unity lại cắt shader vì không scene/prefab nào trỏ tới.

## Bản clone tối ưu (sandbox) — 2026-07-22

User muốn tối ưu code/DB cũ nhưng KHÔNG đụng bản chạy thật. Đã dựng sandbox tách hẳn:

| Thành phần | Bản thật (đồng bộ VPS) | Bản clone (thử nghiệm) |
|---|---|---|
| Mã nguồn server | `C:\Users\Admin\Music\zz1\src` | `D:\NRO_ToiUu\src` |
| CSDL | `z` | `z_thu` (bản sao 48 bảng) |
| Cổng | 14445 | 14446 |
| Game root | game root VPS | `C:\NRO_Thu` (asset dùng chung qua junction) |

Sửa tối ưu **chỉ ở `D:\NRO_ToiUu\src`**. Build + chạy: `D:\NRO_ToiUu\BUILD-VA-CHAY.ps1` (bằng chứng: build sạch 489 file 2026-07-22). Client thử trỏ `127.0.0.1:14446`. Xong, kiểm chứng OK thì mới port thay đổi sang bản thật rồi đồng bộ VPS. **Nguyên tắc: đo trước, sửa sau — không đập lại engine.**

Giftcode: đồ vật vào **hòm thư** (`itemsMailBox`), mở qua NPC **Ông Gohan** → "Hòm Thư (N món)" (`OngGohan.java:64`). Vàng/ngọc vào thẳng nhân vật. Muốn giftcode cho đồ thẳng vào rương thì sửa `GiftCodeService.giftCode()` nhánh `default` (`:127`) — ĐÃ làm trong clone `D:\NRO_ToiUu` (ưu tiên `addItemBag`, đầy rương mới rớt hòm thư), chưa port sang bản thật.

## Lỗi panel admin đã sửa 2026-07-22 (chiều muộn)

- **"Không thể tải danh sách vật phẩm"** (`buff_item.php` endpoint `?ajax_items=1`): JSON bị hỏng do output bẩn/UTF-8 lỗi. Vá 3 lớp: `ob_start()` ngay dòng đầu cho nhánh ajax (trước khi nạp `game_bridge`), tắt `display_errors`, thêm cờ `JSON_INVALID_UTF8_SUBSTITUTE`; và trả `error` trong JSON + client hiện lỗi thật thay vì câu chung. Chỉ `buff_item.php` dính pattern này.
- **Ảnh icon vật phẩm không hiện trong panel**: web gọi `/icon/<icon_id>.png` nhưng `htdocs\icon\` chưa từng được tạo. Kho icon thật ở `zz1\data\icon\x3\<id>.png` (đặt tên đúng theo icon_id; x2/x3/x4 là 3 mức zoom, ~21k file mỗi mức). Sửa bằng **junction** `htdocs\icon` → `data\icon\x3` (không copy). Thư mục `icon` bị DONG-BO loại trừ nên phải tạo junction RIÊNG trên VPS: chạy `bo-cai-vps\TAO-LINK-ICON-VPS.ps1`. Đã tạo junction ở máy nhà, test HTTP 200.
- **Đồng hồ đo panel** (`panel.php` đầu file, `register_shutdown_function`): hiện "Server tạo trang: X ms" góc trái dưới, chỉ GET, chỉ admin — **tạm, gỡ sau khi đo xong**.

## Bản iOS — phần máy Windows đã xong 2026-07-22

User **có MacBook** (hướng dẫn cũ giả định phải đi mượn máy — đã sửa). Chia việc: Windows xuất dự án Xcode, Mac đóng gói và gửi Apple. Không cần cài Unity trên Mac.

Module `iOSSupport` đã cài (2,7 GB, cùng cách với Android: lấy link từ API phát hành của Unity, đối chiếu MD5 Unity công bố + kiểm chữ ký Authenticode trước khi chạy).

```
Unity.exe -batchmode -quit -nographics -projectPath D:\NRO_Project
  -buildTarget iOS -executeMethod KameBuild.BuildIOS -logFile <log>
```

**`-buildTarget iOS` là bắt buộc**, không phải tuỳ chọn. `KameIOSPostProcess.cs` bọc trong `#if UNITY_IOS` — cờ này chỉ bật khi Unity khởi động sẵn ở nền tảng iOS. Thiếu tham số thì script không được biên dịch, `Info.plist` thiếu khai báo mã hoá mà **không báo lỗi gì**. (Rào `#if` là cần thiết: thư viện `UnityEditor.iOS.Xcode` chỉ tồn tại khi đã cài module, bỏ rào ra thì máy chưa cài module hỏng cả build PC lẫn Android.)

`BuildIOS` tự làm 3 việc chống bẫy TestFlight: bật `appleEnableAutomaticSigning` (Mac chỉ cần chọn Team), **tự tăng `buildNumber`** mỗi lần xuất (gửi trùng số là Apple từ chối ngay), và ghi `ITSAppUsesNonExemptEncryption=false` vào `Info.plist` (không có thì mỗi lần gửi đều bị hỏi về mã hoá trước khi phát hành được).

Kết quả: `D:\NRO_Build\iOS` — 2.706 file, 668 MB. Gói gửi Mac: `D:\NRO_iOS_GuiMac.zip` (205 MB), bên trong có `CLAUDE.md` bàn giao cho agent bên Mac (bản gốc: `bo-cai-vps\CLAUDE-cho-may-Mac.md`).

**Nén bằng `tar.exe`, KHÔNG dùng `ZipFile.CreateFromDirectory`.** Thư viện .NET trên Windows ghi dấu phân cách là `\` — sai chuẩn ZIP. Đã dính lỗi này một lần: cả 2.706 mục đều dùng `\`, giải nén trên macOS sẽ ra một đống file phẳng có `\` trong tên, dự án Xcode hỏng sạch mà không báo gì.

Sang Mac phải `chmod +x process_symbols.sh` — nén trên Windows làm mất quyền chạy, thiếu là Xcode báo `Permission denied` giữa chừng.

**`lineMaterial` là code chết** — chỉ được gán rồi không dùng ở đâu. `drawLine()` vẽ bằng `fillRect` và texture 1×1 xoay góc (IMGUI), không đụng tới nó. Nên lỗi này **chưa từng làm mất tính năng nào của người chơi**, chỉ bẩn log. Cứ để nguyên cho an toàn; muốn dọn thì xoá `lineMaterial` + `CreateLineMaterial()` ở **cả Game1 lẫn Game2** và lời gọi ở `Main.cs:254`.

## Những thứ chưa có

- **Bản iOS** — phần Windows đã xong (xem mục dưới). Còn thiếu **tài khoản Apple Developer 99 USD/năm**, user chưa mua. Rủi ro thật: Apple có thể từ chối vì bản quyền Dragon Ball.
- **Trang tải game tự biết file nào đã có**: `tg_download_info()` kiểm `is_file` trong `sv1\down\`, thiếu file thì hiện "Sắp có" và **không đẻ ra nút 404**. Nút dùng đường dẫn rút gọn `/sv1/download?file=...` chứ không phải `download.php?file=...`. Bản PC tên **`NROKAME.zip`** (đổi từ `PC.zip` 2026-07-22, tên cũ trả 404 — sửa ở `download.php`, `tai-game.php`, `down\.htaccess`).

## Đợt sửa 2026-07-22 (chiều) — 4 lỗi bản PC + đăng ký chỉ còn ở web

Spec: `docs\superpowers\specs\2026-07-22-sua-4-loi-ban-pc-design.md`. Đã làm:

- **Đăng nhập sai mật khẩu không còn treo.** `NDVSqlFetcher.login()` trước đây `return null` không gửi gì khi sai mật khẩu → client đợi mãi "Xin chờ" rồi văng. Giờ mọi đường ra (kể cả khối catch) đều `sendThongBaoOK` + `sendLoginFail`. Câu báo kiêm chỉ đường: "…Chưa có tài khoản? Đăng ký tại nrokame.com".
- **Đăng ký CHỈ còn ở web.** `Service.registerAccount()` bị khoá đầu hàm (client sửa đổi gọi thẳng lệnh 42 cũng chỉ nhận câu chỉ ra web); phần thân giữ lại đã sửa lưu mật khẩu qua `auth.PasswordUtil.hash` (trước đây lưu THÔ — lỗi nghiêm trọng). Client: xoá `AccountRegisterScreen.cs`, xoá `requestAccountRegistration` trong `Service.cs`, cả Game1 lẫn Game2.
- **Màn chọn server: 4 nút cố định** Đăng nhập/Đổi tài khoản/Máy chủ/Tùy chọn — bỏ hẳn cơ chế `nCmdPlay` đếm nút động và mọi `cmd[N + nCmdPlay]`. Kèm sửa phím điều hướng (trước wrap tới chỉ số 4 → giờ `cmd.Length - 1`, không sửa là văng game khi bấm phím xuống). `paint()` gọi thẳng `GameCanvas.paintBGGameScr(g)` thay vì tô đen khi `loadScreen=false`.
- **Icon app**: `Assets\Resources\appicon.png` (nền trong suốt user tự tách) thay `appicon.jpg` — GIỮ NGUYÊN GUID `b611de4c…` bằng cách chép nội dung `.meta` cũ sang `appicon.png.meta` (đổi `alphaIsTransparency: 1`). ProjectSettings không phải sửa.
- Môi trường thử ở `C:\NRO_Thu` (DB `z_thu`, cổng 14446, junction asset) — giữ lại để thử lần sau.
- Bản lưu trước sửa: `banluu-server-20260722\`, `banluu-client-20260722\`, `banluu-web-20260722\`.

Bẫy mới ghi nhận: (1) javac một-file lẻ báo "cannot find symbol getSession()" giả — Lombok chỉ sinh hàm khi build đủ 489 file, cứ build đầy đủ mà kiểm. (2) Build PC fail "file being used by another process" = cửa sổ game thử đang mở chiếm .exe. (3) Ảnh map KHÔNG phải một file: ghép từ ~40-50 mảnh PNG 96×96 (mức x4) trong `zz1\data\res\x1..x4`, tên dạng `<bộ>$<số>`, 1.285 mảnh tổng.
- **Trang nạp thẻ** — user đã quyết định bỏ, không làm.
- `zz1\sql\` trống; file schema thật là `zz1\z.sql` (**bản CŨ 47 bảng — đừng dùng**; bản đúng là `goi-len-vps\z-moi-nhat.sql`, 48 bảng).

## `toolcaigame.exe`

`Toocaigame\toolcaigame.exe` là **trình tự động cài môi trường cho VPS Windows** (.NET 8): cài JDK 21, XAMPP, Chrome, Notepad++, WinRAR, UniKey; sửa PATH/JAVA_HOME; đổi timezone; mở port firewall.

**Chưa từng chạy, và nên thận trọng**: nó được build lại từ source decompile (PDB còn nguyên `C:\Users\datvp\...\toolcaigame_decompiled\`), chạy quyền admin, và tải archive từ Google Drive cá nhân — link có thể bị đổi nội dung bất cứ lúc nào. Nếu dùng, tải zip đó về kiểm tra riêng trước, và chỉ chạy trên VPS dùng thử xoá đi làm lại được.
