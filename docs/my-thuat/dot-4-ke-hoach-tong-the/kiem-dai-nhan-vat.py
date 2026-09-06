# -*- coding: utf-8 -*-
"""kiem-dai-nhan-vat.py — kiểm NHANH hình học dải chroma nhân vật đợt 4 (brief 01-nhan-vat.md §7 bước 1).

Không thay `python -m tools.artqa.artqa review` (cổng thật của dự án); đây là bước người vẽ tự soi
trước khi giao, chỉ cần Pillow + numpy (scipy có thì kiểm thêm đảo tím trong hình).

Dùng:
    python kiem-dai-nhan-vat.py <thư-mục-dải> [--o 724] [--le 40] [--chan 4] [--cao 600 700]

Kiểm cho MỌI *.png trong thư mục:
  (1) tên đúng mẫu  <...>_<n>f_chroma.png  hoặc  <...>_portrait_chroma.png
  (2) cỡ ảnh = n×ô × ô  (portrait: vuông, cạnh ≥ 1000)
  (3) viền ảnh ≥ 85 % tím RGB(246,4,242) và lệch màu viền ≤ 12  (đúng ngưỡng artqa.py:23-24)
  (4) mỗi ô có hình (không ô rỗng)
  (5) hình cách mép trái/phải ô ≥ --le px (không cắn ô kề)
  (6) bàn chân chạm đáy ô: mép dưới hộp ≥ ô − --chan px   (không áp cho portrait)
  (7) cao hình trong [--cao]  (ô 724: 600–700 ≈ tóc→chân 640–670 ± tư thế; không áp cho portrait)
  (8) không có mảng tím tách rời bên trong hộp hình (> 0,04 % diện tích ô) — cần scipy, không có thì bỏ qua
KHÔNG đo cỡ đầu — đối chiếu bằng mắt với anh-tham-chieu/THUOC-NHAN-VAT.png.
Thoát mã 1 nếu có file FAIL. In bảng để dán nguyên vào báo cáo.

PHẠM VI: chỉ cho thư mục dải NHÂN VẬT / boss / NPC đứng thẳng, ô 724 (brief 01, 02, 03). Quái (rộng hơn cao),
hiệu ứng hu_, sheet icon chiêu, HUD master KHÔNG kiểm bằng script này — dùng `python -m tools.artqa.artqa review`.
Chạy thử 06/09/2026 trên anh-nguon/vertical-slice-v1: 20/20 FAIL, đúng các lỗi đã biết của lô v1 (nền gradient
lệch 28–33, chân hở đáy 9–49 px ở idle/cast, ô 543 ở run/melee, khung nằm ở hit) → script bắt được lỗi.
"""
from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path

import numpy as np
from PIL import Image

CHROMA = (246, 4, 242)


def la_tim(rgb: np.ndarray) -> np.ndarray:
    """Cùng luật với artqa.py:60-62: R−G ≥ 110 và B−G ≥ 110 → coi là nền."""
    r = rgb[..., 0].astype(np.int16)
    g = rgb[..., 1].astype(np.int16)
    b = rgb[..., 2].astype(np.int16)
    return (r - g >= 110) & (b - g >= 110)


def so_o(ten: str) -> int | None:
    if ten.endswith("_portrait_chroma.png"):
        return 1
    m = re.search(r"_(\d+)f_chroma\.png$", ten)
    return int(m.group(1)) if m else None


