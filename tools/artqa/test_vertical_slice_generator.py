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
    def test_alpha_rect_applies_normalize_center_and_bottom_offsets_for_player_idle_frame_1(self):
        from tools.artqa.generate_vertical_slice import _alpha_rect

        crop_report = {
            "source_size": [2172, 724],
            "source_edges": [0, 724, 1448, 2172],
            "content_boxes": [
                [180, 62, 520, 715],
                [160, 68, 495, 713],
                [124, 65, 461, 715],
            ],
            "output_cell_size": [340, 653],
        }
        rect = _alpha_rect(
            {"id": "player_idle", "frame": 1},
            [884, 68, 10, 10],
            crop_report,
            (1020, 653),
        )

        self.assertEqual(rect, [342, 8, 10, 10])

    def test_alpha_rect_rejects_malformed_source_edges_as_generation_error(self):
        from tools.artqa.generate_vertical_slice import GenerationError, _alpha_rect

        base = {
            "source_size": [2172, 724],
            "source_edges": [0, 724, 1448, 2172],
            "content_boxes": [[180, 62, 520, 715]] * 3,
            "output_cell_size": [340, 653],
        }
        for edges in (["0", 724, 1448, 2172], [-1, 724, 1448, 2172], [0, 724, 1448, 2173]):
            with self.subTest(edges=edges):
                crop_report = deepcopy(base)
                crop_report["source_edges"] = edges
                with self.assertRaisesRegex(GenerationError, "ASSET-ART-004"):
                    _alpha_rect({"id": "player_idle"}, [884, 68, 10, 10], crop_report, (1020, 653))

    def test_alpha_rect_rejects_malformed_output_cell_size_as_generation_error(self):
        from tools.artqa.generate_vertical_slice import GenerationError, _alpha_rect

        base = {
            "source_size": [2172, 724],
            "source_edges": [0, 724, 1448, 2172],
            "content_boxes": [[180, 62, 520, 715]] * 3,
            "output_cell_size": [340, 653],
        }
        for cell_size in (["340", 653], [0, 653], [1021, 653]):
            with self.subTest(cell_size=cell_size):
                crop_report = deepcopy(base)
                crop_report["output_cell_size"] = cell_size
                with self.assertRaisesRegex(GenerationError, "ASSET-ART-004"):
                    _alpha_rect({"id": "player_idle"}, [884, 68, 10, 10], crop_report, (1020, 653))

    def _single_sprite_mapping(self, sprite_id=27078):
        mapping = json.loads(MAPPING.read_text(encoding="utf-8"))
        sprite = next(item for item in mapping["sprites"] if item["id"] == sprite_id)
        return {"sprites": [deepcopy(sprite)]}

    def _write_mapping(self, root, mapping):
        path = root / "mapping.json"
        path.write_text(json.dumps(mapping), encoding="utf-8")
        return path

    def _mapping_for_ids(self, ids):
        mapping = json.loads(MAPPING.read_text(encoding="utf-8"))
        wanted = set(ids)
        mapping["sprites"] = [
            deepcopy(sprite) for sprite in mapping["sprites"] if sprite["id"] in wanted
        ]
        return mapping

    def test_raw_sheet_crop_precedes_alpha_for_all_hoc_vien_sprites(self):
        from tools.artqa.generate_vertical_slice import generate_vertical_slice

        ids = list(range(27041, 27075))
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            mapping = self._mapping_for_ids(ids)
            mapping_path = self._write_mapping(root, mapping)
            index = generate_vertical_slice(mapping_path, REVIEW, SOURCE_ROOT, root / "generated")

            self.assertEqual(index["sprite_count"], len(ids))
            self.assertEqual(set(index["sprites"]), {str(sprite_id) for sprite_id in ids})
            for sprite in mapping["sprites"]:
                record = index["sprites"][str(sprite["id"])]
                self.assertEqual(record["source_rect"], sprite["crop"]["source_rect"])
                raw_crop = root / "generated" / record["raw_crop"]["path"]
                alpha = root / "generated" / record["alpha"]["path"]
                self.assertTrue(raw_crop.is_file())
                self.assertTrue(alpha.is_file())
                with Image.open(raw_crop) as image:
                    self.assertEqual(image.size, tuple(sprite["crop"]["source_rect"][2:]))
                self.assertEqual(record["raw_crop"]["sha256"], hashlib.sha256(raw_crop.read_bytes()).hexdigest())
                self.assertEqual(record["alpha"]["sha256"], hashlib.sha256(alpha.read_bytes()).hexdigest())
                self.assertEqual(set(record["zooms"]), {"x1", "x2", "x3", "x4"})
                self.assertEqual(set(record["evidence"]), {"alpha", "crop", "zooms"})
                for zoom in record["zooms"].values():
                    self.assertTrue((root / "generated" / zoom["path"]).is_file())

    def test_shared_sheet_frame_uses_only_its_raw_crop_for_alpha_derivation(self):
        from tools.artqa.artqa import extract_alpha, normalize_cells
        from tools.artqa.generate_vertical_slice import generate_vertical_slice

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            mapping = self._mapping_for_ids([27041, 27042])
            mapping_path = self._write_mapping(root, mapping)
            index = generate_vertical_slice(mapping_path, REVIEW, SOURCE_ROOT, root / "generated")
            source = SOURCE_ROOT / "master_hoc_vien_x4_chroma.png"
            for sprite in mapping["sprites"]:
                rect = sprite["crop"]["source_rect"]
                with Image.open(source) as sheet:
                    raw_crop = sheet.crop((rect[0], rect[1], rect[0] + rect[2], rect[1] + rect[3])).convert("RGBA")
                expected, _, _ = normalize_cells(extract_alpha(raw_crop), 1)
                x4_path = root / "generated" / index["sprites"][str(sprite["id"])] ["zooms"]["x4"]["path"]
                with Image.open(x4_path) as actual:
                    self.assertEqual(actual.convert("RGBA").tobytes(), expected.tobytes())

    def test_missing_source_rect_fails_closed_for_sprite_27078(self):
        from tools.artqa.generate_vertical_slice import GenerationError, generate_vertical_slice

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            mapping = self._mapping_for_ids([27078])
            mapping["sprites"][0].pop("crop", None)
            mapping_path = self._write_mapping(root, mapping)
            with self.assertRaisesRegex(GenerationError, "ASSET-ART-004.*27078.*source_rect"):
                generate_vertical_slice(mapping_path, REVIEW, SOURCE_ROOT, root / "generated")

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

    def test_mapping_sheet_coordinates_are_proven_by_raw_crop(self):
        from tools.artqa.generate_vertical_slice import generate_vertical_slice

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            mapping_path = self._write_mapping(root, self._mapping_for_ids([27041]))
            index = generate_vertical_slice(mapping_path, REVIEW, SOURCE_ROOT, root / "generated")
            self.assertEqual(index["sprite_count"], 1)
            self.assertTrue((root / "generated" / index["sprites"]["27041"]["alpha"]["path"]).is_file())

    def test_missing_sprite_crop_27078_is_not_replaced_by_a_content_box(self):
        from tools.artqa.generate_vertical_slice import GenerationError, generate_vertical_slice

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            mapping = self._single_sprite_mapping()
            mapping["sprites"][0].pop("crop", None)
            mapping_path = self._write_mapping(root, mapping)
            with self.assertRaisesRegex(GenerationError, "ASSET-ART-004.*27078.*source_rect"):
                generate_vertical_slice(mapping_path, REVIEW, SOURCE_ROOT, root / "generated")

    def test_negative_source_crop_is_rejected_before_coordinate_conversion(self):
        from tools.artqa.generate_vertical_slice import GenerationError, generate_vertical_slice

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            mapping = self._single_sprite_mapping()
            mapping["sprites"][0]["crop"] = {"source_rect": [-1, 0, 2, 2]}
            mapping_path = self._write_mapping(root, mapping)
            with self.assertRaisesRegex(GenerationError, "ASSET-ART-004.*27078.*source_rect negative"):
                generate_vertical_slice(mapping_path, REVIEW, SOURCE_ROOT, root / "generated")

    def test_source_crop_outside_source_canvas_is_rejected_before_alpha_conversion(self):
        from tools.artqa.generate_vertical_slice import GenerationError, generate_vertical_slice

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            mapping = self._single_sprite_mapping()
            mapping["sprites"][0]["crop"] = {"source_rect": [1253, 0, 2, 2]}
            mapping_path = self._write_mapping(root, mapping)
            with self.assertRaisesRegex(GenerationError, "ASSET-ART-004.*27078.*raw source"):
                generate_vertical_slice(mapping_path, REVIEW, SOURCE_ROOT, root / "generated")

    def test_valid_empty_raw_crop_is_derived_without_alpha_coordinate_fallback(self):
        from tools.artqa.generate_vertical_slice import generate_vertical_slice

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            mapping = self._single_sprite_mapping()
            mapping["sprites"][0]["crop"] = {"source_rect": [0, 0, 1, 1]}
            mapping_path = self._write_mapping(root, mapping)
            index = generate_vertical_slice(mapping_path, REVIEW, SOURCE_ROOT, root / "generated")
            self.assertEqual(index["sprites"]["27078"]["source_rect"], [0, 0, 1, 1])

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
