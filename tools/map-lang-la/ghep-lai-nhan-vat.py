# -*- coding: utf-8 -*-
"""Ghep lai 3 manh nhan vat cho dung: het HAI KHUON MAT, hien duoc CHAN.

VAN DE. Bo 42 sprite Naruto khong chia theo cach cua Ngoc Rong:
  - 17 anh "than" (27044-27060) da ve san CA CAI DAU (mat nho, khong bang tran)
  - 3 anh "dau"  (27041-27043) lai la mot cai dau nua, TO HON, co bang tran
Engine ve ca hai -> tren man hinh co hai con mat o hai do cao, hai cai tai.
Cung luc do, 14 anh "chan" (27061-27074) bi anh than ve de len kin -> quan cut
ngang dui, khong thay dep dau.

CACH CHUA (khong can ve lai anh):
  1. Khong ve part dau nua — tro 3 khung cua part 30 sang mot anh TRONG SUOT.
     Anh than da co dau roi nen nhan vat van du mat.
  2. Xep lai chieu cao: chan cham dat, than nam TREN chan chu khong trum len.

     chan:  cao 32 don vi, muon day chan = cy  ->  dy = -32 + off(10) = -22
     than:  cao 40 don vi, muon day than = cy - 12 (lo 12 don vi chan duoi)
                                            ->  dy = -52 + off(17) = -35

     Than ve SAU chan nen no de len phan hong/tay cua anh chan — dung y do
     cua bo sprite (anh chan co san hai ban tay o phia tren).
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
# Do THUC NGHIEM trong game (khong tin cong thuc): voi than -35 / chan -22,
# duong mat co cat ngang HONG nhan vat — chan chim 11 don vi. Ha ca cum
# xuong them 11 de day dep trung mep co.
DY_MOI = {30: None, 57: -46, 58: -33}   # None = part dau, chi doi id

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
            if pid == 30:
                d[vi_id:vi_id + 2] = struct.pack(">h", ID_TRONG)   # dau -> anh trong suot
                bao.setdefault(pid, []).append("id->trong")
            else:
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
    if pid == 30:
        khoi = re.sub(r'"sprite_id":\s*\d+', f'"sprite_id": {ID_TRONG}', khoi)
    else:
        khoi = re.sub(r'"dy":\s*-?\d+', f'"dy": {DY_MOI[pid]}', khoi)
    van_m = van_m[:m.start()] + khoi + van_m[m.end():]
json.loads(van_m)
open(MAPPING, "w", encoding="utf-8", newline="").write(van_m)
print("da sua nguon mapping-id.json")
