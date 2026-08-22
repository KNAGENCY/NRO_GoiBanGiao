# -*- coding: utf-8 -*-
# VE LAI bo cuc map 1 TU DAU: don sach tan du tile Dragon Ball, dat lai bang tile Lang La.
import shutil

MAP = "D:/NRO_Project/Assets/Resources/res/mymap/1.bytes"
shutil.copy2(MAP, MAP + ".bak-truoc-ve-lai")
data = bytearray(open(MAP, "rb").read())
tmw, tmh = data[0], data[1]

def dat(r, c, v):
    if 0 <= r < tmh and 0 <= c < tmw:
        data[2 + r*tmw + c] = v

def lay(r, c):
    return data[2 + r*tmw + c] if (0 <= r < tmh and 0 <= c < tmw) else 0

# ---- 1) DON SACH: xoa MOI tile cu (< 60), giu tile moi ----
don = 0
for r in range(tmh):
    for c in range(tmw):
        if 0 < lay(r, c) < 60:
            dat(r, c, 0); don += 1
print(f"[don] xoa {don} tile cu (tan du ban do Dragon Ball)")

# ---- 2) DIA HINH: duong chan co GO va DOC, khong phang li ----
# do cao mat dat theo cot: 18 = chuan; 17 = go cao; 19 = trung
CAO = {}
for c in range(tmw):
    if 6 <= c <= 11:      CAO[c] = 17     # go cao ben trai
    elif 30 <= c <= 37:   CAO[c] = 17     # go cao ben phai
    else:                 CAO[c] = 18
# lam muot: khong de bac nhay > 1 o
for c in range(1, tmw):
    if abs(CAO[c] - CAO[c-1]) > 1:
        CAO[c] = CAO[c-1] + (1 if CAO[c] > CAO[c-1] else -1)

def bam(c, r):
    x = (c*73856093) ^ (r*19349663)
    x = (x ^ (x >> 13)) * 1274126177
    return (x ^ (x >> 16)) & 0xFFFF

LOI_DI = set(range(20, 27))     # loi dat lat da giua lang

for c in range(1, tmw-1):
    mat = CAO[c]
    trai, phai = CAO.get(c-1, mat), CAO.get(c+1, mat)
    # mat tren: co, hoac loi lat da
    if c in LOI_DI and mat == 18:
        dat(mat, c, 69)                       # duong lat da
    elif trai > mat:          # bac vuong: van dung mat co
        dat(mat, c, 62)
    elif phai > mat:                          # ben phai thap hon -> goc phai
        dat(mat, c, 62)
    else:
        dat(mat, c, 62)                       # mat co thuong
    # than dat xuong day
    for r in range(mat+1, tmh):
        dat(r, c, 61 if bam(c, r) % 100 < 16 else 60)
    # vach o cho doi cao (mat cat dat lo ra)
    if trai > mat:
        for r in range(mat+1, trai+1): dat(r, c, 65)
    if phai > mat:
        for r in range(mat+1, phai+1): dat(r, c, 66)

# ---- 3) VAT THE: dat theo do cao mat dat tung cho ----
#   (r_chan = hang mat dat; vat cao N o -> dat o hang r_chan - N)
def dat_vat(cot, slot, cao_o):
    r_chan = CAO.get(cot, 18)
    dat(r_chan - cao_o, cot, slot)

dat_vat(3,  70, 6)    # nha nho     (6 o)
dat_vat(13, 71, 6)    # nha vua     (6 o)
dat_vat(31, 72, 8)    # nha lon     (8 o)
dat_vat(40, 73, 6)    # nha quan    (6 o)
dat_vat(19, 74, 5)    # cong torii  (5 o)
dat_vat(28, 77, 3)    # gieng       (3 o)
dat_vat(9,  78, 2)    # thung go    (2 o)
dat_vat(45, 79, 3)    # bien hieu   (3 o)
dat_vat(47, 80, 8)    # cay         (8 o)
dat_vat(7,  76, 4)    # den long    (4 o)
for c in [17, 22, 35]:
    dat_vat(c, 75, 2)                       # hang rao (2 o)
for c in [5, 16, 26, 38, 44]:
    dat_vat(c, 81, 1)                       # bui co   (1 o)

open(MAP, "wb").write(bytes(data))

# ---- 4) KIEM LAI ----
d2 = open(MAP, "rb").read()
cells = list(d2[2:2+tmw*tmh])
con_cu = sum(1 for v in cells if 0 < v < 60)
print(f"[kiem] tile cu con sot: {con_cu} (phai = 0)")
print(f"[kiem] file {len(d2)} byte, luoi {d2[0]}x{d2[1]}")
print("\nSo do hang 9-23 (. trong | # dat/co moi | N vat the):")
print("     " + "".join(str(c % 10) for c in range(0, 52)))
for r in range(9, tmh):
    s = ""
    for c in range(tmw):
        v = cells[r*tmw + c]
        s += "." if v == 0 else ("N" if v >= 70 else "#")
    print(f"{r:3d}  {s}")
