# ĐỢT 1 — HƯỚNG DẪN TẠO ẢNH (dành cho AI tạo ảnh)

> **Đọc toàn bộ Phần A trước khi tạo bất kỳ ảnh nào.**
> Sau đó làm tuần tự Nhóm 1 → Nhóm 8. Mỗi dòng trong bảng là một ảnh.
> Tự đặt tên file đúng cột `Tên file`. Xuất PNG nền trong suốt.
> Khi xong tất cả, nén thành 1 file zip giữ nguyên cấu trúc thư mục theo nhóm.

---

# PHẦN A — LUẬT CHUNG (BẮT BUỘC)

## A1. Khối phong cách — dán vào MỌI prompt

```
STYLE: 2D anime game art, bold thick dark outline (2-3px at final size),
flat cel shading with exactly 3 tone steps (base / shadow / highlight),
limited palette, high contrast, clean readable silhouette at small size,
centered composition, transparent background, no drop shadow on ground,
no text, no watermark, no signature, no border frame.
INSPIRATION: ninja anime of the 2000s — original designs, NOT copies of
existing copyrighted characters.
```

## A2. Bảng màu chuẩn — dùng xuyên suốt

| Vai trò | Mã màu |
|---|---|
| Nền tối chủ đạo | `#1B1B23` |
| Nền phụ | `#2E2E3A` |
| Viền sáng | `#4A4A5C` |
| Chữ chính | `#F2EFE6` |
| Cam Lá (màu thương hiệu) | `#E8722C` |
| Cam nhạt | `#F5A65B` |
| Xanh chakra | `#3FA9C9` |
| Đỏ máu | `#C0392B` |
| Vàng kim (quý hiếm) | `#E8B33C` |
| Tím hiếm | `#8E5BC4` |
| Xám kim loại | `#8A8F98` |

**Quy tắc màu theo độ hiếm** (viền quanh vật phẩm):
Thường = xám `#8A8F98` · Hiếm = xanh `#3FA9C9` · Quý = tím `#8E5BC4` · Huyền thoại = vàng `#E8B33C`

## A3. Màu riêng từng làng

| Làng | Màu chính | Màu phụ |
|---|---|---|
| Làng Lá | `#E8722C` cam | `#3D7A3D` xanh lá |
| Làng Cát | `#D9B26A` be cát | `#8C6239` nâu |
| Làng Sương Mù | `#7FA8B8` xanh xám | `#2E4A56` xanh đậm |
| Làng Đá | `#8A7A66` nâu đá | `#4A4238` nâu tối |
| Làng Mây | `#C9C9D6` bạc | `#F0D14A` vàng sét |

> Đợt 1 **chỉ dùng màu Làng Lá**. Bốn làng còn lại để đợt 2.

## A4. Quy tắc kích thước

Tạo ảnh **đúng bằng kích thước ghi ở cột `Cỡ`**. Nếu công cụ không tạo được cỡ nhỏ, hãy tạo ở bội số 8 lần rồi thu nhỏ bằng thuật toán **nearest-neighbor** (không dùng làm mờ).

| Loại | Cỡ |
|---|---|
| Biểu tượng nhỏ | 64×64 |
| Biểu tượng vật phẩm | 64×64 |
| Vũ khí dài | 96×48 |
| Trang phục (mỗi phần) | 128×160 |
| Quái | 96×96 |
| Boss | 160×160 |
| NPC | 96×128 |
| Chân dung hội thoại | 200×200 |
| Hiệu ứng (tấm nhiều khung) | 240×240 mỗi khung |
| Mảnh ghép bản đồ | 96×96 |
| Nút giao diện | 96×48 |
| Khung giao diện | ghi riêng từng dòng |

## A5. Quy tắc đặt tên file

```
<nhóm>_<mã>_x4.png
```

Tiền tố nhóm: `ui_` giao diện · `vp_` vật phẩm · `vk_` vũ khí · `tp_` trang phục · `quai_` quái · `boss_` boss · `npc_` NPC · `cd_` chân dung · `hu_` hiệu ứng · `map_` mảnh ghép bản đồ

Thư mục trong zip:
```
01-giao-dien/  02-vat-pham/  03-trang-phuc/  04-vu-khi/
05-quai/  06-boss/  07-npc/  08-hieu-ung/  09-ban-do/
```

