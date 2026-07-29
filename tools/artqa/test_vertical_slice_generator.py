import hashlib
import json
from copy import deepcopy
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

from PIL import Image


ROOT = Path(__file__).parents[2]
MAPPING = ROOT / "docs/my-thuat/mapping-id.json"
REVIEW = ROOT / "docs/my-thuat/review/vertical-slice-v1.json"
SOURCE_ROOT = Path(r"D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1")


class VerticalSliceGeneratorTests(unittest.TestCase):
    def _single_sprite_mapping(self, sprite_id=27078):
        mapping = json.loads(MAPPING.read_text(encoding="utf-8"))
        sprite = next(item for item in mapping["sprites"] if item["id"] == sprite_id)
        return {"sprites": [deepcopy(sprite)]}

    def _write_mapping(self, root, mapping):
        path = root / "mapping.json"
        path.write_text(json.dumps(mapping), encoding="utf-8")
        return path

    def test_generates_real_four_zoom_outputs_and_provenance_for_all_reviewed_assets(self):
        from tools.artqa.generate_vertical_slice import generate_reviewed_assets

        with tempfile.TemporaryDirectory() as temporary:
            output = Path(temporary) / "generated"
            index = generate_reviewed_assets(REVIEW, SOURCE_ROOT, output)

            self.assertEqual(index["asset_count"], 20)
            self.assertEqual(index["review_id"], "vertical-slice-v1")
            self.assertTrue((output / "generated-index.json").is_file())
            self.assertIn("iruka_portrait_chroma.png", index["sources"])
            for asset, record in index["assets"].items():
                self.assertEqual(set(record["zooms"]), {"x1", "x2", "x3", "x4"})
                for zoom, item in record["zooms"].items():
                    path = output / item["path"]
                    self.assertTrue(path.is_file(), f"{asset} {zoom}")
                    self.assertEqual(item["sha256"], hashlib.sha256(path.read_bytes()).hexdigest())
                    with Image.open(path) as generated:
                        self.assertEqual(item["dimensions"], list(generated.size))
                        self.assertEqual(generated.mode, "RGBA")

    def test_source_mutation_is_rejected_before_any_generated_output_is_created(self):
        from tools.artqa.generate_vertical_slice import GenerationError, generate_reviewed_assets

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            source_root = root / "source"
            source_root.mkdir()
            for source in SOURCE_ROOT.glob("*.png"):
                (source_root / source.name).write_bytes(source.read_bytes())
            (source_root / "qa").mkdir()
            for evidence_dir in ("alpha", "crop-reports", "contact-sheets", "previews"):
                target = source_root / "qa" / evidence_dir
                target.mkdir()
                for evidence in (SOURCE_ROOT / "qa" / evidence_dir).glob("*"):
                    (target / evidence.name).write_bytes(evidence.read_bytes())

            mapping = json.loads(MAPPING.read_text(encoding="utf-8"))
            review = json.loads(REVIEW.read_text(encoding="utf-8"))
            review["source_root"] = str(source_root)
            for record in review["assets"]:
                for evidence in record["evidence"].values():
                    evidence["path"] = str(source_root / "qa" / Path(evidence["path"]).relative_to(SOURCE_ROOT / "qa"))
            mapping_path = root / "mapping.json"
            review_path = root / "review.json"
            mapping_path.write_text(json.dumps(mapping), encoding="utf-8")
            review_path.write_text(json.dumps(review), encoding="utf-8")
            changed = source_root / "iruka_portrait_chroma.png"
            changed.write_bytes(changed.read_bytes() + b"mutation")
            output = root / "generated"

            with self.assertRaisesRegex(GenerationError, "ASSET-ART-001"):
                generate_reviewed_assets(review_path, source_root, output)
            self.assertFalse(output.exists())

    def test_generator_rejects_a_review_without_the_iruka_portrait(self):
        from tools.artqa.generate_vertical_slice import GenerationError, generate_vertical_slice

        with tempfile.TemporaryDirectory() as temporary:
            review = json.loads(REVIEW.read_text(encoding="utf-8"))
            review["assets"] = [item for item in review["assets"] if item["asset"] != "iruka_portrait_chroma.png"]
            review_path = Path(temporary) / "review.json"
            review_path.write_text(json.dumps(review), encoding="utf-8")
            with self.assertRaisesRegex(GenerationError, "ASSET-ART-002"):
                generate_vertical_slice(MAPPING, review_path, SOURCE_ROOT, Path(temporary) / "generated")

    def test_mapping_cannot_be_generated_from_unproven_sheet_coordinates(self):
        from tools.artqa.generate_vertical_slice import GenerationError, generate_vertical_slice

        with tempfile.TemporaryDirectory() as temporary:
            with self.assertRaisesRegex(GenerationError, "ASSET-ART-004.*27041"):
                generate_vertical_slice(MAPPING, REVIEW, SOURCE_ROOT, Path(temporary) / "generated")

    def test_missing_sprite_crop_27078_is_not_replaced_by_a_content_box(self):
        from tools.artqa.generate_vertical_slice import GenerationError, generate_vertical_slice

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            mapping = self._single_sprite_mapping()
            mapping_path = self._write_mapping(root, mapping)
            with self.assertRaisesRegex(GenerationError, "ASSET-ART-004.*27078.*explicit crop"):
                generate_vertical_slice(mapping_path, REVIEW, SOURCE_ROOT, root / "generated")

    def test_negative_source_crop_is_rejected_before_coordinate_conversion(self):
        from tools.artqa.generate_vertical_slice import GenerationError, generate_vertical_slice

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            mapping = self._single_sprite_mapping()
            mapping["sprites"][0]["crop"] = {"source_rect": [-1, 0, 2, 2]}
            mapping_path = self._write_mapping(root, mapping)
            with self.assertRaisesRegex(GenerationError, "ASSET-ART-004.*27078.*negative crop"):
                generate_vertical_slice(mapping_path, REVIEW, SOURCE_ROOT, root / "generated")

    def test_source_crop_outside_source_canvas_is_rejected_before_alpha_conversion(self):
        from tools.artqa.generate_vertical_slice import GenerationError, generate_vertical_slice

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            mapping = self._single_sprite_mapping()
            mapping["sprites"][0]["crop"] = {"source_rect": [1253, 0, 2, 2]}
            mapping_path = self._write_mapping(root, mapping)
            with self.assertRaisesRegex(GenerationError, "ASSET-ART-004.*27078.*source canvas"):
                generate_vertical_slice(mapping_path, REVIEW, SOURCE_ROOT, root / "generated")

    def test_source_crop_must_fit_derived_alpha_after_explicit_coordinate_conversion(self):
        from tools.artqa.generate_vertical_slice import GenerationError, generate_vertical_slice

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            mapping = self._single_sprite_mapping()
            mapping["sprites"][0]["crop"] = {"source_rect": [0, 0, 1, 1]}
            mapping_path = self._write_mapping(root, mapping)
            with self.assertRaisesRegex(GenerationError, "ASSET-ART-004.*27078.*alpha evidence"):
                generate_vertical_slice(mapping_path, REVIEW, SOURCE_ROOT, root / "generated")

    def test_cli_success_returns_sprite_count_for_a_proven_crop(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            mapping = self._single_sprite_mapping()
            mapping["sprites"][0]["crop"] = {"source_rect": [105, 58, 940, 1165]}
            mapping_path = self._write_mapping(root, mapping)
            output = root / "generated"
            result = subprocess.run(
                [
                    sys.executable,
                    str(ROOT / "tools/artqa/generate_vertical_slice.py"),
                    "--mapping", str(mapping_path),
                    "--review", str(REVIEW),
                    "--source-root", str(SOURCE_ROOT),
                    "--output", str(output),
                ],
                cwd=ROOT,
                text=True,
                capture_output=True,
            )
            self.assertEqual(result.returncode, 0, result.stderr)
            self.assertEqual(json.loads(result.stdout)["sprite_count"], 1)
            self.assertTrue((output / "generated-index.json").is_file())


if __name__ == "__main__":
    unittest.main()
