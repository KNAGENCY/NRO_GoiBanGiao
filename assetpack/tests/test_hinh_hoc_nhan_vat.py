"""Chốt kiểm hình học 3 mảnh nhân vật người chơi (part 30 đầu / 57 thân / 58 chân).

TẠI SAO có file này — lỗi đã nổ BA lần và không lần nào test bắt được:

1. Pipeline cũ (`fit_to_canvas`) co TỪNG khung cho vừa hộp 128×160/128×128, nên mỗi
   khung một tỉ lệ: khung nằm rộng 665 px co còn 1/5 khung đứng rộng 340 px. Nhân
   vật phình/teo theo pose, bàn chân trôi -10..-26 px so với `cy` (nhảy lên xuống
   khi chạy), đầu hở +7 px ở cf 20/31 và chồng -22 px ở cf 23/24.
2. `_parts()` trong closure trả `dy` cố định 0 rồi có người sửa tay trong
   `mapping-id.json` — lần sinh closure kế tiếp sẽ âm thầm lùi lại số 0.
3. (05/09/2026, V1) "Đầu tách thân" thực chất là HAI KHUÔN MẶT: mảnh thân cắt ở 28 %
   chiều cao khung chứa cả mặt và vẽ đè lên đầu; 5 dải nguồn khác tỉ lệ (run nhỏ hơn
   idle ~12 %, melee ~17 %); khung NẰM gán nhầm slot thân 15; `REQUIRED_CLIENT_CF`
   lệch slot so với CharInfo. Sửa bằng dữ liệu: tỉ lệ TỪNG khung theo cỡ đầu, cắt thân
   từ cằm − 2 px x1, gán lại slot, dx/dy bằng MILP (`tools/artqa/nhan_vat_toi_uu.py`).

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
ASSET_LOCK = ROOT / "server/internal/game/assets/asset-lock.json"
SPRITE_ROOT = ROOT / "docs/my-thuat/lang-la-client-assets/res"
CHARINFO = ROOT / "tools/nhan-vat/charinfo.json"

# Bảng slot V1 (quyết định 05/09/2026): cf 2..6 là 5 khung CHẠY (`Char.cs:2537`,
# `cf = gameTick % 5 + 2`) nên slot 2..6 = run0..3,run0; slot thân 15 (cf 20/31 = thế
# bắt đầu skill) từng là khung NẰM → cast0; slot 14 = cast2; slot 16 + chân 13 (cf 25
# "sau dịch chuyển") = idle2. Ba khung hit_death_getup KHÔNG dùng (client không vẽ nằm;
# hit0 lệch tâm mặt 15 px).
BANG_THAN_V1 = {
    0: ("idle", 0), 1: ("idle", 1), 2: ("run", 0), 3: ("run", 1), 4: ("run", 2), 5: ("run", 3),
    6: ("run", 0), 7: ("melee", 0), 8: ("melee", 1), 9: ("melee", 2), 10: ("melee", 3),
    11: ("cast", 0), 12: ("cast", 1), 13: ("cast", 2), 14: ("cast", 2), 15: ("cast", 0), 16: ("idle", 2),
}
BANG_CHAN_V1 = {
    0: ("idle", 0), 1: ("idle", 1), 2: ("run", 0), 3: ("run", 1), 4: ("run", 2), 5: ("run", 3),
    6: ("run", 0), 7: ("melee", 0), 8: ("melee", 1), 9: ("melee", 2), 10: ("melee", 3),
    11: ("cast", 0), 12: ("cast", 1), 13: ("idle", 2),
}
BANG_DAU_V1 = {0: ("idle", 0), 1: ("idle", 1), 2: ("idle", 2)}


def _ghep_khung():
    spec = importlib.util.spec_from_file_location(
        "ghep_khung", ROOT / "tools/nhan-vat/ghep-khung.py"
    )
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def _do_bang(zoom: int = 1):
    """Đo lại 33 cf ở mức zoom cho trước: (bàn chân so với cy, khe đầu-thân, khe thân-chân, cao tổng).

    Đơn vị = px của mức zoom đó (x1: 1 px = 1 đơn vị part.bin; x4: 4 px = 1 đơn vị).
    """
    module = _ghep_khung()
    charinfo = json.loads(CHARINFO.read_text(encoding="utf-8"))
    parts = module.doc_part_bin(str(PART_BIN))
    sprites = str(SPRITE_ROOT / f"x{zoom}" / "smallimage")
    cx, cy = 48 * zoom, 96 * zoom
    rows = {}
    for cf in range(len(charinfo)):
        _, vung = module.ghep(cf, charinfo, parts, sprites, (30, 57, 58), zoom, cx, cy)
        dau, than, chan = vung["dau"][1], vung["than"][1], vung["chan"][1]
        rows[cf] = {
            "chan_so_cy": chan[3] - cy,
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

    def test_required_client_cf_va_charinfo_slot_to_cf_cung_mot_su_that(self):
        """Hai bảng slot→cf phải KHỚP nhau — và khớp CharInfo.

        Lịch sử: tới 05/09/2026 `REQUIRED_CLIENT_CF` viết tay sắp theo cf đầu tiên nên
        lệch slot ở thân/chân (`CharInfo[0][2][0] == 1`: cf 0 lấy slot thân 1, bảng cũ
        nói slot 0) và chân thiếu hẳn cf 26..32; test cũ chỉ ghi nhận "lệch" để không
        ai sửa cho khớp bảng sai. Nay bảng SINH từ charinfo.json (assetpack/animation.py)
        nên hai nguồn phải bằng nhau — lệch = ai đó viết tay lại.
        """
        from assetpack.animation import REQUIRED_CLIENT_CF
        from tools.artqa.task6b_closure import charinfo_slot_to_cf, load_charinfo

        charinfo = load_charinfo(CHARINFO)
        for component, name in (("head", "hoc-vien-head"), ("body", "hoc-vien-body"), ("leg", "hoc-vien-leg")):
            by_slot = charinfo_slot_to_cf(charinfo, component)
            self.assertEqual(sorted(by_slot), list(range(len(REQUIRED_CLIENT_CF[name]))), name)
            for slot, cf_list in enumerate(REQUIRED_CLIENT_CF[name]):
                self.assertEqual(by_slot[slot], list(cf_list), f"{name} slot {slot}")
        # Vài giá trị chốt để đọc bằng mắt.
        self.assertEqual(REQUIRED_CLIENT_CF["hoc-vien-body"][1], (0, 1))
        self.assertEqual(REQUIRED_CLIENT_CF["hoc-vien-body"][0], (23, 24))
        self.assertEqual(REQUIRED_CLIENT_CF["hoc-vien-leg"][1], (0, 1, 22))
        self.assertEqual(REQUIRED_CLIENT_CF["hoc-vien-leg"][0], (8,))

    def test_lat_cat_v1_theo_cam_khong_theo_phan_tram(self):
        """Đầu: đỉnh tóc → cằm + 6 px nguồn. Thân: cằm − 2 px x1 (= 8/s_f px nguồn) → 76 %.
        Chân: 66 % → 100 %. Đây là phần sửa "hai khuôn mặt": thân không được chứa mặt."""
        from tools.artqa.task6b_closure import _component_rect_v1

        record = {
            "asset": "x",
            "checks": {"crop": {"source_edges": [0, 700], "content_boxes": [[100, 60, 440, 713]]}},
        }
        khung = {"chin_y": 332.0, "scale": 0.4}  # cằm ở y=332 → 272 px dưới mép trên khung (653 cao)
        dau = _component_rect_v1(record, 0, "head", khung)
        than = _component_rect_v1(record, 0, "body", khung)
        chan = _component_rect_v1(record, 0, "leg", khung)
        self.assertEqual(dau, [100, 60, 340, 272 + 6])
        # 8 / 0.4 = 20 px nguồn trên cằm; đáy = round(653 × 0.76) = 496
        self.assertEqual(than[1], 60 + 272 - 20)
        self.assertEqual(than[1] + than[3], 60 + 496)
        self.assertEqual(chan, [100, 60 + 431, 340, 653 - 431])
        # khung nằm (không có cằm) → thân rơi về 28 % như cũ, đầu về 42 %
        nam = _component_rect_v1(record, 0, "body", {"chin_y": None, "scale": 0.4})
        self.assertEqual(nam[1], 60 + round(653 * 0.28))
        self.assertEqual(_component_rect_v1(record, 0, "head", {"chin_y": None, "scale": 0.4})[3], round(653 * 0.42))


class HangNeoBanChanTests(unittest.TestCase):
    """Ghim hằng neo bàn chân vào DỮ LIỆU đã sinh.

    Điều phối áp đột biến `PLAYER_FEET_OFFSET_X1 = -5` (05/09/2026): 25 test trên
    vẫn xanh — vì chúng đo artifact đã sinh, còn hằng số chỉ tác dụng lúc SINH LẠI.
    Đổi hằng này là dời cả nhân vật lên/xuống so với mặt đất mà không test nào kêu.
    Ba chốt dưới đây buộc ba nơi (hằng trong code, mapping, số đo trên ảnh+part.bin)
    phải cùng nói -11; muốn dời nhân vật thì đổi cả ba, có chủ đích.
    """

    def test_hang_trong_code_bang_mapping_bang_so_do(self):
        from tools.artqa.task6b_closure import PLAYER_FEET_OFFSET_X1

        self.assertEqual(PLAYER_FEET_OFFSET_X1, -11, "cf 0 bản đang chạy đứng ở cy-11")
        mapping = json.loads(MAPPING.read_text(encoding="utf-8"))
        geometry = mapping["build_inputs"]["player_component_geometry"]
        self.assertEqual(geometry["feet_offset_x1"], PLAYER_FEET_OFFSET_X1)
        self.assertEqual(geometry["milp"]["tham_so"]["ban_chan_x1"], PLAYER_FEET_OFFSET_X1)
        rows = _do_bang()
        self.assertEqual(
            rows[0]["chan_so_cy"], PLAYER_FEET_OFFSET_X1,
            "bàn chân cf 0 đo trên part.bin + ảnh phải bằng hằng neo — lệch = artifact chưa sinh lại",
        )


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

    def test_moi_manh_dung_ti_le_rieng_theo_co_dau(self):
        """V1: s_f = s0 / (cỡ đầu khung / cỡ đầu idle). Ảnh trim ≤ vùng cắt × s_f.

        Trước 05/09 chốt "một tỉ lệ chung" (s0 cho mọi khung); đo lại thấy 5 dải nguồn
        vẽ nhân vật khác cỡ (run 0,86–0,91, melee 0,80–0,85 lần idle theo 6 thước cứng
        của đầu) nên một tỉ lệ chung làm đầu idle to hơn thân run/melee → nhân theo cỡ
        đầu từng khung. s0 vẫn là mốc chung (`uniform_scale_x4`).
        """
        scale = self.geometry["uniform_scale_x4"]
        self.assertAlmostEqual(scale, 128 / 340)
        for item in self.geometry["components"]:
            s_f = item["scale_frame"]
            he_so = item["nguon"]["he_so_dau"]
            self.assertIsNotNone(he_so, item["id"])
            self.assertAlmostEqual(s_f, scale / he_so, places=9, msg=item["id"])
            self.assertTrue(0.75 <= he_so <= 1.05, (item["id"], he_so))
            width, height = item["size_x4"]
            rect_width, rect_height = item["source_rect"][2], item["source_rect"][3]
            self.assertLessEqual(width, round(rect_width * s_f) + 1, item["id"])
            self.assertLessEqual(height, round(rect_height * s_f) + 1, item["id"])

    def test_khong_con_nhoi_vao_hop_co_dinh(self):
        """Dấu vân tay của `fit_to_canvas`: mọi ảnh đúng 128×128 hoặc 128×160.

        Pipeline cũ co từng khung cho vừa hộp nên 34 sprite chỉ có 2 kích thước.
        Với tỉ lệ bám theo pose thật, kích thước phải đa dạng.
        """
        sizes = [tuple(item["size_x4"]) for item in self.geometry["components"]]
        for size in sizes:
            self.assertNotIn(size, {(128, 128), (128, 160)})
        self.assertGreaterEqual(len(set(sizes)), 20, sorted(set(sizes)))

    def test_khung_idle_0_giu_dung_co_cu(self):
        """Nhân vật trong game không được đổi cỡ: mốc s0 giữ khung idle 0 = 128 px ở x4.

        Đầu idle0 thật có thể 127/128 px vì s_f = s0 / hệ số (hệ số idle0 ≈ 1,01 do tham
        chiếu là trung bình 3 khung idle) — lệch ≤ 1 px, không phải đổi cỡ. Thân/chân
        trim theo nội dung nên hẹp hơn khung (đo: 107 px) — chỉ đầu mới rộng bằng khung.
        """
        scale = self.geometry["uniform_scale_x4"]
        idle0 = [
            item
            for item in self.geometry["components"]
            if item["action"] == "idle" and item["action_slot"] == 0
        ]
        self.assertEqual(len(idle0), 3)
        for item in idle0:
            self.assertEqual(item["source_rect"][2], 340, item["id"])
            self.assertAlmostEqual(item["scale_frame"] * item["nguon"]["he_so_dau"], scale, places=9)
            if item["component"] == "head":
                self.assertLessEqual(abs(item["size_x4"][0] - 128), 1, item["id"])
        self.assertEqual(round(340 * scale), 128)

    def test_bang_slot_v1_khong_con_khung_nam(self):
        """Gán slot theo quyết định V1 (xem BANG_*_V1 đầu file); hit_death_getup không dùng."""
        by_slot = {(item["component"], item["slot_index"]): item for item in self.geometry["components"]}
        self.assertEqual(len(by_slot), 34)
        for bang, component in ((BANG_DAU_V1, "head"), (BANG_THAN_V1, "body"), (BANG_CHAN_V1, "leg")):
            for slot, (action, action_slot) in bang.items():
                item = by_slot[(component, slot)]
                self.assertEqual((item["action"], item["action_slot"]), (action, action_slot), (component, slot))
        self.assertFalse(any(item["action"] == "hit_death_getup" for item in self.geometry["components"]))

    def test_than_khong_chua_mat_dau_cat_toi_cam(self):
        """Gốc của "hai khuôn mặt": mảnh thân phải bắt đầu dưới cằm (trừ 2 px x1 gộp).

        `moc.chin_off` = cằm so với mép trên vùng đục của mảnh (x1). Thân ≤ BODY_CHIN_MARGIN_X1 + 0,5 → mặt
        (kết thúc ở cằm) không nằm trong mảnh. Đầu: đáy vùng cắt = cằm + 6 px nguồn.
        """
        from tools.artqa.task6b_closure import BODY_CHIN_MARGIN_X1

        for item in self.geometry["components"]:
            moc = item["moc"]
            chin_src = item["nguon"]["chin_y"]
            self.assertIsNotNone(chin_src, item["id"])
            x, y, w, h = item["source_rect"]
            if item["component"] == "body":
                # ngưỡng = lề cằm gộp vào thân (BODY_CHIN_MARGIN_X1) + 0,5 px sai số làm tròn nguồn → x1
                self.assertLessEqual(moc["chin_off"], BODY_CHIN_MARGIN_X1 + 0.5, (item["id"], moc["chin_off"]))
                self.assertGreaterEqual(moc["chin_off"], -1.0, (item["id"], moc["chin_off"]))
            elif item["component"] == "head":
                self.assertEqual(y + h, round(chin_src) + 6, item["id"])
                self.assertEqual(y, item["nguon"]["frame_top"], item["id"])
            else:
                self.assertEqual(y, item["nguon"]["frame_top"] + round(item["nguon"]["frame_height"] * 0.66), item["id"])

    def test_nghiem_milp_trong_mapping_tu_chung_minh_khong_can_giai_lai(self):
        """Nghiệm dx/dy + tham số MILP ghi trong mapping phải thoả mọi ràng buộc cứng
        khi thay số vào (không gọi scipy giải lại). Đây là cách mapping tự chứng minh."""
        from tools.artqa.nhan_vat_toi_uu import ThamSoMILP, kiem_nghiem

        milp = self.geometry["milp"]
        self.assertEqual(milp["phien_ban"], "V1 05/09")
        tham_so = ThamSoMILP.tu_json(milp["tham_so"])
        self.assertEqual(tham_so.khe_than_chan, (-11, -2))
        self.assertEqual(tham_so.khe_dau_than, (-9, 0))
        charinfo = json.loads(CHARINFO.read_text(encoding="utf-8"))
        components = [
            {
                "component": item["component"],
                "slot_index": item["slot_index"],
                "action": item["action"],
                "action_slot": item["action_slot"],
                "moc": item["moc"],
            }
            for item in self.geometry["components"]
        ]
        self.assertEqual(kiem_nghiem(components, charinfo, milp["nghiem"], tham_so), [])
        # nghiệm trong mapping == part_offset từng mảnh
        for item in self.geometry["components"]:
            slot = item["slot_index"]
            self.assertEqual(item["part_offset"]["dx"], milp["nghiem"][item["component"]]["dx"][slot], item["id"])
            self.assertEqual(item["part_offset"]["dy"], milp["nghiem"][item["component"]]["dy"][slot], item["id"])

    def test_dy_dau_slot_0_dung_ghim(self):
        """Ghim dy đầu slot 0 = −38 (−39 vô nghiệm với khe thân–chân ≤ −2 — điều phối đã thử).
        Đầu idle hạ 1 px so với bản trước (−39) là giá đã chấp nhận."""
        from assetpack.part import parse_part

        milp = self.geometry["milp"]
        self.assertEqual(milp["tham_so"]["ghim_dy_dau"], {"0": -38})
        parts = parse_part(PART_BIN.read_bytes())
        self.assertEqual(parts[30].frames[0].dy, -38)

    def test_manh_nam_trong_o_contact_sheet(self):
        """Ô `REQUIRED_ANIMATION_SHEETS` là giới hạn trên; mảnh V1 to hơn (s_f tới 0,47)
        nhưng vẫn phải nằm trong ô — vượt là phải nới ô có chủ đích."""
        from assetpack.animation import REQUIRED_ANIMATION_SHEETS

        for item in self.geometry["components"]:
            cell = REQUIRED_ANIMATION_SHEETS[f"hoc-vien-{item['component']}"]["cell_size"]
            self.assertLessEqual(item["size_x4"][0], cell[0], item["id"])
            self.assertLessEqual(item["size_x4"][1], cell[1], item["id"])

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

    def test_mapping_animations_khop_required_client_cf(self):
        """`animations[].frames[].client_cf` trong mapping là thứ `assetpack build` kiểm
        (validate_animation_closure) — phải theo bảng sinh từ CharInfo, cả sprites[]."""
        from assetpack.animation import REQUIRED_CLIENT_CF

        for item in self.mapping["animations"]:
            for frame in item["frames"]:
                self.assertEqual(
                    frame["client_cf"], list(REQUIRED_CLIENT_CF[item["name"]][frame["slot_index"]]),
                    (item["name"], frame["slot_index"]),
                )
        for sprite in self.mapping["sprites"]:
            if 27041 <= sprite["id"] <= 27074:
                self.assertEqual(
                    sprite["client_cf"],
                    list(REQUIRED_CLIENT_CF[f"hoc-vien-{sprite['component']}"][sprite["slot_index"]]),
                    sprite["id"],
                )

    def test_lock_versions_data_28_va_sha_part_bin(self):
        """part.bin đổi → versions.data phải tăng (27 → 28) để client bỏ cache RMS;
        SHA trong lock phải là SHA của part.bin đang nhúng."""
        import hashlib

        lock = json.loads(ASSET_LOCK.read_text(encoding="utf-8"))
        self.assertEqual(lock["versions"]["data"], 28)
        part = PART_BIN.read_bytes()
        output = next(o for o in lock["outputs"] if o["path"] == "server/part.bin")
        self.assertEqual(output["sha256"], hashlib.sha256(part).hexdigest())
        self.assertEqual(output["bytes"], len(part))


class DoTrenKhungVeTests(unittest.TestCase):
    """Đo lại đúng như client vẽ — đây là thứ người chơi thật sự nhìn thấy."""

    @classmethod
    def setUpClass(cls):
        cls.rows = _do_bang()
        cls.rows_x4 = _do_bang(zoom=4)

    def test_ban_chan_moi_khung_cung_mot_cao_do(self):
        """Trước: -10..-26 px (nhân vật nhảy lên xuống khi chạy). Sau: lệch ≤ 1 px."""
        goc = self.rows[0]["chan_so_cy"]
        for cf, row in self.rows.items():
            self.assertLessEqual(
                abs(row["chan_so_cy"] - goc), 1,
                f"cf {cf}: ban chan {row['chan_so_cy']} vs cf0 {goc}",
            )

    def test_chieu_cao_tong_theo_tu_the_trong_gioi_han_da_chap_nhan(self):
        """Đổi kỳ vọng 05/09 (V1): trước là ±3 px vì "một tỉ lệ chung".

        Với tỉ lệ theo cỡ đầu + bộ xương client giữ đầu gần cố định, khung CHẠY nguồn vẽ
        khom nên nhân vật thấp đi ~10 % khi chạy (đo prototype: 53–55 so 61–62 px) —
        giá đã trình (GHI-CHU.md §3) và chủ dự án CHỐT tối 05/09: "chọn 1 đi, làm luôn cả Iruka"
        (BAT-DAU-TU-DAY.md §5 hàng "Khe đầu–thân", NHAT-KY-AGENT.md 06/09). Không chỉ cf 2..6:
        cf 15/16 (thân run0/cast2 + chân melee2) và cf 3 cũng 53 px. Đo sau sinh 05/09:
        lệch so cf 0 trong [−8, +2]. Khoá dải [−10, +4] để không tệ hơn.
        """
        goc = self.rows[0]["cao_tong"]
        for cf, row in self.rows.items():
            lech = row["cao_tong"] - goc
            self.assertTrue(-10 <= lech <= 4, f"cf {cf}: cao tong {row['cao_tong']} vs cf0 {goc}")

    def test_than_va_chan_luon_chong_nhau_khong_bao_gio_ho(self):
        """Ràng buộc MILP: khe thân–chân ∈ [−11, −2] ở x1 (siết −2 để x4 không hở)."""
        for cf, row in self.rows.items():
            self.assertLessEqual(row["khe_than_chan"], -2, f"cf {cf} ho than-chan")
            self.assertGreaterEqual(row["khe_than_chan"], -11, f"cf {cf} chong qua sau")

    def test_khe_dau_than_khong_ho_o_moi_cf(self):
        """Đổi kỳ vọng 05/09 (V1): trước cho phép hở +5 (8 cf) và +19 (cf 20/31 khung nằm)
        vì coi đó là "ràng buộc của nguồn". Chẩn đoán lại: cái hở thật là hai khuôn mặt,
        và với tỉ lệ theo cỡ đầu + slot 15 = cast0 + MILP thì 33/33 cf chồng [−9, 0]
        (đo prototype [−3, −1]). Không còn ngoại lệ."""
        for cf, row in self.rows.items():
            self.assertLessEqual(row["khe_dau_than"], 0, f"cf {cf} ho dau-than")
            self.assertGreaterEqual(row["khe_dau_than"], -9, f"cf {cf} chong qua sau")

    def test_x4_khong_cf_nao_ho(self):
        """Client Unity vẽ ảnh x4 ở toạ độ x1 × 4: làm tròn zoom có thể mở khe 1–3 px
        ở x4 dù x1 chồng. Siết khe thân–chân ≤ −2 chính là để đây xanh."""
        for cf, row in self.rows_x4.items():
            self.assertLessEqual(row["khe_dau_than"], 0, f"cf {cf} ho dau-than o x4")
            self.assertLessEqual(row["khe_than_chan"], 0, f"cf {cf} ho than-chan o x4")


if __name__ == "__main__":
    unittest.main()
