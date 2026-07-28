"""Layout golden nhỏ cho EffectData: action cố định, frame giữ nguyên thứ tự."""

import json


class MobError(ValueError):
    pass


_ACTIONS = ("idle", "run", "pounce", "hit", "death")
_FRAME_COUNTS = {"idle": 2, "run": 4, "pounce": 2, "hit": 1, "death": 1}


def encode_mob(mob_id: int, frames: dict[str, list[int]]) -> bytes:
    if not isinstance(mob_id, int) or set(frames) != set(_ACTIONS) or any(not values for values in frames.values()):
        raise MobError("ASSET-MOB-003: thiếu action EffectData")
    return json.dumps([mob_id, {key: list(frames[key]) for key in _ACTIONS}], separators=(",", ":")).encode("utf-8")


def parse_mob(data: bytes):
    mob_id, frames = json.loads(data.decode("utf-8"))
    if set(frames) != set(_ACTIONS):
        raise MobError("ASSET-MOB-003: golden EffectData sai action")
    return mob_id, frames


def validate_mob_golden(mob: dict) -> None:
    if not isinstance(mob, dict) or not isinstance(mob.get("frames"), dict):
        raise MobError("ASSET-MOB-003: thiếu golden frame/action")
    if set(mob["frames"]) != set(_ACTIONS) or any(not isinstance(value, list) or len(value) != _FRAME_COUNTS[action] for action, value in mob["frames"].items()):
        raise MobError("ASSET-MOB-003: golden thiếu idle/run/pounce/hit/death")
    if not isinstance(mob.get("bytes_sha256"), str) or len(mob["bytes_sha256"]) != 64:
        raise MobError("ASSET-MOB-003: thiếu SHA Mob bytes")
    if not isinstance(mob.get("data_sha256"), str) or len(mob["data_sha256"]) != 64:
        raise MobError("ASSET-MOB-003: thiếu SHA Mob data")
