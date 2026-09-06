# -*- coding: utf-8 -*-
"""Kiểm bộ dò "hai khuôn mặt" (cột `so mat`) và cửa `--kiem` của tools/nhan-vat/ghep-khung.py.

Mọi ảnh trong test là ẢNH TỔNG HỢP tự dựng (khối màu da / navy / cam / áo) — không đụng
sprite thật (đang được sinh lại bởi việc khác) nên chạy được ở mọi checkout, không cần
D:\\NRO_Project. Phần đầu-cuối (`main --kiem`) dựng part.bin + charinfo + 3 sprite giả
trong thư mục tạm, đúng định dạng nhị phân client đọc.

ghep-khung.py có dấu gạch nối trong tên nên nạp bằng importlib (như
thu-nghiem-khe-dau-than/toi-uu-dxdy.py và assetpack/tests/test_hinh_hoc_nhan_vat.py).
"""
import contextlib
import importlib.util
import io
import json
import struct
import tempfile
import unittest
from pathlib import Path

import numpy as np
from PIL import Image, ImageDraw

ROOT = Path(__file__).resolve().parents[2]
GHEP_KHUNG = ROOT / "tools" / "nhan-vat" / "ghep-khung.py"

# Màu mẫu. DA: h≈29°, s≈0,44, v≈0,91 → trong ngưỡng da (h 20–48, s 0,35–0,8, v>0,55).
DA = (232, 180, 130, 255)
NAVY = (40, 50, 110, 255)      # băng đô / quần / tay áo (h≈231°)
CAM_DAI = (200, 60, 30, 255)   # đai lưng cam (h≈10,6° — sát ngưỡng dưới, phải bị loại)
XANH_AO = (60, 90, 50, 255)    # áo


def _nap():
    spec = importlib.util.spec_from_file_location("ghep_khung", GHEP_KHUNG)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def _anh(rong, cao, khoi):
    """Ảnh RGBA trong suốt rong×cao, tô các khối (l, t, r, b, màu) — r/b không kể."""
    im = Image.new("RGBA", (rong, cao), (0, 0, 0, 0))
    ve = ImageDraw.Draw(im)
    for l, t, r, b, mau in khoi:
        ve.rectangle([l, t, r - 1, b - 1], fill=mau)
    return im


