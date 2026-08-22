# -*- coding: utf-8 -*-
# Do vien trang/den nam TRONG anh tile — nguyen nhan "o vuong vien trang" trong game.
from PIL import Image
import numpy as np

TEN = {60: "dat_1", 61: "dat_2", 62: "mat_co", 63: "goc_trai", 64: "goc_phai",
       65: "vach_trai", 66: "vach_phai", 67: "tile_67", 68: "tile_68", 69: "duong_da"}
THU_MUC = "D:/NRO_Project/Assets/Resources/res/x4/t/"

for i in range(60, 70):
    p = THU_MUC + "1$" + str(i) + ".png"
    im = Image.open(p).convert("RGBA")
    a = np.array(im).astype(int)
    h, w, _ = a.shape
    rgb, alpha = a[:, :, :3], a[:, :, 3]
    canh = [("tren", rgb[0], alpha[0]), ("duoi", rgb[-1], alpha[-1]),
            ("trai", rgb[:, 0], alpha[:, 0]), ("phai", rgb[:, -1], alpha[:, -1])]
    xau = []
    for ten_canh, px, al in canh:
        dac = al > 40
        if dac.sum() == 0:
            continue
        sang = px[dac].mean()
        if sang > 212:
            xau.append(f"{ten_canh}=TRANG({sang:.0f})")
        elif sang < 30:
            xau.append(f"{ten_canh}=DEN({sang:.0f})")
    print(f"{i} {TEN[i]:11s} {w}x{h}  " + (" ".join(xau) if xau else "sach"))
