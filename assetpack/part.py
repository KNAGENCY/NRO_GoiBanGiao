"""Đọc và sinh part.bin theo đúng layout nhị phân mà client Unity tiêu thụ."""

from dataclasses import dataclass
import struct


FRAME_COUNTS = {0: 3, 1: 17, 2: 14}
PHASE_ONE_REPLACEMENTS = {30, 57, 58}
PHASE_ONE_ADDITIONS = {2189, 2190, 2191}


class PartError(ValueError):
    """Lỗi artifact part hoặc tham chiếu sprite trong manifest."""


@dataclass(frozen=True)
class Frame:
    sprite_id: int
    dx: int
    dy: int


@dataclass(frozen=True)
class Part:
    part_type: int
    frames: tuple[Frame, ...]


def _part_fail(message: str):
    raise PartError(f"ASSET-PART-001: {message}")


def _is_int(value):
    return isinstance(value, int) and not isinstance(value, bool)


def _validate_part(part: Part):
    expected_frames = FRAME_COUNTS.get(part.part_type)
    if expected_frames is None:
        _part_fail(f"type part không hợp lệ: {part.part_type}")
    if len(part.frames) != expected_frames:
        _part_fail(
            f"type {part.part_type} cần {expected_frames} frame, nhận {len(part.frames)}"
        )
    for frame in part.frames:
        if not _is_int(frame.sprite_id) or not 0 <= frame.sprite_id <= 32767:
            _part_fail(f"sprite ID ngoài signed short: {frame.sprite_id!r}")
        if not _is_int(frame.dx) or not -128 <= frame.dx <= 127:
            _part_fail(f"dx ngoài signed byte: {frame.dx!r}")
        if not _is_int(frame.dy) or not -128 <= frame.dy <= 127:
            _part_fail(f"dy ngoài signed byte: {frame.dy!r}")


def parse_part(data: bytes) -> tuple[Part, ...]:
    """Parse toàn bộ file, từ chối cả byte dư để endian sai không thể lọt im lặng."""
    if len(data) < 2:
        _part_fail("thiếu short số record")
    record_count = struct.unpack_from(">h", data, 0)[0]
    if record_count < 0:
        _part_fail("số record âm")
    offset = 2
    parts = []
    for index in range(record_count):
        if offset >= len(data):
            _part_fail(f"thiếu type tại record {index}")
        part_type = data[offset]
        offset += 1
        frame_count = FRAME_COUNTS.get(part_type)
        if frame_count is None:
            _part_fail(f"type {part_type} không hợp lệ tại record {index}")
        byte_count = frame_count * 4
        if offset + byte_count > len(data):
            _part_fail(f"thiếu frame tại record {index}")
        frames = tuple(
            Frame(*struct.unpack_from(">hbb", data, offset + frame_index * 4))
            for frame_index in range(frame_count)
        )
        offset += byte_count
        parts.append(Part(part_type, frames))
    if offset != len(data):
        _part_fail(f"dư {len(data) - offset} byte sau record cuối")
    return tuple(parts)


def encode_part(parts) -> bytes:
    """Encode có kiểm tra; format bắt buộc big-endian như DataInputStream của client."""
    parts = tuple(parts)
    if len(parts) > 32767:
        _part_fail(f"quá nhiều record: {len(parts)}")
    output = bytearray(struct.pack(">h", len(parts)))
    for part in parts:
        _validate_part(part)
        output.append(part.part_type)
        for frame in part.frames:
            output.extend(struct.pack(">hbb", frame.sprite_id, frame.dx, frame.dy))
    return bytes(output)


def _manifest_part(entry, sprite_ids) -> Part:
    part_id = entry.get("id")
    part_type = entry.get("type")
    frames = entry.get("frames")
    if not _is_int(part_id) or part_id < 0:
        _part_fail(f"part ID không hợp lệ: {part_id!r}")
    if not _is_int(part_type) or not isinstance(frames, list):
        _part_fail(f"part {part_id} thiếu type hoặc frames")
    parsed_frames = []
    for frame in frames:
        if not isinstance(frame, dict):
            _part_fail(f"frame part {part_id} không phải object")
        sprite_id = frame.get("sprite_id")
        if sprite_id not in sprite_ids:
            raise PartError(f"ASSET-SPRITE-002: part {part_id} trỏ sprite không tồn tại")
        parsed_frames.append(Frame(sprite_id, frame.get("dx"), frame.get("dy")))
    part = Part(part_type, tuple(parsed_frames))
    _validate_part(part)
    return part


def build_part(manifest: dict, legacy_part: bytes) -> bytes:
    """Giữ nguyên legacy trừ 30/57/58 và chỉ nối các part mới sau record 2188."""
    legacy = list(parse_part(legacy_part))
    sprite_ids = {sprite.get("id") for sprite in manifest.get("sprites", []) if isinstance(sprite, dict)}
    supplied = manifest.get("parts", [])
    if not isinstance(supplied, list):
        _part_fail("parts manifest không phải danh sách")
    replacements = {}
    additions = {}
    for entry in supplied:
        if not isinstance(entry, dict):
            _part_fail("part manifest không phải object")
        part_id = entry.get("id")
        part = _manifest_part(entry, sprite_ids)
        if part_id < len(legacy):
            if part_id not in PHASE_ONE_REPLACEMENTS:
                _part_fail(f"Phase 1 không được thay part legacy {part_id}")
            if part_id in replacements:
                _part_fail(f"part {part_id} lặp trong manifest")
            replacements[part_id] = part
        else:
            if part_id in additions:
                _part_fail(f"part {part_id} lặp trong manifest")
            additions[part_id] = part
    if set(replacements) != PHASE_ONE_REPLACEMENTS:
        _part_fail("Phase 1 phải thay đủ part 30/57/58")
    if set(additions) != PHASE_ONE_ADDITIONS:
        _part_fail("part nối mới phải liền sau legacy")
    for part_id, part in replacements.items():
        legacy[part_id] = part
    return encode_part([*legacy, *(additions[part_id] for part_id in sorted(additions))])
