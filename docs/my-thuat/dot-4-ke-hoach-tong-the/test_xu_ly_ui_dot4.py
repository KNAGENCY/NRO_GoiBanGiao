from pathlib import Path
import sys

import pytest
from PIL import Image

sys.path.insert(0, str(Path(__file__).parent))

from xu_ly_ui_dot4 import ValidationError, mask_edge_connected_checkerboard, process_files


def checkerboard(size: tuple[int, int]) -> Image.Image:
    image = Image.new("RGBA", size)
    pixels = image.load()
    for y in range(size[1]):
        for x in range(size[0]):
            value = 255 if ((x // 4) + (y // 4)) % 2 else 242
            pixels[x, y] = (value, value, value, 255)
    return image


def pixel_columns_equal(image: Image.Image, left: int, right: int) -> bool:
    return all(image.getpixel((left, y)) == image.getpixel((right, y)) for y in range(image.height))


def pixel_rows_equal(image: Image.Image, top: int, bottom: int) -> bool:
    return all(image.getpixel((x, top)) == image.getpixel((x, bottom)) for x in range(image.width))


def test_mask_only_removes_neutral_checkerboard_connected_from_edge() -> None:
    image = checkerboard((20, 20))
    # Vật thể màu nằm giữa không được xoá; một điểm trắng kín cũng không phải nền.
    for y in range(5, 15):
        for x in range(5, 15):
            image.putpixel((x, y), (232, 114, 42, 255))
    image.putpixel((10, 10), (255, 255, 255, 255))

    changed = mask_edge_connected_checkerboard(image)

    assert changed > 0
    assert image.getpixel((0, 0))[3] == 0
    assert image.getpixel((6, 6)) == (232, 114, 42, 255)
    assert image.getpixel((10, 10)) == (255, 255, 255, 255)


def test_processes_three_assets_with_backup_seams_hud_and_idempotence(tmp_path: Path) -> None:
    root = tmp_path / "source"
    backup = tmp_path / "backup"
    root.mkdir()

    hud = checkerboard((620, 164))
    hud.paste((83, 45, 27, 255), (20, 16, 300, 132))
    hud.save(root / "ui_hud_khung_trai_x4.png")

    popup = checkerboard((40, 320))
    popup.paste((58, 36, 24, 255), (8, 40, 32, 80))
    popup.paste((58, 36, 24, 255), (8, 120, 32, 200))
    popup.paste((58, 36, 24, 255), (8, 240, 32, 280))
    popup.save(root / "ui_khung_popup_9o_x4.png")

    button = checkerboard((120, 96))
    button.paste((232, 114, 42, 255), (8, 16, 112, 80))
    button.save(root / "ui_nut_thuong_3manh_x4.png")

    first = process_files(root, backup)
    assert {item.name for item in first} == {
        "ui_hud_khung_trai_x4.png",
        "ui_khung_popup_9o_x4.png",
        "ui_nut_thuong_3manh_x4.png",
    }
    assert all(item.changed for item in first)
    assert all((backup / item.name).is_file() for item in first)

    hud_after = Image.open(root / "ui_hud_khung_trai_x4.png").convert("RGBA")
    for box in ((332, 20, 596, 60), (332, 80, 568, 104)):
        assert hud_after.crop(box).getchannel("A").getextrema() == (255, 255)
        assert hud_after.getpixel((box[0] + 4, box[1] + 4)) == (27, 20, 16, 255)

    popup_after = Image.open(root / "ui_khung_popup_9o_x4.png").convert("RGBA")
    assert pixel_columns_equal(popup_after.crop((0, 40, 40, 80)), 0, 39)
    assert pixel_rows_equal(popup_after.crop((0, 120, 40, 160)), 0, 39)
    assert pixel_rows_equal(popup_after.crop((0, 160, 40, 200)), 0, 39)
    assert pixel_columns_equal(popup_after.crop((0, 240, 40, 280)), 0, 39)

    button_after = Image.open(root / "ui_nut_thuong_3manh_x4.png").convert("RGBA")
    assert pixel_columns_equal(button_after, 40, 79)
    assert button_after.getchannel("A").getextrema()[0] == 0

    bytes_before_second_run = {path.name: path.read_bytes() for path in root.glob("*.png")}
    second = process_files(root, backup)
    assert all(not item.changed for item in second)
    assert bytes_before_second_run == {path.name: path.read_bytes() for path in root.glob("*.png")}


def test_refuses_non_checkerboard_input_before_changing_it(tmp_path: Path) -> None:
    root = tmp_path / "source"
    root.mkdir()
    checkerboard((620, 164)).save(root / "ui_hud_khung_trai_x4.png")
    checkerboard((40, 320)).save(root / "ui_khung_popup_9o_x4.png")
    bad = Image.new("RGBA", (120, 96), (30, 20, 10, 255))
    target = root / "ui_nut_thuong_3manh_x4.png"
    bad.save(target)
    original = target.read_bytes()

    with pytest.raises(ValidationError, match="checkerboard"):
        process_files(root, tmp_path / "backup")

    assert target.read_bytes() == original
