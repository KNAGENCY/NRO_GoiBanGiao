import hashlib
import json
import tempfile
import unittest
from copy import deepcopy
from pathlib import Path
from unittest.mock import patch

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

    def test_build_calls_art_gate_before_schema_or_output(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            review = root / "review.json"
            review.write_text("{}", encoding="utf-8")
            with patch("assetpack.build.require_review_gate", side_effect=RuntimeError("gate")):
                with self.assertRaisesRegex(BuildError, "ASSET-REVIEW-GATE-001"):
                    build({}, root / "out", review, base_dir=root)
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

        def png(name, size=(1, 1)):
            path = source_dir / name
            Image.new("RGBA", size, (12, 34, 56, 255)).save(path, format="PNG")
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
            tile_art.append({"index": index, "solid": index == 1, "sources": sources})

        backgrounds = {f"b{index:02d}": {zoom: png(f"bg-{index}-{zoom}.png") for zoom in ("x1", "x2", "x3", "x4")} for index in range(4)}
        hud = {name: {zoom: png(f"{name}-{zoom}.png") for zoom in ("x1", "x2", "x3", "x4")} for name in ("naruto-hp-bg", "naruto-hp-fill", "naruto-chakra-bg", "naruto-chakra-fill")}
        ui = []
        runtime_key_entries = []
        for index in range(8):
            mode = "direct" if index == 0 else "slice_group" if index < 3 else "overlay_required"
            if mode == "direct":
                logical_keys = [f"/naruto/ui/{index}.png"]
                future_key = None
                site_role = "call"
                expression = f"loadImage('{logical_keys[0]}');"
            elif mode == "slice_group":
                logical_keys = [f"/naruto/ui/{index}-{part}.png" for part in ("left", "mid", "right")]
                future_key = None
                site_role = "call"
                expression = " ".join(f"loadImage('{key}');" for key in logical_keys)
            else:
                logical_keys = []
                future_key = f"/naruto/ui/future-{index}.png"
                site_role = "consumer"
                expression = f"consumeOverlay('{future_key}');"
            site_file_name = f"{site_role}-site-{index}.cs"
            call_site = root / "inputs" / site_file_name
            call_site.write_text(expression, encoding="utf-8")
            call_sha = hashlib.sha256(call_site.read_bytes()).hexdigest()
            physical_paths = {zoom: [f"Assets/Resources/res/{zoom}/naruto/ui/{index}.png"] for zoom in ("x1", "x2", "x3", "x4")}
            ui.append({
                "source_asset": f"ui-{index}",
                "source_path": png(f"ui-{index}.png"),
                "mode": mode,
                "logical_keys": logical_keys,
                "future_logical_key": future_key,
                "physical_paths": physical_paths,
            })
            runtime_key_entries.append({
                "source_asset": f"ui-{index}",
                "mode": mode,
                "logical_keys": logical_keys,
                "future_logical_key": future_key,
                "call_sites": [{
                    "file": f"inputs/{site_file_name}",
                    "line": 1,
                    "expression": expression,
                    "file_sha256": call_sha,
                }] if site_role == "call" else [],
                "consumer_sites": [{
                    "file": f"inputs/{site_file_name}",
                    "line": 1,
                    "expression": expression,
                    "file_sha256": call_sha,
                }] if site_role == "consumer" else [],
                "physical_paths": physical_paths,
            })
        runtime_key_map = root / "inputs" / "runtime-key-map.json"
        runtime_key_map.write_text(json.dumps({
            "schema_version": 1,
            "unity_source_scope": "fixture",
            "entries": runtime_key_entries,
        }), encoding="utf-8")

        from assetpack.tests.test_mob import GOLDEN_FRAMES, GOLDEN_IMAGES, GOLDEN_TRAILER
        mob_payload = (Path("assetpack/tests/fixtures/legacy/mob.bin")).read_bytes()
        mob_frames = {"idle": GOLDEN_FRAMES[:2], "run": GOLDEN_FRAMES[2:6], "pounce": GOLDEN_FRAMES[6:8], "hit": GOLDEN_FRAMES[8:9], "death": GOLDEN_FRAMES[9:]}
        mob_bytes = binary("mob.bytes", mob_payload)
        mob_data = binary("mob.data", mob_payload)
        mob = {"id": 0, "bytes_path": mob_bytes, "data_path": mob_data, "bytes_sha256": hashlib.sha256(mob_payload).hexdigest(), "data_sha256": hashlib.sha256(mob_payload).hexdigest(), "images": GOLDEN_IMAGES, "arr_frames": [], "trailer": GOLDEN_TRAILER, "frames": mob_frames, "img": {zoom: png(f"mob-{zoom}.png") for zoom in ("x1", "x2", "x3", "x4")}}
        animation_specs = (("hoc-vien-head", 3, (128, 128)), ("hoc-vien-body", 17, (128, 160)), ("hoc-vien-leg", 14, (128, 128)), ("quai-soi-rung", 10, (96, 96)))
        animation_sheets = []
        for name, count, cell_size in animation_specs:
            animation_sheets.append({"name": name, "slot_count": count, "cell_size": list(cell_size), "cells": [{"slot_index": index, "path": png(f"{name}-{index}.png", size=cell_size)} for index in range(count)]})

        for sprite in sprites:
            sprite["source_sha256"] = hashlib.sha256((root / sprite["source_path"]).read_bytes()).hexdigest()
        cf = {
            "hoc-vien-head": [[0, 1, 7, 8, 11, 12, 17, 18, 19, 20, 21, 22, 23, 24, 30, 31, 32], [2, 3, 4, 5, 6, 9, 10, 13, 14, 15, 16, 26, 27, 28, 29], [25]],
            "hoc-vien-body": [[0, 1], [2, 15, 28], [3], [4], [5], [6], [7, 8, 17], [9, 10], [11, 21, 32], [12, 18, 22], [13, 26], [14, 27], [16, 29], [19, 30], [20, 31], [23, 24], [25]],
            "hoc-vien-leg": [[0, 1, 22], [2], [3], [4], [5], [6], [7, 23], [8], [9], [10], [11], [12], [13, 14, 15, 16, 17, 18, 19, 20, 21, 24], [25]],
            "quai-soi-rung": [[index] for index in range(10)],
        }
        manifest = {
            "schema_version": 1,
            "source_catalog_version": 25,
            "versions": {"runtime_slice": 1, "data": 1, "map": 1, "resource": 1},
            "phase": {"name": "lang-la"},
            "sprites": sprites,
            "parts": mapping["parts"],
            "animations": [{"name": name, "slot_count": count, "frames": [{"slot_index": index, "client_cf": cf[name][index]} for index in range(count)]} for name, count, _ in animation_specs],
            "maps": [{"id": 1, "tile_id": 1, "background_id": 0}],
            "tiles": [{"id": 0, "semantic": "ground", "solid": False}],
            "backgrounds": [{"id": 0, "layers": 4}],
            "mobs": [{"id": 0, "name": "soi-rung"}],
            "skills": [{"id": 27079 + index, "pose_slot": index, "levels": list(range(12)), "effect_id": -1, "arrow_id": -1} for index in range(4)],
            "catalog": {"maps": [{"id": index + 1, "name": name, "spawn": {"x": 100, "y": 384}} for index, name in enumerate(["Làng Lá", "Làng Cát", "Làng Sương Mù", "Làng Đá", "Làng Mây", "Khu rừng"])], "npcs": [{"id": 0, "name": "Iruka", "head": 2189, "body": 2190, "leg": 2191, "portrait": 27078, "x": 172, "y": 384, "dialog": "Chào mừng em đến Làng Lá.", "menu": ["Đóng"]}], "mobs": [{"id": 0, "name": "Sói Rừng", "hp": 100, "range": 100, "speed": 0, "dart": 0, "spawns": [{"id": 0, "x": 160, "y": 384}, {"id": 1, "x": 40, "y": 384}]}], "skill_poses": [0, 1, 2, 3]},
            "ui_bindings": ui,
            "build_inputs": {
                "legacy_part_path": binary("part.bin", (Path("server/internal/game/part.bin")).read_bytes()),
                "tile_set_info_path": binary("tileset.bin", (Path("server/internal/game/tileset.bin")).read_bytes()),
                "tile_art": tile_art,
                "background_sources": backgrounds,
                "hud_sources": hud,
                "mob": mob,
                "map": {"width": 5, "height": 17, "indexes": [1] + [0] * 84, "spawn": [100, 384], "tile_size": [24, 24], "client_solid": [True] + [False] * 84},
                "animation_sheets": animation_sheets,
                "runtime_key_map": "inputs/runtime-key-map.json",
                "tile_set_baseline": {"sha256": hashlib.sha256((root / "inputs/tileset.bin").read_bytes()).hexdigest(), "bytes_hex": (root / "inputs/tileset.bin").read_bytes().hex(), "semantic": {"byte_length": len((root / "inputs/tileset.bin").read_bytes()), "tile_count": (root / "inputs/tileset.bin").read_bytes()[0]}},
            },
        }
        for binding, entry in zip(manifest["ui_bindings"], runtime_key_entries):
            binding["call_sites"] = entry["call_sites"]
            binding["consumer_sites"] = entry["consumer_sites"]
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
            required = ("server/part.bin", "server/tile_set_info.bin", "server/map-1-collision.bin", "server/skill.bin", "server/catalog.json", "client/mymap/1.bytes", "client/x1/t/1/t_01.png", "client/x2/t/1$1.png", "client/x4/bg/b03.png", "client/Mob/0.bytes", "client/x3/Mob/0/img.png", "client/x1/mainImage/naruto-hp-fill.png", "client/x4/naruto/ui/0.png", "client/evidence/tile-contact-x1.png", "client/evidence/full-map-x4.png", "client/evidence/animation/hoc-vien-body.png", "client/evidence/hud-golden-25-50-100-x4.png", "client/evidence/iruka-composite-light-dark-x4.png", "asset-lock.json")
            for relative in required:
                self.assertTrue((out / relative).is_file(), relative)
            self.assertEqual(lock["server_client_pair_sha256"], lock["server_client_pair_sha256"])
            stored = json.loads((out / "asset-lock.json").read_text(encoding="utf-8"))
            self.assertTrue(stored["outputs"])
            self.assertTrue(stored["evidence"])
            self.assertIn("provenance", stored)
            self.assertIn("runtime_key_map", stored["provenance"])
            self.assertTrue(stored["provenance"]["generated_evidence"])

    def test_mob_payload_mutation_is_rejected_even_when_manifest_sha_is_updated(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest = self._complete_manifest(root)
            mob_path = root / manifest["build_inputs"]["mob"]["bytes_path"]
            changed = bytearray(mob_path.read_bytes())
            changed[40] ^= 1
            changed = bytes(changed)
            mob_path.write_bytes(changed)
            manifest["build_inputs"]["mob"]["bytes_sha256"] = hashlib.sha256(changed).hexdigest()
            review = self._write_complete_review(root, manifest)
            with self.assertRaisesRegex(BuildError, "ASSET-MOB-003"):
                build(manifest, root / "out", review, base_dir=root)
            self.assertFalse((root / "out").exists())

    def test_build_generates_production_mob_payload_without_legacy_fixture(self):
        from assetpack.build import build
        from assetpack.mob import parse_mob

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest = self._complete_manifest(root)
            mob = manifest["build_inputs"]["mob"]
            mob["payload_mode"] = "generated"
            for key in ("bytes_path", "data_path", "bytes_sha256", "data_sha256"):
                del mob[key]
            review = self._write_complete_review(root, manifest)

            out = root / "out"
            build(manifest, out, review, base_dir=root)

            generated_bytes = (out / "client/Mob/0.bytes").read_bytes()
            generated_data = (out / "client/Mob/0/data.bytes").read_bytes()
            self.assertEqual(generated_bytes, generated_data)
            parsed = parse_mob(generated_bytes)
            self.assertEqual(
                parsed["frames"],
                [frame for action in ("idle", "run", "pounce", "hit", "death") for frame in mob["frames"][action]],
            )

    def test_ui_binding_uses_zoom_specific_sources(self):
        from assetpack.build import build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest = self._complete_manifest(root)
            expected = {}
            for binding_index, binding in enumerate(manifest["ui_bindings"]):
                sources = {}
                for zoom_index, zoom in enumerate(("x1", "x2", "x3", "x4"), start=1):
                    path = root / "inputs" / f"ui-{binding_index}-{zoom}.png"
                    Image.new("RGBA", (zoom_index, zoom_index), (12, 34, 56, 255)).save(path)
                    sources[zoom] = f"inputs/{path.name}"
                binding["sources"] = sources
                expected[binding["source_asset"]] = sources
            review = self._write_complete_review(root, manifest)

            out = root / "out"
            build(manifest, out, review, base_dir=root)

            first = manifest["ui_bindings"][0]
            for zoom_index, zoom in enumerate(("x1", "x2", "x3", "x4"), start=1):
                physical = first["physical_paths"][zoom][0].replace("\\", "/")
                marker = f"Assets/Resources/res/{zoom}/"
                relative = physical.split(marker, 1)[1]
                with Image.open(out / "client" / zoom / relative) as image:
                    self.assertEqual(image.size, (zoom_index, zoom_index))

    def test_production_source_base_requires_complete_input_sha_closure(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest = self._complete_manifest(root)
            manifest["build_inputs"]["source_base"] = "."
            review = self._write_complete_review(root, manifest)

            with self.assertRaisesRegex(BuildError, "ASSET-PROVENANCE-014"):
                build(manifest, root / "out", review, base_dir=root)

    def test_animation_build_requires_all_canonical_sheets(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest = self._complete_manifest(root)
            manifest["build_inputs"]["animation_sheets"].pop()
            review = self._write_complete_review(root, manifest)
            with self.assertRaisesRegex(BuildError, "ASSET-ANIMATION-008"):
                build(manifest, root / "out", review, base_dir=root)

    def test_animation_cf_and_cell_bounds_are_part_of_required_closure(self):
        from assetpack.build import BuildError, build

        for mutation in ("cf", "cell"):
            with self.subTest(mutation=mutation), tempfile.TemporaryDirectory() as temporary:
                root = Path(temporary)
                manifest = self._complete_manifest(root)
                if mutation == "cf":
                    manifest["animations"][0]["frames"][1]["client_cf"] = []
                else:
                    wrong = root / "inputs" / "hoc-vien-head-0.png"
                    Image.new("RGBA", (129, 128), (1, 2, 3, 255)).save(wrong, format="PNG")
                review = self._write_complete_review(root, manifest)
                with self.assertRaisesRegex(BuildError, "ASSET-ANIMATION-008"):
                    build(manifest, root / "out", review, base_dir=root)
                self.assertFalse((root / "out").exists())

    def test_runtime_key_mutation_is_rejected_by_bound_contract(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest = self._complete_manifest(root)
            manifest["ui_bindings"][0]["logical_keys"] = ["/guessed/not-a-runtime-key.png"]
            review = self._write_complete_review(root, manifest)
            with self.assertRaisesRegex(BuildError, "ASSET-UI-007"):
                build(manifest, root / "out", review, base_dir=root)
            self.assertFalse((root / "out").exists())

    def test_runtime_binding_duplicate_cannot_hide_a_missing_logical_binding(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest = self._complete_manifest(root)
            manifest["ui_bindings"][-1]["source_asset"] = manifest["ui_bindings"][0]["source_asset"]
            review = self._write_complete_review(root, manifest)
            with self.assertRaisesRegex(BuildError, "ASSET-UI-007"):
                build(manifest, root / "out", review, base_dir=root)
            self.assertFalse((root / "out").exists())

    def test_runtime_entry_3_copy_into_entry_4_is_a_collision_even_with_distinct_source_assets(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest = self._complete_manifest(root)
            runtime_map_path = root / manifest["build_inputs"]["runtime_key_map"]
            runtime_map = json.loads(runtime_map_path.read_text(encoding="utf-8"))
            source = runtime_map["entries"][3]
            target = runtime_map["entries"][4]
            for field in ("mode", "logical_keys", "future_logical_key", "physical_paths", "call_sites", "consumer_sites"):
                target[field] = deepcopy(source[field])
                manifest["ui_bindings"][4][field] = deepcopy(source[field])
            runtime_map_path.write_text(json.dumps(runtime_map), encoding="utf-8")
            review = self._write_complete_review(root, manifest)

            with self.assertRaisesRegex(BuildError, "ASSET-UI-007"):
                build(manifest, root / "out", review, base_dir=root)
            self.assertFalse((root / "out").exists())

    def test_runtime_collision_cannot_be_hidden_by_source_or_site_metadata(self):
        from assetpack.build import BuildError, build

        mutations = ("logical_keys", "future_logical_key", "physical_paths")
        for field in mutations:
            with self.subTest(field=field), tempfile.TemporaryDirectory() as temporary:
                root = Path(temporary)
                manifest = self._complete_manifest(root)
                runtime_map_path = root / manifest["build_inputs"]["runtime_key_map"]
                runtime_map = json.loads(runtime_map_path.read_text(encoding="utf-8"))

                source_index, target_index = (1, 2) if field == "logical_keys" else (3, 4)
                source = runtime_map["entries"][source_index]
                target = runtime_map["entries"][target_index]
                target[field] = deepcopy(source[field])
                manifest["ui_bindings"][target_index][field] = deepcopy(source[field])

                if field == "logical_keys":
                    target_site = target["call_sites"][0]
                    target_expression = (
                        " ".join(f"loadImage('{key}');" for key in source["logical_keys"])
                        + " loadImage('/naruto/ui/independent-target-marker.png');"
                    )
                    (root / target_site["file"]).write_text(target_expression, encoding="utf-8")
                    target_site["expression"] = target_expression
                    target_site["file_sha256"] = hashlib.sha256((root / target_site["file"]).read_bytes()).hexdigest()
                    manifest["ui_bindings"][target_index]["call_sites"] = deepcopy(target["call_sites"])
                elif field == "future_logical_key":
                    target_site = target["consumer_sites"][0]
                    target_expression = f"consumeOverlay('{source[field]}'); consumeOverlay('/naruto/ui/independent-target-marker.png');"
                    (root / target_site["file"]).write_text(target_expression, encoding="utf-8")
                    target_site["expression"] = target_expression
                    target_site["file_sha256"] = hashlib.sha256((root / target_site["file"]).read_bytes()).hexdigest()
                    manifest["ui_bindings"][target_index]["consumer_sites"] = deepcopy(target["consumer_sites"])

                runtime_map_path.write_text(json.dumps(runtime_map), encoding="utf-8")
                review = self._write_complete_review(root, manifest)

                with self.assertRaisesRegex(BuildError, "ASSET-UI-007"):
                    build(manifest, root / "out", review, base_dir=root)
                self.assertFalse((root / "out").exists())

    def test_tileset_payload_mutation_is_rejected_even_when_sha_is_updated(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest = self._complete_manifest(root)
            tileset_path = root / manifest["build_inputs"]["tile_set_info_path"]
            changed = bytearray(tileset_path.read_bytes())
            changed[-1] ^= 1
            tileset_path.write_bytes(changed)
            manifest["build_inputs"]["tile_set_info_path"] = manifest["build_inputs"]["tile_set_info_path"]
            manifest["build_inputs"]["tile_set_baseline"]["sha256"] = hashlib.sha256(changed).hexdigest()
            review = self._write_complete_review(root, manifest)
            with self.assertRaisesRegex(BuildError, "ASSET-TILESET-011"):
                build(manifest, root / "out", review, base_dir=root)

    def test_malformed_tileset_returns_asset_error_not_runtime_exception(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest = self._complete_manifest(root)
            (root / manifest["build_inputs"]["tile_set_info_path"]).write_bytes(b"")
            review = self._write_complete_review(root, manifest)
            with self.assertRaisesRegex(BuildError, "ASSET-TILESET-010"):
                build(manifest, root / "out", review, base_dir=root)

    def test_zero_tile_size_returns_stable_build_error(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest = self._complete_manifest(root)
            manifest["build_inputs"]["map"]["tile_size"] = [0, 24]
            review = self._write_complete_review(root, manifest)
            with self.assertRaisesRegex(BuildError, "ASSET-MAP-004"):
                build(manifest, root / "out", review, base_dir=root)

    def test_mutation_solid_tile_without_matching_collision_is_rejected(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest = self._complete_manifest(root)
            manifest["build_inputs"]["tile_art"][0]["solid"] = True
            review = self._write_complete_review(root, manifest)
            with self.assertRaisesRegex(BuildError, "ASSET-COLLISION-006"):
                build(manifest, root / "out", review, base_dir=root)

    def test_mutation_erasing_all_collision_is_rejected(self):
        from assetpack.build import BuildError, build

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            manifest = self._complete_manifest(root)
            manifest["build_inputs"]["map"]["indexes"] = [0] * 85
            manifest["build_inputs"]["map"]["client_solid"] = [False] * 85
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
