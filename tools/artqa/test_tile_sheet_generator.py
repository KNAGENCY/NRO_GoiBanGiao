import hashlib
import json
import shutil
import tempfile
import unittest
from pathlib import Path


ROOT = Path(__file__).parents[2]
SOURCE = Path(r"D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1\generated\naruto-village-tile-sheet-10x10-v2.png")
SOURCE_SHA256 = "2031edef163fc766b9f44df6c30724c0e41c57730b4f98749e3940fae157036b"


class TileSheetGeneratorTests(unittest.TestCase):
    def test_v2_produces_exactly_100_tiles_with_four_zoom_outputs_and_lattice_provenance(self):
        from tools.artqa.generate_tile_sheet import generate_tile_sheet

        with tempfile.TemporaryDirectory() as temporary:
            output = Path(temporary) / "tiles"
            index = generate_tile_sheet(SOURCE, output, expected_sha256=SOURCE_SHA256)

            self.assertEqual(index["source_sha256"], SOURCE_SHA256)
            self.assertEqual(
                {key: index["lattice"][key] for key in ("columns", "rows")},
                {"columns": 10, "rows": 10},
            )
            self.assertTrue(index["lattice"]["all_cells_nonempty"])
            self.assertGreater(index["lattice"]["nonblack_ratio_min"], 0.10)
            self.assertEqual(len(index["tiles"]), 100)
            for tile in index["tiles"]:
                self.assertEqual(set(tile["zooms"]), {"x1", "x2", "x3", "x4"})
                for zoom, record in tile["zooms"].items():
                    path = output / record["path"]
                    self.assertTrue(path.is_file(), f"tile {tile['index']} {zoom}")
                    self.assertEqual(record["sha256"], hashlib.sha256(path.read_bytes()).hexdigest())

            stored = json.loads((output / "generated-index.json").read_text(encoding="utf-8"))
            self.assertEqual(stored, index)

    def test_source_sha_mutation_is_rejected_before_any_tile_output(self):
        from tools.artqa.generate_tile_sheet import TileSheetError, generate_tile_sheet

        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            source = root / SOURCE.name
            shutil.copyfile(SOURCE, source)
            source.write_bytes(source.read_bytes() + b"mutation")
            output = root / "tiles"

            with self.assertRaisesRegex(TileSheetError, "ASSET-TILE-001"):
                generate_tile_sheet(source, output, expected_sha256=SOURCE_SHA256)
            self.assertFalse(output.exists())


if __name__ == "__main__":
    unittest.main()
