# -*- coding: utf-8 -*-
"""ghep-khung.py — ghép đầu/thân/chân của nhân vật đúng phép toán vẽ của client Unity
(Char.paintCharBody, Char.cs ~6057-6084) cho cả 33 khung `cf`, KHÔNG cần mở Unity.

Vì sao có file này (05/09/2026): chơi thử thấy "đầu tách khỏi thân" ở một số khung
chuyển động. Muốn sửa phải ĐO được: khung nào lệch, lệch bao nhiêu px, trước và sau khi
đổi part.bin/sprite. Unity chụp màn hình không tách được từng khung; công cụ này thì có.

Phép toán vẽ (hướng phải, num2=+1, neo TOP|LEFT), k = 0 đầu, 1 chân, 2 thân:
    idx = CharInfo[cf][k][0]
    x   = cx + CharInfo[cf][k][1] + part.pi[idx].dx
    y   = cy - CharInfo[cf][k][2] + part.pi[idx].dy        (y là MÉP TRÊN ảnh)
Thứ tự vẽ trong client: đầu → chân → thân (Char.cs 6081, 6083, 6084).

Dùng:
  python tools/nhan-vat/ghep-khung.py --out ket-qua.png
  Tuỳ chọn: --part-bin, --sprites (thư mục Small<id>.png x1), --charinfo (JSON 33 khung),
            --sprites-closure (thư mục closure kiểu sprites/<id>/x1/sprite.png),
            --head/--body/--leg (id part, mặc định 30/57/58), --zoom (mặc định 1 = ảnh x1).
In ra bảng đo mỗi khung: khe đầu–thân (px), khe thân–chân, chân so với cy, chiều cao tổng
(đỉnh đầu → bàn chân), và SỐ MẶT.
Khe = mép trên vùng ĐỤC của mảnh dưới trừ mép dưới vùng đục của mảnh trên; >0 = HỞ, <0 = chồng.

Bộ dò "hai khuôn mặt" (cột `so mat`, thêm 05/09/2026 chiều): chẩn đoán lại cho thấy "đầu
tách thân" thực ra là mảnh THÂN chứa cả một khuôn mặt (cắt ở 28 % chiều cao khung nguồn rơi
lên trán ở khung chạy/đánh) và được vẽ đè lên đầu. Khe bbox ÂM vẫn có thể là hai mặt, nên
chỉ đo khe là chưa đủ. Cột `so mat` = số vùng da liên thông đủ lớn (≥ NGUONG_MAT_PX_X1 ×
zoom² px) trong HỘP ĐẦU của ảnh ghép: bbox mảnh đầu nới ngang 2 px, kéo đáy xuống tới 60 %
mảnh thân (vai/ngực — nơi mặt thứ hai hay xuất hiện; đai lưng cam và bàn tay hạ thấp nằm
ngoài). Nhân vật bình thường phải có đúng 1.

Cửa kiểm tự động (thoát mã ≠ 0 khi có cf vi phạm, in danh sách cf):
  python tools/nhan-vat/ghep-khung.py --out kq.png \
      --kiem "dau-than=-9,0;than-chan=-11,-2;so-mat=1;chan=-11,-10"
  Khoá cho phép: dau-than, than-chan, chan, so-mat, cao, rong-than. Giá trị `a,b` = đoạn
  [a, b] (kể cả hai đầu); một số = phải bằng đúng. Ô đo không được ("?") cũng là vi phạm.

Chỉ đo một số cf (06/09/2026): `--cf "0,1"`. NPC (Iruka, part 2189/2190/2191) chỉ đảo cf 0/1
(Npc.cs:179–185) nên áp cửa kiểm cho 33 cf là vô nghĩa với NPC:
  python tools/nhan-vat/ghep-khung.py --head 2189 --body 2190 --leg 2191 --cf 0,1 --out kq.png \
      --kiem "dau-than=-9,-1;than-chan=-11,-2;so-mat=1;chan=-11"
"""
import argparse
import json
import os
import struct

import numpy as np
from PIL import Image, ImageDraw

try:  # scipy nhanh hơn; thiếu thì dùng BFS numpy thuần — cùng kết quả (4 láng giềng)
    from scipy import ndimage as _ndi
except ImportError:  # pragma: no cover - chỉ chạy khi máy không có scipy
    _ndi = None

