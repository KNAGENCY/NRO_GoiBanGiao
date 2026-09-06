"""Kiểm thử hợp đồng QA và cổng Task 6A.

Các test dùng ảnh nhỏ tự tạo, không đụng ảnh nguồn dùng chung.
"""

from __future__ import annotations

import json
import hashlib
from pathlib import Path
import re
import subprocess
import sys
import tempfile
import unittest

from PIL import Image
import tools.artqa.artqa as artqa_module

from tools.artqa.artqa import (
    analyze_source_background,
    extract_alpha,
    normalize_cells,
    render_animation_preview,
    review,
    validate_animation,
    verify_review_gate,
)


REQUIRED_ASSET_NAMES = (
    "hu_anthan_10f_chroma.png",
    "hud_hp_chakra_4_chroma.png",
    "iruka_portrait_chroma.png",
    "master_hoc_vien_x4_chroma.png",
    "master_iruka_x4_chroma.png",
    "master_soi_rung_x4_chroma.png",
    "player_cast_3f_chroma.png",
    "player_hit_death_getup_3f_chroma.png",
    "player_idle_3f_chroma.png",
    "player_melee_4f_chroma.png",
    "player_run_4f_chroma.png",
    "skill_icons_cloud_4_chroma.png",
    "skill_icons_leaf_4_chroma.png",
    "skill_icons_mist_4_chroma.png",
    "skill_icons_sand_4_chroma.png",
    "skill_icons_stone_4_chroma.png",
    "wolf_hit_death_2f_chroma.png",
    "wolf_idle_2f_chroma.png",
    "wolf_pounce_2f_chroma.png",
    "wolf_run_4f_chroma.png",
)


class AlphaExtractionTests(unittest.TestCase):
    def test_adaptive_gradient_magenta_extracts_but_raw_source_is_rejected(self):
        image = Image.new("RGB", (25, 7))
        for y in range(image.height):
            for x in range(image.width):
                image.putpixel((x, y), (220 + x, 15 + y, 225 + x))
        image.paste((30, 80, 140), (10, 2, 15, 6))

        source_issues, _ = analyze_source_background(image)
        result = extract_alpha(image)

        self.assertIn("background_gradient_exceeds_static_matte_tolerance", source_issues)
        self.assertEqual(result.getpixel((0, 0))[3], 0)
        self.assertEqual(result.getpixel((12, 3))[3], 255)

    def test_near_white_hud_is_not_misclassified_as_chroma(self):
        image = Image.new("RGB", (7, 7), (250, 248, 246))
        image.paste((170, 20, 30), (2, 2, 5, 5))

        source_issues, _ = analyze_source_background(image)
        result = extract_alpha(image)

        self.assertIn("border_chroma_not_detected", source_issues)
        self.assertEqual(result.getpixel((0, 0))[3], 255)


