"""QA xác định được cho ảnh nguồn vertical slice.

Không sửa PNG nguồn. Chroma chỉ được tách từ vùng nối tới mép ảnh để một chi
tiết màu tím bên trong không bị khoét nhầm; nền chuyển sắc bị báo FAIL vì soft
matte không thể chứng minh được biên alpha lặp lại an toàn.
"""

from __future__ import annotations

import argparse
from collections import deque
from datetime import datetime, timezone
import hashlib
import json
from pathlib import Path
import re
from typing import Iterable

from PIL import Image, ImageDraw
import numpy as np
from scipy import ndimage


CHROMA_EDGE_TOLERANCE = 18
MIN_BORDER_CHROMA_RATIO = 0.85
MAX_BORDER_CHROMA_SPREAD = 12
LIGHT = (245, 242, 235, 255)
DARK = (28, 32, 43, 255)


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def is_chroma(pixel: tuple[int, int, int]) -> bool:
    """Nhận diện nền tím/magenta, không dựa vào tên file `chroma`."""
    red, green, blue = pixel
    return red - green >= 110 and blue - green >= 110 and red >= 150 and blue >= 150


def _border_points(width: int, height: int) -> Iterable[tuple[int, int]]:
    for x in range(width):
        yield x, 0
        yield x, height - 1
    for y in range(1, height - 1):
        yield 0, y
        yield width - 1, y


def extract_alpha(image: Image.Image) -> Image.Image:
    """Flood nền chroma từ mép và tạo soft-matte một pixel ở biên vật thể."""
    rgb = image.convert("RGB")
    pixels = np.asarray(rgb, dtype=np.int16)
    chroma = ((pixels[:, :, 0] - pixels[:, :, 1] >= 110)
              & (pixels[:, :, 2] - pixels[:, :, 1] >= 110)
              & (pixels[:, :, 0] >= 150) & (pixels[:, :, 2] >= 150))
    # Label hoá theo 4-neighbour rồi chỉ giữ thành phần có mặt ở biên: đây là
    # flood-fill được vector hoá, tránh treo hàng phút trên ảnh 1,5 triệu pixel.
    labels, _ = ndimage.label(chroma, structure=np.array(((0, 1, 0), (1, 1, 1), (0, 1, 0))))
    edge_labels = np.unique(np.concatenate((labels[0, :], labels[-1, :], labels[:, 0], labels[:, -1])))
    flooded = (labels > 0) & np.isin(labels, edge_labels[edge_labels > 0])
    matte = np.full(flooded.shape, 255, dtype=np.uint8)
    matte[flooded] = 0
    # Feather chỉ ở rìa. Một pixel bị nền vây kín là chi tiết thật, giữ opaque.
    neighbours = ndimage.convolve(flooded.astype(np.uint8), np.array(((0, 1, 0), (1, 0, 1), (0, 1, 0))), mode="constant")
    soft_edge = (~flooded) & (neighbours > 0) & (neighbours < 3)
    matte[soft_edge] = 192
    result = rgb.convert("RGBA")
    result.putalpha(Image.fromarray(matte, mode="L"))
    return result


def check_cell_boundaries(image: Image.Image, columns: int) -> list[str]:
    if columns < 1:
        return ["invalid_column_count"]
    if image.width % columns:
        return ["canvas_width_not_divisible_by_columns"]
    return []


def _expected_columns(name: str) -> int:
    match = re.search(r"_(\d+)(?:f|_chroma)" , name)
    return int(match.group(1)) if match else 1


def _border_metrics(image: Image.Image) -> tuple[float, int]:
    rgb = image.convert("RGB")
    points = list(_border_points(*rgb.size))
    chroma = [rgb.getpixel(point) for point in points if is_chroma(rgb.getpixel(point))]
    ratio = len(chroma) / len(points)
    if not chroma:
        return ratio, 0
    channels = list(zip(*chroma))
    return ratio, max(max(channel) - min(channel) for channel in channels)


