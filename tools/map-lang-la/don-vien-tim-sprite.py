# -*- coding: utf-8 -*-
"""Don pixel magenta sot lai o mep sprite nhan vat.

Sprite duoc tach khoi nen magenta, con sot vai pixel mau (R cao, B cao, G thap)
bam o vien. Phong to len la thay ngay mot duong vien tim quanh nguoi.

Cach sua: pixel nao la magenta thi lay mau trung binh cua cac hang xom KHONG
magenta ma dap vao (giu nguyen alpha). Neu quanh no toan magenta thi cho trong
suot luon.

LUU Y: sua o ban CHEP trong client. Anh nguon ArtRoot
(docs/my-thuat/anh-nguon/.../production-closure-v11/) van con loi — can sua o do
roi build lai, neu khong lan cam sprite sau se quay nguoc.
"""
from PIL import Image
import numpy as np
import glob
import os
import shutil

MUC = ["x1", "x2", "x3", "x4"]


def la_magenta(a):
    r, g, b, al = a[:, :, 0], a[:, :, 1], a[:, :, 2], a[:, :, 3]
    return (al > 20) & (r > 70) & (b > 70) & (g < r - 22) & (g < b - 22)


tong_anh, tong_px = 0, 0
for muc in MUC:
    thu_muc = f"D:/NRO_Project/Assets/Resources/res/{muc}/smallimage"
    n_anh, n_px = 0, 0
    for p in glob.glob(thu_muc + "/Small*.png"):
        goc = p + ".bak-tim"
        a = np.array(Image.open(goc if os.path.exists(goc) else p).convert("RGBA")).astype(int)
        xau = la_magenta(a)
        if not xau.any():
            continue
        if not os.path.exists(p + ".bak-tim"):
            shutil.copy2(p, p + ".bak-tim")
        h, w, _ = a.shape
        sach = (a[:, :, 3] > 20) & ~xau
        ys, xs = np.where(xau)
        for y, x in zip(ys, xs):
            tong = np.zeros(3, dtype=float)
            dem = 0
            for dy in (-1, 0, 1):
                for dx in (-1, 0, 1):
                    ny, nx = y + dy, x + dx
                    if 0 <= ny < h and 0 <= nx < w and sach[ny, nx]:
                        tong += a[ny, nx, :3]
                        dem += 1
            if dem:
                a[y, x, :3] = (tong / dem).round()
            else:
                a[y, x, 3] = 0          # quanh no toan magenta -> cho trong suot
        Image.fromarray(a.astype(np.uint8), "RGBA").save(p)
        n_anh += 1
        n_px += int(xau.sum())
    print(f"  [{muc}] {n_anh} anh co vien tim, don {n_px} pixel")
    tong_anh += n_anh
    tong_px += n_px

print(f"\ntong: {tong_anh} anh, {tong_px} pixel magenta da don")
