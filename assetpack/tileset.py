"""Bảng vật lý tile mà máy chủ phát cho client (`tile_set_info.bin`, gói `-82`).

TẠI SAO có file này thay vì cứ chép thẳng bảng gốc:

Bảng gốc mua kèm server Ngọc Rồng chỉ khai báo tile **1–39**. Map Làng Lá dùng
tile **60–69 và 83–88**, mà client đọc bảng này để dựng `types[]` — tile không
được khai thì `types[] = 0`, tức "không có gì cả": không sàn, không vách. Cả màn
không có sàn thì **nhân vật rơi xuyên đất xuống đáy màn** (lỗi thật đã phải vá
ngày 22/08/2026).

Trước đây phần "thêm tile Làng Lá" được làm **bằng tay** bằng
`tools/map-lang-la/them-tile-vao-tileset.py` rồi commit thẳng file .bin vào
`server/internal/game/assets/`. Hậu quả: `assetpack build` không hề biết chuyện
đó, cứ chép lại bảng gốc 1696 byte — ai build rồi chép artifact sang máy chủ là
**lùi bảng tile về bản gốc và tái sinh nguyên lỗi cũ, im lặng**.

Nên phép biến đổi được đưa hẳn vào đường ống: nguồn vẫn là MỘT (bảng gốc đã qua
review + khai báo `build_inputs.tile_set_overlay` trong manifest), và kết quả bị
ghim bằng golden nên trôi một byte là build đỏ ngay chứ không đợi vào game mới biết.

Định dạng file (đọc ngược từ client `Controller.cs`, nhánh `case -82`):

    byte  soBoTile
      byte  soNhom
        int32 loai      (big-endian; là MẶT NẠ BIT của thuộc tính vật lý)
        byte  soTile
          byte tile ...

`loai` là mặt nạ bit chứ không phải số thứ tự, nên hai nhóm khác `loai` mô tả hai
thuộc tính khác nhau của cùng một tile — một tile hoàn toàn có thể có mặt ở nhiều
nhóm. Bảng gốc còn có tile **lặp trong cùng một nhóm** (nhóm 0x0008 có 29 hai
lần), nên tuyệt đối KHÔNG được "dọn dẹp" khi đọc/ghi lại: giữ nguyên từng byte
mới so sánh được với bản đang chạy.
"""

import hashlib
import struct
from typing import Any


class TilesetError(ValueError):
    pass


# Bit "đứng được lên trên" trong `types[]` của client. Đây là bit duy nhất mà
# thiếu nó thì nhân vật rơi xuyên đất — mọi chốt kiểm ngữ nghĩa đều xoay quanh nó.
FLOOR_BIT = 0x0002


def decode_tile_set_info(data: bytes) -> list[list[tuple[int, list[int]]]]:
    """Bung file thành `[bộ][ (loại, [tile...]) ]`, GIỮ NGUYÊN thứ tự và bản lặp."""
    if not isinstance(data, (bytes, bytearray)) or not data:
        raise TilesetError("ASSET-TILESET-012: tile_set_info rỗng/malformed")
    data = bytes(data)
    offset = 0

    def _byte() -> int:
        nonlocal offset
        if offset >= len(data):
            raise TilesetError("ASSET-TILESET-012: tile_set_info cụt giữa chừng")
        value = data[offset]
        offset += 1
        return value

    tile_sets: list[list[tuple[int, list[int]]]] = []
    for _ in range(_byte()):
        groups: list[tuple[int, list[int]]] = []
        for _ in range(_byte()):
            if offset + 4 > len(data):
                raise TilesetError("ASSET-TILESET-012: tile_set_info cụt ở trường loại")
            kind = struct.unpack_from(">i", data, offset)[0]
            offset += 4
            count = _byte()
            if offset + count > len(data):
                raise TilesetError("ASSET-TILESET-012: tile_set_info cụt ở danh sách tile")
            groups.append((kind, list(data[offset:offset + count])))
            offset += count
        tile_sets.append(groups)
    if offset != len(data):
        raise TilesetError(f"ASSET-TILESET-012: thừa {len(data) - offset} byte cuối tile_set_info")
    return tile_sets