## A6. Cấm tuyệt đối

```
NEGATIVE (thêm vào mọi prompt):
photorealistic, 3d render, blurry, gradient background, white background,
multiple objects in one image (unless the row says sprite sheet),
text, letters, numbers, logo, watermark, signature, frame border,
guns, modern clothing, western fantasy armor, drop shadow on floor.
```

## A7. Thứ tự làm — QUAN TRỌNG

Làm **Nhóm 1 (Giao diện) trước tiên** và **dừng lại chờ duyệt** sau khi xong 10 ảnh đầu. Giao diện định ra phong cách cho mọi thứ còn lại — sai ở đây thì 350 ảnh sau sai theo.

---

# PHẦN B — DANH SÁCH ẢNH

## NHÓM 1 — GIAO DIỆN (60 ảnh) → `01-giao-dien/`

**Khối phụ cho nhóm này:** `UI element, clean vector-like edges, subtle scroll/parchment and dark lacquered wood motif, ninja scroll aesthetic, symmetrical.`

### 1A. Mười ảnh đầu — LÀM TRƯỚC RỒI DỪNG CHỜ DUYỆT

| # | Tên file | Cỡ | SUBJECT |
|---|---|---|---|
| 1 | `ui_khung_hoithoai_x4.png` | 480×160 | dialogue box, dark lacquered wood panel with orange trim, aged paper inner area, small scroll-roll caps at both ends, 9-slice friendly (flat edges) |
| 2 | `ui_nut_thuong_x4.png` | 96×48 | button, dark wood with orange border, slight bevel, empty center |
| 3 | `ui_nut_nhan_x4.png` | 96×48 | same button but pressed-in state, darker, inset shadow |
| 4 | `ui_nut_khoa_x4.png` | 96×48 | same button but disabled, desaturated grey, small padlock at right |
| 5 | `ui_thanh_theluc_x4.png` | 240×32 | horizontal stat bar, empty dark socket with orange-red fill from left, rounded caps |
| 6 | `ui_thanh_chakra_x4.png` | 240×32 | same bar shape, cyan blue fill |
| 7 | `ui_o_vatpham_x4.png` | 64×64 | empty inventory slot, dark recessed square, thin lighter inner border |
| 8 | `ui_o_dangchon_x4.png` | 64×64 | same slot with bright orange glowing border, selected state |
| 9 | `ui_khung_avatar_x4.png` | 96×96 | circular portrait frame, rope-braid ring, dark wood outer, empty center |
| 10 | `ui_logo_game_x4.png` | 512×256 | game logo mark: a stylised leaf-and-spiral emblem over bold brush-stroke shapes, orange and dark, no letters |

**→ DỪNG. Gửi 10 ảnh này để duyệt phong cách trước khi làm tiếp.**

### 1B. Phần còn lại của giao diện

