export const meta = {
  name: 'nro-kame-website',
  description: 'Dung website Ngoc Rong Kame theo UI kit: nen tang dung chung, 6 trang song song, roi kiem chung doi khang',
  phases: [
    { title: 'Nen tang', detail: 'theme.css, layout.php, data.php, bang news + 30 bai mau' },
    { title: 'Trang', detail: '6 trang dung song song, moi agent mot file rieng' },
    { title: 'Kiem chung', detail: 'Kiem tra cu phap, HTTP, 13 chuc nang, du lieu that' },
  ],
}

const KIT = 'C:/Users/Admin/Music/brand_assets/ui_kit/nro-kame-complete-ui-kit'

const NEN = String.raw`
=== BOI CANH CHUNG — DOC KY TRUOC KHI LAM ===

Du an: website cho server game Ngoc Rong Online private ten "Ngọc Rồng Kame" (nrokame.com).
Toan bo NOI DUNG HIEN THI phai bang TIENG VIET CO DAU. Comment trong code cung tieng Viet.

## Duong dan
- Web goc:        C:/xampp/htdocs
- Thu muc game:   C:/xampp/htdocs/sv1
- Thu muc moi:    C:/xampp/htdocs/sv1/kame
- Bo UI kit goc:  C:/Users/Admin/Music/brand_assets/ui_kit/nro-kame-complete-ui-kit
    (co ui-kit.css, SITEMAP.md, README.md, mockups/ de tham khao bo cuc)
- Anh thuong hieu da lap: C:/xampp/htdocs/img/brand/
    logo.png (599x136) · logo-vuong.png · hero.jpg (1920x1080) · og.jpg
    giay-cu.jpg (texture lap lien mach) · band-boss.jpg · band-congdong.jpg
    fbcover.jpg · favicon-16/32/48/180/192/512.png
- PHP CLI kiem cu phap: C:/xampp/php/php.exe -l <file>
- MySQL: & C:/xampp/mysql/bin/mysql.exe -u root z --default-character-set=utf8mb4 -e "..."

## He mau va chu — LAY TU ui-kit.css, KHONG TU CHE
Bien CSS co san: --kame-orange-600 #f06030 (nut chinh), --kame-orange-700 #d94c22 (hover),
--kame-maroon-700 #600018 (vien/tieu de), --kame-maroon-900 #36000e (footer),
--kame-paper-100 #fff3d5 (nen the), --kame-paper-300 #e8d9bc, --kame-ink #28150f (chu),
--kame-gold #ffc83d, --kame-red #d80000 (boss song/canh bao), --kame-green #3f702a,
--kame-blue #174b89.
Font: tieu de "Roboto Condensed" 900 italic UPPERCASE, noi dung "Be Vietnam Pro".
Nap qua Google Fonts. TUYET DOI KHONG dung Bebas Neue (khong co dau tieng Viet).

Class co san trong ui-kit.css — DUNG LAI, dung viet lai:
kame-container kame-stack kame-cluster kame-grid kame-heading kame-muted
kame-button kame-card kame-badge kame-field kame-label kame-input kame-select
kame-textarea kame-help kame-error kame-tabs kame-tab kame-breadcrumb kame-alert
kame-accordion kame-table-wrap kame-table kame-modal-backdrop kame-modal
kame-pagination kame-page kame-dragonballs kame-dragonball kame-header kame-nav
kame-footer

## HOP DONG API — kame/data.php va kame/layout.php cung cap san
Trang KHONG duoc tu viet SQL. Goi cac ham nay:

layout.php:
  kame_head(string $title, string $desc = '')
  kame_header(string $active)   // trangchu | vutru | tintuc | bxh | hotro
  kame_footer()
  kame_auth_modals()            // 4 modal: dang nhap, dang ky, quen mk, doi mk

data.php:
  kame_server_status(): array          // online, accounts, players, ok
  kame_bosses(?int $limit = null): array  // name,status,state,hp,maxHp,mapName,respawnIn,avatar_id
  kame_boss_meta(): array              // ok, stale, age (giay)
  kame_top_power(int $limit = 8): array   // name, power, avatar_id, planet
  kame_top_wealth(int $limit = 8): array  // name, gold, avatar_id
  kame_news(int $limit, int $offset = 0): array
  kame_news_count(): int
  kame_news_one(string $slug): ?array
  kame_icon_url(int $id): string       // "/icon_asset.php?file=<id>.png"
  kame_placeholder(string $file, string $ratio = '16/9'): string

## LUAT BAT BUOC
1. TUYET DOI KHONG SUA cac file sau — vua duoc va lo hong mat khau trong ngay:
   sv1/dangky.php  sv1/dangnhap.php  sv1/doimatkhau.php  sv1/quenmatkhau.php
   sv1/download.php  sv1/config.php  sv1/boss_status.php  icon_asset.php  index.php
   Chi duoc GOI chung qua form/fetch.
2. KHONG bia so lieu. Cam: "128K" "215K" "78K" "350.842" "5.236" "12 may chu".
   So online / tong tai khoan lay that qua kame_server_status().
   Chua co du lieu thi ghi "Chưa có dữ liệu", khong dien so gia.
3. KHONG nhac Discord o bat ky dau. Cong dong: Zalo, Facebook, YouTube, TikTok —
   chi ghi TEN nen tang, khong ghi so nguoi theo doi.
4. KHONG dung anh nhan vat Dragon Ball. Anh boss/nguoi choi lay TU GAME qua kame_icon_url().
5. Moi trang co dong o chan: "Dành cho người chơi trên 12 tuổi. Chơi quá 180 phút mỗi ngày sẽ hại sức khỏe."
6. Chua co anh trong sv1/kame/assets/ — dung kame_placeholder() cho:
   the-gioi-traidat.jpg the-gioi-namec.jpg the-gioi-xayda.jpg (800x500)
   map-traidat.jpg map-namec.jpg map-xayda.jpg (1400x800)
   tin-mac-dinh.jpg (800x450)  band-taigame.jpg (1920x600)
7. Responsive 3 muc 1180 / 900 / 600 px.
8. Moi file PHP phai qua duoc php -l
9. CHI SUA/TAO file duoc giao. Khong dung vao file cua agent khac.
`

