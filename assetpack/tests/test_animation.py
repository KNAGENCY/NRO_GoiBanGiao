import unittest


def _closure():
    # Ô contact sheet = giới hạn trên (xem REQUIRED_ANIMATION_SHEETS): sprite
    # người chơi không còn bị nhồi vào hộp cố định nên mỗi cell một kích cỡ.
    specs = (
        ("hoc-vien-head", 3, (128, 128)),
        ("hoc-vien-body", 17, (256, 160)),
        ("hoc-vien-leg", 14, (160, 128)),
        ("quai-soi-rung", 10, (96, 96)),
    )
    # Bảng vàng ĐỘC LẬP, chép tay từ `tools/nhan-vat/charinfo.json` (= Char.cs:483) theo
    # đúng chỉ số slot: CharInfo[cf][0][0] đầu, [1][0] chân, [2][0] thân.
    # Đổi 05/09/2026 (V1 khe đầu–thân): bảng cũ viết tay SẮP THEO cf đầu tiên chứ không
    # theo slot — thân slot 0 thật sự phục vụ cf 23/24 (không phải 0/1), slot 1 mới là
    # idle 0/1; chân slot 0 = cf 8, slot 1 = cf 0/1/22, và cf 26..32 (chân slot 8) bị
    # bỏ quên hẳn. Bù dx/dy theo bảng cũ là lệch một khung.
    expected_cf = {
        "hoc-vien-head": [[0, 1, 7, 8, 11, 12, 17, 18, 19, 20, 21, 22, 23, 24, 30, 31, 32], [2, 3, 4, 5, 6, 9, 10, 13, 14, 15, 16, 26, 27, 28, 29], [25]],
        "hoc-vien-body": [[23, 24], [0, 1], [2, 15, 28], [3], [4], [5], [6], [7, 8, 17], [12, 18, 22], [11, 21, 32], [13, 26], [14, 27], [9, 10], [16, 29], [19, 30], [20, 31], [25]],
        "hoc-vien-leg": [[8], [0, 1, 22], [2], [3], [4], [5], [6], [7, 23], [12, 26, 27, 28, 29, 30, 31, 32], [13, 14, 15, 16, 17, 18, 19, 20, 21, 24], [9], [10], [11], [25]],
        "quai-soi-rung": [[i] for i in range(10)],
    }
    sheets, animations = [], []
    for name, count, size in specs:
        sheets.append({"name": name, "slot_count": count, "cell_size": list(size), "cells": [{"slot_index": i, "path": f"{name}-{i}.png"} for i in range(count)]})
        animations.append({"name": name, "slot_count": count, "frames": [{"slot_index": i, "client_cf": expected_cf[name][i]} for i in range(count)]})
    return sheets, animations


