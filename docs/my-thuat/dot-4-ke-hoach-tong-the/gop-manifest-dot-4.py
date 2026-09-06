# -*- coding: utf-8 -*-
"""Gộp 9 manifest hạng mục đợt 4 -> manifest-dot-4.json (một file máy đọc cho cả đợt).

Chạy lại MỖI KHI một manifest hạng mục đổi:
    python docs\\my-thuat\\dot-4-ke-hoach-tong-the\\gop-manifest-dot-4.py

Nguồn sự thật vẫn là từng NN-*.json cạnh brief; file gộp chỉ nối anh[] lại và chuẩn hoá
mã làng. (Bản gốc do agent soát viết trong thư mục tạm của phiên 06/09; chép vào repo để
không mất — bài học "script nằm ngoài repo là script sẽ mất".)
"""
import io, json, os, collections
from collections import Counter, OrderedDict

D  = os.path.dirname(os.path.abspath(__file__))
RA = os.path.join(D, "manifest-dot-4.json")
MA = [("01", "01-nhan-vat",              "Nhân vật người chơi"),
      ("02", "02-npc",                   "NPC"),
      ("03", "03-quai",                  "Quái"),
      ("04", "04-boss",                  "Boss"),
      ("05", "05-map",                   "Map (Làng Lá còn thiếu · Khu rừng · 4 làng)"),
      ("06", "06-vat-pham",              "Icon vật phẩm · trang bị · vũ khí"),
      ("07", "07-ky-nang-hieu-ung",      "Kỹ năng · hiệu ứng"),
      ("08", "08-giao-dien",             "Giao diện (UI / HUD)"),
      ("09", "09-chan-dung-thuong-hieu", "Chân dung · thương hiệu")]

def jl(f):
    return json.loads(io.open(os.path.join(D, f + ".json"), encoding="utf-8").read(),
                      object_pairs_hook=OrderedDict)

J = {ma: jl(ma) for _, ma, _ in MA}

# ---------------------------------------------------------------- chuẩn hoá mã làng
CHUAN = {"la": "la", "lang_la": "la", "làng lá": "la", "Làng Lá": "la",
         "cat": "cat", "lang_cat": "cat", "Làng Cát": "cat",
         "smu": "smu", "lang_smu": "smu", "lang_suong_mu": "smu", "Làng Sương Mù": "smu",
         "da": "da", "lang_da": "da", "Làng Đá": "da",
         "may": "may", "lang_may": "may", "Làng Mây": "may",
         "chung": "chung", "chung_ui": "chung", "rung": "rung"}

MAU_LANG = OrderedDict([
    ("la",   "Làng Lá / hệ Hoả — cam–đỏ ấm trên navy và xanh ô-liu (Leaf / Fire: warm orange-red on dark navy and olive green)"),
    ("cat",  "Làng Cát / hệ Phong — vàng cát, nâu sa mạc, lục nhạt (Sand / Wind: sand-yellow, sun-bleached tan, earth brown, pale-green accents)"),
    ("smu",  "Làng Sương Mù / hệ Thuỷ — xanh dương sâu, lam, xám đá, bạc lạnh (Mist / Water: deep blue, lake-blue, slate grey, cool damp look)"),
    ("da",   "Làng Đá / hệ Thổ — nâu đất, xám đá, be (Stone / Earth: earth brown, rock grey, beige)"),
    ("may",  "Làng Mây / hệ Lôi — CHÀM / xanh-tím (R−G < 100), trắng-xanh, vàng chớp; CẤM tím hồng-magenta vì trùng nền chroma "
             "(Cloud / Lightning: indigo blue-violet, blue-white, lightning-yellow — never hot pink or magenta)"),
    ("rung", "Khu rừng (map 6, không phải làng) — lục ngọc sẫm và rêu, nâu vỏ cây ẩm, sương xanh-lam ở xa"),
    ("chung","Dùng chung 5 làng / dân thường — nâu đất, ô-liu bạc màu, vàng rơm, trắng ngà; không huy hiệu làng"),
])

# ---------------------------------------------------------------- phong cách khoá
PCK = ("detailed 2D game sprite, hand-painted with cel-shading, clean dark outlines, "
       "rich shading (40+ colors), Naruto ninja aesthetic, side/3-4 view facing RIGHT, full body, "
       "NO flat vector shapes, NO background scenery, NO drop shadow on ground.")
bien_the = OrderedDict()
for _, ma, _ in MA:
    for k, v in J[ma].items():
        if isinstance(v, str) and (k.startswith("phong_cach_khoa") or
                                   k in ("luat_icon", "khoi_dai_chroma", "duoi_prompt_chung")):
            bien_the["%s.%s" % (ma, k)] = v

