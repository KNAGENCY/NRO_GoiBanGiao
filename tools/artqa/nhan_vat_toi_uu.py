# -*- coding: utf-8 -*-
"""nhan_vat_toi_uu.py — MILP chọn dx/dy nguyên cho 34 mảnh người chơi (V1, 05/09/2026).

VÌ SAO PHẢI GIẢI TỐI ƯU, KHÔNG CHỈNH TAY
----------------------------------------
Client vẽ nhân vật bằng bộ xương cứng ``CharInfo`` (``Char.cs:483``, 33 cf): mỗi cf
chọn (slot đầu, slot chân, slot thân) và một offset ``(x, y)`` cho từng mảnh. Chỉ có
3 đầu dùng chung cho 33 cf, và mỗi mảnh chỉ có MỘT ``dx/dy`` trong ``part.bin`` dù nó
phục vụ nhiều cf. Đổi 3 dy đầu để khớp mọi cf là VÔ NGHIỆM (cửa sổ 9 px, biên độ cần
12–14 — ``thu-nghiem-khe-dau-than/phan-tich-cua-so-dy.py``); nhưng 17 dy thân (mỗi slot
chỉ phục vụ 1–3 cf) đủ bậc tự do. Đây là bài toán số nguyên có ràng buộc cứng → MILP.

Phép vẽ client (hướng phải, neo TOP|LEFT), k = 0 đầu / 1 chân / 2 thân:
    left = cx + CI[cf][k][1] + dx[slot]
    top  = cy - CI[cf][k][2] + dy[slot]          (top là MÉP TRÊN ảnh)

Ẩn số: dy đầu (3), dy thân (17), dx đầu (3), dx thân (17), dx chân (14). **dy chân CỐ
ĐỊNH**: đáy vùng đục của chân = ``cy + ban_chan_x1`` tại cf neo (cf nhỏ nhất dùng slot)
— giữ đúng quy ước hiện hành (-11) để nhân vật không đổi cao độ so với mặt đất.

Sai số giải phẫu từng cf (đơn vị x1):
    e  = cằm đầu − đường cằm khung thân        (muốn 0; e>0 đầu tụt vào cổ, e<0 hở)
    d  = bàn chân "ảo" của khung thân − bàn chân (thân trượt so với chân)
    ex = tâm mặt đầu − tâm mặt khung thân       (ngang)
    gx = tâm thân − tâm chân                    (ngang; cùng khung: tâm khung, khác khung: tâm hông)
Ràng buộc cứng: khe đầu–thân (bbox) ∈ ``khe_dau_than``, khe thân–chân ∈ ``khe_than_chan``
(siết trên −2 để x4 không hở khi làm tròn), e ∈ ``e``, |d| ≤ ``d_max``, ghim dy đầu
(``ghim_dy_dau``; slot 0 = −38 vì −39 vô nghiệm — điều phối đã thử).
Mục tiêu: M·max|e| + Σ w_cf(|e| + w_d|d| + |ex| + |gx|) + w_neo·Σ|tâm thân − cx|.

Toán chép từ ``tools/nhan-vat/thu-nghiem-khe-dau-than/toi-uu-dxdy.py`` (bản đã sinh
``part-V1-siet-ghim38.bin``); thêm ``kiem_nghiem``/``bang_cf`` để mapping tự chứng minh
nghiệm đã lưu thoả ràng buộc mà không cần scipy giải lại.
"""

from __future__ import annotations

import json
from dataclasses import asdict, dataclass, field
from typing import Any

import numpy as np

COMPONENT_CHARINFO_SLOT = {"head": 0, "leg": 1, "body": 2}
SO_SLOT = {"head": 3, "body": 17, "leg": 14}


class VoNghiem(ValueError):
    """MILP không có nghiệm với ràng buộc đã chốt — báo, không tự nới."""


