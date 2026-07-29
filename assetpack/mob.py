"""Parser strict cho binary ``EffectData.readData`` của client Unity legacy."""

import struct


class MobError(ValueError):
    pass


_ACTIONS = ("idle", "run", "pounce", "hit", "death")
_FRAME_COUNTS = {"idle": 2, "run": 4, "pounce": 2, "hit": 1, "death": 1}
_MAX_IMAGES = 255
_MAX_FRAMES = 1024


class _Reader:
    def __init__(self, data: bytes):
        if not isinstance(data, (bytes, bytearray, memoryview)):
            raise MobError("ASSET-MOB-003: payload không phải binary")
        self.data = bytes(data)
        self.pos = 0

    def take(self, count: int) -> bytes:
        if count < 0 or self.pos + count > len(self.data):
            raise MobError("ASSET-MOB-003: payload thiếu byte theo length")
        value = self.data[self.pos:self.pos + count]
        self.pos += count
        return value

    def u8(self) -> int:
        return self.take(1)[0]

    def s8(self) -> int:
        value = self.u8()
        return value - 256 if value >= 128 else value

    def s16(self) -> int:
        return int.from_bytes(self.take(2), "big", signed=True)


def parse_mob(data: bytes) -> dict:
    """Đọc đúng các field mà ``EffectData.readData(myReader)`` đọc.

    Format không có JSON/header tự đặt: ``sbyte imageCount``, mỗi image gồm
    ``sbyte id + 4 unsigned byte``; tiếp theo là ``short frameCount``, các
    frame/component ``sbyte + short dx + short dy + sbyte imageId``; cuối cùng
    là ``short arrFrameCount`` và các ``short arrFrame``.
    """
    try:
        reader = _Reader(data)
        image_count = reader.s8()
        if not 1 <= image_count <= _MAX_IMAGES:
            raise MobError("ASSET-MOB-003: image header malformed")
        if image_count * 5 + 2 > len(reader.data) - reader.pos:
            raise MobError("ASSET-MOB-003: image length vượt payload")
        images = []
        image_ids = set()
        for _ in range(image_count):
            image_id = reader.s8()
            x, y, width, height = (reader.u8() for _ in range(4))
            if image_id < 0 or image_id in image_ids or width == 0 or height == 0:
                raise MobError("ASSET-MOB-003: image header malformed")
            image_ids.add(image_id)
            images.append({"id": image_id, "x": x, "y": y, "w": width, "h": height})

        frame_count = reader.s16()
        if not 1 <= frame_count <= _MAX_FRAMES:
            raise MobError("ASSET-MOB-003: frame length malformed")
        frames = []
        for _ in range(frame_count):
            component_count = reader.s8()
            if component_count < 1:
                raise MobError("ASSET-MOB-003: frame payload malformed")
            components = []
            for _ in range(component_count):
                dx, dy, image_id = reader.s16(), reader.s16(), reader.s8()
                if image_id not in image_ids:
                    raise MobError("ASSET-MOB-003: frame image reference malformed")
                components.append([dx, dy, image_id])
            frames.append(components)

        arr_count = reader.s16()
        if arr_count < 0:
            raise MobError("ASSET-MOB-003: arrFrame length malformed")
        arr_frames = [reader.s16() for _ in range(arr_count)]
        # Legacy EffectData.readData không đọc hai byte reserved cuối payload;
        # giữ chúng trong golden để không biến lỗi length thành im lặng.
        if len(reader.data) - reader.pos != 2:
            raise MobError("ASSET-MOB-003: payload length/trailer malformed")
        trailer = list(reader.take(2))
        return {"images": images, "frames": frames, "arr_frames": arr_frames, "trailer": trailer}
    except MobError:
        raise
    except (IndexError, OverflowError, TypeError, ValueError) as error:
        raise MobError(f"ASSET-MOB-003: payload EffectData malformed: {error}") from error


def _expected_frames(mob: dict) -> list[list[list[int]]]:
    frames = mob.get("frames")
    if not isinstance(frames, dict) or set(frames) != set(_ACTIONS):
        raise MobError("ASSET-MOB-003: golden thiếu idle/run/pounce/hit/death")
    flattened = []
    for action in _ACTIONS:
        value = frames[action]
        if not isinstance(value, list) or len(value) != _FRAME_COUNTS[action]:
            raise MobError("ASSET-MOB-003: golden sai action/frame count")
        for frame in value:
            if not isinstance(frame, list):
                raise MobError("ASSET-MOB-003: golden frame payload malformed")
            for component in frame:
                if (
                    not isinstance(component, list)
                    or len(component) != 3
                    or any(not isinstance(item, int) or isinstance(item, bool) for item in component)
                    or not -32768 <= component[0] <= 32767
                    or not -32768 <= component[1] <= 32767
                    or not -128 <= component[2] <= 127
                ):
                    raise MobError("ASSET-MOB-003: golden frame payload malformed")
        flattened.extend(value)
    return flattened