# ---------------------------------------------------------------- ảnh tham chiếu
atc, atc_mo_ta = OrderedDict(), OrderedDict()
for _, ma, _ in MA:
    for k, v in J[ma].get("anh_tham_chieu", {}).items():
        f  = v.get("file") if isinstance(v, dict) else v
        mt = v.get("mo_ta") if isinstance(v, dict) else None
        if k in atc and atc[k] != f:
            k = "%s::%s" % (ma[:2], k)          # tránh đè nếu một ngày có va chạm
        atc[k] = f
        if mt: atc_mo_ta.setdefault(k, mt)
thieu_atc = sorted(k for k, v in atc.items() if not (isinstance(v, str) and os.path.exists(v)))

# ---------------------------------------------------------------- gộp anh[]
anh, stt = [], 0
dem_hm, dem_uu = OrderedDict(), Counter()
for so, ma, ten in MA:
    d = J[ma]
    muc = list(d.get("anh", []))
    # 05 có thêm concept[] (tham chiếu, KHÔNG lắp) — README §5 ghi rõ
    for c in d.get("concept", []):
        c = OrderedDict(c); c["nhom"] = "concept"; c["kieu"] = "tham_chieu_khong_lap"
        muc.append(c)
    c_uu = Counter()
    for a in muc:
        stt += 1
        e = OrderedDict()
        e["stt_gop"]  = stt
        e["hang_muc"] = ma
        e["ten_file"] = a.get("ten_file")
        goc = a.get("lang")
        chuan = CHUAN.get(goc) if isinstance(goc, str) else None
        e["lang"] = chuan
        if isinstance(goc, str) and chuan != goc:
            e["lang_goc"] = goc
        for k, v in a.items():
            if k in ("ten_file", "lang", "stt"):
                continue
            e[k] = v
        if "stt" in a:
            e["stt_trong_hang_muc"] = a["stt"]
        anh.append(e)
        kieu_khong_lap = a.get("kieu") == "tham_chieu_khong_lap"
        ts = a.get("trang_thai")
        key = "concept" if kieu_khong_lap else ("da_co" if ts == "da_co" else a.get("uu_tien"))
        c_uu[key] += 1
        dem_uu[key] += 1
    dem_hm[ma] = OrderedDict([("ten", ten), ("tong_anh", len(muc)),
                              ("P1", c_uu["P1"]), ("P2", c_uu["P2"]), ("P3", c_uu["P3"]),
                              ("da_co", c_uu["da_co"]), ("concept_khong_lap", c_uu["concept"]),
                              ("brief_md", os.path.join(D, ma + ".md")),
                              ("manifest_goc", os.path.join(D, ma + ".json"))])

trung = OrderedDict()
for t, c in Counter(a["ten_file"] for a in anh).items():
    if c > 1:
        trung[t] = sorted(set(a["hang_muc"] for a in anh if a["ten_file"] == t))

