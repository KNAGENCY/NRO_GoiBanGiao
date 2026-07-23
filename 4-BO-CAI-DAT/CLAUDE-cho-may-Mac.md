# Hướng dẫn cho trợ lý trên máy Mac — làm bản iOS Ngọc Rồng Kame

Đặt file này tên `CLAUDE.md` ngay cạnh thư mục dự án Xcode khi chép sang Mac.

---

## Bạn đang làm gì

Đây là **dự án Xcode do Unity xuất ra**, không phải mã nguồn game. Việc của bạn
là đóng gói nó thành bản cài iOS rồi đưa lên **TestFlight** để người chơi tải.

**Chủ sở hữu không phải lập trình viên.** Trả lời bằng tiếng Việt đời thường,
nói *tại sao* trước khi nói *làm thế nào*. Đừng bắt họ tự phán đoán kỹ thuật.

Game là server private Ngọc Rồng Online. Bản PC và Android đã phát hành xong,
đang chạy tại `nrokame.com`. iOS là mảnh cuối.

### Đừng sửa mã game ở đây

Mọi thứ trong `Classes/`, `Libraries/`, `Data/` là **máy sinh ra**. Sửa vào đó
sẽ bị xoá sạch ở lần xuất tiếp theo. Muốn đổi gì trong game — kể cả địa chỉ máy
chủ — phải báo về máy Windows sửa mã C# rồi xuất lại.

Sửa được ở đây, và chỉ nên sửa ở đây: cấu hình ký (Signing), Team, ảnh biểu
tượng, `Info.plist`.

---

## Giá trị cố định — không được đổi

| Thứ | Giá trị | Vì sao quan trọng |
|---|---|---|
| Bundle ID | `com.nrokame.game` | Sai một ký tự là không ký được. **Vĩnh viễn** sau khi lên chợ |
| Tên hiển thị | `Ngoc Rong Kame` | Không dấu — cố ý, tránh rắc rối đường dẫn |
| Máy chủ game | `nrokame.com:14445` | Đã nhúng sẵn trong bản build |
| iOS tối thiểu | 12.0 | |
| Thiết bị | iPhone + iPad | |

---

## Cần có trước khi bắt đầu

1. **Xcode** — tải miễn phí từ App Store. Nặng ~10 GB.
2. **Tài khoản Apple Developer 99 USD/năm, đã được Apple duyệt xong.**
   Đăng ký ở `developer.apple.com/programs`, chọn **Individual** (cá nhân).
   Duyệt mất 1–3 ngày. Chưa duyệt xong thì không làm tiếp được.
3. **Bundle ID đã đăng ký**: `developer.apple.com` → Certificates, Identifiers
   & Profiles → Identifiers → **+** → điền đúng `com.nrokame.game`.

Kiểm tra cả ba thứ này **trước** khi động vào Xcode. Thiếu cái nào thì Xcode
báo lỗi đỏ khó hiểu chứ không nói thẳng là bạn thiếu gì.

---

## Quy tắc bảo mật — bắt buộc

Chủ sở hữu có nguyên tắc rõ ràng: **bí mật không rời khỏi máy**. Đã từ chối đưa
khoá ký Android lên cloud dù biết rủi ro mất dữ liệu.

- **Đừng hỏi mật khẩu Apple ID.** Đăng nhập qua Xcode → Settings → Accounts,
  để Xcode giữ trong Keychain. Bạn không cần biết mật khẩu.
- **Đừng yêu cầu app-specific password dán vào chat.** Nếu cần tải lên bằng
  dòng lệnh thì dùng Keychain profile, đừng để mật khẩu trên dòng lệnh — mọi
  thứ trên dòng lệnh đều hiện trong danh sách tiến trình.
- **Đừng đưa chứng chỉ hay khoá ký lên đâu cả** — kể cả để "sao lưu cho chắc".

---

## Các bước

### 0. Ngay sau khi giải nén — cấp lại quyền chạy

Gói này nén trên Windows. Windows không có khái niệm "quyền chạy" như macOS,
nên khi giải nén ra, file script mất quyền đó và Xcode sẽ báo lỗi
`Permission denied` giữa lúc đóng gói.

Chạy ngay trong thư mục dự án:

```bash
chmod +x process_symbols.sh
find . -name "*.sh" -exec chmod +x {} \;
```

Lệnh thứ hai quét cả thư mục con, phòng khi bản xuất sau có thêm script mới.

Cũng nên bỏ các file rác macOS tự sinh nếu gói đi qua nhiều máy:

```bash
xattr -cr .
```

### 1. Mở dự án

Tìm `Unity-iPhone.xcodeproj`, mở bằng Xcode.

Nếu có cả file `.xcworkspace` thì **mở workspace**, không mở project.

### 2. Đặt chữ ký

Chọn mục **Unity-iPhone** ở cột trái → thẻ **Signing & Capabilities**.

- **Automatically manage signing**: đã bật sẵn từ Unity. Nếu chưa thì tích vào.
- **Team**: chọn tài khoản Apple Developer. Đây là thứ duy nhất phải chọn tay.
- **Bundle Identifier**: xác nhận đúng `com.nrokame.game`.

