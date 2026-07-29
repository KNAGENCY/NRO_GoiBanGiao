import hashlib
import json
import tempfile
import unittest
from copy import deepcopy
from pathlib import Path


ROOT = Path(__file__).parents[2]
MAPPING = ROOT / "docs/my-thuat/mapping-id.json"
REVIEW = ROOT / "docs/my-thuat/review/vertical-slice-v1.json"
SHARED_ROOT = ROOT.parents[1]


class CanonicalTask6BPreflightTests(unittest.TestCase):
    def test_fixture_copy_is_explicitly_not_production_approved(self):
        provenance = json.loads(
            (ROOT / "docs/my-thuat/production-inputs/legacy/mob-fixture-copy.provenance.json")
            .read_text(encoding="utf-8")
        )
        self.assertEqual(provenance["source_role"], "committed_test_fixture")
        self.assertFalse(provenance["production_approved"])
        self.assertEqual(provenance["source_sha256"], provenance["independent_disk_match"]["sha256"])

    def test_production_manifest_declares_complete_canonical_closure(self):
        from assetpack.build import REQUIRED_TOP_LEVEL

        manifest = json.loads(MAPPING.read_text(encoding="utf-8"))
        self.assertEqual(set(manifest), REQUIRED_TOP_LEVEL)
        self.assertEqual(manifest["build_inputs"]["source_base"], ".")
        self.assertEqual([sprite["id"] for sprite in manifest["sprites"]], list(range(27041, 27083)))
        self.assertEqual(len(manifest["skills"]), 4)
        self.assertEqual(len(manifest["ui_bindings"]), 8)
        self.assertEqual(len(manifest["build_inputs"]["tile_art"]), 100)

        paths = []
        for sprite in manifest["sprites"]:
            paths.extend(sprite["zooms"].values())
        for tile in manifest["build_inputs"]["tile_art"]:
            paths.extend(tile["sources"].values())
        for sources in manifest["build_inputs"]["background_sources"].values():
            paths.extend(sources.values())
        for sources in manifest["build_inputs"]["hud_sources"].values():
            paths.extend(sources.values())
        self.assertTrue(paths)
        self.assertFalse(any("planned" in Path(path).parts for path in paths))

    def test_manifest_pins_every_generated_zoom_and_legacy_input_sha(self):
        manifest = json.loads(MAPPING.read_text(encoding="utf-8"))
        for sprite in manifest["sprites"]:
            self.assertEqual(set(sprite["zoom_sha256"]), {"x1", "x2", "x3", "x4"})
            self.assertTrue(all(len(value) == 64 for value in sprite["zoom_sha256"].values()))
        inputs = manifest["build_inputs"]
        self.assertEqual(len(inputs["tile_set_info_sha256"]), 64)
        self.assertEqual(len(inputs["legacy_part_sha256"]), 64)
        self.assertEqual(len(inputs["runtime_key_map_sha256"]), 64)
        self.assertEqual(len(inputs["tile_art_sha256"]), 100)

    def test_build_rejects_19_of_20_reviewed_source_provenance_records(self):
        from assetpack.build import (
            BuildError,
            _review_records,
            _validate_production_review_provenance,
            build,
        )

        manifest = deepcopy(json.loads(MAPPING.read_text(encoding="utf-8")))
        removed = manifest["build_inputs"]["source_provenance"].pop(0)
        del manifest["build_inputs"]["production_closure_sha256"][
            removed["source_path"]
        ]
        entries, _ = _review_records(REVIEW)
        with self.assertRaisesRegex(BuildError, "ASSET-PROVENANCE-014"):
            _validate_production_review_provenance(
                REVIEW, entries, manifest, SHARED_ROOT
            )

        with tempfile.TemporaryDirectory() as temporary:
            out = Path(temporary) / "out"
            with self.assertRaisesRegex(BuildError, "ASSET-PROVENANCE-014"):
                build(manifest, out, REVIEW, base_dir=SHARED_ROOT)
            self.assertFalse(out.exists())

    def test_build_rejects_duplicate_extra_malformed_or_wrong_sha_provenance(self):
        from assetpack.build import BuildError, build

        def duplicate(manifest):
            manifest["build_inputs"]["source_provenance"].append(
                deepcopy(manifest["build_inputs"]["source_provenance"][0])
            )

        def extra(manifest):
            relative = (
                "docs/my-thuat/anh-nguon/vertical-slice-v1/qa/alpha/"
                "hud_hp_chakra_4_alpha.png"
            )
            digest = hashlib.sha256((SHARED_ROOT / relative).read_bytes()).hexdigest()
            manifest["build_inputs"]["source_provenance"].append(
                {
                    "source_path": relative,
                    "source_sha256": digest,
                    "role": "reviewed_vertical_slice_source",
                }
            )
            manifest["build_inputs"]["production_closure_sha256"][relative] = digest

        def malformed(manifest):
            del manifest["build_inputs"]["source_provenance"][0]["source_sha256"]

        def wrong_sha(manifest):
            manifest["build_inputs"]["source_provenance"][0][
                "source_sha256"
            ] = "0" * 64

        for name, mutation in (
            ("duplicate", duplicate),
            ("extra", extra),
            ("malformed", malformed),
            ("wrong_sha", wrong_sha),
        ):
            with self.subTest(name=name), tempfile.TemporaryDirectory() as temporary:
                manifest = deepcopy(json.loads(MAPPING.read_text(encoding="utf-8")))
                mutation(manifest)
                out = Path(temporary) / "out"
                with self.assertRaisesRegex(BuildError, "ASSET-PROVENANCE-014"):
                    build(manifest, out, REVIEW, base_dir=SHARED_ROOT)
                self.assertFalse(out.exists())

    def test_build_rejects_missing_production_artifact_review_before_output(self):
        from assetpack.build import BuildError, build

        manifest = json.loads(MAPPING.read_text(encoding="utf-8"))
        review_payload = json.loads(REVIEW.read_text(encoding="utf-8"))
        review_payload.pop("production_artifact_review", None)
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            review = root / "review.json"
            review.write_text(
                json.dumps(review_payload, ensure_ascii=False),
                encoding="utf-8",
            )
            out = root / "out"
            with self.assertRaisesRegex(BuildError, "ASSET-ARTIFACT-REVIEW-016"):
                build(manifest, out, review, base_dir=SHARED_ROOT)
            self.assertFalse(out.exists())


if __name__ == "__main__":
    unittest.main()