out = OrderedDict()
out["meta"] = OrderedDict([
    ("du_an", "Game Naruto (lõi client Ngọc Rồng) — mỹ thuật đợt 4"),
    ("dot", 4),
    ("ngay", "2026-09-06"),
    ("mo_ta",
     "GỘP 9 manifest hạng mục của đợt 4 thành MỘT file máy đọc. Nguồn sự thật vẫn là từng NN-*.json cạnh brief; "
     "file này chỉ nối anh[] lại, thêm khoá hang_muc + stt_gop và CHUẨN HOÁ mã làng về la|cat|smu|da|may|rung|chung "
     "(04-boss và 07 dùng lang_la/lang_cat/lang_suong_mu…, 01 dùng 'Làng Lá' — mã gốc giữ ở lang_goc). "
     "Sinh lại bằng scratchpad\\gop-manifest.py khi một manifest hạng mục đổi."),
    ("thu_muc_ke_hoach", D),
    ("luu_anh_tao_ra_tai", r"D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\dot-4\<hạng mục>\  (mỗi anh[].luu_vao ghi đường dẫn tuyệt đối riêng)"),
    ("tong_anh", len(anh)),
    ("tong_anh_theo_uu_tien", OrderedDict([
        ("P1", dem_uu["P1"]), ("P2", dem_uu["P2"]), ("P3", dem_uu["P3"]),
        ("da_co", dem_uu["da_co"]), ("concept_khong_lap", dem_uu["concept"]),
        ("tong", len(anh)),
        ("ghi_chu", "P1/P2/P3 = số file CẦN VẼ (trang_thai can_ve hoặc ve_lai). da_co = đã có, KHÔNG vẽ. "
                    "concept_khong_lap = 5 bức concept map của 05, vẽ để đo và duyệt, không lắp vào game. "
                    "Cộng thô: một khuôn mặt chân dung xuất hiện ở cả brief chủ thể (01/02/04) và 09 với tên khác nhau "
                    "— xem README.md §6 mục D1 trước khi lấy số này làm số ảnh phải trả tiền."),
    ])),
    ("theo_hang_muc", dem_hm),
    ("cach_dung",
     "Duyệt anh[] → BỎ trang_thai == 'da_co' → làm 3 ảnh mẫu của hạng mục đó trước (INDEX-DOT-4.json → anh_mau_3_truoc), "
     "dừng chờ duyệt → rồi P1 → P2 → P3. Mỗi ảnh: gửi prompt_day_du + các ảnh trong dung_tham_chieu (tra anh_tham_chieu) "
     "tới AI vẽ → lưu vào anh[].luu_vao + ten_file → đo lại cỡ thật so kich_thuoc "
     "(dải chroma: rộng chia hết so_khung, nền tím RGB(246,4,242) ở 4 góc)."),
    ("luu_y", [
        "Mã làng đã chuẩn hoá; mã gốc của manifest hạng mục nằm ở lang_goc (chỉ có khi khác).",
        "kich_thuoc của dòng trang_thai='da_co' KHÔNG phải lúc nào cũng là cỡ file thật: "
        "03-quai ghi cỡ nếu vẽ lại theo chuẩn ô 724 (4 dải wolf_* thật là 1774×887 / 2172×724), "
        "06-vat-pham ghi cỡ chuẩn trong kich-thuoc-anh.json (nguồn vẽ to gấp 4). Muốn cỡ thật thì đo file.",
        "6 mục da_co của 07 (5 skill_icons_* + hu_anthan_10f) có prompt_day_du RỖNG và không có khối prompt ở mục 6 của brief.",
        "223 tên tĩnh mới (05:165 · 08:30 · 09:17 · 04:10 · 06:1) phải NỐI CUỐI "
        "D:\\NRO_GoiBanGiao\\docs\\my-thuat\\kich-thuoc-anh.json trước khi chạy KIEM-TRA-ANH.py / TAO-4-MUC-ZOOM.py; "
        "hiện chưa nối tên nào (đã kiểm 06/09). 04 và 09 đặt HAI mã tra cứu khác nhau cho cùng 3 chân dung boss Làng Lá — "
        "chốt D1 (README §6) rồi mới nối, đừng nối cả hai.",
        "Mọi thứ máy chủ chưa chốt đều ghi 'đề xuất' trong brief: quái 4 làng + Khu rừng (quai.yaml chỉ có sói), "
        "boss 4 làng (boss.yaml chỉ có 2 boss Lá), NPC ngoài Iruka (KH9 chưa có spec), danh mục vật phẩm 4 làng (KH7/KH8), "
        "đường lắp hiệu ứng (chưa có bộ đóng gói EffectData).",
    ]),
    ("nguon_gop", [os.path.join(D, ma + ".json") for _, ma, _ in MA]),
    ("sinh_boi", "gop-manifest.py (critic nhất quán đợt 4)"),
])
out["phong_cach_khoa"] = PCK
out["phong_cach_khoa_ghi_chu"] = (
    "Chuỗi trên là bản khoá cho SPRITE (nhân vật, NPC, quái, boss, chân dung, hiệu ứng). "
    "Map dùng bản môi trường ('detailed 2D game environment art for a Naruto ninja world…'), "
    "UI và icon vật phẩm có bản riêng — xem phong_cach_khoa_bien_the. "
    "Mọi anh[].prompt_day_du đã tự chứa bản đúng của nó, không cần ghép thêm.")
out["phong_cach_khoa_bien_the"] = bien_the
out["mau_lang"] = MAU_LANG
out["mau_lang_theo_hang_muc"] = OrderedDict(
    (ma, OrderedDict((CHUAN.get(k, k), v) for k, v in J[ma].get("mau_lang", {}).items()))
    for _, ma, _ in MA)
out["anh_tham_chieu"] = atc
out["anh_tham_chieu_mo_ta"] = atc_mo_ta
out["anh_tham_chieu_thieu_tren_dia"] = thieu_atc
out["trung_ten_giua_hang_muc"] = trung
out["anh"] = anh

io.open(RA, "w", encoding="utf-8", newline="\n").write(json.dumps(out, ensure_ascii=False, indent=1) + "\n")

# ------------------------------------------------------------------ kiểm lại
kt = json.load(io.open(RA, encoding="utf-8"))
print("Đã ghi %s (%d byte)" % (RA, os.path.getsize(RA)))
print("json.load lại: OK — anh[] = %d" % len(kt["anh"]))
print("theo uu_tien:", json.dumps(kt["meta"]["tong_anh_theo_uu_tien"], ensure_ascii=False)[:220])
print("trùng tên:", kt["trung_ten_giua_hang_muc"])
print("tham chiếu thiếu:", kt["anh_tham_chieu_thieu_tren_dia"])
print("mã làng sau chuẩn hoá:", Counter(a.get("lang") for a in kt["anh"]))
thieu = [a["ten_file"] for a in kt["anh"] if not a.get("luu_vao")]
print("thiếu luu_vao:", thieu)
xau = [a["ten_file"] for a in kt["anh"]
       if a.get("luu_vao") and not os.path.isdir(a["luu_vao"].rstrip("\\/"))]
print("luu_vao trỏ thư mục chưa có:", xau)
