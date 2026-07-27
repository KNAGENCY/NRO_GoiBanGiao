# -*- coding: utf-8 -*-
"""
===========================================================================
 TẠO ĐỦ 4 MỨC ZOOM (x4 / x3 / x2 / x1) TỪ ẢNH GỐC AI VẼ
 -------------------------------------------------------------------------
 VÌ SAO CÓ CÔNG CỤ NÀY

 Game gửi ảnh xuống máy người chơi theo 4 cỡ, tuỳ màn hình to nhỏ:
     x4 = 100%   x3 = 75%   x2 = 50%   x1 = 25%
 Thiếu MỘT mức thôi là người chơi ở cỡ đó **mất ảnh trong im lặng** — máy chủ
 cũ nuốt lỗi, không báo gì (xem BAT-DAU-TU-DAY.md phần mapping ảnh).

 AI vẽ ảnh KHÔNG đẹp ở cỡ nhỏ (96×96 ra bết). Nên quy trình đúng là:
     AI vẽ ở cỡ GỐC LỚN  ->  công cụ này thu về đủ 4 mức
 Công cụ nhận ảnh nguồn **cỡ bất kỳ**, tự tra cỡ chuẩn theo tên file trong
 kich-thuoc-anh.json, rồi sinh cả 4 mức.

 -------------------------------------------------------------------------
 ĐIỂM KỸ THUẬT QUAN TRỌNG — vì sao không dùng lại THU-NHO-ANH.ps1

 Ảnh có nền trong suốt lưu màu ở những chỗ trong suốt một cách tuỳ tiện
 (thường là màu đen). Thu nhỏ kiểu thường sẽ TRỘN màu rác đó vào mép hình,
 sinh ra **viền tối bẩn** quanh nhân vật — nhìn rõ nhất khi đặt lên nền sáng.
 Cách chữa đúng: nhân màu với độ đục trước khi thu (premultiply), thu xong
 chia ngược lại. Công cụ này làm vậy; bản PowerShell cũ thì không.

 -------------------------------------------------------------------------
 CÁCH DÙNG

   python TAO-4-MUC-ZOOM.py <thư-mục-ảnh-nguồn> [thư-mục-xuất]

   Ví dụ:
     python TAO-4-MUC-ZOOM.py anh-nguon anh-game

   Không truyền gì thì mặc định: nguồn "anh-nguon", xuất "anh-game".

 KẾT QUẢ
   anh-game/x4/ten_anh_x4.png
   anh-game/x3/ten_anh_x3.png
   anh-game/x2/ten_anh_x2.png
   anh-game/x1/ten_anh_x1.png

 Không bao giờ ghi đè ảnh nguồn. Ảnh không có trong bảng cỡ -> báo, bỏ qua.
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

# 4 mức zoom: tên mức -> tỉ lệ so với cỡ chuẩn (x4 = 100%)
CAC_MUC = [("x4", 1.00), ("x3", 0.75), ("x2", 0.50), ("x1", 0.25)]

GOC_SCRIPT = Path(__file__).resolve().parent
BANG_CO = GOC_SCRIPT / "kich-thuoc-anh.json"


def thu_nho_giu_mep_sach(im: Image.Image, rong: int, cao: int) -> Image.Image:
    """Thu nhỏ ảnh RGBA mà KHÔNG sinh viền bẩn quanh mép.

    Chỗ trong suốt của ảnh vẫn có màu lưu sẵn (thường là đen). Thu nhỏ kiểu
    thường sẽ trộn màu rác đó vào mép hình -> viền tối quanh nhân vật.

    Chế độ "RGBa" của Pillow là ảnh đã NHÂN SẴN màu với độ đục. Đổi sang nó
    trước khi thu, rồi đổi ngược lại sau, chính là cách chữa chuẩn — và chạy
    bằng mã C của thư viện nên nhanh hơn tự tính từng điểm ảnh rất nhiều.
    """
    return (
        im.convert("RGBA")
        .convert("RGBa")                       # nhân màu với độ đục
        .resize((rong, cao), Image.LANCZOS)
        .convert("RGBA")                       # chia ngược lại
    )


def ten_theo_muc(ten_file: str, muc: str) -> str:
    """quai_soi_rung_x4.png + 'x2'  ->  quai_soi_rung_x2.png

    Tên trong manifest luôn kết thúc bằng _x4; đổi hậu tố cho khớp mức, để
    nhìn tên là biết ngay ảnh cỡ nào (không phải dò theo thư mục).
    """
    goc = Path(ten_file).stem
    duoi = Path(ten_file).suffix or ".png"
    for m, _ in CAC_MUC:
        if goc.endswith("_" + m):
            goc = goc[: -(len(m) + 1)]
            break
    return f"{goc}_{muc}{duoi}"


def main() -> int:
    thu_muc_nguon = Path(sys.argv[1]) if len(sys.argv) > 1 else GOC_SCRIPT / "anh-nguon"
    thu_muc_ra = Path(sys.argv[2]) if len(sys.argv) > 2 else GOC_SCRIPT / "anh-game"

    if not BANG_CO.exists():
        print(f"KHONG TIM THAY BANG CO: {BANG_CO}")
        print("File kich-thuoc-anh.json phai nam canh script nay.")
        return 1
    if not thu_muc_nguon.is_dir():
        print(f"KHONG TIM THAY THU MUC ANH NGUON: {thu_muc_nguon}")
        return 1

    bang = json.loads(BANG_CO.read_text(encoding="utf-8"))
    nguon = sorted(p for p in thu_muc_nguon.rglob("*.png") if thu_muc_ra not in p.parents)
    if not nguon:
        print(f"Khong co file .png nao trong {thu_muc_nguon}")
        return 1

    print(f"\nTim thay {len(nguon)} anh nguon. Dang sinh 4 muc zoom...\n")

    xong = 0
    khong_co_trong_bang = []
    canh_bao_le = []      # cỡ không chia hết cho 4 -> mức x1 bị lệch nửa điểm ảnh
    canh_bao_phong_to = []  # ảnh nguồn NHỎ HƠN cỡ chuẩn -> phải phóng to, sẽ mờ

    for duong in nguon:
        ten = duong.name
        # tra bảng theo tên chuẩn (_x4) dù file nguồn đặt tên mức nào
        ten_tra = ten_theo_muc(ten, "x4")
        co = bang.get(ten_tra) or bang.get(ten)
        if co is None:
            khong_co_trong_bang.append(ten)
            continue
        rong_chuan, cao_chuan = int(co[0]), int(co[1])

        try:
            im = Image.open(duong)
            im.load()
        except Exception as e:  # ảnh hỏng thì báo tên, không làm chết cả lô
            print(f"  Loi mo {ten}: {e}")
            continue

        if im.width < rong_chuan or im.height < cao_chuan:
            canh_bao_phong_to.append((ten, im.size, (rong_chuan, cao_chuan)))
        if rong_chuan % 4 or cao_chuan % 4:
            canh_bao_le.append((ten_tra, (rong_chuan, cao_chuan)))

        for muc, ti_le in CAC_MUC:
            r = max(1, round(rong_chuan * ti_le))
            c = max(1, round(cao_chuan * ti_le))
            ra = thu_muc_ra / muc
            ra.mkdir(parents=True, exist_ok=True)
            anh_moi = thu_nho_giu_mep_sach(im, r, c)
            anh_moi.save(ra / ten_theo_muc(ten, muc), "PNG", optimize=True)
        xong += 1
        if xong % 25 == 0:
            print(f"  ...{xong}/{len(nguon)}")

    print("\n" + "=" * 70)
    print("  KET QUA")
    print("=" * 70)
    print(f"  Da sinh du 4 muc cho : {xong}/{len(nguon)} anh")
    print(f"  Tong file xuat ra    : {xong * 4}")
    print(f"  Thu muc xuat         : {thu_muc_ra}")

    if khong_co_trong_bang:
        print(f"\n  [!] {len(khong_co_trong_bang)} anh KHONG CO trong bang co "
              f"(sai ten so voi manifest) — da bo qua:")
        for t in khong_co_trong_bang[:20]:
            print(f"        {t}")
        if len(khong_co_trong_bang) > 20:
            print(f"        ... va {len(khong_co_trong_bang) - 20} file nua")

    if canh_bao_phong_to:
        print(f"\n  [!] {len(canh_bao_phong_to)} anh NGUON NHO HON co chuan -> phai "
              f"phong to, anh se MO. Nen ve lai o co lon hon:")
        for t, ngu, ch in canh_bao_phong_to[:10]:
            print(f"        {t}: nguon {ngu[0]}x{ngu[1]} < chuan {ch[0]}x{ch[1]}")

    if canh_bao_le:
        print(f"\n  [i] {len(set(canh_bao_le))} anh co co chuan khong chia het cho 4 "
              f"-> muc x1 lam tron, lech nua diem anh (thuong vo hai, nhung neu "
              f"day la dai nhieu khung hinh thi cac khung se lech):")
        for t, ch in list(dict.fromkeys(canh_bao_le))[:10]:
            print(f"        {t}: {ch[0]}x{ch[1]}")

    print("\n  => Buoc tiep: chay KIEM-TRA-ANH.py tren thu muc xuat de kiem that.")
    print("=" * 70)
    return 0


if __name__ == "__main__":
    sys.exit(main())
