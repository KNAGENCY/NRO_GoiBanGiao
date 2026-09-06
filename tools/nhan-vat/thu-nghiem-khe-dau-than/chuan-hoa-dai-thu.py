# -*- coding: utf-8 -*-
"""chuan-hoa-thu.py -- thu nghiem HUONG ANH.

Buoc 1: chuan hoa ti le TUNG KHUNG theo co dau idle_0 (thuoc = cac bo phan cung: cao/rong da mat,
        cao mat, mat->cam, toc->cam, day bang do; do tu do.json cua phan bien 3).
Buoc 2: day khung da chuan hoa qua DUNG phep toan pipeline (task6b_closure.player_component_images,
        monkeypatch bo doc anh + tuy chon luat cat) -> 34 manh x1 + part.bin va dy/dx.
Buoc 3: do bang tools/nhan-vat/ghep-khung.py (cong cu cua du an) va MILP tinh dy tot nhat.
Bien the:
   goc        : k=1, luat cat 28/76 (phai tai hien dung bang do trong de)
   scale-pct  : k tung khung, luat cat 28/76 (chi chuan hoa ti le)
   scale-anat : k tung khung, luat cat theo moc giai phau (dau: toc->cam+co; than: cam+6 -> 76%; chan: 60% -> 100%)
Khong ghi gi vao repo/client; moi thu nam trong wf-huong-anh/.
"""
import copy, json, math, struct, subprocess, sys
from pathlib import Path

import numpy as np
from PIL import Image, ImageDraw
from scipy.optimize import Bounds, LinearConstraint, milp

sys.stdout.reconfigure(encoding="utf-8")
GOC = Path(r"D:\NRO_GoiBanGiao")
sys.path.insert(0, str(GOC))
from tools.artqa import task6b_closure as t6  # noqa: E402

OUT = Path(r"D:\NRO_GoiBanGiao\docs\my-thuat\review\khe-dau-than-2026-09-05") / "thu-nghiem-out"  # 05/09 trỏ thư mục tạm của phiên (đã mất)
DO = Path(r"D:\NRO_GoiBanGiao\docs\my-thuat\review\khe-dau-than-2026-09-05") / "moc-giai-phau-do.json")
ART = GOC / "docs/my-thuat/anh-nguon/vertical-slice-v1"
REVIEW = GOC / "docs/my-thuat/review/vertical-slice-v1.json"
PART_BIN = GOC / "server/internal/game/assets/part.bin"
GHEP = GOC / "tools/nhan-vat/ghep-khung.py"
CHARINFO = json.load(open(GOC / "tools/nhan-vat/charinfo.json", encoding="utf-8"))
SO_KHUNG = {0: 3, 1: 17, 2: 14}
PART_IDS = {"head": 30, "body": 57, "leg": 58}
NAME_OF = {a: n for a, n, _ in t6.PLAYER_ACTIONS}          # action -> file
STRIP_OF = {n: n.replace("_chroma.png", "") for _, n, _ in t6.PLAYER_ACTIONS}  # file -> do.json strip key

# ---------- do dau: he so k tung khung ----------
do = {(r["strip"], r["i"]): r for r in json.load(open(DO, encoding="utf-8"))}
REF = do[("player_idle_3f", 0)]
FEATS = ["face_h", "face_w", "eye_h", "eye_to_chin", "head_h_hair_chin", "band_thick"]
EXCL = {
    ("player_melee_4f", 1): {"face_h", "eye_to_chin", "head_h_hair_chin"},   # tay che cam
    ("player_hit_death_getup_3f", 0): {"eye_h", "eye_to_chin", "band_thick", "head_h_hair_chin"},  # nham mat, ngua dau
}
THAY_THE = {("player_hit_death_getup_3f", 1): ("player_hit_death_getup_3f", 0)}  # khung NAM -> chep hit0


def he_so(strip, i):
    r = do[(strip, i)]
    ex = EXCL.get((strip, i), set())
    rs = [REF[f] / r[f] for f in FEATS if f not in ex and r.get(f)]
    return float(np.median(rs)) if rs else 1.0


