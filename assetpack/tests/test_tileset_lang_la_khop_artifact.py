"""Chốt kiểm: bảng tile client (`tile_set_info.bin`) build ra phải bằng bản đang chạy.

TẠI SAO có file này — quả mìn thứ hai cùng họ với `test_map1_manifest_khop_artifact.py`:

Bảng tile gốc của Ngọc Rồng chỉ khai báo tile **1–39**. Map Làng Lá dùng tile
**60–69 và 83–88**, nên với bảng gốc client đặt `types[] = 0` cho cả màn: không ô
nào là "sàn", **nhân vật rơi xuyên đất xuống y=576** (lỗi đã phải vá ngày 22/08).
Bản 1734 byte đang chạy trong `server/internal/game/assets/` là bản **sửa tay** bằng
`tools/map-lang-la/them-tile-vao-tileset.py`, trong khi `assetpack build` vẫn chép
thẳng bảng gốc 1696 byte từ `build_inputs.tile_set_info_path` — **ai build lại rồi
chép artifact sang máy chủ là tái sinh nguyên lỗi cũ, im lặng, không log**.

Bản vá: phép "thêm tile Làng Lá" được khai báo trong chính manifest
(`build_inputs.tile_set_overlay`) và do đường ống thi hành, nên nguồn vẫn là MỘT
(bảng gốc đã review + khai báo trong manifest) chứ không phải một file nhị phân
sửa tay không ai truy được nguồn.

Chốt kiểm ở đây buộc ba nơi cùng nói một chuyện:
  1. `docs/my-thuat/mapping-id.json` — bảng gốc golden + overlay (thứ build đọc)
  2. `assetpack/tileset.py`          — phép biến đổi
  3. `server/internal/game/assets/tile_set_info.bin` — bảng máy chủ đang chạy

Lệch một cái là đỏ ngay, thay vì đợi người chơi báo "vào game là rơi xuống hố".
"""

import copy
import hashlib
import json
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
MAPPING = ROOT / "docs/my-thuat/mapping-id.json"
# Bảng tile máy chủ đang phát cho client (gói -82). Đây là "sự thật đang chạy".
TILESET = ROOT / "server/internal/game/assets/tile_set_info.bin"

# Bit "đứng được lên trên" trong `types[]` của client. Thiếu bit này ở mọi tile
# của màn = không có sàn = rơi xuyên đất. Đây chính là bit mà quả mìn xoá mất.
BIT_SAN = 0x0002


def _manifest() -> dict:
    return json.loads(MAPPING.read_text(encoding="utf-8"))


def _baseline_bytes(manifest: dict) -> bytes:
    """Bảng gốc 1696 byte lấy từ golden nhúng trong manifest.

    KHÔNG đọc file gốc trên đĩa: `docs/my-thuat/anh-nguon/` bị `.gitignore` loại
    nên một git worktree hoàn toàn không có nó. `bytes_hex` là bản sao byte-đúng
    nằm TRONG git, và `validate_tileset_baseline` đã buộc file trên đĩa phải khớp
    nó ở lúc build — nên dùng golden ở đây là an toàn và chạy được ở mọi checkout.
    """
    return bytes.fromhex(manifest["build_inputs"]["tile_set_baseline"]["bytes_hex"])


