# Danh sách chính xác những gì phải chép lên VPS

Mọi dung lượng dưới đây tôi **đo thật trên máy bạn** ngày 21/07/2026, không phải ước lượng.

**Tổng phải chép: khoảng 862 MB (0,85 GB).**
Ổ C: của VPS có 70 GB, nên dung lượng hoàn toàn không phải vấn đề. Cái đáng lo là
**thời gian chép** — riêng thư mục ảnh có hơn 71.000 file nhỏ, chép qua Remote Desktop
kiểu kéo-thả sẽ rất lâu và hay đứt giữa chừng. Xem mục "Chép bằng cách nào" ở cuối.

---

## PHẦN 1 — Server game Java

Nguồn: `C:\Users\Admin\Music\zz1`

### Bắt buộc chép

| Chép cái gì | Dung lượng | Số file | Vì sao bắt buộc |
|---|---:|---:|---|
| `data\` | **822,1 MB** | 71.296 | Kho ảnh và dữ liệu game mà server gửi xuống cho client. Thiếu là game không hiện được gì. Bên trong nặng nhất: `icon` 545 MB, `img_by_name` 137,4 MB, `res` 51,4 MB, `effect` 37,1 MB, `item_bg_temp` 31,9 MB, `mob` 18,7 MB. **Bên trong đây cũng có `data\config\config.properties`** — file cấu hình phải thay bằng bản VPS. |
| `src\` | **4,9 MB** | 491 | Nghe lạ nhưng **bắt buộc**: trang quản trị đọc thẳng mã nguồn Java bằng regex (`game_bridge.php` dòng 194-236) để hiện tỉ lệ nâng sao Pha Lê. Không có `src\` thì trang đó hiện toàn số 0. |
| `zz1.jar` | **7,58 MB** | 1 | Chính là server. Đây là "fat jar" — tự chứa sẵn driver MySQL, HikariCP, `ga.jar`. |
| `run.bat` | ~4 KB | 1 | File khởi động, có sẵn vòng canh gác tự chạy lại. Dùng nguyên, đừng viết lại. |
| Các file lẻ ở gốc `zz1` | ~20 KB | 12 | `maintenanceConfig.txt`, `phucloi_config.properties`, `panel_boss_overrides.properties`, `panel_star_upgrade_rates.properties`, `panel_exp_status.txt`, `panel_drop_status.txt`, `panel_admin_mode_status.txt`, `panel_xnxx_status.txt`, `panel_boss_persist_status.txt`, `panel_server_status.json`, `build.xml`, `manifest.mf`. Đây là **mọi thiết lập bạn đã chỉnh trong trang quản trị** (tỉ lệ EXP, drop, boss, phúc lợi, nâng sao). Bỏ qua thì game quay về số mặc định trong code. Nhẹ tí xíu, cứ chép hết cho chắc. |
| `lib\` | **11,1 MB** | 2 | *Không cần để chạy* (vì `zz1.jar` đã tự chứa hết). Nhưng cần nếu sau này muốn **build lại server ngay trên VPS**. 11 MB không đáng gì — nên chép. |

**Cộng phần 1: khoảng 845,7 MB.**

### Bỏ qua được

| Bỏ cái gì | Tiết kiệm | Vì sao bỏ |
|---|---:|---|
| `out\` | 6,6 MB | Chứa `done.jar` — bản build tạm lúc thử nghiệm. Không phải bản đang chạy. |
| `dist\` | 5,5 MB | Bản build Ant cũ (`LIB_SALE.jar` + README). Cũ hơn bản đang chạy. |
| `nbproject\` | 0,2 MB | Cấu hình riêng của phần mềm NetBeans trên máy bạn. VPS không cài NetBeans. |
| `build\` | ~0 MB | Rác sinh ra lúc biên dịch (`built-jar.properties`, `verify-sources.txt`...). |
| `logs\` | ~0 MB | Nhật ký chạy ở máy nhà. `run.bat` tự tạo lại thư mục này trên VPS (dòng 44). |
| `sql\` | 0 MB | Thư mục rỗng. |
| `2026-07-21_truoc-hash-matkhau\` | 0 MB | Thư mục rỗng (bản sao lưu trước khi đổi mật khẩu sang mã hoá, đã dọn). |
| `run.bat.goc` | ~2 KB | Bản `run.bat` gốc trước khi sửa. Giữ ở máy nhà làm kỷ niệm, đừng đưa lên. |
| `sources.txt` | 20 KB | **Nên bỏ hẳn.** Đây là danh sách file nguồn cũ, **thiếu 6 file**. Để trên VPS thì lần sau ai đó (kể cả AI) dùng nó để build lại là hỏng. |
| `src\server\ServerManager.java.bak-update-shop` | 22,7 KB | **Nên bỏ hẳn.** File này không được biên dịch (đuôi `.bak-...`), nhưng chứa mã cũ tự gọi lại `run.bat` gây lỗi đẻ cửa sổ chồng chất, và còn tên miền chủ cũ. Trang quản trị lại quét thư mục `src\` bằng regex — để rác trong đó là rủi ro thừa. **Chuyển sang thư mục sao lưu riêng ở máy nhà.** |

Tổng tiết kiệm chỉ ~12,3 MB — **lý do bỏ không phải để nhẹ, mà để sau này không ai
mở nhầm file cũ ra sửa rồi mất công vô ích.**

> **Chú ý: đừng chép nguyên file `src\server\icon.png` đi đâu cả — phải GIỮ.**
> Nó nằm trong `src\` và server cần nó để hiện biểu tượng cửa sổ. Thiếu là server
> văng lỗi ngay khi khởi động (`NullPointerException ... ServerManagerUI.java:53`).

---

## PHẦN 2 — Web PHP

Nguồn: `C:\xampp\htdocs` (đây là bản đang chạy thật; bản `C:\Users\Admin\Music\htdocs`
là bản sao, **đừng lấy bản đó**).

### Bắt buộc chép

| Chép cái gì | Dung lượng | Số file | Ghi chú |
|---|---:|---:|---|
| `sv1\` | **6,33 MB** | 67 | Toàn bộ trang game. **Tên thư mục phải giữ đúng `sv1`** — nó là khoá tra cứu game root trong `game_bridge.php`. Đổi tên là đứt hết liên kết web ↔ server. Bên trong: `kame` 5,30 MB, `admin` 0,74 MB, `down` (rỗng). |
| `img\` | **8,98 MB** | 39 | Ảnh giao diện website. |
| `index.php` | 44 KB | 1 | Trang ngụy trang "EduVerse" ngoài cùng. |
| `privacy-policy.php` | 22 KB | 1 | |
| `ddos_guard.php` | 14 KB | 1 | Lớp chặn spam. |
| `icon_asset.php` | 2 KB | 1 | Cầu nối lấy ảnh icon game từ `zz1\data\icon\x4`. |
| `.htaccess` | 3 KB | 1 | **Dùng bản trong thư mục này**: `htaccess-vps.txt` (nhớ đổi tên thành `.htaccess`). |

**Cộng phần 2: khoảng 15,4 MB.**

### Bỏ qua

| Bỏ cái gì | Tiết kiệm | Vì sao bỏ |
|---|---:|---|
| `zz4.rar` | **5,06 MB** | Một file nén nằm lạc trong thư mục web gốc, không có trang nào dùng đến. `.htaccess` hiện đang chặn tải file `.rar` nên chưa lộ, nhưng không có lý do gì đưa nó lên máy chủ công cộng. **Mở ra xem nó là gì ở máy nhà trước, rồi cất vào chỗ sao lưu.** |
| `sv1\admin\reset_log.txt` | ~1 KB | 3 dòng nhật ký cũ từ lúc chạy thử ở nhà (có IP `::1` và ID admin). Không lộ được ra ngoài vì `.htaccess` trong `admin` đã chặn file `.txt`, nhưng dọn cho sạch. |

---

## PHẦN 3 — Cơ sở dữ liệu ⚠️ ĐỌC KỸ, CHỖ NÀY DỄ HỎNG NHẤT

**KHÔNG được dùng file `zz1\z.sql` để tạo cơ sở dữ liệu trên VPS.**

Ba lý do, đều đã kiểm chứng bằng cách mở file ra đếm:

| | `z.sql` | Cơ sở dữ liệu đang chạy |
|---|---:|---:|
| Số bảng | **47** | **48** |
| Bảng `news` | **không có** | có |
| Cột mật khẩu | chuỗi thường, đọc được | `pbkdf2$10000$…` (83 ký tự) |
| Chữ "nrokame" | **không có ở đâu cả** | có |
| Ngày | 20/07/2026 16:19 | hôm nay |

Cụ thể:

1. **Thiếu bảng `news`** — đây chính là chỗ chênh 47 ↔ 48.
2. **Mật khẩu người chơi trong `z.sql` vẫn là bản chưa mã hoá.** Cơ sở dữ liệu đang
   chạy đã chuyển sang dạng băm PBKDF2. Nhập `z.sql` đè lên là quay ngược việc này.
3. **Không có chữ "nrokame"** — nghĩa là các đường link tải game trong bảng `adminpanel`
   vẫn là bản trước khi đổi thương hiệu.

**Dùng `z.sql` trên VPS = quay ngược lại toàn bộ công sức đổi thương hiệu và mã hoá
mật khẩu, mà không có thông báo lỗi nào.**

> **Đính chính (21/07/2026):** một bản trước của tài liệu này ghi rằng `z.sql` thiếu
> bảng `login_security`. **Câu đó sai** — bảng `login_security` **có** trong `z.sql`
> (dòng 3999). Và nó cũng không liên quan gì đến mã hoá mật khẩu: các cột của nó là
> `failed_count`, `first_failed_at`, `locked_until`, `lock_level` — đây là bảng **khoá
> tài khoản khi đăng nhập sai quá nhiều lần**. Kết luận "đừng dùng `z.sql`" vẫn đúng
> nguyên, nhưng vì 3 lý do đã liệt kê ở trên, không phải vì `login_security`.

### Cách làm đúng: để script tự xuất bản mới

Chạy `2-DONG-GOI-O-MAY-NHA.ps1` — nó xuất bản dump mới **và tự kiểm tra** xem bản vừa
xuất có đủ 48 bảng, có bảng `news`, có `pbkdf2`, có `nrokame` không. Không đủ là nó báo.

Nếu muốn tự gõ tay, đây là lệnh đúng — chạy trên **máy hiện tại**, khi MySQL đang bật:

```powershell
& "C:\xampp\mysql\bin\mysqldump.exe" -u root --default-character-set=utf8mb4 --routines --events --single-transaction --result-file="C:\Users\Admin\Music\bo-cai-vps\goi-len-vps\z-moi-nhat.sql" z
```

> **Chú ý `--result-file`, đừng dùng dấu `>`.** Dấu `>` của PowerShell ghi file theo
> bảng mã của cửa sổ lệnh. Trên máy này bảng mã đang là UTF-8 nên may mà không sao,
> nhưng trên **Windows Server 2019 mặc định thì không phải** — dùng `>` ở đó sẽ làm
> **hỏng toàn bộ chữ tiếng Việt** trong dump (tên vật phẩm, lời thoại NPC…) mà không
> báo lỗi gì. `--result-file` là mysqldump tự ghi thẳng ra file, không qua cửa sổ lệnh.

File tạo ra khoảng **1 MB**. Chép file *đó* lên VPS, không phải `z.sql`.

Trên **VPS**, việc nhập do `3-CAI-GAME-LEN-VPS.ps1` lo. Nếu muốn tự gõ tay:

```powershell
& "C:\xampp\mysql\bin\mysql.exe" -u root -e "CREATE DATABASE z CHARACTER SET utf8mb4;"
& "C:\xampp\mysql\bin\mysql.exe" -u root --default-character-set=utf8mb4 z -e "source C:/duong-dan-toi/z-moi-nhat.sql"
```

> **Chú ý: KHÔNG có dấu `<` trong lệnh trên.** PowerShell **không có** toán tử `<` —
> gõ vào là báo lỗi ngay, chưa kịp gọi mysql:
> `The '<' operator is reserved for future use.`
> (Cách viết `mysql -u root z < file.sql` chỉ chạy trong cửa sổ `cmd` kiểu cũ.)
> Trong PowerShell phải nhờ chính mysql đọc file bằng lệnh `source`, và đường dẫn
> **dùng dấu gạch xuôi `/`**, không phải `\`.

Nhập xong hãy kiểm tra ngay — phải ra đúng **48**:

```powershell
& "C:\xampp\mysql\bin\mysql.exe" -u root -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='z';"
```

> Đặt mật khẩu MySQL **sau khi** nhập xong dữ liệu, rồi mới sửa 3 chỗ chứa mật khẩu.
> Làm ngược lại thì lệnh nhập ở trên sẽ bị từ chối.

Giữ `z.sql` lại ở máy nhà làm bản gốc dự phòng, nhưng **đừng đưa lên VPS** để không ai
lỡ tay nhập nhầm. Script bước 3 có kiểm tra: nếu bạn chỉ vào đúng file tên `z.sql`,
nó sẽ từ chối và giải thích lý do.

---

## PHẦN 4 — File cài đặt cho người chơi (CHƯA CÓ)

Thư mục `C:\xampp\htdocs\sv1\down\` hiện **trống trơn**, chỉ có mỗi file `.htaccess`.
Mọi nút "Tải Game" trên web đang lỗi 404.

Sau khi build lại client Unity (đã đổi IP sang `103.116.38.60`), đặt 3 file vào đó với
**đúng 3 tên này, phân biệt chữ hoa chữ thường**, vì `.htaccess` chỉ cho phép đúng 3 tên:

- `PC.rar`
- `NRO.apk`
- `IPA.ipa`

Đặt sai tên (ví dụ `nro.apk` viết thường) thì người chơi nhận lỗi 403 dù file có thật
trên máy chủ. Dung lượng chưa tính được vì các file này chưa tồn tại.

---

## Bảng tổng kết dung lượng

| Phần | Dung lượng |
|---|---:|
| Server game (`zz1`, đã bỏ rác) | ~845,7 MB |
| Web (`htdocs`, đã bỏ `zz4.rar`) | ~15,4 MB |
| Bản dump CSDL mới | ~1 MB |
| **TỔNG** | **~862 MB** |
| Đã bỏ qua được | ~17,4 MB |
| File cài client cho người chơi | chưa có |

---

## Chép bằng cách nào

Kéo-thả 71.000 file qua Remote Desktop rất hay đứt giữa chừng, và **đứt giữa chừng
nguy hiểm hơn là hỏng hẳn**: bạn sẽ có một thư mục ảnh thiếu vài trăm file, server vẫn
chạy nhưng người chơi thấy vài chỗ trống, rất khó tìm ra.

**Cách an toàn hơn: nén trước rồi chép 1 file duy nhất.**

Trên máy hiện tại (PowerShell), nén phần server game (bỏ luôn rác):

```powershell
$nguon = "C:\Users\Admin\Music\zz1"
$tam   = "C:\Users\Admin\Music\bo-cai-vps\zz1-len-vps"
Remove-Item $tam -Recurse -Force -EA SilentlyContinue
New-Item -ItemType Directory -Force $tam | Out-Null

