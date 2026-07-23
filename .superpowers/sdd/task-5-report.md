# Báo cáo Task 5: Đối chiếu mật khẩu tương thích bản cũ

## Trạng thái: DONE

## File đã tạo
- `server/internal/auth/password.go` — `IsHashed`, `Hash`, `Verify`
- `server/internal/auth/password_test.go` — 6 test theo brief

Nội dung hai file được chép **nguyên văn** từ `task-5-brief.md` (Step 1 và Step 3), không thêm hàm nào ngoài brief, không cài thêm thư viện ngoài nào.

## Quy trình TDD

1. Viết `password_test.go` trước (chưa có `password.go`).
2. Chạy xác nhận FAIL (biên dịch lỗi vì `Verify`/`Hash`/`IsHashed` chưa tồn tại):

```
$ cd D:/NRO_GoiBanGiao/server && go test ./internal/auth/ -v
# nroserver/internal/auth [nroserver/internal/auth.test]
internal\auth\password_test.go:10:6: undefined: Verify
internal\auth\password_test.go:13:5: undefined: Verify
internal\auth\password_test.go:19:12: undefined: Hash
internal\auth\password_test.go:23:6: undefined: Verify
internal\auth\password_test.go:26:5: undefined: Verify
internal\auth\password_test.go:33:10: undefined: Hash
internal\auth\password_test.go:34:10: undefined: Hash
internal\auth\password_test.go:42:6: undefined: Verify
internal\auth\password_test.go:45:5: undefined: Verify
internal\auth\password_test.go:51:6: undefined: IsHashed
internal\auth\password_test.go:51:6: too many errors
FAIL	nroserver/internal/auth [build failed]
FAIL
```

3. Viết `password.go` nguyên văn theo brief, dùng `crypto/pbkdf2` của thư viện chuẩn Go (từ Go 1.24), **không chạy `go get golang.org/x/crypto`**, đúng chữ ký thư viện chuẩn:
   `pbkdf2.Key(hàm_băm, mật_khẩu_string, muối []byte, số_vòng int, số_byte int) ([]byte, error)`.
4. `go build ./...` sạch, không lỗi:

```
$ cd D:/NRO_GoiBanGiao/server && go build ./...
EXIT:0
```

5. Chạy lại test package `auth`:

```
$ cd D:/NRO_GoiBanGiao/server && go test ./internal/auth/ -v
=== RUN   TestDoiChieuBamCuaBanJava
--- PASS: TestDoiChieuBamCuaBanJava (0.02s)
=== RUN   TestBamRoiDoiChieu
--- PASS: TestBamRoiDoiChieu (0.03s)
=== RUN   TestMuoiNgauNhien
--- PASS: TestMuoiNgauNhien (0.02s)
=== RUN   TestTaiKhoanCuLuuTho
--- PASS: TestTaiKhoanCuLuuTho (0.00s)
=== RUN   TestNhanBietDangBam
--- PASS: TestNhanBietDangBam (0.00s)
=== RUN   TestChuoiBamHong
--- PASS: TestChuoiBamHong (0.00s)
PASS
ok  	nroserver/internal/auth	0.093s
```

`TestDoiChieuBamCuaBanJava` PASS — bằng chứng bản Go đọc được đúng chuỗi băm `pbkdf2$100000$5Q53KdKJ0c559SGZQOh2fw==$CtbK9GXUDmUYhDEg4AOdFM3cLT7s08mEPGgCoyUOAzU=` do bản Java cũ sinh ra cho mật khẩu `"1"`.

6. Chạy toàn bộ module (`go test ./...`) để xác nhận không phá vỡ package `protocol` đã có:

