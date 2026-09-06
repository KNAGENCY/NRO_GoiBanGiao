"""Hợp đồng Bước 0: registry lấy đúng từ năm manifest canonical của Đợt 4."""

import json
import unittest
from pathlib import Path


PLAN = Path(__file__).resolve().parent
REGISTRY = PLAN.parent / "kich-thuoc-anh.json"


def _load(name):
    return json.loads((PLAN / name).read_text(encoding="utf-8"))


def canonical_registry_entries():
    boss = _load("04-boss.json")["meta"]["ten_moi_can_them_kich_thuoc"]
    maps = _load("05-map.json")["ten_moi_can_them_kich_thuoc"]
    kunai = _load("06-vat-pham.json")["meta"]["ten_moi_can_them_kich_thuoc"]
    ui = _load("08-giao-dien.json")["ten_moi_kich_thuoc"]
    brand = _load("09-chan-dung-thuong-hieu.json")["meta"]["ten_moi_can_them_kich_thuoc"]
    result = {}
    for item in boss:
        result[item["ten_file"]] = item["kich_thuoc"]
    result.update(maps)
    for item in kunai:
        result[item["ten_file"]] = item["kich_thuoc"]
    result.update(ui)
    for item in brand:
        name, width, height = item.split(" = ")[0], *map(int, item.split(" = ")[1].split())
        result[name] = [width, height]
    return result


class Buoc0RegistryTests(unittest.TestCase):
    def test_canonical_manifest_counts_and_names_are_unique(self):
        entries = canonical_registry_entries()
        boss_meta = _load("04-boss.json")["meta"]
        active_boss = boss_meta["ten_moi_can_them_kich_thuoc"]
        reserve_boss = boss_meta["ten_du_phong_neu_doi_D1"]["danh_sach"]
        self.assertEqual(len(active_boss) + len(reserve_boss), 10)
        self.assertEqual(boss_meta["quyet_dinh_registry_D1"]["chi_noi"], "ten_moi_can_them_kich_thuoc")
        self.assertFalse({item["ten_file"] for item in reserve_boss} & set(entries))
        self.assertEqual(len(_load("05-map.json")["ten_moi_can_them_kich_thuoc"]), 165)
        self.assertEqual(len(_load("06-vat-pham.json")["meta"]["ten_moi_can_them_kich_thuoc"]), 1)
        self.assertEqual(len(_load("08-giao-dien.json")["ten_moi_kich_thuoc"]), 30)
        self.assertEqual(len(_load("09-chan-dung-thuong-hieu.json")["meta"]["ten_moi_can_them_kich_thuoc"]), 17)
        self.assertEqual(len(entries), 216)

    def test_merged_manifest_is_regenerated_from_c1_player_manifest(self):
        merged = _load("manifest-dot-4.json")
        run = next(item for item in merged["anh"] if item["ten_file"] == "player_run_5f_chroma.png")
        hit = next(item for item in merged["anh"] if item["ten_file"] == "player_hit_death_getup_3f_chroma.png")
        self.assertIn("C1 đã chốt", run["ghi_chu"])
        self.assertIn("cùng cỡ đầu", hit["ghi_chu"])

    def test_every_canonical_name_is_appended_once_with_exact_size(self):
        entries = canonical_registry_entries()
        registry = json.loads(REGISTRY.read_text(encoding="utf-8"))

        self.assertEqual(list(registry)[-len(entries):], list(entries))
        self.assertEqual({name: registry.get(name) for name in entries}, entries)


if __name__ == "__main__":
    unittest.main()
