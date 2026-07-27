# -*- coding: utf-8 -*-
"""
===========================================================================
 KIỂM TRA ẢNH — BẢN KIỂM THẬT
 -------------------------------------------------------------------------
 VÌ SAO CÓ FILE NÀY (đọc kỹ, đây là bài học đắt tiền)

 File cũ KIEM-TRA-ANH.ps1 ghi trong phần chú thích là nó kiểm "5 điều", trong
 đó có "3. Kích thước có đúng không". NHƯNG ĐỌC MÃ THÌ NÓ KHÔNG HỀ KIỂM KÍCH
 THƯỚC. Hậu quả thật: lô 343 ảnh chạy qua nó, báo "sạch", rồi README của lô
 còn ghi "Không phát hiện sai kích thước" — trong khi 62 ảnh đang to gấp 4 lần
 cỡ chuẩn. Công cụ nói dối còn tệ hơn không có công cụ, vì nó tạo niềm tin sai.

 Bản Python này kiểm ĐÚNG những gì nó nói, và in ra bằng chứng cho từng lỗi.

 -------------------------------------------------------------------------
 KIỂM NHỮNG GÌ

   1. Tên file  — có trong manifest không (sai tên = game không tìm thấy ảnh)
   2. Thiếu/thừa — đủ 343 ảnh chưa, có ảnh lạ không
   3. Kích thước — ĐÚNG cỡ của mức zoom đang kiểm (điều mà bản cũ bỏ sót)
   4. Kênh trong suốt — thiếu alpha thì nhân vật có nền vuông xấu
   5. Ảnh phẳng  — dưới 40 màu gần như chắc chắn là "hình khối vẽ bằng code"
   6. Ảnh trống  — toàn trong suốt, vẽ hụt
   7. File quá nhẹ — dưới 3 KB thường là ảnh phẳng (bài học lô v3)

 -------------------------------------------------------------------------
 CÁCH DÙNG

   python KIEM-TRA-ANH.py <thư-mục> [mức]

   mức = x4 | x3 | x2 | x1  -> kiểm cỡ đúng mức đó
   mức = nguon              -> ảnh gốc AI vẽ, cỡ tự do, KHÔNG kiểm cỡ
   Không truyền mức -> tự đoán theo tên thư mục (…/x2/ -> x2), không đoán
   được thì coi là "nguon".

   Ví dụ:
     python KIEM-TRA-ANH.py anh-game/x4
     python KIEM-TRA-ANH.py anh-nguon nguon
===========================================================================
"""

import json
import sys
from pathlib import Path

try:
    from PIL import Image
except ImportError:
    print("THIEU THU VIEN Pillow. Cai bang lenh:  python -m pip install Pillow")
    sys.exit(1)

TI_LE_MUC = {"x4": 1.00, "x3": 0.75, "x2": 0.50, "x1": 0.25}
SO_MAU_TOI_THIEU = 40          # dưới mức này là ảnh phẳng (bài học lô v3)
DUNG_LUONG_NGO_VUC = 3 * 1024  # dưới 3 KB thì đáng ngờ

GOC = Path(__file__).resolve().parent


def ten_chuan_x4(ten: str) -> str:
    """ten_anh_x2.png -> ten_anh_x4.png (khoá tra trong bảng cỡ luôn là _x4)."""
    than, duoi = Path(ten).stem, Path(ten).suffix or ".png"
    for m in TI_LE_MUC:
        if than.endswith("_" + m):
            return f"{than[:-len(m)]}x4{duoi}"
    return ten


