# -*- coding: utf-8 -*-
"""THUOC DO MAP — mot lenh, 10 giay, thay cho mot dot soi 45-80 phut.

    python tools/map-lang-la/do-map.py <anh-game.png> [anh-game-2.png ...]

Ca ngay 22/08 minh chay bon dot soi bang agent, moi dot 45-80 phut, de tim ra
nhung con so ma may do duoc trong 10 giay. File nay gom dung nhung phep do da
that su bat duoc loi, moi phep kem MOC so voi concept da duyet.

DUNG NO TRUOC, agent SAU. Agent chi con dung o hai viec ma may khong lam duoc:
nhan xet tham my ("net ve choi voi nen") va doi chieu y do ("concept muon gi").

BAY KHI DO — da mac ca hai, ghi lai de khong mac lai:
  1. Mat na "mau troi" bat ca lang xa mo (cung xanh nhat) -> so khong nhuc nhich
     du anh doi ro. Phai do bang DO CHI TIET, khong do bang mau.
  2. "Do lech chuan cua hang" phai tinh theo chieu ngang cho TUNG KENH mau.
     Tinh std tren ca 3 kenh thi mot hang dong mau (12,95,169) ra std ~64, doc
     nham thanh "con nhieu van" trong khi hang do da phang tuyet doi.
"""
import sys
import os
import numpy as np
from PIL import Image

CONCEPT = "D:/NRO_GoiBanGiao/docs/my-thuat/map-1-lang-la/concept-map1-langla.png"

# Moc lay tu concept da duyet. "tot" = huong can di.
MOC = {
    "troi_toan_man":  (24, "thap hon la tot"),
    "troi_dai_nha":   (0,  "thap hon la tot"),
    "o_tron_giua":    (6,  "thap hon la tot"),
    "dat_phang":      (27, "thap hon la tot"),
    "cao_khoi_nui":   (31, "cao hon la tot"),
    "ke_ngang":       (0,  "thap hon la tot"),
}


def _mat_na_troi(a):
    r, g, b = a[:, :, 0], a[:, :, 1], a[:, :, 2]
    return (b > r + 55) & (b > 135) & (g > r + 15)


def do_mot_anh(p):
    im = Image.open(p).convert("RGB")
    a = np.array(im).astype(int)
    h, w, _ = a.shape
    nho = np.array(im.resize((1000, 500))).astype(int)
    kq = {}

    # 1) TROI chiem bao nhieu man — trong qua thi bức tranh rong
    kq["troi_toan_man"] = _mat_na_troi(nho).mean() * 100

    # 2) O DAI NHA CUA, bao nhieu be ngang van la troi — do "lang co lien khoi khong"
    giua = _mat_na_troi(nho)[int(500 * 0.45):int(500 * 0.72)]
    kq["troi_dai_nha"] = (giua.mean(axis=0) > 0.6).mean() * 100

    # 3) O TRON o dai giua nui va nha — do "chan troi co bi dong khong".
    #    Do bang DO CHI TIET chu khong bang mau: lang xa mo cung xanh nhat nen
    #    mat na mau khong phan biet duoc no voi troi.
    dai = nho[int(500 * 0.42):int(500 * 0.62)].astype(float)
    o = dai.reshape(dai.shape[0] // 10, 10, dai.shape[1] // 10, 10, 3)
    kq["o_tron_giua"] = (o.std(axis=(1, 3)).mean(axis=2) < 6).mean() * 100

    # 4) DAT o phan duoi co phang li khong
    dat = nho[int(500 * 0.75):].astype(float)
    tb = dat.reshape(-1, 3).mean(axis=0)
    kq["dat_phang"] = (np.abs(dat - tb).sum(axis=2) < 30).mean() * 100

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

    # 6) DUONG KE NGANG chay suot man — buoc nhay do sang lon nhat o vung nen.
    #    Bo vung HUD (tren 120px) va vung mat dat (duoi 62% khung).
    #    Chi quet vung NEN: duoi HUD, tren mat dat. Lay qua 62% la dinh mep co —
    #    do la ranh gioi dia hinh THAT, khong phai loi ghep lop.
    #    CHI tinh la duong ke neu no CHAY SUOT — nhieu cot cung nhay mot luc.
    #    Neu chi mot vai cot nhay thi do la ranh gioi TU NHIEN cua vat (dinh nui,
    #    mep mai nha), khong phai loi ghep lop. Da tung doc nham dinh nui thanh
    #    duong ke va di "sua" nham.
    y0, y1 = 120, int(h * 0.55)
    cot = a[:, int(w * 0.2):int(w * 0.95)].mean(axis=2)
    nhay_tung_cot = np.abs(np.diff(cot, axis=0))
    # buoc nhay o PHAN VI 25 cua cac cot: cao thi nghia la DA SO cot deu nhay
    lan_toa = np.percentile(nhay_tung_cot[y0:y1], 25, axis=1)
    kq["ke_ngang"] = float(lan_toa.max()) if len(lan_toa) else 0.0

    # 7) NEN CO LAP KHONG — tu tuong quan dai nen theo chieu ngang
    dai_nen = a[int(h * 0.28):int(h * 0.36)].mean(axis=(0, 2))
    dai_nen = dai_nen - dai_nen.mean()
    n = len(dai_nen)
    tot = 0.0
    for lag in range(200, n // 2):
        x, y = dai_nen[:-lag], dai_nen[lag:]
        m = np.sqrt((x ** 2).sum() * (y ** 2).sum())
        if m > 0:
            tot = max(tot, float((x * y).sum() / m))
    kq["nen_lap"] = tot * 100
    return kq


def main(duong_dan):
    if not duong_dan:
        print(__doc__)
        return 1
    goc = do_mot_anh(CONCEPT) if os.path.exists(CONCEPT) else None
    ket = [(os.path.basename(p), do_mot_anh(p)) for p in duong_dan if os.path.exists(p)]
    if not ket:
        print("khong doc duoc anh nao")
        return 1

    ten_dep = {
        "troi_toan_man": "Troi chiem man",
        "troi_dai_nha":  "Dai nha cua la troi",
        "o_tron_giua":   "O tron o dai giua",
        "dat_phang":     "Dat phang li",
        "cao_khoi_nui":  "Khoi nui cao",
        "ke_ngang":      "Duong ke ngang (buoc nhay)",
        "nen_lap":       "Nen lap lai (tu tuong quan)",
    }
    rong = max(len(v) for v in ten_dep.values())
    print(f"{'':<{rong}} | {'concept':>8} |" + "".join(f" {t[:14]:>14} |" for t, _ in ket))
    print("-" * (rong + 12 + 17 * len(ket)))
    for khoa, nhan in ten_dep.items():
        c = f"{goc[khoa]:8.1f}" if goc and khoa in goc else "       -"
        dong = f"{nhan:<{rong}} | {c} |"
        for _, k in ket:
            v = k[khoa]
            danh_gia = ""
            if khoa in MOC and goc:
                muc = MOC[khoa][0]
                thap_tot = MOC[khoa][1].startswith("thap")
                if khoa == "ke_ngang":
                    dat = v < 10          # duoi 10 la mat thuong khong thay
                elif thap_tot:
                    dat = v <= muc * 1.35 + 3
                else:
                    dat = v >= muc * 0.75
                danh_gia = " OK" if dat else " !!"
            dong += f" {v:11.1f}{danh_gia:>3} |"
        print(dong)
    print("\nOK = da bang hoac hon concept · !! = con chenh")
    print("Duong ke ngang: duoi 10 la mat thuong khong thay. Nen lap: duoi 40 la khong lo.")
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
