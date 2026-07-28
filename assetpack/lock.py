"""Khoá provenance canonical, dùng byte này làm input SHA cho deploy."""

import json


def lock_bytes(lock: dict) -> bytes:
    canonical = dict(lock)
    if "outputs" in canonical:
        canonical["outputs"] = sorted(canonical["outputs"], key=lambda output: output["path"])
    return (json.dumps(canonical, ensure_ascii=False, sort_keys=True, separators=(",", ":")) + "\n").encode("utf-8")