@dataclass(frozen=True)
class ThamSoMILP:
    """Ràng buộc + trọng số; ghi nguyên vào mapping để tái lập."""

    ban_chan_x1: int = -11
    khe_dau_than: tuple[int, int] = (-9, 0)
    khe_than_chan: tuple[int, int] = (-11, -2)
    e: tuple[float, float] = (-2.0, 1.0)
    d_max: float = 5.5
    ghim_dy_dau: dict[int, int] = field(default_factory=lambda: {0: -38})
    w_d: float = 0.5
    M: float = 50.0
    w_cf: dict[int, float] = field(
        default_factory=lambda: {0: 4.0, 1: 4.0, 2: 2.0, 3: 2.0, 4: 2.0, 5: 2.0, 6: 2.0, 7: 2.0, 8: 2.0}
    )
    w_neo: float = 0.05

    def json(self) -> dict[str, Any]:
        d = asdict(self)
        d["khe_dau_than"] = list(self.khe_dau_than)
        d["khe_than_chan"] = list(self.khe_than_chan)
        d["e"] = list(self.e)
        d["ghim_dy_dau"] = {str(k): int(v) for k, v in self.ghim_dy_dau.items()}
        d["w_cf"] = {str(k): float(v) for k, v in self.w_cf.items()}
        return d

    @classmethod
    def tu_json(cls, d: dict[str, Any]) -> "ThamSoMILP":
        return cls(
            ban_chan_x1=int(d["ban_chan_x1"]),
            khe_dau_than=(int(d["khe_dau_than"][0]), int(d["khe_dau_than"][1])),
            khe_than_chan=(int(d["khe_than_chan"][0]), int(d["khe_than_chan"][1])),
            e=(float(d["e"][0]), float(d["e"][1])),
            d_max=float(d["d_max"]),
            ghim_dy_dau={int(k): int(v) for k, v in d.get("ghim_dy_dau", {}).items()},
            w_d=float(d["w_d"]),
            M=float(d["M"]),
            w_cf={int(k): float(v) for k, v in d.get("w_cf", {}).items()},
            w_neo=float(d["w_neo"]),
        )

    def __eq__(self, other):
        if not isinstance(other, ThamSoMILP):
            return NotImplemented
        return json.dumps(self.json(), sort_keys=True) == json.dumps(other.json(), sort_keys=True)

    __hash__ = None  # dict field → không hash


# --------------------------------------------------------------------------- bảng tra
def _bang_manh(components) -> dict[tuple[str, int], dict[str, Any]]:
    bang = {}
    for c in components:
        key = (c["component"], int(c["slot_index"]))
        if key in bang:
            raise ValueError(f"trùng mảnh {key}")
        moc = dict(c["moc"])
        # khung không có cằm (nằm): coi mép trên mảnh + 2 px là "đường cằm"; tâm mặt = tâm khung
        if moc.get("chin_off") is None:
            moc["chin_off"] = 2.0
        if moc.get("face_cx_off") is None:
            moc["face_cx_off"] = moc["center_off"]
        bang[key] = {**c, "moc": moc}
    return bang


def _so_slot(charinfo, component: str) -> int:
    k = COMPONENT_CHARINFO_SLOT[component]
    return max(entry[k][0] for entry in charinfo) + 1


def _cf_theo_slot(charinfo, component: str) -> dict[int, list[int]]:
    k = COMPONENT_CHARINFO_SLOT[component]
    result: dict[int, list[int]] = {}
    for cf, entry in enumerate(charinfo):
        result.setdefault(entry[k][0], []).append(cf)
    return result


def _T(c):  # mép trên vùng đục (x1) so với góc ảnh
    return c["moc"]["bbox_x1"][1]


def _H(c):
    return c["moc"]["bbox_x1"][3] - c["moc"]["bbox_x1"][1]


def _B(c):
    return c["moc"]["bbox_x1"][3]


def _cung_khung(a, b) -> bool:
    return (a["action"], a["action_slot"]) == (b["action"], b["action_slot"])


