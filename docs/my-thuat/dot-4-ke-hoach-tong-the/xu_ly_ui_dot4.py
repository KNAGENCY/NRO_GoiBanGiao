"""Sửa có kiểm soát ba mẫu UI Đợt 4 bị image generator nhúng nền caro.

Chỉ nhận đúng ba tên/cỡ đã đo. Nền bị xoá phải là vùng xám-trắng trung tính
nối từ mép canvas; giấy/gỗ nằm trong vật thể không bao giờ là điểm khởi đầu.
Mọi đầu vào lạ bị từ chối trước khi ghi bất kỳ file nào.
"""

from __future__ import annotations

import argparse
import hashlib
import shutil
from collections import Counter, deque
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
HUD_TEXT_PLATE = (58, 36, 24, 255)
POPUP_CREAM = (255, 229, 199, 255)


class ValidationError(ValueError):
    """Dừng trước khi ghi vì input không giống mẫu generator đã đo."""


@dataclass(frozen=True)
class ProcessResult:
    name: str
    changed: bool
    cleared_pixels: int


@dataclass(frozen=True)
class AxisCadence:
    tile: int
    offset: int
    inverted: bool

    def band(self, index: int) -> bool:
        return bool((((index + self.offset) // self.tile) & 1) ^ self.inverted)


@dataclass(frozen=True)
class CheckerboardPattern:
    gray: tuple[int, int, int]
    bright: tuple[int, int, int]
    horizontal: AxisCadence
    vertical: AxisCadence
    phase: bool

    def expected_bright(self, x: int, y: int) -> bool:
        return self.horizontal.band(x) ^ self.vertical.band(y) ^ self.phase

    def matches(self, pixel: tuple[int, int, int, int], x: int, y: int) -> bool:
        """Chỉ xoá đúng màu, nhịp và pha caro đã chứng minh trên ảnh này."""
        expected = self.bright if self.expected_bright(x, y) else self.gray
        return pixel[3] == 255 and pixel[:3] == expected


def _is_neutral_candidate(pixel: tuple[int, int, int, int]) -> bool:
    """Lọc ứng viên để đo; riêng màu trung tính chưa đủ để kết luận là nền."""
    red, green, blue, alpha = pixel
    return alpha == 255 and min(red, green, blue) >= 230 and max(red, green, blue) - min(red, green, blue) <= 8


def _brightness(pixel: tuple[int, int, int, int] | tuple[int, int, int]) -> int:
    return sum(pixel[:3]) // 3


def _border_lines(image: Image.Image) -> tuple[list[tuple[int, int, int, int]], ...]:
    width, height = image.size
    return (
        [image.getpixel((x, 0)) for x in range(width)],
        [image.getpixel((x, height - 1)) for x in range(width)],
        [image.getpixel((0, y)) for y in range(height)],
        [image.getpixel((width - 1, y)) for y in range(height)],
    )


def _transition_count(line: list[tuple[int, int, int, int]], split: int) -> int:
    """Đếm nhịp đổi sáng/xám liên tiếp; vật thể màu ngắt chuỗi thay vì tạo nhịp giả."""
    transitions = 0
    previous: bool | None = None
    for pixel in line:
        if not _is_neutral_candidate(pixel):
            previous = None
            continue
        bright = _brightness(pixel) > split
        if previous is not None and bright != previous:
            transitions += 1
        previous = bright
    return transitions


def _fit_axis_cadence(
    line: list[tuple[int, int, int, int]], split: int
) -> tuple[AxisCadence, float] | None:
    """Khớp nhịp ô vuông trên một cạnh, cho phép cạnh bị artwork che một đoạn."""
    samples = [
        (index, _brightness(pixel) > split)
        for index, pixel in enumerate(line)
        if _is_neutral_candidate(pixel)
    ]
    if len(samples) < 8:
        return None

    best: tuple[int, int, AxisCadence] | None = None
    max_tile = min(24, max(1, len(line) // 4))
    for tile in range(1, max_tile + 1):
        for offset in range(tile):
            for inverted in (False, True):
                cadence = AxisCadence(tile=tile, offset=offset, inverted=inverted)
                hits = sum(cadence.band(index) == bright for index, bright in samples)
                candidate = (hits, -tile, cadence)
                if best is None or candidate[:2] > best[:2]:
                    best = candidate

    assert best is not None
    hits, _, cadence = best
    return cadence, hits / len(samples)


def _field_evidence(
    samples: list[tuple[int, int, tuple[int, int, int, int]]],
    size: tuple[int, int],
    pattern: CheckerboardPattern,
) -> tuple[int, int]:
    """Đo đúng màu + đúng pha trên toàn mặt ảnh, không chỉ bốn đường biên."""
    width, height = size
    matches = 0
    occupied_bins: set[tuple[int, int]] = set()
    for x, y, pixel in samples:
        if pattern.matches(pixel, x, y):
            matches += 1
            occupied_bins.add((min(3, x * 4 // width), min(3, y * 4 // height)))
    return matches, len(occupied_bins)


def _detect_checkerboard_edge_signature(
    image: Image.Image,
) -> tuple[tuple[int, int, int], tuple[int, int, int], int, tuple[list[tuple[int, int, int, int]], ...]]:
    """Học palette và xác nhận bốn mép có nhịp sáng/xám theo cả hai trục."""
    lines = _border_lines(image)
    border = [pixel for line in lines for pixel in line]
    gray_counts: Counter[tuple[int, int, int]] = Counter()
    bright_counts: Counter[tuple[int, int, int]] = Counter()
    for pixel in border:
        if not _is_neutral_candidate(pixel):
            continue
        color = pixel[:3]
        value = _brightness(pixel)
        if value <= 248:
            gray_counts[color] += 1
        elif value >= 250:
            bright_counts[color] += 1

    if not gray_counts or not bright_counts:
        raise ValidationError("khong co du hai mau cho signature checkerboard spatial")
    gray = gray_counts.most_common(1)[0][0]
    bright = bright_counts.most_common(1)[0][0]
    gray_value = _brightness(gray)
    bright_value = _brightness(bright)
    if bright_value - gray_value < 6:
        raise ValidationError("hai mau khong tach du cho signature checkerboard spatial")

    split = (gray_value + bright_value) // 2
    horizontal = max(_transition_count(lines[0], split), _transition_count(lines[1], split))
    vertical = max(_transition_count(lines[2], split), _transition_count(lines[3], split))
    if horizontal < 4 or vertical < 4:
        raise ValidationError("hai dai mau khong tao signature checkerboard spatial")
    return gray, bright, split, lines


def _detect_checkerboard_pattern(image: Image.Image) -> CheckerboardPattern:
    """Học hai màu nền rồi chứng minh cadence/pha trên toàn mặt ảnh.

    Fail-closed ở đây quan trọng hơn xoá sạch vài pixel xám do resize: màu nào
    không thuộc đúng hai cụm caro đã chứng minh sẽ được giữ nguyên như artwork.
    """
    gray, bright, split, lines = _detect_checkerboard_edge_signature(image)

    width, height = image.size
    horizontal_fits = [
        fit
        for line in (lines[0], lines[1])
        if (fit := _fit_axis_cadence(line, split)) is not None and fit[1] >= 0.85
    ]
    vertical_fits = [
        fit
        for line in (lines[2], lines[3])
        if (fit := _fit_axis_cadence(line, split)) is not None and fit[1] >= 0.85
    ]
    if not horizontal_fits or not vertical_fits:
        raise ValidationError("khong khop duoc nhip checkerboard spatial 2D")

    horizontal_cadence, _ = max(horizontal_fits, key=lambda item: item[1])
    vertical_cadence, _ = max(vertical_fits, key=lambda item: item[1])
    samples = [
        (index % width, index // width, pixel)
        for index, pixel in enumerate(image.get_flattened_data())
        if _is_neutral_candidate(pixel)
    ]
    best: tuple[float, CheckerboardPattern, int, int] | None = None
    for phase in (False, True):
        pattern = CheckerboardPattern(
            gray=gray,
            bright=bright,
            horizontal=horizontal_cadence,
            vertical=vertical_cadence,
            phase=phase,
        )
        matches, occupied_bins = _field_evidence(samples, image.size, pattern)
        ratio = matches / len(samples) if samples else 0.0
        candidate = (ratio, pattern, matches, occupied_bins)
        if best is None or candidate[0] > best[0]:
            best = candidate

    assert best is not None
    ratio, pattern, matches, occupied_bins = best
    minimum_coverage = max(32, width * height // 10)
    if ratio < 0.85 or matches < minimum_coverage or occupied_bins < 8:
        raise ValidationError("khong co du checkerboard spatial 2D dung pha tren mat anh")
    return pattern


def mask_edge_connected_checkerboard(image: Image.Image) -> int:
    """Đổi nền caro đã nhận dạng nối từ mép thành alpha 0, không flood qua vật thể."""
    pattern = _detect_checkerboard_pattern(image)
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
        if (x, y) in seen or not pattern.matches(pixels[x, y], x, y):
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


def _validate_opaque_checkerboard(
    path: Path, image: Image.Image, *, require_spatial_2d: bool = True
) -> None:
    _check_name_and_size(path, image)
    if image.getchannel("A").getextrema() != (255, 255):
        raise ValidationError(f"{path.name} da co alpha; khong dung mau input caro duc")
    width, height = image.size
    border = [image.getpixel((x, 0)) for x in range(width)]
    border += [image.getpixel((x, height - 1)) for x in range(width)]
    border += [image.getpixel((0, y)) for y in range(1, height - 1)]
    border += [image.getpixel((width - 1, y)) for y in range(1, height - 1)]
    neutral_border = sum(_is_neutral_candidate(pixel) for pixel in border)
    if neutral_border * 100 < len(border) * 70:
        raise ValidationError(f"{path.name} khong co du nen checkerboard trung tinh noi tu bien")
    # Nền preview thật có ít nhất hai dải sáng/xám. Một tấm trắng đặc hoặc
    # highlight trắng chạm biên không phải checkerboard và không được xoá.
    brightness = [_brightness(pixel) for pixel in border if _is_neutral_candidate(pixel)]
    bright = sum(value >= 250 for value in brightness)
    gray = sum(230 <= value <= 248 for value in brightness)
    required_band = max(4, len(border) // 100)
    if bright < required_band or gray < required_band:
        raise ValidationError(f"{path.name} khong co signature checkerboard sang/xam")
    if require_spatial_2d:
        _detect_checkerboard_pattern(image)
    else:
        _detect_checkerboard_edge_signature(image)


def _same_column(image: Image.Image, left: int, right: int) -> bool:
    return all(image.getpixel((left, y)) == image.getpixel((right, y)) for y in range(image.height))


def _same_row(image: Image.Image, top: int, bottom: int) -> bool:
    return all(image.getpixel((x, top)) == image.getpixel((x, bottom)) for x in range(image.width))


def _lock_popup_seams(image: Image.Image) -> None:
    # Cells 2/7 tile ngang; 4/5 tile dọc. Chuẩn hoá cả trục lặp để không
    # đẩy vết gãy từ pixel biên vào pixel áp chót.
    for top in (40, 240):
        for y in range(top, top + 40):
            sample = image.getpixel((20, y))
            for x in range(40):
                image.putpixel((x, y), sample)
    for top in (120, 160):
        for x in range(40):
            sample = image.getpixel((x, top + 20))
            for y in range(top, top + 40):
                image.putpixel((x, y), sample)

    # Mặt trong của đủ tám cell phải nối trực tiếp với ruột #FFE5C7.
    for x in range(40):
        image.putpixel((x, 39), POPUP_CREAM)   # cell 1, cạnh trong dưới
        image.putpixel((x, 79), POPUP_CREAM)   # cell 2, cạnh trong dưới
        image.putpixel((x, 119), POPUP_CREAM)  # cell 3, cạnh trong dưới
        image.putpixel((x, 200), POPUP_CREAM)  # cell 6, cạnh trong trên
        image.putpixel((x, 240), POPUP_CREAM)  # cell 7, cạnh trong trên
        image.putpixel((x, 280), POPUP_CREAM)  # cell 8, cạnh trong trên
    for y in range(40):
        image.putpixel((39, y), POPUP_CREAM)         # cell 1, cạnh trong phải
        image.putpixel((0, 80 + y), POPUP_CREAM)     # cell 3, cạnh trong trái
        image.putpixel((39, 120 + y), POPUP_CREAM)   # cell 4, cạnh trong phải
        image.putpixel((0, 160 + y), POPUP_CREAM)    # cell 5, cạnh trong trái
        image.putpixel((39, 200 + y), POPUP_CREAM)   # cell 6, cạnh trong phải
        image.putpixel((0, 280 + y), POPUP_CREAM)    # cell 8, cạnh trong trái


def _lock_button_middle_seam(image: Image.Image) -> None:
    for y in range(96):
        sample = image.getpixel((60, y))
        for x in range(40, 80):
            image.putpixel((x, y), sample)


def _normalize_hud_troughs(image: Image.Image) -> None:
    for left, top, right, bottom in TROUGHS:
        for y in range(top, bottom):
            for x in range(left, right):
                shadow = x - left < 4 or right - 1 - x < 4 or y - top < 4 or bottom - 1 - y < 4
                image.putpixel((x, y), TROUGH_SHADOW if shadow else TROUGH_COLOR)


def _normalize_hud_text_plate(image: Image.Image) -> None:
    # Client đè chân dung/tên/sức mạnh lên vùng này; nền một màu đục giúp chữ
    # đọc được và tránh alpha/carô lọt vào ô chân dung 24..128.
    for y in range(16, 132):
        for x in range(20, 300):
            image.putpixel((x, y), HUD_TEXT_PLATE)


def _verify_processed(path: Path, image: Image.Image) -> None:
    _check_name_and_size(path, image)
    alpha = image.getchannel("A")
    if alpha.getextrema()[0] != 0:
        raise ValidationError(f"{path.name} khong co alpha 0 sau xu ly")
    if path.name == "ui_hud_khung_trai_x4.png":
        plate = image.crop((20, 16, 300, 132))
        if plate.getchannel("A").getextrema() != (255, 255) or set(plate.get_flattened_data()) != {HUD_TEXT_PLATE}:
            raise ValidationError("vung chu/chan dung HUD khong plain opaque")
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
        horizontal_inner = (39, 79, 119, 200, 240, 280)
        vertical_inner = ((39, 0, 40), (0, 80, 120), (39, 120, 160), (0, 160, 200), (39, 200, 240), (0, 280, 320))
        if any(image.getpixel((x, y)) != POPUP_CREAM for y in horizontal_inner for x in range(40)):
            raise ValidationError("popup inner ngang khong ket thuc #FFE5C7")
        if any(image.getpixel((x, y)) != POPUP_CREAM for x, top, bottom in vertical_inner for y in range(top, bottom)):
            raise ValidationError("popup inner doc khong ket thuc #FFE5C7")
    elif path.name == "ui_nut_thuong_3manh_x4.png" and not _same_column(image, 40, 79):
        raise ValidationError("manh giua nut khong tile ngang")


def _validate_processed_lineage(path: Path, backup_dir: Path) -> None:
    """Ảnh đã alpha chỉ được nâng cấp nếu truy được backup caro cùng tên.

    Backup round 2 là ảnh generator đã resize nên cadence không đều tới từng pixel;
    nó chỉ được đọc để chứng minh lineage. Input đục sắp bị mask vẫn bắt buộc qua
    spatial 2D nghiêm ngặt trước khi bất kỳ pixel nào được đổi.
    """
    backup = backup_dir / path.name
    if not backup.is_file():
        raise ValidationError(f"{path.name} da xu ly nhung thieu backup de kiem chung nguon")
    with Image.open(backup) as backup_image:
        _validate_opaque_checkerboard(path, backup_image.convert("RGBA"), require_spatial_2d=False)


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
            _validate_processed_lineage(path, backup_dir)
        else:
            _validate_opaque_checkerboard(path, images[path])
            backup = backup_dir / path.name
            if backup.is_file() and hashlib.sha256(backup.read_bytes()).digest() != hashlib.sha256(path.read_bytes()).digest():
                raise ValidationError(f"{path.name} backup hash khac input opaque moi")

    backup_dir.mkdir(parents=True, exist_ok=True)
    results: list[ProcessResult] = []
    for path in paths:
        image = images[path]
        before = image.tobytes()
        cleared = 0
        if not already_processed[path]:
            backup = backup_dir / path.name
            if not backup.exists():
                shutil.copy2(path, backup)
            cleared = mask_edge_connected_checkerboard(image)
            if cleared == 0:
                raise ValidationError(f"{path.name} khong tim thay pixel checkerboard de xoa")
        if path.name == "ui_hud_khung_trai_x4.png":
            _normalize_hud_text_plate(image)
            _normalize_hud_troughs(image)
        elif path.name == "ui_khung_popup_9o_x4.png":
            _lock_popup_seams(image)
        else:
            _lock_button_middle_seam(image)
        _verify_processed(path, image)
        changed = image.tobytes() != before
        if changed:
            image.save(path, "PNG", optimize=True)
        results.append(ProcessResult(path.name, changed, cleared))
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
