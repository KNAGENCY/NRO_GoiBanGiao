"""Schema và các cổng kiểm tra manifest/baseline của Task 1."""

from dataclasses import dataclass
import hashlib
import json
from pathlib import Path
from typing import Any


class ManifestError(ValueError):
    """Lỗi contract có mã ổn định cho CI và người làm asset."""


@dataclass(frozen=True)
class Sprite:
    id: int
    zooms: dict[str, str]
    source_path: str
    source_sha256: str


@dataclass(frozen=True)
class PartFrame:
    sprite_id: int
    dx: int
    dy: int


@dataclass(frozen=True)
class Part:
    id: int
    frames: list[PartFrame]


@dataclass(frozen=True)
class AnimationFrame:
    slot_index: int
    client_cf: list[int]


@dataclass(frozen=True)
class Animation:
    name: str
    frames: list[AnimationFrame]


@dataclass(frozen=True)
class Map:
    id: int
    tile_id: int
    background_id: int


@dataclass(frozen=True)
class Tile:
    id: int
    semantic: str
    solid: bool


@dataclass(frozen=True)
class Background:
    id: int
    layers: int


@dataclass(frozen=True)
class Mob:
    id: int
    name: str


@dataclass(frozen=True)
class Catalog:
    slice: str
    default_parts: dict[str, int]


@dataclass(frozen=True)
class Version:
    runtime_slice: int
    data: int = 0
    map: int = 0
    resource: int = 0


REQUIRED_LEGACY_KINDS = {
    "part.bin",
    "NR_part",
    "NR_image",
    "legacy_skill",
    "legacy_mob",
    "small_local",
    "client_hardcoded_ids",
}


def _fail(code: str, detail: str) -> None:
    raise ManifestError(f"{code}: {detail}")


def validate_manifest(
    manifest: dict[str, Any],
    *,
    base_dir: Path = Path("."),
    occupied_ids=(),
    previous_runtime_version=None,
) -> None:
    seen_ids = set()
    legacy_ids = set(occupied_ids)
    phase_range = manifest.get("phase", {}).get("sprite_id_range", [27041, 27082])
    if phase_range != [27041, 27082]:
        _fail("ASSET-ID-001", "Phase 1 phải khóa đúng dải 27041..27082")

    iruka_layers = []
    for sprite in manifest.get("sprites", []):
        sprite_id = sprite["id"]
        if sprite_id in seen_ids:
            _fail("ASSET-ID-001", f"sprite ID trùng: {sprite_id}")
        if not 27041 <= sprite_id <= 27082:
            _fail("ASSET-ID-001", f"sprite ID ngoài Phase 1 27041..27082: {sprite_id}")
        if sprite_id in legacy_ids:
            _fail("ASSET-ID-001", f"sprite ID va chạm baseline legacy: {sprite_id}")
        if set(sprite.get("zooms", ())) != {"x1", "x2", "x3", "x4"}:
            _fail("ASSET-ZOOM-003", f"sprite {sprite_id} phải có đủ x1..x4")

        source_hash = sprite.get("source_sha256", "")
        source_path = sprite.get("source_path")
        if (
            not source_path
            or source_hash == "0" * 64
            or len(source_hash) != 64
            or any(char not in "0123456789abcdef" for char in source_hash)
        ):
            _fail("ASSET-SPRITE-002", f"sprite {sprite_id} có source path/SHA không hợp lệ")
        try:
            source_digest = hashlib.sha256((base_dir / source_path).read_bytes()).hexdigest()
            x4_digest = hashlib.sha256((base_dir / sprite["zooms"]["x4"]).read_bytes()).hexdigest()
        except (OSError, KeyError) as error:
            _fail("ASSET-SPRITE-002", f"sprite {sprite_id} thiếu source/x4: {error}")
        if source_digest != source_hash or x4_digest != source_hash:
            _fail("ASSET-SPRITE-002", f"sprite {sprite_id} source/x4 SHA không khớp")

        crop = sprite.get("crop")
        if crop and crop.get("subject") == "iruka":
            rect, overlap, layer = crop.get("source_rect"), crop.get("overlap"), crop.get("layer")
            valid_rect = (
                isinstance(rect, list)
                and len(rect) == 4
                and all(isinstance(value, int) for value in rect)
                and rect[0] >= 0
                and rect[1] >= 0
                and rect[2] > 0
                and rect[3] > 0
            )
            valid_overlap = (
                isinstance(overlap, list)
                and len(overlap) == 2
                and all(isinstance(value, int) and -127 <= value <= 127 for value in overlap)
            )
            if not valid_rect or not valid_overlap or not isinstance(layer, int) or layer < 0:
                _fail("ASSET-CROP-007", "crop Iruka sai kiểu hoặc ngoài phạm vi")
            iruka_layers.append(layer)
        seen_ids.add(sprite_id)

    if iruka_layers and (
        len(iruka_layers) != len(set(iruka_layers)) or iruka_layers != sorted(iruka_layers)
    ):
        _fail("ASSET-CROP-007", "layer Iruka phải duy nhất và tăng dần")

    previous = previous_runtime_version
    if previous is None:
        previous = manifest.get("previous_runtime_slice_version")
    if previous is not None:
        runtime_version = manifest.get("versions", {}).get("runtime_slice")
        if runtime_version is None or runtime_version <= previous:
            _fail("ASSET-VERSION-006", "runtime_slice_version phải tăng")

    for part in manifest.get("parts", []):
        if not 0 <= part["id"] <= 32767:
            _fail("ASSET-PART-004", f"part ID ngoài signed short: {part['id']}")
        for frame in part.get("frames", []):
            if frame.get("sprite_id") not in seen_ids:
                _fail("ASSET-PART-004", f"part {part['id']} trỏ sprite không tồn tại")
            if not -128 <= frame.get("dx", 0) <= 127 or not -128 <= frame.get("dy", 0) <= 127:
                _fail("ASSET-PART-005", f"part {part['id']} có dx/dy ngoài sbyte")

    for animation in manifest.get("animations", []):
        frames = animation.get("frames", [])
        slot_count = animation.get("slot_count", len(frames))
        seen_slots = set()
        for frame in frames:
            slot = frame.get("slot_index")
            client_cf = frame.get("client_cf")
            if (
                not isinstance(slot, int)
                or not 0 <= slot < slot_count
                or slot in seen_slots
                or not isinstance(client_cf, list)
                or not client_cf
                or len(client_cf) != len(set(client_cf))
                or any(not isinstance(value, int) or not -128 <= value <= 127 for value in client_cf)
            ):
                _fail("ASSET-ANIMATION-008", f"animation {animation.get('name', '?')} có slot/cf sai")
            seen_slots.add(slot)
        if seen_slots != set(range(len(frames))):
            _fail("ASSET-ANIMATION-008", f"animation {animation.get('name', '?')} slot không liên tục")


