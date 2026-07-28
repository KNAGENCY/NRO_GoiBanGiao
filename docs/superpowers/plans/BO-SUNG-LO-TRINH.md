# Bổ sung lộ trình — những thứ game Naruto online cần mà lộ trình hiện chưa có

> **Nguồn:** rà soát 4 chiều (tính năng game cũ bị rơi · sẵn sàng vận hành · độ bền
> dữ liệu sống · chống lạm dụng ngoài chiến đấu) do phiên điều phối Fable 5 chạy
> 27/07/2026, đối chiếu với `LO-TRINH.md`. Đây là **đề xuất chờ chủ dự án chốt**, chưa
> đưa vào bảng lộ trình chính.

## Cách đọc

~40 điểm gom lại còn **9 khối**. Mức gấp:
- 🔴 **Phải có trước khi mở cửa** — thiếu là mất dữ liệu người chơi, mất tiền, hoặc sập server.
- 🟡 **Nên có sớm** — game vẫn chạy nhưng người chơi bỏ đi / bị quấy rối.
- 🟢 **Để sau cũng được** — làm dày game, không phải điều kiện sống còn.

🧱 = **việc kiến trúc, KHÔNG nhét sau được** — không quyết bây giờ thì sau này đập đi xây lại cả một tầng.

---

## 🔴🧱 Khối 1 — Lưu trạng thái sống (đề xuất **KH11, ngay sau KH6**)

**Lỗ hổng nguy hiểm nhất toàn lộ trình.** `server/internal/store/store.go` chỉ có hàm TẠO
và ĐỌC nhân vật — **không có hàm nào UPDATE bảng `nhan_vat`**. Nhân vật di chuyển, mất máu,
lên cấp đều chỉ nằm trong RAM. Tắt server hay mất điện là **bay sạch tiến độ**.

Là việc kiến trúc vì: actor chạy vòng lặp 100ms; ghi thẳng database trong đó mà DB chậm thì
cả khu vực đứng hình. Phải thiết kế ghi "gửi ra ngoài rồi vòng lặp chạy tiếp", và phải có
**trước KH7** (vật phẩm) — nếu không mỗi hệ tự chế một kiểu lưu, chắp vá và vẫn mất dữ liệu.

Gồm: hàm UPDATE nhân vật; auto-save 5 phút/lần (thiết kế §2.4/§4.4 đã chốt phải có, chưa ai
lên lịch); lưu khi tắt server sạch; "sự kiện chụp ảnh người chơi" an toàn để không phá luật
actor; ghi DB bất đồng bộ; cột `dang_nhap_cuoi` đã khai schema nhưng chưa nơi nào ghi.

**Kèm — chặn đăng nhập trùng:** hai phiên cùng tài khoản tạo hai bản nhân vật độc lập; phiên
cũ đóng xoá nhầm nhân vật phiên mới khỏi khu vực (âm thầm). Cũng là gốc của lỗi nhân bản đồ
(Khối 3). Giải quyết cùng lúc vì cùng đụng "một tài khoản = một phiên sống".

## 🔴🧱 Khối 2 — Vá sập server + tự bật lại (một phần vào KH6/Kế hoạch B NGAY, phần còn lại KH V)

1. **Không có lưới chắn panic ở goroutine.** `protocol/frame.go` tự cảnh báo một lỗi không
   chặn sẽ "sập cả tiến trình, đá văng TOÀN BỘ người chơi". Lưới hiện chỉ bọc một hàm; chỗ
   mỗi-kết-nối và mỗi-khu-vực (nơi KH6 đang viết) chưa bọc → một lỗi lạ khi đánh nhau là sập
   cả server. **Vá trước khi gộp KH6** — việc nhỏ, cực gấp. *(Cùng họ với lỗ hổng #1 —
   actor chặn 30s ở `Gui()`, §6.10b BAT-DAU-TU-DAY.)*
2. **Không có tự bật lại khi sập** — thiết kế hứa "~10 giây" nhưng repo không có gì cho VPS
   Linux. Server sập 3h sáng là nằm chết tới sáng. → KH V, dùng cơ chế OS (systemd restart).