def dy_chan_co_dinh(components, charinfo, tham_so: ThamSoMILP) -> list[int]:
    """dy từng slot chân: đáy vùng đục = cy + ban_chan tại cf neo (cf nhỏ nhất dùng slot)."""

    bang = _bang_manh(components)
    n = _so_slot(charinfo, "leg")
    leg_cf = _cf_theo_slot(charinfo, "leg")
    result = []
    for l in range(n):
        c = bang[("leg", l)]
        cf0 = leg_cf[l][0]
        # bbox bottom trên màn = cy - ly + dy + B  ⇒ = cy + ban_chan  ⇒ dy = ban_chan - B + ly
        result.append(int(tham_so.ban_chan_x1 - _B(c) + charinfo[cf0][1][2]))
    return result


# --------------------------------------------------------------------------- giải
def giai(components, charinfo, tham_so: ThamSoMILP | None = None) -> dict[str, Any]:
    """Giải MILP; trả ``{"head": {"dx": [...], "dy": [...]}, "body": ..., "leg": ..., "muc_tieu": f}``.

    Raise ``VoNghiem`` khi không có nghiệm với ràng buộc đã chốt.
    """

    from scipy.optimize import Bounds, LinearConstraint, milp

    ts = tham_so or ThamSoMILP()
    bang = _bang_manh(components)
    ncf = len(charinfo)
    nh, nb, nl = (_so_slot(charinfo, k) for k in ("head", "body", "leg"))
    for comp, n in (("head", nh), ("body", nb), ("leg", nl)):
        for s in range(n):
            if (comp, s) not in bang:
                raise ValueError(f"thiếu mảnh {comp} slot {s}")
    dy_l = dy_chan_co_dinh(components, charinfo, ts)

    # chỉ số biến: dy_h, dy_b, dx_h, dx_b, dx_l | 8 biến phụ mỗi cf | t = max|e| | 2 neo mỗi slot thân
    nv = nh + nb + nh + nb + nl
    iyh = lambda s: s
    iyb = lambda b: nh + b
    ixh = lambda s: nh + nb + s
    ixb = lambda b: nh + nb + nh + b
    ixl = lambda l: nh + nb + nh + nb + l
    aux = lambda cf, k: nv + cf * 8 + k
    it = nv + 8 * ncf
    ineo = lambda b, k: it + 1 + 2 * b + k
    n = it + 1 + 2 * nb

    A, lb, ub = [], [], []

    def row(coefs, lo, hi):
        r = np.zeros(n)
        for i, v in coefs:
            r[i] += v
        A.append(r)
        lb.append(lo)
        ub.append(hi)

    cost = np.zeros(n)
    cost[it] = ts.M
    for cf in range(ncf):
        (hs, hx, hy), (ls, lx, ly), (bs, bx, by) = charinfo[cf][0], charinfo[cf][1], charinfo[cf][2]
        h, b, l = bang[("head", hs)], bang[("body", bs)], bang[("leg", ls)]
        # --- dọc ---
        e_const = (-hy + h["moc"]["chin_off"]) - (-by + b["moc"]["chin_off"])
        # e = dy_h - dy_b + e_const = e+ - e-
        row([(iyh(hs), 1), (iyb(bs), -1), (aux(cf, 0), -1), (aux(cf, 1), 1)], -e_const, -e_const)
        row([(iyh(hs), 1), (iyb(bs), -1)], ts.e[0] - e_const, ts.e[1] - e_const)
        # d = (-by + dy_b + feet_off_b) - ban_chan
        d_const = -by + b["moc"]["feet_off"] - ts.ban_chan_x1
        row([(iyb(bs), 1), (aux(cf, 2), -1), (aux(cf, 3), 1)], -d_const, -d_const)
        row([(iyb(bs), 1)], -ts.d_max - d_const, ts.d_max - d_const)
        # khe đầu-thân (bbox) = (-by + dy_b + T_b) - (-hy + dy_h + T_h + H_h)
        k_const = (-by + _T(b)) - (-hy + _T(h) + _H(h))
        row([(iyb(bs), 1), (iyh(hs), -1)], ts.khe_dau_than[0] - k_const, ts.khe_dau_than[1] - k_const)
        # khe thân-chân = (-ly + dy_l + T_l) - (-by + dy_b + T_b + H_b)
        kl_const = (-ly + dy_l[ls] + _T(l)) - (-by + _T(b) + _H(b))
        row([(iyb(bs), -1)], ts.khe_than_chan[0] - kl_const, ts.khe_than_chan[1] - kl_const)
        # t >= |e|
        row([(it, 1), (aux(cf, 0), -1), (aux(cf, 1), -1)], 0, np.inf)
        # --- ngang ---
        ex_const = (hx + h["moc"]["face_cx_off"]) - (bx + b["moc"]["face_cx_off"])
        row([(ixh(hs), 1), (ixb(bs), -1), (aux(cf, 4), -1), (aux(cf, 5), 1)], -ex_const, -ex_const)
        if _cung_khung(b, l):
            gx_const = (bx + b["moc"]["center_off"]) - (lx + l["moc"]["center_off"])
        else:
            gx_const = (bx + b["moc"]["hip_cx_off"]) - (lx + l["moc"]["hip_cx_off"])
        row([(ixb(bs), 1), (ixl(ls), -1), (aux(cf, 6), -1), (aux(cf, 7), 1)], -gx_const, -gx_const)
        w = ts.w_cf.get(cf, 1.0)
        for k in range(8):
            cost[aux(cf, k)] += (ts.w_d if k in (2, 3) else 1.0) * w
    # neo ngang (khử bậc tự do tịnh tiến): tâm khung THÂN tại cf neo của slot ≈ cx, phạt nhẹ
    body_cf = _cf_theo_slot(charinfo, "body")
    for b, cfs in body_cf.items():
        c = bang[("body", b)]
        cf0 = cfs[0]
        k0 = -(charinfo[cf0][2][1] + c["moc"]["center_off"])
        row([(ixb(b), 1), (ineo(b, 0), -1), (ineo(b, 1), 1)], k0, k0)
        cost[ineo(b, 0)] = ts.w_neo
        cost[ineo(b, 1)] = ts.w_neo
    for slot, gia_tri in ts.ghim_dy_dau.items():
        if not 0 <= int(slot) < nh:
            raise ValueError(f"ghim dy đầu slot {slot} ngoài 0..{nh - 1}")
        row([(iyh(int(slot)), 1)], int(gia_tri), int(gia_tri))

    integrality = np.zeros(n)
    integrality[:nv] = 1
    lo = np.full(n, -np.inf)
    hi = np.full(n, np.inf)
    lo[:nv] = -128
    hi[:nv] = 127
    lo[nv:] = 0
    res = milp(
        c=cost,
        constraints=LinearConstraint(np.array(A), lb, ub),
        integrality=integrality,
        bounds=Bounds(lo, hi),
    )
    if res.x is None:
        raise VoNghiem(f"MILP vô nghiệm với ràng buộc đã chốt: {res.message} | tham số {ts.json()}")
    x = np.round(res.x).astype(int)
    return {
        "head": {"dx": [int(x[ixh(s)]) for s in range(nh)], "dy": [int(x[iyh(s)]) for s in range(nh)]},
        "body": {"dx": [int(x[ixb(b)]) for b in range(nb)], "dy": [int(x[iyb(b)]) for b in range(nb)]},
        "leg": {"dx": [int(x[ixl(l)]) for l in range(nl)], "dy": [int(v) for v in dy_l]},
        "muc_tieu": float(res.fun),
    }


