# -*- coding: utf-8 -*-
"""Kiem 24 anh dot 2 truoc khi lap vao game.

Kiem dung 6 diem da ghi trong de bai (DOT-2-ANH-CAN-VE.md §"Sau loi hay mac"),
cong mot so kiem rieng cho tung loai. KHONG duyet vi anh dep — duyet vi do duoc.
"""
from PIL import Image
import numpy as np
import os

THU_MUC = "D:/NRO_GoiBanGiao/docs/my-thuat/map-1-lang-la/anh-da-tao"

# ten file -> (rong, cao) mong doi o muc x4
CHO_DOI = {
    "tile_vach_go_trai_x4.png": (96, 96), "tile_vach_go_phai_x4.png": (96, 96),
    "tile_chan_go_x4.png": (96, 96),
    "tile_dat_lienmach_1_x4.png": (96, 96), "tile_dat_lienmach_2_x4.png": (96, 96),
    "tile_dat_tang_giua_x4.png": (96, 96), "tile_dat_tang_day_x4.png": (96, 96),
    "tile_ke_da_x4.png": (96, 96), "tile_bac_thang_x4.png": (96, 96),
    "tile_bac_thang_dinh_x4.png": (96, 96),
    "tile_co_sang_da_x4.png": (96, 96), "tile_da_sang_co_x4.png": (96, 96),
    "nen_nuihokage_rong_x4.png": (1536, 768),
    "nen_lang_xa_rong_x4.png": (1536, 576),
    "bong_nha_nho_x4.png": (576, 96), "bong_nha_lon_x4.png": (768, 96),
    "bong_quan_x4.png": (768, 96), "bong_torii_x4.png": (480, 96),
    "bong_cay_x4.png": (576, 96), "bong_gieng_x4.png": (288, 96),
    "bong_thung_x4.png": (192, 96), "bong_den_x4.png": (192, 96),
    "tt_cay_nho_x4.png": (384, 384), "tt_bui_ram_x4.png": (288, 192),
}

TILE = [t for t in CHO_DOI if t.startswith("tile_")]
NEN = ["nen_nuihokage_rong_x4.png", "nen_lang_xa_rong_x4.png"]
BONG = [t for t in CHO_DOI if t.startswith("bong_")]
VAT = ["tt_cay_nho_x4.png", "tt_bui_ram_x4.png"]

MAU_TROI = np.array([10, 86, 154])       # mau day anh nen phai nhoa ve


def canh_ban(a):
    """Vien trang/den nam trong anh — loi da vap nhieu lan."""
    rgb, al = a[:, :, :3].astype(float), a[:, :, 3]
    xau = []
    for ten, px, m in (("tren", rgb[0], al[0]), ("duoi", rgb[-1], al[-1]),
                       ("trai", rgb[:, 0], al[:, 0]), ("phai", rgb[:, -1], al[:, -1])):
        dac = m > 40
        if dac.sum() == 0:
            continue
        s = px[dac].mean()
        if s > 212:
            xau.append(f"{ten}=TRANG({s:.0f})")
        elif s < 30:
            xau.append(f"{ten}=DEN({s:.0f})")
    return xau


def lap_duoc(a):
    """Tile phai lap voi chinh no: mep trai phai noi duoc vao mep phai."""
    rgb = a[:, :, :3].astype(float)
    doc = np.abs(rgb[:, 0] - rgb[:, -1]).mean()
    ngang = np.abs(rgb[0] - rgb[-1]).mean()
    return doc, ngang


loi, canh_bao, dat = [], [], []

