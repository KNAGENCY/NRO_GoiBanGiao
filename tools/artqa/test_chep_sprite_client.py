# -*- coding: utf-8 -*-
"""Test tools/nhan-vat/chep-sprite-client.py — chép 37 sprite sang client, sao lưu MỘT LẦN.

Vì sao (06/09/2026): bản cũ chép bản client vào thư mục sao lưu MỌI lần chạy `--that`,
nên lần thứ hai (sau khi client đã nhận bản mới) đè luôn bản lùi trước-05/09 — mất đường
lùi mà không ai báo. Test dựng nguồn/client/sao lưu giả trong thư mục tạm, chạy hai lần và
kiểm bản sao lưu vẫn là bản client ĐẦU TIÊN. Không đụng D:\\NRO_Project.
"""
import contextlib
import importlib.util
import io
import tempfile
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
SCRIPT = ROOT / "tools" / "nhan-vat" / "chep-sprite-client.py"


def _nap():
    spec = importlib.util.spec_from_file_location("chep_sprite_client", SCRIPT)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def _dung(goc: Path, m, noi_dung_client: bytes, noi_dung_nguon: bytes):
    nguon, client, sao_luu = goc / "nguon", goc / "client", goc / "sao-luu"
    for zoom in m.ZOOMS:
        for sid in m.DAI_ID:
            p = nguon / zoom / "smallimage" / f"Small{sid}.png"
            p.parent.mkdir(parents=True, exist_ok=True)
            p.write_bytes(noi_dung_nguon)
            c = client / zoom / "smallimage" / f"Small{sid}.png"
            c.parent.mkdir(parents=True, exist_ok=True)
            c.write_bytes(noi_dung_client)
    return nguon, client, sao_luu


def _chay(m, nguon, client, sao_luu, that):
    argv = ["--nguon", str(nguon), "--client", str(client), "--sao-luu", str(sao_luu)]
    if that:
        argv.append("--that")
    buf = io.StringIO()
    with contextlib.redirect_stdout(buf):
        ma = m.main(argv)
    return ma, buf.getvalue()


class ChepSpriteClientTests(unittest.TestCase):
    def setUp(self):
        self.m = _nap()

    def test_dai_id_gom_34_nguoi_choi_va_3_iruka(self):
        self.assertEqual(list(self.m.DAI_ID), list(range(27041, 27078)))
        self.assertEqual(len(self.m.DAI_ID), 37)

    def test_xem_truoc_khong_ghi_gi(self):
        with tempfile.TemporaryDirectory() as tam:
            nguon, client, sao_luu = _dung(Path(tam), self.m, b"CU", b"MOI")
            ma, out = _chay(self.m, nguon, client, sao_luu, that=False)
            self.assertEqual(ma, 0, out)
            self.assertIn("se chep 148 file", out)  # 37 × 4 zoom
            self.assertIn("sao luu moi 148 / giu 0", out)
            self.assertIn("XEM TRUOC", out)
            self.assertFalse(sao_luu.exists())
            self.assertEqual((client / "x1" / "smallimage" / "Small27075.png").read_bytes(), b"CU")

    def test_that_hai_lan_khong_ghi_de_ban_lui(self):
        with tempfile.TemporaryDirectory() as tam:
            nguon, client, sao_luu = _dung(Path(tam), self.m, b"CU", b"MOI")
            ma, out = _chay(self.m, nguon, client, sao_luu, that=True)
            self.assertEqual(ma, 0, out)
            self.assertIn("sao luu moi 148 / giu 0", out)
            self.assertIn("DA GHI", out)
            self.assertEqual((sao_luu / "x4" / "Small27041.png").read_bytes(), b"CU")
            self.assertEqual((client / "x4" / "smallimage" / "Small27041.png").read_bytes(), b"MOI")
            # lần hai: client đã là MOI; sao lưu phải vẫn là CU
            ma, out = _chay(self.m, nguon, client, sao_luu, that=True)
            self.assertEqual(ma, 0, out)
            self.assertIn("sao luu moi 0 / giu 148", out)
            self.assertEqual((sao_luu / "x4" / "Small27041.png").read_bytes(), b"CU")
            self.assertEqual((sao_luu / "x1" / "Small27077.png").read_bytes(), b"CU")

    def test_sao_luu_mot_phan_da_co_thi_chi_luu_phan_moi(self):
        """Tình huống thật 06/09: 136 file người chơi đã sao lưu, 12 file Iruka chưa."""
        with tempfile.TemporaryDirectory() as tam:
            nguon, client, sao_luu = _dung(Path(tam), self.m, b"CU", b"MOI")
            for zoom in self.m.ZOOMS:
                for sid in range(27041, 27075):
                    p = sao_luu / zoom / f"Small{sid}.png"
                    p.parent.mkdir(parents=True, exist_ok=True)
                    p.write_bytes(b"LUI-CU")
            ma, out = _chay(self.m, nguon, client, sao_luu, that=True)
            self.assertEqual(ma, 0, out)
            self.assertIn("sao luu moi 12 / giu 136", out)
            self.assertEqual((sao_luu / "x2" / "Small27050.png").read_bytes(), b"LUI-CU")
            self.assertEqual((sao_luu / "x2" / "Small27076.png").read_bytes(), b"CU")

    def test_thieu_nguon_thoat_1_khong_ghi(self):
        with tempfile.TemporaryDirectory() as tam:
            nguon, client, sao_luu = _dung(Path(tam), self.m, b"CU", b"MOI")
            (nguon / "x1" / "smallimage" / "Small27077.png").unlink()
            ma, out = _chay(self.m, nguon, client, sao_luu, that=True)
            self.assertEqual(ma, 1, out)
            self.assertIn("THIEU 1 file nguon", out)


if __name__ == "__main__":
    unittest.main()
