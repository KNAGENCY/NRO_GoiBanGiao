# -*- coding: utf-8 -*-
"""Tao 2 o chuyen tiep giua mat co va loi lat da.

Hien co va loi da gap nhau bang MOT NHAT CAT DOC tuyet doi phang — mau nhay 67%
trong dung mot pixel, ranh gioi la duong thang cao 48px, khong mot pixel pha tron.

Cach lam: tron o co (62) va o lat da (69) qua mot ranh gioi LOM CHOM, dung ham
bam on dinh nen ranh gioi khong bao gio doi giua cac lan chay.

Dat vao slot 63 va 64 — hai o "doc cheo" cu, dang khong dung den va DA duoc khai
bao la san trong tile_set_info nen khong phai sua gi them.
"""
from PIL import Image
import numpy as np
import os
import shutil


def bam(x):
    x = (x ^ (x >> 13)) * 1274126177
    return (x ^ (x >> 16)) & 0xFFFFFFFF


for muc in ["x1", "x2", "x3", "x4"]:
    d = f"D:/NRO_Project/Assets/Resources/res/{muc}/t/"
    p_co, p_da = d + "1$62.png", d + "1$69.png"
    if not (os.path.exists(p_co) and os.path.exists(p_da)):
        print(f"[{muc}] thieu anh nguon, bo qua"); continue
    co = np.array(Image.open(p_co).convert("RGBA")).astype(float)
    da = np.array(Image.open(p_da).convert("RGBA")).astype(float)
    h, w, _ = co.shape
    if da.shape != co.shape:
        print(f"[{muc}] hai anh khac co, bo qua"); continue

    # ranh gioi lom chom: moi hang mot vi tri hoi khac, quanh giua o
    # bien do phai DU LON moi thay: +-w/4. Thu +-w/8 truoc, sau khi lam muot
    # chi con +-3px o muc x2 — mat thuong khong phan biet duoc voi duong thang.
    bien = np.array([w * 0.5 + (bam(y * 7919 + 13) % (w // 2)) - w // 4 for y in range(h)])
    # lam muot theo chieu doc de khoi rang cua qua deu
    for _ in range(2):
        bien = np.convolve(np.pad(bien, 2, mode="edge"), np.ones(5) / 5, mode="valid")

    xs = np.arange(w)[None, :]
    # 0 = co, 1 = da; chuyen dan qua mot dai rong ~6% be ngang
    dai = max(2.0, w * 0.10)
    t = np.clip((xs - bien[:, None]) / dai + 0.5, 0.0, 1.0)
    t = t * t * (3 - 2 * t)

    for tile, mat_na in ((63, t), (64, 1.0 - t)):
        p = d + f"1${tile}.png"
        goc = p + ".bak-chuyentiep"
        if os.path.exists(p) and not os.path.exists(goc):
            shutil.copy2(p, goc)
        a = co * (1 - mat_na[:, :, None]) + da * mat_na[:, :, None]
        Image.fromarray(np.clip(a, 0, 255).astype(np.uint8), "RGBA").save(p)
    print(f"[{muc}] tao o chuyen tiep 63 (co->da) va 64 (da->co), {w}x{h}")

print("\nxong — nho chay lai ve-bo-cuc-map.py de dat hai o nay vao hai dau loi di")
