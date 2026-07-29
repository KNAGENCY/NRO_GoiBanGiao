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


def _crop_rect(sprite: dict[str, Any], review: dict[str, Any]) -> list[int]:
    crop = sprite.get("crop")
    if isinstance(crop, dict) and isinstance(crop.get("source_rect"), list):
        rect = crop["source_rect"]
    else:
        boxes = review.get("checks", {}).get("crop", {}).get("content_boxes", [])
        if not boxes:
            _fail("ASSET-ART-001", f"thiếu crop: {sprite.get('id')}")
        left, top, right, bottom = boxes[0]
        rect = [left, top, right - left, bottom - top]
    if len(rect) != 4 or any(not isinstance(value, int) or isinstance(value, bool) for value in rect):
        _fail("ASSET-ART-001", f"crop malformed: {sprite.get('id')}")
    if rect[2] <= 0 or rect[3] <= 0:
        _fail("ASSET-ART-001", f"crop rỗng: {sprite.get('id')}")
    return rect


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

    for sprite in sprites:
        source_name = Path(sprite.get("source_path", "")).name
        record = reviewed.get(source_name)
        if record is None:
            _fail("ASSET-ART-002", f"sprite chưa có source review: {source_name}")
        alpha_path = _resolve_evidence_path(source_root, record["evidence"]["alpha"]["path"])
        with Image.open(alpha_path) as alpha_image:
            alpha_size = alpha_image.size
        rect = _crop_rect(sprite, record)
        if rect[0] + rect[2] > alpha_size[0] or rect[1] + rect[3] > alpha_size[1]:
            _fail("ASSET-ART-004", f"mapping crop không nằm trong alpha evidence: sprite {sprite.get('id')}")
    _fail("ASSET-ART-004", "mapping đã kiểm tra nhưng chưa có sprite binding provenance")


def main() -> int:
    import argparse

    parser = argparse.ArgumentParser(description="Sinh PNG zoom deterministic cho vertical slice")
    parser.add_argument("--mapping", type=Path, required=True)
    parser.add_argument("--review", type=Path, required=True)
    parser.add_argument("--source-root", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    try:
        result = generate_vertical_slice(args.mapping, args.review, args.source_root, args.output)
    except (GenerationError, OSError, json.JSONDecodeError) as error:
        print(error)
        return 1
    print(json.dumps({"sprite_count": result["sprite_count"], "output": str(args.output)}))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
