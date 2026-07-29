"""Sinh các PNG dẫn xuất từ đúng bộ nguồn Task 6A đã review.

Generator này không vẽ nội dung mới và không ghi vào thư mục nguồn. Nó dùng
alpha evidence đã được review, crop rect trong mapping và resize cố định để
build có thể tái lập; mọi hash được tính từ byte thật sau khi ghi.
"""

from __future__ import annotations

import hashlib
import json
from pathlib import Path
import shutil
import tempfile
from typing import Any

from PIL import Image

try:
    from .artqa import extract_alpha, normalize_cells, validate_alpha_output
except ImportError:  # direct CLI execution
    import sys
    sys.path.insert(0, str(Path(__file__).resolve().parents[2]))
    from tools.artqa.artqa import extract_alpha, normalize_cells, validate_alpha_output


class GenerationError(ValueError):
    """Lỗi có mã ổn định cho cổng art và mutation test."""


ZOOMS = ("x1", "x2", "x3", "x4")
ZOOM_FACTORS = {"x1": 0.25, "x2": 0.5, "x3": 0.75, "x4": 1.0}
EVIDENCE_KINDS = ("alpha", "crop_report", "contact_sheet", "animation_preview")


def _sha(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def _fail(code: str, detail: str) -> None:
    raise GenerationError(f"{code}: {detail}")


def _read_json(path: Path) -> dict[str, Any]:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as error:
        _fail("ASSET-ART-001", f"JSON lỗi: {path}: {error}")
    if not isinstance(value, dict):
        _fail("ASSET-ART-001", f"JSON phải là object: {path}")
    return value


def _resolve_evidence_path(source_root: Path, value: str) -> Path:
    candidate = Path(value)
    if not candidate.is_absolute():
        candidate = source_root / candidate
    try:
        resolved = candidate.resolve()
        resolved.relative_to((source_root / "qa").resolve())
    except ValueError:
        _fail("ASSET-ART-001", f"evidence nằm ngoài qa: {value}")
    return resolved


def _reviewed_assets(review: dict[str, Any], source_root: Path) -> dict[str, dict[str, Any]]:
    assets = review.get("assets")
    if not isinstance(assets, list):
        _fail("ASSET-ART-002", "review thiếu assets[]")
    expected = {path.name for path in source_root.glob("*.png")}
    names = {item.get("asset") for item in assets if isinstance(item, dict)}
    if names != expected or len(names) != 20 or "iruka_portrait_chroma.png" not in names:
        _fail("ASSET-ART-002", "review phải khớp đúng 20 PNG nguồn, gồm Iruka portrait")
    records: dict[str, dict[str, Any]] = {}
    for item in assets:
        if not isinstance(item, dict) or not isinstance(item.get("asset"), str):
            _fail("ASSET-ART-002", "record nguồn malformed")
        name = item["asset"]
        if item.get("status") != "PASS" or item.get("source_status") != "PASS":
            _fail("ASSET-ART-002", f"nguồn chưa PASS: {name}")
        source = source_root / name
        if not source.is_file() or _sha(source) != item.get("source_sha256"):
            _fail("ASSET-ART-001", f"source SHA lệch hoặc thiếu: {name}")
        try:
            with Image.open(source) as image:
                source_size = list(image.size)
        except (OSError, ValueError) as error:
            _fail("ASSET-ART-001", f"source không đọc được: {name}: {error}")
        crop = item.get("checks", {}).get("crop", {})
        if source_size != crop.get("source_size"):
            _fail("ASSET-ART-001", f"source dimensions lệch review: {name}")
        evidence = item.get("evidence")
        if not isinstance(evidence, dict) or set(evidence) != set(EVIDENCE_KINDS):
            _fail("ASSET-ART-001", f"evidence thiếu closure: {name}")
        for kind in EVIDENCE_KINDS:
            record = evidence[kind]
            if not isinstance(record, dict) or not isinstance(record.get("path"), str):
                _fail("ASSET-ART-001", f"evidence malformed: {name}/{kind}")
            evidence_path = _resolve_evidence_path(source_root, record["path"])
            if not evidence_path.is_file() or _sha(evidence_path) != record.get("sha256"):
                _fail("ASSET-ART-001", f"evidence SHA lệch hoặc thiếu: {name}/{kind}")
        records[name] = item
    return records


def _crop_rect(sprite: dict[str, Any], source_size: list[int]) -> list[int]:
    crop = sprite.get("crop")
    if not isinstance(crop, dict) or not isinstance(crop.get("source_rect"), list):
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} thiếu explicit crop")
    rect = crop["source_rect"]
    if len(rect) != 4 or any(not isinstance(value, int) or isinstance(value, bool) for value in rect):
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} crop malformed")
    if rect[0] < 0 or rect[1] < 0:
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} negative crop")
    if rect[2] <= 0 or rect[3] <= 0:
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} empty crop")
    if (
        len(source_size) != 2
        or rect[0] + rect[2] > source_size[0]
        or rect[1] + rect[3] > source_size[1]
    ):
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} crop outside source canvas")
    return rect


