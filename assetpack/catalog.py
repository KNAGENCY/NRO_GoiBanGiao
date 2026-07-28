"""Catalog JSON canonical để lock không đổi theo thứ tự dict của Python."""

import json


def encode_catalog(catalog: dict) -> bytes:
    return json.dumps(catalog, ensure_ascii=False, sort_keys=True, separators=(",", ":")).encode("utf-8")


def validate_catalog(catalog: dict) -> None:
    """Catalog runtime không được rỗng hoặc lấy lại danh sách từ code test."""
    if not isinstance(catalog, dict):
        raise ValueError("ASSET-CATALOG-009: catalog phải là object")
    maps = catalog.get("maps")
    if not isinstance(maps, list) or len(maps) != 6 or any(not isinstance(name, str) or not name for name in maps):
        raise ValueError("ASSET-CATALOG-009: catalog phải có sáu map name")
    if not isinstance(catalog.get("npc_templates"), list) or "Iruka" not in catalog["npc_templates"]:
        raise ValueError("ASSET-CATALOG-009: thiếu template Iruka")
    if not isinstance(catalog.get("mob_templates"), list) or "Sói Rừng" not in catalog["mob_templates"]:
        raise ValueError("ASSET-CATALOG-009: thiếu template Sói Rừng")
    poses = catalog.get("skill_poses")
    if not isinstance(poses, list) or len(poses) != 4 or set(poses) != {0, 1, 2, 3}:
        raise ValueError("ASSET-CATALOG-009: thiếu bốn pose skill")
    if not isinstance(catalog.get("spawn"), list) or not catalog["spawn"]:
        raise ValueError("ASSET-CATALOG-009: thiếu spawn")
