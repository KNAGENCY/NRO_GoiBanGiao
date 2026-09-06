# -*- coding: utf-8 -*-
"""nhan_vat_moc.py — mốc giải phẫu đầu nhân vật người chơi (V1, 05/09/2026).

VÌ SAO CÓ MODULE NÀY
--------------------
"Đầu tách khỏi thân" trong game hoá ra là HAI KHUÔN MẶT chồng nhau: pipeline cắt mảnh
THÂN ở 28 % chiều cao khung nguồn, nhưng 5 dải nguồn vẽ nhân vật khác tỉ lệ (dải chạy
nhỏ hơn idle ~12 %, dải đánh ~17 %) và khom hơn, nên ở khung chạy/đánh lát cắt ấy rơi
lên trán → mảnh thân chứa cả mắt–mũi–miệng và client vẽ nó ĐÈ lên đầu idle (thứ tự vẽ
đầu → chân → thân, ``Char.cs`` 6081–6084). Muốn sửa bằng dữ liệu phải đo được, cho
TỪNG khung nguồn:

* **cằm** (mép dưới khối da mặt) — để cắt thân từ cằm xuống, không còn mặt;
* **cỡ đầu** — để nhân tỉ lệ từng khung sao cho đầu idle (3 khung dùng chung cho cả
  33 cf trong bộ xương cứng ``CharInfo``) khớp cỡ với thân của khung đó.

Cỡ đầu đo bằng **median 6 thước cứng** của đầu (rộng/cao mặt, cao mắt, mắt→cằm,
đỉnh tóc→cằm, dày băng đô) — dùng median để một thước bị tay che hay tóc bay không kéo
lệch cả khung. Bộ thước này là đúng bộ đã sinh ra nghiệm chốt trong hồ sơ
``docs/my-thuat/review/khe-dau-than-2026-09-05/`` (``geometry-V1.json``), chép sang
từ ``tools/nhan-vat/thu-nghiem-khe-dau-than/do-moc-giai-phau.py`` + ``sinh-manh-v1.py``.

Toàn bộ số đo là **đơn vị pixel nguồn**: ``*_y`` tuyệt đối theo chiều cao dải,
``face_cx`` tương đối so với mép trái Ô khung (``source_edges[i]``).

Dùng:
    python -m tools.artqa.nhan_vat_moc --out docs/my-thuat/review/nhan-vat-moc-giai-phau.json
"""

from __future__ import annotations

import argparse
import hashlib
import json
import math
from pathlib import Path
from statistics import median
from typing import Any

import numpy as np
from PIL import Image
from scipy import ndimage as ndi

CHROMA = (246, 4, 242)
NGUON_NGUOI_CHOI = (
    "player_idle_3f_chroma.png",
    "player_run_4f_chroma.png",
    "player_melee_4f_chroma.png",
    "player_cast_3f_chroma.png",
    "player_hit_death_getup_3f_chroma.png",
)
NGUON_IDLE = "player_idle_3f_chroma.png"
# 6 thước cứng của đầu (không đổi theo tư thế tay/chân).
THUOC_CUNG = ("face_w", "face_h", "eye_h", "eye_to_chin", "head_h_hair_chin", "band_thick")
# Thước bị hỏng ở khung cụ thể (đo bằng mắt trên dbg_*.png của prototype):
#   melee#1: tay đấm che cằm → cao mặt và mắt→cằm ngắn giả (68 so 96 idle).
#   hit#0 : ngã ngửa, mặt nghiêng → mắt/băng đô/tóc đo sai.
THUOC_LOAI = {
    ("player_melee_4f_chroma.png", 1): frozenset({"face_h", "eye_to_chin"}),
    ("player_hit_death_getup_3f_chroma.png", 0): frozenset(
        {"eye_h", "eye_to_chin", "band_thick", "head_h_hair_chin"}
    ),
}
# Khung mà cằm bị che → ước lượng cằm = đỉnh mặt + cao mặt idle × hệ số đầu.
CAM_BI_CHE = frozenset({("player_melee_4f_chroma.png", 1)})


class MocError(ValueError):
    """Đo mốc giải phẫu thất bại (nguồn lệch SHA, thiếu khung…)."""


