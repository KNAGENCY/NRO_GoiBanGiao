"""Sinh production closure Lang Lá từ đúng nguồn đã review.

Module này cố ý không tự gán PASS cho nguồn mới. Mọi sprite canonical đều được
crop từ một asset đã có trong ``vertical-slice-v1.json``; các input còn lại được
copy/chốt SHA vào một closure bất biến dưới ArtRoot dùng chung.
"""

from __future__ import annotations

import argparse
from collections import deque
import hashlib
import json
import shutil
import tempfile
from pathlib import Path
from typing import Any

import numpy as np
from PIL import Image

from assetpack.animation import REQUIRED_ANIMATION_SHEETS, REQUIRED_CLIENT_CF
from assetpack.build import _source_paths
from tools.artqa.artqa import _adaptive_chroma_mask, extract_alpha
from tools.artqa.generate_tile_sheet import generate_tile_sheet


ZOOMS = ("x1", "x2", "x3", "x4")
ART_PREFIX = Path("docs/my-thuat/anh-nguon/vertical-slice-v1")
CLOSURE_RELATIVE = ART_PREFIX / "generated/production-closure-v9"
APPROVED_UI_RELATIVE = Path(
    "my thuat/dot1-nhom1-duyet/dot1-nhom1-duyet/01-giao-dien"
)
TILE_SOURCE_SHA256 = "2031edef163fc766b9f44df6c30724c0e41c57730b4f98749e3940fae157036b"

PLAYER_ACTIONS = (
    ("idle", "player_idle_3f_chroma.png", 3),
    ("run", "player_run_4f_chroma.png", 4),
    ("melee", "player_melee_4f_chroma.png", 4),
    ("cast", "player_cast_3f_chroma.png", 3),
    ("hit_death_getup", "player_hit_death_getup_3f_chroma.png", 3),
)
WOLF_ACTIONS = (
    ("idle", "wolf_idle_2f_chroma.png", 2),
    ("run", "wolf_run_4f_chroma.png", 4),
    ("pounce", "wolf_pounce_2f_chroma.png", 2),
    ("hit_death", "wolf_hit_death_2f_chroma.png", 2),
)
UI_ASSETS = (
    "ui_logo_game_x4.png",
    "ui_khung_avatar_x4.png",
    "ui_khung_hoithoai_x4.png",
    "ui_nut_thuong_x4.png",
    "ui_nut_nhan_x4.png",
    "ui_nut_khoa_x4.png",
    "ui_o_vatpham_x4.png",
    "ui_o_dangchon_x4.png",
)


class ClosureError(ValueError):
    """Input của production closure thiếu hoặc lệch provenance."""


