"""Ghép strip action theo slot index, không theo thứ tự input dễ bị đảo cell."""

from PIL import Image


class AnimationError(ValueError):
    pass


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