MAC_DINH_PART_BIN = r"D:\NRO_GoiBanGiao\server\internal\game\assets\part.bin"
MAC_DINH_SPRITES = r"D:\NRO_Project\Assets\Resources\res\x1\smallimage"
MAC_DINH_CHARINFO = os.path.join(os.path.dirname(os.path.abspath(__file__)), "charinfo.json")
SO_KHUNG = {0: 3, 1: 17, 2: 14}

# --- bộ dò "hai khuôn mặt" ------------------------------------------------------------
# Ngưỡng màu da (HSV) lấy nguyên từ thu-nghiem-khe-dau-than/do-moc-giai-phau.py, đã đo trên
# 17 khung nguồn: da mặt h 20–48°, s 0,35–0,8, v > 0,55. Đai lưng cam nằm ở h 10–20°, băng
# đô / quần / tay áo navy ở h 195–250° nên không lọt vào mặt nạ.
NGUONG_DA = {"h": (20.0, 48.0), "s": (0.35, 0.8), "v": 0.55}
# Vùng da liên thông phải ≥ ngưỡng này (px ở x1, nhân zoom²) mới được tính là một "mặt".
# Đo 05/09/2026 trên ảnh ghép x1: mặt thật 61–97 px; mặt thứ hai do mảnh thân cắt 28 % mang
# theo 27–88 px; nắm tay / bàn tay / ngón chân ≤ 8 px. 16 nằm giữa, cách cả hai phía ≥ 1,7×.
NGUONG_MAT_PX_X1 = 16
# Hộp đầu: nới ngang mỗi bên (px x1) và phần mảnh thân tính vào hộp (từ mép trên thân xuống).
HOP_DAU_NOI_NGANG_X1 = 2
HOP_DAU_PHAN_THAN = 0.6

# Khoá của --kiem → tên ô đo trong bảng.
KIEM_COT = {
    "dau-than": "khe_dt",
    "than-chan": "khe_tc",
    "chan": "chan_cy",
    "so-mat": "so_mat",
    "cao": "cao_tong",
    "rong-than": "rong_than",
}


def doc_part_bin(duong):
    d = open(duong, "rb").read()
    n = struct.unpack_from(">h", d, 0)[0]
    off = 2
    parts = []
    for _ in range(n):
        t = d[off]
        off += 1
        fr = [struct.unpack_from(">hbb", d, off + k * 4) for k in range(SO_KHUNG[t])]
        off += SO_KHUNG[t] * 4
        parts.append((t, fr))
    return parts


def bbox_duc(im):
    return im.split()[3].getbbox()  # (l, t, r, b) hoặc None


def duong_sprite(thu_muc, sid, kieu_closure, zoom):
    """Đường tới ảnh của sprite `sid`.

    Hai bố cục cùng tồn tại: thư mục client phẳng (`Small<id>.png`) và thư mục
    closure (`sprites/<id>/x<zoom>/sprite.png`). Nhận cả hai để đo được đầu ra
    closure mà không phải chép/đổi tên tay.
    """
    if kieu_closure:
        return os.path.join(thu_muc, str(sid), "x%d" % zoom, "sprite.png")
    return os.path.join(thu_muc, "Small%d.png" % sid)


def ghep(cf, charinfo, parts, sprites, ids, zoom, cx, cy, kieu_closure=False):
    """Trả (ảnh ghép, số đo) cho một khung."""
    o = Image.new("RGBA", (cx * 2, cy + 16 * zoom), (0, 0, 0, 0))
    vung = {}
    thu_tu = [(0, "dau", ids[0]), (1, "chan", ids[2]), (2, "than", ids[1])]  # đúng thứ tự client
    for k, ten, pid in thu_tu:
        if k >= len(charinfo[cf]):
            continue
        idx, dx_ci, dy_ci = charinfo[cf][k]
        t, fr = parts[pid]
        if idx >= len(fr):
            vung[ten] = ("IDX %d/%d" % (idx, len(fr)), None)
            continue
        sid, dx, dy = fr[idx]
        duong = duong_sprite(sprites, sid, kieu_closure, zoom)
        if not os.path.isfile(duong):
            vung[ten] = ("THIEU %d" % sid, None)
            continue
        im = Image.open(duong).convert("RGBA")
        x = cx + (dx_ci + dx) * zoom
        y = cy - dy_ci * zoom + dy * zoom
        o.alpha_composite(im, (x, y))
        bb = bbox_duc(im)
        vung[ten] = (sid, (x + bb[0], y + bb[1], x + bb[2], y + bb[3]) if bb else None)
    return o, vung


