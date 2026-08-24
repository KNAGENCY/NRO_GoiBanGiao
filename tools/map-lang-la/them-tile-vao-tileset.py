# -*- coding: utf-8 -*-
"""Them tile nen Lang La (60-69) vao bang tile_set_info.

Vi sao can: bang goc cua Ngoc Rong chi khai bao tile 1-39. Tile Lang La danh so
60-82 nen client dat types[] = 0 cho ca man -> nhan vat roi xuyen dat xuong day.

Dinh dang file (doc tu client Controller.cs case -82):
    byte  soBo
      byte  soNhom
        int32 loai (big-endian)
        byte  soTile
          byte tile ...

Chi sua BO 1 (map 1 dung tileID = 1). Cac bo khac giu nguyen tung byte.
"""
import struct
import shutil
import hashlib
import json
import os

FILE = "D:/NRO_GoiBanGiao/server/internal/game/assets/tile_set_info.bin"
LOCK = "D:/NRO_GoiBanGiao/server/internal/game/assets/asset-lock.json"

BIT_SAN = 0x0002    # dung duoc len tren
BIT_PHAI = 0x0004   # chan ben phai
BIT_TRAI = 0x0008   # chan ben trai
BIT_TRAN = 0x2000   # chan tu duoi nhay len

# Tile nen Lang La va vai tro
# 60-69 la bo tile nen ban dau; 83-86 la bo bo sung dot 2 (chan go, ke da,
# bac thang, bac thang dinh) — deu la mat dat, deu dung duoc len tren.
MAT_DAT = [60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 83, 84, 85, 86]

# KHOI_DAC: chan di ngang. Phai co ca 62 (mat co) — day la cho tung sai:
# nhan vat dung tren mat dat hang 18 thi client kiem chan ngang o cy - chh,
# tuc HANG 17. O go, hang 17 la MAT CO cua go. Neu 62 khong chan ngang thi
# nguoi choi di thang xuyen vao suon go, trong khi may chu (doc tu bo cuc) lai
# chan -> goi di chuyen bi tu choi, nguoi khac thay dung chet mot cho.
#
# Cho 62 chan ngang KHONG lam kho di lai o cho phang: dung tren hang 18 thi o
# ngang tam than la hang 17, va hang 17 o cho phang la khong khi.
KHOI_DAC = [60, 61, 62, 65, 66, 67, 68, 83, 84]      # dat than + mat co + vach go + ke da
CO_TRAN = [60, 61, 67, 68, 83, 84]                   # khoi dac: khong nhay xuyen len

THEM = {BIT_SAN: MAT_DAT, BIT_PHAI: KHOI_DAC, BIT_TRAI: KHOI_DAC, BIT_TRAN: CO_TRAN}


def doc(data):
    """Bung file thanh [ {loai: [tile,...]}, ... ] giu nguyen THU TU nhom."""
    i = 0
    so_bo = data[i]; i += 1
    bo = []
    for _ in range(so_bo):
        so_nhom = data[i]; i += 1
        nhom = []
        for _ in range(so_nhom):
            loai = struct.unpack(">i", data[i:i + 4])[0]; i += 4
            n = data[i]; i += 1
            nhom.append((loai, list(data[i:i + n]))); i += n
        bo.append(nhom)
    assert i == len(data), f"thua {len(data) - i} byte cuoi"
    return bo


def ghi(bo):
    out = bytearray([len(bo)])
    for nhom in bo:
        out.append(len(nhom))
        for loai, tiles in nhom:
            out += struct.pack(">i", loai)
            assert len(tiles) <= 255, f"nhom 0x{loai:04X} co {len(tiles)} tile, vuot 255"
            out.append(len(tiles))
            out += bytes(tiles)
    return bytes(out)


if not os.path.exists(FILE + ".bak-goc"):
    shutil.copy2(FILE, FILE + ".bak-goc")
goc = open(FILE + ".bak-goc", "rb").read()   # LUON tu ban goc: chay lai khong cong don

bo = doc(goc)
assert ghi(bo) == goc, "doc/ghi khong khop — dinh dang hieu sai, DUNG LAI"
print(f"doc lai byte-dung: {len(goc)} byte, {len(bo)} bo tile")

nhom1 = bo[0]
co_san = {loai for loai, _ in nhom1}
for loai, ds_them in THEM.items():
    if loai not in co_san:                       # nhom chua ton tai -> them moi
        nhom1.append((loai, list(ds_them)))
        print(f"  nhom 0x{loai:04X}: TAO MOI voi {len(ds_them)} tile")
        continue
    for k, (l, tiles) in enumerate(nhom1):
        if l != loai:
            continue
        moi = [t for t in ds_them if t not in tiles]
        nhom1[k] = (l, tiles + moi)
        print(f"  nhom 0x{loai:04X}: {len(tiles)} -> {len(tiles) + len(moi)} tile (them {moi})")
        break

moi = ghi(bo)
open(FILE, "wb").write(moi)
sha = hashlib.sha256(moi).hexdigest()
print(f"\nfile moi: {len(moi)} byte (goc {len(goc)}), sha256 = {sha}")

# Cap nhat asset-lock cho khop, neu khong server tu choi khoi dong.
#
# TUYET DOI khong json.dump lai ca file: dinh dang se doi tren HANG NGHIN dong
# va lam vo cac test doi chieu chuoi nguyen van. Chi va DUNG MOT muc bang regex.
import re

van = open(LOCK, encoding="utf-8").read()
cu = json.loads(van)
truoc = next(o for o in cu["outputs"] if o.get("path") == "server/tile_set_info.bin")
print(f"asset-lock: {truoc['bytes']} byte / {truoc['sha256'][:16]}...  ->  {len(moi)} byte / {sha[:16]}...")

# Muc can va, giu nguyen thu tu khoa nhu file dang co
mau = re.compile(
    r'\{\s*"bytes":\s*' + str(truoc["bytes"]) +
    r',\s*"path":\s*"server/tile_set_info\.bin",\s*"sha256":\s*"' + truoc["sha256"] + r'"\s*\}'
)
# File .bin duoc ghi o HAI cho trong lock: outputs[] va provenance.artifacts[].
# Hai cho phai khop nhau, nen va ca hai.
khop = mau.findall(van)
if len(khop) != 2:
    raise SystemExit(f"tim thay {len(khop)} muc khop, can dung 2 — DUNG LAI de khong pha file")
van_moi = mau.sub(
    '{"bytes":%d,"path":"server/tile_set_info.bin","sha256":"%s"}' % (len(moi), sha),
    van)

kt = json.loads(van_moi)          # phai van la JSON hop le
for cho, ds in (("outputs", kt["outputs"]), ("provenance.artifacts", kt["provenance"]["artifacts"])):
    sau = next(o for o in ds if o.get("path") == "server/tile_set_info.bin")
    assert sau["sha256"] == sha and sau["bytes"] == len(moi), f"va hong o {cho}"
assert len(kt["outputs"]) == len(cu["outputs"]), "so muc outputs thay doi — DUNG LAI"
assert len(kt["provenance"]["artifacts"]) == len(cu["provenance"]["artifacts"]), "so muc artifacts thay doi"
open(LOCK, "w", encoding="utf-8", newline="").write(van_moi)
print(f"da va asset-lock.json — {len(van_moi) - len(van):+d} ky tu, {len(kt['outputs'])} muc giu nguyen")
