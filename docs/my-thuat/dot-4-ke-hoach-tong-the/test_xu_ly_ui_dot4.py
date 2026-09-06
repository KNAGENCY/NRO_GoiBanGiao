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


def max_adjacent_jump(image: Image.Image, axis: str, index: int) -> int:
    if axis == "x":
        return max(abs(a - b) for y in range(image.height) for a, b in zip(image.getpixel((index, y)), image.getpixel((index + 1, y))))
    return max(abs(a - b) for x in range(image.width) for a, b in zip(image.getpixel((x, index)), image.getpixel((x, index + 1))))


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


def test_mask_preserves_neutral_art_highlight_touching_edge() -> None:
    image = checkerboard((64, 64))
    highlight = (244, 244, 244, 255)
    # Highlight thật đi ra tận mép phải và nằm giữa nền caro 255/242. Màu
    # trung tính + nối biên không đủ để kết luận đây là nền.
    for y in range(46, 51):
        for x in range(52, 64):
            image.putpixel((x, y), highlight)

    changed = mask_edge_connected_checkerboard(image)

    assert changed > 0
    assert image.getpixel((0, 0))[3] == 0  # ô nền xám 242 thật bị xoá
    assert image.getpixel((4, 0))[3] == 0  # ô nền trắng 255 thật bị xoá
    assert image.getpixel((32, 32))[3] == 0  # nền thật trong ruột ảnh cũng bị xoá
    assert image.getpixel((55, 48)) == highlight
    assert image.getpixel((63, 48)) == highlight


def test_mask_refuses_two_neutral_bands_without_checkerboard_pattern() -> None:
    image = Image.new("RGBA", (64, 64), (255, 255, 255, 255))
    image.paste((242, 242, 242, 255), (0, 0, 32, 64))
    before = image.tobytes()

    with pytest.raises(ValidationError, match="spatial"):
        mask_edge_connected_checkerboard(image)

    assert image.tobytes() == before


def test_mask_refuses_checkerboard_border_around_solid_gray_interior() -> None:
    image = Image.new("RGBA", (120, 96), (242, 242, 242, 255))
    border_pattern = checkerboard(image.size)
    for x in range(image.width):
        image.putpixel((x, 0), border_pattern.getpixel((x, 0)))
        image.putpixel((x, image.height - 1), border_pattern.getpixel((x, image.height - 1)))
    for y in range(1, image.height - 1):
        image.putpixel((0, y), border_pattern.getpixel((0, y)))
        image.putpixel((image.width - 1, y), border_pattern.getpixel((image.width - 1, y)))
    before = image.tobytes()

    with pytest.raises(ValidationError, match="checkerboard spatial 2D"):
        mask_edge_connected_checkerboard(image)

    assert image.tobytes() == before


def test_processes_three_assets_with_backup_seams_hud_and_idempotence(tmp_path: Path) -> None:
    root = tmp_path / "source"
    backup = tmp_path / "backup"
    root.mkdir()

    hud = checkerboard((620, 164))
    hud.paste((83, 45, 27, 255), (20, 16, 300, 132))
    hud.save(root / "ui_hud_khung_trai_x4.png")

    popup = checkerboard((40, 320))
    for y in range(40, 80):
        for x in range(40):
            popup.putpixel((x, y), (30 + x * 4, 36, 24, 255))
    popup.paste((58, 36, 24, 255), (8, 120, 32, 200))
    popup.paste((58, 36, 24, 255), (8, 240, 32, 280))
    popup.save(root / "ui_khung_popup_9o_x4.png")

    button = checkerboard((120, 96))
    for y in range(16, 80):
        for x in range(8, 112):
            button.putpixel((x, y), (140 + x, 114, 42, 255))
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
    assert hud_after.crop((20, 16, 300, 132)).getchannel("A").getextrema() == (255, 255)
    assert set(hud_after.crop((20, 16, 300, 132)).get_flattened_data()) == {(58, 36, 24, 255)}

    popup_after = Image.open(root / "ui_khung_popup_9o_x4.png").convert("RGBA")
    assert pixel_columns_equal(popup_after.crop((0, 40, 40, 80)), 0, 39)
    assert pixel_rows_equal(popup_after.crop((0, 120, 40, 160)), 0, 39)
    assert pixel_rows_equal(popup_after.crop((0, 160, 40, 200)), 0, 39)
    assert pixel_columns_equal(popup_after.crop((0, 240, 40, 280)), 0, 39)
    assert max_adjacent_jump(popup_after.crop((0, 40, 40, 80)), "x", 38) == 0
    assert max_adjacent_jump(popup_after.crop((0, 120, 40, 160)), "y", 38) == 0
    assert max_adjacent_jump(popup_after.crop((0, 160, 40, 200)), "y", 38) == 0
    assert max_adjacent_jump(popup_after.crop((0, 240, 40, 280)), "x", 38) == 0
    cream = (255, 229, 199, 255)
    assert all(popup_after.getpixel((x, 39)) == cream for x in range(40))
    assert all(popup_after.getpixel((x, 79)) == cream for x in range(40))
    assert all(popup_after.getpixel((x, 119)) == cream for x in range(40))
    assert all(popup_after.getpixel((x, 200)) == cream for x in range(40))
    assert all(popup_after.getpixel((x, 240)) == cream for x in range(40))
    assert all(popup_after.getpixel((x, 280)) == cream for x in range(40))
    assert all(popup_after.getpixel((39, y)) == cream for y in range(40))
    assert all(popup_after.getpixel((0, y)) == cream for y in range(80, 120))
    assert all(popup_after.getpixel((39, y)) == cream for y in range(120, 160))
    assert all(popup_after.getpixel((0, y)) == cream for y in range(160, 200))
    assert all(popup_after.getpixel((39, y)) == cream for y in range(200, 240))
    assert all(popup_after.getpixel((0, y)) == cream for y in range(280, 320))

    button_after = Image.open(root / "ui_nut_thuong_3manh_x4.png").convert("RGBA")
    assert pixel_columns_equal(button_after, 40, 79)
    assert max_adjacent_jump(button_after, "x", 78) == 0
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


def test_refuses_solid_white_even_when_it_touches_every_border(tmp_path: Path) -> None:
    root = tmp_path / "source"
    root.mkdir()
    checkerboard((620, 164)).save(root / "ui_hud_khung_trai_x4.png")
    checkerboard((40, 320)).save(root / "ui_khung_popup_9o_x4.png")
    white = Image.new("RGBA", (120, 96), (255, 255, 255, 255))
    target = root / "ui_nut_thuong_3manh_x4.png"
    white.save(target)
    original = target.read_bytes()

    with pytest.raises(ValidationError, match="signature checkerboard"):
        process_files(root, tmp_path / "backup")

    assert target.read_bytes() == original


def test_refuses_reused_backup_when_opaque_input_hash_changed(tmp_path: Path) -> None:
    root = tmp_path / "source"
    backup = tmp_path / "backup"
    root.mkdir()
    for name, size in {
        "ui_hud_khung_trai_x4.png": (620, 164),
        "ui_khung_popup_9o_x4.png": (40, 320),
        "ui_nut_thuong_3manh_x4.png": (120, 96),
    }.items():
        checkerboard(size).save(root / name)
    process_files(root, backup)

    changed_input = checkerboard((120, 96))
    changed_input.putpixel((60, 48), (10, 20, 30, 255))
    target = root / "ui_nut_thuong_3manh_x4.png"
    changed_input.save(target)
    original = target.read_bytes()

    with pytest.raises(ValidationError, match="backup hash"):
        process_files(root, backup)

    assert target.read_bytes() == original
