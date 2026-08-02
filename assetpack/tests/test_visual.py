import tempfile
import unittest
from pathlib import Path

from PIL import Image


class VisualContractTests(unittest.TestCase):
    def test_background_requires_all_four_layers(self):
        from assetpack.visual import VisualError, validate_background_layers

        with self.assertRaisesRegex(VisualError, "ASSET-BG-005"):
            validate_background_layers({"b00": object(), "b01": object(), "b02": object()})

    def test_hud_golden_clips_hp_and_chakra_at_25_50_100_percent(self):
        from assetpack.visual import render_hud_golden

        hud = {
            "naruto-hp-bg": Image.new("RGBA", (8, 2), (40, 40, 40, 255)),
            "naruto-hp-fill": Image.new("RGBA", (8, 2), (220, 20, 20, 255)),
            "naruto-chakra-bg": Image.new("RGBA", (8, 2), (60, 60, 60, 255)),
            "naruto-chakra-fill": Image.new("RGBA", (8, 2), (20, 80, 220, 255)),
        }
        with tempfile.TemporaryDirectory() as temporary:
            output = Path(temporary) / "hud-golden.png"
            render_hud_golden(hud, output)
            with Image.open(output) as golden:
                golden.load()
                self.assertEqual(golden.size, (24, 4))
                expected_fill_widths = (2, 4, 8)
                for state_index, fill_width in enumerate(expected_fill_widths):
                    left = state_index * 8
                    self.assertEqual(
                        golden.getpixel((left + fill_width - 1, 0)),
                        (220, 20, 20, 255),
                    )
                    if fill_width < 8:
                        self.assertEqual(
                            golden.getpixel((left + fill_width, 0)),
                            (40, 40, 40, 255),
                        )
                    self.assertEqual(
                        golden.getpixel((left + fill_width - 1, 2)),
                        (20, 80, 220, 255),
                    )
                    if fill_width < 8:
                        self.assertEqual(
                            golden.getpixel((left + fill_width, 2)),
                            (60, 60, 60, 255),
                        )

    def test_iruka_contact_sheet_composites_all_layers_on_light_and_dark(self):
        from assetpack.visual import render_iruka_composite_contact_sheet

        layers = {
            "leg": Image.new("RGBA", (4, 4), (0, 0, 0, 0)),
            "body": Image.new("RGBA", (4, 4), (0, 0, 0, 0)),
            "head": Image.new("RGBA", (4, 4), (0, 0, 0, 0)),
        }
        layers["leg"].putpixel((1, 3), (20, 180, 20, 255))
        layers["body"].putpixel((1, 2), (20, 20, 180, 255))
        layers["head"].putpixel((1, 1), (180, 20, 20, 255))
        with tempfile.TemporaryDirectory() as temporary:
            output = Path(temporary) / "iruka-contact.png"
            render_iruka_composite_contact_sheet(layers, output)
            with Image.open(output) as contact:
                contact.load()
                self.assertEqual(contact.size, (8, 4))
                for offset, background in (
                    (0, (245, 245, 245, 255)),
                    (4, (24, 24, 24, 255)),
                ):
                    self.assertEqual(
                        contact.getpixel((offset, 0)),
                        background,
                    )
                    self.assertEqual(
                        contact.getpixel((offset + 1, 1)),
                        (180, 20, 20, 255),
                    )
                    self.assertEqual(
                        contact.getpixel((offset + 1, 2)),
                        (20, 20, 180, 255),
                    )
                    self.assertEqual(
                        contact.getpixel((offset + 1, 3)),
                        (20, 180, 20, 255),
                    )
