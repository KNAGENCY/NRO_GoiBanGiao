"""Chốt kiểm hình học 3 mảnh nhân vật người chơi (part 30 đầu / 57 thân / 58 chân).

TẠI SAO có file này — lỗi đã nổ hai lần và không lần nào test bắt được:

1. Pipeline cũ (`fit_to_canvas`) co TỪNG khung cho vừa hộp 128×160/128×128, nên mỗi
   khung một tỉ lệ: khung nằm rộng 665 px co còn 1/5 khung đứng rộng 340 px. Nhân
   vật phình/teo theo pose, bàn chân trôi -10..-26 px so với `cy` (nhảy lên xuống
   khi chạy), đầu hở +7 px ở cf 20/31 và chồng -22 px ở cf 23/24.
2. `_parts()` trong closure trả `dy` cố định 0 rồi có người sửa tay trong
   `mapping-id.json` — lần sinh closure kế tiếp sẽ âm thầm lùi lại số 0.

Ba nguồn dưới đây phải cùng nói một chuyện, lệch một cái là đỏ ngay:
  1. `docs/my-thuat/mapping-id.json` → `build_inputs.player_component_geometry`
  2. `server/internal/game/assets/part.bin` (thứ máy chủ gửi xuống client)
  3. `docs/my-thuat/lang-la-client-assets/res/x*/smallimage/Small<id>.png` (ảnh thật)

Các phép đo dùng ĐÚNG phép toán vẽ của client (`Char.paintCharBody`), qua
`tools/nhan-vat/ghep-khung.py`, nên không cần mở Unity.
"""

import importlib.util
import json
import os
import unittest
from pathlib import Path

from PIL import Image


ROOT = Path(__file__).resolve().parents[2]
MAPPING = ROOT / "docs/my-thuat/mapping-id.json"
PART_BIN = ROOT / "server/internal/game/assets/part.bin"
SPRITE_ROOT = ROOT / "docs/my-thuat/lang-la-client-assets/res"
CHARINFO = ROOT / "tools/nhan-vat/charinfo.json"
KHUNG_NAM = {20, 31}  # thân lấy từ hai khung "nằm" (hit_death_getup slot 15)


def _ghep_khung():
    spec = importlib.util.spec_from_file_location(
        "ghep_khung", ROOT / "tools/nhan-vat/ghep-khung.py"
    )
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def _do_bang():
    """Đo lại 33 cf: (bàn chân so với cy, khe đầu-thân, khe thân-chân, cao tổng)."""
    module = _ghep_khung()
    charinfo = json.loads(CHARINFO.read_text(encoding="utf-8"))
    parts = module.doc_part_bin(str(PART_BIN))
    sprites = str(SPRITE_ROOT / "x1" / "smallimage")
    rows = {}
    for cf in range(len(charinfo)):
        _, vung = module.ghep(cf, charinfo, parts, sprites, (30, 57, 58), 1, 48, 96)
        dau, than, chan = vung["dau"][1], vung["than"][1], vung["chan"][1]
        rows[cf] = {
            "chan_so_cy": chan[3] - 96,
            "khe_dau_than": than[1] - dau[3],
            "khe_than_chan": chan[1] - than[3],
            "cao_tong": chan[3] - dau[1],
        }
    return rows


