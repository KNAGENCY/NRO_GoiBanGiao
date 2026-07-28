import json
import hashlib
import subprocess
import struct
import sys
import tempfile
import unittest
import zlib
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]


def rgba_png(width=1, height=1):
    def chunk(kind, payload):
        return (
            struct.pack(">I", len(payload))
            + kind
            + payload
            + struct.pack(">I", zlib.crc32(kind + payload))
        )

    rows = b"".join(b"\0" + b"\0\0\0\0" * width for _ in range(height))
    return (
        b"\x89PNG\r\n\x1a\n"
        + chunk(b"IHDR", struct.pack(">IIBBBBB", width, height, 8, 6, 0, 0, 0))
        + chunk(b"IDAT", zlib.compress(rows))
        + chunk(b"IEND", b"")
    )


class ManifestValidationTests(unittest.TestCase):
    def minimal_manifest(self):
        digest = hashlib.sha256(b"runtime-source").hexdigest()
        return {
            "schema_version": 1,
            "versions": {"runtime_slice": 2},
            "sprites": [
                {
                    "id": 27041,
                    "zooms": {"x1": "a", "x2": "b", "x3": "c", "x4": "x4.png"},
                    "source_path": "source.png",
                    "source_sha256": digest,
                    "x4_sha256": digest,
                    "source_dimensions": [1, 1],
                    "source_alpha": False,
                }
            ],
            "parts": [{"id": 30, "frames": [{"sprite_id": 27041, "dx": 0, "dy": 0}]}],
            "animations": [{"name": "idle", "frames": [{"slot_index": 0, "client_cf": [0]}]}],
        }

    def validate(
        self, manifest, *, legacy_ledger=None, version_ledger=None, files=None, art_root=None
    ):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "manifest.json"
            if files is None:
                files = self.runtime_files()
            for relative, content in files.items():
                target = Path(directory) / relative
                target.parent.mkdir(parents=True, exist_ok=True)
                target.write_bytes(content)
            path.write_text(json.dumps(manifest), encoding="utf-8")
            command = [sys.executable, "-m", "assetpack", "validate", "--manifest", str(path)]
            if legacy_ledger is not None:
                ledger_path = Path(directory) / "legacy.json"
                ledger_path.write_text(json.dumps(legacy_ledger), encoding="utf-8")
                command.extend(["--legacy-ledger", str(ledger_path)])
            if version_ledger is not None:
                version_path = Path(directory) / "versions.json"
                version_path.write_text(json.dumps(version_ledger), encoding="utf-8")
                command.extend(["--version-ledger", str(version_path)])
            if art_root is not None:
                command.extend(["--art-root", str(art_root)])
            return subprocess.run(
                command,
                cwd=ROOT,
                text=True,
                capture_output=True,
                check=False,
            )

    def runtime_files(self):
        return {"source.png": b"runtime-source", "x4.png": b"runtime-source"}

    def pin_runtime_source(self, manifest):
        digest = hashlib.sha256(b"runtime-source").hexdigest()
        manifest["sprites"][0].update(
            source_path="source.png",
            source_sha256=digest,
            x4_sha256=digest,
            source_dimensions=[1, 1],
            source_alpha=False,
        )
        manifest["sprites"][0]["zooms"]["x4"] = "x4.png"
        return manifest

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
        collision = self.validate(
            manifest,
            legacy_ledger={"inputs": [], "expected": {"occupied_ids": [27042]}},
        )
        self.assertNotEqual(0, collision.returncode)
        self.assertIn("ASSET-ID-001", collision.stderr)

        manifest["sprites"][0]["id"] = 32768
        overflow = self.validate(manifest)
        self.assertNotEqual(0, overflow.returncode)
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

    def test_legacy_ledger_requires_every_input_closure_and_relative_paths(self):
        required = {
            "part.bin",
            "NR_part",
            "NR_image",
            "legacy_skill",
            "legacy_mob",
            "small_local",
            "client_hardcoded_ids",
        }
        ledger = json.loads(
            (ROOT / "assetpack/tests/fixtures/legacy-baseline.json").read_text(encoding="utf-8")
        )
        self.assertEqual(required, {entry["kind"] for entry in ledger["inputs"]})
        self.assertTrue(all(not Path(entry["path"]).is_absolute() for entry in ledger["inputs"]))

    def test_verify_legacy_fails_when_each_required_input_is_missing_or_changed(self):
        fixture = ROOT / "assetpack/tests/fixtures/legacy-baseline.json"
        ledger = json.loads(fixture.read_text(encoding="utf-8"))
        for entry in ledger["inputs"]:
            with self.subTest(kind=entry["kind"]):
                changed = json.loads(json.dumps(ledger))
                target = next(item for item in changed["inputs"] if item["kind"] == entry["kind"])
                target["sha256"] = "f" * 64
                with tempfile.TemporaryDirectory() as directory:
                    path = Path(directory) / "ledger.json"
                    path.write_text(json.dumps(changed), encoding="utf-8")
                    result = self.verify_legacy(path)
                self.assertNotEqual(0, result.returncode)
                self.assertIn("ASSET-LEGACY-009", result.stderr)

    def test_verify_legacy_rejects_incomplete_required_closure(self):
        fixture = ROOT / "assetpack/tests/fixtures/legacy-baseline.json"
        ledger = json.loads(fixture.read_text(encoding="utf-8"))
        ledger["inputs"].pop()
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "ledger.json"
            path.write_text(json.dumps(ledger), encoding="utf-8")
            result = self.verify_legacy(path)
        self.assertNotEqual(0, result.returncode)
        self.assertIn("ASSET-LEGACY-009", result.stderr)

    def test_rejects_collision_from_independent_legacy_ledger(self):
        manifest = self.pin_runtime_source(self.minimal_manifest())
        ledger = {"inputs": [], "expected": {"occupied_ids": [27041]}}
        result = self.validate(manifest, legacy_ledger=ledger, files=self.runtime_files())
        self.assertNotEqual(0, result.returncode)
        self.assertIn("ASSET-ID-001", result.stderr)

    def test_rejects_sprite_outside_phase_one_allocation(self):
        manifest = self.pin_runtime_source(self.minimal_manifest())
        manifest["sprites"][0]["id"] = 9000
        manifest["parts"][0]["frames"][0]["sprite_id"] = 9000
        result = self.validate(manifest, files=self.runtime_files())
        self.assertNotEqual(0, result.returncode)
        self.assertIn("ASSET-ID-001", result.stderr)

    def test_rejects_missing_placeholder_and_mismatched_runtime_source(self):
        for mutation in ("missing_path", "zero_hash", "source_mismatch", "x4_mismatch"):
            with self.subTest(mutation=mutation):
                manifest = self.pin_runtime_source(self.minimal_manifest())
                files = self.runtime_files()
                if mutation == "missing_path":
                    del manifest["sprites"][0]["source_path"]
                elif mutation == "zero_hash":
                    manifest["sprites"][0]["source_sha256"] = "0" * 64
                elif mutation == "source_mismatch":
                    files["source.png"] = b"changed"
                else:
                    files["x4.png"] = b"different"
                result = self.validate(manifest, files=files)
                self.assertNotEqual(0, result.returncode)
                self.assertIn("ASSET-SPRITE-002", result.stderr)

    def test_rejects_runtime_version_not_above_independent_ledger(self):
        manifest = self.pin_runtime_source(self.minimal_manifest())
        result = self.validate(
            manifest,
            version_ledger={"runtime_slice_version": 2},
            files=self.runtime_files(),
        )
        self.assertNotEqual(0, result.returncode)
        self.assertIn("ASSET-VERSION-006", result.stderr)

    def test_rejects_invalid_iruka_layer_contract(self):
        invalid_crops = [
            [
                {"subject": "iruka", "source_rect": [0, 0, 1, 1], "overlap": [0, 0], "layer": 0},
                {"subject": "iruka", "source_rect": [0, 0, 1, 1], "overlap": [0, 0], "layer": 0},
            ],
            [{"subject": "iruka", "source_rect": [0, 0, 0, 1], "overlap": [0, 0], "layer": 0}],
            [{"subject": "iruka", "source_rect": [0, 0, 1, 1], "overlap": [128, 0], "layer": "0"}],
        ]
        for crops in invalid_crops:
            with self.subTest(crops=crops):
                manifest = self.pin_runtime_source(self.minimal_manifest())
                base = manifest["sprites"][0]
                manifest["sprites"] = []
                for offset, crop in enumerate(crops):
                    sprite = json.loads(json.dumps(base))
                    sprite["id"] = 27041 + offset
                    sprite["crop"] = crop
                    manifest["sprites"].append(sprite)
                manifest["parts"] = []
                result = self.validate(manifest, files=self.runtime_files())
                self.assertNotEqual(0, result.returncode)
                self.assertIn("ASSET-CROP-007", result.stderr)

    def test_rejects_animation_slot_range_and_client_cf_integrity(self):
        invalid_frames = [
            [{"slot_index": -1, "client_cf": [0]}],
            [{"slot_index": 1, "client_cf": [0]}],
            [{"slot_index": 0, "client_cf": [0, 0]}],
            [{"slot_index": 0, "client_cf": [128]}],
            [{"slot_index": "0", "client_cf": [0]}],
        ]
        for frames in invalid_frames:
            with self.subTest(frames=frames):
                manifest = self.pin_runtime_source(self.minimal_manifest())
                manifest["animations"][0]["slot_count"] = 1
                manifest["animations"][0]["frames"] = frames
                result = self.validate(manifest, files=self.runtime_files())
                self.assertNotEqual(0, result.returncode)
                self.assertIn("ASSET-ANIMATION-008", result.stderr)

    def test_rejects_common_runtime_placeholder_for_all_sprites(self):
        # Mutation caught: replacing every real source with one tiny sentinel.
        manifest = self.pin_runtime_source(self.minimal_manifest())
        duplicate = json.loads(json.dumps(manifest["sprites"][0]))
        duplicate["id"] = 27042
        manifest["sprites"].append(duplicate)
        manifest["parts"] = []

        result = self.validate(manifest)

        self.assertNotEqual(0, result.returncode)
        self.assertIn("ASSET-SPRITE-002", result.stderr)

    def test_verifies_source_and_x4_with_independent_paths_and_hashes(self):
        # Mutation caught: checking only source_path, or reusing source SHA for x4.
        manifest = self.pin_runtime_source(self.minimal_manifest())
        manifest["sprites"][0]["zooms"]["x4"] = "rendered-x4.png"
        manifest["sprites"][0]["x4_sha256"] = hashlib.sha256(b"rendered-x4").hexdigest()
        files = {"source.png": b"runtime-source", "rendered-x4.png": b"rendered-x4"}
        accepted = self.validate(manifest, files=files)
        self.assertEqual(0, accepted.returncode, accepted.stderr)

        for field, value in (
            ("source_sha256", "f" * 64),
            ("x4_sha256", "f" * 64),
        ):
            with self.subTest(field=field):
                changed = json.loads(json.dumps(manifest))
                changed["sprites"][0][field] = value
                result = self.validate(changed, files=files)
                self.assertNotEqual(0, result.returncode)
                self.assertIn("ASSET-SPRITE-002", result.stderr)

    def test_rejects_bool_in_typed_integer_fields(self):
        # Mutation caught: Python bool is an int subclass and must not enter binary fields.
        mutations = (
            ("source_rect", lambda manifest: manifest["sprites"][0].update(
                crop={"subject": "iruka", "source_rect": [False, 0, 1, 1],
                      "overlap": [0, 0], "layer": 0}
            ), "ASSET-CROP-007"),
            ("overlap", lambda manifest: manifest["sprites"][0].update(
                crop={"subject": "iruka", "source_rect": [0, 0, 1, 1],
                      "overlap": [False, 0], "layer": 0}
            ), "ASSET-CROP-007"),
            ("layer", lambda manifest: manifest["sprites"][0].update(
                crop={"subject": "iruka", "source_rect": [0, 0, 1, 1],
                      "overlap": [0, 0], "layer": False}
            ), "ASSET-CROP-007"),
            ("slot_index", lambda manifest: manifest["animations"][0]["frames"][0].update(
                slot_index=False
            ), "ASSET-ANIMATION-008"),
            ("client_cf", lambda manifest: manifest["animations"][0]["frames"][0].update(
                client_cf=[False]
            ), "ASSET-ANIMATION-008"),
        )
        for field, mutate, error_code in mutations:
            with self.subTest(field=field):
                manifest = self.pin_runtime_source(self.minimal_manifest())
                mutate(manifest)
                result = self.validate(manifest)
                self.assertNotEqual(0, result.returncode)
                self.assertIn(error_code, result.stderr)

    def test_verifies_declared_source_dimensions_and_alpha(self):
        # Mutation caught: a regenerated source with a changed canvas/alpha contract.
        manifest = self.minimal_manifest()
        png = rgba_png(2, 3)
        digest = hashlib.sha256(png).hexdigest()
        manifest["sprites"][0].update(
            source_path="source.png",
            source_sha256=digest,
            x4_sha256=digest,
            source_dimensions=[2, 3],
            source_alpha=True,
        )
        manifest["sprites"][0]["zooms"]["x4"] = "source.png"
        accepted = self.validate(manifest, files={"source.png": png})
        self.assertEqual(0, accepted.returncode, accepted.stderr)

        for field, value in (("source_dimensions", [3, 2]), ("source_alpha", False)):
            with self.subTest(field=field):
                changed = json.loads(json.dumps(manifest))
                changed["sprites"][0][field] = value
                result = self.validate(changed, files={"source.png": png})
                self.assertNotEqual(0, result.returncode)
                self.assertIn("ASSET-SPRITE-002", result.stderr)


if __name__ == "__main__":
    unittest.main()