class BoDoHaiMatTests(unittest.TestCase):
    """dem_vung_da / mat_na_da / hop_dau / so_mat trên ảnh tổng hợp."""

    @classmethod
    def setUpClass(cls):
        cls.g = _nap()

    def test_hai_khoi_da_tach_roi_dem_duoc_2(self):
        # hai khối 6×6 = 36 px (≥ 16), cách nhau 3 px → hai "mặt"
        im = _anh(40, 40, [(5, 5, 11, 11, DA), (5, 14, 11, 20, DA)])
        self.assertEqual(self.g.dem_vung_da(im, (0, 0, 40, 40), self.g.NGUONG_MAT_PX_X1), 2)

    def test_mot_khoi_da_dem_duoc_1(self):
        im = _anh(40, 40, [(5, 5, 15, 15, DA)])
        self.assertEqual(self.g.dem_vung_da(im, (0, 0, 40, 40), self.g.NGUONG_MAT_PX_X1), 1)

    def test_hai_khoi_da_cham_nhau_la_mot_vung(self):
        # chạm cạnh (4 láng giềng) → một vùng liên thông
        im = _anh(40, 40, [(5, 5, 11, 11, DA), (5, 11, 11, 17, DA)])
        self.assertEqual(self.g.dem_vung_da(im, (0, 0, 40, 40), self.g.NGUONG_MAT_PX_X1), 1)

    def test_khoi_da_nho_hon_nguong_khong_tinh_la_mat(self):
        # nắm tay / ngón chân ≤ 8 px ở x1 (đo 05/09) → khối 3×3 = 9 px phải bị bỏ qua
        im = _anh(40, 40, [(5, 5, 15, 15, DA), (25, 25, 28, 28, DA)])
        self.assertEqual(self.g.dem_vung_da(im, (0, 0, 40, 40), self.g.NGUONG_MAT_PX_X1), 1)
        im2 = _anh(40, 40, [(25, 25, 28, 28, DA)])
        self.assertEqual(self.g.dem_vung_da(im2, (0, 0, 40, 40), self.g.NGUONG_MAT_PX_X1), 0)

    def test_nguong_nhan_zoom_binh_phuong(self):
        # x4: một mặt ~16×16 = 256 px trở lên; 10×10 = 100 px (nắm tay x4) không tính
        vung = {"dau": (1, (0, 0, 160, 160)), "than": (2, (0, 150, 160, 240))}
        nho = _anh(160, 240, [(40, 40, 50, 50, DA)])
        lon = _anh(160, 240, [(40, 40, 60, 60, DA)])
        self.assertEqual(self.g.so_mat(nho, vung, 4), 0)
        self.assertEqual(self.g.so_mat(lon, vung, 4), 1)

    def test_mau_khong_phai_da_khong_lot_mat_na(self):
        # băng đô/quần navy, đai lưng cam, áo xanh: khối to nhưng không phải da → 0 mặt
        im = _anh(40, 40, [(2, 2, 14, 14, NAVY), (16, 2, 28, 14, CAM_DAI), (2, 20, 28, 38, XANH_AO)])
        self.assertEqual(self.g.dem_vung_da(im, (0, 0, 40, 40), self.g.NGUONG_MAT_PX_X1), 0)
        mask = self.g.mat_na_da(im)
        self.assertFalse(mask.any())

    def test_mat_na_da_bo_pixel_trong_suot(self):
        im = _anh(10, 10, [(0, 0, 10, 10, DA)])
        im.putalpha(0)
        self.assertFalse(self.g.mat_na_da(im).any())
        im2 = _anh(10, 10, [(0, 0, 10, 10, DA)])
        self.assertTrue(self.g.mat_na_da(im2).all())

    def test_chi_dem_trong_hop(self):
        # khối da ngoài hộp (ở dưới, chỗ đai lưng/bàn tay hạ) không được tính
        im = _anh(40, 60, [(5, 5, 15, 15, DA), (5, 45, 15, 55, DA)])
        self.assertEqual(self.g.dem_vung_da(im, (0, 0, 40, 30), self.g.NGUONG_MAT_PX_X1), 1)
        self.assertEqual(self.g.dem_vung_da(im, (0, 0, 40, 60), self.g.NGUONG_MAT_PX_X1), 2)
        # hộp vượt ngoài ảnh vẫn phải an toàn
        self.assertEqual(self.g.dem_vung_da(im, (-5, -5, 100, 100), self.g.NGUONG_MAT_PX_X1), 2)
        self.assertEqual(self.g.dem_vung_da(im, (30, 30, 20, 20), self.g.NGUONG_MAT_PX_X1), 0)

    def test_hop_dau_noi_ngang_va_keo_xuong_60_phan_tram_than(self):
        vung = {"dau": (1, (10, 10, 40, 36)), "than": (2, (12, 30, 42, 50))}
        # đáy = max(36, 30 + 0,6×20 = 42) = 42 ; nới ngang 2 px x1
        self.assertEqual(self.g.hop_dau(vung, 1), (8, 10, 42, 42))
        # zoom 2: nới 4 px
        self.assertEqual(self.g.hop_dau(vung, 2), (6, 10, 44, 42))
        # thân cao mà nằm trên cao (đáy thân < đáy đầu) → giữ đáy đầu
        vung2 = {"dau": (1, (10, 10, 40, 36)), "than": (2, (12, 12, 42, 20))}
        self.assertEqual(self.g.hop_dau(vung2, 1), (8, 10, 42, 36))
        # thiếu thân → chỉ bbox đầu ; thiếu đầu → None
        self.assertEqual(self.g.hop_dau({"dau": (1, (10, 10, 40, 36))}, 1), (8, 10, 42, 36))
        self.assertIsNone(self.g.hop_dau({"than": (2, (12, 30, 42, 50))}, 1))
        self.assertIsNone(self.g.hop_dau({"dau": ("THIEU 27041", None)}, 1))

    def test_so_mat_bat_mat_thu_hai_tren_nguc(self):
        # dựng đúng tình huống hiện trạng: mặt thật trong bbox đầu, mặt thứ hai ở phần trên thân
        vung = {"dau": (1, (30, 20, 62, 46)), "than": (2, (30, 44, 62, 70))}
        mot = _anh(96, 112, [(40, 30, 56, 40, DA), (30, 60, 62, 70, CAM_DAI)])
        hai = _anh(96, 112, [(40, 30, 56, 40, DA), (42, 50, 56, 58, DA), (30, 60, 62, 70, CAM_DAI)])
        self.assertEqual(self.g.so_mat(mot, vung, 1), 1)
        self.assertEqual(self.g.so_mat(hai, vung, 1), 2)
        self.assertIsNone(self.g.so_mat(hai, {"than": (2, (30, 44, 62, 70))}, 1))

    def test_gan_nhan_bfs_khop_scipy(self):
        # nhánh dự phòng (không scipy) phải cho cùng số vùng và cùng phân bố diện tích
        rng = np.random.default_rng(5)
        mask = rng.random((30, 40)) > 0.55
        scipy_co = self.g._ndi is not None
        lab1, n1 = self.g.gan_nhan_lien_thong(mask)
        cu = self.g._ndi
        try:
            self.g._ndi = None
            lab2, n2 = self.g.gan_nhan_lien_thong(mask)
        finally:
            self.g._ndi = cu
        self.assertEqual(n2, n1)
        self.assertEqual(sorted(np.bincount(lab1.ravel())[1:]), sorted(np.bincount(lab2.ravel())[1:]))
        self.assertTrue((lab2 > 0).tolist() == mask.tolist())
        if not scipy_co:  # pragma: no cover
            self.skipTest("khong co scipy: chi kiem BFS tu nhat quan")