class CropAndAnimationTests(unittest.TestCase):
    def test_normalize_cells_crops_real_edges_without_cutting_frame(self):
        strip = Image.new("RGBA", (11, 6), (0, 0, 0, 0))
        strip.paste((255, 0, 0, 255), (1, 1, 4, 6))
        strip.paste((0, 255, 0, 255), (7, 2, 11, 6))

        normalized, report, issues = normalize_cells(strip, columns=2)

        self.assertEqual(issues, [])
        self.assertEqual(normalized.width % 2, 0)
        self.assertEqual(report["source_edges"], [0, 6, 11])
        self.assertEqual(report["content_boxes"][1], [1, 2, 5, 6])

    def test_loop_and_timing_reject_duplicate_or_bad_duration(self):
        frames = []
        for x in (0, 0, 6):
            frame = Image.new("RGBA", (10, 6), (0, 0, 0, 0))
            frame.paste((255, 255, 255, 255), (x, 2, x + 2, 6))
            frames.append(frame)

        issues, _ = validate_animation(frames, [100, 0, 100], loop=True)

        self.assertIn("invalid_frame_timing", issues)
        self.assertIn("duplicate_adjacent_frame", issues)
        self.assertIn("loop_boundary_motion_jump", issues)

    def test_animation_preview_is_temporal_gif(self):
        frames = [
            Image.new("RGBA", (4, 4), (255, 0, 0, 255)),
            Image.new("RGBA", (4, 4), (0, 255, 0, 255)),
        ]
        with tempfile.TemporaryDirectory() as tmp:
            output = Path(tmp) / "preview.gif"
            render_animation_preview(frames, output, [80, 120], loop=True)
            with Image.open(output) as preview:
                self.assertEqual(getattr(preview, "n_frames", 1), 2)
                self.assertEqual(preview.info["duration"], 80)

    def test_one_shot_animation_preview_has_no_infinite_loop_metadata(self):
        frames = [
            Image.new("RGBA", (4, 4), (255, 0, 0, 255)),
            Image.new("RGBA", (4, 4), (0, 255, 0, 255)),
        ]
        with tempfile.TemporaryDirectory() as tmp:
            output = Path(tmp) / "one-shot.gif"
            render_animation_preview(frames, output, [80, 120], loop=False)
            with Image.open(output) as preview:
                self.assertNotEqual(preview.info.get("loop"), 0)