| # | Tên file | Cỡ | SUBJECT |
|---|---|---|---|
| 11 | `ui_thanh_kinhnghiem_x4.png` | 240×24 | thin xp bar, gold fill, dark socket |
| 12 | `ui_khung_tuido_x4.png` | 480×480 | inventory panel, dark wood frame, paper interior, grid guides |
| 13 | `ui_khung_ruong_x4.png` | 480×480 | storage chest panel, heavier iron corners |
| 14 | `ui_khung_shop_x4.png` | 480×400 | shop panel, hanging cloth banner at top |
| 15 | `ui_khung_nhiemvu_x4.png` | 400×480 | quest panel shaped like an unrolled scroll |
| 16 | `ui_khung_thongbao_x4.png` | 400×120 | small notification banner, orange edge glow |
| 17 | `ui_khung_canhbao_x4.png` | 400×120 | warning banner, red edge glow, small hazard mark |
| 18 | `ui_nut_dong_x4.png` | 64×64 | close button, round dark wood with an X mark |
| 19 | `ui_nut_caidat_x4.png` | 64×64 | settings button, round, gear-shuriken hybrid mark |
| 20 | `ui_nut_chat_x4.png` | 64×64 | chat button, speech bubble mark |
| 21 | `ui_nut_bando_x4.png` | 64×64 | map button, folded map mark |
| 22 | `ui_nut_tuido_x4.png` | 64×64 | bag button, ninja pouch mark |
| 23 | `ui_nut_kynang_x4.png` | 64×64 | skill button, hand-seal mark |
| 24 | `ui_nut_nhiemvu_x4.png` | 64×64 | quest button, small scroll mark |
| 25 | `ui_nut_toidoi_x4.png` | 64×64 | party button, three small figures mark |
| 26 | `ui_nut_xephang_x4.png` | 64×64 | ranking button, podium mark |
| 27 | `ui_nut_thu_x4.png` | 64×64 | mail button, sealed letter mark |
| 28 | `ui_muiten_len_x4.png` | 48×48 | up arrow, thick, orange |
| 29 | `ui_muiten_xuong_x4.png` | 48×48 | down arrow |
| 30 | `ui_muiten_trai_x4.png` | 48×48 | left arrow |
| 31 | `ui_muiten_phai_x4.png` | 48×48 | right arrow |
| 32 | `ui_tien_ryo_x4.png` | 64×64 | currency icon: an old oval gold coin with a square hole, edge-on stack of two |
| 33 | `ui_tien_ngoc_x4.png` | 64×64 | currency icon: a cut cyan gemstone, faceted |
| 34 | `ui_tien_xu_x4.png` | 64×64 | currency icon: a small purple event token with a spiral mark |
| 35 | `ui_tien_thoivang_x4.png` | 64×64 | currency icon: a gold ingot bar |
| 36 | `ui_lang_la_x4.png` | 64×64 | village emblem: stylised leaf on a metal plate |
| 37 | `ui_lang_cat_x4.png` | 64×64 | village emblem: hourglass-spiral on a metal plate |
| 38 | `ui_lang_suong_x4.png` | 64×64 | village emblem: four wavy mist lines on a metal plate |
| 39 | `ui_lang_da_x4.png` | 64×64 | village emblem: two rock slabs on a metal plate |
| 40 | `ui_lang_may_x4.png` | 64×64 | village emblem: cloud shape on a metal plate |
| 41 | `ui_he_hoa_x4.png` | 64×64 | element icon: stylised flame, red-orange |
| 42 | `ui_he_phong_x4.png` | 64×64 | element icon: swirling wind curl, pale green |
| 43 | `ui_he_thuy_x4.png` | 64×64 | element icon: water droplet with wave, blue |
| 44 | `ui_he_tho_x4.png` | 64×64 | element icon: rock cluster, brown |
| 45 | `ui_he_loi_x4.png` | 64×64 | element icon: lightning bolt, yellow-white |
| 46 | `ui_bac_hocvien_x4.png` | 64×64 | rank badge tier 1: plain cloth headband strip, no metal |
| 47 | `ui_bac_genin_x4.png` | 64×64 | rank badge tier 2: metal plate, one notch |
| 48 | `ui_bac_chunin_x4.png` | 64×64 | rank badge tier 3: metal plate with green vest collar behind |
| 49 | `ui_bac_jonin_x4.png` | 64×64 | rank badge tier 4: metal plate, silver trim |
| 50 | `ui_bac_tinhanh_x4.png` | 64×64 | rank badge tier 5: metal plate with two side wings |
| 51 | `ui_bac_anbu_x4.png` | 64×64 | rank badge tier 6: white animal mask silhouette |
| 52 | `ui_bac_kage_x4.png` | 64×64 | rank badge tier 7: conical hat silhouette, gold |
| 53 | `ui_bac_hiennhan_x4.png` | 64×64 | rank badge tier 8: toad-eye marking motif, orange |
| 54 | `ui_bac_lucdao_x4.png` | 64×64 | rank badge tier 9: nine magatama in a ring, gold on black |
| 55 | `ui_thanh_nap_x4.png` | 400×48 | loading bar, dark socket with orange fill and a small running shuriken at the fill edge |
| 56 | `ui_joystick_nen_x4.png` | 192×192 | virtual joystick base, dark translucent ring with subtle rope texture |
| 57 | `ui_joystick_num_x4.png` | 96×96 | virtual joystick knob, orange dome |
| 58 | `ui_nut_danh_x4.png` | 128×128 | big round attack button, orange, a shuriken mark in the middle |
| 59 | `ui_nen_dangnhap_x4.png` | 1024×576 | login background: a quiet ninja village at dusk seen from a hillside, silhouetted rooftops, large moon, warm orange sky, heavy stylisation, low detail |
| 60 | `ui_nen_chonlang_x4.png` | 1024×576 | village-select background: an unrolled map of five nations on dark wood table, lantern light, low detail |

