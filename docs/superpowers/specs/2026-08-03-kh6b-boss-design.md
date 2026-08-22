# KH6b — Boss (thiết kế)

**Ngày:** 03/08/2026 · **Trạng thái:** thiết kế (chờ duyệt) · **Nghiên cứu nguồn:** [`../research/2026-08-03-kh6b-boss-research.md`](../research/2026-08-03-kh6b-boss-research.md) (đối chiếu Java + client Unity).

## 1. Phát hiện cốt lõi + vì sao

**Boss KHÔNG phải quái — nó là một "người chơi do máy điều khiển"** (`boss.Boss extends Player` bản Java). Với client, boss **là một nhân vật**. Hệ quả then chốt:
- Boss **không có giao thức mạng mới** — tái dùng hết gói người chơi: xuất hiện `-5`, rời `-6`, chat `44`, đánh nhau `-60`/`-45` (KH6). **Chỉ 1 gói mới: `93`** (thông báo chạy chữ toàn server).
- Phần khó của KH6b nằm ở **AI + vòng đời + dữ liệu**, không ở byte.

→ KH6b = một **thực thể server-điều-khiển trong actor khu vực** (không có phiên mạng), chạy máy trạng thái AI, phát các gói người-chơi đã có.

## 2. Quyết định đã chốt với chủ dự án

| # | Quyết định | Chốt | Cách hiện thực |
|---|---|---|---|
| 1 | Xuất hiện | **Cả lịch + trigger** | v1: theo lịch (state REST→RESPAWN theo `secondsRest`). Thêm đường spawn theo trigger (sự kiện/nhiệm vụ gọi `SinhBoss(idLoai, mapId)`). |
| 2 | Pha | **Cả biến-hình-sau-chết + nổi-giận-%máu** | Biến hình = mảng pha `[]PhaBoss` (giết pha → hồi sinh pha mạnh hơn). Nổi giận = mỗi pha có `mocNoiGian []MocHP` (xuống %máu → đổi hệ số dame/def/nhịp, chat). |
| 3 | Chia thưởng | **Theo đóng góp** | Actor lưu `map[idNguoi]tongSatThuong` cho mỗi boss; khi chết chia thưởng theo tỉ lệ. Bản 1 thưởng = **điểm săn boss** (đồ rơi hoãn tới sau KH7). |
| 4 | Thông báo | **Boss lớn** | Gói `93` khi boss lớn (`laBossLon=true`) xuất hiện + bị hạ; qua kênh mới `PhatToanServer` (rate-limit ~1,5s như bản gốc). |

**Quyết định kỹ thuật (em tự chốt):** (a) **một mô hình dữ liệu boss** (không 95 class con) — hành vi đặc biệt khai báo bằng cờ/tham số; (b) boss **dùng chung catalog chiêu KH6** (thêm cờ "chiêu chỉ-boss" sau nếu cần).

## 3. Kiến trúc — boss trong actor khu vực

Điều kiện tiên quyết **ĐÃ ĐẠT**: KH5 đồng bộ nhân vật giữa các phiên (`khuvuc.go` `suKienVao` gửi `-5` hai chiều). Boss là một "char" trong khu vực **nhưng không có phiên**.

- **Thực thể boss** (`Boss` struct): `id` (int âm **runtime duy nhất**, KHÔNG dùng id-loại), `viTri`, `hp/hpMax`, `outfit[6]`, `phaHienTai`, `duLieu []PhaBoss`, `trangThai BossState`, `mocThoiGian` (cho state machine), `satThuongTheoNguoi map[int64]int64`.
- **Nằm trong `KhuVuc`** (thêm `boss []*Boss` cạnh `quai`), **CHỈ luồng `vongLap` chạm** (luật §6.9/§6.10b — KHÔNG luồng 150ms riêng như Java; boss update ngay trong `nhipTim` 100ms).
- **Xuất hiện:** khi boss vào khu vực, actor phát `goiXuatHien(bossNhuNhanVat)` cho cả khu vực (tái dùng `xuathien.go`), **điền outfit thật + typePk=PK_ALL**. Người mới vào cũng nhận `-5` của boss (sửa `suKienVao` để lặp cả `k.boss`).
- **Người đánh boss:** client gửi `-60` với `id` = id boss (âm). Actor định tuyến: id âm + có trong `k.boss` → áp sát thương lên boss (thay vì người chơi), cộng `satThuongTheoNguoi`.
- **Boss đánh/đi/chat:** trong `nhipTim`, boss chọn mục tiêu + chiêu → `phatChoCaKhuVuc(goiDanhNguoiChoi(...))` (tái dùng KH6); ngoài tầm → di chuyển (`phatChoNguoiKhac(goiDiChuyen(...))`); chat → `phatChoCaKhuVuc(goiChat44(...))`.

## 4. Mô hình dữ liệu (data-driven)

