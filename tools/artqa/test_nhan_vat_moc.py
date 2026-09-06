# -*- coding: utf-8 -*-
"""Test cho ``tools/artqa/nhan_vat_moc.py`` — mốc giải phẫu đầu nhân vật (V1 05/09/2026).

Vì sao có test này: chẩn đoán "đầu tách thân" hoá ra là HAI KHUÔN MẶT — mảnh thân cắt ở
28 % chiều cao khung chứa cả mặt. Sửa bằng dữ liệu đòi hỏi đo được CẰM và CỠ ĐẦU từng
khung nguồn. Phép đo này từng nằm trong script thử nghiệm
(``tools/nhan-vat/thu-nghiem-khe-dau-than/do-moc-giai-phau.py``); nay là module import
được, nên phải có test khoá: (1) toán median 6 thước cứng, (2) phép tách màu trên ảnh
tổng hợp biết trước đáp số, (3) đối chiếu với số đã đo trong hồ sơ chẩn đoán.
"""

import json
import unittest
from pathlib import Path

import numpy as np

ROOT = Path(__file__).resolve().parents[2]
ART_ROOT = ROOT / "docs/my-thuat/anh-nguon/vertical-slice-v1"
REVIEW = ROOT / "docs/my-thuat/review/vertical-slice-v1.json"
HO_SO = ROOT / "docs/my-thuat/review/khe-dau-than-2026-09-05/moc-giai-phau-do.json"


def _idle_gia():
    """Trung bình 3 khung idle lấy đúng từ hồ sơ chẩn đoán (làm tròn để đọc)."""
    return {
        "face_w": 167.333,
        "face_h": 96.667,
        "eye_h": 33.667,
        "eye_to_chin": 64.5,
        "head_h_hair_chin": 268.0,
        "band_thick": 25.5,
    }


class HeSoTiLeTests(unittest.TestCase):
    def test_he_so_dau_la_median_6_thuoc_cung(self):
        from tools.artqa.nhan_vat_moc import he_so_dau

        # run0 trong hồ sơ: 148/85/31/58.5/243/24 → tỉ lệ 0.8845 0.8793 0.9208 0.9070
        # 0.9067 0.9412 → median = (0.9067 + 0.9070)/2 ≈ 0.90685 (geometry-V1.json).
        khung = {
            "face_w": 148, "face_h": 85, "eye_h": 31, "eye_to_chin": 58.5,
            "head_h_hair_chin": 243, "band_thick": 24,
        }
        r = he_so_dau(khung, _idle_gia())
        self.assertAlmostEqual(r, 0.90685, places=4)

    def test_he_so_dau_bo_thuoc_bi_loai(self):
        from tools.artqa.nhan_vat_moc import he_so_dau

        khung = {
            "face_w": 100, "face_h": 1, "eye_h": 100, "eye_to_chin": 1,
            "head_h_hair_chin": 100, "band_thick": 1,
        }
        idle = {k: 100.0 for k in khung}
        # Không loại: median của [0.01,0.01,0.01,1,1,1] = (0.01+1)/2
        self.assertAlmostEqual(he_so_dau(khung, idle), 0.505)
        # Loại 2 thước hỏng → còn [0.01,1,1,1] → median = 1.0 (một thước hỏng lẻ không kéo lệch)
        self.assertAlmostEqual(he_so_dau(khung, idle, loai={"face_h", "eye_to_chin"}), 1.0)

    def test_khung_khong_co_mat_tra_none(self):
        from tools.artqa.nhan_vat_moc import he_so_dau

        self.assertIsNone(he_so_dau({"content_h": 235}, _idle_gia()))

    def test_ti_le_khung_bang_ti_le_goc_chia_he_so(self):
        from tools.artqa.nhan_vat_moc import ti_le_khung

        self.assertAlmostEqual(ti_le_khung(128 / 340, 0.5), 2 * 128 / 340)
        self.assertAlmostEqual(ti_le_khung(128 / 340, None), 128 / 340)

    def test_cam_khung_uoc_luong_khi_tay_che_cam(self):
        from tools.artqa.nhan_vat_moc import cam_khung

        idle = _idle_gia()
        # bình thường: cằm = chin_y đo được
        self.assertEqual(cam_khung({"chin_y": 332.0, "face_top": 233.0}, idle, 1.0, bi_che=False), 332.0)
        # tay che cằm (melee1): cằm = đỉnh mặt + cao mặt idle × hệ số
        self.assertAlmostEqual(
            cam_khung({"chin_y": 334.0, "face_top": 266.0}, idle, 0.8522, bi_che=True),
            266.0 + 96.667 * 0.8522,
            places=3,
        )
        self.assertIsNone(cam_khung({"content_h": 1}, idle, 1.0, bi_che=False))


