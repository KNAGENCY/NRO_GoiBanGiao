# -*- coding: utf-8 -*-
"""Test cho ``tools/artqa/nhan_vat_toi_uu.py`` — MILP chọn dx/dy cho 34 mảnh (V1 05/09/2026).

Vì sao có test này: bộ xương cứng của client (``CharInfo``: 3 đầu dùng chung cho 33 cf,
một dx/dy cho mỗi mảnh) khiến "chỉnh tay" dx/dy không hội tụ — chỉnh khung này hở khung
kia. Bài toán là MILP thật (ẩn nguyên, ràng buộc khe cứng, mục tiêu sai số cằm), giải
bằng ``scipy.optimize.milp``. Ba điều phải khoá: (1) nghiệm trả về là số nguyên trong
signed byte và thoả mọi ràng buộc cứng, (2) ``kiem_nghiem`` bắt được nghiệm sai mà KHÔNG
cần giải lại (để mapping tự chứng minh), (3) vô nghiệm phải báo rõ, không im lặng nới.
"""

import unittest


def _charinfo():
    # [head], [leg], [body] — đúng thứ tự k = 0/1/2 của Char.cs.
    return [
        [[0, -13, 34], [0, -8, 10], [0, -9, 16]],
        [[0, -13, 35], [0, -8, 10], [1, -9, 17]],
    ]


def _manh(component, slot, action, action_slot, *, bbox, chin_off, feet_off, face_cx_off, center_off, hip_cx_off):
    return {
        "component": component,
        "slot_index": slot,
        "action": action,
        "action_slot": action_slot,
        "moc": {
            "bbox_x1": list(bbox),
            "chin_off": chin_off,
            "feet_off": feet_off,
            "face_cx_off": face_cx_off,
            "center_off": center_off,
            "hip_cx_off": hip_cx_off,
        },
    }


def _components():
    return [
        # đầu idle0: cao 26, cằm ở hàng 25 (đơn vị x1, so với góc trên-trái ảnh)
        _manh("head", 0, "idle", 0, bbox=(0, 0, 32, 26), chin_off=25.0, feet_off=60.0, face_cx_off=19.0, center_off=16.0, hip_cx_off=16.0),
        # thân idle0/idle1: mép trên = cằm − 2 px → chin_off = 2; bàn chân "ảo" cách mép trên 40/41
        _manh("body", 0, "idle", 0, bbox=(0, 0, 29, 30), chin_off=2.0, feet_off=40.0, face_cx_off=15.0, center_off=14.5, hip_cx_off=14.0),
        _manh("body", 1, "idle", 1, bbox=(0, 0, 29, 30), chin_off=2.0, feet_off=41.0, face_cx_off=15.0, center_off=14.5, hip_cx_off=14.0),
        # chân idle0: đáy ảnh = bàn chân
        _manh("leg", 0, "idle", 0, bbox=(0, 0, 27, 21), chin_off=None, feet_off=21.0, face_cx_off=None, center_off=13.5, hip_cx_off=13.0),
    ]