# --------------------------------------------------------------------------- màu
def _rgb_sang_hsv(a: np.ndarray):
    a = a.astype(np.float32) / 255.0
    r, g, b = a[..., 0], a[..., 1], a[..., 2]
    mx = a.max(-1)
    mn = a.min(-1)
    d = mx - mn
    h = np.zeros_like(mx)
    m = d > 1e-6
    dd = np.where(d == 0, 1, d)
    rm = m & (mx == r)
    gm = m & (mx == g) & ~rm
    bm = m & ~rm & ~gm
    h[rm] = (60 * ((g - b) / dd) % 360)[rm]
    h[gm] = (60 * ((b - r) / dd) + 120)[gm]
    h[bm] = (60 * ((r - g) / dd) + 240)[bm]
    s = np.where(mx > 0, d / np.where(mx == 0, 1, mx), 0)
    return h, s, mx


def _khoi_lon_nhat(mask: np.ndarray, dieu_kien=None):
    lab, n = ndi.label(mask)
    if n == 0:
        return None
    sizes = ndi.sum(mask, lab, range(1, n + 1))
    for k in np.argsort(-sizes):
        khoi = lab == (k + 1)
        if dieu_kien is None or dieu_kien(khoi):
            return khoi
    return None


def _hop(m: np.ndarray):
    ys, xs = np.where(m)
    if len(ys) == 0:
        return None
    return int(xs.min()), int(ys.min()), int(xs.max()) + 1, int(ys.max()) + 1


# --------------------------------------------------------------------------- đo
def do_khung(rgb: np.ndarray, tag: str) -> dict[str, Any]:
    """Đo một khung nguồn (mảng H×W×3 uint8 trên nền chroma).

    Trả dict: ``content_*`` (hộp nội dung), và nếu tìm được mặt: ``face_w/face_h/
    face_area/chin_y/face_top/face_cx``, mắt (``eye_h/eye_w/eye_cy/eye_to_chin``),
    ``hair_top/head_h_hair_chin``, ``band_thick/band_w``, ``feet_to_chin``.
    Khung không có khối da (nằm sấp) chỉ có ``content_*``.
    """

    H, W, _ = rgb.shape
    bg = np.array(CHROMA)
    fg = np.abs(rgb.astype(int) - bg).sum(-1) > 90
    fg = ndi.binary_opening(fg, iterations=1)
    cb = _hop(fg)
    r: dict[str, Any] = {"tag": tag}
    if cb is None:
        return r
    r.update(
        content_h=cb[3] - cb[1], content_w=cb[2] - cb[0],
        content_bottom=cb[3], content_top=cb[1],
    )
    h, s, v = _rgb_sang_hsv(rgb)
    skin = fg & (h > 20) & (h < 48) & (s > 0.35) & (s < 0.8) & (v > 0.55)
    skin = ndi.binary_opening(skin, iterations=1)
    navy = fg & (h > 195) & (h < 250) & (s > 0.28) & (v < 0.45)  # băng đô / quần / tay áo
    # Mặt = khối da lớn nhất có tâm nằm trong 55 % trên của hộp nội dung.
    top_lim = cb[1] + 0.55 * (cb[3] - cb[1])

    def o_tren(khoi):
        ys, _ = np.where(khoi)
        return ys.mean() < top_lim

    face = _khoi_lon_nhat(skin, o_tren)
    if face is None:
        return r
    fx0, fy0, fx1, fy1 = _hop(face)
    r.update(
        face_w=fx1 - fx0, face_h=fy1 - fy0, face_area=int(face.sum()),
        chin_y=fy1, face_top=fy0, face_cx=(fx0 + fx1) / 2,
    )
    # Mắt = lỗ lớn nhất trong khối da (mắt / lông mày / miệng).
    holes = ndi.binary_fill_holes(face) & ~face
    eye = _khoi_lon_nhat(holes)
    if eye is not None:
        eb = _hop(eye)
        r.update(eye_h=eb[3] - eb[1], eye_w=eb[2] - eb[0], eye_cy=(eb[1] + eb[3]) / 2, eye_area=int(eye.sum()))
        r["eye_to_chin"] = fy1 - r["eye_cy"]
    # Đỉnh tóc: hàng fg cao nhất trong dải cột của mặt (±50 % bề rộng mặt).
    xa = max(0, int(fx0 - 0.5 * (fx1 - fx0)))
    xb = min(W, int(fx1 + 0.5 * (fx1 - fx0)))
    col = fg[:, xa:xb].any(1)
    hair_top = int(np.argmax(col))
    r["hair_top"] = hair_top
    r["head_h_hair_chin"] = fy1 - hair_top
    # Băng đô: navy trong khoảng hàng [hair_top, face_top + 30 % face_h], cột của mặt;
    # bề dày = median số hàng navy theo cột ở 1/3 giữa.
    band_region = np.zeros_like(navy)
    band_region[hair_top:int(fy0 + 0.3 * (fy1 - fy0)), fx0:fx1] = True
    bandc = _khoi_lon_nhat(navy & band_region)
    if bandc is not None:
        cols = range(int(fx0 + (fx1 - fx0) / 3), int(fx0 + 2 * (fx1 - fx0) / 3))
        th = [int(bandc[:, c].sum()) for c in cols]
        th = [t for t in th if t > 0]
        r["band_thick"] = float(median(th)) if th else None
        bb = _hop(bandc)
        r["band_w"] = bb[2] - bb[0]
    r["feet_to_chin"] = cb[3] - fy1
    return r