---

## NHÓM 2 — VẬT PHẨM (120 ảnh) → `02-vat-pham/`

**Khối phụ:** `single game item icon, 3/4 top view, floating, centered, no hands, no background.`

### 2A. Trang bị Làng Lá — 5 bộ phận × 8 cấp = 40

Mẫu tên: `vp_<bộphận>_t<cấp>_x4.png`, cỡ **64×64**.

Cấp 1→8 tăng dần độ tinh xảo: `vải thô` → `vải bền` → `da nhẹ` → `lưới thép` → `giáp nhẹ` → `giáp Jonin` → `giáp ANBU` → `áo choàng Kage`.

| Bộ phận | Tên file | SUBJECT gốc (thêm bậc chất liệu theo cấp) |
|---|---|---|
| Áo | `vp_ao_t1..t8_x4.png` | a ninja jacket, orange and dark blue accents, high collar |
| Quần | `vp_quan_t1..t8_x4.png` | ninja trousers with bandage wraps at the ankles |
| Găng | `vp_gang_t1..t8_x4.png` | a pair of fingerless ninja gloves with a metal back plate |
| Giày | `vp_giay_t1..t8_x4.png` | open-toe ninja sandals, dark blue |
| Băng trán | `vp_bangtran_t1..t8_x4.png` | forehead protector, metal plate with leaf emblem on a cloth strip |

### 2B. Tiêu hao — 12

| Tên file | SUBJECT |
|---|---|
| `vp_vien_theluc_1_x4.png` | small red pill in a paper wrapper |
| `vp_vien_theluc_2_x4.png` | medium red pill, glass vial |
| `vp_vien_theluc_3_x4.png` | large red pill, ornate jar |
| `vp_vien_chakra_1_x4.png` | small blue pill in paper |
| `vp_vien_chakra_2_x4.png` | medium blue pill, glass vial |
| `vp_vien_chakra_3_x4.png` | large blue pill, ornate jar |
| `vp_thuoc_giai_x4.png` | green antidote bottle with cork |
| `vp_bang_gac_x4.png` | rolled white bandage with a red cross-less plain tag |
| `vp_binh_luong_x4.png` | soldier pill: dark round pellet in a small tin |
| `vp_mi_ramen_x4.png` | steaming bowl of ramen with narutomaki slice |
| `vp_dango_x4.png` | three-ball dango skewer |
| `vp_binh_nuoc_x4.png` | leather water flask |

### 2C. Nguyên liệu — 15

| Tên file | SUBJECT |
|---|---|
| `vp_manh_thep_x4.png` | jagged steel shard |
| `vp_cuon_trang_x4.png` | blank paper scroll, rolled |
| `vp_muc_x4.png` | small ink stone with brush |
| `vp_long_chim_x4.png` | dark feather |
| `vp_da_chakra_tho_x4.png` | raw dull blue crystal chunk |
| `vp_da_chakra_tinh_x4.png` | polished glowing blue crystal |
| `vp_soi_chakra_x4.png` | a coil of glowing blue thread |
| `vp_da_thu_x4.png` | folded piece of animal hide |
| `vp_go_cung_x4.png` | short log of dense dark wood |
| `vp_day_thep_x4.png` | coil of thin steel wire |
| `vp_bot_no_x4.png` | small pouch of black powder, spilling |
| `vp_ngoc_lua_x4.png` | red ember gem |
| `vp_ngoc_gio_x4.png` | pale green swirl gem |
| `vp_ngoc_nuoc_x4.png` | blue droplet gem |
| `vp_ngoc_set_x4.png` | yellow jagged gem |

### 2D. Cuộn kỹ năng Làng Lá — 8

Cỡ 64×64. `SUBJECT: a technique scroll, rolled parchment tied with a coloured cord, a small emblem seal on the front —` rồi thêm:

