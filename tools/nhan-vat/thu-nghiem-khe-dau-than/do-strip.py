# -*- coding: utf-8 -*-
"""Đo 5 dải nguồn nhân vật: mỗi khung cao bao nhiêu, đầu to bao nhiêu → tách
"vẽ nhỏ hơn" (scale) khỏi "tư thế thấp hơn" (pose)."""
import json, sys
from pathlib import Path
from PIL import Image
sys.path.insert(0, r"D:\NRO_GoiBanGiao")
from tools.artqa.task6b_closure import _review_by_name, _frame_rect, PLAYER_ACTIONS, _load_reviewed_alpha

ART = Path(r"D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1")
review = json.loads((Path(r"D:\NRO_GoiBanGiao\docs\my-thuat\review\vertical-slice-v1.json")).read_text(encoding="utf-8"))
rec = _review_by_name(review)
cache = {}
print("PLAYER_ACTIONS =", PLAYER_ACTIONS)
print("action   fr | frame x,y,w,h            | H    | dau: rong  cao (vung 0-42%) | rong hang rong nhat 0-25% | H/H_idle0  dau/dau_idle0")
H0 = None; W0 = None
rows = []
for action, src, n in PLAYER_ACTIONS:
    alpha = _load_reviewed_alpha(ART, rec, src, cache)
    for i in range(n):
        x, y, w, h = _frame_rect(rec[src], i)
        fr = alpha.crop((x, y, x + w, y + h))
        a = fr.getchannel("A")
        head = a.crop((0, 0, w, round(h * 0.42)))
        hb = head.getbbox()
        top25 = a.crop((0, 0, w, round(h * 0.25)))
        # rộng hàng rộng nhất trong 25% trên (xấp xỉ bề ngang hộp sọ, ít phụ thuộc tư thế)
        px = top25.load(); best = 0
        for yy in range(top25.height):
            xs = [xx for xx in range(w) if px[xx, yy] > 64]
            if xs:
                best = max(best, xs[-1] - xs[0] + 1)
        if H0 is None:
            H0, W0 = h, best
        rows.append((action, i, (x, y, w, h), h, hb[2] - hb[0] if hb else 0, hb[3] - hb[1] if hb else 0, best, h / H0, best / W0))
for r in rows:
    print("%-8s %2d | %-24s | %4d | %4d %4d | %4d | %.3f %.3f" % (r[0], r[1], r[2], r[3], r[4], r[5], r[6], r[7], r[8]))