def _ghi_part_bin(duong, parts):
    """parts: [(type, [(sid, dx, dy), ...]), ...] → nhị phân client (>h n; mỗi part: b type + >hbb*k)."""
    d = bytearray(struct.pack(">h", len(parts)))
    for t, fr in parts:
        d += bytes([t])
        for sid, dx, dy in fr:
            d += struct.pack(">hbb", sid, dx, dy)
    Path(duong).write_bytes(bytes(d))


def _dung_bo_thu(thu_muc, dy_than, than_co_mat):
    """Dựng charinfo 1 cf + part.bin 3 part (id 0 đầu/1 thân/2 chân) + 3 sprite giả.

    cf 0: đầu idx0 dy 40, chân idx0 dy 10, thân idx0 dy 20 (zoom 1, cx 48, cy 96):
      đầu 16×12 → hàng 56..68 ; thân 16×14 → mép trên 76 + dy_than ; chân 16×10 → 86..96.
    dy_than=0 → khe đầu–thân +8 (hở) ; dy_than=-10 → -2 (chồng 2 px).
    than_co_mat=True → mảnh thân mang một khối da 8×8 (mặt thứ hai) ở phần trên.
    """
    thu_muc = Path(thu_muc)
    (thu_muc / "charinfo.json").write_text(json.dumps([[[0, 0, 40], [0, 0, 10], [0, 0, 20]]]), encoding="utf-8")
    g = _nap()
    _ghi_part_bin(thu_muc / "part.bin", [
        (0, [(1, 0, 0)] * g.SO_KHUNG[0]),
        (1, [(2, 0, dy_than)] * g.SO_KHUNG[1]),
        (2, [(3, 0, 0)] * g.SO_KHUNG[2]),
    ])
    _anh(16, 12, [(0, 0, 16, 3, (30, 20, 20, 255)), (0, 3, 16, 12, DA)]).save(thu_muc / "Small1.png")
    khoi_than = [(0, 0, 16, 14, XANH_AO)]
    if than_co_mat:
        khoi_than.append((4, 0, 12, 8, DA))
    _anh(16, 14, khoi_than).save(thu_muc / "Small2.png")
    _anh(16, 10, [(0, 0, 16, 10, NAVY)]).save(thu_muc / "Small3.png")
    return thu_muc


