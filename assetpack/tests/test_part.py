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

# Do doi DOC cua ba manh nhan vat, don vi game. Ghim CUNG o day (khong doc tu
# manifest) de oracle van doc lap: sua nham manifest la test nay do.
#
# Vi sao khong phai 0: anh than chi ve NUA DUOI cua dau, nen ba manh phai chong
# len nhau dung cho hai ban tay trung nhau. Ba so nay do trong game that (dot lam
# map 1 Lang La) - lech mot o la nhan vat cut co hoac lun toi that lung.
DY_DAU, DY_THAN, DY_CHAN = -41, -50, -33
LEGACY_SHA256 = "09349f28a95b46469b691abf87d6b341f3395f53ca8d6d704b48ab706992826b"
BUILT_SHA256 = "4fba8a7038b0687bac5d49bbd5499a43b0d3955d693f30ea573672d140dfd618"


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
            struct.pack(">hbb", sprite_id, 0, DY_DAU)
            for sprite_id in (27041, 27042, 27043)
        )
        expected[57] = b"\x01" + b"".join(
            struct.pack(">hbb", sprite_id, 0, DY_THAN)
            for sprite_id in (
                27044, 27045, 27046, 27047, 27048, 27049, 27050, 27051, 27052,
                27053, 27054, 27055, 27056, 27057, 27058, 27059, 27060,
            )
        )
        expected[58] = b"\x02" + b"".join(
            struct.pack(">hbb", sprite_id, 0, DY_CHAN)
            for sprite_id in (
                27061, 27062, 27063, 27064, 27065, 27066, 27067,
                27068, 27069, 27070, 27071, 27072, 27073, 27074,
            )
        )
        expected.extend((
            record_bytes(0, 27075, 3),
            record_bytes(1, 27076, 17),
            record_bytes(2, 27077, 14),
        ))
        golden = struct.pack(">h", 2192) + b"".join(expected)
        self.assertEqual(golden, actual)
        self.assertEqual(BUILT_SHA256, hashlib.sha256(actual).hexdigest())

    def test_swapping_two_body_slots_breaks_the_independent_part_golden(self):
        mutated = json.loads(json.dumps(self.manifest))
        body = next(part for part in mutated["parts"] if part["id"] == 57)
        body["frames"][0], body["frames"][1] = body["frames"][1], body["frames"][0]

        changed = build_part(mutated, self.legacy)

        self.assertNotEqual(BUILT_SHA256, hashlib.sha256(changed).hexdigest())

    def test_rejects_manifest_missing_required_iruka_part_2191(self):
        missing_2191 = json.loads(json.dumps(self.manifest))
        missing_2191["parts"] = [
            part for part in missing_2191["parts"] if part["id"] != 2191
        ]

        with self.assertRaisesRegex(PartError, "ASSET-PART-001"):
            build_part(missing_2191, self.legacy)

    def test_rejects_manifest_with_extra_iruka_part_2192(self):
        extra_2192 = json.loads(json.dumps(self.manifest))
        extra_part = json.loads(json.dumps(
            next(part for part in extra_2192["parts"] if part["id"] == 2191)
        ))
        extra_part["id"] = 2192
        extra_2192["parts"].append(extra_part)

        with self.assertRaisesRegex(PartError, "ASSET-PART-001"):
            build_part(extra_2192, self.legacy)

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
