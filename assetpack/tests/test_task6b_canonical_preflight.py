import json
import unittest
from pathlib import Path


ROOT = Path(__file__).parents[2]
MAPPING = ROOT / "docs/my-thuat/mapping-id.json"


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


if __name__ == "__main__":
    unittest.main()