## 🔴🧱 Khối 3 — Chống nhân bản đồ, nhét vào đặc tả KH7 TRƯỚC khi viết code

Tên "KH7 — vật phẩm/cửa hàng/giao dịch" che mất 4 cơ chế khác hẳn nhau (game cũ có): cửa
hàng NPC · **giao dịch 2 người** (`Trade.java` 508 dòng) · **ký gửi/chợ đồ**
(`ConsignShopService.java` 598 dòng — KHÔNG phải cửa hàng NPC) · kho đồ. Thiết kế §6.2 vẽ
đúng hướng (mỗi vật phẩm có "căn cước" + bảng lịch sử) **nhưng chưa migrate vào DB** (mới 2
bảng). Viết KH7 mà chưa có căn cước + chưa chặn "hai giao dịch cùng lúc một món đồ" → người
chơi **nhân bản đồ**, kinh tế game chết không cứu được. Phải chốt căn cước + cách khoá trong
đặc tả KH7 **trước dòng code đầu tiên**.

## 🟡 Khối 4 — Chat + kiểm duyệt + kết bạn (đề xuất **KH12**)

Game online không chat là không giữ được người chơi. Game cũ có: chat toàn server
(`ChatGlobalService.java`), nhắn riêng, chat bang, kết bạn/thù (`FriendAndEnemyService.java`
320 dòng). **Không KH nào nhắc tới.** Kèm kiểm duyệt: game cũ lọc từ bậy bằng thay chuỗi cố
định (dễ lách), **không có "cấm nói" (mute)** — grep 489 file = 0. Admin không có nút bịt
miệng kẻ spam link lừa đảo lên chat toàn server. Đề xuất KH "Xã hội + kiểm duyệt": chat, kết
bạn, **mute có thời hạn**, danh sách từ cấm dùng chung (cũng để lọc tên nhân vật).

## 🟡 Khối 5 — Công cụ quản trị sống (rải vào W3 + KH V + Kế hoạch B)

- **Không có hệ tố cáo** (grep "report" = 0). Người chơi bị lừa/hack không có cách báo. → khối
  "tố cáo + công cụ GM", phụ thuộc KH9 + W.
- **Khoá tài khoản chỉ là ô đúng/sai** — không lý do, thời hạn, nhật ký. → W3: thêm lý do +
  thời hạn + nhật ký ai khoá.
- **Không có bảng đèn báo sống/chết server** — game cũ ghi mỗi giây online/CPU/RAM
  (`panel_server_status.json`), server Go mới có 0. Tới lúc sập mới biết. → KH V.

## 🔴 Khối 6 — Nạp tiền thật (mở rộng W3 — rủi ro cao nhất vì đụng tiền)

W3 chỉ có ba chữ "nạp thẻ". Game cũ có cả cụm: thẻ cào, đối soát ngân hàng, quà theo mốc
nạp, top nạp. Đụng tiền thật trực tiếp — sai một ly là mất tiền/uy tín. Cần thiết kế rõ:
người chơi nạp → bằng chứng gì → ai duyệt → cộng vào đâu → **chống cộng hai lần**. Kèm
Giftcode (cùng là "phát tài sản ảo có kiểm soát"). Tách khối riêng, đừng để chung
đăng ký/đăng nhập.

## 🟡 Khối 7 — Nền tảng vận hành VPS (mở rộng KH V, một phần làm được NGAY)

- 🔴 **Sao lưu database tự động** — hệ cũ đã thiếu, giờ vẫn thiếu. **Mâu thuẫn chưa chốt:**
  thiết kế nói production dùng Supabase, tài liệu vận hành nói tự cài Postgres — mà Supabase
  free **không backup tự động**. Ổ cứng hỏng là mất cả server. Phải chốt trước khi mở cửa.
- 🟡 **Log phình ăn hết ổ** — code tự nhận "9GB/giờ" nếu bật nhầm chế độ chi tiết, chỉ chặn
  bằng "nhớ đừng bật". Cần cắt log cũ tự động.
- 🟡 **Quản lý khoá bí mật production** — đang lặp bẫy "mật khẩu rải 3-4 chỗ" của hệ cũ.
- 🟢 **Cập nhật nóng mất trạng thái** — Khối 1 (auto-save) đã đỡ phần lớn.

