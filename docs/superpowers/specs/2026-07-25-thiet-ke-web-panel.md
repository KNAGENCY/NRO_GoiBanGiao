# Thiết kế Web + Panel quản trị

**Mục tiêu:** thay bộ web PHP cũ (`2-WEB/`) bằng web hiện đại: trang người chơi
(đăng ký, đăng nhập, nạp thẻ, bảng xếp hạng) + panel quản trị. Chạy **song song**
với máy chủ game, dùng **chung cơ sở dữ liệu PostgreSQL**.

## Quyết định (đã chốt)

- **Vị trí:** `web/` trong cùng repo (monorepo). Vercel deploy được từ thư mục con.
- **Stack:** **Next.js (App Router) + TypeScript + Tailwind CSS.** CSDL qua thư
  viện `postgres` (postgres.js). Xác thực: PBKDF2 (thư viện chuẩn `node:crypto`)
  đối chiếu bảng `tai_khoan`, phiên qua **cookie ký (JWT `jose`)** httpOnly.
- **Đây là stack riêng** — ràng buộc "không thư viện ngoài" chỉ áp cho máy chủ Go.
  Web dùng npm bình thường (Next.js, postgres, jose, tailwind).
- **Env:** `DATABASE_URL` = chuỗi kết nối PostgreSQL (`postgres://postgres:nrodev123@127.0.0.1:5433/nro`).
  `AUTH_SECRET` = khoá ký JWT.

## Dùng chung CSDL với game — điểm mấu chốt

Web đối chiếu **cùng bảng `tai_khoan`** máy chủ game dùng (Kế hoạch 1):
`id, ten_dang_nhap, mat_khau_bam, bi_khoa, la_admin, tao_luc`. Ràng buộc tên
`^[A-Za-z0-9_]{4,20}$`. Mật khẩu định dạng **`pbkdf2$100000$<muối b64>$<băm b64>`**
(HMAC-SHA256, 100.000 vòng, muối 16 byte, băm 32 byte) — **phải khớp tuyệt đối**
bản Go/PHP để web và game dùng chung một mật khẩu.

> **Lợi ích kép:** web có trang **đăng ký** — mà game hiện CHƯA có đường đăng ký
> (Kế hoạch 1 chỉ làm đăng nhập). Nên web register vừa phục vụ web, vừa là cách
> tạo tài khoản cho game. Băm phải đúng format để đăng nhập game được luôn.

## Trang người chơi

| Đường dẫn | Nội dung |
|---|---|
| `/` | Trang chủ: giới thiệu game Naruto, nút tải, link đăng ký/đăng nhập |
| `/dang-ky` | Đăng ký: tên (4-20 `[A-Za-z0-9_]`), mật khẩu → tạo hàng `tai_khoan` với băm đúng format |
| `/dang-nhap` | Đăng nhập: đối chiếu PBKDF2, đặt cookie phiên |
| `/tai-khoan` | Trang cá nhân (sau đăng nhập): xem nhân vật, đổi mật khẩu |
| `/nap-the` | Nạp thẻ (khung trước, tích hợp cổng thanh toán sau) |
| `/bang-xep-hang` | Bảng xếp hạng theo sức mạnh nhân vật (`nhan_vat.suc_manh`) |

## Panel quản trị (`/admin`, chỉ `la_admin`)

| Đường dẫn | Nội dung |
|---|---|
| `/admin` | Bảng điều khiển: số liệu tổng quan (số tài khoản, nhân vật) |
| `/admin/tai-khoan` | Danh sách người chơi: tìm, **khoá/mở** (`bi_khoa`), xem nhân vật |
| `/admin/cau-hinh` | Chỉnh tỉ lệ EXP/rơi đồ/giá (bảng `cau_hinh_game` — tạo khi cần) |
| `/admin/nhat-ky` | Nhật ký quản trị (bảng `nhat_ky_quan_tri` — tạo khi cần) |

Tính năng nâng cao (giftcode, buff, giao dịch, sự kiện… như panel PHP cũ) làm dần
ở các kế hoạch web sau (W2, W3…).

## Lộ trình web

- **W1 — Nền tảng** (làm trước): scaffold Next.js + kết nối CSDL + tiện ích PBKDF2
  (băm + đối chiếu, khớp format Go) + trang chủ + đăng ký + đăng nhập + phiên cookie.
  Kiểm thử: đăng ký trên web → đăng nhập được **cả trên web lẫn game** (cùng băm).
- **W2 — Trang cá nhân + bảng xếp hạng + panel cơ bản** (danh sách tài khoản, khoá/mở).
- **W3 — Cấu hình game + nạp thẻ + tính năng nâng cao.**

## Bảo mật (nhắc, làm dần)

- Đăng nhập web cũng cần **chống dò mật khẩu** (giới hạn số lần) — cùng nợ với
  máy chủ game (Kế hoạch bảo mật B).
- Panel chỉ cho `la_admin`; kiểm ở server (middleware), không tin client.
- Không log mật khẩu. Cookie httpOnly + Secure (prod).
