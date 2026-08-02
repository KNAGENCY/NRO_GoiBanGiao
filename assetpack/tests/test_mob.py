import hashlib
import unittest
from pathlib import Path


FIXTURE = Path(__file__).parent / "fixtures" / "legacy" / "mob.bin"
GOLDEN_IMAGES = [
    {"id": 0, "x": 0, "y": 0, "w": 12, "h": 16},
    {"id": 1, "x": 12, "y": 0, "w": 13, "h": 16},
    {"id": 2, "x": 25, "y": 0, "w": 12, "h": 16},
    {"id": 3, "x": 37, "y": 0, "w": 13, "h": 17},
    {"id": 4, "x": 50, "y": 0, "w": 17, "h": 20},
    {"id": 5, "x": 67, "y": 0, "w": 17, "h": 21},
    {"id": 6, "x": 116, "y": 0, "w": 19, "h": 22},
    {"id": 7, "x": 135, "y": 0, "w": 21, "h": 22},
    {"id": 8, "x": 84, "y": 0, "w": 15, "h": 21},
    {"id": 9, "x": 99, "y": 0, "w": 17, "h": 21},
]
GOLDEN_FRAMES = [
    [[3, -17, 0], [17, -20, 4], [37, -21, 8]],
    [[2, -17, 1], [17, -21, 5], [36, -21, 9]],
    [[3, -17, 2], [18, -17, 0], [33, -20, 4], [52, -22, 6]],
    [[2, -18, 3], [17, -17, 1], [33, -21, 5], [52, -22, 7]],
    [[2, -20, 4], [22, -22, 6], [43, -21, 8], [63, -17, 0], [81, -17, 2]],
    [[2, -21, 5], [22, -22, 7], [42, -21, 9], [62, -17, 1], [80, -18, 3]],
    [[4, -20, 4], [24, -24, 8]],
    [[4, -21, 5], [23, -24, 9]],
    [[5, -26, 6], [27, -17, 2]],
    [[5, -26, 7], [26, -18, 3]],
]
GOLDEN_TRAILER = [0, 0]


class MobContractTests(unittest.TestCase):
    def test_encoder_builds_deterministic_round_trip_from_naruto_manifest(self):
        from assetpack.mob import encode_mob, parse_mob

        images = [
            {"id": index, "x": index * 24, "y": 0, "w": 24, "h": 24}
            for index in range(10)
        ]
        flat_frames = [[[-12, -24, index]] for index in range(10)]
        mob = {
            "id": 0,
            "payload_mode": "generated",
            "images": images,
            "frames": {
                "idle": flat_frames[:2],
                "run": flat_frames[2:6],
                "pounce": flat_frames[6:8],
                "hit": flat_frames[8:9],
                "death": flat_frames[9:],
            },
            "arr_frames": [],
            "trailer": [0, 0],
        }

        first = encode_mob(mob)
        second = encode_mob(mob)

        self.assertEqual(first, second)
        self.assertEqual(
            parse_mob(first),
            {
                "images": images,
                "frames": flat_frames,
                "arr_frames": [],
                "trailer": [0, 0],
            },
        )
        self.assertNotEqual(
            hashlib.sha256(first).hexdigest(),
            hashlib.sha256(FIXTURE.read_bytes()).hexdigest(),
        )

    def test_real_effectdata_fixture_matches_legacy_read_data_contract(self):
        from assetpack.mob import parse_mob

        parsed = parse_mob(FIXTURE.read_bytes())
        self.assertEqual(parsed["images"], GOLDEN_IMAGES)
        self.assertEqual(parsed["frames"], GOLDEN_FRAMES)
        self.assertEqual(parsed["arr_frames"], [])
        self.assertEqual(parsed["trailer"], GOLDEN_TRAILER)

    def test_non_binary_json_payload_is_rejected_with_stable_asset_code(self):
        from assetpack.mob import MobError, parse_mob

        with self.assertRaisesRegex(MobError, "ASSET-MOB-003"):
            parse_mob(b'[0,{"idle":[0,1]}]')

    def test_payload_header_mutation_is_rejected_even_when_sha_changes(self):
        from assetpack.mob import MobError, validate_mob_golden

        mutated = bytearray(FIXTURE.read_bytes())
        mutated[1] ^= 1
        golden = {
            "id": 0,
            "bytes_sha256": hashlib.sha256(mutated).hexdigest(),
            "data_sha256": hashlib.sha256(FIXTURE.read_bytes()).hexdigest(),
            "images": GOLDEN_IMAGES,
            "frames": {"idle": GOLDEN_FRAMES[:2], "run": GOLDEN_FRAMES[2:6], "pounce": GOLDEN_FRAMES[6:8], "hit": GOLDEN_FRAMES[8:9], "death": GOLDEN_FRAMES[9:]},
            "arr_frames": [],
            "trailer": GOLDEN_TRAILER,
        }
        with self.assertRaisesRegex(MobError, "ASSET-MOB-003"):
            validate_mob_golden(golden, bytes(mutated), FIXTURE.read_bytes())
