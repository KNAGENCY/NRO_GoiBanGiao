# -*- coding: utf-8 -*-
"""Cat bo phan TROI thua o dinh anh lang xa, de lang che duoc chan nui.

VI SAO. Trong concept, chan nui Hokage khong he fade — no bi GIAU SAU mai nha
cua lang xa. Do la ly do bức tranh co chieu sau: troi bi DONG lai o cao do mai
nha, khong roi thang xuong mat co.

Anh lang xa moi ve rat dep nhung 25% tren cung la TROI. Lop nay dat o cao do
186..330 don vi, ma chan nui o 216 — nghia la mai nha dau tien chi bat dau tu
222, THAP HON chan nui. Ket qua: giua nui va nha gan lo ra mot khoang troi rong,
lang trong nhu treo lo lung.

CACH SUA. Cat bo phan troi tren, roi keo phan con lai ve du chieu cao cu. Lang
se bat dau ngay tu mep tren lop, che kin chan nui.
"""
from PIL import Image
import numpy as np
import os
import shutil

NGUON = "D:/NRO_GoiBanGiao/docs/my-thuat/map-1-lang-la/anh-da-tao/nen_lang_xa_rong_x4.png"
RES = "D:/NRO_Project/Assets/Resources/res"
MUC = {"x1": 0.25, "x2": 0.5, "x3": 0.75, "x4": 1.0}

goc = Image.open(NGUON).convert("RGB")
a = np.array(goc).astype(int)
h, w, _ = a.shape

# Tim hang dau tien co MAI NHA / CAY (khong con la troi thuan)
r, g, b = a[:, :, 0], a[:, :, 1], a[:, :, 2]
troi = (b > r + 40) & (b > 120)
ty_le_troi = troi.mean(axis=1)
dau_lang = next((y for y in range(h) if ty_le_troi[y] < 0.80), 0)
print(f"anh goc {w}x{h}: troi thuan tu hang 0 den {dau_lang} ({dau_lang/h*100:.0f}% chieu cao)")

# Cat tu do (chua lai mot chut troi cho mem), roi keo ve du chieu cao cu
chua = int(h * 0.04)
cat = goc.crop((0, max(0, dau_lang - chua), w, h)).resize((w, h), Image.LANCZOS)
print(f"cat tu hang {max(0, dau_lang - chua)}, keo lai ve {w}x{h}")

for muc, ty in MUC.items():
    im = cat if ty == 1.0 else cat.resize((round(w * ty), round(h * ty)), Image.LANCZOS)
    for ten in ("b01", "b11"):
        p = f"{RES}/{muc}/bg/{ten}.png"
        if os.path.exists(p) and not os.path.exists(p + ".bak-cattroi"):
            shutil.copy2(p, p + ".bak-cattroi")
        im.save(p)
    print(f"  [{muc}] ghi b01 + b11 ({im.size[0]}x{im.size[1]})")