for ten in sorted(CHO_DOI):
    p = os.path.join(THU_MUC, ten)
    if not os.path.exists(p):
        loi.append(f"{ten}: THIEU FILE")
        continue
    im = Image.open(p).convert("RGBA")
    a = np.array(im)
    h, w = a.shape[0], a.shape[1]
    rw, rh = CHO_DOI[ten]
    ghi = []

    # 1) kich thuoc dung tung pixel
    if (w, h) != (rw, rh):
        loi.append(f"{ten}: co {w}x{h}, de bai ghi {rw}x{rh}")
        continue

    # 2) vien ban
    v = canh_ban(a)
    if v:
        loi.append(f"{ten}: vien ban — {' '.join(v)}")

    # 3) rieng tile: phai lap duoc + phai DAC (khong trong suot)
    if ten in TILE:
        doc, ngang = lap_duoc(a)
        if doc > 45 or ngang > 45:
            canh_bao.append(f"{ten}: mep khong noi lien khi lap (doc {doc:.0f}, ngang {ngang:.0f})")
        trong = (a[:, :, 3] < 40).mean()
        if trong > 0.12:
            canh_bao.append(f"{ten}: {trong*100:.0f}% trong suot — tile nen phai dac")
        ghi.append(f"lap doc {doc:.0f}/ngang {ngang:.0f}")

    # 4) rieng nen: day phai nhoa ve mau troi, va KHONG duoc ghep doi
    if ten in NEN:
        day = a[-1, :, :3].astype(float).mean(axis=0)
        lech = np.abs(day - MAU_TROI).sum()
        if lech > 45:
            loi.append(f"{ten}: day = ({day[0]:.0f},{day[1]:.0f},{day[2]:.0f}), "
                       f"de bai ghi phai nhoa ve (10,86,154)")
        gon = a[-1, :, :3].astype(float).std(axis=0).mean()
        if gon > 6:
            loi.append(f"{ten}: hang day con van (do gon song {gon:.1f}) — se lo duong ke ngang")
        # ghep doi: nua trai co giong nua phai khong
        rgb = a[:, :, :3].astype(float)
        trai, phai = rgb[:, :w//2], rgb[:, w//2:]
        khac = np.abs(trai - phai).mean()
        if khac < 12:
            loi.append(f"{ten}: NUA TRAI GIONG NUA PHAI (khac {khac:.1f}) — van la anh ghep doi")
        ghi.append(f"day lech {lech:.0f}, gon {gon:.1f}, hai nua khac {khac:.0f}")

    # 5) rieng bong: phai co alpha, dam nhat o giua
    if ten in BONG:
        al = a[:, :, 3].astype(float)
        if al.max() < 40:
            loi.append(f"{ten}: gan nhu trong suot hoan toan, khong thay bong")
        if (al > 250).mean() > 0.5:
            loi.append(f"{ten}: phan lon DUC hoan toan — bong phai mo dan")
        giua = al[h//2, w//2]
        mep = al[h//2, 2]
        if giua <= mep:
            canh_bao.append(f"{ten}: giua ({giua:.0f}) khong dam hon mep ({mep:.0f})")
        ghi.append(f"alpha giua {giua:.0f} mep {mep:.0f}")

    # 6) rieng vat the: chan phai cham mep duoi anh
    if ten in VAT:
        dac = a[:, :, 3] > 40
        hang_co = np.where(dac.any(axis=1))[0]
        ho = h - 1 - hang_co.max() if len(hang_co) else h
        if ho > 2:
            loi.append(f"{ten}: chan HO {ho}px so voi mep duoi — trong game se lo lung")
        ghi.append(f"chan ho {ho}px")

    dat.append(f"  {ten:32s} {w}x{h}" + (f"  [{', '.join(ghi)}]" if ghi else ""))

print(f"KIEM {len(CHO_DOI)} ANH DOT 2\n" + "=" * 70)
print("\nDAT:")
print("\n".join(dat))
if canh_bao:
    print(f"\nCANH BAO ({len(canh_bao)}) — dung duoc nhung nen xem lai:")
    for c in canh_bao:
        print(f"  ! {c}")
if loi:
    print(f"\nLOI ({len(loi)}) — PHAI SUA TRUOC KHI LAP:")
    for l in loi:
        print(f"  X {l}")
else:
    print("\nKhong co loi chan.")
