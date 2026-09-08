# -*- coding: utf-8 -*-
"""Lắp thử đúng ba mẫu UI Đợt 4, có cổng Unity và backup phục hồi được."""

from __future__ import annotations

import argparse
import hashlib
import json
import shutil
import subprocess
import sys
from datetime import datetime, timezone
from pathlib import Path

from PIL import Image


class PreflightError(RuntimeError):
    pass


DEFAULT_CONTRACT = {
    "unity_root": r"D:\NRO_Project",
    "sources": {
        "ui_hud_khung_trai_x4.png": {
            "size": [620, 164],
            "sha256": "1eb658cc3bffeb74e2105efe178f5950db865b3d8949dfcbf98837156cec13fd",
        },
        "ui_khung_popup_9o_x4.png": {
            "size": [40, 320],
            "sha256": "0fc7659d319bb756ba954a51dae9b570bef171b24bfe8f4cd457be2d9e6bc2e6",
        },
        "ui_nut_thuong_3manh_x4.png": {
            "size": [120, 96],
            "sha256": "15edc2e5ebf57e29e05b0276d1d5cacea14952c1abb74103f0275c7228638d6c",
        },
    },
}

RATIOS = {"x1": 0.25, "x2": 0.5, "x3": 0.75, "x4": 1.0}
RUNTIME_LEVELS = ("x2", "x3", "x4")
HUD_SOURCE = "ui_hud_khung_trai_x4.png"
POPUP_SOURCE = "ui_khung_popup_9o_x4.png"
BUTTON_SOURCE = "ui_nut_thuong_3manh_x4.png"


