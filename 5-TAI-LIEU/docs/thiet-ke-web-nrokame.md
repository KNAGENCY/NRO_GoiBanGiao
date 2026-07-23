# Hồ sơ thiết kế website — Ngọc Rồng Kame

Tài liệu dùng để: (1) sinh ảnh mockup, (2) code lại web theo ảnh, (3) biết thứ gì đổi được trong game.

Layout **làm mới hoàn toàn**. Chỉ giữ lại chức năng, không giữ bố cục cũ.

---

# PHẦN 1 — Chức năng bắt buộc phải có

Đây là hợp đồng. Thiết kế mới có thể sắp xếp, gộp, đổi chỗ tùy ý — nhưng **không được thiếu** mục nào.

| # | Chức năng | Ghi chú kỹ thuật |
|---|---|---|
| 1 | **Đăng nhập** | gọi `/sv1/dangnhap` — có khoá luỹ tiến khi sai nhiều lần |
| 2 | **Đăng ký** | gọi `/sv1/dangky` — cần tài khoản, mật khẩu ≥6 ký tự, SĐT Việt Nam |
| 3 | **Quên mật khẩu** | gọi `/sv1/quenmatkhau` — chỉ xác minh, không tự đặt lại |
| 4 | **Đổi mật khẩu** | gọi `/sv1/doimatkhau` |
| 5 | **Tải game** | 3 nền tảng: PC / Android / iOS. Luồng 2 bước: chọn nền tảng → xác nhận → tải |
| 6 | **Theo dõi boss trực tiếp** | gọi `/sv1/boss_status`, tự làm mới. **98 boss**, có tên, HP, bản đồ, giờ hồi sinh |
| 7 | **Bảng xếp hạng** | 2 bảng riêng: `bxh01` (sức mạnh + nhiệm vụ), `bxh02` (tài sản) |
| 8 | **Bình luận cộng đồng** | đăng, sửa, xoá, trả lời, thích. Chống spam: cần nhân vật >10 tỉ sức mạnh |
| 9 | **Trạng thái server** | số người online, tình trạng hoạt động |
| 10 | **Điều khoản sử dụng** | hộp thoại riêng |
| 11 | **Thông báo từ admin** | khối thông báo trên trang |
| 12 | **Liên kết cộng đồng** | Zalo, Fanpage |
| 13 | **Lối vào trang quản trị** | chỉ hiện với tài khoản admin |

**Bắt buộc thêm (chưa có, nên có):**
- Dòng cảnh báo *"Dành cho người chơi trên 12 tuổi. Chơi quá 180 phút mỗi ngày sẽ hại sức khỏe."*
- Trang/mục **Hướng dẫn cài đặt**
- Trang/mục **Danh sách giftcode**

---

# PHẦN 2 — Định hướng thị giác

## Bối cảnh cạnh tranh

Đã khảo sát 10 server đang sống (ChiBi, Nghiện Ngọc Rồng, Awaken, NCD, Hồi Sinh Blue...). Quy luật rõ: **gần như tất cả đều dùng chung nhóm từ hoài niệm** và **cùng một kiểu web đỏ-vàng-cam rực rỡ, dày đặc chữ, nhiều khung viền**.

→ **Cơ hội: làm ngược lại.** Tối, tĩnh, nhiều khoảng trống, ảnh lớn. Trông như trang của một studio game thật, không phải web server lậu.

## ⚠️ Bảng màu và chữ — xem file riêng

**Bảng màu ở phiên bản đầu của tài liệu này đã bị thay thế.** Bản đầu dùng tông tối lạnh hiện đại — sau khi mở tài nguyên game ra xem thì thấy **đá nhau hoàn toàn** với bản sắc thật của game (giấy cũ sepia + cam đỏ + xanh lá, phong cách cuộn giấy phương Đông).

👉 **Dùng bảng màu và font trong [`prompt-key-khoa-visual.md`](prompt-key-khoa-visual.md)** — được đo trực tiếp từ ảnh giao diện và logo game, không phải tự nghĩ ra.

