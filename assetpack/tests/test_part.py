"""Oracle độc lập cho định dạng part.bin của client Unity."""

import hashlib
import json
import struct
import unittest
from pathlib import Path

from assetpack.part import PartError, build_part, encode_part, parse_part


ROOT = Path(__file__).resolve().parents[2]
LEGACY_PART = ROOT / "server" / "internal" / "game" / "part.bin"
MANIFEST = ROOT / "docs" / "my-thuat" / "mapping-id.json"
FRAME_COUNTS = {0: 3, 1: 17, 2: 14}
LEGACY_SHA256 = "09349f28a95b46469b691abf87d6b341f3395f53ca8d6d704b48ab706992826b"


def scan_legacy_records(data: bytes):
    """Đọc byte thô riêng với parser; oracle này không gọi code production."""
    count = struct.unpack_from(">h", data, 0)[0]
    offset = 2
    records = []
    for _ in range(count):
        start = offset
        part_type = data[offset]
        offset += 1
        frame_count = FRAME_COUNTS[part_type]
        offset += frame_count * 4
        records.append(data[start:offset])
    if offset != len(data):
        raise AssertionError(f"oracle dừng tại {offset}, file dài {len(data)}")
    return records


def record_bytes(part_type: int, sprite_id: int, frame_count: int) -> bytes:
    """Golden encoder tối giản, độc lập với encode_part để làm expected fixture."""
    return bytes([part_type]) + b"".join(
        struct.pack(">hbb", sprite_id, 0, 0) for _ in range(frame_count)
    )


class PartBinaryTest(unittest.TestCase):
    def setUp(self):
        self.legacy = LEGACY_PART.read_bytes()
        self.manifest = json.loads(MANIFEST.read_text(encoding="utf-8"))

    def test_legacy_oracle_pins_records_types_frames_bytes_and_sha(self):
        self.assertEqual(87707, len(self.legacy))
        self.assertEqual(LEGACY_SHA256, hashlib.sha256(self.legacy).hexdigest())
        records = scan_legacy_records(self.legacy)
        self.assertEqual(2189, len(records))
        self.assertEqual("00001100000012000000140000", records[0].hex())
        self.assertEqual("01026000000245000002460000024700ff0248000002490000024a0000026100ff024b0000024c0000024d0000024e0100024f000002500000025101000252000000100000", records[57].hex())
        self.assertEqual({13, 57, 69}, {len(record) for record in records})

    def test_legacy_round_trip_is_byte_identical(self):
        self.assertEqual(self.legacy, encode_part(parse_part(self.legacy)))

    def test_build_replaces_only_phase_one_records_and_appends_iruka(self):
        actual = build_part(self.manifest, self.legacy)
        actual_records = scan_legacy_records(actual)
        legacy_records = scan_legacy_records(self.legacy)
        self.assertEqual(2192, len(actual_records))
        for record_id in range(2189):
            if record_id not in {30, 57, 58}:
                self.assertEqual(legacy_records[record_id], actual_records[record_id])

        expected = list(legacy_records)
        expected[30] = b"\x00" + b"".join(
            struct.pack(">hbb", sprite_id, 0, 0) for sprite_id in (27041, 27042, 27043)
        )
        expected[57] = record_bytes(1, 27044, 17)
        expected[58] = record_bytes(2, 27061, 14)
        expected.extend((
            record_bytes(0, 27075, 3),
            record_bytes(1, 27076, 17),
            record_bytes(2, 27077, 14),
        ))
        golden = struct.pack(">h", 2192) + b"".join(expected)
        self.assertEqual(golden, actual)
        self.assertEqual("935c7a32b7b8e40d7d66234611931b5803b348310fcca60e4de324a023c8d354", hashlib.sha256(actual).hexdigest())

    def test_rejects_corrupt_endian_type_and_frame_layout_with_part_code(self):
        corruptions = (
            b"\x8d\x05" + self.legacy[2:],
            b"\x00\x01\x03",
            b"\x00\x01\x00\x00\x01\x00\x00",
        )
        for data in corruptions:
            with self.subTest(data=data.hex()):
                with self.assertRaisesRegex(PartError, "ASSET-PART-001"):
                    parse_part(data)

    def test_rejects_committed_corruption_fixtures_with_part_code(self):
        fixtures = (
            "part-corrupt-endian.bin",
            "part-corrupt-type.bin",
            "part-corrupt-frame-count.bin",
        )
        for name in fixtures:
            with self.subTest(name=name):
                data = (ROOT / "assetpack" / "tests" / "fixtures" / name).read_bytes()
                with self.assertRaisesRegex(PartError, "ASSET-PART-001"):
                    parse_part(data)

    def test_rejects_manifest_frame_count_and_missing_sprite_id(self):
        bad_count = json.loads(json.dumps(self.manifest))
        bad_count["parts"][0]["frames"].pop()
        with self.assertRaisesRegex(PartError, "ASSET-PART-001"):
            build_part(bad_count, self.legacy)

        missing_sprite = json.loads(json.dumps(self.manifest))
        missing_sprite["parts"][0]["frames"][0]["sprite_id"] = 1
        with self.assertRaisesRegex(PartError, "ASSET-SPRITE-002"):
            build_part(missing_sprite, self.legacy)


if __name__ == "__main__":
    unittest.main()