| Tên file | Emblem |
|---|---|
| `vp_cuon_phanthan_x4.png` | two overlapping human silhouettes, white cord |
| `vp_cuon_thaythe_x4.png` | a log crossed with a leaf, brown cord |
| `vp_cuon_anthan_x4.png` | a fading swirl, grey cord |
| `vp_cuon_hoadon_x4.png` | a flame, red cord |
| `vp_cuon_kunaibay_x4.png` | crossed kunai, steel-grey cord |
| `vp_cuon_tuchakra_x4.png` | a footprint on a tree trunk, green cord |
| `vp_cuon_hoiphuc_x4.png` | a healing palm, pale blue cord |
| `vp_cuon_lienhoan_x4.png` | a spiral of three fists, orange cord |

### 2E. Vật phẩm nhiệm vụ — 10

| Tên file | SUBJECT |
|---|---|
| `vp_meo_tora_x4.png` | a grumpy calico cat with a red ribbon on one ear |
| `vp_cuon_troi_x4.png` | scroll with the character-less "sky" seal, white and gold |
| `vp_cuon_dat_x4.png` | scroll with an "earth" seal, brown and dark green |
| `vp_giayphep_thi_x4.png` | an official application form with a red stamp |
| `vp_the_dangky_x4.png` | a wooden registration tag on a cord |
| `vp_bando_xusong_x4.png` | a rolled coastal map, partly unrolled |
| `vp_thu_hokage_x4.png` | a sealed letter with a wax stamp |
| `vp_manh_matna_x4.png` | a broken piece of a white hunter mask |
| `vp_kiem_gay_x4.png` | the broken tip of an oversized cleaver sword |
| `vp_hop_com_x4.png` | a bento box, closed, with a cloth wrap |

### 2F. Hộp / rương / khác — 15

| Tên file | SUBJECT |
|---|---|
| `vp_ruong_go_x4.png` | small wooden chest, closed |
| `vp_ruong_sat_x4.png` | iron-banded chest, closed |
| `vp_tui_ninja_x4.png` | leather ninja pouch |
| `vp_tui_shuriken_x4.png` | thigh holster pouch |
| `vp_hop_qua_x4.png` | gift box with orange ribbon |
| `vp_chia_khoa_x4.png` | old iron key |
| `vp_ve_dichuyen_x4.png` | teleport ticket, paper talisman with a red seal |
| `vp_bua_hosinh_x4.png` | revival talisman, paper strip with black markings |
| `vp_da_maylan_x4.png` | four-leaf clover encased in resin |
| `vp_da_baove_x4.png` | protective stone, grey with a rune |
| `vp_da_nangcap_x4.png` | upgrade stone, glowing orange core in rock |
| `vp_hop_gacha_x4.png` | a sealed summoning box with a spiral mark |
| `vp_the_doiten_x4.png` | a blank name tag on a cord |
| `vp_binh_tayky_x4.png` | a wash-bottle for erasing seals, white |
| `vp_da_thanhtay_x4.png` | a clear crystal shard, faintly white |

---

## NHÓM 3 — TRANG PHỤC NHÂN VẬT (36 ảnh) → `03-trang-phuc/`

**Khối phụ:** `character sprite part for a 2D side-view game, facing right, neutral idle pose, no background, no shadow, consistent lighting from upper-left, the part must align to a shared body skeleton.`

12 bộ, mỗi bộ 3 file: `_dau_` (đầu, 128×128), `_than_` (thân + tay, 128×160), `_chan_` (chân, 128×128).

| # | Mã bộ | Mô tả bộ (dùng cho cả 3 phần) |
|---|---|---|
| 1 | `hocvien` | plain academy outfit: simple shirt, no headband, short dark hair |
| 2 | `genin_co` | basic genin: orange-and-blue jacket, leaf headband |
| 3 | `genin_cai` | improved genin: same but with a shoulder guard and utility belt |
| 4 | `chunin_luoi` | chunin: dark mesh underlayer, sleeveless top, headband |
| 5 | `chunin_giap` | chunin: green flak vest, long sleeves |
| 6 | `jonin_giap` | jonin: green flak vest with scroll pouches, arm wraps, mask optional |
| 7 | `anbu` | ANBU: grey armour plates, white animal mask, arm guards, short cloak |
| 8 | `kage` | Kage: white and orange robe with wide sleeves, conical hat held not worn |
| 9 | `tuluyen` | training gear: sleeveless, bandaged arms, weighted ankle bands |
| 10 | `lehoi` | festival yukata, orange with leaf pattern |
| 11 | `aomua` | straw rain cloak and wide hat |
| 12 | `anthan` | stealth outfit: all dark grey, hood up, face wrap |

