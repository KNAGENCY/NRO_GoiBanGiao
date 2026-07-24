# Kế hoạch 5 — Di chuyển + mô hình actor cho khu vực (Thiết kế)

**Ngày:** 2026-07-24
**Nối tiếp:** Kế hoạch 4 (nhân vật đã vào thế giới; `KhuVuc` tối thiểu đã tồn tại, `xuLyClientOk` gửi `-82/-30/-24`).
**Nghiên cứu nền:** `docs/superpowers/research/kh5-goi-dichuyen-len.md`, `kh5-goi-phat-xuong.md` (đã giải mã + soi chéo đối nghịch 0 sai lệch).
**Ràng buộc mượt/chống gian lận:** `docs/superpowers/specs/2026-07-23-hoa-giai-muot-va-chong-gian-lan.md` §4.3 (nhịp 100ms, server authoritative, client không dự đoán).

---

## 1. Mục tiêu

Nhân vật **đi lại được** trong bản đồ làng; **máy chủ giữ vị trí thật** (authoritative) và từ chối di chuyển gian lận ở mức **vừa phải**; đặt **nền kiến trúc actor** để các kế hoạch sau (quái, người chơi khác, đánh nhau) cắm vào mà không phải viết lại.

Bằng chứng "xong": hai người chơi giả trong cùng một khu vực **thấy nhau di chuyển** qua mạng thật; gói di chuyển vượt ngưỡng bị từ chối; vị trí trên server luôn đúng.

## 2. Kiến trúc — khu vực là một "actor"

`KhuVuc` (đã có ở KH4, hiện chỉ giữ danh sách người chơi) được nâng thành **actor**: một **luồng riêng (goroutine)** sở hữu toàn bộ trạng thái khu vực + một **hộp thư (channel)** để nhận sự kiện + một **nhịp tim 100ms**.

```
   Phiên A (goroutine)          Phiên B (goroutine)
        │  gửi sự kiện                │  gửi sự kiện
        └──────────┐      ┌──────────┘
                   ▼      ▼
              ┌──────────────────┐
              │  hộp thư KhuVuc   │  (channel có đệm)
              └────────┬─────────┘
                       ▼
          Vòng lặp KhuVuc (1 goroutine duy nhất)
          - xử lý sự kiện tuần tự (không tranh chấp)
          - mỗi 100ms: nhịp tim (chỗ cho quái ở KH6)
          - phát tin: gọi s.Gui của từng phiên (đã an toàn nhờ muGui)
```

**Vì sao actor (không phải khoá mutex quanh trạng thái chia sẻ):**
- Trạng thái khu vực (ai đang ở đâu) chỉ **một luồng chạm tới** → không cần khoá, không có tranh chấp, không deadlock. Suy luận dễ.
- Khớp máy chủ Java cũ (`GameLoop` một vòng lặp/nhịp). Nhịp 100ms có sẵn chỗ cho quái/vật phẩm tự di chuyển ở KH6 mà không đổi kiến trúc.
- Phiên đã có `muGui sync.Mutex` bảo vệ ghi ra kết nối → luồng khu vực gọi `s.Gui(...)` phát tin **an toàn ngay**, không cần kênh gửi riêng cho mỗi phiên (YAGNI).

**Sự kiện hộp thư (kiểu):**
| Sự kiện | Khi nào | Dữ liệu |
|---------|---------|---------|
| `SuKienVao` | phiên nhận `-39` (đã tải xong map) | nhân vật + phiên |
| `SuKienDiChuyen` | phiên nhận `-7` (đã qua kiểm hợp lệ) | id nhân vật, x, y đích |
| `SuKienRa` | phiên đóng (ngắt kết nối / đổi map sau này) | id nhân vật |
| *(nhịp tim)* | mỗi 100ms từ `time.Ticker` | — (KH5 chưa dùng; chỗ cho KH6) |

**Vòng lặp** dùng `select` trên `{hộp thư, ticker.C, kênh dừng}`. Kênh dừng để tắt sạch actor khi khu vực rỗng/tắt máy (tránh rò goroutine).

## 3. Vòng đời & sở hữu