def _chay_main(g, thu_muc, kiem, cf=None):
    argv = [
        "--part-bin", str(thu_muc / "part.bin"), "--sprites", str(thu_muc),
        "--charinfo", str(thu_muc / "charinfo.json"), "--head", "0", "--body", "1", "--leg", "2",
        "--out", str(thu_muc / "kq.png"),
    ]
    if kiem is not None:
        argv += ["--kiem", kiem]
    if cf is not None:
        argv += ["--cf", cf]
    buf = io.StringIO()
    with contextlib.redirect_stdout(buf), contextlib.redirect_stderr(io.StringIO()):
        ma = g.main(argv)
    return ma, buf.getvalue()


class CuaKiemTests(unittest.TestCase):
    """--kiem: phân tích chuỗi, chấm từng hàng, mã thoát của main."""

    @classmethod
    def setUpClass(cls):
        cls.g = _nap()

    def test_phan_tich_kiem_doan_va_gia_tri_don(self):
        dk = self.g.phan_tich_kiem("dau-than=-9,0;than-chan=-11,-2;so-mat=1;chan=-11,-10")
        self.assertEqual(dk, {"dau-than": (-9, 0), "than-chan": (-11, -2), "so-mat": (1, 1), "chan": (-11, -10)})
        # thứ tự đảo vẫn ra đoạn đúng; khoảng trắng và dấu ; thừa được bỏ qua
        self.assertEqual(self.g.phan_tich_kiem(" cao=65,50 ; "), {"cao": (50, 65)})

    def test_phan_tich_kiem_tu_choi_khoa_la_va_dinh_dang_sai(self):
        for xau in ("mat=1", "dau-than", "dau-than=a,b", "dau-than=1,2,3", "", ";"):
            with self.subTest(xau=xau), self.assertRaises(ValueError):
                self.g.phan_tich_kiem(xau)

    def test_kiem_hang_bao_dung_o_vi_pham(self):
        dk = self.g.phan_tich_kiem("dau-than=-9,0;than-chan=-11,-2;so-mat=1;chan=-11,-10")
        dat = {"khe_dt": -3, "khe_tc": -5, "chan_cy": -11, "so_mat": 1}
        self.assertEqual(self.g.kiem_hang(dat, dk), [])
        hong = {"khe_dt": 5, "khe_tc": -5, "chan_cy": None, "so_mat": 2}
        loi = self.g.kiem_hang(hong, dk)
        self.assertEqual(len(loi), 3)
        self.assertTrue(any(l.startswith("dau-than=+5") for l in loi))
        self.assertTrue(any(l.startswith("so-mat=+2") for l in loi))
        self.assertTrue(any(l.startswith("chan=?") for l in loi))

    def test_main_kiem_thoat_khac_0_khi_cf_ho_va_hai_mat(self):
        with tempfile.TemporaryDirectory() as tam:
            thu_muc = _dung_bo_thu(tam, dy_than=0, than_co_mat=True)
            ma, out = _chay_main(self.g, thu_muc, "dau-than=-9,0;than-chan=-11,-2;so-mat=1;chan=-11,-10")
            self.assertEqual(ma, 1, out)
            self.assertIn("KIEM: 1/1 cf VI PHAM", out)
            self.assertIn("cf  0:", out)
            self.assertIn("dau-than=+8", out)
            self.assertIn("so-mat=+2", out)
            self.assertIn("chan=+0", out)
            dong_vi_pham = [l for l in out.splitlines() if l.startswith("  cf  0:")][0]
            self.assertNotIn("than-chan", dong_vi_pham)  # 86 - 90 = -4 nằm trong [-11, -2]
            self.assertEqual(dong_vi_pham.count(";"), 2)   # đúng 3 ô vi phạm
            self.assertTrue((thu_muc / "kq.png").is_file())

    def test_main_kiem_thoat_0_khi_dat(self):
        with tempfile.TemporaryDirectory() as tam:
            thu_muc = _dung_bo_thu(tam, dy_than=-10, than_co_mat=False)
            ma, out = _chay_main(self.g, thu_muc, "dau-than=-9,0;so-mat=1;chan=0")
            self.assertEqual(ma, 0, out)
            self.assertIn("KIEM: 1/1 cf dat", out)

    def test_khong_kiem_thi_luon_thoat_0_va_van_in_so_mat(self):
        with tempfile.TemporaryDirectory() as tam:
            thu_muc = _dung_bo_thu(tam, dy_than=0, than_co_mat=True)
            ma, out = _chay_main(self.g, thu_muc, None)
            self.assertEqual(ma, 0, out)
            self.assertNotIn("KIEM:", out)
            hang = [l for l in out.splitlines() if l.startswith(" 0 |")][0]
            self.assertTrue(hang.rstrip().endswith("| 2"), hang)

    def test_phan_tich_cf_danh_sach_va_tu_choi_sai(self):
        self.assertEqual(self.g.phan_tich_cf("0,1", 33), [0, 1])
        self.assertEqual(self.g.phan_tich_cf(" 5 ", 33), [5])
        self.assertEqual(self.g.phan_tich_cf("3,1,2", 33), [3, 1, 2])  # giữ thứ tự người dùng
        for xau in ("", "a", "0,,1", "33", "-1", "0,0"):
            with self.subTest(xau=xau), self.assertRaises(ValueError):
                self.g.phan_tich_cf(xau, 33)

    def test_main_cf_chi_do_va_kiem_cac_cf_duoc_chon(self):
        """NPC chỉ dùng cf 0/1 (Npc.cs:179–185): --cf "0,1" áp cửa kiểm cho đúng 2 cf ấy.
        Dựng charinfo 2 cf: cf 0 đạt, cf 1 hở đầu–thân 8 px → --cf 0 thoát 0, --cf 0,1 thoát 1."""
        with tempfile.TemporaryDirectory() as tam:
            thu_muc = _dung_bo_thu(tam, dy_than=-10, than_co_mat=False)
            (thu_muc / "charinfo.json").write_text(
                json.dumps([[[0, 0, 40], [0, 0, 10], [0, 0, 20]], [[0, 0, 40], [0, 0, 10], [0, 0, 10]]]),
                encoding="utf-8",
            )
            kiem = "dau-than=-9,0;so-mat=1;chan=0"
            ma, out = _chay_main(self.g, thu_muc, kiem, cf="0")
            self.assertEqual(ma, 0, out)
            self.assertIn("KIEM: 1/1 cf dat", out)
            self.assertFalse(any(l.startswith(" 1 |") for l in out.splitlines()), out)
            ma, out = _chay_main(self.g, thu_muc, kiem, cf="0,1")
            self.assertEqual(ma, 1, out)
            self.assertIn("KIEM: 1/2 cf VI PHAM", out)
            self.assertIn("cf  1:", out)
            self.assertTrue(any(l.startswith(" 0 |") for l in out.splitlines()), out)
            # không --cf → đo cả 2 cf, cf 1 vi phạm
            ma, out = _chay_main(self.g, thu_muc, kiem)
            self.assertEqual(ma, 1, out)
            self.assertIn("KIEM: 1/2 cf VI PHAM", out)
            # --cf sai → lỗi tham số (argparse thoát 2)
            with self.assertRaises(SystemExit) as cm:
                _chay_main(self.g, thu_muc, kiem, cf="7")
            self.assertEqual(cm.exception.code, 2)

    def test_bang_giu_cot_cu_va_them_cot_so_mat_o_cuoi(self):
        # script khác (kiem-bang-khe.py cũ, chuan-hoa-dai-thu.py) đọc 7 ô đầu bằng regex — phải còn khớp
        import re
        with tempfile.TemporaryDirectory() as tam:
            thu_muc = _dung_bo_thu(tam, dy_than=0, than_co_mat=True)
            _, out = _chay_main(self.g, thu_muc, None)
        dong = out.splitlines()
        self.assertTrue(dong[0].startswith(
            "cf | dau(sid) than(sid) chan(sid) | khe dau-than | khe than-chan | chan so voi cy | rong than | cao tong"))
        self.assertTrue(dong[0].endswith("| so mat"))
        m = re.match(r"\s*(\d+)\s*\|\s*(\d+)\s+(\d+)\s+(\d+)\s*\|\s*([+-]?\d+)\s*\|\s*([+-]?\d+)\s*\|\s*([+-]?\d+)", dong[1])
        self.assertIsNotNone(m, dong[1])
        self.assertEqual([int(x) for x in m.groups()], [0, 1, 2, 3, 8, -4, 0])


if __name__ == "__main__":
    unittest.main()
