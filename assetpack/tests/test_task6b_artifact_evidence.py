import hashlib
import json
import tempfile
import unittest
from copy import deepcopy
from pathlib import Path

from PIL import Image, ImageChops


ROOT = Path(__file__).parents[2]
SHARED_ROOT = Path(r"D:\NRO_GoiBanGiao")
ART_ROOT = SHARED_ROOT / "docs/my-thuat/anh-nguon/vertical-slice-v1"
REVIEW = ROOT / "docs/my-thuat/review/vertical-slice-v1.json"


class Task6BArtifactEvidenceTests(unittest.TestCase):
    def _artifact_review_fixture(self, root: Path):
        from assetpack.build import production_closure_digest

        evidence = {}
        closure = {}
        for kind, payload in (
            ("hud_golden", b"hud evidence"),
            ("iruka_composite", b"iruka evidence"),
        ):
            relative = f"evidence/{kind}.png"
            path = root / relative
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_bytes(payload)
            digest = hashlib.sha256(payload).hexdigest()
            evidence[kind] = {"path": relative, "sha256": digest}
            closure[relative] = digest
        manifest = {
            "build_inputs": {
                "artifact_evidence": deepcopy(evidence),
                "production_closure_sha256": closure,
            }
        }
        review = {
            "reviewer": "Source Reviewer Fixture",
            "production_artifact_review": {
                "status": "PASS",
                "reviewer": "/root/task6b_visual_review6",
                "timestamp": "2026-07-29T19:35:12+07:00",
                "production_closure_digest_sha256": production_closure_digest(
                    manifest
                ),
                "evidence": {
                    kind: {**record, "status": "PASS"}
                    for kind, record in evidence.items()
                },
            },
        }
        review_path = root / "review.json"
        review_path.write_text(json.dumps(review), encoding="utf-8")
        return manifest, review, review_path

    def test_generator_records_real_hud_and_iruka_evidence_in_canonical_closure(self):
        from tools.artqa import task6b_closure

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)

            def canonical_png(relative: str, color):
                canonical = task6b_closure.CLOSURE_RELATIVE / relative
                local = root / Path(canonical).relative_to(
                    task6b_closure.CLOSURE_RELATIVE
                )
                local.parent.mkdir(parents=True, exist_ok=True)
                Image.new("RGBA", (8, 8), color).save(local, format="PNG")
                return canonical.as_posix()

            hud_sources = {
                name: {
                    "x4": canonical_png(f"hud/{name}/x4/sprite.png", color)
                }
                for name, color in (
                    ("naruto-hp-bg", (40, 40, 40, 255)),
                    ("naruto-hp-fill", (220, 20, 20, 255)),
                    ("naruto-chakra-bg", (60, 60, 60, 255)),
                    ("naruto-chakra-fill", (20, 80, 220, 255)),
                )
            }
            sprites = [
                {
                    "id": sprite_id,
                    "zooms": {
                        "x4": canonical_png(
                            f"sprites/{sprite_id}/x4/sprite.png",
                            color,
                        )
                    },
                }
                for sprite_id, color in (
                    (27075, (180, 20, 20, 200)),
                    (27076, (20, 20, 180, 200)),
                    (27077, (20, 180, 20, 200)),
                )
            ]

            evidence = task6b_closure._generate_production_evidence(
                root, hud_sources, sprites
            )

            self.assertEqual(set(evidence), {"hud_golden", "iruka_composite"})
            for record in evidence.values():
                self.assertEqual(set(record), {"path", "sha256"})
                local = root / Path(record["path"]).relative_to(
                    task6b_closure.CLOSURE_RELATIVE
                )
                self.assertTrue(local.is_file())
                self.assertEqual(
                    record["sha256"],
                    hashlib.sha256(local.read_bytes()).hexdigest(),
                )

    def test_iruka_runtime_layers_share_transform_and_silhouette_does_not_clip(self):
        from tools.artqa import task6b_closure

        review = json.loads(REVIEW.read_text(encoding="utf-8"))
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            sprites = task6b_closure._generate_sprites(review, ART_ROOT, root)
            iruka = {
                sprite["id"]: sprite
                for sprite in sprites
                if sprite["id"] in {27075, 27076, 27077}
            }
            self.assertEqual(set(iruka), {27075, 27076, 27077})
            layers = {}
            for sprite_id, name in (
                (27075, "head"),
                (27076, "body"),
                (27077, "leg"),
            ):
                path = root / Path(
                    iruka[sprite_id]["zooms"]["x4"]
                ).relative_to(task6b_closure.CLOSURE_RELATIVE)
                with Image.open(path) as image:
                    image.load()
                    layers[name] = image.convert("RGBA")
                self.assertEqual(layers[name].size, (128, 160), name)

            output = root / "iruka-light-dark.png"
            from assetpack.visual import render_iruka_composite_contact_sheet

            render_iruka_composite_contact_sheet(layers, output)
            with Image.open(output) as contact:
                contact.load()
                self.assertEqual(contact.size, (256, 160))
                for column, color in enumerate(
                    ((245, 245, 245), (24, 24, 24))
                ):
                    panel = contact.crop((column * 128, 0, (column + 1) * 128, 160))
                    difference = ImageChops.difference(
                        panel.convert("RGB"),
                        Image.new("RGB", panel.size, color),
                    )
                    bbox = difference.getbbox()
                    self.assertIsNotNone(bbox)
                    left, top, right, bottom = bbox
                    self.assertGreaterEqual(left, 2)
                    self.assertGreaterEqual(top, 2)
                    self.assertLessEqual(right, 126)
                    self.assertLessEqual(bottom, 159)

    def test_artifact_review_accepts_independent_pass_over_exact_closure(self):
        from assetpack.build import validate_production_artifact_review

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest, _, review_path = self._artifact_review_fixture(root)
            records = validate_production_artifact_review(
                review_path, manifest, root
            )
            self.assertEqual(
                set(records),
                {"hud_golden", "iruka_composite"},
            )

    def test_artifact_review_rejects_deleted_changed_wrong_sha_nonpass_or_reviewer(self):
        from assetpack.build import (
            BuildError,
            validate_production_artifact_review,
        )

        def deleted(root, manifest, review):
            (root / "evidence/hud_golden.png").unlink()

        def changed(root, manifest, review):
            (root / "evidence/hud_golden.png").write_bytes(b"changed")

        def wrong_sha(root, manifest, review):
            review["production_artifact_review"]["evidence"]["hud_golden"][
                "sha256"
            ] = "0" * 64

        def nonpass(root, manifest, review):
            review["production_artifact_review"]["evidence"]["hud_golden"][
                "status"
            ] = "FAIL"

        def missing_reviewer(root, manifest, review):
            del review["production_artifact_review"]["reviewer"]

        def empty_reviewer(root, manifest, review):
            review["production_artifact_review"]["reviewer"] = ""

        def source_reviewer_reused(root, manifest, review):
            review["production_artifact_review"]["reviewer"] = review["reviewer"]

        def self_signed_alias(root, manifest, review):
            review["production_artifact_review"][
                "reviewer"
            ] = "Task6B closure generator self-signed alias"

        def wrong_closure(root, manifest, review):
            review["production_artifact_review"][
                "production_closure_digest_sha256"
            ] = "0" * 64

        for name, mutation in (
            ("deleted", deleted),
            ("changed", changed),
            ("wrong_sha", wrong_sha),
            ("nonpass", nonpass),
            ("missing_reviewer", missing_reviewer),
            ("empty_reviewer", empty_reviewer),
            ("source_reviewer_reused", source_reviewer_reused),
            ("self_signed_alias", self_signed_alias),
            ("wrong_closure", wrong_closure),
        ):
            with self.subTest(name=name), tempfile.TemporaryDirectory() as temporary:
                root = Path(temporary)
                manifest, review, review_path = self._artifact_review_fixture(root)
                mutation(root, manifest, review)
                review_path.write_text(json.dumps(review), encoding="utf-8")
                with self.assertRaisesRegex(
                    BuildError, "ASSET-ARTIFACT-REVIEW-016"
                ):
                    validate_production_artifact_review(
                        review_path, manifest, root
                    )


if __name__ == "__main__":
    unittest.main()
