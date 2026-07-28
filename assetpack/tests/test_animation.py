import unittest


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