Tóm tắt: nền nâu ấm `#2E1A0C`, khối nội dung giấy cũ `#E8D9BC`, nút bấm cam `#F06030`, đỏ `#D80000` chỉ cho trạng thái sống/khẩn. Tiêu đề **Baloo 2** (tròn mập giống logo game), nội dung **Be Vietnam Pro**. Cả hai đều có đủ dấu tiếng Việt.

## Nguyên tắc bố cục

- **Tối, nhiều khoảng trống.** Khoảng cách giữa các khối lớn hơn cảm giác "đủ".
- **Ảnh là vật liệu chính**, không phải trang trí. Ít nhất 3 khối có ảnh lớn.
- **Không lặp bố cục.** Đừng để 2 khối liền nhau cùng kiểu "chữ trái – ảnh phải".
- **Bo góc thống nhất**: 12px cho thẻ, 8px cho nút.
- **Một hành động chính mỗi khối.** Nút phụ phải trông rõ là phụ.

## Cấm — những thứ làm web trông rẻ tiền

- Gradient tím-xanh kiểu AI, viền phát sáng neon
- Ảnh nền lấp lánh + chữ vàng gradient
- Nhồi 8 thẻ vào một hàng
- Khung viền vàng bao quanh mọi thứ
- Ảnh nhân vật Dragon Ball cắt ghép thô, nền cưa
- Chữ chạy ngang (marquee), đếm ngược nhấp nháy
- Số liệu giả kiểu "99% hài lòng"

---

# PHẦN 3 — Bảy khối, bảy ảnh

**Quy tắc: mỗi khối một ảnh riêng, khổ ngang.** Đừng gộp cả trang vào một ảnh dài — không code được.

| # | Khối | Nhiệm vụ | Bố cục | Nền |
|---|---|---|---|---|
| 1 | **Hero** | Gây ấn tượng + 1 hành động | chữ dưới-trái đè lên ảnh | ảnh tràn viền, phủ tối |
| 2 | **Trạng thái server** | Chứng minh đang sống | dải ngang, số lớn | màu phẳng |
| 3 | **Tải game** | Chuyển đổi chính | 3 cột cân bằng | bề mặt tối + ảnh nhỏ |
| 4 | **Boss đang sống** | Điểm khác biệt | lưới bento không khe hở | tối, ảnh boss |
| 5 | **Bảng xếp hạng** | Tạo ganh đua | ảnh trái 40 / bảng phải 60 | ảnh xử lý 2 tông |
| 6 | **Cộng đồng** | Giữ chân | căn giữa, thoáng | màu phẳng, tối giản |
| 7 | **CTA cuối + chân trang** | Chốt | căn giữa, chữ lớn | ảnh tràn viền tối |

---

# PHẦN 4 — Prompt sinh ảnh

## Prompt nền tảng (dán trước mỗi prompt khối)

```
Website design comp, horizontal 16:9, dark premium game portal for a Vietnamese
Dragon Ball MMO called "Ngọc Rồng Kame".

ART DIRECTION:
- Deep dark mode. Background #0B0C0F, surfaces #171B24, hairlines #2A2F3A.
- Text #EDE9E3 primary, #8A8F9A secondary.
- Single accent: warm amber #F2A03D used sparingly for the primary action and
  key numbers only. Optional rare red #D6452F for "boss alive" state.
- NO purple/blue AI gradients, NO neon glow, NO glassmorphism, NO marquee.
- Typography: bold condensed sans display with FULL Vietnamese diacritics
  (Be Vietnam Pro style). Strong size contrast, tight tracking, short lines.
- Generous negative space. Calm, cinematic, editorial. Feels like a real game
  studio site, not a fan server page.
- Consistent 12px card radius, 8px button radius.
- Photography/illustration: anime-inspired Dragon Ball world, cinematic color
  grade matched to the palette, never bright rainbow.

Clean readable UI layout a developer can code from.
```

## Khối 1 — Hero

