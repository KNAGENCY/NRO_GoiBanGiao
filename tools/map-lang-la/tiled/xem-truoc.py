# -*- coding: utf-8 -*-
"""Ve map ra anh PNG de NHIN NGAY, khong can mo Unity.

    python tools/map-lang-la/tiled/xem-truoc.py [--khung] [--game]

    --khung : ve them luoi o va so cot/hang — dung khi can go toa do.
    --game  : ve tu res/mymap/1.bytes thay vi tu .tmx (de doi chieu hai ben).

Vi sao can: truoc gio muon xem map doi the nao phai chay Unity, bam Play, cho
nap, chup man hinh — moi vong vai phut, nen moi lan sua la mot lan doan mo.
File nay ve trong duoi MOT GIAY.

Anh ra KHONG phai anh game that: khong co 4 lop nen, khong nhan vat, khong HUD.
No de xem BO CUC — dung cai hay xep hong.

Ve theo QUY UOC NEO cua tung ben (xem chung.py): doc .tmx thi neo DUOI-trai
(giong Tiled), doc map game thi neo TREN-trai (giong client). Hai duong deu ra
CUNG mot buc anh — do chinh la phep kiem chung rang cau noi khong lam lech gi.
"""
import os
import sys

from PIL import Image, ImageDraw

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import chung  # noqa: E402

TMX = os.path.join(os.path.dirname(os.path.abspath(__file__)), "map-1-lang-la.tmx")
MAP = "D:/NRO_Project/Assets/Resources/res/mymap/1.bytes"
RA = os.path.join(os.path.dirname(os.path.abspath(__file__)), "xem-truoc.png")

NEN = (150, 196, 224)  # troi nhat, chi de phan biet o trong voi o co tile


def nap_anh_tile():
    kho = {}
    import re
    for f in os.listdir(chung.GOC_TILE):
        m = re.fullmatch(r"1\$(\d+)\.png", f)
        if m:
            kho[int(m.group(1))] = Image.open(os.path.join(chung.GOC_TILE, f)).convert("RGBA")
    return kho


def main(argv):
    tu_game = "--game" in argv
    O = chung.O

    if tu_game:
        if not os.path.exists(MAP):
            print(f"khong thay {MAP}")
            return 1
        w, h, o = chung.doc_map_game(MAP)
        # client neo goc TREN-trai: anh bat dau ngay tai o neo
        dat = [(r, c, o[r * w + c]) for r in range(h) for c in range(w) if o[r * w + c]]
        def toa_do(r, c, t):
            return c * O, r * O
        nguon = "res/mymap/1.bytes (quy uoc client: neo TREN-trai)"
    else:
        if not os.path.exists(TMX):
            print(f"khong thay {TMX} — chay map-sang-tiled.py truoc")
            return 1
        w, h, gid = chung.doc_lop_tile_tmx(TMX)
        dat = [(r, c, gid[r * w + c] - 1) for r in range(h) for c in range(w) if gid[r * w + c]]
        def toa_do(r, c, t):
            # Tiled neo goc DUOI-trai: anh ket thuc o mep duoi cua o neo
            return c * O, (r + 1) * O - chung.cao_o(t) * O
        nguon = f"{os.path.basename(TMX)} (quy uoc Tiled: neo DUOI-trai)"

    kho = nap_anh_tile()
    anh = Image.new("RGBA", (w * O, h * O), NEN + (255,))
    thieu = set()
    for r, c, t in dat:
        im = kho.get(t)
        if im is None:
            thieu.add(t)
            continue
        anh.alpha_composite(im, toa_do(r, c, t))

    if "--khung" in argv:
        v = ImageDraw.Draw(anh)
        for c in range(w + 1):
            v.line([(c * O, 0), (c * O, h * O)], fill=(255, 0, 0, 60), width=1)
        for r in range(h + 1):
            v.line([(0, r * O), (w * O, r * O)], fill=(255, 0, 0, 60), width=1)
        for c in range(0, w, 5):
            v.text((c * O + 4, 2), str(c), fill=(200, 0, 0, 255))
        for r in range(0, h, 2):
            v.text((3, r * O + 4), str(r), fill=(200, 0, 0, 255))

    anh.convert("RGB").save(RA)
    print(f"da ve {RA}  ({w * O}x{h * O} px)")
    print(f"  nguon: {nguon}")
    print(f"  {len(dat)} o co tile / {w * h} o")
    if thieu:
        print(f"  !! {len(thieu)} tile id KHONG co anh: {sorted(thieu)}")
        print("     (trong game se la o trong — day thuong la loi that)")
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
