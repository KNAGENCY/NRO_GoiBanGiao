# -*- coding: utf-8 -*-
"""Nang 3 manh nhan vat (dau/than/chan) len cho chan cham dat, het lun nua nguoi.

Vi sao lun: sprite Naruto moi cao 32-40 don vi, con bang tu the CharInfo trong
client la bang GOC cua Ngoc Rong, tinh cho sprite nho hon. Ket qua: chan ve tu
y = cy-10, cao 32 -> day chan o cy+22, tuc CHIM 22 don vi duoi mat dat.

Sua o dung cho: dy cua tung frame trong part.bin (offset rieng cua sprite),
KHONG dung vao CharInfo (bang tu the dung chung cho moi nhan vat, moi tu the).

Dinh dang part.bin (doc tu client GameScr.readPart):
    short soPart
    moi part:  byte type          (0 -> 3 frame, 1 -> 17 frame, 2 -> 14 frame)
               moi frame: short id, byte dx, byte dy
"""
import hashlib
import json
import os
import re
import shutil
import struct

BIN = "D:/NRO_GoiBanGiao/server/internal/game/assets/part.bin"
LOCK = "D:/NRO_GoiBanGiao/server/internal/game/assets/asset-lock.json"
MAPPING = "D:/NRO_GoiBanGiao/docs/my-thuat/mapping-id.json"

SO_FRAME = {0: 3, 1: 17, 2: 14}
PART_NHAN_VAT = {30, 57, 58}     # dau, than, chan cua nhan vat nguoi choi
NANG = 22                        # do lun do duoc trong game (chan chim 22 don vi)

goc = BIN + ".bak-goc"
if not os.path.exists(goc):
    shutil.copy2(BIN, goc)
d = bytearray(open(goc, "rb").read())      # LUON tu ban goc, chay lai khong cong don

so_part = struct.unpack(">h", d[0:2])[0]
i = 2
doi = 0
bao = []
for pid in range(so_part):
    typ = d[i]; i += 1
    n = SO_FRAME.get(typ)
    if n is None:
        raise SystemExit(f"part {pid} co type {typ} la, khong biet so frame — DUNG LAI")
    for f in range(n):
        vi_tri_dy = i + 2 + 1                       # short id (2) + byte dx (1)
        if pid in PART_NHAN_VAT:
            cu = struct.unpack("b", d[vi_tri_dy:vi_tri_dy + 1])[0]
            moi = cu - NANG
            if not -128 <= moi <= 127:
                raise SystemExit(f"part {pid} frame {f}: dy {cu} - {NANG} = {moi} tran byte — DUNG LAI")
            d[vi_tri_dy] = struct.pack("b", moi)[0]
            doi += 1
            if f < 2:
                bao.append(f"part {pid} frame {f}: dy {cu} -> {moi}")
        i += 4                                      # id(2) + dx(1) + dy(1)

if i != len(d):
    raise SystemExit(f"doc het {i}/{len(d)} byte — hieu sai dinh dang, DUNG LAI")
print(f"doc dung dinh dang: {so_part} part, {len(d)} byte")
for b in bao:
    print("  " + b)
print(f"  ... tong {doi} frame duoc nang {NANG} don vi")

open(BIN, "wb").write(bytes(d))
sha = hashlib.sha256(bytes(d)).hexdigest()
print(f"\npart.bin: {len(d)} byte, sha256 = {sha}")

# --- va asset-lock (hai cho, co the dang ghi hai gia tri khac nhau) ---
van = open(LOCK, encoding="utf-8").read()
cu_json = json.loads(van)
muc_moi = '{"bytes":%d,"path":"server/part.bin","sha256":"%s"}' % (len(d), sha)
da = set()
for ten, ds in (("outputs", cu_json["outputs"]), ("provenance.artifacts", cu_json["provenance"]["artifacts"])):
    x = next((k for k in ds if k.get("path") == "server/part.bin"), None)
    if x is None or (x["bytes"], x["sha256"]) in da:
        continue
    mau = re.compile(r'\{\s*"bytes":\s*' + str(x["bytes"]) +
                     r',\s*"path":\s*"server/part\.bin",\s*"sha256":\s*"' + x["sha256"] + r'"\s*\}')
    n = len(mau.findall(van))
    if n == 0:
        raise SystemExit(f"khong tim thay muc {ten} — DUNG LAI")
    van = mau.sub(muc_moi, van)
    da.add((x["bytes"], x["sha256"]))
    print(f"  {ten}: {x['sha256'][:12]}... -> {sha[:12]}...  ({n} cho)")
kt = json.loads(van)
for ten, ds in (("outputs", kt["outputs"]), ("provenance.artifacts", kt["provenance"]["artifacts"])):
    x = next((k for k in ds if k.get("path") == "server/part.bin"), None)
    if x:
        assert x["sha256"] == sha and x["bytes"] == len(d), f"va hong o {ten}"
assert len(kt["outputs"]) == len(cu_json["outputs"])
open(LOCK, "w", encoding="utf-8", newline="").write(van)
print("da va asset-lock.json")

# --- sua NGUON mapping-id.json cho khop, de build lai khong quay nguoc ---
mp = json.load(open(MAPPING, encoding="utf-8"))
sua = 0
for p in mp["parts"]:
    if p.get("id") in PART_NHAN_VAT:
        for fr in p["frames"]:
            fr["dy"] = fr.get("dy", 0) - NANG
            sua += 1
with open(MAPPING, "w", encoding="utf-8", newline="") as f:
    json.dump(mp, f, ensure_ascii=False, indent=2)
    f.write("\n")
print(f"da sua nguon mapping-id.json: {sua} frame")
