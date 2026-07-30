"""Catalog runtime canonical, đủ dữ liệu để server không suy diễn ID theo vị trí."""

import json


def encode_catalog(catalog: dict) -> bytes:
    return json.dumps(catalog, ensure_ascii=False, sort_keys=True, separators=(",", ":")).encode("utf-8")


def _integer(value) -> bool:
    return isinstance(value, int) and not isinstance(value, bool)


def _point(value) -> bool:
    return isinstance(value, dict) and _integer(value.get("x")) and _integer(value.get("y"))


def validate_catalog(catalog: dict) -> None:
    """Chỉ chấp nhận catalog có binding runtime tường minh, không fallback theo index."""
    if not isinstance(catalog, dict):
        raise ValueError("ASSET-CATALOG-009: catalog phải là object")
    maps = catalog.get("maps")
    if not isinstance(maps, list) or len(maps) != 6:
        raise ValueError("ASSET-CATALOG-009: catalog phải có sáu map")
    map_ids = set()
    for item in maps:
        if not isinstance(item, dict) or not _integer(item.get("id")) or not isinstance(item.get("name"), str) or not item["name"] or not _point(item.get("spawn")):
            raise ValueError("ASSET-CATALOG-009: map thiếu id/name/spawn")
        map_ids.add(item["id"])
    if len(map_ids) != len(maps):
        raise ValueError("ASSET-CATALOG-009: map id trùng")

    npcs = catalog.get("npcs")
    if not isinstance(npcs, list) or len(npcs) != 1:
        raise ValueError("ASSET-CATALOG-009: phải có NPC Iruka")
    npc = npcs[0]
    npc_fields = ("id", "head", "body", "leg", "portrait", "x", "y")
    if (not isinstance(npc, dict) or npc.get("name") != "Iruka" or
            any(not _integer(npc.get(key)) for key in npc_fields) or
            not isinstance(npc.get("dialog"), str) or not npc["dialog"] or
            not isinstance(npc.get("menu"), list) or not npc["menu"] or
            any(not isinstance(menu, str) or not menu for menu in npc["menu"])):
        raise ValueError("ASSET-CATALOG-009: NPC thiếu binding runtime")

    mobs = catalog.get("mobs")
    if not isinstance(mobs, list) or len(mobs) != 1:
        raise ValueError("ASSET-CATALOG-009: phải có mob Sói Rừng")
    mob = mobs[0]
    mob_fields = ("id", "hp", "range", "speed", "dart")
    if (not isinstance(mob, dict) or mob.get("name") != "Sói Rừng" or
            any(not _integer(mob.get(key)) for key in mob_fields) or mob["hp"] <= 0 or
            not isinstance(mob.get("spawns"), list) or len(mob["spawns"]) < 1):
        raise ValueError("ASSET-CATALOG-009: mob thiếu binding runtime")
    spawn_ids = set()
    for spawn in mob["spawns"]:
        if not isinstance(spawn, dict) or not _integer(spawn.get("id")) or not _integer(spawn.get("x")) or not _integer(spawn.get("y")):
            raise ValueError("ASSET-CATALOG-009: mob spawn sai")
        spawn_ids.add(spawn["id"])
    if len(spawn_ids) != len(mob["spawns"]):
        raise ValueError("ASSET-CATALOG-009: mob spawn id trùng")

    poses = catalog.get("skill_poses")
    if not isinstance(poses, list) or len(poses) != 4 or set(poses) != {0, 1, 2, 3}:
        raise ValueError("ASSET-CATALOG-009: thiếu bốn pose skill")