phase('Nen tang')

const nenTang = await agent(
  NEN + String.raw`
=== NHIEM VU: DUNG NEN TANG ===

Tao 4 file va 1 bang database. Day la nen mong cho 6 trang khac — lam that chac,
dat ten dung hop dong API o tren.

## 1. Chep ui-kit.css
Chep nguyen ven tu bo kit sang C:/xampp/htdocs/sv1/kame/ui-kit.css. KHONG SUA NOI DUNG.
Trong file co dong background url("assets/kame_tex_giay-cu_1024x1024_v1.jpg") — duong dan
nay sai. DUNG sua ui-kit.css; thay vao do ghi de bang body{} trong theme.css tro toi
/img/brand/giay-cu.jpg

## 2. Tao sv1/kame/theme.css
- Nap Roboto Condensed (700,900 italic) va Be Vietnam Pro (400,600,800) qua Google Fonts
- body{} ghi de background tro /img/brand/giay-cu.jpg
- Lop hero trang chu: anh /img/brand/hero.jpg, phu toi, chu duoi-trai
- Lop the boss: khung, thanh mau (do khi ACTIVE, xam khi AFK), avatar bo tron
- Lop hang bang xep hang: so thu hang, top 3 to hon mau vang
- Lop o giu cho anh (kame-ph): nen giay cu, vien dut, ten file o giua
- Lop khoi thong ke that (kame-stat)
- Lop dai canh bao 12 tuoi
- Lop menu mobile drawer khi duoi 900px
Dung bien CSS cua ui-kit.css, khong hardcode ma mau.

## 3. Tao sv1/kame/layout.php
kame_header($active): logo trai (/img/brand/logo.png cao 40px), menu giua 5 muc:
  Trang chủ /sv1/trangchu · Vũ trụ /sv1/vu-tru · Tin tức /sv1/tin-tuc
  · Bảng xếp hạng /sv1/bxh01 · Hỗ trợ /sv1/ho-tro
  Phai: nut ĐĂNG NHẬP (ghost, mo modal) + TẢI GAME (primary -> /sv1/tai-game).
  Neu da dang nhap: ten tai khoan + Đăng xuất; neu is_admin them link Quản trị -> /sv1/admin/panel.
  Muc active mau cam. Duoi 900px chuyen drawer.

kame_auth_modals(): 4 modal gui bang fetch toi /sv1/dangnhap /sv1/dangky
  /sv1/quenmatkhau /sv1/doimatkhau.
  DOC KY 4 file do TRUOC de biet chung nhan tham so gi va tra JSON kieu gi.
  Hien ket qua bang kame-alert.

kame_footer(): 4 cot (Khám phá / Hỗ trợ / Cộng đồng / Về chúng tôi), lien ket
  Zalo-Facebook-YouTube-TikTok chi ghi TEN. Duoi cung dong canh bao 12 tuoi
  + "© 2026 Ngọc Rồng Kame". KHONG Discord, KHONG newsletter, KHONG so lieu.

## 4. Tao sv1/kame/data.php
Cai dat day du hop dong API. Ghi chu:
- require_once config.php de dung $pdo
- kame_bosses(): doc C:/xampp/htdocs/sv1/boss_status.json, ghep avatar_id qua bang
  head_avatar. DOC sv1/boss_status.php de theo dung cach no lam (dung sua no).
- kame_boss_meta(): stale = true neu file cu hon 60 giay
- kame_server_status(): doc panel_server_status.json o game root. Tim theo thu tu
  C:/Users/Admin/Music/zz1 roi C:/Users/Administrator/Music/zz1.
  accounts = COUNT(*) FROM account; players = COUNT(*) FROM player.
  Khong doc duoc file thi ok=false.
- kame_top_power(): DOC sv1/bxh01.php de biet chinh xac chi so nao trong cot data_point
  (chuoi JSON) la suc manh — dung lai dung logic do. Ghep avatar qua player.head ->
  head_avatar.avatar_id. planet suy tu cot gender (0 Trái Đất, 1 Namếc, 2 Xayda) —
  kiem chung bang C:/Users/Admin/Music/zz1/src/consts/ConstPlayer.java
- kame_top_wealth(): tham khao sv1/bxh02.php
- kame_placeholder(): tra ve mot div class kame-ph co aspect-ratio va ten file o giua
- MOI ham bao ve bang try/catch, loi thi tra mang rong — KHONG lam trang trang

## 5. Tao bang news + gieo 30 bai mau
CREATE TABLE IF NOT EXISTS news (
  id INT AUTO_INCREMENT PRIMARY KEY,
  slug VARCHAR(160) NOT NULL UNIQUE,
  title VARCHAR(255) NOT NULL,
  excerpt TEXT,
  body LONGTEXT,
  image VARCHAR(255) DEFAULT NULL,
  tag VARCHAR(40) DEFAULT 'tin-tuc',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

Gieo dung 30 bai TIENG VIET CO DAU, noi dung HOP LY voi game NRO that:
- tag da dang: su-kien, cap-nhat, huong-dan, giftcode, bao-tri
- Dua tren dac diem THAT: EXP 10x, ti le roi do 5x, 98 boss, 183 ban do,
  3 hanh tinh Trái Đất/Namếc/Xayda, boss that (Ông già Noel, Thỏ Đại Gia, MR.GOLD),
  su kien Mabư 12h va 22h, Black Ball War 20h, Map War 21h29,
  Đại hội võ thuật moi gio (dang ky truoc phut 25)
- body 3-6 doan, dung the p / ul / h3
- created_at rai deu trong 60 ngay gan day
- image de NULL het
Kiem chung: SELECT COUNT(*) FROM news phai ra 30.

## KIEM CHUNG TRUOC KHI BAO XONG
- php -l sach ca 3 file php
- SELECT COUNT(*) FROM news = 30
- Viet file tam sv1/kame/_test.php goi thu MOI ham trong data.php, chay qua
  http://localhost/sv1/kame/_test.php, xac nhan khong loi va co du lieu that.
  XOA _test.php sau khi kiem xong.

Bao cao: file da tao, ket qua kiem chung, va CHINH XAC chu ky cac ham trong data.php
de 6 agent trang biet ma goi.
`,
  { label: 'nen-tang', phase: 'Nen tang' }
)

