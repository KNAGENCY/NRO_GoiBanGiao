"""Collision server/client có cùng nghĩa tile, không chỉ cùng kích thước byte."""

import hashlib
import struct


class MapError(ValueError):
    pass


def encode_collision(width: int, height: int, indexes, *, spawn_index: int) -> bytes:
    indexes = list(indexes)
    if not 1 <= width <= 255 or not 1 <= height <= 255 or len(indexes) != width * height:
        raise MapError("ASSET-COLLISION-006: header hoặc số cell collision sai")
    if not 0 <= spawn_index < len(indexes) or any(not isinstance(value, int) or not 0 <= value <= 99 for value in indexes):
        raise MapError("ASSET-COLLISION-006: index tile hoặc spawn không hợp lệ")
    return struct.pack(">BB", width, height) + bytes(indexes)


def semantic_hash(solid_values) -> str:
    return hashlib.sha256(bytes(1 if value else 0 for value in solid_values)).hexdigest()


def validate_collision_semantics(tiles, indexes, client_solid) -> None:
    by_id = {entry["id"]: bool(entry["solid"]) for entry in tiles}
    server = [by_id[index] for index in indexes]
    if len(server) != len(client_solid) or semantic_hash(server) != semantic_hash(client_solid):
        raise MapError("ASSET-COLLISION-006: collision client/server khác semantic")
