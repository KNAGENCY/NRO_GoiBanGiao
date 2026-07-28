import unittest


class MapContractTests(unittest.TestCase):
    def test_collision_payload_is_header_then_exactly_width_times_height_cells(self):
        from assetpack.map import MapError, encode_collision, semantic_hash

        payload = encode_collision(2, 2, [0, 1, 2, 3], spawn_index=0)
        self.assertEqual(payload, b"\x02\x02\x00\x01\x02\x03")
        self.assertEqual(semantic_hash([False, True, False, True]), semantic_hash([False, True, False, True]))
        with self.assertRaisesRegex(MapError, "ASSET-COLLISION-006"):
            encode_collision(2, 2, [0, 1, 2, 100], spawn_index=0)

    def test_solid_tile_semantics_must_match_client_collision(self):
        from assetpack.map import MapError, validate_collision_semantics

        with self.assertRaisesRegex(MapError, "ASSET-COLLISION-006"):
            validate_collision_semantics([{"id": 7, "solid": True}], [7], [False])
