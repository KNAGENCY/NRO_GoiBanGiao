# Làm bản iOS cho Ngọc Rồng Kame (qua TestFlight)

Tài liệu này viết cho người **không phải lập trình viên**. Đọc hết một lượt trước
khi tiêu tiền, vì có vài chỗ không quay lại được.

Cập nhật: 22/07/2026

---

## Trước hết: vì sao tôi không làm hộ được

Apple bắt buộc bản iOS phải được **biên dịch trên máy Mac**, bằng phần mềm Xcode
của họ. Đây là quy định kỹ thuật cứng, không có cách lách. Máy bạn đang dùng là
Windows nên dừng ở đây.

Unity trên máy bạn chỉ tạo ra được **mã nguồn Xcode**, chưa phải file cài. Muốn
thành file cài phải đưa mã đó sang máy Mac.

---

## Bạn cần chuẩn bị

| Thứ | Giá | Tình trạng |
|---|---|---|
| Máy Mac | — | ✅ Bạn đã có MacBook |
| Xcode | miễn phí | Tải từ App Store trên MacBook, nặng ~10 GB |
| Tài khoản Apple Developer | **99 USD/năm** (~2,5 triệu) | ❌ Chưa có — đây là thứ duy nhất phải mua |

Chỉ còn tài khoản Apple Developer là chưa có. Không có nó thì không dùng được
TestFlight — đây là quy định của Apple, không lách được.

## Phần việc ở máy Windows: ĐÃ LÀM XONG

Ngày 22/07/2026 tôi đã làm sẵn trên máy này:

- Cài bộ iOS Build Support cho Unity (2,7 GB, kiểm mã băm và chữ ký Unity)
- Xuất ra dự án Xcode: **2.706 file, 668 MB**
- Bật sẵn **ký tự động** để bên Mac chỉ việc chọn Team, khỏi tự tạo chứng chỉ
- Đặt **số hiệu bản build tự tăng** mỗi lần xuất, tránh bị Apple từ chối vì trùng
- Khai sẵn `ITSAppUsesNonExemptEncryption` để mỗi lần gửi khỏi bị hỏi về mã hoá

Gói gửi sang Mac: **`D:\NRO_iOS_GuiMac.zip`**

Bên trong có sẵn file `CLAUDE.md` — đưa cho trợ lý trên MacBook đọc, nó sẽ biết
phải làm gì mà không cần bạn giải thích lại.

---

## TestFlight là gì

Là kênh phát hành **thử nghiệm** của Apple. Người chơi cài app TestFlight từ
App Store, rồi mở link mời của bạn để cài game.

**Điểm cần biết trước khi chọn đường này:**

| | |
|---|---|
| Số người tối đa | 10.000 người thử |
| Thời hạn mỗi bản | **90 ngày** — hết hạn là game ngừng mở được |
| Duyệt | Vẫn phải qua Apple duyệt, nhưng nhẹ hơn App Store nhiều |
| Chi phí | Nằm trong 99 USD/năm, không mất thêm |

**Cứ 90 ngày bạn phải đưa lên bản mới**, nếu không người chơi iOS mất game.
Đây là cái bẫy lớn nhất của TestFlight — nhiều server private chết vì quên.

---

## Các bước làm

### Bước 1 — Đăng ký Apple Developer (làm trên máy nào cũng được)

1. Vào `developer.apple.com/programs`
2. Bấm **Enroll**, đăng nhập bằng Apple ID
3. Chọn **Individual** (cá nhân) — nhanh hơn Company nhiều
4. Nộp 99 USD bằng thẻ quốc tế

Apple duyệt mất **1-3 ngày**, có khi lâu hơn. Làm bước này **trước tiên**, đừng
đợi đến lúc có máy Mac mới đăng ký.

> Chọn **Individual**: tên nhà phát hành sẽ là tên thật của bạn.
> Chọn **Company**: hiện tên công ty, nhưng cần giấy phép kinh doanh và mã
> D-U-N-S — mất thêm vài tuần. Với server game thì Individual là đủ.

### Bước 2 — Tạo App ID

Trên `developer.apple.com` → **Certificates, Identifiers & Profiles** →
**Identifiers** → dấu **+**

Điền **Bundle ID** đúng chuỗi này, không sai một ký tự:

```
com.nrokame.game
```

> Chuỗi này phải khớp tuyệt đối với thứ nằm trong file `ServerListScreen.cs`
> của game. Sai là bản build không ký được, Xcode báo lỗi khó hiểu.

### Bước 3 — Chép gói sang MacBook

Chép `D:\NRO_iOS_GuiMac.zip` sang MacBook bằng USB, ổ cứng rời, hoặc AirDrop
nếu hai máy cùng mạng.

Giải nén ra. Bên trong có:
- `CLAUDE.md` — tài liệu bàn giao
- `NgocRongKame-iOS/` — dự án Xcode

### Bước 4 — Giao cho trợ lý trên MacBook

Mở trợ lý ở thư mục vừa giải nén rồi nói đại ý:

> Đọc CLAUDE.md rồi làm bản iOS đưa lên TestFlight giúp tôi.

Tài liệu đó đã có đủ: các bước, giá trị cố định, những chỗ hay hỏng, và quy tắc
bảo mật (không hỏi bạn mật khẩu Apple ID).

Ba bước dưới đây là để **bạn biết trợ lý sẽ làm gì** — không cần tự làm.

### Bước 5 — Mở bằng Xcode (trên máy Mac)

1. Cấp lại quyền chạy cho script: `chmod +x process_symbols.sh`
   *(nén trên Windows nên quyền này bị mất — thiếu bước này Xcode báo
   `Permission denied` giữa chừng)*
2. Mở `Unity-iPhone.xcodeproj` bằng Xcode
3. Chọn mục **Unity-iPhone** ở cột trái → thẻ **Signing & Capabilities**
4. Mục **Team**: chọn tài khoản Apple Developer vừa đăng ký

Ký tự động đã bật sẵn từ bên Windows, Bundle ID đã đúng `com.nrokame.game`.
Team là thứ **duy nhất** phải chọn tay.

Nếu báo lỗi đỏ ở đây, thường là:
- Chưa đăng nhập Apple ID trong Xcode → menu **Xcode > Settings > Accounts**
- Bundle ID chưa tạo ở Bước 2
- Tài khoản Developer chưa được Apple duyệt xong

### Bước 6 — Đóng gói và gửi lên Apple

1. Trên thanh trên cùng của Xcode, chọn thiết bị là **Any iOS Device**
   *(không phải máy giả lập — chọn máy giả lập là nút Archive bị mờ, bấm không được)*
2. Menu **Product → Archive**, đợi **20-40 phút**. Game dùng IL2CPP nên Xcode
   phải dịch một khối C++ rất lớn. Lâu là bình thường, đừng dừng giữa chừng.
3. Cửa sổ Organizer hiện ra → bấm **Distribute App**
4. Chọn **App Store Connect** → **Upload**

### Bước 7 — Bật TestFlight

1. Vào `appstoreconnect.apple.com`
2. **My Apps** → **+** → **New App**
   - Platform: iOS
   - Bundle ID: chọn `com.nrokame.game`
   - Tên: `Ngọc Rồng Kame`
3. Sang thẻ **TestFlight**, đợi bản vừa gửi xử lý xong (15-60 phút)
4. Điền **Test Information** — phần **What to Test** bắt buộc có chữ
5. Tạo nhóm ở **External Testing** → bấm **Public Link**

Apple duyệt bản đầu mất **1-2 ngày**. Duyệt xong bạn nhận được một địa chỉ dạng:

```
https://testflight.apple.com/join/xxxxxxxx
```

Đưa tôi địa chỉ đó, tôi gắn vào nút tải iOS trên web.

---

## Ba chỗ hay bị Apple từ chối

**1. Không có tài khoản để họ thử**
Phần **App Review Information** phải điền sẵn một tài khoản game (tên + mật khẩu)
để người của Apple đăng nhập vào xem. Không có là bị trả về ngay.

Tạo riêng một tài khoản cho việc này ở `nrokame.com/sv1/dangky`, đừng đưa tài
khoản quản trị.

**2. Máy chủ phải đang bật lúc họ duyệt**
Người duyệt mở game mà không kết nối được là họ đánh trượt. **Đừng tắt VPS
trong lúc chờ duyệt.**

**3. Bản quyền hình ảnh**
Game dùng nhân vật Dragon Ball. Apple **có thể** từ chối vì vi phạm bản quyền
của Toei Animation / Shueisha. Đây là rủi ro thật, tôi không thể hứa là qua được.

Nếu bị từ chối vì lý do này thì gần như không có cách sửa — trừ khi thay toàn bộ
hình ảnh nhân vật. Cân nhắc điều này **trước khi** nộp 99 USD.

---

## Tóm tắt chi phí và thời gian

| Việc | Tiền | Thời gian |
|---|---|---|
| Apple Developer | 99 USD/năm | chờ duyệt 1-3 ngày |
| Máy Mac | 0 (mượn) → 12 triệu (mua) | — |
| Cài iOS module + xuất mã | 0 | ~1 giờ |
| Xcode đóng gói và gửi | 0 | ~1 giờ |
| Apple duyệt TestFlight | 0 | 1-2 ngày |
| **Sau đó, mỗi 90 ngày** | 0 | ~1 giờ, **bắt buộc** |

---

## Lời khuyên thật lòng

Nếu bạn chưa chắc server sẽ có người chơi, **hãy mở cửa bằng bản PC và Android
trước**. Hai bản đó miễn phí, không cần xin phép ai, không có hạn 90 ngày.

Đợi khi có lượng người chơi ổn định và có người hỏi bản iOS thì hãy bỏ 99 USD.
Ở Việt Nam, người chơi game private trên iPhone chiếm phần nhỏ, và nhiều người
trong số đó vẫn cài được qua máy tính.

Tiền đó lúc đầu để dành trả VPS thì có ích hơn.