- KH5 vẫn **một khu vực cho mỗi phiên** như KH4 (chưa gộp người chơi chung map — xem §7 Phạm vi). Nhưng bộ máy phát tin phải **đúng cho N người**: test dùng 2 phiên trỏ vào **cùng một** `KhuVuc` để chứng minh họ thấy nhau. Việc "nhiều phiên chia sẻ một khu vực theo mapId" là một task nhỏ trong kế hoạch (bảng khu vực theo mapId), không phải chương riêng.
- Actor khởi động khi khu vực có người đầu tiên; dừng (đóng kênh dừng) khi người cuối rời — hoặc đơn giản hơn cho KH5: actor sống theo vòng đời tiến trình cho 5 map cố định. **Quyết định:** tạo sẵn 5 actor khu vực (mapId 1-5) lúc khởi động máy chủ, sống suốt đời tiến trình. Đơn giản, không cần logic tạo/huỷ động; 5 goroutine ngủ là chi phí không đáng kể. Tạo/huỷ động để dành khi có map tạm (tổ đội/sự kiện) ở kế hoạch sau.

## 4. Giao thức (đã giải mã byte-đúng, đã soi chéo)

### 4.1 `-7` NHẬN (client → server): yêu cầu di chuyển
| # | Trường | Kiểu đọc | Ghi chú |
|---|--------|----------|---------|
| 1 | cờ bay (khinh công) | `readByte` | dùng cho thành tích "bay"; KH5 chỉ đọc bỏ qua (không nghiệp vụ) |
| 2 | toX | `readShort` | tọa độ đích X |
| 3 | toY | `readShort` | tọa độ đích Y — **optional**: bản gốc đọc trong try/catch riêng; thiếu thì giữ Y hiện tại |

Server **không** reject cứng khi thiếu byte — fallback về vị trí hiện tại (khớp bản gốc).

### 4.2 `-7` PHÁT (server → client): thông báo di chuyển
Gửi cho **người khác** trong khu vực (KHÔNG dội lại chính người vừa đi — client tự vẽ chuyển động của mình).
| # | Trường | Kiểu ghi |
|---|--------|----------|
| 1 | playerId | `writeInt` |
| 2 | x | `writeShort` |
| 3 | y | `writeShort` |

### 4.3 `-39` NHẬN: "đã tải xong bản đồ" → kích xuất hiện
Client gửi `-39` sau khi nhận đủ `-82/-30/-24` và vẽ xong map. Server phản ứng:
1. Gửi `-5` (xuất hiện) của **từng người đang có mặt** trong khu vực cho người mới (để người mới thấy người cũ).
2. Phát `-5` của **người mới** cho những người đang có mặt (để người cũ thấy người mới).
Với một người chơi: đăng ký vào khu vực, không có ai để gửi → 0 byte ra dây (đúng bản gốc).

### 4.4 `-5` PHÁT: người chơi xuất hiện (28 trường)
Thứ tự byte đầy đủ (từ `Zone.infoPlayer`, đã xác minh từng dòng). KH5 dùng giá trị tối thiểu an toàn cho các trường chưa có hệ thống (buff/pet/ngựa/aura = 0):

| # | Trường | Kiểu ghi | Nguồn giá trị KH5 |
|---|--------|----------|-------------------|
| 1 | playerId | writeInt | `nv.ID` |
| 2 | clanId | writeInt | -1 (chưa có bang hội) |
| 3 | level | writeByte | `byte(nv.Cap)` — cấp thật của nhân vật (field `Cap` đã có; nhân vật mới mặc định Cap=1 nên hành vi hiện tại y hệt "dùng 1", nhưng đúng ngay khi có hệ lên cấp) |
| 4 | cờ cố định | writeBoolean | false |
| 5 | typePk | writeByte | 0 |
| 6 | gender | writeByte | `byte(nv.Lang)` |
| 7 | gender (lặp) | writeByte | `byte(nv.Lang)` |
| 8 | head | writeShort | `nv.KieuToc` |
| 9 | name | writeUTF | `nv.Ten` |
| 10 | hp | **writeInt** | `nv.HpToiDa` (client đọc int vì `cn.readInt=true`) |
| 11 | hpMax | **writeInt** | `nv.HpToiDa` |
| 12 | body | writeShort | 0 (chưa có trang bị) |
| 13 | leg | writeShort | 0 |
| 14 | flagbag | writeByte | 0 |
| 15 | hằng -1 | writeByte | -1 |
| 16 | x | writeShort | `nv.ViTriX` |
| 17 | y | writeShort | `nv.ViTriY` |
| 18 | effBuffHp | writeShort | 0 |
| 19 | effBuffMp | writeShort | 0 |
| 20 | numEff | writeByte | 0 |
| 21 | idSpaceShip | writeByte | 0 |
| 22 | isMonkey | writeByte | 0 |
| 23 | mount | writeShort | 0 |
| 24 | cFlag | writeByte | 0 |
| 25 | hằng 0 | writeByte | 0 |
| 26 | aura | writeShort | 0 |
| 27 | effFront | writeByte | 0 |
| 28 | hat | writeShort | 0 |

