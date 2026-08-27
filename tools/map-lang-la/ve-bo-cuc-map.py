# -*- coding: utf-8 -*-
"""DA NGHI HUU — day la CACH CU de xep bo cuc map 1. Doc thi duoc, chay thi khong.

    Chay nham file nay = mat toan bo bo cuc dang co.
    Vi vay no CHAN NGAY o dong dau, phai co co --toi-that-su-muon-ve-lai-tu-dau
    moi chiu chay.

VI SAO NGHI HUU
  93 con so cung nam trong 167 dong code. Muon dich mot can nha sang trai 2 o
  thi sua mot con so roi chay lai — ma chay lai la XOA SACH 1248 o roi ve lai
  tu dau. Do khong phai la SUA, do la SINH LAI: moi lan sinh lai deu doi vai
  thu ngoai y muon. Chinh cai do lam ra cam giac "chap va, xep hinh".

  Thay bang cau noi Tiled (tools/map-lang-la/tiled/): bo cuc song trong mot file
  .tmx sua duoc bang chuot, sua tung o mot, khong xoa gi ca.

VI SAO KHONG XOA HAN FILE NAY
  No la HO SO GOC cua bo cuc map 1 dang chay: doc no la biet vi sao tung cong
  trinh nam o dung cot do, va nhung ghi chu trong no (nen cao, tuong bien hai
  cot, loi lat da phai nam cho phang) la kien thuc that su ton cong moi co.
"""
import shutil
import sys
import time

MAP = "D:/NRO_Project/Assets/Resources/res/mymap/1.bytes"
CO_EP_CHAY = "--toi-that-su-muon-ve-lai-tu-dau"

CANH_BAO = f"""
================================================================================
  DUNG LAI. FILE NAY DA NGHI HUU - VA KHONG CO GI BI THAY DOI CA.
================================================================================

Ban vua chay  ve-bo-cuc-map.py  - cach CU de xep bo cuc map 1.

  Moi lan chay, no XOA SACH 1248 o cua map roi ve lai tu dau bang 93 con so
  cung viet trong code. Nghia la: chay nham mot cai la bay mat toan bo bo cuc
  dang co. Lan nay chua dong vao file map - no dung lai truoc khi kip lam gi.

GIO XEP BO CUC BANG TILED - sua bang chuot, tung o mot, khong xoa gi:

  1) Nhin bo cuc dang co (0,6 giay, khong can mo Unity):
       python tools/map-lang-la/tiled/xem-truoc.py --khung

  2) Dua map dang co ra file sua duoc:
       python tools/map-lang-la/tiled/map-sang-tiled.py

  3) Mo  tools/map-lang-la/tiled/map-1-lang-la.tmx  bang Tiled (mapeditor.org,
     mien phi). Keo tha nha cua cho vua y - co san anh concept lam giay can
     nam duoi de doi chieu. Luu lai.

  4) Dua nguoc ve map game, roi sinh lai va cham:
       python tools/map-lang-la/tiled/tiled-sang-map.py
       python tools/map-lang-la/sinh-va-cham.py

  Huong dan day du:  tools/map-lang-la/tiled/README.md

FILE NAY VAN GIU LAI de doc: no ghi lai vi sao tung cong trinh nam o dung cot
do. Doc thi duoc, chay thi khong.

Neu ban THAT SU muon vut bo cuc hien tai di va ve lai tu dau (rat hiem khi
dung - gan nhu luc nao Tiled cung la thu ban muon), thi go them co nay:

    python tools/map-lang-la/ve-bo-cuc-map.py {CO_EP_CHAY}

================================================================================
"""

if CO_EP_CHAY not in sys.argv:
    print(CANH_BAO)
    # Thoat NGAY, truoc ca dong shutil.copy2 ben duoi: khong doc, khong ghi,
    # khong sao luu — chay nham thi khong mat gi het.
    raise SystemExit(1)

