"""Chốt kiểm: khối `build_inputs.map` trong manifest phải khớp map 1 THẬT.

TẠI SAO có file này — quả mìn im lặng đã nổ một lần:

`assetpack build` sinh `server/map-1-collision.bin` **từ manifest**, không phải từ
bố cục client. Tháng 8/2026 manifest còn ghi map 5×17 (di sản vertical slice) trong
khi map 1 Làng Lá thật là 52×24, nên lưới va chạm chỉ phủ góc trái bản đồ: người
chơi đứng ngoài lưới bị máy chủ **âm thầm từ chối mọi lệnh di chuyển** — không lỗi,
không log, chơi một mình không thấy. Lưới đúng 52×24 được sinh tay bằng
`tools/map-lang-la/sinh-va-cham.py` rồi commit thẳng vào `server/internal/game/assets/`,
NHƯNG manifest thì không ai sửa → **lần chạy `assetpack build` kế tiếp sẽ đè lại bản
sai 5×17** và tái sinh nguyên lỗi cũ.

Vì vậy chốt kiểm ở đây so ba nguồn phải cùng nói một chuyện:
  1. `docs/my-thuat/mapping-id.json`            — thứ `assetpack build` đọc
  2. `server/internal/game/assets/map-1-collision.bin` — artifact máy chủ đang chạy
  3. `docs/my-thuat/lang-la-client-assets/res/mymap/1.bytes` — bản sao bố cục client

Lệch một cái là đỏ ngay, thay vì đợi người chơi báo "đứng yên không đi được".
"""

import json
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
MAPPING = ROOT / "docs/my-thuat/mapping-id.json"
COLLISION = ROOT / "server/internal/game/assets/map-1-collision.bin"
# Bản sao byte thật của bố cục client, có trong git (client Unity `D:\NRO_Project`
# KHÔNG thuộc repo này nên không được phép làm nguồn cho test).
CLIENT_MAP = ROOT / "docs/my-thuat/lang-la-client-assets/res/mymap/1.bytes"


def _map_config() -> dict:
    return json.loads(MAPPING.read_text(encoding="utf-8"))["build_inputs"]["map"]


class Map1ManifestKhopArtifactTests(unittest.TestCase):
    def test_kich_thuoc_manifest_khop_header_client_va_header_va_cham(self):
        """Ba nơi cùng khai 52×24 — đây chính là con số đã lệch và gây lỗi."""
        client = CLIENT_MAP.read_bytes()
        collision = COLLISION.read_bytes()
        config = _map_config()

        # `1.bytes`: byte0 = rộng, byte1 = cao, sau đó W*H ô chỉ số tile.
        self.assertEqual((client[0], client[1]), (52, 24))
        self.assertEqual(len(client), 2 + client[0] * client[1])
        # `map-1-collision.bin`: byte0..3 = W, H, tileW, tileH, sau đó W*H ô.
        self.assertEqual((collision[0], collision[1]), (client[0], client[1]))
        self.assertEqual(len(collision), 4 + collision[0] * collision[1])

        self.assertEqual(config["width"], client[0])
        self.assertEqual(config["height"], client[1])
        self.assertEqual([collision[2], collision[3]], config["tile_size"])
        self.assertEqual(len(config["indexes"]), config["width"] * config["height"])
        self.assertEqual(len(config["client_solid"]), config["width"] * config["height"])

    def test_manifest_sinh_lai_dung_byte_luoi_va_cham_dang_chay(self):
        """Chốt chặn thật sự: build từ manifest KHÔNG được đổi artifact đang chạy.

        So từng byte thay vì chỉ so kích thước — lưới đúng kích thước mà sai ô nào
        đặc cũng đủ khoá người chơi tại chỗ.
        """
        from assetpack.map import encode_collision

        config = _map_config()
        width, height = config["width"], config["height"]
        tile_width, tile_height = config["tile_size"]
        spawn_index = (config["spawn"][1] // tile_height) * width + config["spawn"][0] // tile_width

        sinh_ra = encode_collision(
            width,
            height,
            config["client_solid"],
            spawn_index=spawn_index,
            tile_size=config["tile_size"],
        )
        self.assertEqual(sinh_ra, COLLISION.read_bytes())

    def test_y_nghia_tile_client_va_server_van_khop(self):
        """Mỗi mã tile phải nhất quán: hoặc luôn đặc, hoặc luôn đi được."""
        from assetpack.map import validate_collision_semantics

        manifest = json.loads(MAPPING.read_text(encoding="utf-8"))
        config = manifest["build_inputs"]["map"]
        tile_rows = [
            {"id": item["index"], "solid": item.get("solid", False)}
            for item in manifest["build_inputs"]["tile_art"]
        ]
        # Không raise = client/server cùng nghĩa tile.
        validate_collision_semantics(tile_rows, config["indexes"], config["client_solid"])

        # `tiles` (bảng nghĩa cấp cao) không được mâu thuẫn với `tile_art`.
        art_solid = {item["index"]: item.get("solid", False) for item in manifest["build_inputs"]["tile_art"]}
        for tile in manifest["tiles"]:
            self.assertEqual(tile["solid"], art_solid[tile["id"]], f"tile {tile['id']} lệch solid")
            self.assertEqual(tile["semantic"] == "solid", tile["solid"], f"tile {tile['id']} lệch semantic")

    def test_diem_hoi_sinh_phai_dung_duoc(self):
        """Spawn nằm trong lưới và không nằm trong ô đặc — nếu không, vào game là kẹt."""
        config = _map_config()
        width, height = config["width"], config["height"]
        tile_width, tile_height = config["tile_size"]
        col, row = config["spawn"][0] // tile_width, config["spawn"][1] // tile_height
        self.assertTrue(0 <= col < width, "spawn ra ngoài chiều rộng")
        self.assertTrue(0 <= row < height, "spawn ra ngoài chiều cao")
        self.assertFalse(config["client_solid"][row * width + col], "spawn rơi vào ô đặc")


if __name__ == "__main__":
    unittest.main()