*(Bản gốc còn gửi kèm `-103` cờ và `-8` nếu đang chết. KH5: nhân vật luôn sống, chưa có cờ bang → **bỏ** cả hai; ghi chú trong mã để KH sau thêm.)*

### 4.5 `-6` PHÁT: người chơi rời đi
| # | Trường | Kiểu ghi |
|---|--------|----------|
| 1 | playerId | writeInt |

## 5. Chống gian lận — mức "Vừa phải" (người dùng chọn)

Kiểm **trước khi** chấp nhận `SuKienDiChuyen`, áp cho **mọi map** (bản gốc chỉ áp map 85-91 — lỗ hổng):

1. **Trong biên bản đồ:** `0 ≤ toX ≤ rộng_map` và `0 ≤ toY ≤ cao_map`. Kích thước lấy từ dữ liệu tile map (định dạng đã giải ở KH4 §10; nếu chưa nạp, dùng biên rộng an toàn theo hằng số map, ghi chú TODO nạp thật). Ngoài biên → **từ chối gói, giữ vị trí cũ**, ghi log mức debug (không spam).
2. **Khoảng cách mỗi bước:** `khoảng_cách_Euclid(vị_trí_cũ, đích) ≤ NGƯỠNG_BƯỚC`. Ngưỡng theo **từng bước nhảy** (không theo giây) nên khoan dung với mạng lag. Đề xuất `NGƯỠNG_BƯỚC = 500` (khớp ngưỡng bản gốc dùng cho map sự kiện) — vượt → từ chối gói, giữ vị trí cũ.

**Không** kiểm theo thời gian/tốc độ ở KH5 (đó là mức "Chặt", để dành). Từ chối = **im lặng bỏ gói + giữ authoritative**, KHÔNG ngắt kết nối (tránh đá oan; client sẽ tự đồng bộ lại khi nhận gói khác).

## 6. Điểm chạm & luồng dữ liệu

- **`login.go` / bộ điều phối:** thêm `case -7` và `case -39` vào switch xử lý gói. Mỗi case: đọc/kiểm, rồi **gửi sự kiện vào hộp thư khu vực** của phiên (`b.khuVuc`), KHÔNG tự chạm trạng thái khu vực.
- **`khuvuc.go`:** nâng `KhuVuc` thành actor (hộp thư, vòng lặp, ticker, kênh dừng); các hàm dựng gói `-5/-7/-6`; hàm phát (lặp qua người chơi, gọi `s.Gui`).
- **Kiểm hợp lệ di chuyển:** hàm thuần `kiemDiChuyenHopLe(map, cũ, đích) bool` — dễ test riêng, không cần mạng.
- **`xuLyClientOk` (KH4):** sau khi gửi `-82/-30/-24`, phiên đã có `b.khuVuc`. KH5: khi nhận `-39`, gửi `SuKienVao`.
- **Khởi động máy chủ:** tạo 5 actor khu vực (mapId 1-5), lưu vào một sổ tra (`map[int32]*KhuVuc` hoặc slice) để phiên tìm khu vực theo mapId của nhân vật khi vào.

**Chú ý sở hữu phiên:** `SuKienVao` mang cả `*mynet.Session` để khu vực phát tin ngược. Khu vực giữ `map[int64]*nguoiTrongKhuVuc{ nv, session }`. Khi phiên đóng, gửi `SuKienRa` để khu vực gỡ + phát `-6`.

## 7. Phạm vi (YAGNI) — làm gì / KHÔNG làm

**Làm:** actor khu vực + nhịp 100ms; nhận `-7` + kiểm hợp lệ vừa phải + cập nhật authoritative + phát `-7`; nhận `-39` + phát `-5` hai chiều; `-6` khi rời; 5 khu vực cố định; test 2 người chơi giả thấy nhau qua mạng thật.