def kiem_file(p: Path, o: int, le: int, chan: int, cao: tuple[int, int]) -> tuple[bool, list[str]]:
    loi: list[str] = []
    n = so_o(p.name)
    if n is None:
        return False, [f"tên sai mẫu (cần _<n>f_chroma.png hoặc _portrait_chroma.png)"]
    portrait = p.name.endswith("_portrait_chroma.png")
    im = Image.open(p).convert("RGB")
    W, H = im.size
    rgb = np.asarray(im)
    tim = la_tim(rgb)

    # (2) cỡ
    if portrait:
        if W != H or W < 1000:
            loi.append(f"cỡ {W}×{H}: chân dung phải vuông và ≥ 1000 (mẫu 1254×1254)")
        o_thuc = W
    else:
        if (W, H) != (n * o, o):
            loi.append(f"cỡ {W}×{H}: phải đúng {n * o}×{o} ({n} ô × {o})")
        o_thuc = W // n if n else W

    # (3) viền
    vien = np.concatenate([rgb[0, :, :], rgb[-1, :, :], rgb[:, 0, :], rgb[:, -1, :]])
    vien_tim = la_tim(vien)
    ti_le = float(vien_tim.mean()) if len(vien) else 0.0
    if ti_le < 0.85:
        loi.append(f"viền chỉ {ti_le:.0%} tím (cần ≥ 85 %) — hình chạm mép ảnh hoặc nền không tím")
    if vien_tim.any():
        vt = vien[vien_tim].astype(np.int16)
        lech = int((vt.max(axis=0) - vt.min(axis=0)).max())
        if lech > 12:
            loi.append(f"nền viền lệch màu {lech} (cần ≤ 12) — có gradient/vignette, phải là tím phẳng")
        tb = vt.mean(axis=0)
        if np.abs(tb - np.array(CHROMA)).max() > 6:
            loi.append(f"màu nền trung bình {tuple(int(x) for x in tb)} khác RGB{CHROMA}")

    # (4)-(8) từng ô
    try:
        from scipy import ndimage as ndi  # type: ignore
    except Exception:  # pragma: no cover
        ndi = None
    for i in range(n):
        x0, x1 = i * o_thuc, (i + 1) * o_thuc if i < n - 1 else W
        cell = ~tim[:, x0:x1]
        ys, xs = np.nonzero(cell)
        if len(ys) == 0:
            loi.append(f"ô {i}: rỗng")
            continue
        top, bot, lft, rgt = int(ys.min()), int(ys.max()), int(xs.min()), int(xs.max())
        if lft < le or (x1 - x0 - 1 - rgt) < le:
            loi.append(f"ô {i}: hình cách mép ô trái {lft} / phải {x1 - x0 - 1 - rgt} px (cần ≥ {le}) — nguy cơ cắn ô kề")
        if not portrait:
            ho = H - 1 - bot
            if ho > chan:
                loi.append(f"ô {i}: bàn chân hở đáy ô {ho} px (cần ≤ {chan}) — trong game nhân vật sẽ bay")
            cao_hinh = bot - top + 1
            if not (cao[0] <= cao_hinh <= cao[1]):
                loi.append(f"ô {i}: cao hình {cao_hinh} px (mong {cao[0]}–{cao[1]}) — nghi khác cỡ đầu / khom / nằm")
            rong = rgt - lft + 1
            if rong > cao_hinh:
                loi.append(f"ô {i}: hình rộng {rong} > cao {cao_hinh} — nghi khung NẰM (cấm)")
        if ndi is not None:
            # đảo tím bên trong hộp hình, không nối với nền ngoài
            box = tim[top:bot + 1, x0 + lft:x0 + rgt + 1]
            lab, k = ndi.label(box)
            if k:
                bien = set(np.unique(np.concatenate([lab[0], lab[-1], lab[:, 0], lab[:, -1]]))) - {0}
                dien = np.bincount(lab.ravel())
                dao = sum(int(dien[j]) for j in range(1, k + 1) if j not in bien)
                if dao > 0.0004 * o_thuc * o_thuc:
                    loi.append(f"ô {i}: {dao} px tím nằm trong hình (> 0,04 % ô) — máy sẽ khoét lỗ, đổi màu chỗ đó")
    return not loi, loi


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("thu_muc")
    ap.add_argument("--o", type=int, default=724, help="cạnh ô (mặc định 724)")
    ap.add_argument("--le", type=int, default=40, help="khoảng cách tối thiểu tới mép trái/phải ô")
    ap.add_argument("--chan", type=int, default=4, help="bàn chân hở đáy tối đa (px)")
    ap.add_argument("--cao", type=int, nargs=2, default=(600, 700), help="cao hình cho phép (px)")
    a = ap.parse_args()
    files = sorted(Path(a.thu_muc).glob("*.png"))
    if not files:
        print(f"Không có *.png trong {a.thu_muc}")
        return 1
    hong = 0
    print(f"KIỂM DẢI CHROMA — {a.thu_muc} — ô {a.o}, lề ≥ {a.le}, chân ≤ {a.chan}, cao {a.cao[0]}–{a.cao[1]}")
    for p in files:
        ok, loi = kiem_file(p, a.o, a.le, a.chan, tuple(a.cao))
        im = Image.open(p)
        print(f"{'PASS' if ok else 'FAIL'}  {p.name}  {im.size[0]}×{im.size[1]}")
        for l in loi:
            print(f"      - {l}")
        hong += 0 if ok else 1
    print(f"Tổng: {len(files)} file, {len(files) - hong} PASS, {hong} FAIL")
    return 1 if hong else 0


if __name__ == "__main__":
    sys.exit(main())
