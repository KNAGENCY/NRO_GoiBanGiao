"""Cổng build: review phải PASS và evidence/source SHA phải khớp trước bind."""

import json
import hashlib
from pathlib import Path

from .catalog import encode_catalog
from .lock import lock_bytes


class BuildError(ValueError):
    pass


def _entries(review):
    return review.get("assets", review if isinstance(review, list) else [])


def _review_documents(review_path: Path):
    try:
        paths = sorted(review_path.glob("*.json")) if review_path.is_dir() else [review_path]
        if not paths:
            raise OSError("thư mục review rỗng")
        return [(path, json.loads(path.read_text(encoding="utf-8"))) for path in paths]
    except (OSError, json.JSONDecodeError) as error:
        raise BuildError(f"ASSET-REVIEW-010: thiếu review JSON: {error}") from error


def review_digest(review_path: Path) -> str:
    digest = hashlib.sha256()
    for path, _ in _review_documents(review_path):
        digest.update(path.name.encode("utf-8")); digest.update(path.read_bytes())
    return digest.hexdigest()


def require_approved_review(review_path: Path) -> None:
    entries = [entry for _, review in _review_documents(review_path) for entry in _entries(review)]
    if not entries or any(entry.get("status") != "PASS" for entry in entries):
        raise BuildError("ASSET-REVIEW-010: source chưa được review PASS")
    for entry in entries:
        if not entry.get("source_sha256") or not entry.get("evidence_sha256") or not entry.get("evidence_path"):
            raise BuildError("ASSET-REVIEW-010: thiếu SHA/evidence review")
        evidence = Path(entry["evidence_path"])
        if not evidence.is_file() or hashlib.sha256(evidence.read_bytes()).hexdigest() != entry["evidence_sha256"]:
            raise BuildError("ASSET-REVIEW-010: evidence SHA không khớp")


def build(manifest: dict, out: Path, review_path: Path) -> None:
    require_approved_review(review_path)
    # Việc tạo output chỉ được phép sau gate; Task 6B sẽ truyền manifest đã bind PASS.
    out.mkdir(parents=True, exist_ok=True)
    catalog_path = out / "catalog.json"
    catalog_path.write_bytes(encode_catalog(manifest.get("catalog", {})))
    outputs = []
    for path in sorted(out.rglob("*")):
        if path.is_file():
            outputs.append({"path": path.relative_to(out).as_posix(), "sha256": hashlib.sha256(path.read_bytes()).hexdigest(), "bytes": path.stat().st_size})
    lock = {
        "slice": manifest.get("slice"),
        "versions": manifest.get("versions", {}),
        "review_sha256": review_digest(review_path),
        "inputs": [{"path": path, "sha256": hashlib.sha256(Path(path).read_bytes()).hexdigest()} for path in sorted(manifest.get("review_sources", []))],
        "outputs": outputs,
    }
    (out / "asset-lock.json").write_bytes(lock_bytes(lock))
