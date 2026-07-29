import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


ROOT = Path(__file__).parents[2]
SHARED_ROOT = Path(r"D:\NRO_GoiBanGiao")
MAPPING = ROOT / "docs/my-thuat/mapping-id.json"
REVIEW = ROOT / "docs/my-thuat/review/vertical-slice-v1.json"
REVIEWED_NAMES = {
    item["asset"]
    for item in json.loads(REVIEW.read_text(encoding="utf-8"))["assets"]
}


class Task6BRound2ProductionTests(unittest.TestCase):
    def test_canonical_source_base_resolves_repo_root_relative_artroot_paths(self):
        from assetpack.build import resolve_manifest_base

        manifest = {
            "build_inputs": {"source_base": "."},
            "sprites": [{
                "source_path": "docs/my-thuat/anh-nguon/vertical-slice-v1/iruka_portrait_chroma.png",
            }],
        }

        base_dir = resolve_manifest_base(ROOT / "docs/my-thuat/mapping-id.json", manifest, base_dir=SHARED_ROOT)

        self.assertEqual(base_dir, SHARED_ROOT.resolve())
        source = base_dir / manifest["sprites"][0]["source_path"]
        self.assertTrue(source.is_file())

    def test_canonical_source_base_rejects_mapping_relative_short_source(self):
        from assetpack.build import BuildError, resolve_manifest_base

        manifest = {
            "build_inputs": {"source_base": "."},
            "sprites": [{"source_path": "iruka_portrait_chroma.png"}],
        }

        with self.assertRaisesRegex(BuildError, "ASSET-REVIEW-010"):
            resolve_manifest_base(ROOT / "docs/my-thuat/mapping-id.json", manifest, base_dir=SHARED_ROOT)

    def test_bind_source_path_must_resolve_from_the_declared_build_base(self):
        from assetpack.bind import BindError, validate_reviewed_sources

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            source = root / "art" / "approved.png"
            source.parent.mkdir()
            source.write_bytes(b"approved source")
            manifest = {
                "sprites": [{
                    "source_path": "art/missing.png",
                    "source_sha256": "a" * 64,
                }]
            }
            review = {"assets": [{
                "asset": "missing.png",
                "source_sha256": "a" * 64,
            }]}
            with self.assertRaisesRegex(BindError, "ASSET-REVIEW-010"):
                validate_reviewed_sources(manifest, review, base_dir=root)

    def test_bind_accepts_matching_source_and_review_sha_from_the_same_base(self):
        import hashlib
        from assetpack.bind import validate_reviewed_sources

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            source = root / "art" / "approved.png"
            source.parent.mkdir()
            source.write_bytes(b"approved source")
            digest = hashlib.sha256(source.read_bytes()).hexdigest()
            manifest = {"sprites": [{"source_path": "art/approved.png", "source_sha256": digest}]}
            review = {"assets": [{"asset": "approved.png", "source_sha256": digest}]}
            validate_reviewed_sources(manifest, review, base_dir=root)

    def test_cli_base_dir_reaches_manifest_validation(self):
        with tempfile.TemporaryDirectory() as temporary:
            temporary_root = Path(temporary)
            manifest = temporary_root / "mapping.json"
            review = temporary_root / "review.json"
            manifest.write_text('{"slice": "lang-la"}', encoding="utf-8")
            review.write_text('{"assets": []}', encoding="utf-8")
            result = subprocess.run(
                [
                    sys.executable,
                    "-m",
                    "assetpack",
                    "build",
                    "--manifest",
                    str(manifest),
                    "--slice",
                    "lang-la",
                    "--base-dir",
                    str(temporary_root),
                    "--out",
                    str(temporary_root / "asset-build"),
                    "--review",
                    str(review),
                ],
                cwd=ROOT,
                text=True,
                capture_output=True,
            )
            self.assertEqual(result.returncode, 1)
            self.assertIn("ASSET-SCHEMA-001", result.stderr)

    def test_production_mapping_is_canonical_and_resolves_all_reviewed_sources(self):
        from assetpack.build import REQUIRED_TOP_LEVEL

        manifest = json.loads(MAPPING.read_text(encoding="utf-8"))
        self.assertEqual(set(manifest), REQUIRED_TOP_LEVEL)
        self.assertIn("source_base", manifest["build_inputs"])

        base_dir = (MAPPING.parent / manifest["build_inputs"]["source_base"]).resolve()
        sources = {Path(sprite["source_path"]).name for sprite in manifest["sprites"]}
        self.assertEqual(sources, REVIEWED_NAMES)
        self.assertIn("iruka_portrait_chroma.png", sources)
        for sprite in manifest["sprites"]:
            source_path = Path(sprite["source_path"])
            self.assertFalse(source_path.is_absolute())
            self.assertNotIn("tests", source_path.parts)
            self.assertTrue((base_dir / source_path).is_file(), sprite["source_path"])

    def test_cli_real_production_build_uses_manifest_base_and_creates_artifacts(self):
        manifest = json.loads(MAPPING.read_text(encoding="utf-8"))
        build_inputs = manifest.get("build_inputs", {})
        self.assertIsInstance(build_inputs, dict)
        source_base = build_inputs.get("source_base")
        self.assertIsInstance(source_base, str)
        base_dir = (MAPPING.parent / source_base).resolve()
        with tempfile.TemporaryDirectory() as temporary:
            out = Path(temporary) / "asset-build"
            result = subprocess.run(
                [
                    sys.executable,
                    "-m",
                    "assetpack",
                    "build",
                    "--manifest",
                    str(MAPPING),
                    "--slice",
                    "lang-la",
                    "--base-dir",
                    str(base_dir),
                    "--out",
                    str(out),
                    "--review",
                    str(REVIEW),
                ],
                cwd=ROOT,
                text=True,
                capture_output=True,
            )
            self.assertEqual(result.returncode, 0, result.stderr or result.stdout)
            self.assertTrue((out / "server").is_dir())
            self.assertTrue((out / "client").is_dir())
            self.assertTrue((out / "asset-lock.json").is_file())


if __name__ == "__main__":
    unittest.main()