Tên file ví dụ: `tp_genin_co_dau_x4.png`, `tp_genin_co_than_x4.png`, `tp_genin_co_chan_x4.png`

---

## NHÓM 4 — VŨ KHÍ (20 ảnh) → `04-vu-khi/`

**Khối phụ:** `single weapon, side view, blade or tip pointing to the upper right at 45 degrees, no hand.` Cỡ **96×48** trừ khi ghi khác.

| Tên file | SUBJECT |
|---|---|
| `vk_kunai_thuong_x4.png` | plain steel kunai with ring pommel and cloth grip wrap |
| `vk_kunai_thep_x4.png` | refined kunai, brighter steel, blue wrap |
| `vk_kunai_no_x4.png` | kunai with an explosive paper tag tied to the ring |
| `vk_kunai_ba_x4.png` | three-pronged teleport kunai, black seal marks on the blade |
| `vk_shuriken_x4.png` (64×64) | four-point throwing star |
| `vk_shuriken_lon_x4.png` (96×96) | giant folding four-blade shuriken |
| `vk_senbon_x4.png` | a set of three thin needles, bound |
| `vk_tanto_x4.png` | short straight sword with a square guard |
| `vk_kiem_ngan_x4.png` | short curved sword, dark scabbard beside it |
| `vk_dao_doi_x4.png` | a pair of matched short blades crossed |
| `vk_cuon_no_x4.png` (64×64) | explosive tag scroll, red markings |
| `vk_bom_khoi_x4.png` (64×64) | small black smoke bomb sphere with a fuse |
| `vk_bom_dinh_x4.png` (64×64) | caltrop cluster, four spikes |
| `vk_day_thep_cuon_x4.png` (64×64) | coiled steel wire spool |
| `vk_quat_gio_x4.png` | folding war fan, three purple moons on cream |
| `vk_bua_lon_x4.png` | oversized cleaver sword with an oval hole near the tip |
| `vk_giao_x4.png` | short spear with a leaf-shaped head |
| `vk_chuy_x4.png` | stone-headed war club |
| `vk_cung_x4.png` | short recurve bow, unstrung look |
| `vk_con_roi_tay_x4.png` (96×96) | a small wooden puppet arm segment with joints and a hidden blade |

---

## NHÓM 5 — QUÁI (12 ảnh) → `05-quai/`

**Khối phụ:** `creature sprite for a 2D side-view game, facing right, side profile, idle stance, full body visible, no background, no shadow.` Cỡ **96×96**.

| Tên file | SUBJECT |
|---|---|
| `quai_bia_go_x4.png` | a wooden training dummy post with a rope-wrapped body and painted target rings |
| `quai_soi_rung_x4.png` | a lean grey forest wolf, snarling |
| `quai_lon_rung_x4.png` | a bristly wild boar with chipped tusks |
| `quai_chim_ung_x4.png` | a brown hunting hawk, wings half spread |
| `quai_nhen_doc_x4.png` | a purple spider with yellow eye cluster |
| `quai_ca_sau_x4.png` | a river crocodile, low and long |
| `quai_ran_lon_x4.png` | a thick green snake coiled, head raised |
| `quai_khi_rung_x4.png` | a red-faced forest monkey baring teeth |
| `quai_bo_cap_x4.png` | a dark scorpion, tail raised |
| `quai_ninja_gia_x4.png` | a masked bandit in dark ragged clothes holding a short blade |
| `quai_linh_danhthue_x4.png` | a mercenary in mismatched armour with a spear |
| `quai_nguoi_rom_x4.png` | a straw training effigy animated by faint blue chakra threads |

---

## NHÓM 6 — BOSS (6 ảnh) → `06-boss/`

