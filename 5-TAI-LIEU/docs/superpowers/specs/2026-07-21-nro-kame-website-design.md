# Thiết kế website Ngọc Rồng Kame

**Ngày:** 2026-07-21 · **Trạng thái:** đã được chủ dự án duyệt

## Bối cảnh

Server NRO đã chạy được ở local, đã đổi thương hiệu sang Ngọc Rồng Kame, đã vá lỗ hổng mật khẩu. Chủ dự án cung cấp bộ `nro-kame-complete-ui-kit.zip` gồm 7 mockup, `ui-kit.css` chất lượng production, `SITEMAP.md`, và 3 component SVG.

Web hiện tại là **một trang đơn** (`sv1/trangchu.php`, 89KB) cộng 2 trang bảng xếp hạng. Yêu cầu: **làm lại layout hoàn toàn** theo bộ kit, **giữ nguyên toàn bộ chức năng**.

## Quyết định đã chốt

| Vấn đề | Quyết định |
|---|---|
| Mockup và `ui-kit.css` mâu thuẫn về font/màu | **Theo `ui-kit.css`**. Mockup ghi Bebas Neue + Inter; Bebas Neue **không có dấu tiếng Việt** nên chữ sẽ vỡ. CSS ghi Roboto Condensed + Be Vietnam Pro — cả hai có dấu đầy đủ. |
| 8 nhân vật hư cấu (Son Kai, Lina...) với hệ chỉ số không tồn tại | **Thay bằng top người chơi thật** từ bảng `player`, ảnh đại diện thật qua `head_avatar` → `icon_asset.php` |
| Boss "Vua Hắc Ám" hư cấu | **Thay bằng 98 boss thật** từ `boss_status.json`, ảnh thật qua `avatar_id` |
| Số liệu 128K Discord / 215K Facebook / 350.842 tài khoản | **Bỏ hết số**. Chỉ ghi tên nền tảng. Số online/tài khoản lấy thật từ DB. |
| Discord | **Bỏ hoàn toàn** — chủ dự án không dùng |
| Hệ thống tin tức | Tạo bảng `news` + **gieo 30 bài mẫu**. Không làm trang quản trị — chủ dự án nói "không cần thiết lắm" |
| Form gửi ticket hỗ trợ | **Bỏ.** Thay bằng FAQ + nút vào nhóm Zalo + mục bình luận có sẵn |
| Menu chính | **Trang chủ · Vũ trụ · Tin tức · Bảng xếp hạng · Hỗ trợ** |
| Ảnh nhân vật Dragon Ball nhận ra được trong mockup | Không dùng. Ảnh mới phải là nhân vật gốc hoặc phong cảnh không người |

## Chức năng bắt buộc giữ — 13 mục

Đăng nhập · Đăng ký · Quên mật khẩu · Đổi mật khẩu · Tải game (3 nền tảng, luồng 2 bước) · Theo dõi boss trực tiếp · 2 bảng xếp hạng · Bình luận cộng đồng · Trạng thái server · Điều khoản · Thông báo admin · Liên kết cộng đồng · Lối vào trang quản trị.

Bổ sung: dòng cảnh báo *"Dành cho người chơi trên 12 tuổi. Chơi quá 180 phút mỗi ngày sẽ hại sức khỏe."*

## Kiến trúc

```
C:\xampp\htdocs\
├── index.php               GIỮ NGUYÊN (trang nguỵ trang EduVerse)
├── icon_asset.php          GIỮ NGUYÊN (đang phục vụ 20.924 icon game)
├── img/brand/              ảnh thương hiệu đã lắp
└── sv1/
    ├── kame/                       THƯ MỤC MỚI
    │   ├── ui-kit.css              từ bộ kit, KHÔNG SỬA
    │   ├── theme.css               phần bổ sung
    │   ├── layout.php              header + footer + <head> dùng chung
    │   ├── data.php                tầng truy vấn dữ liệu
    │   └── assets/                 ảnh bổ sung (chủ dự án gửi sau)
    ├── trangchu.php        VIẾT LẠI
    ├── vu-tru.php          MỚI
    ├── tin-tuc.php         MỚI
    ├── tai-game.php        MỚI
    ├── ho-tro.php          MỚI
    ├── bxh01.php bxh02.php GIỮ LOGIC, khoác giao diện mới
    ├── config.php          GIỮ NGUYÊN
    ├── boss_status.php     GIỮ NGUYÊN
    └── dangky/dangnhap/doimatkhau/quenmatkhau/download.php
                            ❌ TUYỆT ĐỐI KHÔNG SỬA
```