def scan_legacy_inputs(root: Path, ledger: dict[str, Any]) -> dict[str, Any]:
    """Đọc đúng tập input đã khóa; không tự tìm cache thay thế."""
    kinds = {entry.get("kind") for entry in ledger.get("inputs", [])}
    if kinds != REQUIRED_LEGACY_KINDS:
        _fail("ASSET-LEGACY-009", "ledger thiếu hoặc thừa closure legacy bắt buộc")
    found = []
    for entry in ledger["inputs"]:
        if Path(entry["path"]).is_absolute():
            _fail("ASSET-LEGACY-009", f"path legacy phải tương đối: {entry['path']}")
        path = root / entry["path"]
        if not path.is_file():
            _fail("ASSET-LEGACY-009", f"thiếu input legacy: {entry['path']}")
        content = path.read_bytes()
        digest = hashlib.sha256(content).hexdigest()
        if digest != entry["sha256"]:
            _fail("ASSET-LEGACY-009", f"SHA legacy lệch: {entry['path']}")
        if len(content) != entry.get("bytes", len(content)):
            _fail("ASSET-LEGACY-009", f"kích thước legacy lệch: {entry['path']}")
        if entry["kind"] == "part.bin":
            if len(content) < 2:
                _fail("ASSET-LEGACY-009", f"part.bin quá ngắn: {entry['path']}")
            expected_records = ledger.get("expected", {}).get("part_records")
            actual_records = int.from_bytes(content[:2], byteorder="big")
            if expected_records is not None and actual_records != expected_records:
                _fail("ASSET-LEGACY-009", "số record part.bin lệch baseline")
        found.append({"path": entry["path"], "kind": entry["kind"], "sha256": digest})
    return {"inputs": found, "occupied_ids": ledger.get("expected", {}).get("occupied_ids", [])}


def verify_legacy_inputs(ledger_path: Path, root: Path) -> None:
    ledger = json.loads(ledger_path.read_text(encoding="utf-8"))
    scan_legacy_inputs(root, ledger)
