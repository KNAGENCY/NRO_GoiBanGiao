# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Trả lời bằng **tiếng Việt đời thường**. Chủ sở hữu **không phải lập trình viên** — luôn nói *tại sao* trước khi nói *làm thế nào*.

## ⚠️ ĐỌC `BAT-DAU-TU-DAY.md` TRƯỚC

[`BAT-DAU-TU-DAY.md`](BAT-DAU-TU-DAY.md) ở gốc repo là **điểm vào duy nhất**: dự án đang làm gì, đã xong tới đâu, còn gì chưa làm, đã vấp những lỗi nào. Đọc nó trước khi động vào bất cứ thứ gì.

**Bắt buộc:** làm xong việc gì đáng kể (một task, một kế hoạch, một bản vá, một phát hiện lỗi) thì **cập nhật `BAT-DAU-TU-DAY.md` ngay** rồi mới commit — cách cập nhật ghi ở §11 của chính file đó. Không cập nhật = phiên sau (người hay agent) mất phương hướng.

## Nhiều agent làm song song

Dự án này có thể có nhiều agent cùng làm. Hai lệnh và một sổ:

| Việc | Dùng |
|---|---|
| **Bắt đầu** một phiên/agent mới | `/vao-viec` — tự đọc bối cảnh, xem ai đang làm gì, đề xuất việc tiếp |
| **Kết thúc** một việc | `/ghi-so` — kiểm chứng, ghi nhật ký, cập nhật tài liệu, commit |
| Xem ai vừa chạm vào gì | [`NHAT-KY-AGENT.md`](NHAT-KY-AGENT.md) — sổ chung, dòng mới nhất ở **cuối** bảng |

**Chống đụng nhau:** mỗi agent một **nhánh git riêng** (`git switch -c <ten-viec>`), không hai agent cùng commit lên `main`. `server/` (game) và `web/` là hai vùng tách biệt — làm song song an toàn; hai agent cùng vào `server/internal/game/` thì chắc chắn đụng.

Agent không phải Claude Code (Codex, Cursor…) đọc [`AGENTS.md`](AGENTS.md).

Hook `SessionStart` trong `.claude/settings.json` tự nạp `git log` + nhật ký agent vào đầu mỗi phiên — **chỉ hoạt động khi mở Claude Code trong thư mục dự án này** (dùng `MO-AGENT.bat` cho chắc).

## Việc đang làm: viết lại thành game Naruto

Dự án thật sự hiện nay là **xây game Naruto mới**, viết lại máy chủ bằng **Go + PostgreSQL** (`server/`) và web/panel bằng **Next.js** (`web/`), dựa trên phần lõi giao thức của server Ngọc Rồng đã mua. Đây là **game mới hoàn toàn** — 5 làng Naruto thay 3 hành tinh Dragon Ball; dấu vết game cũ xoá dần.

| Thư mục | Vai trò |
|---|---|
| `server/`, `web/`, `docs/superpowers/`, `docs/my-thuat/` | 🟢 **Đang sống** — mọi việc mới nằm ở đây |
| `1-MAY-CHU-GAME/`, `2-WEB/`, `3-CO-SO-DU-LIEU/`, `4-BO-CAI-DAT/`, `5-TAI-LIEU/` | 📖 **Chỉ tra cứu** — code Java/PHP cũ, dùng để giải mã giao thức nhị phân mà client Unity nói. Đừng sửa để "chạy game" |

Lộ trình: `docs/superpowers/plans/LO-TRINH.md`. Cách làm việc (subagent + TDD + review đối kháng): §10 của `BAT-DAU-TU-DAY.md`.

---

## Phần dưới đây mô tả DI SẢN (server Java/PHP cũ)

Giữ lại vì đây là nguồn tra cứu giao thức và là bản đang chạy thật. **Không phải hướng dẫn cho việc đang làm.**

**Gói bàn giao** của server private Ngọc Rồng Online, thương hiệu **Ngọc Rồng Kame**, tên miền **nrokame.com**. Không phải một project mà là 4 mảnh dính chặt nhau, đóng gói lại để dựng server mới trên VPS Windows:

