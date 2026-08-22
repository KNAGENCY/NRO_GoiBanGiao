# -*- coding: utf-8 -*-
"""Cam sprite nhan vat ban CHINH THUC (ArtRoot production-closure-v11) vao client.

Vi sao can: part.bin tro toi sprite id 27041-27082, nhung trong Resources cua
client chi co ban thu nghiem cu Small9000/9001/9002. Client tai part.bin moi ->
tro toi 27041 -> khong tim thay anh -> nhan vat lai vo hinh.

Nguon:  docs/my-thuat/anh-nguon/vertical-slice-v1/generated/production-closure-v11/
            sprites/<id>/x<z>/sprite.png
Dich:   D:/NRO_Project/Assets/Resources/res/x<z>/smallimage/Small<id>.png
"""
import os
import shutil

NGUON = ("D:/NRO_GoiBanGiao/docs/my-thuat/anh-nguon/vertical-slice-v1/"
         "generated/production-closure-v11/sprites")
DICH = "D:/NRO_Project/Assets/Resources/res/{muc}/smallimage"
MUC = ["x1", "x2", "x3", "x4"]

ids = sorted(int(x) for x in os.listdir(NGUON) if x.isdigit())
print(f"nguon co {len(ids)} sprite: {ids[0]}..{ids[-1]}")

tong, thieu = 0, []
for muc in MUC:
    thu_muc = DICH.format(muc=muc)
    os.makedirs(thu_muc, exist_ok=True)
    n = 0
    for i in ids:
        src = f"{NGUON}/{i}/{muc}/sprite.png"
        if not os.path.exists(src):
            thieu.append(f"{i}/{muc}")
            continue
        shutil.copy2(src, f"{thu_muc}/Small{i}.png")
        n += 1
    tong += n
    print(f"  [{muc}] chep {n}/{len(ids)} anh")

print(f"\ntong {tong} anh da cam vao client")
if thieu:
    print(f"THIEU {len(thieu)} anh nguon: {thieu[:10]}")
