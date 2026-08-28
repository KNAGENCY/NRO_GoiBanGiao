# KH6b Boss — kế hoạch thi công (BẢN 2, sửa theo 7 mục CHẶN)

> Ngày 22/08/2026. Bản 1 do workflow khảo sát sinh ra; bản này đã áp 7 vấn đề CHẶN
> do 3 lens phản biện tìm ra (đều kiểm chứng bằng code thật).

## Các thay đổi so với bản 1

### CHAN-1+6 (than goi -60 la danh sach id + cdir)
BO han task 'sua parser -60' khoi ban 1. Bang chung client that (Service.cs:1256-1265) da duoc ghi vao dac ta (specs/2026-07-26-giao-thuc-chien-dau.md dong 189 + ghi chu 22/08/2026): server Go doc 1 int nen VO TINH DUNG voi 1 muc tieu. login.go GIU NGUYEN. Task 8 (dinh tuyen -60 len boss) chi doc gia tri id da parse san, khong dung toi tang parser; ghi ro trong brief Task 8 la 'cam sua login.go'.

### CHAN-2 (announce goi tu trong actor -> tu khoa)
Task cu 18 'goi SoKhuVuc.PhatToanServer tu trong vongLap' bi XOA. Thay bang Task 6 MOI: bo phat thanh rieng (loathongbao.go) theo dung khuon ThoGhi (thogi.go) — actor CHI goi loa.Xep(YeuCauLoa{...}) (chi khoa mutex vai ns, KHONG BAO GIO chan), mot goroutine rieng rut hang doi, ap nhip ~1,5s roi moi goi SoKhuVuc.PhatToanServer. Ly do da kiem chung: guiDieuKhien (khuvuc.go:470) cho toi thoiGianChoControl=100ms tren chinh hop thu ma vongLap dang ban -> goi tu trong actor la tu khoa chac chan. Task 12 (announce boss lon) phu thuoc Task 6, khong duoc goi thang.

### CHAN-3 (lich boss lam do test tich hop)
MoSoKhuVuc() (sokhuvuc.go) GIU NGUYEN — cam tuyet doi nhet khoi dong lich boss vao day. Tach thanh Task 13 rieng: bo dieu phoi lich (lichboss.go) dung khuon BoLuuDinhKy (luudinhky.go) — Mo... roi BatDau() hai buoc, dong ho tiem duoc, CHI cmd/nroserver/main.go bat (Task 14). Task 13 phai co MOT test chot: MoSoKhuVuc() don thuan KHONG sinh boss nao, 5 khu vuc van im lang (helper kiem im lang chi tha -11/-10 phai con xanh nguyen).

