# -*- coding: utf-8 -*-
"""Tính lại khe đầu–thân / thân–chân từng cf CHỈ từ số trong mapping-id.json + charinfo.json,
đối chiếu với bảng đo tay, rồi giải bài toán khả thi (cửa sổ dy đầu; MILP dy đầu + dy thân)."""
import json, sys
import numpy as np
sys.stdout.reconfigure(encoding='utf-8')

m = json.load(open(r'D:\NRO_GoiBanGiao\docs\my-thuat\mapping-id.json', encoding='utf-8'))
G = m['build_inputs']['player_component_geometry']['components']
CI = json.load(open(r'D:\NRO_GoiBanGiao\tools\nhan-vat\charinfo.json', encoding='utf-8'))
S = m['build_inputs']['player_component_geometry']['uniform_scale_x4']  # 128/340

# Kích cỡ x1 = size_x4/4 làm tròn (đã kiểm bằng kiem-bbox.py: cao x1 thật = round(h4/4), bbox đục cao = cả ảnh)
comp = {}
for c in G:
    h1 = round(c['size_x4'][1] / 4)
    comp[(c['component'], c['slot_index'])] = dict(id=c['id'], H=h1, dy=c['part_offset']['dy'],
        top_rel=c['placement_x1']['top_rel'], action=c['action'], aslot=c['action_slot'],
        hrect=c['source_rect'][3])
# kiểm: cao x1 thật đo được từ ảnh (kiem-bbox.py)
H_THAT = {27041:26,27042:26,27043:26,27044:30,27045:29,27046:29,27047:24,27048:20,27049:22,27050:23,
          27051:23,27052:24,27053:23,27054:23,27055:28,27056:27,27057:26,27058:24,27059:10,27060:20,
          27061:21,27062:21,27063:21,27064:17,27065:15,27066:16,27067:16,27068:16,27069:16,27070:16,
          27071:16,27072:20,27073:19,27074:18}
for k, v in comp.items():
    assert v['H'] == H_THAT[v['id']], (k, v['H'], H_THAT[v['id']])

# Bảng đo tay từ prompt (ghep-khung.py): cf -> (khe dau-than, khe than-chan, cao tong)
DO_TAY = {0:(-8,-6,62),1:(-8,-5,63),2:(-9,-6,61),3:(-2,-5,60),4:(5,-4,62),5:(3,-5,62),6:(1,-5,61),
 7:(4,-5,64),8:(3,-10,63),9:(-3,-3,63),10:(-2,-6,65),11:(2,-8,62),12:(2,-5,63),13:(2,-5,62),14:(-6,-2,62),
 15:(-9,-2,60),16:(-5,-3,60),17:(4,-10,59),18:(2,-8,60),19:(-1,-4,61),20:(19,-10,61),21:(2,-8,59),
 22:(2,-11,62),23:(-6,-2,64),24:(-6,-7,59),25:(5,-8,61),26:(2,-4,63),27:(-6,-1,63),28:(-9,-1,61),
 29:(-5,-2,61),30:(-1,-3,62),31:(19,-9,62),32:(2,-6,61)}

def hinh_hoc(cf, dyh, dyb, dyl, Hb=None):
    """Trả (khe đầu-thân, khe thân-chân, bàn chân so cy, cao tổng) cho một cf với dy cho trước.
    dyh/dyb/dyl: dict slot->dy. Hb: dict slot->cao thân (ghi đè, dùng khi thay khung slot 15)."""
    (ih, _, cih), (il, _, cil), (ib, _, cib) = CI[cf][0], CI[cf][1], CI[cf][2]
    h, b, l = comp[('head', ih)], comp[('body', ib)], comp[('leg', il)]
    Hbody = (Hb or {}).get(ib, b['H'])
    head_top = -cih + dyh[ih]; head_bot = head_top + h['H']
    body_top = -cib + dyb[ib]; body_bot = body_top + Hbody
    leg_top = -cil + dyl[il]; leg_bot = leg_top + l['H']
    return body_top - head_bot, leg_top - body_bot, leg_bot, leg_bot - head_top

