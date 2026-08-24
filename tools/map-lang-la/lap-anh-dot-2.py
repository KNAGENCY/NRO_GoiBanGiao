# -*- coding: utf-8 -*-
"""Lap 24 anh dot 2 vao client, thu nho san cho 4 muc zoom.

Anh giao o muc x4 (goc). x1/x2/x3 tu thu nho tu do.

Anh nen (b0x) con phai chep sang ban b1x vi client nap theo typeBg = 1 —
loi da vap, xem LOI-DA-GAP.md §4.
"""
from PIL import Image
import os
import shutil

NGUON = "D:/NRO_GoiBanGiao/docs/my-thuat/map-1-lang-la/anh-da-tao"
RES = "D:/NRO_Project/Assets/Resources/res"
MUC = {"x1": 0.25, "x2": 0.5, "x3": 0.75, "x4": 1.0}

# --- tile nen: ten file -> so hieu tile trong map ---
TILE = {
    "tile_dat_lienmach_1_x4.png": 60,   # dat tang mat
    "tile_dat_lienmach_2_x4.png": 61,   # bien the
    "tile_co_sang_da_x4.png":     63,   # chuyen co -> lat da
    "tile_da_sang_co_x4.png":     64,   # chuyen lat da -> co
    "tile_vach_go_trai_x4.png":   65,
    "tile_vach_go_phai_x4.png":   66,
    "tile_dat_tang_giua_x4.png":  67,
    "tile_dat_tang_day_x4.png":   68,
    # slot moi
    "tile_chan_go_x4.png":        83,
    "tile_ke_da_x4.png":          84,
    "tile_bac_thang_x4.png":      85,
    "tile_bac_thang_dinh_x4.png": 86,
}

# --- vat the (cung dung slot tile, chi khac la anh to hon mot o) ---
VAT_THE = {
    "tt_cay_nho_x4.png": 87,
    "tt_bui_ram_x4.png": 88,
}

# --- anh nen ---
NEN = {
    "nen_nuihokage_rong_x4.png": "b02",
    "nen_lang_xa_rong_x4.png":   "b01",
}


def thu_nho(im, ty_le):
    if ty_le == 1.0:
        return im
    w, h = im.size
    return im.resize((max(1, round(w * ty_le)), max(1, round(h * ty_le))), Image.LANCZOS)


def luu(im, duong_dan):
    thu_muc = os.path.dirname(duong_dan)
    os.makedirs(thu_muc, exist_ok=True)
    if os.path.exists(duong_dan) and not os.path.exists(duong_dan + ".bak-dot1"):
        shutil.copy2(duong_dan, duong_dan + ".bak-dot1")
    im.save(duong_dan)


dem = 0
for ten, so in list(TILE.items()) + list(VAT_THE.items()):
    p = os.path.join(NGUON, ten)
    if not os.path.exists(p):
        print(f"  THIEU {ten}"); continue
    goc = Image.open(p).convert("RGBA")
    for muc, ty in MUC.items():
        luu(thu_nho(goc, ty), f"{RES}/{muc}/t/1${so}.png")
        dem += 1
    print(f"  tile {so:2d} <- {ten}  ({goc.size[0]}x{goc.size[1]})")

for ten, dich in NEN.items():
    p = os.path.join(NGUON, ten)
    if not os.path.exists(p):
        print(f"  THIEU {ten}"); continue
    goc = Image.open(p).convert("RGB")
    for muc, ty in MUC.items():
        im = thu_nho(goc, ty)
        luu(im, f"{RES}/{muc}/bg/{dich}.png")
        # ban sao ma client THAT nap (typeBg = 1)
        luu(im, f"{RES}/{muc}/bg/b1{dich[-1]}.png")
        dem += 2
    print(f"  nen {dich} (+b1{dich[-1]}) <- {ten}  ({goc.size[0]}x{goc.size[1]})")

print(f"\nda ghi {dem} file anh vao client")
print("CHUA LAP: 8 anh bong — luoi tile chi co MOT lop nen bong se de len o mat co,")
print("          can lop bgItem rieng (may chu phai gui). Ghi vao ke hoach.")
