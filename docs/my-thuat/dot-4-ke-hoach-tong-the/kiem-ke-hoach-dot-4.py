# -*- coding: utf-8 -*-
"""
kiem-ke-hoach-dot-4.py — đếm lại kế hoạch đợt 4 bằng máy (chỉ đọc 9 manifest JSON).

Chạy:  python D:\\NRO_GoiBanGiao\\docs\\my-thuat\\dot-4-ke-hoach-tong-the\\kiem-ke-hoach-dot-4.py
       (thêm  --khong-ghi  nếu chỉ muốn xem, không ghi INDEX-DOT-4.json)

Làm gì:
  1. Kiểm 9 cặp brief  NN-*.md + NN-*.json  có đủ không.
  2. Đếm anh[] theo uu_tien × trang_thai của từng manifest (số thật, không lấy từ meta).
  3. Tìm TÊN FILE TRÙNG giữa các manifest (một ảnh bị hai brief cùng đòi vẽ).
  4. Kiểm mọi đường dẫn trong anh_tham_chieu có tồn tại trên đĩa không.
  5. Kiểm thư mục LƯU VÀO (meta.luu_vao và anh[].luu_vao) đã tạo chưa.
  6. Ghi INDEX-DOT-4.json (mục lục máy đọc cho tool tự động) cạnh file này.

Không sửa manifest, không sửa ảnh, không commit. Chỉ cần Python 3 chuẩn.
"""
import io
import json
import os
import sys
from collections import Counter, defaultdict

THU_MUC = os.path.dirname(os.path.abspath(__file__))
INDEX_RA = os.path.join(THU_MUC, "INDEX-DOT-4.json")

# Thứ tự, tên hiển thị, kiểu file chính, 3 ảnh mẫu (theo §7 của từng brief).
HANG_MUC = [
    ("01", "01-nhan-vat", "Nhân vật người chơi", "dai_chroma",
     ["player_run_5f_chroma.png", "player_hit_death_getup_3f_chroma.png", "player_cat_idle_3f_chroma.png"]),
    ("02", "02-npc", "NPC", "dai_chroma (master) + tinh_alpha (chân dung)",
     ["master_hokage3_x4_chroma.png", "master_tho_ren_x4_chroma.png", "cd_hokage3_x4.png"]),
    ("03", "03-quai", "Quái", "dai_chroma",
     ["master_lon_rung_x4_chroma.png", "lon_rung_run_4f_chroma.png", "nhen_doc_hit_death_2f_chroma.png"]),
    ("04", "04-boss", "Boss", "dai_chroma",
     ["boss_la_mat_na_vo_idle_3f_chroma.png", "boss_la_ke_la_mat_idle_3f_chroma.png", "boss_la_mat_na_vo_melee_4f_chroma.png"]),
    ("05", "05-map", "Map (Làng Lá còn thiếu · Khu rừng · 4 làng)", "tinh_alpha (nền đặc + tile/vật thể alpha)",
     ["nen_nuihokage_rong2_x4.png", "tile_dat_van_1_x4.png", "map_rung_cay_co_thu_1_x4.png"]),
    ("06", "06-vat-pham", "Icon vật phẩm · trang bị · vũ khí", "tinh_alpha",
     ["vp_kunai_luyentap_x4.png", "vp_cat_ao_t1_x4.png", "vp_may_quan_t1_x4.png"]),
    ("07", "07-ky-nang-hieu-ung", "Kỹ năng · hiệu ứng", "dai_chroma",
     ["hu_la_haohoacau_no_6f_chroma.png", "hu_la_ylieu_8f_chroma.png", "hu_la_haohoacau_bay_3f_chroma.png"]),
    ("08", "08-giao-dien", "Giao diện (UI / HUD)", "tinh_alpha (đúng cỡ, thay file client)",
     ["ui_hud_khung_trai_x4.png", "ui_khung_popup_9o_x4.png", "ui_nut_thuong_3manh_x4.png"]),
    ("09", "09-chan-dung-thuong-hieu", "Chân dung · thương hiệu", "dai_chroma (chân dung) + tinh_alpha đục (thương hiệu)",
     ["nguoi_choi_la_portrait_chroma.png", "ui_icon_app_x4.png", "nhan_gia_mat_na_p1_portrait_chroma.png"]),
]

