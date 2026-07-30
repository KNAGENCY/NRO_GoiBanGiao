"""Collision server/client có cùng nghĩa tile, không chỉ cùng kích thước byte."""

import hashlib
import struct


class MapError(ValueError):
    pass


def validate_tileset_baseline(data: bytes, baseline: dict) -> None:
    """Kiểm baseline bằng bytes vàng và parse semantic tối thiểu độc lập."""
    if not isinstance(data, bytes) or not data:
        raise MapError("ASSET-TILESET-010: tile_set_info rỗng/malformed")
    if not isinstance(baseline, dict):
        raise MapError("ASSET-TILESET-010: thiếu tile-set baseline")
    try:
        expected_bytes = bytes.fromhex(baseline["bytes_hex"])
        expected_sha = baseline["sha256"]
        semantic = baseline["semantic"]
    except (KeyError, TypeError, ValueError) as error:
        raise MapError(f"ASSET-TILESET-010: baseline malformed: {error}") from error
    if not isinstance(expected_sha, str) or len(expected_sha) != 64 or len(expected_bytes) == 0:
        raise MapError("ASSET-TILESET-010: baseline golden không hợp lệ")
    if len(data) != len(expected_bytes) or data != expected_bytes:
        raise MapError("ASSET-TILESET-011: tile_set_info khác baseline golden")
    if hashlib.sha256(data).hexdigest() != expected_sha:
        raise MapError("ASSET-TILESET-011: SHA tile_set_info khác baseline golden")
    if not isinstance(semantic, dict) or semantic.get("byte_length") != len(data) or semantic.get("tile_count") != data[0]:
        raise MapError("ASSET-TILESET-011: semantic tile_set_info khác baseline")


def encode_collision(width: int, height: int, solid_cells, *, spawn_index: int, tile_size) -> bytes:
    solid_cells = list(solid_cells)
    if not 1 <= width <= 255 or not 1 <= height <= 255 or len(solid_cells) != width * height:
        raise MapError("ASSET-COLLISION-006: header hoặc số cell collision sai")
    if (not isinstance(tile_size, (list, tuple)) or len(tile_size) != 2
            or any(not isinstance(value, int) or isinstance(value, bool) or not 1 <= value <= 255 for value in tile_size)):
        raise MapError("ASSET-COLLISION-006: tile size collision sai")
    if not isinstance(spawn_index, int) or isinstance(spawn_index, bool) or not 0 <= spawn_index < len(solid_cells) or any(not isinstance(value, bool) for value in solid_cells):
        raise MapError("ASSET-COLLISION-006: semantic solid hoặc spawn không hợp lệ")
    if solid_cells[spawn_index] or not any(solid_cells):
        raise MapError("ASSET-COLLISION-006: spawn phải walkable và map phải có solid")
    return struct.pack(">BBBB", width, height, tile_size[0], tile_size[1]) + bytes(1 if value else 0 for value in solid_cells)


def semantic_hash(solid_values) -> str:
    return hashlib.sha256(bytes(1 if value else 0 for value in solid_values)).hexdigest()


def validate_collision_semantics(tiles, indexes, client_solid) -> None:
    by_id = {}
    for entry in tiles:
        if not isinstance(entry, dict) or not isinstance(entry.get("id"), int) or isinstance(entry.get("id"), bool) or not isinstance(entry.get("solid"), bool):
            raise MapError("ASSET-COLLISION-006: tile semantic malformed")
        by_id[entry["id"]] = entry["solid"]
    try:
        server = [by_id[index] for index in indexes]
    except (KeyError, TypeError):
        raise MapError("ASSET-COLLISION-006: index tile không có semantic")
    if len(server) != len(client_solid) or any(not isinstance(value, bool) for value in client_solid) or semantic_hash(server) != semantic_hash(client_solid):
        raise MapError("ASSET-COLLISION-006: collision client/server khác semantic")