dyh0 = {s: comp[('head', s)]['dy'] for s in range(3)}
dyb0 = {s: comp[('body', s)]['dy'] for s in range(17)}
dyl0 = {s: comp[('leg', s)]['dy'] for s in range(14)}

print("=== (1) TÍNH LẠI TỪ SỐ vs ĐO TAY ===")
print("cf | slot h/b/l | khe d-t  (do) | khe t-c (do) | chan-cy | cao (do) | khop?")
lech = 0
for cf in range(33):
    kdt, ktc, chan, cao = hinh_hoc(cf, dyh0, dyb0, dyl0)
    d = DO_TAY[cf]
    ok = (kdt, ktc, cao) == d
    lech += (not ok)
    print("%2d | %d/%2d/%2d | %+3d (%+3d) | %+3d (%+3d) | %+3d | %2d (%2d) | %s" % (
        cf, CI[cf][0][0], CI[cf][2][0], CI[cf][1][0], kdt, d[0], ktc, d[1], chan, cao, d[2], "OK" if ok else "LECH"))
print("so cf lech so voi do tay:", lech)

print()
print("=== (2) CỬA SỔ dy ĐẦU (chỉ đổi dy đầu, thân/chân giữ nguyên), yêu cầu khe đầu–thân trong [-9, 0] ===")
# khe(cf) = A(cf) - dy_h  với A = CI_h - CI_b + dy_b - H_h  → dy_h ∈ [A, A+9]
for hs in range(3):
    cfs = [cf for cf in range(33) if CI[cf][0][0] == hs]
    A = {}
    for cf in cfs:
        kdt, *_ = hinh_hoc(cf, dyh0, dyb0, dyl0)
        A[cf] = kdt + dyh0[hs]
    lo = max(A.values()); hi = min(A.values()) + 9
    print("dau slot %d (id %d, dy hien tai %d): cf dung = %s" % (hs, comp[('head', hs)]['id'], dyh0[hs], cfs))
    print("   can dy_h >= max A = %d (cf %s) va dy_h <= min A + 9 = %d (cf %s) -> %s" % (
        lo, [cf for cf in cfs if A[cf] == lo], hi, [cf for cf in cfs if A[cf] == min(A.values())],
        "KHA THI [%d, %d]" % (lo, hi) if lo <= hi else "VO NGHIEM (bien do khe = %d > 9)" % (lo - (hi - 9))))
    cfs2 = [cf for cf in cfs if cf not in (20, 31)]
    if cfs2 != cfs:
        lo2 = max(A[cf] for cf in cfs2); hi2 = min(A[cf] for cf in cfs2) + 9
        print("   bo cf 20/31 (khung nam): can dy_h in [%d, %d] -> %s" % (lo2, hi2,
              "KHA THI" if lo2 <= hi2 else "VAN VO NGHIEM (bien do %d > 9)" % (lo2 - (hi2 - 9))))
    # minimax: dịch Δ tối thiểu hoá vi phạm lớn nhất
    best = None
    for delta in range(-15, 16):
        viol = max(max(0, (A[cf] - (dyh0[hs] + delta))) + max(0, -9 - (A[cf] - (dyh0[hs] + delta))) for cf in cfs2)
        if best is None or viol < best[1]: best = (delta, viol)
    print("   dich dy dau tot nhat (bo 20/31, minimax): delta=%+d -> vi pham lon nhat %d px" % best)
    for delta in sorted({0, best[0], best[0] + 1, best[0] - 1, 4, 5}):
        ke = {cf: A[cf] - (dyh0[hs] + delta) for cf in cfs2}
        ho = {cf: k for cf, k in ke.items() if k > 0}
        chong = {cf: k for cf, k in ke.items() if k < -9}
        print("      delta %+d: HO %s | CHONG QUA %s" % (delta, ho, chong))

print()
print("=== (3) Thay khung NẰM ở slot thân 15 (id 27059) — riêng cf 20/31 ===")
for cf in (20, 31):
    (ih, _, cih), (il, _, cil), (ib, _, cib) = CI[cf][0], CI[cf][1], CI[cf][2]
    print("cf %d: dau idx %d, than idx %d, chan idx %d (CI dy: dau %d, than %d, chan %d)" % (cf, ih, ib, il, cih, cib, cil))
