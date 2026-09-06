"""Sinh production closure Lang Lá từ đúng nguồn đã review.

Module này cố ý không tự gán PASS cho nguồn mới. Mọi sprite canonical đều được
crop từ một asset đã có trong ``vertical-slice-v1.json``; các input còn lại được
copy/chốt SHA vào một closure bất biến dưới ArtRoot dùng chung.
"""

from __future__ import annotations

import argparse
from collections import deque
import dataclasses
import hashlib
import json
import shutil
import tempfile
from pathlib import Path
from typing import Any

import numpy as np
from PIL import Image

from assetpack.animation import REQUIRED_ANIMATION_SHEETS, REQUIRED_CLIENT_CF
from assetpack.build import _source_paths
from assetpack.visual import (
    render_hud_golden,
    render_iruka_composite_contact_sheet,
)
from tools.artqa import nhan_vat_moc, nhan_vat_toi_uu
from tools.artqa.artqa import _adaptive_chroma_mask, extract_alpha
from tools.artqa.generate_tile_sheet import generate_tile_sheet


ZOOMS = ("x1", "x2", "x3", "x4")
ART_PREFIX = Path("docs/my-thuat/anh-nguon/vertical-slice-v1")
CLOSURE_RELATIVE = ART_PREFIX / "generated/production-closure-v11"
APPROVED_UI_RELATIVE = Path(
    "my thuat/dot1-nhom1-duyet/dot1-nhom1-duyet/01-giao-dien"
)
TILE_SOURCE_SHA256 = "2031edef163fc766b9f44df6c30724c0e41c57730b4f98749e3940fae157036b"

PLAYER_ACTIONS = (
    ("idle", "player_idle_3f_chroma.png", 3),
    ("run", "player_run_5f_chroma.png", 5),
    ("melee", "player_melee_4f_chroma.png", 4),
    ("cast", "player_cast_3f_chroma.png", 3),
    ("hit_death_getup", "player_hit_death_getup_3f_chroma.png", 3),
)
# Closure V1 đã ký vẫn cần tái lập được khi review cũ chỉ có run 4 khung. Lối này
# chỉ là tương thích dữ liệu cũ; C1 luôn ưu tiên bộ tên mới ở trên.
LEGACY_PLAYER_ACTIONS = (
    ("idle", "player_idle_3f_chroma.png", 3),
    ("run", "player_run_4f_chroma.png", 4),
    ("melee", "player_melee_4f_chroma.png", 4),
    ("cast", "player_cast_3f_chroma.png", 3),
    ("hit_death_getup", "player_hit_death_getup_3f_chroma.png", 3),
)

# --- Hình học nhân vật người chơi (part 30 đầu / 57 thân / 58 chân) ------------
#
# Client vẽ ba mảnh rời theo Char.paintCharBody (Char.cs ~6057-6084), neo TOP|LEFT
# khi quay phải:
#     x = cx + CharInfo[cf][k][1] + part.pi[idx].dx
#     y = cy - CharInfo[cf][k][2] + part.pi[idx].dy      (y là MÉP TRÊN ảnh)
# Nên ba thứ phải cùng kể một câu chuyện: tỉ lệ ảnh, gốc toạ độ ảnh, và dx/dy.
#
# Bản trước dùng ``fit_to_canvas`` co TỪNG khung cho vừa hộp 128×160/128×128 nên
# mỗi khung có một tỉ lệ khác nhau (khung nằm rộng 665 px co còn 1/5 khung đứng
# rộng 340 px) → nhân vật phình/teo theo pose, bàn chân trôi -10..-26 px so với
# cy, đầu hở +7 px ở cf 20/31. Bản này:
#   * MỘT tỉ lệ ``s`` cho mọi khung, mọi mảnh (giữ khung idle 0 đúng 128 px ở x4,
#     tức nhân vật trong game không đổi cỡ);
#   * gốc toạ độ của mỗi khung = giữa-dưới bbox nội dung CẢ KHUNG ("bàn chân"),
#     nên mọi mảnh của cùng khung nằm trong CÙNG hệ toạ độ như Iruka;
#   * canvas = đúng vùng đã scale sau khi trim alpha (không nhồi vào hộp cố định),
#     vì SmallImage của client vẽ ảnh mọi kích cỡ, và trim làm ảnh đối xứng khi
#     client lật ngang (hướng trái neo TOP|RIGHT);
#   * dx/dy của từng frame part bù đúng CharInfo của cf neo.
PLAYER_REFERENCE_SOURCE = "player_idle_3f_chroma.png"
PLAYER_REFERENCE_FRAME = 0
PLAYER_REFERENCE_X4_WIDTH = 128
# Bàn chân (mép dưới vùng đục của chân) so với cy, đơn vị x1. -11 là đúng giá trị
# của cf 0 trong bản đang chạy: giữ nguyên để nhân vật không đột ngột đổi độ cao
# so với mặt đất; đổi MỘT hằng số này là dời cả nhân vật lên/xuống đều nhau.
PLAYER_FEET_OFFSET_X1 = -11
COMPONENT_CHARINFO_SLOT = {"head": 0, "leg": 1, "body": 2}

# --- V1 (05/09/2026): tỉ lệ TỪNG KHUNG theo cỡ đầu + cắt theo cằm + gán lại slot -----
#
# Chẩn đoán lại "đầu tách thân" (docs/my-thuat/review/khe-dau-than-2026-09-05/GHI-CHU.md):
# thứ nhìn thấy là HAI KHUÔN MẶT — mảnh thân cắt ở 28 % chiều cao khung, ở khung
# chạy/đánh (nguồn vẽ nhỏ hơn idle 12–27 % và khom hơn) lát cắt rơi lên trán nên mảnh
# thân chứa cả mặt và client vẽ nó ĐÈ lên đầu idle. Ba tầng lệch, ba phép sửa dữ liệu:
#   1. s_f = s0 / (cỡ đầu khung / cỡ đầu idle) — đo 6 thước cứng của đầu
#      (tools/artqa/nhan_vat_moc.py); đầu idle dùng chung 33 cf nên thân phải khớp cỡ đầu.
#   2. Cắt ĐẦU từ đỉnh tóc tới cằm + 6 px nguồn; cắt THÂN từ cằm − 2 px x1 tới 76 %;
#      chân giữ 66–100 %. Thân KHÔNG được chứa mặt.
#   3. Gán lại slot theo cách client THẬT dùng (Char.cs:2537 `cf = gameTick % 5 + 2`
#      → cf 2..6 là 5 khung CHẠY; cf 20/31 = thế bắt đầu skill, client không bao giờ vẽ
#      nằm). Ba khung hit_death_getup không dùng (hit0 lệch tâm mặt 15 px).
# dx/dy không còn bù từng mảnh theo cf neo mà giải MỘT LẦN cho cả 34 mảnh bằng MILP
# (tools/artqa/nhan_vat_toi_uu.py): khe đầu–thân ∈ [−9, 0], thân–chân ∈ [−11, −2], ghim
# dy đầu slot 0 = −38 (−39 vô nghiệm). Giá đã chấp nhận: nhân vật thấp đi ~10 % khi chạy.
PLAYER_HEAD_SLOTS_V1 = (("idle", 0), ("idle", 1), ("idle", 2))
PLAYER_BODY_SLOTS_V1 = {
    0: ("idle", 0), 1: ("idle", 1),
    2: ("run", 0), 3: ("run", 1), 4: ("run", 2), 5: ("run", 3), 6: ("run", 4),
    7: ("melee", 0), 8: ("melee", 1), 9: ("melee", 2), 10: ("melee", 3),
    11: ("cast", 0), 12: ("cast", 1), 13: ("cast", 2),
    14: ("hit_death_getup", 0), 15: ("hit_death_getup", 1), 16: ("hit_death_getup", 2),
}
PLAYER_LEG_SLOTS_V1 = {
    0: ("idle", 0), 1: ("idle", 1),
    2: ("run", 0), 3: ("run", 1), 4: ("run", 2), 5: ("run", 3), 6: ("run", 4),
    7: ("melee", 0), 8: ("melee", 1), 9: ("melee", 2), 10: ("melee", 3),
    11: ("cast", 0), 12: ("cast", 1), 13: ("cast", 2),
}
LEGACY_PLAYER_BODY_SLOTS_V1 = {
    0: ("idle", 0), 1: ("idle", 1),
    2: ("run", 0), 3: ("run", 1), 4: ("run", 2), 5: ("run", 3), 6: ("run", 0),
    7: ("melee", 0), 8: ("melee", 1), 9: ("melee", 2), 10: ("melee", 3),
    11: ("cast", 0), 12: ("cast", 1), 13: ("cast", 2),
    14: ("cast", 2), 15: ("cast", 0), 16: ("idle", 2),
}
LEGACY_PLAYER_LEG_SLOTS_V1 = {
    0: ("idle", 0), 1: ("idle", 1),
    2: ("run", 0), 3: ("run", 1), 4: ("run", 2), 5: ("run", 3), 6: ("run", 0),
    7: ("melee", 0), 8: ("melee", 1), 9: ("melee", 2), 10: ("melee", 3),
    11: ("cast", 0), 12: ("cast", 1), 13: ("idle", 2),
}
PLAYER_HEAD_SCALE_RANGE_C1 = (0.95, 1.05)
HEAD_CHIN_MARGIN_SRC = 6      # px nguồn dưới cằm gộp vào ĐẦU (cổ áo che mối nối)
BODY_CHIN_MARGIN_X1 = 2.0     # px x1 trên cằm gộp vào THÂN (đầu vẽ trước, thân đè lên)
BODY_BOTTOM_RATIO = 0.76
LEG_TOP_RATIO = 0.66
# Lát cắt cũ theo phần trăm — chỉ còn dùng cho khung không có cằm (nằm) và Iruka.
HEAD_BOTTOM_RATIO_CU = 0.42
BODY_TOP_RATIO_CU = 0.28
PLAYER_MILP_V1 = nhan_vat_toi_uu.ThamSoMILP(ban_chan_x1=PLAYER_FEET_OFFSET_X1)
PLAYER_GEOMETRY_VERSION = "V1 05/09"
CHARINFO_PATH = (
    Path(__file__).resolve().parents[1] / "nhan-vat" / "charinfo.json"
)

# --- Iruka V1 (06/09/2026): NPC part 2189 đầu / 2190 thân / 2191 chân → 27075/27076/27077 --
#
# Client vẽ NPC bằng ĐÚNG phép toán của người chơi (Npc.cs:305-313, cùng CharInfo) nhưng
# chỉ đảo cf 0/1 (Npc.cs:179-185); CharInfo cf 1 nâng đầu + thân 1 px cùng nhau, chân đứng
# yên. Bản đang chạy sinh Iruka bằng `_generate_iruka_components` (hình học CŨ): ép cả
# khung vào canvas 128×160 → Iruka chỉ 38 px x1 = 62 % người chơi, part 2189–2191 ghi
# dx = dy = 0 → đầu lửng +13 px, thân chìm +29 px. Bản V1 dùng lại bộ hàm người chơi:
#   * s_iruka = s0 / r × IRUKA_TI_LE_DAU. r = cỡ đầu Iruka / cỡ đầu người chơi = median 6
#     thước cứng, đo trên CÙNG tham chiếu trung bình 3 khung idle mà `player_frame_geometry`
#     dùng — nhờ vậy IRUKA_TI_LE_DAU = 1.0 đúng nghĩa "đầu Iruka bằng đầu người chơi trên
#     màn" (đo 06/09: r = 1,6642; nếu lấy idle0 lẻ thì 1,6382).
#   * Lát cắt theo cằm y như người chơi (`_component_rect_v1`), trim alpha, không hộp cố định.
#   * dx/dy KHÔNG cần MILP: một khung, ba mảnh cùng hệ toạ độ → `component_placement` tại
#     cf 0; cf 1 lệch đúng 1 px theo CharInfo, ghi ở `charinfo_residual`. Mọi frame của
#     part 2189 (3) / 2190 (17) / 2191 (14) cùng trỏ một sprite nên cùng dx/dy.
IRUKA_TI_LE_DAU = 1.0
IRUKA_SOURCE = "master_iruka_x4_chroma.png"
IRUKA_IDS = (27075, 27076, 27077)
IRUKA_PART_IDS = {"head": 2189, "body": 2190, "leg": 2191}
IRUKA_CHARINFO_CF = (0, 1)
IRUKA_GEOMETRY_VERSION = "V1 06/09"
# 6 thước cứng cho Iruka: bỏ `head_h_hair_chin` (đuôi tóc dựng của Iruka cao hơn tóc người
# chơi ~100 px nguồn → tỉ lệ 1,37 là thước lạc, kéo median lệch), thay bằng sqrt diện tích
# mặt để vẫn đủ 6 thước cho median.
IRUKA_THUOC_CUNG = ("face_h", "face_w", "face_area_sqrt", "eye_h", "eye_to_chin", "band_thick")

