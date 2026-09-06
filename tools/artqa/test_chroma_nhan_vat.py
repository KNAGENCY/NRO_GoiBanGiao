"""Khe nền magenta còn sót trên mảnh nhân vật — chốt cách xoá (06/09/2026).

TẠI SAO có file này: trong game thấy VỆT HỒNG nhỏ trên nhân vật (khe giữa cánh tay
và thắt lưng khung đứng, trên đầu Iruka). Truy về nguồn:

* `artqa.extract_alpha` chỉ flood magenta CHẠM MÉP ảnh, nên khe KÍN (nách, giữa dải
  khăn trán và tóc) giữ nguyên magenta với alpha 255;
* `remove_large_chroma_islands` chỉ xoá đảo ≥ 0,038 % diện tích (598 px với ảnh
  nguồn của gói này) — chú thích cũ tin đảo 586 px trên master Iruka là "accent
  bông tai" nên cố ý giữ. Đo lại: đảo đó ở x[488,518] y[329,375], màu (245,3,238)
  — magenta khoá màu nguyên chất, là KHE giữa khăn và tóc, không phải bông tai.

Kết quả: 293 px hồng trên 37 mảnh ×4 zoom. Bản vá hạ ngưỡng cho RIÊNG các dải nhân
vật (`NGUON_NHAN_VAT`) và vá màu đốm lẻ để không thủng áo; đường closure v11 cũ
giữ nguyên tham số mặc định nên bất động từng byte.

Ảnh dựng tay trong file này — không cần kho ảnh nguồn.
"""

import unittest

import numpy as np
from PIL import Image

from tools.artqa.artqa import extract_alpha, is_chroma
from tools.artqa.task6b_closure import (
    KHE_NEN_TOI_THIEU,
    NGUON_NHAN_VAT,
    remove_large_chroma_islands,
)

CHROMA = (246, 4, 242)
AO = (28, 42, 74)  # navy áo Naruto


def _anh_thu(khe: tuple[int, int, int, int] | None, dom: tuple[int, int] | None):
    """Nền chroma chạm mép + thân áo đặc; đục thêm một khe và/hoặc một đốm lẻ."""

    anh = Image.new("RGB", (60, 60), CHROMA)
    for y in range(10, 50):
        for x in range(10, 50):
            anh.putpixel((x, y), AO)
    if khe is not None:
        x0, y0, x1, y1 = khe
        for y in range(y0, y1):
            for x in range(x0, x1):
                anh.putpixel((x, y), CHROMA)
    if dom is not None:
        anh.putpixel(dom, CHROMA)
    return anh


def _dem_chroma(anh: Image.Image) -> int:
    arr = np.asarray(anh.convert("RGBA"))
    return sum(
        1
        for y in range(arr.shape[0])
        for x in range(arr.shape[1])
        if arr[y, x, 3] > 0 and is_chroma(tuple(int(v) for v in arr[y, x, :3]))
    )


class KheNenNhanVatTests(unittest.TestCase):
    def test_mac_dinh_v11_van_giu_khe_nho(self):
        """Đường closure cũ KHÔNG được đổi: khe 9 px dưới sàn 16 px vẫn sống (và vẫn hồng).

        Ảnh thử 60×60 nên ngưỡng v11 rơi xuống sàn `max(16, …)`; trên ảnh nguồn thật
        (1254² / 2172×724) ngưỡng đó là 598 px — cùng một cơ chế, khác con số.
        """
        anh = _anh_thu((20, 20, 23, 23), None)
        ket_qua = remove_large_chroma_islands(anh, extract_alpha(anh))
        self.assertEqual(_dem_chroma(ket_qua), 9)

    def test_ha_nguong_thi_khe_thanh_trong_suot(self):
        anh = _anh_thu((20, 20, 23, 23), None)
        ket_qua = remove_large_chroma_islands(
            anh, extract_alpha(anh), dien_tich_khe=KHE_NEN_TOI_THIEU, va_dom=True
        )
        self.assertEqual(_dem_chroma(ket_qua), 0)
        alpha = np.asarray(ket_qua.getchannel("A"))
        self.assertEqual(int(alpha[20:23, 20:23].max()), 0, "khe phải trong suốt hẳn")

    def test_dom_le_duoc_va_mau_chu_khong_thung_ao(self):
        """Đốm 1 px giữa lòng áo: hạ alpha là đục lỗ kim → phải VÁ MÀU, giữ alpha."""
        anh = _anh_thu(None, (30, 30))
        truoc = np.asarray(extract_alpha(anh).getchannel("A"))
        ket_qua = remove_large_chroma_islands(
            anh, extract_alpha(anh), dien_tich_khe=KHE_NEN_TOI_THIEU, va_dom=True
        )
        sau = np.asarray(ket_qua.getchannel("A"))
        self.assertEqual(_dem_chroma(ket_qua), 0)
        self.assertEqual(int(sau[30, 30]), 255, "đốm phải còn đục")
        self.assertEqual(int((truoc > 0).sum()), int((sau > 0).sum()), "không mất pixel đục nào")
        self.assertEqual(tuple(np.asarray(ket_qua.convert("RGB"))[30, 30]), AO)

    def test_khong_va_dom_thi_dom_o_nguyen(self):
        """`va_dom=False` là hành vi cũ — đốm dưới ngưỡng không bị đụng tới."""
        anh = _anh_thu(None, (30, 30))
        ket_qua = remove_large_chroma_islands(anh, extract_alpha(anh), dien_tich_khe=2)
        self.assertEqual(_dem_chroma(ket_qua), 1)

    def test_nguon_nhan_vat_phu_du_5_dai_nguoi_choi_va_iruka(self):
        """Sót một dải là sót vệt hồng trên đúng những mảnh cắt từ dải đó."""
        from tools.artqa.task6b_closure import IRUKA_SOURCE
        from tools.artqa.nhan_vat_moc import NGUON_NGUOI_CHOI

        self.assertEqual(NGUON_NHAN_VAT, frozenset(NGUON_NGUOI_CHOI) | {IRUKA_SOURCE})


    def test_dom_co_lap_khong_co_mau_sach_thi_dat_alpha_0(self):
        """Nhánh dự phòng của `_va_dom_chroma`: quanh đốm không có màu sạch → đục lỗ.

        Kiểm đối kháng 06/09 chỉ ra nhánh này chưa test nào chạm tới. Dữ liệu thật hiện
        không kích hoạt nó (mọi đốm đều tìm được màu lân cận trong bán kính 11), nhưng
        hành vi phải được chốt: thà một pixel trong suốt còn hơn một chấm magenta chói.
        """

        import numpy as np
        from PIL import Image

        from tools.artqa.task6b_closure import remove_large_chroma_islands

        # Ảnh toàn nền chroma, KHÔNG có một pixel màu sạch nào → nhánh vá màu không tìm
        # được màu thay thế. Alpha dựng tay (đục hết) để đi thẳng vào nhánh cần chốt.
        anh = Image.new("RGB", (25, 25), (246, 4, 242))
        tach = Image.new("RGBA", (25, 25), (246, 4, 242, 255))
        ket_qua = remove_large_chroma_islands(anh, tach, dien_tich_khe=8, va_dom=True)
        alpha = np.asarray(ket_qua.convert("RGBA").getchannel("A"))
        self.assertEqual(int(alpha.max()), 0, "moi pixel chroma phai thanh trong suot")


if __name__ == "__main__":
    unittest.main()