# --- bộ dò "hai khuôn mặt" ------------------------------------------------------------
def rgb_sang_hsv(rgb):
    """RGB uint8 (..., 3) → (h độ 0–360, s 0–1, v 0–1), vector hoá numpy."""
    a = np.asarray(rgb).astype(np.float32) / 255.0
    r, g, b = a[..., 0], a[..., 1], a[..., 2]
    mx = a.max(-1)
    mn = a.min(-1)
    d = mx - mn
    h = np.zeros_like(mx)
    m = d > 1e-6
    dd = np.where(d == 0, 1, d)
    rm = m & (mx == r)
    gm = m & (mx == g) & ~rm
    bm = m & ~rm & ~gm
    h[rm] = (60 * ((g - b) / dd) % 360)[rm]
    h[gm] = (60 * ((b - r) / dd) + 120)[gm]
    h[bm] = (60 * ((r - g) / dd) + 240)[bm]
    s = np.where(mx > 0, d / np.where(mx == 0, 1, mx), 0)
    return h, s, mx


def mat_na_da(anh):
    """Mặt nạ bool: pixel ĐỤC có màu da (HSV trong NGUONG_DA) của ảnh RGBA."""
    a = np.asarray(anh.convert("RGBA"))
    duc = a[..., 3] > 0
    h, s, v = rgb_sang_hsv(a[..., :3])
    (h0, h1), (s0, s1), v0 = NGUONG_DA["h"], NGUONG_DA["s"], NGUONG_DA["v"]
    return duc & (h > h0) & (h < h1) & (s > s0) & (s < s1) & (v > v0)


def gan_nhan_lien_thong(mask):
    """Gán nhãn vùng liên thông 4 láng giềng. Trả (mảng nhãn int32, số vùng).

    Dùng scipy.ndimage.label khi có; không có thì BFS thuần — kết quả giống nhau.
    """
    mask = np.asarray(mask, dtype=bool)
    if _ndi is not None:
        lab, n = _ndi.label(mask)
        return lab.astype(np.int32), int(n)
    H, W = mask.shape
    lab = np.zeros((H, W), dtype=np.int32)
    n = 0
    for y0 in range(H):
        for x0 in range(W):
            if not mask[y0, x0] or lab[y0, x0]:
                continue
            n += 1
            ngan = [(y0, x0)]
            lab[y0, x0] = n
            while ngan:
                y, x = ngan.pop()
                for yy, xx in ((y - 1, x), (y + 1, x), (y, x - 1), (y, x + 1)):
                    if 0 <= yy < H and 0 <= xx < W and mask[yy, xx] and not lab[yy, xx]:
                        lab[yy, xx] = n
                        ngan.append((yy, xx))
    return lab, n


def dem_vung_da(anh, hop, nguong_px):
    """Số vùng da liên thông có diện tích ≥ `nguong_px` nằm trong hộp (l, t, r, b) của ảnh.

    Chỉ xét phần mặt nạ trong hộp (vùng vắt qua mép hộp tính diện tích phần bên trong).
    """
    mask = mat_na_da(anh)
    H, W = mask.shape
    l, t, r, b = hop
    l, t, r, b = max(0, int(l)), max(0, int(t)), min(W, int(r)), min(H, int(b))
    if r <= l or b <= t:
        return 0
    lab, n = gan_nhan_lien_thong(mask[t:b, l:r])
    if n == 0:
        return 0
    dien_tich = np.bincount(lab.ravel(), minlength=n + 1)[1:]
    return int((dien_tich >= nguong_px).sum())


def hop_dau(vung, zoom):
    """Hộp dò mặt cho một khung đã ghép, hoặc None khi không có mảnh đầu.

    = bbox mảnh đầu nới ngang HOP_DAU_NOI_NGANG_X1×zoom mỗi bên; đáy kéo xuống tới
    mép trên mảnh thân + HOP_DAU_PHAN_THAN × cao thân (vai/ngực), không thấp hơn đáy đầu.
    """
    d = vung.get("dau", ("-", None))[1]
    if d is None:
        return None
    th = vung.get("than", ("-", None))[1]
    noi = HOP_DAU_NOI_NGANG_X1 * zoom
    day = d[3]
    if th is not None:
        day = max(day, int(round(th[1] + HOP_DAU_PHAN_THAN * (th[3] - th[1]))))
    return (d[0] - noi, d[1], d[2] + noi, day)


def so_mat(anh, vung, zoom):
    """Số "mặt" (vùng da ≥ NGUONG_MAT_PX_X1×zoom² px) trong hộp đầu; None nếu thiếu mảnh đầu."""
    hop = hop_dau(vung, zoom)
    if hop is None:
        return None
    return dem_vung_da(anh, hop, NGUONG_MAT_PX_X1 * zoom * zoom)


