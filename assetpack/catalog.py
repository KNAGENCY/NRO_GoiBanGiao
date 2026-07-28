"""Catalog JSON canonical để lock không đổi theo thứ tự dict của Python."""

import json


def encode_catalog(catalog: dict) -> bytes:
    return json.dumps(catalog, ensure_ascii=False, sort_keys=True, separators=(",", ":")).encode("utf-8")