class TilesetLangLaKhopArtifactTests(unittest.TestCase):
    def test_ban_dang_chay_giai_ma_roi_ma_hoa_lai_dung_tung_byte(self):
        """Hiểu đúng định dạng trước đã — đọc/ghi không khép kín thì mọi so sánh sau đều vô nghĩa."""
        from assetpack.tileset import decode_tile_set_info, encode_tile_set_info

        data = TILESET.read_bytes()
        self.assertEqual(encode_tile_set_info(decode_tile_set_info(data)), data)

    def test_baseline_goc_cong_overlay_sinh_ra_dung_ban_dang_chay(self):
        """Chốt chặn thật sự: build từ manifest KHÔNG được lùi bảng tile về bản gốc.

        So từng byte chứ không chỉ so kích thước — đủ dài mà thiếu đúng một tile
        sàn cũng làm người chơi rơi xuyên đất ở đúng chỗ đó.
        """
        from assetpack.tileset import apply_tile_set_overlay

        manifest = _manifest()
        sinh_ra = apply_tile_set_overlay(
            _baseline_bytes(manifest), manifest["build_inputs"]["tile_set_overlay"]
        )
        self.assertEqual(sinh_ra, TILESET.read_bytes())

    def test_manifest_khai_dung_ket_qua_cua_ban_dang_chay(self):
        """`result` trong manifest là chốt chặn ở LÚC BUILD, phải ghim đúng bản đang chạy."""
        manifest = _manifest()
        result = manifest["build_inputs"]["tile_set_overlay"]["result"]
        data = TILESET.read_bytes()
        self.assertEqual(result["byte_length"], len(data))
        self.assertEqual(result["sha256"], hashlib.sha256(data).hexdigest())
        # Và phải thật sự KHÁC bảng gốc — nếu bằng nhau nghĩa là overlay rỗng,
        # tức quả mìn đã quay lại dưới dạng "khai báo cho có".
        self.assertNotEqual(data, _baseline_bytes(manifest))

    def test_overlay_chi_dung_vao_bo_tile_cua_map_1_khong_dung_bo_khac(self):
        """39 bộ tile còn lại là di sản Ngọc Rồng chưa ai review — không được chạm."""
        from assetpack.tileset import decode_tile_set_info

        manifest = _manifest()
        goc = decode_tile_set_info(_baseline_bytes(manifest))
        chay = decode_tile_set_info(TILESET.read_bytes())
        index = manifest["build_inputs"]["tile_set_overlay"]["tile_set_index"]
        self.assertEqual(len(goc), len(chay))
        for i, (bo_goc, bo_chay) in enumerate(zip(goc, chay)):
            if i == index:
                continue
            self.assertEqual(bo_goc, bo_chay, f"bộ tile {i} bị đổi ngoài ý muốn")
        self.assertNotEqual(goc[index], chay[index], "bộ tile map 1 không hề được thêm gì")

    def test_moi_tile_lang_la_dac_va_co_dung_deu_nam_trong_nhom_san(self):
        """Kiểm ở tầng Ý NGHĨA, không chỉ tầng hash.

        Máy chủ chặn di chuyển theo `tile_art[].solid`; client cho đứng theo bit
        sàn. Tile nào máy chủ bảo "đặc" mà client không khai là sàn thì hai bên
        nói ngược nhau — người chơi rơi xuyên đất rồi bị máy chủ kéo lại. Chỉ soi
        tile MỚI (id > số tile bảng gốc khai) vì tile 1–39 thuộc bảng gốc đã có
        sẵn ngữ nghĩa, không phải việc của Làng Lá.
        """
        from assetpack.tileset import decode_tile_set_info

        manifest = _manifest()
        inputs = manifest["build_inputs"]
        so_tile_goc = inputs["tile_set_baseline"]["semantic"]["tile_count"]
        dac = {tile["index"] for tile in inputs["tile_art"] if tile.get("solid")}
        dung = set(inputs["map"]["indexes"])
        can_co = sorted(t for t in dac & dung if t > so_tile_goc)
        self.assertTrue(can_co, "map không dùng tile Làng Lá đặc nào — kiểm tra vô nghĩa")

        bo = decode_tile_set_info(TILESET.read_bytes())[inputs["tile_set_overlay"]["tile_set_index"]]
        san = {tile for loai, tiles in bo if loai == BIT_SAN for tile in tiles}
        self.assertEqual([t for t in can_co if t not in san], [], "tile đặc thiếu bit sàn")

    def test_bo_bot_mot_tile_lang_la_thi_chot_kiem_do(self):
        """Kiểm thử đột biến: cố tình đánh rơi tile 60 khỏi nhóm sàn.

        Phải đỏ ở CẢ HAI tầng — hash (`result` lệch) và ý nghĩa (thiếu bit sàn).
        Chốt kiểm nào chỉ đỏ ở một tầng là chốt kiểm dễ bị vá vòng.
        """
        from assetpack.tileset import TilesetError, apply_tile_set_overlay, decode_tile_set_info

        manifest = _manifest()
        hong = copy.deepcopy(manifest["build_inputs"]["tile_set_overlay"])
        for nhom in hong["groups"]:
            if nhom["bit"] == BIT_SAN:
                nhom["tiles"] = [t for t in nhom["tiles"] if t != 60]
        with self.assertRaisesRegex(TilesetError, "ASSET-TILESET-013"):
            apply_tile_set_overlay(_baseline_bytes(manifest), hong)

        # Tầng ý nghĩa: bỏ luôn `result` để qua được chốt hash, vẫn phải bị bắt.
        hong_khong_result = copy.deepcopy(hong)
        hong_khong_result["result"] = None
        du_lieu = apply_tile_set_overlay(
            _baseline_bytes(manifest), hong_khong_result, check_result=False
        )
        bo = decode_tile_set_info(du_lieu)[hong["tile_set_index"]]
        san = {tile for loai, tiles in bo if loai == BIT_SAN for tile in tiles}
        self.assertNotIn(60, san)

    def test_overlay_khong_duoc_khai_lai_tile_da_co_san_trong_bang_goc(self):
        """Khai trùng = hoặc bảng gốc đã đổi, hoặc người viết hiểu nhầm. Cả hai đều phải dừng."""
        from assetpack.tileset import TilesetError, apply_tile_set_overlay

        manifest = _manifest()
        hong = copy.deepcopy(manifest["build_inputs"]["tile_set_overlay"])
        # Tile 1 đã nằm sẵn trong nhóm sàn của bảng gốc Ngọc Rồng.
        for nhom in hong["groups"]:
            if nhom["bit"] == BIT_SAN:
                nhom["tiles"] = [1] + nhom["tiles"]
        with self.assertRaisesRegex(TilesetError, "ASSET-TILESET-012"):
            apply_tile_set_overlay(_baseline_bytes(manifest), hong)


if __name__ == "__main__":
    unittest.main()
