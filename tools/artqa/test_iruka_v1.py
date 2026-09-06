# -*- coding: utf-8 -*-
"""Test Iruka V1 (06/09/2026) — NPC Iruka (part 2189 đầu / 2190 thân / 2191 chân,
sprite 27075/27076/27077) sinh cùng hình học với người chơi.

VÌ SAO CÓ TEST NÀY
------------------
Bản đang chạy sinh Iruka bằng ``_generate_iruka_components`` (hình học CŨ): ép cả khung
vào canvas 128×160 nên Iruka chỉ 38 px x1 = 62 % người chơi, và part 2189–2191 ghi
dx = dy = 0 → đầu lửng +13 px, thân chìm +29 px (GHI-CHU 05/09). Bản V1 dùng đúng bộ hàm
của người chơi (``_component_rect_v1`` cắt theo cằm, ``component_placement`` bù CharInfo
cf 0) với tỉ lệ s = s0 / r × IRUKA_TI_LE_DAU, r = cỡ đầu Iruka / cỡ đầu người chơi đo
bằng median 6 thước cứng trên CÙNG tham chiếu trung bình 3 khung idle mà agent A dùng.

Ba lớp test:
  1. toán thuần (không cần ảnh nguồn) — hệ số đầu, lát cắt, residual CharInfo;
  2. sinh từ nguồn (cần kho ảnh ``docs/my-thuat/anh-nguon``) — ghép bằng đúng phép vẽ
     client (``tools/nhan-vat/ghep-khung.py``) cho cf 0/1, so chiều cao với người chơi;
  3. artifact trong repo (mapping ↔ part.bin ↔ ảnh) — chạy được ở mọi checkout.
"""

import contextlib
import hashlib
import importlib.util
import io
import json
import math
import struct
import tempfile
import unittest
from pathlib import Path

from PIL import Image, ImageOps

ROOT = Path(__file__).resolve().parents[2]
ART_ROOT = ROOT / "docs/my-thuat/anh-nguon/vertical-slice-v1"
REVIEW = ROOT / "docs/my-thuat/review/vertical-slice-v1.json"
MAPPING = ROOT / "docs/my-thuat/mapping-id.json"
PART_BIN = ROOT / "server/internal/game/assets/part.bin"
SPRITE_ROOT = ROOT / "docs/my-thuat/lang-la-client-assets/res"
CHARINFO = ROOT / "tools/nhan-vat/charinfo.json"
GHEP_KHUNG = ROOT / "tools/nhan-vat/ghep-khung.py"

# Cửa nghiệm thu Iruka (đơn vị x1), cùng dải với người chơi: khe đầu–thân ∈ [−9, −1]
# (siết trên −1 để x4 không hở), thân–chân ∈ [−11, −2], bàn chân = cy − 11, cao 62–66.
KHE_DAU_THAN = (-9, -1)
KHE_THAN_CHAN = (-11, -2)
BAN_CHAN = -11
CAO = (62, 66)
CO_NGUON = (ART_ROOT / "master_iruka_x4_chroma.png").is_file() and (
    ART_ROOT / "player_idle_3f_chroma.png"
).is_file()


def _nap_ghep_khung():
    spec = importlib.util.spec_from_file_location("ghep_khung", GHEP_KHUNG)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def _ghi_part_bin(duong, parts):
    """parts: [(type, [(sid, dx, dy), ...]), ...] → nhị phân client (>h n; b type + >hbb*k)."""
    d = bytearray(struct.pack(">h", len(parts)))
    for t, fr in parts:
        d += bytes([t])
        for sid, dx, dy in fr:
            d += struct.pack(">hbb", sid, dx, dy)
    Path(duong).write_bytes(bytes(d))


def _do_cf(g, charinfo, parts, thu_muc_sprite, ids, cf, zoom):
    """Ghép một cf như client và trả (ảnh, vùng, số đo + so_mat)."""
    cx, cy = 48 * zoom, 96 * zoom
    anh, vung = g.ghep(cf, charinfo, parts, str(thu_muc_sprite), ids, zoom, cx, cy)
    so_do = g.do_khung(vung, cy)
    so_do["so_mat"] = g.so_mat(anh, vung, zoom)
    return anh, vung, so_do