# ── Khe nền còn sót trên mảnh nhân vật (06/09/2026) ──────────────────────────
# Triệu chứng: VỆT HỒNG nhỏ trên nhân vật trong game (rõ nhất ở khe giữa cánh tay
# và thắt lưng khung đứng, và trên đầu Iruka). Nguyên nhân: `extract_alpha` chỉ
# flood magenta CHẠM MÉP ảnh nên khe kín giữ nguyên magenta ĐỤC (alpha 255), rồi
# `remove_large_chroma_islands` chỉ xoá đảo ≥ 598 px nên khe nhỏ sống sót — 293 px
# hồng trên 37 mảnh ×4 zoom.
#
# Vì sao dám xoá SẠCH trên các dải này: trang phục nhân vật không có tím/hồng
# thật. Histogram 06/09 trên `player_idle` (360 652 px đục) và `master_iruka`
# (323 075 px đục): mọi pixel hue 280–340° đều có kênh lục 1–5 — tức viền răng cưa
# của chính nền khoá màu chứ không phải vải (áo navy, khăn cam, tóc nâu/đen).
# Đo thêm: `is_chroma` nghiêm ngặt ⊂ mặt nạ candidate trên cả 6 dải (0 pixel lọt),
# nên xoá hết đảo candidate là đủ để đếm chroma về 0.
NGUON_NHAN_VAT = frozenset(
    {source for _, source, _ in PLAYER_ACTIONS}
    | {source for _, source, _ in LEGACY_PLAYER_ACTIONS}
    | {IRUKA_SOURCE}
)
# Đảo ≥ 8 px là KHE NỀN thật → alpha 0 (nền vốn phải trong suốt). Đo 06/09 trên
# 6 dải: khe nhỏ nhất 9 px, đốm lẻ lớn nhất 4 px — khoảng trống 4↔9 tách hai loại
# dứt khoát. Đốm < 8 px được VÁ MÀU (giữ alpha) để không thủng áo bằng lỗ kim.
KHE_NEN_TOI_THIEU = 8

WOLF_ACTIONS = (
    ("idle", "wolf_idle_2f_chroma.png", 2),
    ("run", "wolf_run_4f_chroma.png", 4),
    ("pounce", "wolf_pounce_2f_chroma.png", 2),
    ("hit_death", "wolf_hit_death_2f_chroma.png", 2),
)
UI_ASSETS = (
    "ui_logo_game_x4.png",
    "ui_khung_avatar_x4.png",
    "ui_khung_hoithoai_x4.png",
    "ui_nut_thuong_x4.png",
    "ui_nut_nhan_x4.png",
    "ui_nut_khoa_x4.png",
    "ui_o_vatpham_x4.png",
    "ui_o_dangchon_x4.png",
)


class ClosureError(ValueError):
    """Input của production closure thiếu hoặc lệch provenance."""


