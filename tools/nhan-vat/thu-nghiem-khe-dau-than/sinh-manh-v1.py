# -*- coding: utf-8 -*-
"""sinh-manh.py — sinh 34 mảnh x1/x4 theo PHƯƠNG ÁN MỚI, dùng lại hàm của
tools/artqa/task6b_closure.py (alpha, frame rect, zoom). Không ghi vào repo.

Ba thay đổi so với pipeline hiện hành:
  1. Tỉ lệ TỪNG KHUNG s_f = s0 * (cỡ đầu idle / cỡ đầu khung f)   (chuẩn hoá theo đầu)
  2. Mép trên mảnh THÂN cắt theo CẰM của chính khung đó (chin_y - c px x1),
     không theo 28% chiều cao khung; ĐẦU cắt hair_top -> chin + 6 px nguồn.
  3. Gán lại slot: thân/chân slot 2..6 = run0,run1,run2,run3,run0 (cf 2..6 là 5 khung
     chạy theo Char.cs), slot thân 15 (cf 20/31) = cast1 thay khung NẰM, slot 0 tuỳ chọn.
Xuất: sprites/<ten>/Small<id>.png (x1) + x4/Small<id>.png + geometry.json.
"""
from __future__ import annotations
import json, sys, argparse
from pathlib import Path
from statistics import median
from PIL import Image

GOC = Path(r"D:\NRO_GoiBanGiao")
sys.path.insert(0, str(GOC))
from tools.artqa.task6b_closure import (  # noqa: E402
    _review_by_name, _frame_rect, _load_reviewed_alpha, zoom_variants,
    load_charinfo, charinfo_slot_to_cf, PLAYER_REFERENCE_X4_WIDTH,
)

ART = GOC / "docs/my-thuat/anh-nguon/vertical-slice-v1"
REVIEW = GOC / "docs/my-thuat/review/vertical-slice-v1.json"
DO_JSON = Path(r"D:\NRO_GoiBanGiao\docs\my-thuat\review\khe-dau-than-2026-09-05") / "moc-giai-phau-do.json"  # 05/09 trỏ thư mục tạm của phiên (đã mất)
STRIP = {"idle": "player_idle_3f_chroma.png", "run": "player_run_4f_chroma.png",
         "melee": "player_melee_4f_chroma.png", "cast": "player_cast_3f_chroma.png",
         "hit": "player_hit_death_getup_3f_chroma.png"}
TAG = {"idle": "player_idle_3f", "run": "player_run_4f", "melee": "player_melee_4f",
       "cast": "player_cast_3f", "hit": "player_hit_death_getup_3f"}
S0 = PLAYER_REFERENCE_X4_WIDTH / 340.0  # 128/340 như pipeline

# ---- gán slot (mặc định = phương án đề xuất) -------------------------------
HEAD_SLOTS = [("idle", 0), ("idle", 1), ("idle", 2)]
BODY_SLOTS = {
    0: ("idle", 0), 1: ("idle", 1),
    2: ("run", 0), 3: ("run", 1), 4: ("run", 2), 5: ("run", 3), 6: ("run", 0),
    7: ("melee", 0), 8: ("melee", 1), 9: ("melee", 2), 10: ("melee", 3),
    11: ("cast", 0), 12: ("cast", 1), 13: ("cast", 2),
    14: ("hit", 0), 15: ("cast", 1), 16: ("idle", 2),   # slot 16 (chỉ cf 25 "sau dịch chuyển"): đứng idle2 thay hit2 ngồi xổm
}
LEG_SLOTS = {
    0: ("idle", 0), 1: ("idle", 1),
    2: ("run", 0), 3: ("run", 1), 4: ("run", 2), 5: ("run", 3), 6: ("run", 0),
    7: ("melee", 0), 8: ("melee", 1), 9: ("melee", 2), 10: ("melee", 3),
    11: ("cast", 0), 12: ("cast", 1), 13: ("idle", 2),  # slot 13 chỉ cf 25 -> cùng khung idle2 với thân 16
}
# bản "cũ" để đối chứng (thứ tự khung như mapping hiện hành, khung nằm ở slot 15)
BODY_SLOTS_CU = dict(BODY_SLOTS); BODY_SLOTS_CU.update({2: ("idle", 2), 3: ("run", 0), 4: ("run", 1), 5: ("run", 2), 6: ("run", 3), 15: ("hit", 1), 16: ("hit", 2)})
LEG_SLOTS_CU = dict(LEG_SLOTS); LEG_SLOTS_CU.update({2: ("idle", 2), 3: ("run", 0), 4: ("run", 1), 5: ("run", 2), 6: ("run", 3), 13: ("cast", 2)})

IDLE_MEAN = {}


def do_metrics():
    d = json.loads(DO_JSON.read_text(encoding="utf-8"))
    m = {}
    for r in d:
        m[(r["strip"], r["i"])] = r
    return m


