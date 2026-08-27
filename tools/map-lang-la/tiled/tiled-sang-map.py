# -*- coding: utf-8 -*-
"""Doi file .tmx (vua sua trong Tiled) -> map cua game.

    python tools/map-lang-la/tiled/tiled-sang-map.py [--thu]

    --thu : chi kiem, KHONG ghi de map that. Dung cho vong khu hoi.

Vat the to duoc dich NGUOC LEN — xem chu thich QUY UOC NEO o chung.py.

Sau buoc nay CON HAI VIEC, thieu la khong thay gi doi trong game:
  1. python tools/map-lang-la/sinh-va-cham.py
  2. TANG vsData ben client — client giu ban cache theo so nay. Da mac bay
     nay hai lan: sua map xong vao game thay y nguyen, tuong code hong.
"""
import os
import shutil
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import chung  # noqa: E402

MAP = "D:/NRO_Project/Assets/Resources/res/mymap/1.bytes"
TMX = os.path.join(os.path.dirname(os.path.abspath(__file__)), "map-1-lang-la.tmx")


def main(argv):
    chi_thu = "--thu" in argv
    if not os.path.exists(TMX):
        print(f"khong thay {TMX} — chay map-sang-tiled.py truoc")
        return 1

    w, h, gid = chung.doc_lop_tile_tmx(TMX)
    o, va_cham, tran = chung.tiled_sang_game(w, h, gid)

    if va_cham or tran:
        print("!! KHONG doi nguoc duoc ma khong mat du lieu:")
        for t, r, c, khac, r2 in va_cham:
            print(f"   tile {t} o (hang {r}, cot {c}) dich len hang {r2} thi DE LEN tile {khac}")
        for t, r, c, r2 in tran:
            print(f"   tile {t} o (hang {r}, cot {c}) dich len hang {r2} thi ra NGOAI map")
            print("      (vat the cao qua, dat qua sat mep tren)")
        return 1

    cu = open(MAP, "rb").read() if os.path.exists(MAP) else b""
    if cu:
        cw, ch = cu[0], cu[1]
        if (cw, ch) != (w, h):
            print(f"!! kich thuoc doi: map cu {cw}x{ch}, .tmx {w}x{h}")
            print("   Doi kich thuoc map la viec LON — va cham, spawn, tuong bien deu theo no.")
            return 1

    moi = bytes([w, h]) + bytes(o)

    if chi_thu:
        giong = moi == cu
        print(f"map cu : {len(cu)} byte")
        print(f"map moi: {len(moi)} byte")
        print("KHOP TUNG BYTE" if giong else "KHAC NHAU")
        if not giong and cu:
            khac = [i for i in range(min(len(cu), len(moi))) if cu[i] != moi[i]]
            print(f"  {len(khac)} byte khac, vi tri dau: {khac[:8]}")
        return 0 if giong else 2

    if cu:
        shutil.copy2(MAP, MAP + ".bak-tiled")
    with open(MAP, "wb") as f:
        f.write(moi)
    print(f"da ghi {MAP} ({w}x{h}, {sum(1 for x in o if x)} o co tile)")
    print()
    print("CON HAI BUOC NUA, thieu la khong thay gi doi trong game:")
    print("  1. python tools/map-lang-la/sinh-va-cham.py")
    print("  2. TANG vsData ben client (client giu ban cache theo so nay)")
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