def encode_mob(mob: dict) -> bytes:
    """Sinh payload EffectData trực tiếp từ golden Sói Rừng trong manifest."""
    if not isinstance(mob, dict) or mob.get("payload_mode") != "generated":
        raise MobError("ASSET-MOB-003: production Mob phải khai payload_mode generated")
    images = mob.get("images")
    if not isinstance(images, list) or not 1 <= len(images) <= 127:
        raise MobError("ASSET-MOB-003: image golden malformed")
    output = bytearray(struct.pack("b", len(images)))
    image_ids = set()
    for image in images:
        if not isinstance(image, dict):
            raise MobError("ASSET-MOB-003: image golden malformed")
        values = [image.get(key) for key in ("id", "x", "y", "w", "h")]
        image_id, x, y, width, height = values
        if (
            any(not isinstance(value, int) or isinstance(value, bool) for value in values)
            or not 0 <= image_id <= 127
            or image_id in image_ids
            or not 0 <= x <= 255
            or not 0 <= y <= 255
            or not 1 <= width <= 255
            or not 1 <= height <= 255
        ):
            raise MobError("ASSET-MOB-003: image golden malformed")
        image_ids.add(image_id)
        output.extend(struct.pack(">bBBBB", image_id, x, y, width, height))

    frames = _expected_frames(mob)
    if not 1 <= len(frames) <= 32767:
        raise MobError("ASSET-MOB-003: frame golden malformed")
    output.extend(struct.pack(">h", len(frames)))
    for frame in frames:
        if not 1 <= len(frame) <= 127:
            raise MobError("ASSET-MOB-003: frame golden malformed")
        output.extend(struct.pack("b", len(frame)))
        for dx, dy, image_id in frame:
            if image_id not in image_ids:
                raise MobError("ASSET-MOB-003: frame image reference malformed")
            output.extend(struct.pack(">hhb", dx, dy, image_id))

    arr_frames = mob.get("arr_frames")
    trailer = mob.get("trailer")
    if (
        not isinstance(arr_frames, list)
        or len(arr_frames) > 32767
        or any(
            not isinstance(value, int)
            or isinstance(value, bool)
            or not -32768 <= value <= 32767
            for value in arr_frames
        )
        or not isinstance(trailer, list)
        or len(trailer) != 2
        or any(
            not isinstance(value, int)
            or isinstance(value, bool)
            or not 0 <= value <= 255
            for value in trailer
        )
    ):
        raise MobError("ASSET-MOB-003: arrFrame/trailer golden malformed")
    output.extend(struct.pack(">h", len(arr_frames)))
    for value in arr_frames:
        output.extend(struct.pack(">h", value))
    output.extend(bytes(trailer))
    payload = bytes(output)
    parsed = parse_mob(payload)
    expected = {
        "images": images,
        "frames": frames,
        "arr_frames": arr_frames,
        "trailer": trailer,
    }
    if parsed != expected:
        raise MobError("ASSET-MOB-003: encoder không round-trip golden")
    return payload


def validate_mob_golden(mob: dict, bytes_payload: bytes | None = None, data_payload: bytes | None = None) -> None:
    if not isinstance(mob, dict):
        raise MobError("ASSET-MOB-003: thiếu golden frame/action")
    if not isinstance(mob.get("id"), int) or isinstance(mob.get("id"), bool):
        raise MobError("ASSET-MOB-003: mob.id malformed")
    for field, payload in (("bytes_sha256", bytes_payload), ("data_sha256", data_payload)):
        expected_sha = mob.get(field)
        if not isinstance(expected_sha, str) or len(expected_sha) != 64:
            raise MobError(f"ASSET-MOB-003: thiếu SHA Mob {field}")
        if payload is None:
            raise MobError("ASSET-MOB-003: thiếu payload golden để đối chiếu")
        import hashlib
        if hashlib.sha256(payload).hexdigest() != expected_sha:
            raise MobError(f"ASSET-MOB-003: SHA Mob {field} không khớp payload")

    expected_images = mob.get("images")
    expected_arr = mob.get("arr_frames")
    expected_trailer = mob.get("trailer")
    expected_frames = _expected_frames(mob)
    if not isinstance(expected_images, list) or not isinstance(expected_arr, list) or not isinstance(expected_trailer, list):
        raise MobError("ASSET-MOB-003: thiếu golden header/arrFrame")
    try:
        parsed_bytes = parse_mob(bytes_payload)
        parsed_data = parse_mob(data_payload)
    except MobError:
        raise
    expected = {"images": expected_images, "frames": expected_frames, "arr_frames": expected_arr, "trailer": expected_trailer}
    if parsed_bytes != expected or parsed_data != expected or parsed_bytes != parsed_data:
        raise MobError("ASSET-MOB-003: payload EffectData không khớp golden contract")
