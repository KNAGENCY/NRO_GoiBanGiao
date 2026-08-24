# -*- coding: utf-8 -*-
"""Vet cac hang pixel hong o DINH anh nen.

b02 va b03 co 4 hang dau lech han so voi phan con lai (56, 71, 93, 97 roi on
dinh o ~88) — trong game hien thanh mot vach toi manh chay ngang sat dinh man.
Gan chac la du am cua lan cat/ghep anh truoc day.

Cach sua: thay cac hang hong bang hang SACH dau tien.
"""
from PIL import Image
import numpy as np
import os, shutil

NGUONG = 6.0        # lech qua 6 don vi so voi hang chuan = hong

for muc in ["x1", "x2", "x3", "x4"]:
    d = f"D:/NRO_Project/Assets/Resources/res/{muc}/bg/"
    for ten in ("b02", "b03"):
        p = d + ten + ".png"
        if not os.path.exists(p):
            continue
        goc = p + ".bak-dinh"
        if not os.path.exists(goc):
            shutil.copy2(p, goc)
        a = np.array(Image.open(goc).convert("RGB")).astype(float)
        h = a.shape[0]
        chuan = a[6:12].mean()                     # muc sang "that" cua vung dinh
        hong = 0
        while hong < h // 8 and abs(a[hong].mean() - chuan) > NGUONG:
            hong += 1
        if hong == 0:
            print(f"[{muc}] {ten}: dinh sach, khong phai sua")
            continue
        truoc = [round(a[y].mean(), 1) for y in range(hong)]
        a[:hong] = a[hong]
        Image.fromarray(np.clip(a, 0, 255).astype(np.uint8)).save(p)
        # ban sao ma game that dung (typeBg = 1 -> b1x)
        ban_sao = d + ("b1" + ten[-1] + ".png")
        if os.path.exists(ban_sao):
            shutil.copy2(p, ban_sao)
        print(f"[{muc}] {ten}: vet {hong} hang hong {truoc} -> {a[0].mean():.1f} (chuan {chuan:.1f})")