## 🟡 Khối 8 — Kỷ luật kỹ thuật cho nhiều agent (KH nhỏ, làm được NGAY)

- **Không có CI** (`.github/workflows` rỗng) — "test xanh trước khi commit" hoàn toàn dựa tự
  giác, mô hình dễ vỡ nhất khi nhiều agent. *(Cộng bẫy test §G-0: `go test ./...` báo `ok`
  cả khi 27 test bị SKIP vì thiếu DSN.)*
- 🧱 **Không có công cụ migration DB** — mới 2 file SQL áp tay, không lần vết, không lùi được.
  KH7/8/9 sắp đổ **rất nhiều** bảng. Áp migration TRƯỚC khi schema phình.
- **Không quét lỗ hổng thư viện định kỳ** (npm đã có 3 lỗi mức cao).

## 🔴🧱 Khối 10 — Mapping ảnh Naruto vào client (thêm 28/07/2026)

**Tại sao có khối này.** 27-28/07 chủ dự án nghiệm thu Unity: 27/07 client kẹt "Xin chờ",
28/07 vào được game sau chuỗi 4 fix (chi tiết BAT-DAU §12), nhưng **nhân vật không hiển
thị**, **map trống**, **background Dragon Ball**. Nguyên nhân:
bộ 343 ảnh Naruto ở `docs/my-thuat/anh-game/x1..x4/` đã vẽ xong nhưng **chưa map vào
runtime** — server vẫn gửi `part.bin` placeholder Dragon Ball (`dongbo.go:15`), atlas
`Big0..4.png` trong Unity vẫn là Dragon Ball, thư mục `mymap/` chỉ có 4 file map cũ.

**Bốn hạng mục con** (không thể làm ẩu):

1. **Format binary phải reverse-engineer chính xác.** `part.bin` không chứa ảnh — chỉ
   metadata `short partSize` + `[byte type + pi.Length × (short id, sbyte dx, sbyte dy)]`
   với `pi.Length` theo type (head=3, body=17, leg=14, aux=2 — xem `Game1/Part.cs:10-29`).
   `id` là chỉ số vào atlas `NR_image` — cũng phải reverse-engineer format bảng tra
   (`Game1/SmallImage.cs:61-87`). Sai 1 byte = client render sai hoặc `IndexOutOfRange`
   nuốt im lặng như bẫy §6 BAT-DAU.

2. **Bộ 343 ảnh thiếu mapping ID số → tên file.** `manifest-tao-anh.json` chỉ có
   `ten_file/nhom/kich_thuoc/prompt`, không có ID. `HUONG-DAN-TAO-ANH.md:185-188` nói
   "ID = vị trí trong danh sách" — phải **tạo danh sách chính thức riêng** ở server
   (không được sắp lại, chỉ được thêm ở cuối; ngược lại phá cache client).

3. **Background parallax 4-5 lớp** (`GameCanvas.cs:1510-1701`) — mỗi map tải `b<typeBg>{0,1,2,3}.png`
   theo lớp gần-xa cuộn tốc độ khác nhau. Ảnh Naruto (`map_la_co_1_x*.png`) là **single
   image** → phải cắt lớp thủ công **hoặc** rework `paintBackgroundtLayer` để dùng single.

4. **Deploy chưa có script** — `TAO-4-MUC-ZOOM.py` chỉ sinh 4 zoom vào `anh-game/x*/` phẳng.
   Không có script copy sang `D:\NRO_Project\Assets\Resources\res\x*/` theo đúng
   `Big*.png`/`Small<id>.png`/`b<typeBg>{0..3}.png`. Cần viết + test.

**Điểm dừng an toàn.** Không đụng vào `part.bin` cho tới khi đã có toàn bộ 4 hạng mục
byte-exact spec + subagent + review đối kháng — bản Java Bootstrap vẽ nhân vật KHÔNG kiểm
biên (`parts[head]` truy cập thẳng): sai 1 id = client sập vẽ. **Không có POC nhanh** —
mọi cách "bằng tay copy đại" đều thất bại (27/07 đã thử: rename `map_la_co_1.png` → `b1.png`
→ vô nghĩa vì background là parallax, không phải 1 ảnh).

