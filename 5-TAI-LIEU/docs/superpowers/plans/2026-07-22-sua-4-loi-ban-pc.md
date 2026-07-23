# Kế hoạch: Sửa 4 lỗi bản PC + chuyển đăng ký về web

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Đăng nhập sai phải báo lỗi thay vì treo; đăng ký chỉ còn ở web; màn chọn server 4 nút cố định cân giữa, không nền đen; file tải tên `NROKAME.zip`.

**Architecture:** Server Java sửa 2 điểm (hết trả-về-im-lặng, khoá lệnh đăng ký). Client Unity sửa `ServerListScreen` về 4 nút cố định và xoá `AccountRegisterScreen` — làm y hệt ở **cả Game1 lẫn Game2**. Web đổi tên file tải. Kiểm chứng theo cách A: tái hiện lỗi ở máy nhà trước khi sửa.

**Tech Stack:** Java 21 (javac trực tiếp, không Ant), Unity 2022.3.62f2 batchmode, PHP 8/XAMPP, MariaDB 10.4.

## Global Constraints

- KHÔNG có git ở repo này → thay "commit" bằng: tạo bản sao `.bak-YYYYMMDD` trước khi sửa mỗi file, và chạy bước kiểm chứng của task trước khi sang task sau.
- Client có HAI bản song song `Assets\Scripts\Game1\` và `Assets\Scripts\Game2\` — mọi sửa đổi client phải làm ở cả hai, xong task phải diff hai bản để soát.
- Chuỗi hiển thị cho người chơi: tiếng Việt có dấu. Tên miền viết đúng `nrokame.com`.
- Build server: 3 điểm bắt buộc trong CLAUDE.md (LIB_SALE.jar, liệt kê lại toàn bộ .java không BOM, chép icon.png). Build client: luôn truyền `-buildTarget` khớp hàm gọi.
- KHÔNG sửa `src\network\` (code chết decompile). KHÔNG format lại `CombineService.java` (panel web parse bằng regex).
- Server đang chạy trên VPS không được đụng tới cho đến Task 9.

---

### Task 1: Server — hết trả-về-im-lặng khi đăng nhập

**Files:**
- Modify: `C:\Users\Admin\Music\zz1\src\jdbc\daos\NDVSqlFetcher.java:78-91` (nhánh sai mật khẩu), `:173-177` (câu báo sai tài khoản), `:178-184` (khối catch)

**Interfaces:**
- Consumes: `Service.gI().sendThongBaoOK(MySession, String)`, `Service.gI().sendLoginFail(MySession, boolean)`, `al.wrong()` — đều có sẵn, xem cách dùng ở dòng 174-176 cùng file.
- Produces: không đổi chữ ký hàm nào. Mọi đường ra khỏi `login()` đều gửi ít nhất một gói tin cho client.

- [ ] **Step 1: Sao lưu**

```powershell
Copy-Item C:\Users\Admin\Music\zz1\src\jdbc\daos\NDVSqlFetcher.java C:\Users\Admin\Music\zz1\src\jdbc\daos\NDVSqlFetcher.java.bak-20260722
```

- [ ] **Step 2: Sửa nhánh sai mật khẩu (dòng 80-82)**

Thay:

```java
                if (!auth.PasswordUtil.verify(session.pp, storedPassword)) {
                    return null;
                }
```

bằng:

```java
                if (!auth.PasswordUtil.verify(session.pp, storedPassword)) {
                    // Trước đây return null mà không gửi gì -> client đợi mãi,
                    // treo "Xin chờ" rồi tự thoát. Mọi đường ra đều phải báo.
                    Service.gI().sendThongBaoOK(session,
                            "Sai tài khoản hoặc mật khẩu. Chưa có tài khoản? Đăng ký tại nrokame.com");
                    Service.gI().sendLoginFail(session, false);
                    al.wrong();
                    return null;
                }
