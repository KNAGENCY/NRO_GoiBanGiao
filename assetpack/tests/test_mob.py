import unittest


class MobContractTests(unittest.TestCase):
    def test_mob_layout_retains_idle_run_pounce_hit_and_death_frames(self):
        from assetpack.mob import encode_mob, parse_mob

        frames = {"idle": [0, 1], "run": [2, 3], "pounce": [4], "hit": [5], "death": [6]}
        self.assertEqual(parse_mob(encode_mob(7, frames)), (7, frames))