log('Nen tang xong. Bat dau dung 6 trang song song.')

phase('Trang')

const TRANG = [
  {
    key: 'trangchu',
    file: 'sv1/trangchu.php',
    mockup: '00-homepage.png',
    mo: String.raw`
VIET LAI HOAN TOAN sv1/trangchu.php (file cu 89KB). SAO LUU truoc thanh trangchu.php.cu

Bo cuc tu tren xuong:
1. HERO: nen /img/brand/hero.jpg, phu toi. Nhan nho cam "MÁY CHỦ CHÍNH THỨC",
   tieu de lon 2 dong, mot dong phu ngan, 2 nut: CHƠI NGAY (primary — mo modal dang ky
   neu chua dang nhap, tro /sv1/tai-game neu roi) va KHÁM PHÁ (secondary -> /sv1/vu-tru).
   Chu o duoi-trai (anh hero da chua san vung trong o do).
2. DAI THONG KE THAT: kame_server_status(). 3 o: Đang chơi (online), Tài khoản (accounts),
   Boss đang sống (dem tu kame_bosses loc status ACTIVE).
   ok=false thi hien "Máy chủ đang bảo trì". TUYET DOI khong bia so.
3. BAY VIEN NGOC RONG: dung kame-dragonballs / kame-dragonball co san, nhan 1 SAO..7 SAO.
4. BOSS DANG SONG: tieu de + cham do. Luoi 6-8 the boss dau tien tu kame_bosses().
   Moi the: avatar (kame_icon_url), ten, thanh mau hp/maxHp, ten ban do,
   trang thai (Đang sống do / Hồi sinh sau X phút). Nut xem tat ca 98 boss.
   Tu lam moi bang fetch('/sv1/boss_status') moi 10 giay.
   Neu kame_boss_meta stale thi hien nhan "Dữ liệu chưa cập nhật".
5. TIN TUC MOI NHAT: 4 bai tu kame_news(4). Anh dung kame_placeholder('tin-mac-dinh.jpg','800/450').
   Link -> /sv1/tin-tuc?bai=<slug>. Nut Xem tất cả -> /sv1/tin-tuc.
6. BA THE GIOI: 3 the. Anh dung kame_placeholder cho the-gioi-traidat.jpg,
   the-gioi-namec.jpg, the-gioi-xayda.jpg (ti le 800/500).
   Mo ta phai DUNG game that: Trái Đất can bang, co Dịch chuyển tức thời;
   Namếc hoi mau, Trị thương va Đẻ trứng; Xayda burst manh, Biến hình va Tự phát nổ.
7. TOP NGUOI CHOI: kame_top_power(8). Hang, avatar that, ten, suc manh.
   Top 3 to hon mau vang. Link -> /sv1/bxh01. Rong thi "Chưa có dữ liệu — hãy là người đầu tiên".
8. TAI GAME: 3 nut Windows / Android / iOS -> /sv1/tai-game.
9. CONG DONG: Zalo, Facebook, YouTube, TikTok — CHI TEN, khong so.
Ket thuc bang kame_footer().
`,
  },
  {
    key: 'vutru',
    file: 'sv1/vu-tru.php',
    mockup: '05-universe-characters.png',
    mo: String.raw`
TAO MOI sv1/vu-tru.php. Thay HET noi dung hu cau trong mockup bang du lieu THAT.

1. Tieu de + breadcrumb.
2. TAB 3 HANH TINH (kame-tabs): Trái Đất / Namếc / Xayda. Moi tab:
   - Anh ban do: kame_placeholder('map-traidat.jpg','1400/800') (tuong tu namec, xayda)
   - Dia danh THAT lay tu bang map_template:
     mysql.exe -u root z -e "SELECT id,NAME FROM map_template ORDER BY id"
     Chon 6-8 dia danh tieu bieu moi hanh tinh, ghi ten DUNG NHU TRONG DB.
   - Mo ta hanh tinh + bo ky nang THAT: doc bang skill_template va
     C:/Users/Admin/Music/zz1/src/consts/ConstPlayer.java
3. THAY "Các chiến binh nổi bật" (8 nhan vat hu cau) bang "CHIẾN BINH MẠNH NHẤT":
   top 8 nguoi choi that tu kame_top_power(8). Avatar that, ten, hanh tinh, suc manh.
   Rong thi "Chưa có dữ liệu".
4. THAY "Boss huyền thoại — Vua Hắc Ám" (hu cau) bang "BOSS TRONG THẾ GIỚI KAME":
   bang day du 98 boss tu kame_bosses(). Dung kame-table-wrap / kame-table,
   co o tim kiem loc theo ten bang JavaScript thuan.
   Cot: Ảnh · Tên · Trạng thái · Máu · Bản đồ · Hồi sinh.
5. HANH TRINH 6 buoc dua tren nhiem vu THAT: Khoi dau (mo ruong, thu dau than) ->
   Hoc bay -> Bai su -> Gia nhap bang hoi -> Dai chien Fide/Xen -> End game 5 ty suc manh.
6. CTA cuoi + kame_footer().
`,
  },
  {
    key: 'tintuc',
    file: 'sv1/tin-tuc.php',
    mockup: '02-news-listing.png',
    mo: String.raw`
TAO MOI sv1/tin-tuc.php — MOT file lo hai viec.
Tham khao them mockup 03-news-detail.png cho phan chi tiet.

KHONG co ?bai= -> DANH SACH:
- Bai noi bat: bai moi nhat, the lon
- Tab loc theo tag (kame-tabs): Tất cả · Sự kiện · Cập nhật · Hướng dẫn · Giftcode · Bảo trì
  loc bang ?tag=
- O tim kiem ?q= tim trong title va excerpt
- Luoi the bai viet, 9 bai moi trang, dung kame_news(limit,offset) va kame_news_count()
- Phan trang dung kame-pagination / kame-page
- Anh moi bai dung kame_placeholder('tin-mac-dinh.jpg','800/450')
- Rong thi kame-alert bao "Chưa có bài viết"

Co ?bai=<slug> -> CHI TIET:
- kame_news_one(slug). Khong thay -> kame-alert--danger "Không tìm thấy bài viết"
  + nut quay lai. KHONG duoc trang trang.
- Breadcrumb, tieu de, ngay dang, tag, anh dau bai, noi dung body (chua HTML, in nguyen)
- Cuoi bai: 3 bai lien quan cung tag
- KHONG lam binh luan o day (binh luan nam o trang ho-tro)

Bao ve: escape moi $_GET bang htmlspecialchars khi in ra.
`,
  },
  {
    key: 'taigame',
    file: 'sv1/tai-game.php',
    mockup: '04-download.png',
    mo: String.raw`
TAO MOI sv1/tai-game.php.

1. Bang dau trang: kame_placeholder('band-taigame.jpg','1920/600') lam nen, tieu de de len.
2. BA THE NEN TANG Windows / Android / iOS.
   QUAN TRONG — DOC sv1/download.php TRUOC. No chi cho phep 3 ten file:
   NRO.apk, PC.rar, IPA.ipa — va thu muc sv1/down/ HIEN DANG TRONG nen se tra 404.
   Moi the phai kiem tra file co ton tai trong C:/xampp/htdocs/sv1/down/ khong:
   - Co  -> nut TẢI NGAY tro /sv1/download?file=<ten>, hien dung luong that
   - Chua -> nut mo xam, chu "Sắp có", kem kame-alert--info "Bản cài đang được hoàn thiện"
   TUYET DOI khong lam nut chet 404.
3. CAU HINH DE NGHI: bang don gian (Windows 10 64-bit, RAM 4GB, 1GB trong).
   Ghi ro day la con so uoc tinh.
4. BA BUOC CAI DAT: tai -> giai nen -> chay, dung kame-card danh so.
5. FAQ dung kame-accordion, 5-6 cau thuc te: khong tai duoc, may bao virus, quen mat khau,
   choi tren dien thoai, cau hinh yeu, loi ket noi. Cau tra loi phai HUU ICH va that.
6. Nut phu: Đăng ký tài khoản mo modal dang ky.
7. kame_footer().
`,
  },
  {
    key: 'hotro',
    file: 'sv1/ho-tro.php',
    mockup: '06-community-support.png',
    mo: String.raw`
TAO MOI sv1/ho-tro.php. BO form gui ticket va bang theo doi ticket (backend khong co).

1. Bang dau trang: /img/brand/band-congdong.jpg lam nen, tieu de de len.
2. THE CONG DONG: Zalo · Facebook · YouTube · TikTok. Moi the chi TEN nen tang
   + mot cau mo ta ngan + nut Tham gia (link tam nrokame.com/lien-he).
   KHONG so nguoi theo doi. KHONG Discord.
3. TIM KIEM TRO GIUP: o tim kiem loc FAQ ben duoi bang JavaScript thuan.
4. FAQ theo nhom dung kame-accordion, it nhat 4 nhom moi nhom 4-6 cau:
   - Tài khoản: dang ky, quen mat khau, doi mat khau, bi khoa
   - Cài đặt & Kết nối: tai game, khong vao duoc, mat ket noi, cau hinh
   - Trong game: len suc manh, nhiem vu, bang hoi, boss, su kien theo gio
   - Nạp & Vật phẩm: cach nap, khong nhan duoc do, giftcode
   Noi dung DUNG game that. Moc gio that: Mabư 12h va 22h, Black Ball War 20h-21h,
   Map War 21h29, Đại hội võ thuật moi gio (dang ky truoc phut 25).
   EXP 10x, roi do 5x, 98 boss, 183 ban do.
5. BINH LUAN CONG DONG: lay logic binh luan cu tu ban sao luu sv1/trangchu.php.cu
   (agent trang chu tao) hoac tu D:/NRO_Backup. Bang site_comments + site_comment_likes,
   chong spam: can nhan vat tren 10 ty suc manh, 3 bai/ngay, 100 tra loi/ngay.
   Neu khong tim duoc file cu thi doc truc tiep bang site_comments trong DB va cai dat lai.
   GIU NGUYEN ten bang va cot, khong doi schema.
6. kame_footer().
`,
  },
  {
    key: 'bxh',
    file: 'sv1/bxh01.php va sv1/bxh02.php',
    mockup: '01-ui-kit.png',
    mo: String.raw`
KHOAC GIAO DIEN MOI cho sv1/bxh01.php va sv1/bxh02.php.

QUAN TRONG: GIU NGUYEN toan bo LOGIC truy van va tinh toan hien co. Doc ky truoc khi sua.
Chi thay phan HTML/CSS xuat ra. Sao luu bxh01.php.cu va bxh02.php.cu truoc khi sua.

Bo cuc moi:
- kame_head / kame_header / kame_footer + ui-kit.css + theme.css, bo CSS cu
- Tieu de + breadcrumb
- kame-tabs chuyen qua lai: "Sức mạnh" (bxh01) va "Tài sản" (bxh02)
- Bang dung kame-table-wrap / kame-table
- Cot: Hạng · Ảnh · Tên · Hành tinh · Chỉ số
- Top 3: so hang to hon mau vang, hang co nen nhat hon
- Anh dai dien ghep player.head -> head_avatar.avatar_id -> kame_icon_url()
  (neu logic cu chua co avatar thi THEM VAO, day la cai tien duoc phep)
- Rong: "Chưa có dữ liệu — hãy là người đầu tiên"

Luu y: bxh01/bxh02 hien SELECT toan bo bang player roi sap xep bang PHP — voi vai nghin
nguoi choi se rat cham. Neu sua duoc sang ORDER BY + LIMIT trong SQL ma khong lam sai
ket qua thi lam; neu logic qua phuc tap (data_point la JSON) thi GIU NGUYEN va ghi
// TODO hieu nang vao code.
`,
  },
]