def _khung_tong_hop():
    """Nhân vật vẽ bằng hình chữ nhật trên nền chroma; đáp số biết trước."""
    H, W = 400, 300
    rgb = np.zeros((H, W, 3), dtype=np.uint8)
    rgb[:, :] = (246, 4, 242)                 # nền chroma
    rgb[40:80, 100:200] = (40, 30, 20)        # tóc
    rgb[80:95, 100:200] = (30, 40, 110)       # băng đô navy (dày 15)
    rgb[95:195, 110:190] = (230, 180, 130)    # mặt da 80×100, cằm y=195
    rgb[120:130, 130:150] = (0, 0, 0)         # mắt (lỗ trong khối da) cao 10
    rgb[195:330, 105:195] = (30, 40, 110)     # thân navy
    rgb[330:380, 120:180] = (60, 60, 60)      # chân, đáy y=380
    return rgb


class DoKhungTongHopTests(unittest.TestCase):
    def test_do_khung_tim_dung_mat_cam_mat_bang_do(self):
        from tools.artqa.nhan_vat_moc import do_khung

        r = do_khung(_khung_tong_hop(), "tong-hop")
        self.assertEqual(r["content_top"], 40)
        self.assertEqual(r["content_bottom"], 380)
        self.assertEqual(r["face_w"], 80)
        self.assertEqual(r["face_h"], 100)
        self.assertEqual(r["chin_y"], 195)
        self.assertEqual(r["face_top"], 95)
        self.assertAlmostEqual(r["face_cx"], 150.0)
        self.assertEqual(r["eye_h"], 10)
        self.assertAlmostEqual(r["eye_to_chin"], 70.0)
        self.assertEqual(r["hair_top"], 40)
        self.assertEqual(r["head_h_hair_chin"], 155)
        self.assertAlmostEqual(r["band_thick"], 15.0)

    def test_khung_khong_co_da_chi_tra_hop_noi_dung(self):
        from tools.artqa.nhan_vat_moc import do_khung

        rgb = np.zeros((50, 40, 3), dtype=np.uint8)
        rgb[:, :] = (246, 4, 242)
        rgb[10:30, 5:35] = (30, 40, 110)
        r = do_khung(rgb, "nam")
        self.assertEqual(r["content_top"], 10)
        self.assertEqual(r["content_bottom"], 30)
        self.assertNotIn("face_w", r)
        self.assertNotIn("chin_y", r)


@unittest.skipUnless(
    (ART_ROOT / "player_idle_3f_chroma.png").is_file() and HO_SO.is_file(),
    "cần kho ảnh nguồn + hồ sơ chẩn đoán",
)
class DoNguonThatTests(unittest.TestCase):
    """Port phải cho ra ĐÚNG số đã đo trong hồ sơ 05/09 (moc-giai-phau-do.json)."""

    @classmethod
    def setUpClass(cls):
        from tools.artqa.nhan_vat_moc import do_nguon
        from tools.artqa.task6b_closure import _review_by_name

        review = json.loads(REVIEW.read_text(encoding="utf-8"))
        cls.moc = do_nguon(ART_ROOT, _review_by_name(review))
        cls.ho_so = {
            (r["strip"], r["i"]): r
            for r in json.loads(HO_SO.read_text(encoding="utf-8"))
        }

    def test_du_17_khung_nguoi_choi(self):
        self.assertEqual(len(self.moc), 17)

    def test_khop_ho_so_tung_khung(self):
        tag = {
            "player_idle_3f_chroma.png": "player_idle_3f",
            "player_run_4f_chroma.png": "player_run_4f",
            "player_melee_4f_chroma.png": "player_melee_4f",
            "player_cast_3f_chroma.png": "player_cast_3f",
            "player_hit_death_getup_3f_chroma.png": "player_hit_death_getup_3f",
        }
        keys = ("face_w", "face_h", "chin_y", "face_top", "face_cx", "eye_h",
                "eye_to_chin", "hair_top", "head_h_hair_chin", "band_thick",
                "content_top", "content_bottom")
        for (source, i), r in self.moc.items():
            goc = self.ho_so[(tag[source], i)]
            for k in keys:
                with self.subTest(khung=f"{source}#{i}", thuoc=k):
                    if goc.get(k) is None:
                        self.assertIsNone(r.get(k))
                    else:
                        self.assertAlmostEqual(r[k], goc[k], places=3)

    def test_trung_binh_idle_va_he_so_run0_khop_geometry_v1(self):
        from tools.artqa.nhan_vat_moc import he_so_khung, tham_chieu_idle

        idle = tham_chieu_idle(self.moc)
        self.assertAlmostEqual(idle["face_h"], 96.667, places=2)
        # geometry-V1.json: run0 0.90685, melee2 0.80086, cast1 1.00028, idle0 1.01055
        self.assertAlmostEqual(he_so_khung(self.moc, ("player_run_4f_chroma.png", 0), idle), 0.90685, places=4)
        self.assertAlmostEqual(he_so_khung(self.moc, ("player_melee_4f_chroma.png", 2), idle), 0.80086, places=4)
        self.assertAlmostEqual(he_so_khung(self.moc, ("player_cast_3f_chroma.png", 1), idle), 1.00028, places=4)
        self.assertAlmostEqual(he_so_khung(self.moc, ("player_idle_3f_chroma.png", 0), idle), 1.01055, places=4)
        # khung nằm không có mặt → None
        self.assertIsNone(he_so_khung(self.moc, ("player_hit_death_getup_3f_chroma.png", 1), idle))


if __name__ == "__main__":
    unittest.main()