### CHAN-4 (KiemDanhHopLe nhan *nguoiTrongKhuVuc, boss khong truyen vao duoc)
Them Task 5 DON NEN, lam TRUOC moi task boss cham khuvuc.go: doi chu ky thanh KiemDanhHopLe(benDanh BenDanh, mucTieu ViTriCoMau, bayGio) voi BenDanh{X, Y, ChakraHienTai, LanDanhCuoi, Chieu *CapChieu, VoHanChakra bool}. VoHanChakra=true bo qua kiem chakra (bay #7 cua thiet ke: dung so ma thuat 31_07_2002). Thuan refactor: hai cho goi trong khuvuc.go (case suKienDanhQuai, case suKienDanhNguoiChoi) dung adapter tu *nguoiTrongKhuVuc; toan bo chotkiem_test.go phai xanh KHONG doi ky vong nao.

### CHAN-5 (goiXuatHienDayDu thieu truong, HinhDangHien khong phu)
Task 4 chot CHU KY DAY DU, khong lay tu store.NhanVat nua: goiXuatHienDayDu(id int64, ten string, cap byte, hp, hpToiDa int32, x, y int16, hd HinhDangHien) — HinhDangHien{Lang int16, Outfit [6]int16 (head/body/leg/flagBag/aura/effFront), TypePk byte}. goiXuatHien cu tro thanh vo boc mong cho nguoi choi (outfit 0, typePk 0) de KH5 khong vo. Kem do: PhaBoss THEM truong Cap (kiem 1..255, tu choi 0 va >255) trong Task 3, vi truong 3 cua goi -5 la byte level.

### CHAN-7 (hai co che sinh boss chong nhau)
Chia doi trach nhiem ro rang: 'khi nao + o map nao' NAM NGOAI actor (Task 13, lichboss.go) va gui lenh sinh qua hop thu DIEU KHIEN cua khu vuc; may trang thai trong actor (Task 9) CHI lo chuyen pha SAU KHI nhan lenh do — no khong duoc tu quyet dinh spawn, khong doc dong ho lich, khong biet map nao. Task 9 bo han nhanh 'tu hoi sinh theo GiayNghi trong nhipTim'; GiayNghi tro thanh du lieu ma bo dieu phoi doc.

### Thu tu thi cong (chong dung khuvuc.go)
Xep lai: moi task SUA khuvuc.go (7, 8, 9, 10, 11, 12) chay TUAN TU tren cung nhanh, KHONG song song — khuvuc.go la 63KB, hai agent cung vao chac chan dung. Ba task file MOI hoan toan (1, 2, 3) tach ra lam song song trong worktree rieng ngay tu dau.


## Danh sách task (thứ tự thi công)

| # | Task | File | Phụ thuộc | Song song? |
|---|---|---|---|---|
| 1 | GATE: kiem client co chap nhan mot id da -6 quay lai bang -5 khong (Unity MCP doc Controller.cs/findCharInMap) — ghi ket qua thanh van ban | `D:/NRO_GoiBanGiao/docs/superpowers/research/2026-08-22-kh6b-kiem-client-bien-hinh.md` | khong — chay song song ngay tu dau | ✅ file mới |
| 2 | Goi 93 (thong bao chay chu toan server) + goi 44 (chat trong khu vuc) — hai builder thuan + test byte-chinh-xac | `D:/NRO_GoiBanGiao/server/internal/game/goithongbao.go` | khong — chay song song ngay tu dau | ✅ file mới |
| 3 | Mo hinh du lieu boss: Boss, PhaBoss (THEM truong Cap, kiem 1..255), MocNoiGian, KieuXuatHien; bo dem id am duy nhat toan cuc; kep HP ve [0, 2_147_483_647] | `D:/NRO_GoiBanGiao/server/internal/game/boss.go` | khong — chay song song ngay tu dau | ✅ file mới |
| 4 | CHAN-5: goiXuatHienDayDu(id, ten, cap, hp, hpToiDa, x, y, HinhDangHien) + kieu HinhDangHien{Lang, Outfit[6], TypePk}; goiXuatHien cu thanh vo boc mong cho nguoi choi | `D:/NRO_GoiBanGiao/server/internal/game/xuathien.go` | khong (nhung la file SUA — lam tuan tu, khong song song voi task khac cham xuathien.go) | — tuần tự |
| 5 | CHAN-4: tach KiemDanhHopLe thanh ham thuan hai ben — BenDanh{X,Y,ChakraHienTai,LanDanhCuoi,Chieu,VoHanChakra} + ViTriCoMau; sua 2 cho goi trong khuvuc.go bang adapter, test cu xanh nguyen | `D:/NRO_GoiBanGiao/server/internal/game/chotkiem.go (+ 2 cho goi trong khuvuc.go)` | khong — nhung PHAI xong TRUOC task 7-12 vi cham khuvuc.go | — tuần tự |
| 6 | CHAN-2: bo phat thanh ngoai actor (khuon ThoGhi.Xep, KHONG BAO GIO chan) + SoKhuVuc.PhatToanServer + nhip ~1,5s; actor chi Xep yeu cau, goroutine rieng moi phat | `D:/NRO_GoiBanGiao/server/internal/game/loathongbao.go (MOI) + sokhuvuc.go, khuvuc.go (SUA: them truong loa + suKienPhatGoi)` | task 2 (can goiThongBao93) | — tuần tự |
| 7 | Boss vao khu vuc: them k.boss canh k.quai, phat -5 outfit that khi boss hien, sua suKienVao lap CA k.boss cho nguoi moi vao | `D:/NRO_GoiBanGiao/server/internal/game/khuvuc.go` | task 3, 4 | — tuần tự |
| 8 | Dinh tuyen -60: id am + co trong k.boss -> ap sat thuong len boss, cong satThuongTheoNguoi. CAM sua parser login.go (CHAN-1 da giai quyet) | `D:/NRO_GoiBanGiao/server/internal/game/khuvuc.go` | task 5, 7 | — tuần tự |
| 9 | CHAN-7: may trang thai vong doi trong actor (NGHI->HIEN->CHAO->HOAT_DONG->TROI->RA_DI) CHI chay sau khi nhan lenh sinh qua hop thu dieu khien — khong tu quyet 'khi nao / map nao' | `D:/NRO_GoiBanGiao/server/internal/game/khuvuc.go` | task 7 | — tuần tự |
| 10 | AI boss trong nhipTim 100ms: chon muc tieu + danh bang chieu KH6 (qua BenDanh) + di chuyen khi ngoai tam + thoai danh bang goi 44 | `D:/NRO_GoiBanGiao/server/internal/game/khuvuc.go` | task 2, 5, 9 | — tuần tự |
| 11 | Bien hinh qua pha (GIU NGUYEN id) + noi gian theo %mau (MocNoiGian doi he so dame/def/nhip + chat) | `D:/NRO_GoiBanGiao/server/internal/game/khuvuc.go + boss.go` | task 1 (GATE: neu client khong chap nhan id quay lai thi doi sang id moi moi pha TRUOC khi lam), task 9, 10 | — tuần tự |
| 12 | Chia thuong theo dong gop (diem san boss) khi boss chet + announce 93 boss lon luc hien/luc bi ha — phat QUA hang doi loa, khong goi thang | `D:/NRO_GoiBanGiao/server/internal/game/khuvuc.go` | task 6, 9, 11 | — tuần tự |
| 13 | CHAN-3+7: bo dieu phoi lich boss NGOAI actor (khuon BoLuuDinhKy: Mo... + BatDau() hai buoc, dong ho tiem duoc) gui lenh sinh qua hop thu dieu khien; test chot MoSoKhuVuc() KHONG sinh boss nao | `D:/NRO_GoiBanGiao/server/internal/game/lichboss.go` | task 9 (can kieu lenh sinh) | — tuần tự |
| 14 | Noi day that o main.go: bat bo phat thanh + bo dieu phoi lich boss; tat may chu dung sach ca hai theo dung thu tu (LuuTatCa -> DongVaChoXong -> DungTatCa) | `D:/NRO_GoiBanGiao/server/cmd/nroserver/main.go (+ internal/game/tatmaychu.go)` | task 6, 13 | — tuần tự |
| 15 | Boss mau (1-2 con, du lieu thuan Go) + smoke bang client Unity that: thay boss hien, danh duoc, boss danh lai, bien hinh, announce chay chu | `D:/NRO_GoiBanGiao/server/internal/game/noidungboss.go` | task 14 | — tuần tự |

## Ghi chú
PHAM VI: chi ban 1 (KH6b-1, §6 thiet ke). Da loai khoi danh sach: spawn theo trigger, hanh dong dac biet (bom/hap thu/trieu hoi), roi do ra dat (cho KH7), panel boss_status.json, menu -96/-118, chuyen du lieu sang YAML (KH8).

SONG SONG: chi 3 task dau (1, 2, 3) la file MOI hoan toan va khong phu thuoc gi — tach worktree rieng chay cung luc. Tu task 4 tro di lam TUAN TU: task 5-12 deu cham khuvuc.go (63KB, 1500 dong) nen hai agent cung vao la dung chac. Task 6 va 13 tuy la file moi nhung phu thuoc dau ra cua task truoc nen khong xep doc_lap.

GATE THAT SU: task 1 phai co ket qua TRUOC task 11. Neu client tu choi mot id da nhan -6 quay lai bang -5, doi thiet ke bien hinh sang 'moi pha mot id moi (-6 pha cu, -5 pha moi)' va sua brief task 11 truoc khi code — dung code roi moi kiem.

BA KHUON DA CO SAN TRONG REPO, BAM THEO CHU KHONG SANG TAO LAI:
- thogi.go ThoGhi.Xep -> khuon hang doi khong-bao-gio-chan cho task 6 (loa).
- luudinhky.go BoLuuDinhKy -> khuon goroutine dieu phoi ngoai actor, Mo/BatDau hai buoc + dongHo tiem duoc, cho task 13 (lich boss).
- khuvuc.go hoiSinhQuai/hoiMauQuai/aiQuai tach khoi nhipTim -> khuon tach ham thuan de test moc thoi gian ma khong sleep, cho task 9/10.

BAY DA KIEM CHUNG BANG CODE THAT (dua vao brief tung task):
- khuvuc.go:470 guiDieuKhien cho toi 100ms tren hop thu ma vongLap dang ban -> moi loi goi 'dong bo' tu trong actor deu tu khoa.
- chotkiem.go hien doc thang nguoiDanh.nv.* va nguoiDanh.chieuDangChon -> boss (kieu rieng, khong co store.NhanVat) khong the truyen vao.
- xuathien.go dang ghi cung 0 cho body/leg/flagBag/aura/effFront va lay het tu store.NhanVat -> boss khong mac duoc do that.
- specs/2026-07-26-giao-thuc-chien-dau.md dong 189 da dinh chinh -60 (danh sach id + cdir) tu 22/08/2026 — dung mo lai chuyen nay.

VIEC KEM: sau khi chot danh sach nay, ghi thanh file ke hoach D:/NRO_GoiBanGiao/docs/superpowers/plans/2026-08-22-ke-hoach-kh6b-boss.md va cap nhat BAT-DAU-TU-DAY.md (§11) truoc khi commit.

## Trạng thái thật — soát lại bằng CODE ngày 28/08/2026

Bảng task ở trên ghi trạng thái **cũ hơn thực tế**. Đã kiểm từng task bằng cách `grep` chính
mã nguồn (không tin sổ), kết quả:

| Task | Bằng chứng trong code | Trạng thái |
|---|---|---|
| 2 | `lenhThongBaoToanServer` (`thongbao.go`), `lenhChat` (`chat.go`) | ✅ xong |
| 3 | `type Boss struct` (`boss.go`) | ✅ xong |
| 4 | `func goiXuatHienDayDu` + `type HinhDangHien` (`xuathien.go`) | ✅ xong |
| 5 | `type BenDanh` + `KiemDanhHopLe` thuần (`chotkiem.go`) | ✅ xong |
| 6 | `BoPhatThanh` ở **`thongbao.go`** (KHÔNG phải `loathongbao.go`), có từ `6e51dd8` | ✅ xong — soát 28/08 tìm theo **tên file** nên kết luận nhầm |
| 7 | `k.boss = append` + `themBoss` (`khuvuc.go` ~1312) | ✅ xong |
| 8 | định tuyến `-60` id âm → boss, trong `case suKienDanhNguoiChoi` (`khuvuc.go` 1128+, +139 dòng) | ✅ xong 28/08 — dùng lại `KiemDanhHopLe` + `TinhSatThuong` + `bossNhanSatThuong` |
| 9 | `nhipBoss` + `mayTrangThaiBoss` + `bossRoiTheGioi` (`khuvuc.go` 2038+), `Boss.The` nối `TheSongBoss.Tat()` | ✅ xong 28/08 — **trừ nhánh biến hình ở RA_DI**, chờ GATE task 1 |
| **10** | `nhipTim` không xử lý boss | ⬜ **CHƯA LÀM** |
| **11** | `khuvuc.go` có **0 lần** nhắc `MocNoiGian` — tên đó chỉ nằm ở `boss.go` (task 3) | ⬜ **CHƯA LÀM** |
| 12 | `bossNhanSatThuong` + `bossChet` + `traThuongBoss` (`khuvuc.go` 1317+), `PhaBoss.DiemThuong` | ✅ xong 28/08 — **nhưng `bossNhanSatThuong` chưa ai gọi, chờ task 8** |
| 13 | `lichboss.go` + 3 chỗ trong `khuvuc.go` (dòng 400–417, 554–573, 732–739) | ✅ xong 28/08 |
| 14 | `noiBoPhatThanh()` + `GanLoa` 5 khu vực, thứ tự tắt `boLuu → loa → TatMayChu → Close` | ⏳ **nửa xong** — bộ phát thanh đã nối; **lịch boss chưa bật** (có chú thích đánh dấu chỗ) |
| 15 | `noidungboss.go` — 2 boss mẫu (1 con 2 pha + boss lớn + `MocNoiGian`) | ⏳ **nửa xong** — dữ liệu đủ; **smoke client Unity chưa làm** (chỉ chủ dự án làm được) |

**Còn 2 task: `10 → 11`** — cả hai chạm `khuvuc.go` nên **phải tuần tự**.

Task 9 đã để sẵn hai chỗ nối, có đánh dấu trong code:
- bước `BossHoatDong` **để trống** cho task 10, kèm test `TestBossHoatDongDungYen_ChoTask10` ghim nguyên
  trạng "im lặng" — task 10 **buộc phải sửa test đó**, tức có một chỗ bắt dừng lại đọc.
- nhánh **còn pha** ở `BossRaDi` liệt kê đủ những trường phải đặt lại, chờ task 11.

> ⚠️ **Task 1 (GATE hỏi client về biến hình) VẪN CHƯA LÀM** — file research chưa tồn tại.
> Nó **chặn task 11**: chưa biết client có chấp nhận một id đã `-6` quay lại bằng `-5` hay không.
> Làm task 11 mà bỏ qua GATE này là đoán, và đoán sai thì boss biến hình xong **biến mất**.

Task 12, 13, 15 đã xong và gộp `main` (`922eed8`); task 14 xong nửa bộ phát thanh.

> ⚠️ **Hai mối nối đang treo, cố ý:**
> - `bossNhanSatThuong` (task 12) **chưa ai gọi** — task 8 phải gọi nó sau khi tính sát thương,
>   và **tuyệt đối không tự `b.Hp -= ...`** (làm vậy là bỏ qua cả sổ đóng góp lẫn chốt chống chia
>   thưởng hai lần). Câu này đã ghi ngay trong chú thích hàm.
> - `TheSongBoss.Tat()` (task 13) **chưa ai gọi** — task 9 phải gọi khi boss hết pha và rời thế
>   giới, nếu không mỗi mục lịch chỉ sinh **đúng một lần** rồi thôi.
> - Lịch boss **chưa bật** ở `main.go` — task 14 đã để chú thích đánh dấu đúng hai chỗ bật/tắt. Cả 6 lẫn 13 đều chạm `khuvuc.go` (63KB) nên **phải
tuần tự** — đúng cảnh báo ở mục "Thứ tự thi công" bên trên, không giao song song được.

Sẵn sàng cho task 6: `SoKhuVuc.PhatToanServer` đã có ở `sokhuvuc.go:43`; khuôn để bắt chước
là `thogi.go` (hàng đợi không chặn) và `luudinhky.go:111-172` (`BatDau`/`Dung` hai bước).

### Bài học từ lần soát nhầm 28/08

Tôi soát trạng thái bằng cách `grep` theo **tên file kế hoạch dự kiến** (`loathongbao.go`) và
kết luận task 6 chưa làm. Thực tế nó đã xong từ `6e51dd8`, chỉ nằm ở `thongbao.go`.

**Soát trạng thái phải grep theo THỨ CODE LÀM ĐƯỢC, không theo tên file kế hoạch đoán trước.**
Kế hoạch viết trước khi thi công nên tên file trong đó là dự đoán, không phải sự thật.

May là lần soát nhầm đó không vô ích: agent chuyển sang nghiệm thu bằng kiểm thử đột biến và
tìm ra **loa đảo thứ tự thông báo** — bốn thông báo A,B,C,D lên sóng thành A,D,C,B, người chơi
đọc "Boss đã bị hạ" trước "Boss xuất hiện". Đã vá (`0a82053`).

### Soát nhầm LẦN HAI — 28/08, cùng một kiểu sai

Lần một: grep theo **tên file** kế hoạch dự đoán (`loathongbao.go`) → tưởng task 6 chưa làm, thực ra
đã xong ở `thongbao.go`.

Lần hai: grep theo **tên kiểu** trên toàn bộ `*.go` → trúng `boss.go` rồi tưởng task 8–11 đã xong.
Nhưng `boss.go` là sản phẩm của **task 3** (mô hình dữ liệu). Thứ task 8–11 phải làm là **nối chúng
vào actor**, tức phải xuất hiện trong `khuvuc.go` — nơi đó có **0 lần**.

**Cách soát đúng, dùng từ giờ:** grep **trong đúng file mà task đó phải sửa** (cột "File" của bảng
task), không grep toàn thư mục. Một cái tên tồn tại không có nghĩa là nó được gọi.

```bash
# Sai: grep -r "MocNoiGian" *.go          -> trung boss.go, ket luan nham
# Dung: grep -c "MocNoiGian" khuvuc.go    -> 0, tuc CHUA noi day
```

---

## Task 12b — Sổ nợ phần thưởng (28/08/2026, nhánh `so-no-thuong` dựng trên `cf18340`)

> ℹ️ **Ghi chú lịch sử:** nhánh này dựng trên `cf18340` (task 12) khi commit đó còn chưa
> gộp `main`. Trong lúc làm, `main` đã tiến lên: task 12 được gộp (`922eed8`), rồi task 8
> (`15be8df`) và task 9 (`9e6f6c9`). Đã gộp `main` vào nhánh này và giải xung đột —
> `khuvuc.go` tự gộp được, `sokhuvuc.go` chỉ là hai hàm mới đặt cạnh nhau (`LoaDaGan` của
> task 14 và `GanSoNo`/`DongSoNo` của mục này).

**Vấn đề.** Task 12 chốt một hạn chế có chủ đích: người đóng góp sát thương nhưng **đã rời
khu vực / mất kết nối** đúng lúc boss chết thì không nhận được gì — `traThuongBoss` chỉ ghi
một dòng `slog.Warn` kèm số điểm bị bỏ. Lý do khi đó đúng và đến giờ vẫn đúng: actor không
sở hữu con trỏ `*store.NhanVat` của người đã rời (§6.9), còn đẩy một bản chụp cũ qua `ThoGhi`
thì **ghi đè cả dòng**, kéo lùi mọi thứ người đó làm từ lúc rời đi.

**Cách giải.** Mở đường thứ ba: khoản đó thành một **món nợ bền** trong bảng
`thuong_cho_nhan`, đi qua hàng đợi `SoNoThuong` (khuôn `ThoGhi` — actor chỉ chạm một map
trong bộ nhớ vài chục nano-giây, không cham CSDL), và được trả ở **lần đăng nhập sau**.

| Mảnh | File |
|---|---|
| Bảng + ràng buộc | `server/internal/store/migrations/0004_thuong_cho_nhan.sql` |
| Ghi nợ / nhận nợ (một giao dịch) | `server/internal/store/nothuong.go` |
| Hàng đợi ngoài actor | `server/internal/game/sonothuong.go` |
| `traThuongBoss` ghi nợ thay vì bỏ | `server/internal/game/khuvuc.go` (`xepNoThuong`) |
| Trả nợ lúc đăng nhập | `server/internal/game/login.go` (`traNoThuongDangCho`) |
| Nối dây + tắt sạch | `server/internal/game/moluutru.go`, `sokhuvuc.go`, `tatmaychu.go` |

**Ba bẫy đã phải xử — ghi lại kẻo người sau vấp:**

1. **Thứ tự trả nợ không thể đảo.** `store.NhanNoThuong` cộng vào *cột* `kinh_nghiem`, nhưng
   bản `*store.NhanVat` trong RAM (thứ actor sẽ sở hữu từ đó trở đi) vẫn mang số cũ. `ThoGhi`
   ghi đè cả dòng, nên chỉ cần một lần lưu định kỳ 60 giây là con số cũ đó **xoá sạch phần
   vừa cộng** — đúng cái bẫy đã khiến khoản thưởng này không trả được ngay từ đầu, rình sẵn ở
   đầu bên kia. Vì vậy `traNoThuongDangCho` cộng luôn vào `nv` ngay sau khi CSDL cộng xong.
2. **Hàng đợi nợ CỘNG DỒN, không ghi đè.** `ThoGhi` ghi đè vì "bản chụp mới luôn tốt hơn bản
   cũ"; ở sổ nợ mỗi khoản là một món tiền riêng — ghi đè là làm biến mất phần thưởng của
   người chơi. Ghi lỗi thì cả lô được xếp lại (cộng dồn) để thử ở nhịp sau.
3. **Đánh dấu và cộng phải nằm trong MỘT giao dịch.** Tách ra hai giao dịch thì một cú sập
   giữa chừng hoặc ăn mất phần thưởng (đánh dấu xong, chưa cộng) hoặc phát hai lần (cộng
   xong, chưa đánh dấu). Gọi `NhanNoThuong` hai lần liên tiếp phải ra `500` rồi `0`.

**Nối dây có chủ đích tránh `main.go`** (đang thuộc task 14, có agent khác làm): sổ nợ được
mở và cắm trong `MoLuuTru` — đúng chỗ dự án đã chốt là điểm lắp ráp duy nhất — và `TatMayChu`
lấy nó từ `SoKhuVuc` thay vì đổi chữ ký. Ba test khoá riêng việc nối dây, vì bài học của
chính task 12 là `CongSatThuong` từng nằm trong `boss.go` suốt từ task 3 mà không ai gọi tới.

**Nghiệm thu, vòng 1.** 30 test mới, mỗi mảnh đều thấy ĐỎ trước khi viết code; đột biến 8/8
bị giết. `go test ./...` 8/8 gói xanh, `gofmt` + `vet` sạch. `-race` **không chạy được** trên
máy này (thiếu `gcc`, xem §9) — bù bằng `-count=5`.

**Nghiệm thu, vòng 2 — sau khi HAI agent soi đối kháng độc lập.** Vòng 1 chưa đủ. Một agent
soi tính đúng đắn, một agent tấn công chính bộ test bằng 27 đột biến tự nghĩ: **11 bị giết,
16 sống sót**. Kết quả:

- **Một LỖI THẬT không test nào thấy:** lưới chống panic đặt cao hơn một tầng so với chỗ cần.
  `ghiHetHangDoi` lấy hẳn lô nợ ra khỏi hàng đợi rồi mới gọi kho, còn `xepLai` chỉ chạy trên
  đường *trả lỗi* — nên `GhiNoThuong` mà **panic** thì cả lô bốc hơi im lặng, ngược hẳn điều
  chú thích của chính hàm đó cam kết. `ThoGhi` không dính vì nó có lưới **hai tầng**; bản đầu
  chỉ chép tầng ngoài.
- **Ba đường mất/nhân đôi tiền khi CSDL trục trặc:** commit mập mờ → thêm `khoa_lo` + ràng
  buộc duy nhất + `ON CONFLICT DO NOTHING` (kéo theo hàng đợi phải tách hai tầng: lô hỏng thử
  lại NGUYÊN VẸN, khoản mới đi lô riêng); vòng thử lại không có van xả → bỏ lô sau 5 lần;
  `DongVaChoXong` báo "tắt sạch" khi vừa đánh rơi tiền → trả lỗi lần xả cuối, kèm hạn giờ
  riêng trong `TatMayChu`.
- **Chú thích `GanSoNo` chỉ dẫn NGƯỢC thứ tự tắt máy chủ** — chép từ `GanLoa` mà quên sửa.
- **Điểm mù oracle (§6.12)** trong chính test của bản đầu: lấy kỳ vọng từ hằng số mà code
  đang dùng, nên đổi tên hằng thì cả gói `game` vẫn xanh.

Sau vòng sửa: **+20 test**, tổng **50** (7 store + 8 store đường lỗi + 7 hàng đợi +
11 hàng đợi đường hỏng + 7 actor + 7 đăng nhập + 3 nối dây).

**Vòng đột biến thứ hai** (13 cái, tự chạy trên chính bản vá): 10 bị giết, 3 sống sót, đã
bịt cả ba. Cái đầu tiên đáng ghi lại: **test van xả tôi vừa viết cũng sai** — kho giả từ
chối theo *lần gọi* nên lần thứ 6 thành công ở cả hai bản; phải có một lô hỏng **vĩnh
viễn** mới phân biệt được. Và viết lại nó cho thấy hậu quả của việc thiếu van xả **nhẹ hơn**
bản soi mô tả: tách hàng đợi hai tầng (làm vì lý do khác) đã vô tình chữa chỗ "một lô hỏng
làm kẹt cả đường ghi", nên van xả giờ chỉ còn ngăn `choThuLai` phình mãi.

**Về câu hỏi "tính năng có chạy thật không" — đã đổi câu trả lời giữa chừng.** Lúc hai agent
soi (nền `cf18340`), `bossNhanSatThuong` — đường duy nhất nạp `SatThuongTheoNguoi` — **không
có ai gọi ngoài test**, nên `traThuongBoss` chia một map rỗng và sổ nợ chưa từng ghi khoản
nào. Đó đúng là cái bẫy mục này viết ra để cảnh báo, chỉ lùi lên một tầng cao hơn chỗ ba test
nối dây khoá được. **Nhưng trong lúc làm, task 8 đã xong và gộp `main`** (`15be8df` — "người
chơi đánh được boss"), nối đúng mối treo đó: `khuvuc.go` trên `main` gọi `bossNhanSatThuong`
thật khi người chơi đánh trúng boss. Sau khi gộp `main` vào đây thì đường đi đã thông từ đầu
tới cuối. **Vẫn chưa nghiệm thu tay bằng client Unity thật** — đó là việc chỉ chủ dự án làm
được, và không test nào thay thế được.

**Còn nợ lại (không chặn gì):** bảng `thuong_cho_nhan` chỉ tăng, chưa có đường dọn các dòng
`da_nhan = true` cũ; và chỉ hỗ trợ `loai_thuong = 'kinh_nghiem'` — cột để mở nhưng thêm loại
mới phải sửa **cả** ràng buộc CHECK **lẫn** `NhanNoThuong`, nếu không khoản nợ loại mới sẽ
nằm lại vĩnh viễn không ai trả.
