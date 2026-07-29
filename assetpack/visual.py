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


def render_hud_golden(hud, output: Path) -> None:
    """Render HP/chakra ở 25/50/100%, clip fill thay vì co texture."""

    expected = {
        "naruto-hp-bg",
        "naruto-hp-fill",
        "naruto-chakra-bg",
        "naruto-chakra-fill",
    }
    if set(hud) != expected:
        raise VisualError("ASSET-UI-007: HUD golden phải đủ bốn resource")
    rows = (("naruto-hp-bg", "naruto-hp-fill"), ("naruto-chakra-bg", "naruto-chakra-fill"))
    sizes = {hud[name].size for name in expected}
    if len(sizes) != 1:
        raise VisualError("ASSET-UI-007: HUD background/fill phải cùng kích thước")
    width, height = sizes.pop()
    golden = Image.new("RGBA", (width * 3, height * 2))
    for column, ratio in enumerate((0.25, 0.5, 1.0)):
        fill_width = round(width * ratio)
        for row, (background_name, fill_name) in enumerate(rows):
            state = hud[background_name].convert("RGBA").copy()
            fill = hud[fill_name].convert("RGBA")
            state.alpha_composite(fill.crop((0, 0, fill_width, height)), (0, 0))
            golden.alpha_composite(state, (column * width, row * height))
    output.parent.mkdir(parents=True, exist_ok=True)
    golden.save(output, format="PNG", optimize=False, compress_level=9)


def render_iruka_composite_contact_sheet(layers, output: Path) -> None:
    """Ghép part theo layer/neo đáy rồi đặt cạnh nhau trên nền sáng và tối."""

    if set(layers) != {"head", "body", "leg"}:
        raise VisualError("ASSET-PART-001: Iruka composite phải đủ head/body/leg")
    width = max(image.width for image in layers.values())
    height = max(image.height for image in layers.values())
    composite = Image.new("RGBA", (width, height))
    for name in ("leg", "body", "head"):
        layer = layers[name].convert("RGBA")
        composite.alpha_composite(
            layer,
            ((width - layer.width) // 2, height - layer.height),
        )
    contact = Image.new("RGBA", (width * 2, height))
    for column, color in enumerate(
        ((245, 245, 245, 255), (24, 24, 24, 255))
    ):
        background = Image.new("RGBA", (width, height), color)
        background.alpha_composite(composite)
        contact.alpha_composite(background, (column * width, 0))
    output.parent.mkdir(parents=True, exist_ok=True)
    contact.save(output, format="PNG", optimize=False, compress_level=9)