const ketQua = await parallel(TRANG.map(t => () =>
  agent(
    NEN +
      '\n=== KET QUA TU AGENT NEN TANG ===\n' + nenTang +
      '\n\n=== NHIEM VU CUA BAN ===\nFile ban phu trach: ' + t.file +
      '\nMockup tham khao: ' + KIT + '/mockups/' + t.mockup +
      '\n(HAY MO ANH MOCKUP RA XEM truoc khi code — de nam bo cuc.)\n' +
      t.mo + String.raw`
=== TRUOC KHI BAO XONG, TU KIEM ===
1. C:/xampp/php/php.exe -l tren file cua ban — phai sach
2. Truy cap that qua http://localhost/... — HTTP 200, khong Warning/Notice
3. Doc HTML sinh ra: khong con so lieu bia, khong con chu discord, khong con ten
   nhan vat hu cau (Son Kai, Lina, Moru, Kaiken, Zeri, Zoru, Nazu, Mira, Vua Hắc Ám)
4. Co dong canh bao 12 tuoi o chan
5. Chu tieng Viet hien dung dau, khong o vuong

Bao cao ngan: file da tao/sua, ket qua 5 muc kiem tren, va nhung gi ban KHONG lam duoc.
`,
    { label: 'trang:' + t.key, phase: 'Trang' }
  )
))

