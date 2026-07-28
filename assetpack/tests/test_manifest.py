import json
import hashlib
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]


class ManifestValidationTests(unittest.TestCase):
    def minimal_manifest(self):
        return {
            "schema_version": 1,
            "versions": {"runtime_slice": 2},
            "sprites": [
                {
                    "id": 27041,
                    "zooms": {"x1": "a", "x2": "b", "x3": "c", "x4": "d"},
                    "source_sha256": "a" * 64,
                }
            ],
            "parts": [{"id": 30, "frames": [{"sprite_id": 27041, "dx": 0, "dy": 0}]}],
            "animations": [{"name": "idle", "frames": [{"slot_index": 0, "client_cf": [0]}]}],
        }

    def validate(self, manifest):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "manifest.json"
            path.write_text(json.dumps(manifest), encoding="utf-8")
            return subprocess.run(
                [sys.executable, "-m", "assetpack", "validate", "--manifest", str(path)],
                cwd=ROOT,
                text=True,
                capture_output=True,
                check=False,
            )

    def verify_legacy(self, ledger):
        return subprocess.run(
            [sys.executable, "-m", "assetpack", "verify-legacy-inputs", "--ledger", str(ledger)],
            cwd=ROOT,
            text=True,
            capture_output=True,
            check=False,
        )

    def test_rejects_duplicate_sprite_id(self):
        # Mutation caught: removing the duplicate-ID guard accepts two assets
        # with the same runtime identity.
        sprite = self.minimal_manifest()["sprites"][0]
        result = self.validate({"sprites": [sprite, dict(sprite)]})

        self.assertNotEqual(0, result.returncode)
        self.assertIn("ASSET-ID-001", result.stderr)

    def test_rejects_sprite_without_x3(self):
        # Mutation caught: accepting an incomplete four-zoom runtime closure.
        manifest = self.minimal_manifest()
        del manifest["sprites"][0]["zooms"]["x3"]

        result = self.validate(manifest)

        self.assertNotEqual(0, result.returncode)
        self.assertIn("ASSET-ZOOM-003", result.stderr)

    def test_rejects_part_pointing_to_unknown_sprite(self):
        # Mutation caught: allowing part.bin to reference an absent sprite.
        manifest = self.minimal_manifest()
        manifest["parts"][0]["frames"][0]["sprite_id"] = 27042

        result = self.validate(manifest)

        self.assertNotEqual(0, result.returncode)
        self.assertIn("ASSET-PART-004", result.stderr)

    def test_rejects_part_offset_outside_signed_byte(self):
        # Mutation caught: serializing an offset that cannot fit legacy sbyte.
        manifest = self.minimal_manifest()
        manifest["parts"][0]["frames"][0]["dx"] = 128

        result = self.validate(manifest)

        self.assertNotEqual(0, result.returncode)
        self.assertIn("ASSET-PART-005", result.stderr)

    def test_rejects_invalid_source_hash(self):
        # Mutation caught: accepting an unpinned source input.
        manifest = self.minimal_manifest()
        manifest["sprites"][0]["source_sha256"] = "not-a-sha"

        result = self.validate(manifest)

        self.assertNotEqual(0, result.returncode)
        self.assertIn("ASSET-SPRITE-002", result.stderr)

    def test_rejects_runtime_version_that_does_not_increase(self):
        # Mutation caught: shipping a changed runtime closure with stale version.
        manifest = self.minimal_manifest()
        manifest["previous_runtime_slice_version"] = 2

        result = self.validate(manifest)

        self.assertNotEqual(0, result.returncode)
        self.assertIn("ASSET-VERSION-006", result.stderr)

    def test_rejects_iruka_crop_without_overlap_and_layer(self):
        # Mutation caught: accepting a split NPC that can leave seams at joins.
        manifest = self.minimal_manifest()
        manifest["sprites"][0]["crop"] = {"subject": "iruka", "source_rect": [0, 0, 1, 1]}

        result = self.validate(manifest)

        self.assertNotEqual(0, result.returncode)
        self.assertIn("ASSET-CROP-007", result.stderr)

    def test_rejects_animation_frame_without_slot_and_cf(self):
        # Mutation caught: allowing an animation that cannot address CharInfo.
        manifest = self.minimal_manifest()
        manifest["animations"][0]["frames"] = [{}]

        result = self.validate(manifest)

        self.assertNotEqual(0, result.returncode)
        self.assertIn("ASSET-ANIMATION-008", result.stderr)

    def test_verifies_versioned_legacy_input_ledger(self):
        # Mutation caught: silently accepting a different baseline part.bin.
        ledger = ROOT / "assetpack" / "tests" / "fixtures" / "legacy-baseline.json"

        result = self.verify_legacy(ledger)

        self.assertEqual(0, result.returncode, result.stderr)

    def test_rejects_legacy_ledger_with_wrong_part_record_baseline(self):
        # Mutation caught: treating only the part.bin hash as its full baseline.
        fixture = ROOT / "assetpack" / "tests" / "fixtures" / "legacy-baseline.json"
        ledger = json.loads(fixture.read_text(encoding="utf-8"))
        ledger["expected"]["part_records"] = 2188
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "ledger.json"
            path.write_text(json.dumps(ledger), encoding="utf-8")
            result = self.verify_legacy(path)

        self.assertNotEqual(0, result.returncode)
        self.assertIn("ASSET-LEGACY-009", result.stderr)

    def test_accepts_the_committed_minimal_manifest_fixture(self):
        # Mutation caught: a schema change that makes the documented fixture stale.
        fixture = ROOT / "assetpack" / "tests" / "fixtures" / "manifest-minimal.json"
        result = subprocess.run(
            [sys.executable, "-m", "assetpack", "validate", "--manifest", str(fixture)],
            cwd=ROOT,
            text=True,
            capture_output=True,
            check=False,
        )

        self.assertEqual(0, result.returncode, result.stderr)

    def test_rejects_legacy_occupied_id_and_signed_short_overflow(self):
        # Mutation caught: allocating into a legacy identity or beyond int16.
        manifest = self.minimal_manifest()
        manifest["sprites"][0]["id"] = 27042
        manifest["parts"][0]["frames"][0]["sprite_id"] = 27042
        manifest["legacy_occupied_ids"] = [27042]

        collision = self.validate(manifest)
        self.assertIn("ASSET-ID-001", collision.stderr)

        manifest["sprites"][0]["id"] = 32768
        overflow = self.validate(manifest)
        self.assertIn("ASSET-ID-001", overflow.stderr)

    def test_rejects_poc_sprite_identity_from_official_manifest(self):
        # Mutation caught: accidentally carrying 9000..9002 from the POC forward.
        manifest = self.minimal_manifest()
        manifest["sprites"][0]["id"] = 9000
        manifest["parts"][0]["frames"][0]["sprite_id"] = 9000
        manifest["phase"] = {"forbidden_poc_ids": [9000, 9001, 9002]}

        result = self.validate(manifest)

        self.assertNotEqual(0, result.returncode)
        self.assertIn("ASSET-ID-001", result.stderr)


if __name__ == "__main__":
    unittest.main()