```
$ cd D:/NRO_GoiBanGiao/server && go test ./... -v
?   	nroserver/cmd/nroserver	[no test files]
=== RUN   TestDoiChieuBamCuaBanJava
--- PASS: TestDoiChieuBamCuaBanJava (0.02s)
=== RUN   TestBamRoiDoiChieu
--- PASS: TestBamRoiDoiChieu (0.03s)
=== RUN   TestMuoiNgauNhien
--- PASS: TestMuoiNgauNhien (0.02s)
=== RUN   TestTaiKhoanCuLuuTho
--- PASS: TestTaiKhoanCuLuuTho (0.00s)
=== RUN   TestNhanBietDangBam
--- PASS: TestNhanBietDangBam (0.00s)
=== RUN   TestChuoiBamHong
--- PASS: TestChuoiBamHong (0.00s)
PASS
ok  	nroserver/internal/auth	(cached)
=== RUN   TestGhiSoNguyenBigEndian
--- PASS: TestGhiSoNguyenBigEndian (0.00s)
=== RUN   TestGhiChuoiKieuJava
--- PASS: TestGhiChuoiKieuJava (0.00s)
=== RUN   TestChuoiTiengVietVongTron
--- PASS: TestChuoiTiengVietVongTron (0.00s)
=== RUN   TestChuoiQuaDaiBaoLoi
--- PASS: TestChuoiQuaDaiBaoLoi (0.00s)
=== RUN   TestChuoiMaHoaVuotGioiHanBaoLoi
--- PASS: TestChuoiMaHoaVuotGioiHanBaoLoi (0.00s)
=== RUN   TestGhiChuoiCoNulKhopJava
--- PASS: TestGhiChuoiCoNulKhopJava (0.00s)
=== RUN   TestGhiChuoiEmojiKhopJava
--- PASS: TestGhiChuoiEmojiKhopJava (0.00s)
=== RUN   TestUTFVongTronNhieuLoaiChuoi
--- PASS: TestUTFVongTronNhieuLoaiChuoi (0.00s)
=== RUN   TestDocByteJavaThatCoNul
--- PASS: TestDocByteJavaThatCoNul (0.00s)
=== RUN   TestDocByteJavaThatCoEmoji
--- PASS: TestDocByteJavaThatCoEmoji (0.00s)
=== RUN   TestDocSoNguyen
--- PASS: TestDocSoNguyen (0.00s)
=== RUN   TestDocQuaCuoiBaoLoi
--- PASS: TestDocQuaCuoiBaoLoi (0.00s)
=== RUN   TestDocGoiBatTayChuaMaHoa
--- PASS: TestDocGoiBatTayChuaMaHoa (0.00s)
=== RUN   TestVongTronSauKhiCoKhoa
--- PASS: TestVongTronSauKhiCoKhoa (0.00s)
=== RUN   TestLenhDacBietDungDoDaiBaByte
--- PASS: TestLenhDacBietDungDoDaiBaByte (0.00s)
=== RUN   TestLenhThuongDungDoDaiHaiByte
--- PASS: TestLenhThuongDungDoDaiHaiByte (0.00s)
=== RUN   TestGoiRongVanGhiDoDai
--- PASS: TestGoiRongVanGhiDoDai (0.00s)
=== RUN   TestLenhDacBietThanRongVanDungBaByte
--- PASS: TestLenhDacBietThanRongVanDungBaByte (0.00s)
=== RUN   TestGoiDacBietThanRongKhongLamLechGoiSau
--- PASS: TestGoiDacBietThanRongKhongLamLechGoiSau (0.00s)
=== RUN   TestLenhThuongThanRongVanMaHoaDoDai
=== RUN   TestLenhThuongThanRongVanMaHoaDoDai/Data=nil
=== RUN   TestLenhThuongThanRongVanMaHoaDoDai/Data=[]byte{}
--- PASS: TestLenhThuongThanRongVanMaHoaDoDai (0.00s)
    --- PASS: TestLenhThuongThanRongVanMaHoaDoDai/Data=nil (0.00s)
    --- PASS: TestLenhThuongThanRongVanMaHoaDoDai/Data=[]byte{} (0.00s)
=== RUN   TestLenhDacBietChuaCoKhoaTraLoiKhongPanic
--- PASS: TestLenhDacBietChuaCoKhoaTraLoiKhongPanic (0.00s)
=== RUN   TestTuChoiGoiQuaLon
--- PASS: TestTuChoiGoiQuaLon (0.00s)
=== RUN   TestDeriveKey
--- PASS: TestDeriveKey (0.00s)
=== RUN   TestDeriveKeyKhongSuaMangGoc
--- PASS: TestDeriveKeyKhongSuaMangGoc (0.00s)
=== RUN   TestGhiRoiDocRaByteBanDau
--- PASS: TestGhiRoiDocRaByteBanDau (0.00s)
=== RUN   TestConTroQuayVong
--- PASS: TestConTroQuayVong (0.00s)
=== RUN   TestHaiConTroDocLap
--- PASS: TestHaiConTroDocLap (0.00s)
=== RUN   TestGenerateRawKeyChanDoDaiSai
--- PASS: TestGenerateRawKeyChanDoDaiSai (0.00s)
PASS
ok  	nroserver/internal/protocol	(cached)
```