class ToanIrukaTests(unittest.TestCase):
    """Toán thuần — không cần ảnh nguồn."""

    def test_hang_so_iruka(self):
        from tools.artqa import task6b_closure as t

        self.assertEqual(t.IRUKA_TI_LE_DAU, 1.0)
        self.assertEqual(t.IRUKA_SOURCE, "master_iruka_x4_chroma.png")
        self.assertEqual(t.IRUKA_IDS, (27075, 27076, 27077))
        self.assertEqual(t.IRUKA_PART_IDS, {"head": 2189, "body": 2190, "leg": 2191})
        self.assertEqual(t.IRUKA_CHARINFO_CF, (0, 1))
        # 6 thước: bỏ head_h_hair_chin (đuôi tóc Iruka), thay bằng sqrt diện tích mặt.
        self.assertEqual(len(t.IRUKA_THUOC_CUNG), 6)
        self.assertNotIn("head_h_hair_chin", t.IRUKA_THUOC_CUNG)
        self.assertIn("face_area_sqrt", t.IRUKA_THUOC_CUNG)

    def test_he_so_dau_iruka_la_median_6_thuoc_tren_trung_binh_idle(self):
        from tools.artqa.task6b_closure import iruka_he_so_dau

        idle = {
            "face_h": 100, "face_w": 100, "face_area": 10000, "eye_h": 10,
            "eye_to_chin": 50, "band_thick": 10, "head_h_hair_chin": 200,
        }
        # ba khung idle lệch nhau đối xứng → trung bình đúng bằng `idle`
        moc = {
            ("player_idle_3f_chroma.png", 0): {**idle, "face_h": 90, "face_area": 8100},
            ("player_idle_3f_chroma.png", 1): {**idle, "face_h": 110, "face_area": 12100},
            ("player_idle_3f_chroma.png", 2): dict(idle),
        }
        # Iruka: 1.7 / 1.5 / sqrt(25600)=160 → 1.6 / 1.55 / 1.7 / 2.0 ; tóc 2.5 (phải bị bỏ)
        iruka = {
            "face_h": 170, "face_w": 150, "face_area": 25600, "eye_h": 15.5,
            "eye_to_chin": 85, "band_thick": 20, "head_h_hair_chin": 500,
        }
        kq = iruka_he_so_dau(iruka, moc)
        # sorted: 1.5 1.55 1.6 1.7 1.7 2.0 → median (1.6 + 1.7) / 2
        self.assertAlmostEqual(kq["he_so"], 1.65, places=9)
        self.assertAlmostEqual(kq["tham_chieu"]["face_h"], 100.0)
        self.assertAlmostEqual(kq["tham_chieu"]["face_area_sqrt"], (90 + 110 + 100) / 3)
        self.assertEqual(kq["idle_frames"], [0, 1, 2])
        self.assertNotIn("head_h_hair_chin", kq["ti_le_tung_thuoc"])
        self.assertAlmostEqual(kq["ti_le_tung_thuoc"]["band_thick"], 2.0)

    def test_he_so_dau_iruka_tu_choi_khung_khong_co_mat(self):
        from tools.artqa.task6b_closure import ClosureError, iruka_he_so_dau

        idle = {"face_h": 100, "face_w": 100, "face_area": 10000, "eye_h": 10, "eye_to_chin": 50, "band_thick": 10}
        moc = {("player_idle_3f_chroma.png", i): dict(idle) for i in range(3)}
        with self.assertRaises(ClosureError):
            iruka_he_so_dau({"content_h": 1104}, moc)

    def test_lat_cat_iruka_theo_cam_ra_dung_ba_hinh_chu_nhat(self):
        """Khung master_iruka [369,70,481,1104], cằm y=438 (rel 368), s ≈ 0,2262:
        ĐẦU [0, cằm+6] → [369,70,481,374]; THÂN [cằm − 8/s, 0.76H] → [369,403,481,506];
        CHÂN [0.66H, H] → [369,799,481,375]."""
        from tools.artqa.task6b_closure import _component_rect_v1

        record = {
            "asset": "master_iruka_x4_chroma.png",
            "checks": {"crop": {"source_edges": [0, 1254], "content_boxes": [[369, 70, 850, 1174]]}},
        }
        khung = {"chin_y": 438.0, "scale": 0.3764705882352941 / 1.6642}
        self.assertEqual(_component_rect_v1(record, 0, "head", khung), [369, 70, 481, 374])
        self.assertEqual(_component_rect_v1(record, 0, "body", khung), [369, 403, 481, 506])
        self.assertEqual(_component_rect_v1(record, 0, "leg", khung), [369, 799, 481, 375])

    def test_charinfo_residual_cf1_nang_dau_than_1px_chan_dung_yen(self):
        from tools.artqa.task6b_closure import iruka_charinfo_residual, load_charinfo

        charinfo = load_charinfo(CHARINFO)
        self.assertEqual(iruka_charinfo_residual(charinfo, "head"), {"0": [0, 0], "1": [0, -1]})
        self.assertEqual(iruka_charinfo_residual(charinfo, "body"), {"0": [0, 0], "1": [0, -1]})
        self.assertEqual(iruka_charinfo_residual(charinfo, "leg"), {"0": [0, 0], "1": [0, 0]})

    def test_iruka_bang_cf_du_doan_dung_phep_ve_client(self):
        """Dự đoán x1 từ bbox + dx/dy phải trùng công thức Char.paintCharBody."""
        from tools.artqa.task6b_closure import iruka_bang_cf

        charinfo = [
            [[0, -13, 34], [1, -8, 10], [1, -9, 16]],
            [[0, -13, 35], [1, -8, 10], [1, -9, 17]],
        ]

        def manh(component, bbox, dx, dy):
            return {"component": component, "moc": {"bbox_x1": list(bbox)}, "placement": {"dx": dx, "dy": dy}}

        pieces = [
            manh("head", (0, 0, 27, 22), -1, -41),   # top = -34-41 = -75, bottom = -53
            manh("body", (0, 0, 26, 29), -3, -39),   # top = -16-39 = -55, bottom = -26
            manh("leg", (0, 0, 24, 22), -2, -23),    # top = -10-23 = -33, bottom = -11
        ]
        bang = iruka_bang_cf(pieces, charinfo, cfs=(0, 1))
        self.assertEqual([r["cf"] for r in bang], [0, 1])
        self.assertEqual(bang[0]["khe_dau_than"], -55 - (-53))
        self.assertEqual(bang[0]["khe_than_chan"], -33 - (-26))
        self.assertEqual(bang[0]["chan"], -11)
        self.assertEqual(bang[0]["cao"], -11 - (-75))
        # cf 1: đầu và thân cùng lên 1 px → khe đầu–thân giữ, thân–chân bớt chồng 1, cao +1
        self.assertEqual(bang[1]["khe_dau_than"], -2)
        self.assertEqual(bang[1]["khe_than_chan"], -6)
        self.assertEqual(bang[1]["chan"], -11)
        self.assertEqual(bang[1]["cao"], 65)