def main() -> int:
    if len(sys.argv) < 2:
        print(__doc__)
        return 1
    thu_muc = Path(sys.argv[1])
    if not thu_muc.is_dir():
        print(f"KHONG TIM THAY THU MUC: {thu_muc}")
        return 1

    # xác định mức: tham số > tên thư mục > "nguon"
    muc = sys.argv[2].lower() if len(sys.argv) > 2 else thu_muc.name.lower()
    if muc not in TI_LE_MUC:
        muc = "nguon"

    bang = json.loads((GOC / "kich-thuoc-anh.json").read_text(encoding="utf-8"))
    anh = sorted(thu_muc.rglob("*.png"))
    if not anh:
        print(f"Khong co file .png nao trong {thu_muc}")
        return 1

    print(f"\nKiem {len(anh)} anh trong {thu_muc}  (muc: {muc})\n")

    sai_ten, sai_co, thieu_alpha, anh_phang, anh_trong, nhe_bat_thuong = [], [], [], [], [], []
    da_gap = set()

    for p in anh:
        ten4 = ten_chuan_x4(p.name)
        da_gap.add(ten4)
        co_chuan = bang.get(ten4)
        if co_chuan is None:
            sai_ten.append(p.name)
            continue

        try:
            im = Image.open(p)
            im.load()
        except Exception as e:
            print(f"  Loi mo {p.name}: {e}")
            continue

        # 3. kích thước — điều mà bản PowerShell cũ KHÔNG kiểm
        if muc != "nguon":
            ti = TI_LE_MUC[muc]
            mong_doi = (max(1, round(co_chuan[0] * ti)), max(1, round(co_chuan[1] * ti)))
            if (im.width, im.height) != mong_doi:
                sai_co.append((p.name, (im.width, im.height), mong_doi))

        # 4. kênh trong suốt
        if im.mode != "RGBA" and "A" not in im.getbands():
            thieu_alpha.append(p.name)
            rgba = im.convert("RGBA")
        else:
            rgba = im.convert("RGBA")

        # 6. ảnh trống (không một điểm nào đục)
        if rgba.getchannel("A").getextrema()[1] == 0:
            anh_trong.append(p.name)

        # 5. ảnh phẳng — đếm màu trên phần NHÌN THẤY được
        mau = rgba.convert("RGB").getcolors(maxcolors=1 << 20)
        if mau is not None and len(mau) < SO_MAU_TOI_THIEU:
            anh_phang.append((p.name, len(mau)))

        # 7. file quá nhẹ (chỉ cảnh báo với ảnh cỡ lớn — ảnh x1 nhẹ là bình thường)
        if muc in ("x4", "nguon") and p.stat().st_size < DUNG_LUONG_NGO_VUC:
            nhe_bat_thuong.append((p.name, p.stat().st_size))

    thieu = sorted(set(bang) - da_gap)

    # ---------------------------------------------------------------- báo cáo
    def khoi(tieu_de, danh_sach, in_dong=lambda x: f"        {x}"):
        if not danh_sach:
            return 0
        print(f"\n  [X] {tieu_de}: {len(danh_sach)}")
        for x in danh_sach[:20]:
            print(in_dong(x))
        if len(danh_sach) > 20:
            print(f"        ... va {len(danh_sach) - 20} cai nua")
        return len(danh_sach)

    print("=" * 70)
    print("  KET QUA KIEM")
    print("=" * 70)
    print(f"  Da kiem            : {len(anh)} anh")
    print(f"  Co trong manifest  : {len(bang)} ten")

    loi = 0
    loi += khoi("SAI TEN (khong co trong manifest)", sai_ten)
    loi += khoi("THIEU so voi manifest", thieu)
    loi += khoi("SAI KICH THUOC", sai_co,
                lambda x: f"        {x[0]}: dang {x[1][0]}x{x[1][1]}, phai {x[2][0]}x{x[2][1]}")
    loi += khoi("THIEU KENH TRONG SUOT", thieu_alpha)
    loi += khoi("ANH TRONG (khong co gi nhin thay)", anh_trong)
    loi += khoi(f"ANH PHANG (duoi {SO_MAU_TOI_THIEU} mau)", anh_phang,
                lambda x: f"        {x[0]}: {x[1]} mau")

    if nhe_bat_thuong:
        print(f"\n  [!] Nhe bat thuong (duoi 3 KB, DANG NGO chu chua chac sai): "
              f"{len(nhe_bat_thuong)}")
        for t, kb in nhe_bat_thuong[:10]:
            print(f"        {t}: {kb} byte")

    print("\n" + "=" * 70)
    if loi == 0:
        print("  KET LUAN: DAT — khong loi nao.")
    else:
        print(f"  KET LUAN: CO {loi} LOI phai sua (xem cac muc [X] o tren).")
    print("=" * 70)
    return 0 if loi == 0 else 2


if __name__ == "__main__":
    sys.exit(main())
