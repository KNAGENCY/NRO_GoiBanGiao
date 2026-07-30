"""Task 6B: bind provenance đã duyệt vào manifest canonical.

Module này chỉ làm việc với metadata và đường dẫn đã có. Nó không tạo ảnh giả,
không tự sửa review, và không ghi output nếu cổng review không qua.
"""

from __future__ import annotations

import hashlib
import json
from pathlib import Path
from typing import Any

from tools.artqa.artqa import ReviewGateError, require_review_gate
from .build import BuildError, _schema, resolve_manifest_base


class BindError(ValueError):
    """Lỗi bind có mã ổn định để CI và mutation test bắt được."""


EVIDENCE_KINDS = ("alpha", "crop_report", "contact_sheet", "animation_preview")
_KNOWN_CLIENT_CF = {
    "player_idle": [0],
}


def _fail(code: str, detail: str) -> None:
    raise BindError(f"{code}: {detail}")


def _sha(value: Any) -> bool:
    return isinstance(value, str) and len(value) == 64 and all(
        char in "0123456789abcdef" for char in value
    )


def _validate_evidence(record: dict[str, Any]) -> None:
    evidence = record.get("evidence")
    if not isinstance(evidence, dict) or set(evidence) != set(EVIDENCE_KINDS):
        _fail("ASSET-EVIDENCE-012", f"{record.get('asset', '?')} thiếu derived evidence")
    for kind in EVIDENCE_KINDS:
        item = evidence[kind]
        if (
            not isinstance(item, dict)
            or not isinstance(item.get("path"), str)
            or not item["path"]
            or not _sha(item.get("sha256"))
        ):
            _fail("ASSET-EVIDENCE-012", f"{record.get('asset', '?')} evidence {kind} malformed")
    if record.get("evidence_path") != evidence["crop_report"].get("path"):
        _fail("ASSET-EVIDENCE-012", f"{record.get('asset', '?')} evidence_path mismatch")
    if record.get("evidence_sha256") != evidence["crop_report"].get("sha256"):
        _fail("ASSET-EVIDENCE-012", f"{record.get('asset', '?')} evidence_sha256 mismatch")


def _validate_geometry(record: dict[str, Any]) -> None:
    dimensions = record.get("source_dimensions")
    crop = record.get("crop")
    rect = crop.get("source_rect") if isinstance(crop, dict) else None
    if (
        not isinstance(dimensions, list)
        or len(dimensions) != 2
        or any(not isinstance(value, int) or isinstance(value, bool) or value <= 0 for value in dimensions)
        or not isinstance(crop, dict)
        or not isinstance(rect, list)
        or len(rect) != 4
        or any(not isinstance(value, int) or isinstance(value, bool) for value in rect)
        or rect[0] < 0
        or rect[1] < 0
        or rect[2] <= 0
        or rect[3] <= 0
        or rect[0] + rect[2] > dimensions[0]
        or rect[1] + rect[3] > dimensions[1]
    ):
        _fail("ASSET-CROP-007", f"{record.get('asset', '?')} crop/source dimensions invalid")
    anchor = record.get("anchor")
    if (
        not isinstance(anchor, dict)
        or any(
            not isinstance(anchor.get(key), int) or isinstance(anchor.get(key), bool)
            for key in ("x", "y")
        )
    ):
        _fail("ASSET-ANCHOR-009", f"{record.get('asset', '?')} anchor invalid")
    for field in ("overlap", "layer"):
        value = record.get(field)
        if field == "overlap" and (
            not isinstance(value, list) or len(value) != 2
            or any(not isinstance(item, int) or isinstance(item, bool) or item < 0 for item in value)
        ):
            _fail("ASSET-ANCHOR-009", f"{record.get('asset', '?')} overlap invalid")
        if field == "layer" and (not isinstance(value, int) or isinstance(value, bool) or value < 0):
            _fail("ASSET-ANCHOR-009", f"{record.get('asset', '?')} layer invalid")
    frames = crop.get("frames")
    if not isinstance(frames, list) or not frames:
        _fail("ASSET-ANCHOR-009", f"{record.get('asset', '?')} crop frame geometry missing")
    for frame in frames:
        if not isinstance(frame, dict) or not isinstance(frame.get("anchor"), dict) or "overlap" not in frame or "layer" not in frame:
            _fail("ASSET-ANCHOR-009", f"{record.get('asset', '?')} frame anchor geometry missing")


