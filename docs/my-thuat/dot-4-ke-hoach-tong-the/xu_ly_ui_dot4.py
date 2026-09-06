"""Sửa có kiểm soát ba mẫu UI Đợt 4 bị image generator nhúng nền caro.

Chỉ nhận đúng ba tên/cỡ đã đo. Nền bị xoá phải là vùng xám-trắng trung tính
nối từ mép canvas; giấy/gỗ nằm trong vật thể không bao giờ là điểm khởi đầu.
Mọi đầu vào lạ bị từ chối trước khi ghi bất kỳ file nào.
"""

from __future__ import annotations

import argparse
import shutil
from collections import deque
from dataclasses import dataclass
from pathlib import Path

from PIL import Image


SPECS = {
    "ui_hud_khung_trai_x4.png": (620, 164),
    "ui_khung_popup_9o_x4.png": (40, 320),
    "ui_nut_thuong_3manh_x4.png": (120, 96),
}
TROUGH_COLOR = (27, 20, 16, 255)  # #1B1410, vùng game tự vẽ máu/chakra.
TROUGH_SHADOW = (18, 15, 13, 255)
TROUGHS = ((332, 20, 596, 60), (332, 80, 568, 104))


class ValidationError(ValueError):
    """Dừng trước khi ghi vì input không giống mẫu generator đã đo."""


@dataclass(frozen=True)
class ProcessResult:
    name: str
    changed: bool
    cleared_pixels: int


def is_border_checkerboard(pixel: tuple[int, int, int, int]) -> bool:
    """Nhận đúng vùng caro sáng, trung tính; không nhận giấy vàng hay gỗ."""
    red, green, blue, alpha = pixel
    return alpha == 255 and min(red, green, blue) >= 230 and max(red, green, blue) - min(red, green, blue) <= 8


def mask_edge_connected_checkerboard(image: Image.Image) -> int:
    """Đổi nền caro trung tính nối từ mép thành alpha 0, không flood qua vật thể."""
    pixels = image.load()
    width, height = image.size
    queue: deque[tuple[int, int]] = deque()
    seen: set[tuple[int, int]] = set()

    for x in range(width):
        queue.extend(((x, 0), (x, height - 1)))
    for y in range(1, height - 1):
        queue.extend(((0, y), (width - 1, y)))

    while queue:
        x, y = queue.popleft()
        if (x, y) in seen or not is_border_checkerboard(pixels[x, y]):
            continue
        seen.add((x, y))
        for next_x, next_y in ((x - 1, y), (x + 1, y), (x, y - 1), (x, y + 1)):
            if 0 <= next_x < width and 0 <= next_y < height:
                queue.append((next_x, next_y))

    for x, y in seen:
        red, green, blue, _ = pixels[x, y]
        pixels[x, y] = (red, green, blue, 0)
    return len(seen)


def _check_name_and_size(path: Path, image: Image.Image) -> None:
    expected = SPECS.get(path.name)
    if expected is None:
        raise ValidationError(f"ten file khong nam trong batch Task 2: {path.name}")
    if image.size != expected:
        raise ValidationError(f"kich thuoc {path.name} la {image.size}, phai la {expected}")


def _validate_opaque_checkerboard(path: Path, image: Image.Image) -> None:
    _check_name_and_size(path, image)
    if image.getchannel("A").getextrema() != (255, 255):
        raise ValidationError(f"{path.name} da co alpha; khong dung mau input caro duc")
    width, height = image.size
    border = [image.getpixel((x, 0)) for x in range(width)]
    border += [image.getpixel((x, height - 1)) for x in range(width)]
    border += [image.getpixel((0, y)) for y in range(1, height - 1)]
    border += [image.getpixel((width - 1, y)) for y in range(1, height - 1)]
    neutral_border = sum(is_border_checkerboard(pixel) for pixel in border)
    if neutral_border * 100 < len(border) * 70:
        raise ValidationError(f"{path.name} khong co du nen checkerboard trung tinh noi tu bien")


def _same_column(image: Image.Image, left: int, right: int) -> bool:
    return all(image.getpixel((left, y)) == image.getpixel((right, y)) for y in range(image.height))


def _same_row(image: Image.Image, top: int, bottom: int) -> bool:
    return all(image.getpixel((x, top)) == image.getpixel((x, bottom)) for x in range(image.width))


