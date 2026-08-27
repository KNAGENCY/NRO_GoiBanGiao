# -*- coding: utf-8 -*-
"""Doi map cua game -> file .tmx mo duoc bang Tiled.

    python tools/map-lang-la/tiled/map-sang-tiled.py

Chay de dua bo cuc DANG CO vao Tiled — khong mat gi da lam. Tu day tro di bo
cuc song trong .tmx chu khong con nam rai rac duoi dang 93 con so cung trong
ve-bo-cuc-map.py, noi ma sua mot so la ve lai ca 1248 o.

File .tmx sinh ra co hai lop:
  "concept"  anh concept da duyet, nam DUOI cung lam giay can, mo 45%. Xep nha
             de len no la thay ngay minh co bam concept khong. Day dung la thu
             thieu tu dau: concept mot noi, bo cuc mot noi, khong gi buoc chung
             khop nhau nen chung troi dan ra.
  "mat-dat"  lop tile that.

Vat the to bi DICH XUONG khi sang Tiled — xem chu thich QUY UOC NEO o chung.py.
"""
import os
import sys

from xml.sax.saxutils import quoteattr

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import chung  # noqa: E402

MAP = "D:/NRO_Project/Assets/Resources/res/mymap/1.bytes"
CONCEPT = "D:/NRO_GoiBanGiao/docs/my-thuat/map-1-lang-la/concept-map1-langla.png"
RA = os.path.join(os.path.dirname(os.path.abspath(__file__)), "map-1-lang-la.tmx")
BO_TILE = "lang-la.tsx"


def main():
    if not os.path.exists(MAP):
        print(f"khong thay map: {MAP}")
        return 1
    w, h, o = chung.doc_map_game(MAP)
    gid, va_cham, tran = chung.game_sang_tiled(w, h, o)

    if va_cham or tran:
        print("!! KHONG doi duoc sang Tiled ma khong mat du lieu:")
        for t, r, c, khac, r2 in va_cham:
            print(f"   tile {t} o (hang {r}, cot {c}) dich xuong hang {r2} thi DE LEN tile {khac - 1}")
        for t, r, c, r2 in tran:
            print(f"   tile {t} o (hang {r}, cot {c}) dich xuong hang {r2} thi ra NGOAI map")
        print("   Sua bo cuc cho het chong lan roi chay lai.")
        return 1

    hang = []
    for r in range(h):
        hang.append(",".join(str(gid[r * w + c]) for c in range(w)))
    csv = ",\n".join(hang)

    lop_concept = ""
    if os.path.exists(CONCEPT):
        from PIL import Image
        with Image.open(CONCEPT) as im:
            cw, ch = im.size
        rong = w * chung.O
        cao = int(ch * rong / cw)
        y = max(0, h * chung.O - cao)
        src = quoteattr(os.path.relpath(CONCEPT, os.path.dirname(RA)).replace("\\", "/"))
        lop_concept = (
            f' <imagelayer id="1" name="concept" offsetx="0" offsety="{y}" opacity="0.45">\n'
            f'  <image source={src} width="{rong}" height="{cao}"/>\n'
            f" </imagelayer>\n"
        )

    tmx = (
        '<?xml version="1.0" encoding="UTF-8"?>\n'
        f'<map version="1.10" tiledversion="1.11.0" orientation="orthogonal"'
        f' renderorder="right-down" width="{w}" height="{h}"'
        f' tilewidth="{chung.O}" tileheight="{chung.O}" infinite="0"'
        f' nextlayerid="3" nextobjectid="1">\n'
        f' <tileset firstgid="1" source="{BO_TILE}"/>\n'
        f"{lop_concept}"
        f' <layer id="2" name="mat-dat" width="{w}" height="{h}">\n'
        f'  <data encoding="csv">\n{csv}\n</data>\n'
        f" </layer>\n"
        "</map>\n"
    )
    with open(RA, "w", encoding="utf-8", newline="\n") as f:
        f.write(tmx)

    to = sum(1 for x in o if x and chung.cao_o(x) > 1)
    print(f"da ghi {RA}")
    print(f"  map {w}x{h} = {w * h} o, {sum(1 for x in o if x)} o co tile")
    print(f"  {to} vat the to da duoc dich xuong cho khop cach Tiled ve")
    print(f"  lop concept: {'co' if lop_concept else 'KHONG (khong thay anh concept)'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
