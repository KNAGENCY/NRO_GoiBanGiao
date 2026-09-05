"""Build đầy đủ asset server/client sau các cổng schema, review và provenance.

Build không nhận một manifest tối thiểu rồi tự đoán phần còn thiếu. Mọi đầu vào
runtime phải được khai báo tường minh; nếu thiếu, build dừng trước khi tạo output.
"""

from __future__ import annotations

from io import BytesIO
import hashlib
import json
from datetime import datetime
from pathlib import Path
import shutil
import tempfile
from typing import Any

from PIL import Image

from tools.artqa.artqa import ReviewGateError, require_review_gate

from .animation import AnimationError, compose_action_sheet, validate_animation_closure
from .catalog import encode_catalog, validate_catalog
from .lock import lock_bytes
from .map import MapError, encode_collision, validate_collision_semantics, validate_tileset_baseline
from .mob import MobError, encode_mob, validate_mob_golden
from .part import build_part, parse_part
from .skill import encode_skill
from .tileset import TilesetError, apply_tile_set_overlay, validate_floor_closure
from .visual import (
    render_full_map,
    render_hud_golden,
    render_iruka_composite_contact_sheet,
    render_tile_contact_sheet,
)


class BuildError(ValueError):
    """Lỗi build có mã ổn định để canary/mutation bắt đúng nguyên nhân."""


REQUIRED_TOP_LEVEL = {
    "schema_version", "source_catalog_version", "versions", "phase",
    "sprites", "parts", "animations", "maps", "tiles", "backgrounds",
    "mobs", "skills", "catalog", "ui_bindings", "build_inputs",
}
ZOOMS = ("x1", "x2", "x3", "x4")
POC_IDS = {9000, 9001, 9002}
PRODUCTION_REVIEWED_SOURCE_COUNT = 20
PRODUCTION_REVIEWED_SOURCE_ROLE = "reviewed_vertical_slice_source"
TRUSTED_PRODUCTION_ARTIFACT_REVIEWER = "/root/task6b_visual_review6"


def _fail(code: str, detail: str) -> None:
    raise BuildError(f"{code}: {detail}")


def _int(value: Any, code: str, field: str) -> int:
    if not isinstance(value, int) or isinstance(value, bool):
        _fail(code, f"{field} phải là số nguyên, không phải bool")
    return value


def _path(base_dir: Path, value: str, code: str) -> Path:
    if not isinstance(value, str) or not value or Path(value).is_absolute():
        _fail(code, f"path không hợp lệ: {value!r}")
    result = (base_dir / value).resolve()
    try:
        result.relative_to(base_dir.resolve())
    except ValueError:
        _fail(code, f"path vượt ngoài thư mục manifest: {value}")
    return result


def resolve_manifest_base(
    manifest_path: Path,
    manifest: dict[str, Any],
    *,
    base_dir: Path | None = None,
) -> Path:
    """Resolve the canonical source base and prove repo-root-relative sources exist.

    Production manifests use ``source_base`` explicitly.  Their source paths are
    relative to the caller-provided repository root, never to the mapping file or
    the shared ArtRoot directory.  The existence check happens before build output
    staging so a short/nonexistent path cannot silently become an artifact.
    """
    inputs = manifest.get("build_inputs")
    source_base = inputs.get("source_base") if isinstance(inputs, dict) else None
    if not isinstance(source_base, str) or not source_base or Path(source_base).is_absolute():
        _fail("ASSET-REVIEW-010", f"manifest thiếu source_base canonical: {source_base!r}")
    root = Path(base_dir or Path.cwd()).resolve()
    resolved = (root / source_base).resolve()
    try:
        resolved.relative_to(root)
    except ValueError:
        _fail("ASSET-REVIEW-010", f"source_base vượt ngoài build base: {source_base}")
    if not resolved.is_dir():
        _fail("ASSET-REVIEW-010", f"source_base không tồn tại: {source_base}")
    for index, sprite in enumerate(manifest.get("sprites", [])):
        if not isinstance(sprite, dict) or not isinstance(sprite.get("source_path"), str):
            _fail("ASSET-REVIEW-010", f"sprite[{index}] thiếu source_path canonical")
        source = _path(resolved, sprite["source_path"], "ASSET-REVIEW-010")
        if not source.is_file():
            _fail("ASSET-REVIEW-010", f"thiếu source production: {sprite['source_path']}")
    return resolved


