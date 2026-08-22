# -*- coding: utf-8 -*-
"""Lap khe ho giua lop lang xa va mat co.

Bon lop nen KHONG phu het chieu cao man:
    b03 may      -118..  74
    b02 nui         24.. 216
    b01 lang xa    186.. 330
    b00 troi       310.. 358      <- chi cao 48 don vi
    mat co                 432
Tu 358 den 432 (74 don vi) khong lop nao che — engine dien mau don, va o moi khe
giua hai nha gan, TROI roi thang xuong tan mat co. Do la thu lam "lang treo lo
lung" ma dot soi doi chieu concept goi ten.

b00 dang la anh TROI TRON. Thay bang phan duoi cua anh lang xa (suong mo + mai
nha mo) thi khe ho co noi dung, giong concept: chan troi bi DONG o cao do mai nha.
"""
from PIL import Image
import os
import shutil

NGUON = "D:/NRO_GoiBanGiao/docs/my-thuat/map-1-lang-la/anh-da-tao/nen_lang_xa_rong_x4.png"
RES = "D:/NRO_Project/Assets/Resources/res"
MUC = {"x1": 0.25, "x2": 0.5, "x3": 0.75, "x4": 1.0}
CAO_B00 = 192      # muc x4; truoc la 96 — gap doi de phu duoc nhieu hon

goc = Image.open(NGUON).convert("RGB")
w, h = goc.size
# lay dai o khoang 45-80% chieu cao: cho co mai nha mo va suong, chua qua nhat
y0, y1 = int(h * 0.45), int(h * 0.80)
dai = goc.crop((0, y0, w, y1)).resize((w, CAO_B00), Image.LANCZOS)
print(f"cat dai {y0}..{y1} cua anh lang xa, keo ve {w}x{CAO_B00}")

for muc, ty in MUC.items():
    im = dai if ty == 1.0 else dai.resize((round(w * ty), round(CAO_B00 * ty)), Image.LANCZOS)
    for ten in ("b00", "b10"):
        p = f"{RES}/{muc}/bg/{ten}.png"
        if os.path.exists(p) and not os.path.exists(p + ".bak-khe"):
            shutil.copy2(p, p + ".bak-khe")
        im.save(p)
    print(f"  [{muc}] b00 + b10 = {im.size[0]}x{im.size[1]}")
