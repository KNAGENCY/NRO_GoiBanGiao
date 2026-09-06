# -*- coding: utf-8 -*-
"""Phan bien: do dau bang cac bo phan CUNG (da mat, bang do, tam kim loai, mat) thay cho bbox 0-42%."""
import json, sys, math
from pathlib import Path
import numpy as np
from PIL import Image, ImageDraw
from scipy import ndimage as ndi

ART = Path(r"D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1")
OUT = Path(r"D:\NRO_GoiBanGiao\docs\my-thuat\review\khe-dau-than-2026-09-05") / "thu-nghiem-out"  # 05/09 trỏ thư mục tạm của phiên (đã mất)
OUT.mkdir(parents=True, exist_ok=True)
STRIPS = ["player_idle_3f", "player_run_4f", "player_melee_4f", "player_cast_3f", "player_hit_death_getup_3f", "master_iruka_x4", "master_hoc_vien_x4"]

def rgb_to_hsv(a):
    a = a.astype(np.float32) / 255.0
    r, g, b = a[..., 0], a[..., 1], a[..., 2]
    mx = a.max(-1); mn = a.min(-1); d = mx - mn
    h = np.zeros_like(mx)
    m = d > 1e-6
    dd = np.where(d == 0, 1, d)
    rm = m & (mx == r); gm = m & (mx == g) & ~rm; bm = m & ~rm & ~gm
    h[rm] = (60 * ((g - b) / dd) % 360)[rm]
    h[gm] = (60 * ((b - r) / dd) + 120)[gm]
    h[bm] = (60 * ((r - g) / dd) + 240)[bm]
    s = np.where(mx > 0, d / np.where(mx == 0, 1, mx), 0)
    return h, s, mx

def largest_cc(mask, region=None):
    lab, n = ndi.label(mask)
    if n == 0: return None
    sizes = ndi.sum(mask, lab, range(1, n + 1))
    order = np.argsort(-sizes)
    for k in order:
        comp = lab == (k + 1)
        if region is None or region(comp): return comp
    return None

def bbox(m):
    ys, xs = np.where(m)
    if len(ys) == 0: return None
    return xs.min(), ys.min(), xs.max() + 1, ys.max() + 1

def measure_frame(rgb, tag):
    H, W, _ = rgb.shape
    bg = np.array([246, 4, 242])
    fg = np.abs(rgb.astype(int) - bg).sum(-1) > 90
    fg = ndi.binary_opening(fg, iterations=1)
    cb = bbox(fg)
    h, s, v = rgb_to_hsv(rgb)
    skin = fg & (h > 20) & (h < 48) & (s > 0.35) & (s < 0.8) & (v > 0.55)
    skin = ndi.binary_opening(skin, iterations=1)
    light = fg & (v > 0.55) & (s < 0.22)               # tam kim loai + long trang mat
    navy = fg & (h > 195) & (h < 250) & (s > 0.28) & (v < 0.45)  # bang do / quan / tay ao
    # mat = khoi da lon nhat co tam nam trong 55% tren cua hop noi dung
    top_lim = cb[1] + 0.55 * (cb[3] - cb[1])
    def in_top(c):
        ys, xs = np.where(c); return ys.mean() < top_lim
    face = largest_cc(skin, in_top)
    r = {"tag": tag, "content_h": cb[3] - cb[1], "content_w": cb[2] - cb[0], "content_bottom": cb[3], "content_top": cb[1]}
    if face is None:
        return r
    fb = bbox(face)
    fx0, fy0, fx1, fy1 = fb
    r.update(face_w=fx1 - fx0, face_h=fy1 - fy0, face_area=int(face.sum()), chin_y=fy1, face_top=fy0, face_cx=(fx0 + fx1) / 2)
    # lo trong khoi da (mat, long may, mieng)
    filled = ndi.binary_fill_holes(face)
    holes = filled & ~face
    eye = largest_cc(holes)
    if eye is not None:
        eb = bbox(eye); r.update(eye_h=eb[3] - eb[1], eye_w=eb[2] - eb[0], eye_cy=(eb[1] + eb[3]) / 2, eye_area=int(eye.sum()))
        r["eye_to_chin"] = fy1 - r["eye_cy"]
    # dinh toc: hang fg cao nhat trong cot cua mat (+-50% be rong mat)
    xa = max(0, int(fx0 - 0.5 * (fx1 - fx0))); xb = min(W, int(fx1 + 0.5 * (fx1 - fx0)))
    col = fg[:, xa:xb].any(1)
    hair_top = int(np.argmax(col))
    r["hair_top"] = hair_top; r["head_h_hair_chin"] = fy1 - hair_top
    # bang do: navy trong khoang hang [hair_top, face_top+30% face_h], cot cua mat; be day = so hang navy theo tung cot (median 1/3 giua)
    band_region = np.zeros_like(navy); band_region[hair_top:int(fy0 + 0.3 * (fy1 - fy0)), fx0:fx1] = True
    band = navy & band_region
    bandc = largest_cc(band)
    if bandc is not None:
        cols = range(int(fx0 + (fx1 - fx0) / 3), int(fx0 + 2 * (fx1 - fx0) / 3))
        th = [bandc[:, c].sum() for c in cols]
        th = [t for t in th if t > 0]
        r["band_thick"] = float(np.median(th)) if th else None
        bb_ = bbox(bandc); r["band_w"] = bb_[2] - bb_[0]
    # tam kim loai: vung sang trong khoang hang [hair_top, face_top + 15% face_h]
    plate_region = np.zeros_like(light); plate_region[hair_top:int(fy0 + 0.15 * (fy1 - fy0)), max(0, fx0 - 40):min(W, fx1 + 40)] = True
    plate = largest_cc(light & plate_region)
    if plate is not None and plate.sum() > 150:
        pb = bbox(plate); r.update(plate_w=pb[2] - pb[0], plate_h=pb[3] - pb[1], plate_area=int(plate.sum()), plate_cy=(pb[1] + pb[3]) / 2)
    # nam tay / da khac ngoai mat
    others = skin & ~face
    fist = largest_cc(others)
    if fist is not None:
        r["fist_sqrt_area"] = math.sqrt(fist.sum())
    r["feet_to_chin"] = cb[3] - fy1
    # debug
    dbg = Image.fromarray(rgb).convert("RGB"); d = ImageDraw.Draw(dbg)
    d.rectangle([fx0, fy0, fx1, fy1], outline=(255, 0, 0), width=3)
    d.line([xa, hair_top, xb, hair_top], fill=(0, 255, 255), width=3)
    if "eye_cy" in r: d.ellipse([r["face_cx"] - 6, r["eye_cy"] - 6, r["face_cx"] + 6, r["eye_cy"] + 6], fill=(0, 0, 255))
    if plate is not None and "plate_w" in r:
        pb = bbox(plate); d.rectangle([int(x) for x in pb], outline=(0, 255, 0), width=3)
    if bandc is not None:
        bb_ = bbox(bandc); d.rectangle([int(x) for x in bb_], outline=(255, 255, 0), width=2)
    d.rectangle([int(x) for x in cb], outline=(255, 255, 255), width=1)
    ov = np.array(dbg)
    ov[face] = (ov[face] * 0.4 + np.array([255, 0, 0]) * 0.6).astype(np.uint8)
    Image.fromarray(ov).save(OUT / f"dbg_{tag}.png")
    return r

