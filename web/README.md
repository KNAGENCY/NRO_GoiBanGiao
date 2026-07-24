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