Còn báo lỗi đỏ thì thường là một trong ba:
- Chưa đăng nhập Apple ID trong Xcode → Settings → Accounts
- Bundle ID chưa đăng ký trên trang Apple
- Tài khoản Developer chưa được duyệt xong

### 3. Đóng gói

Trên thanh trên cùng, chọn thiết bị đích là **Any iOS Device (arm64)**.

> Chọn máy giả lập là **Product → Archive bị mờ đi**, bấm không được. Đây là
> chỗ hay tắc nhất.

**Product → Archive**. Lần đầu mất **20–40 phút** — game dùng IL2CPP nên Xcode
phải biên dịch một khối C++ rất lớn. Chạy lâu là bình thường, đừng dừng giữa
chừng.

Muốn làm bằng dòng lệnh thì đại ý:

```bash
xcodebuild -project Unity-iPhone.xcodeproj -scheme Unity-iPhone \
  -configuration Release -destination 'generic/platform=iOS' \
  -archivePath ./build/NgocRongKame.xcarchive archive \
  DEVELOPMENT_TEAM=<TEAM_ID>
```

Rồi `xcodebuild -exportArchive` với một `exportOptions.plist` đặt
`method = app-store-connect`.

### 4. Gửi lên Apple

Xong Archive thì cửa sổ **Organizer** tự hiện → **Distribute App** →
**App Store Connect** → **Upload**.

Đường dòng lệnh: sau khi export ra `.ipa`, dùng công cụ tải lên của Xcode.
**Hãy tự kiểm tra bản Xcode trên máy có `xcrun altool` hay
`xcrun iTMSTransporter`** — Apple đã đổi công cụ này vài lần, đừng tin
trí nhớ. Nếu không chắc, dùng Organizer cho lành; nó luôn đúng.

### 5. Bật TestFlight

1. Vào `appstoreconnect.apple.com` → **My Apps** → **+** → **New App**
   - Platform: iOS
   - Bundle ID: chọn `com.nrokame.game`
   - Tên: `Ngọc Rồng Kame`
2. Sang thẻ **TestFlight**, đợi bản vừa gửi xử lý xong (15–60 phút)
3. Điền **Test Information** — ô **What to Test** bắt buộc có chữ
4. **External Testing** → tạo nhóm → bật **Public Link**

Apple duyệt bản đầu mất **1–2 ngày**. Xong sẽ có địa chỉ dạng:

```
https://testflight.apple.com/join/xxxxxxxx
```

**Đưa địa chỉ đó cho chủ sở hữu** — bên máy Windows sẽ gắn vào nút tải iOS
trên web.

---

## Những chỗ hay hỏng

**Số hiệu bản build phải tăng mỗi lần gửi.**
Gửi trùng số là Apple từ chối ngay lúc tải lên. Bên Windows đã đặt tự tăng mỗi
lần xuất, nhưng nếu bạn Archive lại nhiều lần từ cùng một bản xuất thì phải tự
tăng `CFBundleVersion` trong `Info.plist`.

**Ô "does your app use encryption?" đã khai sẵn.**
`Info.plist` đã có `ITSAppUsesNonExemptEncryption = false`. Đúng với game này —
nó nói chuyện với máy chủ bằng TCP thường, không tự viết thuật toán mã hoá nào.
Đừng xoá dòng đó, xoá là mỗi lần gửi lại bị hỏi.

**Apple cần tài khoản để đăng nhập thử.**
Phần **App Review Information** phải điền sẵn một tài khoản game (tên + mật
khẩu) cho người của Apple vào xem. Không có là bị trả về ngay. Bảo chủ sở hữu
tạo riêng một tài khoản ở `nrokame.com/sv1/dangky` — **đừng dùng tài khoản
quản trị**.

**Máy chủ phải đang bật lúc Apple duyệt.**
Người duyệt mở game mà không kết nối được là đánh trượt. Nhắc chủ sở hữu đừng
tắt VPS trong lúc chờ.

**Bản TestFlight hết hạn sau 90 ngày.**
Hết hạn là người chơi iOS mở game không được nữa. Cứ 90 ngày phải gửi bản mới.
Đây là lý do nhiều server private mất người chơi iOS — quên mất hạn này.

**Rủi ro bị Apple từ chối vì bản quyền.**
Game dùng nhân vật Dragon Ball. Apple **có thể** từ chối vì bản quyền của Toei
Animation / Shueisha. Đây là rủi ro thật, không ai hứa trước được. Nếu bị từ
chối vì lý do này thì gần như không sửa được, trừ khi thay toàn bộ hình ảnh
nhân vật.

Nói thẳng điều này với chủ sở hữu **trước khi** họ nộp 99 USD, nếu họ chưa nộp.

---

## Báo lại gì

Khi xong, báo về:
- Địa chỉ TestFlight (`https://testflight.apple.com/join/...`)
- Số hiệu bản build đã gửi
- Team ID đã dùng (để lần sau xuất từ Windows điền sẵn)

Nếu tắc, chép nguyên văn thông báo lỗi của Xcode — đừng diễn giải lại.