**Khối phụ:** `boss character sprite for a 2D side-view game, facing right, imposing stance, full body, no background, no shadow, stronger rim light.` Cỡ **160×160**.

| Tên file | SUBJECT |
|---|---|
| `boss_mizuki_x4.png` | a treacherous academy instructor: silver hair, dark uniform, oversized shuriken on his back, sneering |
| `boss_zabuza_x4.png` | a tall bandaged-face swordsman, bare torso with striped arm warmers, enormous cleaver sword on his back, mist at his feet |
| `boss_haku_x4.png` | a slender masked figure in a pale green kimono, white hunter mask with red markings, senbon needles between fingers |
| `boss_ran_khonglo_x4.png` | a colossal purple serpent rearing up, yellow slit eyes, fangs bared |
| `boss_trom_bitmat_x4.png` | a hulking masked brigand leader with a heavy chain wrapped around one arm |
| `boss_orochimaru_x4.png` | a pale long-black-haired figure with golden snake eyes and purple rope belt, half-hidden by shadow |

---

## NHÓM 7 — NPC (15 ảnh + 15 chân dung) → `07-npc/`

**Khối phụ (sprite):** `NPC sprite for a 2D side-view game, facing right, standing idle, full body, no background.` Cỡ **96×128**.
**Khối phụ (chân dung):** `character bust portrait, shoulders up, facing slightly left, neutral expression, flat cel shading, no background.` Cỡ **200×200**.

Mỗi nhân vật tạo **2 ảnh**: `npc_<mã>_x4.png` và `cd_<mã>_x4.png`.

| Mã | SUBJECT |
|---|---|
| `hokage3` | an elderly village leader in white-and-red robes with a conical hat, smoking a long pipe, kind but weary |
| `iruka` | a young teacher with a scar across the nose, green vest, warm smile |
| `kakashi` | a silver-haired jonin, face mask, one eye covered by a slanted headband, reading a small book |
| `chu_ramen` | a cheerful middle-aged noodle-stand cook with a white cloth headband and apron |
| `ba_vukhi` | an old weapons-shop woman with grey bun hair, arms folded |
| `thay_hocvien` | a strict academy instructor with square glasses and a clipboard |
| `linh_gac` | a bored gate guard in a green vest leaning on a post |
| `nguoi_giao_nv` | a mission-desk clerk surrounded by scroll stacks |
| `y_ta` | a hospital nurse in a pale uniform holding a clipboard |
| `ong_tazuna` | a grizzled old bridge-builder with a straw hat and a sake bottle |
| `co_be_hoa` | a small girl holding a basket of flowers |
| `tho_ren` | a burly blacksmith with a leather apron and hammer |
| `nguoi_dua_thu` | a courier with a shoulder satchel of scrolls, mid-stride |
| `trongtai_thi` | a tournament proctor in dark uniform, one hand raised |
| `nha_chep_cuon` | a thin scribe with ink-stained fingers and a brush behind the ear |

---

## NHÓM 8 — HIỆU ỨNG (10 tấm sprite) → `08-hieu-ung/`

**Khối phụ:** `SPRITE SHEET, exactly 10 frames arranged in a 5 columns × 2 rows grid, each cell 240×240, transparent background, consistent centre point and scale across all frames, no character, no ground, no cell borders, no numbers.` Tổng ảnh **1200×480**.

| Tên file | SUBJECT (chuỗi hoạt ảnh) |
|---|---|
| `hu_phanthan_x4.png` | a puff of white smoke expanding then splitting into three smaller puffs, then fading |
| `hu_thaythe_x4.png` | a swirl of smoke revealing a wooden log outline, then the log dropping and smoke clearing |
| `hu_anthan_x4.png` | a body-shaped shimmer distorting, fading to transparent ripples, then gone |
| `hu_hoadon_x4.png` | a small ember growing into a large orange fireball travelling right, then bursting |
| `hu_kunaibay_x4.png` | a motion-blur streak of a thrown blade with a spark on impact |
| `hu_tuchakra_x4.png` | blue chakra gathering at a point, swirling upward, pulsing twice, settling |
| `hu_hoiphuc_x4.png` | soft green light rising in gentle spirals, sparkles, fading |
| `hu_choang_x4.png` | three small yellow stars orbiting in a ring, looping smoothly |
| `hu_chimang_x4.png` | a sharp white-orange impact flash with radiating spikes, then fading |
| `hu_bangkinh_x4.png` | thin ice mirrors forming in an arc, glinting, then shattering into shards |