```
[PROMPT NỀN TẢNG]

SECTION: Hero. Giant statement scale, first viewport.
COMPOSITION: bottom-left text over a full-bleed cinematic background image.
BACKGROUND: wide anime landscape at dusk — floating rocky islands, distant
temple, a single glowing amber orb in the sky. Dark tonal overlay so text stays
fully readable.
CONTENT: small amber label "MÁY CHỦ CHÍNH THỨC" / huge two-line headline
"NGỌC RỒNG KAME" / one short supporting line / one solid amber primary button
"TẢI GAME" plus one ghost outline button "ĐĂNG KÝ".
Top bar: minimal logo left, 5 thin nav links right, no clutter.
Lots of empty sky above the text block.
```

## Khối 2 — Trạng thái server

```
[PROMPT NỀN TẢNG]

SECTION: Server status strip. Mini scale, calm, mostly negative space.
COMPOSITION: single horizontal row, centered.
BACKGROUND: flat solid #0B0C0F, one hairline above and below.
CONTENT: 4 oversized metrics in amber with tiny grey labels underneath —
"1.248 / Đang chơi", "98 / Boss đang sống", "x10 / Kinh nghiệm",
"Hoạt động / Trạng thái". Thin vertical divider lines between them.
Nothing else. Very airy.
```

## Khối 3 — Tải game

```
[PROMPT NỀN TẢNG]

SECTION: Download. Medium scale, conversion-focused.
COMPOSITION: centered heading, then three equal cards in a row.
BACKGROUND: surface #171B24 panel on dark page.
CONTENT: heading "Tải game" + one short line. Three cards: Windows, Android,
iOS. Each card has a platform glyph, platform name, file size caption, and a
button — the middle (Windows) card uses the solid amber button, the other two
use outline buttons. Equal spacing, generous padding inside cards.
```

## Khối 4 — Boss đang sống

```
[PROMPT NỀN TẢNG]

SECTION: Live boss tracker. Large, content-rich — the signature feature.
COMPOSITION: gapless bento grid, asymmetric tile sizes.
BACKGROUND: dark, with anime boss character art bleeding inside the larger tiles.
CONTENT: heading top-left "Boss đang sống" with a small live dot in red #D6452F.
Grid of boss tiles — one large feature tile with big character art, name, an
amber HP bar, map name; plus 5-6 smaller tiles each with a small portrait, name,
status pill ("Đang sống" red / "Hồi sinh 12:40" grey), and map. No gaps between
tiles. Data-dense but organized.
```

## Khối 5 — Bảng xếp hạng

```
[PROMPT NỀN TẢNG]

SECTION: Leaderboard. Medium editorial.
COMPOSITION: 40% image left, 60% ranking table right (inverted classic).
BACKGROUND: left side a duotone-treated anime warrior image graded to
amber/charcoal; right side flat dark surface.
CONTENT: heading "Bảng xếp hạng" with two small tab pills "Sức mạnh" and
"Tài sản". Ranking rows: rank number (top 3 in amber, oversized), small avatar,
player name, planet tag, power value right-aligned in monospace. Thin hairline
between rows. About 8 rows. Row 1 subtly emphasized.
```

## Khối 6 — Cộng đồng

```
[PROMPT NỀN TẢNG]

SECTION: Community. Mini minimalist, calm breathing section.
COMPOSITION: stacked center, very generous whitespace.
BACKGROUND: flat #0B0C0F, no image.
CONTENT: short centered heading "Cộng đồng người chơi", one supporting line,
two ghost outline buttons side by side "Nhóm Zalo" and "Fanpage". Below, a
single wide comment card showing avatar, player name, timestamp, short Vietnamese
comment text, and a small reply/like row. Just one card — restraint, not a wall.
```

## Khối 7 — CTA cuối + chân trang

```
[PROMPT NỀN TẢNG]

SECTION: Closing CTA and footer.
COMPOSITION: centered statement over full-bleed image, footer bar beneath.
BACKGROUND: dark cinematic anime sky with a single amber dragon ball orb,
heavily darkened.
CONTENT: large centered headline "Bắt đầu hành trình", one short line, one solid
amber button "TẢI GAME NGAY". Beneath, a slim footer: logo left, three link
columns center, and at the very bottom in small grey text the legal line
"Dành cho người chơi trên 12 tuổi. Chơi quá 180 phút mỗi ngày sẽ hại sức khỏe."
```

---

# PHẦN 5 — Thứ đổi được trong game

Chia theo **công sức bỏ ra**, không theo loại.