def encode_tile_set_info(tile_sets) -> bytes:
    """Đóng gói ngược lại. Phải khép kín với `decode_tile_set_info` từng byte."""
    if not isinstance(tile_sets, list) or not 0 <= len(tile_sets) <= 255:
        raise TilesetError("ASSET-TILESET-012: số bộ tile phải là 0..255")
    out = bytearray([len(tile_sets)])
    for groups in tile_sets:
        if not isinstance(groups, list) or not 0 <= len(groups) <= 255:
            raise TilesetError("ASSET-TILESET-012: số nhóm trong bộ tile phải là 0..255")
        out.append(len(groups))
        for kind, tiles in groups:
            if not isinstance(kind, int) or isinstance(kind, bool) or not -2147483648 <= kind <= 2147483647:
                raise TilesetError("ASSET-TILESET-012: loại nhóm phải là int32")
            tiles = list(tiles)
            if len(tiles) > 255:
                raise TilesetError(f"ASSET-TILESET-012: nhóm 0x{kind:04X} có {len(tiles)} tile, vượt 255")
            if any(not isinstance(tile, int) or isinstance(tile, bool) or not 0 <= tile <= 255 for tile in tiles):
                raise TilesetError("ASSET-TILESET-012: mã tile phải là byte 0..255")
            out += struct.pack(">i", kind)
            out.append(len(tiles))
            out += bytes(tiles)
    return bytes(out)


def _kiem_khai_bao(overlay: Any) -> tuple[int, list[dict[str, Any]]]:
    """Soi khai báo overlay trước khi động vào byte — sai chính tả khoá phải nổ, không im lặng."""
    if not isinstance(overlay, dict):
        raise TilesetError("ASSET-TILESET-012: thiếu khai báo tile_set_overlay")
    if set(overlay) - {"tile_set_index", "groups", "result", "note"}:
        raise TilesetError("ASSET-TILESET-012: tile_set_overlay có khoá lạ")
    if not {"tile_set_index", "groups"} <= set(overlay):
        raise TilesetError("ASSET-TILESET-012: tile_set_overlay thiếu tile_set_index/groups")
    index = overlay["tile_set_index"]
    if not isinstance(index, int) or isinstance(index, bool) or index < 0:
        raise TilesetError("ASSET-TILESET-012: tile_set_index phải là số nguyên không âm")
    groups = overlay["groups"]
    if not isinstance(groups, list) or not groups:
        raise TilesetError("ASSET-TILESET-012: tile_set_overlay.groups rỗng")
    seen_bits: set[int] = set()
    for group in groups:
        if not isinstance(group, dict) or not {"bit", "role", "tiles"} <= set(group) or set(group) - {"bit", "role", "tiles", "note"}:
            raise TilesetError("ASSET-TILESET-012: nhóm overlay thiếu bit/role/tiles hoặc có khoá lạ")
        bit = group["bit"]
        if not isinstance(bit, int) or isinstance(bit, bool) or not 1 <= bit <= 2147483647:
            raise TilesetError("ASSET-TILESET-012: bit nhóm overlay phải là int32 dương")
        if bit in seen_bits:
            raise TilesetError(f"ASSET-TILESET-012: nhóm overlay 0x{bit:04X} khai hai lần")
        seen_bits.add(bit)
        if not isinstance(group["role"], str) or not group["role"]:
            raise TilesetError("ASSET-TILESET-012: nhóm overlay thiếu role dạng chữ")
        tiles = group["tiles"]
        if (
            not isinstance(tiles, list)
            or not tiles
            or any(not isinstance(tile, int) or isinstance(tile, bool) or not 0 <= tile <= 255 for tile in tiles)
        ):
            raise TilesetError(f"ASSET-TILESET-012: nhóm overlay 0x{bit:04X} có danh sách tile sai")
        if len(set(tiles)) != len(tiles):
            raise TilesetError(f"ASSET-TILESET-012: nhóm overlay 0x{bit:04X} khai trùng tile")
    return index, groups


