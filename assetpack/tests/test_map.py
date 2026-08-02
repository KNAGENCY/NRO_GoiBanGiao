import unittest


class MapContractTests(unittest.TestCase):
    def test_collision_payload_is_header_then_exactly_width_times_height_cells(self):
        from assetpack.map import MapError, encode_collision, semantic_hash

        payload = encode_collision(2, 2, [False, True, False, True], spawn_index=0, tile_size=(24, 24))
        self.assertEqual(payload, b"\x02\x02\x18\x18\x00\x01\x00\x01")
        self.assertEqual(semantic_hash([False, True, False, True]), semantic_hash([False, True, False, True]))
        with self.assertRaisesRegex(MapError, "ASSET-COLLISION-006"):
            encode_collision(2, 2, [False, True, False, 2], spawn_index=0, tile_size=(24, 24))

    def test_solid_tile_semantics_must_match_client_collision(self):
        from assetpack.map import MapError, validate_collision_semantics

        with self.assertRaisesRegex(MapError, "ASSET-COLLISION-006"):
            validate_collision_semantics([{"id": 7, "solid": True}], [7], [False])