rows = []
for n in STRIPS:
    rep = json.loads((ART / "qa" / "crop-reports" / f"{n}_crop.json").read_text(encoding="utf-8"))
    edges = rep["crop"]["source_edges"]
    im = np.array(Image.open(ART / f"{n}_chroma.png").convert("RGB"))
    for i in range(len(edges) - 1):
        fr = im[:, edges[i]:edges[i + 1]]
        r = measure_frame(fr, f"{n}_{i}")
        r["strip"] = n; r["i"] = i
        rows.append(r)

keys = ["content_h", "face_h", "face_w", "face_area", "eye_h", "eye_to_chin", "head_h_hair_chin", "band_thick", "plate_w", "plate_h", "plate_area", "fist_sqrt_area", "feet_to_chin"]
print("%-28s" % "frame" + "".join("%9s" % k[:8] for k in keys))
for r in rows:
    print("%-28s" % r["tag"] + "".join("%9s" % (("%.0f" % r[k]) if r.get(k) is not None else "-") for k in keys))
json.dump(rows, open(OUT / "do.json", "w"), indent=1, default=float)

# ti le so voi idle (trung binh 3 khung) — dung sqrt cho dien tich
idle = [r for r in rows if r["strip"] == "player_idle_3f"]
def mean(rs, k, sq=False):
    v = [r[k] for r in rs if r.get(k) is not None]
    v = [math.sqrt(x) for x in v] if sq else v
    return sum(v) / len(v) if v else None
print("\n== ti le trung binh moi dai / idle (1.00 = bang idle) ==")
mk = [("content_h", False), ("face_h", False), ("face_w", False), ("face_area", True), ("eye_h", False), ("eye_to_chin", False), ("head_h_hair_chin", False), ("band_thick", False), ("plate_w", False), ("plate_area", True), ("fist_sqrt_area", False)]
print("%-28s" % "dai" + "".join("%9s" % k[:8] for k, _ in mk))
for n in STRIPS:
    rs = [r for r in rows if r["strip"] == n]
    line = "%-28s" % n
    for k, sq in mk:
        a = mean(rs, k, sq); b = mean(idle, k, sq)
        line += "%9s" % (("%.3f" % (a / b)) if a and b else "-")
    print(line)
print("\n== ti le tung khung / idle ==")
print("%-28s" % "khung" + "".join("%9s" % k[:8] for k, _ in mk))
for r in rows:
    line = "%-28s" % r["tag"]
    for k, sq in mk:
        b = mean(idle, k, sq); a = r.get(k)
        if a is not None and sq: a = math.sqrt(a)
        line += "%9s" % (("%.3f" % (a / b)) if a and b else "-")
    print(line)