---

## NHÓM 9 — MẢNH GHÉP BẢN ĐỒ (80 ảnh) → `09-ban-do/`

**Khối phụ:** `seamless tileable terrain tile for a 2D side-scrolling game, flat cel shading, no perspective distortion, edges must tile seamlessly with neighbours, no characters, no objects unless stated, no shadow outside the tile.` Cỡ **96×96**.

### 9A. Làng Lá — mặt đất và đường (12)
`map_la_dat_1..3_x4.png` đất nện nâu · `map_la_co_1..3_x4.png` cỏ xanh ngắn · `map_la_datlat_1..3_x4.png` đường lát đá phẳng · `map_la_soi_1..3_x4.png` sỏi rời

### 9B. Làng Lá — kiến trúc (14)
`map_la_maingoi_1..3_x4.png` mái ngói đỏ cam · `map_la_tuonggo_1..3_x4.png` tường gỗ đứng · `map_la_tuongvua_1..2_x4.png` tường vữa trắng · `map_la_cua_x4.png` cửa gỗ trượt · `map_la_cuaso_x4.png` cửa sổ giấy · `map_la_hangrao_1..2_x4.png` hàng rào tre · `map_la_denlong_x4.png` đèn lồng cam treo · `map_la_conglang_x4.png` (192×192) cổng làng lớn hai cột đỏ

### 9C. Rừng (16)
`map_rung_thancay_1..4_x4.png` thân cây lớn có vỏ sần · `map_rung_tanla_1..4_x4.png` tán lá dày · `map_rung_bui_1..3_x4.png` bụi rậm · `map_rung_da_1..3_x4.png` tảng đá rêu · `map_rung_re_1..2_x4.png` rễ cây nổi

### 9D. Nước và cầu — Xứ Sóng (14)
`map_song_nuoc_1..3_x4.png` mặt nước gợn · `map_song_bo_1..2_x4.png` bờ đất ướt · `map_cau_van_1..3_x4.png` ván cầu gỗ · `map_cau_cot_1..2_x4.png` cột cầu · `map_bien_song_1..2_x4.png` sóng biển · `map_suongmu_1..2_x4.png` lớp sương mờ (bán trong suốt)

### 9E. Học viện và nội thất (10)
`map_hv_sango_1..2_x4.png` sàn gỗ · `map_hv_bang_x4.png` bảng đen · `map_hv_ban_x4.png` bàn học · `map_hv_ke_x4.png` kệ cuộn giấy · `map_hv_chieu_1..2_x4.png` chiếu tatami · `map_hv_tuong_1..2_x4.png` tường trong nhà · `map_hv_bia_x4.png` bia gỗ tập ném

### 9F. Núi và nền xa (14)
`map_nui_da_1..4_x4.png` vách đá · `map_nui_dinh_1..2_x4.png` đỉnh núi · `map_nen_troi_1..3_x4.png` nền trời (chuyển sắc nhẹ) · `map_nen_may_1..3_x4.png` mây xa · `map_nen_nuihokage_x4.png` (384×192) núi khắc mặt các đời thủ lĩnh, silhouette

---

# PHẦN C — KIỂM TRA TRƯỚC KHI GỬI

- [ ] Đã làm 10 ảnh giao diện đầu tiên và **được duyệt** trước khi làm tiếp
- [ ] Mọi ảnh **nền trong suốt**, không nền trắng
- [ ] Không có chữ, số, logo, chữ ký trong ảnh
- [ ] Kích thước đúng cột `Cỡ`
- [ ] Tên file đúng quy tắc, đúng thư mục
- [ ] Các tấm sprite hiệu ứng: đúng lưới 5×2, tâm không xê dịch giữa các khung
- [ ] Mảnh ghép bản đồ: ghép cạnh nhau không thấy đường nối
- [ ] Trang phục: cùng tư thế, cùng hướng ánh sáng, khớp khung xương chung

**Tổng đợt 1: ~360 ảnh.**