log('6 trang xong. Bat dau kiem chung doi khang.')

phase('Kiem chung')

const KIEM_SCHEMA = {
  type: 'object',
  properties: {
    loi: {
      type: 'array',
      items: {
        type: 'object',
        properties: {
          file: { type: 'string' },
          muc_do: { type: 'string', description: 'nghiem-trong | trung-binh | nho' },
          mo_ta: { type: 'string', description: 'Loi la gi, bang tieng Viet' },
          bang_chung: { type: 'string', description: 'Trich dan chinh xac dong code hoac ket qua lenh' },
          cach_sua: { type: 'string' },
        },
        required: ['file', 'muc_do', 'mo_ta', 'bang_chung'],
      },
    },
    tom_tat: { type: 'string' },
  },
  required: ['loi', 'tom_tat'],
}

const GOC_KIEM = [
  {
    key: 'chuc-nang',
    prompt: String.raw`Kiem tra 13 CHUC NANG BAT BUOC con day du va HOAT DONG khong:
dang nhap · dang ky · quen mat khau · doi mat khau · tai game (3 nen tang, luong 2 buoc)
· theo doi boss truc tiep · 2 bang xep hang · binh luan cong dong · trang thai server
· dieu khoan su dung · thong bao admin · lien ket cong dong · loi vao trang quan tri.

Voi moi chuc nang: tim trong code xem no o dau, goi endpoint nao, va THU THAT bang
Invoke-WebRequest neu duoc. Dac biet kiem 4 modal tai khoan co gui dung tham so ma
sv1/dangnhap.php va sv1/dangky.php mong doi khong — DOC 2 file do de doi chieu.
Bao cao chuc nang nao THIEU hoac GOI SAI.`,
  },
  {
    key: 'du-lieu-that',
    prompt: String.raw`San tim SO LIEU BIA va NOI DUNG HU CAU con sot.
Grep moi file moi trong C:/xampp/htdocs/sv1 tim: 128K 215K 78K 350.842 350842
5.236 5236 128.450 128450 "12 máy chủ" discord Discord "Son Kai" Lina Moru Kaiken
Zeri Zoru Nazu Mira "Vua Hắc Ám" "Đại Chiến Võ Đạo" "100.000.000 VND".
Cung kiem: cac con so hien tren trang co that su lay tu database khong hay bi hardcode.
Truy cap that cac trang va doc HTML tra ve.`,
  },
  {
    key: 'khong-pha-hong',
    prompt: String.raw`Xac minh KHONG file cam nao bi sua.
Phai giu nguyen: sv1/dangky.php sv1/dangnhap.php sv1/doimatkhau.php sv1/quenmatkhau.php
sv1/download.php sv1/config.php sv1/boss_status.php icon_asset.php index.php

So sanh voi ban sao luu D:/NRO_Backup/2026-07-21_truoc-hash-matkhau (co config.php,
dangky.php, dangnhap.php, doimatkhau.php, account.php) bang Get-FileHash.
File nao khac ban sao luu -> loi NGHIEM TRONG.
Luu y: config.php va cac file tai khoan DA duoc sua hop le trong ngay de them ma hoa
mat khau PBKDF2 — hay so sanh voi ban HIEN TAI truoc khi chay workflow, khong phai
ban goc. Neu khong chac, kiem bang cach: grep kame_password_verify va kame_password_hash
van con trong dangnhap.php / dangky.php / doimatkhau.php.
Cung thu dang nhap that qua HTTP voi tai khoan "1" mat khau "1".`,
  },
  {
    key: 'ky-thuat',
    prompt: String.raw`Kiem tra CHAT LUONG KY THUAT.
- Chay C:/xampp/php/php.exe -l tren MOI file .php trong sv1/ va sv1/kame/
- Truy cap tung trang qua HTTP, bat moi PHP Warning/Notice/Deprecated trong HTML tra ve
- SQL injection: moi $_GET/$_POST dung trong truy van phai qua prepared statement
- XSS: moi bien in ra HTML phai qua htmlspecialchars
- Duong dan anh: moi img src co tro toi file ton tai that khong
- Font: co nap Roboto Condensed va Be Vietnam Pro khong, co lo dung Bebas Neue khong
- Responsive: co media query cho 900px va 600px khong`,
  },
]

