import hashlib
import json
import tempfile
import unittest
from pathlib import Path

from PIL import Image


class BuildContractTests(unittest.TestCase):
    def _review(self, root: Path, source: Path, *, source_sha=None):
        evidence = root / "golden.json"
        evidence.write_bytes(b"golden evidence")
        review = root / "review.json"
        review.write_text(json.dumps({
            "assets": [{
                "asset": source.name,
                "status": "PASS",
                "source_sha256": source_sha or hashlib.sha256(source.read_bytes()).hexdigest(),
                "evidence_path": str(evidence),
                "evidence_sha256": hashlib.sha256(evidence.read_bytes()).hexdigest(),
            }]
        }), encoding="utf-8")
        return review

    def test_minimal_manifest_is_rejected_before_any_output_is_bound(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            source = root / "source.png"
            source.write_bytes(b"source")
            review = self._review(root, source)
            with self.assertRaisesRegex(BuildError, "ASSET-SCHEMA-001"):
                build({"slice": "lang-la"}, root / "out", review)
            self.assertFalse((root / "out").exists())

    def test_review_source_sha_must_match_the_bytes_on_disk(self):
        from assetpack.build import BuildError, require_approved_review

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            source = root / "source.png"
            source.write_bytes(b"actual source")
            review = self._review(root, source, source_sha="0" * 64)
            manifest = {"sprites": [{"source_path": source.name, "source_sha256": "0" * 64}], "build_inputs": {}}
            with self.assertRaisesRegex(BuildError, "ASSET-REVIEW-010"):
                require_approved_review(review, manifest, root)

    def test_existing_output_is_rejected_without_overwrite(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            source = root / "source.png"
            source.write_bytes(b"source")
            review = root / "review.json"
            review.write_text(json.dumps({"assets": [{"status": "FAIL"}]}), encoding="utf-8")
            out = root / "out"
            out.mkdir()
            sentinel = out / "sentinel"
            sentinel.write_bytes(b"keep")
            with self.assertRaisesRegex(BuildError, "ASSET-BUILD-012"):
                build({"slice": "lang-la"}, out, review)
            self.assertEqual(b"keep", sentinel.read_bytes())

    def test_lock_carries_artifact_metadata_and_evidence_hashes(self):
        from assetpack.lock import lock_bytes

        data = json.loads(lock_bytes({
            "versions": {"runtime_slice": 1},
            "inputs": [{"path": "source.png", "sha256": "a" * 64}],
            "outputs": [{"path": "server/part.bin", "sha256": "b" * 64, "bytes": 10}],
            "evidence": [{"path": "contact.png", "sha256": "c" * 64}],
            "closure": {"sprite_ids": [27041], "part_ids": [30], "frame_counts": {"30": 3}},
        }).decode("utf-8"))
        self.assertEqual(data["outputs"][0]["bytes"], 10)
        self.assertEqual(data["evidence"][0]["sha256"], "c" * 64)
        self.assertEqual(data["closure"]["frame_counts"]["30"], 3)

    def _complete_manifest(self, root: Path):
        mapping = json.loads(Path("docs/my-thuat/mapping-id.json").read_text(encoding="utf-8"))
        source_dir = root / "inputs"
        source_dir.mkdir()

        def png(name):
            path = source_dir / name
            Image.new("RGBA", (1, 1), (12, 34, 56, 255)).save(path, format="PNG")
            return f"inputs/{name}"

        def binary(name, data):
            path = source_dir / name
            path.write_bytes(data)
            return f"inputs/{name}"

        sprites = []
        for sprite in mapping["sprites"]:
            path = png(f"sprite-{sprite['id']}.png")
            sprites.append({"id": sprite["id"], "source_path": path, "source_sha256": "", "zooms": {zoom: path for zoom in ("x1", "x2", "x3", "x4")}})

        tiles = []
        tile_art = []
        for index in range(100):
            sources = {zoom: png(f"tile-{zoom}-{index}.png") for zoom in ("x1", "x2", "x3", "x4")}
            tile_art.append({"index": index, "solid": False, "sources": sources})

        backgrounds = {f"b{index:02d}": {zoom: png(f"bg-{index}-{zoom}.png") for zoom in ("x1", "x2", "x3", "x4")} for index in range(4)}
        hud = {name: {zoom: png(f"{name}-{zoom}.png") for zoom in ("x1", "x2", "x3", "x4")} for name in ("naruto-hp-bg", "naruto-hp-fill", "naruto-chakra-bg", "naruto-chakra-fill")}
        ui = []
        for index in range(8):
            ui.append({
                "source_asset": f"ui-{index}",
                "source_path": png(f"ui-{index}.png"),
                "mode": "direct",
                "logical_keys": [f"/naruto/ui/{index}.png"],
                "future_logical_key": None,
                "physical_paths": {zoom: [f"Assets/Resources/res/{zoom}/naruto/ui/{index}.png"] for zoom in ("x1", "x2", "x3", "x4")},
            })

        mob_payload = json.dumps([0, {"idle": [0, 1], "run": [2, 3, 4, 5], "pounce": [6, 7], "hit": [8], "death": [9]}], separators=(",", ":")).encode()
        mob_bytes = binary("mob.bytes", mob_payload)
        mob_data = binary("mob.data", mob_payload)
        mob = {"id": 0, "bytes_path": mob_bytes, "data_path": mob_data, "bytes_sha256": hashlib.sha256(mob_payload).hexdigest(), "data_sha256": hashlib.sha256(mob_payload).hexdigest(), "frames": {"idle": [0, 1], "run": [2, 3, 4, 5], "pounce": [6, 7], "hit": [8], "death": [9]}, "img": {zoom: png(f"mob-{zoom}.png") for zoom in ("x1", "x2", "x3", "x4")}}
        animation_cells = [{"slot_index": index, "path": png(f"animation-{index}.png")} for index in range(3)]

        for sprite in sprites:
            sprite["source_sha256"] = hashlib.sha256((root / sprite["source_path"]).read_bytes()).hexdigest()
        manifest = {
            "schema_version": 1,
            "source_catalog_version": 25,
            "versions": {"runtime_slice": 1, "data": 1, "map": 1, "resource": 1},
            "phase": {"name": "lang-la"},
            "sprites": sprites,
            "parts": mapping["parts"],
            "animations": [{"name": "idle", "slot_count": 3, "frames": [{"slot_index": 0, "client_cf": [0]}, {"slot_index": 1, "client_cf": [1]}, {"slot_index": 2, "client_cf": [2]}]}],
            "maps": [{"id": 1, "tile_id": 1, "background_id": 0}],
            "tiles": [{"id": 0, "semantic": "ground", "solid": False}],
            "backgrounds": [{"id": 0, "layers": 4}],
            "mobs": [{"id": 0, "name": "soi-rung"}],
            "skills": [{"id": 27079 + index, "pose_slot": index, "levels": list(range(12)), "effect_id": -1, "arrow_id": -1} for index in range(4)],
            "catalog": {"maps": ["Làng Lá", "Làng Cát", "Làng Sương Mù", "Làng Đá", "Làng Mây", "Khu rừng"], "npc_templates": ["Iruka"], "mob_templates": ["Sói Rừng"], "spawn": [[100, 384]], "skill_poses": [0, 1, 2, 3]},
            "ui_bindings": ui,
            "build_inputs": {
                "legacy_part_path": binary("part.bin", (Path("server/internal/game/part.bin")).read_bytes()),
                "tile_set_info_path": binary("tileset.bin", (Path("server/internal/game/tileset.bin")).read_bytes()),
                "tile_art": tile_art,
                "background_sources": backgrounds,
                "hud_sources": hud,
                "mob": mob,
                "map": {"width": 5, "height": 17, "indexes": [0] * 85, "spawn": [100, 384], "tile_size": [24, 24], "client_solid": [False] * 85},
                "animation_sheets": [{"name": "idle", "slot_count": 3, "cell_size": [1, 1], "cells": animation_cells}],
            },
        }
        return manifest

    def _write_complete_review(self, root: Path, manifest):
        from assetpack.build import _source_paths

        evidence = root / "golden.png"
        Image.new("RGBA", (1, 1), (255, 255, 255, 255)).save(evidence, format="PNG")
        entries = []
        for relative in _source_paths(manifest):
            source = root / relative
            entries.append({"asset": source.name, "status": "PASS", "source_sha256": hashlib.sha256(source.read_bytes()).hexdigest(), "evidence_path": str(evidence), "evidence_sha256": hashlib.sha256(evidence.read_bytes()).hexdigest()})
        review = root / "review.json"
        review.write_text(json.dumps({"assets": entries}), encoding="utf-8")
        return review

    def test_build_emits_server_client_artifacts_and_pair_lock(self):
        from assetpack.build import build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest = self._complete_manifest(root)
            review = self._write_complete_review(root, manifest)
            out = root / "out"
            lock = build(manifest, out, review, base_dir=root)
            second = build(manifest, root / "out-second", review, base_dir=root)
            self.assertEqual(lock["server_client_pair_sha256"], second["server_client_pair_sha256"])
            self.assertEqual((out / "asset-lock.json").read_bytes(), (root / "out-second" / "asset-lock.json").read_bytes())
            required = ("server/part.bin", "server/tile_set_info.bin", "server/map-1-collision.bin", "server/skill.bin", "server/catalog.json", "client/mymap/1.bytes", "client/x1/t/1/t_01.png", "client/x2/t/1$1.png", "client/x4/bg/b03.png", "client/Mob/0.bytes", "client/x3/Mob/0/img.png", "client/x1/mainImage/naruto-hp-fill.png", "client/x4/naruto/ui/0.png", "client/evidence/tile-contact-x1.png", "client/evidence/full-map-x4.png", "client/evidence/animation/idle.png", "asset-lock.json")
            for relative in required:
                self.assertTrue((out / relative).is_file(), relative)
            self.assertEqual(lock["server_client_pair_sha256"], lock["server_client_pair_sha256"])
            stored = json.loads((out / "asset-lock.json").read_text(encoding="utf-8"))
            self.assertTrue(stored["outputs"])
            self.assertTrue(stored["evidence"])

    def test_mutation_solid_tile_without_matching_collision_is_rejected(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest = self._complete_manifest(root)
            manifest["build_inputs"]["tile_art"][0]["solid"] = True
            review = self._write_complete_review(root, manifest)
            with self.assertRaisesRegex(BuildError, "ASSET-COLLISION-006"):
                build(manifest, root / "out", review, base_dir=root)

    def test_mutation_deleting_b03_is_rejected_with_background_code(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest = self._complete_manifest(root)
            del manifest["build_inputs"]["background_sources"]["b03"]
            review = self._write_complete_review(root, manifest)
            with self.assertRaisesRegex(BuildError, "ASSET-BG-005"):
                build(manifest, root / "out", review, base_dir=root)

    def test_mutation_bool_and_unknown_skill_id_are_rejected_with_stable_codes(self):
        from assetpack.build import BuildError, build

        for mutation, code in ((lambda m: m["versions"].update(data=True), "ASSET-SCHEMA-001"), (lambda m: m["skills"][0].update(id=9000), "ASSET-SKILL-004")):
            with self.subTest(code=code), tempfile.TemporaryDirectory() as temporary:
                root = Path(temporary)
                manifest = self._complete_manifest(root)
                mutation(manifest)
                review = self._write_complete_review(root, manifest)
                with self.assertRaisesRegex(BuildError, code):
                    build(manifest, root / "out", review, base_dir=root)


if __name__ == "__main__":
    unittest.main()