class GiaiMILPTests(unittest.TestCase):
    def test_nghiem_nguyen_thoa_moi_rang_buoc_va_ban_chan_neo(self):
        from tools.artqa.nhan_vat_toi_uu import ThamSoMILP, bang_cf, giai, kiem_nghiem

        ts = ThamSoMILP(ghim_dy_dau={})
        nghiem = giai(_components(), _charinfo(), ts)
        for phan in ("head", "body", "leg"):
            for truc in ("dx", "dy"):
                for v in nghiem[phan][truc]:
                    self.assertIsInstance(v, int)
                    self.assertTrue(-128 <= v <= 127)
        self.assertEqual(len(nghiem["head"]["dy"]), 1)
        self.assertEqual(len(nghiem["body"]["dy"]), 2)
        self.assertEqual(len(nghiem["leg"]["dy"]), 1)
        self.assertEqual(kiem_nghiem(_components(), _charinfo(), nghiem, ts), [])
        bang = bang_cf(_components(), _charinfo(), nghiem, ts)
        self.assertEqual(len(bang), 2)
        # dy chân cố định: đáy chân = cy − 11 tại cf neo (cf 0) → dy_l = −11 − 21 + 10 = −22
        self.assertEqual(nghiem["leg"]["dy"][0], -22)
        self.assertEqual(bang[0]["chan"], -11)
        for hang in bang:
            self.assertGreaterEqual(hang["khe_dau_than"], ts.khe_dau_than[0])
            self.assertLessEqual(hang["khe_dau_than"], ts.khe_dau_than[1])
            self.assertGreaterEqual(hang["khe_than_chan"], ts.khe_than_chan[0])
            self.assertLessEqual(hang["khe_than_chan"], ts.khe_than_chan[1])
            self.assertGreaterEqual(hang["e"], ts.e[0] - 1e-6)
            self.assertLessEqual(hang["e"], ts.e[1] + 1e-6)
            self.assertLessEqual(abs(hang["d"]), ts.d_max + 1e-6)

    def test_ghim_dy_dau_duoc_ton_trong(self):
        from tools.artqa.nhan_vat_toi_uu import ThamSoMILP, giai, kiem_nghiem

        # không ghim thì tối ưu là −40 (e = 0); ghim −43 khả thi nhưng không tối ưu →
        # nghiệm phải theo ghim, không theo mục tiêu.
        ts = ThamSoMILP(ghim_dy_dau={0: -43})
        nghiem = giai(_components(), _charinfo(), ts)
        self.assertEqual(nghiem["head"]["dy"][0], -43)
        self.assertEqual(kiem_nghiem(_components(), _charinfo(), nghiem, ts), [])

    def test_vo_nghiem_bao_ro_khong_tu_noi(self):
        from tools.artqa.nhan_vat_toi_uu import ThamSoMILP, VoNghiem, giai

        # ghim đầu quá cao (dy −120) mà khe đầu–thân phải ∈ [−9, 0] và thân phải giữ
        # bàn chân ảo trong ±5,5 → không có dy thân nào thoả cả hai.
        ts = ThamSoMILP(ghim_dy_dau={0: -120})
        with self.assertRaises(VoNghiem):
            giai(_components(), _charinfo(), ts)


class KiemNghiemTests(unittest.TestCase):
    def test_bat_nghiem_lech_ma_khong_giai_lai(self):
        from tools.artqa.nhan_vat_toi_uu import ThamSoMILP, giai, kiem_nghiem

        ts = ThamSoMILP(ghim_dy_dau={})
        nghiem = giai(_components(), _charinfo(), ts)
        hong = {k: ({t: list(v) for t, v in val.items()} if isinstance(val, dict) else val) for k, val in nghiem.items()}
        hong["body"]["dy"][0] += 12  # đẩy thân xuống 12 px → hở đầu–thân và trượt chân
        loi = kiem_nghiem(_components(), _charinfo(), hong, ts)
        self.assertTrue(loi)
        self.assertTrue(any("khe_dau_than" in m or "d" in m for m in loi), loi)

    def test_bat_gia_tri_khong_nguyen_va_ngoai_byte(self):
        from tools.artqa.nhan_vat_toi_uu import ThamSoMILP, giai, kiem_nghiem

        ts = ThamSoMILP(ghim_dy_dau={})
        nghiem = giai(_components(), _charinfo(), ts)
        hong = {k: ({t: list(v) for t, v in val.items()} if isinstance(val, dict) else val) for k, val in nghiem.items()}
        hong["head"]["dx"][0] = 3.5
        self.assertTrue(any("nguyen" in m for m in kiem_nghiem(_components(), _charinfo(), hong, ts)))
        hong["head"]["dx"][0] = 200
        self.assertTrue(any("byte" in m for m in kiem_nghiem(_components(), _charinfo(), hong, ts)))

    def test_bat_ghim_bi_pha(self):
        from tools.artqa.nhan_vat_toi_uu import ThamSoMILP, giai, kiem_nghiem

        ts = ThamSoMILP(ghim_dy_dau={0: -43})
        nghiem = giai(_components(), _charinfo(), ts)
        hong = {k: ({t: list(v) for t, v in val.items()} if isinstance(val, dict) else val) for k, val in nghiem.items()}
        hong["head"]["dy"][0] = -42
        self.assertTrue(any("ghim" in m for m in kiem_nghiem(_components(), _charinfo(), hong, ts)))

    def test_tham_so_di_qua_json_khong_mat(self):
        from tools.artqa.nhan_vat_toi_uu import ThamSoMILP

        ts = ThamSoMILP(ghim_dy_dau={0: -38}, d_max=5.5)
        lai = ThamSoMILP.tu_json(ts.json())
        self.assertEqual(lai, ts)
        self.assertEqual(lai.ghim_dy_dau, {0: -38})


if __name__ == "__main__":
    unittest.main()