const kiemTra = await parallel(GOC_KIEM.map(g => () =>
  agent(
    NEN + String.raw`

=== NHIEM VU: KIEM CHUNG DOI KHANG ===
Ban la nguoi kiem tra doc lap. 6 agent khac vua dung xong website. Viec cua ban la
TIM LOI, khong phai khen. Mac dinh gia dinh co loi cho toi khi chung minh duoc la khong.

Chi bao cao loi ban CHUNG MINH duoc bang bang chung cu the (trich dong code, ket qua
lenh, hoac noi dung HTML tra ve). KHONG doan mo ho.

Goc kiem cua ban:
` + g.prompt + String.raw`

Cac trang can kiem: /sv1/trangchu /sv1/vu-tru /sv1/tin-tuc /sv1/tai-game /sv1/ho-tro
/sv1/bxh01 /sv1/bxh02

Tra ve theo dung schema. muc_do: "nghiem-trong" neu lam vo chuc nang hoac lo lot du lieu;
"trung-binh" neu sai yeu cau nhung khong vo; "nho" neu chi la thieu sot nho.`,
    { label: 'kiem:' + g.key, phase: 'Kiem chung', schema: KIEM_SCHEMA }
  )
))

const tatCaLoi = kiemTra.filter(Boolean).flatMap(k => k.loi || [])
const nghiemTrong = tatCaLoi.filter(l => l.muc_do === 'nghiem-trong')

log('Kiem chung xong: ' + tatCaLoi.length + ' loi, trong do ' + nghiemTrong.length + ' nghiem trong')

return {
  so_loi: tatCaLoi.length,
  nghiem_trong: nghiemTrong.length,
  loi: tatCaLoi,
  tom_tat_kiem: kiemTra.filter(Boolean).map(k => k.tom_tat),
}