**Vì sao không đụng 5 file tài khoản:** chúng vừa được vá lỗ hổng mật khẩu trong cùng ngày. Sửa lại là rước rủi ro làm hỏng đăng nhập. Chỉ thay giao diện *gọi* chúng, không thay bản thân chúng.

## Luồng dữ liệu

```
Trang  ──┬─ kame/data.php ─┬─ boss_status.json        (server Java ghi mỗi 1 giây)
         │                 ├─ panel_server_status.json (online thật)
         │                 ├─ bảng player              (top sức mạnh / tài sản)
         │                 ├─ bảng account             (tổng tài khoản)
         │                 ├─ bảng news                (30 bài gieo sẵn)
         │                 └─ bảng site_comments       (bình luận)
         └─ icon_asset.php ─ ảnh gốc trong game
```

`kame/data.php` là **tầng truy vấn duy nhất**. Trang không tự viết SQL — mọi truy vấn đi qua đây để dễ sửa và dễ kiểm.

## Xử lý lỗi

| Tình huống | Hành vi |
|---|---|
| Server Java tắt | Khối boss hiện "Máy chủ đang bảo trì", phần còn lại của trang vẫn chạy |
| `boss_status.json` cũ hơn 60 giây | Hiện nhãn "dữ liệu chưa cập nhật", vẫn hiển thị nội dung cũ |
| Chưa có người chơi | Bảng xếp hạng hiện "Chưa có dữ liệu — hãy là người đầu tiên" |
| Ảnh icon thiếu | Rơi về ảnh mặc định, không vỡ bố cục |
| Bảng `news` chưa tồn tại | Trang tin tức hiện thông báo trống, không lỗi trắng trang |

## Ảnh — trạng thái

**Đã có:** logo ngang/vuông, hero 1920×1080, band-boss, band-congdong, texture giấy cũ, og, fbcover, favicon 6 cỡ.

**Chủ dự án đang tạo, gửi sau — 8 tấm:**
`kame_web_the-gioi-{traidat,namec,xayda}_800x500_v1.jpg` · `kame_web_map-{traidat,namec,xayda}_1400x800_v1.jpg` · `kame_web_tin-mac-dinh_800x450_v1.jpg` · `kame_web_band-taigame_1920x600_v1.jpg`

**Cho tới khi có ảnh:** dùng ô giữ chỗ bằng CSS (nền màu + hoạ tiết giấy cũ + nhãn tên file). Tên file đã chốt nên khi ảnh về chỉ cần chép vào `sv1/kame/assets/`, không phải sửa code.

## Kiểm chứng

1. `php -l` sạch trên mọi file mới
2. Mỗi trang trả HTTP 200
3. Đủ 13 chức năng, đặc biệt: đăng nhập/đăng ký vẫn gọi đúng `/sv1/dangnhap`, `/sv1/dangky`
4. Không còn chuỗi `discord`, không còn số liệu bịa (`128K`, `215K`, `78K`, `350`, `5.236`)
5. Boss và bảng xếp hạng hiển thị dữ liệu thật từ DB
6. Responsive ở 3 mức 1180 / 900 / 600 px

## Ngoài phạm vi

Trang quản trị tin tức · hệ thống ticket · newsletter · Discord · giải đấu · vá client · đổi mật khẩu database (`root` vẫn rỗng — làm khi lên VPS).
