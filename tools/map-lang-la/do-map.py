# -*- coding: utf-8 -*-
"""THUOC DO MAP — mot lenh, 1 giay, thay cho mot dot soi 45-80 phut.

    python tools/map-lang-la/do-map.py <anh-game.png> [anh-game-2.png ...]

    --khong-anh      chi in bang so, khong ve anh chan doan
    --ra <thu-muc>   noi de anh chan doan (mac dinh: tools/map-lang-la/chan-doan)

Ca ngay 22/08 minh chay bon dot soi bang agent, moi dot 45-80 phut, de tim ra
nhung con so ma may do duoc trong 1 giay. File nay gom dung nhung phep do da
that su bat duoc loi, moi phep kem MOC so voi concept da duyet.

DUNG NO TRUOC, agent SAU. Agent chi con dung o hai viec ma may khong lam duoc:
nhan xet tham my ("net ve choi voi nen") va doi chieu y do ("concept muon gi").

VI SAO CO THEM ANH CHAN DOAN
  Bang so noi "Dat phang li 91.4 !!" — biet la hong, nhung KHONG biet hong O
  DAU, nen van phai sua mo. Anh chan doan to mau thang len chinh anh dau vao:
  vung dat bi phang, cot troi lot xuong dai nha, hang pixel co duong ke ngang.
  Con so cho biet CO benh, anh cho biet benh NAM O DAU.

  Bang so van in y het nhu cu — anh chi la thu THEM. Muon nhanh toi da thi
  them --khong-anh.

BAY KHI DO — da mac ca ba, ghi lai de khong mac lai:
  1. Mat na "mau troi" bat ca lang xa mo (cung xanh nhat) -> so khong nhuc nhich
     du anh doi ro. Phai do bang DO CHI TIET, khong do bang mau.
     (Anh chan doan o o so 1 cho THAY tan mat cai bay nay: lang xa cung bi to
     xanh. Do la ly do chi so 3 phai do bang do chi tiet.)
  2. "Do lech chuan cua hang" phai tinh theo chieu ngang cho TUNG KENH mau.
     Tinh std tren ca 3 kenh thi mot hang dong mau (12,95,169) ra std ~64, doc
     nham thanh "con nhieu van" trong khi hang do da phang tuyet doi.
  3. Khong phai duong ngang nao cung la loi. Chi tinh la duong ke khi DA SO cot
     cung nhay (phan vi 25 cua cac cot). Vai cot nhay la ranh gioi tu nhien
     (dinh nui, mep mai). Da doc nham dinh nui thanh duong ke va di "sua" nham.
"""
import os
import sys

import numpy as np
from PIL import Image, ImageDraw, ImageFont

CONCEPT = "D:/NRO_GoiBanGiao/docs/my-thuat/map-1-lang-la/concept-map1-langla.png"
THU_MUC_RA = os.path.join(os.path.dirname(os.path.abspath(__file__)), "chan-doan")

# Moc lay tu concept da duyet. "tot" = huong can di.
MOC = {
    "troi_toan_man":  (24, "thap hon la tot"),
    "troi_dai_nha":   (0,  "thap hon la tot"),
    "o_tron_giua":    (6,  "thap hon la tot"),
    "dat_phang":      (27, "thap hon la tot"),
    "cao_khoi_nui":   (31, "cao hon la tot"),
    "ke_ngang":       (0,  "thap hon la tot"),
}

# Nguong tuyet doi trong so tay (§3). Khac voi MOC o tren: MOC la "so cua
# concept" (so sanh tuong doi), con day la "dat khi nao" cua so tay. Ghi ca hai
# len anh de nguoi doc khong phai mo so tay ra tra.
NGUONG_SO_TAY = {
    "troi_toan_man": "\u2264 30%",
    "troi_dai_nha":  "\u2264 3%",
    "o_tron_giua":   "\u2264 10%",
    "dat_phang":     "\u2264 40%",
    "cao_khoi_nui":  "\u2265 30% khung",
    "ke_ngang":      "< 10",
    "nen_lap":       "< 40",
}