class ToanHinhHocTests(unittest.TestCase):
    """Toán thuần — không cần ảnh nguồn, chạy được ở mọi checkout."""

    def test_ti_le_chung_chot_boi_khung_idle_0(self):
        from tools.artqa.task6b_closure import player_uniform_scale

        self.assertAlmostEqual(player_uniform_scale(340), 128 / 340)
        self.assertAlmostEqual(player_uniform_scale(665), 128 / 665)

    def test_ti_le_tu_choi_be_rong_khong_hop_le(self):
        from tools.artqa.task6b_closure import ClosureError, player_uniform_scale

        for width in (0, -1):
            with self.subTest(width=width), self.assertRaises(ClosureError):
                player_uniform_scale(width)

    def test_dx_dy_bu_dung_charinfo_va_lam_tron_ve_x1(self):
        from tools.artqa.task6b_closure import component_placement

        # 40 px sang phải và 200 px lên trên so với gốc x4 → 10 / -50 ở x1.
        placement = component_placement(
            component_left_x4=140.0,
            component_top_x4=100.0,
            frame_center_x4=100.0,
            frame_feet_x4=300.0,
            charinfo_row=[1, -9, 16],
            feet_offset_x1=0,
        )
        self.assertEqual(placement["left_rel"], 10)
        self.assertEqual(placement["top_rel"], -50)
        # dx = left_rel - CharInfo.dx ; dy = top_rel + CharInfo.dy
        self.assertEqual(placement["dx"], 10 - (-9))
        self.assertEqual(placement["dy"], -50 + 16)

    def test_feet_offset_doi_thang_top_rel_va_dy(self):
        from tools.artqa.task6b_closure import component_placement

        chung = dict(
            component_left_x4=100.0,
            component_top_x4=100.0,
            frame_center_x4=100.0,
            frame_feet_x4=300.0,
            charinfo_row=[1, 0, 0],
        )
        khong = component_placement(**chung, feet_offset_x1=0)
        lech = component_placement(**chung, feet_offset_x1=-11)
        self.assertEqual(lech["top_rel"] - khong["top_rel"], -11)
        self.assertEqual(lech["dy"] - khong["dy"], -11)

    def test_dx_dy_ngoai_signed_byte_bi_tu_choi(self):
        from tools.artqa.task6b_closure import ClosureError, component_placement

        with self.assertRaisesRegex(ClosureError, "signed byte"):
            component_placement(
                component_left_x4=0.0,
                component_top_x4=0.0,
                frame_center_x4=4000.0,
                frame_feet_x4=0.0,
                charinfo_row=[0, 0, 0],
            )

    def test_slot_to_cf_lay_tu_charinfo_chu_khong_phai_bang_required(self):
        """`REQUIRED_CLIENT_CF` LỆCH một bậc ở thân/chân — bù theo nó là sai hết.

        `CharInfo[0][2][0] == 1` nghĩa là cf 0 lấy slot thân **1**, trong khi
        `REQUIRED_CLIENT_CF["hoc-vien-body"][0] == (0, 1)` nói slot 0 phục vụ cf 0/1.
        Ghi lại sự thật ở đây để không ai "sửa cho khớp bảng" rồi lệch một khung.
        """
        from assetpack.animation import REQUIRED_CLIENT_CF
        from tools.artqa.task6b_closure import charinfo_slot_to_cf, load_charinfo

        charinfo = load_charinfo(CHARINFO)
        head = charinfo_slot_to_cf(charinfo, "head")
        body = charinfo_slot_to_cf(charinfo, "body")
        leg = charinfo_slot_to_cf(charinfo, "leg")

        # Đầu thì hai bảng khớp nhau.
        for slot, cf_list in enumerate(REQUIRED_CLIENT_CF["hoc-vien-head"]):
            self.assertEqual(head[slot], list(cf_list), f"head slot {slot}")
        # Thân/chân thì không.
        self.assertEqual(body[1], [0, 1])
        self.assertEqual(body[0], [23, 24])
        self.assertEqual(leg[1], [0, 1, 22])
        self.assertEqual(leg[0], [8])
        self.assertNotEqual(body[0], list(REQUIRED_CLIENT_CF["hoc-vien-body"][0]))
        self.assertNotEqual(leg[0], list(REQUIRED_CLIENT_CF["hoc-vien-leg"][0]))
        # Mọi slot đều có ít nhất một cf dùng — không slot nào chết.
        self.assertEqual(sorted(head), list(range(3)))
        self.assertEqual(sorted(body), list(range(17)))
        self.assertEqual(sorted(leg), list(range(14)))


class CongCuDoTests(unittest.TestCase):
    def test_ghep_khung_doc_duoc_ca_hai_bo_cuc_thu_muc_sprite(self):
        """Thư mục client phẳng và thư mục closure phải cùng đo được."""
        module = _ghep_khung()
        self.assertTrue(
            module.duong_sprite("kho", 27041, False, 1).endswith("Small27041.png")
        )
        closure = module.duong_sprite("kho", 27041, True, 2)
        self.assertIn("27041", closure)
        self.assertTrue(closure.endswith(os.path.join("x2", "sprite.png")))