# --------------------------------------------------------------------------- kiểm
def bang_cf(components, charinfo, nghiem, tham_so: ThamSoMILP | None = None) -> list[dict[str, Any]]:
    """Dự đoán từng cf theo phép vẽ client (x1): e, d, ex, gx, khe, bàn chân, cao tổng."""

    ts = tham_so or ThamSoMILP()
    bang = _bang_manh(components)
    rows = []
    for cf, entry in enumerate(charinfo):
        (hs, hx, hy), (ls, lx, ly), (bs, bx, by) = entry[0], entry[1], entry[2]
        h, b, l = bang[("head", hs)], bang[("body", bs)], bang[("leg", ls)]
        dy_h, dy_b, dy_l = nghiem["head"]["dy"][hs], nghiem["body"]["dy"][bs], nghiem["leg"]["dy"][ls]
        dx_h, dx_b, dx_l = nghiem["head"]["dx"][hs], nghiem["body"]["dx"][bs], nghiem["leg"]["dx"][ls]
        e = (-hy + dy_h + h["moc"]["chin_off"]) - (-by + dy_b + b["moc"]["chin_off"])
        d = (-by + dy_b + b["moc"]["feet_off"]) - ts.ban_chan_x1
        ex = (hx + dx_h + h["moc"]["face_cx_off"]) - (bx + dx_b + b["moc"]["face_cx_off"])
        if _cung_khung(b, l):
            gx = (bx + dx_b + b["moc"]["center_off"]) - (lx + dx_l + l["moc"]["center_off"])
        else:
            gx = (bx + dx_b + b["moc"]["hip_cx_off"]) - (lx + dx_l + l["moc"]["hip_cx_off"])
        khe = (-by + dy_b + _T(b)) - (-hy + dy_h + _T(h) + _H(h))
        khel = (-ly + dy_l + _T(l)) - (-by + dy_b + _T(b) + _H(b))
        chan = -ly + dy_l + _B(l)
        cao = chan - (-hy + dy_h + _T(h))
        rows.append(
            {
                "cf": cf,
                "e": round(float(e), 3),
                "d": round(float(d), 3),
                "ex": round(float(ex), 3),
                "gx": round(float(gx), 3),
                "khe_dau_than": int(khe),
                "khe_than_chan": int(khel),
                "chan": int(chan),
                "cao": int(cao),
                "than": f"{b['action']}{b['action_slot']}",
                "chan_khung": f"{l['action']}{l['action_slot']}",
            }
        )
    return rows