# chỉ đổi dy thân 15 (giữ ảnh nằm cao 10): khe d-t <= 0 cần dy_b15 <= ?, khi đó khe t-c = ?
for dyb15 in (-8, -20, -27, -30, -36):
    dyb = dict(dyb0); dyb[15] = dyb15
    r20 = hinh_hoc(20, dyh0, dyb, dyl0); r31 = hinh_hoc(31, dyh0, dyb, dyl0)
    print("  anh nam giu nguyen, dy_b15=%d: cf20 khe d-t %+d, t-c %+d | cf31 khe d-t %+d, t-c %+d" % (dyb15, r20[0], r20[1], r31[0], r31[1]))
print("  -> ha than xuong de dong khe dau thi mo khe than-chan: anh cao 10 px khong the lap ca hai.")
for cand in (14, 16, 13, 7, 0):
    Hc = comp[('body', cand)]['H']
    # tìm dy_b15 tốt nhất (khe d-t và t-c cùng trong [-9,0]) cho khung chép từ slot cand
    ok = []
    for dyb15 in range(-60, 10):
        dyb = dict(dyb0); dyb[15] = dyb15
        rs = [hinh_hoc(cf, dyh0, dyb, dyl0, Hb={15: Hc}) for cf in (20, 31)]
        if all(-9 <= r[0] <= 0 and -9 <= r[1] <= 0 for r in rs): ok.append(dyb15)
    print("  chep khung slot %2d (id %d, cao %2d) vao slot 15: dy_b15 kha thi = %s" % (cand, comp[('body', cand)]['id'], Hc, ok))

print()
print("=== (4) MILP: dy đầu (3) + dy thân (17) tự do, dy chân GIỮ NGUYÊN (bàn chân không đổi) ===")
from scipy.optimize import milp, LinearConstraint, Bounds
def giai(lo_hb, hi_hb, lo_bl, hi_bl, Hb_override=None, leg_free=0, target_hb=-8, target_bl=-6, chi_cf=None):
    """Biến: dyh[3], dyb[17], dyl[14], u_hb[33], u_bl[33] (|khe - target|). Ràng buộc khe trong cửa sổ.
    leg_free: cho dy chân lệch tối đa bấy nhiêu px so hiện tại (0 = cố định)."""
    cfs = list(range(33)) if chi_cf is None else chi_cf
    nH, nB, nL, nC = 3, 17, 14, len(cfs)
    n = nH + nB + nL + 2 * nC
    iH = lambda s: s; iB = lambda s: nH + s; iL = lambda s: nH + nB + s
    iU = lambda k: nH + nB + nL + k; iV = lambda k: nH + nB + nL + nC + k
    cobj = np.zeros(n); rows = []; lo = []; hi = []
    for k, cf in enumerate(cfs):
        (ih, _, cih), (il, _, cil), (ib, _, cib) = CI[cf][0], CI[cf][1], CI[cf][2]
        Hh = comp[('head', ih)]['H']; Hbody = (Hb_override or {}).get(ib, comp[('body', ib)]['H']); Hl = comp[('leg', il)]['H']
        # khe_hb = (cih - cib - Hh) + dyb - dyh
        c_hb = cih - cib - Hh
        r = np.zeros(n); r[iB(ib)] = 1; r[iH(ih)] = -1; rows.append(r); lo.append(lo_hb - c_hb); hi.append(hi_hb - c_hb)
        # khe_bl = (-cil + cib - Hbody) + dyl - dyb
        c_bl = -cil + cib - Hbody
        r = np.zeros(n); r[iL(il)] = 1; r[iB(ib)] = -1; rows.append(r); lo.append(lo_bl - c_bl); hi.append(hi_bl - c_bl)
        # |khe_hb - target| <= u
        for sg in (1, -1):
            r = np.zeros(n); r[iB(ib)] = sg; r[iH(ih)] = -sg; r[iU(k)] = -1; rows.append(r); lo.append(-np.inf); hi.append(-sg * (c_hb - target_hb))
            r = np.zeros(n); r[iL(il)] = sg; r[iB(ib)] = -sg; r[iV(k)] = -1; rows.append(r); lo.append(-np.inf); hi.append(-sg * (c_bl - target_bl))
        cobj[iU(k)] = 1; cobj[iV(k)] = 1
    lb = np.full(n, -128.0); ub = np.full(n, 127.0)
    for s in range(nL):
        lb[iL(s)] = dyl0[s] - leg_free; ub[iL(s)] = dyl0[s] + leg_free
    for k in range(nC): lb[iU(k)] = 0; lb[iV(k)] = 0
    integ = np.zeros(n); integ[:nH + nB + nL] = 1
    res = milp(cobj, constraints=LinearConstraint(np.array(rows), lo, hi), bounds=Bounds(lb, ub), integrality=integ)
    if res.status != 0:
        return None
    x = np.round(res.x).astype(int)
    return ({s: int(x[iH(s)]) for s in range(nH)}, {s: int(x[iB(s)]) for s in range(nB)}, {s: int(x[iL(s)]) for s in range(nL)}, res.fun)

