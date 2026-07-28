"""Schema typed và các kiểm tra không cần sinh artifact ở Task 1."""

from dataclasses import dataclass
import hashlib
import json
from pathlib import Path
from typing import Any


class ManifestError(ValueError):
    """Lỗi contract có mã ổn định để CI và người làm asset cùng đọc."""


@dataclass(frozen=True)
class Sprite:
    id: int
    zooms: dict[str, str]
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


def _fail(code: str, detail: str) -> None:
    raise ManifestError(f"{code}: {detail}")


def validate_manifest(manifest):
    seen_ids = set()
    for sprite in manifest.get("sprites", []):
        sprite_id = sprite["id"]
        if sprite_id in seen_ids:
            _fail("ASSET-ID-001", f"sprite ID trùng: {sprite_id}")
        if not 0 <= sprite_id <= 32767:
            _fail("ASSET-ID-001", f"sprite ID ngoài signed short: {sprite_id}")
        if sprite_id in set(manifest.get("legacy_occupied_ids", [])):
            _fail("ASSET-ID-001", f"sprite ID va chạm baseline legacy: {sprite_id}")
        if sprite_id in set(manifest.get("phase", {}).get("forbidden_poc_ids", [])):
            _fail("ASSET-ID-001", f"sprite ID POC bị cấm trong manifest chính: {sprite_id}")
        if set(sprite.get("zooms", ())) != {"x1", "x2", "x3", "x4"}:
            _fail("ASSET-ZOOM-003", f"sprite {sprite_id} phải có đủ x1..x4")
        source_hash = sprite.get("source_sha256", "")
        if len(source_hash) != 64 or any(char not in "0123456789abcdef" for char in source_hash):
            _fail("ASSET-SPRITE-002", f"sprite {sprite_id} có source SHA-256 không hợp lệ")
        crop = sprite.get("crop")
        if crop and crop.get("subject") == "iruka":
            if "overlap" not in crop or "layer" not in crop:
                _fail("ASSET-CROP-007", "crop Iruka phải khóa overlap và layer")
        seen_ids.add(sprite_id)

    if "previous_runtime_slice_version" in manifest:
        runtime_version = manifest.get("versions", {}).get("runtime_slice")
        if runtime_version is None or runtime_version <= manifest["previous_runtime_slice_version"]:
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
        for frame in animation.get("frames", []):
            if "slot_index" not in frame or not frame.get("client_cf"):
                _fail("ASSET-ANIMATION-008", f"animation {animation.get('name', '?')} thiếu slot/cf")


def scan_legacy_inputs(root: Path, ledger: dict[str, Any]) -> dict[str, Any]:
    """Đọc đúng tập input đã khóa; không tự tìm cache thay thế khi thiếu file."""
    found = []
    for entry in ledger.get("inputs", []):
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
