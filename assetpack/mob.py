"""Layout golden nhỏ cho EffectData: action cố định, frame giữ nguyên thứ tự."""

import json


class MobError(ValueError):
    pass


_ACTIONS = ("idle", "run", "pounce", "hit", "death")


def encode_mob(mob_id: int, frames: dict[str, list[int]]) -> bytes:
    if not isinstance(mob_id, int) or set(frames) != set(_ACTIONS) or any(not values for values in frames.values()):
        raise MobError("ASSET-MOB-003: thiếu action EffectData")
    return json.dumps([mob_id, {key: list(frames[key]) for key in _ACTIONS}], separators=(",", ":")).encode("utf-8")


def parse_mob(data: bytes):
    mob_id, frames = json.loads(data.decode("utf-8"))
    if set(frames) != set(_ACTIONS):
        raise MobError("ASSET-MOB-003: golden EffectData sai action")
    return mob_id, frames