```
type PhaBoss struct {
    Ten        string
    Lang       int16       // lang Naruto thay gender Dragon Ball
    Outfit     [6]int16    // head, body, leg, flagBag, aura, effFront
    Dame       int64
    HpChon     []int64     // spawn chon ngau nhien mot gia tri
    MapVao     []int32     // chon ngau nhien khu vuc
    BoChieu    [][2]int32  // {templateId, level}
    ThoaiHien  []string    // chat luc xuat hien
    ThoaiDanh  []string    // chat luc dang danh
    ThoaiChet  []string    // chat luc chet
    GiayNghi   int32       // truoc khi hoi sinh (pha DEFAULT)
    KieuHien   KieuXuatHien // DEFAULT / BIEN_HINH / CUNG_BOSS_KHAC / GOI_BOI_BOSS_KHAC
    MocNoiGian []MocNoiGian // (quyet dinh 2) xuong %HP -> doi he so
    HanhDongDacBiet *HanhDongDacBiet // (hoan) bom / hap thu / trieu hoi
}
type MocNoiGian struct { PhanTramHP int; HeSoDame, HeSoDef float64; Chat string }
```
Boss cụ thể = `[]PhaBoss` khai báo trong `noidungboss.go` (như `BossesData.java`, nhưng data-thuần Go; sau KH8 chuyển sang YAML). HP kẹp `[0, 2_147_483_647]` (§6.17), chặn override > 2 tỷ.

## 5. Vòng đời (máy trạng thái, chạy trong actor tick 100ms)

Rút gọn từ 9 pha Java, giữ đủ ý (chi tiết research §3):
`NGHI` (chờ `GiayNghi`/trigger) → `HIEN` (chọn map, phát `-5`, announce `93` nếu boss lớn) → `CHAO` (phát `ThoaiHien` qua `44`) → `HOAT_DONG` (mỗi tick: chọn mục tiêu + chiêu/di chuyển + thỉnh thoảng `ThoaiDanh`; theo dõi `MocNoiGian`) → (chết) → `TROI` (phát `ThoaiChet`) → `RA_DI` (còn pha → về `HIEN` pha+1 giữ **cùng id**; hết pha → `-6` + về `NGHI`).

Chết (hp≤0 trong lúc áp sát thương): announce `93` "ai đó đã hạ boss" (nếu boss lớn), **chia thưởng theo `satThuongTheoNguoi`**, xoá map đóng góp, chuyển `TROI`.

## 6. Phạm vi

**Bản 1 (KH6b-1) — làm:** thực thể boss trong actor + xuất hiện `-5` (outfit thật) + đi lại + đánh người bằng chiêu KH6 + người đánh boss (định tuyến `-60`) + **biến hình qua pha** + **nổi giận theo %máu** + **chia thưởng theo đóng góp** (điểm) + **announce `93` boss lớn** + spawn theo **lịch**.

**Bản 2 (KH6b-2) — sau:** spawn theo **trigger** (sự kiện/nhiệm vụ), **hành động đặc biệt** (gồng bom diện rộng, hấp thụ, triệu hồi lính/phân thân — nền `bossAppearTogether`).

**Hoãn (phụ thuộc việc khác / để sau):** rơi đồ ra đất (cần **KH7 vật phẩm**); panel `boss_status.json` + lệnh respawn/edit; menu boss `-96`/`-118`; chuyển dữ liệu boss sang YAML (**KH8**).

## 7. Bẫy — đọc trước khi code (từ research §11)

1. **Đừng nhầm** `mob.BigBoss` (quái to, KH6) với `boss.Boss` (Player-AI, KH6b).
2. **`id` runtime DUY NHẤT** mỗi boss (bộ đếm âm toàn cục) — bản gốc spawn 5 con cùng `id=-1` là bug; client `findCharInMap(id)` cần id không trùng.
3. **Biến hình giữ NGUYÊN id** (pha 1 hiện lại cùng id, outfit mới) — ⚠️ **PHẢI kiểm client** chấp nhận một id đã `-6` quay lại bằng `-5` **trước khi** làm biến hình (dùng Unity MCP đọc `Controller.cs`).
4. **HP kẹp int32**, chặn cộng dồn tràn (hấp thụ), chặn override > 2 tỷ.
5. **Boss AI trong actor** (không luồng riêng) — nếu không sẽ tranh chấp `viTri` (bẫy §6.9).
6. **Announce `93` toàn server** cần kênh `PhatToanServer` (khác mọi gói combat chỉ trong khu vực) + hàng đợi/nhịp chống ngập.
7. Boss **né/khiên/vô hạn chakra** — mô hình bằng cờ tường minh, đừng bê số ma thuật `31_07_2002`.

## 8. Ngoài phạm vi (YAGNI)

- 95 class con (dùng data-driven).
- Toàn bộ ~120 boss nội dung (KH6b làm **khung + 1-2 boss mẫu**; thêm boss là dữ liệu, không phải code — nhất là sau KH8/YAML).
- Cân bằng chi tiết (chỉ số cụ thể) — chốt khi có content.

## 9. Thứ tự thi công (chi tiết ở plan)

1. **Kiểm client biến-hình-giữ-id** bằng Unity MCP (bẫy #3) — gate: nếu client không chấp nhận, đổi thiết kế (dùng id mới mỗi pha + `-6`/`-5`).
2. Kênh `PhatToanServer` + gói `93` (test byte).
3. Thực thể boss + xuất hiện `-5` outfit thật + định tuyến `-60` người-đánh-boss (TDD).
4. Máy trạng thái + AI đánh/đi/chat trong actor (TDD, test nhịp).
5. Biến hình pha + nổi giận %máu + chia thưởng đóng góp (TDD).
6. Boss mẫu + smoke bằng client thật (Unity MCP).
