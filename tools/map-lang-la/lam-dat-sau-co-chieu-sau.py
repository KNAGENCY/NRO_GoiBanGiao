# -*- coding: utf-8 -*-
"""Lam tang dat sau co chieu sau, het canh "man nau phang li co vai mieng va vuong".

Hien trang: gan 1/4 man hinh duoi chi co 2 mau o — 60 (dat phang) va 61 (soi).
O soi khac han o dat nen noi len nhu mieng va vuong sac canh, con lai la mot mang
nau chet.

Cach chua bang code (khong can ve moi):
  - 61 tro thanh bien the NHE cua 60 (lech sang 7%) -> rai day duoc ma khong lo o
  - 67, 68 (dang bo trong) tro thanh dat TOI DAN theo do sau
  - them nhieu hat rat nhe, khac nhau tung o, de mat khong bat duoc quy luat

Ket qua: cang xuong sau cang toi, giong dat that, khong con mang phang.
"""
from PIL import Image
import numpy as np
import os
import shutil

# tile dich -> (he so sang, muc nhieu)
# Chenh sang giua cac o phai NHO. Thu 48% thi tang dat bien thanh ban co lo ro
# tung o vuong 48px — te hon ca mang phang ban dau. Mat nguoi bat ranh o rat
# nhay khi hai o ke nhau lech qua ~8%.
CONG_THUC = {
    61: (1.035, 4),   # bien the rat nhe cua dat chuan, rai xen ke
    67: (0.945, 4),   # tang giua
    68: (0.885, 4),   # tang day
}


def bam(x):
    x = (x ^ (x >> 13)) * 1274126177
    return (x ^ (x >> 16)) & 0xFFFFFFFF


for muc in ["x1", "x2", "x3", "x4"]:
    d = f"D:/NRO_Project/Assets/Resources/res/{muc}/t/"
    nguon = d + "1$60.png"
    if not os.path.exists(nguon):
        print(f"[{muc}] khong co tile 60, bo qua"); continue
    goc = np.array(Image.open(nguon).convert("RGBA")).astype(float)
    h, w, _ = goc.shape
    for tile, (he_so, nhieu) in CONG_THUC.items():
        p = d + f"1${tile}.png"
        if os.path.exists(p) and not os.path.exists(p + ".bak-datsau"):
            shutil.copy2(p, p + ".bak-datsau")
        a = goc.copy()
        a[:, :, :3] *= he_so
        # nhieu hat: moi pixel lech mot chut theo mot ham bam on dinh
        if nhieu:
            r = np.fromiter((bam(y * 8191 + x * 131 + tile * 7919) % (2 * nhieu + 1) - nhieu
                             for y in range(h) for x in range(w)), dtype=float, count=h * w)
            a[:, :, :3] += r.reshape(h, w, 1)
        Image.fromarray(np.clip(a, 0, 255).astype(np.uint8), "RGBA").save(p)
    print(f"[{muc}] tao lai tile 61/67/68 tu tile 60 ({w}x{h})")

print("\nxong — nho chay lai ve-bo-cuc-map.py de dung tang moi")
