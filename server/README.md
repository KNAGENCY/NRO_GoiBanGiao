# Máy chủ game (Go)

Máy chủ mô phỏng lại giao thức nhị phân của client Unity/J2ME cũ (xem
`internal/protocol`), có nghiệp vụ đăng nhập (xem `internal/game`), và lưu
tài khoản trên PostgreSQL (xem `internal/store`).

## Chạy lần đầu

### 1. Tạo cơ sở dữ liệu và nạp lược đồ

CSDL `nro` (dùng thật) và `nro_test` (dùng cho kiểm thử) cần được tạo trước
trên PostgreSQL. Sau đó nạp lược đồ:

```bash
D:\PostgreSQL\bin\psql.exe -h 127.0.0.1 -p 5433 -U postgres -d nro -f internal/store/migrations/0001_tai_khoan.sql
```

### 2. Tạo một tài khoản để thử

Tên đăng nhập chỉ được chứa chữ cái Latin, chữ số và gạch dưới, dài **4–20
ký tự** (`^[A-Za-z0-9_]{4,20}$`).

Mật khẩu là `1`, dùng chuỗi băm PBKDF2 sẵn có (khớp định dạng
`internal/auth`):

```sql
INSERT INTO tai_khoan (ten_dang_nhap, mat_khau_bam) VALUES
('thutest', 'pbkdf2$100000$5Q53KdKJ0c559SGZQOh2fw==$CtbK9GXUDmUYhDEg4AOdFM3cLT7s08mEPGgCoyUOAzU=');
```

Muốn tạo tài khoản khác với mật khẩu khác, dùng `auth.Hash("mật khẩu")`
trong một chương trình Go nhỏ, hoặc viết thẳng SQL với chuỗi băm tự sinh.

### 3. Cấu hình và chạy máy chủ

Biến môi trường:

| Biến           | Bắt buộc | Ý nghĩa                                          | Mặc định  |
|----------------|----------|---------------------------------------------------|-----------|
| `NRO_DSN`      | Có       | Chuỗi kết nối PostgreSQL                           | (không có)|
| `NRO_DIA_CHI`  | Không    | Địa chỉ TCP để lắng nghe                           | `:14445`  |

PowerShell:

```powershell
$env:NRO_DSN = "postgres://postgres:nrodev123@127.0.0.1:5433/nro"
$env:NRO_DIA_CHI = ":14445"
go run ./cmd/nroserver
```

cmd.exe:

```bat
set NRO_DSN=postgres://postgres:nrodev123@127.0.0.1:5433/nro
set NRO_DIA_CHI=:14445
go run ./cmd/nroserver
```

Khi chạy đúng, máy chủ in ra (theo đúng thứ tự):

```
level=INFO msg="dang ket noi CSDL..."
level=INFO msg="da noi CSDL thanh cong."
level=INFO msg="may chu game dang lang nghe" dia_chi=:14445
level=INFO msg="may chu san sang. Nhan Ctrl+C de dung may chu an toan."
```

Nếu thiếu `NRO_DSN` hoặc không nối được CSDL, máy chủ in lỗi tiếng Việt giải
thích rõ đang thiếu gì / nên kiểm tra gì, rồi thoát với mã khác 0 — **không**
mở cổng lắng nghe trong tình trạng nửa vời.

### 4. Tắt máy chủ

Nhấn `Ctrl+C` (hoặc gửi tín hiệu `SIGTERM`). Máy chủ sẽ đóng kết nối CSDL rồi
mới thoát, in ra:

```
level=INFO msg="da nhan tin hieu dung, dang dong may chu..." tin_hieu=interrupt
level=INFO msg="da dong CSDL. Thoat may chu."
```

## Kiểm chứng bằng client giả (`thuclient`)

Không cần mở game Unity thật để kiểm tra máy chủ có đăng nhập đúng không —
dùng `cmd/thuclient`, một chương trình độc lập làm đúng những gì client thật
làm khi kết nối (bắt tay trao khóa, gửi gói đăng nhập) và in kết quả bằng
tiếng Việt dễ hiểu.

Chạy máy chủ ở một cửa sổ, rồi ở cửa sổ khác:

```powershell
go run ./cmd/thuclient -dia-chi 127.0.0.1:14445 -ten thutest -mat-khau 1
```

Tham số dòng lệnh:

| Tham số      | Bắt buộc | Ý nghĩa                          | Mặc định            |
|--------------|----------|-----------------------------------|----------------------|
| `-dia-chi`   | Không    | Địa chỉ máy chủ (`host:port`)     | `127.0.0.1:14445`    |
| `-ten`       | Có       | Tên đăng nhập                     | (không có)           |
| `-mat-khau`  | Có       | Mật khẩu                          | (không có)           |

Mật khẩu đúng → thấy dòng `DANG NHAP THANH CONG!`. Mật khẩu sai hoặc tài
khoản không tồn tại → thấy `THONG BAO TU MAY CHU: Sai tài khoản hoặc mật
khẩu.` rồi `DANG NHAP THAT BAI`, chương trình tự thoát ngay (không treo).

## Chạy toàn bộ kiểm thử

```bash
go test ./... -v
```

Kiểm thử trong `internal/store` cần biến `NRO_TEST_DSN` (trỏ vào một CSDL
PostgreSQL riêng để kiểm thử, ví dụ `nro_test`); không đặt thì các kiểm thử
đó tự động bị bỏ qua (SKIP), các gói còn lại vẫn chạy và không phụ thuộc
CSDL.

```powershell
$env:NRO_TEST_DSN = "postgres://postgres:nrodev123@127.0.0.1:5433/nro_test"
go test ./... -v
```
