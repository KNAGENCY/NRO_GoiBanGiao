# Chuẩn giao diện Kame-Admin — guideline cho agent restyle panel

Ngày: 2026-07-22 · Đã duyệt hướng Kiểu 2 (nền vân giấy, header nâu, viền cam).

## Mục tiêu

Restyle mọi trang trong `C:\xampp\htdocs\sv1\admin\` để dùng chung hệ thiết kế
Kame, thay cho CSS tự viết mỗi trang một kiểu. **Chỉ đổi giao diện — không đổi
một dòng logic nào.**

## Ba file nền (ĐÃ CÓ, không sửa)

- `sv1/kame/ui-kit.css` — class chung: `.kame-button`, `.kame-button--secondary`,
  `.kame-button--ghost`, `.kame-button--danger`, `.kame-card`, `.kame-badge`,
  `.kame-badge--success/hot/info`, `.kame-field`, `.kame-label`, `.kame-input`,
  `.kame-select`, `.kame-textarea`, `.kame-help`, `.kame-error`, `.kame-alert`,
  `.kame-alert--success`, `.kame-muted`, `.kame-heading`.
- `sv1/kame/theme.css` — biến `--kame-*` (màu, khoảng cách, bo góc) + nền vân giấy.
- `sv1/admin/kame-admin.css` — riêng admin: `.kadmin`, `.kadmin-bar`,
  `.kadmin-bar__title`, `.kadmin-bar__actions`, `.kadmin-body`, `.kadmin-card`,
  `.kadmin-card__head`, `.kadmin-toolbar`, `.kadmin-toolbar__grow`,
  `.kadmin-table-wrap`, `.kadmin-table` (td `.num`/`.act`), `.kadmin-stats`,
  `.kadmin-stat`, `.kadmin-stat__label`, `.kadmin-stat__value`, `.kadmin-back`,
  `.kadmin-icon`.

## Khung chung (ĐÃ CÓ) — `sv1/admin/admin_layout.php`

Hàm dùng trong mỗi trang:
- `admin_head('Tên trang')` — mở `<html><head>` nạp 3 CSS + Tabler icon, mở `.kadmin`.
- `admin_bar('Tiêu đề', 'ten-icon-tabler')` — thanh header nâu + nút "Quay lại Panel".
- `admin_foot()` — đóng khung.

## Trang mẫu vàng (ĐÃ LÀM) — chép y khuôn từ đây

`sv1/admin/maintenance.php` — đọc để thấy cấu trúc chuẩn:
phần logic PHP giữ nguyên → `admin_head()` → `admin_bar()` →
`<div class="kadmin-body"> ...nội dung dùng .kadmin-card... </div>` → `admin_foot()`.

## Công thức restyle một trang

1. **Sao lưu**: chép file sang `C:\Users\Admin\Music\banluu-web-20260722\<ten>.bak`.
2. **Thêm** `require_once __DIR__ . '/admin_layout.php';` cạnh các require đầu file.
3. **Giữ nguyên 100%**: `session_start`, kiểm đăng nhập, kiểm `is_admin`, mọi câu
   truy vấn `$pdo`, xử lý `$_POST`, CSRF, biến dữ liệu, và **mọi thẻ `<script>`**.
4. **Thay khối `<!DOCTYPE>...</html>`**: bỏ `<head><style>` cũ, thay bằng
   `admin_head()` + `admin_bar()` + `<div class="kadmin-body">` + `admin_foot()`.
5. **Chuyển markup nội dung** sang class Kame: bảng → `.kadmin-table` trong
   `.kadmin-table-wrap`; form input → `.kame-input`/`.kame-select`; nút → `.kame-button`
   (nút xoá/nguy hiểm → `--danger`, nút phụ → `--secondary`); thẻ khối → `.kadmin-card`
   + `.kadmin-card__head`; thông báo → `.kame-alert`; số liệu tổng → `.kadmin-stat`.
6. **Giữ nguyên mọi `name=`, `id=`, `action=`, `onclick=`, `data-*`** — JS và form
   phụ thuộc vào chúng. Chỉ đổi `class=` và bọc thẻ, không đổi thuộc tính chức năng.

## Quy tắc CỨNG — vi phạm là hỏng

- ❌ KHÔNG sửa PHP logic, truy vấn, xử lý POST, CSRF, JS chức năng.
- ❌ KHÔNG hardcode màu (`#fff`, `red`, `rgb(...)`). Chỉ dùng biến `--kame-*` và class.
- ❌ KHÔNG tự chế class CSS mới. Thiếu thành phần thì báo lại, KHÔNG tự thêm vào
  `kame-admin.css` (để tránh mỗi agent thêm một kiểu). 
- ❌ KHÔNG đổi `name/id/action/onclick/data-*`.
- ✅ Chỉ đổi `class=`, cấu trúc thẻ hiển thị, và phần `<head>`.

## Kiểm chứng bắt buộc trước khi báo xong

1. `& 'C:\xampp\php\php.exe' -l <file>` → "No syntax errors".
2. Đếm: mọi `name=`, `id=`, `<script>`, câu `$pdo->` trong bản mới phải bằng bản cũ
   (grep đối chiếu) — chứng minh không mất chức năng.
3. Ghi rõ: đã đổi những gì (class/markup), giữ nguyên những gì (logic/JS).

## Danh sách trang (16 trang, maintenance.php đã xong)

Đơn giản: `buff_vnd.php`, `itemshop.php`, `gioihan.php`, `noitai.php`.
Vừa: `account.php`, `giaodich.php`, `taocode.php`, `phucloi.php`, `item.php`,
`buff_item.php`, `map.php`.
Lớn/phức tạp (làm sau, cẩn thận): `giftcode.php`, `ngoc.php`, `vang.php`,
`player.php`, `panel.php` (trang chính, để cuối).

## Ghi chú

- `game_bridge.php`, `logout.php`, `reset.php`, `panel_check.php`, `index.php`,
  `ajax_create_giftcode.php` KHÔNG phải trang hiển thị — bỏ qua.
- Làm trên htdocs máy nhà, xong xuôi cả loạt mới đồng bộ VPS.
- Đồng hồ đo trong `panel.php` là tạm — khi restyle `panel.php` thì giữ lại.