| Thư mục | Nội dung | Công nghệ |
|---|---|---|
| `1-MAY-CHU-GAME/` | `zz1.jar` + 489 file `.java` trong `src/` + `data/` (asset + config) | Java 21, project Ant/NetBeans tên `NROOCEAN` |
| `2-WEB/htdocs/` | Toàn bộ web: trang chủ, đăng ký, admin panel | PHP 8 (XAMPP/Apache) |
| `3-CO-SO-DU-LIEU/z-moi-nhat.sql` | Schema + dữ liệu gốc, **48 bảng** | MariaDB 10.4, schema `z` |
| `4-BO-CAI-DAT/` | Script cài đặt VPS đánh số 1→5 + `cau-hinh/` | .bat / PowerShell |
| `5-TAI-LIEU/` | `CLAUDE.md` (tài liệu kỹ thuật sâu), `TINH-TRANG-DU-AN.md`, specs | |

**Source client Unity KHÔNG có trong gói** (nằm ở `D:\NRO_Project`, nhiều GB). `nrokame.keystore` (khoá ký APK) cũng không — mất nó là người chơi phải gỡ game cài lại, không ai khôi phục được.

> **`5-TAI-LIEU/CLAUDE.md` là tài liệu quan trọng nhất** — chi tiết build server, build client Unity 3 nền tảng, lịch sử sửa lỗi, tình trạng bảo mật. Đọc nó khi đụng vào bất cứ thứ gì sâu hơn mức "dựng cho chạy". Lưu ý: nó viết theo **đường dẫn máy gốc** (`C:\Users\Admin\Music\zz1`, `C:\xampp\htdocs`, `bo-cai-vps`, `goi-len-vps`) — trong gói này lần lượt là `1-MAY-CHU-GAME`, `2-WEB\htdocs`, `4-BO-CAI-DAT`.

## Lệnh thường dùng

```bat
REM Chạy server — working directory PHẢI là 1-MAY-CHU-GAME (run.bat tự "cd /d %~dp0")
run.bat

REM Import DB (tạo mới)
C:\xampp\mysql\bin\mysql.exe -u root -e "CREATE DATABASE z CHARACTER SET utf8mb4;"
C:\xampp\mysql\bin\mysql.exe -u root z < 3-CO-SO-DU-LIEU\z-moi-nhat.sql
```

Không có test tự động, không có linter, không có staging. Kiểm chứng = chạy thật rồi vào game.

**Cài lên VPS — chạy đúng thứ tự số**, mỗi file tự giải thích ở phần đầu:
`1-CAI-MOI-TRUONG.bat` (JDK 21 + XAMPP + dịch vụ + tường lửa, chạy trên VPS) → `2-DONG-GOI-O-MAY-NHA.ps1` (xuất dump DB mới + nén zz1/htdocs, chạy ở máy nhà) → `3-CAI-GAME-LEN-VPS.ps1` (tạo DB, đặt mật khẩu ở cả 4 chỗ) → `4-TU-DONG-BAT.bat` (AutoAdminLogon + Startup, **chỉ chạy sau khi server đã chạy tay ngon lành**). `5-GO-RA-LAM-LAI.bat` tháo ngược bước 1 và 4. `DONG-BO.ps1` đồng bộ máy nhà → VPS (chạy trên VPS, tự loại trừ 3 file cấu hình).

**Build lại server** (không cần Ant, dùng `javac` JDK 21) — 3 điểm bắt buộc, thiếu cái nào cũng fail: copy `lib\ga.jar` → `lib\LIB_SALE.jar`; liệt kê **toàn bộ 489 file** `.java` vào danh sách **không BOM** (`sources.txt` là danh sách cũ, thiếu 6 file); chép `src\server\icon.png` vào thư mục class trước khi đóng gói. Chi tiết + lệnh đầy đủ ở `5-TAI-LIEU/CLAUDE.md`.

## Ba mảnh nói chuyện với nhau thế nào

Không có API, không có socket giữa web và server. Chỉ ba kênh — hiểu sai chỗ này là sửa sai chỗ:

1. **Database `z` dùng chung** — Java và PHP đọc/ghi cùng bảng. Đây là lý do admin panel sửa được game đang chạy.
2. **`panel_cmd.txt`** — web ghi → Java poll mỗi 2 giây rồi **xoá file**. Lệnh phân cách bằng `|`: `BUFF_VND`, `BUFF_ITEM`, `KICK_PLAYER`, `ADMIN_MODE`, `MAINTENANCE_NOW`, `SET_EXP`, `SET_DROP`, `SPAWN_BOSS`, `EDIT_BOSS`…
3. **File trạng thái JSON** — Java ghi → web đọc mỗi giây: `boss_status.json` (vào `htdocs\sv1\`), `panel_server_status.json` (vào game root). Các file `panel_*.txt/.properties/.json` ở gốc `1-MAY-CHU-GAME` chính là mặt cắt trạng thái này.

Phân chia dữ liệu: **logic + chỉ số game nằm trong MySQL**; `1-MAY-CHU-GAME\data\` là **kho asset client** server stream qua socket theo 4 mức zoom `x1`–`x4`.

## Bẫy — đọc trước khi sửa bất cứ thứ gì

- **Server bắt buộc có màn hình.** `main()` (`src\server\ServerManager.java:167`) chỉ mở cửa sổ Swing; toàn bộ khởi động nằm trong constructor `ServerManagerUI`. Không chạy headless, **không chạy được dạng Windows Service** → khoá vào VPS Windows có AutoAdminLogon.
- **Phải là Java 21.** `zz1.jar` bytecode major 65 dù config ghi `javac.target=17`. `java` trên PATH có thể là JRE 8 → `UnsupportedClassVersionError`.
- **Tầng mạng không có source.** `network.Network`, `network.Session`, `network.SessionManager`, `utils.Logger` nằm trong `lib\ga.jar` đã biên dịch. Thư mục `src\network\` là **bản decompile chết** — sửa vào đó vô tác dụng.
- **Đổi tên thư mục `sv1` là đứt toàn bộ liên kết** — tên thư mục web là khoá tra cứu game root (`sv1\admin\game_bridge.php`).
- **Panel đọc source Java bằng regex** (`game_bridge.php:194-236` parse `CombineService.java`) → phải deploy cả `src\` cạnh jar, và sẽ vỡ nếu format lại code đó.
- **Game root phải ghi được** — `data\update_data\part` bị ghi đè mỗi lần khởi động từ bảng `part`.
- **`1-MAY-CHU-GAME\z.sql` là bản CŨ 47 bảng — ĐỪNG DÙNG.** Thiếu bảng `login_security` (băm mật khẩu) và còn link thương hiệu cũ. Bản đúng là `3-CO-SO-DU-LIEU\z-moi-nhat.sql`. Dùng nhầm = quay ngược mọi thứ mà không có thông báo lỗi nào. `1-MAY-CHU-GAME\sql\` thì trống.
- **Lombok**: `javac` một file lẻ sẽ báo "cannot find symbol getSession()" **giả**. Cứ build đủ 489 file mà kiểm.
- **Dòng `server.sv1` trong `data\config\config.properties` là dòng quan trọng nhất cả file** — đây là địa chỉ server gửi xuống cho client. Để `127.0.0.1` thì máy bạn test vẫn vào được nhưng **không người chơi nào vào được**. Xoá dòng này server lỗi ngay lúc khởi động (không có khoá `server.ip` thay thế).
- **Mật khẩu MySQL phải khớp ở 3 nơi**, sót chỗ nào triệu chứng khác hẳn nhau: `data\config\config.properties` (`database.pass`) → server không chạy; `sv1\config.php:5` → web chết toàn bộ; `sv1\boss_status.php:35` → **web vẫn chạy, chỉ mất ảnh mặt boss**, lỗi bị nuốt trong `catch`, khó chẩn đoán nhất.
- **`.htaccess` nạp `ddos_guard.php` qua `auto_prepend_file` bằng đường dẫn Windows tuyệt đối** (`C:/xampp/htdocs/...`) → đúng trên VPS Windows nếu XAMPP ở đúng `C:\xampp`, nhưng **trên Linux im lặng không chạy**, mất sạch lớp bảo vệ mà không báo lỗi.
- **`index.php` là trang nguỵ trang** (site học tập giả "EduVerse"). Trang game thật ở `/sv1/trangchu`.
- **Không mở cổng 3306 ra internet.** Mật khẩu root đang rỗng trong gói và MariaDB đang nghe mọi card mạng. Chỉ mở 80 + 14445.

## Chưa làm / đang treo (của bản Java/PHP cũ)

> Việc còn lại của **dự án Naruto đang làm** nằm ở `BAT-DAU-TU-DAY.md` §5 (còn gì chưa làm) và §7 (nợ bảo mật). Phần dưới chỉ nói về bản cũ.


Chưa có HTTPS (mật khẩu đăng nhập đi qua mạng dạng thô) · chưa có sao lưu DB tự động · `server.debug=true` nên đổi `false` trước ngày mở cửa (log phình ăn hết ổ C:) · `$SITE_OPEN_TIME` trong `sv1\config.php` vẫn là `2026-07-01` (quá khứ, đồng hồ đếm ngược hiện số âm) · link Zalo/Facebook đang trỏ vòng về chính website · iOS còn thiếu tài khoản Apple Developer.

Tình trạng đầy đủ: `5-TAI-LIEU/TINH-TRANG-DU-AN.md`.