# --- bảng đo + cửa kiểm ---------------------------------------------------------------
def do_khung(vung, cy):
    """Các số đo hình học của một khung từ bbox đục (None khi thiếu mảnh)."""
    d = vung.get("dau", ("-", None))
    th = vung.get("than", ("-", None))
    ch = vung.get("chan", ("-", None))
    return {
        "khe_dt": th[1][1] - d[1][3] if d[1] and th[1] else None,
        "khe_tc": ch[1][1] - th[1][3] if th[1] and ch[1] else None,
        "chan_cy": ch[1][3] - cy if ch[1] else None,
        "rong_than": th[1][2] - th[1][0] if th[1] else None,
        "cao_tong": ch[1][3] - d[1][1] if d[1] and ch[1] else None,
    }


def phan_tich_kiem(chuoi):
    """Phân tích chuỗi --kiem "khoa=a,b;khoa=v;..." → {khoa: (min, max)}.

    Một giá trị = phải bằng đúng; hai giá trị = đoạn [a, b]. Khoá lạ / định dạng sai →
    ValueError với lời giải thích.
    """
    dieu_kien = {}
    for muc in (m.strip() for m in chuoi.split(";")):
        if not muc:
            continue
        if "=" not in muc:
            raise ValueError("--kiem: muc '%s' thieu dau '=' (vd dau-than=-9,0)" % muc)
        khoa, gia_tri = (p.strip() for p in muc.split("=", 1))
        if khoa not in KIEM_COT:
            raise ValueError("--kiem: khoa la '%s'; cho phep: %s" % (khoa, ", ".join(KIEM_COT)))
        try:
            so = [int(x) for x in gia_tri.split(",")]
        except ValueError:
            raise ValueError("--kiem: gia tri '%s' cua %s phai la so nguyen" % (gia_tri, khoa))
        if len(so) == 1:
            lo = hi = so[0]
        elif len(so) == 2:
            lo, hi = min(so), max(so)
        else:
            raise ValueError("--kiem: %s chi nhan 1 hoac 2 so, nhan %d" % (khoa, len(so)))
        dieu_kien[khoa] = (lo, hi)
    if not dieu_kien:
        raise ValueError("--kiem: chuoi rong")
    return dieu_kien


def phan_tich_cf(chuoi, so_cf):
    """Phân tích --cf "0,1" → [0, 1]: cf cần đo/kiểm, giữ thứ tự, không trùng, trong 0..so_cf-1."""
    if chuoi is None or not chuoi.strip():
        raise ValueError("--cf: chuoi rong (vd --cf 0,1)")
    ket = []
    for muc in chuoi.split(","):
        muc = muc.strip()
        if not muc:
            raise ValueError("--cf: co muc rong trong '%s'" % chuoi)
        try:
            cf = int(muc)
        except ValueError:
            raise ValueError("--cf: '%s' khong phai so nguyen" % muc)
        if not 0 <= cf < so_cf:
            raise ValueError("--cf: cf %d ngoai 0..%d" % (cf, so_cf - 1))
        if cf in ket:
            raise ValueError("--cf: cf %d lap" % cf)
        ket.append(cf)
    return ket


def kiem_hang(so_do, dieu_kien):
    """Danh sách lỗi (chuỗi) của một hàng đo so với điều kiện; rỗng = đạt."""
    loi = []
    for khoa, (lo, hi) in dieu_kien.items():
        v = so_do.get(KIEM_COT[khoa])
        if v is None:
            loi.append("%s=? (khong do duoc)" % khoa)
        elif v < lo or v > hi:
            loi.append("%s=%+d ngoai [%d, %d]" % (khoa, v, lo, hi))
    return loi


