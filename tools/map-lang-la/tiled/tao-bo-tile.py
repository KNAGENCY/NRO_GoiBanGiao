# -*- coding: utf-8 -*-
"""Sinh BO TILE cho Tiled tu chinh anh tile cua client.

    python tools/map-lang-la/tiled/tao-bo-tile.py

Vi sao can file nay: Tiled phai biet MOI o vuong trong map ung voi anh nao thi
no moi ve ra duoc man hinh giong game. Bo tile nay chinh la bang tra do.

MOT DIEU QUAN TRONG ve cach danh so:
  Tiled goi so cua mot o la "gid", bat dau tu 1 va 0 nghia la o TRONG.
  Game goi so cua mot o la "tile id", 0 cung nghia la o TRONG.
  Nen chi can dat firstgid = 1 la co ngay phep doi mot-doi-mot:

      gid = tile_id + 1        (o trong: gid = 0, tile_id = 0)

  Khong bang tra, khong cho nao de sai lech. Day la ly do bo tile nay dung
  kieu "collection of images" (moi tile mot file anh rieng) chu khong phai
  mot tam atlas: atlas ep moi tile cung co, ma tile cua minh KHONG cung co.

TILE KHONG CUNG CO — cho de nham nhat:
  58 o thuong la 96x96 (= 24 don vi game x4). Nhung nha/cay/cong la MOT tile
  voi anh to toi 384x384, neo o o cua no roi tran len tren. Tiled ve tile
  trong lop tile theo goc DUOI-TRAI y het client, nen anh to tran len tren
  dung nhu trong game — nhin thay ngay, khong phai doan.
"""
import os
import re
import sys
from xml.sax.saxutils import quoteattr

from PIL import Image

GOC_TILE = "D:/NRO_Project/Assets/Resources/res/x4/t"
RA = os.path.join(os.path.dirname(os.path.abspath(__file__)), "lang-la.tsx")

# Muc x4: mot o luoi = 24 don vi game x 4 = 96 pixel.
O = 96


def doc_tile():
    """Tra [(tile_id, ten_file, rong, cao)] sap theo tile_id."""
    ds = []
    for f in os.listdir(GOC_TILE):
        m = re.fullmatch(r"1\$(\d+)\.png", f)
        if not m:
            continue
        with Image.open(os.path.join(GOC_TILE, f)) as im:
            w, h = im.size
        ds.append((int(m.group(1)), f, w, h))
    return sorted(ds)


def main():
    if not os.path.isdir(GOC_TILE):
        print(f"khong thay thu muc tile: {GOC_TILE}")
        return 1
    ds = doc_tile()
    if not ds:
        print("khong doc duoc tile nao")
        return 1

    goc_tuong_doi = os.path.relpath(GOC_TILE, os.path.dirname(RA)).replace("\\", "/")
    d = [
        '<?xml version="1.0" encoding="UTF-8"?>',
        # tilecount/columns cua collection-of-images chi la goi y hien thi.
        f'<tileset version="1.10" tiledversion="1.11.0" name="lang-la"'
        f' tilewidth="{O}" tileheight="{O}" tilecount="{len(ds)}" columns="0">',
        ' <grid orientation="orthogonal" width="%d" height="%d"/>' % (O, O),
    ]
    for tid, ten, w, h in ds:
        nguon = quoteattr(f"{goc_tuong_doi}/{ten}")
        d.append(f' <tile id="{tid}">')
        d.append(f'  <image width="{w}" height="{h}" source={nguon}/>')
        d.append(" </tile>")
    d.append("</tileset>")

    with open(RA, "w", encoding="utf-8", newline="\n") as f:
        f.write("\n".join(d) + "\n")

    to = [x for x in ds if (x[2], x[3]) != (O, O)]
    print(f"da ghi {RA}")
    print(f"  {len(ds)} tile, trong do {len(to)} tile TO hon mot o:")
    for tid, ten, w, h in to[:12]:
        print(f"     id {tid:3}  {w}x{h}  ({w // O}x{h // O} o)")
    if len(to) > 12:
        print(f"     ... con {len(to) - 12} cai nua")
    return 0


if __name__ == "__main__":
    sys.exit(main())
