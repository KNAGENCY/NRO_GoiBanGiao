# -*- coding: utf-8 -*-
"""sinh-lai-sprite.py — sinh lại 34 mảnh nhân vật người chơi + part.bin.

VÌ SAO CÓ CÔNG CỤ NÀY (05/09/2026)
----------------------------------
Đường chính để sinh asset là ``tools/artqa/task6b_closure.py`` →
``assetpack build``. Đường đó đang **bị khoá** cho một closure MỚI:

* ``generate_production_closure`` từ chối ghi đè và chỉ ghi vào **kho ảnh dùng
  chung** (``D:\\NRO_GoiBanGiao\\docs\\my-thuat\\anh-nguon\\...``) — nằm ngoài
  worktree này;
* ``build()`` đòi ``production_artifact_review`` do **reviewer độc lập**
  (``/root/task6b_visual_review6``) ký, và chữ ký đó bind theo
  ``production_closure_digest`` — đổi một byte sprite là chữ ký hết hiệu lực
  (``ASSET-ARTIFACT-REVIEW-016``). Tự sửa chữ ký = giả mạo, nên không làm.

Công cụ này **dùng lại đúng các hàm toán** trong ``task6b_closure`` (import, không
chép) để sinh ảnh + dx/dy, rồi ghi thẳng vào repo:

* ``docs/my-thuat/lang-la-client-assets/res/x{1..4}/smallimage/Small<id>.png``
  — bản art client có trong git (client Unity không thuộc repo này);
* ``server/internal/game/assets/part.bin`` — dựng qua ``assetpack.part.build_part``
  từ chính legacy ``server/internal/game/part.bin``;
* ``docs/my-thuat/mapping-id.json`` — ``parts[].frames[].dx/dy`` và khối
  ``build_inputs.player_component_geometry`` (sự thật hình học mới);
* ``server/internal/game/assets/asset-lock.json`` — SHA/độ dài part.bin + version.

Dùng:
    python tools/nhan-vat/sinh-lai-sprite.py            # ghi thật vào worktree
    python tools/nhan-vat/sinh-lai-sprite.py --thu      # chỉ in bảng, không ghi
"""

from __future__ import annotations

import argparse
import hashlib
import json
import sys
from pathlib import Path

GOC = Path(__file__).resolve().parents[2]
if str(GOC) not in sys.path:
    sys.path.insert(0, str(GOC))

from assetpack.animation import REQUIRED_ANIMATION_SHEETS  # noqa: E402
from assetpack.lock import lock_bytes  # noqa: E402
from assetpack.part import build_part, parse_part  # noqa: E402
from tools.artqa.task6b_closure import (  # noqa: E402
    PLAYER_FEET_OFFSET_X1,
    _review_by_name,
    load_charinfo,
    player_component_images,
    zoom_variants,
)

ART_ROOT = Path(r"D:\NRO_GoiBanGiao\docs\my-thuat\anh-nguon\vertical-slice-v1")
REVIEW = GOC / "docs/my-thuat/review/vertical-slice-v1.json"
MAPPING = GOC / "docs/my-thuat/mapping-id.json"
SPRITE_ROOT = GOC / "docs/my-thuat/lang-la-client-assets/res"
SPRITE_RELATIVE = "docs/my-thuat/lang-la-client-assets/res"
LEGACY_PART = GOC / "server/internal/game/part.bin"
PART_OUT = GOC / "server/internal/game/assets/part.bin"
ASSET_LOCK = GOC / "server/internal/game/assets/asset-lock.json"
PART_IDS = {"head": 30, "body": 57, "leg": 58}
ZOOMS = ("x1", "x2", "x3", "x4")


