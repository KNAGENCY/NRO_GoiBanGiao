import hashlib
import json
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch


class Task6BBindTests(unittest.TestCase):
    def _record(self, *, asset="player_idle_3f_chroma.png", slot=0):
        return {
            "asset": asset,
            "source_sha256": "a" * 64,
            "source_dimensions": [100, 100],
            "crop": {"source_rect": [0, 0, 100, 100], "output_dimensions": [100, 100]},
            "anchor": {"x": 50, "y": 99},
            "action_strip": "player_idle",
            "slot_index": slot,
            "client_cf": [0],
            "evidence": {
                "alpha": {"path": "alpha.png", "sha256": "b" * 64},
                "crop_report": {"path": "crop.json", "sha256": "c" * 64},
                "contact_sheet": {"path": "contact.png", "sha256": "d" * 64},
                "animation_preview": {"path": "preview.gif", "sha256": "e" * 64},
            },
        }

    def test_bind_calls_review_gate_and_fails_closed_without_output(self):
        from assetpack.bind import BindError, bind_manifest

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            mapping = root / "mapping.json"
            review = root / "review.json"
            runtime = root / "runtime.json"
            output = root / "canonical.json"
            mapping.write_text("{}", encoding="utf-8")
            review.write_text("{}", encoding="utf-8")
            runtime.write_text("{}", encoding="utf-8")
            with patch("assetpack.bind.require_review_gate", side_effect=RuntimeError("gate")):
                with self.assertRaisesRegex(BindError, "ASSET-REVIEW-GATE-001"):
                    bind_manifest(mapping, review, runtime, art_root=root, output=output)
            self.assertFalse(output.exists())

    def test_baseline_mapping_is_rejected_as_noncanonical_build_manifest(self):
        from assetpack.bind import BindError, validate_canonical_manifest

        mapping = json.loads(
            (Path(__file__).parents[2] / "docs/my-thuat/mapping-id.json").read_text(
                encoding="utf-8"
            )
        )
        with self.assertRaisesRegex(BindError, "ASSET-SCHEMA-001"):
            validate_canonical_manifest(mapping)

    def test_two_slots_cannot_share_one_action_slot(self):
        from assetpack.bind import BindError, validate_bind_records

        records = [self._record(slot=0), self._record(asset="player_idle_3f_chroma.png", slot=0)]
        with self.assertRaisesRegex(BindError, "ASSET-ANIMATION-008"):
            validate_bind_records(records, review_document={"assets": []}, runtime_map={"entries": []})

    def test_source_sha_change_after_review_is_rejected(self):
        from assetpack.bind import BindError, validate_bind_records

        record = self._record()
        record["review_source_sha256"] = "f" * 64
        with self.assertRaisesRegex(BindError, "ASSET-REVIEW-010"):
            validate_bind_records([record], review_document={"assets": []}, runtime_map={"entries": []})

    def test_wrong_crop_anchor_and_client_cf_are_rejected(self):
        from assetpack.bind import BindError, validate_bind_records

        for mutation, code in (
            (lambda record: record["crop"].update(source_rect=[0, 0, 0, 100]), "ASSET-CROP-007"),
            (lambda record: record["anchor"].update(x="50"), "ASSET-ANCHOR-009"),
            (lambda record: record.update(client_cf=[1]), "ASSET-ANIMATION-008"),
        ):
            with self.subTest(code=code):
                record = self._record()
                mutation(record)
                with self.assertRaisesRegex(BindError, code):
                    validate_bind_records([record], review_document={"assets": []}, runtime_map={"entries": []})

    def test_missing_runtime_key_and_missing_evidence_are_rejected(self):
        from assetpack.bind import BindError, validate_bind_records

        record = self._record()
        record["evidence"].pop("contact_sheet")
        with self.assertRaisesRegex(BindError, "ASSET-EVIDENCE-012"):
            validate_bind_records([record], review_document={"assets": []}, runtime_map={"entries": []})

        record = self._record()
        with self.assertRaisesRegex(BindError, "ASSET-UI-007"):
            validate_bind_records([record], review_document={"assets": []}, runtime_map={"entries": []})


if __name__ == "__main__":
    unittest.main()
