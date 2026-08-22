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

# ---- 1) DON SACH: xoa TOAN BO o (ke ca vat the cu cua cac lat truoc) ----
#      Chi xoa tile < 60 la KHONG DU: vat the cu (>= 70) se con nguyen va
#      chong len vat the moi -> nha long nha, quan long quan.
don = sum(1 for r in range(tmh) for c in range(tmw) if lay(r, c) != 0)
for r in range(tmh):
    for c in range(tmw):
        dat(r, c, 0)
print(f"[don] xoa sach {don} o (ca tile nen lan vat the cu)")

# ---- 2) DIA HINH: duong chan co GO va DOC, khong phang li ----
# do cao mat dat theo cot: 18 = chuan; 17 = go cao; 19 = trung
CAO = {}
for c in range(tmw):
    if 23 <= c <= 27:     CAO[c] = 17     # go duoi cong torii (cong dung tren be cao)
    elif 34 <= c <= 41:   CAO[c] = 17     # go duoi nha lon (khop dung cot 34-41)
    else:                 CAO[c] = 18
# lam muot: khong de bac nhay > 1 o
for c in range(1, tmw):
    if abs(CAO[c] - CAO[c-1]) > 1:
        CAO[c] = CAO[c-1] + (1 if CAO[c] > CAO[c-1] else -1)

def bam(c, r):
    x = (c*73856093) ^ (r*19349663)
    x = (x ^ (x >> 13)) * 1274126177
    return (x ^ (x >> 16)) & 0xFFFF

# Loi lat da phai nam o cho PHANG. De o 20-26 thi no chui xuong duoi go
# (go 23-27), chi con lo ra mot mang da o mat cat vach — trong rat la.
LOI_DI = set(range(8, 15))      # loi dat lat da truoc day nha, cho phang

for c in range(1, tmw-1):
    mat = CAO[c]
    trai, phai = CAO.get(c-1, mat), CAO.get(c+1, mat)
    # mat tren: co, hoac loi lat da
    if c in LOI_DI and mat == 18:
        # hai dau loi di dung o CHUYEN TIEP (63/64) de co lan dan sang da,
        # khong con la mot nhat cat doc phang li
        if c == min(LOI_DI):    dat(mat, c, 63)   # co -> da
        elif c == max(LOI_DI):  dat(mat, c, 64)   # da -> co
        else:                   dat(mat, c, 69)   # duong lat da
    elif trai > mat:          # bac vuong: van dung mat co
        dat(mat, c, 62)
    elif phai > mat:                          # ben phai thap hon -> goc phai
        dat(mat, c, 62)
    else:
        dat(mat, c, 62)                       # mat co thuong
    # than dat xuong day — cang sau cang toi, xen ke bien the cho khoi phang
    for r in range(mat+1, tmh):
        # ranh gioi tang GON SONG theo cot, khong ke thang bang thuoc
        song = (bam(c, 777) % 3) - 1        # -1, 0 hoac 1
        sau = (r - mat) + song
        if sau <= 2:    nen, phu = 60, 61   # sat mat: dat sang
        elif sau <= 4:  nen, phu = 67, 60   # tang giua
        else:           nen, phu = 68, 67   # tang day: toi nhat
        dat(r, c, phu if bam(c, r) % 100 < 34 else nen)
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

# BO CUC — moi cong trinh mot khoang rieng, khong cai nao cam vao cai nao.
#
# Map chi rong 52 o ma truoc do nhoi 4 cong trinh lon (6+8+8+8 = 30 o) cong
# torii va cay, nen cai nao cung dinh cai nao: torii cam vao nha vua, cay tran
# ra ngoai map, den long nam trong long nha nho, hai hang rao cam xuyen cot
# torii va dan phang len tuong nha lon.
#
# Bo NHA VUA de lay lai khoang tho, va dua CAY vao giua map — truoc do ca man
# chi co mot cay dinh sat bien phai, khuc giua khong co mang xanh nao.
#
#          cot      vat                 ghi chu
dat_vat(2,  70, 6)   # 2-7    nha nho
dat_vat(9,  78, 2)   # 9-10   thung go
dat_vat(12, 76, 4)   # 12-13  den long
dat_vat(15, 80, 8)   # 15-20  cay          (dua vao giua map)
dat_vat(23, 74, 5)   # 23-27  cong torii   (dung tren go 23-27)
dat_vat(29, 77, 3)   # 29-31  gieng
dat_vat(34, 72, 8)   # 34-41  nha lon      (dung tron tren go 34-41)
dat_vat(43, 73, 6)   # 43-48  quan mi
for c in [21, 32]:
    dat_vat(c, 75, 2)                      # hang rao lam lan can hai chan go
for c in [8, 14, 28, 42, 49]:
    dat_vat(c, 81, 1)                      # bui co o cac khe, khong cham cong trinh

# Tuong VO HINH o bien: khong co thi nguoi choi roi xuong ho o mep map.
#
# Dat HAI COT moi ben chu khong phai mot. Ly do: camera dung o cmxLim =
# (tmw-1)*24 - gW, tuc chi nhin duoc toi x = 1224; ma than nhan vat ve SANG PHAI
# diem neo (moc TOP|LEFT khi quay phai). Neu tuong chi o cot cuoi thi nhan vat
# van di duoc toi sat 1224 va truot HAN ra ngoai mep phai — nguoi choi mat dau
# chinh minh, chi con mot mau bong o mep.
#
# Tile 39 da duoc khai bao dac ca bon huong trong bang goc, anh cua no da lam
# trong suot nen tuong khong ve ra gi.
# Chi dat tuong o phan TREN mat dat. Neu dat tuong ca cot thi phan duoi khong
# co tile dat nao -> mep man hinh ho ra mot dai mau troi, nhin nhu map bi khuyet.
# Chan ngang cua client kiem o hang (mat dat - 1) nen tuong tu hang do tro len
# la du giu nguoi choi lai.
for c in (0, 1, tmw - 2, tmw - 1):
    mat = CAO.get(c, 18)
    for r in range(tmh):
        if r < mat:
            dat(r, c, 39)                                  # tuong vo hinh
        else:
            dat(r, c, 61 if bam(c, r) % 100 < 34 else 60)  # dat, de mep nhin lien mach

open(MAP, "wb").write(bytes(data))

# ---- 4) KIEM LAI ----
d2 = open(MAP, "rb").read()
cells = list(d2[2:2+tmw*tmh])
# Tile 39 la TUONG co chu dich (dat o buoc tren), khong phai tan du.
con_cu = sum(1 for v in cells if 0 < v < 60 and v != 39)
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
