"""Kiểm tra closure nền b00..b03 trước khi render contact/full-map."""

from pathlib import Path
from PIL import Image


class VisualError(ValueError):
    pass


def validate_background_layers(layers) -> None:
    if set(layers) != {"b00", "b01", "b02", "b03"}:
        raise VisualError("ASSET-BG-005: background phải đủ b00..b03")


def render_tile_contact_sheet(tiles, output: Path) -> None:
    if set(tiles) != set(range(100)):
        raise VisualError("ASSET-ZOOM-003: contact sheet phải đủ 100 tile")
    width = max(image.width for image in tiles.values())
    height = max(image.height for image in tiles.values())
    sheet = Image.new("RGBA", (width * 10, height * 10))
    for index in range(100):
        image = tiles[index].convert("RGBA")
        sheet.alpha_composite(image, ((index % 10) * width, (index // 10) * height))
    output.parent.mkdir(parents=True, exist_ok=True)
    sheet.save(output, format="PNG")


def render_full_map(tiles, indexes, width: int, height: int, output: Path) -> None:
    if len(indexes) != width * height or any(index not in tiles for index in indexes):
        raise VisualError("ASSET-MAP-004: full-map thiếu tile/index")
    tile_width = max(image.width for image in tiles.values())
    tile_height = max(image.height for image in tiles.values())
    canvas = Image.new("RGBA", (tile_width * width, tile_height * height))
    for position, index in enumerate(indexes):
        canvas.alpha_composite(tiles[index].convert("RGBA"), ((position % width) * tile_width, (position // width) * tile_height))
    output.parent.mkdir(parents=True, exist_ok=True)
    canvas.save(output, format="PNG")
