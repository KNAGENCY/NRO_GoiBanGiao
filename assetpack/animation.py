"""Ghép strip action theo slot index, không theo thứ tự input dễ bị đảo cell."""

from PIL import Image


class AnimationError(ValueError):
    pass


REQUIRED_ANIMATION_SHEETS = {
    "hoc-vien-head": {"slot_count": 3, "cell_size": (128, 128)},
    "hoc-vien-body": {"slot_count": 17, "cell_size": (128, 160)},
    "hoc-vien-leg": {"slot_count": 14, "cell_size": (128, 128)},
    "quai-soi-rung": {"slot_count": 10, "cell_size": (96, 96)},
}

REQUIRED_CLIENT_CF = {
    "hoc-vien-head": (
        (0, 1, 7, 8, 11, 12, 17, 18, 19, 20, 21, 22, 23, 24, 30, 31, 32),
        (2, 3, 4, 5, 6, 9, 10, 13, 14, 15, 16, 26, 27, 28, 29),
        (25,),
    ),
    "hoc-vien-body": (
        (0, 1), (2, 15, 28), (3,), (4,), (5,), (6,), (7, 8, 17), (9, 10),
        (11, 21, 32), (12, 18, 22), (13, 26), (14, 27), (16, 29), (19, 30),
        (20, 31), (23, 24), (25,),
    ),
    "hoc-vien-leg": (
        (0, 1, 22), (2,), (3,), (4,), (5,), (6,), (7, 23), (8,), (9,), (10,),
        (11,), (12,), (13, 14, 15, 16, 17, 18, 19, 20, 21, 24), (25,),
    ),
    "quai-soi-rung": tuple((index,) for index in range(10)),
}


def validate_animation_closure(sheets, animations, *, check_paths=None, check_cell_size=None) -> None:
    """Khóa cả bốn strip runtime, kể cả metadata không thể suy từ PNG."""
    if not isinstance(sheets, list) or len(sheets) != len(REQUIRED_ANIMATION_SHEETS) or any(not isinstance(sheet, dict) for sheet in sheets):
        raise AnimationError("ASSET-ANIMATION-008: thiếu required animation_sheets closure")
    sheet_names = [sheet.get("name") for sheet in sheets]
    if any(not isinstance(name, str) for name in sheet_names):
        raise AnimationError("ASSET-ANIMATION-008: sheet name phải là chuỗi exact")
    if len(set(sheet_names)) != len(sheet_names) or set(sheet_names) != set(REQUIRED_ANIMATION_SHEETS):
        raise AnimationError("ASSET-ANIMATION-008: sheet name phải exact và duy nhất")
    if not isinstance(animations, list) or len(animations) != len(REQUIRED_ANIMATION_SHEETS) or any(not isinstance(item, dict) for item in animations):
        raise AnimationError("ASSET-ANIMATION-008: thiếu animation contract tương ứng")
    animation_names = [item.get("name") for item in animations]
    if any(not isinstance(name, str) for name in animation_names):
        raise AnimationError("ASSET-ANIMATION-008: animation name phải là chuỗi exact")
    if len(set(animation_names)) != len(animation_names) or set(animation_names) != set(REQUIRED_ANIMATION_SHEETS):
        raise AnimationError("ASSET-ANIMATION-008: animation name phải exact và duy nhất")
    by_animation = {item["name"]: item for item in animations}
    for sheet in sheets:
        name = sheet["name"]
        expected = REQUIRED_ANIMATION_SHEETS[name]
        slot_count = sheet.get("slot_count")
        if not isinstance(slot_count, int) or isinstance(slot_count, bool) or slot_count != expected["slot_count"]:
            raise AnimationError(f"ASSET-ANIMATION-008: {name} slot_count sai")
        cell_size = sheet.get("cell_size")
        if not isinstance(cell_size, (list, tuple)) or tuple(cell_size) != expected["cell_size"]:
            raise AnimationError(f"ASSET-ANIMATION-008: {name} sai slot count/cell size")
        cells = sheet.get("cells")
        if not isinstance(cells, list) or len(cells) != expected["slot_count"]:
            raise AnimationError(f"ASSET-ANIMATION-008: {name} thiếu frame closure")
        cell_slots = [cell.get("slot_index") if isinstance(cell, dict) else None for cell in cells]
        if any(not isinstance(slot, int) or isinstance(slot, bool) for slot in cell_slots) or set(cell_slots) != set(range(expected["slot_count"])):
            raise AnimationError(f"ASSET-ANIMATION-008: {name} slot không đủ/duy nhất")
        contract = by_animation[name]
        slot_count = contract.get("slot_count")
        if not isinstance(slot_count, int) or isinstance(slot_count, bool) or slot_count != expected["slot_count"]:
            raise AnimationError(f"ASSET-ANIMATION-008: {name} slot_count contract sai")
        frames = contract.get("frames")
        if not isinstance(frames, list) or len(frames) != expected["slot_count"]:
            raise AnimationError(f"ASSET-ANIMATION-008: {name} frame count contract sai")
        frame_slots = [frame.get("slot_index") if isinstance(frame, dict) else None for frame in frames]
        if any(not isinstance(slot, int) or isinstance(slot, bool) for slot in frame_slots) or set(frame_slots) != set(range(expected["slot_count"])):
            raise AnimationError(f"ASSET-ANIMATION-008: {name} contract slot sai")
        for frame in frames:
            client_cf = frame.get("client_cf") if isinstance(frame, dict) else None
            slot = frame.get("slot_index") if isinstance(frame, dict) else None
            expected_cf = REQUIRED_CLIENT_CF[name][slot] if isinstance(slot, int) and not isinstance(slot, bool) and 0 <= slot < expected["slot_count"] else None
            if not isinstance(client_cf, list) or tuple(client_cf) != expected_cf:
                raise AnimationError(f"ASSET-ANIMATION-008: {name} client_cf malformed")
        for cell in cells:
            path = cell.get("path") if isinstance(cell, dict) else None
            if not isinstance(path, str) or not path:
                raise AnimationError(f"ASSET-ANIMATION-008: {name} thiếu path cell")
            if check_paths is not None and not check_paths(path).is_file():
                raise AnimationError(f"ASSET-ANIMATION-008: thiếu cell {path}")
            if check_cell_size is not None and not check_cell_size(path, expected["cell_size"]):
                raise AnimationError(f"ASSET-ANIMATION-008: {name} cell vượt canvas/boundary")


def compose_action_sheet(cells, slot_count: int, cell_size: tuple[int, int]) -> Image.Image:
    if slot_count not in {3, 10, 14, 17}:
        raise AnimationError("ASSET-ANIMATION-008: slot count phải là 3/17/14/10")
    if len(cells) != slot_count or {slot for slot, _ in cells} != set(range(slot_count)):
        raise AnimationError("ASSET-ANIMATION-008: slot index phải đủ và duy nhất")
    width, height = cell_size
    sheet = Image.new("RGBA", (width * slot_count, height))
    for slot, cell in cells:
        if cell.size != cell_size:
            raise AnimationError("ASSET-ANIMATION-008: canvas cell không đúng")
        sheet.alpha_composite(cell.convert("RGBA"), (slot * width, 0))
    return sheet