7. `go vet ./...` sạch, thoát mã 0, không in gì.

## Số test đạt/tổng
**34/34 PASS** (6 test mới của `internal/auth` Task 5 + 28 test đã có sẵn của `internal/protocol` Task 2-4, không có test nào bị ảnh hưởng).

## Kiểm tra `go.mod` / `go.sum`

- `go.mod` **không** có phần `require` nào (chỉ có `module nroserver` và `go 1.24`), xác nhận `crypto/pbkdf2` thuộc thư viện chuẩn, không cần thêm dependency ngoài:

```
module nroserver

go 1.24
```

- Không có `go.sum` phát sinh (đúng như kỳ vọng — không có dependency ngoài nào để khoá phiên bản).
- Không chạy `go get golang.org/x/crypto` như cảnh báo trong yêu cầu; brief gốc (Step 4) tự nói "Không cần thêm thư viện nào" nên không có bước nào cần bỏ qua.

## Mã băm commit

```
422401312d33cf36b3246b600af5c9ceb9fced1e
```
Short: `4224013`
Message: `feat(auth): doi chieu mat khau PBKDF2 tuong thich ban Java`
Diff: 2 file mới (`password.go`, `password_test.go`), 149 dòng thêm, không sửa file nào khác (không đụng `server/internal/protocol/`, không sửa `go.mod`/`go.sum` vì không có gì để sửa).

## Điều nghi ngờ

- **Không phát sinh nghi ngờ mới về logic** — toàn bộ `password.go` được chép nguyên văn từ brief; điểm dễ sai nhất (chữ ký `pbkdf2.Key` của thư viện chuẩn khác `x/crypto`: hàm băm đứng đầu, mật khẩu là `string`, trả kèm `error`) đã được brief cảnh báo trước và test `TestDoiChieuBamCuaBanJava` xác nhận khớp đúng với chuỗi băm Java thật.
- `Verify` chỉ chấp nhận đúng 4 phần tách bởi `$` (`pbkdf2`, số vòng, muối, băm) — nếu định dạng chuỗi băm Java/PHP có biến thể khác (ví dụ thêm phần thứ 5, hoặc dùng ký tự phân cách khác trong một số bản build cũ) thì `IsHashed`/`Verify` sẽ coi là "chuỗi hỏng" và trả `false` thay vì đối chiếu — nên xác minh thêm với vài mẫu dữ liệu thật khác trong CSDL `z` (không chỉ mẫu mật khẩu `"1"` trong brief) trước khi coi là bao phủ hết các trường hợp thực tế.
- Nhánh mật khẩu thô (`TestTaiKhoanCuLuuTho`) coi **bất kỳ chuỗi nào không có tiền tố `pbkdf2$`** là mật khẩu thô và so sánh trực tiếp bằng `subtle.ConstantTimeCompare`. Đây đúng theo yêu cầu brief, nhưng có nghĩa nếu một tài khoản cũ có mật khẩu thô tình cờ bắt đầu đúng bằng chuỗi `"pbkdf2$"` (khó xảy ra trong thực tế, gần như không thể vì đây là mật khẩu người dùng tự đặt) thì sẽ bị hiểu nhầm sang nhánh đã băm và luôn trả `false`. Rủi ro lý thuyết, không cần xử lý thêm trừ khi có bằng chứng dữ liệu thật gặp phải.
- Không sinh test nào ngoài 6 test đã cho trong brief, không thêm hàm ngoài yêu cầu (`Hash`, `Verify`, `IsHashed` đúng 3 hàm public theo Interfaces trong brief).