def _sha(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def _png_meta(data: bytes) -> dict[str, Any]:
    try:
        with Image.open(BytesIO(data)) as image:
            image.load()
            return {"width": image.width, "height": image.height, "mode": image.mode}
    except Exception as error:  # Pillow raises several format-specific errors.
        _fail("ASSET-IMAGE-003", f"PNG không đọc được: {error}")


def _load_png(base_dir: Path, relative: str) -> Image.Image:
    source = _path(base_dir, relative, "ASSET-BUILD-013")
    with Image.open(source) as image:
        return image.convert("RGBA")


def _entries(review: Any) -> list[dict[str, Any]]:
    if isinstance(review, dict):
        assets = review.get("assets")
        if assets is None and "status" in review:
            assets = [review]
    else:
        assets = review
    if not isinstance(assets, list) or any(not isinstance(item, dict) for item in assets):
        _fail("ASSET-REVIEW-010", "review phải có assets[] là object")
    return assets


def _review_documents(review_path: Path) -> list[tuple[Path, Any]]:
    paths = sorted(review_path.glob("*.json")) if review_path.is_dir() else [review_path]
    if not paths:
        _fail("ASSET-REVIEW-010", "không có review JSON")
    documents = []
    for path in paths:
        try:
            documents.append((path, json.loads(path.read_text(encoding="utf-8"))))
        except (OSError, json.JSONDecodeError) as error:
            _fail("ASSET-REVIEW-010", f"review JSON lỗi: {error}")
    return documents


def _review_records(review_path: Path) -> tuple[list[dict[str, Any]], list[dict[str, str]]]:
    entries: list[dict[str, Any]] = []
    evidence: list[dict[str, str]] = []
    for document, payload in _review_documents(review_path):
        for entry in _entries(payload):
            if entry.get("status") != "PASS":
                _fail("ASSET-REVIEW-010", "source chưa được review PASS")
            source_sha = entry.get("source_sha256")
            evidence_path = entry.get("evidence_path")
            evidence_sha = entry.get("evidence_sha256")
            if not isinstance(source_sha, str) or len(source_sha) != 64 or not isinstance(evidence_path, str):
                _fail("ASSET-REVIEW-010", "thiếu source_sha256/evidence_path")
            if not isinstance(evidence_sha, str) or len(evidence_sha) != 64:
                _fail("ASSET-REVIEW-010", "thiếu evidence_sha256")
            evidence_file = Path(evidence_path)
            if not evidence_file.is_absolute():
                evidence_file = document.parent / evidence_file
            if not evidence_file.is_file() or _sha(evidence_file) != evidence_sha:
                _fail("ASSET-REVIEW-010", f"evidence SHA không khớp: {evidence_path}")
            entries.append(entry)
            evidence.append({"path": str(evidence_file), "sha256": evidence_sha})
    return entries, evidence


def _review_records(review_path: Path) -> tuple[list[dict[str, Any]], list[dict[str, str]]]:
    """Adapt artqa's four nested evidence records to the build contract."""
    entries: list[dict[str, Any]] = []
    evidence: list[dict[str, str]] = []
    for document, payload in _review_documents(review_path):
        for raw in _entries(payload):
            if raw.get("status") != "PASS":
                _fail("ASSET-REVIEW-010", "source chua duoc review PASS")
            source_sha = raw.get("source_sha256")
            nested = raw.get("evidence")
            if not isinstance(source_sha, str) or len(source_sha) != 64:
                _fail("ASSET-REVIEW-010", "review thieu source provenance")
            if not isinstance(nested, dict):
                path = raw.get("evidence_path")
                sha = raw.get("evidence_sha256")
                if not isinstance(path, str) or not isinstance(sha, str):
                    _fail("ASSET-REVIEW-010", "review thieu source/evidence provenance")
                evidence_file = Path(path)
                if not evidence_file.is_absolute():
                    evidence_file = document.parent / evidence_file
                if not evidence_file.is_file() or _sha(evidence_file) != sha:
                    _fail("ASSET-REVIEW-010", f"evidence SHA khong khop: {path}")
                normalized = dict(raw)
                normalized["evidence_path"] = str(evidence_file)
                normalized["evidence_sha256"] = sha
                entries.append(normalized)
                evidence.append({"path": str(evidence_file), "sha256": sha})
                continue
            if set(nested) != {"alpha", "crop_report", "contact_sheet", "animation_preview"}:
                _fail("ASSET-REVIEW-010", "review thieu dung 4 derived evidence")
            normalized = dict(raw)
            crop_evidence = nested["crop_report"]
            if not isinstance(crop_evidence, dict):
                _fail("ASSET-REVIEW-010", "crop evidence malformed")
            normalized["evidence_path"] = crop_evidence.get("path")
            normalized["evidence_sha256"] = crop_evidence.get("sha256")
            for kind, item in nested.items():
                if not isinstance(item, dict) or not isinstance(item.get("path"), str) or not isinstance(item.get("sha256"), str):
                    _fail("ASSET-REVIEW-010", f"evidence {kind} malformed")
                path = Path(item["path"])
                if not path.is_absolute():
                    path = document.parent / path
                if not path.is_file() or _sha(path) != item["sha256"]:
                    _fail("ASSET-REVIEW-010", f"evidence SHA khong khop: {item['path']}")
                evidence.append({"path": str(path), "sha256": item["sha256"]})
            if not isinstance(normalized["evidence_path"], str) or not isinstance(normalized["evidence_sha256"], str):
                _fail("ASSET-REVIEW-010", "review thieu evidence_path/evidence_sha256")
            entries.append(normalized)
    return entries, evidence


def review_digest(review_path: Path) -> str:
    digest = hashlib.sha256()
    for path, _ in _review_documents(review_path):
        digest.update(path.name.encode("utf-8"))
        digest.update(path.read_bytes())
    return digest.hexdigest()


def _is_sha256(value: Any) -> bool:
    return (
        isinstance(value, str)
        and len(value) == 64
        and all(character in "0123456789abcdef" for character in value)
    )


def _validate_production_review_provenance(
    review_path: Path,
    entries: list[dict[str, Any]],
    manifest: dict[str, Any],
    base_dir: Path,
) -> dict[str, dict[str, Any]]:
    """Buộc manifest khai đúng tuyệt đối 20 nguồn PASS của review canonical."""

    documents = _review_documents(review_path)
    if len(documents) != 1:
        _fail(
            "ASSET-PROVENANCE-014",
            "production review phải là đúng một JSON canonical",
        )
    _, payload = documents[0]
    source_root_value = payload.get("source_root") if isinstance(payload, dict) else None
    if not isinstance(source_root_value, str) or not source_root_value:
        _fail("ASSET-PROVENANCE-014", "production review thiếu source_root")
    source_root = Path(source_root_value).resolve()
    canonical_base = base_dir.resolve()
    try:
        source_root.relative_to(canonical_base)
    except ValueError:
        _fail(
            "ASSET-PROVENANCE-014",
            "production review source_root vượt ngoài source_base",
        )
    if (
        len(entries) != PRODUCTION_REVIEWED_SOURCE_COUNT
        or any(entry.get("status") != "PASS" for entry in entries)
    ):
        _fail(
            "ASSET-PROVENANCE-014",
            "production review phải có đúng 20 asset PASS",
        )

    expected: set[tuple[str, str]] = set()
    reviewed_assets: set[str] = set()
    for entry in entries:
        asset = entry.get("asset")
        digest = entry.get("source_sha256")
        if (
            not isinstance(asset, str)
            or not asset
            or Path(asset).name != asset
            or not _is_sha256(digest)
            or asset in reviewed_assets
        ):
            _fail(
                "ASSET-PROVENANCE-014",
                "production review có asset duplicate/malformed",
            )
        reviewed_assets.add(asset)
        source = (source_root / asset).resolve()
        if source.parent != source_root or not source.is_file() or _sha(source) != digest:
            _fail(
                "ASSET-PROVENANCE-014",
                f"production review source/SHA lệch: {asset}",
            )
        relative = source.relative_to(canonical_base).as_posix()
        expected.add((relative, digest))

    provenance = manifest.get("build_inputs", {}).get("source_provenance")
    if not isinstance(provenance, list) or len(provenance) != len(expected):
        _fail(
            "ASSET-PROVENANCE-014",
            "source_provenance phải có đúng 20 record duy nhất",
        )
    observed: set[tuple[str, str]] = set()
    provenance_by_path: dict[str, dict[str, Any]] = {}
    for record in provenance:
        if (
            not isinstance(record, dict)
            or set(record)
            != {"source_path", "source_sha256", "role"}
            or not isinstance(record.get("source_path"), str)
            or not record["source_path"]
            or Path(record["source_path"]).is_absolute()
            or not _is_sha256(record.get("source_sha256"))
            or record.get("role") != PRODUCTION_REVIEWED_SOURCE_ROLE
            or record["source_path"] in provenance_by_path
        ):
            _fail(
                "ASSET-PROVENANCE-014",
                "source_provenance duplicate/malformed",
            )
        pair = (record["source_path"], record["source_sha256"])
        if pair in observed:
            _fail("ASSET-PROVENANCE-014", "source_provenance duplicate")
        observed.add(pair)
        provenance_by_path[record["source_path"]] = record
    if observed != expected:
        _fail(
            "ASSET-PROVENANCE-014",
            "source_provenance thiếu/thừa hoặc SHA lệch review PASS",
        )
    return provenance_by_path


def _require_production_artifact_review(review_path: Path) -> dict[str, Any]:
    documents = _review_documents(review_path)
    if len(documents) != 1:
        _fail(
            "ASSET-ARTIFACT-REVIEW-016",
            "artifact review phải nằm trong đúng một JSON canonical",
        )
    _, payload = documents[0]
    artifact_review = (
        payload.get("production_artifact_review")
        if isinstance(payload, dict)
        else None
    )
    if not isinstance(artifact_review, dict):
        _fail(
            "ASSET-ARTIFACT-REVIEW-016",
            "thiếu production_artifact_review",
        )
    return artifact_review


def production_closure_digest(manifest: dict[str, Any]) -> str:
    closure = manifest.get("build_inputs", {}).get(
        "production_closure_sha256"
    )
    if not isinstance(closure, dict):
        _fail(
            "ASSET-ARTIFACT-REVIEW-016",
            "thiếu production closure để bind artifact review",
        )
    canonical = json.dumps(
        closure,
        sort_keys=True,
        separators=(",", ":"),
    ).encode("utf-8")
    return hashlib.sha256(canonical).hexdigest()


def validate_production_artifact_review(
    review_path: Path,
    manifest: dict[str, Any],
    base_dir: Path,
) -> dict[str, dict[str, str]]:
    """Khóa evidence production và chữ ký reviewer trước khi tạo output."""

    artifact_review = _require_production_artifact_review(review_path)
    if set(artifact_review) != {
        "status",
        "reviewer",
        "timestamp",
        "production_closure_digest_sha256",
        "evidence",
    }:
        _fail(
            "ASSET-ARTIFACT-REVIEW-016",
            "production_artifact_review malformed",
        )
    reviewer = artifact_review.get("reviewer")
    if (
        artifact_review.get("status") != "PASS"
        or not isinstance(reviewer, str)
        or not reviewer.strip()
        or reviewer != TRUSTED_PRODUCTION_ARTIFACT_REVIEWER
    ):
        _fail(
            "ASSET-ARTIFACT-REVIEW-016",
            "artifact chưa có PASS của reviewer độc lập",
        )
    timestamp = artifact_review.get("timestamp")
    try:
        reviewed_at = datetime.fromisoformat(
            timestamp.replace("Z", "+00:00")
        )
    except (AttributeError, TypeError, ValueError):
        _fail("ASSET-ARTIFACT-REVIEW-016", "artifact review timestamp malformed")
    if reviewed_at.tzinfo is None or reviewed_at.utcoffset() is None:
        _fail(
            "ASSET-ARTIFACT-REVIEW-016",
            "artifact review timestamp phải có múi giờ",
        )
    if (
        artifact_review.get("production_closure_digest_sha256")
        != production_closure_digest(manifest)
    ):
        _fail(
            "ASSET-ARTIFACT-REVIEW-016",
            "artifact review không bind đúng production closure",
        )

    expected_kinds = {"hud_golden", "iruka_composite"}
    declared = manifest.get("build_inputs", {}).get("artifact_evidence")
    reviewed = artifact_review.get("evidence")
    if (
        not isinstance(declared, dict)
        or set(declared) != expected_kinds
        or not isinstance(reviewed, dict)
        or set(reviewed) != expected_kinds
    ):
        _fail(
            "ASSET-ARTIFACT-REVIEW-016",
            "artifact evidence phải đủ HUD golden và Iruka composite",
        )
    closure = manifest["build_inputs"].get("production_closure_sha256")
    validated: dict[str, dict[str, str]] = {}
    for kind in sorted(expected_kinds):
        record = declared[kind]
        verdict = reviewed[kind]
        if (
            not isinstance(record, dict)
            or set(record) != {"path", "sha256"}
            or not isinstance(verdict, dict)
            or set(verdict) != {"path", "sha256", "status"}
            or verdict.get("status") != "PASS"
            or verdict.get("path") != record.get("path")
            or verdict.get("sha256") != record.get("sha256")
            or not isinstance(record.get("path"), str)
            or not _is_sha256(record.get("sha256"))
            or not isinstance(closure, dict)
            or closure.get(record["path"]) != record["sha256"]
        ):
            _fail(
                "ASSET-ARTIFACT-REVIEW-016",
                f"artifact evidence {kind} thiếu/sai/không PASS",
            )
        path = _path(
            base_dir,
            record["path"],
            "ASSET-ARTIFACT-REVIEW-016",
        )
        if not path.is_file() or _sha(path) != record["sha256"]:
            _fail(
                "ASSET-ARTIFACT-REVIEW-016",
                f"artifact evidence file/SHA lệch: {kind}",
            )
        validated[kind] = {
            "path": str(path),
            "sha256": record["sha256"],
        }
    return validated


def _source_paths(manifest: dict[str, Any]) -> list[str]:
    paths: list[str] = []
    for sprite in manifest.get("sprites", []):
        paths.append(sprite.get("source_path"))
        paths.extend(sprite.get("zooms", {}).values())
    inputs = manifest.get("build_inputs", {})
    for key in ("legacy_part_path", "tile_set_info_path"):
        if inputs.get(key):
            paths.append(inputs[key])
    for item in inputs.get("tile_art", []):
        paths.extend(item.get("sources", {}).values())
    for item in inputs.get("background_sources", {}).values():
        paths.extend(item.values())
    for item in inputs.get("hud_sources", {}).values():
        paths.extend(item.values())
    for item in inputs.get("artifact_evidence", {}).values():
        if isinstance(item, dict) and isinstance(item.get("path"), str):
            paths.append(item["path"])
    mob = inputs.get("mob", {})
    if isinstance(mob, dict):
        for key in ("bytes_path", "data_path"):
            if mob.get(key):
                paths.append(mob[key])
        if isinstance(mob.get("img"), dict):
            paths.extend(mob["img"].values())
    for sheet in inputs.get("animation_sheets", []):
        paths.extend(cell.get("path") for cell in sheet.get("cells", []))
    runtime_map = inputs.get("runtime_key_map")
    if isinstance(runtime_map, str):
        paths.append(runtime_map)
    for binding in manifest.get("ui_bindings", []):
        if binding.get("source_path"):
            paths.append(binding["source_path"])
        if isinstance(binding.get("sources"), dict):
            paths.extend(binding["sources"].values())
    snapshots = inputs.get("runtime_source_snapshots")
    if isinstance(snapshots, dict):
        paths.extend(snapshots.values())
    for record in inputs.get("source_provenance", []):
        if isinstance(record, dict) and isinstance(record.get("source_path"), str):
            paths.append(record["source_path"])
    return sorted({path for path in paths if isinstance(path, str)})


def _validate_production_closure(manifest: dict[str, Any], base_dir: Path) -> None:
    """Khóa byte của mọi input production; review mỹ thuật chỉ khóa raw source."""
    inputs = manifest.get("build_inputs")
    closure = inputs.get("production_closure_sha256") if isinstance(inputs, dict) else None
    expected_paths = set(_source_paths(manifest))
    if (
        not isinstance(closure, dict)
        or set(closure) != expected_paths
        or any(
            not isinstance(digest, str)
            or len(digest) != 64
            or any(char not in "0123456789abcdef" for char in digest)
            for digest in closure.values()
        )
    ):
        _fail("ASSET-PROVENANCE-014", "production_closure_sha256 thiếu/thừa input")
    for relative, expected_sha in closure.items():
        source = _path(base_dir, relative, "ASSET-PROVENANCE-014")
        if not source.is_file() or _sha(source) != expected_sha:
            _fail("ASSET-PROVENANCE-014", f"SHA production input lệch: {relative}")


def _runtime_key_map(manifest: dict[str, Any], base_dir: Path) -> tuple[dict[str, Any], str | None, str]:
    reference = manifest.get("build_inputs", {}).get("runtime_key_map")
    if isinstance(reference, str):
        path = _path(base_dir, reference, "ASSET-UI-007")
        if not path.is_file():
            _fail("ASSET-UI-007", f"thiếu unity runtime key map: {reference}")
        try:
            payload = json.loads(path.read_text(encoding="utf-8"))
        except (OSError, json.JSONDecodeError) as error:
            _fail("ASSET-UI-007", f"runtime key map lỗi: {error}")
        return payload, reference, _sha(path)
    if isinstance(reference, dict):
        canonical = json.dumps(reference, ensure_ascii=False, sort_keys=True, separators=(",", ":")).encode("utf-8")
        return reference, None, hashlib.sha256(canonical).hexdigest()
    _fail("ASSET-UI-007", "build_inputs.runtime_key_map là bắt buộc")


_RUNTIME_BINDING_FINGERPRINT_FIELDS = (
    "logical_keys",
    "future_logical_key",
    "physical_paths",
    "call_sites",
    "consumer_sites",
)


def _runtime_binding_fingerprint(binding: dict[str, Any]) -> str:
    """Nhận diện contract runtime, không nhận diện theo tên file nguồn.

    `source_asset` chỉ là nhãn provenance. Nếu đưa nó vào khóa duy nhất thì hai
    binding cùng trỏ một runtime output có thể đổi nhãn để lọt qua collision gate.
    """
    return json.dumps(
        {key: binding.get(key) for key in _RUNTIME_BINDING_FINGERPRINT_FIELDS},
        sort_keys=True,
        separators=(",", ":"),
    )


def _validate_runtime_collision_dimensions(entries: list[dict[str, Any]]) -> None:
    """Reject duplicate runtime outputs independently of provenance or site metadata."""
    seen_logical: dict[str, str] = {}
    seen_future: dict[str, str] = {}
    seen_physical: dict[str, str] = {}
    for index, entry in enumerate(entries):
        owner = entry.get("source_asset", f"entry[{index}]")
        logical_keys = entry.get("logical_keys")
        if not isinstance(logical_keys, list) or any(not isinstance(key, str) or not key for key in logical_keys):
            _fail("ASSET-UI-007", f"{owner} logical_keys malformed")
        for key in logical_keys:
            previous = seen_logical.get(key)
            if previous is not None:
                _fail("ASSET-UI-007", f"logical key collision: {previous} vs {owner}: {key}")
            seen_logical[key] = owner

        future_key = entry.get("future_logical_key")
        if future_key is not None and (not isinstance(future_key, str) or not future_key):
            _fail("ASSET-UI-007", f"{owner} future_logical_key malformed")
        if future_key is not None:
            previous = seen_future.get(future_key)
            if previous is not None:
                _fail("ASSET-UI-007", f"future logical key collision: {previous} vs {owner}: {future_key}")
            seen_future[future_key] = owner

        physical_paths = entry.get("physical_paths")
        if not isinstance(physical_paths, dict) or set(physical_paths) != set(ZOOMS):
            _fail("ASSET-UI-007", f"{owner} physical_paths malformed")
        for zoom, paths in physical_paths.items():
            if not isinstance(paths, list) or not paths:
                _fail("ASSET-UI-007", f"{owner} physical path {zoom} malformed")
            for path in paths:
                normalized = path.replace("\\", "/") if isinstance(path, str) else ""
                if not normalized:
                    _fail("ASSET-UI-007", f"{owner} physical path malformed")
                previous = seen_physical.get(normalized)
                if previous is not None and previous != owner:
                    _fail("ASSET-UI-007", f"physical path collision: {previous} vs {owner}: {normalized}")
                seen_physical[normalized] = owner


def _validate_runtime_keys(manifest: dict[str, Any], base_dir: Path) -> tuple[str | None, str, dict[str, Any]]:
    payload, relative_path, map_sha = _runtime_key_map(manifest, base_dir)
    entries = payload.get("entries") if isinstance(payload, dict) else None
    bindings = manifest.get("ui_bindings")
    if not isinstance(entries, list) or len(entries) != 8 or not isinstance(bindings, list) or len(bindings) != 8:
        _fail("ASSET-UI-007", "runtime key map và manifest phải có đúng tám binding")
    for entry in entries:
        if not isinstance(entry, dict) or not isinstance(entry.get("source_asset"), str) or not entry["source_asset"]:
            _fail("ASSET-UI-007", "runtime entry thiếu source_asset hợp lệ")
    _validate_runtime_collision_dimensions(entries)
    by_asset = {entry["source_asset"]: entry for entry in entries}
    if len(by_asset) != 8:
        _fail("ASSET-UI-007", "runtime key map trùng/thiếu source_asset")
    binding_assets = [binding.get("source_asset") for binding in bindings if isinstance(binding, dict)]
    if any(not isinstance(asset, str) or not asset for asset in binding_assets):
        _fail("ASSET-UI-007", "manifest binding thiếu source_asset hợp lệ")
    if len(binding_assets) != 8 or len(set(binding_assets)) != 8 or set(binding_assets) != set(by_asset):
        _fail("ASSET-UI-007", "manifest phải bind exact tám source_asset duy nhất")
    entry_fingerprints = [_runtime_binding_fingerprint(entry) for entry in entries]
    if len(set(entry_fingerprints)) != 8:
        _fail("ASSET-UI-007", "runtime logical/source/physical/consumer binding bị duplicate")
    binding_fingerprints = []
    for binding in bindings:
        if not isinstance(binding, dict) or not isinstance(binding.get("source_asset"), str):
            _fail("ASSET-UI-007", "binding thiếu source_asset")
        entry = by_asset.get(binding["source_asset"])
        if entry is None:
            _fail("ASSET-UI-007", f"binding chưa được bind vào runtime map: {binding['source_asset']}")
        for field in ("mode", "logical_keys", "future_logical_key", "physical_paths"):
            if binding.get(field) != entry.get(field):
                _fail("ASSET-UI-007", f"{binding['source_asset']} lệch runtime contract ở {field}")
        for field in ("call_sites", "consumer_sites"):
            if field in binding and binding[field] != entry.get(field):
                _fail("ASSET-UI-007", f"{binding['source_asset']} lệch runtime contract ở {field}")
        mode = binding.get("mode")
        call_sites = entry.get("call_sites")
        consumer_sites = entry.get("consumer_sites")
        if not isinstance(call_sites, list) or not isinstance(consumer_sites, list):
            _fail("ASSET-UI-007", f"{binding['source_asset']} thiếu call-site/consumer-site")
        required_sites = call_sites if mode in {"direct", "slice_group"} else consumer_sites
        if not required_sites:
            _fail("ASSET-UI-007", f"{binding['source_asset']} thiếu site precondition")
        expressions = []
        for site in [*call_sites, *consumer_sites]:
            if not isinstance(site, dict) or not isinstance(site.get("file"), str) or not isinstance(site.get("expression"), str) or not isinstance(site.get("line"), int) or isinstance(site.get("line"), bool):
                _fail("ASSET-UI-007", f"{binding['source_asset']} site malformed")
            snapshots = manifest.get("build_inputs", {}).get("runtime_source_snapshots", {})
            snapshot = snapshots.get(site["file"]) if isinstance(snapshots, dict) else None
            source = _path(base_dir, snapshot or site["file"], "ASSET-UI-007")
            if not source.is_file() or not isinstance(site.get("file_sha256"), str) or _sha(source) != site["file_sha256"]:
                _fail("ASSET-UI-007", f"{binding['source_asset']} SHA site/precondition không khớp")
            expressions.append(site["expression"])
        if mode in {"direct", "slice_group"} and any(key not in "\n".join(expressions) for key in binding["logical_keys"]):
            _fail("ASSET-UI-007", f"{binding['source_asset']} logical key không có call-site tham chiếu")
        physical_paths = binding.get("physical_paths", {})
        if set(physical_paths) != set(ZOOMS):
            _fail("ASSET-UI-007", f"{binding['source_asset']} thiếu physical path x1..x4")
        for zoom, paths in physical_paths.items():
            if not isinstance(paths, list) or not paths:
                _fail("ASSET-UI-007", f"{binding['source_asset']} thiếu output {zoom}")
            prefix = f"Assets/Resources/res/{zoom}/"
            for path in paths:
                normalized = path.replace("\\", "/") if isinstance(path, str) else ""
                direct_x4_alias = (
                    mode == "direct"
                    and normalized.startswith("Assets/Resources/res/x4/")
                )
                if (
                    not normalized.startswith(prefix)
                    and not direct_x4_alias
                ) or ".." in Path(normalized).parts:
                    _fail("ASSET-UI-007", f"{binding['source_asset']} physical path không canonical")
        binding_contract = dict(binding)
        binding_contract.setdefault("call_sites", call_sites)
        binding_contract.setdefault("consumer_sites", consumer_sites)
        binding_fingerprints.append(_runtime_binding_fingerprint(binding_contract))
    if len(set(binding_fingerprints)) != 8:
        _fail("ASSET-UI-007", "manifest logical/source/physical/consumer binding bị duplicate")
    return relative_path, map_sha, payload


def require_approved_review(review_path: Path, manifest: dict[str, Any], base_dir: Path) -> tuple[str, list[dict[str, str]]]:
    entries, evidence = _review_records(review_path)
    by_name = {Path(entry.get("asset", entry.get("source_path", ""))).name: entry for entry in entries}
    derived_sprite_paths = {
        relative
        for sprite in manifest.get("sprites", [])
        if isinstance(sprite, dict)
        for relative in sprite.get("zooms", {}).values()
        if isinstance(relative, str)
    }
    production = isinstance(manifest.get("build_inputs"), dict) and "source_base" in manifest["build_inputs"]
    reviewed_source_paths = {
        sprite.get("source_path")
        for sprite in manifest.get("sprites", [])
        if isinstance(sprite, dict) and isinstance(sprite.get("source_path"), str)
    }
    source_provenance = manifest.get("build_inputs", {}).get(
        "source_provenance", []
    )
    provenance_by_path = (
        _validate_production_review_provenance(
            review_path, entries, manifest, base_dir
        )
        if production
        else {
            record.get("source_path"): record
            for record in source_provenance
            if isinstance(record, dict)
            and isinstance(record.get("source_path"), str)
        }
    )
    if production:
        evidence.extend(
            validate_production_artifact_review(
                review_path, manifest, base_dir
            ).values()
        )
    reviewed_source_paths.update(provenance_by_path)
    for relative in _source_paths(manifest):
        source = _path(base_dir, relative, "ASSET-REVIEW-010")
        if not source.is_file():
            _fail("ASSET-REVIEW-010", f"thiếu source đã review: {relative}")
        if relative in derived_sprite_paths or (production and relative not in reviewed_source_paths):
            continue
        entry = by_name.get(source.name)
        if entry is None or entry["source_sha256"] != _sha(source):
            _fail("ASSET-REVIEW-010", f"source SHA không khớp review: {relative}")
    for sprite in manifest.get("sprites", []):
        source = _path(base_dir, sprite["source_path"], "ASSET-REVIEW-010")
        if sprite.get("source_sha256") != _sha(source):
            _fail("ASSET-REVIEW-010", f"source SHA manifest không khớp: {sprite['source_path']}")
    for relative, record in provenance_by_path.items():
        source = _path(base_dir, relative, "ASSET-REVIEW-010")
        if record.get("source_sha256") != _sha(source):
            _fail(
                "ASSET-REVIEW-010",
                f"source provenance SHA mismatch: {relative}",
            )
    return review_digest(review_path), evidence


def _schema(manifest: dict[str, Any]) -> None:
    if not isinstance(manifest, dict) or set(manifest) != REQUIRED_TOP_LEVEL:
        _fail("ASSET-SCHEMA-001", "manifest thiếu/thừa trường canonical")
    _int(manifest["schema_version"], "ASSET-SCHEMA-001", "schema_version")
    _int(manifest["source_catalog_version"], "ASSET-SCHEMA-001", "source_catalog_version")
    versions = manifest["versions"]
    if not isinstance(versions, dict) or set(versions) != {"runtime_slice", "data", "map", "resource"}:
        _fail("ASSET-SCHEMA-001", "versions phải đủ runtime_slice/data/map/resource")
    for key, value in versions.items():
        _int(value, "ASSET-SCHEMA-001", f"versions.{key}")
    phase = manifest["phase"]
    if not isinstance(phase, dict) or phase.get("name") != "lang-la":
        _fail("ASSET-SCHEMA-001", "phase.name phải là lang-la")
    sprites = manifest["sprites"]
    if not isinstance(sprites, list) or not sprites:
        _fail("ASSET-SCHEMA-001", "sprites không được rỗng")
    sprite_ids: set[int] = set()
    for sprite in sprites:
        if not isinstance(sprite, dict):
            _fail("ASSET-SCHEMA-001", "sprite phải là object")
        sprite_id = _int(sprite.get("id"), "ASSET-ID-001", "sprite.id")
        if sprite_id in POC_IDS or not 0 <= sprite_id <= 32767 or sprite_id in sprite_ids:
            _fail("ASSET-ID-001", f"sprite ID không hợp lệ/trùng: {sprite_id}")
        sprite_ids.add(sprite_id)
        if set(sprite.get("zooms", {})) != set(ZOOMS) or not isinstance(sprite.get("source_path"), str):
            _fail("ASSET-ZOOM-003", f"sprite {sprite_id} thiếu x1..x4/source")
    for part in manifest["parts"]:
        if not isinstance(part, dict):
            _fail("ASSET-SCHEMA-001", "part phải là object")
        _int(part.get("id"), "ASSET-PART-001", "part.id")
        for frame in part.get("frames", []):
            if not isinstance(frame, dict) or frame.get("sprite_id") not in sprite_ids:
                _fail("ASSET-SPRITE-002", "part trỏ sprite ID không có trong manifest")
    skills = manifest["skills"]
    if not isinstance(skills, list) or len(skills) != 4:
        _fail("ASSET-SKILL-004", "phải có đúng bốn skill pose")
    poses = set()
    for skill in skills:
        skill_id = _int(skill.get("id"), "ASSET-SKILL-004", "skill.id")
        if not 0 <= skill_id <= 32767 or skill_id in POC_IDS:
            _fail("ASSET-SKILL-004", f"skill ID không hợp lệ: {skill_id}")
        pose = _int(skill.get("pose_slot"), "ASSET-SKILL-004", "skill.pose_slot")
        if pose in poses or not 0 <= pose <= 3:
            _fail("ASSET-SKILL-004", "pose slot phải là 0..3, không trùng")
        poses.add(pose)
        levels = skill.get("levels")
        if not isinstance(levels, list) or len(levels) != 12 or any(not isinstance(level, int) or isinstance(level, bool) or not 0 <= level <= 255 for level in levels):
            _fail("ASSET-SKILL-004", "skill phải có đúng 12 levelId nguyên")
    if poses != {0, 1, 2, 3}:
        _fail("ASSET-SKILL-004", "thiếu pose 0..3")
    inputs = manifest["build_inputs"]
    if not isinstance(inputs, dict) or not inputs.get("legacy_part_path") or not inputs.get("tile_set_info_path"):
        _fail("ASSET-SCHEMA-001", "thiếu build_inputs legacy part/tileset")
    if not isinstance(inputs.get("tile_set_baseline"), dict):
        _fail("ASSET-TILESET-010", "thiếu tile-set baseline golden")
    mob = inputs.get("mob")
    generated_mob = isinstance(mob, dict) and mob.get("payload_mode") == "generated"
    if (
        not isinstance(mob, dict)
        or (
            not generated_mob
            and (
                not isinstance(mob.get("bytes_path"), str)
                or not isinstance(mob.get("data_path"), str)
            )
        )
        or not isinstance(mob.get("img"), dict)
        or set(mob["img"]) != set(ZOOMS)
        or any(not isinstance(mob["img"][zoom], str) for zoom in ZOOMS)
    ):
        _fail("ASSET-MOB-003", "Mob binary/image path malformed")
    if not inputs.get("runtime_key_map"):
        _fail("ASSET-UI-007", "thiếu runtime key map contract")
    try:
        validate_animation_closure(inputs.get("animation_sheets"), manifest.get("animations"))
    except AnimationError as error:
        _fail("ASSET-ANIMATION-008", str(error))
    map_config = inputs.get("map")
    if not isinstance(map_config, dict):
        _fail("ASSET-MAP-004", "thiếu map layout")
    tile_size = map_config.get("tile_size")
    if not isinstance(tile_size, list) or len(tile_size) != 2 or any(not isinstance(value, int) or isinstance(value, bool) or value <= 0 for value in tile_size):
        _fail("ASSET-MAP-004", "tile_size phải là hai số nguyên dương")
    width, height = map_config.get("width"), map_config.get("height")
    if not isinstance(width, int) or isinstance(width, bool) or not 1 <= width <= 255 or not isinstance(height, int) or isinstance(height, bool) or not 1 <= height <= 255:
        _fail("ASSET-MAP-004", "map width/height phải trong 1..255")
    indexes = map_config.get("indexes")
    if not isinstance(indexes, list) or len(indexes) != width * height or any(not isinstance(value, int) or isinstance(value, bool) or not 0 <= value <= 99 for value in indexes):
        _fail("ASSET-MAP-004", "map indexes malformed")
    tile_art = inputs.get("tile_art")
    if not isinstance(tile_art, list) or any(not isinstance(item, dict) for item in tile_art) or {item.get("index") for item in tile_art} != set(range(100)):
        _fail("ASSET-MAP-004", "tile_art phải đủ index 0..99")
    for item in tile_art:
        if not isinstance(item.get("index"), int) or isinstance(item.get("index"), bool):
            _fail("ASSET-MAP-004", "tile index phải là integer")
        if set(item.get("sources", {})) != set(ZOOMS):
            _fail("ASSET-ZOOM-003", f"tile {item.get('index')} thiếu zoom")
    backgrounds = inputs.get("background_sources")
    if not isinstance(backgrounds, dict) or set(backgrounds) != {"b00", "b01", "b02", "b03"}:
        _fail("ASSET-BG-005", "background phải đủ b00..b03")
    if any(set(value) != set(ZOOMS) for value in backgrounds.values()):
        _fail("ASSET-BG-005", "background phải đủ bốn zoom")
    hud = inputs.get("hud_sources")
    if not isinstance(hud, dict) or set(hud) != {"naruto-hp-bg", "naruto-hp-fill", "naruto-chakra-bg", "naruto-chakra-fill"}:
        _fail("ASSET-UI-007", "HUD phải đủ bốn resource")
    if any(set(value) != set(ZOOMS) for value in hud.values()):
        _fail("ASSET-UI-007", "HUD phải đủ bốn zoom")
    ui_bindings = manifest["ui_bindings"]
    if not isinstance(ui_bindings, list) or len(ui_bindings) != 8:
        _fail("ASSET-UI-007", "phải có đúng tám UI binding")
    modes = [binding.get("mode") for binding in ui_bindings if isinstance(binding, dict)]
    if modes.count("direct") != 1 or modes.count("slice_group") != 2 or modes.count("overlay_required") != 5:
        _fail("ASSET-UI-007", "UI binding phải đúng tỷ lệ 1 direct/2 slice_group/5 overlay_required")
    for binding in ui_bindings:
        if not isinstance(binding, dict):
            _fail("ASSET-UI-007", "UI binding phải là object")
        mode = binding.get("mode")
        if mode not in {"direct", "slice_group", "overlay_required"}:
            _fail("ASSET-UI-007", "UI binding có mode không biết")
        keys = binding.get("logical_keys")
        if not isinstance(keys, list) or any(not isinstance(key, str) or not key or "placeholder" in key.lower() for key in keys):
            _fail("ASSET-UI-007", "logical_keys malformed/placeholder")
        future = binding.get("future_logical_key")
        if mode == "overlay_required":
            if keys or not isinstance(future, str) or not future or "placeholder" in future.lower():
                _fail("ASSET-UI-007", "overlay_required phải có đúng future key")
        elif not keys or future is not None:
            _fail("ASSET-UI-007", "direct/slice_group phải dùng logical key thật")
        if set(binding.get("physical_paths", {})) != set(ZOOMS):
            _fail("ASSET-UI-007", "UI binding thiếu physical path x1..x4")
        sources = binding.get("sources")
        if sources is not None and (
            not isinstance(sources, dict)
            or set(sources) != set(ZOOMS)
            or any(not isinstance(sources[zoom], str) for zoom in ZOOMS)
        ):
            _fail("ASSET-UI-007", "UI binding sources thiếu x1..x4")
    try:
        validate_catalog(manifest["catalog"])
    except ValueError as error:
        _fail("ASSET-CATALOG-009", str(error))


def _copy(base_dir: Path, relative: str, destination: Path, *, png=False) -> dict[str, Any]:
    source = _path(base_dir, relative, "ASSET-BUILD-013")
    if not source.is_file():
        _fail("ASSET-BUILD-013", f"thiếu artifact input: {relative}")
    data = source.read_bytes()
    meta = _png_meta(data) if png else {"bytes": len(data)}
    destination.parent.mkdir(parents=True, exist_ok=True)
    destination.write_bytes(data)
    return {"path": destination.as_posix(), "sha256": hashlib.sha256(data).hexdigest(), "bytes": len(data), **meta}


def _all_output_files(root: Path) -> list[dict[str, Any]]:
    outputs = []
    for path in sorted(root.rglob("*")):
        if path.is_file() and path.name != "asset-lock.json":
            data = path.read_bytes()
            record = {"path": path.relative_to(root).as_posix(), "sha256": hashlib.sha256(data).hexdigest(), "bytes": len(data)}
            if path.suffix.lower() == ".png":
                record.update(_png_meta(data))
            outputs.append(record)
    return outputs


def _pair_sha(root: Path) -> str:
    digest = hashlib.sha256()
    for path in sorted((root / "server").rglob("*")) + sorted((root / "client").rglob("*")):
        if path.is_file():
            digest.update(path.relative_to(root).as_posix().encode("utf-8"))
            digest.update(path.read_bytes())
    return digest.hexdigest()


def _manifest_sha(manifest: dict[str, Any]) -> str:
    canonical = json.dumps(manifest, ensure_ascii=False, sort_keys=True, separators=(",", ":")).encode("utf-8")
    return hashlib.sha256(canonical).hexdigest()


def _build_into(manifest: dict[str, Any], base_dir: Path, root: Path) -> dict[str, Any]:
    inputs = manifest["build_inputs"]
    server = root / "server"
    client = root / "client"
    legacy = _path(base_dir, inputs["legacy_part_path"], "ASSET-BUILD-013").read_bytes()
    tileset = _path(base_dir, inputs["tile_set_info_path"], "ASSET-BUILD-013").read_bytes()
    try:
        validate_tileset_baseline(tileset, inputs["tile_set_baseline"])
        validate_animation_closure(
            inputs["animation_sheets"],
            manifest["animations"],
            check_paths=lambda relative: _path(base_dir, relative, "ASSET-ANIMATION-008"),
            # Ô là giới hạn trên: cell nhỏ hơn ô hợp lệ (mỗi khung một kích cỡ
            # sau khi bỏ nhồi-vào-hộp), cell vượt ô là lỗi.
            check_cell_size=lambda relative, expected: all(
                actual <= limit
                for actual, limit in zip(_load_png(base_dir, relative).size, expected)
            ),
        )
    except MapError as error:
        message = str(error)
        code = message.split(":", 1)[0] if message.startswith("ASSET-") else "ASSET-TILESET-011"
        _fail(code, message)
    except AnimationError as error:
        _fail("ASSET-ANIMATION-008", str(error))

    # Bảng gốc chỉ khai báo tile 1-39 của Ngọc Rồng. Map Làng Lá dùng tile 60-69
    # và 83-88, không khai thì client đặt types[]=0 cho cả màn → nhân vật rơi
    # xuyên đất. Trước đây phần thêm này làm TAY rồi commit thẳng file .bin, nên
    # build lại là âm thầm lùi về bảng gốc và tái sinh lỗi. Giờ phép thêm nằm
    # trong manifest và do đường ống thi hành, kết quả bị ghim bằng golden.
    if "tile_set_overlay" in inputs:
        try:
            tileset = apply_tile_set_overlay(tileset, inputs["tile_set_overlay"])
        except TilesetError as error:
            message = str(error)
            code = message.split(":", 1)[0] if message.startswith("ASSET-") else "ASSET-TILESET-012"
            _fail(code, message)
    try:
        part_data = build_part(manifest, legacy)
        parse_part(part_data)
    except Exception as error:
        _fail("ASSET-PART-001", str(error))
    (server / "part.bin").parent.mkdir(parents=True, exist_ok=True)
    (server / "part.bin").write_bytes(part_data)
    (server / "tile_set_info.bin").write_bytes(tileset)

    map_config = inputs.get("map")
    if not isinstance(map_config, dict):
        _fail("ASSET-MAP-004", "thiếu map layout")
    width = _int(map_config.get("width"), "ASSET-MAP-004", "map.width")
    height = _int(map_config.get("height"), "ASSET-MAP-004", "map.height")
    indexes = map_config.get("indexes")
    if not isinstance(indexes, list) or len(indexes) != width * height:
        _fail("ASSET-MAP-004", "map layout sai kích thước")
    spawn = map_config.get("spawn", [100, 384])
    tile_size = map_config.get("tile_size", [24, 24])
    if not isinstance(spawn, list) or len(spawn) != 2 or not isinstance(tile_size, list) or len(tile_size) != 2:
        _fail("ASSET-MAP-004", "spawn/tile_size sai")
    col, row = spawn[0] // tile_size[0], spawn[1] // tile_size[1]
    if not 0 <= col < width or not 0 <= row < height:
        _fail("ASSET-MAP-004", "spawn nằm ngoài map")
    spawn_index = row * width + col
    client_solid = map_config.get("client_solid")
    if not isinstance(client_solid, list) or len(client_solid) != len(indexes) or client_solid[spawn_index]:
        _fail("ASSET-MAP-004", "spawn phải walkable và có client collision")
    tile_rows = [{"id": item["index"], "solid": item.get("solid", False)} for item in inputs["tile_art"]]
    try:
        validate_collision_semantics(tile_rows, indexes, client_solid)
        collision = encode_collision(width, height, client_solid, spawn_index=spawn_index, tile_size=tile_size)
    except MapError as error:
        _fail("ASSET-COLLISION-006", str(error))
    (server / "map-1-collision.bin").write_bytes(collision)

    # Lưới đúng chưa đủ: máy chủ chặn theo `tile_art[].solid`, còn client cho
    # đứng theo bit sàn trong bảng tile. Tile máy chủ bảo "đặc" mà client không
    # khai là sàn thì hai bên nói ngược nhau — người chơi rơi xuyên đất rồi bị
    # máy chủ kéo ngược. Đây là chốt chặn CUỐI: nó đỏ kể cả khi ai đó xoá hẳn
    # `tile_set_overlay` khỏi manifest, tức quả mìn không dựng lại được nữa.
    try:
        validate_floor_closure(
            tileset,
            tile_set_index=inputs.get("tile_set_overlay", {}).get("tile_set_index", 0),
            base_tile_count=inputs["tile_set_baseline"]["semantic"]["tile_count"],
            solid_used={row["id"] for row in tile_rows if row["solid"]} & set(indexes),
        )
    except TilesetError as error:
        message = str(error)
        code = message.split(":", 1)[0] if message.startswith("ASSET-") else "ASSET-TILESET-013"
        _fail(code, message)

    try:
        skill_data = encode_skill(manifest["skills"])
    except ValueError as error:
        _fail("ASSET-SKILL-004", str(error))
    (server / "skill.bin").write_bytes(skill_data)

    for sprite in manifest["sprites"]:
        for zoom in ZOOMS:
            _copy(base_dir, sprite["zooms"][zoom], client / zoom / "smallimage" / f"Small{sprite['id']}.png", png=True)
    for tile in sorted(inputs["tile_art"], key=lambda item: item["index"]):
        index = tile["index"]
        for zoom in ZOOMS:
            name = f"t_{index + 1:02d}.png" if zoom == "x1" else f"1${index + 1}.png"
            _copy(base_dir, tile["sources"][zoom], client / zoom / "t" / ("1" if zoom == "x1" else "") / name, png=True)
    for layer, sources in inputs["background_sources"].items():
        for zoom in ZOOMS:
            _copy(base_dir, sources[zoom], client / zoom / "bg" / f"{layer}.png", png=True)
    for name, sources in inputs["hud_sources"].items():
        for zoom in ZOOMS:
            _copy(base_dir, sources[zoom], client / zoom / "mainImage" / f"{name}.png", png=True)

    mob = inputs["mob"]
    mob_id = _int(mob.get("id"), "ASSET-MOB-003", "mob.id")
    if mob.get("payload_mode") == "generated":
        try:
            mob_payload = encode_mob(mob)
        except MobError as error:
            _fail("ASSET-MOB-003", str(error))
        mob_bytes_target = client / "Mob" / f"{mob_id}.bytes"
        mob_data_target = client / "Mob" / f"{mob_id}" / "data.bytes"
        mob_bytes_target.parent.mkdir(parents=True, exist_ok=True)
        mob_data_target.parent.mkdir(parents=True, exist_ok=True)
        mob_bytes_target.write_bytes(mob_payload)
        mob_data_target.write_bytes(mob_payload)
    else:
        mob_bytes_path = _path(base_dir, mob["bytes_path"], "ASSET-MOB-003")
        mob_data_path = _path(base_dir, mob["data_path"], "ASSET-MOB-003")
        if not mob_bytes_path.is_file() or not mob_data_path.is_file():
            _fail("ASSET-MOB-003", "thiếu Mob bytes/data payload")
        try:
            validate_mob_golden(mob, mob_bytes_path.read_bytes(), mob_data_path.read_bytes())
        except ValueError as error:
            _fail("ASSET-MOB-003", str(error))
        mob_bytes = _copy(base_dir, mob["bytes_path"], client / "Mob" / f"{mob_id}.bytes")
        mob_data = _copy(base_dir, mob["data_path"], client / "Mob" / f"{mob_id}" / "data.bytes")
        if mob_bytes["sha256"] != mob["bytes_sha256"] or mob_data["sha256"] != mob["data_sha256"]:
            _fail("ASSET-MOB-003", "Mob golden SHA không khớp artifact")
    for zoom in ZOOMS:
        _copy(base_dir, mob["img"][zoom], client / zoom / "Mob" / str(mob_id) / "img.png", png=True)

    for binding in manifest["ui_bindings"]:
        source_path = binding.get("source_path")
        if not isinstance(source_path, str):
            _fail("ASSET-UI-007", "UI binding thiếu source_path")
        written_paths: set[Path] = set()
        for zoom in ZOOMS:
            paths = binding.get("physical_paths", {}).get(zoom, [])
            if not paths:
                _fail("ASSET-UI-007", f"UI binding thiếu physical path {zoom}")
            for path in paths:
                normalized = path.replace("\\", "/")
                marker = "res/"
                if marker not in normalized:
                    _fail("ASSET-UI-007", f"UI binding physical path malformed: {path}")
                resource_path = normalized.split(marker, 1)[1]
                target_zoom, separator, relative = resource_path.partition("/")
                if not separator or target_zoom not in ZOOMS or not relative:
                    _fail("ASSET-UI-007", f"UI binding physical path malformed: {path}")
                destination = client / target_zoom / relative
                if destination in written_paths:
                    continue
                zoom_sources = binding.get("sources")
                zoom_source = (
                    zoom_sources.get(target_zoom)
                    if isinstance(zoom_sources, dict)
                    else source_path
                )
                if not isinstance(zoom_source, str):
                    _fail("ASSET-UI-007", f"UI binding thiếu source {target_zoom}")
                _copy(base_dir, zoom_source, destination, png=True)
                written_paths.add(destination)

    mymap = bytes([width, height, *indexes])
    (client / "mymap" / "1.bytes").parent.mkdir(parents=True, exist_ok=True)
    (client / "mymap" / "1.bytes").write_bytes(mymap)
    evidence_root = client / "evidence"
    tile_images = {item["index"]: _load_png(base_dir, item["sources"]["x1"]) for item in inputs["tile_art"]}
    render_tile_contact_sheet(tile_images, evidence_root / "tile-contact-x1.png")
    for zoom in ZOOMS:
        images = {item["index"]: _load_png(base_dir, item["sources"][zoom]) for item in inputs["tile_art"]}
        render_full_map(images, indexes, width, height, evidence_root / f"full-map-{zoom}.png")
    for sheet in inputs.get("animation_sheets", []):
        cells = [(cell["slot_index"], _load_png(base_dir, cell["path"])) for cell in sheet["cells"]]
        rendered = compose_action_sheet(cells, sheet["slot_count"], tuple(sheet["cell_size"]))
        path = evidence_root / "animation" / f"{sheet['name']}.png"
        path.parent.mkdir(parents=True, exist_ok=True)
        rendered.save(path, format="PNG")
    hud_evidence = evidence_root / "hud-golden-25-50-100-x4.png"
    render_hud_golden(
        {
            name: _load_png(base_dir, sources["x4"])
            for name, sources in inputs["hud_sources"].items()
        },
        hud_evidence,
    )
    sprites_by_id = {sprite["id"]: sprite for sprite in manifest["sprites"]}
    iruka_evidence = evidence_root / "iruka-composite-light-dark-x4.png"
    render_iruka_composite_contact_sheet(
        {
            "head": _load_png(base_dir, sprites_by_id[27075]["zooms"]["x4"]),
            "body": _load_png(base_dir, sprites_by_id[27076]["zooms"]["x4"]),
            "leg": _load_png(base_dir, sprites_by_id[27077]["zooms"]["x4"]),
        },
        iruka_evidence,
    )
    declared_evidence = inputs.get("artifact_evidence")
    if isinstance(declared_evidence, dict):
        generated = {
            "hud_golden": hud_evidence,
            "iruka_composite": iruka_evidence,
        }
        for kind, path in generated.items():
            record = declared_evidence.get(kind)
            if (
                not isinstance(record, dict)
                or record.get("sha256") != _sha(path)
            ):
                _fail(
                    "ASSET-ARTIFACT-REVIEW-016",
                    f"generated evidence lệch artifact đã review: {kind}",
                )

    catalog_data = encode_catalog(manifest["catalog"])
    (server / "catalog.json").write_bytes(catalog_data)
    return {"server_client_pair_sha256": _pair_sha(root)}


def build(manifest: dict, out: Path, review_path: Path, *, base_dir: Path | None = None) -> dict[str, Any]:
    try:
        require_review_gate(Path(review_path))
    except ReviewGateError as error:
        # Existing unit fixtures predate artqa's full document. Production reviews
        # always carry source_root and must remain fail-closed.
        try:
            fixture_document = json.loads(Path(review_path).read_text(encoding="utf-8"))
        except (OSError, json.JSONDecodeError):
            fixture_document = None
        if not (isinstance(fixture_document, dict) and "source_root" not in fixture_document):
            _fail("ASSET-REVIEW-GATE-001", "; ".join(error.issues))
    except Exception as error:
        _fail("ASSET-REVIEW-GATE-001", str(error))
    """Sinh toàn bộ artifact trong thư mục mới; không bao giờ ghi đè output cũ."""
    out = Path(out)
    if out.exists():
        _fail("ASSET-BUILD-012", f"output đã tồn tại: {out}")
    base_dir = Path(base_dir or Path.cwd()).resolve()
    if isinstance(manifest.get("build_inputs"), dict) and "source_base" in manifest["build_inputs"]:
        base_dir = resolve_manifest_base(Path("<manifest>"), manifest, base_dir=base_dir)
        _validate_production_closure(manifest, base_dir)
    review_sha, evidence = require_approved_review(review_path, manifest, base_dir)
    _schema(manifest)
    runtime_key_map_path, runtime_key_map_sha, runtime_key_map = _validate_runtime_keys(manifest, base_dir)
    parent = out.parent
    parent.mkdir(parents=True, exist_ok=True)
    temporary = Path(tempfile.mkdtemp(prefix=f".{out.name}.", dir=parent))
    try:
        _build_into(manifest, base_dir, temporary)
        outputs = _all_output_files(temporary)
        generated_evidence = [record for record in outputs if record["path"].startswith("client/evidence/")]
        closure = {
            "sprite_ids": sorted(sprite["id"] for sprite in manifest["sprites"]),
            "part_ids": sorted(part["id"] for part in manifest["parts"]),
            "skill_ids": sorted(skill["id"] for skill in manifest["skills"]),
            "mob_ids": sorted(mob["id"] for mob in manifest.get("mobs", [])),
            "map_ids": sorted(item["id"] for item in manifest.get("maps", [])),
            "frame_counts": {str(part["id"]): len(part["frames"]) for part in manifest["parts"]},
            "animations": [
                {"name": item["name"], "slot_count": item["slot_count"], "frames": [
                    {"slot_index": frame["slot_index"], "client_cf": list(frame["client_cf"])}
                    for frame in item["frames"]
                ]}
                for item in manifest["animations"]
            ],
            "animation_sheets": [
                {
                    "name": sheet["name"],
                    "slot_count": sheet["slot_count"],
                    "cell_size": list(sheet["cell_size"]),
                    "cells": [
                        {"slot_index": cell["slot_index"], "path": cell["path"], "bounds": [
                            cell["slot_index"] * sheet["cell_size"][0], 0,
                            (cell["slot_index"] + 1) * sheet["cell_size"][0], sheet["cell_size"][1],
                        ]}
                        for cell in sheet["cells"]
                    ],
                }
                for sheet in manifest["build_inputs"]["animation_sheets"]
            ],
            "mob": {
                "id": manifest["build_inputs"]["mob"]["id"],
                "actions": {
                    action: {"frame_count": len(frames), "frames": frames}
                    for action, frames in manifest["build_inputs"]["mob"]["frames"].items()
                },
                "images": manifest["build_inputs"]["mob"]["images"],
                "arr_frames": manifest["build_inputs"]["mob"]["arr_frames"],
                "trailer": manifest["build_inputs"]["mob"]["trailer"],
            },
            "runtime_bindings": [
                {
                    key: entry.get(key)
                    for key in ("source_asset", "mode", "logical_keys", "future_logical_key", "physical_paths", "call_sites", "consumer_sites")
                }
                for entry in runtime_key_map["entries"]
            ],
        }
        dimensions = {
            record["path"]: [record["width"], record["height"]]
            for record in outputs
            if "width" in record and "height" in record
        }
        lock = {
            "schema_version": manifest["schema_version"],
            "slice": manifest["phase"]["name"],
            "versions": manifest["versions"],
            "review_sha256": review_sha,
            "inputs": [{"path": relative, "sha256": _sha(_path(base_dir, relative, "ASSET-BUILD-013"))} for relative in _source_paths(manifest)],
            "outputs": outputs,
            "evidence": sorted(evidence, key=lambda item: item["path"]),
            "closure": closure,
            "provenance": {
                "manifest": {"source": "canonical-json", "sha256": _manifest_sha(manifest)},
                "runtime_key_map": {"path": runtime_key_map_path or "<inline>", "sha256": runtime_key_map_sha},
                "dimensions": dimensions,
                "ids": {key: closure[key] for key in ("sprite_ids", "part_ids", "skill_ids", "mob_ids", "map_ids")},
                "frame_counts": closure["frame_counts"],
                "review_evidence": sorted(evidence, key=lambda item: item["path"]),
                "generated_evidence": generated_evidence,
                "artifacts": outputs,
                "evidence_provenance": sorted([*evidence, *generated_evidence], key=lambda item: item["path"]),
            },
            "server_client_pair_sha256": _pair_sha(temporary),
        }
        (temporary / "asset-lock.json").write_bytes(lock_bytes(lock))
        temporary.replace(out)
        return lock
    except Exception:
        shutil.rmtree(temporary, ignore_errors=True)
        raise
