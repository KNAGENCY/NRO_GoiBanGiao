"""Kiểm thử hợp đồng cho QA ảnh nguồn Task 6A.

Mỗi test nêu đột biến nó phải bắt: bỏ alpha chroma hoặc lẫn ranh ô sprite.
"""

from pathlib import Path
import tempfile
import unittest

from PIL import Image

from tools.artqa.artqa import check_cell_boundaries, extract_alpha


class AlphaExtractionTests(unittest.TestCase):
    def test_border_connected_magenta_becomes_transparent_with_soft_edge(self):
        """Bắt đột biến bỏ flood/alpha: nền magenta nối biên phải trong suốt."""
        image = Image.new("RGB", (5, 5), (255, 0, 255))
        image.putpixel((2, 2), (30, 80, 140))

        result = extract_alpha(image)

        self.assertEqual(result.getpixel((0, 0))[3], 0)
        self.assertEqual(result.getpixel((2, 2))[3], 255)
        self.assertLess(result.getpixel((1, 2))[3], 255)


class CellBoundaryTests(unittest.TestCase):
    def test_rejects_non_divisible_sprite_strip_width(self):
        """Bắt đột biến bỏ kiểm tra cell: strip 10px không thể chia 3 frame."""
        image = Image.new("RGBA", (10, 4), (0, 0, 0, 0))

        issues = check_cell_boundaries(image, columns=3)

        self.assertIn("canvas_width_not_divisible_by_columns", issues)


if __name__ == "__main__":
    unittest.main()