def kiem_nghiem(components, charinfo, nghiem, tham_so: ThamSoMILP | None = None) -> list[str]:
    """Danh sách vi phạm (rỗng = nghiệm hợp lệ). Không giải lại — chỉ thay số vào ràng buộc."""

    ts = tham_so or ThamSoMILP()
    loi: list[str] = []
    for phan in ("head", "body", "leg"):
        for truc in ("dx", "dy"):
            for i, v in enumerate(nghiem[phan][truc]):
                if isinstance(v, bool) or not isinstance(v, (int, np.integer)):
                    if isinstance(v, float) and float(v).is_integer():
                        pass
                    else:
                        loi.append(f"{phan}[{i}].{truc} = {v!r} khong nguyen")
                        continue
                if not -128 <= int(v) <= 127:
                    loi.append(f"{phan}[{i}].{truc} = {v} ngoai signed byte")
    if loi:
        return loi
    dy_l_mong = dy_chan_co_dinh(components, charinfo, ts)
    for l, v in enumerate(nghiem["leg"]["dy"]):
        if int(v) != dy_l_mong[l]:
            loi.append(f"leg[{l}].dy = {v}, mong doi {dy_l_mong[l]} (ban chan cy{ts.ban_chan_x1:+d} tai cf neo)")
    for slot, gia_tri in ts.ghim_dy_dau.items():
        if int(nghiem["head"]["dy"][int(slot)]) != int(gia_tri):
            loi.append(f"head[{slot}].dy = {nghiem['head']['dy'][int(slot)]}, ghim = {gia_tri}")
    for hang in bang_cf(components, charinfo, nghiem, ts):
        cf = hang["cf"]
        if not ts.khe_dau_than[0] <= hang["khe_dau_than"] <= ts.khe_dau_than[1]:
            loi.append(f"cf {cf}: khe_dau_than {hang['khe_dau_than']} ngoai {list(ts.khe_dau_than)}")
        if not ts.khe_than_chan[0] <= hang["khe_than_chan"] <= ts.khe_than_chan[1]:
            loi.append(f"cf {cf}: khe_than_chan {hang['khe_than_chan']} ngoai {list(ts.khe_than_chan)}")
        if not ts.e[0] - 1e-6 <= hang["e"] <= ts.e[1] + 1e-6:
            loi.append(f"cf {cf}: e {hang['e']} ngoai {list(ts.e)}")
        if abs(hang["d"]) > ts.d_max + 1e-6:
            loi.append(f"cf {cf}: d {hang['d']} vuot d_max {ts.d_max}")
    return loi
