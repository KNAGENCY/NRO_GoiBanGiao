import unittest


def _closure():
    specs = (
        ("hoc-vien-head", 3, (128, 128)),
        ("hoc-vien-body", 17, (128, 160)),
        ("hoc-vien-leg", 14, (128, 128)),
        ("quai-soi-rung", 10, (96, 96)),
    )
    expected_cf = {
        "hoc-vien-head": [[0, 1, 7, 8, 11, 12, 17, 18, 19, 20, 21, 22, 23, 24, 30, 31, 32], [2, 3, 4, 5, 6, 9, 10, 13, 14, 15, 16, 26, 27, 28, 29], [25]],
        "hoc-vien-body": [[0, 1], [2, 15, 28], [3], [4], [5], [6], [7, 8, 17], [9, 10], [11, 21, 32], [12, 18, 22], [13, 26], [14, 27], [16, 29], [19, 30], [20, 31], [23, 24], [25]],
        "hoc-vien-leg": [[0, 1, 22], [2], [3], [4], [5], [6], [7, 23], [8], [9], [10], [11], [12], [13, 14, 15, 16, 17, 18, 19, 20, 21, 24], [25]],
        "quai-soi-rung": [[i] for i in range(10)],
    }
    sheets, animations = [], []
    for name, count, size in specs:
        sheets.append({"name": name, "slot_count": count, "cell_size": list(size), "cells": [{"slot_index": i, "path": f"{name}-{i}.png"} for i in range(count)]})
        animations.append({"name": name, "slot_count": count, "frames": [{"slot_index": i, "client_cf": expected_cf[name][i]} for i in range(count)]})
    return sheets, animations


class AnimationContractTests(unittest.TestCase):
    def test_slot_index_places_cells_in_the_expected_sheet_not_input_order(self):
        from assetpack.animation import compose_action_sheet
        from PIL import Image

        red = Image.new("RGBA", (1, 1), (255, 0, 0, 255))
        blue = Image.new("RGBA", (1, 1), (0, 0, 255, 255))
        green = Image.new("RGBA", (1, 1), (0, 255, 0, 255))
        sheet = compose_action_sheet([(1, blue), (0, red), (2, green)], 3, (1, 1))
        self.assertEqual(sheet.getpixel((0, 0)), (255, 0, 0, 255))
        self.assertEqual(sheet.getpixel((1, 0)), (0, 0, 255, 255))

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
