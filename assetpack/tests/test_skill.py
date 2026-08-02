import unittest


class SkillContractTests(unittest.TestCase):
    def test_unity_skill_layout_is_big_endian_with_twelve_levels_and_minus_one_sentinels(self):
        from assetpack.skill import encode_skill, parse_skill

        record = {"id": 9, "pose_slot": 3, "levels": list(range(12)), "effect_id": -1, "arrow_id": -1}
        encoded = encode_skill([record])
        self.assertEqual(encoded[:4], b"\x00\x01\x00\x09")
        self.assertEqual(parse_skill(encoded), (record,))

    def test_effect_80_is_rejected_instead_of_becoming_a_hidden_client_load(self):
        from assetpack.skill import SkillError, encode_skill

        with self.assertRaisesRegex(SkillError, "ASSET-SKILL-004"):
            encode_skill([{"id": 9, "pose_slot": 3, "levels": list(range(12)), "effect_id": 80, "arrow_id": -1}])
