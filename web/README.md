# Web NRO (Naruto Online)

Trang web người chơi + panel quản trị cho game Naruto Online. Xây bằng
Next.js (App Router) + TypeScript + Tailwind CSS.

**Đây là W1 — chỉ nền tảng:** scaffold dự án, kết nối CSDL, băm mật khẩu
tương thích máy chủ game, trang chủ, đăng ký, đăng nhập, phiên đăng nhập
(cookie). Trang cá nhân, nạp thẻ, bảng xếp hạng, panel quản trị làm ở các
đợt sau (W2, W3).

## Dùng chung cơ sở dữ liệu với game

Web **không có CSDL riêng** — nó đọc/ghi thẳng vào cùng CSDL PostgreSQL
`nro` mà máy chủ game (`server/`, Go) dùng, cùng bảng `tai_khoan`. Vì vậy:

- Đăng ký trên web tạo tài khoản chơi được **luôn trong game** (và ngược
  lại, tài khoản tạo trong game đăng nhập được **luôn trên web**).
- Mật khẩu băm bằng PBKDF2-HMAC-SHA256, định dạng chuỗi lưu
  `pbkdf2$100000$<muối base64>$<băm base64>` — **khớp tuyệt đối** với
  `server/internal/auth/password.go`. Xem `lib/matkhau.ts`.
- CSDL phải được tạo và nạp lược đồ trước (xem `server/README.md`, mục
  "Tạo cơ sở dữ liệu và nạp lược đồ") — web KHÔNG tự tạo bảng.

## Cách chạy

### 1. Cài phụ thuộc

```bash
cd web
npm install
```

### 2. Tạo file môi trường

Sao chép `.env.local.example` thành `.env.local` rồi chỉnh nếu cần:

```bash
cp .env.local.example .env.local
```

| Biến           | Ý nghĩa                                                        |
|----------------|-----------------------------------------------------------------|
| `DATABASE_URL` | Chuỗi kết nối PostgreSQL — mặc định trỏ tới CSDL `nro` của game |
| `AUTH_SECRET`  | Khoá bí mật ký cookie phiên (JWT) — **đổi thành chuỗi ngẫu nhiên dài** khi lên production |

### 3. Chạy dev server

```bash
npm run dev
```

Mở http://localhost:3000.

### 4. Build kiểm tra biên dịch

```bash
npm run build
```

## Phiên đăng nhập — thu hồi được

Cookie phiên là JWT sống 7 ngày. Nếu chỉ kiểm chữ ký thì trong 7 ngày đó
**nội dung token là sự thật duy nhất**, và mọi quyết định của quản trị viên
đều không ăn: khoá tài khoản không đuổi được người đang online, hạ quyền
admin không gỡ được quyền, đăng xuất chỉ xoá cookie ở đúng một trình duyệt.

Cách chữa hiện tại: **mỗi lần đọc phiên đều hỏi cơ sở dữ liệu.**

| Thứ | Lấy từ đâu |
|---|---|
| `bi_khoa` | CSDL, mỗi lần đọc phiên → khoá là ra ngay |
| `la_admin` | CSDL, mỗi lần đọc phiên → hạ quyền là mất quyền ngay |
| `the_he_phien` | CSDL, so với số ghi trong token → đăng xuất/đổi mật khẩu giết mọi bản sao token |

**Giá phải trả: thêm đúng 1 truy vấn cho mỗi request *có* cookie phiên**
(tra cứu theo khoá chính, dưới 1 ms). Khách chưa đăng nhập không tốn thêm gì.
`docPhien()` được bọc trong React `cache()` nên layout + page + server action
trong cùng một request chỉ hỏi CSDL một lần.

Lý do chọn hướng này, và **những hướng đã cân nhắc rồi loại** (chỉ kiểm route
nhạy cảm · nhớ đệm ngắn hạn · dùng số thế hệ mà không hỏi CSDL · bảng `phien`
riêng), ghi đầy đủ ở đầu file [`lib/xacthucphien.ts`](lib/xacthucphien.ts).

Hai đánh đổi đã biết trước:

- **Đăng xuất ở một thiết bị là đăng xuất ở mọi thiết bị** (huỷ theo tài khoản
  chứ không theo từng phiên). Đổi lại được miễn phí chức năng "thoát khỏi mọi
  thiết bị".
- **Cookie phát hành trước bản vá này bị từ chối** (không có `the_he_phien` thì
  không thu hồi được) → người đang đăng nhập phải đăng nhập lại đúng một lần.

Cột `the_he_phien` do
[`server/internal/store/migrations/0003_the_he_phien.sql`](../server/internal/store/migrations/0003_the_he_phien.sql)
tạo ra — máy chủ Go tự áp migration khi khởi động, không phải sửa code Go.

### Kiểm thử

```bash
npm run test          # chạy cả hai bộ dưới đây
npm run test:phien    # riêng phần phiên — CẦN PostgreSQL đang bật
```

`test:phien` là **kiểm thử tích hợp, cần CSDL thật** và **cố tình chết ầm ĩ**
nếu không kết nối được — chứ không âm thầm bỏ qua (§9 của `BAT-DAU-TU-DAY.md`
ghi lại vụ 15 test đỏ nằm im nhiều ngày vì kiểu bỏ qua im lặng đó).

```powershell
D:\PostgreSQL\bin\pg_ctl.exe -D D:\PostgreSQL\data start
$env:NRO_TEST_DSN = "postgres://postgres:nrodev123@127.0.0.1:5432/nro_test"
```

Không đặt `NRO_TEST_DSN` thì nó dùng `DATABASE_URL` trong `.env.local`.
**Cổng PostgreSQL không cố định giữa các lần khởi động** — tự kiểm, đừng tin
con số ghi sẵn.

## Kiểm tương thích mật khẩu với máy chủ game

`lib/matkhau.ts` có bộ kiểm nội bộ (không cần Go):

```bash
npm run test:matkhau
```

Kiểm **hai chiều thật sự** với máy chủ Go (Node băm → Go đối chiếu, và
ngược lại) được thực hiện thủ công một lần khi dựng W1 — xem báo cáo
`../.superpowers/sdd/web-w1-report.md`. Không có chương trình Go tạm nào
được để lại trong `server/` (đã xoá sau khi kiểm xong).

## Cấu trúc thư mục

```
web/
├── app/
│   ├── page.tsx           # Trang chủ
│   ├── dang-ky/           # Đăng ký (page.tsx + server action actions.ts)
│   ├── dang-nhap/         # Đăng nhập (page.tsx + server action actions.ts)
│   ├── layout.tsx
│   └── globals.css        # Tailwind
├── lib/
│   ├── db.ts               # Kết nối PostgreSQL (postgres.js)
│   ├── matkhau.ts           # Băm/đối chiếu PBKDF2, khớp format máy chủ Go
│   ├── matkhau.test.ts      # Kiểm thử nội bộ cho matkhau.ts
│   └── phien.ts              # Cookie phiên JWT (jose)
├── .env.local.example
└── package.json
```