```

- [ ] **Step 3: Đồng bộ câu báo nhánh không-có-tài-khoản (dòng ~174)**

Thay chuỗi `"Thông tin tài khoản hoặc mật khẩu không chính xác"` bằng
`"Sai tài khoản hoặc mật khẩu. Chưa có tài khoản? Đăng ký tại nrokame.com"`.
(Hai nhánh cố ý dùng cùng một câu — không để lộ tài khoản nào tồn tại.)

- [ ] **Step 4: Khối catch cũng phải báo (dòng 178-184)**

Thay:

```java
        } catch (Exception e) {
            Logger.error(session.uu);
            if (player != null) {
                player.dispose();
                player = null;
            }
            Logger.logException(NDVSqlFetcher.class, e);
        } finally {
```

bằng:

```java
        } catch (Exception e) {
            Logger.error(session.uu);
            if (player != null) {
                player.dispose();
                player = null;
            }
            Logger.logException(NDVSqlFetcher.class, e);
            // Lỗi bất ngờ cũng không được im lặng — client sẽ treo y hệt.
            Service.gI().sendThongBaoOK(session, "Không đăng nhập được, vui lòng thử lại sau.");
            Service.gI().sendLoginFail(session, false);
        } finally {
```

- [ ] **Step 5: Kiểm nhanh bằng biên dịch một-file**

```powershell
& "C:\Program Files\Java\jdk-21\bin\javac.exe" -nowarn -encoding UTF-8 --release 17 -cp "C:\Users\Admin\Music\zz1\lib\LIB_SALE.jar;C:\Users\Admin\Music\zz1\src" -d "$env:TEMP\chk1" C:\Users\Admin\Music\zz1\src\jdbc\daos\NDVSqlFetcher.java
```

Mong đợi: không có dòng `error:` nào (cảnh báo được phép). Build đầy đủ nằm ở Task 3.

---

### Task 2: Server — khoá lệnh đăng ký trong game + sửa chỗ băm

**Files:**
- Modify: `C:\Users\Admin\Music\zz1\src\services\Service.java:397` (đầu hàm `registerAccount`), `:456` (chỗ setString mật khẩu)

**Interfaces:**
- Consumes: `sendThongBaoOK(MySession, String)` (hàm cùng lớp), `auth.PasswordUtil.hash(String)` (Task 1 không đổi nó).
- Produces: `registerAccount()` giữ nguyên chữ ký; từ nay luôn từ chối và trả lời bằng thông báo. Client cũ bấm Đăng ký sẽ nhận câu chỉ ra web — không treo.

- [ ] **Step 1: Sao lưu**

```powershell
Copy-Item C:\Users\Admin\Music\zz1\src\services\Service.java C:\Users\Admin\Music\zz1\src\services\Service.java.bak-20260722
```

- [ ] **Step 2: Chặn ngay đầu hàm (sau dòng 397-398)**

Ngay sau `Connection connection = null;` chèn:

```java
        // 2026-07-22: Đăng ký chỉ còn ở web (nrokame.com). Khoá ở server chứ
        // không chỉ gỡ nút client — client bị sửa vẫn gọi thẳng được lệnh này.
        // Phần thân bên dưới giữ lại phòng sau này muốn bật lại.
        if (true) {
            sendThongBaoOK(session,
                    "Vui lòng đăng ký tài khoản tại nrokame.com rồi quay lại đăng nhập.");
            return;
        }
```

(`if (true)` để javac không báo "unreachable code" cho phần thân giữ lại.)

- [ ] **Step 3: Sửa chỗ băm trong phần thân giữ lại (dòng ~456)**

Thay `statement.setString(2, pass);` bằng
`statement.setString(2, auth.PasswordUtil.hash(pass));`
— nếu sau này bật lại thì không tái tạo lỗi mật khẩu thô.

- [ ] **Step 4: Kiểm nhanh biên dịch một-file**

```powershell
& "C:\Program Files\Java\jdk-21\bin\javac.exe" -nowarn -encoding UTF-8 --release 17 -cp "C:\Users\Admin\Music\zz1\lib\LIB_SALE.jar;C:\Users\Admin\Music\zz1\src" -d "$env:TEMP\chk2" C:\Users\Admin\Music\zz1\src\services\Service.java
```

Mong đợi: không có `error:`.

---

### Task 3: Build lại zz1.jar

**Files:**
- Create: `C:\Users\Admin\Music\zz1\zz1.jar` (ghi đè, sao lưu trước)

**Interfaces:**
- Consumes: mã đã sửa ở Task 1-2.
- Produces: jar mới cho Task 4 (thử máy nhà) và Task 9 (VPS).

- [ ] **Step 1: Sao lưu jar đang chạy**

```powershell
Copy-Item C:\Users\Admin\Music\zz1\zz1.jar C:\Users\Admin\Music\zz1\zz1.jar.bak-20260722
```

- [ ] **Step 2: Build theo đúng quy trình đã kiểm chứng trong CLAUDE.md**

```powershell
cd C:\Users\Admin\Music\zz1
$out = "$env:TEMP\zz1-build"; $list = "$env:TEMP\zz1-files.txt"
Remove-Item $out -Recurse -Force -ErrorAction SilentlyContinue
$files = Get-ChildItem src -Recurse -Filter *.java | ForEach-Object { $_.FullName }
[IO.File]::WriteAllLines($list, $files)   # KHÔNG BOM — WriteAllLines mặc định là đúng
& "C:\Program Files\Java\jdk-21\bin\javac.exe" -nowarn -encoding UTF-8 --release 17 -cp "lib\LIB_SALE.jar" -d $out "@$list"
Copy-Item src\server\icon.png $out\server\
```

Lưu ý: `Get-ChildItem src -Recurse` sẽ nhặt cả file `.bak-20260722` **không** — filter `*.java` đã loại chúng, nhưng KHÔNG được đặt bản sao lưu đuôi `.java` ở trong `src\`.

- [ ] **Step 3: Gộp ga.jar + đóng gói (đúng cách đã làm 2026-07-21)**

```powershell
cd $out
& "C:\Program Files\Java\jdk-21\bin\jar.exe" xf C:\Users\Admin\Music\zz1\lib\ga.jar
Remove-Item META-INF\*.SF, META-INF\*.DSA, META-INF\*.RSA -ErrorAction SilentlyContinue
"Main-Class: server.ServerManager" | Out-File "$env:TEMP\manifest.txt" -Encoding ascii
& "C:\Program Files\Java\jdk-21\bin\jar.exe" cfm C:\Users\Admin\Music\zz1\zz1.jar "$env:TEMP\manifest.txt" .
```

- [ ] **Step 4: Kiểm chứng jar**

```powershell
& "C:\Program Files\Java\jdk-21\bin\jar.exe" tf C:\Users\Admin\Music\zz1\zz1.jar | Select-String 'server/ServerManager.class|server/icon.png|jdbc/daos/NDVSqlFetcher.class' 
```

Mong đợi: đủ cả 3 dòng.

---

### Task 4: Tái hiện lỗi ở máy nhà → xác nhận sửa xong (cách A)

**Files:**
- Create: game root thử `C:\NRO_Thu\` (mklink /J trỏ asset nặng), DB thử `z_thu`
- Không sửa mã trong task này — chỉ chạy và quan sát.

**Interfaces:**
- Consumes: `zz1.jar` mới (Task 3), quy trình test-root trong CLAUDE.md.
- Produces: bằng chứng TRƯỚC (client cũ + jar cũ treo khi sai mật khẩu) và SAU (jar mới hiện câu báo). Ghi nhận nguyên nhân nền đen cho Task 6.

- [ ] **Step 1: Dựng DB thử**

```powershell
C:\xampp\mysql\bin\mysql.exe -u root -e "DROP DATABASE IF EXISTS z_thu; CREATE DATABASE z_thu CHARACTER SET utf8mb4;"
Get-Content C:\Users\Admin\Music\goi-len-vps\z-moi-nhat.sql | C:\xampp\mysql\bin\mysql.exe -u root z_thu
C:\xampp\mysql\bin\mysql.exe -u root z_thu -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='z_thu';"
```

Mong đợi: 48 bảng.

- [ ] **Step 2: Dựng game root thử** — copy `config\` + `update_data\` thật, mklink /J các thư mục nặng (`data\icon`, `data\img_by_name`, `data\big_image`…), sửa trong bản config đó: `server.port=14446`, `database.name=z_thu`. Chép `zz1.jar` CŨ (`.bak-20260722`) vào trước.

- [ ] **Step 3: BẰNG CHỨNG TRƯỚC** — chạy jar cũ (`java -jar zz1.jar` bằng JDK 21, working dir = `C:\NRO_Thu`), build client thử trỏ `127.0.0.1:14446` (sửa tạm `ListIP`/`linkDefault` ở cả Game1+Game2, build `KameBuild.BuildWindows`, KHÔNG phát hành bản này). Đăng nhập tài khoản có thật với **mật khẩu sai**.

Mong đợi: treo "Xin chờ" rồi văng ra — đúng hiện tượng user tả. Nếu KHÔNG tái hiện được: DỪNG, ghi lại chính xác điều quan sát được, xem nhật ký server, báo lại trước khi đi tiếp.

- [ ] **Step 4: Cùng phiên đó, quan sát nền đen** — mở client lần đầu (xoá thư mục lưu `%LOCALAPPDATA%Low\Kame\` trước để giả lập người chơi mới), ghi nhận: màn chọn server có nền đen không, xuất hiện lúc nào, `loadScreen` bật khi nào. Kết luận nguyên nhân ghi vào ghi chú cho Task 6.

- [ ] **Step 5: BẰNG CHỨNG SAU** — thay `zz1.jar` MỚI vào `C:\NRO_Thu`, chạy lại, đăng nhập sai mật khẩu.

Mong đợi: hiện hộp thoại "Sai tài khoản hoặc mật khẩu. Chưa có tài khoản? Đăng ký tại nrokame.com", quay về màn đăng nhập, KHÔNG treo. Đăng nhập đúng mật khẩu: vào game bình thường. Bấm Đăng ký (client cũ còn nút): nhận "Vui lòng đăng ký tài khoản tại nrokame.com…".

- [ ] **Step 6: Kiểm mật khẩu thô không còn được tạo**

```powershell
C:\xampp\mysql\bin\mysql.exe -u root z_thu -e "SELECT username, LEFT(password,7) FROM account ORDER BY id DESC LIMIT 5;"
```

Mong đợi: không có dòng mới nào ngoài định dạng `pbkdf2$` (vì lệnh đăng ký đã bị khoá — không tạo được tài khoản mới từ game).

---

### Task 5: Client — 4 nút cố định, xoá màn đăng ký (Game1 + Game2)

**Files:**
- Modify: `D:\NRO_Project\Assets\Scripts\Game1\ServerListScreen.cs` (initCommand :180-276, paint :391+:396, perform case 3 :834-840, case 10100 :879-881, switchToMe :739-745, switchToMe2 :766-772)
- Delete: `D:\NRO_Project\Assets\Scripts\Game1\AccountRegisterScreen.cs`
- Lặp lại y hệt cho `Game2\` (số dòng có thể lệch vài dòng — tìm theo nội dung).

**Interfaces:**
- Consumes: `LoginScr`, `Rms`, `Command(caption, listener, idAction, p)`, `mResources.change_account`, `mResources.option`, `mResources.server` — đều có sẵn.
- Produces: `cmd` luôn đúng 4 phần tử: `[0]`=Đăng nhập(action 3), `[1]`=Đổi tài khoản(7), `[2]`=Máy chủ(17), `[3]`=Tùy chọn(8). Mọi chỗ khác trong file tra nút theo chỉ số cố định này, không còn `nCmdPlay`.

- [ ] **Step 1: Sao lưu cả hai file ở cả hai Game**

```powershell
foreach ($g in 'Game1','Game2') {
  Copy-Item "D:\NRO_Project\Assets\Scripts\$g\ServerListScreen.cs" "D:\NRO_Project\Assets\Scripts\$g\ServerListScreen.cs.bak-20260722"
  Copy-Item "D:\NRO_Project\Assets\Scripts\$g\AccountRegisterScreen.cs" "D:\NRO_Project\Assets\Scripts\$g\AccountRegisterScreen.cs.bak-20260722"
}
Move-Item "D:\NRO_Project\Assets\Scripts\Game1\ServerListScreen.cs.bak-20260722","D:\NRO_Project\Assets\Scripts\Game2\ServerListScreen.cs.bak-20260722","D:\NRO_Project\Assets\Scripts\Game1\AccountRegisterScreen.cs.bak-20260722","D:\NRO_Project\Assets\Scripts\Game2\AccountRegisterScreen.cs.bak-20260722" "C:\Users\Admin\Music\banluu-client-20260722\" -Force
```

(Bản `.bak` phải chuyển RA NGOÀI `Assets\` — đuôi `.cs.bak…` không sao nhưng Unity vẫn tạo `.meta` rác; để hẳn ngoài dự án cho sạch. Tạo thư mục trước nếu chưa có.)

- [ ] **Step 2: Viết lại initCommand (thay toàn bộ :180-276)**

```csharp
    	private void initCommand()
    	{
    		// 4 nút cố định. Trước đây số nút thay đổi theo nCmdPlay (có tài
    		// khoản lưu hay không) khiến vị trí xê dịch và chỉ số nút phải cộng
    		// trừ ở khắp nơi — nguồn lỗi. Đăng ký giờ chỉ còn ở web.
    		nCmdPlay = 0;
    		cmd = new Command[4];
    		cmd[0] = new Command("Đăng nhập", this, 3, null);
    		cmd[1] = new Command(mResources.change_account, this, 7, null);
    		cmd[2] = new Command(string.Empty, this, 17, null); // caption đặt ở switchToMe
    		cmd[3] = new Command(mResources.option, this, 8, null);
    		int num = GameCanvas.hh - 15 * cmd.Length + 28;
    		for (int i = 0; i < cmd.Length; i++)
    		{
    			cmd[i].y = num;
    			cmd[i].setType();
    			cmd[i].x = (GameCanvas.w - cmd[i].w) / 2;
    			num += 30;
    		}
    	}
```

Giữ nguyên công thức vị trí cũ — nó đã căn giữa ngang và giờ ổn định dọc vì
`cmd.Length` cố định. Không sáng tác công thức mới.

- [ ] **Step 3: perform() — case 3 bỏ đường mở màn đăng ký, xoá case 10100**

Trong case 3, thay khối :836-840:

```csharp
					if (Rms.loadRMSString("acc") == null || Rms.loadRMSString("acc").Equals(string.Empty))
					{
						new AccountRegisterScreen().switchToMe();
						return;
					}
```

bằng: (không thay gì — XOÁ hẳn 5 dòng đó; phần dưới của case 3 đã tự xử lý
trường hợp chưa có tài khoản lưu: mở LoginScr với ô trống). Xoá tiếp toàn bộ:

```csharp
			case 10100:
				new AccountRegisterScreen().switchToMe();
				return;
```

- [ ] **Step 4: Đổi mọi chỉ số nút động về tĩnh**

Tìm toàn file các mẫu sau và thay:
- `cmd[2 + nCmdPlay]` → `cmd[2]` (paint :391, :396; switchToMe :739; switchToMe2 :766)
- `cmd[1 + nCmdPlay]` → `cmd[1]` (:741, :768)
- Khối `if (cmd.Length == 5 + nCmdPlay) { cmd[4 + nCmdPlay].caption = mResources.option; }` → `cmd[3].caption = mResources.option;` (:742-745, :769-772)

Sau đó kiểm không còn sót:

```powershell
Select-String -Path D:\NRO_Project\Assets\Scripts\Game1\ServerListScreen.cs -Pattern 'nCmdPlay\s*[+\]]'
```

Mong đợi: 0 kết quả (biến `nCmdPlay` còn được gán `= 0` thì không sao).

- [ ] **Step 5: Xoá AccountRegisterScreen.cs (+ file .meta cùng tên) ở cả hai Game, rồi quét tham chiếu sót**

```powershell
Remove-Item D:\NRO_Project\Assets\Scripts\Game1\AccountRegisterScreen.cs*, D:\NRO_Project\Assets\Scripts\Game2\AccountRegisterScreen.cs*
Get-ChildItem D:\NRO_Project\Assets -Recurse -Include *.cs | Select-String -Pattern 'AccountRegisterScreen|requestAccountRegistration' | Select-Object Path, LineNumber, Line
```

Mong đợi: chỉ còn `Service.cs` (hàm `requestAccountRegistration` — giờ không ai
gọi). Xoá luôn hàm đó ở `Game1\Service.cs:798-…` và bản Game2 tương ứng
(khối `public void requestAccountRegistration … }` trọn vẹn), chạy lại lệnh
quét đến khi ra 0 kết quả.

- [ ] **Step 6: Lặp lại Step 2-5 cho Game2** (tìm theo nội dung, không tin số dòng), rồi diff hai bản:

```powershell
fc.exe /N D:\NRO_Project\Assets\Scripts\Game1\ServerListScreen.cs D:\NRO_Project\Assets\Scripts\Game2\ServerListScreen.cs
```

Mong đợi: hai file vốn không giống nhau 100%, nhưng các vùng vừa sửa phải
giống nhau từng chữ — soát bằng mắt phần chênh lệch.

- [ ] **Step 7: Kiểm biên dịch bằng build PC thử**

```powershell
D:\Unity\2022.3.62f2\Editor\Unity.exe -batchmode -quit -nographics -projectPath D:\NRO_Project -buildTarget Win64 -executeMethod KameBuild.BuildWindows -logFile "$env:TEMP\kiem-client.log"
Select-String -Path "$env:TEMP\kiem-client.log" -Pattern 'error CS|\[KameBuild\]'
```

Mong đợi: không có `error CS`, có dòng `[KameBuild] XONG`.

---

### Task 6: Client — nền đen màn chọn server (sửa theo quan sát Task 4)

**Files:**
- Modify: `D:\NRO_Project\Assets\Scripts\Game1\ServerListScreen.cs:313-321` (+ bản Game2)

**Interfaces:**
- Consumes: kết luận nguyên nhân từ Task 4 Step 4. `GameCanvas.paintBGGameScr(g)` — hàm vẽ nền có sẵn, chính nhánh else đang dùng.
- Produces: màn chọn server không bao giờ đen kịt.

- [ ] **Step 1: Đối chiếu quan sát Task 4.** Ứng viên chính (nghi ngờ từ khảo sát): khối

```csharp
    		if (!loadScreen)
    		{
    			g.setColor(0);
    			g.fillRect(0, 0, GameCanvas.w, GameCanvas.h);
    		}
    		else
    		{
    			GameCanvas.paintBGGameScr(g);
    		}