def _sha(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def _review_by_name(review: dict[str, Any]) -> dict[str, dict[str, Any]]:
    assets = review.get("assets")
    if not isinstance(assets, list):
        raise ClosureError("ASSET-CLOSURE-015: review thiếu assets[]")
    result = {}
    for record in assets:
        if not isinstance(record, dict) or not isinstance(record.get("asset"), str):
            raise ClosureError("ASSET-CLOSURE-015: review asset malformed")
        if record.get("status") != "PASS":
            continue
        result[record["asset"]] = record
    return result


def _frame_rect(record: dict[str, Any], frame_index: int) -> list[int]:
    crop = record.get("checks", {}).get("crop", {})
    edges = crop.get("source_edges")
    boxes = crop.get("content_boxes")
    if (
        not isinstance(edges, list)
        or not isinstance(boxes, list)
        or frame_index < 0
        or frame_index >= len(boxes)
        or len(edges) != len(boxes) + 1
    ):
        raise ClosureError(
            f"ASSET-CLOSURE-015: crop review malformed: {record.get('asset')}"
        )
    left, top, right, bottom = boxes[frame_index]
    return [
        edges[frame_index] + left,
        top,
        right - left,
        bottom - top,
    ]


def _component_rect(
    record: dict[str, Any], frame_index: int, component: str
) -> list[int]:
    x, y, width, height = _frame_rect(record, frame_index)
    # Lát cắt theo % chiều cao khung — CHỈ còn dùng cho Iruka và khung không có cằm.
    # Với người chơi, 28 % ở khung chạy/đánh rơi lên trán → thân chứa mặt (hai khuôn
    # mặt); xem `_component_rect_v1`.
    if component == "head":
        top = 0
        bottom = max(1, round(height * HEAD_BOTTOM_RATIO_CU))
    elif component == "body":
        top = round(height * BODY_TOP_RATIO_CU)
        bottom = max(top + 1, round(height * BODY_BOTTOM_RATIO))
    elif component == "leg":
        top = round(height * LEG_TOP_RATIO)
        bottom = height
    else:
        raise ClosureError(f"ASSET-CLOSURE-015: component lạ: {component}")
    return [x, y + top, width, bottom - top]


def _component_rect_v1(
    record: dict[str, Any], frame_index: int, component: str, khung: dict[str, Any]
) -> list[int]:
    """Lát cắt V1 theo CẰM của chính khung (đơn vị nguồn), không theo % chiều cao.

    ``khung`` = một mục của ``nhan_vat_moc.hinh_hoc_khung`` (``chin_y`` tuyệt đối theo
    dải, ``scale`` = s_f). Khung không có cằm (nằm) rơi về lát cắt % cũ.
    """

    x, y, width, height = _frame_rect(record, frame_index)
    chin = khung.get("chin_y")
    scale = float(khung["scale"])
    if chin is None:
        return _component_rect(record, frame_index, component)
    chin_rel = float(chin) - y
    if component == "head":
        top, bottom = 0.0, chin_rel + HEAD_CHIN_MARGIN_SRC
    elif component == "body":
        top = chin_rel - BODY_CHIN_MARGIN_X1 * 4.0 / scale
        bottom = height * BODY_BOTTOM_RATIO
    elif component == "leg":
        top, bottom = height * LEG_TOP_RATIO, float(height)
    else:
        raise ClosureError(f"ASSET-CLOSURE-015: component lạ: {component}")
    top_i = min(max(0, round(top)), height - 1)
    bottom_i = min(max(top_i + 1, round(bottom)), height)
    return [x, y + top_i, width, bottom_i - top_i]


def player_frame_geometry(
    records: dict[str, dict[str, Any]],
    moc: dict[tuple[str, int], dict[str, Any]],
    *,
    idle_frames: tuple[int, ...] = (0, 1, 2),
) -> tuple[dict[tuple[str, int], dict[str, Any]], dict[str, float], float]:
    """Từng khung nguồn: hệ số đầu, s_f, cằm, tâm mặt. Trả (hình_học, tham_chiếu_idle, s0)."""

    reference = _frame_rect(records[PLAYER_REFERENCE_SOURCE], PLAYER_REFERENCE_FRAME)
    s0 = player_uniform_scale(reference[2])
    tham_chieu = nhan_vat_moc.tham_chieu_idle(moc, frames=idle_frames)
    return nhan_vat_moc.hinh_hoc_khung(moc, tham_chieu, s0), tham_chieu, s0


def _player_actions_for_records(records: dict[str, dict[str, Any]]) -> tuple[tuple[str, str, int], ...]:
    """Ưu tiên C1; chỉ rơi về V1 để tái lập artifact V1 đã được ký."""

    if all(source_name in records for _, source_name, _ in PLAYER_ACTIONS):
        return PLAYER_ACTIONS
    if all(source_name in records for _, source_name, _ in LEGACY_PLAYER_ACTIONS):
        return LEGACY_PLAYER_ACTIONS
    raise ClosureError("ASSET-CLOSURE-015: review thiếu đủ bộ dải người chơi C1 hoặc V1")


def _player_slot_maps(actions: tuple[tuple[str, str, int], ...]) -> tuple[dict[int, tuple[str, int]], dict[int, tuple[str, int]]]:
    if actions == PLAYER_ACTIONS:
        return PLAYER_BODY_SLOTS_V1, PLAYER_LEG_SLOTS_V1
    return LEGACY_PLAYER_BODY_SLOTS_V1, LEGACY_PLAYER_LEG_SLOTS_V1


def _require_c1_same_head_size(
    geometry: dict[tuple[str, int], dict[str, Any]],
    actions: tuple[tuple[str, str, int], ...],
) -> None:
    """C1 từ chối run/hit có đầu lệch quá ±5 % so với idle chuẩn.

    V1 chỉ dùng bù tỷ lệ bằng toán nên không thể áp tiêu chuẩn mới vào artifact đã ký.
    C1 đã chốt vẽ lại run + hit, vì vậy phải chặn ngay ở closure thay vì lại co giãn
    từng pose rồi để nhân vật thấp hoặc xuất hiện hai khuôn mặt trong game.
    """

    if actions != PLAYER_ACTIONS:
        return
    minimum, maximum = PLAYER_HEAD_SCALE_RANGE_C1
    for action, source_name, frame_count in actions:
        if action not in {"run", "hit_death_getup"}:
            continue
        for frame_index in range(frame_count):
            scale = geometry.get((source_name, frame_index), {}).get("he_so_dau")
            if scale is None or not minimum <= float(scale) <= maximum:
                raise ClosureError(
                    "ASSET-CLOSURE-015: cùng cỡ đầu C1 lệch ở "
                    f"{source_name}#{frame_index}: {scale!r}; "
                    f"cần trong [{minimum}, {maximum}] so với idle"
                )


def derive_player_slot_specs(
    review: dict[str, Any],
    *,
    moc: dict[tuple[str, int], dict[str, Any]],
    records: dict[str, dict[str, Any]] | None = None,
    idle_frames: tuple[int, ...] = (0, 1, 2),
) -> list[dict[str, Any]]:
    """Trả đúng 3/17/14 slot theo bảng V1, lát cắt theo cằm, tỉ lệ từng khung.

    ``moc`` = mốc giải phẫu của 17 khung nguồn (``nhan_vat_moc.do_nguon``). Mỗi spec
    mang thêm ``scale_frame`` (s_f) và ``nguon`` (cằm, tâm mặt, hệ số đầu — đơn vị nguồn).
    """

    records = records if records is not None else _review_by_name(review)
    actions = _player_actions_for_records(records)
    by_action: dict[str, tuple[str, dict[str, Any]]] = {}
    for action, source_name, expected_count in actions:
        record = records.get(source_name)
        if record is None:
            raise ClosureError(
                f"ASSET-CLOSURE-015: nguồn player chưa PASS: {source_name}"
            )
        boxes = record.get("checks", {}).get("crop", {}).get("content_boxes", [])
        if len(boxes) != expected_count:
            raise ClosureError(
                f"ASSET-CLOSURE-015: frame count lệch review: {source_name}"
            )
        by_action[action] = (source_name, record)
    hinh_hoc, tham_chieu, s0 = player_frame_geometry(records, moc, idle_frames=idle_frames)
    _require_c1_same_head_size(hinh_hoc, actions)
    body_slots, leg_slots = _player_slot_maps(actions)

    groups = (
        ("head", 27041, dict(enumerate(PLAYER_HEAD_SLOTS_V1))),
        ("body", 27044, body_slots),
        ("leg", 27061, leg_slots),
    )
    slots = []
    for component, first_id, bang in groups:
        for slot_index in range(len(bang)):
            action, frame_index = bang[slot_index]
            source_name, record = by_action[action]
            key = (source_name, frame_index)
            if key not in hinh_hoc:
                raise ClosureError(f"ASSET-CLOSURE-015: thiếu mốc giải phẫu khung {key}")
            khung = hinh_hoc[key]
            x, y, width, height = _frame_rect(record, frame_index)
            slots.append(
                {
                    "id": first_id + slot_index,
                    "component": component,
                    "slot_index": slot_index,
                    "action": action,
                    "action_slot": frame_index,
                    "frame_index": frame_index,
                    "source_path": source_name,
                    "source_sha256": record["source_sha256"],
                    "source_rect": _component_rect_v1(record, frame_index, component, khung),
                    "scale_frame": float(khung["scale"]),
                    "nguon": {
                        "frame_top": y,
                        "frame_height": height,
                        "edge_x": int(moc[key]["edge_x"]),
                        "he_so_dau": khung["he_so_dau"],
                        "chin_y": khung["chin_y"],
                        "cam_uoc_luong": bool(khung["cam_uoc_luong"]),
                        "face_cx": khung["face_cx"],
                        "hair_top": khung["hair_top"],
                        "thuoc_loai": list(khung["thuoc_loai"]),
                    },
                    "anchor": "feet-bottom-center",
                    "layer": {"head": 2, "body": 1, "leg": 0}[component],
                    "overlap_ratio": {
                        "head_body": 0.14,
                        "body_leg": 0.10,
                    },
                    "client_cf": list(
                        REQUIRED_CLIENT_CF[f"hoc-vien-{component}"][slot_index]
                    ),
                }
            )
    return slots


def load_charinfo(path: Path | None = None) -> list[list[list[int]]]:
    """Bảng khung xương 33 cf trích từ ``Char.cs`` (nguồn sự thật của client)."""

    payload = json.loads(Path(path or CHARINFO_PATH).read_text(encoding="utf-8"))
    if not isinstance(payload, list) or len(payload) != 33:
        raise ClosureError("ASSET-CLOSURE-015: charinfo phải có đúng 33 cf")
    return payload


def charinfo_slot_to_cf(
    charinfo: list[list[list[int]]], component: str
) -> dict[int, list[int]]:
    """Slot (chỉ số frame trong part) → các cf THẬT dùng nó, theo CharInfo.

    Không dùng ``REQUIRED_CLIENT_CF`` cho việc này: bảng đó lệch một bậc so với
    ``CharInfo`` ở thân và chân (ví dụ nó nói slot thân 0 phục vụ cf 0/1, còn
    ``CharInfo[0][2][0]`` = 1 nên cf 0 thật sự lấy slot thân 1). Bù dx/dy phải
    dựa trên cái client THẬT SỰ tra, nếu không mọi phép bù đều lệch một khung.
    """

    if component not in COMPONENT_CHARINFO_SLOT:
        raise ClosureError(f"ASSET-CLOSURE-015: component lạ: {component}")
    slot = COMPONENT_CHARINFO_SLOT[component]
    result: dict[int, list[int]] = {}
    for cf, entry in enumerate(charinfo):
        result.setdefault(entry[slot][0], []).append(cf)
    return result


def player_uniform_scale(
    reference_frame_width: int,
    target_x4_width: int = PLAYER_REFERENCE_X4_WIDTH,
) -> float:
    """s0 = 128/340 là tỉ lệ MỐC (khung idle 0).

    Từ V1 05/09 KHÔNG còn dùng chung cho mọi khung: mỗi khung vẽ với
    ``scale_frame = s0 / he_so_dau`` (hệ số cỡ đầu khung / cỡ đầu idle, xem
    ``nhan_vat_moc``), để đầu idle dùng chung 33 cf khớp với thân/chân của khung đó.
    """

    if reference_frame_width <= 0 or target_x4_width <= 0:
        raise ClosureError("ASSET-CLOSURE-015: bề rộng khung tham chiếu phải > 0")
    return target_x4_width / reference_frame_width


def component_placement(
    *,
    component_left_x4: float,
    component_top_x4: float,
    frame_center_x4: float,
    frame_feet_x4: float,
    charinfo_row: list[int],
    feet_offset_x1: int = PLAYER_FEET_OFFSET_X1,
) -> dict[str, int]:
    """dx/dy (đơn vị x1) cho một frame part, bù đúng CharInfo của cf neo.

    ``*_x4`` là toạ độ trong ảnh x4 đã scale; chia 4 để về đơn vị x1 mà part.bin
    dùng. ``left_rel``/``top_rel`` là vị trí mong muốn của góc trên-trái ảnh so
    với (cx, cy); âm là sang trái / lên trên.
    """

    left_rel = round((component_left_x4 - frame_center_x4) / 4)
    top_rel = round((component_top_x4 - frame_feet_x4) / 4) + feet_offset_x1
    dx = left_rel - charinfo_row[1]
    dy = top_rel + charinfo_row[2]
    if not -128 <= dx <= 127 or not -128 <= dy <= 127:
        raise ClosureError(
            f"ASSET-CLOSURE-015: dx/dy ngoài signed byte: dx={dx} dy={dy}"
        )
    return {"left_rel": left_rel, "top_rel": top_rel, "dx": dx, "dy": dy}


def zoom_variants(x4: Image.Image) -> dict[str, Image.Image]:
    """x4 → bốn mức zoom; x1 đúng bằng 1/4 nên đo trên x1 khớp part.bin."""

    result = {}
    for zoom_index, zoom in enumerate(ZOOMS, start=1):
        if zoom == "x4":
            result[zoom] = x4
            continue
        size = (
            max(1, round(x4.width * zoom_index / 4)),
            max(1, round(x4.height * zoom_index / 4)),
        )
        result[zoom] = x4.resize(size, Image.Resampling.NEAREST)
    return result


def _hip_cx_x1(piece: Image.Image, component: str) -> float:
    """Tâm hông (x1, so với mép trái ảnh): thân = 12 hàng x4 dưới cùng, chân = 12 hàng trên cùng.

    Dùng để căn ngang thân–chân khi hai mảnh lấy từ HAI khung khác nhau (không có
    "tâm khung" chung để so).
    """

    alpha = np.asarray(piece.getchannel("A")) > 0
    rows = alpha[-12:] if component == "body" else alpha[:12]
    xs = np.where(rows.any(axis=0))[0]
    if len(xs) == 0:
        return piece.width / 8.0
    return float(xs.mean() / 4.0)


def player_component_images(
    review: dict[str, Any],
    *,
    art_root: Path,
    records: dict[str, dict[str, Any]] | None = None,
    alpha_cache: dict[str, Image.Image] | None = None,
    charinfo: list[list[list[int]]] | None = None,
    feet_offset_x1: int = PLAYER_FEET_OFFSET_X1,
    moc: dict[tuple[str, int], dict[str, Any]] | None = None,
    tham_so_milp: nhan_vat_toi_uu.ThamSoMILP | None = None,
    idle_frames: tuple[int, ...] = (0, 1, 2),
) -> list[dict[str, Any]]:
    """34 mảnh người chơi (V1): ảnh x4 + mốc giải phẫu + dx/dy giải bằng MILP.

    Trả về thứ tự slot 27041..27074. Mỗi phần tử là spec của
    ``derive_player_slot_specs`` cộng thêm:
      * ``image_x4`` — mảnh đã cắt theo cằm, scale s_f, trim alpha;
      * ``moc`` — mốc giải phẫu đơn vị x1 so với góc trên-trái ảnh (``chin_off``,
        ``feet_off``, ``hair_off``, ``face_cx_off``, ``center_off``, ``hip_cx_off``,
        ``bbox_x1``) — chính là đầu vào của MILP;
      * ``placement`` — ``dx/dy`` từ nghiệm MILP (một lần cho cả 34 mảnh) và
        ``left_rel/top_rel`` tương ứng tại cf neo;
      * ``charinfo_cf`` (cf thật dùng slot), ``charinfo_anchor_cf`` (cf nhỏ nhất),
        ``charinfo_residual`` — nay là sai số giải phẫu dự đoán TỪNG cf (e/d/ex/gx, khe);
      * ``milp`` — tham số + nghiệm + bảng cf (chung cho mọi phần tử);
      * ``shared_transform`` — s0, s_f, gốc khung.

    ``moc`` bỏ trống thì đo thẳng từ ảnh nguồn (``nhan_vat_moc.do_nguon``).
    """

    records = records if records is not None else _review_by_name(review)
    alpha_cache = alpha_cache if alpha_cache is not None else {}
    charinfo = charinfo if charinfo is not None else load_charinfo()
    tham_so = tham_so_milp if tham_so_milp is not None else dataclasses.replace(
        PLAYER_MILP_V1, ban_chan_x1=feet_offset_x1
    )
    if tham_so.ban_chan_x1 != feet_offset_x1:
        raise ClosureError(
            f"ASSET-CLOSURE-015: MILP ban_chan_x1={tham_so.ban_chan_x1} khác feet_offset_x1={feet_offset_x1}"
        )
    actions = _player_actions_for_records(records)
    if moc is None:
        try:
            moc = nhan_vat_moc.do_nguon(
                art_root, records,
                sources=tuple(source_name for _, source_name, _ in actions),
            )
        except nhan_vat_moc.MocError as error:
            raise ClosureError(str(error)) from error
    cf_by_slot = {
        component: charinfo_slot_to_cf(charinfo, component)
        for component in COMPONENT_CHARINFO_SLOT
    }
    specs = derive_player_slot_specs(review, moc=moc, records=records, idle_frames=idle_frames)
    hinh_hoc, tham_chieu, s0 = player_frame_geometry(records, moc, idle_frames=idle_frames)

    frames: dict[tuple[str, int], dict[str, Any]] = {}

    def frame_at(source_name: str, frame_index: int, scale: float):
        key = (source_name, frame_index)
        if key not in frames:
            alpha = _load_reviewed_alpha(
                art_root, records, source_name, alpha_cache
            )
            x, y, width, height = _frame_rect(records[source_name], frame_index)
            resized = alpha.crop((x, y, x + width, y + height)).resize(
                (max(1, round(width * scale)), max(1, round(height * scale))),
                Image.Resampling.NEAREST,
            )
            bbox = resized.getchannel("A").getbbox()
            if bbox is None:
                raise ClosureError(
                    f"ASSET-CLOSURE-015: khung nguồn rỗng: {source_name}#{frame_index}"
                )
            khung = hinh_hoc[key]
            frames[key] = {
                "resized": resized,
                "rect": (x, y, width, height),
                "center_x4": (bbox[0] + bbox[2]) / 2.0,
                "feet_x4": float(bbox[3]),
                "chin_x4": None if khung["chin_y"] is None else (float(khung["chin_y"]) - y) * scale,
                "face_cx_x4": None
                if khung["face_cx"] is None
                else (float(khung["face_cx"]) + moc[key]["edge_x"] - x) * scale,
                "hair_top_x4": None if khung["hair_top"] is None else (float(khung["hair_top"]) - y) * scale,
            }
        return frames[key]

    pieces = []
    for spec in specs:
        component = spec["component"]
        source_name = spec["source_path"]
        frame_index = spec["frame_index"]
        scale = spec["scale_frame"]
        frame = frame_at(source_name, frame_index, scale)
        resized = frame["resized"]
        frame_x, frame_y = frame["rect"][0], frame["rect"][1]
        rect_x, rect_y, rect_width, rect_height = spec["source_rect"]
        left = max(0, round((rect_x - frame_x) * scale))
        top = max(0, round((rect_y - frame_y) * scale))
        right = min(resized.width, round((rect_x - frame_x + rect_width) * scale))
        bottom = min(resized.height, round((rect_y - frame_y + rect_height) * scale))
        if right <= left or bottom <= top:
            raise ClosureError(
                f"ASSET-CLOSURE-015: vùng mảnh rỗng sau scale: {spec['id']}"
            )
        piece = resized.crop((left, top, right, bottom))
        bbox = piece.getchannel("A").getbbox()
        if bbox is None:
            raise ClosureError(
                f"ASSET-CLOSURE-015: sprite alpha rỗng: {spec['id']}"
            )
        piece = piece.crop(bbox)
        component_left = left + bbox[0]
        component_top = top + bbox[1]
        x1 = zoom_variants(piece)["x1"]
        bbox_x1 = x1.getchannel("A").getbbox()
        if bbox_x1 is None:
            raise ClosureError(f"ASSET-CLOSURE-015: ảnh x1 rỗng: {spec['id']}")

        def off_y(v):
            return None if v is None else float((v - component_top) / 4.0)

        def off_x(v):
            return None if v is None else float((v - component_left) / 4.0)

        cf_list = cf_by_slot[component].get(spec["slot_index"], [])
        if not cf_list:
            raise ClosureError(
                f"ASSET-CLOSURE-015: slot {spec['id']} không cf nào dùng"
            )
        pieces.append(
            {
                **spec,
                "image_x4": piece,
                "moc": {
                    "bbox_x1": [int(v) for v in bbox_x1],
                    "chin_off": off_y(frame["chin_x4"]),
                    "feet_off": off_y(frame["feet_x4"]),
                    "hair_off": off_y(frame["hair_top_x4"]),
                    "face_cx_off": off_x(frame["face_cx_x4"]),
                    "center_off": off_x(frame["center_x4"]),
                    "hip_cx_off": _hip_cx_x1(piece, component),
                },
                "charinfo_cf": cf_list,
                "charinfo_anchor_cf": cf_list[0],
                "_frame": frame,
                "_origin_x4": [int(component_left), int(component_top)],
            }
        )

    try:
        nghiem = nhan_vat_toi_uu.giai(pieces, charinfo, tham_so)
    except nhan_vat_toi_uu.VoNghiem as error:
        raise ClosureError(f"ASSET-CLOSURE-015: {error}") from error
    bang = nhan_vat_toi_uu.bang_cf(pieces, charinfo, nghiem, tham_so)
    vi_pham = nhan_vat_toi_uu.kiem_nghiem(pieces, charinfo, nghiem, tham_so)
    if vi_pham:
        raise ClosureError(f"ASSET-CLOSURE-015: nghiệm MILP không tự kiểm được: {vi_pham[:3]}")
    milp = {
        "phien_ban": PLAYER_GEOMETRY_VERSION,
        "tham_so": tham_so.json(),
        "nghiem": {k: nghiem[k] for k in ("head", "body", "leg")},
        "muc_tieu": nghiem["muc_tieu"],
        "bang_cf": bang,
        "tham_chieu_idle": tham_chieu,
        "idle_frames": list(idle_frames),
    }

    result = []
    for piece in pieces:
        component = piece["component"]
        slot = piece["slot_index"]
        charinfo_index = COMPONENT_CHARINFO_SLOT[component]
        anchor_row = charinfo[piece["charinfo_anchor_cf"]][charinfo_index]
        dx = int(nghiem[component]["dx"][slot])
        dy = int(nghiem[component]["dy"][slot])
        frame = piece.pop("_frame")
        origin = piece.pop("_origin_x4")
        residual = {
            str(row["cf"]): {k: row[k] for k in ("e", "d", "ex", "gx", "khe_dau_than", "khe_than_chan")}
            for row in bang
            if row["cf"] in piece["charinfo_cf"]
        }
        result.append(
            {
                **piece,
                "placement": {
                    "dx": dx,
                    "dy": dy,
                    "left_rel": anchor_row[1] + dx,
                    "top_rel": dy - anchor_row[2],
                },
                "charinfo_residual": residual,
                "milp": milp,
                "shared_transform": {
                    "source_frame": list(frame["rect"]),
                    "scale": s0,
                    "scale_frame": piece["scale_frame"],
                    "frame_center_x4": frame["center_x4"],
                    "frame_feet_x4": frame["feet_x4"],
                    "component_origin_x4": origin,
                    "feet_offset_x1": feet_offset_x1,
                },
            }
        )
    return result


def _thuoc_iruka(khung: dict[str, Any], ten: str) -> float | None:
    """Giá trị một thước trong IRUKA_THUOC_CUNG; `face_area_sqrt` = sqrt(face_area)."""

    if ten == "face_area_sqrt":
        v = khung.get("face_area")
        return None if v is None else float(np.sqrt(float(v)))
    v = khung.get(ten)
    return None if v is None else float(v)


def iruka_he_so_dau(
    khung_iruka: dict[str, Any],
    moc: dict[tuple[str, int], dict[str, Any]],
    *,
    idle_frames: tuple[int, ...] = (0, 1, 2),
    source: str = PLAYER_REFERENCE_SOURCE,
    thuoc: tuple[str, ...] = IRUKA_THUOC_CUNG,
) -> dict[str, Any]:
    """r = cỡ đầu Iruka / cỡ đầu người chơi (median 6 thước, tham chiếu trung bình idle).

    ``khung_iruka`` = một mục của ``nhan_vat_moc.do_dai`` trên master Iruka; ``moc`` = mốc
    17 khung người chơi (chỉ dùng các khung idle). Tham chiếu là TRUNG BÌNH các khung idle
    — cùng tham chiếu với ``player_frame_geometry`` — nên r = 1 nghĩa là đầu hai nhân vật
    bằng nhau trên màn.
    """

    if khung_iruka.get("face_w") is None:
        raise ClosureError("ASSET-CLOSURE-015: khung Iruka không tìm được mặt (da)")
    tham_chieu: dict[str, float] = {}
    ti_le: dict[str, float] = {}
    for ten in thuoc:
        vals = []
        for i in idle_frames:
            key = (source, i)
            if key not in moc:
                raise ClosureError(f"ASSET-CLOSURE-015: thiếu mốc khung idle {key}")
            v = _thuoc_iruka(moc[key], ten)
            if v is not None:
                vals.append(v)
        if not vals:
            raise ClosureError(f"ASSET-CLOSURE-015: idle thiếu thước {ten}")
        tham_chieu[ten] = float(sum(vals) / len(vals))
        v_iruka = _thuoc_iruka(khung_iruka, ten)
        if v_iruka is None or tham_chieu[ten] <= 0:
            continue
        ti_le[ten] = v_iruka / tham_chieu[ten]
    if not ti_le:
        raise ClosureError("ASSET-CLOSURE-015: Iruka không đo được thước nào")
    return {
        "he_so": float(np.median(list(ti_le.values()))),
        "kieu": "trung_binh_idle",
        "idle_frames": list(idle_frames),
        "thuoc": list(thuoc),
        "tham_chieu": tham_chieu,
        "ti_le_tung_thuoc": ti_le,
    }


def iruka_charinfo_residual(
    charinfo: list[list[list[int]]],
    component: str,
    cfs: tuple[int, ...] = IRUKA_CHARINFO_CF,
    anchor_cf: int = 0,
) -> dict[str, list[int]]:
    """Độ lệch (Δleft, Δtop, x1) của mảnh tại từng cf so với cf neo — do CharInfo, không do dx/dy.

    Client: left = cx + CI[cf][k][1] + dx ; top = cy − CI[cf][k][2] + dy. dx/dy chung cho mọi
    cf nên phần lệch giữa cf chỉ còn CI. Với cf 1: đầu/thân [0, −1] (lên 1 px), chân [0, 0].
    """

    k = COMPONENT_CHARINFO_SLOT[component]
    goc = charinfo[anchor_cf][k]
    return {
        str(cf): [int(charinfo[cf][k][1] - goc[1]), int(-(charinfo[cf][k][2] - goc[2]))]
        for cf in cfs
    }


def iruka_bang_cf(
    pieces: list[dict[str, Any]],
    charinfo: list[list[list[int]]],
    cfs: tuple[int, ...] = IRUKA_CHARINFO_CF,
) -> list[dict[str, Any]]:
    """Dự đoán (x1) khe đầu–thân, thân–chân, bàn chân, cao tổng theo phép vẽ client cho từng cf.

    Cùng công thức với ``nhan_vat_toi_uu.bang_cf`` nhưng không cần slot (mọi frame part
    Iruka cùng dx/dy). Dùng để mapping tự chứng minh và để test đối chiếu với ghep-khung.
    """

    by = {p["component"]: p for p in pieces}

    def top_bottom(component, cf):
        k = COMPONENT_CHARINFO_SLOT[component]
        p = by[component]
        bbox = p["moc"]["bbox_x1"]
        top = -charinfo[cf][k][2] + p["placement"]["dy"]
        return top + bbox[1], top + bbox[3]

    rows = []
    for cf in cfs:
        h_t, h_b = top_bottom("head", cf)
        b_t, b_b = top_bottom("body", cf)
        l_t, l_b = top_bottom("leg", cf)
        rows.append(
            {
                "cf": int(cf),
                "khe_dau_than": int(b_t - h_b),
                "khe_than_chan": int(l_t - b_b),
                "chan": int(l_b),
                "cao": int(l_b - h_t),
            }
        )
    return rows


def iruka_component_images(
    review: dict[str, Any],
    *,
    art_root: Path,
    records: dict[str, dict[str, Any]] | None = None,
    alpha_cache: dict[str, Image.Image] | None = None,
    charinfo: list[list[list[int]]] | None = None,
    moc: dict[tuple[str, int], dict[str, Any]] | None = None,
    moc_iruka: dict[tuple[str, int], dict[str, Any]] | None = None,
    feet_offset_x1: int = PLAYER_FEET_OFFSET_X1,
    ti_le_dau: float = IRUKA_TI_LE_DAU,
    idle_frames: tuple[int, ...] = (0, 1, 2),
) -> list[dict[str, Any]]:
    """3 mảnh Iruka (V1) cùng khuôn với ``player_component_images``: đầu, thân, chân.

    Mỗi phần tử: ``id``, ``part_id``, ``component``, ``source_rect``, ``image_x4`` (cắt theo
    cằm, scale s_iruka, trim alpha), ``moc`` (mốc x1 so với góc ảnh), ``placement``
    (``left_rel/top_rel/dx/dy`` tại cf 0), ``charinfo_cf`` [0, 1], ``charinfo_anchor_cf`` 0,
    ``charinfo_residual``, ``bang_cf`` (dự đoán cf 0/1, chung), ``shared_transform``
    (``scale``, ``he_so_dau`` r, ``ti_le_dau``, ``tham_chieu``, gốc khung).

    ``moc`` (17 khung người chơi) / ``moc_iruka`` bỏ trống thì đo thẳng từ ảnh nguồn.
    """

    records = records if records is not None else _review_by_name(review)
    alpha_cache = alpha_cache if alpha_cache is not None else {}
    charinfo = charinfo if charinfo is not None else load_charinfo()
    record = records.get(IRUKA_SOURCE)
    if record is None:
        raise ClosureError(f"ASSET-CLOSURE-015: nguồn Iruka chưa PASS: {IRUKA_SOURCE}")
    if PLAYER_REFERENCE_SOURCE not in records:
        raise ClosureError(f"ASSET-CLOSURE-015: nguồn player chưa PASS: {PLAYER_REFERENCE_SOURCE}")
    try:
        if moc is None:
            moc = nhan_vat_moc.do_dai(art_root, records, PLAYER_REFERENCE_SOURCE)
        if moc_iruka is None:
            moc_iruka = nhan_vat_moc.do_dai(art_root, records, IRUKA_SOURCE)
    except nhan_vat_moc.MocError as error:
        raise ClosureError(str(error)) from error
    key = (IRUKA_SOURCE, 0)
    if key not in moc_iruka:
        raise ClosureError(f"ASSET-CLOSURE-015: thiếu mốc giải phẫu khung {key}")
    khung_iruka = moc_iruka[key]
    if khung_iruka.get("chin_y") is None:
        raise ClosureError("ASSET-CLOSURE-015: Iruka không đo được cằm")

    he_so = iruka_he_so_dau(khung_iruka, moc, idle_frames=idle_frames)
    r = he_so["he_so"]
    if r <= 0:
        raise ClosureError(f"ASSET-CLOSURE-015: hệ số đầu Iruka không hợp lệ: {r}")
    s0 = player_uniform_scale(_frame_rect(records[PLAYER_REFERENCE_SOURCE], PLAYER_REFERENCE_FRAME)[2])
    scale = s0 / r * float(ti_le_dau)
    khung = {"chin_y": float(khung_iruka["chin_y"]), "scale": scale}

    alpha = _load_reviewed_alpha(art_root, records, IRUKA_SOURCE, alpha_cache)
    x, y, width, height = _frame_rect(record, 0)
    resized = alpha.crop((x, y, x + width, y + height)).resize(
        (max(1, round(width * scale)), max(1, round(height * scale))),
        Image.Resampling.NEAREST,
    )
    frame_bbox = resized.getchannel("A").getbbox()
    if frame_bbox is None:
        raise ClosureError(f"ASSET-CLOSURE-015: khung nguồn rỗng: {IRUKA_SOURCE}#0")
    frame_center_x4 = (frame_bbox[0] + frame_bbox[2]) / 2.0
    frame_feet_x4 = float(frame_bbox[3])
    chin_x4 = (float(khung_iruka["chin_y"]) - y) * scale
    face_cx_x4 = (
        None
        if khung_iruka.get("face_cx") is None
        else (float(khung_iruka["face_cx"]) + int(khung_iruka["edge_x"]) - x) * scale
    )
    hair_top_x4 = None if khung_iruka.get("content_top") is None else (float(khung_iruka["content_top"]) - y) * scale

    pieces = []
    for component in ("head", "body", "leg"):
        rect = _component_rect_v1(record, 0, component, khung)
        rect_x, rect_y, rect_width, rect_height = rect
        left = max(0, round((rect_x - x) * scale))
        top = max(0, round((rect_y - y) * scale))
        right = min(resized.width, round((rect_x - x + rect_width) * scale))
        bottom = min(resized.height, round((rect_y - y + rect_height) * scale))
        if right <= left or bottom <= top:
            raise ClosureError(f"ASSET-CLOSURE-015: vùng mảnh Iruka rỗng sau scale: {component}")
        piece = resized.crop((left, top, right, bottom))
        bbox = piece.getchannel("A").getbbox()
        if bbox is None:
            raise ClosureError(f"ASSET-CLOSURE-015: sprite Iruka alpha rỗng: {component}")
        piece = piece.crop(bbox)
        component_left = left + bbox[0]
        component_top = top + bbox[1]
        x1 = zoom_variants(piece)["x1"]
        bbox_x1 = x1.getchannel("A").getbbox()
        if bbox_x1 is None:
            raise ClosureError(f"ASSET-CLOSURE-015: ảnh x1 Iruka rỗng: {component}")
        k = COMPONENT_CHARINFO_SLOT[component]
        placement = component_placement(
            component_left_x4=float(component_left),
            component_top_x4=float(component_top),
            frame_center_x4=frame_center_x4,
            frame_feet_x4=frame_feet_x4,
            charinfo_row=charinfo[0][k],
            feet_offset_x1=feet_offset_x1,
        )

        def off_y(v):
            return None if v is None else float((v - component_top) / 4.0)

        def off_x(v):
            return None if v is None else float((v - component_left) / 4.0)

        pieces.append(
            {
                "id": IRUKA_IDS[("head", "body", "leg").index(component)],
                "part_id": IRUKA_PART_IDS[component],
                "name": f"iruka-{component}",
                "component": component,
                "slot_index": 0,
                "action": "idle",
                "action_slot": 0,
                "frame_index": 0,
                "source_path": IRUKA_SOURCE,
                "source_sha256": record["source_sha256"],
                "source_rect": rect,
                "scale_frame": scale,
                "nguon": {
                    "frame_top": y,
                    "frame_height": height,
                    "edge_x": int(khung_iruka["edge_x"]),
                    "he_so_dau": r,
                    "chin_y": float(khung_iruka["chin_y"]),
                    "cam_uoc_luong": False,
                    "face_cx": None if khung_iruka.get("face_cx") is None else float(khung_iruka["face_cx"]),
                    "hair_top": None if khung_iruka.get("content_top") is None else int(khung_iruka["content_top"]),
                },
                "image_x4": piece,
                "moc": {
                    "bbox_x1": [int(v) for v in bbox_x1],
                    "chin_off": off_y(chin_x4),
                    "feet_off": off_y(frame_feet_x4),
                    "hair_off": off_y(hair_top_x4),
                    "face_cx_off": off_x(face_cx_x4),
                    "center_off": off_x(frame_center_x4),
                    "hip_cx_off": _hip_cx_x1(piece, component),
                },
                "placement": placement,
                "charinfo_cf": list(IRUKA_CHARINFO_CF),
                "charinfo_anchor_cf": 0,
                "charinfo_residual": iruka_charinfo_residual(charinfo, component),
                "anchor": "feet-bottom-center",
                "layer": {"head": 2, "body": 1, "leg": 0}[component],
                "overlap_ratio": {"head_body": 0.14, "body_leg": 0.10},
                "provenance_role": "approved_iruka_master_component",
                "shared_transform": {
                    "source_frame": [x, y, width, height],
                    "scale": scale,
                    "scale_frame": scale,
                    "uniform_scale_x4": s0,
                    "he_so_dau": r,
                    "ti_le_dau": float(ti_le_dau),
                    "tham_chieu": {
                        "kieu": he_so["kieu"],
                        "idle_frames": he_so["idle_frames"],
                        "thuoc": he_so["thuoc"],
                        "gia_tri": he_so["tham_chieu"],
                        "ti_le_tung_thuoc": he_so["ti_le_tung_thuoc"],
                    },
                    "frame_bbox_x4": [int(v) for v in frame_bbox],
                    "frame_center_x4": frame_center_x4,
                    "frame_feet_x4": frame_feet_x4,
                    "component_origin_x4": [int(component_left), int(component_top)],
                    "feet_offset_x1": feet_offset_x1,
                },
            }
        )
    bang = iruka_bang_cf(pieces, charinfo)
    for piece in pieces:
        piece["bang_cf"] = bang
    return pieces


def fit_to_canvas(image: Image.Image, size: tuple[int, int]) -> Image.Image:
    """Fit RGBA theo tỉ lệ, neo giữa-dưới và không cắt alpha."""

    source = image.convert("RGBA")
    bbox = source.getchannel("A").getbbox()
    if bbox is None:
        raise ClosureError("ASSET-CLOSURE-015: crop alpha rỗng")
    source = source.crop(bbox)
    width, height = size
    scale = min(width / source.width, height / source.height)
    target = (
        max(1, round(source.width * scale)),
        max(1, round(source.height * scale)),
    )
    resized = source.resize(target, Image.Resampling.NEAREST)
    canvas = Image.new("RGBA", size)
    canvas.alpha_composite(
        resized,
        ((width - resized.width) // 2, height - resized.height),
    )
    return canvas


def _va_dom_chroma(
    mau: np.ndarray, candidate: np.ndarray, alpha: np.ndarray, dom: list[tuple[int, int]]
) -> np.ndarray:
    """Thay màu từng đốm chroma lẻ bằng trung bình lân cận không-chroma, GIỮ alpha.

    Vá màu chứ không hạ alpha vì đốm lẻ nằm trong lòng hình (nhiễu răng cưa lọt
    vào giữa mảng áo); xoá alpha ở đó là đục một lỗ kim trên áo. Chỉ khi quanh đó
    không có pixel sạch nào trong bán kính 11 mới chịu hạ alpha.
    """

    height, width = candidate.shape
    sach = (~candidate) & (alpha > 0)
    for x, y in dom:
        for ban_kinh in (2, 5, 11):
            y0, y1 = max(0, y - ban_kinh), min(height, y + ban_kinh + 1)
            x0, x1 = max(0, x - ban_kinh), min(width, x + ban_kinh + 1)
            o = sach[y0:y1, x0:x1]
            if o.any():
                mau[y, x] = np.rint(mau[y0:y1, x0:x1][o].mean(axis=0)).astype(mau.dtype)
                break
        else:
            alpha[y, x] = 0
    return mau


def remove_large_chroma_islands(
    source: Image.Image,
    extracted: Image.Image,
    *,
    dien_tich_khe: int | None = None,
    va_dom: bool = False,
) -> Image.Image:
    """Xóa matte magenta nội vùng, giữ accent màu nhỏ của nhân vật.

    ``dien_tich_khe`` bỏ trống = ngưỡng v11 (0,038 % diện tích ảnh, sàn 16 px) —
    đường closure cũ giữ NGUYÊN từng byte. Truyền số nhỏ hơn để xoá cả khe nhỏ;
    ``va_dom`` thì đảo dưới ngưỡng được vá màu thay vì bỏ mặc (xem
    ``_va_dom_chroma``).
    """

    rgb = np.asarray(source.convert("RGB"), dtype=np.int16)
    high = np.maximum(rgb[:, :, 0], rgb[:, :, 2])
    low = np.minimum(rgb[:, :, 0], rgb[:, :, 2])
    candidate = (
        (low - rgb[:, :, 1] >= 70)
        & (high >= 150)
        & (high - rgb[:, :, 1] >= 100)
        & (np.abs(rgb[:, :, 0] - rgb[:, :, 2]) <= 95)
    )
    alpha = np.asarray(extracted.convert("RGBA").getchannel("A")).copy()
    remaining = candidate & (alpha > 0)
    visited = np.zeros(remaining.shape, dtype=bool)
    # Chú thích cũ (v11) tin rằng đảo 586 px trên master Iruka là "accent bông tai"
    # nên chốt ngưỡng 598 px để giữ nó lại. ĐO LẠI 06/09/2026: đảo đó ở
    # x[488,518] y[329,375], màu (245,3,238)/(255,7,255) — magenta khoá màu nguyên
    # chất, chính là KHE giữa dải khăn trán và tóc; không có bông tai nào. Ngưỡng
    # cũ vẫn là mặc định để đường v11 bất động, còn nhân vật thì gọi với
    # ``dien_tich_khe`` nhỏ (xem ``NGUON_NHAN_VAT`` / ``KHE_NEN_TOI_THIEU``).
    threshold = (
        max(16, round(remaining.size * 0.00038))
        if dien_tich_khe is None
        else dien_tich_khe
    )
    height, width = remaining.shape
    dom: list[tuple[int, int]] = []
    for start_y, start_x in zip(*np.where(remaining & ~visited)):
        if visited[start_y, start_x]:
            continue
        component = []
        queue = deque([(int(start_x), int(start_y))])
        visited[start_y, start_x] = True
        while queue:
            x, y = queue.popleft()
            component.append((x, y))
            for nx, ny in ((x - 1, y), (x + 1, y), (x, y - 1), (x, y + 1)):
                if (
                    0 <= nx < width
                    and 0 <= ny < height
                    and remaining[ny, nx]
                    and not visited[ny, nx]
                ):
                    visited[ny, nx] = True
                    queue.append((nx, ny))
        if len(component) >= threshold:
            for x, y in component:
                alpha[y, x] = 0
        elif va_dom:
            dom.extend(component)
    result = extracted.convert("RGBA")
    if dom:
        mau = np.asarray(result.convert("RGB")).copy()
        mau = _va_dom_chroma(mau, candidate, alpha, dom)
        return Image.fromarray(np.dstack([mau, alpha]), mode="RGBA")
    result.putalpha(Image.fromarray(alpha, mode="L"))
    return result


def _save_zooms(
    x4: Image.Image, temporary: Path, relative_stem: Path
) -> tuple[dict[str, str], dict[str, str]]:
    paths = {}
    hashes = {}
    for zoom, image in zoom_variants(x4).items():
        destination = temporary / relative_stem / zoom / "sprite.png"
        destination.parent.mkdir(parents=True, exist_ok=True)
        image.save(destination, format="PNG", optimize=False, compress_level=9)
        repo_relative = CLOSURE_RELATIVE / relative_stem / zoom / "sprite.png"
        paths[zoom] = repo_relative.as_posix()
        hashes[zoom] = _sha(destination)
    return paths, hashes


def _load_reviewed_alpha(
    art_root: Path,
    records: dict[str, dict[str, Any]],
    source_name: str,
    cache: dict[str, Image.Image],
) -> Image.Image:
    if source_name not in records:
        raise ClosureError(
            f"ASSET-CLOSURE-015: source chưa có PASS: {source_name}"
        )
    source = art_root / source_name
    record = records[source_name]
    if not source.is_file() or _sha(source) != record.get("source_sha256"):
        raise ClosureError(
            f"ASSET-CLOSURE-015: SHA source lệch review: {source_name}"
        )
    if source_name not in cache:
        with Image.open(source) as image:
            image.load()
            if source_name in NGUON_NHAN_VAT:
                cache[source_name] = remove_large_chroma_islands(
                    image,
                    extract_alpha(image),
                    dien_tich_khe=KHE_NEN_TOI_THIEU,
                    va_dom=True,
                )
            else:
                cache[source_name] = remove_large_chroma_islands(
                    image, extract_alpha(image)
                )
    return cache[source_name]


def _sprite_record(
    spec: dict[str, Any],
    *,
    art_root: Path,
    records: dict[str, dict[str, Any]],
    alpha_cache: dict[str, Image.Image],
    temporary: Path,
    target_size: tuple[int, int],
) -> dict[str, Any]:
    source_name = spec["source_path"]
    source = art_root / source_name
    alpha = _load_reviewed_alpha(art_root, records, source_name, alpha_cache)
    x, y, width, height = spec["source_rect"]
    if x < 0 or y < 0 or x + width > alpha.width or y + height > alpha.height:
        raise ClosureError(
            f"ASSET-CLOSURE-015: crop ngoài source: {spec['id']}"
        )
    x4 = fit_to_canvas(alpha.crop((x, y, x + width, y + height)), target_size)
    if x4.getchannel("A").getbbox() is None:
        raise ClosureError(f"ASSET-CLOSURE-015: sprite alpha rỗng: {spec['id']}")
    zooms, zoom_hashes = _save_zooms(
        x4, temporary, Path("sprites") / str(spec["id"])
    )
    record = {
        "id": spec["id"],
        "name": spec.get("name", f"{spec.get('component', 'asset')}-{spec['id']}"),
        "crop": {
            "asset_id": spec["id"],
            "source_rect": spec["source_rect"],
        },
        "zooms": zooms,
        "source_path": (ART_PREFIX / source_name).as_posix(),
        "source_sha256": _sha(source),
        "x4_sha256": zoom_hashes["x4"],
        "source_dimensions": list(alpha.size),
        "source_alpha": False,
        "target_x4": list(target_size),
        "zoom_sha256": zoom_hashes,
        "anchor": spec.get("anchor", "bottom-center"),
        "layer": spec.get("layer", 0),
        "overlap_ratio": spec.get("overlap_ratio", {}),
    }
    for key in (
        "component",
        "slot_index",
        "action",
        "action_slot",
        "client_cf",
        "provenance_role",
    ):
        if key in spec:
            record[key] = spec[key]
    return record


def player_sprite_record(
    component: dict[str, Any],
    *,
    art_root: Path,
    records: dict[str, dict[str, Any]],
    temporary: Path,
) -> dict[str, Any]:
    """Record manifest cho một mảnh người chơi đã đặt trong hệ toạ độ chung."""

    source_name = component["source_path"]
    source = art_root / source_name
    x4 = component["image_x4"]
    with Image.open(source) as image:
        source_dimensions = list(image.size)
    zooms, zoom_hashes = _save_zooms(
        x4, temporary, Path("sprites") / str(component["id"])
    )
    record = {
        "id": component["id"],
        "name": f"{component['component']}-{component['id']}",
        "crop": {
            "asset_id": component["id"],
            "source_rect": component["source_rect"],
        },
        "zooms": zooms,
        "source_path": (ART_PREFIX / source_name).as_posix(),
        "source_sha256": _sha(source),
        "x4_sha256": zoom_hashes["x4"],
        "source_dimensions": source_dimensions,
        "source_alpha": False,
        "target_x4": list(x4.size),
        "zoom_sha256": zoom_hashes,
        # Gốc toạ độ là bàn chân của CẢ KHUNG, không phải đáy canvas.
        "anchor": "feet-bottom-center",
        "layer": component["layer"],
        "overlap_ratio": component["overlap_ratio"],
        "component": component["component"],
        "slot_index": component["slot_index"],
        "action": component["action"],
        "action_slot": component["action_slot"],
        "client_cf": component["client_cf"],
        "client_cf_charinfo": component["charinfo_cf"],
        "charinfo_anchor_cf": component["charinfo_anchor_cf"],
        "charinfo_residual": component["charinfo_residual"],
        "part_offset": {
            "dx": component["placement"]["dx"],
            "dy": component["placement"]["dy"],
        },
        "placement_x1": {
            "left_rel": component["placement"]["left_rel"],
            "top_rel": component["placement"]["top_rel"],
        },
        "shared_transform": component["shared_transform"],
    }
    return record


def _single_frame_spec(
    records: dict[str, dict[str, Any]],
    *,
    sprite_id: int,
    source_name: str,
    frame_index: int,
    name: str,
    provenance_role: str,
) -> dict[str, Any]:
    record = records.get(source_name)
    if record is None:
        raise ClosureError(f"ASSET-CLOSURE-015: source chưa PASS: {source_name}")
    return {
        "id": sprite_id,
        "name": name,
        "source_path": source_name,
        "source_sha256": record["source_sha256"],
        "source_rect": _frame_rect(record, frame_index),
        "provenance_role": provenance_role,
    }


def _generate_iruka_components(
    *,
    record: dict[str, Any],
    art_root: Path,
    alpha_cache: dict[str, Image.Image],
    temporary: Path,
) -> list[dict[str, Any]]:
    """Giữ ba part trong cùng hệ tọa độ để runtime ghép lại đúng một Iruka.

    HÌNH HỌC CŨ (closure v11): ép cả khung vào canvas 128×160 (Iruka 38 px x1 = 62 %
    người chơi), lát cắt theo % chiều cao, không có dx/dy. Chỉ còn phục vụ đường
    ``_generate_sprites`` → closure v11 mà test cũ và ``sprites[].zooms`` trong mapping
    vẫn trỏ tới. Ảnh Iruka ĐANG DÙNG (client) sinh bằng ``iruka_component_images`` qua
    ``tools/nhan-vat/sinh-lai-sprite.py``.
    """

    source_name = IRUKA_SOURCE
    source = art_root / source_name
    alpha = _load_reviewed_alpha(
        art_root, {source_name: record}, source_name, alpha_cache
    )
    frame_x, frame_y, frame_width, frame_height = _frame_rect(record, 0)
    frame = alpha.crop(
        (
            frame_x,
            frame_y,
            frame_x + frame_width,
            frame_y + frame_height,
        )
    )
    canvas_size = (128, 160)
    content_limit = (120, 152)
    scale = min(
        content_limit[0] / frame_width,
        content_limit[1] / frame_height,
    )
    resized_size = (
        max(1, round(frame_width * scale)),
        max(1, round(frame_height * scale)),
    )
    offset = (
        (canvas_size[0] - resized_size[0]) // 2,
        canvas_size[1] - 4 - resized_size[1],
    )
    result = []
    for offset_id, component in enumerate(("head", "body", "leg")):
        component_x, component_y, component_width, component_height = (
            _component_rect(record, 0, component)
        )
        layer = Image.new("RGBA", frame.size)
        local_box = (
            component_x - frame_x,
            component_y - frame_y,
            component_x - frame_x + component_width,
            component_y - frame_y + component_height,
        )
        layer.alpha_composite(frame.crop(local_box), local_box[:2])
        resized = layer.resize(resized_size, Image.Resampling.NEAREST)
        canvas = Image.new("RGBA", canvas_size)
        canvas.alpha_composite(resized, offset)
        sprite_id = 27075 + offset_id
        zooms, zoom_hashes = _save_zooms(
            canvas, temporary, Path("sprites") / str(sprite_id)
        )
        result.append(
            {
                "id": sprite_id,
                "name": f"iruka-{component}",
                "crop": {
                    "asset_id": sprite_id,
                    "source_rect": [
                        component_x,
                        component_y,
                        component_width,
                        component_height,
                    ],
                },
                "zooms": zooms,
                "source_path": (ART_PREFIX / source_name).as_posix(),
                "source_sha256": _sha(source),
                "x4_sha256": zoom_hashes["x4"],
                "source_dimensions": list(alpha.size),
                "source_alpha": False,
                "target_x4": list(canvas_size),
                "zoom_sha256": zoom_hashes,
                "anchor": "bottom-center",
                "layer": {"head": 2, "body": 1, "leg": 0}[component],
                "overlap_ratio": {"head_body": 0.14, "body_leg": 0.10},
                "component": component,
                "slot_index": 0,
                "action": "idle",
                "action_slot": 0,
                "provenance_role": "approved_iruka_master_component",
                "shared_transform": {
                    "source_frame": [
                        frame_x,
                        frame_y,
                        frame_width,
                        frame_height,
                    ],
                    "resized_size": list(resized_size),
                    "canvas_offset": list(offset),
                },
            }
        )
    return result


def _generate_sprites(
    review: dict[str, Any], art_root: Path, temporary: Path
) -> list[dict[str, Any]]:
    records = _review_by_name(review)
    cache: dict[str, Image.Image] = {}
    sprites = []
    for component in player_component_images(
        review, art_root=art_root, records=records, alpha_cache=cache
    ):
        sprites.append(
            player_sprite_record(
                component,
                art_root=art_root,
                records=records,
                temporary=temporary,
            )
        )

    sprites.extend(
        _generate_iruka_components(
            record=records["master_iruka_x4_chroma.png"],
            art_root=art_root,
            alpha_cache=cache,
            temporary=temporary,
        )
    )

    portrait = _single_frame_spec(
        records,
        sprite_id=27078,
        source_name="iruka_portrait_chroma.png",
        frame_index=0,
        name="iruka-portrait",
        provenance_role="approved_iruka_portrait",
    )
    sprites.append(
        _sprite_record(
            portrait,
            art_root=art_root,
            records=records,
            alpha_cache=cache,
            temporary=temporary,
            target_size=(200, 200),
        )
    )

    skill_record = records["skill_icons_leaf_4_chroma.png"]
    for slot in range(4):
        spec = _single_frame_spec(
            records,
            sprite_id=27079 + slot,
            source_name="skill_icons_leaf_4_chroma.png",
            frame_index=slot,
            name=f"skill-leaf-{slot}",
            provenance_role="approved_skill_icon",
        )
        spec.update({"slot_index": slot, "action": "skill_pose"})
        sprites.append(
            _sprite_record(
                spec,
                art_root=art_root,
                records=records,
                alpha_cache=cache,
                temporary=temporary,
                target_size=(64, 64),
            )
        )
    if [sprite["id"] for sprite in sprites] != list(range(27041, 27083)):
        raise ClosureError("ASSET-CLOSURE-015: sprite ID closure không liên tục")
    return sprites


def _generate_wolf_cells(
    review: dict[str, Any], art_root: Path, temporary: Path
) -> list[str]:
    records = _review_by_name(review)
    cache: dict[str, Image.Image] = {}
    paths = []
    slot = 0
    for action, source_name, count in WOLF_ACTIONS:
        record = records.get(source_name)
        if record is None or len(
            record.get("checks", {}).get("crop", {}).get("content_boxes", [])
        ) != count:
            raise ClosureError(
                f"ASSET-CLOSURE-015: wolf frame closure lệch: {source_name}"
            )
        alpha = _load_reviewed_alpha(
            art_root, records, source_name, cache
        )
        for frame_index in range(count):
            x, y, width, height = _frame_rect(record, frame_index)
            cell = fit_to_canvas(
                alpha.crop((x, y, x + width, y + height)), (96, 96)
            )
            destination = (
                temporary / "animation" / "quai-soi-rung" / f"{slot:02d}.png"
            )
            destination.parent.mkdir(parents=True, exist_ok=True)
            cell.save(destination, format="PNG", compress_level=9)
            paths.append(
                (
                    CLOSURE_RELATIVE
                    / "animation"
                    / "quai-soi-rung"
                    / f"{slot:02d}.png"
                ).as_posix()
            )
            slot += 1
    return paths


def _animation_closure(
    sprites: list[dict[str, Any]], wolf_paths: list[str]
) -> tuple[list[dict[str, Any]], list[dict[str, Any]]]:
    by_id = {sprite["id"]: sprite for sprite in sprites}
    ranges = {
        "hoc-vien-head": list(range(27041, 27044)),
        "hoc-vien-body": list(range(27044, 27061)),
        "hoc-vien-leg": list(range(27061, 27075)),
    }
    sheets = []
    animations = []
    for name, required in REQUIRED_ANIMATION_SHEETS.items():
        if name == "quai-soi-rung":
            cell_paths = wolf_paths
        else:
            cell_paths = [by_id[sprite_id]["zooms"]["x4"] for sprite_id in ranges[name]]
        sheets.append(
            {
                "name": name,
                "slot_count": required["slot_count"],
                "cell_size": list(required["cell_size"]),
                "cells": [
                    {"slot_index": index, "path": path}
                    for index, path in enumerate(cell_paths)
                ],
            }
        )
        animations.append(
            {
                "name": name,
                "slot_count": required["slot_count"],
                "frames": [
                    {
                        "slot_index": index,
                        "client_cf": list(REQUIRED_CLIENT_CF[name][index]),
                    }
                    for index in range(required["slot_count"])
                ],
            }
        )
    return sheets, animations


def _generate_mob_atlases(
    shared_root: Path, temporary: Path, wolf_paths: list[str]
) -> dict[str, str]:
    # EffectData dùng toạ độ logical x1; atlas zoom lớn giữ nguyên layout nhân zoom.
    x4_cells = []
    for relative in wolf_paths:
        with Image.open(shared_root / relative) as image:
            image.load()
            x4_cells.append(image.convert("RGBA"))
    paths = {}
    for zoom_index, zoom in enumerate(ZOOMS, start=1):
        cell_size = 24 * zoom_index
        atlas = Image.new("RGBA", (cell_size * 10, cell_size))
        for slot, x4 in enumerate(x4_cells):
            cell = (
                x4
                if zoom == "x4"
                else x4.resize((cell_size, cell_size), Image.Resampling.NEAREST)
            )
            atlas.alpha_composite(cell, (slot * cell_size, 0))
        destination = temporary / "mob" / zoom / "img.png"
        destination.parent.mkdir(parents=True, exist_ok=True)
        atlas.save(destination, format="PNG", compress_level=9)
        paths[zoom] = (CLOSURE_RELATIVE / "mob" / zoom / "img.png").as_posix()
    return paths


def _largest_interior_alpha_component(image: Image.Image) -> Image.Image:
    """Crop đúng thanh HUD, bỏ đường viền artboard chạm mép quadrant."""

    source = image.convert("RGBA")
    alpha = np.asarray(source.getchannel("A"))
    visible = alpha > 0
    height, width = visible.shape
    visited = np.zeros_like(visible, dtype=bool)
    components: list[list[tuple[int, int]]] = []

    for start_y, start_x in zip(*np.where(visible & ~visited)):
        if visited[start_y, start_x]:
            continue
        component: list[tuple[int, int]] = []
        queue: deque[tuple[int, int]] = deque([(int(start_x), int(start_y))])
        visited[start_y, start_x] = True
        touches_edge = False
        while queue:
            x, y = queue.popleft()
            component.append((x, y))
            touches_edge = (
                touches_edge
                or x in (0, width - 1)
                or y in (0, height - 1)
            )
            for next_x, next_y in (
                (x - 1, y),
                (x + 1, y),
                (x, y - 1),
                (x, y + 1),
            ):
                if not (0 <= next_x < width and 0 <= next_y < height):
                    continue
                if not visible[next_y, next_x] or visited[next_y, next_x]:
                    continue
                visited[next_y, next_x] = True
                queue.append((next_x, next_y))
        if not touches_edge:
            components.append(component)

    if not components:
        raise ClosureError(
            "ASSET-CLOSURE-015: HUD không có component nội vùng"
        )

    component = max(components, key=len)
    xs = [point[0] for point in component]
    ys = [point[1] for point in component]
    bbox = (min(xs), min(ys), max(xs) + 1, max(ys) + 1)
    selected = np.zeros_like(alpha, dtype=np.uint8)
    for x, y in component:
        selected[y, x] = alpha[y, x]
    result = source.copy()
    result.putalpha(Image.fromarray(selected, mode="L"))
    return result.crop(bbox)


def _extract_hud_bar(quadrant: Image.Image) -> Image.Image:
    extracted = extract_alpha(quadrant)
    alpha = np.asarray(extracted.getchannel("A")).copy()
    alpha[_adaptive_chroma_mask(quadrant.convert("RGB"))] = 0
    extracted.putalpha(Image.fromarray(alpha, mode="L"))
    return _largest_interior_alpha_component(extracted)


def _generate_hud(
    art_root: Path,
    review: dict[str, Any],
    temporary: Path,
) -> dict[str, dict[str, str]]:
    records = _review_by_name(review)
    alpha = _load_reviewed_alpha(
        art_root, records, "hud_hp_chakra_4_chroma.png", {}
    )
    middle_x, middle_y = alpha.width // 2, alpha.height // 2
    boxes = {
        "naruto-hp-bg": (0, 0, middle_x, middle_y),
        "naruto-hp-fill": (middle_x, 0, alpha.width, middle_y),
        "naruto-chakra-bg": (0, middle_y, middle_x, alpha.height),
        "naruto-chakra-fill": (middle_x, middle_y, alpha.width, alpha.height),
    }
    result = {}
    for name, box in boxes.items():
        quadrant = alpha.crop(box)
        # Bỏ viền trắng phân ô của master trước khi trim alpha nội dung thanh.
        quadrant = quadrant.crop(
            (8, 8, quadrant.width - 8, quadrant.height - 8)
        )
        pixels = np.asarray(quadrant.convert("RGBA")).copy()
        red = pixels[:, :, 0].astype(np.int16)
        green = pixels[:, :, 1].astype(np.int16)
        blue = pixels[:, :, 2].astype(np.int16)
        # HUD dùng gradient magenta rộng hơn chroma chuẩn. Chỉ ở bốn quadrant
        # này, xóa pixel có đồng thời R/B cao và G thấp; red/blue fill vẫn giữ.
        matte = (
            (red >= 120)
            & (blue >= 120)
            & (green + 45 < np.minimum(red, blue))
        )
        pixels[:, :, 3][matte] = 0
        quadrant = Image.fromarray(pixels, mode="RGBA")
        content = _extract_hud_bar(quadrant)
        x4 = content.resize((240, 32), Image.Resampling.NEAREST)
        zooms, _ = _save_zooms(x4, temporary, Path("hud") / name)
        result[name] = zooms
    return result


def _generate_ui(
    shared_root: Path, temporary: Path
) -> dict[str, dict[str, str]]:
    approved = shared_root / APPROVED_UI_RELATIVE
    result = {}
    for name in UI_ASSETS:
        source = approved / name
        if not source.is_file():
            raise ClosureError(f"ASSET-CLOSURE-015: thiếu UI đã duyệt: {source}")
        with Image.open(source) as image:
            image.load()
            x4 = image.convert("RGBA")
        zooms, _ = _save_zooms(
            x4, temporary, Path("ui") / Path(name).stem
        )
        result[name] = zooms
    return result


def _copy_file(source: Path, temporary: Path, relative: Path) -> str:
    if not source.is_file():
        raise ClosureError(f"ASSET-CLOSURE-015: thiếu input: {source}")
    destination = temporary / relative
    destination.parent.mkdir(parents=True, exist_ok=True)
    shutil.copyfile(source, destination)
    return (CLOSURE_RELATIVE / relative).as_posix()


def _generate_production_evidence(
    temporary: Path,
    hud_sources: dict[str, dict[str, str]],
    sprites: list[dict[str, Any]],
) -> dict[str, dict[str, str]]:
    """Render đúng artifact production cần reviewer nhìn, ngay trong closure."""

    def load(relative: str) -> Image.Image:
        local = temporary / Path(relative).relative_to(CLOSURE_RELATIVE)
        with Image.open(local) as image:
            image.load()
            return image.convert("RGBA")

    hud = {name: load(sources["x4"]) for name, sources in hud_sources.items()}
    hud_path = temporary / "evidence/hud-golden-25-50-100-x4.png"
    render_hud_golden(hud, hud_path)

    by_id = {sprite["id"]: sprite for sprite in sprites}
    iruka_layers = {
        "head": load(by_id[27075]["zooms"]["x4"]),
        "body": load(by_id[27076]["zooms"]["x4"]),
        "leg": load(by_id[27077]["zooms"]["x4"]),
    }
    iruka_path = temporary / "evidence/iruka-composite-light-dark-x4.png"
    render_iruka_composite_contact_sheet(iruka_layers, iruka_path)

    return {
        "hud_golden": {
            "path": (
                CLOSURE_RELATIVE / "evidence/hud-golden-25-50-100-x4.png"
            ).as_posix(),
            "sha256": _sha(hud_path),
        },
        "iruka_composite": {
            "path": (
                CLOSURE_RELATIVE
                / "evidence/iruka-composite-light-dark-x4.png"
            ).as_posix(),
            "sha256": _sha(iruka_path),
        },
    }


def _copy_runtime_inputs(
    project_root: Path, shared_root: Path, temporary: Path
) -> tuple[str, str, str, dict[str, str]]:
    part = _copy_file(
        project_root / "server/internal/game/part.bin",
        temporary,
        Path("inputs/part.bin"),
    )
    tileset = _copy_file(
        project_root / "server/internal/game/tileset.bin",
        temporary,
        Path("inputs/tileset.bin"),
    )
    runtime_map = _copy_file(
        project_root / "docs/my-thuat/unity-runtime-key-map.json",
        temporary,
        Path("inputs/unity-runtime-key-map.json"),
    )
    snapshots = {}
    payload = json.loads(
        (project_root / "docs/my-thuat/unity-runtime-key-map.json").read_text(
            encoding="utf-8"
        )
    )
    sites = {
        site["file"]: site["file_sha256"]
        for entry in payload["entries"]
        for site in [*entry["call_sites"], *entry["consumer_sites"]]
    }
    for site, expected_sha in sorted(sites.items()):
        source = Path("D:/NRO_Project") / site
        if not source.is_file() or _sha(source) != expected_sha:
            raise ClosureError(
                f"ASSET-CLOSURE-015: runtime site lệch precondition: {site}"
            )
        relative = Path("runtime-snapshots") / Path(site)
        snapshots[site] = _copy_file(source, temporary, relative)
    return part, tileset, runtime_map, snapshots


def _background_sources(shared_root: Path) -> dict[str, dict[str, str]]:
    bases = {
        "b00": "map_nen_troi_1",
        "b01": "map_nen_may_1",
        "b02": "map_nen_may_2",
        "b03": "map_nen_nuihokage",
    }
    result = {}
    for layer, base in bases.items():
        result[layer] = {}
        for zoom in ZOOMS:
            relative = Path("docs/my-thuat/anh-game") / zoom / f"{base}_{zoom}.png"
            if not (shared_root / relative).is_file():
                raise ClosureError(
                    f"ASSET-CLOSURE-015: thiếu background: {relative}"
                )
            result[layer][zoom] = relative.as_posix()
    return result


def _tile_art(
    art_root: Path, temporary: Path
) -> tuple[list[dict[str, Any]], list[dict[str, Any]]]:
    source = art_root / "generated/naruto-village-tile-sheet-10x10-v2.png"
    index = generate_tile_sheet(
        source,
        temporary / "tile-art",
        expected_sha256=TILE_SOURCE_SHA256,
    )
    rows = []
    hashes = []
    for tile in index["tiles"]:
        sources = {
            zoom: (
                CLOSURE_RELATIVE / "tile-art" / tile["zooms"][zoom]["path"]
            ).as_posix()
            for zoom in ZOOMS
        }
        rows.append({"index": tile["index"], "solid": False, "sources": sources})
        hashes.append(
            {
                "index": tile["index"],
                **{
                    zoom: tile["zooms"][zoom]["sha256"]
                    for zoom in ZOOMS
                },
            }
        )
    return rows, hashes


def _parts(sprites: list[dict[str, Any]]) -> list[dict[str, int]]:
    """dx/dy lấy từ chính hình học đã sinh, không còn là hằng số chép tay.

    Trước đây hàm này trả dy cố định 0 rồi có người sửa tay trong mapping
    (dy -41/-50/-33) — nghĩa là lần sinh closure kế tiếp sẽ âm thầm lùi lại số 0.
    Giờ mỗi frame lấy đúng ``part_offset`` mà ``player_component_images`` tính.
    """

    by_id = {sprite["id"]: sprite for sprite in sprites}

    def frames(ids: list[int]) -> list[dict[str, int]]:
        rows = []
        for sprite_id in ids:
            offset = by_id.get(sprite_id, {}).get("part_offset", {"dx": 0, "dy": 0})
            rows.append(
                {
                    "sprite_id": sprite_id,
                    "dx": offset["dx"],
                    "dy": offset["dy"],
                }
            )
        return rows

    return [
        {"id": 30, "type": 0, "frames": frames(list(range(27041, 27044)))},
        {"id": 57, "type": 1, "frames": frames(list(range(27044, 27061)))},
        {"id": 58, "type": 2, "frames": frames(list(range(27061, 27075)))},
        {"id": 2189, "type": 0, "frames": frames([27075] * 3)},
        {"id": 2190, "type": 1, "frames": frames([27076] * 17)},
        {"id": 2191, "type": 2, "frames": frames([27077] * 14)},
    ]


def _skills() -> list[dict[str, Any]]:
    return [
        {
            "id": 27079 + slot,
            "pose_slot": slot,
            "levels": list(range(12)),
            "effect_id": -1,
            "arrow_id": -1,
        }
        for slot in range(4)
    ]


def _catalog() -> dict[str, Any]:
    return {
        "maps": [
            "Làng Lá",
            "Làng Cát",
            "Làng Sương Mù",
            "Làng Đá",
            "Làng Mây",
            "Khu rừng",
        ],
        "npc_templates": ["Iruka"],
        "mob_templates": ["Sói Rừng"],
        "spawn": [[100, 384]],
        "skill_poses": [0, 1, 2, 3],
    }


def _ui_bindings(
    project_root: Path,
    ui_sources: dict[str, dict[str, str]],
) -> list[dict[str, Any]]:
    runtime = json.loads(
        (project_root / "docs/my-thuat/unity-runtime-key-map.json").read_text(
            encoding="utf-8"
        )
    )
    bindings = []
    for entry in runtime["entries"]:
        name = entry["source_asset"]
        binding = dict(entry)
        binding["source_path"] = (
            APPROVED_UI_RELATIVE / name
        ).as_posix()
        binding["source_sha256"] = _sha(
            project_root.parents[1] / APPROVED_UI_RELATIVE / name
            if project_root.parent.name == ".worktrees"
            else project_root / APPROVED_UI_RELATIVE / name
        )
        binding["sources"] = ui_sources[name]
        bindings.append(binding)
    return bindings


def generate_production_closure(
    *,
    project_root: Path,
    shared_root: Path,
    review_path: Path,
    mapping_path: Path,
) -> dict[str, Any]:
    """Sinh closure atomically và ghi mapping canonical trong worktree."""

    project_root = project_root.resolve()
    shared_root = shared_root.resolve()
    art_root = shared_root / ART_PREFIX
    output = shared_root / CLOSURE_RELATIVE
    if output.exists():
        raise ClosureError(
            f"ASSET-CLOSURE-015: closure đã tồn tại, không ghi đè: {output}"
        )
    review = json.loads(review_path.read_text(encoding="utf-8"))
    output.parent.mkdir(parents=True, exist_ok=True)
    temporary = Path(
        tempfile.mkdtemp(prefix=".production-closure-v9-", dir=output.parent)
    )
    try:
        sprites = _generate_sprites(review, art_root, temporary)
        wolf_paths = _generate_wolf_cells(review, art_root, temporary)
        # wolf_paths trỏ tới tên canonical cuối; lúc compose dùng file trong temp.
        local_wolf_paths = [
            (temporary / Path(path).relative_to(CLOSURE_RELATIVE))
            for path in wolf_paths
        ]
        mob_paths = {}
        for zoom_index, zoom in enumerate(ZOOMS, start=1):
            cell_size = 24 * zoom_index
            atlas = Image.new("RGBA", (cell_size * 10, cell_size))
            for slot, path in enumerate(local_wolf_paths):
                with Image.open(path) as image:
                    image.load()
                    x4 = image.convert("RGBA")
                cell = (
                    x4
                    if zoom == "x4"
                    else x4.resize(
                        (cell_size, cell_size), Image.Resampling.NEAREST
                    )
                )
                atlas.alpha_composite(cell, (slot * cell_size, 0))
            destination = temporary / "mob" / zoom / "img.png"
            destination.parent.mkdir(parents=True, exist_ok=True)
            atlas.save(destination, format="PNG", compress_level=9)
            mob_paths[zoom] = (
                CLOSURE_RELATIVE / "mob" / zoom / "img.png"
            ).as_posix()

        hud_sources = _generate_hud(art_root, review, temporary)
        artifact_evidence = _generate_production_evidence(
            temporary, hud_sources, sprites
        )
        ui_sources = _generate_ui(shared_root, temporary)
        tile_art, tile_hashes = _tile_art(art_root, temporary)
        part_path, tileset_path, runtime_map, snapshots = _copy_runtime_inputs(
            project_root, shared_root, temporary
        )
        sheets, animations = _animation_closure(sprites, wolf_paths)
        ui_bindings = _ui_bindings(project_root, ui_sources)
        tileset_local = temporary / Path(tileset_path).relative_to(CLOSURE_RELATIVE)
        tileset_bytes = tileset_local.read_bytes()
        mob_images = [
            {"id": index, "x": index * 24, "y": 0, "w": 24, "h": 24}
            for index in range(10)
        ]
        mob_frames_flat = [[[-12, -24, index]] for index in range(10)]
        mob_frames = {
            "idle": mob_frames_flat[:2],
            "run": mob_frames_flat[2:6],
            "pounce": mob_frames_flat[6:8],
            "hit": mob_frames_flat[8:9],
            "death": mob_frames_flat[9:],
        }
        manifest = {
            "schema_version": 1,
            "source_catalog_version": 2,
            "versions": {
                "runtime_slice": 1,
                "data": 1,
                "map": 1,
                "resource": 1,
            },
            "phase": {
                "name": "lang-la",
                "sprite_id_range": [27041, 27082],
                "reserved_sprite_id_range": [27041, 27999],
                "forbidden_poc_ids": [9000, 9001, 9002],
            },
            "sprites": sprites,
            "parts": _parts(sprites),
            "animations": animations,
            "maps": [{"id": 1, "tile_id": 1, "background_id": 0}],
            "tiles": [
                {
                    "id": index,
                    "semantic": "walkable" if index == 0 else "decorative",
                    "solid": False,
                }
                for index in range(100)
            ],
            "backgrounds": [{"id": 0, "layers": 4}],
            "mobs": [{"id": 0, "name": "soi-rung"}],
            "skills": _skills(),
            "catalog": _catalog(),
            "ui_bindings": ui_bindings,
            "build_inputs": {
                "source_base": ".",
                "source_provenance": [
                    {
                        "source_path": (ART_PREFIX / record["asset"]).as_posix(),
                        "source_sha256": record["source_sha256"],
                        "role": "reviewed_vertical_slice_source",
                    }
                    for record in sorted(
                        review["assets"], key=lambda item: item["asset"]
                    )
                ],
                "legacy_part_path": part_path,
                "legacy_part_sha256": _sha(
                    temporary / Path(part_path).relative_to(CLOSURE_RELATIVE)
                ),
                "tile_set_info_path": tileset_path,
                "tile_set_info_sha256": _sha(tileset_local),
                "tile_set_baseline": {
                    "sha256": _sha(tileset_local),
                    "bytes_hex": tileset_bytes.hex(),
                    "semantic": {
                        "byte_length": len(tileset_bytes),
                        "tile_count": tileset_bytes[0],
                    },
                },
                "tile_art": tile_art,
                "tile_art_sha256": tile_hashes,
                "background_sources": _background_sources(shared_root),
                "hud_sources": hud_sources,
                "artifact_evidence": artifact_evidence,
                "mob": {
                    "id": 0,
                    "payload_mode": "generated",
                    "images": mob_images,
                    "frames": mob_frames,
                    "arr_frames": [],
                    "trailer": [0, 0],
                    "img": mob_paths,
                },
                "map": {
                    "width": 5,
                    "height": 17,
                    "indexes": [0] * 85,
                    "spawn": [100, 384],
                    "tile_size": [24, 24],
                    "client_solid": [False] * 85,
                },
                "animation_sheets": sheets,
                "runtime_key_map": runtime_map,
                "runtime_key_map_sha256": _sha(
                    temporary / Path(runtime_map).relative_to(CLOSURE_RELATIVE)
                ),
                "runtime_source_snapshots": snapshots,
            },
        }
        # Tất cả path canonical lúc này trỏ tới output cuối; hash file trong temp.
        closure_hashes = {}
        for relative in _source_paths(manifest):
            path = shared_root / relative
            if relative.startswith(CLOSURE_RELATIVE.as_posix() + "/"):
                path = temporary / Path(relative).relative_to(CLOSURE_RELATIVE)
            if not path.is_file():
                raise ClosureError(
                    f"ASSET-CLOSURE-015: input closure thiếu: {relative}"
                )
            closure_hashes[relative] = _sha(path)
        manifest["build_inputs"]["production_closure_sha256"] = closure_hashes

        index = {
            "schema_version": 1,
            "generator": "tools.artqa.task6b_closure",
            "review_sha256": _sha(review_path),
            "mapping_sha256": hashlib.sha256(
                json.dumps(
                    manifest,
                    ensure_ascii=False,
                    sort_keys=True,
                    separators=(",", ":"),
                ).encode("utf-8")
            ).hexdigest(),
            "source_count": len(closure_hashes),
            "sprite_ids": [sprite["id"] for sprite in sprites],
        }
        (temporary / "generated-index.json").write_text(
            json.dumps(index, ensure_ascii=False, indent=2) + "\n",
            encoding="utf-8",
        )
        temporary.replace(output)
        mapping_path.write_text(
            json.dumps(manifest, ensure_ascii=False, indent=2) + "\n",
            encoding="utf-8",
        )
        return manifest
    except Exception:
        shutil.rmtree(temporary, ignore_errors=True)
        raise


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--project-root", required=True, type=Path)
    parser.add_argument("--shared-root", required=True, type=Path)
    parser.add_argument("--review", required=True, type=Path)
    parser.add_argument("--mapping", required=True, type=Path)
    args = parser.parse_args()
    generate_production_closure(
        project_root=args.project_root,
        shared_root=args.shared_root,
        review_path=args.review,
        mapping_path=args.mapping,
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