# Sao luu co DAU THOI GIAN. Ban cu ghi de len dung mot ten .bak-truoc-ve-lai,
# nen chay hai lan la ban luu "truoc khi ve lai" bi chinh lan chay thu hai de
# len — mat luon cai bo cuc minh dinh cuu.
ban_luu = f"{MAP}.bak-truoc-ve-lai-{time.strftime('%Y%m%d-%H%M%S')}"
shutil.copy2(MAP, ban_luu)
print(f"[sao luu] {ban_luu}")
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
    if 16 <= c <= 28:     CAO[c] = 17     # nen cao cua khu trung tam (torii + nha lon)
    elif False:           CAO[c] = 17     # (go thu hai da gop vao khu trung tam)
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
    # Mat cat o cho doi cao. Concept ve lang dung tren KE DA XEP, khong phai
    # vach dat tro — nay da co anh ke da (84) va chan go (83) nen dung dung.
    if trai > mat:
        for r in range(mat+1, trai+1):
            dat(r, c, 84 if r == mat+1 else 83)      # ke da o tren, chan go o duoi
    if phai > mat:
        for r in range(mat+1, phai+1):
            dat(r, c, 84 if r == mat+1 else 83)

# ---- 3) VAT THE: dat theo do cao mat dat tung cho ----
#   (r_chan = hang mat dat; vat cao N o -> dat o hang r_chan - N)
def dat_vat(cot, slot, cao_o):
    r_chan = CAO.get(cot, 18)
    dat(r_chan - cao_o, cot, slot)

# BO CUC — XEP DAY DAC nhu CONCEPT DA DUYET.
#
# Do lai concept va game moi thay minh da di lac: o dai nha cua, concept co 0%
# be ngang la troi (nha che kin, lang la MOT KHOI LIEN), con game co 39-55% —
# cac cong trinh dung roi rac tren nen troi. Do la cai chu du an goi la
# "nhu dang choi xep hinh".
#
# Vong truoc con bo bot nha va gian vat the ra "cho khoi cham nhau" — dung huong
# ky thuat nhung SAI huong bo cuc, vi trong concept nha CHONG LAN nhau la binh
# thuong va chinh cho do tao chieu sau.
#
# Nay xep sat lai, cho chong 1-2 o, dung lai anh nha (concept cung co nhieu nha
# giong nhau), phu kin 48/48 o dung duoc.
#
#          cot      vat            ghi chu
dat_vat(2,  70, 6)   #  2-7    nha nho
dat_vat(8,  71, 6)   #  8-15   nha vua
dat_vat(16, 74, 5)   # 16-20   cong torii   (tren nen cao)
dat_vat(21, 72, 8)   # 21-28   nha lon      (tren nen cao)
dat_vat(29, 77, 3)   # 29-31   gieng
dat_vat(32, 73, 6)   # 32-39   quan mi
dat_vat(40, 80, 8)   # 40-45   cay
dat_vat(44, 70, 6)   # 44-49   nha nho lan hai
# vat nho dat XEN vao chan cong trinh — trong concept den long va hang rao cung
# nam de len chan nha, do la mot phan cua ve day dac
dat_vat(7,  76, 4)   # den long truoc nha vua
dat_vat(31, 76, 4)   # den long truoc quan
for c in (15, 30, 43):
    dat_vat(c, 75, 2)                      # hang rao xen giua cac nha
for c in (6, 20, 28, 39, 49):
    dat_vat(c, 81, 1)                      # bui co o chan cong trinh
dat_vat(19, 78, 2)   # thung go canh torii

# BAC THANG len nen cao (concept co cau thang da dan len cong torii).
# Dat o mep TRAI cua nen cao: mot bac o chan, mot bac dinh noi vao mat co.
dat(CAO.get(15, 18), 15, 85)          # bac thang o chan nen
dat(CAO.get(16, 17), 16, 86)          # bac dinh, noi vao mat nen cao

# Hai cay nho + bui ram moi, rai vao cho con thua de bot troi
dat_vat(13, 87, 4)   # cay nho (4 o)
dat_vat(38, 88, 2)   # bui ram (2 o)

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