def resize_premult(im, k):
    if abs(k - 1.0) < 1e-9:
        return im.convert("RGBA")
    a = np.asarray(im.convert("RGBA")).astype(np.float32) / 255.0
    pm = a.copy()
    pm[..., :3] *= a[..., 3:4]
    w, h = im.size
    nw, nh = max(1, round(w * k)), max(1, round(h * k))
    pm_im = Image.fromarray((pm * 255).round().astype(np.uint8), "RGBA").resize((nw, nh), Image.LANCZOS)
    b = np.asarray(pm_im).astype(np.float32) / 255.0
    al = b[..., 3:4]
    rgb = np.where(al > 0, b[..., :3] / np.maximum(al, 1e-6), 0)
    out = np.concatenate([np.clip(rgb, 0, 1), al], -1)
    return Image.fromarray((out * 255).round().astype(np.uint8), "RGBA")


# ---------- dung dai da chuan hoa + review gia ----------
review = json.loads(REVIEW.read_text(encoding="utf-8"))
records = t6._review_by_name(review)
_cache = {}
ALPHA = {n: t6._load_reviewed_alpha(ART, records, n, _cache) for _, n, _ in t6.PLAYER_ACTIONS}
_orig_component_rect = t6._component_rect
_orig_loader = t6._load_reviewed_alpha


def dung_bien_the(scale_theo_khung: bool, thay_khung_nam: bool):
    """Tra (review2, NORM: file->RGBA strip, LM: (file,i)->landmarks, KS: (file,i)->k)."""
    NORM, LM, KS = {}, {}, {}
    review2 = copy.deepcopy(review)
    by_name = {a["asset"]: a for a in review2["assets"]}
    for action, fname, n in t6.PLAYER_ACTIONS:
        strip = STRIP_OF[fname]
        rec = records[fname]
        frames = []
        for i in range(n):
            src_strip, src_i = (strip, i)
            if thay_khung_nam and (strip, i) in THAY_THE:
                src_strip, src_i = THAY_THE[(strip, i)]
            src_fname = src_strip + "_chroma.png"
            x, y, w, h = t6._frame_rect(records[src_fname], src_i)
            crop = ALPHA[src_fname].crop((x, y, x + w, y + h))
            k = he_so(src_strip, src_i) if scale_theo_khung else 1.0
            im = resize_premult(crop, k)
            r = do[(src_strip, src_i)]
            chin_rel = (r["chin_y"] - y) * k if r.get("chin_y") is not None else None
            frames.append((im, k, chin_rel, (src_strip, src_i)))
        CW = max(f[0].width for f in frames) + 20
        CH = max(f[0].height for f in frames) + 20
        canvas = Image.new("RGBA", (CW * n, CH), (0, 0, 0, 0))
        boxes, edges = [], [c * CW for c in range(n + 1)]
        for i, (im, k, chin_rel, src) in enumerate(frames):
            canvas.alpha_composite(im, (i * CW + 10, 10))
            boxes.append([10, 10, 10 + im.width, 10 + im.height])
            LM[(fname, i)] = {"chin": chin_rel, "H": im.height, "k": k, "src": src}
            KS[(fname, i)] = k
        NORM[fname] = canvas
        rec2 = by_name[fname]
        rec2["checks"]["crop"]["source_edges"] = edges
        rec2["checks"]["crop"]["content_boxes"] = boxes
        rec2["checks"]["crop"]["source_size"] = [canvas.width, canvas.height]
    return review2, NORM, LM, KS


NECK_SRC, MARGIN_SRC, LEG_TOP_PCT = 55, 6, 0.60


def make_anat_cut(LM):
    def _component_rect_anat(record, frame_index, component):
        x, y, w, h = t6._frame_rect(record, frame_index)
        lm = LM[(record["asset"], frame_index)]
        chin = lm["chin"] if lm["chin"] is not None else h * 0.42
        if component == "head":
            top, bottom = 0, min(h, round(chin + NECK_SRC))
        elif component == "body":
            top, bottom = round(chin + MARGIN_SRC), max(round(chin + MARGIN_SRC) + 1, round(h * 0.76))
        elif component == "leg":
            top, bottom = round(h * LEG_TOP_PCT), h
        else:
            raise ValueError(component)
        return [x, y + top, w, bottom - top]
    return _component_rect_anat