def in_nghiem(ten, sol, Hb=None):
    if sol is None:
        print(ten, "-> VO NGHIEM"); return
    dyh, dyb, dyl, f = sol
    print(ten, "-> CO NGHIEM, tong |khe - muc tieu| = %.0f" % f)
    print("   dy dau  moi:", dyh, " (cu", dyh0, ")")
    print("   dy than moi:", {s: dyb[s] for s in range(17)})
    print("   dy than doi:", {s: dyb[s] - dyb0[s] for s in range(17) if dyb[s] != dyb0[s]})
    if dyl != dyl0: print("   dy chan doi:", {s: dyl[s] - dyl0[s] for s in range(14) if dyl[s] != dyl0[s]})
    ke = [hinh_hoc(cf, dyh, dyb, dyl, Hb) for cf in range(33)]
    print("   khe dau-than tung cf :", " ".join("%+d" % k[0] for k in ke))
    print("   khe than-chan tung cf:", " ".join("%+d" % k[1] for k in ke))
    print("   cao tong tung cf     :", " ".join("%d" % k[3] for k in ke))
    print("   ban chan-cy          :", sorted(set(k[2] for k in ke)))

in_nghiem("S2  giu anh NAM slot 15, khe hb & bl trong [-9,0]", giai(-9, 0, -9, 0))
in_nghiem("S2b giu anh NAM slot 15, bo cf 20/31 khoi rang buoc", giai(-9, 0, -9, 0, chi_cf=[cf for cf in range(33) if cf not in (20, 31)]))
Hb14 = {15: comp[('body', 14)]['H']}
in_nghiem("S3  slot 15 = chep hit0 (cao 24), khe hb & bl trong [-9,0]", giai(-9, 0, -9, 0, Hb_override=Hb14), Hb14)
in_nghiem("S3' slot 15 = chep hit0, khe hb [-9,-3], bl [-9,-2] (chat hon)", giai(-9, -3, -9, -2, Hb_override=Hb14), Hb14)
in_nghiem("S3'' slot 15 = chep hit0, khe hb [-9,-5], bl [-9,-3] (rat chat)", giai(-9, -5, -9, -3, Hb_override=Hb14), Hb14)
in_nghiem("S4  nhu S3' nhung cho dy chan lech <=1 (ban chan +-1)", giai(-9, -3, -9, -2, Hb_override=Hb14, leg_free=1), Hb14)

print()
print("=== (5) Muc tieu 'khe thiet ke' theo khung nguon cua THAN (chong 14% cao khung x s/4) ===")
for cf in range(33):
    ib = CI[cf][2][0]; b = comp[('body', ib)]
    Hframe = b['hrect'] / 0.48
    T = -0.14 * Hframe * S / 4
    kdt = hinh_hoc(cf, dyh0, dyb0, dyl0)[0]
    ih = CI[cf][0][0]
    Hhead_frame = 0.42 * Hframe * S / 4
    print("cf %2d than %s%d (H khung nguon %.0f) muc tieu khe %.1f | hien tai %+d | sai %+.1f | dau idle cao 26 vs dau 'dung khung' %.1f" % (
        cf, b['action'], b['aslot'], Hframe, T, kdt, kdt - T, Hhead_frame))