@unittest.skipUnless(CO_NGUON, "cần kho ảnh nguồn docs/my-thuat/anh-nguon/vertical-slice-v1")
class SinhIrukaTuNguonTests(unittest.TestCase):
    """Sinh 3 mảnh từ master_iruka rồi ghép như client — không ghi gì vào repo."""

    @classmethod
    def setUpClass(cls):
        from tools.artqa import nhan_vat_moc
        from tools.artqa.task6b_closure import (
            _review_by_name,
            iruka_component_images,
            load_charinfo,
        )

        cls.g = _nap_ghep_khung()
        cls.review = json.loads(REVIEW.read_text(encoding="utf-8"))
        cls.records = _review_by_name(cls.review)
        cls.charinfo = load_charinfo(CHARINFO)
        cls.moc = nhan_vat_moc.do_nguon(ART_ROOT, cls.records)
        cls.iruka = iruka_component_images(
            cls.review, art_root=ART_ROOT, records=cls.records, charinfo=cls.charinfo, moc=cls.moc
        )
        cls.tam = tempfile.TemporaryDirectory()
        cls.thu_muc = Path(cls.tam.name)
        # part.bin thử: id 0 đầu / 1 thân / 2 chân, MỌI frame cùng dx/dy (như 2189/2190/2191)
        by = {p["component"]: p for p in cls.iruka}
        _ghi_part_bin(cls.thu_muc / "part.bin", [
            (0, [(27075, by["head"]["placement"]["dx"], by["head"]["placement"]["dy"])] * 3),
            (1, [(27076, by["body"]["placement"]["dx"], by["body"]["placement"]["dy"])] * 17),
            (2, [(27077, by["leg"]["placement"]["dx"], by["leg"]["placement"]["dy"])] * 14),
        ])
        from tools.artqa.task6b_closure import zoom_variants

        for p in cls.iruka:
            for zoom, im in zoom_variants(p["image_x4"]).items():
                (cls.thu_muc / zoom).mkdir(exist_ok=True)
                im.save(cls.thu_muc / zoom / f"Small{p['id']}.png")
        cls.parts = cls.g.doc_part_bin(str(cls.thu_muc / "part.bin"))

    @classmethod
    def tearDownClass(cls):
        cls.tam.cleanup()

    def test_tra_dung_3_manh_theo_thu_tu_dau_than_chan(self):
        self.assertEqual([p["id"] for p in self.iruka], [27075, 27076, 27077])
        self.assertEqual([p["component"] for p in self.iruka], ["head", "body", "leg"])
        self.assertEqual([p["part_id"] for p in self.iruka], [2189, 2190, 2191])
        for p in self.iruka:
            self.assertEqual(p["source_path"], "master_iruka_x4_chroma.png")
            self.assertEqual(p["charinfo_cf"], [0, 1])
            self.assertEqual(p["charinfo_anchor_cf"], 0)
            self.assertIsInstance(p["image_x4"], Image.Image)
            self.assertIsNotNone(p["image_x4"].getchannel("A").getbbox())
            for k in ("left_rel", "top_rel", "dx", "dy"):
                self.assertIsInstance(p["placement"][k], int)

    def test_ti_le_dau_bang_nguoi_choi_cung_tham_chieu_trung_binh_idle(self):
        st = self.iruka[0]["shared_transform"]
        self.assertEqual(st["ti_le_dau"], 1.0)
        # đo 06/09: r = 1,6642 trên trung bình 3 idle (idle0 lẻ: 1,6382)
        self.assertAlmostEqual(st["he_so_dau"], 1.6642, places=3)
        self.assertEqual(st["tham_chieu"]["kieu"], "trung_binh_idle")
        self.assertEqual(st["tham_chieu"]["idle_frames"], [0, 1, 2])
        self.assertAlmostEqual(st["uniform_scale_x4"], 128 / 340)
        self.assertAlmostEqual(st["scale"], st["uniform_scale_x4"] / st["he_so_dau"] * st["ti_le_dau"])
        for p in self.iruka:
            self.assertEqual(p["shared_transform"]["scale"], st["scale"])
            self.assertEqual(p["scale_frame"], st["scale"])

    def test_lat_cat_that_khop_ba_hinh_chu_nhat_da_chot(self):
        rects = {p["component"]: p["source_rect"] for p in self.iruka}
        self.assertEqual(rects["head"], [369, 70, 481, 374])
        self.assertEqual(rects["body"], [369, 403, 481, 506])
        self.assertEqual(rects["leg"], [369, 799, 481, 375])
        # thân không chứa mặt: cằm ≤ 2,5 px x1 dưới mép trên vùng đục
        than = next(p for p in self.iruka if p["component"] == "body")
        from tools.artqa.task6b_closure import BODY_CHIN_MARGIN_X1  # ngưỡng theo hằng thiết kế, không phải số ma thuật
        self.assertLessEqual(than["moc"]["chin_off"], BODY_CHIN_MARGIN_X1 + 0.5)
        self.assertGreaterEqual(than["moc"]["chin_off"], -1.0)

    def test_part_offset_dung_cong_thuc_component_placement_cf0(self):
        from tools.artqa.task6b_closure import COMPONENT_CHARINFO_SLOT, PLAYER_FEET_OFFSET_X1, component_placement

        for p in self.iruka:
            st = p["shared_transform"]
            ox, oy = st["component_origin_x4"]
            mong = component_placement(
                component_left_x4=ox,
                component_top_x4=oy,
                frame_center_x4=st["frame_center_x4"],
                frame_feet_x4=st["frame_feet_x4"],
                charinfo_row=self.charinfo[0][COMPONENT_CHARINFO_SLOT[p["component"]]],
                feet_offset_x1=PLAYER_FEET_OFFSET_X1,
            )
            self.assertEqual(p["placement"], mong, p["component"])
            self.assertEqual(st["feet_offset_x1"], PLAYER_FEET_OFFSET_X1)

    def test_ghep_cf0_cf1_x1_dat_cua(self):
        for cf in (0, 1):
            _, _, so_do = _do_cf(self.g, self.charinfo, self.parts, self.thu_muc / "x1", (0, 1, 2), cf, 1)
            with self.subTest(cf=cf):
                self.assertTrue(KHE_DAU_THAN[0] <= so_do["khe_dt"] <= KHE_DAU_THAN[1], so_do)
                self.assertTrue(KHE_THAN_CHAN[0] <= so_do["khe_tc"] <= KHE_THAN_CHAN[1], so_do)
                self.assertEqual(so_do["chan_cy"], BAN_CHAN, so_do)
                self.assertTrue(CAO[0] <= so_do["cao_tong"] <= CAO[1], so_do)
                self.assertEqual(so_do["so_mat"], 1, so_do)

    def test_ghep_x4_khong_khe_duong(self):
        for cf in (0, 1):
            _, _, so_do = _do_cf(self.g, self.charinfo, self.parts, self.thu_muc / "x4", (0, 1, 2), cf, 4)
            with self.subTest(cf=cf):
                self.assertLessEqual(so_do["khe_dt"], 0, so_do)
                self.assertLessEqual(so_do["khe_tc"], 0, so_do)
                self.assertTrue(-47 <= so_do["chan_cy"] <= -41, so_do)  # −44 ± làm tròn zoom
                self.assertEqual(so_do["so_mat"], 1, so_do)

    def test_bang_cf_du_doan_khop_ghep_that(self):
        """Bảng dự đoán trong spec (ghi vào mapping) phải trùng số đo ghép — mapping tự chứng minh."""
        bang = {r["cf"]: r for r in self.iruka[0]["bang_cf"]}
        for cf in (0, 1):
            _, _, so_do = _do_cf(self.g, self.charinfo, self.parts, self.thu_muc / "x1", (0, 1, 2), cf, 1)
            self.assertEqual(bang[cf]["khe_dau_than"], so_do["khe_dt"], cf)
            self.assertEqual(bang[cf]["khe_than_chan"], so_do["khe_tc"], cf)
            self.assertEqual(bang[cf]["chan"], so_do["chan_cy"], cf)
            self.assertEqual(bang[cf]["cao"], so_do["cao_tong"], cf)

    def test_cao_iruka_bang_nguoi_choi_trong_3px(self):
        """Người chơi đo trên artifact repo (part 30/57/58 + ảnh x1) tại cf 0."""
        parts_nc = self.g.doc_part_bin(str(PART_BIN))
        _, _, nc = _do_cf(self.g, self.charinfo, parts_nc, SPRITE_ROOT / "x1" / "smallimage", (30, 57, 58), 0, 1)
        _, _, ir = _do_cf(self.g, self.charinfo, self.parts, self.thu_muc / "x1", (0, 1, 2), 0, 1)
        self.assertIsNotNone(nc["cao_tong"])
        self.assertLessEqual(abs(ir["cao_tong"] - nc["cao_tong"]), 3, (ir["cao_tong"], nc["cao_tong"]))
        self.assertEqual(ir["chan_cy"], nc["chan_cy"])

    def test_guong_trai_khop_0px(self):
        """Npc.cs:305–316: hướng trái vẽ x = cx − CI_dx − dx, lật ngang, neo TOP|RIGHT.
        Ảnh trim alpha nên bản trái phải là gương đúng từng pixel của bản phải."""
        cx, cy = 48, 96
        phai = Image.new("RGBA", (cx * 2, cy + 16), (0, 0, 0, 0))
        trai = Image.new("RGBA", (cx * 2, cy + 16), (0, 0, 0, 0))
        for k, pid in ((0, 0), (1, 2), (2, 1)):  # thứ tự vẽ client: đầu, chân, thân
            idx, dx_ci, dy_ci = self.charinfo[0][k]
            sid, dx, dy = self.parts[pid][1][idx]
            im = Image.open(self.thu_muc / "x1" / f"Small{sid}.png").convert("RGBA")
            y = cy - dy_ci + dy
            phai.alpha_composite(im, (cx + dx_ci + dx, y))
            # neo TOP|RIGHT tại x = cx − CI_dx − dx → mép trái = x − w
            trai.alpha_composite(ImageOps.mirror(im), (cx - dx_ci - dx - im.width, y))
        self.assertEqual(ImageOps.mirror(phai).tobytes(), trai.tobytes())


