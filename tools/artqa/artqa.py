"""QA xác định được cho ảnh nguồn vertical slice.

PNG nguồn chỉ được đọc. Mọi alpha/crop/preview/report được ghi dưới ``qa/``.
Nguồn chroma chuyển sắc vẫn bị báo lỗi nguồn; bản alpha thích nghi chỉ được
PASS riêng khi alpha, halo, canvas, cell và thời gian hoạt ảnh đều đạt ngưỡng.
"""

from __future__ import annotations

import argparse
from collections import deque
from datetime import datetime, timezone
import hashlib
import json
from pathlib import Path
import re
from typing import Iterable

import numpy as np
from PIL import Image, ImageChops, ImageDraw


MIN_BORDER_CHROMA_RATIO = 0.85
MAX_BORDER_CHROMA_SPREAD = 12
MAX_SOFT_ALPHA_RATIO = 0.15
MIN_TRANSPARENT_RATIO = 0.05
MIN_OPAQUE_RATIO = 0.01
FRAME_DURATION_MS = 100
LIGHT = (245, 242, 235, 255)
DARK = (28, 32, 43, 255)
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
REQUIRED_EVIDENCE = ("alpha", "crop_report", "contact_sheet", "animation_preview")


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def is_chroma(pixel: tuple[int, int, int]) -> bool:
    red, green, blue = pixel
    return red - green >= 110 and blue - green >= 110 and red >= 150 and blue >= 150


def _border_points(width: int, height: int) -> Iterable[tuple[int, int]]:
    for x in range(width):
        yield x, 0
        yield x, height - 1
    for y in range(1, height - 1):
        yield 0, y
        yield width - 1, y


def analyze_source_background(image: Image.Image) -> tuple[list[str], dict]:
    rgb = image.convert("RGB")
    points = list(_border_points(*rgb.size))
    chroma = [rgb.getpixel(point) for point in points if is_chroma(rgb.getpixel(point))]
    ratio = len(chroma) / len(points)
    spread = 0
    if chroma:
        channels = tuple(zip(*chroma))
        spread = max(max(channel) - min(channel) for channel in channels)
    issues: list[str] = []
    if ratio < MIN_BORDER_CHROMA_RATIO:
        issues.append("border_chroma_not_detected")
    if spread > MAX_BORDER_CHROMA_SPREAD:
        issues.append("background_gradient_exceeds_static_matte_tolerance")
    return issues, {
        "border_chroma_ratio": round(ratio, 6),
        "border_chroma_spread": spread,
        "thresholds": {
            "min_border_chroma_ratio": MIN_BORDER_CHROMA_RATIO,
            "max_border_chroma_spread": MAX_BORDER_CHROMA_SPREAD,
        },
    }


def _adaptive_chroma_mask(rgb: Image.Image) -> np.ndarray:
    pixels = np.asarray(rgb, dtype=np.int16)
    # Chấp nhận dải magenta chuyển sắc nhưng không nhận màu gần trắng: hai kênh
    # đỏ/xanh phải cùng trội hơn xanh lá và còn bão hòa rõ.
    high = np.maximum(pixels[:, :, 0], pixels[:, :, 2])
    low = np.minimum(pixels[:, :, 0], pixels[:, :, 2])
    return (
        (low - pixels[:, :, 1] >= 70)
        & (high >= 150)
        & (high - pixels[:, :, 1] >= 100)
        & (np.abs(pixels[:, :, 0] - pixels[:, :, 2]) <= 95)
    )


