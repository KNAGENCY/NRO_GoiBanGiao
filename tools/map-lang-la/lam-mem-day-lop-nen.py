# -*- coding: utf-8 -*-
"""Lam mem DAY lop b01 (lang xa) cho khop dinh lop b00 ngay duoi no.

Cac lop nen la PNG DAC (khong co kenh trong suot), nen mep lop nay xen phang lop
kia -> hien ra mot vach ngang chay suot man hinh. Da xu ly xong cho noi
nui/troi (b02 <-> b01); day la cho noi con lai: b01 <-> b00.
"""
from PIL import Image
import numpy as np
import os, shutil

for muc in ["x1", "x2", "x3", "x4"]:
    d = f"D:/NRO_Project/Assets/Resources/res/{muc}/bg/"
    p01, p00 = d + "b01.png", d + "b00.png"
    if not (os.path.exists(p01) and os.path.exists(p00)):
        print(f"[{muc}] thieu anh, bo qua"); continue
    goc = p01 + ".bak-day"
    if not os.path.exists(goc):
        shutil.copy2(p01, goc)
    dinh00 = np.array(Image.open(p00).convert("RGB")).astype(float)[:6].reshape(-1, 3).mean(axis=0)
    a = np.array(Image.open(goc).convert("RGB")).astype(float)
    h, w, _ = a.shape
    n = max(4, int(h * 0.18))
    truoc = a[-1].mean(axis=0).copy()
    for i in range(n):
        t = (i + 1) / n
        t = t * t * (3 - 2 * t)
        y = h - n + i
        a[y] = a[y] * (1 - t) + dinh00 * t
    Image.fromarray(np.clip(a, 0, 255).astype(np.uint8)).save(p01)
    shutil.copy2(p01, d + "b11.png")
    print(f"[{muc}] b01 day ({truoc[0]:.0f},{truoc[1]:.0f},{truoc[2]:.0f}) -> "
          f"({a[-1].mean(axis=0)[0]:.0f},{a[-1].mean(axis=0)[1]:.0f},{a[-1].mean(axis=0)[2]:.0f})  "
          f"khop dinh b00 ({dinh00[0]:.0f},{dinh00[1]:.0f},{dinh00[2]:.0f}), fade {n} hang")