def _alpha_rect(
    sprite: dict[str, Any],
    rect: list[int],
    crop_report: dict[str, Any],
    alpha_size: tuple[int, int],
) -> list[int]:
    """Convert a source-sheet rect into the reviewed derived-alpha space."""
    source_size = crop_report.get("source_size")
    edges = crop_report.get("source_edges")
    boxes = crop_report.get("content_boxes")
    cell_size = crop_report.get("output_cell_size")
    if (
        not isinstance(source_size, list)
        or not isinstance(edges, list)
        or not isinstance(boxes, list)
        or not isinstance(cell_size, list)
        or len(source_size) != 2
        or len(edges) < 2
        or len(boxes) != len(edges) - 1
        or len(cell_size) != 2
        or any(not isinstance(value, int) or isinstance(value, bool) or value <= 0 for value in source_size)
        or any(not isinstance(value, int) or isinstance(value, bool) for value in edges)
        or edges[0] != 0
        or edges[-1] != source_size[0]
        or any(left >= right for left, right in zip(edges, edges[1:]))
        or any(edge < 0 or edge > source_size[0] for edge in edges)
        or any(not isinstance(value, int) or isinstance(value, bool) or value <= 0 for value in cell_size)
        or cell_size[0] * (len(edges) - 1) != alpha_size[0]
        or cell_size[1] != alpha_size[1]
    ):
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} crop coordinate evidence malformed")

    if (
        not isinstance(rect, list)
        or len(rect) != 4
        or any(not isinstance(value, int) or isinstance(value, bool) for value in rect)
    ):
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} crop coordinate malformed")
    source_left, source_top, width, height = rect
    source_right, source_bottom = source_left + width, source_top + height
    frame_index = None
    for index, (left, right) in enumerate(zip(edges, edges[1:])):
        if source_left >= left and source_right <= right:
            frame_index = index
            break
    if frame_index is None:
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} crop crosses source frame boundary")

    box = boxes[frame_index]
    if not all(
        isinstance(item, list)
        and len(item) == 4
        and all(isinstance(value, int) and not isinstance(value, bool) for value in item)
        for item in boxes
    ):
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} alpha content box malformed")
    box_left, box_top, box_right, box_bottom = box
    frame_width = box_right - box_left
    frame_height = box_bottom - box_top
    frame_source_width = edges[frame_index + 1] - edges[frame_index]
    if (
        box_left < 0
        or box_top < 0
        or box_right <= box_left
        or box_bottom <= box_top
        or box_right > frame_source_width
        or box_bottom > source_size[1]
        or frame_width > cell_size[0]
        or frame_height > cell_size[1]
    ):
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} alpha content box outside source cell")
    normalize_x = (cell_size[0] - frame_width) // 2
    normalize_y = cell_size[1] - frame_height
    alpha_rect = [
        frame_index * cell_size[0] + normalize_x + source_left - edges[frame_index] - box_left,
        normalize_y + source_top - box_top,
        width,
        height,
    ]
    if (
        alpha_rect[0] < 0
        or alpha_rect[1] < 0
        or alpha_rect[2] <= 0
        or alpha_rect[3] <= 0
        or alpha_rect[0] + alpha_rect[2] > alpha_size[0]
        or alpha_rect[1] + alpha_rect[3] > alpha_size[1]
    ):
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} crop outside alpha evidence")
    return alpha_rect


def _save_zoom(image: Image.Image, path: Path, root: Path, factor: float) -> dict[str, Any]:
    width = max(1, round(image.width * factor))
    height = max(1, round(image.height * factor))
    resized = image if factor == 1.0 else image.resize((width, height), Image.Resampling.LANCZOS)
    resized = resized.convert("RGBA")
    path.parent.mkdir(parents=True, exist_ok=True)
    resized.save(path, format="PNG", optimize=False, compress_level=9)
    return {
        "path": path.relative_to(root).as_posix(),
        "sha256": _sha(path),
        "dimensions": [resized.width, resized.height],
    }