class ArtifactKhopNhauTests(unittest.TestCase):
    """mapping ↔ part.bin ↔ ảnh thật trong repo (không cần kho ảnh nguồn)."""

    @classmethod
    def setUpClass(cls):
        cls.mapping = json.loads(MAPPING.read_text(encoding="utf-8"))
        cls.geometry = cls.mapping["build_inputs"]["player_component_geometry"]

    def test_moi_manh_dung_dung_mot_ti_le(self):
        scale = self.geometry["uniform_scale_x4"]
        self.assertAlmostEqual(scale, 128 / 340)
        for item in self.geometry["components"]:
            width, height = item["size_x4"]
            rect_width, rect_height = item["source_rect"][2], item["source_rect"][3]
            # Ảnh đã trim alpha nên phải ≤ vùng cắt đã scale bằng tỉ lệ CHUNG.
            # Vượt = đã co bằng một tỉ lệ khác (phình lên cho vừa hộp).
            self.assertLessEqual(width, round(rect_width * scale) + 1, item["id"])
            self.assertLessEqual(height, round(rect_height * scale) + 1, item["id"])

    def test_khong_con_nhoi_vao_hop_co_dinh(self):
        """Dấu vân tay của `fit_to_canvas`: mọi ảnh đúng 128×128 hoặc 128×160.

        Pipeline cũ co từng khung cho vừa hộp nên 34 sprite chỉ có 2 kích thước.
        Với một tỉ lệ chung, kích thước phải bám theo pose thật.
        """
        sizes = [tuple(item["size_x4"]) for item in self.geometry["components"]]
        for size in sizes:
            self.assertNotIn(size, {(128, 128), (128, 160)})
        self.assertGreaterEqual(len(set(sizes)), 20, sorted(set(sizes)))

    def test_khung_idle_0_giu_dung_co_cu(self):
        """Nhân vật trong game không được đổi cỡ: khung idle 0 vẫn rộng 128 px ở x4."""
        scale = self.geometry["uniform_scale_x4"]
        idle0 = [
            item
            for item in self.geometry["components"]
            if item["action"] == "idle" and item["action_slot"] == 0
        ]
        self.assertEqual(len(idle0), 3)
        for item in idle0:
            self.assertEqual(item["source_rect"][2], 340, item["id"])
        self.assertEqual(round(340 * scale), 128)

    def test_anh_that_dung_kich_thuoc_khai_va_du_bon_muc_zoom(self):
        for item in self.geometry["components"]:
            for zoom, relative in item["paths"].items():
                path = ROOT / relative
                self.assertTrue(path.is_file(), relative)
                with Image.open(path) as image:
                    size = image.size
                if zoom == "x4":
                    self.assertEqual(list(size), item["size_x4"], relative)
                else:
                    index = int(zoom[1])
                    self.assertEqual(
                        list(size),
                        [
                            max(1, round(item["size_x4"][0] * index / 4)),
                            max(1, round(item["size_x4"][1] * index / 4)),
                        ],
                        relative,
                    )

    def test_part_bin_mang_dung_dx_dy_cua_mapping(self):
        from assetpack.part import parse_part

        parts = parse_part(PART_BIN.read_bytes())
        khai = {item["id"]: item["part_offset"] for item in self.geometry["components"]}
        for part in self.mapping["parts"]:
            if part["id"] not in (30, 57, 58):
                continue
            thuc = parts[part["id"]]
            self.assertEqual(len(thuc.frames), len(part["frames"]), part["id"])
            for index, frame in enumerate(part["frames"]):
                self.assertEqual(thuc.frames[index].sprite_id, frame["sprite_id"])
                self.assertEqual(thuc.frames[index].dx, frame["dx"], (part["id"], index))
                self.assertEqual(thuc.frames[index].dy, frame["dy"], (part["id"], index))
                self.assertEqual(khai[frame["sprite_id"]]["dx"], frame["dx"])
                self.assertEqual(khai[frame["sprite_id"]]["dy"], frame["dy"])


class DoTrenKhungVeTests(unittest.TestCase):
    """Đo lại đúng như client vẽ — đây là thứ người chơi thật sự nhìn thấy."""

    @classmethod
    def setUpClass(cls):
        cls.rows = _do_bang()

    def test_ban_chan_moi_khung_cung_mot_cao_do(self):
        """Trước: -10..-26 px (nhân vật nhảy lên xuống khi chạy). Sau: lệch ≤ 1 px."""
        goc = self.rows[0]["chan_so_cy"]
        for cf, row in self.rows.items():
            self.assertLessEqual(
                abs(row["chan_so_cy"] - goc), 1,
                f"cf {cf}: ban chan {row['chan_so_cy']} vs cf0 {goc}",
            )

    def test_chieu_cao_tong_khong_doi_qua_3px(self):
        """Nhân vật không được phình/teo theo pose."""
        goc = self.rows[0]["cao_tong"]
        for cf, row in self.rows.items():
            self.assertLessEqual(
                abs(row["cao_tong"] - goc), 3,
                f"cf {cf}: cao tong {row['cao_tong']} vs cf0 {goc}",
            )

    def test_than_va_chan_luon_chong_nhau_khong_bao_gio_ho(self):
        """Trước có khe hở tới +8 px giữa thân và chân; giờ phải luôn ≤ 0."""
        for cf, row in self.rows.items():
            self.assertLessEqual(row["khe_than_chan"], 0, f"cf {cf} ho than-chan")
            self.assertGreaterEqual(row["khe_than_chan"], -12, f"cf {cf} chong qua sau")

    def test_khe_dau_than_trong_gioi_han_da_do(self):
        """Ràng buộc còn lại của NGUỒN, không phải của thuật toán.

        Năm strip nguồn vẽ nhân vật cao thấp khác nhau (khung idle cao 653 px,
        khung chạy f1 chỉ 459 px = 70%). Với MỘT tỉ lệ chung và neo bàn chân,
        thân ở khung nhỏ không thể vừa nối cổ vừa giữ chân đúng chỗ: 8 cf hở
        +1..+5 px. Hai cf 20/31 lấy thân từ khung NẰM nên hở +19 px — đã biết
        và chấp nhận (ngoại lệ ghi rõ trong kế hoạch).
        Trước khi sửa dải này là -22..+7 px; test khoá để không tệ hơn.
        """
        for cf, row in self.rows.items():
            gioi_han = 19 if cf in KHUNG_NAM else 5
            self.assertLessEqual(row["khe_dau_than"], gioi_han, f"cf {cf} ho dau-than")
            self.assertGreaterEqual(row["khe_dau_than"], -12, f"cf {cf} chong qua sau")


if __name__ == "__main__":
    unittest.main()