def head_ratio(m, strip, i):
    """cỡ đầu khung / cỡ đầu idle (median 6 số đo cứng, loại số đo nhiễu)."""
    keys = ["face_w", "face_h", "eye_h", "eye_to_chin", "head_h_hair_chin", "band_thick"]
    bad = {("player_melee_4f", 1): {"face_h", "eye_to_chin"},
           ("player_hit_death_getup_3f", 0): {"eye_h", "eye_to_chin", "band_thick", "head_h_hair_chin"}}
    r = m[(TAG[strip], i)]
    if "face_w" not in r:
        return None  # khung nằm
    ratios = []
    for k in keys:
        if k in bad.get((TAG[strip], i), set()):
            continue
        ratios.append(r[k] / IDLE_MEAN[k])
    return median(ratios)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--ten", default="A")
    ap.add_argument("--c", type=float, default=2.0, help="px x1 hàm dưới cằm gộp vào thân")
    ap.add_argument("--cu", action="store_true", help="giữ gán slot cũ (đối chứng)")
    ap.add_argument("--khong-chuan-hoa", action="store_true", help="tỉ lệ chung s0 (đối chứng)")
    ap.add_argument("--slot0", default=None, help="idle0|hit0|cast1 ... đổi thân slot 0")
    ap.add_argument("--gan", default="", help="đổi thân theo slot, vd '14=cast2,15=cast0' (chân đổi kèm nếu slot<14 và --gan-chan)")
    ap.add_argument("--gan-chan", default="", help="đổi chân theo slot, vd '9=cast1'")
    ap.add_argument("--fudge", default="", help="nhân thêm tỉ lệ từng khung, vd 'run1=1.05'")
    ap.add_argument("--dau", default="idle", help="idle | run0 : nguồn 3 đầu")
    ap.add_argument("--out", default=str(Path(__file__).parent / "sprites"))
    a = ap.parse_args()

    global IDLE_MEAN
    m = do_metrics()
    for k in ["face_w", "face_h", "eye_h", "eye_to_chin", "head_h_hair_chin", "band_thick"]:
        IDLE_MEAN[k] = sum(m[("player_idle_3f", i)][k] for i in range(3)) / 3.0

    review = json.loads(REVIEW.read_text(encoding="utf-8"))
    records = _review_by_name(review)
    cache = {}
    body_slots = dict(BODY_SLOTS_CU if a.cu else BODY_SLOTS)
    leg_slots = dict(LEG_SLOTS_CU if a.cu else LEG_SLOTS)
    if a.slot0:
        body_slots[0] = (a.slot0[:-1], int(a.slot0[-1]))
    for kv in filter(None, a.gan.split(",")):
        k, v = kv.split("="); body_slots[int(k)] = (v[:-1], int(v[-1]))
    for kv in filter(None, a.gan_chan.split(",")):
        k, v = kv.split("="); leg_slots[int(k)] = (v[:-1], int(v[-1]))
    fudge = {}
    for kv in filter(None, a.fudge.split(",")):
        k, v = kv.split("="); fudge[(k[:-1], int(k[-1]))] = float(v)
    head_slots = list(HEAD_SLOTS)
    if a.dau == "run0":
        head_slots[1] = ("run", 0)

    out = Path(a.out) / a.ten
    (out / "x4").mkdir(parents=True, exist_ok=True)

    # --- tỉ lệ và mốc giải phẫu từng khung -------------------------------
    frames = {}

    def frame(strip, i):
        key = (strip, i)
        if key in frames:
            return frames[key]
        name = STRIP[strip]
        alpha = _load_reviewed_alpha(ART, records, name, cache)
        fx, fy, fw, fh = _frame_rect(records[name], i)
        r = 1.0 if a.khong_chuan_hoa else head_ratio(m, strip, i)
        if r is None:
            r = 1.0
        s = S0 / r * fudge.get((strip, i), 1.0)
        resized = alpha.crop((fx, fy, fx + fw, fy + fh)).resize(
            (max(1, round(fw * s)), max(1, round(fh * s))), Image.Resampling.NEAREST)
        bb = resized.getchannel("A").getbbox()
        met = m[(TAG[strip], i)]
        chin = met.get("chin_y")
        if (TAG[strip], i) == ("player_melee_4f", 1):  # tay che cằm -> ước lượng từ face_top
            chin = met["face_top"] + IDLE_MEAN["face_h"] * r
        # do.json: y tuyệt đối trong dải, x tính theo Ô khung -> cộng mép ô (source_edges[i])
        edge = records[name]["checks"]["crop"]["source_edges"][i]
        info = dict(strip=strip, i=i, fx=fx, fy=fy, fw=fw, fh=fh, s=s, ratio=r, img=resized,
                    center_x4=(bb[0] + bb[2]) / 2.0, feet_x4=float(bb[3]),
                    chin_x4=(chin - fy) * s if chin is not None else None,
                    face_cx_x4=(met["face_cx"] + edge - fx) * s if "face_cx" in met else None,
                    hair_top_x4=(met["content_top"] - fy) * s)
        frames[key] = info
        return info

    def cut(strip, i, comp, sid, slot):
        f = frame(strip, i)
        s, H = f["s"], f["fh"]
        if comp == "head":
            top_src, bot_src = 0, (f["chin_x4"] / s) + 6  # hair_top -> cằm + 6 px nguồn
        elif comp == "body":
            if f["chin_x4"] is None:                        # khung nằm: không có cằm -> 28% như cũ
                top_src = H * 0.28
            else:
                top_src = (f["chin_x4"] / s) - a.c * 4 / s  # cằm - c px x1
            bot_src = H * 0.76
        else:
            top_src, bot_src = H * 0.66, H
        top = max(0, round(top_src * s)); bottom = min(f["img"].height, round(bot_src * s))
        piece = f["img"].crop((0, top, f["img"].width, bottom))
        bb = piece.getchannel("A").getbbox()
        piece = piece.crop(bb)
        left_x4, top_x4 = bb[0], top + bb[1]
        zooms = zoom_variants(piece)
        x1 = zooms["x1"]
        x1.save(out / f"Small{sid}.png"); zooms["x4"].save(out / "x4" / f"Small{sid}.png")
        bb1 = x1.getchannel("A").getbbox()
        # tâm hông (x1, so với mép trái ảnh): thân = 12 hàng x4 dưới cùng, chân = 12 hàng x4 trên cùng
        import numpy as np
        al = np.asarray(piece.getchannel("A")) > 0
        rows = al[-12:] if comp == "body" else al[:12]
        xs = np.where(rows.any(axis=0))[0]
        hip_cx = (xs.mean() / 4) if len(xs) else piece.width / 8
        return dict(
            id=sid, component=comp, slot=slot, strip=strip, frame=i, scale=s, ratio=f["ratio"],
            size_x4=list(piece.size), size_x1=list(x1.size), bbox_x1=list(bb1),
            # mốc giải phẫu, đơn vị x1, so với góc trên-trái ảnh x1
            chin_off=(f["chin_x4"] - top_x4) / 4 if f["chin_x4"] is not None else None,
            feet_off=(f["feet_x4"] - top_x4) / 4,
            hair_off=(f["hair_top_x4"] - top_x4) / 4,
            face_cx_off=(f["face_cx_x4"] - left_x4) / 4 if f["face_cx_x4"] is not None else None,
            center_off=(f["center_x4"] - left_x4) / 4,
            hip_cx_off=float(hip_cx),
            crop_top_x4=top, crop_bottom_x4=bottom,
        )

    comps = []
    for k, (st, i) in enumerate(head_slots):
        comps.append(cut(st, i, "head", 27041 + k, k))
    for k in range(17):
        st, i = body_slots[k]
        comps.append(cut(st, i, "body", 27044 + k, k))
    for k in range(14):
        st, i = leg_slots[k]
        comps.append(cut(st, i, "leg", 27061 + k, k))

    geo = dict(ten=a.ten, c=a.c, chuan_hoa=not a.khong_chuan_hoa, s0=S0, fudge={f"{k[0]}{k[1]}": v for k, v in fudge.items()},
               body_slots={k: f"{v[0]}{v[1]}" for k, v in body_slots.items()},
               leg_slots={k: f"{v[0]}{v[1]}" for k, v in leg_slots.items()},
               head_slots=[f"{v[0]}{v[1]}" for v in head_slots],
               frames={f"{k[0]}{k[1]}": dict(s=v["s"], ratio=v["ratio"], fh=v["fh"]) for k, v in frames.items()},
               components=comps)
    (out / "geometry.json").write_text(json.dumps(geo, indent=1, ensure_ascii=False), encoding="utf-8")
    print("khung  ratio_dau  s_f     H_khung_norm(x1)")
    for k, v in sorted(frames.items()):
        print("%-8s %.3f  %.4f  %5.1f" % (f"{k[0]}{k[1]}", v["ratio"], v["s"], v["fh"] * v["s"] / 4))
    print("id    manh slot khung   x1(WxH)  chin_off feet_off face_cx center")
    for c in comps:
        print("%5d %-4s %2d  %-7s %2dx%-2d  %s %7.1f %s %6.1f" % (
            c["id"], c["component"], c["slot"], f"{c['strip']}{c['frame']}", c["size_x1"][0], c["size_x1"][1],
            "%7.1f" % c["chin_off"] if c["chin_off"] is not None else "   -   ", c["feet_off"],
            "%6.1f" % c["face_cx_off"] if c["face_cx_off"] is not None else "   -  ", c["center_off"]))
    print("da ghi", out)


if __name__ == "__main__":
    main()
