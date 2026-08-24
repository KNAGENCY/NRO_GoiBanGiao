# -*- coding: utf-8 -*-
"""Ghep lai 3 manh nhan vat cho dung: het HAI KHUON MAT, hien duoc CHAN.

VAN DE. Bo 42 sprite Naruto khong chia theo cach cua Ngoc Rong:
  - 17 anh "than" (27044-27060) da ve san CA CAI DAU (mat nho, khong bang tran)
  - 3 anh "dau"  (27041-27043) lai la mot cai dau nua, TO HON, co bang tran
Engine ve ca hai -> tren man hinh co hai con mat o hai do cao, hai cai tai.
Cung luc do, 14 anh "chan" (27061-27074) bi anh than ve de len kin -> quan cut
ngang dui, khong thay dep dau.

DA THU MOT CACH SAI, ghi lai de khong lam lai: tro part dau sang anh TRONG SUOT.
Het hai mat that, nhung anh than chi ve NUA DUOI cai dau (cat phang ngay tren
chan may) — xoa manh dau la xoa luon so, toc gai va BANG TRAN LA KONOHA, tuc
mat dau hieu nhan dien cua ninja Lang La. Nhin anh la thay dinh dau bi vat phang.

CACH CHUA DUNG (khong can ve lai anh): giu manh dau, nhung dat cho no PHU DUNG
len phan dau da ve san trong anh than.

  Bang tu the CharInfo[0] (Char.cs:483) cho tuong quan GOC giua ba manh:
      dau  offset 34   |  than offset 16  |  chan offset 10
  tuc dau cao hon than 18 don vi, than cao hon chan 6. Anh ve tai
  y0 = cy - offset + dy, nen giu nguyen tuong quan thi ba manh van ghep dung —
  chi can dich CA CUM cung mot luong.

  1. Dich ca cum len 11 don vi de day chan trung mep co (do trong game, xem
     LOI-DA-GAP.md §7 — co mot do lech he thong ~11-12 don vi giua toa do logic
     va toa do ve, chua truy ra goc).
  2. Rieng manh dau lech them: cho CAM cua anh dau trung CAM cua dau-trong-than.
     Tinh ra dau dy = -24 trong khi than/chan = -33.

  Than ve SAU chan nen de len phan hong/tay cua anh chan — dung y do cua bo
  sprite (anh chan co san hai ban tay o phia tren).
"""
import hashlib
import json
import os
import re
import shutil
import struct

import numpy as np
from PIL import Image

BIN = "D:/NRO_GoiBanGiao/server/internal/game/assets/part.bin"
LOCK = "D:/NRO_GoiBanGiao/server/internal/game/assets/asset-lock.json"
MAPPING = "D:/NRO_GoiBanGiao/docs/my-thuat/mapping-id.json"

ID_TRONG = 27083          # id moi cho anh trong suot (chua ai dung)
SO_FRAME = {0: 3, 1: 17, 2: 14}
# Tuong quan GOC (dau 34 / than 16 / chan 10) chi dung voi sprite Ngoc Rong, vi
# anh than cua ho TRONG SUOT o phan duoi nen chan lo ra. Anh than moi (27044-60)
# ve DAC toi tan day anh, giu tuong quan goc la than che kin chan.
#
# Nen tinh lai tu cho hai anh PHAI trung nhau — hai ban tay ve o ca anh than lan
# anh chan:
#     tay trong than: hang 60-79 cua anh 80  -> 30..40 don vi tu dinh anh
#     tay trong chan: hang 14-30 cua anh 64  ->  7..15 don vi tu dinh anh
#     cho trung:  (cy-16+dy_than) + 30 = (cy-10+dy_chan) + 7
#              => dy_than = dy_chan - 17
#
# chan dy=-33 (day chan trung mep co, ke ca do lech he thong ~11 don vi)
# than dy=-50 (tay trung, va day than o cy-26 nen chan lo ra 15 don vi: du dep
#              va ong chan)
# dau  dy=-41 (cam cua anh dau trung cam cua dau-trong-than)
DY_MOI = {30: -41, 57: -50, 58: -33}

