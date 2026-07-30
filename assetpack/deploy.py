"""Cổng phát hành cặp artifact server/client đã bị khóa theo SHA.

Build chỉ chứng minh artifact được sinh đúng. Module này chứng minh hai nửa đem
đi chạy là cùng một cặp, phiên bản đã tăng khi byte runtime đổi, và client chỉ
được chép vào Unity snapshot có overlay.
"""

from __future__ import annotations

import hashlib
import json
from pathlib import Path
import shutil
import subprocess
import tempfile
from typing import Any, Callable


class DeployError(ValueError):
    """Lỗi deploy có mã ổn định để test mutation bắt đúng nguyên nhân."""


Preflight = Callable[[Path, Path, Path, Path], int]

_REQUIRED_SERVER = {
    "server/part.bin",
    "server/skill.bin",
    "server/tile_set_info.bin",
    "server/map-1-collision.bin",
    "server/catalog.json",
}
_VERSION_KEYS = {"runtime_slice", "data", "map", "resource"}


def _fail(code: str, detail: str) -> None:
    raise DeployError(f"{code}: {detail}")


def _sha_bytes(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def _sha(path: Path) -> str:
    return _sha_bytes(path.read_bytes())


def _read_json(path: Path, code: str) -> dict[str, Any]:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as error:
        _fail(code, f"không đọc được {path}: {error}")
    if not isinstance(value, dict):
        _fail(code, f"{path} phải chứa JSON object")
    return value


def _inside(candidate: Path, root: Path) -> bool:
    try:
        candidate.resolve().relative_to(root.resolve())
        return True
    except ValueError:
        return False


def _lock_outputs(lock: dict[str, Any]) -> dict[str, dict[str, Any]]:
    raw = lock.get("outputs")
    if not isinstance(raw, list):
        _fail("ASSET-RUNTIME-SHA-018", "asset-lock thiếu outputs[]")
    outputs: dict[str, dict[str, Any]] = {}
    for item in raw:
        if not isinstance(item, dict):
            _fail("ASSET-RUNTIME-SHA-018", "output trong lock không phải object")
        relative = item.get("path")
        sha256 = item.get("sha256")
        size = item.get("bytes")
        if (
            not isinstance(relative, str)
            or not relative
            or "\\" in relative
            or Path(relative).is_absolute()
            or ".." in Path(relative).parts
            or not isinstance(sha256, str)
            or len(sha256) != 64
            or not isinstance(size, int)
            or isinstance(size, bool)
            or size < 0
            or relative in outputs
        ):
            _fail("ASSET-RUNTIME-SHA-018", f"output lock không hợp lệ: {relative!r}")
        outputs[relative] = item
    missing = sorted(_REQUIRED_SERVER - set(outputs))
    if missing:
        _fail("ASSET-RUNTIME-SHA-018", f"thiếu runtime server: {missing}")
    return outputs


def _verify_versions(lock: dict[str, Any]) -> dict[str, int]:
    versions = lock.get("versions")
    if not isinstance(versions, dict) or set(versions) != _VERSION_KEYS:
        _fail("ASSET-VERSION-015", "lock phải đủ runtime_slice/data/map/resource")
    if any(
        not isinstance(value, int)
        or isinstance(value, bool)
        or value < 1
        or value > 255
        for value in versions.values()
    ):
        _fail("ASSET-VERSION-015", "version runtime phải là byte 1..255")
    return versions


def _read_artifact_lock(root: Path) -> tuple[dict[str, Any], bytes]:
    path = root / "asset-lock.json"
    try:
        raw = path.read_bytes()
        lock = json.loads(raw.decode("utf-8"))
    except (OSError, UnicodeDecodeError, json.JSONDecodeError) as error:
        _fail("ASSET-RUNTIME-SHA-018", f"asset-lock không đọc được: {error}")
    if not isinstance(lock, dict):
        _fail("ASSET-RUNTIME-SHA-018", "asset-lock phải là object")
    _lock_outputs(lock)
    _verify_versions(lock)
    return lock, raw


def _artifact_file(
    relative: str,
    server_root: Path,
    client_root: Path,
) -> Path:
    if relative.startswith("server/"):
        root = server_root
    elif relative.startswith("client/"):
        root = client_root
    else:
        _fail("ASSET-RUNTIME-SHA-018", f"output ngoài server/client: {relative}")
    path = (root / relative).resolve()
    if not _inside(path, root):
        _fail("ASSET-RUNTIME-SHA-018", f"output vượt artifact root: {relative}")
    return path


def _verify_outputs(
    lock: dict[str, Any],
    server_root: Path,
    client_root: Path,
) -> dict[str, dict[str, Any]]:
    outputs = _lock_outputs(lock)
    pair = hashlib.sha256()
    # build._pair_sha sắp Path theo quy tắc của Windows (không phân biệt hoa
    # thường), nên phải giữ đúng oracle đó thay vì sort chuỗi ASCII.
    ordered = sorted(
        (path for path in outputs if path.startswith("server/")),
        key=Path,
    )
    ordered += sorted(
        (path for path in outputs if path.startswith("client/")),
        key=Path,
    )
    for relative in ordered:
        path = _artifact_file(relative, server_root, client_root)
        if not path.is_file():
            _fail("ASSET-RUNTIME-SHA-018", f"thiếu output: {relative}")
        data = path.read_bytes()
        record = outputs[relative]
        if len(data) != record["bytes"] or _sha_bytes(data) != record["sha256"]:
            _fail("ASSET-RUNTIME-SHA-018", f"SHA/length lệch: {relative}")
        pair.update(relative.encode("utf-8"))
        pair.update(data)
    expected_pair = lock.get("server_client_pair_sha256")
    if pair.hexdigest() != expected_pair:
        _fail("ASSET-LOCK-MISMATCH-017", "pair SHA server/client không khớp lock")
    return outputs


def _ledger_versions(path: Path) -> dict[str, int]:
    payload = _read_json(path, "ASSET-VERSION-015")
    mapping = {
        "runtime_slice": payload.get("runtime_slice_version"),
        "data": payload.get("data_version"),
        "map": payload.get("map_version"),
        "resource": payload.get("resource_version"),
    }
    if any(not isinstance(value, int) or isinstance(value, bool) for value in mapping.values()):
        _fail("ASSET-VERSION-015", "version ledger thiếu bốn mốc số nguyên")
    return mapping


def _category(relative: str) -> str:
    if relative == "server/part.bin" or "/smallimage/" in relative:
        return "data"
    if relative == "server/skill.bin" or "/mainImage/" in relative:
        return "resource"
    return "map"


def _validate_version_bump(
    lock: dict[str, Any],
    ledger_path: Path,
    previous_lock_path: Path | None,
) -> None:
    current = _verify_versions(lock)
    ledger = _ledger_versions(ledger_path)
    if previous_lock_path is None:
        for key in sorted(_VERSION_KEYS):
            if current[key] <= ledger[key]:
                _fail(
                    "ASSET-VERSION-015",
                    f"versions.{key}={current[key]} chưa tăng qua ledger {ledger[key]}",
                )
        return

    previous = _read_json(previous_lock_path, "ASSET-VERSION-015")
    previous_versions = _verify_versions(previous)
    old_outputs = {
        path: item["sha256"] for path, item in _lock_outputs(previous).items()
    }
    new_outputs = {
        path: item["sha256"] for path, item in _lock_outputs(lock).items()
    }
    changed = {
        path
        for path in set(old_outputs) | set(new_outputs)
        if old_outputs.get(path) != new_outputs.get(path)
    }
    for key in _VERSION_KEYS:
        if current[key] < previous_versions[key]:
            _fail("ASSET-VERSION-015", f"versions.{key} bị hạ")
    if not changed:
        return
    if current["runtime_slice"] <= previous_versions["runtime_slice"]:
        _fail("ASSET-VERSION-015", "runtime đổi byte nhưng runtime_slice không tăng")
    for category in {_category(path) for path in changed}:
        if current[category] <= previous_versions[category]:
            _fail(
                "ASSET-VERSION-015",
                f"runtime {category} đổi byte nhưng version không tăng",
            )


def _client_destination(snapshot: Path, relative: str) -> Path | None:
    if not relative.startswith("client/"):
        return None
    client_relative = relative.removeprefix("client/")
    first = client_relative.split("/", 1)[0]
    if first not in {"x1", "x2", "x3", "x4", "mymap", "Mob"}:
        return None
    return snapshot / "Assets" / "Resources" / "res" / client_relative


def _to_physical(relative: str) -> str:
    return "Assets/Resources/res/" + relative.removeprefix("client/")


def _resource(
    outputs: dict[str, dict[str, Any]],
    kind: str,
    name: str,
    paths: list[str],
) -> dict[str, Any]:
    zooms = []
    for zoom, relative in enumerate(paths, 1):
        record = outputs.get(relative)
        if (
            record is None
            or not relative.startswith(f"client/x{zoom}/")
            or not isinstance(record.get("width"), int)
            or isinstance(record.get("width"), bool)
            or record["width"] <= 0
            or not isinstance(record.get("height"), int)
            or isinstance(record.get("height"), bool)
            or record["height"] <= 0
        ):
            _fail(
                "ASSET-PREFLIGHT-022",
                f"metadata width/height hoac zoom sai: {relative}",
            )
        width, height = record["width"], record["height"]
        zooms.append(
            {
                "zoom": zoom,
                "physical_path": _to_physical(relative),
                "expected_width": width,
                "expected_height": height,
            }
        )
    return {"kind": kind, "name": name, "zooms": zooms}


def _find_zoom_paths(
    outputs: dict[str, dict[str, Any]],
    predicate: Callable[[str], bool],
) -> list[str]:
    paths = []
    for zoom in range(1, 5):
        matches = sorted(
            path
            for path in outputs
            if path.startswith(f"client/x{zoom}/") and predicate(path)
        )
        if not matches:
            _fail("ASSET-PREFLIGHT-022", f"thiếu resource {predicate!r} ở x{zoom}")
        paths.append(matches[0])
    return paths


def _exact_zoom_paths(
    outputs: dict[str, dict[str, Any]],
    relative: str,
) -> list[str]:
    paths = []
    for zoom in range(1, 5):
        expected = f"client/x{zoom}/{relative}"
        if expected not in outputs:
            _fail("ASSET-PREFLIGHT-022", f"thieu resource canonical: {expected}")
        paths.append(expected)
    return paths


def _canonical_closure_id(lock: dict[str, Any], key: str) -> int:
    closure = lock.get("closure")
    values = closure.get(key) if isinstance(closure, dict) else None
    if (
        not isinstance(values, list)
        or not values
        or any(not isinstance(value, int) or isinstance(value, bool) for value in values)
    ):
        _fail("ASSET-PREFLIGHT-022", f"closure thieu {key} canonical")
    return min(values)


def _resource_contract(
    lock: dict[str, Any],
    outputs: dict[str, dict[str, Any]],
) -> dict[str, Any]:
    sprite_id = _canonical_closure_id(lock, "sprite_ids")
    skill_id = _canonical_closure_id(lock, "skill_ids")
    mob_id = _canonical_closure_id(lock, "mob_ids")
    resources = [
        _resource(
            outputs,
            "small",
            f"iruka-player-small-{sprite_id}",
            _exact_zoom_paths(outputs, f"smallimage/Small{sprite_id}.png"),
        ),
        _resource(
            outputs,
            "background",
            "lang-la-background",
            _exact_zoom_paths(outputs, "bg/b00.png"),
        ),
        _resource(
            outputs,
            "mob",
            f"soi-rung-{mob_id}",
            _exact_zoom_paths(outputs, f"Mob/{mob_id}/img.png"),
        ),
        _resource(
            outputs,
            "hud",
            "naruto-hud",
            _exact_zoom_paths(outputs, "mainImage/naruto-hp-bg.png"),
        ),
        _resource(
            outputs,
            "skill",
            f"skill-pose-{skill_id}",
            _exact_zoom_paths(outputs, f"smallimage/Small{skill_id}.png"),
        ),
    ]
    return {
        "schema_version": 1,
        "tile_id": 1,
        "tile_count": 100,
        "resources": resources,
    }


def _validate_runtime_map(
    runtime_map_path: Path,
    outputs: dict[str, dict[str, Any]],
    lock: dict[str, Any],
) -> None:
    runtime_provenance = lock.get("provenance", {}).get("runtime_key_map", {})
    expected_sha = (
        runtime_provenance.get("sha256")
        if isinstance(runtime_provenance, dict)
        else None
    )
    if (
        not isinstance(expected_sha, str)
        or len(expected_sha) != 64
        or not runtime_map_path.is_file()
        or _sha(runtime_map_path) != expected_sha
    ):
        _fail(
            "ASSET-RUNTIME-SHA-018",
            "runtime key map không còn đúng SHA đã khóa lúc build",
        )
    payload = _read_json(runtime_map_path, "ASSET-PREFLIGHT-022")
    entries = payload.get("entries")
    if not isinstance(entries, list) or len(entries) != 8:
        _fail("ASSET-PREFLIGHT-022", "runtime key map phải có đúng tám binding")
    # Game1 có call-site lịch sử dùng cả mainImage/mainimage. Trên snapshot
    # Windows chúng là cùng file; vẫn cấm hai output khác nhau đụng nhau sau
    # case-fold để không biến việc nới case thành một lỗ collision.
    folded_outputs: dict[str, str] = {}
    for relative in outputs:
        folded = relative.casefold()
        if folded in folded_outputs and folded_outputs[folded] != relative:
            _fail(
                "ASSET-PREFLIGHT-022",
                f"hai output chỉ khác hoa/thường: {folded_outputs[folded]} / {relative}",
            )
        folded_outputs[folded] = relative
    for entry in entries:
        paths = entry.get("physical_paths") if isinstance(entry, dict) else None
        if not isinstance(paths, dict):
            _fail("ASSET-PREFLIGHT-022", "runtime binding thiếu physical_paths")
        for zoom in ("x1", "x2", "x3", "x4"):
            values = paths.get(zoom)
            if not isinstance(values, list) or not values:
                _fail("ASSET-PREFLIGHT-022", f"runtime binding thiếu {zoom}")
            for physical in values:
                prefix = "Assets/Resources/res/"
                if not isinstance(physical, str) or not physical.startswith(prefix):
                    _fail("ASSET-PREFLIGHT-022", f"physical path sai: {physical!r}")
                relative = "client/" + physical.removeprefix(prefix)
                if relative.casefold() not in folded_outputs:
                    _fail("ASSET-PREFLIGHT-022", f"runtime output chưa bị lock: {relative}")


def unity_preflight(
    unity_exe: Path,
    snapshot: Path,
    runtime_map: Path,
    contract: Path,
    log_path: Path,
) -> int:
    """Chạy đúng EditMode batch preflight của Task 5 trên snapshot đã deploy."""

    command = [
        str(unity_exe),
        "-batchmode",
        "-nographics",
        "-projectPath",
        str(snapshot),
        "-executeMethod",
        "NroE2E.NroResourcePreflightTests.RunBatch",
        "-runtime-key-map",
        str(runtime_map.resolve()),
        "-resource-contract",
        str(contract.resolve()),
        "-logFile",
        str(log_path.resolve()),
    ]
    try:
        completed = subprocess.run(command, timeout=600, check=False)
    except (OSError, subprocess.TimeoutExpired) as error:
        _fail("ASSET-UNITY-PREFLIGHT-023", str(error))
    return completed.returncode


def deploy(
    server_artifact: Path,
    client_artifact: Path,
    release_root: Path,
    snapshot: Path,
    version_ledger: Path,
    runtime_map: Path,
    *,
    previous_lock_path: Path,
    protected_unity_root: Path,
    preflight: Preflight,
) -> Path:
    """Kiểm trước, deploy cặp lock vào snapshot, rồi xuất release immutable."""

    server_artifact = Path(server_artifact).resolve()
    client_artifact = Path(client_artifact).resolve()
    release_root = Path(release_root).resolve()
    snapshot = Path(snapshot).resolve()
    runtime_map = Path(runtime_map).resolve()
    if previous_lock_path is None:
        _fail("ASSET-VERSION-015", "deploy phai co previous asset lock")
    if protected_unity_root is None:
        _fail("ASSET-DEPLOY-SOURCE-020", "deploy phai co Unity root can bao ve")
    previous_lock_path = Path(previous_lock_path).resolve()
    protected_unity_root = Path(protected_unity_root).resolve()
    server_lock, server_lock_raw = _read_artifact_lock(server_artifact)
    client_lock, client_lock_raw = _read_artifact_lock(client_artifact)
    if _sha_bytes(server_lock_raw) != _sha_bytes(client_lock_raw):
        _fail("ASSET-LOCK-MISMATCH-017", "server/client dùng hai asset-lock khác nhau")
    outputs = _verify_outputs(server_lock, server_artifact, client_artifact)
    _validate_version_bump(server_lock, Path(version_ledger), previous_lock_path)
    _validate_runtime_map(runtime_map, outputs, server_lock)
    resource_contract = _resource_contract(server_lock, outputs)

    if _inside(snapshot, protected_unity_root):
        _fail("ASSET-DEPLOY-SOURCE-020", "không được deploy vào Unity gốc")
    if not (snapshot / ".nro-e2e-overlay.json").is_file():
        _fail("ASSET-DEPLOY-SNAPSHOT-021", "đích không phải snapshot đã phủ overlay")
    active_marker = snapshot / ".nro-asset-lock.json"
    if active_marker.exists():
        _fail("ASSET-DEPLOY-ACTIVE-019", "snapshot đang giữ một artifact chạy")

    lock_sha = _sha_bytes(server_lock_raw)
    release = release_root / lock_sha
    if release.exists():
        _fail("ASSET-DEPLOY-ACTIVE-019", f"release immutable đã tồn tại: {release}")
    release_root.mkdir(parents=True, exist_ok=True)
    temporary = Path(tempfile.mkdtemp(prefix=f".{lock_sha}.", dir=release_root))
    touched: list[tuple[Path, bytes | None]] = []
    try:
        for relative in sorted(outputs):
            source = _artifact_file(relative, server_artifact, client_artifact)
            release_file = temporary / relative
            release_file.parent.mkdir(parents=True, exist_ok=True)
            shutil.copyfile(source, release_file)
            destination = _client_destination(snapshot, relative)
            if destination is None:
                continue
            old = destination.read_bytes() if destination.is_file() else None
            touched.append((destination, old))
            destination.parent.mkdir(parents=True, exist_ok=True)
            shutil.copyfile(source, destination)

        (temporary / "asset-lock.json").write_bytes(server_lock_raw)
        shutil.copyfile(runtime_map, temporary / "runtime-key-map.json")
        contract_path = temporary / "resource-contract.json"
        contract_path.write_text(
            json.dumps(
                resource_contract,
                ensure_ascii=False,
                sort_keys=True,
                separators=(",", ":"),
            )
            + "\n",
            encoding="utf-8",
        )
        log_path = temporary / "preflight.log"
        code = preflight(snapshot, runtime_map, contract_path, log_path)
        log = log_path.read_text(encoding="utf-8", errors="replace") if log_path.is_file() else ""
        if code != 0 or "RESOURCE-PREFLIGHT-OK" not in log:
            _fail(
                "ASSET-UNITY-PREFLIGHT-023",
                f"Unity exit={code}, thiếu RESOURCE-PREFLIGHT-OK",
            )

        marker = {
            "schema_version": 1,
            "lock_sha256": lock_sha,
            "server_client_pair_sha256": server_lock["server_client_pair_sha256"],
            "versions": server_lock["versions"],
        }
        active_marker.write_text(
            json.dumps(marker, ensure_ascii=False, sort_keys=True) + "\n",
            encoding="utf-8",
        )
        temporary.replace(release)
        return release
    except Exception:
        # Precondition phía trên chứng minh marker chưa tồn tại trước lượt này;
        # nếu publish release lỗi sau preflight thì phải rollback cả dấu active.
        active_marker.unlink(missing_ok=True)
        for destination, old in reversed(touched):
            if old is None:
                destination.unlink(missing_ok=True)
            else:
                destination.write_bytes(old)
        shutil.rmtree(temporary, ignore_errors=True)
        raise