def sinh_components(review2, NORM, LM, cut):
    t6._load_reviewed_alpha = lambda art_root, recs, name, cache: NORM[name]
    t6._component_rect = make_anat_cut(LM) if cut == "anat" else _orig_component_rect
    try:
        recs2 = t6._review_by_name(review2)
        comps = t6player = t6.player_component_images(review2, art_root=ART, records=recs2, charinfo=CHARINFO)
    finally:
        t6._load_reviewed_alpha = _orig_loader
        t6._component_rect = _orig_component_rect
    assert len(comps) == 34
    return comps


# ---------- ghi sprite x1 + part.bin ----------
def ghi_sprites(comps, thu_muc):
    thu_muc.mkdir(parents=True, exist_ok=True)
    for c in comps:
        t6.zoom_variants(c["image_x4"])["x1"].save(thu_muc / f"Small{c['id']}.png")


def va_part_bin(dxy, duong_ra):
    """dxy: {(component, slot): (dx, dy)}"""
    d = bytearray(PART_BIN.read_bytes())
    n = struct.unpack_from(">h", d, 0)[0]
    off, offs = 2, []
    for _ in range(n):
        t = d[off]
        offs.append((t, off + 1))
        off += 1 + SO_KHUNG[t] * 4
    for comp, pid in PART_IDS.items():
        t, o = offs[pid]
        for k in range(SO_KHUNG[t]):
            if (comp, k) in dxy:
                dx, dy = dxy[(comp, k)]
                struct.pack_into(">bb", d, o + k * 4 + 2, dx, dy)
    duong_ra.write_bytes(d)


def chay_ghep(part_path, spr_dir, out_png):
    r = subprocess.run([sys.executable, str(GHEP), "--part-bin", str(part_path), "--sprites", str(spr_dir), "--out", str(out_png)],
                       capture_output=True, text=True, encoding="utf-8", errors="replace")
    if r.returncode != 0:
        print(r.stdout, r.stderr)
        raise SystemExit("ghep-khung loi")
    rows = {}
    for line in r.stdout.splitlines():
        p = [s.strip() for s in line.split("|")]
        if len(p) >= 7 and p[0].isdigit():
            rows[int(p[0])] = dict(khe_dt=int(p[2]), khe_tc=int(p[3]), chan=int(p[4]), cao=int(p[6]))
    return rows


# ---------- hinh hoc thuan so + MILP ----------
def geo(comps):
    return {(c["component"], c["slot_index"]): dict(H=t6.zoom_variants(c["image_x4"])["x1"].height,
                                                       dy=c["placement"]["dy"], dx=c["placement"]["dx"], id=c["id"])
            for c in comps}


def khe(cf, G, dyh, dyb, dyl):
    (ih, _, cih), (il, _, cil), (ib, _, cib) = CHARINFO[cf][0], CHARINFO[cf][1], CHARINFO[cf][2]
    head_top = -cih + dyh[ih]; head_bot = head_top + G[("head", ih)]["H"]
    body_top = -cib + dyb[ib]; body_bot = body_top + G[("body", ib)]["H"]
    leg_top = -cil + dyl[il]; leg_bot = leg_top + G[("leg", il)]["H"]
    return body_top - head_bot, leg_top - body_bot, leg_bot, leg_bot - head_top


