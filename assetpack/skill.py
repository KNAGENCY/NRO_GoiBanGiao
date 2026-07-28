"""Skill binary độc lập: short count/id big-endian, 12 level và sentinel -1."""

import struct


class SkillError(ValueError):
    pass


def _is_int(value):
    return isinstance(value, int) and not isinstance(value, bool)


def _check(record):
    if (not _is_int(record.get("id")) or not _is_int(record.get("pose_slot"))
            or record.get("effect_id") == 80 or len(record.get("levels", [])) != 12
            or any(not _is_int(level) or not 0 <= level <= 255 for level in record["levels"])
            or record.get("effect_id") not in range(-1, 80) or record.get("arrow_id") not in range(-1, 128)):
        raise SkillError("ASSET-SKILL-004: 12 level/pose/effect/arrow không hợp lệ")


def encode_skill(records) -> bytes:
    records = list(records)
    if len(records) > 32767:
        raise SkillError("ASSET-SKILL-004: quá nhiều skill")
    out = bytearray(struct.pack(">H", len(records)))
    for record in records:
        _check(record)
        out.extend(struct.pack(">Hbb", record["id"], record["pose_slot"], record["effect_id"]))
        out.extend(bytes(record["levels"]))
        out.extend(struct.pack("b", record["arrow_id"]))
    return bytes(out)


def parse_skill(data: bytes):
    if len(data) < 2:
        raise SkillError("ASSET-SKILL-004: thiếu count")
    count = struct.unpack_from(">H", data)[0]
    offset, records = 2, []
    for _ in range(count):
        if offset + 17 > len(data):
            raise SkillError("ASSET-SKILL-004: count/endian sai")
        skill_id, pose, effect = struct.unpack_from(">Hbb", data, offset)
        levels = list(data[offset + 4:offset + 16])
        arrow = struct.unpack_from("b", data, offset + 16)[0]
        record = {"id": skill_id, "pose_slot": pose, "levels": levels, "effect_id": effect, "arrow_id": arrow}
        _check(record)
        records.append(record); offset += 17
    if offset != len(data):
        raise SkillError("ASSET-SKILL-004: byte thừa")
    return tuple(records)