TEN_DEP = {
    "troi_toan_man": "Troi chiem man",
    "troi_dai_nha":  "Dai nha cua la troi",
    "o_tron_giua":   "O tron o dai giua",
    "dat_phang":     "Dat phang li",
    "cao_khoi_nui":  "Khoi nui cao",
    "ke_ngang":      "Duong ke ngang (buoc nhay)",
    "nen_lap":       "Nen lap lai (tu tuong quan)",
}


# ===========================================================================
#  PHAN 1 — DO
#  Giu nguyen y het ban cu. Chi them tham so `chi_tiet`: khi bat, ham ghi lai
#  CHO NAO trong anh da lam ra tung con so (bien `ct`) de phan ve dung lai.
#  Cac con so KHONG doi mot ly nao — day la dieu kien de bang so cu con tin duoc.
# ===========================================================================

def _mat_na_troi(a):
    r, g, b = a[:, :, 0], a[:, :, 1], a[:, :, 2]
    return (b > r + 55) & (b > 135) & (g > r + 15)


def do_mot_anh(p, chi_tiet=False):
    im = Image.open(p).convert("RGB")
    a = np.array(im).astype(int)
    h, w, _ = a.shape
    nho = np.array(im.resize((1000, 500))).astype(int)
    kq = {}
    # `ct` la "ban do vet thuong": moi khoa la du lieu khong gian cua mot chi so.
    # Toa do trong ct dung hai he — he anh nho (500x1000) va he anh goc (h x w) —
    # nen khoa nao dung he nao deu ghi ro o cho luu.
    ct = {"anh": im, "w": w, "h": h} if chi_tiet else None

    # 1) TROI chiem bao nhieu man — trong qua thi buc tranh rong
    troi = _mat_na_troi(nho)
    kq["troi_toan_man"] = troi.mean() * 100
    if ct is not None:
        ct["troi"] = troi                        # he anh nho 500x1000

    # 2) O DAI NHA CUA, bao nhieu be ngang van la troi — do "lang co lien khoi khong"
    y_dau, y_cuoi = int(500 * 0.45), int(500 * 0.72)
    giua = troi[y_dau:y_cuoi]
    cot_troi = giua.mean(axis=0) > 0.6
    kq["troi_dai_nha"] = cot_troi.mean() * 100
    if ct is not None:
        ct["dai_nha"] = (y_dau, y_cuoi, cot_troi)   # he anh nho

    # 3) O TRON o dai giua nui va nha — do "chan troi co bi dong khong".
    #    Do bang DO CHI TIET chu khong bang mau: lang xa mo cung xanh nhat nen
    #    mat na mau khong phan biet duoc no voi troi.
    yd0, yd1 = int(500 * 0.42), int(500 * 0.62)
    dai = nho[yd0:yd1].astype(float)
    o = dai.reshape(dai.shape[0] // 10, 10, dai.shape[1] // 10, 10, 3)
    # std tinh RIENG tung kenh (axis 1 va 3 = trong long o), roi moi trung binh
    # ba kenh — bay so 2 o dau file.
    o_tron = o.std(axis=(1, 3)).mean(axis=2) < 6
    kq["o_tron_giua"] = o_tron.mean() * 100
    if ct is not None:
        ct["o_tron"] = (yd0, yd1, o_tron)        # o_tron: luoi (10, 100) o vuong

    # 4) DAT o phan duoi co phang li khong
    y_dat = int(500 * 0.75)
    dat = nho[y_dat:].astype(float)
    tb = dat.reshape(-1, 3).mean(axis=0)
    phang = np.abs(dat - tb).sum(axis=2) < 30
    kq["dat_phang"] = phang.mean() * 100
    if ct is not None:
        ct["dat_phang"] = (y_dat, phang)         # he anh nho

    # 5) KHOI NUI cao bao nhieu % khung — nui la khoi chu dao cua boi canh
    #    Mat na phai loai duoc GO (nha cua): da thi XAM (R~G~B), go thi R>G>B ro.
    #    Va chi quet nua TREN khung — nui khong bao gio nam duoi mat dat.
    r, g, b = a[:, :, 0], a[:, :, 1], a[:, :, 2]
    xam = (np.abs(r - g) < 26) & (np.abs(g - b) < 34)
    sang = (r > 95) & (r < 225)
    da = xam & sang
    da[int(h * 0.55):] = False
    hang = np.where(da.mean(axis=1) > 0.12)[0]
    kq["cao_khoi_nui"] = (hang.max() - hang.min()) / h * 100 if len(hang) else 0.0
    if ct is not None:
        ct["nui"] = (da, int(hang.min()) if len(hang) else None,
                     int(hang.max()) if len(hang) else None)   # he anh goc

    # 6) DUONG KE NGANG chay suot man — buoc nhay do sang lon nhat o vung nen.
    #    Bo vung HUD (tren 120px) va vung mat dat (duoi 62% khung).
    #    Chi quet vung NEN: duoi HUD, tren mat dat. Lay qua 62% la dinh mep co —
    #    do la ranh gioi dia hinh THAT, khong phai loi ghep lop.
    #    CHI tinh la duong ke neu no CHAY SUOT — nhieu cot cung nhay mot luc.
    #    Neu chi mot vai cot nhay thi do la ranh gioi TU NHIEN cua vat (dinh nui,
    #    mep mai nha), khong phai loi ghep lop. Da tung doc nham dinh nui thanh
    #    duong ke va di "sua" nham.
    y0, y1 = 120, int(h * 0.55)
    x0, x1 = int(w * 0.2), int(w * 0.95)
    cot = a[:, x0:x1].mean(axis=2)
    nhay_tung_cot = np.abs(np.diff(cot, axis=0))
    # buoc nhay o PHAN VI 25 cua cac cot: cao thi nghia la DA SO cot deu nhay
    lan_toa = np.percentile(nhay_tung_cot[y0:y1], 25, axis=1)
    kq["ke_ngang"] = float(lan_toa.max()) if len(lan_toa) else 0.0
    if ct is not None:
        ct["ke_ngang"] = (y0, y1, x0, x1, lan_toa)   # he anh goc

    # 7) NEN CO LAP KHONG — tu tuong quan dai nen theo chieu ngang
    yn0, yn1 = int(h * 0.28), int(h * 0.36)
    dai_nen = a[yn0:yn1].mean(axis=(0, 2))
    dai_nen = dai_nen - dai_nen.mean()
    n = len(dai_nen)
    tot = 0.0
    lag_tot = 0          # nho lai lag nao trung nhat — chi de VE, khong doi so
    for lag in range(200, n // 2):
        x, y = dai_nen[:-lag], dai_nen[lag:]
        m = np.sqrt((x ** 2).sum() * (y ** 2).sum())
        if m > 0:
            v = float((x * y).sum() / m)
            if v > tot:
                tot, lag_tot = v, lag
    kq["nen_lap"] = tot * 100
    if ct is not None:
        ct["nen_lap"] = (yn0, yn1, lag_tot, dai_nen)   # he anh goc
    return kq, ct


def danh_gia(khoa, v, co_concept):
    """Tra ' OK' / ' !!' / ''.

    Tach ra thanh ham rieng de BANG SO va ANH CHAN DOAN dung chung dung mot
    phep danh gia — hai ben khong bao gio noi khac nhau ve cung mot chi so.
    """
    if khoa not in MOC or not co_concept:
        return ""
    muc, huong = MOC[khoa]
    if khoa == "ke_ngang":
        dat = v < 10          # duoi 10 la mat thuong khong thay
    elif huong.startswith("thap"):
        dat = v <= muc * 1.35 + 3
    else:
        dat = v >= muc * 0.75
    return " OK" if dat else " !!"


# ===========================================================================
#  PHAN 2 — VE ANH CHAN DOAN
#  Bay o: 1 o anh goc de doi chieu + 6 o chi so co the to len anh + 1 o chi so
#  nen lap. Moi o mot mau rieng, co dong chu giai ngay tren dau — de nguoi
#  khong lap trinh mo ra la hieu ngay "cho do dang bi gi".
# ===========================================================================

RONG_O = 800          # be ngang moi o anh trong bang chan doan
LE = 20
KHE = 18

DO = (255, 60, 60)
CAM = (255, 150, 40)
VANG = (255, 230, 0)
LUC = (60, 230, 110)
# Mat na troi phai to bang mau TIM CANH SEN, khong phai mau xanh: to xanh len
# nen troi xanh thi khong nhin thay gi — ban dau da lam vay va o so 1 coi nhu
# trong. Mau nao khong co san trong anh thi mau do moi chi duoc cho.
TIM = (255, 60, 220)
TRANG = (255, 255, 255)
NEN_BANG = (26, 28, 34)

_FONT = {}


def _chu_font(co):
    if co not in _FONT:
        for ten in ("C:/Windows/Fonts/segoeui.ttf", "C:/Windows/Fonts/arial.ttf"):
            try:
                _FONT[co] = ImageFont.truetype(ten, co)
                break
            except OSError:
                continue
        else:
            try:
                _FONT[co] = ImageFont.load_default(co)
            except TypeError:
                _FONT[co] = ImageFont.load_default()
    return _FONT[co]


def _chu(v, xy, s, co=17, mau=TRANG):
    """Chu co vien den — de doc duoc ca khi de len vung sang lan vung toi."""
    v.text(xy, s, font=_chu_font(co), fill=mau, stroke_width=2, stroke_fill=(15, 15, 15))


def _cat_dong(s, font, rong):
    """Cat mot cau thanh nhieu dong vua be ngang `rong`.

    Khong cat thi dong chu giai cua o ben trai TRAN sang o ben phai va de len
    chu cua no — da bi that o ban dau, doc ra thanh mot mo chu chong nhau.
    """
    dong, hien = [], ""
    for t in s.split():
        thu = (hien + " " + t).strip()
        if hien and font.getlength(thu) > rong:
            dong.append(hien)
            hien = t
        else:
            hien = thu
    if hien:
        dong.append(hien)
    return dong


def _nen_mo(o, x0, y0, x1, y1, muc=0.42):
    """Lam toi mot mang de duong cong / chu ve len tren con doc duoc."""
    PW, PH = o.size
    x0, x1 = int(max(0, x0)), int(min(PW, x1))
    y0, y1 = int(max(0, y0)), int(min(PH, y1))
    if x1 <= x0 or y1 <= y0:
        return
    a = np.array(o).astype(float)
    a[y0:y1, x0:x1] *= muc
    o.paste(Image.fromarray(a.astype(np.uint8)))


def _to_vung(o, mat_na, mau, dam=0.55):
    """To mau len nhung pixel mat_na bat.

    mat_na co the o do phan giai khac (anh nho 500x1000, hay anh goc) — phong
    ve co o anh bang NEAREST chu khong phai LANCZOS: lam nhoe ranh gioi vung
    loi thi mat luon cai chinh xac ma anh nay sinh ra de chi.
    """
    if mat_na is None or not mat_na.any():
        return
    m = Image.fromarray(mat_na.astype(np.uint8) * 255).resize(o.size, Image.NEAREST)
    m = np.array(m) > 127
    a = np.array(o).astype(float)
    a[m] = a[m] * (1 - dam) + np.array(mau, dtype=float) * dam
    o.paste(Image.fromarray(a.astype(np.uint8)))


def _ke_dut(v, x0, y, x1, mau, buoc=16, day=2):
    x = x0
    while x < x1:
        v.line([(x, y), (min(x + buoc // 2, x1), y)], fill=mau, width=day)
        x += buoc


def _mui_ten_doc(v, x, y0, y1, mau):
    v.line([(x, y0), (x, y1)], fill=mau, width=3)
    for y, d in ((y0, 1), (y1, -1)):
        v.line([(x - 7, y + 8 * d), (x, y), (x + 7, y + 8 * d)], fill=mau, width=3)


def _mui_ten_ngang(v, y, x0, x1, mau):
    v.line([(x0, y), (x1, y)], fill=mau, width=3)
    for x, d in ((x0, 1), (x1, -1)):
        v.line([(x + 8 * d, y - 7), (x, y), (x + 8 * d, y + 7)], fill=mau, width=3)


# --- tung o mot -----------------------------------------------------------

def _ve_troi_toan_man(o, v, ct):
    _to_vung(o, ct["troi"], TIM)


def _ve_troi_dai_nha(o, v, ct):
    y_dau, y_cuoi, cot_troi = ct["dai_nha"]
    mn = np.zeros((500, 1000), bool)
    mn[y_dau:y_cuoi] = cot_troi[None, :]      # to nguyen COT trong dai nha cua
    _to_vung(o, mn, DO, 0.6)
    PW, PH = o.size
    for yy in (y_dau, y_cuoi):
        _ke_dut(v, 0, yy / 500 * PH, PW, VANG)
    _chu(v, (10, y_dau / 500 * PH + 6), "dải nhà cửa (45–72% chiều cao)", 16, VANG)


def _ve_o_tron_giua(o, v, ct):
    yd0, yd1, o_tron = ct["o_tron"]
    mn = np.zeros((500, 1000), bool)
    mn[yd0:yd1] = np.repeat(np.repeat(o_tron, 10, axis=0), 10, axis=1)
    _to_vung(o, mn, DO, 0.55)
    PW, PH = o.size
    # Vien tung o chi ve khi so o loi con it — nhieu qua thi vien bien thanh
    # mot tam luoi, che mat chinh cai anh nen nhin.
    if o_tron.sum() <= 150:
        ch, cw = 10 / 500 * PH, 10 / 1000 * PW
        for i, j in zip(*np.where(o_tron)):
            x, y = j * cw, (yd0 + i * 10) / 500 * PH
            v.rectangle([x, y, x + cw, y + ch], outline=DO, width=2)
    for yy in (yd0, yd1):
        _ke_dut(v, 0, yy / 500 * PH, PW, VANG)
    _chu(v, (10, yd0 / 500 * PH + 6), "dải giữa (42–62% chiều cao)", 16, VANG)


def _ve_dat_phang(o, v, ct):
    y_dat, phang = ct["dat_phang"]
    mn = np.zeros((500, 1000), bool)
    mn[y_dat:] = phang
    _to_vung(o, mn, CAM, 0.6)
    PW, PH = o.size
    _ke_dut(v, 0, y_dat / 500 * PH, PW, VANG)
    _chu(v, (10, y_dat / 500 * PH + 6), "vùng đất (1/4 dưới khung)", 16, VANG)


def _ve_cao_khoi_nui(o, v, ct):
    da, r_dinh, r_chan = ct["nui"]
    _to_vung(o, da, LUC, 0.5)
    PW, PH = o.size
    h = ct["h"]
    moc = int(h * 0.55)
    _ke_dut(v, 0, 0.55 * PH, PW, (170, 170, 170))
    # Ghi ben PHAI: vach 55% gan nhu luon nam sat vach "chan khoi nui", hai
    # dong chu de trai se de len nhau.
    s = "chỉ quét nửa trên — núi không nằm dưới đất"
    _chu(v, (PW - _chu_font(15).getlength(s) - 10, 0.55 * PH + 5), s, 15, (200, 200, 200))
    if r_dinh is None:
        _chu(v, (12, PH / 2), "KHÔNG tìm thấy khối núi nào", 20, DO)
        return
    yd, yc = r_dinh / h * PH, r_chan / h * PH
    for yy in (yd, yc):
        v.line([(0, yy), (PW, yy)], fill=VANG, width=3)
    _mui_ten_doc(v, PW - 74, yd, yc, VANG)
    _chu(v, (12, yd - 26), f"đỉnh khối núi (hàng {r_dinh})", 16, VANG)
    # Neu chan nui cham dung mep quet 55% thi con so nay bi CAT bot chu khong
    # phai la nui het o do — noi thang ra, khong de nguoi doc hieu nham.
    ghi = f"chân khối núi (hàng {r_chan})"
    if r_chan >= moc - 3:
        ghi += " — chạm mép quét, đá còn kéo xuống dưới"
    _chu(v, (12, yc - 26), ghi, 16, VANG)


def _ve_ke_ngang(o, v, ct):
    y0, y1, x0, x1, lan_toa = ct["ke_ngang"]
    PW, PH = o.size
    sx, sy = PW / ct["w"], PH / ct["h"]
    v.rectangle([x0 * sx, y0 * sy, x1 * sx, y1 * sy], outline=VANG, width=2)
    _chu(v, (x0 * sx + 6, y0 * sy + 6), "vùng quét (dưới HUD, trên mặt đất)", 15, VANG)

    # Ho so buoc nhay tung hang, ve doc ben trai. Nhin duong nay biet ngay
    # "mot hang vot han len" (dung la duong ke) hay "ca vung deu lam ram"
    # (chi la van canh binh thuong).
    K, GOC = 6.0, 12          # 6 px cho moi don vi buoc nhay
    _nen_mo(o, 0, y0 * sy, GOC + 27 * K, y1 * sy)   # nen toi, khong thi duong trang chim vao mai ngoi
    v.line([(GOC + 10 * K, y0 * sy), (GOC + 10 * K, y1 * sy)], fill=VANG, width=1)
    diem = [(GOC + min(float(t), 26.0) * K, (y0 + i) * sy) for i, t in enumerate(lan_toa)]
    if len(diem) > 1:
        v.line(diem, fill=(255, 255, 255), width=2)
    _chu(v, (GOC + 10 * K + 4, y1 * sy - 24), "ngưỡng 10", 14, VANG)

    # Hang nao DA SO cot cung nhay tu 10 tro len moi la duong ke that.
    for i in np.where(lan_toa >= 10)[0]:
        v.line([(x0 * sx, (y0 + i) * sy), (x1 * sx, (y0 + i) * sy)], fill=DO, width=3)
    if len(lan_toa):
        i = int(np.argmax(lan_toa))
        manh = float(lan_toa[i])
        yy = (y0 + i) * sy
        mau = DO if manh >= 10 else VANG
        _ke_dut(v, x0 * sx, yy, x1 * sx, mau, 20, 3)
        ghi = f"hàng {y0 + i} · bước nhảy {manh:.1f}"
        if manh < 10:
            ghi += "  (mạnh nhất, nhưng dưới ngưỡng — chưa phải lỗi)"
        _chu(v, (x0 * sx + 6, yy + 6), ghi, 16, mau)


def _ve_nen_lap(o, v, ct):
    yn0, yn1, lag, prof = ct["nen_lap"]
    PW, PH = o.size
    sx, sy = PW / ct["w"], PH / ct["h"]
    ytb = (yn0 + yn1) / 2 * sy
    cao = max(40.0, (yn1 - yn0) * sy)          # dai nen rat mong, noi cao len de nhin duoc
    y_tren, y_duoi = ytb - cao, ytb + cao

    # Nen toi mo sau duong cong, khong thi duong trang chim vao troi sang.
    _nen_mo(o, 0, y_tren, PW, y_duoi, 0.45)
    # Vach VANG danh dau dung dai nen that (28-36%), con vung toi chi la CHO VE
    # duong cong — noi rong ra cho nhin duoc. Ve khung quanh vung toi thi nguoi
    # doc tuong dai nen day gap doi thuc te.
    for yy in (yn0 * sy, yn1 * sy):
        _ke_dut(v, 0, yy, PW, VANG)
    _chu(v, (10, yn1 * sy + 4), "dải nền được so (28–36% chiều cao)", 15, VANG)

    bien = np.abs(prof).max() or 1.0
    def diem(mang, lech_px):
        return [((i + lech_px) * sx, ytb - float(t_) / bien * (cao * 0.85))
                for i, t_ in enumerate(mang)]

    v.line(diem(prof, 0), fill=TRANG, width=2)
    if lag:
        # Ve chinh duong do dich sang trai `lag` px: cho nao no NAM DE len duong
        # trang nghia la doan nen do lap lai sau `lag` px.
        v.line(diem(prof, -lag), fill=CAM, width=2)
        x_a = PW * 0.12
        _mui_ten_ngang(v, y_tren - 16, x_a, x_a + lag * sx, CAM)
        _chu(v, (x_a + 6, y_tren - 44), f"lặp lại sau {lag} px", 16, CAM)
    else:
        _chu(v, (12, y_duoi + 8), "không tìm thấy bước lặp nào", 16, LUC)


VE = {
    "troi_toan_man": _ve_troi_toan_man,
    "troi_dai_nha":  _ve_troi_dai_nha,
    "o_tron_giua":   _ve_o_tron_giua,
    "dat_phang":     _ve_dat_phang,
    "cao_khoi_nui":  _ve_cao_khoi_nui,
    "ke_ngang":      _ve_ke_ngang,
    "nen_lap":       _ve_nen_lap,
}

NHAN = {
    "troi_toan_man": ("1. Trời chiếm màn",
                      "Tím = pixel bị tính là trời. LƯU Ý: đây là mặt nạ theo MÀU nên nó bắt cả "
                      "làng xa mờ (bẫy 1 ở đầu file) — nhìn để hiểu con số, đừng chữa máy móc "
                      "theo nó."),
    "troi_dai_nha":  ("2. Dải nhà cửa là trời",
                      "Đỏ = cột mà trời lọt xuống tận dải nhà cửa. Làng liền một khối thì dải "
                      "này không còn cột đỏ nào."),
    "o_tron_giua":   ("3. Ô trơn ở dải giữa",
                      "Đỏ = ô 10×10 không có chi tiết. Đo bằng ĐỘ CHI TIẾT chứ không bằng màu, "
                      "nên làng xa mờ vẫn được tính là có nội dung."),
    "dat_phang":     ("4. Đất phẳng lì",
                      "Cam = pixel gần y hệt màu trung bình của cả vùng đất; chỗ KHÔNG tô là chỗ "
                      "còn vân. Càng nhiều cam thì lòng đất càng trống, cần thêm vân đất / kè đá."),
    "cao_khoi_nui":  ("5. Khối núi cao",
                      "Lục = pixel đá (xám, sáng vừa). Hai vạch vàng là đỉnh và chân khối núi — "
                      "khoảng cách giữa chúng chính là con số."),
    "ke_ngang":      ("6. Đường kẻ ngang",
                      "Đỏ = hàng mà ĐA SỐ cột cùng nhảy độ sáng (bẫy 3: vài cột nhảy chỉ là đỉnh "
                      "núi, không phải lỗi). Đường trắng bên trái là bước nhảy của từng hàng, "
                      "vạch vàng dọc là ngưỡng 10."),
    "nen_lap":       ("7. Nền lặp lại",
                      "Trắng = độ sáng dải nền theo chiều ngang; cam = chính nó dịch đi N px. "
                      "Hai đường đè khít nhau càng nhiều thì ảnh nền càng lặp."),
}


def ve_anh_chan_doan(duong_dan, kq, ct, goc, thu_muc):
    """Ve bang 8 o: 1 o anh goc + 7 o chi so, moi o to len dung cho lam hong."""
    im = ct["anh"]
    cao_o = max(1, round(im.height * RONG_O / im.width))
    khoa_ds = list(TEN_DEP)
    tong = 1 + len(khoa_ds)
    hang = (tong + 1) // 2
    cao_dau = 74

    # Cat truoc moi dong chu giai theo be ngang cua o, roi moi biet phai chua
    # bao nhieu dong — de vua khit chu khong tran sang o ben canh.
    f_giai = _chu_font(15)
    giai_ds = {k: _cat_dong(g, f_giai, RONG_O - 4) for k, (_, g) in NHAN.items()}
    giai_ds[None] = _cat_dong("Không tô gì — để so sánh với bảy ô còn lại.", f_giai, RONG_O - 4)
    so_dong = max(len(x) for x in giai_ds.values())
    cao_nhan = 8 + 28 + so_dong * 20 + 6

    W = LE * 2 + RONG_O * 2 + KHE
    H = cao_dau + LE + hang * (cao_nhan + cao_o) + (hang - 1) * KHE + LE
    bang = Image.new("RGB", (W, H), NEN_BANG)
    vb = ImageDraw.Draw(bang)

    ten = os.path.basename(duong_dan)
    vb.text((LE, 16), f"CHỖ NÀO LÀM HỎNG CHỈ SỐ — {ten}  ({im.width}×{im.height})",
            font=_chu_font(26), fill=TRANG)
    vb.text((LE, 48), "Mỗi ô là một chỉ số. Vùng tô màu chính là chỗ làm ra con số đó. "
                      "Số bên phải tên ô: ảnh này · concept · sổ tay.",
            font=_chu_font(16), fill=(175, 180, 190))

    nen_o = im.resize((RONG_O, cao_o), Image.LANCZOS)
    for k, khoa in enumerate([None] + khoa_ds):
        i, j = divmod(k, 2)
        x = LE + j * (RONG_O + KHE)
        y = cao_dau + LE + i * (cao_nhan + cao_o + KHE)
        o = nen_o.copy()
        v = ImageDraw.Draw(o)
        if khoa is None:
            tieu_de, mau_de = "0. Ảnh gốc", TRANG
        else:
            VE[khoa](o, v, ct)
            gia = kq[khoa]
            dg = danh_gia(khoa, gia, goc is not None).strip()
            so = f"{gia:.1f}"
            if goc is not None:
                so += f"  ·  concept {goc[khoa]:.1f}"
            so += f"  ·  sổ tay {NGUONG_SO_TAY[khoa]}"
            tieu_de = f"{NHAN[khoa][0]}:  {so}" + (f"   [{dg}]" if dg else "")
            mau_de = DO if dg == "!!" else (LUC if dg == "OK" else TRANG)
        f_de = _chu_font(21 if _chu_font(21).getlength(tieu_de) <= RONG_O else 18)
        vb.text((x, y + 8), tieu_de, font=f_de, fill=mau_de)
        for n, dong in enumerate(giai_ds[khoa]):
            vb.text((x, y + 38 + n * 20), dong, font=f_giai, fill=(175, 180, 190))
        bang.paste(o, (x, y + cao_nhan))
        vb.rectangle([x, y + cao_nhan, x + RONG_O - 1, y + cao_nhan + cao_o - 1],
                     outline=(70, 74, 84), width=1)

    os.makedirs(thu_muc, exist_ok=True)
    ra = os.path.join(thu_muc, os.path.splitext(ten)[0] + "-chan-doan.png")
    bang.save(ra)
    return ra


# ===========================================================================
#  PHAN 3 — CHAY
# ===========================================================================

def main(argv):
    ve_anh = "--khong-anh" not in argv
    thu_muc = THU_MUC_RA
    if "--ra" in argv:
        i = argv.index("--ra")
        if i + 1 >= len(argv):
            print("--ra can mot thu muc di kem")
            return 1
        thu_muc = argv[i + 1]
        argv = argv[:i] + argv[i + 2:]
    duong_dan = [x for x in argv if not x.startswith("--")]

    if not duong_dan:
        print(__doc__)
        return 1
    goc = do_mot_anh(CONCEPT)[0] if os.path.exists(CONCEPT) else None
    ket = []
    for p in duong_dan:
        if not os.path.exists(p):
            continue
        kq, ct = do_mot_anh(p, chi_tiet=ve_anh)
        ket.append((os.path.basename(p), p, kq, ct))
    if not ket:
        print("khong doc duoc anh nao")
        return 1

    rong = max(len(v) for v in TEN_DEP.values())
    print(f"{'':<{rong}} | {'concept':>8} |" + "".join(f" {t[:14]:>14} |" for t, _, _, _ in ket))
    print("-" * (rong + 12 + 17 * len(ket)))
    for khoa, nhan in TEN_DEP.items():
        c = f"{goc[khoa]:8.1f}" if goc and khoa in goc else "       -"
        dong = f"{nhan:<{rong}} | {c} |"
        for _, _, k, _ in ket:
            dong += f" {k[khoa]:11.1f}{danh_gia(khoa, k[khoa], bool(goc)):>3} |"
        print(dong)
    print("\nOK = da bang hoac hon concept · !! = con chenh")
    print("Duong ke ngang: duoi 10 la mat thuong khong thay. Nen lap: duoi 40 la khong lo.")

    if ve_anh:
        # Chu o day de ASCII: console Windows khong in duoc dau tieng Viet,
        # chu co dau se ra mot mo ky tu rac.
        print("\nAnh chan doan (to mau thang len anh, cho biet CHO NAO lam chi so hong):")
        for _, p, kq, ct in ket:
            print("  " + ve_anh_chan_doan(p, kq, ct, goc, thu_muc))
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
