# -*- coding: utf-8 -*-
"""Dung tuong VO HINH o hai cot bien map.

Cot 0 va cot 51 cua map 1 dang RONG hoan toan ca 24 hang. Hau qua:
  - nguoi choi di toi mep thi roi xuong ho o cot bien
  - camera dung truoc mep nen nhan vat di qua do la BIEN MAT khoi man hinh

Moi map goc cua Ngoc Rong deu dat tile 39 lam tuong o hai cot nay. Tile 39 da
duoc khai bao DAC ca bon huong trong bang goc (0x2 san, 0x4 chan phai,
0x8 chan trai, 0x2000 tran) nen khong phai sua tile_set_info.

Chi can lam anh cua no TRONG SUOT — tuong chan duoc nguoi ma khong ve ra gi.
Tile 39 hien khong duoc dung o o nao khac trong map nen doi anh la an toan.
"""
from PIL import Image
import numpy as np
import os
import shutil

MAP = "D:/NRO_Project/Assets/Resources/res/mymap/1.bytes"
TILE_TUONG = 39

# 1) anh tile 39 -> trong suot
for muc, canh in (("x1", 24), ("x2", 48), ("x3", 72), ("x4", 96)):
    p = f"D:/NRO_Project/Assets/Resources/res/{muc}/t/1${TILE_TUONG}.png"
    if os.path.exists(p) and not os.path.exists(p + ".bak-tuong"):
        shutil.copy2(p, p + ".bak-tuong")
    Image.fromarray(np.zeros((canh, canh, 4), dtype=np.uint8), "RGBA").save(p)
print(f"tile {TILE_TUONG} -> anh trong suot o 4 muc zoom (tuong vo hinh)")

# 2) dat tuong o hai cot bien
d = bytearray(open(MAP, "rb").read())
w, h = d[0], d[1]
n = 0
for r in range(h):
    for c in (0, w - 1):
        if d[2 + r * w + c] == 0:
            d[2 + r * w + c] = TILE_TUONG
            n += 1
open(MAP, "wb").write(bytes(d))
print(f"dat {n} o tuong vao cot 0 va cot {w-1}")

kt = open(MAP, "rb").read()
o = list(kt[2:2 + w * h])
print(f"kiem: cot 0 = {sorted(set(o[r*w] for r in range(h)))}, "
      f"cot {w-1} = {sorted(set(o[r*w + w-1] for r in range(h)))}")
