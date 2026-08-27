# -*- coding: utf-8 -*-
"""Thu dung chung cua cau noi Tiled — quan trong nhat la QUY UOC NEO.

=============================================================================
 HAI BEN NEO VAT THE TO O HAI GOC KHAC NHAU. Doc ky truoc khi sua bat cu gi.
=============================================================================

Vat the to (nha, cay, cong torii) la MOT tile id voi anh cao nhieu o. Cau hoi
la: o luoi ghi so do nam o dau cua cai anh?

  client game : goc TREN-trai. Nha cao 6 o dat o hang 12 -> chiem hang 12..17.
  Tiled       : goc DUOI-trai. Cung so do o hang 12 -> chiem hang 7..12.

Lech nhau dung (chieu_cao_o - 1) hang. Khong xu ly thi mo Tiled ra se thay
nha treo lo lung tren troi — va do la LOI CUA CONG CU, khong phai loi cua map.

Da tu mac dung bay nay: ve xem truoc bang quy uoc duoi-trai, nhin thay 16 vat
the deu "treo cach mat dat 1 den 7 o", suyt bao la map hong. Do lai bang quy
uoc tren-trai thi 16/16 CHAM DAT chinh xac.

Nen cau noi doi toa do khi qua lai:

    game -> tiled :  hang_tiled = hang_game + cao_o - 1
    tiled -> game :  hang_game  = hang_tiled - cao_o + 1

Tile thuong cao 1 o nen cao_o - 1 = 0: khong doi gi. Chi vat the to bi dich.
"""
import os
import re

from PIL import Image

GOC_TILE = "D:/NRO_Project/Assets/Resources/res/x4/t"
O = 96  # mot o luoi o muc x4 = 24 don vi game x 4

_cao = None


def cao_o(tile_id):
    """Chieu cao cua tile tinh bang SO O. Tile khong co anh coi nhu 1 o."""
    global _cao
    if _cao is None:
        _cao = {}
        for f in os.listdir(GOC_TILE):
            m = re.fullmatch(r"1\$(\d+)\.png", f)
            if not m:
                continue
            with Image.open(os.path.join(GOC_TILE, f)) as im:
                _cao[int(m.group(1))] = max(1, im.height // O)
    return _cao.get(tile_id, 1)


def kich_thuoc_tile():
    """Tra {tile_id: (rong_px, cao_px)} cho moi tile co anh."""
    ds = {}
    for f in os.listdir(GOC_TILE):
        m = re.fullmatch(r"1\$(\d+)\.png", f)
        if not m:
            continue
        with Image.open(os.path.join(GOC_TILE, f)) as im:
            ds[int(m.group(1))] = im.size
    return ds


def doc_map_game(duong_dan):
    """Doc res/mymap/N.bytes -> (rong, cao, [tile_id...])."""
    d = open(duong_dan, "rb").read()
    w, h = d[0], d[1]
    o = list(d[2:2 + w * h])
    if len(o) != w * h:
        raise ValueError(f"file map hong: can {w * h} o, doc duoc {len(o)}")
    return w, h, o


def doc_lop_tile_tmx(duong_dan):
    """Doc lop tile CSV dau tien trong .tmx -> (rong, cao, [gid...])."""
    s = open(duong_dan, encoding="utf-8").read()
    m = re.search(r"<map\b[^>]*\bwidth=\"(\d+)\"[^>]*\bheight=\"(\d+)\"", s)
    if not m:
        raise ValueError("khong doc duoc kich thuoc map trong .tmx")
    w, h = int(m.group(1)), int(m.group(2))
    m = re.search(r"<layer\b[^>]*>.*?<data\s+encoding=\"csv\">(.*?)</data>", s, re.S)
    if not m:
        raise ValueError(
            "khong thay lop tile ma hoa CSV. Trong Tiled: Map > Map Properties > "
            "Tile Layer Format = CSV, roi luu lai."
        )
    gid = [int(x) for x in re.findall(r"-?\d+", m.group(1))]
    if len(gid) != w * h:
        raise ValueError(f"lop tile co {len(gid)} o, map khai bao {w}x{h} = {w * h}")
    return w, h, gid


def game_sang_tiled(w, h, o):
    """Dich vat the to XUONG (cao_o - 1) hang de Tiled ve dung nhu game."""
    ra = [0] * (w * h)
    va_cham = []
    tran = []
    for r in range(h):
        for c in range(w):
            t = o[r * w + c]
            if not t:
                continue
            r2 = r + cao_o(t) - 1
            if r2 >= h:
                tran.append((t, r, c, r2))
                continue
            if ra[r2 * w + c]:
                va_cham.append((t, r, c, ra[r2 * w + c], r2))
                continue
            ra[r2 * w + c] = t + 1  # gid = tile_id + 1
    return ra, va_cham, tran


def tiled_sang_game(w, h, gid):
    """Chieu nguoc lai: dich vat the to LEN (cao_o - 1) hang."""
    ra = [0] * (w * h)
    va_cham = []
    tran = []
    for r in range(h):
        for c in range(w):
            g = gid[r * w + c]
            if not g:
                continue
            if g & 0xE0000000:
                raise ValueError(
                    f"o (hang {r}, cot {c}) bi LAT hoac XOAY trong Tiled. "
                    "Game khong ho tro — hay dung tile khac thay vi lat."
                )
            t = g - 1
            if not 0 <= t <= 255:
                raise ValueError(f"o (hang {r}, cot {c}) co tile id {t}, ngoai khoang 0..255")
            r2 = r - cao_o(t) + 1
            if r2 < 0:
                tran.append((t, r, c, r2))
                continue
            if ra[r2 * w + c]:
                va_cham.append((t, r, c, ra[r2 * w + c], r2))
                continue
            ra[r2 * w + c] = t
    return ra, va_cham, tran