def sha256(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for block in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(block)
    return digest.hexdigest()


def image_size(path: Path) -> tuple[int, int]:
    try:
        with Image.open(path) as image:
            return image.size
    except Exception as error:
        raise PreflightError(f"UI-DOT4-PNG-INVALID {path}: {error}") from error


def scaled(size: list[int], level: str) -> tuple[int, int]:
    ratio = RATIOS[level]
    return round(size[0] * ratio), round(size[1] * ratio)


def zoom_name(source_name: str, level: str) -> str:
    return source_name.removesuffix("_x4.png") + f"_{level}.png"


def target_paths(unity_root: Path) -> list[Path]:
    names = (
        "myTexture2dpanel.png",
        "myTexture2dimgPopup.png",
        "btn0left.png",
        "btn0mid.png",
        "btn0right.png",
    )
    return [
        unity_root / "Assets" / "Resources" / "res" / level / "mainimage" / name
        for level in RUNTIME_LEVELS
        for name in names
    ]


def load_contract(path: Path | None) -> dict:
    if path is None:
        return DEFAULT_CONTRACT
    try:
        contract = json.loads(path.read_text(encoding="utf-8"))
    except Exception as error:
        raise PreflightError(f"UI-DOT4-CONTRACT-INVALID {path}: {error}") from error
    if not isinstance(contract.get("sources"), dict):
        raise PreflightError("UI-DOT4-CONTRACT-INVALID thieu sources")
    return contract


def verify_sources(source_root: Path, contract: dict) -> None:
    for name, expected in contract["sources"].items():
        path = source_root / name
        if not path.is_file():
            raise PreflightError(f"UI-DOT4-SOURCE-MISSING {path}")
        actual_size = image_size(path)
        if actual_size != tuple(expected["size"]):
            raise PreflightError(
                f"UI-DOT4-SOURCE-SIZE {name} actual={actual_size} expected={tuple(expected['size'])}"
            )
        actual_sha = sha256(path)
        if actual_sha != expected["sha256"]:
            raise PreflightError(
                f"UI-DOT4-SOURCE-SHA {name} actual={actual_sha} expected={expected['sha256']}"
            )


def load_unity_state(args: argparse.Namespace) -> dict:
    if args.unity_state_file:
        try:
            return json.loads(args.unity_state_file.read_text(encoding="utf-8"))
        except Exception as error:
            raise PreflightError(f"UI-DOT4-UNITY-STATE-INVALID: {error}") from error
    if not args.unity_bridge or not args.unity_state_code:
        raise PreflightError("UI-DOT4-UNITY-STATE-MISSING")
    result = subprocess.run(
        [sys.executable, str(args.unity_bridge), "--code", str(args.unity_state_code)],
        text=True,
        encoding="utf-8",
        errors="replace",
        capture_output=True,
        check=False,
        timeout=140,
    )
    if result.returncode != 0:
        raise PreflightError(
            f"UI-DOT4-UNITY-MCP-UNAVAILABLE exit={result.returncode} {result.stdout}{result.stderr}"
        )
    marker = "UI_DOT4_STATE="
    mcp_text = result.stdout.strip()
    try:
        envelope = json.loads(mcp_text)
        wrapped_result = envelope.get("data", {}).get("result")
        if isinstance(wrapped_result, str):
            mcp_text = wrapped_result
    except (json.JSONDecodeError, AttributeError):
        pass
    start = mcp_text.find(marker)
    if start < 0:
        raise PreflightError(f"UI-DOT4-UNITY-STATE-INVALID raw={result.stdout!r}")
    raw = mcp_text[start + len(marker) :].splitlines()[0].strip()
    try:
        return json.loads(raw)
    except Exception as error:
        raise PreflightError(f"UI-DOT4-UNITY-STATE-INVALID raw={raw!r}") from error


def verify_unity_state(unity_root: Path, state: dict, contract: dict) -> None:
    if state.get("is_playing") or state.get("is_playing_or_will_change"):
        raise PreflightError("UI-DOT4-UNITY-PLAYING")
    if state.get("is_compiling") or state.get("is_updating"):
        raise PreflightError("UI-DOT4-UNITY-IMPORT-LOCKED")
    expected_data = (unity_root / "Assets").resolve()
    try:
        actual_data = Path(state["data_path"]).resolve()
    except Exception as error:
        raise PreflightError("UI-DOT4-UNITY-STATE-INVALID thieu data_path") from error
    if actual_data != expected_data:
        raise PreflightError(
            f"UI-DOT4-UNITY-PROJECT actual={actual_data} expected={expected_data}"
        )
    pinned_root = contract.get("unity_root")
    if pinned_root and unity_root.resolve() != Path(pinned_root).resolve():
        raise PreflightError(
            f"UI-DOT4-UNITY-ROOT actual={unity_root.resolve()} expected={Path(pinned_root).resolve()}"
        )


def verify_zoom_outputs(zoom_root: Path, contract: dict) -> None:
    for name, expected in contract["sources"].items():
        for level in RATIOS:
            path = zoom_root / level / zoom_name(name, level)
            if not path.is_file():
                raise PreflightError(f"UI-DOT4-ZOOM-MISSING {path}")
            actual = image_size(path)
            wanted = scaled(expected["size"], level)
            if actual != wanted:
                raise PreflightError(
                    f"UI-DOT4-ZOOM-SIZE {path.name} actual={actual} expected={wanted}"
                )


def stage_outputs(zoom_root: Path, staging_root: Path, contract: dict) -> None:
    verify_zoom_outputs(zoom_root, contract)
    for level in RUNTIME_LEVELS:
        output = staging_root / "Assets" / "Resources" / "res" / level / "mainimage"
        output.mkdir(parents=True, exist_ok=True)

        shutil.copy2(
            zoom_root / level / zoom_name(HUD_SOURCE, level),
            output / "myTexture2dpanel.png",
        )

        popup_path = zoom_root / level / zoom_name(POPUP_SOURCE, level)
        with Image.open(popup_path).convert("RGBA") as popup:
            cell_height = popup.height // 8
            if cell_height * 8 != popup.height:
                raise PreflightError(f"UI-DOT4-POPUP-SLICE {popup_path}")
            rebuilt = Image.new("RGBA", popup.size)
            for index in range(8):
                top = index * cell_height
                rebuilt.paste(popup.crop((0, top, popup.width, top + cell_height)), (0, top))
            rebuilt.save(output / "myTexture2dimgPopup.png", "PNG", optimize=True)

        button_path = zoom_root / level / zoom_name(BUTTON_SOURCE, level)
        with Image.open(button_path).convert("RGBA") as button:
            piece_width = button.width // 3
            if piece_width * 3 != button.width:
                raise PreflightError(f"UI-DOT4-BUTTON-SLICE {button_path}")
            for index, name in enumerate(("btn0left.png", "btn0mid.png", "btn0right.png")):
                left = index * piece_width
                button.crop((left, 0, left + piece_width, button.height)).save(
                    output / name, "PNG", optimize=True
                )
    verify_staging(staging_root, contract)
    print("UI-DOT4-STAGED count=15")


def relative_target(path: Path, root: Path) -> str:
    return path.relative_to(root).as_posix()


def verify_staging(staging_root: Path, contract: dict) -> None:
    expected_sizes = {}
    for level in RUNTIME_LEVELS:
        expected_sizes[level, "myTexture2dpanel.png"] = scaled(
            contract["sources"][HUD_SOURCE]["size"], level
        )
        expected_sizes[level, "myTexture2dimgPopup.png"] = scaled(
            contract["sources"][POPUP_SOURCE]["size"], level
        )
        button_full = scaled(contract["sources"][BUTTON_SOURCE]["size"], level)
        for name in ("btn0left.png", "btn0mid.png", "btn0right.png"):
            expected_sizes[level, name] = (button_full[0] // 3, button_full[1])
    files = sorted(staging_root.rglob("*.png"))
    if len(files) != 15:
        raise PreflightError(f"UI-DOT4-STAGING-COUNT actual={len(files)} expected=15")
    for (level, name), wanted in expected_sizes.items():
        path = staging_root / "Assets" / "Resources" / "res" / level / "mainimage" / name
        if not path.is_file():
            raise PreflightError(f"UI-DOT4-STAGING-MISSING {path}")
        if image_size(path) != wanted:
            raise PreflightError(
                f"UI-DOT4-STAGING-SIZE {path} actual={image_size(path)} expected={wanted}"
            )


def create_backup(unity_root: Path, backup_root: Path) -> None:
    targets = target_paths(unity_root)
    missing = [path for path in targets if not path.is_file()]
    if missing:
        raise PreflightError("UI-DOT4-TARGET-MISSING " + ", ".join(map(str, missing)))
    manifest_path = backup_root / "backup-manifest.json"
    if manifest_path.exists():
        verify_backup(unity_root, backup_root)
        print("UI-DOT4-BACKUP-VERIFIED count=15 existing=True")
        return
    files = []
    for target in targets:
        relative = target.relative_to(unity_root)
        backup = backup_root / relative
        if backup.exists():
            raise PreflightError(f"UI-DOT4-BACKUP-UNTRACKED {backup}")
        backup.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(target, backup)
        files.append(
            {
                "relative_path": relative.as_posix(),
                "sha256": sha256(target),
                "size": list(image_size(target)),
            }
        )
    manifest = {
        "created_at_utc": datetime.now(timezone.utc).isoformat(),
        "unity_root": str(unity_root.resolve()),
        "files": files,
    }
    backup_root.mkdir(parents=True, exist_ok=True)
    manifest_path.write_text(json.dumps(manifest, ensure_ascii=False, indent=2), encoding="utf-8")
    verify_backup(unity_root, backup_root)
    print("UI-DOT4-BACKUP-CREATED count=15")


def verify_backup(unity_root: Path, backup_root: Path) -> dict:
    manifest_path = backup_root / "backup-manifest.json"
    if not manifest_path.is_file():
        raise PreflightError("UI-DOT4-BACKUP-MISSING manifest")
    try:
        manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    except Exception as error:
        raise PreflightError(f"UI-DOT4-BACKUP-INVALID {error}") from error
    expected = {relative_target(path, unity_root) for path in target_paths(unity_root)}
    records = {item.get("relative_path"): item for item in manifest.get("files", [])}
    if set(records) != expected:
        raise PreflightError(
            f"UI-DOT4-BACKUP-MISSING expected={len(expected)} actual={len(records)}"
        )
    for relative, item in records.items():
        backup = backup_root / Path(relative)
        if not backup.is_file() or sha256(backup) != item.get("sha256"):
            raise PreflightError(f"UI-DOT4-BACKUP-SHA {relative}")
        if list(image_size(backup)) != item.get("size"):
            raise PreflightError(f"UI-DOT4-BACKUP-SIZE {relative}")
    return manifest


def install(unity_root: Path, backup_root: Path, staging_root: Path, contract: dict) -> None:
    backup = verify_backup(unity_root, backup_root)
    verify_staging(staging_root, contract)
    backup_records = {item["relative_path"]: item for item in backup["files"]}
    installed = []
    for target in target_paths(unity_root):
        relative = target.relative_to(unity_root)
        source = staging_root / relative
        shutil.copy2(source, target)
        installed.append(
            {
                "relative_path": relative.as_posix(),
                "backup_sha256": backup_records[relative.as_posix()]["sha256"],
                "installed_sha256": sha256(target),
                "size": list(image_size(target)),
            }
        )
    (backup_root / "installation-manifest.json").write_text(
        json.dumps(
            {
                "installed_at_utc": datetime.now(timezone.utc).isoformat(),
                "unity_root": str(unity_root.resolve()),
                "files": installed,
            },
            ensure_ascii=False,
            indent=2,
        ),
        encoding="utf-8",
    )
    print("UI-DOT4-INSTALLED count=15")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("action", choices=("stage", "backup", "install", "verify"))
    parser.add_argument("--source-root", type=Path, required=True)
    parser.add_argument("--zoom-root", type=Path, required=True)
    parser.add_argument("--unity-root", type=Path, required=True)
    parser.add_argument("--backup-root", type=Path, required=True)
    parser.add_argument("--staging-root", type=Path, required=True)
    parser.add_argument("--contract", type=Path)
    parser.add_argument("--unity-state-file", type=Path)
    parser.add_argument("--unity-bridge", type=Path)
    parser.add_argument("--unity-state-code", type=Path)
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    try:
        contract = load_contract(args.contract)
        verify_sources(args.source_root, contract)
        if args.action == "stage":
            stage_outputs(args.zoom_root, args.staging_root, contract)
            return 0
        state = load_unity_state(args)
        verify_unity_state(args.unity_root, state, contract)
        if args.action == "backup":
            create_backup(args.unity_root, args.backup_root)
        elif args.action == "install":
            install(args.unity_root, args.backup_root, args.staging_root, contract)
        else:
            verify_zoom_outputs(args.zoom_root, contract)
            verify_staging(args.staging_root, contract)
            verify_backup(args.unity_root, args.backup_root)
            print("UI-DOT4-VERIFIED count=15")
        return 0
    except (PreflightError, subprocess.TimeoutExpired) as error:
        print(error)
        return 2


if __name__ == "__main__":
    raise SystemExit(main())
