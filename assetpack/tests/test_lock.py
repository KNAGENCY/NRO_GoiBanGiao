import json
import hashlib
import tempfile
import unittest
from pathlib import Path


class LockAndBuildGateTests(unittest.TestCase):
    def test_review_fail_blocks_build_before_artifact_is_bound(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            (root / "review.json").write_text(json.dumps({"assets": [{"status": "FAIL"}]}), encoding="utf-8")
            with self.assertRaisesRegex(BuildError, "ASSET-REVIEW-010"):
                build({"slice": "lang-la"}, root / "out", root / "review.json")
            self.assertFalse((root / "out").exists())

    def test_lock_is_canonical_and_two_separate_builds_have_the_same_sha(self):
        from assetpack.lock import lock_bytes

        left = lock_bytes({"versions": {"data": 1}, "outputs": [{"path": "b", "sha256": "2"}, {"path": "a", "sha256": "1"}]})
        right = lock_bytes({"outputs": [{"sha256": "1", "path": "a"}, {"sha256": "2", "path": "b"}], "versions": {"data": 1}})
        self.assertEqual(left, right)

    def test_approved_build_writes_catalog_and_lock_with_identical_tree_hashes(self):
        from assetpack.build import build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            evidence = root / "evidence.png"; evidence.write_bytes(b"golden")
            source = root / "source.png"; source.write_bytes(b"source")
            review = {"assets": [{"status": "PASS", "source_sha256": hashlib.sha256(source.read_bytes()).hexdigest(), "evidence_path": str(evidence), "evidence_sha256": hashlib.sha256(evidence.read_bytes()).hexdigest()}]}
            review_path = root / "review.json"; review_path.write_text(json.dumps(review), encoding="utf-8")
            manifest = {"slice": "lang-la", "versions": {"data": 1}, "catalog": {"maps": ["Làng Lá"], "npcs": ["Iruka"], "mobs": ["Sói"], "spawn": [1], "skill_poses": [3]}, "review_sources": [str(source)]}
            first, second = root / "one", root / "two"
            build(manifest, first, review_path); build(manifest, second, review_path)
            def tree_hash(directory):
                digest = hashlib.sha256()
                for path in sorted(directory.rglob("*")):
                    if path.is_file(): digest.update(path.relative_to(directory).as_posix().encode()); digest.update(path.read_bytes())
                return digest.hexdigest()
            self.assertEqual(tree_hash(first), tree_hash(second))
            lock = json.loads((first / "asset-lock.json").read_text(encoding="utf-8"))
            self.assertTrue(lock["review_sha256"])
            self.assertTrue((first / "catalog.json").is_file())
