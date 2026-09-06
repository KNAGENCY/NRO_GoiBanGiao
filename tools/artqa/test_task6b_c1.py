"""Chốt C1 Đợt 4: dải chạy 5 khung và mọi đầu cùng cỡ idle."""

import unittest

from tools.artqa.task6b_closure import ClosureError, PLAYER_ACTIONS, PLAYER_BODY_SLOTS_V1, PLAYER_LEG_SLOTS_V1, derive_player_slot_specs


def _record(name, frames):
    return {
        "asset": name,
        "source_sha256": "test",
        "checks": {"crop": {
            "source_edges": [index * 100 for index in range(frames + 1)],
            "content_boxes": [[0, 0, 100, 100] for _ in range(frames)],
        }},
    }


def _moc(records, bad_run_frame=None):
    result = {}
    for _, name, frames in PLAYER_ACTIONS:
        for index in range(frames):
            size = 130 if name == "player_run_5f_chroma.png" and index == bad_run_frame else 100
            result[(name, index)] = {
                "edge_x": index * 100,
                "face_w": size,
                "face_h": size,
                "eye_h": size / 10,
                "eye_to_chin": size / 2,
                "head_h_hair_chin": size * 2,
                "band_thick": size / 4,
                "chin_y": index * 100 + 70,
                "face_top": index * 100 + 10,
                "face_cx": 50,
                "content_top": index * 100,
                "content_bottom": index * 100 + 100,
            }
    return result


class C1PlayerPipelineTests(unittest.TestCase):
    def test_c1_uses_five_run_frames_and_hit_slots(self):
        self.assertIn(("run", "player_run_5f_chroma.png", 5), PLAYER_ACTIONS)
        self.assertEqual(
            [PLAYER_BODY_SLOTS_V1[index] for index in range(2, 7)],
            [("run", index) for index in range(5)],
        )
        self.assertEqual(
            [PLAYER_LEG_SLOTS_V1[index] for index in range(2, 7)],
            [("run", index) for index in range(5)],
        )
        self.assertEqual(
            [PLAYER_BODY_SLOTS_V1[index] for index in range(14, 17)],
            [("hit_death_getup", index) for index in range(3)],
        )

    def test_c1_rejects_a_run_frame_with_different_head_size(self):
        records = {name: _record(name, frames) for _, name, frames in PLAYER_ACTIONS}
        review = {"assets": list(records.values())}

        with self.assertRaisesRegex(ClosureError, "cùng cỡ đầu"):
            derive_player_slot_specs(review, records=records, moc=_moc(records, bad_run_frame=2))


if __name__ == "__main__":
    unittest.main()
