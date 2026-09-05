# -*- coding: utf-8 -*-
"""toi-uu.py — chọn dx/dy nguyên cho 34 mảnh (part 30/57/58) bằng MILP.

Ẩn số: dy_h[3], dy_b[17], dx_h[3], dx_b[17], dx_l[14]. dy_l[14] CỐ ĐỊNH: bàn chân
(mép dưới chân) = cy-11 tại cf neo nhỏ nhất (đúng quy ước hiện hành).
Phép vẽ client: top = cy - CI_y(cf) + dy[slot]; left = cx + CI_x(cf) + dx[slot].

Sai số giải phẫu từng cf:
  e(cf)  = cằm đầu - đường cằm của khung thân  (muốn 0; e>0 đầu tụt vào cổ, e<0 hở)
  d(cf)  = bàn chân "ảo" của khung thân - (cy-11) (muốn 0; thân trượt so với chân)
  ex(cf) = tâm mặt đầu - tâm mặt khung thân   (ngang)
  gx(cf) = tâm khung thân - tâm khung chân     (ngang)
Ràng buộc cứng: khe đầu-thân (bbox) trong [KHE_MIN, 0]; khe thân-chân trong [-11, 0].
Mục tiêu: M*max|e| + sum|e| + w_d*sum|d| + sum|ex| + sum|gx|.
Ghi part.bin vá (bản sao server part.bin) + bảng dự đoán.
"""
from __future__ import annotations
import json, struct, sys, argparse, shutil
from pathlib import Path
import numpy as np
from scipy.optimize import milp, LinearConstraint, Bounds

GOC = Path(r"D:\NRO_GoiBanGiao")
sys.path.insert(0, str(GOC))
sys.path.insert(0, str(GOC / "tools" / "nhan-vat"))
from tools.artqa.task6b_closure import load_charinfo, charinfo_slot_to_cf  # noqa: E402
import importlib.util
spec = importlib.util.spec_from_file_location("ghep", GOC / "tools/nhan-vat/ghep-khung.py")
ghep = importlib.util.module_from_spec(spec); spec.loader.exec_module(ghep)