**Đề xuất trình tự** (KH riêng, ước ~5-8 task như KH6/KH11):
1. Phase 0: subagent phân tích byte-exact **cả 3 format**: `NR_image`, `part.bin`,
   `Big*.png` packing scheme — dùng `1-MAY-CHU-GAME/data/update_data/part` làm bằng chứng.
2. Phase 1: tạo `docs/my-thuat/mapping-id.json` — danh sách ID chính thức, chốt với chủ dự án.
3. Phase 2: viết script Python: pack ảnh → `Big*.png` + `NR_image` + `part.bin` byte-exact.
4. Phase 3: script deploy: rename + copy ảnh vào Unity Resources đúng thư mục.
5. Phase 4: cắt lớp background (79 ảnh `map_*`) HOẶC rework render (chọn 1). Chốt phương án.
6. Phase 5: test byte-by-byte + client Unity render đúng.

**Blocker của các việc khác.** Task 15 (nghiệm thu tay KH6) và Task 9 (nghiệm thu tay KH11)
đều **không hoàn thành được** khi nhân vật không hiện. Vì vậy khối này chặn cả hai nghiệm
thu.

---

## 🟢 Khối 9 — Hệ tính năng làm dày game (phần lớn để sau)

- **Bang hội/gia tộc** — `ClanService.java` 1085 dòng, hệ lớn nhất. Giữ chân lâu dài nhưng
  nặng → KH riêng (KH13), đừng nhét KH7.
- **Thú cưng/triệu hồi** — `Pet.java` 1540 dòng. Thiết kế có ý "thú triệu hồi" (§7.1), để sau.
- **Thành tựu + danh hiệu** — khác "thang bậc nhẫn giả" (sức mạnh cốt lõi). Để sau.
- **Khung sự kiện định kỳ** — game cũ có 7 loại sự kiện + giải đấu. **Cần quyết sớm:** "Kỳ
  thi Chunin" (§7.1/§7.3) là nhiệm vụ 1 lần (KH9) hay sự kiện định kỳ (cần khung riêng)? Kẻo
  KH9 làm nửa vời rồi viết lại.
- **Nhiệm vụ phụ + điểm danh + phúc lợi hằng ngày** — dễ bị bỏ sót vì tên "KH9" gợi chỉ 36
  nhiệm vụ chính. Ghi rõ vào phạm vi KH9.
- **Chế tác + cường hoá** — `CombineService.java` 698 dòng. KH8 chỉ lo *số liệu*, không lo
  *luồng*. Ghi rõ luồng chế tác thuộc KH7 hay KH8.
- **Boss sống qua restart** — game cũ coi boss khác quái thường. Chốt trong đặc tả KH6b.

---

## Kết luận

Hai lỗ hổng lớn nhất đều **kiến trúc, không nhét sau được**: (1) **server không lưu gì cả** —
không một dòng code cập nhật nhân vật xuống DB, tiến độ bốc hơi mỗi lần tắt máy; (2) **KH7
chưa có căn cước vật phẩm + chống nhân bản** — mở cửa với lỗ này thì kinh tế game bị bơm đồ
giả, không cứu được. Cả hai phải thiết kế **trước** khi viết tiếp.

Lộ trình cũng **thiếu hẳn tầng "vận hành sống"** mà game cũ đã có: chống sập/tự bật lại, sao
lưu DB (còn mâu thuẫn Supabase vs tự host), giám sát online, kiểm duyệt chat/mute, hệ tố cáo.

**Tối thiểu trước ngày mở cửa:** KH11 (lưu trạng thái) + vá panic trong KH6 + đặc tả chống
nhân bản trong KH7 + backup DB + kiểm duyệt chat cơ bản. Bang hội/thú cưng/thành tựu/sự kiện
là đồ trang trí, để sau. Việc **nên khởi động ngay hôm nay** vì càng để càng đắt: **migration
DB + CI** — trước khi vài KH nữa đổ hàng chục bảng vào một schema không ai lần vết được.