# Chép các phần bắt buộc
Copy-Item "$nguon\data" $tam -Recurse
Copy-Item "$nguon\src"  $tam -Recurse
Copy-Item "$nguon\lib"  $tam -Recurse
Copy-Item "$nguon\zz1.jar","$nguon\run.bat" $tam
Get-ChildItem $nguon -File -Include *.txt,*.properties,*.json,*.xml,*.mf |
    Where-Object { $_.Name -ne 'sources.txt' } | Copy-Item -Destination $tam

# Bo file rac trong src
Remove-Item "$tam\src\server\ServerManager.java.bak-update-shop" -EA SilentlyContinue

# Nen lai (mat vai phut)
Compress-Archive -Path "$tam\*" -DestinationPath "C:\Users\Admin\Music\bo-cai-vps\zz1-len-vps.zip" -Force
```

Và nén phần web:

```powershell
$w = "C:\Users\Admin\Music\bo-cai-vps\web-len-vps"
Remove-Item $w -Recurse -Force -EA SilentlyContinue
New-Item -ItemType Directory -Force $w | Out-Null
Copy-Item "C:\xampp\htdocs\sv1","C:\xampp\htdocs\img" $w -Recurse
Copy-Item "C:\xampp\htdocs\index.php","C:\xampp\htdocs\privacy-policy.php","C:\xampp\htdocs\ddos_guard.php","C:\xampp\htdocs\icon_asset.php" $w
Remove-Item "$w\sv1\admin\reset_log.txt" -EA SilentlyContinue
Compress-Archive -Path "$w\*" -DestinationPath "C:\Users\Admin\Music\bo-cai-vps\web-len-vps.zip" -Force
```

> File `.htaccess` bắt đầu bằng dấu chấm nên `Copy-Item` thường **bỏ sót nó**.
> Đừng lo — bạn sẽ dùng bản `htaccess-vps.txt` trong thư mục này, chép tay lên VPS
> rồi đổi tên thành `.htaccess`.

Đưa 2 file zip lên VPS bằng cách kéo-thả vào cửa sổ Remote Desktop, hoặc tải lên
Google Drive rồi mở trình duyệt trên VPS tải xuống (thường nhanh và ổn định hơn).

**Sau khi giải nén trên VPS, kiểm tra ngay số file** — phải khớp với máy nhà:

```powershell
(Get-ChildItem "C:\Users\Administrator\Music\zz1\data" -Recurse -File).Count
```

Phải ra **71.296**. Lệch nhiều là giải nén thiếu, làm lại.

---

## Đặt thư mục game ở đâu trên VPS

Trên Windows Server tài khoản thường tên `Administrator`, nên nơi **ít phải sửa code
nhất** là:

```
C:\Users\Administrator\Music\zz1
```

Lý do rất cụ thể: nhiều file đã có sẵn nhánh dự phòng trỏ đúng đường dẫn này —
`sv1\kame\data.php` dòng 23, `icon_asset.php` dòng 7, `game_bridge.php` dòng 40-41 và
81-82. Đặt đúng chỗ này là chúng tự chạy, không phải sửa gì.

Nếu bạn muốn để chỗ khác (ví dụ `C:\NgocRongKame\zz1`) thì phải sửa thêm **5 chỗ**:
4 file PHP vừa nêu, cộng `game_bridge.php` **dòng 131** — dòng này là lỗi thật, nó
**không có** nhánh Administrator (chỉ có `...\Admin\Music\zz1` và `...\Admin\Music\j`),
nên **dù đặt đúng chỗ khuyến nghị vẫn phải sửa dòng 131**. Bỏ qua nó thì bấm Lưu ở
trang quản trị Phúc Lợi sẽ báo thành công nhưng cấu hình không tới được server game.