PART_GOC = GOC / "server/internal/game/assets/part.bin"
FEET = -11
import os
KL_MAX = int(os.environ.get('KL_MAX','0'))


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--geo", required=True)
    ap.add_argument("--out-part", required=True)
    ap.add_argument("--khe-min", type=int, default=-9)
    ap.add_argument("--e-min", type=float, default=-2.0)
    ap.add_argument("--e-max", type=float, default=1.0)
    ap.add_argument("--d-max", type=float, default=4.0)
    ap.add_argument("--w-d", type=float, default=0.5)
    ap.add_argument("--M", type=float, default=50.0)
    ap.add_argument("--w-cf", default="0:4,1:4,2:2,3:2,4:2,5:2,6:2,7:2,8:2",
                    help="trọng số theo cf (đứng/chạy/đánh hay hiện nhất)")
    a = ap.parse_args()
    w_cf = {int(k): float(v) for k, v in (kv.split(":") for kv in filter(None, a.w_cf.split(",")))}

    geo = json.loads(Path(a.geo).read_text(encoding="utf-8"))
    comps = {(c["component"], c["slot"]): c for c in geo["components"]}
    for c in comps.values():  # khung nằm không có cằm: coi mép trên mảnh + c là "đường cằm"
        if c["chin_off"] is None:
            c["chin_off"] = geo.get("c", 2.0)
        if c["face_cx_off"] is None:
            c["face_cx_off"] = c["center_off"]
    ci = load_charinfo()
    ncf = len(ci)
    # bbox thật của ảnh x1 (trên/dưới/trái/phải so với góc ảnh)
    def H(c): return c["bbox_x1"][3] - c["bbox_x1"][1]
    def T(c): return c["bbox_x1"][1]
    def L(c): return c["bbox_x1"][0]
    def R(c): return c["bbox_x1"][2]

    # dy_l cố định theo cf neo (cf nhỏ nhất dùng slot)
    leg_cf = charinfo_slot_to_cf(ci, "leg")
    dy_l, dx_l_anchor = {}, {}
    for l, cfs in leg_cf.items():
        c = comps[("leg", l)]
        cf0 = cfs[0]
        dy_l[l] = FEET - c["bbox_x1"][3] + ci[cf0][1][2]  # bbox bottom = cy - CI + dy + bottom -> = -11

    # chỉ số biến: dy_h 0..2, dy_b 3..19, dx_h 20..22, dx_b 23..39, dx_l 40..53
    nv = 3 + 17 + 3 + 17 + 14
    iyh = lambda s: s; iyb = lambda b: 3 + b; ixh = lambda s: 20 + s; ixb = lambda b: 23 + b; ixl = lambda l: 40 + l
    # biến phụ: e+,e-,d+,d-,ex+,ex-,gx+,gx- cho mỗi cf (8*ncf) + t (max|e|)
    n = nv + 8 * ncf + 1
    it = n - 1
    def aux(cf, k): return nv + cf * 8 + k

    A, lb, ub = [], [], []
    def row(coefs, lo, hi):
        r = np.zeros(n);
        for i, v in coefs: r[i] += v
        A.append(r); lb.append(lo); ub.append(hi)

    cost = np.zeros(n)
    cost[it] = a.M
    for cf in range(ncf):
        (hs, hx, hy), (ls, lx, ly), (bs, bx, by) = ci[cf][0], ci[cf][1], ci[cf][2]
        h, b, l = comps[("head", hs)], comps[("body", bs)], comps[("leg", ls)]
        # --- dọc ---
        # cằm đầu = -hy + dy_h + chin_off_h ; đường cằm thân = -by + dy_b + chin_off_b
        e_const = (-hy + h["chin_off"]) - (-by + b["chin_off"])
        # e = dy_h - dy_b + e_const = e+ - e-
        row([(iyh(hs), 1), (iyb(bs), -1), (aux(cf, 0), -1), (aux(cf, 1), 1)], -e_const, -e_const)
        row([(iyh(hs), 1), (iyb(bs), -1)], a.e_min - e_const, a.e_max - e_const)
        # d = (-by + dy_b + feet_off_b) - FEET
        d_const = -by + b["feet_off"] - FEET
        row([(iyb(bs), 1), (aux(cf, 2), -1), (aux(cf, 3), 1)], -d_const, -d_const)
        row([(iyb(bs), 1)], -a.d_max - d_const, a.d_max - d_const)
        # khe đầu-thân bbox = (-by + dy_b + T_b) - (-hy + dy_h + T_h + H_h)  in [khe_min, 0]
        k_const = (-by + T(b)) - (-hy + T(h) + H(h))
        row([(iyb(bs), 1), (iyh(hs), -1)], a.khe_min - k_const, 0 - k_const)
        # khe thân-chân = (-ly + dy_l + T_l) - (-by + dy_b + T_b + H_b) in [-11, 0]
        kl_const = (-ly + dy_l[ls] + T(l)) - (-by + T(b) + H(b))
        row([(iyb(bs), -1)], -11 - kl_const, KL_MAX - kl_const)
        # t >= |e|
        row([(it, 1), (aux(cf, 0), -1), (aux(cf, 1), -1)], 0, np.inf)
        # --- ngang ---
        # tâm mặt đầu = hx + dx_h + face_cx_h ; tâm mặt thân = bx + dx_b + face_cx_b
        fb = b["face_cx_off"] if b["face_cx_off"] is not None else b["center_off"]
        ex_const = (hx + h["face_cx_off"]) - (bx + fb)
        row([(ixh(hs), 1), (ixb(bs), -1), (aux(cf, 4), -1), (aux(cf, 5), 1)], -ex_const, -ex_const)
        # thân-chân ngang: cùng khung -> giữ đúng bố cục nguồn (tâm khung); khác khung -> căn tâm HÔNG
        same = (b["strip"], b["frame"]) == (l["strip"], l["frame"])
        if same:
            gx_const = (bx + b["center_off"]) - (lx + l["center_off"])
        else:
            gx_const = (bx + b["hip_cx_off"]) - (lx + l["hip_cx_off"])
        row([(ixb(bs), 1), (ixl(ls), -1), (aux(cf, 6), -1), (aux(cf, 7), 1)], -gx_const, -gx_const)
        for k in range(8):
            cost[aux(cf, k)] += (a.w_d if k in (2, 3) else 1.0) * w_cf.get(cf, 1.0)
    # neo ngang (khử bậc tự do tịnh tiến): tâm khung THÂN tại cf neo của slot ~ cx, phạt nhẹ (0.05)
    body_cf = charinfo_slot_to_cf(ci, "body")
    ig = n  # thêm biến phụ cho neo
    extra = 2 * 17
    n2 = n + extra
    A2 = [np.concatenate([r, np.zeros(extra)]) for r in A]
    A[:] = A2
    cost = np.concatenate([cost, np.full(extra, 0.05)])
    for b, cfs in body_cf.items():
        c = comps[("body", b)]; cf0 = cfs[0]
        r = np.zeros(n2); r[ixb(b)] = 1; r[ig + 2 * b] = -1; r[ig + 2 * b + 1] = 1
        A.append(r); k0 = -(ci[cf0][2][1] + c["center_off"]); lb.append(k0); ub.append(k0)
    n = n2

    pin = os.environ.get('PIN_HEAD0')
    if pin is not None:
        r = np.zeros(n); r[iyh(0)] = 1; A.append(r); lb.append(int(pin)); ub.append(int(pin))
        print('GHIM dy_h[0] =', pin)
    integrality = np.zeros(n); integrality[:nv] = 1
    lo = np.full(n, -np.inf); hi = np.full(n, np.inf)
    lo[:nv] = -128; hi[:nv] = 127
    lo[nv:] = 0
    res = milp(c=cost, constraints=LinearConstraint(np.array(A), lb, ub), integrality=integrality, bounds=Bounds(lo, hi))
    if res.x is None:
        print("VO NGHIEM:", res.message); return 2
    x = np.round(res.x).astype(int)
    dy_h = [x[iyh(s)] for s in range(3)]; dy_b = [x[iyb(b)] for b in range(17)]
    dx_h = [x[ixh(s)] for s in range(3)]; dx_b = [x[ixb(b)] for b in range(17)]; dx_l = [x[ixl(l)] for l in range(14)]
    print("dy_h", dy_h, "dx_h", dx_h)
    print("dy_b", dy_b); print("dx_b", dx_b)
    print("dy_l", [dy_l[l] for l in range(14)]); print("dx_l", dx_l)

    # bảng dự đoán
    print("cf | e(cam) d(than) ex gx | khe dau-than | khe than-chan | chan | cao | than/chan")
    worst = 0; sum_e = sum_d = 0.0
    for cf in range(ncf):
        (hs, hx, hy), (ls, lx, ly), (bs, bx, by) = ci[cf][0], ci[cf][1], ci[cf][2]
        h, b, l = comps[("head", hs)], comps[("body", bs)], comps[("leg", ls)]
        e = (-hy + dy_h[hs] + h["chin_off"]) - (-by + dy_b[bs] + b["chin_off"])
        d = (-by + dy_b[bs] + b["feet_off"]) - FEET
        fb = b["face_cx_off"] if b["face_cx_off"] is not None else b["center_off"]
        ex = (hx + dx_h[hs] + h["face_cx_off"]) - (bx + dx_b[bs] + fb)
        if (b["strip"], b["frame"]) == (l["strip"], l["frame"]):
            gx = (bx + dx_b[bs] + b["center_off"]) - (lx + dx_l[ls] + l["center_off"])
        else:
            gx = (bx + dx_b[bs] + b["hip_cx_off"]) - (lx + dx_l[ls] + l["hip_cx_off"])
        khe = (-by + dy_b[bs] + T(b)) - (-hy + dy_h[hs] + T(h) + H(h))
        khel = (-ly + dy_l[ls] + T(l)) - (-by + dy_b[bs] + T(b) + H(b))
        chan = -ly + dy_l[ls] + l["bbox_x1"][3]
        cao = chan - (-hy + dy_h[hs] + T(h))
        worst = max(worst, abs(e)); sum_e += abs(e); sum_d += abs(d)
        print("%2d | %+5.1f %+5.1f %+5.1f %+5.1f | %+3d | %+3d | %+3d | %2d | %s/%s" % (
            cf, e, d, ex, gx, khe, khel, chan, cao, f"{b['strip']}{b['frame']}", f"{l['strip']}{l['frame']}"))
    print("max|e| = %.1f ; TB|e| = %.2f ; TB|d| = %.2f ; muc tieu = %.2f" % (worst, sum_e / ncf, sum_d / ncf, res.fun))

    # vá part.bin: records 30 (đầu) 57 (thân) 58 (chân)
    data = bytearray(PART_GOC.read_bytes())
    nrec = struct.unpack_from(">h", data, 0)[0]
    off = 2; offs = {}
    for r in range(nrec):
        t = data[off]; offs[r] = (t, off + 1); off += 1 + ghep.SO_KHUNG[t] * 4
    def patch(pid, frames):
        t, o = offs[pid]
        assert ghep.SO_KHUNG[t] == len(frames), (pid, t, len(frames))
        for k, (sid, dx, dy) in enumerate(frames):
            struct.pack_into(">hbb", data, o + k * 4, sid, int(dx), int(dy))
    patch(30, [(27041 + s, dx_h[s], dy_h[s]) for s in range(3)])
    patch(57, [(27044 + b, dx_b[b], dy_b[b]) for b in range(17)])
    patch(58, [(27061 + l, dx_l[l], dy_l[l]) for l in range(14)])
    Path(a.out_part).write_bytes(bytes(data))
    json.dump(dict(dy_h=[int(v) for v in dy_h], dy_b=[int(v) for v in dy_b], dy_l=[int(dy_l[l]) for l in range(14)],
                   dx_h=[int(v) for v in dx_h], dx_b=[int(v) for v in dx_b], dx_l=[int(v) for v in dx_l]),
              open(str(a.out_part) + ".json", "w"), indent=1)
    print("da ghi", a.out_part)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
