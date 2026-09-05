# -*- coding: utf-8 -*-
"""chep-sprite-client.py — chép 34 sprite nhân vật mới sang client Unity.

Công cụ này KHÔNG được gọi tự động. Client Unity (`D:\\NRO_Project`) nằm ngoài
repo và thường đang mở Play — ghi đè ảnh lúc đó là Unity reimport giữa chừng.

Mặc định chỉ IN ra sẽ làm gì (`--dry-run`). Chỉ khi thêm `--that` mới ghi thật:

    python tools/nhan-vat/chep-sprite-client.py            # xem trước
    python tools/nhan-vat/chep-sprite-client.py --that     # ghi thật (Unity đã dừng Play)

Trước khi ghi, bản cũ trong client được sao lưu vào
`docs/my-thuat/lang-la-client-assets/smallimage-truoc-05-09/x{n}/` (trong worktree,
không đụng client) để còn đường lùi.
"""

from __future__ import annotations

import argparse
import shutil
from pathlib import Path

GOC = Path(__file__).resolve().parents[2]
NGUON = GOC / "docs/my-thuat/lang-la-client-assets/res"
SAO_LUU = GOC / "docs/my-thuat/lang-la-client-assets/smallimage-truoc-05-09"
CLIENT = Path(r"D:\NRO_Project\Assets\Resources\res")
ZOOMS = ("x1", "x2", "x3", "x4")
DAI_ID = range(27041, 27075)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--nguon", type=Path, default=NGUON)
    parser.add_argument("--client", type=Path, default=CLIENT)
    parser.add_argument("--sao-luu", type=Path, default=SAO_LUU)
    parser.add_argument(
        "--that",
        action="store_true",
        help="ghi that vao client (mac dinh chi in ra, khong ghi gi)",
    )
    args = parser.parse_args()
    ghi = args.that

    thieu = []
    chep = 0
    luu = 0
    for zoom in ZOOMS:
        for sprite_id in DAI_ID:
            ten = f"Small{sprite_id}.png"
            nguon = args.nguon / zoom / "smallimage" / ten
            dich = args.client / zoom / "smallimage" / ten
            if not nguon.is_file():
                thieu.append(str(nguon))
                continue
            if dich.is_file():
                luu_vao = args.sao_luu / zoom / ten
                if ghi:
                    luu_vao.parent.mkdir(parents=True, exist_ok=True)
                    shutil.copyfile(dich, luu_vao)
                luu += 1
            if ghi:
                dich.parent.mkdir(parents=True, exist_ok=True)
                shutil.copyfile(nguon, dich)
            chep += 1

    if thieu:
        print("THIEU %d file nguon, vi du: %s" % (len(thieu), thieu[0]))
        print("Chay truoc: python tools/nhan-vat/sinh-lai-sprite.py")
        return 1
    print("nguon : %s" % args.nguon)
    print("dich  : %s" % args.client)
    print("sao luu ban cu -> %s" % args.sao_luu)
    print("se chep %d file (%d file cu duoc sao luu)" % (chep, luu))
    if not ghi:
        print()
        print("DAY LA XEM TRUOC — chua ghi gi. Dung Play trong Unity roi chay lai voi --that.")
    else:
        print("DA GHI. Mo Unity, cho reimport xong roi vao game kiem.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
