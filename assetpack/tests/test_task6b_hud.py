import importlib
import importlib.util
import json
import tempfile
import unittest
from pathlib import Path

from PIL import Image


ROOT = Path(__file__).parents[2]
SHARED_ROOT = Path(r"D:\NRO_GoiBanGiao")
ART_ROOT = SHARED_ROOT / "docs/my-thuat/anh-nguon/vertical-slice-v1"
REVIEW = ROOT / "docs/my-thuat/review/vertical-slice-v1.json"
HUD_NAMES = (
    "naruto-hp-bg",
    "naruto-hp-fill",
    "naruto-chakra-bg",
    "naruto-chakra-fill",
)


def _load_closure_generator():
    try:
        return importlib.import_module("tools.artqa.task6b_closure")
    except ModuleNotFoundError:
        # Vòng fix 4 đã hoàn nguyên source khi visual gate đỏ, nhưng bytecode
        # của chính generator đó còn lại và là oracle RED trung thực trước fix.
        bytecode = ROOT / "tools/artqa/__pycache__/task6b_closure.cpython-312.pyc"
        spec = importlib.util.spec_from_file_location(
            "tools.artqa.task6b_closure_recovered",
            bytecode,
        )
        if spec is None or spec.loader is None:
            raise
        module = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(module)
        return module


class Task6BHudProductionTests(unittest.TestCase):
    def test_hud_x4_uses_the_bar_content_not_the_quadrant_dividers(self):
        generator = _load_closure_generator()
        review = json.loads(REVIEW.read_text(encoding="utf-8"))

        with tempfile.TemporaryDirectory() as temporary:
            output = Path(temporary)
            generator._generate_hud(ART_ROOT, review, output)
            images = {
                name: Image.open(output / "hud" / name / "x4/sprite.png").convert("RGBA")
                for name in HUD_NAMES
            }

            for name, image in images.items():
                self.assertEqual(image.size, (240, 32), name)
                bbox = image.getchannel("A").getbbox()
                self.assertIsNotNone(bbox, name)
                self.assertGreaterEqual(bbox[2] - bbox[0], 180, name)
                self.assertEqual(image.getchannel("A").getextrema(), (0, 255), name)

            for resource in ("hp", "chakra"):
                background = images[f"naruto-{resource}-bg"]
                fill = images[f"naruto-{resource}-fill"]
                center = (background.width // 2, background.height // 2)
                self.assertEqual(background.getchannel("A").getpixel(center), 0, resource)
                self.assertGreater(fill.getchannel("A").getpixel(center), 0, resource)
                self.assertNotEqual(background.tobytes(), fill.tobytes(), resource)


if __name__ == "__main__":
    unittest.main()