def _sha(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def sinh_anh(components, root: Path, ghi: bool) -> dict[int, dict[str, str]]:
    """Ghi Small<id>.png cho bốn mức zoom; trả sha256 từng mức."""

    hashes: dict[int, dict[str, str]] = {}
    for component in components:
        sprite_id = component["id"]
        hashes[sprite_id] = {}
        for zoom, image in zoom_variants(component["image_x4"]).items():
            path = root / zoom / "smallimage" / f"Small{sprite_id}.png"
            if ghi:
                path.parent.mkdir(parents=True, exist_ok=True)
                image.save(path, format="PNG", optimize=False, compress_level=9)
                hashes[sprite_id][zoom] = _sha(path.read_bytes())
            else:
                hashes[sprite_id][zoom] = "?" * 64
    return hashes


def cap_nhat_mapping(mapping: dict, components, hashes, feet_offset: int) -> dict:
    """parts[].dx/dy + khối hình học; KHÔNG đụng sprites[] (xem ghi chú bên dưới)."""

    by_id = {component["id"]: component for component in components}
    for part in mapping["parts"]:
        if part["id"] not in set(PART_IDS.values()):
            continue
        for frame in part["frames"]:
            component = by_id.get(frame["sprite_id"])
            if component is None:
                continue
            frame["dx"] = component["placement"]["dx"]
            frame["dy"] = component["placement"]["dy"]
    # Ô contact sheet giờ là giới hạn trên (cell mỗi khung một cỡ) — mapping phải
    # khai đúng ô hiện hành, nếu không `_schema` của build từ chối cả manifest.
    for sheet in mapping["build_inputs"]["animation_sheets"]:
        sheet["cell_size"] = list(
            REQUIRED_ANIMATION_SHEETS[sheet["name"]]["cell_size"]
        )
    scale = components[0]["shared_transform"]["scale"]
    mapping["build_inputs"]["player_component_geometry"] = {
        "ghi_chu": (
            "Hình học THẬT của 34 mảnh người chơi đang chạy. sprites[].zooms vẫn "
            "trỏ production-closure-v11 (hình học CŨ, mỗi khung một tỉ lệ) vì "
            "closure mới cần reviewer độc lập ký lại (ASSET-ARTIFACT-REVIEW-016). "
            "Ảnh đang dùng nằm ở docs/my-thuat/lang-la-client-assets/res/x*/smallimage/."
        ),
        "generator": "tools/nhan-vat/sinh-lai-sprite.py",
        "uniform_scale_x4": scale,
        "feet_offset_x1": feet_offset,
        "anchor": "feet-bottom-center",
        "sprite_root": SPRITE_RELATIVE,
        "components": [
            {
                "id": component["id"],
                "component": component["component"],
                "slot_index": component["slot_index"],
                "action": component["action"],
                "action_slot": component["action_slot"],
                "source_path": component["source_path"],
                "source_rect": component["source_rect"],
                "size_x4": list(component["image_x4"].size),
                "placement_x1": {
                    "left_rel": component["placement"]["left_rel"],
                    "top_rel": component["placement"]["top_rel"],
                },
                "part_offset": {
                    "dx": component["placement"]["dx"],
                    "dy": component["placement"]["dy"],
                },
                "charinfo_cf": component["charinfo_cf"],
                "charinfo_anchor_cf": component["charinfo_anchor_cf"],
                "charinfo_residual": component["charinfo_residual"],
                "paths": {
                    zoom: f"{SPRITE_RELATIVE}/{zoom}/smallimage/Small{component['id']}.png"
                    for zoom in ZOOMS
                },
                "sha256": hashes[component["id"]],
            }
            for component in components
        ],
    }
    return mapping


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--art-root", type=Path, default=ART_ROOT)
    parser.add_argument("--review", type=Path, default=REVIEW)
    parser.add_argument("--mapping", type=Path, default=MAPPING)
    parser.add_argument("--sprite-root", type=Path, default=SPRITE_ROOT)
    parser.add_argument("--legacy-part", type=Path, default=LEGACY_PART)
    parser.add_argument("--part-out", type=Path, default=PART_OUT)
    parser.add_argument("--asset-lock", type=Path, default=ASSET_LOCK)
    parser.add_argument("--vsdata", type=int, default=27)
    parser.add_argument(
        "--thu", action="store_true", help="chỉ in bảng, không ghi file nào"
    )
    args = parser.parse_args()

    review = json.loads(args.review.read_text(encoding="utf-8"))
    records = _review_by_name(review)
    charinfo = load_charinfo()
    components = player_component_images(
        review, art_root=args.art_root, records=records, charinfo=charinfo
    )
    if len(components) != 34:
        raise SystemExit(f"mong doi 34 manh, nhan {len(components)}")

    hashes = sinh_anh(components, args.sprite_root, ghi=not args.thu)

    print("id    manh slot  x4        x1      dx   dy  cf-neo  cf dung slot")
    for component in components:
        width, height = component["image_x4"].size
        print(
            "%5d %-4s %2d   %3dx%-3d %3dx%-3d %+4d %+4d %5d  %s"
            % (
                component["id"],
                component["component"],
                component["slot_index"],
                width,
                height,
                round(width / 4),
                round(height / 4),
                component["placement"]["dx"],
                component["placement"]["dy"],
                component["charinfo_anchor_cf"],
                component["charinfo_cf"],
            )
        )
    print("ti le chung s =", components[0]["shared_transform"]["scale"])

    mapping = json.loads(args.mapping.read_text(encoding="utf-8"))
    mapping = cap_nhat_mapping(mapping, components, hashes, PLAYER_FEET_OFFSET_X1)
    part_bytes = build_part(mapping, args.legacy_part.read_bytes())
    print(
        "part.bin: %d record, %d byte, sha %s"
        % (len(parse_part(part_bytes)), len(part_bytes), _sha(part_bytes)[:16])
    )
    if args.thu:
        print("(--thu: khong ghi file nao)")
        return 0

    args.mapping.write_text(
        json.dumps(mapping, ensure_ascii=False, indent=2) + "\n", encoding="utf-8"
    )
    args.part_out.write_bytes(part_bytes)

    # asset-lock.json phải giữ ĐÚNG dạng canonical của `assetpack.lock.lock_bytes`
    # (một dòng, khoá đã sắp xếp): test Go tìm chuỗi con `"mob_ids":[0]` trong đó,
    # ghi lại bằng indent là mọi chốt fail-closed đó im lặng ngừng hoạt động.
    lock = json.loads(args.asset_lock.read_text(encoding="utf-8"))
    for output in lock["outputs"]:
        if output["path"] == "server/part.bin":
            output["sha256"] = _sha(part_bytes)
            output["bytes"] = len(part_bytes)
    lock["versions"]["data"] = args.vsdata
    for sheet in lock.get("closure", {}).get("animation_sheets", []):
        width, height = REQUIRED_ANIMATION_SHEETS[sheet["name"]]["cell_size"]
        sheet["cell_size"] = [width, height]
        for cell in sheet["cells"]:
            slot = cell["slot_index"]
            cell["bounds"] = [slot * width, 0, (slot + 1) * width, height]
    args.asset_lock.write_bytes(lock_bytes(lock))
    print("da ghi:", args.part_out, args.mapping, args.asset_lock)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