def milp_fit(G, t_hb, t_bl, lo_hb, hi_hb, lo_bl, hi_bl):
    """dy dau + dy than nguyen (chan co dinh); muc tieu min tong |khe - muc tieu|; rang buoc cung [lo, hi] (None = bo)."""
    nH, nB, nL, nC = 3, 17, 14, 33
    n = nH + nB + nH + nB + 2 * nC  # dyh, dyb, |d dyh|, |d dyb| (khong dung), u, v
    iH = lambda s: s; iB = lambda s: nH + s
    iU = lambda k: nH + nB + nH + nB + k; iV = lambda k: iU(nC) + k
    dyh0 = {s: G[("head", s)]["dy"] for s in range(nH)}
    dyb0 = {s: G[("body", s)]["dy"] for s in range(nB)}
    dyl0 = {s: G[("leg", s)]["dy"] for s in range(nL)}
    rows, lo, hi = [], [], []
    cobj = np.zeros(n)
    for cf in range(nC):
        (ih, _, cih), (il, _, cil), (ib, _, cib) = CHARINFO[cf][0], CHARINFO[cf][1], CHARINFO[cf][2]
        c_hb = cih - cib - G[("head", ih)]["H"]              # khe_hb = dyb - dyh + c_hb
        c_bl = -cil + dyl0[il] + cib - G[("body", ib)]["H"]  # khe_bl = -dyb + c_bl
        if lo_hb is not None:
            r = np.zeros(n); r[iB(ib)] = 1; r[iH(ih)] = -1; rows.append(r); lo.append(lo_hb - c_hb); hi.append(hi_hb - c_hb)
        if lo_bl is not None:
            r = np.zeros(n); r[iB(ib)] = -1; rows.append(r); lo.append(lo_bl - c_bl); hi.append(hi_bl - c_bl)
        for sg in (1, -1):
            r = np.zeros(n); r[iB(ib)] = sg; r[iH(ih)] = -sg; r[iU(cf)] = -1; rows.append(r); lo.append(-np.inf); hi.append(-sg * (c_hb - t_hb))
            r = np.zeros(n); r[iB(ib)] = -sg; r[iV(cf)] = -1; rows.append(r); lo.append(-np.inf); hi.append(-sg * (c_bl - t_bl))
        cobj[iU(cf)] = 1; cobj[iV(cf)] = 1
    lb = np.full(n, -128.0); ub = np.full(n, 127.0)
    for i in range(nH + nB, n): lb[i] = 0
    integ = np.zeros(n); integ[: nH + nB] = 1
    res = milp(cobj, constraints=LinearConstraint(np.array(rows), lo, hi), bounds=Bounds(lb, ub), integrality=integ)
    if res.status != 0:
        return None
    x = np.round(res.x).astype(int)
    dyh = {s: int(x[iH(s)]) for s in range(nH)}; dyb = {s: int(x[iB(s)]) for s in range(nB)}
    return dyh, dyb, dyl0


def bang(G, dyh, dyb, dyl):
    return [khe(cf, G, dyh, dyb, dyl) for cf in range(33)]


def in_bang(ten, ke, t_hb=None, t_bl=None):
    hb = [k[0] for k in ke]; bl = [k[1] for k in ke]
    print(f"  {ten}")
    print("    khe dau-than :", " ".join(f"{v:+d}" for v in hb))
    print("    khe than-chan:", " ".join(f"{v:+d}" for v in bl))
    print(f"    dau-than: min {min(hb)} max {max(hb)} | so cf HO(>0): {sum(1 for v in hb if v > 0)} | than-chan: min {min(bl)} max {max(bl)} | cao tong {min(k[3] for k in ke)}..{max(k[3] for k in ke)} | ban chan {sorted(set(k[2] for k in ke))}")
    if t_hb is not None:
        e1 = [abs(v - t_hb) for v in hb]; e2 = [abs(v - t_bl) for v in bl]
        print(f"    lech so muc tieu ({t_hb}/{t_bl}): dau-than TB {np.mean(e1):.1f} max {max(e1)} (cf {int(np.argmax(e1))}); than-chan TB {np.mean(e2):.1f} max {max(e2)} (cf {int(np.argmax(e2))})")


