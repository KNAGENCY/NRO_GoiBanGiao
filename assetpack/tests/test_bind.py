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
            "crop": {"source_rect": [0, 0, 100, 100], "output_dimensions": [100, 100], "frames": [{"anchor": {"x": 50, "y": 99}, "overlap": [0, 0], "layer": 0}]},
            "anchor": {"x": 50, "y": 99},
            "overlap": [0, 0],
            "layer": 0,
            "action_strip": "player_idle",
            "slot_index": slot,
            "client_cf": [0],
            "evidence": {
                "alpha": {"path": "alpha.png", "sha256": "b" * 64},
                "crop_report": {"path": "crop.json", "sha256": "c" * 64},
                "contact_sheet": {"path": "contact.png", "sha256": "d" * 64},
                "animation_preview": {"path": "preview.gif", "sha256": "e" * 64},
            },
            "evidence_path": "crop.json",
            "evidence_sha256": "c" * 64,
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

    def test_production_mapping_is_now_a_canonical_bind_manifest(self):
        from assetpack.bind import validate_canonical_manifest

        mapping = json.loads(
            (Path(__file__).parents[2] / "docs/my-thuat/mapping-id.json").read_text(
                encoding="utf-8"
            )
        )
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

    def test_review_adapter_keeps_flat_evidence_provenance_per_asset(self):
        from assetpack.bind import _review_records

        document = {"assets": [{
            "asset": "player_idle_3f_chroma.png",
            "source_sha256": "a" * 64,
            "checks": {"source": {"dimensions": [2172, 724]}, "crop": {
                "source_size": [2172, 724],
                "source_edges": [0, 724, 1448, 2172],
                "content_boxes": [[180, 62, 520, 715], [160, 68, 495, 713], [124, 65, 461, 715]],
                "output_canvas_size": [1020, 653],
            }, "animation": {"frame_count": 3}},
            "evidence": {
                kind: {"path": f"{kind}.bin", "sha256": value * 64}
                for kind, value in (("alpha", "b"), ("crop_report", "c"), ("contact_sheet", "d"), ("animation_preview", "e"))
            },
        }]}
        record = _review_records(document)[0]
        self.assertEqual(record["evidence_path"], "crop_report.bin")
        self.assertEqual(record["evidence_sha256"], "c" * 64)

    def test_crop_report_derives_anchor_overlap_layer_for_each_frame(self):
        from assetpack.bind import _review_records

        document = {"assets": [{
            "asset": "player_idle_3f_chroma.png",
            "source_sha256": "a" * 64,
            "checks": {"source": {"dimensions": [2172, 724]}, "crop": {
                "source_size": [2172, 724],
                "source_edges": [0, 724, 1448, 2172],
                "content_boxes": [[180, 62, 520, 715], [160, 68, 495, 713], [124, 65, 461, 715]],
                "output_canvas_size": [1020, 653],
            }, "animation": {"frame_count": 3}},
            "evidence": {kind: {"path": f"{kind}.bin", "sha256": value * 64}
                        for kind, value in (("alpha", "b"), ("crop_report", "c"), ("contact_sheet", "d"), ("animation_preview", "e"))},
        }]}
        record = _review_records(document)[0]
        self.assertEqual(record["anchor"], {"x": 350, "y": 715})
        self.assertEqual(record["crop"]["frames"][1]["anchor"], {"x": 327, "y": 713})
        self.assertEqual(record["crop"]["frames"][2]["layer"], 2)
        self.assertIn("overlap", record["crop"]["frames"][0])

    def test_missing_crop_report_box_is_not_replaced_by_placeholder(self):
        from assetpack.bind import BindError, _review_records

        document = {"assets": [{
            "asset": "broken.png", "source_sha256": "a" * 64,
            "checks": {"source": {"dimensions": [10, 10]}, "crop": {"source_size": [10, 10], "content_boxes": []}},
            "evidence": {},
        }]}
        with self.assertRaisesRegex(BindError, "ASSET-CROP-007"):
            _review_records(document)

    def test_manifest_source_must_be_in_exact_reviewed_set(self):
        from assetpack.bind import BindError, validate_reviewed_sources

        with self.assertRaisesRegex(BindError, "ASSET-REVIEW-010"):
            validate_reviewed_sources(
                {"sprites": [{"source_path": "new.png"}]},
                {"assets": [{"asset": "approved.png", "source_sha256": "a" * 64}]},
            )

    def test_runtime_map_rejects_duplicate_binding_identity(self):
        from assetpack.bind import BindError, _validate_runtime_map

        entries = [{"source_asset": "same"} for _ in range(8)]
        with self.assertRaisesRegex(BindError, "ASSET-UI-007"):
            _validate_runtime_map({"entries": entries})


if __name__ == "__main__":
    unittest.main()
