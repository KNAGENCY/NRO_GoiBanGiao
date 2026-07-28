import unittest


class ImageContractTests(unittest.TestCase):
    def test_four_zoom_outputs_use_nearest_neighbour_and_exact_resource_paths(self):
        from assetpack.images import build_zooms
        from PIL import Image

        image = Image.new("RGBA", (1, 1), (12, 34, 56, 78))
        zooms = build_zooms(image, "Small27041", {"x1": 1, "x2": 2, "x3": 3, "x4": 4})
        self.assertEqual({key: value.size for key, value in zooms.items()}, {"x1": (1, 1), "x2": (2, 2), "x3": (3, 3), "x4": (4, 4)})
        self.assertEqual(zooms["x4"].getpixel((3, 3)), (12, 34, 56, 78))