# ---------- chay ----------
def main():
    OUT.mkdir(parents=True, exist_ok=True)
    print("=== He so k tung khung (idle_0 = 1.000) ===")
    for action, fname, n in t6.PLAYER_ACTIONS:
        strip = STRIP_OF[fname]
        print(f"  {strip:28s}", " ".join(f"{he_so(strip, i):.3f}" for i in range(n)))
    print()
    S = 128 / 340  # ti le chung (idle_0 khong doi k=1 nen giu nguyen)
    print("=== Tu the: cam->ban chan sau chuan hoa (x1 px) so voi bo xuong client (~35.5 = 61 tong - 25.5 dau) ===")
    for action, fname, n in t6.PLAYER_ACTIONS:
        strip = STRIP_OF[fname]
        for i in range(n):
            r = do[(strip, i)]
            if r.get("chin_y") is None:
                print(f"  {strip}_{i}: NAM, khong do"); continue
            k = he_so(strip, i)
            ftc = r["feet_to_chin"] * k * S / 4; H = r["content_h"] * k
            print(f"  {strip}_{i}: k={k:.3f} cao khung {H:5.0f} px (idle 655) | cam->chan {ftc:5.1f} px x1 | thieu so bo xuong {35.5 - ftc:+5.1f} px")
    print()

    ket_qua = {}
    for ten, scale, thay, cut in [("goc", False, False, "pct"), ("scale-pct", True, True, "pct"), ("scale-anat", True, True, "anat")]:
        print(f"################ BIEN THE {ten} ################")
        review2, NORM, LM, KS = dung_bien_the(scale, thay)
        for fname, im in NORM.items():
            im.save(OUT / f"strip-{ten}-{STRIP_OF[fname]}.png")
        comps = sinh_components(review2, NORM, LM, cut)
        G = geo(comps)
        spr = OUT / f"sprites-{ten}"
        ghi_sprites(comps, spr)
        dxy = {(c["component"], c["slot_index"]): (c["placement"]["dx"], c["placement"]["dy"]) for c in comps}
        pb = OUT / f"part-{ten}.bin"; va_part_bin(dxy, pb)
        rows = chay_ghep(pb, spr, OUT / f"ghep-{ten}.png")
        dyh = {s: G[("head", s)]["dy"] for s in range(3)}; dyb = {s: G[("body", s)]["dy"] for s in range(17)}; dyl = {s: G[("leg", s)]["dy"] for s in range(14)}
        ke = bang(G, dyh, dyb, dyl)
        # doi chieu cong thuc thuan so voi ghep-khung
        lech = [(cf, rows[cf]["khe_dt"], ke[cf][0], rows[cf]["khe_tc"], ke[cf][1]) for cf in range(33) if rows[cf]["khe_dt"] != ke[cf][0] or rows[cf]["khe_tc"] != ke[cf][1]]
        print(f"  ghep-khung vs cong thuc: {33 - len(lech)}/33 khop" + (f" LECH: {lech}" if lech else ""))
        print("  kich co manh x1 (H): dau", [G[("head", s)]["H"] for s in range(3)], "than", [G[("body", s)]["H"] for s in range(17)], "chan", [G[("leg", s)]["H"] for s in range(14)])
        in_bang("dy pipeline (neo cf nho nhat) -- do bang ghep-khung", [(rows[cf]["khe_dt"], rows[cf]["khe_tc"], rows[cf]["chan"], rows[cf]["cao"]) for cf in range(33)])
        if cut == "pct":
            t_hb, t_bl, win = -8, -6, (-9, 0, -11, 0)
        else:
            t_hb, t_bl, win = -4, -9, (-6, -2, -12, -6)
        fit = milp_fit(G, t_hb, t_bl, None, None, None, None)
        if fit:
            ke2 = bang(G, *fit)
            in_bang(f"MILP dy dau+than, chan co dinh, KHONG rang buoc cung, muc tieu {t_hb}/{t_bl}", ke2, t_hb, t_bl)
            dyh2, dyb2, _ = fit
            dxy2 = {("head", s): (dxy[("head", s)][0], dyh2[s]) for s in range(3)}
            dxy2.update({("body", s): (dxy[("body", s)][0], dyb2[s]) for s in range(17)})
            dxy2.update({("leg", s): dxy[("leg", s)] for s in range(14)})
            pb2 = OUT / f"part-{ten}-milp.bin"; va_part_bin(dxy2, pb2)
            rows2 = chay_ghep(pb2, spr, OUT / f"ghep-{ten}-milp.png")
            print("    (ghep-khung xac nhan:", "khop" if all(rows2[cf]["khe_dt"] == ke2[cf][0] and rows2[cf]["khe_tc"] == ke2[cf][1] for cf in range(33)) else "LECH", ")")
            print("    dy dau moi", dyh2, "| dy than moi", dyb2)
        fit_c = milp_fit(G, t_hb, t_bl, *win)
        print(f"  MILP voi cua so cung dau-than [{win[0]},{win[1]}] than-chan [{win[2]},{win[3]}]: " + ("CO NGHIEM" if fit_c else "VO NGHIEM"))
        if fit_c:
            in_bang("   nghiem cua so cung", bang(G, *fit_c), t_hb, t_bl)
        ket_qua[ten] = dict(k={f"{STRIP_OF[k[0]]}_{k[1]}": v for k, v in KS.items()}, geo={f"{c}:{s}": v for (c, s), v in G.items()},
                            ghep=rows, milp=(fit and {"dyh": fit[0], "dyb": fit[1]}))
        print()
    json.dump(ket_qua, open(OUT / "ket-qua.json", "w", encoding="utf-8"), indent=1, default=str)
    print("da ghi", OUT)


if __name__ == "__main__":
    main()
