# -*- coding: utf-8 -*-
"""Lam mem DAY cua moi lop nen, cho khop dinh lop ke duoi.

VI SAO. Cac lop nen la PNG DAC (khong co kenh trong suot) va duoc ve voi chieu
cao co dinh. Duoi day moi lop, engine dien not man hinh bang MOT MAU DON
(fillRect voi mau hang day cua anh). Neu day anh con noi dung — may, mai ngoi,
tan cay — thi tren duong day la noi dung, duoi la mang mau phang: hien ra thanh
mot duong ke ngang chay suot man hinh, dung yen khi cuon camera.

Da gap ba duong nhu vay:
    y=344  day lop nui Hokage (b02)      -> chua bang cach fade chan nui
    y=141  day lop may       (b03)
    y=561  day lop lang xa   (b01)

CACH CHUA. Cho phan duoi moi lop nhoa dan ve dung mau DINH cua lop ke duoi, va
20% cuoi cung la mau do THUAN (het van). Nho vay cho giap giua "anh da phang"
va "mang mau don" khong con ranh gioi.

Chay lai duoc nhieu lan: luon doc tu ban .bak-day.
"""
from PIL import Image
import numpy as np
import os
import shutil

# lop -> lop ke DUOI no (lay mau dinh cua lop duoi lam dich)
CAP = [("b03", "b02"), ("b02", "b01"), ("b01", "b00")]
TY_LE_FADE = 0.34      # 34% duoi cua anh dung de chuyen
TY_LE_THUAN = 0.20     # 20% cuoi cung la mau thuan, khong con van


def mau_dinh(duong_dan):
    a = np.array(Image.open(duong_dan).convert("RGB")).astype(float)
    return a[:6].reshape(-1, 3).mean(axis=0)


for muc in ["x1", "x2", "x3", "x4"]:
    d = f"D:/NRO_Project/Assets/Resources/res/{muc}/bg/"
    for tren, duoi in CAP:
        p_tren, p_duoi = d + tren + ".png", d + duoi + ".png"
        if not (os.path.exists(p_tren) and os.path.exists(p_duoi)):
            continue
        goc = p_tren + ".bak-day"
        if not os.path.exists(goc):
            shutil.copy2(p_tren, goc)

        dich = mau_dinh(p_duoi)
        a = np.array(Image.open(goc).convert("RGB")).astype(float)
        h, w, _ = a.shape
        bd = int(h * (1 - TY_LE_FADE))
        kt = int(h * (1 - TY_LE_THUAN))
        # do do DONG DEU theo chieu ngang cho tung kenh — khong phai std giua 3 kenh mau
        lech_truoc = a[-1].std(axis=0).mean()

        for y in range(bd, h):
            if y >= kt:
                a[y] = dich
            else:
                t = (y - bd) / max(1, kt - bd)
                t = t * t * (3 - 2 * t)
                a[y] = a[y] * (1 - t) + dich * t

        Image.fromarray(np.clip(a, 0, 255).astype(np.uint8)).save(p_tren)
        # ban sao ma game that dung: typeBg = 1 -> b1x
        ban_sao = d + "b1" + tren[-1] + ".png"
        shutil.copy2(p_tren, ban_sao)
        print(f"[{muc}] {tren} day: do gon song {lech_truoc:5.1f} -> {a[-1].std(axis=0).mean():4.1f}"
              f"   khop dinh {duoi} ({dich[0]:.0f},{dich[1]:.0f},{dich[2]:.0f})")
    print()