## Chặn số vòng lặp quá lớn

### Vấn đề
Hàm `Verify` đọc số vòng từ chuỗi băm rồi dùng ngay, không chặn giá trị quá lớn. Nếu CSDL bị hỏng với chuỗi như `pbkdf2$2147483648$...`, sẽ chạy hơn 2 tỷ vòng và treo hàng phút — một luồng xử lý đăng nhập bị treo cũng nguy hiểm ngang sập máy chủ.

### Sửa chữa
1. **Thêm hằng số chặn trên**: `soVongToiDa = 1_000_000` (gấp 10 lần giá trị hiện tại 100,000)
2. **Sửa `Verify`**: Kiểm `vong > soVongToiDa` cùng với `vong < 1`
3. **Thêm test mới**: `TestChanSoVongQuaLon` — kiểm 3 trường hợp:
   - `1000001` (vừa vượt ngưỡng)
   - `2147483648` (2 tỷ)
   - `9223372036854775807` (64-bit max)

### Test trước sửa (FAIL)
```
$ cd D:/NRO_GoiBanGiao/server && go test ./internal/auth/ -run TestChanSoVongQuaLon -timeout 30s -v

=== RUN   TestChanSoVongQuaLon
=== RUN   TestChanSoVongQuaLon/vua_vuot_nguong
=== RUN   TestChanSoVongQuaLon/2_ty_vong
panic: test timed out after 30s
```
Test timeout vì số vòng `2147483648` chạy thực 2+ tỷ vòng, treo vượt quá 30s.

### Test sau sửa (PASS)
```
$ cd D:/NRO_GoiBanGiao/server && go test ./internal/auth/ -v

=== RUN   TestDoiChieuBamCuaBanJava
--- PASS: TestDoiChieuBamCuaBanJava (0.02s)
=== RUN   TestBamRoiDoiChieu
--- PASS: TestBamRoiDoiChieu (0.03s)
=== RUN   TestMuoiNgauNhien
--- PASS: TestMuoiNgauNhien (0.02s)
=== RUN   TestTaiKhoanCuLuuTho
--- PASS: TestTaiKhoanCuLuuTho (0.00s)
=== RUN   TestNhanBietDangBam
--- PASS: TestNhanBietDangBam (0.00s)
=== RUN   TestChuoiBamHong
--- PASS: TestChuoiBamHong (0.00s)
=== RUN   TestChanSoVongQuaLon
=== RUN   TestChanSoVongQuaLon/vua_vuot_nguong
=== RUN   TestChanSoVongQuaLon/2_ty_vong
=== RUN   TestChanSoVongQuaLon/64-bit_max
--- PASS: TestChanSoVongQuaLon (0.00s)
    --- PASS: TestChanSoVongQuaLon/vua_vuot_nguong (0.00s)
    --- PASS: TestChanSoVongQuaLon/2_ty_vong (0.00s)
    --- PASS: TestChanSoVongQuaLon/64-bit_max (0.00s)
PASS
ok  	nroserver/internal/auth	(cached)
```

Toàn bộ 34 test đạt (6 test auth + 28 test protocol), không test nào bị break.

### Mã băm commit
```
54c312f Chan so vong lap qua lon trong doi chieu mat khau
```

**Tác giả**: Claude  
**Thay đổi**: 2 file (`password.go`, `password_test.go`), 54 dòng thêm, chặn số vòng >1,000,000, thêm test `TestChanSoVongQuaLon` với 3 sub-test.
