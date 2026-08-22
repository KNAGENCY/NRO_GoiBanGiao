# -*- coding: utf-8 -*-
"""Tao anh chan dung RIENG cho nhan vat nguoi choi, vua o xanh tren HUD.

Truoc do dung Small27078 — nhung do la anh 'iruka-portrait', chinh la chan dung
NPC thay Iruka trong catalog. Nguoi choi va NPC dung chung mot khuon mat.

Anh moi cat tu chinh anh DAU cua nhan vat (27041) nen dung mat nguoi choi, va
thu nho cho vua o xanh tren HUD (o do cao khoang 28 don vi).
"""
from PIL import Image
import numpy as np
import os

NGUON = 27041          # anh dau nhan vat
DICH = 27084           # id moi cho chan dung nguoi choi
CANH = {"x1": 26, "x2": 52, "x3": 78, "x4": 104}   # 26 don vi game

for muc, canh in CANH.items():
    d = f"D:/NRO_Project/Assets/Resources/res/{muc}/smallimage/"
    p_nguon = d + f"Small{NGUON}.png"
    if not os.path.exists(p_nguon):
        print(f"[{muc}] khong co Small{NGUON}, bo qua"); continue
    im = Image.open(p_nguon).convert("RGBA")
    a = np.array(im)
    ys, xs = np.where(a[:, :, 3] > 20)
    # cat sat noi dung roi dat vao khung vuong, chua le 1 pixel
    cat = im.crop((xs.min(), ys.min(), xs.max() + 1, ys.max() + 1))
    w, h = cat.size
    canh_noi_dung = max(w, h)
    khung = Image.new("RGBA", (canh_noi_dung, canh_noi_dung), (0, 0, 0, 0))
    khung.alpha_composite(cat, ((canh_noi_dung - w) // 2, (canh_noi_dung - h) // 2))
    khung.resize((canh, canh), Image.LANCZOS).save(d + f"Small{DICH}.png")
    print(f"[{muc}] Small{DICH}.png {canh}x{canh}  (cat tu Small{NGUON} vung {w}x{h})")

print(f"\nxong — nho doi AvatarNhanVat trong server thanh {DICH}")