def _lock_popup_seams(image: Image.Image) -> None:
    # Cells 2/7 tile ngang; 4/5 tile dọc. Chỉ sao đúng đường biên lặp 1 px.
    for top in (40, 240):
        for y in range(top, top + 40):
            image.putpixel((39, y), image.getpixel((0, y)))
    for top in (120, 160):
        for x in range(40):
            image.putpixel((x, top + 39), image.getpixel((x, top)))


def _lock_button_middle_seam(image: Image.Image) -> None:
    for y in range(96):
        image.putpixel((79, y), image.getpixel((40, y)))


def _normalize_hud_troughs(image: Image.Image) -> None:
    for left, top, right, bottom in TROUGHS:
        for y in range(top, bottom):
            for x in range(left, right):
                shadow = x - left < 4 or right - 1 - x < 4 or y - top < 4 or bottom - 1 - y < 4
                image.putpixel((x, y), TROUGH_SHADOW if shadow else TROUGH_COLOR)


def _verify_processed(path: Path, image: Image.Image) -> None:
    _check_name_and_size(path, image)
    alpha = image.getchannel("A")
    if alpha.getextrema()[0] != 0:
        raise ValidationError(f"{path.name} khong co alpha 0 sau xu ly")
    if path.name == "ui_hud_khung_trai_x4.png":
        for left, top, right, bottom in TROUGHS:
            if any(image.getpixel((x, y)) != TROUGH_COLOR for y in range(top + 4, bottom - 4) for x in range(left + 4, right - 4)):
                raise ValidationError("ranh HUD khong flat #1B1410")
    elif path.name == "ui_khung_popup_9o_x4.png":
        if not _same_column(image.crop((0, 40, 40, 80)), 0, 39):
            raise ValidationError("popup top edge khong tile ngang")
        if not _same_row(image.crop((0, 120, 40, 160)), 0, 39):
            raise ValidationError("popup left edge khong tile doc")
        if not _same_row(image.crop((0, 160, 40, 200)), 0, 39):
            raise ValidationError("popup right edge khong tile doc")
        if not _same_column(image.crop((0, 240, 40, 280)), 0, 39):
            raise ValidationError("popup bottom edge khong tile ngang")
    elif path.name == "ui_nut_thuong_3manh_x4.png" and not _same_column(image, 40, 79):
        raise ValidationError("manh giua nut khong tile ngang")


def process_files(source_dir: Path, backup_dir: Path) -> list[ProcessResult]:
    """Xử lý đúng ba final; lần hai chỉ xác nhận bất biến và không ghi lại."""
    source_dir = Path(source_dir)
    backup_dir = Path(backup_dir)
    paths = [source_dir / name for name in SPECS]
    if any(not path.is_file() for path in paths):
        missing = ", ".join(path.name for path in paths if not path.is_file())
        raise ValidationError(f"thieu final Task 2: {missing}")

    images = {path: Image.open(path).convert("RGBA") for path in paths}
    already_processed = {path: images[path].getchannel("A").getextrema()[0] == 0 for path in paths}
    for path in paths:
        if already_processed[path]:
            _verify_processed(path, images[path])
        else:
            _validate_opaque_checkerboard(path, images[path])

    backup_dir.mkdir(parents=True, exist_ok=True)
    results: list[ProcessResult] = []
    for path in paths:
        if already_processed[path]:
            results.append(ProcessResult(path.name, False, 0))
            continue
        backup = backup_dir / path.name
        if not backup.exists():
            shutil.copy2(path, backup)
        image = images[path]
        cleared = mask_edge_connected_checkerboard(image)
        if cleared == 0:
            raise ValidationError(f"{path.name} khong tim thay pixel checkerboard de xoa")
        if path.name == "ui_hud_khung_trai_x4.png":
            _normalize_hud_troughs(image)
        elif path.name == "ui_khung_popup_9o_x4.png":
            _lock_popup_seams(image)
        else:
            _lock_button_middle_seam(image)
        _verify_processed(path, image)
        image.save(path, "PNG", optimize=True)
        results.append(ProcessResult(path.name, True, cleared))
    return results


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("source_dir", type=Path)
    parser.add_argument("backup_dir", type=Path)
    args = parser.parse_args()
    for result in process_files(args.source_dir, args.backup_dir):
        print(f"{result.name}: changed={result.changed} cleared={result.cleared_pixels}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
