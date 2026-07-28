"""Kiểm thử hợp đồng QA và cổng Task 6A.

Các test dùng ảnh nhỏ tự tạo, không đụng ảnh nguồn dùng chung.
"""

from __future__ import annotations

import json
from pathlib import Path
import tempfile
import unittest

from PIL import Image

from tools.artqa.artqa import (
    analyze_source_background,
    extract_alpha,
    normalize_cells,
    render_animation_preview,
    validate_animation,
    verify_review_gate,
)


class AlphaExtractionTests(unittest.TestCase):
    def test_adaptive_gradient_magenta_extracts_but_raw_source_is_rejected(self):
        image = Image.new("RGB", (25, 7))
        for y in range(image.height):
            for x in range(image.width):
                image.putpixel((x, y), (220 + x, 15 + y, 225 + x))
        image.paste((30, 80, 140), (10, 2, 15, 6))

        source_issues, _ = analyze_source_background(image)
        result = extract_alpha(image)

        self.assertIn("background_gradient_exceeds_static_matte_tolerance", source_issues)
        self.assertEqual(result.getpixel((0, 0))[3], 0)
        self.assertEqual(result.getpixel((12, 3))[3], 255)

    def test_near_white_hud_is_not_misclassified_as_chroma(self):
        image = Image.new("RGB", (7, 7), (250, 248, 246))
        image.paste((170, 20, 30), (2, 2, 5, 5))

        source_issues, _ = analyze_source_background(image)
        result = extract_alpha(image)

        self.assertIn("border_chroma_not_detected", source_issues)
        self.assertEqual(result.getpixel((0, 0))[3], 255)


class CropAndAnimationTests(unittest.TestCase):
    def test_normalize_cells_crops_real_edges_without_cutting_frame(self):
        strip = Image.new("RGBA", (11, 6), (0, 0, 0, 0))
        strip.paste((255, 0, 0, 255), (1, 1, 4, 6))
        strip.paste((0, 255, 0, 255), (7, 2, 11, 6))

        normalized, report, issues = normalize_cells(strip, columns=2)

        self.assertEqual(issues, [])
        self.assertEqual(normalized.width % 2, 0)
        self.assertEqual(report["source_edges"], [0, 6, 11])
        self.assertEqual(report["content_boxes"][1], [1, 2, 5, 6])

    def test_loop_and_timing_reject_duplicate_or_bad_duration(self):
        frames = []
        for x in (0, 0, 6):
            frame = Image.new("RGBA", (10, 6), (0, 0, 0, 0))
            frame.paste((255, 255, 255, 255), (x, 2, x + 2, 6))
            frames.append(frame)

        issues, _ = validate_animation(frames, [100, 0, 100], loop=True)

        self.assertIn("invalid_frame_timing", issues)
        self.assertIn("duplicate_adjacent_frame", issues)
        self.assertIn("loop_boundary_motion_jump", issues)

    def test_animation_preview_is_temporal_gif(self):
        frames = [
            Image.new("RGBA", (4, 4), (255, 0, 0, 255)),
            Image.new("RGBA", (4, 4), (0, 255, 0, 255)),
        ]
        with tempfile.TemporaryDirectory() as tmp:
            output = Path(tmp) / "preview.gif"
            render_animation_preview(frames, output, [80, 120])
            with Image.open(output) as preview:
                self.assertEqual(getattr(preview, "n_frames", 1), 2)
                self.assertEqual(preview.info["duration"], 80)


class GateTests(unittest.TestCase):
    def _write_review(self, root: Path) -> Path:
        evidence = root / "alpha.png"
        evidence.write_bytes(b"alpha")
        digest = __import__("hashlib").sha256(evidence.read_bytes()).hexdigest()
        review = {
            "review_id": "vertical-slice-v1",
            "assets": [{
                "asset": "a.png",
                "required": True,
                "status": "PASS",
                "evidence": {
                    "alpha": {"path": str(evidence), "sha256": digest},
                    "crop_report": {"path": str(evidence), "sha256": digest},
                    "contact_sheet": {"path": str(evidence), "sha256": digest},
                    "animation_preview": {"path": str(evidence), "sha256": digest},
                },
            }],
            "summary": {"total": 1, "pass": 1, "fail": 0},
        }
        path = root / "review.json"
        path.write_text(json.dumps(review), encoding="utf-8")
        return path

    def test_gate_rejects_evidence_sha_mismatch(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            review = self._write_review(root)
            (root / "alpha.png").write_bytes(b"tampered")

            ok, issues = verify_review_gate(review)

            self.assertFalse(ok)
            self.assertTrue(any("sha256_mismatch" in issue for issue in issues))

    def test_gate_cannot_bypass_a_required_fail_via_summary(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            review = self._write_review(root)
            data = json.loads(review.read_text(encoding="utf-8"))
            data["assets"][0]["status"] = "FAIL"
            data["summary"] = {"total": 1, "pass": 1, "fail": 0}
            review.write_text(json.dumps(data), encoding="utf-8")

            ok, issues = verify_review_gate(review)

            self.assertFalse(ok)
            self.assertIn("required_asset_not_pass:a.png", issues)


if __name__ == "__main__":
    unittest.main()