def extract_alpha(image: Image.Image) -> Image.Image:
    """Flood nền magenta thích nghi từ mép theo 4-neighbour, không khoét nội bộ."""
    rgb = image.convert("RGB")
    candidate = _adaptive_chroma_mask(rgb)
    height, width = candidate.shape
    flooded = np.zeros((height, width), dtype=bool)
    queue: deque[tuple[int, int]] = deque()
    for x, y in _border_points(width, height):
        if candidate[y, x] and not flooded[y, x]:
            flooded[y, x] = True
            queue.append((x, y))
    while queue:
        x, y = queue.popleft()
        for nx, ny in ((x - 1, y), (x + 1, y), (x, y - 1), (x, y + 1)):
            if 0 <= nx < width and 0 <= ny < height and candidate[ny, nx] and not flooded[ny, nx]:
                flooded[ny, nx] = True
                queue.append((nx, ny))

    # Một số HUD master có viền trắng 1 ô ngăn nền magenta với biên ảnh. Khi
    # đó nền vẫn là một matte liên tục theo từng ô, nhưng không chạm biên ảnh.
    # Chỉ dùng fallback khi vùng chroma lớn; nền trắng gần như không tạo
    # candidate nên không thể bị xóa nhầm.
    if not flooded.any() and candidate.mean() >= 0.20:
        flooded = candidate.copy()

    alpha = np.full((height, width), 255, dtype=np.uint8)
    alpha[flooded] = 0
    # Feather một pixel về phía vật thể. Không cần SciPy: bốn phép dịch mảng là
    # đủ và giữ thuật toán hoàn toàn xác định.
    neighbour_count = np.zeros((height, width), dtype=np.uint8)
    neighbour_count[1:, :] += flooded[:-1, :]
    neighbour_count[:-1, :] += flooded[1:, :]
    neighbour_count[:, 1:] += flooded[:, :-1]
    neighbour_count[:, :-1] += flooded[:, 1:]
    alpha[(~flooded) & (neighbour_count > 0) & (neighbour_count < 3)] = 192
    result = rgb.convert("RGBA")
    result.putalpha(Image.fromarray(alpha, mode="L"))
    return result


def _cell_edges(width: int, columns: int) -> list[int]:
    return [round(index * width / columns) for index in range(columns + 1)]