def apply_tile_set_overlay(base: bytes, overlay: Any, *, check_result: bool = True) -> bytes:
    """Thêm các tile Làng Lá vào bảng gốc theo khai báo trong manifest.

    Chỉ THÊM, không bao giờ sửa hay xoá: bảng gốc là di sản Ngọc Rồng chưa ai
    review lại, đụng vào là đổi vật lý của 39 bộ tile khác mà không ai biết.

    `check_result=False` chỉ dành cho test đột biến (muốn xem cái bảng sai nó ra
    thế nào); đường build thật luôn để mặc định True để trôi một byte là dừng.
    """
    index, groups = _kiem_khai_bao(overlay)
    tile_sets = decode_tile_set_info(base)
    if index >= len(tile_sets):
        raise TilesetError(f"ASSET-TILESET-012: bảng gốc không có bộ tile {index}")
    target = list(tile_sets[index])

    for group in groups:
        bit, tiles = group["bit"], list(group["tiles"])
        vi_tri = next((i for i, (kind, _) in enumerate(target) if kind == bit), None)
        if vi_tri is None:
            # Nhóm chưa tồn tại trong bảng gốc → tạo mới, xếp cuối bộ.
            target.append((bit, tiles))
            continue
        kind, co_san = target[vi_tri]
        trung = [tile for tile in tiles if tile in co_san]
        if trung:
            # Bảng gốc đã khai tile này rồi. Hoặc bảng gốc vừa đổi, hoặc người
            # viết overlay hiểu nhầm — cả hai đều phải dừng để có người nhìn lại,
            # chứ không âm thầm bỏ qua như script sửa tay ngày xưa.
            raise TilesetError(
                f"ASSET-TILESET-012: nhóm 0x{bit:04X} khai lại tile đã có sẵn trong bảng gốc: {trung}"
            )
        target[vi_tri] = (kind, co_san + tiles)

    tile_sets[index] = target
    data = encode_tile_set_info(tile_sets)

    if check_result:
        result = overlay.get("result")
        if (
            not isinstance(result, dict)
            or set(result) != {"byte_length", "sha256"}
            or result["byte_length"] != len(data)
            or not isinstance(result["sha256"], str)
            or result["sha256"] != hashlib.sha256(data).hexdigest()
        ):
            raise TilesetError(
                "ASSET-TILESET-013: bảng tile sau overlay khác golden đã ghim "
                f"(ra {len(data)} byte / {hashlib.sha256(data).hexdigest()})"
            )
    return data


def validate_floor_closure(data: bytes, *, tile_set_index: int, base_tile_count: int, solid_used) -> None:
    """Chốt ngữ nghĩa: tile MỚI nào máy chủ coi là đặc thì client phải khai là sàn.

    TẠI SAO cần kiểm ở tầng ý nghĩa chứ không chỉ so hash: hash chỉ bắt được
    "khác bản đã ghim", còn cái này bắt được "khai báo tự nó vô lý". Máy chủ chặn
    di chuyển theo `tile_art[].solid`; client cho đứng theo bit sàn. Tile máy chủ
    bảo "đặc" mà client không khai sàn = hai bên nói ngược nhau = người chơi rơi
    xuyên đất rồi bị máy chủ kéo ngược lại.

    Chỉ soi tile có mã **lớn hơn** số tile bảng gốc khai báo: 1..39 là di sản
    Ngọc Rồng, ngữ nghĩa của chúng không phải việc của Làng Lá và không được
    dùng làm cớ để chặn build.
    """
    tile_sets = decode_tile_set_info(data)
    if not 0 <= tile_set_index < len(tile_sets):
        raise TilesetError(f"ASSET-TILESET-012: không có bộ tile {tile_set_index}")
    floor = {tile for kind, tiles in tile_sets[tile_set_index] if kind == FLOOR_BIT for tile in tiles}
    thieu = sorted(tile for tile in solid_used if tile > base_tile_count and tile not in floor)
    if thieu:
        raise TilesetError(
            "ASSET-TILESET-013: tile máy chủ coi là đặc nhưng client không khai là sàn "
            f"(nhân vật sẽ rơi xuyên đất): {thieu}"
        )