## ⚡ Đổi ngay, không build gì — chỉ thay file ảnh

Đây là phần **nhiều người không biết**. Thư mục `zz1\data\` là kho ảnh mà server **truyền thẳng cho client lúc chơi**. Thay file PNG là người chơi thấy ngay lần đăng nhập sau.

| Thư mục | Nội dung | Dung lượng |
|---|---|---|
| `data\icon\x1..x4` | **Toàn bộ icon vật phẩm** | 545 MB |
| `data\img_by_name\x1..x4` | Ảnh theo tên | 137 MB |
| `data\res\x1..x4` | Tài nguyên chung | 51 MB |
| `data\effect` + `data\effdata` | Hiệu ứng chiêu thức | 37 MB |
| `data\item_bg_temp` | Nền khung vật phẩm | 32 MB |
| `data\mob\x1..x4` | Ảnh quái | 19 MB |

⚠️ Phải thay **đủ cả 4 mức** `x1`–`x4` (4 độ phân giải), thiếu mức nào thì máy có màn hình tương ứng bị trống ảnh.

## ⚡ Đổi ngay qua trang quản trị hoặc database

Không cần build, không cần khởi động lại:

- **Vàng, ngọc, vật phẩm** của từng người chơi
- **Shop**: hàng bán, giá cả
- **Boss**: máu, sát thương, tên, vị trí, giờ hồi sinh
- **Tỉ lệ EXP và rơi đồ** (đang 10x / 5x)
- **Giftcode**: tạo, sửa, đặt phần thưởng
- **Vật phẩm**: chỉ số, tên, mô tả — bảng `item_template` có 1771 dòng
- **Bản đồ**: tên, quái xuất hiện — 183 bản đồ
- **NPC**: 111 NPC
- **Nhiệm vụ**, **thành tựu**, **giới hạn sức mạnh**
- **Thông báo trong game** (bảng `notify`)
- **Bảo trì**, kick người chơi, chế độ admin

## 🔧 Cần build lại server — khoảng 1–3 phút, người chơi văng ~1 phút

- Công thức tính toán (sát thương, ép đồ, tỉ lệ thành công)
- Cơ chế sự kiện mới
- Logic bang hội, võ đài, phó bản
- Nội dung chat của bot
- Danh hiệu, hiệu ứng
- Sửa lỗi
- **Chuyển mật khẩu sang mã hoá** (phải sửa cùng lúc cả bên web)

## 🐢 Cần build lại client Unity — 30–60 phút, người chơi phải tải lại

- **Tên game** trên cửa sổ và thanh tác vụ
- **Logo** (`Assets\Logo\1.png`) và **icon ứng dụng** (`Assets\Resources\appicon.jpg`)
- **Màn hình đăng nhập, đăng ký, chọn nhân vật**
- **Giao diện trong game**: bố cục nút, màu sắc, khung
- **Font chữ** (`Assets\Resources\fontsys\x2..x4`)
- **Địa chỉ server** — `Assets\Scripts\Game1\ServerListScreen.cs:43`
- Thêm màn hình / tính năng mới phía người chơi

**Quy tắc: gom thay đổi client lại làm một lần.** Mỗi lần build là mỗi lần bắt tất cả người chơi tải lại — làm nhiều lần liên tục là mất người chơi.

---

# PHẦN 6 — Cách dùng tài liệu này

1. Dán **prompt nền tảng + prompt khối 1** vào công cụ sinh ảnh → được ảnh Hero
2. Lặp lại cho khối 2 đến 7 → tổng **7 ảnh ngang**
3. Xem lại: 7 ảnh có trông như **cùng một website** không? Nếu lệch màu/lệch font thì sinh lại khối bị lệch
4. Gửi cả 7 ảnh cho Claude → code lại web theo ảnh, giữ nguyên 13 chức năng ở Phần 1

**Mẹo:** nếu công cụ sinh ảnh làm hỏng chữ tiếng Việt (rất hay xảy ra) thì kệ nó — ảnh chỉ để định hướng **bố cục, màu, tỉ lệ**. Chữ thật sẽ do code sinh ra, không lấy từ ảnh.