class AnimationContractTests(unittest.TestCase):
    def test_required_client_cf_sinh_dung_tu_charinfo(self):
        """Bảng slot→cf phải là cái client THẬT SỰ tra (CharInfo), không phải bảng viết tay.

        Bảng cũ sắp theo cf đầu tiên nên thân/chân lệch slot và chân thiếu hẳn cf 26..32
        (xem chú thích ở `_closure`). Test này khoá cả hai: bảng module == bảng vàng chép
        tay ở đây, và == suy ra trực tiếp từ charinfo.json.
        """
        import json
        from pathlib import Path

        from assetpack.animation import REQUIRED_CLIENT_CF, charinfo_client_cf

        _, animations = _closure()
        for item in animations:
            for frame in item["frames"]:
                self.assertEqual(
                    list(REQUIRED_CLIENT_CF[item["name"]][frame["slot_index"]]),
                    frame["client_cf"],
                    f"{item['name']} slot {frame['slot_index']}",
                )
        charinfo = json.loads(
            (Path(__file__).resolve().parents[2] / "tools/nhan-vat/charinfo.json").read_text(encoding="utf-8")
        )
        derived = charinfo_client_cf(charinfo)
        for name in ("hoc-vien-head", "hoc-vien-body", "hoc-vien-leg"):
            self.assertEqual(REQUIRED_CLIENT_CF[name], derived[name], name)
            # mỗi cf xuất hiện đúng một lần trong mỗi bảng (33 cf, không slot nào chết)
            flat = [cf for cfs in REQUIRED_CLIENT_CF[name] for cf in cfs]
            self.assertEqual(sorted(flat), list(range(33)), name)
        self.assertEqual(REQUIRED_CLIENT_CF["hoc-vien-body"][0], (23, 24))
        self.assertEqual(REQUIRED_CLIENT_CF["hoc-vien-body"][1], (0, 1))
        self.assertEqual(REQUIRED_CLIENT_CF["hoc-vien-leg"][0], (8,))
        self.assertEqual(REQUIRED_CLIENT_CF["hoc-vien-leg"][8], (12, 26, 27, 28, 29, 30, 31, 32))

    def test_slot_index_places_cells_in_the_expected_sheet_not_input_order(self):
        from assetpack.animation import compose_action_sheet
        from PIL import Image

        red = Image.new("RGBA", (1, 1), (255, 0, 0, 255))
        blue = Image.new("RGBA", (1, 1), (0, 0, 255, 255))
        green = Image.new("RGBA", (1, 1), (0, 255, 0, 255))
        sheet = compose_action_sheet([(1, blue), (0, red), (2, green)], 3, (1, 1))
        self.assertEqual(sheet.getpixel((0, 0)), (255, 0, 0, 255))
        self.assertEqual(sheet.getpixel((1, 0)), (0, 0, 255, 255))

    def test_cell_nho_hon_o_duoc_dan_giua_duoi_va_cell_vuot_o_bi_tu_choi(self):
        """Ô là giới hạn trên: mỗi khung một kích cỡ nên cell không còn bằng ô."""
        from assetpack.animation import AnimationError, compose_action_sheet
        from PIL import Image

        nho = Image.new("RGBA", (2, 1), (255, 0, 0, 255))
        sheet = compose_action_sheet([(0, nho), (1, nho), (2, nho)], 3, (4, 2))
        self.assertEqual(sheet.size, (12, 2))
        # dán giữa-dưới: x = 0*4 + (4-2)//2 = 1, y = 2-1 = 1
        self.assertEqual(sheet.getpixel((1, 1)), (255, 0, 0, 255))
        self.assertEqual(sheet.getpixel((1, 0)), (0, 0, 0, 0))

        to = Image.new("RGBA", (5, 1), (255, 0, 0, 255))
        with self.assertRaisesRegex(AnimationError, "ASSET-ANIMATION-008"):
            compose_action_sheet([(0, to), (1, to), (2, to)], 3, (4, 2))

    def test_charinfo_cf_is_an_independent_golden_contract(self):
        from assetpack.animation import AnimationError, validate_animation_closure

        sheets, animations = _closure()
        animations[0]["frames"][0]["client_cf"] = [99]
        with self.assertRaisesRegex(AnimationError, "ASSET-ANIMATION-008"):
            validate_animation_closure(sheets, animations)

    def test_duplicate_sheet_or_animation_name_is_rejected(self):
        from assetpack.animation import AnimationError, validate_animation_closure

        sheets, animations = _closure()
        sheets[-1]["name"] = sheets[0]["name"]
        with self.assertRaisesRegex(AnimationError, "ASSET-ANIMATION-008"):
            validate_animation_closure(sheets, animations)

        sheets, animations = _closure()
        animations[-1]["name"] = animations[0]["name"]
        with self.assertRaisesRegex(AnimationError, "ASSET-ANIMATION-008"):
            validate_animation_closure(sheets, animations)

    def test_malformed_sheet_or_animation_name_types_return_animation_error(self):
        from assetpack.animation import AnimationError, validate_animation_closure

        for collection in ("sheets", "animations"):
            for bad_name in ([], {}, None, 7):
                with self.subTest(collection=collection, bad_name=bad_name):
                    sheets, animations = _closure()
                    target = sheets if collection == "sheets" else animations
                    target[0]["name"] = bad_name
                    with self.assertRaisesRegex(AnimationError, "ASSET-ANIMATION-008"):
                        validate_animation_closure(sheets, animations)

    def test_slot_count_must_be_present_integer_and_canonical_for_both_closures(self):
        from assetpack.animation import AnimationError, validate_animation_closure

        for collection in ("sheets", "animations"):
            for mutation in ("missing", "bool", "float", "wrong"):
                with self.subTest(collection=collection, mutation=mutation):
                    sheets, animations = _closure()
                    target = sheets[0] if collection == "sheets" else animations[0]
                    if mutation == "missing":
                        target.pop("slot_count")
                    elif mutation == "bool":
                        target["slot_count"] = True
                    elif mutation == "float":
                        target["slot_count"] = 3.0
                    else:
                        target["slot_count"] = 4
                    with self.assertRaisesRegex(AnimationError, "ASSET-ANIMATION-008"):
                        validate_animation_closure(sheets, animations)