class GateTests(unittest.TestCase):
    def _write_review(self, root: Path) -> Path:
        source_root = root / "source"
        evidence_root = source_root / "qa"
        evidence_root.mkdir(parents=True)
        (source_root / "a.png").write_bytes(b"source")
        evidence = evidence_root / "alpha.png"
        evidence.write_bytes(b"alpha")
        digest = __import__("hashlib").sha256(evidence.read_bytes()).hexdigest()
        review = {
            "review_id": "vertical-slice-v1",
            "source_root": str(source_root),
            "assets": [{
                "asset": "a.png",
                "required": True,
                "status": "PASS",
                "evidence": {
                    "alpha": {"path": str(evidence), "sha256": digest},
                    "crop_report": {"path": str(evidence), "sha256": digest},
                    "contact_sheet": {"path": str(evidence), "sha256": digest},
                    "animation_preview": {"path": str(evidence), "sha256": digest},
                },
            }],
            "summary": {"total": 1, "pass": 1, "fail": 0},
        }
        path = root / "review.json"
        path.write_text(json.dumps(review), encoding="utf-8")
        return path

    def test_gate_rejects_evidence_sha_mismatch(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            review = self._write_review(root)
            (root / "alpha.png").write_bytes(b"tampered")

            ok, issues = verify_review_gate(review)

            self.assertFalse(ok)
            self.assertTrue(any("sha256_mismatch" in issue for issue in issues))

    def test_gate_cannot_bypass_a_required_fail_via_summary(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            review = self._write_review(root)
            data = json.loads(review.read_text(encoding="utf-8"))
            data["assets"][0]["status"] = "FAIL"
            data["summary"] = {"total": 1, "pass": 1, "fail": 0}
            review.write_text(json.dumps(data), encoding="utf-8")

            ok, issues = verify_review_gate(review)

            self.assertFalse(ok)
            self.assertIn("required_asset_not_pass:a.png", issues)

    def _write_strict_review(self, root: Path) -> Path:
        source_root = root / "vertical-slice-v1"
        qa_root = source_root / "qa"
        qa_root.mkdir(parents=True)
        assets = []
        for index, name in enumerate(REQUIRED_ASSET_NAMES):
            source = source_root / name
            Image.new("RGB", (4, 4), (220, 10 + index, 225)).save(source)
            stem = Path(name).stem.replace("_chroma", "")
            match = re.search(r"_(\d+)(?:f|_chroma)", name)
            columns = int(match.group(1)) if match else 1
            alpha = qa_root / "alpha" / f"{stem}_alpha.png"
            crop = qa_root / "crop-reports" / f"{stem}_crop.json"
            contact = qa_root / "contact-sheets" / f"{stem}_light-dark.png"
            preview = qa_root / "previews" / f"{stem}_animation.gif"
            for path in (alpha, crop, contact, preview):
                path.parent.mkdir(parents=True, exist_ok=True)
            alpha_image = Image.new("RGBA", (4 * columns, 4), (255, 0, 0, 0))
            for frame_index in range(columns):
                alpha_image.putpixel((frame_index * 4 + 1, 1), (255, 0, 0, 255))
            alpha_image.save(alpha)
            Image.new("RGBA", (8 * columns, 28), (0, 0, 0, 255)).save(contact)
            frames = [Image.new("RGBA", (4, 4), (frame_index + 1, 0, 0, 255)) for frame_index in range(columns)]
            frames[0].save(preview, save_all=True, append_images=frames[1:], duration=[100] * columns, loop=0)
            def evidence(path: Path, kind: str) -> dict:
                return {
                    "path": str(path),
                    "sha256": hashlib.sha256(path.read_bytes()).hexdigest(),
                    "asset": name,
                    "kind": kind,
                }
            crop_document = {
                "asset": name,
                "crop": {"output_canvas_size": [4 * columns, 4], "output_cell_size": [4, 4]},
                "alpha": {"transparent_ratio": 1.0, "opaque_ratio": 0.0, "soft_alpha_ratio": 0.0},
                "animation": {"frame_count": columns, "durations_ms": [100] * columns},
                "evidence": {
                    "alpha": hashlib.sha256(alpha.read_bytes()).hexdigest(),
                    "contact_sheet": hashlib.sha256(contact.read_bytes()).hexdigest(),
                    "animation_preview": hashlib.sha256(preview.read_bytes()).hexdigest(),
                },
            }
            crop.write_text(json.dumps(crop_document), encoding="utf-8")
            assets.append({
                "asset": name,
                "required": True,
                "status": "PASS",
                "source_status": "PASS",
                "source_sha256": hashlib.sha256(source.read_bytes()).hexdigest(),
                "issues": [],
                "evidence": {
                    "alpha": evidence(alpha, "alpha"),
                    "crop_report": evidence(crop, "crop_report"),
                    "contact_sheet": evidence(contact, "contact_sheet"),
                    "animation_preview": evidence(preview, "animation_preview"),
                },
                "checks": {
                    "crop": {"output_canvas_size": [4 * columns, 4], "output_cell_size": [4, 4]},
                    "animation": {"frame_count": columns, "durations_ms": [100] * columns},
                },
            })
        document = {
            "review_id": "vertical-slice-v1",
            "reviewer": "test",
            "timestamp": "2026-07-29T00:00:00Z",
            "source_root": str(source_root),
            "assets": assets,
            "summary": {"total": 20, "pass": 20, "fail": 0, "raw_source_pass": 20, "raw_source_fail": 0, "accepted_source_pass": 20, "accepted_source_fail": 0},
        }
        path = root / "review.json"
        path.write_text(json.dumps(document), encoding="utf-8")
        return path

    def _write_generated_review(self, root: Path, names=REQUIRED_ASSET_NAMES) -> Path:
        source_root = root / "vertical-slice-v1"
        source_root.mkdir(parents=True)
        for index, name in enumerate(names):
            match = re.search(r"_(\d+)(?:f|_chroma)", name)
            columns = int(match.group(1)) if match else 1
            image = Image.new("RGB", (columns * 100, 100), (220 + index % 20, 10, 225))
            for frame_index in range(columns):
                color = (30 + frame_index, 80, 140 + frame_index)
                for y in range(10, 91):
                    for x in range(10, 91):
                        if not (40 <= x <= 59 and y >= 40):
                            image.putpixel((frame_index * 100 + x, y), color)
            image.save(source_root / name)
        review_path = root / "review.json"
        review(source_root, review_path, "test reviewer")
        return review_path

    def test_gate_requires_reviewer_and_timestamp(self):
        with tempfile.TemporaryDirectory() as tmp:
            review_path = self._write_generated_review(Path(tmp))
            data = json.loads(review_path.read_text(encoding="utf-8"))
            data.pop("reviewer")
            data["timestamp"] = None
            review_path.write_text(json.dumps(data), encoding="utf-8")

            ok, issues = verify_review_gate(review_path)

            self.assertFalse(ok)
            self.assertTrue(any("reviewer_invalid" in issue for issue in issues))
            self.assertTrue(any("timestamp_invalid" in issue for issue in issues))

    def test_gate_requires_raw_source_status_summary_and_issues_to_match_source(self):
        with tempfile.TemporaryDirectory() as tmp:
            review_path = self._write_generated_review(Path(tmp))
            data = json.loads(review_path.read_text(encoding="utf-8"))
            first = data["assets"][0]
            first["raw_source_status"] = "FAIL"
            first["source_issues"] = ["forged_raw_issue"]
            data["summary"]["raw_source_pass"] = 20
            data["summary"]["raw_source_fail"] = 0
            review_path.write_text(json.dumps(data), encoding="utf-8")

            ok, issues = verify_review_gate(review_path)

            self.assertFalse(ok)
            self.assertTrue(any("raw_source_status_mismatch" in issue for issue in issues))
            self.assertTrue(any("raw_source_issues_mismatch" in issue for issue in issues))
            self.assertTrue(any("raw_source_summary" in issue for issue in issues))

    def test_gate_recomputes_all_evidence_content_not_only_hashes(self):
        evidence_kinds = ("alpha", "crop_report", "contact_sheet", "animation_preview")
        for kind in evidence_kinds:
            with self.subTest(kind=kind), tempfile.TemporaryDirectory() as tmp:
                root = Path(tmp)
                review_path = self._write_generated_review(root)
                data = json.loads(review_path.read_text(encoding="utf-8"))
                asset = data["assets"][0]
                item = asset["evidence"][kind]
                path = Path(item["path"])
                if kind == "alpha":
                    with Image.open(path) as image:
                        changed = image.convert("RGBA")
                    changed.putpixel((0, 0), (1, 2, 3, 0))
                    changed.save(path)
                elif kind == "crop_report":
                    document = json.loads(path.read_text(encoding="utf-8"))
                    document["crop"]["output_cell_size"][0] += 1
                    path.write_text(json.dumps(document), encoding="utf-8")
                elif kind == "contact_sheet":
                    with Image.open(path) as image:
                        changed = image.convert("RGBA")
                    changed.putpixel((0, 0), (1, 2, 3, 255))
                    changed.save(path)
                else:
                    with Image.open(path) as image:
                        frames = []
                        for index in range(image.n_frames):
                            image.seek(index)
                            frame = image.convert("RGBA")
                            frame.putpixel((0, 0), (1, 2, 3, 255))
                            frames.append(frame)
                        durations = [100] * len(frames)
                    frames[0].save(path, save_all=True, append_images=frames[1:], duration=durations, loop=0)
                item["sha256"] = hashlib.sha256(path.read_bytes()).hexdigest()
                if kind == "alpha":
                    crop_path = Path(asset["evidence"]["crop_report"]["path"])
                    crop_document = json.loads(crop_path.read_text(encoding="utf-8"))
                    crop_document["evidence"]["alpha"] = item["sha256"]
                    crop_path.write_text(json.dumps(crop_document), encoding="utf-8")
                    asset["evidence"]["crop_report"]["sha256"] = hashlib.sha256(crop_path.read_bytes()).hexdigest()
                review_path.write_text(json.dumps(data), encoding="utf-8")

                ok, issues = verify_review_gate(review_path)

                self.assertFalse(ok)
                self.assertTrue(any("evidence_content" in issue or "evidence_canonical" in issue for issue in issues))

    def test_gate_rejects_one_shot_gif_that_advertises_infinite_loop(self):
        with tempfile.TemporaryDirectory() as tmp:
            review_path = self._write_generated_review(Path(tmp))
            data = json.loads(review_path.read_text(encoding="utf-8"))
            asset = next(item for item in data["assets"] if item["asset"] == "hu_anthan_10f_chroma.png")
            preview_path = Path(asset["evidence"]["animation_preview"]["path"])
            with Image.open(preview_path) as image:
                frames = [image.seek(index) or image.convert("RGBA") for index in range(image.n_frames)]
                frames = [frame.copy() for frame in frames]
            frames[0].save(preview_path, save_all=True, append_images=frames[1:], duration=[100] * len(frames), loop=0)
            asset["evidence"]["animation_preview"]["sha256"] = hashlib.sha256(preview_path.read_bytes()).hexdigest()
            crop_path = Path(asset["evidence"]["crop_report"]["path"])
            crop_document = json.loads(crop_path.read_text(encoding="utf-8"))
            crop_document["evidence"]["animation_preview"] = asset["evidence"]["animation_preview"]["sha256"]
            crop_path.write_text(json.dumps(crop_document), encoding="utf-8")
            asset["evidence"]["crop_report"]["sha256"] = hashlib.sha256(crop_path.read_bytes()).hexdigest()
            review_path.write_text(json.dumps(data), encoding="utf-8")

            ok, issues = verify_review_gate(review_path)

            self.assertFalse(ok)
            self.assertTrue(any("animation_loop_semantics" in issue for issue in issues))

    def test_gate_malformed_null_assets_and_evidence_fail_closed(self):
        for malformed in ("asset", "evidence"):
            with self.subTest(malformed=malformed), tempfile.TemporaryDirectory() as tmp:
                root = Path(tmp)
                source_root = root / "source"
                (source_root / "qa").mkdir(parents=True)
                data = {
                    "review_id": "vertical-slice-v1",
                    "reviewer": "test",
                    "timestamp": "2026-07-29T00:00:00Z",
                    "source_root": str(source_root),
                    "assets": [None] if malformed == "asset" else [{"asset": "a.png", "evidence": None}],
                    "summary": {},
                }
                review_path = root / "review.json"
                review_path.write_text(json.dumps(data), encoding="utf-8")
                result = subprocess.run(
                    [sys.executable, "-m", "tools.artqa.artqa", "gate", "--review", str(review_path)],
                    capture_output=True,
                    text=True,
                    check=False,
                )
                self.assertEqual(result.returncode, 2)
                self.assertIn('"allowed": false', result.stdout)

    def test_gate_accepts_new_c1_asset_set_when_every_evidence_is_valid(self):
        with tempfile.TemporaryDirectory() as tmp:
            review = self._write_generated_review(
                Path(tmp),
                ("player_run_5f_chroma.png", "player_hit_death_getup_3f_chroma.png"),
            )

            ok, issues = verify_review_gate(review)

            self.assertTrue(ok, issues)

    def test_gate_rejects_empty_or_duplicate_asset_names(self):
        with tempfile.TemporaryDirectory() as tmp:
            review = self._write_generated_review(Path(tmp), ("player_run_5f_chroma.png",))
            data = json.loads(review.read_text(encoding="utf-8"))
            data["assets"] = []
            data["summary"] = {"total": 0, "pass": 0, "fail": 0,
                               "raw_source_pass": 0, "raw_source_fail": 0,
                               "accepted_source_pass": 0, "accepted_source_fail": 0}
            review.write_text(json.dumps(data), encoding="utf-8")

            ok, issues = verify_review_gate(review)

            self.assertFalse(ok)
            self.assertIn("required_assets_missing", issues)

            duplicate_root = Path(tmp) / "duplicate"
            duplicate_root.mkdir()
            review = self._write_generated_review(duplicate_root, ("player_run_5f_chroma.png",))
            data = json.loads(review.read_text(encoding="utf-8"))
            data["assets"].append(dict(data["assets"][0]))
            data["summary"] = {"total": 2, "pass": 2, "fail": 0,
                               "raw_source_pass": 2, "raw_source_fail": 0,
                               "accepted_source_pass": 2, "accepted_source_fail": 0}
            review.write_text(json.dumps(data), encoding="utf-8")

            ok, issues = verify_review_gate(review)

            self.assertFalse(ok)
            self.assertIn("asset_names_invalid_or_duplicate", issues)

    def test_gate_requires_source_status_and_source_sha(self):
        with tempfile.TemporaryDirectory() as tmp:
            review = self._write_strict_review(Path(tmp))
            data = json.loads(review.read_text(encoding="utf-8"))
            data["assets"][0]["source_status"] = "FAIL"
            data["assets"][1]["source_sha256"] = "0" * 64
            review.write_text(json.dumps(data), encoding="utf-8")
            ok, issues = verify_review_gate(review)
            self.assertFalse(ok)
            self.assertTrue(any("source_status_mismatch" in issue for issue in issues))
            self.assertTrue(any("source_sha256_mismatch" in issue for issue in issues))

    def test_gate_rejects_nonempty_issues_and_duplicate_evidence(self):
        with tempfile.TemporaryDirectory() as tmp:
            review = self._write_strict_review(Path(tmp))
            data = json.loads(review.read_text(encoding="utf-8"))
            first = data["assets"][0]
            first["issues"] = ["manual_failure"]
            first["evidence"]["crop_report"] = first["evidence"]["alpha"]
            review.write_text(json.dumps(data), encoding="utf-8")
            ok, issues = verify_review_gate(review)
            self.assertFalse(ok)
            self.assertTrue(any("derived_issues_mismatch" in issue for issue in issues))
            self.assertTrue(any("evidence_duplicate" in issue for issue in issues))

    def test_gate_rejects_evidence_outside_qa_root_or_wrong_kind(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            review = self._write_strict_review(root)
            data = json.loads(review.read_text(encoding="utf-8"))
            outside = root / "outside.png"
            outside.write_bytes(b"not-qa")
            item = data["assets"][0]["evidence"]["contact_sheet"]
            item["path"] = str(outside)
            item["sha256"] = hashlib.sha256(outside.read_bytes()).hexdigest()
            item["kind"] = "alpha"
            review.write_text(json.dumps(data), encoding="utf-8")
            ok, issues = verify_review_gate(review)
            self.assertFalse(ok)
            self.assertTrue(any("evidence_outside_qa_root" in issue for issue in issues))
            self.assertTrue(any("evidence_kind_mismatch" in issue for issue in issues))

    def test_bind_precondition_and_cli_fail_closed(self):
        with tempfile.TemporaryDirectory() as tmp:
            review = self._write_strict_review(Path(tmp))
            data = json.loads(review.read_text(encoding="utf-8"))
            data["assets"][0]["status"] = "FAIL"
            review.write_text(json.dumps(data), encoding="utf-8")
            precondition = getattr(artqa_module, "require_review_gate", None)
            self.assertTrue(callable(precondition), "missing bind precondition API")
            with self.assertRaises(Exception):
                precondition(review)
            result = subprocess.run(
                [sys.executable, "-m", "tools.artqa.artqa", "gate", "--review", str(review)],
                capture_output=True,
                text=True,
                check=False,
            )
            self.assertEqual(result.returncode, 2)

    def test_bind_precondition_cli_allows_only_a_strict_review(self):
        with tempfile.TemporaryDirectory() as tmp:
            review = self._write_generated_review(Path(tmp))
            result = subprocess.run(
                [sys.executable, "-m", "tools.artqa.artqa", "bind-precondition", "--review", str(review)],
                capture_output=True,
                text=True,
                check=False,
            )
            self.assertEqual(result.returncode, 0)
            self.assertIn('"allowed": true', result.stdout)


if __name__ == "__main__":
    unittest.main()