def _anchors_and_loop(alpha: Image.Image, columns: int) -> list[str]:
    if columns <= 1 or alpha.width % columns:
        return []
    cell_width = alpha.width // columns
    boxes = []
    for index in range(columns):
        crop = alpha.crop((index * cell_width, 0, (index + 1) * cell_width, alpha.height))
        boxes.append(crop.getchannel("A").getbbox())
    if any(box is None for box in boxes):
        return ["empty_animation_cell"]
    bottoms = [box[3] for box in boxes if box]
    if max(bottoms) - min(bottoms) > max(4, alpha.height // 20):
        return ["anchor_bottom_drift"]
    # Chỉ strip có từ 3 frame mới kiểm tra loop. So sánh bbox đầu/cuối là điều kiện
    # cần, không tự nhận là animation mượt.
    if columns >= 3:
        first, last = boxes[0], boxes[-1]
        assert first is not None and last is not None
        if abs(first[1] - last[1]) > alpha.height // 8 or abs(first[3] - last[3]) > alpha.height // 8:
            return ["loop_boundary_anchor_mismatch"]
    return []


def _render_on_background(alpha: Image.Image, background: tuple[int, int, int, int]) -> Image.Image:
    rendered = Image.new("RGBA", alpha.size, background)
    rendered.alpha_composite(alpha)
    return rendered


def _contact_sheet(entries: list[tuple[str, Image.Image]], background: tuple[int, int, int, int]) -> Image.Image:
    thumb_size = (240, 160)
    columns = 4
    rows = (len(entries) + columns - 1) // columns
    sheet = Image.new("RGBA", (columns * 250, rows * 190), background)
    draw = ImageDraw.Draw(sheet)
    for index, (name, alpha) in enumerate(entries):
        thumb = _render_on_background(alpha, background)
        thumb.thumbnail(thumb_size, Image.Resampling.LANCZOS)
        x, y = (index % columns) * 250, (index // columns) * 190
        sheet.alpha_composite(thumb, (x + (240 - thumb.width) // 2, y + 4))
        draw.text((x + 5, y + 168), name[:34], fill=(20, 20, 20, 255) if background == LIGHT else (235, 235, 235, 255))
    return sheet


def _preview(alpha: Image.Image, columns: int) -> Image.Image:
    """Preview tĩnh tất cả frame theo thứ tự để reviewer thấy điểm nối loop."""
    if columns < 2 or alpha.width % columns:
        return _render_on_background(alpha, DARK)
    cell_width = alpha.width // columns
    preview = Image.new("RGBA", (cell_width * columns, alpha.height), DARK)
    for index in range(columns):
        preview.alpha_composite(alpha.crop((index * cell_width, 0, (index + 1) * cell_width, alpha.height)), (index * cell_width, 0))
    return preview


def review(art_root: Path, review_path: Path, reviewer: str) -> dict:
    evidence_root = art_root / "qa"
    alpha_root = evidence_root / "alpha"
    preview_root = evidence_root / "previews"
    report_root = evidence_root / "reports"
    for directory in (alpha_root, preview_root, report_root):
        directory.mkdir(parents=True, exist_ok=True)
    reviewed: list[dict] = []
    contact_entries: list[tuple[str, Image.Image]] = []
    for source in sorted(art_root.glob("*.png")):
        with Image.open(source) as raw:
            alpha = extract_alpha(raw)
        columns = _expected_columns(source.name)
        ratio, spread = _border_metrics(raw)
        issues = check_cell_boundaries(raw, columns)
        if ratio < MIN_BORDER_CHROMA_RATIO:
            issues.append("border_chroma_not_detected")
        if spread > MAX_BORDER_CHROMA_SPREAD:
            issues.append("background_gradient_exceeds_static_matte_tolerance")
        issues.extend(_anchors_and_loop(alpha, columns))
        alpha_path = alpha_root / source.name.replace("_chroma", "_alpha")
        alpha.save(alpha_path)
        preview_path = preview_root / source.name.replace("_chroma.png", "_animation-preview.png")
        _preview(alpha, columns).save(preview_path)
        asset_report = {
            "asset": source.name,
            "source_sha256": sha256(source),
            "status": "PASS" if not issues else "FAIL",
            "issues": issues,
            "checks": {"columns": columns, "border_chroma_ratio": round(ratio, 6), "border_chroma_spread": spread},
            "alpha_path": str(alpha_path),
            "preview_path": str(preview_path),
        }
        report_path = report_root / source.name.replace(".png", ".json")
        report_path.write_text(json.dumps(asset_report, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
        reviewed.append({
            "asset": source.name,
            "status": asset_report["status"],
            "source_sha256": asset_report["source_sha256"],
            "evidence_path": str(report_path),
            "evidence_sha256": sha256(report_path),
            "issues": issues,
        })
        contact_entries.append((source.name, alpha))
    light = evidence_root / "contact-sheet-light.png"
    dark = evidence_root / "contact-sheet-dark.png"
    _contact_sheet(contact_entries, LIGHT).save(light)
    _contact_sheet(contact_entries, DARK).save(dark)
    document = {
        "review_id": "vertical-slice-v1",
        "reviewer": reviewer,
        "timestamp": datetime.now(timezone.utc).isoformat().replace("+00:00", "Z"),
        "source_root": str(art_root),
        "evidence": [
            {"path": str(light), "sha256": sha256(light)},
            {"path": str(dark), "sha256": sha256(dark)},
        ],
        "assets": reviewed,
        "summary": {"total": len(reviewed), "pass": sum(x["status"] == "PASS" for x in reviewed), "fail": sum(x["status"] == "FAIL" for x in reviewed)},
    }
    review_path.parent.mkdir(parents=True, exist_ok=True)
    review_path.write_text(json.dumps(document, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    return document


def main() -> int:
    parser = argparse.ArgumentParser(description="QA PNG nguồn vertical slice (không sửa nguồn)")
    parser.add_argument("--art-root", type=Path, required=True)
    parser.add_argument("--review", type=Path, required=True)
    parser.add_argument("--reviewer", default="Codex Task 6A")
    args = parser.parse_args()
    result = review(args.art_root, args.review, args.reviewer)
    print(json.dumps(result["summary"], ensure_ascii=False))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