def main(argv=None):
    ap = argparse.ArgumentParser()
    ap.add_argument("--part-bin", default=MAC_DINH_PART_BIN)
    ap.add_argument("--sprites", default=MAC_DINH_SPRITES)
    ap.add_argument(
        "--sprites-closure",
        default=None,
        help="thu muc closure kieu sprites/<id>/x<zoom>/sprite.png (thay cho --sprites)",
    )
    ap.add_argument("--charinfo", default=MAC_DINH_CHARINFO)
    ap.add_argument("--head", type=int, default=30)
    ap.add_argument("--body", type=int, default=57)
    ap.add_argument("--leg", type=int, default=58)
    ap.add_argument("--zoom", type=int, default=1)
    ap.add_argument(
        "--cf",
        default=None,
        help='chi do/kiem cac cf nay, vd "0,1" (NPC chi dung cf 0/1); mac dinh ca 33',
    )
    ap.add_argument("--out", required=True)
    ap.add_argument(
        "--kiem",
        default=None,
        help='cua kiem, vd "dau-than=-9,0;than-chan=-11,-2;so-mat=1;chan=-11,-10" — '
        "thoat ma 1 va in cf vi pham",
    )
    a = ap.parse_args(argv)
    dieu_kien = None
    if a.kiem is not None:
        try:
            dieu_kien = phan_tich_kiem(a.kiem)
        except ValueError as e:
            ap.error(str(e))

    charinfo = json.load(open(a.charinfo))
    cf_list = list(range(len(charinfo)))
    if a.cf is not None:
        try:
            cf_list = phan_tich_cf(a.cf, len(charinfo))
        except ValueError as e:
            ap.error(str(e))
    parts = doc_part_bin(a.part_bin)
    ids = (a.head, a.body, a.leg)
    kieu_closure = a.sprites_closure is not None
    thu_muc_sprite = a.sprites_closure if kieu_closure else a.sprites
    z = a.zoom
    cx, cy = 48 * z, 96 * z
    o_w, o_h = cx * 2, cy + 16 * z
    cot = min(11, len(cf_list))
    hang = (len(cf_list) + cot - 1) // cot
    sheet = Image.new("RGBA", (o_w * cot, (o_h + 12) * hang), (40, 40, 40, 255))
    ve = ImageDraw.Draw(sheet)
    vi_pham = []
    print("cf | dau(sid) than(sid) chan(sid) | khe dau-than | khe than-chan | chan so voi cy | rong than | cao tong | so mat")
    for vi_tri, cf in enumerate(cf_list):
        o, vung = ghep(cf, charinfo, parts, thu_muc_sprite, ids, z, cx, cy, kieu_closure)
        r, c = divmod(vi_tri, cot)
        ox, oy = c * o_w, r * (o_h + 12)
        # vạch cy (mặt đất) và cx
        ve.line([(ox, oy + cy), (ox + o_w, oy + cy)], fill=(90, 200, 90, 255))
        ve.line([(ox + cx, oy), (ox + cx, oy + o_h)], fill=(70, 70, 120, 255))
        sheet.alpha_composite(o, (ox, oy))
        ve.text((ox + 2, oy + o_h), "cf %d" % cf, fill=(255, 255, 255, 255))
        so_do = do_khung(vung, cy)
        so_do["so_mat"] = so_mat(o, vung, z)
        hop = hop_dau(vung, z)
        if hop is not None:  # vẽ hộp dò mặt: vàng mờ khi đúng 1 mặt, đỏ khi khác 1
            mau = (230, 200, 60, 110) if so_do["so_mat"] == 1 else (255, 60, 60, 220)
            ve.rectangle([ox + hop[0], oy + hop[1], ox + hop[2] - 1, oy + hop[3] - 1], outline=mau)
        d, th, ch = vung.get("dau", ("-", None)), vung.get("than", ("-", None)), vung.get("chan", ("-", None))
        print("%2d | %s %s %s | %s | %s | %s | %s | %s | %s" % (cf, d[0], th[0], ch[0],
              "%+d" % so_do["khe_dt"] if so_do["khe_dt"] is not None else "?",
              "%+d" % so_do["khe_tc"] if so_do["khe_tc"] is not None else "?",
              "%+d" % so_do["chan_cy"] if so_do["chan_cy"] is not None else "?",
              so_do["rong_than"] if so_do["rong_than"] is not None else "?",
              so_do["cao_tong"] if so_do["cao_tong"] is not None else "?",
              so_do["so_mat"] if so_do["so_mat"] is not None else "?"))
        if dieu_kien is not None:
            loi = kiem_hang(so_do, dieu_kien)
            if loi:
                vi_pham.append((cf, loi))
    sheet.save(a.out)
    print("da ghi", a.out)
    if dieu_kien is not None:
        if vi_pham:
            print("KIEM: %d/%d cf VI PHAM (%s)" % (len(vi_pham), len(cf_list), a.kiem))
            for cf, loi in vi_pham:
                print("  cf %2d: %s" % (cf, "; ".join(loi)))
            return 1
        print("KIEM: %d/%d cf dat (%s)" % (len(cf_list), len(cf_list), a.kiem))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