def _raw_source_rect(sprite: dict[str, Any], source_size: list[int]) -> list[int]:
    rect = sprite.get("source_rect")
    crop = sprite.get("crop")
    nested = crop.get("source_rect") if isinstance(crop, dict) else None
    if rect is None:
        rect = nested
    if not isinstance(rect, list):
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} source_rect missing")
    if nested is not None and nested != rect:
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} source_rect mismatch")
    if len(rect) != 4 or any(not isinstance(value, int) or isinstance(value, bool) for value in rect):
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} source_rect malformed")
    if rect[0] < 0 or rect[1] < 0:
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} source_rect negative")
    if rect[2] <= 0 or rect[3] <= 0:
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} source_rect empty")
    if len(source_size) != 2 or rect[0] + rect[2] > source_size[0] or rect[1] + rect[3] > source_size[1]:
        _fail("ASSET-ART-004", f"sprite {sprite.get('id')} source_rect outside raw source")
    return rect


def _derive_raw_sprite(
    sprite: dict[str, Any],
    reviewed: dict[str, dict[str, Any]],
    source_root: Path,
    temporary: Path,
) -> dict[str, Any]:
    sprite_id = sprite["id"]
    source_path = sprite.get("source_path")
    if not isinstance(source_path, str) or not source_path or Path(source_path).is_absolute():
        _fail("ASSET-ART-004", f"sprite {sprite_id} source_path invalid")
    source_file = (source_root / source_path).resolve()
    try:
        source_file.relative_to(source_root)
    except ValueError:
        _fail("ASSET-ART-004", f"sprite {sprite_id} source_path outside ArtRoot")
    if not source_file.is_file():
        _fail("ASSET-ART-004", f"sprite {sprite_id} raw source missing")
    record = reviewed.get(Path(source_path).name)
    if record is None:
        _fail("ASSET-ART-002", f"sprite source review missing: {source_path}")
    source_sha = _sha(source_file)
    if source_sha != record.get("source_sha256") or sprite.get("source_sha256") != source_sha:
        _fail("ASSET-ART-004", f"sprite {sprite_id} source SHA mismatch")

    with Image.open(source_file) as raw_image:
        raw_image.load()
        source_size = list(raw_image.size)
        declared_size = sprite.get("source_dimensions")
        if declared_size is not None and declared_size != source_size:
            _fail("ASSET-ART-004", f"sprite {sprite_id} raw source dimensions mismatch")
        rect = _raw_source_rect(sprite, source_size)
        raw_crop = raw_image.crop((rect[0], rect[1], rect[0] + rect[2], rect[1] + rect[3])).convert("RGBA")

    # Canonical order: crop raw source coordinates, then run the shared chroma/alpha
    # extraction and single-cell normalization on that crop. No full-sheet alpha
    # evidence coordinates are used here.
    extracted = extract_alpha(raw_crop)
    normalized, crop_report, crop_issues = normalize_cells(extracted, 1)
    alpha_issues, alpha_metrics = validate_alpha_output(normalized, 1)
    stem = str(sprite_id)
    sprite_root = temporary / "sprites" / stem
    raw_path = sprite_root / "raw-crop.png"
    alpha_path = sprite_root / "alpha.png"
    crop_path = sprite_root / "crop.json"
    raw_path.parent.mkdir(parents=True, exist_ok=True)
    raw_crop.save(raw_path, format="PNG", optimize=False, compress_level=9)
    normalized.save(alpha_path, format="PNG", optimize=False, compress_level=9)
    crop_document = {
        "sprite_id": sprite_id,
        "source_path": source_path,
        "source_sha256": source_sha,
        "source_dimensions": source_size,
        "source_rect": rect,
        "raw_crop_dimensions": list(raw_crop.size),
        "alpha_dimensions": list(normalized.size),
        "crop": crop_report,
        "alpha": alpha_metrics,
        "issues": crop_issues + alpha_issues,
    }
    crop_path.write_text(json.dumps(crop_document, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    zooms = {}
    for zoom in ZOOMS:
        zooms[zoom] = _save_zoom(normalized, sprite_root / zoom / "sprite.png", temporary, ZOOM_FACTORS[zoom])
    raw_crop_record = {
        "path": raw_path.relative_to(temporary).as_posix(),
        "sha256": _sha(raw_path),
        "dimensions": list(raw_crop.size),
    }
    alpha_record = {
        "path": alpha_path.relative_to(temporary).as_posix(),
        "sha256": _sha(alpha_path),
        "dimensions": list(normalized.size),
    }
    crop_record = {
        "path": crop_path.relative_to(temporary).as_posix(),
        "sha256": _sha(crop_path),
        "source_rect": rect,
        "raw_dimensions": list(raw_crop.size),
        "alpha_dimensions": list(normalized.size),
    }
    return {
        "source_path": source_path,
        "source_sha256": source_sha,
        "source_dimensions": source_size,
        "source_rect": rect,
        "raw_crop": raw_crop_record,
        "alpha": alpha_record,
        "crop": crop_record,
        "zooms": zooms,
        "evidence": {
            "alpha": alpha_record,
            "crop": crop_record,
            "zooms": zooms,
        },
    }


def generate_reviewed_assets(
    review_path: Path,
    source_root: Path,
    output_root: Path,
) -> dict[str, Any]:
    """Sinh x1..x4 cho 20 asset-derived canvas đã được review.

    Đây là output provenance, chưa phải sprite ID binding. Alpha evidence đã
    chứa crop canvas canonical nên không được cắt lại theo tọa độ chưa review.
    """
    review = _read_json(Path(review_path))
    source_root = Path(source_root).resolve()
    output_root = Path(output_root).resolve()
    if output_root.exists():
        _fail("ASSET-ART-003", f"output đã tồn tại: {output_root}")
    if not source_root.is_dir():
        _fail("ASSET-ART-001", f"source root không tồn tại: {source_root}")
    reviewed = _reviewed_assets(review, source_root)
    parent = output_root.parent
    parent.mkdir(parents=True, exist_ok=True)
    temporary = Path(tempfile.mkdtemp(prefix=f".{output_root.name}.", dir=parent))
    try:
        index: dict[str, Any] = {
            "schema_version": 1,
            "review_id": review.get("review_id"),
            "source_root": str(source_root),
            "asset_count": len(reviewed),
            "sources": {},
            "assets": {},
        }
        for source_name, record in sorted(reviewed.items()):
            alpha_record = record["evidence"]["alpha"]
            alpha_path = _resolve_evidence_path(source_root, alpha_record["path"])
            with Image.open(alpha_path) as alpha_image:
                alpha = alpha_image.convert("RGBA")
            expected_dimensions = record["checks"]["crop"].get("output_canvas_size")
            if list(alpha.size) != expected_dimensions:
                _fail("ASSET-ART-001", f"alpha dimensions lệch crop review: {source_name}")
            stem = Path(source_name).stem
            zooms = {}
            for zoom in ZOOMS:
                path = temporary / "assets" / stem / zoom / "asset.png"
                zooms[zoom] = _save_zoom(alpha, path, temporary, ZOOM_FACTORS[zoom])
            index["sources"][source_name] = {
                "source_sha256": record["source_sha256"],
                "source_dimensions": record["checks"]["crop"]["source_size"],
                "alpha_evidence_path": alpha_record["path"],
                "alpha_evidence_sha256": alpha_record["sha256"],
                "crop_evidence_path": record["evidence"]["crop_report"]["path"],
                "crop_evidence_sha256": record["evidence"]["crop_report"]["sha256"],
            }
            index["assets"][source_name] = {
                "source_sha256": record["source_sha256"],
                "alpha_dimensions": list(alpha.size),
                "zooms": zooms,
            }
        (temporary / "generated-index.json").write_text(
            json.dumps(index, ensure_ascii=False, indent=2, sort_keys=True) + "\n",
            encoding="utf-8",
        )
        temporary.replace(output_root)
        return index
    except Exception:
        shutil.rmtree(temporary, ignore_errors=True)
        raise


def generate_vertical_slice(
    mapping_path: Path,
    review_path: Path,
    source_root: Path,
    output_root: Path,
) -> dict[str, Any]:
    """Kiểm tra mapping có thể bind vào alpha evidence đã review hay chưa.

    Cổng này cố ý chưa sinh sprite khi mapping crop không chứng minh được từ
    alpha canvas. Nhờ vậy đường dẫn planned hoặc crop sheet sai không biến
    thành output giả chỉ vì Pillow vẫn có thể cắt một vùng nào đó.
    """
    mapping = _read_json(Path(mapping_path))
    review = _read_json(Path(review_path))
    source_root = Path(source_root).resolve()
    output_root = Path(output_root).resolve()
    if output_root.exists():
        _fail("ASSET-ART-003", f"output đã tồn tại: {output_root}")
    if not source_root.is_dir():
        _fail("ASSET-ART-001", f"source root không tồn tại: {source_root}")
    reviewed = _reviewed_assets(review, source_root)
    sprites = mapping.get("sprites")
    if not isinstance(sprites, list) or not sprites:
        _fail("ASSET-ART-002", "mapping thiếu sprites[]")

    parent = output_root.parent
    parent.mkdir(parents=True, exist_ok=True)
    temporary = Path(tempfile.mkdtemp(prefix=f".{output_root.name}.", dir=parent))
    try:
        index: dict[str, Any] = {
            "schema_version": 1,
            "sprite_count": len(sprites),
            "mapping_sha256": _sha(Path(mapping_path)),
            "review_id": review.get("review_id"),
            "sprites": {},
        }
        for sprite in sprites:
            if isinstance(sprite, dict) and isinstance(sprite.get("id"), int):
                index["sprites"][str(sprite["id"])] = _derive_raw_sprite(
                    sprite, reviewed, source_root, temporary
                )
                continue
            if not isinstance(sprite, dict) or not isinstance(sprite.get("id"), int):
                _fail("ASSET-ART-004", "sprite record malformed")
            source_path = sprite.get("source_path")
            source_name = Path(source_path).name if isinstance(source_path, str) else ""
            record = reviewed.get(source_name)
            if record is None:
                _fail("ASSET-ART-002", f"sprite chưa có source review: {source_name}")
            crop_report = record["checks"]["crop"]
            source_size = crop_report["source_size"]
            if sprite.get("source_dimensions") not in (None, source_size):
                _fail("ASSET-ART-004", f"sprite {sprite['id']} source canvas mismatch")
            if sprite.get("source_sha256") != record.get("source_sha256"):
                _fail("ASSET-ART-004", f"sprite {sprite['id']} source SHA mismatch")
            rect = _crop_rect(sprite, source_size)
            alpha_path = _resolve_evidence_path(source_root, record["evidence"]["alpha"]["path"])
            with Image.open(alpha_path) as alpha_image:
                alpha = alpha_image.convert("RGBA")
            expected_alpha_size = crop_report.get("output_canvas_size")
            if list(alpha.size) != expected_alpha_size:
                _fail("ASSET-ART-004", f"sprite {sprite['id']} alpha canvas mismatch")
            alpha_rect = _alpha_rect(sprite, rect, crop_report, alpha.size)
            cropped = alpha.crop((
                alpha_rect[0],
                alpha_rect[1],
                alpha_rect[0] + alpha_rect[2],
                alpha_rect[1] + alpha_rect[3],
            ))
            stem = str(sprite["id"])
            zooms = {}
            for zoom in ZOOMS:
                path = temporary / "sprites" / stem / zoom / "sprite.png"
                zooms[zoom] = _save_zoom(cropped, path, temporary, ZOOM_FACTORS[zoom])
            index["sprites"][str(sprite["id"])] = {
                "source_path": source_path,
                "source_sha256": record["source_sha256"],
                "source_rect": rect,
                "alpha_evidence_path": record["evidence"]["alpha"]["path"],
                "alpha_evidence_sha256": record["evidence"]["alpha"]["sha256"],
                "alpha_rect": alpha_rect,
                "zooms": zooms,
            }
        (temporary / "generated-index.json").write_text(
            json.dumps(index, ensure_ascii=False, indent=2, sort_keys=True) + "\n",
            encoding="utf-8",
        )
        temporary.replace(output_root)
        return index
    except Exception:
        shutil.rmtree(temporary, ignore_errors=True)
        raise


def main() -> int:
    import argparse
    import sys

    parser = argparse.ArgumentParser(description="Sinh PNG zoom deterministic cho vertical slice")
    parser.add_argument("--mapping", type=Path, required=True)
    parser.add_argument("--review", type=Path, required=True)
    parser.add_argument("--source-root", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    try:
        result = generate_vertical_slice(args.mapping, args.review, args.source_root, args.output)
    except (GenerationError, OSError, json.JSONDecodeError) as error:
        print(error, file=sys.stderr)
        return 1
    print(json.dumps({"sprite_count": result["sprite_count"], "output": str(args.output)}))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