def _validate_animation(record: dict[str, Any], occupied: set[tuple[str, int]]) -> None:
    action = record.get("action_strip")
    slot = record.get("slot_index")
    client_cf = record.get("client_cf")
    if not isinstance(action, str) or not action or not isinstance(slot, int) or isinstance(slot, bool) or slot < 0:
        _fail("ASSET-ANIMATION-008", f"{record.get('asset', '?')} action/slot invalid")
    key = (action, slot)
    if key in occupied:
        _fail("ASSET-ANIMATION-008", f"duplicate action slot: {action}[{slot}]")
    occupied.add(key)
    if (
        not isinstance(client_cf, list)
        or not client_cf
        or any(not isinstance(value, int) or isinstance(value, bool) for value in client_cf)
        or len(client_cf) != len(set(client_cf))
    ):
        _fail("ASSET-ANIMATION-008", f"{record.get('asset', '?')} client_cf invalid")
    expected = _KNOWN_CLIENT_CF.get(action)
    if expected is not None and client_cf != expected:
        _fail("ASSET-ANIMATION-008", f"{action} client_cf không đúng contract")


def _validate_runtime_map(runtime_map: dict[str, Any]) -> None:
    entries = runtime_map.get("entries") if isinstance(runtime_map, dict) else None
    if not isinstance(entries, list) or len(entries) != 8:
        _fail("ASSET-UI-007", "runtime key map phai co dung 8 binding")
    owners = [entry.get("source_asset") for entry in entries if isinstance(entry, dict)]
    if len(owners) != 8 or any(not isinstance(owner, str) or not owner for owner in owners) or len(set(owners)) != 8:
        _fail("ASSET-UI-007", "runtime key map source_asset phai exact 8 va duy nhat")


def validate_reviewed_sources(
    manifest: dict[str, Any],
    review_document: dict[str, Any],
    *,
    base_dir: Path | None = None,
) -> None:
    """Reject a new manifest source unless its exact basename is reviewed."""
    reviewed_by_name = {
        Path(item.get("asset", "")).name: item
        for item in review_document.get("assets", [])
        if isinstance(item, dict) and isinstance(item.get("asset"), str)
    }
    reviewed = set(reviewed_by_name)
    sources = {
        Path(sprite.get("source_path", "")).name
        for sprite in manifest.get("sprites", [])
        if isinstance(sprite, dict) and isinstance(sprite.get("source_path"), str)
    }
    if not sources or not sources.issubset(reviewed):
        _fail("ASSET-REVIEW-010", f"manifest source chua co trong exact reviewed set: {sorted(sources - reviewed)}")
    if base_dir is None:
        return
    base_dir = Path(base_dir).resolve()
    for sprite in manifest.get("sprites", []):
        relative = sprite.get("source_path")
        if not isinstance(relative, str) or not relative or Path(relative).is_absolute():
            _fail("ASSET-REVIEW-010", f"source path không canonical: {relative!r}")
        source = (base_dir / relative).resolve()
        try:
            source.relative_to(base_dir)
        except ValueError:
            _fail("ASSET-REVIEW-010", f"source path vượt build base: {relative}")
        if not source.is_file():
            _fail("ASSET-REVIEW-010", f"thiếu source đã review: {relative}")
        digest = hashlib.sha256(source.read_bytes()).hexdigest()
        reviewed_item = reviewed_by_name.get(source.name)
        if reviewed_item is None or digest != sprite.get("source_sha256") or digest != reviewed_item.get("source_sha256"):
            _fail("ASSET-REVIEW-010", f"source SHA không khớp review: {relative}")


def validate_bind_records(
    records: list[dict[str, Any]],
    *,
    review_document: dict[str, Any],
    runtime_map: dict[str, Any],
) -> None:
    """Kiểm tra các metadata bind trước khi manifest có thể được ghi."""
    if not isinstance(records, list) or not records:
        _fail("ASSET-EVIDENCE-012", "bind records rỗng")
    reviewed = {
        Path(item.get("asset", "")).name: item
        for item in review_document.get("assets", [])
        if isinstance(item, dict) and isinstance(item.get("asset"), str)
    }
    occupied: set[tuple[str, int]] = set()
    for record in records:
        if not isinstance(record, dict):
            _fail("ASSET-EVIDENCE-012", "bind record phải là object")
        _validate_evidence(record)
        review_sha = record.get("review_source_sha256")
        if review_sha is not None and review_sha != record.get("source_sha256"):
            _fail("ASSET-REVIEW-010", f"source SHA đổi sau review: {record.get('asset', '?')}")
        reviewed_item = reviewed.get(Path(record.get("asset", "")).name)
        if reviewed_item is not None and reviewed_item.get("source_sha256") != record.get("source_sha256"):
            _fail("ASSET-REVIEW-010", f"source SHA lệch review: {record.get('asset', '?')}")
        _validate_geometry(record)
        _validate_animation(record, occupied)
    _validate_runtime_map(runtime_map)


def validate_canonical_manifest(manifest: dict[str, Any]) -> None:
    """Áp đúng schema mà ``assetpack.build`` sẽ nhận, không nới lỏng contract."""
    try:
        _schema(manifest)
    except BuildError as error:
        _fail(str(error).split(":", 1)[0], str(error).split(":", 1)[1].strip())


