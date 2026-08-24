# -*- coding: utf-8 -*-
"""Xoa vien trang/den nam TRONG anh tile nen (60-69) o ca 4 muc zoom.

Cach lam: KHONG cat anh (se lech ti le), ma keo hang/cot SACH gan nhat de len mep.
Nho vay tile ghep lien nhau khong con khe trang/den.
"""
from PIL import Image
import numpy as np
import os, shutil

MUC = ["x1", "x2", "x3", "x4"]
TILE = list(range(60, 70))


def ban(px, al):
    """Mot hang/cot co 'ban' khong: qua trang (>215) hoac qua den (<42)."""
    dac = al > 40
    if dac.sum() == 0:
        return False
    s = px[dac].mean()
    return s > 212 or s < 30


def sua(duong_dan):
    im = Image.open(duong_dan).convert("RGBA")
    a = np.array(im).astype(int)
    h, w, _ = a.shape
    gioi_han = max(1, min(h, w) // 8)   # vien khong bao gio day qua 1/8 anh
    da_sua = []

    # tren
    k = 0
    while k < gioi_han and ban(a[k, :, :3], a[k, :, 3]): k += 1
    if k: a[:k] = a[k]; da_sua.append(f"tren{k}")
    # duoi
    k = 0
    while k < gioi_han and ban(a[h-1-k, :, :3], a[h-1-k, :, 3]): k += 1
    if k: a[h-k:] = a[h-1-k]; da_sua.append(f"duoi{k}")
    # trai
    k = 0
    while k < gioi_han and ban(a[:, k, :3], a[:, k, 3]): k += 1
    if k: a[:, :k] = a[:, k:k+1]; da_sua.append(f"trai{k}")
    # phai
    k = 0
    while k < gioi_han and ban(a[:, w-1-k, :3], a[:, w-1-k, 3]): k += 1
    if k: a[:, w-k:] = a[:, w-1-k:w-k]; da_sua.append(f"phai{k}")

    if da_sua:
        if not os.path.exists(duong_dan + ".bak-vien"):
            shutil.copy2(duong_dan, duong_dan + ".bak-vien")
        Image.fromarray(a.astype(np.uint8), "RGBA").save(duong_dan)
    return da_sua


tong = 0
for m in MUC:
    thu_muc = f"D:/NRO_Project/Assets/Resources/res/{m}/t/"
    dong = []
    for i in TILE:
        p = thu_muc + "1$" + str(i) + ".png"
        if not os.path.exists(p):
            dong.append(f"{i}:THIEU"); continue
        kq = sua(p)
        if kq:
            dong.append(f"{i}:{'+'.join(kq)}"); tong += 1
    print(f"[{m}] " + ("  ".join(dong) if dong else "khong co gi phai sua"))
print(f"\nTong {tong} anh da xoa vien.")
