"""Ghép strip action theo slot index, không theo thứ tự input dễ bị đảo cell."""

import json
from pathlib import Path

from PIL import Image


class AnimationError(ValueError):
    pass


# Bảng khung xương 33 cf trích từ `Char.cs:483` của client (nguồn sự thật duy nhất về
# "cf nào tra slot nào"). Nằm ở tools/nhan-vat vì công cụ đo `ghep-khung.py` dùng chung.
CHARINFO_PATH = Path(__file__).resolve().parents[1] / "tools" / "nhan-vat" / "charinfo.json"
# CharInfo[cf][k]: k = 0 đầu, 1 chân, 2 thân (đúng thứ tự trong Char.cs).
_CHARINFO_SLOT = {"hoc-vien-head": 0, "hoc-vien-leg": 1, "hoc-vien-body": 2}
_SLOT_COUNT = {"hoc-vien-head": 3, "hoc-vien-body": 17, "hoc-vien-leg": 14}


def charinfo_client_cf(charinfo):
    """slot → tuple cf THẬT dùng slot đó, suy thẳng từ CharInfo cho cả ba mảnh.

    Vì sao (05/09/2026): bảng `REQUIRED_CLIENT_CF` cũ viết tay và SẮP THEO cf đầu tiên
    thay vì theo slot — thân slot 0 thật sự phục vụ cf 23/24 (bảng cũ nói 0/1), slot 1
    mới là idle 0/1; chân slot 0 = cf 8 (bảng cũ nói 0/1/22) và cf 26..32 (chân slot 8)
    bị bỏ quên hẳn. Ai bù dx/dy theo bảng cũ là lệch một khung — đó là một trong ba tầng
    lệch gây "đầu tách thân". Sinh từ dữ liệu để không bao giờ trôi nữa.
    """
    if not isinstance(charinfo, list) or len(charinfo) != 33:
        raise AnimationError("ASSET-ANIMATION-008: charinfo phải có đúng 33 cf")
    result = {}
    for name, k in _CHARINFO_SLOT.items():
        by_slot = {}
        for cf, entry in enumerate(charinfo):
            by_slot.setdefault(int(entry[k][0]), []).append(cf)
        count = _SLOT_COUNT[name]
        if sorted(by_slot) != list(range(count)):
            raise AnimationError(f"ASSET-ANIMATION-008: {name} slot trong CharInfo phải phủ đủ 0..{count - 1}")
        result[name] = tuple(tuple(by_slot[slot]) for slot in range(count))
    return result


def _nap_charinfo(path=CHARINFO_PATH):
    try:
        return json.loads(Path(path).read_text(encoding="utf-8"))
    except (OSError, ValueError) as error:
        raise AnimationError(f"ASSET-ANIMATION-008: không đọc được charinfo {path}: {error}") from error


# `cell_size` là Ô TỐI ĐA của contact sheet bằng chứng, KHÔNG phải kích thước bắt
# buộc của từng cell.
#
# Vì sao đổi (05/09/2026): sprite người chơi từng bị nhồi vào hộp cố định
# 128×128/128×160 nên mỗi khung có một tỉ lệ riêng — nhân vật phình/teo theo pose,
# đầu rời khỏi thân. Bản sáng 05/09 giữ một tỉ lệ chung; V1 chiều 05/09 mỗi khung một
# tỉ lệ theo cỡ đầu (s0 / hệ số đầu, tools/artqa/nhan_vat_moc.py) — cell vì thế có
# kích thước khác nhau. Client `SmallImage` vẽ ảnh
# mọi kích cỡ, chỉ contact sheet mới cần ô. Ô phải bao được cell rộng nhất; cell
# vượt ô vẫn là lỗi (fail-closed) chứ không bị cắt im lặng.
REQUIRED_ANIMATION_SHEETS = {
    "hoc-vien-head": {"slot_count": 3, "cell_size": (128, 128)},
    "hoc-vien-body": {"slot_count": 17, "cell_size": (256, 160)},
    "hoc-vien-leg": {"slot_count": 14, "cell_size": (160, 128)},
    "quai-soi-rung": {"slot_count": 10, "cell_size": (96, 96)},
}

# Ba bảng người chơi SINH từ CharInfo (xem `charinfo_client_cf`); sói là 1 slot / 1 cf.
REQUIRED_CLIENT_CF = {
    **charinfo_client_cf(_nap_charinfo()),
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
    """Dán cell vào ô, căn giữa-dưới. Cell nhỏ hơn ô là bình thường; vượt ô là lỗi."""
    if slot_count not in {3, 10, 14, 17}:
        raise AnimationError("ASSET-ANIMATION-008: slot count phải là 3/17/14/10")
    if len(cells) != slot_count or {slot for slot, _ in cells} != set(range(slot_count)):
        raise AnimationError("ASSET-ANIMATION-008: slot index phải đủ và duy nhất")
    width, height = cell_size
    sheet = Image.new("RGBA", (width * slot_count, height))
    for slot, cell in cells:
        if cell.width > width or cell.height > height:
            raise AnimationError("ASSET-ANIMATION-008: canvas cell không đúng")
        sheet.alpha_composite(
            cell.convert("RGBA"),
            (slot * width + (width - cell.width) // 2, height - cell.height),
        )
    return sheet
