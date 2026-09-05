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
(đỉnh đầu → bàn chân).
Khe = mép trên vùng ĐỤC của mảnh dưới trừ mép dưới vùng đục của mảnh trên; >0 = HỞ, <0 = chồng.
"""
import argparse
import json
import os
import struct

from PIL import Image, ImageDraw

MAC_DINH_PART_BIN = r"D:\NRO_GoiBanGiao\server\internal\game\assets\part.bin"
MAC_DINH_SPRITES = r"D:\NRO_Project\Assets\Resources\res\x1\smallimage"
MAC_DINH_CHARINFO = os.path.join(os.path.dirname(os.path.abspath(__file__)), "charinfo.json")
SO_KHUNG = {0: 3, 1: 17, 2: 14}


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


def main():
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
    ap.add_argument("--out", required=True)
    a = ap.parse_args()

    charinfo = json.load(open(a.charinfo))
    parts = doc_part_bin(a.part_bin)
    ids = (a.head, a.body, a.leg)
    kieu_closure = a.sprites_closure is not None
    thu_muc_sprite = a.sprites_closure if kieu_closure else a.sprites
    z = a.zoom
    cx, cy = 48 * z, 96 * z
    o_w, o_h = cx * 2, cy + 16 * z
    cot = 11
    hang = (len(charinfo) + cot - 1) // cot
    sheet = Image.new("RGBA", (o_w * cot, (o_h + 12) * hang), (40, 40, 40, 255))
    ve = ImageDraw.Draw(sheet)
    print("cf | dau(sid) than(sid) chan(sid) | khe dau-than | khe than-chan | chan so voi cy | rong than | cao tong")
    for cf in range(len(charinfo)):
        o, vung = ghep(cf, charinfo, parts, thu_muc_sprite, ids, z, cx, cy, kieu_closure)
        r, c = divmod(cf, cot)
        ox, oy = c * o_w, r * (o_h + 12)
        # vạch cy (mặt đất) và cx
        ve.line([(ox, oy + cy), (ox + o_w, oy + cy)], fill=(90, 200, 90, 255))
        ve.line([(ox + cx, oy), (ox + cx, oy + o_h)], fill=(70, 70, 120, 255))
        sheet.alpha_composite(o, (ox, oy))
        ve.text((ox + 2, oy + o_h), "cf %d" % cf, fill=(255, 255, 255, 255))
        d, th, ch = vung.get("dau", ("-", None)), vung.get("than", ("-", None)), vung.get("chan", ("-", None))
        khe_dt = th[1][1] - d[1][3] if d[1] and th[1] else None
        khe_tc = ch[1][1] - th[1][3] if th[1] and ch[1] else None
        chan_cy = ch[1][3] - cy if ch[1] else None
        rong_than = th[1][2] - th[1][0] if th[1] else None
        cao_tong = ch[1][3] - d[1][1] if d[1] and ch[1] else None
        print("%2d | %s %s %s | %s | %s | %s | %s | %s" % (cf, d[0], th[0], ch[0],
              "%+d" % khe_dt if khe_dt is not None else "?",
              "%+d" % khe_tc if khe_tc is not None else "?",
              "%+d" % chan_cy if chan_cy is not None else "?",
              rong_than if rong_than is not None else "?",
              cao_tong if cao_tong is not None else "?"))
    sheet.save(a.out)
    print("da ghi", a.out)


if __name__ == "__main__":
    main()