**KHÔNG làm (kế hoạch sau):** quái/NPC di chuyển (KH6 dùng nhịp 100ms có sẵn); đánh nhau/skill (KH6); vật phẩm rơi/nhặt (KH7); pet/ngựa/aura/bang hội (điền 0, ghi chú); đổi map giữa các khu vực (khi có nhiều map thật); gói `-103/-8`; kiểm gian lận theo thời gian (mức "Chặt").

## 8. Xử lý lỗi & biên

- **Client chậm không được làm kẹt cả khu vực:** `s.Gui` ghi TCP có thể chặn nếu client nghẽn. Đặt **hạn ghi (write deadline)** hợp lý trên kết nối khi phát broadcast, hoặc bọc phát trong cơ chế không để một phiên treo vòng lặp khu vực. KH5 tối thiểu: dựa vào hạn ghi sẵn có của phiên; nếu `s.Gui` lỗi → ghi log + đánh dấu phiên để gỡ (gửi `SuKienRa`), KHÔNG panic vòng lặp.
- **Hộp thư đầy:** channel có đệm (ví dụ 256). Nếu đầy (phiên gửi quá nhanh — dấu hiệu flood), **bỏ gói di chuyển** (giữ authoritative) thay vì chặn phiên. Ghi nhận để chống flood ở Kế hoạch B.
- **Sự kiện tới sau khi actor dừng:** gửi vào channel đã đóng gây panic — dùng `select` với kênh dừng, hoặc `recover` mềm; KH5 (actor sống suốt đời tiến trình) ít gặp, nhưng vẫn phòng.
- **Nhân vật ở mapId không có khu vực:** không thể xảy ra (mapId luôn 1-5), nhưng nếu tra khu vực trả nil → log + không đăng ký (giống nhánh `bd==nil` KH4).

## 9. Kiểm thử

- **Đơn vị (không mạng):** `kiemDiChuyenHopLe` — trong biên/ngoài biên/vượt ngưỡng bước; dựng gói `-5/-7/-6` byte-đúng (đọc lại bằng Reader, khẳng định hết buffer như KH4).
- **Actor (không mạng thật, dùng phiên giả hoặc channel):** gửi `SuKienVao` × 2 → mỗi người nhận `-5` của người kia; gửi `SuKienDiChuyen` của A → B nhận `-7 [idA][x][y]`, A KHÔNG nhận lại; `SuKienRa` của A → B nhận `-6 [idA]`.
- **Đầu-cuối (mạng + CSDL thật, nối tiếp `TestVersionDanceDayDu`):** 2 client giả cùng vào mapId 1 → sau `-39`, mỗi client nhận `-5` của người kia; client A gửi `-7` → client B nhận `-7`; client A gửi `-7` vượt ngưỡng → B KHÔNG nhận, vị trí A trên server không đổi.
- **Đua dữ liệu:** chạy các test actor với `-race` để chứng minh không tranh chấp (mô hình actor phải sạch race).

## 10. Rủi ro & quyết định đã chốt

| Vấn đề | Quyết định |
|--------|-----------|
| Khoá mutex hay actor? | **Actor** (một luồng sở hữu trạng thái) — dễ suy luận, khớp Java, nền cho quái KH6 |
| Kênh gửi riêng mỗi phiên? | **Không** — `muGui` của phiên đã bảo vệ ghi; gọi trực tiếp `s.Gui` (YAGNI) |
| Tạo/huỷ khu vực động? | **Không** — 5 actor cố định suốt đời tiến trình (đơn giản) |
| Độ chặt chống gian lận | **Vừa phải** — biên + khoảng cách/bước, mọi map; không kiểm theo thời gian |
| Dội `-7` lại cho chính mình? | **Không** — client tự vẽ (khớp bản gốc, giữ mượt) |
| Client chậm làm kẹt khu vực | Hạn ghi + gỡ phiên lỗi; đệm hộp thư; không để treo vòng lặp |

---

**Bước tiếp:** dùng kỹ năng `writing-plans` để chia đặc tả này thành các task nhỏ (TDD, commit thường xuyên), rồi thực thi bằng `subagent-driven-development` như Kế hoạch 4.
