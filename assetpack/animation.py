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


def validate_animation_closure(sheets, animations, *, check_paths=None, check_cell_size=None) -> None:
    """Khóa cả bốn strip runtime, kể cả metadata không thể suy từ PNG."""
    if not isinstance(sheets, list) or {sheet.get("name") for sheet in sheets if isinstance(sheet, dict)} != set(REQUIRED_ANIMATION_SHEETS):
        raise AnimationError("ASSET-ANIMATION-008: thiếu required animation_sheets closure")
    if not isinstance(animations, list) or {item.get("name") for item in animations if isinstance(item, dict)} != set(REQUIRED_ANIMATION_SHEETS):
        raise AnimationError("ASSET-ANIMATION-008: thiếu animation contract tương ứng")
    by_animation = {item["name"]: item for item in animations}
    for sheet in sheets:
        if not isinstance(sheet, dict) or sheet.get("name") not in REQUIRED_ANIMATION_SHEETS:
            raise AnimationError("ASSET-ANIMATION-008: tên sheet không thuộc closure")
        name = sheet["name"]
        expected = REQUIRED_ANIMATION_SHEETS[name]
        if sheet.get("slot_count") != expected["slot_count"] or tuple(sheet.get("cell_size", ())) != expected["cell_size"]:
            raise AnimationError(f"ASSET-ANIMATION-008: {name} sai slot count/cell size")
        cells = sheet.get("cells")
        if not isinstance(cells, list) or len(cells) != expected["slot_count"]:
            raise AnimationError(f"ASSET-ANIMATION-008: {name} thiếu frame closure")
        if {cell.get("slot_index") for cell in cells if isinstance(cell, dict)} != set(range(expected["slot_count"])):
            raise AnimationError(f"ASSET-ANIMATION-008: {name} slot không đủ/duy nhất")
        contract = by_animation[name]
        frames = contract.get("frames")
        if contract.get("slot_count", expected["slot_count"]) != expected["slot_count"] or not isinstance(frames, list) or len(frames) != expected["slot_count"]:
            raise AnimationError(f"ASSET-ANIMATION-008: {name} frame count contract sai")
        if {frame.get("slot_index") for frame in frames if isinstance(frame, dict)} != set(range(expected["slot_count"])):
            raise AnimationError(f"ASSET-ANIMATION-008: {name} contract slot sai")
        for frame in frames:
            client_cf = frame.get("client_cf") if isinstance(frame, dict) else None
            if not isinstance(client_cf, list) or not client_cf or len(client_cf) != len(set(client_cf)) or any(not isinstance(value, int) or isinstance(value, bool) or not -128 <= value <= 127 for value in client_cf):
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