# "Đã có" theo cách từng brief đếm (gồm cả ảnh đã có nằm NGOÀI anh[] như 51 ảnh map 1,
# 100 vp + 20 vk của lô 343, 62 ui_*). Số trong anh[] tự đếm ở dưới.
DA_CO_THEO_BRIEF = {"01": 2, "02": 2, "03": 5, "04": 0, "05": 51, "06": 120, "07": 6, "08": 63, "09": 1}


def doc_json(duong_dan):
    with io.open(duong_dan, "r", encoding="utf-8") as f:
        return json.load(f)


def lay_luu_vao(meta, anh):
    ds = set()
    for k in ("luu_vao", "luu_vao_master", "luu_vao_chan_dung"):
        v = meta.get(k)
        if isinstance(v, str) and v.strip():
            ds.add(v.rstrip("\\/"))
    for a in anh:
        v = a.get("luu_vao")
        if isinstance(v, str) and v.strip():
            ds.add(v.rstrip("\\/"))
    return sorted(ds)


def main():
    ghi = "--khong-ghi" not in sys.argv
    loi = 0
    index = {"meta": {
        "mo_ta": "Mục lục máy đọc của kế hoạch mỹ thuật đợt 4 — sinh bởi kiem-ke-hoach-dot-4.py từ 9 manifest JSON cùng thư mục. "
                 "Nguồn sự thật là từng NN-*.json; file này chỉ gom lại để tool tự động biết có bao nhiêu hạng mục, "
                 "manifest nào, lưu vào đâu, làm 3 ảnh mẫu nào trước.",
        "thu_muc": THU_MUC,
        "cach_dung": "Với mỗi hang_muc[]: mở manifest_json → duyệt anh[] → bỏ trang_thai == 'da_co' → làm anh_mau[] trước, "
                     "dừng chờ duyệt → rồi P1 → P2 → P3. Đọc README.md §5 để biết khoá nào tên khác nhau giữa các manifest.",
        "sinh_boi": os.path.basename(__file__),
    }, "hang_muc": [], "tong": {}, "trung_ten_giua_manifest": [], "anh_tham_chieu_thieu": [], "thu_muc_luu_vao_thieu": []}

    tong = Counter()
    ten_o_dau = defaultdict(list)          # ten_file -> [ma hạng mục]
    tham_chieu_thieu = []
    luu_vao_thieu = []

    print("=== KIỂM KẾ HOẠCH ĐỢT 4 — %s" % THU_MUC)
    for ma, goc, ten, kieu, anh_mau in HANG_MUC:
        md = os.path.join(THU_MUC, goc + ".md")
        js = os.path.join(THU_MUC, goc + ".json")
        co_md, co_js = os.path.isfile(md), os.path.isfile(js)
        if not (co_md and co_js):
            print("[X] %s: thiếu %s%s" % (goc, "" if co_md else "MD ", "" if co_js else "JSON"))
            loi += 1
            continue
        d = doc_json(js)
        meta = d.get("meta", {})
        anh = d.get("anh", [])
        dem = Counter()
        for a in anh:
            uu, tt = a.get("uu_tien", "?"), a.get("trang_thai", "?")
            dem[(uu, tt)] += 1
            ten_o_dau[a.get("ten_file", "?")].append(ma)
        can_ve = {p: sum(n for (uu, tt), n in dem.items() if uu == p and tt != "da_co") for p in ("P1", "P2", "P3")}
        da_co_trong = sum(n for (uu, tt), n in dem.items() if tt == "da_co")
        for p in ("P1", "P2", "P3"):
            tong[p] += can_ve[p]
        tong["da_co_theo_brief"] += DA_CO_THEO_BRIEF[ma]

        # ảnh tham chiếu
        atc = d.get("anh_tham_chieu", {})
        thieu_o_day = []
        if isinstance(atc, dict):
            for k, v in atc.items():
                p = v.get("file") if isinstance(v, dict) else v
                if isinstance(p, str) and p.startswith(("D:", "C:")) and not os.path.exists(p):
                    thieu_o_day.append((k, p))
        for k, p in thieu_o_day:
            tham_chieu_thieu.append({"hang_muc": goc, "khoa": k, "duong_dan": p})

        # thư mục lưu
        ds_luu = lay_luu_vao(meta, anh)
        for th in ds_luu:
            if "vertical-slice-v1" in th:
                continue  # chỗ để của ảnh da_co, không phải thư mục nộp
            if not os.path.isdir(th):
                luu_vao_thieu.append({"hang_muc": goc, "thu_muc": th})

        # tên mới cần nối vào kich-thuoc-anh.json (mỗi manifest để một chỗ khác nhau)
        # 05 và 08 để dạng dict {ten_file: [w, h]}; 04/06/09 để list trong meta — đếm cả hai kiểu
        ten_moi = d.get("ten_moi_can_them_kich_thuoc") or meta.get("ten_moi_can_them_kich_thuoc") \
            or d.get("ten_moi_kich_thuoc") or []
        so_ten_moi = len(ten_moi) if isinstance(ten_moi, (list, dict)) else 0

        brief_md = meta.get("brief_md") or meta.get("brief") or meta.get("file_brief") or md
        print("[OK] %s — %s" % (goc, ten))
        print("     cần vẽ: P1 %d · P2 %d · P3 %d | da_co trong anh[]: %d (brief tính đã có: %d) | tổng anh[]: %d"
              % (can_ve["P1"], can_ve["P2"], can_ve["P3"], da_co_trong, DA_CO_THEO_BRIEF[ma], len(anh)))
        print("     chi tiết: " + " · ".join("%s/%s=%d" % (uu, tt, n) for (uu, tt), n in sorted(dem.items())))
        print("     lưu vào: " + " | ".join(ds_luu))
        if thieu_o_day:
            loi += len(thieu_o_day)
            for k, p in thieu_o_day:
                print("     [X] ảnh tham chiếu THIẾU: %s -> %s" % (k, p))
        index["hang_muc"].append({
            "stt": ma, "ma": goc, "ten": ten, "kieu_file_chinh": kieu,
            "brief_md": brief_md, "manifest_json": js,
            "luu_vao": ds_luu,
            "so_anh_can_ve": can_ve, "da_co_trong_manifest": da_co_trong, "da_co_theo_brief": DA_CO_THEO_BRIEF[ma],
            "tong_anh_trong_manifest": len(anh),
            "anh_mau_3_truoc": anh_mau,
            "so_ten_moi_can_noi_kich_thuoc_anh": so_ten_moi,
        })

    # trùng tên giữa các manifest
    trung = {t: sorted(set(ms)) for t, ms in ten_o_dau.items() if len(set(ms)) > 1}
    print("\n=== TÊN FILE TRÙNG GIỮA CÁC MANIFEST: %d" % len(trung))
    for t in sorted(trung):
        print("     %s  <- %s" % (t, ", ".join(trung[t])))
        index["trung_ten_giua_manifest"].append({"ten_file": t, "hang_muc": trung[t]})

    print("\n=== ẢNH THAM CHIẾU THIẾU TRÊN ĐĨA: %d" % len(tham_chieu_thieu))
    index["anh_tham_chieu_thieu"] = tham_chieu_thieu
    print("=== THƯ MỤC LƯU VÀO CHƯA TẠO: %d" % len(luu_vao_thieu))
    for x in luu_vao_thieu:
        print("     [!] %s -> %s" % (x["hang_muc"], x["thu_muc"]))
    index["thu_muc_luu_vao_thieu"] = luu_vao_thieu

    index["tong"] = {"can_ve": {"P1": tong["P1"], "P2": tong["P2"], "P3": tong["P3"]},
                     "da_co_theo_brief": tong["da_co_theo_brief"],
                     "ghi_chu": "Tổng cộng thô theo từng brief; chân dung xuất hiện ở 01/02/04 VÀ 09 với tên khác nhau "
                                "(xem README.md §6 mục D1) nên số P1/P2 thật thấp hơn khi chốt vẽ mỗi chân dung một lần."}
    print("\n=== TỔNG (thô, cộng theo brief): P1 %d · P2 %d · P3 %d · đã có %d"
          % (tong["P1"], tong["P2"], tong["P3"], tong["da_co_theo_brief"]))

    if ghi:
        with io.open(INDEX_RA, "w", encoding="utf-8") as f:
            json.dump(index, f, ensure_ascii=False, indent=1)
        print("Đã ghi %s" % INDEX_RA)
    print("Lỗi cứng: %d" % loi)
    return 1 if loi else 0


if __name__ == "__main__":
    sys.exit(main())