```

tô đen toàn màn khi chưa tải dữ liệu lần đầu.

- [ ] **Step 2: Nếu quan sát khớp**, thay khối trên bằng:

```csharp
    		// Trước đây khi chưa tải dữ liệu (loadScreen=false) tô đen toàn màn
    		// -> lần mở đầu tiên người chơi nhìn thấy màn đen. Vẽ nền game trong
    		// mọi trường hợp; paintBGGameScr tự xử lý khi ảnh nền chưa sẵn sàng.
    		GameCanvas.paintBGGameScr(g);
```

Nếu quan sát KHÔNG khớp (nền đen đến từ chỗ khác): DỪNG, ghi lại điều thấy
được kèm ảnh chụp, báo user trước khi sửa khác đi.

- [ ] **Step 3: Áp dụng y hệt cho Game2, build thử lại như Task 5 Step 7, rồi chạy client thử ở môi trường Task 4** (xoá `%LOCALAPPDATA%Low\Kame\` để giả lập lần mở đầu).

Mong đợi: lần mở đầu không còn màn đen; 4 nút thẳng hàng giữa màn hình.

---

### Task 7: Web — đổi tên file tải thành NROKAME.zip

**Files:**
- Modify: `C:\xampp\htdocs\sv1\download.php:8`, `C:\xampp\htdocs\sv1\tai-game.php:92`, `C:\xampp\htdocs\sv1\down\.htaccess:5,14-17`
- Rename: `C:\xampp\htdocs\sv1\down\PC.zip` → `NROKAME.zip` (làm CUỐI, sau khi bản PC mới build xong ở Task 8 — không đổi tên file cũ rồi build đè)

**Interfaces:**
- Consumes: `tg_download_info()` tự kiểm `is_file` — đổi tên xong thẻ Windows tự "Sẵn sàng", không phải sửa gì thêm.
- Produces: đường tải mới `/sv1/download?file=NROKAME.zip`. Đường cũ `?file=PC.zip` trả 404 (chấp nhận — chưa phát link cũ rộng rãi).

- [ ] **Step 1: Sao lưu 3 file** (cùng mẫu `.bak-20260722`, để ngoài thư mục web: `C:\Users\Admin\Music\banluu-web-20260722\`).

- [ ] **Step 2: `download.php`** — thay dòng `'PC.zip' => 'application/zip',` bằng `'NROKAME.zip' => 'application/zip',`.

- [ ] **Step 3: `tai-game.php`** — thay `'file'   => 'PC.zip',` bằng `'file'   => 'NROKAME.zip',`.

- [ ] **Step 4: `down\.htaccess`** — dòng 5: `^(?:PC\.zip|NRO\.apk|IPA\.ipa)$` → `^(?:NROKAME\.zip|NRO\.apk|IPA\.ipa)$`; khối `<Files "PC.zip">…</Files>` → `<Files "NROKAME.zip">` và `filename=NROKAME.zip`.

- [ ] **Step 5: Kiểm chứng (sau khi Task 8 đặt file mới vào chỗ)**

```powershell
Invoke-WebRequest -Uri 'http://localhost/sv1/download?file=NROKAME.zip' -Method Head -UseBasicParsing | Select-Object -Expand Headers
```

Mong đợi: 200, `Content-Type: application/zip`, `filename="NROKAME.zip"`. Và
`?file=PC.zip` phải trả 404. Trang `/tai-game` thẻ Windows "Sẵn sàng", tên tệp
hiện `NROKAME.zip`.

---

### Task 8: Build lại 3 bản phát hành

**Files:**
- Create: `D:\NRO_Build\PC\` → đóng gói `C:\xampp\htdocs\sv1\down\NROKAME.zip`; `D:\NRO_Build\Android\NgocRongKame.apk`; `D:\NRO_Build\iOS\` → `D:\NRO_iOS_GuiMac.zip`

**Interfaces:**
- Consumes: client đã sửa (Task 5-6) VỚI địa chỉ server đã trả về `nrokame.com` (hoàn nguyên sửa tạm `127.0.0.1:14446` của Task 4 Step 3 — kiểm bằng grep trước khi build!).
- Produces: 3 gói phát hành cho Task 9.

- [ ] **Step 1: KIỂM địa chỉ server trước khi build**

```powershell
Select-String -Path D:\NRO_Project\Assets\Scripts\Game1\ServerListScreen.cs, D:\NRO_Project\Assets\Scripts\Game2\ServerListScreen.cs -Pattern 'ListIP|linkDefault' | Select-String -Pattern '127\.0\.0\.1|14446'
```

Mong đợi: 0 kết quả. Có kết quả = còn địa chỉ thử, PHẢI hoàn nguyên đã.

- [ ] **Step 2: Build PC** (lệnh như Task 5 Step 7). Đóng gói: loại thư mục `*_BackUpThisFolder_ButDontShipItWithYourGame`, nén bằng **tar.exe** (không dùng ZipFile .NET — lỗi dấu `\`):

```powershell
cd D:\NRO_Build\PC
tar.exe -a -c -f C:\xampp\htdocs\sv1\down\NROKAME.zip --exclude '*BackUpThisFolder*' *
```

Kiểm: mở zip đếm mục chứa `\` phải = 0; dung lượng ~69 MB.

- [ ] **Step 3: Build APK**

```powershell
D:\Unity\2022.3.62f2\Editor\Unity.exe -batchmode -quit -nographics -projectPath D:\NRO_Project -buildTarget Android -executeMethod KameBuild.BuildAndroid -logFile "$env:TEMP\apk2.log"
```

Rồi **USER ký lại** (bấm đúp `bo-cai-vps\KY-APK.bat` — cần mật khẩu khoá, trợ lý không làm thay được). Script tự chép vào `sv1\down\NRO.apk` + `goi-len-vps\`. Kiểm sau ký: `apksigner verify --print-certs` không còn `Android Debug`.

- [ ] **Step 4: Xuất lại gói iOS**

```powershell
D:\Unity\2022.3.62f2\Editor\Unity.exe -batchmode -quit -nographics -projectPath D:\NRO_Project -buildTarget iOS -executeMethod KameBuild.BuildIOS -logFile "$env:TEMP\ios2.log"
cd D:\NRO_iOS_GuiMac; Remove-Item NgocRongKame-iOS -Recurse -Force
robocopy D:\NRO_Build\iOS D:\NRO_iOS_GuiMac\NgocRongKame-iOS /E /NP /NDL /NJH /NJS
cd D:\NRO_iOS_GuiMac; tar.exe -a -c -f D:\NRO_iOS_GuiMac.zip CLAUDE.md NgocRongKame-iOS
```

(robocopy mã thoát 0-7 là thành công.) Kiểm: zip mở ra 0 mục chứa `\`.

- [ ] **Step 5: Đưa dự án về lại buildTarget Android** (nền tảng hay dùng):

```powershell
D:\Unity\2022.3.62f2\Editor\Unity.exe -batchmode -quit -nographics -projectPath D:\NRO_Project -buildTarget Android -logFile "$env:TEMP\vetarget.log"
```

---

### Task 9: Đưa lên VPS + kiểm đầu-cuối

**Files:**
- Không sửa file nào — chỉ đồng bộ và kiểm.

**Interfaces:**
- Consumes: mọi thứ Task 1-8. `DONG-BO.ps1` đồng bộ zz1.jar, src, htdocs (gồm `down\`), bo-cai-vps; tự bỏ qua 3 file mật khẩu.
- Produces: nrokame.com chạy bản mới.

- [ ] **Step 1: USER chạy trên VPS** (qua Remote Desktop có chia sẻ ổ C):

```powershell
powershell -ExecutionPolicy Bypass -File "\\tsclient\C\Users\Admin\Music\bo-cai-vps\DONG-BO.ps1"
```

Mong đợi: jar khớp MD5 mới, 3 file mật khẩu "con nguyen".

- [ ] **Step 2: USER khởi động lại server game trên VPS** (bấm nút trên cửa sổ server hoặc tắt-bật `run.bat`) — jar mới chỉ có tác dụng sau khi khởi động lại. Chú ý bẫy PAUSE trong run.bat (CLAUDE.md).

- [ ] **Step 3: Kiểm từ máy nhà**

```powershell
Invoke-WebRequest -Uri 'http://nrokame.com/sv1/download?file=NROKAME.zip' -Method Head -UseBasicParsing | Select-Object StatusCode
Test-NetConnection nrokame.com -Port 14445 | Select-Object TcpTestSucceeded
```

Mong đợi: 200 và True.

- [ ] **Step 4: Kiểm người-chơi-thật:** tải `NROKAME.zip` từ web, giải nén, mở game — không màn đen, 4 nút giữa màn; đăng nhập sai → câu báo chỉ ra web; đăng ký tài khoản mới TRÊN WEB → đăng nhập game bằng tài khoản đó → vào được. Đây là tiêu chí nghiệm thu cuối cùng, USER tự bấm.

---

## Ghi chú tổng

- Sau khi mọi task xong: cập nhật CLAUDE.md (tên file tải mới, đăng ký web-only, cơ chế 4 nút cố định) và bộ nhớ trợ lý.
- Các bản `.bak-20260722` giữ ít nhất tới khi Task 9 nghiệm thu xong.
- Task 33 cũ (cài thử APK lên điện thoại) vẫn treo — nhắc user làm với bản APK MỚI sau Task 8.
