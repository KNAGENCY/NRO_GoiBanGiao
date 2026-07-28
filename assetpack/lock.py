"""Khoá provenance canonical, dùng byte này làm input SHA cho deploy."""

import json


def lock_bytes(lock: dict) -> bytes:
    canonical = dict(lock)
    for key in ("inputs", "outputs", "evidence"):
        if key in canonical:
            canonical[key] = sorted(canonical[key], key=lambda item: item["path"])
    return (json.dumps(canonical, ensure_ascii=False, sort_keys=True, separators=(",", ":")) + "\n").encode("utf-8")