class ArtifactIrukaTests(unittest.TestCase):
    """mapping ↔ part.bin ↔ ảnh thật trong repo — chạy được ở mọi checkout."""

    @classmethod
    def setUpClass(cls):
        cls.mapping = json.loads(MAPPING.read_text(encoding="utf-8"))
        cls.geometry = cls.mapping["build_inputs"]["iruka_component_geometry"]
        cls.charinfo = json.loads(CHARINFO.read_text(encoding="utf-8"))
        cls.g = _nap_ghep_khung()
        cls.parts = cls.g.doc_part_bin(str(PART_BIN))

    def test_khoi_iruka_geometry_rieng_khong_noi_vao_player(self):
        self.assertEqual(len(self.mapping["build_inputs"]["player_component_geometry"]["components"]), 34)
        self.assertEqual(len(self.geometry["components"]), 3)
        self.assertEqual([c["id"] for c in self.geometry["components"]], [27075, 27076, 27077])
        self.assertEqual([c["part_id"] for c in self.geometry["components"]], [2189, 2190, 2191])
        self.assertEqual(self.geometry["ti_le_dau"], 1.0)
        self.assertEqual(self.geometry["tham_chieu"]["kieu"], "trung_binh_idle")
        self.assertTrue(1.60 <= self.geometry["he_so_dau"] <= 1.72, self.geometry["he_so_dau"])
        self.assertEqual(self.geometry["feet_offset_x1"], -11)
        self.assertAlmostEqual(self.geometry["scale_x4"], self.geometry["uniform_scale_x4"] / self.geometry["he_so_dau"])
        self.assertIn("V1", self.geometry["phien_ban"])
        # sprites[] 27075–27077 vẫn trỏ closure v11 (reviewer chưa ký closure mới) — ghi chú phải nói rõ
        self.assertIn("closure", self.geometry["ghi_chu"].lower())
        for sprite in self.mapping["sprites"]:
            if sprite["id"] in (27075, 27076, 27077):
                self.assertIn("production-closure-v11", sprite["zooms"]["x4"])

    def test_moi_frame_part_2189_2190_2191_cung_dx_dy_khop_geometry_va_part_bin(self):
        khai = {c["part_id"]: c for c in self.geometry["components"]}
        for part in self.mapping["parts"]:
            if part["id"] not in khai:
                continue
            c = khai[part["id"]]
            thuc = self.parts[part["id"]]
            self.assertEqual(len(thuc[1]), len(part["frames"]), part["id"])
            for i, frame in enumerate(part["frames"]):
                self.assertEqual(frame["sprite_id"], c["id"], (part["id"], i))
                self.assertEqual((frame["dx"], frame["dy"]), (c["part_offset"]["dx"], c["part_offset"]["dy"]), (part["id"], i))
                self.assertEqual(thuc[1][i], (c["id"], frame["dx"], frame["dy"]), (part["id"], i))
            self.assertNotEqual((c["part_offset"]["dx"], c["part_offset"]["dy"]), (0, 0), "dx=dy=0 là bản lỗi cũ")

    def test_anh_that_du_4_zoom_dung_kich_thuoc_va_sha(self):
        for c in self.geometry["components"]:
            for zoom, relative in c["paths"].items():
                path = ROOT / relative
                self.assertTrue(path.is_file(), relative)
                self.assertEqual(hashlib.sha256(path.read_bytes()).hexdigest(), c["sha256"][zoom], relative)
                with Image.open(path) as im:
                    size = list(im.size)
                if zoom == "x4":
                    self.assertEqual(size, c["size_x4"], relative)
                else:
                    k = int(zoom[1])
                    self.assertEqual(size, [max(1, round(c["size_x4"][0] * k / 4)), max(1, round(c["size_x4"][1] * k / 4))], relative)
            # không còn hộp 128×160 cố định
            self.assertNotEqual(tuple(c["size_x4"]), (128, 160), c["id"])

    def test_ghep_artifact_cf0_cf1_dat_cua_x1_va_x4(self):
        ids = (2189, 2190, 2191)
        for cf in (0, 1):
            _, _, so_do = _do_cf(self.g, self.charinfo, self.parts, SPRITE_ROOT / "x1" / "smallimage", ids, cf, 1)
            with self.subTest(cf=cf, zoom=1):
                self.assertTrue(KHE_DAU_THAN[0] <= so_do["khe_dt"] <= KHE_DAU_THAN[1], so_do)
                self.assertTrue(KHE_THAN_CHAN[0] <= so_do["khe_tc"] <= KHE_THAN_CHAN[1], so_do)
                self.assertEqual(so_do["chan_cy"], BAN_CHAN, so_do)
                self.assertTrue(CAO[0] <= so_do["cao_tong"] <= CAO[1], so_do)
                self.assertEqual(so_do["so_mat"], 1, so_do)
            _, _, x4 = _do_cf(self.g, self.charinfo, self.parts, SPRITE_ROOT / "x4" / "smallimage", ids, cf, 4)
            with self.subTest(cf=cf, zoom=4):
                self.assertLessEqual(x4["khe_dt"], 0, x4)
                self.assertLessEqual(x4["khe_tc"], 0, x4)
                self.assertEqual(x4["so_mat"], 1, x4)

    def test_bang_cf_trong_mapping_khop_do_that(self):
        bang = {r["cf"]: r for r in self.geometry["bang_cf"]}
        self.assertEqual(sorted(bang), [0, 1])
        for cf in (0, 1):
            _, _, so_do = _do_cf(self.g, self.charinfo, self.parts, SPRITE_ROOT / "x1" / "smallimage", (2189, 2190, 2191), cf, 1)
            self.assertEqual(bang[cf]["khe_dau_than"], so_do["khe_dt"], cf)
            self.assertEqual(bang[cf]["khe_than_chan"], so_do["khe_tc"], cf)
            self.assertEqual(bang[cf]["chan"], so_do["chan_cy"], cf)
            self.assertEqual(bang[cf]["cao"], so_do["cao_tong"], cf)

    def test_cao_iruka_bang_nguoi_choi_tren_artifact(self):
        _, _, nc = _do_cf(self.g, self.charinfo, self.parts, SPRITE_ROOT / "x1" / "smallimage", (30, 57, 58), 0, 1)
        _, _, ir = _do_cf(self.g, self.charinfo, self.parts, SPRITE_ROOT / "x1" / "smallimage", (2189, 2190, 2191), 0, 1)
        self.assertLessEqual(abs(ir["cao_tong"] - nc["cao_tong"]), 3, (ir["cao_tong"], nc["cao_tong"]))
        self.assertEqual(self.geometry["cao_x1"], ir["cao_tong"])


if __name__ == "__main__":
    unittest.main()