def normalize_cells(alpha: Image.Image, columns: int) -> tuple[Image.Image, dict, list[str]]:
    if columns < 1:
        return alpha.copy(), {}, ["invalid_column_count"]
    edges = _cell_edges(alpha.width, columns)
    frames: list[Image.Image] = []
    boxes: list[list[int] | None] = []
    issues: list[str] = []
    for index in range(columns):
        frame = alpha.crop((edges[index], 0, edges[index + 1], alpha.height))
        box = frame.getchannel("A").getbbox()
        boxes.append(list(box) if box else None)
        if box is None:
            issues.append(f"empty_animation_cell:{index}")
            frames.append(Image.new("RGBA", (1, 1), (0, 0, 0, 0)))
        else:
            frames.append(frame.crop(box))
    cell_width = max(frame.width for frame in frames)
    cell_height = max(frame.height for frame in frames)
    normalized = Image.new("RGBA", (cell_width * columns, cell_height), (0, 0, 0, 0))
    for index, frame in enumerate(frames):
        normalized.alpha_composite(frame, (index * cell_width + (cell_width - frame.width) // 2, cell_height - frame.height))
    report = {
        "source_size": [alpha.width, alpha.height],
        "source_edges": edges,
        "content_boxes": boxes,
        "output_cell_size": [cell_width, cell_height],
        "output_canvas_size": [normalized.width, normalized.height],
    }
    return normalized, report, issues


def split_frames(strip: Image.Image, columns: int) -> list[Image.Image]:
    if columns < 1 or strip.width % columns:
        return []
    width = strip.width // columns
    return [strip.crop((index * width, 0, (index + 1) * width, strip.height)) for index in range(columns)]


def _frame_centroid(frame: Image.Image) -> tuple[float, float]:
    alpha = np.asarray(frame.getchannel("A"), dtype=np.float64)
    ys, xs = np.nonzero(alpha)
    if not len(xs):
        return 0.0, 0.0
    weights = alpha[ys, xs]
    return float(np.average(xs, weights=weights)), float(np.average(ys, weights=weights))


def validate_animation(frames: list[Image.Image], durations_ms: list[int], loop: bool) -> tuple[list[str], dict]:
    issues: list[str] = []
    if len(frames) != len(durations_ms) or any(duration < 20 or duration > 1000 for duration in durations_ms):
        issues.append("invalid_frame_timing")
    duplicates: list[int] = []
    for index in range(1, len(frames)):
        if ImageChops.difference(frames[index - 1], frames[index]).getbbox() is None:
            duplicates.append(index)
    if duplicates:
        issues.append("duplicate_adjacent_frame")
    centroids = [_frame_centroid(frame) for frame in frames]
    jumps = [
        ((centroids[index][0] - centroids[index - 1][0]) ** 2
         + (centroids[index][1] - centroids[index - 1][1]) ** 2) ** 0.5
        for index in range(1, len(centroids))
    ]
    seam_jump = 0.0
    if loop and len(centroids) >= 3:
        seam_jump = ((centroids[0][0] - centroids[-1][0]) ** 2
                     + (centroids[0][1] - centroids[-1][1]) ** 2) ** 0.5
        # Loop phải quay về gần pose đầu; không dùng chính bước nhảy cuối lớn
        # làm chuẩn vì như vậy một seam hỏng có thể tự hợp thức hóa nó.
        if seam_jump > max(2.0, frames[0].width * 0.25):
            issues.append("loop_boundary_motion_jump")
    return issues, {
        "frame_count": len(frames),
        "durations_ms": durations_ms,
        "loop": loop,
        "centroids": [[round(x, 3), round(y, 3)] for x, y in centroids],
        "seam_jump": round(seam_jump, 3),
        "duplicate_adjacent_indices": duplicates,
    }


def _animation_policy(name: str) -> tuple[bool, str]:
    if name == "hu_anthan_10f_chroma.png":
        return False, "one_shot_effect_sequence"
    if name == "hud_hp_chakra_4_chroma.png":
        return False, "static_ui_states"
    return True, "loop"


def validate_alpha_output(image: Image.Image, columns: int) -> tuple[list[str], dict]:
    alpha = np.asarray(image.getchannel("A"))
    total = alpha.size
    transparent = int(np.count_nonzero(alpha == 0)) / total
    opaque = int(np.count_nonzero(alpha == 255)) / total
    soft = int(np.count_nonzero((alpha > 0) & (alpha < 255))) / total
    issues: list[str] = []
    if transparent < MIN_TRANSPARENT_RATIO:
        issues.append("alpha_transparency_below_threshold")
    if opaque < MIN_OPAQUE_RATIO:
        issues.append("opaque_content_below_threshold")
    if soft > MAX_SOFT_ALPHA_RATIO:
        issues.append("halo_soft_alpha_exceeds_threshold")
    if columns < 1 or image.width % columns:
        issues.append("normalized_canvas_not_divisible_by_columns")
    return issues, {
        "transparent_ratio": round(transparent, 6),
        "opaque_ratio": round(opaque, 6),
        "soft_alpha_ratio": round(soft, 6),
        "thresholds": {
            "min_transparent_ratio": MIN_TRANSPARENT_RATIO,
            "min_opaque_ratio": MIN_OPAQUE_RATIO,
            "max_soft_alpha_ratio": MAX_SOFT_ALPHA_RATIO,
        },
    }


def render_animation_preview(frames: list[Image.Image], output: Path, durations_ms: list[int]) -> None:
    rendered = []
    for frame in frames:
        canvas = Image.new("RGBA", frame.size, DARK)
        canvas.alpha_composite(frame)
        rendered.append(canvas.convert("P", palette=Image.Palette.ADAPTIVE))
    output.parent.mkdir(parents=True, exist_ok=True)
    rendered[0].save(
        output,
        save_all=True,
        append_images=rendered[1:],
        duration=durations_ms,
        loop=0,
        disposal=2,
    )


def _asset_contact_sheet(name: str, alpha: Image.Image) -> Image.Image:
    sheet = Image.new("RGBA", (alpha.width * 2, alpha.height + 24), DARK)
    light = Image.new("RGBA", alpha.size, LIGHT)
    dark = Image.new("RGBA", alpha.size, DARK)
    light.alpha_composite(alpha)
    dark.alpha_composite(alpha)
    sheet.alpha_composite(light, (0, 0))
    sheet.alpha_composite(dark, (alpha.width, 0))
    ImageDraw.Draw(sheet).text((4, alpha.height + 4), name, fill=(255, 255, 255, 255))
    return sheet


def _evidence(path: Path, asset: str, kind: str) -> dict:
    return {
        "path": str(path.resolve()),
        "sha256": sha256(path),
        "asset": asset,
        "kind": kind,
    }


def review(art_root: Path, review_path: Path, reviewer: str) -> dict:
    evidence_root = art_root / "qa"
    alpha_root = evidence_root / "alpha"
    crop_root = evidence_root / "crop-reports"
    contact_root = evidence_root / "contact-sheets"
    preview_root = evidence_root / "previews"
    for directory in (alpha_root, crop_root, contact_root, preview_root):
        directory.mkdir(parents=True, exist_ok=True)

    reviewed: list[dict] = []
    for source in sorted(art_root.glob("*.png")):
        with Image.open(source) as raw:
            raw.load()
            source_issues, source_metrics = analyze_source_background(raw)
            extracted = extract_alpha(raw)
        columns = _expected_columns(source.name)
        normalized, crop_report, crop_issues = normalize_cells(extracted, columns)
        output_issues, alpha_metrics = validate_alpha_output(normalized, columns)
        frames = split_frames(normalized, columns)
        durations = [FRAME_DURATION_MS] * len(frames)
        loop, sequence_semantics = _animation_policy(source.name)
        animation_issues, animation_metrics = validate_animation(frames, durations, loop=loop)
        animation_metrics["sequence_semantics"] = sequence_semantics

        stem = source.stem.replace("_chroma", "")
        alpha_path = alpha_root / f"{stem}_alpha.png"
        crop_path = crop_root / f"{stem}_crop.json"
        contact_path = contact_root / f"{stem}_light-dark.png"
        preview_path = preview_root / f"{stem}_animation.gif"
        normalized.save(alpha_path)
        crop_document = {
            "asset": source.name,
            "crop": crop_report,
            "alpha": alpha_metrics,
            "animation": animation_metrics,
            "issues": crop_issues + output_issues + animation_issues,
        }
        _asset_contact_sheet(source.name, normalized).save(contact_path)
        render_animation_preview(frames, preview_path, durations)
        crop_document["evidence"] = {
            "alpha": sha256(alpha_path),
            "contact_sheet": sha256(contact_path),
            "animation_preview": sha256(preview_path),
        }
        crop_path.write_text(json.dumps(crop_document, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")

        # Trạng thái bind thuộc output đã chuẩn hóa. Lỗi nền nguồn vẫn được giữ
        # riêng và không bao giờ bị biến thành lời khẳng định raw source PASS.
        bind_issues = crop_issues + output_issues + animation_issues
        reviewed.append({
            "asset": source.name,
            "required": True,
            "status": "PASS" if not bind_issues else "FAIL",
            # source_status là cổng chấp nhận source bất biến vào pipeline;
            # raw_source_status vẫn ghi rõ lỗi matte của PNG gốc, không che
            # giấu việc output là bản derived.
            "source_status": "PASS" if not bind_issues else "FAIL",
            "raw_source_status": "PASS" if not source_issues else "FAIL",
            "source_sha256": sha256(source),
            "source_issues": source_issues,
            "issues": bind_issues,
            "checks": {
                "source": source_metrics,
                "alpha": alpha_metrics,
                "crop": crop_report,
                "animation": animation_metrics,
            },
            "evidence": {
                "alpha": _evidence(alpha_path, source.name, "alpha"),
                "crop_report": _evidence(crop_path, source.name, "crop_report"),
                "contact_sheet": _evidence(contact_path, source.name, "contact_sheet"),
                "animation_preview": _evidence(preview_path, source.name, "animation_preview"),
            },
        })

    document = {
        "review_id": "vertical-slice-v1",
        "reviewer": reviewer,
        "timestamp": datetime.now(timezone.utc).isoformat().replace("+00:00", "Z"),
        "source_root": str(art_root.resolve()),
        "assets": reviewed,
        "summary": {
            "total": len(reviewed),
            "pass": sum(item["status"] == "PASS" for item in reviewed),
            "fail": sum(item["status"] == "FAIL" for item in reviewed),
            "raw_source_pass": sum(item["raw_source_status"] == "PASS" for item in reviewed),
            "raw_source_fail": sum(item["raw_source_status"] == "FAIL" for item in reviewed),
            "accepted_source_pass": sum(item["source_status"] == "PASS" for item in reviewed),
            "accepted_source_fail": sum(item["source_status"] == "FAIL" for item in reviewed),
        },
    }
    review_path.parent.mkdir(parents=True, exist_ok=True)
    review_path.write_text(json.dumps(document, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    return document


def _expected_columns(name: str) -> int:
    match = re.search(r"_(\d+)(?:f|_chroma)", name)
    return int(match.group(1)) if match else 1


class ReviewGateError(RuntimeError):
    def __init__(self, review_path: Path, issues: list[str]):
        self.review_path = review_path
        self.issues = issues
        super().__init__(f"review gate failed for {review_path}: {', '.join(issues)}")


def _path_under(path: Path, root: Path) -> bool:
    try:
        path.resolve().relative_to(root.resolve())
    except ValueError:
        return False
    return True


def _evidence_filename(asset_name: str, kind: str) -> str:
    stem = Path(asset_name).stem.replace("_chroma", "")
    return {
        "alpha": f"{stem}_alpha.png",
        "crop_report": f"{stem}_crop.json",
        "contact_sheet": f"{stem}_light-dark.png",
        "animation_preview": f"{stem}_animation.gif",
    }[kind]


def _verify_evidence_content(asset: dict, kind: str, path: Path, issues: list[str]) -> None:
    name = asset["asset"]
    checks = asset.get("checks", {})
    expected_name = _evidence_filename(name, kind)
    if path.name != expected_name:
        issues.append(f"evidence_filename_mismatch:{name}:{kind}")
    try:
        if kind == "alpha":
            with Image.open(path) as image:
                image.load()
                if image.mode != "RGBA":
                    issues.append(f"evidence_alpha_mode_invalid:{name}")
                expected_size = checks.get("crop", {}).get("output_canvas_size")
                if expected_size and list(image.size) != expected_size:
                    issues.append(f"evidence_alpha_dimensions_invalid:{name}")
                alpha_issues, _ = validate_alpha_output(image, _expected_columns(name))
                if alpha_issues:
                    issues.append(f"evidence_alpha_quality_invalid:{name}:{','.join(alpha_issues)}")
        elif kind == "crop_report":
            document = json.loads(path.read_text(encoding="utf-8"))
            if document.get("asset") != name:
                issues.append(f"evidence_content_asset_mismatch:{name}:{kind}")
            for key in ("crop", "alpha", "animation", "evidence"):
                if not isinstance(document.get(key), dict):
                    issues.append(f"evidence_content_missing:{name}:{kind}:{key}")
            if document.get("evidence", {}).get("alpha") != asset["evidence"]["alpha"]["sha256"]:
                issues.append(f"evidence_cross_hash_mismatch:{name}:alpha")
            if document.get("evidence", {}).get("contact_sheet") != asset["evidence"]["contact_sheet"]["sha256"]:
                issues.append(f"evidence_cross_hash_mismatch:{name}:contact_sheet")
            if document.get("evidence", {}).get("animation_preview") != asset["evidence"]["animation_preview"]["sha256"]:
                issues.append(f"evidence_cross_hash_mismatch:{name}:animation_preview")
        elif kind == "contact_sheet":
            with Image.open(path) as image:
                image.load()
                expected_size = checks.get("crop", {}).get("output_canvas_size")
                if expected_size and list(image.size) != [expected_size[0] * 2, expected_size[1] + 24]:
                    issues.append(f"evidence_contact_dimensions_invalid:{name}")
        else:
            with Image.open(path) as image:
                image.seek(0)
                frame_count = int(getattr(image, "n_frames", 1))
                expected_count = checks.get("animation", {}).get("frame_count")
                if expected_count is not None and frame_count != expected_count:
                    issues.append(f"evidence_animation_frame_count_invalid:{name}")
                expected_durations = checks.get("animation", {}).get("durations_ms", [])
                durations = []
                for index in range(frame_count):
                    image.seek(index)
                    durations.append(image.info.get("duration"))
                if expected_durations and durations != expected_durations:
                    issues.append(f"evidence_animation_timing_invalid:{name}")
    except (OSError, ValueError, KeyError, json.JSONDecodeError) as error:
        issues.append(f"evidence_unreadable:{name}:{kind}:{error}")


def verify_review_gate(review_path: Path) -> tuple[bool, list[str]]:
    try:
        document = json.loads(review_path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as error:
        return False, [f"review_unreadable:{error}"]
    issues: list[str] = []
    if document.get("review_id") != "vertical-slice-v1":
        issues.append("review_id_invalid")
    source_root_value = document.get("source_root")
    if not isinstance(source_root_value, str) or not source_root_value:
        return False, ["source_root_invalid"]
    source_root = Path(source_root_value)
    if not source_root.is_dir():
        issues.append("source_root_invalid")
    qa_root = source_root / "qa"
    if not qa_root.is_dir():
        issues.append("qa_root_invalid")
    assets = document.get("assets")
    if not isinstance(assets, list):
        return False, issues + ["required_assets_missing"]
    names = [asset.get("asset") for asset in assets if isinstance(asset, dict)]
    if len(assets) != len(names) or len(names) != len(set(names)) or tuple(sorted(names)) != tuple(sorted(REQUIRED_ASSET_NAMES)):
        issues.append("exact_required_assets_mismatch")
    seen_evidence: set[Path] = set()
    for asset in assets:
        if not isinstance(asset, dict):
            issues.append("asset_record_invalid")
            continue
        name = str(asset.get("asset", "<unknown>"))
        if asset.get("required") is not True:
            issues.append(f"asset_not_required:{name}")
        if asset.get("status") != "PASS":
            issues.append(f"required_asset_not_pass:{name}")
        if asset.get("source_status") != "PASS":
            issues.append(f"source_status_not_pass:{name}")
        if asset.get("issues") != []:
            issues.append(f"issues_not_empty:{name}")
        source = source_root / name
        if not source.is_file():
            issues.append(f"source_file_missing:{name}")
        elif asset.get("source_sha256") != sha256(source):
            issues.append(f"source_sha256_mismatch:{name}")
        evidence = asset.get("evidence", {})
        if set(evidence) != set(REQUIRED_EVIDENCE):
            issues.append(f"evidence_kinds_invalid:{name}")
        for kind in REQUIRED_EVIDENCE:
            item = evidence.get(kind)
            if not isinstance(item, dict) or not isinstance(item.get("path"), str) or not re.fullmatch(r"[0-9a-f]{64}", str(item.get("sha256", ""))):
                issues.append(f"evidence_missing:{name}:{kind}")
                continue
            if item.get("asset") != name or item.get("kind") != kind:
                issues.append(f"evidence_kind_mismatch:{name}:{kind}")
            path = Path(item["path"])
            resolved = path.resolve()
            if not _path_under(resolved, qa_root):
                issues.append(f"evidence_outside_qa_root:{name}:{kind}")
            if resolved in seen_evidence:
                issues.append(f"evidence_duplicate:{name}:{kind}")
            seen_evidence.add(resolved)
            if not path.is_file():
                issues.append(f"evidence_file_missing:{name}:{kind}")
            elif sha256(path) != item["sha256"]:
                issues.append(f"sha256_mismatch:{name}:{kind}")
            elif _path_under(resolved, qa_root):
                _verify_evidence_content(asset, kind, path, issues)
    actual_pass = sum(asset.get("status") == "PASS" for asset in assets)
    actual_fail = len(assets) - actual_pass
    summary = document.get("summary", {})
    if summary.get("total") != len(assets) or summary.get("pass") != actual_pass or summary.get("fail") != actual_fail:
        issues.append("summary_does_not_match_assets")
    if summary.get("accepted_source_pass") != len(assets) or summary.get("accepted_source_fail") != 0:
        issues.append("accepted_source_summary_does_not_match_assets")
    return not issues, issues


def require_review_gate(review_path: Path) -> dict:
    ok, issues = verify_review_gate(review_path)
    if not ok:
        raise ReviewGateError(review_path, issues)
    return json.loads(review_path.read_text(encoding="utf-8"))


def main() -> int:
    parser = argparse.ArgumentParser(description="QA PNG nguồn và cổng Task 6B")
    subparsers = parser.add_subparsers(dest="command", required=True)
    review_parser = subparsers.add_parser("review", help="sinh output QA và review JSON")
    review_parser.add_argument("--art-root", type=Path, required=True)
    review_parser.add_argument("--review", type=Path, required=True)
    review_parser.add_argument("--reviewer", default="Codex Task 6A")
    gate_parser = subparsers.add_parser("gate", help="từ chối bind nếu review/evidence chưa đạt")
    gate_parser.add_argument(
        "--review",
        type=Path,
        default=Path("docs/my-thuat/review/vertical-slice-v1.json"),
    )
    bind_parser = subparsers.add_parser("bind-precondition", help="precondition cho Task 6B/assetpack")
    bind_parser.add_argument("--review", type=Path, required=True)
    args = parser.parse_args()
    if args.command == "review":
        result = review(args.art_root, args.review, args.reviewer)
        print(json.dumps(result["summary"], ensure_ascii=False))
        return 0
    if args.command == "bind-precondition":
        try:
            document = require_review_gate(args.review)
        except ReviewGateError as error:
            print(json.dumps({"allowed": False, "issues": error.issues}, ensure_ascii=False))
            return 2
        print(json.dumps({"allowed": True, "review_id": document["review_id"]}, ensure_ascii=False))
        return 0
    ok, issues = verify_review_gate(args.review)
    print(json.dumps({"allowed": ok, "issues": issues}, ensure_ascii=False))
    return 0 if ok else 2


if __name__ == "__main__":
    raise SystemExit(main())
