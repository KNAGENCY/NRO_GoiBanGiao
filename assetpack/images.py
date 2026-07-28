"""Chuẩn hoá sprite theo bốn zoom bằng nội suy gần nhất của Unity."""

from PIL import Image


class ImageError(ValueError):
    pass


def build_zooms(image: Image.Image, _resource_name: str, scales: dict[str, int]) -> dict[str, Image.Image]:
    required = {"x1", "x2", "x3", "x4"}
    if set(scales) != required or any(not isinstance(value, int) or value < 1 for value in scales.values()):
        raise ImageError("ASSET-ZOOM-003: cần đủ x1..x4 với hệ số dương")
    source = image.convert("RGBA")
    return {zoom: source.resize((source.width * scale, source.height * scale), Image.Resampling.NEAREST) for zoom, scale in scales.items()}