def _read_json(path: Path) -> dict[str, Any]:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as error:
        _fail("ASSET-BIND-001", f"JSON lỗi {path}: {error}")
    if not isinstance(value, dict):
        _fail("ASSET-BIND-001", f"JSON phải là object: {path}")
    return value


def _review_records(document: dict[str, Any]) -> list[dict[str, Any]]:
    records = []
    for asset in document.get("assets", []):
        if not isinstance(asset, dict):
            continue
        checks = asset.get("checks", {})
        crop = checks.get("crop", {})
        boxes = crop.get("content_boxes")
        source_size = crop.get("source_size")
        edges = crop.get("source_edges")
        cell_size = crop.get("output_cell_size") or crop.get("output_canvas_size")
        if not isinstance(boxes, list) or not boxes:
            _fail("ASSET-CROP-007", f"{asset.get('asset', '?')} canonical crop report has no content box")
        if not isinstance(source_size, list) or not isinstance(edges, list) or not isinstance(cell_size, list):
            _fail("ASSET-CROP-007", f"{asset.get('asset', '?')} canonical crop report malformed")
        frames = []
        for index, box in enumerate(boxes):
            if not isinstance(box, list) or len(box) != 4 or not all(isinstance(value, int) for value in box):
                _fail("ASSET-CROP-007", f"{asset.get('asset', '?')} canonical content box malformed")
            x0, y0, x1, y1 = box
            if x1 <= x0 or y1 <= y0:
                _fail("ASSET-CROP-007", f"{asset.get('asset', '?')} canonical content box empty")
            frames.append({
                "slot_index": index,
                "source_rect": [x0, y0, x1 - x0, y1 - y0],
                "anchor": {"x": (x0 + x1) // 2, "y": y1},
                "overlap": [x0, max(0, cell_size[0] - x1)],
                "layer": index,
            })
        evidence = asset.get("evidence", {})
        if not isinstance(evidence, dict) or set(evidence) != set(EVIDENCE_KINDS):
            _fail("ASSET-EVIDENCE-012", f"{asset.get('asset', '?')} derived evidence missing")
        first = frames[0]
        record = {
            "asset": asset.get("asset"),
            "review_source_sha256": asset.get("source_sha256"),
            "source_sha256": asset.get("source_sha256"),
            "source_dimensions": source_size,
            "crop": {"source_rect": first["source_rect"], "output_dimensions": crop.get("output_canvas_size"), "frames": frames},
            "anchor": first["anchor"],
            "overlap": first["overlap"],
            "layer": first["layer"],
            "action_strip": Path(asset.get("asset", "")).stem,
            "slot_index": 0,
            "client_cf": [0],
            "dimensions": crop.get("output_canvas_size"),
            "animation": checks.get("animation", {}),
            "evidence": evidence,
            "evidence_path": evidence["crop_report"].get("path"),
            "evidence_sha256": evidence["crop_report"].get("sha256"),
        }
        records.append(record)
    return records


def bind_manifest(
    mapping: Path,
    review: Path,
    runtime_map: Path,
    *,
    art_root: Path,
    output: Path | None = None,
) -> dict[str, Any]:
    """Bind review metadata, only committing output after the gate and checks pass."""
    try:
        review_document = require_review_gate(Path(review))
    except ReviewGateError as error:
        _fail("ASSET-REVIEW-GATE-001", "; ".join(error.issues))
    except Exception as error:
        _fail("ASSET-REVIEW-GATE-001", str(error))
    manifest = _read_json(Path(mapping))
    runtime = _read_json(Path(runtime_map))
    if isinstance(manifest.get("build_inputs"), dict) and "source_base" in manifest["build_inputs"]:
        art_root = resolve_manifest_base(Path(mapping), manifest, base_dir=art_root)
    validate_reviewed_sources(manifest, review_document, base_dir=art_root)
    records = _review_records(review_document)
    validate_bind_records(records, review_document=review_document, runtime_map=runtime)
    manifest.setdefault("build_inputs", {})
    manifest["build_inputs"].setdefault("task6b", {})
    manifest["build_inputs"]["task6b"].update(
        {
            "review_id": review_document.get("review_id"),
            "review_sha256": hashlib.sha256(Path(review).read_bytes()).hexdigest(),
            "derived_asset_count": len(records),
            "derived_evidence": [record["evidence"] for record in records],
            "runtime_key_map": runtime,
        }
    )
    validate_canonical_manifest(manifest)
    if output is not None:
        output = Path(output)
        if output.exists():
            _fail("ASSET-BIND-002", f"output da ton tai: {output}")
        output.parent.mkdir(parents=True, exist_ok=True)
        temporary = output.with_name(f".{output.name}.tmp")
        temporary.write_text(json.dumps(manifest, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
        temporary.replace(output)
    return manifest