def _sha(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def do_dai(art_root: Path, records: dict[str, dict[str, Any]], source_name: str) -> dict[tuple[str, int], dict[str, Any]]:
    """Đo mọi khung của một dải nguồn đã PASS review; khoá = (tên nguồn, chỉ số khung)."""

    record = records.get(source_name)
    if record is None:
        raise MocError(f"ASSET-CLOSURE-015: nguồn chưa PASS: {source_name}")
    source = Path(art_root) / source_name
    if not source.is_file() or _sha(source) != record.get("source_sha256"):
        raise MocError(f"ASSET-CLOSURE-015: SHA source lệch review: {source_name}")
    edges = record.get("checks", {}).get("crop", {}).get("source_edges")
    if not isinstance(edges, list) or len(edges) < 2:
        raise MocError(f"ASSET-CLOSURE-015: crop review thiếu source_edges: {source_name}")
    with Image.open(source) as image:
        im = np.array(image.convert("RGB"))
    result = {}
    for i in range(len(edges) - 1):
        r = do_khung(im[:, edges[i]:edges[i + 1]], f"{source_name}#{i}")
        r["source"] = source_name
        r["i"] = i
        r["edge_x"] = int(edges[i])
        result[(source_name, i)] = r
    return result


def do_nguon(art_root: Path, records: dict[str, dict[str, Any]], sources=NGUON_NGUOI_CHOI) -> dict[tuple[str, int], dict[str, Any]]:
    """Đo mọi khung của các dải người chơi (mặc định 5 dải = 17 khung)."""

    moc: dict[tuple[str, int], dict[str, Any]] = {}
    for name in sources:
        moc.update(do_dai(art_root, records, name))
    return moc


# --------------------------------------------------------------------------- tỉ lệ
def tham_chieu_idle(moc, source: str = NGUON_IDLE, frames=(0, 1, 2)) -> dict[str, float]:
    """Cỡ đầu tham chiếu = trung bình 3 khung idle (6 thước cứng)."""

    ref = {}
    for k in THUOC_CUNG:
        vals = [moc[(source, i)][k] for i in frames if moc[(source, i)].get(k) is not None]
        if not vals:
            raise MocError(f"ASSET-CLOSURE-015: idle thiếu thước {k}")
        ref[k] = float(sum(vals) / len(vals))
    return ref


def he_so_dau(khung: dict[str, Any], tham_chieu: dict[str, float], loai=frozenset()) -> float | None:
    """cỡ đầu khung / cỡ đầu tham chiếu = median 6 thước cứng (bỏ thước bị loại).

    ``None`` khi khung không có mặt (khung nằm sấp).
    """

    if khung.get("face_w") is None:
        return None
    ratios = []
    for k in THUOC_CUNG:
        if k in loai:
            continue
        v = khung.get(k)
        if v is None:
            continue
        ratios.append(v / tham_chieu[k])
    if not ratios:
        return None
    return float(median(ratios))


def he_so_khung(moc, key: tuple[str, int], tham_chieu: dict[str, float]) -> float | None:
    """Hệ số đầu của khung ``key`` với bảng loại thước ``THUOC_LOAI``."""

    return he_so_dau(moc[key], tham_chieu, THUOC_LOAI.get(key, frozenset()))


def ti_le_khung(ti_le_goc: float, he_so: float | None) -> float:
    """s_f = s0 / (cỡ đầu khung / cỡ đầu idle); khung không có mặt giữ s0."""

    if he_so is None or he_so <= 0:
        return float(ti_le_goc)
    return float(ti_le_goc) / float(he_so)


def cam_khung(khung: dict[str, Any], tham_chieu: dict[str, float], he_so: float | None, *, bi_che: bool) -> float | None:
    """Toạ độ y (nguồn) của cằm; khung bị tay che cằm thì ước lượng từ đỉnh mặt."""

    if khung.get("chin_y") is None:
        return None
    if bi_che:
        return float(khung["face_top"]) + tham_chieu["face_h"] * float(he_so if he_so else 1.0)
    return float(khung["chin_y"])


def hinh_hoc_khung(moc, tham_chieu: dict[str, float], ti_le_goc: float) -> dict[tuple[str, int], dict[str, Any]]:
    """Từng khung: hệ số đầu, tỉ lệ s_f, cằm, tâm mặt, đỉnh tóc (đều đơn vị nguồn)."""

    result = {}
    for key, khung in moc.items():
        he_so = he_so_khung(moc, key, tham_chieu)
        result[key] = {
            "he_so_dau": he_so,
            "scale": ti_le_khung(ti_le_goc, he_so),
            "chin_y": cam_khung(khung, tham_chieu, he_so, bi_che=key in CAM_BI_CHE),
            "cam_uoc_luong": key in CAM_BI_CHE,
            "face_cx": None if khung.get("face_cx") is None else float(khung["face_cx"]),
            "hair_top": None if khung.get("content_top") is None else int(khung["content_top"]),
            "thuoc_loai": sorted(THUOC_LOAI.get(key, ())),
        }
    return result


# --------------------------------------------------------------------------- JSON
def _so(v):
    if isinstance(v, (np.integer,)):
        return int(v)
    if isinstance(v, (np.floating,)):
        return float(v)
    return v


def xuat_json(moc, tham_chieu, hinh_hoc, ti_le_goc: float, records: dict[str, dict[str, Any]]) -> dict[str, Any]:
    """Gói kết quả đo thành JSON tái lập được (kèm SHA nguồn)."""

    sources = sorted({k[0] for k in moc})
    return {
        "generator": "tools/artqa/nhan_vat_moc.py",
        "ghi_chu": (
            "V1 05/09/2026 — mốc giải phẫu đầu nhân vật người chơi đo tự động từ ảnh nguồn "
            "chroma (tách màu da/navy). Đơn vị px nguồn; *_y tuyệt đối theo dải, face_cx so "
            "với mép trái ô khung. Cỡ đầu = median 6 thước cứng; tỉ lệ từng khung s_f = "
            "s0 / (cỡ đầu khung / cỡ đầu idle)."
        ),
        "ti_le_goc_x4": ti_le_goc,
        "thuoc_cung": list(THUOC_CUNG),
        "tham_chieu_idle": tham_chieu,
        "nguon": {name: records[name]["source_sha256"] for name in sources},
        "khung": [
            {
                **{k: _so(v) for k, v in moc[key].items()},
                **{f"hh_{k}": _so(v) for k, v in hinh_hoc[key].items()},
            }
            for key in sorted(moc, key=lambda k: (sources.index(k[0]), k[1]))
        ],
    }


def main() -> int:
    from tools.artqa.task6b_closure import (
        PLAYER_REFERENCE_FRAME,
        PLAYER_REFERENCE_SOURCE,
        _frame_rect,
        _review_by_name,
        player_uniform_scale,
    )

    goc = Path(__file__).resolve().parents[2]
    parser = argparse.ArgumentParser()
    parser.add_argument("--art-root", type=Path, default=goc / "docs/my-thuat/anh-nguon/vertical-slice-v1")
    parser.add_argument("--review", type=Path, default=goc / "docs/my-thuat/review/vertical-slice-v1.json")
    parser.add_argument("--out", type=Path, default=goc / "docs/my-thuat/review/nhan-vat-moc-giai-phau.json")
    args = parser.parse_args()
    records = _review_by_name(json.loads(args.review.read_text(encoding="utf-8")))
    moc = do_nguon(args.art_root, records)
    tham_chieu = tham_chieu_idle(moc)
    s0 = player_uniform_scale(_frame_rect(records[PLAYER_REFERENCE_SOURCE], PLAYER_REFERENCE_FRAME)[2])
    hinh_hoc = hinh_hoc_khung(moc, tham_chieu, s0)
    payload = xuat_json(moc, tham_chieu, hinh_hoc, s0, records)
    args.out.parent.mkdir(parents=True, exist_ok=True)
    args.out.write_text(json.dumps(payload, ensure_ascii=False, indent=1) + "\n", encoding="utf-8")
    print("khung                                 he_so  s_f     cam_y")
    for key in sorted(hinh_hoc):
        hh = hinh_hoc[key]
        print("%-36s %s  %.4f  %s" % (
            f"{key[0]}#{key[1]}",
            "%.4f" % hh["he_so_dau"] if hh["he_so_dau"] is not None else "  -   ",
            hh["scale"],
            "%.1f%s" % (hh["chin_y"], "*" if hh["cam_uoc_luong"] else "") if hh["chin_y"] is not None else "-",
        ))
    print("da ghi", args.out)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
