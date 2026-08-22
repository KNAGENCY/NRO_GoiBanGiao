# -*- coding: utf-8 -*-
"""Sinh lai map-1-collision.bin TU CHINH bo cuc client, khong go tay.

Doc thang res/mymap/1.bytes roi suy ra o nao dac. Nho vay doi bo cuc xong chay
lai la va cham tu khop — khong bao gio lech hang nua.

Dinh dang file: byte 0..3 = W, H, tileW, tileH; sau do W*H o, 1 = dac.
"""
import hashlib
import json
import os
import re
import shutil

MAP = "D:/NRO_Project/Assets/Resources/res/mymap/1.bytes"
BIN = "D:/NRO_GoiBanGiao/server/internal/game/assets/map-1-collision.bin"
LOCK = "D:/NRO_GoiBanGiao/server/internal/game/assets/asset-lock.json"

TILE_NEN = range(60, 70)     # 60-69 la tile nen (dat/co/da); >= 70 la vat the trang tri

d = open(MAP, "rb").read()
W, H = d[0], d[1]
o = list(d[2:2 + W * H])
print(f"bo cuc client: {W}x{H}")

cells = bytearray(W * H)
for c in range(W):
    for r in range(H):
        v = o[r * W + c]
        if v in TILE_NEN:
            cells[r * W + c] = 1
# Hai tuong bien: client de trong nhung khong co tuong thi nhan vat chay tuot ra ngoai.
for r in range(H):
    cells[r * W] = 1
    cells[r * W + (W - 1)] = 1

# In lai de mat nguoi kiem duoc
print("     " + "".join(str(i % 10) for i in range(W)))
for r in range(14, H):
    print(f"{r:3d}  " + "".join("X" if cells[r * W + c] else "." for c in range(W)))

data = bytes([W, H, 24, 24]) + bytes(cells)
if not os.path.exists(BIN + ".bak"):
    shutil.copy2(BIN, BIN + ".bak")
cu = open(BIN, "rb").read()
open(BIN, "wb").write(data)
sha = hashlib.sha256(data).hexdigest()
print(f"\nfile: {len(cu)} -> {len(data)} byte, sha256 = {sha}")

# Va asset-lock. File nay duoc ghi o HAI cho: outputs[] va provenance.artifacts[].
#
# LUU Y: hai cho co the dang ghi HAI GIA TRI KHAC NHAU. Dot Lat 4 chi cap nhat
# outputs, de provenance con ghi ban cu 89 byte (luoi 5x17) — day la "no manifest"
# ghi trong NHAT-KY-AGENT. Nen phai va TUNG CHO theo dung gia tri no dang co,
# khong dung chung mot mau.
van = open(LOCK, encoding="utf-8").read()
cu_json = json.loads(van)
moi_muc = '{"bytes":%d,"path":"server/map-1-collision.bin","sha256":"%s"}' % (len(data), sha)

da_va = set()
for ten, ds in (("outputs", cu_json["outputs"]), ("provenance.artifacts", cu_json["provenance"]["artifacts"])):
    x = next((i for i in ds if i.get("path") == "server/map-1-collision.bin"), None)
    if x is None:
        print(f"  ({ten}: khong ghi file nay, bo qua)")
        continue
    van_truoc = (x["bytes"], x["sha256"])
    if van_truoc in da_va:
        continue                       # hai cho cung gia tri -> mot lan sub la du
    mau = re.compile(
        r'\{\s*"bytes":\s*' + str(x["bytes"]) +
        r',\s*"path":\s*"server/map-1-collision\.bin",\s*"sha256":\s*"' + x["sha256"] + r'"\s*\}')
    n = len(mau.findall(van))
    if n == 0:
        raise SystemExit(f"khong tim thay muc {ten} trong van ban — DUNG LAI")
    van = mau.sub(moi_muc, van)
    da_va.add(van_truoc)
    print(f"  {ten}: {x['bytes']} byte / {x['sha256'][:12]}... -> {len(data)} byte / {sha[:12]}...  ({n} cho)")

van_moi = van
kt = json.loads(van_moi)
for ten, ds in (("outputs", kt["outputs"]), ("provenance.artifacts", kt["provenance"]["artifacts"])):
    x = next((i for i in ds if i.get("path") == "server/map-1-collision.bin"), None)
    if x is None:
        continue
    assert x["sha256"] == sha and x["bytes"] == len(data), f"va hong o {ten}"
assert len(kt["outputs"]) == len(cu_json["outputs"]), "so muc outputs thay doi — DUNG LAI"
open(LOCK, "w", encoding="utf-8", newline="").write(van_moi)
print(f"da va asset-lock.json ({len(kt['outputs'])} muc giu nguyen)")