# --- 1) tao anh trong suot o ca 4 muc zoom ---
for muc, canh in (("x1", 32), ("x2", 64), ("x3", 96), ("x4", 128)):
    d = f"D:/NRO_Project/Assets/Resources/res/{muc}/smallimage/"
    os.makedirs(d, exist_ok=True)
    Image.fromarray(np.zeros((canh, canh, 4), dtype=np.uint8), "RGBA").save(d + f"Small{ID_TRONG}.png")
print(f"da tao anh trong suot Small{ID_TRONG}.png o 4 muc zoom")

# --- 2) sua part.bin ---
goc = BIN + ".bak-ghep"
if not os.path.exists(goc):
    shutil.copy2(BIN, goc)
d = bytearray(open(BIN, "rb").read())        # tu ban HIEN TAI (da co dy -22)

so_part = struct.unpack(">h", d[0:2])[0]
i = 2
bao = {}
for pid in range(so_part):
    typ = d[i]; i += 1
    n = SO_FRAME.get(typ)
    if n is None:
        raise SystemExit(f"part {pid} type {typ} la — DUNG LAI")
    for f in range(n):
        vi_id, vi_dy = i, i + 3
        if pid in DY_MOI:
            moi = DY_MOI[pid]
            d[vi_dy] = struct.pack("b", moi)[0]
            bao.setdefault(pid, []).append(f"dy->{moi}")
        i += 4
if i != len(d):
    raise SystemExit(f"doc het {i}/{len(d)} byte — DUNG LAI")
for pid, ds in bao.items():
    print(f"  part {pid}: {len(ds)} khung, {ds[0]}")

open(BIN, "wb").write(bytes(d))
sha = hashlib.sha256(bytes(d)).hexdigest()
print(f"\npart.bin {len(d)} byte, sha256 = {sha}")

# --- 3) va asset-lock (hai cho) ---
van = open(LOCK, encoding="utf-8").read()
cu = json.loads(van)
muc_moi = '{"bytes":%d,"path":"server/part.bin","sha256":"%s"}' % (len(d), sha)
da = set()
for ten, ds in (("outputs", cu["outputs"]), ("provenance.artifacts", cu["provenance"]["artifacts"])):
    x = next((k for k in ds if k.get("path") == "server/part.bin"), None)
    if x is None or (x["bytes"], x["sha256"]) in da:
        continue
    mau = re.compile(r'\{\s*"bytes":\s*' + str(x["bytes"]) +
                     r',\s*"path":\s*"server/part\.bin",\s*"sha256":\s*"' + x["sha256"] + r'"\s*\}')
    if not mau.findall(van):
        raise SystemExit(f"khong tim thay muc {ten} — DUNG LAI")
    van = mau.sub(muc_moi, van)
    da.add((x["bytes"], x["sha256"]))
    print(f"  {ten}: {x['sha256'][:12]}... -> {sha[:12]}...")
kt = json.loads(van)
for ds in (kt["outputs"], kt["provenance"]["artifacts"]):
    x = next((k for k in ds if k.get("path") == "server/part.bin"), None)
    if x:
        assert x["sha256"] == sha and x["bytes"] == len(d)
open(LOCK, "w", encoding="utf-8", newline="").write(van)
print("da va asset-lock.json")

# --- 4) sua nguon mapping-id.json cho khop (va dung 34 dong) ---
van_m = open(MAPPING, encoding="utf-8").read()
for pid in (30, 57, 58):
    m = re.search(r'\{\s*"id":\s*' + str(pid) + r',\s*"type":\s*\d+,\s*"frames":\s*\[.*?\n      \]\n    \}',
                  van_m, re.S)
    if not m:
        raise SystemExit(f"khong tim thay khoi part {pid} trong mapping-id — DUNG LAI")
    khoi = m.group(0)
    khoi = re.sub(r'"dy":\s*-?\d+', f'"dy": {DY_MOI[pid]}', khoi)
    van_m = van_m[:m.start()] + khoi + van_m[m.end():]
json.loads(van_m)
open(MAPPING, "w", encoding="utf-8", newline="").write(van_m)
print("da sua nguon mapping-id.json")
