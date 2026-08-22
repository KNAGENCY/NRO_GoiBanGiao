# -*- coding: utf-8 -*-
"""Lam mem CHO GIAP giua cac lop nen, cho het duong ke ngang chay suot man.

VI SAO CO DUONG KE. Lop nen la PNG DAC ve voi chieu cao co dinh. Duoi day moi
lop, engine dien not man hinh bang MOT MAU DON (mau hang day cua anh). Neu day
anh con noi dung — may, mai ngoi, tan cay — thi tren duong day la noi dung, duoi
la mang mau phang: ra mot duong ke ngang, dung yen khi cuon camera.

MOT MAU NOI DUY NHAT cho moi cho giap. Da thu cach "day lop tren khop dinh lop
duoi" theo day chuyen va HONG: mau bi keo dan qua tung lop, cuoi cung chan nui
Hokage bi nhuom mau lang xa — buoc nhay o cho giap tang tu 13,5 len 27,4.
Ca bon lop deu la vat the tren nen TROI, nen tat ca cho giap dung chung mau troi.

Chay lai duoc nhieu lan: luon doc tu ban .bak-noi.
"""
from PIL import Image
import numpy as np
import os
import shutil

MAU_NOI = np.array([10.0, 86.0, 154.0])   # mau troi, dung cho MOI cho giap
FADE_DAY, THUAN_DAY, FADE_TREN = 0.30, 0.15, 0.22


def mo(p):
    goc = p + ".bak-noi"
    if not os.path.exists(goc):
        shutil.copy2(p, goc)
    return np.array(Image.open(goc).convert("RGB")).astype(float)


def luu(a, p, d, ten):
    Image.fromarray(np.clip(a, 0, 255).astype(np.uint8)).save(p)
    shutil.copy2(p, d + "b1" + ten[-1] + ".png")   # ban client THAT nap (typeBg=1)


for muc in ["x1", "x2", "x3", "x4"]:
    d = f"D:/NRO_Project/Assets/Resources/res/{muc}/bg/"
    bao = []
    for ten in ["b03", "b02", "b01", "b00"]:
        p = d + ten + ".png"
        if not os.path.exists(p):
            continue
        a = mo(p)
        h = a.shape[0]
        gon = a[-1].std(axis=0).mean()
        # day nhoa ve mau noi (tru lop cuoi b00 — duoi no khong con lop nao)
        if ten != "b00":
            bd, kt = int(h * (1 - FADE_DAY)), int(h * (1 - THUAN_DAY))
            for y in range(bd, h):
                t = 1.0 if y >= kt else (lambda u: u * u * (3 - 2 * u))((y - bd) / max(1, kt - bd))
                a[y] = a[y] * (1 - t) + MAU_NOI * t
        # mep tren bat dau tu mau noi (tru lop tren cung b03 — tren no la khoang khong)
        if ten != "b03":
            n = max(4, int(h * FADE_TREN))
            for y in range(n):
                t = (lambda u: u * u * (3 - 2 * u))(y / n)
                a[y] = MAU_NOI * (1 - t) + a[y] * t
        luu(a, p, d, ten)
        bao.append(f"{ten} {gon:4.1f}->{a[-1].std(axis=0).mean():3.1f}")
    print(f"[{muc}] " + " | ".join(bao))