def _sha(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def _review_by_name(review: dict[str, Any]) -> dict[str, dict[str, Any]]:
    assets = review.get("assets")
    if not isinstance(assets, list):
        raise ClosureError("ASSET-CLOSURE-015: review thiếu assets[]")
    result = {}
    for record in assets:
        if not isinstance(record, dict) or not isinstance(record.get("asset"), str):
            raise ClosureError("ASSET-CLOSURE-015: review asset malformed")
        if record.get("status") != "PASS":
            continue
        result[record["asset"]] = record
    return result


def _frame_rect(record: dict[str, Any], frame_index: int) -> list[int]:
    crop = record.get("checks", {}).get("crop", {})
    edges = crop.get("source_edges")
    boxes = crop.get("content_boxes")
    if (
        not isinstance(edges, list)
        or not isinstance(boxes, list)
        or frame_index < 0
        or frame_index >= len(boxes)
        or len(edges) != len(boxes) + 1
    ):
        raise ClosureError(
            f"ASSET-CLOSURE-015: crop review malformed: {record.get('asset')}"
        )
    left, top, right, bottom = boxes[frame_index]
    return [
        edges[frame_index] + left,
        top,
        right - left,
        bottom - top,
    ]


def _component_rect(
    record: dict[str, Any], frame_index: int, component: str
) -> list[int]:
    x, y, width, height = _frame_rect(record, frame_index)
    # Chia theo cơ thể, có overlap 10% để không tạo khe khi runtime ghép part.
    if component == "head":
        top = 0
        bottom = max(1, round(height * 0.42))
    elif component == "body":
        top = round(height * 0.28)
        bottom = max(top + 1, round(height * 0.76))
    elif component == "leg":
        top = round(height * 0.66)
        bottom = height
    else:
        raise ClosureError(f"ASSET-CLOSURE-015: component lạ: {component}")
    return [x, y + top, width, bottom - top]


def derive_player_slot_specs(review: dict[str, Any]) -> list[dict[str, Any]]:
    """Trả đúng 3/17/14 slot từ các action strip đã review."""

    records = _review_by_name(review)
    frames: list[dict[str, Any]] = []
    for action, source_name, expected_count in PLAYER_ACTIONS:
        record = records.get(source_name)
        if record is None:
            raise ClosureError(
                f"ASSET-CLOSURE-015: nguồn player chưa PASS: {source_name}"
            )
        boxes = record.get("checks", {}).get("crop", {}).get("content_boxes", [])
        if len(boxes) != expected_count:
            raise ClosureError(
                f"ASSET-CLOSURE-015: frame count lệch review: {source_name}"
            )
        for frame_index in range(expected_count):
            frames.append(
                {
                    "action": action,
                    "action_slot": frame_index,
                    "source_path": source_name,
                    "record": record,
                    "frame_index": frame_index,
                }
            )

    idle = frames[:3]
    without_hit = frames[:14]
    groups = (
        ("head", 27041, idle),
        ("body", 27044, frames),
        ("leg", 27061, without_hit),
    )
    slots = []
    for component, first_id, selected in groups:
        for slot_index, frame in enumerate(selected):
            record = frame["record"]
            slots.append(
                {
                    "id": first_id + slot_index,
                    "component": component,
                    "slot_index": slot_index,
                    "action": frame["action"],
                    "action_slot": frame["action_slot"],
                    "source_path": frame["source_path"],
                    "source_sha256": record["source_sha256"],
                    "source_rect": _component_rect(
                        record, frame["frame_index"], component
                    ),
                    "anchor": "bottom-center",
                    "layer": {"head": 2, "body": 1, "leg": 0}[component],
                    "overlap_ratio": {
                        "head_body": 0.14,
                        "body_leg": 0.10,
                    },
                    "client_cf": list(
                        REQUIRED_CLIENT_CF[f"hoc-vien-{component}"][slot_index]
                    ),
                }
            )
    return slots


def fit_to_canvas(image: Image.Image, size: tuple[int, int]) -> Image.Image:
    """Fit RGBA theo tỉ lệ, neo giữa-dưới và không cắt alpha."""

    source = image.convert("RGBA")
    bbox = source.getchannel("A").getbbox()
    if bbox is None:
        raise ClosureError("ASSET-CLOSURE-015: crop alpha rỗng")
    source = source.crop(bbox)
    width, height = size
    scale = min(width / source.width, height / source.height)
    target = (
        max(1, round(source.width * scale)),
        max(1, round(source.height * scale)),
    )
    resized = source.resize(target, Image.Resampling.NEAREST)
    canvas = Image.new("RGBA", size)
    canvas.alpha_composite(
        resized,
        ((width - resized.width) // 2, height - resized.height),
    )
    return canvas


def remove_large_chroma_islands(
    source: Image.Image, extracted: Image.Image
) -> Image.Image:
    """Xóa matte magenta nội vùng lớn, giữ accent màu nhỏ của nhân vật."""

    rgb = np.asarray(source.convert("RGB"), dtype=np.int16)
    high = np.maximum(rgb[:, :, 0], rgb[:, :, 2])
    low = np.minimum(rgb[:, :, 0], rgb[:, :, 2])
    candidate = (
        (low - rgb[:, :, 1] >= 70)
        & (high >= 150)
        & (high - rgb[:, :, 1] >= 100)
        & (np.abs(rgb[:, :, 0] - rgb[:, :, 2]) <= 95)
    )
    alpha = np.asarray(extracted.convert("RGBA").getchannel("A")).copy()
    remaining = candidate & (alpha > 0)
    visited = np.zeros(remaining.shape, dtype=bool)
    # Với master Iruka, matte gap nhỏ nhất có 633 px còn accent bông tai 586 px;
    # ngưỡng theo diện tích này tách được hai vùng mà không sửa tay nguồn.
    threshold = max(16, round(remaining.size * 0.00038))
    height, width = remaining.shape
    for start_y, start_x in zip(*np.where(remaining & ~visited)):
        if visited[start_y, start_x]:
            continue
        component = []
        queue = deque([(int(start_x), int(start_y))])
        visited[start_y, start_x] = True
        while queue:
            x, y = queue.popleft()
            component.append((x, y))
            for nx, ny in ((x - 1, y), (x + 1, y), (x, y - 1), (x, y + 1)):
                if (
                    0 <= nx < width
                    and 0 <= ny < height
                    and remaining[ny, nx]
                    and not visited[ny, nx]
                ):
                    visited[ny, nx] = True
                    queue.append((nx, ny))
        if len(component) >= threshold:
            for x, y in component:
                alpha[y, x] = 0
    result = extracted.convert("RGBA")
    result.putalpha(Image.fromarray(alpha, mode="L"))
    return result


def _save_zooms(
    x4: Image.Image, temporary: Path, relative_stem: Path
) -> tuple[dict[str, str], dict[str, str]]:
    paths = {}
    hashes = {}
    for zoom_index, zoom in enumerate(ZOOMS, start=1):
        size = (
            max(1, round(x4.width * zoom_index / 4)),
            max(1, round(x4.height * zoom_index / 4)),
        )
        image = (
            x4
            if zoom == "x4"
            else x4.resize(size, Image.Resampling.NEAREST)
        )
        destination = temporary / relative_stem / zoom / "sprite.png"
        destination.parent.mkdir(parents=True, exist_ok=True)
        image.save(destination, format="PNG", optimize=False, compress_level=9)
        repo_relative = CLOSURE_RELATIVE / relative_stem / zoom / "sprite.png"
        paths[zoom] = repo_relative.as_posix()
        hashes[zoom] = _sha(destination)
    return paths, hashes


def _load_reviewed_alpha(
    art_root: Path,
    records: dict[str, dict[str, Any]],
    source_name: str,
    cache: dict[str, Image.Image],
) -> Image.Image:
    if source_name not in records:
        raise ClosureError(
            f"ASSET-CLOSURE-015: source chưa có PASS: {source_name}"
        )
    source = art_root / source_name
    record = records[source_name]
    if not source.is_file() or _sha(source) != record.get("source_sha256"):
        raise ClosureError(
            f"ASSET-CLOSURE-015: SHA source lệch review: {source_name}"
        )
    if source_name not in cache:
        with Image.open(source) as image:
            image.load()
            cache[source_name] = remove_large_chroma_islands(
                image, extract_alpha(image)
            )
    return cache[source_name]


def _sprite_record(
    spec: dict[str, Any],
    *,
    art_root: Path,
    records: dict[str, dict[str, Any]],
    alpha_cache: dict[str, Image.Image],
    temporary: Path,
    target_size: tuple[int, int],
) -> dict[str, Any]:
    source_name = spec["source_path"]
    source = art_root / source_name
    alpha = _load_reviewed_alpha(art_root, records, source_name, alpha_cache)
    x, y, width, height = spec["source_rect"]
    if x < 0 or y < 0 or x + width > alpha.width or y + height > alpha.height:
        raise ClosureError(
            f"ASSET-CLOSURE-015: crop ngoài source: {spec['id']}"
        )
    x4 = fit_to_canvas(alpha.crop((x, y, x + width, y + height)), target_size)
    if x4.getchannel("A").getbbox() is None:
        raise ClosureError(f"ASSET-CLOSURE-015: sprite alpha rỗng: {spec['id']}")
    zooms, zoom_hashes = _save_zooms(
        x4, temporary, Path("sprites") / str(spec["id"])
    )
    record = {
        "id": spec["id"],
        "name": spec.get("name", f"{spec.get('component', 'asset')}-{spec['id']}"),
        "crop": {
            "asset_id": spec["id"],
            "source_rect": spec["source_rect"],
        },
        "zooms": zooms,
        "source_path": (ART_PREFIX / source_name).as_posix(),
        "source_sha256": _sha(source),
        "x4_sha256": zoom_hashes["x4"],
        "source_dimensions": list(alpha.size),
        "source_alpha": False,
        "target_x4": list(target_size),
        "zoom_sha256": zoom_hashes,
        "anchor": spec.get("anchor", "bottom-center"),
        "layer": spec.get("layer", 0),
        "overlap_ratio": spec.get("overlap_ratio", {}),
    }
    for key in (
        "component",
        "slot_index",
        "action",
        "action_slot",
        "client_cf",
        "provenance_role",
    ):
        if key in spec:
            record[key] = spec[key]
    return record


def _single_frame_spec(
    records: dict[str, dict[str, Any]],
    *,
    sprite_id: int,
    source_name: str,
    frame_index: int,
    name: str,
    provenance_role: str,
) -> dict[str, Any]:
    record = records.get(source_name)
    if record is None:
        raise ClosureError(f"ASSET-CLOSURE-015: source chưa PASS: {source_name}")
    return {
        "id": sprite_id,
        "name": name,
        "source_path": source_name,
        "source_sha256": record["source_sha256"],
        "source_rect": _frame_rect(record, frame_index),
        "provenance_role": provenance_role,
    }


def _generate_sprites(
    review: dict[str, Any], art_root: Path, temporary: Path
) -> list[dict[str, Any]]:
    records = _review_by_name(review)
    cache: dict[str, Image.Image] = {}
    sprites = []
    targets = {"head": (128, 128), "body": (128, 160), "leg": (128, 128)}
    for spec in derive_player_slot_specs(review):
        sprites.append(
            _sprite_record(
                spec,
                art_root=art_root,
                records=records,
                alpha_cache=cache,
                temporary=temporary,
                target_size=targets[spec["component"]],
            )
        )

    iruka_record = records["master_iruka_x4_chroma.png"]
    for offset, component in enumerate(("head", "body", "leg")):
        spec = {
            "id": 27075 + offset,
            "name": f"iruka-{component}",
            "component": component,
            "slot_index": 0,
            "action": "idle",
            "action_slot": 0,
            "source_path": "master_iruka_x4_chroma.png",
            "source_sha256": iruka_record["source_sha256"],
            "source_rect": _component_rect(iruka_record, 0, component),
            "anchor": "bottom-center",
            "layer": {"head": 2, "body": 1, "leg": 0}[component],
            "overlap_ratio": {"head_body": 0.14, "body_leg": 0.10},
            "provenance_role": "approved_iruka_master_component",
        }
        sprites.append(
            _sprite_record(
                spec,
                art_root=art_root,
                records=records,
                alpha_cache=cache,
                temporary=temporary,
                target_size=targets[component],
            )
        )

    portrait = _single_frame_spec(
        records,
        sprite_id=27078,
        source_name="iruka_portrait_chroma.png",
        frame_index=0,
        name="iruka-portrait",
        provenance_role="approved_iruka_portrait",
    )
    sprites.append(
        _sprite_record(
            portrait,
            art_root=art_root,
            records=records,
            alpha_cache=cache,
            temporary=temporary,
            target_size=(200, 200),
        )
    )

    skill_record = records["skill_icons_leaf_4_chroma.png"]
    for slot in range(4):
        spec = _single_frame_spec(
            records,
            sprite_id=27079 + slot,
            source_name="skill_icons_leaf_4_chroma.png",
            frame_index=slot,
            name=f"skill-leaf-{slot}",
            provenance_role="approved_skill_icon",
        )
        spec.update({"slot_index": slot, "action": "skill_pose"})
        sprites.append(
            _sprite_record(
                spec,
                art_root=art_root,
                records=records,
                alpha_cache=cache,
                temporary=temporary,
                target_size=(64, 64),
            )
        )
    if [sprite["id"] for sprite in sprites] != list(range(27041, 27083)):
        raise ClosureError("ASSET-CLOSURE-015: sprite ID closure không liên tục")
    return sprites


def _generate_wolf_cells(
    review: dict[str, Any], art_root: Path, temporary: Path
) -> list[str]:
    records = _review_by_name(review)
    cache: dict[str, Image.Image] = {}
    paths = []
    slot = 0
    for action, source_name, count in WOLF_ACTIONS:
        record = records.get(source_name)
        if record is None or len(
            record.get("checks", {}).get("crop", {}).get("content_boxes", [])
        ) != count:
            raise ClosureError(
                f"ASSET-CLOSURE-015: wolf frame closure lệch: {source_name}"
            )
        alpha = _load_reviewed_alpha(
            art_root, records, source_name, cache
        )
        for frame_index in range(count):
            x, y, width, height = _frame_rect(record, frame_index)
            cell = fit_to_canvas(
                alpha.crop((x, y, x + width, y + height)), (96, 96)
            )
            destination = (
                temporary / "animation" / "quai-soi-rung" / f"{slot:02d}.png"
            )
            destination.parent.mkdir(parents=True, exist_ok=True)
            cell.save(destination, format="PNG", compress_level=9)
            paths.append(
                (
                    CLOSURE_RELATIVE
                    / "animation"
                    / "quai-soi-rung"
                    / f"{slot:02d}.png"
                ).as_posix()
            )
            slot += 1
    return paths


def _animation_closure(
    sprites: list[dict[str, Any]], wolf_paths: list[str]
) -> tuple[list[dict[str, Any]], list[dict[str, Any]]]:
    by_id = {sprite["id"]: sprite for sprite in sprites}
    ranges = {
        "hoc-vien-head": list(range(27041, 27044)),
        "hoc-vien-body": list(range(27044, 27061)),
        "hoc-vien-leg": list(range(27061, 27075)),
    }
    sheets = []
    animations = []
    for name, required in REQUIRED_ANIMATION_SHEETS.items():
        if name == "quai-soi-rung":
            cell_paths = wolf_paths
        else:
            cell_paths = [by_id[sprite_id]["zooms"]["x4"] for sprite_id in ranges[name]]
        sheets.append(
            {
                "name": name,
                "slot_count": required["slot_count"],
                "cell_size": list(required["cell_size"]),
                "cells": [
                    {"slot_index": index, "path": path}
                    for index, path in enumerate(cell_paths)
                ],
            }
        )
        animations.append(
            {
                "name": name,
                "slot_count": required["slot_count"],
                "frames": [
                    {
                        "slot_index": index,
                        "client_cf": list(REQUIRED_CLIENT_CF[name][index]),
                    }
                    for index in range(required["slot_count"])
                ],
            }
        )
    return sheets, animations


def _generate_mob_atlases(
    shared_root: Path, temporary: Path, wolf_paths: list[str]
) -> dict[str, str]:
    # EffectData dùng toạ độ logical x1; atlas zoom lớn giữ nguyên layout nhân zoom.
    x4_cells = []
    for relative in wolf_paths:
        with Image.open(shared_root / relative) as image:
            image.load()
            x4_cells.append(image.convert("RGBA"))
    paths = {}
    for zoom_index, zoom in enumerate(ZOOMS, start=1):
        cell_size = 24 * zoom_index
        atlas = Image.new("RGBA", (cell_size * 10, cell_size))
        for slot, x4 in enumerate(x4_cells):
            cell = (
                x4
                if zoom == "x4"
                else x4.resize((cell_size, cell_size), Image.Resampling.NEAREST)
            )
            atlas.alpha_composite(cell, (slot * cell_size, 0))
        destination = temporary / "mob" / zoom / "img.png"
        destination.parent.mkdir(parents=True, exist_ok=True)
        atlas.save(destination, format="PNG", compress_level=9)
        paths[zoom] = (CLOSURE_RELATIVE / "mob" / zoom / "img.png").as_posix()
    return paths


def _largest_interior_alpha_component(image: Image.Image) -> Image.Image:
    """Crop đúng thanh HUD, bỏ đường viền artboard chạm mép quadrant."""

    source = image.convert("RGBA")
    alpha = np.asarray(source.getchannel("A"))
    visible = alpha > 0
    height, width = visible.shape
    visited = np.zeros_like(visible, dtype=bool)
    components: list[list[tuple[int, int]]] = []

    for start_y, start_x in zip(*np.where(visible & ~visited)):
        if visited[start_y, start_x]:
            continue
        component: list[tuple[int, int]] = []
        queue: deque[tuple[int, int]] = deque([(int(start_x), int(start_y))])
        visited[start_y, start_x] = True
        touches_edge = False
        while queue:
            x, y = queue.popleft()
            component.append((x, y))
            touches_edge = (
                touches_edge
                or x in (0, width - 1)
                or y in (0, height - 1)
            )
            for next_x, next_y in (
                (x - 1, y),
                (x + 1, y),
                (x, y - 1),
                (x, y + 1),
            ):
                if not (0 <= next_x < width and 0 <= next_y < height):
                    continue
                if not visible[next_y, next_x] or visited[next_y, next_x]:
                    continue
                visited[next_y, next_x] = True
                queue.append((next_x, next_y))
        if not touches_edge:
            components.append(component)

    if not components:
        raise ClosureError(
            "ASSET-CLOSURE-015: HUD không có component nội vùng"
        )

    component = max(components, key=len)
    xs = [point[0] for point in component]
    ys = [point[1] for point in component]
    bbox = (min(xs), min(ys), max(xs) + 1, max(ys) + 1)
    selected = np.zeros_like(alpha, dtype=np.uint8)
    for x, y in component:
        selected[y, x] = alpha[y, x]
    result = source.copy()
    result.putalpha(Image.fromarray(selected, mode="L"))
    return result.crop(bbox)


def _extract_hud_bar(quadrant: Image.Image) -> Image.Image:
    extracted = extract_alpha(quadrant)
    alpha = np.asarray(extracted.getchannel("A")).copy()
    alpha[_adaptive_chroma_mask(quadrant.convert("RGB"))] = 0
    extracted.putalpha(Image.fromarray(alpha, mode="L"))
    return _largest_interior_alpha_component(extracted)


def _generate_hud(
    art_root: Path,
    review: dict[str, Any],
    temporary: Path,
) -> dict[str, dict[str, str]]:
    records = _review_by_name(review)
    alpha = _load_reviewed_alpha(
        art_root, records, "hud_hp_chakra_4_chroma.png", {}
    )
    middle_x, middle_y = alpha.width // 2, alpha.height // 2
    boxes = {
        "naruto-hp-bg": (0, 0, middle_x, middle_y),
        "naruto-hp-fill": (middle_x, 0, alpha.width, middle_y),
        "naruto-chakra-bg": (0, middle_y, middle_x, alpha.height),
        "naruto-chakra-fill": (middle_x, middle_y, alpha.width, alpha.height),
    }
    result = {}
    for name, box in boxes.items():
        quadrant = alpha.crop(box)
        # Bỏ viền trắng phân ô của master trước khi trim alpha nội dung thanh.
        quadrant = quadrant.crop(
            (8, 8, quadrant.width - 8, quadrant.height - 8)
        )
        pixels = np.asarray(quadrant.convert("RGBA")).copy()
        red = pixels[:, :, 0].astype(np.int16)
        green = pixels[:, :, 1].astype(np.int16)
        blue = pixels[:, :, 2].astype(np.int16)
        # HUD dùng gradient magenta rộng hơn chroma chuẩn. Chỉ ở bốn quadrant
        # này, xóa pixel có đồng thời R/B cao và G thấp; red/blue fill vẫn giữ.
        matte = (
            (red >= 120)
            & (blue >= 120)
            & (green + 45 < np.minimum(red, blue))
        )
        pixels[:, :, 3][matte] = 0
        quadrant = Image.fromarray(pixels, mode="RGBA")
        content = _extract_hud_bar(quadrant)
        x4 = content.resize((240, 32), Image.Resampling.NEAREST)
        zooms, _ = _save_zooms(x4, temporary, Path("hud") / name)
        result[name] = zooms
    return result


def _generate_ui(
    shared_root: Path, temporary: Path
) -> dict[str, dict[str, str]]:
    approved = shared_root / APPROVED_UI_RELATIVE
    result = {}
    for name in UI_ASSETS:
        source = approved / name
        if not source.is_file():
            raise ClosureError(f"ASSET-CLOSURE-015: thiếu UI đã duyệt: {source}")
        with Image.open(source) as image:
            image.load()
            x4 = image.convert("RGBA")
        zooms, _ = _save_zooms(
            x4, temporary, Path("ui") / Path(name).stem
        )
        result[name] = zooms
    return result


def _copy_file(source: Path, temporary: Path, relative: Path) -> str:
    if not source.is_file():
        raise ClosureError(f"ASSET-CLOSURE-015: thiếu input: {source}")
    destination = temporary / relative
    destination.parent.mkdir(parents=True, exist_ok=True)
    shutil.copyfile(source, destination)
    return (CLOSURE_RELATIVE / relative).as_posix()


def _copy_runtime_inputs(
    project_root: Path, shared_root: Path, temporary: Path
) -> tuple[str, str, str, dict[str, str]]:
    part = _copy_file(
        project_root / "server/internal/game/part.bin",
        temporary,
        Path("inputs/part.bin"),
    )
    tileset = _copy_file(
        project_root / "server/internal/game/tileset.bin",
        temporary,
        Path("inputs/tileset.bin"),
    )
    runtime_map = _copy_file(
        project_root / "docs/my-thuat/unity-runtime-key-map.json",
        temporary,
        Path("inputs/unity-runtime-key-map.json"),
    )
    snapshots = {}
    payload = json.loads(
        (project_root / "docs/my-thuat/unity-runtime-key-map.json").read_text(
            encoding="utf-8"
        )
    )
    sites = {
        site["file"]: site["file_sha256"]
        for entry in payload["entries"]
        for site in [*entry["call_sites"], *entry["consumer_sites"]]
    }
    for site, expected_sha in sorted(sites.items()):
        source = Path("D:/NRO_Project") / site
        if not source.is_file() or _sha(source) != expected_sha:
            raise ClosureError(
                f"ASSET-CLOSURE-015: runtime site lệch precondition: {site}"
            )
        relative = Path("runtime-snapshots") / Path(site)
        snapshots[site] = _copy_file(source, temporary, relative)
    return part, tileset, runtime_map, snapshots


def _background_sources(shared_root: Path) -> dict[str, dict[str, str]]:
    bases = {
        "b00": "map_nen_troi_1",
        "b01": "map_nen_may_1",
        "b02": "map_nen_may_2",
        "b03": "map_nen_nuihokage",
    }
    result = {}
    for layer, base in bases.items():
        result[layer] = {}
        for zoom in ZOOMS:
            relative = Path("docs/my-thuat/anh-game") / zoom / f"{base}_{zoom}.png"
            if not (shared_root / relative).is_file():
                raise ClosureError(
                    f"ASSET-CLOSURE-015: thiếu background: {relative}"
                )
            result[layer][zoom] = relative.as_posix()
    return result


def _tile_art(
    art_root: Path, temporary: Path
) -> tuple[list[dict[str, Any]], list[dict[str, Any]]]:
    source = art_root / "generated/naruto-village-tile-sheet-10x10-v2.png"
    index = generate_tile_sheet(
        source,
        temporary / "tile-art",
        expected_sha256=TILE_SOURCE_SHA256,
    )
    rows = []
    hashes = []
    for tile in index["tiles"]:
        sources = {
            zoom: (
                CLOSURE_RELATIVE / "tile-art" / tile["zooms"][zoom]["path"]
            ).as_posix()
            for zoom in ZOOMS
        }
        rows.append({"index": tile["index"], "solid": False, "sources": sources})
        hashes.append(
            {
                "index": tile["index"],
                **{
                    zoom: tile["zooms"][zoom]["sha256"]
                    for zoom in ZOOMS
                },
            }
        )
    return rows, hashes


def _parts() -> list[dict[str, Any]]:
    def frames(ids: list[int], dy: int) -> list[dict[str, int]]:
        return [{"sprite_id": sprite_id, "dx": 0, "dy": dy} for sprite_id in ids]

    return [
        {"id": 30, "type": 0, "frames": frames(list(range(27041, 27044)), 0)},
        {"id": 57, "type": 1, "frames": frames(list(range(27044, 27061)), 0)},
        {"id": 58, "type": 2, "frames": frames(list(range(27061, 27075)), 0)},
        {"id": 2189, "type": 0, "frames": frames([27075] * 3, 0)},
        {"id": 2190, "type": 1, "frames": frames([27076] * 17, 0)},
        {"id": 2191, "type": 2, "frames": frames([27077] * 14, 0)},
    ]


def _skills() -> list[dict[str, Any]]:
    return [
        {
            "id": 27079 + slot,
            "pose_slot": slot,
            "levels": list(range(12)),
            "effect_id": -1,
            "arrow_id": -1,
        }
        for slot in range(4)
    ]


def _catalog() -> dict[str, Any]:
    return {
        "maps": [
            "Làng Lá",
            "Làng Cát",
            "Làng Sương Mù",
            "Làng Đá",
            "Làng Mây",
            "Khu rừng",
        ],
        "npc_templates": ["Iruka"],
        "mob_templates": ["Sói Rừng"],
        "spawn": [[100, 384]],
        "skill_poses": [0, 1, 2, 3],
    }


def _ui_bindings(
    project_root: Path,
    ui_sources: dict[str, dict[str, str]],
) -> list[dict[str, Any]]:
    runtime = json.loads(
        (project_root / "docs/my-thuat/unity-runtime-key-map.json").read_text(
            encoding="utf-8"
        )
    )
    bindings = []
    for entry in runtime["entries"]:
        name = entry["source_asset"]
        binding = dict(entry)
        binding["source_path"] = (
            APPROVED_UI_RELATIVE / name
        ).as_posix()
        binding["source_sha256"] = _sha(
            project_root.parents[1] / APPROVED_UI_RELATIVE / name
            if project_root.parent.name == ".worktrees"
            else project_root / APPROVED_UI_RELATIVE / name
        )
        binding["sources"] = ui_sources[name]
        bindings.append(binding)
    return bindings


def generate_production_closure(
    *,
    project_root: Path,
    shared_root: Path,
    review_path: Path,
    mapping_path: Path,
) -> dict[str, Any]:
    """Sinh closure atomically và ghi mapping canonical trong worktree."""

    project_root = project_root.resolve()
    shared_root = shared_root.resolve()
    art_root = shared_root / ART_PREFIX
    output = shared_root / CLOSURE_RELATIVE
    if output.exists():
        raise ClosureError(
            f"ASSET-CLOSURE-015: closure đã tồn tại, không ghi đè: {output}"
        )
    review = json.loads(review_path.read_text(encoding="utf-8"))
    output.parent.mkdir(parents=True, exist_ok=True)
    temporary = Path(
        tempfile.mkdtemp(prefix=".production-closure-v9-", dir=output.parent)
    )
    try:
        sprites = _generate_sprites(review, art_root, temporary)
        wolf_paths = _generate_wolf_cells(review, art_root, temporary)
        # wolf_paths trỏ tới tên canonical cuối; lúc compose dùng file trong temp.
        local_wolf_paths = [
            (temporary / Path(path).relative_to(CLOSURE_RELATIVE))
            for path in wolf_paths
        ]
        mob_paths = {}
        for zoom_index, zoom in enumerate(ZOOMS, start=1):
            cell_size = 24 * zoom_index
            atlas = Image.new("RGBA", (cell_size * 10, cell_size))
            for slot, path in enumerate(local_wolf_paths):
                with Image.open(path) as image:
                    image.load()
                    x4 = image.convert("RGBA")
                cell = (
                    x4
                    if zoom == "x4"
                    else x4.resize(
                        (cell_size, cell_size), Image.Resampling.NEAREST
                    )
                )
                atlas.alpha_composite(cell, (slot * cell_size, 0))
            destination = temporary / "mob" / zoom / "img.png"
            destination.parent.mkdir(parents=True, exist_ok=True)
            atlas.save(destination, format="PNG", compress_level=9)
            mob_paths[zoom] = (
                CLOSURE_RELATIVE / "mob" / zoom / "img.png"
            ).as_posix()

        hud_sources = _generate_hud(art_root, review, temporary)
        ui_sources = _generate_ui(shared_root, temporary)
        tile_art, tile_hashes = _tile_art(art_root, temporary)
        part_path, tileset_path, runtime_map, snapshots = _copy_runtime_inputs(
            project_root, shared_root, temporary
        )
        sheets, animations = _animation_closure(sprites, wolf_paths)
        ui_bindings = _ui_bindings(project_root, ui_sources)
        tileset_local = temporary / Path(tileset_path).relative_to(CLOSURE_RELATIVE)
        tileset_bytes = tileset_local.read_bytes()
        mob_images = [
            {"id": index, "x": index * 24, "y": 0, "w": 24, "h": 24}
            for index in range(10)
        ]
        mob_frames_flat = [[[-12, -24, index]] for index in range(10)]
        mob_frames = {
            "idle": mob_frames_flat[:2],
            "run": mob_frames_flat[2:6],
            "pounce": mob_frames_flat[6:8],
            "hit": mob_frames_flat[8:9],
            "death": mob_frames_flat[9:],
        }
        manifest = {
            "schema_version": 1,
            "source_catalog_version": 2,
            "versions": {
                "runtime_slice": 1,
                "data": 1,
                "map": 1,
                "resource": 1,
            },
            "phase": {
                "name": "lang-la",
                "sprite_id_range": [27041, 27082],
                "reserved_sprite_id_range": [27041, 27999],
                "forbidden_poc_ids": [9000, 9001, 9002],
            },
            "sprites": sprites,
            "parts": _parts(),
            "animations": animations,
            "maps": [{"id": 1, "tile_id": 1, "background_id": 0}],
            "tiles": [
                {
                    "id": index,
                    "semantic": "walkable" if index == 0 else "decorative",
                    "solid": False,
                }
                for index in range(100)
            ],
            "backgrounds": [{"id": 0, "layers": 4}],
            "mobs": [{"id": 0, "name": "soi-rung"}],
            "skills": _skills(),
            "catalog": _catalog(),
            "ui_bindings": ui_bindings,
            "build_inputs": {
                "source_base": ".",
                "source_provenance": [
                    {
                        "source_path": (ART_PREFIX / record["asset"]).as_posix(),
                        "source_sha256": record["source_sha256"],
                        "role": "reviewed_vertical_slice_source",
                    }
                    for record in sorted(
                        review["assets"], key=lambda item: item["asset"]
                    )
                ],
                "legacy_part_path": part_path,
                "legacy_part_sha256": _sha(
                    temporary / Path(part_path).relative_to(CLOSURE_RELATIVE)
                ),
                "tile_set_info_path": tileset_path,
                "tile_set_info_sha256": _sha(tileset_local),
                "tile_set_baseline": {
                    "sha256": _sha(tileset_local),
                    "bytes_hex": tileset_bytes.hex(),
                    "semantic": {
                        "byte_length": len(tileset_bytes),
                        "tile_count": tileset_bytes[0],
                    },
                },
                "tile_art": tile_art,
                "tile_art_sha256": tile_hashes,
                "background_sources": _background_sources(shared_root),
                "hud_sources": hud_sources,
                "mob": {
                    "id": 0,
                    "payload_mode": "generated",
                    "images": mob_images,
                    "frames": mob_frames,
                    "arr_frames": [],
                    "trailer": [0, 0],
                    "img": mob_paths,
                },
                "map": {
                    "width": 5,
                    "height": 17,
                    "indexes": [0] * 85,
                    "spawn": [100, 384],
                    "tile_size": [24, 24],
                    "client_solid": [False] * 85,
                },
                "animation_sheets": sheets,
                "runtime_key_map": runtime_map,
                "runtime_key_map_sha256": _sha(
                    temporary / Path(runtime_map).relative_to(CLOSURE_RELATIVE)
                ),
                "runtime_source_snapshots": snapshots,
            },
        }
        # Tất cả path canonical lúc này trỏ tới output cuối; hash file trong temp.
        closure_hashes = {}
        for relative in _source_paths(manifest):
            path = shared_root / relative
            if relative.startswith(CLOSURE_RELATIVE.as_posix() + "/"):
                path = temporary / Path(relative).relative_to(CLOSURE_RELATIVE)
            if not path.is_file():
                raise ClosureError(
                    f"ASSET-CLOSURE-015: input closure thiếu: {relative}"
                )
            closure_hashes[relative] = _sha(path)
        manifest["build_inputs"]["production_closure_sha256"] = closure_hashes

        index = {
            "schema_version": 1,
            "generator": "tools.artqa.task6b_closure",
            "review_sha256": _sha(review_path),
            "mapping_sha256": hashlib.sha256(
                json.dumps(
                    manifest,
                    ensure_ascii=False,
                    sort_keys=True,
                    separators=(",", ":"),
                ).encode("utf-8")
            ).hexdigest(),
            "source_count": len(closure_hashes),
            "sprite_ids": [sprite["id"] for sprite in sprites],
        }
        (temporary / "generated-index.json").write_text(
            json.dumps(index, ensure_ascii=False, indent=2) + "\n",
            encoding="utf-8",
        )
        temporary.replace(output)
        mapping_path.write_text(
            json.dumps(manifest, ensure_ascii=False, indent=2) + "\n",
            encoding="utf-8",
        )
        return manifest
    except Exception:
        shutil.rmtree(temporary, ignore_errors=True)
        raise


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--project-root", required=True, type=Path)
    parser.add_argument("--shared-root", required=True, type=Path)
    parser.add_argument("--review", required=True, type=Path)
    parser.add_argument("--mapping", required=True, type=Path)
    args = parser.parse_args()
    generate_production_closure(
        project_root=args.project_root,
        shared_root=args.shared_root,
        review_path=args.review,
        mapping_path=args.mapping,
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
