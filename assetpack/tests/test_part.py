"""Oracle độc lập cho định dạng part.bin của client Unity."""

import hashlib
import json
import re
import struct
import unittest
from pathlib import Path

from assetpack.part import PartError, build_part, encode_part, parse_part


ROOT = Path(__file__).resolve().parents[2]
LEGACY_PART = ROOT / "server" / "internal" / "game" / "part.bin"
MANIFEST = ROOT / "docs" / "my-thuat" / "mapping-id.json"
FRAME_COUNTS = {0: 3, 1: 17, 2: 14}

# Do doi cua tung KHUNG ba manh nhan vat, don vi game (x1). Ghim CUNG o day
# (khong doc tu manifest) de oracle van doc lap: sua nham manifest la test nay do.
#
# Truoc 05/09/2026 o day chi co ba hang so (-41/-50/-33) vi moi khung bi nhoi vao
# cung mot hop 128x128/128x160, nen moi khung chung mot do doi. Sang (05/09): moi
# khung giu ti le CHUNG va duoc neo theo ban chan cua chinh no. Chieu 05/09 (V1 khe
# dau-than): dx/dy la NGHIEM MILP cho ca 34 manh (tools/artqa/nhan_vat_toi_uu.py) —
# than cat tu cam xuong (het "hai khuon mat"), ti le tung khung theo co dau, slot
# than/chan 2..6 = 5 khung chay, slot 15 het khung nam; dy dau slot 0 ghim -38.
# Bang nay chep tay tu bang `sinh-lai-sprite.py` in ra (khong doc tu manifest) de
# oracle van doc lap. Doi mot cap so o day nghia la nhan vat dich cho khi ve - phai co y.
DX_DY = {
    # dau (part 30) - ba khung idle
    27041: (-4, -38), 27042: (-4, -32), 27043: (-7, -36),
    # than (part 57): idle0 idle1 run0 run1 run2 run3 run0 melee0..3 cast0 cast1 cast2 cast2 cast0 idle2
    27044: (-5, -26), 27045: (-6, -33), 27046: (-9, -26), 27047: (-8, -24),
    27048: (-12, -26), 27049: (-8, -26), 27050: (-9, -26), 27051: (-2, -24),
    27052: (0, -19), 27053: (-10, -18), 27054: (-19, -24), 27055: (-5, -24),
    27056: (-6, -18), 27057: (-3, -25), 27058: (-8, -30), 27059: (-8, -29),
    27060: (-3, -24),
    # chan (part 58): idle0 idle1 run0 run1 run2 run3 run0 melee0..3 cast0 cast1 idle2
    27061: (7, -9), 27062: (-7, -22), 27063: (-9, -19), 27064: (-7, -17),
    27065: (-15, -18), 27066: (-7, -18), 27067: (-9, -20), 27068: (-8, -14),
    27069: (-16, -16), 27070: (0, -22), 27071: (-13, -6), 27072: (-10, -6),
    27073: (-9, -9), 27074: (-9, -15),
    # Iruka (part 2189 dau / 2190 than / 2191 chan) — 06/09/2026: truoc do dx=dy=0 (dau lung
    # +13, than chim +29 px). Nay component_placement tai cf 0 voi ti le s0/1.6642; moi frame
    # cua mot part cung tro mot sprite nen cung dx/dy. Chep tay tu bang IRUKA sinh-lai-sprite.py.
    27075: (-1, -39), 27076: (-3, -39), 27077: (-2, -22),
}
LEGACY_SHA256 = "09349f28a95b46469b691abf87d6b341f3395f53ca8d6d704b48ab706992826b"
# sha256 part.bin sau V1 05/09 + Iruka 06/09 (in boi sinh-lai-sprite.py: "part.bin: ... sha ...").
BUILT_SHA256 = "27075681b1a5426884268dddccf32c500e8a25ece854b36e69427abd99605ea4"


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

        def khung(part_type: int, ids) -> bytes:
            return bytes([part_type]) + b"".join(
                struct.pack(">hbb", sprite_id, *DX_DY[sprite_id])
                for sprite_id in ids
            )

        expected = list(legacy_records)
        expected[30] = khung(0, (27041, 27042, 27043))
        expected[57] = khung(1, range(27044, 27061))
        expected[58] = khung(2, range(27061, 27075))
        # Iruka: 3/17/14 frame cung tro mot sprite, cung dx/dy (khong con 0/0 nhu truoc 06/09)
        expected.extend((
            khung(0, [27075] * 3),
            khung(1, [27076] * 17),
            khung(2, [27077] * 14),
        ))
        self.assertNotEqual(record_bytes(0, 27075, 3), expected[2189], "Iruka dx=dy=0 la ban loi cu")
        golden = struct.pack(">h", 2192) + b"".join(expected)
        self.assertEqual(golden, actual)
        self.assertEqual(BUILT_SHA256, hashlib.sha256(actual).hexdigest())

    def test_swapping_two_body_slots_breaks_the_independent_part_golden(self):
        mutated = json.loads(json.dumps(self.manifest))
        body = next(part for part in mutated["parts"] if part["id"] == 57)
        body["frames"][0], body["frames"][1] = body["frames"][1], body["frames"][0]

        changed = build_part(mutated, self.legacy)

        self.assertNotEqual(BUILT_SHA256, hashlib.sha256(changed).hexdigest())

    def test_committed_part_bin_khop_golden_va_lock(self):
        """File part.bin ĐÃ COMMIT (server/internal/game/assets) phải đúng golden và đúng lock.

        Trước 06/09 oracle chỉ kiểm build_part(manifest) — đổi 1 byte trong file đã commit
        thì test_part vẫn xanh (đột biến M10 của kiểm đối kháng). Đây là cửa cho chính file
        máy chủ nhúng (go:embed dongbo.go). Cũng ghim lock.versions.data == vsData trong dongbo.go.
        """
        assets = ROOT / "server" / "internal" / "game" / "assets"
        committed = (assets / "part.bin").read_bytes()
        self.assertEqual(BUILT_SHA256, hashlib.sha256(committed).hexdigest())
        lock = json.loads((assets / "asset-lock.json").read_text(encoding="utf-8"))
        out = [o for o in lock["outputs"] if o["path"] == "server/part.bin"]
        self.assertEqual(1, len(out))
        self.assertEqual(BUILT_SHA256, out[0]["sha256"])
        self.assertEqual(len(committed), out[0]["bytes"])
        # versions.data của lock phải đúng vsData máy chủ nhúng (dongbo.go); mapping["versions"]
        # là khối của closure v11 đã ký nên KHÔNG so ở đây (test_deploy ghim nó riêng).
        dongbo = (ROOT / "server" / "internal" / "game" / "dongbo.go").read_text(encoding="utf-8")
        vs_data = int(re.search(r"vsData\s+byte\s*=\s*(\d+)", dongbo).group(1))
        self.assertEqual(vs_data, lock["versions"]["data"])

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
