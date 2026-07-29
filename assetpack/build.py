"""Build đầy đủ asset server/client sau các cổng schema, review và provenance.

Build không nhận một manifest tối thiểu rồi tự đoán phần còn thiếu. Mọi đầu vào
runtime phải được khai báo tường minh; nếu thiếu, build dừng trước khi tạo output.
"""

from __future__ import annotations

from io import BytesIO
import hashlib
import json
from pathlib import Path
import shutil
import tempfile
from typing import Any

from PIL import Image

from .animation import AnimationError, compose_action_sheet, validate_animation_closure
from .catalog import encode_catalog, validate_catalog
from .lock import lock_bytes
from .map import MapError, encode_collision, validate_collision_semantics, validate_tileset_baseline
from .mob import validate_mob_golden
from .part import build_part, parse_part
from .skill import encode_skill
from .visual import render_full_map, render_tile_contact_sheet


class BuildError(ValueError):
    """Lỗi build có mã ổn định để canary/mutation bắt đúng nguyên nhân."""


REQUIRED_TOP_LEVEL = {
    "schema_version", "source_catalog_version", "versions", "phase",
    "sprites", "parts", "animations", "maps", "tiles", "backgrounds",
    "mobs", "skills", "catalog", "ui_bindings", "build_inputs",
}
ZOOMS = ("x1", "x2", "x3", "x4")
POC_IDS = {9000, 9001, 9002}


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


def review_digest(review_path: Path) -> str:
    digest = hashlib.sha256()
    for path, _ in _review_documents(review_path):
        digest.update(path.name.encode("utf-8"))
        digest.update(path.read_bytes())
    return digest.hexdigest()


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
    return sorted({path for path in paths if isinstance(path, str)})


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


def _validate_runtime_keys(manifest: dict[str, Any], base_dir: Path) -> tuple[str | None, str, dict[str, Any]]:
    payload, relative_path, map_sha = _runtime_key_map(manifest, base_dir)
    entries = payload.get("entries") if isinstance(payload, dict) else None
    bindings = manifest.get("ui_bindings")
    if not isinstance(entries, list) or len(entries) != 8 or not isinstance(bindings, list) or len(bindings) != 8:
        _fail("ASSET-UI-007", "runtime key map và manifest phải có đúng tám binding")
    for entry in entries:
        if not isinstance(entry, dict) or not isinstance(entry.get("source_asset"), str) or not entry["source_asset"]:
            _fail("ASSET-UI-007", "runtime entry thiếu source_asset hợp lệ")
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
            source = _path(base_dir, site["file"], "ASSET-UI-007")
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
                if not normalized.startswith(prefix) or ".." in Path(normalized).parts:
                    _fail("ASSET-UI-007", f"{binding['source_asset']} physical path không canonical")
        binding_contract = dict(binding)
        binding_contract.setdefault("call_sites", call_sites)
        binding_contract.setdefault("consumer_sites", consumer_sites)
        binding_fingerprints.append(_runtime_binding_fingerprint(binding_contract))
        if mode == "overlay_required" and any(binding["future_logical_key"] not in site["expression"] for site in consumer_sites):
            _fail("ASSET-UI-007", f"{binding['source_asset']} consumer precondition thiếu future logical key")
    if len(set(binding_fingerprints)) != 8:
        _fail("ASSET-UI-007", "manifest logical/source/physical/consumer binding bị duplicate")
    return relative_path, map_sha, payload


def require_approved_review(review_path: Path, manifest: dict[str, Any], base_dir: Path) -> tuple[str, list[dict[str, str]]]:
    entries, evidence = _review_records(review_path)
    by_name = {Path(entry.get("asset", entry.get("source_path", ""))).name: entry for entry in entries}
    for relative in _source_paths(manifest):
        source = _path(base_dir, relative, "ASSET-REVIEW-010")
        if not source.is_file():
            _fail("ASSET-REVIEW-010", f"thiếu source đã review: {relative}")
        entry = by_name.get(source.name)
        if entry is None or entry["source_sha256"] != _sha(source):
            _fail("ASSET-REVIEW-010", f"source SHA không khớp review: {relative}")
    for sprite in manifest.get("sprites", []):
        source = _path(base_dir, sprite["source_path"], "ASSET-REVIEW-010")
        if sprite.get("source_sha256") != _sha(source):
            _fail("ASSET-REVIEW-010", f"source SHA manifest không khớp: {sprite['source_path']}")
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
    if (
        not isinstance(mob, dict)
        or not isinstance(mob.get("bytes_path"), str)
        or not isinstance(mob.get("data_path"), str)
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
            check_cell_size=lambda relative, expected: _load_png(base_dir, relative).size == expected,
        )
    except MapError as error:
        message = str(error)
        code = message.split(":", 1)[0] if message.startswith("ASSET-") else "ASSET-TILESET-011"
        _fail(code, message)
    except AnimationError as error:
        _fail("ASSET-ANIMATION-008", str(error))
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
        collision = encode_collision(width, height, indexes, spawn_index=spawn_index)
    except MapError as error:
        _fail("ASSET-COLLISION-006", str(error))
    (server / "map-1-collision.bin").write_bytes(collision)

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
    mob_bytes_path = _path(base_dir, mob["bytes_path"], "ASSET-MOB-003")
    mob_data_path = _path(base_dir, mob["data_path"], "ASSET-MOB-003")
    if not mob_bytes_path.is_file() or not mob_data_path.is_file():
        _fail("ASSET-MOB-003", "thiếu Mob bytes/data payload")
    try:
        validate_mob_golden(mob, mob_bytes_path.read_bytes(), mob_data_path.read_bytes())
    except ValueError as error:
        _fail("ASSET-MOB-003", str(error))
    mob_id = _int(mob.get("id"), "ASSET-MOB-003", "mob.id")
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
        for zoom in ZOOMS:
            paths = binding.get("physical_paths", {}).get(zoom, [])
            if not paths:
                _fail("ASSET-UI-007", f"UI binding thiếu physical path {zoom}")
            for path in paths:
                normalized = path.replace("\\", "/")
                marker = f"res/{zoom}/"
                if marker in normalized:
                    normalized = normalized.split(marker, 1)[1]
                _copy(base_dir, source_path, client / zoom / normalized, png=True)

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

    catalog_data = encode_catalog(manifest["catalog"])
    (server / "catalog.json").write_bytes(catalog_data)
    return {"server_client_pair_sha256": _pair_sha(root)}


def build(manifest: dict, out: Path, review_path: Path, *, base_dir: Path | None = None) -> dict[str, Any]:
    """Sinh toàn bộ artifact trong thư mục mới; không bao giờ ghi đè output cũ."""
    out = Path(out)
    if out.exists():
        _fail("ASSET-BUILD-012", f"output đã tồn tại: {out}")
    base_dir = Path(base_dir or Path.cwd()).resolve()
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
